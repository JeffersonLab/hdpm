package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Create the fetch command
var cmdFetch = &cobra.Command{
	Use:   "fetch [PACKAGE...]",
	Short: "Fetch packages",
	Long: `Fetch packages.

Download and unpack packages into the $GLUEX_TOP directory.
If GLUEX_TOP is not set, packages are unpacked into the
current working directory.`,
	Example: `1. hdpm fetch sim-recon --deps
2. hdpm fetch root geant4
3. hdpm fetch --all
4. hdpm fetch sim-recon -d --xml https://halldweb.jlab.org/dist/version.xml`,
	Run: runFetch,
}

var force, deps, all bool

func init() {
	cmdHDPM.AddCommand(cmdFetch)

	cmdFetch.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
	cmdFetch.Flags().BoolVarP(&force, "force", "f", false, "Do not skip cernlib on macOS")
	cmdFetch.Flags().BoolVarP(&deps, "deps", "d", false, "Include dependencies")
	cmdFetch.Flags().BoolVarP(&all, "all", "a", false, "Fetch all packages in the package settings")
}

func runFetch(cmd *cobra.Command, args []string) {
	pkgInit()
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
	} else if deps {
		args = addDeps(args)
	}
	printPackages(args)

	// Change package versions to XMLfile versions
	if XML != "" {
		versionXML(XML)
	}
	// Change package versions to versions passed on command line
	changeVersions(args, versions)

	// Set proxy env. variables if on JLab CUE
	setenvJLabProxy()

	// Fetch packages
	mkcd(PD)
	for _, pkg := range packages {
		if !pkg.in(args) {
			continue
		}
		if runtime.GOOS == "darwin" && !force && pkg.Name == "cernlib" {
			fmt.Printf("macOS detected: Skipping %s\n", pkg.Name)
			continue
		}
		pkg.config()
		if pkg.isFetched() {
			fmt.Printf("Path exists: %s\n", pkg.Path)
			continue
		}
		pkg.fetch()
	}
}

func (p *Package) fetch() {
	if p.isFetched() {
		return
	}
	switch strings.HasSuffix(p.URL, ".tar.gz") || strings.HasSuffix(p.URL, ".tgz") {
	case true:
		if p.Name != "cernlib" {
			if p.usesCMake() {
				fetchTarfile(p.URL, p.Path+"/src")
			} else {
				fetchTarfile(p.URL, p.Path)
			}
		} else if p.Version == "2005" {
			p.mkcd()
			fetchTarfile(strings.Replace(p.URL, ".2005.corr.2014.04.17", "-2005-all-new", 1), "")
			fetchTarfile(strings.Replace(p.URL, "corr", "install", 1), "")
			run("curl", "-OL", p.URL)
			run("mv", "-f", path.Base(p.URL), "cernlib.2005.corr.tgz")
		}
	case false:
		if strings.Contains(p.URL, "svn") {
			if p.Version != "master" && !strings.Contains(p.URL, "tags") {
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", "-r", p.Version, p.URL, p.Path)
			} else {
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", p.URL, p.Path)
			}
		}
		if strings.Contains(p.URL, "git") && !strings.Contains(p.URL, "archive") {
			run("git", "clone", "-b", p.Version, p.URL, p.Path)
		}
	}
	fmt.Println()
}

func (p *Package) isFetched() bool {
	return isPath(p.Path)
}

func fetchTarfile(url, path string) {
	file := filepath.Base(url)
	fmt.Printf("Downloading %s ...\n", file)
	if strings.Contains(url, "https://") || strings.Contains(url, "http://") {
		run("curl", "-OL", url)
	} else {
		run("cp", "-p", url, ".")
	}
	fmt.Printf("\nUnpacking %s ...\n", file)
	if path != "" {
		mk(path)
		tar := exec.Command("tar", "tf", file)
		head := exec.Command("head", "-n1")
		tarOut, _ := tar.StdoutPipe()
		tar.Start()
		head.Stdin = tarOut
		headOut, _ := head.Output()
		ncomp := "2"
		if !strings.HasPrefix(string(headOut), "./") {
			ncomp = "1"
		}
		run("tar", "xf", file, "-C", path, "--strip-components="+ncomp)
	} else {
		run("tar", "xf", file)
	}
	os.Remove(file)
}
