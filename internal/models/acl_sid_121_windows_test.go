//go:build windows

package models

// The live half of ticket 121 AC#4 / R-109-3: turning what icacls printed on a
// real object into an identity, so that the two ACL measurements in this package
// (handoff_window_109_windows_test.go's "who holds the write in the hand-off
// window" and no_seal_ruling_windows_test.go's "this class is inherit-wide on
// purpose") judge grants the way internal/winsec does, by SID.
//
// The pure parser and its synthetic-dump cases are in acl_sid_121_test.go, which
// is untagged so both CI legs can read them. This file only supplies the
// directory service the parser needs for a real descriptor, which by definition
// exists only where the descriptor does.

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// sidTranslateCache keeps name -> SID resolutions for the test binary's
// lifetime. Each resolution is a subprocess, and a handful of distinct trustees
// covers every object these two cases probe.
var sidTranslateCache = struct {
	sync.Mutex
	m map[string]string
}{m: map[string]string{}}

// liveSIDResolver answers "which principal is this trustee text". A bare numeric
// SID needs no answer and is handled by aceLine.sidIdentity before this is ever
// consulted; anything else must translate or the case fails loudly, because a
// skipped ACE is the false green this whole file exists to remove.
func liveSIDResolver(t *testing.T) func(string) (string, error) {
	t.Helper()
	return func(account string) (string, error) {
		sidTranslateCache.Lock()
		defer sidTranslateCache.Unlock()
		if sid, ok := sidTranslateCache.m[account]; ok {
			return sid, nil
		}
		out, err := runPowerShell(fmt.Sprintf(
			"(New-Object System.Security.Principal.NTAccount '%s')."+
				"Translate([System.Security.Principal.SecurityIdentifier]).Value",
			strings.ReplaceAll(account, "'", "''")))
		if err != nil {
			return "", fmt.Errorf("cannot resolve trustee %q to a SID: %w (%s)", account, err, strings.TrimSpace(out))
		}
		sid := strings.TrimSpace(out)
		if !isNumericSID(sid) {
			return "", fmt.Errorf("resolution of %q produced %q, not a SID", account, sid)
		}
		sidTranslateCache.m[account] = sid
		return sid, nil
	}
}

// runPowerShell is the translation backend, the same one internal/winsec's
// acl_windows_test.go uses for the same purpose.
func runPowerShell(script string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	err := cmd.Run()
	return buf.String(), err
}

// aceTextForSID returns the ACE lines of path's descriptor that belong to the
// principal named by account, rendered as "trustee rights" so a failing log line
// still shows what icacls actually said. account is only ever a label used to
// reach a SID; the comparison itself never looks at it.
func aceTextForSID(t *testing.T, path, account string) []string {
	t.Helper()
	want, err := liveSIDResolver(t)(account)
	if err != nil {
		t.Fatalf("%v", err)
	}
	owned, err := aceLinesForSID(icaclsRun(t, path), path, want, liveSIDResolver(t))
	if err != nil {
		t.Fatalf("reading the descriptor of %s: %v", path, err)
	}
	out := make([]string, 0, len(owned))
	for _, a := range owned {
		out = append(out, a.Trustee+" "+a.Rights)
	}
	return out
}

// allTrusteesOn returns every trustee on the object, whatever its identity
// resolves to. Used to state, in a log line a human can re-check against icacls,
// that the set the SID filter chose from was not empty.
func allTrusteesOn(t *testing.T, path string) []string {
	t.Helper()
	lines := parseACELines(icaclsRun(t, path), path)
	out := make([]string, 0, len(lines))
	for _, a := range lines {
		out = append(out, a.Trustee+" "+a.Rights)
	}
	return out
}
