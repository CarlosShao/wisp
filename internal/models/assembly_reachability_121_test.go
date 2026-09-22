package models

// Ticket 121 AC#3 - the assembly-reachability gate, the one entry in this file
// whose subject is not this package.
//
// WHAT IT IS FOR. This repository has been bitten six times by the same shape:
// a capability lands, its own tests go green, and nothing in the composition
// root ever calls it, so the thing nobody can see is that the capability is not
// in the binary. Ticket 110 had internal/winsec outside every CI step, ticket
// 114's composer request had no production caller, ticket 117's security notice
// had no production listener, and ticket 109's hand-off guard - green inside
// internal/models, with the only non-test call to Manager.VerifyInstalled at
// internal/models/bridge.go:58 - sat behind a Run that had zero production
// callers, which is how it earned R-109-4 and this ticket. Every one of those
// was "somebody remembered to look". This case is the point where looking
// becomes a red test.
//
// WHY THE HOST IS internal/models AND NOT cmd/wisp. cmd/wisp's own test binary
// needs the sherpa DLLs staged (ticket 98's load hole), so its tests are run by
// scripts/wisp-cli-tests.sh on the windows leg only - a gate there has ZERO
// denominator on the ubuntu leg, which is exactly the half that catches an edge
// written into a //go:build windows file. internal/models is in scripts/portable-
// tests.sh's core scope, which the ubuntu leg runs with the strict runner
// (tools/d22scan/runtests.sh: forces -v -count=1, any top-level SKIP is fatal,
// zero PASS and zero FAIL is fatal). So this file has no t.Skip anywhere: an
// unavailable toolchain must be a failure this step reports, never a green
// nobody read.
//
// ---------------------------------------------------------------------------
// THE LIST (verbatim, this is the whole policy this case enforces):
//
//	capabilityPackages = { internal/models, internal/statemachine }
//
// That is two entries out of the thirteen root-module packages that are not in
// cmd/wisp's dependency graph today (measured in
// docs/evidence/s1/121-preflight-deps-reachability.md on 7699ec3, three GOOS
// identical: go list ./... = 33, in-graph = 20, out-of-graph = 13). Of those 13
// only 4 were "real signals a ticket could fix": models, statemachine, audio,
// ball. The other 9 are tautologically red and stay out of the list for good:
// 4 packages that hold nothing but a doc.go (internal/agent/scheduler,
// internal/session, internal/speech, internal/watchdog - no exported symbol to
// reach), 3 package main (cmd/balldebug, cmd/llmrecord, tools/signmodels - a
// main package can never be anybody's dependency), and 2 test fixtures
// (internal/llm/adaptertest, whose exported signatures take *testing.T, and
// internal/llm/golden, which imports net/http/httptest). Pinning all 13 would
// have landed this case permanently red, and a permanently red gate is read as
// noise, which is worse than no gate. internal/audio and internal/ball are the
// two real ones left out on purpose: neither is this ticket's to wire (the
// orchestrator registered both as explicit deferrals, A91(3)), and audio's only
// consumer, internal/speech, is still a doc.go.
//
// WHAT "capability" MEANS HERE, stated because nothing in the code says it.
// "能力类包" is NOT a property derived from this repository - there is no
// marker, no build tag, no PLAN table column the test can read, and the
// pre-flight (§10 item 9) confirmed no document in the repo defines the term.
// This slice IS the definition: a hand-maintained list, curated by the
// orchestrator, one entry per package somebody has actually built a claim
// about. Adding an entry is a deliberate act, and the rule the orchestrator set
// is that each ticket adds itself, at closing time, if and only if its own AC
// claimed a production path; nothing gets swept in here on the theory that it
// looks like a capability. Deleting an entry from this list goes through the
// same person who put it here.
//
// THE PER-GOOS LEG IS NOT DECORATION. capabilityPackages is asserted once per
// GOOS (windows, linux, darwin) because the failure ticket 121 AC#2 could have
// shipped is an edge that exists in a //go:build windows file: the windows leg
// of the graph would contain internal/models, the linux leg would not, and CI
// runs those two legs in separate jobs where nobody reads them as the same
// sentence. The orchestrator made per-GOOS a hard criterion for this case.
//
// THE -e, AND WHAT IT BUYS AND COSTS (read this before "fixing" the rc).
// GOOS=linux go list -deps ./cmd/wisp exits 1 on this repository, and it exits
// 1 for a reason that has nothing to do with any package in this list: the
// third-party github.com/k2-fsa/sherpa-onnx-go-linux has "build constraints
// exclude all Go files" for that GOOS. The instrument therefore has to add -e
// to get a package set at all, on every GOOS, and the command's exit code is
// NOT the criterion here - the parsed set is. That is the honest shape of this
// gate and it is why every invocation below passes -e.
//
// What -e can hide: go list -e keeps going past a broken edge, so a package
// whose only route into the graph ran through a broken edge would not be
// reported as missing - the failure direction is a false red, not a false
// green. The windows leg closes that gap for real, with an instrument that does
// not need -e: windowsGraphMustMatchPlainRun below re-runs the same query
// WITHOUT -e and requires the same set. The linux and darwin legs cannot make
// that claim, because there is no non--e form of the query that succeeds on
// them; on those two GOOS the shape "reachable only through the sherpa edge"
// is NOT tested, and is named here rather than left to a reader who would
// otherwise assume -e and no -e are the same instrument.

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// modulePrefix is this module's import-path prefix, and the filter that turns
// "every package in the transitive closure, stdlib and third party included"
// into "this repository's packages". The 266-line unfiltered count is not the
// denominator (the pre-flight's §1.1 坑登记): filtered is 20 before AC#2, 22
// after.
const modulePrefix = "github.com/CarlosShao/wisp"

// capabilityPackages is AC#3's list. See the block comment above: it is
// hand-maintained, it is intentionally two entries long, and expanding it is a
// decision that belongs to whoever closes a ticket, not to this file.
var capabilityPackages = []string{
	modulePrefix + "/internal/models",
	modulePrefix + "/internal/statemachine",
}

// graphGOOS lists the platforms the assertion runs on. All three are metadata
// queries: listing dependencies for another GOOS needs no cross toolchain and
// no target machine, so this is not a "we ran Go on macOS" claim in disguise.
var graphGOOS = []string{"windows", "linux", "darwin"}

// TestAC3CapabilityPackagesReachTheCompositionRoot is the gate.
func TestAC3CapabilityPackagesReachTheCompositionRoot(t *testing.T) {
	root := moduleRootForTest(t)

	for _, goos := range graphGOOS {
		goos := goos
		t.Run("GOOS="+goos, func(t *testing.T) {
			graph := depsGraphForGOOS(t, root, goos, true)

			// Instrument sanity BEFORE policy: an empty set is the soundest way
			// for a filter to pass nothing and look like a real answer, and it
			// is what a missing toolchain, a read-only build cache or a
			// GOFLAGS surprise all produce. Reported as an instrument failure,
			// never as "every entry of the list is missing", because those two
			// read the same in a log and mean different things to a human.
			if len(graph) == 0 {
				t.Fatalf("GOOS=%s: go list -e -deps ./cmd/wisp returned no %s/* packages at all - "+
					"this is an instrument failure, NOT a statement about the list below. "+
					"Check the go toolchain on PATH, the build cache, and that this test runs "+
					"inside the module (%s).", goos, modulePrefix, root)
			}
			// The matcher must be able to answer "no". Nothing here is a
			// package of this module, so a set that contains it means the
			// filter matched too much and every green below is fake. This is
			// the leg that cannot drift with policy: no future ticket wiring
			// something real can make it red.
			if graph["github.com/CarlosShao/wisp/internal/notARealPackage"] {
				t.Fatalf("GOOS=%s: the graph contains a path that is not a package in this module, "+
					"so the membership test below proves nothing. Parsed set had %d entries.",
					goos, len(graph))
			}

			for _, pkg := range capabilityPackages {
				if !graph[pkg] {
					t.Errorf("%s is NOT in GOOS=%s's go list -deps ./cmd/wisp. That is the R-109-4 "+
						"shape: capability code that exists in the tree and not in the binary. If a "+
						"build tag is the cause, note that this case runs the same query per GOOS on "+
						"purpose (windows-only edges are what it exists to catch). Parsed set: %d of "+
						"this module's packages.", pkg, goos, len(graph))
				}
			}

			if goos == "windows" {
				windowsGraphMustMatchPlainRun(t, root, graph)
			}
		})
	}
}

// windowsGraphMustMatchPlainRun is the -e control leg, and the only place this
// case can answer the pre-flight's limitation that "-e keeps going past a broken
// edge": on the platform where the query succeeds WITHOUT -e, the two sets must
// be identical. A difference means something in cmd/wisp's closure stopped
// resolving, which is precisely the thing -e would have let this file print as
// green.
func windowsGraphMustMatchPlainRun(t *testing.T, root string, withE map[string]bool) {
	t.Helper()
	plainOut, err := goListDeps(root, "windows", false)
	if err != nil {
		t.Errorf("GOOS=windows go list -deps ./cmd/wisp WITHOUT -e failed (%v) - the plain query is "+
			"expected to succeed on windows, so this is not the known third-party sherpa rc=1 case; the "+
			"-e/no-e equivalence check is therefore unavailable and is reported rather than skipped.", err)
		return
	}
	plain := modulePackageSet(plainOut)
	if diff := setDiff(withE, plain); diff != "" {
		t.Errorf("GOOS=windows: the -e set and the plain set differ, so -e is papering over an "+
			"unresolvable edge somewhere in cmd/wisp's closure:\n%s", diff)
	} else {
		t.Logf("GOOS=windows: -e set == plain set (%d of this module's packages), so the -e used on the "+
			"linux/darwin legs is not hiding a broken edge on the one leg where it can be checked", len(withE))
	}
}

// depsGraphForGOOS runs the query for one GOOS and returns this module's
// packages as a set. A failing query is fatal for that GOOS subtest: with -e the
// command is expected to produce a package set on every platform, including
// linux where the third-party sherpa constraint makes the no--e form exit 1.
func depsGraphForGOOS(t *testing.T, root, goos string, useE bool) map[string]bool {
	t.Helper()
	out, err := goListDeps(root, goos, useE)
	if err != nil {
		t.Fatalf("go list -deps ./cmd/wisp (GOOS=%s, -e=%v) failed: %v\noutput: %s",
			goos, useE, err, truncate(out))
	}
	return modulePackageSet(out)
}

// goListDeps resolves the go toolchain, runs the query in the module root for
// one GOOS, and returns stdout. The toolchain is looked up before exec so a
// host without one fails with a sentence about the toolchain instead of a
// generic os/exec error that reads like a policy violation.
func goListDeps(root, goos string, useE bool) (string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("no go toolchain on PATH (%w) - this case shells out to `go list` on "+
			"purpose and cannot be replaced by a static list of packages", err)
	}
	args := []string{"list"}
	if useE {
		args = append(args, "-e")
	}
	args = append(args, "-deps", "./cmd/wisp")

	cmd := exec.Command(goBin, args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOOS="+goos)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%v: %w", strings.Join(args, " "), err)
	}
	// stderr is appended to stdout's text so a caller that fails sees the
	// sherpa constraint message that explains a truncated closure.
	return stdout.String() + stderr.String(), nil
}

// modulePackageSet filters a `go list` dump down to this module's own packages.
// Exact-line matching, not substring: a strings.Contains over the whole dump
// would answer "present" for internal/model-metadata appearing inside some
// third-party path, which is the sort of green this file exists to refuse.
func modulePackageSet(out string) map[string]bool {
	set := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, modulePrefix+"/") {
			continue
		}
		set[line] = true
	}
	return set
}

// setDiff renders what one set has that the other lacks, for the -e control leg.
func setDiff(a, b map[string]bool) string {
	var sb strings.Builder
	for k := range a {
		if !b[k] {
			sb.WriteString("  only in -e set: " + k + "\n")
		}
	}
	for k := range b {
		if !a[k] {
			sb.WriteString("  only in plain set: " + k + "\n")
		}
	}
	return sb.String()
}

// moduleRootForTest walks up from the test's working directory to the module
// root. Go runs a package's tests with that package's directory as the cwd, so
// without this the query would resolve ./cmd/wisp against internal/models.
func moduleRootForTest(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("no go.mod above %s - this case must run inside the wisp module to have a "+
				"dependency graph to read", dir)
		}
		dir = parent
	}
}

func truncate(s string) string {
	if len(s) <= 400 {
		return s
	}
	return s[:400] + "...(truncated)"
}
