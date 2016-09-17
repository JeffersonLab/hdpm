package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Create the save command
var cmdSave = &cobra.Command{
	Use:   "save TEMPLATE",
	Short: "Save the current build settings",
	Long: `
Save the current build settings as a new build template.

This ensures that they are not lost after switching to a different template.
Switch between templates by using the select command.

Pass the name of the new template as the only argument.
Names of predefined templates are not allowed.

Predefined templates: master, jlab, workshop-2016

Usage examples:
1. hdpm save test
2. hdpm save root5
3. hdpm save whatever
`,
	Run: runSave,
}

func init() {
	cmdHDPM.AddCommand(cmdSave)
}

func runSave(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nSaving settings to the current working directory ...")
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: Pass name of new template as single argument")
		os.Exit(2)
	}
	arg := args[0]
	dir := filepath.Join(packageDir(), "templates", arg)
	mk(dir)
	for _, pkg := range packages {
		pkg.write(dir)
	}
}
