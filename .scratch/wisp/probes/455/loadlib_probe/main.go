//go:build windows

// loadlib_probe -- answers the one Windows-loader question that decides
// whether staging next to the exe is enough (ticket 455, item 2, bullet 3):
//
//	when a DLL is loaded by an explicit LoadLibrary/Explicit LoadLibraryEx call
//	from an exe that lives in ANOTHER directory, is the loading DLL's own
//	directory searched for ITS dependencies (onnxruntime.dll)?
//
// Usage: loadlib_probe.exe <path-to-sherpa-onnx-c-api.dll> <flags-hex>
//   flags 0            -> default search order
//   flags 0x8          -> LOAD_WITH_ALTERED_SEARCH_PATH
//   flags 0x00000100   -> LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR (needs SetDefaultDllDirectories)
//
// Prints the loaded onnxruntime.dll path + size, so the answer is a path, not
// a version string.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type loadedModule struct {
	base uintptr
	path string
}

func loadedModules() ([]loadedModule, error) {
	proc := windows.CurrentProcess()
	buf := make([]windows.Handle, 2048)
	var needed uint32
	cb := uint32(len(buf)) * uint32(unsafe.Sizeof(buf[0]))
	if err := windows.EnumProcessModules(proc, &buf[0], cb, &needed); err != nil {
		return nil, err
	}
	n := int(needed / uint32(unsafe.Sizeof(buf[0])))
	if n > len(buf) {
		n = len(buf)
	}
	name := make([]uint16, windows.MAX_PATH)
	out := make([]loadedModule, 0, n)
	for _, m := range buf[:n] {
		if err := windows.GetModuleFileNameEx(proc, m, &name[0], uint32(len(name))); err != nil {
			continue
		}
		out = append(out, loadedModule{base: uintptr(m), path: windows.UTF16ToString(name)})
	}
	return out, nil
}

func sizeOf(path string) string {
	st, err := os.Stat(path)
	if err != nil {
		return "stat: " + err.Error()
	}
	return fmt.Sprintf("%d B", st.Size())
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: loadlib_probe.exe <path-to-sherpa-onnx-c-api.dll> <flags-hex>")
		os.Exit(2)
	}
	dll := os.Args[1]
	flags64, err := strconv.ParseUint(strings.TrimPrefix(os.Args[2], "0x"), 16, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad flags:", os.Args[2])
		os.Exit(2)
	}
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	fmt.Printf("== loadlib probe: flags=0x%x ==\n", flags64)
	fmt.Printf("exe        : %s\n", exe)
	fmt.Printf("cwd        : %s\n", cwd)
	fmt.Printf("loading    : %s\n", dll)

	// LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR only takes effect after the process
	// opts into the extended search set; do so when that flag is requested.
	if flags64&0x00000100 != 0 {
		if err := windows.SetDefaultDllDirectories(windows.LOAD_LIBRARY_SEARCH_DEFAULT_DIRS); err != nil {
			fmt.Println("SetDefaultDllDirectories:", err)
		}
	}
	p, err := windows.UTF16PtrFromString(dll)
	if err != nil {
		fmt.Println("UTF16PtrFromString:", err)
		os.Exit(1)
	}
	procLoadLibraryExW := windows.NewLazySystemDLL("kernel32.dll").NewProc("LoadLibraryExW")
	h, _, callErr := procLoadLibraryExW.Call(uintptr(unsafe.Pointer(p)), 0, uintptr(flags64))
	if h == 0 {
		fmt.Println("LoadLibraryEx error:", callErr)
		os.Exit(1)
	}
	mods, err := loadedModules()
	if err != nil {
		fmt.Println("EnumProcessModules:", err)
		os.Exit(1)
	}
	for _, m := range mods {
		base := strings.ToLower(filepath.Base(m.path))
		if strings.Contains(base, "onnxruntime") || strings.Contains(base, "sherpa-onnx") {
			fmt.Printf("module     : %-26s path=%s\n             size=%s\n",
				filepath.Base(m.path), m.path, sizeOf(m.path))
		}
	}
}
