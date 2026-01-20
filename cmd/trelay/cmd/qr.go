package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var (
	qrOutput string
	qrSize   int
	qrOpen   bool
)

var qrCmd = &cobra.Command{
	Use:   "qr <slug>",
	Short: "Generate QR code for a link",
	Long: `Generate a QR code image for a shortened link.

Examples:
  trelay qr my-link
  trelay qr my-link --output qr.png
  trelay qr my-link --size 512 --open`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cli.GetClient()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		// Get link to verify it exists
		link, err := client.GetLink(args[0])
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		// Build short URL
		cfg, err := cli.LoadConfig()
		if err != nil {
			cli.Error(err.Error())
			return err
		}
		shortURL := fmt.Sprintf("%s/%s", cfg.APIURL, link.Slug)

		// Generate QR code
		outputFile := qrOutput
		if outputFile == "" {
			outputFile = fmt.Sprintf("%s-qr.png", link.Slug)
		}

		if err := qrcode.WriteFile(shortURL, qrcode.Medium, qrSize, outputFile); err != nil {
			cli.Error(fmt.Sprintf("Failed to generate QR code: %v", err))
			return err
		}

		cli.Success(fmt.Sprintf("QR code saved to: %s", outputFile))

		if qrOpen {
			if err := openFile(outputFile); err != nil {
				cli.Error(fmt.Sprintf("Failed to open file: %v", err))
			}
		}

		return nil
	},
}

func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(qrCmd)

	qrCmd.Flags().StringVarP(&qrOutput, "output", "o", "", "Output file path (default: <slug>-qr.png)")
	qrCmd.Flags().IntVarP(&qrSize, "size", "s", 256, "QR code size in pixels")
	qrCmd.Flags().BoolVar(&qrOpen, "open", false, "Open QR code after generation")
}
