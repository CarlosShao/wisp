# 101 — 票 90 的持久化存储 `internal/perm` **在生产里零 importer** ⇒ "档位重启后还在"（M3）今天其实还没通

**Status:** open（2026-09-21 17:1x 编排者建；来源=`agent-ticket90b` 交件时**自己点名**的台账缺口）
→ **ready-for-review**（2026-09-21 `agent-ticket101` 交件：AC#1–AC#5 五框都有真实读数，见下；
裁决表仍归验收方，`docs/evidence/s1/101-*.md` 本代理未写）
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

- [x] **AC#1** `grep -rn "wisp/internal/perm" --include=*.go cmd/ internal/ | grep -v _test.go` **非零命中**
      且落在**启动装配路径**上（贴 file:line）。这条就是本票存在的理由：**"有人调用它"必须是量出来的，不是宣布的**。
      读数（`7286297`，本机实测 1 命中，非零）：`cmd/wisp/run.go:53`（import）⇒ 装配路径上
      `run.go:216` `config.NewManager` → `run.go:314` `perm.New(perm.Options{Manager, Confirm, Logf})`
      → `run.go:332` 启动审计 `MODE-READ` → `run.go:343` `Modes: rt.modes`（进决策链的那一行）。
- [x] **AC#2** 端到端两条用例，**分开、不许合并**（票 90 的 AC#3b 边界，我原样搬过来）：
      读数：三条各自一个测试函数，`cmd/wisp/run_mode101_test.go`，"重启"= 同一 data dir 上第二次
      `runTextTask`（新进程面 / 新 Manager / 新 bridge，同一份盘上字节）。
      (a) 手动改成"全自动" ⇒ 重启后读回**仍是全自动**（并仍触发 M4 的那一次 L2 强确认口径）；
          `TestTicket101ManualSwitchSurvivesRestart`：切 `ask_high_risk` 问 0 次、切 `auto_approve`
          问 1 次（R20/M4 非对称），重启后 `MODE-SILENCE mode=auto_approve assessed=L1 effective=L0`
          且 L1 写开 0 张卡；对照半：同一 boot 里工作区外的 `fs.read`（R2 红线）**仍然**被问、被拒。
      (b) **从未手动改过** ⇒ 重启后读回**默认档"每步都问"**；
          `TestTicket101UntouchedConfigRestartsAtDefault`：文件里根本没有 `permission_mode` 键，
          冷启动 3 次都在 `ask_every_step`，每次 L1 写开 1 张卡（与 (a) 的 0 张正好相反），
          且 3 次重启后键**依然不存在**（读默认 ≠ 选默认）。
      (c) 顺带回归票 90 的那条：**会话授权不能跨重启**，别在接线时把它带成能跨。
          `TestTicket101SessionGrantDoesNotCrossRestart`：在 `auto_approve` 档下插一条**未过期**的
          `GrantScopeSession` 行（B 档 `.env`），同会话内已被拒、重启后新会话仍被拒（`Options.Confirmations`
          在装配根保持 nil）；对照半：同一 boot 的普通工作区内写**照自动档执行**，所以"被拒"不是"全都问"。
- [x] **AC#3** **响亮失败面**：档位存储损坏/版本不认识/权限读不到 ⇒ **必须回到最严档并写审计**，
      不许"读不到就按上一次缓存的宽松值"。给三态各自一条用例与真实读数。
      读数：`TestTicket101UnreadableModeFailsLoudlyAndStrict/{存储损坏,版本不认识,权限读不到}` 三条子用例
      各自 PASS —— 盘上先放着 `auto_approve`，再把存储弄坏：退出码 2、`[audit] perm: MODE-READ-FAILED
      ... mode=ask_every_step`、用户可见"配置未就绪"，且 `onRuntime` 从未被叫到（决策链根本没装配），
      日志里不出现 `mode=auto_approve`。非空对照：同一 fixture 不弄坏时 exit 0、装配完成、档=auto_approve。
      ⚠ 口径说明（不扩写、不遮掩）：本票的"最严档"落法是 **不装配 + 退出 2 + 审计点名 ask_every_step**，
      比"继续跑但按最严档"更严；票 90 的 `validateRisk`（未知值让 load 失败，票 83 规则）原样保留。
- [x] **AC#4** 变异：把装配那行**注释掉** ⇒ AC#2 的 (a) 必须红
      （证明这条线不是"恰好也绿"）。锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，
      **编译失败不算变异**；**变异只在 `/tmp` 仓外快照里做**（目录带会话后缀）。
      读数：快照 `/tmp/wisp-t101-mut-agentticket101`（`git archive HEAD` + 拷 `third_party/sherpa-onnx`），
      锚点 `run.go:343` `Modes: rt.modes,` 被注释；同链 grep：主树 `run.go:343: Modes: rt.modes,` /
      快照 `run.go:343: // MUTATION(t101 AC#4): Modes: rt.modes,`。快照**编译通过**（`ok`/`FAIL` 有测试读数，
      非 build 失败），`go test -run TestTicket101 ./cmd/wisp/` rc=1：
      **(a) FAIL** —— `run_mode101_test.go:305: auto_approve still opened 1 card(s) for an L1 write;
      the档 was read but never reached the decision chain`；(b)(c) 与另两条仍 PASS（正是"默认档=fail-closed"
      的形状：只有 (a) 能区分这条线）。快照已删除；主树 `git diff --quiet` rc=0、
      `git status --porcelain cmd/wisp/ internal/perm/` 空。
- [x] **AC#5** 门禁：`gofmt -l`/`gofumpt -l` 空、`go vet ./cmd/wisp/ ./internal/perm/` rc=0、
      `go test -count=2 ./internal/perm/` rc=0 且逐条点名 SKIP/FAIL。
      读数：`gofmt -l cmd/wisp internal/perm` 空（rc=0）、`gofumpt -l` 同（rc=0）、
      `go vet ./cmd/wisp/ ./internal/perm/` rc=0、`go test -count=2 ./internal/perm/` rc=0
      = 28 PASS / **0 SKIP / 0 FAIL**（逐条点名：无 SKIP、无 FAIL，15 个顶层测试 ×2 减去缓存复用）。
      ⚠ `go test ./cmd/wisp/` 在本机是**加载期 `0xc0000135`（缺 sherpa dll，票 98 的账）** ⇒
      **你这条不能拿它当判据**，判据要么走票 98 的注入命令，要么显式登记"宿主包本机不可测"，
      **不许默默跳过**（这正是票 98 要收的那个洞）。**收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。
      —— 本票**走了票 98 的注入命令**（`PATH="$PWD/third_party/sherpa-onnx:$PATH"`）：
      不带注入 `go test ./cmd/wisp/` = `exit status 0xc0000135` FAIL（复现票 98），
      带注入 = 全包 `ok 40.903s` rc=0；`-count=2 -run TestTicket101` = 16 PASS / 0 SKIP / 0 FAIL rc=0。
      `sh scripts/d22scan.sh`：首跑 rc=1（我在 `run_mode101_test.go:41` 注释里放了 `⚠`，ban #8 抓到）
      ⇒ 已改掉，复跑 **rc=0 clean**（live scope：bans #1-5 internal/=197、cmd/=20，ban #6 frontend/=40，
      ban #7 internal/tools/=17，ban #8 design/=16、frontend/=40、internal/=337、cmd/=26）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④）；票面 append-only（改行前先读；标题前插段落要重抄标题，删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字；15 次工具调用内交回第一枚 checkpoint；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/config/`+`internal/agent/`+`internal/perm/`（票 90 的两个会话刚交件，**语义别改**）、
`internal/winsec/`+`internal/secret/`（`agent-ticket89b` 正在写退回单）、`frontend/`+`internal/panel/`（票 92 地界）。
**`cmd/wisp/run.go` 可能同时是票 77/92 的落点** ⇒ 动它之前先在票面登记"需要谁协调"，别抢。

## Progress log（append-only）

- 2026-09-21（**agent-ticket101 交件**）：AC#1–AC#5 五框都有读数（见上，逐条附真实数字），
  Status 建议 `ready-for-review`（裁决表按票面留给验收方，`docs/evidence/s1/101-*.md` 我没写）。
  落了什么：`cmd/wisp/run.go` 装配根改用 `*config.Manager`（唯一真源、可原子写回）→ `perm.New` →
  `tools.Options.Modes`；M4 的那一次确认接真 `approval.Gate`（D47 自注册 `host:mode-switch`）；
  启动写 `MODE-READ origin=startup`，读不到写 `MODE-READ-FAILED ... mode=ask_every_step` + 不装配链 + 退出 2；
  `Options.Confirmations` **保持 nil**（从盘上读的只有 mode）。
  **语义零改动**：`internal/risk/mode.go`、`internal/config/permmode.go`、`internal/perm/store.go`、
  `internal/tools/mode.go` 一行未碰（`git show --stat 7286297` 只有 `cmd/wisp/run.go` + 一个新测试文件 + 本票面）。
  三条**交回时点名**的账（不是我藏起来的东西，是要别人接着做的事）：
  ① **票 90 有一处"注释说做了、代码没做"**：`perm.New` 的文档写"records what the mode the process started
  with, so the audit trail ... has a first line to compare against"，但它调的是 `record()`（只进内存历史），
  不是 `audit()`（才写 sink）⇒ 启动那一行**从来没进过审计**。我在装配根用 `MODE-READ` 自己补了这行
  （AC#2(a)(b) 与 AC#3 都断言它），**没有去改 perm**（语义不是我的地界）。
  → 归票 90 的验收判：要么 `New` 改走 `audit()`，要么把那句注释改掉。
  ② **`Store.Set` 至今没有生产调用者**：面板 composer 是票 92 的（M5 输入口 + 只许"显示 + 发起请求"），
  原生点击通道在票 77。⇒ 今天真机上"手动改档"只有两条路：改 `config.toml` 的 `risk.permission_mode`
  （票 90 已定：启动时文件即操作者声明），或等 92 接上 `rt.modes.Set`。**我没有假装它有调用者**——
  这正是本票立案的 memory 第 8 条同一个形状，只是这次读侧已经接上。
  ③ **控制台跑法切不到全自动**：`wisp run` 没有原生点击通道 ⇒ `confirmModeSwitch` 必然被拒
  （`TestTicket101ModeSwitchUsesTheRealL2Gate` 钉住：卡片确实开 1 张、档不变、审计写 `result=refused-confirm`、
  盘上不落 `auto_approve`）。这是 fail-closed 的正确形状，但**owner 那边要它"可用"得等 77/92 的原生侧**，
  别把这条读成"功能坏了"。
  next= 非实现者按上表 1:1 出裁决表；票 92 落 `Set` 的调用者时**复用** `assembleRuntime` 已有的
  `rt.modes` 与 `runSpec.modeConfirm` seam，别再开第二条装配路径；票 77 若动 `run.go` 同一段，以 `7286297` 为准。
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
