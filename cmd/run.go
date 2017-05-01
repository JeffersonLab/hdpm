package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Create the run command
var cmdRun = &cobra.Command{
	Use:   "run [COMMAND]",
	Short: "Run a command in the GlueX offline environment",
	Long: `Run a command in the GlueX offline environment.

A bash shell is started if no command is given.
Enclose multi-word commands in quotes.`,
	Example: `1. hdpm run
2. hdpm run "hd_root -PPLUGINS=omega_hists file.evio"
3. hdpm run root`,
	Run: runRun,
}

func init() {
	cmdHDPM.AddCommand(cmdRun)
}

func runRun(cmd *cobra.Command, args []string) {
	if runtime.GOOS == "darwin" {
		fmt.Fprintln(os.Stderr, "run: macOS is unsupported (due to SIP).")
		os.Exit(2)
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments: Enclose multi-word commands in quotes.\n")
		cmd.Usage()
		os.Exit(2)
	}
	pkgInit()
	setEnv()
	arg := "bash"
	if len(args) == 1 {
		arg = args[0]
	}
	c := command("sh", "-c", arg)
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		log.Fatalln(err)
	}
}
