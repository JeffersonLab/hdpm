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
)

type Package struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	URL        string   `json:"url"`
	Path       string   `json:"path"`
	Cmds       []string `json:"cmds"`
	Deps       []string `json:"deps"`
	IsPrebuilt bool     `json:"isPrebuilt"`
}

// Default package settings
var masterPackages = []Package{
	{Name: "cmake", Version: "3.6.2",
		URL:        "https://cmake.org/files/v3.6/cmake-[VER]-Linux-x86_64.tar.gz",
		Path:       "cmake/[VER]",
		Cmds:       []string{""},
		Deps:       []string{""},
		IsPrebuilt: true},
	{Name: "xerces-c", Version: "3.1.4",
		URL:        "http://archive.apache.org/dist/xerces/c/3/sources/xerces-c-[VER].tar.gz",
		Path:       "xerces-c/[VER]",
		Cmds:       []string{"./configure --prefix=[PATH]", "make", "make install"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "cernlib", Version: "2005",
		URL:        "http://www-zeuthen.desy.de/linear_collider/cernlib/new/cernlib.2005.corr.2014.04.17.tgz",
		Path:       "cernlib",
		Cmds:       []string{""},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "root", Version: "6.06.08",
		URL:        "https://root.cern.ch/download/root_v[VER].source.tar.gz",
		Path:       "root/[VER]",
		Cmds:       []string{"cmake -Droofit=ON -DCMAKE_INSTALL_PREFIX=[PATH] ../root", "cmake --build . -- -j8", "cmake --build . --target install"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "amptools", Version: "0.9.2",
		URL:        "http://downloads.sourceforge.net/project/amptools/AmpTools_v[VER].tgz",
		Path:       "amptools/[VER]",
		Cmds:       []string{"cd AmpTools; make; cd ../", "cd AmpPlotter; make"},
		Deps:       []string{"root"},
		IsPrebuilt: false},
	{Name: "geant4", Version: "10.02.p02",
		URL:  "http://geant4.cern.ch/support/source/geant4.[VER].tar.gz",
		Path: "geant4/[VER]",
		Cmds: []string{"cmake -DCMAKE_INSTALL_PREFIX=[PATH] ../geant4",
			"make -j8", "make install"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "evio", Version: "4.4.6",
		URL:        "https://coda.jlab.org/drupal/system/files/coda/evio/evio-4.4/evio-[VER].tgz",
		Path:       "evio/[VER]",
		Cmds:       []string{"scons --prefix=[PATH] install"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "rcdb", Version: "0.00",
		URL:        "https://github.com/JeffersonLab/rcdb/archive/v[VER].tar.gz",
		Path:       "rcdb/[VER]",
		Cmds:       []string{"cd cpp; scons"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "ccdb", Version: "1.06.01",
		URL:        "https://github.com/JeffersonLab/ccdb/archive/v[VER].tar.gz",
		Path:       "ccdb/[VER]",
		Cmds:       []string{"scons"},
		Deps:       []string{""},
		IsPrebuilt: false},
	{Name: "jana", Version: "0.7.5p2",
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
		Deps:       []string{"xerces-c", "cernlib", "root", "evio", "ccdb", "jana", "hdds"},
		IsPrebuilt: false},
	{Name: "gluex_root_analysis", Version: "master",
		URL:        "https://github.com/JeffersonLab/gluex_root_analysis/archive/[VER].tar.gz",
		Path:       "gluex_root_analysis/[VER]",
		Cmds:       []string{"./make_all.sh"},
		Deps:       []string{"root", "sim-recon"},
		IsPrebuilt: false},
}

var jlabVersions = map[string]string{
	"cmake":               "3.6.2",
	"xerces-c":            "3.1.2",
	"cernlib":             "2005",
	"root":                "5.34.34",
	"amptools":            "",
	"geant4":              "",
	"evio":                "4.3.1",
	"rcdb":                "0.00",
	"ccdb":                "1.06.01",
	"jana":                "0.7.5p2",
	"hdds":                "master",
	"sim-recon":           "master",
	"gluex_root_analysis": "master",
}

var ws16Package = Package{
	Name: "gluex_workshops", Version: "master",
	URL:        "https://github.com/JeffersonLab/gluex_workshops",
	Path:       "gluex_workshops/[VER]",
	Cmds:       []string{"cd physics_workshop_2016/session2/omega_ref; scons install", "cd physics_workshop_2016/session2/omega_skim_rest; scons install", "cd physics_workshop_2016/session2/omega_solutions; scons install", "cd physics_workshop_2016/session3b/omega_skim_tree; scons install", "cd physics_workshop_2016/session5b/p2gamma_workshop; scons install"},
	Deps:       []string{"xerces-c", "cernlib", "root", "evio", "ccdb", "jana", "hdds", "sim-recon", "gluex_root_analysis"},
	IsPrebuilt: false,
}
var ws16Versions = map[string]string{
	"cmake":               "3.6.2",
	"xerces-c":            "3.1.3",
	"cernlib":             "2005",
	"root":                "5.34.36",
	"amptools":            "0.9.1",
	"geant4":              "",
	"evio":                "4.4.6",
	"rcdb":                "0.00",
	"ccdb":                "1.06.01",
	"jana":                "0.7.5p1",
	"hdds":                "physics_workshop_2016",
	"sim-recon":           "physics_workshop_2016",
	"gluex_root_analysis": "physics_workshop_2016",
	"gluex_workshops":     "master",
}

var packages []Package
var OS string

func init() {
	envInit()
	OS = osrelease()
	dir := filepath.Join(packageDir(), "settings")
	for _, pkg := range masterPackages {
		if isPath(dir) {
			pkg.read()
		}
		pkg.config("master")
		packages = append(packages, pkg)
	}
}

func getPackage(name string) Package {
	for _, p := range packages {
		if name == p.Name {
			return p
		}
	}
	return Package{}
}

var packageNames = []string{
	"cmake",
	"xerces-c",
	"cernlib",
	"root",
	"amptools",
	"geant4",
	"evio",
	"rcdb",
	"ccdb",
	"jana",
	"hdds",
	"sim-recon",
	"gluex_root_analysis",
	"gluex_workshops",
}

func packageDir() string {
	pdir := os.Getenv("GLUEX_TOP")
	if pdir == "" {
		pdir, _ = os.Getwd()
	}
	return pdir
}

func jlabPackageDir() string {
	return "/group/halld/Software/builds/" + OS
}

var jsep = map[string]string{
	"cmake":               "-",
	"xerces-c":            "-",
	"cernlib":             "",
	"root":                "_",
	"amptools":            "_",
	"geant4":              "-",
	"evio":                "-",
	"rcdb":                "_",
	"ccdb":                "_",
	"jana":                "_",
	"hdds":                "-",
	"sim-recon":           "-",
	"gluex_root_analysis": "-",
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

func (p *Package) config(arg string) {
	if p.Version == "" {
		p.URL = ""
		p.Path = ""
	}
	if p.Version == "latest" {
		p.Version = "master"
	}
	p.Path = strings.Replace(p.Path, "[OS]", OS, -1)
	p.Path = strings.Replace(p.Path, "[VER]", p.Version, -1)
	if !strings.HasPrefix(p.Path, "/") && p.Path != "" {
		p.Path = filepath.Join(packageDir(), p.Path)
	}
	if p.Name == "evio" {
		major_minor := ver_i(p.Version, 0) + "." + ver_i(p.Version, 1)
		re := regexp.MustCompile("4.[0-9]")
		if !strings.Contains(p.URL, major_minor) {
			p.URL = re.ReplaceAllString(p.URL, major_minor)
		}
	}
	p.URL = strings.Replace(p.URL, "[VER]", p.Version, -1)

	// Template-dependent config.
	if arg == "jlab" {
		p.Version = jlabVersions[p.Name]
		jp := filepath.Join(jlabPackageDir(), p.Name, p.Name+jsep[p.Name]+p.Version)
		if isPath(jp) {
			p.Path = jp
			p.IsPrebuilt = true
		}
		if p.Name == "cernlib" && isPath(filepath.Join(jlabPackageDir(), p.Name)) {
			p.Path = filepath.Join(jlabPackageDir(), p.Name)
			p.IsPrebuilt = true
		}
	}
	if arg == "workshop-2016" {
		p.Version = ws16Versions[p.Name]
	}

	if p.Name == "rcdb" && strings.Contains(OS, "gcc4.4") {
		p.IsPrebuilt = true
	}
	if p.Version == "master" && strings.Contains(p.URL, "https://github.com/JeffersonLab/"+p.Name+"/archive/") {
		p.URL = "https://github.com/JeffersonLab/" + p.Name
	}
	if p.Name == "jana" && p.Version == "master" {
		p.URL = "https://phys12svn.jlab.org/repos/JANA"
	}
	var cmds []string
	for _, cmd := range p.Cmds {
		if p.Path != "" {
			cmds = append(cmds, strings.Replace(cmd, "[PATH]", p.Path, -1))
		}
	}
	p.Cmds = cmds
	if p.Name == "root" && ver_i(p.Version, 0) == "6" && strings.Contains(os.Getenv("PATH"), "/opt/rh/devtoolset-3/root/usr/bin") && (strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6")) {
		if len(p.Cmds) > 0 && !strings.Contains(p.Cmds[0], "./configure") {
			p.Cmds = nil
			p.Cmds = append(p.Cmds, "./configure --enable-roofit")
			p.Cmds = append(p.Cmds, "make -j8 && make clean")
		}
	}
	if runtime.GOOS == "darwin" {
		path := filepath.Join("/usr/local/Cellar/xerces-c", p.Version)
		if p.Name == "xerces-c" && isPath(path) && isPath(packageDir()+"/dist") {
			p.IsPrebuilt = true
			p.Path = path
		}
	}
}

func (p *Package) template() {
	p.URL = strings.Replace(p.URL, p.Version, "[VER]", -1)
	var cmds []string
	for _, cmd := range p.Cmds {
		cmds = append(cmds, strings.Replace(cmd, p.Path, "[PATH]", -1))
	}
	p.Cmds = cmds
	p.Path = strings.Replace(p.Path, packageDir(), "", -1)
	p.Path = strings.Replace(p.Path, "/"+p.Name, p.Name, -1)
	p.Path = strings.Replace(p.Path, p.Version, "[VER]", -1)
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
	p.template()
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

func (p *Package) read() {
	dir := filepath.Join(packageDir(), "settings")
	b, err := ioutil.ReadFile(dir + "/" + p.Name + ".json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(b, &p)
	if !in(packageNames, p.Name) {
		fmt.Fprintf(os.Stderr, "%s: unknown package name\n", p.Name)
		os.Exit(2)
	}
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
	dir := filepath.Join(packageDir(), "settings")
	for _, p1 := range packages {
		for _, p2 := range v.Packs {
			if p1.Name == p2.Name {
				p1.Version = p2.Version
			}
		}
		p1.write(dir)
	}
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

func isPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func isSymLink(path string) bool {
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
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Server release 6.*") {
				release = "_RHEL6"
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Workstation release 7.") {
				release = "_RHEL7"
			} else if strings.HasPrefix(rs, "Red Hat Enterprise Linux Server release 7.*") {
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
			s := make([]string, 2)
			for _, l := range strings.Split(rs, "\n") {
				if strings.HasPrefix(l, "DISTRIB_ID=") {
					s = append(s, strings.TrimLeft(l, "DISTRIB_ID="))
				}
				if strings.HasPrefix(l, "DISTRIB_RELEASE=") {
					s = append(s, strings.TrimLeft(l, "DISTRIB_RELEASE="))
				}
			}
			release = "_" + s[0] + s[1]
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
