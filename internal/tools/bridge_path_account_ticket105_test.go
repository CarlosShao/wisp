package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 105, AC#2 — ticket 102's rewrite account gets its first PRODUCTION
// reader.
//
// What ticket 102 left behind: RewrittenRoots()/Roots()/UnusableRoots() had
// zero non-test call sites, so "the operator can see which root was moved" was
// a comment rather than a behaviour. This file does not test a function; it
// runs ONE REAL fs.read through a composed Bridge (the same composition
// cmd/wisp/run.go builds: NewPathCanonicalizer -> tools.New(Options{Paths,
// Logf})) with the audit sink writing to an actual file, then reads that file
// back off disk and asserts the CONTENT of the record: the config spelling,
// the tree expansion moved it onto, the construct that did it, and how many
// roots are authorizing right now.
//
// Mutation for the acceptance side: delete the b.log("tools: PATH-ACCOUNT ...")
// block in Bridge.book(). Every case below goes red on a missing/garbled
// record, not on a symbol.

// fileAudit is the production sink shape: cmd/wisp's agentRuntime.auditf writes
// "[audit] " + format to its writer. Here the writer is a file on disk, so the
// assertion is about what an operator would find in the log.
func fileAudit(t *testing.T) (path string, logf func(string, ...any), readBack func() string) {
	t.Helper()
	path = filepath.Join(t.TempDir(), "wisp-audit.log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	logf = func(format string, args ...any) {
		fmt.Fprintf(f, "[audit] "+format+"\n", args...)
	}
	return path, logf, func() string {
		if err := f.Sync(); err != nil {
			t.Fatalf("sync audit file: %v", err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read back the audit file %s: %v", path, err)
		}
		return string(data)
	}
}

// accountBridge composes the real fs pair over `allowed` with the audit sink on
// disk, mirroring cmd/wisp/run.go's Options.
func accountBridge(t *testing.T, allowed []string, logf func(string, ...any)) *Bridge {
	t.Helper()
	pc := NewPathCanonicalizer(allowed, nil)
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: pc}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry:   reg,
		Paths:      pc,
		Gate:       NoGate{},
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       logf,
	})
}

// auditLine returns the single line of `log` carrying `marker`, or "".
func auditLine(log, marker string) string {
	for _, l := range strings.Split(log, "\n") {
		if strings.Contains(l, marker) {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

// TestRealToolCallWritesRewriteAccountIntoAudit is AC#2's end-to-end half: one
// fs.read, and the account of ticket 102's book appears in the audit file with
// its content intact.
func TestRealToolCallWritesRewriteAccountIntoAudit(t *testing.T) {
	dir := t.TempDir()
	tree := filepath.Join(dir, "proj")
	if err := os.MkdirAll(tree, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(tree, "a.txt")
	if err := os.WriteFile(target, []byte("accounted"), 0o600); err != nil {
		t.Fatal(err)
	}

	const name = "WISP105_ACCOUNT_ROOT"
	t.Setenv(name, dir)
	spelled := "%" + name + "%" + pathSep + "proj"

	path, logf, readBack := fileAudit(t)
	b := accountBridge(t, []string{spelled}, logf)

	out, err := b.Execute(t.Context(), req("fs.read", argsFor(target)))
	if err != nil {
		t.Fatalf("Execute(fs.read): %v", err)
	}
	if out.IsError {
		t.Fatalf("precondition: the read must actually run, got %q", out.Text)
	}
	if !strings.Contains(out.Text, "accounted") {
		t.Fatalf("precondition: fs.read returned %q, want the file's bytes", out.Text)
	}

	log := readBack()
	line := auditLine(log, "tools: PATH-ACCOUNT")
	if line == "" {
		t.Fatalf("no PATH-ACCOUNT record in the audit file %s after a real fs.read; "+
			"ticket 102's account has a producer but nobody reads it out to the operator:\n%s", path, log)
	}
	for _, want := range []struct{ needle, meaning string }{
		{"%WISP105_ACCOUNT_ROOT%", "the config spelling that got moved"},
		{tree, "the tree expansion actually authorized"},
		{"(env)", "the construct that substituted"},
		{"roots=1", "how many roots are authorizing right now"},
	} {
		if !strings.Contains(line, want.needle) {
			t.Errorf("PATH-ACCOUNT record = %q, want it to carry %s (%q)", line, want.meaning, want.needle)
		}
	}

	// The record is the account, not an error report: the call it belongs to is
	// still booked, and the two lines name the same tool.
	if auditLine(log, "tools: call ") == "" {
		t.Errorf("the per-call record vanished next to the new one:\n%s", log)
	}
}

// TestUnusableRootIsVisibleInTheAuditRecord covers the other half of the book:
// an [fs] allowed_dirs entry that C26 expanded onto a tree it cannot confirm
// authorizes nothing, and the operator has to be able to read WHY.
func TestUnusableRootIsVisibleInTheAuditRecord(t *testing.T) {
	dir := t.TempDir()
	const name = "WISP105_GHOST_ROOT"
	t.Setenv(name, filepath.Join(dir, "not-there"))
	spelled := "%" + name + "%" + pathSep + "proj"

	path, logf, readBack := fileAudit(t)
	b := accountBridge(t, []string{spelled}, logf)

	// A real tool call on a path outside every root: refused or not, the call is
	// booked, and that is what carries the record.
	if _, err := b.Execute(t.Context(), req("fs.read", argsFor(filepath.Join(dir, "x.txt")))); err != nil {
		t.Fatalf("Execute(fs.read): %v", err)
	}
	line := auditLine(readBack(), "tools: PATH-ACCOUNT")
	if line == "" {
		t.Fatalf("no PATH-ACCOUNT record in %s:\n%s", path, readBack())
	}
	if !strings.Contains(line, "not confirmed on disk") || !strings.Contains(line, "%"+name+"%") {
		t.Errorf("unusable half of the account missing from the record: %q (want %q and the drop reason)",
			line, spelled)
	}
	if !strings.Contains(line, "roots=0") {
		t.Errorf("record = %q, want roots=0: a dropped root must not read as authorization", line)
	}
}

// TestAccountRecordIsWrittenEvenWithNothingToReport is the anti-"only when it
// complains" case: a plainly spelled allowlist has an empty account, and the
// record must still be there so a reader can tell "read and clean" apart from
// "never read".
func TestAccountRecordIsWrittenEvenWithNothingToReport(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(target, []byte("plain"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, logf, readBack := fileAudit(t)
	b := accountBridge(t, []string{dir}, logf)
	if _, err := b.Execute(t.Context(), req("fs.read", argsFor(target))); err != nil {
		t.Fatalf("Execute(fs.read): %v", err)
	}
	line := auditLine(readBack(), "tools: PATH-ACCOUNT")
	if line == "" {
		t.Fatalf("a clean account wrote no record, so an operator cannot tell it apart from an unwired one:\n%s", readBack())
	}
	if !strings.Contains(line, "roots=1") || !strings.Contains(line, "rewritten=[]") ||
		!strings.Contains(line, "unusable=[]") {
		t.Errorf("record = %q, want roots=1 with both account lists present and empty", line)
	}
}
