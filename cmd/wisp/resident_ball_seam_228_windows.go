//go:build windows

package main

// Ticket 228 AC#9: the seam that makes the ball-REFUSED arm reachable.
//
// Every symbol here exists only so the AC#9 case walks that arm; nothing in this
// file runs unless a test calls it, and the production assembly root never passes
// it (startResidentBall's newBall hook stays nil there, which is ball.New).
//
// WHY A SEAM AT ALL. The refusal branch is where the ticket's double-statement
// defect lived, and on a developer's desktop the ball always comes up - so before
// this seam the branch could only be seen by editing the source, which is exactly
// what the ticket reports ("只有突变能现形、CI 里零读数"). A defect CI cannot watch is
// a defect that comes back, so the arm gets a way in: a factory that answers the
// error ball.New would answer on a host with no desktop, through the SAME
// production path (verdict, booking, boot report), with no desktop required.

import "github.com/CarlosShao/wisp/internal/ball"

// withBallNew228 hands the ball host a factory answering fn's built ball. An
// error makes startResidentBall take its refused arm for real.
func withBallNew228(fn func(ball.Options) (*ball.Ball, error)) ballHostHook {
	return func(h *panelHostHooks) { h.newBall = fn }
}
