//go:build !windows

package internal

func SetWindowProtection(hwnd uintptr, enable bool) bool {
	return false
}

func StartGlobalHotkeys(onAction func(action string)) {
}
