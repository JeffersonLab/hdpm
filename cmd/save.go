package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// Create the save command
var cmdSave = &cobra.Command{
	Use:   "save NAME",
	Short: "Save the current package settings",
	Long: `
Save the current package settings.

Give a name for the package settings as the only argument.

Saved settings are restored by using the select command:
hdpm select NAME

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
		fmt.Fprintln(os.Stderr, "Error: Pass name/id of settings as single argument")
		os.Exit(2)
	}
	arg := args[0]
	dir := filepath.Join(packageDir(), ".saved-settings", arg)
	if isPath(dir) {
		t := time.Now().Round(time.Second)
		os.Rename(dir, dir+"_"+t.Format(time.RFC3339))
	}
	mk(dir)
	for _, pkg := range packages {
		pkg.write(dir)
	}
}
