package plugin

// Ticket 51 block B - the C24 GojaHostAPI contract tests (acceptance criterion
// "C24 surface").
//
// What is under test: internal/plugin/c24.go - the frozen table, the string
// scan, and the runtime object the VM actually sees. The three claims:
//
//  1. the table is the owner-approved six-group set (A1541 item 1: http, fs,
//     clipboard, notify, sysinfo, state) and nothing else;
//  2. every listed entry reaches the wired host with its frozen capability and
//     risk floor, and a plugin only sees the entries its own declaration covers;
//  3. an entry outside the table, or outside the plugin's declaration, is
//     undefined - logged by the static scan AND by the runtime trap, because a
//     computed access (`wisp[g][m]`) is invisible to a string scan.
//
// The tests speak to the surface through goja (a real VM, real JS), never by
// calling the Go methods behind it: "e2e" here means the JS text a plugin would
// ship.

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// allC24Caps is the frozen C3 token set (SPEC-07 §2) spelled out instead of
// imported: internal/plugin cannot import internal/tools (tools -> agent ->
// memory -> plugin), so the twelve-word list is copied here on purpose. A token
// that is not in this list is not a capability, and a C24 row naming one would
// be a contract violation.
var allC24Caps = []string{
	"fs.read", "fs.write", "net", "clipboard", "shell", "sysinfo",
	"screen", "input", "window", "notify", "memory",
}

// c24Recorder is the host side of the surface: it records every call the VM
// makes and answers with a canned reply, so "did the entry reach the host?" is
// a reading rather than an inference.
type c24Recorder struct {
	mu    sync.Mutex
	calls []C24Request
	reply func(C24Request) C24Reply
}

func (h *c24Recorder) CallC24(_ context.Context, req C24Request) (C24Reply, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.calls = append(h.calls, req)
	if h.reply != nil {
		return h.reply(req), nil
	}
	return C24Reply{Value: map[string]any{"entry": req.Entry}}, nil
}

func (h *c24Recorder) Calls() []C24Request {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]C24Request(nil), h.calls...)
}

func (h *c24Recorder) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.calls)
}

// c24Log captures the surface's own log lines (audit lines and
// contract-violation lines) without touching the global logger.
type c24Log struct {
	mu    sync.Mutex
	lines []string
}

func (l *c24Log) Logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *c24Log) All() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.lines...)
}

func (l *c24Log) CountContaining(sub string) int {
	n := 0
	for _, line := range l.All() {
		if strings.Contains(line, sub) {
			n++
		}
	}
	return n
}

func (l *c24Log) Joined() string { return strings.Join(l.All(), "\n") }

// runC24 executes one snippet with a surface built from the given declaration.
func runC24(t *testing.T, host *c24Recorder, log *c24Log, declared []string, source string) (GojaResult, error) {
	t.Helper()
	opts := []C24Option{}
	if log != nil {
		opts = append(opts, WithC24Logf(log.Logf))
	}
	surface := NewC24Surface("demo", "echo", declared, host, opts...)
	r := NewGojaRuntime(WithGojaRegistry(observe.NewRegistry()))
	return r.Execute(context.Background(), GojaCall{
		PluginID: "demo",
		Source:   source,
		Surface:  surface,
	})
}

// TestC24SurfaceIsFrozenAndOwnerApproved pins the table itself: the owner
// approved exactly these nine entries in six groups (A1541 item 1), so a future
// edit that adds, drops, renames or re-levels one is a contract change and must
// be a human approval - this test is where that shows up first.
func TestC24SurfaceIsFrozenAndOwnerApproved(t *testing.T) {
	t.Parallel()
	want := []struct {
		path string
		cap  string
		li   risk.Level
	}{
		{"wisp.http.request", "net", risk.L1},
		{"wisp.fs.read", "fs.read", risk.L0},
		{"wisp.fs.write", "fs.write", risk.L1},
		{"wisp.clipboard.read", "clipboard", risk.L0},
		{"wisp.clipboard.write", "clipboard", risk.L1},
		{"wisp.notify.send", "notify", risk.L0},
		{"wisp.sysinfo.get", "sysinfo", risk.L0},
		{"wisp.state.get", "", risk.L0},
		{"wisp.state.set", "", risk.L0},
	}
	if len(C24SurfaceTable) != len(want) {
		t.Fatalf("C24 table has %d entries, the owner-approved set has %d", len(C24SurfaceTable), len(want))
	}
	for i, w := range want {
		got := C24SurfaceTable[i]
		if got.Path != w.path || got.Capability != w.cap || got.Level != w.li {
			t.Errorf("C24 row %d = {%q %q %v}, want {%q %q %v}",
				i, got.Path, got.Capability, got.Level, w.path, w.cap, w.li)
		}
	}
	groups := map[string]bool{}
	seen := map[string]bool{}
	for _, e := range C24SurfaceTable {
		parts := strings.Split(e.Path, ".")
		if len(parts) != 3 || parts[0] != c24RootName {
			t.Errorf("entry %q is not <root>.<group>.<method>", e.Path)
			continue
		}
		groups[parts[1]] = true
		if e.Capability != "" && !containsString(allC24Caps, e.Capability) {
			t.Errorf("entry %q names capability %q, which is outside the frozen C3 set", e.Path, e.Capability)
		}
		if seen[e.Path] {
			t.Errorf("entry %q is listed twice", e.Path)
		}
		seen[e.Path] = true
		if e.Level == risk.Deny {
			t.Errorf("entry %q carries Deny: the C24 column is the R1 floor, Deny belongs to C19/R3", e.Path)
		}
		if strings.TrimSpace(e.Doc) == "" {
			t.Errorf("entry %q has no doc line (the C24 table is also the plugin-facing reference)", e.Path)
		}
		if _, ok := C24EntryFor(e.Path); !ok {
			t.Errorf("C24EntryFor(%q) says missing while the table carries it", e.Path)
		}
	}
	wantGroups := []string{"clipboard", "fs", "http", "notify", "state", "sysinfo"}
	gotGroups := make([]string, 0, len(groups))
	for g := range groups {
		gotGroups = append(gotGroups, g)
	}
	sort.Strings(gotGroups)
	if !reflect.DeepEqual(gotGroups, wantGroups) {
		t.Fatalf("groups = %v, want the owner-approved six %v", gotGroups, wantGroups)
	}
	if _, ok := C24EntryFor("wisp.screen.capture"); ok {
		t.Errorf("wisp.screen.capture is not in the approved set but C24EntryFor found it")
	}
}

// TestC24ScanRuler pins the static half of the double check. The ruler is
// deliberately TEXTUAL (C24: "JS 侧出现任何未列出的宿主入口即判契约违规（字符串
// 扫描可验）"), so a mention inside a comment is also a hit - the scan reports,
// it does not silently pass.
func TestC24ScanRuler(t *testing.T) {
	t.Parallel()
	granted := []C24Entry{mustEntry(t, "wisp.fs.read"), mustEntry(t, "wisp.state.get")}

	cases := []struct {
		name     string
		source   string
		wantPath []string
		wantKind []string
	}{
		{
			name:   "granted entry is clean",
			source: `wisp.fs.read({path: "x"})`,
		},
		{
			name:     "table entry outside the declaration",
			source:   `wisp.fs.write({path: "x", data: "y"})`,
			wantPath: []string{"wisp.fs.write"},
			wantKind: []string{C24ViolationUndeclared},
		},
		{
			name:     "entry outside the table is unknown",
			source:   `wisp.fetch("x")`,
			wantPath: []string{"wisp.fetch"},
			wantKind: []string{C24ViolationUnknown},
		},
		{
			name:     "a builtin tool name is not a C24 entry",
			source:   `wisp.screen.capture({})`,
			wantPath: []string{"wisp.screen.capture"},
			wantKind: []string{C24ViolationUnknown},
		},
		{
			name:   "the group namespace itself is legal",
			source: `var g = wisp.fs; var s = wisp.state;`,
		},
		{
			name:   "a longer name is not a prefix match",
			source: `wisp.fs.read.deeper({})`,
			// wisp.fs.read.deeper is neither an entry nor a proper prefix.
			wantPath: []string{"wisp.fs.read.deeper"},
			wantKind: []string{C24ViolationUnknown},
		},
		{
			name:   "the bare root is not a violation",
			source: `var w = typeof wisp; var q = wispfoo.bar;`,
		},
		{
			name:     "textual ruler catches a comment mention too",
			source:   "// uses wisp.http.request\nwisp.fs.read({path: 'x'})",
			wantPath: []string{"wisp.http.request"},
			wantKind: []string{C24ViolationUndeclared},
		},
		{
			name:     "duplicates collapse",
			source:   `wisp.fs.write({}); wisp.fs.write({});`,
			wantPath: []string{"wisp.fs.write"},
			wantKind: []string{C24ViolationUndeclared},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ScanC24Source(tc.source, granted)
			if len(got) != len(tc.wantPath) {
				t.Fatalf("scan found %d violation(s) %v, want %d", len(got), got, len(tc.wantPath))
			}
			for i := range tc.wantPath {
				if got[i].Path != tc.wantPath[i] {
					t.Errorf("violation %d path = %q, want %q", i, got[i].Path, tc.wantPath[i])
				}
				if got[i].Kind != tc.wantKind[i] {
					t.Errorf("violation %d kind = %q, want %q", i, got[i].Kind, tc.wantKind[i])
				}
				if !strings.Contains(got[i].String(), "TIER2-CONTRACT-VIOLATION") {
					t.Errorf("violation %d does not print the auditable token: %q", i, got[i].String())
				}
			}
		})
	}
}

// TestC24SurfaceExposesExactlyTheGrantedEntries is the visibility half: the
// declared capability decides which entries exist, and state.* exists for every
// plugin because it is the plugin's own data (capability-free by contract).
//
// The two negative probes are written as COMPUTED accesses on purpose. A plain
// `typeof wisp.fs.write` would also be a static-ruler hit (the ruler is
// textual), so the source here must not spell the entries it is proving absent;
// the runtime trap is what answers, and this test keeps the static result empty.
func TestC24SurfaceExposesExactlyTheGrantedEntries(t *testing.T) {
	t.Parallel()
	host := &c24Recorder{}
	out, err := runC24(t, host, nil, []string{"fs.read"}, `(function(){
		var nm = function(a, b){ return a + b };
		return JSON.stringify({
			root: typeof wisp,
			granted: typeof wisp.fs.read,
			ungranted: typeof wisp.fs[nm("wr", "ite")],
			otherGroup: typeof wisp.http[nm("req", "uest")],
			stateGet: typeof wisp.state.get,
			stateSet: typeof wisp.state.set,
		})
	})()`)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(out.Value), &got); err != nil {
		t.Fatalf("result %q is not JSON: %v", out.Value, err)
	}
	want := map[string]string{
		"root": "object", "granted": "function", "ungranted": "undefined",
		"otherGroup": "undefined", "stateGet": "function", "stateSet": "function",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("surface face = %v, want %v", got, want)
	}
	if len(out.ContractViolations) != 0 {
		t.Fatalf("a computed probe produced a STATIC violation: %v", out.ContractViolations)
	}
	if n := host.Count(); n != 0 {
		t.Fatalf("typeof probes reached the host %d time(s)", n)
	}
}

// TestC24AddsNoGlobalsBeyondTheSurface is the "no globals beyond the surface"
// claim measured the only way it can be: diff the VM's own global names against
// the same VM with no surface. The diff must be exactly ["wisp"], which also
// proves no timer, no fetch, no require and no process appeared.
func TestC24AddsNoGlobalsBeyondTheSurface(t *testing.T) {
	t.Parallel()
	const names = `(function(){ return JSON.stringify(Object.getOwnPropertyNames(globalThis).sort()) })()`

	baseline, err := NewGojaRuntime(WithGojaRegistry(observe.NewRegistry())).Execute(context.Background(),
		GojaCall{PluginID: "globals", Source: names})
	if err != nil {
		t.Fatalf("baseline globals: %v", err)
	}
	withSurface, err := runC24(t, &c24Recorder{}, nil, allC24Caps, names)
	if err != nil {
		t.Fatalf("surface globals: %v", err)
	}

	var want, got []string
	if err := json.Unmarshal([]byte(baseline.Value), &want); err != nil {
		t.Fatalf("baseline list %q: %v", baseline.Value, err)
	}
	if err := json.Unmarshal([]byte(withSurface.Value), &got); err != nil {
		t.Fatalf("surface list %q: %v", withSurface.Value, err)
	}
	var added []string
	for _, n := range got {
		if !containsString(want, n) {
			added = append(added, n)
		}
	}
	if !reflect.DeepEqual(added, []string{"wisp"}) {
		t.Fatalf("surface added globals %v, want exactly [wisp]", added)
	}

	banned := []string{
		"fetch", "XMLHttpRequest", "WebSocket", "require", "process", "module",
		"exports", "console", "setTimeout", "setInterval", "setImmediate",
	}
	for _, b := range banned {
		if containsString(got, b) {
			t.Errorf("the VM exposes banned global %q: goja ships no host object and the surface adds none", b)
		}
	}
}

// TestC24EveryEntryReachesTheHost is AC#3's "all listed entries work e2e (mock
// plugin each)": one JS snippet per row, a real VM, a recording host.
func TestC24EveryEntryReachesTheHost(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path    string
		source  string
		wantArg any
	}{
		{
			"wisp.http.request",
			`(function(){ return wisp.http.request({url: "https://api.example.com/x", method: "GET"}) })()`,
			map[string]any{"url": "https://api.example.com/x", "method": "GET"},
		},
		{
			"wisp.fs.read",
			`(function(){ return wisp.fs.read({path: "D:/work/notes.txt"}) })()`,
			map[string]any{"path": "D:/work/notes.txt"},
		},
		{
			"wisp.fs.write",
			`(function(){ return wisp.fs.write({path: "D:/work/notes.txt", data: "hi"}) })()`,
			map[string]any{"path": "D:/work/notes.txt", "data": "hi"},
		},
		{
			"wisp.clipboard.read",
			`(function(){ return wisp.clipboard.read({}) })()`,
			map[string]any{},
		},
		{
			"wisp.clipboard.write",
			`(function(){ return wisp.clipboard.write({text: "hi"}) })()`,
			map[string]any{"text": "hi"},
		},
		{
			"wisp.notify.send",
			`(function(){ return wisp.notify.send({title: "t", body: "b"}) })()`,
			map[string]any{"title": "t", "body": "b"},
		},
		{
			"wisp.sysinfo.get",
			`(function(){ return wisp.sysinfo.get({keys: ["battery"]}) })()`,
			map[string]any{"keys": []any{"battery"}},
		},
		{
			"wisp.state.get",
			`(function(){ return wisp.state.get({key: "counter"}) })()`,
			map[string]any{"key": "counter"},
		},
		{
			"wisp.state.set",
			`(function(){ return wisp.state.set({key: "counter", value: 3}) })()`,
			map[string]any{"key": "counter", "value": float64(3)},
		},
	}
	covered := map[string]bool{}
	for _, tc := range cases {
		covered[tc.path] = true
	}
	for _, e := range C24SurfaceTable {
		if !covered[e.Path] {
			t.Fatalf("entry %q has no e2e case: every row of the frozen table must be proven to reach the host", e.Path)
		}
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			host := &c24Recorder{}
			out, err := runC24(t, host, nil, allC24Caps, tc.source)
			if err != nil {
				t.Fatalf("%s: %v", tc.path, err)
			}
			if len(out.ContractViolations) != 0 {
				t.Fatalf("%s produced contract violations: %v", tc.path, out.ContractViolations)
			}
			calls := host.Calls()
			if len(calls) != 1 {
				t.Fatalf("%s reached the host %d time(s), want 1", tc.path, len(calls))
			}
			got := calls[0]
			want, _ := C24EntryFor(tc.path)
			if got.Entry != want.Path || got.Capability != want.Capability || got.Level != want.Level {
				t.Errorf("host saw {%q %q %v}, want the frozen row {%q %q %v}",
					got.Entry, got.Capability, got.Level, want.Path, want.Capability, want.Level)
			}
			if got.PluginID != "demo" || got.Tool != "echo" {
				t.Errorf("host call lost its owner: plugin=%q tool=%q", got.PluginID, got.Tool)
			}
			if !reflect.DeepEqual(got.Args, tc.wantArg) {
				t.Errorf("host args = %#v, want %#v", got.Args, tc.wantArg)
			}
			if !strings.Contains(out.Value, want.Path) {
				t.Errorf("the host's reply did not come back to JS: %q", out.Value)
			}
		})
	}
}

// TestC24UndeclaredEntryIsUndefinedAndLoggedAtRuntime is the runtime half of
// the double check: a COMPUTED access is invisible to the string scan, so only
// the proxy trap can log it. Both halves are measured here.
func TestC24UndeclaredEntryIsUndefinedAndLoggedAtRuntime(t *testing.T) {
	t.Parallel()
	host := &c24Recorder{}
	log := &c24Log{}
	source := `(function(){
		var g = "fs"; var m = "write";
		return typeof wisp[g][m];
	})()`
	out, err := runC24(t, host, log, []string{"fs.read"}, source)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if got := strings.Trim(out.Value, `"`); got != "undefined" {
		t.Fatalf("computed access returned %q, want undefined", got)
	}
	if len(out.ContractViolations) != 0 {
		t.Fatalf("the static ruler cannot see a computed access, yet it reported %v", out.ContractViolations)
	}
	if n := log.CountContaining("TIER2-CONTRACT-VIOLATION"); n != 1 {
		t.Fatalf("runtime trap logged %d violation line(s), want exactly 1:\n%s", n, log.Joined())
	}
	line := log.Joined()
	if !strings.Contains(line, "path=wisp.fs.write") || !strings.Contains(line, "kind=undeclared-entry") ||
		!strings.Contains(line, "access=dynamic") {
		t.Fatalf("runtime violation line is not auditable: %q", line)
	}
	if host.Count() != 0 {
		t.Fatalf("a refused entry reached the host")
	}

	// An unknown group is seen by BOTH rulers: the scan flags the literal
	// `wisp.screen`, and the proxy trap flags the access itself. Two lines,
	// one per ruler, both naming the same path - which is what the double check
	// is for.
	log2 := &c24Log{}
	out2, err := runC24(t, &c24Recorder{}, log2, []string{"fs.read"},
		`(function(){ return typeof wisp.screen })()`)
	if err != nil {
		t.Fatalf("unknown group call: %v", err)
	}
	if strings.Trim(out2.Value, `"`) != "undefined" {
		t.Fatalf("unknown group is not undefined: %q", out2.Value)
	}
	if len(out2.ContractViolations) != 1 || !strings.Contains(out2.ContractViolations[0], "path=wisp.screen") {
		t.Fatalf("static violations = %v, want one naming wisp.screen", out2.ContractViolations)
	}
	if n := log2.CountContaining("kind=unknown-entry"); n != 2 {
		t.Fatalf("unknown group logged %d line(s), want one per ruler:\n%s", n, log2.Joined())
	}
	both := log2.Joined()
	if !strings.Contains(both, "access=static") || !strings.Contains(both, "access=dynamic") {
		t.Fatalf("the two lines do not name which ruler saw it: %q", both)
	}
}

// TestC24StaticViolationIsLoggedAndStillRuns: a reference to an entry outside
// the table is reported by the scan (and logged) but the call itself runs - the
// reference is undefined at runtime, exactly as the contract says, and a plugin
// that only mentions a name in a dead branch is not executed for it.
func TestC24StaticViolationIsLoggedAndStillRuns(t *testing.T) {
	t.Parallel()
	log := &c24Log{}
	out, err := runC24(t, &c24Recorder{}, log, allC24Caps,
		`(function(){ if (false) { wisp.screen.capture({}) } return "ran" })()`)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if got := strings.Trim(out.Value, `"`); got != "ran" {
		t.Fatalf("call did not run: %q", got)
	}
	if len(out.ContractViolations) != 1 || !strings.Contains(out.ContractViolations[0], "path=wisp.screen.capture") {
		t.Fatalf("static violations = %v, want one naming wisp.screen.capture", out.ContractViolations)
	}
	if n := log.CountContaining("TIER2-CONTRACT-VIOLATION"); n != 1 {
		t.Fatalf("violation was reported but not logged (%d line(s))", n)
	}
}

// TestC24RefusalFromTheHostBecomesAJSThrow: a host refusal is the plugin's to
// catch (like fs throwing in Node), and an uncaught one is a tool-class error
// so the model can self-correct instead of the task dying.
func TestC24RefusalFromTheHostBecomesAJSThrow(t *testing.T) {
	t.Parallel()
	host := &c24Recorder{reply: func(req C24Request) C24Reply {
		return C24Reply{Refusal: "path D:/x is outside the authorized directories"}
	}}
	out, err := runC24(t, host, nil, []string{"fs.read"}, `(function(){
		try { wisp.fs.read({path: "D:/x"}) } catch (e) { return "caught: " + e.message }
		return "not reached"
	})()`)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if !strings.Contains(out.Value, "caught: ") || !strings.Contains(out.Value, "outside the authorized directories") {
		t.Fatalf("the host refusal did not reach the plugin as an exception: %q", out.Value)
	}

	uncaught, err := runC24(t, host, nil, []string{"fs.read"}, `(function(){ return wisp.fs.read({path: "D:/x"}) })()`)
	if err == nil {
		t.Fatalf("an uncaught host refusal must fail the call, got %q", uncaught.Value)
	}
	if class := observe.ClassOfOrInternal(err); class != "tool" {
		t.Fatalf("uncaught refusal class = %q, want tool", class)
	}
	if !strings.Contains(err.Error(), "outside the authorized directories") {
		t.Fatalf("the refusal text was lost: %v", err)
	}
}

// TestC24StateIsCapabilityFree: the plugin's own KV needs no capability (C24 /
// SPEC-07 §5: "插件自身 KV，无 capability——插件自己的数据").
func TestC24StateIsCapabilityFree(t *testing.T) {
	t.Parallel()
	host := &c24Recorder{}
	out, err := runC24(t, host, nil, nil, `(function(){
		return JSON.stringify([
			wisp.state.get({key: "counter"}),
			wisp.state.set({key: "counter", value: 1}),
		])
	})()`)
	if err != nil {
		t.Fatalf("call with no declarations: %v", err)
	}
	if len(out.ContractViolations) != 0 {
		t.Fatalf("state.* was reported as a violation: %v", out.ContractViolations)
	}
	calls := host.Calls()
	if len(calls) != 2 {
		t.Fatalf("state calls reached the host %d time(s), want 2", len(calls))
	}
	for i, c := range calls {
		if c.Capability != "" {
			t.Errorf("state call %d carried capability %q, want capability-free", i, c.Capability)
		}
		if c.Level != risk.L0 {
			t.Errorf("state call %d level = %v, want L0", i, c.Level)
		}
	}
}

// TestC24SurfaceDiesWithTheCall: the surface is installed into the per-call VM,
// so the next call of the same plugin cannot see the previous call's grants.
func TestC24SurfaceDiesWithTheCall(t *testing.T) {
	t.Parallel()
	r := NewGojaRuntime(WithGojaRegistry(observe.NewRegistry()))
	first := NewC24Surface("demo", "echo", []string{"fs.read"}, &c24Recorder{})
	if _, err := r.Execute(context.Background(), GojaCall{
		PluginID: "demo", Source: `(function(){ return typeof wisp.fs.read })()`, Surface: first,
	}); err != nil {
		t.Fatalf("first call: %v", err)
	}

	second := NewC24Surface("demo", "echo", nil, &c24Recorder{})
	out, err := r.Execute(context.Background(), GojaCall{
		PluginID: "demo", Source: `(function(){ return typeof wisp.fs.read })()`, Surface: second,
	})
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if got := strings.Trim(out.Value, `"`); got != "undefined" {
		t.Fatalf("the second call still saw the first call's grant: %q", got)
	}
}

// TestC24RejectsBadArguments: the uniform calling convention is one
// JSON-serializable object, and anything else is a JS TypeError the plugin
// author can read - never a silently coerced call.
func TestC24RejectsBadArguments(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		source string
		want   string
	}{
		{"two arguments", `(function(){ wisp.fs.read({path: "a"}, {path: "b"}) })()`, "at most one"},
		{"a function is not data", `(function(){ wisp.fs.read(function(){}) })()`, "JSON-serializable"},
		{"a cyclic object", `(function(){ var a = {}; a.self = a; wisp.fs.read(a) })()`, "circular"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			host := &c24Recorder{}
			_, err := runC24(t, host, nil, allC24Caps, tc.source)
			if err == nil {
				t.Fatalf("bad argument was accepted")
			}
			if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.want)) {
				t.Fatalf("error %q does not explain the shape (%q)", err, tc.want)
			}
			if class := observe.ClassOfOrInternal(err); class != "tool" {
				t.Fatalf("class = %q, want tool", class)
			}
			if host.Count() != 0 {
				t.Fatalf("a call with a bad argument still reached the host")
			}
		})
	}
}

// TestC24SurfaceWithoutHostFailsClosed: a surface with no host is an assembly
// bug, and it must not run the plugin's code (nothing it calls could work).
func TestC24SurfaceWithoutHostFailsClosed(t *testing.T) {
	t.Parallel()
	r := NewGojaRuntime(WithGojaRegistry(observe.NewRegistry()))
	surface := NewC24Surface("demo", "echo", allC24Caps, nil)
	_, err := r.Execute(context.Background(), GojaCall{
		PluginID: "demo", Source: `(function(){ return "should not run" })()`, Surface: surface,
	})
	if err == nil {
		t.Fatalf("a host-less surface was accepted")
	}
	if class := observe.ClassOfOrInternal(err); class != "config" {
		t.Fatalf("class = %q, want config (an assembly bug, not a plugin bug)", class)
	}
}

// --- helpers ---------------------------------------------------------------

func mustEntry(t *testing.T, path string) C24Entry {
	t.Helper()
	e, ok := C24EntryFor(path)
	if !ok {
		t.Fatalf("C24EntryFor(%q) missing", path)
	}
	return e
}
