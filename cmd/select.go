package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Create the select command
var cmdSelect = &cobra.Command{
	Use:   "select [NAME]",
	Short: "Select package settings",
	Long: `
Select package settings.

Default settings:
hdpm select master
hdpm select

Alternate usage:
hdpm select --xml XMLFILE-URL | XMLFILE-PATH

The XMLfile versions are applied on top of the current settings.

Shortcut URL ids:
latest   : https://halldweb.jlab.org/dist/version.xml
jlab-dev : https://halldweb.jlab.org/dist/version_jlab.xml
jlab     : https://halldweb.jlab.org/dist/version_jlab.xml

The JLab shortcuts will also set the paths of dependencies to the halld
group installations on the JLab CUE.

JLab development settings (for JLab CUE use only):
hdpm select --xml jlab-dev

If you use "jlab" instead of "jlab-dev", hdds and sim-recon will also
be set to the latest prebuilt packages on the JLab CUE.
Run "hdpm env" to write the env-setup scripts to the env-setup directory.

Usage examples:
1. hdpm select my-saved-settings
2. hdpm select --xml latest
3. hdpm select --xml https://halldweb.jlab.org/dist/version_1.27.xml
`,
	Run: runSelect,
}

var XML string

func init() {
	cmdHDPM.AddCommand(cmdSelect)

	cmdSelect.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
}

func runSelect(cmd *cobra.Command, args []string) {
	pkgInit()
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nWriting settings to the current working directory ...")
	}
	if XML != "" {
		versionXML(XML)
		return
	}
	arg := "master"
	if len(args) == 1 {
		arg = args[0]
	}
	if arg != "master" {
		tdir := filepath.Join(PD, ".saved-settings")
		if isPath(tdir + "/" + arg) {
			os.RemoveAll(SD)
			run("cp", "-pr", tdir+"/"+arg, SD)
			write_text(SD+"/.id", arg)
			return
		} else {
			fmt.Fprintf(os.Stderr, "\nError:\n%s does not exist.\n",
				tdir+"/"+arg)
			fmt.Fprintf(os.Stderr, "\nSaved settings: %v\n", strings.Join(readDir(tdir), ", "))
			os.Exit(2)
		}
	}
	mk(SD)
	write_text(SD+"/.id", arg)
	for _, pkg := range masterPackages {
		pkg.write(SD)
	}
}
