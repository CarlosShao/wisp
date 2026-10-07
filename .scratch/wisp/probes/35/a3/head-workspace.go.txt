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
//
// Both ends of the root carry C26's ticket 102 account (ticket 181 AC#7): a
// bare string cannot, and a snapshot built from one can only ever answer
// "nothing was rewritten" - which is the constant-false field this ticket
// measured and named. The account type is risk.Result because that is the
// structure the seam already speaks (ResolveWorkspace returns it) and the only
// one internal/panel and internal/tools may both hold without importing each
// other.
type PathScope interface {
	// WorkspaceRoot reports the current narrowed root with its account
	// attached; the zero Result (Canonical == "") means unset.
	WorkspaceRoot() risk.Result
	// ResolveWorkspace applies C26 + existence + allowlist containment and
	// returns the full C26 Result so the account is readable here.
	ResolveWorkspace(input string) (risk.Result, error)
	// SetWorkspaceRoot narrows scope to an already-resolved root, taking the
	// account that describes it.
	SetWorkspaceRoot(root string, acc risk.Result) error
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

// WorkspaceViewFromRoot renders the current narrowing for a snapshot, from the
// account the path scope holds for it (ticket 181 AC#7).
//
// This is the producer side of WorkspaceView.Rewritten. Before this ticket the
// function took a bare string and filled three of the view's six fields, so the
// two account booleans and the user's own spelling were not "read as false",
// they were unreachable - which made every consumer of the account
// (ReadGitForWorkspace, ProjectInstructionLoadRequestFor) a branch no running
// process could take, exactly as 181-r2 measured and 181-r2 could not fix
// because the fix sits on this side of the seam.
//
// The reason sentence differs from the one RequestWorkspaceSwitch writes on a
// rewritten success (workspace.go's own :109-111), and that is deliberate: that
// one narrates a REFUSAL that is still being produced upstream, this one
// narrates a narrowing that is IN FORCE. Two facts, two sentences; folding them
// into one string would make whichever is false the more likely to read well.
func WorkspaceViewFromRoot(res risk.Result) WorkspaceView {
	if strings.TrimSpace(res.Canonical) == "" {
		return UnsetWorkspaceView()
	}
	view := WorkspaceView{
		Set: true, Spelling: res.Spelling, Canonical: res.Canonical,
		Reparse: res.Reparse, Rewritten: res.Rewritten,
		Reason: "已收窄到该工作区：范围外的路径按 R2 判定",
	}
	if res.Rewritten {
		view.Reason = "已收窄到该工作区，但这条路径被 C26 的展开步改写过（" +
			strings.Join(res.Rewrites, "+") + "）：下面的坐标 " + res.Canonical +
			" 是改写后的树，不是使用者命名的那串 " + res.Spelling + "，范围外的路径按 R2 判定"
	}
	if res.Reparse {
		view.Reason += "；这条路径穿过了 [fs] reparse_point_exceptions 上明列的联接点"
	}
	return view
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
		} else if sErr := scope.SetWorkspaceRoot(actable, res); sErr != nil {
			err = sErr
		}
	}
	if err != nil {
		if audit != nil {
			audit("workspace: SWITCH-REFUSED from=%q requested=%q err=%v result=refused "+
				"detail=%q", previous.Canonical, input, err,
				"范围未变：仍按上一次生效的工作区（或 [fs] allowed_dirs）判定")
		}
		return currentAfter(scope, previous), fmt.Errorf("panel: workspace switch refused: %w", err)
	}
	if audit != nil {
		audit("workspace: SWITCH from=%q to=%q spelling=%q reparse=%v rewritten=%v rewrites=%v result=ok",
			previous.Canonical, res.Canonical, res.Spelling, res.Reparse, res.Rewritten,
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
// request. A scope that answers inconsistently is surfaced, not smoothed over -
// and since ticket 181 AC#7 "inconsistently" covers the account as well as the
// coordinate: the same tree coming back with a different rewrite book is a
// scope that is not sure what it authorised, and the panel must not pick the
// answer it likes. risk.Result holds a slice, so it is not comparable with ==
// and an ordinary string compare would silently ignore the account half.
func currentAfter(scope PathScope, previous risk.Result) WorkspaceView {
	now := scope.WorkspaceRoot()
	if !sameC26Answer(previous, now) {
		return WorkspaceView{Reason: fmt.Sprintf(
			"拒绝后范围或其账户仍然变了（%q -> %q，账户 %s -> %s），请重启会话",
			previous.Canonical, now.Canonical,
			accountSummary(previous), accountSummary(now))}
	}
	return WorkspaceViewFromRoot(now)
}

// sameC26Answer compares two answers a path scope gave about its own root,
// coordinate and account together.
func sameC26Answer(a, b risk.Result) bool {
	return a.Canonical == b.Canonical && a.Spelling == b.Spelling &&
		a.Reparse == b.Reparse && a.Rewritten == b.Rewritten &&
		a.Resolved == b.Resolved &&
		strings.Join(a.Rewrites, "+") == strings.Join(b.Rewrites, "+")
}

// accountSummary is the account half of the audit sentence above, spelled out
// rather than dumped as a struct so the line stays readable in a log.
func accountSummary(res risk.Result) string {
	return fmt.Sprintf("reparse=%v rewritten=%v rewrites=%s spelling=%q",
		res.Reparse, res.Rewritten, strings.Join(res.Rewrites, "+"), res.Spelling)
}
