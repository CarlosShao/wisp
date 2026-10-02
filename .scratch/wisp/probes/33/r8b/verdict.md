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

## ③ 改法（逐枚 `文件:行`，行号一律**在 `49448d7f` 上现读**；HEAD 会漂，引之前先复核）

尺：`git diff --numstat 3ef6cf41 49448d7f -- internal` ⇒
`internal/ball/sta_release_windows_test.go 196/2`、`internal/ball/sta_windows.go 1/1`（`ball_windows.go` **不在表里＝零改动**）。

**产码（一枚，一行换一行）**
- `internal/ball/sta_windows.go:72`：`defer s.releaseCOM() // TEMP pre-fix shape for the winlive baseline, reverted immediately`
  → **`defer s.releaseThread()`**。
  净效果 `-66 字节`（`wc -c`：`ebe3bd57` 上 9,680 → 现在 9,614），语义是：`start()` 的三条出口
  （create 失败 / 泵 break / panic 往上走）现在都先跑完"拆自己那扇窗 → 把自己队列派空 → CoUninitialize"，
  **仍锁着线程时**跑完，之后才轮到 `:65` 的 `defer runtime.UnlockOSThread()`（LIFO）。
  同文件 `:113-116` 那段"收摊的 unwind 是 deferred releaseThread"从谎话变成事实，本腿一字未改它。
- ⛔ **零枚常量、零枚阈值改动**：`pmRemove`、`releasePumpCap=4096`、`observe` 那六枚名册、
  `internal/observe/thresholds.go` 一字节未动。
- `internal/ball/ball_windows.go:964`（`b.sta.forgetWindow()`）**原样保留**：它是死腿那 5 行里唯一在生效的一处，
  而本腿 ④ 的 M3b 证明它此前**无人看守**（摘掉全绿）——所以现在给它补了牙，而不是删它。

**判据件（`internal/ball/sta_release_windows_test.go`，逐枚新增点）**
| 行 | 新增的是什么 | 为什么要本腿来加 |
|---|---|---|
| `:26`、`:31` | import `bytes`／`log/slog`／`strings`／`sync` | 下面两枚牙都要捕获产码在收摊时说的话 |
| `:343-348` | `run` 结构加字段 `releaseLogged string` | 见下 `:484` 那枚 |
| `:412-426` | 把 `s.start(...)` 整段包在一枚只收 ERROR 的 slog 捕获里，`start()` 返回即还原 | `start()` 不返回到泵退出＋deferred `releaseThread` 跑完为止，所以这段窗口恰好罩住整个 unwind |
| `:452` | 那行 `t.Logf` 末尾多打 `what the release unwind logged=%q` | 读数要进日志，不许只在断言里存在 |
| `:484-487` | **牙 #1**：`if r.releaseLogged != "" { t.Errorf(...) }` | M3b 实测：只摘 `ball_windows.go:964` 那枚 `forgetWindow()`，本文件**其余每一枚断言全绿**（窗真没了、队列真空、`staThread.hwnd` 最终也归零，因为 `releaseThread` 自己会清）⇒ 那 5 行是无人看守的承重墙。补这条之后：Close 拆窗又不记账 ⇒ 收摊那步会对一枚已死的句柄伸手 ⇒ `DestroyWindow` 返回 0 ⇒ 产码自己喊出 `could not destroy its own window` ⇒ 当场红。⛔ 这条不改产码行为，只把已有的日志升成判据 |
| `:495-515` | `t33r8LockedWriter`（带锁的 `bytes.Buffer`） | `-race` 下不许把"捕获日志"变成数据竞争 |
| `:517-595` | `t33r8PumpLeg` 结构＋`t33r8RunPumpLeg(plant int)`：在**本测试自己锁住**的线程上种 N 枚**线程级**消息、只调一次产码 `releaseOwnQueueToQuiet()`、记下返回值／队列 head／它说了什么，然后自己收尾 | 上限那一格要么有读数，要么就是"没人验过的 4096" |
| `:598-654` | **牙 #2**：`TestReleasePumpCapIsALoudReadingNotAGreen`（两腿：超上限＋控制腿） | 见 ④ 末两小节的读数 |

⛔ 未 `t.Skip` 任何用例、⛔ 未放宽任何断言、⛔ 未把任何红降成 `t.Logf`；② 那两枚红转绿是**接上产码**的结果，不是改判据的结果（改的只有 `:72` 那一行，见上表）。

---

## ④ 核心那一问：`releaseThread` 到底有没有把"带着别人状态的线程还池"排掉

**改前（必红＝② 那一发，`ebe3bd57` 未改动态）**：顶层 52 PASS／3 FAIL，两枚红就是这两枚钉。
**改后（必绿，在 `49448d7f`）**：
`go test ./internal/ball/ -count=1 -v -run 'TestSTARelease|TestReleasePumpCap'` ⇒ **rc=0**、3 枚全 PASS
（日志 `.scratch/wisp/probes/33/r8b/logs/fix-targeted-with-lognail.txt`）。关键两行逐字：

```
33-r8 create-error door: the instant start() returned: ball windows=0 all windows=0 classes=[] queue head=empty (start err=plant: renderer init failed plant err=<nil>)
33-r8 pump-exit door=quit-only: ... AT RELEASE: ball windows=0 all=0 classes=[] head=empty staThread.hwnd=0x0 | plant PostMessage rc=1 | this file's own pump moved 0 | after own teardown: ball windows=0 head=empty | what the release unwind logged=""
```

⇒ 同一枚线程，**收摊那一刻**窗数 0、队列报空、记录归零、无抱怨。
⚠ `all windows=0`（连 `classes=[]`）顺带量到一形：拆掉球窗之后那枚 Go 运行时的 `IME` 辅助窗也不在这条线程上了 ⇒
"还池时裸数窗＝0"这一维在这一形下真做得到（② 里我提过"全机数不得、只能数本包类名"，两件事不冲突：断言仍只问 `WispBallWindow`）。

**六发突变（逐发还原＋md5 核；载具 `.scratch/wisp/probes/33/r8b/mutate.py`，逐字日志 `logs/mutation-<名>.txt`）**
尺＝`go test ./internal/ball/ -count=1 -v -run 'TestSTARelease|TestReleasePumpCap'`，⛔ 不带 tag（默认档，秒级）。
表里 `FAIL/PASS` 是**顶层** `--- FAIL`／`--- PASS` 枚数；六发跑完 `git status --porcelain internal` 只剩本腿两枚在写文件，
`MUT-M` 残留 `grep -c` ＝ **0/0/0**（三枚文件各自），`final md5` 两枚产码与起手逐字相同。

| 突变 | 拔掉的是哪一步 | rc | FAIL 枚数 | 指名红用例 | 逐字红句（截要害） |
|---|---|---|---|---|---|
| **M0** `defer s.releaseCOM()` | 整条 duty 不接（＝`ebe3bd57` 出货形状） | 1 | 2 | 两枚门钉 | 门 1：`the instant start() returned: ball windows=1 ... queue head=msg 0x82F8 on hwnd 0x35C06DA`；门 2：`door=quit-only: ... still owning 1 live WispBallWindow window(s)` |
| **M1** `if false && hwnd != 0` | 第 1 步（拆窗） | 1 | 2 | 两枚门钉 | 门 1：`ball windows=1`（而 `queue head=empty`——泵仍把消息派掉了，只是窗没拆）；门 2 `quit-only`：`still owning 1 live WispBallWindow` |
| **M2** 摘 `releaseOwnQueueToQuiet()` | 第 2 步（派空自己队列） | 1 | **1** | 只 `TestSTAReleaseAfterFailedCreateHandsBackNoWindow` | `the instant start() returned: ball windows=0 ... queue head=msg 0x82F9 on hwnd 0x0` ⇒ **窗口级那枚植物被 `DestroyWindow` 一起带走了，剩下的是线程级的**：这一发独立复认了死腿注释里 M2 的归因（"只有线程级植物咬得住泵那一半"），也复认了 33-r7 那条"拆窗会顺手清掉写给该窗的消息" |
| **M3** 摘 `s.hwnd = 0` | 第 1 步里的"同时销账" | 1 | **1** | 只 `TestSTAReleaseAfterPumpExitDispatchesItsQueue` | `door=quit-only: ... with staThread.hwnd still holding 0x13C00808` |
| **M3b** `if false { b.sta.forgetWindow() }` | `Ball.Close` 那一侧的销账（`ball_windows.go:964`） | 1 | **1** | 只 `TestSTAReleaseAfterPumpExitDispatchesItsQueue` | `door=real-Close: the release unwind said something while unwinding ("... level=ERROR msg=\"ball: the ui-sta owner could not destroy its own window before releasing the thread; the handle stays alive in this process\" ...")` ⇒ **这一发在本腿补牙之前是 rc=0／全绿**（第一版网格逐字回显抄在 `logs/mutation-report-run1-before-tooth.txt`：
`M3b-close-does-not-forget rc=0 FAIL=0 PASS=3 restored_md5=yes build_failed=False | (no red)`；
⚠ 那一次的 `logs/mutation-M3b-*.txt` 与 `mutation-report.txt` 都被第二次运行覆写了，**只剩这份转录**，见 ⑦ 第 8 条）|
| **M4** `const pmRemove = 0` | 让派空变成"只看不拿" | 1 | 2 | 门 1 ＋ 上限钉 | 门 1：`queue head=msg 0x82F9 on hwnd 0x0`；上限钉控制腿：`planted under the cap, must count what it moved and report a quiet queue; it returned pumped=4096 head-after="msg 0x82F9 on hwnd 0x0"` ＋ `a pump that reached an empty queue must not warn; it logged "... hit its cap ..."` |

⇒ **判定：不是恒真判据。**六发各打红，且 M2/M3/M3b **各自只打红一枚**——分得开"哪一步被拔掉"，
不是那种"动哪都红"的钝尺；`door=real-Close` 那一支在全部改动里**一直是绿的**（真出货路径），
所以这套钉也不是"只测得出坏"。⛔ 没有任何一发靠放宽断言变绿。

**上限那一格（`releasePumpCap`）：报的是响亮读数，不是绿——在位的绿的读数**
`logs/capnail-targeted.txt` 逐字（`49448d7f`，rc=0）：

```
cap leg (planted 4160, want pump to stop at the cap): tid=30848 posted=4160 postErr=<nil> pumped=4096 head-after="msg 0x82F9 on hwnd 0x0" windows(ball/all)=0/0 classes=[] | own drain moved 64, head-final="empty" | logged="time=... level=WARN msg=\"ball: the ui-sta release pump hit its cap before the thread's queue was empty; the thread goes back to the Go pool non-empty\" cap=4096 pumped=4096"
control leg (planted 8, under the cap): tid=30848 posted=8 postErr=<nil> pumped=8 head-after="empty" windows(ball/all)=0/0 classes=[] | own drain moved 0, head-final="empty" | logged=""
```

三件事就此有读数：① 超上限时返回值**恰好等于** `releasePumpCap`（不是"它说空了"），且队列**仍然报得出消息**
⇒ 调用方不可能把它误读成"泵到静"；② 它**说在嘴上**（WARN 原文进判据，`strings.Contains` 两条都断）；
③ 同一枚函数在种得少时答 `pumped=8／empty／无日志` ⇒ 上限腿不是"永远答 4096"的瞎尺。
⚠ 这一枚钉问的是**产码那枚 helper 本身**（直接调 `releaseOwnQueueToQuiet()`），⛔ 它不钉"门上的调用者会不会忽略这个返回值"——
`releaseThread` 现在就是忽略返回值的那个调用者（它不做"还没空就别还池"的二次判定）。这一格归 ⑦ 第 3 条。
⛔ 没有为变绿改 `4096` 这一字节。

**两枚"顺带量到"的形状，具名免得下一程误归因**
- 门 2 的 `quit-only` 那一支在 ② 里 `head=empty` **不是仪器坏**：WM_QUIT 在队列取用优先级里排在窗口消息之后，
  外层泵 break 之前已经把窗口消息派完了。⇒ 队列那一半牙齿天然只咬得住 create-error 门（那里永不泵）。M2 一发把这句话说实了。
- ⛔ **一枚没复现的形状（不许读成"不会中"）**：`go test -tags=winlive` 整包在 **M0（不接线）**那一形下
  也是 `FAIL=3`（只有这两枚钉自己红，逐字 `logs/mutation-winline-M0-not-wired.txt`，顶层 67 PASS／3 FAIL／0 SKIP），
  **没有任何一枚别人的 winlive 用例被顺序依赖打红、整包日志里 `grep -c panic` ＝ 0**。
  ⇒ 本腿**没有**在这台机器上复现 86 条那一格"后面的用例第一次建窗就 panic"。
  所以这枚修法的凭据＝"钉住的不变式（本表 6 发）＋ 33-r7 那四形对照读数"，⛔ 不许写成"顺序依赖毒源已实测排除"。
  反过来说，"不接线也没别人红"也**不能**读成"不必接"：票 33 裁定（三）第 5 格原话那句"别把没中读成不会中"本腿照抄并自负。

**一枚正向的复用读数（这一维本腿真量到了）**：把收摊钉与真球 winlive 用例挤进同一进程、跑两遍，
`go test ./internal/ball/ -tags=winlive -count=2 -run 'TestSTARelease|TestReleasePumpCap|TestBallLive'`
⇒ **rc=0、顶层 18 PASS／0 FAIL**（9 枚 × 2 遍，逐名两两成对，`logs/pool-reuse-stress.txt`）。
同一进程里"收摊钉先跑、真建球窗的用例后跑"这个序**逐遍没坏**，两枚 `TestSTARelease…` 各自把线程还池之后，
`TestBallLiveLifecycle`／`TestBallLiveIdleBorderTransition` 这些要真建窗的用例照样建得起。
⚠ 主张只到这一层：**本腿没有量"后跑那枚是不是恰好落到刚被还回来的同一枚 M 上"**（那是调度器决定的，钉它会造出
一枚随 Go 版本抖的假门，见 ⑦ 第 7 条）。⇒ 这一发读的是"还回去的线程用得掉"，⛔ 不是"不还会坏在哪"。

---

## ⑤ winlive 名册（`-tags=winlive`，逐名，⛔ 仅本机可量）

命令逐字（`49448d7f`，PATH 带 sherpa＋build，`-timeout 900s`，日志 `logs/winelive-full.txt`）：
`go test ./internal/ball/ -tags=winlive -count=1 -v -timeout 900s`
⇒ **rc=1**、顶层 **PASS=69／FAIL=1／SKIP=0**（含子用例口径 75／1）、耗时 `11.346s`、`grep -c panic` ＝ **0**。
顶层枚数对账：默认档 56（② 基线）＋ winlive 新增 **14** ＝ **70** ⇒ 与 `69+1` 逐枚吻合，一枚不多一枚不少。

**winlive 那一族逐名（14 枚，全部 PASS）**：

| 用例 | 结果 | 用时 |
|---|---|---|
| `TestBallLiveLifecycle` | PASS | 1.91s |
| `TestBallLivePositionPersistence` | PASS | 0.44s |
| `TestBallLiveIdleBorderTransition` | PASS | 4.46s |
| `TestBallLiveEdgeDock` | PASS | 0.22s |
| `TestBallLiveEdgeDockHover` | PASS | 0.27s |
| `TestBallLiveAudioLiquidGate` | PASS | 1.11s |
| `TestLiveHotkeyRebindEndToEnd` | PASS | 0.64s |
| `TestLiveHotkeyOccupiedVsNotAttempted` | PASS | 0.04s |
| `TestLiveMuteHotkeyEndToEnd` | PASS | 0.05s |
| `TestLiveSleepingZeroTimerHandles` | PASS | 1.57s |
| `TestLiveClickSummonsAndDragDoesNot` | PASS | 0.04s |
| `TestLiveConfirmingCancelAndEscReturned` | PASS | 0.06s |
| `TestLiveNeverStealsFocus` | PASS | 0.20s |
| `TestLiveTransparentCornerFallsThrough` | PASS | 0.05s |

本腿那三枚在同一发里也各跑一次：`TestSTAReleaseAfterFailedCreateHandsBackNoWindow` PASS 0.01s、
`TestSTAReleaseAfterPumpExitDispatchesItsQueue` PASS 0.02s、`TestReleasePumpCapIsALoudReadingNotAGreen` PASS 0.03s。
唯一那枚 FAIL 还是 `TestC21TableColourRowsMatchTokensCSS`（② 第 3 行那枚常红，与本程无关）。

⚠ **口径与归属写死三句**：
1. 这 14 枚（含真球那 6 枚，会在本机桌面真起窗、真注册热键）是**本腿自己在 `49448d7f` 上跑的**，
   ⛔ 不是转述 `33-r8` 那句 "Winlive smoke passes on the real ball path"（那句在本腿之前盘上没有名册；
   现在有了名册，也不许再拿那句自述当凭据）。
2. `winlive` 在 CI 里**零岗位**（票 33 裁定 7 第四次同一形，本腿不重开）⇒ 本节全部读数＝**〔仅本机可量、CI 永看不见〕**，
   引用本节的任何表必须带这句。
3. ⛔ **题面那句"票 33 那 19 枚球／像素属〔仅本机可量〕"在本包射程内量不到**：本包 winlive 顶层枚数实测 **14**
   （尺＝`grep -hn "^func Test" internal/ball/{live_windows_test.go,live_guard_windows_test.go,hotkey_live_test.go,interaction_live_test.go}`，
   与跑出来的名册逐枚对得上）。那 19 枚若指 `cmd/wisp` 那一族，⛔ 本腿不许跑那枚包（`198-r1` 在飞，会洗掉它的读数）⇒ **归编排者量**，见 ⑦ 第 5 条。

---

## ⑥ 门禁读数（终态，`49448d7f`；每条带逐字命令与退码）

| 门 | 命令逐字 | 读数 |
|---|---|---|
| build | `go build ./internal/ball/` | **rc=0**（无输出） |
| vet | `go vet ./internal/ball/` | **rc=0**（`logs/vet-1.txt` 空）⚠ 见下方那条坑 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/ball/sta_windows.go internal/ball/ball_windows.go internal/ball/sta_release_windows_test.go` | **输出为空** ⇒ 三枚文件都已格式化（含未改的 `ball_windows.go`） |
| 整包（默认档） | `go test ./internal/ball/ -count=1 -v` | **rc=1**、顶层 **55 PASS／1 FAIL／0 SKIP**（含子用例 61／1）＝`logs/fix-fullpkg.txt`；那 1 枚 FAIL 是 ② 第 3 行那枚常红，⛔ 不是"本腿把它调剩下的" |
| 整包（winlive） | `go test ./internal/ball/ -tags=winlive -count=1 -v` | **rc=1**、顶层 **69／1／0** ＝ ⑤ 全表 |
| 定向单发 | `go test ./internal/ball/ -count=1 -v -run 'TestSTARelease\|TestReleasePumpCap'` | **rc=0**、3 枚 PASS |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**、末行逐字 `d22scan: clean - no D22 ban violations`；分母（**是"被扫文件枚数"，不是违规枚数**）：`bans #1-5 internal/=226`、`cmd/=36`、`ban #7 internal/tools/=23`、`ban #8 internal/=481`、`cmd/=86`、`ban #6/#8 frontend/=85`、`ban #8 design/=39`。全文 `logs/d22scan-final.txt`（同一发里 d22scan 自带测试 `ok tools/d22scan 19.085s`，顶层 PASS=34／FAIL=0） |
| staticcheck | **本腿未跑** | ⛔ 不出"干净的绿"：本机版（2025.1.1）与 CI 钉的（2026.2.1）不同版、解不开 go1.27 产物，票 33 裁定 8 已把这格定成〔未复认、交 CI〕。本腿不跑它拿假结论 ⇒ 这一格**没有人签过** |

**写面闸门（终态自证）**：`git status --porcelain internal` 在 `49448d7f` 之后＝**0 行**（本腿两枚文件已入库）；
`git show --name-only 49448d7f` ＝ **恰两枚路径**（`internal/ball/sta_windows.go`＋`internal/ball/sta_release_windows_test.go`）
⇒ 这才是"本腿提交了什么"的正确尺。⚠ 别用 `git diff --name-only 3ef6cf41 49448d7f` 问同一件事：
那一发现读**含别人的路径**（`docs/reports/pending-and-issues.md`、票 250／251 的票面、`probes/114|167|224` 若干枚），
因为两枚 commit 之间 HEAD 被别人推进过（本腿起手 `ebe3bd57` → 现在已漂到 `93f2b0c9` 之后）——
**共享工作树里"区间两端的差"含别人的账**，这条本腿差一步就写进结论，具名留着当下一程的教训。
⛔ `cmd/wisp` 里 `198-r1` 在飞的那几枚（`M cmd/wisp/run.go`、`?? cmd/wisp/firstrun*.go`）本腿一字未碰、未 `git add`、未跑其测试。
票面 AC 框现读（尺按 10-02 换形后的 `^[[:space:]]*- [ ]`）：**13 枚未勾／1 枚已勾**（09:3x 在 `93f2b0c9` 之后的工作树上现量，
与 10-01 12:22 编排者那次复尺同值），本腿一枚未动、也未新增格；
`git log --oneline -1 -- .scratch/wisp/issues/33-panel-host-c27.md` 不是本腿的提交（票面最后一枚动它的是别人的账）。

⚠ **本腿踩到并记下的一枚仪器坑**：`go vet ./... 2>&1 | head -20; echo "vet rc=$?"` 打出来的是 **`head` 的退码**，
不是 vet 的——那一发恰好也报 rc=0，但同一次 `go test` 是 build failed。⇒ 退码必须**不接管道**取，
或接管道后用 `${PIPESTATUS[0]}`。本表所有退码都是**不接管道**取的（日志先落盘再 `grep`）。

---

## ⑦ 判不动的地方（编号；每条写为什么＋谁能判）

1. **`TestC21TableColourRowsMatchTokensCSS` 这枚红本腿判不了也不该修**：射程在 `design/assets/tokens.css`，
   那是本腿的零读取区，红因是工作树缺该文件（票面 10-01 11:18 那节已记为已知常红，
   并明写那句 `the CSS leg of this check must never skip` 是它自己的规矩）。⇒ 谁能判＝编排者／界面侧那枚 agent。
2. **"顺序依赖毒源真的被排掉了"这一维，本包证不到端到端**：④ 末节那发负读数——winlive 整包在**不接线**那一形下
   也只有这两枚钉自己红、零枚别人的用例被连带、零 panic。⇒ 我能主张的只有"还池那一刻线程是干净的（6 发突变各打红）"
   ＋"还回去的线程下一条真建窗用例用得掉（`-count=2` 18/18）"。谁能判＝非实现者验收腿（`33-v2` 之后那一发）
   在安静机器上按 86 条那一格重跑；⛔ 本腿不许把它读成"已排除"。
3. **上限那一步"没泵到静就别还池"没有人钉**：`releaseThread` 现在**丢掉** `releaseOwnQueueToQuiet()` 的返回值
   （`sta_windows.go:159` 那句是裸调用），也就是说真撞了 4096 上限时，线程**照样还池**，只是盘上多一条 WARN。
   本腿钉的是"那一步会响亮读数"（④ 那两腿），⛔ 没钉"撞了上限就该改变收摊行为"——
   改不改属**产品形状决策**（要不要把一条 WARN 升成拒绝还池？还池失败往哪儿报？），且 `4096` 这一枚预算值
   该不该是 4096 属 D32/SLO 那一档。⇒ 谁能判＝编排者（本腿一字节未动那两个数）。
4. **`PostTask` 的注释与代码相反，本腿未动它**：`sta_windows.go:240-247` 注释逐字
   `// Window not created yet (or gone): run inline as last resort so callers cannot deadlock on a dead thread.`，
   代码做的是 `delete(s.tasks, id)` 然后 **return**（闭包**没跑**）。它与本程只有一处关系：
   我在 ④ 里主张"收摊那步不可能复活一枚球任务闭包"，凭的是 `DestroyWindow` 会带走写给该窗的消息
   （M2 那一发的 `head` 由 `0x82F8 on hwnd` 变成只剩 `0x82F9 on hwnd 0x0`，实测）。
   ⛔ 本腿不修那枚注释（不在射程、且改它要连带改还池语义）。谁能判＝票 07／64 那条线的主人或编排者。
5. **题面"19 枚球／像素"本腿量不到**（⑤ 第 3 句）：本包 winlive 顶层枚数＝14。谁能判＝编排者
   （要么在 `cmd/wisp` 那一族量，要么把那行改成本包实数 14；⛔ 本腿不跑 `cmd/wisp`）。
6. **staticcheck 这一格今天无人签**：本机版不同版、CI 钉版自陈导入期即崩（票 33 裁定 8）。谁能判＝CI。
7. **还池后 Go 是否真的复用同一枚 OS 线程：不属可断言的维度**。④ 那条 `-count=2` 读数只支持"复用没坏"，
   而"落到同一枚 M 上"是调度器决定，钉它会造出一枚随 Go 版本抖的假门。⇒ 谁能判＝没人（本腿主张：**不该钉**，
   要真复现只能按 86 条那一格的整包序去撞，见第 2 条）。
8. **"补牙之前 M3b 不红"那一发不可再生**：`mutate.py` 用固定输出名，第二次运行覆写了
   `logs/mutation-report.txt` 与 `logs/mutation-M3b-close-does-not-forget.txt`，所以"补牙前 `(no red)`"
   只剩本腿从终端回显逐字抄的 `logs/mutation-report-run1-before-tooth.txt` 一份（那份已注明自己的来源与射程）。
   要独立复现它＝把 `sta_release_windows_test.go:484` 那条断言注掉再跑 M3b——本腿**没跑**那一发
   （跑了就得再动一次写面文件，且它证的是"我的牙之前不存在"，不是任何产品事实）。
   ⇒ 谁能判＝验收腿愿意自己注掉再看一次；本腿主张的强度到此为止（登记在此，⛔ 不当已复认）。
9. **票 33 那 14 格 AC 与本腿的关系**：本腿只做了裁定（三）第 5 格那一件事，
   现读票面 AC 框 `^[[:space:]]*- \[ \]` ＝ **13**、`^[[:space:]]*- \[x\]` ＝ **1**（与 10-01 12:22 那次复尺同值，一枚未动）。
   ⇒ 谁能判＝`33-v2`／编排者。⛔ 本腿不据此主张票 33 任何一格"完成"。
