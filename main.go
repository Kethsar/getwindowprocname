package main

import (
	"log"

	"github.com/Kethsar/w32"
)

var windows []w32.HWND // For storing the window handles from a call to w32.EnumWindows()

func main() {
	windows = make([]w32.HWND, 0, 10)
	x, y, ok := w32.GetCursorPos()
	if !ok {
		log.Fatalln("Could not get cursor position")
	}

	w32.EnumWindows(enumProc, 0)

	// Whittle down the array of window handles to only valid windows
	validWindows := make([]w32.HWND, 0, 10)
	for _, h := range windows {
		if IsValidWindow(h) {
			validWindows = append(validWindows, h)
		}
	}

	// iterate array calling GetWindowRect
	// Check rect for mouse position
	for _, h := range validWindows {
		rect := w32.GetWindowRect(h)
		curIn := (y >= int(rect.Top) && y <= int(rect.Bottom)) &&
			(x >= int(rect.Left) && x <= int(rect.Right))

		if !curIn {
			continue
		}

		_, procID := w32.GetWindowThreadProcessId(h)
		hProc, err := w32.OpenProcess(w32.PROCESS_QUERY_INFORMATION, false, uintptr(procID))
		if err != nil {
			log.Println(err)
			continue
		}

		procName := w32.QueryFullProcessImageName(hProc)
		w32.CloseHandle(hProc)

		log.Println("Handle:", h, "; Process ID:", procID, "; Process Name:", procName, "; RECT:", rect)
	}
}

func enumProc(hwnd w32.HWND, lparam w32.LPARAM) w32.LRESULT {
	windows = append(windows, hwnd)
	return w32.LRESULT(1) // Something non-zero for true to continue enumeration
}

func IsValidRect(r *w32.RECT) bool {
	return (r.Bottom-r.Top) > 0 && (r.Right-r.Left) > 0
}

func IsWindowCloaked(hwnd w32.HWND) bool {
	var dwmEnabled w32.BOOL

	// I think technically this should only be called on Vista or higher
	// But this particular program is only meant for W10 so don't bother checking
	w32.DwmIsCompositionEnabled(&dwmEnabled)
	if dwmEnabled == 0 {
		return false
	}

	cloaked, ret := w32.DwmGetWindowAttribute(hwnd, w32.DWMWA_CLOAKED)

	// The type assertion on cloaked could be dangerous
	// Except given the attribute we are checking, it can only be a *w32.DWORD
	return (ret == 0) &&
		(*(cloaked.(*w32.DWORD)) != 0)
}

func IsValidWindow(hwnd w32.HWND) bool {
	rect := w32.GetWindowRect(hwnd)

	return w32.IsWindowVisible(hwnd) &&
		!IsWindowCloaked(hwnd) &&
		(w32.GetWindowText(hwnd) != "") &&
		IsValidRect(rect)
}
