# 票 295 `AC#4` 成对复跑（编排者亲跑，2026-10-13 16:22–16:50）——"改前⇄改后"各 2 发、同树、逐名作差

件主＝编排者（⛔ 实现者；本件是 10-12 收腿押下的那半把"差集尺"的补交付）。
所有原始日志在 `C:\Users\swq\AppData\Local\wisp-test-gotmp\orch295-pre-20261013\logs\`（⛔ 入仓；本件是摘要，引用带文件名与 sha）。

## 0. 设计（为什么是这样一把尺）

- **同一枚工作树 @`28b083d4`**（`git status --porcelain -- cmd/wisp`＝0 行；四枚修复件工作树 sha256 ⇄ `HEAD` blob 逐枚相等）。
- **改前侧**＝`go test -overlay=ov-pre.json`：把修复动过的**4 枚件**逐字节换回 `ef8d16b3`（＝`d5ef9cb9^`）的 blob。alt 四枚 sha256 ⇄ `295-r1` 的 `05-pristine-hashes.txt` 第一块（"AC#1 之前的 pristine"）**逐枚相等**：
  - `771c3a62091ea0f4902cdaf2ec350a3062791722dd400d9a82d5ffdf248bbc25` resident_audio_windows.go
  - `4097b03a3ab25ea79a4e043599d126dd906c95669811b0f2fdf5e495542cc327` resident_audio_247_windows_test.go
  - `cc7ac51edff3751b44fa981391ddd48bb83f910b59766e93405d23df20969043` resident_audio_247_live_windows_test.go
  - `2f74d148e4c77aa02e37d79ac9112c4c14e6c7b1ef5612dccc2fbf536f7d68b4` resident_mute_290_windows_test.go
- **改后侧**＝不挂 overlay（c1/c2），同一棵树。
- 四发同命令 `go test -count=1 -v ./cmd/wisp/ ./internal/audio/`、同 `PATH`（sherpa + build）、`GOTMPDIR` 逐发独立目录（全在 C: 卷专用根）；发射前 `tasklist` 家族命中＝0。
- 脚本＝`.scratch/wisp/probes/orch/20261013/pre-collect-295.sh`（入仓）；overlay sha256＝`27cdcb0dc072422f2d0c4b70bab6326b119442ff7eb99b300188a62adeeb29d0`。

## 1. 读数（逐字，源＝`logs/orch-295-pre-collect.txt`）

| 发 | side | rc | lines | sha256(前 12) | census |
|---|---|---|---|---|---|
| a | 改前(overlay) | 1 | 2731 | `31b2f5f94fa4` | 452 PASS / 2 FAIL / 2 SKIP / 457 RUN |
| b | 改前(overlay) | 1 | 2721 | `a2e931a880ba` | 452 PASS / 2 FAIL / 2 SKIP / 457 RUN |
| c1 | 改后 | 1 | 2733 | `16755230a525` | 452 PASS / 2 FAIL / 2 SKIP / 457 RUN |
| c2 | 改后 | 1 | 2720 | `61f9404bb950` | 452 PASS / 2 FAIL / 2 SKIP / 457 RUN |

- 每发两包：`FAIL github.com/CarlosShao/wisp/cmd/wisp`（396–407s）⇄ `ok …/internal/audio`（15.9–16.6s）。
- flake 对拉：`diff roster-a roster-b`＝**0 行**；`diff roster-c1 roster-c2`＝**0 行**。
- 交集：`comm -12` ⇒ pre∩＝**456 行**（452 PASS / 2 FAIL / 2 SKIP）、post∩＝**456 行**（同四数）。
- ★**pre∩ ⇄ post∩ 双向作差＝0 行**（`diff` 空、`comm -13`＋`comm -23` 皆空）⇒ **修复未改变任何用例的任何颜色**。
- 两侧同红＝`TestAC4FocusReturnToPriorWindowGap33r5` ＋ `TestAC14GoSideEvalPushReachesThePage`（= 台账 `:14768` 真窗族的本波稳态组合）；两侧同跳＝`TestAC247LiveMicrophoneLevelsReachTheBallSeam` ＋ `TestLiveWasapiSmoke`（环境门）。
- 旁参（⛔ 不承重）：同族 07:57 三发（`orch295-20261013-075707`，head=08d5a39c）＝pristine 红 `AC13ColdStart`＋`AC14GoSideEvalPush`；M1 红 4 枚（+`BothMuteGestures`＋`AC4Focus`）；M2 红 7 枚（再 +`OutcomeIsReadOffTheGate`＋两枚 secret 家族）。⇒ 真窗族颜色随会话漂、单波内稳。

## 2. ★"坏世界真被造出来"的正负控（⛔ 让 §1 的空差集变成恒真式）

- 正控＝`ov-check.json`（与 `ov-pre.json` **同键同形**，只把 `resident_audio_windows.go` 那枚值换成 `package wisp` 的坏件）：编译 **rc=1**，报 `found packages main (always_write_no_clobber_226_test.go) and wisp (resident_audio_windows.go)` ⇒ **替换内容真到达编译器 = overlay 的键形与值读都真生效**（件 `logs/ov-check-applied.txt`）。
- 负控＝同树不挂 overlay：**rc=0**（件 `logs/ov-check-control.txt`）。
- ⇒ §1 的"改前侧"确实编译的是 `ef8d16b3` 形状的 4 枚件；空差集是**真读数**、⛔ 是台面故障。

## 3. `AC#4` 另三数（同一台面，`logs/` 同根）

- `GOFLAGS= go build ./...` ⇒ **rc=0**（`gate-build.txt` 空输出）。
- `gofumpt -l cmd/wisp` ⇒ **2 枚**：`cmd\wisp\models.go`、`cmd\wisp\panel_transport_35r2_test.go`（`gate-gofumpt.txt`）。两枚⛔ 是 4 枚修复件（无交集）⇒ 修复净增 0 枚；两枚均系既存（`models.go`＝CRLF 假枚那族；`panel_transport_35r2_test.go`＝票 298 残余）。
- `sh scripts/d22scan.sh`：
  - 工作树直跑 ⇒ **rc=1**，唯一红因＝`frontend/node_modules.bak-oldpkg/**`（distinct 路径 **35** 枚，全在该**未跟踪**依赖备份目录里；**0 枚**碰任何 tracked 路径、0 枚碰修复面）（`gate-d22scan.txt`）。
  - 按票面"**纯净树** rc=0"判据，在 `git archive HEAD` 仓外导出树（`orch295-clean-20261013`，`.go` 枚数＝1003）重跑 ⇒ **rc=0 clean**（`d22scan-clean-tree.txt`）：分母逐字 `bans#1-5 internal/=231, cmd/=40; ban#6 frontend/=107; ban#7 internal/tools/=23; ban#8 design/=30, frontend/=107, internal/=533, cmd/=121`。

## 4. 归属与残余

- 工作树那发 rc=1 **⛔ 本票的账**：`node_modules.bak-oldpkg`＝依赖树备份（合并/换锁那次工程动作的产物，未跟踪）；要不要清/挪＝另议（⛔ 本件处置）。
- `Q-295-r1-archive-set`（腿的两棵导出树七项）＝**以本件替代而闭**：本件把"改前侧"在同一棵树上造出来并自带正控，比导出树更可控；腿原件留档（`probes/295/r1/logs/`）。
