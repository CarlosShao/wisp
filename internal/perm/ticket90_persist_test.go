package perm

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 90, part 2: the persistence boundary R20/M3 drew, and the two AC
// cases that must NEVER be merged into one.
//
//	AC#3(a)  a manually chosen mode survives a restart           (perm: YES)
//	AC#3(b)  a never-touched config reads the default after start (perm: YES)
//	AC#3b    a D45 session grant does NOT survive a restart      (grant: NO)
//
// They are three test functions on purpose. The orchestrator's ruling says a
// single "persistence" case is how a permanent免审通行证 gets in: one green
// checkbox that silently covers two opposite requirements can only ever be
// satisfied by the looser one.

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

func writeConfig(t *testing.T, dir, mode string) string {
	t.Helper()
	c := config.NewDefaults()
	c.Risk.PermissionMode = mode
	path := filepath.Join(dir, "config.toml")
	if err := config.SaveFile(path, c); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	return path
}

func openManager(t *testing.T, path string) *config.Manager {
	t.Helper()
	m, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager(%s): %v", path, err)
	}
	return m
}

// t90tool is a declared-L1 tool: with no facts and no paths its verdict is
// exactly L1, which is the level the modes disagree about.
type t90tool struct{ mu sync.Mutex }

func (*t90tool) Name() string        { return "t90.l1" }
func (*t90tool) Description() string { return "ticket 90 fixture" }
func (*t90tool) Parameters() tools.JSONSchema {
	return tools.JSONSchema(`{"type":"object"}`)
}

func (*t90tool) Execute(context.Context, json.RawMessage, func(string)) (tools.Result, error) {
	return tools.Result{Text: "ran"}, nil
}

type t90gate struct {
	mu        sync.Mutex
	windows   int
	approvals int
}

func (g *t90gate) PendingWindow(context.Context, tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.windows++
	return tools.AnswerTimeout, ""
}

func (g *t90gate) PendingApproval(context.Context, tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	return tools.AnswerAllow, ""
}

func (g *t90gate) counts() (int, int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.windows, g.approvals
}

// askOnce runs one call through a real bridge fed by a real config file and
// returns which gate branch it reached. This is the "改了这个键，行为跟着变"
// witness ticket 83 asks for: nothing here reads the mode directly.
func askOnce(t *testing.T, path string) (int, int) {
	t.Helper()
	mgr := openManager(t, path)
	st, err := New(Options{Manager: mgr})
	if err != nil {
		t.Fatalf("perm.New: %v", err)
	}
	gate := &t90gate{}
	reg := tools.NewRegistry()
	if err := reg.Register(tools.Entry{
		Tool: &t90tool{},
		Decl: tools.Decl{
			Capabilities: []tools.Capability{tools.CapFSRead},
			Needs:        []tools.Capability{tools.CapFSRead},
			Declared:     risk.L1,
			Resident:     true,
			Provider:     tools.KindBuiltin,
		},
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	b := tools.New(tools.Options{Registry: reg, Gate: gate, Modes: st})
	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t90-persist", CorrelationID: "t90-persist", CallID: "c1",
		Name: "t90.l1", Args: json.RawMessage(`{}`),
	}); err != nil {
		t.Fatalf("bridge.Execute: %v", err)
	}
	return gate.counts()
}

// ---------------------------------------------------------------------------
// AC#3(a): the mode crosses the restart boundary
// ---------------------------------------------------------------------------

// TestTicket90ManualSwitchSurvivesRestart is R20/M3 as the owner set it: the
// 档 that was manually picked is the档 the next process starts with. A
// never-manually-changed config is the OTHER test below, not a case here.
func TestTicket90ManualSwitchSurvivesRestart(t *testing.T) {
	dir := sealableTempDir124(t)
	path := writeConfig(t, dir, risk.ModeAskEveryStepName)
	mgr := openManager(t, path)
	var confirmCalls int
	st, err := New(Options{
		Manager: mgr,
		Confirm: func(context.Context, Switch) error { confirmCalls++; return nil },
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := st.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("start mode = %v, want the default ask_every_step", got)
	}
	for _, want := range []risk.Mode{risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		if err := st.Set(context.Background(), want, "test", "owner"); err != nil {
			t.Fatalf("Set(%s): %v", want, err)
		}
		// "Restart": a brand-new Manager over the same file, i.e. the only
		// thing a real process restart has to read from.
		restarted := openManager(t, path)
		if got := restarted.Config().PermissionMode(); got != want {
			t.Errorf("after restart mode = %v, want the manually chosen %v", got, want)
		}
	}
	if confirmCalls != 1 {
		t.Errorf("confirmations asked = %d, want 1 (only the auto_approve switch)", confirmCalls)
	}
	// The persisted form is the config key, and it is a real string in the
	// file - not a runtime-only value that a restart could not see.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(raw), "permission_mode") ||
		!strings.Contains(string(raw), "auto_approve") {
		t.Errorf("config.toml does not carry the chosen档:\n%s", raw)
	}
}

// TestTicket90UntouchedConfigStartsAtTheDefault is AC#3(b): never manually
// changed => a restart reads the first档, and the key is absent-or-default.
func TestTicket90UntouchedConfigStartsAtTheDefault(t *testing.T) {
	dir := sealableTempDir124(t)
	c := config.NewDefaults()
	if c.Risk.PermissionMode != risk.ModeAskEveryStepName {
		t.Fatalf("the schema default is %q, want %q (R20/M2)",
			c.Risk.PermissionMode, risk.ModeAskEveryStepName)
	}
	path := filepath.Join(dir, "config.toml")
	if err := config.SaveFile(path, c); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	for i := 0; i < 3; i++ {
		mgr := openManager(t, path)
		st, err := New(Options{Manager: mgr})
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		if got := st.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Fatalf("cold start #%d mode = %v, want ask_every_step", i+1, got)
		}
		// No switch happened, so no confirmation was ever needed either.
		if _, a := askOnce(t, path); a != 0 {
			t.Fatalf("cold start #%d reached an L2 card", i+1)
		}
	}
}

// ---------------------------------------------------------------------------
// AC#3b: the session grant does NOT cross it - the opposite expectation
// ---------------------------------------------------------------------------

// TestTicket90SessionGrantDoesNotSurviveRestart is AC#3b, and it is a separate
// test from AC#3(a) on purpose (the orchestrator's boundary clause under R20:
// "改模式能跨重启，会话授权不能跨重启，两条各一条用例，不许合并成一条").
//
// It pins the half the ticket does not build (ticket 49 wires the grant's
// lifecycle) with the machinery that does exist: an approval_grant row is
// keyed to a session id, so the next process - which has a new session - finds
// nothing, while the same two restarts above leave the mode in place.
func TestTicket90SessionGrantDoesNotSurviveRestart(t *testing.T) {
	ctx := context.Background()
	dir := sealableTempDir124(t)

	// Same restart pair as AC#3(a): grant vs mode, one mechanism each.
	store, err := memory.Open(dir, memory.WithLogger(slog.New(slog.DiscardHandler)))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	now := time.Now().Unix()
	id, err := store.InsertGrant(ctx, memory.ApprovalGrant{
		Scope:     memory.GrantScopeSession,
		Tool:      "fs.write",
		Pattern:   filepath.Join(dir, "*"),
		SessionID: "session-before-restart",
		CreatedAt: now,
		ExpiresAt: now + 3600, // still live by its own clock: only the SESSION dies
	})
	if err != nil {
		t.Fatalf("InsertGrant: %v", err)
	}
	if id == 0 {
		t.Fatal("InsertGrant returned no id")
	}
	live, err := store.ListGrantsBySession(ctx, "session-before-restart")
	if err != nil || len(live) != 1 {
		t.Fatalf("within the session the grant must be visible: %d rows, err %v", len(live), err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// "Restart": a new process, therefore a new session id, over the same DB.
	reopened, err := memory.Open(dir, memory.WithLogger(slog.New(slog.DiscardHandler)))
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer reopened.Close()
	next, err := reopened.ListGrantsBySession(ctx, "session-after-restart")
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(next) != 0 {
		t.Errorf("a session grant survived the restart (%d rows for the new session)", len(next))
	}
	all, err := reopened.ListGrants(ctx)
	if err != nil {
		t.Fatalf("ListGrants: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("grant rows after restart = %d, want the one old-session row kept for audit", len(all))
	} else if all[0].SessionID != "session-before-restart" {
		t.Errorf("the surviving row belongs to %q, not the dead session", all[0].SessionID)
	}

	// Deliberately NO mode assertion here. R20/M3's boundary clause says the
	// two expectations are two cases, not one: this test fails ONLY if a session
	// grant becomes durable, and TestTicket90ManualSwitchSurvivesRestart fails
	// ONLY if the mode stops being durable. Wiring the contrast into this test
	// would make AC#3b depend on AC#3a's mechanism, which is exactly the merge
	// the ruling forbids ("合并成一条持久化" = 给永久免审通行证开门).
}

// ---------------------------------------------------------------------------
// the ticket-83 rule: the key is not a liar
// ---------------------------------------------------------------------------

// TestTicket90ConfigKeyChangesWhatTheChainAsks is the ticket-83 requirement in
// behavioral form: flipping the key in the file and reloading changes what the
// decision chain asks - and an unparsable value stops the load instead of
// quietly meaning something else.
func TestTicket90ConfigKeyChangesWhatTheChainAsks(t *testing.T) {
	for _, tc := range []struct {
		key         string
		wantWindows int
		wantApprove int
	}{
		{risk.ModeAskEveryStepName, 1, 0},
		{risk.ModeAskHighRiskName, 0, 0},
		{risk.ModeAutoApproveName, 0, 0},
	} {
		dir := sealableTempDir124(t)
		path := writeConfig(t, dir, tc.key)
		w, a := askOnce(t, path)
		if w != tc.wantWindows || a != tc.wantApprove {
			t.Errorf("permission_mode=%q -> windows=%d approvals=%d, want %d/%d",
				tc.key, w, a, tc.wantWindows, tc.wantApprove)
		}
	}
	// A value the reader cannot map must fail the LOAD (票 83: 响亮失败), and the
	// error must name the key path.
	dir := sealableTempDir124(t)
	path := writeConfig(t, dir, risk.ModeAskEveryStepName)
	c := config.NewDefaults()
	c.Risk.PermissionMode = "everything_off"
	if err := config.SaveFile(path, c); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	if _, _, err := config.LoadFile(path, nil); err == nil {
		t.Fatal("an unparsable permission_mode loaded without error; it would then silently mean " +
			"the default, which is the lying-key shape ticket 83 exists to stop")
	} else if !strings.Contains(err.Error(), "risk.permission_mode") {
		t.Errorf("error %q does not name the key path", err)
	}
}

// TestTicket90HandEditLooseningGoesThroughD36 is the boundary between the two
// write channels. R20/M3 persists the mode; it does not weaken D36 (放宽锁定段要
// L2 重新确认): a hand-edited loosening is applied only when the host's
// ConfirmLocked hook approves, and the running mode tracks that verdict - which
// is why Store reads the mode from the live section rather than caching it.
//
// The cold-start half is the contract's, not this ticket's: at startup
// config.toml IS the operator's declaration (SPEC-03 §2), exactly as for every
// other locked key. That is also why R20/M4 puts the confirmation on the SWITCH,
// where a caller is identifiable, and not on the value.
func TestTicket90HandEditLooseningGoesThroughD36(t *testing.T) {
	dir := sealableTempDir124(t)
	path := writeConfig(t, dir, risk.ModeAskEveryStepName)
	mgr := openManager(t, path)
	var seenSection string
	var seenKeys []string
	mgr.ConfirmLocked = func(section string, keys []string) bool {
		seenSection, seenKeys = section, keys
		return true // the host's L2 re-confirmation was granted
	}
	st, err := New(Options{
		Manager: mgr,
		Confirm: func(context.Context, Switch) error { return nil },
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Hand-edit the file to the loosest档.
	hand := config.NewDefaults()
	hand.Risk.PermissionMode = risk.ModeAutoApproveName
	time.Sleep(15 * time.Millisecond)
	if err := config.SaveFile(path, hand); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	time.Sleep(15 * time.Millisecond)
	rep, err := mgr.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	t90dump(t, path, rep)
	var reported bool
	for _, d := range rep.Locked {
		if d.Section == "risk" {
			reported = d.Direction == config.DirLoosen && d.Approved &&
				containsKey(d.Keys, "risk.permission_mode")
		}
	}
	if !reported {
		t.Errorf("the reload did not report risk.permission_mode as an approved loosening: %+v",
			rep.Locked)
	}
	if seenSection != "risk" || !containsKey(seenKeys, "risk.permission_mode") {
		t.Errorf("the D36 hook saw section=%q keys=%v; a mode loosening must go through it",
			seenSection, seenKeys)
	}
	if got := st.PermissionMode(); got != risk.ModeAutoApprove {
		t.Errorf("after an approved loosening the running mode = %v, want auto_approve", got)
	}
	// A denied loosening must NOT move the running mode either: same read,
	// opposite verdict.
	mgr.ConfirmLocked = func(string, []string) bool { return false }
	denied := config.NewDefaults()
	denied.Risk.PermissionMode = risk.ModeAskHighRiskName
	time.Sleep(15 * time.Millisecond)
	if err := config.SaveFile(path, denied); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	time.Sleep(15 * time.Millisecond)
	if _, err := mgr.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload (denied): %v", err)
	}
	// ask_high_risk from auto_approve is a TIGHTENING, so D36 applies it without
	// asking - the point of this half is only that the running mode always
	// equals what the config layer concluded.
	if got := st.PermissionMode(); got != risk.ModeAskHighRisk {
		t.Errorf("running mode = %v, want ask_high_risk (a tightening is hot-applied)", got)
	}
	// And the programmatic switch (which asked its own L2) takes effect at once.
	if err := st.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got := st.PermissionMode(); got != risk.ModeAutoApprove {
		t.Errorf("the confirmed switch did not take effect: %v", got)
	}
}

func t90dump(t *testing.T, path string, rep *config.Report) {
	t.Helper()
	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	t.Logf("t90dump file=%d bytes err=%v report=%+v locked=%+v",
		st.Size(), err, rep, rep.Locked)
}

// containsKey is a tiny membership helper for the audit assertions.
func containsKey(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}
