package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	deletePermanent bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <slug>",
	Short: "Delete a link",
	Long: `Delete a shortened link. By default, this is a soft delete.
Use --permanent to permanently remove the link.

Examples:
  trelay delete my-link
  trelay delete my-link --permanent`,
	Aliases: []string{"rm", "remove"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		if err := client.DeleteLink(args[0], deletePermanent); err != nil {
			cli.Error(err.Error())
			return err
		}

		action := "deleted"
		if deletePermanent {
			action = "permanently deleted"
		}

		cli.Success(fmt.Sprintf("Link '%s' %s", args[0], action))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deletePermanent, "permanent", false, "Permanently delete (cannot be restored)")
}
