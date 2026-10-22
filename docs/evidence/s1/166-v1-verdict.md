# 票 166 终裁腿 `166-v1` · `AC#2`「四件最小出口」终裁 · 2026-10-22

腿＝**非实现者**（裁决者≠实现者）。⛔ 落任何 `.go`（仓内零字节；我的尺只落在仓外 `D:\tmp\wisp\20261022\166v1\`）、
⛔ 翻框、⛔ 改票名／`Status:`／`-done`／台账、⛔ push。被裁对象＝`166-r1`（`17cfe8cf`／`b4dfe165`／`b4102f1e`，
件 `docs/evidence/s1/166-r1-impl.md` 135 行）。台面：主检出 `D:\work\workspace\projects plans\Wisp`（`HEAD=e94a5b17` 时点，
`internal/ cmd/ scripts/` 干净）；`go1.27.1 windows/amd64`；`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`。
读数件目录 `D:\tmp\wisp\20261022\166v1\`（`logs/` 30 枚、`trees/` 6 棵、我的尺 245 行、`pristine/`），**只建不删**。

## 0. 三行判语

① **那处行为变更应当成立**（`ORDER BY` 逐字节未动、撤销把 `updated_at` 放回原值、撤销把行放回原位；⛔ 触契约文字）。
② **`AC#2` 能翻** —— 但票尾必须新增"四件生产调用者＝0/0/0/0/0、用户仍回不到昨天那次对话"这条具名残余（⛔ 抹）。
③ 它自报的**四件一件都不响**只在"票前基线 `2f41a80cb`"这个指称下成立；在"本腿修前 `ce84d53d`"下前三件**会响**、第四件**红而非哑**。

## 1️⃣ 票面前提新鲜度（对 HEAD blob 复量）

| 票面逐字 | 现跑 | 判 |
|---|---|---|
| `AC#6` 地基 `manage.go:93` 写 `m.now().Unix()` | 今天 **`:106`**，偏 **+13**。成因＝10-18 装配腿的守卫块**净 +13 行**（单一 hunk `@@ -48,0 +49,13 @@`；先插 `2f936413` +15、`7404c42f` 收成 +13；今日块体 `:49-60`，`var _ ManageStore` 在 `:60`） | **腐坏 1 处，它报得对** |
| `AC#6` 地基 `dao_session.go:125` 排 `ORDER BY updated_at DESC, id DESC` | `:125` 逐字命中 | 对 |
| `AC#6` 地基 `session.go:96` 由 `rand.Read` 铸 id | `:96` 逐字命中 | 对 |
| 六格＝五勾一未勾、`AC#2` 唯一未勾 | 勾 5／未勾 1（`:25`） | 对 |
| `Status:` 的派单前置（票 158 验收＋161 结案） | `158-…-done.md`、`161-…-done.md` 都在盘 | 理由已消失 |
| 基线 `2f41a80cb`／`5b6338f6`／`ce84d53d` | 三枚都在（`5b6338f6` 与 `ce84d53d` 在本腿两枚文件上 blob **同值**：`dao_session.go= f1dc8e83`） | 对 |

★`Status:` 行逐字仍写"立而不派"——**未随派单更新**；⛔ 我改（票名／`Status:` 不属我写面）。

## 2️⃣ ★★核心判点一：那处行为变更**应当成立**

| 判据 | 我的读数 |
|---|---|
| `ORDER BY` 有没有被碰（`AC#6`ⓑ 那一格） | **没碰**。四条 `ORDER BY` 文本在 `ce84d53d`／`17cfe8cf`／`HEAD` 三处逐字节同（合并 sha256 `a9b7715b…`）；两条排序 SQL（`:125` List／`:144` Search）所在 113–146 行 `diff` **空**；`17cfe8cf` 五个 hunk＝`+191` 注释／`-192` `now :=`／`202` stamp／`205` SQL／`+263` 注释，无一触及排序 ⇒ **⛔ 碰 `AC#6`ⓑ（tie-break）** |
| 撤销后 `updated_at` 是否回到归档前值 | **是**。我的尺 `i4h` 判树绿；修前树红：`undo rewrote the activity stamp: updated_at 1792680288, want 1792659720` |
| 撤销后行是否回到原位置 | **是**。`i4i` 比"归档前序列 vs 撤销后序列"逐 ID 相等（判树绿）；修前树红：`before [c5a5… 0c2b…] after [0c2b… c5a5…]` |
| 契约面 | `updated_at` 在 `docs/PLAN.md`／`docs/specs/**` 只对 `profile`／`memory` 两表有文字（`PLAN:2706`、`SPEC-02:38/40/159`、`SPEC-05:143`）；**`session.updated_at` 无任何契约文字** ⇒ 这处行为变更 ⛔ 触契约、⛔ 动阈值/golden |
| 语义 | `updated_at`＝最后**活动**；归档是掩码不是活动（326 AC#1 的"追加轮次才刷"照旧）；撤销＝精确逆＝主流 harness 的"归档/取消归档回到时间序原位、⛔ 顶到最前" |

**判语：成立、该落。** 两处残余具名（⛔ 抹）：① `RenameSession` **仍**写 `updated_at` ⇒ 改名仍会把行顶到最前，
这是另一条语义问题、⛔ 本裁范围；② 若归档期间有人 `AppendSessionTurn`（326 AC#1 会刷），撤销恢复的是"归档那一刻"的值——
"精确逆"只对归档这一跳成立，而这正是对的（新消息本就是活动）。

## 3️⃣ ★★核心判点二：四件在未修码上「响不响」——两个指称都报

- **指称甲＝票前基线 `2f41a80cb`**（`git archive` 仓外树，1003 枚 `.go`）：七枚符号 `git grep -c` 全 **rc=1／0 命中**
  （`ListSessions`／`SearchSessions`／`RenameSession`／`SetSessionArchived`／`NewManager`／`ManagerOptions`／`SessionListOptions`）；
  `internal/session` 无 `manage.go`、`internal/memory` 无 `dao_session*.go`；`go doc ./internal/memory Store` 63 行里唯一 session 行＝
  `ListGrantsBySession`（D45 台账）；我的探针在该树 `go vet ./internal/session/` **rc=1 `undefined: Manager`**（RUN=0，编都编不过）。
  正控同树：`ListGrantsBySession`／`NewLedger` 各 1 文件命中；同尺打 HEAD 五符号全命中 ⇒ **尺灵**。⇒ 四件**都不响**。
- **指称乙＝本腿修前 `ce84d53d`**：我的四枚用例 **Item1/2/3 PASS**（前三件会响）、**Item4 FAIL 3 行**（`i4c`／`i4h`／`i4i`）。
  ⇒ 它那句"四件一件都不响"只在甲下成立；⛔ 读成"四件在修前全哑"。
- ★**尺的射程更正（复核成立）**：Go 1.27.1 上 `go test -run '^TestZZNoSuchName166V1$'` **rc=0**，输出
  `ok … [no tests to run]`（`-v` 时另有 `testing: warning: no tests to run` + `PASS`）⇒ **⛔ 拿 rc 判响不响**。

## 4️⃣ 我自造的反向判据（⛔ 只看它的 T1–T8；台面＝仓外导出树，物理改字节、串行、逐发还原）

我的尺＝`zz_166v1_verdict_test.go`（**245 行、4 枚用例**，只存在于仓外；判树 4/4 PASS rc=0）。四发变异：

| 发 | 摘掉／改掉哪一跳 | 变异读数（逐字红行） |
|---|---|---|
| **M-A** | 未修码本体（`ce84d53d` 树，码**未手改**） | rc=1；`:208 i4c archive rewrote the activity stamp: updated_at 1792659720 -> 1792680288`；`:229 i4h undo rewrote the activity stamp: updated_at 1792680288, want 1792659720`；`:236 i4i undo changed the by-time order: before [sess_c5a5… sess_0c2b…], after [sess_0c2b… sess_c5a5…]` |
| **M-B** | `dao_session.go:201` `if archived` → `if true`（撤销不清 `archived_at`） | rc=1；`i4g undo left archived_at set: 1792680441`／`i4i`（撤销后序列少一行）／`i4j undo did not return the conversation to search: []` |
| **M-C** | `:125` `ORDER BY updated_at **ASC**`（顺序轴反向） | rc=1；`i1a by-time list=… want …`／`i1b by-workspace list`／`i1c titles out of order`＋Item4 `premise` Fatal |
| **M-D** | `:137` `if false && !opts.IncludeArchived`（搜的归档掩码摘掉） | rc=1；`i4e archived row still surfaces in search: [sess_0dd6…]` |

- **三读数**：变异 rc=1（红名如上）／**还原 rc=0**（`cmp=0`）／**原样 rc=0**（4/4 PASS）。
- **还原核**：`internal/memory/dao_session.go` sha256 `4bb05526ce8651e94190125505b4f87ec359d50826e9aebdccf75a8b1aec8268`
  （**64 位＝SHA-256，`sha256sum` 算**）＝提交版；SHA-1 `75c5685ed688bce306e996b4ecd19e8deaebb39f`（**40 位＝`git hash-object` 算**）＝`HEAD` blob。
  ★两算法分开比、⛔ 混比（射程清单第 9 例）。修前 blob sha256 `ec6ee530b190795093829eca07e6003bf866db5da3e64c47e6d0b0b1bcd1d2b2`
  与 `166-r1` 自留的 `blob-dao-session-before.go` **同算法同值**（`cmp` SAME）⇒ 它那三枚 64 位哈希与我这把尺同算法、可比。

## 5️⃣ 门禁复跑（四数分开、⛔ 相减）

台面＝主检出；G1–G5 于 **22:51–22:57** 跑，跑前 `git status --porcelain -- internal/ cmd/ scripts/` **空**（当时 `HEAD=e94a5b17`）。
★此后别的腿又落了两笔（`c0159ac2`／`74f29a44`，含 `internal/config/tiers_223r3_test.go` 等 4 枚 `.go`）⇒ 下表只对**该时点**成立；
G6乙 是 23:2x 的工作树，故它量到的与下表**不是同一棵树**（诊断面，⛔ 当判据）。

| 门 | 命令 | rc／读数 |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **rc=0**（0 行输出） |
| G2 | `go vet ./internal/session/ ./internal/memory/` | **rc=0**（0 行） |
| G3 | `go test -count=1` 两包各一发 | `session` **rc=0 顶层 RUN=27 PASS=27 FAIL=0 SKIP=0**（子测试 RUN=17）／`memory` **rc=0 顶层 RUN=45 PASS=43 FAIL=0 SKIP=2**（子测试 RUN=35） |
| G4a | `GOTMPDIR=D:/tmp/… sh scripts/d22scan.sh` | **rc=2**（1 条 GUARD／0 条 clean／**0 条 `=== RUN`**）＝守卫拒跑，见 §8 |
| G4b | `sh scripts/d22scan.sh`（GOTMPDIR 未设） | **rc=0** `clean - no D22 ban violations`；正控 `runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`；我这发 ban #8 `internal/`=**620**、`cmd/`=148（`166-r1` 22:12 那发报 619——差 1 是**盘上 `.go` 枚数随时点**，⛔ 相减、⛔ 读成门禁差） |
| G5 | `sh scripts/check-path-length-budget.sh` | **rc=0** `VERDICT GREEN`（denominator 17594、over-the-hat=61、in-roster=61、wall interval=0） |
| G6甲 | `git archive 17cfe8cf` 树 + `gofumpt -l`（v0.12.0，**判据**） | 20 行、**全部 `.scratch/`**（non-scratch=0）；**本腿两枚 0 命中**（归属判据用 `[\\/]`）；树内 `.go`=1157 |
| G6乙 | 工作树 `gofumpt -l`（**只诊断**） | 54 行、本腿两枚 0 命中；3 行非 `.scratch`＝`cmd\wisp\models.go`、`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`；rc=2＝有不可解析样本（诊断面⛔ 当判据） |

★G3 的 RUN 尺形：`^=== RUN` 会把子测试算进去（44／80）；我报的是**顶层**（`^=== RUN   [A-Za-z0-9_]+$`）。
★甲乙量的是**两个对象**（blob 树 vs 文件系统）⇒ 54≠20 **⛔ 相减**。

## 6️⃣ 「新增红 0」独立核

- **两侧＝两棵真树（⛔ overlay）**：`judged`=`17cfe8cf`（1157 枚 `.go`）／`prefix`=`ce84d53d`（1156 枚）；`diff -rq` 只有两处：
  `internal/memory/dao_session.go` 差、`internal/session/ticket166r1_ac2_gaps_test.go` 只在 judged ⇒ 差集只可能有本腿因果。
- **我现跑**：after **22** 名／before **21** 名／**19 同名**。only-before＝`TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1`＋`TestGateHalfDuplexClosedZeroFrames`；
  only-after＝`TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`／`TestTenOpsInOneToolCallGetOneConfirm`／`TestVetoTravelsBackToTheLoopAsAFailure`。
- **五名定点复跑两侧**：四名**两侧都 PASS**（含 audio 那枚；它在 `166-r1` 那对里是 **after 独有**、在我这对里是 **before 独有** ⇒ 两对独立读数方向相反＝抖动）；
  `…342r1` **两侧都 FAIL**，红行 `ancestorworld_342_348_windows_test.go:386 … settle the world back: remove …\data\sub: The directory is not empty`（清理面、两树同形）。
  ⇒ **归因本腿的新增红＝0**。
- ★**我的 22/21 ⛔ 与它的 10/10 横比**：我的两树是 `git archive` 导出树（**无 `.git`**、在 D: 卷）。同尺同树对照：
  `TestReadGitOnThisRepositoryIsSelfConsistent`（panel）主检出 **PASS**／导出树 **FAIL**（`not_a_repo`）；`TestTicket252P1AllowlistLegs`（tools）
  主检出 **PASS**／导出树 **FAIL**（8.3 短名前置条件缺失）⇒ **差集有效、绝对数⛔ 横比**。
- **结构面**：`SetSessionArchived` 的唯一调用者是 `Manager.Archive/Unarchive`；`internal/session` 的生产导入者只有
  `cmd/wisp/run.go`、`cmd/wisp/resident_approval_windows.go`（只碰 `Ledger`/`Mint`）⇒ 两枚文件**没有能力**让 `{session, memory}` 之外变色。

## 7️⃣ 产品面残余 → `AC#2` **能翻**

- **我现量的五个零**：①`session.NewManager` **0**（`rc=1`；正控 `func NewManager` 3 命中、`session.NewLedger` 在 `cmd/wisp/run.go:491` 命中）；
  ②非 `_test.go` 的 Manager 方法调用点 **0**（同尺 `_test.go` **70** 命中）；③`memory.Session` in `internal/panel`=**0**、`cmd/wisp`=**0**；
  ④`ListSessions` 两处=**0**；⑤`SessionByID` 两处=**0**（另 `SetSessionArchived` 两处=**0**）；CLI `"sessions"|"history"` in `cmd/wisp/main.go`=**0**。
- **判语：`AC#2` 能翻。** 依据＝票面 `Status:` 自陈"本票只做 **Go 侧**"＋`AC#2` 逐字要的是"四件最小出口＋每件答响不响"：
  四件在 Go 侧都在、都经 `ManageStore` 装配（`manage.go:60` 编译期守卫＋七枚 `TestAC2*` 逐名 PASS，我另起树复跑同名全绿）、
  每件都有**今天能转红**的用例（我四发变异逐件红，§4），`AC#2` 逐字"归档带撤销"的撤销**现在是真撤销**。
- **硬条件（⛔ 抹）**：①票尾必须新增一条具名残余："四件生产调用者＝**0/0/0/0/0**；**用户仍回不到昨天那次对话**；
  接线＝`cmd/wisp`＋`internal/panel`（均归别的腿）"；②⛔ 把 `AC#2` 成立读成"产品面成立"；③⛔ 拿 `go build` rc=0 或"用例绿"冒充接线；
  ④本格与接线票**分开记**（`-done` 与否由编排者定，本腿⛔ 改票）。

## 8️⃣ G4 第一发 `rc=2` 的解释核：**成立**

守卫在 `tools/d22scan/runtests.sh:169-201`（ticket 311 `AC#1` shape B）：参照卷取 `%LOCALAPPDATA%`（C:），`GOTMPDIR` 落在别的卷 ⇒
打三行 GUARD 后 `exit 2`（`:200`）。我复跑：`GOTMPDIR=D:/tmp/…` **rc=2**（1 GUARD／0 clean／**0 `=== RUN`**＝测试进程一次都没起）；
摘掉 `GOTMPDIR` **rc=0**。⇒ **⛔ 把"守卫拒跑"读成"门禁红"**。

## 9️⃣ 写面纪律（逐笔 `git show --name-only --format=""`）

| 提交 | 路径 | `.go` 枚 | 越界 |
|---|---|---|---|
| `17cfe8cf` | `internal/memory/dao_session.go`、`internal/session/ticket166r1_ac2_gaps_test.go` | **2** | 0 |
| `b4dfe165` | `docs/evidence/s1/166-r1-impl.md` | 0 | 0 |
| `b4102f1e` | `docs/evidence/s1/166-r1-impl.md` | 0 | 0 |

越界尺：`^(internal/session/|internal/memory/|docs/evidence/s1/166-r1-[^/]*\.md$)` ⇒ 三笔**全 0**；对 `AC#5` 禁区 pattern
（`PLAN.md`／`specs/`／`risk/`／`panel/`／`agent/approval/`／`thresholds.go`／golden／`allowlist.txt`／`frontend/`／`design/`）也**全 0**。
**正控（同尺）**：`394cc212`（223-r3）`.go`=**1**、越界=**1** ⇒ 尺看得见 `.go`、也看得见越界；`17cfe8cf` `.go`=2／越界=0 ⇒ 尺不虚报。

## 10. 残余清单（⛔ 抹；全部不属本腿写面）

1. **产品面不成立**：四件生产调用者 0/0/0/0/0（§7 五把尺）。**记法**：票尾具名 bullet（编排者代笔）＋本表 §7＋新接线票
   （`cmd/wisp` 装配行＋CLI 出口＋`internal/panel` 会话面），⛔ 写成本格"已完成"。
2. `RenameSession` 仍刷 `updated_at`（§2 残余①）＝另一条语义问题。
3. 归档期间 `AppendSessionTurn` 会刷 `updated_at`（§2 残余②）。
4. `Status:` 行仍写"立而不派"（§1）＝票面待编排者更新。
5. `166-r1` 的件里"父提交 `5b6338f6`"不是真父（真父 `ce84d53d`）；两者在本腿两枚文件上 blob **同值** ⇒ 它的 overlay 比对**不受影响**、⛔ 判它没还原。

## 附：本件的尺（自陈射程）

- 我的反向判据尺＝`zz_166v1_verdict_test.go`（245 行，仓外）；量的是**生产 `Manager`＋真 `wisp.db`**，⛔ mock；时钟步进 ⛔ 靠同秒 tie-break。
- 我的名册尺＝`go test -count=1 ./internal/...`（30 包、两棵导出树）；名册取**主检出式样**（`git archive` 树），⛔ 取工作树。
- 本件提交号**不写进本件**（防自我推翻）；见 `git log` 末笔（主题逐字以 `票 166 终裁腿 166-v1` 开头）。
