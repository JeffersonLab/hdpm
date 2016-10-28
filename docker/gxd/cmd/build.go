package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Create the build command
var cmdBuild = &cobra.Command{
	Use:   "build [TAG...] [STAGE...]",
	Short: "Build docker images",
	Long: `Build docker images.
	
tags: c6, c7, u14, u16
stages: base, deps, sim-recon

All tags/stages will be built if no arguments are given.

Usage examples:
1. gxd build
2. gxd build sim-recon c6
`,
	Run: runBuild,
}

var rmi bool

func init() {
	cmdGXD.AddCommand(cmdBuild)

	cmdBuild.Flags().BoolVarP(&rmi, "rmi", "", false, "Remove old image before build.")
}

func runBuild(cmd *cobra.Command, args []string) {
	var names = map[string]string{
		"c6":  "centos6",
		"c7":  "centos7",
		"u14": "ubuntu14",
		"u16": "ubuntu16",
	}

	var tags = []string{"c6", "c7", "u14", "u16"}
	var stages = []string{"base", "deps", "sim-recon"}

	var ts, ss []string
	for _, arg := range args {
		if !in(tags, arg) && !in(stages, arg) {
			fmt.Fprintf(os.Stderr, "%s: Unknown tag or stage\n", arg)
			os.Exit(2)
		}
		if in(tags, arg) {
			ts = append(ts, arg)
		}
		if in(stages, arg) {
			ss = append(ss, arg)
		}
	}

	if len(ts) == 0 {
		ts = tags
	}
	if len(ss) == 0 {
		ss = stages
	}

	wd := workDir()
	for _, tag := range ts {
		dir := wd + "/buildfiles/" + names[tag]
		for _, stage := range ss {
			name, file := stage, "Dockerfile"
			if in([]string{"base", "deps"}, stage) {
				name, file = "hd"+stage, "Dockerfile-"+stage
			}
			if rmi {
				run("docker", "rmi", name+":"+tag)
			}
			log := output("docker", "build", "--no-cache", "-t", name+":"+tag, "-f", dir+"/"+file, dir)
			write_text(wd+"/.log-"+name+"-"+tag, log)
		}
	}
}
