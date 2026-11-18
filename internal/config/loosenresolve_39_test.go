package config

// Ticket 39 AC#2's resolve teeth (write leg 2026-11-18): the queued
// loosening's two endings, both judged on the FILE (read back), never on a
// return value.
//
//   - accepted card  - the real write lands, and exactly one audit line names
//     the card id and the old->new value of every changed key;
//   - refused / timed out card - nothing is written: the file read back still
//     carries the old value, and the audit line says why nothing moved.
//
// Plus the one-card-one-write rule (a second resolve is refused by name) and
// the unknown-id refusal.

import (
	"strings"
	"testing"
)

// acceptQueue39 is a LoosenQueue that hands back the request's own RequestID,
// exactly the contract the production adapter (cmd/wisp) keeps.
type acceptQueue39 struct{ reqs []LoosenRequest }

func (q *acceptQueue39) QueueLoosen(req LoosenRequest) (string, error) {
	q.reqs = append(q.reqs, req)
	return req.RequestID, nil
}

// AC#2, accepted arm: the card's allow turns into a REAL write, and the one
// audit line carries the card id plus old->new for the key it changed.
func Test39AcceptedLoosenCardWritesAndAuditsOldNewWithCardID(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	q := &acceptQueue39{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	before := fileBytes(t, path)
	res, err := m.WriteSection(SectionWrite{
		Section: "plugins", Key: "plugins.zed.enabled",
		Source: SourceNative, Origin: "cli.test",
		Apply: func(c *Config) error {
			entry := c.Plugins.Entries["zed"]
			entry.Enabled = true
			c.Plugins.Entries["zed"] = entry
			return nil
		},
	})
	if err != nil {
		t.Fatalf("queueing the native loosening: %v", err)
	}
	if res.Written || res.PendingID == "" {
		t.Fatalf("the queue step must write nothing and hand out an id: %+v", res)
	}
	if string(fileBytes(t, path)) != string(before) {
		t.Fatal("config.toml changed while the card was still unanswered")
	}
	if reqs := q.reqs; len(reqs) != 1 || reqs[0].RequestID == "" {
		t.Fatalf("the queue did not receive a request id: %+v", q.reqs)
	}
	cardID := q.reqs[0].RequestID

	if _, err := m.ResolveQueuedLoosen(cardID, true); err != nil {
		t.Fatalf("resolving an accepted card: %v", err)
	}
	got := readLoaded(t, path)
	if !got.Plugins.Entries["zed"].Enabled {
		t.Error("an accepted card did not write the loosening: plugins.zed.enabled is still false on disk")
	}
	if !m.Config().Plugins.Entries["zed"].Enabled {
		t.Error("the in-memory config did not move with the accepted write")
	}
	lines := auditLinesContaining(t, audit, "SECTION-WRITE")
	var accepted []string
	for _, l := range lines {
		if strings.Contains(l, "outcome=written") && strings.Contains(l, "card="+cardID) {
			accepted = append(accepted, l)
		}
	}
	if len(accepted) != 1 {
		t.Fatalf("accepted-write SECTION-WRITE lines = %d, want exactly 1 (card=%s):\n%s",
			len(accepted), cardID, strings.Join(lines, "\n"))
	}
	for _, want := range []string{
		"section=plugins", "key=plugins.zed.enabled", "direction=loosen",
		"source=native", "origin=cli.test", "changes=",
		"plugins.zed.enabled:false->true",
	} {
		if !strings.Contains(accepted[0], want) {
			t.Errorf("the accepted write's audit line is missing %q:\n%s", want, accepted[0])
		}
	}
}

// AC#2, refused arm: the card's reject (and, identically, a timeout - the gate
// delivers both as "not allow") writes NOTHING. The file is read back and must
// still carry the old value.
func Test39RefusedLoosenCardKeepsOldValueOnDisk(t *testing.T) {
	m, path, _ := newTestManager(t, pluginsBase215)
	audit := &auditLog{}
	m.LoosenQueue = &acceptQueue39{}
	m.Audit = audit.f

	res, err := m.WriteSection(SectionWrite{
		Section: "plugins", Key: "plugins.zed.enabled",
		Source: SourceNative, Origin: "cli.test",
		Apply: func(c *Config) error {
			entry := c.Plugins.Entries["zed"]
			entry.Enabled = true
			c.Plugins.Entries["zed"] = entry
			return nil
		},
	})
	if err != nil || res.PendingID == "" {
		t.Fatalf("queueing: res=%+v err=%v", res, err)
	}
	cardID := res.PendingID

	if _, err := m.ResolveQueuedLoosen(cardID, false); err != nil {
		t.Fatalf("resolving a refused card must not error: %v", err)
	}
	got := readLoaded(t, path)
	if got.Plugins.Entries["zed"].Enabled {
		t.Error("a REFUSED card still wrote the loosening: plugins.zed.enabled is true on disk")
	}
	if m.Config().Plugins.Entries["zed"].Enabled {
		t.Error("the in-memory config moved on a refused card")
	}
	lines := auditLinesContaining(t, audit, "SECTION-WRITE")
	var denied []string
	for _, l := range lines {
		if strings.Contains(l, "outcome=denied") && strings.Contains(l, "card="+cardID) {
			denied = append(denied, l)
		}
	}
	if len(denied) != 1 {
		t.Fatalf("denied SECTION-WRITE lines = %d, want exactly 1:\n%s", len(denied), strings.Join(lines, "\n"))
	}
	if !strings.Contains(denied[0], "key=plugins.zed.enabled") {
		t.Errorf("the denied line does not name the key:\n%s", denied[0])
	}
}

// One card authorizes one write: the registry entry is spent on the first
// resolve, and a replay of the same card id is refused by name.
func Test39OneCardAuthorizesOneWrite(t *testing.T) {
	m, _, _ := newTestManager(t, pluginsBase215)
	m.LoosenQueue = &acceptQueue39{}
	res, err := m.WriteSection(SectionWrite{
		Section: "plugins", Key: "plugins.zed.enabled",
		Source: SourceNative, Origin: "cli.test",
		Apply: func(c *Config) error {
			entry := c.Plugins.Entries["zed"]
			entry.Enabled = true
			c.Plugins.Entries["zed"] = entry
			return nil
		},
	})
	if err != nil || res.PendingID == "" {
		t.Fatalf("queueing: res=%+v err=%v", res, err)
	}
	if _, err := m.ResolveQueuedLoosen(res.PendingID, true); err != nil {
		t.Fatalf("first resolve: %v", err)
	}
	_, err = m.ResolveQueuedLoosen(res.PendingID, true)
	if err == nil {
		t.Fatal("a second resolve of the same card id succeeded; one card must authorize one write")
	}
	if !strings.Contains(err.Error(), res.PendingID) {
		t.Errorf("the replay refusal does not name the card id: %v", err)
	}
}

// A card id that was never registered (or whose queue failed the hand-off) is
// refused, not defaulted into a write.
func Test39UnknownCardIDIsRefused(t *testing.T) {
	m, _, _ := newTestManager(t, pluginsBase215)
	if _, err := m.ResolveQueuedLoosen("host:config-loosen#404", true); err == nil {
		t.Fatal("resolving an unknown card id succeeded")
	} else if !strings.Contains(err.Error(), "host:config-loosen#404") {
		t.Errorf("the refusal does not name the id: %v", err)
	}
}

// Drift is fail-closed: a key the queued Apply touches can sit on the NEUTRAL
// side at queue time and on the LOOSEN side at resolve time, because the file
// moved while the card hung. The card named only the loosening keys it showed;
// a completion that would loosen more than that is refused and nothing is
// written (the planted drift itself is left exactly as it stands).
func Test39AcceptedCardDoesNotAuthorizeKeysItNeverNamed(t *testing.T) {
	base := settingsBase + `
[plugins]
tier2_enabled = false

[plugins.zed]
enabled = false
capabilities = ["fs.read"]
`
	m, path, _ := newTestManager(t, base)
	q := &acceptQueue39{}
	audit := &auditLog{}
	m.LoosenQueue = q
	m.Audit = audit.f

	// One owned key ("plugins.zed", the whole entry subtree), one Apply that
	// touches enabled AND capabilities. At queue time the file already carries
	// capabilities = ["fs.read"], so the card names exactly one loosening key.
	res, err := m.WriteSection(SectionWrite{
		Section: "plugins", Key: "plugins.zed",
		Source: SourceNative, Origin: "cli.test",
		Apply: func(c *Config) error {
			entry := c.Plugins.Entries["zed"]
			entry.Enabled = true
			entry.Capabilities = []string{"fs.read"}
			c.Plugins.Entries["zed"] = entry
			return nil
		},
	})
	if err != nil || res.PendingID == "" {
		t.Fatalf("queueing: res=%+v err=%v", res, err)
	}
	cardID := q.reqs[0].RequestID
	if len(q.reqs[0].Keys) != 1 || q.reqs[0].Keys[0] != "plugins.zed.enabled" {
		t.Fatalf("the card was raised for %v, want only plugins.zed.enabled", q.reqs[0].Keys)
	}
	// While the card is up, the capabilities line is removed by hand (a
	// tighten the reload path or the operator is free to make). Now the same
	// Apply re-adds the list against a base without it: a loosening the card
	// never named.
	hand := readLoaded(t, path)
	entry := hand.Plugins.Entries["zed"]
	entry.Capabilities = nil
	hand.Plugins.Entries["zed"] = entry
	if err := SaveFile(path, hand); err != nil {
		t.Fatalf("planting the drift: %v", err)
	}

	_, err = m.ResolveQueuedLoosen(cardID, true)
	if err == nil {
		t.Fatal("a completion that loosens keys the card never named was accepted")
	}
	if !strings.Contains(err.Error(), "plugins.zed.capabilities") {
		t.Errorf("the drift refusal does not name the unapproved key: %v", err)
	}
	got := readLoaded(t, path)
	if got.Plugins.Entries["zed"].Enabled {
		t.Error("the drifted completion still wrote plugins.zed.enabled")
	}
	if len(got.Plugins.Entries["zed"].Capabilities) != 0 {
		t.Error("the planted drift itself was disturbed; the refusal must leave the file as it stands")
	}
	if lines := auditLinesContaining(t, audit, "outcome=refused-drift"); len(lines) != 1 {
		t.Errorf("the drift refusal is not audited exactly once: %+v", lines)
	}
}
