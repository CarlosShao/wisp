# 33-r6 面板宿主 C27 — 顺序依赖 WM_QUIT 毒线程修复（写码腿）

票 33 `.scratch/wisp/issues/33-panel-host-c27.md`。本腿只修 HEAD 上那一枚**真红**：
`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（整包序红、隔离单跑绿）。

## 起手锚点（同发取，`date` + `git log` + `rev-parse` + `status --porcelain`）

- 时刻：`2026-10-01 14:22 +08`
- HEAD：`05f91799e1c88dfa2f96580e82e5cd9f5383d7a7`
  （标题＝`ledger(A503) + 撤票249: ... 抓到 1 枚真红（顺序依赖的 WM_QUIT 毒线程）...`）
- 分支：`dev`
- 起手名册：`.scratch/wisp/probes/33/r6/status-start.txt`（**242 行**，本腿不写"必须为空"，闸门＝终态名册等于起手名册＋只我这几枚路径）
- 台件目录：`.scratch/wisp/probes/33/r6/`

## 判据进度表

| 格 | 内容 | 状态 |
|---|---|---|
| ① | 最小复现（哪一枚先跑导致、组跑逐字读数） | 已填 |
| ② | 定性：产品缺陷 vs 测试卫生缺陷（两边说法） | 已填 |
| ③ | 修法 + 反控红句 | 已填 |
| ④ | 整包逐名红册（PASS/FAIL/SKIP 三数 + 逐名红） | 进行中（fullpack.log） |
| ⑤ | 门禁四数（build / vet / d22scan / gofumpt） | 待填 |
| ⑤⑥⑦ | 我可能写错的条目 / 我跑了哪些尺 / 判不动的地方（甲/乙/不做＋现量） | 已填（骨架轮） |


---

## ① 最小复现（哪一枚先跑导致、逐字读数）

所有读数用同法编译的测试二进制 `.scratch/wisp/probes/33/r6/wisp.test`（DLL 走
`third_party/sherpa-onnx` + `build`），HEAD `05f91799`（起手）。逐时刻：

**对照（隔离单跑）——本腿未复跑编排者的三发，自己复跑一发同形：**
```
14:24:22  AC13 alone, count=1:  --- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (1.47s)
```

**组 A（一枚 harness 窗 + AC13，count=3）：两枚全绿——单枚 harness 不足以致红。**
```
14:24:34-41  RealWindowHop PASS ×3 ; AC13 PASS ×3 (1.32 / 0.71 / 0.80s)
```

**组 B（两枚 harness 窗 + AC13，count=5）——AC13 命中目标红，逐字：**
```
14:25:39-26:21 (HEAD 05f91799)
  iter1: FAIL TestPanelHostRealWindowHopAndLifecycle (1.55s) ; PASS AC4Focus ; PASS AC13
  iter2: FAIL TestPanelHostRealWindowHopAndLifecycle (0.08s) ; PASS AC4Focus ; PASS AC13
  iter3: PASS RealWindowHop ; PASS AC4Focus ; --- FAIL: TestAC13... (15.01s)   <-- 目标红，签名逐字
  iter4: PASS RealWindowHop ; PASS AC4Focus ; --- FAIL: TestAC13... (15.01s)   <-- 同上
  iter5: 全 PASS
```
⇒ **先跑者是 `panel_host_windows_test.go` 的两枚 harness 真窗用例
（`TestPanelHostRealWindowHopAndLifecycle` / `TestAC4FocusReturnToPriorWindowGap33r5`）。**
它们各自 `mgr.Destroy()`（post `WM_CLOSE`）后，harness 协程泵到 `WM_DESTROY → wndproc
w.Terminate()（裸 PostQuitMessage）→ WM_QUIT` 落进自己线程队列，再 `runtime.UnlockOSThread()`
（:90）把这条**带 quit 的线程还给 Go 池**；AC13 的 `panel-sta` 若被调度器分到同一 M，
`bringUp → Embed` 的 `GetMessageW`（chromium.go:100-108）立刻取到 quit、`r==0` 早退、
`inited` 仍 0、`e.webview` 未赋值 → `Init`（:131）nil 解引用 panic → `showAndWait` 等满 15s → 红。

**为什么组跑是概率性的**（命中 2/5，另一次 1/6）：池里既有被毒的 M 也有干净的 M，AC13 落到哪个
由调度器定；整包几十枚真窗 ⇒ 被毒线程多、命中高，故"整包必红、隔离必绿"。`RealWindowHop`
自己在 iter1/2 也红＝毒在 count 之间跨轮传播（它下一轮起窗也落到上一轮的毒线程）。

**确定性最小样本（不靠调度器运气）——同一条锁死线程上把 quit 定植：**
临时探针（已删，`.scratch/wisp/probes/33/r6/` 不留痕）逐字读数 `14:32:59`：
```
A (poison PostQuitMessage, 同线程, 不排干, NewWithOptions):  panicked="runtime error: invalid memory address or nil pointer dereference"  nilWindow=true
B (同 poison, 先 PeekMessage(WM_QUIT,PM_REMOVE) 排干, 再建窗): panicked=<nil>  nilWindow=false
```
⇒ **一枚待建窗线程上"有没有 pending WM_QUIT"就是红/绿的唯一变量，且"建窗前把 quit 排干"确定性地转绿。**
这把整包那枚 race 钉死成了可复现的单线程事实。

## ② 定性：产品缺陷 vs 测试卫生缺陷（两边说法）

**触发者是测试卫生缺陷。** 唯一"泵过一扇 webview 窗又 `UnlockOSThread` 把线程还给池"的地方是
`panel_host_windows_test.go` 的 `hostThreadHarness`（:90）。产品侧 `residentPanel.loop`
锁而不解（`panel_resident_windows.go:202`，全程无 Unlock），线程随协程退出被销毁，它自己那条
WM_QUIT 一起没了——**出货拓扑今天不会复发这一枚红**。`notify_windows.go` 的 Unlock 落在 message-only
窗（非 go-webview2 的 wndproc），无 quit 源，安全。

**但根 Hazard 属产品鲁棒性缺口，且 shipping 面可达，故修法落产码。** `go-webview2` 每次销毁一扇
它托管的窗都**免费**往当时泵关闭消息的线程投一枚 WM_QUIT（`webview.go:242-243` + `:381-383`）。
`bringUp` 隐含假设"来建窗的线程队列是 quit-clean 的"，而这条假设对任何**复用**过窗口线程的调用方
都不成立。`PanelManager` 文档化的 recreate 面（`Destroy` 之后"later Show creates a fresh window"）
在同一条仍在跑的线程上正是"上一扇窗已 Destroy、下一扇 bringUp 同线程"＝同形 panic。
编排者那句"谁从错协程调一次 Terminate 就毒掉日后建窗的线程"方向对但机制更准：不是"错协程"，是
"任何泵过窗关闭的线程被复用"。⇒ 落产码 `bringUp`，既治整包这枚测试触发的红，也补 shipping 的隐患。

## ③ 修法 + 反控红句

**修法（最小面，产码 `cmd/wisp/panel_host_windows.go`）：**
`bringUp` 进 `webview2.NewWithOptions` **之前**调 `drainStaleQuitBeforeCreate()`——
`PeekMessageW(0x12,0x12, PM_REMOVE)` 循环只摘本线程队列里的 `WM_QUIT`（不碰其它消息、不碰别的线程、
只在建窗前），带 64 枚上限、命中上限即 `slog.Warn` 自陈。出货路径队列本就没有 quit，此调用是 no-op。
未碰依赖、未 `go mod`、未动 `internal/**`。

**钉（确定性）：`TestAC13BringUpSurvivesAReusedThreadQuit`（同包）**——在一条锁死的线程上
`PostQuitMessage(0)` 定植 quit，再跑产品自己的 `bringUp`。修在＝建起窗、`created=true`、无 panic；
反控＝把 `bringUp` 里那一行 `drainStaleQuitBeforeCreate()` 删掉。

**反控逐字读数（HEAD `e06e715e`，临时删该行后跑，`14:39` 段）：**
```
--- FAIL: TestAC13BringUpSurvivesAReusedThreadQuit (0.14s)
    bringUp on a thread carrying a pending WM_QUIT runtime error: invalid memory address or nil pointer dereference
    (the recovered panic) ... Init (chromium.go:131) dereferenced the nil control. Fix: bringUp must call drainStaleQuitBeforeCreate
```
即"红因是当场 panic（不是 15s 超时）"，钉有牙。反控后已 `git cat-file blob HEAD:` 还原该文件、
md5 `3a16c129…` 复确、`git diff --stat HEAD` 空。
带修正常态：`14:50` 同钉 count=6 全 `PASS ... created=true`；`TestAC13` 冷启动整包复跑见 ④。

---

## ⑤⑥⑦ 自我对抗三节（先写满）

### ⑤ 我可能写错的条目（逐条，附"如果错了后果"）

- (a) **"归因＝hostThreadHarness 的 `runtime.UnlockOSThread()` 把带 WM_QUIT 的线程还给池"**：
  这是从依赖源码（`webview.go:242-243` 的 `case WMDestroy: w.Terminate()` + `:381-383` 的裸
  `PostQuitMessage`）＋编排者假设推出来的**形状**，本腿尚未用最小样本复现成确定性红。
  若错（真凶是别的形，如共享 dataPath 抢用、并发 bringUp、`CoInitialize` 那一支），
  那我的修法（去 Unlock）治不了真红 ⇒ 判据④那枚具名红仍在。
- (b) **"出货路径不受此害"**：只核了 `residentPanel.loop` 锁而不解（`panel_resident_windows.go:202`，无 Unlock），
  且 `notify_windows.go` 的 Unlock 落在 message-only 窗（非 webview wndproc）。
  未逐行读 live 测试与其它起窗点 ⇒ 若还有别的"锁了又解且那线程起过 webview 窗"，出货面也可能中。
- (c) **"AC13 单跑三发全绿"**：引编排者 `.scratch/wisp/probes/orchestrator/33r5-ac13-alone.txt`，**本腿未复跑**。
- (d) **毒线程复用是调度器决定的、非硬保证**：组跑可能偶发不复用同一 M ⇒ "红"也可能变"绿"。
  我的钉必须能在"没复用到时"也响（否则反控不成立）。

### ⑥ 我跑了哪些尺 / 哪些没跑

- 已跑（起手）：`ls .scratch/wisp/probes/33/`（名册 r6 未占）、`git status --porcelain`（存 status-start.txt）、
  `grep UnlockOSThread cmd/wisp/*.go`、读依赖源码 `chromium.go`/`webview.go`、读三份产码/测试文件。
- 待跑：最小复现组跑（①）、修法后整包（④）、门禁四数（⑤）。
- ⛔ 不跑：`staticcheck`（本机版与 CI 钉版不同 ⇒ 标〔未复认〕，不拿假结论）。

### ⑦ 判不动的地方（逐条：甲＝按我读数、乙＝别的形、不做）

- 甲：**去 `hostThreadHarness` 的 `runtime.UnlockOSThread()`**（让窗口线程随协程退出被销毁，
  带走的 WM_QUIT 也一起没了）——与产品 `residentPanel.loop` 同一不变式，落在我写面（`cmd/wisp/**` 测试）。
- 乙（**若甲读数不成立才转**）：把不变式落进产码——但产码 `bringUp`/resident 本就锁而不解，
  真正"锁了又解"只在测试 harness ⇒ 若定性为产码面要写死，只能加注释＋一枚会响的钉，不能加 Unlock 逻辑。
- 不做：**不**碰依赖（`pkg/edge`、module cache）；**不** `go mod`；**不**搬 AC13 进 winlive/t.Skip；
  **不**放宽那 15s；**不**把断言降级成 Logf。
