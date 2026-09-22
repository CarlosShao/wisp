package winsec

import (
	"errors"
	"fmt"
	"io/fs"
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

// RewriteAccounted is what ticket 108's AC#4 adds to the seam, and the reason it
// is an optional capability rather than a new method on C26Resolver: the account
// this reports already exists (internal/risk's Result.Rewritten, ticket 102), it
// simply used to stop at the install point, so winsec could see that an answer
// was a clean spelling but never ask whether it still named the tree the caller
// pointed at (R-103-1's second half). A resolver that cannot answer that question
// may not be installed at all, and ResolvePath refuses to act on its answer even
// if something got the seam past that rule.
//
// ResolveAccounted is the same question Resolve asks, with the ledger attached.
// ResolvePath calls it *instead of* Resolve, never both, so the pipeline runs
// once per seal either way; the two answers agreeing is checked where the two
// legs can be compared (the conformance probe), not on the hot path.
type RewriteAccounted interface {
	ResolveAccounted(input string) (path string, rewritten bool, err error)
}

var (
	resolverMu sync.RWMutex
	resolver   C26Resolver
	// seamLatched records that the seam has been given a resolver, which under
	// ticket 108's AC#1 is a one-way door: see SetPathResolver.
	seamLatched bool
)

// SetPathResolver installs the C26 pipeline. internal/risk calls it from init.
//
// Ticket 103's AC#1 is the reason this function is not a plain assignment, and
// ticket 108's AC#1 is the reason the fallback to nil is gone. Ticket 103 added
// three guards - single use, conformance, loud refusal - and ticket 103's
// acceptance agent walked through the first one by changing the order of calls:
// SetPathResolver(nil) was allowed unconditionally, so "free the seam, then
// install the fake" satisfied every rule the guard had (PROBE P1b: the fake was
// accepted, SealFile returned nil, and an explicit S-1-1-0 grant on somebody
// else's file was stripped - the exact outcome AC#1 claimed to prevent).
//
// So the seam is now one-way, and the only directions it can be turned are the
// ones that can narrow it:
//
//  1. nil is refused once a resolver has ever been installed. The old comment
//     justified the fallback with "nil means the built-in floor, which can
//     refuse and rewrite nothing" - true on its own, and useless in sequence,
//     because freeing the seam is exactly what re-opens it to a fake. A
//     one-time install that can be undone is not one-time;
//  2. single use is a latch, not a occupancy check: a *different* resolver is
//     refused whether the slot currently reads nil or not;
//  3. conformance: the candidate must answer every hostile shape in
//     resolverProbeShapes by refusing it, or by handing back a spelling the floor
//     itself accepts *and* keeping that answer inside the tree the shape names
//     (resolverTreeOwnershipFailure, which is what a "constant" fake cannot do);
//  4. loud refusal with an audit record: a rejected install leaves the incumbent
//     in place and writes an ERROR-level slog record naming the candidate and
//     the reason. A guard that fails quietly is indistinguishable from no guard.
//
// It deliberately does not panic: init() order across packages means a panic
// here would take a binary down over a wiring mistake, whereas a refusal leaves
// it on the floor, which is only ever narrower.
//
// What the guard cannot do is read intent, so it is not the whole story:
// ResolvePath re-runs the floor on whatever the installed resolver answers and
// refuses an answer the resolver accounts as a rewrite, which is what makes a
// fake that somehow got in inert rather than merely unlikely. Nothing in this
// package can undo an install; tests that need to put the seam back reach
// SetSeamForTest in export_test.go, which is not compiled into a binary.
func SetPathResolver(r C26Resolver) {
	// The candidate is exercised *before* the lock is taken: probing a resolver
	// means calling back into code this package does not own, and holding the
	// write lock across that would let a resolver that calls ResolvePath from
	// another goroutine block every seal in the process.
	var reason string
	if r != nil {
		reason = resolverConformanceFailure(r)
	}
	resolverMu.Lock()
	defer resolverMu.Unlock()
	name := func(candidate C26Resolver) string {
		if candidate == nil {
			return "<floor>"
		}
		return fmt.Sprintf("%T", candidate)
	}
	if r == nil {
		if resolver == nil && !seamLatched {
			slog.Info("winsec: sealing path resolver left at the built-in floor",
				"reason", "nothing has ever been installed, so nil names no change")
			return
		}
		slog.Error("winsec: refusing to release the sealing path resolver",
			"installed", name(resolver),
			"reason", "the seam is one-way: falling back to nil would free it for a forged resolver to install, which is ticket 103's probe P1b",
			"consequence", "the incumbent resolver stays in place; nothing was uninstalled")
		return
	}
	if reason != "" {
		slog.Error("winsec: refusing to install a path resolver into the sealing seam",
			"resolver", name(r),
			"reason", reason,
			"consequence", "the incumbent resolver, or the built-in floor, stays in place")
		return
	}
	if seamLatched {
		if incumbent := name(resolver); incumbent != name(r) {
			slog.Error("winsec: refusing to replace the already installed sealing path resolver",
				"installed", incumbent,
				"attempted", name(r),
				"reason", "the seam is single-use and one-way: there is no order of calls that frees it")
			return
		}
		// The identity branch ticket 103 left unpinned (R-103-4: zero coverage).
		// Reinstalling the resolver that is already in place is a no-op that
		// keeps the FIRST instance, which is what makes the seam idempotent in
		// the only direction that cannot weaken it.
		slog.Debug("winsec: sealing path resolver installed again identically",
			"resolver", name(r), "kept", "the first instance")
		return
	}
	resolver = r
	seamLatched = true
	slog.Info("winsec: sealing path resolver installed", "resolver", name(r),
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
// makes it safe to run them from an init()-time setter. The tree-ownership leg
// below is deliberately the same kind of probe: two spellings, no mkdir.
func resolverProbeShapes() []string {
	sep := string(filepath.Separator)
	shapes := []string{resolverProbeRoot() + sep + ".." + sep + "wisp-103-conformance-probe"}
	if runtime.GOOS == "windows" {
		shapes = append(shapes, "wisp103"+sep+"conformance-probe")
	}
	return shapes
}

// resolverProbeRoot is the spelling these probes are built on, and ticket 125's
// AC#2 ruling written where the code has to carry it.
//
// Both probe sites used to take os.TempDir() verbatim. On a machine whose temp
// dir is *spelled* through a symlink - a container with TMPDIR under a linked
// /var, macOS where /tmp and /var are links, any box whose temp path simply
// contains one - that hands the guard a root which the built-in floor refuses
// for its own sake. The candidate then answers the hostile shape honestly (it
// folds the parent pointer and keeps the rest of the spelling the OS gave it),
// the floor refuses that answer, and the guard concludes the *candidate* is the
// problem: it refuses to install C26, and the whole process drops to the floor
// over a spelling nobody chose (measured pre-fix, Linux container,
// `installed=<nil>` with the ERROR line; `risk.c26Pipeline` in the same container
// with a real temp dir). That is the gatekeeper reading the OS's own legitimate
// shape as an attack, and it is the same instrument family as ticket 124's
// harness reds - it only grew on the guard's own leg.
//
// So the probe root is resolved before it is used. That is ticket 119's already
// approved discipline - the layer that asks the OS resolves what the OS answered
// (internal/proc's SealableRoot) - applied to the only OS read this package makes
// about its *own* fixture. It is deliberately not a call into internal/proc:
// envfork.go's boundary note ("nothing in internal/winsec calls it") is the
// package doc ticket 113 AC#6 wrote, and the floor must not start depending on
// the layer above it to keep that boundary checkable.
//
// What this does not change, and what the legs in seam_probe_root_125_other_test.go
// measure instead of assert: no caller's path passes through here; ResolvePath
// still re-runs the floor on every answer the installed resolver gives; a
// spelling nobody resolved is still refused by platformVerifyPlacement; and the
// probes keep their hostility, because "<root>/../<name>" is still a parent
// pointer the floor refuses outright, so a candidate that answers it by passing
// it through is still refused - now on a root that names a real tree on every
// platform rather than on one that depends how the machine spells /tmp.
//
// The direction on failure is "hand back what os.TempDir said", i.e. exactly
// today's behavior. An unresolvable temp dir must never become a skipped probe:
// a probe that does not run is a guard that passes for free.
func resolverProbeRoot() string {
	return resolveProbeRoot(os.TempDir())
}

// resolveProbeRoot walks to the longest prefix of path that the filesystem agrees
// to, asks that prefix for its own real spelling, and re-joins the components that
// do not exist yet unchanged. It is the same walk proc.SealableRoot performs, for
// the same reason, on material this package owns rather than on a caller's root:
// no cleaning, no absolutizing, no case folding, and nothing that could turn a
// refusal into an approval - the answer is only ever fed to a check that can
// still say no.
func resolveProbeRoot(path string) string {
	if path == "" {
		return path
	}
	var missing []string
	cur := path
	for {
		_, err := os.Lstat(cur)
		if err == nil {
			break
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return path // unreadable, not absent: not this guard's to reinterpret
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return path
		}
		missing = append(missing, filepath.Base(cur))
		cur = parent
	}
	real, err := filepath.EvalSymlinks(cur)
	if err != nil {
		return path
	}
	for i := len(missing) - 1; i >= 0; i-- {
		real = filepath.Join(real, missing[i])
	}
	return real
}

// resolverConformanceFailure returns the reason r must not be installed, or ""
// when it may.
func resolverConformanceFailure(r C26Resolver) string {
	if _, ok := r.(RewriteAccounted); !ok {
		if _, isFloor := r.(builtinVerifier); !isFloor {
			return "it cannot answer the tree-ownership question at all: without ResolveAccounted this package would have to seal whatever tree its answers name, which is ticket 103's R-103-1"
		}
	}
	for _, probe := range resolverProbeShapes() {
		out, rewritten, err := resolveAccounted(r, probe)
		if err != nil || rewritten {
			continue // refusing, or honestly accounting for a move, are both jobs
		}
		if out == probe {
			return fmt.Sprintf("it answered the hostile shape %q by passing it through unchanged, which is the bypass this seam exists to close", probe)
		}
		if _, fErr := (builtinVerifier{}).Resolve(out); fErr != nil {
			return fmt.Sprintf("it answered %q with %q, a spelling the built-in floor itself refuses: %v", probe, out, fErr)
		}
	}
	return resolverTreeOwnershipFailure(r)
}

// resolveAccounted asks one resolver the one question, preferring the leg that
// carries ticket 102's account so the pipeline runs once per seal and not twice.
// Something that cannot answer is reported as a move: the only way into the seam
// is to be able to say otherwise (checked by resolverConformanceFailure), so a
// caller reaching this helper with anything else has not earned the benefit of
// the doubt.
func resolveAccounted(r C26Resolver, input string) (string, bool, error) {
	if acc, ok := r.(RewriteAccounted); ok {
		return acc.ResolveAccounted(input)
	}
	if _, isFloor := r.(builtinVerifier); isFloor {
		p, err := r.Resolve(input)
		return p, false, err
	}
	return "", true, fmt.Errorf("%w: %s does not implement RewriteAccounted", ErrUnresolvedPath, resolverLabel(r))
}

// resolverTreeOwnershipFailure is AC#4's install-time leg, and the direct answer
// to R-103-1: ticket 103's probe asked only whether an answer is a *clean
// spelling*, so a resolver that answers every input with one fixed clean
// absolute path walked straight through (P1b's fake). Containment is checked
// relationally instead of against a template, so this stays an opinion about the
// candidate's answers and not a second normalizer:
//
//	ask the candidate for a directory and for something inside that directory.
//	An honest pipeline answers with the inner thing inside the outer thing, or
//	refuses, or says "I moved it". A "constant" fake answers the same string
//	twice, and a tree-moving fake answers something that is not inside the tree
//	it just named for the parent.
//
// Refusal on either probe is acceptable, as it is everywhere in this file.
func resolverTreeOwnershipFailure(r C26Resolver) string {
	sep := string(filepath.Separator)
	parent := resolverProbeRoot()
	return treeOwnershipFailureForPair(r, parent, parent+sep+"wisp-108-tree-ownership-probe")
}

// treeOwnershipFailureForPair is the leg above with the probe pair supplied, so
// ticket 112's instrument can plant the CI runner's own shape - a parent that
// exists and is spelled in its 8.3 short form, and a child that does not exist
// under that spelling - on a machine whose own temp path needs no short name.
//
// Why a second witness is needed at all (measured, run 35595651898): the honest
// pipeline answers an EXISTING object with the filesystem's own real path, which
// expands 8.3, and a MISSING leaf with the caller's lexical spelling, which does
// not. On a box whose profile directory is longer than eight characters those two
// answers name the same tree in two different spellings
// (C:\Users\RUNNER~1\AppData\Local\Temp and C:\Users\runneradmin\AppData\Local\Temp),
// and the single containment check below read that as a moved seal and refused to
// install C26 - which is how one machine's username took the whole sealing
// resolver down to the built-in floor. Comparing spellings is not something this
// package may do (D22 ban #2 forbids a second normalizer), so instead of
// loosening the comparison the probe now asks the candidate one more question and
// accepts only when the candidate vouches for its own answer:
//
//	either the child's answer sits inside the answer for the child's stated
//	parent (ticket 108's rule, unchanged and still the first thing tried), or
//	the candidate, asked about the parent prefix of the answer it just gave,
//	names that very same tree it already named for the probe parent.
//
// An honest resolver passes the second witness precisely because the two
// spellings are one object, and it has to say so consistently out loud. A
// constant fake is caught earlier, by the equal-answers leg. A tree-moving fake
// is caught here as well: asked what its own answer's parent is, it has to place
// that parent inside the tree it named for the probe parent, and the whole point
// of the fake is that it does not. Refusal or an honest rewrite account on any
// leg still counts as narrow, as it does everywhere in this file; nothing here
// can be satisfied by passing an answer through unchanged, and
// ResolvePath keeps re-running the floor on whatever the installed resolver
// answers, so a candidate that contradicts itself between two calls is still
// refused at the moment it matters rather than at install time only.
func treeOwnershipFailureForPair(r C26Resolver, parent, child string) string {
	childAns, childMoved, cErr := resolveAccounted(r, child)
	if cErr != nil || childMoved {
		return "" // refusing, or honestly accounting for a move, are both narrow
	}
	parentAns, parentMoved, pErr := resolveAccounted(r, parent)
	if pErr != nil || parentMoved {
		return ""
	}
	if childAns == parentAns {
		return fmt.Sprintf("it answered %q and %q with the same single spelling %q, i.e. it does not resolve the tree the caller named at all",
			parent, child, childAns)
	}
	if answerInsideTree(childAns, parentAns) {
		return ""
	}
	anchor := filepath.Dir(childAns)
	if anchor == childAns {
		// Nothing to vouch for: the answer is a root, and a root cannot sit inside
		// the tree the probe named. The original refusal stands.
		return fmt.Sprintf("it answered %q with %q, which is not inside the tree %q it answered for that path's own parent %q: the seam may not be used to move a seal into another tree",
			child, childAns, parentAns, parent)
	}
	anchorAns, anchorMoved, aErr := resolveAccounted(r, anchor)
	if aErr != nil || anchorMoved {
		return ""
	}
	if sameTree(anchorAns, parentAns) || answerInsideTree(anchorAns, parentAns) {
		return ""
	}
	return fmt.Sprintf("it answered %q with %q, which is neither inside the tree %q it answered for that path's own parent %q nor inside the tree that answer names once the candidate is asked about it directly (%q for %q): the seam may not be used to move a seal into another tree",
		child, childAns, parentAns, parent, anchorAns, anchor)
}

// sameTree reports whether two spellings name one object by the platform's own
// component rule, with no normalization performed on either side - the same
// rule answerInsideTree uses, minus the requirement of a deeper path.
//
// Ticket 126 put the volume segment into that rule, which is where it belongs
// and never was. pathComponents starts after filepath.VolumeName on purpose
// (winsec.go:344's comment is the reason: a volume is not an ancestor of
// anything, and Lstat("C:") names whatever directory the process happens to be
// standing in), but this comparison inherited that start point as if it were a
// statement about identity instead of a statement about ancestors. The
// consequence is that C:\store-44440\artifact.txt and
// D:\store-44440\artifact.txt have always been "one tree" here.
//
// Two production faces read that verdict, which is why the fix is in the
// comparison and not in either caller:
//
//   - noticeNamesTree (winsec_windows.go:106) attributes a seal's notice to the
//     callers asking "was my tree reported on?". Measured on this box with two
//     real volumes, a seal that ran on C: produced one notice and that notice
//     was attributed to a never-sealed tree on D: - forward leg green, so this
//     is a misattribution and not an everything-is-unattributed artifact;
//   - the install-time tree-ownership leg below (:325) treats a matching second
//     witness as "narrow", i.e. as grounds to INSTALL the resolver that offered
//     it. Measured with a candidate that answers a probe parent with
//     D:\wisp126-seam\probe-tree and that same tree's own parent with
//     C:\wisp126-seam\probe-tree: before this change the seam admitted it.
//
// That second face is why the change is only allowed in this direction. Both
// comparisons are used exclusively as grounds to PASS something, so narrowing
// them can refuse a candidate that was previously admitted, and can never admit
// one that was previously refused. Nothing here widens a guard, and the
// honest-resolver legs of tickets 112 and 115 - which compare two spellings of
// one object on one volume - keep their verdicts, because a real object's volume
// letter does not change between its 8.3 spelling, its long spelling and its
// mixed-case spelling (measured: the installed pipeline answers
// "c:\Windows\explorer.exe" as "C:\Windows\explorer.exe").
func sameTree(a, b string) bool {
	if !sameVolume(a, b) {
		return false
	}
	compsA, compsB := pathComponents(a), pathComponents(b)
	if len(compsA) != len(compsB) {
		return false
	}
	for i := range compsA {
		if foldSegment(compsA[i]) != foldSegment(compsB[i]) {
			return false
		}
	}
	return true
}

// sameVolume reports whether two spellings name something on the same volume,
// under the platform's own case rule and nothing else. On POSIX
// filepath.VolumeName is the empty string for every input, so this is true for
// every pair and the two comparisons above keep exactly the verdicts they had
// before ticket 126 - the change is Windows-only in effect, not in build tag.
//
// The comparison is by the volume segment as spelled, with no normalization: a
// drive letter and its \\?\ form, or two UNC shares that redirect to one store,
// read as different volumes here. That is the narrow direction, and it is the
// same disagreement platformVerifyPlacement already refuses at the landing floor
// (placement_windows.go:33's extended-length and UNC legs), so nothing that
// reaches a seal can be lost on it.
func sameVolume(a, b string) bool {
	return foldSegment(filepath.VolumeName(a)) == foldSegment(filepath.VolumeName(b))
}

// foldSegment applies the platform's own case rule to one path segment. Windows
// names objects case-insensitively and POSIX does not, which is the rule
// sameTree and answerInsideTree were already applying per segment; it is
// factored out only so the volume segment is folded by the same rule and this
// package does not grow a second one.
func foldSegment(s string) string {
	if os.PathSeparator == '\\' {
		return strings.ToLower(s)
	}
	return s
}

// answerInsideTree reports whether child names something inside parent, by
// component, over the separator set the platform actually uses (pathPieces'
// rule). No normalization is performed on either side, and on Windows the
// comparison is case-insensitive because that is the platform's own rule for
// naming an object - not because a mismatch would be treated as a rewrite.
//
// Ticket 126 added the volume leg for the same reason it added it to sameTree,
// and the acceptance round of ticket 118 measured this function answering the
// same way there: D:\store\artifact is not inside C:\store, it is a different
// tree that happens to be spelled alike below its volume. Both callers of this
// function are the seam guard's grounds to pass a candidate (:311 and :325), so
// this too only ever refuses more.
func answerInsideTree(child, parent string) bool {
	if !sameVolume(child, parent) {
		return false
	}
	childComps, parentComps := pathComponents(child), pathComponents(parent)
	if len(childComps) <= len(parentComps) {
		return false
	}
	for i, want := range parentComps {
		if foldSegment(childComps[i]) != foldSegment(want) {
			return false
		}
	}
	return true
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
	p, moved, err := resolveAccounted(r, input)
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
	// AC#4: ticket 102's account, now readable on this side of the seam. A
	// resolver that says "this answer moved off the tree the caller named" is
	// neither second-guessed nor believed - the seal is refused, the same direction
	// internal/risk's Result.Actable() takes, but refused *here*, so it holds for
	// whatever is installed rather than only for the wiring in internal/risk.
	if moved {
		return ResolvedPath{}, fmt.Errorf("winsec: refusing to seal %s: the installed %s answered %q and accounted that answer as a rewrite off the tree the caller named: %w",
			input, resolverLabel(r), p, ErrUnresolvedPath)
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
	if problem, found := lexicalTraversal(input); found {
		return "", fmt.Errorf("%w: %s %s", ErrUnresolvedPath, input, problem)
	}
	return platformVerifyPlacement(input)
}

// ResolveAccounted is the floor answering AC#4's question. The answer is always
// "not moved", and it is a fact about this implementation rather than a claim it
// makes per call: builtinVerifier has no rewrite branch at all, it hands back the
// input unchanged or refuses (see the note at the top of this file on why the
// floor is not a second PathResolver).
func (b builtinVerifier) ResolveAccounted(input string) (string, bool, error) {
	p, err := b.Resolve(input)
	return p, false, err
}

// lexicalTraversal is the floor's portable shape check, stated over the
// separators the platform actually uses rather than over two independent
// string.Split passes: a run of separators inside the path is an empty segment
// (the spelling Clean would have eaten, which is the thing being refused), and a
// "." or ".." component is a pointer whose target only the filesystem can name.
//
// The separator set is the same rule pathPieces uses, and for the same reason
// (ticket 108's AC#2): on Windows '/' is a separator the OS honours, so
// "C:/a/../b" has to be seen; on POSIX a backslash is an ordinary character in a
// name, so treating it as a separator would refuse the real directory named
// `a\b` for a "." that is only a "." in a string.
func lexicalTraversal(path string) (string, bool) {
	nativeIsBackslash := os.PathSeparator == '\\'
	isSep := func(c byte) bool {
		return c == os.PathSeparator || (nativeIsBackslash && c == '/')
	}
	i := len(filepath.VolumeName(path))
	leading := true
	for i < len(path) {
		if isSep(path[i]) {
			if leading {
				i++
				continue
			}
			// Two separators in a run, not at the start: an empty segment.
			j := i
			for j < len(path) && isSep(path[j]) {
				j++
			}
			if j-i > 1 {
				return fmt.Sprintf("has an empty path segment (index %d)", i), true
			}
			i = j
			continue
		}
		start := i
		for i < len(path) && !isSep(path[i]) {
			i++
		}
		switch seg := path[start:i]; seg {
		case ".", "..":
			// Collapsing ".." lexically is exactly how "seal A" becomes
			// "modify B" behind a link: the filesystem resolves the parent
			// pointer of the object it reached, a string does not.
			return fmt.Sprintf("traverses %q", seg), true
		}
		leading = false
	}
	return "", false
}
