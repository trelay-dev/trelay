package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	statsExport string
)

var statsCmd = &cobra.Command{
	Use:   "stats <slug>",
	Short: "View link statistics",
	Long: `View click statistics for a link.

Examples:
  trelay stats my-link
  trelay stats my-link -o json
  trelay stats my-link --export csv`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		stats, err := client.GetStats(args[0])
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		format := cli.OutputFormat(outputFormat)
		if statsExport != "" {
			format = cli.OutputFormat(statsExport)
		}

		return cli.PrintStats(stats, format)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().StringVar(&statsExport, "export", "", "Export format (json, csv)")
}
