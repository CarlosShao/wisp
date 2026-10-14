# 253-v3 — 锚、机制链自核、运行窗口与足迹（验收腿；非实现者）

**身份**：验收腿 `253-v3`（`agent=253-v3`）。被裁对象＝落地腿 `253-r6`（产码 `2298f7d2`；交付 `89d29fce`）。
**判语**：见 `docs/evidence/s1/ticket253-253-v3-accept-v1.md`。⛔ 本腿未翻任何 AC 框、未动 `Status:`／票名、零 push。

## 1. 被裁件的盘上读数（我现跑）

| 件 | sha256（我现跑） | 与 r6 自报 |
|---|---|---|
| `cmd/wisp/panel_control_windows.go`（改后，HEAD blob `84ff38c9…`） | `2fe49342912fa4ab054880f2cf914ad9ed2b0e9dacaacef8f8ccd6a650469ede` | 逐字相同 |
| `cmd/wisp/panel_unbound_253r6_windows_test.go`（270 行） | `53ed0947406dd5a80285f2570fdeb95a6be7c22a296bfc1dcc3af46a10e87e26` | 逐字相同 |
| 改前 blob（`git show '2298f7d2^:…'`，反控用） | `ae714aeeafe2f24c584f2cb4eb14a36c3a75f252fc3340cbc969d88087de1e46` | 逐字相同 |

## 2. 机制链自核（行号＝现读 HEAD）

- `browser.MessageCallback = c.msgcb`＝`panel_control_windows.go:339`；库侧 `edge.Chromium.MessageReceived` → `MessageCallback`（模块缓存 `pkg/edge/chromium.go:233-241`）⇒ 管道已是本仓代码（299-r1 起），票面「静默面属第三方」在生产路径上过期——成立。
- `msgcb`＝`:711-755`：invalid-JSON 支 `:713-720`（slog.Warn+return）；无绑定门支 `:728-740`（`boundDoor` `:704-709` 先判，slog.Warn＋`window._rpc[id].reject(驳回句)`）；有绑定门 `:742-754`＝库原路径一字未改。
- `callbinding`＝`:757-812`，无门仍回 `(nil, nil)`（`:761-763`，对上游忠实；库 `webview.go:139-168` 我现读逐行比对，无门支与库同形）。
- **改前静默 null 的必要条件（自核，非引件）**：封套是合法 JSON ⇒ `json.Unmarshal` 成功；`method` 不在 `bindings` ⇒ `callbinding` 回 `(nil,nil)`、`err==nil` ⇒ 落到 `json.Marshal(nil)=="null"` 支 ⇒ `Eval("…resolve(null); …")`；该支沿途**无任何日志调用**。⇒ 页面（有槽）拿到 `null`、Go 侧零记录。反控两发自跑证实（§4）。
- **三分界（逐点）**：
  - 旧门 `window.external.invoke(envelope)`：35/303 钩子（`panel_host_windows.go:813-832`）把 `ex.invoke` 包成 `inside=true` 的转发（`nativeInvoke`＝库 `chromium.go:112` 的 `s=>chrome.webview.postMessage(s)`，属性在调用时读）⇒ `inside==true` ⇒ `native.call(cw, msg)` 逐字节出原生口 ⇒ 原封到达 `msgcb`。**本次落点**。
  - 折入路 `chrome.webview.postMessage(envelope)`（页面自己的、非 invoke 内）：`inside==false` 且 `wispDispatch` 可调用 ⇒ 折成 `{id,method:"wispDispatch",params:[封套]}` ⇒ 绑定门 ⇒ `dispatchRaw` → composer（35/303 已覆盖，有声）。my (D) 正控与 r6 的 folded 读数两向核过。
  - P3 未折入世界：`typeof window.wispDispatch !== "function"` ⇒ 同走原生口 ⇒ 同落点。我自种真窗一发（(E)）。
- 绑定门全集（生产）＝`wispDispatch`（`panel_host_windows.go:842`）＋瞬态 `wispProbeRT`（`:892`）；`panel.mode.request`／`legacy.echo` 均不在其中 ⇒ 旧封套解成 `{id:0,method:"panel.mode.request",params:nil}` 无门可答。
- ⚠ **一条与 r6 件不符的通道读数（本腿更正）**：stdlib `slog.SetDefault` 会把 `log` 包默认输出桥进当前 default handler（Go 源 `log/slog/logger.go:62-73` `handlerWriter`），resident 侧 `installLogSink` 已把 tee 设为 default（`resident_windows.go:101` → `logsink.go:153-161`）⇒ 改前 `log.Printf("invalid RPC message…")` **并非「非 resident 的 slog/JSONL 通道」**，它会被桥进 JSONL（INFO、无 attrs）。我的 (C) 改前读数当场证（capture 收到该行）。

## 3. 运行配方与窗口

- 环境：`PATH=/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\253v3'`（反斜杠形）；每枚真窗用例**单发独立 invocation**（-count=1 -v），无两窗连发。
- 窗口：起手锚 `9065c659`（15:2x 现读）→ 执行窗内 HEAD 前移 3 笔（`f260b8a6`／`037bfb7f`／`125e2255`，全为编排者件）；`git log --name-only 9065c659..HEAD` 核过**零 `cmd/wisp`／`internal/panel` 重叠**⇒ 被裁码窗口内无人动。
- 五次真窗一发（各一枚单窗用例）：

| # | 用例 | 树 | rc | raw |
|---|---|---|---|---|
| 1 | `Test253r6TheSilentDropBranchIsLoud`（r6 件，我复跑） | 改后 | 0 | `logs/05-…` |
| 2 | `Test253v3NewShapesOnTheRealHost`（我五形四读数版） | 改后 | 0 | `logs/06-…` |
| 3 | 同上（四形版） | 改前 overlay | 1 | `logs/07-reverse-control-run.txt` |
| 4 | 同上（五形版，＋P3） | 改后 | 0 | `logs/08-…` |
| 5 | 同上（五形版） | 改前 overlay | 1 | `logs/09-…` |

- 突变共 2 窗（加 overlay 测试件；换产码为改前 blob），每窗跑完**逐字节还原**：产码 sha256 回 `2fe49342…`、`git hash-object` == HEAD blob `84ff38c9…`、`git status --porcelain cmd/wisp internal/panel` 空、overlay 测试件移出仓外（master 留本目录）。
- 我自造件：`panel_unbound_253v3_windows_test.go`（overlay 用；master 在本目录，⚠ 未提交进 `cmd/wisp`）。

## 4. 反控逐字（我自己的形）

- 改前 overlay 下：`(A) id=977 legacy.echo ⇒ "RESOLVED(null)"`、`(E) 门不可调用 ⇒ "RESOLVED(null)"`、`(B) ⇒ "getterHit=true"（无记录）`、`(C) ⇒ capture 收到桥过来的 "invalid RPC message: …"`、`(D) 正控 ⇒ "RESOLVED('')"`；全集 rc=1。
- 改后：`(A) ⇒ "REJECTED('面板请求被拒绝 [legacy.echo]：…宿主已记录本次丢弃')"`＋slog `method=legacy.echo id=977 params=1`；`(B) ⇒ "getterHit=true"`＋slog（页面仍拿不到东西＝结构边界）；`(C) ⇒ slog.Warn "…not an RPC frame… err=invalid character 'v' after top-level value bytes=34"`；`(D) ⇒ "RESOLVED('')"＋handler 恰 1＋零新告警`；`(E) ⇒ REJECTED + 恰多一条 slog`；全集 rc=0。

## 5. 盘上件

- 本文件；`logs/01-gates.txt`（build/vet/gofumpt/sha256/d22scan/长度门）·`02-families.txt`（headless 15）·`03-realwindow-postfix.txt`·`04-reverse-control.txt`（提交）；raw 全文（`01-gates-build-vet.txt`／`02-headless-family15.txt`／`03-d22scan.txt`／`04-length-gate.txt`／`05…09`）留盘未提交，文件名如上表。
- 足迹：只 1 笔 commit（evidence＋本目录；显式 pathspec）；⛔ 未 push。
