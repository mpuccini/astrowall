package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configDir = "astrowall"
const configFile = "config.json"

// Config holds persistent application settings.
type Config struct {
	APIKey             string `json:"api_key"`
	PreviousWallpaper  string `json:"previous_wallpaper,omitempty"`
}

// Path returns the full path to the config file.
func Path() string {
	dir, _ := os.UserConfigDir() // ~/.config on Linux, ~/Library/Application Support on macOS
	return filepath.Join(dir, configDir, configFile)
}

// Load reads the config file. Returns a zero Config (not an error) if the file
// doesn't exist yet, so callers can always fall through to defaults.
func Load() Config {
	var cfg Config
	data, err := os.ReadFile(Path())
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg)
	return cfg
}

// Save writes the config to disk, creating the directory if needed.
func Save(cfg Config) error {
	p := Path()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}
