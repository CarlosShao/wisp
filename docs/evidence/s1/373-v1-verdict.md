# 票 373 / 终裁腿 `373-v1`（非实现者）— 裁决：**〔甲〕同删＋先父后子合格**，★但「孤儿类已全关」⛔随票翻

- 腿＝`373-v1`（**非实现者裁腿**）。⛔ 落任何 `.go`、⛔ 翻框／改票名／`Status:`／`-done`／台账、⛔ push。写面＝本件＋仓外 `D:\tmp\wisp\20261024\373v1\**`（before 导出树／probe 2 枚／mut 4 枚／overlay 6 枚／日志 20 枚）。
- 开工 HEAD `00c7468dd`；落笔 HEAD `5d64b118c`（期间 377-a1 两笔＋`A1232` 只动 `docs/**`、`.scratch/**`；`git diff --numstat 00c7468dd..HEAD -- internal/memory cmd/wisp` ＝**空**）。被裁＝`93fdc0496`＋`d03c500da`＋`00c7468dd`（`git show --name-status` 逐笔见 §12）。
- ★★**工作树声明（停手上报项）**：开工 13:00 我现跑 `git status --porcelain -- internal/memory cmd/wisp` ＝**空**；**13:07:19 起 `cmd/wisp/run.go` 被别的腿改脏**（`M`、`5+/1-`、`consoleApprovalUI.Prompt` 的「窗口」行，L2 改读 Deadline），另有 `?? cmd/wisp/approval_card_window_376_test.go`（13:08:50）⇒ **376 腿在飞**。★我⛔碰、⛔清；后果是 373-r1 的「还原核 5/5 逐位相同」我只证 **4/5**（§8）。`internal/memory` 4 枚 `.go` 全 `SAME`。
- ★**before 口径（我改的，具名）**：票面让我用 `git archive HEAD` 当 before，但 HEAD 已含 373-r1 ⇒ **before＝`git archive 93fdc0496^`（`28b15313d`）**，after＝现树；否则 before≡after，一个数也量不出。probe 全在仓外，经 `go test -overlay` 注入（`cmd/wisp` 那枚只注入**测试文件**，⛔ 改被测树）。
- **尺的射程**（落笔前五问）：①对象＝`Open(dir)` 建出的**真 v4 库**（⛔手抄 DDL）；`fkc`＝业务形状连接（`dsn()`，`foreign_keys` 关）上 `PRAGMA foreign_key_check` 的**行数**。②形状＝整数＋`child|rowid|parent|fkid`＋`tool_call` 的 `(task_id,seq)` 多重集。③代码 vs 注释：契约命中用现跑 `git grep`；`d22scan` 数**盘上**文件（故 `cmd/` 150→151 是别人那枚未跟踪测试，不是本笔）。④总数核对：每 fixture 先读 `fkc=0` 基线；E1 另有**确定性**判据（触发器+审计表），⛔只靠 2ms 采样。⑤谁算的：整数由 SQLite `COUNT(*)` 给；只比 `sha256sum`（64 位）与 `git show HEAD:<f>` 的字节。

## §1 十一项逐项判语（读数在 §3–§11）

| # | 项 | 判语 |
|---|---|---|
| 1 | **票面前提新鲜度＋`373-a1` 三修正** | ★**前提在 before 树为真、对 HEAD 已陈**：before 实测三处留孤儿 3/2/2/1（§3）；HEAD 行号全漂（`:118`→`129`、`:134`→`163`、retention 父删 `:146`→`:163`、`:197` 泛型→`:218`）。★**三修正全部成立**：①retention 是**两条独立语句**（父 `started_at`／子 `COALESCE(ended_at,started_at,decided_at)`、同 tx `:147`）⇒「只删时间戳非 NULL 的子」是**半句**；②`CASCADE` 是**死条款**（P8 实测 `on_delete=NO ACTION`；`dsnPragmas` 四枚无 `foreign_keys`；DDL 里 `ON DELETE` 命中 **0**，唯一的 1 处是 `retention.go:228` 的注释）；③`schema.go:353 DROP TABLE task_log` **非孤儿源**（同一步/tx、`:348-352` 逐列搬、`:476` 重建前先查名册、`migrationDSN` 强 FK OFF）⇒ 我复核**成立** |
| 2 | **核心一：三处级联** | ★**合格**（§3：before→after 四发 3/2/2/1 → 0/0/0/0，且「留了谁」逐次点名） |
| 3 | **核心二：先父后子 E1** | ★**合格**（§4：改前 `max_fkc=1`、触发器记到 child→parent **反序 44ms**；改后 `0`／0 条） |
| 4 | **核心三：事务边界** | ★**合格**（§5：`RAISE(ABORT)` 拒父删 ⇒ 子删回滚 3 行仍在；摘触发器 ⇒ 1 行。retention 同法 `1→0`） |
| 5 | **「删的是谁的」** | ★**合格**（§3 表：P1 删 3 行全量／P2 只删 A 的两子**留 `task-B#1`**／retention 只删到期父的子**留 `task-live#1`**；P5b 证 not-found 时子删**回滚**） |
| 6 | **`366` 夹具改 raw 播种** | ★**改法对**（§6：仍测「legacy 孤儿库被 3→4 迁移记录且不删行」3/2/1/1；不再测「今天的生产路径会造孤儿」，后者已由 §3 的 0 从正面钉住） |
| 7 | **件一（`Journal` 保持 `nil`／`:1334` 才补开）** | ★**对，且是唯一不越权的选法**（§7：子环路写面无字变动=数据模型不被动；`rt.loopOpt` 在 `:1197` 先定型、`:1209-1213` 只在**副本**上装 ⇒ 断言成立）。★代价就是 §9-1 那第二类孤儿源**仍在** |
| 8 | **反向判据 M1–M4＋还原核** | ★**成立**（§8：每发只红对应那一枚，`fkc` 回 2/3/2/1、E1 回 `max_fkc=1`；还原核在最后一发变异**之后**量） |
| 9 | **四数＋格式甲乙** | ★memory **73/0/2 → 78/0/2**（rc=0）、`cmd/wisp` 定向 **10/0/0 → 10/0/0**（rc=0）；`gofmt -l` **空**、`gofumpt -l` **空**（§11） |
| 10 | **边界** | ★**0 越界**（§10：6 文件；`ddlV1–V4`／`schema.go`／`open.go`／`dsnPragmas`／`winsec`／`panel`／`perm` 全 **0 diff**；`cmd/wisp` 只 `run.go`） |
| 11 | **三条上报该不该在本票做** | ★**1 ⇒ 另立 378（⛔本票；票面 `378` 已在 `.scratch` 立案）／2 记／3 记**，三条判语全同意（§9） |

## §2 `AC#2` 能翻／⛔翻 ＋ 缺什么

- ★★★**按 `AC#1` 裁形的范围（三处级联＋件一件二件三）＝能翻**：四发读数、E1、事务边界、M1–M4、还原核、门禁都过（§3–§11）。
- ★★**缺 A（⛔随票翻的一句）**：票面「孤儿这个类**只有一个总根**」「本票**一次解决**这个类的所有来源」「票 366 的审计会**长期读 0**」——**量出来是假的**：E2（`task.spawn` 真入口）⇒ 子任务 `fs.read` 行 `parent_rows=0`，跑完 `task_log=1／tool_call=2／fkc=1`（§4）。⇒ **该句必须改口径或挂 378**，⛔把本票当「类已关」。
- ★**缺 B（契约文字）**：`PurgePrivacy(PrivacyTaskLog)` 连带删 `tool_call` 已落地，但 `docs/specs/**` 一字未动（写面外）⇒ 要另派/摆机主。
- ★**缺 C（一件没在本笔钉）**：`cascade_373_test.go` 未单独钉「三列全 NULL 的 legacy 子」那一形（我用仓外 probe P4 量到 `1/0/fkc=1 → 1/1/fkc=0`；行覆盖在 368 的用例里）⇒ 不是洞，属**可加不可加**的边。
- ★**缺 D（纪律项，⛔实现者的）**：`run.go` 现脏使 5/5 还原核只能证 4/5（§8）。

## §3 ★★三处级联：我自己造的（before→after，含「留了谁」）

夹具：`task-A`(2 子)/`task-B`(1 子) 走生产 DAO；retention 另造 `task-old`(31d)/`task-live` 对控。probe＝`probe/v373_probe_test.go`（仓外，`-overlay`），尺＝`COUNT(*)`＋`fkc`＋多重集。

| 那一发（生产入口） | before `task_log`/`tool_call`/`fkc` | after | 删掉的（多重集）／**留下的** |
|---|---|---|---|
| P1 `PurgePrivacy(PrivacyTaskLog)`（`rows=2` 两树同） | 0／**3**／**3** | 0／**0**／**0** | `[A#1 A#2 B#1]` 全量／**`[]`** |
| P2 `DeletePrivacyItem(PrivacyTaskLog,"task-A")` | 1／**3**／**2** | 1／**1**／**0** | `[A#1 A#2]`／**`[task-B#1]`** |
| P3 retention（父 31d、**两子 5d**、活对控） | 1／3／**2**，`1/0` | 1／1／**0**，`1/2` | `[old#1 old#2]`／**`[task-live#1]`** |
| P3b 同上但**单子** 5d（票面那个 1/1 算术） | 1／2／**1**，`1/0` | 1／1／**0**，`1/1` | `[old#1]`／`[task-live#1]` |
| P4 retention legacy（三列全 NULL 子） | 1／2／**1**，`1/0`；第二发仍 `0/0` | 1／1／**0**，`1/1`；第二发 `0/0` | `[null#1]`／`[task-live#1]` |

★判语：**删对了、且没删多**——P2 只动 A、P3/P3b/P4 都不碰活父的活子；★**`P3` 是票面没写的第二机制**（父子窗口不齐）我复现了；★票面的「1/1」是**单子夹具的算术**，两子就是 `1/2`（⛔当常数背）。★P3/P4 第二发 `0/0` ⇒ 不是把孤儿推迟，是没了。

## §4 ★★★先父后子 E1：我自己造的（尺比它硬）＋「时间戳不能当判据」判语

- 尺：`runTextTask` 真入口＋一次模型发出的 `fs.read`；**三条**独立取证：①SQLite 触发器 `BEFORE INSERT ON tool_call` 把「父不在就插子」**逐条**记进审计表（**确定性**：⛔采样漏不掉，且**不拦**那次插入，被测写面零改动）；②对称触发器记每次 `task_log` INSERT ⇒ **审计表 rowid 顺序＝两次写的事务顺序**（都过 `s.write` 单写队列）；③2ms 采样 `COUNT`＋`fkc`。
- **before**（`28b15313d` 导出树）：`poll_samples=45`、**`max_fkc_during_run=1`**、`child_without_parent_seen=true`；审计 rowid 1＝`child_without_parent`@`05:14:08.026`，rowid 2＝`parent_opened`@`.070` ⇒ **反序 44ms**；跑完 `1/1/fkc=0`（**跑完就自愈**，所以只看终点读不出病）。
- **after**（现树）：`poll_samples=29`、**`max_fkc=0`**、`child_without_parent=false`；审计只有 rowid 1＝`parent_opened`@`05:20:20.459`；跑完 `1/1/fkc=0`。
- ★★**「时间戳是秒粒度、不能当判据」＝对，而且理由比它说的更强**：①我的实测里**子行的 `started_at` 是 NULL**（`bridge.go` 是事后 `book`，只落 `decided_at`/`ended_at`）⇒ 比较连定义都没有；②before 那次两次写相隔 **44ms** 却同秒（`05:14:08`）⇒ 「同秒」确实分不出前后；③所以它那句「两值都同秒（`1792816889`／`1792816918`）」的数字**不是证据**（两个数差 29s，本身还跨了两发），**结论对、举证弱**。★判据只认「事务顺序（触发器/审计）或运行期 `fkc`」。
- ★**尺的射程**：它的 2ms 采样形状是 `children>0 && parents==0`，**看不见第二类孤儿**——E2 里存在无父的子行时该布尔仍是 `false`（§9-1）。⇒ 采样只能当辅证。

## §5 ★★★事务边界那一发：我自己造的（`RAISE(ABORT)`）

- 夹具 `A`(2 子)`/B`(1 子)；种 `CREATE TRIGGER v373_block BEFORE DELETE ON task_log BEGIN SELECT RAISE(ABORT,'v373 blocked'); END`。
- **after 树**：`DeletePrivacyItem(task-A)` ⇒ 报错 `constraint failed: v373 blocked (1811)`，库仍 `task_log=2／tool_call=**3**／fkc=0` ⇒ **子删已跑、父删被拒 ⇒ 子删回滚**；`DROP TRIGGER` 再发同调用 ⇒ `task_log=1／tool_call=**1**`。★**「3」与「1」之差就是回滚**。
- rejection 同理（`retention.go`）：被拒后库仍 `1/1`，摘触发器再跑 ⇒ `0/0`。★诚实补一刀：**被拒那一发 `res.ToolCallDeleted=1` 但库里一行没少**（计数在闭包内先赋值）——它的 §8 已披露这条，我复现了。
- 第二条（⛔靠触发器）：`DeleteTaskLog("task-ghost")`（父不存在、子存在）⇒ `ErrNotFound`、**子行仍在**（`removed=[]`）⇒ 证明同 tx，不是「顺序对就完事」。
- ★**before 树同一发是 `3→3`**（它压根没有子删）⇒ 这一发**只有在 after 树才有判别力**，这正是它设计对的地方。

## §6 `366` 夹具改 raw 播种：仍测什么／不再测什么

- **仍测**（逐条）：①`purge-tasklogs`／`delete-tasklog`／`retention` 三条**legacy 库形状**（0 父 3 子／1 父 3 子／1 枚孤儿）照样被 `3→4` 迁移**记录**成 `3/2/1/1`；②**不删行**、不 gate、不回滚；③`manual-insert` 那条第五形不变。读数我复跑过（§11 的 78/0/2 含这 4 枚）。
- **不再测**（逐条）：①「今天的 `PurgeTaskLogs`／`DeleteTaskLog`／retention 会造孤儿」——**这正是本笔修掉的**，留着它只能靠断言失败来「测」；②retention 那发里「aging 子窗自己删掉的老时间戳子」也不再经 `RunRetentionOnce` 走（改成手写孤儿子），**该行为仍在 `dao_toolcall_368_test.go`／`retention_test.go` 覆盖**。
- ★判语：**改法对**——同 `368-r1` 的先例（raw 播种 legacy 形状），且新 fixture 的**终态逐列同旧 fixture**（旧的那发只是把 `task-old` 父和老化子删掉，两者都不出现在终态）。★唯一代价：raw 播种绕过 DAO，将来若 v3 DDL 冻结被破就会漂；v3 是冻结的，可接受。

## §7 件一：`rt.loopOpt.Journal` 保持 `nil`／`:1334` 才补开

- ★**`nil` 正确**：`rt.loopOpt, rt.loopOptSet = opts, true`（`:1197`，`opts.Journal` 就是 `nil`）**先**定型；`:1209` 的 `loopOpts := opts` 是**副本**，`:1213` 只改副本 ⇒ `task.spawn` 读到的 `BaseOptions`（`run.go:817`→`internal/tools/subagent_197.go:306`）仍是 `Journal:nil`。★**「子环路写入面一字不变」成立**（E2 佐证：子任务**没有** `task_log` 行）。
- ★**`:1334` 只在 `!opened()` 才补开，正确**：三种到达方式逐条对——①控制词（`loop.go:349` 早退）不派工具⇒无子可孤；②开行写失败＝E1 形，补写后 `fkc` 回 0（与改前**每个**任务同形）；③正常路径只 `Finish`，⛔ 写第二行（journal 的 `FinishTaskLog` 惰性，`started_at` 从「写完时刻」变「开始时刻」）。`opened()` 在 `bg.Wait()` 之后读（`close(t.done)` 是屏障，且 `atomic.Bool`）。
- ★**我加的一条残余（它未报）**：开行提前后，**进程崩在 start 与 close 之间会留 `state='running'`**（改前崩了**一行都没有**）。这不是缺陷（主流 harness 正是先落行），但本仓**没有开机对账**把 `running` 收成 `interrupted`（`interrupted` 只由 `StatusTruncated` 映射写入）⇒ 建议随 378 或另立小格，⛔本笔。

## §8 ★★反向判据 M1–M4（我自己造）＋还原核

| 变异（仓外副本＋`-overlay`，工作树零突变） | 摘掉什么 | 我的读数 |
|---|---|---|
| M1 | `DeleteTaskLog` 的 `DELETE FROM tool_call WHERE task_id=?` | **只 P2 红**：`fkc=2`、留 `[A#1 A#2 B#1]`；P1/P3/P3b/P4 全绿 |
| M2 | `PurgeTaskLogs` 的 `DELETE FROM tool_call` | **只 P1 红**：`fkc=3`、`tool_call=3`；其余全绿 |
| M3 | retention 的 cascade 调用（真摘，赋 `int64(0), error(nil)`） | **只 retention 红**：P3 `fkc=2`、P3b `fkc=1`、P4 `fkc=1`；P1/P2 绿 |
| M4 | `loopOpts.Journal = rootRows` | **只 E1 红**：`max_fkc=1`、`child_without_parent=1`（审计 `child`@`05:16:30.010`→`parent`@`.023`） |

- ★**还原核（时机＝最后一发变异 M4 @13:16:30 之后）**：`sha256sum` vs `git show HEAD:<f>` ⇒ **4/5 SAME**（`dao_tasklog.go 9a22f389…`／`retention.go a771c0f9…`／`migration_fkcheck_366_test.go 86e19845…`／`cascade_373_test.go 4f28956b…`）；★**`run.go` 判 DIFF，原因见件头：别人（376）的 `5+/1-` 未提交改**，`git diff HEAD -- cmd/wisp/run.go` **只有那一个 hunk**，373 的 hunk 逐字未动。`git status --porcelain -- internal/memory` ＝**空**；`go test ./internal/memory` **rc=0（78/0/2）**；E1 probe **PASS**（`max_fkc=0`）。
- ★**13:3x 收口复核（376-r1 以 `ceef4d069` 落库之后）**：`cmd/wisp/run.go` 现 `git status` **空**、与 HEAD `sha256` **SAME（`565ff346…`）**、373 的三处 hunk（`taskLifecycleJournal`／`loopOpts.Journal = rootRows`／`rootRows.opened()`）**仍在** ⇒ **还原核 5/5 逐位相同成立**（HEAD 的 `run.go` 另含 376 那一笔，与本笔无关）。

## §9 ★★三条上报：该不该在本票做

1. ★★★**子环路不写 `task_log` 行 ⇒ `task.spawn` 的子任务 `tool_call` 无父**。★**我自造 E2 实证**：审计 rowid 1＝`parent_opened`（root）、rowid 2＝`child_without_parent`（子任务）；库 `task_log=1／tool_call=2／fkc=1`；逐行 `spawn` 行 `parent_rows=1`、子 `fs.read` 行 `parent_rows=**0**`。★**判：⛔本票做**——给每个子任务新增 `task_log` 行会改**隐私页列表**与 **retention 面**（父龄按谁算、子任务行算谁的隐私域），那是**数据模型决定**不是机械补写；它「只具名上报、未动」**处置对**；**另立 378 对**（票面已立案）。★★但**票面「类已关」必须改口**（§2 缺 A）。
2. ★**`PurgeTaskLogs` 语义扩到 `tool_call` 已落地**，契约文字（`SPEC-02 §5` 一带）未写。★**判：记、⛔本票**——`docs/specs/**` 不在其写面；代码侧证据我复核过：`PurgeToolCalls` 仍在（`dao_toolcall.go:207`／`privacy.go:179`）⇒ 「只想清取证」的窄门还在，语义扩展不是单向黑洞。
3. ★**残余：开行写失败会回落到「父后到」**。★**判：记、⛔本票**——这是 store 故障路径，且补写把 `fkc` 回补到 0（与改前对**每个**任务的行为同形）；真「父从未落」只剩这一条。★比「写不进就拒绝跑任务」更合本仓的「大声降级」政策，我**不同意**把它升级成本票缺陷。

## §10 边界（逐项，`git show --name-status` 现跑）

- 三笔＝`93fdc0496`（`cmd/wisp/run.go`／`internal/memory/{dao_tasklog,retention,migration_fkcheck_366_test,cascade_373_test}.go`／本件的前身 `373-r1-impl.md`）＋`d03c500da`／`00c7468dd`（**只** `docs/evidence/s1/373-r1-impl.md`）。**总数 6 文件**。
- `git diff --numstat 93fdc0496^..93fdc0496`：`run.go 98/4`、`dao_tasklog.go 32/3`、`retention.go 40/3`、`migration_fkcheck_366_test.go 42/36`、`cascade_373_test.go 367/0`、件 128/0。
- ★**禁区别名逐条 0 diff**：`internal/memory/schema.go`（`ddlV1`–`ddlV4` 文本与语句数）、`internal/memory/open.go`（`dsnPragmas`／`migrationDSN`）、`internal/winsec/**`、`internal/panel/**`、`internal/perm/**` ⇒ `git diff --name-only 93fdc0496^..HEAD | grep -Ei "schema\.go|open\.go|winsec/|panel/|internal/perm/|ddl|dsn"` ＝**NONE**。`cmd/wisp/**` 只有 `run.go`。⛔ 任何 FK／schema／`fkc` 硬门槛改动＝**0**。

## §11 门禁 rc（全部 `cmd > out 2>&1; rc=$?`，⛔走管道）

| 门 | 读数 |
|---|---|
| `go test -count=1 -v ./internal/memory/` | before **73/0/2** rc=0 ⇒ after **78/0/2** rc=0（+5＝新用例 5 顶层） |
| `cmd/wisp` 定向（`PATH="$PWD/third_party/sherpa-onnx:$PATH"`，`-run "TestRunTextTask\|Test365\|TestHostDispatch…\|TestComposedGate…\|TestMissingBlob…"`） | before **10/0/0** rc=0 ⇒ after **10/0/0** rc=0（⛔整包，会死 `MARK299`） |
| `go build ./...` | **rc=0**（log 0 字节） |
| `go vet ./internal/memory/ ./cmd/wisp/` | **rc=0**（log 0 字节） |
| `d22scan`（`cd tools/d22scan && go run . -root ../../`） | **rc=0**，`clean - no D22 ban violations` |
| 格式甲 `gofmt -l <5 文件>` | **空**，rc=0 |
| 格式乙 `gofumpt -l <5 文件>` | **空**，rc=0（⛔ `-l .`：会扫 `.scratch` 探针 ⇒ rc=2） |
| ⭐新用例的牙：把 `cascade_373_test.go` 注入 **before 树** | **5/5 红**，红行逐字 `:109/114/118`、`:141/145/149/165`、`:198/203/208/211/222`、`:253/260`、`:328/361` ⇒ 与它 §9 的表逐行对上 |

## §12 提交＋逐笔越界

- 本件＝**1 枚文件**（`docs/evidence/s1/373-v1-verdict.md`）；`git add`／`commit` 两处**逐个文件名字面点名**、⛔glob／`-A`；提交前 `git diff --cached --name-status` 现跑＝**只有本件（`A`）**，无 `D`；**⛔ push**。★**本件提交＝`5977c17fe`（正文笔）＋本补录笔**，两笔各只本件 1 文件；落件时 HEAD `5d64b118c`。
- ★**越界 0**：本腿仓内写 `internal/**`／`cmd/**`／`scripts/**`／`.scratch/wisp/**`／`docs/specs/**`／`docs/PLAN.md`／`docs/reports/**` **各 0 字节**；仓外只有 `D:\tmp\wisp\20261024\373v1\**`。★别人的脏（`.scratch/**` 大量、`cmd/wisp/run.go` 的 376 改、`?? cmd/wisp/approval_card_window_376_test.go`、`377-a1` 在读）**⛔碰、⛔清**。
- ★**落点核**：落笔前 `ls docs/evidence/s1/373-*` ＝`373-a1-census.md`／`373-r1-impl.md` 两枚 ⇒ `373-v1-verdict.md` **空号，未改号**。
