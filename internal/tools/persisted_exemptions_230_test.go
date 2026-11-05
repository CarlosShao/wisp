package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 230 AC#4 (the tools-side seam): Options.PersistedExemptions is a
// SECOND B-tier override source beside Options.Confirmations. It is read by the
// same single production call site (mode.go's readBlacklist) and handed to
// risk.GateStanding as its own map, so the security line can name which source
// allowed the use.
//
// What this test measures, one reading per run:
//   - with no source at all, the note says "unlockable by one L2 confirm";
//   - with ONLY the persisted list armed, the note says already unlocked, and
//     exactly one security line says source=config:risk.blacklist_overrides
//     (which is only possible because the two sources are not merged);
//   - with ONLY a confirmation on record, the same note and a line saying
//     source=l2-confirm (the pre-existing source is untouched);
//   - in every case the routing still asks once: an override record changes the
//     note, never the card (R3 stays silence-immune, ticket 90's AC#5 half).
func TestPersistedExemptionsAreASecondSource(t *testing.T) {
	dir := t.TempDir()
	pem := filepath.Join(dir, "server.pem")
	if err := os.WriteFile(pem, []byte("material"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	canon := t90Canonical(t, pem)
	root := t90Canonical(t, dir)
	key := risk.OverrideKey(canon)

	var logs []string
	prev := risk.Logf
	risk.Logf = func(format string, args ...any) {
		logs = append(logs, fmt.Sprintf(format, args...))
	}
	t.Cleanup(func() { risk.Logf = prev })

	run := func(confirmations, persisted map[string]bool) (BlacklistNote, int, int) {
		gate := &t90Gate{}
		e, _ := t90Entry("t230.b", risk.L0, []string{"path"}, nil)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		b := New(Options{
			Registry: reg, Gate: gate,
			Paths:               NewPathCanonicalizer([]string{root}, nil),
			Confirmations:       func() map[string]bool { return confirmations },
			PersistedExemptions: func() map[string]bool { return persisted },
		})
		if _, err := b.Execute(context.Background(), t90Req("t230.b", pem)); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		w, a := gate.counts()
		return gate.last().Blacklist, w, a
	}

	// 1. No source at all: the fail-closed default is unchanged.
	logs = nil
	note, w, a := run(nil, nil)
	if note.Class != risk.ClassB || len(note.Unlockable) != 1 || len(note.AlreadyUnlocked) != 0 {
		t.Fatalf("with no source the note changed: %+v", note)
	}
	if w+a != 1 {
		t.Errorf("with no source the B-tier path asked %d windows + %d approvals, want exactly one", w, a)
	}

	// 2. The persisted list ALONE unlocks the note, and the line says so.
	logs = nil
	note, w, a = run(nil, map[string]bool{key: true})
	if note.Class != risk.ClassB || len(note.AlreadyUnlocked) != 1 || len(note.Unlockable) != 0 {
		t.Fatalf("a persisted exemption did not reach the note (source not wired): %+v", note)
	}
	overrideLines := lines230(logs, "B-list OVERRIDE")
	if len(overrideLines) != 1 {
		t.Fatalf("a persisted exemption emitted %d override lines, want exactly 1: %q", len(overrideLines), logs)
	}
	if !strings.Contains(overrideLines[0], "source=config:risk.blacklist_overrides") {
		t.Errorf("the override line does not name the config source (sources were merged?): %q", overrideLines[0])
	}
	if strings.Contains(overrideLines[0], "l2-confirm") {
		t.Errorf("a config-sourced allow claims an L2 confirm: %q", overrideLines[0])
	}
	if w+a != 1 {
		t.Errorf("a persisted exemption changed the ROUTING (%d/%d): the note is evidence, "+
			"not an authority, and only an L2 answer may silence R3", w, a)
	}

	// 3. The pre-existing source still works and still says l2-confirm.
	logs = nil
	note, _, _ = run(map[string]bool{key: true}, nil)
	if len(note.AlreadyUnlocked) != 1 {
		t.Fatalf("a confirmation on record stopped unlocking the note: %+v", note)
	}
	overrideLines = lines230(logs, "B-list OVERRIDE")
	if len(overrideLines) != 1 || !strings.Contains(overrideLines[0], "source=l2-confirm") {
		t.Fatalf("the confirmation source must emit exactly one line naming source=l2-confirm, got %q", logs)
	}

	// 4. Both sources armed: still one line, naming the human-answered one.
	logs = nil
	note, _, _ = run(map[string]bool{key: true}, map[string]bool{key: true})
	if len(note.AlreadyUnlocked) != 1 {
		t.Fatalf("both sources armed must still report the one unlocked path: %+v", note)
	}
	overrideLines = lines230(logs, "B-list OVERRIDE")
	if len(overrideLines) != 1 || !strings.Contains(overrideLines[0], "source=l2-confirm") {
		t.Fatalf("with both sources armed, want exactly one line naming l2-confirm, got %q", logs)
	}
}

// lines230 keeps the log lines that mention one tag (so a test can count the
// override lines without depending on the order of any other line).
func lines230(lines []string, tag string) []string {
	var out []string
	for _, l := range lines {
		if strings.Contains(l, tag) {
			out = append(out, l)
		}
	}
	return out
}
