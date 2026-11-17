# 票 49 — Session grants (D45-2) 判决表

**写腿**: agent-97 · 2026-11-17 · 基锚 `26ee1c081`（票 49 派单存档提交）
**判决**: 实现腿（非裁决腿）—— 本表只记录实现读数与自测读数；缺口审计与对抗验收归非实现者。

## 0. 现读（照派单，基线确认）

- `grep -rn "grant\|Grant" internal/agent internal/panel internal/store --include=*.go`（前 30 行现跑）：
  票 49 在本腿开工前**未开工**——上一提交 `26ee1c081` 只含派单文件本身（`git show --stat 26ee1c081` = 1 文件 17 行）。
  已存在的 grant 面＝票 224 的桥：`internal/agent/approval/gate.go`（allowSession 写侧）、`internal/tools/grant.go`
  （GrantSource 读侧）、`internal/session/grants.go`（Ledger）、`internal/memory/dao_misc.go`（DAO，注释明写
  「Lifecycle POLICY ... belongs to tickets 49/44」）。
- 票面 9-19 行号当今天：以 HEAD `26ee1c081` 现读为准，未照 2026-09-19 的旧票面行号当现状。

## 1. 逐格实现与读数

### 格 1 — Grant lifecycle（create → 静默 → session 结束失效，逐调用审计行）

实现：
- `internal/session/grants.go`：`RevokeAll`（会话收口 stamp，行保留不删）+ `Revoke`（单行撤销）；
- `cmd/wisp/run.go` `agentRuntime.close()`：关停序列新增 `SESSION-GRANT-SWEEP`（会话结束失效的生产落点），
  revoke 行数进审计行；行保留（SPEC-02 §4 审计窗），不删除。

读数（载具单跑 `TestTicket49SessionEndSweepRevokesEveryLiveGrant`，PASS）：
- 一笔「本会话内允许」的授权行 + 一笔被授权调用 → `tool_call` 行落盘 `decision=allow_session_grant`、
  `grant_id=<授权行 id>`；
- run 收口后同一 session id 现读 `approval_grant`：**1 行、`revoked_at` 非空**（会话结束失效落盘）；
- 审计含 `close: SESSION-GRANT-SWEEP session=<id> revoked=1`。

### 格 2 — Revoke 立即生效（红/绿两读数）

实现：`session.Ledger.Revoke`（stamped revoked_at）+ 操作者动词 `revoke-grant <编号>`（console，
cmd/wisp approval_reply.go）。

读数（载具单跑 `TestTicket49RevokeVerbStopsTheNextCoveredCall`，PASS）：
- 绿（撤前）：同形调用 **0 卡**、`tool_call.decision=allow_session_grant`、`grant_id=<行 id>`；
- 红（撤后）：同形调用 **≥1 卡**（L1 窗口重新弹）、`decision≠allow_session_grant`、`grant_id=NULL`。

### 格 3 — Limits（三支）

① L2 永不被 grant 自动放行（有 grant 也仍弹卡）：
- 读数（`internal/tools` 单跑 `TestTicket49L2WithALiveGrantStillShowsTheCard`，PASS）：
  R8 覆盖写（L2）+ **全覆盖** grant 源（promiscuous，Covering/TargetCovering 全答 true）→ **仍 1 张 L2 卡**，
  grant 源 0 次被查询（L2 路由根本不进 grant 检查——route 的 switch 语句形态）。
- 票 224 的 `TestTicket224SessionGrantNeverCoversL2`（正控：告知源 L2 也不查）继续绿。

② tainted exfil 无视 grant 反而升级：
- 读数（`TestTicket49TaintedPathWriteIgnoresEveryGrant` / `TestTicket49TaintedExfilIgnoresTheTargetGrant`，PASS）：
  真 C19 判定器 + t90Taint 全命中探测器 → R4 升级为 L2+SessionOverrideBlocked，全覆盖 grant 源在场
  **仍弹卡**，`tool_call.grant_id=NULL`；input.type 变体下 grant 源 0 次被查询（`targetBindingBlocked`
  在查询前拒绝）。

③ config 层限制逐项：
- `[risk] permission_mode`（放宽档）：读数（`TestTicket49ModeSilenceNeverBooksAGrantLink`，PASS）——
  `auto_approve` 静默 L1 的调用 `decision=allow`、`grant_id=NULL`、**永不** `allow_session_grant`
  （grant 只答「还存在的问题」；`!sil.Silenced` 是闸）。
- `[fs]`/B 档豁免（risk.blacklist_overrides / 单文件豁免）：走 C19 判定与 D36 重新确认，与会话授权无交集
  （票 230/226 现有面，未动）。
- 🔒 三节（risk/fs/net/plugins）放宽走 D36 重新确认：现有 `config_reload.go` 面未动，会话授权不参与。

### 格 4 — input.type（三支各一枚牙）

实现（§15 第 7 项定案逐条）：
- 绑定键＝**解析后的目标进程名**（hwnd→pid→小写映像名，`resolveInputTarget`/`SystemBridge.ResolveTarget`，
  与执行时同一管线）；pid 回退名（`pid:<n>`）不落授权（`isGrantableTarget`，bridge 侧拒查）。
- per-process：`session.Ledger.RecordTarget/TargetCovering`——行=(tool, 进程名, session)。
- all-apps 危险档：`RecordTargetAllApps` 是 **唯一** 能写 `*` 行的门；操作者动词为两步序列
  `session-expand`（印后果 + 30s 窗）→ `session-all`（窗口内才可答）；未展开即答＝拒绝并点名展开步骤。
- 换目标重新确认：judge 时解析 → 卡片印「目标应用：<进程名>」；**执行前 byte 界 swap guard** 重解析，
  与授权名不符 → 拒绝（文案含「目标应用已切换」「重新确认」），不花授权。
- 卡片（cardrender.go）：input.type 卡印目标应用行；`SessionGrant=true` 时印三档选项（all-apps 标危险档）。

读数：
- `internal/session`（真 SQLite，4 枚 PASS）：绑定行答自己点名的进程；`RecordTarget` 拒 `*`；
  `RecordTargetAllApps` 只收 input.type；Word 行不答 终端（隔离）。
- `internal/tools`（真判定链，PASS）：无 grant/无 resolver/解析失败 → L2 照弹（四分支）；
  有 grant → **0 卡**、`decision=allow_session_grant`、`grant_id=<行 id>`；
  R4 命中 → 升级无视 grant；换目标 → lookup 半 + byte 半双双重弹/拒；pid 回退名 → 不查授权照弹。
- `cmd/wisp` 装配（载具单跑，PASS）：红（首注入弹 L2 卡，grant_id=NULL）→ 绿（RecordTarget 后同形调用
  0 卡、`allow_session_grant`、grant_id=<行 id>）→ 隔离（hwnd 9=cmd.exe 重新弹卡）；危险档
  `session-expand → session-all` 端到端落 `*` 行（未展开先答被拒）。

### 格 5 — Batch＋grant（10 操作 + opt-in → 一张卡 → 整 pattern 一次授予）

实现：`approval.PatternForPaths`（卡片自身路径的最长公共祖先目录 + `/*`，拒标记路径/裸盘根/无公共目录）
+ `AllowSessionPattern`（同一 nonce、一张卡、一行 pattern 行，审计行带「整批模式一次授予」）。

读数：
- `internal/agent/approval`（`TestTicket49TenOpBatchOneCardOnePatternGrant`，PASS）：10 路径批 → **1 张卡**
  → `session-pattern` → **恰 1 行** `C:/work/*`（不是 10 行逐路径）。
- `cmd/wisp` 装配（`TestTicket49AssembledBatchOneCardOnePatternRow`，PASS）：10 路径经装配 gate 出
  **1 张卡**，操作者动词 `session-pattern` 答复 → 落 1 行 `<dir>/*`（tool=fs.write），审计含
  「整批模式一次授予」。
- 非聚合卡走 pattern 门 → 放行调用但 **0 行** + 拒因入审计（`TestTicket49PatternGrantOnANonBatchCardWritesNothing`）。

### 格 6 — Audit（grant_id 关联 SQL 读数）

实现：224 已落 `tool_call.grant_id` 列与 `DecisionAllowGrant` 词表；本票把「每笔 grant-backed 执行」
的逐调用覆盖补齐（input.type 降级分支同样走 `book(... grantID)`）。

读数（真库 SQL）：
- `TestTicket49GrantUseBooksGrantIDOnEveryRow`（tools）：两笔被授权调用 → **两行** tool_call，
  每行 `decision=allow_session_grant` 且 `grant_id=31`（逐调用，非每 grant 一行）。
- `cmd/wisp`（`TestTicket49SessionEndSweepRevokesEveryLiveGrant` / `...RevokeVerb...` /
  `...TargetBinding...`，PASS）：
  `SELECT decision, grant_id FROM tool_call WHERE task='t49-…'` →
  被授权行 `(allow_session_grant, <grant 行 id>)`、未授权行 `(<非 grant 判定>, NULL)`。

## 2. 契约与禁区

- **未改任何 C-## 契约文字**、未动 D43/D45 文本、未动 PLAN/specs/thresholds/golden/allowlist/design。
- `tools.Decision` 新增 `TargetProcess`（描述性事实，非授权通道）；票 90/224 的 reflection 钉
  （Decision 无 allow/grant 字段）继续绿——名字不含 allow/grant/approve，票 148 白名单把该字段记入
  panelVerdictDroppedFields（面板不携出）。
- ticket 224 的「Bridge 恰一枚 grant 形字段」钉（`TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter`）
  继续绿：input.type 读侧走**同一** `b.grants` 源（接口断言发现 TargetCovering 能力），无第二通道。
- `frontend/**`、`docs/contracts/**`、C17 白名单：零触碰。
- 新协程：零新增（test goroutine 走 fakes_test 的 spawn；生产侧无新 goroutine）。

## 3. 自测读数（全部现跑）

| 面 | 命令 | 读数 |
|---|---|---|
| session | `go test ./internal/session/` | ok（含 8 枚 TestTicket49*） |
| tools | `go test ./internal/tools/` | ok（含 12 枚 TestTicket49*，票 224/90/148 钉全绿） |
| approval | `go test ./internal/agent/approval/` | ok（含 4 枚 TestTicket49*，148 白名单按其自述流程同步） |
| memory | `go test ./internal/memory/` | ok（未动，回归） |
| cmd/wisp 票 49 | 载具 `-test.run=TestTicket49` | ok：5 枚全 PASS |
| cmd/wisp 整包 | 载具全跑 | 见 §4（已知死法区分） |

载具规约遵守：一次一发、跑前 `Get-Process wisp.test` 现查（无则发）、`D:\work\soft\Git\bin\bash.exe`、
未裸跑整包（诊断性单测裸跑遭遇 0xc0000135 onnx 绑定死法即改载具，未据其下任何结论）。

## 4. 整包死法登记（照派单 §14 预案）

整包载具跑按 A1893/A1882 处置：若 WebView2 区段 exit 139/绑定死法命中，以票 49 五枚 + 票 224 三枚 +
票 148 一枚的单跑复量为准，如实写——本表以单跑读数为准，整包读数只作回归参考，不把既有死法记到本票头上。

## 5. 未完成 / 边界（如实）

- 面板侧 grants.list/grants.revoke 入站方法**未建**（SPEC-08:171 的面板路由）——按
  `internal/panel/l2_grant_boundary_test.go` 的 banned-root 表，Go 侧入站 grant 路由会红；操作者出口今天为
  console 动词（`grants` / `revoke-grant <编号>`），与票 224 的处置一致。这是登记，不是省略。
- 批量 opt-in 的「一张卡」读数使用队列型确认（PendingApproval）作为答案面；3 秒 L1 窗口自身无应答动词，
  这是本仓 D45-1 形状（窗口到点即执行）。opt-in 授权的答案面落在队列卡上——与现树一致，未新造通道。
- `resolveTargetOverride`（runSpec 测试侧解析覆盖）仅测试填；生产恒 nil，真机走 SystemBridge 同源解析。
