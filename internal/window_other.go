//go:build !windows

package internal

import "os/exec"

func HideWindow(cmd *exec.Cmd) {}

func SetWindowProtection(hwnd uintptr, enable bool) bool {
	return false
}

func StartGlobalHotkeys(onAction func(action string)) {
}
