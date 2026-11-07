# 246 AC#8 — boot 时序 flake 真修：`signal.Notify(bootExit)` 钉到 `runResident` 第一句

- 腿：`246ac8`（写腿，编排者派单 2026-11-07）；票：`.scratch/wisp/issues/246-the-process-that-runs-tasks-is-not-the-process-holding-the-ball-so-no-veto-channel-has-an-executor.md` AC#8。
- 地界：本腿只改 `cmd/wisp/resident_windows.go` 一枚产码文件；`internal/ball`、`internal/winsec`、`internal/agent`、`internal/proc`、thresholds/golden、`frontend/**`、`design/**` 零碰（收件读数见 §7）。
- **AC#8 勾未翻**（派单纪律：翻勾归编排者）；票面"地界属票 228、本格不改"的旧路由已被本次派单改道，归属裁定留编排者在票面落格。

## 1. flake 归属与落点（现跑定位，非转述）

| 事实 | 读数 |
|---|---|
| 用例 | `TestAC228ExitRequestDuringBootStillLeavesThroughD38E`（`cmd/wisp/resident_ball_228_windows_test.go:184`，真机：起无参 wisp.exe，每 50ms 轮询 sink 文件**一出现立刻**发 CTRL_BREAK） |
| 未注册窗口（修前） | `installLogSink` 打开 sink 文件（修前 `resident_windows.go:121`）→ `signal.Notify(bootExit,…)`（修前 `:230`）之间：crash-marker 扫描（真开 SQLite）＋ `writeRunMarker` ＋ `markUpdatePendingHandoff` ＋ `autostartSync` ＋ D38(e) 关停 defer 注册 |
| 翻车机制 | break 落进窗口 ⇒ 信号无人注册 ⇒ Go 控制台 handler 返回 0 ⇒ 默认处置杀进程 `0xc000013a`；D38(e) trail、托盘移除、热键归还全丢 |
| 观测对齐 | 窗口内 sink 文件已存在 ⇒ 测试的触发条件已满足；红发读数「sink 1 file(s)、**0 records**」与"文件先出现、记录还没写、Notify 还没注册"三段完全咬合 |

## 2. 修前正控（会翻的证明）

- 命令：`go test ./cmd/wisp -run '^TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=5 -timeout 30m`（env：unset GOTMPDIR/TMPDIR/TEMP/TMP；PATH 加 build/third_party/sherpa-onnx；起跑前 `tasklist` wisp.exe/balldebug.exe 计数 0）。
- 读数：**5 发 1 红**（翻车率 1/5），总时长 46.3s。
- 红发细节（2026-11-07 12:00:37）：`exit status 0xc000013a`；sink 目录 `1 file(s), 0 record(s)`；随后 `a boot-time exit produced no "wisp: persistent log sink installed" record (records: [])`。其余 4 发绿（install=1，shutdown trail 在场）。
- 注：修前产码本就带着"注册在 sink 之后"的形状——这一发红就是该形状的活体对照，等效于变异试验"MUT：把注册从最前挪走会红"。

## 3. 修法（真修，非 skip/放宽/加大超时）

- 改动：把 `bootExit := make(chan os.Signal, 1)` ＋ `signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)` ＋ `defer signal.Stop(bootExit)` 从修前 `:229-231` 提到 `runResident` **第一句**（`printVersions` 之前）；原位注释改写为指向新位置与归因（AC#8），`residentBootWindowExit228` 的消费点、`resident_tray_stop_228_windows_test.go:367` 钉的函数、票 228 AC#11 托盘挂载注释块均一字未动。
- 为什么这是真修：未注册窗口从"[sink→Notify] 约 300ms+"缩到"main 派发＋attachParentConsole 的几条指令"——测试的触发条件（sink 文件出现）从此**必然落在注册之后**，不是降低碰撞概率，是消灭可观察的暴露面。
- 语义对账：① 消费点不变——`residentBootWindowExit228(rt, bootExit, bootStop)`（`select` 判定 boot 窗口退出）照旧；② `defer signal.Stop` 成为**第一个注册的 defer**，LIFO 展开＝最后注销 ⇒ 整条 D38(e) 链执行期间注册仍在场（信号进缓冲，不打断 teardown）；③ os/signal 向所有注册者转发，`internal/proc` `RunEventLoop` 自己的 `NotifyContext`（`boot_windows.go:183`）照常收得到；④ `wisp run`/`doctor` 控制台腿与 `cmd/balldebug` 不经 `runResident`，零影响。

## 4. 修后正控（不翻的证明）

- 命令：`go test ./cmd/wisp -run '^TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=10 -timeout 30m`（同 env；起跑前 wisp.exe 计数 0）。
- 读数：**10 发 10 绿，零枚 0xc000013a，总时长 57.3s**（门禁要求 5 发，实际按票面 AC#8 判据的 10 发跑）。

## 5. 门禁读数表

| 尺 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp` | 空（rc 0） |
| vet | `go vet ./cmd/wisp` | rc 0 |
| build | `go build ./cmd/wisp` | rc 0 |
| 定向 test | `-run '^TestAC228ExitRequestDuringBootStillLeavesThroughD38E$' -count=10` | 10/10 绿（§4） |
| d22scan | `sh scripts/d22scan.sh` | rc 0；real repo scan `clean - no D22 ban violations`；自检 45 PASS / 0 FAIL |
| 整包序 | `go test ./cmd/wisp -timeout 30m`（修后跑两遍） | 两遍各 1 红＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`（见 §6）；其余含本票五枚 commit 相关用例全绿 |

## 6. 整包序那枚红的归因（AC#5 纪律：归因，不顺手修）

- 红名＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`（`config_reload_223_test.go:436/439`：`wisp run` 控制台腿的卡片渲染断言——与本腿改的 `runResident`（no-args 常驻腿）不同代码路径、不同支）。
- 隔离复跑：`-run '^TestTicket223HandEditedFsLooseningCostsAnL2Card$' -count=3` ＝ **3/3 绿**（2.89s/3.35s/2.50s）。
- 判：整包序才翻的包级 flake，非本票造成；本腿零触碰其路径，不顺手修（归票 223 那一族）。

## 7. 禁区对账（可核命令读数）

| 禁区 | 尺 | 读数 |
|---|---|---|
| 只动一枚产码文件 | `git status --porcelain -- cmd internal`（提交前） | 仅 `M cmd/wisp/resident_windows.go`＋本腿两枚新件，0 枚别的突变 |
| thresholds/golden | `git status --porcelain -- internal/observe docs` | 0 行 |
| internal/ball、winsec、agent、proc | `git status --porcelain -- internal` | 0 行 |
| frontend/design | `git status --porcelain -- frontend design` | 0 行 |
| LF | `git ls-files --eol cmd/wisp/resident_windows.go` | `i/lf w/lf`；`grep -qP "\r$"` 无命中 |
| git 纪律 | 提交用两处显式 pathspec；无 add -A/amend/reset/rebase/stash/checkout ./clean；只 commit 不 push | 收件时以 `git log --stat -1` 复认 |
