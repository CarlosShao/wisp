# 35-r2 起手锚（票 35 AC#6 `:52` 页面<->宿主传输：甲③ 落地 + 判据换行为尺）

腿＝`35-r2`（产码腿，写面只有 `cmd/wisp/**` 与本草证目录）。
起手时刻 `2026-10-07 19:24:38 +08`；起手 HEAD 顶＝`a7f9781c`（编排者 A682 那笔，票 `:150` 已具名改裁甲③）。

## 现量读数（全部本腿自跑）

| 尺 | 读数 |
|---|---|
| `git status --porcelain \| wc -l` | 753（共享工作树既有脏，非本腿造成） |
| `git diff --stat a7f9781c -- cmd/wisp/ internal/panel/ frontend/` | 空（三处起手都干净） |
| `cmd/wisp/panel_host_windows.go` CR 数 工作树／`git show HEAD:` | 0 ／ 0（本机 autocrlf 幻影与本文件无关） |
| `cmd/wisp/panel_transport_35r1_test.go` CR 数 工作树／HEAD | 0 ／ 0 |
| 转发钩子常量行号 | `:648`（`const panelPostMessageForwardInit = \`(function () {`）＝简报所写一致 |
| `w.Init(...)` 注册行号 | `:671` ＝ v1 突变那一发的锚 |
| `dispatchRaw` | `:677` ＝ v1 具名的漂移后现量（票面旧值 `:630` 作废） |
| `panelDispatchBinding` | `:80` ＝ `"wispDispatch"` |
| C17 守卫名册（`internal/panel/bridge.go:42-45`＋`:66-67`） | **六枚**：`panel.mode.request` / `panel.workspace.request` / `panel.attachment.add` / `panel.message.send` / `config.get` / `config.set`（与简报"六枚"一致；v1 亦量得六枚） |

## 库原文（本腿自己 `sed -n` 取，⛔ 不是转述；`github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`）

- `pkg/edge/chromium.go:112` 逐字＝`e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")` ⇒ 库桩唯一出口＝`chrome.webview.postMessage` **属性**，且箭头函数在**调用时**取该属性当前值。
- `webview.go:450` `func (w *webview) Bind`；`:462-478` 注入的桩逐字（本腿取回原文）：同步 `var seq = RPC.nextSeq++;` → `new Promise(...)` 里 `RPC[seq] = {resolve, reject}` → `window.external.invoke(JSON.stringify({id: seq, method: name, params: Array.prototype.slice.call(arguments)}))` → `return promise;`。
- `webview.go:131-135` `type rpcMessage struct { ID int; Method string; Params []json.RawMessage }`（⚠ 页面裸信封 `{"method":…,"requestId":…}` 能**成功**解析成 `ID=0`＋未绑名，不是 JSON 错误）。
- `webview.go:139-160` `msgcb`＋`:163-168` `callbinding`：未绑名 ⇒ `return nil, nil`（**不是 error**）⇒ 走 `:157` 的 **`window._rpc[id].resolve(…)`** 成功支（⛔ 不是票 `:58-59` 旧句的 `reject`；A682 §2 已具名更正）。页面没有 `_rpc[0]` ⇒ 该 Eval 在浏览器里是一枚被吞掉的 TypeError ⇒ 请求就地死。

## 病与裁形（本腿要做的事）

已落形状（`fb2fb802` 的 `:648` 那段 JS）＝死循环：钩子把 `cw.postMessage` 换成"调 `window.wispDispatch(message)`"，而 `wispDispatch` 的出口正是 `cw.postMessage` ⇒ 原生出口一次都不被调。
编排者 10-07 改裁**甲③**（票 `:150`，`a7f9781c`）：转发钩子**先把原生出口存成局部引用**，库自己产生的 RPC 帧**原样交回原生出口**，只有页面的裸信封才折进 `wispDispatch`。
本腿选**形ⓐ（再入守卫）**，理由写进 `impl.md` §1。

判据换成**行为模型**（新文件 `cmd/wisp/panel_transport_35r2_test.go`）：
一个内置的 JS 子集解释器**逐字执行**三处脚本——①`chromium.go:112` 库垫片原文、②`webview.go:462-478` 库桩原文（由 `Bind` 名册名按库自己的拼接式生成）、③**产码里那枚 `panelPostMessageForwardInit` 运行时字符串本身**（⛔ 不抄文本）——外加 `msgcb`/`callbinding` 按 `webview.go:131-168` 复刻，并带**再入深度帽**。
`panel_transport_35r1_test.go` 一字不动（既有到达性断言保留）。

## 本机可测性（照简报）

`go test ./cmd/wisp/` 载入即死 `exit status 0xc0000135`（缺原生 DLL）⇒ harness＝`go test -c -o /d/tmp/wisp35r2/panel.test.exe ./cmd/wisp/` ＋ 三枚 DLL（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）拷到 exe 旁 ＋ CWD＝`cmd/wisp`。
⛔ 不开任何真实 WebView2 窗（真窗那一发归 `35-v2`）；⛔ 不 push；⛔ 不动 `frontend/**`／`internal/panel/bridge.go`／任何 AC 框；证据件一律 `.txt`/`.md`（⛔ `.out`，根 `.gitignore:8` 全仓 `*.out`）。

（完）
