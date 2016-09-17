package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Create the run command
var cmdRun = &cobra.Command{
	Use:   "run [COMMAND]",
	Short: "Run a command in the GlueX offline environment",
	Long: `
Run a command in the GlueX offline environment.

A bash shell is started by default.
Enclose multi-word commands in double quotes.

Usage examples:
1. hdpm run
2. hdpm run "hd_root -PPLUGINS=omega_hists file.evio"
3. hdpm run root
`,
	Run: runRun,
}

func init() {
	cmdHDPM.AddCommand(cmdRun)
}

func runRun(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "darwin" {
		fmt.Fprintln(os.Stderr, "Info: macOS is unsupported.")
		os.Exit(2)
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: Enclose multi-word commands in double quotes.")
		os.Exit(2)
	}
	setenv("CCDB_USER", os.Getenv("USER"))
	arg := "bash"
	if len(args) == 1 {
		arg = args[0]
	}
	env("")
	run("sh", "-c", arg)
}