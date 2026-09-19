package memory

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// db-writer — the single-writer goroutine (D35 rule 3, D38b).
//
// Every write transaction is enqueued here; nothing else touches the writer
// connection. No SQLITE_BUSY retry exists anywhere: serialization happens by
// construction (one goroutine + one connection + BEGIN IMMEDIATE), readers
// never block under WAL.
//
// Lifecycle (the roster lists db-writer among the 6 resident names, but D20
// forbids resident connections; keeping the goroutine itself non-resident
// when the store is idle keeps RosterReport at or below the frozen baseline 6
// without any roster exemption):
//
//   - Lazy start: the first enqueue when no writer is running spawns it via
//     observe.Registry.Spawn (name "db-writer", owner "memory").
//   - Idle exit: the loop returns as soon as the queue is momentarily empty.
//     The exit decision and every enqueue hold the same mutex, so a command
//     is either (a) already in the queue when the writer does its final empty
//     check — it drains it, or (b) appended after — the enqueuer sees
//     running=false and spawns a fresh writer. No lost writes, no timers, no
//     wait-while-holding-lock deadlocks.

// flushTimeout bounds how long Close waits for the writer to drain (C11
// disposal budgets live at 3s; the DB flush stays well inside it).
const flushTimeout = 2 * time.Second

type writeCommand struct {
	ctx  context.Context
	fn   func(ctx context.Context, tx *sql.Tx) error
	done chan error
}

type writeQueue struct {
	db     *sql.DB // writer pool (MaxOpenConns=1)
	reg    *observe.Registry
	logger *slog.Logger

	mu      sync.Mutex
	items   []writeCommand
	running bool
	gen     uint64 // bumped on every spawn; guards the exit backstop
	closed  bool
}

func newWriteQueue(db *sql.DB, reg *observe.Registry, logger *slog.Logger) *writeQueue {
	return &writeQueue{db: db, reg: reg, logger: logger}
}

// enqueue submits one write transaction and blocks until it committed,
// rolled back, or ctx expired. The command runs even if ctx expires while
// waiting; the caller just stops watching (SQLite transactions cannot be
// cancelled mid-flight from the outside without a per-command ctx, which the
// command's own statements do honor).
func (q *writeQueue) enqueue(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	cmd := writeCommand{ctx: ctx, fn: fn, done: make(chan error, 1)}
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return ErrStoreClosed
	}
	q.items = append(q.items, cmd)
	if !q.running {
		q.running = true
		q.gen++
		g := q.gen
		q.reg.Spawn("db-writer", "memory", nil, func(ctx context.Context) {
			q.loop(g)
		})
	}
	q.mu.Unlock()

	select {
	case err := <-cmd.done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("memory: write abandoned by caller: %w", ctx.Err())
	}
}

// loop is the db-writer body (generation gen): drain the queue, exit when
// empty.
//
// Panic discipline (D38b, adversarial MAJOR-1): every command runs under its
// OWN recover (exec) — a panicking write fn becomes an internal-class error
// for its caller, its transaction is rolled back (releasing the single
// writer connection), and the loop continues. Without the per-command
// boundary one panic would (a) kill the loop with running=true, wedging
// every future Store.write, and (b) leak the checked-out Tx, deadlocking
// BeginTx on the MaxOpenConns=1 pool.
//
// The deferred reset below is the backstop for any abnormal loop death: it
// re-opens the queue so the next enqueue spawns a fresh writer — but only if
// no newer writer took over in the meantime (gen guard).
func (q *writeQueue) loop(gen uint64) {
	defer func() {
		q.mu.Lock()
		if q.gen == gen {
			q.running = false
		}
		q.mu.Unlock()
	}()
	for {
		q.mu.Lock()
		if len(q.items) == 0 {
			q.running = false
			q.mu.Unlock()
			return
		}
		cmd := q.items[0]
		q.items = q.items[1:]
		q.mu.Unlock()

		cmd.done <- q.exec(cmd)
	}
}

// exec runs one command as a single transaction on the writer connection,
// with a per-command recover boundary (the registry's recover stays as the
// last-resort backstop; this boundary keeps the loop alive so the queue can
// never wedge on a panicking command).
func (q *writeQueue) exec(cmd writeCommand) (err error) {
	tx, err := q.db.BeginTx(cmd.ctx, nil)
	if err != nil {
		return fmt.Errorf("memory: begin write tx: %w", err)
	}
	// Release the sole writer connection on EVERY exit path: after Commit a
	// Rollback is a no-op (ErrTxDone); after a panic it is what keeps the
	// queue alive (a leaked Tx would deadlock the next BeginTx forever).
	defer tx.Rollback()
	defer func() {
		if rec := recover(); rec != nil {
			stack := string(debug.Stack())
			err = observe.New(observe.ClassInternal,
				fmt.Sprintf("write command panicked: %v", rec))
			q.logger.Error("write command panic recovered (loop continues)",
				"component", "memory",
				"error_class", string(observe.ClassInternal),
				"recovered", fmt.Sprint(rec),
				"stack", stack)
		}
	}()
	if err := cmd.fn(cmd.ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("memory: commit write tx: %w", err)
	}
	return nil
}

// close marks the queue closed and waits (bounded) for the running writer to
// drain and exit. Never runs commands: everything already queued completes.
func (q *writeQueue) close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	q.mu.Unlock()

	tm := observe.NewTimeout(flushTimeout)
	for {
		q.mu.Lock()
		idle := !q.running && len(q.items) == 0
		q.mu.Unlock()
		if idle {
			return
		}
		if tm.Expired() {
			q.logger.Error("db-writer flush timeout: pending writes may not be committed",
				"component", "memory", "pending", len(q.items))
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// idle reports whether no writer is running and nothing is queued (tests).
func (q *writeQueue) idle() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return !q.running && len(q.items) == 0
}
