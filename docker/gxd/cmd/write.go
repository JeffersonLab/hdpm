package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Create the write command
var cmdWrite = &cobra.Command{
	Use:   "write [TAG...]",
	Short: "Write the docker-image id",
	Long: `Write the docker-image id.

tags: c6, c7, u16

The ids of all tags will be written if no arguments are given.`,
	Example: `1. gxd write c6`,
	Run:     runWrite,
}

func init() {
	cmdGXD.AddCommand(cmdWrite)
}

func runWrite(cmd *cobra.Command, args []string) {
	var tags = []string{"c6", "c7", "u16"}

	for _, arg := range args {
		if !in(tags, arg) {
			fmt.Fprintf(os.Stderr, "%s: Unknown tag\n", arg)
			os.Exit(2)
		}
	}

	if len(args) == 0 {
		args = tags
	}

	wd := workDir()
	for _, tag := range tags {
		if !in(args, tag) {
			continue
		}
		s := output("docker", "inspect", "--format='{{.Id}}'", "hddeps:"+tag)
		write_text(wd+"/.id-deps-"+tag, strings.Split(s, ":")[1][0:5])
	}
}
