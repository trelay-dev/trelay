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
	qrCopy   bool
)

var qrCmd = &cobra.Command{
	Use:   "qr <slug>",
	Short: "Generate QR code for a link",
	Long: `Generate a QR code image for a shortened link.

Examples:
  trelay qr my-link
  trelay qr my-link --output qr.png
  trelay qr my-link --size 512 --open
  trelay qr my-link --copy`,
	Args: cobra.ExactArgs(1),
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

		cfg, err := cli.LoadConfig()
		if err != nil {
			cli.Error(err.Error())
			return err
		}
		shortURL := fmt.Sprintf("%s/%s", cfg.APIURL, link.Slug)

		outputFile := qrOutput
		if outputFile == "" {
			outputFile = fmt.Sprintf("%s-qr.png", link.Slug)
		}

		if err := qrcode.WriteFile(shortURL, qrcode.Medium, qrSize, outputFile); err != nil {
			cli.Error(fmt.Sprintf("Failed to generate QR code: %v", err))
			return err
		}

		cli.Success(fmt.Sprintf("QR code saved to: %s", outputFile))

		if qrCopy {
			if err := copyToClipboard(outputFile); err != nil {
				cli.Error(fmt.Sprintf("Failed to copy to clipboard: %v", err))
			} else {
				cli.Success("QR code copied to clipboard")
			}
		}

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

func copyToClipboard(imagePath string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`set the clipboard to (read (POSIX file "%s") as TIFF picture)`, imagePath))
		return cmd.Run()
	case "linux":
		cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", "-i", imagePath)
		return cmd.Run()
	case "windows":
		// Windows clipboard for images requires more complex handling
		return fmt.Errorf("clipboard copy not supported on Windows")
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func init() {
	rootCmd.AddCommand(qrCmd)

	qrCmd.Flags().StringVarP(&qrOutput, "output", "o", "", "Output file path (default: <slug>-qr.png)")
	qrCmd.Flags().IntVarP(&qrSize, "size", "s", 256, "QR code size in pixels")
	qrCmd.Flags().BoolVar(&qrOpen, "open", false, "Open QR code after generation")
	qrCmd.Flags().BoolVar(&qrCopy, "copy", false, "Copy QR code to clipboard")
}
