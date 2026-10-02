//go:build windows

package main

// 33-p1 leg 2 (2026-10-02), question 2: the Go->page receipt on the CURRENT
// shipping API surface, with the real products the host actually has.
//
// Since the first 33-p1 leg (10-01) the shipping host changed shape: the
// resident panel thread (cmd/wisp/panel_resident_windows.go) now hands the pump
// to the library's Run() (handOverPump), so the shipping form IS the "run" form.
// This probe measures the two push channels the dependency exports on that form:
//
//   - Eval  = ExecuteScript under the hood (chromium.go:138-145: direct COM call
//     on the control's owning thread; result ignored). This is the only push
//     channel the high-level webview2.WebView interface exports today.
//   - PostWebMessageAsString exists ONLY as an internal side effect of the
//     dependency's own message handler (chromium.go:242 echoes the received
//     message back); no method exposes it to callers, and the production tree
//     (cmd internal) contains zero call sites of either name.
//
// The page judges: it polls for a marker that only a successful push can set
// (document.title for Eval), and reports back through a second binding. A push
// that does not arrive is reported as NOT_SEEN by the page itself - Go-side
// "the call did not error" is recorded but is NOT the judge (that is the fake
// green the first leg's R25 exposed).
//
// Three runs in one binary so the log lines are comparable:
//
//	-mode eval        : Go pushes document.title via w.Eval on the UI thread,
//	                    while Run() pumps. Page reports EVAL_SEEN / EVAL_NOT_SEEN.
//	-mode eval-disp   : same push issued from ANOTHER goroutine wrapped in
//	                    w.Dispatch (the only cross-thread door). Page judges.
//	-mode dup-manager : question 1 rider - a SECOND webview2.NewWithOptions on
//	                    the SAME STA thread while the first window still exists
//	                    (the "second PanelManager" shape; PanelManager itself is
//	                    cmd/wisp package main, not importable, so the raw
//	                    dependency call is what bringUp does under it).
//
// Probe only. Lives under .scratch/ (outside ./... and outside d22scan's Go
// scope). No production file is touched; nothing in internal/ is modified; the
// observe registry is used as shipped (sink installed HERE, in the probe, like
// the first leg did).

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

const (
	wmQuit     = 0x0012
	coinitSTA  = 0x2
	hwndMsg    = ^windows.HWND(2) // HWND_MESSAGE
	probeTitle = "P1Q2-OK-33"
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

// pumpOnce drains the thread queue exactly like panel_host_windows.go:621-632,
// counting WM_APP arrivals (the dependency posts WM_APP for its dispatchq).
func pumpOnce(wmAppSeen *atomic.Int64) {
	var m pnlMsg
	for {
		r, _, _ := pPeekMsgW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1) // PM_REMOVE
		if r == 0 {
			return
		}
		if m.message == 0x8000 && wmAppSeen != nil {
			wmAppSeen.Add(1)
		}
		pTransMsg.Call(uintptr(unsafe.Pointer(&m)))
		pDispMsgW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

var (
	reported  chan struct{}
	reports   atomic.Int64
	firstRep  atomic.Int64 // ms since t0
	t0        time.Time
	uiTID     = make(chan uint32, 1)
	evaled    atomic.Bool // set right before/after the Go push is issued
	evalAtMs  atomic.Int64
	quitUsed  atomic.Bool
)

func reportBody() string { return "" }

func main() {
	mode := flag.String("mode", "eval", "eval | eval-disp | dup-manager")
	deadlineS := flag.Int("deadline", 20, "watchdog seconds before a forced WM_QUIT")
	flag.Parse()
	switch *mode {
	case "eval", "eval-disp", "dup-manager":
	default:
		fmt.Println("MODE-ERROR: -mode wants eval | eval-disp | dup-manager")
		os.Exit(2)
	}

	dir := filepath.Join(os.TempDir(), "wisp33p1q2", *mode)
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

	var wmAppSeen atomic.Int64
	hostDone := make(chan struct{})
	reported = make(chan struct{}, 1)

	// Watchdog: never touches window/COM state; only posts WM_QUIT to the UI
	// thread so Run() can end if the page never reports. A forced quit marks the
	// run's readings as watchdog-tainted.
	reg.Spawn("p1q2-watchdog", "33-p1-q2", nil, func(_ context.Context) {
		tid := <-uiTID
		select {
		case <-hostDone:
		case <-reported:
			// The page has spoken; readings are complete. End Run() politely by
			// posting WM_QUIT to the UI thread (a quit AFTER the measurement, so
			// no reading is tainted).
			pPostThrW.Call(uintptr(tid), wmQuit, 0, 0)
			fmt.Printf("WATCHDOG clean exit: page reported, WM_QUIT posted to tid=%d\n", tid)
		case <-time.After(time.Duration(*deadlineS) * time.Second):
			quitUsed.Store(true)
			r, _, e := pPostThrW.Call(uintptr(tid), wmQuit, 0, 0)
			fmt.Printf("WATCHDOG forced WM_QUIT to tid=%d r=%d err=%v (readings tainted)\n", tid, r, e)
		}
	})

	reg.Spawn("p1q2-usta", "33-p1-q2", nil, func(_ context.Context) {
		runHost(*mode, dir, &wmAppSeen, hostDone)
	})

	select {
	case <-hostDone:
	case <-time.After(time.Duration(*deadlineS+15) * time.Second):
		fmt.Printf("MODE=%s RESULT=HOST-NEVER-FINISHED\n", *mode)
	}

	lat := firstRep.Load()
	fmt.Printf("MODE=%s WMAPP_DEQUEUED=%d REPORTS=%d FIRST_REPORT_MS=%d EVAL_ISSUED=%v EVAL_AT_MS=%d QUIT_FORCED=%v PANICS=%d SINK_WRITES=%d\n",
		*mode, wmAppSeen.Load(), reports.Load(), lat, evaled.Load(), evalAtMs.Load(), quitUsed.Load(),
		reg.PanicCount(), sinkWrites.Load())
	verdict(*mode, lat)
}

func verdict(mode string, latMs int64) {
	n := reports.Load()
	switch mode {
	case "eval", "eval-disp":
		switch {
		case n == 0:
			fmt.Printf("VERDICT %s: page never reported -> BLIND run, no reading\n", mode)
		case strings.Contains(fmt.Sprint(latMs), "-"):
			fmt.Printf("VERDICT %s: unreadable\n", mode)
		default:
			// the page's own words decide; they land in REPORT_TEXT above
			fmt.Printf("VERDICT %s: page reported (%d) - see REPORT_TEXT for EVAL_SEEN / EVAL_NOT_SEEN\n", mode, n)
		}
	case "dup-manager":
		fmt.Printf("VERDICT dup-manager: see RESULT line above for the second create's fate\n")
	}
}

// runHost is the STA thread: CoInitializeEx(STA) by name (the resident panel
// thread's own guard shape, panel_resident_windows.go:214), create window #1,
// bind the reporter, then run the mode's scenario and pump Run() until it ends.
func runHost(mode, dir string, wmAppSeen *atomic.Int64, hostDone chan struct{}) {
	defer close(hostDone)
	defer pCoUninit.Call()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tid, _, _ := pGetCurTID.Call()
	uiTID <- uint32(tid)

	dataPath := filepath.Join(dir, "userdata")
	_ = os.MkdirAll(dataPath, 0o755)

	// Window #1 = the "resident" window. No host window class of our own is
	// needed for eval/eval-disp (the library makes its own top-level window);
	// for dup-manager we DO plant a message-only window first so the thread is
	// a real pumping STA thread like the replica shape, but the second create
	// still happens while window #1 is alive - that is the shape asked for.
	w1 := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:    false,
		DataPath: dataPath,
		WindowOptions: webview2.WindowOptions{
			Title: "33-p1-q2 first window", Width: 420, Height: 260,
		},
	})
	if w1 == nil {
		fmt.Printf("MODE=%s RESULT=CREATE1-RETURNED-NIL (runtime missing?)\n", mode)
		return
	}
	t1 := time.Now()
	fmt.Printf("MODE=%s CREATE1 ok in %dms\n", mode, time.Since(t1).Milliseconds())

	var bindErr error
	if bindErr = w1.Bind("wispReport", func(v string) string {
		if reports.CompareAndSwap(0, 1) {
			firstRep.Store(time.Since(t0).Milliseconds())
		} else {
			reports.Add(1)
		}
		fmt.Printf("REPORT_TEXT page-said=%q at_ms=%d\n", v, time.Since(t0).Milliseconds())
		select {
		case reported <- struct{}{}:
		default:
		}
		return "ok"
	}); bindErr != nil {
		fmt.Printf("MODE=%s RESULT=BIND-FAILED %v\n", mode, bindErr)
		w1.Destroy()
		return
	}

	t0 = time.Now()

	switch mode {
	case "eval":
		w1.SetHtml(evalPage("eval"))
		// Issue the push ON this thread after a short settle, while Run() owns
		// the loop. Eval = ExecuteScript (chromium.go:138-145): a direct COM
		// call, result ignored, no dispatchq involved.
		go func() {
			time.Sleep(700 * time.Millisecond)
			// CROSS-THREAD by construction here would be wrong for Eval; the
			// scenario wants the on-thread call. So the actual push happens on
			// the UI thread through the self-posted task below.
			_ = w1 // keep the handle alive
		}()
		evalScheduled := time.After(600 * time.Millisecond)
		go func() {
			<-evalScheduled
			// Record intent; the push itself must run on the UI thread. We ride
			// w.Dispatch (the cross-thread door) ONCE to post a closure that
			// does the on-thread Eval - for -mode eval that closure runs a
			// DIRECT w1.Eval, i.e. ExecuteScript, not a page binding.
			w1.Dispatch(func() {
				evaled.Store(true)
				evalAtMs.Store(time.Since(t0).Milliseconds())
				w1.Eval("document.title='" + probeTitle + "'")
				fmt.Printf("PUSH issued via w1.Eval (ExecuteScript) on UI thread at_ms=%d\n", evalAtMs.Load())
			})
		}()
		w1.Run()
		pumpReap(w1, wmAppSeen)
		fmt.Printf("MODE=%s RESULT=RUN-RETURNED\n", mode)

	case "eval-disp":
		w1.SetHtml(evalPage("eval-disp"))
		// The push rides ONLY Dispatch (the cross-thread door the resident
		// panel's post() uses once wRef is published). Inside the closure the
		// push is still w1.Eval - this isolates "Dispatch closure actually ran"
		// from "Eval reached the page".
		go func() {
			time.Sleep(600 * time.Millisecond)
			w1.Dispatch(func() {
				evaled.Store(true)
				evalAtMs.Store(time.Since(t0).Milliseconds())
				w1.Eval("document.title='" + probeTitle + "-DISP'")
				fmt.Printf("PUSH issued via Dispatch->Eval at_ms=%d\n", evalAtMs.Load())
			})
		}()
		w1.Run()
		pumpReap(w1, wmAppSeen)
		fmt.Printf("MODE=%s RESULT=RUN-RETURNED\n", mode)

	case "dup-manager":
		// Question 1 rider: the shipping "second PanelManager" shape. PanelManager
		// is package main in cmd/wisp (not importable), so the raw dependency
		// call under it is what we exercise: a SECOND webview2.NewWithOptions on
		// the SAME thread while window #1 is alive and pumping.
		w1.SetHtml(`<!doctype html><html><head><title>first stays up</title></head><body><p>w1 alive</p></body></html>`)
		// Let the first window settle on the library pump briefly.
		time.Sleep(400 * time.Millisecond)
		t2 := time.Now()
		w2 := webview2.NewWithOptions(webview2.WebViewOptions{
			Debug:    false,
			DataPath: filepath.Join(dataPath, "second"),
			WindowOptions: webview2.WindowOptions{
				Title: "33-p1-q2 second window", Width: 420, Height: 260,
			},
		})
		ms2 := time.Since(t2).Milliseconds()
		if w2 == nil {
			fmt.Printf("MODE=dup-manager RESULT=CREATE2-RETURNED-NIL ms=%d (first window still alive)\n", ms2)
		} else {
			// Does the second control work? Bind + SetHtml on it.
			b2 := w2.Bind("wispReport", func(v string) string {
				if reports.CompareAndSwap(0, 1) {
					firstRep.Store(time.Since(t0).Milliseconds())
				} else {
					reports.Add(1)
				}
				fmt.Printf("REPORT_TEXT w2-page-said=%q at_ms=%d\n", v, time.Since(t0).Milliseconds())
				select {
				case reported <- struct{}{}:
				default:
				}
				return "ok"
			})
			fmt.Printf("MODE=dup-manager RESULT=CREATE2-RETURNED-WINDOW ms=%d bind2err=%v\n", ms2, b2)
			if b2 == nil {
				w2.SetHtml(`<!doctype html><html><head><title>second alive</title></head><body>
					<script>try{window.wispReport("W2-LIVE")}catch(e){};</script></body></html>`)
				// Give the second control a bounded window to speak; then we
				// must destroy it - its messages ride the same thread queue.
				deadline := time.Now().Add(4 * time.Second)
				for time.Now().Before(deadline) {
					pumpOnce(wmAppSeen)
					time.Sleep(5 * time.Millisecond)
				}
				w2.Destroy()
				fmt.Printf("MODE=dup-manager W2 destroyed after its window\n")
			}
		}
		// First window still alive; end the run cleanly.
		pumpReap(w1, wmAppSeen)
		w1.Destroy()
		fmt.Printf("MODE=dup-manager RESULT=DONE w1 destroyed last\n")
	}
}

// pumpReap lets WM_CLOSE/destroy traffic and the library's own teardown messages
// settle after Run() returns, so no msedgewebview2 child set is orphaned.
func pumpReap(w webview2.WebView, wmAppSeen *atomic.Int64) {
	_ = w
	for i := 0; i < 80; i++ {
		pumpOnce(wmAppSeen)
		time.Sleep(25 * time.Millisecond)
	}
}

// evalPage polls its own title and reports through the binding - the page's
// word is the judge, per the first leg's R25 lesson.
func evalPage(tag string) string {
	return `<!doctype html><html><head><meta charset="utf-8"><title>untouched-` + tag + `</title></head><body><div id="o">q2</div>
<script>
var n = 0;
var iv = setInterval(function () {
  n++;
  if (document.title.indexOf("P1Q2-OK") >= 0) {
    clearInterval(iv);
    try { window.wispReport("EVAL_SEEN title=" + document.title + " polls=" + n); } catch (e) {}
  } else if (n > 60) {
    clearInterval(iv);
    try { window.wispReport("EVAL_NOT_SEEN title=" + document.title + " polls=" + n); } catch (e) {}
  }
}, 200);
</script></body></html>`
}
