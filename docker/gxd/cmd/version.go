package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Create the version command
var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the gxd version number",
	Long: `
Show the gxd version number.
`,
	Run: runVersion,
}

const VERSION = "dev"

func init() {
	cmdGXD.AddCommand(cmdVersion)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("gxd version %s\n", VERSION)
}
