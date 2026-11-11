# 票 163 验收裁决表 v1（非实现者只读裁决腿 163v1）

- 裁决人：163v1（编排者派单 · 非实现者 · 只读）
- 派单：`.scratch/wisp/probes/orch/20261111/txt/163v1-dispatch.txt`
- 裁决时刻：2026-11-11 18:3x +08，工作树 HEAD 附近（163 三枚 commit：`75f9479fa` 产码＋测试、`821f4a2d3` 组合根接线、`7bcc2ae08` unwired 文案收尾）
- 方法：读产码＋测试＋组合根＋unwired 现状，逐格引 file:line；**未跑任何 go test/go build（派单硬约束）**，测试格判据＝测试代码以已提交文本在册＋163-impl-r1（2026-11-10）/163-closeout-N（2026-11-11）两条 Progress log 记录的门禁绿读数；硬约束禁区（approval/**、frontend、thresholds 等）以三枚 163 commit 的 `git show --stat` 文件清单为凭。

## 逐格裁决

| 格 | 判 | 判据（file:line / 证据件名） |
|---|---|---|
| **AC#2** 埋标记取退出码、不靠时间；答"摘掉标记是否存在一发超时从此判错" | **PASS** | ①标记协议：每发写 `<命令行> 2>&1` + `echo <随机标记> %errorlevel%`（`internal/tools/shell_session.go:282`），标记＝每发 crypto/rand 12 字节（`shell.go:293-299`）；取码扫描 `findMarkerCode`（`shell_session.go:371-389`，LastIndex＋数字载荷校验）。②不靠时间：整个家族零 `time.Now`/`time.Since`/裸 `go func(`（grep 现量仅命中 `shell.go:59` 的注释）；等待三源 select＝notify token / 流结束 / 本发 ctx（`shell_session.go:312-351`），超时即 C22 ctx.Done，无墙钟差。③摘标记问题的答案在码里成文：**是——摘掉标记则流无边界、本发超时的输出与下一发的输出不可分**（`shell_session.go:50-58` 注释逐字）；代码按"该世界不可达"建造：超时路径 `dropLocked()` 回收整树＋本发**不报退出码、不归属输出**＋`reset:true`（`shell_session.go:344-350`），渲染层超时文案明说"失去同步，不是命令失败"（`shell_session.go:487-494`）。④测试钉住：正向 `Test163SessionMarkerCarriesExitCode`（`internal/tools/shell_session_test.go:195`，码 7/码 0 两臂＋标记不泄漏）；反向（a1 活样本陷阱）`Test163SessionMarkerIsNotAGuess`（`shell_session_test.go:222`——输出里埋 `exit 1` 诱饵、标记不到 ⇒ 不认码、会话回收、下一发全新会话）；单发形反猜对照 `internal/tools/shell_exec_live_windows_test.go`（在册，`shell_exec.go:24-26` 指名）。⑤一次性的退出码走 Wait 状态字段，不走文本（`shell_exec.go:129-149`＋`exitCodeOf` `shell_exec.go:170-179`）。 |
| **AC#3** 会话寿命绑任务，任务取消/结束会话跟着收，有变异证明 | **PASS** | ①名册按宿主 task id 键控，`ShellSessions`（`shell_session.go:75-134`）；唯一收集器 `TerminateForTask`＝关 kill-on-close 树＋收尸＋删行，无第二收口（`shell_session.go:110-123`，注释逐字 "the ONLY collector"）。②任务结束边界：bridge `Options.TaskCleanup` 缝（`internal/tools/bridge.go:124-136`），`CloseTask` 在 C25 close **之前**先跑清道夫并落 TASK-JANITOR 审计行（`bridge.go:1035-1047`）。③组合根真接线：`shellSessions := tools.NewShellSessions()`＋`BuiltinShellEntries(ShellDeps{... Sessions: shellSessions ...})`（`cmd/wisp/run.go:862-875`）＋`Options.TaskCleanup: shellSessions.TerminateForTask`（`run.go:1390`），同一枚对象既被工具看也被清道夫看。④变异证明有牙：`Test163AC3TaskEndCollectsSessionTeeth` 正向（CloseTask 后 count=0、树被关）＋作用域臂（B 任务会话不受 A 边界牵连）（`internal/tools/shell_session_teeth_test.go:155-195`）；负向对照 `Test163AC3UnwiredJanitorIsTheDefect`（拆掉清道夫接线则会话在 CloseTask 后存活——恰是上面那枚测试会翻红的形状，`shell_session_teeth_test.go:202-217`）；单元臂 `Test163SessionRosterTerminateKillsTree`（`shell_session_test.go:275-306`，含异任务 id 收不到、二次收 0）。⑤树的 OS 底座：`NewTree`→`proc.OpenJobScope` 平台拆分 `cmd/wisp/shell_tree_163_windows.go`／`shell_tree_163_other.go`（非 Windows fail-closed 拒绝，`shell.go:84-92` 缝定义、`run.go:851-859` 说明）。 |
| **AC#4** 审批形状不松：会话≠授权一次一直放行；不动 `internal/agent/approval/**` | **PASS** | ①每发独立过全门：Bridge.Execute 每调用都走 C3→C19→路由→C22（`bridge.go:144-153` 流程注释＋dispatch 实现），Facts 每发现填 argv＋allowlist 交 R6（`shell.go:200-207`），两条 Decl 均声 L2 地板（`shell.go:163-192`）；会话存在本身不带任何权限（`shell.go:47-56` 注释逐字 "a session is NOT an authorization"）。②有牙测试：`Test163AC4SessionNeverCarriesAuthority`——两发命令 gate 访问计数＝2、拒绝臂（gate 拒 ⇒ 什么都不跑、拒绝理由透传、会话不被拒绝销毁、再批仍能用）（`shell_session_teeth_test.go:108-148`）；`Test163AC4GrantsNeverCoverTheSessionRoute`——D45 授权源放得再宽也不覆盖 L2 路线，gate 照访（`shell_session_teeth_test.go:222-248`）。③禁区零触碰：163 全部三枚 commit（`75f9479fa`/`821f4a2d3`/`7bcc2ae08`）的 `git show --stat` 文件清单里**无任何 `internal/agent/approval/**` 路径**（清单见本文 §禁区核验），approval 目录现量 grep 也只有历史票的 `shell.run` 测试字面，无 163 形。 |

## AC#5 残格现状（信息格，待机主，不算欠账）

派单要的"unwired.go 里 shell_* missing 文案"现状已**被收尾腿 Q（`7bcc2ae08`）改写**，与 11-10 进度 log 记的"三行未动、等 owner"相比已前进一格：

- **守卫表里三行 shell_* 已不在**：`internal/config/unwired.go:61-72` 的注释逐字说明 `risk.shell_enabled`／`risk.shell_allowlist` 两行因组合根接线（腿 N）落地而撤行——"load-time error must not outlive the consumer that justifies it"；两枚键在 `lockedKeyDisposition` 里改记 consumed 句（`unwired.go:122`＝shell_enabled 喂 `ShellDeps.Enabled`、`unwired.go:124`＝shell_allowlist 喂 `ShellDeps.Allowlist`→R6）。
- **仍守卫的 shell_* 只剩一行**：`risk.allow_shell_string`，守卫行 `unwired.go:73-78`（missing 文案在 `unwired.go:75`：string mode 未实现、`Facts.ShellString` 无生产写入者恒 false、R6 字符串判语永不跑），处置登记 `unwired.go:123`。这一行等的是机主对"字符串模式做不做"的定夺，属**诚实行**，不是 163 的欠账。
- ⚠ 给编排者的观察（非裁决）：11-10/11-11 两条 Progress log 都把"unwired 三行文案"记为"待 owner 一句话"，而 `7bcc2ae08` 已把其中两行撤掉、一行改写。若 owner 对那句措辞的批准尚无台账落点，建议编排者核对 `A` 账是否有对应裁定记录；本腿只报盘上现状，不裁这一步的合规性。
- 票面 AC#2/AC#3/AC#4 复选框今日仍未翻勾、Status 未结（本腿禁改票面，如实记录）；`frontend/**`、`design/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**` 在三枚 163 commit 的文件清单里**零命中**（AC#5 零字节名单核验通过）。

## 禁区核验（三枚 163 commit `git show --stat` 现量）

| commit | 触及 |
|---|---|
| `75f9479fa` | 票面 +2 行、`internal/tools/{bridge.go,shell.go,shell_exec.go,shell_session.go,shell_*_test.go,shell_live_windows_test.go,spawn_roster_census_229_test.go}`、`docs/evidence/s1/163-persistent-shell-session-impl-r1.md` —— 无 approval/frontend/thresholds |
| `821f4a2d3` | `cmd/wisp/{run.go,shell_tree_163_windows.go,shell_tree_163_other.go,...}`、票 388 件 —— 无禁区 |
| `7bcc2ae08` | 仅 `internal/config/unwired.go`＋`unwired_test.go` —— 无禁区 |

## 裁决局限（如实）

- 未复跑任何测试/构建（派单硬约束 ⛔ go test/go build）。"测试绿"一栏采信 Progress log 里 163-impl-r1（`go test ./internal/tools/` ok 32.3s、`./internal/proc/` ok 3.0s、d22scan clean）与 163-closeout-N（vet/build rc0、d22scan clean、path-budget GREEN、整包红集 ⊆ HEAD 基线红集）的**已记录读数**；本表对格的判据全部落在盘上可复读的代码与测试文本，后人可用一把尺复跑推翻。
- 本腿零 git add/commit/push、零产码/票面改动；唯一写面＝本证据件。
