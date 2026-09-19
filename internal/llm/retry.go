package llm

import (
	"context"
	"math"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Retry / backoff (SPEC-05 sec 3.4: "the retry loop lives outside the agent
// core"). RetryingProvider wraps any LlmProvider:
//
//   - exponential backoff, at most Max attempts (D37 ceiling 3);
//   - 429 honors the upstream retry-after hint verbatim (no backoff on top);
//   - 401/403 (auth), quota/budget and other non-retryable classes surface
//     immediately;
//   - a failure AFTER any data event reached the consumer is never retried
//     (a reconnect would duplicate payload); the stream ends mid-flight with
//     the failure marked instead;
//   - cancellation is never retried.
//
// The wrapper only decides IF a retry happens; the terminal Error event for
// the final failure is emitted by the inner provider and passed through.

// RetryOptions configures RetryingProvider. Zero fields take the defaults
// documented per field.
type RetryOptions struct {
	// Max is the number of RETRIES after the initial attempt (D37 ceiling 3).
	// 0 = default 3; negative disables retries entirely.
	Max int
	// Base is the initial backoff (default 1s; config llm.retry.backoff_ms).
	Base time.Duration
	// MaxDelay caps the backoff ladder (default 15s).
	MaxDelay time.Duration
	// Sleep is the wall-clock abstaining delay function (tests inject a
	// recorder). It must respect ctx cancellation. nil = time-based sleep.
	Sleep func(ctx context.Context, d time.Duration) error
	// Now is the time source for Retry-After HTTP-date parsing.
	Now func() time.Time
}

func (o *RetryOptions) fill() {
	if o.Max == 0 {
		o.Max = observe.MaxBackoffAttempts
	}
	if o.Max < 0 {
		o.Max = 0
	}
	if o.Base <= 0 {
		o.Base = time.Second
	}
	if o.MaxDelay <= 0 {
		o.MaxDelay = 15 * time.Second
	}
	if o.Sleep == nil {
		o.Sleep = sleepCtx
	}
	if o.Now == nil {
		o.Now = time.Now
	}
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// RetryingProvider is the C5 decorator adding the retry ladder.
type RetryingProvider struct {
	inner LlmProvider
	opts  RetryOptions
}

// NewRetrying wraps inner with the retry policy.
func NewRetrying(inner LlmProvider, opts RetryOptions) *RetryingProvider {
	opts.fill()
	return &RetryingProvider{inner: inner, opts: opts}
}

// Info passes through to the inner provider.
func (r *RetryingProvider) Info() ProviderInfo { return r.inner.Info() }

// Stream implements LlmProvider with the retry ladder. Terminal events of an
// attempt (Error/Stop/Done) are held back until the wrapper knows the attempt
// is final: a swallowed attempt's terminals are discarded, a final attempt's
// are flushed in order, keeping the C6 contract intact for the consumer.
func (r *RetryingProvider) Stream(ctx context.Context, req *Request, emit func(StreamEvent) error) error {
	attempt := 0
	for {
		delivered := false
		var terminal []StreamEvent // buffered Error/Stop/Done of this attempt

		err := r.inner.Stream(ctx, req, func(ev StreamEvent) error {
			if ev.IsData() {
				delivered = true
				return emit(ev)
			}
			terminal = append(terminal, ev)
			return nil
		})

		if ctx.Err() != nil && !delivered && err == nil {
			// Cancelled before anything was delivered: still forward the
			// cancelled terminal so the consumer sees a well-formed stream.
			return flushTerminals(terminal, emit)
		}

		if err == nil {
			return flushTerminals(terminal, emit)
		}

		// emit errors: the consumer stopped listening - never retry.
		if ctx.Err() != nil || classOf(err) == observe.ClassCancelled {
			return err
		}

		if !retryWanted(err, classOf(err)) || delivered {
			// Final failure: surface the buffered terminals, then the error.
			if ferr := flushTerminals(terminal, emit); ferr != nil {
				return ferr
			}
			return err
		}

		// Swallow this attempt's terminals (consumer saw no data) and retry.
		attempt++
		if attempt > r.opts.Max {
			if ferr := flushTerminals(terminal, emit); ferr != nil {
				return ferr
			}
			return err
		}
		delay := r.delayFor(err, classOf(err), attempt)
		if serr := r.opts.Sleep(ctx, delay); serr != nil {
			if ferr := flushTerminals(terminal, emit); ferr != nil {
				return ferr
			}
			return err // ctx died while backing off: report the real failure
		}
	}
}

func flushTerminals(terminal []StreamEvent, emit func(StreamEvent) error) error {
	for _, ev := range terminal {
		if err := emit(ev); err != nil {
			return err
		}
	}
	return nil
}

// retryWanted applies the D37 policy with the seam refinements:
//   - no-retry sentinels (TLS interception) never retry;
//   - rate_limit requires a retry-after hint (observe.Error.Retryable already
//     encodes that; NewHTTPError guarantees the hint or the fallback);
//   - cancelled never retries.
func retryWanted(err error, class observe.ErrorClass) bool {
	if isNoRetry(err) {
		return false
	}
	if class == observe.ClassCancelled {
		return false
	}
	if class == observe.ClassRateLimit {
		e := observeErr(err)
		if e == nil {
			return false
		}
		return e.Retryable()
	}
	return class.Retryable()
}

// delayFor computes the sleep before the next attempt: 429 with a hint waits
// exactly the hint; everything else climbs base*2^(attempt-1), capped.
func (r *RetryingProvider) delayFor(err error, class observe.ErrorClass, attempt int) time.Duration {
	if class == observe.ClassRateLimit {
		if e := observeErr(err); e != nil && e.RetryAfter > 0 {
			return e.RetryAfter
		}
	}
	d := time.Duration(float64(r.opts.Base) * math.Pow(2, float64(attempt-1)))
	if d > r.opts.MaxDelay {
		return r.opts.MaxDelay
	}
	return d
}

// classOf extracts the D37 class from an error chain (observe.ClassOf).
func classOf(err error) observe.ErrorClass {
	c, _ := observe.ClassOf(err)
	return c
}

// observeErr returns the first *observe.Error in the chain, or nil.
func observeErr(err error) *observe.Error {
	var e *observe.Error
	if errorsAs(err, &e) {
		return e
	}
	return nil
}

// errorsAs is a local errors.As (kept dependency-light for the seam tests).
func errorsAs(err error, target **observe.Error) bool {
	for err != nil {
		if e, ok := err.(*observe.Error); ok {
			*target = e
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}
