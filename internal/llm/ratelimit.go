package llm

import (
	"context"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Local rate self-limiting (SPEC-05 sec 3.3a): per-provider rpm/tpm token
// buckets let Wisp "behave before the provider's 429 does". The buckets are
// monotonic-clock based (D42#9), goroutine-safe, and never enforce a quota -
// that is C23/CostMeter territory (ticket 44). They only pace requests.

// RateLimiter is the pacing primitive. Acquire blocks until tokens are
// available or ctx is done (returns ctx.Err()).
type RateLimiter interface {
	Acquire(ctx context.Context, tokens int) error
}

// RateLimits carries the parsed catalog limits (0 = unlimited).
type RateLimits struct {
	RPM int
	TPM int
}

// Active reports whether any limit is configured.
func (r RateLimits) Active() bool { return r.RPM > 0 || r.TPM > 0 }

// tokenBucket is a continuous-refill bucket (rate tokens per minute,
// capacity = rate, burst-friendly). Refill is computed lazily from the
// monotonic clock on every call; no background ticker exists.
type tokenBucket struct {
	mu       sync.Mutex
	capacity float64   // tokens the bucket holds when full (= per-minute rate)
	tokens   float64   // current level
	refillPM float64   // tokens per minute
	last     time.Time // monotonic reading of the last update
}

func newTokenBucket(perMinute int, now time.Time) *tokenBucket {
	if perMinute <= 0 {
		return nil
	}
	return &tokenBucket{
		capacity: float64(perMinute),
		tokens:   float64(perMinute), // start full: a cold client is not throttled
		refillPM: float64(perMinute),
		last:     now,
	}
}

// take removes up to want tokens, waiting until they accrue.
func (b *tokenBucket) take(ctx context.Context, want int, sleep func(context.Context, time.Duration) error) error {
	if b == nil {
		return nil
	}
	wantF := float64(want)
	if wantF > b.capacity {
		wantF = b.capacity // a single request can never exceed one full bucket
	}
	for {
		b.mu.Lock()
		now := time.Now()
		b.refill(now)
		if b.tokens >= wantF {
			b.tokens -= wantF
			b.mu.Unlock()
			return nil
		}
		need := wantF - b.tokens
		// Time until the deficit refills at refillPM tokens/minute.
		wait := time.Duration(need / b.refillPM * float64(time.Minute))
		b.mu.Unlock()

		if wait > 250*time.Millisecond {
			wait = 250 * time.Millisecond // re-check periodically (ctx cancel)
		}
		if err := sleep(ctx, wait); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// adjust returns delta tokens to the bucket (delta > 0) or charges them
// (delta < 0). Used to reconcile a pre-flight reservation with real usage.
func (b *tokenBucket) adjust(delta float64) {
	if b == nil || delta == 0 {
		return
	}
	b.mu.Lock()
	now := time.Now()
	b.refill(now)
	b.tokens += delta
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	if b.tokens < 0 {
		b.tokens = 0
	}
	b.mu.Unlock()
}

func (b *tokenBucket) refill(now time.Time) {
	elapsed := now.Sub(b.last)
	if elapsed <= 0 {
		return
	}
	b.last = now
	b.tokens += elapsed.Minutes() * b.refillPM
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}

// BucketLimiter pairs the rpm and tpm buckets of one provider.
type BucketLimiter struct {
	rpm, tpm *tokenBucket
	sleep    func(context.Context, time.Duration) error
}

// NewBucketLimiter builds the limiter for one provider (nil-safe: limits of
// 0 produce a limiter that never blocks).
func NewBucketLimiter(limits RateLimits) *BucketLimiter {
	now := time.Now()
	return &BucketLimiter{
		rpm:   newTokenBucket(limits.RPM, now),
		tpm:   newTokenBucket(limits.TPM, now),
		sleep: sleepCtx,
	}
}

// Acquire reserves one request (rpm) and est tokens (tpm).
func (l *BucketLimiter) Acquire(ctx context.Context, estTokens int) error {
	if l == nil {
		return nil
	}
	if err := l.rpm.take(ctx, 1, l.sleep); err != nil {
		return err
	}
	return l.tpm.take(ctx, estTokens, l.sleep)
}

// Reconcile adjusts the tpm bucket by raw delta (reserved surplus returns,
// under-estimates charge). Exposed for tests and callers that already know
// the delta.
func (l *BucketLimiter) Reconcile(delta int) {
	if l == nil {
		return
	}
	l.tpm.adjust(float64(delta))
}

// LimiterProvider decorates a provider with the local pacing (and the
// post-response usage reconcile). It sits between the chain element and the
// retry wrapper: retries already consumed the 429 lesson, the limiter only
// gates NEW attempts. To keep that ordering simple it wraps the OUTERMOST
// provider (see Resolver.BuildChain).
type LimiterProvider struct {
	inner    LlmProvider
	limiter  *BucketLimiter
	reserved int
}

// NewLimiterProvider wraps inner with limiter (nil limiter = pass-through).
func NewLimiterProvider(inner LlmProvider, limiter *BucketLimiter) *LimiterProvider {
	return &LimiterProvider{inner: inner, limiter: limiter}
}

// Info passes through.
func (p *LimiterProvider) Info() ProviderInfo { return p.inner.Info() }

// Stream gates the request, then streams, then reconciles the tpm bucket
// against the real usage carried by C6 Usage events: the pre-flight estimate
// is refunded or topped up exactly once (usage is monotonic, so each event
// only charges the increment over what was already accounted).
func (p *LimiterProvider) Stream(ctx context.Context, req *Request, emit func(StreamEvent) error) error {
	if p.limiter == nil {
		return p.inner.Stream(ctx, req, emit)
	}
	reserved := req.EstimateTokens()
	if err := p.limiter.Acquire(ctx, reserved); err != nil {
		return observe.Wrap(observe.ClassCancelled, err, "cancelled while self-throttling")
	}
	accounted := reserved
	firstUsage := true
	return p.inner.Stream(ctx, req, func(ev StreamEvent) error {
		if ev.Type == EvUsage {
			actual := ev.Usage.InputTokens + ev.Usage.OutputTokens
			if actual > 0 && (firstUsage || actual > accounted) {
				// First event reconciles the pre-flight reservation with the
				// real usage (refund or charge); later events charge only the
				// increment (usage is monotonic across a stream).
				p.limiter.Reconcile(accounted - actual)
				accounted = actual
				firstUsage = false
			}
		}
		return emit(ev)
	})
}
