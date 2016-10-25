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
	Short: "Write id of docker images",
	Long: `Write id of docker images.
	
tags: c6, c7, u14, u16

All tags will be written if no arguments are given.

Usage examples:
1. gxd write -u USER c6
`,
	Run: runWrite,
}

func init() {
	cmdGXD.AddCommand(cmdWrite)

	cmdWrite.Flags().StringVarP(&USER, "user", "u", "", "Docker username")
}

func runWrite(cmd *cobra.Command, args []string) {
	if USER == "" {
		fmt.Fprint(os.Stderr, "Please pass Docker username as a flag.\n")
		os.Exit(2)
	}
	var names = map[string]string{
		"c6": "centos6",
		"c7": "centos7",
		"u14": "ubuntu14",
		"u16": "ubuntu16",
	}

	var tags = []string{"c6", "c7", "u14", "u16"}

	for _, arg := range args {
		if !in(tags, arg) {
			fmt.Fprintf(os.Stderr, "%s: Unknown tag\n", arg)
			os.Exit(2)
		}
	}

	if len(args) == 0 {
		args = tags
	}

	name := "sim-recon-deps"
	wd := workDir()
	for _, tag := range tags {
		if !in(args, tag) {
			continue
		}
		repo := "quay.io" + "/" + USER + "/" + name + ":" + names[tag]
		if tag != "c6" && tag != "c7" {
			repo = USER + "/" + name + ":" + names[tag]
		}
		s := output("docker", "inspect", "--format='{{.Id}}'", repo)
		write_text(wd+"/.id-deps-"+tag, strings.Split(s, ":")[1][0 : 5])
	}
}
