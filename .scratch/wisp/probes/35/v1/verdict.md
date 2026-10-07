# 35-v1 对抗验收判语（票 35 AC#6 `:52` / AC#7 `:63`）

腿＝`35-v1`（非实现者）。起手锚 commit `4b5c756a`；起手 HEAD 顶＝`ce6a4080`，date `Wed Oct 7 18:47:21 CST 2026`，porcelain=753（共享树既有脏）。
逐字读数全在 `logs/`（同名引用）。⛔ 本件不翻任何 AC 框、未动 `frontend/**`/`docs/**`/SLO/golden/thresholds、未 push。

---

## G1 定向突变（牙在不在）＝**成立（两发反形都红，恢复绿）**

四发实测（全部 DLL-harness：`go test -c` + 三枚原生 DLL 置 exe 旁 + CWD=`cmd/wisp`）：

| 发 | exe | 读数 | rc |
|---|---|---|---|
| 基线（原码） | panel-base.exe | `--- PASS: TestPagePostMessageEnvelopeReachesDispatchRawViaTransport (0.00s)` | 0 |
| 突变① `:671`→`_ = panelPostMessageForwardInit` | panel-mut1.exe | `panel_transport_35r1_test.go:209: AC#6 RED: installPanelTransport registered no page->wispDispatch forwarding hook …` | 1 |
| 突变② `:671`→`w.Init("/* no forwarding */")` | panel-mut2.exe | 同 `:209` 红句（装了钩子但钩子无三条词面 ⇒ forwardingInstalled 仍假） | 1 |
| 恢复 | panel-restored.exe | `--- PASS` | 0 |

还原三行尺（`logs/g1-md5-*.txt`、`logs/g1-diff-exitcode.txt`）：
md5-before＝`652ac6d6ff0db769f7bed7172a27dc72`；md5-after＝`652ac6d6ff0db769f7bed7172a27dc72`（逐字节同）；`git diff --exit-code -- cmd/wisp/panel_host_windows.go` 回显 `git-diff-exit-code=0`。

⚠ 具名形状：两发突变红都落在**正控门 #1**（`:209` 的 `forwardingInstalled()`），`:217`/`:224` 的到达性断言在两发里都没轮到执行。这不是盲区（判据要求"删转发⇒红"满足），但要说清：这枚用例对"转发语句在不在"敏感，靠的是**词面检测**（`chrome.webview`+`postMessage`+`window.wispDispatch` 三串都在 Init 文本里）＋ fake 的投递模型，而不是真浏览器投递。下一格攻的正是后者。

## G2 fake 保真度（对拉真库源码）＝**基本同形，但"转发已装"那一支比真的宽松——该用例在"帧形在真浏览器里可达"这一面恒真**

真库＝`github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808`（`go.mod:19` 在册，`cmd/wisp/panel_host_windows.go:64` 别名 `webview2` 直用，非 fork）。逐支对拉（原文在 `logs/lib-msgcb.txt`、`logs/lib-bind.txt`）：

1. **msgcb 解析**（`webview.go:139`，wiring `:103` `chromium.MessageCallback = w.msgcb`）：`json.Unmarshal` 失败 ⇒ `log.Printf` + return（**door 不叫**）；fake（测试 `:172-175`）同样不叫，只少了日志 ⇒ 对判据无差。
2. **未绑方法**：真库 `callbinding`（`webview.go:163-168`）对不在册名回 `(nil,nil)` ⇒ msgcb 走 **`window._rpc[0].resolve(null)`** 支（`webview.go:157`），**不是票 `:58-59` 写的 `reject(…)`**（reject 只在参数形状/序列化坏时，`:149/:153`）；fake（`:176-181`）计 `rpcDeadSlot`。两型 **door 都不叫** ⇒ 死亡结局同形；票面那句 reject 措辞与盘上库差一支（具名报，见末节）。
3. **桩帧形**：库注入的 `window[name]` 桩（`webview.go:460-480`）发 `window.external.invoke(JSON.stringify({id: seq, method: name, params: slice(arguments)}))`＝`{id,method:"wispDispatch",params:[raw]}`；fake `deliverPostMessage`（`:144-160`）造同形单层帧 ⇒ **形状本身一致**。
4. **★不宽松的一支**：库的 `window.external.invoke` 不是原生通道——是库自己在 document-created 注入的垫片（**`pkg/edge/chromium.go:112`：`window.external={invoke:s=>window.chrome.webview.postMessage(s)}`**），每次调用**动态**取 `chrome.webview.postMessage` 属性。而 `panelPostMessageForwardInit` 恰恰**重新赋值了这个属性**。⇒ 真浏览器里：页面 `postMessage(env)` → 包装器 → `wispDispatch(env)` → 桩 → `external.invoke(frame1)` → `postMessage(frame1)` **＝又进包装器** → `wispDispatch(frame1)` → … **静态读＝自-reference 环，一帧都到不了原生**；fake 则绕过这条回路、直接把干净单层帧喂 msgcb。⇒ **具名：这枚用例的绿证明的是"Go 侧注册＋下游到达"这一面（真且有牙）；在"这条转发在真 WebView2 里能不能把帧送出去"那一面，fake 比真的宽松，它的恒真面就在这条支上。**
5. `Init` 语义：库 `webview.Init`（`webview.go:435`）＝`browser.Init`＝`AddScriptToExecuteOnDocumentCreated`（`chromium.go:130-136`），对**之后创建的每个文档**生效；fake 的 `Init` 只存字符串。方向与语义一致。

## G3 真 WebView2 那一发＝**〔取不到・已给机制〕——本机今天没拿到"页面 postMessage→Go"的逐字到达，且静态读有自-reference 隐患**

- 既有载体不在这条边上：`git grep -l 'tags winlive' -- cmd/wisp` 六枚文件（`logs/winelive-files.txt`）全是 resident/geometry 族；四枚真窗用例（`TestPanelHostRealWindowHopAndLifecycle` `panel_host_windows_test.go:614`、AC13 `panel_resident_windows_test.go:314`、AC14 `:812/:855`）的入向全走 **`Eval("window.wispDispatch(…)")`**（`logs/grep-envelope.txt` `:301/:801/:805/:808/:864`＝票 :146 点名的 N6 形状，测试自己扮演页面）⇒ **没有任何一枚真窗载体调用过页面自己的 `chrome.webview.postMessage`**。
- 我写了台件（仅 `/d/tmp/wisp35v1/rig/`，⛔ 未进仓、零产码）：同库同版本、逐字 `forwardInit`＋逐字 `panel.ts:246` 信封、Go 门打印到达。两发：转发形 `logs/g3-rig.txt`＝`RIG: timeout after 12s, total door arrivals=0` + `RIG-NO-ARRIVAL`；**对照形（N6 直调 wispDispatch，r1 报过真机绿的那条）`logs/g3-rig-ctrl.txt`＝同样 0 到达** ⇒ 台件没复刻出产品载体的泵/Show 拓扑（wisp 的 panel 线程有自己的 pump 与 `DataPath`/`WindowOptions` 设定），**台件读数不能反推产品真窗**。
- 进程残骸尺：起前 msedgewebview2 计数＝12，两发跑完均＝12（`logs/g3-rig*.txt` 尾行），一窗一发、无残留。
- ⇒ **结论口径：今天这台机取不到那一发〔机制已给〕**；且 G2-4 的静态读（`chromium.go:112` 垫片与转发钩子争同一枚属性）意味着真机那发**有可能红的不是"没装钩子"而是"钩子把库自己的出口也接走了"**。AC#6 原文要的是"transport agreement"——这半枚**未经证实**，⛔ 不许把 fake 绿写成"浏览器里通"。

## G4 顺序证＝**成立（现量行号＋库语义；票 :144 行号漂移具名）**

- `bringUp`：`installPanelTransport` 调用在 `:408`；首个文档 `firstRoundTripLocked` 在 `:425`、`serveEntry` 在 `:427`（其内 `w.SetHtml` 现量 `:467`；`serveNotBuiltNoticeLocked` 的 SetHtml `:448`）⇒ **转发注册先于任何 SetHtml**。
- 存活语义：`Init`＝`AddScriptToExecuteOnDocumentCreated`（`chromium.go:130-136`），窗口已建后加入对**后续每个新文档**生效，活得过 `SetHtml` 重建；钩子自带 `__wispForwardInstalled` 防同文档重挂。
- 票 `:144` 引的 `:449/:468/:681` 现量＝`:448/:467/:671` 区（漂移 −1/−1/−10）；票 `:52/:63` 引的 `dispatchRaw @ :630` 现量＝`:677`（**+47**）；票 `:54` 引的 `w.Bind @ :405` 现量＝bind 进 `installPanelTransport` 体内 `:665`，`bringUp` 处剩 `:408` 调用行。

## G5 三条禁区＝**都在（三把尺全现量）**

1. `git diff --stat ce6a4080 -- frontend/` 与 `git diff --stat 10899b77 -- frontend/` 输出行数均＝**0**。
2. 名册：`git diff --stat ce6a4080 -- internal/panel/`＝**0**。盘上名册**六枚**（非简报说的四枚）：`bridge.go:42-45` `panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send` ＋ `:66-67` `config.get`/`config.set`，`knownComposerMethod`（`:148`）逐一列这六枚。逐字对拉 `panel.ts` 正控：`git show HEAD:frontend/src/lib/panel.ts` `:246`＝`sendRequest("panel.mode.request", { to });`（`logs/panelts-246.txt`）⇒ (a) 的"逐字"成立。
3. 无"答掉 RPC 解析器却不经 dispatchRaw"新支路：`_rpc|wispDispatch` 在 `cmd/wisp` 共 23 命中（`logs/grep-envelope.txt`），逐条看过——新增只在 `panel_host_windows.go:645/653/658`（注释+钩子文本）与测试文件自身，`panel_resident_windows_test.go` 那五处是既有 N6 形状；调用形状尺 `git grep -n 'm\.dispatchRaw('`＝**恰 1 枚**（`:666`，installPanelTransport 闭包内），构造/定义不计调用。

## 门禁与名册（本腿自跑）

- `sh scripts/d22scan.sh`：rc=0，末行 `d22scan: clean - no D22 ban violations … ban #8 cmd/=105`（`logs/d22scan.txt`）。
- `go test -count=1 ./internal/...` 起手与终态各一发：rc 均 1；逐名红集合 **`diff` 回显 `DIFF-EMPTY`**，六枚＝已知五枚（`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestC21DesignTokensFourWayAgree`/`TestC21TableColourRowsMatchTokensCSS`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`）＋间歇门 `TestResolvePerCallBudget`（出现，不计账）。
- `cmd/wisp` 裸 `go test` 的 `0xc0000135`：简报现量＋r1 起终两发同现；本腿**未重跑裸发**（harness `go test -c` 编译四发全 rc=0），按第 108 条标〔该读数非本腿现量〕。

## 总结论

**票 35 AC#6 今天不该翻。** 成立的部分：牙（两发反形都在 `:209` 红、恢复绿、md5/diff 还原坐实）、时序（转发先于首个文档且活得过重建）、三条禁区零违例、门禁零新红。**不成立/取不到的部分**：AC#6 原文要的是"页面↔宿主传输协议"这条**边**——fake 在"转发已装"支比真库宽松（G2-4：库的 `external.invoke` 出口与被钩子覆写的 `postMessage` 是同一枚属性，静态读为自-reference 环），真机那一发本腿〔取不到・已给机制〕。**翻框前置条件**＝一发真机读数：用 `bringUp` 的真实载体（含其泵/Show 拓扑）让**页面**调 `postMessage`、Go 门逐字收到信封；若红，转发形应改为"先存 `var native = cw.postMessage.bind(cw)` 再包装"之类仍属 Go 侧一文件的形（不触页面、不改名册）。

**AC#7 的 (d) 三支守卫红文案今天覆盖 0/3**（现量：测试文件内 `不是面板 composer|伪造|串台|缺少 requestId` 命中数＝**0**）；(a)(b)(c) 三支本腿复核**成立**——(a) panel.ts:246 逐字对拉过，(b) 入口是 `installPanelTransport`+fake 投递、零直调 `dispatchRaw`、零 `Eval(wispDispatch)`，(c) `:229-231` 断言 requestId/method/source 到达下游。**简报口径"r1 覆盖 (a)(b)(c)、(d) 零支"＝与我现量一致，非照抄。**

## 简报与盘面冲突（具名报回，一律按盘面做）

1. 简报 G5②说"四个在册方法名"——盘上名册**六枚**（`bridge.go:41-68`＋`:148`；票 `:73` 亦写"6 枚名册"）。我按六枚对。
2. 票 `:58-59` 说未注册信封被 msgcb 以 `window._rpc[0].reject(…)` 应答——盘上库对**未绑方法**走的是 `resolve(null)` 支（`webview.go:147-148`+`:157`），reject 只在参数形状坏时。死亡结局（door 不叫）不变，判据不受影响；票句与盘面差一支，记录在此。
3. 简报说"票面指 `webview.go:139-160`"——盘面 msgcb 起 `:139` 至 `:161`（含收括），wiring `:103` ✔。
4. 简报 G1 说红那发含 "`panel_transport_35r1_test.go:2xx:` 那行"——现量红行＝`:209`（两发同线），✔。
5. 本腿全程未见需要顶回的其他处；r1 `impl.md` 各读数我逐把复跑，除上面第 2 条（它没写、我补量）外无冲突。

（完；证据 logs/ 23 枚同 commit 入库。）
