# 224-c2 只读普查：票 224 读格「本机绿 / CI 红」的边界在哪

代号 `224-c2`。只读腿：零 `go test`／零 `go build`／零 `go vet`／不执行任何 exe。
本文件所有读数都来自**已落盘的 CI 日志**与**源码**，凡"必须真跑才答得出"的格一律具名写成
「待落地腿自量」并给确切尺（§3）。

## 本轮（骨架发）状态

本发只交 §0 起手锚 + §5 判不动格的前 4 条 + 一条已经把整格判死的主读数（见 §0.4）。
§1／§2／§3／§4 的正文在随后的增量发里补，本节以下位置在补之前是空的——空的和"待填"不是一回事：
下面每一格都有出处或明确的判不动理由。

---

## §0 起手锚

### 0.1 时刻与 HEAD（取数时刻 2026-10-02 08:54 +08，共享树里的枚数会漂）

- 分支 `dev`，HEAD `0c1ec3a22a734dcd7e6acbeeea34a2cdff9b68f9`
  （`ledger(A512)`，2026-10-02 08:51:42 +0800）。
- 编排者给的推送区间 `0589fd9c..8ae4c23e`：`git merge-base --is-ancestor` 复认两枚都是 HEAD 的祖先，
  区间实测 **231 枚** commit（`git rev-list --count 0589fd9c..8ae4c23e`，取数时刻同上）。
- CI 读出来的那一发 run 的时间戳（日志内嵌）：新 run `2026-10-01T16:17:05Z`，
  基线 run `2026-09-30T23:20:14Z`。

### 0.2 `git status --porcelain cmd internal`（取数时刻 2026-10-02 08:54 +08）

```
 M cmd/wisp/panel_inbound.go
 M cmd/wisp/run.go
 M internal/panel/bridge.go
 M internal/panel/composer.go
 M internal/panel/composer_dispatch.go
 M internal/panel/composer_dispatch_test.go
 M internal/panel/l2_grant_boundary_test.go
 M internal/panel/pump.go
?? cmd/wisp/panel_config_store.go
?? internal/ball/sta_release_windows_test.go
?? internal/config/settings.go
?? internal/config/settings_248_test.go
?? internal/panel/config_handlers.go
```

非空＝别人在飞的腿（`248-r1`／`33-r8`／`250-r1`），本枚一律不动、不提交、不 `stash`。
⚠ 由这条引出的口径：本枚引用的 `cmd/wisp/run.go` 行号是**工作树读数**（该枚文件正在被改），
`git diff --numstat HEAD -- cmd/wisp/run.go` ＝ `19 2`；其余被引文件
（`cmd/wisp/ticket224_assembly_test.go`／`internal/tools/bridge.go`／`internal/tools/grant.go`／
`internal/tools/paths.go`／`internal/session/grants.go`／`internal/risk/rules_gateway.go`／
`internal/agent/approval/gate.go`／`cmd/wisp/approval_reply_201_test.go`）都不在这份名单里，
行号＝HEAD＝CI 那一发。

`cmd/wisp/ticket224_assembly_test.go` 在工作树里与 8ae4c23e **逐字节等长**（574 行，
`git show 8ae4c23e:... | wc -l` = 574，`git diff --stat HEAD -- <该文件>` 空）⇒ 题面那三行行号
读的是同一份文件，无漂移。

### 0.3 复认题面那三行断言（逐字，源码侧）

`cmd/wisp/ticket224_assembly_test.go`：

- `:315-316` `if measured.decision != agent.DecisionAllowGrant { t.Errorf("covered call booked decision=%q, want %q", …) }`
  —— 报错落在 `:316`（t.Errorf 那一行），不是 `:315`。
- `:318-320` grant_id 断言 → 报错行 `:319`。
- `:322-324` GRANT-USE 审计行断言 → 报错行 `:324`（条件是 `:322-323`）。

⇒ 题面三行的**内容**逐字对上，行号差一枚（315 vs 316）。见 §6-U1。

### 0.4 本发主读数：CI 那一发的完整红名册是七行，不是三行

出处 `.scratch/wisp/probes/orchestrator/ci-delta-1/logs_new_all.txt`（test-windows job 原始日志），
`--- FAIL` 在 `:3418`，中间这枚用例的**全部** file:line 行是：

| 日志行 | 报出的源码行 | 源码里的形状 | 题面提到过吗 |
|---|---|---|---|
| `:3408` | `cmd/wisp/ticket224_assembly_test.go:269` | `t.Errorf("the control call errored: the L1 window running out means execute")` | ⛔ 没提 |
| `:3409` | `…:293` | `t.Errorf("the granted call errored: %v", measErr)` → 打出 `true` | ⛔ 没提 |
| `:3410` | `…:311` | `t.Errorf("a live session grant left %d cards … 「命中授权就不再弹卡」 is not happening …")` | ⚠ 题面说这句"没响" |
| `:3411` | `…:316` | decision 断言 | ✔ |
| `:3412` | `…:319` | grant_id 断言 | ✔ |
| `:3413` | `…:324` | GRANT-USE 断言 | ✔ |
| `:3416` | `…:327` | `os.Stat(granted)` 断言：文件没写出来 | ⛔ 没提 |

（`:3414`/`:3415` 是 `:324` 那句 `%s` 里带的审计行内容，不是独立断言。）

⇒ **推翻题面两处**：
1. `:311` 那句（"命中授权就不再弹卡"）**响了**，而且是七行里最先响的那一句之一（见 §6-U2）；
2. 红名不止三行——`269`（**对照那发自己也 errored**）／`293`／`327`（覆盖那发根本没写出文件）三句题面没有。
   这三句不是噪音，它们把方向从"授权没命中"改到"这一发走的是哪条路由"。

### 0.5 判死读数（本文件的结论，出处逐条）

**这一格在 CI 上红的直接原因不是授权那一跳，而是被装配的 run 在 runner 上把一次
「授权目录内的 fs.write」判成了 L2，于是那一发走的是审批队列（300/20 秒）而不是 L1 窗口（2 秒），
而 grant 那一支在源码里只挂在 L1 分支上。**

四条独立读数（前三条是本枚现量，第四条是 CI 日志里同族旧红）：

1. `decision="timeout"` 这个列值在 `internal/tools/bridge.go` 里**只有一个写入点**：
   `:457-460` 的 `case risk.L2: … case AnswerTimeout: dec.DecisionColumn = agent.DecisionTimeout`。
   L1 分支（`:437-442`）拿到 `AnswerTimeout` 写的是 `DecisionAllow`（`"allow"`），
   逐字注释 "The L1 window running out unopposed MEANS EXECUTE (SPEC-06 §2)"。
   ⇒ CI 那发既然读到 `timeout`，它就在 L2 分支上。
2. 时长对得上 L2 而**对不上** L1：`cmd/wisp/approval_reply_201_test.go:139-141` 把
   `[risk] confirm_timeout_sec = 20` 写进配置（那是 **L2 审批超时**，见
   `internal/config/schema.go:449-450`），而 L1 窗口另有其键 `risk.l1_window_sec`
   默认 `2`（`internal/config/schema.go:451-453`），并且 `internal/agent/approval/gate.go:137-145`
   把它硬夹在 `[MinL1Window=2s, MaxL1Window=3s]`（`internal/agent/approval/queue.go:120-122`）——
   **L1 路径在物理上不可能等 20 秒。**
   CI 实测：开机完成 `16:16:22.498`（`:3406`）→ GRANT-RECORD 的 `expires_at=1793463402`
   ＝ `2026-10-31T16:16:42Z`，`session/grants.go:163` 用的是 `now+30d`（`clockCeilingDefault`
   = `memory.GrantAuditTTL`，`:108`）⇒ 落行时刻＝`2026-10-01T16:16:42Z`＝开机后 **19.5 秒**；
   最后一条 WAL 日志 `16:17:02.589`（`:3417`）＝再 **20.6 秒**。41.72s ≈ 20+20+开机。
   本机同码：`--- PASS … (3.20s)`（`.scratch/wisp/probes/33/r9/09-fullpack-1.txt:1325`，
   取数时刻日志内嵌 2026-10-01T17:55:28→31），`3.17s`（同目录 `11-fullpack-2.txt:1324`），
   `3.23s`（`.scratch/wisp/probes/245/v1/gate-cmdwisp-full.txt:1095`，2026-09-30）。
   ⇒ 本机两发调用一共只等了约 2 秒＝**只可能是一枚 2 秒的 L1 窗口**。
3. 判成 L2 的那一步在授权目录比对上，而 runner 的临时根是 **8.3 短名拼法**
   `C:\Users\RUNNER~1\AppData\Local\Temp\…`（CI 日志里每一枚 temp 路径都是这个拼法，
   本枚用例的 `pattern` 逐字含 `RUNNER~1`，`:3410`/`:3415`）。
   `internal/tools/paths.go:133-163 InAllowlist` 要求**两次**包含：`:140` 对词法 canonical，
   `:152-155` 还要对 `resolvedForm()`（`filepath.EvalSymlinks`，`:247-267`）解析过的形状再包含一次；
   任一次拼法不一致（短名 vs 长名）就返回 false ⇒ `internal/risk/rules_gateway.go:44-51` 出
   `R2: 目标路径在授权目录之外` ⇒ L2。
   同一发 CI 日志里有**现成的正面读数**：`TestPathResolverShortNameAListDenied` 红在
   `internal/risk/pathresolver_junction_windows_test.go:135`
   ——「`got "C:\\Users\\runneradmin\\…"` want `"C:\\Users\\RUNNER~1\\…"`」
   （日志 `:4064`）＝**同一台机器上 `risk.Resolve` 会把短名展开成长名**，而那枚测试的期望值
   （来自 `t.TempDir()`）是短名。两枚形状合起来就是"root 与 target 拼法不同"的现场。
4. 同族**早已存在、与票 224 无关**的 CI 红：`cmd/wisp/run_test.go:378`
   `TestComposedGateBlocksAWriteForTwoSeconds` 在基线 run（`0589fd9c` 那一发）与新 run **两发都红**，
   `--- FAIL … (301.70s)`（新）/ `(301.12s)`（基线），报错逐字：
   `run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`
   （新 run 日志 `:3099`/`:3095`；基线日志 `:3709`）。
   300 秒＝`confirm_timeout_sec` 的 default（`internal/config/schema.go:450`），
   这句话是**只有 L2 审批路由才会印的句子**，出现在一枚"断言 L1 两秒窗口"的用例上＝
   同一根因的第二个现场，而且它比票 224 早一发就红。

⇒ 于是"本机绿/CI 红"这一格**判死**：不是会话授权那一跳缺实现（题面 (e) 支在这一发被推翻），
也不是等满某个现成 deadline 之外的时序玄事（(c) 支有确切常量对上：20s 与 300s 两个都是
`confirm_timeout_sec`，41.72s＝两个 20s）。剩下真正没判死的是**拼法在哪一步分家**——
见 §5。

---

## §1 这台装配机上「命中授权」那一跳到底有几环

（本发未覆盖，随后增量发补；§0.5 已给出这一跳的成败不取决于本节，取决于路由等级。）

---

## §2 本机绿／CI 红的候选与各自那把尺

（本发未覆盖；§0.5 已把 (c) 判死为"对上了现成常量、且不是 224 的实现缺口"，
(a)(b)(d)(e) 的逐支判语在增量发里给，出处已在 §0.5 列全。）

---

## §3 落地腿 224-r3 的派单料

（本发未覆盖，随后增量发补。）

---

## §4 撞钉预检

（本发未覆盖，随后增量发补。）

---

## §5 判不动的地方

- **J1（拼法在 InAllowlist 的哪一步分家）**：§0.5 第 3 条只证到"两次包含里至少有一次在 runner 上
  不接受短名拼法"。到底是 (i) `NewPathCanonicalizer` 把**存在的**根目录（`…\002`）经
  `risk.Resolve` 记成了长名、而**叶子还不存在**的目标文件被记回短名；还是 (ii) 两边都停在短名、
  断在 `:152 resolvedForm()` 那一步——**只读判不动**：两形的现场读数（`Roots()`、
  `Canonicalize(target)`、`EvalSymlinks`）都要真跑才拿得到。
  需要谁裁：必须真跑（落地腿 `224-r3` 现量，尺在 §3）。
- **J2（这算不算产品缺陷）**：owner 的真实机器上 `%TEMP%` 是长名（本机日志逐字
  `C:\Users\swq\AppData\Local\Temp\…`），所以"授权目录内的写在真实用户机上判成 L2"这一形
  今天**只发生在 runner 形状下**。它是"CI 环境畸形"还是"短名拼法的用户会撞上真缺陷"，
  取决于 D7 自用期是否承诺覆盖 8.3 拼法——**规格没写这一条**（未定义即停，AGENTS.md §2）。
  需要谁裁：编排者上报 owner。
- **J3（票 224 的 AC 勾有没有勾早）**：本枚**不碰任何 AC 勾选框**（边界 4）。
  §0.5 只支持"那一跳的实现被这发 CI 读数证明'还没被 CI 验过'，但不支持它没接"。
  判"这一格算不算已验"要的是**在能判红的本机载具上**（§3）重跑，不是这份日志。
  需要谁裁：非实现者验收腿。
- **J4（`301.70s` 那一族里其余几枚的归因）**：`TestTicket101*`／`TestTicket223*`／
  `TestAC246ResidentPipelineAsksThroughTheOneGate` 在新旧两发的红名册里都有（
  `.scratch/wisp/probes/orchestrator/ci-delta-1/fails_base.tsv` / `fails_new.tsv`），
  但本发只逐字读了 `TestComposedGateBlocksAWriteForTwoSeconds` 一句，其余几句的报错文本没读，
  ⛔ 不许把"同族"当"同因"。**本枚只登记，不判**。
  需要谁裁：必须真跑／或下一枚差集腿逐枚读报错。
