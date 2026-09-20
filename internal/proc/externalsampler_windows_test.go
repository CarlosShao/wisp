//go:build windows

package proc

import (
	"os"
	"testing"
	"time"
)

// Ticket 66 AC#3: the out-of-tree reader is what puts the D32 CPU gate back
// into a state where it means something. These tests pin the two properties
// the gate depends on - the subject's counters really come from outside it,
// and the reader refuses to produce a number when the view is not trustworthy.

func TestExternalSamplerReadsSubjectFromOutside(t *testing.T) {
	h := startHelper(t, nil, "job-child")
	t.Cleanup(func() {
		_ = h.cmd.Process.Kill()
		_ = h.cmd.Wait()
	})
	pid := uint32(h.cmd.Process.Pid)

	r := NewExternalSampler(pid)
	if got := r.PIDs(); len(got) != 1 || got[0] != pid {
		t.Fatalf("PIDs() = %v, want [%d]", got, pid)
	}
	m, err := r.ReadTree()
	if err != nil {
		t.Fatalf("ReadTree: %v", err)
	}
	if m.PIDs != 1 {
		t.Errorf("PIDs = %d, want 1 (the subject only, never the observer)", m.PIDs)
	}
	if m.PrivateWorkingSetBytes <= 0 {
		t.Errorf("subject private working set = %d, want > 0", m.PrivateWorkingSetBytes)
	}
	if m.Handles <= 0 || m.Threads <= 0 {
		t.Errorf("subject decoded as handles=%d threads=%d, want both > 0", m.Handles, m.Threads)
	}
	if m.CPUTotalNanos < 0 {
		t.Errorf("subject CPU nanos = %d, want >= 0", m.CPUTotalNanos)
	}
	// The observer is not in the measured set, so it can never bill its own
	// ReadTree cost (or its own scratch buffer) into the subject's numbers.
	if m.SelfWriteOps != 0 {
		t.Errorf("SelfWriteOps = %d, want 0 for an out-of-tree read", m.SelfWriteOps)
	}
	if m.TCPConnections != 0 {
		t.Errorf("subject TCP connections = %d, want 0 for a test helper", m.TCPConnections)
	}
}

func TestExternalSamplerRefusesToMeasureItself(t *testing.T) {
	r := NewExternalSampler(uint32(os.Getpid()))
	if _, err := r.ReadTree(); err == nil {
		t.Fatal("out-of-tree read of the observer itself must fail (that is the A15 self-measurement)")
	}
}

func TestExternalSamplerFailClosedOnDeadSubject(t *testing.T) {
	// A pid that cannot exist: the sampler must error, never return a
	// zero-footprint reading that a gate could mistake for an idle subject.
	r := NewExternalSampler(0xFFFFFFF0)
	if _, err := r.ReadTree(); err == nil {
		t.Fatal("missing subject must fail closed")
	}
	if _, err := NewExternalSampler().ReadTree(); err == nil {
		t.Fatal("an empty subject set must fail closed")
	}
}

func TestExternalSamplerTracksSubjectGrowth(t *testing.T) {
	// Two reads of a live subject, ~250ms apart, must both succeed and report
	// a non-decreasing CPU total: this is the exact pair the window mean is
	// built from, so a reader that returns a constant or a negative delta here
	// would silently report 0% for any process.
	h := startHelper(t, nil, "job-child")
	t.Cleanup(func() {
		_ = h.cmd.Process.Kill()
		_ = h.cmd.Wait()
	})
	r := NewExternalSampler(uint32(h.cmd.Process.Pid))
	first, err := r.ReadTree()
	if err != nil {
		t.Fatalf("first ReadTree: %v", err)
	}
	time.Sleep(250 * time.Millisecond)
	second, err := r.ReadTree()
	if err != nil {
		t.Fatalf("second ReadTree: %v", err)
	}
	if second.CPUTotalNanos < first.CPUTotalNanos {
		t.Errorf("subject CPU total went backwards: %d -> %d", first.CPUTotalNanos, second.CPUTotalNanos)
	}
	if second.PrivateWorkingSetBytes <= 0 {
		t.Errorf("second read private working set = %d, want > 0", second.PrivateWorkingSetBytes)
	}
}
