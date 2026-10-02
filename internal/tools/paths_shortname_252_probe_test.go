//go:build windows

package tools

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 252, AC#1 (probe leg 252-p1). READINGS ONLY - this file judges nothing
// and fixes nothing.
//
// The question it answers is narrow: InAllowlist (paths.go:133) demands TWO
// containment checks before it authorizes anything -
//
//	leg 1 (paths.go:140): rootsContain(p.roots, foldPath(canonical))
//	leg 2 (paths.go:152): resolvedForm(canonical) must answer, and its folded
//	                      form must also sit under a root
//
// and the ticket's premise is that a physical path spelled two ways (the long
// name and its 8.3 short alias) cannot satisfy both at once. The premise does
// NOT say which leg answers first, and "obviously both fail" is not a reading
// (ticket 252 AC#1 forbids using it as one), so every case below prints:
//
//	Roots(), foldPath(canonical), leg1, resolvedForm(canonical) as (rf, ok),
//	leg2, the final InAllowlist boolean, and the derived FIRST FAILING LEG.
//
// Consistency is asserted, not assumed: final == leg1 && leg2 (the workspace
// leg at paths.go:159 is inert here because no workspace was narrowed). If that
// identity ever breaks, this file goes red instead of quietly describing a
// shape production does not have.
//
// Short names come from the OS (GetShortPathNameW), never typed by hand, and
// where the machine refuses to produce a distinguishable one the case FAILS
// LOUDLY with its prerequisite spelled out - never skips. That is this
// repository's standing rule for 8.3 fixtures
// (bridge_junction_windows_test.go:68-90, shortNameOf).
//
// No file is created anywhere: every "new file" asked here is a name that does
// not exist yet, which is the exact shape the ticket says is misjudged.

// shortName252 asks the OS for the 8.3 spelling of an EXISTING path. A reply
// that is the same string as the long spelling means this volume gives this
// name no usable alias, and then a probe that needs two spellings measures
// nothing - so it fails instead of passing vacuously.
func shortName252(t *testing.T, long string) string {
	t.Helper()
	p16, err := syscall.UTF16PtrFromString(long)
	if err != nil {
		t.Fatalf("UTF16PtrFromString(%q): %v", long, err)
	}
	n, err := syscall.GetShortPathName(p16, nil, 0)
	if err != nil || n == 0 {
		t.Fatalf("前置条件缺失：%s 拿不到 8.3 短名（GetShortPathNameW: %v, n=%d）。"+
			"需要该卷开启短名生成（查 fsutil 8dot3name query %s，改它要管理员权限）。"+
			"本探针不许 skip，也不许用手工打出来的 ~1 代替真短名。",
			long, err, n, filepath.VolumeName(long))
	}
	buf := make([]uint16, n)
	if n, err = syscall.GetShortPathName(p16, &buf[0], n); err != nil || n == 0 {
		t.Fatalf("GetShortPathNameW(%q) 第二遍: %v (n=%d)", long, err, n)
	}
	short := syscall.UTF16ToString(buf[:n])
	if strings.EqualFold(short, long) {
		t.Fatalf("前置条件缺失：%s 的短名回显 %s 与长名同串（该名本身已合 8.3 或卷未生成别名），"+
			"两形对比无从谈起。本探针不许 skip。", long, short)
	}
	return short
}

// aliasOf252 is the non-fatal sibling used by the environment census: it reports
// what the OS says without demanding a difference, because "this machine gives
// no short alias here" is itself the reading AC#6 needs.
func aliasOf252(t *testing.T, p string) (string, string) {
	t.Helper()
	p16, err := syscall.UTF16PtrFromString(p)
	if err != nil {
		return "", "UTF16PtrFromString: " + err.Error()
	}
	n, err := syscall.GetShortPathName(p16, nil, 0)
	if err != nil || n == 0 {
		return "", "GetShortPathNameW: " + err.Error()
	}
	buf := make([]uint16, n)
	if n, err = syscall.GetShortPathName(p16, &buf[0], n); err != nil || n == 0 {
		return "", "GetShortPathNameW pass2: " + err.Error()
	}
	short := syscall.UTF16ToString(buf[:n])
	if strings.EqualFold(short, p) {
		return short, "same-as-long (no usable alias)"
	}
	return short, "DIFFERENT spelling exists"
}

// legs252 is one full reading of InAllowlist's internals.
type legs252 struct {
	roots    []string
	folded   string
	leg1     bool
	rf       string
	rfOK     bool
	foldedRF string
	leg2     bool
	final    bool
}

// measure252 reproduces the two legs verbatim (same package, same unexported
// helpers) so the boolean can be attributed, then asks the real InAllowlist.
func measure252(pc *PathCanonicalizer, canonical string) legs252 {
	l := legs252{roots: pc.Roots()}
	l.folded = foldPath(canonical)
	l.leg1 = rootsContain(l.roots, l.folded)
	rf, ok := resolvedForm(canonical)
	l.rf, l.rfOK = rf, ok
	if ok {
		l.foldedRF = foldPath(rf)
		l.leg2 = rootsContain(l.roots, l.foldedRF)
	}
	l.final = pc.InAllowlist(canonical)
	return l
}

func (l legs252) firstFailingLeg() string {
	switch {
	case !l.leg1:
		return "LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))"
	case !l.leg2:
		return "LEG 2 resolvedForm (paths.go:152-153)"
	default:
		return "none: both legs hold"
	}
}

func (l legs252) log(t *testing.T, indent string) {
	t.Helper()
	t.Logf("%sRoots()                       : %q", indent, l.roots)
	t.Logf("%sfoldPath(canonical)           : %q", indent, l.folded)
	t.Logf("%sLEG1 rootsContain(roots,folded): %v", indent, l.leg1)
	t.Logf("%sresolvedForm(canonical)        : %q ok=%v", indent, l.rf, l.rfOK)
	t.Logf("%sfoldPath(resolvedForm)         : %q", indent, l.foldedRF)
	t.Logf("%sLEG2 rootsContain(roots,foldRF): %v", indent, l.leg2)
	t.Logf("%sInAllowlist(canonical)         : %v", indent, l.final)
	t.Logf("%sFIRST FAILING LEG              : %s", indent, l.firstFailingLeg())
}

// one252 is a single probe case: how the roots are spelled, what is asked.
type one252 struct {
	name     string
	rootsCfg []string
	asked    string
	// wantFinal is -1 = no expectation (a reading), 0 = expect false,
	// 1 = expect true (a positive control: red means the carrier is broken).
	wantFinal int
}

func run252(t *testing.T, c one252) {
	t.Helper()
	t.Logf("CASE %s", c.name)
	t.Logf("  roots as configured : %q", c.rootsCfg)
	t.Logf("  asked raw           : %q", c.asked)

	pc := NewPathCanonicalizer(c.rootsCfg, nil)
	if u := pc.UnusableRoots(); len(u) != 0 {
		t.Fatalf("carrier broken: NewPathCanonicalizer dropped a configured root: %q", u)
	}
	if len(pc.Roots()) != len(c.rootsCfg) {
		t.Fatalf("carrier broken: %d roots configured, %d effective: %q",
			len(c.rootsCfg), len(pc.Roots()), pc.Roots())
	}

	res, rerr := risk.Resolve(c.asked, nil)
	t.Logf("  risk.Resolve        : err=%v Canonical=%q Resolved=%v Rewritten=%v Rewrites=%q",
		rerr, res.Canonical, res.Resolved, res.Rewritten, res.Rewrites)

	canonical, cerr := pc.Canonicalize(c.asked)
	if cerr != nil {
		t.Fatalf("carrier broken: Canonicalize(%q): %v (no verdict was reached, so no leg can be attributed)", c.asked, cerr)
	}
	t.Logf("  Canonicalize        : %q", canonical)
	if rerr == nil && canonical != res.Canonical {
		t.Logf("  NOTE: Canonicalize and risk.Resolve disagree: %q vs %q", canonical, res.Canonical)
	}

	l := measure252(pc, canonical)
	l.log(t, "  ")

	if l.final != (l.leg1 && l.leg2) {
		t.Errorf("leg breakdown does not explain the boolean: InAllowlist=%v but leg1=%v leg2=%v "+
			"(the workspace leg at paths.go:159 is inert in this probe - if it is not, the carrier changed under us)",
			l.final, l.leg1, l.leg2)
	}
	switch c.wantFinal {
	case -1:
	case 0:
		if l.final {
			t.Errorf("expected refusal for %s and got authorization", c.name)
		}
	case 1:
		if !l.final {
			t.Errorf("POSITIVE CONTROL RED for %s: same-shape ask was refused, so this probe's carrier is broken and no negative reading from it can be trusted", c.name)
		}
	}
}

// repoRoot252 derives the repository root from the test's own working directory
// (go test runs with cwd = the package directory), never from a typed literal.
func repoRoot252(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	root := filepath.Dir(filepath.Dir(wd))
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("derived repo root %q is not there: %v", root, err)
	}
	t.Logf("cwd=%q derived repo root=%q", wd, root)
	return root
}

// TestTicket252P1AllowlistLegs is AC#1 proper: the two directions plus the
// same-shape positive controls, on one real tree of this machine.
func TestTicket252P1AllowlistLegs(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	leaf := "q252p1-never-created.txt"
	t.Logf("long spelling : %q", long)
	t.Logf("short spelling: %q", short)
	t.Logf("equal?        : %v", strings.EqualFold(long, short))

	cases := []one252{
		// AC#1 Q1: roots long, asked short, leaf missing.
		{name: "Q1 roots=LONG asked=SHORT+missing-leaf", rootsCfg: []string{long}, asked: short + `\` + leaf, wantFinal: -1},
		// AC#1 Q2: roots short, asked long, leaf missing (ticket 107's shape).
		{name: "Q2 roots=SHORT asked=LONG+missing-leaf", rootsCfg: []string{short}, asked: long + `\` + leaf, wantFinal: -1},
		// AC#1 Q3 positive controls: one shape on both sides.
		{name: "P1 roots=LONG asked=LONG+missing-leaf", rootsCfg: []string{long}, asked: long + `\` + leaf, wantFinal: 1},
		{name: "P2 roots=LONG asked=SHORT of an EXISTING dir", rootsCfg: []string{long}, asked: short, wantFinal: 1},
		{name: "P3 roots=SHORT asked=SHORT of an EXISTING dir", rootsCfg: []string{short}, asked: short, wantFinal: 1},
		// Both sides configured short, leaf missing: the spelling is the same on
		// both sides of the CONFIG, which is not the same as being the same on
		// both sides of the comparison. Logged, not judged.
		{name: "Q3b roots=SHORT asked=SHORT+missing-leaf", rootsCfg: []string{short}, asked: short + `\` + leaf, wantFinal: -1},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) { run252(t, c) })
	}
}

// TestTicket252P1OwnerMachineReachability is AC#6's cell: the same question asked
// against paths that exist on this development machine, under real subdirectories
// of the repository, with a file name that does not exist yet.
func TestTicket252P1OwnerMachineReachability(t *testing.T) {
	long := repoRoot252(t)
	leaf := "q252p1-new-note.md"
	// A real, existing workspace directory BELOW the root: this is the shape an
	// fs.write against the owner's own project tree actually carries.
	for _, sub := range []string{"docs", filepath.Join("internal", "tools")} {
		dir := filepath.Join(long, sub)
		if _, err := os.Stat(dir); err != nil {
			t.Fatalf("carrier broken: expected workspace dir %q is not there: %v", dir, err)
		}
		shortDir := shortName252(t, dir)
		t.Logf("workspace dir=%q short=%q", dir, shortDir)
		run252(t, one252{name: "W-long " + sub, rootsCfg: []string{long}, asked: dir + `\` + leaf, wantFinal: 1})
		run252(t, one252{name: "W-short " + sub, rootsCfg: []string{long}, asked: shortDir + `\` + leaf, wantFinal: -1})
	}
}

// TestTicket252P1EnvSpellingCensus answers "who could hand production a short
// spelling on THIS machine": every environment spelling the C26 pipeline's step 1
// (risk.expandAccounted, via %VAR%/$VAR/~) can substitute, plus the process
// working directory that lexCanonical anchors relative paths onto.
func TestTicket252P1EnvSpellingCensus(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	names := []string{"USERPROFILE", "APPDATA", "LOCALAPPDATA", "TEMP", "TMP", "CD", "HOME"}
	for _, n := range names {
		v := os.Getenv(n)
		if v == "" {
			t.Logf("env %-13s: (unset)", n)
			continue
		}
		alias, verdict := aliasOf252(t, v)
		_, serr := os.Lstat(v)
		t.Logf("env %-13s: value=%q alias=%q -> %s lstat_err=%v", n, v, alias, verdict, serr)
	}
	alias, verdict := aliasOf252(t, wd)
	t.Logf("process cwd     : value=%q alias=%q -> %s", wd, alias, verdict)
	alias, verdict = aliasOf252(t, os.TempDir())
	t.Logf("os.TempDir()    : value=%q alias=%q -> %s", os.TempDir(), alias, verdict)
	home, herr := os.UserHomeDir()
	alias, verdict = "", "n/a"
	if herr == nil && home != "" {
		alias, verdict = aliasOf252(t, home)
	}
	t.Logf("os.UserHomeDir(): value=%q err=%v alias=%q -> %s", home, herr, alias, verdict)
	// The expansion leg itself, on the two spellings of the repo root: does step 1
	// of C26 touch a plain long or short path at all (Rewritten must be false for
	// both, or the roots would be dropped as unusable, not merely re-folded).
	root := repoRoot252(t)
	for _, s := range []string{root, shortName252(t, root)} {
		res, rerr := risk.Resolve(s, nil)
		t.Logf("Resolve(%q) -> Canonical=%q Resolved=%v Rewritten=%v err=%v", s, res.Canonical, res.Resolved, res.Rewritten, rerr)
	}
}
