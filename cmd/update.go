package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Create the update command
var cmdUpdate = &cobra.Command{
	Use:   "update [PACKAGE...]",
	Short: "Update selected Git/SVN repositories",
	Long: `Update selected Git/SVN repositories.

Update all repos if no arguments are given.

Usage examples:
1. hdpm update
2. hdpm update sim-recon
3. hdpm update rcdb hdds
`,
	Run: runUpdate,
}

func init() {
	cmdHDPM.AddCommand(cmdUpdate)
}

func runUpdate(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nUpdating packages in the current working directory ...")
	}
	for _, arg := range args {
		if !in(packageNames, arg) {
			fmt.Printf("%s: unknown package name\n", arg)
			os.Exit(2)
		}
	}
	if len(args) == 0 {
		args = packageNames
	} else {
		args = addDeps(args)
	}
	mkcd(packageDir())
	for _, pkg := range packages {
		if !pkg.in(args) {
			continue
		}
		if !pkg.isFetched() {
			fmt.Printf("%s/%s does not exist\n", pkg.Name, pkg.Version)
			continue
		}
		pkg.cd()
		pkg.update()
	}
}

func (p *Package) update() {
	if strings.Contains(p.URL, "svn") && !strings.Contains(p.URL, "tags") {
		fmt.Printf("Updating %s to svn revision %s.\n", p.Name, p.Version)
		if p.Version != "master" {
			run("svn", "update", "--non-interactive", "-r"+p.Version)
		} else {
			run("svn", "update")
		}
	}
	if strings.Contains(p.URL, "git") && !strings.Contains(p.URL, "archive") {
		fmt.Printf("\nUpdating %s branch of %s.\n", p.Version, p.Name)
		run("git", "checkout", p.Version)
		run("git", "pull")
	}
}
