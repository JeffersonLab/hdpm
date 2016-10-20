package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Create the build command
var cmdBuild = &cobra.Command{
	Use:   "build [PACKAGE...]",
	Short: "Build packages and dependencies",
	Long: `Build packages and dependencies.
	
Display build information if a package is already built.

Alternate usage:
hdpm build --xml XMLFILE-URL | XMLFILE-PATH
hdpm build DIRECTORY

All packages in the package settings will be built if
no arguments are given.

Usage examples:
1. hdpm build
2. hdpm build geant4 amptools
3. hdpm build --xml https://halldweb.jlab.org/dist/version.xml
4. hdpm build sim-recon/master/src/plugins/Analysis/pi0omega
`,
	Run: runBuild,
}

func init() {
	cmdHDPM.AddCommand(cmdBuild)

	cmdBuild.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
}

func runBuild(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nInstalling packages to the current working directory ...")
	}
	// Build a sim-recon subdirectory if passed as argument
	cwd, _ := os.Getwd()
	if len(args) == 1 && isPath(filepath.Join(cwd, args[0])) {
		dir := filepath.Join(cwd, args[0])
		if isPath(filepath.Join(dir, "SConstruct")) || isPath(filepath.Join(dir, "SConscript")) {
			cd(dir)
			env("")
			run("scons", "-u", "install")
			return
		}
	}
	// Parse args
	versions := extractVersions(args)
	args = extractNames(args)
	for _, arg := range args {
		if !in(packageNames, arg) {
			fmt.Fprintf(os.Stderr, "%s: Unknown package name\n", arg)
			os.Exit(2)
		}
	}
	if len(args) == 0 {
		args = packageNames
	} else {
		args = addDeps(args)
	}
	// Change package versions to XMLfile versions
	if XML != "" {
		versionXML(XML)
	}
	// Change package versions to versions passed on command line
	changeVersions(args, versions)
	// Set environment variables
	env("")
	// Fetch and build packages
	mkcd(packageDir())
	isBuilt := false
	for _, pkg := range packages {
		if !pkg.in(args) {
			continue
		}
		if runtime.GOOS == "darwin" && pkg.Name == "cernlib" {
			fmt.Printf("macOS detected: Skipping %s\n", pkg.Name)
			continue
		}
		pkg.fetch()
		if pkg.IsPrebuilt {
			fmt.Printf("Skipping prebuilt package: %s\n", pkg.Name)
			continue
		}
		pkg.build(&isBuilt)
	}
}

func (p *Package) build(isBuilt *bool) {
	ti := time.Now().Round(time.Second)
	fname := p.Path + "/success.hdpm"
	if p.in([]string{"jana", "hdds", "sim-recon"}) {
		fname = p.Path + "/" + OS + "/success.hdpm"
	}
	if isPath(fname) {
		printStats(fname, isBuilt)
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
		if p.usesCMake() {
			setenvPath(getPackage("cmake").Path)
			mkcd("../" + p.Name + "-build")
			run("mv", p.Path, "../"+p.Name)
			mk(p.Path)
		}
		for _, cmd := range p.Cmds {
			run("sh", "-c", cmd)
		}
		if p.usesCMake() {
			run("rm", "-rf", "../"+p.Name+"-build", "../"+p.Name)
		}
	} else {
		prep_cernlib_patches()
		run("sh", "-c", "patch < Install_cernlib.patch")
		run("./Install_cernlib")
		run("rm", "-rf", "2005/build", "2005/src")
		run("mv", "2005", "../")
		cd("../")
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

func printStats(fname string, isBuilt *bool) {
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

func addDeps(args []string) []string {
	var deps []string
	for _, arg := range args {
		for _, pkg := range packages {
			if pkg.Name != arg {
				continue
			}
			for _, dep := range pkg.Deps {
				if !in(args, dep) && dep != "" {
					deps = append(deps, dep)
				}
			}
		}
	}
	if len(deps) > 0 {
		fmt.Printf("Dependencies: %s\n", strings.Join(deps, ", "))
	}
	deps = append(deps, args...)
	return deps
}

func (p *Package) taggedDeps() string {
	var deps []string
	for _, dep := range p.Deps {
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
	ldd := "ldd"
	oe := "so"
	if runtime.GOOS == "darwin" {
		ldd = "otool -L"
		oe = "dylib"
	}
	xerces_c := getPackage("xerces-c")
	cernlib := getPackage("cernlib")
	hdds := getPackage("hdds")
	var cmds = map[string]*exec.Cmd{
		"xerces-c":            commande(ldd, xerces_c.Path+"/lib/libxerces-c."+oe),
		"cernlib":             commande("ls", "-lh", cernlib.Path+"/"+cernlib.Version+"/lib/libgeant321.a"),
		"root":                commande("root", "-b", "-q", "-l"),
		"evio":                commande("evio2xml"),
		"rcdb":                commande("rcdb"),
		"ccdb":                commande("ccdb"),
		"jana":                commande("jana"),
		"hdds":                commande(ldd, hdds.Path+"/"+OS+"/"+"/lib/libhdds.so"),
		"sim-recon":           commande("hd_root"),
		"gluex_root_analysis": commande("root", "-b", "-q", "-l"),
	}
	for _, dep := range p.Deps {
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
			if !shlib.in(user.Deps) {
				continue
			}
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
