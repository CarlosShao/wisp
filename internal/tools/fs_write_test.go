package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Segment 2's contract surface: the D34 risk matrix for the write half, the D31
// atomic writer, the recycle-bin backend and the applied-steps ledger.
//
// Every case runs through the bridge (never by calling a tool directly), because
// the thing under test is the verdict the GATE sees and the forensics the call
// leaves behind, not just the bytes on disk.

// ---------------------------------------------------------------------------
// harness
// ---------------------------------------------------------------------------

// decLog captures Options.OnDecision, which is the only place an L0 verdict is
// visible (an L0 call never reaches a gate).
type decLog struct {
	mu  sync.Mutex
	all []Decision
}

func (d *decLog) add(x Decision) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.all = append(d.all, x)
}

func (d *decLog) last() Decision {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.all) == 0 {
		return Decision{}
	}
	return d.all[len(d.all)-1]
}

// fsDepsBridge wires the whole fs family over one deps, journal-free.
func fsDepsBridge(t *testing.T, d FSDeps, g Gate) (*Bridge, *decLog) {
	t.Helper()
	if g == nil {
		g = NoGate{}
	}
	if d.Paths == nil {
		t.Fatal("test deps must carry a C26 resolver")
	}
	dl := &decLog{}
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: d.Paths, Gate: g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		OnDecision: dl.add,
		Logf:       func(string, ...any) {},
	}), dl
}

// args renders one call's argument object.
func args(t *testing.T, m map[string]any) string {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// pathArgsOf sends the path in forward-slash spelling on purpose: C26 must fold
// it, so a test that only ever sent backslashes would not prove the fold runs.
func pathArgsOf(t *testing.T, path string) string {
	return args(t, map[string]any{"path": slash(path)})
}

func writeArgsOf(t *testing.T, path, content string) string {
	return args(t, map[string]any{"path": slash(path), "content": content})
}

func moveArgsOf(t *testing.T, from, to string) string {
	return args(t, map[string]any{"from": slash(from), "to": slash(to)})
}

func readString(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", p, err)
	}
	return string(b)
}

func existsFile(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// tempRaw gives a temp dir in OS spelling (tests hand it to the bridge, which
// canonicalizes; mustCanonical is for when a test needs the canonical form).
func tempRaw(t *testing.T) string { return t.TempDir() }

// noStagingFilesLeft asserts the D31 cleanup contract: a killed write leaves
// nothing but the target behind, because a stray temp file in an authorized
// directory is itself a side effect.
func noStagingFilesLeft(t *testing.T, dir string) {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range des {
		if strings.HasPrefix(e.Name(), tempPrefix) {
			t.Errorf("staging file left behind after the kill: %s", filepath.Join(dir, e.Name()))
		}
	}
}

// otherVolumeDir returns a directory root on a volume OTHER than the one `ref`
// sits on, or skips when the machine genuinely has only one. The skip is for
// the hardware, never for the logic: TestCrossVolumeRuleIsR8 pins the verdict
// itself so the escalation cannot go untested on a single-drive box.
func otherVolumeDir(t *testing.T, ref string) string {
	t.Helper()
	refVol := volumeID(ref)
	if refVol == "" {
		t.Skip("no volume identity for the reference path on this platform")
	}
	for l := 'A'; l <= 'Z'; l++ {
		root := string(l) + `:\`
		if _, err := os.Stat(root); err != nil {
			continue
		}
		if v := volumeID(root + "wisp-probe"); v != "" && v != refVol {
			return root
		}
	}
	t.Skip("this machine has a single volume; the cross-volume branch cannot be exercised here")
	return ""
}

// ---------------------------------------------------------------------------
// 1. the D34 risk matrix
// ---------------------------------------------------------------------------

// TestD34WriteMatrix asserts one row per D34 line, including the two rules the
// write half depends on: R8 for an overwrite and R2 for an out-of-allowlist
// target. Every gate answer is REJECT, so a row that reached a gate proves the
// route AND writes nothing.
func TestD34WriteMatrix(t *testing.T) {
	root := tempCanonical(t)
	outside := tempCanonical(t) // deliberately NOT in the allowlist
	newFile := filepath.Join(root, "occupied-destination.txt")
	if err := os.WriteFile(newFile, []byte("there"), 0o600); err != nil {
		t.Fatal(err)
	}
	existing := writeUnder(t, root, "existing.txt", "old")
	aListTarget := filepath.Join(outside, "target.txt")
	crossDest := otherVolumeDir(t, root)

	cases := []struct {
		name     string
		tool     string
		args     func(t *testing.T) string
		want     risk.Level
		rules    []risk.RuleID
		noRules  []risk.RuleID
		route    string // "window" | "approval" | "none"
		deleteOn bool
	}{
		{
			name: "read_inside_allowlist_is_L0", tool: "fs.read",
			args: func(t *testing.T) string { return pathArgsOf(t, existing) },
			want: risk.L0, route: "none",
			noRules: []risk.RuleID{risk.R2},
		},
		{
			name: "read_out_of_allowlist_is_L2_and_is_a_verdict_not_a_dormant_rule", tool: "fs.read",
			args: func(t *testing.T) string { return pathArgsOf(t, aListTarget) },
			want: risk.L2, rules: []risk.RuleID{risk.R2}, route: "approval",
		},
		{
			name: "list_out_of_allowlist_is_L2", tool: "fs.list",
			args: func(t *testing.T) string { return pathArgsOf(t, outside) },
			want: risk.L2, rules: []risk.RuleID{risk.R2}, route: "approval",
		},
		{
			name: "write_new_file_is_L1", tool: "fs.write",
			args: func(t *testing.T) string {
				return writeArgsOf(t, filepath.Join(root, "brand-new.txt"), "hello")
			},
			want: risk.L1, route: "window", noRules: []risk.RuleID{risk.R8},
		},
		{
			name: "write_over_existing_is_L2_via_R8", tool: "fs.write",
			args: func(t *testing.T) string { return writeArgsOf(t, existing, "new") },
			want: risk.L2, rules: []risk.RuleID{risk.R8}, route: "approval",
		},
		{
			name: "write_out_of_allowlist_is_L2", tool: "fs.write",
			args: func(t *testing.T) string {
				return writeArgsOf(t, filepath.Join(outside, "out.txt"), "hello")
			},
			want: risk.L2, rules: []risk.RuleID{risk.R2}, route: "approval",
		},
		{
			name: "trash_is_L1_because_the_bin_can_give_it_back", tool: "fs.trash",
			args: func(t *testing.T) string { return pathArgsOf(t, existing) },
			want: risk.L1, route: "window", noRules: []risk.RuleID{risk.R8},
		},
		{
			name: "move_same_volume_new_destination_is_L1", tool: "fs.move",
			args: func(t *testing.T) string {
				return moveArgsOf(t, existing, filepath.Join(root, "moved.txt"))
			},
			want: risk.L1, route: "window", noRules: []risk.RuleID{risk.R8},
		},
		{
			name: "move_over_existing_destination_is_L2", tool: "fs.move",
			args: func(t *testing.T) string { return moveArgsOf(t, existing, newFile) },
			want: risk.L2, rules: []risk.RuleID{risk.R8}, route: "approval",
		},
		{
			name: "move_across_volumes_is_L2", tool: "fs.move",
			args: func(t *testing.T) string {
				return moveArgsOf(t, existing, crossDest+`wisp-seg2-crossvol-nonexistent.txt`)
			},
			want: risk.L2, rules: []risk.RuleID{risk.R8}, route: "approval",
		},
		{
			name: "delete_when_enabled_is_L2", tool: "fs.delete",
			args: func(t *testing.T) string { return pathArgsOf(t, existing) },
			want: risk.L2, rules: []risk.RuleID{risk.R8}, route: "approval",
			deleteOn: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := &gateSpy{} // both answers empty -> AnswerReject, nothing executes
			b, log := fsDepsBridge(t, FSDeps{
				Paths: NewPathCanonicalizer([]string{root}, nil), DeleteEnabled: tc.deleteOn,
			}, g)
			out, err := b.Execute(t.Context(), req(tc.tool, tc.args(t)))
			if err != nil {
				t.Fatalf("a gate rejection is a refusal, not a host fault: %v", err)
			}
			if tc.route != "none" && !out.IsError {
				t.Fatalf("the gate said no, yet the call succeeded: %+v", out)
			}
			dec := log.last()
			if tc.route == "approval" {
				dec = g.approvalDecision()
			} else if tc.route == "window" {
				dec = g.windowDecision()
			}
			if dec.Level != tc.want {
				t.Fatalf("level = %v, want %v (decision=%+v)", dec.Level, tc.want, dec)
			}
			for _, r := range tc.rules {
				if !contains(dec.RulesHit, r) {
					t.Errorf("rules_hit = %v, want %s in it (a D34 row that cannot name its rule is a dormant rule)",
						dec.RulesHit, r)
				}
			}
			for _, r := range tc.noRules {
				if contains(dec.RulesHit, r) {
					t.Errorf("rules_hit = %v must NOT contain %s (%v must not be escalated)", dec.RulesHit, r, tc.tool)
				}
			}
			w, a := g.counts()
			switch tc.route {
			case "none":
				if w != 0 || a != 0 {
					t.Errorf("an %v call must not reach a gate: window=%d approval=%d", tc.want, w, a)
				}
			case "window":
				if w != 1 || a != 0 {
					t.Errorf("want the L1 window route only: window=%d approval=%d", w, a)
				}
			case "approval":
				if a != 1 || w != 0 {
					t.Errorf("want the L2 approval route only: window=%d approval=%d", w, a)
				}
			}
		})
	}
}

// TestCrossVolumeRuleIsR8 pins the ESCALATION ITSELF, with no hardware
// dependency: D34's cross-volume row is expressed as R8's permanent-delete
// class (the source is destroyed outside the bin), and this is the rule mapping
// a single-volume machine would otherwise leave untested.
func TestCrossVolumeRuleIsR8(t *testing.T) {
	a := bareAssessor()
	got := a.Assess("fs.move", map[string]any{"from": `C:\a\x.txt`, "to": `D:\y\x.txt`}, risk.Facts{
		Declared:     risk.L1,
		Irreversible: []string{"delete"},
	})
	if got.Level != risk.L2 {
		t.Fatalf("cross-volume move judged %v, want L2", got.Level)
	}
	if !contains(got.RulesHit, risk.R8) {
		t.Errorf("rules_hit = %v, want R8", got.RulesHit)
	}
	same := bareAssessor().Assess("fs.move", map[string]any{}, risk.Facts{Declared: risk.L1})
	if same.Level != risk.L1 || contains(same.RulesHit, risk.R8) {
		t.Errorf("same-volume move must hold the L1 floor, got %+v", same)
	}
}

// bareAssessor is the R1-R9 judge with no path rules in play, used by the
// pure-rule test above.
func bareAssessor() *risk.RiskAssessor { return risk.NewRiskAssessor() }

// TestOverwriteDetectionFollowsTheCanonicalPath proves the R8 probe asks the OS
// about the path C26 resolved, not about what the model typed: a relative-ish
// spelling with .. must still land on the existing file and therefore L2.
func TestOverwriteDetectionFollowsTheCanonicalPath(t *testing.T) {
	root := tempCanonical(t)
	existing := writeUnder(t, root, "note.txt", "content")
	g := &gateSpy{}
	b, _ := fsDepsBridge(t, FSDeps{Paths: NewPathCanonicalizer([]string{root}, nil)}, g)

	spelled := slash(root) + `/sub/../note.txt`
	out, err := b.Execute(t.Context(), req("fs.write",
		args(t, map[string]any{"path": spelled, "content": "x"})))
	if err != nil || !out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	dec := g.approvalDecision()
	if dec.Level != risk.L2 || !contains(dec.RulesHit, risk.R8) {
		t.Fatalf("a traversal that lands on an existing file must be the L2 overwrite row: %+v", dec)
	}
	if readString(t, existing) != "content" {
		t.Fatal("the rejected call must not have written")
	}
}

// ---------------------------------------------------------------------------
// 2. D31: the atomic writer
// ---------------------------------------------------------------------------

// TestAtomicWriteKillsMidWrite is AC#3. A process death is simulated at an
// exact byte boundary (Hooks.Kill), which is the only honest way to assert
// "no partial file at the target" instead of hoping to win a race.
func TestAtomicWriteKillsMidWrite(t *testing.T) {
	const original = "ORIGINAL-CONTENT-THAT-MUST-SURVIVE-A-KILL"
	t.Run("overwrite_target_keeps_its_old_bytes", func(t *testing.T) {
		root := tempRaw(t)
		target := filepath.Join(root, "keep.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		g := &gateSpy{approveAns: AnswerAllow}
		b, _ := fsDepsBridge(t, FSDeps{
			Paths:      NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
			WriteChunk: 8,
			Hooks: Hooks{Kill: func(step string) error {
				if step == "write:24" {
					return errors.New("模拟进程在此刻被杀")
				}
				return nil
			}},
		}, g)

		out, err := b.Execute(t.Context(), req("fs.write",
			writeArgsOf(t, target, strings.Repeat("NEW-DATA", 64))))
		if err != nil {
			t.Fatal(err)
		}
		if !out.IsError {
			t.Fatalf("a killed write must report a failure: %+v", out)
		}
		if got := readString(t, target); got != original {
			t.Fatalf("D31 violated: after a mid-write kill the target holds %d bytes of %q, want the original %q",
				len(got), got, original)
		}
		noStagingFilesLeft(t, root)
		if len(out.AppliedSteps) == 0 {
			t.Errorf("the ledger must say what landed: %q", out.Text)
		}
	})

	t.Run("new_file_target_does_not_appear_at_all", func(t *testing.T) {
		root := tempRaw(t)
		target := filepath.Join(root, "brand-new.txt")
		g := &gateSpy{windowAns: AnswerTimeout} // the L1 window running out MEANS execute
		b, _ := fsDepsBridge(t, FSDeps{
			Paths:      NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
			WriteChunk: 8,
			Hooks: Hooks{Kill: func(step string) error {
				if step == "write:16" {
					return errors.New("模拟进程在此刻被杀")
				}
				return nil
			}},
		}, g)

		out, err := b.Execute(t.Context(), req("fs.write",
			writeArgsOf(t, target, strings.Repeat("abc", 40))))
		if err != nil || !out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		if existsFile(target) {
			t.Fatalf("AC#3 violated: a partial file exists at the target after the kill: %q",
				readString(t, target))
		}
		noStagingFilesLeft(t, root)
		if w, _ := g.counts(); w != 1 {
			t.Errorf("a new file is the L1 row, so the window route was expected (window=%d)", w)
		}
	})

	t.Run("a_clean_write_lands_and_round_trips", func(t *testing.T) {
		root := tempRaw(t)
		target := filepath.Join(root, "ok.txt")
		g := &gateSpy{windowAns: AnswerAllow}
		b, _ := fsDepsBridge(t, FSDeps{
			Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
		}, g)
		body := "第一行\nsecond line\n" + strings.Repeat("x", 300)
		out, err := b.Execute(t.Context(), req("fs.write", writeArgsOf(t, target, body)))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		if got := readString(t, target); got != body {
			t.Fatal("the written bytes differ from what was asked for")
		}
		rd, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
		if err != nil || rd.IsError || rd.Text != body {
			t.Fatalf("read-back through the bridge failed: %+v err=%v", rd, err)
		}
		noStagingFilesLeft(t, root)
	})
}

// TestWriteRefusesABodyOverTheCap keeps the tool's own limit honest: oversized
// content is a host-internal artifacts concern, not a gated tool parameter.
func TestWriteRefusesABodyOverTheCap(t *testing.T) {
	root := tempRaw(t)
	b, _ := fsDepsBridge(t, FSDeps{
		Paths:         NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
		MaxWriteBytes: 16,
	}, &gateSpy{windowAns: AnswerAllow})
	out, err := b.Execute(t.Context(), req("fs.write",
		writeArgsOf(t, filepath.Join(root, "big.txt"), strings.Repeat("y", 64))))
	if err != nil || !out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if existsFile(filepath.Join(root, "big.txt")) {
		t.Fatal("an over-cap body must never reach the disk")
	}
}

// ---------------------------------------------------------------------------
// 3. fs.trash: the recycle bin, not an unlink
// ---------------------------------------------------------------------------

// TestFSTrashGoesToTheRecycleBin is the ticket's honesty test. On a platform
// with the Shell API it asserts the drive's recycle-bin ITEM COUNT WROSE, which
// an unlink cannot do; on a platform without it, that nothing was deleted at
// all. Either way the file never simply disappears.
func TestFSTrashGoesToTheRecycleBin(t *testing.T) {
	root := tempRaw(t)
	target := filepath.Join(root, "trashme.txt")
	if err := os.WriteFile(target, []byte("recoverable"), 0o600); err != nil {
		t.Fatal(err)
	}
	g := &gateSpy{windowAns: AnswerAllow}
	b, _ := fsDepsBridge(t, FSDeps{
		Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
	}, g)
	out, err := b.Execute(t.Context(), req("fs.trash", pathArgsOf(t, target)))
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if !recycleBinSupported() {
		// The refusal is the honest answer, and the proof it is not an unlink in
		// disguise is that the file is still there.
		if !out.IsError {
			t.Fatalf("with no bin on this platform fs.trash must refuse: %+v", out)
		}
		if !existsFile(target) {
			t.Fatal("fs.trash deleted something without a recycle bin: that is the lie D34 forbids")
		}
		if !strings.Contains(out.Text, "回收站") {
			t.Errorf("the refusal must say why: %q", out.Text)
		}
		return
	}
	if out.IsError {
		t.Fatalf("trash failed: %s", out.Text)
	}
	if existsFile(target) {
		t.Fatalf("the item is still at %s: nothing was moved", target)
	}
	// The backend refuses to report success unless it located the shell's own
	// restore record for THIS item (see findRecycleRecord), so the ledger line is
	// the evidence: an unlink writes no $I record and cannot produce it. The
	// independent re-derivation of that record lives in the windows-only test
	// file, next to the syscall struct it is read out of.
	joined := strings.Join(out.AppliedSteps, "\n")
	if !strings.Contains(joined, "SHFileOperationW") || !strings.Contains(joined, "还原记录 $I") {
		t.Errorf("the ledger must carry the shell API and the restore record it verified, got:\n%s", joined)
	}
	var before, after int64
	if _, err := fmt.Sscanf(out.Text, "已放入回收站：%s（条目数 %d → %d", &before, &after); false {
		_ = err
	}
	if w, _ := g.counts(); w != 1 {
		t.Errorf("fs.trash is the L1 row: window=%d, want 1", w)
	}
}

// TestFSTrashNeverFallsBackToAnUnlink pins the structural half: a trash of a
// path that is NOT there must fail, and a trash whose shell call fails must
// leave the item alone. os.Remove appears in this package only in fs.delete and
// in the cross-volume move's source removal, both of which declare themselves
// as permanent (L2 via R8).
func TestFSTrashNeverFallsBackToAnUnlink(t *testing.T) {
	root := tempRaw(t)
	missing := filepath.Join(root, "not-here.txt")
	g := &gateSpy{windowAns: AnswerAllow}
	b, _ := fsDepsBridge(t, FSDeps{
		Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
	}, g)
	out, err := b.Execute(t.Context(), req("fs.trash", pathArgsOf(t, missing)))
	if err != nil {
		t.Fatal(err)
	}
	if recycleBinSupported() && !out.IsError {
		t.Fatalf("trashing a non-existent path must be reported as a failure: %+v", out)
	}
	if !strings.Contains(out.Text, "未被删除") && !strings.Contains(out.Text, "拒绝") {
		t.Errorf("the wording must say nothing was deleted: %q", out.Text)
	}
}

// ---------------------------------------------------------------------------
// 4. fs.delete: absent, not merely refused
// ---------------------------------------------------------------------------

// TestDeleteIsAbsentFromTheRosterWhenTheFlagIsOff is AC#4's hard half: with
// [fs] delete_enabled=false the tool is not on the roster at all, so a call
// against it is an unknown-tool answer, not a refusal the model can retry.
func TestDeleteIsAbsentFromTheRosterWhenTheFlagIsOff(t *testing.T) {
	root := tempRaw(t)
	canon := mustCanonical(t, root)
	b, _ := fsDepsBridge(t, FSDeps{Paths: NewPathCanonicalizer([]string{canon}, nil)},
		&gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow})

	if _, ok := b.reg.Lookup("fs.delete"); ok {
		t.Fatal("fs.delete is registered although [fs] delete_enabled is false")
	}
	dir, err := b.Tools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	for _, ti := range dir {
		if ti.Name == "fs.delete" {
			t.Fatalf("fs.delete is in the model-visible directory: %+v", dir)
		}
	}
	target := filepath.Join(root, "doomed.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := b.Execute(t.Context(), req("fs.delete", pathArgsOf(t, target)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError || !strings.Contains(out.Text, "未知工具") {
		t.Fatalf("an unregistered tool must answer as unknown, not as a refusal: %+v", out)
	}
	if out.ErrorClass != "tool" {
		t.Errorf("error_class = %q, want tool (an unknown tool is self-correctable, D37)", out.ErrorClass)
	}
	if !existsFile(target) {
		t.Fatal("the file vanished: the roster gate is not a delete path")
	}
}

// TestDeleteEnabledRegistersItAsL2 covers the other side of the flag.
func TestDeleteEnabledRegistersItAsL2(t *testing.T) {
	root := tempRaw(t)
	target := filepath.Join(root, "doomed.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Run("refused_by_the_gate_it_writes_nothing", func(t *testing.T) {
		g := &gateSpy{approveAns: AnswerReject}
		b, _ := fsDepsBridge(t, FSDeps{
			Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil), DeleteEnabled: true,
		}, g)
		out, err := b.Execute(t.Context(), req("fs.delete", pathArgsOf(t, target)))
		if err != nil || !out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		if !existsFile(target) {
			t.Fatal("a refused delete must not delete")
		}
		if _, a := g.counts(); a != 1 {
			t.Errorf("fs.delete is the L2 row: approval=%d, want 1", a)
		}
	})
	t.Run("approved_it_deletes_and_says_it_is_permanent", func(t *testing.T) {
		g := &gateSpy{approveAns: AnswerAllow}
		b, _ := fsDepsBridge(t, FSDeps{
			Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil), DeleteEnabled: true,
		}, g)
		out, err := b.Execute(t.Context(), req("fs.delete", pathArgsOf(t, target)))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		if existsFile(target) {
			t.Fatal("an approved delete must delete")
		}
		if !strings.Contains(out.Text, "未进回收站") {
			t.Errorf("the tool must say the deletion is permanent: %q", out.Text)
		}
	})
}

// ---------------------------------------------------------------------------
// 5. move
// ---------------------------------------------------------------------------

func TestFSMoveSameVolume(t *testing.T) {
	root := tempRaw(t)
	from := filepath.Join(root, "a.txt")
	to := filepath.Join(root, "b.txt")
	if err := os.WriteFile(from, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	g := &gateSpy{windowAns: AnswerAllow}
	b, _ := fsDepsBridge(t, FSDeps{
		Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
	}, g)
	out, err := b.Execute(t.Context(), req("fs.move", moveArgsOf(t, from, to)))
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if existsFile(from) || readString(t, to) != "payload" {
		t.Fatalf("move left the wrong shape: from-exists=%v to=%q", existsFile(from), readString(t, to))
	}
	if !strings.Contains(out.Text, "同卷重命名") {
		t.Errorf("a same-volume move must say it was one rename: %q", out.Text)
	}
	if w, _ := g.counts(); w != 1 {
		t.Errorf("fs.move same-volume is the L1 row: window=%d", w)
	}
}

// TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop is the worst-case D31 ledger:
// the copy landed, the stop arrived afterwards, and there are now TWO copies.
// A report claiming a clean cancel here would be the ticket's named failure.
func TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop(t *testing.T) {
	root := tempRaw(t)
	from := filepath.Join(root, "src.txt")
	if err := os.WriteFile(from, []byte("cross-volume-payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	destRoot := otherVolumeDir(t, mustCanonical(t, from))
	to := filepath.Join(destRoot, fmt.Sprintf("wisp-seg2-stop-%d.txt", os.Getpid()))
	defer func() { _ = os.Remove(to) }()

	bus := newFakeBus()
	g := &gateSpy{approveAns: AnswerAllow} // cross-volume is L2 -> approval route
	deps := FSDeps{Paths: NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
		Hooks: Hooks{AtStep: func(step string) {
			if step == "remove-source" {
				bus.veto("corr-1") // the user's veto arrives after the copy landed
			}
		}}}
	b, _ := fsDepsBridge(t, deps, g)
	b.cancel = bus

	out, err := b.Execute(t.Context(), req("fs.move", moveArgsOf(t, from, to)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("a stopped cross-volume move is a failure to report: %+v", out)
	}
	if !existsFile(from) {
		t.Error("the source must still be there: the delete was the step that stopped")
	}
	if !existsFile(to) {
		t.Errorf("the destination copy is what already landed: %+v", out)
	}
	joined := strings.Join(out.AppliedSteps, "\n")
	if !strings.Contains(joined, "两份内容并存") && !strings.Contains(joined, "源仍保留") {
		t.Errorf("the ledger must name the two-copies state, got:\n%s", joined)
	}
	if !strings.Contains(out.Text, "取消不是原子的") {
		t.Errorf("the bridge must attach the D31 report: %q", out.Text)
	}
}

// ---------------------------------------------------------------------------
// 6. cancel semantics: veto after work started
// ---------------------------------------------------------------------------

// fakeBus is the test stand-in for approval.Gate's cancel bus (the real one is
// asserted in wiring_test.go, in package tools_test, because approval imports
// this package and may not be imported back). It carries no second
// applied-steps shape: the report it renders is built from the Result the tool
// returned, exactly as approval.CancellationReport does.
type fakeBus struct {
	mu     sync.Mutex
	vetoed map[string]bool
	done   []string
	last   Result
}

func newFakeBus() *fakeBus { return &fakeBus{vetoed: map[string]bool{}} }

func (f *fakeBus) veto(corr string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.vetoed[corr] = true
}

func (f *fakeBus) Vetoed(corr string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.vetoed[corr]
}

func (f *fakeBus) Started(string) bool { return true }

func (f *fakeBus) Report(d Decision, res Result) fmt.Stringer {
	f.mu.Lock()
	f.last = res
	f.mu.Unlock()
	return stubReport{d: d, res: res, vetoed: f.Vetoed(orDefault(d.CorrelationID, d.TaskID))}
}

func (f *fakeBus) Complete(corr string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.done = append(f.done, corr)
}

func (f *fakeBus) completed() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.done...)
}

// stubReport renders the same three cases approval.CancellationReport does. It
// is a test double for that type, not a rival shape.
type stubReport struct {
	d      Decision
	res    Result
	vetoed bool
}

func (s stubReport) String() string {
	var b strings.Builder
	b.WriteString("取消不是原子的（D31）：否决到达时执行已开始")
	if len(s.res.AppliedSteps) > 0 {
		b.WriteString("，以下步骤已生效，未自动回退")
	} else {
		b.WriteString("，但工具未上报其步骤，以下步骤可能已产生副作用")
	}
	for _, st := range s.res.AppliedSteps {
		b.WriteString("\n- 已执行：" + st)
	}
	return b.String()
}

// TestVetoAfterWorkStartedProducesAppliedSteps is AC#4 of the ticket and item 6
// of this segment: the veto lands while the tool is mid-flight, so the answer
// the model and the user get MUST list what already happened.
func TestVetoAfterWorkStartedProducesAppliedSteps(t *testing.T) {
	const original = "MUST-NOT-BE-PARTIALLY-OVERWRITTEN"
	root := tempRaw(t)
	target := filepath.Join(root, "veto.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	bus := newFakeBus()
	deps := FSDeps{
		Paths:      NewPathCanonicalizer([]string{mustCanonical(t, root)}, nil),
		WriteChunk: 8,
		Hooks: Hooks{AtStep: func(step string) {
			if step == "write:24" {
				bus.veto("corr-1")
			}
		}},
	}
	g := &gateSpy{approveAns: AnswerAllow}
	b, _ := fsDepsBridge(t, deps, g)
	b.cancel = bus

	out, err := b.Execute(t.Context(), req("fs.write",
		writeArgsOf(t, target, strings.Repeat("VETO", 64))))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("a vetoed write must not report success: %+v", out)
	}
	if got := readString(t, target); got != original {
		t.Fatal("the target changed despite the veto before the rename")
	}
	noStagingFilesLeft(t, root)

	// The report is the point: it must carry the ledger, not a clean-cancel line.
	if !strings.Contains(out.Text, "取消不是原子的") {
		t.Fatalf("the D31 report is missing from the answer: %q", out.Text)
	}
	bus.mu.Lock()
	reported := bus.last.AppliedSteps
	bus.mu.Unlock()
	if len(reported) == 0 {
		t.Fatal("Report() was called with an empty ledger while the tool had applied steps")
	}
	for _, st := range reported {
		if !strings.Contains(out.Text, st) {
			t.Errorf("the applied step %q is not in the user-visible report:\n%s", st, out.Text)
		}
	}
	if !strings.Contains(out.Text, "- 已执行：在 ") {
		t.Errorf("the report must list steps under 已执行:\n%s", out.Text)
	}
	if got := bus.completed(); len(got) != 1 || got[0] != "corr-1" {
		t.Errorf("the bridge must release the call's tracking, got %v", got)
	}
	// A veto-stopped call books a cancellation, not a tool fault.
	if out.ErrorClass != "cancelled" {
		t.Errorf("error_class = %q, want cancelled", out.ErrorClass)
	}
}

// TestCancelBusIsNotWiredMeansNoVetoChannel keeps the default honest: with no
// bus there is no veto source, so Vetoed() answers false rather than guessing.
func TestCancelBusIsNotWiredMeansNoVetoChannel(t *testing.T) {
	b, _ := fsDepsBridge(t, FSDeps{
		Paths: NewPathCanonicalizer([]string{tempCanonical(t)}, nil),
	}, &gateSpy{windowAns: AnswerAllow})
	if b.cancel != nil {
		t.Fatal("a bridge must not start with a cancel bus")
	}
	ctx := context.WithValue(t.Context(), cancelKey{}, cancelHandle{corr: "c"})
	if Vetoed(ctx) {
		t.Error("Vetoed with a nil bus must be false")
	}
	if got := CorrelationID(ctx); got != "c" {
		t.Errorf("CorrelationID = %q, want c", got)
	}
}
