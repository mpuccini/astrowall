//go:build windows

package background

import (
	"syscall"
	"unsafe"
)

func setWindows(filePath string) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("SystemParametersInfoW")

	pathPtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}

	const spiSetDeskWallpaper = 0x0014
	const spifUpdateINIFile = 0x01
	const spifSendChange = 0x02

	ret, _, callErr := proc.Call(
		uintptr(spiSetDeskWallpaper),
		uintptr(0),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(spifUpdateINIFile|spifSendChange),
	)
	if ret == 0 {
		return callErr
	}
	return nil
}
