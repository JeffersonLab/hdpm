package cmd

import "github.com/spf13/cobra"

// Create the add command
var cmdAdd = &cobra.Command{
	Use:   "add PACKAGE...",
	Short: "Add packages to the current settings",
	Long: `Add packages to the current package settings.

Restore one or more packages to the default settings.

To see available packages use:
1. hdpm show -e   (for extra packages)
2. hdpm show -m   (for default/master packages)`,
	Example: `1. hdpm add hd_utilities
2. hdpm add cmake gluex_workshops
3. hdpm add sim-recon   (to restore default settings)`,
	Run: runAdd,
}

func init() {
	cmdHDPM.AddCommand(cmdAdd)
}

func runAdd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		exitNoPackages(cmd)
	}
	pkgInit()
	newDir := false
	if !isPath(SD) {
		newDir = true
		mk(SD)
		s := newSettings("master", "Default settings of hdpm version "+VERSION)
		s.write(SD)
		for _, pkg := range masterPackages {
			pkg.write(SD)
		}
	}
	for _, pkg := range extraPackages {
		if !pkg.in(args) {
			continue
		}
		pkg.write(SD)
	}
	if newDir {
		return
	}
	for _, pkg := range masterPackages {
		if !pkg.in(args) {
			continue
		}
		pkg.write(SD)
	}
}
