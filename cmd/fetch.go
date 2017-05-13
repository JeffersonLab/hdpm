package cmd

import (
	"fmt"
	"log"
	"net/http"
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

Download and unpack packages into the $GLUEX_TOP directory.`,
	Aliases: []string{"get", "pull", "update"},
	Example: `1. hdpm fetch
2. hdpm fetch sim-recon --deps
3. hdpm fetch root geant4
4. hdpm fetch sim-recon -d --xml https://halldweb.jlab.org/dist/version.xml`,
	Run: runFetch,
}

var force, deps, noCheckURL bool

func init() {
	cmdHDPM.AddCommand(cmdFetch)

	cmdFetch.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
	cmdFetch.Flags().BoolVarP(&force, "force", "f", false, "Do not skip cernlib on macOS")
	cmdFetch.Flags().BoolVarP(&deps, "deps", "d", false, "Include dependencies")
	cmdFetch.Flags().BoolVarP(&noCheckURL, "no-check-url", "", false, "Do not check URL")
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
	if len(args) == 0 {
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
			if !pkg.isRepo() {
				fmt.Printf("Path exists: %s\n", pkg.Path)
				continue
			}
			pkg.update()
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
		switch {
		case strings.Contains(p.URL, "svn"):
			if p.Version != "master" && !strings.Contains(p.URL, "tags") {
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", "-r", p.Version, p.URL, p.Path)
			} else {
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", p.URL, p.Path)
			}
		case strings.Contains(p.URL, "git") && !strings.Contains(p.URL, "archive"):
			run("git", "clone", "-b", p.Version, p.URL, p.Path)
		default:
			p.mkcd()
			err := fetchURL(p.URL)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	fmt.Println()
}

func (p *Package) isFetched() bool {
	return isPath(p.Path)
}

func fetchTarfile(url, path string) {
	if err := fetchTarfileError(url, path); err != nil {
		log.SetPrefix("fetch failed: ")
		log.SetFlags(0)
		log.Fatalln(err)
	}
}

func checkURL(url string) {
	if noCheckURL {
		return
	}
	resp, err := http.Head(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "URL-check failed: %s\n", url)
		log.SetPrefix("fetch failed: HTTP status: ")
		log.SetFlags(0)
		log.Fatalln(resp.Status)
	}
}

func fetchURL(url string) error {
	fmt.Printf("Downloading %s ...\n", filepath.Base(url))
	var err error
	if strings.Contains(url, "https://") || strings.Contains(url, "http://") {
		checkURL(url)
		err = runE("curl", "-OL", url)
	} else {
		err = runE("cp", "-p", url, ".")
	}
	return err
}

func fetchTarfileError(url, path string) error {
	err := fetchURL(url)
	if err != nil {
		return err
	}
	file := filepath.Base(url)
	fmt.Printf("\nUnpacking %s ...\n", file)
	defer os.Remove(file)
	if path == "" {
		return runE("tar", "xf", file)
	}
	tar := exec.Command("tar", "tf", file)
	head := exec.Command("head", "-n1")
	tarOut, err := tar.StdoutPipe()
	if err != nil {
		return err
	}
	tar.Start()
	head.Stdin = tarOut
	headOut, err := head.Output()
	if err != nil {
		return err
	}
	ncomp := "2"
	if !strings.HasPrefix(string(headOut), "./") {
		ncomp = "1"
	}
	mk(path)
	return runE("tar", "xf", file, "-C", path, "--strip-components="+ncomp)
}

func (p *Package) update() {
	p.cd()
	if strings.Contains(p.URL, "svn") && !strings.Contains(p.URL, "tags") {
		fmt.Printf("%s: Updating to svn revision %s ...\n", p.Name, p.Version)
		if p.Version != "master" {
			run("svn", "update", "--non-interactive", "-r"+p.Version)
		} else {
			run("svn", "update")
		}
	}
	if strings.Contains(p.URL, "git") && !strings.Contains(p.URL, "archive") {
		fmt.Printf("%s: Updating %s branch ...\n", p.Name, p.Version)
		run("git", "checkout", p.Version)
		run("git", "pull")
	}
}

func (p *Package) isRepo() bool {
	return isPath(p.Path+"/.git") || isPath(p.Path+"/.svn")
}
