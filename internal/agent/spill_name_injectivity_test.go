package agent

// Ticket 79 — the artifact NAME must be injective in the tool-call id, and the
// write that lands it must be exclusive.
//
// The defect: artifactName kept only [A-Za-z0-9_-], so the genuinely different
// model-supplied ids `p/q`, `p\q` and `pq` all became tool-output-pq.txt, and
// writeFileExclusive - which was os.WriteFile, i.e. O_CREATE|O_TRUNC with no
// O_EXCL anywhere - let the later call replace the earlier call's bytes. The
// earlier Spill.Path had already been handed to the model and is re-readable via
// fs.read, so the model could re-read what it believed was its own tool result
// and get someone else's output. That is a C25 provenance break.
//
// What this file pins, in order of "what would a regression look like":
//  1. the collision triples from the report get three DIFFERENT names (AC#1);
//  2. a re-read of the FIRST call's path still returns the FIRST call's bytes
//     after a colliding second call has spilled - the property the model relies
//     on, stated as behavior rather than as a string comparison;
//  3. names round-trip: an independent decoder in this test turns the disk name
//     back into the id, so "injective" is proven and not merely printed;
//  4. writeFileExclusive fails with fs.ErrExist on an occupied name (AC#2) and
//     leaves the incumbent bytes alone;
//  5. a same-id retry still succeeds and still overwrites - the documented
//     behavior ticket 79 keeps, deliberately, and deliberately pins.
//  6. upper-case ids do not fold onto lower-case ones on a case-insensitive FS.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// spill79Budget spills anything over 40 tokens, so a test payload can be small.
var spill79Budget = Budgets{
	RawOutputCapBytes: 1 << 20,
	SpillTokens:       40,
	SpillHeadTokens:   10,
	SpillTailTokens:   10,
}

func spill79Payload(mark string) string { return strings.Repeat(mark, 200) }

// ---------------------------------------------------------------------------
// AC#1 — distinct ids, distinct names.
// ---------------------------------------------------------------------------

// TestArtifactNameDoesNotFoldDistinctIDs is AC#1: the exact collision triples
// from the report, each pair in its own named subtest.
func TestArtifactNameDoesNotFoldDistinctIDs(t *testing.T) {
	cases := []struct {
		sub string
		ids []string
	}{
		{"p_slash_vs_pq", []string{"p/q", "pq"}},
		{"p_backslash_vs_pq", []string{`p\q`, "pq"}},
		{"p_slash_vs_p_backslash", []string{"p/q", `p\q`}},
		{"the_whole_report_triple", []string{"p/q", `p\q`, "pq"}},
		{"dot_and_dotdot_vs_empty", []string{".", "..", ""}},
		{"colon_and_unc_shape", []string{`a:b`, `\\a\b`}},
		{"underscore_dash_and_percent", []string{"a_b", "a-b", "a%b"}},
		{"case_folding", []string{"pq", "PQ"}},
		{"digits_vs_letters", []string{"o0", "0o"}},
	}
	for _, tc := range cases {
		t.Run(tc.sub, func(t *testing.T) {
			seen := map[string]string{}
			for _, id := range tc.ids {
				name := artifactName(id, 1)
				if prev, dup := seen[name]; dup {
					t.Errorf("ids %q and %q both name %q: the encoder folds distinct ids "+
						"onto one disk name again, which is ticket 79's clobbering bug", prev, id, name)
				}
				seen[name] = id
				if filepath.Base(name) != name {
					t.Errorf("name %q for id %q is not a bare file name", name, id)
				}
				for _, bad := range []string{"/", "\\", ":", ".."} {
					if strings.Contains(name, bad) {
						t.Errorf("name %q for id %q contains %q", name, id, bad)
					}
				}
			}
			if len(seen) != len(tc.ids) {
				t.Errorf("%d distinct names for %d ids: %v", len(seen), len(tc.ids), seen)
			}
			t.Logf("names: %v", seen)
		})
	}
}

// TestSpilledBytesSurviveANameThatUsedToCollide is the defect as the model
// experiences it: call A spills, call B (a different id that A's name folded
// onto) spills, and A's path - already handed to the model - must still hold
// A's bytes. Under the old strip-and-truncate pair this returned B's payload.
func TestSpilledBytesSurviveANameThatUsedToCollide(t *testing.T) {
	dir := t.TempDir()
	sp := NewSpiller(dir, spill79Budget)

	for _, pair := range []struct {
		sub      string
		idA, idB string
	}{
		{"slash_id_then_bare_id", "p/q", "pq"},
		{"bare_id_then_backslash_id", "pq", `p\q`},
		{"dotdot_id_then_word_id", "../../escape", "escape"},
	} {
		t.Run(pair.sub, func(t *testing.T) {
			a, err := sp.Prepare(pair.idA, spill79Payload("AAAA"))
			if err != nil {
				t.Fatalf("Prepare(%q): %v", pair.idA, err)
			}
			if !a.Spilled {
				t.Fatal("A did not spill, so nothing was written to defend")
			}
			b, err := sp.Prepare(pair.idB, spill79Payload("BBBB"))
			if err != nil {
				t.Fatalf("Prepare(%q): %v", pair.idB, err)
			}
			if a.Name == b.Name || a.Path == b.Path {
				t.Fatalf("%q and %q still share the name %q", pair.idA, pair.idB, a.Name)
			}
			first, err := os.ReadFile(a.Path)
			if err != nil {
				t.Fatalf("re-reading A's path: %v", err)
			}
			if !strings.HasPrefix(string(first), "AAAA") || strings.Contains(string(first), "BBBB") {
				t.Errorf("A's artifact holds %q: A's bytes were replaced by B's output (C25 break)",
					first[:min(24, len(first))])
			}
			second, err := os.ReadFile(b.Path)
			if err != nil {
				t.Fatalf("re-reading B's path: %v", err)
			}
			if !strings.HasPrefix(string(second), "BBBB") {
				t.Errorf("B's artifact holds %q, want B's own payload", second[:min(24, len(second))])
			}
			// And the stub A handed the model must name the file A can still read.
			if !strings.Contains(a.Text, a.Path) {
				t.Error("A's stub does not name A's path")
			}
		})
	}
}

// decodeArtifactName turns a disk name back into the tool-call id, independently
// of the encoder: %XX (uppercase hex only) is reversed and nothing else is
// touched. ok=false means the name is not of the decodable form (the sequence
// fallback, or a long id carrying a digest tail).
func decodeArtifactName(name string) (id string, ok bool) {
	const pfx, sfx = "tool-output-", ".txt"
	if !strings.HasPrefix(name, pfx) || !strings.HasSuffix(name, sfx) {
		return "", false
	}
	core := strings.TrimSuffix(strings.TrimPrefix(name, pfx), sfx)
	var b strings.Builder
	for i := 0; i < len(core); i++ {
		c := core[i]
		if c != '%' {
			if !isSpill79Literal(c) {
				return "", false // any other byte must have been escaped
			}
			b.WriteByte(c)
			continue
		}
		if i+2 >= len(core) || !isSpill79UpperHex(core[i+1]) || !isSpill79UpperHex(core[i+2]) {
			return "", false
		}
		v := byte(hexVal(core[i+1]))<<4 | hexVal(core[i+2])
		if isSpill79Literal(v) {
			return "", false // a literal that got escaped: the encoder is not canonical
		}
		b.WriteByte(v)
		i += 2
	}
	return b.String(), true
}

func isSpill79Literal(c byte) bool {
	return c == '_' || c == '-' || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z')
}

func isSpill79UpperHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')
}

func hexVal(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	return c - 'A' + 10
}

// TestArtifactNameRoundTripsToTheExactID proves injectivity rather than restating
// it: for every id whose encoding fits the readable budget, the disk name decodes
// back to the id, byte for byte. Two ids sharing a name would have to share a
// decode, so a clean sweep here is a 1:1 proof over the corpus.
func TestArtifactNameRoundTripsToTheExactID(t *testing.T) {
	ids := []string{
		"pq", "p/q", `p\q`, "p_q", "p-q", "PQ", "p1", "", "....", "..", "!!!",
		"a:b", `\\fileserver\share\payload`, `C:\Windows\System32\drop`,
		"../../escape", "%", "%2F", "%zz", "tool-output-x.txt", "\x00\x7f\xff",
		"e/\u00e9", "call_01ABCdef", "a b\tc\nd", "CON", "nul.txt", "x.txt ",
	}
	seen := map[string]string{}
	for _, id := range ids {
		name := artifactName(id, 7)
		if prev, dup := seen[name]; dup {
			t.Errorf("ids %q and %q share the name %q", prev, id, name)
		}
		seen[name] = id
		if id == "" {
			// The one documented exception: the sequence fallback, whose token is
			// spelled with a dot the encoder can never emit for a non-empty id.
			if name != "tool-output-seq.7.txt" {
				t.Errorf("empty id produced %q, want the sequence fallback name", name)
			}
			continue
		}
		got, ok := decodeArtifactName(name)
		if !ok {
			t.Errorf("id %q -> %q does not decode back to a literal id", id, name)
			continue
		}
		if got != id {
			t.Errorf("id %q -> name %q -> decodes to %q, want the id back unchanged", id, name, got)
		}
	}
	if len(seen) != len(ids) {
		t.Errorf("only %d distinct names for %d ids", len(seen), len(ids))
	}
}

// TestArtifactNameIsCanonicalAndBounded pins the two properties the encoding
// relies on for a case-insensitive filesystem: %XX is the only upper-case a name
// can hold (so %2f and %2F can never be two spellings of one name, and a
// case-folded match between two different names is impossible), and the name
// stays inside a Windows path component even for an absurd id.
func TestArtifactNameIsCanonicalAndBounded(t *testing.T) {
	escapes := regexp.MustCompile(`%[0-9A-F]{2}`)
	long := strings.Repeat("A", 40) + "/" + strings.Repeat("b", 200)
	name := artifactName(long, 3)
	if len(name) > 140 {
		t.Errorf("name for a 240-byte id is %d bytes (%q); it must stay bounded", len(name), name)
	}
	if !strings.HasSuffix(name, ".txt") {
		t.Errorf("name %q lost its extension", name)
	}
	if !regexp.MustCompile(`-[0-9a-f]{16}\.txt$`).MatchString(name) {
		t.Errorf("over-long id name %q has no digest tail: cutting the readable part "+
			"alone would fold distinct long ids onto one name again", name)
	}
	core := strings.TrimSuffix(strings.TrimPrefix(name, "tool-output-"), ".txt")
	stripped := escapes.ReplaceAllString(core, "")
	if strings.ToLower(stripped) != stripped {
		t.Errorf("name %q holds upper case outside a %%XX escape, so a case-insensitive "+
			"filesystem would fold it onto a different id's name", name)
	}
	// The escape-boundary cut: no name may end mid-escape.
	if n := len(core); n >= 1 && core[n-1] == '%' {
		t.Errorf("name %q ends in a dangling %%", name)
	}
	// Distinct over-long ids must not fold either.
	two := artifactName(strings.Repeat("A", 40)+"/"+strings.Repeat("b", 200)+"!", 3)
	if two == name {
		t.Errorf("two different 240-byte ids share the bounded name %q", name)
	}
}

// ---------------------------------------------------------------------------
// AC#2 — the write really is exclusive, and a same-id retry is specified.
// ---------------------------------------------------------------------------

func TestWriteFileExclusiveIsExclusive(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "artifact.txt")

	if err := writeFileExclusive(p, []byte("FIRST")); err != nil {
		t.Fatalf("first exclusive write failed: %v", err)
	}
	// The name has to mean what it says: a second writer gets told, not absorbed.
	err := writeFileExclusive(p, []byte("SECOND"))
	if err == nil {
		t.Fatal("writeFileExclusive overwrote an existing file: O_EXCL is missing again, " +
			"which is half of ticket 79 (the other half is the folded name)")
	}
	if !errors.Is(err, fs.ErrExist) {
		t.Errorf("second write = %v, want fs.ErrExist", err)
	}
	body, rerr := os.ReadFile(p)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if string(body) != "FIRST" {
		t.Errorf("the incumbent's bytes are %q, want FIRST untouched by a refused write", body)
	}
	// NOTE (found while writing this, deliberately not asserted): the 0o600 the
	// create asks for is decorative on Windows - the artifact comes back
	// -rw-rw-rw- because the mode is inherited from the directory's ACL, so an
	// assertion here would be measuring Go's errno mapping rather than anything
	// this project controls. Recorded in ticket 79's Progress log instead.
}

// TestSpillSameIDRetryOverwrites is the other half of AC#2: making the primitive
// exclusive must not turn a RETRIED call into an error path, because ticket 79's
// names are injective - the only way this call can find its own name occupied is
// by being the same logical id again. Documented and pinned: last writer wins,
// the call succeeds, one file, no temp left behind.
func TestSpillSameIDRetryOverwrites(t *testing.T) {
	dir := t.TempDir()
	sp := NewSpiller(dir, spill79Budget)

	a, err := sp.Prepare("call_retry", spill79Payload("AAAA"))
	if err != nil || !a.Spilled {
		t.Fatalf("first spill: spilled=%v err=%v", a.Spilled, err)
	}
	b, err := sp.Prepare("call_retry", spill79Payload("BBBB"))
	if err != nil {
		t.Fatalf("a retried call must not become an error path: %v", err)
	}
	if b.Name != a.Name || b.Path != a.Path {
		t.Errorf("retry of the same id changed the name: %q then %q", a.Name, b.Name)
	}
	body, err := os.ReadFile(a.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(body), "BBBB") {
		t.Errorf("same-id retry left %q in the artifact, want the newer payload (last-writer-wins)",
			body[:min(24, len(body))])
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("a retry left %d entries (%v), want exactly the one artifact and no temp file",
			len(names), names)
	}
}

// TestSpillAcrossRestartsKeepsRetrySemantics: the documented behavior must not
// depend on this process having seen the id before (a restart with a repeated id
// is a retry too, not a failure).
func TestSpillAcrossRestartsKeepsRetrySemantics(t *testing.T) {
	dir := t.TempDir()
	if _, err := NewSpiller(dir, spill79Budget).Prepare("call_restart", spill79Payload("AAAA")); err != nil {
		t.Fatal(err)
	}
	sp2, err := NewSpiller(dir, spill79Budget).Prepare("call_restart", spill79Payload("CCCC"))
	if err != nil {
		t.Fatalf("a fresh Spiller hitting the same id errored instead of replacing: %v", err)
	}
	body, err := os.ReadFile(sp2.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(body), "CCCC") {
		t.Errorf("artifact holds %q after a restart-retry, want CCCC", body[:min(24, len(body))])
	}
}
