# 102 — C26 的 `expandInput` 会把含 `%VAR%` 或前导 `~` 的拼写**改写到另一棵树**、还**封它并返回 nil**（票 94 验收代理的 PROBE B2 实测）

**Status:** open（2026-09-21 17:4x 编排者建；来源=`acceptor-ticket94` 的残留 **R-a**，它判"这条比票 94 本票的判据更值钱"）
**Type:** **安全缺陷（fail-open 形状）**——不是整理、不是命名
**Blocks:** 票 94 挂 `-done` 的条件之一（验收代理原话：两句过头文案要降级，其中一句就是"路径已解析"）
· **Blocked by:** nothing（`internal/risk/` 此刻无人写；票 89 的修复在 `winsec`/`secret`，不撞）
**Packages:** `internal/risk/pathresolver*.go`（⚠ **冻结文件**，见下面 D22 段）+ `internal/winsec/` 的调用侧。

## 实测（验收代理在 `/tmp/wisp94acc-94` 快照里做的 PROBE B2，原文在它报告里，file:line 自己去取）

`expandInput` 在处理输入路径时会做**环境变量展开**一类改写。
⇒ 含 `%VAR%` 或前导 `~` 的输入被**改写到另一棵树**上；调用方拿到的是 **nil（成功）**，
而 winsec 那头**照原样把"它以为的那一棵"封了**——即"封 A 的意图 + 判 B 的事实 + 报告成功"三件同时发生。
这**正是** `ban #6`/`PLAN.md:1588` 那一族最怕的形状：**判定所依据的形式与人看到的形式不是同一个东西**。

⚠ **别把它读成"winsec 的 bug"**：`filepath.Abs` 那条已经被票 94 修掉了；这一条在 **C26 自己**的入口规范化里，
影响面比 winsec 宽得多（所有走 `Resolve` 的判定都吃它）。

## 为什么要先判一件事再动手

环境变量/家目录展开本身**不是错**——它甚至可能是契约要求的行为。
真正的问题是 **"改写之后的树"和"被授权/被封存/被展示的树"之间没有任何绑定**：
调用方无从知道自己拿到的结果对应哪一个拼写，也没被告知发生过改写。
⇒ 修法有两个方向，**你必须先判再选，并给理由**：
- **(A) 不许改写**：入口拒绝带 `%VAR%`/`~` 的输入（响亮失败），让调用方自己先展开成绝对形式；
- **(B) 改写要记账**：`Result` 里带上"我做过哪种改写"的显式字段，
  **且凡是要拿它做安全决定（授权、封存、放行）的路径必须消费这个字段**（不消费就是新的洞）。
⚠ **不许选**第三条："把展开留着但只在 winsec 侧拦"——那是同一个洞再挖一个补丁，
`ban #1-5` 的教训就是**"一个遍历同时服务拒绝与放行"必须拆**（memory 第 22 条同族）。

## D22 段（这条决定本票能不能自己收口）

`internal/risk/pathresolver*.go` 是**冻结契约文件承载的实现**。判据（R21④ 与 memory 第 20 条，我给自己上的锁）：
- **如果契约文本已经要求"路径必须先规范化/解析"，那修实现不是契约变更** ⇒ 你直接修，commit 正文引 `PLAN.md`/`SPEC-06` 的**原句 + file:line**（**不许写"契约要求"却引不出句子**——本仓已因为我这句加戏订正过一次，见 R21④）。
- **如果你要改契约文本的语义**（例如明写"禁止任何展开"）⇒ **停手**，把要改的那一句的原句、拟改句、影响面写进票面交回，**owner 批准之前不许动那个文件**。

## AC（1:1，裁决表 `docs/evidence/s1/102-*.md` 由验收方出，不是自裁）

- [ ] **AC#1** **先复现再修**：把验收代理的 PROBE B2 做成**包内可重跑的用例**（不是临时探针）：
      输入含 `%VAR%`/前导 `~` ⇒ 断言"结果树与请求树**必须同一棵**，否则必须报错"。
      ⚠ 这条用例在**修之前必须红**（红名+断言原文进票面）——它既是复现也是回归网。
- [ ] **AC#2** 选定 (A) 还是 (B)，并给**每个生产调用点**的处置表（`grep` 出所有 `Resolve(` 非测试命中，逐点写"要不要消费改写标记/是否拒绝"）。**不许只改 winsec 那条腿。**
- [ ] **AC#3** 变异：把修法退回"照旧展开且不记账" ⇒ AC#1 红；另做一发 **反向**：把 `Result` 的标记字段变成恒 `false` ⇒ **也必须有用例红**（否则 (B) 是纸面承诺）。
- [ ] **AC#4** 影响面回归：`go test -count=2 ./internal/risk/ ./internal/winsec/ ./internal/tools/ ./internal/memory/` 全绿，
      并逐条点名 SKIP/FAIL；`sh scripts/d22scan.sh` 纯净树 rc=0、`ban #6`/`ban #8` 对 `frontend/` 的数**不许降**。
      ⚠ `TestResolvePerCallBudget` 那条 1ms 墙钟预算**在负载下会假失败**（票 86）：量到红先判是不是它，**别把它记成回归，也别拿它当借口掩盖真红**。
- [ ] **AC#5** 若走 (B)：给"谁消费了这个字段"一条**静态可检查**的判据（例：`Screen`/`writeGate`/`PrivateDirAll` 三处必须读它），
      而不是"注释里提醒后人记得读"（票 96/票 90 两次教训：**靠注释与习惯的方向性，最后都要重做**）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；
**不 push**；不在仓内建 worktree（A38④，快照目录带会话后缀）；变异只在 `/tmp` 仓外快照做；
票面 append-only（改行前先读；标题前插段落要重抄标题，删除列必须 0）；四种假绿逐条点名；**编译失败不算变异**；
数字不达标写 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报；
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、`rules_gateway.go`、
`tools/d22scan/**`、`allowlist.txt`（5 行非注释，只许变短）。
⚠ 共树：`internal/winsec/`+`internal/secret/` 刚被 `agent-ticket89b` 改过（验收在跑），
`cmd/wisp/`+`internal/perm/` 刚被票 101 改过（验收在跑）⇒ **动这两个包之前先在票面登记**。

## AC#1 修前红名与断言原文（复现先于任何改码；本段写于改码之前）

用例：`internal/risk/pathresolver_expansion_test.go`（包内可重跑，不是临时探针）。
命令：`go test -count=1 -v -run TestC26ExpansionMustNotRewriteOntoAnotherTree ./internal/risk/`
读数：**`=== RUN` 3 条（1 父 + 2 子用例）· 2 子用例全 `--- FAIL` · 6 条断言红 · 0 SKIP**。

红名：
- `TestC26ExpansionMustNotRewriteOntoAnotherTree/percent_env_var_spelling`
- `TestC26ExpansionMustNotRewriteOntoAnotherTree/leading_tilde_spelling`

断言原文（helper `requireSameTreeOrRefusal`）：

```
FAIL-OPEN: <leg> acted on a tree the caller did not name and reported success (err = nil)
  caller's spelling : %s
  requested tree    : %s
  tree acted upon   : %s
  invariant         : 结果树与请求树必须同一棵，否则必须报错
```

本机实测三行（就是 PROBE B2 的形状，`err = nil`）：

```
caller's spelling : C:\…\a%WISP102PROBE_TARGET%b\artifacts
tree acted upon   : C:\…\aelsewhereb\artifacts                 ← 另一棵树
前导 ~           : ~\wisp102-tilde-probe\artifacts → C:\Users\swq\wisp102-tilde-probe\artifacts
```

三条生产腿各红两次（`%VAR%` 与 `~` 各一）：`c26Pipeline.Resolve`（装进 winsec 的那条适配器）、
`winsec.ResolvePath`（封存调用方看到的"成功/失败"）、`syncSet.resolveTarget`（fs.write 的 sync 判定）。
⚠ "请求树"用的是包内已获准的词法步 `lexCanonical`（`filepath.Abs`+`Clean` 仍只在 `pathresolver.go` 里），
没有新增第二套规范化，也没有给 D22 添表面。

## 修法判定：选 (B)；D22 未触发停手（契约确实 mandates 展开，原句在下面）

契约原句（逐字，能引出来才继续动手）：
- `docs/specs/SPEC-06-security-gatekeeping.md:50`：`展开(env / ~) → 绝对化 → Clean → 打开句柄取 GetFinalPathNameByHandle(VOLUME_NAME_DOS) 真实路径`
- `docs/PLAN.md:2375`：**修复 = C26 `PathResolver`，唯一入口**：
  `展开(env/~)` → `绝对化` → `Clean` → **打开句柄取 `GetFinalPathNameByHandle(VOLUME_NAME_DOS)` 得到真实路径**
- `docs/PLAN.md:1376`：**C26** **`PathResolver`** | **唯一**的路径规范化入口：展开 → 绝对化 → 取句柄真实路径 …
  **白名单与黑名单都必须经它**

⇒ 展开是契约**明写的 pipeline 第一步**（不是实现自己加的）。(A)「入口直接拒绝 `%VAR%`/`~`」等于把
契约的一步删掉，而且 `SPEC-06 §4.1` 的 A 档黑名单本身就用 `~/.git-credentials`、
`%APPDATA%\wisp\config.toml` 这种拼写（`docs/specs/SPEC-06-security-gatekeeping.md:63-66`），
入口拒绝展开会让"黑名单经 C26"这条落不了地 ⇒ **(A) 是改契约语义，本票不选，也不需要 owner 批文本**。
选 **(B)**：展开照旧（契约要求），但**改写必须显式记账**，凡拿它做安全决定的路径必须读这个账。
本票只动 `internal/risk/` 的实现与 `internal/tools/` 的消费侧，
`docs/PLAN.md`、`docs/specs/*.md` **一字未动** ⇒ D22 那道闸没有把我拦住，因为这是**修实现不是改契约**。

**需要谁协调：无。** `internal/winsec/`、`internal/secret/`、`cmd/wisp/`、`internal/perm/` 本票**只调用不改**
——拒绝发生在 `internal/risk/winsec_c26.go` 这条适配器腿上，winsec 原样传播 error（`resolve.go:113-119` 的
`%w` 包装），封存因此失败而不是"封了另一棵还报成功"；测试 import winsec 是只读引用
（`risk → winsec` 边本来就存在，见 `winsec_c26.go:3`）。

## 未勾原因：以下 AC 还在做

- **AC#2** 处置表（写完实现后逐点登记，含"要不要消费改写标记/是否拒绝"）。
- **AC#3** 双向变异（/tmp 仓外快照）。
- **AC#4** 门禁读数。
- **AC#5** "谁消费了这个字段"的静态判据。

next= 实现 (B)：`Result` 加改写记账 + 一个必须读账才能拿到可行动树形的闸 + 四条腿逐个消费 + 静态判据用例。

## Progress log（append-only）

- 2026-09-21 17:4x（编排者）：建票。来源是 `acceptor-ticket94` 的 R-a（它把它排在四条新账的第一条，
  并明确写"**不能原样对 owner 说'路径已解析'**"）。我给本票加了两道它没写的闸：
  ① **先判"展开是不是契约要求的"**再选修法（否则很容易把该保留的行为当 bug 删掉）；
  ② **D22 段**：动冻结实现文件时"修实现 ≠ 改契约"这条要自己引原句，引不出就停手交回。
  ⚠ 还有一条我明写在 AC#2 里：**不许只修 winsec 那条腿**——入口规范化的洞对所有判定都成立，
  只补一个调用方就是"同一不变式只修一半"（票 75/82 都栽过这个）。
  next= 派单（`internal/risk/` 现在无人写）；本票**优先于票 93/99/85**那三张门禁票，因为它是 fail-open。
