package main

// Ticket 128 AC#2: the judgment cases for the ruling "refuse to start".
//
// WHAT WAS MEASURED FIRST (AC#1, docs/evidence/s1/128-ac1-consequences.md). With
// %APPDATA% unset this package's resolveDataDir answered base = "." and three legs
// disagreed about what that means: `wisp run` silently moved its whole state into
// the start-up directory (two same-named wisp-<date>-001.jsonl files in two
// trees, config.toml, a DPAPI secrets\ dir, and a second start in another
// directory reading an empty config), the resident leg refused with rc=1, and
// `wisp secret list` refused with rc=2.
//
// WHAT THESE CASES ANSWER, in AC#2's own words: same shape, every leg refuses, the
// reason is readable by a human, and - the part AC#1 could only observe as a
// consequence - the start-up directory gains not one byte.
//
// HOW THE SHAPE IS REACHED IN-PROCESS. os.UserConfigDir() fails only when the OS
// has no answer (no %APPDATA% on Windows, neither $XDG_CONFIG_HOME nor $HOME on
// POSIX); a CI runner cannot produce that for one case without moving every other
// case's temp dir with it. So the package reads the OS answer through the one
// variable it owns, userConfigDir, which production binds to os.UserConfigDir and
// never rebinds. The real, unbuilt-here shape is measured a second time in
// dataroot_128_windows_test.go against a compiled wisp.exe and a genuinely absent
// APPDATA - which is also where the resident leg can be driven at all, since it
// ends in os.Exit.
//
// THE SECOND AMBIENT PREMISE (ticket 128 AC#4, measured on both sides). Rebinding
// the seam only puts a leg on the refusal path if the leg ever asks the OS at all,
// and that is decided by WISP_ENV: resolveDataDir answers env "test" with
// proc.TestDataDir() before it reads userConfigDir. The test-windows job exports
// WISP_ENV=test for the whole run, so these legs ran green on a dev host and red
// on the runner until pinEnvThatAsksTheOS128 existed. Readings:
// docs/evidence/s1/128-ac4-gates-and-ci-divergence.md.
//
// AC#3 MUTATION ANCHOR: restore `base = "."` in resolveDataDir (cmd/wisp/doctor.go,
// the branch under `base, err := userConfigDir()`) and these named cases go red:
// TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128 (no error returned) and
// TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (every leg). Shortening the
// marker floor is a separate mutation and reddens only ShortenableList; every reading
// is in docs/evidence/s1/128-ac2-refusal-and-ac3-mutation.md.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
)

// planted128Cause is the OS's own words on a Windows box with no APPDATA, spelled
// the way syscall.Getenv spells it (AC#1 read it live off a real process).
var planted128Cause = errors.New("%AppData% is not defined")

// failConfigDir128 rebinds this package's single OS read of the config root for
// one case and puts it back through Cleanup, so no case can leave the refusal
// armed for the next one.
func failConfigDir128(t *testing.T) {
	t.Helper()
	prev := userConfigDir
	userConfigDir = func() (string, error) { return "", planted128Cause }
	t.Cleanup(func() { userConfigDir = prev })
}

// okConfigDir128 asserts the premise the refusal cases depend on: the seam is
// bound to the real os.UserConfigDir unless the case just rebinded it, and
// portable.txt is not sitting next to the test binary (that branch is checked
// first and would resolve the root without ever asking the OS).
func assertNoPortableOverride128(t *testing.T) {
	t.Helper()
	if exeDir := executableDir(); exeDir != "" {
		if _, err := os.Stat(filepath.Join(exeDir, "portable.txt")); err == nil {
			t.Fatalf("premise broke: a portable.txt sits next to the test binary (%s), so resolveDataDir takes the portable branch and never asks the OS for a config dir - these cases would assert nothing", exeDir)
		}
	}
}

// rescueMarkers128 are the three things AC#2's booked cost (A105 item 7: the run
// leg goes from "degraded but usable" to "not usable at all") demands of the
// refusing line: which variable is missing, where the data root belongs, and what
// to set to get it back. A list this short is exactly the shape R-121-2 is about,
// so TestAC2RefusalMarkersAreNotAShortenableList128 pins its length and its
// strict-prefix behaviour rather than trusting it.
var rescueMarkers128 = []string{
	"用户配置目录不可得",       // the verdict, and the OS read that failed
	"当前工作目录",          // the statement that "." is no longer a fallback
	"APPDATA",         // the Windows variable to set
	"XDG_CONFIG_HOME", // the POSIX variable to set
	"可写目录",            // what a correct value looks like
	"数据根本应是",          // where the data root belongs
}

// refusalLeg128 is one production consumer of the data root, driven in-process.
type refusalLeg128 struct {
	// entry is the name of the function that resolves the data root. The table
	// below is checked against this package's own call graph (see
	// TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128), which is what keeps
	// shortening it from being silent.
	entry string
	// drive runs the leg and reports its process exit code plus everything it
	// printed. A leg that cannot produce a code (the two helpers) maps
	// "refused" onto 2, the class the command legs return.
	drive func(t *testing.T) (int, string)
}

func refusalLegs128() []refusalLeg128 {
	return []refusalLeg128{
		{entry: "runTextTask", drive: func(t *testing.T) (int, string) {
			var out bytes.Buffer
			rc := runTextTask(runSpec{argv: []string{"票 128 探针"}, stdout: io.Discard, stderr: &out})
			return rc, out.String()
		}},
		{entry: "cmdModels", drive: func(t *testing.T) (int, string) {
			var out bytes.Buffer
			rc := cmdModels([]string{"list"}, modelsIO{stdout: io.Discard, stderr: &out})
			return rc, out.String()
		}},
		{entry: "cmdProviders", drive: func(t *testing.T) (int, string) {
			var out bytes.Buffer
			rc := cmdProviders([]string{"list"}, providersIO{stdout: io.Discard, stderr: &out})
			return rc, out.String()
		}},
		{entry: "cmdDoctor", drive: func(t *testing.T) (int, string) {
			var ok bool
			text := captureStdout128(t, func() { ok = cmdDoctor() })
			if ok {
				return 0, text
			}
			return 1, text
		}},
		{entry: "resolveSecretLayout", drive: func(t *testing.T) (int, string) {
			l, err := resolveSecretLayout(buildinfo.EnvDev)
			if err != nil {
				return 2, err.Error()
			}
			return 0, l.DataDir
		}},
	}
}

// TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128 is the verdict itself,
// per env fork. Every assertion here is the shape AC#1 measured as a disk write:
// dev forks to wisp-dev, prod to wisp, and "." appears nowhere.
func TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128(t *testing.T) {
	assertNoPortableOverride128(t)
	for _, env := range []string{"dev", "prod"} {
		t.Run(env, func(t *testing.T) {
			failConfigDir128(t)
			cwd := t.TempDir()
			t.Chdir(cwd)
			dir, err := resolveDataDir(env)
			if err == nil {
				t.Fatalf("AC#2 RED: resolveDataDir(%q) returned %q with no error - that is the pre-ticket-128 behaviour of falling back to the start-up directory (AC#1 measured it as two same-named log files, a moved config.toml, a second DPAPI store and an empty config on the next start)", env, dir)
			}
			if dir != "" {
				t.Errorf("AC#2 RED: a refusing resolution still handed back %q, want the empty string", dir)
			}
			if !errors.Is(err, errDataDirUnresolved) {
				t.Errorf("AC#2 RED: refusal is not of the class errDataDirUnresolved: %v", err)
			}
			for _, marker := range rescueMarkers128 {
				if !strings.Contains(err.Error(), marker) {
					t.Errorf("AC#2 RED: the refusal does not say %q, so nobody can act on it: %v", marker, err)
				}
			}
			want := filepath.Join("<用户配置目录>", "wisp")
			if env == "dev" {
				want = filepath.Join("<用户配置目录>", "wisp-dev")
			}
			if !strings.Contains(err.Error(), want) {
				t.Errorf("AC#2 RED: the refusal does not name the tree the data root belongs in (%s); it reads: %v", want, err)
			}
			if !strings.Contains(err.Error(), planted128Cause.Error()) {
				t.Errorf("AC#2 RED: the refusal drops the OS's own reason (%s), which is the only part that says which variable is empty: %v", planted128Cause, err)
			}
			assertDirEmpty128(t, cwd, "resolveDataDir("+env+")")
		})
	}
}

// TestAC2TestDataDirBranchStillResolves128 keeps the refusal from swallowing the
// harness branch: env "test" returns proc.TestDataDir() and asks the OS nothing.
// AC#1 could not reach base = "." under WISP_ENV=test for exactly this reason
// (A105 item 5a), and the fix must not change that either.
func TestAC2TestDataDirBranchStillResolves128(t *testing.T) {
	assertNoPortableOverride128(t)
	failConfigDir128(t)
	dir, err := resolveDataDir("test")
	if err != nil {
		t.Fatalf("AC#2: the test-env resolution refused as well: %v - the harness root is declared, not asked for", err)
	}
	if dir != proc.TestDataDir() {
		t.Errorf("AC#2: test data dir = %q, want %q", dir, proc.TestDataDir())
	}
}

// TestAC2RefusalMarkersAreNotAShortenableList128 is the R-121-2 nail for this
// ticket's own instrument: a marker list can only be trusted while dropping an
// entry from the front is caught. The floor says how many entries the list has to
// have, and the negative leg proves the check is a prefix-sensitive one rather
// than "is any of these present".
func TestAC2RefusalMarkersAreNotAShortenableList128(t *testing.T) {
	if len(rescueMarkers128) < 6 {
		t.Fatalf("AC#2 RED (the instrument, not the code): rescueMarkers128 holds %d entries; AC#2's booked cost requires at least 6 (verdict, no-fallback statement, both variable spellings, the writable-directory rule, the tree the root belongs in). Shortening this list is how a refusal stops being actionable and the gate keeps reading green", len(rescueMarkers128))
	}
	full := strings.Join(rescueMarkers128, "|")
	// Negative leg: every strict prefix of the marker list must fail the very
	// matcher the legs use, so a check that quietly became "at least one marker
	// present" cannot pass here.
	for cut := 1; cut < len(rescueMarkers128); cut++ {
		partial := strings.Join(rescueMarkers128[:cut], "|")
		if allMarkersPresent128(partial, rescueMarkers128) {
			t.Errorf("AC#2 RED (the instrument): a text holding only the first %d of %d markers satisfied allMarkersPresent128, so the marker check is not the conjunctive one AC#2 needs", cut, len(rescueMarkers128))
		}
	}
	if !allMarkersPresent128(full, rescueMarkers128) {
		t.Fatal("AC#2 RED (the instrument): the complete marker text failed its own matcher, so nothing below can be believed")
	}
}

// allMarkersPresent128 is the matcher the legs use, factored out so the case above
// can attack it.
func allMarkersPresent128(text string, markers []string) bool {
	for _, m := range markers {
		if !strings.Contains(text, m) {
			return false
		}
	}
	return true
}

// pinEnvThatAsksTheOS128 pins the one ambient variable that decides whether any
// leg below ever asks the OS for its config dir, and proves the pin took.
//
// WHY, measured on both sides rather than reasoned out: .github/workflows/ci.yml
// gives the whole test-windows job `WISP_ENV: test`, and resolveDataDir answers
// env "test" with proc.TestDataDir() *before* it reads userConfigDir, so on the
// runner the seam this case rebinds was never consulted. Four legs went red over
// their missing markers there while the fifth, resolveSecretLayout, stayed green -
// it passes buildinfo.EnvDev explicitly and so never looks at the environment.
// That is the whole of the "green on the laptop, red on CI" report this case
// earned: run 35967768017 versus the same command on a dev host where WISP_ENV is
// unset and buildinfo.DefaultEnv is "dev".
//
// The case PLANTS the runner's value first and only then pins over it, so the pin
// has a tooth on every platform: delete the pin and this case goes red on a laptop
// too, instead of staying green here and going red only where nobody is looking.
// Nothing below changes what a leg has to print (that is AC#2's ruling), and no
// failure here is turned into a skip.
func pinEnvThatAsksTheOS128(t *testing.T) {
	t.Helper()
	ambient, hadAmbient := os.LookupEnv("WISP_ENV")
	t.Setenv("WISP_ENV", "test")                   // the value the windows job exports
	t.Setenv("WISP_ENV", string(buildinfo.EnvDev)) // the pin that has to beat it
	failConfigDir128(t)
	if got := buildinfo.EnvString(); got != string(buildinfo.EnvDev) {
		t.Fatalf("premise broke: the pin did not hold - WISP_ENV resolves to %q (ambient on entry was %q, present=%t), so no leg below ever asks the OS and its markers prove nothing", got, ambient, hadAmbient)
	}
	dir, err := resolveDataDir(buildinfo.EnvString())
	if !errors.Is(err, errDataDirUnresolved) {
		t.Fatalf("premise broke: with WISP_ENV=%q and a failed OS read, resolveDataDir returned (%q, %v) instead of errDataDirUnresolved - the resolution stopped asking userConfigDir, and the five legs below would then be watching nothing", buildinfo.EnvString(), dir, err)
	}
	if dir != "" {
		t.Fatalf("premise broke: a refusing resolution handed back %q", dir)
	}
}

// TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 is the three-legs-in-one
// picture AC#2 was ruled on. The leg table is validated against this package's own
// call graph before a single leg runs: every production function that resolves the
// data root has to be driven here, and a row that names a function which no longer
// does is stale. That is what makes deleting a row, or moving a resolution
// somewhere the gate cannot see, loud instead of silent.
func TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128(t *testing.T) {
	assertNoPortableOverride128(t)

	want := map[string]bool{}
	for _, f := range resolveDataDirConsumers128(t) {
		want[f] = true
	}
	// resolveSecretLayout is not a caller of resolveDataDir - it is the second,
	// independent reader of the same OS answer, and AC#1's third leg. It is named
	// here rather than left out, so the set below is checked for equality rather
	// than containment.
	want["resolveSecretLayout"] = true

	got := map[string]bool{}
	for _, leg := range refusalLegs128() {
		if got[leg.entry] {
			t.Errorf("AC#2 RED (the instrument): leg %q is registered twice, which means the table no longer mirrors anything", leg.entry)
		}
		got[leg.entry] = true
	}
	var missing, stale []string
	for k := range want {
		if !got[k] {
			missing = append(missing, k)
		}
	}
	for k := range got {
		if !want[k] {
			stale = append(stale, k)
		}
	}
	sort.Strings(missing)
	sort.Strings(stale)
	if len(missing) > 0 {
		t.Fatalf("AC#2 RED: %d function(s) resolve the data root in this package and no leg here drives them (%s). Either a new consumer landed without a refusal case, or the call graph moved the resolution out of sight - both leave a leg free to fall back to the start-up directory again",
			len(missing), strings.Join(missing, ", "))
	}
	if len(stale) > 0 {
		t.Errorf("AC#2 RED (the instrument): %d registered leg(s) no longer resolve the data root (%s); the table is now asserting a shape the code does not have",
			len(stale), strings.Join(stale, ", "))
	}

	pinEnvThatAsksTheOS128(t)
	for _, leg := range refusalLegs128() {
		t.Run(leg.entry, func(t *testing.T) {
			cwd := t.TempDir()
			t.Chdir(cwd)
			rc, out := leg.drive(t)
			if rc == 0 {
				t.Fatalf("AC#2 RED: leg %q exited 0 with no data root - it either resolved one out of the start-up directory or did nothing at all. Output:\n%s", leg.entry, out)
			}
			if !allMarkersPresent128(out, rescueMarkers128) {
				var absent []string
				for _, m := range rescueMarkers128 {
					if !strings.Contains(out, m) {
						absent = append(absent, m)
					}
				}
				t.Errorf("AC#2 RED: leg %q refuses, but the line a human sees does not say %s (missing markers). Output:\n%s", leg.entry, strings.Join(absent, " / "), out)
			}
			// The consequence AC#1 measured: the start-up directory growing a
			// wisp-dev\ with logs\, config.toml and secrets\ inside it.
			assertDirEmpty128(t, cwd, "leg "+leg.entry)
		})
	}
}

// resolveDataDirConsumers128 reads this package's source for every production
// function that calls resolveDataDir. No leg list is written in this file: the
// list is the call graph, which is the only thing that grows when somebody adds a
// consumer.
func resolveDataDirConsumers128(t *testing.T) []string {
	t.Helper()
	pkg, err := loadMainPackage131(".")
	if err != nil {
		t.Fatalf("AC#2 RED (the instrument): %v", err)
	}
	// Non-vacuity, R-121-2's other half: a loader that parsed nothing, or parsed
	// the directory under a name that hides the files, would return an empty
	// caller set and make the equality check above trivially satisfiable by an
	// empty table.
	if len(pkg.funcs) < 40 {
		t.Fatalf("AC#2 RED (the instrument): the package parse found only %d production functions, so its caller set cannot be trusted as a denominator", len(pkg.funcs))
	}
	var out []string
	for name, fn := range pkg.funcs {
		if name == "resolveDataDir" {
			continue // the resolver itself is not one of its consumers
		}
		if fn.calls["resolveDataDir"] {
			out = append(out, name)
		}
	}
	if len(out) < 4 {
		t.Fatalf("AC#2 RED (the instrument): only %d functions in this package resolve the data root (want at least 4: the run, models and providers commands plus doctor's own check). Someone replaced resolveDataDir with an inlined guess and the gate would otherwise watch nothing", len(out))
	}
	sort.Strings(out)
	return out
}

// assertDirEmpty128 fails when dir holds anything at all, named or nested.
func assertDirEmpty128(t *testing.T, dir, who string) {
	t.Helper()
	var found []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p != dir {
			rel, rerr := filepath.Rel(dir, p)
			if rerr != nil {
				rel = p
			}
			found = append(found, fmt.Sprintf("%s (%v)", rel, d.IsDir()))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("AC#2: walking %s after %s: %v", dir, who, err)
	}
	if len(found) > 0 {
		t.Errorf("AC#2 RED: %s wrote into the start-up directory instead of refusing - AC#1's搬家 result in one line: %s",
			who, strings.Join(found, ", "))
	}
}

// captureStdout128 runs fn with os.Stdout pointing at a pipe and returns what was
// printed. cmdDoctor writes to stdout directly, and its self-rescuing line is the
// one AC#2's cost ledger promises the operator.
func captureStdout128(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	prev := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = prev })
	done := make(chan string, 1)
	go func() {
		var sb strings.Builder
		buf := make([]byte, 4096)
		for {
			n, rerr := r.Read(buf)
			sb.Write(buf[:n])
			if rerr != nil {
				break
			}
		}
		done <- sb.String()
	}()
	fn()
	w.Close()
	os.Stdout = prev
	out := <-done
	r.Close()
	return out
}
