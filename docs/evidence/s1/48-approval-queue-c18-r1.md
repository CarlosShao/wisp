# 48 — ApprovalQueue full (C18) 交付裁决表 r1

- 腿：agent-96（写腿，本波唯一 cmd/wisp 写腿）
- 日期：2026-11-16
- 票：`.scratch/wisp/issues/48-approval-queue-full.md`（5 格全做）
- 派单：`.scratch/wisp/probes/orch/20261116/txt/48-approval-queue-dispatch.txt`
- 结论：**5/5 格全绿；零生产码改动**（五格的生产机制已在库中，本腿只落行为牙）；不涉及 C18 契约文字变更。
- ⚠ 本表由实现腿自书（裁决者≠实现者未满足），对抗验收需由另一 agent 复跑落档。

## 0. 交付物（全部为测试文件 + 本表，无生产码 diff）

| 件 | 内容 |
|---|---|
| `internal/agent/approval/ticket48_c18_full_test.go` | 格 1/2/3/5(内部面) 四枚测试：`TestTicket48FIFOFiveTasksHeadOnlyCorrelationRouting`、`TestTicket48TimeoutLadderWarn270Reject300Replay`、`TestTicket48TimeoutLaddersIndependentPerItem`、`TestTicket48FailClosedBridgeDeathMidPending`、`TestTicket48DepthNumbersSyncUnder100msAndEmptyQueueTransitions` |
| `cmd/wisp/approval_queue_48_windows_test.go` | 格 4/5(面板面) 两枚测试：`Test48SuspendFixtureVoidsPendingAndResumeRefusesStale`、`Test48PanelRosterJoinDepthSyncUnder100ms` |
| `docs/evidence/s1/48-approval-queue-c18-r1.md` | 本表 |

写面核对：`internal/agent/approval/**`、`cmd/wisp/**`、`docs/evidence/s1/48-*` — 均在派单写面清单内。`internal/panel/**` 与 `internal/agent/**` 生产文件零字节改动（格 5 面板半边复用生产 join 纯函数，不需要新码）。

## 1. 格 1 — 5 并发任务 × L2 请求（FIFO / 队头单显 / correlationId 路由）

测试：`TestTicket48FIFOFiveTasksHeadOnlyCorrelationRouting`（PASS，0.03s）。

- **FIFO 逐任务次序读数**：入队次序 = 队列读出面 `LiveApprovals` =
  `[corr48-1 corr48-2 corr48-3 corr48-4 corr48-5]`，Position = 1..5 逐枚对上
  （测试输出：`48 AC#1 FIFO admission order: [corr48-1 corr48-2 corr48-3 corr48-4 corr48-5]`）。
- **只显队头**：`Panel().Head()` 在深度 5 时只回 corr48-1、Depth=1；其余 4 张经 `View(corr)` 各自可寻址且 Depth=各自位次；队列清空后 Head 返回 ok=false。
- **决定各自回原任务**：乱序答复（2拒→4拒→1拒→5许→3许）逐枚落到**该卡自己的**停车调用上：拒的三枚收到 AnswerReject + 带各自 corr 的理由句，许的两枚收到 AnswerAllow；每步断言其余四枚的 channel 上无交叉投递（mis-routed reply = fail）。
- 反向先钉：卡1 的 grant 花在卡3 → `ErrBadGrant`、队列不动；未知 corr 拒绝 → `ErrUnknownCorrelation`。
- 审计面：5 枚答案各有一条 `ANSWER-REJECT/ANSWER-ALLOW corr=<各自的>`；伪造尝试有 `FORGED-OR-STALE` 行。

## 2. 格 2 — 超时梯（270s 警告 / 300s 自动拒 / replay 补完；独立计时）

测试：`TestTicket48TimeoutLadderWarn270Reject300Replay`、`TestTicket48TimeoutLaddersIndependentPerItem`（均 PASS）。判据 = 注入时钟 `fakeClock`（包内 ticket-11 模式），**无一处墙钟 sleep 当判据**（A1886）。

三读数（单卡梯）：
1. **270s 警告可见**：Advance(270s) 后恰好 1 枚 `EventWarning`，corr 正确、`Remaining=30s`、文案含「30 秒后自动拒绝」；此刻卡仍在队（深度 1、未答复）。
2. **300s 自动拒**：再 Advance(30s) 后停车调用收到 `AnswerTimeout`，理由含「300」「拒绝」；审计面落 `ANSWER-EXPIRED … decision=timeout->reject after=300s`；队列深度 0。卡面自身 Deadline 字段=300s。
3. **replay 补完原请求**：`Gate.Replay` 对已超时卡返回原 Decision（tool/reason/paths 全同），新 corr=`corr48-timeout#replay`；重放卡带**新 grant**，且旧卡死 grant 花上去= `ErrBadGrant`（陈旧令牌不得跨越），新 grant 放行 → 停车调用收到 `AnswerAllow`（补完），队列清零。

独立计时（双卡相隔 100s）：A 在 T=270s 自警、T=300s 自灭（深度降 1、B 无恙）；B 在自己的 T=270s（墙面 T=370s）自警、自己的 T=300s（墙面 T=400s）自灭。两枚警告 corr 各自正确、次序 = A 先 B 后。

## 3. 格 3 — fail-closed：pending 中杀桥 ⇒ 自动拒 + 任务存活

测试：`TestTicket48FailClosedBridgeDeathMidPending`（PASS）。红句（判据句）：

> **「任务上下文已结束，审批请求已作废并按拒绝处理」** — 桥（派发 ctx）在 pending 中被杀后，停车调用收到 `AnswerReject`，理由句即上句（含「作废」「拒绝」，never allow）。

配套三支：
- **任务存活**：task root ctx 在桥死后仍活着（`taskRoot.Done()` 未触发——C18「拒绝后任务 root ctx 不取消」）；死卡上的迟到拒绝与迟到放行（携带死卡 capture 的 grant）都回 `ErrUnknownCorrelation`，grant 随卡死。
- **界面死 ⇒ 下卡 fail-closed**：UI.Prompt 返回错误后请求立即被拒，理由含「不可达」「fail-closed」，卡片未入队（深度 0）。
- **队满 ⇒ 溢出请求 fail-closed**：MaxPending=2 塞满后第 3 枚请求被拒，理由含「已满」，未被展示（展示数=2）。

## 4. 格 4 — 挂起夹具（pending 审批作废 + 通知；resume 不接受陈旧审批）

测试：`Test48SuspendFixtureVoidsPendingAndResumeRefusesStale`（PASS，cmd/wisp，windows 构建面）。夹具 = 票 43 的 `newSuspendFixture43`（真 gate、真队列项、真 `runSuspendSequence`），本腿只加牙。

- **作废读数**：两张 pending L2（FIFO 位次 1/2）经生产挂起序列后 `voided=2`，审计行 `SUSPEND-L2-VOIDED refused_l2=2 voided_l1=0`。
- **通知读数**：stdout 出现 **「错过 2 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）」**（计数 N=2 与作废数一致）。
- **队列/账本**：挂起后队列深度 0、账本 Pending=0、`Look("corr-48-a")` 不在。
- **两枚停车调用**都收到 `AnswerReject` + 含「挂起」的理由。
- **resume 后不接受陈旧审批**：(a) 挂起前 capture 的 grant 在恢复后花费 → `ErrUnknownCorrelation`；(b) 指名死 corr 的拒绝 → `ErrUnknownCorrelation`；(c) 恢复后新卡拿到**新 grant**（与死 grant 字节不同）且能正常放行（AnswerAllow）——证明拒的是「陈旧」，不是「一切审批」。

## 5. 格 5 — 徽标/面板深度同步 + 空队列迁移（给数字，可判事实）

测试：`TestTicket48DepthNumbersSyncUnder100msAndEmptyQueueTransitions`（内部面）、`Test48PanelRosterJoinDepthSyncUnder100ms`（面板 join 面）。均 PASS。

- **深度数字**：MaxPending=12 夹具逐层填到深度 12（两位数）：每层读出 卡面 Prompt.Depth=当前层数（badge 的唯一生产来源，ballCardUI 原样 post）、`Head().Depth=1`（队头单显数字恒 1）、`LiveApprovals` 末行 Position=12。逐层排空：每步队列深度与头卡 corr 都在**答复返回后的第一次读**就对上（下一步 corr48-deep-(i+1) 即成新队头）。
- **同步判据（可判事实，⛔墙钟猜想）**：answer→projection 的同步性以「**额外轮询数 = 0**」判——答复（拒/许）返回后的**第一次**生产面读取（Depth / LiveApprovals→TaskRosterSectionFrom join / Head）就已见结算后状态；>0 次额外轮询才见 = 投影滞后 = 红。三次排水读数：`extra polls=0 (want 0), wall reading 0s/511.6µs/0s (sub-tick prints 0s), bound 100ms`；allow 方向同样 extra polls=0。墙钟读数同时记录且全部 << 100ms，但**判据是零额外轮询**（墙钟在亚 tick 区间没有分辨率，A1886 同理不拿它当判据）。
- **面板 join 逐任务读数**：loop 生产形状 corr（`<taskID>#<callID>`）经 `NativeVerdict.CardView` + `panel.TaskRosterSectionFrom`（生产 pump 同一纯 hop）点亮各自 task 行；3 卡全亮（blocked=3）；逐枚结算后该行在第一次读即熄、队头推进。
- **空队列迁移**：结算至空 → `Head()=ok=false`、`Depth()=0`（badge 输入=0，即 ballCardUI 将贴出的数）、`LiveApprovals=[]`、join 后无一行 BlockedOnApproval；空队列上的答复 → `ErrUnknownCorrelation`。

## 6. 整包与门禁

- `go test ./internal/agent/approval`（全量，含新 5 枚）**ok 5.99s**。
- `go test ./internal/panel`（全量）**ok 21.41s**。
- `go vet` 三包（approval / panel / cmd/wisp）clean；`gofmt -l` 两枚新件 clean。
- `tools/d22scan -root ../../`：**clean - no D22 ban violations**。
- `go test ./cmd/wisp` 整包经载具 `scripts/wisp-cli-staged-tests.sh` **一发**（跑前 `Get-Process wisp.test` CLEAR）：
  467s，日志 `.scratch/wisp/probes/orch/20261116/txt/48-fullrun-r1.log`。
  **RUN=512 / PASS=378 / FAIL=0 / SKIP=1**（SKIP=TestAC6CSPOnTheNotBuiltNoticeIsReadFromItsRealResponse33lane2，票 33 lane2 自带环境门，与本腿无关）。
  进程在 WebView2/panel 区段以 **exit 139** 死掉（台账**票 450** 记录的同一死法："local full-package run dies at about four hundred seconds inside the webview region"，状态未完成、非本腿引入）。
- A1893 复量：截断点那一枚 `TestPanelThreadIsSTAAndExitsCleanly` **单跑经载具 1.29s PASS** —— 整包 139 是进程级死，不是该测试自身红。
- 直跑 `go test ./cmd/wisp` 出现 `exit status 0xc0000135`（DLL 装载，System32 旧 onnxruntime 遮蔽）——A1888 说的裸跑不可判；本腿未把裸跑读数当结论，全部结论取自载具跑。

## 7. A1882 整包锚点判读

- **终点锚 `TestZZ315A1`：整包未命中**（进程死在锚点区段之前）。载具单跑复量：三枚全中 ——
  `TestZZ315A1RootLoopFirstBridgeCall` PASS 0.40s、`TestZZ315A1SiblingChildFirstBridgeCall` PASS 0.36s、`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` PASS 0.29s（stage-cli/wisp-cli-tests.log）。
- **最后一条 `=== RUN`（截断点）**：`=== RUN   TestPanelThreadIsSTAAndExitsCleanly`（整包日志 line≈1944；其 PASS/FAIL 终止行未来得及打印，进程随即死）。单跑复量 PASS 1.29s。
- **新测试位置**：`Test48SuspendFixtureVoidsPendingAndResumeRefusesStale` 与 `Test48PanelRosterJoinDepthSyncUnder100ms` 在整包日志 line 502/511 处 `--- PASS`，在截断点**之前**；`TestTicket48*` 五枚属 `internal/agent/approval` 包，不在 cmd/wisp 整包分母里（由 approval 包全量跑覆盖，见 §6 第一行）。
- 环境红名单：本次整包 `--- FAIL` 计数 = **0**，`TestAC4Focus…`/`TestAC246…` 未出现红（不需要援引环境红豁免；Wisp 球 PID 49348 在跑未被触碰、未杀）。

## 8. 硬约束自查

- ⛔ `thresholds.go`/golden/`allowlist.txt`/`docs/PLAN.md`/`docs/specs/**`/D43/`design/**`：未动（diff 全为新增测试件 + 本表）。
- ⛔ 放宽断言：无（全部为新增行为牙；既有测试零改动）。
- ⛔ push：未 push（只 commit）。
- 新协程：全部经 `observe.Registry.Spawn`（名字用 D38 名册的 `approval-waiter`；观察日志里 `lifecycle-43-*` 的 WARN 来自票 43 既有夹具，非本腿新增）。
- ⛔ 反向动 C17 白名单：未动（本腿零生产码；面板 join 走的 `NativeVerdict.CardView`/`TaskRosterSectionFrom` 都是既有出站面）。
- ⛔ 提交 dist 字节：未提交。
- 契约文字（C18 语义）：无需变更 —— 五格在现有 C18 实现上全部可达（超时=300s 一律拒、队头单显、深度计数、原生-only allow、fail-closed、一键重放、按 corr 路由全部已实现并有既有牙）。**无「要改契约」的格子**。
- 挂起语义（§16.11#10 的 D43 表附录文字）：票面写明「implement + document as the D43 table addendum; human review at 54」——挂起作废**行为**已实现并有牙（票 43 落地 + 本腿格 4），D43 表的**文字附录**属契约文档面（需人工批准面），本腿未动、留给 54 评审；这不算「哪格没打完」，是票面自己划的人工评审点。

## 9. 复跑尺

```
go test ./internal/agent/approval -run 'TestTicket48' -count=1 -v
& D:\work\soft\Git\bin\bash.exe -c 'scripts/wisp-cli-staged-tests.sh -test.run "Test48SuspendFixtureVoidsPendingAndResumeRefusesStale|Test48PanelRosterJoinDepthSyncUnder100ms" -test.v'
```

（整包：`scripts/wisp-cli-staged-tests.sh`，串行单发，跑前 `Get-Process wisp.test`。本波整包在 webview 区段死于 exit 139（票 450 已知未完成项）；cmd/wisp 包内全部测试在多次分段单发中 0 红。）

## 10. 裁决者注

- 本表由实现腿自书；按 D22 双角色规则，对抗验收须由另一 agent 复跑 §9 两把尺 + 整包后另立表。
- 整包 exit 139 的量测留痕：`.scratch/wisp/probes/orch/20261116/txt/48-fullrun-r1.log`（完整日志，含截断点行号）。
