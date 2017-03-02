package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Create the save command
var cmdSave = &cobra.Command{
	Use:   "save NAME",
	Short: "Save the current package settings",
	Long: `Save the current package settings.

Give a name for the package settings as the only argument.

Saved settings are restored by using the select command:
  hdpm select NAME`,
	Example: `1. hdpm save test
2. hdpm save root5
3. hdpm save whatever`,
	Run: runSave,
}

var comment string

func init() {
	cmdHDPM.AddCommand(cmdSave)

	cmdSave.Flags().StringVarP(&comment, "comment", "c", "", "Comment")
}

func runSave(cmd *cobra.Command, args []string) {
	pkgInit()
	if os.Getenv("GLUEX_TOP") == "" {
		fmt.Println("GLUEX_TOP environment variable is not set.\nSaving settings to the current working directory ...")
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Give a name for the package settings as the only argument.\n")
		cmd.Usage()
		os.Exit(2)
	}
	arg := args[0]
	dir := filepath.Join(PD, ".saved-settings", arg)
	type shift struct {
		current string
		next    string
	}
	var saved []shift
	d := dir
	for isPath(d) {
		s := strings.Split(d, "@t")
		sh := shift{}
		if len(s) > 1 {
			b := ""
			for n := 0; n < len(s)-1; n++ {
				b += s[n]
			}
			i, err := strconv.Atoi(s[len(s)-1])
			if err != nil {
				log.Fatalln(err)
			}
			sh = shift{d, b + "@t" + strconv.Itoa(i+1)}
		} else {
			sh = shift{d, d + "@t1"}
		}
		saved = append(saved, sh)
		d = sh.next
	}
	for n := len(saved) - 1; n >= 0; n-- {
		nd := saved[n].next
		os.Rename(saved[n].current, nd)
		s0 := &Settings{}
		if isPath(nd + "/.info.json") {
			s0.read(nd)
		}
		s0.Name = filepath.Base(nd)
		s0.write(nd)
	}
	mk(dir)
	for _, pkg := range packages {
		pkg.write(dir)
	}
	s := newSettings(arg, comment)
	s.write(SD)
	s.write(dir)
}
