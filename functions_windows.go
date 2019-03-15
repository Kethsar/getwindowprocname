// +build windows

package main

import (
	"context"
	"log"
	"net"
	"regexp"
	"strings"

	pb "github.com/Kethsar/getwindowprocname/proto"
	"google.golang.org/grpc"

	"github.com/Kethsar/w32"
)

// dummy struct to implement the grpc server
type remoteProcServer struct{}

var (
	windows   []w32.HWND                          // For storing the window handles from a call to w32.EnumWindows()
	procRegex = regexp.MustCompile(`\\([^\\]+)$`) // To get the base process name
)

func (s *remoteProcServer) GetWindowInfo(ctx context.Context, cur *pb.Cursor) (*pb.WindowInfo, error) {
	winfo := getWindowInfo(int(cur.GetX()), int(cur.GetY()))
	return winfo, nil
}

func startServer() {
	if len(c.Port) < 2 { // colon and number
		log.Fatalln("Port in config either too short or missing")
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterRemoteProcServer(s, &remoteProcServer{})

	err = s.Serve(lis)
	if err != nil {
		log.Println(err)
	}
}

// Get the process name for the window current under the cursor
func getWindowInfo(x, y int) *pb.WindowInfo {
	procName := ""
	windows = make([]w32.HWND, 0, 10)
	winfo := &pb.WindowInfo{
		ClientRect:        &pb.Rect{},
		MonitorResolution: &pb.Resolution{},
	}
	defer clearWindows() // Immediately clear the windows array after the function returns to prevent keeping them in memory unnecessarily

	if x == -1 && y == -1 {
		cx, cy, ok := w32.GetCursorPos()

		if !ok {
			log.Println("Could not get cursor position")
			return winfo
		}
		x = cx
		y = cy
	}

	// Enum all windows with enumProc as the callback
	w32.EnumWindows(enumProc, 0)

	// Whittle down the array of window handles to only valid windows
	validWindows := make([]w32.HWND, 0, 10)
	for _, h := range windows {
		if IsValidWindow(h) {
			validWindows = append(validWindows, h)
		}
	}

	if len(validWindows) < 1 { // No real valid windows found, somehow
		return winfo
	}

	for _, h := range validWindows {
		rect := w32.GetWindowRect(h)

		if !CursorInRect(rect, x, y) {
			continue
		}

		// Get the process ID of the window, and attempt to open a handle to it
		_, procID := w32.GetWindowThreadProcessId(h)
		hProc, err := w32.OpenProcess(w32.PROCESS_QUERY_INFORMATION, false, uintptr(procID))
		if err != nil {
			continue
		}

		// Get the full process/program name, including path and extension
		procName = w32.QueryFullProcessImageName(hProc)
		w32.CloseHandle(hProc)

		if procName != "" {
			crect := w32.GetClientRect(h)
			cleft, ctop := w32.ClientToScreen(h, 0, 0)
			crect.Bottom += int32(ctop)
			crect.Right += int32(cleft)

			winfo.ClientRect.Left = int32(cleft)
			winfo.ClientRect.Top = int32(ctop)
			winfo.ClientRect.Bottom = crect.Bottom
			winfo.ClientRect.Right = crect.Right
			break
		}
	}

	// check if we found a process
	matches := procRegex.FindStringSubmatch(procName)
	if len(matches) < 1 {
		return winfo
	}

	// Extract the process/program name and remove its extension
	procName = matches[1]
	extIndex := strings.LastIndex(procName, ".")
	if extIndex > 0 { // If the only dot is the first character in the string, don't just blank the string
		procName = procName[:extIndex]
	}

	monWidth := w32.GetSystemMetrics(w32.SM_CXSCREEN)
	monHeight := w32.GetSystemMetrics(w32.SM_CYSCREEN)

	winfo.MonitorResolution.Width = int32(monWidth)
	winfo.MonitorResolution.Height = int32(monHeight)
	winfo.ProcName = procName

	return winfo
}

func enumProc(hwnd w32.HWND, lparam w32.LPARAM) w32.LRESULT {
	windows = append(windows, hwnd)
	return w32.LRESULT(1) // Something non-zero for true to continue enumeration
}

// Do the rectangle dimensions actually form a rectangle?
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

// Check if the window has an actual window open that is displayable in some way
func IsValidWindow(hwnd w32.HWND) bool {
	rect := w32.GetWindowRect(hwnd)

	return w32.IsWindowVisible(hwnd) &&
		!IsWindowCloaked(hwnd) &&
		(w32.GetWindowText(hwnd) != "") &&
		IsValidRect(rect)
}

func CursorInRect(rect *w32.RECT, x, y int) bool {
	return (y >= int(rect.Top) && y <= int(rect.Bottom)) &&
		(x >= int(rect.Left) && x <= int(rect.Right))
}

func clearWindows() {
	windows = nil
}
