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
	// onGet lets the browser model wrap one property's reads (chrome.webview.postMessage)
	// without changing what the scripts see. nil for ordinary objects.
	onGet func(name string, cur jsValue) jsValue
}

func newJSObject() *jsObject { return &jsObject{props: map[string]jsValue{}} }

func (o *jsObject) has(name string) bool { _, ok := o.props[name]; return ok }

func (o *jsObject) set(name string, v jsValue) {
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

func (p *jsParser) peek() jsTok    { return p.toks[p.i] }
func (p *jsParser) next() jsTok    { t := p.toks[p.i]; p.i++; return t }
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
	cw.set("postMessage", &jsFunction{name: "native", host: func(_ jsValue, args []jsValue) jsValue {
		bf.nativeCalls++
		msg := ""
		if len(args) > 0 {
			msg = jsToString(args[0])
		}
		bf.nativeFrames = append(bf.nativeFrames, msg)
		bf.msgcb(msg)
		return nil
	}})
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
	return fmt.Sprintf("nativeExitCalls=%d doorRounds=%d maxPostDepth=%d capTripped=%v capDepth=%d unparsable=%d unboundSlots=%d evalThrew=%d resolved=%d rejected=%d scriptThrows=%v",
		bf.nativeCalls, bf.doorRounds, bf.maxPostDepth, bf.capTripped, bf.capDepth,
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
