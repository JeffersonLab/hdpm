package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Create the pull command
var cmdPull = &cobra.Command{
	Use:   "pull [TAG...]",
	Short: "Pull docker images",
	Long: `Pull docker images.
	
tags: c6, c7, u14, u16

All tags will be pulled if no arguments are given.

Usage examples:
1. gxd pull -u USER c6
`,
	Run: runPull,
}

func init() {
	cmdGXD.AddCommand(cmdPull)

	cmdPull.Flags().StringVarP(&USER, "user", "u", "", "Docker username")
}

func runPull(cmd *cobra.Command, args []string) {
	if USER == "" {
		fmt.Fprint(os.Stderr, "Please pass Docker username as flag.\n")
		os.Exit(2)
	}
	var names = map[string]string{
		"c6":  "centos6",
		"c7":  "centos7",
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
	for _, tag := range tags {
		if !in(args, tag) {
			continue
		}
		repo := "quay.io" + "/" + USER + "/" + name + ":" + names[tag]
		if tag != "c6" && tag != "c7" {
			repo = USER + "/" + name + ":" + names[tag]
		}
		run("docker", "pull", repo)
		run("docker", "rmi", name+":"+tag)
		run("docker", "tag", repo, name+":"+tag)
		run("docker", "rmi", repo)
	}
}
