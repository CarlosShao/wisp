package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// C5 LlmProvider - the three-protocol seam (D8 / SPEC-05 sec 3.1). Agent core
// code depends on this interface only; protocol differences (wire shapes, SSE
// dialects, error payloads) stop at the adapters registered in this package's
// registry.
//
// Implementations MUST honor the C6 stream contract documented in events.go
// (exactly one terminal pair, then Done) and MUST classify every failure into
// a D37 observe.ErrorClass at the seam (SPEC-05 sec 3: "errors classified at
// the seam").

// CacheMode is the prompt-caching capability of a provider (SPEC-05 sec 4.1:
// "C5 must expose cache-breakpoint position as a provider capability").
// Anthropic carries explicit cache_control markers; OpenAI-compatible
// providers cache implicitly by byte-stable prefixes; some providers do not
// cache at all.
type CacheMode string

const (
	CacheNone     CacheMode = "none"
	CacheImplicit CacheMode = "implicit" // byte-stable prefix caching only
	CacheExplicit CacheMode = "explicit" // Request.CacheBreakpoints honored
)

// CacheSupport is the caching capability bit-block of a provider/model pair.
type CacheSupport struct {
	Mode CacheMode
	// Breakpoints reports whether Request.CacheBreakpoints has an effect.
	Breakpoints bool
}

// ProviderInfo describes the provider/model a stream request will hit.
type ProviderInfo struct {
	Protocol string // config protocol enum (openai-chat | openai-responses | anthropic)
	Provider string // catalog provider name
	Model    string // model id
	Cache    CacheSupport
	// MaxContextWindow of the model (0 = unknown); D15 thresholds scale by it.
	MaxContextWindow int
}

// Request is one normalized LLM call (C7-based). It is consumed exactly once
// by Stream and must not be mutated by the provider.
type Request struct {
	Model string

	// System is the cache-stable system prompt (SPEC-05 sec 4.1: the prefix
	// blocks 1/5/6/2 are byte-stable; the time/focus block belongs at the end
	// of Messages instead).
	System []Content

	// Messages are the conversation turns (user / assistant / tool).
	Messages []Message

	// Tools the model may call; nil/empty declares none.
	Tools []ToolDef

	// ToolChoice is nil for the provider default (auto).
	ToolChoice *ToolChoice

	Temperature       *float64
	MaxOutputTokens   int
	ThinkingIntensity string // off|low|medium|high (config enum; "" = default)
	StopSequences     []string

	// CacheBreakpoints mark positions AFTER which the prefix is stable and
	// cacheable (index into Messages; -1 = after System). Only honored by
	// CacheExplicit providers (see CacheSupport).
	CacheBreakpoints []int

	// ResolveBytes resolves ImagePart.BytesRef into raw bytes at wire-encode
	// time. Required only for requests carrying image parts; the resolved
	// bytes are encoded into the wire request and must never be logged.
	ResolveBytes func(ref string) (mimeType string, data []byte, err error)
}

// Validate performs the cheap structural checks the seam owns (adapters do
// their own protocol-specific checks).
func (r *Request) Validate() error {
	if r == nil {
		return observe.New(observe.ClassInternal, "llm: nil request")
	}
	if r.Model == "" {
		return observe.New(observe.ClassConfig, "llm: request.model is empty")
	}
	for _, m := range r.Messages {
		if len(m.Content) == 0 {
			return observe.New(observe.ClassInternal,
				fmt.Sprintf("llm: message with role %q has empty content", m.Role))
		}
	}
	return nil
}

// EstimateTokens is the pre-flight token estimate used by local TPM
// self-throttling (internal/llm ratelimit). It is a coarse heuristic
// (utf-8 bytes / 4) over text and tool schemas; exact accounting happens via
// Usage events after the response.
func (r *Request) EstimateTokens() int {
	n := 0
	count := func(s string) { n += len(s) }
	for _, p := range r.System {
		if t, ok := p.(TextPart); ok {
			count(t.Text)
		}
	}
	for _, m := range r.Messages {
		for _, p := range m.Content {
			switch c := p.(type) {
			case TextPart:
				count(c.Text)
			case ToolUsePart:
				count(string(c.Input))
			case ToolResultPart:
				for _, ip := range c.Content {
					if t, ok := ip.(TextPart); ok {
						count(t.Text)
					}
				}
			}
		}
	}
	for _, t := range r.Tools {
		count(t.Name + t.Description + string(t.Parameters))
	}
	tokens := n / 4
	if r.MaxOutputTokens > 0 {
		tokens += r.MaxOutputTokens
	}
	return tokens
}

// LlmProvider is the C5 seam. See the package and events.go docs for the
// event contract.
type LlmProvider interface {
	// Stream runs one normalized request and emits C6 events until Done.
	// emit may return an error to stop consumption; Stream then aborts and
	// returns that error verbatim.
	Stream(ctx context.Context, req *Request, emit func(StreamEvent) error) error

	// Info describes the provider/model the next Stream will hit.
	Info() ProviderInfo
}

// ProviderFactory builds a provider from resolved endpoint options.
type ProviderFactory func(ep EndpointOptions) (LlmProvider, error)

// EndpointOptions are the resolved endpoint settings a factory needs. They
// come from the config catalog (ticket 05) via the package Resolver.
type EndpointOptions struct {
	Provider string // catalog provider name (diagnostics only)
	Model    string
	Protocol string
	BaseURL  string
	APIKey   string // resolved plaintext; never logged
	Compat   struct {
		Loose             bool
		AllowMissingUsage bool
		ExtraHeaders      map[string]string
	}
	// HTTPClient overrides the transport (proxy/TLS tests); nil = default
	// client with proxy support (ProxyFunc) applied.
	HTTPClient *http.Client
	// Retry-after fallback for 429 without a header, and backoff base/ceiling
	// are owned by the retry wrapper, not the endpoint.
}

// registry maps protocol enum -> factory. Adapters self-register from their
// own package init (ticket 09: openai-chat; ticket 11: openai-responses,
// anthropic).
var (
	regMu    sync.RWMutex
	registry = map[string]ProviderFactory{}
)

// RegisterProtocol installs a factory for a config protocol enum. Registering
// the same protocol twice panics (programming error, caught by tests).
func RegisterProtocol(protocol string, f ProviderFactory) {
	regMu.Lock()
	defer regMu.Unlock()
	if _, dup := registry[protocol]; dup {
		panic(fmt.Sprintf("llm: protocol %q registered twice", protocol))
	}
	registry[protocol] = f
}

// Protocols lists the registered protocol enums, sorted.
func Protocols() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(registry))
	for p := range registry {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// NewProvider builds the provider for ep.Protocol. Unregistered protocols
// yield a ClassConfig error naming the ticket that owns the adapter
// (ticket 11 for openai-responses / anthropic) so a chain referencing a
// not-yet-implemented protocol fails with an actionable message instead of
// panicking.
func NewProvider(ep EndpointOptions) (LlmProvider, error) {
	regMu.RLock()
	f, ok := registry[ep.Protocol]
	regMu.RUnlock()
	if !ok {
		return nil, observe.New(observe.ClassConfig, fmt.Sprintf(
			"llm: no adapter registered for protocol %q (provider %q): openai-responses and anthropic adapters land with ticket 11",
			ep.Protocol, ep.Provider))
	}
	return f(ep)
}

// UnmarshalArgs is the shared helper adapters use to normalize tool argument
// JSON: empty fragments become empty objects, anything else must be valid
// JSON (agent core feeds ToolUsePart.Input straight into tool dispatch).
func UnmarshalArgs(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return json.RawMessage("{}"), nil
	}
	if !json.Valid(raw) {
		return nil, observe.New(observe.ClassProvider, "llm: tool arguments are not valid JSON")
	}
	return raw, nil
}
