package config

// Ticket 215a's write-path teeth (AC#3: 开关跨档位一致，动作进审计).
//
// The claim these cases pin: the plugin switch writes through the SAME
// protected section-write entry every settings field uses (WriteSection,
// which takes the D36 direction verdict before a byte moves), and the toggle
// itself is audited in EVERY outcome - written, queued, refused, unchanged.
// The direction verdict is the security architecture: switching a plugin OFF
// is a tighten and lands immediately; switching one ON is a loosening of the
// locked [plugins] section, refused for every panel-sourced attempt and
// queued to the L2 card for native ones. Every case below asserts the file
// on disk, not just the return value, because a write the file did not
// record is a lie no receipt can fix.

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// pluginsBase215 is settingsBase plus a [plugins] section with one enabled
// entry and one disabled entry: the toggle cases need a row to switch off
// (a real tighten) and a row to try switching on (a real loosening).
const pluginsBase215 = settingsBase + `
[plugins]
tier2_enabled = true

[plugins.clip]
enabled = true
capabilities = ["fs.read"]

[plugins.zed]
enabled = false
`

func auditLinesContaining(t *testing.T, audit *auditLog, marker string) []string {
	t.Helper()
	var out []string
	for _, line := range audit.all() {
		if strings.Contains(line, marker) {
			out = append(out, line)
		}
	}
	return out
}

// AC#3, forward: switching OFF (the tighten direction) writes immediately -
// the file carries enabled = false afterwards, the in-memory config agrees,
// the receipt names the exact key path, and BOTH audit lines landed (the
// toggle's PLUGIN-TOGGLE with outcome=written and the direction verdict's
// SECTION-WRITE with direction=tighten outcome=written).
func Test215PluginToggleOffWritesImmediatelyAndAudits(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.Audit = audit.f

	before := fileBytes(t, path)
	written, err := m.SetPluginEnabled("clip", false, SourcePanel, "panel.config.set")
	if err != nil {
		t.Fatalf("switching a plugin OFF is a tighten and must write: %v", err)
	}
	if want := []string{"plugins.clip.enabled"}; !reflect.DeepEqual(written, want) {
		t.Errorf("written key paths = %v, want %v: the receipt names the file's own key", written, want)
	}
	if m.Config().Plugins.Entries["clip"].Enabled {
		t.Error("the in-memory config still carries the plugin as enabled")
	}
	if got := readLoaded(t, path); got.Plugins.Entries["clip"].Enabled {
		t.Error("config.toml still carries the plugin as enabled after a written toggle")
	}
	// One key, and one key only: an unrelated row must not have moved.
	if got := readLoaded(t, path); got.Plugins.Entries["zed"].Enabled {
		t.Error("an unrelated plugin row moved during the toggle")
	}
	if string(before) == string(fileBytes(t, path)) {
		t.Error("config.toml is byte-identical after a written toggle - nothing landed")
	}

	toggle := auditLinesContaining(t, audit, "PLUGIN-TOGGLE")
	if len(toggle) != 1 {
		t.Fatalf("PLUGIN-TOGGLE audit lines = %d, want exactly 1: %+v", len(toggle), toggle)
	}
	for _, want := range []string{"key=plugins.clip.enabled", "on=false", "source=panel", "origin=panel.config.set", "outcome=written"} {
		if !strings.Contains(toggle[0], want) {
			t.Errorf("the PLUGIN-TOGGLE line does not name %q: %s", want, toggle[0])
		}
	}
	section := auditLinesContaining(t, audit, "SECTION-WRITE")
	if len(section) != 1 || !strings.Contains(section[0], "direction=tighten") || !strings.Contains(section[0], "outcome=written") {
		t.Errorf("the direction verdict line is missing or wrong: %+v", section)
	}
}

// AC#3, the security arm: switching ON from the PANEL is a loosening of the
// locked [plugins] section and is refused structurally - nothing written,
// file byte-identical, and both audit lines say refused (the attempt itself
// must be in the audit trail even when it did nothing).
func Test215PluginToggleOnFromPanelIsRefusedNotWrittenAndAudited(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.Audit = audit.f

	before := fileBytes(t, path)
	written, err := m.SetPluginEnabled("zed", true, SourcePanel, "panel.config.set")
	if err == nil {
		t.Fatalf("a panel-sourced loosening must be refused, not written: %+v", written)
	}
	if !errors.Is(err, ErrPanelLoosenRefused) {
		t.Errorf("the refusal is not the structural panel-loosening one: %v", err)
	}
	if written != nil {
		t.Errorf("a refused write returned a receipt: %v", written)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Errorf("config.toml changed on a refused loosening:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if m.Config().Plugins.Entries["zed"].Enabled {
		t.Errorf("the in-memory config moved on a refused write: zed.enabled = true")
	}

	toggle := auditLinesContaining(t, audit, "PLUGIN-TOGGLE")
	if len(toggle) != 1 || !strings.Contains(toggle[0], "outcome=refused") {
		t.Errorf("the refused attempt is not audited as refused: %+v", toggle)
	}
	section := auditLinesContaining(t, audit, "SECTION-WRITE")
	if len(section) != 1 || !strings.Contains(section[0], "outcome=refused-panel") {
		t.Errorf("the SECTION-WRITE line does not say the panel was refused: %+v", section)
	}
}

// AC#3, the native arm: switching ON from a native surface is a loosening
// that is QUEUED to the L2 card - one request handed to the queue naming the
// section and the key, nothing written, the audit line saying queued.
func Test215PluginToggleOnFromNativeIsQueuedToTheL2Card(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	q := &testLoosenQueue{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	before := fileBytes(t, path)
	written, err := m.SetPluginEnabled("zed", true, SourceNative, "cli.test")
	if err != nil {
		t.Fatalf("a native loosening must queue, not fail: %v", err)
	}
	if written != nil {
		t.Errorf("a queued loosening returned a receipt: %v", written)
	}
	reqs := q.snapshot()
	if len(reqs) != 1 {
		t.Fatalf("queued requests = %d, want 1", len(reqs))
	}
	if reqs[0].Section != "plugins" || !reflect.DeepEqual(reqs[0].Keys, []string{"plugins.zed.enabled"}) {
		t.Errorf("the queued request names the wrong target: %+v", reqs[0])
	}
	if reqs[0].Source != SourceNative {
		t.Errorf("the queued request lost its source: %+v", reqs[0])
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Error("config.toml changed on a queued loosening")
	}

	toggle := auditLinesContaining(t, audit, "PLUGIN-TOGGLE")
	if len(toggle) != 1 || !strings.Contains(toggle[0], "outcome=queued") {
		t.Errorf("the queued attempt is not audited as queued: %+v", toggle)
	}
}

// The fail-closed twin of the queued case: no queue wired, no loosening
// anywhere - refused, file untouched, audited as refused.
func Test215PluginToggleOnNativeWithoutQueueIsRefusedFailClosed(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.Audit = audit.f

	before := fileBytes(t, path)
	_, err := m.SetPluginEnabled("zed", true, SourceNative, "cli.test")
	if !errors.Is(err, ErrLoosenQueueUnwired) {
		t.Errorf("want the fail-closed no-queue refusal, got: %v", err)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Error("config.toml changed on a refused loosening")
	}
	if lines := auditLinesContaining(t, audit, "PLUGIN-TOGGLE"); len(lines) != 1 || !strings.Contains(lines[0], "outcome=refused") {
		t.Errorf("the refused attempt is not audited as refused: %+v", lines)
	}
}

// The house rule: a settings write cannot create a row. An id the config
// does not declare is refused by name, the file is untouched, and the
// attempt is still audited (an attempt that did nothing is a fact, not a
// silence).
func Test215PluginToggleUnknownIDRefusesAndNeverCreatesARow(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.Audit = audit.f

	before := fileBytes(t, path)
	_, err := m.SetPluginEnabled("ghost", false, SourcePanel, "panel.config.set")
	if err == nil {
		t.Fatal("a toggle for an undeclared plugin id must be refused")
	}
	if !strings.Contains(err.Error(), "ghost") {
		t.Errorf("the refusal does not name the id: %v", err)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Error("config.toml changed on a refused unknown-id toggle (a row was created)")
	}
	if got := readLoaded(t, path); len(got.Plugins.Entries) != 2 {
		t.Errorf("the config now declares %d plugin rows, want the original 2: %+v", len(got.Plugins.Entries), got.Plugins.Entries)
	}
	if lines := auditLinesContaining(t, audit, "PLUGIN-TOGGLE"); len(lines) != 1 || !strings.Contains(lines[0], "outcome=refused") {
		t.Errorf("the unknown-id attempt is not audited: %+v", lines)
	}
}

// The neutral case WriteSection says nothing about: toggling to the value
// the file already carries writes nothing, but the toggle attempt itself is
// still audited with outcome=unchanged - "三档都写" includes the no-op.
func Test215PluginToggleNoOpIsAuditedUnchanged(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.Audit = audit.f

	before := fileBytes(t, path)
	written, err := m.SetPluginEnabled("clip", true, SourcePanel, "panel.config.set")
	if err != nil {
		t.Fatalf("a no-op toggle must not fail: %v", err)
	}
	if written != nil {
		t.Errorf("a no-op toggle returned a receipt: %v", written)
	}
	if after := fileBytes(t, path); string(after) != string(before) {
		t.Error("config.toml changed on a no-op toggle")
	}
	if lines := auditLinesContaining(t, audit, "PLUGIN-TOGGLE"); len(lines) != 1 || !strings.Contains(lines[0], "outcome=unchanged") {
		t.Errorf("the no-op attempt is not audited as unchanged: %+v", lines)
	}
}
