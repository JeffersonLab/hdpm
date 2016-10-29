package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Create the clean command
var cmdClean = &cobra.Command{
	Use:   "clean [PACKAGE...]",
	Short: "Clean build products of selected packages",
	Long: `
Clean/remove build products of selected packages.

If no arguments are given, the following packages will be cleaned:
ccdb, jana, hdds, sim-recon, gluex_root_analysis

Alternate usage:
hdpm clean DIRECTORY

Usage examples:
1. hdpm clean
2. hdpm clean sim-recon
3. hdpm clean sim-recon/master/src/plugins/Analysis/pi0omega
`,
	Run: runClean,
}

var obliterate bool

func init() {
	cmdHDPM.AddCommand(cmdClean)

	cmdClean.Flags().BoolVarP(&obliterate, "obliterate", "", false, "Clean for distribution, obliterate source code!")
}

func runClean(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nCleaning packages in the current working directory ...")
	}
	// Clean a sim-recon subdirectory if passed as argument
	cwd, _ := os.Getwd()
	if len(args) == 1 && isPath(filepath.Join(cwd, args[0])) {
		dir := filepath.Join(cwd, args[0])
		if isPath(filepath.Join(dir, "SConstruct")) || isPath(filepath.Join(dir, "SConscript")) {
			cd(dir)
			env("")
			run("scons", "-u", "-c", "install")
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
	}
	// Change package versions to versions passed on command line
	changeVersions(args, versions)
	// Set environment variables
	env("")
	for _, pkg := range packages {
		pkg.config()
		if !pkg.in(args) || !pkg.isFetched() || pkg.IsPrebuilt || pkg.inDist() {
			continue
		}
		pkg.cd()
		if obliterate {
			pkg.distclean()
		} else {
			pkg.clean()
		}
	}
}

func (p *Package) clean() {
	if p.Name == "ccdb" {
		run("rm", "-f", "success.hdpm", ".sconsign.dblite")
		run("scons", "-c")
	}
	if p.in([]string{"jana", "hdds", "sim-recon", "gluex_root_analysis"}) {
		run("rm", "-f", "success.hdpm", ".sconsign.dblite", "src/.sconsign.dblite")
		run("rm", "-rf", OS)
		if isPath("src") {
			run("rm", "-rf", "src/."+OS)
		}
		if p.Name == "gluex_root_analysis" {
			for _, dir := range []string{"libraries/DSelector", "programs/MakeDSelector", "programs/tree_to_amptools"} {
				run("rm", "-rf", dir+"/"+OS)
			}
		}
	}
}

func (p *Package) distclean() {
	if p.Name == "root" && strings.Contains(p.Cmds[0], "./configure") {
		if !isPath("Makefile") {
			return
		}
		run("cp", "-p", "success.hdpm", "../")
		run("make", "dist")
		cd("../")
		run("rm", "-rf", p.Version)
		files := glob(filepath.Dir(p.Path) + "/*.tar.gz")
		for _, file := range files {
			run("tar", "xf", file)
			os.Remove(file)
		}
		run("mv", "success.hdpm", p.Version)
	} else {
		if p.Name == "sim-recon" {
			run("rm", "-rf", "src/."+OS, "src/.sconsign.dblite", "."+OS)
		} else {
			run("rm", "-rf", "src", "."+OS)
		}
		rmGlob(p.Path + "/*.*gz")
		rmGlob(p.Path + "/*.contents")
		rmGlob(p.Path + "/.g*")
		rmGlob(p.Path + "/.s*")
		rmGlob(p.Path + "/setenv.*")
		rmGlob(p.Path + "/" + OS + "/setenv.*")
	}
}
