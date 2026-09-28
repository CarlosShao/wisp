# 180-c1 ＋ 182-c1 — 面板那一维的 Go 侧现量普查（只读·零产码）

> 性质：**只读普查**。本件只交「界面上那一维，Go 里到底有没有真值、有没有读者、归属哪枚票」的
> 现量与代价。**不勾任何 AC 框、不动任何产码、不选落点、不提"顺手把批准接上"**
> （由面板侧来源的 L2「允许」是已定案禁止项，`AGENTS §1.2` ＋ 台账 `Q-49` 那一族）。
> 权威来源：`.scratch/wisp/dispatches/2026-09-28-145x-readonly-180-c1-plus-182-c1-panel-field-census.md`
> ＋ 票 180 ＋ 票 182。本件不新造规矩。

## 1. 起手锚＋写面闸门

- 起手时刻：`2026-09-28 14:40 +0800`
- 起手 HEAD：`e9ef94d0`（Mon Sep 28 14:35:19 2026 +0800）
  ⚠ 派单写的编排者锚点 HEAD＝`03c01d71`；起手实测 HEAD 已是 `e9ef94d0`（共享工作树，别人先走了）。
  本程所有行号一律**现跑**，不抄任何票面历史号。
- 起手写面闸门：`git status --porcelain -- internal/ cmd/` ⇒ **空输出**（干净）。
- 终态写面闸门：见 §10 末（待落）。

## 2. 票 180 AC#1 — 面板窗口尺寸归属三问

三把尺全部现跑（锚 `eba0b89b` 之后的工作树，本程 14:4x）。

**① 配置里那枚字段**（尺 `grep -n "width\|Width" internal/config/schema.go` ⇒ 命中 2 行）
- `internal/config/schema.go:527-528` 逐字：
  `// Width is the panel width in px.` ＋ `Width int \`toml:"width" default:"640"\``
- 所在结构体：`PanelSection`（`internal/config/schema.go:525-535`，五枚键全在场：
  `Enabled` `Width` `Height` `KeepAliveInSession` `Scale`）。

**② 生产码里谁读它**（尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go`）
- ⇒ **空输出（rc=1）**，命中 **0 枚**。票面 `039efb47` 的"读者＝零枚"到今天**仍然成立**。
- 顺手把同段另三枚也跑了：`grep -rn "\.Height\b\|\.Scale\b" --include=*.go internal/ cmd/ | grep -v _test.go`
  ⇒ **唯一命中 `internal/agent/budgets.go:94` 的 `b.Scale`**，那是 `CostSection` 的预算缩放，
  与 `PanelSection.Scale` 无关 ⇒ `Panel.Height`／`Panel.Scale` 生产读者同样 **0 枚**。
- `PanelSection` 整体在生产里**唯一**被碰到的地方是热重载的**通用拷贝**：
  `internal/config/manager.go:211` 逐字 `{"panel", &cur.Panel, &fresh.Panel, func() { cur.Panel = fresh.Panel }},`
  ⇒ 它按指针搬整枚 struct，**不读 `Width` 的值、也不产生任何效果**。

**③ 窗口／WebView 的实际尺寸今天从哪来**
- 尺 `grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test` ⇒ **空输出**。
- 尺 `grep -rn "WebView2\|CreateWindow\|bounds\|Rect" --include=*.go internal/panel/ internal/winsec/ cmd/ | grep -v _test.go`
  ⇒ 命中只有两类：(a) **注释**说这棵树里根本没有宿主；(b) `cmd/wisp/notify_windows.go:149` 的
  `pCreateWindowExW.Call(...)`——同一函数 `:153` 的错误串逐字 `CreateWindowExW(message-only): %v`，
  即那是**通知用的 message-only 隐藏窗口**，不是面板。
- 注释现量（逐字，三处）：
  - `internal/panel/doc.go:16`：`// DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge).`
  - `internal/panel/pump.go:16`：`// tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer,`
  - `cmd/wisp/run.go:438`：`// have is a page to go to: this tree carries no WebView2 host (tickets 33/35),`
- ⇒ **今天没有任何人决定面板尺寸，因为面板窗口本身不存在**（票 33 的 host、票 35 的 bridge 都还是
  `ready-for-agent`；`ls .scratch/wisp/issues/` 里 `33-panel-host-c27.md`／`35-panel-bridge-c17.md` 均**无 `-done` 后缀**）。

**②→③ 之间缺的不是"一跳"，是整条链**：`cfg.Panel.Width` → 宿主 → `SetWindowPos/尺寸`。
中间那一环（WebView2 host）从未落地，所以 `Width` 与 `Height`/`Scale` 一样是**悬空字段**。

**三档结论：＝「规格要求生效」那一档**（⇒ 按票 180 AC#2，正解是把那一跳接上；但**接哪一跳归编排者选落点，本程不选**）。
规格原句（不转述，逐字）：
- `docs/PLAN.md:2726`：`**三档生效级别**：\`hot\`（立即生效）· \`reload\`（需重载子系统，如换 ASR 模型）· \`restart\`（需重启进程）。`
- `docs/PLAN.md:2743`（D36 配置树那一行）：
  `| \`[panel]\` | \`enabled\` \`width\` \`height\` \`keep_alive_in_session\`(true) \`scale\` | \`hot\` |`
- `docs/specs/SPEC-03-config-secrets-envs.md:39`：
  `| \`[panel]\` | \`enabled(bool)=true\` \`width(int)=640\` \`height(int)\` \`keep_alive_in_session(bool)=true\` \`scale(float)\` | hot |`

⇒ `width` 被明确列进 `hot`，而 `hot` 的定义逐字就是「立即生效」⇒ **规格要求它生效，本票属"补实现"支，不是"改规格文字（＝人工批准）"支**。
反向自查（规格有没有另一处只要求"显示"）：尺 `grep -rn "width\|Width" docs/specs/SPEC-03*.md docs/specs/SPEC-08*.md`
⇒ SPEC-08 里唯一命中是 `SPEC-08-ui-ball-panel.md:218` 的 SVG `stroke-width`，**与面板窗口无关**；
`SPEC-08` 全文**没有一处**给面板窗口定尺寸。

## 3. 票 180 AC#4 — 同族「带 `default` 但生产零读者」枚数与逐枚名

尺（现跑，逐枚）：对 `internal/config/schema.go` 里每一枚带 `default:` 的字段名跑
`grep -rn "\.<字段名>\b" --include=*.go internal/ cmd/ | grep -v _test.go`，取命中数。
原始读数落 `.scratch/wisp/probes/180/c1/same-family-census.txt`（61 行）与
`probes/180/c1/field-to-struct-map.txt`（69 行，字段→struct 映射）。

- 基数：`schema.go` 里带 `default:` 的**行＝69**，去重后的**字段名枚数＝61**；
  非测试 Go 文件总数＝231（`internal/` ＋ `cmd/`）。
- **零生产读者（total=0）＝18 枚**，逐枚名（限定名，映射表现跑）：

| # | 限定名 | 所在 struct |
|---|---|---|
| 1 | `PanelSection.Width` | `PanelSection`（票 180 本尊） |
| 2 | `PanelSection.KeepAliveInSession` | `PanelSection` |
| 3 | `BallSection.ClickThrough` | `BallSection` |
| 4 | `BallSection.HideOnFullscreen` | `BallSection` |
| 5 | `SessionSection.WarmTimeoutSec` | `SessionSection` |
| 6 | `SessionSection.SettlingSec` | `SessionSection` |
| 7 | `SessionSection.ConversationIdleSec` | `SessionSection` |
| 8 | `AudioSection.InputDevice` | `AudioSection` |
| 9 | `AudioSection.SampleRate` | `AudioSection` |
| 10 | `PrivacySection.DiagnosticsOptIn` | `PrivacySection` |
| 11 | `PrivacySection.RetentionDays` | `PrivacySection` |
| 12 | `MemorySection.L1Enabled` | `MemorySection` |
| 13 | `MemorySection.L1Max` | `MemorySection` |
| 14 | `MemorySection.L3RetentionDays` | `MemorySection` |
| 15 | `AECConfig.EchoRef` | `AECConfig`（`voice.aec.echo_ref`） |
| 16 | `RollConfig.SizeMB` | `RollConfig`（`observe.roll.*`） |
| 17 | `RollConfig.Days` | `RollConfig` |
| 18 | `ObserveSection.SLOSampleIntervalSec` | `ObserveSection` |

- ⚠ **18 是下界不是全量**：粗尺 `\.<名>` 对**同名异处**的字段会**多报**读者
  （`Enabled`/`Mode`/`Size`/`Level`/`Max`/`Provider` 这类泛名，命中很可能来自别的 struct），
  所以"非 0"不等于"这一枚有读者"；反过来"0"是**铁的零读者**。⇒ 真实枚数 **≥18**。
  票 83 那张表用 I-A~I-G 七把仪器处理的正是这个盲区（见下）。
- **逐枚归口＝不用另立名单：这 18 枚全部已在「票 83（done）」的 AC#1 全量键表里登记过**。
  尺 `grep -rliE "<键名>" .scratch/wisp/issues/*.md` ⇒ 命中 `83-config-keys-that-lie-must-fail-loudly-done.md`，
  逐字（`issues/83-...md:155`）：
  `| \`panel.{enabled,width,height,keep_alive_in_session,scale}\` | \`schema.go:517-525\` | 无（I-A） | 📋 票 33/34/36 |`
  其余各枚同样在表内有行：`:131` ball 两枚 → 📋 票 62/65/68 ＋ 票 39；`:133` session 三枚 → 📋 票 28；
  `:134` `voice.*`（含 aec）→ 📋 票 15/26/27/41/59/60/61；`:135` audio 两枚 → 📋 票 13 尾巴／票 26；
  `:152` privacy 两枚 → 📋 票 45／票 08；`:154` memory 三枚 → 📋 票 29；`:161` observe/roll 三枚 → 📋 票 08/42/66。
  ⚠ 该表引的 `schema.go:517-525` 等行号**已被注释级改动推漂**（今天 `PanelSection` 在 `:525-535`），
  本程一律用现跑号，不沿用票 83 的历史号。
- ⇒ **对票 180 现量第 3 条的复核结论**：`internal/config/unwired_test.go` 的"名册"**射程只有四枚 locked section**——
  `unwired_test.go:311 TestEveryLockedSectionKeyIsAccountedFor` 里 `collect(...)` 只喂
  `risk` / `fs` / `net` / `plugins`（`grep -n "collect(" internal/config/unwired_test.go` 现量），
  守卫表 `lockedKeyDisposition` 在 `internal/config/unwired.go:119`；`:217` 那行 `width = 800` 属
  `TestUnwiredGuardIsScopedNotABigStick` 的"诚实配置不该报错"样本，**不是登记**。
  ⇒ `[panel] width` 不在名册射程内是**设计如此**，本票第二格因此**既不是"补登记"、也不是"新发现"**：
  发现与归口都已在**票 83 那张表**里；真正缺的是**把效果接上**那一格（AC#2，落点归编排者）。

## 4. 票 182 AC#1 — 「任务监控」那一栏的堆数与逐堆名

- 堆数（现量）：TBD
- 逐堆名 ＋ 「我在谁那儿看到的这一堆」：TBD

## 5. 票 182 AC#2 — 逐堆三档定性（有源／无源但规格要求／规格真空）

TBD

## 6. 票 182 AC#3 — 逐堆归属指认表（含「没票认领」那一列）

TBD

## 7. 票 182 AC#5 — 雷区清单（逐堆：有无面板发起的批准／写动作）

TBD

## 8. 面板要的字段清单（交编排者转前端会话）

TBD

## 9. 本程没测什么（逐名）

TBD

## 10. 门禁终态

TBD

## 11. 被拒调用＋零删除自证＋工具调用终值

TBD

## 12. next＝落地腿还缺什么、哪几枚要人先批准

TBD
