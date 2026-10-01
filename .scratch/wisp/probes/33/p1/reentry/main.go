//go:build windows

package main

// 33-p1 probe, question 1: re-entrant WebView2 bring-up from inside a window
// callback on a thread that is ALREADY pumping messages.
//
// Ticket 33 ruling J1 put the panel host on the ball's ui-sta thread, and
// 33-r1 reported that this panics on this box. 33-a2 read the shape but could
// not run, so three outcomes stayed open: (a) panic, (b) permanent block in the
// library's nested pump, (c) it returns but the message order is wrong /
// callbacks re-enter. This program separates the three with counters instead of
// adjectives.
//
// THREAD SHAPE = [REPLICA] of internal/ball's ui-sta, mirrored line for line
// from sta_windows.go:47-96 / :127-143 and win32_windows.go:99-102,160-183:
//
//	runtime.LockOSThread -> GetCurrentThreadId -> CoInitializeEx(STA)
//	-> RegisterClassExW -> CreateWindowExW -> GetMessageW/Translate/Dispatch loop
//	-> wmAppTask (= WM_APP+0x201) task posting, run inside the window proc
//
// It is NOT the ball's own thread: b.fire / staThread / PostTask are all
// unexported and internal/ball is write-banned for this leg, so there is no
// programmatic way to hang a closure on the real ui-sta. What this measures is
// the physics of the shape, not the fate of that one callback.
//
// The re-entrant call sits where the shipping seam would put it: inside a task
// executed by the window proc, i.e. inside the outer DispatchMessageW frame.
//
// MODES
//
//	-baseline              create with no outer pump on the stack (the nested
//	                       pump is the only pump) -> gives the duration that
//	                       "slow" looks like, the control for branch (b)
//	-reentrant             create from inside the window proc while the outer
//	                       pump is running -> the shipping question
//	-quitduringcreate      inside the window proc, but WM_QUIT is posted first
//	                       so the library's loop breaks at chromium.go:106 with
//	                       e.webview still nil -> forces branch (a) on purpose
//	                       and records where the stack lands
//
// Panic capture uses the repo's own instrument as-is: the thread runs under
// observe.Registry.Spawn, whose run() recovers and hands a PanicEvent carrying
// debug.Stack() to the sink (internal/observe/goroutine.go:294-326, :153, :241).
// Nothing in internal/observe is modified, and no observe test is run here.

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

const (
	wmApp       = 0x8000
	wmAppTask   = wmApp + 0x201 // internal/ball/win32_windows.go:102
	wmDestroy   = 0x0002
	wmNcCreate  = 0x0081
	wmClose     = 0x0010
	wmQuit      = 0x0012
	pmRemove    = 0x0001
	coinitMT    = 0x0000 // COINIT_MULTITHREADED
	coinitSTA   = 0x2    // internal/ball/win32_windows.go:149 coinitApartmentThreaded
	hwndMessage = ^windows.HWND(2)
)

var (
	modUser32    = windows.NewLazySystemDLL("user32.dll")
	modOle32     = windows.NewLazySystemDLL("ole32.dll")
	modKernel32  = windows.NewLazySystemDLL("kernel32.dll")
	pRegClsExW   = modUser32.NewProc("RegisterClassExW")
	pCreateWndW  = modUser32.NewProc("CreateWindowExW")
	pDefWndProcW = modUser32.NewProc("DefWindowProcW")
	pDestroyWnd  = modUser32.NewProc("DestroyWindow")
	pGetMsgW     = modUser32.NewProc("GetMessageW")
	pPeekMsgW    = modUser32.NewProc("PeekMessageW")
	pTransMsg    = modUser32.NewProc("TranslateMessage")
	pDispMsgW    = modUser32.NewProc("DispatchMessageW")
	pPostMsgW    = modUser32.NewProc("PostMessageW")
	pPostQuit    = modUser32.NewProc("PostQuitMessage")
	pPostThrW    = modUser32.NewProc("PostThreadMessageW")
	pCoInitEx    = modOle32.NewProc("CoInitializeEx")
	pCoUninit    = modOle32.NewProc("CoUninitialize")
	pGetModHdlW  = modKernel32.NewProc("GetModuleHandleW")
	pGetCurThrID = modKernel32.NewProc("GetCurrentThreadId")
)

// mirrors internal/ball/win32_windows.go:160-173 (layout, not names).
type wndClassEx struct {
	size      uint32
	style     uint32
	wndProc   uintptr
	clsExtra  int32
	winExtra  int32
	instance  windows.Handle
	icon      windows.Handle
	cursor    windows.Handle
	bkgnd     windows.Handle
	menuName  uintptr
	className *uint16
	iconSm    windows.Handle
}

// mirrors internal/ball/win32_windows.go:175-183.
type staMsg struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	_       uint32
}

// thread is the replica of staThread (its fields, not its type).
type thread struct {
	mu      sync.Mutex
	nextID  uint64
	tasks   map[uint64]func()
	hwnd    windows.HWND
	tid     uint32
	w       webview2.WebView
	seqMu   sync.Mutex
	seq     []string
	outerIt atomic.Int64 // iterations of the OUTER pump
	inside  atomic.Int32 // 1 while the window-proc callback is running
	// The ordering evidence: tasks posted BEFORE the blocking call.
	nestedR  atomic.Int64 // ran while inside==1 -> a pump other than the outer took them
	afterR   atomic.Int64 // ran after the callback returned, in order
	pingRan  atomic.Int64 // watchdog ping: proves a pump is still consuming this queue
	createR  atomic.Int64 // 0 = not returned, 1 = returned a window, 2 = returned nil
	createMs atomic.Int64
}

var (
	probe       *thread
	probedPtr   uintptr
	logf        = stdoutWriter
	stallSecGlb int
)

func stdoutWriter(format string, a ...any) { fmt.Printf(format+"\n", a...) }

func (t *thread) say(format string, a ...any) {
	t.seqMu.Lock()
	defer t.seqMu.Unlock()
	t.seq = append(t.seq, fmt.Sprintf(format, a...))
}

func (t *thread) allSeq() []string {
	t.seqMu.Lock()
	defer t.seqMu.Unlock()
	return append([]string{}, t.seq...)
}

func main() {
	modeBaseline := flag.Bool("baseline", false, "no outer pump on the stack")
	modeReent := flag.Bool("reentrant", false, "create from inside the window proc")
	modeQuit := flag.Bool("quitduringcreate", false, "inside the window proc, WM_QUIT posted first")
	modeSkip := flag.Bool("skipcreate", false, "control: same callback position, no blocking create at all")
	stallSec := flag.Int("stall", 20, "seconds after which a non-returning create is called blocked")
	comMode := flag.String("com", "sta", "COM init on the thread before creating: sta | mta | none")
	logPath := flag.String("log", "", "verbatim event log path")
	flag.Parse()
	if *comMode != "sta" && *comMode != "mta" && *comMode != "none" {
		fmt.Println("COM-ERROR: -com wants sta | mta | none")
		os.Exit(2)
	}

	mode := ""
	for m, on := range map[string]bool{
		"baseline": *modeBaseline, "reentrant": *modeReent,
		"quitduringcreate": *modeQuit, "skipcreate-control": *modeSkip,
	} {
		if on {
			mode = m
		}
	}
	n := 0
	for _, b := range []bool{*modeBaseline, *modeReent, *modeQuit, *modeSkip} {
		if b {
			n++
		}
	}
	if n != 1 {
		fmt.Println("MODE-ERROR: exactly one of -baseline / -reentrant / -quitduringcreate / -skipcreate")
		os.Exit(2)
	}
	stallSecGlb = *stallSec

	if *logPath != "" {
		fh, err := os.Create(*logPath)
		if err != nil {
			fmt.Println("LOG-OPEN-FAILED:", err)
			os.Exit(2)
		}
		logf = func(format string, a ...any) {
			fmt.Fprintf(fh, format+"\n", a...)
			fmt.Fprintf(os.Stdout, format+"\n", a...)
			_ = fh.Sync()
		}
	}

	dir := filepath.Join(os.TempDir(), "wisp33p1", "reentry-"+mode+"-"+*comMode)
	_ = os.MkdirAll(dir, 0o755)
	sinkFile := filepath.Join(dir, "panic-sink.txt")

	reg := observe.NewRegistry()
	var sinkWrites atomic.Int64
	reg.SetPanicSink(func(ev observe.PanicEvent) {
		sinkWrites.Add(1)
		body := fmt.Sprintf("Goroutine=%s Owner=%s Recovered=%s At=%s\nSTACK:\n%s\n",
			ev.Goroutine, ev.Owner, ev.Recovered, ev.At, ev.Stack)
		_ = os.WriteFile(sinkFile, []byte(body), 0o644)
		logf("SINK-EVENT goroutine=%s recovered=%q file=%s", ev.Goroutine, ev.Recovered, sinkFile)
	})

	probe = &thread{tasks: map[uint64]func(){}}
	probedPtr = windows.NewCallback(probeWndProc)

	h := reg.Spawn("probe33p1-usta", "33-p1", nil, func(_ context.Context) { runSTA(mode, dir, *comMode) })

	watchdogExit := make(chan struct{})
	reg.Spawn("probe33p1-stallwatch", "33-p1", nil, func(_ context.Context) {
		defer close(watchdogExit)
		tick := time.NewTicker(500 * time.Millisecond)
		defer tick.Stop()
		start := time.Now()
		for range tick.C {
			if probe.createR.Load() != 0 {
				return // it came back (or came back nil); nothing to classify
			}
			elapsed := time.Since(start)
			if elapsed < time.Duration(*stallSec)*time.Second {
				continue
			}
			// Still not back. Is any pump on this thread still consuming messages?
			before := probe.pingRan.Load()
			probe.postTask(func() { probe.pingRan.Add(1) })
			time.Sleep(2 * time.Second)
			after := probe.pingRan.Load()
			logf("WATCHDOG mode=%s create_returned=false waited=%s tasks_taken_by_nested_pump=%d ping_delta=%d outer_pump_iters=%d",
				mode, elapsed.Truncate(time.Millisecond), probe.nestedR.Load(), after-before, probe.outerIt.Load())
			r, _, e := pPostThrW.Call(uintptr(probe.tid), wmQuit, 0, 0)
			logf("WATCHDOG abort: posted WM_QUIT to tid=%d r=%d err=%v (a FORCED stop, not a natural return)", probe.tid, r, e)
			return
		}
	})

	select {
	case <-h.Done():
		err := h.Err()
		logf("STA-HANDLE-RETURNED err=%v", err)
	case <-time.After(120 * time.Second):
		logf("STA-HANDLE-TIMEOUT 120s (thread never came back even after the watchdog)")
	}
	select {
	case <-watchdogExit:
	case <-time.After(5 * time.Second):
		logf("WATCHDOG-STILL-ALIVE")
	}

	logf("FINAL mode=%s create_returned=%v create_nil=%v create_ms=%d tasks_taken_by_nested_pump=%d tasks_after_callback=%d outer_iters=%d ping=%d panics_total=%d sink_writes=%d sink_file=%s",
		mode, probe.createR.Load() == 1, probe.createR.Load() == 2, probe.createMs.Load(),
		probe.nestedR.Load(), probe.afterR.Load(), probe.outerIt.Load(), probe.pingRan.Load(),
		reg.PanicCount(), sinkWrites.Load(), sinkFile)
	for _, s := range probe.allSeq() {
		logf("SEQ %s", s)
	}
	switch mode {
	case "skipcreate-control":
		if probe.nestedR.Load() == 0 && probe.afterR.Load() == 5 {
			logf("VERDICT skipcreate-control: ORDER KEPT - all 5 tasks ran after the callback, the nested pump took none. This is the shape the reentrant run SHOULD have looked like.")
		} else {
			logf("VERDICT skipcreate-control: UNEXPECTED nested=%d after=%d - the ordering instrument itself is suspect, re-read every reentrant number with that in mind",
				probe.nestedR.Load(), probe.afterR.Load())
		}
	case "baseline":
		if probe.createR.Load() == 1 {
			logf("VERDICT baseline: bring-up with this thread as the ONLY pump returned in %dms (sink=%d panic_count=%d)",
				probe.createMs.Load(), sinkWrites.Load(), reg.PanicCount())
		} else {
			logf("VERDICT baseline: bring-up did not return even with no outer pump on the stack -> the re-entry question is moot, look elsewhere")
		}
	case "reentrant":
		switch {
		case sinkWrites.Load() > 0 || reg.PanicCount() > 0:
			logf("VERDICT reentrant: branch (a) PANIC from the natural shape (no WM_QUIT was forced) -> sink file %s", sinkFile)
		case probe.createR.Load() == 0:
			logf("VERDICT reentrant: branch (b) BLOCK - bring-up never returned in %d stall seconds; the watchdog then had to force a WM_QUIT", *stallSec)
		case probe.nestedR.Load() > 0:
			logf("VERDICT reentrant: branch (c) ORDER-REVERSAL - it returned in %dms, but %d tasks posted before it ran INSIDE the nested pump (outer pump frozen at iter=%d the whole time)",
				probe.createMs.Load(), probe.nestedR.Load(), probe.outerIt.Load())
		default:
			logf("VERDICT reentrant: NONE of the three - it returned in %dms and the outer queue was untouched", probe.createMs.Load())
		}
	case "quitduringcreate":
		if sinkWrites.Load() > 0 {
			logf("VERDICT quitduringcreate: branch (a) PANIC reachable, and the observe sink captured debug.Stack() -> %s", sinkFile)
		} else {
			logf("VERDICT quitduringcreate: no sink event; either the panic escaped the recover boundary (see stderr + exit code) or the forced WM_QUIT took another path")
		}
	}
}

// runSTA mirrors sta_windows.go:47-96 shape, then drives the mode.
func runSTA(mode, dir, com string) {
	defer pCoUninit.Call()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tid, _, _ := pGetCurThrID.Call()
	probe.tid = uint32(tid)
	logf("STA tid=%d locked=true com=%s", probe.tid, com)

	switch com {
	case "sta":
		r1, _, _ := pCoInitEx.Call(0, coinitSTA)
		logf("CoInitializeEx(0x2 COINIT_APARTMENTTHREADED) r=%d", r1)
	case "mta":
		r1, _, _ := pCoInitEx.Call(0, coinitMT)
		logf("CoInitializeEx(0x0 COINIT_MULTITHREADED) r=%d", r1)
	default:
		logf("CoInitializeEx SKIPPED (com=none): this thread never asked for an apartment")
	}

	cls := windows.StringToUTF16Ptr("Wisp33P1ReplicaSTA")
	wc := wndClassEx{
		size:      uint32(unsafe.Sizeof(wndClassEx{})),
		wndProc:   probedPtr,
		instance:  windows.Handle(moduleHandle()),
		className: cls,
	}
	if r, _, e := pRegClsExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		logf("RegisterClassExW FAILED r=%d err=%v - aborting replica", r, e)
		return
	}
	// Message-only window: a recorded deviation from the ball's 1x1 WS_POPUP
	// top-level. Same thread queue, same pump, same window-proc dispatch path,
	// but no desktop surface for the replica itself.
	hwnd, _, err := pCreateWndW.Call(
		0, uintptr(unsafe.Pointer(cls)), 0,
		0, 0, 0, 0, 0,
		uintptr(hwndMessage), 0, uintptr(moduleHandle()), 0)
	if hwnd == 0 {
		logf("CreateWindowExW(message-only) FAILED err=%v - aborting replica", err)
		return
	}
	probe.hwnd = windows.HWND(hwnd)
	logf("replica window hwnd=%p class=Wisp33P1ReplicaSTA message_only=true", probe.hwnd)

	if mode == "baseline" {
		// No outer pump on the stack: call bring-up directly, then start the
		// pump afterwards only to reap the window.
		probeCallback(mode, dir)()
	} else {
		probe.postTask(probeCallback(mode, dir))
	}

	// The outer pump - the co-resident the library's nested pump would have to
	// share. Identical shape to sta_windows.go:78-92.
	var m staMsg
	for {
		r, _, _ := pGetMsgW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		c := int32(r)
		if c == 0 {
			probe.say("outer GetMessageW returned 0 (WM_QUIT) at iter=%d", probe.outerIt.Load())
			break
		}
		if c == -1 {
			probe.say("outer GetMessageW error at iter=%d", probe.outerIt.Load())
			break
		}
		probe.outerIt.Add(1)
		pTransMsg.Call(uintptr(unsafe.Pointer(&m)))
		pDispMsgW.Call(uintptr(unsafe.Pointer(&m)))
	}

	// Reap the control so no msedgewebview2 set is orphaned on this box (other
	// legs count the process tree; a probe that leaks children is a lie).
	if probe.w != nil {
		probe.w.Destroy()
		for i := 0; i < 120; i++ {
			var pm staMsg
			r, _, _ := pPeekMsgW.Call(uintptr(unsafe.Pointer(&pm)), 0, 0, 0, pmRemove)
			if r == 0 {
				time.Sleep(25 * time.Millisecond)
				continue
			}
			pTransMsg.Call(uintptr(unsafe.Pointer(&pm)))
			pDispMsgW.Call(uintptr(unsafe.Pointer(&pm)))
		}
	}
	pDestroyWnd.Call(uintptr(probe.hwnd))
	probe.say("replica window destroyed")
}

// probeCallback is the task that runs inside the window proc, i.e. inside the
// outer DispatchMessageW frame - the position the shipping seam would put
// bringUp in.
func probeCallback(mode, dir string) func() {
	return func() {
		probe.inside.Store(1)
		probe.say("CALLBACK enter outer_iters=%d tid=%d mode=%s", probe.outerIt.Load(), probe.tid, mode)

		// Five follow-up tasks posted BEFORE the blocking call. Whether they run
		// before bring-up returns is the ordering question.
		for i := 1; i <= 5; i++ {
			i := i
			probe.postTask(func() {
				if probe.inside.Load() == 1 {
					probe.nestedR.Add(1)
					probe.say("task P%d ran INSIDE the callback (a pump other than the outer took it) outer_iters=%d",
						i, probe.outerIt.Load())
				} else {
					probe.afterR.Add(1)
					probe.say("task P%d ran after the callback, in order (outer_iter=%d)", i, probe.outerIt.Load())
				}
			})
		}

		if mode == "skipcreate-control" {
			// Control for the ordering claim: identical task posting and identical
			// callback position, but nothing blocking in between. If the shape is
			// honest, all five tasks must run AFTER the callback, in order.
			probe.createR.Store(1)
			probe.createMs.Store(0)
			probe.say("SKIP-CREATE control: no bring-up issued; the outer pump should keep the order")
			probe.inside.Store(0)
			probe.say("CALLBACK exit outer_iters=%d", probe.outerIt.Load())
			pPostQuit.Call(0)
			return
		}

		if mode == "quitduringcreate" {
			// Force chromium.go:106 (r == 0 -> break) before the COM callback can
			// set inited, so :112 -> :130-131 runs with e.webview still nil.
			r, _, e := pPostThrW.Call(uintptr(probe.tid), wmQuit, 0, 0)
			probe.say("FORCED WM_QUIT posted to tid=%d r=%d err=%v (artificial: nothing in the resident topology posts it here at this moment)", probe.tid, r, e)
		}

		t0 := time.Now()
		w := webview2.NewWithOptions(webview2.WebViewOptions{
			Debug:     false,
			DataPath:  dir,
			AutoFocus: false,
			WindowOptions: webview2.WindowOptions{
				Title:  "33-p1 reentry " + mode,
				Width:  420,
				Height: 260,
			},
		})
		ms := time.Since(t0).Milliseconds()
		probe.createMs.Store(ms)
		if w == nil {
			probe.createR.Store(2)
			probe.say("bring-up RETURNED nil after %dms (host-side error shape, no panic)", ms)
			probe.inside.Store(0)
			return
		}
		probe.w = w
		probe.createR.Store(1)
		probe.say("bring-up RETURNED a window after %dms; still inside the same callback frame=%v", ms, probe.inside.Load() == 1)

		// "Driving" it from the same frame, the way bringUp does after create.
		if berr := w.Bind("wispPing", func(s string) string { return "pong:" + s }); berr != nil {
			probe.say("Bind from inside the callback FAILED: %v", berr)
		} else {
			probe.say("Bind from inside the callback OK")
		}
		w.SetHtml(`<!doctype html><html><body><div id="o">33-p1 reentry</div></body></html>`)
		probe.say("SetHtml from inside the callback issued")
		probe.inside.Store(0)
		probe.say("CALLBACK exit outer_iters=%d", probe.outerIt.Load())

		pPostQuit.Call(0) // let this thread's pump end so the replica can reap
	}
}

// postTask mirrors sta_windows.go:127-143: assign an id, store the closure,
// PostMessageW with the id in wParam.
func (t *thread) postTask(fn func()) {
	t.mu.Lock()
	t.nextID++
	id := t.nextID
	t.tasks[id] = fn
	t.mu.Unlock()
	if t.hwnd == 0 {
		return
	}
	pPostMsgW.Call(uintptr(t.hwnd), wmAppTask, uintptr(id), 0)
}

func (t *thread) runTask(id uint64) {
	t.mu.Lock()
	fn := t.tasks[id]
	delete(t.tasks, id)
	t.mu.Unlock()
	if fn != nil {
		fn()
	}
}

func probeWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmAppTask:
		probe.runTask(uint64(wParam))
		return 0
	case wmClose:
		pDestroyWnd.Call(hwnd)
		return 0
	case wmDestroy:
		return 0
	}
	r, _, _ := pDefWndProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func moduleHandle() uintptr {
	h, _, _ := pGetModHdlW.Call(0)
	return h
}
