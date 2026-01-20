package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var getCmd = &cobra.Command{
	Use:   "get <slug>",
	Short: "Get details of a link",
	Long: `Get detailed information about a specific link.

Examples:
  trelay get my-link
  trelay get my-link -o json`,
	Aliases: []string{"info", "show"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		link, err := client.GetLink(args[0])
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		return cli.PrintLink(link, cli.OutputFormat(outputFormat))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
