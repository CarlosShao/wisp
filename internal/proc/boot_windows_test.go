//go:build windows

package proc

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
)

// uniqueLayout builds a prod-style layout with per-test-run names so tests
// never touch the real prod/dev mutexes.
func uniqueLayout(t *testing.T) Layout {
	t.Helper()
	return Layout{
		Env:               buildinfo.EnvProd,
		DataDir:           t.TempDir(),
		MutexName:         fmt.Sprintf(`Local\wisp-test-boot-%d-%d`, os.Getpid(), time.Now().UnixNano()),
		ActivateEventName: fmt.Sprintf(`Local\wisp-test-boot-%d-%d.activate`, os.Getpid(), time.Now().UnixNano()),
		MutexEnabled:      true,
	}
}

// TestBootTestEnvBootsAllPackages: the ticket-03 init-time self-checks pass
// for the test env (real data dir, no mutex) and shutdown walks the frozen
// 10-step order with the Job close executed.
func TestBootTestEnvBootsAllPackages(t *testing.T) {
	rt, err := Boot(buildinfo.EnvTest, WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot(test): %v", err)
	}
	if rt.Layout.DataDir == "" {
		t.Fatal("test env must resolve a concrete data dir (ticket 06)")
	}
	if rt.Instance != nil {
		t.Fatal("test env must not acquire a mutex")
	}
	if rt.Job == nil || rt.Job.Closed() {
		t.Fatal("Job Object missing after boot")
	}
	if rt.Registry == nil {
		t.Fatal("registry missing after boot")
	}

	records := rt.Shutdown(false)
	if len(records) != 10 {
		t.Fatalf("shutdown records = %d, want 10", len(records))
	}
	jobRec := records[StepCloseJob-1]
	if jobRec.Step != StepCloseJob || jobRec.Skipped || jobRec.Err != nil {
		t.Fatalf("job close record = %+v", jobRec)
	}
	if !rt.Job.Closed() {
		t.Fatal("job still open after shutdown")
	}
	if rt.Instance != nil {
		t.Fatal("instance should stay nil in test env")
	}
}

// TestBootSingleInstanceConflict runs the D42#7 second-instance rule through
// Boot itself: the second Boot in the same session detects the first and
// reports ErrAlreadyRunning; after the first shuts down, Boot succeeds again.
func TestBootSingleInstanceConflict(t *testing.T) {
	layout := uniqueLayout(t)

	first, err := Boot(buildinfo.EnvProd, WithLayout(layout), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("first Boot: %v", err)
	}
	if first.Instance == nil {
		t.Fatal("first boot must own the mutex")
	}

	_, err = Boot(buildinfo.EnvProd, WithLayout(layout), WithRegistry(observe.NewRegistry()))
	if !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("second Boot err = %v, want ErrAlreadyRunning", err)
	}

	first.Shutdown(false)
	_, err = Boot(buildinfo.EnvProd, WithLayout(layout), WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("Boot after first shutdown: %v", err)
	}
}
