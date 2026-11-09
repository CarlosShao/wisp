# 票 411 AC#4 事实收集腿 —— 9 腿真红逐枚「事实＋三档建议＋代码行引用」（2026-11-09）

- 身份：本腿**只查事实、不下最终裁定**；三档＝甲真缺陷／乙合法新码待重登基线／丙 pattern 不精确。
  裁定归 154/158/161/171/178/195 家族裁决者（票 411 AC#4）。
- 锚点：`cb2287a5`（＝修后真读数的锚；`git diff --stat cb2287a5 HEAD -- internal cmd` 零改动，
  锚树＝现 HEAD 树，产码行号引用在锚树现读，非转述）。
- 方法：按仪器 `.scratch/wisp/probes/154/gate-clauses.sh` 的词表与 pair 算术**逐字复刻**，在锚树上
  现跑 `git grep` 复算每腿名册；9 腿复算值与 flip 电池基线发存档
  （`probes/154/r20261108/30-flip-final.txt`，锚 cb2287a5）**9/9 对平**：
  G1=1／G1b=4／G2=2／G5=2／G5pos=1／G5neg=14／G6neg=1／G7=6／G7neg=7。
- 门禁（本腿零产码改动 ⇒ gofumpt/vet/build 无对象）：`d22scan` 现跑 **clean rc=0**
  （bans#1-5 internal=387＋cmd=64；ban#8 internal=966＋cmd=244 全数过）。
- 工作树状态核过：internal/＋cmd/ 无别人在飞的 tracked 改动（仅 2 枚未跟踪 test 件，不在锚树内）。

## 0. 对平声明

复刻器一处自纠在案（防下一位再踩）：GNU grep 3.0 的 `-Ew` 对 `Defer\(Named\)?`（多转义的 `\(`）
**不匹配**（现量 rc=1），对仪器原文 `Defer(Named)?`（无转义）匹配正常。仪器原文没错；首次复算
G7 得 10≠7 就是这一味引起的，换回仪器逐字词表后 6/7 全对平。

## 1. G1（run 腿，quiet/want_n 0，实测 1）——建议：丙（pattern 不精确，加 `-w`）

- 唯一代码命中（锚树原行）：`internal/plugin/manifestslot/tier2.go:46`
  `out, err := ex.Execute(ctx, plugin.Tier2ToolRequest{`。
- 类型事实（现读定义）：
  - `ToolRequest`＝`internal/agent/tools.go:54`（字段 TaskID/CorrelationID/CallID/Name/Args/Timeout，
    注释「ToolRequest is one dispatch」——**agent 派发请求**，G1 想抓的形状）。
  - `Tier2ToolRequest`＝`internal/plugin/tier2_loader.go:80`（字段 PluginID/Tool/Declared/Params，
    注释「one manifest tool call, as the loader receives it」——**C24 插件层的另一个类型**，
    不是 agent.ToolRequest 的别名/嵌入/子类型）。
- 机理：G1 词表 `-nE 'ToolRequest\{'` **无 -w**，`Tier2ToolRequest{` 里的子串 `ToolRequest{` 被命中。
- 修法实测（正反两向）：`grep -Ew 'ToolRequest\{'` 对 tier2.go:46 那行 **rc=1（排除）**；对真
  `agent.ToolRequest{` 形状行 **rc=0（保留）** ⇒ 加 `-w` 即修、不丢真命中。
- 名册里另一行 `internal/tools/bridge.go:984` 是整行注释，已被 run 腿分子剔掉（hits=1 正是剔注释后的数）。
- 为什么是丙不是乙：这一枚响**不是**基线过期，是词表子串误命中了 G1 语义之外的类型；重登基线
  只会把误命中固化进判据。若裁定者认为 Tier-2 构造本身也该被某条腿盯住，那是新腿/改词表语义的
  裁决后动作，不是 G1 现词表的本意。

## 2. G1b（run 腿，quiet/want_n 0，实测 4）——建议：乙（合法新执行面，重登 0→4）

四枚命中逐枚（锚树现读）：

| # | 位置 | 原行 | 是什么 | 落地 |
|---|---|---|---|---|
| 1 | `internal/plugin/manifestslot/tier2.go:46` | `ex.Execute(ctx, plugin.Tier2ToolRequest{` | Tier2Run 适配器→`Tier2Executor.Execute`（C24 面，见 tier2.go:34-38 注释） | 11-05 `9d27a8a7f`（票 51-B） |
| 2 | `internal/plugin/tier2_loader.go:195` | `out, err := e.gate.Execute(ctx, GojaCall{` | `Tier2Executor`→`Tier2Gate.Execute`（gate-first 门） | 11-05 `9d27a8a7f` |
| 3 | `internal/plugin/tier2_gate.go:210` | `return rt.Execute(ctx, call)` | `Tier2Gate.Execute` 方法体（:204 定义）内转 `GojaRuntime.Execute` | 11-05 `eafb13d5f`（票 51-ac5） |
| 4 | `internal/tools/subagent_197.go:616` | `return p.inner.Execute(ctx, req)` | `subagentToolProvider`（`agent.ToolProvider` **装饰器**，:582 `inner agent.ToolProvider`）在拒掉嵌套 spawn 后**委托 inner**＝seam 本层，不越桥 | **09-28** `b9fa815b4` |

- 时间线事实：G1b want_n 0 是 171 r2（09-27，f1b99a70）现量登记的；**10-14 锚 `605838b07` 的存档门
  读数（`probes/178/r1/logs/gate-post.txt`）里 G1b 已经 BAD 实测=1**（＝第 4 枚）⇒ 这枚基线过期
  **早于 411 修**，不是修尺修出来的；10-14 后新增的是 Tier-2 三枚（11-03~11-05）。
- 定性：第 4 枚＝装饰器委托（与 loop.go:875/:879 `l.opt.Tools.Execute` 同层，不构成旁路）；
  第 1-3 枚＝Tier-2/C24 插件执行栈的自有调用链（Tier2Run→Executor→Gate→GojaRuntime），链上有
  自己的门（tier2_gate.go:28-36 注释「Execute is the call door: gate first」；默认关＋无 L2 确认
  不得开＋逐插件开关，:175-198 现读在案）。四枚都不是「绕开 bridge/loop 的裸执行」。
- 建议乙、重登 0→4。附一枚留给裁定者的口径注：tier2_gate.go:210 位于 Tier2Gate.Execute 的
  **定义本体里**（与 loop.go/bridge.go 被排除的理由同形）；若裁定者按「定义链免点」的同一口径
  处理，重登值是 3 而非 4——两读都写在这里，裁定者现量拍。

## 3. G2（run 腿，ring/want_n 1，实测 2）——建议：乙（配对完整的合法新码，重登 1→2）

- 锚树现读：`cmd/wisp/run.go:1514` `rt.bridge.OpenTask(taskID)` 与 `:1520`
  `rt.bridge.CloseTask(taskID)` **同在 `admitTask`（:1508）里**：开在任务起点，关在返回的
  清理闭包里（:1521-1527）；loop 端收口＝`internal/agent/loop.go:409-413`
  `if r := l.opt.AdmitTask(taskID); r != nil { revokeAdmission = r }` ＋ `defer revokeAdmission()`
  ⇒ **配对完整、端到端有收尾**（「WHY THIS HOOK OWNS THE CLOSE」注释 :1475 在案）。
- 落地＝票 315（`bfc74ae9d`，10-14，commit 标题即「admitTask 任务开始补 OpenTask，与同钩子
  revoke 里那枚 CloseTask 同 id/同所有者成对」）。
- 时间线：10-14 锚存档 gate-post.txt 里 G2 已 BAD 实测=2 ⇒ 基线 1 从 10-14 起过期，早于 411 修。
  名册另 3 行（run.go:1473/:1475、subagent_197.go:360）全是注释、分子已剔。

## 4. G5（pair 腿 ring/want_n 1，实测 2）与 G5pos（quiet/0，实测 1）——建议：乙（同在册活形状的新实例，重登 1→2）

- 新点名的文件＝`cmd/wisp/resident_present_windows.go`（10-31 落 `c15888b28`，票 30 装配缝腿）。
- 开：`:80` `if prov.OpenScope(presentClipboardScope) == nil { return nil }`——构造期一次性开，
  开不了就 nil（fail-closed，:66-70 注释「returns nil only when the gate itself cannot exist」）。
- **Close 的真实落点＝同文件无、全树无（现跑全 cmd/internal 的 OpenScope 调用点核对）。这不是漏，
  是设计**：注释块 :40-47 逐句写明——这是**进程生命周期 scope**（与
  `resident_voice_engines_windows.go:89-99` 的 voiceAnnounceScope 同形），**不是任务 scope 的影子**；
  理由＝admitTask 的 revoke 在 loop.Run 返回前就关任务 scope，而 CheckText 对未绑定 id fail-closed，
  闸若指着已关的任务 scope 会对每次剪贴板写入误拒。
- 同形状的在册样板（G5 基线 1 当初就是为它登记的）：`cmd/wisp/panel_assets.go:243`
  `_ = prov.OpenScope(taintSourceScopeID)`，其注释 :240-242：「进程生命周期探针没有可关它的边界，
  `_ =` 把这句说出来」。家族第三枚同形＝`resident_voice_windows.go:601/:607` attachProvenance 开
  voice-announce scope（句柄 nil 检查后不留）——它被 G5 尺读成 balanced 只因文件里别的
  泛用 `.Close()`（:776/:804/:817/:824，引擎收尾）；这正是仪器 handle 味明写的
  CLOSE-BY-GENERIC-CLOSE 粗糙处（泛用 Close 全仓 42 枚文件，171 r2 现量）。
- 结论建议：与在册活形状逐位同形（构造期开、句柄不留、文件内无收尾、注释具名）⇒ 乙，
  重登 1→2。是否翻「进程生命周期 scope 不关」这个产品语义本身＝票 158 登记时已接受过的判断，
  翻案归裁定者（翻案则 panel_assets.go 也要一起翻，不是只治新的一枚）。
- G5pos 的 1 枚＝同一文件（正控射程只排 panel_assets.go 一枚，机制正确），随主尺同批重登；
  正控的排除名单要不要加 resident_present_windows.go＝重登时的同格动作（171 r3「谁改口径同格核下来」）。

## 5. G5neg（pair 腿 ring/want_n 8，实测 14）——负控按设计响；漂移＝纯新增文件，建议：乙（重登 8→14）

- 10-14 锚（`605838b07`，存档 gate-post.txt 里该腿恰好 ok＝实测 8）复算名册 8 枚现跑在案：
  `cmd/wisp/panel_assets.go`、`internal/risk/provenance.go`、`provenance_syncdirs_other_test.go`、
  `provenance_test.go`、`syncdirs_redteam_windows_test.go`、`syncdirs_test.go`、
  `syncdirs_windows_test.go`、`taintmatch_test.go`。
- 今天 14 枚＝老 8 枚**全数在场**＋新 6 枚：`resident_present_windows.go`（10-31）、
  `resident_voice_markseam_394_windows_test.go`（11-01）、`internal/risk/markseam_394_test.go`（11-01）、
  `resident_voice_ttsinherit_382wr2_windows_test.go`（11-01）、
  `internal/plugin/disposal_risk_403_test.go`（11-04）、
  **`internal/risk/provenance_160_boundary_test.go`（10-22，票 160 落地腿 `c4d65e4ed`）**。
- ★勘误（对派单）：派单把 +6 里的一枚点名为 `provenance_syncdirs_other_test`——该文件 09-21
  落地（`ea01b9ff5`）、**10-14 时已在 8 枚名册里**，不是漂移来源；真正的新第 6 枚是
  `provenance_160_boundary_test.go`（`git log --diff-filter=A` 与锚树 diffset 现跑双证）。
- 盲区核对：盲区 9 枚文件（approval×4＋manifestslot×4）**一枚都没进 G5neg 名册** ⇒ 这条腿的
  漂移与射程变宽无关，纯粹是 10-14 后新增文件（含正控污染源 resident_present）。
- flip 电池 stale 发 3 FAIL（G5neg 基线抬到 9 仍 14>9 走 BAD）＝同一枚 finding 本体
  （`30-flip-final.txt:99-105`），重登 8→14 后该发自洽。

## 6. G6neg（piles 腿 quiet/want_n 0，实测 1）——②只开堆＝一枚测试文件；建议：乙（重登 0→1），或测试卫生补关归 0——两案摆裁定者

- ②堆唯一文件＝`internal/tools/web_bridge_test.go`（11-02 落 `547d9da00`，票 22）：
  `:186` `b.OpenTask("task-22-clean")`，同文件无 CloseTask。
- 现读事实：这是**故意的只开**——注释 :181-185 逐句写明「OpenTask is the composition root's duty
  at task admission; calling it here keeps the test's view equal to the production shape」，且未开
  scope 的 fail-close 行为已由 `Test22UnopenedScopeFailsClosed` 单独断言（:183-185 引用）。
- 泄漏面核实：桥是**每测新建**的（`webFixtureBridge` :73-108，`b := New(Options{...})`）⇒ 不关的
  scope 死在自家 fixture 里，无跨测状态泄漏。
- 定性：票 178 的②堆判语针对的是**生产码**「只开」（防作弊正控注入的也是生产形状）；测试里
  自备的开侧不是那个形状。二案：⑴ 乙——重登 0→1（piles 基线从此含测试自备形状）；
  ⑵ 测试卫生——测试尾补 `b.CloseTask("task-22-clean")` 让②堆回 0（一行、无行为变化）。
  两案都成立，归裁定者拍；判不出的问题：无（形状全静态可读）。

## 7. G7（pair 腿 ring rev/want_n 1，实测 6）与 G7neg（ring rev/want_n 2，实测 7）——红全在反向；5 枚新增逐枚全找到收尾落点 ⇒ 建议乙

先对平：两腿正向（Defer 无 DisposalScope）均 **0 枚**；G7 反向 6 枚＝基线在册的
`internal/memory/retention.go`（161-v2 判「形状使然＋词表缺一味」，正控地位）＋新 5 枚；
G7neg 反向 7 枚＝retention.go＋`retention_test.go` 两枚老名册＋同 5 枚新枚
（disposal.go/disposal_test/403_test 因有 Defer 定义与调用而 balanced，现算在案）。

新 5 枚逐枚（锚树现读，全为 10-26~11-05 落码）：

| # | 文件 | DisposalScope 在哪 | **Close/Dispose 真实落点** | 判 |
|---|---|---|---|---|
| 1 | `cmd/wisp/resident_voice_bargein_windows.go`（10-31 `4c21ac42c`） | 构造**注入**的字段（:239 字段、:310 参数 `scope *plugin.DisposalScope`） | 用途＝:458 `m.handle = m.scope.Go(voiceBargeinRoster, ...)` 把 barge-in 常驻 goroutine **挂上 owner 的 scope**；收尾在 owner＝`resident_voice_windows.go:690` 建 "voice-15e" scope、`:751` `stop()` 里 `rv.scope.Dispose()`（D38(e) 步骤 4 关停钩子，:735-739 注释） | **委托，非真漏** |
| 2 | `cmd/wisp/resident_voice_windows.go`（10-26 `ca8834d72`） | :690 `rv.scope = plugin.NewDisposalScope("voice-15e", nil, opts...)` | 同文件 :751 `res := rv.scope.Dispose()`（stop()，stopped 幂等防重入 :742-747）；引擎收尾另有显式 Close（:776/:804/:817/:824） | 同文件有收尾（Dispose），非真漏 |
| 3 | `internal/plugin/goja_runtime.go`（11-03 `8a2a40b0f`） | :257 per-call scope（`call.Scope==nil` 时自建 "tool:goja:…"） | **同函数 defer**：:258-263 `defer func(){ if res := scope.Dispose(); res.Incomplete {…Warn…} }()` | 同函数收尾，非真漏 |
| 4 | `internal/plugin/manifestslot/command.go`（11-03 `5438d104e`） | :225 每调用 scope（"plugin-command:…/…"，父＝注入的 ScopeParent，:65-66/:101-104） | :260 `_ = scope.Dispose()` **顺位调用（非 defer）**，在 proc.Wait 之后；panic 窗口由 lifecycle 的 ScopeParent 父联级兜底（uninstall 时 :456 Dispose 级联 join） | 收尾在场（顺位＋父级两道）；「顺位而非 defer」一枚小口径注留给裁定者 |
| 5 | `internal/plugin/manifestslot/lifecycle.go`（11-03 `d99b1a970`） | :84 每插件 scope 字段；:247-252 `newScope` 建（父＝opt.ScopeParent，:66） | :456 `rep.Disposal = mp.scope.Dispose()`（uninstall 生命周期事件；Incomplete 显式报错 :458-459）；:172-192 注释明写调用 ctx 绑定到该 scope、Dispose join 它 | Manager 持有、生命周期事件收尾，非真漏 |

- 共同语义（为什么「无 Defer」≠「无收尾」）：G7 反向词表抓的是「拿着 scope 却没挂收尾」；
  C11 里 **Defer 是可选的登记面，Dispose 才是收尾的执行点**——不登记任何 Defer、靠 Dispose
  做 ctx 取消/联级 join 的 scope 是合法形状（disposal.go 定义本体现读为证：Defer 返回 bool 供
  登记失败告警，Dispose 跑的是已登记的收尾集）。五枚的收尾全部走 Dispose（同函数 defer×1、
  同函数顺位×1、owner 文件生命周期方法×1、构造注入委托×1、Manager 生命周期事件×1）。
- 建议：乙——五枚全部找到落点，无真漏 ⇒ G7 重登 1→6、G7neg 重登 2→7。
  留给裁定者的另一扇门（裁决后才动的改判据动作，本件不代拍）：若裁定者认为 G7 反向的合方
  也该认 `Dispose()`（而不只认 Defer(Named)?），那五枚会整批退场、基线回 1/2——这是改词表
  语义，按 171 r3 规矩改口径同格核基线，另票另程。

## 8. 派单具名的两处勘误（事实层）

1. **G5neg 的 +6 名单**：新第 6 枚是 `internal/risk/provenance_160_boundary_test.go`（10-22），
   不是派单写的 `provenance_syncdirs_other_test`（09-21 落地、10-14 已在册）。见 §5。
2. **「10-14 后新落」的时间线**：G1b 的 subagent_197.go:616（09-28）与 G2 的 run.go 对（10-14）
   在 10-14 锚存档读数里**就已在响**（gate-post.txt：G1b=1 BAD、G2=2 BAD）——这两枚基线过期与
   411 修尺无关，修尺只是让全树读数第一次完整可见；纯 10-14 后新增的是 Tier-2 三枚（11-03~05）、
   resident_present（10-31）、G7 五枚（10-26~11-05）、web_bridge_test（11-02）与 G5neg 测试四枚。
   盲区身份另注：tier2.go／command.go／lifecycle.go 是「盲区文件＋新文件」双重身份；
   goja_runtime.go 在 internal/plugin/（两层深，旧射程本可见），是 11-03 新落码不是盲区复活。

## 9. 汇总表（9 腿 × 三档建议）

| 腿 | 基线→实测 | 建议 | 重登值 | 一句话依据 |
|---|---|---|---|---|
| G1 | 0→1 | **丙** | （加 -w 后回 0） | `Tier2ToolRequest` 是 plugin 层另一个类型，子串误命中；`-Ew` 实测正反两向都对 |
| G1b | 0→4 | 乙 | 4（或 3，定义本体口径由裁定者拍） | 1 枚装饰器委托＋3 枚 Tier-2 自有门控链，无裸执行 |
| G2 | 1→2 | 乙 | 2 | admitTask 开/关同文件、loop 端 defer 收口；10-14 起基线已过期 |
| G5 | 1→2 | 乙 | 2 | resident_present 与在册活形状 panel_assets 逐位同形，注释具名「进程生命周期」 |
| G5pos | 0→1 | 乙 | （随 G5；正控排除名单同格补） | 同一枚文件污染正控，机制正确 |
| G5neg | 8→14 | 乙 | 14 | 老 8 枚全在＋新 6 枚（含勘误后的 160_boundary_test）；盲区零贡献 |
| G6neg | 0→1 | 乙（或测试卫生补关归 0） | 1 或 0 | 测试故意只开＋每测新桥无跨测泄漏；两案都成立，裁定者拍 |
| G7 | 1→6 | 乙 | 6 | 5 枚新增全部找到 Dispose 落点（同函数×2、owner×1、委托×1、Manager×1） |
| G7neg | 2→7 | 乙 | 7 | 同 G7 五枚；disposal 定义族 balanced，对平在案 |

- **甲＝0 枚**：9 腿里没有一枚查成「生产码真漏」。唯一两处值得裁定者另眼看的口径注
  （不是缺陷断言）：G1b 的 tier2_gate.go:210 定义本体免不免点（§2）；command.go 的 Dispose
  顺位而非 defer（§7 表 #4，父联级兜底在案）。
- 需运行时才能定的问题：**无**（9 腿全部是静态可读形状；G7 五枚的 Dispose 都挂在具名的
  生命周期路径上，路径存在性由代码与注释现读，非运行时断言）。
