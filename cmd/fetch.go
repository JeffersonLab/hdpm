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
	Short: "Fetch packages and dependencies",
	Long: `Fetch packages and dependencies.
	
Download and unpack packages into the $GLUEX_TOP directory.
If GLUEX_TOP is not set, packages are unpacked into the
current working directory.

All packages in the build template will be fetched if
no arguments are given.

Usage examples:
1. hdpm fetch
2. hdpm fetch root cmake
`,
	Run: runFetch,
}

func init() {
	cmdHDPM.AddCommand(cmdFetch)

	cmdFetch.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
}

func runFetch(cmd *cobra.Command, args []string) {
	if XML != "" {
		versionXML(XML)
	}
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nInstalling packages to the current working directory ...")
	}
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
	// Set http proxy env. variable if on JLab CUE
	if isPath("/w/work/halld/home") {
		setenv("http_proxy", "http://jprox.jlab.org:8081")
		setenv("https_proxy", "https://jprox.jlab.org:8081")
	}
	mkcd(packageDir())
	for _, pkg := range packages {
		if !pkg.in(args) {
			continue
		}
		ver, ok := versions[pkg.Name]
		pkg.changeVersion(ver, ok)
		if runtime.GOOS == "darwin" &&
			(pkg.Name == "cernlib" || pkg.Name == "cmake") {
			fmt.Printf("macOS detected: Skipping %s\n", pkg.Name)
			continue
		}
		if pkg.isFetched() {
			fmt.Printf("%s/%s exists\n", pkg.Name, pkg.Version)
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
			fetchTarfile(p.URL, p.Path)
		} else if p.Version == "2005" {
			p.mkcd()
			fetchTarfile(strings.Replace(p.URL, ".2005.corr.2014.04.17", "-2005-all-new", 1), "") // get the "all" file
			fetchTarfile(strings.Replace(p.URL, "corr", "install", 1), "")                        // get the "install" file
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
	fmt.Printf("Unpacking %s ...\n", file)
	if path != "" {
		mk(path)
		tar := exec.Command("tar", "tf", file)
		head := exec.Command("head")
		tarOut, _ := tar.StdoutPipe()
		tar.Start()
		head.Stdin = tarOut
		headOut, _ := head.Output()
		ncomp := "2"
		if !strings.HasPrefix(string(headOut), ".") {
			ncomp = "1"
		}
		run("tar", "xf", file, "-C", path, "--strip-components="+ncomp)
	} else {
		run("tar", "xf", file)
	}
	os.Remove(file)
}
