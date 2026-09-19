package llm

import (
	"context"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// text_chain fallback (SPEC-05 sec 3.3a / sec 3.4): an ordered walk over the
// provider catalog. 401/403 (auth), exhausted quota (budget), 5xx after the
// retry ladder (provider) and persistent network failures all step to the
// next element; the chain is exhausted -> Error(provider), task ctx kept
// (D40#4 restartability is the agent core's job, ticket 10).
//
// The failover is NEVER silent: ChainRunner.OnFailover carries every switch
// with its cause class so the ball/panel layer can surface it (SPEC-05 sec
// 3.3a: "quota-exhaustion switches must produce a user-visible event"). The
// C6 event set is contract-exact, so failover is signaled by callback, not
// by a synthetic stream event.

// FailoverEvent describes one text_chain step.
type FailoverEvent struct {
	Index   int            // index of the element being LEFT
	From    string         // "provider/model" being left
	To      string         // "provider/model" being tried next ("" when chain exhausted)
	Cause   *observe.Error // classified reason for the switch
	Attempt int            // 1-based element position
	Total   int            // chain length
}

// ChainRunner streams a request across an ordered provider list.
type ChainRunner struct {
	// Elements are the chain members in priority order (already built
	// providers; each carries its own retry wrapper).
	Elements []LlmProvider
	// Names[i] renders Elements[i] as "provider/model" for events and logs.
	Names []string
	// OnFailover receives every switch (may be nil in tests). Called
	// synchronously before the next element starts.
	OnFailover func(FailoverEvent)
}

// Stream walks the chain.
func (r *ChainRunner) Stream(ctx context.Context, req *Request, emit func(StreamEvent) error) error {
	var lastErr *observe.Error
	for i, p := range r.Elements {
		delivered := false
		err := p.Stream(ctx, req, func(ev StreamEvent) error {
			if ev.IsData() {
				delivered = true
			}
			return emit(ev)
		})

		if err == nil {
			return nil
		}
		oe := observeErr(err)
		if oe == nil {
			oe = observe.Wrap(observe.ClassInternal, err, "chain element returned unclassified error")
		}

		// Consumer stopped listening or the user cancelled: stop everything.
		if ctx.Err() != nil || oe.Class == observe.ClassCancelled {
			return err
		}

		lastErr = oe
		if i == len(r.Elements)-1 {
			break // nothing left to try
		}

		// Payload already reached the consumer from THIS element: a next
		// element would append a fresh response to a half-delivered one.
		// The partial turn stays (marked incomplete via its terminal Error /
		// Stop{error} events) and no failover happens.
		if delivered {
			return err
		}

		// Failover decision: every provider-side terminal failure switches
		// (auth / budget / provider-exhausted / network / rate-limit after
		// the element's own retry ladder). Cancelled already returned above.
		if r.OnFailover != nil {
			r.OnFailover(FailoverEvent{
				Index:   i,
				From:    r.Names[i],
				To:      r.Names[i+1],
				Cause:   oe,
				Attempt: i + 1,
				Total:   len(r.Elements),
			})
		}
	}

	final := lastErr
	if final == nil {
		final = observe.New(observe.ClassProvider, "llm: text_chain exhausted without a classified failure")
	}
	return final
}

// NamesFor renders chain element names (helper for constructors).
func NamesFor(providers []string, models []string) []string {
	out := make([]string, 0, len(providers))
	for i := range providers {
		if i < len(models) {
			out = append(out, fmt.Sprintf("%s/%s", providers[i], models[i]))
		} else {
			out = append(out, providers[i])
		}
	}
	return out
}
