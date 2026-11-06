package config

// Ticket 65: the [ball] style key (D36 hot key; owner-approved direction
// A1581, docs/evidence/s1/65-style-direction.md §4). The vocabulary lives in
// two packages on purpose - internal/config must not import internal/ball -
// so this file carries the drift pin in BOTH directions: the config enum
// equals ball.StyleNames(), the default is glass, and an unknown value is a
// load error naming the vocabulary (never a silent fall back).

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
)

// TestBallStyleVocabularyMatchesBall pins the two-vocabulary drift: the
// config enum and ball.StyleNames() must name exactly the same families in
// the same order, so a style added on one side cannot silently pass load
// while the renderer (or the panel) cannot say it.
func TestBallStyleVocabularyMatchesBall(t *testing.T) {
	got := vocabList(ballStyles)
	want := strings.Join(ball.StyleNames(), "|")
	if got != want {
		t.Fatalf("config ballStyles = %q, want ball.StyleNames() = %q - the two vocabularies drifted", got, want)
	}
}

// TestBallStyleDefaultIsGlass: a config without a style line (every file
// written before ticket 65) loads as glass - the default look the owner
// signed into the quantified spec.
func TestBallStyleDefaultIsGlass(t *testing.T) {
	c := validConfig(t)
	if c.Ball.Style != BallStyleGlass {
		t.Fatalf("ball.style default = %q, want %q", c.Ball.Style, BallStyleGlass)
	}
	for _, good := range ballStyles {
		c := validConfig(t)
		c.Ball.Style = good
		if err := validate(c); err != nil {
			t.Fatalf("ball.style %q must validate: %v", good, err)
		}
	}
}

// TestValidateBallStyleEnum: an unknown style is refused at load with the
// vocabulary in the message (same rule as every other enum in validate.go -
// "it fell back" and "it was never written" must look different).
func TestValidateBallStyleEnum(t *testing.T) {
	for _, bad := range []string{"aurora", "nebula", "", "GLASS", "Glass ", "jarvis2"} {
		c := validConfig(t)
		c.Ball.Style = bad
		err := validate(c)
		if err == nil {
			t.Fatalf("ball.style %q must be rejected", bad)
		}
		if !strings.Contains(err.Error(), "ball.style") || !strings.Contains(err.Error(), "glass|jarvis") {
			t.Fatalf("the rejection must name the key and the vocabulary, got: %v", err)
		}
	}
}
