# 33-a3 verdict — `*Locked` 命名与锁纪律相反那笔旧账的事实量齐（只读普查，⛔ 不选形、⛔ 不改产码）

- 腿：`33-a3`（只读普查）· 分支 `dev` · 起手锚 `e6c3be17`（commit `d1a6ffef` 落本腿第 1 笔）
- **读数时刻：快照 A `11:23:32 +08`／快照 B `11:27:59 +08`，两次逐枚同号＝未漂**；收尾 HEAD `8b32060b5570cc5c17990475a83030fc2518c29e`（`255-r1: panel width reaches a live window (ticket 255 AC#4, form A)`）
- ⇒ 本腿读数期间在飞的 `255-r1` **已在 11:27 前落进 HEAD**，`git status --short -- cmd internal` 对 `cmd/wisp` 现量**零脏面**
  ⇒ **本件所有 `cmd/wisp/panel_host_windows.go` 行号＝HEAD 现量，可直接引。** ⚠ 此刻在飞的是 `35-r7`（`internal/panel/composer_dispatch.go`＋未跟踪的 `inbound_raw_leak_35r7_test.go`），与本腿名册无交集（本腿只读 `internal/panel/pump.go`，该文件未脏）。
- ⚠ **行号按本腿现量**，`33-v4` 给的号已全部漂移（载体文件被在飞腿 `255-r1` 改了 97 行）：
  `firstRoundTripLocked` `:751` → **`:840`**、`setPriorFocusLocked` `:536` → **`:625`**、
  `serveNotBuiltNoticeLocked` `:460` → **`:460`（未漂）**。凡引用行号请连本段一起引，⛔ 别单引号。

---

## 0. 工具面自报（编排者要求）

```
$ go list ./...            # 只读，35 枚包
rc=0
$ go list ./... | grep -c scratch
0        rc=1（命中 0＝.scratch 树不在编译面内）
$ go env GOVERSION GOMODCACHE GOFLAGS
go1.27.1 / D:\work\base\gopath\pkg\mod / (GOFLAGS 空)
rc=0
$ git log -1 --format=%H ; git status --short -- cmd internal
rc=0
```
⛔ 未跑 `go build` / `go vet` / `go test` / 任何测试二进制（`255-r1` 独占同机 Go 编译面）。
本腿全部读数来自 `grep` / `awk` / `git` / `go list`，**零执行**。

---

## Q1｜名册：全仓带 `Locked` 后缀的方法 23 枚（产码）＋22 枚（`.scratch` 里不可编译的拷贝）

分母尺写法（⛔ 裸 grep 符号名）：

```
声明尺：grep -rn --include='*.go' -E '\bfunc \([^)]*\) [A-Za-z_][A-Za-z0-9_]*Locked\(' .   → rc=0，45 命中
        其中 .scratch 22 枚（grep -c 'scratch' = 22，rc=0）／产码 23 枚（grep -vc = 23，rc=0）
        拷贝名册＝6(probes-242-r1/backup)＋2(197/r3b/pre)＋6(242/r1/backup)＋6(259/r2/before)＋1(33/r8b/logs/pristine)＋1(62/v1/m1)
自由函数尺：grep -rn --include='*.go' -E '\bfunc [A-Za-z_][A-Za-z0-9_]*Locked\(' .          → rc=1，零枚
字符串字面尺：grep -rn --include='*.go' -E '"[A-Za-z_]*Locked[A-Za-z_]*"' cmd internal tools → rc=1，零枚
自取锁尺（逐枚）：对上表 23 枚声明行号逐枚 `awk -v s=<decl行> 'NR>=s{print} NR>s && /^}/{exit}'` 抽函数体，
        再 `grep -qE '\.mu\.(Lock|RLock)\(|\.lk\.'` → 输出 23 行中 **SELF-LOCKS 只有 2 行**（`panel_host_windows.go:460`／`:751`（当时读数，现 `:840`）），
        其余 **21 行全部 no-lock** ⇒ **第三枚自取锁的 `*Locked` 零枚**（该支 rc=1）
调用尺（每一枚）：<接收者变量>.<方法名>(   —— 带左括号，所以构造点／同名字段／注释不计
```
`<类型全名>＋<调用形状>` 逐枚（**非测试调用者枚数／各自此刻持不持锁**）。三栏＝①注释写没写"caller must hold"②函数体第一句是不是自己取锁③非测试调用者。

| # | `<类型全名>.<方法>` | `<file>:<line>`（现量） | ①注释 | ②自取锁 | ③非测试调用者（形状 `x.M(`） |
|---|---|---|---|---|---|
| 1 | `*PanelManager.serveNotBuiltNoticeLocked` | `cmd/wisp/panel_host_windows.go:460` | 无（:456-459 只讲用途） | **是**（`:461 m.mu.Lock()`） | **1**：`:451`（`coldStartPageHandover` 内）**不持锁** |
| 2 | `*PanelManager.firstRoundTripLocked` | `cmd/wisp/panel_host_windows.go:840` | 无（:824-839 讲射程，未提锁） | **是**（`:841`… 现量函数体第一句 `m.mu.Lock()`） | **1**：`:448`（同上）**不持锁**；＋1 测试 `panel_pageover_33r10_windows_test.go:222` 不持锁 |
| 3 | `*PanelManager.setPriorFocusLocked` | `cmd/wisp/panel_host_windows.go:625` | **有**（`:624` "The caller must hold m.mu."） | 否 | **1**：`:525`，落在 `Show` 的 `:522-526` 持锁段内＝**持锁** ✓ |
| 4 | `*approval.Queue.indexLocked` | `internal/agent/approval/queue.go:210` | 有（:209 "Caller holds q.mu."） | 否 | 1：`:180`（enqueue）持锁 ✓ |
| 5 | `*Queue.unindexLocked` | `…queue.go:223` | 有（:222） | 否 | 1：`:302`（在 `dropLocked` 体内＝**传递持锁**）✓ |
| 6 | `*Queue.lookupForAllowLocked` | `…queue.go:246` | 有（:245） | 否 | 2：`:385`（allowScoped，:384 Lock）✓／`:457`（sessionSubject，:455-456 Lock+defer）✓ |
| 7 | `*Queue.lookupForRefusalLocked` | `…queue.go:265` | 有（:264） | 否 | 1：`:493`（reject，:492-498 段内）✓ |
| 8 | `*Queue.dropLocked` | `…queue.go:294` | 有（:293） | 否 | 2：`:320`（deliver :314-315）✓／`:353`（grantNonce :351-354 段内）✓；＋1 测试 `ticket259_denial_rulers_test.go:216`（:215 显式 Lock）✓ |
| 9 | `*Queue.viewLocked` | `…queue.go:570` | **无（整条没注释）** | 否 | 2：`:567`（view :561-562）✓／`:588`（head :583-584）✓ |
| 10 | `*HalfDuplexGate.openInnerLocked` | `internal/audio/gate.go:247` | **未写持锁**（:244-246 讲失败上报） | 否 | 3：`:109`（Start）／`:147`（SetSpeaking）／`:183`（SetMuted）全部在 `g.mu.Lock()` 段内 ✓ |
| 11 | `*HalfDuplexGate.closeInnerLocked` | `internal/audio/gate.go:255` | **无注释** | 否 | 2：`:153`／`:189` 持锁 ✓ |
| 12 | `*Ball.applyStateLocked` | `internal/ball/ball_windows.go:342` | 有，但**否掉锁**：:340-341 "(The name is historical: the state fields are STA-thread-owned, not mutex-shared.)" | 否 | 2：`:263`（`createOnSTA`，STA 线程，**不持 `b.mu`**）／`:312`（`b.sta.PostTask(func(){…})` 闭包内，**不持锁**） |
| 13 | `*config.Manager.planLocked` | `internal/config/manager.go:314` | **未写 mutex 持锁**；:309 那句 "handles one locked section" 里的 "locked" 指 **D36 的 locked section（权限锁）**，不是 mutex | 否 | 4：`:248/:253/:258/:263`（都在 `plan()` 内；`plan()` 唯一调用点 `manager.go:169` 落在 `:155-173` 的 `m.mu` 段内＝**传递持锁**）✓ |
| 14 | `*rollingWriter.ensureFileLocked` | `internal/observe/logging.go:225` | 未写持锁（:223-224 讲滚动规则） | 否 | 2：`:215`（Write :210-211）✓／**`:199`（`newRollingWriterClock` 构造体内＝不持锁**，对象尚未发布） |
| 15 | `*rollingWriter.closeFileLocked` | `internal/observe/logging.go:254` | **无注释** | 否 | 2：`:230`（ensureFileLocked 体内＝传递）✓／`:285`（Close :282-283）✓ |
| 16 | `*rollingWriter.sweepLocked` | `internal/observe/logging.go:317` | 未写持锁（:314-316 讲保留窗口） | 否 | 2：**`:202`（构造体＝不持锁**）／`:304`（flushLoop :303-305 段内）✓；＋1 测试 `logging_test.go:224`（**不持锁**） |
| 17 | `*StreamLog.enforceBoundLocked` | `internal/panel/pump.go:577` | 有（:572-576 "Caller holds s.mu."） | 否 | 2：`:528`（Append :522-523）✓／`:551`（Close :545-546）✓ |
| 18 | `*StreamLog.clampLocked` | `internal/panel/pump.go:600` | 有（:599） | 否 | 2：`:534`（Append 段内）✓／`:583`（enforceBoundLocked 体内＝传递）✓ |
| 19 | `*ShutdownHookSet.setSlotLocked` | `internal/proc/shutdown_hooks.go:110` | 未写持锁（:106-109 讲穷尽性） | 否 | 1：`:102`（Register :87-88 Lock+defer）✓ |
| 20 | `*Machine.rearmLocked` | `internal/statemachine/machine.go:178` | 有（:177 "Called with m.mu held"） | 否 | 2：`:158`（Dispatch :113 起段内）✓／**`:73`（`New` 构造体内＝不持锁**，尚未发布） |
| 21 | `*Machine.stopTimerLocked` | `internal/statemachine/machine.go:205` | **无注释** | 否 | 2：`:173`（Close :170-171）✓／`:179`（rearmLocked 体内＝传递）✓ |
| 22 | `*Registry.addLocked` | `internal/tools/registry.go:153` | 有（:152 "callers hold the write lock"） | 否 | 2：`:110`（Register :99-100）✓／`:194`（RegisterProvider :184-185）✓ |
| 23 | `*Registry.listLocked` | `internal/tools/registry.go:224` | **无注释** | 否 | 1：`:221`（List :219-220 `RLock`+defer）✓ |

三栏汇总（⛔ 不是"两枚的事"，是全族的账）：

- ① **明写 caller-holds 的＝10 枚**（#3,4,5,6,7,8,17,18,20,22）；① **完全没有注释的＝5 枚**（#9,11,15,21,23）；① **有注释但没写锁契约的＝8 枚**（#1,2,10,12,13,14,16,19）。`10+5+8=23` ✓
- ② **函数体自取锁的＝2 枚**（#1、#2）。尺＝逐枚抽函数体到行首 `}` 再 grep `.mu.(Lock|RLock)`：
  `…/verdict §Q1尺2`，命中只有这两枚，**第三枚零枚**（rc 见下）。
- ③ 非测试调用者合计 **41 处**；测试调用者 **3 处**（`dropLocked`×1、`sweepLocked`×1、`firstRoundTripLocked`×1）。
  41 处里**此刻不持锁的＝7 处**：#2×1、#1×1（这两枚是"不许持锁"的正解）、#12×2（设计上无线程锁）、
  #14/#16/#20 各 1 处（构造期，对象尚未发布）。⇒ **另有一族"名字没写、却是 caller-must-hold"的没出现在②里，见 Q6。**

`.scratch` 那 22 枚是同一批产码文件的**前像/备份拷贝**（`259/r2/before/queue.go`、`242/r1/backup/queue.go`、
`probes-242-r1/backup/queue.go`、`197/r3b/pre/pump.go`、`33/r8b/logs/pristine-ball_windows.go`、`62/v1/m1/ball_windows.go`），
`go list ./...` 零命中＝**目录以 `.` 开头，go 工具根本不载入**，⛔ 不进编译面、⛔ 不该进任何"全仓枚数"的分母。

---

## Q2｜今天到底有没有一条真死锁路 —— **没有。这是一枚陷阱，不是现行 bug**

逐枚判（只两枚自取锁，所以只需判这两枚的调用点能否在持锁状态下到达）：

**#2 `*PanelManager.firstRoundTripLocked`（`panel_host_windows.go:840`，第一句 `m.mu.Lock()`）**

- 非测试调用点**只有 1 处**：`:448`，位于 `coldStartPageHandover`（`:447-454`）。
- 该函数体**整段零 `m.mu.` 语句**（尺：`awk 'NR>=447 && NR<=454' | grep -c 'm\.mu\.'` ＝ **0**，rc=1）。
- 它的唯一调用点 `:421`（`bringUp`）落在 `m.mu.Unlock()`（`:417`）与 `m.mu.Lock()`（`:423`）**之间**＝不持锁。
- 测试调用点 `panel_pageover_33r10_windows_test.go:222`：接收者是 `NewPanelManager(...)` 现造的实例，
  前一行 `:221 attachSink33r10(...)`，**无任何 `Lock`** ⇒ 不持锁。

**#1 `*PanelManager.serveNotBuiltNoticeLocked`（`:460`，`:461` 自取锁）**

- 非测试调用点**只有 1 处**：`:451`，同一个不持锁的 `coldStartPageHandover` 内；测试调用者**零枚**。

⇒ **判定：今天不存在一条"持 `m.mu` 又调自取锁的 `*Locked`"的路径。它是命名与锁纪律相反留下的陷阱，⛔ 不是现行死锁。**
（⚠ 这一格按编排者要求的分寸写：`33-v4` 的原话"谁照名字在持锁状态下调这两枚，就是自死锁"＝**条件句**，本腿证实条件句的前件今天在任何调用点都不成立。）

**但有一条相邻的真锁序边，具名在此（⛔ 不是自死锁，是 AB-BA 的原料）**：

- `cmd/wisp/panel_resident_windows.go:323-324`：`rp.mu.Lock()` 之内调 `rp.mgr.currentWindow()`，
  而 `currentWindow`（`panel_host_windows.go:303-304`）**取 `m.mu`** ⇒ 本进程存在一条 **`rp.mu → m.mu`** 的持锁边。
- 反向边今天**零枚**。尺（⛔ 不含测试文件，逐枚带 rc）：
  `grep -n 'm\.mu\.' cmd/wisp/panel_host_windows.go` ＝ **42 行**（rc=0；其中 1 行是注释 `:624` ⇒ 可执行锁句 41）；
  `grep -rn --include='*.go' '\bm\.mu\.' cmd/wisp \| grep -v 'panel_host_windows.go' \| grep -vc '_test.go'` ＝ **0**（rc=1）
  ⇒ `PanelManager.mu`（声明在 `panel_host_windows.go:147`）**只在同一枚文件内被取放**，该文件的产码零处引用 `residentPanel`
  （`grep -c 'residentPanel' cmd/wisp/panel_host_windows.go` ＝ 0，rc=1）⇒ **不存在 `m.mu → rp.mu` 的反向持锁边。**
- ⚠ 这一格的意义是给 Q4 的 ⓑ 用的：**把 `m.mu` 的持有段从"三句"扩成"跨一次最多 5 秒的消息泵"**，
  就是把一条今天单向的锁序边变成一条有 5 秒窗口去撞的边。

---

## Q3｜这台机/这套仓有没有仪器看得见它 —— **零覆盖（具名）**

| 候选仪器 | 现量 | 结论 |
|---|---|---|
| `tools/d22scan` 的 ban 名册 | `main.go:10-40` 列 **8 项**：goroutine / pathresolver-bypass / plaintext-key / wallclock-timeout / mirror-hash / panel-approval / internal-artifact-tool / emoji。`grep -iE 'lock\|mutex\|deadlock' tools/d22scan/main.go` 命中的只有 `wallclock`/`block comment` 那几行 | ⛔ **锁命名／持锁调用／重入不在射程** |
| `-race` 门 | `grep -rn -i 'race' .github/workflows/ci.yml scripts/*.sh scripts/*.ps1`（去掉 `trace`/`graceful`/`brace` 类无关词）＝**零命中**；`-race` 只出现在 `docs/evidence/**` 的历史验收记录里；全仓唯一带 race tag 的文件＝`internal/risk/pathresolver_budget_norace_test.go:1`（`//go:build !race`） | ⛔ **本仓 CI/脚本今天没有 `-race` 这一道门** |
| 就算有 `-race`，它能不能看见死锁 | ⛔ **看不见**。`-race` 报的是**数据竞争**（两发未同步访问，其中至少一发是写）；死锁没有非法访问可言，只有**挂住**。死锁在测试二进制里的表现分两形：① 全进程所有 goroutine 都永久阻塞 ⇒ Go 运行时 `fatal error: all goroutines are asleep - deadlock!`（这形会红，且**不需要** `-race`）；② 只要还有任何一枚活着的计时器/别的线程（本仓有 `time.AfterFunc`（`machine.go:187`）、STA 线程、`observe` 的 flusher、WebView2 自己的线程），运行时的判死探测器就**不响**，于是表现为**挂到 `go test` 超时**才红 | 本文件恰是第②形：`firstRoundTripLocked` 里就有一枚 `time.AfterFunc` 之外的 5 秒 `deadline` 循环＋`pnlPumpOnce()`，`Machine` 还常驻一枚 timer ⇒ 真踩的话大概率是**超时红**，不是"all goroutines asleep" |
| 超时那把尺的口径 | `grep -rn -- '-timeout' .github/workflows/ci.yml scripts/*.sh scripts/*.ps1` ＝**零命中**；`tools/d22scan/runtests.sh:75` 是 `go test -v -count=1 "$@"`，不带 `-timeout` | ⇒ 真死锁会**烧满 `go test` 默认 10 分钟**才红，⚠ 且红句是超时栈、不含"命名"二字 |
| 有没有任何测试按内容锚盯 `*Locked` 的锁形状 | 编译面（`cmd`＋`internal`＋`tools`）里**字符串字面量**出现任何 `*Locked` 名的：`grep -rn --include='*.go' -E '"[A-Za-z_]*Locked[A-Za-z_]*"' cmd internal tools` ＝ **0 行，rc=1**；AST 锚只盯 `windowOptions`／`bringUp`／`NewWithOptions`／`WindowOptions`（`panel_geometry_255_test.go:223/227/293/294`） | ⛔ **没有任何仪器看得见这件事**：既无静态尺、也无运行时尺 |

⇒ **Q3 结论＝零覆盖，具名。** 这件事今天"绿"的唯一原因是**没人碰它**，不是因为有人守着。

---

## Q4｜三形各自的代价（⛔ 本腿不选形）

先给一条**贯穿三形的共同事实**：这两枚的**真实行为**是"短持锁读 `m.w` → 解锁 → 跨锁做 `SetHtml`/泵消息"，
与在飞腿 `255-r1` 刚加进同一文件的 **`requestGeometryOnReshow`（`:598`，无 `Locked` 后缀、`:605` 自取锁、调用点 `Show:519` 在锁外）完全同形**。
⇒ 今天这个文件里有 **3 枚"自取锁、必须锁外调"的 helper，其中 2 枚的名字反着说**。

### ⓐ 改名（去掉两枚的 `Locked` 后缀）

- **动几枚文件**：产码 **1**（`cmd/wisp/panel_host_windows.go`：声明 2 枚 `:460`/`:840` ＋调用 2 处 `:448`/`:451` ＋注释 5 处
  `:430`/`:797`/`:824` 与 `:456`＝现量 `firstRoundTripLocked` 5 命中／`serveNotBuiltNoticeLocked` 3 命中）；
  测试 **1**（`cmd/wisp/panel_pageover_33r10_windows_test.go`：调用 `:222` ＋注释 `:23`/`:55`/`:216`/`:238`＝5 命中）。⇒ **合计 2 枚文件**（`git mv` 之外零新文件，⛔ 不动分母名册）。
- **会不会打到既有钉/既有红**：
  - ⚠ **顶回编排者一句**（见文末"顶回"§第 1 条）：`cmd/wisp/config_readers_255.go` 与 `panel_geometry_255_test.go` 这两把内容锚尺
    **锚不到这两个名字**——`evidenceCite` 只认 `path.go:LINE [token]`（`config_readers_255.go:228`），
    `panel_geometry_255_test.go` 的 AST 只查 `windowOptions`/`bringUp`/`NewWithOptions`/`WindowOptions`。⇒ **改名的红面是编译级的 `:222` 那一发，不是内容锚。**
  - 行号类 cite：`config_readers_255.go:181` 那枚受检 cite 是 `panel_host_windows.go:262 [Width:  uint(width)]`，
    本腿 `11:23:49` 复量 `:262` ＝ `Width:  uint(width),` **对得上**；改名**不增删行** ⇒ 该 cite 与 `:392` 那句散文都不动。⛔ 安全。
  - 票面会变旧：`.scratch/wisp/issues/33-panel-host-c27.md` 里 `firstRoundTripLocked` **4 处**（`:39`/`:321`/`:348`/`:358`）
    是按名字写的射程条款；台账 `docs/reports/pending-and-issues.md` 另有 `10204`/`10292`/`13737` 三处点名。
    ⇒ 改名后这些**人类文本**指不到符号（不改判据、不改代码语义，但**要落账**；票面/台账本腿没动，⛔ 也不该由腿动）。
  - ⚠ **在飞冲突面**：`255-r1` 此刻正在改**同一枚文件**（`git status`：`M cmd/wisp/panel_host_windows.go`＋`config_readers_255.go`＋`panel_resident_windows.go`＋`panel_geometry_255_test.go`，
    未跟踪 `?? cmd/wisp/panel_reshow_255r1_windows_test.go`）。⇒ ⓐ 现在落地＝**与在飞腿同文件同行区改**，
    这正是本项目"改名中间态把在飞腿吓自停"那一枚先例的形状 ⇒ **ⓐ 的时间窗只能在 255-r1 结线之后开。**
- **新代码会不会踩进陷阱**：ⓐ 之后 `requestGeometryOnReshow` 与两枚改名后的 helper **三枚同形同名规**（都无后缀、都自取锁），
  ⇒ ⓐ 是唯一一形**顺手把 255-r1 刚造的那枚也一并归位**的（它今天名字对、行为对，只是孤例）。

### ⓑ 保住后缀、把取锁挪到调用点（真做成 caller-must-hold）

- **动几枚文件**：产码 **1**（`panel_host_windows.go`：两枚声明去掉自带 `Lock/Unlock` ＋ `coldStartPageHandover` 加持锁段）＋测试 **1**（`:222` 那发也得包锁）。
  ⇒ 纸面 2 枚，**但它是三形里唯一"动锁形状"的**，`33-v4` 已按这一条把自己排除（"改它＝动锁形状，超出验收腿射程"）。
- **代价不是机械搬运，本腿量到三格硬的**：
  1. `firstRoundTripLocked` 的主体是**最多 5 秒的消息泵**（`:868 deadline := t0.Add(5 * time.Second)` 起的 `for` ＋ `pnlPumpOnce()`）。
     真做成 caller-must-hold ⇒ **`m.mu` 被跨在一次 5 秒泵上**；同文件里 `IsShown`/`IsCreated`/`windowHandle`/`currentWindow`/`LastColdMs` 全都排队等 `m.mu`，
     而 `panel_resident_windows.go:323-324` 那条 `rp.mu → m.mu` 边会把这 5 秒**传染给常驻面板线程的 `rp.mu`**。
     ⇒ ⓑ 要么接受"锁持有跨 5 秒"，要么**先把"读 `m.w`"和"泵"拆成两枚**——那就不再是"挪锁"，是重做函数边界。
  2. `coldStartPageHandover` 里 `:450 m.serveEntry()` **自己也取 `m.mu`**（`:480`）。⇒ 持锁段必须**只**裹住 `:448` 与 `:451` 两处、
     在 `:450` 之前放开；写成"整个 handover 持锁"＝**当场造出真死锁**（`serveEntry` 重入 `m.mu`）。这一格要写进任何派单。
  3. 它**增删行** ⇒ 一切行号 cite 的漂移风险。本文件唯一受检 cite 在 `:262`（在改动区**之上**），⛔ 不被 ⓑ 打中；
     但散文 cite `reached from the create at :392`（同一条 roster 句里）会随插入漂，且 `panel_inbound_guards_35r3_test.go:43`
     的注释锚 `panel_host_windows.go:694-696`、`panel_transport_35r2_test.go:2041/:2076`、`panel_transport_live_35v2_windows_test.go:182`
     四处**散文行号**今天全指 `:694-696`/`:680`/`:638` 一带 ⇒ ⓑ 在 `:448-:840` 区间插行会把它们全部写旧（无仪器、但后续程会照它们读码）。
- **新代码会不会踩进陷阱**：ⓑ **会新造一处**——`requestGeometryOnReshow`（`:598`）今天锁外调（`Show:519`）。
  如果 ⓑ 的"形制统一"顺手把它也改成 caller-must-hold，那 `Show` 就变成
  `:507-508 短持锁读 created` → `:519 取几何+SetSize` 的两段持锁，而 `windowOptions()`（`:237`）读 `m.geometry` 不取锁＝侥幸安全；
  ⚠ 一旦有人把 `m.geometry` 从"构造后不改的字段"（`:183` 注释）改成可变字段，这条侥幸就没了。⇒ ⓑ 与 255-r1 的**语义耦合**必须显式写死。

### ⓒ 只加一枚会响的钉（⛔ 不动产码一行）

- **动几枚文件**：新增 **1** 枚 `cmd/wisp/*_test.go` ⇒ 三形里唯一**加文件**的。⚠ 加文件的代价本仓记过账：
  它会挪 `d22scan` **ban #8 的分母**（`internal/`/`cmd/` 面的 emoji 命中数按现量入账，票 143/141 都为此改过台账数字），
  ⛔ 且⚠ 按 `AGENTS.md §1.5`／编排者本轮纪律"⛔ 新建 `.sh`/`.ps1`/`.txt`/`.out`"，新文件**只能是 `.go` 测试**、落点必须问一句归谁。
- **可复用的既有形制（本腿找到一枚真的先例，具名）**：`internal/config/manager_223_test.go:39-51` 的
  `completedWithin(budget, fn)`＝把可疑调用放到独立 goroutine、`select` 上 `time.After(budget)`，
  `:113-120` 用它的返回值把"钩子仍在锁内"判成**红**而不是挂住。文件 `:20-25` 还**具名解释**了为什么这里允许裸 `go`（d22scan bans #1-5 不扫 `_test.go`）。
  ⇒ ⓒ 若要钉，**本仓已有现成尺形**，⛔ 不必新造判据语言；⚠ 但它是**墙钟预算**（3 s / 10 s），
  而 `AGENTS.md §1.2` 禁的是"产码里用墙钟时间差实现超时"——测试面今天不受该 ban（bans #1-5 不扫 `_test.go`），
  这条豁免在 `manager_223_test.go:20-25` 是**写死并归过口**的，ⓒ 复用时要照抄那段自陈。
- **静态那形（AST 钉）**：本仓已有能力（`panel_geometry_255_test.go` 用 `go/parser` + `ast.Inspect` 在同一文件上取证），
  ⇒ 钉"声明名以 `Locked` 结尾 ⇒ 体内不得出现 `mu.Lock`"技术上 1 枚文件就够。
  ⚠ 自指风险：`config_readers_255.go:228` 那把 cite 尺**会命中注释里写的 token**（`evidenceCite` 不区分注释），
  同类病本项目记过（"注释里写会被自己扫的现在时计数"）⇒ 静态钉的注释⛔不得出现 `xxx.go:NNN [token]` 形状。
- **既有钉会不会被 ⓒ 点红**：会，且只有一处已知——`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`
  对 `citedRowsFloor = 8`（`config_receipt_255_test.go:220`）设了下限，ⓒ 若顺手往 roster 加行只增不减 ⇒ 安全；
  但 ⓒ **新增测试名**不需要进 `scripts/portable-tests.sh` 的 not-in-scope 名册（那张表只登记"声明不在射程"的用例，`:20-35`），
  ⚠ 唯一要求是 ⛔ 不得 `t.Skip`（`runtests.sh` 把"任何 top-level SKIP"判红）。
- **新代码会不会踩进陷阱**：ⓒ 对 255-r1 的 `requestGeometryOnReshow` **零改动**，且钉会把
  "名字带 Locked 却自取锁"从**读码约定**变成**红了会响的事** ⇒ 唯一三形里"今天不碰在飞腿"的一形。

---

## Q5｜查重 —— **票池里没有这格框，⛔ 不新立一票；具名归口＝票 33**

尺（按"这句话在说什么"扫，⛔ 只搜符号名；长行**整行读**，未 `cut`）：

```
$ grep -rn -iE '死锁|deadlock|reentr|不可重入|自取锁|自己取锁|自锁|must hold|持锁|重入' .scratch/wisp/issues/ | wc -l
26     rc=0            （14 枚文件命中）
$ grep -rn -iE '…同上…' docs/reports/pending-and-issues.md | wc -l
21     rc=0
$ grep -rln -E 'Locked|自取锁|不可重入|重入|self-deadlock|self-locking|死锁' .scratch/wisp/issues/ | wc -l
19     rc=0
```

逐条判定（命中≠同账，按机制链分）：

| 命中处 | 它在说什么 | 与本账同一条机制链吗 |
|---|---|---|
| 票 33 `AC#13`（`:39`）／`:321`／`:348`／`:358` | 冷启动**最终显示哪份文档**，`firstRoundTripLocked` 只是被点名的那一步 | ⛔ 不同事（射程/次序），但**同一枚文件、同一枚函数** ⇒ 天然归口 |
| 票 222 `:24` | spawn 占许可等孩子的**互等**；它自己写"⛔ '永挂'这一支我不下结论……真形状是周期性松绑＋父侧报错，不是死锁" | ⛔ 不同机制（许可饥饿，且已被作者否掉"死锁"定性） |
| 票 195 `:60` | "不许在入口里等一张卡＝**死锁不是死红**"，是设计约束；它的凭据是 `manager_223_test.go:55` | ⛔ 不同机制，但**同一族解法**（锁外做会阻塞的那一跳） |
| 票 47 `:25-26` | 任务调度器的"两任务互等 ⇒ 超时＋显式报错"deadlock guard | ⛔ 不同机制（任务图，不是同一枚 mutex） |
| 票 130 `:64` | `slog.Default().Handler()` 当 mirror **会死锁**，已修＋留守卫用例 | ⛔ 不同机制；但是**"自取锁族"结过一次账的先例** |
| 票 87 `:29`／`:124` | "不要把这张票读成修一个死锁"＋当年设计名 `resolveLocked(name, strict)`（**今天不存在**：全仓 `func (…)` 名册里零枚 `resolveLocked`） | ⛔ 已定性的旧账；`resolveLocked` 是没落地的旧设计名 |
| 票 266（标题含"门**自锁**钉"） | `scripts/*.ps1` 的门钉可被删掉而零外牙——"self-locking" 指的是**仪器锁自己**，⛔ 与 mutex 无关 | ⛔ **术语撞名，必须具名排除**（否则下一个读台账的人会以为已经有一票管这件事） |
| `pending-and-issues.md:9216` | `manager.go` 头注释"Callbacks run OUTSIDE the lock"**曾对 `ConfirmLocked` 是假话** | ⛔ 已修（票 223），但是**同族第二枚具名旧账**（见 Q6 第 1 条） |
| `pending-and-issues.md:13737`（**A697**）／`:13809`（**A702**） | 编排者把"锁语义"这一格**具名交给 `33-v4` 裁**，随后 A702 收下 v4 并**把本腿（`33-a3`）写进在飞名册** | ✅ **就是这笔账本身**，已在台账，⛔ 不需要新条目 |
| `.scratch/wisp/probes/33/v4/verdict.md:54-69` | 这笔账的**原始出处**（本任务转述一致，行号现量＝`:54-69`，其中判定在 `:60-65`、判语在 `:69`） | ✅ 出处 |

⇒ **查重结论：票池里没有已经存在的框/票管这件事**（14 枚命中文件逐枚都不是这条机制链）；
**台账已有两笔（A697 交棒／A702 收棒＋派本腿）**，**原始记录在 `33-v4` verdict `:54-69`**。
⇒ **归口建议（⛔ 由编排者落笔，本腿不改票面）：票 33**——载体同一枚文件、且 `AC#13` 框里今天就写着 `firstRoundTripLocked` 的名字；
⛔ 不新立一票（本项目定式：同一机制链只记一次）。⚠ 若编排者判定"命名纪律是全仓的"（Q1 ①栏里 **13 枚没写 caller-holds**、
**5 枚零注释**），那要立的也不是"修这两枚"，而是"这后缀在本仓到底是什么意思"那一枚**口径**——**这格超出本腿射程，我只把料备齐。**

---

## Q6｜别的同形 —— **有，五枚，具名；不是零枚**

**甲形（名字没后缀、却要求持锁——编排者判定的"更危险那一形"）：3 枚，全在 `internal/config`**

| 枚 | `<file>:<line>` | 契约文本 | 今天守不守 |
|---|---|---|---|
| `*reloadPlan.commit` | `internal/config/manager.go:228` | `:227` "Caller holds m.mu." | ✓ 唯一调用点 `:193` 在 `:192-196` 持锁段内 |
| `*Manager.writeAllowedDirs` | `internal/config/allowdirs.go:139` | `:132` "Caller holds m.mu." | ✓ 调用点 `:82`（`:74-75` 持锁）／`:124`（`:122-123` 持锁） |
| `*Manager.mergeWrite` | `internal/config/writeguard.go:111` | `:83` "Caller holds m.mu." | ✓ 调用点 `allowdirs.go:142`（传递）／`permmode.go:79`（`:74-75` 持锁）／`settings.go:303`（`:262-263` 持锁） |

⇒ 名字**一律没提示**，守不守全靠注释＋读码；⛔ 三枚今天**没有一处违规**，也**没有任何仪器盯**（同 Q3 的零覆盖）。

**乙形（`Locked` 后缀、但那个"Locked"根本不是 mutex 意思）：3 处**

1. `config.Manager.ConfirmLocked`（**字段不是方法**，`internal/config/manager.go:51`）——
   `:43` 的 "locked" 指 **D36 的 locked section**，而 `:47` 明写 **"It runs OUTSIDE mu"**；
   它的调用点 `:172` 取快照、`:180` **在 `:173` 解锁之后**才调＝**契约与后缀相反**。
   ⚠ 这笔账**结过一次**：台账 `:9216`＋票 223＋`manager_223_test.go:53-120` 就是"名字说锁内、真要在锁外"的旧案，
   修法正是"改行为、留名字"（票 223 之后名字未动）⇒ **本仓对同一族病有过一次定形判例**，裁定这两枚时该引。
2. `*Ball.applyStateLocked`（`internal/ball/ball_windows.go:342`）——`:340-341` 自陈
   "the name is historical: the state fields are STA-thread-owned, not mutex-shared"：后缀在此**命名的是线程亲和，不是锁**。
3. `*Manager.planLocked`（`internal/config/manager.go:314`）——`:309` "handles one locked section" 里 "locked" 同样指 D36 档；
   它**事实上**要求持 `m.mu`（经 `plan()` 传递），但**注释没说**⇒ 一枚"后缀来源是别的词、锁契约只靠传递链成立"的中间态。

**丙形（`Locked` 后缀、零"caller-holds"注释）：13 枚**——见 Q1 ①栏的 5 枚"完全没有注释"（#9/#11/#15/#21/#23）
＋8 枚"有注释但没写锁契约"。⇒ **"Locked 后缀在本仓到底意味着什么"这件事今天是三套语义并存**（must-hold / self-locking / D36-locked / thread-affinity），
这枚事实是 Q4 三形之外的**第四种可能**（先定口径再动代码），⛔ 由编排者裁，本腿只把枚数摆出来。

**丁形（反例：处理得很干净的一形，具名以免被当成同形）**
`internal/tools/paths.go:236` `rootsContain` 的 `:234-235`："**It takes no lock so the locked callers can compose it (InAllowlist) without re-entering `p.mu`**"
＝显式声明"我不取锁，所以持锁调用方可以安全组合"。⇒ 本仓**有**"用注释把重入面钉死"的写法先例，`ⓐ/ⓒ` 若要补注释，形制照它。

---

## 顶回（编排者要我具名报回的句子；没有则写"零枚"——**本腿有 3 枚**）

1. **⛔ 不成立**：派单里"⚠ 特别注意 `cmd/wisp/config_readers_255.go`、`panel_geometry_255_test.go` 这类按内容锚解析源码的仪器，**改名会让它们红**"。
   现量：`evidenceCite = regexp.MustCompile(`+"`"+`([\w./-]+\.go):(\d+) \[([^\]]+)\]`+"`"+`)（`config_readers_255.go:228`）只认 `file:LINE [token]`，
   全仓字符串字面量里的 `*Locked` 名＝**零命中**（尺见 Q3 表第 4 行，rc=1）；`panel_geometry_255_test.go` 的 AST 锚只查
   `windowOptions`/`bringUp`/`NewWithOptions`/`WindowOptions`。⇒ **改名的真红面是 `panel_pageover_33r10_windows_test.go:222` 那一发的编译依赖，不是内容锚。**
   （⚠ 但派单的**方向**对：那两把尺确实会因**增删行**红——那是 ⓑ 的账，见 Q4 ⓑ 第 3 格。）
2. **⛔ 我的定性词不许重于证据这条对派单同样生效，本腿据现量收窄一句**：派单写"今天到底有没有一条真死锁路"时把答案预分成"有/没有"，
   现量是"**没有**"（条件句前件在 3 个调用点上全部不成立，逐跳凭据见 Q2）。⇒ 任何把这笔账写成"**现行 bug**"的票面措辞，本腿不背书。
3. **⚠ 补一笔派单没问的**：`panel_resident_windows.go:323-324` 那条 **`rp.mu → m.mu`** 持锁边今天单向、安全，
   但它是 ⓑ 的真实代价来源（"锁跨 5 秒泵"会把 5 秒传染给常驻面板线程的 `rp.mu`）。
   ⛔ 派单只把 ⓑ 的代价问成"动几枚文件＋会不会打到既有钉"，这一格不在那两问里，但裁定必须知道 ⇒ 具名补上。

**另：一条我自己要顶自己的**：`11:20-11:23` 我两次量到 roster 的两枚受检 cite **不匹配**
（`panel_host_windows.go:262` 读到 `Title:`、`panel_resident_windows.go:207` 读到 `}`），
`11:23:49` 复量**两枚都对上了** ⇒ 那是**读到在飞腿写入的中间态**，⛔ 不是缺陷、⛔ 不记账。
（同一机制的镜像病：⚠ 在有人正在写的树上取行号锚，每一发都要带时间戳并复量一次才敢写进结论。）

---

## 收尾自证（本腿只写了自己那两件；尺于 `11:30:26 +08` 现跑）

```
$ git show --stat --format='%h %s' d1a6ffef | head
 .scratch/wisp/probes/33/a3/00-anchor.md | 37 ++++++++++++++++++++++++++++++++
 1 file changed, 37 insertions(+)                                        rc=0

$ git show --stat --format='%h %s' 824746eb | head
 .scratch/wisp/probes/33/a3/verdict.md   | 305 ++++++++++++++++++++++++++
 1 file changed, 305 insertions(+)                                       rc=0

$ wc -c .scratch/wisp/probes/33/a3/00-anchor.md .scratch/wisp/probes/33/a3/verdict.md
 1670 .scratch/wisp/probes/33/a3/00-anchor.md
32923 .scratch/wisp/probes/33/a3/verdict.md                              rc=0

$ git status --short -- cmd internal docs .scratch/wisp/issues
?? internal/panel/inbound_raw_leak_35r7_test.go     ← 35-r7 在飞，⛔ 不是本腿的
                                                                    rc=0
```
⇒ 本腿两枚 commit（`d1a6ffef`／`824746eb`）**各只带 `probes/33/a3/` 一枚路径**，
`cmd`＋`internal`＋`docs`＋票池对本腿**零脏面**（唯一那行未跟踪件属 `35-r7`）。
⚠ 本节自身的落笔 commit 记不进本表（同一枚文件不能列出收录它的那一刀——台账 HEAD 已把这条写成规矩）。

⛔ 未 push、⛔ 未开真窗、⛔ 未改 `docs/reports/pending-and-issues.md`／`HANDOVER.md`／任何票面、⛔ 未勾任何 AC 框、
⛔ 未新建 `.sh`/`.ps1`/`.txt`/`.out`、⛔ 未删任何临时件、⛔ 未跑 build/vet/test（只跑过 `go list`/`go env`，见 §0）。
