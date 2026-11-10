package config

// Ticket 413's settings tests: the guarded [ball].size setter and the read
// face the size slider reads. The band itself (44..160, A1636) and the load
// gate's sentence are pinned in validate_test.go's own band test; these cover
// the WRITE path the slider rides and the round trip AC#4 names (write 90,
// read back 90). The band constants and the refusal sentence are single-sourced
// in validate.go (ballSizeMinPx/ballSizeMaxPx/ballSizeRangeError), so the
// exact-sentence assertions below pin the pairing, not a copy.

import (
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TestSetBallSizeWritesThroughGuardedPath is AC#4's round trip, config half: a
// 90 write lands on disk under the owned key path, the receipt names ball.size,
// and the read face answers 90 back.
func TestSetBallSizeWritesThroughGuardedPath(t *testing.T) {
	m, path, _ := newTestManager(t, matrixBase)
	written, err := m.SetBallSize(90)
	if err != nil {
		t.Fatalf("SetBallSize(90): %v", err)
	}
	if len(written) == 0 || written[0] != "ball.size" {
		t.Fatalf("receipt = %v, want [ball.size]", written)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// The guarded writer re-lays the file in canonical form; the key line is
	// what the next start (and ballSizeSeed413) will load.
	if !strings.Contains(string(b), "size = 90") {
		t.Fatalf("config.toml does not carry size = 90 after the write:\n%s", b)
	}
	// The read face answers what the process would act on: the written size.
	if v := m.BallSizeOptions(); v.Size != 90 {
		t.Fatalf("BallSizeOptions().Size = %d, want 90 (in-memory view)", v.Size)
	}
}

// TestSetBallSizeRefusesOutsideRange is AC#4's refusal half: 43 and 161 (and
// 0, the "unset" shape a slider could hand back) are refused BEFORE the write
// with validate.go's own sentence verbatim, and config.toml stays
// byte-identical. The Detail equality is deliberate and is the reverse
// criterion: deleting the setter's range guard turns this red, because the
// refusal would then come back as the generic whole-validate wrapper instead
// of the setter's own named sentence.
func TestSetBallSizeRefusesOutsideRange(t *testing.T) {
	m, path, _ := newTestManager(t, matrixBase)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []int{43, 161, 0} {
		written, err := m.SetBallSize(bad)
		if err == nil {
			t.Fatalf("SetBallSize(%d) must be refused", bad)
		}
		oe, ok := err.(*observe.Error)
		if !ok {
			t.Fatalf("SetBallSize(%d) refusal is %T, want *observe.Error", bad, err)
		}
		want := validateGoSentence413(bad)
		if oe.Detail != want {
			t.Fatalf("SetBallSize(%d) refusal detail = %q, want validate.go's own sentence %q verbatim", bad, oe.Detail, want)
		}
		if written != nil {
			t.Fatalf("a refused write must return no receipt, got %v", written)
		}
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatalf("a refused write must leave config.toml byte-identical")
	}
}

// validateGoSentence413 spells the load gate's own sentence for one value -
// the same string ballSizeRangeError formats, spelled out here so the test
// pins the sentence against a literal, not against the function it asserts.
func validateGoSentence413(v int) string {
	switch v {
	case 43:
		return "config.toml: ball.size 43 out of range [44, 160]"
	case 161:
		return "config.toml: ball.size 161 out of range [44, 160]"
	default:
		return "config.toml: ball.size 0 out of range [44, 160]"
	}
}

// TestBallSizeOptionsReader: the reader answers the file's own size and the
// schema band (the constants validate.go judges by, never a copy), and a nil
// Manager answers Size 0 with the band intact - a caller answers "not
// configured", never a made-up diameter.
func TestBallSizeOptionsReader(t *testing.T) {
	m, _, _ := newTestManager(t, matrixBase) // matrixBase carries [ball] size = 56
	v := m.BallSizeOptions()
	if v.Size != 56 {
		t.Fatalf("read = %d, want the file's own 56", v.Size)
	}
	if v.Min != 44 || v.Max != 160 {
		t.Fatalf("band = [%d, %d], want [44, 160] (the schema band validate.go judges by)", v.Min, v.Max)
	}
	var nilManager *Manager
	nv := nilManager.BallSizeOptions()
	if nv.Size != 0 || nv.Min != 44 || nv.Max != 160 {
		t.Fatalf("nil Manager read = %+v, want Size 0 with the band intact", nv)
	}
}
