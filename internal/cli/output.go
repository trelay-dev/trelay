package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

// OutputFormat defines the output format type.
type OutputFormat string

const (
	OutputFormatTable OutputFormat = "table"
	OutputFormatJSON  OutputFormat = "json"
	OutputFormatCSV   OutputFormat = "csv"
)

// PrintLinks outputs links in the specified format.
func PrintLinks(links []Link, format OutputFormat) error {
	switch format {
	case OutputFormatJSON:
		return printJSON(links)
	case OutputFormatCSV:
		return printLinksCSV(links)
	default:
		return printLinksTable(links)
	}
}

// PrintLink outputs a single link.
func PrintLink(link *Link, format OutputFormat) error {
	switch format {
	case OutputFormatJSON:
		return printJSON(link)
	default:
		return printLinkDetails(link)
	}
}

// PrintStats outputs statistics.
func PrintStats(stats *ClickStats, format OutputFormat) error {
	switch format {
	case OutputFormatJSON:
		return printJSON(stats)
	case OutputFormatCSV:
		return printStatsCSV(stats)
	default:
		return printStatsTable(stats)
	}
}

func printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func printLinksTable(links []Link) error {
	if len(links) == 0 {
		fmt.Println("No links found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SLUG\tURL\tCLICKS\tCREATED\tTAGS")
	fmt.Fprintln(w, "----\t---\t------\t-------\t----")

	for _, l := range links {
		url := l.OriginalURL
		if len(url) > 50 {
			url = url[:47] + "..."
		}

		tags := ""
		if len(l.Tags) > 0 {
			for i, tag := range l.Tags {
				if i > 0 {
					tags += ", "
				}
				tags += tag
			}
		}

		created := l.CreatedAt
		if len(created) > 10 {
			created = created[:10]
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			l.Slug, url, l.ClickCount, created, tags)
	}

	return w.Flush()
}

func printLinkDetails(link *Link) error {
	fmt.Printf("Slug:        %s\n", link.Slug)
	fmt.Printf("URL:         %s\n", link.OriginalURL)
	fmt.Printf("Clicks:      %d\n", link.ClickCount)
	fmt.Printf("Created:     %s\n", link.CreatedAt)
	fmt.Printf("Updated:     %s\n", link.UpdatedAt)

	if link.Domain != "" {
		fmt.Printf("Domain:      %s\n", link.Domain)
	}

	if link.HasPassword {
		fmt.Printf("Password:    Yes\n")
	}

	if link.ExpiresAt != nil {
		fmt.Printf("Expires:     %s\n", *link.ExpiresAt)
	}

	if len(link.Tags) > 0 {
		fmt.Printf("Tags:        %v\n", link.Tags)
	}

	return nil
}

func printLinksCSV(links []Link) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	headers := []string{"slug", "url", "clicks", "created_at", "tags"}
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, l := range links {
		tags := ""
		for i, tag := range l.Tags {
			if i > 0 {
				tags += ";"
			}
			tags += tag
		}

		record := []string{
			l.Slug,
			l.OriginalURL,
			strconv.FormatInt(l.ClickCount, 10),
			l.CreatedAt,
			tags,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func printStatsTable(stats *ClickStats) error {
	fmt.Printf("Total Clicks: %d\n\n", stats.TotalClicks)

	if len(stats.ClicksByDay) > 0 {
		fmt.Println("Clicks by Day:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tCLICKS")
		for _, d := range stats.ClicksByDay {
			fmt.Fprintf(w, "%s\t%d\n", d.Date, d.Clicks)
		}
		w.Flush()
		fmt.Println()
	}

	if len(stats.TopReferrers) > 0 {
		fmt.Println("Top Referrers:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "REFERRER\tCLICKS")
		for _, r := range stats.TopReferrers {
			ref := r.Referrer
			if len(ref) > 50 {
				ref = ref[:47] + "..."
			}
			fmt.Fprintf(w, "%s\t%d\n", ref, r.Clicks)
		}
		w.Flush()
	}

	return nil
}

func printStatsCSV(stats *ClickStats) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{"total_clicks", strconv.FormatInt(stats.TotalClicks, 10)})
	w.Write([]string{})
	w.Write([]string{"date", "clicks"})

	for _, d := range stats.ClicksByDay {
		w.Write([]string{d.Date, strconv.FormatInt(d.Clicks, 10)})
	}

	return nil
}

func PrintFolders(folders []Folder, format OutputFormat) error {
	switch format {
	case OutputFormatJSON:
		return printJSON(folders)
	case OutputFormatCSV:
		return printFoldersCSV(folders)
	default:
		return printFoldersTable(folders)
	}
}

func printFoldersTable(folders []Folder) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPARENT\tCREATED")
	fmt.Fprintln(w, "--\t----\t------\t-------")

	for _, f := range folders {
		parent := "-"
		if f.ParentID != nil {
			parent = strconv.FormatInt(*f.ParentID, 10)
		}

		created := f.CreatedAt
		if len(created) > 10 {
			created = created[:10]
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", f.ID, f.Name, parent, created)
	}

	return w.Flush()
}

func printFoldersCSV(folders []Folder) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{"id", "name", "parent_id", "created_at"})
	for _, f := range folders {
		parent := ""
		if f.ParentID != nil {
			parent = strconv.FormatInt(*f.ParentID, 10)
		}
		w.Write([]string{strconv.FormatInt(f.ID, 10), f.Name, parent, f.CreatedAt})
	}

	return nil
}

func Success(message string) {
	fmt.Printf("✓ %s\n", message)
}

func Warning(message string) {
	fmt.Printf("⚠ %s\n", message)
}

func Error(message string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", message)
}
