# 134 AC#4 取证 — 采样有效性前置检查真触发（shape C，2026-09-23）

被验版本：本文件所在提交（工作树 = `scripts/slo-check.ps1` 的本次改动，锚定 sha `ac6f31c`）。
执行机器：本机（`COMPUTERNAME` 见读数），**不是** runner 的 `_work` 目录 —— 见文末"取证边界"。

## 1. 前置检查的三态读数（真跑，rc 都在同一行末尾）

### 1a. 并发负载 = 真 `go.exe`（名称探针）⇒ FAIL machine-contended

负载源：`go vet -a ./internal/...`（`-a` 强制重编全部依赖，实测 `go.exe` pid 47420 起来）。

```
after 1x0.5s: toolchain procs=1
go.exe                       47420 Console                    1     41,716 K
=== gate ===
slo-check.ps1: clearing 1 stale report file(s) from C:\Users\swq\AppData\Local\Temp\wisp134\out
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: FAIL machine-contended - subset=full refused to sample, no numbers were produced
slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: go.exe pid=47420 started=2026-09-23 10:27:17 path=D:\work\base\go\bin\go.exe
slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: compile.exe pid=1640 started=2026-09-23 10:27:20 path=D:\work\base\go\pkg\tool\windows_amd64\compile.exe
slo-check.ps1: machine-contended reason: machine-wide cpu utilisation 100% over a 1s window (>= 50%)
slo-check.ps1: machine-contended: 3 reason(s), 0 state file(s) written, slo-report.json NOT written
GATE-RC=1
```

三条 reason 同时命中：两枚外来工具链进程 + 机器整体 CPU 100%。**零个数字产出**（见 §2）。

### 1b. 并发负载 = "runner 目录下的进程"（路径探针）⇒ FAIL machine-contended

造法：把 `ping.exe` 拷成 `…\Temp\wisp134\_work\wisp\wisp\foreignping.exe` 并跑起来，
再只给门禁 `GITHUB_WORKSPACE=<那个假 _work\wisp\wisp>`（模拟 runner 的 checkout 根）。

```
slo-check.ps1: clearing 1 stale report file(s) from C:\Users\swq\AppData\Local\Temp\wisp134\out
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: FAIL machine-contended - subset=full refused to sample, no numbers were produced
slo-check.ps1: machine-contended reason: process running under the runner work root: foreignping.exe pid=8392 started=2026-09-23 10:26:11 path=C:\Users\swq\AppData\Local\Temp\wisp134\_work\wisp\wisp\foreignping.exe
slo-check.ps1: machine-contended: 1 reason(s), 0 state file(s) written, slo-report.json NOT written
GATE-RC=1
```

⚠ 这一发第一次跑时**没**报 contended（`precheck ok … cpu max 23%`）——原因是造负载的那条
bash 链整体被 `&` 吞进后台、门禁先起跑（仪器坑，不是判据坑）；重发前先 `tasklist` 确认 pid
活着才拿到上面这发读数。**空输出/绿灯先怀疑仪器。**

### 1c. 机器安静 ⇒ 放行（同一枚检查的绿色半边，防"永远拒绝"）

```
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 22%
slo-check.ps1: sampling state Sleeping for 6s
```

（放行后继续原步序 —— 这一发 `-WispExe` 指的是 `$env:TEMP\wisp134\notwisp.exe`（`ping.exe` 的拷贝），
所以六个状态全 `exit=1 pass=False`、`all_pass=False`、`GATE-RC=1`：**故意不拿它冒充真 SLO 结论**，
只为证明前置检查不会把安静机器也拒了。）

## 2. "不产出任何数字"的自证

- 拒绝路径的最后一行是 `0 state file(s) written, slo-report.json NOT written`，且 `exit 1`。
- 进入检查之前脚本会清 `OutDir` 里遗留的 `*.json`（`clearing 1 stale report file(s)`），
  因为 runner 复用 `E:\work\base\actions-runner\_work\wisp\wisp`，上一轮的 `slo-report.json`
  会留在 `build\slo\` 里冒充本轮结论。
- 检查点位置：插在 `wisp.exe` 就绪之后、`$states = @(...)` 取样循环之前 ⇒ 拒绝时**一个状态都没采**。

## 3. 两个阈值未动的自证

- 前置检查只**新增拒绝理由**，不读也不写任何 D32 度量。
- `git diff --numstat -- scripts/slo-check.ps1` = **`174  0`**（增 174 行、**删 0 行**）⇒ 纯插入。
- 阈值所在文件本轮零改动（`git status --porcelain` 里 `internal/**` 一个字都没有）。原文照抄：

```
internal/observe/thresholds.go:19memCapSleeping     int64 = 25 << 20      # D32 RSS <= 25MB
internal/observe/thresholds.go:26cpuLimitSleeping     = 0.5 // % of all-core mean, 1min window
```

## 4. 取证边界（如实标未验证）

- 上面三发是**本机**跑的 `scripts/slo-check.ps1`（同一枚脚本、同一枚入口），
  不是 `slo-full` 那枚 job 的步级读数 —— 本票**不 push**，改动还没进 GitHub 侧的 `dev`，
  所以"run id + job id + step 名 + 结论"这一格**欠编排者 push 之后回读**（票面 AC#2 同一格）。
- `GITHUB_STEP_SUMMARY` 那一段（拒绝时往 step summary 写一张 reason 表）在本地没有该环境变量
  ⇒ **未验证**，只在 CI 里才会被执行；未执行不等于执行失败，但不许拿它当已生效。
- "安静机器 ⇒ 放行"用的 `-WispExe` 是假的（§1c），所以**完整六态在改动后能否仍然跑通**
  这一格也**未验证**，等 push 后的真 run。
