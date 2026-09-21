//go:build windows

package tools

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

// A18 (ticket 73 flipped ③): a REAL `taskkill /F` during a staged write, and
// what is on disk afterwards.
//
// The existing `TestAtomicWriteKillsMidWrite` kills IN-PROCESS via
// `Hooks.Kill`, which returns an error and therefore RUNS the Go cleanup — so
// its "no staging file left" assertion is unobservable under a real fault
// (that is A18's exact objection). This file kills a real child OS process
// with a real `taskkill /F` and records what is true now:
//
//	① the target is complete-or-absent (D31 holds: staging + one os.Rename),
//	② each interrupted write leaves exactly one `.wisp-tmp-*` behind — the
//	   residue is REAL, and this test still asserts it exists BEFORE anything
//	   can sweep it (by the third subtest exactly one remains: the second
//	   child's own write already reclaimed the first child's orphan),
//	③ the NEXT write through the bridge reclaims it (internal/tools/fs_staging
//	   go's sweeper, ticket 73): the directory ends with zero `.wisp-tmp-*`
//	   entries and the destination file is still byte-exact.
//
// ③ used to be pinned as "the orphan is still there" (761447f), because the
// owner had not picked between self-healing and declaring the residue
// acceptable. Registry Q-16 picked self-healing, so this file now pins the
// opposite; ② is kept as the positive control that makes ③ mean something — a
// sweep that ran before the residue existed would look identical to a sweep
// that never ran if only ③ were asserted.
//
// Prerequisites: Windows with `taskkill.exe` (System32 on every SKU) and a
// test binary that can re-exec itself. Missing either fails LOUDLY — this file
// contains no t.Skip.

const (
	// envA18Dir marks a child process and hands it its working directory.
	envA18Dir = "WISP_A18_CHILD_DIR"
	// envA18Target is the file the child writes; envA18Signal is where it
	// reports that the staging file exists; envA18Step names the boundary.
	envA18Target = "WISP_A18_CHILD_TARGET"
	envA18Signal = "WISP_A18_CHILD_SIGNAL"
	envA18Step   = "WISP_A18_CHILD_STEP"
	// a18StagedStep is the boundary the child blocks at. fs.write runs in
	// 4-byte chunks here, so by the time this step is announced one chunk has
	// already gone into the staging file: the kill interrupts a real mid-write.
	a18StagedStep = "write:8"
	// a18ChildTest is the test function both halves live in; the child re-runs
	// only this one.
	a18ChildTest = "TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile"
)

// a18ChildBody is the child half: a real fs.write through a real bridge whose
// gate approves, blocked forever at a mid-write boundary. The parent kills it.
// If the block is ever bypassed the child says so and exits nonzero, because a
// kill that interrupts nothing would make the whole characterization vacuous.
func a18ChildBody(t *testing.T, dir string) {
	t.Helper()
	target := os.Getenv(envA18Target)
	signal := os.Getenv(envA18Signal)
	step := os.Getenv(envA18Step)
	if step == "" {
		step = a18StagedStep
	}
	deps := FSDeps{
		Paths:      NewPathCanonicalizer([]string{mustCanonical(t, dir)}, nil),
		WriteChunk: 4,
		Hooks: Hooks{AtStep: func(s string) {
			if s != step {
				return
			}
			// "Bytes are staged and the target is still untouched."
			if err := os.WriteFile(signal, []byte(s), 0o600); err != nil {
				fmt.Fprintf(os.Stderr, "A18 child signal write failed: %v\n", err)
				os.Exit(4)
			}
			select {} // block until the OS kills this process
		}},
	}
	b, _ := fsDepsBridge(t, deps,
		&gateSpy{windowAns: AnswerAllow, approveAns: AnswerAllow})
	rq := req("fs.write", writeArgsOf(t, target, strings.Repeat("NEW-", 64)))
	if _, err := b.Execute(t.Context(), rq); err != nil {
		fmt.Fprintf(os.Stderr, "A18 child write returned: %v\n", err)
	}
	fmt.Fprintln(os.Stderr, "A18 child wrote without blocking: the seam did not fire")
	os.Exit(5)
}

// spawnA18Writer starts the child, waits until it reports a staged file, then
// kills it with a REAL taskkill /F. It returns the child's own output.
func spawnA18Writer(t *testing.T, dir, target string) string {
	t.Helper()
	signal := filepath.Join(filepath.Dir(dir), ".a18-staged-"+filepath.Base(target))
	if err := os.Remove(signal); err != nil && !os.IsNotExist(err) {
		t.Fatalf("清理信号文件 %s: %v", signal, err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^"+a18ChildTest+"$", "-test.timeout=90s")
	cmd.Env = append(os.Environ(),
		envA18Dir+"="+dir, envA18Target+"="+target,
		envA18Signal+"="+signal, envA18Step+"="+a18StagedStep)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		t.Fatalf("启动写盘子进程失败: %v", err)
	}
	pid := cmd.Process.Pid
	killed, waited := false, false
	reap := func() {
		if !killed { // never orphan a child blocked in select{}
			_, _ = exec.Command("taskkill", "/F", "/PID", fmt.Sprint(pid)).CombinedOutput()
			killed = true
		}
		if !waited {
			waited = true
			_ = cmd.Wait()
		}
	}
	defer reap()

	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(signal); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("子进程 %d 在 30s 内没有报告暂存文件（它没跑到写盘边界）：%s", pid, buf.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
	out, err := exec.Command("taskkill", "/F", "/PID", fmt.Sprint(pid)).CombinedOutput()
	if err != nil {
		t.Fatalf("前置条件缺失：真 taskkill 跑不动（taskkill /F /PID %d: %v: %s）。\n"+
			"需要 Windows 自带的 System32\\taskkill.exe 可用，且当前用户能终止自己的进程。\n"+
			"这条用例不许退化成进程内 Hooks.Kill——那正是 A18 登记的假象来源。", pid, err, out)
	}
	killed = true
	waited = true
	_ = cmd.Wait() // also drains the copy goroutines, so buf is safe to read now
	s := buf.String()
	if strings.Contains(s, "wrote without blocking") {
		t.Fatalf("子进程没被卡在写盘中间，kill 无意义：%s", s)
	}
	return s
}

// stagingFiles lists Wisp's staging files in dir — the residue A18 is about.
func stagingFiles(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", dir, err)
	}
	var out []string
	for _, e := range des {
		if strings.HasPrefix(e.Name(), tempPrefix) {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out
}

// TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile
// is the characterization plus ticket 73's fix: D31's promise survives a real
// kill, the residue is real, and the next write through the bridge reclaims it.
func TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile(t *testing.T) {
	if dir := os.Getenv(envA18Dir); dir != "" {
		a18ChildBody(t, dir)
		return
	}
	work := tempRaw(t)
	auth := filepath.Join(work, "auth")
	if err := os.MkdirAll(auth, 0o700); err != nil {
		t.Fatal(err)
	}
	existing := filepath.Join(auth, "existing.txt")
	const oldBytes = "OLD-BYTES：这一版内容必须逐字活着穿过一次真 kill"
	if err := os.WriteFile(existing, []byte(oldBytes), 0o600); err != nil {
		t.Fatal(err)
	}
	fresh := filepath.Join(auth, "fresh.txt")

	t.Run("overwrite_branch_target_keeps_every_old_byte", func(t *testing.T) {
		spawnA18Writer(t, auth, existing)
		got, err := os.ReadFile(existing)
		if err != nil {
			t.Fatalf("目标文件在真 kill 之后读不到了: %v", err)
		}
		if string(got) != oldBytes {
			t.Fatalf("D31 violated under a REAL kill: 目标成了半截内容（%d 字节）", len(got))
		}
	})

	t.Run("new_file_branch_target_never_appears", func(t *testing.T) {
		spawnA18Writer(t, auth, fresh)
		if _, err := os.Stat(fresh); err == nil {
			t.Fatal("半截的新文件出现在了目标路径上")
		} else if !os.IsNotExist(err) {
			t.Fatalf("Stat(%s): %v", fresh, err)
		}
	})

	// The A18 half, flipped by ticket 73: two real kills left two staging files
	// behind (asserted FIRST, as the positive control), and the next write
	// through the bridge reclaims exactly those orphans.
	t.Run("the_next_write_reclaims_the_staging_residue", func(t *testing.T) {
		before := stagingFiles(t, auth)
		if len(before) != 1 {
			t.Fatalf("真 kill 之后的 .wisp-tmp-* 残留 = %v, want 1 个。"+
				"这条是 ③ 的阳性对照：kill 本身没留下残口的话，\"扫干净了\"就什么都不是。"+
				"1 而不是 761447f 时代的 2，是因为第二个子进程自己写盘前就把第一个的孤儿扫掉了"+
				"（清扫挂在每次写盘上，见 fs_write.go 的 reclaimStaging）⇒ 到这一步盘上必须"+
				"只剩第二个 kill 的那一个残留，多一个少一个都说明行为变了", before)
		}
		for _, p := range before {
			st, err := os.Stat(p)
			if err != nil {
				t.Fatalf("Stat(%s): %v", p, err)
			}
			t.Logf("残留暂存文件 %s：%d 字节（子进程被杀时已经写进暂存区的字节数）",
				filepath.Base(p), st.Size())
		}
		// "Next startup", as the code exists today: a brand-new canonicalizer,
		// registry, assessor and bridge over the same root, doing a successful
		// write into that same directory.
		b, _, _ := allowAllBridge(t, mustCanonical(t, auth), false)
		out, err := b.Execute(t.Context(), req("fs.write", writeArgsOf(t,
			filepath.Join(auth, "after-restart.txt"), "重启后的第一次成功写盘")))
		if err != nil {
			t.Fatal(err)
		}
		if out.IsError {
			t.Fatalf("新桥的合法写入被拒了，这条对照就没意义: %+v", out)
		}
		if after := stagingFiles(t, auth); len(after) != 0 {
			t.Fatalf("下一次写盘没有清扫历史残留 (%v -> %v)：A18 的残口仍然开着——"+
				"清扫器（internal/tools/fs_staging.go）被摘掉了或者归因失败", before, after)
		}
		if _, err := os.Stat(filepath.Join(auth, "after-restart.txt")); err != nil {
			t.Fatalf("新桥的正常写入没落地: %v", err)
		}
		// The sweep deletes files in the same directory as the destination, so
		// "the destination is still whole" is part of THIS assertion, not
		// something subtest ① is assumed to have covered.
		got, err := os.ReadFile(existing)
		if err != nil {
			t.Fatalf("清扫之后目标读不到了: %v", err)
		}
		if string(got) != oldBytes {
			t.Fatalf("清扫器动了目标文件：existing.txt 不再是 OLD-BYTES（%d 字节）", len(got))
		}
	})
}
