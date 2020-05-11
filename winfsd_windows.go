package winfsd

import (
	"debug/pe"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// lazyProc is a copy of syscall.LazyProc from Go's src/syscall/dll_windows.go.
type lazyProc struct {
	mu   sync.Mutex
	Name string
	l    *syscall.LazyDLL
	proc *proc
}

// proc is a copy of syscall.Proc from Go's src/syscall/dll_windows.go.
type proc struct {
	Dll  *syscall.DLL
	Name string
	addr uintptr
}

// createFileWAddr is set by init to the address of kernel32!CreateFileW.
var createFileWAddr uintptr

// createFileWTrampoline is implemented in winfsd_windows_amd64.s.
func createFileWTrampoline()

func init() {
	var module windows.Handle
	null := unsafe.Pointer(uintptr(0))
	if err := windows.GetModuleHandleEx(0, (*uint16)(null), &module); err != nil {
		panic("failed to get module handle: " + err.Error())
	}
	windows.CloseHandle(module)

	exe, err := os.Executable()
	if err != nil {
		panic("failed to get executable: " + err.Error())
	}
	peFile, err := pe.Open(exe)
	if err != nil {
		panic("failed to open executable: " + err.Error())
	}
	peFile.Close()
	var symbol *pe.Symbol
	for _, s := range peFile.Symbols {
		if s.Name == "syscall.procCreateFileW" {
			symbol = s
			break
		}
	}
	if symbol == nil {
		panic("symbol not found")
	}

	// The value of a module handle is the module's load address.
	baseAddr := uintptr(module)
	dataSectionAddr := baseAddr + uintptr(peFile.Section(".data").VirtualAddress)
	procCreateFileWAddr := dataSectionAddr + uintptr(symbol.Value)
	procCreateFileWPtr := unsafe.Pointer(procCreateFileWAddr)

	// Save the address of kernel32!CreateFileW.  Calling Addr also sets the
	// proc field, needed below.
	createFileWAddr = (*(**syscall.LazyProc)(procCreateFileWPtr)).Addr()

	// Install createFileWTrampoline.
	createFileWTrampolineAddr := reflect.ValueOf(createFileWTrampoline).Pointer()
	(*(**lazyProc)(procCreateFileWPtr)).proc.addr = createFileWTrampolineAddr
}
