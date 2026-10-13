package panel

// Ticket 186's judgements for the one write leg the owner's 09-28 ruling (Q-64)
// asked for: panel.workspace.request -> WorkspaceWriteHandler ->
// RequestWorkspaceSwitch. The scope here is scripted (scriptedScope,
// workspace_test.go), so each case pins one decision THIS file owns: route it
// through the existing C17 name, never re-decide a refusal, audit a refusal that
// never reached the switch leg, and - the standing one - prove by instrument
// that the leg cannot move a tree on disk.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// workspaceEnvelope186 is the raw string the page really posts for this route
// (frontend/src/lib/panel.ts's requestWorkspaceChange builds exactly these
// fields), so the routing case below travels the production reader rather than
// a hand-built struct.
const workspaceEnvelope186 = `{"method":"panel.workspace.request","requestId":"rid-186","source":"panel-composer","path":"D:/work/wisp"}`

// workspaceReq186 is one parsed envelope in the shape bridge.go hands a handler
// (method already checked by the reader, the target on Path).
func workspaceReq186(path string) ComposerRequest {
	return ComposerRequest{Method: MethodWorkspaceRequest, RequestID: "rid-186", Path: path}
}

func TestWorkspaceHandlerAppliesACleanSwitchThroughTheOneLeg(t *testing.T) {
	scope := &scriptedScope{res: risk.Result{
		Spelling: `D:\work\wisp`, Canonical: `D:\work\wisp`, Resolved: true,
	}}
	var log auditLog
	h := &WorkspaceWriteHandler{Scope: scope, Audit: log.f, Actor: "test"}
	if err := h.HandleWorkspaceRequest(context.Background(), workspaceReq186(`D:\work\wisp`)); err != nil {
		t.Fatalf("a clean in-scope switch was refused: %v", err)
	}
	if len(scope.setArgs) != 1 || scope.setArgs[0] != `D:\work\wisp` {
		t.Fatalf("the scope was narrowed to %v, want exactly one call with the resolved root", scope.setArgs)
	}
	// The success line is workspace.go's own SWITCH line. A second line written
	// here would be two records of one fact, which is how the two drift.
	if !strings.Contains(log.joined(), "SWITCH from=") {
		t.Errorf("no SWITCH audit line reached the sink:\n%s", log.joined())
	}
	if strings.Contains(log.joined(), "WORKSPACE-REFUSED") {
		t.Errorf("a successful switch also booked a refusal line:\n%s", log.joined())
	}
}

func TestWorkspaceHandlerRefusesWithoutAScopeAndNamesTheLayer(t *testing.T) {
	var log auditLog
	h := &WorkspaceWriteHandler{Audit: log.f}
	err := h.HandleWorkspaceRequest(context.Background(), workspaceReq186(`D:\work\wisp`))
	if !errors.Is(err, ErrNoWorkspaceScope) {
		t.Fatalf("err = %v, want ErrNoWorkspaceScope - a handler with no scope must refuse, not panic and not succeed", err)
	}
	if !strings.Contains(log.joined(), "WORKSPACE-REFUSED") {
		t.Errorf("the refusal was not audited - a dropped request that leaves no line looks handled:\n%s", log.joined())
	}
}

func TestWorkspaceHandlerRefusesAMisroutedMethod(t *testing.T) {
	scope := &scriptedScope{}
	var log auditLog
	h := &WorkspaceWriteHandler{Scope: scope, Audit: log.f}
	req := workspaceReq186(`D:\work\wisp`)
	req.Method = MethodModeRequest
	if err := h.HandleWorkspaceRequest(context.Background(), req); err == nil {
		t.Fatal("a mode request was answered by the workspace handler")
	}
	if len(scope.setArgs) != 0 {
		t.Errorf("a misrouted request still narrowed the scope: %v", scope.setArgs)
	}
}

// TestWorkspaceHandlerPropagatesC26sRefusalVerbatim is AC#3(iii)'s other half:
// the refusal the user reads is the one the resolver made, and the switch leg's
// own SWITCH-REFUSED line is the only line booked for it.
func TestWorkspaceHandlerPropagatesC26sRefusalVerbatim(t *testing.T) {
	scope := &scriptedScope{err: risk.ErrReparseDenied}
	var log auditLog
	h := &WorkspaceWriteHandler{Scope: scope, Audit: log.f}
	err := h.HandleWorkspaceRequest(context.Background(), workspaceReq186(`D:\work\link`))
	if !errors.Is(err, risk.ErrReparseDenied) {
		t.Fatalf("err = %v, want the resolver's own sentinel to survive the hop", err)
	}
	if !strings.Contains(log.joined(), "SWITCH-REFUSED") {
		t.Errorf("no SWITCH-REFUSED line reached the sink:\n%s", log.joined())
	}
	if strings.Contains(log.joined(), "WORKSPACE-REFUSED") {
		t.Errorf("the refusal was booked twice (the switch leg DID run):\n%s", log.joined())
	}
	if len(scope.setArgs) != 0 {
		t.Errorf("a refused switch narrowed the scope anyway: %v", scope.setArgs)
	}
}

// TestTheWorkspaceRouteIsTheC17NameAndReachesTheHandler drives the real envelope
// through the real reader and router: the capability has exactly one entrance,
// and it is the whitelisted method - not a second channel the page could reach.
func TestTheWorkspaceRouteIsTheC17NameAndReachesTheHandler(t *testing.T) {
	scope := &scriptedScope{res: risk.Result{
		Spelling: `D:/work/wisp`, Canonical: `D:/work/wisp`, Resolved: true,
	}}
	var log auditLog
	d := &ComposerDispatch{
		Workspace: &WorkspaceWriteHandler{Scope: scope, Audit: log.f},
		Audit:     log.f,
	}
	if _, err := d.Handle(context.Background(), workspaceEnvelope186); err != nil {
		t.Fatalf("the whitelisted route did not reach the handler: %v", err)
	}
	if len(scope.setArgs) != 1 {
		t.Fatalf("the routed request did not narrow the scope: %v", scope.setArgs)
	}
	if got := whitelistMethodsFromSource(t, filepath.Join(panelRepoRoot(t), "internal", "panel", "bridge.go")); !slicesContains(got, MethodWorkspaceRequest) {
		t.Errorf("the route the page sends is not in the C17 whitelist read off bridge.go: %v", got)
	}
}

func TestWorkspaceScopeHolderRefusesBeforeTheRuntimeExists(t *testing.T) {
	h := NewWorkspaceScopeHolder()
	if h.Bound() {
		t.Fatal("a fresh holder reports a bound scope")
	}
	if got := h.WorkspaceRoot(); got.Canonical != "" {
		t.Errorf("unbound holder reports %q as the narrowed root", got.Canonical)
	}
	if _, err := h.ResolveWorkspace(`D:\work\wisp`); !errors.Is(err, ErrNoWorkspaceScope) {
		t.Errorf("unbound ResolveWorkspace err = %v, want ErrNoWorkspaceScope", err)
	}
	if err := h.SetWorkspaceRoot(`D:\work\wisp`, risk.Result{}); !errors.Is(err, ErrNoWorkspaceScope) {
		t.Errorf("unbound SetWorkspaceRoot err = %v, want ErrNoWorkspaceScope", err)
	}
	// And through the handler, which is how a user meets it.
	var log auditLog
	handler := &WorkspaceWriteHandler{Scope: h, Audit: log.f}
	if err := handler.HandleWorkspaceRequest(context.Background(), workspaceReq186(`D:\work\wisp`)); !errors.Is(err, ErrNoWorkspaceScope) {
		t.Errorf("handler err = %v, want the unbound holder's own refusal", err)
	}
}

func TestWorkspaceScopeHolderBindsTheRuntimesOwnScopeOnce(t *testing.T) {
	runtimeScope := &scriptedScope{res: risk.Result{
		Spelling: `D:\work\wisp`, Canonical: `D:\work\wisp`, Resolved: true,
	}}
	h := NewWorkspaceScopeHolder()
	if err := h.Bind(runtimeScope); err != nil {
		t.Fatalf("bind the runtime's scope: %v", err)
	}
	if err := h.Bind(runtimeScope); err != nil {
		t.Errorf("re-binding the same scope must be a no-op: %v", err)
	}
	if err := h.Bind(&scriptedScope{}); err == nil {
		t.Error("a second, DIFFERENT scope was accepted: the panel would narrow one tree and the runtime judge with another")
	}
	var log auditLog
	handler := &WorkspaceWriteHandler{Scope: h, Audit: log.f}
	if err := handler.HandleWorkspaceRequest(context.Background(), workspaceReq186(`D:\work\wisp`)); err != nil {
		t.Fatalf("a switch through the holder was refused: %v", err)
	}
	if len(runtimeScope.setArgs) != 1 {
		t.Fatalf("the runtime's own scope was not the one narrowed: %v", runtimeScope.setArgs)
	}
}

// TestTheRealSwitchLegReachesTheSuccessPathWithACleanAccount is the behavioural
// half of the 10-08 19:1x readings (ticket 186's tail): with the REAL resolver
// (tools.PathCanonicalizer, not a script) the handler now reaches workspace.go's
// success path, and the value that path fills for `Rewritten` can only ever be
// false - a rewritten spelling is refused by the resolver long before the
// panel-side second read at workspace.go:120 is asked to confirm it.
func TestTheRealSwitchLegReachesTheSuccessPathWithACleanAccount(t *testing.T) {
	root := mkTree(t, "switch leg tree", map[string]string{".git/HEAD": "ref: refs/heads/dev\n"})
	scope := tools.NewPathCanonicalizer([]string{root}, nil)
	var log auditLog
	h := &WorkspaceWriteHandler{Scope: scope, Audit: log.f, Actor: "test"}

	// ① the fill point's caller count is 1 now (this handler), and it is reached.
	if err := h.HandleWorkspaceRequest(context.Background(), workspaceReq186(root)); err != nil {
		t.Fatalf("the real resolver refused a clean switch to its own allowlisted root: %v", err)
	}
	acc := scope.WorkspaceRoot()
	if acc.Canonical == "" || acc.Rewritten {
		t.Fatalf("after the switch the scope reports %+v, want the clean coordinate with rewritten=false", acc)
	}
	// ③ so the success branch's rewritten account (workspace.go:144) has no
	// producer: the only way to reach it is a rewritten spelling, and
	// ② that spelling is refused by the resolver itself, with C26's own reason.
	// The variable is DEFINED and points at the very root that is in scope, so the
	// refusal cannot be read as "the directory is missing" - it is the account, and
	// nothing else, that stops it.
	t.Setenv("WISP186R1", root)
	if _, err := scope.ResolveWorkspace(`%WISP186R1%`); err == nil {
		t.Error("the real resolver authorised a rewritten spelling: the panel-side Actable gate would then be the only thing standing, and workspace.go:120 would stop being a second read")
	} else if !errors.Is(err, risk.ErrRewrittenPath) {
		t.Errorf("rewritten spelling refused with %v, want risk.ErrRewrittenPath so the two gates agree on the reason", err)
	}
}

// ---------------------------------------------------------------------------
// The AC#3(i) standing instrument: this leg narrows a scope; it cannot move a
// tree. "Applied in place while work is in flight is safe" rests on that bound,
// so the bound is measured on the leg's own sources rather than remembered.

// repoMutatingCallRe matches the shapes a leg would need to actually change
// something on disk: a process spawn (what `git checkout` costs, and the 09-27
// incident's mechanism) or a write. Reads (Stat/ReadFile/ReadDir) are not here -
// the git read face legitimately uses them.
var repoMutatingCallRe = regexp.MustCompile(
	`\bexec\.Command\b|"os/exec"|\bos\.WriteFile\b|\bos\.Create\b|\bos\.OpenFile\b|` +
		`\bos\.Remove\b|\bos\.RemoveAll\b|\bos\.Rename\b|\bos\.Mkdir\b|\bos\.MkdirAll\b|\bos\.Chmod\b`)

// workspaceSwitchLegFiles are the production files a switch request can travel
// through. A hit in any of them is the shape this ticket refuses.
var workspaceSwitchLegFiles = []string{"workspace.go", "workspace_handlers.go", "git.go"}

func repoMutatingHits186(t *testing.T, dir string, files []string) []string {
	t.Helper()
	var hits []string
	for _, name := range files {
		p := filepath.Join(dir, name)
		for _, hit := range reMatchesInFile(repoMutatingCallRe, p) {
			hits = append(hits, hit)
		}
	}
	return hits
}

func TestTheWorkspaceSwitchLegCannotMutateARepo(t *testing.T) {
	legDir := filepath.Join(panelRepoRoot(t), "internal", "panel")
	if hits := repoMutatingHits186(t, legDir, workspaceSwitchLegFiles); len(hits) > 0 {
		t.Errorf("the workspace switch leg now contains a repo-mutating shape (%s): Q-64 asked for \"which tree am I scoped to\", "+
			"and a spawn or a write here is how the 09-27 incident (HANDOVER §4.0u / A328) is reproduced. "+
			"If a future leg really needs it, that is a ruling to write down here, not a line to add quietly", strings.Join(hits, "; "))
	}
	// The positive control: the same scan over a copy of the leg with one
	// planted line must fire. Planted into a copy, never into the tree.
	planted := t.TempDir()
	for _, name := range workspaceSwitchLegFiles {
		data, err := os.ReadFile(filepath.Join(legDir, name))
		if err != nil {
			t.Fatalf("read %s for the planted copy: %v", name, err)
		}
		body := string(data)
		if name == workspaceSwitchLegFiles[0] {
			body += "\nfunc plantedCheckout() { _ = exec.Command(\"git\", \"checkout\", \"dev\") }\n"
		}
		if err := os.WriteFile(filepath.Join(planted, name), []byte(body), 0o600); err != nil {
			t.Fatalf("write the planted copy of %s: %v", name, err)
		}
	}
	if hits := repoMutatingHits186(t, planted, workspaceSwitchLegFiles); len(hits) == 0 {
		t.Error("the scan stayed silent on a copy that calls exec.Command(\"git\", \"checkout\", ...): the instrument is blind, and " +
			"the clean reading above is decoration")
	}
}
