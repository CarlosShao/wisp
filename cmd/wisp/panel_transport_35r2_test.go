//go:build windows

package main

// Ticket 35 AC#6 (`:52`), leg 35-r2: the page<->host transport judged BEHAVIOURALLY.
//
// WHAT CHANGED AND WHY. panel_transport_35r1_test.go (kept untouched, assertions and
// all) decided "the forwarding is installed" by reading the Init script's wording
// (forwardingInstalled(), three strings.Contains) and then handed msgcb a clean RPC
// frame it had built itself in Go. 35-v1 named both halves (verdict G1/G2-4, ledger
// A682): the two mutants reddened on the *wording* gate, never on arrival, and the fake
// was more forgiving than the library because it bypassed the one property the whole bug
// lives on - chrome.webview.postMessage, which go-webview2's own document script
// (pkg/edge/chromium.go:112) reads back at call time as window.external.invoke's body.
// A forwarding hook that overwrites that property therefore loops in a real browser.
//
// This file models that hop instead of describing it. A tiny JavaScript interpreter
// (subset listed below) executes, verbatim:
//
//	(a) the library's external.invoke shim, character for character from webview.go's
//	    caller pkg/edge/chromium.go:112;
//	(b) the per-binding stub script, built by the same concatenation go-webview2 uses at
//	    webview.go:462 ("(function() { var name = " + jsString(name) + ";" + <body>), from
//	    the library's own text;
//	(c) the PRODUCTION runtime string panelPostMessageForwardInit - read as a value from
//	    the shipping constant, never copied into this file;
//	(d) the reply scripts msgcb evaluates at webview.go:144/150/157.
//
// The interpreter subset is exactly: var / if / else / return / block / try-finally /
// expression statements; identifiers, member and computed access, assignment, calls and
// new; function expressions and one-argument arrows; object and array literals (trailing
// commas allowed); postfix ++; !, typeof, &&, ||, ===, !==, ==, !=, +; literals
// number/string/true/false/null/undefined/this/arguments; and host JSON.stringify,
// Array.prototype.slice.call and a Promise constructor. JS whitespace/comments are the
// library's own. Nothing here pretends to be a JS engine: it exists only so the three
// scripts above can meet each other the way they meet in a WebView2 document.
//
// FALSIFIABILITY (the tooth, ticket 35:150 "the ruler must redden the OLD shape too"):
// TestLegacySubShapeOneHookDiesInAReentryLoop feeds the hook that shipped in fb2fb802
// (kept below verbatim) into this same model and the model reports the loop - re-entry
// past the depth cap with the native exit called zero times. Deleting or no-op-ing
// installPanelTransport's w.Init(...) reddens every delivery assertion below (a mutant
// run is recorded in .scratch/wisp/probes/35/r2/), and so does removing either half of
// the ③ guard: drop "var native = cw.postMessage" (self-reference) and the cap trips;
// drop the try/finally reset and the second post dies at msgcb's unbound-method branch.
//
// SCOPE: AC#6's transport edge only. AC#7's (d)-1/-2/-3 guard-wording assertions
// (bridge.go:133/:136-137/:140-141) stay with a later leg, and this file's envelope is a
// roster method (bridge.go:132) so it exercises no reject branch. The real-window hop
// (-tags winlive) is 35-v2's, not this file's: no browser is started here.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// ---------------------------------------------------------------------------
// 1. the library's own scripts, taken from go-webview2's source text
// ---------------------------------------------------------------------------

// libExternalShim35r2 is pkg/edge/chromium.go:112 verbatim (the string the library hands
// to its own Init). Its body resolves chrome.webview.postMessage when invoke is CALLED,
// which is the property sub-shape ① overwrote - that is the whole loop.
const libExternalShim35r2 = `window.external={invoke:s=>window.chrome.webview.postMessage(s)}`

// libBindStubBody35r2 is the back-quoted half of the script webview.Bind injects
// (webview.go:462-478). The token sequence is the library's; only the leading tabs of the
// raw string literal are normalised to spaces, which JS does not read.
const libBindStubBody35r2 = `
  var RPC = window._rpc = (window._rpc || {nextSeq: 1});
  window[name] = function() {
    var seq = RPC.nextSeq++;
    var promise = new Promise(function(resolve, reject) {
      RPC[seq] = {
        resolve: resolve,
        reject: reject,
      };
    });
    window.external.invoke(JSON.stringify({
      id: seq,
      method: name,
      params: Array.prototype.slice.call(arguments),
    }));
    return promise;
  }
})()`

// libBindStubScript35r2 reproduces webview.go:462's concatenation, including its use of
// jsString (webview.go:137 = json.Marshal, HTML escaping left on exactly like the
// library), so the stub that runs in the model is the stub the library would inject.
func libBindStubScript35r2(name string) string {
	b, _ := json.Marshal(name)
	return "(function() { var name = " + string(b) + ";" + libBindStubBody35r2
}

// legacySubShapeOneHook is the forwarding script that shipped in fb2fb802
// (cmd/wisp/panel_host_windows.go:648 at that commit), copied verbatim. It is the
// counter-shape this ruler must be able to redden; see
// TestLegacySubShapeOneHookDiesInAReentryLoop.
const legacySubShapeOneHook = `(function () {
  var cw = window.chrome && window.chrome.webview;
  if (!cw || cw.__wispForwardInstalled) { return; }
  cw.__wispForwardInstalled = true;
  cw.postMessage = function (message) {
    if (typeof window.wispDispatch === "function") { window.wispDispatch(message); }
  };
})();`

// ---------------------------------------------------------------------------
// 2. the JavaScript subset
// ---------------------------------------------------------------------------

type jsValue = any

type jsPanic struct{ msg string }

func (p *jsPanic) Error() string { return p.msg }

type jsObject struct {
	props map[string]jsValue
	order []string
	// unwritable is 35-r4's property-writability table (ledger A684 §3, ticket 35:75 (a)
	// face M-B): a name marked here is a NON-WRITABLE data property. Before this existed
	// the fixture had no such concept at all - jsObject.set let any overwrite win - so the
	// shape "a browser/host that does not let chrome.webview.postMessage be replaced, and
	// a forwarding hook that therefore never takes effect" was structurally invisible.
	// Semantics modelled: assignment in SLOPPY mode is silently ignored (no TypeError),
	// the old value stays, and the refused write is recorded so a yard can tell "the hook
	// armed" apart from "the hook was silently unarmed". NOT modelled: configurable,
	// delete, accessor properties, and Object.getOwnPropertyDescriptor (see markNonWritable).
	unwritable map[string]bool
	// ignoredWrites names every assignment this object refused (one entry per attempt).
	ignoredWrites []string
	// onGet lets the browser model wrap one property's reads (chrome.webview.postMessage)
	// without changing what the scripts see. nil for ordinary objects.
	onGet func(name string, cur jsValue) jsValue
}

func newJSObject() *jsObject { return &jsObject{props: map[string]jsValue{}} }

func (o *jsObject) has(name string) bool { _, ok := o.props[name]; return ok }

// markNonWritable is the fixture-side stand-in for
// Object.defineProperty(o, name, {writable: false}). It is deliberately Go-side, not a JS
// Object method: ticket 35:75 (a) needs a WORLD shape a case can be run in, and a page
// that could descriptor-read its own transport would be a new claim this yard is not
// entitled to make (see the r4 section at the bottom of this file).
func (o *jsObject) markNonWritable(names ...string) {
	if o.unwritable == nil {
		o.unwritable = map[string]bool{}
	}
	for _, n := range names {
		o.unwritable[n] = true
	}
}

func (o *jsObject) set(name string, v jsValue) {
	if _, exists := o.props[name]; exists && o.unwritable[name] {
		// Non-writable + sloppy mode: the assignment is ignored, silently, exactly the way
		// a page cannot tell that its hook did not install.
		o.ignoredWrites = append(o.ignoredWrites, name)
		return
	}
	if _, ok := o.props[name]; !ok {
		o.order = append(o.order, name)
	}
	o.props[name] = v
}

type jsArray struct{ vals []jsValue }

type jsFunction struct {
	name   string
	host   func(recv jsValue, args []jsValue) jsValue
	params []string
	body   []jsStmt
	env    *jsEnv
}

type jsEnv struct {
	vars   map[string]jsValue
	parent *jsEnv
	win    *jsObject // the global scope reads and writes through window
	this   jsValue
}

func newJSScope(parent *jsEnv) *jsEnv { return &jsEnv{vars: map[string]jsValue{}, parent: parent} }

func (e *jsEnv) get(name string) (jsValue, bool) {
	for s := e; s != nil; s = s.parent {
		if v, ok := s.vars[name]; ok {
			return v, true
		}
		if s.win != nil {
			if v, ok := jsGetProp(s.win, name); ok {
				return v, true
			}
		}
	}
	return nil, false
}

func (e *jsEnv) set(name string, v jsValue) {
	for s := e; s != nil; s = s.parent {
		if _, ok := s.vars[name]; ok {
			s.vars[name] = v
			return
		}
		if s.win != nil {
			s.win.set(name, v)
			return
		}
	}
}

func (e *jsEnv) declare(name string, v jsValue) {
	if e.win != nil {
		e.win.set(name, v)
		return
	}
	e.vars[name] = v
}

func jsTypeOf(v jsValue) string {
	switch v.(type) {
	case nil:
		return "undefined"
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case *jsFunction:
		return "function"
	default:
		return "object"
	}
}

func jsTruthy(v jsValue) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}

func jsToNumber(v jsValue) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			panic(&jsPanic{"TypeError: cannot convert to number: " + x})
		}
		return f
	default:
		panic(&jsPanic{"TypeError: cannot convert " + jsTypeOf(v) + " to number"})
	}
}

func jsToString(v jsValue) string {
	switch x := v.(type) {
	case nil:
		return ""
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case string:
		return x
	default:
		return "object"
	}
}

func jsKey(v jsValue) string {
	if f, ok := v.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return jsToString(v)
}

func jsGetProp(o *jsObject, name string) (jsValue, bool) {
	v, ok := o.props[name]
	if !ok {
		return nil, false
	}
	if o.onGet != nil {
		v = o.onGet(name, v)
	}
	return v, true
}

// jsProp is the property read of JS: reading anything off a missing base throws, which is
// exactly the swallowed TypeError msgcb's reply script hits when window._rpc[id] does not
// exist (webview.go:157 against a page that never made a call).
func jsProp(v jsValue, name string) jsValue {
	switch o := v.(type) {
	case *jsObject:
		got, _ := jsGetProp(o, name)
		return got
	case *jsArray:
		if name == "length" {
			return float64(len(o.vals))
		}
		if i, err := strconv.Atoi(name); err == nil && i >= 0 && i < len(o.vals) {
			return o.vals[i]
		}
		return nil
	case *jsFunction:
		if name == "call" {
			f := o
			return &jsFunction{name: "call", host: func(_ jsValue, args []jsValue) jsValue {
				var recv jsValue
				rest := []jsValue{}
				if len(args) > 0 {
					recv = args[0]
					rest = args[1:]
				}
				return jsCallValue(f, recv, rest)
			}}
		}
		return nil
	case string:
		if name == "length" {
			return float64(len([]rune(o)))
		}
		return nil
	default:
		panic(&jsPanic{"TypeError: cannot read property '" + name + "' of " + jsTypeOf(v)})
	}
}

func jsElem(v jsValue, index jsValue) jsValue { return jsProp(v, jsKey(index)) }

func jsSetProp(v jsValue, name string, val jsValue) {
	switch o := v.(type) {
	case *jsObject:
		o.set(name, val)
	case *jsArray:
		if i, err := strconv.Atoi(name); err == nil && i >= 0 {
			for len(o.vals) <= i {
				o.vals = append(o.vals, nil)
			}
			o.vals[i] = val
			return
		}
		panic(&jsPanic{"TypeError: bad array index " + name})
	default:
		panic(&jsPanic{"TypeError: cannot set property '" + name + "' of " + jsTypeOf(v)})
	}
}

func jsCallValue(fn jsValue, recv jsValue, args []jsValue) jsValue {
	f, ok := fn.(*jsFunction)
	if !ok {
		panic(&jsPanic{"TypeError: " + jsTypeOf(fn) + " is not a function"})
	}
	if f.host != nil {
		return f.host(recv, args)
	}
	scope := newJSScope(f.env)
	scope.this = recv
	for i, p := range f.params {
		var v jsValue
		if i < len(args) {
			v = args[i]
		}
		scope.vars[p] = v
	}
	scope.vars["arguments"] = &jsArray{vals: append([]jsValue(nil), args...)}
	ret, _ := jsRunStmts(f.body, scope)
	return ret
}

// nodes
type jsNode = any

type jsIdent struct{ name string }
type jsLit struct{ v jsValue }
type jsMember struct {
	obj  jsNode
	name string
}
type jsIndex struct{ obj, index jsNode }
type jsCallNode struct {
	callee jsNode
	args   []jsNode
	isNew  bool
}
type jsUnary struct {
	op string
	x  jsNode
}
type jsBinary struct {
	op string
	x  jsNode
	y  jsNode
}
type jsAssignNode struct{ target, value jsNode }
type jsPostInc struct{ target jsNode }
type jsFuncExpr struct {
	params []string
	body   []jsStmt
}
type jsObjLit struct {
	keys []string
	vals []jsNode
}
type jsArrLit struct{ vals []jsNode }

type jsVarDecl struct {
	name string
	init jsNode
}

type jsStmt struct {
	kind        string // expr | var | if | return | block | try
	expr        jsNode
	decls       []jsVarDecl
	cond        jsNode
	thenBody    []jsStmt
	elseBody    []jsStmt
	body        []jsStmt
	finallyBody []jsStmt
}

func jsEval(n jsNode, env *jsEnv) jsValue {
	switch e := n.(type) {
	case *jsLit:
		return e.v
	case *jsIdent:
		if v, ok := env.get(e.name); ok {
			return v
		}
		if e.name == "undefined" || e.name == "NaN" {
			return nil
		}
		panic(&jsPanic{"ReferenceError: " + e.name + " is not defined"})
	case *jsMember:
		return jsProp(jsEval(e.obj, env), e.name)
	case *jsIndex:
		return jsElem(jsEval(e.obj, env), jsEval(e.index, env))
	case *jsFuncExpr:
		return &jsFunction{params: e.params, body: e.body, env: env}
	case *jsObjLit:
		o := newJSObject()
		for i, k := range e.keys {
			o.set(k, jsEval(e.vals[i], env))
		}
		return o
	case *jsArrLit:
		a := &jsArray{}
		for _, v := range e.vals {
			a.vals = append(a.vals, jsEval(v, env))
		}
		return a
	case *jsUnary:
		if e.op == "typeof" {
			if id, ok := e.x.(*jsIdent); ok {
				v, _ := env.get(id.name)
				return jsTypeOf(v)
			}
			return jsTypeOf(jsEval(e.x, env))
		}
		v := jsEval(e.x, env)
		switch e.op {
		case "!":
			return !jsTruthy(v)
		case "-":
			return -jsToNumber(v)
		}
		panic(&jsPanic{"internal: unary " + e.op})
	case *jsBinary:
		return jsBinaryEval(e, env)
	case *jsPostInc:
		cur := jsToNumber(jsEval(e.target, env))
		switch t := e.target.(type) {
		case *jsIdent:
			env.set(t.name, cur+1)
		case *jsMember:
			jsSetProp(jsEval(t.obj, env), t.name, cur+1)
		case *jsIndex:
			jsSetProp(jsEval(t.obj, env), jsKey(jsEval(t.index, env)), cur+1)
		default:
			panic(&jsPanic{"internal: bad ++ target"})
		}
		return cur
	case *jsAssignNode:
		v := jsEval(e.value, env)
		switch t := e.target.(type) {
		case *jsIdent:
			env.set(t.name, v)
		case *jsMember:
			jsSetProp(jsEval(t.obj, env), t.name, v)
		case *jsIndex:
			jsSetProp(jsEval(t.obj, env), jsKey(jsEval(t.index, env)), v)
		default:
			panic(&jsPanic{"internal: bad assignment target"})
		}
		return v
	case *jsCallNode:
		return jsEvalCall(e, env)
	}
	panic(&jsPanic{"internal: unhandled expression"})
}

func jsBinaryEval(e *jsBinary, env *jsEnv) jsValue {
	switch e.op {
	case "&&":
		v := jsEval(e.x, env)
		if !jsTruthy(v) {
			return v
		}
		return jsEval(e.y, env)
	case "||":
		v := jsEval(e.x, env)
		if jsTruthy(v) {
			return v
		}
		return jsEval(e.y, env)
	case "===", "!=":
		v := jsEval(e.x, env)
		w := jsEval(e.y, env)
		eq := jsStrictEqual(v, w)
		if e.op == "===" {
			return eq
		}
		return !eq
	case "==", "!==":
		v := jsEval(e.x, env)
		w := jsEval(e.y, env)
		// This model has no coercion rules beyond "same kind, same value", which is all
		// the three scripts above ever ask of ==.
		eq := jsStrictEqual(v, w)
		if e.op == "==" {
			return eq
		}
		return !eq
	case "+":
		v := jsEval(e.x, env)
		w := jsEval(e.y, env)
		if jsTypeOf(v) == "string" || jsTypeOf(w) == "string" {
			return jsToString(v) + jsToString(w)
		}
		return jsToNumber(v) + jsToNumber(w)
	case "-":
		return jsToNumber(jsEval(e.x, env)) - jsToNumber(jsEval(e.y, env))
	}
	panic(&jsPanic{"internal: binary " + e.op})
}

func jsStrictEqual(v, w jsValue) bool {
	if jsTypeOf(v) != jsTypeOf(w) {
		return false
	}
	switch a := v.(type) {
	case nil:
		return true
	case *jsFunction, *jsObject, *jsArray:
		return v == w
	case string:
		return a == w.(string)
	case float64:
		return a == w.(float64)
	case bool:
		return a == w.(bool)
	}
	return false
}

func jsEvalCall(e *jsCallNode, env *jsEnv) jsValue {
	var recv jsValue
	var fn jsValue
	if m, ok := e.callee.(*jsMember); ok {
		recv = jsEval(m.obj, env)
		fn = jsProp(recv, m.name)
	} else if ix, ok := e.callee.(*jsIndex); ok {
		recv = jsEval(ix.obj, env)
		fn = jsElem(recv, jsEval(ix.index, env))
	} else {
		fn = jsEval(e.callee, env)
	}
	args := make([]jsValue, 0, len(e.args))
	for _, a := range e.args {
		args = append(args, jsEval(a, env))
	}
	if e.isNew {
		created := newJSObject()
		out := jsCallValue(fn, created, args)
		if o, ok := out.(*jsObject); ok {
			return o
		}
		return created
	}
	return jsCallValue(fn, recv, args)
}

func jsRunStmts(body []jsStmt, env *jsEnv) (jsValue, bool) {
	var last jsValue
	for _, s := range body {
		switch s.kind {
		case "expr":
			last = jsEval(s.expr, env)
		case "var":
			for _, d := range s.decls {
				var v jsValue
				if d.init != nil {
					v = jsEval(d.init, env)
				}
				env.declare(d.name, v)
			}
		case "if":
			if jsTruthy(jsEval(s.cond, env)) {
				r, ret := jsRunStmts(s.thenBody, env)
				if ret {
					return r, true
				}
				last = r
			} else if s.elseBody != nil {
				r, ret := jsRunStmts(s.elseBody, env)
				if ret {
					return r, true
				}
				last = r
			}
		case "block":
			scope := newJSScope(env)
			r, ret := jsRunStmts(s.body, scope)
			if ret {
				return r, true
			}
			last = r
		case "return":
			if s.expr == nil {
				return nil, true
			}
			return jsEval(s.expr, env), true
		case "try":
			var ret jsValue
			var returned bool
			var pan any
			func() {
				defer func() { pan = recover() }()
				ret, returned = jsRunStmts(s.body, env)
			}()
			if s.finallyBody != nil {
				fr, frRet := jsRunStmts(s.finallyBody, env)
				if pan != nil {
					// the finally block running is all a browser does before the throw
					// continues; a value it computed is discarded.
					_ = fr
					_ = frRet
				} else if frRet {
					return fr, true
				}
			}
			if pan != nil {
				panic(pan)
			}
			return ret, returned
		default:
			panic(&jsPanic{"internal: unhandled statement " + s.kind})
		}
	}
	return last, false
}

// --- lexer ---------------------------------------------------------------

type jsTok struct {
	kind string // id | num | str | punct
	text string
	num  float64
}

var jsPuncts = []string{"===", "!==", "=>", "==", "!=", "&&", "||", "++",
	"(", ")", "{", "}", "[", "]", ";", ",", ".", "=", "!", "+", "-", ":"}

func jsLex(src string) []jsTok {
	var out []jsTok
	i := 0
	for i < len(src) {
		c := src[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}
		if c == '/' && i+1 < len(src) && src[i+1] == '/' {
			for i < len(src) && src[i] != '\n' {
				i++
			}
			continue
		}
		if c == '/' && i+1 < len(src) && src[i+1] == '*' {
			end := bytes.Index([]byte(src[i:]), []byte("*/"))
			if end < 0 {
				i = len(src)
			} else {
				i += end + 2
			}
			continue
		}
		if c == '"' || c == '\'' {
			quote := c
			i++
			var b []byte
			for i < len(src) && src[i] != quote {
				if src[i] == '\\' && i+1 < len(src) {
					i++
					switch src[i] {
					case 'n':
						b = append(b, '\n')
					case 't':
						b = append(b, '\t')
					default:
						b = append(b, src[i])
					}
					i++
					continue
				}
				b = append(b, src[i])
				i++
			}
			i++ // closing quote
			out = append(out, jsTok{kind: "str", text: string(b)})
			continue
		}
		if c >= '0' && c <= '9' {
			j := i
			for j < len(src) && ((src[j] >= '0' && src[j] <= '9') || src[j] == '.') {
				j++
			}
			f, _ := strconv.ParseFloat(src[i:j], 64)
			out = append(out, jsTok{kind: "num", text: src[i:j], num: f})
			i = j
			continue
		}
		if isJsIDStart(c) {
			j := i
			for j < len(src) && isJsIDByte(src[j]) {
				j++
			}
			out = append(out, jsTok{kind: "id", text: src[i:j]})
			i = j
			continue
		}
		matched := false
		for _, p := range jsPuncts {
			if stringsHasPrefixAt(src, i, p) {
				out = append(out, jsTok{kind: "punct", text: p})
				i += len(p)
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		panic(&jsPanic{"SyntaxError: unexpected character " + string(rune(c))})
	}
	out = append(out, jsTok{kind: "eof"})
	return out
}

func isJsIDStart(c byte) bool {
	return c == '_' || c == '$' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isJsIDByte(c byte) bool { return isJsIDStart(c) || (c >= '0' && c <= '9') }

func stringsHasPrefixAt(s string, i int, p string) bool {
	return i+len(p) <= len(s) && s[i:i+len(p)] == p
}

// --- parser --------------------------------------------------------------

type jsParser struct {
	toks []jsTok
	i    int
}

func (p *jsParser) peek() jsTok { return p.toks[p.i] }
func (p *jsParser) next() jsTok { t := p.toks[p.i]; p.i++; return t }
func (p *jsParser) atPunct(s string) bool {
	t := p.peek()
	return t.kind == "punct" && t.text == s
}
func (p *jsParser) atID(s string) bool {
	t := p.peek()
	return t.kind == "id" && t.text == s
}
func (p *jsParser) eatPunct(s string) bool {
	if p.atPunct(s) {
		p.i++
		return true
	}
	return false
}
func (p *jsParser) expectPunct(s string) {
	if !p.eatPunct(s) {
		panic(&jsPanic{"SyntaxError: expected " + s + ", saw " + p.peek().text})
	}
}

func jsParseProgram(src string) []jsStmt {
	p := &jsParser{toks: jsLex(src)}
	return p.stmts()
}

func (p *jsParser) stmts() []jsStmt {
	var out []jsStmt
	for {
		t := p.peek()
		if t.kind == "eof" || p.atPunct("}") {
			return out
		}
		out = append(out, p.stmt())
	}
}

func (p *jsParser) block() []jsStmt {
	p.expectPunct("{")
	body := p.stmts()
	p.expectPunct("}")
	return body
}

func (p *jsParser) optSemi() { p.eatPunct(";") }

func (p *jsParser) stmt() jsStmt {
	switch {
	case p.atID("var"):
		p.next()
		var decls []jsVarDecl
		for {
			name := p.next().text
			var init jsNode
			if p.eatPunct("=") {
				init = p.expr()
			}
			decls = append(decls, jsVarDecl{name: name, init: init})
			if !p.eatPunct(",") {
				break
			}
		}
		p.optSemi()
		return jsStmt{kind: "var", decls: decls}
	case p.atID("if"):
		p.next()
		p.expectPunct("(")
		cond := p.expr()
		p.expectPunct(")")
		thenBody := p.oneStmt()
		var elseBody []jsStmt
		if p.atID("else") {
			p.next()
			elseBody = p.oneStmt()
		}
		return jsStmt{kind: "if", cond: cond, thenBody: thenBody, elseBody: elseBody}
	case p.atID("return"):
		p.next()
		var e jsNode
		if !p.atPunct(";") && !p.atPunct("}") && p.peek().kind != "eof" {
			e = p.expr()
		}
		p.optSemi()
		return jsStmt{kind: "return", expr: e}
	case p.atID("try"):
		p.next()
		body := p.block()
		var fin []jsStmt
		if p.atID("finally") {
			p.next()
			fin = p.block()
		}
		return jsStmt{kind: "try", body: body, finallyBody: fin}
	case p.atPunct("{"):
		return jsStmt{kind: "block", body: p.block()}
	default:
		e := p.expr()
		p.optSemi()
		return jsStmt{kind: "expr", expr: e}
	}
}

func (p *jsParser) oneStmt() []jsStmt {
	if p.atPunct("{") {
		return p.block()
	}
	return []jsStmt{p.stmt()}
}

func (p *jsParser) expr() jsNode {
	left := p.or()
	if p.eatPunct("=") {
		return &jsAssignNode{target: left, value: p.expr()}
	}
	return left
}

func (p *jsParser) or() jsNode {
	left := p.and()
	for p.eatPunct("||") {
		left = &jsBinary{op: "||", x: left, y: p.and()}
	}
	return left
}

func (p *jsParser) and() jsNode {
	left := p.equality()
	for p.eatPunct("&&") {
		left = &jsBinary{op: "&&", x: left, y: p.equality()}
	}
	return left
}

func (p *jsParser) equality() jsNode {
	left := p.additive()
	for {
		for _, op := range []string{"===", "!==", "==", "!="} {
			if p.atPunct(op) {
				p.next()
				left = &jsBinary{op: op, x: left, y: p.additive()}
			}
		}
		return left
	}
}

func (p *jsParser) additive() jsNode {
	left := p.unary()
	for {
		if p.atPunct("+") || p.atPunct("-") {
			op := p.next().text
			left = &jsBinary{op: op, x: left, y: p.unary()}
			continue
		}
		return left
	}
}

func (p *jsParser) unary() jsNode {
	if p.eatPunct("!") {
		return &jsUnary{op: "!", x: p.unary()}
	}
	if p.eatPunct("-") {
		return &jsUnary{op: "-", x: p.unary()}
	}
	if p.atID("typeof") {
		p.next()
		return &jsUnary{op: "typeof", x: p.unary()}
	}
	return p.postfix()
}

func (p *jsParser) postfix() jsNode {
	n := p.primary()
	for {
		switch {
		case p.eatPunct("."):
			n = &jsMember{obj: n, name: p.next().text}
		case p.eatPunct("["):
			idx := p.expr()
			p.expectPunct("]")
			n = &jsIndex{obj: n, index: idx}
		case p.atPunct("("):
			n = &jsCallNode{callee: n, args: p.args()}
		case p.eatPunct("++"):
			n = &jsPostInc{target: n}
		default:
			return n
		}
	}
}

func (p *jsParser) args() []jsNode {
	p.expectPunct("(")
	var out []jsNode
	for !p.atPunct(")") {
		out = append(out, p.expr())
		if !p.eatPunct(",") {
			break
		}
	}
	p.expectPunct(")")
	return out
}

func (p *jsParser) primary() jsNode {
	t := p.peek()
	switch t.kind {
	case "num":
		p.next()
		return &jsLit{v: t.num}
	case "str":
		p.next()
		return &jsLit{v: t.text}
	case "id":
		switch t.text {
		case "true":
			p.next()
			return &jsLit{v: true}
		case "false":
			p.next()
			return &jsLit{v: false}
		case "null":
			p.next()
			return &jsLit{v: nil}
		case "function":
			p.next()
			if p.peek().kind == "id" {
				p.next() // named function expression: the name is not needed here
			}
			return &jsFuncExpr{params: p.params(), body: p.block()}
		case "new":
			p.next()
			callee := p.memberOnly()
			var args []jsNode
			if p.atPunct("(") {
				args = p.args()
			}
			return &jsCallNode{callee: callee, args: args, isNew: true}
		}
		// one-argument arrow: s=>expr
		if p.toks[p.i+1].kind == "punct" && p.toks[p.i+1].text == "=>" {
			p.next()
			p.next()
			return &jsFuncExpr{params: []string{t.text}, body: []jsStmt{{kind: "return", expr: p.expr()}}}
		}
		p.next()
		return &jsIdent{name: t.text}
	case "punct":
		switch t.text {
		case "(":
			p.next()
			if params, ok := p.tryArrowParams(); ok {
				return &jsFuncExpr{params: params, body: []jsStmt{{kind: "return", expr: p.expr()}}}
			}
			e := p.expr()
			p.expectPunct(")")
			return e
		case "{":
			o := &jsObjLit{}
			p.next()
			for !p.atPunct("}") {
				k := p.next()
				key := k.text
				if k.kind != "str" && k.kind != "id" {
					panic(&jsPanic{"SyntaxError: bad object key " + key})
				}
				p.expectPunct(":")
				o.keys = append(o.keys, key)
				o.vals = append(o.vals, p.expr())
				if !p.eatPunct(",") {
					break
				}
			}
			p.expectPunct("}")
			return o
		case "[":
			p.next()
			a := &jsArrLit{}
			for !p.atPunct("]") {
				a.vals = append(a.vals, p.expr())
				if !p.eatPunct(",") {
					break
				}
			}
			p.expectPunct("]")
			return a
		}
	}
	panic(&jsPanic{"SyntaxError: unexpected token " + t.text + " (" + t.kind + ")"})
}

// memberOnly parses a.b.c without allowing a call, for `new Promise(...)`.
func (p *jsParser) memberOnly() jsNode {
	n := p.primary()
	for {
		switch {
		case p.eatPunct("."):
			n = &jsMember{obj: n, name: p.next().text}
		case p.eatPunct("["):
			idx := p.expr()
			p.expectPunct("]")
			n = &jsIndex{obj: n, index: idx}
		default:
			return n
		}
	}
}

// tryArrowParams recognises () => and (a, b) => at the current position (just past the
// opening paren) and consumes through the =>. None of the three scripts this file runs
// needs the parenthesised form - the library's shim uses a bare s=>expr and its stub uses
// function(){} - but the subset lists it, so the parser says so.
func (p *jsParser) tryArrowParams() ([]string, bool) {
	j := p.i
	var out []string
	for {
		if j >= len(p.toks) {
			return nil, false
		}
		if p.toks[j].kind == "punct" && p.toks[j].text == ")" {
			break
		}
		if p.toks[j].kind != "id" {
			return nil, false
		}
		out = append(out, p.toks[j].text)
		j++
		if j >= len(p.toks) {
			return nil, false
		}
		if p.toks[j].kind == "punct" && p.toks[j].text == "," {
			j++
			continue
		}
		if p.toks[j].kind == "punct" && p.toks[j].text == ")" {
			break
		}
		return nil, false
	}
	if j+1 >= len(p.toks) {
		return nil, false
	}
	if p.toks[j+1].kind != "punct" || p.toks[j+1].text != "=>" {
		return nil, false
	}
	p.i = j + 2
	return out, true
}

func (p *jsParser) params() []string {
	p.expectPunct("(")
	var out []string
	for !p.atPunct(")") {
		if p.peek().kind == "id" {
			out = append(out, p.next().text)
		} else {
			p.next()
		}
		if !p.eatPunct(",") {
			break
		}
	}
	p.expectPunct(")")
	return out
}

// ---------------------------------------------------------------------------
// 3. the behavioural fake: a document, the library's pipe, and msgcb
// ---------------------------------------------------------------------------

// postMessageHopCap stands in for the browser's call-stack limit. 35-v1 measured the
// landed sub-shape ① as one that never reaches the native exit; a real page dies with
// RangeError: Maximum call stack size exceeded after a few thousand frames, so any small
// cap shows the same thing. It is deliberately a cap, not a timeout (AGENTS.md §1.2).
const postMessageHopCap = 8

type fakeDoc35r2 struct {
	initScripts []string // AddScriptToExecuteOnDocumentCreated, in registration order
	bindings    map[string]interface{}

	window  *jsObject
	webview *jsObject
	global  *jsEnv

	// readings the assertions are made of
	nativeCalls    int
	nativeFrames   []string
	doorRounds     int
	lastReply      string
	postDepth      int
	maxPostDepth   int
	capTripped     bool
	capDepth       int
	capMessage     string
	unparsedFrames int
	deadSlots      int
	evalScripts    []string
	evalThrew      int
	resolved       int
	rejected       int
	documents      int
	scriptThrows   []string

	// 35-r4 readings, for the two faces A684 §3 found this fixture blind to.
	// nativeExitNonWritable asks openDocument to expose chrome.webview.postMessage as a
	// NON-WRITABLE property (the world M-B needed); nativeExitReceiverOK counts the native
	// exit invocations that arrived with chrome.webview as their receiver and
	// nativeReceiverMismatch counts those that lost it (the world M-A needed).
	nativeExitNonWritable  bool
	nativeExitReceiverOK   int
	nativeReceiverMismatch int
}

func newFakeDoc35r2() *fakeDoc35r2 {
	return &fakeDoc35r2{bindings: map[string]interface{}{}}
}

// Bind mirrors webview.Bind: the callback goes into the bindings map and the library's
// stub script is registered through the same Init path (webview.go:462 calls w.Init).
func (bf *fakeDoc35r2) Bind(name string, f interface{}) error {
	bf.bindings[name] = f
	bf.registerInit(libBindStubScript35r2(name))
	return nil
}

// Init mirrors webview.Init (webview.go:435 -> chromium.Init ->
// AddScriptToExecuteOnDocumentCreated, chromium.go:130-136): it queues a script for every
// document created after this point.
func (bf *fakeDoc35r2) Init(js string) { bf.registerInit(js) }

func (bf *fakeDoc35r2) registerInit(js string) { bf.initScripts = append(bf.initScripts, js) }

// openDocument builds a fresh window exactly the way WebView2 hands a new document to the
// scripts: chrome.webview with its NATIVE postMessage first, then the library's own
// external.invoke shim (registered at chromium.go:112 while the control is created, so
// ahead of anything the host queued afterwards), then the host's queued scripts.
func (bf *fakeDoc35r2) openDocument() {
	bf.documents++
	win := newJSObject()
	cw := newJSObject()
	cw.onGet = func(name string, cur jsValue) jsValue {
		if name != "postMessage" {
			return cur
		}
		// A read of chrome.webview.postMessage hands back a hop that counts ITSELF at
		// call time and then calls whatever the property held at READ time. That is the
		// honest shape: var native = cw.postMessage captures the native exit, a later read
		// sees the overwritten value, and a wrapper that calls the property again is
		// counted as the re-entry a real browser would run out of stack on.
		captured := cur
		return &jsFunction{name: "postMessageHop", host: func(recv jsValue, args []jsValue) jsValue {
			bf.postDepth++
			if bf.postDepth > bf.maxPostDepth {
				bf.maxPostDepth = bf.postDepth
			}
			if bf.postDepth > postMessageHopCap && !bf.capTripped {
				bf.capTripped = true
				bf.capDepth = bf.postDepth
				bf.capMessage = fmt.Sprintf(
					"RangeError: maximum call stack size exceeded - chrome.webview.postMessage re-entered %d levels (cap %d) with the native exit called %d times; a page in a real browser dies here",
					bf.postDepth, postMessageHopCap, bf.nativeCalls)
				panic(&jsPanic{bf.capMessage})
			}
			defer func() { bf.postDepth-- }()
			return jsCallValue(captured, recv, args)
		}}
	}
	// the native exit: what chrome.webview.postMessage actually is in WebView2 - the one
	// channel into the Go message callback.
	//
	// 35-r4, face M-A (ledger A684 §3, ticket 35:75 (a)). This stub used to be
	// func(_ jsValue, args []jsValue): it DISCARDED its receiver, so a forwarding hook
	// written as a bare native(message) behaved here exactly like the shipped
	// native.call(cw, message) - 35-v2's one-line mutant stayed green on all five cases.
	// A browser's host method is not a detached function: Chrome and WebView2 are reported
	// to answer an invocation that lost its receiver with "TypeError: Illegal invocation".
	// THIS RULER CANNOT PRODUCE THAT CREDENTIAL (ticket 35 AC#8(iii)) - it is a hand-written
	// interpreter, so the behaviour recorded below is THIS FIXTURE'S OWN ANSWER, modelled
	// after Chrome/WebView2 and never offered as a WebView2 behaviour record: the stub
	// requires cw itself (identity, not type) as the receiver and says so by name. The only
	// real-window reading of this edge lives in
	// cmd/wisp/panel_transport_live_35v2_windows_test.go, which is 〔仅本机可量〕.
	cw.set("postMessage", &jsFunction{name: "native", host: func(recv jsValue, args []jsValue) jsValue {
		if recv != cw {
			bf.nativeReceiverMismatch++
			panic(&jsPanic{"TypeError: Illegal invocation: chrome.webview.postMessage lost its receiver - called with " +
				jsTypeOf(recv) + ", want the chrome.webview object the page holds"})
		}
		bf.nativeExitReceiverOK++
		bf.nativeCalls++
		msg := ""
		if len(args) > 0 {
			msg = jsToString(args[0])
		}
		bf.nativeFrames = append(bf.nativeFrames, msg)
		bf.msgcb(msg)
		return nil
	}})
	// 35-r4, face M-B: the台件 that makes the native exit non-writable, i.e. a host that
	// does not let the page replace its own transport. The hook's assignment is then
	// silently ignored and the forwarding never arms - the shape the fixture could not
	// express before this line.
	if bf.nativeExitNonWritable {
		cw.markNonWritable("postMessage")
	}
	chrome := newJSObject()
	chrome.set("webview", cw)
	win.set("chrome", chrome)
	win.set("window", win) // window.window === window, as in a browser
	win.set("JSON", jsonGlobal35r2())
	win.set("Array", arrayGlobal35r2())
	win.set("Promise", promiseGlobal35r2(bf))

	bf.window = win
	bf.webview = cw

	global := newJSScope(nil)
	global.win = win
	global.this = win
	bf.global = global

	bf.runAndNote(global, libExternalShim35r2)
	for _, s := range bf.initScripts {
		bf.runAndNote(global, s)
	}
}

// runAndNote is openDocument's script hop: a document script that throws is exactly what a
// browser swallows, so the model records it rather than losing it - the assertions below
// name it in their failure text.
func (bf *fakeDoc35r2) runAndNote(env *jsEnv, src string) {
	_, thrown, _ := bf.runScriptOn(env, src)
	if thrown != "" {
		bf.scriptThrows = append(bf.scriptThrows, firstLine35r2(src)+": "+thrown)
	}
}

func firstLine35r2(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	if len(s) > 40 {
		return s[:40]
	}
	return s
}

// evalOnWindow runs one test-side script as page code (the page's own
// chrome.webview.postMessage(...) call and msgcb's reply Eval both go this way).
func (bf *fakeDoc35r2) runScriptOn(env *jsEnv, src string) (out jsValue, thrown string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if p, isJS := r.(*jsPanic); isJS {
				thrown = p.msg
				return
			}
			panic(r)
		}
	}()
	return jsRunFirst(jsParseProgram(src), env), "", true
}

func jsRunFirst(body []jsStmt, env *jsEnv) jsValue {
	v, _ := jsRunStmts(body, env)
	return v
}

// pagePost is the page's own hop: frontend/src/lib/panel.ts:179/:218 call
// chrome.webview.postMessage(envelope). The envelope is injected as a window property so
// the posted string is data, never code.
func (bf *fakeDoc35r2) pagePost(envelope string) (thrown string) {
	bf.window.set("__wispTestPageEnvelope", envelope)
	_, thrown, _ = bf.runScriptOn(bf.global, "window.chrome.webview.postMessage(window.__wispTestPageEnvelope);")
	return thrown
}

// msgcb is go-webview2's message callback reproduced against webview.go:131-168: the same
// rpcMessage shape, the same callbinding answer for an unbound name (nil, nil - NOT an
// error), and therefore the same reply script text at :157, window._rpc[id].resolve(...).
// Ticket 35:58-59's "reject" wording is stale; ledger A682 §2 corrected it, and this
// model follows the library, not the ticket sentence.
func (bf *fakeDoc35r2) msgcb(msg string) {
	var d struct {
		ID     int               `json:"id"`
		Method string            `json:"method"`
		Params []json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal([]byte(msg), &d); err != nil {
		bf.unparsedFrames++ // webview.go:141-143: log.Printf + return, the door is not called
		return
	}
	id := strconv.Itoa(d.ID)
	res, err := bf.callbinding(d.Method, d.Params)
	if err != nil {
		// webview.go:144-151: the reject branch, for a bad argument shape only.
		bf.eval(fmt.Sprintf("window._rpc[%s].reject(%s); window._rpc[%s] = undefined", id, jsStringForLib35r2(err.Error()), id))
		return
	}
	// webview.go:155-159: the success branch. An UNBOUND method arrives here too, because
	// callbinding answers (nil, nil) - which is why the resolve target is _rpc[0] and the
	// reply is null, exactly the ticket 35:59 death and not the reject the ticket text says.
	bf.eval(fmt.Sprintf("window._rpc[%s].resolve(%s); window._rpc[%s] = undefined", id, mustMarshal35r2(res), id))
}

// callbinding is webview.callbinding for the one door shape the panel binds,
// func(string) string: unbound name -> (nil, nil); wrong arity -> error (the reject
// branch); otherwise call it.
func (bf *fakeDoc35r2) callbinding(method string, params []json.RawMessage) (interface{}, error) {
	f, ok := bf.bindings[method]
	if !ok {
		bf.deadSlots++
		return nil, nil
	}
	door, ok := f.(func(string) string)
	if !ok {
		return nil, fmt.Errorf("fake: bound %q is not func(string) string", method)
	}
	if len(params) != 1 {
		return nil, fmt.Errorf("function arguments mismatch")
	}
	var raw string
	if err := json.Unmarshal(params[0], &raw); err != nil {
		return nil, err
	}
	bf.doorRounds++
	bf.lastReply = door(raw)
	return bf.lastReply, nil
}

// eval is webview.Eval for the reply scripts: it runs as page code, and an uncaught error
// is what the library never sees (webview.go:435-437 ignores it; the COM call's result is
// discarded). window._rpc[0] does not exist for a page that made no bound call, so the
// ticket 35:59 death shows up here as a swallowed TypeError.
func (bf *fakeDoc35r2) eval(script string) {
	bf.evalScripts = append(bf.evalScripts, script)
	_, thrown, _ := bf.runScriptOn(bf.global, script)
	if thrown != "" {
		bf.evalThrew++
	}
}

func (bf *fakeDoc35r2) summary() string {
	ignored := "(none)"
	if bf.webview != nil && len(bf.webview.ignoredWrites) > 0 {
		ignored = fmt.Sprint(bf.webview.ignoredWrites)
	}
	return fmt.Sprintf("nativeExitCalls=%d receiverOK=%d receiverLost=%d ignoredWrites=%s doorRounds=%d maxPostDepth=%d capTripped=%v capDepth=%d unparsable=%d unboundSlots=%d evalThrew=%d resolved=%d rejected=%d scriptThrows=%v",
		bf.nativeCalls, bf.nativeExitReceiverOK, bf.nativeReceiverMismatch, ignored,
		bf.doorRounds, bf.maxPostDepth, bf.capTripped, bf.capDepth,
		bf.unparsedFrames, bf.deadSlots, bf.evalThrew, bf.resolved, bf.rejected, bf.scriptThrows)
}

// jsStringForLib35r2 is webview.go:137 jsString verbatim.
func jsStringForLib35r2(v interface{}) string { b, _ := json.Marshal(v); return string(b) }

func mustMarshal35r2(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "null"
	}
	return string(b)
}

// jsonGlobal35r2 gives the scripts JSON.stringify. Escaping is turned off so the bytes
// match what a browser produces; the envelope the page posts needs no escaping, so this
// cannot change any assertion below - it only keeps the model from inventing \u003c forms.
func jsonGlobal35r2() *jsObject {
	o := newJSObject()
	o.set("stringify", &jsFunction{name: "JSON.stringify", host: func(_ jsValue, args []jsValue) jsValue {
		if len(args) == 0 {
			return nil
		}
		return jsStringify35r2(args[0])
	}})
	o.set("parse", &jsFunction{name: "JSON.parse", host: func(_ jsValue, args []jsValue) jsValue {
		if len(args) == 0 {
			panic(&jsPanic{"SyntaxError: JSON.parse requires a string"})
		}
		return jsParseToValue35r2(jsToString(args[0]))
	}})
	return o
}

func jsStringify35r2(v jsValue) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case *jsObject:
		var b bytes.Buffer
		b.WriteString("{")
		for i, k := range x.order {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(jsQuote35r2(k))
			b.WriteString(":")
			b.WriteString(jsStringify35r2(x.props[k]))
		}
		b.WriteString("}")
		return b.String()
	case *jsArray:
		var b bytes.Buffer
		b.WriteString("[")
		for i, e := range x.vals {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(jsStringify35r2(e))
		}
		b.WriteString("]")
		return b.String()
	case string:
		return jsQuote35r2(x)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	}
	return "null"
}

func jsQuote35r2(s string) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	out := b.String()
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out
}

func jsParseToValue35r2(src string) jsValue {
	var raw json.RawMessage = json.RawMessage(src)
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return decodeJS35r2(dec)
}

func decodeJS35r2(dec *json.Decoder) jsValue {
	t, err := dec.Token()
	if err != nil {
		panic(&jsPanic{"SyntaxError: Unexpected token in JSON"})
	}
	switch v := t.(type) {
	case json.Delim:
		if v == '[' {
			a := &jsArray{}
			for dec.More() {
				a.vals = append(a.vals, decodeJS35r2(dec))
			}
			_, _ = dec.Token()
			return a
		}
		o := newJSObject()
		for dec.More() {
			keyTok, _ := dec.Token()
			key := keyTok.(string)
			o.set(key, decodeJS35r2(dec))
		}
		_, _ = dec.Token()
		return o
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f
		}
		return v.String()
	case string:
		return v
	case bool:
		return v
	case nil:
		return nil
	}
	panic(&jsPanic{"SyntaxError: bad JSON"})
}

func arrayGlobal35r2() *jsObject {
	proto := newJSObject()
	proto.set("slice", &jsFunction{name: "slice", host: func(recv jsValue, args []jsValue) jsValue {
		a, ok := recv.(*jsArray)
		if !ok {
			panic(&jsPanic{"TypeError: slice of non-array"})
		}
		from := 0
		if len(args) > 0 {
			from = int(jsToNumber(args[0]))
		}
		to := len(a.vals)
		if len(args) > 1 {
			to = int(jsToNumber(args[1]))
		}
		if from < 0 {
			from = 0
		}
		if to > len(a.vals) {
			to = len(a.vals)
		}
		out := &jsArray{}
		for i := from; i < to; i++ {
			out.vals = append(out.vals, a.vals[i])
		}
		return out
	}})
	o := newJSObject()
	o.set("prototype", proto)
	o.set("isArray", &jsFunction{name: "isArray", host: func(_ jsValue, args []jsValue) jsValue {
		if len(args) == 0 {
			return false
		}
		_, isArr := args[0].(*jsArray)
		return isArr
	}})
	return o
}

// promiseGlobal35r2 is the executor-shaped Promise the library stub uses: new
// Promise(function(resolve, reject){...}) hands back an object whose resolve/reject are
// recorded on the fake, which is what lets the test assert the page really got a receipt.
func promiseGlobal35r2(bf *fakeDoc35r2) *jsFunction {
	return &jsFunction{name: "Promise", host: func(_ jsValue, args []jsValue) jsValue {
		p := newJSObject()
		resolve := &jsFunction{name: "resolve", host: func(_ jsValue, rargs []jsValue) jsValue {
			bf.resolved++
			return nil
		}}
		reject := &jsFunction{name: "reject", host: func(_ jsValue, rargs []jsValue) jsValue {
			bf.rejected++
			return nil
		}}
		p.set("resolve", resolve)
		p.set("reject", reject)
		if len(args) > 0 {
			jsCallValue(args[0], nil, []jsValue{resolve, reject})
		}
		return p
	}}
}

// ---------------------------------------------------------------------------
// 4. the four delivery assertions
// ---------------------------------------------------------------------------

// expectDelivered is the arrival check AC#6 needs, shared by every shape so the same yard
// measures the good one and the counter-shape: the native exit ran exactly once, the frame
// that reached msgcb is the library's own RPC frame carrying the page's envelope as
// params[0] verbatim, dispatchRaw's downstream saw the requestId, and nothing died in an
// unbound slot.
func expectDelivered(t *testing.T, bf *fakeDoc35r2, spy *modeSpy35r1, wantPosts int) {
	t.Helper()
	if bf.capTripped {
		t.Fatalf("AC#6 RED: the transport re-entered chrome.webview.postMessage past the depth cap (%s) - the page never leaves the page. %s",
			bf.summary(), bf.firstCapMessage())
	}
	if bf.nativeCalls != wantPosts {
		t.Fatalf("AC#6 RED: the native exit (chrome.webview.postMessage's original value) was called %d time(s), want %d. %s",
			bf.nativeCalls, wantPosts, bf.summary())
	}
	if bf.doorRounds != wantPosts {
		t.Fatalf("AC#6 RED: the bound door fired %d time(s), want %d. %s", bf.doorRounds, wantPosts, bf.summary())
	}
	if bf.deadSlots != 0 || bf.evalThrew != 0 {
		t.Fatalf("AC#6 RED: %d frame(s) died at msgcb's unbound-method branch and %d reply script(s) threw. %s",
			bf.deadSlots, bf.evalThrew, bf.summary())
	}
	for i, frame := range bf.nativeFrames {
		var d struct {
			ID     int               `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal([]byte(frame), &d); err != nil {
			t.Fatalf("AC#6 RED: frame %d that reached the native exit is not a library RPC frame: %v (%q)", i, err, frame)
		}
		if d.Method != panelDispatchBinding {
			t.Errorf("AC#6 RED: frame %d names method %q, want %q (the forwarding did not route into the bound door)", i, d.Method, panelDispatchBinding)
		}
		if len(d.Params) != 1 {
			t.Fatalf("AC#6 RED: frame %d carries %d params, want 1 (%q)", i, len(d.Params), frame)
		}
		var raw string
		if err := json.Unmarshal(d.Params[0], &raw); err != nil {
			t.Fatalf("AC#6 RED: frame %d param 0 is not a string: %v", i, err)
		}
		if raw != pageModeRequestEnvelope {
			t.Errorf("AC#6 RED: the string that arrived inside the RPC frame is not the page's envelope verbatim: got %q want %q", raw, pageModeRequestEnvelope)
		}
	}
	if len(spy.reqs) != wantPosts {
		t.Fatalf("AC#6 RED: dispatchRaw's downstream ran %d time(s), want %d. %+v", len(spy.reqs), wantPosts, spy.reqs)
	}
	for _, r := range spy.reqs {
		if r.RequestID != pageModeRequestID {
			t.Errorf("AC#6 RED: the router was reached without the page's requestId: got %q want %q", r.RequestID, pageModeRequestID)
		}
		if r.Method != panel.MethodModeRequest {
			t.Errorf("AC#6 RED: routed method %q, want %q", r.Method, panel.MethodModeRequest)
		}
		if r.Source != panel.ComposerRequestSource {
			t.Errorf("AC#6 RED: source %q, want %q", r.Source, panel.ComposerRequestSource)
		}
	}
	if bf.resolved != wantPosts {
		t.Errorf("AC#6 RED: the page got %d receipt(s) via window._rpc[id].resolve, want %d (msgcb's reply scripts: %v)", bf.resolved, wantPosts, bf.evalScripts)
	}
}

func (bf *fakeDoc35r2) firstCapMessage() string {
	if bf.capMessage == "" {
		return "(no cap message recorded)"
	}
	return bf.capMessage
}

// TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3 is the green half of the hard
// criterion: the shipping wiring (installPanelTransport, whose forwarding text is read
// from panelPostMessageForwardInit as a value) is installed, a document is opened, and the
// page posts panel.ts:246's envelope. It must leave the page through the native exit ONCE,
// as the library's own frame, and arrive at dispatchRaw's downstream carrying the page's
// requestId.
func TestPagePostMessageEnvelopeReachesNativeExitViaShapeA3(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")

	bf := newFakeDoc35r2()
	ctx := context.Background()
	if err := mgr.installPanelTransport(bf, ctx); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()

	if thrown := bf.pagePost(pageModeRequestEnvelope); thrown != "" {
		t.Fatalf("AC#6 RED: the page's own postMessage threw in the document: %s. %s", thrown, bf.summary())
	}
	if bf.maxPostDepth > postMessageHopCap {
		t.Errorf("AC#6 RED: chrome.webview.postMessage nested %d levels deep; sub-shape ③ must leave the page in one hop", bf.maxPostDepth)
	}
	expectDelivered(t, bf, spy, 1)
	t.Logf("shape ③ delivery: %s | frame=%q | door reply=%q", bf.summary(), bf.firstFrameOrNone(), bf.lastReply)
}

func (bf *fakeDoc35r2) firstFrameOrNone() string {
	if len(bf.nativeFrames) == 0 {
		return "(none)"
	}
	return bf.nativeFrames[0]
}

// TestShapeA3SecondPagePostStillDelivers catches a leaked re-entry flag: the wrapper must
// reset `inside` on the way out (the try/finally), or the page's SECOND envelope would be
// forwarded to the native exit as itself, land in msgcb's unbound-method branch, and die
// exactly the way ticket 35:59 describes.
func TestShapeA3SecondPagePostStillDelivers(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")
	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()

	for i := 0; i < 2; i++ {
		if thrown := bf.pagePost(pageModeRequestEnvelope); thrown != "" {
			t.Fatalf("AC#6 RED: page post #%d threw: %s. %s", i+1, thrown, bf.summary())
		}
	}
	expectDelivered(t, bf, spy, 2)
}

// TestShapeA3ForwardingHookIsIdempotentInOneDocument pins the hook's own guard: a second
// run of the same Init queue in one document (AddScriptToExecuteOnDocumentCreated plus a
// host that re-queues) must not wrap twice, because a wrapper whose `native` is another
// wrapper is the loop all over again.
func TestShapeA3ForwardingHookIsIdempotentInOneDocument(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")
	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()
	// re-run the whole queue against the SAME document
	for _, s := range bf.initScripts {
		bf.runScriptOn(bf.global, s)
	}
	if thrown := bf.pagePost(pageModeRequestEnvelope); thrown != "" {
		t.Fatalf("AC#6 RED: after re-running the Init queue the page's post threw: %s. %s", thrown, bf.summary())
	}
	expectDelivered(t, bf, spy, 1)
}

// TestLegacySubShapeOneHookDiesInAReentryLoop is the red half of the hard criterion, made
// a permanent nail: the forwarding that shipped in fb2fb802 (verbatim above) is fed to this
// same behavioural model, and the model reports the loop instead of a delivery - the native
// exit called ZERO times while chrome.webview.postMessage re-enters past the cap. That is
// 35-v1's G2-4 reading, now measurable without a browser. The counter-shape's own delivery
// assertions are run first and are expected to fail, which is what shows this yard is not
// reading a string.
func TestLegacySubShapeOneHookDiesInAReentryLoop(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")

	bf := newFakeDoc35r2()
	// the counter-shape in the position the production hook occupies: bound door first,
	// then the forwarding script, exactly the order installPanelTransport uses.
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	// Drop the shipping forwarding (the last script queued is panelPostMessageForwardInit)
	// and put sub-shape ① in its place.
	bf.initScripts = bf.initScripts[:len(bf.initScripts)-1]
	bf.registerInit(legacySubShapeOneHook)
	bf.openDocument()

	thrown := bf.pagePost(pageModeRequestEnvelope)
	if !bf.capTripped {
		t.Fatalf("COUNTER-SHAPE NOT RED: sub-shape ① delivered without tripping the re-entry cap. %s (page throw=%q)", bf.summary(), thrown)
	}
	if bf.nativeCalls != 0 {
		t.Fatalf("COUNTER-SHAPE READING WRONG: the native exit was called %d time(s), want 0 for the looping hook. %s", bf.nativeCalls, bf.summary())
	}
	if bf.doorRounds != 0 || len(spy.reqs) != 0 {
		t.Fatalf("COUNTER-SHAPE READING WRONG: Go saw doorRounds=%d routerRuns=%d, want 0/0 - the loop must not deliver. %s", bf.doorRounds, len(spy.reqs), bf.summary())
	}
	t.Logf("sub-shape ① under the behavioural yard: re-entered %d levels (cap %d), native exit called %d times, door fired %d times, page error=%q",
		bf.capDepth, postMessageHopCap, bf.nativeCalls, bf.doorRounds, thrown)

	// And the same yard must refuse to certify it as a delivery.
	if bf.capTripped == false || bf.nativeCalls != 0 {
		t.Fatalf("the yard failed to reject the counter-shape")
	}
}

// TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot is the ruler's own calibration (⛔ it is
// not evidence for AC#6 - it installs NO forwarding): with only the library's pipe, the
// page's raw envelope does reach the native exit, msgcb parses it into rpcMessage with
// ID=0 (webview.go:131-135 answers a no-id envelope happily), callbinding finds nothing
// bound and answers (nil, nil), and the reply script targets window._rpc[0], which the page
// never created. That is ticket 35:52-59's death, modelled rather than quoted: the door
// fires zero times and the reply throws a TypeError nobody sees.
func TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")

	bf := newFakeDoc35r2()
	// Bind the door the way the host does, but register no forwarding script at all.
	if err := bf.Bind(panelDispatchBinding, func(raw string) string {
		reply, _ := mgr.dispatchRaw(context.Background(), raw)
		return reply
	}); err != nil {
		t.Fatalf("Bind: %v", err)
	}
	bf.openDocument()

	if thrown := bf.pagePost(pageModeRequestEnvelope); thrown != "" {
		t.Fatalf("calibration: the page's raw post threw before msgcb: %q", thrown)
	}
	if bf.nativeCalls != 1 {
		t.Fatalf("calibration: want the raw envelope to reach the native exit once, got %d", bf.nativeCalls)
	}
	if bf.doorRounds != 0 || len(spy.reqs) != 0 {
		t.Fatalf("calibration: an unbound envelope must not reach the door or the router (door=%d router=%d)", bf.doorRounds, len(spy.reqs))
	}
	if bf.deadSlots != 1 {
		t.Fatalf("calibration: want exactly one frame in msgcb's unbound-method branch, got %d", bf.deadSlots)
	}
	if len(bf.evalScripts) != 1 || bf.evalScripts[0] != `window._rpc[0].resolve(null); window._rpc[0] = undefined` {
		t.Fatalf("calibration: the reply script is %v, want the library's resolve(null) form (webview.go:157)", bf.evalScripts)
	}
	if bf.evalThrew != 1 {
		t.Fatalf("calibration: want the reply to throw against a page with no _rpc[0] (swallowed TypeError), got evalThrew=%d", bf.evalThrew)
	}
	if bf.resolved != 0 || bf.rejected != 0 {
		t.Fatalf("calibration: nothing may be settled here, got resolved=%d rejected=%d", bf.resolved, bf.rejected)
	}
}

// ---------------------------------------------------------------------------
// 7. 35-r4: the two faces A684 §3 found this fixture blind to (ticket 35:75 (a))
// ---------------------------------------------------------------------------
//
// WHAT THIS SECTION IS FOR. 35-v2 attacked this yard with two one-line mutants of the
// shipped hook and BOTH stayed green (ledger A684 §3): M-A replaced native.call(cw,
// message) with a bare native(message), M-B deleted the 'typeof ... !== "function"' half
// of the guard. Neither reading was the interpreter's fault: the native exit stub
// discarded its receiver and jsObject.set had no concept of a non-writable property, so
// the two shapes had nothing to be wrong *in*. Both halves are now modelled -
// a receiver-strict native exit, and a per-object writability table - and each face gets
// its own named red below.
//
// WHAT THIS BUYS AND WHAT IT DOES NOT. It buys "a hook that drops this, or whose
// assignment is silently refused, is no longer certified as a delivery". It does NOT buy
// "WebView2 really binds the receiver / really lets the page replace postMessage" -
// A684 §1 puts that格 with ticket 35:52's real-window evidence, and no assertion here
// stands in for it. Object.defineProperty is NOT implemented as a JS-visible method (only
// the Go-side markNonWritable台件), so a page that descriptor-reads or descriptor-redefines
// its own transport still lives outside this yard; descriptor READS
// (Object.getOwnPropertyDescriptor), configurable and delete are named as un-modelled in
// .scratch/wisp/probes/35/r4/impl.md rather than quietly skipped.

// runShippedHookInOneWorld35r4 installs the SHIPPING wiring (installPanelTransport: the
// bound door, then the production forwarding string read as a value) into a document,
// posts the page's own envelope once, and hands back what the caller asserts on.
// nonWritable selects the world where chrome.webview.postMessage is exposed as a
// non-writable data property.
func runShippedHookInOneWorld35r4(t *testing.T, nonWritable bool) (*fakeDoc35r2, *modeSpy35r1, string) {
	t.Helper()
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")
	bf := newFakeDoc35r2()
	bf.nativeExitNonWritable = nonWritable
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()
	return bf, spy, bf.pagePost(pageModeRequestEnvelope)
}

// TestForwardingHookMustNotLoseTheNativeExitReceiver is face M-A. The shipped hook hands
// the native exit its receiver (native.call(cw, message)); 35-v2's mutant dropped it with
// one character-level edit (native(message)) and all five cases above stayed green,
// because the stub took func(_ jsValue, ...) and could not see the difference. The stub
// now demands cw itself - the rule THIS FIXTURE answers with, modelled after what Chrome
// and WebView2 are reported to do (this yard never observes a browser; AC#8(iii)) - so the
// receiver-less shape dies here by name instead of passing silently.
//
// TOOTH: overlay-mutating the shipped hook to native(message) reddens the FIRST branch
// below with the thrown "TypeError: Illegal invocation: chrome.webview.postMessage lost
// its receiver" in its text; overlay-mutating the stub back to discarding recv turns the
// mutant green again - both readings are recorded in .scratch/wisp/probes/35/r4/logs.
func TestForwardingHookMustNotLoseTheNativeExitReceiver(t *testing.T) {
	bf, spy, thrown := runShippedHookInOneWorld35r4(t, false)

	if thrown != "" {
		t.Fatalf("M-A RED (receiver lost at the native exit): the page's own postMessage threw inside the document, which under this model means the forwarding hook invoked chrome.webview.postMessage without its receiver: %s. %s",
			thrown, bf.summary())
	}
	if bf.nativeReceiverMismatch != 0 {
		t.Fatalf("M-A RED (receiver lost at the native exit): %d native-exit invocation(s) reached the host without chrome.webview as their receiver. THIS FIXTURE answers that with Illegal invocation - modelled after what Chrome and WebView2 are reported to do, NOT a WebView2 behaviour record (AC#8(iii); the real-window reading is the 〔仅本机可量〕 winlive rig) - and on that model the page's letter never leaves the page, so this hook must not be certified as a delivery. %s",
			bf.nativeReceiverMismatch, bf.summary())
	}
	if bf.nativeExitReceiverOK != 1 {
		t.Fatalf("M-A READING WRONG: the native exit was honoured with its receiver %d time(s), want exactly 1. Zero means the receiver face was never exercised here, so this case would prove nothing. %s",
			bf.nativeExitReceiverOK, bf.summary())
	}
	expectDelivered(t, bf, spy, 1)
	t.Logf("M-A face has teeth: 1 native-exit call, all with chrome.webview as receiver (receiverLost=%d); a bare native(message) trips it", bf.nativeReceiverMismatch)
}

// TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit is face M-B's writability
// half: the fixture used to have NO writable/configurable concept (grep of this file: zero
// occurrences), so "the host does not let the page replace its own transport, and the
// forwarding hook therefore never arms" could not be expressed at all - an override always
// won, and the yard reported a delivery.
//
// In the non-writable world the assignment must be REFUSED, SILENTLY (sloppy mode: no
// TypeError), and the consequence must be the loud reading: the page's letter reaches the
// native exit unforwarded, dies at msgcb's unbound-method branch (ticket 35:52-59's death)
// and never reaches the Go door. The writable world is run through the same helper as the
// control, so this case is not an artefact of a台件 nobody checked.
func TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit(t *testing.T) {
	bf, spy, thrown := runShippedHookInOneWorld35r4(t, true)

	if thrown != "" {
		t.Fatalf("M-B1 RED (wrong death): a non-writable chrome.webview.postMessage must refuse the hook's assignment SILENTLY (sloppy mode, no TypeError), but the page's own post threw: %s. %s",
			thrown, bf.summary())
	}
	if got := bf.webview.ignoredWrites; len(got) != 1 || got[0] != "postMessage" {
		t.Fatalf("M-B1 RED (fixture blind to writability): the hook's assignment to a property this world marks NON-WRITABLE took effect unnoticed - ignoredWrites=%v, want exactly [postMessage]. This is the A684 §3 face this case exists to make bite: without a refused write the yard cannot tell an armed hook from a silently ignored override. %s",
			got, bf.summary())
	}
	if bf.nativeCalls != 1 || bf.nativeReceiverMismatch != 0 {
		t.Fatalf("M-B1 READING WRONG: the page's letter must still leave through the untouched native exit exactly once (and with its receiver): got nativeExitCalls=%d receiverLost=%d. %s",
			bf.nativeCalls, bf.nativeReceiverMismatch, bf.summary())
	}
	if bf.doorRounds != 0 || len(spy.reqs) != 0 {
		t.Fatalf("M-B1 RED: with the native exit NOT writable the page's letter still reached the Go door (doorRounds=%d routerRuns=%d) - the model let a silently refused override behave like a successful hook. %s",
			bf.doorRounds, len(spy.reqs), bf.summary())
	}
	if bf.deadSlots != 1 || bf.evalThrew != 1 {
		t.Fatalf("M-B1 RED: an unarmed forwarding must leave the raw envelope to die at msgcb's unbound-method branch (unboundSlots=1, swallowed reply TypeError=1), got unboundSlots=%d evalThrew=%d. %s",
			bf.deadSlots, bf.evalThrew, bf.summary())
	}
	if len(bf.nativeFrames) != 1 || bf.nativeFrames[0] != pageModeRequestEnvelope {
		t.Fatalf("M-B1 READING WRONG: want the page's envelope to reach the native exit UNFORWARDED exactly once, got %v", bf.nativeFrames)
	}

	bf2, spy2, thrown2 := runShippedHookInOneWorld35r4(t, false)
	if thrown2 != "" {
		t.Fatalf("M-B1 CONTROL RED: the same shipped hook in a writable world threw: %s. %s", thrown2, bf2.summary())
	}
	if len(bf2.webview.ignoredWrites) != 0 {
		t.Fatalf("M-B1 CONTROL RED: the writable world refused %d assignment(s) (%v) - the台件 is not selecting on writability. %s",
			len(bf2.webview.ignoredWrites), bf2.webview.ignoredWrites, bf2.summary())
	}
	expectDelivered(t, bf2, spy2, 1)
	t.Logf("M-B1 face has teeth: non-writable world refused the hook (doorRounds=0, unboundSlots=1), writable world armed it (doorRounds=%d)", bf2.doorRounds)
}

// TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent is M-B's other half as
// 35-v2 hit it: deleting 'typeof window.wispDispatch !== "function"' from the shipped
// guard changed nothing, because in every world this yard ran the door was always bound.
// The half-branch exists for the document whose binding did not register, so that world is
// now built (the Bind stub script is dropped, the forwarding hook is kept) and the guard
// has work to do: the hook must fall back to the native exit and forward byte-for-byte,
// not call a non-function.
//
// TOOTH: overlay-deleting the typeof half makes the wrapper reach window.wispDispatch,
// which is undefined here, and the page's post throws a TypeError the case reports by name
// (logs/mut-typof-half.txt).
func TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")
	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	// installPanelTransport queues the library's Bind stub first and the forwarding hook
	// second (panel_host_windows.go:693-699, the order this yard already models); dropping
	// script 0 leaves a document with the hook and NO window.wispDispatch.
	bf.initScripts = bf.initScripts[1:]
	if len(bf.initScripts) != 1 {
		t.Fatalf("calibration: want the forwarding hook alone queued, got %d script(s)", len(bf.initScripts))
	}
	bf.openDocument()
	if bf.window.has(panelDispatchBinding) {
		t.Fatalf("calibration: the door %q is present in the page, so this world does not exercise the guard's typeof half at all", panelDispatchBinding)
	}

	thrown := bf.pagePost(pageModeRequestEnvelope)
	if thrown != "" {
		t.Fatalf("M-B2 RED (guard's typeof half gone): with no door in the page the forwarding hook must fall back to the native exit; it instead invoked a non-function and the page's post threw: %s. %s",
			thrown, bf.summary())
	}
	if bf.capTripped {
		t.Fatalf("M-B2 RED: the door-absent fallback re-entered chrome.webview.postMessage past the cap (%s). %s", bf.firstCapMessage(), bf.summary())
	}
	if bf.nativeCalls != 1 || bf.nativeExitReceiverOK != 1 || bf.nativeReceiverMismatch != 0 {
		t.Fatalf("M-B2 READING WRONG: the fallback must reach the native exit exactly once, with its receiver: nativeExitCalls=%d receiverOK=%d receiverLost=%d. %s",
			bf.nativeCalls, bf.nativeExitReceiverOK, bf.nativeReceiverMismatch, bf.summary())
	}
	if len(bf.nativeFrames) != 1 || bf.nativeFrames[0] != pageModeRequestEnvelope {
		t.Fatalf("M-B2 READING WRONG: the door-absent fallback must forward the page's bytes untouched, got %v", bf.nativeFrames)
	}
	if bf.doorRounds != 0 || len(spy.reqs) != 0 || bf.deadSlots != 1 || bf.evalThrew != 1 {
		t.Fatalf("M-B2 READING WRONG: an absent door must leave the envelope dead at msgcb's unbound branch (doorRounds=0 routerRuns=0 unboundSlots=1 evalThrew=1), got doorRounds=%d routerRuns=%d unboundSlots=%d evalThrew=%d. %s",
			bf.doorRounds, len(spy.reqs), bf.deadSlots, bf.evalThrew, bf.summary())
	}
	t.Logf("door absent: the guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=%d, no throw)", len(bf.nativeFrames))
}

// TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable is ticket 35 AC#8(i),
// the face `35-v4` measured as a REAL tautology (恒绿): replacing the shipped guard's
// typeof half (cmd/wisp/panel_host_windows.go:680, expanded as
// typeof window.wispDispatch !== "function") with the truthiness test !window.wispDispatch
// left all twelve cases green - reading .scratch/wisp/probes/35/v4/logs/run-v-guard-truthy.txt.
// The two branches only diverge in a world where the door NAME exists but holds something
// uncallable, and this yard had no such world: the door's only construction point is the
// library's Bind stub (:1229 -> libBindStubScript35r2 at :97), which injects a function.
// So the case that documents itself as guarding "not call a non-function" (:2019-2021 above)
// had only ever proved the undefined half of that sentence.
//
// THE WORLD: window.wispDispatch is a TRUTHY NON-FUNCTION (a string). Staged Go-side, the
// same way 35-r4 stages the writability world with markNonWritable (:149) - deliberately not
// as a JS Object method: Object.defineProperty stays invisible to this interpreter (named in
// the r4 un-modelled list), and AC#8(i) asks for no descriptor. The guard reads
// window.<binding> at POST time rather than at hook-install time, so overwriting the property
// between openDocument and pagePost is the whole world this face needs.
//
// TOOTH: overlay-replacing the guard's typeof half with the truthiness test makes this case
// the ONLY red one in the family - the wrapper reaches window.wispDispatch, the interpreter
// answers a string invocation with "TypeError: string is not a function", and the case
// reports it by name. Under the shipped guard it is green. Both readings are in
// .scratch/wisp/probes/35/r5/logs.
func TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")
	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()

	// Calibration first: this face is only this face if the door NAME is there. Without
	// this check the台件 could silently degrade into the :2026 world and re-green the same
	// mutant that face already reddens.
	if !bf.window.has(panelDispatchBinding) {
		t.Fatalf("M-B3 CALIBRATION RED: this world must hold the door name %q so the guard has a value to test; the page has none, which is the :2026 door-absent face, not AC#8(i)'s.", panelDispatchBinding)
	}
	if _, isFunction := bf.window.props[panelDispatchBinding].(*jsFunction); !isFunction {
		t.Fatalf("M-B3 CALIBRATION RED: the Bind stub did not put a function behind %q, so the台件 never had a callable door to make uncallable (got %s).", panelDispatchBinding, jsTypeOf(bf.window.props[panelDispatchBinding]))
	}
	bf.window.set(panelDispatchBinding, "not-a-function")
	if got := jsTypeOf(bf.window.props[panelDispatchBinding]); got == "undefined" || got == "function" {
		t.Fatalf("M-B3 CALIBRATION RED: the台件 was supposed to leave a TRUTHY NON-FUNCTION in the page, got typeof window.%s == %s.", panelDispatchBinding, got)
	}

	thrown := bf.pagePost(pageModeRequestEnvelope)
	if thrown != "" {
		t.Fatalf("M-B3 RED (guard reduced to a truthiness test): a door that is present but NOT callable must still fall back to the native exit; this hook invoked window.%s anyway and the page's post threw: %s. %s",
			panelDispatchBinding, thrown, bf.summary())
	}
	if bf.capTripped {
		t.Fatalf("M-B3 RED: the not-callable-door fallback re-entered chrome.webview.postMessage past the cap (%s). %s", bf.firstCapMessage(), bf.summary())
	}
	if bf.nativeCalls != 1 || bf.nativeExitReceiverOK != 1 || bf.nativeReceiverMismatch != 0 {
		t.Fatalf("M-B3 READING WRONG: the fallback must reach the native exit exactly once, with its receiver: nativeExitCalls=%d receiverOK=%d receiverLost=%d. %s",
			bf.nativeCalls, bf.nativeExitReceiverOK, bf.nativeReceiverMismatch, bf.summary())
	}
	if len(bf.nativeFrames) != 1 || bf.nativeFrames[0] != pageModeRequestEnvelope {
		t.Fatalf("M-B3 READING WRONG: a present-but-uncallable door must still forward the page's bytes untouched, got %v", bf.nativeFrames)
	}
	if bf.doorRounds != 0 || len(spy.reqs) != 0 || bf.deadSlots != 1 || bf.evalThrew != 1 {
		t.Fatalf("M-B3 READING WRONG: an uncallable door must leave the envelope dead at msgcb's unbound branch (doorRounds=0 routerRuns=0 unboundSlots=1 evalThrew=1), got doorRounds=%d routerRuns=%d unboundSlots=%d evalThrew=%d. %s",
			bf.doorRounds, len(spy.reqs), bf.deadSlots, bf.evalThrew, bf.summary())
	}
	t.Logf("M-B3 face has teeth: door present as typeof===%q, guard's typeof half routed the page's envelope to the native exit byte-for-byte (frames=%d, no throw)",
		jsTypeOf(bf.window.props[panelDispatchBinding]), len(bf.nativeFrames))
}
