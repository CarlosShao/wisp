# 249 — 常驻进程的 panic 记录面**没装**：协程观察者今天兜得住炸、但盘上不会留下任何栈

Status: OPEN（编排者 10-01 **13:13** 立，料全部出自探针腿 `33-p1` 与本编排者 13:1x 的现跑尺；台账 `A502`）→ **WITHDRAWN（10-01 14:2x 编排者自撤；出处＝台账 `A503`；⛔ 任何腿不许据此开工）** — 撤回原因：**上面标题与「现量」表第 2 行的那句前提是错的**，我立票时没去读决定它的那三行码；原文一字不抹，判定与复认尺全在文末「撤回理由」一节。撤回前它未派出、零枚提交、零格被勾。
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

---

## 撤回理由（10-01 **14:2x**，我自己复跑尺抓到我自己写错的前提；⛔ 原话不抹，就地留这一节）

**我错在哪一句**：标题与「现量」表第 2 行写的是「`SetPanicSink` 生产零调用者 ⇒ **真常驻炸时盘上不会留下任何栈**」。这句话的**前半是真的、后半是假的**，而我没去读决定它的那两行码就立了票。

**现读复认（三条尺，同发取，`14:2x`）**：
1. `internal/observe/goroutine.go:236`（`NewRegistry` 体内）逐字 `sink: defaultPanicSink` ⇒ **默认就装着**；`:244-245` 更写明 `SetPanicSink(nil)` 会把默认装回来 ⇒ `:314` 那句 `if sink != nil` **在生产里永远不成立**（不可能走空）。
2. `internal/observe/goroutine.go:167-176` 的 `defaultPanicSink` 逐字在 `slog.Error("goroutine panic recovered", …, "stack", ev.Stack, …)` ⇒ **栈进 slog 记录**。
3. `cmd/wisp/logsink.go:145-162` 的 `installLogSink` 把进程级 slog 默认设成 **tee**：`primary`＝**脱敏（redact）之后**的 JSONL 文件管线、`mirror`＝stderr。⇒ 只要数据根解析出来了，**那一枚 `debug.Stack()` 是真的落到数据根内的日志文件里**。

**活样本（不是推的）**：我自己刚在 HEAD 整包复跑 `cmd/wisp` 时抓到一条真 panic，逐字含栈：`level=ERROR msg="goroutine panic recovered" goroutine=panel-sta owner="panel host (ticket 33)" root=panel-host error_class=internal recovered="runtime error: invalid memory address or nil pointer dereference" stack="goroutine 1164 [running, locked to thread]: runtime/debug.Stack() … pkg/edge.(*Chromium).Init(…) chromium.go:131 … cmd/wisp.(*PanelManager).bringUp(…) panel_host_windows.go:226"`（原文在 `.scratch/wisp/probes/orchestrator/33r5-head-roster.txt`）。这枚 panic 本身是真缺陷、已转给 **`33-r6`**（那是另一格，与本票前提无关，⛔ 不许混着读）。

**还剩下的那一小块价值有多低（具名，免得被读成"完全没必要"）**：出货没装的只有「把这条记录改写成 ticket 08 那套 `PanicEvent` JSONL schema」——同一条信息、同一个文件、换个字段形状。⛔ 这不值一枚接线票，更⛔ 不值得 owner 一句话。

**我的账（三条，逐条可核）**：
- ① 我把一枚腿的**推断**（`33-p1` ⑤-13：它只看到 `:314` 的 `if sink != nil`、没看到 `:236` 的默认值）当成了读数，并且**我自己复跑的那把尺只数了调用者枚数**（`grep SetPanicSink` 排 test＝只命中定义）——那把尺证的是"没人替换默认"，⛔ 它证不了"没有默认"。我当时把这两件事写成了一句，越界了。
- ② 这正是我记忆里第 74／83 条那一味**反过来发生**的一次：不是我采纳别人的"做不到"，而是**我自己写了一句"拿不到"**，而它比建议更容易被无条件服从。⇒ 规矩补一句：**"某样东西拿不到／盘上没有"这类否定句，落笔前必须去读"它由谁产生、谁投递、谁落地"那三行，读不全就只许写〔未证〕并禁止任何腿据此开工。**
- ③ 立票快、撤票慢，这一来一回白占了台账一格与排程一行（`A502` 第 P4 条、任务清单）。⇒ 以后凡"新缺一枚接线票"，我先跑的那把尺应该是**能力问句**（"今天炸一次，栈在不在盘上？"＝跑一发看日志），而不是**调用者枚数**（那只说明谁替换了默认）。

⛔ 本票**不再重开**，除非有人证明 `installLogSink` 在真常驻路径上没跑到（那是另一枚缺陷，归 `33-r6`／票 228 那一族，不该由这张纸背着）。撤销口令在这里语义＝"维持撤回"；若要翻案，得先交**一发真跑读数**证明盘上取不到栈。
