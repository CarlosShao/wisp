# 253-r3 — 只读预检：票 253 AC#1（可达性钉改问能力 + "只 postMessage 不 wispDispatch"正控）备料

- 时刻：2026-10-02 22:31 +0800；锚点 HEAD `0ce6cb91`（dev）。收笔锚另记于文末。
- 性质：只读。⛔ 未跑任何 Go 命令（`go build/test/vet` 零发）；⛔ 票面勾选框零碰；⛔ 产码零改；⛔ `frontend/**`／`design/**` 零读（grep 根显式 = `internal/panel`、`cmd/wisp`、`.scratch/wisp/issues`、`docs/reports`、上游 module cache 只读引用）。
- 写点唯一 = 本文件。占位词 0（每格都是读数或〔量不到〕）。

---

## 一问：词面钉现状 — `postMessage` 在 internal/panel + cmd/wisp 测试里的每一枚

grep 逐枚（`grep -rn -i "postMessage" internal/panel cmd/wisp --include=*_test.go`，本腿复跑）：

### 1a. internal/panel 下命中 9 行，按"钉什么"分类

**钉"页面必须走 postMessage 词面"的（词面尺，认字符串）：**

1. `internal/panel/bridge_test.go:157` — `strings.Contains(text, "window.postMessage(")`（还有 `new WebSocket`/`fetch(`），断言原文：
   `t.Error("frontend opened a second channel besides the WebView2 host bridge - ticket 92 forbids it")`
   ⇒ 读 `frontend/src/lib/panel.ts` 的**文本**，禁第二通道。词面。
2. `internal/panel/bridge_test.go:161-163` — `if n := strings.Count(text, "bridge.postMessage"); n != 2 {`
   断言原文：`t.Errorf("postMessage call sites = %d, want 2 (the approval request and sendRequest) - a third one means a second envelope was invented", n)`
   ⇒ **这是钉"postMessage 调用点恰两枚"的核心词面钉**（TestFrontendComposerRequestsMatchTheEnvelope，`:131`）。
3. `internal/panel/bridge_test.go:20`（文件头注释）— `bridge.postMessage(JSON.stringify({method:"panel.approval.request", correlationId, outcome}))`，仅注释。非断言。

**钉"调用点在哪枚文件/有几枚"的结构词面尺（同一族，认 `\.\s*postMessage\s*\(` 正则）：**

4. `internal/panel/composer_test.go:381` — `var hostBridgeCallRe = regexp.MustCompile(`+"`"+`\.\s*postMessage\s*\(`+"`"+`)`，注释逐字：`hostBridgeCallRe matches a CALL into a postMessage-style host bridge. The leading dot is required, so panel.ts's own interface declaration "postMessage(message: string): void;" is not counted as a call site.`
   ⇒ **票 253 票面 §现量-1 点名的那把尺（composer_test.go:381〔已验，非待验〕）就是它**。
5. `internal/panel/composer_test.go:460-465`（scanRendererHostDoors 内）— `if hostBridgeCallRe.MatchString(line) { rep.callSites = append(...); if !inPanelLib { rep.outside = append(...) } }`
6. `internal/panel/composer_test.go:509-512`（TestTheRendererHoldsExactlyOneDoorToTheHost）— `if len(rep.callSites) != 2 { t.Errorf("host call sites in frontend/src = %d, want 2 (the approval request and sendRequest)..."` ⇒ 词面（正则认调用形），断言的是 **frontend/src 树**，不是 Go 侧收没收到。
7. `internal/panel/composer_test.go:497`（头注释）与 `:548`、`:552`（TestPlantedRendererDoorShapesGoRed 的 plantA/plantB 文本：`window.chrome.webview.postMessage(JSON.stringify({ method: "panel.mode.set", to: "auto_approve" }))`）⇒ 这枚是**正控**：把"postMessage 词面+禁路由"种进假树要求尺红。注意它证明的是**词面尺有牙**，不是"Go 侧收到"。

**与入向无关的（Win32 PostMessageW，别混进撞钉清单）：**

8. `cmd/wisp/panel_resident_windows_test.go:596-598` — `t33r9NailPostMessage = pnlModUser32.NewProc("PostMessageW")`；`:669` 调用发 WM_CLOSE。这是 33-r9 的关窗钉，**与页面→Go 通道零关系**。
9. `cmd/wisp/panel_host_windows_test.go:226` — 注释：`records that the router reached a handler for a real postMessage`（recordingModeHandler 头注释）。只提词，无断言。

### 1b. 判读：今天没有任何一枚钉在钉"wispDispatch"

全仓 `internal/panel` + `cmd/wisp` 的 `*_test.go` 里 **`wispDispatch` 只在 cmd/wisp/panel_resident_windows_test.go 出现**（`:199` 注释、`:301`、`:801`、`:805`、`:808`、`:864`，全部是测试**自己用**这扇门回话，见二问）。⇒ 生产绑定 `cmd/wisp/panel_host_windows.go:82` `const panelDispatchBinding = "wispDispatch"` / `:319` `w.Bind(panelDispatchBinding, ...)` **零枚测试断言它存在**；既有"可达性"阅读全部落在 1a 的词面尺上。**票面 §现量-1 的判读成立**："这条链路可达"那一格今天由 `composer_test.go:381` 的正则 + `bridge_test.go:161` 的计数在判，两把都认 JS 词面、不认 Go 侧行为。

---

## 二问：能力形正控的接缝 — 票 33 "页面自己报回 REPLIED" 的定式 + 现成假页面

### 2a. REPLIED 定式（file:line + 断言原文）

唯一实例：`cmd/wisp/panel_resident_windows_test.go`，票 33 AC#14 nail 1（TestAC14AwaitedBindingReplyReachesThePage，`:812`）。

- 造话（`ac14AwaitJS`，`:786-810`）：页面侧脚本逐发 `await Promise.race([Promise.resolve(window.wispDispatch(ENVS[i])).then(function(v){ return 'REPLIED'; }), new Promise(function(res){ setTimeout(function(){ res('TIMEOUT-2S'); }, 2000); })])`（`:804-806` 原文），再把结果用**同一扇门**报回：`window.wispDispatch('{"method":"panel.mode.request","requestId":"ac14r-i",...,"to":"<got>"}')`（`:808`）。`:801` 先发 beacon（requestId=`ac14r-beacon`、to=`SCRIPT-RAN`）证明脚本跑过。
- 判读（`:826-828` 原文）：`if strings.Contains(joined, "NOT_RESOLVED") || strings.Contains(joined, "TIMEOUT-2S") || strings.Contains(joined, "THREW") || !strings.Contains(joined, "REPLIED") { t.Errorf("the awaited JS binding reply did not reach the page: the page reported %q. Go did receive the calls (%d reached the handler)..."` ＋配对断言 `if calls != 3`（`:830-832`，Go 侧 handler 必须也被到达 3 次——两半必须互证，"one instrument may not stand in for both"）。
- **这就是票 253 AC#1 说的"由页面自己报回"的形**：判据载体＝requestId 前缀区分"真请求"（`ac14-`）与"报告"（`ac14r-`），`countRequests`（`:837-845`）+ `awaitReport`（`:215-229`，超时＝t.Fatalf 判测量失败）。

同文件另两枚同族（页面自报，但不是 await 形）：`probeJS`（`:295-303`，AC#13，`window.wispDispatch(reportJSEnv("ac13-probe", "b"))` 位图回话）与 `TestAC14GoSideEvalPushReachesThePage`（`:855-871`，`:864`）。

### 2b. 现成的"假页面"零件（可复用清单）

- **recordingModeHandler**（`cmd/wisp/panel_host_windows_test.go:227-243`）：`type recordingModeHandler struct { mu sync.Mutex; reqs []panel.ComposerRequest }`，`HandleModeRequest` 追加并收下，`count()`/`all()`（resident `:206-210`）读回。**这是唯一的"假 handler"**，`startPanelForTest`（resident `:59-67`）拿它组装真装配：`disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}`。
- **真宿主装配架**（不是假页面、是"页面自己会说话"的现成环境）：
  - `startPanelForTest`（resident `:59-67`）+ `showAndWait`（`:104-118`）＋ `evalOnPanelThread`（`:166-189`）＝向活文档 Eval 任意 JS 的现成缝。
  - `reportJSEnv`（resident `:250-253`）：拼一枚带回的 composer 信封（method=`panel.mode.request`，source=`panel-composer`，payload 进 `to`）。
  - `firstRoundTripLocked`（`cmd/wisp/panel_host_windows.go:570-615`）＝生产侧已有的"绑一枚一次性探测 binding `wispProbeRT`＋SetHtml 探测页＋泵到回话"的样板——**能力形新尺若需要"第二枚 binding 收报"，这里是现成形状**（但注意 AC#13 已钉"探测页不得最后留在屏上"，重排/复用要过 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`，resident `:314-332`）。
  - `TestPanelHostRealWindowHopAndLifecycle`（panel_host_windows_test.go `:507` 起，`:554` `mgr.dispatchRaw(ctx, env)`）＝在真窗上从 Go 侧直接驱 router 的缝（H3 半），配 recordingModeHandler 计数。
- internal/panel 侧的"假页面"**不存在**：`composer_dispatch_test.go` 的 hop33ModeWriter/hop33Capture/hop33RefusingHandler（`:52-107`）与 config_route_248_test.go 的 cfg248Leg（`:26-56`）都是 **handler 侧替身**，只证 `ComposerDispatch.Handle` 行为；该包没有任何持有 WebView/文档句柄的桩。⇒ **能力形可达性钉（页面→Go→页面闭环）只能落 cmd/wisp 的 resident/host 测试面**；internal/panel 里能落的最多是"词面正控"那一半。

---

## 三问：撞钉清单 — 新尺（能力型可达性）与既有断言谁互斥

读法：逐枚给断言内部，判"新尺落地时谁会红、红的归因是谁"。⛔ 票 248 两枚已知红（见 ⓕ）不算新尺凭据、不许修。

### ⓐ composer_test.go 词面双尺（最可能被新正控洗红的邻居）

- `TestTheRendererHoldsExactlyOneDoorToTheHost`（`:502-526`）：断言 `len(rep.callSites) != 2`（`:509`）、outside=0（`:505`）、assembly=0（`:513`）、computed=0（`:517`）、unknown=0（`:521`）。**新尺不动它**——它的射程是 frontend/src 文本。⚠ 但 AC#1 的正控若做成"往真树种一枚只-postMessage 文件"，绝不能种进 `frontend/src`（两层禁令本就禁），必须像 `TestPlantedRendererDoorShapesGoRed`（`:533`）那样种 **t.TempDir 副本**，否则这枚照 `:505`（outside=0）红，红因＝新正控污染真树，不是缺陷本体。**归因写死：正控载体一律仓外。**
- `TestFrontendComposerRequestsMatchTheEnvelope`（bridge_test.go `:131-165`）：`bridge.postMessage` 计数==2（`:161`）。同样只读 panel.ts 文本；新尺不触。⚠ 它与 `:157` 的"第二通道"禁令意味着**正控文案里连 `window.postMessage(` 字样都只许活在仓外/测试字符串里**。
- `composerRouteLiterals()`（`:409-419`）闭集五名（四 panel.*＋panel.approval.request）：新尺不得往真树加任何 `panel.*` 字面量，否则 `:521` unknown 红。

### ⓑ 票 33/145/146 名册差集尺（AC#2 将来要动的，本腿只登记形状）

- **票 33 族（capability 化已完成的先例）**：`TestPanelHostIsAttachedAndNamesTheWindowHops`（composer_dispatch_test.go `:434-540`）——**从词面尺改成能力尺的现成范本**（头注释 `:397-433` 逐字记录了 A386/A489 两次反转与"正控 A/B＋负控 C＋反向正控 D"四件套）。AC#1 的新尺可直接镜像它的 teeth 结构。零互斥。
- **票 181/33 名册双向尺**：`TestGitDimensionHasNoModelCallableTool`（internal/panel/git_test.go `:362-390`）— `whitelistMethodsFromSource`（`:407-423`，用 `:394` `panelMethodRe = regexp.MustCompile(`+"`"+`"(panel\.[a-z0-9_.-]+)"`+"`"+`)` 抽 bridge.go 字面）要求与四常量 `equalStrings`（`:385`）。**这就是票 253 票面 §现量-2 说的"panelMethodRe（git_test.go:394）"〔已验〕**。⚠ 正控 `TestPlantedGitToolShapesGoRed`（`:484-532`）在 `:509` 断 `len(whitelistMethodsFromSource(bridge.go)) != 4` want 4 —— **这枚今天就在锁死"bridge.go 恰四枚 panel.*"**。票 253 AC#2 若把 `config.get/set` 纳入某把尺的射程，这枚的 want-4 与报错文案必须同 commit 具名解冻（先例：票 248/台账 A487 对 l2 冻结锚的解冻形，以及 A487-J2 把 git_test.go 那两枚钉的改写已裁过一次"甲＋乙合并"）。**本腿不裁、只点名。**
- **票 145 差集尺**：`pump_test.go:123` 与 `:291` 两枚 `strings.Join(sortedCopy(keys),",") != "composer,generatedAt,pending,results"`（快照四键恒等钉，A388 具名解冻史在册）。AC#1 新尺**不碰快照键**，零互斥；AC#2 若走快照面则必撞，须同款"两行具名＋正控"解冻。
- **票 146**：`internal/agent/approval/ticket146_liveapprovals_backing_test.go` ＋ `pending_read.go` 的 LiveApprovals 深浅拷贝尺。射程在 approval 包的出向读面，与入向可达性**零交集**。

### ⓒ 冻结件射程（l2_grant_boundary_test.go —— ⛔ 一字不动区）

- `TestAnsweredPanelRoutesCarryNoApprovalDecision`（`:1229-1306`）：`answered := sortedSet(pkg.answered)` 全部必须过 `knownComposerMethod`（`:1240-1242`）＋`Method*` 常量双向（`:1293-1303`）。**能力形新尺若要在 Go 侧新增任何"回答"面（比如让宿主对不认识形状回错误），只要不进 `knownComposerMethod` 名册就零互斥**；反之若有人图省事把"具名拒收"做成第六枚方法名 ⇒ `:1301` 红（"a route written straight into the guard's case list"），红因＝把 AC#3 做成了新方法，**这是该红的**。
- `guardRosterOf`（`:2030-`）派生名册（A487 解冻后的能力形）：要求 guard case 与 Method* 常量双向相等（`:2093` default 分支即 problem、`:2099` 非常量 label 即 problem）。`config.get/set` 走常量（bridge.go `:66-67`）⇒ 名册尺**看不见但不违规**——这正是票 253 §现量-2 的洞：**两把名册尺的正则都只认 `panel.` 前缀字面**（`composer_test.go:394` `routeLiteralRe = "panel\.[A-Za-z.]+"`；`git_test.go:394` `"(panel\.[a-z0-9_.-]+)"`），而 `MethodConfigGet/Set` 是**常量**且 `knownComposerMethod` 已答 ⇒ `TestAnsweredPanelRoutesCarryNoApprovalDecision` 的双向对账**今天绿**（answered 集里六个名字全部有常量、全被 guard 答），名册尺对"页面侧能拼出什么"全盲。⇒ **AC#2 的"种一枚没进名册的假腿今天绿"缺陷本体 confirmed**：往真 bridge.go 塞 `MethodX = "frobnicate"`（非 panel. 前缀、非常量互斥面之外任何名）会被 `guardRosterOf` 抓（label 非常量→problem）——但塞**常量** `"config.anything"` 进 guard case + 常量声明成对加 ⇒ 双向仍相等、`git_test.go:385` 的 panelMethodRe 抽不到它、`routeLiteralRe` 也抽不到 ⇒ **全绿穿过，洞属实**。
- `TestRealGuardRefusesEveryAssemblableApprovalRouteName`（`:1530`）与 `grantRoutePrefixes`（`:1316-1319`）：**prefix 集 8 枚全 `panel.*`**；`config.*` 不在 grid 里 ⇒ approval 词表对 config 名全盲（票面 §现量-2 的另一面）。改它＝冻结件编辑，⛔ 本票不得碰，只能单立尺。
- `TestPlantedGrantWiringGoesRedInASnapshot`（`:2168`）：plant B/E 都挂在 `guardAnchor` 上；新尺不种 l2 面 ⇒ 零互斥。

### ⓓ 宿主/ Resident 面的既有钉（AC#1 新尺的落点邻居）

- `TestAC14AwaitedBindingReplyReachesThePage`（resident `:812`）：**新尺与它共享门与 handler**。若新尺复用 `recordingModeHandler`，必须用**不同 requestId 前缀**（`countRequests` 按 `strings.HasPrefix(req.RequestID, prefix)` 分桶，`:837-845`），否则两枚尺互相计数污染。红因归新尺自己写脏桶名。
- `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（resident `:314`）：新尺若 SetHtml 自己的探测页 ⇒ `:330` 红（冷启动结束于探测页）。**归因＝新尺违反 AC#13 次序裁**，修法走"Eval 注入进现有活文档"（`evalOnPanelThread` 现成）而不是再 SetHtml。
- `TestPanelHostRealWindowHopAndLifecycle`（host_windows_test `:507`）：`:558` `h.count() != 1` 精确计数。新尺若在**同一**装配上多塞信封 ⇒ 这枚红。归因＝共享 handler 桶污染，同 ⓓ-1 处置。
- `firstRoundTripLocked` 的 `wispProbeRT`（host_windows.go `:582`）：第二枚 binding 已存在且是生产码。AC#1 判据"问能力"若走"页面对**不存在的方法**调用必得 reject/可读错误"那半，注意上游 `callbinding`（上游 webview.go `:162-169`）：**未注册方法 `f, ok := w.bindings[d.Method]; if !ok { return nil, nil }` ⇒ resolve(null) 静默**；而**错误路径**（binding 返回 error）走 `window._rpc[id].reject(...)`（`:145-147`）⇒ **"页面拿到可读错误"半格在现有依赖上是可实现的**（让 wispDispatch 的 Go 闭包返回 error），"旧封套 postMessage 直发"那半（页面不包 {id,method,params} 直发裸 JSON）则走 `msgcb` 的 `log.Printf("invalid RPC message")`（`:139-142`）＝**stderr 一行、Promise 悬死**——票面 §现量-3 "webview.go:162-168 resolve(null) 零审计"〔已验，行号实为 `:139-169`，票面行号漂 23 行〕。⇒ AC#3 的 Go 侧落点在 wispDispatch 闭包内（`panel_host_windows.go:319-322` 现在吞 error：`reply, _ := m.dispatchRaw(ctx, raw)`），上游零改动可达。
- `TestAC9ComposerDispatchHasAProductionCaller`（panel_inbound_33_test.go `:182-198`）：AST 找 `<var>.Handle` 调用点 ≥1。新尺若新增第二扇"直驱 Handle"的门（比如绕过 dispatchRaw）会让枚数**涨但不红**（断言是 `len(sites)==0` 才红）⇒ 零互斥；但 `main.go` usage 断言（`:192-197`）意味着任何新 CLI 面要进 usage 块。

### ⓔ 票 248 名下两枚已知红（⛔ 不算凭据、不许修，只登记隔离面）

台账 A533（`docs/reports/pending-and-issues.md:10830`）逐名：`internal/panel` 整包 4 FAIL＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`（末枚＝Q-52 已知常红），全部**页面契约族**（Go 快照长胖、frontend 类型没跟上／design/assets 路径被别队挪走），**与票 253 三格零交集**。⚠ 其中 ⓐ-1 的 `TestComposerContractTypesMatchFrontend`（composer_test.go `:48-81`）**在册常红** ⇒ 它对"新尺落地后谁红"的读数**不可用作对照**（它红是既知的，红上加红无法归因）；新尺的对照面只能选在册绿的那批。

### ⓕ 其余擦边，逐枚判不撞

- `TestUnlistedMethodNameIsRefusedAndAudited` / `TestRosterMismatchBackstopRefusesInsteadOfAccepting` / `TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped`（composer_dispatch_test.go `:201/:231/:249`）：router 行为尺，新尺不动 router ⇒ 零互斥。⚠ `TestAcceptedRequestInventsNoReplyLine`（`:157-169`）钉"accepted 路径 reply==''"——AC#3 若让"收到但形状不认识"回错误文本，走的是 **refusal 路径**（ErrComposerRequest 系），不碰 accepted 空 reply 判 ⇒ 零互斥，但 AC#3 实现时这条必须抄进派单防误伤。
- `TestAC1LockedFamilyFieldsAreRefusedBeforeTheLeg`（config_route_248_test.go `:215`）等 248 家族：settings 面，与可达性/名册射程两格零交。
- `TestAC1InboundLegAnswersSettingsRouteEndToEnd`（cmd/wisp/panel_config_248_test.go `:413-452`）：stdin 腿端到端，`config.get/set` **词面**在测试字符串里（`:416-420`）——AC#2 新尺若扩名册尺射程去数"非 panel. 前缀方法名"，这枚测试文件里的字面量会被新尺数到吗？**取决于新尺扫描根**：现有两把尺的根都是产码文件（bridge.go / frontend/src），不扫 `_test.go`（`gitToolNamesUnder` 的正控 `:490-495` 明确 `_test.go` 不算）⇒ 零互斥，照抄这个"测试文件不进分母"的规矩即可。
- `TestComposerRenderFixtureTellsTheTruth`（composer_test.go `:677`）／`TestFrontendHasNoEmoji`（frontend_hygiene_test.go `:246`）／tokens_fourway（`:439`）：frontend 面，零交。
- `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`（panel_host_gate_test.go `:77`）：embed 面，零交。
- `TestCleanCheckoutBuilds_AC11`（`:385`）：`git archive HEAD` ⇒ 新尺测试文件入库后它拷的树变大，判据只问 build 成败 ⇒ 零互斥。

---

## 四、量不到的格子（本腿禁 Go 命令，以下〔量不到〕如实登记）

1. 〔量不到〕两把词面尺 + 新能力尺**同跑一遍**的 PASS/FAIL 名册——需要 `go test ./internal/panel/ ./cmd/wisp/`；只读腿禁跑。落地腿起手应复跑在册绿基线（panel 在册 4 红见 ⓔ）。
2. 〔量不到〕上游 `callbinding` resolve(null) 的**运行时**形态（页面侧 Promise 实际悬多久、reject 是否可达）——静态读数已给（三问 ⓓ），动态验证要真窗（`TestPanelHostRealWindowHopAndLifecycle` 那一档），归票 253 落地/验收腿。
3. 〔量不到〕`frontend/src/lib/panel.ts` 今天**实际**有几个 postMessage 调用点／panel.* 字面量——两把尺的绿依赖它，但读 `frontend/**` 是本腿禁区，只确认尺的断言文本（want 2 / 五名闭集）。
4. 〔量不到〕A487 解冻先例之外，`git_test.go:385/:509` 两枚 want-4 锚**是否可依同款程序具名解冻**——程序存在（A487 五样齐），批准是编排者的事，本腿不裁（⛔ 不建议修法选形）。

## 五、推翻前人哪句

1. **票 253 票面 §现量-3 的行号**："依赖模块 webview.go:162-168〔待验〕" → 实际静默面是 **webview.go `:139-169`**（`msgcb` 的 invalid-JSON `log.Printf` `:139-142` ＋ `callbinding` 的 `if !ok { return nil, nil }` `:162-169`）；票面引的 162-168 恰是 callbinding 本体，**漏了上半的 invalid-JSON 分支**——旧封套若连 {id,method,params} 形都不是，走的是 `:139-142`（log 后 return，Promise 永悬），不是 resolve(null)。两条静默路都在，票面结论不变、行号补全。
2. **票 253 票面 §现量-2 的"掉在 routeLiteralRe（composer_test.go:394〔待验〕）与 panelMethodRe（cmd/wisp/git_test.go:394〔待验〕）之外"**：**panelMethodRe 不在 cmd/wisp**，在 **`internal/panel/git_test.go:394`**（cmd/wisp 下零命中，grep 已验）。票面结论（config.get/set 不带前缀、两把名册尺抽不到）**成立不变**，出处行号修正。
3. **前四程（253-r1/r1b/r2）无产码无 commit**（台账 A542/A544 在册，盘上 `.scratch/wisp/probes/253/r1/` 空目录、r2 无目录）——本腿无"前人读数"可推翻，以上两条只对票面票面文本。

## 收笔锚

- 2026-10-02 22:4x +0800；HEAD `0ce6cb91`（写盘期间无新提交进来；porcelain 增量仅本文件与既有他票脏面，本腿零产码零冻结件触碰）。
