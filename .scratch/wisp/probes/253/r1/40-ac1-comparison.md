# 253-r1 · 我的做法 vs 票面 AC#1 逐字判据（40-ac1-comparison）

时刻 `2026-10-08 16:53 +0800`。回报要求这一句单独说清，所以单开一文件，逐字对逐字。

## 票面原文（两处，逐字）

`AC#1` 本体＝`.scratch/wisp/issues/253-panel-inbound-three-ruler-holes.md:16`：

> **AC#1 可达性那枚钉改成问能力**：判据不许再认某个 JS 调用的**词面**，要问**"页面发起的一次调用，Go 侧有没有收到并回执"**（载体沿用票 33 那族"由页面自己报回"的定式：`REPLIED` 那形）。⚠ 正控必配：造一发"只 `postMessage`、不 `wispDispatch`"的假腿 ⇒ 指名用例**必须红**。⛔ 不许把两种拼法都列进白名单了事（那是把尺调宽，不是把门修好）。

同票 §8 编排者改判窄义＝`:42`（**票面自己已经推翻 `:16` 的宽义口径**，逐字）：

> **AC#1 的增量改窄（推翻本票排程节里我给的宽义口径）**：`cmd/wisp` **今天确实已有"页面发起、Go 侧回执"的行为用例**——`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`cmd/wisp/panel_resident_windows_test.go:314`）与 `TestAC14AwaitedBindingReplyReachesThePage`（`:812`，两半互证：`countRequests==3` 且页面 `REPLIED`）。⇒ 真洞只在窄义：**没有任何词面/AST 尺断言产码里 `w.Bind("wispDispatch")` 存在**。

再往下一道＝台账 `A717` §5（10-08 14:1x，派单体里引的就是它）：能力形那半今天已有真产码行为尺（`cmd/wisp/panel_transport_35r1_test.go:195→:203` 直调 `installPanelTransport`）⇒ **代笔只准补"绑定字面量"那一窄形**。

## 出入（四条，一和二是对立面，我照票面 §8＋`A717` 做，不照 `:16` 字面做）

1. **具名出入①**：`:16` 第一句字面是"判据**不许再认词面**"，而**我这枚交付的正是一枚词面/AST 尺**（`os.ReadFile`＋`go/parser`，读产码字节，一枚产码函数都不调用）。
   ⇒ 我照的不是 `:16` 的字面，是同票 `:42` 那句自我推翻＋`A717` §5 的收窄：`:42` 把"真洞"重写成"没有任何词面/AST 尺断言产码里 `w.Bind(...)` 存在"，而能力形那一半 `:42`/`A717` 都判给**已存在**的 35r1/35r2/35r3。照 `:16` 字面做＝这一格什么都不建（因为"问能力"那台仪器今天已在盘上），所以我把"出入"读成**票面自身已改判**，不是我违抗票面。⛔ 我不翻框、也不自称结案。
2. **具名出入②**：`:16` 要的载体是 `REPLIED` 那形（页面自己报回），我这枚**没有**任何 `REPLIED`、不开窗、不发封套。那一形今天在哪＝`cmd/wisp/panel_resident_windows_test.go:812 TestAC14AwaitedBindingReplyReachesThePage`（`:42` 逐字点名）与 `panel_transport_35r1_test.go:195`。⛔ 我不重造（重造＝`issues/README.md:80` 机主原话"不要重复劳动"那一格）。
3. **具名出入③（正控那一发怎么配的）**：`:16` 的正控是"造一发只 `postMessage`、不 `wispDispatch` 的假腿 ⇒ 指名用例必须红"，那是**激励层**的红，词面尺没法"发"东西。我给的是同一坏法的**词面对应物**，两发都盘上种过、都 rc=1：
   - **M3**＝把 `:792` 的 `fmt.Sprintf` 实参换成字面量 `"wispOther"`（＝转发钩子还在装、但它叫的门**不再由那枚被绑定的常量生成**），红句逐字 `panel_dispatch_binding_roster_253r1_windows_test.go:398: AC#1 RED: installPanelTransport calls w.Init([panel_host_windows.go:808]) with no script built (via fmt.Sprintf) from "panelDispatchBinding" ...`；
   - **M1**＝`:802` 绑一枚名册外的名（`"wispStaleDoor"`），两面红（`:377`＋`:383`）。
   ⇒ "真发一封只 `postMessage` 的假腿"那一格仍归 35r1（它红句 `:209` 逐字 `AC#6 RED: installPanelTransport registered no page->wispDispatch forwarding hook`），⛔ 我没把它算成自己的。
4. **具名出入④（白名单那句我守住了）**：`:16` 末句"⛔ 不许把两种拼法都列进白名单了事"。我的 `doorRoster253r1` **只一枚** `"wispDispatch"`；`cmd/wisp/panel_host_windows.go:852` 那枚 `wispProbeRT`（`firstRoundTrip` 的一次性探针门）**没有**被列进名册当白名单——它落在名册射程**之外**（尺只走 `installPanelTransport` 的函数体），我把这一格作为**恒真面 #3 具名公开**（附实证：干净树 rc=0 而 `:852` 同存），⛔ 不拿"它也红／它不算"当两不得。
5. 另一条必须说白的**看不见**：`:16` 问的是"Go 侧有没有收到**并回执**"。我这枚对"**回执**"那一半**零断言**（不看闭包体、不看 reply 走没走回页面）——恒真面 #6/#1 具名。

## 一句话结论

我这枚尺＝票面 `:42`＋`A717` §5 收窄后的那一枚**词面/AST 绑定锁步尺**，⛔ 不等于 `:16` 原句的"改成问能力"；`:16` 那半的凭据仍在能力形那三枚现存尺上，而**这两半合起来才叫 `AC#1`**。所以我按派单**不翻框**，把"这格还缺哪一发读数"留给非实现者验收腿（`253-v1`／`259-v1` 那族）裁。
