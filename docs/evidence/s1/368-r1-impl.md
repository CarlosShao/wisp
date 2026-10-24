# 票 368 / 腿 `368-r1` 落地（写腿）— DAO 层给三列时间戳兜底：**booked 行带 `decided_at`+`ended_at`、pending 行带 `started_at`**

- 腿＝`368-r1`（**写腿**）；开工 anchor HEAD＝`75c2395c3b9dc21d6467dbda5b1ec6b60465cf5b`；★裁形＝**〔丁〕在 DAO 层 `InsertToolCall` 兜底**（编排者已拍，本腿未改形）。
- 提交：**`c0eea973b`**（父＝`29c9d2b20`，编排者清索引残留那一笔，只动 `docs/reports/pending-and-issues.md`）；★`git show --name-status` 只有本件 4 枚路径；⛔ push。
- 写面＝`internal/memory/**`（3 枚 `.go`：`dao_toolcall.go`／新增 `dao_toolcall_368_test.go`／`migration_fkcheck_366_test.go` 里**一枚夹具的播种行**）＋本件。★`internal/{agent,tools,cmd}/**`、`scripts/**`、`ddlV*`、`thresholds.go`、golden、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**` **零字节**。
- ★★三枚硬边界**一条未越**：①`retention.go:154-155` 那条 `COALESCE` SQL **一字未动**；②**未加** `foreign_key_check` 硬门槛；③`ddlV*`／12 枚名册／`196-r2` 护栏／`196-r3` 的 `migrationDSN` **一字未动**。
- 探针／变异**全在仓外**：`D:\tmp\wisp\20261024\368r1\{tree-before,tree-after,probes,mut,out}`；★变异一律 `-overlay` 换 `dao_toolcall.go` ⇒ **无工作树突变窗口**（`git hash-object` 变异前后同一 `5b4e0e02…`）。
- 命令环境：`GOFLAGS=`、★**只** `unset GOTMPDIR`（`TMPDIR/TEMP/TMP` 未动 —— 本包会中假红）；rc 一律 `cmd > out 2>&1; rc=$?`（Git Bash 的 `$?`，⛔ 走管道）。

## §1 为什么是〔丁〕（★推翻票面〔甲〕〔乙〕〔丙〕的那三条读数）

普查件 `docs/evidence/s1/368-a1-census.md`（149 行、`bac7d5a29`）**四条**直接决定了形：

1. ★★**§7 第六条路径**：`tools.Bridge.book`（`internal/tools/bridge.go:1171`，构造 `:1210-1220`，`:1225` **`if _, err := b.j.InsertToolCall(ctx, tc); err != nil`** ⇒ 返回 id 丢掉）**三列一个都不设**、而 `Decision/Outcome` 有值 ⇒ **这行不可能有更新者**。
2. ★★**§7 生产装配**：全仓 `InsertToolCall` 非测试调用点**只有两个**（本腿现跑 = `agent/journal.go:75`、`tools/bridge.go:1225`）；而 `cmd/wisp/run.go` 今天给 `tools.New{Journal: mem}`（**`:777`**，普查记 `:769`，漂 +8）＋给 `agent.Options{Journal: nil}`（**`:1118`**，普查记 `:1080`，漂 +38；注释逐字 "the loop's own tool_call rows are deliberately NOT booked: the bridge already writes the authoritative one"）⇒ **今天生产上的 `tool_call` 行只由 `book` 写**。
3. ★★**§8 判据 3**：〔甲〕改 `COALESCE` 判据＝**读侧**替写侧擦屁股；〔乙〕让 `startCall` 写 `started_at`＝**只治 loop 侧**，而 loop 侧今天 `Journal: nil` ⇒ **两形对 `book` 的行零效果**。
4. ★**§2 附读数**：`started_at` **生产上零写入者**（本腿现跑：全仓 `UPDATE tool_call` 只有两条 —— `dao_toolcall.go:113` 写 `decided_at`、`:142` 写 `ended_at`；**没有一条写 `started_at`**）。
   ⇒ ★**"某列永远 NULL"是写者的责任、⛔ 读侧 `COALESCE` 的责任** ⇒ 落点只能是 `InsertToolCall` 本身：**一处改动，两个写者（含将来任何写者）全覆盖**。

## §2 「三列分别该在什么时候有值」——我先量、后判

**量到的既有语义**（不是猜）：`memory/models.go:81-85` 把 `Decision`／`Outcome` 的 `""` 明写成 **pending**；`DecideToolCall`（`:106-129`）写 `decided_at`＝决定时刻；`FinishToolCall`（`:136-158`）写 `ended_at`＝结束时刻；**`started_at` 无任何写入者**（上面 §1.4）。★另有一条**读者侧**旁证：`privacy.go:103-108` 取时刻的优先序是 **`started_at` → `decided_at` → unix 0** ⇒ 今天 pending 行在隐私页上显示成 **1970**，作者原本就指望 `started_at` 是"这一发什么时候发生"的第一臂。

**判决（兜底规则＝行自己的词汇，⛔ 一律同一个值）**：

| 行在插入时的状态 | 已经发生的事件 | 兜底写哪一列 |
|---|---|---|
| `Decision != ""` | 已决定（`book` 的行） | `decided_at = now` |
| `Outcome != ""` | 已结束（`book` 的行） | `ended_at = now` |
| 两者都 `""` | 尚未决定、尚未结束 ⇒ **插入就是开始**（`startCall` 的行） | `started_at = now` |
| 两者都 `""` 但调用方已给 `ended_at` | 事后补记的行：插入时刻会**晚于**它自己的结束 | **⛔ 不写**（宁可 NULL 也不自相矛盾） |

★**`startCall` 那发"该不该变"＝该变**：它插入时 DECISION/OUTCOME 都还没有 ⇒ 插入时刻就是它**唯一真拥有**的瞬间；不写它，`COALESCE` 仍读到 NULL ⇒ 这类行照样永生，`retention.go:150-153` 注释承诺的 "no NULL combination may make a row immortal" 就**仍是假话**。★**`book` 那发"不该有 `started_at`"**：`Bridge` 只看到裁决与结局（调用可能几分钟前就开始了），给它编一个开始时刻就是让列说谎 ⇒ 它 `started_at` 保持 NULL（本件 §4 读数）。
★规则本身**不允许**任何插入形状三列全 NULL：pending 拿到 `started_at`、decided 拿到 `decided_at`、finished 拿到 `ended_at`、事后补记的行自带 `ended_at`。

## §3 改了什么（逐行）

| 文件 | 行 | 改动 |
|---|---|---|
| `internal/memory/dao_toolcall.go` | `:16-17` | `InsertToolCall` 的 doc 一句改为"并填上行自己证明已经发生过的时间戳" |
| 同上 | **`:22`** | 新增一跳 `stampToolCallInsert(&tc, observe.NowWallUTC().Unix())`（★**反向判据就是摘掉这一行**） |
| 同上 | `:46-82` | 规则的长注释：为什么在 DAO、覆盖哪两个写者、三列各按什么状态填、调用方给的值永不覆盖 |
| 同上 | `:83-100` | `func stampToolCallInsert(tc *ToolCall, now int64)`（上表四行判据，**18 行**） |
| 同上 | `:102-103` | `func unixPtr(v int64) *int64`（可空时间戳装箱，**2 行**） |
| `internal/memory/migration_fkcheck_366_test.go` | `:139`、`:151-161` | ★**披露**：票 366 的 "retention" 夹具原来靠 `InsertToolCall` 造成"三列全 NULL 的子行"；现在 DAO **再也造不出这种行**，故改成**手写 pre-368 那条 INSERT**（同一列、同一字节形状），**断言与期望数一个没动**（`task_log_deleted=1`／`toolcall_deleted=1`、孤儿 1、子行 4）——该形状对**从旧版本升上来的库**仍然真实，迁移审计必须继续看得见它 |
| `internal/memory/dao_toolcall_368_test.go` | 新件 **321 行** | 7 枚用例（§8 逐名） |

`git diff --numstat`：`dao_toolcall.go` **+62/−1**、`migration_fkcheck_366_test.go` **+18/−6**、新件 **+321**。★**未**动 `retention.go`／`schema.go`／`open.go`／`privacy.go`／`internal/tools/**`。

## §4 `bridge.book` 那一发（★真 `Bridge` ＋ 真 `Store` 当 journal，一发 L0 `fs.read`；仓外探针 `probe_368_test.go`，`tree-before`=HEAD、`tree-after`=HEAD＋本腿 `memory` 字节）

夹具：父 `task_log` 在"现在"，探针把 retention 的 `Now` 钉到 **+31 天** ⇒ 子行"现在插入"即被判 31 天老。

| 读数 | 改前（`tree-before`＝`75c2395c3` 的 blob，本腿未改它） | **改后** |
|---|---|---|
| 行三列（`decision=allow outcome=success` 不变） | `decided_at=<NULL> started_at=<NULL> ended_at=<NULL>` | ★**`decided_at=1792812005`／`started_at=<NULL>`／`ended_at=1792812005`** |
| retention #1 | `tasklog_deleted=1`／**`toolcall_deleted=0`**、`fkc` **0→1** | ★**`1`／`1`**、`tool_call=0`、**`fkc=0`** |
| retention #2 | `0`／`0`、`fkc` **仍 1**、孤儿行**一字未动** | **`0`／`0`**、`tool_call=0`、**`fkc=0`**（没有第二轮可删的） |
| 手动把该行 `ended_at` 写成 40 天前再跑 | `toolcall_deleted=1`、`fkc` **0**（⇒ 不是尺坏，是 NULL） | 无幸存行可手动补（★控制自动失效，探针逐字记 `no-survivor`） |

★两列**同一秒**不是"一律填同一个值"：`book` 只有**一个**观察时刻（裁决与结局它一起写），三列里 `started_at` 仍 **NULL** ⇒ `COALESCE` 的臂序语义不动。

## §5 `startCall` 那一发（`agent/journal.go:75-82` 的逐字字面量，⛔ 走 loop）

| 读数 | 改前 | **改后** |
|---|---|---|
| 行三列 | `decided_at=<NULL> started_at=<NULL> ended_at=<NULL>` | ★**`decided_at=<NULL>`／`started_at=1792812005`／`ended_at=<NULL>`**（**不是**三列全 NULL 了，★只多这一列） |
| retention #1／#2 | `1/0`、`fkc 1` ／ `0/0`、`fkc 仍 1` | ★**`1/1`、`fkc=0`** ／ `0/0`、`fkc=0` |

★**该变，且只该变 `started_at`**（理由见 §2 末）：`decision`／`outcome` 仍 pending ⇒ 给它们时间戳就是假账；★loop 正常收尾时 `decide`／`finish` 后写的时间戳 **≥ 插入时刻** ⇒ 行内顺序 `started_at ≤ decided_at ≤ ended_at` 成立（用例 §7-⑤ 钉住）。

## §6 正控：**已有值的行不被覆盖** ＋ 「正常老化仍工作」

| 夹具 | 改前 | 改后 |
|---|---|---|
| **P3 三列预置**（调用方给 `decided/started/ended`，比 now 老 400/401/402 天） | 回读**逐字相等**（`1758252006/1758165606/1758079206`）；retention `1/1`、`fkc 0` | ★**逐字相等、一个秒都没被改写**（同上三值）；retention `1/1`、`fkc 0` |
| **P4 正常老化**（只给老 `ended_at`，无 decision/outcome，父 31 天） | `started_at=<NULL>`；retention **`1/1`**、`fkc 0` | ★`started_at=<NULL>`（§2 末行判据：已给 `ended_at` ⇒ 不编开始）；retention **`1/1`**、`fkc 0` |
| **P2 loop 全生命周期**（insert → decide → finish） | `started_at=<NULL>`、`decided=ended`；retention `1/1` | 三列**都有值**（同一秒）；retention `1/1`、`fkc 0` |

⇒ ★**尺没坏**：三列原本有值的行照旧与父同清（普查 §5 的正控在本件 P3 复现）、且**兜底 ⛔ 覆盖**它们；`fkc` 在所有"该清"的夹具上归 **0**。

## §7 新用例逐名 ＋「为什么今天响」（`internal/memory/dao_toolcall_368_test.go`）

| # | 用例（行） | 钉什么 | 今天响吗 |
|---|---|---|---|
| ① | `TestTicket368BookedRowCarriesTheInstantsItProves`（`:81`） | `book` 形状 ⇒ `decided_at`/`ended_at` 非 NULL 且**同一秒**、`started_at` **必须 NULL** | ★**响**：HEAD 上两列都是 NULL（摘掉 `:22` 那一跳也响） |
| ② | `TestTicket368BookedRowIsNotImmortal`（`:117`） | **两轮** retention ⇒ `1/1`→`0/0`、`tool_call=0`、`fkc=0` | ★**响**：HEAD 上第一轮 `tool_call_deleted=0`、`tool_call=1`、`fkc=1` |
| ③ | `TestTicket368PendingRowAgesFromItsInsert`（`:154`） | pending 行拿到 `started_at`、**不带** `decided_at`/`ended_at`，与父同清 | ★**响**：HEAD 上 `started_at` 是 NULL |
| ④ | `TestTicket368ExplicitTimestampsAreNotOverwritten`（`:191`） | **反向牙**：三列预置 ＋ pending 预置 `started_at` ＋ 只给 `ended_at` 那发 ⇒ 一律逐字保留／`started_at` 必须 NULL | 绿（控制）；★M2／M3 变异下**响** |
| ⑤ | `TestTicket368ColumnOrderSurvivesTheWrites`（`:257`） | 全生命周期后三列都有值且 `started ≤ decided ≤ ended`，并按 `ended_at` 老化 | ★**响**：HEAD 上 `started_at` 是 NULL |
| ⑥ | `TestTicket368NormalAgingStillWorks`（`:294`） | 尺的正控：只带老 `ended_at` 的子行与父同清（raw SQL 播，⛔ 走 DAO） | 绿（两侧都该绿） |
| ⑦ | `TestTicket368WindowMatchesTheContract`（`:317`） | 文件里的 30 天常量 == `TaskLogTTL` | 绿 |

## §8 反向判据：三读数 ＋ 还原核（★时机在**最后一发变异之后**）

三枚变异都是 `-overlay` 换 `dao_toolcall.go`（工作树零改动）：**M1**＝摘掉新加那一跳（`:22`）；**M2**＝兜底改成无条件覆盖；**M3**＝三列一律填（票面禁的那形）。

| 读数 | 结果 |
|---|---|
| ★**变异红** | **M1** `rc=1`：4 枚转红（①②③⑤），**逐名红行** `dao_toolcall_368_test.go:94, 98, 132, 136, 139, 148, 169, 276`。**M2** `rc=1`：1 枚（④），红行 `:210, :216, :230, :249`。**M3** `rc=1`：2 枚（①④），红行 `:106, :249` |
| ★**还原绿** | 变异后未动任何文件，`rc=0`、**7/7 PASS** |
| ★**原样绿** | 变异前同一条命令 `rc=0`、**7/7 PASS** |
| ★**还原核** | 最后一发变异（M3）**之后**量：`git hash-object internal/memory/dao_toolcall.go` ＝ **`5b4e0e027c944a78eb43378a2ebbdb9200f105dc`**（与变异前同一），且"逐名判定"归一化 `diff` **rc=0**（两次只差耗时数字） |

★判据的**方向**：M1 红＝**新加那一跳是载荷**；M2 红＝**"不覆盖"也是载荷**（不是摆设）；M3 红＝**"⛔ 一律填"这条边界真有牙**。

## §9 四数（★改前／改后分开，⛔ 相减）

| 发 | 命令 | rc | memory | tools | 墙钟 |
|---|---|---|---|---|---|
| **改前**（HEAD 字节，10:57） | `go test -count=1 ./internal/memory/ ./internal/tools/` | **0** | `ok … 46.396s` | `ok … 39.585s` | 54.1s |
| **改后**（终版字节，11:2x） | 同上 | **0** | `ok … 36.475s` | `ok … 30.791s` | — |

★两发都是**一条命令两个包一次跑**；★耗时**不比较**（机器负载在变，`368-a1` 同包 63.8s／`366-a1` 20.5s 就是同一现象）。

## §10 门禁 rc（终版字节）

| 门 | 命令 | rc | 输出 |
|---|---|---|---|
| G1 build | `GOFLAGS= go build ./...` | **0** | **0 字节** |
| G2 vet | `GOFLAGS= go vet ./internal/memory/ ./internal/tools/ ./internal/agent/` | **0** | **0 字节** |
| G3 四数 | 见 §9 | **0** | `ok` ×2 |
| G4 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations` |
| G5 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN` |
| G6 格式·**甲（入库 blob，判据）** | `git cat-file blob :<path> > /d/tmp/…/blob-N.go`（**索引 blob**＝提交后 `HEAD:<path>` 同一对象，复核见 §11）；`gofumpt -l` ／ `-d` **三枚 blob 副本** | **0** | **空**：`-l` 无行、`-d` **0 字节**；三枚 staged blob `5b4e0e02`／`c62d7b41`／`6100af22` **＝** worktree `git hash-object` 同值 |
| G7 格式·**乙（工作树，只作诊断）** | `gofumpt -l internal/memory/dao_toolcall.go internal/memory/dao_toolcall_368_test.go internal/memory/migration_fkcheck_366_test.go` | **0** | **空** |

★尺的射程（落笔前五问）：①**甲尺量的是 blob**（提交后的对象），乙尺量工作树；②`gofumpt -l` 的路径在 Windows 上是反斜杠 ⇒ 判归属一律 `[\\/]`；③本件**三条绝对路径逐字点名**、⛔ `gofumpt -l .`（会扫进 `.scratch/` 里故意写坏的探针 ⇒ `rc=2`，⛔ 当格式读数）；④"0 行输出"＋`rc=0` 两把尺同号；⑤行号/耗时全是本腿现跑，⛔ 抄普查。
★另：三枚 `.go` 全部 **LF、无 BOM**（`file` ＝ `ASCII text`、`grep -c $'\r'` ＝ **0**）；`env` 里 `TMPDIR/TEMP/TMP` 未动 ⇒ `os.TempDir()` 非空（本包 7 枚用例都在 `sealableTempDir124` 下建库，空 `TMP` 会假红）。

## §11 纪律格 ＋ 残余

- **索引核对（开工）**：`git diff --cached --name-status` ＝ **一条** `A .scratch/wisp/issues/114-…-done.md`（★**别人的** `.scratch` 暂存，**不在**本腿写面；**上报、未碰**）。★它后来被编排者自己那一笔 `29c9d2b20` 清掉；★本腿 `add`/`commit` **两处显式 pathspec、逐字点名**、⛔ `-A` ⇒ **提交后复核** `git show --name-status c0eea973b` 只有本件三枚 `.go` ＋ 本件，那条暂存项**始终没进**本腿的提交。
- **改前读数的归属**：开工 anchor `75c2395c3` 与提交父 `29c9d2b20` 之间 `git diff --stat -- internal/ cmd/`＝**空** ⇒ §4–§9 的"改前"读数与被提交的父代码是**同一份**（`29c9d2b20` 只动台账）。
- **提交后复核（甲尺的对象）**：`HEAD:<path>` 三枚 blob ＝ `5b4e0e02`／`c62d7b41`／`6100af22`，与工作树 `git hash-object` **同值**；把三枚 `HEAD:` blob 导出成副本跑 `gofumpt -l`（空、`rc=0`）／`-d`（**0 字节**）⇒ §10 的 G6 与 G7 是**同一份字节**。
- **工作树**：开工/收尾 `git status --porcelain -- internal/memory/ docs/evidence/s1/` 只有本腿三枚 `.go`；★工作树里别人的脏（`.scratch/**` 的 `M`／`?? $null`／`-`／`-dummy` 等）**一个没碰、一个没清**。
- **临时件**：只建不删，全在 `D:\tmp\wisp\20261024\368r1\**`；★`.out` 一枚没有（日志一律 `.txt`／`.log`）。
- ★**残余（⛔ 本票的活，逐条具名）**：
  1. ★**retention 仍零生产调用者**（本腿现跑：`StartRetentionJob`／`RunRetentionOnce` 非测试命中只有 `retention.go` 自身）⇒ 今天**一次都不跑**；本腿治的是"**一旦接线，孤儿不再积累**"，⛔ "今天已经在清"。
  2. ★**pre-368 库里的存量三列全 NULL 行**：本腿**不回溯修**（那是数据迁移，不是写者修复；票 366 的审计**照旧记录**它们）。要不要补一个 `UPDATE … WHERE started_at IS NULL AND ended_at IS NULL AND decided_at IS NULL` 的版本化步骤 ⇒ **另开票／另拍**。
  3. ★`bridge.book` 仍把 id 丢掉（`bridge.go:1225`）；现在那行自带终局时间戳 ⇒ 不需要更新者，但将来若要"改一行"就得留 id（`internal/tools` 不在本腿写面，⛔ 动）。
