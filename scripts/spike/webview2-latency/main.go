// webview2-latency measures go-webview2 cold/hot window bring-up latency
// (P11, D32 16.3.4 "panel cold launch <=1500ms / hot <=200ms"):
//
//	cold    = NewWithOptions (creates Win32 window + shows it + Embed waits
//	          for CoreWebView2 environment/controller) + SetHtml + one
//	          Dispatch round trip on the main loop ("usable" point).
//	hot     = ShowWindow(SW_SHOW) of the hidden window + one message pump +
//	          one Dispatch round trip (the C27 single-window reuse path).
//	recreate= destroy + NewWithOptions again in the same process (environment
//	          re-created; NOT shared - pkg/edge builds a fresh environment
//	          per window, which is exactly why C27 single-window reuse exists).
//
// Modes: full (one process, cold + hot x10 + recreate x10), coldonly (one
// cold, for the driver), driver (spawns children for >=10 true colds).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/scripts/spike/common"
	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

type fullRun struct {
	Mode        string             `json:"mode"`
	Machine     common.MachineInfo `json:"machine"`
	StartedAt   string             `json:"startedAtUtc"`
	DataPath    string             `json:"dataPath"`
	ColdMs      float64            `json:"coldMs"`
	ColdNote    string             `json:"coldNote"`
	HotShowMs   []float64          `json:"hotShowMs"`
	HotP50      float64            `json:"hotShowP50Ms"`
	HotP95      float64            `json:"hotShowP95Ms"`
	HotAliveP50 float64            `json:"hotAliveP50Ms"`
	Recreates   []float64          `json:"recreateMs"`
	RecreateP50 float64            `json:"recreateP50Ms"`
	RecreateP95 float64            `json:"recreateP95Ms"`
	Error       string             `json:"error,omitempty"`
}

type coldOnly struct {
	Mode   string  `json:"mode"`
	ColdMs float64 `json:"coldMs"`
	Error  string  `json:"error,omitempty"`
}

type driverOut struct {
	Mode          string    `json:"mode"`
	Children      int       `json:"children"`
	WarmupChildMs float64   `json:"warmupChildMs_discarded"`
	Colds         []float64 `json:"coldMs"`
	ColdP50       float64   `json:"coldP50Ms"`
	ColdP95       float64   `json:"coldP95Ms"`
	FullRun       *fullRun  `json:"fullRun"`
}

var (
	modUser32         = windows.NewLazySystemDLL("user32.dll")
	pShowWindow       = modUser32.NewProc("ShowWindow")
	pUpdateWindow     = modUser32.NewProc("UpdateWindow")
	pPeekMessageW     = modUser32.NewProc("PeekMessageW")
	pTranslateMessage = modUser32.NewProc("TranslateMessage")
	pDispatchMessageW = modUser32.NewProc("DispatchMessageW")
)

const swHide = 0
const swShow = 5

func pumpOnce() {
	type msg struct {
		hwnd    windows.HWND
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}
	var m msg
	for {
		r, _, _ := pPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1)
		if r == 0 {
			return
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func pumpFor(d time.Duration) {
	t0 := time.Now()
	for time.Since(t0) < d {
		pumpOnce()
		time.Sleep(5 * time.Millisecond)
	}
}

// dispatchRT measures one Dispatch round trip on the main loop (proves the
// browser side is alive and the loop is responsive).
func dispatchRT(w webview2.WebView, timeout time.Duration) float64 {
	t0 := time.Now()
	alive := make(chan struct{})
	w.Dispatch(func() { close(alive) })
	for {
		pumpOnce()
		select {
		case <-alive:
			return float64(time.Since(t0).Microseconds()) / 1000.0
		default:
		}
		if time.Since(t0) > timeout {
			return -1
		}
	}
}

func makeWebview(dataPath string) webview2.WebView {
	return webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		DataPath:  dataPath,
		AutoFocus: false,
		WindowOptions: webview2.WindowOptions{
			Title:  "Wisp S0 webview2 spike",
			Width:  420,
			Height: 260,
		},
	})
}

const pageHTML = `<!doctype html><html><head><meta charset="utf-8"><style>body{font-family:Segoe UI,sans-serif;background:#101017;color:#e8e8f0;margin:0;padding:16px}h2{margin:0 0 8px}p{color:#9aa}</style></head><body><h2>L2 confirm card (spike)</h2><p>go-webview2 latency probe</p><button onclick="document.getElementById('x').textContent='clicked'">Allow</button><span id="x"></span></body></html>`

// bringUp measures the full cold path: create+show+embed -> html -> alive.
func bringUp(dataPath string) (webview2.WebView, float64, error) {
	t0 := time.Now()
	w := makeWebview(dataPath)
	if w == nil {
		return nil, 0, fmt.Errorf("NewWithOptions returned nil")
	}
	createMs := float64(time.Since(t0).Microseconds()) / 1000.0

	w.SetHtml(pageHTML)
	pumpFor(120 * time.Millisecond)
	aliveMs := dispatchRT(w, 3*time.Second)
	if aliveMs < 0 {
		return w, createMs, fmt.Errorf("dispatch round trip timed out after create")
	}
	totalMs := float64(time.Since(t0).Microseconds()) / 1000.0
	_ = totalMs
	return w, createMs + aliveMs, nil
}

func pctl(v []float64, p float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	return s[int(math.Round(p*float64(len(s)-1)))]
}

func runFull(dataPath string) fullRun {
	rep := fullRun{Mode: "full", Machine: common.GetMachineInfo(),
		StartedAt: time.Now().UTC().Format(time.RFC3339), DataPath: dataPath,
		ColdNote: "NewWithOptions(create+show+embed) + SetHtml + first dispatch round trip"}

	w, coldMs, err := bringUp(dataPath)
	if err != nil {
		rep.Error = err.Error()
		return rep
	}
	rep.ColdMs = coldMs

	// hot cycles: hide -> show (C27 single-window reuse = the hot path)
	hwnd := uintptr(w.Window())
	for i := 0; i < 10; i++ {
		pShowWindow.Call(hwnd, swHide)
		pumpFor(80 * time.Millisecond)
		t0 := time.Now()
		pShowWindow.Call(hwnd, swShow)
		pUpdateWindow.Call(hwnd)
		pumpOnce()
		showMs := float64(time.Since(t0).Microseconds()) / 1000.0
		aliveMs := dispatchRT(w, 3*time.Second)
		rep.HotShowMs = append(rep.HotShowMs, showMs)
		if i == 0 {
			rep.HotAliveP50 = aliveMs
		}
		pumpFor(80 * time.Millisecond)
	}
	rep.HotP50 = pctl(rep.HotShowMs, 0.5)
	rep.HotP95 = pctl(rep.HotShowMs, 0.95)

	// recreate cycles: destroy -> full create again (fresh environment each
	// time; documents the cost C27 avoids)
	w.Destroy()
	pumpFor(300 * time.Millisecond)
	for i := 0; i < 10; i++ {
		w2, cMs, err := bringUp(dataPath)
		if err != nil {
			rep.Error = "recreate: " + err.Error()
			break
		}
		rep.Recreates = append(rep.Recreates, cMs)
		w2.Destroy()
		pumpFor(300 * time.Millisecond)
	}
	rep.RecreateP50 = pctl(rep.Recreates, 0.5)
	rep.RecreateP95 = pctl(rep.Recreates, 0.95)
	return rep
}

func runColdOnly(dataPath string) coldOnly {
	c := coldOnly{Mode: "coldonly"}
	w, coldMs, err := bringUp(dataPath)
	if err != nil {
		c.Error = err.Error()
		return c
	}
	c.ColdMs = coldMs
	w.Destroy()
	return c
}

func runDriver(children int, dataPath string) driverOut {
	d := driverOut{Mode: "driver", Children: children}
	exe, _ := os.Executable()

	// warmup child: creates the profile dir; discarded from cold stats.
	warm := runChild(exe, dataPath)
	if warm.Error == "" {
		d.WarmupChildMs = warm.ColdMs
	}
	for i := 0; i < children; i++ {
		c := runChild(exe, dataPath)
		if c.Error == "" {
			d.Colds = append(d.Colds, c.ColdMs)
		}
	}
	d.ColdP50 = pctl(d.Colds, 0.5)
	d.ColdP95 = pctl(d.Colds, 0.95)

	d.FullRun = runChildFull(exe, dataPath)
	return d
}

func runChild(exe, dataPath string) coldOnly {
	var c coldOnly
	cmd := exec.Command(exe, "-mode", "coldonly", "-datapath", dataPath)
	out, err := cmd.Output()
	if err != nil {
		c.Error = err.Error()
		return c
	}
	json.Unmarshal(out, &c)
	return c
}

func runChildFull(exe, dataPath string) *fullRun {
	cmd := exec.Command(exe, "-mode", "full", "-datapath", dataPath)
	out, err := cmd.Output()
	if err != nil {
		return &fullRun{Error: err.Error()}
	}
	var f fullRun
	json.Unmarshal(out, &f)
	return &f
}

func main() {
	mode := flag.String("mode", "full", "full | coldonly | driver")
	children := flag.Int("children", 12, "driver: number of cold children")
	dataPath := flag.String("datapath", filepath.Join(os.TempDir(), "wisp-spike-webview2-profile"), "WebView2 user data folder")
	out := flag.String("out", "", "path to write JSON")
	flag.Parse()

	var b []byte
	switch *mode {
	case "full":
		f := runFull(*dataPath)
		b, _ = json.MarshalIndent(f, "", "  ")
	case "coldonly":
		c := runColdOnly(*dataPath)
		b, _ = json.Marshal(c)
	case "driver":
		d := runDriver(*children, *dataPath)
		b, _ = json.MarshalIndent(d, "", "  ")
	}
	fmt.Println(string(b))
	if *out != "" {
		os.WriteFile(*out, b, 0644)
	}
}
