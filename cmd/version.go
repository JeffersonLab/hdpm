package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Create the version command
var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the hdpm version number",
	Long: `
Show the hdpm version number.
`,
	Run: runVersion,
}

const VERSION = "0.9.0"

func init() {
	cmdHDPM.AddCommand(cmdVersion)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("hdpm version %s\n", VERSION)
}