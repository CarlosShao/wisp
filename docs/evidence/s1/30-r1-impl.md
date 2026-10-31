# 票 30 落地腿 `30-r1` —— D10 结果呈现路由（落地件，写腿）

> **尺**：`HEAD=9d99af5c4`（`dev`，现跑 2026-10-31 11:0x–11:5x +08）。★开工三核：分支＝`dev`；★**HEAD 与派单锚点 `281862907` 不等，但该锚是当前 HEAD 的祖先**（其后是编排者停车点两笔＋票 42 落地一笔，与派单"42-r1 已收口"自洽）＝**锚点过期而非旧提交/detached**，本腿未切任何东西，按 dev 顶开工。索引脏面＝`cmd/wisp` 两件（`41a2-r1` 在飞）＋台账一件（编排者）＋零散 untracked，**本写面（`internal/present`）开工前为空**。
> **写面**：`internal/present/**` **新增 8 枚／51,224 B**（`doc.go` 1,709／`tier.go` 6,676／`summary.go` 3,451／`dispatch.go` 10,302／`tier_test.go` 3,731／`summary_test.go` 4,158／`dispatch_test.go` 11,382／`e2e_present_30_test.go` 9,815）＋ 本文件。**⛔ 零字节**：`cmd/wisp/**`（`41a2-r1` 独占，本腿未触）、`internal/risk/**`、`internal/speech/**`、`internal/session/**`、`internal/panel/**`（C17 桥未动）、**`internal/agent/**`（`prompt.go` 等逐枚 `git hash-object` ＝ HEAD 同值；摘要协议因此做在分发器侧、prompt ⑥ 未改一字）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、specs、PLAN。零配置键、零 thresholds、零 golden。
> **临时件（只建不删）**：`D:\tmp\wisp\20261028\30r1\{logs,mut,overlay,bin,export}`。★三发变异全走**仓外 mutant ＋ `-overlay`** ⇒ **工作树零突变窗口**（§5）；d22scan 正控走**仓外 `git clone --shared` 导出树**（§7）。

## 0 一句话

`internal/present` 新包：D10 四档分发器（从上到下第一条命中、结构化标记优先于长度）＋ 同响应 ≤60 字口播摘要协议（末行 `口播摘要：…` 提取，无则首句兜底，**永不二次请求**）＋ 工件经 `memory.Store.PutArtifact`（名字守卫＋私有 DACL＋写路径配额同一张脸）＋ tier-3 双通道（notify/badge/task 无条件并排，toast 死了角标照更，D42#8）＋ tier-4 才许碰剪贴板（结构性单分支）。通道全是窄端口，`nil` 端口＝跳过不炸；真 TTS/球/面板/通知的接线具名交回 `cmd/wisp`（§8）。

## 1 ★★★ 票面 6 格 AC 逐条（原文准引 ＋ 判定 ＋ 读数）

| AC | 票面原文（逐字摘） | 判定 | 读数 |
|---|---|---|---|
| AC1 | `Tier matrix tests: crafted outputs for every tier + boundary values (60/61, 400/401, 800) + structured-vs-length precedence cases.` | **✅ 翻** | `TestClassifySendsEveryTierWhereD10PutsIt`：**18 案**含 60/61 汉字、400/401 汉字、800/801 字符（拉丁、低汉字）四组边界＋空白填充案（**700 字母＋200 空格＝short**：空白不计数）＋「40 字表格＝structured」＋代码栅栏压长度（900 字符＋栅栏）＋列表/双链接/双路径案＋**单链接不构成链接清单**（链接 href 不被二次计数）案＋**散文里单个路径留在 tier-4**（D10 的"文件路径清单"取 ≥2 行；≥1 行会让"文件在 C:\x\a.txt 里"被迫落文件——具名解释，§8⑤）。整案＋`TestClassifyStructuredBeatsLengthBothWays` 双向全绿 |
| AC2 | `Single-request rule: mockllm request count == 1 per result (summary present in output).` | **✅ 翻** | E2E（§3）：mockllm 三案各 `RouteCount("chat")==1`（**服务器侧计数**）、C5 缝案 `calls==1`；**且逐案断言 `res.Text` 里带着同一响应的 `口播摘要：…`**。本包零第二次 LLM 调用（分发器无任何 provider 引用，结构性成立） |
| AC3 | `Artifacts: file written under artifacts dir, path in logs + result payload; quota respected.` | **✅ 翻（带边界）** | **真 `memory.Store`**（`t.TempDir`）：文件真落 `ArtifactsDir()` 盘上、内容＝display 文本逐字；路径进 **receipt**（`rc.Artifact`）＋**日志**（`present: artifact written … path=…`，buffer handler 现读）＋**播报行**（`。全文已存：<全路径>`）。★配额：写只走 `PutArtifact`（**写路径配额检查就在这张脸里**，artifacts.go:81）；500MB LRU 机的力学属票 04 自己的测试（`internal/memory` artifacts 族），本腿注入**配额写失败**证明拒绝被具名上报、呈现不断（§4）。★边界：LRU 逐出读数不在本腿（那是 04 的面） |
| AC4 | `Clipboard: tier 4 writes; tiers 1–3 never touch clipboard.` | **✅ 翻** | tier-4 案：剪贴板**恰一次**、内容恰＝display；tier-1/2/3 案（单测＋E2E）：`clipboard calls == 0`。**结构性**：`WriteClipboard` 只存在于 `case TierShort:` 分支（dispatch.go route()）——**M3 变异**把它提到全档共享段 ⇒ 三条红（§5），牙是真的 |
| AC5 | `Tier 3 dual-channel: toast sent AND badge/task-list updated (event trace), toast-blocked fixture → badge still updates.` | **✅ 翻** | `TestTierThreeDualChannelSurvivesAToastFailure`：`Notify` 注入报错后 **badge 恰 1 次、task 恰 1 次**，receipt trace 逐字 `[{notify failed …} {badge sent} {task sent} {panel sent} {speak sent}]`、`OK=false`；E2E tier-3 案同读数。**M2 变异**摘掉两调用 ⇒ `badge updates = 0, want 1 (dual-channel is unconditional)` 红（§5）。★D42#8 的"无条件"是结构：badge/task 的调用**不在** notify 的错误分支里，物理上无"toast 成功才更角标"这条路 |
| AC6 | `E2E: four results through the real pipeline; ball/panel/notify/audio observed per tier.` | **✅ 翻（带具名边界）** | §3：真 `agent.Loop`（真 C5 缝、真 turn 组装、真 `EvDone`）×4 档；面板/通知/角标/任务表/剪贴板/播报行/工件盘上文件逐案现读。★**边界**：ball 的 Speaking 态与真实出声在 `cmd/wisp` 语音腿（`26-r3` 已接 `announce()` 门），本腿观察到的是 **`Speaker` 端口被递了哪句话**——真球/真音频读数归接线腿（§8①），与 `26-r1` AC1 的判定同形 |

★票面自身：`Status: ready-for-agent`、`Claimed by: —`、Progress log 0 行——**⛔ 本腿不改票面**（翻框/Status 是编排者收尾动作）。

## 2 ★★ 契约引文出处现核（第 20 处定式）

- `PLAN.md:242` 起 **D10 结果呈现层**：四档表逐字对上（含"结构化标记优先于长度""摘要由 LLM 在同一次响应里给出""写剪贴板只在第 4 档做""⚠ 第 3 档……必须同时更新悬浮球角标"）。`SPEC-05 §7` 同表同义。`D42#8`（`PLAN.md:3047` 附近）：Focus Assist ⇒ "总是同时更新悬浮球角标 + 面板任务列表，不依赖通知"。**票面引文与 PLAN 逐字一致，第 19 处（前提腐烂）不触发**：票面前提"prompt ⑥ 段已有 ≤60 字口播摘要指令"现核成立（`prompt.go:67` 逐字在）。
- **阈值是契约常量**：`shortLimit=60`／`longHanziLimit=400`／`longCharLimit=800`／`MaxSpokenRunes=60` 全是 `tier.go` 包内常量，注释写明"contract constants, not tuning knobs；无配置键也不许加"。**⛔ thresholds.go 一字未动**（那是 SLO 的文件）。

## 3 ★★★ 核心行为端到端（真回路读数，`Test30FourResultsThroughTheRealPipeline` 全绿）

真 `agent.Loop`（真 C5 缝＋真 turn 组装＋真 `EvDone`）×四档，每案先 `proc.Reset`、跑完读 `RouteCount`：

| 案 | 走法 | 读数（现跑逐字） |
|---|---|---|
| tier-1 表格＋摘要 | **mockllm 子进程**（openai-chat 真适配器走真 HTTP；非 golden 回声模式把最后一条 user 文本回回来——零工具＋空 scene 时回路不追加 volatile suffix，回声＝样本本身） | `RouteCount==1`；`res.Text` 带 `口播摘要：表格已存档，请查收`；**tier=structured**；`display`＝回声体剥掉末行；**工件真落盘**且内容＝display；播报＝`表格已存档，请查收。全文已存：<全路径>`；剪贴板 0 次 |
| tier-2 402 汉字＋摘要 | **C5 进程内缝假体**（`scriptedC5`：TextDelta→StopEndTurn→Done，`calls` 计数） | `calls==1`；**tier=long**；播报前缀＝`长文已归档`；工件在盘、内容＝剥 marker 后的长文；剪贴板 0 次。★为什么不用 mockllm：**其回声硬上限 400 字节**（`tools/mockllm/chat.go` 的 `len(answer) > 400` 截断——402 汉字＝1,206 B 物理回不来），这是本腿 §8⑥ 的具名发现；缝注入是 AGENTS 1.3 的许可注入面，回路由（loop/router）与另三案**同一条** |
| tier-3 61 汉字 | mockllm 回声 | `RouteCount==1`；**tier=medium**；notify 恰 1、**badge 恰 1、task 恰 1**（双通道无条件）、panel 恰 1（无工件）、播报＝首句 ≤60 rune、剪贴板 0 次 |
| tier-4 短句 | mockllm 回声 | `RouteCount==1`；**tier=short**；剪贴板**恰 1 次且内容恰＝display**、播报＝display 全文、notify/badge/task/panel **全 0** |

全套 `-v`：**19 top-level PASS／0 FAIL／0 SKIP**；`-race` 全绿（`ok … 4.868s`）。

## 4 ★★ 正控（具名行为、不崩）

| 夹具 | 读数 |
|---|---|
| 空白结果 `"   "` | **tier=none**、全通道 0 次调用、receipt 注 `all skipped: empty result: nothing to present`、`OK=true` |
| 全端口 nil（无装配宿主） | **逐通道 skipped `not wired`** 注记、不 panic、`OK=true`（skip 不烂账） |
| 工件配额写失败（注入错误） | receipt 注 `artifact failed: memory: enforce artifacts quota…`、`OK=false`；**面板照常显示、播报改挂诚实尾巴 `。全文未能保存`**——拒绝被具名，不被吞 |
| announce 拒收（`Speak` 返 false） | receipt 注 `speak failed: announce not taken (voice chain not assembled or queue full)`、`OK=false` |
| 工件名 | `result-<taskid>-<hash8>.md` 过 `memory.ValidArtifactName`（含 `..\evil:task` 消毒案）；同内容同名（幂等）、异内容异名 |

## 5 ★★★ 反向判据：三发变异（仓外 mutant ＋ `-overlay`）＋ 还原核（时机＝最后一发之后）

| 发 | 变异 | rc | 逐字红句（节选） |
|---|---|---|---|
| **M1** | `tier.go` `Classify` 首行插 `return TierNone, "mutation M1: routing removed"`（**摘掉路由分发判定**） | **1** | 矩阵全线红；`dispatch_test.go:154: tier = none, want structured`；E2E：`e2e_present_30_test.go:122: tier = none (mutation M1: routing removed), want structured`——四档全塌 |
| **M2** | `dispatch.go` tier-3 的 `if d.ch.Badge != nil {`／`if d.ch.TaskList != nil {` 改 `if false {`（**摘掉双通道**） | **1** | `badge updates = 0, want 1 (dual-channel is unconditional)`／`task list entries = 0, want 1`；trace 显示 `notify failed` 旁没有 `badge sent`；E2E tier-3：`badge=0 task=0` |
| **M3** | `dispatch.go` `route()` 开头（switch 之前）插入无条件剪贴板调用（**剪贴板逃出 tier-4 分支**） | **1** | `clipboard called 1 times on tier 1, want 0: long results never take the clipboard`；tier-3 `clipboard called on tier 3, want 0`；tier-4 `clipboard = 2 calls`；连空结果案都红（空结果也被写了剪贴板） |

★**还原核（三发全跑完后，同一会话）**：①无 `-overlay` 重跑 **rc=0、19 PASS／0 FAIL**；②`git hash-object` 与变异前逐字同值——`tier.go=75f527dd92f63e5b…`、`dispatch.go=9073f5cd6ed71d00…`（变异前后各打一遍，逐字相同）；③树内变异标记 **0 命中**、`git status --porcelain -- internal/present` 只有本腿 8 枚新件。

## 6 ★ 四数 ＋ 两尺

| 面 | 读数 |
|---|---|
| `internal/present` **改前** | `git ls-tree HEAD -- internal/present` ＝ **0 文件**（新包） |
| `internal/present` **改后**（盘上） | **8 文件／51,224 B**（§0 逐枚） |
| staging 全套 | `go test -count=1 -v ./internal/present` ＝ **19 PASS／0 FAIL／0 SKIP**；`-race` 同绿 |
| 邻接包 | `internal/agent` 全套 **rc=0**（`ok … 3.746s`）；`internal/memory -run Artifact` **rc=0**；`go build ./...` **rc=0**；`go vet ./internal/present` **rc=0** |
| 甲 `gofmt -l internal/present` | **空**，rc 0 |
| 乙 `gofumpt -l internal/present` | **空**，rc 0 |
| 全仓两尺（`internal cmd tools`） | 恰 **2 条既存假红**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`（CRLF/`autocrlf` 造，**⛔ 本腿责任、⛔ 本腿修**；先例 26-r1 记的第三条 `cmd\wisp\models.go` 今日不在红名单——那是 `41a2-r1` 的在飞面，本腿只报现状）。本腿 8 枚两尺皆空 |

## 7 ★ 门禁 d22scan（全仓 ＋ 导出树正控，第 21 处定式）

- **全仓**（含本腿 8 枚）：`go run . -root <repo>` **rc=0 `clean`**；分母：bans#1-5 internal/＝**274**、bans#1-5 cmd/＝**45**、ban#8 internal/＝**706**（含注释与 _test.go）、ban#8 cmd/＝**161**、ban#6 frontend/＝116、ban#7 internal/tools/＝25、ban#8 design/＝39。
- **正控**（仓外导出树 `D:\tmp\wisp\20261028\30r1\export`，`git clone --shared` @dev）：植入真 U+2705 字符串 ⇒ **rc=1**，逐字 `internal/present/zz_d22scan_probe.go:5: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)`。★第一发把 `\u2705` 写成了字面转义，门禁报 clean——**顺手证明了"门禁没红"必须追问夹具真假**：转义序列不是字形，重植真字符才咬红。

## 8 ★ 边界／交回父代理（逐条自陈）

1. **`cmd/wisp` 接线具名交回**（本腿库侧做完、零装配）：`present.New(present.Channels{...})` 一次装配——`Speak`→`rt.voiceAnnounce` 的同款闭包（ChTTS 闸在 announce 路径上，`26-r3` HOP 1，本腿不重复设闸）；`Notify`→`rt.notify`；`Badge`→`ball.SetBadgeText`；`TaskList`→面板任务表入口（**票 36 的面**，今天可能还没有真实现，nil＝跳过不炸）；`Panel`→C17 桥的结果展示（**桥本身一字未动**）；`Clipboard`→真写剪贴板的适配器；`Artifacts`→`*memory.Store` **原样即插**（方法集逐字吻合，零适配器）。调用点＝`loop.Run` 返回后对 `res.Status==completed` 的结果 `Dispatch(ctx, present.Input{TaskID: res.TaskID, Text: res.Text})`。
2. **剪贴板泄密闸是适配器义务**：`ClipboardWriter` 文档注释逐字写了**必须**先 `risk.Provenance.CheckText(scope, risk.ChClipboard, text)` 再写真写（闸的脸在 risk 公开面、scope 是宿主事实，本腿不持有）。此为交回清单的一部分，不是已守事实。
3. **AC2 措辞与实现的一处偏差如实报**：票面说 "mockllm request count"——三案真用 mockllm 子进程并读其**服务器侧计数**；**402 汉字案**因 mockllm 回声 400 字节硬上限（§8⑥）改走 C5 进程内缝假体（同 loop 同 router，`calls` 在缝上计）。这是"注入面选错就得撒谎"与"按缝注入"之间的选择，本腿选后者并具名。
4. **tier-3 播报＝首句 ≤60（D10 逐字）**，模型若给了 `口播摘要：` 行会被剥出显示、但不冒充首句——两规则各归各，测试双向钉住。
5. **D10 两处"清单"量词的具名解释**：`文件路径清单` 取 **≥2 行**（单路径在散文里＝tier-4，否则"文件在 C:\x\a.txt 里"会被迫落文件、违背"一句话短结果"档）；`≥2 个链接` 计 **markdown 链接＋裸 http(s)**（先剥 markdown 跨度再数裸 URL，href 不二次计数——单链接其 href 是 http 的陷阱案在测试里）。
6. **发现（非本腿缺陷）**：mockllm 非 golden 回声模式 `len(answer) > 400` 截断在 **字节** 上（`tools/mockllm/chat.go`）——多字节文本被拦腰切。归 mockllm 自己的票，本腿只登记不修。
7. **零 goroutine／零墙钟超时／零裸 `go func(`**：`Dispatch` 全同步跑在调用方 goroutine；慢通道是宿主适配器的队列（D38b），本包注释逐字声明。零配置键、零 emoji、`filepath` 只用 `Join`（无 `Clean/Abs`）。
8. **`DEFERRED` 面**：本包未引入任何 `DEFERRED` 标记，无需登记表动作。

## 9 交回父代理

**未定义即停：未触发**（①写面与硬边界无正面冲突——装配不需要动 `cmd/wisp`，端口即缝；②票面未要求改 C17/D10 契约文字；③票面前提全在（§2 现核）；④引文出处无错位）。★开工三核的 HEAD 偏差已在页首具名（锚点过期，未切任何东西）。★**六格全部落在写面内**，AC6 带 §3 的具名边界（真球/真音频读数随 `cmd/wisp` 接线腿交回）。
