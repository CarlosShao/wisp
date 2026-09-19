package proc

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// Layout is the per-environment resource fork (SPEC-03 §5.2). All three envs
// resolve to a complete layout here; ticket 06 completed the test-env fork
// (data dir via WISP_TEST_DATA_DIR or %TEMP%\wisp-test-<pid>, no mutex,
// injectable endpoints) and the portable-mode override (SPEC-02 §6).
type Layout struct {
	Env               buildinfo.Env
	DataDir           string
	MutexName         string // "" when MutexEnabled is false
	ActivateEventName string // "" when MutexEnabled is false
	MutexEnabled      bool   // test env registers no mutex (parallel tests)
	Portable          bool   // portable.txt override applied (SPEC-02 §6)

	// Field-level per-env defaults (SPEC-03 §5.2). Consumption lands in later
	// tickets: the endpoints feed provider defaults (09/11), AutoUpdateChecks
	// feeds the update scheduler (44/56).
	DefaultLLMBaseURL string // dev: http://127.0.0.1:18080/v1; prod/test: ""
	DefaultMirrorURL  string // dev: http://127.0.0.1:18081; prod/test: ""
	AutoUpdateChecks  bool   // prod: true; dev/test: false
}

// Injected / marker names (SPEC-03 §5.2, SPEC-02 §6).
const (
	// TestDataDirEnv is the test-env data-dir override (explicit injection).
	TestDataDirEnv = "WISP_TEST_DATA_DIR"

	// PortableMarker is the file whose presence next to the exe switches the
	// whole data dir to the exe-relative portable dir.
	PortableMarker = "portable.txt"

	// DevLLMBaseURL / DevMirrorBaseURL are the dev-env default endpoints
	// (compose mock-llm / model-mirror, SPEC-11).
	DevLLMBaseURL    = "http://127.0.0.1:18080/v1"
	DevMirrorBaseURL = "http://127.0.0.1:18081"
)

// LayoutFor maps an environment onto its resource names, given the user
// config root (os.UserConfigDir, i.e. %APPDATA% on Windows; injected as a
// parameter so this stays a pure function). The `Local\` namespace makes the
// mutex per-session (D42#7): a second logged-in user can run their own Wisp.
// The portable-mode override (portable.txt -> exe-relative data dir) is
// applied by DefaultLayout / ApplyPortableOverride and takes precedence over
// the DataDir chosen here (SPEC-02 §6).
func LayoutFor(env buildinfo.Env, userConfigRoot string) (Layout, error) {
	switch env {
	case buildinfo.EnvProd:
		return Layout{
			Env:               env,
			DataDir:           filepath.Join(userConfigRoot, "wisp"),
			MutexName:         `Local\wisp-single-instance`,
			ActivateEventName: `Local\wisp-single-instance.activate`,
			MutexEnabled:      true,
			AutoUpdateChecks:  true,
		}, nil
	case buildinfo.EnvDev:
		return Layout{
			Env:               env,
			DataDir:           filepath.Join(userConfigRoot, "wisp-dev"),
			MutexName:         `Local\wisp-dev-single-instance`,
			ActivateEventName: `Local\wisp-dev-single-instance.activate`,
			MutexEnabled:      true,
			DefaultLLMBaseURL: DevLLMBaseURL,
			DefaultMirrorURL:  DevMirrorBaseURL,
			AutoUpdateChecks:  false, // dev machines are not update guinea pigs
		}, nil
	case buildinfo.EnvTest:
		// No mutex registration (tests run in parallel, SPEC-03 §5.2) and no
		// endpoint defaults: the harness injects them explicitly. The data
		// dir comes from WISP_TEST_DATA_DIR or %TEMP%\wisp-test-<pid>.
		return Layout{
			Env:          env,
			DataDir:      TestDataDir(),
			MutexEnabled: false,
		}, nil
	default:
		return Layout{}, fmt.Errorf("proc: unknown env %q", env)
	}
}

// TestDataDir resolves the test-env data dir (SPEC-03 §5.2): the
// WISP_TEST_DATA_DIR injection when set, else %TEMP%\wisp-test-<pid> (one dir
// per test process, so parallel test binaries never share state).
func TestDataDir() string {
	if dir := os.Getenv(TestDataDirEnv); dir != "" {
		return dir
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("wisp-test-%d", os.Getpid()))
}

// PortableDataDirName is the exe-relative data dir under portable mode: dev
// keeps its own dir ("data-dev") so debugging never touches real portable
// data (SPEC-02 §6 + SPEC-03 §5.2).
func PortableDataDirName(env buildinfo.Env) (string, error) {
	switch env {
	case buildinfo.EnvDev:
		return "data-dev", nil
	case buildinfo.EnvProd, buildinfo.EnvTest:
		return "data", nil
	default:
		return "", fmt.Errorf("proc: unknown env %q", env)
	}
}

// ApplyPortableOverride applies the portable-mode rule (SPEC-02 §6): when
// <exeDir>\portable.txt exists, the layout's data dir becomes the exe-relative
// portable dir and takes precedence over the per-env fork dir. Precedence
// within the test env: an explicit WISP_TEST_DATA_DIR injection wins over the
// marker so CI stays deterministic regardless of where the test binary sits.
// Returns (layout, portable=true) when the override applied.
func ApplyPortableOverride(env buildinfo.Env, l Layout, exeDir string) (Layout, bool, error) {
	if env == buildinfo.EnvTest && os.Getenv(TestDataDirEnv) != "" {
		return l, false, nil
	}
	name, err := PortableDataDirName(env)
	if err != nil {
		return l, false, err
	}
	marker := filepath.Join(exeDir, PortableMarker)
	if _, err := os.Stat(marker); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return l, false, fmt.Errorf("proc: portable marker %s: %w", marker, err)
		}
		return l, false, nil
	}
	l.DataDir = filepath.Join(exeDir, name)
	l.Portable = true
	return l, true, nil
}

// DefaultLayout resolves LayoutFor against the OS user config dir and then
// applies the portable-mode override (portable.txt next to the running exe;
// SPEC-02 §6). It is the single resolution order used by Boot, Summarize and
// the CLI.
func DefaultLayout(env buildinfo.Env) (Layout, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Layout{}, fmt.Errorf("proc: user config dir: %w", err)
	}
	l, err := LayoutFor(env, dir)
	if err != nil {
		return Layout{}, err
	}
	if exe, err := os.Executable(); err == nil {
		l, _, err = ApplyPortableOverride(env, l, filepath.Dir(exe))
		if err != nil {
			return Layout{}, err
		}
		// os.Executable failing is not fatal here: without it there is no exe
		// dir to find a portable.txt in, so the fork dir above stands.
	}
	return l, nil
}

// Summary is the environment-identity + data-dir digest consumed by the ball
// tooltip and panel title (SPEC-03 §5.2: the env must be visible as
// "Wisp · dev" so debugging never mistakes dev data for real data; badge
// rendering is tickets 07/35).
type Summary struct {
	Env               buildinfo.Env
	DataDir           string
	Portable          bool
	MutexName         string
	MutexEnabled      bool
	DefaultLLMBaseURL string
	DefaultMirrorURL  string
	AutoUpdateChecks  bool
}

// Summary projects the layout onto the UI-facing digest.
func (l Layout) Summary() Summary {
	return Summary{
		Env:               l.Env,
		DataDir:           l.DataDir,
		Portable:          l.Portable,
		MutexName:         l.MutexName,
		MutexEnabled:      l.MutexEnabled,
		DefaultLLMBaseURL: l.DefaultLLMBaseURL,
		DefaultMirrorURL:  l.DefaultMirrorURL,
		AutoUpdateChecks:  l.AutoUpdateChecks,
	}
}

// Summarize resolves the effective environment layout (per-env fork + portable
// override) and summarizes it for the UI badge surfaces.
func Summarize(env buildinfo.Env) (Summary, error) {
	l, err := DefaultLayout(env)
	if err != nil {
		return Summary{}, err
	}
	return l.Summary(), nil
}

// EnvBadge is the visible env identity (SPEC-03 §5.2). prod shows the bare
// product name: the "· <env>" suffix always means "this is not your real
// data". Product name is the frozen brand string until buildinfo carries a
// display name.
func (s Summary) EnvBadge() string {
	if s.Env == buildinfo.EnvProd || s.Env == "" {
		return "Wisp"
	}
	return "Wisp · " + string(s.Env)
}
