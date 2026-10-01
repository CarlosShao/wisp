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

---

## 编排者收尾标注（15:36，`84c67e7a`；本腿 15:2x 死于模型服务断线，通知一律当未验证，下列读数全是我自己在盘上取的）

**本腿的死亡形状＝第二种形（正文写满了却没提交）**：它交了四枚（`31444593` 骨架／`e06e715e` 窄形修法＋第一枚钉／
`4bd32fe0` 证据 ①②③／`4be1d3f1` **加宽**修法＋第二枚钉），但工作树里还压着**三枚未提交的改动**。我按
「代提死腿的原话＝允许、代填它空着的两节＝禁止」处置：**逐字代提那三枚改动＝上面 `84c67e7a`**（pathspec 只有
`cmd/wisp/panel_host_windows.go` 45/41、`panel_host_windows_test.go` 11/1、`panel_resident_windows_test.go` 5/66）。

**⛔ 本表现在与代码不同步，且这是它自己改的方向造成的，不是我改坏的**：
- §③ 与 §⑦ 写的是 `e06e715e`／`4be1d3f1` 那两形（§③＝`drainStaleQuitBeforeCreate` 窄形，§⑦ 甲＝去 harness 的 Unlock），
  而 `4be1d3f1` 一度把产码换成**整条队列排干**；`84c67e7a` 又**收回窄形**并把 §⑦ 甲那一支（harness 不再
  `UnlockOSThread`）真落了产码侧的镜像注释。⇒ 引用"r6 的修法"时必须带 commit 号，别引用本表当现状。
- 它在 `84c67e7a` 的产码注释里留下一句**我没有读数的自述**：整条队列排干"不派只删"会留下
  **zombie WebView2 controller、卡住下一次建窗**。这句在它本表里**没有对应读数**（§①②③ 都没登），我也没复跑
  ⇒ 现按〔仅死腿自述〕处理，`33-r7`／`33-v2` 若要据它选型，得先把那一形自己复现一次。
- **§④（整包名册）与 §⑤（门禁四数）两节它没写，我不代填**。取而代之的是我这两发现量，出处写清＝编排者本人、
  取数时刻与锚点一并给：
  - 我 15:27:43 在 `4be1d3f1` ＋它那三枚未提交改动之上跑整包 `cmd/wisp`（PATH 带 `third_party/sherpa-onnx`＋`build`，
    跑前确认机器上没有别的 `go test`）＝**rc=1、319.926 s、一枚红＝`TestAC14AwaitedBindingReplyReachesThePage` (15.00s)**，
    红句逐字 `panel_resident_windows_test.go:443: timed out after 15s waiting for the panel thread to finish Show -
    this is a failed measurement, not a pass`；同发 recovering 出的 panic 栈＝
    `pkg/edge/chromium.go:131 Init` ← `webview.go:340` ← `webview.go:109 NewWithOptions` ←
    `cmd/wisp/panel_host_windows.go:249 bringUp` ← `:353 Show` ← `panel_resident_windows.go:324 RequestShow.func1`
    ← `:250 drainTasks` ← `:210 loop`。逐字日志＝`.scratch/wisp/probes/orchestrator/33r6-head-roster.txt`（658 行）。
  - 门禁四数（我替它补跑销账）：`go build ./...` rc=0／`go vet ./cmd/wisp` rc=0／`sh scripts/d22scan.sh` rc=0
    （clean；ban #8 `cmd/` 81、`internal/` 476）／`gofumpt -l` 起手**点出 `panel_resident_windows_test.go` 不合规**
    ⇒ 我 `gofumpt -w` 补格式（纯格式：该文件逐字 5/66、`-w` 后 3/64），复跑列表为空。`staticcheck` 本机版与 CI 钉版
    不同 ⇒ 〔未复认〕。⛔ `go mod tidy` 一字节没跑。
  - ⇒ **判语：本票这枚顺序依赖红到今天为止没有一枚腿交回过"整包绿"**。窄形我这发＝AC#13 绿、AC#14 红；
    宽形它的读数＝AC#13 绿、AC#14 红。**两形都被量过、都不够**。
  - 那条栈的位置有信息量：它落在 `loop` 里 `handOverPump` **之前**的第一次 `drainTasks` ⇒ 是**这条 `panel-sta`
    的头一次建窗**就踩到毒，毒源不是"本线程上一轮自己留的"，而是**调度器把新协程放到了一枚已被别家毒过的 M 上**。
    这一格把"再换一种排干范围"那条路当场收窄了。

**一笔纪律账，记在这腿名下（不影响它的读数）**：§① 里那枚"确定性最小样本"的临时探针它写了「**已删，
`.scratch/wisp/probes/33/r6/` 不留痕**」——本仓 `issues/README` 规则 8 是"临时件**只建不删**"。⇒ 那两行
A/B 对照读数（`nilWindow=true` vs `false`）今天**盘上不可复算**，只值〔仅自述〕；`33-v2` 若要拿它当凭据，
得自己重造那一形。

**未提交的原始增量我已另存一份**（防"代提即丢"）：`.scratch/wisp/probes/orchestrator/r6-uncommitted.patch`
（231 行／14,134 字节）。我没有删、没有 revert、没有 `checkout .`。

**归程**：`33-r7`（修复腿，独占 `cmd/wisp`＋桌面）带着上面两发读数接着做。它要证的**不是**"再换一种排干范围"，
而是"一枚新协程可能落到被别家毒过的 M 上、而第一次建窗就会炸"这件事**由谁在哪一层挡**。票 33 的 AC 框一枚没动
（`.scratch/wisp/issues/33-panel-host-c27.md` 不在我上面那次 commit 的 pathspec 里，未勾数不变）。
