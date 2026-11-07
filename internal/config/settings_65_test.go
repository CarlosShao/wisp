package config

// Ticket 65's settings-write tests: the guarded [ball].style setter and the
// reader face the style picker reads. The schema side of this key (vocabulary,
// enum refusal, glass default) is pinned in ball_style_65_test.go and is not
// re-tested here; these cover the WRITE path the settings page rides.

import (
	"os"
	"strings"
	"testing"
)

// TestSetBallStyleWritesThroughGuardedPath: a jarvis write lands on disk under
// the owned key path, and the receipt names ball.style (writeOneKey's receipt
// is the key paths the FILE changed by, not a summary).
func TestSetBallStyleWritesThroughGuardedPath(t *testing.T) {
	m, path, _ := newTestManager(t, matrixBase)
	written, err := m.SetBallStyle("jarvis")
	if err != nil {
		t.Fatalf("SetBallStyle(jarvis): %v", err)
	}
	if len(written) == 0 || written[0] != "ball.style" {
		t.Fatalf("receipt = %v, want [ball.style]", written)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// The guarded writer re-lays the file in canonical form (single-quoted
	// strings - see the ticket-226 merge note); the key line is what the next
	// start will load.
	if !strings.Contains(string(b), `style = 'jarvis'`) {
		t.Fatalf("config.toml does not carry style = %q after the write:\n%s", "jarvis", b)
	}
	// The read face answers what the process would act on: the written family.
	if v := m.BallStyleOptions(); v.Style != "jarvis" {
		t.Fatalf("BallStyleOptions().Style = %q, want jarvis (in-memory view)", v.Style)
	}
}

// TestSetBallStyleRefusesOutsideVocabulary: a word the schema does not carry is
// refused BEFORE the write, with the field name and the vocabulary in the
// message, and config.toml stays byte-identical (the guarded path's own rule).
func TestSetBallStyleRefusesOutsideVocabulary(t *testing.T) {
	m, path, _ := newTestManager(t, matrixBase)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"aurora", "GLASS", "Jarvis", ""} {
		written, err := m.SetBallStyle(bad)
		if err == nil {
			t.Fatalf("SetBallStyle(%q) must be refused", bad)
		}
		if !strings.Contains(err.Error(), "ball.style") || !strings.Contains(err.Error(), "glass|jarvis") {
			t.Fatalf("the refusal must name the key and the vocabulary, got: %v", err)
		}
		if written != nil {
			t.Fatalf("a refused write must return no receipt, got %v", written)
		}
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("a refused write must leave config.toml byte-identical")
	}
}

// TestSetBallStyleTrimsEnvelopeSpace: the envelope's surrounding space is
// trimmed (the one normalization), and the trimmed word is what lands - while
// a case-folded word ("Jarvis") is refused, because the load-time enum is exact
// lowercase and a setter that accepted it would write a file the next start
// refuses (the trim comment's whole point).
func TestSetBallStyleTrimsEnvelopeSpace(t *testing.T) {
	m, _, _ := newTestManager(t, matrixBase)
	if _, err := m.SetBallStyle("  jarvis  "); err != nil {
		t.Fatalf("surrounding space must be trimmed: %v", err)
	}
	if _, err := m.SetBallStyle("Jarvis"); err == nil {
		t.Fatalf("a case-folded word must be refused: the load enum is exact lowercase")
	}
}

// TestBallStyleOptionsReader: the reader answers the schema vocabulary in
// schema order, and the glass default flows from the schema default (no
// literal here) - a config without a style line reads as glass.
func TestBallStyleOptionsReader(t *testing.T) {
	m, _, _ := newTestManager(t, matrixBase)
	v := m.BallStyleOptions()
	if v.Style != "glass" {
		t.Fatalf("default read = %q, want glass (the schema default, not a local literal)", v.Style)
	}
	if got, want := strings.Join(v.Vocabulary, "|"), "glass|jarvis"; got != want {
		t.Fatalf("vocabulary = %q, want %q (the schema's own list, in order)", got, want)
	}
	// A nil Manager answers an empty style with the vocabulary intact - a
	// caller answers "not configured", never a made-up family.
	var nilManager *Manager
	nv := nilManager.BallStyleOptions()
	if nv.Style != "" || strings.Join(nv.Vocabulary, "|") != "glass|jarvis" {
		t.Fatalf("nil Manager read = %+v, want empty style with the vocabulary intact", nv)
	}
}
