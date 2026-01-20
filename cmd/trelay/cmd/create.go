package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	createSlug     string
	createDomain   string
	createPassword string
	createTTL      int
	createTags     []string
	createBulk     bool
)

var createCmd = &cobra.Command{
	Use:   "create <url>",
	Short: "Create a shortened link",
	Long: `Create a new shortened link for the given URL.

Examples:
  trelay create https://example.com
  trelay create https://example.com --slug my-link
  trelay create https://example.com --password secret --ttl 24
  trelay create https://example.com --tags project,docs
  trelay create https://example.com --domain short.example.com

Bulk create from stdin:
  cat urls.txt | trelay create --bulk
  echo -e "https://a.com\nhttps://b.com" | trelay create --bulk`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		if createBulk {
			return createBulkLinks(client)
		}

		if len(args) == 0 {
			return fmt.Errorf("URL is required (or use --bulk to read from stdin)")
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

		cfg, _ := cli.LoadConfig()
		baseURL := cfg.APIURL

		cli.Success(fmt.Sprintf("Created link: %s", link.Slug))
		fmt.Printf("Short URL:    %s/%s\n", baseURL, link.Slug)
		fmt.Printf("Original URL: %s\n", link.OriginalURL)

		return nil
	},
}

func createBulkLinks(client *cli.Client) error {
	scanner := bufio.NewScanner(os.Stdin)
	created := 0
	failed := 0

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" || strings.HasPrefix(url, "#") {
			continue
		}

		req := cli.CreateLinkRequest{
			URL:      url,
			Domain:   createDomain,
			Password: createPassword,
			TTLHours: createTTL,
			Tags:     createTags,
		}

		link, err := client.CreateLink(req)
		if err != nil {
			cli.Error(fmt.Sprintf("Failed to create link for %s: %v", url, err))
			failed++
			continue
		}

		if outputFormat == "json" {
			cli.PrintLink(link, cli.OutputFormatJSON)
		} else {
			fmt.Printf("%s -> %s\n", url, link.Slug)
		}
		created++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	if outputFormat != "json" {
		fmt.Printf("\nCreated: %d, Failed: %d\n", created, failed)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createSlug, "slug", "s", "", "Custom slug for the link")
	createCmd.Flags().StringVarP(&createDomain, "domain", "d", "", "Custom domain for the link")
	createCmd.Flags().StringVarP(&createPassword, "password", "p", "", "Password protect the link")
	createCmd.Flags().IntVarP(&createTTL, "ttl", "t", 0, "Time-to-live in hours (0 = no expiration)")
	createCmd.Flags().StringSliceVar(&createTags, "tags", nil, "Tags for the link (comma-separated)")
	createCmd.Flags().BoolVar(&createBulk, "bulk", false, "Read URLs from stdin (one per line)")
}
