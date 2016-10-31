package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Create the push command
var cmdPush = &cobra.Command{
	Use:   "push [TAG...]",
	Short: "Push docker images",
	Long: `Push docker images.
	
tags: c6, c7, u14, u16

All tags will be pushed if no arguments are given.

Usage examples:
1. gxd push -u USER c6
`,
	Run: runPush,
}

var USER string

func init() {
	cmdGXD.AddCommand(cmdPush)

	cmdPush.Flags().StringVarP(&USER, "user", "u", "", "Docker username")
}

func runPush(cmd *cobra.Command, args []string) {
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
		run("docker", "tag", "hddeps:"+tag, repo)
		run("docker", "push", repo)
		run("docker", "rmi", repo)
	}
}