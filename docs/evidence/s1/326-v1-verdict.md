# 票 326 非实现者终裁（裁决腿 `326-v1`）——`AppendSessionTurn` 刷新 `session.updated_at`

- 腿＝`326-v1`（裁决者·只读·非实现者）｜实现方＝`326-r1`（五笔 `86298e24`→`b4adfd0e`）｜票＝`.scratch/wisp/issues/326-append-session-turn-does-not-refresh-session-updated-so-a-live-chat-never-moves-up.md`
- 起手 HEAD＝`45c4c03e`（2026-10-16 12:07:49 +08）｜交件 HEAD＝本腿三笔 `c2e1c07a`＋`ed9c81d3`＋本笔（别腿在其间亦前进，逐笔名册见文末尺①）
- **⛔ 翻框**：本腿交回时票框尺仍是**未勾 6／已勾 0**（`grep -c '^- \[ \]'`＝6 rc=0 ⇄ `grep -c '^- \[x\]'`＝0 rc=1）；翻框归编排者
- 全部读数与尺＝`.scratch/wisp/probes/326/v1/00-readings.md`＋`logs/`（15 枚原始件）；突变只打在 `git archive` 的**仓外**导出树盘上字节（⛔ `-overlay`、⛔ 共享工作树当改前）

## 判语总表（只用 成立／不成立／不可判〔缺哪一发读数〕）

| 格 | 判语 | 决定性读数（我自己跑的） |
|---|---|---|
| AC#0 | **成立** | 票面四行逐枚现跑，行号今天＝写侧五处、排序两行 `:125`/`:144`（与编排者读数一致）；"追加只写子表"那行已翻面（同尺在 `86298e24^` 命中 0、交付态命中 1） |
| AC#1 | **成立** | 锚＝`dao_session.go:300` 那句 `UPDATE session SET updated_at=?`，事务句柄＝`:279` 闭包的 `tx *sql.Tx`（`writer.go:135 BeginTx`／`:142 defer tx.Rollback()`／`:155-157` 报错即整枚回滚）；去牙⇒**6 枚红**；真未修树⇒票面那一发逐字复现；反形⇒**本腿新增确定性仪器抓得住** |
| AC#2 | **成立** | `ORDER BY updated_at DESC, id DESC` 两行**逐字节未动**（diff 尺 0 命中 rc=1、正控 4 行；两行切片 sha256 `6a4e5322…` 与锚同值）；`rand.Read`／`func Mint` 在本笔 diff 里 0 命中⇒**同秒那一半没顺手治** |
| AC#3 | **成立** | 交付态我自己跑 **5 发×8 枚杀点＝40 枚全 `ALL`、`HALF-STATE` 0 枚**；这台件**⛔ 恒过**：去牙 8/8 红、反形逐发 6/6 FAIL、未修码 8/8 红；迁移族 diff 0 命中＋三枚兄弟台件同发绿 |
| AC#4 | **成立** | 两包四数逐字复现实现方申报（`internal/memory RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP=2 rc=0`⇄未修码 `78/38/3/2 rc=1`；`internal/session RUN=37 PASS_TOP=20 FAIL_TOP=0 SKIP=0 rc=0`⇄未修码 `37/19/1/0 rc=1`）；门禁四把各落 rc；一律显式 `-timeout 30m`；⛔ 跑 `./cmd/wisp`、⛔ 跑 `./...` |
| AC#5 | **成立** | 说实话这一格属实且我加深一跳：会话列表今天⛔ 接进产品（`cmd/`＋`internal/panel/` 里 `ListSessions`/`SessionStore` **零调用者**＝编排者尺①复认），且 `Manager.AppendTurn` 生产调用点 **0 枚**⇒"值对了"在界面上看不见，缺的那一跳属票 166 的接线缺口 |

## AC#1 三发的落点（编排者点名要我裁死的一格）

**② 去牙**＝物理删掉 `:300` 那一行，`sha256=ee48c311d1925a2b7188203239a051cba968668884804e9f8e1b731ab3a9dd95`（与实现方 M01 的号**同一个**，我独立跑出来的）。必须红的用例＝**6 枚**：
`internal/memory` 里 `TestAppendSessionTurnSurfacesTheConversation`／`TestAppendSessionTurnRefreshUsesStoredClockFamily`／`TestAppendSessionTurnCrashAtomicAcrossKill`／`TestAppendSessionTurn326v1RefreshSharesTheTurnTransaction`／`TestAppendSessionTurn326v1TriggerBlocksNothingElse`（后两枚＝本腿件），
`internal/session` 里 `TestSessionAppendTurnMovesItBackToTheFront`。红句逐字（全量＝`logs/10-*`、`logs/90-*`）：

```
dao_session_test.go:118: appending a turn did not reach session.updated_at: sess-a.updated_at = 1777593660, want the appended turn's stamp 1777593780 (before=1777593660, delta=0)
dao_session_test.go:122: continuing a chat does not surface it: list = [sess-b sess-a], want sess-a first after one more turn in it
manage_test.go:218: continuing a conversation does not surface it: list = [sess_0a8435cd9b9403cc24b1db99d8b98448 sess_e3311e847dfd60297fda7b6c32de1bfe], want [sess_e3311e847dfd60297fda7b6c32de1bfe sess_0a8435cd9b9403cc24b1db99d8b98448] after one more turn in sess_e3311e847dfd60297fda7b6c32de1bfe
manage_test.go:222: AppendTurn wrote the child row only: sess_e3311e847dfd60297fda7b6c32de1bfe.updated_at stayed 1777593660, want it past 1777593660
dao_session_append_326_test.go:86: half state after kill at 2 commits: session_turn has 5 rows (maxSeq 5, newest turn stamp 1777636804) while session.updated_at = 1777633200 — the touch and the turn did not share one transaction (want NOTHING or ALL, got a 3604s lag or a torn seq)
```
还原哈希＝`ec6ee530b190795093829eca07e6003bf866db5da3e64c47e6d0b0b1bcd1d2b2`（＝`45c4c03e` 与 `86298e24` 的 `dao_session.go` blob，`cmp` rc=0），还原后六枚复跑全 `--- PASS`。
★"未修码上这一发必红"票面那一发我**另建了真未修树**（`git archive 86298e24^`，`dao_session.go` sha `38bff760eabb5bf9…`＝`fd01e796` 同名件 `cmp` rc=0）复跑：红名册 3＋1 枚、四数 `78/38/3/2 rc=1` 与 `37/19/1/0 rc=1` ⇒ 实现方申报的号**逐字复现**。

**③ 反形（刷新挪到事务外＝先提交 INSERT、再单独 `UPDATE`）**——三族尺的实测颜色：

| 尺 | 反形下（sha `c31e1750690993a7…`） |
|---|---|
| 三枚 AC#1 终态用例 | **全 PASS**，同行证据照样 `updated_at … -> … (delta=120)` ⇒ 实现方自报的"这族尺对那一形盲"**复现成立**（⛔ 隐瞒、⛔ 夸大） |
| 真杀台件 | **逐发 6/6 枚 FAIL**；逐枚杀点命中 **23/48＝48%**（他报 56%，同量级、随时序漂）⇒ 概率性，⛔ 当终态判据 |
| ★本腿新建仪器 `dao_session_326v1_test.go` | **2/2 枚 FAIL**：`the touch and the turn do not share one transaction: the refresh failed (… 326v1 injected refresh failure (1811)) yet 1 session_turn row(s) stayed committed while session.updated_at = 1777593600` |

⇒ **〔缺世界〕这一格今天不再缺。** 实现方上交的两条路（(α) 接受 56%／(β) 另立一票开"tx 内第 N 句后注入失败"的测试接缝）**都不必走**：接缝已存在且⛔ 动产码——用真实 SQLite 的 `BEFORE UPDATE … RAISE(ABORT)` 触发器把母表刷新强制打失败，"两条语句同一枚事务"就等价于"刷新失败⇒子表⛔ 留行"，0.14 秒、⛔ 杀进程、⛔ 依赖墙钟。交付态绿（回滚到零行），未修码红，反形红。
★必须说清的代价：AC#1 判"成立"里"同一枚事务"那一半**今天靠的是本腿这枚新仪器**；编排者若把它搬走／删掉，那一半就退回"只有 48% 概率牙的杀台件"，我就该改判。⇒ 请拍：**留**（推荐，名册代价＝`internal/memory` 顶层 `41→43`、`RUN 78→80`）或**搬走并另立一票**。

**"这台件是不是恒过"**（派单那一问，我照做了 MU-D：刷新留在事务内、两条语句之间 `os.Exit`，门控 `WISP_326V1_EXIT_MID_TX=1`，sha `21dd5d710b6b2093…`）：
- 门控关（同一枚二进制，正控）＝`ok … 2.043s` ⇒ 那枚 `if` 本身⛔ 改行为
- 门控开＝**没报 `HALF-STATE`**，红在这句：`dao_session_append_326_test.go:86: child announced 0 commits within 60s, want >= 2 (台件故障，不是判据)` ⇒ `--- FAIL (60.03s)`
- 我的判读（⛔ 顺着派单话术写）：**派单预期的那个颜色在这形物理不可达，也不该可达**——未提交事务被 `os.Exit` 掐死正是"要么全有要么全无"的**全无**那一半。这台件在这形没有放它过，而是被它自己的前提控制（announce 计数）拦红；结合三形实测（去牙 8/8、反形逐发 6/6、未修码 8/8），**这台件⛔ 恒过＝有牙**，AC#3 判成立不靠它一枚，靠"三形都拦得住＋交付态 40 枚杀点 0 命中"。

## AC#4 四数＋门禁（尺口径写死；两把⛔ 混用）

`RUN=`＝`grep -c '^=== RUN'`（含子测试）⇄ `PASS_TOP/FAIL_TOP/SKIP=`＝`grep -c '^--- PASS|FAIL|SKIP'`（只数顶层）。包级一律 `-count=1 -timeout 30m`。
| 把 | rc | 读数 |
|---|---|---|
| `go build ./...`（导出树，七件套带 `frontend`，745 枚 `.go`） | 0 | 干净 |
| `go vet ./internal/memory/ ./internal/session/` | 0 | 零行 |
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/memory internal/session`（真身） | 0 | 列出 **0 枚**；对照：裸名 `gofumpt`＝`command not found`；⛔ 拿 `gofmt -l` 冒充 |
| `sh scripts/d22scan.sh`（仓根工作树，含别腿在飞字节＝台面读数） | 0 | `clean - no D22 ban violations`（分母 `internal/=240`、`cmd/=41`、`frontend/=116`、`internal/tools/=25`、`design/=39`、`internal/=577`、`cmd/=139`）；**rc=2 GUARD 那一发本腿⛔ 撞上** |
| 同尺打在导出树 | 1 | **台面故障⛔ 计入判据**：红的是读运行库的三枚 `TestScannerSelfScanOfRealRepoIsGreen`／`TestScopeReportMatchesRealCoverage`／`TestRealRepoLedgerIsHonest`（树里 `^\.git$` 枚数＝0）⇒ 具名到行、⛔ 代码红 |
| `sh scripts/check-path-length-budget.sh --with-self-test`（仓根） | 0 | `VERDICT GREEN`；分母 14264、越帽 60＝名册 60、名册外 0、最长相对 185、runner 最坏 229（帽 165） |
| 同尺打在导出树 | 2 | 尺自述拒答 `not inside a git work tree`＝**GUARD 射程拒绝，⛔ 读成红** |
| `go test ./internal/memory/ -count=1 -v`（交付字节） | 0 | `RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP=2` |
| 同上（＋本腿判据件，＝HEAD `c2e1c07a` 态） | 0 | `RUN=80 PASS_TOP=43 FAIL_TOP=0 SKIP=2` |
| `go test ./internal/session/ -count=1 -v` | 0 | `RUN=37 PASS_TOP=20 FAIL_TOP=0 SKIP=0` |
两枚 `SKIP` 点名＝`TestSubprocessCrashWriter`＋`TestAppendSessionTurnCrashChild`（父／子角色互斥，形状照抄既有件；实现笔新增 `t.Skip` 枚数＝2、两句都在**新建件**里，尺＝`^+.*t\.Skip`）。
时钟那一格＝**成立**：`:300` 写的是 `t.CreatedAt`；整函数体只 1 枚读时钟句子＝`:274 orNow(t.CreatedAt)`，位置在闭包**外**；闭包体 `:279-:302` 内读时钟 **0 枚 rc=1**（正控：同尺打 `:170-200` 命中 2）；`dao_session.go` 亦⛔ `import "time"`。

## AC#5 沿到消费者（⛔ 把"库里的值对了"说成"用户看得见"，也⛔ 反过来判它修错）

复认编排者两把尺：① `git grep -nE 'ListSessions|SessionStore' HEAD -- cmd internal/panel internal/session` 滤 `_test`＝**2 行**（`manage.go:40` 接口声明、`:104` 唯一转调）⇒ `cmd/`＋`internal/panel/` 零调用者；② 排序两行今天仍在 `:125`／`:144`。
票面三把尺我逐字复跑（各带正控）：`session\.NewManager` 全树 **0 枚 rc=1**（正控 `func NewManager`＝`manage.go:68`）；`"sessions"|"history"` 在 `cmd/wisp/main.go` **0 枚 rc=1**（正控＝同文件 `case "…"` 名册 `run/providers/doctor/secret/models/slo` 命中得到）；`memory\.Session` 在 `cmd/ internal/panel/ internal/tools/ internal/agent/` **0 枚 rc=1**（正控 `internal/session/manage.go` 18 枚）。
★比票面深一跳（我数的，⛔ 抄）：`AppendSessionTurn` 生产调用点 **1 枚**＝`internal/session/manage.go:150`（⛔ 拿 `go build` rc=0 冒充接线）；上游 `Manager.AppendTurn`（`manage.go:149`）生产调用点 **0 枚**；`cmd/` 里 `session.*` 标识符名册 8 枚全是**审批台账那一族**（`SessionID`/`Ledger`/`Record`/`NewLedger`/`Mint`/`ID`/`LedgerOptions`/`Covering`），会话 `Manager` 族 **0 枚** ⇒ 包被 import⛔ 等于列表接上。⇒ 今天界面上⛔ 看得见这一格，成因是**两跳**（读者没接＋写侧调用者也没接），都属票 166 的接线缺口，都撞 `cmd/wisp/**`＋`internal/panel/**`＝本波别的腿的面，本票⛔ 碰。

## 复认上交料（只复认存在，⛔ 替编排者裁）

1. **tie-break**：两行逐字节未动（尺见 AC#2），`id DESC` 与 `rand.Read` 铸随机 hex 的形状今天照旧 ⇒ "本票没顺手治它"**成立**；实现方那份暴露面（同秒 200 对硬币 45%⇄跨秒 10 对全对）我**⛔ 复跑**（那是"要不要治"的输入、⛔ 是"修得对不对"的判据）⇒ 缺的那一发＝`cp zz326r1-tiebreak_bench_test.go` 进导出树后 `go test -bench` 两发，要就说。
2. **缺母行那一形**（`UPDATE … WHERE id=?` 命中 0 行时净效果仍是"只多一行子表"）：票面⛔ 说、本层注释明写"stores rows，策略归 `internal/session`" ⇒ **未定义即停，⛔ 自填**成立；附尺复跑＝`git grep -n 'foreign_keys' -- '*.go'` 0 命中 rc=1 ⇒ ddlV3 的 `REFERENCES session(id)` 今天不生效。判⛔ 判 `ErrNotFound`＝编排者拍。
3. **回填／迟到戳**（`updated_at = t.CreatedAt` 被过去的戳写回）：今天唯一生产调用者 `manage.go:149-150` 传 `m.now().Unix()` ⇒ 生产不可达；防它要 `MAX(updated_at,?)`＝⛔ 票面 ⓐ 形状 ⇒ 交编排者拍。

## 交付缺陷具名（实现方那件 172 行／17568 字节，我逐条对尺后只有这些）

- **D1（射程标注缺）**：`docs/evidence/s1/326-append-updated.md` §7 写"`ListSessions`／`SearchSessions` 的非测试命中＝只有 DAO 自己的定义 ⇒ 产码调用点 **0 枚**"——严格讲 `internal/session/manage.go:104`/`:113` 也是非测试（产码）转调句子，枚数 **2**；"0 枚"只在 `cmd/`＋`internal/panel/`＋`internal/tools/`＋`internal/agent/` 这个射程内成立。数没错、**缺射程**＝违反本仓"引任何数必带尺＋射程"那条。
- **D2（引文不全）**：§2 未修码红句只贴了 `manage_test.go:218`/`:222` 两行；我那一发实测同枚用例红句是 **4 行**（还含 `:214` 同行证据与 `:226` "the touch and the turn disagree on the clock"）。⛔ 造假，是节选，但申报里没写"节选"。
- **D3（关键颜色漏报）**：§4 只报反形"18/32＝56% 命中"，**没报逐发颜色**＝我数的是逐发 **6/6 枚 FAIL**（4 轮 14/32 ＋另两发 5/8、4/8）。这一枚直接影响编排者拍 (α)/(β)：杀台件对**整发 CI 颜色**其实是稳的，不稳的是逐枚杀点归因；申报缺这一维。
- 复认⛔ 缺陷的两条（写明以免被当缺陷追）：申报"AC#1 三枚终态用例对反形⛔ 响"＝**实测成立**；申报四数／门禁 rc／还原哈希 `ec6ee530…`＝**逐字复现**。lag 那族数字我复跑是 `3602s..3649s`（他报 `3602..3646`）＝随时序漂、同形，⛔ 计缺陷。

## 禁区自证／没裁的格（具名）

- **本腿没做的事（禁区自证）**／没裁的格（具名）：⛔ 动产码、⛔ 改任何既有断言（尺＝`86298e24` 四件删除行**全 0** ＋ `git diff 86298e24^ 86298e24 | grep -cE '^-.*\.(Errorf|Fatalf)\('`＝**0 枚 rc=1** ＋ `git diff --numstat 86298e24 HEAD -- internal/memory internal/session` 只列本腿那枚新件 125 行）、⛔ 删用例、⛔ 放宽、⛔ 翻框、⛔ push、⛔ 跑 `./cmd/wisp` 整包与 `./...`、⛔ 杀 PID 21892、⛔ 碰脏 `go.mod`（**仍脏**：mtime `Oct 10 22:28`、`git diff --stat`＝1 insertion/1 deletion，归编排者处置）
- 我读过的既有件逐枚**只读没动**：`internal/memory/dao_session.go`、`dao_session_test.go`、`dao_session_append_326_test.go`、`schema_test.go`、`tempdir_resolved_124_test.go`、`dao_tasklog.go`、`open.go`、`writer.go`、`internal/session/manage.go`、`manage_test.go`、`internal/session/session.go`、`scripts/d22scan.sh`、`scripts/check-path-length-budget.sh`、票 326、`docs/evidence/s1/326-append-updated.md`
- ⛔ 裁的格：**没有**（六枚全裁）。⛔ 复跑的读数两发已具名于上交料 1 与本节 D3 的对照（tie-break 暴露面 bench、实现方 M02/M03 两枚"刷成常量／刷成 `created_at`"反形我没重跑，它们的判据作用已被我的去牙＋反形＋die-mid-tx 三发覆盖；要字节级复现缺的那一发＝同一导出树 `sed` 换 `SET updated_at=0` 与 `SET updated_at=created_at` 各一发四色对照）

## 交件自证（三把尺）

① 逐笔提交＋`git show --name-only`（只应含写面 1/2/3/4）：`c2e1c07a`＝`internal/memory/dao_session_326v1_test.go`（写面 4 的新建件）；`ed9c81d3`＝`.scratch/wisp/probes/326/v1/**`（17 枚）；本笔＝`docs/evidence/s1/326-v1-verdict.md`；末笔＝票 326 末尾**一条** bullet。逐枚名册由编排者现跑 `git show --name-only <号>` 复核。
② 绝对路径：`D:\work\workspace\projects plans\Wisp\docs\evidence\s1\326-v1-verdict.md`（本件，字节与行数见票尾提交后的现量）；读数件 `…\.scratch\wisp\probes\326\v1\00-readings.md`＝**132 行／20290 字节**；日志 15 枚逐枚非空（`wc -c` 全 >0，`logs/wp-with-326v1.txt` 最大）。
③ 决定性读数：去牙红句＋还原哈希 `ec6ee530…`（AC#1②）；反形红句（AC#1③）；杀台件交付态 **40 枚杀点 0 命中**、三形对照（AC#3）；两包四数 `78/41/0/2`·`37/20/0/0` 与未修码 `78/38/3/2`·`37/19/1/0`，一律 `-timeout 30m`（AC#4）；门禁 rc 逐枚 0/0/0/0（含两把导出树红的**具名台面故障**）；能力类调用点枚数＝`AppendSessionTurn` 生产 1 枚、`Manager.AppendTurn` 生产 0 枚、`session.NewManager` 生产 0 枚（AC#5）。
