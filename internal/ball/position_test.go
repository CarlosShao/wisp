package ball

import (
	"os"
	"path/filepath"
	"testing"
)

func monitorsFixture() []MonitorRect {
	return []MonitorRect{
		{Device: `\\.\DISPLAY1`, Primary: true, Screen: Rect{0, 0, 1920, 1080}, Work: Rect{0, 0, 1920, 1040}},
		{Device: `\\.\DISPLAY2`, Screen: Rect{1920, 0, 3200, 1080}, Work: Rect{1920, 0, 3200, 1040}},
	}
}

// TestPositionStoreRoundTrip: save/load per device; corrupt file recovers.
func TestPositionStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ballpos.json")
	s, err := OpenPositionStore(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, ok := s.Get(`\\.\DISPLAY2`); ok {
		t.Fatal("fresh store must be empty")
	}
	want := PosEntry{X: 2000, Y: 300}
	if err := s.Put(`\\.\DISPLAY2`, want); err != nil {
		t.Fatalf("put: %v", err)
	}

	// Re-open (fresh process view).
	s2, err := OpenPositionStore(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	got, ok := s2.Get(`\\.\DISPLAY2`)
	if !ok || got != want {
		t.Fatalf("round trip: got %+v ok=%v, want %+v", got, ok, want)
	}
	if _, ok := s2.Get(`\\.\DISPLAY1`); ok {
		t.Fatal("per-monitor isolation broken")
	}

	// Corrupt file: recover with an empty store, no error.
	corrupt := filepath.Join(t.TempDir(), "bad.json")
	if err := writeFile(corrupt, "{not json"); err != nil {
		t.Fatal(err)
	}
	s3, err := OpenPositionStore(corrupt)
	if err != nil {
		t.Fatalf("corrupt open: %v", err)
	}
	if _, ok := s3.Get("x"); ok {
		t.Fatal("corrupt store must start empty")
	}
}

func writeFile(path, content string) error {
	return writeFileBytes(path, []byte(content))
}

// TestResolvePosition covers the multi-monitor restore matrix (SPEC-08 §2:
// 位置按显示器保存；不可见 → 回主屏).
func TestResolvePosition(t *testing.T) {
	mons := monitorsFixture()
	edge := 70

	t.Run("saved visible on its monitor: kept", func(t *testing.T) {
		m, p := ResolvePosition(mons, `\\.\DISPLAY2`, PosEntry{2500, 200}, true, edge)
		if m.Device != `\\.\DISPLAY2` || p.X != 2500 || p.Y != 200 {
			t.Fatalf("got %s %+v", m.Device, p)
		}
	})

	t.Run("saved partially off right edge: clamped", func(t *testing.T) {
		_, p := ResolvePosition(mons, `\\.\DISPLAY2`, PosEntry{3180, 200}, true, edge)
		if p.X != 3200-edge {
			t.Fatalf("clamp X = %d, want %d", p.X, 3200-edge)
		}
	})

	t.Run("saved monitor detached: primary default", func(t *testing.T) {
		_, p := ResolvePosition(mons, `\\.\DISPLAY9`, PosEntry{9999, 9999}, true, edge)
		want := DefaultPosition(mons[0], edge)
		if p != want {
			t.Fatalf("fallback = %+v, want %+v", p, want)
		}
	})

	t.Run("nothing saved: primary default", func(t *testing.T) {
		_, p := ResolvePosition(mons, "", PosEntry{}, false, edge)
		want := DefaultPosition(mons[0], edge)
		if p != want {
			t.Fatalf("default = %+v, want %+v", p, want)
		}
	})

	t.Run("saved point on no monitor: that-monitor default impossible -> primary", func(t *testing.T) {
		// Device still exists but the point drifted far off it: the ball
		// lands at that monitor's default spot (visible > remembered).
		_, p := ResolvePosition(mons, `\\.\DISPLAY2`, PosEntry{-5000, -5000}, true, edge)
		want := DefaultPosition(mons[1], edge)
		if p != want {
			t.Fatalf("monitor default = %+v, want %+v", p, want)
		}
	})

	t.Run("primary detection when listed second", func(t *testing.T) {
		rev := []MonitorRect{mons[1], mons[0]}
		_, p := ResolvePosition(rev, "", PosEntry{}, false, edge)
		want := DefaultPosition(rev[1], edge)
		if p != want {
			t.Fatalf("primary default = %+v, want %+v", p, want)
		}
	})
}

// TestDefaultPositionInsideWorkArea guards the factory spot.
func TestDefaultPositionInsideWorkArea(t *testing.T) {
	m := monitorsFixture()[0]
	p := DefaultPosition(m, 70)
	if p.X < m.Work.L || p.Y < m.Work.T || p.X+70 > m.Work.R || p.Y+70 > m.Work.B {
		t.Fatalf("default spot %+v outside work area %+v", p, m.Work)
	}
}

// writeFileBytes is the os helper shared by the store tests.
func writeFileBytes(path string, b []byte) error {
	return os.WriteFile(path, b, 0o644)
}
