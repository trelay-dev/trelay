package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	listSearch string
	listTags   []string
	listFolder int64
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all links",
	Long: `List all shortened links with optional filtering.

Examples:
  trelay list
  trelay list --search example
  trelay list --tags project,docs
  trelay list --folder 1
  trelay list --limit 10 --offset 20`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		opts := cli.ListLinksOptions{
			Search: listSearch,
			Tags:   listTags,
			Limit:  listLimit,
			Offset: listOffset,
		}

		if cmd.Flags().Changed("folder") {
			opts.FolderID = &listFolder
		}

		links, err := client.ListLinks(opts)
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		return cli.PrintLinks(links, cli.OutputFormat(outputFormat))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listSearch, "search", "s", "", "Search in slug and URL")
	listCmd.Flags().StringSliceVar(&listTags, "tags", nil, "Filter by tags (comma-separated)")
	listCmd.Flags().Int64VarP(&listFolder, "folder", "f", 0, "Filter by folder ID")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 50, "Maximum number of results")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "Offset for pagination")
}
