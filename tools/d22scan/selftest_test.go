package main

// Ticket 161 AC#2 / AC#4: the self-test's own gate.
//
// Three things are asserted here, and the second is the reason the first is not
// decoration:
//   1. every case in selfCases shows what its line claims (runSelfTest exits 0);
//   2. the roster audit has teeth - a table with one violation sample deleted is
//      REFUSED, measured rather than promised;
//   3. the -self-test flag is wired to os.Exit in the built binary, so "the CI
//      step runs it" is a process reading and not a function call a future main()
//      can quietly drop (the same reasoning TestBuiltBinaryGoesRedEndToEnd gives
//      for verdict()).
//
// No test here calls t.Skip, and none of them touches a ban's matcher: they feed
// files to the existing walks and read what came back.

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// selfTestOutput runs the entry point the CI step runs, in-process.
func selfTestOutput(t *testing.T) (string, int) {
	t.Helper()
	var buf bytes.Buffer
	rc := runSelfTest(&buf, &buf, selfBansSourceFlag)
	return buf.String(), rc
}

// TestSelfTestEntryPassesEveryCase is AC#2's "both directions" reading, inside the
// step CI already runs (tools/d22scan/runtests.sh via the lint job's positive
// control), so the self-test has two independent ways to be executed.
func TestSelfTestEntryPassesEveryCase(t *testing.T) {
	out, rc := selfTestOutput(t)
	if rc != 0 {
		t.Fatalf("-self-test exited %d, want 0. Output:\n%s", rc, out)
	}
	// A case line is the only line carrying both a verdict column and the " | "
	// separator before its evidence, so counting them counts CASES, not notes. A
	// mismatch means a case stopped printing at all, which is how an instrument
	// disappears without anyone deleting it.
	counted := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, " OK ") && strings.Contains(line, " | ") {
			counted++
		}
	}
	if counted != len(selfCases) {
		t.Errorf("printed %d case lines for %d cases - one case produced no verdict at all:\n%s", counted, len(selfCases), out)
	}
	if !strings.Contains(out, "clean - all") {
		t.Errorf("missing the clean line:\n%s", out)
	}
}

// TestSelfTestRosterAuditRejectsAHollowTable is AC#4's in-CI reading: it deletes
// ingredients from a COPY of the table and requires the audit to say so. Without
// this test, "removing a sample turns CI red" is only a sentence in a comment.
func TestSelfTestRosterAuditRejectsAHollowTable(t *testing.T) {
	numbered, emitted, err := readRosters(selfBansSourceFlag)
	if err != nil {
		t.Fatalf("reading this tool's own roster: %v", err)
	}
	if len(numbered) == 0 || len(emitted) == 0 {
		t.Fatalf("empty roster read from %s (numbered=%d emitted=%d): the self-test would be auditing nothing", selfBansSourceFlag, len(numbered), len(emitted))
	}
	// Positive control first: the shipped table must audit clean.
	if holes := auditSelfCases(selfCases, numbered, emitted); len(holes) != 0 {
		t.Fatalf("the shipped table must audit clean, got %d hole(s):\n  %s", len(holes), strings.Join(holes, "\n  "))
	}

	without := func(drop func(selfCase) bool) []selfCase {
		out := make([]selfCase, 0, len(selfCases))
		for _, c := range selfCases {
			if !drop(c) {
				out = append(out, c)
			}
		}
		return out
	}

	// (a) the violating sample of one tag is deleted: exactly AC#4's question.
	// The tag is taken from the table rather than hardcoded, so this stays a real
	// test if a ban is ever added or renumbered.
	target := selfCases[0].tag
	joined := ""
	holes := auditSelfCases(without(func(c selfCase) bool { return c.tag == target && c.want == wantRing }), numbered, emitted)
	joined = strings.Join(holes, "\n  ")
	if len(holes) == 0 {
		t.Fatalf("deleting every expect-ring sample of %q was accepted: the audit is decoration", target)
	}
	if !strings.Contains(joined, target) {
		t.Errorf("the report must name the tag whose sample went missing (%q), got:\n  %s", target, joined)
	}
	if !strings.Contains(joined, "only an expect-silent sample") {
		t.Errorf("expected the one-direction wording for %q, got:\n  %s", target, joined)
	}

	// (b) the whole pair of one tag is deleted.
	holes = auditSelfCases(without(func(c selfCase) bool { return c.tag == "mirror-hash" }), numbered, emitted)
	if !strings.Contains(strings.Join(holes, "\n"), "NO sample at all") {
		t.Errorf("deleting both samples of mirror-hash must read as no sample at all, got:\n  %s", strings.Join(holes, "\n  "))
	}

	// (c) a sample for a tag the tool cannot emit - a rename that only touched one
	// side - must be refused too, or the table would keep "testing" a dead ban.
	ghost := append(append([]selfCase{}, selfCases...), selfCase{
		tag: "ban-that-was-removed", want: wantRing, file: "internal/probe/ghost.go",
		src: "package probe\n", summary: "negative control",
	})
	holes = auditSelfCases(ghost, numbered, emitted)
	if !strings.Contains(strings.Join(holes, "\n"), "ban-that-was-removed") {
		t.Errorf("a case naming a tag with no emission site must be reported, got:\n  %s", strings.Join(holes, "\n  "))
	}

	// (d) duplicated samples are not coverage of anything.
	dup := append(append([]selfCase{}, selfCases...), selfCases[0])
	holes = auditSelfCases(dup, numbered, emitted)
	if !strings.Contains(strings.Join(holes, "\n"), "same path with the same bytes") {
		t.Errorf("seeding the same bytes twice must be reported as one sample testing nothing new, got:\n  %s", strings.Join(holes, "\n  "))
	}

	// (e) a roster with a gap in the numbering (a deleted `//\t<N>` line) must not
	// inherit the remaining numbers silently - same lesson AGENTS.md writes for the
	// D-table ("47 枚、无缺号").
	gapped := []banEntry{{num: 1, tag: "bare-goroutine"}, {num: 3, tag: "plaintext-key"}}
	holes = auditSelfCases(selfCases, gapped, emitted)
	if !strings.Contains(strings.Join(holes, "\n"), "not contiguous") {
		t.Errorf("a gap in the ban numbering must be reported, got:\n  %s", strings.Join(holes, "\n  "))
	}
}

// TestSelfTestRefusesToAuditAnEmptyRoster: the entry point's authority comes from
// reading main.go. Point it at anything else and it must say so loudly (exit 2)
// and never "pass" against a roster of zero - the empty-instrument rule ticket 71
// AC#4 applies to the walks, applied here to the self-test itself.
func TestSelfTestRefusesToAuditAnEmptyRoster(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		name string
		body string
	}{
		{"not-a-roster", "package main\n\nfunc main() {}\n"},
		{"block-without-numbers", "// Bans (fixture)\n//\n// nothing numbered here\npackage main\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := filepath.Join(dir, c.name+".go")
			if err := os.WriteFile(path, []byte(c.body), 0o644); err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if rc := runSelfTest(&buf, &buf, path); rc != 2 {
				t.Fatalf("rc=%d want 2 for a source with no usable roster:\n%s", rc, buf.String())
			}
			if !strings.Contains(buf.String(), "FATAL") {
				t.Errorf("a roster failure must read as FATAL, got:\n%s", buf.String())
			}
		})
	}
	var buf bytes.Buffer
	if rc := runSelfTest(&buf, &buf, filepath.Join(dir, "does-not-exist.go")); rc != 2 {
		t.Fatalf("rc=%d want 2 when the roster source cannot be read at all", rc)
	}
}

// TestSelfTestFlagIsWiredInTheBuiltBinary is the process-level reading of the CI
// step. It also proves the mode is independent of -root: a nonsense -root must
// neither fail the self-test (it scans no tree) nor let it "pass" by scanning one.
func TestSelfTestFlagIsWiredInTheBuiltBinary(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("no go toolchain to build the scanner with: %v", err)
	}
	bin := filepath.Join(t.TempDir(), "d22scan-selftest"+exeSuffix())
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
	}
	for _, args := range [][]string{
		{"-self-test"},
		{"-self-test", "-root", filepath.Join("definitely", "not", "a", "wisp", "root")},
	} {
		cmd := exec.Command(bin, args...)
		cmd.Dir = "."
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		rc := 0
		if err != nil {
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				t.Fatalf("running %v failed: %v", args, err)
			}
			rc = ee.ExitCode()
		}
		joined := stdout.String() + stderr.String()
		if rc != 0 {
			t.Fatalf("rc=%d want 0 for %v\nstdout:\n%s\nstderr:\n%s", rc, args, stdout.String(), stderr.String())
		}
		if !strings.Contains(joined, "d22scan -self-test: clean - all") {
			t.Errorf("%v printed no self-test clean line:\n%s", args, joined)
		}
	}
}
