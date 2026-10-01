# 33-r7 面板宿主 C27 — 常驻建窗踩到别人留下的消息：层归属修复腿（写码腿）

票 33 `.scratch/wisp/issues/33-panel-host-c27.md`。本腿只修那一枚顺序依赖红：
整包跑时 `TestAC14AwaitedBindingReplyReachesThePage` 红、隔离单跑绿。
AC 框（`- [ ]`／`- [x]`）一枚没碰，勾与不勾归编排者／验收腿。

## 起手锚点（同发取，`date` + `git log -1 --format=%H` + `rev-parse --abbrev-ref HEAD`）

- 时刻：`2026-10-01 15:41:26 +08`
- HEAD：`416d9d56977a89a070e29e23391b1a2f02bdc9b5`（标题 `docs(33-r6 收尾两枚落盘)：A504＋它表尾的编排者标注……`）
- 分支：`dev`
- 台件目录：`.scratch/wisp/probes/33/r7/`（现读 `ls .scratch/wisp/probes/33/`＝a1 a2 h1 p1 r1 r2 r3 r4 r5 r6 ⇒ r7 未占）
- 证据件号：`ls docs/evidence/s1/ | grep -i r7` 零命中 ⇒ `33-panel-host-c27-r7.md` 未占
- 起手桌面枚数（`tasklist //FI "IMAGENAME eq msedgewebview2.exe"`，15:41）：6 枚（本机 owner 也可能自己开着 Edge ⇒ 只报数，不据此指控谁）

## 判据进度表

| 格 | 内容 | 状态 |
|---|---|---|
| ① | 最小确定性复现（谁先跑、留下什么、`bringUp` 读到什么）＋逐字读数＋命令 | 已填 |
| ② | 定性：产品缺陷 vs 测试卫生缺陷（两边说法都写＋我选哪支、凭什么） | 已填 |
| ③ | 修法＋反控（定向突变必须当场红，抄原始红句逐字） | 已填（两枚反控各 3/3、2/2 红） |
| ④ | 整包逐名红册（PASS/FAIL/SKIP 三数＋逐名红，目标 rc=0 且红册空） | 发 1 已填（rc=0／240-0-0／红册空）；发 2 待补 |
| ⑤ | 门禁四数（build / vet / d22scan / gofumpt） | 已填（build/vet 终态复跑待补同发读数） |
| ⑥ | 我可能写错的条目（附"如果错了后果"） | 已填（8 条） |
| ⑦ | 判不动的地方（逐条 甲／乙／不做 ＋现量＋为什么判不了） | 已填（甲我选／乙复算过不选／丙＋两格具名交回） |
| 交件判语 | 修好了没有＋归口哪一层 | 待填 |

---

## ① 最小确定性复现（谁先跑、留下什么、`bringUp` 读到什么）

### (a) 不靠调度器运气的那一发：同一条锁死线程上定植"活窗＋一枚没派的 `WM_CLOSE`"

台件＝`.scratch/wisp/probes/33/r7/zz_33r7_probe_windows_test.go.bak`（探针源码，`go test` 不认 `.bak`，
所以它不进名册；要复算：`cp` 成 `cmd/wisp/zz_33r7_probe_windows_test.go` 再按下面命令跑，跑完 `mv` 回来，
⛔ 不删）。产码＝HEAD `416d9d56`（`bringUp` 只有窄形 `drainStaleQuitBeforeCreate`）。

命令（一发一个进程、一种 mode，因为一次失败的建窗能把整枚二进制带走，见 (c)）：

```
GOFLAGS= go test -c -o .scratch/wisp/probes/33/r7/r7.test ./cmd/wisp
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP33R7MODE=<mode> \
  ./.scratch/wisp/probes/33/r7/r7.test -test.run TestP7ProbePlantThenMode -test.v -test.timeout 120s
```

步骤：锁线程 → `bringUp#1`（起真窗）→ 泵到空（让定植不含歧义）→ `mgr.Destroy()`（只 **POST** `WM_CLOSE`，不派）
→ 按 mode 处置 → `bringUp#2`（产品自己的 `bringUp`）。逐字读数（`modes-grid.txt`，每 mode 三发）：

```
### mode=none iter=1 rc=0
    MODE none              panicked=runtime error: invalid memory address or nil pointer dereference created=false elapsed=0ms
       | tid=21212
       | bringUp#1 ok, pumped 0 message(s) to idle, thread windows=3
       | after Destroy (undispatched): head=WM_CLOSE(0x10) hwnd=0x220EDE GetQueueStatus mask=0x0
       | no purge (HEAD shape): bringUp runs its own WM_QUIT-only drain
       | pre-create peek: filtered WM_QUIT found=false (msg=WM_NULL) | unfiltered head=WM_CLOSE(0x10) | thread windows=3
```

`bringUp` 读到的东西＝**一枚没派的 `WM_CLOSE`（0x10）＋它所属的那扇还活着的窗**；
`none`／`dispatch` 两形 6/6 panic，与票面上那枚 recovered 栈同一签名（`pkg/edge.(*Chromium).Init` → `chromium.go:131`）。

| mode（建窗前做了什么） | 三发读数 | 计数 |
|---|---|---|
| `none`＝HEAD 的窄形（只 `PeekMessageW(WM_QUIT..WM_QUIT, PM_REMOVE)`） | `panicked=... nil pointer dereference created=false elapsed=0ms` | 0/3 建起 |
| `dispatch`＝先把那枚 `WM_CLOSE` **派**出去（孤儿窗真被拆：`thread windows now=0`），再跑 HEAD 的窄排干 | `quit-drain removed 0` ＋ `pre-create peek: filtered WM_QUIT found=false` ＃ 紧接着照样 `panicked=... nil pointer` | 0/3 建起 |
| `wide`＝`PeekMessageW(0,0,PM_REMOVE)` 清整条、只删不派 | `wide purge removed 1: WM_CLOSE(0x10)` → `created=true` 两发，但 `thread windows now=4`（上一任的窗留在世上）；第三发**挂死**（见 (c)） | 2/3 建起（1 发挂死） |
| `pumpAfterDestroy`＝派到底、泵到 `PeekMessage` 报空 | `pumped-and-dispatched 3(4/2) message(s); thread windows now=0` ＋ `unfiltered head=empty` → `created=true` | 3/3 建起 |

⇒ **`bringUp` 层的"排干范围"这一维被穷举过了**：窄形看不见闩锁起来的 quit（`removed=0` 却照样被 `GetMessageW` 交出来），
宽形要么留下一扇没人认领的窗（每建一次多一扇）要么吃掉本次建窗自己要用的消息，
唯一 sound 的那一形（泵到底）做的事是**把孤儿窗拆掉**——那是线程的 owner 才有的知识，不是 `bringUp` 的。

### (b) "谁先跑、留下什么"：家族 A/B 一发里 TID 逐枚对上（不是整包那发；整包那发的顺序在末尾引前人日志）

`mutation-2-no-seal.txt`（＝本腿修法里把 owner 侧那道封（`releaseThreadClean`）拿掉、其余全在位，
`-run 'TestAC13|TestAC14|TestPanelThread|TestBallPanelGestures' -count=3`）第 2 轮逐字：

```
    panel_resident_windows_test.go:482: AC#13 reused-thread release: tid=16272 dispatched 0 message(s) before unlocking; windows left on that thread=0 queue head=
    panel_resident_windows_test.go:566: AC#13 queued-close: tid=16272 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0x0, plantQueued=false), then bringUp#2 err=first bringUp (the control this plant needs): panel host: refusing to create the panel window: this thread still has an undispatched WM_CLOSE queued for hwnd 0x143082A. ...
```

⇒ **同一次运行里、同一个 OS 线程号（16272）**：泄漏者没泵就还池，下一条用例的第一次建窗落到它上面，
产码的检查看到的正是上一任欠着没派的那枚 `WM_CLOSE`（hwnd `0x143082A`）。
第 3 轮同一枚泄漏者还污染的**另一条**路：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 红 20.00s，
红句就是票上那句 `timed out after 15s waiting for the panel thread to finish Show`，
而那发日志里 `refusing to create` 0 次、`goroutine panic recovered` 0 次＝**既不拒绝也不 panic 的第三种形：挂在建窗自己的 `GetMessageW` 上**（与 (a) 的 `wide` 挂死同一去处）。

票上原先那发整包读数（`.scratch/wisp/probes/33/r6/fullpack.log:601-614`，14:52，HEAD `4bd32fe0`＋未提交的窄形改动，
**本腿只读该文件、未复跑**）顺序与此完全一致：
`--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (1.08s)` 紧接着 `=== RUN TestAC14AwaitedBindingReplyReachesThePage`
→ `goroutine panic recovered ... panel_host_windows.go:241 bringUp ... loop(:210)`。
⇒ 编排者 15:27 那条"头一次建窗就踩到别人留的毒"的判读**成立**，但毒源不是"别的线程被别家毒过"这种无主形状：
`panel-sta` 落到的是一枚**测试自己还池的、还欠着一枚 `WM_CLOSE` 的线程**，而 33-r6 那枚钉就是最后写它的人。

### (c) 顺带量到的两条后果（都不是"测试难看"那么轻）

- 五种 mode 挤在一发进程里跑：第 2 发直接 **SIGSEGV，rc=139**，`plant-modes-2.txt` 只落到第 1 发的读数
  ⇒ 一次 panic 之后留下的半初始化 controller 能让**整枚测试二进制**没掉，不是"这枚用例红"。
- `wide` 第 2 发：`modes/wide-2.txt` rc=2，`panic: test timed out after 2m0s`，卡住的协程栈顶＝
  `pkg/edge/chromium.go:100` 的 `GetMessageW`（`Embed` 的循环）⇒ 33-r6 那句"整条排干会卡住下一次建窗"
  **本腿自己复现出来了**（1/3），另有一发整包读数同形（`.scratch/wisp/probes/33/r6/fullpack2.log:640`，
  `pre-create message purge hit its cap` removed=256，show 永不完成；那是 `4be1d3f1` 宽形的整包日志，本腿未复跑）。

## ② 定性：产品缺陷 vs 测试卫生缺陷（两边说法）

**测试卫生缺陷（我选这支作触发者）**：全仓"泵过一扇 webview 窗又把 M 还给池"的点，现读只有两枚——
`cmd/wisp/panel_resident_windows_test.go:468`（`TestAC13BringUpSurvivesAReusedThreadQuit`，33-r6 新加的钉自己）
与 `cmd/wisp/panel_host_windows_test.go:90`（`hostThreadHarness`，HEAD 已改成锁到底、注释写明理由）。
出货拓扑没有这一形：`panel_resident_windows.go:202` 锁而不解、协程退出时线程随它销毁、窗与 quit 一起没。
`cmd/wisp/notify_windows.go:139` 的解锁落在 message-only 窗上，解锁前 `DestroyWindow` 已同步跑完，
且本仓 `PostQuitMessage` 的调用点现读只有依赖的 `Terminate`、`terminateOnThisThread`、
`internal/ball/sta_windows.go:161` 与 testdata 那枚独立程序——前三枚都不在 `notify` 那条线程上。
凭据＝(b) 那发 TID 逐枚对上的读数（泄漏者→下一条用例的建窗被具名拒绝）。

**产品缺陷那一支（不是没道理，但它要说的是另一件事）**：`bringUp` 隐含一条从没写下来的前提
"来建窗的线程没有别人欠着的消息"，而 `PanelManager` 文档化的 recreate 面（`Destroy` 之后 `Show` 重建）
在同一条线程上正是这个形状。这一支的说法是"产码该把前提变成自己能检查的"。
**我只走到这里就停**：读数显示检查得到（过滤式 peek 能读到那枚 close）、**拒绝得住（不 panic、不挂死、具名）**，
但**修不干净**——四种排干范围都被量过（①(a) 表），其中 sound 的那一形做的事是拆掉别人的窗。
所以产码这一侧落的是"检查＋具名拒绝＋不动别人的队列"，owner 那一侧落的是"泵到底再还池"。

⇒ 归口：**甲（线程的 owner 那一层）挡住它**；`bringUp` 那一层只负责"下一次再有这种东西时，
当场说得清是谁欠的"，它不是挡的那一层。①(b) 那发就是这句话的读数：泄漏被拿掉封时，
拒绝立刻点名 hwnd；泄漏被封住时，包内不再有毒可踩。

## ③ 修法＋反控红句

**改了什么（两枚文件，全在 `cmd/wisp/**`）**

1. `cmd/wisp/panel_resident_windows_test.go`（甲，落地的就是这一枚）：
   新增 `pumpThreadToQuiet` / `threadWindowCount` / `threadQueueHead` / `releaseThreadClean` 四枚工具，
   `TestAC13BringUpSurvivesAReusedThreadQuit` 在 `Destroy` 之后、`UnlockOSThread` 之前**把线程派到空**，
   并加两条断言：还池前 `windows left == 0` 且 `queue head == empty`，红句指名"这条线程被还池后
   下一条用例的建窗会踩到什么"。原来三条断言一枚没删、15 秒没放宽。
2. `cmd/wisp/panel_host_windows.go`（乙里读数支持的那半）：`bringUp` 在建窗前跑完窄形排干之后，
   用 `staleCloseQueued()`（`PeekMessageW([WM_CLOSE..WM_CLOSE], PM_NOREMOVE)`，不删不派）检查，
   脏就**返回具名错误、根本不建窗**（没有重试循环、没有清理）；`bringUp` 上面那段注释换成 ① 的四形读数，
   把 33-r6 那句〔仅死腿自述〕的"zombie controller"换成本腿自己的复算（1/3 挂死 + `thread windows 3->4`）。
   新增钉 `TestAC13BringUpRefusesAThreadWithAQueuedClose`：定植同一形状，断言
   无 panic、错误里含 `WM_CLOSE`、`created=false`、**拒绝后那枚 close 仍在队列里**（证明产码没吃别人的消息）、
   以及本用例自己还池前线程干净。

**反控 1＝拿掉产码那一处检查**（`git` 未提交的临时删；快照与还原：
`.scratch/wisp/probes/33/r7/panel_host_windows.go.keep` md5 `00336ed0704e51f4cdb043a6d0b3b852`，
跑完 `cp` 回去、md5 复确一致）。`-run 'TestAC13|TestAC14|TestPanelThread|TestBallPanelGestures' -count=2`
逐字红句（`mutation-1-no-check.txt`，2/2 都红，其余 6 枚全绿）：

```
--- FAIL: TestAC13BringUpRefusesAThreadWithAQueuedClose (1.19s)
    panel_resident_windows_test.go:574: bringUp on a thread carrying a queued WM_CLOSE runtime error: invalid memory address or nil pointer dereference instead of refusing it by name. ... Fix: bringUp must call staleCloseQueued before NewWithOptions and return an error
    panel_resident_windows_test.go:577: bringUp did not name the refusal: err=<nil>, want an error that says a WM_CLOSE is queued on this thread. ...
```

同一发里 `--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.13s)`／`(0.88s)`
⇒ **绿 AC#14 的是甲那一道封，不是这枚检查**；检查的牙只由它自己那枚钉兜（当场红、panic 原文抄在上面，
不是 15 秒超时）。另注：这发红句里 `windows=0 queue head=` 那些字段是零值——协程已经 panic 展开，
后面的赋值没执行，别把"queue head=" 读成"队列为空"。

**反控 2＝拿掉测试那一道封**（同一手法，快照 `panel_resident_windows_test.go.keep` md5 `7451332c7ef1b15b889a586e948f8473`）：
`mutation-2-no-seal.txt`，`-count=3` 里 `TestAC13BringUpSurvivesAReusedThreadQuit` **3/3 红**：

```
    panel_resident_windows_test.go:497: this test is about to hand OS thread 7544 back to the pool with its queue non-empty (). A queued WM_CLOSE on that thread becomes a WM_QUIT inside the next create's own pump (webview.go:242-243 + 381-383): measured 3/3 panic, and no pre-create quit drain can see it while other messages are still queued
--- FAIL: TestAC13BringUpSurvivesAReusedThreadQuit (0.93s)
```

（`queue non-empty ()` 里的空串同样是零值：这一形里 `releaseThreadClean` 整句没跑，
`head` 从没被赋过值——断言仍然响，是因为"没跑"本身就等于"没干净"。）
并且这一发的第 2、3 轮把它毒到的**下游**也拍下来了：`tid=16272` 被下一条用例的建窗**具名拒绝**、
`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 红 20.00s（①(b) 逐字）。

**带修常态**：见 ④（整包名册）与 `.scratch/wisp/probes/33/r7/nails-fix.txt`
（两枚钉 `-count=3` 全 `PASS`，读数是 `dispatched 3(4) message(s) ... windows left=0 queue head=empty`
与 `refusing to create ... created=false ... the close was still queued after the refusal=true`）。

## ④ 整包逐名红册

命令（两发同一条，PATH 少了它就是"没跑"而不是绿）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -v -timeout 25m
```

### 发 1（`.scratch/wisp/probes/33/r7/fullpack-1.txt`，共 1347 行：`ok` 在 1329、`rc=0` 在 1330、尾 `date` 1331、1332-1347 是 tasklist 块；首枚日志行 `16:18:31`、尾行 `Thu Oct  1 16:21:55 CST 2026`，那发起手的 `date` 落进了任务 stdout 而非本文件）

- **rc=0**，`ok  	github.com/CarlosShao/wisp/cmd/wisp	204.702s`
- 三数（只认 `--- FAIL`/`--- PASS`/`--- SKIP` 行）：**PASS=240、FAIL=0、SKIP=0**
- **逐名红册＝空**（`grep -cE '^--- FAIL|^    --- FAIL'` = 0）
- 桌面枚数：起手（15:41）6 枚 → 这发终态 13 枚。逐枚查过归属：6 枚属 `SearchHost`（任务栏搜索，起手那 6 枚就是它）、
  6 枚属 owner 自己的 `clipsync-desktop.exe`、**1 枚是我这发之前探针跑留下的孤儿**（PID 3112，
  `--user-data-dir=C:\Users\swq\AppData\Local\Temp\wisp-33r7-plant-dispatch\EBWebView`，父进程 16064 已不存在）。
  16:23 我已 `taskkill //PID 3112 //F` 收掉它，收后 `Get-CimInstance` 里再无 `wisp*` 目录的 webview 进程（13→12）。
  ⇒ 不是"机器上本来就该有 13 枚"，也不是"跑了整包就多 7 枚"：**多出来的一枚是我的，已经清掉**。
- 面板家族逐名（这一发里真跑了、不是被跳过）：

```
--- PASS: TestPanelHostRealWindowHopAndLifecycle (1.10s)
--- PASS: TestAC4FocusReturnToPriorWindowGap33r5 (0.50s)
--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.97s)
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (0.80s)
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (1.02s)
--- PASS: TestAC14AwaitedBindingReplyReachesThePage (0.49s)
--- PASS: TestAC14GoSideEvalPushReachesThePage (0.50s)
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (0.34s)
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
--- PASS: TestBallPanelGesturesReachThePanelThread (0.34s)
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (0.50s)
```

  两枚新落地的读数在这一发里也在（不是我另跑的定向读数）：

```
AC#13 reused-thread release: tid=7412 dispatched 3 message(s) before unlocking; windows left on that thread=0 queue head=empty
AC#13 queued-close: tid=13284 first bringUp ok, pumped 0 to idle, Destroy left a queued close (hwnd 0xFE0E82, plantQueued=true), then bringUp#2 err=panel host: refusing t...
```

### 发 2（`.scratch/wisp/probes/33/r7/fullpack-2.txt`）：**rc=1、两枚红，都不是本案那一形**

- 时刻：`Thu Oct  1 16:24:24 CST 2026`（文件首行 `date`）→ `Thu Oct  1 16:28:32 CST 2026`（`rc=1` 之后那行 `date`），
  **rc=1**，`FAIL	github.com/CarlosShao/wisp/cmd/wisp	245.550s`（发 1 是 204.702s ⇒ 整包慢了两成）
- 三数：**PASS=238、FAIL=2、SKIP=0**
- 逐名红（红句逐字）：

```
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (1.41s)
    panel_host_windows_test.go:600: hot re-show 292.1 ms exceeds D32 panel hot budget 200 ms
--- FAIL: TestPanelHostLatencyPercentilesAC2 (0.00s)
    panel_host_windows_test.go:817: hot re-show P95 292.063 ms over 1 runs exceeds the D32 panel hot budget 200 ms
```

  同一发里的读数上下文：`cold bring-up measured on this box: 967.450 ms`、
  `hot re-show measured on this box: 292.063 ms`（发 1 同一枚用例是 **60.854 ms**）、
  `machine-wide msedgewebview2=19`（发 1 是 20，而这台机器 15:41 起手时是 6）。
  两枚红是**同一个因**（同一枚 292 ms 的样本被两把尺各读一次），不是两件事。
- **本案那一族在这发里全绿**（这就是为什么我说这两枚红跟毒无关）：

```
--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (1.11s)
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (0.89s)
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (1.01s)
--- PASS: TestAC14AwaitedBindingReplyReachesThePage (0.62s)
--- PASS: TestAC14GoSideEvalPushReachesThePage (0.55s)
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (0.43s)
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
--- PASS: TestBallPanelGesturesReachThePanelThread (0.41s)
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (0.51s)
```

- 我这发改的东西**不在热点路径上**（这是读码，不是读数）：`Show` 在 `created==true` 时直接跳过 `bringUp`
  （`cmd/wisp/panel_host_windows.go:352` 那个 `if !created`），所以热重显既不走 quit 排干也不走 close 检查。
- 复测（跑完后趁机器安静立刻单发定向，`latency-recheck.txt`，16:29:52，
  `-run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2' -count=3`）：
  热重显 **37.089 / 36.387 / 40.290 ms**，两枚用例 3/3 `--- PASS`；冷建窗 730.309 / 403.401 / 323.014 ms。
- 我的判读：**负载噪声**（同一枚尺、同一份码，安静时 37-40 ms，那一发 292 ms，且整包同时慢两成、
  机器上 webview 进程从 6 涨到 19-20）。**但这枚红我不取消、不解释成"不算"**：阈值一字节没动、
  断言一条没放宽，三发数都在这张表上，归口留给编排者（要按哪一发算 AC#2/D32 那一行不是写码腿的权）。
- 发 3 是为此补的：见下。

### 发 3（`.scratch/wisp/probes/33/r7/fullpack-3.txt`，共 1347 行；行号：1＝`Thu Oct  1 16:30:09 CST 2026`、1329＝`PASS`、1330＝`ok` 行、1331＝`rc=0`、1332＝`Thu Oct  1 16:33:34 CST 2026`、1333 起是 tasklist 块）

- **rc=0**，`ok  	github.com/CarlosShao/wisp/cmd/wisp	201.376s`
- 三数：**PASS=240、FAIL=0、SKIP=0**；**逐名红册＝空**
  （`grep -cE '^--- FAIL|^    --- FAIL'` = 0）
- 发 2 那两枚红的那一把尺，这一发回到帽内：`hot re-show measured on this box: 86.076 ms`（预算 200 ms）、
  `cold bring-up measured on this box: 758.174 ms`（预算 1500 ms）、`our tree webview=7`、
  `machine-wide msedgewebview2=19`——**同一发里机器并没有更安静**（19 枚与发 2 的 19 枚同级），
  所以 86 ms vs 292 ms 的差别我只能说"那一发恰好挤了一下"，说不出是谁挤的
- 本案那一族逐名：

```
--- PASS: TestPanelHostRealWindowHopAndLifecycle (1.00s)
--- PASS: TestAC4FocusReturnToPriorWindowGap33r5 (0.52s)
--- PASS: TestPanelHostLatencyPercentilesAC2 (0.00s)
--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.97s)
--- PASS: TestAC13BringUpSurvivesAReusedThreadQuit (0.98s)
--- PASS: TestAC13BringUpRefusesAThreadWithAQueuedClose (1.06s)
--- PASS: TestAC14AwaitedBindingReplyReachesThePage (0.59s)
--- PASS: TestAC14GoSideEvalPushReachesThePage (0.45s)
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (0.43s)
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
--- PASS: TestBallPanelGesturesReachThePanelThread (0.37s)
--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (0.48s)
```

- 桌面枚数：终态 12 枚，逐枚按 `--user-data-dir` 归属查过，**属于我这发的＝0 枚**
  （`Get-CimInstance ... | grep -ci wisp` = 0；那 12 枚分属 `SearchHost` 与 owner 自己的 `clipsync-desktop`）。
  发 1 终态里我那枚孤儿（PID 3112）已在 16:23 收掉，之后两发都没再留下带 `wisp*` 目录的进程。
- 三发汇总：**发 1 rc=0／240-0-0**、**发 2 rc=1／238-2-0（两枚红＝同一枚 292 ms 的热重显被两把尺各读一次）**、
  **发 3 rc=0／240-0-0**。判据④要的"rc=0 且红册为空"**达成过两发**；发 2 那两枚我不判成本案的回归，
  理由与"如果错了后果"写在 ⑥(i)，归口留给编排者。

## ⑤ 门禁四数

| 尺 | 命令 | 读数 |
|---|---|---|
| build | `GOFLAGS= go build ./...` | rc=0（16:16 段，gofumpt `-w` 之前） |
| vet | `GOFLAGS= go vet ./cmd/wisp` | rc=0（同上；`-w` 之后另发复跑，见下方终态复跑） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**、`d22scan: clean - no D22 ban violations`；分母 live scope：bans #1-5 internal/=224、bans #1-5 cmd/=34、ban #6 frontend/=85、ban #7 internal/tools/=23、ban #8 design/=39、frontend/=85、internal/=476、**cmd/=81（注释与 `_test.go` 都在射程）**；同发 `runtests.sh: OK - packages=[./...] PASS=34 FAIL=0 SKIP=0`（含 `tools/d22scan` 自己的 16.248s）；忽略过滤只跳过 1 枚 `frontend/dist/assets/` 下的文件（按 git index 判定） |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/panel_host_windows.go cmd/wisp/panel_resident_windows_test.go` | 起手点出 `panel_resident_windows_test.go` 不合规 ⇒ `-w` 后**列表为空**。改动是纯格式：与 `-w` 前的快照逐枚 diff＝只删掉第 590 行一枚空行（快照 `.scratch/wisp/probes/33/r7/panel_resident_windows_test.go.pregofumpt`） |
| staticcheck | — | **〔未复认〕**：本机版与 CI 钉版不同，不出结论 |

## ⑥ 我可能写错的条目（逐条，附"如果错了后果"）

- (a) **"整包那枚红的毒源＝`TestAC13BringUpSurvivesAReusedThreadQuit` 释放回池的那条线程"**：
  机制我用自己的台件复算过（同锁线程定植"活窗＋一枚没派的 `WM_CLOSE`"⇒ `bringUp` 3/3 panic，
  `.scratch/wisp/probes/33/r7/modes-grid.txt`），但"包内确实是这一枚用例把线程还了池、且 AC#14 落在它上面"
  是**归因**。它有两枚凭据：一是 ①(b) 那发把 TID 逐枚对上了（泄漏者还池 `tid=16272`，下一条用例的第一次建窗
  就落在 16272 并被产码具名拒绝、连 hwnd 都点了），二是③的反控 2 里同一形状 3/3 响。
  **但"票上原先那发红也必然是这一枚用例造成的"这一句仍然只到"顺序相邻＋形状相同"级**——
  我复跑的是 33-r6 那枚钉自己（前人整包日志我只读未复跑）。
  若错（真凶是别的还池点，或根本是 dataPath 复用/并发建窗），判据④仍会红，且我的修法（甲封点）不解决问题；
  实际三发名册里有两发整包红册为空，这一支就算不完整也不影响绿，但会影响"下一枚同类红该往哪儿找"。
- (b) **"`WM_QUIT` 是闩锁（latched），过滤式 `PeekMessageW([WM_QUIT,WM_QUIT])` 在队列里还有别的消息时看不见它"**：
  读数支持（mode=dispatch 三发：drain removed=0、filtered found=false，紧接着的建窗照样 panic），
  但"闩锁"是**我对机制的命名**，不是 Win32 给我回的字段。若机制其实是别的
  （例如销毁过程在 drain 之后才异步投 quit），那"任何排干范围都不可靠"这句仍成立（因为读数是 removed=0 + panic），
  而我在 ② 里写的"闩锁"这个词要降格成假设。
- (c) **"整条队列排干会吃掉本次建窗自己的完成回调 ⇒ 挂在 `GetMessageW`"**：机制是假设。
  盘上事实只有两条：mode=wide 三发里 1 发 rc=2（`-test.timeout 120s` 到点，栈顶＝
  `pkg/edge/chromium.go:100` 的 GetMessageW，`modes/wide-2.txt`）；整包 `4be1d3f1` 那发
  `pre-create message purge hit its cap` removed=256 且 show 永不完成（`.scratch/wisp/probes/33/r6/fullpack2.log:640`，
  本腿未复跑＝引前人读数）。若"吃掉回调"这句错，两条读数仍然只支持"宽排干不可靠"，不支持我给的因。
- (d) **"出货拓扑今天不会中这一形（面板线程锁到底、队列建窗前必空）"**：我只读了 `panel_resident_windows.go:202`
  的锁而不解、`notify_windows.go:138-139`（message-only 窗、解锁前 `DestroyWindow`、全仓 `PostQuitMessage` 只有
  依赖的 `Terminate` 与我们自己的 `terminateOnThisThread` 两处）。**没**逐行读 `internal/ball` 的 ui-sta 与
  托盘那条腿。若还有别的"在面板线程上起过窗又投过 close"的路径，出货面也会中，且我的定性（测试卫生）要改。
- (e) **"产品侧 `bringUp` 拒绝建窗不会误杀合法状态"**：判据只认"队列里有一枚 `WM_CLOSE`"。
  `TestAC4PriorFocusSurvivesARefusedPanelSample` 在面板线程上先起了编辑器窗再 `showAndWait`，
  那一发合法地带着别的消息建窗——若过滤式 peek 在那一发里读到 `WM_CLOSE`，我就会把一枚绿用例改红。
  读数见 ⑤/④（整包若在那枚用例红，就是这支错了）。
- (f) **"`pumpAfterDestroy` 3/3 绿＝owner 侧排干是 sound 的修法"**：三发都是**同进程同线程连着建第二扇窗**。
  它没有覆盖"owner 排干后把线程还给池、再由**另一枚协程**拿去建窗"这一形（那是包内真实形状，
  受调度器支配，我只能靠整包名册验）。若差别要命，判据④会红。
- (g) **写面外还有一枚 owner 我没封到，也封不动**：`internal/ball/sta_windows.go:161` 的 `quit()` 投
  `PostQuitMessage(0)`，而 `:52-53` 的 `start` 是 `LockOSThread` + `defer runtime.UnlockOSThread()`
  ——球的那条 ui-sta 收摊时把一条**可能带着闩锁 quit 的线程还给 Go 池**，这正是甲层要封的形状，而它在
  我写面外（⛔ 派单禁动 `internal/**`）。今天这一发整包没中它（红册为空），我说不出为什么没中——
  **别把"没中"读成"不会中"**。若它其实会毒到后来的 `panel-sta`，后果＝又一枚顺序依赖红，
  且我这发的两道修都拦不住它（close 检查看不见闩锁 quit，见 (b)）。归口在 ⑦，交编排者裁。
- (h) **`releaseThreadClean` 之后还池的线程"必然干净"这件事只在这台机器上量过**：
  `nails-fix.txt` 3/3 是
  `windows left=0 queue head=empty`、`dispatched 3-4 message(s)`。若哪天 WebView2 在窗销毁后还继续异步投消息，
  泵到空这一步会撞 4096 的帽（撞帽会 `slog.Warn` 自陈）或把 `head` 留成非空 ⇒ 我新加的两条断言会红——
  红得对（那种线程本来就不该还池），但那是**新红**，不是回归，别把它当成"33-r7 把绿改坏了"。
- (i) **发 2 那两枚红我判成"负载噪声"**：凭据是三样——同一发里本案那一族 9 枚全绿、
  趁机器安静复测 `-count=3` 得 37.089/36.387/40.290 ms（`latency-recheck.txt`）、
  以及读码说热点路径不经过我改的任何一处（`Show` 的 `if !created`）。
  **如果这判错了**（真是我这发的东西把热重显推到 292 ms，比如我新建又拆掉的窗让 WebView2 运行时整体变慢），
  后果＝D32 面板那一行的达标情况被我这句话掩盖了，而阈值与断言我都没动、也没有权去动 ⇒
  这一格只能由编排者按发 3 与后续 CI 读数裁，我把三发数都留在 ④ 上、不替它挑一发好看的。

## ⑦ 判不动的地方（逐条：甲＝按我读数／乙＝别的形／不做 ＋现量＋为什么判不了）

- **甲（我选的）**：层＝**线程的 owner**。谁的线程谁负责"还池／再建窗前把它泵到静"，
  出货侧已经是这一形（`loop` 锁到底）；包内唯一破口是测试把带活窗的线程还了池。
  现量＝`modes-grid.txt` 四形对照（none 3/3 panic、dispatch 3/3 panic、wide 2/3 建起但 +1 泄漏窗且 1/3 挂死、
  pumpAfterDestroy 3/3 建起且 windows 3→0）。
- **乙（复算过，不选）**：产码 `bringUp` 自己把线程弄干净。三种排干我都量过：窄形拦不住 close 引发的闩锁 quit
  （3/3 panic），宽形要么泄漏上一任的窗（`thread windows now=4`）要么吃掉本次建窗的消息挂死（1/3 rc=2 +
  fullpack2 的 cap=256），"派完再窄排"照样 3/3 panic。
  **我只保留了乙里可判定的那一半**：`bringUp` 把"队列里有没有 close"变成**能检查、能拒绝**的前提
  （检查＋具名拒绝＝读数支持；清理＝读数反对）。代价写在产码注释与 ③：脏线程上面板**不开**并留一句名，
  而不是半初始化的 controller（那枚能整机 segfault，见 ①）。
- **丙**：读数没指向别处，但**"为什么这条 M 被分给新协程"这一维我没判**——它属 Go 运行时调度，
  不是 cmd/wisp 能约束的；我只把"落到脏 M 上会怎样"变成可检查的。
- **判不动／要人拍板的一格（具名交回）**：`internal/ball/sta_windows.go` 的 `quit()`（`:161` 投 `PostQuitMessage`）
  ＋ `start`（`:52-53` `LockOSThread`/`defer UnlockOSThread`）＝球那条 ui-sta 收摊时把线程连同 quit 一起还池。
  这属甲层（owner 侧），但在 `internal/**`＝**本腿写面外**，我没读它、没改它、也没读数说明它今天会不会命中。
  三条路摆给编排者：甲′＝派一枚能写 `internal/ball` 的腿把同一道"泵到空再还池"落进去；
  乙′＝承认 `bringUp` 拦不住闩锁 quit、把它记成 DEFERRED 一栏；不做＝按今天的名册（④那发没中）先走。
  为什么我判不了：要复现它得让球的 ui-sta 收摊与面板首建窗在**同一条 M**上相遇，那要么动 `internal/**`
  要么在 `cmd/wisp` 里造一枚跨包的定植器（后者测的是我自己造的假现场，不是它的代码路径）。
- **这发故意没做的一格（不是没想到，是范围）**：`showAndWait` 的等待条件是
  `IsShown() || startUpErr() != nil`，而 `Show` 返回的错误只落在日志里
  （`panel_resident_windows.go:324-327`），所以"建窗被具名拒绝"在读数上仍然表现为**15 秒超时**——
  ①(b) 那发 20.00s 的红就是这个形状。把它变成"当场红＋红句里带那句名"要给 `residentPanel` 加一枚
  lastShowErr 字段（产码面第三处改动），收益只是把已经能命名的东西命名得快一点，票 33 的判据没要它。
  ⇒ 记在这里交回，不由我这一发顺手加。
- **不做（范围与禁令，逐条）**：⛔ 不碰依赖（`pkg/edge` 那三行没有 `if`，改不了它，只能不制造前提）、
  ⛔ 不动 `internal/**`、⛔ 不动 `shutdown.go`/`shutdown_hooks.go`、⛔ 不搬 `winlive`、⛔ 不加 `t.Skip`
  到新用例、⛔ 不放宽那 15 秒、⛔ 不把断言降级成 `t.Logf`、⛔ 不删 AC13/AC14 任何一枚、
  ⛔ 不给 `bringUp` 加"重试到绿"的循环（拒绝＝一次性具名错误，没有循环）、⛔ 不改票面 AC 框、
  ⛔ 不动台账／HANDOVER／票面（编排者独占）。
