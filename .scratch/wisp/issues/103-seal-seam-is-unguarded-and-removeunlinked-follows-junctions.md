# 103 — 密封的**两处旁路**：`SetPathResolver` 谁都能装（装了个橡皮图章就静默重写外来 DACL），`RemoveUnlinked` 能沿 junction 删别人真文件且返回 nil（票 94 验收的 R-c / R-b）

**Status:** open（2026-09-21 17:4x 编排者建；来源=`acceptor-ticket94` 的两条残留，它自己给的方向是"立案，不救不判"）
**Type:** 安全边界（一个是**接缝无守卫**，一个是**今天够不到的陷阱**——两者不同档，同票不同判据）
**Blocks:** 票 94 挂 `-done` 的条件之一 · **Blocked by:** nothing
**Packages:** `internal/winsec/`（seam 的注册口与 `RemoveUnlinked`）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、
              `internal/risk/pathresolver*.go`（**入口展开那条归票 102，别在这儿顺手改**）、`rules_gateway.go`、
              `tools/d22scan/**`、`allowlist.txt`。

## 两条读数的原文位置

`docs/evidence/s1/94-adversarial-acceptance.md` 的 **PROBE A**（`SetPathResolver` 无守卫）与 **PROBE F / PROBE W**
（`RemoveUnlinked` 沿 junction 删目标 + `removeStray` 的 `WalkDir` **不下降** ⇒ 今天不可达）。

- **R-c（今天可利用）**：`SetPathResolver` 没有任何守卫 ⇒ 包外可以装一个"什么都说 OK"的解析器，
  之后 winsec **静默重写外来主体的 DACL** 并报告成功。
  ⚠ 这是票 94 为了绕开 `winsec → risk` **传递依赖成环**（`risk → observe → secret → winsec`，`go list -deps` 实测）
  才引入的装配形状 ⇒ **形状本身是必要的，缺的是守卫**。别把 seam 拆掉倒回 `filepath.Abs`（那是把票 94 白做）。
- **R-b（今天够不到，但是陷阱）**：`RemoveUnlinked` 收到一个 junction 时可以删掉**别人树里的真文件**并返回 nil。
  验收代理自己实测了"为什么今天够不到"（`removeStray` 用 `WalkDir` 且**不下降进链接**，PROBE W）
  ⇒ 所以本票的判据是**把它变成"一旦有人接上就会红"的 tripwire**，而不是"现在就有人踩了"。
  ⚠ 修它要动 `internal/memory`（票 18/79 地界）⇒ **AC#2 明写不许顺手扩界**。

## AC（1:1，裁决表 `docs/evidence/s1/103-*.md` 由验收方出）

- [ ] **AC#1（R-c）** seam 只能被**装一次**、且装的必须是**真解析器**：给注册口加守卫
      （幂等/一次性 + 类型上不给伪造留门，或伪造时**响亮失败并审计**）。
      判据用例两条腿：**装第二个 ⇒ 红**；**装一个恒说 OK 的 ⇒ winsec 必须拒**，
      且"什么都没发生"不算绿（要能证明**外来 DACL 没被改**，取 SID 级读数）。
- [ ] **AC#2（R-b）** 给 `RemoveUnlinked` 一条**平台用例**：输入一个 junction/符号链接 ⇒ **必须拒、返回错误且不删目标**；
      并在 `internal/memory/removeStray` 一侧补一条**"下降进链接就会红"的 tripwire**（断言遍历不下降）。
      ⚠ 修 `RemoveUnlinked` 的**语义**若需要动 `internal/memory` ⇒ **停手登记交回**，本票只做守卫与 tripwire。
- [ ] **AC#3** 变异三向：① 去掉注册守卫 ⇒ AC#1 红；② 把"伪造解析器"改成"注册了但被拒" ⇒ **必须仍有一条用例绿**
      （证明它红在守卫而不是红在噪声）；③ 把 tripwire 的"不下降"改成"下降" ⇒ 该用例红。
      锚点=承载行为那一行，**同链 grep 证落地**，`go build` rc=0 先量到（**编译失败不算变异**），
      还原在 `/tmp` 快照里做并证 `diff -q` 干净。
- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/` rc=0 且逐条点名 SKIP/FAIL
      （⚠ 票 89 的 `TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs` 刚落地，别把它改松）；
      `gofmt`/`gofumpt` 空、`go vet` 与 `GOOS=linux go vet` **按包** rc=0、
      **收尾前必跑 `sh scripts/d22scan.sh`** 纯净树 rc=0（A64②）。
      ⚠ `go test ./cmd/wisp/` 本机需按**票 98** 的注入命令跑；`TestResolvePerCallBudget` 负载下会假红（票 86）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；
**不 push**；不在仓内建 worktree（A38④，快照带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字；**真机测 junction 只在临时目录造、测完清掉，绝不删真数据**；
15 次工具调用内交回第一枚 checkpoint；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/winsec/` 刚被 `agent-ticket89b` 改过、它的**复验正在跑** ⇒ 你与它可能同文件；
开工前先 `git log --oneline -5 -- internal/winsec/` 看有没有未结的验收，并在票面登记"我等谁让路"。

## Progress log（append-only）

- 2026-09-21 17:4x（编排者）：建票。**两条我合成一张但判据分开**，理由：它们是同一个包、同一个"密封动作"的两处旁路，
  分两张会互相等（都要碰 `internal/winsec/`）；但**档位不同**我写在标题里了——
  **R-c 今天可利用**、**R-b 今天是陷阱**（验收代理自己实测了"够不到"）⇒
  AC#1 是修，AC#2 只是守卫 + tripwire，**不许把两者混成一个"都已修"的读数**。
  另一个我明写的克制点：**入口展开那条（R-a）不在这张票**，它归**票 102**，因为它在冻结的 C26 实现里、
  影响面比 winsec 宽得多——放一张票会让"守卫"级别的活被"契约级"的活拖住（票 89/95 的分票理由同形）。
  next= 排 in **票 102 之后**（fail-open 优先），但可与它并行——只要不撞同一文件。
