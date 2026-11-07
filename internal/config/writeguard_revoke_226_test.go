package config

// Ticket 226 AC#5 - the 落库 half of revoking a long-lived rule, measured on the
// reading that decides whether the rule is still in force.
//
// TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber (writeguard_226_test.go)
// already pins the no-clobber properties of the write. What this file adds is
// the effect reading AC#5's tooth asks for: after the revocation, the list the
// NEXT START loads no longer grants the revoked directory - asserted through
// LoadFile and through a fresh NewManager, the two shapes a boot actually goes
// through, not through a field on the writer's own struct. The before-side
// control loads the same start content fresh and finds the entry still there,
// so the absence after the write is attributable to the revocation and not to a
// fixture that never granted anything.

import (
	"reflect"
	"strings"
	"testing"
)

func TestAC5RevokedRuleIsGoneForTheNextStart(t *testing.T) {
	start := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data", "D:\\notes"]`, 1)

	// Control: the start content this case revokes from really grants both
	// directories when a fresh start loads it.
	m0, _, _ := newTestManager(t, start)
	if !reflect.DeepEqual(m0.Config().FS.AllowedDirs, []string{`D:\data`, `D:\notes`}) {
		t.Fatalf("the fixture does not grant what it revokes: %v", m0.Config().FS.AllowedDirs)
	}

	m, path, _ := newTestManager(t, start)
	if err := m.SetAllowedDirs([]string{`D:\notes`}); err != nil {
		t.Fatalf("SetAllowedDirs (the revoke write): %v", err)
	}

	// LoadFile is what a boot parses; the revoked entry must be out of it.
	got := readLoaded(t, path)
	if reflect.DeepEqual(got.FS.AllowedDirs, []string{`D:\data`, `D:\notes`}) {
		t.Fatalf("the revoked entry is still what the file grants: %v", got.FS.AllowedDirs)
	}
	if !reflect.DeepEqual(got.FS.AllowedDirs, []string{`D:\notes`}) {
		t.Errorf("the revocation took more than the one named entry: %v", got.FS.AllowedDirs)
	}

	// A fresh NewManager on the same path is the next start itself; its view is
	// what the [fs] gate is built from (run.go:624), so the rule is really out
	// of force only when this view has lost the entry too.
	m2, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("next-start load: %v", err)
	}
	if reflect.DeepEqual(m2.Config().FS.AllowedDirs, []string{`D:\data`, `D:\notes`}) {
		t.Fatalf("the next start would still grant the revoked entry: %v", m2.Config().FS.AllowedDirs)
	}
	if !reflect.DeepEqual(m2.Config().FS.AllowedDirs, []string{`D:\notes`}) {
		t.Errorf("the next start lost entries the revocation did not name: %v", m2.Config().FS.AllowedDirs)
	}

	// The running writer's own memory agrees with the file: file and memory
	// cannot diverge after a successful revoke.
	if !reflect.DeepEqual(m.Config().FS.AllowedDirs, []string{`D:\notes`}) {
		t.Errorf("memory and file diverged after the revoke: %v", m.Config().FS.AllowedDirs)
	}
}
