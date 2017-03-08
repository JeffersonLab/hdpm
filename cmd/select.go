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
	Use:   "select NAME",
	Short: "Select package settings",
	Long: `Select package settings.

The settings files are written to the $GLUEX_TOP/settings directory.`,
	Example: `1. hdpm select master (for default settings)
2. hdpm select my-saved-settings

Usage:
  hdpm select --xml XMLFILE-URL | XMLFILE-PATH
Examples:
1. hdpm select --xml latest
2. hdpm select --xml https://halldweb.jlab.org/dist/version_1.27.xml

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
Run "hdpm env" to write the env-setup scripts to the env-setup directory.`,
	Run: runSelect,
}

var XML string

func init() {
	cmdHDPM.AddCommand(cmdSelect)

	cmdSelect.Flags().BoolVarP(&showList, "list", "l", false, "List saved package settings")
	cmdSelect.Flags().BoolVarP(&rm, "rm", "", false, "Remove saved package settings")
	cmdSelect.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
}

func runSelect(cmd *cobra.Command, args []string) {
	pkgInit()
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nWriting settings to the current working directory ...")
	}
	tdir := filepath.Join(PD, ".saved-settings")
	if showList {
		dirs := readDir(tdir)
		if len(dirs) > 0 {
			fmt.Println("Saved settings")
			fmt.Printf("%s\n", strings.Join(dirs, ", "))
		}
		for _, dir := range dirs {
			s := &Settings{}
			if isPath(tdir + "/" + dir + "/.info.json") {
				s.read(tdir + "/" + dir)
			} else {
				s.Name = dir
			}
			fmt.Println(strings.Repeat("-", 80))
			fmt.Printf("%s\n", s.Name)
			if s.Comment != "" {
				fmt.Printf("%s\n", s.Comment)
			}
			fmt.Printf("%s\n", s.Timestamp)
		}
		return
	}
	if XML != "" {
		versionXML(XML)
		return
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No settings were specified on the command line.\n")
		cmd.Usage()
		os.Exit(2)
	}
	if rm {
		for _, dir := range readDir(tdir) {
			if !in(args, dir) {
				continue
			}
			if isPath(tdir + "/" + dir) {
				os.RemoveAll(tdir + "/" + dir)
				fmt.Printf("Removed: %s\n", dir)
			}
		}
		return
	}
	arg := "master"
	if len(args) >= 1 {
		arg = args[0]
	}
	if arg != "master" {
		if isPath(tdir + "/" + arg) {
			os.RemoveAll(SD)
			run("cp", "-pr", tdir+"/"+arg, SD)
			return
		} else {
			fmt.Fprintf(os.Stderr, "Unknown settings id:\n%s does not exist.\n",
				tdir+"/"+arg)
			fmt.Fprintf(os.Stderr, "\nSaved settings: %s\n", strings.Join(readDir(tdir), ", "))
			os.Exit(2)
		}
	}
	mk(SD)
	s := newSettings(arg, "Default settings of hdpm version "+VERSION)
	s.write(SD)
	for _, pkg := range masterPackages {
		pkg.write(SD)
	}
}
