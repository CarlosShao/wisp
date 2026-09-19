//go:build windows

package proc

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
)

// uniqueForkLayout takes the pinned per-env fork layout and appends a
// per-test-run suffix to the mutex/event names. The test proves the NAME-FORK
// property on real Windows mutexes without ever touching the real prod/dev
// mutexes (a dev machine may genuinely have Wisp running; taking its mutex
// would break the user's session and flake the test).
func uniqueForkLayout(t *testing.T, env buildinfo.Env) Layout {
	t.Helper()
	l, err := LayoutFor(env, t.TempDir())
	if err != nil {
		t.Fatalf("LayoutFor(%s): %v", env, err)
	}
	suffix := fmt.Sprintf("-itest-%d-%d", os.Getpid(), time.Now().UnixNano())
	l.MutexName += suffix
	l.ActivateEventName += suffix
	l.DataDir = t.TempDir()
	return l
}

// TestMutexNamesPerEnv is the ticket-06 integration acceptance (Windows
// runner): the per-env mutex names are actually distinct kernel objects - a
// prod instance and a dev instance coexist (dev debugging never blocks the
// real instance), while the same name still excludes within one env, so the
// isolation comes exactly from the forked names.
func TestMutexNamesPerEnv(t *testing.T) {
	prod := uniqueForkLayout(t, buildinfo.EnvProd)
	dev := uniqueForkLayout(t, buildinfo.EnvDev)

	if prod.MutexName == dev.MutexName {
		t.Fatalf("fork matrix broken: %q == %q", prod.MutexName, dev.MutexName)
	}
	// The names being exercised are the pinned fork names (suffixed, never the
	// bare ones a real instance could own).
	if !strings.HasPrefix(prod.MutexName, `Local\wisp-single-instance`) {
		t.Errorf("prod fork mutex %q does not extend the pinned prod name", prod.MutexName)
	}
	if !strings.HasPrefix(dev.MutexName, `Local\wisp-dev-single-instance`) {
		t.Errorf("dev fork mutex %q does not extend the pinned dev name", dev.MutexName)
	}

	// Prod and dev instances coexist: distinct names, no cross-blocking.
	prodRT, err := Boot(buildinfo.EnvProd, WithLayout(prod), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot(prod fork): %v", err)
	}
	t.Cleanup(func() { prodRT.Shutdown(false) })
	if prodRT.Instance == nil || prodRT.Instance.MutexName != prod.MutexName {
		t.Fatalf("prod instance = %+v, want mutex %q", prodRT.Instance, prod.MutexName)
	}

	devRT, err := Boot(buildinfo.EnvDev, WithLayout(dev), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot(dev fork) while prod holds its mutex: %v", err)
	}
	t.Cleanup(func() { devRT.Shutdown(false) })
	if devRT.Instance == nil || devRT.Instance.MutexName != dev.MutexName {
		t.Fatalf("dev instance = %+v, want mutex %q", devRT.Instance, dev.MutexName)
	}

	// The same name still excludes across env declarations: the isolation is
	// carried by the name fork itself.
	_, err = Boot(buildinfo.EnvDev, WithLayout(prod), WithRegistry(observe.NewRegistry()))
	if !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("Boot with prod's mutex name err = %v, want ErrAlreadyRunning", err)
	}

	// Release, then the name is free again.
	devRT.Shutdown(false)
	free := uniqueForkLayout(t, buildinfo.EnvDev)
	free.MutexName = dev.MutexName
	free.ActivateEventName = dev.ActivateEventName
	rt, err := Boot(buildinfo.EnvDev, WithLayout(free), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot after dev release: %v", err)
	}
	rt.Shutdown(false)
}
