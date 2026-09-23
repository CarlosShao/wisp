//go:build windows

package main

// Ticket 128 AC#2, the real-process half. The in-process cases in
// dataroot_128_test.go inject the OS's "no config dir" answer; this file measures
// the same shape the way AC#1 did it - a compiled wisp.exe, APPDATA genuinely
// absent from the child's environment, and the working directory set to a fresh
// empty tree - because two of the legs cannot be driven any other way: the
// resident leg ends in os.Exit, and its refusal is produced three packages away
// from this one.
//
// What each row has to show, in AC#2's words: the leg refuses (non-zero exit),
// the line a human reads names the missing variable and the fix, and the
// directory it started in stays empty. The last row is the one that answers AC#1:
// before this ticket `wisp run` wrote wisp-dev\logs\wisp-<date>-001.jsonl,
// wisp-dev\config.toml and a sealed wisp-dev\secrets\ into the start-up directory
// while the other two legs refused, so the same machine ended up with two
// same-named logs in two trees and a second start that read an empty config.
//
// The exe is built once for the whole table by ticket 118's helper
// (buildWispForTest in secret_argv_windows_test.go), which also colocates the
// sherpa/onnxruntime DLLs - without that the child dies at load time with
// 0xC0000135 and zero output, which reads exactly like "the code path never ran".

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// realProcessLeg128 is one command line of the shipped binary.
type realProcessLeg128 struct {
	name string
	argv []string
	// markers are the pieces the refusing line has to contain. Every leg except
	// the resident one now speaks with this package's single voice
	// (dataDirUnresolved128); the resident leg's sentence is produced by
	// internal/proc's DefaultLayout and is checked for that instead.
	markers []string
}

func TestAC2RealProcessRefusesOnEveryLegWithoutAppData128(t *testing.T) {
	exe := buildWispForTest(t)
	env := appDataFreeEnv128(t)
	// Premise, and the one thing that would silently turn every row green: a
	// portable.txt next to the exe resolves the data root without ever asking
	// the OS for a config dir.
	if _, err := os.Stat(filepath.Join(filepath.Dir(exe), "portable.txt")); err == nil {
		t.Fatal("premise broke: portable.txt next to the built exe, so the portable branch wins and no leg below ever reads APPDATA")
	}

	for _, leg := range []realProcessLeg128{
		{name: "run", argv: []string{"run", "票 128 探针"}, markers: rescueMarkers128},
		{name: "secret-list", argv: []string{"secret", "list"}, markers: rescueMarkers128},
		{name: "doctor", argv: []string{"doctor"}, markers: rescueMarkers128},
		{name: "resident", argv: nil, markers: []string{"user config dir", "%AppData% is not defined"}},
	} {
		t.Run(leg.name, func(t *testing.T) {
			cwd := t.TempDir()
			rc, out := runWispNoAppData128(t, exe, env, cwd, leg.argv)
			// The reading the交回件 quotes, verbatim, before anything is asserted.
			t.Logf("AC#2 real process, leg %q: APPDATA unset, WISP_ENV=dev, cwd=%s -> rc=%d\n%s",
				leg.name, cwd, rc, out)
			if rc == 0 {
				t.Fatalf("AC#2 RED: leg %q exited 0 with APPDATA absent from its environment, so it resolved a data root somewhere instead of refusing:\n%s", leg.name, out)
			}
			for _, marker := range leg.markers {
				if !strings.Contains(out, marker) {
					t.Errorf("AC#2 RED: leg %q refused but its line never says %q:\n%s", leg.name, marker, out)
				}
			}
			assertDirEmpty128(t, cwd, "real leg "+leg.name)
		})
	}
}

// appDataFreeEnv128 builds the child environment with APPDATA removed (compared
// case-insensitively: this is Windows, where the variable arrives as APPDATA,
// AppData or appdata depending on who set it) and WISP_ENV=dev pinned.
//
// WISP_ENV=test is not an option here and that is a measurement, not a preference:
// resolveDataDir returns proc.TestDataDir() before it ever asks the OS, so the test
// env cannot reach the branch this ticket is about (A105 item 5a).
func appDataFreeEnv128(t *testing.T) []string {
	t.Helper()
	kept := make([]string, 0, len(os.Environ()))
	var removed int
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.EqualFold(name, "APPDATA") {
			removed++
			continue
		}
		kept = append(kept, kv)
	}
	if removed == 0 {
		t.Fatal("premise broke: this process has no APPDATA to remove, so an appdata-free child proves nothing - the shape has to be created, not assumed")
	}
	// WISP_TEST_DATA_DIR would fork the resolution somewhere else entirely.
	for i, kv := range kept {
		if strings.EqualFold(strings.SplitN(kv, "=", 2)[0], "WISP_TEST_DATA_DIR") {
			kept = append(kept[:i], kept[i+1:]...)
			break
		}
	}
	return append(kept, "WISP_ENV=dev")
}

// runWispNoAppData128 starts the built binary in cwd with argv/env and returns its
// exit code plus combined output.
//
// The 120 s context is a guard, not a measurement: every leg here is expected to
// exit on its own in well under a second, and the only shape that can hang is a
// regression that lets the resident leg boot an event loop nobody asked for. No
// timing claim is taken from this function anywhere (A103: no wall-clock or RSS
// readings were needed for this ticket).
func runWispNoAppData128(t *testing.T, exe string, env []string, cwd string, argv []string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, argv...)
	cmd.Dir = cwd
	cmd.Env = env
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	rc := 0
	if err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("starting %s %v in %s: %v", exe, argv, cwd, err)
		}
		rc = ee.ExitCode()
	}
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("AC#2 RED: leg never exited within the 120 s guard budget, so it did not refuse - it booted something (resident event loop?) and had to be killed by the context. Partial output:\n%s", out.String())
	}
	return rc, strings.ToValidUTF8(out.String(), "?")
}
