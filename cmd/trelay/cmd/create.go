package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	createSlug     string
	createDomain   string
	createPassword string
	createTTL      int
	createTags     []string
)

var createCmd = &cobra.Command{
	Use:   "create <url>",
	Short: "Create a shortened link",
	Long: `Create a new shortened link for the given URL.

Examples:
  trelay create https://example.com
  trelay create https://example.com --slug my-link
  trelay create https://example.com --password secret --ttl 24
  trelay create https://example.com --tags project,docs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		req := cli.CreateLinkRequest{
			URL:      args[0],
			Slug:     createSlug,
			Domain:   createDomain,
			Password: createPassword,
			TTLHours: createTTL,
			Tags:     createTags,
		}

		link, err := client.CreateLink(req)
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		if outputFormat == "json" {
			return cli.PrintLink(link, cli.OutputFormatJSON)
		}

		cli.Success(fmt.Sprintf("Created link: %s", link.Slug))
		fmt.Printf("Short URL:    %s/%s\n", "http://localhost:8080", link.Slug)
		fmt.Printf("Original URL: %s\n", link.OriginalURL)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createSlug, "slug", "s", "", "Custom slug for the link")
	createCmd.Flags().StringVarP(&createDomain, "domain", "d", "", "Custom domain for the link")
	createCmd.Flags().StringVarP(&createPassword, "password", "p", "", "Password protect the link")
	createCmd.Flags().IntVarP(&createTTL, "ttl", "t", 0, "Time-to-live in hours (0 = no expiration)")
	createCmd.Flags().StringSliceVar(&createTags, "tags", nil, "Tags for the link (comma-separated)")
}
