# 197-a2 §A — 五枚未勾格逐格现状（只读普查；锚 HEAD `8e4f6f83`／分支 dev）

> 只读腿，⛔ 零产码、⛔ 零 go 命令、⛔ 未翻框未动 Status。每行＝结论＋尺（命令现跑）。

## A0 票面那句「Go 侧零子代理实体」今天真伪 — **早已过期**

| 事实 | 现跑读数 | 尺 |
|---|---|---|
| Go 侧「子代理实体」非测试文件 | **9 枚**：`internal/panel/composer.go` · `internal/panel/pump.go` · `internal/panel/subagent_roster_197.go` · `internal/streamkey/streamkey.go` · `internal/tools/bridge.go` · `internal/tools/subagent_197.go` · `internal/tools/task.go` · `cmd/wisp/panel_pump.go` · `cmd/wisp/run.go` | `grep -rln "subagent\|Subagent\|SubAgent" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go \| wc -l` |
| 中文「子代理」非测试 | **3 枚**：`internal/panel/subagent_roster_197.go` · `internal/tools/subagent_197.go` · `internal/tools/task.go` | `grep -rln "子代理" --include=*.go internal/ cmd/ \| grep -v _test.go \| wc -l` |
| `TaskRoster` 方法枚数 | **13 枚**（票面锚 `Record:123` 已漂到 `:215`）：`Record:215` `WatchRow:248` `Look:268` `Count:281` `Descendants:305` `PublishSubagent:342` `MarkRoot:355` `TryAcquireSubagentSlot:363` `InFlightSubagents:390` `RunningSubagentIDs:402` `AttachCancel:420` `DetachCancel:434` `Cancel:447` | `grep -n "func (r \*TaskRoster)" internal/tools/task.go` |
| 流键单点铸造 | 字面量唯一住 `internal/streamkey/streamkey.go:29` `SubagentPrefix = "subagent:"`；`internal/panel/pump.go:483` 与 `internal/tools/subagent_197.go:47` 各为 alias | `grep -rn "subagent:" --include=*.go internal/ cmd/ \| grep -v _test.go` |

⇒ 票名尾巴那句「go-side-has-zero-subagent-entity」是**历史状态**（与票面 §「要 owner 带给界面那支的一跳」第一条同口径）。

## A1｜AC#0 先答「这一层归哪个切片」— 未勾（答在票面正文，缺非实现者裁决表）

- **今天真有什么**：票面 §「AC#0 的裁定」已逐字答＝**S5（Agent 内核）为家、S7（并发）为闸**；并判「`:1447` 那行描述的是早已过去的一条时序，不是一条还在生效的禁令」。
- **缺什么**：一枚**非实现者**裁决表把这一格落成凭据（`docs/evidence/s1/` 下今天无 `197-*` 对应件）。⛔ 这一格不欠 Go 产码。
- 尺：票面第 21–24 行逐字；`docs/evidence/s1/` 名册（`ls docs/evidence/s1/ \| grep 197`）。

## A2｜AC#1 实体层有生产调用者 — 未勾（**盘上已有生产调用者**）

- **今天真有什么**：`cmd/wisp/run.go:780` `for _, e := range tools.BuiltinSubagentEntries(tools.SubagentDeps{…})` → `:791 reg.Register(e)`＝`task.spawn` 进真注册表；依赖束是**运行中的对象**（同一枚 `rt.tasks` 名册、同一枚 `rt.bridge`、`rt.c25`、`rt.stream`、`rt.loopOpt`，源码注释 `:770-779` 逐字）。
- **res 常驻进程也走同一条装配**：`cmd/wisp/resident_task_source_windows.go:462` `code := src.run.execute(root.Ctx, text)` ⇒ 面板/球的常驻进程共用 `run.execute()`，故 `task.spawn` 在常驻侧**可达**。
- **缺什么**：非实现者裁决（AC#0/AC#1 同批）。
- 尺：`grep -rn "BuiltinSubagentEntries" --include=*.go internal/ cmd/ \| grep -v _test.go`＝仅 `run.go:780`；`sed -n '770,800p' cmd/wisp/run.go`。

## A3｜AC#2 名册与状态同源 — 未勾（**盘上同源**）

- **今天真有什么**：`internal/tools/subagent_197.go:99-104` 四枚态**只取自 D43**：`subagentStateRunning=statemachine.StateThinking`／`Settled=Settling`／`Stopped=Muted`／`Failed=Error`；有钉 `internal/panel/subagent_roster_197_test.go:246 TestAStatusOutsideD43NeverReachesTheWire`。
- **缺什么**：非实现者裁决；判据要能区分「填了真状态」与「字段恒空」（票 181 `AC#7` 那枚空转教训）。
- 尺：`sed -n '95,110p' internal/tools/subagent_197.go`；`grep -n "^func Test" internal/panel/subagent_roster_197_test.go`。

## A4｜AC#3 每子代理一条流、不许合并 — 未勾（**r2 已把语义重裁并落钉**）

- **今天真有什么**：`197-r2`（`94ea50f7`）把溢出从「合并最老两行」改成「显式截断每条流自己的中段」；结构判据 `internal/panel/subagent_stream_197_test.go` 五枚，其中 `TestNPlusOneSubagentsEachGetTheirOwnStream:57`＝票面第 2 层要的那发「N＋1 枚没有任何一条被并进别人的」。
- **缺什么**：非实现者裁决（「N+1 不并流」这枚判据的对抗验收）。
- 尺：`grep -n "^func Test" internal/panel/subagent_stream_197_test.go`；`git show --stat 94ea50f7`。

## A5｜AC#6 取消语义 — 未勾（**两半已做，第三半拆票；欠一枚 goroutine 名**）

- **做到了的两半**：① 父取消不级联＝`internal/tools/subagent_197.go:327` `context.WithCancel(context.WithoutCancel(ctx))`；② 这件事写给模型看了＝`subagent_197.go:401` 结果文本含「父任务 … 还能用 task.cancel 单独停它」。**无裸 `go func(`**（`197-c2` §① 现跑零命中）。
- **孩子停止出口这半——今天已由票 221 落地**（且**不在本票格内**）：`internal/tools/task.go:600 BuiltinTaskEntries` 现返回 `taskOutput`＋`taskCancel` 两枚；`TaskRoster.Cancel` 现**有生产调用者**＝`internal/tools/task.go:742` `stopped, why := t.d.Roster.Cancel(target)`。⛔ 票 221 文件名带 `-done`（防重领键），其 `Status:` 行写「待派」＝**过期 Status 行**，以盘上码为准。
- **⚠ 仍欠（本票 AC#6 链上的优点）**：收尾协程的名字 `subagent-finish-<id>` **不在 D38 名册里** ⇒ `grep -c "subagent" internal/observe/goroutine.go`＝**0**；spawn 点 `internal/tools/subagent_197.go:383-384`（`observe.NewRootFrom` / `observe.Default.Spawn("subagent-finish-"+bg.ID, "tools", …)`）。⇒ 已立**票 229**（`待派`），后果＝唯一泄漏报警被日常噪声淹掉。
- **缺什么**：非实现者裁决；且 229 那枚名字要么改到名册内形状（甲形，可直接做）、要么动 D38b 表（乙形，须人工批准）。
- 尺：`grep -c "subagent" internal/observe/goroutine.go`；`sed -n '735,745p' internal/tools/task.go`；`sed -n '597,604p' internal/tools/task.go`。

## A6｜附：票 220 那格（AC#4 已勾但只勾了「有一格」那半）

- 票 220（`待派`，**0 勾／5 未勾**）＝名册 `blockedOnApproval` **只报 L2**、同一枚任务第二张卡静默读成「没被卡住」。⇒ 票 197 `AC#4` 勾的是「载体里有这一格且拿得到」，**不是**这一格读得对。
- 尺：`head -5 .scratch/wisp/issues/220-*.md`；`grep -c "^- \[x\]" .scratch/wisp/issues/220-*.md`。

## A7｜附：界面那一跳今天仍是**红的**（跨边界，归 owner，⛔ 本编队写面）

- Go `Snapshot` 现发 **6 枚键**：`internal/panel/composer.go:57-91`（`pending` `results` `composer` `generatedAt` `instructions` `tasks`）。
- 界面 `PanelSnapshot` 只声明 **4 枚键**：`frontend/src/lib/panel.ts:129-137`（`pending` `results` `composer` `generatedAt`）⇒ 缺 `instructions`／`tasks`。
- 判据 `internal/panel/composer_test.go:48 TestComposerContractTypesMatchFrontend`（`:73` 双向对账、`:74` 红句「Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare」）⇒ **今天为红**（`10-13` 合并的 `36b44ee4` 那棵树 panel.ts 仍缺这两枚）。⛔ 恢复绿只有一跳＝页面侧补 `tasks?:`／`instructions?:`，Go 侧零写面。
- 尺：`grep -n "interface" frontend/src/lib/panel.ts`；`sed -n '122,140p' frontend/src/lib/panel.ts`。
