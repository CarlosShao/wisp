package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// web_bridge_test.go — AC row 5 first clause ("Taint-in-query -> L2 (R4
// integration)") and the D30 attack-chain itself: web.fetch result is
// C25-marked by the bridge, and a web.search query carrying a fragment of
// it is escalated to L2 BEFORE the provider dials. Composition is the real
// one from helpers_test.go (real C26, real C19, real C25), plus the ticket-
// 22 family; only the gate is scripted, same posture as every fs suite.

type scriptGate struct {
	mu   sync.Mutex
	seen []Decision
	// approve, when set, allows L2 through.
	approve bool
}

func (g *scriptGate) PendingWindow(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.seen = append(g.seen, d)
	return AnswerTimeout, "" // L1 window expiry means EXECUTE
}

func (g *scriptGate) PendingApproval(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.seen = append(g.seen, d)
	if g.approve {
		return AnswerAllow, ""
	}
	return AnswerReject, "L2 审批未通过（测试脚本）"
}

func (g *scriptGate) levels() []risk.Level {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]risk.Level, 0, len(g.seen))
	for _, d := range g.seen {
		out = append(out, d.Level)
	}
	return out
}

func (g *scriptGate) find(tool string) []Decision {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []Decision
	for _, d := range g.seen {
		if d.Tool == tool {
			out = append(out, d)
		}
	}
	return out
}

// webFixtureBridge wires fs + the ticket-22 family over the local guard and
// the recording first-visit/Launch seams.
func webFixtureBridge(t *testing.T, allowed []string) (*Bridge, *risk.Provenance, *scriptGate, *recShared) {
	t.Helper()
	var paths *PathCanonicalizer
	if len(allowed) > 0 {
		paths = NewPathCanonicalizer(allowed, nil)
	} else {
		paths = NewPathCanonicalizer(nil, nil)
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	gate := &scriptGate{}
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	rs := newRecShared(t, func(d *WebDeps) { d.Paths = paths })
	webEntries := []Entry{
		{Tool: rs.searchTool(nil), Decl: WebSearchDecl()},
		{Tool: rs.fetchTool(), Decl: WebFetchDecl(WebDeps{})},
		{Tool: rs.openTool(), Decl: WebOpenDecl(WebDeps{})},
		{Tool: rs.fileTool(), Decl: FileOpenDecl()},
		{Tool: rs.launchTool(), Decl: AppLaunchDecl()},
	}
	for _, e := range webEntries {
		if err := reg.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	b := New(Options{
		Registry: reg, Paths: paths, Provenance: prov, Gate: gate,
		Logf: func(string, ...any) {},
	})
	return b, prov, gate, rs
}

func bridgeCall(t *testing.T, b *Bridge, task, tool string, params map[string]any) (string, agent.ToolOutcome) {
	t.Helper()
	out, err := b.Execute(context.Background(), reqFor(task, tool, params))
	if err != nil {
		t.Fatalf("bridge fault for %s: %v", tool, err)
	}
	return out.Text, out
}

func reqFor(task, tool string, params map[string]any) agent.ToolRequest {
	b, _ := json.Marshal(params)
	return agent.ToolRequest{TaskID: task, CallID: "call-" + tool, Name: tool, Args: b}
}

func Test22AttackChainFetchThenTaintedSearchEscalates(t *testing.T) {
	dir := t.TempDir()
	secret := filepath.Join(dir, "git-credentials")
	const frag = "ghp_ZEBRATWENTYTWOSECRET99"
	if err := os.WriteFile(secret, []byte("line1\n"+frag+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	b, prov, gate, rs := webFixtureBridge(t, []string{dir})
	fp := &fakeProvider{egress: "api.anysearch.com"}
	rs.w.provider = fp
	const task = "task-22-chain"

	// step 1: fs.read the credential file (L0, marked)
	text, _ := bridgeCall(t, b, task, "fs.read", map[string]any{"path": secret})
	if !strings.Contains(text, frag) {
		t.Fatalf("fs.read did not return the secret: %s", text)
	}
	if len(prov.ScopeTaints(task)) == 0 {
		t.Fatal("fs.read result was not C25-marked")
	}

	// step 2: web.search whose QUERY carries the fragment (the leak door)
	out, err := b.Execute(context.Background(), reqFor(task, "web.search", map[string]any{
		"query": "summarize " + frag,
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("tainted query must NOT pass unsilenced; got: %s", out.Text)
	}
	cards := gate.find("web.search")
	if len(cards) == 0 {
		t.Fatal("no L2 card for the tainted search")
	}
	last := cards[len(cards)-1]
	if last.Level != risk.L2 {
		t.Errorf("verdict level = %v, want L2", last.Level)
	}
	var hitR4 bool
	for _, r := range last.RulesHit {
		if r == risk.R4 {
			hitR4 = true
		}
	}
	if !hitR4 {
		t.Errorf("RulesHit = %v, want to include R4", last.RulesHit)
	}
	if !strings.Contains(last.Reason, "fs.read") || !strings.Contains(last.Reason, "git-credentials") {
		t.Errorf("card must NAME the source (D30/D33 wording), got: %s", last.Reason)
	}
	// the effect never happened: the provider was not called
	if fp.count() != 0 {
		t.Errorf("provider was dialed despite the L2 refusal: %d calls", fp.count())
	}

	// step 3: same query on an UNTAINTED task passes (scope honesty).
	// OpenTask is the composition root's duty at task admission (ticket 12);
	// calling it here keeps the test's view equal to the production shape:
	// a NEVER-OPENED scope fail-closes to R4 while another scope holds
	// taints (provenance.scopeMarks, adversarial report M-1) — and that
	// fail-closed answer is itself asserted by Test22UnopenedScopeFailsClosed.
	b.OpenTask("task-22-clean")
	out2, err := b.Execute(context.Background(), reqFor("task-22-clean", "web.search", map[string]any{
		"query": "summarize " + frag,
	}))
	if err != nil {
		t.Fatal(err)
	}
	if out2.IsError {
		t.Errorf("untainted task should run: %s", out2.Text)
	}
	if fp.count() != 1 {
		t.Errorf("provider calls = %d, want 1", fp.count())
	}
}

func Test22WebFetchResultIsMarkedAndBlocksFollowOn(t *testing.T) {
	dir := t.TempDir()
	b, prov, _, rs := webFixtureBridge(t, []string{dir})
	page := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "PAGEFRAGMENTABCDEFG 4242 payload")
	}))
	const task = "task-22-mark"
	out := execCallViaBridge(t, b, task, "web.fetch", map[string]any{
		"url": "http://page.test:" + portOf(page.URL) + "/p",
	})
	if out.IsError {
		t.Fatalf("local-guard fetch failed: %s", out.Text)
	}
	marks := prov.ScopeTaints(task)
	var found bool
	for _, m := range marks {
		if m.Tool == "web.fetch" && strings.Contains(m.Origin, "page.test") {
			found = true
		}
	}
	if !found {
		t.Fatalf("web.fetch result not marked with URL origin: %+v", marks)
	}
	// follow-on exfil through notify-shaped text (here: web.search query)
	fp := &fakeProvider{egress: "api.anysearch.com"}
	rs.w.provider = fp
	res, err := b.Execute(context.Background(), reqFor(task, "web.search", map[string]any{
		"query": "look up PAGEFRAGMENTABCDEFG",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Errorf("fragment carried by the page must escalate the follow-on search")
	}
	if fp.count() != 0 {
		t.Errorf("provider dialed despite R4")
	}
}

func Test22PrivateTargetRefusedByTheBridgeLayerToo(t *testing.T) {
	// R5 (verdict layer) and the netGuard (connect layer) are two doors.
	// Through the BRIDGE a private literal is L2 before Execute, so the
	// tool does not even run — pin that, plus the direct-tool depth.
	dir := t.TempDir()
	var hits int
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
	}))
	b, _, gate, rs := webFixtureBridge(t, []string{dir})
	out, err := b.Execute(context.Background(), reqFor("task-22-r5", "web.fetch", map[string]any{
		"url": "http://127.0.0.1:" + portOf(srv.URL) + "/x",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("bridge must refuse the private target: %s", out.Text)
	}
	cards := gate.find("web.fetch")
	if len(cards) == 0 || cards[0].Level != risk.L2 {
		t.Errorf("want an L2 (R5) card, got %+v", cards)
	}
	if hits != 0 {
		t.Errorf("R5 refused yet the server was hit %d times", hits)
	}
	// tool-direct depth (no bridge): the guard refuses on its own.
	res := execCall(t, rs.fetchTool(), map[string]any{"url": "http://169.254.169.254/latest/meta-data/"})
	if !res.IsError || !strings.Contains(res.Text, "SSRF") {
		t.Errorf("tool-direct private literal: %s", res.Text)
	}
}

func Test22DeclsRosterMatchesD34Rows(t *testing.T) {
	want := map[string]struct {
		level risk.Level
		caps  []Capability
	}{
		"web.search": {risk.L0, []Capability{CapNet}},
		"web.fetch":  {risk.L0, []Capability{CapNet}},
		"web.open":   {risk.L1, []Capability{CapNet}},
		"file.open":  {risk.L1, []Capability{CapShell}},
		"app.launch": {risk.L1, []Capability{CapShell}},
	}
	got := map[string]Decl{}
	for _, e := range BuiltinWebEntries(WebDeps{}) {
		got[e.Tool.Name()] = e.Decl
	}
	if len(got) != len(want) {
		t.Fatalf("BuiltinWebEntries returned %d tools, want %d", len(got), len(want))
	}
	for name, w := range want {
		d, ok := got[name]
		if !ok {
			t.Errorf("%s not registered", name)
			continue
		}
		if d.Declared != w.level {
			t.Errorf("%s Declared = %v want %v", name, d.Declared, w.level)
		}
		if !sameCaps(d.Capabilities, w.caps) || !sameCaps(d.Needs, w.caps) {
			t.Errorf("%s caps = %v/%v want %v", name, d.Capabilities, d.Needs, w.caps)
		}
		if d.Provider != KindBuiltin {
			t.Errorf("%s provider = %v", name, d.Provider)
		}
	}
}

func sameCaps(a, b []Capability) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// execCallViaBridge is bridgeCall without the text-only signature.
func execCallViaBridge(t *testing.T, b *Bridge, task, tool string, params map[string]any) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(context.Background(), reqFor(task, tool, params))
	if err != nil {
		t.Fatalf("bridge fault %s: %v", tool, err)
	}
	return out
}
