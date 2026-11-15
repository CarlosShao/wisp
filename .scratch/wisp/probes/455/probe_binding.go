//go:build windows

// probe_binding.go -- ticket 455 recon probe.
//
// Answers, from INSIDE the process, which onnxruntime.dll the loader actually
// bound: the absolute path (GetModuleFileNameEx over the real module list) and
// the byte size on disk, plus the version the linked sherpa C API reports.
// The evidence is the PATH and the SIZE, not the version string alone -- the
// red sentence only prints the shadowing copy's version, never its path.
//
// Build shape (mirrors scripts/wisp-cli-staged-tests.sh):
//   go build -o <stage>/probe.exe .
//   copy third_party/sherpa-onnx/*.dll next to probe.exe
//   run it with cwd=cmd/wisp
//
// Run it once with the DLLs staged next to the exe (carrier shape) and once
// from a directory with no DLLs but third_party/sherpa-onnx on PATH (the raw
// `go test` shape): the two readings are the whole point.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"golang.org/x/sys/windows"
)

type loadedModule struct {
	base uintptr
	path string
}

// loadedModules enumerates the process's loaded images (same instrument the
// repo already trusts: internal/speech/cgo_boundary_windows_test.go:34-59).
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
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	fmt.Println("== ticket 455 probe: which onnxruntime.dll is bound ==")
	fmt.Printf("exe        : %s\n", exe)
	fmt.Printf("exe dir    : %s\n", filepath.Dir(exe))
	fmt.Printf("cwd        : %s\n", cwd)
	fmt.Printf("GetOnnxruntimeVersion(): %s\n", sherpa.GetOnnxruntimeVersion())

	mods, err := loadedModules()
	if err != nil {
		fmt.Println("EnumProcessModules:", err)
		os.Exit(1)
	}
	want := []string{"onnxruntime", "sherpa-onnx"}
	found := 0
	for _, m := range mods {
		base := strings.ToLower(filepath.Base(m.path))
		for _, w := range want {
			if strings.Contains(base, w) {
				fmt.Printf("module     : %-26s base=%#x\n             path=%s\n             size=%s\n",
					filepath.Base(m.path), m.base, m.path, sizeOf(m.path))
				found++
				break
			}
		}
	}
	if found == 0 {
		fmt.Println("module     : NONE of onnxruntime/sherpa-onnx is mapped")
	}
	for _, e := range filepath.SplitList(os.Getenv("PATH")) {
		if strings.Contains(strings.ToLower(e), "sherpa") {
			fmt.Printf("PATH sherpa: %s\n", e)
		}
	}
}
