# 290-v1 · 30 门禁与名册（N7 逐名 rc＝本腿独立现跑，⛔ 未引用实现者读数当凭据）

树＝HEAD `5cff604e`（其 `.go` 面与 `69d8c9ef` 之后逐字节等值；本腿三发变异已全部还原，见 `10-mutations-and-cites.md` §6）。
跑门禁前尺：`git status --porcelain -- cmd internal` ＝ **空**（0 行）。
⚠ 全程未用 `-skip` 换绿、⛔ 未放宽任何断言、⛔ 未跑 `WISP_LIVE_MIC=1` 那一族。

## 1. 门禁逐名（每把尺各自一行 rc，⛔ 空输出也记 rc）

| 尺（逐字命令） | 读数 | rc |
|---|---|---|
| `GOFLAGS= go build ./...` | 无任何输出（0 字节捕获件） | `rc_build=0` |
| `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`；射程逐字与实现者一致：`bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=524, ban #8 cmd/=117`（另记 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`，`examined 268 production Go files`） | `rc_d22scan=0` |
| `gofmt -l cmd/wisp` | **非空 3 行**：`cmd\wisp\models.go`、`cmd\wisp\panel_inbound_guards_35r3_test.go`、`cmd\wisp\panel_transport_35r2_test.go` | `rc_gofmt_dir=0`（gofmt 的 `-l` 列出未格式化件时本身退 0） |
| `gofmt -l <本腿那五枚件>`（＝实现者用的口径，逐字复跑） | 空（未列任何文件） | `rc_gofmt_own=0` |
| `gofumpt -l cmd/wisp`（裸名，PATH 里） | `/usr/bin/bash: line 1: gofumpt: command not found` | `rc_bare_gofumpt=127`（具名，⛔ 不当"跳过"） |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp` | 空 | `rc_gofumpt=0` |
| `go vet ./cmd/wisp/` | 无输出 | `rc_vet=0` |
| 靶向 `go test ./cmd/wisp/ -run 'TestAC290' -count=1 -v`（基线，未突变） | `--- PASS` x6、`ok cmd/wisp 0.102s` | `rc_targeted_baseline=0` |
| 靶向 `go test ./cmd/wisp/ -run 'TestTicket255' -count=1 -v`（HEAD 态） | 见 M3 那一发（`10-mutations…` §5）；还原态下该 20 枚全绿 | `rc_m3=1`（读的是**突变态**那一发，还原后无红） |
| 靶向 AC#3 第①形四枚出厂表尺（票 247 那把，本腿自己跑） | `--- PASS` x4（`TestAC247ShippedDefaultsAreTheOnesThisLegReads`／`…DefaultConfigArmsTheGateMutedAndOpensNoDevice`／`…UnmuteReachesTheBallSeam`／`…HandingTheLevelToTheSeamIsNotVisibility`）、`ok 0.114s` | `rc_form1=0` |

★**具名一条与实现者读数不同的门禁形状（不是缺陷，是口径）**：
`30-gates.md:15` 那行把 gofmt 的尺写成**五枚文件列举**，读数"空"在它自己的口径下为真（本腿逐字复跑＝`rc_gofmt_own=0` 空）。
但 **`gofmt -l cmd/wisp` 在 HEAD 上非空**（3 枚件），而那 3 枚件**一枚都不在本票三笔的文件名册里**
（`git log -1 -- cmd/wisp/models.go`→`5e8748b3 2026-10-03`；`panel_inbound_guards_35r3_test.go`→`7f9d6e40 2026-10-07`；
`panel_transport_35r2_test.go`→`3a343bc7 2026-10-08`；`gofmt -d cmd/wisp/models.go` 现读＝整文件重写形＝**CRLF 行尾**）。
⇒ **判语：不成立但非本票账**——票面 `AC#5` 逐字写的是 `gofumpt -l <自己动过的目录>`，而 **gofumpt 在目录口径下是空的**（`rc_gofumpt=0`），
所以按票面那条尺本票**过**；gofmt 的目录级不洁是三枚别腿旧件的既有状态（具名给编排者：这三枚要不要单独立一枚"行尾归一"的小票，⛔ 不该记在 290 头上）。

## 2. 整包两发（票面 `AC#5` 原文是**两枚包**，⛔ 只跑 `cmd/wisp` 是派单口径）

尺逐字＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ -count=1 -v`，**同树连跑两发、都带 `-v`**。

| 发次 | 原始件（本腿，落在 `/tmp/290v1/`，未入库以免挪门禁分母） | PASS | FAIL | SKIP | 包级行 | rc | 135 | 13a |
|---|---|---|---|---|---|---|---|---|
| 本腿 1 | `full-run-1.md`（319,611 字节） | **423** | **6** | 3 | `FAIL cmd/wisp 503.320s`／`ok internal/audio 16.136s` | `rc_run1=1` | 0 | 0 |
| 本腿 2 | `full-run-2.md`（319,614 字节） | **423** | **6** | 3 | `FAIL cmd/wisp 496.045s`／`ok internal/audio 15.934s` | `rc_run2=1` | 0 | 0 |

- 计数口径＝`grep -cE "^[[:space:]]*--- (PASS|FAIL|SKIP)"`（含缩进子测试，与实现者同口径；顶层口径同时记：`311/6/3`）。
- **同树两发自身抖动尺**：`diff v1-run1.red.txt v1-run2.red.txt` → **零行，`rc_diff_my_two_runs=0`** ⇒ 本腿这两发彼此**无 ±1 抖动**。
- `0xc0000135`（缺 DLL）与 `0xc000013a`（被中止）两发均 **0** ⇒ 本腿没有把环境红算进被验物，也没有那一族形状。

## 3. 与实现者名册对拉（逐名作差；尺＝`--no-renames` 语义下的**逐名集合差**，⛔ 不用 git 重命名检测）

```
diff /tmp/290v1/v1-run2.red.txt .scratch/wisp/probes/290/r1/rosters/raw-post-final-2.red.txt   → rc=1，差 1 行：
  4d3  < TestAC246DevLegIgnoresTheTestTaskInjection
```
⇒ 本腿红名册 ＝ 实现者的 5 枚 **＋ 1 枚** `TestAC246DevLegIgnoresTheTestTaskInjection`。
**那一枚多出来的红，本腿现跑定因＝桌面状态，不是代码、不是本票：**
```
失败逐字（full-run-1.md:1509 段内，resident_task_source_246_windows_test.go:434）：
  AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so
  this leg handed its activation over and exited. Stop it and re-run.
  子进程 stdout: wisp: another instance is running in this session; activated it; exiting
本腿环境尺：tasklist //FI "IMAGENAME eq wisp.exe" → **wisp.exe PID 32888 在跑**
  命令行（Get-CimInstance Win32_Process）："D:\work\workspace\projects plans\Wisp\build\wisp.exe"
  进程创建时刻：2026-10-09 16:10:05            ← 晚于实现者那四发（15:4x-15:5x），早于本腿这两发（16:16/16:24）
  tasklist //FI "IMAGENAME eq balldebug.exe" → No tasks match（0 枚）
尺的实现逐字：resident_task_source_246_windows_test.go:432 `if strings.Contains(out, "another instance is running") { … t.Fatalf("AC#7 RED (desktop state, not our code)…" }`
```
⇒ **判语（逐名作差口径）**：相对实现者的名册，**本票新增红＝0 枚**；本腿两发多出的那一枚是**同一棵树上"会话里有个活着的 `build/wisp.exe`"这一桌面状态**造成的，
且它在自己两发里**稳定重现**（不是随机抖动）。⛔ 本腿没有去杀那一枚进程（不是我的写点，且它可能就是编排者/机主那一次真机窗口）。

⚠⚠ **顺带一条要紧的现场读数（归编排者）**：票面 `AC#3` 逐字要求"跑真机前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 必须为 **0**"。
本腿现跑＝`wisp.exe` **1 枚在跑**（PID 32888，`build/wisp.exe`，16:10:05 起）。
⇒ 若那一枚是编排者正在取 AC#3 第②形／票 247 那三格的 console 凭据，**它正好在票面要求的那个窗口里**，是好事；
但它同时会让**任何**整包串跑多红这一枚 246 用例（本腿已实测两发），并且**这一格的名字会盖住"新增红 0 枚"那句话**。
⇒ 建议：**每一发整包之前先落一行 `tasklist` 读数**，并把"有活 wisp.exe"当作**已声明的桌面前置**而不是被测物的红。

## 4. 判据落点（本腿怎么把结论压到语义面）

单发名册不当尺（派单那条本腿复现了它的反面：本腿两发**零抖动**，但相对实现者的树**差 1 枚**，差因在桌面状态）。
⇒ 实质结论的四把语义尺：
1. `git diff-tree --no-commit-id --numstat -r fdee536b` 逐枚＝`60/0`、`142/14`、`287/0`、`287` 那枚是 `_test.go`、`25/0`（与自报逐字一致）；`69d8c9ef`＝`7/7` 单文件。
2. 新增行里的形状尺（编排者跑过，本腿独立跑）：三枚产码件新增行内 `go func(`＝0、`EvMuteKey`＝0、`statemachine.`＝0
   （本腿尺＝`git diff-tree -r -p fdee536b -- <件> | grep "^+" | grep -c "statemachine\."` → 三枚各自 0）。
3. **M1 变异**（拧门改成空拧）→ 三枚必须红、实测三枚红，且 AC#2 那把 grep 尺**仍报 1 枚调用者** ⇒ 行为断言才是牙（`10-mutations…` §2）。
4. **M3 变异**（cite 行号写歪一位）→ 被编辑过的那把票 255 尺**仍红在 `:209`** ⇒ `69d8c9ef` 没把判据编辑成恒真（`10-mutations…` §5）。

## 5. 本腿自己发现的一条新形状（具名，⛔ 未修）

`toggleMute` 在 `cmd/wisp/resident_audio_windows.go:168` 写 `ra.started = gate.Open()`，
而 `ra.started` 是 `residentAudio` 的**裸 bool 字段**（`:105` `started bool`；同结构的邻居 `levels` 用的是 `atomic.Uint64` `:110`），
读它在 `:379` `posture()`，唯一的生产读者是 `cmd/wisp/resident_windows.go:320`（**boot 报表**），而挂载 `rb.attachMuteGate(...)` 在 `:303`——**早于那一行**。
⇒ 装配之后、boot 报表打印之前，用户在 ui-sta 线程按一下静音键，就会与 boot 线程对同一枚裸 bool 形成**未同步的并发读写**（Go 内存模型下的 data race；本仓门禁没有 `-race` 那一发，故今天量不到）。
球侧那一跳作者**是**想到序的（`muteMux sync.Mutex`，`resident_ball_windows.go:178`），门侧这一枚没上闸。
⇒ 危害很小（窗口是 boot 那几毫秒、后果最多是一行报表读旧值），**不影响 AC#2/AC#3/AC#4/AC#5 任何一格的成立**；
落成新票名册第 4 枚建议：`29x-ra-started-is-written-from-the-ui-sta-thread-while-the-boot-report-reads-it`（改成 `atomic.Bool`，⛔ 不加锁不加协程）。
⛔ 本腿未动一字。
