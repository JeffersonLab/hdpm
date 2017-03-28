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

The settings files are written to the $GLUEX_TOP/.hdpm/settings directory.`,
	Example: `1. hdpm select master (for default settings)
2. hdpm select my-saved-settings
3. hdpm select -l
4. hdpm select --rm test@t1 test@t4

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
Run "hdpm env" to write the env scripts to the .hdpm/env directory.`,
	Run: runSelect,
}

var XML string
var useGroupPath bool

func init() {
	cmdHDPM.AddCommand(cmdSelect)

	cmdSelect.Flags().BoolVarP(&showList, "list", "l", false, "List all saved package settings")
	cmdSelect.Flags().BoolVarP(&rm, "rm", "", false, "Remove one or more saved package settings")
	cmdSelect.Flags().StringVarP(&XML, "xml", "", "", "Version XMLfile URL or path")
	cmdSelect.Flags().BoolVarP(&useGroupPath, "group", "g", false, "Use group packages by passing a .hdpm XML")
}

func runSelect(cmd *cobra.Command, args []string) {
	pkgInit()
	tdir := filepath.Join(HD, "saved-settings")
	if isPath(PD+"/.saved-settings") && !isPath(tdir) {
		mk(HD)
		os.Rename(PD+"/.saved-settings", tdir)
	}
	if showList {
		dirs := readDir(tdir)
		if len(dirs) > 0 {
			fmt.Println("Saved settings")
			fmt.Printf("%s\n", strings.Join(dirs, " "))
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
			if s.Timestamp != "" {
				fmt.Printf("%s\n", s.Timestamp)
			}
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
				rmGlob(HD + "/env/" + dir + ".*")
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
			fmt.Fprintf(os.Stderr, "%s: Unknown settings id\n", arg)
			dirs := readDir(tdir)
			if len(dirs) > 0 {
				fmt.Fprintf(os.Stderr, "Saved settings: %s\n", strings.Join(dirs, ", "))
			}
			os.Exit(2)
		}
	}
	os.RemoveAll(SD)
	mk(SD)
	s := newSettings(arg, "Default settings of hdpm version "+VERSION)
	s.write(SD)
	for _, pkg := range masterPackages {
		pkg.write(SD)
	}
}
