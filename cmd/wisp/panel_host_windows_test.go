//go:build windows

package main

// Local, real-host tests for ticket 33 AC#1..AC#4. These create an actual
// WebView2 window, so they belong to the winlive tier, which has NO CI job
// (recorded in the evidence table): they run only on this desktop session and
// every reading is stamped with the machine time + HEAD. They are gated by the
// //go:build windows tag (a real platform fork, not a fake skip); they never use
// testing.Short() or an environment escape hatch to pretend green, and a window
// that cannot be created here is a red, not a skip.

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	"golang.org/x/sys/windows"
)

// hostThreadHarness runs the host on a private locked OS thread for the local
// measurement / round-trip tests. It is deliberately a TEST harness - ban #1
// targets production goroutines, and production bringUp runs on the caller's
// ui-sta thread, never here. The goroutine owns itself, recovers anything the
// pump raises, and stops cleanly when stopHostThread is called.
type hostThreadHarness struct {
	stopping atomic.Bool
}

func (hh *hostThreadHarness) bringUp(mgr *PanelManager, ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		// owner: test panel-host thread; recover below.
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panel host pump panicked: %v", r)
			}
		}()
		runtime.LockOSThread()
		errCh <- mgr.bringUp(ctx)
		for !hh.stopping.Load() {
			pnlPumpOnce()
			time.Sleep(5 * time.Millisecond)
		}
		runtime.UnlockOSThread()
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return context.DeadlineExceeded
	}
}

func (hh *hostThreadHarness) stopHostThread() { hh.stopping.Store(true) }

// recordingModeHandler is the test's stand-in for the production mode writer: it
// records that the router reached a handler for a real postMessage, and accepts.
type recordingModeHandler struct {
	mu   sync.Mutex
	reqs []panel.ComposerRequest
}

func (h *recordingModeHandler) HandleModeRequest(_ context.Context, req panel.ComposerRequest) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.reqs = append(h.reqs, req)
	return nil
}

func (h *recordingModeHandler) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.reqs)
}

// ---------------------------------------------------------------------------
// AC#1 / AC#3 rulers, re-tooled by 33-r4.
//
// Why both live in this shape now: the acceptance leg (docs/evidence/s1/
// 33-panel-host-c27-v1.md §A#23-#26) showed the pre-33-r4 versions could not go
// red - the socket ruler asked GetExtendedTcpTable for the wrong table class,
// filtered on the wrong state constant, only looked at IPv4, and swallowed a
// failed read as "zero listeners"; the child-count ruler counted every
// msedgewebview2.exe on the machine instead of the ones this process started,
// so other people's browser sessions were in AC#1's denominator. Neither claim
// is taken on faith here: TestAC3ListeningSocketRulerSeesItsOwnListener below is
// a positive control that opens a real loopback listener and requires the ruler
// to count it.
//
// ⛔ Nothing in panel_host_windows.go was touched by this leg - the production
// half of ticket 33 belongs to 33-r2.
// ---------------------------------------------------------------------------

// webviewExeName is the browser process a WebView2 window spawns.
const webviewExeName = "msedgewebview2.exe"

// procRow is one line of the CreateToolhelp32Snapshot process table.
type procRow struct {
	pid  uint32
	ppid uint32
	name string
}

// treeReading is one process-tree sample. Two denominators are carried on
// purpose: TreeWebview is what the assertions use ("OUR tree"), MachineNamed is
// only a reading (what the pre-33-r4 ruler counted, which 33-v1 §A#26 measured at
// 14 rows belonging to other sessions on this box).
type treeReading struct {
	Root         uint32
	TreePIDs     []uint32
	DirectKids   int
	TreeWebview  int
	MachineNamed int
}

// snapshotProcesses reads the machine's process table once. An unreadable table
// is a FAILED MEASUREMENT and stops the test: it must never be reported as
// "nothing is running" (that silent 0 is the shape 33-v1 §A#24 named).
func snapshotProcesses(t *testing.T) []procRow {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		t.Fatalf("process snapshot: %v (the measurement failed - this is not a panel verdict)", err)
	}
	defer windows.CloseHandle(snap)
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	rows := make([]procRow, 0, 256)
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		rows = append(rows, procRow{pid: pe.ProcessID, ppid: pe.ParentProcessID, name: windows.UTF16ToString(pe.ExeFile[:])})
	}
	if len(rows) == 0 {
		t.Fatalf("process snapshot walked 0 rows - a broken ruler is not a green AC#1")
	}
	return rows
}

// readWebviewTree samples the process tree rooted at root (root itself excluded
// from TreePIDs, WebView2 children included at any depth - the browser spawns a
// grandchild set, so a direct-children-only walk would read "no children" while
// six processes are alive).
func readWebviewTree(t *testing.T, root uint32) treeReading {
	t.Helper()
	rows := snapshotProcesses(t)
	byParent := make(map[uint32][]procRow, len(rows))
	r := treeReading{Root: root}
	for _, row := range rows {
		byParent[row.ppid] = append(byParent[row.ppid], row)
		if strings.EqualFold(row.name, webviewExeName) {
			r.MachineNamed++
		}
	}
	seen := map[uint32]bool{root: true}
	queue := []uint32{root}
	for len(queue) > 0 && len(r.TreePIDs) < 4096 {
		cur := queue[0]
		queue = queue[1:]
		kids := byParent[cur]
		if cur == root {
			r.DirectKids = len(kids)
		}
		for _, k := range kids {
			if seen[k.pid] {
				continue
			}
			seen[k.pid] = true
			r.TreePIDs = append(r.TreePIDs, k.pid)
			if strings.EqualFold(k.name, webviewExeName) {
				r.TreeWebview++
			}
			queue = append(queue, k.pid)
		}
	}
	return r
}

// pidSetOfTree returns root plus every pid under it - the scope AC#3 has to be
// read over, because a listening socket opened by a browser child is still "our
// panel host listening" (33-v1 §A#25 asked for the tree, not just os.Getpid()).
func pidSetOfTree(root uint32, tree []uint32) map[uint32]bool {
	set := make(map[uint32]bool, len(tree)+1)
	set[root] = true
	for _, p := range tree {
		set[p] = true
	}
	return set
}

const (
	// tcpTableOwnerPIDAll is TCP_TABLE_OWNER_PID_ALL (5). The pre-33-r4 file
	// carried 4, which is TCP_TABLE_OWNER_PID_CONNECTIONS - a different row
	// layout, so state and pid were read out of the wrong bytes. The sampler
	// already in this repo spells it correctly (internal/proc/
	// treemetrics_windows.go:57 = 5), which is the in-repo cross-check 33-v1
	// §A#24 used.
	tcpTableOwnerPIDAll = 5
	// mibTcpStateListen is MIB_TCP_STATE_LISTEN (2). The pre-33-r4 file carried
	// 10 = MIB_TCP_STATE_ESTAB, so no listening socket could ever have matched.
	mibTcpStateListen = 2
	// errInsufficientBuffer is Win32 ERROR_INSUFFICIENT_BUFFER, what
	// GetExtendedTcpTable returns while rewriting the needed size.
	errInsufficientBuffer = 122
)

// tcpRowLayout gives (row size, state offset, owning-pid offset) for one address
// family's *_OWNER_PID row. Both rows were MEASURED on this box by 33-r4 (a
// repo-external probe that opens a 127.0.0.1 and a [::1] listener and searches
// the layout that attributes them to its own pid) - the IPv6 row size is NOT the
// 48 the header documents imply:
//
//	AF_INET : 24 bytes, state at 0, pid at 20  -> own listener attributed (1 row)
//	AF_INET6: 56 bytes, state at 48, pid at 52 -> own listener attributed (1 row)
//	AF_INET6 at 48/40/44 (the textbook layout)  -> 0 rows attributed, i.e. blind
//
// 33-v1 §B#14 hit the same wall and stopped; the reading that unblocks it is here.
func tcpRowLayout(family uint32) (rowSize, stateOff, pidOff uintptr, err error) {
	switch family {
	case windows.AF_INET:
		// MIB_TCPROW_OWNER_PID: state, localAddr, localPort, remoteAddr,
		// remotePort, owningPid -> 6 x uint32 = 24 bytes.
		return 24, 0, 20, nil
	case windows.AF_INET6:
		// MIB_TCP6ROW_OWNER_PID as THIS table actually lays it out on this host:
		// 56 bytes per row, state at 48, owning pid at 52 (measured, not read off
		// a header). The AF_INET6 count is reported only - asserting on it is the
		// orchestrator's scope call (33-r4 dispatch: do not widen the range alone,
		// and internal/proc's own sampler only reads AF_INET).
		return 56, 48, 52, nil
	default:
		return 0, 0, 0, fmt.Errorf("tcpRowLayout: unsupported address family %d", family)
	}
}

// listenRowsForPIDs counts LISTEN rows owned by any pid in pids for one address
// family. It returns (matches, tableRows, err). err is non-nil when the table
// could not be read at all, and the caller must fail the test on it: "the tool
// did not answer" and "nothing is listening" are different facts and the
// pre-33-r4 ruler made them look identical (33-v1 §A#24).
func listenRowsForPIDs(family uint32, pids map[uint32]bool) (int, int, error) {
	rowSize, stateOff, pidOff, err := tcpRowLayout(family)
	if err != nil {
		return 0, 0, err
	}
	modiphlp := windows.NewLazySystemDLL("iphlpapi.dll")
	procTable := modiphlp.NewProc("GetExtendedTcpTable")
	size := uint32(64 << 10)
	var last uintptr
	for attempt := 0; attempt < 6; attempt++ {
		buf := make([]byte, size)
		want := size
		r1, _, _ := procTable.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&want)),
			1, // bOrder TRUE
			uintptr(family),
			tcpTableOwnerPIDAll,
			0,
		)
		last = r1
		if r1 != 0 {
			if r1 == uintptr(errInsufficientBuffer) && want > size {
				size = want
			}
			continue
		}
		entries := *(*uint32)(unsafe.Pointer(&buf[0]))
		need := uint64(4) + uint64(entries)*uint64(rowSize)
		if need > uint64(want) || want > uint32(len(buf)) {
			return 0, 0, fmt.Errorf("GetExtendedTcpTable(family=%d) reported %d rows but wrote %d of %d buffer bytes - refusing to read a truncated table as 'no listeners'", family, entries, want, len(buf))
		}
		n := 0
		for i := uint32(0); i < entries; i++ {
			base := unsafe.Add(unsafe.Pointer(&buf[0]), 4+uintptr(i)*rowSize)
			state := *(*uint32)(unsafe.Add(base, stateOff))
			own := *(*uint32)(unsafe.Add(base, pidOff))
			if state == mibTcpStateListen && pids[own] {
				n++
			}
		}
		return n, int(entries), nil
	}
	return 0, 0, fmt.Errorf("GetExtendedTcpTable(family=%d) failed 6 times (last code %d, buffer %d bytes): listening count UNVERIFIED", family, last, size)
}

// panelHostLatency records every cold/hot pair this process measures, so a
// -count=N run produces one aggregate reading instead of N unrelated log lines.
// The pre-33-r4 file could only ever print a single pair and hard-coded a commit
// into the text of its own measurement line (33-v1 §A#21).

// latencyRecord accumulates the per-process cold/hot pairs.
type latencyRecord struct {
	mu     sync.Mutex
	cold   []float64
	hot    []float64
	stamps []string
	head   string
}

var panelHostLatency latencyRecord

func (l *latencyRecord) record(head string, coldMs, hotMs float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.head = head
	l.cold = append(l.cold, coldMs)
	l.hot = append(l.hot, hotMs)
	l.stamps = append(l.stamps, time.Now().Format(time.RFC3339))
}

// nearestRankPercentile is the nearest-rank definition (index ceil(q/100*n),
// 1-based, on a sorted copy). With n=10 that makes P95 the second-largest
// sample, so one tail run cannot hide behind the median.
func nearestRankPercentile(vals []float64, q int) float64 {
	if len(vals) == 0 {
		return -1
	}
	cp := append([]float64(nil), vals...)
	sort.Float64s(cp)
	idx := int(math.Ceil(float64(q) / 100.0 * float64(len(cp))))
	if idx < 1 {
		idx = 1
	}
	if idx > len(cp) {
		idx = len(cp)
	}
	return cp[idx-1]
}

// TestPanelHostRealWindowHopAndLifecycle is the winlive measurement for ticket
// 33 AC#1 (lifecycle / single-window reuse / process-tree child set), AC#2 (cold
// and hot latency on a real window), AC#3 (no listening socket) and AC#4 (focus
// round trip). It runs on a private locked OS thread (hostThreadHarness) and only
// on this desktop: winlive has no CI job.
//
// Do not parallelize this test or any sibling in this file: the AC#3 socket ruler
// and the AC#4 foreground ruler both read machine-global state, and
// TestAC3ListeningSocketRulerSeesItsOwnListener deliberately opens a listener.
func TestPanelHostRealWindowHopAndLifecycle(t *testing.T) {
	// The WebView2 user-data folder is created here and deliberately NOT deleted
	// (repo rule: throwaways are only ever created). t.TempDir's auto-RemoveAll
	// races the browser's own profile-DB teardown and fails on a locked file - that
	// is the teardown, not a product assertion. The AC#1 exit deadline below is the
	// thing that actually bounds the child processes.
	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
	h := &recordingModeHandler{}
	disp := &panel.ComposerDispatch{Mode: h}
	mgr := NewPanelManager(disp, nil, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	self := uint32(os.Getpid())
	baselineTree := readWebviewTree(t, self)

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("bringUp on the test thread: %v (WebView2 window creation is available on this box - the S0 spike coldonly returned coldMs=862.431 at %s; a nil here means the host, not the environment)", err, time.Now().Format(time.RFC3339))
	}
	defer hh.stopHostThread()
	defer mgr.Destroy()

	if !mgr.IsCreated() {
		t.Fatalf("host reports not created after bring-up")
	}
	if hwnd := mgr.windowHandle(); hwnd == 0 {
		t.Fatalf("host created but HWND is 0")
	}

	// AC#2 cold: the bring-up stamp is a real measurement on this machine. The
	// commit is read at run time - the pre-33-r4 line printed the literal text
	// "HEAD 7a41db9b", so eleven runs taken on a different commit all claimed to
	// have been taken on that one (33-v1 §A#21).
	head := gitHeadShortForTest(t)
	coldMs := mgr.LastColdMs()
	t.Logf("cold bring-up measured on this box: %.3f ms (HEAD %s at read time, %s)", coldMs, head, time.Now().Format(time.RFC3339))
	if coldMs <= 0 {
		t.Fatalf("cold bring-up did not produce a browser round trip (got %.3f); the message channel did not come up on the real window", coldMs)
	}
	if coldMs > 1500 {
		t.Errorf("cold bring-up %.1f ms exceeds D32 panel cold budget 1500 ms", coldMs)
	}

	// H3 + H10: a real composer envelope reaches the router through the ONE door.
	env := `{"method":"panel.mode.request","requestId":"rt-33r1","source":"panel-composer","to":"ask_every_step"}`
	reply, err := mgr.dispatchRaw(ctx, env)
	if err != nil {
		t.Fatalf("dispatchRaw on a valid envelope: %v (reply=%q)", err, reply)
	}
	if h.count() != 1 {
		t.Fatalf("the page-facing door did not reach the mode handler: calls=%d, want 1 (H3 landed on a real host)", h.count())
	}

	// AC#1 hide-don't-destroy, nailed by IDENTITY: a hide then re-show must hand
	// back the SAME HWND, and the browser child set INSIDE THIS PROCESS TREE must
	// not change. 33-v1 §A#26 judged that neither half was actually checked - the
	// handles were only ever compared against 0, and the child count was a
	// machine-wide count by executable name.
	hwndBeforeHide := mgr.windowHandle()
	before := readWebviewTree(t, self)
	mgr.Hide()
	if mgr.IsShown() {
		t.Fatalf("IsShown true after Hide")
	}
	if !mgr.IsCreated() {
		t.Fatalf("Hide destroyed the window: hide-don't-destroy is broken")
	}
	if err := mgr.HotShow(ctx); err != nil {
		t.Fatalf("HotShow: %v", err)
	}
	hwndAfterReShow := mgr.windowHandle()
	if hwndAfterReShow != hwndBeforeHide {
		t.Errorf("hide -> re-show did not reuse the same window: HWND 0x%x became 0x%x (AC#1 words 'second show within session reuses window'; a fresh HWND is a second browser, not a re-show)", hwndBeforeHide, hwndAfterReShow)
	}
	after := readWebviewTree(t, self)
	if before.TreeWebview != after.TreeWebview {
		t.Errorf("hide -> re-show changed the msedgewebview2 set inside THIS process tree: %d -> %d (tree root pid %d; single-window reuse broken). Machine-wide count went %d -> %d and is reported only - it counts other sessions too",
			before.TreeWebview, after.TreeWebview, self, before.MachineNamed, after.MachineNamed)
	}
	if !mgr.IsShown() {
		t.Fatalf("IsShown false after re-show")
	}
	t.Logf("AC#1 denominators (head %s): our tree webview=%d direct-children=%d tree-pids=%d | machine-wide msedgewebview2=%d | same HWND 0x%x across hide->re-show=%t",
		head, after.TreeWebview, after.DirectKids, len(after.TreePIDs), after.MachineNamed, hwndAfterReShow, hwndAfterReShow == hwndBeforeHide)

	// AC#2 hot: re-show latency on the reused window.
	hotMs := mgr.LastHotMs()
	t.Logf("hot re-show measured on this box: %.3f ms", hotMs)
	if hotMs <= 0 {
		t.Errorf("hot re-show did not produce a latency reading")
	} else if hotMs > 200 {
		t.Errorf("hot re-show %.1f ms exceeds D32 panel hot budget 200 ms", hotMs)
	}
	panelHostLatency.record(head, coldMs, hotMs)

	// AC#3: neither this process nor anything it started holds a LISTEN socket.
	// The ruler now reads the right table (TCP_TABLE_OWNER_PID_ALL) and the right
	// state (MIB_TCP_STATE_LISTEN), fails the test when the table cannot be read
	// instead of reporting 0, and covers the process tree rather than only
	// os.Getpid() (33-v1 §A#24/#25).
	pids := pidSetOfTree(self, after.TreePIDs)
	v4, rows4, err4 := listenRowsForPIDs(windows.AF_INET, pids)
	if err4 != nil {
		t.Fatalf("AC#3 could not read the IPv4 table: %v", err4)
	}
	if v4 != 0 {
		t.Errorf("the panel host opened %d IPv4 LISTEN socket(s) across pid set %d (self %d + %d tree pids); D29 forbids a localhost server", v4, len(pids), self, len(after.TreePIDs))
	}
	// The IPv6 family is MEASURED AND REPORTED here, not asserted: the scope
	// decision belongs to the orchestrator (33-r4 dispatch: "don't widen the
	// range on your own"), and internal/proc's own sampler only reads AF_INET
	// (treemetrics_windows.go:277). TestAC3ListeningSocketRulerSeesItsOwnListener
	// below is what tells this reading whether the IPv6 row layout is trustworthy.
	v6, rows6, err6 := listenRowsForPIDs(windows.AF_INET6, pids)
	t.Logf("AC#3 listening-socket ruler (head %s): IPv4 rows=%d LISTEN-owned-by-tree=%d ; IPv6 rows=%d LISTEN-owned-by-tree=%d err=%v",
		head, rows4, v4, rows6, v6, err6)

	// AC#4's focus round trip is NOT asserted here: it moved to
	// TestAC4FocusReturnToPriorWindowGap33r2 below, which owns its own window so
	// that its (currently expected-red) assertions cannot blur into AC#1/AC#2/AC#3
	// readings. What stays here is the pre-33-r4 statement that was already real:
	// the panel is created once and reused.

	// Destroy for good leaves no window (recreate path starts clean), and AC#1's
	// second clause: the WebView children THIS process started exit within 2s.
	// This is a bounded wait on a monotonic deadline, not a wall-clock timeout.
	mgr.Destroy()
	if mgr.IsCreated() {
		t.Fatalf("IsCreated still true after Destroy")
	}
	if hwnd := mgr.windowHandle(); hwnd != 0 {
		t.Fatalf("HWND nonzero after Destroy")
	}
	deadline := time.Now().Add(2 * time.Second)
	finalTree := readWebviewTree(t, self)
	for finalTree.TreeWebview > baselineTree.TreeWebview && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		finalTree = readWebviewTree(t, self)
	}
	if finalTree.TreeWebview > baselineTree.TreeWebview {
		t.Errorf("after Destroy the WebView2 children this process started did not exit within 2s: our tree went baseline %d -> now %d (tree pids %d). The machine-wide count (%d -> %d) is reported only and is NOT the denominator",
			baselineTree.TreeWebview, finalTree.TreeWebview, len(finalTree.TreePIDs), baselineTree.MachineNamed, finalTree.MachineNamed)
	}
	t.Logf("AC#1 after Destroy (head %s): our tree webview baseline %d -> now %d, machine-wide %d -> %d",
		head, baselineTree.TreeWebview, finalTree.TreeWebview, baselineTree.MachineNamed, finalTree.MachineNamed)
}

// TestAC4FocusReturnToPriorWindowGap33r2 is ticket 33 AC#4's ruler ("editor ->
// panel -> hide -> focus back in editor"), re-tooled by 33-r4. The name carries
// the leg that owns the fix because THIS TEST IS EXPECTED TO BE RED until the
// product's focus hop is right: 33-v1 §A#27 judged the version it replaces a ruler
// that could not go red under any focus behaviour at all (the return hop was only
// ever t.Logf'd, and `prior` was sampled AFTER HotShow, so in all eleven of 33-v1's
// runs `prior` was the panel's own HWND - the panel handing focus back to the
// window it had just hidden).
//
// Three statements, in the order the ticket words the hop:
//  1. instrument self-check - the foreground window this process sees BEFORE it
//     has a panel must not be the panel's own handle. Without this, statements 2
//     and 3 could be reading a stuck value; a failure here is a failed
//     measurement (t.Fatalf), not a product verdict;
//  2. what the host RECORDS as the previous foreground when it shows the panel
//     must not be the panel itself (read off the product's own field,
//     panel_host_windows.go:269 records / :311-313 restores);
//  3. after Hide, the foreground must not still be the hidden panel window -
//     D29 says focus goes back to the recorded prior window.
//
// Reading on this box (head 0d717917, 2026-10-01 11:4x-11:5x): every foreground
// sample in the hop - while hidden, after Show, after Hide - came back as the
// panel's own HWND. The panel takes the foreground and nothing ever gives it
// back, so statements 2 and 3 are red for a reason that lives in the production
// file, which 33-r4 is forbidden to touch. ⛔ Do not relax these to green: 33-r4's
// deliverable is "the ruler rings", and this is the one that finally does.
func TestAC4FocusReturnToPriorWindowGap33r2(t *testing.T) {
	foregroundBefore := uintptr(windows.GetForegroundWindow())

	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
	h := &recordingModeHandler{}
	mgr := NewPanelManager(&panel.ComposerDispatch{Mode: h}, nil, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("bringUp for the AC#4 focus hop: %v", err)
	}
	defer hh.stopHostThread()
	defer mgr.Destroy()

	panelHwnd := mgr.windowHandle()
	if panelHwnd == 0 {
		t.Fatalf("host created but HWND is 0 - the focus hop cannot be measured")
	}
	if foregroundBefore == panelHwnd {
		t.Fatalf("the foreground instrument cannot distinguish windows on this box: the foreground recorded BEFORE this process had a panel is already the panel's HWND 0x%x - statements 2 and 3 below would be vacuous, so this is a failed measurement, not a green", panelHwnd)
	}

	// Hide FIRST, then read the prior foreground: this is the sampling-order fix.
	// The pre-33-r4 file read it after HotShow, i.e. after the panel had already
	// taken the foreground.
	mgr.Hide()
	prior := uintptr(windows.GetForegroundWindow())
	if err := mgr.Show(ctx); err != nil {
		t.Fatalf("Show for the focus hop: %v", err)
	}
	afterShow := uintptr(windows.GetForegroundWindow())
	recorded := uintptr(mgr.prevFocus)
	mgr.Hide()
	afterHide := uintptr(windows.GetForegroundWindow())
	head := gitHeadShortForTest(t)
	t.Logf("AC#4 focus hop (head %s): foreground before any panel 0x%x | foreground while hidden (prior) 0x%x | after Show 0x%x | panel hwnd 0x%x | prevFocus recorded at Show 0x%x | after Hide 0x%x",
		head, foregroundBefore, prior, afterShow, panelHwnd, recorded, afterHide)

	if afterShow != panelHwnd {
		t.Errorf("the panel did not take the foreground on Show: foreground 0x%x, panel hwnd 0x%x (AC#4 says only the panel takes focus when shown)", afterShow, panelHwnd)
	}
	if recorded == panelHwnd {
		t.Errorf("the host recorded the panel ITSELF as the previous foreground window (prevFocus 0x%x == the panel) - because this panel was already foreground from the earlier show, the handle of the window the user actually came from is now lost, and Hide hands focus back to the window it just hid. Product side: panel_host_windows.go:269 records unconditionally, :311-313 restores it. RED BY DESIGN, owner 33-r2; ⛔ not to be relaxed in the test", recorded)
	}
	if afterHide == panelHwnd {
		t.Errorf("after Hide the foreground window is STILL the hidden panel (0x%x); the pre-Show foreground was 0x%x. No window got focus back, so D29's 'focus returns to the recorded previous foreground window' is not happening on a real host. RED BY DESIGN, owner 33-r2", afterHide, prior)
	} else if prior != 0 && prior != panelHwnd && afterHide != prior {
		t.Errorf("Hide did not hand focus back to the window that held it before Show: expected 0x%x, foreground is 0x%x", prior, afterHide)
	}
}

// TestPanelHostLatencyPercentilesAC2 aggregates every cold/hot pair measured in
// this process and gates them on the same D32 budget numbers the single-run
// assertions use (cold 1500 ms, hot 200 ms) - at P50 and P95 instead of one
// sample. It exists because ticket 33's AC#2 words "10-run P50/P95 into SLO
// appendix" and the pre-33-r4 instrument could only emit one pair per run
// (33-v1 §A#20: grep for P50/P95 in this package = 0 hits).
//
// It is declared AFTER the lifecycle test on purpose: go test runs the tests of a
// file in source order, and the aggregate is read from the samples that file's
// earlier tests recorded. Run with -count=10 to get a ten-sample appendix:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= \
//	  go test ./cmd/wisp/ -count=10 -v -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'
//
// The appendix in docs/evidence/s1/33-panel-host-c27-r4.md is the deliverable;
// nothing here writes into docs/SLO.md (that file is on the do-not-touch list).
func TestPanelHostLatencyPercentilesAC2(t *testing.T) {
	panelHostLatency.mu.Lock()
	cold := append([]float64(nil), panelHostLatency.cold...)
	hot := append([]float64(nil), panelHostLatency.hot...)
	stamps := append([]string(nil), panelHostLatency.stamps...)
	head := panelHostLatency.head
	panelHostLatency.mu.Unlock()

	if len(cold) == 0 {
		t.Skipf("no cold/hot sample recorded in this process: the lifecycle test did not run in this binary (use -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'). Named skip - an empty aggregate is not a green latency gate")
	}
	t.Logf("AC#2 appendix (head %s, n=%d, same-process consecutive runs):", head, len(cold))
	for i := range cold {
		stamp := "stamp-unavailable"
		if i < len(stamps) {
			stamp = stamps[i]
		}
		t.Logf("  sample %2d: cold=%.3f ms hot=%.3f ms at %s", i+1, cold[i], hot[i], stamp)
	}
	coldP50 := nearestRankPercentile(cold, 50)
	coldP95 := nearestRankPercentile(cold, 95)
	hotP50 := nearestRankPercentile(hot, 50)
	hotP95 := nearestRankPercentile(hot, 95)
	t.Logf("AC#2 percentiles over %d runs (nearest-rank): cold P50=%.3f P95=%.3f (budget 1500) | hot P50=%.3f P95=%.3f (budget 200)",
		len(cold), coldP50, coldP95, hotP50, hotP95)
	// P11's re-evaluation line (cold > 2000 ms) is reported, never a threshold of
	// its own here - docs/SLO.md and internal/observe/thresholds.go are untouched.
	t.Logf("AC#2 P11 line (cold > 2000 ms) - max cold observed over these runs: %.3f ms", maxOf(cold))
	if coldP95 > 1500 {
		t.Errorf("cold P95 %.3f ms over %d runs exceeds the D32 panel cold budget 1500 ms (the single-run assertion is not the only gate: the tail is what AC#2 asks for)", coldP95, len(cold))
	}
	if hotP95 > 200 {
		t.Errorf("hot re-show P95 %.3f ms over %d runs exceeds the D32 panel hot budget 200 ms", hotP95, len(hot))
	}
}

func maxOf(vals []float64) float64 {
	m := -1.0
	for _, v := range vals {
		if v > m {
			m = v
		}
	}
	return m
}

// TestAC3ListeningSocketRulerSeesItsOwnListener is 33-r4's positive control for
// the AC#3 ruler: it opens a REAL loopback listening socket from THIS process and
// requires the ruler to count it, then closes it and requires the count to fall
// back. The ruler this replaces stayed green while the host genuinely held
// 127.0.0.1:62971 (33-v1 §A#23, ground truth from PowerShell in the same run), so
// a control that plants a listener inside the test itself is the only form that
// makes "class 5 / state 2 / row layout" a measured fact rather than a claim.
//
// The test file is allowed to import net: the L1 gate
// (TestPanelHostOpensNoListeningSocketL1) parses the PRODUCTION file only, which
// is exactly what J7 ruled.
//
// ⛔ Do not widen or narrow the assertions here to make them pass. If this goes
// red, the ruler is blind and AC#3 has no instrument.
func TestAC3ListeningSocketRulerSeesItsOwnListener(t *testing.T) {
	self := uint32(os.Getpid())
	pids := pidSetOfTree(self, nil)

	base4, rows4, err := listenRowsForPIDs(windows.AF_INET, pids)
	if err != nil {
		t.Fatalf("AC#3 ruler cannot read the IPv4 table at all: %v", err)
	}
	ln4, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("the positive control could not open a loopback listener on this box: %v (without a real listener this ruler's teeth cannot be measured - that is a failed measurement, not a pass)", err)
	}
	defer func() { _ = ln4.Close() }()

	seen4, _, err := listenRowsForPIDs(windows.AF_INET, pids)
	if err != nil {
		t.Fatalf("AC#3 ruler could not re-read the IPv4 table: %v", err)
	}
	t.Logf("AC#3 positive control: IPv4 table rows=%d, this pid owned %d before / %d while %s is listening", rows4, base4, seen4, ln4.Addr())
	if seen4 < base4+1 {
		t.Errorf("the AC#3 ruler did not see the loopback listener THIS test just opened on %s: counted %d, baseline %d, want at least %d. The ruler is blind again - check the table class (want TCP_TABLE_OWNER_PID_ALL=5), the state constant (want MIB_TCP_STATE_LISTEN=2) and the row layout", ln4.Addr(), seen4, base4, base4+1)
	}

	// The IPv6 family: an INSTRUMENT self-check, not an AC#3 gate. Whether the
	// product assertion should cover IPv6 is the orchestrator's scope call (33-r4
	// dispatch: "don't widen the range on your own"; internal/proc's own sampler
	// reads AF_INET only - treemetrics_windows.go:277). What is nailed here is that
	// the ruler's IPv6 row layout is not blind, because a reported-but-blind number
	// is how this file shipped its last false green (33-v1 §B#14 stopped exactly
	// here).
	base6, _, errb6 := listenRowsForPIDs(windows.AF_INET6, pids)
	if errb6 != nil {
		t.Fatalf("AC#3 ruler cannot read the IPv6 table at all: %v", errb6)
	}
	if ln6, err6 := net.Listen("tcp", "[::1]:0"); err6 == nil {
		defer func() { _ = ln6.Close() }()
		seen6, rows6, errv6 := listenRowsForPIDs(windows.AF_INET6, pids)
		if errv6 != nil {
			t.Fatalf("AC#3 ruler could not re-read the IPv6 table: %v", errv6)
		}
		t.Logf("AC#3 positive control, IPv6 family: listening on %s, table rows=%d, this pid owns %d LISTEN rows (baseline %d)", ln6.Addr(), rows6, seen6, base6)
		if seen6 < base6+1 {
			t.Errorf("the ruler's IPv6 layout is blind: it did not attribute the [::1] listener THIS test just opened (%s) to this pid - counted %d, baseline %d, want at least %d. The reported AF_INET6 number in the lifecycle test is then meaningless; check the row layout (measured on this box: 56 bytes, state at 48, pid at 52)", ln6.Addr(), seen6, base6, base6+1)
		}
		_ = ln6.Close()
		gone6, _, err := listenRowsForPIDs(windows.AF_INET6, pids)
		if err != nil {
			t.Fatalf("AC#3 ruler could not re-read the IPv6 table after Close: %v", err)
		}
		if gone6 != base6 {
			t.Errorf("after Close the IPv6 table still counts %d LISTEN rows for this pid (baseline %d)", gone6, base6)
		}
	} else {
		t.Logf("AC#3 IPv6 control did not run on this box (no [::1] listener available): %v", err6)
	}

	if err := ln4.Close(); err != nil {
		t.Fatalf("closing the control listener: %v", err)
	}
	gone4, _, err := listenRowsForPIDs(windows.AF_INET, pids)
	if err != nil {
		t.Fatalf("AC#3 ruler could not re-read the IPv4 table after Close: %v", err)
	}
	if gone4 != base4 {
		t.Errorf("after Close the ruler still counts %d IPv4 LISTEN socket(s) for this pid (baseline %d) - the ruler is counting something it cannot attribute, which is the same defect in the other direction", gone4, base4)
	}
	t.Logf("AC#3 positive control verdict: ruler counts a real listener (%d -> %d) and stops after Close (%d)", base4, seen4, gone4)
}
