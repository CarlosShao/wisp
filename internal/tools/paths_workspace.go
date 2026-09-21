package tools

// The workspace leg of the composer (ticket 92 AC#3).
//
// "Choose a workspace" is a permission input wearing a picker: the tree that
// gets selected is the tree R2 judges every later write against. So the rules
// here are the strict ones, and they are the SAME primitives the config roots
// already go through - this file adds no second resolver:
//
//   - resolution is risk.Resolve, C26's only entry, so a junction/symlink that
//     is not on the exception list returns risk.ErrReparseDenied and the switch
//     is over. That denial's implementation is Windows-only today
//     (risk.pathresolver_other.go carries its own DEFERRED note for macOS and
//     Linux); this ticket neither weakens nor pretends to strengthen it - what
//     it refuses to do is add a second path leg that could dodge it. The
//     tests below therefore gate the reparse case by GOOS and register the
//     skip, and cover the other three refusals on every platform.
//   - ticket 102's account is read before anything is authorized: a spelling
//     whose %VAR% or ~ expansion moved it onto another tree is refused outright
//     here (unlike a config root, where an operator wrote the line and the
//     rewrite is recorded and reported). A request that arrives from a renderer
//     has no such provenance, so "act only on the tree you named" wins.
//   - existence is confirmed with a capability every platform has
//     (res.Resolved on the leg that has it, treeOnDisk otherwise), never by
//     treating "the handle query is a stub here" as "this folder does not
//     exist" - ticket 107's exact failure.
//   - the selected tree must already be inside the [fs] allowed_dirs roots, so
//     a workspace switch can only ever NARROW scope (see InAllowlist below),
//     and an empty allowlist still authorizes nothing.

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// WorkspaceRoot returns the narrowed workspace root, or "" when the session is
// running on the configured roots alone.
func (p *PathCanonicalizer) WorkspaceRoot() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.workspace
}

// ResolveWorkspace runs the three checks a workspace request must pass before
// it may narrow scope, and hands back the whole C26 Result so the caller can
// read the account itself rather than trust a summary.
//
// Every refusal is an error; none of them leaves a workspace set.
func (p *PathCanonicalizer) ResolveWorkspace(input string) (risk.Result, error) {
	if strings.TrimSpace(input) == "" {
		return risk.Result{}, fmt.Errorf("tools: 工作区路径为空")
	}
	res, err := p.resolve(input)
	if err != nil {
		// risk.ErrReparseDenied propagates verbatim: AC#3(ii) is judged on the
		// reason the caller can show the user, and a re-label here would be the
		// first place the product stopped agreeing with C26.
		return res, err
	}
	// The account leg. Actable() is the one reader of ticket 102's book, and
	// authorizing a tree IS acting on it.
	if _, aErr := res.Actable(); aErr != nil {
		return res, aErr
	}
	// Platform-portable existence: res.Resolved speaks only where handle
	// resolution exists, so a false there is silence, not a denial.
	if !res.Resolved && !treeOnDisk(res.Canonical) {
		return res, fmt.Errorf("tools: 工作区 %s 在磁盘上不是一存在的目录", res.Canonical)
	}
	if !p.inRoots(res.Canonical) {
		return res, fmt.Errorf("tools: 工作区 %s 不在 [fs] allowed_dirs 授权的目录内，拒绝收窄", res.Canonical)
	}
	return res, nil
}

// SetWorkspaceRoot narrows subsequent judgements to root, which MUST already
// have come from ResolveWorkspace. A root outside the allowlist is refused
// here as well: this is the last place a caller can still get it wrong, and the
// narrowing must be structurally impossible to widen.
func (p *PathCanonicalizer) SetWorkspaceRoot(root string) error {
	if strings.TrimSpace(root) == "" {
		return fmt.Errorf("tools: 工作区根为空（用 ClearWorkspace 取消收窄）")
	}
	f := foldPath(root)
	if !p.inRoots(f) {
		return fmt.Errorf("tools: 拒绝把工作区设为未授权的目录 %q", root)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.workspace = f
	return nil
}

// ClearWorkspace removes the narrowing (back to the config roots alone).
func (p *PathCanonicalizer) ClearWorkspace() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.workspace = ""
}

// inRoots is the allowlist question WITHOUT the workspace narrowing - the check
// a candidate workspace itself has to pass. Callers pass either a canonical
// path or an already-folded one; foldPath is idempotent.
func (p *PathCanonicalizer) inRoots(canonical string) bool {
	f := foldPath(canonical)
	if f == "" {
		return false
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return rootsContain(p.roots, f)
}
