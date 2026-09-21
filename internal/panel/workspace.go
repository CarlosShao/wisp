package panel

// The workspace request (ticket 92 AC#3, AC#4).
//
// The composer may show which folder the session is scoped to and ask for a
// different one. It may not decide what that folder resolves to: the whole point
// of the native leg is that C26 gets the vote. This file is that leg's handler,
// and it is deliberately small - it owns no path logic, it calls the one
// canonicalizer the decision chain already uses (internal/tools'
// PathCanonicalizer, which risk.Facts judges through), and it writes the audit
// line for both outcomes.
//
// Two accounts are read here rather than assumed:
//
//   - reparse/junction: propagated from risk.Resolve as risk.ErrReparseDenied
//     (C26's own sentinel, unchanged wording), so the refusal the user sees is
//     the refusal the resolver made. Its implementation is Windows-only today -
//     internal/risk carries a DEFERRED stub for the other platforms, and this
//     ticket neither leans on it nor weakens it.
//   - expansion (ticket 102): a workspace whose %VAR% or ~ spelling moved onto
//     another tree is refused, because selecting a workspace AUTHORIZES a tree.
//     The panel never gets to be the party that quietly renamed what the user
//     pointed at, so res.Actable() is checked a second time on this side of the
//     boundary: a resolver that started answering "yes" to a rewritten path
//     would otherwise become an approval by silence.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// PathScope is the native path surface a workspace switch needs, and nothing
// else. *tools.PathCanonicalizer satisfies it in production; a test can hold it
// to a scripted set of verdicts without building a whole allowlist.
type PathScope interface {
	// WorkspaceRoot reports the current narrowed root ("" when unset).
	WorkspaceRoot() string
	// ResolveWorkspace applies C26 + existence + allowlist containment and
	// returns the full C26 Result so the account is readable here.
	ResolveWorkspace(input string) (risk.Result, error)
	// SetWorkspaceRoot narrows scope to an already-resolved root.
	SetWorkspaceRoot(root string) error
}

// AuditFunc is the sink the runtime uses for its "[audit]" lines. It is a
// function, not a logger type, so the handler stays testable and the assembly
// root keeps one format for every permission event.
type AuditFunc func(format string, args ...any)

// UnsetWorkspaceView renders "no workspace chosen" with a reason, which AC#4
// needs: an input row that shows nothing invites the user to believe the config
// roots are invisible rather than absent.
func UnsetWorkspaceView() WorkspaceView {
	return WorkspaceView{Reason: "未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"}
}

// WorkspaceViewFromRoot renders the current narrowing for a snapshot.
func WorkspaceViewFromRoot(root string) WorkspaceView {
	if strings.TrimSpace(root) == "" {
		return UnsetWorkspaceView()
	}
	return WorkspaceView{
		Set: true, Canonical: root,
		Reason: "已收窄到该工作区：范围外的路径按 R2 判定",
	}
}

// RequestWorkspaceSwitch performs one workspace switch and audits it.
//
// On any refusal the previous scope stays in force (the narrowing is only ever
// widened by a config edit, never by a failed request), the reason is audited,
// and the returned view describes the scope still in effect - so the panel
// cannot end up showing a folder the runtime is not using.
func RequestWorkspaceSwitch(scope PathScope, input string, audit AuditFunc) (WorkspaceView, error) {
	if scope == nil {
		return UnsetWorkspaceView(), fmt.Errorf("panel: no path scope attached - the workspace request was not applied")
	}
	previous := scope.WorkspaceRoot()
	res, err := scope.ResolveWorkspace(input)
	if err == nil {
		// The account is read here too, on the native side of the boundary, so
		// a rewrite cannot pass on the resolver's word alone.
		if actable, aErr := res.Actable(); aErr != nil {
			err = aErr
		} else if sErr := scope.SetWorkspaceRoot(actable); sErr != nil {
			err = sErr
		}
	}
	if err != nil {
		if audit != nil {
			audit("workspace: SWITCH-REFUSED from=%q requested=%q err=%v result=refused "+
				"detail=%q", previous, input, err,
				"范围未变：仍按上一次生效的工作区（或 [fs] allowed_dirs）判定")
		}
		return currentAfter(scope, previous), fmt.Errorf("panel: workspace switch refused: %w", err)
	}
	if audit != nil {
		audit("workspace: SWITCH from=%q to=%q spelling=%q reparse=%v rewritten=%v rewrites=%v result=ok",
			previous, res.Canonical, res.Spelling, res.Reparse, res.Rewritten,
			strings.Join(res.Rewrites, "+"))
	}
	view := WorkspaceView{
		Set: true, Spelling: res.Spelling, Canonical: res.Canonical,
		Reparse: res.Reparse, Rewritten: res.Rewritten,
		Reason: "已收窄到该工作区：范围外的路径按 R2 判定",
	}
	if res.Rewritten {
		view.Reason = "展开改写了这条路径（" + strings.Join(res.Rewrites, "+") + "），已拒绝授权改写后的树"
	}
	return view, nil
}

// currentAfter reports the view of whatever scope is in force after a failed
// request. A scope that answers inconsistently is surfaced, not smoothed over.
func currentAfter(scope PathScope, previous string) WorkspaceView {
	now := scope.WorkspaceRoot()
	if now != previous {
		return WorkspaceView{Reason: fmt.Sprintf(
			"拒绝后范围仍然变了（%q -> %q），请重启会话", previous, now)}
	}
	return WorkspaceViewFromRoot(now)
}
