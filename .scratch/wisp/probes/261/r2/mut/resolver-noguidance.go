package llm

// Ticket 261 leg r2 - mutation self-proof for fix B. This is a copy of
// internal/llm/resolver.go as of HEAD cf95c799 with ONLY the guidance
// half-sentence removed (the message back to the pre-r2 wording). It reaches
// the build through -overlay (mut/overlay.json); the tracked tree never
// changes. Expected reading: the r2 guidance assertions in
// enabled_gate_261_r1_test.go and enabled_reach_261_test.go go RED; every
// other reading in the two files stays as it was.

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

type RoleName string

const (
	RoleChat          RoleName = "chat"
	RoleMemoryExtract RoleName = "memory_extract"
	RoleHandoff       RoleName = "handoff"
	RoleSummarize     RoleName = "summarize"
)

type KeyResolver interface {
	Resolve(ref string) (string, error)
}

type RoleSettings struct {
	Temperature       float64
	ThinkingIntensity string
	MaxOutputTokens   int
}

type Endpoint struct {
	Provider      string
	Model         string
	Protocol      string
	BaseURL       string
	APIKey        string
	Compat        config.Compat
	RPM           int
	TPM           int
	Capabilities  config.Capabilities
	ContextWindow int
	MaxOutput     int
}

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

func (e Endpoint) String() string {
	return fmt.Sprintf("%s/%s (%s %s)", e.Provider, e.Model, e.Protocol, e.BaseURL)
}

type Resolver struct {
	Providers map[string]config.Provider
	TextChain []string
	Roles     config.Roles
	Keys      KeyResolver
}

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

func splitChainElement(el string) (provider, model string, ok bool) {
	i := strings.Index(el, "/")
	if i <= 0 || i == len(el)-1 {
		return "", "", false
	}
	return el[:i], el[i+1:], true
}

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
	if !spec.Enabled {
		// MUTATION: the guidance half-sentence (fix B) is removed here; the
		// message is byte-identical to the pre-r2 wording.
		return Endpoint{}, observe.New(observe.ClassConfig,
			fmt.Sprintf("llm: model %q of provider %q is disabled (enabled=false); re-enable it or remove the entry", model, provider))
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

func secretRefShape(ref string) string { return ref }

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

type ChainBuildOptions struct {
	ProxyMode  string
	ProxyURL   string
	HTTPClient func(ep Endpoint) *http.Client
	RetryBase  time.Duration
	RetryMax   int
}

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
	retry := NewRetrying(inner, RetryOptions{Base: opts.RetryBase, Max: opts.RetryMax})
	return NewLimiterProvider(retry, NewBucketLimiter(RateLimits{RPM: ep.RPM, TPM: ep.TPM})), nil
}

func (r *Resolver) DiscoveredModels(provider string) []string {
	p, ok := r.Providers[provider]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(p.Models))
	for id, spec := range p.Models {
		if !spec.Enabled {
			continue
		}
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
