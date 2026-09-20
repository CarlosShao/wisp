package tools

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// shared test wiring. Every test here goes through the SAME composition the
// future cmd/wisp wiring will use (ticket 12): a real C26 canonicalizer over a
// real [fs] allowed_dirs list, the real C19 assessor, the real C25 engine and a
// real SQLite store for the rows. Nothing is stubbed except the gate, which is
// ticket 21's to build.

// fsBridge returns a bridge over the real fs pair rooted at allowed (canonical,
// already inside the allowlist).
func fsBridge(t *testing.T, allowed string) (*Bridge, *PathCanonicalizer) {
	t.Helper()
	b, _ := fsBridgeWith(t, nil, nil, allowed)
	return b, b.Paths()
}

// fsBridgeWith is fsBridge with an optional journal and gate. The C25 engine is
// always wired, because "results are taint-marked" is a hard rule, not a flag.
func fsBridgeWith(t *testing.T, j agent.Journal, g Gate, allowed ...string) (*Bridge, *risk.Provenance) {
	t.Helper()
	if g == nil {
		g = NoGate{}
	}
	var paths *PathCanonicalizer
	if len(allowed) == 0 {
		paths = NewPathCanonicalizer(nil, nil)
	} else {
		paths = NewPathCanonicalizer(allowed, nil)
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths, Provenance: prov, Journal: j, Gate: g,
		Logf: func(string, ...any) {},
	}), prov
}

// openStore opens a real wisp.db under dir.
func openStore(t *testing.T, dir string) *memory.Store {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := memory.Open(dir, memory.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// mustStartTask books the task_log row the tool_call rows hang off.
func mustStartTask(t *testing.T, s *memory.Store, id string) {
	t.Helper()
	if err := s.StartTaskLog(context.Background(), memory.TaskLog{
		ID: id, State: "running", QueryText: "ticket 20 bridge test",
	}); err != nil {
		t.Fatalf("StartTaskLog: %v", err)
	}
}

// tempCanonical gives a canonicalized temp dir: tests compare against the form
// the resolver produces, never against what os gave them, so a case or 8.3
// difference cannot make a verdict test flaky.
func tempCanonical(t *testing.T) string {
	t.Helper()
	return mustCanonical(t, t.TempDir())
}

// mustCanonical runs the frozen C26 pipeline on one path.
func mustCanonical(t *testing.T, p string) string {
	t.Helper()
	c, err := NewPathCanonicalizer(nil, nil).Canonicalize(p)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", p, err)
	}
	return c
}

// writeUnder puts a file inside dir and returns its canonical path.
func writeUnder(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return mustCanonical(t, p)
}

// argsFor renders one fs call's arguments. The path is sent in forward-slash
// form on purpose: C26 must normalize it, so a test that relied on the OS
// spelling would be testing the wrong thing.
func argsFor(path string) string {
	b, _ := json.Marshal(map[string]any{"path": filepath.ToSlash(path)})
	return string(b)
}
