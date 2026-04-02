# go-apod-bg

A CLI tool that fetches NASA's Astronomy Picture of the Day (APOD) and sets it as your desktop wallpaper. Works on Linux, macOS, and Windows.

## Requirements

- Go 1.25+
- A NASA API key (free at https://api.nasa.gov/) — a rate-limited `DEMO_KEY` is used as fallback

## Build

```bash
make build        # build for your platform
make all          # cross-compile for all platforms (output in dist/)
make clean        # remove build artifacts
```

## Usage

```bash
# Save your API key
./go-apod-bg configure

# Set today's APOD as wallpaper
./go-apod-bg update

# Use a specific date
./go-apod-bg update --date 2025-12-24

# Skip interactive prompts
./go-apod-bg update --auto

# Restore your previous wallpaper
./go-apod-bg restore
```

You can also pass `--api-key YOUR_KEY` or set the `NASA_API_KEY` environment variable instead of using `configure`.

## Platform Support

- **Linux** — GNOME, KDE, XFCE, MATE, Cinnamon, Sway, Hyprland, with fallback to feh/swaybg
- **macOS** — via AppleScript
- **Windows** — via Windows API

## Project Structure

```
cmd/           CLI command definitions (update, configure, restore)
internal/
  api/         NASA APOD API client and image downloader
  background/  Cross-platform wallpaper management
  config/      JSON config file handling
  utils/       Date parsing helpers
```

Images are cached in `~/Pictures/NASA/`. Configuration is stored in `~/.config/go-apod-bg/config.json` (Linux/macOS) or the equivalent user config directory on Windows.

## License

MIT
