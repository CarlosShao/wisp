package ball

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Per-monitor ball position persistence (SPEC-08 §2: "位置按显示器分别保存；
// 拔插后落到不可见区域 → 自动回主屏可见位置").
//
// The store is a small JSON file keyed by monitor device name; the file path
// is injected (tests use t.TempDir, production wires the user data dir). The
// visibility fallback is the pure ResolvePosition function - testable
// without any window.

// PosEntry is a saved window position in physical screen pixels (top-left).
type PosEntry struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type posFile struct {
	Version   int                 `json:"version"`
	Positions map[string]PosEntry `json:"positions"`
}

// PositionStore persists ball positions per monitor.
type PositionStore struct {
	mu   sync.Mutex
	path string
	data posFile
}

// OpenPositionStore loads (or initializes) the store at path.
func OpenPositionStore(path string) (*PositionStore, error) {
	s := &PositionStore{path: path, data: posFile{Version: 1, Positions: map[string]PosEntry{}}}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(b, &s.data); err != nil {
		// Corrupt file: start fresh rather than fail the ball.
		s.data = posFile{Version: 1, Positions: map[string]PosEntry{}}
	}
	if s.data.Positions == nil {
		s.data.Positions = map[string]PosEntry{}
	}
	return s, nil
}

// Get returns the saved position for a monitor device name.
func (s *PositionStore) Get(device string) (PosEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.data.Positions[device]
	return p, ok
}

// Put saves (and flushes) the position for a monitor device name.
func (s *PositionStore) Put(device string, e PosEntry) error {
	s.mu.Lock()
	s.data.Positions[device] = e
	b, err := json.MarshalIndent(s.data, "", "  ")
	s.mu.Unlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// MonitorRect is the testable shape of a connected monitor.
type MonitorRect struct {
	Device  string // szDevice, e.g. `\\.\DISPLAY1`
	Primary bool
	Screen  Rect // full monitor rect, physical px
	Work    Rect // work area (minus taskbar), physical px
}

// Rect is an integer rectangle (screen coordinates, physical pixels).
type Rect struct {
	L, T, R, B int
}

func (r Rect) Width() int  { return r.R - r.L }
func (r Rect) Height() int { return r.B - r.T }

func (r Rect) intersects(o Rect) bool {
	return r.L < o.R && o.L < r.R && r.T < o.B && o.T < r.B
}

// DefaultPosition returns the factory spot: bottom-right of the primary work
// area with a small margin (the ball's usual corner).
func DefaultPosition(primary MonitorRect, edge int) PosEntry {
	const margin = 24
	return PosEntry{
		X: primary.Work.R - edge - margin,
		Y: primary.Work.B - edge - margin,
	}
}

// ResolvePosition decides where to place a window of `edge` px:
// saved position when it is visible on the monitor it belongs to, the saved
// monitor's work area clamped inside when partially off-screen, and the
// primary monitor's default spot when the saved monitor is gone (simulated
// detach / topology change, D42#3) or the position lands nowhere visible.
func ResolvePosition(monitors []MonitorRect, savedDevice string, saved PosEntry, savedKnown bool, edge int) (MonitorRect, PosEntry) {
	primary := monitors[0]
	for _, m := range monitors {
		if m.Primary {
			primary = m
			break
		}
	}

	if savedKnown {
		for _, m := range monitors {
			if m.Device != savedDevice {
				continue
			}
			win := Rect{saved.X, saved.Y, saved.X + edge, saved.Y + edge}
			if m.Work.intersects(win) {
				// Clamp fully inside the monitor's work area when it hangs
				// over an edge.
				x, y := saved.X, saved.Y
				if win.R > m.Work.R {
					x = m.Work.R - edge
				}
				if win.B > m.Work.B {
					y = m.Work.B - edge
				}
				if x < m.Work.L {
					x = m.Work.L
				}
				if y < m.Work.T {
					y = m.Work.T
				}
				return m, PosEntry{x, y}
			}
			// Saved monitor exists but the point is off it entirely: fall
			// back to that monitor's default spot.
			return m, DefaultPosition(m, edge)
		}
	}
	// Saved monitor gone (or nothing saved): primary, default spot.
	return primary, DefaultPosition(primary, edge)
}
