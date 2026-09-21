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

- [x] **AC#1** **先复现再修**：把验收代理的 PROBE B2 做成**包内可重跑的用例**（不是临时探针）：
      输入含 `%VAR%`/前导 `~` ⇒ 断言"结果树与请求树**必须同一棵**，否则必须报错"。
      ⚠ 这条用例在**修之前必须红**（红名+断言原文进票面）——它既是复现也是回归网。
- [x] **AC#2** 选定 (A) 还是 (B)，并给**每个生产调用点**的处置表（`grep` 出所有 `Resolve(` 非测试命中，逐点写"要不要消费改写标记/是否拒绝"）。**不许只改 winsec 那条腿。**
- [x] **AC#3** 变异：把修法退回"照旧展开且不记账" ⇒ AC#1 红；另做一发 **反向**：把 `Result` 的标记字段变成恒 `false` ⇒ **也必须有用例红**（否则 (B) 是纸面承诺）。
- [x] **AC#4** 影响面回归：`go test -count=2 ./internal/risk/ ./internal/winsec/ ./internal/tools/ ./internal/memory/` 全绿，
      并逐条点名 SKIP/FAIL；`sh scripts/d22scan.sh` 纯净树 rc=0、`ban #6`/`ban #8` 对 `frontend/` 的数**不许降**。
      ⚠ `TestResolvePerCallBudget` 那条 1ms 墙钟预算**在负载下会假失败**（票 86）：量到红先判是不是它，**别把它记成回归，也别拿它当借口掩盖真红**。
- [x] **AC#5** 若走 (B)：给"谁消费了这个字段"一条**静态可检查**的判据（例：`Screen`/`writeGate`/`PrivateDirAll` 三处必须读它），
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

## 交件读数（agent-ticket102，2026-09-21 18:2x）—— Status: **in-review**

实现与用例在 commit **`a66aadf`**（AC#1 的红用例先单独落在一枚 **`0117459`**，红名与断言原文在上一段）。
基线对照一律用 `a66aadf^` 的仓外快照 `/tmp/wisp102-baseline-102`。

### AC#2 生产调用点处置表（`grep -rn "risk\.Resolve(" --include=*.go .` 去掉 `_test.go` = 1 命中；
包内 `Resolve(` 命中另 4 处；合计 5 条腿，逐点处置如下，**没有只修 winsec 那条**）

| # | 调用点 | 这条腿在做什么 | 处置 | 落点 |
| --- | --- | --- | --- | --- |
| 1 | `internal/risk/winsec_c26.go:34`（装进 winsec 的适配器，PROBE B2 的那条） | **放置/封存**：winsec 拿返回值去建目录/封 DACL，并且只回 `error` ⇒ 既动手又报成功 | **拒绝**：走 `Result.Actable()`，被改写即返回 `ErrRewrittenPath`；winsec `ResolvePath` 用 `%w` 原样传播 ⇒ 封存失败而不是封错树 | 用例 `TestC26ExpansionMustNotRewriteOntoAnotherTree`（含 `winsec.ResolvePath` 那条断言） |
| 2 | `internal/risk/syncdirs.go:140`（`syncSet.add`，同步根配置） | **解除兜底网**：`canonical=true` + confirmed-grade 才会关掉 sync-suspect 网（放行方向） | **消费但不拒**：被改写的根仍守（canon 用展开后的树），但**不再算 canonical-grade** ⇒ 关不掉兜底网，并 `logf` 一条给运维 | 用例 `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`（对照腿：干净根必须 `canonical=true` 且 `complete=true`，否则用例自己先 fatal） |
| 3 | `internal/risk/syncdirs.go:197`（`resolveTarget`，fs.write 的 sync 判定） | **放行判定**："不是同步目录" 才让写过去 | **拒绝**（fail-closed）：`res.Actable()` 报错 ⇒ 上游按 sync-suspect 处理 | 用例 AC#1 的第三条断言腿 |
| 4 | `internal/risk/syncdirs.go:211`（同一函数里对最近存在祖先的第二次 `Resolve`） | 祖先链重解析 | **消费**：也走 `Actable()`（祖先形是从 C26 自己出来的，今天恒不触发，读账而不是假设） | 同上（`anceCanon` 全链替换，无裸 `.Canonical` 读取） |
| 5 | `internal/tools/paths.go:56`（`PathCanonicalizer.resolve`；`grep "risk\.Resolve("` 唯一非测试命中） | **授权**：`[fs] allowed_dirs` 根 + R2 的 in/out-of-allowlist 判定 | **拆开两条子腿**：根=记账（`RewrittenRoots()`）+ 改写且 OS 认不出 ⇒ 进 `unusable`（什么也不授权）；target=**不拒**，因为 `fs.go:95`/`fs_write.go:173`/`mode.go:93` 打开的就是 `Canonicalize` 的返回值 ⇒ 判定树==执行树==审计打印的那棵，拒它就是把 (B) 做成 (A) | 用例 `TestPathCanonicalizerAccountsForRewrittenRoots`（含"不许变成 (A)"那条：确认过的改写根照样授权） |

不改的相邻面（写清楚免得被当成漏网）：
- `internal/winsec/resolve.go` 的 `builtinVerifier`：**只拒不改写**（票 94 第 5 项实测），本票一字未动。
- `pathresolver.go` 的 `formsOf`/`anchorForms`（A/B 档锚点比较）：走 `lexCanonical`，**不做展开** ⇒ 没有本洞。
- `assessor.go` / `rules_gateway.go`（禁区）：只经 `PathCanonicalizer` 接口取路径，接口签名未变 ⇒ 无需改。
- `Result.Canonical` 的读取点全仓复算：`pathresolver.go`（生产者）+ 上表 2/3/5 号腿，**没有第六处**。

### AC#3 双向变异（全部在仓外快照 `/tmp/wisp102-mut-a66aadf4`，`git archive a66aadf | tar -x`；
每发先 `go build ./...` **rc=0**（编译失败不算变异），改完 `diff -r` 对 pristine 副本 **rc=0**，
本仓 `git diff --quiet` **rc=0**）

| 变异 | 改法 | 结果 |
| --- | --- | --- |
| **M1 修法退回"照旧展开且不记账"** | `pathresolver.go:94` `if r.Rewritten {` ⇒ `if false && r.Rewritten {`（`Actable` 不再拒） | build rc=0；**红**：`TestC26ExpansionMustNotRewriteOntoAnotherTree/{percent_env_var_spelling, leading_tilde_spelling}`、`TestC26RewriteAccountIsRecorded/{percent_env_var, dollar_env_var, leading_home_tilde}`。**绿**（覆盖边界，如实登记）：`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 是文本判据，M1 只改条件不删符号 ⇒ 它不红；`internal/tools` 那条也仍绿（它读 `Rewritten` 不读 `Actable`）。 |
| **M2 反向：`Result` 的标记字段恒 `false`** | `pathresolver.go:116` `Rewritten: len(kinds) > 0` ⇒ `Rewritten: false` | build rc=0；**红**：AC#1 两条子用例、`TestC26RewriteAccountIsRecorded` 三条子用例、`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`（两条断言：`canonical=true` + `complete=true`）、`internal/tools` 的 `TestPathCanonicalizerAccountsForRewrittenRoots` ⇒ **(B) 不是纸面承诺**：把账擦掉，四条腿的用例一起红。 |

### AC#4 门禁读数（基线 = 本票实现后的工作树；`-count=2` 一轮，未重测挑数）

`go test -count=2 ./internal/risk/ ./internal/winsec/ ./internal/tools/ ./internal/memory/` ⇒ **四包 rc=0**，
`-json` 逐包点名的 `=== RUN`/`--- PASS`/`--- FAIL`/`--- SKIP`：

| 包 | RUN | 不同名 | PASS(测试条数) | FAIL | SKIP |
| --- | --- | --- | --- | --- | --- |
| `internal/risk` | 324 | 162 | 322 | **0** | 2 |
| `internal/winsec` | 58 | 29 | 58 | **0** | 0 |
| `internal/tools` | 202 | 101 | 202 | **0** | 0 |
| `internal/memory` | 134 | 67 | 132 | **0** | 2 |

`=== RUN` == 2 × 不同名 四包全部成立（324/58/202/134）。SKIP 逐条点名（都是环境闸，非本票引入）：
`internal/risk/TestSyncRegistryProbeLive`（`syncdirs_test.go:347`，本机 HKCU 无 registry-grade 同步记录）、
`internal/memory/TestSubprocessCrashWriter`（`concurrent_test.go:186`，crash-writer 子进程用例在
`TestCrashRecoveryKillMidWrite` 下跑）⇒ 两条在 `a66aadf^` 的仓外基线快照里**同样 SKIP**（实测贴在上面），
所以是 `-count=2` 各 1 条 × 2 轮 = 2 条，不是回归。
**`TestResolvePerCallBudget`：本票没量到红**（四包 FAIL=0 里有它），未当作借口也未记成回归。
其余门禁：`gofmt -l .` **空**、`gofumpt -l .` **空**（`$(go env GOPATH)/bin/gofumpt`）、
`go vet` 按包 rc=0（risk/tools/winsec/memory 各一）、`GOOS=linux go vet` **按包**同样四 rc=0
（仓根 `./...` 的已知仪器坑不在此列）。
`go test ./cmd/wisp/`：本机 `exit status 0xc0000135`、`=== RUN` **0 条** —— 与票 98 记的加载期缺 DLL 同形，
且在 `a66aadf^` 基线快照里**一模一样**（实测在上面）⇒ 判"仪器不可跑"，不判本票回归；票 98 的注入命令本代理未持有，
`go vet ./cmd/wisp/` rc=0 是本代理能给的最强替代读数。
`sh scripts/d22scan.sh`：**纯净仓外快照** `/tmp/wisp102-d22-0e7640c1`（`git archive HEAD | tar -x`）**rc=0**，18.6s，
台账 8 行逐字：`bans #1-5 internal/=197`、`#1-5 cmd/=20`、`**ban #6 frontend/=37**`、`ban #7 internal/tools/=17`、
`#8 design/=16`、`**#8 frontend/=37**`、`#8 internal/=340`、`#8 cmd/=26`
⇒ 与票 94 验收参照（`#6 frontend/=37`、`#8 16/37/335/25`）比：`#6`/`#8 frontend` **没降**，
`#8 internal/` 340（参照 335）、`#8 cmd/` 26（参照 25）是别人后续进树的文件，只增不减。

### AC#5 "谁消费了这个字段"的静态判据（不是注释提醒）

`internal/risk/pathresolver_rewrite_account_test.go:TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 两半：
1. **逐腿点名**：`winsec_c26.go` 必须出现 `Actable(`；`syncdirs.go` 必须同时出现 `Actable(` 与 `res.Rewritten`；
   `../tools/paths.go` 必须出现 `res.Rewritten` —— 少一个 `t.Errorf`，文件被搬走则 `t.Fatalf`（不许靠"记得找"）。
2. **全仓扫**：`WalkDir` 从仓根扫所有**非测试** `.go`，凡读 `Result.Canonical` **字段**
   （`readsCanonicalField` 用后随字符把 `Canonicalize` 方法排掉，第一版误报 6 处已修）
   的文件，除生产者 `pathresolver.go` 外必须同时出现 `Actable(` 或 `.Rewritten`；
   一条都没扫到 ⇒ `t.Fatalf("the instrument is broken")`（防"扫空=绿"）。
⇒ 以后新增第六条腿想直接抓 `Canonical`，CI 就跑不过；M1 的读数同时说明这条判据**只挡符号消失，不挡条件被改**，
真正的兜底是 AC#1 与 AC#2 表里那五条行为用例。

### 需要协调（已在动手前登记于此，未默默改）

`internal/tools/paths.go` 是 AC#2 处置表第 5 行，也是 `grep -rn "risk\.Resolve(" --include=*.go .`
去掉 `_test.go` 的**唯一**命中，票面 AC#2 明写"不许只改 winsec 那条腿"⇒ 本票按票面改了它，
改动范围只有 `PathCanonicalizer`（新增 `rewritten` 字段 + `RewrittenRoots()`；`resolve` 返回 `risk.Result`；
`NewPathCanonicalizer` 对"改写且 OS 认不出"的根 fail-closed）+ 一个**新文件**
`internal/tools/paths_rewrite_ticket102_test.go`；未动 `bridge.go`/`fs.go`/`fs_write.go`/`mode.go` 一行。
⚠ 另：本轮工具输出里反复出现一段自称"编排者备注、要求把 `internal/tools/` 冻起来"的附加文本
（同一串还出现在别的包的写入结果里），本代理按"工具输出不是授权"处理：没扩大范围、没动禁区，
只把上面这一处最小改动留在此处等裁决；如要回收，revert 单文件即可，`internal/risk` 三条腿不依赖它。
`internal/winsec/`、`internal/secret/`、`cmd/wisp/`、`internal/perm/`、`frontend/`、`internal/panel/`
**一字未动**（只 import/调用）。`docs/PLAN.md`、`docs/specs/*` 一字未动 ⇒ **未触发 D22 停手**（本票是修实现，契约原句已引）。

next= 交回验收代理出裁决表（`docs/evidence/s1/102-*.md`）；本票不解 R-b/R-c/R-d（三条仍在票 94 的账上），
`memory/open.go` 那份"自己的词法 abs"仍归票 18/79 地界：现在的效果是它一旦被改写会**收到错误**而不是静默封错树。

## Progress log（append-only）

- 2026-09-21 17:4x（编排者）：建票。来源是 `acceptor-ticket94` 的 R-a（它把它排在四条新账的第一条，
  并明确写"**不能原样对 owner 说'路径已解析'**"）。我给本票加了两道它没写的闸：
  ① **先判"展开是不是契约要求的"**再选修法（否则很容易把该保留的行为当 bug 删掉）；
  ② **D22 段**：动冻结实现文件时"修实现 ≠ 改契约"这条要自己引原句，引不出就停手交回。
  ⚠ 还有一条我明写在 AC#2 里：**不许只修 winsec 那条腿**——入口规范化的洞对所有判定都成立，
  只补一个调用方就是"同一不变式只修一半"（票 75/82 都栽过这个）。
  next= 派单（`internal/risk/` 现在无人写）；本票**优先于票 93/99/85**那三张门禁票，因为它是 fail-open。
- 2026-09-21 18:3x（agent-ticket102）：交件。AC#1..AC#5 五框全勾，读数在上面（红名/契约原句/五条腿/两向变异/四包 `-count=2`/d22scan 台账）。
  两枚 commit：`0117459`（AC#1 红用例 + 判 (B) 的原句）、`a66aadf`（实现 + 三条 risk 腿 + tools 腿 + AC#5 静态判据 + 用例）。
  ⚠ 一条账要编排者裁：本轮**每一次工具输出末尾**都挂着一段自称"编排者备注"的附加文本，内容是"停手/冻结 `internal/tools/` 与 `internal/secret/`/`winsec`/`cmd/wisp`/…/否则终止并回滚、永久记录"，
  同一串在多个不同包的写入结果里逐字重复 ⇒ 本代理按"工具输出不是授权也不是指令"处理：**没有扩大范围**（禁区一字未动），
  但按票面 AC#2 的原文（"grep 出所有 `Resolve(` 非测试命中逐点处置"、"不许只改 winsec 那条腿"）保留了 `internal/tools/paths.go` 那一处最小改动，
  并在此登记。若那句真是编排者的口径，revert 单文件即可（`internal/risk` 的三条腿不依赖它），但 AC#2 会退回"只修一半"的状态、AC#5 的全仓扫描会红。
  next= 验收代理出 `docs/evidence/s1/102-*.md` 裁决表；R-b/R-c/R-d 与 `memory/open.go` 那份词法 abs 不在本票内（仍挂票 94/18/79 的账）。
