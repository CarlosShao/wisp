package proc

import (
	"fmt"
	"os"
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
	if !prod.AutoUpdateChecks {
		t.Error("prod must default auto-update checks ON (SPEC-03 §5.2)")
	}
	if prod.DefaultLLMBaseURL != "" || prod.DefaultMirrorURL != "" {
		t.Errorf("prod endpoints are user-configured, got LLM=%q mirror=%q", prod.DefaultLLMBaseURL, prod.DefaultMirrorURL)
	}
	if prod.Portable {
		t.Error("fork layout must not claim portable before the override runs")
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
	if dev.AutoUpdateChecks {
		t.Error("dev must default auto-update checks OFF (SPEC-03 §5.2)")
	}
	if dev.DefaultLLMBaseURL != DevLLMBaseURL {
		t.Errorf("dev DefaultLLMBaseURL = %q, want %q", dev.DefaultLLMBaseURL, DevLLMBaseURL)
	}
	if dev.DefaultMirrorURL != DevMirrorBaseURL {
		t.Errorf("dev DefaultMirrorURL = %q, want %q", dev.DefaultMirrorURL, DevMirrorBaseURL)
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

// TestLayoutForTestEnv pins the ticket-06 test-env fork (SPEC-03 §5.2): data
// dir from WISP_TEST_DATA_DIR or %TEMP%\wisp-test-<pid>, no mutex, no endpoint
// defaults (the harness injects them explicitly).
func TestLayoutForTestEnv(t *testing.T) {
	t.Run("explicit WISP_TEST_DATA_DIR", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv(TestDataDirEnv, dir)
		l, err := LayoutFor(buildinfo.EnvTest, t.TempDir())
		if err != nil {
			t.Fatalf("LayoutFor(test): %v", err)
		}
		if l.DataDir != dir {
			t.Errorf("test DataDir = %q, want the WISP_TEST_DATA_DIR injection %q", l.DataDir, dir)
		}
		if l.MutexEnabled || l.MutexName != "" {
			t.Errorf("test env must not register a mutex: %+v", l)
		}
		if l.DefaultLLMBaseURL != "" || l.DefaultMirrorURL != "" {
			t.Errorf("test endpoints must be injected explicitly, got %+v", l)
		}
		if l.AutoUpdateChecks {
			t.Error("test env must not check for updates")
		}
		if l.Env != buildinfo.EnvTest {
			t.Errorf("layout env = %q", l.Env)
		}
	})

	t.Run("falls back to TEMP per pid", func(t *testing.T) {
		t.Setenv(TestDataDirEnv, "")
		l, err := LayoutFor(buildinfo.EnvTest, t.TempDir())
		if err != nil {
			t.Fatalf("LayoutFor(test): %v", err)
		}
		// Ticket 118 AC#9. This leg used to compare DataDir with
		// filepath.Join(os.TempDir(), "wisp-test-<pid>"), i.e. with whatever
		// spelling TMPDIR happens to carry. Ticket 119 made TestDataDir resolve that
		// root, so on the shape 119 exists for - TMPDIR behind a symlink, which is
		// macOS's /var and any Linux host that links its temp tree - the two sides
		// stopped being the same string and a correct implementation read as a
		// failure. The contract is not a spelling, so it is stated as the contract:
		// the per-pid leaf the caller is promised, on a root that reaches the temp
		// tree by identity and has no link left anywhere in it. The last part is
		// what the placement floor (internal/winsec/winsec_other.go, ticket 113)
		// refuses a spelling with a link in it *for*, and it is read from the
		// filesystem here rather than from SealableRoot: an expectation computed by
		// the function under test cannot fail.
		leaf := fmt.Sprintf("wisp-test-%d", os.Getpid())
		if got := filepath.Base(l.DataDir); got != leaf {
			t.Errorf("test DataDir = %q, want the per-pid leaf %q under the temp tree", l.DataDir, leaf)
		}
		parent := filepath.Dir(l.DataDir)
		real, err := filepath.EvalSymlinks(os.TempDir())
		if err != nil {
			t.Fatalf("EvalSymlinks(%q): %v", os.TempDir(), err)
		}
		gotInfo, gotErr := os.Stat(parent)
		wantInfo, wantErr := os.Stat(real)
		if gotErr != nil || wantErr != nil {
			t.Fatalf("stat the root (%v) against the temp tree %q (%v)", gotErr, real, wantErr)
		}
		if !os.SameFile(gotInfo, wantInfo) {
			t.Errorf("the test data root %q is not the tree this machine's temp dir reaches (%q): TestDataDir must resolve os.TempDir(), not substitute a directory of its own", parent, real)
		}
		if link := firstLinkInPath(parent); link != "" {
			t.Errorf("the test data root %q still reaches itself through the link at %q, so the seal that follows refuses it (ticket 113's placement floor); the layer that asks the OS must resolve what the OS answered (ticket 119)", parent, link)
		}
		t.Logf("test DataDir = %q (TMPDIR = %q, which reaches %q)", l.DataDir, os.TempDir(), real)
	})
}

// firstLinkInPath answers, from the filesystem and not from a string, which prefix
// of path reaches itself through a symlink ("" when none does). It walks the way
// SealableRoot walks - ask the OS about this spelling, then about its parent, until
// the parents stop differing - so the claim it supports is "no link is left in this
// root", which is the property internal/winsec's placement floor checks before it
// will chmod anything.
func firstLinkInPath(path string) string {
	for cur := path; ; {
		if info, err := os.Lstat(cur); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return ""
		}
		cur = parent
	}
}

// TestLayoutForkMatrix is the SPEC-03 §6 environment-fork acceptance: the
// three envs resolve to mutually distinct data dirs / mutex names / endpoint
// defaults, and dev and prod are mutually invisible.
func TestLayoutForkMatrix(t *testing.T) {
	root := t.TempDir()
	t.Setenv(TestDataDirEnv, filepath.Join(root, "injected-test-dir"))

	layouts := make(map[buildinfo.Env]Layout, 3)
	for _, env := range []buildinfo.Env{buildinfo.EnvProd, buildinfo.EnvDev, buildinfo.EnvTest} {
		l, err := LayoutFor(env, root)
		if err != nil {
			t.Fatalf("LayoutFor(%s): %v", env, err)
		}
		layouts[env] = l
	}

	dirs := map[buildinfo.Env]string{}
	for env, l := range layouts {
		dirs[env] = l.DataDir
		for other, otherDir := range dirs {
			if env != other && l.DataDir == otherDir {
				t.Errorf("%s and %s share data dir %q", env, other, l.DataDir)
			}
		}
	}
	if layouts[buildinfo.EnvProd].MutexName == layouts[buildinfo.EnvDev].MutexName {
		t.Error("prod and dev share a mutex name")
	}
	if layouts[buildinfo.EnvTest].MutexEnabled {
		t.Error("test env must not register a mutex")
	}
	if layouts[buildinfo.EnvDev].DefaultLLMBaseURL == "" || layouts[buildinfo.EnvDev].DefaultMirrorURL == "" {
		t.Error("dev must carry the mock-llm/mock-mirror default endpoints")
	}
	if layouts[buildinfo.EnvProd].DefaultLLMBaseURL != "" || layouts[buildinfo.EnvTest].DefaultLLMBaseURL != "" {
		t.Error("prod/test endpoints must not be defaulted (user config / explicit injection)")
	}

	// Mutual invisibility, the actual property debugging depends on: no env's
	// data dir may be nested inside another's.
	for a, da := range dirs {
		for b, db := range dirs {
			if a != b && strings.HasPrefix(da, db+string(filepath.Separator)) {
				t.Errorf("%s data dir %q is nested inside %s data dir %q", a, da, b, db)
			}
		}
	}
}

// TestPortableOverride pins the SPEC-02 §6 portable rule: portable.txt next to
// the exe relocates the data dir to the exe-relative portable dir (dev keeps
// its own data-dev), taking precedence over the env fork dir; an explicit
// WISP_TEST_DATA_DIR injection wins over the marker in the test env.
func TestPortableOverride(t *testing.T) {
	exeDir := t.TempDir()
	root := t.TempDir()

	t.Run("no marker is a no-op", func(t *testing.T) {
		base, err := LayoutFor(buildinfo.EnvProd, root)
		if err != nil {
			t.Fatalf("LayoutFor: %v", err)
		}
		got, applied, err := ApplyPortableOverride(buildinfo.EnvProd, base, exeDir)
		if err != nil || applied {
			t.Fatalf("ApplyPortableOverride = (%+v, %v, %v), want unapplied", got, applied, err)
		}
		if got.DataDir != base.DataDir {
			t.Errorf("data dir changed without the marker: %q", got.DataDir)
		}
	})

	for _, tc := range []struct {
		env      buildinfo.Env
		wantName string
	}{
		{buildinfo.EnvProd, "data"},
		{buildinfo.EnvTest, "data"},
		{buildinfo.EnvDev, "data-dev"},
	} {
		base, err := LayoutFor(tc.env, root)
		if err != nil {
			t.Fatalf("LayoutFor(%s): %v", tc.env, err)
		}
		marker := filepath.Join(exeDir, PortableMarker)
		if err := os.WriteFile(marker, nil, 0o644); err != nil {
			t.Fatalf("write marker: %v", err)
		}
		got, applied, err := ApplyPortableOverride(tc.env, base, exeDir)
		if err != nil {
			t.Fatalf("ApplyPortableOverride(%s): %v", tc.env, err)
		}
		if !applied {
			t.Fatalf("%s: portable marker not applied", tc.env)
		}
		if want := filepath.Join(exeDir, tc.wantName); got.DataDir != want {
			t.Errorf("%s portable DataDir = %q, want %q", tc.env, got.DataDir, want)
		}
		if !got.Portable {
			t.Errorf("%s: Portable flag not set", tc.env)
		}
		if err := os.Remove(marker); err != nil {
			t.Fatalf("remove marker: %v", err)
		}
	}

	t.Run("explicit test data dir beats the marker", func(t *testing.T) {
		injected := t.TempDir()
		t.Setenv(TestDataDirEnv, injected)
		if err := os.WriteFile(filepath.Join(exeDir, PortableMarker), nil, 0o644); err != nil {
			t.Fatalf("write marker: %v", err)
		}
		base, err := LayoutFor(buildinfo.EnvTest, root)
		if err != nil {
			t.Fatalf("LayoutFor(test): %v", err)
		}
		got, applied, err := ApplyPortableOverride(buildinfo.EnvTest, base, exeDir)
		if err != nil || applied {
			t.Fatalf("ApplyPortableOverride(test) = (%+v, %v, %v), want unapplied", got, applied, err)
		}
		if got.DataDir != injected {
			t.Errorf("test DataDir = %q, want injected %q", got.DataDir, injected)
		}
	})
}

// TestSummaryBadge pins the env-identity contract consumed by ball/panel
// (SPEC-03 §5.2): prod shows the bare product name, every other env carries
// the visible "Wisp · <env>" suffix.
func TestSummaryBadge(t *testing.T) {
	for env, want := range map[buildinfo.Env]string{
		buildinfo.EnvProd: "Wisp",
		buildinfo.EnvDev:  "Wisp · dev",
		buildinfo.EnvTest: "Wisp · test",
	} {
		s := Layout{Env: env}.Summary()
		if got := s.EnvBadge(); got != want {
			t.Errorf("EnvBadge(%s) = %q, want %q", env, got, want)
		}
	}

	l, err := LayoutFor(buildinfo.EnvDev, t.TempDir())
	if err != nil {
		t.Fatalf("LayoutFor: %v", err)
	}
	s := l.Summary()
	if s.Env != l.Env || s.DataDir != l.DataDir || s.MutexName != l.MutexName ||
		s.MutexEnabled != l.MutexEnabled || s.DefaultLLMBaseURL != l.DefaultLLMBaseURL ||
		s.AutoUpdateChecks != l.AutoUpdateChecks {
		t.Errorf("Summary() did not project the layout: %+v vs %+v", s, l)
	}
}

// TestLayoutForUnknownEnv rejects anything outside the enum.
func TestLayoutForUnknownEnv(t *testing.T) {
	if _, err := LayoutFor(buildinfo.Env("staging"), t.TempDir()); err == nil {
		t.Fatal("unknown env must fail")
	}
}

// TestDefaultLayoutUnknownEnvFails keeps DefaultLayout from silently
// defaulting an invalid env.
func TestDefaultLayoutUnknownEnvFails(t *testing.T) {
	if _, err := DefaultLayout(buildinfo.Env("nope")); err == nil {
		t.Fatal("DefaultLayout with unknown env must fail")
	}
	if _, err := DefaultLayout(buildinfo.EnvTest); err != nil {
		t.Fatalf("DefaultLayout(test): %v", err)
	}
}
