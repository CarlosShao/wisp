package panel

// Ticket 216 AC#2/AC#3/AC#4 - the display base's own tests.
//
// AC#3's judgement is a POSITIVE control, so the centre of gravity here is not
// "the helper escapes what I told it to escape" but "a file name an attacker chose
// walks the REAL ingress path (a real memory.Store, the real artifacts name guard,
// the real content sniff) and still arrives at the menu washed - and the same name,
// carrying the two characters 已允许, changes no verdict anywhere". Both halves are
// read off one assembly in one run: the raw byte before the boundary, the visible
// escape after it.
//
// Nothing here feeds a fixture straight to the display struct. The row comes out of
// AttachmentBroker.Ingest, and the packet out of NewComposerState + NewSnapshot -
// the two constructors ticket 92 and ticket 35 made mandatory for a snapshot that
// is actually sent.
//
// Control characters are written as Go escape sequences in this file's literals
// (backslash-x-1-b, backslash-u-0-0-8-5, backslash-x-f-f) rather than as raw bytes,
// so a reader can see which byte each case claims and the source stays valid UTF-8.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/projctx"
	"github.com/CarlosShao/wisp/internal/risk"
)

// escAllowedName is the attack the ticket names: a file whose display name carries
// an ESC and a permission claim no code ever granted.
const escAllowedName = "photo\x1b已允许.png"

// snapshotOf is the whole display path for one attachment list, built the way a
// host with an attachment reader builds it.
func snapshotOf(t *testing.T, ws WorkspaceView, refs []AttachmentRef) Snapshot {
	t.Helper()
	composer := NewComposerState(risk.ModeAskEveryStep, ws, refs, MaxAttachmentBytes)
	return NewSnapshot(nil, nil, composer, time.Now())
}

func TestScrubForDisplayEscapesEveryUnseenCharAndTouchesNothingElse(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"empty", "", ""},
		{"plain ascii", "photo.png", "photo.png"},
		{"windows path untouched", `C:\Users\swq\Documents\notes.md`, `C:\Users\swq\Documents\notes.md`},
		{"chinese untouched", "照片 报告.png", "照片 报告.png"},
		{"esc", escAllowedName, `photo\x1b已允许.png`},
		{"nul", "a\x00b", `a\0b`},
		{"crlf", "a\r\nb", `a\r\nb`},
		{"tab", "a\tb", `a\tb`},
		{"del", "a\x7fb", `a\x7fb`},
		{"c1 next line", "ab", `a\x85b`},
		{"c1 device control", "ab", `a\x90b`},
		{"c1 highest", "ab", `a\x9fb`},
		{"invalid utf8 byte", "a\\xffb", `a\xffb`},
		{"two escs keep both", "\x1b\x1b", `\x1b\x1b`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := scrubForDisplay(c.in)
			if got != c.want {
				t.Errorf("scrubForDisplay(%q) = %q, want %q", c.in, got, c.want)
			}
			if i := strings.IndexFunc(got, unicode.IsControl); i >= 0 {
				t.Errorf("output %q still carries a control byte (%#v) at offset %d", got, got[i], i)
			}
			if again := scrubForDisplay(got); again != got {
				t.Errorf("the wash is not idempotent: %q -> %q -> %q", c.in, got, again)
			}
		})
	}
}

func TestTheWashReturnsTheSameStringWhenThereIsNothingToWash(t *testing.T) {
	// Not a nicety: "washed" must never mean "re-encoded", or the identity of an
	// ordinary path becomes an accident of this function's encoder.
	in := `D:\work\workspace\projects plans\Wisp`
	if got := scrubForDisplay(in); got != in {
		t.Errorf("scrubForDisplay changed an ordinary path: %q -> %q", in, got)
	}
	if !needsDisplayWash(escAllowedName) {
		t.Error("needsDisplayWash missed an ESC")
	}
	if needsDisplayWash(in) {
		t.Error("needsDisplayWash fired on a clean path")
	}
}

func TestTheSnapshotBoundaryWashesBothRosterFieldsAndCopiesTheRows(t *testing.T) {
	hostilePath := "C:\\data\\\u0085reports"
	refs := []AttachmentRef{
		{ID: "att-1", Name: escAllowedName, MIME: "image/png", Kind: "image", Stored: true},
	}
	snap := snapshotOf(t, WorkspaceView{Set: true, Spelling: hostilePath, Canonical: hostilePath}, refs)

	if got := snap.Composer.Attachments[0].Name; strings.ContainsRune(got, 0x1b) || !strings.Contains(got, `\x1b`) {
		t.Errorf("the attachment row reached the menu as %q, want the ESC shown as a visible escape", got)
	}
	if got := snap.Composer.Workspace.Canonical; strings.ContainsRune(got, 0x85) || !strings.Contains(got, `\x85`) {
		t.Errorf("the canonical path reached the menu as %q, want the C1 shown as a visible escape", got)
	}
	if refs[0].Name != escAllowedName {
		t.Errorf("the boundary rewrote the caller's own row (%q): the wash must copy, not mutate", refs[0].Name)
	}
	// The page reads bytes, not the Go value, so the packet is round-tripped.
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	var wire struct {
		Composer struct {
			Workspace   WorkspaceView   `json:"workspace"`
			Attachments []AttachmentRef `json:"attachments"`
		} `json:"composer"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(wire.Composer.Attachments[0].Name, "已允许") {
		t.Errorf("washing lost the readable half of the name: %q", wire.Composer.Attachments[0].Name)
	}
	for _, s := range []string{wire.Composer.Attachments[0].Name, wire.Composer.Workspace.Canonical} {
		if strings.IndexFunc(s, unicode.IsControl) >= 0 {
			t.Errorf("the packet the page renders still carries a control character: %q", s)
		}
	}
}

func TestAHostileAttachmentNameWalksTheRealBrokerIntoAWashedMenuRow(t *testing.T) {
	broker := openBroker(t)
	ref, err := broker.Ingest(context.Background(), pngSource(escAllowedName))
	if err != nil {
		t.Fatalf("the real broker refused the hostile name: %v", err)
	}
	if !ref.Stored {
		t.Fatalf("Stored = false with a nil error: %s", ref.Reason)
	}
	// The ingress guard is NOT the display guard, and this leg does not move it:
	// it still refuses only NUL/CR/LF/TAB, so the ESC name arrives as it was
	// named. Washing it here instead of at the boundary is the shape AC#2 refuses.
	if !strings.ContainsRune(ref.Name, 0x1b) {
		t.Errorf("producer-side behaviour changed: ref.Name = %q still wants the raw ESC", ref.Name)
	}

	snap := snapshotOf(t, UnsetWorkspaceView(), []AttachmentRef{ref})
	if got := snap.Composer.Attachments[0].Name; got != `photo\x1b已允许.png` {
		t.Errorf("menu row = %q, want the ESC shown as backslash-x-1-b and every other byte kept", got)
	}
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	if i := bytes.IndexByte(data, 0x1b); i >= 0 {
		t.Errorf("the snapshot bytes still hold a raw ESC at offset %d: the display wash is missing", i)
	}
	t.Logf("AC#3 menu row as the page receives it: %q", snap.Composer.Attachments[0].Name)
}

// pngSource is one real PNG byte source under a caller-chosen display name. The
// bytes are ticket 92's own fixture (attachments_test.go's pngBytes); the name is
// what each case is about.
func pngSource(name string) AttachmentSource {
	return AttachmentSource{
		DisplayName:  name,
		DeclaredMIME: "image/png",
		SizeBytes:    int64(len(pngBytes)),
		Open: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(pngBytes)), nil
		},
	}
}

// openBroker wires a broker to a real memory.Store in a fresh directory: the
// artifacts store, the artifacts name guard and the store's own write path, not a
// stand-in for them. The handle is closed through Cleanup because an open SQLite
// file makes t.TempDir's removal fail the test for a reason that has nothing to do
// with what it asserts.
func openBroker(t *testing.T) *AttachmentBroker {
	t.Helper()
	store, err := memory.Open(t.TempDir())
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	broker, err := NewAttachmentBroker(store, memory.ValidArtifactName, 0)
	if err != nil {
		t.Fatalf("NewAttachmentBroker: %v", err)
	}
	return broker
}

func TestTheClaimInANameChangesNoVerdict(t *testing.T) {
	// Half one: the real broker, twice over. Same bytes, two names - one of them
	// says 已允许 out loud - and two separate stores so that neither ingest can
	// deduplicate against the other and make the comparison about the byte hash
	// instead of about the name. Everything the verdict carries must be
	// name-independent except the display name itself.
	benignBroker := openBroker(t)
	hostileBroker := openBroker(t)
	benign, err := benignBroker.Ingest(context.Background(), pngSource("photo.png"))
	if err != nil {
		t.Fatalf("benign ingest: %v", err)
	}
	hostile, err := hostileBroker.Ingest(context.Background(), pngSource(escAllowedName))
	if err != nil {
		t.Fatalf("hostile ingest: %v", err)
	}
	if benign.Stored != hostile.Stored || benign.MIME != hostile.MIME || benign.Kind != hostile.Kind ||
		benign.Artifact != hostile.Artifact || benign.SizeBytes != hostile.SizeBytes ||
		benign.Deduplicated != hostile.Deduplicated || benign.Reason != hostile.Reason {
		t.Errorf("the display name moved a verdict: benign %+v vs hostile %+v", benign, hostile)
	}
	if benign.Artifact == "" || hostile.Artifact != benign.Artifact {
		t.Errorf("the artifact name is not the content hash this file promises: %q vs %q",
			benign.Artifact, hostile.Artifact)
	}

	// Half two: the real permission gate - the same assessor the tool bridge runs
	// (approval.go's header), reached through the panel's own card renderer. A card
	// for a call whose argument carries the hostile name must get the verdict the
	// same call gets without it, raw or washed alike.
	assessor := risk.NewRiskAssessor()
	cardOf := func(arg string) ApprovalCardView {
		return NewApprovalCardView(assessor, ApprovalSubject{
			CorrelationID: "corr-216",
			Tool:          "fs.write",
			Args:          []string{arg},
			Facts: risk.Facts{
				Declared:     risk.L1,
				Paths:        []string{arg},
				Irreversible: []string{"overwrite"},
			},
		})
	}
	plain := cardOf("photo.png")
	hostileCard := cardOf(escAllowedName)
	washedCard := cardOf(scrubForDisplay(escAllowedName))
	t.Logf("real assessor verdicts: plain=%s/%v hostile=%s/%v washed=%s/%v",
		plain.Level, plain.RulesHit, hostileCard.Level, hostileCard.RulesHit, washedCard.Level, washedCard.RulesHit)
	if hostileCard.Level != plain.Level || washedCard.Level != plain.Level {
		t.Errorf("a name that says 已允许 changed the risk level: plain=%q hostile=%q washed=%q",
			plain.Level, hostileCard.Level, washedCard.Level)
	}
	if strings.Join(hostileCard.RulesHit, ",") != strings.Join(plain.RulesHit, ",") {
		t.Errorf("the rules hit differ only because of the name: plain=%v hostile=%v", plain.RulesHit, hostileCard.RulesHit)
	}
	if hostileCard.Level == "L0" {
		t.Errorf("the card downgraded an irreversible overwrite because of a file name")
	}
	if hostileCard.DecidedBy != "native" {
		t.Errorf("DecidedBy = %q: the verdict did not come from the native gate", hostileCard.DecidedBy)
	}
}

func TestWorkspaceAbsentAndWorkspaceUnreadAreTwoSentences(t *testing.T) {
	// AC#4 on this ticket's own field: one empty row must not be able to say both
	// "nobody chose a workspace" and "nobody read the dimension".
	unread := NewSnapshotPump(PumpSources{}).Snapshot()
	readButNone := NewSnapshotPump(PumpSources{Workspace: UnsetWorkspaceView}).Snapshot()
	blank := NewSnapshot(nil, nil, ComposerState{}, time.Now())

	for _, pair := range []struct {
		what string
		got  WorkspaceView
	}{
		{"unread", unread.Composer.Workspace},
		{"read-but-none", readButNone.Composer.Workspace},
		{"blank section", blank.Composer.Workspace},
	} {
		if pair.got.Set {
			t.Errorf("%s claims a workspace is set: %+v", pair.what, pair.got)
		}
		if strings.TrimSpace(pair.got.Reason) == "" {
			t.Errorf("%s renders an empty row with no sentence: AC#4 forbids the silence", pair.what)
		}
	}
	if unread.Composer.Workspace.Reason == readButNone.Composer.Workspace.Reason {
		t.Errorf("the two truths collapsed into one sentence: %q", unread.Composer.Workspace.Reason)
	}
	if unread.Composer.Workspace.Reason != workspaceNoReadingReason ||
		readButNone.Composer.Workspace.Reason != workspaceNotChosenReason {
		t.Errorf("sentence mismatch: unread=%q none=%q",
			unread.Composer.Workspace.Reason, readButNone.Composer.Workspace.Reason)
	}
	if !strings.Contains(unread.Composer.Workspace.Reason, InstructionsStatusNotRun) {
		t.Errorf("the no-reading sentence does not name the precedent state %q it borrows: %q",
			InstructionsStatusNotRun, unread.Composer.Workspace.Reason)
	}
	t.Logf("AC#4 pair: 没有=%q / 读不到=%q",
		readButNone.Composer.Workspace.Reason, unread.Composer.Workspace.Reason)
}

func TestTheFiveInstructionStatesAreStillTheDiskSpelling(t *testing.T) {
	// The census's one correction, nailed so it cannot rot: the fifth state is
	// not_run (the dispatch wrote "na"; 216-a1 pushed back, the orchestrator re-ran
	// the ruler, and AC#4 therefore runs on not_run).
	for spelling, got := range map[string]string{
		"loaded":     InstructionsStatusLoaded,
		"none_found": InstructionsStatusNoneFound,
		"off":        InstructionsStatusOff,
		"refused":    InstructionsStatusRefused,
		"not_run":    InstructionsStatusNotRun,
	} {
		if got != spelling {
			t.Errorf("instruction state spelled %q on disk, want %q", got, spelling)
		}
	}
	// And the two halves AC#4 cares about stay two sections, not one missing key:
	// none_found says "we looked and there was nothing", not_run says "nobody has
	// looked yet".
	looked := InstructionsSectionFromBundle(&projctx.Bundle{})
	never := InstructionsSectionFromBundle(nil)
	if looked.Status == never.Status {
		t.Errorf("none_found and not_run collapsed into %q", looked.Status)
	}
	if never.Status != InstructionsStatusNotRun || strings.TrimSpace(never.Reason) == "" {
		t.Errorf("not_run arrived without its sentence: %+v", never)
	}
	if len(never.Files) != 0 {
		t.Errorf("not_run also invented files: %+v", never.Files)
	}
}

// ---- ticket 327: the invisible formatting runes ----
//
// The next two tests are the AC#2 pair, built to fail in DIFFERENT places so that
// a half-toothed implementation cannot pass both:
//   - TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired can only go red if
//     the WRITER loses its Cf branch - its input also carries an ESC, so the judge
//     answers "wash" whatever it thinks of Cf, and every byte except the Cf escapes
//     is already right;
//   - TestNeedsDisplayWashNamesEveryWashedCfRune can only go red if the JUDGE loses
//     its Cf term - it calls needsDisplayWash and nothing else.
// Tearing out both halves reddens both; tearing out one reddens exactly one.

func TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired(t *testing.T) {
	in := "a\x1b\u202eb\u200dc\xad"
	want := `a\x1b\u202eb\u200dc\xad`
	got := scrubForDisplay(in)
	if got != want {
		t.Errorf("scrubForDisplay(%q) = %q, want %q", in, got, want)
	}
	if again := scrubForDisplay(got); again != got {
		t.Errorf("the wash is not idempotent: %q -> %q -> %q", in, got, again)
	}
	if i := strings.IndexFunc(got, func(r rune) bool {
		return isDisplayWashCf(r) || unicode.IsControl(r)
	}); i >= 0 {
		t.Errorf("output %q still carries an invisible rune %U at offset %d", got, got[i], i)
	}
}

func TestNeedsDisplayWashNamesEveryWashedCfRune(t *testing.T) {
	for _, c := range cfRoster {
		in := "photo" + string(c.r) + ".png"
		if !needsDisplayWash(in) {
			t.Errorf("needsDisplayWash did not see %s (%U) in %q: the fast path would return it raw", c.name, c.r, in)
		}
	}
}

// cfRoster is the display boundary's invisible list, spelled out so a reader can
// count it: nine directional runes, six spacing ones, no Unicode category ranges.
var cfRoster = []struct {
	r    rune
	name string
}{
	{0x00AD, "soft hyphen"},
	{0x200B, "zero width space"},
	{0x200C, "zero width non-joiner"},
	{0x200D, "zero width joiner"},
	{0x2060, "word joiner"},
	{0xFEFF, "byte order mark"},
	{0x202A, "left-to-right embedding"},
	{0x202B, "right-to-left embedding"},
	{0x202C, "pop directional formatting"},
	{0x202D, "left-to-right override"},
	{0x202E, "right-to-left override"},
	{0x2066, "left-to-right isolate"},
	{0x2067, "right-to-left isolate"},
	{0x2068, "first strong isolate"},
	{0x2069, "pop directional isolate"},
}

// TestEveryNamedCfRuneIsEscapedAndNeverDeleted is AC#3's shape: the escaped half of
// the name must still be there, the invisible half must be readable, and the two
// visible halves must not have been spliced into one new name.
func TestEveryNamedCfRuneIsEscapedAndNeverDeleted(t *testing.T) {
	if len(cfRoster) != 15 {
		t.Fatalf("the roster is %d runes, the ruling names fifteen", len(cfRoster))
	}
	for _, c := range cfRoster {
		in := "photo" + string(c.r) + "报告.png"
		got := scrubForDisplay(in)
		if strings.ContainsRune(got, c.r) {
			t.Errorf("%s (%U) survived the wash as itself in %q", c.name, c.r, got)
		}
		want := "photo" + escapeOf(c.r) + "报告.png"
		if got != want {
			t.Errorf("scrubForDisplay(%q) = %q, want %q", in, got, want)
		}
		if again := scrubForDisplay(got); again != got {
			t.Errorf("%s: the wash is not idempotent: %q -> %q", c.name, got, again)
		}
		if !strings.HasPrefix(got, "photo") || !strings.HasSuffix(got, "报告.png") {
			t.Errorf("%s: washing spliced the two visible halves into a different name: %q", c.name, got)
		}
	}
}

// escapeOf is the boundary's own vocabulary, restated here so a wrong table is a
// red test rather than a shared bug: writeControlEscape prints \n \r \t \0 as
// mnemonics, \xNN below U+0100 and \uXXXX above it.
func escapeOf(r rune) string {
	switch r {
	case '\n':
		return `\n`
	case '\r':
		return `\r`
	case '\t':
		return `\t`
	}
	if r <= 0xff {
		return fmt.Sprintf(`\x%02x`, r)
	}
	return fmt.Sprintf(`\u%04x`, r)
}

// TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut pins the other half of AC#1:
// the invisible runes the census measured as reachable but the ruling did not name
// stay raw, so widening the set is a ruling and not an edit a later leg can make
// quietly.
func TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut(t *testing.T) {
	for _, r := range []rune{0x2061, 0x2062, 0x2063, 0x2064, 0x061C, 0x180E, 0x206A} {
		in := "photo" + string(r) + ".png"
		if needsDisplayWash(in) {
			t.Errorf("U+%04X is not on the ruled list but needsDisplayWash fired on it", r)
		}
		if got := scrubForDisplay(in); got != in {
			t.Errorf("U+%04X is not on the ruled list but the wash changed it: %q -> %q", r, in, got)
		}
	}
}

// TestABidiNameWalksTheRealBoundaryIntoAWashedRow is AC#0's field pair run through
// the two constructors a snapshot actually goes through: the disk-spelled workspace
// path and the real attachment broker.
func TestABidiNameWalksTheRealBoundaryIntoAWashedRow(t *testing.T) {
	const bidiName = "photo\u202e已允许.png"
	broker := openBroker(t)
	ref, err := broker.Ingest(context.Background(), pngSource(bidiName))
	if err != nil {
		t.Fatalf("the real broker refused the bidi name: %v", err)
	}
	if !strings.ContainsRune(ref.Name, 0x202e) {
		t.Fatalf("producer-side behaviour changed: ref.Name = %q still wants the raw U+202E", ref.Name)
	}
	// Spelling is the user's own string and composer.go says it goes out verbatim,
	// so it is NOT on this roster; the Cf lands where the ticket's entry A puts it,
	// on the C26-resolved Canonical that the user never typed.
	snap := snapshotOf(t,
		WorkspaceView{Set: true, Spelling: `C:\data\reports`, Canonical: "C:\\data\\\u202ereports"},
		[]AttachmentRef{ref})

	if got := snap.Composer.Attachments[0].Name; got != `photo\u202e已允许.png` {
		t.Errorf("menu row = %q, want the override shown as a visible escape", got)
	}
	if got := snap.Composer.Workspace.Canonical; strings.ContainsRune(got, 0x202e) || !strings.Contains(got, `\u202e`) {
		t.Errorf("canonical path = %q, want the override shown as a visible escape", got)
	}
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	if i := bytes.Index(data, []byte{0xe2, 0x80, 0xae}); i >= 0 {
		t.Errorf("the snapshot bytes still hold raw U+202E at offset %d: the display wash is missing", i)
	}
	if i := bytes.Index(data, []byte(`\u202e`)); i < 0 {
		t.Errorf("the snapshot bytes carry no visible escape anywhere; row=%q canonical=%q",
			snap.Composer.Attachments[0].Name, snap.Composer.Workspace.Canonical)
	}
	t.Logf("327 menu row as the packet carries it: %q", snap.Composer.Attachments[0].Name)
}

// TestTheBidiClaimInANameChangesNoVerdict is AC#4: the same call, the same real
// assessor, three spellings of the argument - the honest name, the name that says
// 已允许 behind a right-to-left override, and that name after the display wash. The
// verdicts must be letter-for-letter identical; only the string a human reads moves.
func TestTheBidiClaimInANameChangesNoVerdict(t *testing.T) {
	const hostile = "photo\u202e已允许.png"
	assessor := risk.NewRiskAssessor()
	cardOf := func(arg string) ApprovalCardView {
		return NewApprovalCardView(assessor, ApprovalSubject{
			CorrelationID: "corr-327",
			Tool:          "fs.write",
			Args:          []string{arg},
			Facts: risk.Facts{
				Declared:     risk.L1,
				Paths:        []string{arg},
				Irreversible: []string{"overwrite"},
			},
		})
	}
	plain := cardOf("photo.png")
	hostileCard := cardOf(hostile)
	washedCard := cardOf(scrubForDisplay(hostile))
	t.Logf("327 verdicts: plain=%s/%v hostile=%s/%v washed=%s/%v",
		plain.Level, plain.RulesHit, hostileCard.Level, hostileCard.RulesHit, washedCard.Level, washedCard.RulesHit)
	if hostileCard.Level != plain.Level || washedCard.Level != plain.Level {
		t.Errorf("a name with an invisible override changed the risk level: plain=%q hostile=%q washed=%q",
			plain.Level, hostileCard.Level, washedCard.Level)
	}
	if strings.Join(hostileCard.RulesHit, ",") != strings.Join(plain.RulesHit, ",") {
		t.Errorf("the rules hit differ only because of the name: plain=%v hostile=%v", plain.RulesHit, hostileCard.RulesHit)
	}
	if washedCard.Level == "L0" {
		t.Errorf("the card downgraded an irreversible overwrite because of a file name")
	}
	if washedCard.DecidedBy != "native" {
		t.Errorf("DecidedBy = %q: the verdict did not come from the native gate", washedCard.DecidedBy)
	}
	if !strings.Contains(scrubForDisplay(hostile), `\u202e`) {
		t.Errorf("the display string reached the human unwashed: %q", scrubForDisplay(hostile))
	}
}
