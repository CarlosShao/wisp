//go:build windows

package main

// 33-p1 probe, question 2: does the "Go -> page" hop have a delivery agent?
//
// Ticket 33 AC#14 / docs/PLAN.md C17 ("Go -> front-end event push; replies must
// be routed by correlationId"). 33-a2 read the dependency and inferred that
// webview.go's dispatchq is only drained inside Run(); it could not run, so the
// hop stays unproven. This program settles it with a real window.
//
// The page calls a bound Go function and AWAITS its return value. Only the page
// can tell whether the reply arrived, so the page reports the awaited value back
// through a SECOND binding. Go-side "the binding was called" is NOT the judge --
// that is exactly what the shipping firstRoundTripLocked already proves.
//
// Two forms, one flag:
//
//	-selfpump : host pumps like cmd/wisp/panel_host_windows.go pnlPumpOnce does
//	            (PeekMessageW(0,0,0,PM_REMOVE) + TranslateMessage + DispatchMessageW).
//	            This is ticket 33's current shape.
//	-run      : host thread enters the library's own Run() message loop.
//
// The selfpump pump additionally DECODES msg.hwnd and msg.message, which the
// shipping pump never reads, to count how many WM_APP thread messages the pump
// dequeues and whether any of them carries a window handle.
//
// Probe only. Not in the default build graph (.scratch/ is outside the "./..."
// match), not scanned by d22scan (its Go scope is internal/ and cmd/ only).

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

const (
	wmApp    = 0x8000 // internal/w32/w32.go:96
	wmQuit   = 0x0012 // internal/w32/w32.go:92
	pmRemove = 0x0001
)

var (
	modUser32 = windows.NewLazySystemDLL("user32.dll")
	pPeekMsgW = modUser32.NewProc("PeekMessageW")
	pTransMsg = modUser32.NewProc("TranslateMessage")
	pDispMsgW = modUser32.NewProc("DispatchMessageW")
	pPostThrW = modUser32.NewProc("PostThreadMessageW")
)

// pnlMsg shape, copied field-for-field from cmd/wisp/panel_host_windows.go:89-97
// (itself a mirror of internal/ball/win32_windows.go:175-183). Unlike the
// shipping pump this one reads hwnd and message - that is the measurement.
type pnlMsg struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
	_       uint32
}

type counters struct {
	echoCalls       int
	wmAppDequeued   int
	wmAppWithHwnd   int
	otherDequeued   int
	roundsWant      int
	reports         []string
	firstReportAtMS float64
	t0              time.Time
}

func main() {
	selfpump := flag.Bool("selfpump", false, "host pumps Peek/Translate/Dispatch (ticket 33 shape)")
	runForm := flag.Bool("run", false, "host thread enters the library Run() loop")
	rounds := flag.Int("rounds", 3, "sequential awaited binding calls")
	evalCheck := flag.Bool("evalcheck", false, "selfpump form, but the page reports whether a direct Go-side Eval reaches it")
	deadlineS := flag.Int("deadline", 30, "max seconds to observe")
	flag.Parse()

	form := "run"
	if *selfpump {
		form = "selfpump"
	}
	if *selfpump && *runForm {
		fmt.Println("FORM-CONFLICT: pick exactly one of -selfpump / -run")
		os.Exit(2)
	}
	if !*selfpump && !*runForm {
		fmt.Println("FORM-MISSING: -selfpump or -run is required")
		os.Exit(2)
	}

	// The watchdog never touches window or COM state; it only posts WM_QUIT to
	// the UI thread so the Run() loop can end. It goes through the observe
	// registry because ban #1 forbids bare goroutines in production, and this
	// probe wants the same discipline (name deliberately outside the D38 roster).
	reg := observe.NewRegistry()
	uiTID := make(chan uint32, 1)
	stop := make(chan struct{})
	quitSent := make(chan string, 1)
	reg.Spawn("probe33p1-watchdog", "33-p1", nil, func(_ context.Context) {
		tid := <-uiTID
		select {
		case <-stop:
			quitSent <- "not-sent (host finished on its own)"
		case <-time.After(time.Duration(*deadlineS) * time.Second):
			r, _, e := pPostThrW.Call(uintptr(tid), wmQuit, 0, 0)
			quitSent <- fmt.Sprintf("sent-after-%ds-r=%d-%v", *deadlineS, r, e)
		}
	})

	c := &counters{roundsWant: *rounds, t0: time.Now()}
	done := make(chan struct{})

	goHost(form, c, done, uiTID, stop, *deadlineS, *evalCheck)

	select {
	case <-done:
	case <-time.After(time.Duration(*deadlineS+10) * time.Second):
		fmt.Printf("FORM=%s RESULT=HOST-NEVER-FINISHED\n", form)
	}
	close(stop)
	qs := <-quitSent

	fmt.Printf("FORM=%s ECHO_CALLS=%d ROUNDS_WANT=%d WMAPP_DEQUEUED=%d WMAPP_WITH_HWND=%d OTHER_DEQUEUED=%d REPORTS=%d REPORT_FIRST=%q REPORT_LATENCY_MS=%.0f QUIT=%s PANICS=%d\n",
		form, c.echoCalls, c.roundsWant, c.wmAppDequeued, c.wmAppWithHwnd, c.otherDequeued,
		len(c.reports), firstOr(c.reports, "<none>"), c.firstReportAtMS, qs, reg.PanicCount())
	if len(c.reports) > 0 {
		fmt.Printf("FORM=%s REPORT_RAW=%s\n", form, strings.Join(c.reports, "|"))
	}
	verdict(c, form, *evalCheck)
}

func firstOr(rs []string, fb string) string {
	if len(rs) == 0 {
		return fb
	}
	return rs[0]
}

// verdict turns the page's own words into a yes/no on the Go->page hop.
func verdict(c *counters, form string, evalCheck bool) {
	if evalCheck {
		seen := 0
		for _, r := range c.reports {
			if strings.Contains(r, "EVAL_SEEN") {
				seen++
			}
		}
		if len(c.reports) == 0 {
			fmt.Printf("VERDICT %s+evalcheck: page never reported -> BLIND run, no reading\n", form)
		} else if seen > 0 {
			fmt.Printf("VERDICT %s+evalcheck: a direct Go-side Eval DID reach the page (%v) even though the binding reply did not -> push and reply are two different hops\n", form, c.reports)
		} else {
			fmt.Printf("VERDICT %s+evalcheck: even a direct Eval did NOT reach the page -> %v\n", form, c.reports)
		}
		return
	}
	resolved := 0
	for _, r := range c.reports {
		if strings.Contains(r, "echo:") {
			resolved++
		}
	}
	switch {
	case len(c.reports) == 0:
		fmt.Printf("VERDICT %s: page never reported back at all -> the JS side of this probe is BLIND, reading void, do not conclude anything from this run\n", form)
	case resolved == 0:
		fmt.Printf("VERDICT %s: receipt DID NOT reach the page (awaited value never resolved in %d reports)\n", form, len(c.reports))
	default:
		fmt.Printf("VERDICT %s: receipt REACHED the page in %d of %d reports\n", form, resolved, len(c.reports))
	}
}

// goHost pins this goroutine to an OS thread (the thread that will own the
// window), creates the control and then drives it in the requested form.
func goHost(form string, c *counters, done chan struct{}, uiTID chan uint32, stop chan struct{}, deadlineSec int, evalCheck bool) {
	defer close(done)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("HOST-PANIC form=%s: %v\n", form, r)
		}
	}()
	runtime.LockOSThread()

	tid := windows.GetCurrentThreadId()
	uiTID <- tid

	dataPath := filepath.Join(os.TempDir(), "wisp33p1", "receipt-"+form)
	_ = os.MkdirAll(dataPath, 0o755)

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		DataPath:  dataPath,
		AutoFocus: false,
		WindowOptions: webview2.WindowOptions{
			Title:  "33-p1 receipt probe",
			Width:  420,
			Height: 260,
		},
	})
	if w == nil {
		fmt.Printf("FORM=%s RESULT=CREATE-RETURNED-NIL (runtime missing or blocked)\n", form)
		return
	}
	defer func() {
		w.Destroy()
		// Let WM_CLOSE -> DestroyWindow -> WM_DESTROY reach the control so the
		// msedgewebview2 child set exits instead of being orphaned on the box.
		for i := 0; i < 120; i++ {
			pumpOnce(c)
			time.Sleep(25 * time.Millisecond)
		}
	}()

	if err := w.Bind("wispEcho", func(s string) string {
		c.echoCalls++
		return "echo:" + s
	}); err != nil {
		fmt.Printf("FORM=%s RESULT=BIND-ECHO-FAILED: %v\n", form, err)
		return
	}
	if err := w.Bind("wispReport", func(v string) string {
		c.reports = append(c.reports, v)
		if c.firstReportAtMS == 0 {
			c.firstReportAtMS = float64(time.Since(c.t0).Milliseconds())
		}
		if form == "run" {
			// The page has said its piece; end the library loop the same way the
			// library itself does (webview.go:381-383 -> PostQuitMessage).
			w.Terminate()
			return "ok"
		}
		select {
		case stop <- struct{}{}:
		default:
		}
		return "ok"
	}); err != nil {
		fmt.Printf("FORM=%s RESULT=BIND-REPORT-FAILED: %v\n", form, err)
		return
	}

	if evalCheck {
		w.SetHtml(evalPage())
	} else {
		w.SetHtml(page(c.roundsWant))
	}
	c.t0 = time.Now()

	if form == "run" {
		w.Run()
		fmt.Printf("FORM=run RESULT=RUN-RETURNED (WM_QUIT consumed the library loop)\n")
		return
	}

	// selfpump form: the ticket 33 shape, with the message fields decoded so the
	// WM_APP arrivals are visible.
	deadline := time.Now().Add(time.Duration(deadlineSec) * time.Second)
	polls := 0
	evalIssued := false
	for {
		polls++
		pumpOnce(c)
		if evalCheck && !evalIssued && time.Since(c.t0) > 800*time.Millisecond {
			// Go-initiated push straight through browser.Eval (webview.go:439-440),
			// which never goes near dispatchq. Issued ON this thread, because Eval
			// is a direct COM call on the control's owning thread.
			w.Eval("document.title='EVAL-OK-33P1'")
			evalIssued = true
			fmt.Printf("FORM=selfpump EVAL_ISSUED at=%s (direct w.Eval on the UI thread, bypassing Dispatch/dispatchq)\n", time.Since(c.t0).Truncate(time.Millisecond))
		}
		if len(c.reports) >= 1 {
			break
		}
		if time.Now().After(deadline) {
			fmt.Printf("FORM=selfpump RESULT=DEADLINE (no report from the page in %ds, polls=%d)\n", deadlineSec, polls)
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// pumpOnce is the shipping pnlPumpOnce (panel_host_windows.go:400-411) with the
// two fields that pump never reads now decoded and counted.
func pumpOnce(c *counters) {
	var m pnlMsg
	for {
		r, _, _ := pPeekMsgW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, pmRemove)
		if r == 0 {
			return
		}
		switch m.message {
		case wmApp:
			c.wmAppDequeued++
			if m.hwnd != 0 {
				c.wmAppWithHwnd++
			}
		default:
			c.otherDequeued++
		}
		pTransMsg.Call(uintptr(unsafe.Pointer(&m)))
		pDispMsgW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

// evalPage reports whether a Go-side Eval (a push that does NOT ride dispatchq)
// is visible to the page while the host is in the ticket 33 self-pump shape.
func evalPage() string {
	return `<!doctype html><html><head><meta charset="utf-8"><title>untouched</title></head><body><div id="o">eval probe</div>
<script>
var n = 0;
var iv = setInterval(function () {
  n++;
  if (document.title.indexOf("EVAL-OK") >= 0) {
    clearInterval(iv);
    try { window.wispReport("EVAL_SEEN title=" + document.title + " polls=" + n); } catch (e) {}
  } else if (n > 40) {
    clearInterval(iv);
    try { window.wispReport("EVAL_NOT_SEEN title=" + document.title + " polls=" + n); } catch (e) {}
  }
}, 200);
</script></body></html>`
}

// through a second binding. A per-call race against a timer makes the FAILURE
// mode observable: an unresolved promise reports NOT_RESOLVED instead of
// silently hanging the script (which would blind the probe).
func page(rounds int) string {
	return `<!doctype html><html><head><meta charset="utf-8"><title>p1</title></head><body><div id="o">probe</div>
<script>
(async function () {
  var out = [];
  for (var i = 0; i < ` + fmt.Sprint(rounds) + `; i++) {
    var v;
    try {
      var p = window.wispEcho("ping-" + i);
      var to = new Promise(function (res) { setTimeout(function () { res("NOT_RESOLVED"); }, 2000); });
      v = await Promise.race([p, to]);
    } catch (e) { v = "REJECT:" + String(e); }
    out.push(i + "=" + v);
  }
  try { window.wispReport(out.join(";")); } catch (e) { window.wispReport("REPORT-THREW:" + String(e)); }
})();
</script></body></html>`
}
