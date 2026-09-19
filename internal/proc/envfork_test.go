package proc

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// TestLayoutForProdDev pins the frozen SPEC-03 §5.2 defaults for prod and dev.
func TestLayoutForProdDev(t *testing.T) {
	root := t.TempDir()

	prod, err := LayoutFor(buildinfo.EnvProd, root)
	if err != nil {
		t.Fatalf("LayoutFor(prod): %v", err)
	}
	if prod.DataDir != filepath.Join(root, "wisp") {
		t.Errorf("prod DataDir = %q, want %q", prod.DataDir, filepath.Join(root, "wisp"))
	}
	if prod.MutexName != `Local\wisp-single-instance` {
		t.Errorf("prod MutexName = %q", prod.MutexName)
	}
	if !prod.MutexEnabled || prod.ActivateEventName == "" {
		t.Errorf("prod must register a mutex with an activation event: %+v", prod)
	}

	dev, err := LayoutFor(buildinfo.EnvDev, root)
	if err != nil {
		t.Fatalf("LayoutFor(dev): %v", err)
	}
	if dev.DataDir != filepath.Join(root, "wisp-dev") {
		t.Errorf("dev DataDir = %q, want %q", dev.DataDir, filepath.Join(root, "wisp-dev"))
	}
	if dev.MutexName != `Local\wisp-dev-single-instance` {
		t.Errorf("dev MutexName = %q", dev.MutexName)
	}

	// dev and prod must be mutually invisible (SPEC-03 §5.2).
	if prod.DataDir == dev.DataDir {
		t.Error("prod and dev share a data dir")
	}
	if prod.MutexName == dev.MutexName {
		t.Error("prod and dev share a mutex")
	}
	if !strings.HasPrefix(prod.MutexName, `Local\`) {
		t.Error("mutex must be per-session (Local namespace, D42#7)")
	}
}

// TestLayoutForTestEnv pins what ticket 03 guarantees for the test env: no
// mutex, and the concrete values explicitly deferred to ticket 06.
func TestLayoutForTestEnv(t *testing.T) {
	l, err := LayoutFor(buildinfo.EnvTest, t.TempDir())
	if !errors.Is(err, ErrTestLayoutDeferred) {
		t.Fatalf("LayoutFor(test) err = %v, want ErrTestLayoutDeferred", err)
	}
	if l.MutexEnabled || l.MutexName != "" {
		t.Errorf("test env must not register a mutex: %+v", l)
	}
	if l.Env != buildinfo.EnvTest {
		t.Errorf("layout env = %q", l.Env)
	}
}

// TestLayoutForUnknownEnv rejects anything outside the enum.
func TestLayoutForUnknownEnv(t *testing.T) {
	if _, err := LayoutFor(buildinfo.Env("staging"), t.TempDir()); err == nil {
		t.Fatal("unknown env must fail")
	}
}
