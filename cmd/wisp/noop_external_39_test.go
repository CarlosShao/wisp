package main

// Ticket 39 AC#4's route teeth (2026-11-18), driven through the REAL panel
// settings door (ConfigWriteHandler -> newConfigStore -> config's guarded
// write path), judged on the FILE and on the wire, never on a return value:
//
//	Test39NoopSaveLeavesFileByteIdenticalAndSaysSo
//	    a save whose value the file already carries moves nothing: same
//	    sha256, same mtime, comments and blank lines intact; the receipt
//	    says explicitly that nothing changed and the reload event carries
//	    an empty written list. The positive control re-saves a REAL change
//	    on the same key so the harness is proven able to tell the two apart.
//	Test39ExternalDriftRefusesTheStaleSave
//	    the page holds a fingerprint from its last config.get, the file is
//	    edited externally, and the stale save is refused (refused-drift,
//	    盘上已变) with the disk left byte-identical; a fresh fingerprint
//	    saves fine; the legacy no-fingerprint shape still works unchanged.
//	Test39ExternalEditErrorsCarryRealLinesAndNameTheKeys
//	    an externally broken file answers errors[]{line,message} through the
//	    read face AND the wire payload: 1-based lines, never swallowed to 0
//	    or 1, each row naming its key; the running config stays up.

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// commentedBase39 has hand-written comments and blank lines: the formatting a
// no-op save must never re-lay.
const commentedBase39 = `# 手写注释：no-op 保存之后这行必须还在。
schema_version = 2

[ball]
# 球的直径，手写的注释
size = 90
`

func write39(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func stat39(t *testing.T, path string) time.Time {
	t.Helper()
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	return st.ModTime()
}

// eventWire39 is the reload event's wire shape, parsed strictly off the reply
// line (the frontend's parseConfigReloadEvent lives in TypeScript; this half
// needs the raw keys).
type eventWire39 struct {
	Event    string   `json:"event"`
	Tier     string   `json:"tier"`
	Written  []string `json:"written"`
	Sections []string `json:"sections"`
	BaseSHA  string   `json:"baseSHA"`
}

func event39(t *testing.T, reply string) eventWire39 {
	t.Helper()
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, panel.ConfigReloadEventMarker) {
			var ev eventWire39
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, panel.ConfigReloadEventMarker)), &ev); err != nil {
				t.Fatalf("the reload event line is not JSON: %v\n%s", err, line)
			}
			if ev.Event != "config.reload" {
				t.Fatalf("the event line is not a config.reload event: %s", line)
			}
			return ev
		}
	}
	t.Fatalf("the accepted write carried no %s event line:\n%s", panel.ConfigReloadEventMarker, reply)
	return eventWire39{}
}

func Test39NoopSaveLeavesFileByteIdenticalAndSaysSo(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.toml"
	write39(t, path, commentedBase39)
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	leg := newConfigStore(mgr, nil, nil, "test-39")
	h := &panel.ConfigWriteHandler{Store: leg, Actor: "test-39"}

	// The read face answers the fingerprint the page will echo back.
	view, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	if view.BaseSHA == "" {
		t.Fatal("the read face answered no fingerprint; the drift handshake has no read half")
	}
	base := view.BaseSHA

	before := mustRead39(t, path)
	beforeMtime := stat39(t, path)

	// The no-op save through the real route: ball.size is already 90.
	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-noop-39", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallSize, ConfigValue: "90", ConfigBaseSHA: base,
	}, "")
	if err != nil {
		t.Fatalf("the no-op save was refused: %v", err)
	}
	if !strings.Contains(reply, "本次没有任何键路径发生变化") {
		t.Errorf("the receipt does not say the save changed nothing:\n%s", reply)
	}
	if strings.Contains(reply, "实际落盘的键路径") {
		t.Errorf("a no-op receipt claims key paths landed:\n%s", reply)
	}
	ev := event39(t, reply)
	if len(ev.Written) != 0 {
		t.Errorf("the no-op's reload event lists written = %v, want an empty list", ev.Written)
	}
	if ev.BaseSHA != base {
		t.Errorf("the no-op's reload event carries baseSHA %q, want the unchanged fingerprint %q", ev.BaseSHA, base)
	}

	// The file: byte-identical, mtime unchanged, comments intact.
	after := mustRead39(t, path)
	if string(after) != string(before) {
		t.Errorf("a no-op save re-laid the file:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if !stat39(t, path).Equal(beforeMtime) {
		t.Errorf("mtime moved across a no-op save (%v -> %v)", beforeMtime, stat39(t, path))
	}
	for _, want := range []string{"# 手写注释：no-op 保存之后这行必须还在。", "# 球的直径，手写的注释"} {
		if !strings.Contains(string(after), want) {
			t.Errorf("a no-op save dropped the comment %q:\n%s", want, after)
		}
	}
	now, _ := mgr.FileSHA256()
	if now != base {
		t.Errorf("sha256 moved across a no-op save (%s -> %s)", base, now)
	}

	// The positive control on the SAME key: a real change moves the bytes, so
	// the harness above is proven able to tell a no-op from a write.
	reply2, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-noop-39b", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallSize, ConfigValue: "120", ConfigBaseSHA: base,
	}, "")
	if err != nil {
		t.Fatalf("the real change was refused: %v", err)
	}
	if !strings.Contains(reply2, "ball.size") {
		t.Errorf("a real save's receipt does not name the key it wrote:\n%s", reply2)
	}
	ev2 := event39(t, reply2)
	if len(ev2.Written) != 1 || ev2.Written[0] != "ball.size" {
		t.Errorf("the real change's event written = %v, want [ball.size]", ev2.Written)
	}
	if strings.Contains(string(mustRead39(t, path)), "# 手写注释") {
		// Expected: a REAL write re-lays the file canonically (mergeWrite's
		// documented cost) - this control only proves bytes moved.
		t.Log("the real write re-laid the file canonically; the comments are gone (mergeWrite's documented cost, ticket 226)")
	}
	after2 := mustRead39(t, path)
	if string(after2) == string(before) {
		t.Error("a real value change left the file byte-identical: the harness cannot tell a write from a no-op")
	}
}

func Test39ExternalDriftRefusesTheStaleSave(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.toml"
	write39(t, path, commentedBase39)
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	leg := newConfigStore(mgr, nil, nil, "test-39")
	h := &panel.ConfigWriteHandler{Store: leg, Actor: "test-39"}

	view, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	base := view.BaseSHA

	// The external edit: a hand change behind the page's back (a real value
	// change to a parseable file, plus a hand comment).
	time.Sleep(10 * time.Millisecond)
	external := strings.Replace(commentedBase39, "size = 90", "size = 80", 1)
	external = "# 外部手改的注释行\n" + external
	write39(t, path, external)
	externalBytes := mustRead39(t, path)
	externalMtime := stat39(t, path)
	externalSHA, _ := mgr.FileSHA256()
	if externalSHA == base {
		t.Fatal("the external edit did not move the file's fingerprint; the scenario is empty")
	}

	// The stale save: the page still echoes the OLD fingerprint and the OLD
	// value (90). It must be refused, and the disk must not move.
	_, err = h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-drift-39", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallSize, ConfigValue: "90", ConfigBaseSHA: base,
	}, "")
	if err == nil {
		t.Fatal("a stale save over an externally edited file was accepted - the drift handshake is dead")
	}
	if !strings.Contains(err.Error(), "盘上已变") || !strings.Contains(err.Error(), "refused-drift") {
		t.Errorf("the drift refusal does not say 盘上已变/refused-drift: %v", err)
	}
	if string(mustRead39(t, path)) != string(externalBytes) {
		t.Errorf("the refused save touched the file:\n%s", mustRead39(t, path))
	}
	if !stat39(t, path).Equal(externalMtime) {
		t.Error("the refused save moved the file's mtime")
	}

	// The fresh save: the page re-reads (new fingerprint) and saves the disk's
	// own current value; it goes through.
	view2, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	if view2.BaseSHA == base {
		t.Error("the re-read returned the stale fingerprint")
	}
	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-drift-39b", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallSize, ConfigValue: "80", ConfigBaseSHA: view2.BaseSHA,
	}, "")
	if err != nil {
		t.Fatalf("the fresh save was refused: %v", err)
	}
	ev := event39(t, reply)
	nowSHA, _ := mgr.FileSHA256()
	if ev.BaseSHA != nowSHA {
		t.Errorf("the accepted save's event carries baseSHA %q, want the post-write fingerprint %q", ev.BaseSHA, nowSHA)
	}

	// The legacy shape, pinned as deliberate: a page with NO fingerprint
	// (pre-39 envelope) still writes over a drifted disk - mergeWrite's own
	// guarded merge keeps every FOREIGN key's value, and the route never
	// breaks an older caller. The fingerprint check only arms for pages that
	// hold one.
	reply2, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-drift-39c", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallSize, ConfigValue: "100",
	}, "")
	if err != nil {
		t.Fatalf("the legacy no-fingerprint write was refused: %v", err)
	}
	if !strings.Contains(reply2, "ball.size") {
		t.Errorf("the legacy write's receipt does not name its key:\n%s", reply2)
	}
}

func Test39ExternalEditErrorsCarryRealLinesAndNameTheKeys(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/config.toml"
	// The process loads a GOOD config (NewManager refuses a broken file at
	// construction - the external breakage is a reload-time state), and THEN
	// the file is hand-broken on disk behind its back: exactly the shape the
	// editor page must see while the Manager keeps its last good config.
	write39(t, path, "schema_version = 2\n\n[app]\ntheme = \"dark\"\n")
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	// The unknown key sits on line 5 (1-based) of the file now on disk.
	write39(t, path, "schema_version = 2\n\n[app]\ntheme = \"dark\"\nbogus_key = 1\n")
	leg := newConfigStore(mgr, nil, nil, "test-39")
	h := &panel.ConfigWriteHandler{Store: leg, Actor: "test-39"}

	view, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings over the broken file: %v", err)
	}
	if len(view.LoadErrors) == 0 {
		t.Fatal("the unknown key was swallowed: the read face answered no errors")
	}
	foundUnknown := false
	for _, e := range view.LoadErrors {
		// The tooth: a decoder row's line is never swallowed to 0 or 1, and
		// the unknown key is named.
		if e.Line == 0 || e.Line == 1 {
			t.Errorf("a decoder row lost its line number (line=%d, message=%q)", e.Line, e.Message)
		}
		if e.Line == 5 && strings.Contains(e.Message, "bogus_key") {
			foundUnknown = true
		}
	}
	if !foundUnknown {
		t.Errorf("no error row names bogus_key at line 5: %+v", view.LoadErrors)
	}

	// The wire carries the same rows: the payload's errors[] with line numbers.
	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigGet, RequestID: "rid-err-39", Source: panel.ComposerRequestSource,
	}, "")
	if err != nil {
		t.Fatalf("config.get refused: %v", err)
	}
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, panel.SettingsPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, panel.SettingsPayloadMarker)
		}
	}
	if payloadLine == "" {
		t.Fatalf("config.get carried no payload line:\n%s", reply)
	}
	var payload struct {
		Errors []struct {
			Line    int    `json:"line"`
			Message string `json:"message"`
		} `json:"errors"`
		Sections []struct{} `json:"sections"`
	}
	if err := json.Unmarshal([]byte(payloadLine), &payload); err != nil {
		t.Fatalf("payload is not JSON: %v\n%s", err, payloadLine)
	}
	wireOK := false
	for _, e := range payload.Errors {
		if e.Line == 5 && strings.Contains(e.Message, "bogus_key") {
			wireOK = true
		}
	}
	if !wireOK {
		t.Errorf("the wire's errors[] does not carry bogus_key at line 5: %s", payloadLine)
	}
	if len(payload.Sections) == 0 {
		t.Error("the section tree vanished while the disk file was broken (the Manager keeps its last good config)")
	}

	// The type-mismatch shape: a value of the wrong type points at ITS line.
	// Same rig: good load first, external breakage second.
	dir2 := t.TempDir()
	path2 := dir2 + "/config.toml"
	write39(t, path2, "schema_version = 2\n\n[app]\ntheme = \"dark\"\n")
	mgr2, err := config.NewManager(path2, nil)
	if err != nil {
		t.Fatalf("NewManager (type arm): %v", err)
	}
	// Line 4 (1-based) carries the wrong type.
	write39(t, path2, "schema_version = 2\n\n[ball]\nsize = \"not-a-number\"\n")
	leg2 := newConfigStore(mgr2, nil, nil, "test-39")
	view2, err := leg2.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings over the type-mismatched file: %v", err)
	}
	foundType := false
	for _, e := range view2.LoadErrors {
		if e.Line == 4 {
			foundType = true
		}
		if e.Line == 0 || e.Line == 1 {
			t.Errorf("a decoder row lost its line number (line=%d, message=%q)", e.Line, e.Message)
		}
	}
	if !foundType {
		t.Errorf("no error row points at the type-mismatched line 4: %+v", view2.LoadErrors)
	}
}

func mustRead39(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	return data
}
