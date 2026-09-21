//go:build windows

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// dllFileVersion reads the VS_FIXEDFILEINFO FileVersion of a PE file through
// the Win32 version API.
var (
	modVersion              = windows.NewLazySystemDLL("version.dll")
	procGetFileVersionInfoW = modVersion.NewProc("GetFileVersionInfoW")
	procGetFileVersionSizeW = modVersion.NewProc("GetFileVersionInfoSizeW")
	procVerQueryValueW      = modVersion.NewProc("VerQueryValueW")
)

type vsFixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

func dllFileVersion(path string) (string, bool) {
	p16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", false
	}
	size, _, _ := procGetFileVersionSizeW.Call(uintptr(unsafe.Pointer(p16)), 0)
	if size == 0 {
		return "", false
	}
	buf := make([]byte, size)
	r1, _, _ := procGetFileVersionInfoW.Call(uintptr(unsafe.Pointer(p16)), 0, size, uintptr(unsafe.Pointer(&buf[0])))
	if r1 == 0 {
		return "", false
	}
	root, err := windows.UTF16PtrFromString("\\")
	if err != nil {
		return "", false
	}
	var ffi *vsFixedFileInfo
	var ffiLen uint32
	r1, _, _ = procVerQueryValueW.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(root)),
		uintptr(unsafe.Pointer(&ffi)), uintptr(unsafe.Pointer(&ffiLen)))
	if r1 == 0 || ffi == nil || ffi.Signature != 0xFEEF04BD {
		return "", false
	}
	return fmt.Sprintf("%d.%d.%d.%d",
		ffi.FileVersionMS>>16, ffi.FileVersionMS&0xffff,
		ffi.FileVersionLS>>16, ffi.FileVersionLS&0xffff), true
}
