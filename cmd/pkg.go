package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Settings struct {
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
}

func newSettings(name, comment string) *Settings {
	s := &Settings{}
	s.Name = name
	s.Comment = comment
	t := time.Now().Round(time.Second)
	s.Timestamp = t.Format(time.RFC3339)
	return s
}

func (s *Settings) read(dir string) {
	b, err := ioutil.ReadFile(dir + "/.info.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(b, &s)
}

func (s *Settings) write(dir string) {
	f, err := os.Create(dir + "/.info.json")
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(f, "%s\n", data)
	f.Close()
}

type Package struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	URL        string   `json:"url"`
	Path       string   `json:"path"`
	Cmds       []string `json:"cmds,omitempty"`
	Deps       []string `json:"deps,omitempty"`
	IsPrebuilt bool     `json:"isPrebuilt,omitempty"`
}

// Default package settings
var masterPackages = [...]Package{
	{Name: "hdpm", Version: "latest",
		URL:        "https://halldweb.jlab.org/dist/hdpm/hdpm-[VER].linux.tar.gz",
		Path:       "hdpm/[VER]",
		Cmds:       nil,
		Deps:       nil,
		IsPrebuilt: true},
	{Name: "cmake", Version: "3.7.2",
		URL:        "https://cmake.org/files/v3.7/cmake-[VER]-Linux-x86_64.tar.gz",
		Path:       "cmake/[VER]",
		Cmds:       nil,
		Deps:       nil,
		IsPrebuilt: true},
	{Name: "xerces-c", Version: "3.1.4",
		URL:        "http://archive.apache.org/dist/xerces/c/3/sources/xerces-c-[VER].tar.gz",
		Path:       "xerces-c/[VER]",
		Cmds:       []string{"./configure --prefix=[PATH]", "make", "make install"},
		Deps:       nil,
		IsPrebuilt: false},
	{Name: "cernlib", Version: "2005",
		URL:        "http://www-zeuthen.desy.de/linear_collider/cernlib/new/cernlib.2005.corr.2014.04.17.tgz",
		Path:       "cernlib",
		Cmds:       nil,
		Deps:       nil,
		IsPrebuilt: false},
	{Name: "root", Version: "6.08.04",
		URL:  "https://root.cern.ch/download/root_v[VER].source.tar.gz",
		Path: "root/[VER]",
		Cmds: []string{"cmake -Droofit=ON -DCMAKE_INSTALL_PREFIX=[PATH] ../src", "cmake --build . -- -j8",
			"cmake --build . --target install", "cd ..; rm -rf build src"},
		Deps:       []string{"cmake"},
		IsPrebuilt: false},
	{Name: "amptools", Version: "0.9.2",
		URL:        "http://downloads.sourceforge.net/project/amptools/AmpTools_v[VER].tgz",
		Path:       "amptools/[VER]",
		Cmds:       []string{"cd AmpTools; make", "cd AmpPlotter; make"},
		Deps:       []string{"root"},
		IsPrebuilt: false},
	{Name: "geant4", Version: "10.02.p02",
		URL:  "http://geant4.cern.ch/support/source/geant4.[VER].tar.gz",
		Path: "geant4/[VER]",
		Cmds: []string{"cmake -DCMAKE_INSTALL_PREFIX=[PATH] -DXERCESC_ROOT_DIR=${XERCESCROOT} -DGEANT4_USE_RAYTRACER_X11=ON -DGEANT4_USE_OPENGL_X11=ON -DGEANT4_BUILD_MULTITHREADED=ON -DGEANT4_INSTALL_DATA=ON ../src",
			"make -j8", "make install", "cd ..; rm -rf build src"},
		Deps:       []string{"cmake", "xerces-c"},
		IsPrebuilt: false},
	{Name: "evio", Version: "4.4.6",
		URL:        "https://coda.jlab.org/drupal/system/files/coda/evio/evio-4.4/evio-[VER].tgz",
		Path:       "evio/[VER]",
		Cmds:       []string{"scons --prefix=[PATH] install"},
		Deps:       nil,
		IsPrebuilt: false},
	{Name: "rcdb", Version: "0.01",
		URL:        "https://github.com/JeffersonLab/rcdb/archive/v[VER].tar.gz",
		Path:       "rcdb/[VER]",
		Cmds:       []string{"cd cpp; scons"},
		Deps:       nil,
		IsPrebuilt: false},
	{Name: "ccdb", Version: "1.06.02",
		URL:        "https://github.com/JeffersonLab/ccdb/archive/v[VER].tar.gz",
		Path:       "ccdb/[VER]",
		Cmds:       []string{"scons"},
		Deps:       nil,
		IsPrebuilt: false},
	{Name: "jana", Version: "0.7.7p1",
		URL:        "https://www.jlab.org/JANA/releases/jana_[VER].tgz",
		Path:       "jana/[VER]",
		Cmds:       []string{"scons -u -j8 install"},
		Deps:       []string{"xerces-c", "root", "ccdb"},
		IsPrebuilt: false},
	{Name: "hdds", Version: "master",
		URL:        "https://github.com/JeffersonLab/hdds/archive/[VER].tar.gz",
		Path:       "hdds/[VER]",
		Cmds:       []string{"scons -u install"},
		Deps:       []string{"xerces-c"},
		IsPrebuilt: false},
	{Name: "sim-recon", Version: "master",
		URL:        "https://github.com/JeffersonLab/sim-recon/archive/[VER].tar.gz",
		Path:       "sim-recon/[VER]",
		Cmds:       []string{"scons -u -j8 install DEBUG=0"},
		Deps:       []string{"cernlib", "evio", "rcdb", "jana", "hdds"},
		IsPrebuilt: false},
	{Name: "hdgeant4", Version: "master",
		URL:        "https://github.com/JeffersonLab/HDGeant4",
		Path:       "hdgeant4/[VER]",
		Cmds:       []string{"ln -sfn G4.10.02.p02fixes src/G4fixes", ". ${G4ROOT}/share/Geant4-10.2.2/geant4make/geant4make.sh; make"},
		Deps:       []string{"geant4", "sim-recon"},
		IsPrebuilt: false},
	{Name: "gluex_root_analysis", Version: "master",
		URL:        "https://github.com/JeffersonLab/gluex_root_analysis/archive/[VER].tar.gz",
		Path:       "gluex_root_analysis/[VER]",
		Cmds:       []string{"./make_all.sh"},
		Deps:       []string{"sim-recon"},
		IsPrebuilt: false},
	{Name: "gluex_workshops", Version: "master",
		URL:        "https://github.com/JeffersonLab/gluex_workshops",
		Path:       "gluex_workshops/[VER]",
		Cmds:       []string{"cd physics_workshop_2016/session2/omega_ref; scons install", "cd physics_workshop_2016/session2/omega_skim_rest; scons install", "cd physics_workshop_2016/session2/omega_solutions; scons install", "cd physics_workshop_2016/session3b/omega_skim_tree; scons install", "cd physics_workshop_2016/session5b/p2gamma_workshop; scons install"},
		Deps:       []string{"gluex_root_analysis"},
		IsPrebuilt: false},
}

// Packages to use
var packages []Package
var packageNames []string

// OS release
var OS string

// Package directory
var PD string

// JLab package directory
const JPD = "/group/halld/Software/builds"

// Settings directory
var SD string

func pkgInit() {
	PD = os.Getenv("GLUEX_TOP")
	if PD == "" {
		PD, _ = os.Getwd()
	}

	OS = osrelease()
	if strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6") {
		if strings.Contains(OS, "gcc") {
			setenvGCC()
		}
	}

	SD = filepath.Join(PD, "settings")

	rsd := isPath(SD)
	for _, pkg := range masterPackages {
		if rsd {
			if isPath(SD + "/" + pkg.Name + ".json") {
				pkg = read(pkg.Name)
				packages = append(packages, pkg)
				packageNames = append(packageNames, pkg.Name)
			}
		} else {
			packages = append(packages, pkg)
			packageNames = append(packageNames, pkg.Name)
		}
	}
}

func read(name string) (p Package) {
	b, err := ioutil.ReadFile(SD + "/" + name + ".json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(b, &p)
	p.Name = name
	return p
}

func getPackage(name string) Package {
	for _, p := range packages {
		if name == p.Name {
			p.config()
			return p
		}
	}
	return Package{}
}

var jsep = map[string]string{
	"hdpm":                "-",
	"cmake":               "-",
	"xerces-c":            "-",
	"cernlib":             "",
	"root":                "-",
	"amptools":            "_",
	"geant4":              ".",
	"evio":                "-",
	"rcdb":                "_",
	"ccdb":                "_",
	"jana":                "_",
	"hdds":                "-",
	"sim-recon":           "-",
	"hdgeant4":            "-",
	"gluex_root_analysis": "-",
	"gluex_workshops":     "-",
}

func ver_i(ver string, i int) string {
	if i < 0 || i > 2 {
		return ver
	}
	for _, v := range strings.Split(ver, "-") {
		if strings.Contains(v, ".") {
			return strings.Split(v, ".")[i]
		}
	}
	return ver
}
func (p *Package) configBinary() {
	s := ""
	if strings.Contains(OS, "Darwin_macosx10.12") {
		s = "macosx64-10.12-clang80"
	}
	if strings.Contains(OS, "Fedora24") {
		s = "Linux-fedora24-x86_64-gcc6.1"
	}
	if strings.Contains(OS, "CentOS7") || strings.Contains(OS, "RHEL7") {
		s = "Linux-centos7-x86_64-gcc4.8"
	}
	if strings.Contains(OS, "Ubuntu14") || strings.Contains(OS, "LinuxMint17") {
		s = "Linux-ubuntu14-x86_64-gcc4.8"
	}
	if strings.Contains(OS, "Ubuntu16") || strings.Contains(OS, "LinuxMint18") {
		s = "Linux-ubuntu16-x86_64-gcc5.4"
	}
	if s == "" {
		fmt.Fprintf(os.Stderr, "%s: ROOT binary distribution not available\n", OS)
		return
	}
	p.URL = "https://root.cern.ch/download/root_v" + p.Version + "." + s + ".tar.gz"
	p.IsPrebuilt = true
	p.Cmds, p.Deps = nil, nil
}

func (p *Package) config() {
	if p.Version == "" {
		p.URL = ""
		p.Path = ""
		p.Cmds, p.Deps = nil, nil
		p.IsPrebuilt = true
		return
	}

	if p.Name == "evio" || p.Name == "cmake" {
		p.configMajorMinorInURL()
	}
	if p.Name == "hdpm" && runtime.GOOS == "darwin" {
		p.URL = strings.Replace(p.URL, "linux", "macOS", -1)
	}
	p.URL = strings.Replace(p.URL, "[VER]", p.Version, -1)

	p.Path = strings.Replace(p.Path, "[OS]", OS, -1)
	p.Path = strings.Replace(p.Path, "[VER]", p.Version, -1)
	if !strings.HasPrefix(p.Path, "/") && p.Path != "" {
		p.Path = filepath.Join(PD, p.Path)
	}

	if p.Version == "master" && strings.Contains(p.URL, "https://github.com/JeffersonLab/"+p.Name+"/archive/") {
		p.URL = "https://github.com/JeffersonLab/" + p.Name
	}
	if p.Name == "jana" && p.Version == "master" {
		p.URL = "https://phys12svn.jlab.org/repos/JANA"
	}

	p.configCmds("[PATH]", p.Path)

	if p.Name == "root" && ver_i(p.Version, 0) == "6" && strings.Contains(os.Getenv("PATH"), "/opt/rh/devtoolset-3/root/usr/bin") && (strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6")) {
		if len(p.Cmds) > 0 && !strings.Contains(p.Cmds[0], "./configure") {
			p.Cmds = nil
			p.Cmds = append(p.Cmds, "./configure --enable-roofit")
			p.Cmds = append(p.Cmds, "make -j8; make clean")
		}
	}

	if p.inDist() {
		p.IsPrebuilt = true
	}
}

func (p *Package) configMajorMinorInURL() {
	major := ver_i(p.Version, 0)
	major_minor := major + "." + ver_i(p.Version, 1)
	re := regexp.MustCompile(major + ".[0-9]")
	if !strings.Contains(p.URL, major_minor) {
		p.URL = re.ReplaceAllString(p.URL, major_minor)
	}
}

func (p *Package) configCmds(oldPath, newPath string) {
	var cmds []string
	for _, cmd := range p.Cmds {
		if p.Path != "" {
			cmds = append(cmds, strings.Replace(cmd, oldPath, newPath, -1))
		}
	}
	p.Cmds = cmds
}

func (p *Package) template() {
	if p.Version == "" {
		return
	}
	p.URL = strings.Replace(p.URL, p.Version, "[VER]", -1)
	p.configCmds(p.Path, "[PATH]")
	p.Path = strings.Replace(p.Path, PD+"/", "", -1)
	p.Path = strings.Replace(p.Path, p.Version, "[VER]", -1)
	p.Path = strings.Replace(p.Path, OS, "[OS]", -1)
}

func write_text(fname, text string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(f, text)
	f.Close()
}

func (p *Package) write(dir string) {
	f, err := os.Create(dir + "/" + p.Name + ".json")
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(f, "%s\n", data)
	f.Close()
}

func exitUnknownPackage(arg string) {
	fmt.Fprintf(os.Stderr, "%s: Unknown package name\n", arg)
	os.Exit(2)
}

func exitNoPackages(cmd *cobra.Command) {
	fmt.Fprintln(os.Stderr, "No packages were specified on the command line.\n")
	cmd.Usage()
	os.Exit(2)
}

func printPackages(args []string) {
	if len(args) > 0 {
		var a []string
		for _, n := range packageNames {
			if in(args, n) {
				a = append(a, n)
			}
		}
		fmt.Printf("Packages: %s\n", strings.Join(a, ", "))
	}
}

func extractNames(args []string) []string {
	var names []string
	for _, arg := range args {
		if strings.Contains(arg, "@") {
			names = append(names, strings.Split(arg, "@")[0])
		} else {
			names = append(names, arg)
		}
	}
	return names
}

func extractVersions(args []string) map[string]string {
	versions := make(map[string]string)
	unchanged := true
	for _, arg := range args {
		if strings.Contains(arg, "@") {
			versions[strings.Split(arg, "@")[0]] = strings.Split(arg, "@")[1]
			unchanged = false
		} else {
			versions[arg] = ""
		}
	}
	if unchanged {
		versions = nil
	}
	return versions
}

func changeVersions(names []string, versions map[string]string) {
	if len(versions) == 0 {
		return
	}
	mk(SD)
	s := newSettings("master", "Default settings of hdpm version "+VERSION)
	if !isPath(SD + "/.info.json") {
		s.write(SD)
	}
	var pkgs []Package
	for _, pkg := range packages {
		if !pkg.in(names) {
			pkgs = append(pkgs, pkg)
			pkg.write(SD)
			continue
		}
		ver, ok := versions[pkg.Name]
		if pkg.Name == "root" && ver == "binary" {
			pkg.configBinary()
		} else {
			pkg.changeVersion(ver, ok)
		}
		pkgs = append(pkgs, pkg)
		pkg.write(SD)
	}
	packages = pkgs
}

func (p *Package) changeVersion(ver string, ok bool) {
	if ok && ver != "" {
		p.Version = ver
		if strings.HasPrefix(p.Path, JPD) {
			p.Path = filepath.Join(p.Name, "[VER]")
			p.IsPrebuilt = false
		}
	}
}

func extractVersion(arg string) string {
	ver := ""
	if strings.Contains(arg, "@") {
		ver = strings.Split(arg, "@")[1]
	}
	return ver
}

func (p *Package) jlabPathConfig(dirtag string) {
	dir := filepath.Join(JPD, OS, p.Name)
	jp := filepath.Join(dir, p.Name+jsep[p.Name]+p.Version)
	if dirtag != "" {
		jp += "^" + dirtag
	}
	if isPath(jp) {
		p.Path = jp
		p.IsPrebuilt = true
	}
	if p.Name == "cernlib" && isPath(dir+"/"+p.Version) {
		p.Path = dir
		p.IsPrebuilt = true
	}
	p.template()
}

func versionXML(file string) {
	msg :=
		`Version XMLfile Directory
URL:  https://halldweb.jlab.org/dist
Path: /group/halld/www/halldweb/html/dist
`
	fmt.Print(msg)
	if file == "latest" {
		file = "https://halldweb.jlab.org/dist/version.xml"
	}
	jlab := strings.HasPrefix(file, "jlab")
	jdev := strings.HasPrefix(file, "jlab-dev")
	if jlab || jdev {
		ver := extractVersion(file)
		if ver != "" {
			file = "https://halldweb.jlab.org/dist/version_" + ver + "_jlab.xml"
		} else {
			file = "https://halldweb.jlab.org/dist/version_jlab.xml"
		}
	}
	wasurl := false
	if strings.Contains(file, "https://") || strings.Contains(file, "http://") {
		wasurl = true
		fmt.Printf("\nDownloading %s ...\n", file)
		run("curl", "-OL", file)
		file = filepath.Base(file)
	}
	type pack struct {
		Name       string `xml:"name,attr"`
		Version    string `xml:"version,attr"`
		WordLength string `xml:"word_length,attr"`
		DirTag     string `xml:"dirtag,attr"`
		URL        string `xml:"url,attr"`
		Branch     string `xml:"branch,attr"`
	}
	type vXML struct {
		XMLName xml.Name `xml:"gversions"`
		Packs   []pack   `xml:"package"`
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	var v *vXML
	xml.Unmarshal(b, &v)
	mk(SD)
	id, c := "master", "Version XML"
	if jlab && !jdev {
		id, c = "jlab", "JLab CUE"
	}
	if jlab && jdev {
		id, c = "jlab-dev", "JLab-dev CUE"
	}
	s := newSettings(id, c)
	if isPath(SD + "/.info.json") {
		s.read(SD)
	}
	s.write(SD)
	var pkgs []Package
	for _, p1 := range packages {
		for _, p2 := range v.Packs {
			if p1.Name == p2.Name {
				if p2.Version != "" {
					p1.Version = p2.Version
				}
				if p2.Version == "" && p2.Branch != "" {
					p1.Version = p2.Branch
				}
				if p2.URL != "" {
					p1.URL = p2.URL
				}
				if jlab || jdev {
					if jdev && (p1.Name == "hdds" || p1.Name == "sim-recon") {
						p1.Version = "master"
						if strings.HasPrefix(p1.Path, JPD) {
							p1.Path = filepath.Join(p1.Name, "[VER]")
						}
						p1.IsPrebuilt = false
					} else {
						p1.jlabPathConfig(p2.DirTag)
					}
				} else {
					if p2.DirTag != "" {
						p1.Path += "^" + p2.DirTag
					}
				}
			}
		}
		if p1.Name == "cmake" && isPath("/apps/cmake/cmake-3.5.1") {
			p1.Version = "3.5.1"
			p1.Path = "/apps/cmake/cmake-[VER]"
		}
		pkgs = append(pkgs, p1)
		p1.write(SD)
	}
	packages = pkgs
	fmt.Println("\nThe XMLfile versions have been applied to your current settings.")
	if wasurl {
		os.Remove(file)
	}
}

func readFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.TrimRight(string(b), "\n")
}

func in(args []string, item string) bool {
	for _, arg := range args {
		if item == arg {
			return true
		}
	}
	return false
}

func (p *Package) in(args []string) bool {
	return in(args, p.Name)
}

func mk(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalln(err)
	}
}

func cd(path string) {
	if err := os.Chdir(path); err != nil {
		log.Fatalln(err)
	}
}

func glob(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln(err)
	}
	return files
}

func rmGlob(pattern string) {
	files := glob(pattern)
	for _, file := range files {
		os.RemoveAll(file)
	}
}

func mkcd(path string) {
	mk(path)
	cd(path)
}

func (p *Package) mkcd() {
	mkcd(p.Path)
}

func (p *Package) cd() {
	cd(p.Path)
}

func (p *Package) inDist() bool {
	if isSymLink(p.Path) {
		s := readLink(p.Path)
		if strings.HasPrefix(s, "../.dist/") || strings.HasPrefix(s, ".dist/") {
			return true
		}
	}
	return false
}

func isPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func isSymLink(path string) bool {
	if !isPath(path) {
		return false
	}
	stat, err := os.Lstat(path)
	if err != nil {
		log.Fatalln(err)
	}
	return stat.Mode()&os.ModeSymlink == os.ModeSymlink
}

func readLink(path string) string {
	link, err := os.Readlink(path)
	if err != nil {
		log.Fatalln(err)
	}
	return link
}

func relPath(base, target string) string {
	p, err := filepath.Rel(base, target)
	if err != nil {
		log.Fatalln(err)
	}
	return p
}

func readDir(path string) []string {
	if !isPath(path) {
		return nil
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}

func contains(list []string, sublist []string) bool {
	found := 0
	for _, sl := range sublist {
		for _, l := range list {
			if sl == l {
				found++
			}
		}
	}
	return found == len(sublist)
}

func setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Fatalln(err)
	}
}

func unsetenv(key string) {
	if err := os.Unsetenv(key); err != nil {
		log.Fatalln(err)
	}
}

func run(name string, args ...string) {
	if err := command(name, args...).Run(); err != nil {
		log.Fatalln(err)
	}
}

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func commande(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

func output(name string, args ...string) string {
	c := exec.Command(name, args...)
	b, err := c.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}
	return strings.TrimRight(string(b), "\n")
}

func osrelease() string {
	uname := output("uname")
	release := ""
	switch uname {
	case "Linux":
		if isPath("/etc/fedora-release") {
			rs := readFile("/etc/fedora-release")
			if strings.HasPrefix(rs, "Fedora release") {
				release = "_Fedora" + strings.Split(rs, " ")[2]
			} else {
				fmt.Fprintln(os.Stderr, "unrecognized Fedora release")
				release = "_Fedora"
			}
		} else if isPath("/etc/redhat-release") {
			rs := readFile("/etc/redhat-release")
			if strings.HasPrefix(rs, "Red Hat Enterprise Linux Workstation release 6.") {
				release = "_RHEL6"
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Server release 6.") {
				release = "_RHEL6"
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Workstation release 7.") {
				release = "_RHEL7"
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Server release 7.") {
				release = "_RHEL7"
			} else if strings.HasPrefix(rs, "CentOS release 6.") {
				release = "_CentOS6"
			} else if strings.HasPrefix(rs, "CentOS Linux release 7.") {
				release = "_CentOS7"
			} else if strings.HasPrefix(rs, "Scientific Linux release 6.") {
				release = "_SL6"
			} else {
				fmt.Fprintln(os.Stderr, "unrecognized Red Hat release")
				release = "_RH"
			}
		} else if isPath("/etc/lsb-release") {
			rs := readFile("/etc/lsb-release")
			release = "_"
			for _, l := range strings.Split(rs, "\n") {
				if strings.HasPrefix(l, "DISTRIB_ID=") {
					release += strings.TrimPrefix(l, "DISTRIB_ID=")
				}
				if strings.HasPrefix(l, "DISTRIB_RELEASE=") {
					release += strings.TrimPrefix(l, "DISTRIB_RELEASE=")
				}
			}
		} else {
			fmt.Fprintln(os.Stderr, "unrecognized Linux release")
			release = "_Linux"
		}
	case "Darwin":
		rs := output("sw_vers", "-productVersion")
		release = "_macosx10." + strings.Split(rs, ".")[1]
	}
	ccv := output("cc", "-dumpversion")
	ccverbose := output("cc", "-v")
	cct := "cc"
	if strings.Contains(ccverbose, "gcc version") {
		cct = "gcc"
		ccv = output("gcc", "-dumpversion")
	} else if strings.Contains(ccverbose, "clang version") {
		cct = "clang"
		for _, l := range strings.Split(ccverbose, "\n") {
			if strings.HasPrefix(l, "clang version") {
				ccv = strings.Split(l, " ")[2]
			}
		}
	} else if strings.Contains(ccverbose, "Apple LLVM version") {
		cct = "llvm"
		for _, l := range strings.Split(ccverbose, "\n") {
			if strings.HasPrefix(l, "Apple LLVM version") {
				ccv = strings.Split(l, " ")[3]
			}
		}
	}
	proc := output("uname", "-p")
	if strings.Contains(ccverbose, "Target: x86_64") || strings.Contains(ccverbose, "Target: i686-apple-darwin") {
		proc = "x86_64"
	}
	if proc == "unknown" {
		proc = output("uname", "-m")
	}
	return uname + release + "-" + proc + "-" + cct + ccv
}
