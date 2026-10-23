//go:build !windows

package observe

// The POSIX half of ticket 132's instrument. There is no descriptor to read here:
// winsec's seal is a chmod (winsec_other.go, sealFile 0o600 / sealDir 0o700), so
// the reading is the mode bit that survives on the object itself. What the case
// can prove is therefore narrower than on Windows and says so: "no group or other
// access on the object the pipeline made". On this platform a file's mode is set
// by the same call that creates it (os.OpenFile with 0o600) and the seal repeats
// that mode, so the mode alone cannot tell the seal from the create - the
// per-object discriminator for the file half lives on Windows, and on both
// platforms the enumeration gate (sealcalls_132_test.go) is what turns red when
// the seal call is removed rather than the mode skipped.

import (
	"fmt"
	"os"
	"testing"
)

// judge132 reads the mode and reports whether the object carries no group or
// other access at all.
func judge132(t *testing.T, path string) (bool, string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	mode := info.Mode().Perm()
	if mode&0o077 != 0 {
		return false, fmt.Sprintf("%s: mode %04o still grants group/other access", path, mode)
	}
	want := os.FileMode(0o600)
	if info.IsDir() {
		want = 0o700
	}
	if mode != want {
		return false, fmt.Sprintf("%s: mode %04o, want %04o", path, mode, want)
	}
	return true, fmt.Sprintf("%s: mode %04o", path, mode)
}

func assertSealedDir132(t *testing.T, path string) { assertSealed132(t, path) }

func assertSealedFile132(t *testing.T, path string) { assertSealed132(t, path) }

func assertSealed132(t *testing.T, path string) {
	t.Helper()
	if ok, detail := judge132(t, path); !ok {
		t.Errorf("ticket 132: not sealed: %s", detail)
	} else {
		t.Logf("ticket 132: sealed: %s", detail)
	}
}

// widenDir132 opens the fixture up to the mode os.MkdirAll(0o755) used to leave
// behind, which is what an older build's log tree carries on this platform.
func widenDir132(t *testing.T, dir string) {
	t.Helper()
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if dirIsSealed132(t, dir) {
		t.Fatalf("fixture is wrong: %s still reads as sealed after chmod 0755, so the repair case would prove nothing", dir)
	}
}

func dirIsSealed132(t *testing.T, dir string) bool {
	t.Helper()
	ok, _ := judge132(t, dir)
	return ok
}
