package config

// Ticket 195 - the direction verdict on the programmatic write path.
//
// The claim these cases pin is narrow and it is the ticket's whole remaining
// deliverable (orchestrator ruling 2026-10-02 sections 1 and 2): a write that
// would LOOSEN a locked section must not reach config.toml through the
// programmatic entry, while a write that TIGHTENS one must land normally. The
// control at the bottom runs the SAME scenario through the write body as it
// stood before this ticket, and shows the loosening landing on disk - without
// that, every "nothing was written" assertion below could be green because the
// scenario never carried a real loosening in the first place.

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// sectionWriteTightBase is settingsBase plus one locked-section value that is
// already on: the tightening cases start here so "wrote it back to false" is a
// real change of direction and not a default.
const sectionWriteTightBase = settingsBase + `
[fs]
delete_enabled = true
`

// testLoosenQueue records what the entry handed over. It never blocks: the
// contract the entry carries is "queue and return" (ruling N9).
type testLoosenQueue struct {
	mu   sync.Mutex
	reqs []LoosenRequest
	next int
}

func (q *testLoosenQueue) QueueLoosen(req LoosenRequest) (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.reqs = append(q.reqs, req)
	q.next++
	return fmt.Sprintf("card-%d", q.next), nil
}

func (q *testLoosenQueue) snapshot() []LoosenRequest {
	q.mu.Lock()
	defer q.mu.Unlock()
	return append([]LoosenRequest(nil), q.reqs...)
}

// auditLog records the attribution lines the entry wrote (AC#2's third
// assertion: the audit says who asked for what).
type auditLog struct {
	mu    sync.Mutex
	lines []string
}

func (a *auditLog) f(format string, args ...any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lines = append(a.lines, fmt.Sprintf(format, args...))
}

func (a *auditLog) all() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.lines...)
}

func fileBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	return data
}

// AC#2, forward: a native request to loosen [fs] is NOT written. The file is
// byte-identical, the in-memory config is unchanged, one request is queued and
// the audit line names the source and the origin.
func TestAC2LooseningIsQueuedAndNeverWritten(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	q := &testLoosenQueue{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	before := fileBytes(t, path)
	res, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply:   func(c *Config) error { c.FS.DeleteEnabled = true; return nil },
	})
	if err != nil {
		t.Fatalf("a native loosening must be queued, not fail: %v", err)
	}
	if res.Written {
		t.Errorf("the loosening was reported as written: %+v", res)
	}
	if res.Direction != DirLoosen {
		t.Errorf("direction = %q, want %q", res.Direction, DirLoosen)
	}
	if res.PendingID != "card-1" {
		t.Errorf("pending id = %q, want card-1: the entry must report the queued request", res.PendingID)
	}
	if want := []string{"fs.delete_enabled"}; !reflect.DeepEqual(res.Keys, want) {
		t.Errorf("queued keys = %v, want %v", res.Keys, want)
	}

	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("config.toml changed on a queued loosening:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if m.Config().FS.DeleteEnabled {
		t.Errorf("the in-memory config carries the unconfirmed loosening")
	}

	reqs := q.snapshot()
	if len(reqs) != 1 {
		t.Fatalf("queued requests = %d, want exactly 1", len(reqs))
	}
	got := reqs[0]
	if got.Section != "fs" || got.Source != SourceNative || got.Origin != "cli.test" {
		t.Errorf("queued request = %+v, want section fs / source native / origin cli.test", got)
	}
	if want := []string{"fs.delete_enabled"}; !reflect.DeepEqual(got.Keys, want) {
		t.Errorf("queued request keys = %v, want %v", got.Keys, want)
	}

	lines := audit.all()
	if len(lines) != 1 {
		t.Fatalf("audit lines = %v, want exactly 1", lines)
	}
	for _, want := range []string{"section=fs", "fs.delete_enabled", "source=native", "origin=cli.test", "outcome=queued", "card-1"} {
		if !strings.Contains(lines[0], want) {
			t.Errorf("audit line %q does not carry %q", lines[0], want)
		}
	}
}

// AC#2, reverse: a TIGHTENING of the same section lands without any card. The
// ticket names this half explicitly - a gate that also swallowed the fail-closed
// direction would bind approval and tightening into one thing.
func TestAC2TighteningIsWrittenWithoutACard(t *testing.T) {
	m, path, _ := newTestManager(t, sectionWriteTightBase)
	q := &testLoosenQueue{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	res, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply:   func(c *Config) error { c.FS.DeleteEnabled = false; return nil },
	})
	if err != nil {
		t.Fatalf("a tightening must be written, not refused: %v", err)
	}
	if !res.Written {
		t.Errorf("the tightening was not written: %+v", res)
	}
	if res.Direction != DirTighten {
		t.Errorf("direction = %q, want %q", res.Direction, DirTighten)
	}
	if res.PendingID != "" {
		t.Errorf("a tightening raised a card (%q): no approval is owed for the fail-closed direction", res.PendingID)
	}
	if reqs := q.snapshot(); len(reqs) != 0 {
		t.Errorf("a tightening was queued: %+v", reqs)
	}
	got := readLoaded(t, path)
	if got.FS.DeleteEnabled {
		t.Errorf("the tightening never reached the file")
	}
	if m.Config().FS.DeleteEnabled {
		t.Errorf("the in-memory config kept the looser value")
	}

	lines := audit.all()
	if len(lines) != 1 || !strings.Contains(lines[0], "direction=tighten") || !strings.Contains(lines[0], "outcome=written") {
		t.Errorf("audit lines = %v, want one written/tighten line", lines)
	}
}

// AC#3: a panel-sourced loosening of a locked section is refused outright.
// Nothing is written, nothing is queued, and the audit says why.
func TestAC3PanelLooseningOfALockedSectionIsRefused(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	q := &testLoosenQueue{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	before := fileBytes(t, path)
	res, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourcePanel,
		Origin:  "panel.config.set",
		Apply:   func(c *Config) error { c.FS.DeleteEnabled = true; return nil },
	})
	if !errors.Is(err, ErrPanelLoosenRefused) {
		t.Fatalf("panel loosening error = %v, want ErrPanelLoosenRefused", err)
	}
	if res.Written || res.PendingID != "" {
		t.Errorf("a refused panel loosening reported progress: %+v", res)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("config.toml changed on a refused panel loosening")
	}
	if reqs := q.snapshot(); len(reqs) != 0 {
		t.Errorf("a refused panel loosening still reached the queue: %+v", reqs)
	}
	lines := audit.all()
	if len(lines) != 1 || !strings.Contains(lines[0], "outcome=refused-panel") {
		t.Errorf("audit lines = %v, want one refused-panel line", lines)
	}
}

// AC#3's boundary: the refusal is scoped to LOOSENING. A panel-sourced
// tightening is a fail-closed change and is written like any other - refusing it
// too would be the over-block the ticket warns against.
func TestAC3PanelTighteningIsNotSwallowed(t *testing.T) {
	m, path, _ := newTestManager(t, sectionWriteTightBase)
	m.LoosenQueue = &testLoosenQueue{}

	res, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourcePanel,
		Origin:  "panel.config.set",
		Apply:   func(c *Config) error { c.FS.DeleteEnabled = false; return nil },
	})
	if err != nil {
		t.Fatalf("a panel tightening was refused: %v", err)
	}
	if !res.Written || readLoaded(t, path).FS.DeleteEnabled {
		t.Errorf("the panel tightening did not land: %+v", res)
	}
}

// The nil queue is a deny, the same rule ConfirmLocked carries: a loosening with
// nowhere to go is not written either.
func TestAC2LooseningWithNoQueueWiredFailsClosed(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	before := fileBytes(t, path)

	res, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply:   func(c *Config) error { c.FS.DeleteEnabled = true; return nil },
	})
	if !errors.Is(err, ErrLoosenQueueUnwired) {
		t.Fatalf("error = %v, want ErrLoosenQueueUnwired", err)
	}
	if res.Written {
		t.Errorf("an unwired loosening was written: %+v", res)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("config.toml changed with no queue wired")
	}
}

// AC#4's reverse: the generic entry does not have a shortcut to
// SetPermissionMode. A mode change offered as a loosening is queued like any
// other locked key, and the file keeps its old mode.
func TestAC4PermissionModeLooseningIsQueuedNotApplied(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	q := &testLoosenQueue{}
	m.LoosenQueue = q
	before := fileBytes(t, path)

	res, err := m.WriteSection(SectionWrite{
		Section: "risk",
		Key:     "risk.permission_mode",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply: func(c *Config) error {
			c.Risk.PermissionMode = "auto_approve"
			return nil
		},
	})
	if err != nil {
		t.Fatalf("queueing a mode loosening: %v", err)
	}
	if res.Written || res.PendingID == "" {
		t.Errorf("the mode loosening did not go to the queue: %+v", res)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("the mode loosening was written to config.toml")
	}
	if got := m.Config().Risk.PermissionMode; got != "ask_every_step" {
		t.Errorf("the in-memory mode changed to %q without a card", got)
	}
	if reqs := q.snapshot(); len(reqs) != 1 || !reflect.DeepEqual(reqs[0].Keys, []string{"risk.permission_mode"}) {
		t.Errorf("queued requests = %+v, want one risk.permission_mode request", reqs)
	}
}

// One call, one key. An Apply that moves a second key is refused before
// anything is written, so a receipt line can never cover more than it says.
func TestSectionWriteRefusesToCarryASecondKey(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	before := fileBytes(t, path)

	_, err := m.WriteSection(SectionWrite{
		Section: "fs",
		Key:     "fs.delete_enabled",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply: func(c *Config) error {
			c.FS.DeleteEnabled = true
			c.FS.AllowedDirs = append(c.FS.AllowedDirs, `D:\second-key`)
			return nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "outside the one it owns") {
		t.Fatalf("error = %v, want the unrelated-key refusal", err)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("config.toml changed on a refused multi-key write")
	}
}

// The section must be the key's own first segment, and it must exist.
func TestSectionWriteRefusesAMismatchedOrUnknownSection(t *testing.T) {
	cases := []struct{ name, section, key string }{
		{"section does not match the key", "llm", "fs.delete_enabled"},
		{"section does not exist", "nosuchsection", "nosuchsection.key"},
		{"no section named", "", "fs.delete_enabled"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, path, _ := newTestManager(t, settingsBase)
			before := fileBytes(t, path)
			_, err := m.WriteSection(SectionWrite{
				Section: tc.section,
				Key:     tc.key,
				Source:  SourceNative,
				Origin:  "cli.test",
				Apply:   func(c *Config) error { c.FS.DeleteEnabled = true; return nil },
			})
			if err == nil {
				t.Fatalf("accepted a write whose section is %q and whose key is %q", tc.section, tc.key)
			}
			if after := fileBytes(t, path); string(after) != string(before) {
				t.Errorf("config.toml changed on a refused write")
			}
		})
	}
}

// The unlocked half of the same entry: a hot key writes normally, with no card
// and a neutral verdict, so the entry is a section write and not a locked-only
// special case.
func TestSectionWriteWritesAnUnlockedSectionNormally(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	q := &testLoosenQueue{}
	m.LoosenQueue = q

	res, err := m.WriteSection(SectionWrite{
		Section: "ball",
		Key:     "ball.size",
		Source:  SourceNative,
		Origin:  "cli.test",
		Apply:   func(c *Config) error { c.Ball.Size = 64; return nil },
	})
	if err != nil {
		t.Fatalf("unlocked write: %v", err)
	}
	if !res.Written || res.Direction != DirNeutral {
		t.Errorf("unlocked write = %+v, want written/neutral", res)
	}
	if got := readLoaded(t, path).Ball.Size; got != 64 {
		t.Errorf("ball.size = %d, want 64", got)
	}
	if reqs := q.snapshot(); len(reqs) != 0 {
		t.Errorf("an unlocked write was queued: %+v", reqs)
	}
}

// CONTROL (the mutation this ticket removes): the write body as it stood before
// ticket 195 wrote the same loosening straight to config.toml - no verdict, no
// card. It is kept here so the AC#2 cases above cannot be green for the wrong
// reason (a scenario that never carried a loosening).
func TestSectionWriteControlOldBodyPersistsALoosening(t *testing.T) {
	m, path, _ := newTestManager(t, settingsBase)
	before := fileBytes(t, path)

	if _, err := writeOneKeyAsIs195(m, "fs.delete_enabled", func(base *Config) bool {
		base.FS.DeleteEnabled = true
		return true
	}, nil); err != nil {
		t.Fatalf("the pre-195 body refused the loosening: %v", err)
	}
	if after := fileBytes(t, path); string(after) == string(before) {
		t.Fatalf("the control did not carry a loosening at all: config.toml is unchanged")
	}
	if !readLoaded(t, path).FS.DeleteEnabled {
		t.Errorf("the control's loosening never reached the file")
	}
}

// writeOneKeyAsIs195 is settings.go's pre-195 body, verbatim (HEAD 7eb0e3f4b):
// pre-read, apply, validate, guarded merge write, report. The only thing it does
// not do is ask what direction the change is - which is the gap ticket 195
// closes, and the reason this copy has to stay in the tree as a control.
func writeOneKeyAsIs195(m *Manager, ownedKey string, apply func(base *Config) bool,
	check func(base *Config) error,
) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pre, preExists, err := m.readCurrentFile()
	if err != nil {
		return nil, err
	}
	candidate := deepCopyConfig(pre)
	if !preExists {
		candidate = deepCopyConfig(m.cur)
		if candidate == nil {
			candidate = NewDefaults()
		}
	}
	if check != nil {
		if err := check(candidate); err != nil {
			return nil, err
		}
	}
	if !apply(candidate) {
		return nil, errors.New("config: " + ownedKey + " does not exist in the file this write is based on; nothing was written")
	}
	if err := validate(candidate); err != nil {
		return nil, err
	}

	memOld := deepCopyConfig(m.cur)
	if memOld == nil {
		return nil, errors.New("config: no loaded config in this Manager, refusing to write one")
	}
	apply(m.cur)
	if err := m.mergeWrite(ownedKey, func(base *Config) { apply(base) }); err != nil {
		m.cur = memOld
		return nil, err
	}

	return m.writtenKeyPaths(pre, preExists, ownedKey), nil
}
