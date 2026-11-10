//go:build windows

package internal

import (
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func HideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procSetWindowDisplayAffinity = user32.NewProc("SetWindowDisplayAffinity")
	procRegisterHotKey           = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey         = user32.NewProc("UnregisterHotKey")
	procGetMessage               = user32.NewProc("GetMessageW")
)

const (
	WDA_NONE               = 0x00000000
	WDA_MONITOR            = 0x00000001
	WDA_EXCLUDEFROMCAPTURE = 0x00000011

	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
	MOD_NOREPEAT = 0x4000

	VK_RETURN = 0x0D
	VK_H      = 0x48
	VK_OEM_2  = 0xBF // '? /' key
	VK_X      = 0x58

	WM_HOTKEY = 0x0312
)

func SetWindowProtection(hwnd uintptr, enable bool) bool {
	affinity := uintptr(WDA_NONE)
	if enable {
		affinity = uintptr(WDA_EXCLUDEFROMCAPTURE)
	}
	r, _, _ := procSetWindowDisplayAffinity.Call(hwnd, affinity)
	if r == 0 && enable {
		// Fallback for older Win10 (< 2004)
		r, _, _ = procSetWindowDisplayAffinity.Call(hwnd, uintptr(WDA_MONITOR))
	}
	return r != 0
}

func StartGlobalHotkeys(onAction func(action string)) {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// 1: Ctrl+Shift+Enter (Smart Assist)
		procRegisterHotKey.Call(0, 1, MOD_CONTROL|MOD_SHIFT|MOD_NOREPEAT, VK_RETURN)
		// 2: Ctrl+Enter (What should I say?)
		procRegisterHotKey.Call(0, 2, MOD_CONTROL|MOD_NOREPEAT, VK_RETURN)
		// 3: Ctrl+H (LeetCode solve)
		procRegisterHotKey.Call(0, 3, MOD_CONTROL|MOD_NOREPEAT, VK_H)
		// 4: Ctrl+Shift+/ (Toggle hide)
		procRegisterHotKey.Call(0, 4, MOD_CONTROL|MOD_SHIFT|MOD_NOREPEAT, VK_OEM_2)
		// 5: Ctrl+Shift+X (Quit)
		procRegisterHotKey.Call(0, 5, MOD_CONTROL|MOD_SHIFT|MOD_NOREPEAT, VK_X)

		var msg struct {
			hwnd    uintptr
			message uint32
			wParam  uintptr
			lParam  uintptr
			time    uint32
			pt      struct{ x, y int32 }
		}

		for {
			r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if int32(r) <= 0 {
				break
			}
			if msg.message == WM_HOTKEY {
				switch msg.wParam {
				case 1:
					onAction("assist")
				case 2:
					onAction("say")
				case 3:
					onAction("leetcode")
				case 4:
					onAction("toggle_hide")
				case 5:
					onAction("quit")
				}
			}
		}
	}()
}
