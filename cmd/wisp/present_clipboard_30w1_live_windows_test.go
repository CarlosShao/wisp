//go:build windows && winlive

package main

// 票 30 装配缝腿 30-w-r1 的 live 层：真 Win32 剪贴板写半边的往返证明（写进、读回、
// 逐字同）。这一层会真的改本机剪贴板一次，所以跟票 43/228 的 live 层同一规矩：
// 只在 -tags winlive 下单独跑，不进默认套件：
//
//	go test -tags winlive ./cmd/wisp -run TestLive30Wr1 -v

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	liveClipUser32 = windows.NewLazySystemDLL("user32.dll")
	liveClipOpen   = liveClipUser32.NewProc("OpenClipboard")
	liveClipClose  = liveClipUser32.NewProc("CloseClipboard")
	liveClipGet    = liveClipUser32.NewProc("GetClipboardData")

	liveClipK32 = windows.NewLazySystemDLL("kernel32.dll")
	liveMove    = liveClipK32.NewProc("RtlMoveMemory")
	liveLock    = liveClipK32.NewProc("GlobalLock")
	liveUnlock  = liveClipK32.NewProc("GlobalUnlock")
	liveSize    = liveClipK32.NewProc("GlobalSize")
)

// readClipboardTextWindows reads the clipboard's CF_UNICODETEXT back as text.
// The block is moveable, so GlobalLock names its address and RtlMoveMemory
// copies out through it - both stay uintptr on this side, which is the same
// vet-clean shape the write half uses (no uintptr ever turns back into a
// Pointer in this package). It is the live probe's own reader, not a
// production path.
func readClipboardTextWindows() (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// OpenClipboard fails while another window holds the clipboard open, and
	// the moment right after a write is exactly when shell-side viewers race
	// it. The budget is a fixed attempt count with a fixed delay - no clock
	// read decides anything, which is the same discipline the notify balloon
	// hold already uses.
	opened := false
	for attempt := 0; attempt < 40; attempt++ {
		if r, _, _ := liveClipOpen.Call(0); r != 0 {
			opened = true
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if !opened {
		return "", errors.New("OpenClipboard stayed refused for the whole attempt budget (another window holds the clipboard)")
	}
	defer liveClipClose.Call()
	h, _, callErr := liveClipGet.Call(clipCFUnicodeText)
	if h == 0 {
		return "", fmt.Errorf("GetClipboardData: %v", callErr)
	}
	size, _, callErr := liveSize.Call(h)
	if size == 0 {
		return "", fmt.Errorf("GlobalSize: %v", callErr)
	}
	locked, _, callErr := liveLock.Call(h)
	if locked == 0 {
		return "", fmt.Errorf("GlobalLock: %v", callErr)
	}
	defer liveUnlock.Call(h)
	buf := make([]uint16, size/unsafe.Sizeof(uint16(0)))
	liveMove.Call(uintptr(unsafe.Pointer(&buf[0])), locked, size)
	return windows.UTF16ToString(buf), nil
}

func TestLive30Wr1RealClipboardWriteRoundTrip(t *testing.T) {
	const marker = "wisp 30wr1 clipboard live probe（真写半边的往返证明）"
	if err := writeClipboardWindows(marker); err != nil {
		t.Fatalf("writeClipboardWindows: %v", err)
	}
	got, err := readClipboardTextWindows()
	if err != nil {
		t.Fatalf("read the clipboard back: %v", err)
	}
	if got != marker {
		t.Errorf("clipboard round trip = %q, want %q", got, marker)
	}
}
