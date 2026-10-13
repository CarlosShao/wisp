package panel

// Ticket 186 - the native handler for one panel.workspace.request, and the
// late-bound scope it narrows.
//
// WHY THIS FILE EXISTS. Ticket 92 built the composer's workspace leg
// (workspace.go's RequestWorkspaceSwitch: resolve through C26, read ticket 102's
// rewrite account, narrow, audit both outcomes) and ticket 181 built the read
// face (git.go, and the account-aware snapshot producer). What neither built is
// the door: RequestWorkspaceSwitch had ZERO production callers (186-a1 §B, and
// this leg re-ran the ruler before writing - see the probe log), the router's
// Workspace socket (composer_dispatch.go) was declared and left nil, and the
// route the panel already sends - panel.workspace.request, C17's second name at
// bridge.go:43 - therefore ended in ErrNoHandlerAttached. Owner's 09-28 ruling
// (Q-64, verbatim in ledger A368) is the one that says the button has to work, so
// this file is the handler that socket was declared for.
//
// WHAT IT IS NOT, in the order the mistakes would be tempting:
//
//   - It is NOT a second switch leg. The narrowing, the C26 account read and BOTH
//     audit lines stay in workspace.go's RequestWorkspaceSwitch, and this file
//     calls it rather than re-deriving any part of it. AC#3(ii) is exactly this
//     property: the tree the panel ends up showing and the tree the runtime
//     judges with come out of the same function over the same scope object, so
//     there is no way for the packet to say A while the tools read B.
//   - It is NOT a git client. Nothing in the switch path spawns a process, runs
//     `git`, checks out a branch or writes to any tree. What Q-64 asked for is
//     "which tree am I scoped to"; a checkout leg is the shape of the 09-27
//     incident (seven commits landed on the frontend branch, HANDOVER §4.0u /
//     A328), and it is not this ticket's - the panel's git surface stays a READ
//     face plus this one request (AC#5, pinned by git_test.go's
//     TestGitDimensionHasNoModelCallableTool and by
//     TestTheWorkspaceSwitchLegCannotMutateARepo below).
//   - It does NOT decide the refusal. Every sentence a user can read on a
//     refusal is C26's or workspace.go's own (risk.ErrReparseDenied verbatim,
//     ErrRewrittenPath's own wording); a re-labelled refusal would be the first
//     place the product stopped agreeing with the resolver.
//
// AC#3(i) - "what does this request do while work is in flight" - is answered
// IN PLACE rather than by a veto, and the criterion is a bound on what one
// request can do:
//
//	no queue, no refusal-while-busy, no pre-emption. The switch is applied to the
//	scope the moment it arrives, and the reason that is safe is measurable rather
//	than argued: the upper bound of one request is ONE field on the path scope
//	(Scope.SetWorkspaceRoot, a mutex-guarded assignment in
//	internal/tools/paths_workspace.go) plus one audit line. Every path judgement
//	is made per call by R2/R3 against whatever scope is in force at that moment,
//	so a call already judged keeps the targets it was judged with and a later
//	call is judged against the narrower tree - nothing already running is
//	re-judged, and nothing already judged is re-labelled. Widening is
//	structurally impossible (inRoots gates every candidate; ClearWorkspace is
//	internal), so the incident's shape - "the tree under in-flight work was
//	replaced" - cannot be produced here at all.
//
// The standing instrument for that last sentence is
// TestTheWorkspaceSwitchLegCannotMutateARepo (workspace_handlers_test.go): it
// plants a process spawn and a write into a copy of the switch leg and requires
// the scan to go red, so "this leg cannot move a tree" is a reading, not a
// promise. If a future leg DOES need to run `git` (worktree creation, say), that
// instrument is where the ruling has to be written down - not by deleting it.

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ErrNoWorkspaceScope is the refusal for a workspace request that arrived in a
// process where the runtime's path scope was not (or not yet) handed over. It is
// not a malformed request: the envelope was fine and the machine is incomplete,
// and collapsing those two would tell the user their own click was wrong.
var ErrNoWorkspaceScope = errors.New("panel: 工作区范围未接入")

// WorkspaceWriteHandler is the native handler for one panel.workspace.request
// (composer_dispatch.go's Workspace socket).
//
// Scope is the ONE path scope this process judges with - not a copy, not a
// second canonicalizer built from the same config. The resident leg hands it the
// same object its task pipeline gave its tool bridge (see
// WorkspaceScopeHolder below); that identity is what keeps AC#3(ii) true
// across the panel/runtime boundary.
type WorkspaceWriteHandler struct {
	// Scope is the runtime's path scope. nil = every request is refused with a
	// name, never a nil dereference and never a silence.
	Scope PathScope
	// Audit is the project's existing "[audit]" sink (workspace.go's AuditFunc).
	// The switch's own lines are written by RequestWorkspaceSwitch through this
	// sink; this file adds a line only for a refusal that never reached it.
	Audit AuditFunc
	// Actor names who triggered it, for the refusal line. Empty stays "unknown":
	// an anonymous request is auditable as anonymous, which beats inventing one.
	Actor string
	// mu serialises two switches that arrive together (the page's dispatch and a
	// second caller on this host). RequestWorkspaceSwitch reads the previous
	// scope, resolves, then sets: two interleaved calls could otherwise write two
	// audit lines whose "from" values describe each other's results. Serialising
	// makes each line true about the moment it describes.
	mu sync.Mutex
}

// HandleWorkspaceRequest answers one parsed panel.workspace.request
// (bridge.go's MethodWorkspaceRequest). It returns the refusal it made, never a
// swallowed one.
func (h *WorkspaceWriteHandler) HandleWorkspaceRequest(ctx context.Context, req ComposerRequest) error {
	if req.Method != "" && req.Method != MethodWorkspaceRequest {
		err := fmt.Errorf("panel: %q 不是工作区请求，本处理器不回答（requestId=%q）", req.Method, req.RequestID)
		h.record(req, err)
		return err
	}
	if h == nil || h.Scope == nil {
		err := fmt.Errorf("%w: 面板的工作区请求在本机没有可换的范围，requestId=%q 被拒绝", ErrNoWorkspaceScope, req.RequestID)
		h.record(req, err)
		return err
	}
	// The empty path is NOT refused here on purpose: ResolveWorkspace already
	// owns that sentence ("工作区路径为空"), and a second parser at this layer is
	// how a request comes to be judged by two vocabularies. Everything below
	// happens inside the one leg, including both audit lines.
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := RequestWorkspaceSwitch(h.Scope, req.Path, h.Audit); err != nil {
		// Returned verbatim: RequestWorkspaceSwitch already prefixes its own
		// sentence, and a second prefix here would make the reader parse two
		// layers to find one fact.
		return err
	}
	return nil
}

// record writes the one audit line a refusal deserves when the switch leg itself
// was never reached. A nil sink is tolerated (the audit is a sink, not the gate)
// but never turns a refusal into a silence at this level - the error still goes
// back to the caller.
func (h *WorkspaceWriteHandler) record(req ComposerRequest, err error) {
	if h == nil || h.Audit == nil {
		return
	}
	h.Audit("panel: WORKSPACE-REFUSED request=%q method=%q origin=%q actor=%q err=%v detail=%q",
		req.RequestID, req.Method, PanelModeOrigin, h.actor(), err,
		"范围未变：工作区切换腿没有被调用")
}

// actor is the audit's who, defaulted rather than invented.
func (h *WorkspaceWriteHandler) actor() string {
	if h == nil || h.Actor == "" {
		return "unknown"
	}
	return h.Actor
}

// WorkspaceScopeHolder is a PathScope whose concrete scope is bound later, by
// the assembly that owns the task pipeline.
//
// WHY LATE BINDING IS THE HONEST SHAPE HERE, and not a convenience. In the
// resident leg the panel host is assembled at BOOT (resident_windows.go), while
// the one tools.PathCanonicalizer this process judges with is created later, when
// the task pipeline is assembled by the task source. Writing disp.Workspace after
// the fact would be a field write racing the page's own dispatch goroutine, and
// giving the panel its OWN canonicalizer would build a second truth source for
// "which tree are we in" - the shape AC#3(ii) forbids and the shape ledger A481
// refused for the approval gate one ticket earlier. So the handler holds this
// holder from the start, and the task source binds the runtime's scope into it:
// after that line, the tree the panel narrows IS the tree the tools judge with.
//
// UNBOUND IS A REFUSAL WITH A NAME. A process with no task pipeline (the CLI
// `wisp panel-inbound` seam, or a resident boot where the pipeline refused to
// assemble) has no tree to switch to; saying that is the difference between a
// control wired to nothing and a control that answers. It is deliberately not a
// zero-value "no workspace": an unset workspace and an unattached scope are two
// different facts and the panel must be able to read which one it is looking at.
type WorkspaceScopeHolder struct {
	mu    sync.RWMutex
	scope PathScope
}

// NewWorkspaceScopeHolder builds the unbound holder one assembly hands its panel
// chain and later fills with the runtime's scope.
func NewWorkspaceScopeHolder() *WorkspaceScopeHolder { return &WorkspaceScopeHolder{} }

// Bind attaches the runtime's scope. A second bind of a DIFFERENT scope is
// refused rather than silently replacing the tree every later judgement reads;
// re-binding the same scope is a no-op, so an assembly that walks its boot twice
// cannot corrupt the definition.
func (h *WorkspaceScopeHolder) Bind(scope PathScope) error {
	if scope == nil {
		return fmt.Errorf("panel: 拒绝把工作区范围绑定到 nil")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.scope != nil && h.scope != scope {
		return fmt.Errorf("panel: 工作区范围已经绑定过另一枚 scope，拒绝替换（换树必须是收窄，不是换一根管子）")
	}
	h.scope = scope
	return nil
}

// Bound reports whether a runtime scope has been handed over yet.
func (h *WorkspaceScopeHolder) Bound() bool {
	if h == nil {
		return false
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.scope != nil
}

// WorkspaceRoot reports the narrowing in force, or the zero Result while unbound
// (the holder has no tree to report, and the zero Result is this seam's own way
// of saying "nothing was narrowed").
func (h *WorkspaceScopeHolder) WorkspaceRoot() risk.Result {
	if h == nil {
		return risk.Result{}
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.scope == nil {
		return risk.Result{}
	}
	return h.scope.WorkspaceRoot()
}

// ResolveWorkspace refuses by name while unbound: there is no allowlist to
// resolve against, and a resolver invented here would be the second leg this
// file refuses to be.
func (h *WorkspaceScopeHolder) ResolveWorkspace(input string) (risk.Result, error) {
	if h == nil {
		return risk.Result{}, fmt.Errorf("%w: 持有者为 nil", ErrNoWorkspaceScope)
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.scope == nil {
		return risk.Result{}, fmt.Errorf("%w: 本进程还没有装配任务管线，没有可换的工作区范围", ErrNoWorkspaceScope)
	}
	return h.scope.ResolveWorkspace(input)
}

// SetWorkspaceRoot narrows the bound scope. While unbound it refuses with the
// same sentence ResolveWorkspace gives, so no caller can reach a state where a
// narrowing was recorded nowhere.
func (h *WorkspaceScopeHolder) SetWorkspaceRoot(root string, acc risk.Result) error {
	if h == nil {
		return fmt.Errorf("%w: 持有者为 nil", ErrNoWorkspaceScope)
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.scope == nil {
		return fmt.Errorf("%w: 本进程还没有装配任务管线，没有可换的工作区范围", ErrNoWorkspaceScope)
	}
	return h.scope.SetWorkspaceRoot(root, acc)
}
