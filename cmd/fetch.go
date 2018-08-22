package cmd

import (
	"fmt"
	"log"
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
Update Git/SVN repositories.`,
	Aliases: []string{"get", "pull", "update"},
	Example: `1. hdpm fetch
2. hdpm fetch halld_recon --deps
3. hdpm fetch root geant4
4. hdpm fetch halld_recon -d --xml https://halldweb.jlab.org/dist/version.xml`,
	Run: runFetch,
}

var deps, noCheckURL bool

func init() {
	cmdHDPM.AddCommand(cmdFetch)

	cmdFetch.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
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
		if runtime.GOOS == "darwin" && pkg.Name == "cernlib" {
			fmt.Printf("macOS detected: Skipping %s\n", pkg.Name)
			continue
		}
		pkg.config()
		if pkg.isFetched() {
			pkg.update()
		} else {
			pkg.fetch()
		}
	}
}

func (p *Package) fetch() {
	if strings.HasSuffix(p.URL, ".tar.gz") || strings.HasSuffix(p.URL, ".tgz") {
		switch p.Name {
		case "cernlib":
			p.mkcd()
			fetchTarfile(strings.Replace(p.URL, ".2005.corr.2014.04.17", "-2005-all-new", 1), "")
			fetchTarfile(strings.Replace(p.URL, "corr", "install", 1), "")
			run("curl", "-OL", p.URL)
			run("mv", "-f", path.Base(p.URL), "cernlib.2005.corr.tgz")
			cd(PD)
		default:
			if p.usesCMake() {
				fetchTarfile(p.URL, p.Path+"/src")
			} else {
				fetchTarfile(p.URL, p.Path)
			}
		}
	} else {
		switch {
		case strings.Contains(p.URL, "git") && !strings.Contains(p.URL, "archive") &&
			!strings.Contains(p.URL, "releases"):
			run("git", "clone", "-b", p.Version, p.URL, p.Path)
		case strings.Contains(p.URL, "svn") && !strings.Contains(p.URL, "tags"):
			switch p.Version {
			case "master":
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", p.URL, p.Path)
			default:
				run("svn", "checkout", "--non-interactive", "--trust-server-cert", "-r", p.Version, p.URL, p.Path)
			}
		default:
			if err := fetchURL(p.URL); err != nil {
				log.SetPrefix("fetch failed: ")
				log.SetFlags(0)
				log.Fatalln(err)
			}
			mk(p.Path)
			f := filepath.Base(p.URL)
			os.Rename(f, filepath.Join(p.Path, f))
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

func checkURL(url string) error {
	if noCheckURL {
		return nil
	}
	s := output("curl", "-ILsS", url)
	if strings.Contains(s, "200 OK") ||
		(strings.Contains(s, "403 Forbidden") && strings.Contains(s, "AmazonS3")) {
		return nil
	}
	status := "unknown"
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "HTTP/") {
			status = strings.Join(strings.Fields(line)[1:], " ")
		}
	}
	return fmt.Errorf("HTTP status: %s", status)
}

func fetchURL(url string) error {
	fmt.Printf("Downloading %s ...\nURL: %s\n", filepath.Base(url), url)
	var err error
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		err = checkURL(url)
		if err != nil {
			return err
		}
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
	fmt.Printf("Unpacking %s ...\nPath: %s\n", file, path)
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
	switch {
	case isPath(p.Path + "/.git"):
		p.cd()
		fmt.Printf("%s: Updating %s branch ...\n", p.Name, p.Version)
		run("git", "checkout", p.Version)
		run("git", "pull")
		cd(PD)
	case isPath(p.Path+"/.svn") && strings.Contains(p.URL, "svn") && !strings.Contains(p.URL, "tags"):
		p.cd()
		fmt.Printf("%s: Updating to svn revision %s ...\n", p.Name, p.Version)
		switch p.Version {
		case "master":
			run("svn", "update")
		default:
			run("svn", "update", "--non-interactive", "-r"+p.Version)
		}
		cd(PD)
	default:
		fmt.Printf("Path exists: %s\n", p.Path)
	}
}

func (p *Package) isRepo() bool {
	return isPath(p.Path+"/.git") || (isPath(p.Path+"/.svn") &&
		strings.Contains(p.URL, "svn") && !strings.Contains(p.URL, "tags"))
}
