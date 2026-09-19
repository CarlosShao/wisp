package proc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/buildinfo"
)

// Layout is the per-environment resource fork (SPEC-03 §5.2). Ticket 03
// freezes the shape plus the prod/dev defaults; the test-env values (data
// dir via WISP_TEST_DATA_DIR or %TEMP%\wisp-test-<pid>, injectable endpoints,
// portable-mode override precedence) are decided and implemented by ticket 06
// and are deliberately NOT preimplemented here.
type Layout struct {
	Env               buildinfo.Env
	DataDir           string // resolved for prod/dev; "" for test until ticket 06
	MutexName         string // "" when MutexEnabled is false
	ActivateEventName string // "" when MutexEnabled is false
	MutexEnabled      bool   // test env registers no mutex (parallel tests)
}

// ErrTestLayoutDeferred marks the test-env fork as ticket-06 scope.
var ErrTestLayoutDeferred = errors.New(
	"proc: test-env layout (data dir, injectable endpoints) is decided by ticket 06")

// LayoutFor maps an environment onto its resource names, given the user
// config root (os.UserConfigDir, i.e. %APPDATA% on Windows; injected as a
// parameter so this stays a pure function). The `Local\` namespace makes the
// mutex per-session (D42#7): a second logged-in user can run their own Wisp.
// Portable mode (portable.txt -> exe-relative data dir) is ticket 06 and
// takes precedence over these defaults there.
func LayoutFor(env buildinfo.Env, userConfigRoot string) (Layout, error) {
	switch env {
	case buildinfo.EnvProd:
		return Layout{
			Env:               env,
			DataDir:           filepath.Join(userConfigRoot, "wisp"),
			MutexName:         `Local\wisp-single-instance`,
			ActivateEventName: `Local\wisp-single-instance.activate`,
			MutexEnabled:      true,
		}, nil
	case buildinfo.EnvDev:
		return Layout{
			Env:               env,
			DataDir:           filepath.Join(userConfigRoot, "wisp-dev"),
			MutexName:         `Local\wisp-dev-single-instance`,
			ActivateEventName: `Local\wisp-dev-single-instance.activate`,
			MutexEnabled:      true,
		}, nil
	case buildinfo.EnvTest:
		// No mutex registration in test env (SPEC-03 §5.2). DataDir stays
		// empty: the concrete value (WISP_TEST_DATA_DIR / %TEMP%\wisp-test-<pid>)
		// is ticket 06.
		return Layout{Env: env, MutexEnabled: false}, ErrTestLayoutDeferred
	default:
		return Layout{}, fmt.Errorf("proc: unknown env %q", env)
	}
}

// DefaultLayout resolves LayoutFor against the OS user config dir.
func DefaultLayout(env buildinfo.Env) (Layout, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return Layout{}, fmt.Errorf("proc: user config dir: %w", err)
	}
	return LayoutFor(env, dir)
}
