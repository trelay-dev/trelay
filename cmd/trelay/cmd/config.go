package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/aftaab/trelay/internal/cli"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `View and modify CLI configuration settings.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  api-url    The Trelay API URL (default: http://localhost:8080)
  api-key    Your API key for authentication

Examples:
  trelay config set api-url https://short.example.com
  trelay config set api-key tr_abc123...`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := cli.LoadConfig()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		switch key {
		case "api-url":
			cfg.APIURL = value
		case "api-key":
			cfg.APIKey = value
		default:
			cli.Error(fmt.Sprintf("Unknown configuration key: %s", key))
			return fmt.Errorf("unknown key: %s", key)
		}

		if err := cli.SaveConfig(cfg); err != nil {
			cli.Error(err.Error())
			return err
		}

		cli.Success(fmt.Sprintf("Configuration '%s' updated", key))
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Get a configuration value.

Examples:
  trelay config get api-url
  trelay config get api-key`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg, err := cli.LoadConfig()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		var value string
		switch key {
		case "api-url":
			value = cfg.APIURL
		case "api-key":
			if cfg.APIKey != "" {
				value = cfg.APIKey[:8] + "..." // Mask API key
			} else {
				value = "(not set)"
			}
		default:
			cli.Error(fmt.Sprintf("Unknown configuration key: %s", key))
			return fmt.Errorf("unknown key: %s", key)
		}

		fmt.Println(value)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration",
	Long:  `Display all current configuration settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cli.LoadConfig()
		if err != nil {
			cli.Error(err.Error())
			return err
		}

		fmt.Printf("API URL: %s\n", cfg.APIURL)

		if cfg.APIKey != "" {
			fmt.Printf("API Key: %s...\n", cfg.APIKey[:8])
		} else {
			fmt.Println("API Key: (not set)")
		}

		path, _ := cli.ConfigPath()
		fmt.Printf("Config:  %s\n", path)

		// System Check for Linux CLI features
		if runtime.GOOS == "linux" {
			fmt.Println("\nSystem Check (Linux):")
			_, wlErr := exec.LookPath("wl-copy")
			_, xclipErr := exec.LookPath("xclip")
			if wlErr != nil && xclipErr != nil {
				fmt.Println("  [!] Warning: Neither 'wl-copy' nor 'xclip' found.")
				fmt.Println("      The --copy flag in 'qr' command will not work.")
				fmt.Println("      Fix: sudo apt install xclip  # or wl-clipboard")
			} else {
				fmt.Println("  [✓] Clipboard support is ready.")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configShowCmd)
}
