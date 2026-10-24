# 票 368 / 腿 `368-v1` 终裁（非实现者）— 〔丁〕成立：形对、牙真、**「不编 `book.started_at`」这条承重判断判对**；`AC#2` 可翻（按〔丁〕射程），存量与「父删即子删」留给级联那一票

- 腿＝`368-v1`（**终裁腿**，非实现者）。被裁＝`368-r1`：`c0eea973b`（码＋件）／`92ceaecdf`（件）／`7fc90d21b`（件）。
- 我的两棵树：`tree-before`＝`git archive 75c2395c3`（r1 anchor）· `tree-after`＝`git archive 7fc90d21b`（r1 交付）；`git diff --stat 75c2395c3 c0eea973b^ -- internal cmd scripts`＝**空** ⇒ 两棵树的「改前」是同一份字节。★判的三枚 blob 在本腿观察到的 **5 个 HEAD**（`7fc90d21b`／`255d42c75`／`b4ee24cb`／`9965e86a9`／`16b1d3f3`）上**逐枚同值**（`5b4e0e02`／`c62d7b41`／`6100af22`）⇒ 下面读数对今天的 HEAD 成立（本腿期间别人推了 4 笔，HEAD 漂了 4 次，⛔ 我未跟）。
- 写面：本件＋仓外 `D:\tmp\wisp\20261024\368v1\**`。**仓内 `.go`＝0 枚、0 字节**（判据：①本腿开工/收尾 `git status --porcelain -- internal/ cmd/ scripts/` 均**空**；②本腿只对仓内一个路径执行过写命令＝本件；③探针三枚在 `368v1\probes\**`，只进 `git archive` 导出树）。★工作树里别人的脏（`.scratch/**` 的 `M`／`??`）**一个没碰、一个没清**。
- ★尺的射程（五问）：①"行三列"＝**新连接** `SELECT decided_at/started_at/ended_at`（`sql.NullInt64`）⇒ `<NULL>` 是 **SQL NULL**；②`fkc`＝`PRAGMA foreign_key_check`（FK 关也答，`368-a1` §1 已证）；③retention＝**生产的** `RunRetentionOnce`，`Now` 钉到 **+31 天**（插入即 >30 天）；④"四数"是 go test 自报的两个包耗时，**分开报、⛔ 相减**；⑤hash＝`git hash-object`（SHA-1、40 位）。

## §1 ① 票面前提新鲜度（对 HEAD blob 现跑）＋`368-a1` 三条

| 票面/普查的说法 | HEAD 现跑 | 判 |
|---|---|---|
| `journal.go:75-82` 的字面量不设三列 | `:75-82` **逐字未变**（只 `TaskID/Seq/Tool/ArgsJSON/RiskLevel/CorrelationID`） | ✅ 字面量 premise 真 |
| 票面「`startCall` 插的行三列**全 NULL**」 | **已过期**：DAO 兜底后同一条行＝`decided=<NULL> started_at=1792813647 ended=<NULL>`（before＝三列全 NULL） | ⚠ 票面这句**只对改前字节成立** |
| `retention.go:154-155` 的 `COALESCE` | blob `2846718e`，与 anchor **同值**；末句逐字 `DELETE … WHERE COALESCE(…) IS NOT NULL AND COALESCE(…) < ?`；`age==ttl` 保留 | ✅ 未动 |
| `368-a1` §1.1「1/1、留 1 孤儿、fkc 0→1」 | 我的 raw 三-NULL 子行（§6）：retention#1 **1/0**、left **1**、fkc **1**；#2 **0/0**、fkc **1** ⇒ **永生复现**。★「1/1」的 1 是旁边老时间戳子行（普查自己已改口），我这一发没放它 ⇒ 1/0 | ✅ 事实对，口径按普查的修正 |
| `368-a1` §1.2 `StartRetentionJob` 零生产调用者 | 非测试命中只有 `retention.go` 自身（`:95/:98/:100/:104/:108/:233-235`） | ✅ 仍真 |
| `368-a1` §2 附「`started_at` 生产零写入者」 | 改前真（`UPDATE tool_call` 只有 `:113` 写 decided_at、`:142` 写 ended_at）；**改后 INSERT 自己写 started_at** | ✅ 改前真 ⚠ **改后该句失效**（`COALESCE` 第二臂不再死） |
| `368-a1` §7 两个 `InsertToolCall` 调用点 ＋ `run.go` 装配 | `journal.go:75`／`bridge.go:1225` 仍是两个；`run.go:777 Journal: mem`／`1118 Journal: nil`（普查记 `:769`／`:1080` ⇒ **+8／+38**） | ✅ 事实真、行号漂 |
| `loop.go` 三个 `startCall`、`schema.go:83`、`open.go:413` | **`:671`／`:910`／`:942`**、`schema.go:83`（`REFERENCES task_log(id)`，**无 ON DELETE**）、`open.go:413` | ✅ 与 r1 报的漂移一致 |
| 新漂移（r1 自造） | `dao_toolcall.go` 的 INSERT 列名从 `:24-31` → **`:26-33`**（+2，r1 加的 doc＋跳） | ⚠ 普查的这枚锚点已过期 |

## §2 ② 核心判点一：`bridge.book` 那一发（★我自己造，真 `Bridge`＋真 `Store`，一发 L0 `fs.read`）

| 读数 | 改前（`tree-before`） | **改后（`tree-after`）** |
|---|---|---|
| 行（`decision=allow outcome=success`） | `decided_at=<NULL> started_at=<NULL> ended_at=<NULL>` | ★`decided_at=1792813579`／`started_at=<NULL>`／`ended_at=1792813579` |
| retention #1 | `tasklog=1`／**`toolcall=0`**、rows **1**、fkc **`[tool_call|rowid=1|parent=task_log|fkid=0]`** | ★**`1`／`1`**、rows **0**、fkc **`[]`** |
| retention #2 | **`0`／`0`**、rows **1**、fkc **仍 1** | **`0`／`0`**、rows **0**、fkc **`[]`** |
| 正控（把幸存行 `ended_at` 手写 40 天前） | **`0`／`1`**、rows **0**、fkc **`[]`** ⇒ 尺没坏 | `no-survivor`（无幸存行可补） |

⇒ ★两列**同一秒**不是「一律同一个值」：`book` 只有一个观察时刻（裁决与结局一起写），三列里 `started_at` 仍 NULL ⇒ `COALESCE` 臂序不动。★`fkc 0→1`（改前）与 `0`（改后）都是我自己量的。

## §3 ③④ 两条判断的判语（本腿承重）

- ★**③ `startCall`「该变、且只该变 `started_at`」＝判对**。①插入发生在 gate/执行**之前**（`loop.go:671`），插入时刻是写者**唯一观测到**的瞬间；②`decision`/`outcome` 的 `""` 在 `models.go:81/85` 明写＝**pending** ⇒ 给它们时间戳就是给「没发生的事」记时刻；③严格尺（§4）量到三列是**三个不同秒** `1792814490 < 1792814491 < 1792814492` ⇒ 不是把一个值写三遍；④`privacy.go:103-108` 优先读 `started_at` ⇒ pending 行不再显示 1970；⑤retention：pending 行 `1/1`、fkc 0（改前 `1/0`、两轮删不掉）。★射程：`started_at` 的语义＝「这一发进流水线的时刻」，审批等待算在 span 内（主流 trace 就是这么切）；要「工具体真正开始」得另立写者/列，⛔ 本票的活。
- ★**④ `book.started_at` 该留 NULL：「编一个就是让列说谎」＝判对（全腿最承重的一步）**。`book` 只观测到裁决与结局，**没观测到开始**；给它一个「插入时刻减若干」的猜测，是没有观测支撑的**倒填**。★注意这条的不对称，它才是关键：`decided_at`/`ended_at` 记的是**行自己已经断言发生的事**（`Decision`/`Outcome` 非空）在**写者能担保的最早时刻**（=晚界，保守）；`started_at` 则会给一件**行没有任何列断言**的事编一个**早界**，且会污染时长语义。更硬的一条：retention 只需要**一列**非空，写 `decided_at` 就够破 NULL ⇒ 填 `started_at` **买不到任何好处**、只买到一句假话。★消费者也不需要它：`privacy.go` 落到 `decided_at`、`ListToolCalls` 的 `COALESCE(started_at, decided_at, 0)` 两处都有兜底（现跑读过）。

## §4 ⑤ 四分支逐条 ＋ ⑥ 永不覆盖／顺序

| 分支（我一发一发造，真 DAO） | 改前 | **改后** | 判 |
|---|---|---|---|
| 两者皆空（插入即开始） | 三列全 NULL | ★`started_at=1792813647`，另两列 NULL | ✅ 该变且只该变它 |
| 两者皆空＋调用方已给 `started_at` | 保留 | 保留（不入第二值） | ✅ |
| **两者皆空＋调用方已给 `ended_at`（事后补记）** | 只有 `ended_at=1758080847` | ★**`started_at` 仍 NULL**、`ended_at` 保留 | ✅ **这条我专门造了：不写** |
| `Decision` 有值 | 三列 NULL | `decided_at=1792813647`，另两列 NULL | ✅ |
| `Outcome` 有值 | 三列 NULL | `ended_at=1792813647`，另两列 NULL | ✅ |
| 两者都有（book 形） | 三列 NULL | `decided_at=ended_at`，`started_at` NULL | ✅ |
| 两者皆空＋已给 `decided_at`（边界） | decided 保留、另两列 NULL | decided 保留＋`started_at=now` | ✅ 不覆盖；⚠ 此刻 `started_at > decided_at`（调用方输入自相矛盾，两个生产写者都不产这种形；规则只承诺「不再三列全 NULL」，⛔ 不承诺校验输入顺序） |
| 列名有值＋已给 `ended_at`（边界） | decided NULL | `decided_at=now > ended_at`（同上，输入矛盾） | ✅ 非缺陷，写进射程 |

- ★**全称量尺**：把 DAO 输入空间枚举成 `decision{空,allow} × outcome{空,success} × 时间戳 8 种子集` ＝ **32 形**，逐形真插入＋raw 回读 ⇒ **改后 32/32 至少一列非空、0 形三列全 NULL**；改前 **28/32**，全 NULL 的 4 形恰是「没给任何时间戳」的四种基础形。★「调用方给的值永不覆盖」：32 形里**逐列逐字**比对，**0 处被改写**。
- ★**⑥ 正控**：预置 `decided/started/ended`＝400/401/402 天前 ⇒ 回读 `1758253647/1758167247/1758080847`，**一个秒都没动**（改前改后同值）。
- ★**⑦ 顺序**：DAO 全生命周期（insert→decide→finish）三列 `1792813647`（同一秒，因调用是瞬时的）；**严格发**（每次隔 1.1s）＝`1792814490 ≤ 1792814491 ≤ 1792814492`；★**真 loop**（golden `tool-then-text`，一发 `echo`，真 store 当 Journal）＝`started=decided=ended=1792813593`，三列全有值、retention `1/1`、fkc 0；**改前**同一发 `started_at=<NULL>`（2/3）⇒ 这是今天响的牙。

## §5 ⑧ 三变异（全 `-overlay` 换 `dao_toolcall.go`；工作树零写入）

| 变异（我自己的字节） | rc | 转红用例（逐名） | 逐名红行 |
|---|---|---|---|
| **M1** 摘掉新加那一跳（`:22`） | **1** | ①②③⑤ `BookedRowCarriesTheInstantsItProves`／`BookedRowIsNotImmortal`／`PendingRowAgesFromItsInsert`／`ColumnOrderSurvivesTheWrites` | `:94 :98 :132 :136 :139 :148 :169 :276`（**8 行**） |
| **M2** 无条件覆盖（删掉每处 `== nil` 守卫） | **1** | ④ `ExplicitTimestampsAreNotOverwritten` | `:210 :216 :230 :249`（**4 行**） |
| **M3** 三列一律填（每个 NULL 列都塞同一个 `now`，⛔ 按行词汇） | **1** | ①③④ `BookedRowCarries…`／`PendingRowAges…`／`ExplicitTimestamps…` | `:106 :174 :249`（**3 行**） |

⇒ ★M1 红＝新跳是载荷；M2 红＝「不覆盖」是载荷；M3 红＝「⛔ 一律填」边界有牙。★**M3 与 r1 报的 2 红/2 行不同**：r1 的 M3 更窄（只动 booked 支），我的 M3 是**票面禁的那句的完整形**（每列都填），因此**多红 ③**；③ 本来就钉「pending 行不许有 decided/ended」，红得对。
★**原样绿**：变异前同一条命令 `rc=0`、**7/7 PASS**；★**还原绿**：最后一发变异（M3）之后 **`rc=0`、7/7 PASS**；★**还原核**（时机＝M3 之后量）：`git hash-object` 三枚＝`5b4e0e02/c62d7b41/6100af22`，仓内文件 vs `tree-after` **逐名 `cmp` rc=0 ×3**，`git diff --numstat HEAD --` 三枚 **空**，`git status --porcelain -- internal/memory/` **空** ⇒ **无工作树突变窗口**。

## §6 ⑨ 它改票 366 那枚夹具（手写 pre-368 INSERT）——判「它还测不测原来那件事」

- ★**diff 只有写者换了**：同列同值（`task-old/seq=1/read_file/{}/L0/NULL/NULL/NULL/NULL/NULL/NULL/c-old-null/NULL`），断言与期望数**一个没动**（`task_log` 1／`tool_call` 4／孤儿 1／子行 `task-A,task-A,task-B,task-old`）。
- ★★**决定性读数**：把 **pre-r1 的夹具文件**用 `-overlay` 换回去、跑同一条用例，再查那一行的形状 ⇒ `task-old/seq=1` 的 `started_at=1792814025`（**不再是 legacy 形**）；r1 的版本 ⇒ `decided_at/started_at/ended_at` **三列全 NULL**（我打印的四行见日志）。★且**两版都 rc=0**：旧夹具**不会变红**，它只是**悄悄不再测那件事** ⇒ 这次改动**不是「为了变绿」，是为了保住被测对象**。
- ★**判：改法对，且比原写法更真**（raw INSERT 不再受将来 DAO 变化影响；「旧库升上来的行」这个对象对迁移审计仍然真实）。★代价（该披露）：夹具不再证明「生产写者造得出该形状」——而这恰是 r1 要推翻的命题，⛔ 不是损失。★366 全套 4/4 子用例绿（新夹具）。

## §7 ⑩ 四数／格式／边界／门禁（逐门自己跑）

- ★**四数**（`go test -count=1 ./internal/memory/ ./internal/tools/`，**两包分开报、⛔ 相减**）：**改前** memory **27.133s**／tools **24.540s**（rc=0，pre-r1 字节，经 3 文件 `-overlay` 在**仓路径**里跑）；**改后** memory **26.587s**／tools **24.559s**（rc=0）。★与 r1 的 46.396/39.585→36.475/30.791 **不比较**（机器负载；本腿 `368-a1` 同包 63.8s 就是同一现象）。★射程：先在导出树跑「改前」时 `TestTicket252R1Alignment…` 因 `D:\tmp` 无 8.3 别名**拒绝 skip 而红**（环境前置，⛔ 代码红）⇒ 改在仓路径＋overlay 重测，避免一次假红。
- **G1** `GOFLAGS= go build ./...` **rc=0**、**0 字节**；**G2** `go vet ./internal/memory/ ./internal/tools/ ./internal/agent/` **rc=0**、0 字节；**G4** `sh scripts/d22scan.sh` **rc=0**：`clean - no D22 ban violations`（internal/ 641 Go 文件、cmd/ 150 逐枚过）；**G5** `sh scripts/check-path-length-budget.sh` **rc=0**：`VERDICT GREEN`、17688 tracked、over-budget 61/roster 61/not-in-roster 0、两单位同选 61。
- **格式甲（入库 blob＝判据）**：`git cat-file blob HEAD:` 三枚导出（`git hash-object` 回读＝`5b4e0e02/c62d7b41/6100af22`）⇒ `gofumpt -l` **空、rc=0**；`gofumpt -d` **0 字节、rc=0**。
- **格式乙（工作树＝只作诊断）**：`gofumpt -l internal/memory/dao_toolcall.go internal/memory/dao_toolcall_368_test.go internal/memory/migration_fkcheck_366_test.go` ⇒ **空、rc=0**；归属判定用 `[\\/]`（命中 `internal[\\/]memory` 的行＝0）。
- **边界零 diff**（r1 交付 `c0eea973b^..c0eea973b`）：`retention.go`／`schema.go`／`open.go`／`internal/tools/**`／`cmd/**`／`scripts/**`／`.github/**`／`docs/PLAN.md`／`docs/specs/**`／`docs/reports/**` **逐目录 `git diff --stat` 全空**；`retention.go/schema.go/open.go/models.go/bridge.go/journal.go/loop.go/run.go/`两脚本在 `75c2395c3`→HEAD **blob 同值**（`SAME` ×10）。★我自己的仓内 `.go`＝**0**。

## §8 ⑪ 编排者那两条裁 ＋ 提交号／逐笔越界

- ★**（a）⛔ 做 backfill ＝判对**。四条理由（都可复跑）：①**没有可依据的值**——行里没有"创建时刻"列，任何 `UPDATE` 只能拿父 `started_at`／迁移时刻去**编**，正是它拒给 `book.started_at` 编值的同一条理由；②若拿迁移时刻，legacy 行会**当场越过 30 天**、下一轮 retention 就删——"修"变成**静默删取证**；③366 的迁移审计**就是要持续看见这个形**（§6 已量），backfill 会把它从用户库里抹掉，等于拆掉审计的物证；④它**治不了别的孤儿源**（见下 (b) 的 coexist 读数）。★补一句：存量**已经**成孤儿的行，将来的级联也够不到 ⇒ 它们需要一次**删**（孤儿清扫），⛔ 不是一次**编时间的 UPDATE**——"不 backfill"不等于"存量不管"。
- ★**（b）父删时级联删子 ＝判对，而且是 r1 缺的那半**。★现跑：父删只有三处——`DeleteTaskLog`（`dao_tasklog.go:118`）、`PurgeTaskLogs`（`:134`）、retention（`retention.go:146`→`:197`）；FK **无 `ON DELETE`**（`schema.go:83`）且 DDL 冻结/本票禁改 ⇒ **必须在同一事务里做应用层级联**，⛔ 加 `foreign_key_check` 硬门槛（票面硬提醒 1 仍成立）。★**coexist 读数**（我造）：父 31 天、子 `started_at` 2 天前／`ended_at` 1 天前 ⇒ **改后三列全有值**（`started=1792641870`／`decided=1792814670`／`ended=1792728270`）、retention#1 `tasklog=1 toolcall=0`、fkc **1**；`Now+30 天` 再跑 `1/0/[]` ⇒ **一个可长达 30 天的孤儿窗口，且三列没有一个是 NULL** ⇒ `started_at` 兜底**治不了它**（这不是 NULL 问题）⇒ 级联不是 backfill 的"替代品"，是**另一半**。
- ★**该披露的三条代价**（⛔ 当反对理由）：①`PurgePrivacy(PrivacyTaskLog)` 会连带删 `tool_call`，而 `PrivacyToolCall` 是**另一个隐私域**（`privacy.go:152-159/170-179`）——产品语义要点头；②366 的 `purge-tasklogs`／`delete-tasklog` 两枚夹具会**不再造孤儿** ⇒ 得像 r1 这次一样改 raw 播种（同一条"legacy 形只能手写"的逻辑）；③若只做"每处父删时删子"，**插悬空行那条**（票面第五条路径）仍不覆盖；写成 `DELETE FROM tool_call WHERE task_id NOT IN (SELECT id FROM task_log)` 的一次**孤儿清扫**才连它一起收（⛔ 只删、⛔ 编时间）。
- ★**它列的残余三条逐条核：全对**——retention 零生产调用者（现跑真）、存量三-NULL 行不回溯修（我 §6 两轮 `1/0` 永生死不了）、`bridge.book` 仍丢 id（`bridge.go:1225` `if _, err :=` 逐字）。
- **提交号／逐笔越界**：`c0eea973b`＝`A docs/evidence/s1/368-r1-impl.md`＋`M internal/memory/dao_toolcall.go`＋`A internal/memory/dao_toolcall_368_test.go`＋`M internal/memory/migration_fkcheck_366_test.go`（**全在 `internal/memory/**`＋件**，越界 0）；`92ceaecdf`／`7fc90d21b` **只动件**。★本腿首笔＝**`66b6a8b93`**（`git show --name-status` 只有 `A docs/evidence/s1/368-v1-verdict.md` 一枚；提交前 `git diff --cached --name-status` 同一枚、无 `D`／无别人的 `A`；两处逐字点名单文件、⛔ `-A`）；本件哈希由紧随的一笔补录（⛔ amend）。

## §9 交件② `AC#2` 能翻／⛔翻 ＋ 缺什么

- ★**能翻（按〔丁〕的射程）**。判据四条：①**两种生产写者的新行都不再永生**（真 `Bridge`／真 loop／32 形枚举，`1/1`＋fkc 0）；②**票面那一发**（31 天父＋NULL 子）在改后 `1/1`、fkc 0，且改前 `1/0`、两轮删不掉；③**牙是真的**（M1/M2/M3 全红，红行逐名；7 枚用例里 ①②③⑤ 今天响）；④门禁全绿、还原核在最后一发变异之后量。
- ★★**但票面 `AC#2` 第二句「父被清时它跟着清」不是全称**：coexist 发（父 31 天、子 1 天）父清了、子**不清**，fkc 1（§8）。⇒ 翻框时**必须写清射程**（"新写入的行不再永生；父删⇒子删只有在子自己的 `COALESCE` 也过线时才成立"），或者把这半句**交给级联那一票**。若把 AC 原句当全称读，则**⛔ 不能翻**——这是本腿唯一一处需要编排者落字的判断。
- ★**翻框还缺的（不阻塞〔丁〕，但必须留在残余）**：①**存量**三-NULL 行仍永生（升级库；r1 已披露）；②`bridge.book` 仍丢 id（将来要"改一行"才需要）；③retention 仍零生产调用者（今天一次不跑）；④孤儿清扫/级联尚不存在（编排者已裁方向，未开票）。
