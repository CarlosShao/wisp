package tools

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/risk"
)

// C4 ToolProvider (SPEC-07 §2): four slots - Builtin, Manifest (Tier-1),
// Goja (Tier-2), MCP. Only Builtin has an implementation here.
type Kind string

// The C4 slots (SPEC-07 §2).
const (
	KindBuiltin  Kind = "builtin"
	KindManifest Kind = "manifest"
	KindGoja     Kind = "goja"
	KindMCP      Kind = "mcp"
)

// ErrSlotNotLanded is what the unimplemented C4 slots answer. It is a distinct,
// checkable signal rather than an empty list, so a composition root that tries
// to wire Tier-1/Tier-2/MCP discovers the gap at construction time instead of
// shipping a provider that silently offers zero tools.
var ErrSlotNotLanded = errors.New(
	"tools: provider slot declared but not implemented (manifest=ticket 50, " +
		"goja=ticket 51, mcp=REJECTED per D13/16.9#7)")

// SlotErr records WHY a slot is still empty, so the reason travels with the
// error rather than living only in a comment.
type SlotErr struct {
	Kind   Kind
	Ticket string // who lands it, or "REJECTED"
}

func (e SlotErr) Error() string {
	if e.Ticket == "" {
		return fmt.Sprintf("%v: slot %s", ErrSlotNotLanded, e.Kind)
	}
	return fmt.Sprintf("%v: slot %s (%s)", ErrSlotNotLanded, e.Kind, e.Ticket)
}

// Unwrap lets errors.Is(err, ErrSlotNotLanded) work through the wrap.
func (e SlotErr) Unwrap() error { return ErrSlotNotLanded }

// providerSlots is the declared C4 surface: every slot the spec names, whether
// or not anything lives behind it.
var providerSlots = []SlotErr{
	{Kind: KindManifest, Ticket: "ticket 50 lands Tier-1 manifests"},
	{Kind: KindGoja, Ticket: "ticket 51 lands the Tier-2 goja runtime"},
	{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"},
}

// Provider is one C4 slot's view of its tools. Builtin implements it here;
// tickets 50/51 add the Manifest and Goja implementations and register them
// against the same Registry, which is why nothing below either slot has to
// know about plugin plumbing.
type Provider interface {
	Kind() Kind
	// List returns the slot's entries in stable order.
	List() []Entry
	// Lookup resolves one tool name inside this slot.
	Lookup(name string) (Entry, bool)
}

// Entry is one registered tool plus its host-side declaration.
type Entry struct {
	Tool Tool
	Decl Decl
}

// Registry is the aggregate tool directory the bridge consults. Safe for
// concurrent use: registration happens at startup while execution is
// per-task-concurrent.
type Registry struct {
	mu     sync.RWMutex
	byName map[string]Entry
	order  []string
	slots  map[Kind]Provider
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{byName: map[string]Entry{}, slots: map[Kind]Provider{}}
}

// Register validates and adds one entry. Validation failures are returned as
// errors and leave the registry untouched - a tool that cannot state its
// capability face or its name shape is a wiring bug, and refusing it at
// registration is what stops the bridge from ever seeing such a call.
func (r *Registry) Register(e Entry) error {
	if err := r.validate(e); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	name := e.Tool.Name()
	if _, dup := r.byName[name]; dup {
		return fmt.Errorf("tools: register %q: duplicate tool name (C1 requires global uniqueness)", name)
	}
	e.Decl.Capabilities = dedupe(e.Decl.Capabilities)
	e.Decl.Needs = dedupe(e.Decl.Needs)
	if e.Decl.Provider == "" {
		e.Decl.Provider = KindBuiltin
	}
	r.addLocked(e)
	return nil
}

// validate runs the registration-time contract checks on one entry without
// touching the registry.
func (r *Registry) validate(e Entry) error {
	if e.Tool == nil {
		return errors.New("tools: register: nil Tool")
	}
	name := e.Tool.Name()
	if err := validToolName(name); err != nil {
		return fmt.Errorf("tools: register %q: %w", name, err)
	}
	for _, p := range e.Decl.PathParams {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("tools: register %q: empty PathParams entry", name)
		}
	}
	if err := validCaps(e.Decl.Capabilities); err != nil {
		return fmt.Errorf("tools: register %q: declared capabilities: %w", name, err)
	}
	if err := validCaps(e.Decl.Needs); err != nil {
		return fmt.Errorf("tools: register %q: needed capabilities: %w", name, err)
	}
	// The registration-time half of the C3 rule. It is not the enforcement
	// point (the bridge re-checks per call, because the declared set is a
	// host-side view that can shrink between registration and the call), but
	// failing here turns a mis-declared builtin into a startup error instead
	// of a mystery rejection at the first user request.
	if miss := newCapSet(e.Decl.Capabilities...).missing(e.Decl.Needs); len(miss) > 0 {
		return fmt.Errorf("tools: register %q: needs undeclared capability %s", name, joinCaps(miss))
	}
	if e.Decl.Declared > risk.L2 {
		return fmt.Errorf("tools: register %q: declared risk %v above L2", name, e.Decl.Declared)
	}
	if len(e.Tool.Parameters()) == 0 {
		return fmt.Errorf("tools: register %q: no parameter schema", name)
	}
	return nil
}

// addLocked inserts one already-validated entry; callers hold the write lock.
func (r *Registry) addLocked(e Entry) {
	name := e.Tool.Name()
	r.byName[name] = e
	r.order = append(r.order, name)
	sort.Strings(r.order)
}

// RegisterProvider attaches a C4 slot's provider. The three unlanded slots
// answer ErrSlotNotLanded rather than being silently refused, so the gap is
// discoverable from the composed registry.
func (r *Registry) RegisterProvider(p Provider) error {
	if p == nil {
		return errors.New("tools: register provider: nil")
	}
	for _, s := range providerSlots {
		if s.Kind == p.Kind() {
			return s // SlotErr: the slot exists in the contract and nothing lands it yet
		}
	}
	if p.Kind() != KindBuiltin {
		return fmt.Errorf("tools: register provider: unknown slot %q", p.Kind())
	}
	// Validate every entry first (Register would refuse a bad one, but it also
	// takes the write lock, so the checks are split out to avoid re-entering).
	var ok []Entry
	for _, e := range p.List() {
		if err := r.validate(e); err != nil {
			return err
		}
		ok = append(ok, e)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.slots[p.Kind()] = p
	for _, e := range ok {
		if _, dup := r.byName[e.Tool.Name()]; !dup {
			if e.Decl.Provider == "" {
				e.Decl.Provider = KindBuiltin
			}
			e.Decl.Capabilities = dedupe(e.Decl.Capabilities)
			e.Decl.Needs = dedupe(e.Decl.Needs)
			r.addLocked(e)
		}
	}
	return nil
}

// UnlandedSlots returns one error per C4 slot with no implementation.
func UnlandedSlots() []error {
	out := make([]error, 0, len(providerSlots))
	for _, s := range providerSlots {
		out = append(out, s)
	}
	return out
}

// Lookup resolves one tool name.
func (r *Registry) Lookup(name string) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.byName[name]
	return e, ok
}

// List returns every entry in name order.
func (r *Registry) List() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.listLocked()
}

func (r *Registry) listLocked() []Entry {
	out := make([]Entry, 0, len(r.order))
	for _, n := range r.order {
		if e, ok := r.byName[n]; ok {
			out = append(out, e)
		}
	}
	return out
}

// validCaps checks every token is in the frozen C3 set.
func validCaps(caps []Capability) error {
	for _, c := range caps {
		if !c.Valid() {
			return fmt.Errorf("unknown capability %q (C3 is a frozen set of %d tokens)",
				c, len(AllCapabilities))
		}
	}
	return nil
}

// validToolName enforces the C1 'namespace.action' shape: exactly one dot,
// non-empty lowercase segments.
func validToolName(name string) error {
	if name == "" {
		return errors.New("empty tool name")
	}
	i := strings.IndexByte(name, '.')
	if i <= 0 || i == len(name)-1 {
		return errors.New("name must be 'namespace.action'")
	}
	if strings.Contains(name[i+1:], ".") {
		return errors.New("name has more than one '.')")
	}
	for _, seg := range []string{name[:i], name[i+1:]} {
		for _, r := range seg {
			ok := r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' || r == '-'
			if !ok {
				return fmt.Errorf("name segment %q has a disallowed character", seg)
			}
		}
	}
	return nil
}

// builtinProvider is the C4 Builtin slot: a static list assembled by the
// composition root.
type builtinProvider struct {
	entries []Entry
}

// NewBuiltinProvider builds the Builtin slot over the given entries.
func NewBuiltinProvider(entries ...Entry) Provider {
	return &builtinProvider{entries: entries}
}

func (p *builtinProvider) Kind() Kind { return KindBuiltin }

// List implements Provider.
func (p *builtinProvider) List() []Entry {
	out := make([]Entry, len(p.entries))
	copy(out, p.entries)
	return out
}

// Lookup implements Provider.
func (p *builtinProvider) Lookup(name string) (Entry, bool) {
	for _, e := range p.entries {
		if e.Tool.Name() == name {
			return e, true
		}
	}
	return Entry{}, false
}
