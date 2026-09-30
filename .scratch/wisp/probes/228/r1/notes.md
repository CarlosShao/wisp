# 228-r1 现场笔记（实现腿：把球装进常驻那条腿）

工单：`.scratch/wisp/issues/228-ball-and-tray-are-not-in-the-process-that-runs-tasks.md`
射程（编排者 09-30 派单写死）：**只碰 AC#1（宿主存在性有生产点数）与 AC#7（三处过期文案）**。
⛔ 不动 AC#2／AC#3／AC#4／AC#5，不设 `runSpec.replyVeto`，不把 `NewChannels()` 换成 `DefaultChannels()`，
不碰工单勾选框，不写台账。

## 一、要做什么

1. **AC#1**：`cmd/wisp` 那条真正会跑任务的腿（无参数＝`runResident`，见
   `cmd/wisp/resident_windows.go:24`）今天零引用 `internal/ball`。本轮让这条腿**自己**把球起来：
   分层窗口＋托盘＋四枚全局热键＋`ui-sta` 线程，全部走 `internal/ball` 现成 API，
   **零新增模块依赖**（球是原生 GDI＋Direct2D＋`UpdateLayeredWindow`）。
   - 落点：新建 `cmd/wisp/resident_ball_windows.go`（`//go:build windows`），由 `runResident` 调用；
     `balldebug` 只当参考实现，**不改它一行**。
   - 起球失败**不致命**：球起不来时这条腿必须继续跑事件循环，并把"这个进程没有球＋原因"
     响亮地说出来（与 `installLogSink` 失败同一姿态，见同文件 `:58-65`）。
   - 收口：本轮新起的每一枚 goroutine 走 `observe.Registry.Spawn`（本球宿主**不新起任何协程**，
     `ui-sta` 由 `ball.New` 内部以在册名 `ui-sta` 起）；`Ball.Close()` 挂在
     `rt.Shutdown(false)` 之前一层 defer（D38(e) 第 2 步＝"hotkey + wake-word listening stops"
     属球这条腿的活，`internal/proc/shutdown.go:18`）。
2. **AC#7**：三处打给人看／写进日志的过期文案改成**带条件的事实句**，⛔ 不留将来式、
   不新增扫注释票号的词面型仪器。三处（工单 AC#7 射程更正那节列的原文位置）：
   - `cmd/wisp/resident_windows.go:81`（`fmt.Printf`，运行时输出）
   - `cmd/wisp/main.go:25`（在 `const usage = ` 原始串里＝`wisp -h` 的用户正文）
   - `internal/proc/boot_windows.go:127`（`slog.Info`，跨了本票地界，交件时具名声明）
   ⚠ 已知约束：`cmd/wisp/resident_sink_nail_127_windows_test.go:97` 用字面量
   `residentReachedLoop = "empty event loop running"` 断言那行输出，所以改后的句子
   **必须仍含 "empty event loop running"**（这条腿今天依旧不跑任务，短语仍然为真）；
   放宽断言＝禁区。
   ⚠ 另一枚仪器：`cmd/wisp/leg_dispatch_gate_133_test.go` 的 `usageBare133 = ^  wisp {2,}\S`
   会读 usage 块，所以 `wisp` 那行的行首形状不得改。

### 本轮球的行为边界（写清楚，免得被读成"已经能干活"）

- 有：窗口真出现在桌面、可拖拽／可贴边（球自己的窗口逻辑）、托盘图标＋四枚菜单项在、
  四枚热键真注册（默认 `Ctrl+Alt+Q/M/P`＋`Esc`，来自 `ball.DefaultHotkeys()`，
  与 `ApplyHotkeyDefaults` 的「宁可看得见也不能哑」口径一致）、状态＝`Sleeping`（真话：这条腿没事干）。
- 没有：⛔ **本轮不把任何手势推进 D43 状态机**。这条腿今天没有任务通路、没有麦克风、
  没有审批门；把点击变成 `Listening` 就是"对外声称有一条它没有的取消／聆听路径"，
  与编排者作废 `replyVeto` 那条同一形状。故十枚 `Events` 回调全部**具名落日志**
  （经 `installLogSink` 落盘，票 117/127 的监听器在），说清"手势到了、这条腿没有执行者"。
- 托盘「退出」今天也没有执行者（详见 §三 待办与交回报告的"没做完"节）：
  要给它执行者，只有两支——给 `internal/proc` 加一枚"请求停机"钩子，
  或在常驻进程里多挂一枚**不在 D38b 名册**的等待协程。两支都超出本轮射程，按「未定义即停」交回编排者。

## 二、起手名册（2026-09-30 现跑，终态判据＝"红名册等于本名册"，不是全绿）

- `git rev-parse --short HEAD` = `3ee9b3ff`
- `git status --porcelain -- internal cmd` = **空**
- `git status --porcelain .scratch/wisp/issues/` = ` M .scratch/wisp/issues/228-ball-and-tray-are-not-in-the-process-that-runs-tasks.md`
  （编排者自己改的票面，**不是我改的**；我不动任何勾选框）
- AC#1 尺（改前）：`grep -rln "wisp/internal/ball" --include=*.go . | grep -v .scratch`
  → **1 行**：`./cmd/balldebug/main.go`
- 桌面三连（动桌面之前）：`tasklist //FI "IMAGENAME eq balldebug.exe"` ＝ 0 枚；
  `tasklist //FI "IMAGENAME eq wisp.exe"` ＝ 0 枚；
  `gh run list --limit 3` ＝ **网络不通**（`Get "https://api.github.com/...": EOF`）⇒ 读不到 slo-full，
  所以本轮**不做任何 CPU／句柄类测量**，只做功能改动。
- 全量基线（红名册来源）：
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1`
  → 原始输出 `.scratch/wisp/probes/228/r1/baseline-test.txt`（红名册跑完逐名记在这里）。

## 三、待办

- [x] 骨架件落盘并 commit（本文件）
- [ ] AC#1：写 `cmd/wisp/resident_ball_windows.go`（宿主＋十枚回调＋带条件的状态句）
- [ ] AC#1：`runResident` 里接上（起球 → 打印事实句 → 跑循环 → defer 收口）
- [ ] AC#1：钉死它——`cmd/wisp` 加 AST 可达性闸门（删掉那枚调用必须红），
      另加一档 `winlive` 真桌面用例（真窗口＋真 pid 归属）
- [ ] AC#7：三处文案改完（含 `internal/proc` 那一行，交件具名声明）
- [ ] 复跑：`go build ./...`、`go vet ./cmd/wisp ./internal/ball ./internal/proc`、
      `PATH=... go test ./cmd/wisp ./internal/ball ./internal/proc ./internal/observe -count=1`、
      全量 `./cmd/wisp ./internal/...`、`go run ./tools/d22scan`
- [ ] 真桌面验收：`build/wisp.exe` 无参数起（默认档，非 winlive），看窗口／托盘／热键读数，
      再用 CTRL_BREAK 让它走 D38(e) 收口；跑之前再核一遍桌面三连
- [ ] 交回报告：commit 清单／AC#1 两发读数／三处逐字新句／逐名测试结果／
      终态名册对比／"没做完"一节（非空）

### 已知会缺、且本轮不补的（照实交回，不默默收窄）

1. 托盘「退出」无执行者（见 §一 末段，两支行都由编排者裁）。
2. 球的尺寸／贴边点击穿透等 `[ball]` 配置项**没读**：常驻腿今天整个 `config.toml` 都没接
   （`cmd/wisp/config_reload.go` 属 `run` 那条腿），所以尺寸＝`ball` 的编译默认 56px。
3. 位置持久化（`ball.OpenPositionStore`）＝ nil：落盘路径没有规格依据，不自己发明。
4. 手势→D43 状态机→任务通路（AC#2／AC#3）、`replyVeto` 与 `DefaultChannels()`（编排者已作废）、
   `rt.close()` 两枚 join（AC#4）：全部本轮不动。
