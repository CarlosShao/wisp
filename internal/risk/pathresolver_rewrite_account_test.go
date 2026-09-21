package risk

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 102 fix (B): the expansion step of SPEC-06 §4 stays (the contract puts
// 「展开(env / ~)」 first in the pipeline on purpose), but every substitution it
// performs is ACCOUNTED for on the Result, and every leg that turns a Result
// into an action reads that account. These tests are the regression net for the
// two ways the fix can rot: the account stops being filled (mutation 94-style:
// Rewritten forced false), or a consumer stops reading it (a future leg that
// grabs Result.Canonical directly and forgets).

func TestC26RewriteAccountIsRecorded(t *testing.T) {
	root := t.TempDir()

	t.Run("percent_env_var", func(t *testing.T) {
		const name = "WISP102ACC_ENV"
		t.Setenv(name, "elsewhere")
		in := filepath.Join(root, "a%"+name+"%b") + sepStr + "artifacts"
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("Resolve(classification leg) must keep working: %v", err)
		}
		requireRewritten(t, res, in, "env")
	})

	t.Run("dollar_env_var", func(t *testing.T) {
		const name = "WISP102ACC_DOLLAR"
		t.Setenv(name, "elsewhere")
		in := filepath.Join(root, "a$"+name+"b")
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		requireRewritten(t, res, in, "env")
	})

	t.Run("leading_home_tilde", func(t *testing.T) {
		in := `~` + sepStr + `wisp102-acc-home` + sepStr + `artifacts`
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if _, err := os.UserHomeDir(); err != nil {
			t.Skipf("no home to expand into on this host: %v", err)
		}
		requireRewritten(t, res, in, "home")
	})

	t.Run("clean_absolute_is_not_marked", func(t *testing.T) {
		in := filepath.Join(root, "plain", "artifacts")
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if res.Rewritten || len(res.Rewrites) != 0 {
			t.Errorf("Rewritten=%v Rewrites=%v for a spelling nothing expanded; a mark that lies by itself is as useless as one stuck at false",
				res.Rewritten, res.Rewrites)
		}
		if res.Spelling != in {
			t.Errorf("Spelling = %q, want the input as handed in %q", res.Spelling, in)
		}
		got, err := res.Actable()
		if err != nil {
			t.Errorf("Actable() refused a path that was never rewritten: %v", err)
		}
		if got != res.Canonical {
			t.Errorf("Actable() = %q, want Canonical %q", got, res.Canonical)
		}
	})

	t.Run("unexpandable_construct_is_not_marked", func(t *testing.T) {
		// The fix must not turn into a blanket ban on '%' (option (A), which the
		// contract forbids): an unknown %CONSTRUCT% passes through untouched and
		// stays actable.
		in := filepath.Join(root, `a%WISP102ACC_UNSET%b`)
		res, err := Resolve(in, nil)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}
		if res.Rewritten {
			t.Errorf("Rewritten=true for %q although nothing expanded (canonical %q)", in, res.Canonical)
		}
		if _, err := res.Actable(); err != nil {
			t.Errorf("Actable() = %v, want the untouched lexical path accepted", err)
		}
	})
}

func requireRewritten(t *testing.T, res Result, in, kind string) {
	t.Helper()
	if !res.Rewritten {
		t.Errorf("Rewritten=false for %q (canonical %q): expansion substituted a value and nobody was told — that is ticket 102's fail-open",
			in, res.Canonical)
	}
	if res.Spelling != in {
		t.Errorf("Spelling = %q, want the caller's input %q", res.Spelling, in)
	}
	found := false
	for _, k := range res.Rewrites {
		if k == kind {
			found = true
		}
	}
	if !found {
		t.Errorf("Rewrites = %v, want it to name %q", res.Rewrites, kind)
	}
	if got, err := res.Actable(); err == nil || got != "" {
		t.Errorf("Actable() = (%q, %v), want ErrRewrittenPath for a rewritten spelling", got, err)
	}
}

// TestC26RewrittenSyncRootDoesNotDisarmSuspectNet is the syncdirs leg (AC#2
// disposition 2). The expanded tree EXISTS on disk, so the handle step resolves
// it: the only thing that can keep this root from counting as canonical-grade
// evidence is the rewrite account. Flip Result.Rewritten to a constant false and
// this test is the one that goes red.
func TestC26RewrittenSyncRootDoesNotDisarmSuspectNet(t *testing.T) {
	elsewhere := t.TempDir()
	const name = "WISP102ACC_SYNCROOT"
	t.Setenv(name, elsewhere)
	resolved := filepath.Join(elsewhere, "wisp102-sub")
	if err := os.MkdirAll(resolved, 0o700); err != nil {
		t.Fatalf("mkdir expanded root: %v", err)
	}
	spelled := "%" + name + "%" + sepStr + "wisp102-sub"

	control := &syncSet{}
	control.add(SyncRoot{Provider: "fixture", Path: resolved, Source: "env"})
	control.finalize()
	if len(control.roots) != 1 || !control.roots[0].canonical {
		t.Skipf("platform gives no handle-resolved form for an existing directory (%+v): the assertion below would prove nothing",
			control.roots)
	}
	if !control.complete {
		t.Fatalf("control set not complete: finalize() ignores canonical, so the rewritten assertion would prove nothing")
	}

	rewritten := &syncSet{}
	rewritten.add(SyncRoot{Provider: "fixture", Path: spelled, Source: "env"})
	rewritten.finalize()
	if len(rewritten.roots) != 1 {
		t.Fatalf("root not added: %+v", rewritten.roots)
	}
	if rewritten.roots[0].canonical {
		t.Errorf("canonical=true for a rewritten sync root (%s): it would disarm the sync-suspect net on a tree the config never spelled",
			rewritten.roots[0].canon)
	}
	if rewritten.complete {
		t.Errorf("complete=true: a rewritten root alone switched off the fail-closed sync-suspect net")
	}
	if rewritten.roots[0].canon != normPath(resolved) && rewritten.roots[0].canon != normPath(control.roots[0].canon) {
		t.Logf("guard still binds to %s (expanded form); the config spelled %s", rewritten.roots[0].canon, spelled)
	}
}

// TestC26RewriteAccountIsConsumedAtEverySecurityLeg is AC#5: the static
// criterion. It does not replace the behavioural tests above - it fails when a
// NEW consumer starts reading Result.Canonical directly, or when one of the
// three security legs stops naming the account it reads.
func TestC26RewriteAccountIsConsumedAtEverySecurityLeg(t *testing.T) {
	// Every leg, and the account it must read out loud.
	legs := []struct {
		file    string
		needles []string
		why     string
	}{
		{"winsec_c26.go", []string{"Actable("}, "sealing/placement: acts on a tree and reports success"},
		{"syncdirs.go", []string{"Actable(", "res.Rewritten"}, "sync roots and the fs.write sync verdict"},
		{filepath.Join("..", "tools", "paths.go"), []string{"res.Rewritten"}, "allowlist roots authorize a tree"},
	}
	for _, leg := range legs {
		src := readSource(t, leg.file)
		for _, n := range leg.needles {
			if !strings.Contains(src, n) {
				t.Errorf("%s no longer reads the ticket 102 account (%q missing): %s", leg.file, n, leg.why)
			}
		}
	}

	// Repo-wide: nothing may read Result.Canonical without reading the account.
	// pathresolver.go is excluded because it WRITES the field.
	producer := "pathresolver.go"
	hits := canonicalReaders(t)
	if len(hits) == 0 {
		t.Fatalf("the scan found no .Canonical reader at all: the instrument is broken, not the code")
	}
	for _, h := range hits {
		if h.file == producer {
			continue
		}
		src := readSource(t, h.file)
		if !strings.Contains(src, "Actable(") && !strings.Contains(src, ".Rewritten") {
			t.Errorf("%s:%d reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer",
				h.file, h.line)
		}
	}
}

type canonRead struct {
	file string
	line int
}

// canonicalReaders walks the module from the repo root and lists every
// non-test Go file that reads a `.Canonical` field.
func canonicalReaders(t *testing.T) []canonRead {
	t.Helper()
	var out []canonRead
	repo := filepath.Join("..", "..")
	err := filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // unreadable corners of a working tree are not findings
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", ".scratch", "node_modules", "dist", "build", "testdata", "design", "docs":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// tools/d22scan carries its own seeded fixtures about unrelated bans.
		if strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(filepath.Join(repo, "tools"))) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for line := 1; sc.Scan(); line++ {
			if readsCanonicalField(sc.Text()) {
				out = append(out, canonRead{file: path, line: line})
			}
		}
		return sc.Err()
	})
	if err != nil {
		t.Fatalf("walk %s: %v", repo, err)
	}
	return out
}

// readsCanonicalField matches a read of the Result.Canonical FIELD and not the
// unrelated PathCanonicalizer.Canonicalize METHOD: the discriminator is the
// character right after ".Canonical", which must not continue the identifier.
func readsCanonicalField(line string) bool {
	for rest := line; ; {
		i := strings.Index(rest, ".Canonical")
		if i < 0 {
			return false
		}
		tail := rest[i+len(".Canonical"):]
		if tail == "" {
			return true
		}
		if c := tail[0]; c != '_' && (c < '0' || c > '9') && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			return true
		}
		rest = tail
	}
}

func readSource(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (the static criterion must not be skippable by moving a file)", path, err)
	}
	return string(b)
}
