package manifestslot

import (
	"context"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The D46 command budget rides on context.WithDeadline, which is the Go
// standard library's monotonic deadline (time.Until on a time.Time that
// carries a monotonic reading). observe.NewTimeout supplies that reading, so
// the budget cannot be shortened or extended by an NTP correction or a
// sleep/resume (D42#1, D22 ban on wall-clock interval math).
//
// No wall-clock subtraction appears anywhere in the command runtime: the
// deadline is enforced by the ctx firing, and the process is killed because the
// ctx was cancelled - never because two timestamps were compared.

// withMonotonicBudget derives a child ctx that expires after d.
func withMonotonicBudget(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	total := observe.NewTimeout(d)
	return context.WithDeadline(parent, total.Deadline())
}
