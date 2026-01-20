package cmd

import (
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

const Version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:     "trelay",
	Short:   "Trelay - Developer-first URL shortener",
	Version: Version,
	Long: `Trelay is a developer-first, privacy-respecting URL manager.
Create, manage, and analyze shortened URLs from the command line.

Examples:
  trelay create https://example.com
  trelay create https://example.com --slug my-link
  trelay list
  trelay stats my-link
  trelay delete my-link`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json, csv)")
}
