# 33-r8 悬浮球 ui-sta — 线程的 owner 在还池前把自己起过的窗拆掉、把队列泵到静（写码腿）

票 33 收尾的一格。本腿只做一件事：**把"线程的 owner 在把线程还给 Go 池之前先把它泵到静"
这道不变式落到球这条 ui-sta 线程上**。AC 框（`- [ ]`／`- [x]`）一枚没碰（票 33 现状 13 未勾／1 已勾，
勾与不勾归编排者／验收腿）；`docs/reports/**` 未动；`cmd/wisp/**`、`internal/panel/**`、`internal/config/**`
未动（那是 248-r1 的写面）；`shutdown.go`／`shutdown_hooks.go` 未动。

## 起手锚点（同发取，`date` + `git log -1 --format='%h %ad %s'` + `git status --porcelain`）

- 时刻：`2026-10-02 08:25:25 +08`
- HEAD：`a35f7f5e`（`Fri Oct 2 00:04:33 2026 +0800` · 标题 `ledger(A511)：33-r9 死在最后一步……`）
- 分支：`dev`
- 台件目录：`.scratch/wisp/probes/33/r8/`（现读 `ls .scratch/wisp/probes/33/`＝a1 a2 h1 n1 p1 r1 r2 r3 r4 r5 r6 r7 r9 v2 winc1
  ⇒ **r8 未占**；r9 是那一枚死在最后一步的腿）
- 证据件号：`ls docs/evidence/s1/ | grep -i 33` 名册里 **没有** `33-ball-thread-owner-r8.md` ⇒ 未占
- 起手 `git status --porcelain` 全仓名册：逐字存 `.scratch/wisp/probes/33/r8/logs/status-start.txt`
  （**本仓工作树永远不干净**：`design/` 下是 owner 自己的未提交件，不还原、不提交、不删。
  闸门＝终态名册＝起手名册＋只本腿这几枚路径，⛔ 不是"必须为空"）
- 起手 `internal/ball` porcelain＝0（量法：`git status --porcelain internal/ball/`；与编排者 08:2x 的读数同发一致）
- 尺与环境：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（本腿每一次 `go test` 都带着它）

## 本腿写面（只有这些）

| 路径 | 动作 |
|---|---|
| `internal/ball/sta_release_windows_test.go` | 新建（枚＋读数仪器） |
| `internal/ball/sta_windows.go` | 改（owner 侧收摊那一步） |
| `internal/ball/ball_windows.go` | 改（Close 里同步清 owner 自己的窗记录；只有这一处，见 ③） |
| `docs/evidence/s1/33-ball-thread-owner-r8.md` | 本文件 |
| `.scratch/wisp/probes/33/r8/**` | 台件（只建不删） |

## 判据进度表

| 格 | 内容 | 状态 |
|---|---|---|
| ① | 现量：球这条 ui-sta 今天到底留下什么（两枚形状×逐字读数×命令） | 已填（见 ①） |
| ② | 关键判据：**这条脏线程今天会不会真的毒到后来者**（答"中"还是"没中"＋凭据；⛔ 不许把没中读成不会中） | 待填 |
| ③ | 修法（选哪一条＋另一条为什么被否；反控＝定向突变逐字红句） | 待填 |
| ④ | 钉（"还池前窗数＝0 且队列头报空"＋同批钉住它自己的收尾）；名册与三数 | 待填 |
| ⑤ | 门禁四数（build / vet / d22scan / gofumpt；staticcheck＝未复认） | 待填 |
| ⑥ | 整包名册（`./internal/ball/` 全跑三数＋逐名红；⛔ 不跑 `cmd/wisp`） | 待填 |
| ⑦ | 我可能写错的条目（附"如果错了后果"） | 待填 |
| ⑧ | 判不动的地方（逐条 甲／乙／不做 ＋现量＋为什么判不了） | 待填 |
| 交件判语 | 修好了没有＋归口哪一层 | 待填 |

---

## ① 现量：球这条 ui-sta 今天留下什么

### (a) 形状从"读码"变成"读数"的仪器

`internal/ball/sta_release_windows_test.go`（`//go:build windows`，**默认层**，与 33-r5/r7/r9 那一家同层：
真窗、默认层、造不出窗的机器是红不是跳）。两枚用例都**把 `staThread.start` 跑在测试自己锁死的线程上**，
理由不是风格：`PeekMessage` 只能看**调用线程**的队列，所以"还池那一刻队列里有没有东西"这个问题，
世界上只有 owner 自己在那个时刻能答。`EnumThreadWindows` 可以跨线程问，所以窗数由报告侧再对一遍。

命令：

```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
GOFLAGS= go test ./internal/ball/ -count=1 -timeout 90s -run 'TestSTARelease' -v
```

### (b) 起手形状复核（编排者给的行号当待验断言，本腿自己复认＝**全部对上**）

- `internal/ball/sta_windows.go:52-53`：`runtime.LockOSThread()` 之后紧跟 `defer runtime.UnlockOSThread()` ⇒ 对。
- `internal/ball/sta_windows.go:156-162`：`quit()` 先 `PostMessageW(s.hwnd, wmNull)` 再 `pPostQuitMessage.Call(0)` ⇒ 对。
- 补一条本腿读码补上的：`Ball.Close`（`ball_windows.go:935-966`）把 `b.hwnd` 归零、`activeBall` CAS 清空，
  但**从不清 `staThread.hwnd`** ⇒ owner 自己的窗记录在收摊后仍指着一扇已销毁的窗（这一条的后果见 ③ 的第二处改动）。

### (c) 逐字读数（HEAD `a35f7f5e`，未改产码；台件 `.scratch/wisp/probes/33/r8/logs/red-3.txt`）

```
=== RUN   TestSTAReleaseAfterFailedCreateHandsBackNoWindow
    33-r8 create-error door: tid=32692 planted hwnd=0x170D7A | after plant: ball windows=1 all windows=2 classes=[WispBallWindow IME] queue head=empty
    33-r8 create-error door: the instant start() returned: ball windows=1 all windows=2 classes=[WispBallWindow IME] queue head=empty (start err=plant: renderer init failed plant err=<nil>)
    sta_release_windows_test.go:259: staThread.start is about to hand OS thread 32692 back to the Go pool while it still owns 1 live WispBallWindow window(s) ...
--- FAIL: TestSTAReleaseAfterFailedCreateHandsBackNoWindow (0.02s)
=== RUN   TestSTAReleaseAfterPumpExitDispatchesItsQueue
    33-r8 pump-exit door=kept-window: tid=32692 hwnd=0x180D7A createErr=<nil> | after create: ball windows=1 all=2 classes=[WispBallWindow IME] head=empty | AT RELEASE: ball windows=1 all=2 classes=[WispBallWindow IME] head=empty | plant PostMessage rc=1 | this file's own pump moved 0 | after own teardown: ball windows=0 head=empty
    33-r8 pump-exit door=close-shaped: tid=32692 hwnd=0x190D7A createErr=<nil> | after create: ball windows=1 all=2 classes=[WispBallWindow IME] head=empty | AT RELEASE: ball windows=0 all=0 classes=[] head=empty | plant PostMessage rc=0 | this file's own pump moved 0 | after own teardown: ball windows=0 head=empty
    sta_release_windows_test.go:411: door=kept-window: the pump ended and staThread.start is about to hand OS thread 32692 back to the Go pool still owning 1 live WispBallWindow window(s) ...
--- FAIL: TestSTAReleaseAfterPumpExitDispatchesItsQueue (0.01s)
```

### (d) 读数说了什么（三句，逐句凭读数）

1. **窗那一半＝中，两枚门都是。** `create-error` 门（生产可达：`createOnSTA` 在
   `ball_windows.go:219` 就把 `s.hwnd` 记下了，之后 `newRenderer`(:239)／`addTrayIcon`(:247) 还能失败）
   还池那一刻 `ball windows=1`；`kept-window` 门（泵因为 `Close` 以外的理由退出：`GetMessageW` 返回 -1 那一支，
   或这条线程天生带一枚闩锁 quit 时泵的第一圈）同样 `ball windows=1`。
   ⇒ **今天确实会把一条还欠着一扇活窗的线程还给 Go 池**，红因具名，不靠整包顺序撞运气。
2. **队列那一半＝今天没中，而且原因已量到。** 在 `kept-window` 门上本腿定植了一枚 `WM_APP+0x2F8`
   （`PostMessageW` 回答 `rc=1`＝已入队），落点是"quit 之后"；还池前读队列头＝`empty`、
   本文件自己的泵又移走 `0` 枚 ⇒ **球自己的泵在 break 之前会把 quit 之后的消息先派干**
   （这正是 33-r7 那句"闩锁的 quit 在有别的消息排在前头时对窄形 peek 不可见"的另一面）。
   `close-shaped` 门（=生产 `Ball.Close` 的顺序）还池前 `ball windows=0 all=0 head=empty`
   ⇒ **发货那条路径今天读数干净**。"队列头报空"这一枚断言今天两边都是绿的，
   它钉的是形状（未来任何在还池前留在队列里的东西都会被它抓住），本腿不假装它今天响过。
3. **`all windows=2 classes=[WispBallWindow IME]`**：一条 Go 线程的原生窗数不是 0（多出来那枚是 user32 的
   `IME` 窗口，`live_guard_windows_test.go` 里前人以"the Go runtime keeps a hidden helper window"记过这件事，
   本腿读出它的**类名**＝`IME`）。所以"还池前窗数＝0"只能是**球自己起过的那扇**＝0（按类名数），
   ⛔ 不能拿原始计数当分母——那会把一条永远不可能为 0 的数写成必红断言。
   另一发读数：`close-shaped` 门销毁之后 `classes=[]`（`IME` 那枚不在 `EnumThreadWindows` 的返回里），
   说明**类名清单本身会随窗的生死变空**，这枚断言的分母是实的。

### (e) 定植的坑（本腿自己踩到，写下来免得下一枚腿重踩）

- 第一稿把还池前的残留消息定植成 `WM_NULL`：`PostMessageW` 回了 `rc=1`，但队列头读 `empty`、
  本文件自己的泵移走 0 枚 ⇒ **`WM_NULL` 是这类仪器看不见的定植**，等于没定植。换成一个没人占的
  `WM_APP` 槽（`0x82F8`）才算真入队。（读数：`.scratch/wisp/probes/33/r8/logs/red-2.txt`）
- 第二稿挂死：`t33r8PumpExit` 的第一圈卡在 `<-done`，`go test -timeout 35s` 的栈 dump 显示
  泵线程在 `sta_windows.go:78`（`GetMessageW`）睡、驱动协程在 `<-done` 睡。根因＝
  `ballWndProc` 按 `activeBall` 的 hwnd 路由（`ball_windows.go:540-547`），测试没立 `activeBall`
  ⇒ `WM_APP_TASK` 被 `DefWindowProc` 吃掉、任务闭包永远不跑。**这一发读数顺手证实了一件事**：
  `Ball.Close` 里 `activeBall.CompareAndSwap(b, nil)` 与 `b.hwnd = 0` 都发生在 `sta.quit()` 之前，
  所以**发货路径上、还池前被派掉的任务闭包必然走 `DefWindowProc` 那一支**＝不会拿一条已经收摊的
  Ball 去跑旧闭包（这就是"owner 侧泵到静"不需要额外 `releasing` 闸的理由；栈 dump 存
  `.scratch/wisp/probes/33/r8/logs/red-2-hangdump.txt`）。

---

## ② 关键判据：这条脏线程今天会不会真的毒到后来者

（待填：本腿要把"会不会毒到"量成读数，而不是推论。已定案的读数口径见 ①(b)：PeekMessage 只有 owner 能读。）

## ③ 修法

（待填）

## ④ 钉

（待填）

## ⑤ 门禁四数

（待填）

## ⑥ 整包名册

（待填）

## ⑦ 我可能写错的条目

（待填 — 第 100 轮前必须写满）

## ⑧ 判不动的地方

（待填 — 第 100 轮前必须写满）

## 交件判语

（待填）
