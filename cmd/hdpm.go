package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Create the hdpm command
var cmdHDPM = &cobra.Command{
	Use:   "hdpm [COMMAND] [ARGS]",
	Short: "A tool for managing GlueX offline software and dependencies",
	Long:  `hdpm is a tool for managing GlueX offline software and dependencies.`,
}

// Execute a hdpm command
func Execute() {
	if err := cmdHDPM.Execute(); err != nil {
		os.Exit(-1)
	}
}
