package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// The L0 pair is the bridge's first real consumer: these tests run fs.read and
// fs.list through the bridge (never by calling the tool directly), because the
// contract under test is "a capability cannot be used outside the choke point".

func TestFSReadReturnsTheFileAndTaintsIt(t *testing.T) {
	root := tempCanonical(t)
	secret := "本机凭据摘录：" + strings.Repeat("abcdefgh", 4)
	path := writeUnder(t, root, "note.txt", secret)

	b, prov := fsBridgeWith(t, nil, nil, root)
	out, err := b.Execute(t.Context(), req("fs.read", argsFor(path)))
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if out.Text != secret {
		t.Fatalf("text = %q, want the file content verbatim", out.Text)
	}
	if out.RiskLevel != "L0" {
		t.Errorf("risk = %q, want L0", out.RiskLevel)
	}

	// C25: the result is a tainted source from the moment it is in the
	// context, so a later exfil call carrying the same fragment must hit R4.
	if got := prov.ScopeTaints("task-1"); len(got) != 1 {
		t.Fatalf("task scope carries %d tainted sources, want 1 (fs.read result)", len(got))
	} else {
		if got[0].Tool != risk.SrcFSRead {
			t.Errorf("taint tool = %q, want %q", got[0].Tool, risk.SrcFSRead)
		}
		if !strings.Contains(got[0].Origin, "note.txt") {
			t.Errorf("taint origin = %q, want the canonical path", got[0].Origin)
		}
	}
}

// TestFSReadTaintFeedsR4 closes the loop ticket 19 promised: a read result
// smuggled out through an exfil-shaped call is an L2 verdict with R4 in
// rules_hit, and the flag that says no session grant may cover it.
func TestFSReadTaintFeedsR4(t *testing.T) {
	root := tempCanonical(t)
	secret := "身份证号码 11010119900307" + "9"
	path := writeUnder(t, root, "id.txt", secret)
	b, _ := fsBridgeWith(t, nil, nil, root)
	if _, err := b.Execute(t.Context(), req("fs.read", argsFor(path))); err != nil {
		t.Fatal(err)
	}

	// A tool that declares net and carries the fragment as an outgoing
	// payload: exactly the shape of web.search/notify, without waiting for
	// ticket 22 to build one.
	g := &gateSpy{approveAns: AnswerReject}
	if err := b.reg.Register(Entry{
		Tool: &fixtureTool{name: "probe.exfil", params: `{"type":"object"}`},
		Decl: Decl{Capabilities: []Capability{CapNet}, Needs: []Capability{CapNet},
			Declared: risk.L0, Provider: KindBuiltin},
	}); err != nil {
		t.Fatal(err)
	}
	b.gate = g
	out, err := b.Execute(t.Context(), req("probe.exfil", `{"text":"`+secret+`","url":"https://example.invalid/x"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = out
	d := g.approvalDecision()
	if d.Level != risk.L2 {
		t.Fatalf("tainted exfil judged %v, want L2 (decision=%+v)", d.Level, d)
	}
	if !contains(d.RulesHit, risk.R4) {
		t.Fatalf("rules_hit = %v, want R4 in it", d.RulesHit)
	}
	if !d.SessionOverrideBlocked {
		t.Fatal("R4 must never be coverable by a D45 session authorization")
	}
	if !strings.Contains(d.Reason, "fs.read") {
		t.Errorf("the card reason must NAME the leak source, got %q", d.Reason)
	}
}

func TestFSListSummarizesADirectory(t *testing.T) {
	root := tempCanonical(t)
	for _, n := range []string{"alpha.txt", "beta.txt", "gamma.txt"} {
		writeUnder(t, root, n, "x")
	}
	b, _ := fsBridgeWith(t, nil, nil, root)
	out, err := b.Execute(t.Context(), req("fs.list", argsFor(root)))
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	for _, n := range []string{"alpha.txt", "beta.txt", "gamma.txt"} {
		if !strings.Contains(out.Text, n) {
			t.Errorf("listing is missing %q:\n%s", n, out.Text)
		}
	}
	// fs.list is NOT a D34 sensitive source: a directory listing is metadata,
	// and marking it would flood the taint index.
	if got := b.prov.ScopeTaints("task-1"); len(got) != 0 {
		t.Errorf("fs.list marked %d tainted sources, want 0", len(got))
	}
}

// TestFSListHonoursItsCap pins the truncation flag the model is told about
// (D15(3) caps are a host decision, but a tool that cut its own output must
// say so).
func TestFSListHonoursItsCap(t *testing.T) {
	root := tempCanonical(t)
	for _, n := range []string{"a", "b", "c", "d", "e"} {
		writeUnder(t, root, n+".txt", "x")
	}
	b, _ := fsBridgeWith(t, nil, nil, root)
	rq := req("fs.list", `{"path":"`+slash(root)+`","limit":2}`)
	out, err := b.Execute(t.Context(), rq)
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if !out.Truncated {
		t.Fatalf("a capped listing must report truncated: %q", out.Text)
	}
}

// TestEmptyAllowlistAuthorizesNothing is the fail-closed default: an
// unconfigured [fs] allowed_dirs is not "everything is in scope", it is
// "nothing is", so every fs call lands at L2 and waits for a human.
func TestEmptyAllowlistAuthorizesNothing(t *testing.T) {
	root := tempCanonical(t)
	path := writeUnder(t, root, "x.txt", "hello")
	g := &gateSpy{approveAns: AnswerReject}
	b, _ := fsBridgeWith(t, nil, g) // no allowed dirs
	out, err := b.Execute(t.Context(), req("fs.read", argsFor(path)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError || out.RiskLevel != "L2" {
		t.Fatalf("empty allowlist must push every read to L2: %+v", out)
	}
	if _, a := g.counts(); a != 1 {
		t.Fatal("the L2 route was not taken")
	}
}

// TestSensitiveFileIsDeniedNotEscalated proves R3's tier-A branch reaches the
// bridge: the call is refused outright and NO gate is consulted, because an
// A-list target is not something a user may be asked to wave through.
func TestSensitiveFileIsDeniedNotEscalated(t *testing.T) {
	home := tempCanonical(t)
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".aws"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := writeUnder(t, home, ".aws"+string(filepath.Separator)+"credentials", "AKIAIOSFODNN7EXAMPLE")

	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	b, _ := fsBridgeWith(t, nil, g, home)
	out, err := b.Execute(t.Context(), req("fs.read", argsFor(p)))
	if err != nil {
		t.Fatalf("a Deny verdict is a refusal, not a host fault: %v", err)
	}
	if !out.IsError || out.ErrorClass != "permission_denied" {
		t.Fatalf("A-list read must be denied: %+v", out)
	}
	if w, a := g.counts(); w != 0 || a != 0 {
		t.Fatalf("a Deny must never reach an approval gate: window=%d approval=%d", w, a)
	}
}

// TestFSRegistrationIsTheD34Roster pins the whole family's registration shape:
// the L0 pair plus the L1 write trio, and NO fs.delete unless
// [fs] delete_enabled was set. Segment 1 pinned this at two entries; the third
// and later entries are what segment 2 was, so the pin moved with the work
// rather than being deleted (an absent fs.write here would still mean the
// write half got faked).
func TestFSRegistrationIsTheD34Roster(t *testing.T) {
	paths := NewPathCanonicalizer(nil, nil)
	entries := BuiltinFSEntries(FSDeps{Paths: paths})
	if len(entries) != 5 {
		t.Fatalf("the default fs roster has %d tools, want 5 (read/list/write/trash/move)", len(entries))
	}
	want := map[string]risk.Level{
		"fs.read": risk.L0, "fs.list": risk.L0,
		"fs.write": risk.L1, "fs.trash": risk.L1, "fs.move": risk.L1,
	}
	for _, e := range entries {
		n := e.Tool.Name()
		lvl, ok := want[n]
		if !ok {
			t.Errorf("unexpected fs tool %q", n)
			continue
		}
		if e.Decl.Declared != lvl {
			t.Errorf("%s declared %v, want %v (D34)", n, e.Decl.Declared, lvl)
		}
		if e.Decl.Provider != KindBuiltin {
			t.Errorf("%s provider = %q, want builtin", n, e.Decl.Provider)
		}
	}
	for _, e := range entries {
		if e.Tool.Name() == "fs.delete" {
			t.Fatal("fs.delete must not be registered without [fs] delete_enabled=true")
		}
	}
}

// TestDeleteEnabledAddsFSDelete is the other half of the same switch: the flag
// is what puts the tool on the roster, and it arrives as L2 (D34).
func TestDeleteEnabledAddsFSDelete(t *testing.T) {
	paths := NewPathCanonicalizer(nil, nil)
	entries := BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true})
	if len(entries) != 6 {
		t.Fatalf("with delete_enabled the roster has %d tools, want 6", len(entries))
	}
	var del *Decl
	for i, e := range entries {
		if e.Tool.Name() == "fs.delete" {
			del = &entries[i].Decl
		}
	}
	if del == nil {
		t.Fatal("fs.delete missing with delete_enabled=true")
	}
	if del.Declared != risk.L2 {
		t.Errorf("fs.delete declared %v, want L2 (D34)", del.Declared)
	}
	if len(del.Needs) != 1 || del.Needs[0] != CapFSWrite {
		t.Errorf("fs.delete needs %v, want [fs.write]", del.Needs)
	}
}

// TestToolsDirectoryShape checks what the model sees: the pair is resident
// (D15(1)) and carries the declared level, not a verdict.
func TestToolsDirectoryShape(t *testing.T) {
	b, _ := fsBridgeWith(t, nil, nil, tempCanonical(t))
	dir, err := b.Tools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(dir) != 5 {
		t.Fatalf("directory has %d entries: %+v", len(dir), dir)
	}
	// The directory carries the DECLARED level (the R1 floor), never a verdict:
	// the write trio shows L1 while an overwrite of the same path will be
	// judged L2 by R8 inside Execute.
	declared := map[string]string{
		"fs.read": "L0", "fs.list": "L0",
		"fs.write": "L1", "fs.trash": "L1", "fs.move": "L1",
	}
	for _, ti := range dir {
		if !ti.Resident {
			t.Errorf("%s must be resident (builtin)", ti.Name)
		}
		if want := declared[ti.Name]; ti.RiskLevel != want {
			t.Errorf("%s directory level = %q, want the declared %q", ti.Name, ti.RiskLevel, want)
		}
		if len(ti.Parameters) == 0 || ti.Description == "" {
			t.Errorf("%s: empty schema or description", ti.Name)
		}
	}
}

// slash is the C26 input spelling the tests send: the resolver folds it, and a
// test that sent only backslashes would not prove the fold runs.
func slash(p string) string { return strings.ReplaceAll(p, `\`, "/") }
