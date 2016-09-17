package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Create the select command
var cmdSelect = &cobra.Command{
	Use:   "select [TEMPLATE]",
	Short: "Select a build template",
	Long: `
Select a build template.
The master (default) template is selected if no argument is given.

This command is used to select the desired packages and settings to
use for the next builds.

Predefined templates: master, jlab, workshop-2016

Usage examples:
1. hdpm select (same as: hdpm select master)
2. hdpm select jlab
3. hdpm select workshop-2016
`,
	Run: runSelect,
}

var XML string

func init() {
	cmdHDPM.AddCommand(cmdSelect)

	cmdSelect.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
}

func runSelect(cmd *cobra.Command, args []string) {
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nWriting settings to the current working directory ...")
	}
	arg := "master"
	if len(args) == 1 {
		arg = args[0]
	}
	if arg == "master" || arg == "jlab" {
		packages = masterPackages
	}
	if arg == "workshop-2016" {
		packages = append(packages, ws16Package)
	}
	dir := filepath.Join(packageDir(), "settings")
	mk(dir)
	write_text(dir+"/.id", arg)
	for _, pkg := range packages {
		pkg.config(arg)
		pkg.write(dir)
	}
	if XML != "" {
		versionXML(XML)
	}
}
