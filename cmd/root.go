package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/marco/go-apod-bg/internal/api"
	"github.com/marco/go-apod-bg/internal/background"
	"github.com/marco/go-apod-bg/internal/config"
	"github.com/marco/go-apod-bg/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-apod-bg",
	Short: "NASA Astronomy Picture of the Day wallpaper setter",
	Long:  "Downloads the NASA Astronomy Picture of the Day and sets it as your desktop wallpaper.",
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Download today's APOD and set it as wallpaper",
	RunE:  runUpdate,
}

var (
	flagDate   string
	flagAuto   bool
	flagAPIKey string
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Save your NASA API key to the config file",
	RunE:  runConfigure,
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore the wallpaper that was set before the last update",
	RunE:  runRestore,
}

func init() {
	updateCmd.Flags().StringVar(&flagDate, "date", "", "date in YYYYMMDD or YYYY-MM-DD format (default: today)")
	updateCmd.Flags().BoolVar(&flagAuto, "auto", false, "skip all interactive prompts")
	updateCmd.Flags().StringVar(&flagAPIKey, "api-key", "", "NASA API key (or set NASA_API_KEY env var)")
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(restoreCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runConfigure(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	fmt.Print("NASA API key")
	if cfg.APIKey != "" {
		fmt.Printf(" (current: %s...%s)", cfg.APIKey[:4], cfg.APIKey[len(cfg.APIKey)-4:])
	}
	fmt.Print(": ")

	var key string
	fmt.Scanln(&key)
	if key == "" {
		fmt.Println("No changes made.")
		return nil
	}

	cfg.APIKey = key
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}
	fmt.Printf("API key saved to %s\n", config.Path())
	return nil
}

func runRestore(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	if cfg.PreviousWallpaper == "" {
		return fmt.Errorf("no previous wallpaper saved — run 'update' at least once first")
	}

	if _, err := os.Stat(cfg.PreviousWallpaper); err != nil {
		return fmt.Errorf("previous wallpaper no longer exists: %s", cfg.PreviousWallpaper)
	}

	if err := background.Set(cfg.PreviousWallpaper, true); err != nil {
		return err
	}

	fmt.Printf("Restored wallpaper: %s\n", cfg.PreviousWallpaper)
	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	date := time.Now()
	if flagDate != "" {
		var err error
		date, err = utils.ParseDate(flagDate)
		if err != nil {
			return fmt.Errorf("invalid date: %w", err)
		}
	}

	client := api.NewClient(flagAPIKey)

	filePath, err := client.DownloadImage(date)
	if err != nil {
		return err
	}

	// Save current wallpaper before changing it
	if prev := background.GetCurrent(); prev != "" {
		cfg := config.Load()
		cfg.PreviousWallpaper = prev
		config.Save(cfg)
	}

	if err := background.Set(filePath, flagAuto); err != nil {
		return err
	}

	fmt.Println("The image has successfully been set as background.")
	return nil
}
