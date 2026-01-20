package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	deletePermanent bool
	deleteBulk      bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <slug> [slug2...]",
	Short: "Delete one or more links",
	Long: `Delete shortened links. By default, this is a soft delete.
Use --permanent to permanently remove the links.
Use --bulk with comma-separated slugs for bulk deletion.

Examples:
  trelay delete my-link
  trelay delete my-link --permanent
  trelay delete link1 link2 link3
  trelay delete --bulk link1,link2,link3`,
	Aliases: []string{"rm", "remove"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		var slugs []string
		if deleteBulk && len(args) == 1 {
			slugs = strings.Split(args[0], ",")
			for i := range slugs {
				slugs[i] = strings.TrimSpace(slugs[i])
			}
		} else {
			slugs = args
		}

		if len(slugs) == 1 {
			if err := client.DeleteLink(slugs[0], deletePermanent); err != nil {
				cli.Error(err.Error())
				return err
			}

			action := "deleted"
			if deletePermanent {
				action = "permanently deleted"
			}
			cli.Success(fmt.Sprintf("Link '%s' %s", slugs[0], action))
		} else {
			result, err := client.BulkDeleteLinks(slugs, deletePermanent)
			if err != nil {
				cli.Error(err.Error())
				return err
			}

			if len(result.Deleted) > 0 {
				cli.Success(fmt.Sprintf("Deleted %d links: %s", len(result.Deleted), strings.Join(result.Deleted, ", ")))
			}
			if len(result.Failed) > 0 {
				cli.Error(fmt.Sprintf("Failed to delete %d links: %s", len(result.Failed), strings.Join(result.Failed, ", ")))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deletePermanent, "permanent", false, "Permanently delete (cannot be restored)")
	deleteCmd.Flags().BoolVar(&deleteBulk, "bulk", false, "Enable bulk deletion with comma-separated slugs")
}
