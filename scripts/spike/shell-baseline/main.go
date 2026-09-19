// shell-baseline measures RSS baselines 1/4/5 (pure Go, NO cgo):
//
//	(1) empty Go process
//	(4) + layered window + Direct2D + DirectWrite (real visible window, one
//	    frame presented)
//	(5) + tray icon + global hotkey + Job Object (KILL_ON_JOB_CLOSE)
//
// Stages are cumulative in one process: each stage reports the absolute
// private working set and the delta to the previous stage.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/CarlosShao/wisp/scripts/spike/common"
	"golang.org/x/sys/windows"
)

type stage struct {
	Stage           string  `json:"stage"`
	PrivateWS_MB    float64 `json:"privateWS_MB"`
	SharedWS_MB     float64 `json:"sharedWS_MB"`
	TotalWS_MB      float64 `json:"totalWS_MB"`
	PrivateCommitMB float64 `json:"privateCommit_MB"`
	GDIObjects      uint32  `json:"gdiObjects"`
	UserObjects     uint32  `json:"userObjects"`
	Handles         uint32  `json:"handles"`
	Threads         uint32  `json:"threads"`
	DeltaPrev_MB    float64 `json:"deltaPrev_MB"`
	MinPrivateWSMB  float64 `json:"minPrivateWS_MB"`
	MaxPrivateWSMB  float64 `json:"maxPrivateWS_MB"`
	Error           string  `json:"error,omitempty"`
	DurationMs      int64   `json:"stageDurationMs,omitempty"`
}

type result struct {
	Program   string             `json:"program"`
	Kind      string             `json:"kind"` // pure-go | cgo
	Machine   common.MachineInfo `json:"machine"`
	StartedAt string             `json:"startedAtUtc"`
	Stages    []stage            `json:"stages"`
	GoVersion string             `json:"goVersion"`
	GOGC      string             `json:"gogc"`
	Note      string             `json:"note"`
}

func mbOf(b uint64) float64 { return float64(b) / (1 << 20) }

func sampleStage(name, kind string, prev *stage, run func() error) stage {
	st := stage{Stage: name}
	t0 := time.Now()
	var runErr error
	if run != nil {
		runErr = run()
	}
	if runErr != nil {
		st.Error = runErr.Error()
	}
	common.SettleGC()
	minS, medS, maxS := common.SampleStable(5, 200)
	st.PrivateWS_MB = mbOf(medS.PrivateWorkingSet)
	st.SharedWS_MB = mbOf(medS.SharedWorkingSet)
	st.TotalWS_MB = mbOf(medS.WorkingSet)
	st.PrivateCommitMB = mbOf(medS.PrivateCommit)
	st.GDIObjects = medS.GDIObjects
	st.UserObjects = medS.UserObjects
	st.Handles = medS.Handles
	st.Threads = medS.Threads
	st.MinPrivateWSMB = mbOf(minS.PrivateWorkingSet)
	st.MaxPrivateWSMB = mbOf(maxS.PrivateWorkingSet)
	if prev != nil {
		st.DeltaPrev_MB = st.PrivateWS_MB - prev.PrivateWS_MB
	}
	st.DurationMs = time.Since(t0).Milliseconds()
	_ = kind
	return st
}

func main() {
	out := flag.String("out", "", "path to write JSON (also printed to stdout)")
	flag.Parse()

	res := result{
		Program:   "shell-baseline",
		Kind:      "pure-go",
		Machine:   common.GetMachineInfo(),
		StartedAt: time.Now().UTC().Format(time.RFC3339),
		GoVersion: "go1.27.1",
		GOGC:      os.Getenv("GOGC"),
		Note:      "stages cumulative; privateWS = Task Manager private working set (QueryWorkingSetEx)",
	}

	var win *common.D2DWindow
	var tray *common.TrayIcon
	var msgWin windows.HWND = 0
	var job windows.Handle = 0
	hotOK, hotMods, hotVK := false, uint32(0), uint32(0)

	// Stage 1: empty Go (before any window/shell resource is created).
	fmt.Fprintln(os.Stderr, "TRACE stage1 begin")
	res.Stages = append(res.Stages, sampleStage("1-empty-go", "pure-go", nil, nil))
	fmt.Fprintln(os.Stderr, "TRACE stage1 end")
	prev := &res.Stages[0]

	// Stage 4: layered window + D2D + DWrite.
	st4 := sampleStage("4-d2d-dwrite-layered-window", "pure-go", prev, func() error {
		w, err := common.CreateD2DWindow("Wisp S0 baseline (stage 4)")
		if err != nil {
			return err
		}
		win = w
		// render a few more frames to reach steady state
		for i := 0; i < 10; i++ {
			win.RenderFrame(fmt.Sprintf("Wisp S0 baseline frame %d", i))
			common.SleepMs(50)
		}
		common.SleepMs(700) // let DWM composite
		return nil
	})
	res.Stages = append(res.Stages, st4)
	prev = &res.Stages[1]

	// Stage 5: tray + global hotkey + Job Object.
	st5 := sampleStage("5-tray-hotkey-jobobject", "pure-go", prev, func() error {
		hw, err := common.HiddenMessageWindow("wisp-spike-shell")
		if err != nil {
			return err
		}
		msgWin = hw
		t, err := common.AddTrayIcon(hw)
		if err != nil {
			return fmt.Errorf("tray: %w", err)
		}
		tray = t
		hotOK, hotMods, hotVK = common.RegisterGlobalHotKey(hw)
		j, err := common.AttachToKillOnCloseJob()
		if err != nil {
			return fmt.Errorf("job: %w", err)
		}
		job = j
		common.SleepMs(500) // tray icon settles
		return nil
	})
	res.Stages = append(res.Stages, st5)

	// Stage 5b: sanity - render more frames on the same window (leak watch).
	res.Stages = append(res.Stages, sampleStage("5b-render-30-more-frames", "pure-go", prev, func() error {
		if win != nil {
			for i := 0; i < 30; i++ {
				win.RenderFrame(fmt.Sprintf("steady frame %d", i))
				common.SleepMs(33)
			}
		}
		return nil
	}))

	// Teardown (everything must be destroyed - no resident leftovers).
	if win != nil {
		win.TearDown()
	}
	if tray != nil {
		tray.Remove()
	}
	if msgWin != 0 {
		common.UnregisterGlobalHotKey(msgWin)
	}
	// job handle intentionally kept until process exit (KILL_ON_JOB_CLOSE)
	_ = job
	_ = hotOK
	_ = hotMods
	_ = hotVK

	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
	if *out != "" {
		os.WriteFile(*out, b, 0644)
	}
}
