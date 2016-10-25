package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Create the gxd command
var cmdGXD = &cobra.Command{
	Use:   "gxd [COMMAND] [ARGS]",
	Short: "A tool for managing docker builds of GlueX offline software",
	Long: `
gxd is a tool for managing docker builds of GlueX offline software.
`,
}

// Execute a gxd command
func Execute() {
	if err := cmdGXD.Execute(); err != nil {
		os.Exit(-1)
	}
}
