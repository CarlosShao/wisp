# 33-a3 — 票 33 AC#14 与票 35/145 的关系：三处原文并排（逐字，⛔ 概括；锚 c315980b）

派单："上交给你答不了的一件"——把 AC#14 那句'最后一公里'原文与票 35/145 的关系**逐字读回**并排放，谁覆盖谁由编排者裁。本腿⛔ 裁下一步（翻 33 格／另立接线票皆不归我）。

---
## 原文一：票 33 `AC#14`（`.scratch/wisp/issues/33-panel-host-c27.md:42`，逐字）

- [ ] **AC#14（10-01 11:4x 编排者追加，来路＝只读普查 `33-a2` §⑦-H＋我 11:44 亲自复认的依赖源码；账 `A497`）：本票的宿主必须说清"Go→页面"这一跳由谁投递——今天的答案是"没有人"。** 现读凭据（模块 `github.com/jchv/go-webview2@.../webview.go`，机主读数非转述）：`Dispatch(f)` 在 `:443-448` 只做两件事＝把闭包追加进 `w.dispatchq`（`:445`）＋`PostThreadMessageW(w.mainthread, WMApp)`；而全文件里 `dispatchq` 唯一的取出点在 `Run()` 内部（`:351` 起，`:362-366` 那个 `if msg.Message == w32.WMApp` 分支）⇒ 不调 `Run()`，队列永远不空。本仓宿主对 `Run()`／`Dispatch()` 的调用者＝**0 枚**（`33-a2` R17/R18）。⚠ 这条为什么不是实现细节：`docs/PLAN.md:1367` 的 **C17** 逐字要求「Go→前端事件推送；回复必须按 correlationId 路由」——"推送"这一半今天没有投递者；特别是页面调 Go 的绑定并 `await` 回话那一跳，库内部就是走 `Dispatch`，绕不开（`Eval` 确实绕得开队列：`:439-440` 直接 `w.browser.Eval(js)`，但那是"Go 主动塞 JS"，不是"把绑定的返回值送回页面"）。

（框状态：`[ ]` 未勾；本腿⛔ 翻框、⛔ 改字。）

---
## 原文二：票 35（`.scratch/wisp/issues/35-panel-bridge-c17.md:93`，逐字；10-07 读数）

- **Go→页（出向）＝整条不存在**：ⓐ`grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd/wisp internal/panel`（剔 test）＝**0**，正控＝同尺全仓只命中 `.scratch/wisp/probes/33/p1/q2/main.go` 一枚探针；ⓑ泵 marshal 出的字节唯一去处＝`cmd/wisp/run.go:725 Out: rt.bookPanelSnapshot` → `panel_pump.go:321` 只写 `rt.lastSnap/lastSnapBytes` ＋ 一行日志，而读者 `(rt *agentRuntime).lastPanelSnapshot()` 带括号数＝**10 处调用全在 `_test.go`、非 test 0 处**；ⓒ`func (m *PanelManager)` 的导出方法 8 枚（`IsCreated/IsShown/LastColdMs/LastHotMs/Show/HotShow/Hide/Destroy`）没有一枚能把数据送进页面。

（补一行同票 `:391`，10-07 普查尺栏，逐字摘录关键段：出向原语非 test 调用者 `\.Eval(`＝**0**；`lastPanelSnapshot()` 非 test **只有定义**（`cmd/wisp/panel_pump.go:383`）⇒"写了零调用者"；并注：`panel.resync` 在 `SPEC-08:163` 标"Go→前端推送"、本不必入册，但一旦铸进 `bridge.go` 的 `Method*` 常量就落进 `inbound_roster_253_test.go` 的 DECLARED 尺⇒`wantFullInboundRosterSize253 = 6` 当场红。）

---
## 原文三：票 145（`.scratch/wisp/issues/145-...grow-the-go-side-carrier.md:208`，逐字）

- ⓑ记账值的读者：`(rt *agentRuntime).lastPanelSnapshot()` 带括号数调用点＝**10 枚，全部在 `_test.go`，非 test 0 枚**（`instructions_200r2_test.go` 3／`panel_pump_test.go` 3／`subagent_*_197_test.go` 4）。⇒ 泵 marshal 出的字节今天只被测试读过，没有任何生产读者。

---
## 射程对照（事实观察，⛔ 谁先谁后由编排者裁）

三枚断言的是**同一枚事实的不同深度**，覆盖关系如下（本腿只量）：

- **票 145** 只钉"数据侧读者"：`lastPanelSnapshot()` 生产 0 调用者（记账有、无人读）。
- **票 35 `:93/391`** 把它扩成"出向整跳"：字节唯一去处＝`bookPanelSnapshot` 只写字段＋日志，且出向原语（`PostWebMessage`/`EvaluateScript`/`.Eval`/COM）非 test＝0，`PanelManager` 8 枚导出方法无一送数据进页面。
- **票 33 `AC#14`** 把矛头指到"投递者＋库机制"这一层：不是缺 API，是缺**把绑定返回值/推送送回页面**的泵——并点出 `Dispatch` 队列仅在 `Run()` 内取出这一库级依赖。
- 三者交叠的"最后一公里"＝`pump.go` 错误文案自陈"（最后一公里在票 33/35 名下）"，即 33/35 两票对这一跳的归属是**并列**的，不是互相作废。

## 本腿现量对三枚原文的两处过期更正（锚 c315980b，尺见 logs）

1. **AC#14 的"本仓宿主对 `Run()`/`Dispatch()` 调用者＝0 枚"在锚点已不成立**：现量 `Run()` 生产调用者存在＝`cmd/wisp/panel_resident_windows.go:320 w.Run()`（＋ `:381 w.Dispatch(fn)`）；`Dispatch` 亦在 `msgcb`/host 内多处用（`panel_control_windows.go:736-751`、`panel_host_windows.go:689`）。⇒ AC#14 依赖的"队列永远不空"前提**已被常驻 runner 满足**；它那句"Go→页面投递者＝没有人"里，**RPC 回执那一跳其实已通**（`msgcb`→`Eval window._rpc.resolve` 在 `:752`），死的只剩"出向主动推快照"那一跳（本腿 `a3-survey.md` 必答题断点 `cmd/wisp/panel_pump.go:424`）。
2. **票 35/145 的行号与名册数是 10-07 读数、锚点已漂**：`run.go:725`→现 `:732`；`panel_pump.go:321/383`→现 `bookPanelSnapshot:360`、`lastPanelSnapshot:424`；`wantFullInboundRosterSize253 = 6`→现 `= 7`（`inbound_roster_253_test.go:63`，自 10-07 起多了 `config.models.list`）。⇒ 引用这三枚原文时须带"10-07 读数"前缀，别按当前行号直接引。

（本件仅并排原文＋量出覆盖/过期关系，⛔ 不裁"翻 33 格"还是"另立接线票"。）
