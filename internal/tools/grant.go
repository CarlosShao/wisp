package tools

import (
	"context"
	"errors"
	"strings"
)

// ---------------------------------------------------------------------------
// D45 scoped session grants: the READ side (ticket 224)
//
// This file exists for the same reason mode.go does, and mode.go:9-18 states it
// better than I can: the permission chain may read a host decision through an
// injected seam, once per call, and nothing in this package may store it or
// change it. A435 approved exactly this shape ("批一枚只读 seam，照现成
// ModeSource 那形由装配根注入") because the alternative - handing
// internal/tools a *memory.Store - would put a database handle inside the
// enforcement layer, and internal/tools has never had one (perm/approval/tools
// reference memory.Store zero times; the only reason this package imports
// internal/memory at all is the tool_call row shape).
//
// What the seam is NOT: it is not an authority the caller can set. Compare the
// two shapes this repository already refuses:
//
//   - agent.ToolRequest carries no mode and no answer, pinned by
//     TestTicket90ModeCarriesNoAllowAuthority (internal/tools/ticket90_test.go:454)
//     because the CALLER of Execute supplies arguments, never a verdict;
//   - tools.Decision carries no allow/grant/approve field, pinned by the same
//     test at :444.
//
// So the session identity lives on the injected source, which the composition
// root minted, and never on a per-call struct. A caller that wanted a different
// session would have to swap the bridge's GrantSource, which is a construction-
// time act visible in the audit log, not an argument.
// ---------------------------------------------------------------------------

// GrantSource is the read side of a D45 session grant. It is satisfied by
// internal/session's ledger and by nothing in this package.
type GrantSource interface {
	// Covering reports whether the host session this source is bound to already
	// holds a live grant for (tool, every path in paths), and returns the id of
	// the row that covered it so the tool_call line can name it
	// (tool_call.grant_id, SPEC-02 §3).
	//
	// Three obligations, all fail-closed:
	//   - it must answer false, not guess true, when it cannot look (a dead
	//     database, a torn row, a panic);
	//   - it must answer false for a call whose paths it cannot key with, which
	//     includes a call with no paths at all - a grant binds
	//     (工具, 路径模式, 会话) and a pathless call has no middle term;
	//   - it must only ever answer for the session it was minted for. That is
	//     the boundary ticket 224's AC#3/AC#4 exist to keep honest: the way this
	//     method stays correct is not a check it runs but the identity it holds.
	Covering(ctx context.Context, tool string, paths []string) (grantID int64, ok bool)
}

// sessionGrantID asks the seam about one call, recovering from anything it
// throws. It returns 0 (never covered) when no source is wired, which is the
// state every composition predating ticket 224 is in, and the state
// cmd/wisp's assembly is in for a bridge built without a session.
//
// The recover is not decoration: this value decides whether a human gets asked
// before something touches their disk, so a broken source must mean "ask", the
// same way permissionMode's broken source means "the strictest档".
func (b *Bridge) sessionGrantID(ctx context.Context, tool string, paths []string) (id int64) {
	if b.grants == nil {
		return 0
	}
	defer func() {
		if rec := recover(); rec != nil {
			b.log("tools: GRANT-READ-FAILED recovering=%v; 按未授权处理（继续询问）", rec)
			id = 0
		}
	}()
	// Only canonicalized paths are quotable. A path C26 could not resolve comes
	// back from displayPaths with a marker attached (bridge.go:874), and a
	// stored row must never be able to match one: refusing to ask here would be
	// a normalization bypass wearing an authorization story.
	resolved := make([]string, 0, len(paths))
	for _, p := range paths {
		c, err := b.canonical(p)
		if err != nil {
			b.log("tools: GRANT-SKIP tool=%s path=%q 无法经 C26 规范化，不查授权", tool, p)
			return 0
		}
		if strings.TrimSpace(c) == "" {
			return 0
		}
		resolved = append(resolved, c)
	}
	if len(resolved) == 0 {
		return 0
	}
	id, ok := b.grants.Covering(ctx, tool, resolved)
	if !ok {
		return 0
	}
	return id
}

// canonical is the C26 route one raw path takes, exposed as a predicate so the
// grant check cannot invent its own normalization. It shares the seam the
// assessor and displayPaths already use; it adds no second resolver.
func (b *Bridge) canonical(p string) (string, error) {
	if b.paths == nil {
		return "", errNoCanonicalizer
	}
	return b.paths.Canonicalize(p)
}

// errNoCanonicalizer names the state "this bridge was built without C26". A
// bridge in that state judges every path fail-closed as out-of-scope L2
// (Options.Paths, bridge.go:35-38), so it has no canonicalized path for a grant
// row to be keyed on either: no resolver, no match, no silence.
var errNoCanonicalizer = errors.New("tools: 未装配 C26 路径规范化器，会话授权无从比对")
