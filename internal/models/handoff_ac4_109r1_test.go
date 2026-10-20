package models

// Ticket 109 rework leg 109-r1, cell B half 1: R-109-2 turned into a break line of
// AC#4's own nail (the acceptance table's minimal set,
// docs/evidence/s1/109-adversarial-acceptance.md:44 - "把 R-109-2/R-109-3 两条界线的
// 洞补进 AC#4 的用例").
//
// WHICH HOLE THIS IS, AND WHERE IT WAS STILL OPEN. The table's PROBE2 (:110-118)
// measured `VerifyInstalled` returning nil with `attacker.onnx` sitting in the
// install directory and with a whole extra `sub/w.onnx`, because VerifyDir iterated
// entry.InstalledFiles() only (:114). Ticket 121 AC#4 closed that for the hand-off
// verb it reached, and its own legs are TestAC41..TestAC43 in verify_tree_121_test.go.
//
// What stayed open is AC#4's side of the job. Ticket 109's nail
// (TestAC4InstallDirWidthIsDeliberateAndHasABreakLine) answers "这里为什么可以宽、
// 宽到什么程度算破", and all three of its break lines are about *who holds a right*.
// Not one bounds the width in bytes. The table says so in as many words (:116-117):
// "AC#4 三条界量的是'谁持有权利'，一条都没量'这个目录里有多少字节被钉过'". These legs
// put the byte face INTO that nail - the格 the退回 asked for - on the one verb a
// caller cannot step around since cell A landed: Manager.Ensure's own return paths.
//
// WHY THE ENSURE VERB AND NOT VerifyInstalled. `VerifyInstalled` is a verb somebody
// has to remember to call. `Ensure` is the one every caller already uses, and R-109-4's
// complaint was that the guard only travelled with the caller who happened to ask. A
// break line that fires only when somebody thinks of the checking function is not a
// break line, it is a convention.
//
// WHY THESE GO RED BEFORE THIS LEG. The intruder is planted from the Progress hook,
// which runs after each branch's last destination-side check at HEAD - the
// fresh-install path hashes inside staging and then moves (downloader.go:300/:710),
// the cache-hit path verifies before its PhaseDone tick - and before Ensure returns.
// So at HEAD the tree comes back with the plant standing in it. Readings, red names
// and assertion text are registered in the ticket face.

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// ac4Rig is one clean installed fixture plus the parts needed to re-enter it while
// planting bytes: the same manifest, the same store root.
type ac4Rig struct {
	t     *testing.T
	m     *Manifest
	dir   string // <store>/kws-fixture
	store string
}

func newAC4Rig(t *testing.T) *ac4Rig {
	t.Helper()
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("clean install: %v", err)
	}
	return &ac4Rig{t: t, m: m, dir: dir, store: store}
}

// ensurePlanting runs one Ensure of the already-installed fixture while writing the
// given rel-path -> body map into the tree the moment Ensure announces completion,
// i.e. after its last destination-side check and before it returns. It hands back
// what Ensure returned.
func (r *ac4Rig) ensurePlanting(plants map[string]string) (string, error) {
	r.t.Helper()
	fired := false
	hook := func(ev ProgressEvent) {
		if ev.Phase != PhaseDone || ev.ModelID != "kws-fixture" || fired {
			return
		}
		fired = true
		for _, rel := range sortedKeys(plants) {
			p := filepath.Join(r.dir, filepath.FromSlash(rel))
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				r.t.Errorf("plant %s: %v", rel, err)
				continue
			}
			if err := os.WriteFile(p, []byte(plants[rel]), 0o644); err != nil {
				r.t.Errorf("plant %s: %v", rel, err)
			}
		}
		r.t.Logf("planted %d unnamed entr%s inside Ensure's post-verification window: %v",
			len(plants), plural(len(plants) == 1, "y", "ies"), sortedKeys(plants))
	}
	mgr, err := NewManager(Options{
		DataDir: r.store, Manifest: r.m, VerifySignature: true, Attempts: 2, BackoffBase: 0, Progress: hook,
	})
	if err != nil {
		r.t.Fatal(err)
	}
	out, err := mgr.Ensure(context.Background(), "kws-fixture")
	if !fired {
		r.t.Fatalf("the hook never fired, so nothing was planted and this leg proves nothing (err=%v)", err)
	}
	return out, err
}

// sortedKeys keeps the planted list, and therefore the log line, stable across runs.
func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TestAC41NailByteFaceUnnamedFileBreaksTheWidth is AC#4's byte break line, first
// shape: the file PROBE2 named, placed by a writer holding the inherited grant the
// nail's first界 expressly ALLOWS, with every named byte of the tree still true.
func TestAC41NailByteFaceUnnamedFileBreaksTheWidth(t *testing.T) {
	r := newAC4Rig(t)
	if err := verifyFileHash(filepath.Join(r.dir, "model.onnx"), tinyArchiveMembers[0].SHA256, tinyArchiveMembers[0].SizeBytes); err != nil {
		t.Fatalf("the fixture's named bytes are already wrong, so this leg measures nothing: %v", err)
	}

	_, err := r.ensurePlanting(map[string]string{"attacker.onnx": "bytes no gate ever read"})
	if err == nil {
		t.Fatalf("AC#4/R-109-2 byte face: Ensure handed back a tree holding attacker.onnx as available. The nail's " +
			"first界 sanctions the inherited write that placed it, so the width has to stop somewhere in bytes too; " +
			"model.onnx still hashes exactly as the signed manifest pins it, which is precisely why hashing the named " +
			"subset cannot answer 'how wide is this'.")
	}
	if !strings.Contains(err.Error(), "attacker.onnx") {
		t.Errorf("the refusal must name the offender or the operator cannot act on it, got: %v", err)
	}
	t.Logf("byte break line (file) holds: %v", err)
}

// TestAC42NailByteFaceExtraDirectoryBreaksTheWidth is the second shape PROBE2
// measured: an entire extra directory.
func TestAC42NailByteFaceExtraDirectoryBreaksTheWidth(t *testing.T) {
	r := newAC4Rig(t)
	_, err := r.ensurePlanting(map[string]string{"sub/w.onnx": "a whole shelf nobody pinned"})
	if err == nil {
		t.Fatal("AC#4/R-109-2 byte face: Ensure handed back a tree holding an entire extra directory (sub/w.onnx) - " +
			"the second PROBE2 shape, and the reason 'verified' has to mean the whole tree rather than the named subset")
	}
	if !strings.Contains(err.Error(), "sub") || !strings.Contains(err.Error(), "w.onnx") {
		t.Errorf("the refusal must name the directory and the byte reached through it, got: %v", err)
	}
	t.Logf("byte break line (directory) holds: %v", err)
}

// TestAC43NailByteFaceCountIsExact is the third leg of the same界, and the one that
// stops the answer from being "refuse the first stray entry and call the rest clean":
// "宽到什么程度算破" is only measurable if the verdict carries the number.
func TestAC43NailByteFaceCountIsExact(t *testing.T) {
	r := newAC4Rig(t)
	_, err := r.ensurePlanting(map[string]string{
		"a1.onnx":      "x",
		"a2.onnx":      "x",
		"a3.onnx":      "x",
		"sub/deep.bin": "x",
	})
	if err == nil {
		t.Fatal("AC#4/R-109-2 byte face: three loose unnamed files plus an unnamed directory and its byte verified clean through Ensure")
	}
	text := err.Error()
	for _, want := range []string{"a1.onnx", "a2.onnx", "a3.onnx", "sub", "deep.bin"} {
		if !strings.Contains(text, want) {
			t.Errorf("the refusal does not name %q: %v", want, err)
		}
	}
	// 5 entries: three loose files, sub/, sub/deep.bin.
	if !strings.Contains(text, "5 entr") {
		t.Errorf("AC#4 byte face: the verdict must carry the exact count (5 here) because that number is what the "+
			"break line is measured in, got: %v", err)
	}
	t.Logf("byte break line (exact count) holds: %v", err)
}

// TestAC44NailByteFaceCleanTreeStillPasses is the reverse leg - AC#2's discipline
// applied to AC#4: a break line that always fires measures nothing. The same tree
// with nothing planted must still hand back, named subdirectory and all.
func TestAC44NailByteFaceCleanTreeStillPasses(t *testing.T) {
	r := newAC4Rig(t)
	if _, err := os.Stat(filepath.Join(r.dir, "dict", "inner.txt")); err != nil {
		t.Fatalf("the fixture did not install the named subdirectory, so this leg never tested it: %v", err)
	}
	if _, err := r.ensurePlanting(map[string]string{}); err != nil {
		t.Fatalf("AC#4 reverse leg: the byte face refused a tree it was handed with nothing planted: %v", err)
	}
	mgr, err := NewManager(Options{DataDir: r.store, Manifest: r.m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("VerifyInstalled on the same clean tree: %v", err)
	}
	t.Logf("clean tree still passes the byte face (named dir dict/ intact): %s", r.dir)
}
