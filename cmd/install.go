package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Create the install command
var cmdInstall = &cobra.Command{
	Use:   "install [COMMIT]",
	Short: "Install binary distribution of sim-recon and deps",
	Long: `
Install binary distribution of sim-recon and dependencies.

Alternate Usage:
1. hdpm install TARFILE-URL | TARFILE-PATH

Usage Examples:
1. hdpm install -l
2. hdpm install
3. hdpm install -c
`,
	Run: runInstall,
}

var showList bool
var cleanLinks bool

func init() {
	cmdHDPM.AddCommand(cmdInstall)

	cmdInstall.Flags().BoolVarP(&showList, "list", "l", false, "List available binary distribution tarfiles.")
	cmdInstall.Flags().BoolVarP(&cleanLinks, "clean", "c", false, "Clean/remove symbolic links.")
}

func runInstall(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nInstalling packages to the current working directory ...")
	}
	arg := ""
	if len(args) >= 1 {
		arg = args[0]
	}
	distDir := filepath.Join(PD, ".dist")
	if !cleanLinks {
		fetchDist(arg)
	}
	OS = strings.Replace(OS, "RHEL", "CentOS", -1)
	OS = strings.Replace(OS, "LinuxMint17", "Ubuntu14", -1)
	OS = strings.Replace(OS, "LinuxMint18", "Ubuntu16", -1)
	if !cleanLinks {
		fmt.Println("\nLinking distribution binaries into " + PD + " ...")
	} else {
		fmt.Println("Removing symlinks in " + PD + " ...")
	}
	for _, pkg := range packages {
		pkg.install()
	}
	// Link env-setup scripts
	mk(PD + "/env-setup")
	rmGlob(PD + "/env-setup/dist.*")
	if cleanLinks {
		return
	}
	for _, sh := range []string{"sh", "csh"} {
		if isPath(distDir + "/env-setup/master." + sh) {
			run("ln", "-s", distDir+"/env-setup/master."+sh, PD+"/env-setup/dist."+sh)
		}
	}
}

func (p *Package) install() {
	pd := filepath.Join(PD, ".dist", p.Name)
	if !isPath(pd) {
		fmt.Printf("Not in distribution: %s\n", p.Name)
		return
	}
	v := dirVersion(pd)
	vd := v
	if p.Name == "hdds" || p.Name == "sim-recon" {
		v = distVersion(pd + "/" + v + "/" + OS)
	}
	if p.Name == "gluex_root_analysis" {
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
		removeSymLinks(d)
		return
	}
	if isPath(pi) {
		fmt.Printf("  Already installed: %s/%s\n", p.Name, v)
		return
	}
	mk(d)
	removeSymLinks(d)
	run("ln", "-s", pd, pi)
}

func removeSymLinks(dir string) {
	for _, lc := range readDir(dir) {
		file := dir + "/" + lc
		if isSymLink(file) {
			if strings.HasPrefix(readLink(file), PD+"/.dist/") {
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
Filename format:    sim-recon-<commit>-<id_deps>-<os_tag>.tar.gz
Available OS tags:  c6 (CentOS 6), c7 (CentOS 7),
                    u14 (Ubuntu 14.04), u16 (Ubuntu 16.04)
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
	if strings.Contains(OS, "Ubuntu14") || strings.Contains(OS, "LinuxMint17") {
		tag = "u14"
	}
	if strings.Contains(OS, "Ubuntu16") || strings.Contains(OS, "LinuxMint18") {
		tag = "u16"
	}
	if tag == "" {
		fmt.Fprintf(os.Stderr, "%s: Unsupported operating system\n", OS)
		os.Exit(2)
	}
	dir := filepath.Join(PD, ".dist")
	if showList {
		fmt.Println("Available tarfiles")
	}
	URL := latestURL(tag, arg)
	if showList {
		os.Exit(0)
	}
	isurl := strings.Contains(URL, "https://") || strings.Contains(URL, "http://")
	if isurl && !strings.Contains(URL, "https://halldweb.jlab.org/dist") {
		fmt.Fprintf(os.Stderr, "%s is an unfamiliar URL.\n", URL)
	}
	if !isurl && !strings.Contains(URL, "/group/halld/www/halldweb/html/dist") {
		fmt.Fprintf(os.Stderr, "%s is an unfamiliar path.\n", URL)
	}
	parts := strings.Split(URL, "-")
	if len(parts) != 5 || !strings.Contains(URL, "sim-recon") {
		fmt.Fprintf(os.Stderr, "%s: unsupported filename format.\n", URL)
		os.Exit(2)
	}
	commit := parts[2]
	idDeps := parts[3]
	tag2 := strings.Split(parts[len(parts)-1], ".")[0]
	if tag2 != tag {
		fmt.Fprintf(os.Stderr, "Warning: %s is for %s distribution.\nYou are on %s.\n", URL, tag2, tag)
	}
	urlDeps := strings.Replace(URL, commit, "deps", -1)
	deps := []string{"xerces-c", "cernlib", "root", "evio", "ccdb", "jana"}
	updateDeps := !contains(readDir(dir), deps)
	update := !isPath(filepath.Join(dir, "sim-recon")) || !isPath(filepath.Join(dir, "hdds"))

	if updateDeps || (isPath(dir+"/.id-deps-"+tag2) && idDeps != readFile(dir+"/.id-deps-"+tag2)) {
		os.RemoveAll(dir)
		updateDeps = true
		update = true
		mkcd(dir)
		fetchTarfile(urlDeps, dir)
		fetchTarfile(URL, dir)
	} else if update || commit != currentCommit(dir) {
		os.RemoveAll(dir + "/sim-recon")
		os.RemoveAll(dir + "/hdds")
		mkcd(dir)
		fetchTarfile(URL, dir)
		update = true
	} else {
		fmt.Printf("Already up-to-date, at commit=%s\n", commit)
	}
	if update {
		rmGlob(dir + "/version_*")
		run("touch", dir+"/version_sim-recon-"+commit+"_deps-"+idDeps)
	}
	if updateDeps {
		updateEnvScript(dir + "/env-setup/master.sh")
		updateEnvScript(dir + "/env-setup/master.csh")
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Environment setup")
	fmt.Println("source " + dir + "/env-setup/master.[c]sh")
	// Check consistency between commit records
	for _, d := range readDir(dir + "/sim-recon/master") {
		if strings.Contains(d, "Linux_") || strings.Contains(d, "Darwin_") {
			v := distVersion(dir + "/sim-recon/master/" + d)
			if commit != v {
				fmt.Fprintf(os.Stderr, "Inconsistent commits: %s and %s\n", commit, v)
				os.Exit(2)
			}
			for _, n := range []string{"jana", "hdds", "sim-recon", "gluex_root_analysis"} {
				p := filepath.Join(dir, n, readDir(filepath.Join(dir, n))[0])
				if OS != d && !isPath(filepath.Join(p, OS)) {
					run("ln", "-s", p+"/"+d, p+"/"+OS)
				}
			}
		}
	}
}

func currentCommit(path string) string {
	commit := ""
	for _, file := range readDir(path) {
		if strings.HasPrefix(file, "version_sim-recon-") {
			commit = strings.TrimPrefix(file, "version_sim-recon-")
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
	top_i, top_f := "GLUEX_TOP=.+", "GLUEX_TOP="+gx
	jli, jlf := "\\${GLUEX_TOP}/julia-.{5,7}/bin:", ""
	if strings.HasSuffix(path, ".csh") {
		set = "setenv"
		top_i, top_f = "GLUEX_TOP .+", "GLUEX_TOP "+gx
	}
	tobeReplaced := [2]string{top_i, jli}
	replacement := [2]string{top_f, jlf}
	for i := 0; i < len(tobeReplaced); i++ {
		re := regexp.MustCompile(tobeReplaced[i])
		data = re.ReplaceAllString(data, replacement[i])
	}
	data = strings.Replace(data, "/opt/rh/", "${GLUEX_TOP}/opt/rh/", -1)
	res_path := "/u/group/halld/www/halldweb/html/resources"
	if isPath(res_path) {
		data = strings.Replace(data, "#"+set+" JANA_RES", set+" JANA_RES", -1)
		data = strings.Replace(data, "/path/to/resources", res_path, -1)
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
			ok = strings.Contains(file, "-"+tag+".t") && strings.Contains(file, arg)
		} else {
			ok = strings.Contains(file, "-"+tag+".t") && !strings.Contains(file, "deps")
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
