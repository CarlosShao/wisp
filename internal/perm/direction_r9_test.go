package perm

// R9's shape pin: the predicate that decides "does this switch need a card" is the
// DIRECTION, and the three-way gate is one function with one channel term.
//
// WHY A SEPARATE SOURCE-LEVEL PIN. store_r9_test.go measures what Set DOES, and a
// behaviour test can only see the paths it exercises. The reverse mutation R9
// names - putting RequiresCard back to `to == risk.ModeAutoApprove` - would leave
// the auto_approve rows of that table green and only the middle-tier rows red, and
// a future change could reintroduce the target rule as an extra term here while
// every door still agrees. This file pins the definition itself, over the complete
// 3x3 matrix, so the middle tier is a first-class cell rather than a side effect.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// r9Ladder is the mode vocabulary in risk.Mode's declared order.
func r9Ladder() []risk.Mode {
	return []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove}
}

// TestR9RequiresCardIsTheDirection pins RequiresCard == IsWidening for all nine
// ordered pairs, and that the no-change cells cost nothing.
func TestR9RequiresCardIsTheDirection(t *testing.T) {
	ladder := r9Ladder()
	if got := len(risk.ModeNames()); got != len(ladder) {
		t.Fatalf("the vocabulary holds %d names, this pin walks %d", got, len(ladder))
	}
	cells := 0
	for _, from := range ladder {
		for _, to := range ladder {
			cells++
			if got, want := RequiresCard(from, to), IsWidening(from, to); got != want {
				t.Errorf("RequiresCard(%s, %s) = %v, want IsWidening = %v: the cost rule must BE the "+
					"direction rule (R9)", from, to, got, want)
			}
			if from == to && RequiresCard(from, to) {
				t.Errorf("RequiresCard(%s, %s) = true; a no-change costs nothing", from, to)
			}
			if from != to && RequiresCard(from, to) == RequiresCard(to, from) {
				t.Errorf("%s <-> %s: both directions read the same cost (%v); a loosening and a "+
					"tightening cannot be the same rule", from, to, RequiresCard(from, to))
			}
		}
	}
	if cells != 9 {
		t.Fatalf("walked %d cells, want 9 (3x3)", cells)
	}
}

// TestR9JudgeCardIsOneThreeWayGate pins the gate as a function of exactly two
// facts - the switch and whether this layer holds a channel - for all nine pairs.
// A loosening without a channel is GateNoChannel, with one is GateAsk; a
// tightening or a no-change is GateNoCard either way.
func TestR9JudgeCardIsOneThreeWayGate(t *testing.T) {
	ladder := r9Ladder()
	for _, from := range ladder {
		for _, to := range ladder {
			loosens := IsWidening(from, to)
			for _, channel := range []bool{false, true} {
				want := GateNoCard
				if loosens && !channel {
					want = GateNoChannel
				} else if loosens && channel {
					want = GateAsk
				}
				if got := JudgeCard(from, to, channel); got != want {
					t.Errorf("JudgeCard(%s, %s, channel=%v) = %v, want %v", from, to, channel, got, want)
				}
			}
		}
	}
}
