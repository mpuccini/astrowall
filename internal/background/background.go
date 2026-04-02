package background

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const minFileSizeKB = 500

// GetCurrent returns the current wallpaper path, or empty string if it can't
// be determined. The returned path has the "file://" prefix stripped.
func GetCurrent() string {
	de := strings.ToLower(os.Getenv("DESKTOP_SESSION"))
	if de == "" {
		de = strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	}

	var out []byte
	var err error

	switch runtime.GOOS {
	case "linux":
		switch {
		case strings.Contains(de, "gnome"),
			strings.Contains(de, "ubuntu"),
			strings.Contains(de, "unity"),
			strings.Contains(de, "awesome-gnome"):
			out, err = exec.Command("gsettings", "get",
				"org.gnome.desktop.background", "picture-uri").Output()

		case strings.Contains(de, "mate"):
			out, err = exec.Command("gsettings", "get",
				"org.mate.background", "picture-filename").Output()

		case strings.Contains(de, "cinnamon"):
			out, err = exec.Command("gsettings", "get",
				"org.cinnamon.desktop.background", "picture-uri").Output()

		case strings.Contains(de, "xfce"):
			out, err = exec.Command("xfconf-query", "-c", "xfce4-desktop",
				"-p", "/backdrop/screen0/monitor0/workspace0/last-image").Output()
		}

	case "darwin":
		out, err = exec.Command("osascript", "-e",
			`tell application "System Events" to get picture of current desktop`).Output()
	}

	if err != nil || len(out) == 0 {
		return ""
	}

	path := strings.TrimSpace(string(out))
	path = strings.Trim(path, "'\"")
	path = strings.TrimPrefix(path, "file://")
	return path
}

// Set sets the desktop wallpaper to the given file path.
// If the file is smaller than 500KB: in interactive mode the user is prompted,
// in auto mode the image is skipped.
func Set(filePath string, auto bool) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("the background file was not found: %w", err)
	}

	sizeKB := info.Size() / 1024
	if sizeKB < minFileSizeKB {
		if auto {
			fmt.Printf("Skipping: image is too small (%d kb, minimum %d kb).\n", sizeKB, minFileSizeKB)
			return nil
		}
		fmt.Printf("The file size of this background is relatively small (%d kb), are you sure you want to continue? [y/N] ", sizeKB)
		var answer string
		fmt.Scanln(&answer)
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	switch runtime.GOOS {
	case "linux":
		return setLinux(filePath)
	case "darwin":
		return setMac(filePath)
	case "windows":
		return setWindows(filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func setLinux(filePath string) error {
	de := strings.ToLower(os.Getenv("DESKTOP_SESSION"))
	if de == "" {
		de = strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	}

	switch {
	case strings.Contains(de, "gnome"),
		strings.Contains(de, "ubuntu"),
		strings.Contains(de, "unity"),
		strings.Contains(de, "awesome-gnome"):
		uri := "file://" + filePath
		// Set for both light and dark themes
		if err := exec.Command("gsettings", "set",
			"org.gnome.desktop.background", "picture-uri", uri).Run(); err != nil {
			return fmt.Errorf("gsettings failed: %w", err)
		}
		// GNOME 42+ also uses picture-uri-dark
		_ = exec.Command("gsettings", "set",
			"org.gnome.desktop.background", "picture-uri-dark", uri).Run()
		return nil

	case strings.Contains(de, "kde"),
		strings.Contains(de, "plasma"):
		script := fmt.Sprintf(`
var allDesktops = desktops();
for (i=0;i<allDesktops.length;i++) {
    d = allDesktops[i];
    d.wallpaperPlugin = "org.kde.image";
    d.currentConfigGroup = Array("Wallpaper", "org.kde.image", "General");
    d.writeConfig("Image", "file://%s")
}`, filePath)
		return exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell",
			"org.kde.PlasmaShell.evaluateScript", script).Run()

	case strings.Contains(de, "xfce"):
		return exec.Command("xfconf-query", "-c", "xfce4-desktop",
			"-p", "/backdrop/screen0/monitor0/workspace0/last-image",
			"-s", filePath).Run()

	case strings.Contains(de, "mate"):
		return exec.Command("gsettings", "set",
			"org.mate.background", "picture-filename", filePath).Run()

	case strings.Contains(de, "cinnamon"):
		uri := "file://" + filePath
		return exec.Command("gsettings", "set",
			"org.cinnamon.desktop.background", "picture-uri", uri).Run()

	case strings.Contains(de, "sway"):
		return exec.Command("swaymsg", "output", "*", "bg", filePath, "fill").Run()

	case strings.Contains(de, "hyprland"):
		return exec.Command("hyprctl", "hyprpaper", "wallpaper", ","+filePath).Run()

	default:
		// Fallback to feh (works for i3, awesome, spectrwm, wmii, etc.)
		if _, err := exec.LookPath("feh"); err == nil {
			return exec.Command("feh", "--bg-center", filePath).Run()
		}
		// Try swaybg as last resort
		if _, err := exec.LookPath("swaybg"); err == nil {
			return exec.Command("swaybg", "-i", filePath, "-m", "fill").Run()
		}
		return fmt.Errorf("unable to change the background: could not detect desktop environment (DESKTOP_SESSION=%q) and no fallback tool (feh, swaybg) found", de)
	}
}

func setMac(filePath string) error {
	script := fmt.Sprintf(`tell application "System Events" to tell every desktop to set picture to "%s"`, filePath)
	return exec.Command("osascript", "-e", script).Run()
}
