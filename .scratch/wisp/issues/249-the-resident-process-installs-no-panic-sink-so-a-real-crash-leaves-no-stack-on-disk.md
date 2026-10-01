# 249 — 常驻进程的 panic 记录面**没装**：协程观察者今天兜得住炸、但盘上不会留下任何栈

Status: OPEN（编排者 10-01 **13:13** 立，料全部出自探针腿 `33-p1` 与本编排者 13:1x 的现跑尺；台账 `A502`）
实现者：`249-r1`（写码腿，**按住**，见「排程」） · 裁决者：`249-v1`（必须 ≠ 实现者） · 归口：可观测性／票 08 的 JSONL 记录面

## 为什么现在立这枚票（每条带尺与取数时刻；⚠ 引用前先重跑，别把这几行当常量）

| # | 事实 | 尺与出处 |
|---|---|---|
| 1 | 协程观察者**有** recover：`internal/observe/goroutine.go:294-326` 里那枚 `defer` 会 `recover()`＋取 `debug.Stack()`＋交给 sink。⇒ 进程不会死，那条协程 return | 〔现读码〕＋`33-p1` R37/R38 真跑复认（`rc=0`、`PANICS=1`、`STA-HANDLE-RETURNED err=... nil pointer dereference`） |
| 2 | **但 `SetPanicSink` 生产零调用者** ⇒ `:314` 那句 `if sink != nil` 走空，**栈不落盘**。出货时炸一次，盘上只剩 `:317-320` 那句**含 recover 值、不含栈**的 error 文本 | 〔我 13:12 现跑〕`grep -rn SetPanicSink --include=*.go cmd internal`（排 `_test.go`）⇒ **只命中定义自身 `internal/observe/goroutine.go:240-241`** |
| 3 | 装了两行 sink 之后，原始栈确实能拿到——**现成形状可抄**：`%TEMP%\wisp33p1\reentry-quitduringcreate-sta\panic-sink.txt`，头帧逐字 `runtime/debug.Stack()` → `observe.(*Registry).run.func1() goroutine.go:302` → `panic(...)` → `edge.(*Chromium).Init() chromium.go:131` → `Embed :112` → `CreateWithOptions webview.go:340` → `NewWithOptions :109` | `33-p1` R37/R38（**人造形状**：它要先 `PostThreadMessageW(tid, WM_QUIT)` 才走到那一支；自然形两发 `panics_total=0`，见 R34/R35） |
| 4 | 这一格与"跑任务的进程里有球"是同一族的形状：**机制齐、生产没接**（同形先例＝语音链路的 `Sink` 与 `SetLoaded`、权限侧的 `SealDir` 生产零调用者） | 项目记忆「240-c1 三分类」／票 132／票 247 现量表 |

## What to build

装配根（`cmd/wisp` 那侧建 `Registry` 的地方）把 panic sink **接上真记录面**：一条协程炸掉时，`recover` 值＋`debug.Stack()`＋协程名＋owner＋RFC3339 时刻，**落到数据根内的日志**（票 08 那条 JSONL 管线，⛔ 不新造落点、⛔ 不写 `%TEMP%`）。接完之后，owner 在我这台机器上炸一次，我**能在盘上找回那枚栈**——这就是本票的唯一目的。

## Acceptance criteria（⛔ 勾归非实现者；产码腿一枚都不许碰复选框）

- [ ] **AC#0**：先把落点摆开再动手——sink 写去哪、谁给它句柄、崩在**极早期**（日志面还没起来）时那一发怎么办（⛔ 不许默默丢；要么有具名兜底、要么在票面登记"这一窗口的栈拿不到"）。
- [ ] **AC#1 生产者半有人真调用**：`SetPanicSink` 的**非 `*_test.go` 调用者 ≥1 枚**，且调用点在装配根（尺：`grep -rn "SetPanicSink" --include=*.go cmd internal` 排 test 后**必须不再是只命中定义**）。⚠ 只"import 了 observe"不算（记忆里那条"imported ≠ assembled"）。
- [ ] **AC#2 会响的钉（能力型，⛔ 词面型）**：一发用例让一条**在册**协程真 panic，断**盘上那枚记录里有栈**——要含 `debug.Stack()` 才有的多帧形状（至少 `panic(` 那一帧＋一个函数帧），⛔ 只断"有一行 error 文本"不算，那正是没装 sink 时也会有的东西。
- [ ] **AC#3 反向证仪器有牙**：同一枚用例在**不装 sink** 的副本／构造态下**必须拿不到栈**（证 AC#2 断的是"装上去那一跳"，不是"Go 天生会打栈"）。反控要带还原自证（仓外副本，或 `git cat-file blob HEAD:<path> >` ＋ md5 相同）。
- [ ] **AC#4 早期崩溃那一档**：日志面还没建好时炸一发，本票定的兜底行为要**可观察**（具名读数，⛔ 不许是"应该有"）。
- [ ] **AC#5 门禁四数**：`GOFLAGS= go build ./...`、`GOFLAGS= go vet ./...`、`sh scripts/d22scan.sh`、`gofumpt -l <自己动过的目录>`，逐条真实 rc；staticcheck 若本机版与 CI 钉版不同 ⇒ 标〔未复认〕、⛔ 不拿假绿。
- [ ] **AC#6 越界检查**：`git diff --name-only <起手锚>..HEAD` 里**不许出现** `internal/observe/goroutine.go` 那枚六协程名册／基线（`ResidentNames`、`ResidentOverBaseline`、`ClassifyGoroutine`）的任何改动——sink 是**注册**动作，不是名册动作。

## 禁区

- ⛔ 一字不动：`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；⛔ 不许为变绿放宽任何断言。
- ⛔ 不动 `internal/observe` 的**六枚常驻协程名册与基线**（D38(e) 那一族；要动先停手上报由我改派）。
- ⛔ 栈里可能带路径／参数 ⇒ 落盘格式**不许把凭据值写进去**；本票的样例行一律只出现变量名／字段名。
- ⛔ 顺手做别的：不接语音 sink（票 247）、不动面板线程（票 33 / `33-r5`）、不动停机十步（票 228）。
- `frontend/**`／`design/**` 零读零写零转述。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件；临时件只建不删。

## 排程

按在 **`33-r5` → 票 248 落地腿 → 票 228 AC#2／AC#11** 之后（它要动 `cmd/wisp` 的装配根，与那三发同包 ⇒ 串行）。⚠ 起腿前必跑两把尺：`git status --porcelain -- cmd internal`（有没有别的写腿）＋ `ls .scratch/wisp/probes/249/`（代号没被用过才写进题面）。
