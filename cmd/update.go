package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Create the update command
var cmdUpdate = &cobra.Command{
	Use:   "update [PACKAGE...]",
	Short: "Update selected Git/SVN repositories",
	Long:  `Update selected Git/SVN repositories.`,
	Example: `1. hdpm update
2. hdpm update sim-recon
3. hdpm update sim-recon --deps
4. hdpm update rcdb hdds`,
	Run: runUpdate,
}

func init() {
	cmdHDPM.AddCommand(cmdUpdate)

	cmdUpdate.Flags().BoolVarP(&deps, "deps", "d", false, "Include dependencies")
}

func runUpdate(cmd *cobra.Command, args []string) {
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

	// Change package versions to versions passed on command line
	changeVersions(args, versions)

	// Set proxy env. variables if on JLab CUE
	setenvJLabProxy()

	// Update packages
	mkcd(PD)
	for _, pkg := range packages {
		pkg.config()
		if !pkg.in(args) || !pkg.isRepo() {
			continue
		}
		pkg.cd()
		pkg.update()
	}
}

func (p *Package) update() {
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
