package winsec

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// Ticket 94's finding is that this package used to decide *which tree to seal*
// from a lexically absolutized string (the filepath.Abs at winsec.go:126 as
// ticket 89 landed it). That is not a style problem: sealing IS the security
// decision, and a lexical pass resolves nothing - not a junction in the middle
// of the path, not an 8.3 short name, not a \\?\ prefix, not a trailing dot or
// space. Each of those turns "I am sealing A" into "I am modifying B", and the
// caller is then told the wide tree is private. That is the class of bug C26's
// PathResolver exists to close (SPEC-06 §4), which is why D22 ban #2 fails the
// build on it.
//
// So the sealing entry points below no longer accept a path the filesystem has
// not agreed to. Two things stand between a caller's string and os.Mkdir:
//
//  1. the installed C26 pipeline (C26Resolver, set by internal/risk at init -
//     see internal/risk/winsec_c26.go), which is the normal case in every
//     production binary and which rewrites the spelling into the OS's own answer;
//  2. when nothing is installed, a built-in verifier that can only REFUSE, never
//     rewrite (builtinVerifier): an absolute, segment-clean, reparse-free spelling
//     passes unchanged and anything else fails loudly.
//
// Why a hook at all, and why the direction is risk -> winsec: internal/winsec
// cannot import internal/risk, because the graph already runs
// risk -> observe -> secret -> winsec (observe/redact.go calls secret.RedactSecret,
// secret/store.go calls winsec.PrivateDirAll), so the reverse edge closes a cycle
// and `go build` refuses. The hook keeps C26 a single implementation - winsec
// asks, risk resolves.
//
// And why the built-in verifier exists rather than a plain refusal when unset:
// a refusal would make the safety guarantee depend on which packages a binary
// happens to link, and internal/secret / internal/memory do not link internal/risk
// at all (measured: their suites go red on 20+ tests). The verifier is therefore
// the floor - it never invents a spelling, so it is not a second normalizer and
// C26 stays the only path normalization entry point (SPEC-06 §4). Linking risk
// only ever makes the answer stronger, never weaker.

// ErrUnresolvedPath reports that a spelling could not be shown to name the tree
// it names, so nothing was created and nothing was sealed. It is the built-in
// verifier's error; when C26 is installed the refusal carries C26's own sentinel
// (risk.ErrReparseDenied) instead, wrapped.
var ErrUnresolvedPath = errors.New("winsec: path is not provably resolved, refusing to seal")

// C26Resolver is the shape of the C26 pipeline as this package needs it: one
// input spelling in, the filesystem's answer out. An error means "this path may
// not be traversed" (SPEC-06 §4 default-deny on reparse points) and every
// sealing entry point propagates it verbatim.
type C26Resolver interface {
	Resolve(input string) (string, error)
}

var (
	resolverMu sync.RWMutex
	resolver   C26Resolver
)

// SetPathResolver installs the C26 pipeline. internal/risk calls it from init.
// Passing nil restores the built-in verifier; there is no way to install a
// resolver that rewrites nothing and refuses nothing, because the built-in one
// always runs the shape and reparse checks below.
func SetPathResolver(r C26Resolver) {
	resolverMu.Lock()
	defer resolverMu.Unlock()
	resolver = r
}

// PathResolverInstalled reports the installed pipeline, or nil when the built-in
// verifier is what a seal will get. Tests use it to prove the wiring is live
// rather than assuming it.
func PathResolverInstalled() C26Resolver {
	resolverMu.RLock()
	defer resolverMu.RUnlock()
	return resolver
}

// ResolvedPath is a path whose spelling has been put in front of the filesystem
// and approved. Its zero value is unusable and nothing outside this package can
// build one that carries a path: the only mint is ResolvePath, which runs C26 or
// the built-in verifier. That is what makes "this was resolved" a type-level fact
// instead of a comment, and it is why the sealing walk below takes ResolvedPath
// rather than a string - filepath.Abs cannot come back without deleting a
// parameter type first.
type ResolvedPath struct {
	path string
}

// String returns the approved spelling.
func (p ResolvedPath) String() string { return p.path }

// IsZero reports the unusable zero value.
func (p ResolvedPath) IsZero() bool { return p.path == "" }

// ResolvePath puts input in front of the filesystem. With C26 installed the
// result is the OS's canonical form (8.3 expanded, reparse targets resolved,
// \\?\ and UNC normalized, reparse traversal denied). Without it the result is
// input unchanged and only if input can be proved to already name that tree.
// Either way an error means nothing was created and nothing was sealed.
func ResolvePath(input string) (ResolvedPath, error) {
	resolverMu.RLock()
	r := resolver
	resolverMu.RUnlock()
	if r == nil {
		r = builtinVerifier{}
	}
	p, err := r.Resolve(input)
	if err != nil {
		// Propagated as-is: errors.Is against the resolver's own sentinel has to
		// keep holding, so a refusal stays attributable to C26 rather than to
		// this package's opinion.
		return ResolvedPath{}, fmt.Errorf("winsec: refusing to seal %s: %w", input, err)
	}
	if p == "" {
		return ResolvedPath{}, fmt.Errorf("winsec: refusing to seal %s: %w: resolver returned an empty path",
			input, ErrUnresolvedPath)
	}
	return ResolvedPath{path: p}, nil
}

// resolveString is ResolvePath for the exported entry points, which keep their
// string signature because the callers that would have to change live in other
// tickets' packages (internal/memory, internal/secret, internal/agent).
func resolveString(input string) (string, error) {
	p, err := ResolvePath(input)
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

// builtinVerifier is the floor that runs when no C26 pipeline is linked. It
// performs no normalization at all - it either passes a spelling through
// untouched or refuses it - which is the difference between this and a second
// PathResolver, and the reason adding it does not violate SPEC-06 §4.
type builtinVerifier struct{}

func (builtinVerifier) Resolve(input string) (string, error) {
	if input == "" {
		return "", ErrUnresolvedPath
	}
	if !filepath.IsAbs(input) {
		return "", fmt.Errorf("%w: %s is not absolute", ErrUnresolvedPath, input)
	}
	for _, sep := range []string{`/`, `\`} {
		if !strings.Contains(input, sep) {
			continue
		}
		for i, seg := range strings.Split(strings.Trim(input, sep), sep) {
			switch seg {
			case "":
				// A doubled separator is a spelling Clean would have eaten, which
				// is the thing being refused here: whatever the OS reads it as,
				// this package will not act on it.
				return "", fmt.Errorf("%w: %s has an empty path segment (index %d)", ErrUnresolvedPath, input, i)
			case ".", "..":
				// Collapsing ".." lexically is exactly how "seal A" becomes
				// "modify B" behind a link: the filesystem resolves the parent
				// pointer of the object it reached, a string does not.
				return "", fmt.Errorf("%w: %s traverses %q", ErrUnresolvedPath, input, seg)
			}
		}
	}
	return platformVerifyPlacement(input)
}
