//go:build !windows

package background

import "fmt"

func setWindows(filePath string) error {
	return fmt.Errorf("windows background setting is not supported on this platform")
}
