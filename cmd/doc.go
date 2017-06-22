package cmd

import (
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Create the doc command
var cmdDoc = &cobra.Command{
	Use:   "doc",
	Short: "Generate Markdown documentation for hdpm",
	Long: `Generate Markdown documentation for hdpm.

Markdown files are written to the doc directory by default.
These files form the Commands documentation on the GitHub wiki.
https://github.com/JeffersonLab/hdpm/wiki/hdpm`,
	Example: `1. hdpm doc
2. hdpm doc -d mydocs`,
	Hidden: true,
	Run:    runDoc,
}

var docDir string

func init() {
	cmdHDPM.AddCommand(cmdDoc)

	cmdDoc.Flags().StringVarP(&docDir, "dir", "d", "doc", "Documentation directory")
}

func runDoc(cmd *cobra.Command, args []string) {
	cmd.Root().DisableAutoGenTag = true
	filePrepender := func(s string) string { return "" }
	linkHandler := func(name string) string {
		url := "https://github.com/JeffersonLab/hdpm/wiki/"
		base := strings.TrimSuffix(name, path.Ext(name))
		return url + strings.Replace(base, "_", "-", 1)
	}
	mk(docDir)
	doc.GenMarkdownTreeCustom(cmd.Root(), docDir, filePrepender, linkHandler)
}
