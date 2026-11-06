package main

// Ticket 123: the programmable answer leg the ticket's AC#2 asks for.
//
// The four `cmd/wisp` cases that went red on the hosted runner
// (run 35606321404, step `cmd/wisp CLI tests`) all share one shape: a call
// they dispatch through the composed bridge ends up on the L2 approval
// queue on that machine, the queue item has nobody to answer it, and C18's
// deadline auto-rejects it. On a developer machine the same calls take the
// L1 window route, where an unopposed window means EXECUTE and no human is
// ever involved - which is why the same tree was green locally and red on
// the runner.
//
// The ticket rules out every fix that touches the gate, the deadline or the
// verdict semantics: the repair is to give the TEST a programmable answer
// leg, so "somebody answered" stops being an accident of the machine the
// suite happens to run on. AGENTS.md section 1.3 names the reply stream of
// `wisp run` as the sanctioned injection seam for exactly this, and
// cmd/wisp/approval_reply.go's reply loop is the production side of that
// seam: it reads `yes <corr>` / `no <corr>` lines and routes them onto the
// approval queue. This file is the test-side half: a scripted io.Reader
// that emits its lines on that stream.
//
// Two properties make one script correct across machine shapes:
//
//   - An answer aimed at a correlation that is not (yet) on the queue is
//     refused by the reply loop as unknown (approval_reply.go: a refused
//     reply is printed and the loop keeps going), so a one-shot script
//     would race the card it means to answer. The script therefore
//     re-emits its lines cyclically every t123ReplyInterval until Close.
//     Re-emitting an already answered line is equally harmless: a settled
//     correlation is no longer pending and every later copy is refused the
//     same way.
//   - An L1 confirmation window is not a queue item, so a `yes`/`no` line
//     can never answer one: on the L1 route every line here is refused as
//     an unknown correlation and the window keeps its own polarity
//     (unopposed expiry means EXECUTE). The lines only ever bite on the L2
//     route, which is the route that needs an answerer.

import (
	"io"
	"sync"
	"testing"
	"time"
)

// t123ReplyInterval is how often the script re-emits its next line. Small
// enough that an answered card resolves the call within one interval of the
// card being displayed; large enough that a boot full of refused lines does
// not flood the run's stdout.
const t123ReplyInterval = 250 * time.Millisecond

// t123replyScript is the scripted answer stream. It implements io.Reader,
// which is the type runSpec.reply takes.
type t123replyScript struct {
	ch   chan string
	done chan struct{}
	once sync.Once
}

// newT123ReplyScript starts the feeder for the given lines. Empty script is
// legal and behaves as an already-closed stream (immediate EOF), which keeps
// call sites uniform.
func newT123ReplyScript(lines ...string) *t123replyScript {
	s := &t123replyScript{ch: make(chan string), done: make(chan struct{})}
	if len(lines) == 0 {
		s.Close()
		return s
	}
	// Test-side feeder. The D22 bare-goroutine ban is scoped to production
	// files (tools/d22scan: the walk covers internal/ and cmd/ EXCLUDING
	// _test.go for ban 1), and every existing cmd/wisp test file already
	// feeds its answer streams from plain goroutines on this shape.
	go func() {
		defer close(s.ch)
		i := 0
		for {
			select {
			case <-s.done:
				return
			case s.ch <- lines[i%len(lines)]:
				i++
			}
			select {
			case <-s.done:
				return
			case <-time.After(t123ReplyInterval):
			}
		}
	}()
	return s
}

// Read yields the next scripted line, newline-terminated. After Close it
// yields io.EOF so the reply loop's scanner stands down instead of leaking.
func (s *t123replyScript) Read(p []byte) (int, error) {
	select {
	case <-s.done:
		return 0, io.EOF
	case line, ok := <-s.ch:
		if !ok {
			return 0, io.EOF
		}
		if len(p) < len(line)+1 {
			return 0, io.EOF
		}
		return copy(p, line+"\n"), nil
	}
}

// Close stands the feeder down. The run's own reply listener also stops at
// run end; Close is what keeps the test itself leak-free.
func (s *t123replyScript) Close() { s.once.Do(func() { close(s.done) }) }

// t123script builds a script whose lifetime is the test's: Cleanup closes
// it even when an assertion aborts the case mid-boot.
func t123script(t *testing.T, lines ...string) *t123replyScript {
	t.Helper()
	s := newT123ReplyScript(lines...)
	t.Cleanup(s.Close)
	return s
}
