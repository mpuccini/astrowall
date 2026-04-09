# astrowall

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
./astrowall configure

# Set today's APOD as wallpaper
./astrowall update

# Use a specific date
./astrowall update --date 2025-12-24

# Skip interactive prompts
./astrowall update --auto

# Restore your previous wallpaper
./astrowall restore
```

You can also pass `--api-key YOUR_KEY` or set the `NASA_API_KEY` environment variable instead of using `configure`.

## Installation

```bash
# Build and install binary + systemd timer
make install

# Or, if you already have the binary installed:
make install-systemd

# Remove everything
make uninstall
```

By default the binary is installed to `/usr/local/bin`. You can override this with `make install INSTALL_DIR=~/.local/bin`.

### Automatic daily wallpaper (Linux)

The `install-systemd` target sets up a systemd user timer that:

- Runs `astrowall update --auto` once per day
- Catches up after sleep/shutdown (`Persistent=true`)
- Retries with backoff on transient failures (network, API errors)
- Does not retry when there is nothing to do (e.g. APOD is a video)

Useful commands:

```bash
# Check timer status
systemctl --user status astrowall.timer

# Trigger a manual run
systemctl --user start astrowall.service

# View logs
journalctl --user -u astrowall.service
```

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
systemd/       Systemd service and timer unit files
```

Images are cached in `~/Pictures/NASA/`. Configuration is stored in `~/.config/astrowall/config.json` (Linux/macOS) or the equivalent user config directory on Windows.

## License

MIT
