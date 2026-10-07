# 35-a2 — Transport-cost survey (read-only leg, ticket 35 / C17 panel bridge)

Leg: `35-a2` (read-only cost census). Shape NOT chosen here (选形归编排者).
Scope = ticket 35's new AC frame "Transport agreement page↔host" (uncheckable by this leg, ⛔ not flipped).

---

## §0 起手锚 (landing gate — first commit, before any long-running command)

```
$ date
Wed Oct  7 11:55:16 CST 2026

$ git rev-parse --short HEAD
28872c9a

$ git branch --show-current
dev

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
 M docs/reports/pending-and-issues.md
 M scripts/build.ps1
```

Pre-existing dirty entries NOT touched by this leg (registered as-is, per §Git 纪律):
`design/**`, `.gitignore`, several `.scratch/**`, `docs/reports/pending-and-issues.md`,
`scripts/build.ps1` (another leg `274-r1` is writing it).

### 8 anchors, re-read live at HEAD `28872c9a` (verbatim)

**Anchor 1 — `go.mod:15-21`** (the dependency marking):
```
15:
16: require (
17:	github.com/dustin/go-humanize v1.0.1 // indirect
18:	github.com/google/uuid v1.6.0 // indirect
19:	github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect
20:	github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect
21:	github.com/k2-fsa/sherpa-onnx-go-linux v1.13.8 // indirect
```

**Anchor 2 — library root `webview.go:99-110`** (`MessageCallback` wiring):
```
99: 	w.bindings = map[string]interface{}{}
100: 	w.autofocus = options.AutoFocus
101:
102: 	chromium := edge.NewChromium()
103: 	chromium.MessageCallback = w.msgcb
104: 	chromium.DataPath = options.DataPath
105: 	chromium.SetPermission(edge.CoreWebView2PermissionKindClipboardRead, edge.CoreWebView2PermissionStateAllow)
106:
107: 	w.browser = chromium
108: 	w.mainthread, _, _ = w32.Kernel32GetCurrentThreadID.Call()
109: 	if !w.CreateWithOptions(options.WindowOptions) {
110: 		return nil
```

**Anchor 3 — library `webview.go:139-160`** (`msgcb`, the RPC parser):
```
139: func (w *webview) msgcb(msg string) {
140: 	d := rpcMessage{}
141: 	if err := json.Unmarshal([]byte(msg), &d); err != nil {
142: 		log.Printf("invalid RPC message: %v", err)
143: 		return
144: 	}
145:
146: 	id := strconv.Itoa(d.ID)
147: 	if res, err := w.callbinding(d); err != nil {
148: 		w.Dispatch(func() {
149: 			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
150: 		})
151: 	} else if b, err := json.Marshal(res); err != nil {
152: 		w.Dispatch(func() {
153: 			w.Eval("window._rpc[" + id + "].reject(" + jsString(err.Error()) + "); window._rpc[" + id + "] = undefined")
154: 		})
155: 	} else {
156: 		w.Dispatch(func() {
157: 			w.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")
158: 		})
159: 	}
160: }
```

**Anchor 4 — library `webview.go:450-482`** (`Bind`, the injected script):
```
450: func (w *webview) Bind(name string, f interface{}) error {
451: 	v := reflect.ValueOf(f)
452: 	if v.Kind() != reflect.Func {
453: 		return errors.New("only functions can be bound")
454: 	}
455: 	if n := v.Type().NumOut(); n > 2 {
456: 		return errors.New("function may only return a value or a value+error")
457: 	}
458: 	w.m.Lock()
459: 	w.bindings[name] = f
460: 	w.m.Unlock()
461:
462: 	w.Init("(function() { var name = " + jsString(name) + ";" + `
463: 		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
464: 		window[name] = function() {
465: 		  var seq = RPC.nextSeq++;
466: 		  var promise = new Promise(function(resolve, reject) {
467: 			RPC[seq] = {
468: 			  resolve: resolve,
469: 			  reject: reject,
470: 			};
471: 		  });
472: 		  window.external.invoke(JSON.stringify({
473: 			id: seq,
474: 			method: name,
475: 			params: Array.prototype.slice.call(arguments),
476: 		  }));
477: 		  return promise;
478: 		}
479: 	})()`)
480:
481: 	return nil
482: }
```

**Anchor 5 — library `webview.go:435-448`** (the four outbound primitives):
```
435: func (w *webview) Init(js string) {
436: 	w.browser.Init(js)
437: }
438:
439: func (w *webview) Eval(js string) {
440: 	w.browser.Eval(js)
441: }
442:
443: func (w *webview) Dispatch(f func()) {
444: 	w.m.Lock()
445: 	w.dispatchq = append(w.dispatchq, f)
446: 	w.m.Unlock()
447: 	_, _, _ = w32.User32PostThreadMessageW.Call(w.mainthread, w32.WMApp, 0, 0)
448: }
```

**Anchor 6 — `pkg/edge/chromium.go:233-248`** (`MessageReceived`, echo back to page):
```
233: func (e *Chromium) MessageReceived(sender *ICoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
234: 	var message *uint16
235: 	_, _, _ = args.vtbl.TryGetWebMessageAsString.Call(
236: 		uintptr(unsafe.Pointer(args)),
237: 		uintptr(unsafe.Pointer(&message)),
238: 	)
239: 	if e.MessageCallback != nil {
240: 		e.MessageCallback(w32.Utf16PtrToString(message))
241: 	}
242: 	_, _, _ = sender.vtbl.PostWebMessageAsString.Call(
243: 		uintptr(unsafe.Pointer(sender)),
244: 		uintptr(unsafe.Pointer(message)),
245: 	)
246: 	windows.CoTaskMemFree(unsafe.Pointer(message))
247: 	return 0
248: }
```

**Anchor 7 — `pkg/edge/corewebview2.go:104-112`** (COM vtable names; note `AsJSON` casing):
```
104: 	CapturePreview                         ComProc
105: 	Reload                                 ComProc
106: 	PostWebMessageAsJSON                   ComProc
107: 	PostWebMessageAsString                 ComProc
108: 	AddWebMessageReceived                  ComProc
109: 	RemoveWebMessageReceived               ComProc
110: 	CallDevToolsProtocolMethod             ComProc
111: 	GetBrowserProcessID                    ComProc
112: 	GetCanGoBack                           ComProc
```

**Anchor 8 — `cmd/wisp/panel_host_windows.go`**, four live reads:

`:78-82` (binding name):
```
78:
79: const (
80: 	panelDispatchBinding = "wispDispatch"
81: 	// panelWidthPx / panelHeightPx are what the host asks WebView2 for when NOBODY
82: 	// handed it a geometry source (票 255 AC#4: these two are now "the value when
```
`:400-412` (the one inbound door, and the comment the orchestrator already flagged as stale):
```
400: 	// The one inbound door. JS -> Go via a bound function over the WebView2
401: 	// message channel; the returned string is the receipt that reaches the page.
402: 	// No new inbound method name is invented here: the raw envelope is parsed by
403: 	// panel.ParseComposerRequest inside Handle, whose whitelist is fixed at the
404: 	// four existing methods (bridge.go:42-45, ticket 248 owns config methods).
405: 	bindErr := w.Bind(panelDispatchBinding, func(raw string) string {
406: 		reply, _ := m.dispatchRaw(ctx, raw)
407: 		return reply
408: 	})
409: 	if bindErr != nil {
410: 		w.Destroy()
411: 		return fmt.Errorf("panel host: bind dispatch door: %w", bindErr)
412: 	}
```
`:440-470` (`SetHtml` call sites — the Init-timing question):
```
440: // the user is never left looking at the round-trip probe page, and it still opens
441: // no socket and loads nothing over the network.
442: func (m *PanelManager) serveNotBuiltNoticeLocked() {
443: 	m.mu.Lock()
444: 	m.mu.Unlock()  [tail: w := m.w / if w == nil { return }]
445: 	w.SetHtml(`<!doctype html>... panel assets unavailable ...`)
...   serveEntry(): w.SetHtml(string(data))
```
(full literal text of `:440-470` is in `logs/anchor-reads.txt`; the two `SetHtml` statements are at
`panel_host_windows.go:449` and `:468`, third one recorded in §2.)
`:625-635` (`dispatchRaw`):
```
625: 		}
626: 	}
627:
628: // dispatchRaw runs one page envelope through the router. Kept separate from the
629: // bind closure so a test can call the same code path the page reaches.
630: func (m *PanelManager) dispatchRaw(ctx context.Context, raw string) (string, error) {
631: 	if m.disp == nil {
632: 		return "panel host: no inbound router attached", fmt.Errorf("panel host: nil router")
633: 	}
634: 	if ctx == nil {
635: 		ctx = context.Background()
```
`:665-672` (`wispProbeRT` second binding):
```
665: 	// A second, one-shot binding the probe page calls; its Go side just signals
666: 	// the round trip completed. Using a binding keeps the probe on the same
667: 	// WebView2 channel the real door uses.
668: 	if err := w.Bind("wispProbeRT", func() string {
669: 		select {
670: 		case <-done:
671: 		default:
672: 			close(done)
```

**Anchor 9 — `internal/panel/bridge.go:38-52`** (the roster, 6 entries):
```
38: // test that does not exist in this repository; ticket 35's snapshot pump was the
39: // step that checked it, and the behaviour was already covered - only the name
40: // was wrong.)
41: const (
42: 	MethodModeRequest      = "panel.mode.request"
43: 	MethodWorkspaceRequest = "panel.workspace.request"
44: 	MethodAttachmentAdd    = "panel.attachment.add"
45: 	MethodMessageSend      = "panel.message.send"
46: 	// MethodConfigGet and MethodConfigSet are ticket 248 AC#1's two settings
47: 	// routes. Three facts about the naming, because each one is a decision
48: 	// somebody could later "fix" by accident:
49: 	//
50: 	//   - The names are NOT invented here. They are the two spellings the
51: 	//     product spec already carries (docs/specs/SPEC-08-ui-ball-panel.md:168,
52: 	//     read-only for this leg), so this file follows a document rather than
```

**Anchor 10 — frontend on the `dev` tree** (line counts re-measured:
`panel.ts` 311, `main.tsx` 61, `App.tsx` 117 — the brief's page line numbers below are dev-true):
```
frontend/src/lib/panel.ts:139-155
139: /** Shape of the object WebView2 installs on the host page (C17). */
140: interface WispHostBridge {
141:   postMessage(message: string): void;
142: }
143:
144: declare global {
145:   interface Window {
146:     /** Installed by WebView2's AddHostObjectToScript / postMessage pipe. */
147:     wispBridge?: WispHostBridge;
148:     chrome?: { webview?: WispHostBridge };
149:   }
150: }
151:
152: function hostBridge(): WispHostBridge | null {
153:   if (typeof window === "undefined") return null;
154:   return window.wispBridge ?? window.chrome?.webview ?? null;
155: }
```
```
frontend/src/lib/panel.ts:162-166 (the "push" promise)
162: /**
163:  * Send one bridge request and return nothing: responses arrive as a fresh
164:  * PanelSnapshot push, never as a return value, because a panel that keeps a
165:  * copy of the answer would be a second state holder (PLAN.md:1044).
166:  */
```
```
frontend/src/lib/panel.ts:179-186 (sender 1)   /   :218-225 (sender 2)
179:   bridge.postMessage(
180:     JSON.stringify({
181:       method: "panel.approval.request",
182:       correlationId,
183:       outcome,
184:     }),
185:   );
186: }
...
218:   bridge.postMessage(
219:     JSON.stringify({
220:       method,
221:       requestId: nextRequestId(),
222:       source: "panel-composer",
223:       ...payload,
224:     }),
225:   );
```
```
frontend/src/main.tsx:56-60
56: createRoot(document.getElementById("root")!).render(
57:   <StrictMode>
58:     {harness === "1" ? <AppHarness /> : harness === "2" ? <Showcase /> : <App />}
59:   </StrictMode>,
60: );
```
```
frontend/src/App.tsx:49-54 (EMPTY)  /  :69-73 (default prop)
49: const EMPTY: PanelSnapshot = {
50:   pending: [],
51:   results: [],
52:   composer: EMPTY_COMPOSER,
53:   generatedAt: "",
54: };
...
69: export default function App({
70:   snapshot = EMPTY,
71: }: {
72:   snapshot?: PanelSnapshot;
73: }) {
```

**Anchor 11 — `internal/panel/composer_dispatch_test.go:455-465`** (the empty stub):
```
455: 		writeHostCarrier(t, a, filepath.Join("internal", "panel"), "host_windows.go", `package panel
456:
457: type CoreWebView2 struct{}
458:
459: func (w *CoreWebView2) WebMessageReceived(sender, args any) error { return nil }
460:
461: func (w *CoreWebView2) PostWebMessageAsJson(json string) error { return nil }
462: `)
463: 		if got := hostChannelCapabilityHits(t, a); len(got) == 0 {
464: 			t.Errorf("POSITIVE CONTROL RED (the ruler, not the product): a tree whose production source hosts a "+
465: 				"CoreWebView2 and takes WebMessageReceived answered 0 capability hits, so the check above is a "+
```

### Corrections to the brief's page anchors (dev-true, found while re-reading)
The brief and `35-a1` write `hostBridge()` at `panel.ts:146-153` and the push promise at `:159-163`.
Live dev reads: `hostBridge()` is `:152-155` (`:146` is the stale comment line inside `declare global`),
and the "return nothing / PanelSnapshot push" sentence is `:162-166`. `:141`/`:179`/`:218` in the AC
frame are correct as written.

### Library-vs-stub spelling finding (new this leg, affects §2 shape 乙 and §5 stub cell)
The real COM vtable field is `PostWebMessageAsJSON` (capital `JSON`, `corewebview2.go:106`).
The stub at `composer_dispatch_test.go:461` is spelled `PostWebMessageAsJson` and is declared on a
locally synthesised `type CoreWebView2 struct{}` inside a `writeHostCarrier` fixture string —
⛔ it is NOT an implementation of the library's type. Recorded, not touched.
