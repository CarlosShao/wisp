# 票 326 裁决腿 `326-v1` 读数件（非实现者·只读·⛔ 翻框）

- 起手 HEAD＝`45c4c03e`（12:07:49 +08）｜本腿唯一新产码面字节＝`c2e1c07a`（一枚新测试件，写面 4）｜交件时 HEAD 已随别腿前进（现量见末节）
- 票框尺（12:07 现跑，射程＝票 326 文件）：`grep -c '^- \[ \]'`＝**6**（rc=0）⇄ `grep -c '^- \[x\]'`＝**0**（rc=1）⇒ **本腿交回仍是 6／0，一枚没翻**
- 台面：两棵 `git archive` **仓外**导出树 ＋ 原生反斜杠 `GOTMPDIR='C:\Users\swq\AppData\Local\Temp\wisp326v1'`（先 `mkdir -p`）
  - `/d/tmp/w326v1/a`＝`git archive 45c4c03e` 七件套（`go.mod go.sum internal cmd tools scripts frontend`）＝**745 枚 `.go`**，交付态
  - `/d/tmp/w326v1/prefix`＝`git archive 86298e24^` 七件套＝**742 枚 `.go`**，未修码（`dao_session.go` sha256 `38bff760eabb5bf9…`，`cmp` 对 `fd01e796` 同名件 **rc=0**＝实现方那棵未修树同一个 blob）
  - 突变全部打在**盘上字节**（⛔ `-overlay`、⛔ 共享工作树），驱动＝`mutate326v1.py`，每发落 `sha_before/sha_after` 且**回读盘上字节**复核
- 本节以下每个数都带【尺＋射程＋量的/算的】。原始逐发＝本目录 `logs/`（10/11/20/30/31/40/41/42/43/50/60/70/80/90 十四枚）

## 1. 实现方五笔写面（逐枚 `git show --name-only`，⛔ 区间尺）

| 笔 | 名册（枚数） | 是否在派单五面内 |
|---|---|---|
| `86298e24` | `internal/memory/dao_session.go`／`dao_session_append_326_test.go`／`dao_session_test.go`／`internal/session/manage_test.go`（4） | 是（1／2 号面） |
| `d3f98938` | `.scratch/wisp/probes/326/r1/**`（21） | 是（3 号面） |
| `a60da735` | `docs/evidence/s1/326-append-updated.md`＋`probes/326/r1/msg-c3.txt`（2） | 是（1／3 号面） |
| `f2787393` | 票 326 文件＋`msg-c4.txt`（2） | 是（4／3 号面） |
| `b4adfd0e` | `probes/326/r1/70-roster-audit.txt`＋`msg-c5.txt`（2） | 是（3 号面） |

零枚 `cmd/wisp/**`、零枚 `internal/panel/**`（尺＝名册逐枚，rc=0）。

三把边界尺（都现跑）：
- `git show --numstat --format='' 86298e24`＝`28/0`｜`181/0`｜`120/0`｜`60/0` ⇒ 四件**删除行数全为 0**（算的）⇒ 这一笔⛔ 动过任何既有句子；`git diff 86298e24^ 86298e24 | grep -cE '^-.*\.(Errorf|Fatalf)\('`＝**0 枚 rc=1**
- `git diff --numstat 86298e24 HEAD -- internal/memory internal/session`＝**空**（rc=0）⇒ 交付字节＝我今天读的字节（当时 HEAD `45c4c03e`）
- `git diff --numstat fd01e796 ce8df8b6 -- internal/memory internal/session`＝**空**（rc=0）⇒ 实现方"成对读数只差那一枚改动"的同屏控制**成立**
- 禁区尺＝`git diff --stat 86298e24^ HEAD -- 三枚冻结件 internal/observe/thresholds.go docs/PLAN.md go.mod go.sum`＝**空 rc=0**（零字节差）
- 新增 `t.Skip` 枚数＝**2**（`^+.*t\.Skip` 命中，rc=0），两句都在**新建件**里（父／子角色互斥，形状＝既有 `TestSubprocessCrashWriter`），⛔ 盖住任何判据格

## 2. AC#0 起手复算（票面四行逐枚现跑，⛔ 用票上的号）

| 票面那行 | 我的尺（射程） | 我的读数 | rc |
|---|---|---|---|
| 追加只写子表 | `sed -n '273,307p' dao_session.go \| grep -c 'UPDATE session SET updated_at'`（交付态函数体）⇄ 同尺打 `86298e24^` 同段 | **1**（已翻面＝本票修的就是它）⇄ **0** | 0 ⇄ 1 |
| 真正写 `updated_at` 的处所 | `grep -n 'updated_at' internal/memory/dao_session.go`（整文件） | 10 行＝列名册 1＋建 1＋两枚 `ORDER BY` 2＋改名 `:177`＋归档 `:198`＋分叉 `:231`＋**新增 `:300`**＋注释 2 ⇒ 写侧**五处**（原四处＋说话） | 0 |
| 列表按这一列排 | `grep -n 'ORDER BY updated_at DESC, id DESC'`（同文件） | **2 枚＝`:125`（列表）与 `:144`（搜索）** ⇒ 与编排者 12:1x 读数**一致** | 0 |
| 时间源一枚族 | `grep -n 'observe.NowWallUTC'`＝`:175/:192/:218`；`orNow`＝`dao_tasklog.go:14-19`（零值那支 `observe.NowWallUTC().Unix()`）；`manage.go:77` 默认 `time.Now` | 秒精度、同一族；事务闭包体（`:279-:302`，起止行现抽）里读时钟的句子 **0 枚 rc=1**（尺＝那 24 行 `grep -nE 'NowWallUTC\|time\.Now\|orNow'`）；整函数体只 1 枚＝`:274 t.CreatedAt = orNow(t.CreatedAt)`，位置在闭包**外** | 0 |

## 3. AC#1 刷新真发生在写入那一枚事务里（三发全跑）

**① 锚（内容＋行号）**：`internal/memory/dao_session.go:300`
＝`_, err = tx.ExecContext(ctx, ` + "`UPDATE session SET updated_at=? WHERE id=?`" + `, t.CreatedAt, t.SessionID)`；
所属事务句柄＝**同文件 `:279` 那枚闭包的 `tx *sql.Tx` 形参**（`err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {`），
该闭包由 `open.go:280 s.write → writer.go:67 enqueue → writer.go:135 q.db.BeginTx` 拉起，`writer.go:155-157` 句子函数报错即 `return err`、`writer.go:142 defer tx.Rollback()` 整枚回滚。
全文件那枚 `UPDATE session SET updated_at` **只有 1 枚**（grep -c＝1，正控＝对 pristine 副本同尺 1 枚）。值＝`t.CreatedAt`（`:274 orNow` 归一）。

**② 去牙（MU-A，⛔ 只报号）**：物理删掉 `:300` 那一行 ⇒ `sha256(ee48c311d1925a2b…)`＝实现方 M01 同值（我独立跑出同一个号）。**必须红的用例＝6 枚**（5 枚在 `internal/memory`、1 枚在 `internal/session`），红句逐字（全文＝`logs/10-muA-de-tooth.txt`）：

```
dao_session_test.go:118: appending a turn did not reach session.updated_at: sess-a.updated_at = 1777593660, want the appended turn's stamp 1777593780 (before=1777593660, delta=0)
dao_session_test.go:122: continuing a chat does not surface it: list = [sess-b sess-a], want sess-a first after one more turn in it
dao_session_test.go:157: refresh used a different clock than the turn it wrote: updated_at = 1775001600, want turn.created_at 1792124215
dao_session_test.go:161: updated_at did not move forward: 1775001600 -> 1775001600 (appending is a touch)
manage_test.go:222: AppendTurn wrote the child row only: sess_ceed….updated_at stayed 1777593660, want it past 1777593660
dao_session_append_326_test.go:86: half state after kill at 2 commits: session_turn has 3 rows … (want NOTHING or ALL, got a 3602s lag or a torn seq)
dao_session_326v1_test.go:63: AppendSessionTurn succeeded while the session touch was forced to fail — the refresh error is being swallowed (ticket 326 landing criterion ③)
```
还原：`restored sha256=ec6ee530b190795093829eca07e6003bf866db5da3e64c47e6d0b0b1bcd1d2b2 cmp_pristine=MATCH`、`cmp` 对 `45c4c03e` blob **rc=0**、六枚复跑全 `--- PASS`（`logs/11-restore-after-muA.txt`）。

**票面那一发"未修码上必红"我也自己复跑了**（⛔ 抄实现方的号）：未修树 `/d/tmp/w326v1/prefix`＝`86298e24^`＋交付的三枚测试件 ⇒
红名册＝`TestAppendSessionTurnSurfacesTheConversation`／`TestAppendSessionTurnRefreshUsesStoredClockFamily`／`TestAppendSessionTurnCrashAtomicAcrossKill`（8/8 枚杀点 `HALF-STATE`，`lag=3603s..3649s`）／`TestSessionAppendTurnMovesItBackToTheFront`，
四数＝`internal/memory RUN=78 PASS_TOP=38 FAIL_TOP=3 SKIP=2 rc=1`、`internal/session RUN=37 PASS_TOP=19 FAIL_TOP=1 SKIP=0 rc=1`（**与实现方申报的号逐字相同**，`logs/90-*`）。
★一条差异具名：我的去牙发与真未修码**⛔ 逐字节相同**——`diff` 显示去牙形还多带 `86298e24` 补的那枚 `if err != nil { return err }`（`LastInsertId` 的错误检查）＋注释（尺＝剥注释后函数体 `cmp` rc=1，差在 line 21）。枚枚到"红句"这一层两把尺咬合，形状不等价这一点不许糊过去。

**③ ★反形（MU-B，编排者最想裁死的一格）**：把刷新挪出事务＝提交 `INSERT` 之后另起一枚 `s.write`（突变全文＝`logs/20-*` 首屏）⇒ `sha256(c31e1750690993a76270b86891da3499550d692b3de3ce0d52432aa9aa13b1b5)`。实测两族颜色：

| 那族尺 | 反形下的颜色 | 结论 |
|---|---|---|
| 三枚 AC#1 终态用例（`TestAppendSessionTurnSurfacesTheConversation`／`…RefreshUsesStoredClockFamily`／`manage_test.go TestSessionAppendTurnMovesItBackToTheFront`） | **全 PASS**，且同行证据仍带 `updated_at 1777593660 -> 1777593780 (delta=120)` | ⛔ 抓得住＝实现方自报的盲区，我复现成立 |
| 真杀台件 `TestAppendSessionTurnCrashAtomicAcrossKill` | **逐发 6/6 枚 FAIL**；逐枚杀点 `HALF-STATE` **23/48＝48%**（4 轮 14/32 ＋作废那两发 5/8＋4/8） | 抓得住但⛔ 确定（命中随时序漂，⛔ 当判据的终态色） |
| 本腿新建确定性仪器（`dao_session_326v1_test.go`，枚一枚 SQLite `BEFORE UPDATE … RAISE(ABORT)` 把母表刷新强制打失败） | **2/2 枚 FAIL**，红句：`the touch and the turn do not share one transaction: the refresh failed (… constraint failed: 326v1 injected refresh failure (1811)) yet 1 session_turn row(s) stayed committed while session.updated_at = 1777593600` | ★**这一形现在有确定性仪器了**（0.14s、⛔ 杀进程、⛔ 依赖墙钟） |

交付态这两枚仪器**绿**（`append err=memory: append session turn: constraint failed: … | session_turn rows=0 | updated_at=1777593600`＝失败即整枚回滚）；未修码上**红**（`logs/90-*` 末节）。⇒ 〔缺世界〕那一格**不再缺**；派单说的 (α) 接受 56%／(β) 另立一票开测试接缝，两枚都⛔ 需要——接缝我用现成的（真实 SQLite 触发器），⛔ 动 `Store.write`／DAO 函数体形状。

## 4. AC#2 排序两行逐字节不动（复认＋⛔ 替编排者裁）

- 尺1＝`git diff 86298e24^ 86298e24 -- <四件>` 对 `^[+-].*ORDER BY updated_at DESC, id DESC` ⇒ **命中 0 枚 rc=1**；diff 本体 **442 行**（⇒ 这把尺⛔ 空转）
- 正控＝同一把尺打在仓外副本的假改动（删 `, id DESC`）⇒ `grep -n ORDER BY` 两侧对照差 **4 行**（2 删 2 增）rc=0 ⇒ "0 命中"⛔ 是尺坏
- 逐字节＝`sed -n '125p;144p'` 两行切片的 sha256：锚 `fd01e796` 与今天**同值** `6a4e53226589a6fb7005e5f8647c9eea75fcfc5fb0d43c2f710ed0ac2274b3ef`；行号今天仍是 `:125`／`:144`（量在交付树与工作树各一次，同值）
- 尺2＝`git diff 86298e24^ 86298e24 | grep -E '^[+-].*(rand\.Read|func Mint|id DESC)'` ⇒ 命中 **2 行，两句都是 `+//` 开头的注释**（`id DESC` 出现在注释里）；`rand.Read`／`func Mint` **0 命中**，`internal/session/session.go:96` 的铸 id 句子今天照旧（正控＝该词在库里能命中）⇒ **同秒先后那一半病本票没顺手治**
- 上交给我裁的料（我只复认存在、⛔ 裁要不要治）：`326-r1` 的暴露面读数＝同秒 200 对 ⇒ 硬币 90/200＝45%；跨秒 10 对 ⇒ 10/10 对。**这两发我⛔ 复跑**（那要跑它的 bench 件 `zz326r1-tiebreak_bench_test.go`，属"要不要治"的输入，⛔ 属"修得对不对"的判据）⇒ 具名＝**未复跑**，需要时缺的那一发＝在交付树里 `cp` 该 bench 件后 `go test -bench` 两发

## 5. AC#3 崩溃原子（我自己跑的杀台件＋"是不是恒过"那一问）

- 形状核：两条语句同处 `:279` 那一枚闭包 ⇒ 同一枚 `tx`；`writer.go:142/155-160`＝报错即回滚、成功才 `Commit`
- 真杀台件我自己跑了 **5 发**（父 `cmd.Process.Kill()`，逐枚杀点报色）：交付态 **40 枚杀点（5 发 × 8 枚）＝`HALF-STATE` 0 枚**、逐发 `--- PASS`（5.68s／1.34s／3 发计数 0,0,0／1.57s）
- 这台件**⛔ 恒过**（三形实测都拦得住）：
  - 打在**去牙**（MU-A）⇒ **8/8 枚杀点 `HALF-STATE`**，`lag=3602s..3649s`（红句逐字见 `logs/10-*`）
  - 打在**反形**（MU-B）⇒ 逐发 **6/6 枚 FAIL**、点命中 48%
  - 打在**未修码**（真 prefix 树）⇒ **8/8 枚 `HALF-STATE`**、`lag=3603s..3649s`
- ★派单那一问（"把刷新挪回事务内但故意在两条语句之间 `os.Exit`，它必须报 `HALF-STATE`"）我照做了（MU-D，sha `21dd5d710b6b2093…`，门控 `WISP_326V1_EXIT_MID_TX=1` ⇒ 同一枚二进制可开可关）：
  - 门控**关**（正控）＝`ok … 2.043s`（与交付态同色 ⇒ 那枚 `if` 本身⛔ 改行为）
  - 门控**开**＝**没报 `HALF-STATE`**，红在这一句：`dao_session_append_326_test.go:86: child announced 0 commits within 60s, want >= 2 (台件故障，不是判据)` ⇒ `--- FAIL (60.03s)`
  - 我的判读（⛔ 顺着派单话术写）：**派单预期的颜色在这形拿不到，也不该拿到**——未提交事务被 `os.Exit` 掐死正是"要么全无"的**全无**那一半，`HALF-STATE` 在这形是物理不可达；台件用**自己的前提控制**（announce 计数）拒绝假绿，⇒ 这台件**有牙、且⛔ 恒过**。要说它⛔ 抓得住的，是"整枚事务⛔ 提交"这一形，而那本来就是合法终态
- 迁移那族⛔ 破：`git diff 86298e24^ 86298e24 -- internal/memory/dao_session.go | grep -cE '^[+-].*(applyV3|ddlV3|migrationChain|SchemaVersionTarget)'`＝**0 枚 rc=1**；同发兄弟台件三枚绿＝`TestSessionMigrationCrashAtomicAndReplay (0.17s)`／`TestCrashRecoveryKillMidWrite (0.19s)`／本票杀件 `(1.57s)`，`TestSubprocessCrashWriter`＝`--- SKIP`（父子角色互斥，既有形状）

## 6. AC#4 门禁＋两包四数（每把带 rc；整包⛔ 跑）

尺口径：`RUN=`＝`grep -c '^=== RUN'`（**含子测试**）⇄ `PASS_TOP=/FAIL_TOP=/SKIP=`＝`grep -c '^--- PASS|FAIL|SKIP'`（**只数顶层**）；两把⛔ 混用。**跑包一律显式 `-timeout 30m`**（派单点名默认 10 分钟对本仓不够；本腿⛔ 一发 `panic: test timed out`，也⛔ 跑 `./cmd/wisp` 与 `go test ./...`）。

| 把 | 台面／射程 | rc | 读数 |
|---|---|---|---|
| `go build ./...` | 导出树 a（745 枚 `.go`） | **0** | 干净 |
| `go vet ./internal/memory/ ./internal/session/` | 导出树 a（含本腿新件） | **0** | 零行输出 |
| `gofumpt -l internal/memory internal/session` | 导出树 a，**真身** `/d/work/base/gopath/bin/gofumpt.exe` | **0** | 列出 **0 枚**；对照＝裸名 `gofumpt` → `command not found`（⛔ 在 PATH）；⛔ 拿 `gofmt -l` 冒充 |
| `sh scripts/d22scan.sh` | **仓根工作树**（含别腿在飞字节＝台面读数） | **0** | `clean - no D22 ban violations`；分母 bans#1-5 `internal/=240`、`cmd/=41`｜#6 `frontend/=116`｜#7 `internal/tools/=25`｜#8 `design/=39`、`frontend/=116`、`internal/=577`、`cmd/=139`（比实现方那次各涨若干＝别人加了件） |
| 同尺 | 导出树 a | 1 | **台面故障、⛔ 计入判据**：红的是 `runtests.sh` 正向控制里三枚**读运行库**的用例 `TestScannerSelfScanOfRealRepoIsGreen`／`TestScopeReportMatchesRealCoverage`／`TestRealRepoLedgerIsHonest`（树里 `ls -a \| grep -c '^\.git$'`＝**0**＝无 git 树）⇒ 具名＝这三枚那一行红，⛔ 代码红 |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 仓根 | **0** | `VERDICT GREEN`；分母 **14264** 枚跟踪路径、越帽 **60**＝名册 **60**、名册外 **0**、最长相对 185、runner 最坏 229（帽 165） |
| 同尺 | 导出树 a | **2** | 尺自述拒答：`not inside a git work tree - refusing to guess a denominator`＝**GUARD 射程拒绝，⛔ 读成红** |
| `go test ./internal/memory/ -count=1 -v` | 导出树 a，**未加**本腿件（＝交付字节原样） | **0** | `RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP=2`（17.401s）⇔ 与实现方申报**逐字同值** |
| 同尺 | 导出树 a，**加了**本腿件（HEAD `c2e1c07a` 态） | **0** | `RUN=80 PASS_TOP=43 FAIL_TOP=0 SKIP=2`（净增＝本腿 2 枚顶层用例） |
| `go test ./internal/session/ -count=1 -v` | 导出树 a | **0** | `RUN=37 PASS_TOP=20 FAIL_TOP=0 SKIP=0`（1.933s／2.269s 两发同数） |
| SKIP 名册 | 两发各一次 | — | `TestSubprocessCrashWriter`＋`TestAppendSessionTurnCrashChild`（父／子互斥；⛔ 新加 skip 逃判据，见第 1 节尺2） |
| 时钟那一格 | `:300` 写的值＝`t.CreatedAt`（`:274 orNow` 归一，位置在闭包**外**）；闭包体 `:279-:302` 内读时钟的句子 **0 枚 rc=1**（正控＝同尺打 `:170-200` 命中 2 枚）；`dao_session.go` 亦⛔ `import "time"`（`grep -c '"time"'`＝0 rc=1） | 1 | ⛔ 事务内再读时钟＝**成立** |

## 7. AC#5 沿到消费者（说实话；两把编排者读数我都复认）

- 复认编排者尺①：`git grep -nE 'ListSessions|SessionStore' HEAD -- cmd internal/panel internal/session` 滤 `_test` ⇒ **2 行**＝`internal/session/manage.go:40`（接口声明）＋`:104`（唯一转调 store 的句子）⇒ **`cmd/` 与 `internal/panel/` 里零调用者**（我这一把⛔ 滤也 2 行，该射程内无测试件命中）
- 复认编排者尺②：库里那两行排序今天仍在 `dao_session.go:125`／`:144`
- 票面三把尺逐字复跑（各带正控）：`git grep -c 'session\.NewManager' -- '*.go'`＝**0 rc=1**（正控 `func NewManager`＝`manage.go:68` 1 枚）；`grep -cE '"sessions"|"history"' cmd/wisp/main.go`＝**0 rc=1**（正控＝同文件 `case "…"` 名册 `run/providers/doctor/secret/models/slo`，这把尺在该文件命中得到）；`git grep -n 'memory\.Session' -- cmd/ internal/panel/ internal/tools/ internal/agent/`＝**0 rc=1**（正控 `internal/session/manage.go` 18 枚）
- ★比票面更深一跳（我数出来的，具名）：`AppendSessionTurn` 的**生产调用点＝1 枚**（`manage.go:150`，⛔ 拿 `go build` rc=0 冒充接线）；它的上游 `Manager.AppendTurn`（`manage.go:149`）的**生产调用点＝0 枚**（尺＝`git grep -n 'AppendTurn' HEAD -- '*.go'` 滤 `_test` ⇒ 只剩 `manage.go:148/149` 两句自身）；`session.NewManager` 生产调用点 **0 枚** ⇒ 会话列表今天没接进产品，**连"追加这一轮"在生产码里也还没人调**（今天会写 `session_turn` 的句子只有 `:237` 分叉拷贝与 `:287` 追加两处 DAO 内部，全树生产字节里无第三个写入点）
- `cmd/` 里 `session.*` 的标识符名册（算的，逐枚）：`SessionID`8／`Ledger`7／`Record`4／`NewLedger`3／`Mint`2／`ID`1／`LedgerOptions`1／`Covering`1 ＋ `session.txt` 字符串 1 ⇒ **全是审批台账那一族**（`run.go:282/476/481`、`resident_approval_windows.go:243/285`）；会话 `Manager` 那一族 **0 枚** ⇒ "包被 import 了"⛔ 等于"列表接上了"，这句必须这么说
- ⇒ 判语：值对了、界面上今天**仍看不见**，缺的那一跳＝装配根构造 `session.NewManager` ＋一枚读者（子命令或面板快照），落 `cmd/wisp/**`＋`internal/panel/**`＝本波别的腿的面，**属票 166 的接线缺口**；本票修的是 DAO 那一格，⛔ 因为界面看不见就判它修错

## 8. 本腿写面自证（逐笔）＋交件三把尺

- `c2e1c07a` 的 `git show --name-only`＝**只含 1 枚** `internal/memory/dao_session_326v1_test.go`（写面 4 允许的新建件；⛔ 改任何既有测试件、⛔ 动产码）
- 其余两枚（判语件／票尾 bullet）见各自提交；本读数件＝第三笔
- 我读过的**既有**件逐枚"只读没动"：`internal/memory/dao_session.go`／`dao_session_test.go`／`dao_session_append_326_test.go`／`schema_test.go`／`tempdir_resolved_124_test.go`／`dao_tasklog.go`／`open.go`／`writer.go`／`internal/session/manage.go`／`manage_test.go`／`session.go`／票 326 文件／`docs/evidence/s1/326-append-updated.md`／`scripts/d22scan.sh`／`scripts/check-path-length-budget.sh` ⇒ 尺＝`git diff --numstat 86298e24 HEAD -- internal/memory internal/session` 只列 `125 0 internal/memory/dao_session_326v1_test.go`（**⛔ 本腿新件之外一枚都不差**）
- ⛔ 放宽任何既有断言（尺＝`86298e24` 四件删除行全 0 ＋ 我名下删除行 0）
- 工作树 `go.mod` **仍脏**（mtime `Oct 10 22:28`、`git diff --stat` ＝ 1 insertion/1 deletion、另带 `LF will be replaced by CRLF` 警告）＝存量、归属不明 ⇒ 本腿⛔ 碰／⛔ 提交／⛔ revert，**具名报"还脏"归编排者处置**
- 台面缺件补一枚具名：七件套在 `45c4c03e` 上 `go build ./...` **rc=0**（本票射程⛔ 需要 `models/` 与 `build/*.dll`；`220-v1` 报的 `models/` 缺件只影响 `internal/models` 那族用例，本腿⛔ 碰到）
