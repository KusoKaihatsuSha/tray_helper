package clipwindow

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"

	"golang.org/x/sys/windows"
)

const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260

	unicodeTextFormat   = 13 // CF_UNICODETEXT
	fixedMemory         = 0x0000
	normalShow          = 9
	aOpen               = "OpenClipboard"
	aClose              = "CloseClipboard"
	aClear              = "EmptyClipboard"
	aGet                = "GetClipboardData"
	aSet                = "SetClipboardData"
	aAlloc              = "GlobalAlloc"
	aFree               = "GlobalFree"
	aLock               = "GlobalLock"
	aUnlock             = "GlobalUnlock"
	aMove               = "RtlMoveMemory"
	aWindowName         = "GetWindowTextW"
	aWindowFname        = "GetForegroundWindow"
	aWindowNameLen      = "GetWindowTextLengthW"
	aWindowEnum         = "EnumWindows"
	aCurrentWindowId    = "GetCurrentThreadId"
	aWindowId           = "GetWindowThreadProcessId"
	aWindowVisible      = "IsWindowVisible"
	aWindowGet          = "GetWindow"
	aWindowTop          = "SetForegroundWindow"
	aWindowShow         = "ShowWindow" // "ShowWindowAsync"
	aSetActiveWindow    = "SetActiveWindow"
	aSetFocus           = "SetFocus"
	aAnimateWindow      = "AnimateWindow"
	aAttachThreadInput  = "AttachThreadInput"
	aBringWindowToTop   = "BringWindowToTop"
	aKeybdEvent         = "keybd_event"
	aEnableWindow       = "EnableWindow"
	aGetWindowRect      = "GetWindowRect"
	aGetWindowInfo      = "GetWindowInfo"
	aMouseEvent         = "mouse_event"
	aFlashWindow        = "FlashWindow"
	aSwitchToThisWindow = "SwitchToThisWindow"
	aOpenIcon           = "OpenIcon"
	aCloseWindow        = "CloseWindow"
	aGetSystemMetrics   = "GetSystemMetrics"

	aSetWindowPos = "SetWindowPos"

	aProcClose                = "CloseHandle"
	aCreateToolhelp32Snapshot = "CreateToolhelp32Snapshot"
	aProcess32First           = "Process32FirstW"
	aProcess32Next            = "Process32NextW"
	aOpenProcess              = "GetProcessId"
)

var (
	user32    = syscall.MustLoadDLL("user32")
	kernel32  = syscall.NewLazyDLL("kernel32")
	clipboard = make(map[string]*syscall.Proc, 5)
	memory    = make(map[string]*syscall.LazyProc, 5)
	window    = make(map[string]*syscall.Proc, 3)
	procss    = make(map[string]*syscall.LazyProc, 4)
)

// Process of exec
type Process struct {
	exe  string
	pid  int
	ppid int
	h    uintptr
}

// PROCESSENTRY32 - process of exec (win handler)
type PROCESSENTRY32 struct {
	ExeFile           [MAX_PATH]uint16
	PriorityClassBase int32
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	Flags             uint32
	DefaultHeapID     uintptr
}

// RECT - process of exec (win handler)
type RECT struct {
	Left, Top, Right, Bottom int32
}

type HWND unsafe.Pointer
type DWORD uint32
type UINT uint32
type WORD uint16
type ATOM WORD

// WINDOWINFO - process of exec (win handler)
type WINDOWINFO struct {
	CbSize          DWORD
	RcWindow        RECT
	RcClient        RECT
	DwStyle         DWORD
	DwExStyle       DWORD
	DwWindowStatus  DWORD
	CxWindowBorders UINT
	CyWindowBorders UINT
	AtomWindowType  ATOM
	WCreatorVersion WORD
}

type POINT struct {
	X, Y int32
}

// WinWin - Window in WindowsOS
type WinWin struct {
	Title     string
	ProcessID uintptr
	Handle    uintptr
}

// Procs - processes interface
type Procs interface {
	*syscall.Proc | *syscall.LazyProc
	Call(...uintptr) (uintptr, uintptr, error)
}

func init() {
	clipboard[aOpen] = user32.MustFindProc(aOpen)
	clipboard[aClose] = user32.MustFindProc(aClose)
	clipboard[aClear] = user32.MustFindProc(aClear)
	clipboard[aGet] = user32.MustFindProc(aGet)
	clipboard[aSet] = user32.MustFindProc(aSet)
	memory[aAlloc] = kernel32.NewProc(aAlloc)
	memory[aFree] = kernel32.NewProc(aFree)
	memory[aLock] = kernel32.NewProc(aLock)
	memory[aUnlock] = kernel32.NewProc(aUnlock)
	memory[aMove] = kernel32.NewProc(aMove)
	window[aWindowName] = user32.MustFindProc(aWindowName)
	window[aWindowNameLen] = user32.MustFindProc(aWindowNameLen)
	window[aWindowFname] = user32.MustFindProc(aWindowFname)
	window[aWindowEnum] = user32.MustFindProc(aWindowEnum)
	window[aWindowId] = user32.MustFindProc(aWindowId)
	window[aWindowVisible] = user32.MustFindProc(aWindowVisible)
	window[aWindowGet] = user32.MustFindProc(aWindowGet)
	window[aWindowTop] = user32.MustFindProc(aWindowTop)
	window[aWindowShow] = user32.MustFindProc(aWindowShow)
	window[aAnimateWindow] = user32.MustFindProc(aAnimateWindow)
	window[aSetActiveWindow] = user32.MustFindProc(aSetActiveWindow)
	window[aSetFocus] = user32.MustFindProc(aSetFocus)
	procss[aCurrentWindowId] = kernel32.NewProc(aCurrentWindowId)
	window[aAttachThreadInput] = user32.MustFindProc(aAttachThreadInput)
	window[aSetWindowPos] = user32.MustFindProc(aSetWindowPos)
	window[aBringWindowToTop] = user32.MustFindProc(aBringWindowToTop)
	window[aKeybdEvent] = user32.MustFindProc(aKeybdEvent)
	window[aEnableWindow] = user32.MustFindProc(aEnableWindow)
	window[aGetWindowRect] = user32.MustFindProc(aGetWindowRect)
	window[aGetWindowInfo] = user32.MustFindProc(aGetWindowInfo)
	window[aMouseEvent] = user32.MustFindProc(aMouseEvent)
	window[aFlashWindow] = user32.MustFindProc(aFlashWindow)
	window[aSwitchToThisWindow] = user32.MustFindProc(aSwitchToThisWindow)
	window[aOpenIcon] = user32.MustFindProc(aOpenIcon)
	window[aCloseWindow] = user32.MustFindProc(aCloseWindow)
	window[aGetSystemMetrics] = user32.MustFindProc(aGetSystemMetrics)
	procss[aProcClose] = kernel32.NewProc(aProcClose)
	procss[aCreateToolhelp32Snapshot] = kernel32.NewProc(aCreateToolhelp32Snapshot)
	procss[aProcess32First] = kernel32.NewProc(aProcess32First)
	procss[aProcess32Next] = kernel32.NewProc(aProcess32Next)

	procss[aOpenProcess] = kernel32.NewProc(aOpenProcess)
}

// call proc and print error
func call[T Procs](o map[string]T, procName string, p ...uintptr) uintptr {
	result, _, err := o[procName].Call(p...)
	if result == 0 && err != nil {
		checkErr := helpers.FindError[syscall.Errno](err)
		if checkErr != nil && !errors.As(checkErr, new(syscall.Errno)) {
			// exclude "The operation completed successfully"
			helpers.ToLog(fmt.Sprintf("%s null result of calling: %v\n", procName, err))
		}
	}
	return result
}

// paramAlloc calculate alloc data
func paramAlloc(text string) (uintptr, uintptr) {
	newText, err := syscall.UTF16FromString(text)
	if err != nil {
		return 0, 0
	}
	needAlloc := 0
	for _, v := range newText {
		needAlloc += int(unsafe.Sizeof(v))
	}
	return uintptr(unsafe.Pointer(&newText[0])), uintptr(needAlloc)
}

// Get get the clipboard data
func Get() string {
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	call(clipboard, aOpen, 0)
	defer call(clipboard, aClose)
	mem := call(clipboard, aGet, unicodeTextFormat)
	p := call(memory, aLock, mem)
	defer call(memory, aUnlock, mem)
	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(p))) // irresistibly linter
}

// Clear the clipboard data
func Clear() {
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	call(clipboard, aOpen, 0)
	defer call(clipboard, aClose)
	call(clipboard, aClear, 0)
}

// Set the clipboard data
func Set(text string) {
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	startAlloc, lenAlloc := paramAlloc(text)
	// init
	call(clipboard, aOpen, 0)
	defer call(clipboard, aClose)
	call(clipboard, aClear, 0)
	// allocate
	mem := call(memory, aAlloc, fixedMemory, lenAlloc)
	call(memory, aMove, call(memory, aLock, mem), startAlloc, lenAlloc)
	defer call(memory, aUnlock, mem)
	// set string
	checkSet := call(clipboard, aSet, unicodeTextFormat, mem)
	if checkSet == 0 {
		call(memory, aFree, mem)
	}
}

// CheckTopWindow check top(focus) window
func (w *WinWin) CheckTopWindow() bool {
	if w == nil {
		return false
	}
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	t := time.Now().Add(time.Second)
	for time.Since(t) <= 0 {

		activeThread := call(window, aWindowId, call(window, aWindowFname))
		targetThread := call(window, aWindowId, uintptr(w.Handle))

		if activeThread == targetThread {
			return true
		}
		<-time.After(10 * time.Millisecond)

	}
	return false
}

// SetTopWindow top(focus) window
func (w *WinWin) SetTopWindow() bool {
	if w == nil {
		return false
	}
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	t := time.Now().Add(time.Second)
	for time.Since(t) <= 0 {

		// --- old version
		// activeThread := call(window, aWindowId, call(window, aWindowFname))
		// currentThread := call(procss, aCurrentWindowId)
		// targetThread := call(window, aWindowId, uintptr(w.Handle))

		// if currentThread != activeThread {
		// 	call(window, aAttachThreadInput, currentThread, activeThread, 1)
		// }
		// if targetThread != currentThread {
		// 	call(window, aAttachThreadInput, targetThread, currentThread, 1)
		// }

		ok := call(window, aWindowTop, uintptr(w.Handle)) > 0
		time.Sleep(100 * time.Millisecond)
		call(window, aBringWindowToTop, uintptr(w.Handle))
		time.Sleep(100 * time.Millisecond)
		call(window, aWindowShow, uintptr(w.Handle), normalShow)
		time.Sleep(100 * time.Millisecond)
		call(window, aEnableWindow, uintptr(w.Handle), 1)
		time.Sleep(100 * time.Millisecond)
		call(window, aSetActiveWindow, uintptr(w.Handle))
		time.Sleep(100 * time.Millisecond)
		call(window, aFlashWindow, uintptr(w.Handle), 1)
		time.Sleep(100 * time.Millisecond)

		if ok {
			return true
		}

		// --- old version
		// if currentThread != activeThread {
		// 	call(window, aAttachThreadInput, currentThread, activeThread, 0)
		// }
		// if targetThread != currentThread {
		// 	call(window, aAttachThreadInput, targetThread, currentThread, 0)
		// }

		// activeThread = call(window, aWindowId, call(window, aWindowFname))
		// targetThread = call(window, aWindowId, uintptr(w.Handle))

		// if activeThread == targetThread {
		//   return true
		// }
		<-time.After(400 * time.Millisecond)

	}
	return false
}

// SetTopClickCenterWindow top(focus) window with click
func (w *WinWin) SetTopClickCenterWindow() bool {
	if ok := w.SetTopWindow(); !ok {
		return ok
	}
	x := call(window, aGetSystemMetrics, 0) / 2
	y := call(window, aGetSystemMetrics, 1) / 2
	x = (x * 65536) / call(window, aGetSystemMetrics, 0)
	y = (y * 65536) / call(window, aGetSystemMetrics, 1)
	time.Sleep(100 * time.Millisecond)
	call(window, aMouseEvent, 0x8000|0x0001, x, y, 0)
	call(window, aMouseEvent, 0x0002)
	call(window, aMouseEvent, 0x0004)
	time.Sleep(100 * time.Millisecond)
	return true
}

// FindWindow window by the pid or the part title
func (w *WinWin) FindWindow(titleOrPid string) bool {
	runtime.LockOSThread() // Using CALL in one thread
	defer runtime.UnlockOSThread()
	var hwnd syscall.Handle
	var findText = ""
	var windowPid uintptr = 0
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		call(window, aWindowName, uintptr(h), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
		findText = syscall.UTF16ToString(b)
		call(window, aWindowId, uintptr(h), uintptr(unsafe.Pointer(&windowPid)))

		//
		if strings.Contains(findText, titleOrPid) && call(window, aWindowGet, uintptr(h), uintptr(4)) == 0 && call(window, aWindowVisible, uintptr(h)) == 1 {
			hwnd = h
			return 0 // stop enum
		}
		//
		if strings.EqualFold(fmt.Sprintf("%d", windowPid), titleOrPid) && call(window, aWindowGet, uintptr(h), uintptr(4)) == 0 && call(window, aWindowVisible, uintptr(h)) == 1 {
			hwnd = h
			return 0 // stop enum
		}
		return 1 // continue enum
	})
	call(window, aWindowEnum, cb, 0, 0)
	if hwnd != 0 {
		w.Title = findText
		w.ProcessID = windowPid
		w.Handle = uintptr(hwnd)
		return true
	}
	return false
}

// FindWindow check windows by pid/title
func FindWindow[T int | string](titleOrPid T) *WinWin {
	var w = new(WinWin)
	w.FindWindow(fmt.Sprintf("%v", titleOrPid))
	return w
}

// SetTopWindow top(focus) window
func SetTopWindow[T int | string](titleOrPid T) bool {
	return FindWindow(titleOrPid).SetTopWindow()
}

// SetTopWindowClick top(focus) window
func SetTopWindowClick[T int | string](titleOrPid T) bool {
	return FindWindow(titleOrPid).SetTopClickCenterWindow()
}

// CheckTopWindow top(focus) window
func CheckTopWindow[T int | string](titleOrPid T) bool {
	return FindWindow(titleOrPid).CheckTopWindow()
}

// Find window
func Find(pid int, title string) <-chan int {
	c := make(chan int, 1)
	if pid != 0 {
		if ww := FindWindow(pid); *ww != (WinWin{}) {
			c <- int(ww.ProcessID)
		} else {
			return FindTitle(title)
		}
	} else {
		return FindTitle(title)
	}
	return c
}

// FindPid window
func FindPid(pid int) <-chan int {
	c := make(chan int, 1)
	if pid != 0 {
		if ww := FindWindow(pid); *ww != (WinWin{}) {
			c <- int(ww.ProcessID)
		}
	}
	return c
}

// FindTitle window
func FindTitle(title string) <-chan int {
	c := make(chan int, 1)
	if w := FindWindow(title); *w != (WinWin{}) {
		c <- int(w.ProcessID)
	}
	return c
}

// LostPid window
func LostPid(pid int) <-chan struct{} {
	c := make(chan struct{}, 1)
	if pid != 0 {
		if ww := FindWindow(pid); *ww == (WinWin{}) {
			c <- struct{}{}
		}
	}
	return c
}

// LostTitle window
func LostTitle(title string) <-chan struct{} {
	c := make(chan struct{}, 1)
	if w := FindWindow(title); *w == (WinWin{}) {
		c <- struct{}{}
	}
	return c
}

// FindProcess window. TODO
func FindProcess(pid int) (*Process, error) {
	ps, err := processes()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		if p.pid == pid {
			return &p, nil
		}
	}
	return nil, nil
}

// newWindowsProcess window. TODO
func newWindowsProcess(e *PROCESSENTRY32, h uintptr) Process {
	end := 0
	for e.ExeFile[end] != 0 {
		end++
	}
	return Process{
		pid:  int(e.ProcessID),
		ppid: int(e.ParentProcessID),
		exe:  syscall.UTF16ToString(e.ExeFile[:end]),
		h:    h,
	}
}

// processes window. TODO
func processes() ([]Process, error) {
	handle := call(procss, aCreateToolhelp32Snapshot, 0x00000002, 0)
	defer call(procss, aProcClose)
	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret := call(procss, aProcess32First, handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("error retrieving process info")
	}
	results := make([]Process, 0, 50)
	for {
		n := newWindowsProcess(&entry, handle)
		results = append(results, n)
		ret := call(procss, aProcess32Next, handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
	return results, nil
}
