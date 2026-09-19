package llm

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Provider catalog consumption (SPEC-05 sec 3.3a; the config side is ticket
// 05's catalog schema v2). The Resolver turns a config snapshot into runnable
// endpoints: roles -> provider/model, text_chain -> ordered elements, preset
// defaults already merged by the config loader, api_key_ref resolved through
// a KeyResolver (plaintext keys exist only inside Endpoint, never in logs).

// RoleName is one of the four fixed agent roles.
type RoleName string

const (
	RoleChat          RoleName = "chat"
	RoleMemoryExtract RoleName = "memory_extract"
	RoleHandoff       RoleName = "handoff"
	RoleSummarize     RoleName = "summarize"
)

// KeyResolver resolves config api_key_ref values ("dpapi:<id>" | "env:NAME").
// *secret.Store satisfies it directly (ticket 06).
type KeyResolver interface {
	Resolve(ref string) (string, error)
}

// Role carries the call-time sampling settings resolved with a role.
type RoleSettings struct {
	Temperature       float64
	ThinkingIntensity string
	MaxOutputTokens   int
}

// Endpoint is one fully resolved provider/model target.
type Endpoint struct {
	Provider string // catalog name (diagnostics)
	Model    string
	Protocol string
	BaseURL  string
	APIKey   string // resolved plaintext; NEVER log
	Compat   config.Compat
	RPM      int
	TPM      int

	// Model catalog fields the seam consumes.
	Capabilities  config.Capabilities
	ContextWindow int
	MaxOutput     int
}

// EndpointOptions converts an Endpoint into the factory input.
func (e Endpoint) EndpointOptions() EndpointOptions {
	var o EndpointOptions
	o.Provider = e.Provider
	o.Model = e.Model
	o.Protocol = e.Protocol
	o.BaseURL = e.BaseURL
	o.APIKey = e.APIKey
	o.ContextWindow = e.ContextWindow
	o.Compat.Loose = e.Compat.Loose
	o.Compat.AllowMissingUsage = e.Compat.AllowMissingUsage
	o.Compat.ExtraHeaders = e.Compat.ExtraHeaders
	return o
}

// String renders the endpoint WITHOUT the key (log-safe).
func (e Endpoint) String() string {
	return fmt.Sprintf("%s/%s (%s %s)", e.Provider, e.Model, e.Protocol, e.BaseURL)
}

// Resolver consumes a validated config snapshot (post-preset, as returned by
// the ticket 05 loader) plus a key resolver.
type Resolver struct {
	Providers map[string]config.Provider
	TextChain []string
	Roles     config.Roles
	Keys      KeyResolver // may be nil: only keyless endpoints resolve
}

// NewResolver builds a Resolver from a *config.Config snapshot.
func NewResolver(cfg *config.Config, keys KeyResolver) *Resolver {
	if cfg == nil {
		return &Resolver{}
	}
	return &Resolver{
		Providers: cfg.LLM.Providers,
		TextChain: cfg.LLM.TextChain,
		Roles:     cfg.LLM.Roles,
		Keys:      keys,
	}
}

// splitChainElement splits "provider/model" at the FIRST '/' (openrouter
// model ids keep their slashes in the model part).
func splitChainElement(el string) (provider, model string, ok bool) {
	i := strings.Index(el, "/")
	if i <= 0 || i == len(el)-1 {
		return "", "", false
	}
	return el[:i], el[i+1:], true
}

// resolveEndpoint materializes one provider/model pair.
func (r *Resolver) resolveEndpoint(provider, model string) (Endpoint, error) {
	p, ok := r.Providers[provider]
	if !ok {
		return Endpoint{}, observe.New(observe.ClassConfig,
			fmt.Sprintf("llm: unknown provider %q in catalog", provider))
	}
	spec, ok := p.Models[model]
	if !ok {
		return Endpoint{}, observe.New(observe.ClassConfig,
			fmt.Sprintf("llm: unknown model %q of provider %q", model, provider))
	}
	ep := Endpoint{
		Provider:      provider,
		Model:         model,
		Protocol:      p.Protocol,
		BaseURL:       p.BaseURL,
		Compat:        p.Compat,
		RPM:           p.RPM,
		TPM:           p.TPM,
		Capabilities:  spec.Capabilities,
		ContextWindow: spec.ContextWindow,
		MaxOutput:     spec.MaxOutputTokens,
	}
	if ref := p.APIKeyRef; ref != "" {
		if r.Keys == nil {
			return Endpoint{}, observe.New(observe.ClassAuth,
				fmt.Sprintf("llm: provider %q has api_key_ref %q but no key resolver is wired (Unconfigured)", provider, secretRefShape(ref)))
		}
		key, err := r.Keys.Resolve(ref)
		if err != nil || key == "" {
			return Endpoint{}, observe.New(observe.ClassAuth,
				fmt.Sprintf("llm: resolving api_key_ref for provider %q failed (Unconfigured)", provider))
		}
		ep.APIKey = key
	}
	return ep, nil
}

// secretRefShape renders a ref without its value for logs ("env:NAME" keeps
// the NAME - names are not secrets - and "dpapi:<id>" keeps the blob id,
// which is likewise not the secret itself).
func secretRefShape(ref string) string { return ref }

// ResolveChain resolves the ordered text_chain into runnable endpoints.
func (r *Resolver) ResolveChain() ([]Endpoint, error) {
	eps := make([]Endpoint, 0, len(r.TextChain))
	for _, el := range r.TextChain {
		p, m, ok := splitChainElement(el)
		if !ok {
			return nil, observe.New(observe.ClassConfig,
				fmt.Sprintf("llm: text_chain element %q must be \"provider/model\"", el))
		}
		ep, err := r.resolveEndpoint(p, m)
		if err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}
	return eps, nil
}

// ResolveRole resolves one of the four fixed roles. Unset roles (both
// provider and model empty in config) fall back to the FIRST text_chain
// element with the role's temperature; when no chain exists either, a
// ClassConfig error names the role.
func (r *Resolver) ResolveRole(name RoleName) (Endpoint, RoleSettings, error) {
	var role config.Role
	switch name {
	case RoleChat:
		role = r.Roles.Chat
	case RoleMemoryExtract:
		role = r.Roles.MemoryExtract
	case RoleHandoff:
		role = r.Roles.Handoff
	case RoleSummarize:
		role = r.Roles.Summarize
	default:
		return Endpoint{}, RoleSettings{}, observe.New(observe.ClassConfig,
			fmt.Sprintf("llm: unknown role %q", name))
	}

	settings := RoleSettings{
		Temperature:       role.Temperature,
		ThinkingIntensity: role.ThinkingIntensity,
		MaxOutputTokens:   role.MaxOutputTokens,
	}

	if role.Provider != "" && role.Model != "" {
		ep, err := r.resolveEndpoint(role.Provider, role.Model)
		return ep, settings, err
	}

	// Unset role: first chain element carries the model.
	for _, el := range r.TextChain {
		p, m, ok := splitChainElement(el)
		if !ok {
			continue
		}
		ep, err := r.resolveEndpoint(p, m)
		if err != nil {
			return Endpoint{}, RoleSettings{}, err
		}
		return ep, settings, nil
	}
	return Endpoint{}, RoleSettings{}, observe.New(observe.ClassConfig,
		fmt.Sprintf("llm: role %q is unset and text_chain is empty (Unconfigured)", name))
}

// BuildChain resolves the text_chain and constructs one provider per element
// (protocol -> registered adapter), wrapped with its local rate limiter and
// the retry ladder. Elements whose protocol has no adapter yet yield a
// ClassConfig error naming ticket 11 (see NewProvider).
func (r *Resolver) BuildChain(opts ChainBuildOptions) ([]LlmProvider, []string, error) {
	eps, err := r.ResolveChain()
	if err != nil {
		return nil, nil, err
	}
	out := make([]LlmProvider, 0, len(eps))
	names := make([]string, 0, len(eps))
	for _, ep := range eps {
		p, err := BuildEndpointProvider(ep, opts)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, p)
		names = append(names, ep.Provider+"/"+ep.Model)
	}
	return out, names, nil
}

// ChainBuildOptions are the process-wide knobs for endpoint construction.
type ChainBuildOptions struct {
	// ProxyMode / ProxyURL come from the [net] section.
	ProxyMode string
	ProxyURL  string
	// HTTPClient overrides the transport per endpoint (tests); nil = default
	// client honoring the proxy options.
	HTTPClient func(ep Endpoint) *http.Client
	// RetryBase overrides the backoff base (tests); 0 = default 1s.
	RetryBase time.Duration
}

// BuildEndpointProvider constructs the provider stack for one endpoint:
// adapter -> retry ladder -> rate limiter (outermost).
func BuildEndpointProvider(ep Endpoint, opts ChainBuildOptions) (LlmProvider, error) {
	epo := ep.EndpointOptions()
	if opts.HTTPClient != nil {
		if c := opts.HTTPClient(ep); c != nil {
			epo.HTTPClient = c
		}
	} else if opts.ProxyMode != "" || opts.ProxyURL != "" {
		epo.HTTPClient = NewDefaultHTTPClient(opts.ProxyMode, opts.ProxyURL)
	}

	inner, err := NewProvider(epo)
	if err != nil {
		return nil, err
	}
	retry := NewRetrying(inner, RetryOptions{Base: opts.RetryBase})
	return NewLimiterProvider(retry, NewBucketLimiter(RateLimits{RPM: ep.RPM, TPM: ep.TPM})), nil
}

// DiscoveredModels lists catalog model ids (sorted) that exist in the
// resolver's registry for one provider - the import diff helper the GUI and
// the discovery command use.
func (r *Resolver) DiscoveredModels(provider string) []string {
	p, ok := r.Providers[provider]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(p.Models))
	for id := range p.Models {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
