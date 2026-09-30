# 228-r1 现场笔记（实现腿：把球装进常驻那条腿）

工单：`.scratch/wisp/issues/228-ball-and-tray-are-not-in-the-process-that-runs-tasks.md`
射程（编排者 09-30 派单写死）：**只碰 AC#1（宿主存在性有生产点数）与 AC#7（三处过期文案）**。
⛔ 不动 AC#2／AC#3／AC#4／AC#5，不设 `runSpec.replyVeto`，不把 `NewChannels()` 换成 `DefaultChannels()`，
不碰工单勾选框，不写台账。

## 一、要做什么（做完的形状）

1. **AC#1**：让无参数那条腿（`cmd/wisp/resident_windows.go` → `runResident`）自己把球起来：
   分层窗口＋托盘＋四枚全局热键＋`ui-sta` 线程，全走 `internal/ball` 现成 API，**零新增模块依赖**。
   - 落点：新建 `cmd/wisp/resident_ball_windows.go`（`//go:build windows`）；
     `cmd/balldebug` 只当参考实现，**没改它一行**。
   - 起球失败**不致命**：响亮落日志＋打印"这个进程没有球＋原因"，事件循环照跑。
   - 收口：本轮**零新起协程**（`ui-sta` 由 `ball.New` 内部以在册名 `observe.Registry.Spawn`）；
     `rb.stop()`（＝`Ball.Close()`，内部 join ui-sta）作为**最后注册的 defer**，LIFO 先于
     `rt.Shutdown(false)` 跑 ⇒ D38(e) 第 2 步（`internal/proc/shutdown.go:18`
     "hotkey + wake-word listening stops"）由持有热键的宿主做，早于第 9 步 Job 关闭。
2. **AC#7**：三处打给人看／写进日志的过期文案改成带条件的事实句（逐字见 §五）。

### 本轮球的行为边界（写清楚，免得被读成"已经能干活"）

- 有：窗口真在桌面（winlive 按 pid 认领到 `WispBallWindow`，hwnd=54396232）、可拖拽／贴边、
  托盘图标＋四枚菜单项在、四枚热键真注册（`Ctrl+Alt+Q/M/P`＋`Esc`，来自 `ball.DefaultHotkeys()`）、
  状态＝`Sleeping`（真话：这条腿没事干）。
- 没有：⛔ 本轮**不把任何手势推进 D43 状态机**。这条腿今天没有任务通路、没有麦克风、没有审批门、
  没有面板宿主，所以十枚 `Events` 回调一律**具名落日志＋打印**，说清"手势到了、这条腿没有执行者"。
  把点击做成 `Listening` ＝让进程声称一条它没有的聆听路径，与编排者作废 `replyVeto`（账 `A472`）同一形状。

## 二、起手名册（2026-09-30 15:50 现跑；终态判据＝"红名册等于本名册"，不是全绿）

- `git rev-parse --short HEAD` = `3ee9b3ff`
- `git status --porcelain -- internal cmd` = **空**
- `git status --porcelain .scratch/wisp/issues/` = ` M .../228-...md`（编排者自己改的票面，我不动任何勾选框）
- AC#1 尺（改前）：`grep -rln "wisp/internal/ball" --include=*.go . | grep -v .scratch`
  → **1 行**：`./cmd/balldebug/main.go`（cmd/wisp＝**0**）
- 桌面三连：`balldebug.exe` 0 枚 / `wisp.exe` 0 枚 / `gh run list --limit 3` **网络不通**
  （`Get "https://api.github.com/...": EOF`）⇒ 读不到 slo-full ⇒ 本轮**零 CPU／句柄类测量**，只做功能改动。
- 全量基线红名册（`PATH=... go test ./cmd/wisp ./internal/... -count=1`，原始输出
  `.scratch/wisp/probes/228/r1/baseline-test.txt`，15:50）＝**4 枚红**：
  1. `cmd/wisp` `TestTicket223HandEditedFsLooseningCostsAnL2Card`（全仓并跑下的争用型，197-v1 曾整包 0 红）
  2. `internal/ball` `TestC21TableColourRowsMatchTokensCSS`（别人删了 `design/assets/**`，不在我地界）
  3. `internal/panel` 四枚：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／
     `TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／
     `TestC21DesignTokensFourwayAgree`（同上，`design/**` 禁读禁修）
  4. `internal/risk` `TestResolvePerCallBudget`（争用型假红，A432／票 226 AC#6 已登记）

## 三、做完了什么（三发 commit）

| commit | 碰的文件 |
|---|---|
| `90cf65b2` probes 骨架件 | `.scratch/wisp/probes/228/r1/notes.md` |
| `1739451d` feat+test AC#1 | `cmd/wisp/resident_ball_windows.go`(新) `cmd/wisp/resident_windows.go` `cmd/wisp/resident_ball_228_test.go`(新) `cmd/wisp/resident_ball_228_windows_test.go`(新) `cmd/wisp/resident_ball_live_228_windows_test.go`(新) |
| `15d2b2cd` docs+fix AC#7 | `cmd/wisp/main.go` `internal/proc/boot_windows.go`（⚠ 具名：**动了 `internal/proc` 这一处**，只改字符串＋属性，零新增依赖边） |
| 本发（polish） | `cmd/wisp/resident_windows.go`（boot 期收到退出请求时不再谎称"loop running"）＋ `cmd/wisp/resident_ball_windows.go`（两处注释/属性措辞）＋ 本 notes |

## 四、门禁读数（每一发命令与终态，逐名）

- 尺（改后）：`grep -rln "wisp/internal/ball" --include=*.go . | grep -v .scratch`
  → `./cmd/balldebug/main.go`、`./cmd/wisp/resident_ball_windows.go`（**生产**，就是 AC#1 要的那枚）、
  `./cmd/wisp/resident_ball_228_test.go`（尺本身引用的是字符串字面量，⛔ 不算宿主）。
  ⇒ **cmd/wisp 那条真会跑任务的腿：0 → 1**。
- `gofmt -l cmd internal` ＝ 空；`go vet ./cmd/wisp`、`go vet -tags winlive ./cmd/wisp`、`go build ./...` 全 rc=0。
- `bash scripts/d22scan.sh`（含它自己的正控 34 包 PASS=34/FAIL=0）⇒ 扫描 **clean**，
  ban #8 实扫 `cmd/` 69 枚、`internal/` 473 枚 Go 文件（我新增的四枚在里面，字符串零禁字）。
- 定向（真桌面，16:2x，`-tags winlive`）：`go test ./cmd/wisp -run 'TestAC228|TestAC1ResidentLeg|TestLive228' -tags winlive -count=1`
  → 8 枚全 PASS：`TestAC228ResidentLegIsTheBallHost`／`TestAC228BallHostAnswersEveryGesture`／
  `TestAC228ResidentLegReportsAndBooksItsBall`／`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`／
  `TestLive228ResidentLegOwnsABallWindowOnTheDesktop`／票 127 那三枚 `TestAC1ResidentLeg*`。
- **跑的哪一档（具名）**：`-tags winlive` 只有 `TestLive228...` 一枚（真 `RegisterHotKey`＋真窗口，
  本机亲跑，读数 hwnd=54396232、hotkeys live 4/4）；其余全部是默认档（无 tag）。
  ⚠ 本文件新写的 live 用例与 `internal/ball` 的 winlive 守卫（`requireQuietBallDesktop`，
  它对不是自己创建的 `WispBallWindow` 判红）**不能同屏并跑** ⇒ 只 `go test -tags winlive ./cmd/wisp -run TestLive228` 定向跑，
  头部注释已把这条规则写进文件。
- 全量终态（`PATH=... go test ./cmd/wisp ./internal/... -count=1`，16:1x，
  原始输出 `.scratch/wisp/probes/228/r1/final-test.txt`）＝**3 枚红**：
  `internal/ball` 1 枚、`internal/panel` 4 枚、`internal/risk` 1 枚 —— **全部在起手名册内，零新增红**；
  起手名册里那枚 `cmd/wisp TestTicket223HandEditedFsLooseningCostsAnL2Card` 这一发是绿的（争用型，两次读数不同因）。
  ⇒ 判据"红名册等于起手名册"以"无新增红名"成立（少一枚是争用型读数差，已具名）。
  ⚠ 这发全量跑在 polish 那一小改**之前**；polish 之后另跑了整包 `./cmd/wisp`（`final-cmdwisp.txt`）核同一件事。
- 本机 `staticcheck` 跑不了（缓存版本与 go1.27.1 的 export data 版本不兼容：
  `cannot decode "internal/cpu", export data version 4 > 2`），⛔ 没为它放宽任何东西；CI 用 `@latest`。

### 突变正控（四发，全部实跑，不是推理）

| 编号 | 种下去的改动 | 读数 |
|---|---|---|
| M1 | `rb := startResidentBall(...)` 换成 `rb := &residentBall{}` | `TestAC228ResidentLegIsTheBallHost` **红**："runResident calls none of the ball hosts" |
| M2 | 删掉 `defer rb.stop()` 一行 | 同枚 **红**："defers no ball teardown … D38(e) step 2" |
| M3 | 删掉 `OnDragEnd:` 一枚回调 | `TestAC228BallHostAnswersEveryGesture` **红**："leaves 1 of 10 gesture callbacks unset: [OnDragEnd]" |
| M4 | 注掉 `signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)` | `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` **红**（exit 0xc000013a、日志零条）**且**票 127 的 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` **红**（同一 0xc000013a）|

M1-M4 全部逐发还原（`cp` 回临时副本，⛔ 未用 checkout／stash），还原后 `git status --porcelain -- cmd` 只剩当时未提交的改动。

## 五、AC#7 三处改后逐字新句

1. `cmd/wisp/resident_windows.go`（`fmt.Printf`，运行时；随 AC#1 那发落）：
   `wisp: empty event loop running (no task is taken by this loop yet); <球的事实句> (Ctrl+C exits cleanly)`
   其中 `<球的事实句>` 是 `residentBall.statusLine()` 的两种取值之一（同一语句写死，不可能与真窗口漂移）：
   - `the floating ball window is up in this process (tray icon added, hotkeys live 4/4: summon=Ctrl+Alt+Q live, mute=…, cancel=…, panel=…)`
   - `this process has NO floating ball window: <Win32 返回的错>`
   ⚠ 保留 `empty event loop running` 子串＝票 127 那枚尺的字面量（`residentReachedLoop`），不是笔误。
   ⚠ polish：boot 期间就已收到退出请求时，改印 `wisp: exit request (interrupt) arrived during boot; the event loop was never entered; <球的事实句>`，
   不再先说"loop running"再不进循环。
2. `cmd/wisp/main.go`（`const usage` 原始串里，`wisp -h` 用户正文）：
   ```
     wisp             GUI resident process: boots the runtime skeleton, hosts the
                      floating ball window, its tray icon and its four global hot
                      keys in this same process, then parks in an event loop that
                      takes no task yet (Ctrl+C stops it; the tray's Exit item has
                      no stop path attached to it, and a ball that cannot be
                      created is reported and does not stop the boot)
   ```
   ⚠ 行首 `  wisp   ` 形状不动＝`leg_dispatch_gate_133_test.go` 的 `usageBare133` 按它核账。
3. `internal/proc/boot_windows.go`（`slog.Info`，运行时日志）：
   `activation requested by second launch` ＋属性 `handled_here="booked and cleared only"`、
   `bring_to_front="no path from this loop to a window; the ball, when a process has one, is its host's"`。

三处都无"要等某票"的将来式；⛔ 未新增任何扫注释票号的词面型仪器（派单尺 `grep "ticket 07"` 仍是派单用的尺）。
⛔ 没动的同类残留：`cmd/wisp/main.go:8` 包文档注释仍写着 "The floating ball GUI is ticket 07."——
它是**注释**，按票 228 AC#7 射程更正那一节明写"逐处判历史出处 vs 过期承诺不塞进本格，另立票 243"，我照令留着。

## 六、没做完的（非空；含两枚要编排者裁的）

1. **AC#2／AC#3／AC#4／AC#5 按令一格未动**：托盘四枚按钮没有一枚通向答复；`bookWaitingState` 没碰球；
   `rt.close()` 那两枚 join（`replyHandle`＋`reloadHandle`）没做；面板侧 L2「允许」的判红尺没写。
2. **托盘「退出」今天没有执行者**（本腿新暴露的一格，按读数说清利害，不夸大）：
   现在出厂的 `build/wisp.exe` 经 `objdump -p` 读 PE 头＝**Subsystem 00000003 (Windows CUI)**，
   不是 windowsgui ⇒ 双击也会带一个控制台，Ctrl+C 这条路**今天仍在**（编排者 09-30 16:1x 为此新立**票 244**，
   因为 `SPEC-11:50` 要求"无参数＝GUI 子系统"而 `docs/BUILD.md:87` 把那一切"推迟给票 07"、票 07 结案时没带走）。
   ⇒ **后果**：本腿让球与托盘进了常驻进程，但托盘那枚「退出」只落一条
   `tray exit requested … outcome=ignored` 日志；**等票 244 翻成 windowsgui 的那一天，双击启动就只剩这一枚按钮能停它**，
   这一格从"不便"升级成"没有出口"。要给它执行者只有两支，两支都超出"只碰 AC#1／AC#7"：
   - 甲：给 `internal/proc` 加一枚"请求停机"钩子（`Runtime` 多一个 channel／`Shutdown` 多一个 hook 位），
     常驻腿按 D38(e) 第 2 步的位置用它；代价＝动 `internal/proc` 的对外形状，且第 2 步的归属要说清。
   - 乙：在常驻进程里多挂一枚**不在 D38b 六枚名册里**的等待协程（`observe.Registry.Spawn`，名字会进
     `RosterReport.Unknown` ⇒ 每次开机一条 "goroutine outside the D38 roster (leak symptom)" WARN；
     现例证：`boot_windows.go:108` 只把 `ResidentOverBaseline` 当故障，`Unknown` 在产品线零消费者）。
     代价＝常驻名册要么加一枚、要么永久带一条假泄漏 WARN。
   ⛔ 我按「未定义即停」两支都没选；建议与票 244 一起裁（同一枚"GUI 那条腿怎么停"的问题）。
3. **`[hotkey]`／`[ball]` 配置没读**：常驻腿今天整个 `config.toml` 都没接（`cmd/wisp/config_reload.go` 属 `run` 那条腿），
   所以热键＝`ball.DefaultHotkeys()` 编译默认、球尺寸＝56px 编译默认。票 64 那座桥（`ball.NewHotkeyReloader`）
   在 `cmd/balldebug -config` 里有现成形状，接过来是另一格。
   ⚠ 顺带一条事实供裁定用：默认 `cancel=Esc` 走的是**全局** `RegisterHotKey`
   （`internal/ball/hotkey_windows.go:377 registerAll`，本机实测读数 `cancel=Esc live`），
   所以球活着的时候整台机器的 Esc 归它——这是 `DefaultHotkeys()` 的既有设计
   （票 64 的口径「宁可看得见也不能哑」），不是本腿新造，但**今天第一次跑在常驻进程里**。
4. **球位置不持久**：`Options.Store`＝nil；落盘路径没有任何规格/PLAN 依据（`grep` SPEC-08/PLAN 无命名），
   不自己发明。拖完重启会回原位，日志里那句 `ball moved: …not remembered` 说了这件事。
5. **视觉模式没翻**：`ball.EnablePrototypeVisuals` 没调用 ⇒ 常驻腿显示的是 `SPEC-08 §2.1` 冻结那一套（小球＋状态环），
   不是票 62 那套 liquid 原型（`balldebug` 默认开原型）。翻不翻是票 62/143 的口味签收，不归本腿。
6. `ball.New` 失败时**已创建的窗口/托盘不会被销毁**（`internal/ball/ball_windows.go:182 createOnSTA` 中途 return error，
   而 `Close()` 只认已 `activeBall.Store` 的成功路径）——本腿只是不致命地报出来；这枚泄漏在 `internal/ball` 里，
   ⛔ 没动，具名留给下一位。
7. 全量那发跑在 polish 之前（polish 只改一条打印分支＋两处注释），整包 `./cmd/wisp` 在 polish 后单独复跑过；
   **未**再跑第二遍全仓 `./internal/...`（零 `internal/` 文件在 polish 里被改）。
