package plugin

// C24 GojaHostAPI - ticket 51, BLOCK B.
//
// C24 freezes three columns and nothing else (PLAN.md:1380):
//
//	entry (method name) · capability (C3 token) · RiskLevel
//
// The owner approved the initial set on 2026-11-05 (ledger A1541 item 1: "票51
// AC#3（C24 GojaHostAPI 六组白名单）批准"), which is the six groups in
// C24SurfaceTable: http, fs, clipboard, notify, sysinfo, state. Opening this
// set - adding an entry, dropping one, changing a capability or a level - is a
// contract change (AGENTS.md §0 "改契约＝人工批准"); TestC24SurfaceIsFrozenAnd-
// OwnerApproved is the ruler that rings first.
//
// What C24 does NOT freeze, and this file therefore does not invent as a
// contract: the SHAPE of each entry's argument and reply. Every entry takes one
// optional JSON-serializable object argument and returns whatever the host
// answers; the meaning of that object is the host implementation's business
// (the C24 table has no argument column to read). The uniform calling
// convention is what makes "undeclared entry access" a single, testable rule.
//
// The three enforcement layers, in the order a call passes them:
//
//  1. C3 at the bridge - the manifest's declared capabilities travel on
//     tools.Decl.Needs, and the bridge refuses an undeclared one BEFORE Execute
//     (tools.Bridge.checkCaps). This file does not repeat that check; it feeds
//     it (§c24 capability face below).
//  2. This surface - only entries whose capability the tool DECLARED are
//     installed into the VM. An entry outside the declaration is `undefined`
//     at runtime (accessing it throws a TypeError the plugin author can read)
//     and is logged as a contract violation.
//  3. C19 at the bridge - C24Request carries the entry's frozen RiskLevel as
//     the R1 floor, and the host judges the concrete target (paths, allowlist,
//     taint) before acting. A declaration is input only; C19 may raise it.
//
// A host implementation can also CRASH (a nil map write, a broken assumption).
// goja does not turn such a panic into a JS exception - exceptionFromValue
// answers nil for a non-Exception value and handleThrow re-panics it - so the
// value travels out of RunProgram as a raw Go panic. bind() recovers it at the
// exact seam, attributes it (C24HostPanic: entry, capability, risk floor, tool,
// plugin, crash-site stack), records the attempt on Calls() and logs one
// C24-HOSTPANIC line, then re-panics so GojaRuntime.Execute's calling-goroutine
// boundary owns the recovery. It is deliberately NOT converted into a JS
// exception: a plugin that could catch a host crash could keep hammering a host
// whose state is unknown. Ticket 51 AC#4; the e2e ruler is
// panic_isolation_51_test.go.
//
// Isolation: the surface is built per call and installed into the per-call VM
// (goja_runtime.go), so state written by one call - including a plugin that
// monkey-patches wisp.fs.read - dies with that VM.
//
// There is no memory ceiling here and none is claimed: goja has no such
// mechanism (D33/F5). The only ceilings are block A's wall clock and result cap.
//
// NOT in this file (and why): the enable gate and its L2 re-confirmation are
// tier2_gate.go (AC#5); reading a plugin's bundle and turning a manifest call
// into a GojaCall is tier2_loader.go; the host behind C24Host - the object that
// owns risk assessment, path canonicalization, taint and the tool_call row - is
// wired by the composition root (it lives in internal/tools, which this package
// cannot import: tools -> agent -> memory -> plugin).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"sync"

	"github.com/dop251/goja"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// c24RootName is the single global the surface adds. Everything the VM can
// reach of the host lives under it; there is no `fetch`, no timer, no module
// loader (TestC24AddsNoGlobalsBeyondTheSurface diffs the global list).
const c24RootName = "wisp"

// The frozen entry paths (C24 table order; a doc/reading consumer may rely on
// it, so append, never reorder).
const (
	C24HTTPRequest    = "wisp.http.request"
	C24FSRead         = "wisp.fs.read"
	C24FSWrite        = "wisp.fs.write"
	C24ClipboardRead  = "wisp.clipboard.read"
	C24ClipboardWrite = "wisp.clipboard.write"
	C24NotifySend     = "wisp.notify.send"
	C24SysinfoGet     = "wisp.sysinfo.get"
	C24StateGet       = "wisp.state.get"
	C24StateSet       = "wisp.state.set"
)

// C24Entry is one row of the frozen C24 table.
//
// Capability is a C3 token spelled as a string because internal/plugin cannot
// import internal/tools (import cycle), not because the token is free-form: the
// empty string is reserved for the capability-free entries (the plugin's own
// KV), and any other value must be in the frozen eleven-token set.
type C24Entry struct {
	Path       string
	Capability string
	Level      risk.Level
	Doc        string
}

// C24SurfaceTable is the owner-approved initial surface. RiskLevel is the R1
// floor for that entry, taken from D34's row for the equivalent builtin so the
// face has one story:
//
//   - fs.read L0 / clipboard.read L0 / sysinfo.get L0: read-only (D34 rows
//     fs.read, clipboard.read, system.get);
//   - fs.write L1 / clipboard.write L1: a reversible write, and clipboard.write
//     is explicitly an exfiltration channel (D34 "算外泄通道", F4);
//   - http.request L1: same reasoning as clipboard.write - the surface cannot
//     tell a GET from a POST before the call, so the floor is the egress-able
//     one (D34 web.open is L1 for the same "leaves the machine" reason; a bare
//     fetch's L0 is reachable only through C19's allowlist gate, which the host
//     applies per target);
//   - notify.send L0 (D34 notify);
//   - state.get/set L0 and capability-free: the plugin's own data (SPEC-07 §5).
//
// No row is Deny: Deny belongs to C19/R3 (a tier-A path, a deny-listed host),
// which fuses per call and takes the max. Writing Deny here would claim the
// entry is unconditionally forbidden, which is the opposite of why it exists.
var C24SurfaceTable = []C24Entry{
	{
		C24HTTPRequest, "net", risk.L1,
		"One outbound HTTP(S) request to a host inside the plugin's net_allowlist; returns the host's reply.",
	},
	{
		C24FSRead, "fs.read", risk.L0,
		"Read a file through the host's C26 path gate; an out-of-allowlist path is raised by C19/R2, never by the surface.",
	},
	{
		C24FSWrite, "fs.write", risk.L1,
		"Write a file through the host's C26 path gate; overwriting an existing file is raised by C19/R8.",
	},
	{
		C24ClipboardRead, "clipboard", risk.L0,
		"Read the clipboard; the result is a C25 taint source, so it cannot later leave unnoticed.",
	},
	{
		C24ClipboardWrite, "clipboard", risk.L1,
		"Write the clipboard - an exfiltration channel (D34/F4), hence the L1 floor.",
	},
	{
		C24NotifySend, "notify", risk.L0,
		"Send a system notification; a URL inside the content is taint-marked like D34's notify.",
	},
	{
		C24SysinfoGet, "sysinfo", risk.L0,
		"Read host facts (battery, network, volume, focused window title - the title is tainted).",
	},
	{
		C24StateGet, "", risk.L0,
		"Read this plugin's own key/value state. Capability-free: it is the plugin's data, not the user's.",
	},
	{
		C24StateSet, "", risk.L0,
		"Write this plugin's own key/value state. Capability-free, same reason as state.get.",
	},
}

// C24EntryFor looks one frozen entry up by its JS path.
func C24EntryFor(path string) (C24Entry, bool) {
	for _, e := range C24SurfaceTable {
		if e.Path == path {
			return e, true
		}
	}
	return C24Entry{}, false
}

// The two ways a reference can be illegal. Both are contract violations; they
// are distinguished because the fix differs (add the entry to the table = human
// approval; declare the capability in the manifest = plugin author).
const (
	C24ViolationUnknown    = "unknown-entry"
	C24ViolationUndeclared = "undeclared-entry"
)

// The two kinds of ruler that can see a reference: the static string scan and
// the runtime proxy trap. A computed access (`wisp[g][m]`) is invisible to the
// first, which is exactly why both exist (AC#3 "scan + runtime double-check").
const (
	c24AccessStatic  = "static"
	c24AccessDynamic = "dynamic"
)

// C24Violation is one contract violation, printable as an audit line.
type C24Violation struct {
	Path       string
	Kind       string
	Capability string
	Access     string
}

// String renders the violation as the single line an audit log can grep for.
func (v C24Violation) String() string {
	cap := v.Capability
	if cap == "" {
		cap = "none"
	}
	return fmt.Sprintf("TIER2-CONTRACT-VIOLATION kind=%s path=%s capability=%s access=%s",
		v.Kind, v.Path, cap, v.Access)
}

// c24PathRe matches a `wisp.<group>.<method>...` reference in source text.
// Deliberately textual (C24: 字符串扫描可验): it also reports a mention inside a
// comment, which is the expected false-positive budget of a string ruler.
var c24PathRe = regexp.MustCompile(`\bwisp(?:\.[A-Za-z_$][A-Za-z0-9_$]*)+`)

// ScanC24Source is the static half of the double check: every host-entry
// reference in the source is judged against the frozen table and the tool's
// granted face. Granted entries pass, namespaces (`wisp.fs`) pass, everything
// else is a violation. Results are deduplicated and path-sorted so two runs of
// the same source read the same.
func ScanC24Source(source string, granted []C24Entry) []C24Violation {
	if !strings.Contains(source, c24RootName) {
		return nil
	}
	seen := map[string]bool{}
	var out []C24Violation
	for _, path := range c24PathRe.FindAllString(source, -1) {
		v, bad := judgeC24Path(path, granted, c24AccessStatic)
		if !bad {
			continue
		}
		key := v.Kind + "|" + v.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// judgeC24Path decides one reference. A name that is a proper namespace prefix
// of a table entry (`wisp.fs`) is legal - it is how a plugin holds a group - and
// a name that merely starts with the root (`wispfoo.bar`) never matches the
// regex in the first place.
func judgeC24Path(path string, granted []C24Entry, access string) (C24Violation, bool) {
	if e, known := C24EntryFor(path); known {
		for _, g := range granted {
			if g.Path == e.Path {
				return C24Violation{}, false
			}
		}
		return C24Violation{Path: e.Path, Kind: C24ViolationUndeclared, Capability: e.Capability, Access: access}, true
	}
	for _, e := range C24SurfaceTable {
		if strings.HasPrefix(e.Path, path+".") {
			return C24Violation{}, false
		}
	}
	return C24Violation{Path: path, Kind: C24ViolationUnknown, Access: access}, true
}

// C24Request is one host entry call, as the surface hands it over.
//
// Level is the frozen C24 column - the R1 floor for this entry - never a
// verdict: the host's assessor fuses it with R1..R9 and the concrete target.
type C24Request struct {
	PluginID   string
	Tool       string
	Entry      string
	Capability string
	Level      risk.Level
	// Args is the single JSON-decoded argument, nil when the plugin passed none.
	Args any
}

// C24Reply is the host's answer.
//
// Refusal is how the host says no: the surface turns it into a JS exception the
// plugin may catch, which is the mainstream host-API behaviour (fs throws in
// Node) and keeps a refusal distinct from a bug. A non-empty Refusal is NOT an
// error class of its own - the host has already booked the reason (audit row,
// C19 verdict); this field only carries the text to the plugin.
type C24Reply struct {
	Value   any    `json:"value,omitempty"`
	Refusal string `json:"refusal,omitempty"`
	// Truncated says the host capped its own payload. The surface logs it; the
	// 256KB result cap that applies to the tool outcome is block A's.
	Truncated bool `json:"truncated,omitempty"`
}

// C24Host is the Go side of the C24 table - one implementation, injected. It is
// the ONLY way a host entry reaches the outside world: no entry has a direct
// dependency on a tool, a file or a socket, so the capability check, the C19
// verdict, the taint mark and the tool_call row all sit on this one seam.
type C24Host interface {
	CallC24(ctx context.Context, req C24Request) (C24Reply, error)
}

// C24HostPanic is a host implementation that CRASHED while serving one entry.
//
// It exists because "something panicked" is not an actionable failure: the value
// carries the entry, its frozen capability and risk floor, the tool and plugin
// that asked, and the stack of the crash SITE (not the unwinding envelope), so
// the log line an operator greps and the error the model reads both say which
// host entry is broken. Recovery itself is not done here - this value is
// re-panicked on purpose and caught by GojaRuntime.Execute on the calling
// goroutine (ticket 51 AC#4).
type C24HostPanic struct {
	PluginID   string
	Tool       string
	Entry      string
	Capability string
	Level      risk.Level
	// Recovered is the value the host panicked with (a string, an error, a
	// runtime.Error - whatever the host implementation sent).
	Recovered any
	// Stack is the stack at the crash site, captured before unwinding.
	Stack string
}

// Error renders the attribution as one line. Plugin and tool are named because
// the same entry can be granted to several tools of the same plugin.
func (p *C24HostPanic) Error() string {
	if p == nil {
		return "host entry panicked"
	}
	return fmt.Sprintf("host entry %s (capability %s, risk floor %s) of plugin %q tool %q panicked inside its host implementation: %v",
		p.Entry, capabilityOrNone(p.Capability), p.Level, p.PluginID, p.Tool, p.Recovered)
}

// asC24HostPanic digs the attribution out of a recovered panic value. It is a
// typed assertion first (goja re-panics the exact value) with an errors.As
// fallback for a host that wrapped one.
func asC24HostPanic(v any) *C24HostPanic {
	if v == nil {
		return nil
	}
	if hp, ok := v.(*C24HostPanic); ok {
		return hp
	}
	if err, ok := v.(error); ok {
		var hp *C24HostPanic
		if errors.As(err, &hp) {
			return hp
		}
	}
	return nil
}

// C24Option configures a surface.
type C24Option func(*C24Surface)

// WithC24Logf replaces the audit sink (default: slog at Warn). One sink carries
// the per-call lines, the refusals and the contract violations because a host
// that wants severity routing can decide it from the line's own prefix.
func WithC24Logf(f func(format string, args ...any)) C24Option {
	return func(s *C24Surface) {
		if f != nil {
			s.logf = f
		}
	}
}

// C24Surface is one call's binding of the frozen table: the entries this tool
// may use, the host behind them, and the audit trail of what happened.
//
// It is built per call (never shared, never cached) and installed into that
// call's VM. The zero value is not usable; use NewC24Surface.
type C24Surface struct {
	pluginID string
	tool     string
	host     C24Host
	granted  []C24Entry
	logf     func(format string, args ...any)

	mu         sync.Mutex
	calls      []string
	violations []C24Violation
}

// NewC24Surface resolves the tool's declaration against the frozen table. A
// capability-free entry is always granted; every other entry needs its
// capability among declared. A nil host is accepted here and refused at install
// time, so the failure is one readable config error instead of a nil panic in a
// goja callback.
func NewC24Surface(pluginID, tool string, declared []string, host C24Host, opts ...C24Option) *C24Surface {
	s := &C24Surface{
		pluginID: pluginID,
		tool:     tool,
		host:     host,
		logf: func(format string, args ...any) {
			slog.Warn(fmt.Sprintf(format, args...))
		},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	have := make(map[string]bool, len(declared))
	for _, c := range declared {
		have[c] = true
	}
	for _, e := range C24SurfaceTable {
		if e.Capability == "" || have[e.Capability] {
			s.granted = append(s.granted, e)
		}
	}
	return s
}

// PluginID / Tool name the owner of this surface (the audit row is keyed on
// both).
func (s *C24Surface) PluginID() string { return s.pluginID }

// Tool is the manifest tool this surface serves.
func (s *C24Surface) Tool() string { return s.tool }

// Granted is the resolved face, in table order (a copy: the caller cannot widen
// the surface it was given).
func (s *C24Surface) Granted() []C24Entry {
	if s == nil {
		return nil
	}
	return append([]C24Entry(nil), s.granted...)
}

// Calls lists the entries the VM actually reached, in call order.
func (s *C24Surface) Calls() []string {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

// Violations lists the violations the RUNTIME trap caught (the static scan's
// findings are on GojaResult.ContractViolations).
func (s *C24Surface) Violations() []C24Violation {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]C24Violation(nil), s.violations...)
}

func (s *C24Surface) record(entry string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, entry)
}

func (s *C24Surface) logViolation(v C24Violation) {
	s.mu.Lock()
	s.violations = append(s.violations, v)
	s.mu.Unlock()
	s.logf("%s", v.String())
}

// c24Groups returns the six group names in table order.
func c24Groups() []string {
	var out []string
	seen := map[string]bool{}
	for _, e := range C24SurfaceTable {
		g := c24Group(e.Path)
		if g == "" || seen[g] {
			continue
		}
		seen[g] = true
		out = append(out, g)
	}
	return out
}

func c24Group(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) < 3 {
		return ""
	}
	return parts[1]
}

func c24MethodName(path string) string {
	parts := strings.Split(path, ".")
	return parts[len(parts)-1]
}

func (s *C24Surface) entriesIn(group string) []C24Entry {
	var out []C24Entry
	for _, e := range s.granted {
		if c24Group(e.Path) == group {
			out = append(out, e)
		}
	}
	return out
}

// c24StringifySource turns a JS argument into JSON with the engine's own
// serializer. Doing it in JS rather than with Value.Export keeps the semantics
// the plugin author expects (a function argument or a cycle is an error, an
// object is what JSON.stringify says it is) and keeps Go from walking a
// hostile object graph by hand.
const c24StringifySource = `(function(v){ return JSON.stringify(v) })`

const c24StringifyURL = "plugin://wisp/tier2-c24-args.js"

// installC24Surface builds the `wisp` global inside the call's VM. Only granted
// entries exist as properties; the Proxy trap is what turns an access to
// anything else into a logged contract violation while still answering
// `undefined` (a plugin that asks for a capability it does not have gets the
// same result as one asking for a name that does not exist - both are refused
// the same way, and both are logged).
func installC24Surface(ctx context.Context, vm *goja.Runtime, s *C24Surface) error {
	if s == nil {
		return nil
	}
	if s.host == nil {
		return observe.New(observe.ClassConfig, fmt.Sprintf(
			"plugin %q: a tier2 call carries a C24 surface with no host wired - "+
				"refusing to run JS whose only host entries could never answer (fail closed)", s.pluginID))
	}
	stringify, err := c24Stringifier(vm)
	if err != nil {
		return observe.Wrap(observe.ClassInternal, err,
			fmt.Sprintf("plugin %q: cannot compile the C24 argument serializer", s.pluginID))
	}

	root := vm.NewObject()
	for _, group := range c24Groups() {
		target := vm.NewObject()
		for _, e := range s.entriesIn(group) {
			fn := vm.ToValue(s.bind(ctx, vm, stringify, e))
			if err := target.DefineDataProperty(c24MethodName(e.Path), fn,
				goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
				return observe.Wrap(observe.ClassInternal, err,
					fmt.Sprintf("plugin %q: cannot define %s", s.pluginID, e.Path))
			}
		}
		groupName := group
		proxy := vm.NewProxy(target, &goja.ProxyTrapConfig{
			Get: func(t *goja.Object, prop string, _ goja.Value) goja.Value {
				if v := c24Inherited(t, prop); v != nil {
					return v
				}
				if v, bad := judgeC24Path(c24RootName+"."+groupName+"."+prop, s.granted, c24AccessDynamic); bad {
					s.logViolation(v)
				}
				return goja.Undefined()
			},
		})
		if err := root.DefineDataProperty(groupName, vm.ToValue(proxy),
			goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
			return observe.Wrap(observe.ClassInternal, err,
				fmt.Sprintf("plugin %q: cannot define the %s group", s.pluginID, groupName))
		}
	}

	rootProxy := vm.NewProxy(root, &goja.ProxyTrapConfig{
		Get: func(t *goja.Object, prop string, _ goja.Value) goja.Value {
			if v := c24Inherited(t, prop); v != nil {
				return v
			}
			if v, bad := judgeC24Path(c24RootName+"."+prop, s.granted, c24AccessDynamic); bad {
				s.logViolation(v)
			}
			return goja.Undefined()
		},
	})
	if err := vm.GlobalObject().DefineDataProperty(c24RootName, vm.ToValue(rootProxy),
		goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE); err != nil {
		return observe.Wrap(observe.ClassInternal, err,
			fmt.Sprintf("plugin %q: cannot define the %s global", s.pluginID, c24RootName))
	}
	return nil
}

// c24Inherited returns a property's own or inherited value, or nil when the
// name does not exist at all. The prototype walk matters for the audit log: a
// plugin calling `wisp.toString` is touching Object.prototype, not a host
// entry, and logging that as an unknown entry would be a false violation.
func c24Inherited(t *goja.Object, prop string) goja.Value {
	if v := t.Get(prop); v != nil {
		return v
	}
	if p := t.Prototype(); p != nil {
		return p.Get(prop)
	}
	return nil
}

// scanC24Call is the static half of AC#3's double check, run by
// GojaRuntime.Execute before any VM exists. It reports every reference to a
// host entry that the frozen table or the tool's declaration does not cover,
// logs each one (through the surface's audit sink when there is a surface), and
// hands the lines back so they travel on GojaResult.
//
// The call still runs: a reference to an entry that was never installed is
// `undefined` in the VM, which is exactly the contract's answer, and a plugin
// that merely mentions a name in a dead branch is not executed for it.
func scanC24Call(call GojaCall) []string {
	var granted []C24Entry
	if call.Surface != nil {
		granted = call.Surface.Granted()
	}
	found := ScanC24Source(call.Source, granted)
	if len(found) == 0 {
		return nil
	}
	out := make([]string, 0, len(found))
	for _, v := range found {
		out = append(out, v.String())
		if call.Surface != nil {
			call.Surface.logf("%s", v.String())
			continue
		}
		slog.Warn(v.String())
	}
	return out
}

// c24Stringifier compiles the argument serializer into the call's own VM.
func c24Stringifier(vm *goja.Runtime) (goja.Callable, error) {
	prg, err := goja.Compile(c24StringifyURL, c24StringifySource, false)
	if err != nil {
		return nil, err
	}
	v, err := vm.RunProgram(prg)
	if err != nil {
		return nil, err
	}
	fn, ok := goja.AssertFunction(v)
	if !ok {
		return nil, fmt.Errorf("the C24 argument serializer did not evaluate to a function")
	}
	return fn, nil
}

// callHost is the single seam between one JS entry and the host implementation.
// Beyond the call it owns one thing: a host CRASH is recorded and attributed
// before the panic leaves this frame. Recording first matters - the call that
// takes the host down is exactly the call an operator needs on file - and so
// does attributing first: the goja frames that unwind the VM do not promise to
// keep a panic value intact, and a bare "tier 2 panicked" line would leave
// nobody able to name the broken entry.
//
// The panic is re-raised as *C24HostPanic, never converted into a JS exception:
// see the file header. Recovery belongs to GojaRuntime.Execute, on the calling
// goroutine, once.
func (s *C24Surface) callHost(ctx context.Context, e C24Entry, args any) (reply C24Reply, err error) {
	defer func() {
		if p := recover(); p != nil {
			hp := &C24HostPanic{
				PluginID:   s.pluginID,
				Tool:       s.tool,
				Entry:      e.Path,
				Capability: e.Capability,
				Level:      e.Level,
				Recovered:  p,
				Stack:      string(debug.Stack()),
			}
			s.record(e.Path)
			s.logf("tier2: C24-HOSTPANIC plugin=%s tool=%s entry=%s capability=%s level=%s panic=%v",
				s.pluginID, s.tool, e.Path, capabilityOrNone(e.Capability), e.Level, p)
			slog.Error("tier2 C24 host implementation panicked",
				"plugin", s.pluginID, "tool", s.tool, "entry", e.Path,
				"capability", capabilityOrNone(e.Capability), "panic", fmt.Sprint(p), "stack", hp.Stack)
			panic(hp)
		}
	}()
	reply, err = s.host.CallC24(ctx, C24Request{
		PluginID:   s.pluginID,
		Tool:       s.tool,
		Entry:      e.Path,
		Capability: e.Capability,
		Level:      e.Level,
		Args:       args,
	})
	s.record(e.Path)
	return reply, err
}

// bind wraps one frozen entry as the JS function the plugin calls. Everything
// the page can observe - the argument decode, the host call, the refusal, the
// audit line - happens here, on the calling goroutine, inside the VM.
func (s *C24Surface) bind(ctx context.Context, vm *goja.Runtime, stringify goja.Callable, e C24Entry) func(goja.FunctionCall) goja.Value {
	return func(fc goja.FunctionCall) goja.Value {
		args, err := c24Arguments(stringify, e, fc)
		if err != nil {
			panic(vm.NewTypeError("%s", err.Error()))
		}
		reply, herr := s.callHost(ctx, e, args)
		if herr != nil {
			s.logf("tier2: C24-HOSTERR plugin=%s tool=%s entry=%s err=%v",
				s.pluginID, s.tool, e.Path, herr)
			panic(vm.NewGoError(herr))
		}
		s.logf("tier2: C24-CALL plugin=%s tool=%s entry=%s capability=%s level=%s truncated=%v",
			s.pluginID, s.tool, e.Path, capabilityOrNone(e.Capability), e.Level, reply.Truncated)
		if reply.Refusal != "" {
			s.logf("tier2: C24-REFUSE plugin=%s tool=%s entry=%s reason=%s",
				s.pluginID, s.tool, e.Path, reply.Refusal)
			panic(vm.NewGoError(fmt.Errorf("%s refused: %s", e.Path, reply.Refusal)))
		}
		if reply.Value == nil {
			return goja.Undefined()
		}
		return vm.ToValue(reply.Value)
	}
}

// c24Arguments enforces the one calling convention the surface has: at most one
// JSON-serializable argument. A second argument, a function, a BigInt or a
// cyclic object is a TypeError - the plugin author's bug, reported as such,
// never silently coerced into a host call.
func c24Arguments(stringify goja.Callable, e C24Entry, fc goja.FunctionCall) (any, error) {
	if len(fc.Arguments) > 1 {
		return nil, fmt.Errorf("%s takes at most one object argument, got %d", e.Path, len(fc.Arguments))
	}
	if len(fc.Arguments) == 0 {
		return nil, nil
	}
	v := fc.Arguments[0]
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return nil, nil
	}
	out, err := stringify(goja.Undefined(), v)
	if err != nil {
		return nil, fmt.Errorf("%s argument is not JSON-serializable (host entries take data): %v", e.Path, err)
	}
	if out == nil || goja.IsUndefined(out) {
		return nil, fmt.Errorf("%s argument is not JSON-serializable (host entries take data, not functions)", e.Path)
	}
	var decoded any
	if err := json.Unmarshal([]byte(out.String()), &decoded); err != nil {
		return nil, fmt.Errorf("%s argument is not JSON-serializable: %v", e.Path, err)
	}
	return decoded, nil
}

func capabilityOrNone(c string) string {
	if c == "" {
		return "none"
	}
	return c
}
