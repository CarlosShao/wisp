// xy-verdict produces the decision inputs for the D25 X/Y topology verdict:
//
//	-role idle-y      : cgo binary (sherpa DLLs mapped at start) + the full
//	                    shell stack (layered D2D/DWrite window + tray + global
//	                    hotkey + Job Object), no session created. This IS the
//	                    path-Y idle process shape; compare its private working
//	                    set against the pinned 25 MB rule (D25).
//	-role idle-x      : same shell stack WITHOUT any cgo import - the path-X
//	                    main-process idle shape, compared against 40 MB
//	                    (D32 16.3.2). Requires the pure-Go build (no tag).
//	-role unload-test : load the ASR paraformer int8 session, run one warmup
//	                    decode, dispose, debug.FreeOSMemory, then sample every
//	                    250 ms for -budget-ms and report the time until private
//	                    WS returns below (pre + 2 MB). D32 16.3.5 rule: if
//	                    path X cannot settle within 10 s, verdict is forced Y.
//
// Build: the speech-session code is behind the `cgo_sherpa` tag so the same
// source compiles both flavors (run.ps1 compiles twice).
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

type unloadReport struct {
	Role            string             `json:"role"`
	Which           string             `json:"which"`
	Machine         common.MachineInfo `json:"machine"`
	StartedAt       string             `json:"startedAtUtc"`
	LoadMs          int64              `json:"loadMs"`
	WarmupMs        int64              `json:"warmupMs"`
	PreWSMB         float64            `json:"prePrivateWS_MB"`
	LoadedWSMB      float64            `json:"loadedPrivateWS_MB"`
	PostDisposeWSMB float64            `json:"postDisposePrivateWS_MB"`
	SettleBudgetMs  int64              `json:"settleBudgetMs"`
	SettledMs       *int64             `json:"settledMs"`
	Settled         bool               `json:"settled"`
	Samples         []float64          `json:"samplesPrivateWS_MB"`
	Note            string             `json:"note"`
	Error           string             `json:"error,omitempty"`
}

type idleReport struct {
	Role            string             `json:"role"`
	Kind            string             `json:"kind"`
	Machine         common.MachineInfo `json:"machine"`
	StartedAt       string             `json:"startedAtUtc"`
	PrivateWSMB     float64            `json:"privateWS_MB"`
	SharedWSMB      float64            `json:"sharedWS_MB"`
	TotalWSMB       float64            `json:"totalWS_MB"`
	PrivateCommitMB float64            `json:"privateCommit_MB"`
	GDIObjects      uint32             `json:"gdiObjects"`
	UserObjects     uint32             `json:"userObjects"`
	Handles         uint32             `json:"handles"`
	Threads         uint32             `json:"threads"`
	ShellParts      string             `json:"shellParts"`
	Note            string             `json:"note"`
	Error           string             `json:"error,omitempty"`
}

type sessionHandle struct {
	loadMs   int64
	warmupMs int64
	dispose  func()
}

func mbOf(b uint64) float64 { return float64(b) / (1 << 20) }

func buildShell() (*common.D2DWindow, *common.TrayIcon, windows.HWND, windows.Handle, error) {
	w, err := common.CreateD2DWindow("Wisp X/Y verdict shell")
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("d2d window: %w", err)
	}
	for i := 0; i < 10; i++ {
		w.RenderFrame(fmt.Sprintf("verdict frame %d", i))
		common.SleepMs(50)
	}
	hw, err := common.HiddenMessageWindow("wisp-spike-verdict")
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("msg window: %w", err)
	}
	tray, err := common.AddTrayIcon(hw)
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("tray: %w", err)
	}
	ok, _, _ := common.RegisterGlobalHotKey(hw)
	_ = ok
	job, err := common.AttachToKillOnCloseJob()
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("job: %w", err)
	}
	common.SleepMs(500)
	return w, tray, hw, job, nil
}

func runIdle(kind string, outPath string) {
	rep := idleReport{
		Role: "idle-" + kind, Machine: common.GetMachineInfo(),
		StartedAt:  time.Now().UTC().Format(time.RFC3339),
		ShellParts: "layered window + Direct2D + DirectWrite + tray + global hotkey + job object",
	}
	switch kind {
	case "y":
		rep.Kind = "cgo"
	case "x":
		rep.Kind = "pure-go"
	}
	w, tray, hw, job, err := buildShell()
	if err != nil {
		rep.Error = err.Error()
	} else {
		common.SettleGC()
		_, med, _ := common.SampleStable(7, 200)
		rep.PrivateWSMB = mbOf(med.PrivateWorkingSet)
		rep.SharedWSMB = mbOf(med.SharedWorkingSet)
		rep.TotalWSMB = mbOf(med.WorkingSet)
		rep.PrivateCommitMB = mbOf(med.PrivateCommit)
		rep.GDIObjects = med.GDIObjects
		rep.UserObjects = med.UserObjects
		rep.Handles = med.Handles
		rep.Threads = med.Threads
		w.TearDown()
		tray.Remove()
		common.UnregisterGlobalHotKey(hw)
		_ = job
	}
	b, _ := json.MarshalIndent(rep, "", "  ")
	fmt.Println(string(b))
	if outPath != "" {
		os.WriteFile(outPath, b, 0644)
	}
}

func runUnload(which, modelsDir string, budgetMs int64, outPath string) {
	rep := unloadReport{
		Role: "unload-test", Which: which, Machine: common.GetMachineInfo(),
		StartedAt: time.Now().UTC().Format(time.RFC3339), SettleBudgetMs: budgetMs,
		Note: "settled = private WS back below pre+2MB; 250ms sampler; dispose includes debug.FreeOSMemory (C11 step 4)",
	}
	defer func() {
		b, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(b))
		if outPath != "" {
			os.WriteFile(outPath, b, 0644)
		}
	}()

	common.SettleGC()
	_, med, _ := common.SampleStable(5, 200)
	pre := med.PrivateWorkingSet
	rep.PreWSMB = mbOf(pre)

	session, err := openSession(which, modelsDir) // nil-safe on pure-go builds
	if err != nil {
		rep.Error = err.Error()
		return
	}
	if session == nil {
		rep.Error = "no session implementation for kind " + which
		return
	}
	rep.LoadMs = session.loadMs
	rep.WarmupMs = session.warmupMs

	common.SettleGC()
	_, m, _ := common.SampleStable(3, 150)
	rep.LoadedWSMB = mbOf(m.PrivateWorkingSet)

	session.dispose()
	common.SettleGC()
	_, m2, _ := common.SampleStable(1, 100)
	rep.PostDisposeWSMB = mbOf(m2.PrivateWorkingSet)

	threshold := pre + 2<<20
	n := int(budgetMs / 250)
	for i := 1; i <= n; i++ {
		common.SleepMs(250)
		s := common.SampleMem()
		rep.Samples = append(rep.Samples, mbOf(s.PrivateWorkingSet))
		if s.PrivateWorkingSet <= threshold {
			t := int64(i) * 250
			rep.SettledMs = &t
			rep.Settled = true
			return
		}
	}
}

func main() {
	role := flag.String("role", "idle-y", "idle-y | idle-x | unload-test")
	which := flag.String("which", "asr", "unload-test session kind: asr")
	modelsDir := flag.String("models", "third_party/spike-models", "spike models dir")
	budget := flag.Int64("budget-ms", 15000, "settle observation window")
	out := flag.String("out", "", "path to write JSON")
	flag.Parse()

	switch *role {
	case "idle-y":
		runIdle("y", *out)
	case "idle-x":
		runIdle("x", *out)
	case "unload-test":
		runUnload(*which, *modelsDir, *budget, *out)
	default:
		fmt.Fprintln(os.Stderr, "unknown role", *role)
		os.Exit(2)
	}
}
