package tools

import (
	"context"
	"fmt"
)

// CancelBus is the D31 seam between this package and the approval layer's veto
// bookkeeping. It exists because tools.Gate's two callbacks are PRE-execution
// only: once the bridge has handed off, nothing in the gate can reach into the
// running tool, and nothing in the tool can tell whether the user tried to
// stop it.
//
// The single applied-steps shape in this repository is
// approval.CancellationReport (internal/agent/approval/report.go), so this
// package deliberately does NOT define another one. Report returns an opaque
// fmt.Stringer for exactly that reason: the bridge can render and forward the
// host's report without owning a second representation of it, and the adapter
// that plugs approval.Gate.Bus() in lives on the approval side (which may
// import this package, not the other way round).
type CancelBus interface {
	// Vetoed reports whether a veto arrived for this call (before or after
	// handoff - the bridge only ever asks about a started call).
	Vetoed(correlationID string) bool
	// Started reports whether the call was handed off to execution.
	Started(correlationID string) bool
	// Report builds the applied-steps report for one finished call.
	Report(d Decision, res Result) fmt.Stringer
	// Complete releases one call's tracking; the bridge defers it.
	Complete(correlationID string)
}

// cancelKey is the context key carrying the running call's cancel handle.
type cancelKey struct{}

// cancelHandle binds a running tool to its own correlation id. A nil bus means
// "no approval layer wired", which makes every veto answer false while still
// letting the tool read its correlation id for logging.
type cancelHandle struct {
	corr string
	bus  CancelBus
}

func withCancel(ctx context.Context, h cancelHandle) context.Context {
	return context.WithValue(ctx, cancelKey{}, h)
}

func cancelOf(ctx context.Context) (cancelHandle, bool) {
	h, ok := ctx.Value(cancelKey{}).(cancelHandle)
	return h, ok
}

// CorrelationID returns the id of the call the running tool belongs to, or ""
// when the tool was invoked outside the bridge.
func CorrelationID(ctx context.Context) string {
	h, _ := cancelOf(ctx)
	return h.corr
}

// Vetoed reports whether the user vetoed this call. A running side-effecting
// tool polls it at every stoppable step boundary (approval.Bus contract,
// step 1). False when no cancel bus is wired.
func Vetoed(ctx context.Context) bool {
	h, ok := cancelOf(ctx)
	if !ok || h.bus == nil {
		return false
	}
	return h.bus.Vetoed(h.corr)
}

// Stopped is the one check a step boundary needs: either the task's context is
// gone or the user vetoed. The reason is user-visible wording, never empty
// when it reports true.
func Stopped(ctx context.Context) (bool, string) {
	if err := ctx.Err(); err != nil {
		return true, "上下文已结束：" + err.Error()
	}
	if Vetoed(ctx) {
		return true, "用户在确认窗口结束后仍投了否决票（D31：取消不是原子的）"
	}
	return false, ""
}
