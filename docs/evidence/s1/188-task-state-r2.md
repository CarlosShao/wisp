# 188-r2 — 后台任务状态维落到 `TaskOutput`（派单 §G，裁定 `A394` (i)）

写腿子代理，2026-09-28 18:1x +08。工单＝票 188 `AC#2`＋`AC#3`。派单＝
`.scratch/wisp/dispatches/2026-09-28-181x-wave2-impl-145r3-and-188r2.md` §G。
起手自取锚（`git rev-parse --short HEAD` 现量）＝`fccaf3e3`；本程第一枚 commit `62d7d3b4`
实际落在并发波次的 `4db5f3f6`（§F `145-r3` 格②）之后——共享工作树，非我改写。

**AC 框一枚未勾**（勾要非实现者表，`SPEC-12 §4.3` #1/#3）。

---

## 1. 字段形状与它取的 D43 名字（逐字出处）

| 件 | 行 | 内容 |
|---|---|---|
| `internal/tools/task.go` | `:124` | `State statemachine.State`（`TaskOutput` 的第三枚字段；旧两枚 `Text:97`／`ArtifactPath:101` 一字未动） |
| `internal/tools/task.go` | `:142` | `func (o TaskOutput) StateAnswer() (statemachine.State, string)`——读侧唯一出口，空值与非法值各回一句 fail-closed 文案 |
| `internal/tools/task_backfill.go` | `:75` | `State statemachine.State`（`TaskBackfill` 的新字段，零值＝今天的生产形状） |
| `internal/tools/task_backfill.go` | `:113` | 写侧门：`case b.State != "" && !statemachine.Valid(b.State)` ⇒ 整条记录**拒绝落表**并带原因 |
| `internal/tools/task_backfill.go` | `:127` | `rec := TaskOutput{Text: text, State: b.State}` |

类型不是 `string`：`statemachine.State` 的取值集就是 D43 的仓内逐字拷贝——
`internal/statemachine/states.go:11-31`（`FirstRun`/`Sleeping`/`Armed`/`Muted`/`Listening`/`Thinking`/
`Acting`/`Speaking`/`Warm`/`Conversation`/`Confirming`/`AwaitingApproval`/`Settling`/`Downloading`/
`Unconfigured`/`NoNetwork`/`Error`/`Queued`/`Stuck`/`WatchdogAlert`，共 20 枚），
该文件标题行 `:8` 自陈 "names exactly as in D43 / SPEC-08"，判定者 `Valid` 在 `:39`。
D43 本体＝`docs/PLAN.md:3055` 标题＋`:3061-3100` 的 40 条转移，**本程零字节未动**（`git diff` 空）。

本程用到的取值只有测试里的 `Acting`／`Error`／`Settling` 三枚，逐字取自上面那张表——
**没有新造任何名字**，`states.go` 也未被本程改动一字。

**`run.go` 未被本程写入**（§F 禁区）。它的生产调用点 `cmd/wisp/run.go:640-647` 用 keyed literal 构造
`TaskBackfill{Roster, Spills}`，所以新字段编得过、行为逐字节不变（`go build ./cmd/wisp/ ./internal/agent/` rc=0）。

## 2. `AC#3` 那枚钉在哪＋正控那一发红在哪条用例

判据文件＝`internal/tools/task_state_188_test.go`（本程新建，465 行）。

- **钉子本体**＝`TestTaskState188AC3PanelHasNoWriteLeg`（`:341`），尺＝`scanTaskStateWriteLegs`（`:249`）：
  AST 扫 `internal/panel/` 的**非测试** Go 源（`filepath.WalkDir`＋`go/parser`，跳过 `_test.go`），
  四种"写这一维"的形状——① 调 `Record`／`Backfill`（`rosterWriteMethods:241`）、
  ② 构造带 `State` 键的 `TaskOutput`（含 `tools.TaskOutput` 限定名）、③ 给 `.State` 赋值、
  ④ `Method*` 白名单常量的值声称 `task`／`state` 轴。
  **分母尺**：解析到的非测试文件数 <10 直接 `Fatalf`（`:346` 现量＝**11 枚**，0 命中）。
  ⇒ "没有发现"不可能等于"什么都没读"。
- **正控（五枚假腿，逐枚必须红）**＝`:392-425` 的 `plants` 表，种在 `t.TempDir()`（**仓内零新文件**），
  实测逐枚红：`roster-record-call`／`roster-backfill-call`／`taskoutput-state-key`／
  `state-field-assign`／`whitelist-method-claim`（`--- PASS` 那次跑里五枚 `t.Logf` 全打印了红因）。
  还带**分母自证**：每枚假腿所在目录必须 `examined==1`，否则 `Fatalf`（尺没吃到这枚假腿就报错，不算绿）。
- **反向控（门不能 aimed 太宽）**＝`:430-453`：同一把尺扫一枚"面板只读 `Look`／`Count`"的读腿 ⇒ **必须不红**。
  这一格是留给票 145 的快照泵读口的：面板只许显示＋发起请求（R20／票 92），
  把读腿也点掉就是本票自己的下游死在自家门上。
- **生产者那一半**＝`TestTaskState188AC2BackfillIsTheGoSideProducer`（`:170`）：Go 侧写腿存在且带门。
  (a) 零 `State`＝生产形状，记录落表、`StateAnswer` 说"未登记"；
  (b) `State: statemachine.StateActing` ⇒ 落表；
  (c) 逐枚喂旧词 ⇒ **拒绝落表**＋原因里带上那枚词与 `D43`＋`roster.Look` 查不到＋`Count` 不变。

## 3. 有没有把 `done/cancelled` 那族渗进来（应为零，给尺）

零。三把尺：

1. `legacyTaskStateWords`（`:96-101`）＝ `done`/`cancelled`/`running`/`succeeded`/`interrupted`/`error`/
   `finished`/`failed`/`stuck`/`queued`/`pending` 十一枚；`TestTaskState188AC2VocabularyIsD43Only`（`:76`）
   逐枚要求 `statemachine.Valid` 为假**且** `StateAnswer` 给出拒收文案（含 `D43` 字样）。
   同一把尺先枚数钉死 20 枚 D43 名字、逐枚要求被接受——两向都响，不是单边期望。
2. 写侧：上面 2(c) 的十一枚逐一喂进 `Backfill` ⇒ 全部被拒（旧词进不了这张表）。
3. 现量尺（证明我没有顺手"统一"memory 那一维）：
   `git diff -- internal/memory/dao_test.go internal/memory/privacy_test.go internal/agent/forensics_test.go internal/agent/loop_golden_test.go cmd/wisp/run.go` ⇒ **空**；
   `./internal/memory/` ok 13.319s、`./internal/agent/` ok 1.784s（那四组钉着的断言颜色未变）。
   `TaskLog.State` 与 `TaskOutput.State` 之间**一张映射表都没画**（票 196 射程，`A394` 令本票不统一）。

## 4. 门禁读数原样

```
go test -count=1 ./internal/tools/ ./internal/memory/
  ok  github.com/CarlosShao/wisp/internal/tools   17.675s
  ok  github.com/CarlosShao/wisp/internal/memory  13.319s        （上一程 12.810s，同色）
go test -count=1 ./internal/agent/                 ok  1.784s
go build ./cmd/wisp/ ./internal/agent/             rc=0
gofumpt -l internal/tools/ internal/memory/        空（本程判据先被点了一次，-w 后复跑全绿）
sh scripts/d22scan.sh                              clean - no D22 ban violations
  ban #8 internal/ examined 442 Go files（上一程 441 ⇒ 增量＝本程判据一枚）
  bans #1-5 internal/=211 cmd/=24 ban #6 frontend/=85 ban #7 internal/tools/=21 ban #8 cmd/=47 design/=39
```

`internal/panel/` 今天在册红**三枚，逐名照实记、未修未当绿**（本程写面外）：

```
--- FAIL: TestComposerContractTypesMatchFrontend
--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme
--- FAIL: TestC21DesignTokensFourWayAgree
FAIL github.com/CarlosShao/wisp/internal/panel  1.205s
```

`gate-clauses`：派单写的 `scripts/gate-clauses.sh` 不存在（`No such file or directory`），
尺实际在 `.scratch/wisp/probes/154/gate-clauses.sh`，本程跑到了：BAD 名册**仍只 `G6neg` 一枚**，
但 `声明=ring 基线=1枚 实测=3枚 因=新增未成对`（`OpenTask`↔`CloseTask` 不过滤文件那一味，`:466-473`）。
本程三枚文件 `OpenTask` 命中 0（`grep -rln "OpenTask\|CloseTask" --include=*.go internal/ cmd/` ⇒
`internal/tools/bridge.go`／`bridge_scope_open_ticket158_test.go`／`pointer_183_cli_seam_test.go`／
`pointer_185_cli_seam_test.go`／`ticket175r2_stamp_live_test.go`／`cmd/wisp/run.go`／
`cmd/wisp/task_scope_close_151_test.go`）⇒ 那两枚增量**不在本程写面**，本程不修不当绿，
也未在"起手锚点减本程改动"的树上复算，所以归属只到"不是我写的字面量"这一层。

## 5. 面板要的键清单（只写在这里，界面侧由 owner 自己带）

`AC#4` 分账：本票交"源"，快照载体归票 145。若那一格将来要显示后台任务状态，可读的真源是
`TaskRoster.Look(id)` 拿到的 `TaskOutput.State`，并且**必须一起带上 `StateAnswer()` 的第二返回值**
（"未登记"／"非 D43 名"两种形状），否则界面会把"宿主没判过"渲染成"已完成"。
本程**没有**给面板加任何写腿，`internal/panel/` 非测试码这一维命中仍是 0。

## 6. 纪律面

- **预算：≤30 枚调用，实耗约 36 枚＝超 6**。超支具名：`gate-clauses` 与 `scripts/` 路径定位 3 枚、
  `gofumpt` 不在 PATH 改走 `go env GOPATH` 重取 2 枚、判据自身两发缺陷各修各验 4 枚
  （`ast.Inspect` 的 nil 收尾节点 panic；`tools.TaskOutput` 限定名让第三枚正控假绿）、
  起手名册核对与首枚 commit 失败重来 2 枚。
- **首枚 commit 失败一次**：未跟踪的判据文件未 `git add` ⇒ `error: pathspec ... did not match any file(s) known to git`，
  **rc=1，没有落任何 commit**（`git log` 其时顶枚仍是 §F 的 `4db5f3f6`）；补 `git add internal/tools/task_state_188_test.go`
  单枚显式路径后重提，未误带别家 staged 件。
- commit 全部显式 pathspec；禁 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean` 全部遵守；
  **仓内零删除**（`git show --stat` 里那 1 处删除是 `task_backfill.go` 被替换的那一行赋值，不是文件）；**只 commit 未 push**。
- `frontend/**`／`design/**` 写面 0；`go.mod`／`go.sum`／`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／
  `internal/panel/bridge.go:42-45` 一字未动；`-overlay` 未用（正控载体走 `t.TempDir()`，仓外）。
- 被拒调用：0 次。

## 7. 没测到什么

1. `TaskOutput.State` 的**生产写入者仍为零**：`cmd/wisp/run.go:640-647` 没带 `State`（本程禁写 `run.go`）。
   因此今天任何真任务的这条记录读起来都是"宿主没有登记这一维"，这是**诚实的空**，不是"任务已完成"。
2. 没测这一维进快照后的样子（归票 145）；没测 `task.output` 的模型可读文本带这一维（stub 措辞按
   `task.go:251` 是 PLAN.md:2564／票 164 `AC#3` 冻结件，本程一字未改）。
3. 门的覆盖边界：不经 `Record`／`Backfill`／`TaskOutput` 字面量／`.State` 赋值／`Method*` 常量的触碰方式
   （例如反射、或面板自建一套同义载体）不在这把尺的射程内。
4. `AC#1` 子代理那一维没碰（`A389` 已摆 owner，票 197 另立）。
5. 两套词表的映射表一张都没画（票 196／`A394` 令本票不统一）。
6. 未在"起手锚点减去本程改动"的树上复跑 `gate-clauses`，所以 `G6neg` 那两枚增量只做到"排除本程字面量"，
   没做到逐枚归人。
