# 票 33 · 接管腿 `33-r8b` · 判语件

工单：`.scratch/wisp/issues/33-panel-host-c27.md`（本腿射程＝编排者裁定（三）第 5 格那一发：
把 33-r7 在面板侧量清的"泵到静再还池"落进球的 `ui-sta` 收摊那一跳）。
本腿是**写码腿**，不是验收腿：⛔ 票面任何 AC 框一枚未勾、⛔ `docs/reports/pending-and-issues.md` 一字未写、
⛔ 未改名 `-done`。勾与判语归编排者与非实现者验收腿。

---

## ① 起手锚＋写面声明

| 尺 | 读数 | 取数时刻／命令（逐字） |
|---|---|---|
| 锚点 | `ebe3bd57`（分支 `dev`，父 `8ae4c23e`） | `git log -1` ＋ `git rev-parse --short=8 HEAD` 同发 09:19 |
| 日期 | `Fri Oct 2 09:19:05 CST 2026 +08` | `date` 与上一行同发 |
| 脏项 | `git status --porcelain internal cmd` ＝ **0 行** | 09:19:05 同发 |
| 运行期 PATH | `ls third_party/sherpa-onnx/*.dll \| wc -l` ＝ **3**（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`，**无 `bin/` 子目录**） | `ls third_party/sherpa-onnx/` 现读 |
| 派单口径 | 每次跑测试前置 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` | 与 ② 那发同批 |

**写面声明（本腿只写这些）**：`internal/ball/**` ＋ 本件与 `.scratch/wisp/probes/33/r8b/logs/**`。
⛔ 未写、未跑、未"顺手修"：`cmd/wisp/**`、`internal/config/**`（`198-r1` 在飞）、`scripts/**`（`250-r1` 在飞）、
`internal/panel/**`（另有主人）、`frontend/**`、`design/**`（后两枚**零读取**：本腿所有 `grep`/`git show`
一律显式给根 `internal`／`cmd`／`tools`／`docs`／`.scratch`，⛔ 无一枚以 `.` 为根）。
一字未动：`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、
`internal/observe/thresholds.go`、任何 golden、`tools/d22scan/allowlist.txt`、三枚冻结测试件。

**本腿起手时（`ebe3bd57` 上）现读的题目事实，逐条复尺——含推翻编排者题面两处**：

1. **行数断言成立**：`git show --numstat ebe3bd57` ＝ `sta_windows.go 106/0`、`ball_windows.go 5/0`、
   `sta_release_windows_test.go 460/0` ⇒ 产码 106+5 ＝ **111 行纯新增**、测试 **460 行**，与题面一致。
   `wc -l -c` 现量（同一枚 sha）：`sta_windows.go` **268 行／9,680 字节**、`ball_windows.go` **1,008／32,016**、
   `sta_release_windows_test.go` **460／21,469**。
2. **`ball_windows.go` 那 5 行的落点**：`b.sta.forgetWindow()` 逐字在 **`internal/ball/ball_windows.go:964`**，
   在 `func (b *Ball *Close)` 的 `if b.hwnd != 0 { … }` 块内（题面写 `:960 起` ＝ hunk 头 `@@ -957,6 +957,11 @@`
   之后第一枚新增行是 960，**注释四行 960–963 ＋ 那句调用 964**，两种指认都指向同一处，不算冲突）。
3. ⛔ **推翻题面第一句（要紧的那枚）**：`releaseThread` **在本仓零枚调用者**。尺（逐字、在 `ebe3bd57` 上现跑）：
   `grep -rn "releaseThread\|releaseOwnQueueToQuiet\|forgetWindow\|releaseCOM" internal/ tools/ cmd/`
   ⇒ `internal/ball/**` 命中里 `releaseThread` **只出现在它自己的定义 `sta_windows.go:147` 与四处注释**
   （`:119/:131/:146/:164`），**零枚调用点**；`releaseOwnQueueToQuiet` 只有定义 `:185` ＋ 被 `:159` 那一处（在未被调用的 `releaseThread` 体内）。
   唯一在位的调用是 `sta_windows.go:72`：

   ```go
   defer s.releaseCOM() // TEMP pre-fix shape for the winlive baseline, reverted immediately
   ```

   同文件 `:113-116` 那段注释却写着"收摊的 unwind 是 deferred releaseThread"。⇒ **那句 `reverted immediately`
   没有被执行**：产码停在改法之前的形状上，新落的那 106 行是**编译得过、永不执行**的一块死码。
   ⇒ 后果：票 33 裁定（三）第 5 格那件事**尚未发生**，而它留下的两枚新用例是**当场红的**（见 ②）。
   ⇒ 这条不是"未验证"那么轻：**未验证 ＋ 未接线**。⛔ 编排者题面那句"死腿留在库里的未验证半成品"要按这个强度重读。
4. ⛔ **推翻题面第二句**：题面"`go vet ./internal/ball/` rc=0 ⇒ 能编译"成立，但那枚 rc=0 **同时是死码的掩护**——
   Go 不因未使用的**方法**报错（只对未使用的局部变量与 import 报错），所以"vet 干净"这一格对"有没有接上"是**瞎的**。
   下一程要判这一族，尺只能是"调用点枚数"或整包颜色，不能是 vet。

---

## ② 基线名册（未改动态，`ebe3bd57`，⛔ 本腿当时一枚产码未改）

命令逐字（起于 09:19:5x，PATH 已带 sherpa＋build）：
`go test ./internal/ball/ -count=1 -v > .scratch/wisp/probes/33/r8b/logs/baseline-plain.txt 2>&1`
退码 **rc=1**，末行逐字 `FAIL	github.com/CarlosShao/wisp/internal/ball	0.147s`。

**四数（顶层口径 `^--- PASS`／`^--- FAIL`／`^--- SKIP`，与含子用例口径分列，⛔ 不混）**：

| 口径 | PASS | FAIL | SKIP | 合计顶层 |
|---|---|---|---|---|
| 顶层（`^--- `） | **52** | **3** | **0** | 55 |
| 含子用例（`--- ` 全量，`grep -c --`） | **58** | **3** | 0 | — |
| 仅缩进子用例（`^[[:space:]]+--- PASS`） | 6 | 0 | — | — |

**红名册逐名（3 枚）**：

| # | 用例名 | 归因（本腿判） |
|---|---|---|
| 1 | `TestSTAReleaseAfterFailedCreateHandsBackNoWindow` | **`ebe3bd57` 那枚死腿半成品自己造成**（准确说：它没接线）。红因逐字见下 |
| 2 | `TestSTAReleaseAfterPumpExitDispatchesItsQueue` | 同上，同一枚根因（`releaseThread` 未被调用） |
| 3 | `TestC21TableColourRowsMatchTokensCSS` | **非本程、非本包可控**：红因逐字 `tokens_table_test.go:1468: read design/assets/tokens.css: open …design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip`——`design/**` 不在这台机器的工作树里（票 33 收件（10-01 11:18）那节末条已把这枚常红记为已知：`git status` 里那 16 枚 ` D` 的 `design/**`）。⛔ 本腿一字不碰 `design/**`，也不许任何人把这枚红读成"内部尺坏了"（那句 `must never skip` 是它自己写的规矩） |

**两枚红用例的逐字读数（这是"钉子有牙"的证据，不是"钉坏了"的证据）**：

- 门 1（create-error）：`the instant start() returned: ball windows=1 all windows=2 classes=[WispBallWindow IME] queue head=msg 0x82F8 on hwnd 0x75B063A (start err=plant: renderer init failed plant err=<nil>)`
  ⇒ 两条 `t.Errorf` 各响一次：**还池时窗数 1**（要 0）、**队列 head 非空**（要 empty）。
  ⚠ 同发一条量到 `classes=[WispBallWindow IME]`、`all windows=2`：Go 运行时自己在那条线程上挂了一扇 `IME` 辅助窗，
  所以"窗数＝0"这一格**只能问本包那一枚类名**（用例正是这么写的），问全机会恒红。
- 门 2（pump-exit，`door=quit-only`）：`AT RELEASE: ball windows=1 all=2 classes=[WispBallWindow IME] head=empty staThread.hwnd=0x3FA0B0A | plant PostMessage rc=1 | this file's own pump moved 0`
  ⇒ 响的两条是 **窗数 1** 与 **`staThread.hwnd` 仍持句柄**；`head=empty` **这一条本来就没响**，
  原因本腿量清了：`quit()` 之后投的那枚 `0x82F8` 是一枚**窗口**消息，而 WM_QUIT 在队列里的取用优先级排在窗口消息**之后**，
  所以外层泵在 break 之前已经把它 `DispatchMessageW` 掉了（同发 `this file's own pump moved 0` ＝ 泵退出时确实已经空）。
  ⇒ **具名结论：队列那一半牙齿只咬得住 create-error 那一门**（那里永不泵）；
  这也正是死腿注释里那枚 M2 归因的独立复认（它把线程级投递 `PostThreadMessageW` 当作唯一能咬住的植物）。
- 门 2 的 `door=real-Close` 一支：**逐字全绿**（`AT RELEASE: ball windows=0 all=0 classes=[] head=empty staThread.hwnd=0x0`）
  ⇒ 这条是**正控形状已在位**：出货路径（`Ball.Close` 派 `DestroyWindow` ＋ `forgetWindow` ＋ `quit`）本来就干净，
  所以那枚钉不是恒红、也不是"只测得出坏"，它分得开 clean 与 dirty 两形。
  ⚠ 注意 `door=real-Close` 那发 `staThread.hwnd=0x0` 靠的正是 `ball_windows.go:964` 那枚 `forgetWindow()`——
  那 5 行**是唯一在生效的产码改动**，它不在死码里。

---

## ③ 改法（逐枚 `文件:行`）

**本枚 commit 时点的状态：产码零改动**（② 那一发之后本腿没碰过任何 `.go`）。
下一步要落的改法只有一枚，形状＝把 `:72` 那行 TEMP 换回它注释里承诺的形态：
`internal/ball/sta_windows.go:72` 的 `defer s.releaseCOM()` → `defer s.releaseThread()`（并删那句 TEMP 注释）。
落完后的逐枚行号、`gofumpt` 读数、以及"到底改了几字节"会在本节重写，⛔ 本节此刻不含任何引用别节的表格。

---

## ④ 核心那一问的突变读数

本节此刻**空着是刻意的**：本腿尚未落 ③，突变（M1–M4）要在改后那一发之上取，
在改前取只会重跑 ② 的两枚红。取数计划（每发都要求"改前必绿／拔掉必红"，且指名用例）：
拔掉 `releaseThread` 的 DestroyWindow 那一支 → 门 1 与门 2 的窗数断言必红；
拔掉 `releaseOwnQueueToQuiet()` 那一行 → 门 1 的队列断言必红；
拔掉 `s.hwnd = 0` 与 `forgetWindow()` → `staThread.hwnd` 那条必红；
`pmRemove` 改 `PM_NOREMOVE` → 上限路径必响。
读数与逐字红句在本节落表，⛔ 不落表不许在 ③ 里说"已验"。

---

## ⑤ winlive 名册

尚未跑。`-tags=winlive` 那一族在本包共 **14 枚用例**（尺＝起手现读，`grep -hn "^func Test" internal/ball/{live_windows_test.go,live_guard_windows_test.go,hotkey_live_test.go,interaction_live_test.go}`，
build 头逐字 `//go:build windows && winlive`）：`TestBallLiveLifecycle`／`TestBallLivePositionPersistence`／
`TestBallLiveIdleBorderTransition`／`TestBallLiveEdgeDock`／`TestBallLiveEdgeDockHover`／`TestBallLiveAudioLiquidGate`／
`TestLiveHotkeyRebindEndToEnd`／`TestLiveHotkeyOccupiedVsNotAttempted`／`TestLiveMuteHotkeyEndToEnd`／
`TestLiveSleepingZeroTimerHandles`／`TestLiveClickSummonsAndDragDoesNot`／`TestLiveConfirmingCancelAndEscReturned`／
`TestLiveNeverStealsFocus`／`TestLiveTransparentCornerFallsThrough`。
⚠ 这 14 枚会在本机桌面上真起窗／真注册热键——照常跑，⛔ 不因为"CI 不跑"就不一名不留。名册（逐名 PASS/FAIL/SKIP）在本节补。

---

## ⑥ 门禁读数

本节此刻只有一条（起手时点）：`go vet ./internal/ball/` 由编排者在 `ebe3bd57` 上跑过 rc=0，
**本腿不复用那一发当凭据**，终态自己重跑 `go build`／`go vet`／`gofumpt -l`／`d22scan` 并逐枚贴读数。
⚠ `staticcheck` 本机版解不开 go1.27 产物 ⇒ 即便跑也标〔未复认〕，⛔ 不许据此出"干净的绿"。

---

## ⑦ 判不动的地方

编号列在本节，此刻 2 条（都是起手就能具名的）：

1. **`TestC21TableColourRowsMatchTokensCSS` 这枚红本腿判不了也不该修**：射程在 `design/assets/tokens.css`，
   那是本腿的零读取区，且它的红因是工作树缺文件（票面已记为已知常红）。谁能判＝编排者／界面侧那枚 agent。
2. **"还池后下一条落到该 M 上的用例真的不再 panic"这一维，本包内证不出来**：② 只钉住"收摊那一刻线程是干净的"，
   而 86 条那格的完整形状是"单跑永远绿、整包才红"——要真复现顺序依赖，得在同一枚整包里让**后续**用例在那条线程上再建真窗。
   本包 14 枚 winlive 用例就是那个"后续"，但它们与默认档不共跑（`//go:build windows && winlive`），
   所以顺序依赖的端到端复现只有编排者那两发整包能判。谁能判＝非实现者验收腿（`33-v2` 之后那一发）。
