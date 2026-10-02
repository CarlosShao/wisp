//go:build windows

package main

// 33-p1 leg 2 (2026-10-02), question 3: cold and hot bring-up latency measured
// on this machine TODAY, through the same measurement the shipping host uses.
//
// The shipping measurement lives in cmd/wisp/panel_host_windows.go
// (PanelManager.lastColdMs = window create + first browser round trip, stamped
// by firstRoundTripLocked; lastHotMs = Show-after-Hide on the reused window).
// PanelManager is package main, not importable, so this probe re-measures the
// SAME two quantities with the SAME definitions on its own thread:
//
//	cold = NewWithOptions -> SetHtml(probe page) -> the probe binding actually
//	       called back on this thread (page -> Go hop proves the channel), ms.
//	hot  = a further Show/after-Hide equivalent is NOT what this file can do
//	       (Hide/Show are PanelManager methods). The hot equivalent here =
//	       SetHtml round trip AGAIN on the live window (channel already up),
//	       i.e. the reuse part of the hot path minus the Win32 show calls.
//	       Labelled [hot-approx] for exactly that reason.
//
// Readings print verbatim; no threshold is asserted here (the D32 budgets
// 1500/200 belong to the SLO files; this probe only reports).
//
// Clean-exit discipline: Destroy on the owning thread, pump to quiet, exit.

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

const (
	wmQuit    = 0x0012
	coinitSTA = 0x2
	rounds    = 3
)

var (
	modUser32   = windows.NewLazySystemDLL("user32.dll")
	modOle32    = windows.NewLazySystemDLL("ole32.dll")
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")
	pPeekMsgW   = modUser32.NewProc("PeekMessageW")
	pTransMsg   = modUser32.NewProc("TranslateMessage")
	pDispMsgW   = modUser32.NewProc("DispatchMessageW")
	pPostThrW   = modUser32.NewProc("PostThreadMessageW")
	pCoInitEx   = modOle32.NewProc("CoInitializeEx")
	pCoUninit   = modOle32.NewProc("CoUninitialize")
	pGetCurTID  = modKernel32.NewProc("GetCurrentThreadId")
)

type pnlMsg struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	_       uint32
}

func pumpOnce() {
	var m pnlMsg
	for {
		r, _, _ := pPeekMsgW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1)
		if r == 0 {
			return
		}
		pTransMsg.Call(uintptr(unsafe.Pointer(&m)))
		pDispMsgW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

var (
	uiTID       = make(chan uint32, 1)
	rtDone      = make(chan struct{}, 1)
	rtCalled    atomic.Int64
	roundStamp  atomic.Int64
	hostDone    = make(chan struct{})
	quitPending atomic.Bool
)

const probePage = `<!doctype html><html><head><meta charset="utf-8"></head><body>
<script>if(window.wispProbeRT)window.wispProbeRT();</script></body></html>`

func main() {
	deadlineS := flag.Int("deadline", 25, "watchdog seconds")
	flag.Parse()

	dir := filepath.Join(os.TempDir(), "wisp33p1q3", "latency")
	_ = os.MkdirAll(dir, 0o755)
	sinkFile := filepath.Join(dir, "panic-sink.txt")

	reg := observe.NewRegistry()
	var sinkWrites atomic.Int64
	reg.SetPanicSink(func(ev observe.PanicEvent) {
		sinkWrites.Add(1)
		body := fmt.Sprintf("Goroutine=%s Owner=%s Recovered=%s At=%s\nSTACK:\n%s\n",
			ev.Goroutine, ev.Owner, ev.Recovered, ev.At, ev.Stack)
		_ = os.WriteFile(sinkFile, []byte(body), 0o644)
		fmt.Printf("SINK-EVENT goroutine=%s recovered=%q file=%s\n", ev.Goroutine, ev.Recovered, sinkFile)
	})

	reg.Spawn("p1q3-watchdog", "33-p1-q3", nil, func(_ context.Context) {
		tid := <-uiTID
		select {
		case <-hostDone:
		case <-time.After(time.Duration(*deadlineS) * time.Second):
			quitPending.Store(true)
			r, _, e := pPostThrW.Call(uintptr(tid), wmQuit, 0, 0)
			fmt.Printf("WATCHDOG forced WM_QUIT to tid=%d r=%d err=%v (readings tainted)\n", tid, r, e)
		}
	})

	reg.Spawn("p1q3-host", "33-p1-q3", nil, func(_ context.Context) { runHost(dir) })

	select {
	case <-hostDone:
	case <-time.After(time.Duration(*deadlineS+15) * time.Second):
		fmt.Printf("RESULT=HOST-NEVER-FINISHED\n")
	}
	fmt.Printf("FINAL panics=%d sink_writes=%d tainted=%v sink_file=%s\n",
		reg.PanicCount(), sinkWrites.Load(), quitPending.Load(), sinkFile)
}

func runHost(dir string) {
	defer close(hostDone)
	defer pCoUninit.Call()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tid, _, _ := pGetCurTID.Call()
	uiTID <- uint32(tid)
	if r, _, _ := pCoInitEx.Call(0, coinitSTA); r != 0 && r != 1 {
		fmt.Printf("CoInitializeEx FAILED r=%d\n", r)
		return
	}
	fmt.Printf("HOST tid=%d com=sta\n", tid)

	// ---- COLD (shipping definition: create + first page->Go round trip) ----
	t0 := time.Now()
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:    false,
		DataPath: filepath.Join(dir, "userdata"),
		WindowOptions: webview2.WindowOptions{Title: "33-p1-q3 latency", Width: 420, Height: 260},
	})
	if w == nil {
		fmt.Printf("COLD create RETURNED-NIL after %s\n", time.Since(t0).Truncate(time.Millisecond))
		return
	}
	createMs := float64(time.Since(t0).Microseconds()) / 1000.0
	fmt.Printf("COLD create alone: %.3f ms\n", createMs)

	_ = w.Bind("wispProbeRT", func() string {
		rtCalled.Add(1)
		select {
		case rtDone <- struct{}{}:
		default:
		}
		return "ok"
	})

	w.SetHtml(probePage)
	deadline := t0.Add(10 * time.Second)
	for {
		pumpOnce()
		select {
		case <-rtDone:
			coldMs := float64(time.Since(t0).Microseconds()) / 1000.0
			fmt.Printf("COLD (create+roundtrip) = %.3f ms (roundtrip calls=%d)\n", coldMs, rtCalled.Load())
			roundStamp.Store(time.Now().UnixNano())

			// ---- HOT-approx: channel already up; re-drive a round trip ----
			// 2 more rounds, each timed the same way.
			for i := 1; i <= rounds-1; i++ {
				rtCalled.Store(0)
				th := time.Now()
				w.SetHtml(probePage + fmt.Sprintf("<!-- %d -->", i))
				dl := time.Now().Add(5 * time.Second)
				got := false
				for {
					pumpOnce()
					if rtCalled.Load() > 0 {
						got = true
						break
					}
					if time.Now().After(dl) {
						break
					}
					time.Sleep(2 * time.Millisecond)
				}
				if got {
					fmt.Printf("HOT-approx round %d = %.3f ms\n", i, float64(time.Since(th).Microseconds())/1000.0)
				} else {
					fmt.Printf("HOT-approx round %d = NO-ROUNDTRIP in 5s\n", i)
				}
			}
			break
		default:
		}
		if rtCalled.Load() > 0 {
			break
		}
		if time.Now().After(deadline) {
			fmt.Printf("COLD NO-ROUNDTRIP in 10s (channel never came up)\n")
			break
		}
		time.Sleep(2 * time.Millisecond)
	}

	w.Destroy()
	pumpReap(60)
	fmt.Printf("TEARDOWN destroyed+pumped\n")
}

func pumpReap(n int) {
	for i := 0; i < n; i++ {
		pumpOnce()
		time.Sleep(25 * time.Millisecond)
	}
}
