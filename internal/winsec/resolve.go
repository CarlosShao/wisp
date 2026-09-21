package winsec

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
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
// Passing nil restores the built-in verifier.
//
// Ticket 103's AC#1 is the reason this function is no longer a plain assignment.
// PROBE A of ticket 94's acceptance installed, from *outside* this package, a
// resolver that answers every hostile spelling with "yes", and winsec then
// silently rewrote somebody else's DACL and reported success. The seam itself is
// not the bug - winsec cannot import risk without closing a cycle (see the note
// at the top of this file) - so the fix is three guards on the way in:
//
//  1. single use: while a resolver is installed it may not be replaced by a
//     different one. Falling back to nil is always allowed, because nil means the
//     built-in floor, which can refuse and rewrite nothing;
//  2. conformance: the candidate must answer every hostile shape in
//     resolverProbeShapes either by refusing it or by handing back a spelling the
//     floor itself accepts. A pass-through rubber stamp can do neither - which is
//     precisely why it is dangerous;
//  3. loud refusal with an audit record: a rejected install leaves the incumbent
//     in place and writes an ERROR-level slog record naming the candidate and the
//     reason. A guard that fails quietly is indistinguishable from no guard.
//
// It deliberately does not panic: init() order across packages means a panic here
// would take a binary down over a wiring mistake, whereas a refusal leaves it on
// the floor, which is only ever narrower.
//
// What the guard cannot do is read intent, so it is not the whole story:
// ResolvePath re-runs the floor on whatever the installed resolver answers, which
// is what makes a fake that somehow got in inert rather than merely unlikely.
func SetPathResolver(r C26Resolver) {
	if r == nil {
		resolverMu.Lock()
		was := fmt.Sprintf("%T", resolver)
		resolver = nil
		resolverMu.Unlock()
		slog.Info("winsec: sealing path resolver reset to the built-in floor", "was", was)
		return
	}
	name := fmt.Sprintf("%T", r)
	if reason := resolverConformanceFailure(r); reason != "" {
		slog.Error("winsec: refusing to install a path resolver into the sealing seam",
			"resolver", name,
			"reason", reason,
			"consequence", "the incumbent resolver, or the built-in floor, stays in place")
		return
	}
	resolverMu.Lock()
	defer resolverMu.Unlock()
	if resolver != nil {
		if incumbent := fmt.Sprintf("%T", resolver); incumbent != name {
			slog.Error("winsec: refusing to replace the already installed sealing path resolver",
				"installed", incumbent,
				"attempted", name,
				"reason", "the seam is single-use: reset it to nil first, the only direction that can narrow")
			return
		}
		slog.Debug("winsec: sealing path resolver installed again identically", "resolver", name)
		return
	}
	resolver = r
	slog.Info("winsec: sealing path resolver installed", "resolver", name,
		"probes_passed", len(resolverProbeShapes()))
}

// resolverProbeShapes are spellings that no resolver which deserves the name may
// answer by passing them straight through, one per way the answer would name a
// different tree than the caller: which object a relative spelling names depends
// on where the process happens to be standing, and which object "parent\.."
// names depends on what the parent is *behind a link*. Both are refused by the
// built-in floor, so the only acceptable answers are a refusal or a rewrite.
//
// They need no fixtures, no privilege and no filesystem mutation, which is what
// makes it safe to run them from an init()-time setter.
func resolverProbeShapes() []string {
	sep := string(filepath.Separator)
	shapes := []string{os.TempDir() + sep + ".." + sep + "wisp-103-conformance-probe"}
	if runtime.GOOS == "windows" {
		shapes = append(shapes, "wisp103"+sep+"conformance-probe")
	}
	return shapes
}

// resolverConformanceFailure returns the reason r must not be installed, or ""
// when it may.
func resolverConformanceFailure(r C26Resolver) string {
	for _, probe := range resolverProbeShapes() {
		out, err := r.Resolve(probe)
		if err != nil {
			continue // refusing a hostile shape is the job, not a failure
		}
		if out == probe {
			return fmt.Sprintf("it answered the hostile shape %q by passing it through unchanged, which is the bypass this seam exists to close", probe)
		}
		if _, fErr := (builtinVerifier{}).Resolve(out); fErr != nil {
			return fmt.Sprintf("it answered %q with %q, a spelling the built-in floor itself refuses: %v", probe, out, fErr)
		}
	}
	return ""
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
// rather than a string - a caller cannot produce the value by hand, it has to
// survive the floor.
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
	// AC#1's second half, and the reason "the seam is guarded" is not load-bearing
	// on the guard alone: whatever is installed, the spelling that reaches
	// os.Mkdir / os.Chmod has to pass the floor *for itself*. A resolver that
	// answers with a path carrying a link in its ancestor chain, a "..", an empty
	// segment or a \\?\ prefix is not describing the tree the caller named, so the
	// seal is refused here rather than taken on the resolver's word.
	if _, fErr := (builtinVerifier{}).Resolve(p); fErr != nil {
		return ResolvedPath{}, fmt.Errorf("winsec: refusing to seal %s: the installed %s answered %q, a spelling the floor itself refuses: %w",
			input, resolverLabel(r), p, fErr)
	}
	return ResolvedPath{path: p}, nil
}

// resolverLabel attributes a refusal to the thing that produced it, so a fake
// answer cannot be reported as this package's own opinion.
func resolverLabel(r C26Resolver) string {
	if _, isFloor := r.(builtinVerifier); isFloor {
		return "built-in floor"
	}
	return fmt.Sprintf("%T", r)
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
