package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// withInjectedSealFailure replaces the platform's descriptor application with a
// failure. AC#5 is about what the *callers* do when the seal cannot be applied,
// and the only honest way to ask that question is to make it fail.
func withInjectedSealFailure(t *testing.T, err error) {
	t.Helper()
	orig := applyDescriptor
	applyDescriptor = func(string, bool) error { return err }
	t.Cleanup(func() { applyDescriptor = orig })
}

// TestAC5FailedSealRefusesTheWrite is the "failure direction may only tighten"
// criterion: an injected seal failure must produce (a) an error naming
// ErrNotSealable, (b) no sensitive bytes left on disk, and (c) no half-made
// directory tree reported as private. "Could not set the ACL, wrote it wide
// anyway" is the exact behaviour this package exists to make impossible.
func TestAC5FailedSealRefusesTheWrite(t *testing.T) {
	injected := errors.New("injected: descriptor could not be applied")
	dir := SealableTempDirForTest124(t)

	t.Run("exclusive artifact", func(t *testing.T) {
		withInjectedSealFailure(t, injected)
		p := filepath.Join(dir, "artifact.txt")
		err := PrivateFileExclusive(p, []byte("tool output that must not land wide"))
		if err == nil {
			t.Fatal("the write succeeded although the seal was refused")
		}
		if !errors.Is(err, ErrNotSealable) {
			t.Fatalf("error does not name the refusal: %v", err)
		}
		assertNoBytesOnDisk(t, p)
	})

	t.Run("replacement write", func(t *testing.T) {
		withInjectedSealFailure(t, injected)
		p := filepath.Join(dir, "blob.bin")
		err := PrivateFile(p, []byte("dpapi-shaped secret"), 0o600)
		if err == nil {
			t.Fatal("PrivateFile accepted a refused seal")
		}
		if !errors.Is(err, ErrNotSealable) {
			t.Fatalf("error does not name the refusal: %v", err)
		}
		assertNoBytesOnDisk(t, p)
	})

	t.Run("directory chain", func(t *testing.T) {
		withInjectedSealFailure(t, injected)
		p := filepath.Join(dir, "a", "b", "c")
		if err := PrivateDirAll(p, 0o700); err == nil {
			t.Fatal("PrivateDirAll reported a private directory it could not seal")
		} else if !errors.Is(err, ErrNotSealable) {
			t.Fatalf("error does not name the refusal: %v", err)
		}
	})

	t.Run("seal file and dir direct", func(t *testing.T) {
		p := filepath.Join(dir, "plain.txt")
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		withInjectedSealFailure(t, injected)
		if err := SealFile(p); !errors.Is(err, ErrNotSealable) {
			t.Fatalf("SealFile: %v", err)
		}
		if err := SealDir(dir); !errors.Is(err, ErrNotSealable) {
			t.Fatalf("SealDir: %v", err)
		}
	})
}

// assertNoBytesOnDisk is the load-bearing half of the injection test: an error
// that leaves the artifact readable-by-everyone on the floor has not refused
// anything.
func assertNoBytesOnDisk(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err == nil {
		t.Errorf("refused write left %d bytes on disk at %s", info.Size(), filepath.Base(path))
		return
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("stat after refusal: %v", err)
	}
}

// TestAC5FailureIsNotSwallowedByTheHappyPath guards the other direction: with no
// injection the same calls must succeed, so the two tests above cannot both be
// satisfied by an implementation that always errors or one that always ignores.
func TestAC5FailureIsNotSwallowedByTheHappyPath(t *testing.T) {
	dir := SealableTempDirForTest124(t)
	p := filepath.Join(dir, "artifact.txt")
	if err := PrivateFileExclusive(p, []byte("written while the seal works")); err != nil {
		t.Fatalf("PrivateFileExclusive: %v", err)
	}
	if got, err := os.ReadFile(p); err != nil || string(got) != "written while the seal works" {
		t.Fatalf("read back %q, %v", got, err)
	}
	if err := PrivateFileExclusive(p, []byte("second")); !errors.Is(err, fs.ErrExist) {
		t.Errorf("exclusive create must keep failing on an occupied name, got %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Errorf("the failed exclusive create damaged the existing artifact: %v", err)
	}
}
