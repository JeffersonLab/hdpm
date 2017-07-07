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

tags: c6, c7, u16

All tags will be pulled if no arguments are given.`,
	Example: `1. gxd pull -u USER c6`,
	Run:     runPull,
}

func init() {
	cmdGXD.AddCommand(cmdPull)

	cmdPull.Flags().StringVarP(&USER, "user", "u", "", "Docker username")
	cmdPull.Flags().BoolVarP(&rmi, "rmi", "", false, "Remove old image before pull")
}

func runPull(cmd *cobra.Command, args []string) {
	if USER == "" {
		exitNoUsername(cmd)
	}
	var names = map[string]string{
		"c6":  "centos6",
		"c7":  "centos7",
		"u16": "ubuntu16",
	}

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

	name := "sim-recon-deps"
	for _, tag := range tags {
		if !in(args, tag) {
			continue
		}
		repo := USER + "/" + name + ":" + names[tag]
		run("docker", "pull", repo)
		if rmi {
			run("docker", "rmi", name+":"+tag)
		}
		run("docker", "tag", repo, name+":"+tag)
		run("docker", "rmi", repo)
	}
}
