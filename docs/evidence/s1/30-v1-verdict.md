# 票 30 终裁腿 `30-v1` —— D10 结果呈现路由（终裁裁决表，非实现者）

> **尺**：实现 commit `eaeba3a92`（`internal/present` 8 枚 51,224 B＋落地件 `30-r1-impl.md`）。★全套门禁在**仓外导出树** `D:\tmp\wisp\20261028\30v1\export`（`git archive eaeba3a92` 现导）现跑，2026-10-31 12:0x–12:1x +08，`unset GOTMPDIR`、rc 全部重定向现取。★工作树零触碰：`git status --porcelain -- internal/present` 全程空；在飞腿脏面（41a2-r1/381-r2→`cmd/wisp`、59-r1→`internal/audio`）未碰、未归因。★导出树基线 `git hash-object`：`tier.go=75f527dd92f63e5b…`／`dispatch.go=9073f5cd6ed71d00…`，与 30-r1 记录逐字同值。★契约引文出处现核（第 20 处定式）：D10=`PLAN.md:242` 现跑仍是 D10 标题行；D42#8=`PLAN.md:3045` 起；票面前提「长结果必须附不超过 60 字的口播摘要」在 `internal/agent/prompt.go` style 段逐字在。
> **实现者六格自评不当验收**：本表每一格只认本腿自己复跑的读数。

## 1 ★★ present 全套复跑（必判①）

`go test ./internal/present/ -count=1 -v`＝**rc=0、19 top-level PASS／0 FAIL／0 SKIP**（另 23 子测试 PASS、`ok … 1.540s`）；`-race` 一发＝**rc=0**。与 30-r1 自报 19/0/0＋race 绿对上。

## 2 ★★★ 三发 mutation 复刻（必判②，仓外 mutant＋`-overlay`，树零写）

| 发 | 变异 | rc | 逐字红句（本腿现跑节选） |
|---|---|---|---|
| M1 | `tier.go` `Classify` 首行插 `return TierNone, "mutation M1: routing removed"` | **1** | `dispatch_test.go:154: tier = none, want structured`；矩阵四档全塌（60/61、400/401、800/801 案全红）；E2E 四子案全红 |
| M2 | `dispatch.go` tier-3 `if d.ch.Badge != nil {`／`if d.ch.TaskList != nil {` ⇒ `if false {` | **1** | `dispatch_test.go:225: badge updates = 0, want 1 (dual-channel is unconditional)`＋`:228 task list entries = 0`；trace＝`notify failed` 旁只有 `badge skipped`；E2E `badge=0 task=0` |
| M3 | `dispatch.go` `route()` switch 前插无条件 `WriteClipboard(display)` | **1** | `dispatch_test.go:180: clipboard called 1 times on tier 1, want 0`；`:234: clipboard called on tier 3, want 0`；tier-4 `clipboard = 2 calls` |

★**还原核（最后一发之后）**：无 overlay 重跑 **rc=0**；`git hash-object` 两件与基线逐字同值；导出树无 `mutation M` 标记；仓库 `internal/present` status 空。**三发全红、红句正中判据＝测试牙是真的。**

## 3 ★★ 正控四案复核（必判③；断言逐条自 dispatch_test.go 现读）

| 案 | 测试 | 读数 |
|---|---|---|
| 空白结果 `"   "` | TestEmptyResultRoutesNowhereAndPanicsNowhere | tier=none、六通道全 0、OK=true |
| 全端口 nil | TestUnwiredChannelsAreSkippedNotFatal | panel/speak/artifact 逐通道 `skipped/not wired`、不 panic、skip 不烂账 |
| 配额写失败 | TestArtifactFailureIsNamedAndPresentationSurvives | 注入 `memory: enforce artifacts quota…` ⇒ receipt 具名 failed、面板照常 1 次、播报诚实尾巴 `。全文未能保存`、OK=false |
| announce 拒收 | TestRefusedAnnounceIsReportedNotSwallowed | `Speak=false` ⇒ `speak failed: announce not taken…` 具名、OK=false |

## 4 ★★ AC2 的 mockllm 400 字节截断绕行合法性（必判④）

**判：合法，AC#2 可翻。** 依据（全部现核）：
- 事实：`tools/mockllm/chat.go:159-162` 逐字 `answer = "echo: "+lastUser.text; if len(answer) > 400 { answer = answer[:400] }` —— **字节**截断，402 汉字（1,206 B）物理回不来，且改 mockllm 在本腿写面之外。
- 注入面：AGENTS 1.3 许可缝里 **C5 LlmProvider** 是具名缝；E2E 的 `scriptedC5` 实现的正是 `llm.LlmProvider`，经 `agent.Options.Provider` 注入，四案**同一个 `agent.New` 回路、同一台 router**，402 案的 `calls==1` 在缝上计。这不是「用 mock 代替真的假报完成」：另三案走真 mockllm 子进程＋真 openai-chat 适配器读**服务器侧** `RouteCount("chat")==1`；被测的 `internal/present` 自身零 provider 引用（dispatch.go imports 无 llm/agent）——「永不二次请求」是结构性事实，缝上计数只是把票面的 single completion 在 402 案同样钉住。
- 诚实性：30-r1 在 §8③ 具名自报偏差而非藏；测试注释（e2e_present_30_test.go:21-31）同文具名。

## 5 ★★ D10 两处量词的具名解释核（必判⑤）

**判：成立。** D10 原文（`PLAN.md:261`）：「≥3 项的列表 / **≥2 个链接** / **文件路径清单**」——链接量词 D10 自己写了 ≥2；「清单」取 **≥2 行**是与之同构的最小复数读法，负向理由具名且测试钉死：`TestClassifySendsEveryTierWhereD10PutsIt` 里「散文单路径 `文件在 C:\x\a.txt 里` ⇒ tier-4」（否则一句话短结果被迫落文件、违背 D10 第 4 档本意）＋「两个路径行 ⇒ structured」。链接**先剥再数**（`linkCount` 先删 markdown 跨度再数裸 http(s)，href 不二次计数）由两案钉死：「单个 markdown 链接 ⇒ short」「md 链接＋裸 URL 计 2 ⇒ structured」。已在 30-r1 §8⑤＋测试名留痕。

## 6 ★★ 具名交回核（必判⑥）

- **装配缝**：`present.New(Channels{...})`（dispatch.go:66-88）七端口＋Logger；`git grep 'internal/present' HEAD -- cmd/wisp`＝**0 命中**（装配尚未发生，缝确实开着，归接线腿）。
- **Speak 不重复设闸**：`Speaker` 端口文档逐字（dispatch.go:17-19）「The ChTTS exfil gate lives on that announce path in the host (ticket 26 wiring) and is NOT duplicated here」——对照 `wire-r1-impl.md:10/18`：闸＝`risk.ChTTS`，位置＝`rt.voiceAnnounce` 播报路径（出声前、球 Speaking 之前）。present 侧 `speak()` 只递线＋把 false 记 `failed`，无第二道闸。**接线腿要点**：闭包必须原样走 `rt.voiceAnnounce`（既有 ChTTS 闸＋球 Speaking 联动都在那条路上）；D10 第 4 档「悬浮球」由 announce 路径点亮球（present 刻意不在 tier-4 发 Badge——E2E 读数 `notify/badge/task/panel 全 0` 自洽）。
- **Clipboard 端口义务**：端口文档逐字（dispatch.go:28-31）「The adapter MUST run risk.Provenance.CheckText(scope, risk.ChClipboard, text) before the real write」；真实 API 现核在（`internal/risk/provenance.go:826 CheckText`／`:123 ChClipboard`）。义务在端口文档＝交回事实成立（守不守是适配器/接线腿的事，30-r1 未把它报成已守）。

## 7 ★ 零回归复跑（必判⑦，同一导出树）

`go test ./internal/agent/... -count=1`＝**rc=0**（ok 3.160s）；`go test ./internal/memory -count=1 -run Artifact`＝**rc=0**；`go build ./...`＝**rc=0**。

## 8 ★ d22scan 导出树＋正控（必判⑧）

- 干净扫（`git archive` 导出树，无 .git ⇒ 扫描器 loud 模式全路径）：**rc=0 clean**；分母 **bans#1-5 internal/=274、cmd/=45、ban#8 internal/=706、cmd/=161 —— 与 30-r1 所报四数逐字对上**。（ban#8 design/ 文本件本树=14 vs 30-r1 报 39：其扫的是工作树、含机主外部 agent 未提交的 `design/**`；四枚 .go 分母不受影响，A207 本就注明分母机依。）
- 正控：导出树植真 U+2705 字符串 ⇒ **rc=1**，逐字 `internal/present/zz_d22scan_probe.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)`。

## 9 ★★★ 六格终裁（必判⑨）

| AC | 票面要求（缩写） | 终裁 | 依据（本腿复跑） |
|---|---|---|---|
| AC1 | 每档矩阵＋边界 60/61、400/401、800＋结构压长度 | **翻** | §1 全绿；tier_test 19 案含全部 D10 边界（60/61、400/401、800/801、空白填充、40 字表格、栅栏压 900 字符、双向 precedence）；M1 红 |
| AC2 | 单请求规则：每结果请求数==1、摘要随响应 | **翻（带具名绕行，§4 判合法）** | E2E 真 loop×4 案：mockllm 三案 `RouteCount==1`（服务器侧）、402 案 C5 缝 `calls==1`；tier1/tier2 案断言 `res.Text` 带同一响应 `口播摘要：…`；包内零 provider 引用 |
| AC3 | 工件落 artifacts、路径进日志＋载荷、配额被尊重 | **翻（带具名边界）** | 真 `memory.Store`：盘上文件内容＝display 逐字；路径进 receipt＋log（buffer 现读）＋播报尾巴；写只走 `PutArtifact`（配额在写路径同脸，`internal/memory/artifacts.go:71/82` 现核）；LRU 力学归票 04（30-r1 具名，本腿认） |
| AC4 | 剪贴板只 tier-4 | **翻** | 结构性单分支（`case TierShort:` 内唯一 `WriteClipboard`）＋M3 红＋tier1/2/3 案 calls==0 |
| AC5 | tier-3 双通道（toast 死角标照更） | **翻** | badge/task 调用在 notify 错误分支之外（dispatch.go:171-217 结构）；M2 红逐字；toast-blocked 夹具 badge=1/task=1 |
| AC6 | 四结果走真管线逐档观察 | **翻（带具名边界）** | 真 `agent.Loop`（真 C5 缝、真 turn 组装、真 EvDone）×4 档逐通道现读；真球/真出声读数随 `cmd/wisp` 接线腿（Speaker 端口只证明递了哪句话；与 26-r1 AC1 判定同形，30-r1 §8① 已具名） |

**结案判定：六格全翻，可带残余结案。** 残余三条逐条：
1. mockllm 回声 400 字节**字节级**截断（`tools/mockllm/chat.go:160`）⇒ 归 tools/mockllm 所有者立票（多字节文本拦腰切；golden 模式不受影响）。
2. D10 两条量词解释（路径清单 ≥2 行／链接先剥再数）⇒ 已在 30-r1 §8⑤＋测试名＋本表 §5 三处留痕，无需改契约；后续接线/面板腿引用时以此为准。
3. `cmd/wisp` 装配缝（`present.New(Channels{...})`＋`loop.Run` 后对 completed 结果 `Dispatch`；`Speak→rt.voiceAnnounce` 闭包不重复设闸；Clipboard 适配器先 `CheckText(scope,ChClipboard,text)`）⇒ 归接线腿（41a2-r1/381-r2 之后的波次）。

★本腿未定义即停：未触发（派单九项全部可裁、无歧义项）。★写面：本文件一件；临时件全在 `D:\tmp\wisp\20261028\30v1\`（logs/mut/overlay/export）。
