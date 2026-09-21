//go:build windows

package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// Ticket 73's sweeper is a DELETE primitive pointed at a directory the model
// chose. Three things have to be true of it and each gets its own case here:
// it only ever touches files it can attribute to itself (AC#2), it never
// traverses a reparse point (AC#3), and a file a live process is holding
// survives without turning into a failed write (AC#4).
//
// Every case runs the REAL bridge (real C26 canonicalizer, real C19 assessor,
// allow-all gate), because the sweeper hangs off fs.write's own path and a unit
// test of sweepStagingOrphans alone would not prove it is called at the place
// that owns the target directory.
//
// No t.Skip anywhere: a fixture this machine cannot build must read as a
// failure with its prerequisite spelled out (see bridge_junction_windows_test.go
// for why this repository insists on that).

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// sweepChildStep* are the re-exec contract of the handle-holding child below.
const (
	envSweepHoldFile   = "WISP_SWEEP_HOLD_FILE"
	envSweepHoldSignal = "WISP_SWEEP_HOLD_SIGNAL"
	sweepHoldTest      = "TestSweepSparesATempFileHeldByAnotherProcess"
)

// deadPID returns the pid of a process this test started and already reaped,
// after PROVING the sweeper classifies it as gone. The proof is not decoration:
// if the OS recycles that pid between here and the sweep, the case that uses it
// must say so rather than pass on a coincidence.
func deadPID(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("前置条件缺失：起不了一个用来烧 pid 的子进程（%v）——"+
			"本用例需要能创建并等待一个真进程", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		// `cmd /c exit 0` exits 0; anything else is the fixture failing, not
		// the behavior under test.
		t.Fatalf("烧 pid 的子进程 %d 没正常退出: %v", pid, err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		alive, err := stagingCreatorAlive(fmt.Sprint(pid))
		if err != nil {
			t.Fatalf("前置条件缺失：判活本身报错（pid %d: %v）⇒ 清扫器会把它当活的，"+
				"下面的孤儿断言就没有意义", pid, err)
		}
		if !alive {
			return pid
		}
		if time.Now().After(deadline) {
			t.Fatalf("pid %d 退出 5s 后仍被判为存活（OpenProcess 还打得开）⇒ "+
				"这条用例的孤儿形状造不出来，不是清扫器的错，是 fixture 前提不成立", pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// plantOrphan puts an ATTRIBUTABLE staging file in dir as if a killed session
// had left it, and returns its path.
func plantOrphan(t *testing.T, dir string, pid int, random string) string {
	t.Helper()
	name := stagingNameOf(stagingOwner(), pid, random)
	if !stagingAttributable(name) {
		t.Fatalf("fixture 造出了自己都不认的名字 %q：命名方案与匹配器不同源了", name)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("ORPHAN"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// writeFileOf puts content in dir/name and returns the path.
func writeFileOf(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// bridgeWrite runs one fs.write through the real bridge over allowed and fails
// loudly if the call itself fails - AC#4 needs "the write still succeeded", and
// AC#2/#3 need "the sweep really happened inside a write that worked".
func bridgeWrite(t *testing.T, allowed, path, content string) {
	t.Helper()
	b, _, _ := allowAllBridge(t, mustCanonical(t, allowed), false)
	out, err := b.Execute(t.Context(), req("fs.write", writeArgsOf(t, path, content)))
	if err != nil {
		t.Fatalf("fs.write 走桥失败: %v", err)
	}
	if out.IsError {
		t.Fatalf("fs.write 被拒了，后面的清扫断言就没有意义: %+v", out)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("写盘没落地 %s: %v", path, err)
	}
}

// ---------------------------------------------------------------------------
// AC#2 — attribution, not "delete everything that looks like ours"
// ---------------------------------------------------------------------------

// TestSweepReclaimsOnlyItsOwnStagingFiles is AC#2: a foreign file that shares
// Wisp's staging PREFIX survives, while the attributable orphan in the very
// same directory does not. The foreign files are named in the assertions, so a
// future "just sweep the prefix" change names what it destroyed.
func TestSweepReclaimsOnlyItsOwnStagingFiles(t *testing.T) {
	dir := filepath.Join(tempRaw(t), "auth")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	dead := deadPID(t)
	orphan := plantOrphan(t, dir, dead, "4242")

	// Each of these is either the wrong shape or the wrong owner, and one is the
	// right shape owned by a process that is STILL RUNNING (this test binary).
	foreign := map[string]string{
		".wisp-tmp-notes.txt":                                     "别人的文件，只是撞了前缀",
		".wisp-tmp-.txt":                                          "空字段也不许被当成随机段",
		".wisp-tmp-deadbeef-1-1":                                  "形状全对，owner 不是本程序",
		".wisp-tmp-" + stagingOwner() + "-notanumber-1":           "pid 段不是数字",
		".wisp-tmp-" + stagingOwner() + "-1":                      "少一段，不是一等公民的形状",
		stagingNameOf(stagingOwner(), os.Getpid(), "7") + "-tail": "随机段后面再挂东西",
	}
	for name, body := range foreign {
		writeFileOf(t, dir, name, body)
	}
	liveOwn := plantOrphan(t, dir, os.Getpid(), "1") // ours, but WE are alive

	bridgeWrite(t, dir, filepath.Join(dir, "note.txt"), "触发清扫的那一次写盘")

	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("该扫掉的孤儿没被扫掉 %s（err=%v）⇒ 清扫器根本没跑", orphan, err)
	}
	if _, err := os.Stat(liveOwn); err != nil {
		t.Fatalf("创建者进程还活着（本测试二进制自己的 pid），孤儿位 %s 却被删了: %v",
			liveOwn, err)
	}
	for name, body := range foreign {
		p := filepath.Join(dir, name)
		got, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("清扫器删掉了不可归因的文件 %s: %v", p, err)
		}
		if string(got) != body {
			t.Fatalf("清扫器改写了 %s 的内容: %q", p, got)
		}
	}
	if got, want := stagingFiles(t, dir), len(foreign)+1; len(got) != want {
		t.Fatalf("清扫后目录里 .wisp-tmp-* 数量 = %d, want %d（%d 个外人文件 + 1 个"+
			"本程序活 pid 的孤儿）: %v", len(got), want, len(foreign), got)
	}
}

// TestSweepNamingSchemeRejectsNearMisses pins the string half of attribution
// without touching the disk, so AC#2's failure message can never be blamed on
// an Lstat surprise.
func TestSweepNamingSchemeRejectsNearMisses(t *testing.T) {
	dead := deadPID(t)
	ok := stagingNameOf(stagingOwner(), dead, "abc123")
	cases := []struct {
		name string
		want bool
		why  string
	}{
		{ok, true, "本程序 + 真 pid + 随机段"},
		{stagingNameOf("00000000", dead, "1"), false, "owner 段不是本程序的 token"},
		{stagingNameOf(stagingOwner(), 0, "1"), false, "pid 段是 0"},
		{".wisp-tmp-" + stagingOwner() + "-" + fmt.Sprint(dead) + "-", false, "随机段为空"},
		{".wisp-tmp-" + stagingOwner() + "-" + fmt.Sprint(dead), false, "只有两段"},
		{".wisp-tmp-" + stagingOwner() + "-" + fmt.Sprint(dead) + "-1-2", false, "三段以上"},
		{".wisp-tmp-" + stagingOwner() + "-" + fmt.Sprint(dead) + `-a\b`, false, "随机段带路径分隔符"},
		{"readme.md", false, "根本不是暂存名"},
	}
	for _, tc := range cases {
		if got := stagingAttributable(tc.name); got != tc.want {
			t.Errorf("stagingAttributable(%q) = %v, want %v（%s）", tc.name, got, tc.want, tc.why)
		}
	}
	if pid := stagingAttribution(ok); pid != fmt.Sprint(dead) {
		t.Errorf("stagingAttribution(%q) = %q, want 创建者 pid %d", ok, pid, dead)
	}
}

// ---------------------------------------------------------------------------
// AC#3 — a reparse point is never traversed
// ---------------------------------------------------------------------------

// TestSweepNeverDeletesThroughARealJunction is AC#3: an ORPHAN-SHAPED, fully
// attributable NAME that is a real NTFS junction pointing outside every allowed
// root. The sweep must leave the junction's target alone - and the target file
// is proven to exist BEFORE the write, so "still there" cannot mean "never was".
func TestSweepNeverDeletesThroughARealJunction(t *testing.T) {
	allowed := tempRaw(t)
	outside := tempRaw(t)
	outsideCanon := mustCanonical(t, outside)
	secret := "junction 对面的真内容，清扫器不许碰到"
	victim := writeFileOf(t, outside, "victim.txt", secret)
	victimCanon := mustCanonical(t, victim)

	// The link's NAME is the orphan shape: a sweeper that keys off the name and
	// then unlinks what it finds would take this as its own.
	dead := deadPID(t)
	link := filepath.Join(allowed, stagingNameOf(stagingOwner(), dead, "777"))
	mkRealJunction(t, link, outside)
	t.Cleanup(func() { _ = os.Remove(link) }) // drop the link before TempDir cleanup

	// POSITIVE CONTROL ①: the junction really reaches the bytes.
	got, err := os.ReadFile(filepath.Join(link, "victim.txt"))
	if err != nil || string(got) != secret {
		t.Fatalf("阳性对照失败：junction 没通向 %s（%v）⇒ 后面的\"没删\"什么也没证明",
			victimCanon, err)
	}
	// POSITIVE CONTROL ②: a plain attributable orphan in the SAME directory gets
	// reclaimed by this very write, so a pass below cannot mean "the sweep never
	// ran".
	orphan := plantOrphan(t, allowed, dead, "778")

	bridgeWrite(t, allowed, filepath.Join(allowed, "note.txt"), "触发清扫的写盘")

	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("同目录的普通孤儿都没被扫掉（err=%v）⇒ 清扫没跑，上面的\"没删\"就是空跑", err)
	}
	after, err := os.ReadFile(victim)
	if err != nil {
		t.Fatalf("junction 对面的文件在清扫之后读不到了 %s: %v", victimCanon, err)
	}
	if string(after) != secret {
		t.Fatalf("清扫器改写了 junction 对面的文件 %s", victimCanon)
	}
	if _, err := os.Lstat(link); err != nil {
		t.Fatalf("清扫器连 junction 本身都删了（reparse point 一律不碰）: %v", err)
	}
	if _, err := os.Stat(outsideCanon); err != nil {
		t.Fatalf("授权目录之外的那个目录本身出问题了: %v", err)
	}
}

// ---------------------------------------------------------------------------
// AC#4 — a staging file another process is holding open
// ---------------------------------------------------------------------------

// sweepHoldChildBody is the child half: open the named file with a share mode
// that REFUSES FILE_SHARE_DELETE and block. That is the exact shape of a live
// writer's in-flight staging file - the handle exists, so the unlink fails.
func sweepHoldChildBody(t *testing.T) {
	path, signal := os.Getenv(envSweepHoldFile), os.Getenv(envSweepHoldSignal)
	if path == "" || signal == "" {
		fmt.Fprintln(os.Stderr, "SWEEP child: 缺 env")
		os.Exit(3)
	}
	p16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SWEEP child: UTF16PtrFromString(%q): %v\n", path, err)
		os.Exit(3)
	}
	h, err := windows.CreateFile(p16, windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil,
		windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SWEEP child: CreateFile(%q): %v\n", path, err)
		os.Exit(4)
	}
	buf := []byte("HELD")
	var written uint32
	if err := windows.WriteFile(h, buf, &written, nil); err != nil {
		fmt.Fprintf(os.Stderr, "SWEEP child: WriteFile(%q): %v\n", path, err)
		os.Exit(4)
	}
	if err := os.WriteFile(signal, []byte("held"), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "SWEEP child: signal: %v\n", err)
		os.Exit(4)
	}
	select {} // hold the handle until the OS kills us
}

// TestSweepSparesATempFileHeldByAnotherProcess is AC#4. Two shapes, because the
// sweeper has two independent guards and each must be seen biting:
//
//	A: the name embeds the LIVE holder's pid ⇒ the liveness check skips it;
//	B: the name embeds a DEAD pid while a live process holds the handle open ⇒
//	 only the retry-then-skip path can save it, and the write must still succeed.
//
// B then kills the holder and writes again: the file goes, which proves the
// survival in B was the HANDLE and not the name.
func TestSweepSparesATempFileHeldByAnotherProcess(t *testing.T) {
	if os.Getenv(envSweepHoldFile) != "" {
		sweepHoldChildBody(t)
		return
	}
	dir := filepath.Join(tempRaw(t), "auth")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}

	// ---- A: live creator pid, nobody holding it.
	live := plantOrphan(t, dir, os.Getpid(), "1")
	bridgeWrite(t, dir, filepath.Join(dir, "a.txt"), "A：创建者还活着")
	if _, err := os.Stat(live); err != nil {
		t.Fatalf("创建者进程（本测试二进制）还活着，它的暂存文件却被删了: %v", err)
	}

	// ---- B: dead pid in the name, a real process holding the handle.
	held := plantOrphan(t, dir, deadPID(t), "2")
	proc := startSweepHolder(t, held)
	defer func() {
		_, _ = exec.Command("taskkill", "/F", "/PID", fmt.Sprint(proc.Process.Pid)).CombinedOutput()
		_ = proc.Wait()
	}()

	bridgeWrite(t, dir, filepath.Join(dir, "b.txt"), "B：文件被别的活进程占着")
	if _, err := os.Stat(held); err != nil {
		t.Fatalf("正被活进程持有的暂存文件被删了（名字里的 pid 已死 ⇒ 只有重试后跳过这一道防线）: %v", err)
	}
	// The child wrote through its handle with OPEN_EXISTING (no truncate), so the
	// first four bytes are its own and the tail is still the planted ORPHAN:
	// "HELD..." is the proof the survivor is the file the child really has open.
	if got := readString(t, held); !strings.HasPrefix(got, "HELD") {
		t.Fatalf("被占用文件的内容不是持有者写的那份: %q", got)
	}

	// Now release it: the NEXT write must reclaim it, so B's survival above was
	// the handle and not a name that happened not to match. Real taskkill, same
	// hammer A18 uses - the point is the handle goes away, not how.
	out, err := exec.Command("taskkill", "/F", "/PID", fmt.Sprint(proc.Process.Pid)).CombinedOutput()
	if err != nil {
		t.Fatalf("结束持有者失败（后面的\"释放后就能扫掉\"对照做不了）: %v: %s", err, out)
	}
	_ = proc.Wait()
	// The delete can still be pending for a moment after the last handle closes,
	// so wait until the file is openable again (Go opens with FILE_SHARE_DELETE,
	// which is exactly what the holder above deliberately did NOT do).
	deadline := time.Now().Add(5 * time.Second)
	for {
		f, err := os.OpenFile(held, os.O_RDONLY, 0o600)
		if err == nil {
			_ = f.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("持有者已退出 5s，%s 仍打不开 ⇒ 无法验证\"释放后能扫掉\": %v", held, err)
		}
		time.Sleep(20 * time.Millisecond)
	}
	bridgeWrite(t, dir, filepath.Join(dir, "c.txt"), "C：占用者已经死了")
	if _, err := os.Stat(held); !os.IsNotExist(err) {
		t.Fatalf("持有者退出后下一次写盘仍没扫掉它（err=%v）⇒ 上一条的\"幸存\"是名字不对，不是句柄", err)
	}
}

// startSweepHolder re-execs this test binary to hold `path` open and waits for
// it to report that the handle is live.
func startSweepHolder(t *testing.T, path string) *exec.Cmd {
	t.Helper()
	signal := filepath.Join(filepath.Dir(path), ".sweep-held-"+filepath.Base(path))
	if err := os.Remove(signal); err != nil && !os.IsNotExist(err) {
		t.Fatalf("清理信号文件 %s: %v", signal, err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^"+sweepHoldTest+"$", "-test.timeout=120s")
	cmd.Env = append(os.Environ(), envSweepHoldFile+"="+path, envSweepHoldSignal+"="+signal)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("重启测试二进制当持有者失败: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(signal); err == nil {
			t.Cleanup(func() { _ = os.Remove(signal) })
			return cmd
		}
		if time.Now().After(deadline) {
			_, _ = exec.Command("taskkill", "/F", "/PID", fmt.Sprint(cmd.Process.Pid)).CombinedOutput()
			_ = cmd.Wait()
			t.Fatalf("持有进程 %d 在 30s 内没报告句已打开：%s", cmd.Process.Pid, buf.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
}
