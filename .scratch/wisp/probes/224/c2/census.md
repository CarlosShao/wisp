# 224-c2 只读普查：票 224 读格「本机绿 / CI 红」的边界到哪为止

代号 `224-c2`（`.scratch/wisp/probes/224/` 起手现量只有 `c1`／`r2`／`v1`／`v2`，`c2` 未被占用，复认）。
只读腿：零 `go test`／零 `go build`／零 `go vet`／不执行任何 exe／零写入他人路径。
本文件每一条读数都来自**已落盘的 CI 日志**或**源码**（都带全路径＋行号），
凡"必须真跑才答得出"的一律写「待落地腿自量」并把尺配好（§3），不硬凑。

**一句话结论**：这一格在 CI 上红的直接原因**不是会话授权那一跳**，而是被装配的 `wisp run`
在 runner 上把一次「授权目录内、文件还不存在」的 `fs.write` 判成了 **L2**——
于是那一发走的是审批队列（20 秒、到点自动拒绝），而 grant 那一支在源码里**只挂在 L1 分支上**，
连被查询的机会都没有。授权那一跳的接线、身份、落盘、读回在 CI 上都有正面读数，没有一处缺件。

---

## §0 起手锚

### 0.1 两次取数（共享树里 HEAD 与在飞名单都会漂，所以都带时刻）

| 时刻 (+08) | HEAD | `git status --porcelain cmd internal` |
|---|---|---|
| 08:54:13 | `0c1ec3a2`（`ledger(A512)` 08:51:42） | 13 条：`M cmd/wisp/panel_inbound.go`、`M cmd/wisp/run.go`、`M internal/panel/{bridge,composer,composer_dispatch,composer_dispatch_test,l2_grant_boundary_test,pump}.go`、`?? cmd/wisp/panel_config_store.go`、`?? internal/ball/sta_release_windows_test.go`、`?? internal/config/settings.go`、`?? internal/config/settings_248_test.go`、`?? internal/panel/config_handlers.go` |
| 09:10:02 | `530ba16f`（＝本枚骨架发，父 `8b050774`） | 4 条：`M cmd/wisp/panel_config_248_test.go`、`M internal/ball/{ball_windows,sta_windows}.go`、`?? internal/ball/sta_release_windows_test.go` |

⇒ 中间那一档（09:0x）**248-r1 把自己那批发进去了**（`0d87a681 ticket248(248-r1 产码)`），
所以我 08:54 看到的 `cmd/wisp/run.go`／`internal/panel/*` 脏态已经不是脏态。
连带两件事写在这里：
1. ⚠ 本枚对 `cmd/wisp/run.go` 的行号**一律取 `git show 8ae4c23e:cmd/wisp/run.go`**（＝CI 那一发的版本），
   不是工作树版本；工作树/HEAD 的 run.go 比它多 `+19/-2`。
2. `internal/panel/l2_grant_boundary_test.go` 在 08:54 的脏态**不是越界改动**：
   `0d87a681` 的正文逐字写了「冻结件 l2_grant_boundary_test.go 只动 A487 具名解冻的 :2051/:2139 两枚锚」。
   本枚没有读那两枚锚的内容（不在射程），只登记"它被具名解冻过"这一事实。
   ⛔ 本枚对该文件零写入。

推送区间复认：`git merge-base --is-ancestor` ⇒ `0589fd9c`、`8ae4c23e` 都是 HEAD 祖先；
`git rev-list --count 0589fd9c..8ae4c23e` ＝ **231 枚**（08:54 时刻）。
`git rev-list --count origin/dev..HEAD` ＝ **21 枚**（09:10 时刻，含本枚）。

### 0.2 题面那句「这枚文件是本次 push 新带进 CI 的」＝复认为真

- `git log --diff-filter=A -- cmd/wisp/ticket224_assembly_test.go` → 只有 `1c601fab`
  （`test(224-r2 N#1+N#2)`，2026-09-30 13:16:03 +0800），且 `1c601fab` 在 `0589fd9c..8ae4c23e` 区间内；
- `git cat-file -e 0589fd9c:cmd/wisp/ticket224_assembly_test.go` → **ABSENT**。
- 该文件在工作树与 8ae4c23e 逐字节同长（574 行，`git diff --stat HEAD -- <该文件>` 空）
  ⇒ 题面引的行号与 CI 那一发是同一份文件。

### 0.3 复认题面三行：内容逐字对，行号各差一枚；而且红名册是**七行**不是三行

源码侧（`cmd/wisp/ticket224_assembly_test.go`）：`:315` 是条件行、`:316` 才是 `t.Errorf` 行；
`:318` 条件、`:319` 报错；`:322-323` 条件、`:324` 报错。⇒ 编排者引的 `:315/:318/:322-323`
是**条件区间**，Go 报的是**报错行**（见 §6-U1）。

CI 侧那一枚用例的**全部** `file:line:` 行，出处
`.scratch/wisp/probes/orchestrator/ci-delta-1/logs_new_all.txt`（test-windows job 原始日志，
`--- FAIL` 在 `:3418`，耗时 41.72s 逐字在同一行）：

| 日志行 | 报出的源码行 | 源码里的句子 | 题面提到 |
|---|---|---|---|
| `:3408` | `cmd/wisp/ticket224_assembly_test.go:269` | `the control call errored: the L1 window running out means execute` | ⛔ 没提 |
| `:3409` | `…:293` | `the granted call errored: true` | ⛔ 没提 |
| `:3410` | `…:311` | `a live session grant left 1 cards on screen for the covered call: 「命中授权就不再弹卡」 is not happening in the assembled run (ledger id=sess_1a0d59d3b4b3afed9833b3516fd42598, pattern="C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking81779933\\002\\granted-in-this-session.txt", row=1)` | ⚠ 题面说"没响" |
| `:3411` | `…:316` | `covered call booked decision="timeout", want "allow_session_grant"` | ✔ |
| `:3412` | `…:319` | `covered call booked grant_id=<nil>, want the covering row 1: …` | ✔ |
| `:3413` | `…:324` | `audit is missing the GRANT-USE line naming row 1; got:` | ✔ |
| `:3416` | `…:327` | `the covered call should have written C:\Users\RUNNER~1\…\granted-in-this-session.txt: … cannot find the file specified` | ⛔ 没提 |

（日志 `:3414`/`:3415` 是 `:324` 那句 `%s` 里带出的两行审计原文，不是独立断言：
`[audit] wisp run: SESSION-MINT id=sess_1a0d59d3b4b3afed9833b3516fd42598`、
`[audit] session: GRANT-RECORD id=1 tool=fs.write pattern="C:\\Users\\RUNNER~1\\…" session=sess_1a0d59d3b4b3afed9833b3516fd42598 expires_at=1793463402`。）

**同一发里另外三枚 224 用例的读数**（口径：只认 `--- FAIL`）：
`--- PASS: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (1.50s)`（日志 `:3397`）、
`TestTicket224ProductionSessionDoesNotSurviveRestart` 不在 CI 红名册里
（`.scratch/wisp/probes/orchestrator/ci-delta-1/fails_new.tsv` 全 **26 行**（`awk END{print NR}` 现量，10-02 09:2x）
只有本枚那一条 224）。

### 0.4 本机侧读数（都是已归档的日志，本枚没有复跑）

- `--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.20s)`
  ＝ `.scratch/wisp/probes/33/r9/09-fullpack-1.txt:1325`（日志内嵌 2026-10-01T17:55:28→31）
- 同枚 `(3.17s)` ＝ `.scratch/wisp/probes/33/r9/11-fullpack-2.txt:1324`（2026-10-01T18:04）
- 同枚 `(3.23s)` ＝ `.scratch/wisp/probes/245/v1/gate-cmdwisp-full.txt:1095`（2026-09-30T17:58）
- 同枚出现在 `.scratch/wisp/probes/246/v2/roster-cmd-wisp-green.txt:219`（PASS 名册）
- 本机 temp 根逐字＝`C:\Users\swq\AppData\Local\Temp\…`（同一批日志，`swq` 只有三枚字母
  ⇒ 该卷不会给它生成 8.3 短名，两形天然相等）

⇒ 本机三发 3.17/3.20/3.23s，CI 一发 41.72s。**这枚用例只有本机半边被验过**这一句成立，
但成立的比题面更远：见 §5-J3 与 §6-U5。

---

## §1 「命中授权」那一跳在这台装配机上有几环

按**环**列，每环给 `file:line` ＋「今天生产码里这一环的调用者枚数」。
行号口径：`cmd/wisp/run.go` ＝ `git show 8ae4c23e:cmd/wisp/run.go`；其余文件工作树＝HEAD＝CI 那一发。

| # | 环 | 生产码位置 | 今天有几个真调用者 | 本机/CI 现量 |
|---|---|---|---|---|
| H1 | **铸身份** | `internal/session/session.go:94-100` `Mint()` ＝ `crypto/rand` + `hex`，⛔ 不读任何 env（`grep os.Getenv internal/session/` 非测试＝**0 命中**）；装配点 `cmd/wisp/run.go:447`，失败落 `:449-451` SESSION-MINT-FAILED | 1（`cmd/wisp/run.go:447`） | CI 现场 `sess_1a0d59d3…`（logs_new_all.txt:3414）＝铸成功 |
| H2 | **建 ledger ＋ 宣布** | `cmd/wisp/run.go:452-456`（`Store: mem`、`Logf: rt.auditf`）、`:464` 打 `SESSION-MINT id=` | 1 | CI 打了那行 ✓ |
| H3 | **产出（谁写行）** | 真点击路：`internal/agent/approval/gate.go:647-648` `NativeAPI.AllowSession` → `:652 allowSession` → `:678 g.grants.Record(ctx, tool, p)`（卡片每印一条路径一行）。入口只有两枚：CLI 那个词 `cmd/wisp/approval_reply.go:259`，与 `internal/agent/approval/replies.go:351→:363`。**面板侧不存在**：`replies.go:347` 逐字"There is deliberately no PanelAllowSession" | **1 条真入口在跑**（CLI reply 流），面板/球侧 0 | 第一枚用例就是走 H3 真路（`:119` 写 `"session "+corr` 进 `runSpec.reply`）→ CI `PASS 1.50s` |
| H3' | 测试侧种行（**绕过了 H3**） | `cmd/wisp/ticket224_assembly_test.go:280` 直接 `rt.session.Record(ctx, "fs.write", c)`；`c` 来自 `:275 rt.paths.Canonicalize`（生产规范化器，不是手拼） | — | CI 打出 `GRANT-RECORD id=1` ✓ |
| H4 | **落盘** | `internal/session/grants.go:157 store.InsertGrant` → `internal/memory/dao_misc.go:17`（唯一写者；`dao_misc.go:21/:27-29` 拒空 pattern／拒零 expires_at） | 1 | 行确实落了：CI 后续 `:330 ListGrantsBySession` 读得回（`:334` 那条断言**没响**） |
| H5 | **投递（注入缝）** | `cmd/wisp/run.go:547-552` `grantRead/grantWrite`（`:544-546` 逐字讲 typed-nil guard）→ `:592 Grants: grantWrite`（approval.Options，`internal/agent/approval/gate.go:59/:88/:154`）／`:726 Grants: grantRead`（tools.Options，`internal/tools/bridge.go:102/:149/:194`） | 1 个装配点，两侧各 1 次消费 | 两侧都在场（`GRANT-RECORD` 由 ledger 的 logf 打出 ⇒ H5 的写侧接上了） |
| H6 | **查（读侧唯一入口）** | `internal/tools/bridge.go:334-342`：`if !sil.Silenced && sil.Level == risk.L1 { grantID = b.sessionGrantID(…) }` → `internal/tools/grant.go:63-97`（`:78-92` 逐条 `b.canonical` 后才是可引用形状，`:94` 调 Covering）→ `internal/session/grants.go:189-231 Covering`（`:198 ListGrantsBySession`、`:238-243 grantLive`、`:322-337 patternCovers`＝`path.Match` 只给带 `*` 的行） | `Covering` 生产调用者 **1**（`grant.go:94`）；`sessionGrantID` **1**（`bridge.go:335`） | ⚠ **CI 上这一环从没被够到**：`logs_new_all.txt` 里既无 `GRANT-HIT`／`GRANT-USE` 也无 `GRANT-SKIP`／`GRANT-READ-FAILED`——`Covering` 未命中不打日志（`grants.go:226-228` 直接 return），所以"没日志"本身不能区分，但 H0（下表）说明那一发根本没进这个 `if` |
| H7 | **判定入账** | `internal/tools/bridge.go:417-472 route()`；L1 分支 `:427-436`（`grantID != 0` ⇒ `DecisionAllowGrant`，常量在 `internal/agent/journal.go:32`）；L1 拿到 `AnswerTimeout` 写的是 `DecisionAllow`（`:439-442`）；**`agent.DecisionTimeout`（"timeout"）在全仓只有一个写入点＝`:457-460` 的 L2 分支** | — | CI 读到 `"timeout"` ⇒ 那一发在 **L2 分支** |
| H8 | **取证列** | `internal/tools/bridge.go:1066-1123 book()`：`:1069-1073` 那行 `tools: call … risk=… decision=… in_allowlist_scope=… grant_id=…`；`:1116-1119 if grantID != 0 { tc.GrantID = &g }`；列本身冻结在 `internal/memory/schema.go:83`（`bridge.go:1058-1065` 逐字讲了这段为什么是参数不是字段） | 1 | CI `grant_id=<nil>`＝H6/H7 没给值，不是 H8 断 |
| H9 | **审计句** | `internal/tools/bridge.go:337` 的 `tools: GRANT-USE …` 在 H6 命中之后才打 | 1 | CI 缺这行＝H6 没命中 |
| **H0** | **等级从哪来（今天真正断的那一环）** | `internal/tools/bridge.go:291` `Assess(…)` → `internal/risk/rules_gateway.go:32-52 rulePathAllowlist`：`:37 Canonicalize` 失败 ⇒ R2（无法规范化措辞），`:45 InAllowlist(canonical)` 假 ⇒ R2「目标路径在授权目录之外」⇒ **L2**（`internal/risk/assessor.go:16/:32/:77`）。而 `InAllowlist` ＝ `internal/tools/paths.go:133-163`，**两次包含**：`:140` 词法、`:152-155` 还要 `resolvedForm()`（`paths.go:247-267`，`filepath.EvalSymlinks`）；根是 `paths.go:53-96 NewPathCanonicalizer` 用 `risk.Resolve` 规范化过的 | `InAllowlist` 生产调用者 **4**：`rules_gateway.go:45`、`bridge.go:1145 inScope`、`internal/tools/task.go:838`、接口声明 `assessor.go:141` | ⬅ 本枚判死的正是这一环 |

### 1.1 「谁产出／谁投递／谁落盘」三处读完之后的那句要紧话

H3 那一列就是这条规矩的用场：`approval_grant` 行的**生产产出者今天只有 CLI 那一枚词**
（`approval_reply.go:259`），面板/球侧 0 枚（`replies.go:347` 是设计而非遗漏）。
所以读格用 `:280` 手工 `Record` 种行**不是**"测试替生产做了中间一环"——
被它替掉的只有 H3（点击→行），而 H3 由同文件第一枚用例走真路单独覆盖（且 CI 绿）。
读格没有替掉的也没有 H5/H6/H7：它注入的是 `rt.session`、查询走的是 `rt.bridge`，中间零 stub。

⚠ 反过来，这台装配机上**「一次写在授权目录内的调用」今天到不了 L1** 这件事，
是 H0 而不是 224 的环——见 §2。

---

## §2 本机绿／CI 红的候选，逐支配尺

先给一刀判死全部候选的那根**公共尺**（纯读码＋已落盘 CI 日志，⛔ 不需要跑）：

> **R-0（读码尺）**：`decision="timeout"` 这个列值在 `internal/tools/bridge.go` 里只有一个写入点——
> `:457-460` 的 `case risk.L2: … case AnswerTimeout: dec.DecisionColumn = agent.DecisionTimeout`。
> L1 分支 `:437-442` 拿到 `AnswerTimeout` 写的是 `DecisionAllow`（逐字注释
> "The L1 window running out unopposed MEANS EXECUTE (SPEC-06 §2)"）。
> CI 那发读到 `"timeout"` ⇒ **它走的是 L2 审批路由**。

三条独立佐证（同一发日志，本枚现量）：

- **时长**：`cmd/wisp/approval_reply_201_test.go:103-162` 这个载具把 `l2Wait=20s` 写成
  `[risk] confirm_timeout_sec = 20`（`:139-142`），而那个键**是 L2 审批超时**
  （`internal/config/schema.go:449-450` 逐字 "how long an L2 confirmation card stays open"）；
  L1 窗口另有其键 `risk.l1_window_sec` 默认 `2`（`schema.go:451-453`），且被
  `internal/agent/approval/gate.go:137-145` 硬夹进 `[MinL1Window=2s, MaxL1Window=3s]`
  （`internal/agent/approval/queue.go:120-122`）⇒ **L1 路由物理上不可能产生 20 秒等待**。
  时间线闭合：开机完成 `16:16:22.498`（logs_new_all.txt:3406）→ H3' 落行的 `created_at`
  ＝ `expires_at 1793463402 − 30d`（`internal/memory/retention.go:36 GrantAuditTTL`，
  行写于 `grants.go:163`）＝ **`2026-10-01T16:16:42Z`**（本枚 `date -u -d @1790871402` 现算）
  ＝ 开机后 19.5s ＝ 对照那发等满一发 20s；再 `16:17:02.589`（`:3417`）＝覆盖那发又等 20.6s。
  41.72s ≈ 开机 + 2×20s。本机 3.2s ＝ 开机 + 一枚 2s 的 L1 窗口。
- **对照那发自己也红**：`:269` 响＝**没种授权的对照调用也 errored**。L1 到点是 allow（不该 err），
  L2 到点是 auto-reject（一定 err）⇒ 又一处 L2 指纹。
- **写出文件缺失**：`:327` 响＝那一发从未执行 ⇒ 与 L2 auto-reject 同向。

下面逐支。

### (a) 会话身份来源依环境而变 —— **判死（推翻）**

- 代码：`internal/session/session.go:94-100` 的 `Mint()` 只有 `crypto/rand.Read` + `hex`，
  没有 `COMPUTERNAME`／`USERPROFILE`／用户名／profile 路径／`%APPDATA%` 任何一处
  （`grep -rn "os.Getenv|USERPROFILE|COMPUTERNAME|APPDATA|TempDir" internal/session/` 非测试文件 **0 命中**）。
- CI 反证：`pattern`/`session` 两枚键在同发日志里逐字相等（`sess_1a0d59d3b4b3afed9833b3516fd42598`
  既出现在 `:3410` 的 ledger id、也出现在 `:3415` 的 GRANT-RECORD 的 `session=`），
  且**`:330-336` 那条"按铸造 id 读得回恰好一行"的断言在 CI 没响** ⇒ 身份与读回都对得上。
- 权限档也不依环境：`perm.New` 只有 `Manager` 一个来源（`internal/perm/store.go:123-147`，
  `:124-126` 无 Manager 直接报错），装配点 `cmd/wisp/run.go:628-632`，读的是本 run 自己那份
  `config.toml`（载具写在 `h.dir` 里）；`ask_every_step` 默认值不改判级
  （`internal/risk/mode.go:168-169` 原样返回 `d.Level`）⇒ 模式不是本机/CI 的差异源。
- ⇒ 这一支**不需要落地腿再跑**就可以划掉。

### (b) 「落盘根／拼法在 CI 上是另一棵树」—— **拆成两半：DB 那半判死，路径拼法那半是根因（判死）**

- **DB 半（判死）**：行就写在 `h.dir`，被 `memory.Open(h.dir)` 与 `h.openStore()` 按铸造 id 读回
  （`:330` 那条断言未响，见上）。⇒ 题面 (b) 里"grant 写进 A 根、查询按 B 根"这一形**不成立**。
- **拼法半（＝根因，本枚判死）**：runner 的临时根是 **8.3 短名拼法**
  `C:\Users\RUNNER~1\AppData\Local\Temp\…`（CI 日志里每一枚 temp 路径都是这形；
  本枚用例的 `pattern` 逐字含 `RUNNER~1`，logs_new_all.txt:3410/:3415）。机制三行：
  1. `internal/risk/pathresolver.go:105-140 Resolve`：`:127-135` 若 `resolveHandle` 成功
     （`internal/risk/pathresolver_windows.go:25-52`，`CreateFile(OPEN_EXISTING)` +
     `GetFinalPathNameByHandle(VOLUME_NAME_DOS)`，`:21-23` 逐字"This step inherently expands
     8.3 short names"）⇒ Canonical ＝**长名**；
     `:136-139` 若路径**还不存在** ⇒ 纯词法兜底，Canonical **保留输入拼法（短名）**。
  2. 授权目录根 `…\002` 是 `t.TempDir()` 建出来的**存在目录** ⇒ `NewPathCanonicalizer`
     （`internal/tools/paths.go:53-96`）把它记成**长名**；
     而被判的那发目标文件 `granted-in-this-session.txt` **此刻还不存在**
     （它正是那一发要写出来的东西）⇒ `Canonicalize` 回**短名**——
     CI 现场直接证到这一条：`:3410` 的 `pattern` 就是短名，而它正是 `:275 Canonicalize` 的输出。
  3. `InAllowlist`（`paths.go:133-163`）把两根形做**与**：`:140 rootsContain(roots_LONG, fold(target_SHORT))`
     ＝假 ⇒ `rules_gateway.go:45-50` 出 `R2: 目标路径在授权目录之外` ⇒ **L2**。
- **同发正面读数（不是本枚推的）**：`TestPathResolverShortNameAListDenied` 在同一发 CI 红在
  `internal/risk/pathresolver_junction_windows_test.go:135`，逐字
  `got "C:\\Users\\runneradmin\\AppData\\Local\\Temp\\…\\.git-credentials"
  want "C:\\Users\\RUNNER~1\\…\\.git-credentials"`（logs_new_all.txt:4064）
  ⇒ **同一台机器上 `Resolve` 确实把短名展开成长名**，而那枚测试的期望串（同样来自 `t.TempDir()`）
  停在短名。该文件 `:104-107` 还逐字写着"this case passes on a dev box and fails on the
  windows-latest runner"——这个形状在仓里**不是新事**。
- **仓里已经记过这条环境差异**：`internal/risk/pathresolver.go:283-300`（ticket 72 那段）
  逐字"On GitHub's windows-latest runner USERPROFILE is the 8.3 short form (C:\Users\RUNNER~1)
  while GetFinalPathNameByHandle hands back the long one (C:\Users\runneradmin)"。
  ⚠ **但那次修复只给了 A/B 黑名单那一侧**（`pathForms.raw/.real/.certain`，`:302-324`），
  `[fs] allowed_dirs` 这条腿（`internal/tools/paths.go:133-163`）**没有同一套双形比较**——
  这就是缺口的位置。判死：本枚读到的是"缺一条腿"，不是"缺授权那一跳"。

### (c) 时序（等满某个 deadline）—— **判死：对得上现成常量，且对的是 L2 的常量**

见上面 R-0 的第二条佐证：41.72s ＝ `2×confirm_timeout_sec(20) ＋ 开机/收尾`。
⛔ 不是"墙钟玄事"，也不需要凑：本仓那族 `300.0x s` 的读数（`confirm_timeout_sec` 的
default 300，`internal/config/schema.go:450`）在同一发里就有现成现场——
`cmd/wisp/run_test.go:378`（见 (e) 的第三个 witness），
区别只是载具写没写那枚键。⇒ 「如果 41.72s 对不上任何现成常量要具名写对不上」这一条**没被触发**：它对上。
⚠ 但**方向**题面猜反了：它不是"那一枚按钮没被命中所以在等"，而是"那一发压根不经过等 2 秒的那条路"。

### (d) 测试自带的载具在 CI 上走不通 —— **作为原因判死；作为形状成立但归到 (e)**

- 载具不假设有人在按确认：读格**没有**接 reply 源（`cmd/wisp/ticket224_assembly_test.go:233-295`
  全程不碰 `h.reply`，而 `approval_reply_201_test.go:91-93` 逐字写"nil means this host wired no
  answer source"），它依赖的是**「L1 窗口到点＝执行」这条极性**（SPEC-06 §2），
  在无人应答下本来就该绿。
- 载具不假设任何 `WISP_*` 环境变量：全仓生产码只读 `WISP_ENV`
  （`cmd/wisp/slo_windows.go:251`、`internal/buildinfo/{env.go:36,buildinfo.go:38}`）与
  `WISP_MODELS_MANIFEST`（`internal/models/manifest.go:325`），都不在这条链上。
- 票 123 那一族"假设有人按确认"的形状在这一枚上**不适用**（它连 answer 侧都没装）。
- ⚠ 题面 (d) 的**内核**其实指对了东西，只是位置错：载具确实"假设了等级是 L1"，
  而这一枚假设在 runner 上塌了——塌在 H0，不塌在载具。（载具的另一处脆弱点在 §3 的尺里要用到：
  它的非空判定是 `cards==1`/`cards==0` 的**弹卡计数**，计数分不出 L1 窗口和 L2 审批卡——
  CI 那发 `control.cards==1` 是**因为 L2 也弹了一张卡**才恰好通过的。）

### (e) 真缺陷：那一跳在生产装配里根本没接／本机绿是测试替生产做了中间一环 —— **判死为"不是这一支"，但结论比题面更硬**

- 接线齐：H5（`run.go:547-552`＋`:592`＋`:726`，typed-nil guard 逐字在 `:544-546`）、
  H6（`bridge.go:334-342` → `grant.go:63-97` → `grants.go:189-231`）、
  H7/H8（`bridge.go:427-436`、`:1116-1119`）都在生产码里，且 224-r2 的 commit 正文
  （`1c601fab`）逐字记过 M4/M5/M1 三枚突变各自把哪一枚打红（"删掉 Record 调用 ⇒ 本枚红"、
  "`bridge.go:1116 的 if grantID != 0 改 if false` ⇒ 真列那枚红"），
  说明那一支不是空转。⇒ 题面"如果成立最要紧"的那一支**在这一发不成立**。
- 但**"CI 半边从没验过"比题面说的更大**：在 runner 的短名 temp 形状下，
  `cmd/wisp` 这台装配机上**没有任何一发调用能走到 L1 分支**（H6 的 `if` 在 L2 上压根不执行）。
  所以票 224 读格的 CI 分母是 0，不是 1。
- 三个**独立、早于/外于票 224**的现场（本枚从已落盘日志逐字读）：
  1. `cmd/wisp/run_test.go:378`（`TestComposedGateBlocksAWriteForTwoSeconds`，
     该用例 `:384-388` 断言"2 秒 L1 窗口真的挡住"＋"cards==1"）——
     **基线发（`0589fd9c`）与新发两发都红**：`--- FAIL … (301.12s)`（logs_base_all.txt:3709）/
     `--- FAIL … (301.70s)`（logs_new_all.txt:3099），报错逐字
     `an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`；
     本机同码绿：`--- PASS … (3.46s)`（`.scratch/wisp/probes/33/r9/09-fullpack-1.txt:1015`）。
     300＝该载具没写 `confirm_timeout_sec` ⇒ default。
  2. `cmd/wisp/run_mode101_test.go:506`（`TestTicket101SessionGrantDoesNotCrossRestart`，
     也是**两发都红**）：逐字
     `control write should NOT error under auto_approve (审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝);
     the control half is broken, so the refusal above proves nothing`（logs_new_all.txt:3010）。
     这一枚尤其值钱：`auto_approve` 本该把授权目录内的 L1 **静默放行**（`internal/risk/mode.go:174-176`），
     它却撞上了一张 L2 卡——因为 **R2 在"任何模式都不许静默"的红线名单里**
     （`mode.go:36-41` 逐字列出 R2）。⇒ 同一条"根形≠目标形"在三枚不同断言上各撞一次。
  3. 本枚读格自己（§0.3 七行）。
- ⇒ 结论落点：**票 224 的实现格没被判为缺陷；被判为缺陷射程的是 H0 那一环（`InAllowlist` 无双形比较），
  它是票 18/72/105/107 那一族授权目录路径腿的账，不是票 224 的账。**
  "AC 有没有勾早"这一问只有 J3 那一句可说（见 §5）。

### 2.1 本枚自破的尺（写进结论前自己跑过一遍的代价）

- 我第一把找"8.3 短名覆盖"的尺是 `grep -rl "GetShortPathName\|RUNNER~1\|8.3" --include=*_test.go`，
  它把 `internal/tools/grant_test.go` 也报了进来——那里命中的是 **`SPEC-06 §8.3`**（条款号），
  不是 8.3 短名。换成 `grep -rn "GetShortPathName\|短名\|short name"` 后分母从 15 枚变 10 枚，
  且 `internal/tools/**` 里真讲短名的只有 `bridge_junction_windows_test.go`（15 次）。
  ⛔ 结论：**"文件名里含 8.3"这把尺在本仓必指错**（§8.3／§18.3 这类条款号会冒充版本号）。
- 另一把：`git status --porcelain` 在 08:54 报 `M internal/panel/l2_grant_boundary_test.go`，
  而 `git diff` 对该文件**空输出**（`git ls-files --eol` ＝ `i/lf w/lf attr/text eol=lf`）——
  两小时后该脏态整体消失（248-r1 提交）。⇒ "porcelain 报了 M"不等于"内容变了"，
  短引状态必须配 `git diff --numstat` 复核，否则会把 EOL/竞态写进结论。

---

## §3 落地腿 224-r3 的派单料（确切尺＋本机可判红的载具）

⛔ 本枚一律没跑。下面每一条都写清"跑哪一发／期望读数"。

### 3.0 前置：PATH 那把尺（复认＋一处改法）

- 复认：缺 sherpa DLL 时 `go test ./cmd/wisp/` 在**加载期**死，`exit status 0xc0000135`、
  `=== RUN` 0 条＝**用例根本没跑**，不是绿。权威文本在
  `scripts/wisp-cli-tests.sh:6-21`（四发读数）与 `:101-109`；
  另有 `docs/evidence/s1/121-adversarial-acceptance.md:280`（同一物两种写法：rc=1 的字符串 vs
  3221225781＝0xC0000135）、`docs/evidence/s1/102-adversarial-acceptance.md:148`
  （`rc=1`＋`=== RUN` 0 条的识别法）。
- 现量在位件（本枚只 ls）：`third_party/sherpa-onnx/` ＝ `onnxruntime.dll`、
  `sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`，与 `deps.toml:30/:34/:38` 三枚
  `[sherpa-onnx.dll.*]` pin 一一对应；`build/` 目录也在（含同三枚 dll ＋ `wisp.exe` 等 8 枚 exe）。
- ⚠ **对编排者给的那把尺的一处改法**：题面尺
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 里 `$PWD` 在 Git Bash 下是
  `/d/work/workspace/projects plans/Wisp` 形＝`scripts/wisp-cli-tests.sh:101-109` 逐字认的**可用形**
  （该注释指出的陷阱是 `pwd -W` 的盘符形 `D:/…`，那个才 0xc0000135）⇒ 前一半复认可用；
  后一半 `$PWD/build` 建议**去掉**：`build/` 里躺着 `wisp.exe`／`wisp228.exe`／`wisp77.exe`
  等 8 枚**过期产物**，把它们排到 PATH 前面只可能让"按名字 exec"的测试拿到旧 exe，
  对本用例零收益。⇒ 建议尺：`PATH="$PWD/third_party/sherpa-onnx:$PATH"`。
  ⛔ 本枚**没有**实测这一处（"会不会真有名冲突"要跑才知道，见 J5）；给的是理由，不是读数。

### 3.1 尺 A：单发跑这一枚用例（本机基线，期望绿）

```bash
cd "D:/work/workspace/projects plans/Wisp"
PATH="$PWD/third_party/sherpa-onnx:$PATH" \
  go test -count=1 -v -run '^TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking$' ./cmd/wisp/ \
  > .scratch/wisp/probes/224/r3/a-baseline-local.txt 2>&1; echo "rc=$?"
```
期望：`rc=0`、`--- PASS: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (3.1s ~ 3.3s)`
（三枚已归档读数 3.17/3.20/3.23，见 §0.4），`--- FAIL` **0 条**，`=== RUN` 恰 1 条。
判读口径：只认 `--- FAIL`；`^--- PASS` 与含子用例的 `    --- PASS` 是两枚分母；
`t.Logf` 也带 `file:line:` 前缀（同一发 CI 日志里 `pathresolver_expansion_test.go:95/:98/:101`
三行挂在 `--- PASS: TestC26ExpansionMustNotRewriteOntoAnotherTree`（logs_new_all.txt:4058）底下，
就是这一坑的现成标本）⇒ ⛔ 不许把带 `:行号:` 的行当红。
⚠ 跑之前先看 `.scratch/wisp/probes/*/`＋`git status` 有没有在飞的计时腿（此刻 33-r8/248-r1/250-r1）；
`-count=1` 别省（缓存会把红藏起来）。

### 3.2 尺 B（**决定性、最便宜、不需要任何环境改造**）：正控——人为让那一发判成 L2，看能不能在本机复现 CI 的七行

在**本机的长名 temp 下**（本机默认形），只把载具里 `[fs] allowed_dirs` 指到**别处**
（例：把 `approval_reply_201_test.go:137` 那行 `allowed_dirs` 换成另一枚 `t.TempDir()`，
或把 `ticket224_assembly_test.go:236` 的 `granted` 挪出 `h.dir`），
跑尺 A 同一条命令。
期望读数（如果 §2 判死成立）：
`--- FAIL … (41.x s)`（±1s），报错行**恰好是 §0.3 那七行**
（`:269`/`:293`/`:311`/`:316`/`:319`/`:324`/`:327`），且 `h.err` 里**没有** `tools: GRANT-USE`、
**没有** `session: GRANT-HIT`、也**没有** `GRANT-SKIP`／`GRANT-READ-FAILED`。
⇒ 这一把尺一刀切开"授权那一跳坏没坏"与"等级判没判错"：它**不碰任何授权代码**就复现全部七行。
⚠ 这是探针不是修法（⛔ 不许把 allowed_dirs 改掉提交）。
另配一枚反向正控：同载具里把 `granted` 文件**在种授权之前先写出来**（照
`internal/tools/grant_test.go:196-199`／`:247-253` 的既有形状——单元版那两枚用例都是**先建文件再种行**：
`:194 TestTicket224LiveGrantStopsTheL1Question` 与 `:245 TestTicket224PartialPathCoverageStillAsks`），
在尺 C 的短名 TMP 下应当**不再复现**；若"文件存在"能翻绿，就一刀钉死"存在与否"是自变量。

### 3.3 尺 C（本机模拟 runner 的 `%TMP%` 形状＝让常驻判据在本机判红）

本仓既有定式：常驻判据必须能在本机判红，不许挂〔仅本机可量〕。载具＝把临时根换成**同一棵树的 8.3 短名拼法**。

```bash
# 0) 前提：该卷必须开着 8.3 短名生成（runner 上实测开着：同发日志里
#    TestPathResolverShortNameAListDenied 没 skip，而是拿到了短名并展开了）
fsutil 8dot3name query D:
cmd //c "dir /X D:\\work\\workspace"          # 看真实短名列
# 1) 造一枚长名（>8 字符，必得短名）目录，读出它的 8.3 拼法
mkdir -p /d/wisp224c2-shortharness-root
cmd //c "dir /X D:\\" | grep -i wisp224
# 2) 用短名拼法当 TMP/TEMP，跑同一枚用例
TMP='D:\WISP224~1' TEMP='D:\WISP224~1' \
PATH="$PWD/third_party/sherpa-onnx:$PATH" \
  go test -count=1 -v -run '^TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking$' ./cmd/wisp/ \
  > .scratch/wisp/probes/224/r3/c-short-tmp.txt 2>&1; echo "rc=$?"
```
期望：`--- FAIL … (41.x s)` ＋ §0.3 那七行；且 `h.dir`/`pattern` 逐字带 `WISP224~1`
（＝CI 日志里 `RUNNER~1` 的同位物）。
反控（同命令、TMP 换回长名拼法）：回到尺 A 的 `--- PASS (3.1~3.3s)`。
交叉 witness（同载具下顺带验，成本 300s，建议只跑一次或另发）：
`-run '^TestComposedGateBlocksAWriteForTwoSeconds$'` 应当在本机翻成
`run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认）…`
＋`(301.x s)` ＝ CI 两发读数的逐字复刻。
⚠ 若第 0 步量到该卷 8.3 关闭：⛔ 不许 skip 后当绿——照 `internal/tools/bridge_junction_windows_test.go:77-88`
的既有规矩"**必须红并写明缺什么**"（那文件逐字："改它需要管理员，所以这条用例在关闭了 8.3 的卷上必须红"），
并退回尺 B（它不需要 8.3 就能复现七行）。⇒ 载具分两档：尺 B 证"判成 L2 就复现"，尺 C 证"runner 那台机器就是这么判的"。

### 3.4 尺 D：一刀判死 J1（`InAllowlist` 里到底是哪条腿失配）—— 一次 `-overlay` 探针，零计时

现成先例：`.scratch/wisp/probes/224/v1/overlay.json` ＋ `v1/probe_a_prod_test.go`、
`v2/dialect_probe_test.go`（224-v1/v2 两枚验收腿就是用 overlay 在本机探生产的）。
探针只打印四个数，不断言：
`rt.paths.Roots()`、`Canonicalize(<存在的根>)`、`Canonicalize(<还不存在的目标>)`、
`InAllowlist(<后者>)`；短名 TMP 下期望：roots＝**长名**、存在的根＝长名、
缺失叶子＝**短名**、`InAllowlist`＝**false**。
再把 `paths.go:140` 与 `:152-155` 两条腿分别单独代入（拿 roots／拿 `resolvedForm`），
就知道修哪一条。⛔ 不要靠"看有没有 GRANT-SKIP 日志"判——`sessionGrantID` 在 L2 分支压根不被调用
（`bridge.go:334`），而 `Covering` 未命中**不打日志**（`grants.go:226-228`），两形在同一段日志里长得一样。

### 3.5 派给 224-r3 的一句现成可用尺（不用等 J1）

被装配的那一发**已经**把答案印在自己的审计行里：`internal/tools/bridge.go:1069-1073` 每发调用打
`tools: call … risk=<L1|L2> decision=… in_allowlist_scope=<true|false> grant_id=…`，
外加 `:1086-1091` 的 `tools: PATH-ACCOUNT … roots=%d rewritten=[…] unusable=[…]`。
它们走 `Options.Logf → agentRuntime.auditf`（＝测试里的 `h.err`）。
今天 `cmd/wisp/ticket224_assembly_test.go:206-217 grantLinesOf` 只把含 `GRANT`/`SESSION-MINT`/`REPLY`
的行留下，**`risk=` 那行被自己的过滤器洗掉了**——所以 CI 红名册看不见等级。
⇒ 落地腿只要在失败分支里多打一行（或临时在探针里 dump `h.err`），就能拿到
`risk=L2 rules_hit=[R2] in_allowlist_scope=false` 这条现场；这也顺手回答"CI 那发到底什么等级"，
不需要真跑 CI。⚠ 这是**判据可读性**，不是修法；改不改 `grantLinesOf` 归落地腿与它的验收腿。

---

## §4 撞钉预检（逐枚读断言，不是 grep 新符号名）

**前提结论**：如果 224-r3 去补 §2(e) 那一支（即"承认 CI 走的是 L2"并想让授权在那儿也生效），
它会立刻撞死在下面第一、第二枚钉上；如果它去补真正的缺口（H0 的路径形），
则要过第三～第六枚。**没有一枚钉禁止"把缺失叶子的路径折叠到根的同形"**——那是可以动的方向。

1. **`internal/tools/grant_test.go:301 TestTicket224SessionGrantNeverCoversL2`（非冻结，行为型）**
   两枚断言逐字：`:321-323` "a session grant covered an L2, which SPEC-06 §8.3 bullet 1 forbids"（approval 卡数必须＝1）；
   `:324-328` `grD.asked()` 必须**空**——"the bridge queried the grant source %d times for an L2 verdict;
   **the L2 branch must not reach the check at all**, so a future edit cannot make it answer yes by accident"。
   存在理由：D45-3／PLAN.md:2148「L2 永不进入任何持久授权（含会话级）」，
   并且它特意写成**正控**（`:298-300`："the grant source is told, in so many words, that this exact
   tool and path are covered"）。⇒ 落地腿若把 `sessionGrantID` 从 `bridge.go:334` 那个 `if` 里搬出来
   （哪怕只是"先问再判级"），**就算判定没变也会因 `asked()` 非空而红**。这是最锋利的一枚。
   配套：`:348 TestTicket224SessionGrantNeverCoversDeny`、`:421 TestTicket224PathlessCallIsNeverGranted`
   （`:430-435` 路径为空的调用**连问都不许问**，且 window 必须＝1）、
   `cmd/wisp/ticket224_assembly_test.go:532-538`（重启后不得出现 `allow_session_grant`／`grant_id` 必须 NULL）。
2. **`internal/tools/bridge.go` 内的等级几何**（不是测试，是它的注释合同）：`:322-332` 逐字把
   "grant 只替一张 L1 卡 answering"写成契约，`:411-416` 逐字"grantID … only ever consulted on the L1
   branch. L2 and Deny have no path to it, which is SPEC-06 §8.3's first bullet expressed as a switch
   statement"。⇒ 与第 1 条同向；任何"让 CI 变绿"的改动如果动了这段几何，验收腿按 D45-3 直接判失败。
3. **`internal/tools/ticket90_test.go:432 TestTicket90ModeCarriesNoAllowAuthority`（行为＋反射双型）**
   `:444-452` 扫 `tools.Decision` **每一枚字段名**，含 `allow`/`approve`/`grant`（子串，不分大小写）
   或类型为 `Answer` ⇒ 红；`:454-464` 扫 `agent.ToolRequest`，含 `mode`/`allow`/`approve`/`auto` ⇒ 红。
   存在理由（逐字）："an allow must only ever come from the gate's own answer (SPEC-06 §9)"、
   "the caller of Execute supplies arguments, never a verdict"；
   `internal/tools/grant.go:26-29` 与 `bridge.go:1062-1065` 都点名引用这两行，说明"grant 走参数不走字段"
   就是被这两枚钉钉住的既有形状。
   ⚠ **射程的精确边界（别按题面照抄）**：`SessionID` 这个名字**不在**禁用子串里——
   往 `Decision` 加 `SessionID` 不会触发它，只有 `Grant*`/`Allow*`/`Approve*`/`Auto*`/类型为 `Answer` 才触发；
   同理 `ToolRequest` 那半枚扫描**不含** `grant`。⇒ 落地腿若"顺手"把会话 id 挂到 `Decision` 上，
   **这枚反射钉不会响**，得靠第 4 条那枚行为钉与 §6 的台账规矩拦——这一点题面写宽了（见 §6-U8）。
4. **`internal/perm/ticket90_persist_test.go`（三枚冻结件之一，⛔ 一字不许动）行为型两枚**
   - `:291 TestTicket90ConfigKeyChangesWhatTheChainAsks`：`:300-306` 期望表逐字
     `{ask_every_step, 1, 0} / {ask_high_risk, 0, 0} / {auto_approve, 0, 0}`，
     其中 `w/a` 来自 `:96-129 askOnce` 里那枚 gate spy 的 **windows/approvals 计数**；
     存在理由＝票 83「键不许说谎」：**改这个键，被问的那条腿必须跟着变**（`:94-95` 逐字
     "nothing here reads the mode directly"）。⇒ 任何把"某类调用从 window 改判到 approval"
     或反向的路由改动，都会在这一枚上红，而且是**分腿计数**红，不是"没弹卡"红。
   - `:181 TestTicket90UntouchedConfigStartsAtTheDefault`：`:199` 冷启动档位必须是 default、
     `:203` "cold start #%d **reached an L2 card**" 是 `t.Fatalf`——它同样按**腿**记账。
     ⚠ 这两枚今天能过，是因为 `askOnce` 建的桥**没有 Paths resolver**（`:118`
     `tools.New(tools.Options{Registry, Gate, Modes})`）、工具参数是 `{}`（`:124`）
     ⇒ R2 休眠、判级停在声明的 L1。**如果落地腿"顺手"给 perm 的桥补上 resolver/路径，
     这两枚会因 R2→L2 立刻红**，而那是冻结件——只能改产品码，不能改它们。
   - `:220 TestTicket90SessionGrantDoesNotSurviveRestart`：`:239` 用 `InsertGrant` **手写 fixture 行**、
     `:263/:267/:272` 断新会话读 0 行＋旧行留着；`:274-280` 逐字声明"这里**故意不做**模式断言"
     （AC#3a 与 AC#3b 是两格，不许合并成一条"持久化"用例）。
     它的 fixture session id 用的是 `session-before-restart`/`session-after-restart` 两枚字面量
     （`cmd/wisp/ticket224_assembly_test.go:411-413` 与 `:436-439/:491-493` 是那枚"测试字面量锁"的另一半）
     ⇒ 落地腿若把生产铸出的 id 改成可派生的形状，同时会撞
     `cmd/wisp/ticket224_assembly_test.go:432-435/:484-487` 与票 224-r2 记录的 M1 反控。
5. **`internal/tools/bridge_junction_windows_test.go`（短名形状的两枚既有钉）**
   `:475 TestBridgeRefusesTheRealShortNameOfAnAListFile`——`:492-494` 先要求
   "短名与长名必须解析到同一枚 canonical，否则用例前提不成立"，`:517-519` 要求卡面路径＝**解析后的长路径**；
   `:527 TestShortNameSpellingGetsTheSameVerdictAsTheLongOne`——`:541-558` 对 short/long **两形各跑一遍同一期望**
   （L2＋R2、卡面必须长名、reason 必须含"目标路径在授权目录之外"）。
   ⚠ 关键射程：这两枚用的目标文件**都是先写出来的**（`:484-487`、`:532-535`）⇒ 两形都能 handle-resolve，
   所以它们**盖不到"叶子还不存在"那一形**（＝本枚用例的形状）。
   落地腿若把 `Canonicalize` 改成"能开多少层就折叠多少层、缺失尾巴原样接回"
   （正是 `paths.go:247-267 resolvedForm` 与 `risk` 里 `pathForms.real` 的既有做法，
   `:309-312` 逐字"components below the deepest openable one are re-appended verbatim"），
   这两枚仍然绿——**它们是第一档正控**，改动必须让它们保持绿。
   另外 `:67-88 shortNameOf` 的"拿不到真短名就红并写明缺什么"是本仓对 8.3 载具的既有规矩。
6. **面板侧那枚（题面没列，但补 H0 时可能被牵连）**：`internal/session/grants.go:39-44` 逐字警告
   SPEC-08:171 的 `grants.list`/`grants.revoke` 若实现成**入向 Go 方法**会直接死在
   `internal/panel/l2_grant_boundary_test.go` 的 banned-root 表上（A435 的「落点」条）。
   ⛔ 本枚**没有**读那枚冻结件的内容（只登记这条指路），因为它的行号在 09:0x 被 `0d87a681`
   动过（A487 具名解冻 `:2051/:2139`）——落地腿要引它时**必须重新数行号**，别引旧数。
7. **`cmd/wisp/run_mode101_test.go:390-448`（与 AC#2 极性相反的常驻钉，行为型）**：
   `:400` 用字面量 `const session = "session-before-restart"` 手写一枚**活的**授权行
   （`:409-418 rt.store.InsertGrant`，`ExpiresAt = now+3600` 逐字注 "live by its OWN clock: only the session may die"），
   然后要求那枚 B 档 `.env` 写**仍被拒**；`:427-441` 是票 224 N#3 改过的那半（注释半），
   逐字承认「Ticket 224 made that false: the composition root now injects the read side into the bridge
   (`cmd/wisp/run.go:653 Grants: grantRead`) and the write side into the gate (run.go:529) …
   So the conditional that is actually doing the refusing is this one: `tools.GrantSource` answers
   only for the identity it holds, this row is keyed by a string no production mint can produce」，
   ⚠ 那句注释里自报的两个锚点 `cmd/wisp/run.go:653`／`run.go:529` **已经漂了两档**
   （那两行写于 `ba5db093`，2026-09-30）：CI 那一发（`git show 8ae4c23e:cmd/wisp/run.go`）上是
   `Grants: grantRead` ＝ **`:726`**、`Grants: grantWrite` ＝ **`:592`**；
   到本枚 09:10 的 HEAD（`0d87a681` 之后）又变成 **`:758`/`:618`**。⇒ 落地腿引这枚注释时要重数行号，
   ⛔ 别照抄任何一份——包括本枚这一份（再漂也是正常的）。
   并留了一句给落地腿当心：`:444-448`「What this guard still has to keep true … a B-tier `.env` write
   stays refused when nothing authorized it. ⚠ Read the observation for what it is, though —
   `firstErr` comes from an approval card running out (the L2 route auto-rejects…)」。
   存在理由：**授权梯度不许被"能读回一行"这件事顺手放宽**——它把"命中授权就不再弹卡"的反面钉住了。
   ⛔ 落地腿若为了让读格绿而改变"哪种行能被命中"的规则（例如允许非铸造形状的 id、或让 pattern
   折叠得更宽），第一下顶红的就是这枚；同族另一枚是 `:506`（本枚 §2(e) witness 2，
   `control write should NOT error under auto_approve`），它顶的是"控制半坏了"这一层。
   ⚠ 另：这两枚都在 `cmd/wisp`，所以本枚只读文本、没有读数；它们的当前绿／红名册见 §5-J4。
8. **身份形状锁（顶红"用派生 id 蒙过去"那一手）**：`internal/session/session.go:41 idBytes=16`、
   `:47 idPrefix="sess_"`、`:59/:69` 的 `Valid()`＝`sess_` + 32 位小写 hex；
   钉在 `internal/session/grants_test.go:108 TestTicket224MintIsRandomAndValid`、
   `:152 TestTicket224TestLiteralsAreOutsideTheMintedShape`、`:350 TestTicket224CoveringDoesNotCrossToANewSession`，
   对面是 `cmd/wisp/ticket224_assembly_test.go:411-413`（两枚字面量常量）＋`:432-439/:491-493`（双向锁）。
   票面 `:58`/N#2 那条"全仓没有仪器能区分重启的两种含义"的判定上限也挂在同一族上
   （`ticket224_assembly_test.go:31-49` 逐字写着这个上限）。
   ⇒ 落地腿若改用"可重算的会话键"来让 CI 变绿，会同时顶红这四枚，且按票面 AC#4（`:38`，未勾）
   属于**本票判失败**的形状。
9. **机器扫的那把**：`tools/d22scan` 禁 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策
   （`AGENTS.md` §1.2、`PLAN.md:1288`）。⚠ 本枚**没跑** d22scan（跑它要 build）。
   直接后果：修 `paths.go:247-267 resolvedForm` 或 `paths.go:140` 时若新增 `filepath.Abs/Clean`
   而不在许可的 resolver 之外——这一条**由落地腿自己跑门禁确认**，期望 rc=0 且零命中。

---

## §5 判不动的地方（需要谁裁）

- **J1（`InAllowlist` 里先失配的是哪条腿）**：本枚只判到"两根形不一致 ⇒ R2 ⇒ L2"，
  而 `paths.go:140`（词法 roots 比对）与 `:152-155`（`resolvedForm` 二次比对）**都能单独**造成 false。
  只读推得最远的是：`:140` **更可能**先失配（因为同一发日志证明 `Canonicalize(缺失叶子)`＝短名，
  而根目录存在、按 `pathresolver.go:127-135` 必被 handle 折叠成长名）；
  但"根被折叠成长名"这一步是从**同发另一枚用例**（存在的文件）外推的，
  ⛔ 直接对 `…\002` 这枚目录现量没人做过。判它要一枚打印 `Roots()` 的探针（§3.4 尺 D）。
  **需要谁裁：必须真跑（落地腿 224-r3，尺 D）。**
- **J2（这算不算产品缺陷／要不要现在修）**：owner 的真机 `%TEMP%` 是长名
  （本机日志逐字 `C:\Users\swq\AppData\Local\Temp\…`），所以今天这条分裂**只在短名拼法下成立**。
  但短名拼法**并不是 CI 独有的怪形**：任何用户手写在 `config.toml` 的 `[fs] allowed_dirs`
  里出现 `C:\Progra~1\…` 之类 8.3 拼法，就会让**该根下所有新建文件**的写都落到 L2
  （fail-closed 方向：多问一张卡，不是放行）。这算"可接受的严格"还是"授权目录形同没设"，
  `PLAN.md`／SPEC-06 §4 没写"授权目录允许 8.3 拼法吗"这一条 ⇒ **未定义即停**。
  **需要谁裁：编排者上报 owner**（碰 D45-3／SPEC-06 §4 的射程，agent 不得自行假设）。
- **J3（票 224 的 AC 有没有勾早）—— 本枚能判的那半已判完，剩下一格归验收腿**：
  本枚**一枚勾选框都没碰**（边界 4）。现读的票面状态
  （`.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md`，取数 10-02 09:2x）：
  `AC#1 :35`／`AC#2 :36`／`AC#3 :37`／`AC#4 :38` **四格未勾**，只有 `AC#5 :39` 已勾；
  N 侧 `:69` N#1、`:70` N#2、`:71` N#3、`:72` N#4、`:74` N#6 已勾，`:73` N#5 未勾。
  框尺实测（用 `^[[:space:]]*- \[ \]` 那把，⚠ 不用 `^- [ ]`——本仓已记过它会整枚漏掉缩进项）＝
  **未勾 5／已勾 6**；而票面 `:78` 那句自称「本票现为 **6 勾／7 未勾**」⇒ 两处差 2 枚未勾，
  最可能是 `N#7 :81`／`N#8 :82`／`N#9 :83` 是**没有框的项目符号**、被那句按事项数进去了。
  ⛔ 本枚不判哪一边对、不改票面，只把这一处算账差登记给编排者。
  ⇒ **题面 (e) 那一支担心的"某格勾早了"，在"读格"这一格上不存在**：AC#2（命中授权就不再弹卡）本来就没勾。
  可判死的只到这一句：**"CI 半边从没验过"成立，且比题面写的更强**——不是"这一枚用例没在 CI 跑过"，
  而是"这一枚用例要验的那条 L1 分支，在 runner 的 temp 形状下于 `cmd/wisp` 装配机上**一发都到不了**"
  （§2(e) 的三个 witness）。⇒ 该格今天**只有本机半边**是既有仪器支持的；
  另一半要么在 §3 的载具上判红，要么具名承认"CI 上这一支不可测"并登记。
  ⚠ 连带读数（只登记不裁）：已勾的 **N#2**（`:70`，凭据＝`TestTicket224ProductionSessionDoesNotSurviveRestart`）
  在 CI 上确实绿，可它那三件控制里的 ③ 是按**弹卡计数**判的
  （`cmd/wisp/ticket224_assembly_test.go:510-512` 计数、`:527-531` 断 `cardsBoot < 1` 才红），
  而 L2 审批卡也算一张卡 ⇒ **CI 那一枚绿分不出"重启后又问了"与"重启后被判成了 L2"**；
  N#1 里"真列 `grant_id` 至少一枚用例"那半（`:69`）落在
  `internal/tools/ticket224_grantid_column_test.go`，而 `internal/tools` 整包在这两发红名册里零命中
  ⇒ 那一枚 CI 绿，但不经过 `cmd/wisp` 那台装配机。
  "算不算已验／这几格要不要退回"归**非实现者验收腿**（D22 双角色），不归本枚。
  **需要谁裁：验收腿，且必须先有 §3 尺 B/尺 C 的本机判红载具，不能只凭这份日志翻框。**
- **J4（同族其余 CI 红的归因）**：新发红名册 **26 行**（`fails_new.tsv`：`--scope=windows` 13＋`wisp-cli-tests.sh` 13），
  基线 **25 行**（`fails_base.tsv`：`--scope=core` 4＋`--scope=windows` 12＋`wisp-cli-tests.sh` 9）；
  按测试名去重取差＝只在基线 4 枚（全是面板/前端契约族）、只在新发 5 枚、两发都红 21 枚
  （现量尺与名册见 `.scratch/wisp/probes/224/c2/01-ci-raw-readings.txt` 末段）。
  ⚠ 我在这两处的第一版数写成了"27 行／24 枚同族"，是没数就写——`awk NR`/`comm` 现量后改成上表（§6-U12）。
  本枚**只逐字读了**三枚的报错文本——`run_test.go:378`、`run_mode101_test.go:506`、
  `pathresolver_junction_windows_test.go:135`。其余（`TestAListWinsWhereBothTablesHit`、
  `TestCanonicalInputGainsNoSecondForm`、`TestClassifyAnchorSpellingIsNotVerdict`、
  `TestBListDefaultDenyAndOverride`（logs_new_all.txt:4073 逐字给了 `{Allow:false NeedL2:true Class:B …}`）、
  `TestPathResolverUNCAListDenied`/`…ExtendedLengthPrefix…`、五枚 `TestSync*`、
  `TestResolvePerCallBudget`、`TestAC246ResidentPipelineAsksThroughTheOneGate`、
  `TestPanelHostRealWindowHopAndLifecycle`、`TestRunPacketCarriesTheLoadedInstructionFiles`、
  `TestTicket101{ManualSwitch,UntouchedConfig}…`、`TestTicket223{HandEditedFsLoosening,PermissionDenied}…`）
  ⛔ **只登记名册，不判同因**——"同族"不等于"同因"。
  **需要谁裁：下一枚差集腿逐枚读报错，或落地腿真跑。**
- **J5（PATH 里带不带 `build/`）**：§3.0 建议去掉 `$PWD/build`，理由是 `build/` 里有 8 枚过期 exe
  （含 `wisp.exe`）；但"本枚用例或其邻居是否真按名字从 PATH exec 过 wisp.exe"
  要跑才知道（`grep -rn 'exec.Command("wisp' cmd internal` 不在本枚射程内跑完）。
  **需要谁裁：必须真跑（落地腿）。** ⛔ 本枚没跑，只给了理由。
- **J6（`0xc0000135` 尺的确切形）**：题面尺的 `$PWD/…` 段复认可用
  （`scripts/wisp-cli-tests.sh:101-109` 认的就是 shell 形，盘符形才是陷阱），
  但那条注释同时警告**带空格的盘符形**会被 MSYS 重写坏——本仓路径含空格
  （`projects plans`），⛔ 本枚**没有**实测"带空格的 `/d/…` 形是否 100% 稳"（要跑）。
  **需要谁裁：必须真跑。** 落地腿第一次跑尺 A 若拿到 `=== RUN` 0 条，先按 §3.0 判"没跑"而不是判绿。
- **J7（`245`/`246`/`33` 那几批本机绿是不是同一枚 HEAD）**：§0.4 的四枚绿读数来自
  2026-09-30 与 2026-10-01 的不同工作树（票 245/246/33 各自的腿），**不是同一枚 commit 的复跑**；
  严格说它们是"历史上绿过四次"，不是"当前 HEAD 绿"。当前 HEAD 的本机读数只有落地腿能给。
  **需要谁裁：必须真跑（尺 A）。**

---

## §6 我推翻／改到编排者题面之处（具名逐条）

- **U1（行号差一枚）**：题面写 `:315`／`:318`／`:322-323`，Go 报出的是 `:316`／`:319`／`:324`
  （`:315`/`:318`/`:322-323` 是 `if` 条件行）。内容逐字全对。
  ⇒ 引这一族时统一按"报错行"数，别按条件区间。
- **U2（":310 那句没响"＝错）**：那句**响了**，在 `cmd/wisp/ticket224_assembly_test.go:311`，
  逐字含「命中授权就不再弹卡」，出处 `.scratch/wisp/probes/orchestrator/ci-delta-1/logs_new_all.txt:3410`。
  而且它是七行里信息量最大的一句（它把 `ledger id`／`pattern`／`row` 三个键一起带出来了，
  本枚的短名结论就从这里的 `pattern` 形状抽出来的）。
- **U3（"三行"→实际七行）**：漏掉的 `:269`（**对照那发自己也 errored**）／`:293`／`:327`（文件从未写出）
  三句不是噪音，它们把方向从"授权没命中"改到"这一发走的不是 L1 路由"。
  ⛔ 只按三行派单会引导落地腿去查授权链，而授权链是好的（§1、§2(e)）。
- **U4（41.72s 的归因）**：题面猜"像是一直等到某个 deadline 才拿到 `timeout`"——方向对、位置错。
  它等的不是 **L1 确认窗口**（那一枚被 `gate.go:137-145` 硬夹在 2–3 秒），
  而是 **L2 审批队列**的 `confirm_timeout_sec`（载具自己在 `approval_reply_201_test.go:139-142` 写成 20）。
  题面提醒的"本仓那族 `300.0x s`＝C18 审批超时常量"同一发就有现场（`run_test.go:378`＝301.70s），
  所以**不需要具名写"对不上"**：它对一个现成常量，只不过对的是 L2 的那个。
- **U5（"新带进 CI 的文件"复认为真，但这格比题面大）**：文件确为区间内新增（§0.2）。
  可"CI 半边从没验过"的实际范围不是这一枚用例，而是**整条 L1 分支在 `cmd/wisp` 装配机上从没被 CI 执行过**
  （§2(e) 三 witness，其中 `run_test.go:378` 与 `run_mode101_test.go:506` 在**基线发就已经红**）。
  ⇒ 这不是回归、也不是票 224 引入的缺口；票 224 只是第一个把"L1 分支必须走到"写成断言的票。
  这一句对派单有实际后果：**主链上"离能用还剩几枚"的账里，这一格不该记在票 224 名下。**
- **U6（候选 (a) 整支可以现在划掉）**：题面请落地腿去查"会话身份取 `COMPUTERNAME`/用户名/profile/`%APPDATA%`"。
  `Mint()` 是纯 `crypto/rand`（`session.go:94-100`），`internal/session` 整包零 env 读取，
  且 CI 那发的行被铸造 id 读得回（`:334` 断言未响）⇒ **不必再派这一支**。
- **U7（候选 (b) 要拆）**：题面把 (b) 写成"DB／grant 行的落盘根在 CI 上是另一棵树"。
  落盘根那一半被同一处读数否掉（行就在 `h.dir`，读得回）；
  真正成立的是**路径拼法**那一半（8.3 短名 vs handle 长名），它属于票 18/72/105/107 那条腿，
  不属于会话授权的账。
- **U8（§4 题面那条禁令写宽了）**：题面说"⛔ 不许把 `SessionID`／grant 塞进 `tools.Decision` 或
  `agent.ToolRequest`，那枚反射钉的射程盖得到"。逐字读 `internal/tools/ticket90_test.go:444-464`：
  `Decision` 侧的禁用子串是 `allow`/`approve`/`grant` ＋类型为 `Answer`；`ToolRequest` 侧是
  `mode`/`allow`/`approve`/`auto`（**不含 `grant`/`session`**）。
  ⇒ `SessionID` 这个名字**两枚扫描都盖不到**；盖得到的是 `GrantID`/`Granted` 这类。
  别指望这枚钉替你拦"把会话 id 挂上去"那一手——那一手要靠 `grant.go:26-29` 那段证词与验收腿。
- **U9（PATH 尺的一处改法）**：见 §3.0——`$PWD/third_party/sherpa-onnx` 段复认可用（形是对的），
  `$PWD/build` 段建议去掉（`build/` 里 8 枚过期 exe 含 `wisp.exe`），理由给足、实测归落地腿（J5）。
- **U10（口径补一条）**：题面提醒"`t.Logf` 也带 `file:line:` 前缀别当成红"——同一发日志里就有标本：
  `pathresolver_expansion_test.go:95/:98/:101` 三行挂在 `--- PASS: TestC26ExpansionMustNotRewriteOntoAnotherTree`
  之下（logs_new_all.txt:4055-4058）。补一句更阴的：`internal/risk/pathresolver_junction_windows_test.go:110-112`
  那两枚 `t.Errorf` 是**"打印 USERPROFILE/A-tier anchor 的加性诊断"**（`:107-109` 自称 "Strictly additive"），
  它们会在**红之前**先响一行——按 `file:line:` 数红名会把诊断行错算成断言。
- **U11（本枚自破的尺，主动登记）**：我第一把找短名覆盖的尺含 `8.3` 一枚分支，它命中的是
  `SPEC-06 §8.3` 这种条款号（`internal/tools/grant_test.go:108/:291`），把 15 枚文件算进了分母；
  换成 `GetShortPathName|短名|short name` 后是 10 枚。⇒ **中文/数字混排的短引在本仓会冒充命中**，
  与题面"多分支 grep 混中文时零命中先怀疑尺"是同一课的**反方向**（这次是"命中过多"）。
- **U12（本枚自破的第二把尺，主动登记）**：J4 与 §0.3 的第一版把红名册写成「27 行」、共红数写成凭印象的数，
  ⛔ 那是 `cat` 过一遍没上尺。现量（`awk 'END{print NR}'` ＋ `cut -f3 | sort -u` ＋ `comm -12/-13/-23`）＝
  **新发 26 行／基线 25 行；按测试名去重：两发共红 21 枚、只在新发 5 枚、只在基线 4 枚**，
  正文两处已按现量改写，尺与名册附在 `.scratch/wisp/probes/224/c2/01-ci-raw-readings.txt` 末段可直接复跑。
  ⇒ 「引『某区间几枚红』必须连口径一起写」这一条本枚自己也踩了一次。
