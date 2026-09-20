//go:build windows

package main

// Differential screenshot evidence (ticket 62 gate): the SAME screen region is
// captured twice - once with a live ball window in a given state, once with
// the ball process gone - and the two frames are subtracted. A self-composited
// mock is not evidence; this is. Every state gets:
//
//	NN-<State>-alive.png   composited desktop, ball present
//	NN-<State>-dead.png    same region, ball process terminated
//	NN-<State>-diff.png    |alive-dead| amplified (what the eye must see)
//
// plus one measured row: CPU (all-core and single-core), private bytes,
// working set, handle count and whether an animation timer is armed (read
// from the live ball object, not from the policy table).
//
// The ball runs in a CHILD balldebug process so the parent can terminate it
// and re-shoot the identical region. Termination is graceful - CTRL_BREAK to
// the child's own process group, so the child's teardown destroys the window
// and the tray icon; Kill is only the fallback. runDiff does not return while
// a child it started is still alive.

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// diffOpts configures one differential evidence run.
type diffOpts struct {
	dir       string        // output directory (created if missing)
	states    []string      // states to shoot, in order
	dwell     time.Duration // settle time after the state is applied
	sample    time.Duration // CPU sampling window (ball alive)
	x, y      int           // where to place the ball window (physical px)
	margin    int           // region padding around the window rect
	amplify   int           // diff PNG gain (1 = raw delta)
	size      int           // child orb size px (0 = default from the ball)
	look      string        // liquid treatment handed to the child ("" = default)
	frozen    bool          // measure the frozen SPEC-08 visuals, not the prototype
	level     float64       // synthetic envelope handed to the child (-level)
	killGrace time.Duration // wait for graceful exit before Kill
}

// procStat is one snapshot of the child's resource counters.
type procStat struct {
	privBytes uint64 // private COMMIT (C30 tree metric, PrivateUsage)
	privWS    uint64 // private WORKING SET - the D32 "Memory (private working set)"
	workSet   uint64 // total working set
	gdi       uint32
	user      uint32
	handles   uint32
	cpuNanos  int64
	sampledAt time.Duration
}

// stateRow is the measured evidence line for one state.
type stateRow struct {
	state     string
	px, py    int
	edge      int
	changed3  int
	changed8  int
	changed24 int
	pixels    int
	maxDelta  uint8
	meanDelta float64
	coreShare float64 // % of one logical core
	allShare  float64 // % of all logical cores (the D32 Sleeping basis)
	privMB    float64 // private commit (C30 tree basis)
	privWSMB  float64 // private working set (D32 Sleeping gate basis)
	wsMB      float64
	gdi       uint32
	userObjs  uint32
	handles   uint32
	timers    string
	alive     string
	dead      string
	diff      string
	teardown  string
	box       [4]int
}

var (
	modKernel32X = windows.NewLazySystemDLL("kernel32.dll")
	modPSAPI2    = windows.NewLazySystemDLL("psapi.dll")

	pFindWindowW      = modUser32S.NewProc("FindWindowW") // user32, not kernel32
	pSetCursorPos     = modUser32S.NewProc("SetCursorPos")
	pGetSystemMetrics = modUser32S.NewProc("GetSystemMetrics")
	pGetProcessTimes  = modKernel32X.NewProc("GetProcessTimes")
	pOpenProcessX     = modKernel32X.NewProc("OpenProcess")
	pCloseHandleX     = modKernel32X.NewProc("CloseHandle")
	pGetHandleCountX  = modKernel32X.NewProc("GetProcessHandleCount")
	pGetProcMemInfo   = modPSAPI2.NewProc("GetProcessMemoryInfo")
	pGetGuiResources  = modUser32S.NewProc("GetGuiResources")
	pGenConsoleCtrl   = modKernel32X.NewProc("GenerateConsoleCtrlEvent")
)

const (
	processQueryInformation = 0x0400
	processVMRead           = 0x0010
	processTerminate        = 0x0001
	ctrlBreakEvent          = 1

	// envBallTitle carries the child's private window title, so the parent
	// FindWindows exactly its own child and no foreign Wisp ball.
	envBallTitle = "WISP_BALLDEBUG_TITLE"
)

// utf16Ptr returns nil on error so the syscall degrades instead of crashing.
func utf16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

// ballWindow finds the ball window of one title (class is fixed).
func ballWindow(title string) windows.HWND {
	h, _, _ := pFindWindowW.Call(uintptr(unsafe.Pointer(utf16Ptr("WispBallWindow"))),
		uintptr(unsafe.Pointer(utf16Ptr(title))))
	return windows.HWND(h)
}

// windowRect returns the screen rect of a window.
func windowRect(hwnd windows.HWND) (shotRect, bool) {
	var r shotRect
	ok, _, _ := pGetWindowRectS.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r, ok != 0 && r.r > r.l && r.b > r.t
}

// grab BitBlts a screen region (CAPTUREBLT so layered windows composite in)
// into a top-down RGBA image.
func grab(l, t, w, h int32) (*image.RGBA, error) {
	hdcScreen, _, _ := pGetDCS.Call(0)
	if hdcScreen == 0 {
		return nil, fmt.Errorf("GetDC(NULL) failed")
	}
	defer pReleaseDCS.Call(0, hdcScreen)

	hdcMem, _, _ := pCreateCompatDC.Call(hdcScreen)
	if hdcMem == 0 {
		return nil, fmt.Errorf("CreateCompatibleDC failed")
	}
	defer pDeleteDCS.Call(hdcMem)

	var bmi shotBMI
	bmi.header.size = 40
	bmi.header.width = uint32(w)
	bmi.header.height = uint32(h) // positive = bottom-up DIB
	bmi.header.planes = 1
	bmi.header.bitCount = 32
	var bits uintptr
	hbm, _, err := pCreateDIBSect.Call(hdcMem, uintptr(unsafe.Pointer(&bmi)), 0,
		uintptr(unsafe.Pointer(&bits)), 0, 0)
	if hbm == 0 {
		return nil, fmt.Errorf("CreateDIBSection: %v", err)
	}
	defer pDeleteObjectS.Call(hbm)
	prev, _, _ := pSelectObjectS.Call(hdcMem, hbm)
	defer pSelectObjectS.Call(hdcMem, prev)

	rr, _, _ := pBitBltS.Call(hdcMem, 0, 0, uintptr(w), uintptr(h), hdcScreen,
		uintptr(l), uintptr(t), 0x00CC0020|0x40000000)
	if rr == 0 {
		return nil, fmt.Errorf("BitBlt failed")
	}
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	src := (*[1 << 28]uint8)(winHeapPtrLocal(bits))
	stride := int(w) * 4
	for y := 0; y < int(h); y++ {
		srcRow := (int(h) - 1 - y) * stride // bottom-up source rows
		for x := 0; x < int(w); x++ {
			i := srcRow + x*4
			img.SetRGBA(x, y, color.RGBA{R: src[i+2], G: src[i+1], B: src[i], A: 255})
		}
	}
	return img, nil
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// openProc opens the child for counters and termination.
func openProc(pid uint32) (windows.Handle, error) {
	h, _, err := pOpenProcessX.Call(processQueryInformation|processVMRead|processTerminate, 0, uintptr(pid))
	if h == 0 {
		return 0, fmt.Errorf("OpenProcess(%d): %v", pid, err)
	}
	return windows.Handle(h), nil
}

type filetimePair struct{ lo, hi uint32 }

// sample reads private bytes / working set / handles / cpu time of the child.
func sample(pid uint32, h windows.Handle, since time.Duration) (procStat, error) {
	var st procStat
	var creation, exit, kernel, user filetimePair
	r, _, err := pGetProcessTimes.Call(uintptr(h),
		uintptr(unsafe.Pointer(&creation)), uintptr(unsafe.Pointer(&exit)),
		uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if r == 0 {
		return st, fmt.Errorf("GetProcessTimes: %v", err)
	}
	var pmc struct {
		cb                 uint32
		pageFaultCount     uint32
		peakWorkingSetSize uintptr
		workingSetSize     uintptr
		quotaPeakPagedPool uintptr
		quotaPagedPool     uintptr
		quotaPeakNonPaged  uintptr
		quotaNonPaged      uintptr
		pagefileUsage      uintptr
		peakPagefileUsage  uintptr
		privateUsage       uintptr
	}
	pmc.cb = uint32(unsafe.Sizeof(pmc))
	r, _, err = pGetProcMemInfo.Call(uintptr(h), uintptr(unsafe.Pointer(&pmc)), uintptr(pmc.cb))
	if r == 0 {
		return st, fmt.Errorf("GetProcessMemoryInfo: %v", err)
	}
	var nh uint32
	r, _, err = pGetHandleCountX.Call(uintptr(h), uintptr(unsafe.Pointer(&nh)))
	if r == 0 {
		return st, fmt.Errorf("GetProcessHandleCount: %v", err)
	}
	kernelTicks := int64(kernel.hi)<<32 | int64(kernel.lo)
	userTicks := int64(user.hi)<<32 | int64(user.lo)
	st.cpuNanos = kernelTicks + userTicks
	st.privBytes = uint64(pmc.privateUsage)
	st.workSet = uint64(pmc.workingSetSize)
	st.handles = nh
	st.privWS = privateWorkingSetFor(pid)
	if r, _, _ := pGetGuiResources.Call(uintptr(h), 0); r != 0xFFFFFFFF {
		st.gdi = uint32(r) // GR_GDIOBJECTS
	}
	if r, _, _ := pGetGuiResources.Call(uintptr(h), 1); r != 0xFFFFFFFF {
		st.user = uint32(r) // GR_USEROBJECTS
	}
	st.sampledAt = since
	return st, nil
}

// privateWorkingSetFor returns a process's private working set (bytes) via
// NtQuerySystemInformation(SystemProcessInformation).WorkingSetPrivateSize -
// the kernel field behind Task Manager's "Memory (private working set)"
// column, which is the D32 Sleeping gate basis (same source the ticket 02
// spike used; QueryWorkingSetEx silently writes nothing on this machine).
//
// The scan buffer is a Go slice because the PARENT measures the CHILD: unlike
// the spike's self-measurement, this scratch cannot inflate any number in the
// table. SYSTEM_PROCESS_INFORMATION (x64, stable since Vista): +0x00
// NextEntryOffset, +0x04 NumberOfThreads, +0x08 WorkingSetPrivateSize,
// +0x50 UniqueProcessId.
func privateWorkingSetFor(targetPID uint32) uint64 {
	const (
		systemProcessInformation = 5
		statusInfoLengthMismatch = 0xC0000004
		entryHeaderFloor         = 0xA8
	)
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	pNtQSI := ntdll.NewProc("NtQuerySystemInformation")
	for _, size := range []int{1 << 20, 4 << 20, 16 << 20, 64 << 20} {
		buf := make([]byte, size)
		var retLen uint32
		r1, _, _ := pNtQSI.Call(systemProcessInformation, uintptr(unsafe.Pointer(&buf[0])), uintptr(size),
			uintptr(unsafe.Pointer(&retLen)))
		if r1 != 0 {
			if uint32(r1) == statusInfoLengthMismatch && size < 64<<20 {
				continue
			}
			return 0
		}
		if int(retLen) > size {
			retLen = uint32(size)
		}
		result := uint64(0)
		off := uint64(0)
		for off+entryHeaderFloor <= uint64(retLen) {
			next := uint64(binary.LittleEndian.Uint32(buf[off:]))
			pid := binary.LittleEndian.Uint64(buf[off+0x50:])
			private := int64(binary.LittleEndian.Uint64(buf[off+0x08:]))
			if uint32(pid) == targetPID {
				if private >= 0 {
					result = uint64(private)
				}
				break
			}
			if next == 0 {
				break
			}
			off += next
		}
		return result
	}
	return 0
}

// parkCursor moves the mouse away so the cursor sprite and any hover tooltip
// cannot pollute the diff.
func parkCursor(x, y int) {
	pSetCursorPos.Call(uintptr(x), uintptr(y))
}

// defaultPlacement picks mid-primary-screen: clear wallpaper on the owner's
// 3440x1440 white desktop, away from the tray corner and any IME candidate
// bar, so the alive/dead pair captures uncontested background.
func defaultPlacement() (int, int) {
	cx, _, _ := pGetSystemMetrics.Call(0) // SM_CXSCREEN
	cy, _, _ := pGetSystemMetrics.Call(1) // SM_CYSCREEN
	if cx == 0 || cy == 0 {
		return 600, 400
	}
	return int(int32(cx) / 2), int(int32(cy) / 2)
}

// runDiff is the -diff entry point.
func runDiff(o diffOpts) error {
	if err := os.MkdirAll(o.dir, 0o755); err != nil {
		return fmt.Errorf("balldebug: -diff mkdir: %w", err)
	}
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("balldebug: os.Executable: %w", err)
	}
	if len(o.states) == 0 {
		return fmt.Errorf("balldebug: -diff needs at least one state (-diff-states)")
	}
	if o.x < 0 || o.y < 0 {
		o.x, o.y = defaultPlacement()
	}
	if o.margin == 0 {
		o.margin = 120
	}
	if o.amplify == 0 {
		o.amplify = 6
	}
	if o.killGrace == 0 {
		o.killGrace = 8 * time.Second
	}
	fmt.Printf("diff: place=(%d,%d) margin=%dpx cpu-sample=%s cores=%d amp=%dx out=%s\n",
		o.x, o.y, o.margin, o.sample, runtime.NumCPU(), o.amplify, o.dir)

	rows := make([]stateRow, 0, len(o.states))
	for i, s := range o.states {
		row, err := runDiffForState(exe, o, s, i)
		if err != nil {
			_ = writeDiffTable(filepath.Join(o.dir, "diff-table.txt"), rows)
			return err
		}
		rows = append(rows, row)
		fmt.Printf("diff: %-16s cpu1c=%6.3f%% cpuAll=%6.3f%% privWS=%6.2fMB commit=%6.2fMB ws=%6.2fMB handles=%4d gdi=%d user=%d timers=%-5s "+
			"px>=3:%d/px>=8:%d/px>=24:%d of %d max=%d mean=%.3f box=(%d,%d %dx%d) %s\n",
			row.state, row.coreShare, row.allShare, row.privWSMB, row.privMB, row.wsMB,
			row.handles, row.gdi, row.userObjs, row.timers,
			row.changed3, row.changed8, row.changed24, row.pixels, row.maxDelta, row.meanDelta,
			row.box[0], row.box[1], row.box[2], row.box[3], row.teardown)
	}
	if err := writeDiffTable(filepath.Join(o.dir, "diff-table.txt"), rows); err != nil {
		return err
	}
	fmt.Printf("diff: table -> %s\n", filepath.Join(o.dir, "diff-table.txt"))
	return nil
}

// runDiffForState drives one state: child up, alive grab, counters, graceful
// teardown, dead grab, diff.
func runDiffForState(exe string, o diffOpts, state string, idx int) (stateRow, error) {
	var row stateRow
	row.state = state
	title := fmt.Sprintf("wisp62-%d-%s", os.Getpid(), state)
	statusPath := filepath.Join(o.dir, ".status-"+state)
	_ = os.Remove(statusPath)

	childArgs := []string{"-state", state, "-hold", "-status", statusPath,
		"-x", fmt.Sprint(o.x), "-y", fmt.Sprint(o.y),
		"-cycle-ms", fmt.Sprint(o.dwell.Milliseconds())}
	if o.size > 0 {
		childArgs = append(childArgs, "-size", fmt.Sprint(o.size))
	}
	if o.look != "" {
		childArgs = append(childArgs, "-look", o.look)
	}
	if o.frozen {
		childArgs = append(childArgs, "-frozen")
	}
	if o.level > 0 {
		childArgs = append(childArgs, "-level", fmt.Sprint(o.level))
	}

	cmd := exec.Command(exe, childArgs...)
	cmd.Env = append(os.Environ(), envBallTitle+"="+title)
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000200} // NEW_PROCESS_GROUP
	if err := cmd.Start(); err != nil {
		return row, fmt.Errorf("balldebug: spawn child: %w", err)
	}
	pid := uint32(cmd.Process.Pid)
	// Belt and braces: whatever happens below, the child this run created is
	// gone before we return.
	defer terminateChild(cmd, o.killGrace)

	if !waitUntil(time.Now().Add(30*time.Second), func() bool {
		hwnd := ballWindow(title)
		if hwnd == 0 {
			return false
		}
		r, ok := windowRect(hwnd)
		if !ok {
			return false
		}
		row.px, row.py = int(r.l), int(r.t)
		row.edge = int(r.r - r.l)
		return true
	}) {
		return row, fmt.Errorf("balldebug: child ball window (%q) never appeared", title)
	}
	if exited(cmd) {
		return row, fmt.Errorf("balldebug: child died early")
	}
	// The child writes its status file once the state is applied and rendered.
	if !waitUntil(time.Now().Add(10*time.Second), func() bool { return fileHas(statusPath, "ready=1") }) {
		return row, fmt.Errorf("balldebug: child never reported ready (state=%s)", state)
	}

	parkCursor(4, 4)
	time.Sleep(o.dwell)

	hwnd := ballWindow(title)
	if hwnd == 0 {
		return row, fmt.Errorf("balldebug: window vanished before the alive grab")
	}
	rct, ok := windowRect(hwnd)
	if !ok {
		return row, fmt.Errorf("balldebug: window rect vanished")
	}
	l, t := rct.l-int32(o.margin), rct.t-int32(o.margin)
	w, h := (rct.r-rct.l)+2*int32(o.margin), (rct.b-rct.t)+2*int32(o.margin)
	// The grab-time rect is the truth for this row (the discovery-time rect can
	// predate the -x/-y debug placement).
	row.px, row.py = int(rct.l), int(rct.t)
	row.edge = int(rct.r - rct.l)

	name := fmt.Sprintf("%02d-%s", idx+1, state)
	row.alive, row.dead, row.diff = name+"-alive.png", name+"-dead.png", name+"-diff.png"

	alive, err := grab(l, t, w, h)
	if err != nil {
		return row, fmt.Errorf("balldebug: alive grab: %w", err)
	}

	hproc, err := openProc(pid)
	if err != nil {
		return row, err
	}
	defer pCloseHandleX.Call(uintptr(hproc))

	t0 := time.Now()
	s0, err := sample(pid, hproc, 0)
	if err != nil {
		return row, err
	}
	time.Sleep(o.sample)
	s1, err := sample(pid, hproc, time.Since(t0))
	if err != nil {
		return row, err
	}
	wall := float64(s1.sampledAt.Nanoseconds())
	cpuTicks := float64(s1.cpuNanos - s0.cpuNanos)
	row.coreShare = 100 * cpuTicks / wall
	row.allShare = row.coreShare / float64(runtime.NumCPU())
	row.privMB = float64(s1.privBytes) / (1024 * 1024)
	row.privWSMB = float64(s1.privWS) / (1024 * 1024)
	row.wsMB = float64(s1.workSet) / (1024 * 1024)
	row.gdi, row.userObjs = s1.gdi, s1.user
	row.handles = s1.handles
	if fileHas(statusPath, "timers=true") {
		row.timers = "yes"
	} else if fileHas(statusPath, "timers=false") {
		row.timers = "no"
	} else {
		row.timers = "?"
	}

	if err := writePNG(filepath.Join(o.dir, row.alive), alive); err != nil {
		return row, err
	}

	terminateChild(cmd, o.killGrace)
	row.teardown = exitWord(cmd)

	if !waitUntil(time.Now().Add(6*time.Second), func() bool { return ballWindow(title) == 0 }) {
		return row, fmt.Errorf("balldebug: child window still on screen after teardown")
	}
	time.Sleep(250 * time.Millisecond)
	dead, err := grab(l, t, w, h)
	if err != nil {
		return row, fmt.Errorf("balldebug: dead grab: %w", err)
	}
	if err := writePNG(filepath.Join(o.dir, row.dead), dead); err != nil {
		return row, err
	}
	img, stat := amplifyDiff(alive, dead, o.amplify)
	if err := writePNG(filepath.Join(o.dir, row.diff), img); err != nil {
		return row, err
	}
	row.changed3, row.changed8, row.changed24 = stat.changed3, stat.changed8, stat.changed24
	row.pixels, row.maxDelta, row.meanDelta = stat.pixels, stat.maxD, stat.meanD
	row.box = [4]int{stat.box.Min.X, stat.box.Min.Y, stat.box.Dx(), stat.box.Dy()}
	_ = os.Remove(statusPath)
	return row, nil
}

// terminateChild asks the child to shut itself down (CTRL_BREAK to its own
// process group so its teardown runs) and only force-kills after the grace.
func terminateChild(cmd *exec.Cmd, grace time.Duration) {
	if cmd.Process == nil || exited(cmd) {
		return
	}
	pid := cmd.Process.Pid
	pGenConsoleCtrl.Call(ctrlBreakEvent, uintptr(pid))
	ok := waitUntil(time.Now().Add(grace), func() bool { return exited(cmd) })
	if !ok {
		_ = cmd.Process.Kill()
		waitUntil(time.Now().Add(5*time.Second), func() bool { return exited(cmd) })
	}
	_ = cmd.Wait()
}

// exitWord reports how the child really left (terminateChild already Wait()ed,
// so the status is read from ProcessState, never by a second Wait).
func exitWord(cmd *exec.Cmd) string {
	st := cmd.ProcessState
	if st == nil {
		if err := cmd.Wait(); err != nil {
			return "wait-error"
		}
		st = cmd.ProcessState
	}
	if st == nil {
		return "unknown"
	}
	if st.ExitCode() == 0 {
		return "clean"
	}
	return fmt.Sprintf("exit:%d", st.ExitCode())
}

func exited(cmd *exec.Cmd) bool {
	return cmd.ProcessState != nil && cmd.ProcessState.Exited()
}

func waitUntil(deadline time.Time, cond func() bool) bool {
	for {
		if cond() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(60 * time.Millisecond)
	}
}

func fileHas(path, needle string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(b), needle)
}

// diffStat summarises one alive-vs-dead comparison.
type diffStat struct {
	pixels    int
	changed3  int
	changed8  int
	changed24 int
	maxD      uint8
	meanD     float64
	box       image.Rectangle
}

// amplifyDiff takes the max channel delta per pixel and writes a gain-boosted
// grayscale PNG, so a 2/255 change is still visible to the eye.
func amplifyDiff(alive, dead *image.RGBA, gain int) (*image.Gray, diffStat) {
	b := alive.Bounds()
	st := diffStat{box: image.Rectangle{
		Min: image.Point{X: b.Dx(), Y: b.Dy()},
		Max: image.Point{X: -1, Y: -1},
	}}
	out := image.NewGray(b)
	var sum float64
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			a := alive.RGBAAt(x, y)
			d := dead.RGBAAt(x, y)
			delta := absDiff(a.R, d.R)
			if v := absDiff(a.G, d.G); v > delta {
				delta = v
			}
			if v := absDiff(a.B, d.B); v > delta {
				delta = v
			}
			st.pixels++
			sum += float64(delta)
			if delta >= 3 {
				st.changed3++
			}
			if delta >= 8 {
				st.changed8++
			}
			if delta >= 24 {
				st.changed24++
				if x < st.box.Min.X {
					st.box.Min.X = x
				}
				if y < st.box.Min.Y {
					st.box.Min.Y = y
				}
				if x > st.box.Max.X {
					st.box.Max.X = x
				}
				if y > st.box.Max.Y {
					st.box.Max.Y = y
				}
			}
			if delta > st.maxD {
				st.maxD = delta
			}
			v := uint32(delta) * uint32(gain)
			if v > 255 {
				v = 255
			}
			out.SetGray(x, y, color.Gray{Y: uint8(v)})
		}
	}
	if st.box.Max.X < st.box.Min.X {
		st.box = image.Rectangle{}
	} else {
		st.box.Max.X++
		st.box.Max.Y++
	}
	st.meanD = sum / float64(maxInt(st.pixels, 1))
	return out, st
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var diffCols = []string{
	"state", "alive", "dead", "diff", "x", "y", "edge", "timers",
	"cpu_pct_1core", "cpu_pct_allcores", "privateWorkingSet_MB", "privateCommit_MB",
	"workset_MB", "handles", "gdiObjects", "userObjects",
	"px_delta_ge3", "px_delta_ge8", "px_delta_ge24", "px_total", "max_delta",
	"mean_delta", "bbox_x", "bbox_y", "bbox_w", "bbox_h", "teardown",
}

func writeDiffTable(path string, rows []stateRow) error {
	var sb strings.Builder
	sb.WriteString(strings.Join(diffCols, "\t") + "\n")
	for _, r := range rows {
		cells := []string{
			r.state, r.alive, r.dead, r.diff, fmt.Sprint(r.px), fmt.Sprint(r.py),
			fmt.Sprint(r.edge), r.timers,
			fmt.Sprintf("%.3f", r.coreShare), fmt.Sprintf("%.3f", r.allShare),
			fmt.Sprintf("%.2f", r.privWSMB), fmt.Sprintf("%.2f", r.privMB),
			fmt.Sprintf("%.2f", r.wsMB), fmt.Sprint(r.handles),
			fmt.Sprint(r.gdi), fmt.Sprint(r.userObjs),
			fmt.Sprint(r.changed3), fmt.Sprint(r.changed8), fmt.Sprint(r.changed24),
			fmt.Sprint(r.pixels), fmt.Sprint(r.maxDelta), fmt.Sprintf("%.3f", r.meanDelta),
			fmt.Sprint(r.box[0]), fmt.Sprint(r.box[1]), fmt.Sprint(r.box[2]), fmt.Sprint(r.box[3]),
			r.teardown,
		}
		sb.WriteString(strings.Join(cells, "\t") + "\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
