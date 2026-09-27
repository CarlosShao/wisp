//go:build windows

package tools

// ticket 162 AC#4b (dispatch r3 cell ②): a REAL `taskkill /F` in the middle of
// an fs.edit write, and what is true on disk afterwards.
//
// Why this file exists at all. AC#4 (r2, and 162-v1's own re-run of it) kills
// through `Hooks.Kill`, which RETURNS AN ERROR - so it walks Go's own cleanup:
// stageAndRename's `discard()` closes the staging file and os.Remove's it, and
// the ledger gets a "删除临时文件 …（目标从头到尾未被改动）" line. That is a real
// property of a graceful cancel and 162-v1 judged it 成立; it is NOT observable
// under a fault that runs no Go code. This file is that other half:
//
//	① the target keeps every original byte (byte-for-byte, not "no error"), and
//	  it is provably NOT a prefix of what the edit intended;
//	② the `.wisp-tmp-*` residue is REAL, named one by one (path, size, exact
//	  bytes), and stays on disk: this test never deletes anything
//	  (只建不删), and it asserts the residue BEFORE any later write can sweep it
//	  away - sweeping is fs_staging.go's job and ticket 73/A18's judgement, and a
//	  second write in this directory would silently erase this cell's evidence.
//
// Shape borrowed from the family the repo already has
// (internal/tools/bridge_a18_kill_windows_test.go, read-only here): re-exec this
// very test binary as a child, block it at a named write boundary, kill it with
// taskkill. `stagingFiles` is that file's helper, reused rather than duplicated -
// 票面 :39 明令 "删除临时件要与 fs.write 同一把尺", and listing the residue is the
// same measurement.
//
// Denominator (dispatch r3 asks the question before the answer): this case has a
// denominator ON THIS MACHINE, a Windows dev box - the proof is
// TestA18RealTaskkill… PASSing in the same package in the same run (see the r3
// evidence file §0/§5 for the timestamps), and this file's own PASS below. What it
// does NOT have is a denominator anywhere else: the build tag keeps it out of
// non-Windows compilation entirely, so on a Linux runner the case is 未经验证, not
// 通过. There is no t.Skip in this file: a missing taskkill.exe fails LOUDLY.

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	// envFSEditKill* mark the child and hand it its fixture. Named apart from
	// A18's because both files can be compiled into one binary and A18's child
	// re-run filter would otherwise see foreign env.
	envFSEditKillDir    = "WISP_162R3_CHILD_DIR"
	envFSEditKillTarget = "WISP_162R3_CHILD_TARGET"
	envFSEditKillSignal = "WISP_162R3_CHILD_SIGNAL"
	// fsEditKillStep is where the child blocks. With WriteChunk 4 the boundary
	// named "write:8" is announced BEFORE the second chunk goes out, so at the
	// moment of the signal the staging file holds exactly len(fsEditKillIntended()[:4])
	// bytes and the target still holds 100% of its original bytes.
	fsEditKillStep = "write:8"
	// fsEditKillChildTest is the one function both halves live in.
	fsEditKillChildTest = "TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue"

	// The fixture, shared by both halves so the parent's byte expectations ARE the
	// child's bytes. The replacement is far longer than the 5 bytes it names, so
	// the intended output is more than twice the file: a 4-byte chunk boundary is
	// genuinely mid-write, and its marker prefix cannot be confused with anything
	// in the original.
	fsEditKillOriginal = "alpha\nbravo\ncharlie\ndelta\necho\n"
	fsEditKillOld      = "alpha"
	fsEditKillNew      = "ALPHA-KILLED-MID-WITHOUT-RUNNING-ANY-GO-CLEANUP\n"
)

// fsEditKillIntended is the byte string the edit WANTED to land: the same single
// literal replacement both halves ask for, computed once so the assertions below
// ("the target is not a prefix of this") are arithmetic, not prose.
func fsEditKillIntended() string {
	return strings.Replace(fsEditKillOriginal, fsEditKillOld, fsEditKillNew, 1)
}

// fsEditKillArgs renders the child's call. Kept in one place so the child cannot
// drift from what the parent claims was interrupted.
func fsEditKillArgs(t *testing.T, target string) string {
	t.Helper()
	return editArgs(t, target, map[string]any{"old": fsEditKillOld, "new": fsEditKillNew})
}

// fsEditKillChildBody is the child half: one real fs.edit through a real bridge,
// blocked forever at fsEditKillStep. If the block is ever bypassed the child
// says so and exits nonzero - a kill that interrupts nothing would make every
// assertion in this file vacuous.
func fsEditKillChildBody(t *testing.T, dir string) {
	t.Helper()
	target := os.Getenv(envFSEditKillTarget)
	signal := os.Getenv(envFSEditKillSignal)
	deps := func(d *FSDeps) {
		d.WriteChunk = 4
		d.Hooks = Hooks{AtStep: func(step string) {
			if step != fsEditKillStep {
				return
			}
			if err := os.WriteFile(signal, []byte(step), 0o600); err != nil {
				fmt.Fprintf(os.Stderr, "162r3 child could not signal: %v\n", err)
				os.Exit(4)
			}
			select {} // until the OS kills this process
		}}
	}
	b, _ := fsEditBridgeDeps(t, mustCanonical(t, dir), deps)
	if _, err := b.Execute(t.Context(), req("fs.edit", fsEditKillArgs(t, target))); err == nil {
		fmt.Fprintln(os.Stderr, "162r3 child finished the edit without blocking: the seam never fired")
		os.Exit(5)
	}
	fmt.Fprintln(os.Stderr, "162r3 child returned instead of being killed")
	os.Exit(6)
}

// spawnFSEditKillWriter starts the child, waits for it to report the boundary,
// kills it with a real taskkill /F and reaps it.
func spawnFSEditKillWriter(t *testing.T, dir, target, signal string) string {
	t.Helper()
	// `signal` lives under this run's own t.TempDir(), so it cannot exist yet and
	// there is nothing to clean: this test issues NO delete of its own (只建不删),
	// and the residue assertion below is the only thing on disk it cares about.
	cmd := exec.Command(os.Args[0], "-test.run=^"+fsEditKillChildTest+"$", "-test.timeout=90s")
	cmd.Env = append(os.Environ(),
		envFSEditKillDir+"="+dir, envFSEditKillTarget+"="+target, envFSEditKillSignal+"="+signal)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("启动写盘子进程失败: %v", err)
	}
	pid := cmd.Process.Pid
	killed, waited := false, false
	defer func() {
		if !killed { // never orphan a child blocked in select{}
			_, _ = exec.Command("taskkill", "/F", "/PID", fmt.Sprint(pid)).CombinedOutput()
			killed = true
		}
		if !waited {
			waited = true
			_ = cmd.Wait()
		}
	}()

	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(signal); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("子进程 %d 在 30s 内没有报告 %s 这道边界（它没跑到写盘中间）：%s", pid, fsEditKillStep, buf.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
	out, err := exec.Command("taskkill", "/F", "/PID", fmt.Sprint(pid)).CombinedOutput()
	if err != nil {
		t.Fatalf("前置条件缺失：真 taskkill 跑不动（taskkill /F /PID %d: %v: %s）。\n"+
			"需要 Windows 自带的 System32\\taskkill.exe，且当前用户能终止自己的进程。\n"+
			"这一发不许退化成进程内 Hooks.Kill——那正是 AC#4b 要补的那一半。", pid, err, out)
	}
	killed = true
	waited = true
	_ = cmd.Wait() // drains the copy goroutines, so buf is readable now
	s := buf.String()
	if strings.Contains(s, "without blocking") || strings.Contains(s, "returned instead") {
		t.Fatalf("子进程没被卡在写盘中间，kill 无意义：%s", s)
	}
	return s
}

// TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue is AC#4b:
// (a) the target holds the ORIGINAL bytes, byte for byte, and (b) the residue is
// named instead of swept under a "clean directory" claim.
func TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue(t *testing.T) {
	if dir := os.Getenv(envFSEditKillDir); dir != "" {
		fsEditKillChildBody(t, dir)
		return
	}
	work := sealableTempDir124(t)
	dir := filepath.Join(work, "edit")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "keep.txt")
	if err := os.WriteFile(target, []byte(fsEditKillOriginal), 0o600); err != nil {
		t.Fatal(err)
	}
	intended := fsEditKillIntended()
	// The fixture has to answer the question the AC asks: the original and the
	// intended bytes must be tellable apart AT THE FRONT, or "not a prefix of the
	// new content" would be true of any content on disk.
	if strings.HasPrefix(fsEditKillOriginal, intended[:8]) || strings.HasPrefix(intended, fsEditKillOriginal) {
		t.Fatalf("unusable fixture: original and intended cannot be told apart")
	}
	signal := filepath.Join(work, ".162r3-staged")

	childOut := spawnFSEditKillWriter(t, dir, target, signal)
	t.Logf("AC#4b 子进程输出（应只有 taskkill 的强制终止，没有任何 Go 清理）：%q", childOut)

	// -----------------------------------------------------------------------
	// (a) the headline, at byte level: the target is the pre-edit file.
	// -----------------------------------------------------------------------
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("真 kill 之后目标读不到了: %v", err)
	}
	if string(got) != fsEditKillOriginal {
		t.Fatalf("AC#4b(a) violated：硬杀之后目标不是原文。got %d 字节 % x\nwant %d 字节 % x\n"+
			"（既不是原文也不是全量新内容＝半个新内容＝损坏）",
			len(got), []byte(got), len(fsEditKillOriginal), []byte(fsEditKillOriginal))
	}
	if bytes.HasPrefix(got, []byte(intended)) || bytes.Contains(got, []byte(fsEditKillNew)) {
		t.Fatalf("AC#4b(a) violated: the target carries the new content, so the kill interrupted nothing: %q", got)
	}
	if len(got) != len(fsEditKillOriginal) || bytes.Equal(got, []byte(intended)) {
		t.Fatalf("AC#4b(a) violated: size %d, want %d (intended is %d)", len(got), len(fsEditKillOriginal), len(intended))
	}
	t.Logf("AC#4b(a) 字节级读数：目标 %d 字节＝改前原文（%q），要落的 %d 字节 %q 一字未出现在盘上目标里",
		len(got), got, len(intended), intended)

	// -----------------------------------------------------------------------
	// (b) the residue, named. Asserted HERE and NOW: the next write into this
	// directory would run fs_staging.go's sweeper and make this unreadable, and
	// deleting it is not this test's to do (只建不删).
	// -----------------------------------------------------------------------
	residue := stagingFiles(t, dir)
	if len(residue) != 1 {
		t.Fatalf("AC#4b(b) violated：真硬杀之后 %s 目录里的 .wisp-tmp-* 残件 = %d 枚（%v），want 1 枚。\n"+
			"0 枚＝taskkill 没打断任何事，或者有别的东西跑了清理；多枚＝计数尺坏了。\n"+
			"这一枚是硬杀不跑 Go 清理的直接凭据（Hooks.Kill 那一族会在台账里写\"删除临时文件\"，见 AC#4）。",
			dir, len(residue), residue)
	}
	var sizes []int
	for _, p := range residue {
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			t.Fatalf("ReadFile(%s): %v", p, rerr)
		}
		st, serr := os.Stat(p)
		if serr != nil {
			t.Fatalf("Stat(%s): %v", p, serr)
		}
		sizes = append(sizes, len(b))
		// The child blocked at "write:8" with WriteChunk 4, so the staging file
		// holds the FIRST 4 bytes of the intended content and nothing else: the
		// half-written bytes exist, they just exist in the right file.
		if string(b) != intended[:4] {
			t.Fatalf("残件内容不是预期那 4 字节：%s 存了 % x，want % x（%q）",
				p, []byte(b), []byte(intended[:4]), intended[:4])
		}
		t.Logf("AC#4b(b) 残件点名：%s（目录 %s 内）：%d 字节、mtime %s，内容 %q ＝要落的新内容的前 4 字节；"+
			"本程不删（只建不删），它是可回收的垃圾不是损坏",
			filepath.Base(p), dir, st.Size(), st.ModTime().Format("15:04:05"), string(b))
	}
	if filepath.Dir(residue[0]) != dir {
		t.Fatalf("残件不该落在目标目录之外：%v", residue)
	}

	// (c) the instrument a user hits next sees an intact file, not a half one.
	// Read AFTER (b) on purpose: fs.read is not a write path and does not sweep,
	// but if it ever starts to, the residue count above is already on record.
	b, _ := fsEditBridge(t, mustCanonical(t, dir))
	rd, rerr := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if rerr != nil || rd.IsError || rd.ErrorClass != "" {
		t.Fatalf("硬杀后读目标必须是干净成功：%+v err=%v", rd, rerr)
	}
	if !strings.Contains(rd.Text, "alpha") || strings.Contains(rd.Text, "ALPHA-KILLED") {
		t.Fatalf("fs.read 看见了不该看见的东西：%q", rd.Text)
	}
	t.Logf("AC#4b(c) 读数：事后 fs.read IsError=false ErrorClass=\"\"，内容 %q（原文完整）；"+
		"目录里 %d 枚残件（%v，各 %v 字节）仍在盘上等人回收", rd.Text, len(residue), residue, sizes)
}
