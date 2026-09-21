# 101 — 票 90 的持久化存储 `internal/perm` **在生产里零 importer** ⇒ "档位重启后还在"（M3）今天其实还没通

**Status:** open（2026-09-21 17:1x 编排者建；来源=`agent-ticket90b` 交件时**自己点名**的台账缺口）
**Type:** 能力已实现但没接线（memory 第 8 条那个形状的又一例：**测试证明它会工作，生产里没人叫它**）
**Blocks:** owner 要的 M3（档位持久化）**在真机上是否成立** · **Blocked by:** nothing
**Packages:** `cmd/wisp/`（装配根：启动时读档、把 mode 注进链）、`internal/perm/`（存储本体，别改语义）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、
              `internal/risk/pathresolver*.go`、`rules_gateway.go`、`tools/d22scan/**` 与 `allowlist.txt`。

## 现场（90b 交回的原文要点）

票 90 把三档语义、红线守卫、审计、持久化存储都做完了，AC#1–AC#5 + AC#3b 六框有读数
（含它补的那 5 轮双向变异）。但它自己登记了两条**装配缺口**：

1. **`internal/perm` 生产零 importer**——存储包写好了、测试绿，**没有任何生产代码 import 它**；
2. **`cmd/wisp/run.go` 里没有 `Modes:` 注入**——也就是启动时**没人把档位塞进决策链**。

90b 的判词是"代价被 fail-closed 兜住"（读不到档 ⇒ 落回最严的"每步都问"）——**这个判法我认可**，
方向是对的：缺线不会导致"意外宽松"。**但**它同时意味着 **owner 的 M3（手动选过就一直按那档）今天不成立**：
用户改了档、重启，会因为没人读档而回到默认。**这是一个功能没通，不是一个功能有洞。**

## AC（1:1，裁决表 `docs/evidence/s1/101-*.md` 由验收方出）

- [ ] **AC#1** `grep -rn "wisp/internal/perm" --include=*.go cmd/ internal/ | grep -v _test.go` **非零命中**
      且落在**启动装配路径**上（贴 file:line）。这条就是本票存在的理由：**"有人调用它"必须是量出来的，不是宣布的**。
- [ ] **AC#2** 端到端两条用例，**分开、不许合并**（票 90 的 AC#3b 边界，我原样搬过来）：
      (a) 手动改成"全自动" ⇒ 重启后读回**仍是全自动**（并仍触发 M4 的那一次 L2 强确认口径）；
      (b) **从未手动改过** ⇒ 重启后读回**默认档"每步都问"**；
      (c) 顺带回归票 90 的那条：**会话授权不能跨重启**，别在接线时把它带成能跨。
- [ ] **AC#3** **响亮失败面**：档位存储损坏/版本不认识/权限读不到 ⇒ **必须回到最严档并写审计**，
      不许"读不到就按上一次缓存的宽松值"。给三态各自一条用例与真实读数。
- [ ] **AC#4** 变异：把装配那行**注释掉** ⇒ AC#2 的 (a) 必须红
      （证明这条线不是"恰好也绿"）。锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，
      **编译失败不算变异**；**变异只在 `/tmp` 仓外快照里做**（目录带会话后缀）。
- [ ] **AC#5** 门禁：`gofmt -l`/`gofumpt -l` 空、`go vet ./cmd/wisp/ ./internal/perm/` rc=0、
      `go test -count=2 ./internal/perm/` rc=0 且逐条点名 SKIP/FAIL。
      ⚠ `go test ./cmd/wisp/` 在本机是**加载期 `0xc0000135`（缺 sherpa dll，票 98 的账）** ⇒
      **你这条不能拿它当判据**，判据要么走票 98 的注入命令，要么显式登记"宿主包本机不可测"，
      **不许默默跳过**（这正是票 98 要收的那个洞）。**收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④）；票面 append-only（改行前先读；标题前插段落要重抄标题，删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字；15 次工具调用内交回第一枚 checkpoint；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/config/`+`internal/agent/`+`internal/perm/`（票 90 的两个会话刚交件，**语义别改**）、
`internal/winsec/`+`internal/secret/`（`agent-ticket89b` 正在写退回单）、`frontend/`+`internal/panel/`（票 92 地界）。
**`cmd/wisp/run.go` 可能同时是票 77/92 的落点** ⇒ 动它之前先在票面登记"需要谁协调"，别抢。

## Progress log（append-only）

- 2026-09-21（**agent-ticket101 开工登记**，动 `cmd/wisp/run.go` 之前，按本票 Rules 第 49 行）：
  ① 现场：`git status --porcelain cmd/wisp/` **空** ⇒ 此刻 `cmd/wisp/` 没有别人的未提交改动，
  我不覆盖任何东西（工作树别处的 `internal/secret/`+`internal/winsec/`+票 89 票面属 `agent-ticket89b`，本票不碰不提交）。
  ② 我要动 `cmd/wisp/run.go` 的哪几行、为什么：
  - `:189` `config.LoadFile` → `config.NewManager`（同一 load 管线，1 行）：mode 的唯一真源必须是**可写回**的
    Manager，否则 `perm.Store.Set` 没有可持久化的对象；`cfg` 改为 `mgr.Config()` 的快照。
  - `:262` 起的 `tools.Options{…}` 增加 **`Modes: rt.modes`（1 行 = AC#4 的变异锚点）**：把档位注进决策链。
  - `:228`/`:255` 之后新增 ~16 行：`perm.New(perm.Options{Manager, Confirm, Logf})` 与 AC#3 的响亮失败面。
  - `agentRuntime` 加字段 `modes *perm.Store`；`runSpec` 加一个注入位 `modeConfirm`（默认实现 =
    `approval.Gate.PendingApproval` 真实 L2 卡片，测试用来代表"原生侧点了一次允许"）。
  - **不碰** `cmd/wisp/` 其余文件；`internal/perm`/`internal/risk`/`internal/config`/`internal/tools`
    的语义**一行不改**（只在 `cmd/wisp/` 新增测试文件）。
  ③ 需要谁协调：**票 92**（面板 composer 是 `Set` 的下一个生产调用者；本票只装读侧 + 把 M4 的 L2 通道接上
  真 gate，`Set` 本身在本票只有装配根持有、无 UI 入口 ⇒ 我不会假装它有）· **票 77**（宿主）：若随后要动
  `run.go` 同一段，以本票提交后的 `assembleRuntime` 为准，新增注入点请复用 `runSpec` 的 seam。
  ④ 我自己钉死的边界（票面 AC#2(c) 的原因）：本票**不**给 bridge 接任何 grant 来源
  （`tools.Options.Confirmations` 保持 nil）——从盘上读回来的只有 mode，**没有** session grant。
- 2026-09-21 17:1x（编排者）：建票。`agent-ticket90b` 交件时把这两条写在收尾段（**它没藏**，
  还说"由票 92/77 落，代价被 fail-closed 兜住"）。我的判断：**fail-closed 兜住 ≠ 功能通**——
  兜住的是"不会意外宽松"，没兜住的是"owner 要的那条 M3 今天没生效"。
  所以单独立案，**不塞回票 90**：90 的六框判据是"语义/守卫/审计/持久化 API"，它确实做到了；
  缺的是装配根那一行，那是另一个交付面（memory 第 8 条的对策：**"做一个能力"和"把它接上"要么同票、要么当场立案**——
  这次是同票做不到（`cmd/wisp/` 被别人的地界压着），所以立案）。
  next= 等 `agent-ticket89b` 交件（它此刻在 `internal/secret/`+`internal/winsec/`）⇒ 与本票无文件冲突，
  但 `cmd/wisp/` 要留给票 92 的话就先派本票；两票撞车时**本票优先**（它挡的是 owner 已拍板的功能）。
