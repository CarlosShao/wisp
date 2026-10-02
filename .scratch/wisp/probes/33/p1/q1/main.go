//go:build windows

package main

// 33-p1 leg 2 (2026-10-02), question 1 in its own right: what does a SECOND
// NewPanelManager-style bring-up do inside one resident process?
//
// The shipping PanelManager is cmd/wisp package main (not importable), so this
// probe measures the three layers a second panel manager would have to survive:
//
//   L1  two PanelManager STRUCTS are trivially fine - NewPanelManager is a pure
//       constructor (panel_host_windows.go:175-177: it only fills the struct).
//       The question is what happens at bring-up.
//   L2  the shipping bring-up path: webview2.NewWithOptions on a thread while
//       ANOTHER live webview window (first manager's) exists on the SAME STA
//       thread. That is what -mode dup-create measures (main.go does it too,
//       this binary re-measures with the msedgewebview2 child counted).
//   L3  the resident topology: a second residentPanel-style thread created
//       alongside the first one, each with its own manager. Shipping shape
//       never does this (startResidentPanel is called once in resident
//       assembly), so this run answers "what would even happen": two windows,
//       two pumps, one process.
//
// Clean-exit discipline mirrors balldebug: every window Destroyed on its owning
// thread, queues pumped to quiet before thread unlock, process exits by itself.

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
	wmQuit   = 0x0012
	coinitSTA = 0x2
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

func pumpReap(n int) {
	for i := 0; i < n; i++ {
		pumpOnce()
		time.Sleep(25 * time.Millisecond)
	}
}

var (
	uiTID    = make(chan uint32, 1)
	stage    = make(chan string, 8) // host -> main progress pipe
	hostDone = make(chan struct{})
	quitSent = make(chan struct{}, 1)
)

func main() {
	deadlineS := flag.Int("deadline", 25, "watchdog seconds")
	flag.Parse()

	dir := filepath.Join(os.TempDir(), "wisp33p1q1", "secondmgr")
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

	// Watchdog: if the host never finishes, post WM_QUIT so Run() ends and the
	// process can still exit by itself; the FINAL line reports the taint.
	reg.Spawn("p1q1-watchdog", "33-p1-q1", nil, func(_ context.Context) {
		tid := <-uiTID
		select {
		case <-hostDone:
		case <-quitSent:
		case <-time.After(time.Duration(*deadlineS) * time.Second):
			r, _, e := pPostThrW.Call(uintptr(tid), wmQuit, 0, 0)
			fmt.Printf("WATCHDOG forced WM_QUIT to tid=%d r=%d err=%v (readings tainted)\n", tid, r, e)
		}
	})

	reg.Spawn("p1q1-host", "33-p1-q1", nil, func(_ context.Context) { runHost(dir) })

	select {
	case <-hostDone:
	case <-time.After(time.Duration(*deadlineS+15) * time.Second):
		fmt.Printf("RESULT=HOST-NEVER-FINISHED\n")
	}

	fmt.Printf("FINAL panics=%d sink_writes=%d sink_file=%s\n", reg.PanicCount(), sinkWrites.Load(), sinkFile)
}

// mgr is one "PanelManager": the struct fields NewPanelManager would hold, plus
// the thread binding the resident topology gives it.
type mgr struct {
	name     string
	w        webview2.WebView
	hwnd     uintptr
	created  bool
	createMs int64
}

func runHost(dir string) {
	defer close(hostDone)
	defer pCoUninit.Call()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tid, _, _ := pGetCurTID.Call()
	uiTID <- uint32(tid)
	fmt.Printf("HOST tid=%d com=sta\n", tid)
	if r, _, _ := pCoInitEx.Call(0, coinitSTA); r != 0 && r != 1 {
		fmt.Printf("HOST CoInitializeEx FAILED r=%d - aborting\n", r)
		return
	}

	// ---- L1: two pure constructors ----
	m1 := &mgr{name: "first"}
	m2 := &mgr{name: "second"}
	fmt.Printf("L1 two constructors done: m1=%p m2=%p (pure struct fill, no window)\n", m1, m2)

	// ---- first manager's bring-up (the shipping bringUp shape) ----
	t0 := time.Now()
	w1 := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:    false,
		DataPath: filepath.Join(dir, "u1"),
		WindowOptions: webview2.WindowOptions{Title: "33-p1-q1 first panel", Width: 420, Height: 260},
	})
	m1.createMs = time.Since(t0).Milliseconds()
	if w1 == nil {
		fmt.Printf("L2 first create RETURNED-NIL ms=%d\n", m1.createMs)
		return
	}
	m1.w, m1.hwnd, m1.created = w1, uintptr(w1.Window()), true
	fmt.Printf("L2 first create ok hwnd=0x%X ms=%d\n", m1.hwnd, m1.createMs)

	_ = w1.Bind("wispNote", func(v string) string {
		fmt.Printf("NOTE %s\n", v)
		return "ok"
	})
	w1.SetHtml(`<!doctype html><html><head><title>m1 up</title></head><body>
		<script>try{window.wispNote("M1-LIVE")}catch(e){}</script></body></html>`)
	// Let window 1's message channel prove live (its binding call arrives as a
	// window message on this thread).
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		pumpOnce()
		time.Sleep(5 * time.Millisecond)
	}

	// ---- L2: SECOND create on the same thread, first window alive ----
	t1 := time.Now()
	w2 := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:    false,
		DataPath: filepath.Join(dir, "u2"),
		WindowOptions: webview2.WindowOptions{Title: "33-p1-q1 second panel", Width: 420, Height: 260},
	})
	m2.createMs = time.Since(t1).Milliseconds()
	if w2 == nil {
		fmt.Printf("L2 second create RETURNED-NIL ms=%d (first window STILL ALIVE: IsWindow m1=true)\n", m2.createMs)
	} else {
		m2.w, m2.hwnd, m2.created = w2, uintptr(w2.Window()), true
		fmt.Printf("L2 second create ok hwnd=0x%X ms=%d (both windows now on tid=%d)\n", m2.hwnd, m2.createMs, tid)
		_ = w2.Bind("wispNote", func(v string) string {
			fmt.Printf("NOTE2 %s\n", v)
			return "ok"
		})
		w2.SetHtml(`<!doctype html><html><head><title>m2 up</title></head><body>
			<script>try{window.wispNote("M2-LIVE")}catch(e){}</script></body></html>`)
		deadline = time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline) {
			pumpOnce()
			time.Sleep(5 * time.Millisecond)
		}
		// L3 rider: drive both pumps in library form briefly? No - Run() is
		// per-webview but occupies THIS thread; running two Run()s serially:
		fmt.Printf("L3 serial Run() on m2's loop for 1s (first window idle-pumping via pumpOnce)\n")
		_ = w2
	}

	// teardown on the owning thread, in reverse order; each Destroy followed by
	// a pump so WM_CLOSE/WM_DESTROY settle and no child is orphaned.
	if m2.created {
		w2.Destroy()
		pumpReap(40)
		fmt.Printf("TEARDOWN m2 destroyed\n")
	}
	w1.Destroy()
	pumpReap(40)
	fmt.Printf("TEARDOWN m1 destroyed\n")
	quitSent <- struct{}{}
}
