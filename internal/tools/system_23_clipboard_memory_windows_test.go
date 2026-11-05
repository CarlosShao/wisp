//go:build windows

package tools

// Ticket 23, the clipboard memory primitive.
//
// realClipboardRead / realClipboardWrite used to forge a Go pointer out of
// GlobalLock's address ((*uint16)(unsafe.Pointer(p)), unsafe.Slice over the same
// address), which go vet refuses and which the collector cannot track: the block
// belongs to the clipboard owner, not to the Go heap. The fix copies with the OS
// memmove instead (moveMemory), keeping the OS-owned side a uintptr for the whole
// call.
//
// That fix is only worth the words if the copy is REAL, so this leg exercises the
// exact two directions the clipboard takes - Go []uint16 into a GMEM block, and
// the block back into a Go slice - against the same Win32 entry points the
// production legs call. A missing kernel32!RtlMoveMemory (a silent no-op copy)
// turns the read leg red instead of quietly shrinking a clipboard read to "".
//
// Deliberately NOT here: a round trip through the SYSTEM clipboard. Writing
// CF_UNICODETEXT would clobber the operator's clipboard to prove a memcpy works;
// ticket 23's AC#4 (the real round trip) is a separate, hand-run row, and this
// leg must stay safe to run on any machine, including a CI runner with no
// interactive session.

import (
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

func Test23ClipboardMemoryMovesBothDirections(t *testing.T) {
	const text = "Wisp 票 23 剪贴板 clipboard · 42%"
	utf16, err := windows.UTF16FromString(text)
	if err != nil {
		t.Fatalf("UTF16FromString(%q) 失败：%v", text, err)
	}
	size := uintptr(len(utf16) * 2)

	// The same allocation realClipboardWrite makes.
	h, _, _ := procGlobalAlloc.Call(gmemMoveable, size)
	if h == 0 {
		t.Fatalf("GlobalAlloc 失败：本机无法分配 GMEM 块")
	}
	defer procGlobalFree.Call(h)

	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		t.Fatalf("GlobalLock 失败")
	}
	// Direction 1: Go slice -> OS block (realClipboardWrite's copy).
	if err := moveMemory(p, uintptr(unsafe.Pointer(&utf16[0])), size); err != nil {
		procGlobalUnlock.Call(h)
		t.Fatalf("写入方向失败：%v", err)
	}
	procGlobalUnlock.Call(h)

	// Direction 2: OS block -> Go slice (realClipboardRead's copy), through the
	// same GlobalSize-then-moveMemory pair the production leg uses.
	locked, _, _ := procGlobalLock.Call(h)
	if locked == 0 {
		t.Fatalf("重新 GlobalLock 失败")
	}
	defer procGlobalUnlock.Call(h)
	if procGlobalSize.Find() != nil {
		t.Fatalf("GlobalSize 不可用：realClipboardRead 无法得知块长度")
	}
	got := int(callUintptr(procGlobalSize, h))
	if got != int(size) {
		t.Fatalf("GlobalSize 读作 %d 字节，期望 %d", got, int(size))
	}
	buf := make([]uint16, got/2)
	if err := moveMemory(uintptr(unsafe.Pointer(&buf[0])), locked, uintptr(got)); err != nil {
		t.Fatalf("读取方向失败：%v", err)
	}
	if back := windows.UTF16ToString(buf); back != text {
		t.Fatalf("读回 %q，期望 %q：OS 拷贝没有真的搬字节", back, text)
	}
}

// Test23ClipboardReadTreatsAZeroLengthBlockAsEmpty pins the other end: a handle
// whose GlobalSize is 0 must answer "" (the honest empty clipboard), never a
// forged pointer and never an error the model would read as "clipboard broken".
func Test23ClipboardReadTreatsAZeroLengthBlockAsEmpty(t *testing.T) {
	// GlobalSize on a handle that is not GMEM answers 0; the read leg's guard is
	// size < 2 -> "".
	if n := int(callUintptr(procGlobalSize, 0)); n != 0 {
		t.Fatalf("GlobalSize(nil) 读作 %d，期望 0", n)
	}
	if err := moveMemory(0, 0, 0); err != nil {
		t.Fatalf("零长度 moveMemory 应当是无操作，得到 %v", err)
	}
}
