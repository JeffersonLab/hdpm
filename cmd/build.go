package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Create the build command
var cmdBuild = &cobra.Command{
	Use:   "build [PACKAGE...]",
	Short: "Build packages and dependencies",
	Long:  `Build packages and dependencies.`,
	Example: `1. hdpm build sim-recon
2. hdpm build geant4 amptools xerces-c
3. hdpm build --all
4. hdpm build sim-recon --xml https://halldweb.jlab.org/dist/version.xml
5. hdpm build gluex_root_analysis -i

Usage:
  hdpm build DIRECTORY
Example:
  hdpm build sim-recon/master/src/plugins/Analysis/pi0omega`,
	Run: runBuild,
}

var jobs string
var showInfo bool

func init() {
	cmdHDPM.AddCommand(cmdBuild)

	cmdBuild.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
	cmdBuild.Flags().BoolVarP(&all, "all", "a", false, "Build all packages in the package settings")
	cmdBuild.Flags().StringVarP(&jobs, "jobs", "j", "", "Number of jobs to run in parallel")
	cmdBuild.Flags().BoolVarP(&showInfo, "info", "i", false, "Show current build information and exit")
}

func runBuild(cmd *cobra.Command, args []string) {
	pkgInit()
	// Build a sim-recon subdirectory if passed as argument
	cwd, _ := os.Getwd()
	if len(args) == 1 && (isPath(filepath.Join(cwd, args[0])) || isPath(args[0])) {
		dir := filepath.Join(cwd, args[0])
		if filepath.IsAbs(args[0]) {
			dir = args[0]
		}
		if isPath(filepath.Join(dir, "SConstruct")) || isPath(filepath.Join(dir, "SConscript")) {
			cd(dir)
			setEnv()
			if jobs == "" {
				run("scons", "-u", "install")
			} else {
				run("scons", "-u", "-j"+jobs, "install")
			}
			return
		}
	}
	// Parse args
	versions := extractVersions(args)
	args = extractNames(args)
	for _, arg := range args {
		if !in(packageNames, arg) {
			exitUnknownPackage(arg)
		}
	}
	if len(args) == 0 && !all {
		exitNoPackages(cmd)
	}
	if all {
		args = packageNames
	}
	args = addDeps(args)
	fmt.Printf("Packages: %s\n", strings.Join(args, ", "))

	// Change package versions to XMLfile versions
	if XML != "" {
		versionXML(XML)
	}
	// Change package versions to versions passed on command line
	changeVersions(args, versions)

	writeVersionXML()

	// Set environment variables
	setEnv()

	// Fetch and build packages
	mkcd(PD)
	isBuilt := false
	for _, arg := range args {
		for _, pkg := range packages {
			if pkg.Name != arg {
				continue
			}
			if runtime.GOOS == "darwin" && pkg.Name == "cernlib" {
				fmt.Printf("macOS detected: Skipping %s\n", pkg.Name)
				continue
			}
			pkg.config()
			pkg.fetch()
			if pkg.IsPrebuilt {
				fmt.Printf("Prebuilt package: %s\n", pkg.Name)
				continue
			}
			pkg.build(&isBuilt)
		}
	}
}

func (p *Package) build(isBuilt *bool) {
	ti := time.Now().Round(time.Second)
	fname := p.Path + "/success.hdpm"
	if p.in([]string{"jana", "hdds", "sim-recon"}) {
		fname = p.Path + "/" + OS + "/success.hdpm"
	}
	if isPath(fname) && (showInfo || !p.isRepo()) {
		printInfo(fname, isBuilt)
		return
	}
	fmt.Printf("\n%s: Checking dependencies ...\n", p.Name)
	p.checkDeps()
	fmt.Printf("Building %s-%s ...\n", p.Name, p.gitVersion())
	p.cd()
	du_i := strings.Fields(output("du", "-sh", p.Path))[0]
	if p.Name != "cernlib" {
		if p.Name == "sim-recon" {
			cd("src")
		}
		if p.Name == "hdgeant4" {
			setEnv() // For geant4 env variables
		}
		if p.usesCMake() {
			p.configCMake()
			if p.Name == "geant4" && isPath("src/LICENSE") {
				run("cp", "-p", "src/LICENSE", ".")
			}
			mkcd("build")
		}
		numCPU := runtime.NumCPU()
		for _, cmd := range p.Cmds {
			cmd = applyNumCPU(cmd, numCPU)
			run("sh", "-c", cmd)
		}
		p.cd()
	} else {
		prep_cernlib_patches()
		run("sh", "-c", "patch < Install_cernlib.patch")
		run("./Install_cernlib")
		run("rm", "-rf", "2005/build", "2005/src")
		run("mv", "2005", "..")
		cd("..")
		run("rm", "-rf", "cernlib")
		p.mkcd()
		run("mv", "../2005", ".")
	}
	du_f := strings.Fields(output("du", "-sh", p.Path))[0]
	tf := time.Now().Round(time.Second)
	text := p.Name + "-" + p.gitVersion() + "\n" +
		tf.Format(time.RFC3339) + "\n" +
		"# build duration\n" +
		tf.Sub(ti).String() + "\n" +
		"# disk use, final minus initial\n" +
		du_f + "B - " + du_i + "B\n" +
		"# dependencies\n" +
		p.taggedDeps()
	write_text(fname, text)
	fmt.Println()
}

func printInfo(fname string, isBuilt *bool) {
	if !(*isBuilt) {
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-18s%-18s%-18s%-18s\n", "package", "build time", "disk use", "timestamp")
		fmt.Println(strings.Repeat("-", 80))
	}
	text := readFile(fname)
	d := strings.Split(text, "\n")
	if strings.HasPrefix(d[0], "gluex_root_analysis") {
		d[0] = strings.Replace(d[0], "gluex_root_analysis", "gxRootAna", -1)
	}
	fmt.Printf("%-18s%-18s%-18s%-18s\n", d[0], d[3], d[5], d[1])
	*isBuilt = true
}

func (p *Package) usesCMake() bool {
	for _, cmd := range p.Cmds {
		if strings.HasPrefix(cmd, "cmake") {
			return true
		}
	}
	return false
}

func (p *Package) configCMake() {
	ver := getCMakeVersion("cmake")
	if strings.HasPrefix(ver, "3.") {
		return
	} else {
		ver = getCMakeVersion("cmake3")
	}
	if strings.HasPrefix(ver, "3.") {
		for _, cmd := range p.Cmds {
			if strings.HasPrefix(cmd, "cmake3") {
				return
			}
		}
		p.configCmds("cmake", "cmake3")
	}
}

func getCMakeVersion(name string) string {
	ver := outputnf(name, "--version")
	if ver != "" {
		f := strings.Fields(ver)
		if len(f) >= 3 {
			ver = f[2]
		}
	}
	return ver
}

func applyNumCPU(cmd string, numCPU int) string {
	if numCPU >= 8 && jobs == "" {
		return cmd
	}
	if strings.Contains(cmd, " -j") {
		j := `-j\s??[1-9][0-9]??`
		re := regexp.MustCompile(j)
		s := strconv.Itoa(numCPU)
		if jobs != "" {
			s = jobs
		}
		cmd = re.ReplaceAllString(cmd, "-j"+s)
	}
	return cmd
}

func (p *Package) gitVersion() string {
	if !isPath(p.Path + "/.git") {
		return p.Version
	}
	dir, _ := os.Getwd()
	p.cd()
	ver := output("git", "log", "-1", "--format=%h")
	cd(dir)
	return ver
}

func dependents(arg string) []string {
	var names []string
	for _, p := range packages {
		if p.Name == arg {
			continue
		}
		p.configDeps()
		if in(p.Deps, arg) {
			names = append(names, p.Name)
		}
	}
	return names
}

func addDeps(args []string) []string {
	var names []string
	used := make(map[string]bool)
	var walk func(*Package)
	walk = func(p *Package) {
		if used[p.Name] || p.Name == "" {
			return
		}
		used[p.Name] = true
		for _, d := range p.Deps {
			pd := getPackage(d)
			walk(&pd)
		}
		names = append(names, p.Name)
	}
	for _, arg := range args {
		for _, pkg := range packages {
			if pkg.Name != arg {
				continue
			}
			pkg.configDeps()
			walk(&pkg)
		}
	}
	return names
}

func (p *Package) taggedDeps() string {
	var deps []string
	for _, dep := range addDeps(p.Deps) {
		if dep != "" {
			pdep := getPackage(dep)
			deps = append(deps, dep+"-"+pdep.gitVersion())
		}
	}
	if len(deps) == 0 {
		deps = append(deps, "none_listed")
	}
	return strings.Join(deps, ",")
}

func (p *Package) checkDeps() {
	xerces_c := getPackage("xerces-c")
	cernlib := getPackage("cernlib")
	hdds := getPackage("hdds")
	var cmds = map[string]*exec.Cmd{
		"xerces-c":            commande("ldd", xerces_c.Path+"/lib/libxerces-c.so"),
		"cernlib":             commande("ls", "-lh", cernlib.Path+"/"+cernlib.Version+"/lib/libgeant321.a"),
		"root":                commande("root", "-b", "-q", "-l"),
		"evio":                commande("evio2xml"),
		"rcdb":                commande("rcdb"),
		"ccdb":                commande("ccdb"),
		"jana":                commande("jana"),
		"hdds":                commande("ldd", hdds.Path+"/"+OS+"/"+"/lib/libhdds.so"),
		"sim-recon":           commande("hd_root"),
		"gluex_root_analysis": commande("root", "-b", "-q", "-l"),
	}
	if runtime.GOOS == "darwin" {
		cmds["xerces-c"] = commande("otool", "-L", xerces_c.Path+"/lib/libxerces-c.dylib")
		cmds["hdds"] = commande("otool", "-L", hdds.Path+"/"+OS+"/"+"/lib/libhdds.so")
	}
	for _, dep := range addDeps(p.Deps) {
		if cmds[dep] == nil {
			continue
		}
		if err := cmds[dep].Run(); err != nil {
			log.Fatalln(err)
		}
	}
	if p.Name != "sim-recon" {
		return
	}
	amptools := getPackage("amptools")
	if !isPath(filepath.Join(amptools.Path, "success.hdpm")) {
		unsetenv("AMPTOOLS")
		unsetenv("AMPPLOTTER")
	}
	// Check version compatibility of deps
	for _, shlib := range packages {
		if !shlib.in([]string{"xerces-c", "root", "ccdb"}) {
			continue
		}
		name_ver := shlib.Name + "-" + shlib.Version
		for _, user := range packages {
			if !user.in([]string{"amptools", "jana", "hdds"}) {
				continue
			}
			if !shlib.in(addDeps(user.Deps)) {
				continue
			}
			user.config()
			path := user.Path
			if user.Name != "amptools" {
				path = filepath.Join(path, OS)
			}
			if !isPath(filepath.Join(path, "success.hdpm")) {
				continue
			}
			user_name_ver := user.Name + "-" + user.Version
			text := readFile(filepath.Join(path, "success.hdpm"))
			d := strings.Split(text, "\n")
			record := d[len(d)-1]
			if !strings.Contains(record, name_ver) {
				fmt.Println("Error: " + name_ver + " is incompatible with " + user_name_ver + ".\n\t" + user_name_ver + " depends on " + record + ".\n\tRebuild " + user_name_ver + " against " + name_ver + ", or use required " + shlib.Name + " version.")
				os.Exit(2)
			}
		}
	}
}
