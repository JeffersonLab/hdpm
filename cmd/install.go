package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Create the install command
var cmdInstall = &cobra.Command{
	Use:     "install [PACKAGE...]",
	Short:   "Install packages and dependencies",
	Long:    `Install packages and dependencies.`,
	Aliases: []string{"build"},
	Example: `1. hdpm install
2. hdpm install halld_recon
3. hdpm install geant4 amptools xerces-c
4. hdpm install halld_recon --xml https://halldweb.jlab.org/dist/version.xml
5. hdpm install gluex_root_analysis -i
6. hdpm install --dist

Usage:
  hdpm install DIRECTORY
Example:
  hdpm install halld_recon/master/src/plugins/Analysis/pi0omega`,
	Run: runInstall,
}

var jobs string
var showInfo, dist bool

func init() {
	cmdHDPM.AddCommand(cmdInstall)

	cmdInstall.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
	cmdInstall.Flags().StringVarP(&jobs, "jobs", "j", "", "Number of jobs to run in parallel")
	cmdInstall.Flags().BoolVarP(&showInfo, "info", "i", false, "Show current build information and exit")
	cmdInstall.Flags().BoolVarP(&dist, "dist", "", false, "Install binary distribution of GlueX software")
	cmdInstall.Flags().BoolVarP(&noCheckURL, "no-check-url", "", false, "Do not check URL")
}

func runInstall(cmd *cobra.Command, args []string) {
	if dist || showList || cleanLinks || osTag != "" {
		installDist(args)
		return
	}
	pkgInit()
	// Build a halld_recon subdirectory if passed as argument
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
	if len(args) == 0 {
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
	isFirst := true
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
			if !pkg.isFetched() {
				pkg.fetch()
			}
			if pkg.isInstalled() {
				fmt.Printf("Already installed: %s\n", pkg.Path)
				continue
			}
			pkg.install(&isFirst)
		}
	}
}

func (p *Package) isInstalled() bool {
	return p.IsPrebuilt || p.Name == "hd_utilities" || p.Name == "gluex_MCwrapper" || p.Name == "virtualenv"
}

func (p *Package) install(isFirst *bool) {
	ti := time.Now().Round(time.Second)
	fname := p.Path + "/success.hdpm"
    if p.in([]string{"jana", "hdds", "halld_recon", "halld_sim"}) {
		fname = p.Path + "/" + OS + "/success.hdpm"
	}
	if isPath(fname) && (showInfo || !p.isRepo()) {
		printInfo(fname, isFirst)
		return
	}
	fmt.Printf("\n%s: Checking dependencies ...\n", p.Name)
	p.checkDeps()
	fmt.Printf("Installing %s-%s ...\n", p.Name, p.gitVersion())
	p.cd()
	du_i := strings.Fields(output("du", "-sh", p.Path))[0]
	switch p.Name {
	case "cernlib":
		prep_cernlib_patches()
		run("sh", "-c", "patch < Install_cernlib.patch")
		run("./Install_cernlib")
		run("rm", "-rf", "2005/build", "2005/src")
		run("mv", "2005", "..")
		cd("..")
		run("rm", "-rf", "cernlib")
		p.mkcd()
		run("mv", "../2005", ".")
	case "pypwa":
		venv := getPackage("virtualenv")
		run("python", venv.Path+"/virtualenv.py", ".")
		whl := filepath.Base(p.URL)
		run("bash", "-c", ". bin/activate; pip install --trusted-host=pypi.python.org "+whl)
		os.Remove(whl)
	default:
		if p.Name == "halld_recon" || p.Name == "halld_sim" {
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
	}
	du_f := strings.Fields(output("du", "-sh", p.Path))[0]
	tf := time.Now().Round(time.Second)
	text := p.Name + "-" + p.gitVersion() + "\n" +
		tf.Format(time.RFC3339) + "\n" +
		"# install duration\n" +
		tf.Sub(ti).String() + "\n" +
		"# disk use, final minus initial\n" +
		du_f + "B - " + du_i + "B\n" +
		"# dependencies\n" +
		p.taggedDeps()
	write_text(fname, text)
	fmt.Println()
}

func printInfo(fname string, isFirst *bool) {
	if *isFirst {
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-18s%-18s%-18s%-18s\n", "package", "install time", "disk use", "timestamp")
		fmt.Println(strings.Repeat("-", 80))
	}
	text := readFile(fname)
	d := strings.Split(text, "\n")
	if strings.HasPrefix(d[0], "gluex_root_analysis") {
		d[0] = strings.Replace(d[0], "gluex_root_analysis", "gxRootAna", 1)
	}
	fmt.Printf("%-18s%-18s%-18s%-18s\n", d[0], d[3], d[5], d[1])
	*isFirst = false
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
	if p.IsPrebuilt || !isPath(p.Path+"/.git") {
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
	if p.Name == "pypwa" || len(p.Deps) == 0 {
		return
	}
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
		"halld_recon":         commande("hd_root"),
		"halld_sim":           commande("mcsmear"),
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
	if p.Name != "halld_recon" && p.Name != "halld_sim" {
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
		shlib.config()
		name_ver := shlib.Name + "-" + shlib.gitVersion()
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
			user_name_ver := user.Name + "-" + user.gitVersion()
			text := readFile(filepath.Join(path, "success.hdpm"))
			d := strings.Split(text, "\n")
			record := d[len(d)-1]
			if !strings.Contains(record, name_ver) {
				fmt.Fprintln(os.Stderr, "Error: "+name_ver+" is incompatible with "+user_name_ver+".\n\t"+user_name_ver+" depends on "+record+".\n\tRebuild "+user_name_ver+" against "+name_ver+", or use required "+shlib.Name+" version.")
				os.Exit(2)
			}
		}
	}
}

var showList, cleanLinks bool
var osTag string

func init() {
	cmdInstall.Flags().BoolVarP(&showList, "list", "l", false, "List available binary distribution tarfiles")
	cmdInstall.Flags().BoolVarP(&cleanLinks, "clean", "c", false, "Clean/remove symbolic links")
	cmdInstall.Flags().StringVarP(&osTag, "tag", "t", "", "Force selection of a OS tag (binary dist)")
}

func installDist(args []string) {
	pkgInit()
	arg := ""
	if len(args) >= 1 {
		arg = args[0]
	}
	distDir := filepath.Join(HD, "dist")
	if !cleanLinks {
		fetchDist(arg)
	}
	OS = strings.Replace(OS, "RHEL", "CentOS", 1)
	OS = strings.Replace(OS, "LinuxMint18", "Ubuntu16", 1)
	if !cleanLinks {
		fmt.Println("\nLinking distribution binaries into " + PD + " ...")
	} else {
		fmt.Println("Removing symlinks in " + PD + " ...")
	}
	for _, pkg := range packages {
		pkg.symlink()
	}
	// Link env scripts
	mk(HD + "/env")
	rmGlob(HD + "/env/dist.*")
	if cleanLinks {
		return
	}
	for _, sh := range []string{"sh", "csh"} {
		if isPath(distDir + "/.hdpm/env/master." + sh) {
			s := distDir + "/.hdpm/env/master." + sh
			l := HD + "/env/dist." + sh
			run("ln", "-s", relPath(filepath.Dir(l), s), l)
		}
	}
}

func (p *Package) symlink() {
	pd := filepath.Join(HD, "dist", p.Name)
	if !isPath(pd) {
		fmt.Printf("Not in distribution: %s\n", p.Name)
		return
	}
	v := dirVersion(pd)
	vd := v
	if p.Name == "hdds" || p.Name == "halld_recon" || p.Name == "halld_sim" {
		v = distVersion(pd + "/" + v + "/" + OS)
	}
	if p.in([]string{"hdgeant4", "gluex_root_analysis"}) {
		v = distVersion(pd + "/" + v)
	}
	pi := filepath.Join(PD, p.Name, v)
	if p.Name == "cernlib" {
		pi = filepath.Dir(pi)
	} else {
		pd += "/" + vd
	}
	d := filepath.Dir(pi)
	if cleanLinks && isPath(d) {
		removeSymlinks(d)
		return
	}
	if isPath(pi) {
		fmt.Printf("  Already installed: %s/%s\n", p.Name, v)
		return
	}
	mk(d)
	removeSymlinks(d)
	run("ln", "-s", relPath(d, pd), pi)
}

func removeSymlinks(dir string) {
	for _, lc := range readDir(dir) {
		file := dir + "/" + lc
		if isSymlink(file) {
			s := readLink(file)
			if strings.HasPrefix(s, "../.hdpm/dist/") || strings.HasPrefix(s, ".hdpm/dist/") {
				os.Remove(file)
			}
		}
	}
}

func distVersion(dir string) string {
	text := readFile(dir + "/success.hdpm")
	firstLine := strings.Split(text, "\n")[0]
	a := strings.Split(firstLine, "-")
	return a[len(a)-1]
}

func dirVersion(dir string) string {
	ver := ""
	for _, v := range readDir(dir) {
		if v != "success.hdpm" {
			ver = v
		}
	}
	return ver
}

func fetchDist(arg string) {
	msg :=
		`Tarfile base URL:   https://halldweb.jlab.org/dist
Path on JLab CUE:   /group/halld/www/halldweb/html/dist
Filename format:    halld_recon-<commit>-<id_deps>-<os_tag>.tar.gz
Available OS tags:  c6 (CentOS 6), c7 (CentOS 7), u16 (Ubuntu 16.04)
`
	fmt.Print(msg)
	fmt.Println(strings.Repeat("-", 80))
	tag := ""
	if strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6") {
		tag = "c6"
	}
	if strings.Contains(OS, "CentOS7") || strings.Contains(OS, "RHEL7") {
		tag = "c7"
	}
	if strings.Contains(OS, "Ubuntu16") || strings.Contains(OS, "LinuxMint18") {
		tag = "u16"
	}
	if osTag != "" {
		tag = osTag
	}
	if tag == "" {
		fmt.Fprintf(os.Stderr, "%s: Binary distribution does not exist for this OS\n", OS)
		os.Exit(2)
	}
	if !in([]string{"c6", "c7", "u16"}, tag) {
		fmt.Fprintf(os.Stderr, "%s: Unknown OS tag\n", tag)
		os.Exit(2)
	}
	dir := filepath.Join(HD, "dist")
	if showList {
		fmt.Println("Available tarfiles")
	}
	URL := latestURL(tag, arg)
	if showList {
		os.Exit(0)
	}
	parts := strings.Split(URL, "-")
	if len(parts) != 5 || !strings.Contains(URL, "halld_") {
		fmt.Fprintf(os.Stderr, "%s: Unsupported filename format\n", URL)
		os.Exit(2)
	}
	commit := parts[2]
	idDeps := parts[3]
	tag2 := strings.Split(parts[len(parts)-1], ".")[0]
	if tag2 != tag {
		fmt.Fprintf(os.Stderr, "Warning: %s is for %s distribution.\nYou are on %s.\n", URL, tag2, tag)
	}
	urlDeps := strings.Replace(URL, commit, "deps", 1)
	deps := []string{"xerces-c", "cernlib", "root", "evio", "ccdb", "jana"}
	updateDeps := !contains(readDir(dir), deps)
	update := !isPath(filepath.Join(dir, "halld_recon")) ||!isPath(filepath.Join(dir, "halld_sim")) || !isPath(filepath.Join(dir, "hdds"))

	if updateDeps || (isPath(dir+"/.id-deps-"+tag2) && idDeps != readFile(dir+"/.id-deps-"+tag2)) {
		os.RemoveAll(dir)
		updateDeps = true
		update = true
		mkcd(dir)
		fetchTarfile(urlDeps, dir)
		fetchTarfile(URL, dir)
	} else if update || commit != currentCommit(dir) {
		for _, n := range []string{"halld_recon", "halld_sim", "hdds", "hdgeant4", "gluex_root_analysis"} {
			os.RemoveAll(filepath.Join(dir, n))
		}
		mkcd(dir)
		fetchTarfile(URL, dir)
		update = true
	} else {
		fmt.Printf("Already up-to-date, at commit=%s\n", commit)
	}
	if update {
		rmGlob(dir + "/version_*")
		run("touch", dir+"/version_halld_recon-"+commit+"_deps-"+idDeps)
	}
	if updateDeps && isPath(dir+"/.hdpm/env/master.sh") {
		updateEnvScript(dir + "/.hdpm/env/master.sh")
		updateEnvScript(dir + "/.hdpm/env/master.csh")
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Environment setup")
	fmt.Println("source $GLUEX_TOP/.hdpm/env/dist.[c]sh")
	// Check consistency between commit records
	for _, d := range readDir(dir + "/halld_recon/master") {
		if strings.HasPrefix(d, "Linux_") || strings.HasPrefix(d, "Darwin_") {
			c := distVersion(dir + "/halld_recon/master/" + d)
			if commit != c {
				fmt.Fprintf(os.Stderr, "Inconsistent commits: %s and %s\n", commit, c)
				os.Exit(2)
			}
			for _, n := range []string{"jana", "hdds", "halld_recon", "halld_sim", "gluex_root_analysis"} {
				v := dirVersion(filepath.Join(dir, n))
				if v == "" {
					continue
				}
				p := filepath.Join(dir, n, v)
				if OS != d && !isPath(filepath.Join(p, OS)) {
					run("ln", "-s", d, p+"/"+OS)
				}
			}
		}
	}
}

func currentCommit(path string) string {
	commit := ""
	for _, file := range readDir(path) {
		if strings.HasPrefix(file, "version_halld_recon-") {
			commit = strings.TrimPrefix(file, "version_halld_recon-")
			commit = strings.Split(commit, "_")[0]
			break
		}
	}
	return commit
}

func updateEnvScript(path string) {
	data := readFile(path)
	gx := filepath.Dir(filepath.Dir(path))
	set := "export"
	top_i, top_f := "GLUEX_TOP=.+", "GLUEX_TOP=\""+gx+"\""
	if strings.HasSuffix(path, ".csh") {
		set = "setenv"
		top_i, top_f = "GLUEX_TOP .+", "GLUEX_TOP \""+gx+"\";"
	}
	re := regexp.MustCompile(top_i)
	data = re.ReplaceAllString(data, top_f)
	res_path := "/u/group/halld/www/halldweb/html/resources"
	if isPath(res_path) {
		data = strings.Replace(data, "#"+set+" JANA_RES", set+" JANA_RES", 1)
		data = strings.Replace(data, "/path/to/resources", res_path, 1)
	}
	write_text(path, data)
}

func latestURL(tag string, arg string) string {
	files := make(map[time.Time]string)
	latest_file := ""
	var latest_dt time.Time
	page := output("curl", "-s", "https://halldweb.jlab.org/dist/")
	lines := strings.Split(page, "\n")
	for _, line := range lines {
		re := regexp.MustCompile("href=\".{25,50}\"")
		r := re.FindString(line)
		if r == "" {
			continue
		}
		file := r[6 : len(r)-1]
		ok := false
		if arg != "" {
			ok = strings.HasPrefix(file, "halld_recon-"+arg) && strings.HasSuffix(file, "-"+tag+".tar.gz")
		} else {
			ok = strings.HasPrefix(file, "halld_recon-") && strings.HasSuffix(file, "-"+tag+".tar.gz") &&
				!strings.HasPrefix(file, "halld_recon-deps-")
		}
		if ok {
			re = regexp.MustCompile("(\\d{4})-(\\d{2})-(\\d{2})")
			r = re.FindString(line)
			s1 := strings.Split(r, "-")
			re = regexp.MustCompile("(\\d{2}):(\\d{2})")
			r = re.FindString(line)
			s2 := strings.Split(r, ":")
			dt := mkDate(s1, s2)
			files[dt] = file
			if latest_dt.Before(dt) {
				latest_dt = dt
				latest_file = file
			}
		}
	}
	if latest_file == "" {
		fmt.Fprintf(os.Stderr, "File not found at https://halldweb.jlab.org/dist for %s OS tag.\n", tag)
		os.Exit(2)
	}
	if !showList {
		fmt.Printf("Chosen tarfile: %s    %v\n", latest_file, latest_dt)
	} else {
		var times []time.Time
		for t, _ := range files {
			times = append(times, t)
		}
		sort.Sort(ts(times))
		for _, t := range times {
			fmt.Printf("%s        %v\n", files[t], t)
		}
	}
	return "https://halldweb.jlab.org/dist/" + latest_file
}

// Sort by time (satisfy sort interface)
type ts []time.Time

func (a ts) Len() int           { return len(a) }
func (a ts) Less(j, i int) bool { return a[i].Before(a[j]) }
func (a ts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func mkDate(s1 []string, s2 []string) time.Time {
	s1 = append(s1, s2...)
	var n []int
	for _, s := range s1 {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalln(err)
		}
		n = append(n, i)
	}
	return time.Date(n[0], time.Month(n[1]), n[2], n[3], n[4], 0, 0, time.Local)
}
