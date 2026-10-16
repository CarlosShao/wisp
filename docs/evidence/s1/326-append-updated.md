# 票 326 — `AppendSessionTurn` 刷新 `session.updated_at`（落地腿 326-r1 交件读数）

- 腿＝`326-r1`（落地腿·写）｜票＝`.scratch/wisp/issues/326-append-session-turn-does-not-refresh-session-updated-so-a-live-chat-never-moves-up.md`
- 起树锚（本腿动手前 HEAD）＝`fd01e796`｜交码提交＝`86298e24`（父 `ce8df8b6`）
- 现量时刻＝2026-10-16 10:53 → 11:26 +08（每件抬头各带自己的时间戳）
- **本表⛔ 翻任何 AC 框**（票框尺起手＝未勾 6 枚／已勾 0 枚，交完仍是 6／0）。翻框＝编排者按**非实现者**裁决定。
- 原始读数目录＝`.scratch/wisp/probes/326/r1/`（`00`／`20`／`40`／`60` 四枚读数件＋`logs/` 十枚原始输出＋两枚可重放脚本）。

## 0. 落点（票面候选 ⓐ／ⓑ，本腿现读后选 ⓐ）

`internal/memory/dao_session.go` 的 `AppendSessionTurn` 在**同一枚 `s.write` 事务**里补一条
`UPDATE session SET updated_at=? WHERE id=?`，值＝`t.CreatedAt`（已经 `orNow` 归一）。
ⓑ 排除＝仓里没有可复用的"touch"助手：既有四处写侧（`:84` 建／`:177` 改名／`:198` 归档／`:231` 分叉）各自内联一条语句，
抽一枚助手＝动既有代码面、超出本票射程。理由与两处"⛔ 自造"的边界全文＝`probes/326/r1/60-open-shapes-and-landing.txt` §2／§5。

## 1. AC#0 起手复算（四行票面现量逐枚现跑）

| 尺（射程） | 读数 | rc |
|---|---|---|
| `AppendSessionTurn` 函数体内 `UPDATE session`（该函数体 27 行整块） | **0 枚** | 1 |
| `grep -n 'updated_at' internal/memory/dao_session.go`（整文件） | 7 行＝列名册 1＋建 1＋两枚 `ORDER BY` 2＋改名 1＋归档 1＋分叉 1 ⇒ 写侧**四处**，说话不在里面 | 0 |
| `grep -n 'ORDER BY updated_at DESC, id DESC'`（同文件） | **2 枚**（`:125` 列表／`:144` 搜索） | 0 |
| 时间源：`observe.NowWallUTC().Unix()`×3（:175/:192/:218）＋`orNow`（dao_tasklog.go:14）＋`manage.go:77` 默认 `time.Now` | 秒精度、同一族 | 0/0/0 |

票框尺＝未勾 **6 枚**（rc=0）／已勾 **0 枚**（rc=1）。原始件＝`logs/00-step0-measure.txt`、`logs/append-body.txt`、`00-setup-and-trees.txt`。

## 2. AC#1 成对读数（首要判据）

形状＝会话 A 先建、会话 B 后建（同一枚步进时钟、相隔一分钟），然后在 **A** 里追加一轮。
两发⛔ 跨台面：红＝基线树 `/d/tmp/w326/base-fd01e79`（`git archive fd01e796 …`，`.go` 740 枚证非空，
`dao_session.go` sha256 `38bff760eabb5bf9…`＝该锚 blob）；绿＝交付树 `/d/tmp/w326/after-86298e2`
（`git archive 86298e24 …`，`.go` 743 枚，`dao_session.go` sha256 `ec6ee530b1907950…`＝HEAD blob）。
**同屏控制**＝`git diff --numstat fd01e796 ce8df8b6 -- internal/memory internal/session` 输出为**空** ⇒ 两棵树里本票射程内
除那一枚 `UPDATE` 外逐字节相同。同一枚 `GOTMPDIR=C:\Users\swq\AppData\Local\Temp\wisp326r1`（原生反斜杠、先 `mkdir -p`）。

**改前那一发（未修码上必红，红句逐字）：**
```
dao_session_test.go:114: 326 AC#1 same-row evidence: sess-a.updated_at 1777593660 -> 1777593660 (delta=0) for turn.created_at=1777593780 | list = [sess-b sess-a]
dao_session_test.go:118: appending a turn did not reach session.updated_at: sess-a.updated_at = 1777593660, want the appended turn's stamp 1777593780 (before=1777593660, delta=0)
dao_session_test.go:122: continuing a chat does not surface it: list = [sess-b sess-a], want sess-a first after one more turn in it
manage_test.go:218: continuing a conversation does not surface it: list = [sess_c78d0af96790c68cba5051c90bc817a5 sess_71a846bb847b450056f828cf8eb63e3f], want [sess_71a846bb847b450056f828cf8eb63e3f sess_c78d0af96790c68cba5051c90bc817a5] after one more turn in sess_71a846bb847b450056f828cf8eb63e3f
manage_test.go:222: AppendTurn wrote the child row only: sess_71a846bb847b450056f828cf8eb63e3f.updated_at stayed 1777593660, want it past 1777593660
```
（同一发里还有第三枚 DAO 用例红句：`updated_at did not move forward: 1775001600 -> 1775001600 (appending is a touch)`。）
原始件＝`logs/10-red-before-fix.txt`。

**改后那一发（同一台面、成对绿）：**
```
dao_session_test.go:114: 326 AC#1 same-row evidence: sess-a.updated_at 1777593660 -> 1777593780 (delta=120) for turn.created_at=1777593780 | list = [sess-a sess-b]
--- PASS: TestAppendSessionTurnSurfacesTheConversation (0.13s)
manage_test.go:214: 326 AC#1 outlet reading: sess_12ae….updated_at 1777593660 -> 1777593780 (delta=120) for turn.created_at=1777593780 | list = [sess_12ae… sess_1cce…]
--- PASS: TestSessionAppendTurnMovesItBackToTheFront (0.08s)
```
★**同行证据**（⛔ 只报"列表顺序对了"）＝同一发里的 `updated_at 1777593660 -> 1777593780 (delta=120)`
与 `turn.created_at=1777593780` **同值**：追加那一轮的戳被逐字写进了母表那一列。另一枚用例钉 `CreatedAt==0` 那支
（`orNow` 走 `observe.NowWallUTC`）⇒ `updated_at == 子表读回的 turn.created_at` 且 `> 追加前`。原始件＝`logs/12-green-after-fix.txt`。

**整包颜色（同一把尺、逐包四数）：**
| 包 | 未修码（基线树＋本腿测试件） | 交付（86298e24 树） |
|---|---|---|
| `internal/memory` | `RUN=78 PASS_TOP=38 FAIL_TOP=3 SKIP=2` rc=1（红的正是本腿三枚） | `RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP=2` rc=0 |
| `internal/session` | `RUN=37 PASS_TOP=19 FAIL_TOP=1 SKIP=0` rc=1 | `RUN=37 PASS_TOP=20 FAIL_TOP=0 SKIP=0` rc=0 |
（`RUN=`＝`^=== RUN`，`PASS_TOP=`＝`^--- PASS` 只数顶层，⛔ 混用。`SKIP=2` 里一枚是本腿 `TestAppendSessionTurnCrashChild`
的父／子角色互斥 skip，形状照抄既有 `TestSubprocessCrashWriter`，⛔ 为绕过任何判据。原始件＝`logs/13-baseline-wholepackage.txt`、`logs/30-ac4-gates-exported.txt`。）

## 3. AC#2 ⛔ 顺手治 tie-break（判据＝那两行逐字节不动）

- 尺＝`git diff fd01e796 -- <本腿四件> | grep -E '^[+-].*ORDER BY updated_at DESC, id DESC'` ⇒ **命中 0 枚、rc=1**；
  diff 本体 255 行（⇒ 尺⛔ 空转）。逐字节对撞＝那两行切片 sha256 在锚与盘上**同值** `6a4e53226589a6fb…`。
- 正控＝同一把尺打在仓外副本的假改动（删掉 `, id DESC`）上 ⇒ **命中 2 枚**、rc=0 ⇒ "0 命中"⛔ 是尺坏。
- ⛔ 动签名：`grep -E '^[-+].*func \(s \*Store\) AppendSessionTurn'` 对 diff ⇒ **0 枚、rc=1**；正控（对新增语句）⇒ 1 枚。
- **上交的料（只测不改）**：修完之后同秒先后的暴露面。同一秒连发 200 对 ⇒ touch **0/200 跨秒**、200/200 仍并列 ⇒
  次序落到随机 hex 的倒序，硬币 **90/200＝45%** 判给没被碰的那枚（与 `166-v1` 量的 45.5% 同量级）；
  跨秒 10 对 ⇒ **10/10 跨秒、0/10 判错**。⇒ 本票把"永远错"改成了"跨秒即对、同秒仍掷硬币"。
  件＝`probes/326/r1/zz326r1-tiebreak_bench_test.go`（⛔ 进产码面：`.scratch` 以点开头的目录 go 工具链整目录跳过；
  跑法是 `cp` 进导出树），读数＝`logs/50-ac2-tiebreak-exposure.txt`。**要动那两行或 `session.id` 的形状，⛔ 本腿，由编排者另裁。**

## 4. 牙（判据仪器自己）

四枚被牙用例＝`TestAppendSessionTurnSurfacesTheConversation`／`TestAppendSessionTurnRefreshUsesStoredClockFamily`／
`TestAppendSessionTurnCrashAtomicAcrossKill`（真杀进程）／`TestSessionAppendTurnMovesItBackToTheFront`（Manager 出口面）。
全部突变打在**交付导出树的盘上字节**（⛔ `-overlay`、⛔ 共享工作树），逐发记 sha256、跑完 `cp` 还原复跑。

| 发 | 改的字节 | sha256(dao_session.go) | 颜色 |
|---|---|---|---|
| M00 正控 | 交付态未动 | `ec6ee530b1907950…` | 四枚全绿 |
| **M01 去牙** | 那条 `UPDATE` **整行物理删除** | `ee48c311d1925a2b…` | **四枚全红** |
| M02 反形 | 刷成常量 `0` | `e3b23fe9d34c3b37…` | **四枚全红**：`updated_at 1777593660 -> 0 (delta=-1777593660)`；真杀台件 8/8 `HALF-STATE` |
| M03a ⚠ | sed 吃掉行尾反引号 ⇒ `[build failed]` | `646a0ce00a8c500b…` | **台面故障、⛔ 计入判据**（重做为 M03b） |
| M03b 反形 | `SET updated_at=created_at` | `1aa46bc0c5cd7fb2…` | **四枚全红**：`…stayed 1777593660`；真杀台件 8/8 `lag=3607s..3645s => HALF-STATE` |
| M04 反形 | `UPDATE` **挪出事务**（提交后另起一枚 `s.write`） | `3c3fda6ce0605e46…` | 真杀台件 **18/32 枚杀点判 `HALF-STATE`、4/4 轮 FAIL**；AC#1 那三枚**放行** |
| M99 还原 | `cp` 回交付字节 | `ec6ee530b1907950…` | 四枚复绿；`cmp` 对 pristine **rc=0**；对 `git show 86298e24:…` 同 sha |
| M99b | 交付字节跑 M04 同轮数（32 枚杀点） | `ec6ee530b1907950…` | **`HALF-STATE` 0 枚**、4/4 轮 PASS ⇒ 那 18 枚红⛔ 是仪器自造 |

★**具名的不敏感**：**AC#1 那三枚尺对"刷新挪到事务外"⛔ 敏感**（只看终态、⛔ 看事务边界）；管这一形的是真杀台件，
命中率 **18/32＝56%**（同秒戳相等时那一形会躲过去）。要钉成确定性红＝需要一枚"事务内第 N 句后可注入失败"的接缝，
那要动 `Store.write`／DAO 函数体形状＝超本票射程 ⇒ 交编排者裁：(α) 接受"56% 命中＋交付态 0 命中"；(β) 另立一票开测试接缝。
去牙那形（M01）与两枚反形（M02／M03b）都是**四枚全红**，⛔ 只靠顺序：同屏那行 `delta=0` 与 `updated_at = 0`／`= created_at` 的红句是同枚读数的两面。
可重放脚本＝`probes/326/r1/de-tooth.sh`＋`de-tooth2.sh`；原始输出＝`logs/40-mutations-batch1.txt`＋`logs/41-mutations-batch2.txt`。

## 5. AC#3 原子性＋迁移不破

- **同一枚事务**：`AppendSessionTurn` 的两条语句都在 `s.write` 的那一枚 tx 里（`open.go:280`：`fn` 返回 error ⇒ 整枚回滚）。
- **真杀进程台件**（沿票 166 AC#3 立过的形状，兄弟＝既有 `TestCrashRecoveryKillMidWrite`）：父进程 `cmd.Process.Kill()`
  杀掉正在追加轮次的子进程（子进程每枚提交后 announce 一行），再从生产路径 `Open(dir)` 重开（WAL 恢复＋启动 checkpoint），
  读数只允许 `NOTHING`／`ALL` 两形。子进程的轮次戳**每秒一枚**（`base+i`）而不是取墙钟——否则整个循环落在一秒内，
  "根本没刷新"与"刷新在另一枚事务里"都满足任何时间戳比较，台件会因错误的理由变绿（这一形本腿真的踩到：
  第一版台件在未修码上判 `ALL`，改成调用方自铸戳后才判出 `HALF-STATE`；两版读数都在 `logs/`）。
- 交付态：**8 枚杀点（5 枚定点＋3 枚派生）全判 `ALL`、`lag=0s`、`integrity_check=ok`、⛔ 一枚 `HALF-STATE`**；
  ×4 轮＝32 枚杀点 0 命中。未修码：同样 8 枚杀点**全判 `HALF-STATE`**（`lag=3602s..3646s`）。
- **迁移那族⛔ 破**：`git diff fd01e796 86298e24 -- dao_session.go | grep -cE '^[+-].*(applyV3|ddlV3|migrationChain|SchemaVersionTarget)'`
  ⇒ **0 枚、rc=1**；既有那枚 `--- PASS: TestSessionMigrationCrashAtomicAndReplay (0.55s)` 在交付态照旧绿，
  兄弟 `--- PASS: TestCrashRecoveryKillMidWrite (0.19s)` 照旧绿。

## 6. AC#4 门禁六把（各落 rc）

第一组＝交付导出树（⛔ 别腿在飞字节）；第二组＝仓根共享工作树（含别腿字节 ⇒ 只当台面读数）。
| 把 | 树 | rc | 读数 |
|---|---|---|---|
| `go build ./...` | after-86298e2（743 枚 .go） | **0** | 干净（六件套之外补抽 `frontend`，⛔ 写面、只当分母） |
| `go build ./...` | 仓根 `f4b62861` | **0** | 台面读数（含 `220-r2` 已提交字节） |
| `go vet ./internal/memory/ ./internal/session/` | after 树 | **0** | 零行输出 |
| `sh scripts/d22scan.sh` | 仓根 | **0** | `clean - no D22 ban violations`；bans #1-5 internal/=239、cmd/=41｜#6 frontend/=116｜#7 internal/tools/=25｜#8 design/=39、frontend/=116、internal/=573、cmd/=139。`GOTMPDIR` 在 C 卷＝`%LOCALAPPDATA%` 同卷 ⇒ 票 311 GUARD 放行（⛔ 这次 rc=2） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 仓根 | **0** | `VERDICT GREEN`；分母 14046 枚跟踪路径、越帽 60＝名册 60、not-in-roster 0、worst 229；positive control PASSED |
| 格式化尺 `gofumpt` **绝对路径真身** | 两棵树 | **0** | `/d/work/base/gopath/bin/gofumpt.exe -l internal/memory internal/session` ⇒ **列出 0 枚**；本腿四件单独再跑 ⇒ 0 枚。裸名 `gofumpt` 对照＝`command not found`（⛔ 在 PATH）；⛔ 拿 `gofmt -l` 的 0 冒充（同四件它也是 0，此处两把⛔ 混，报的是真身那把） |
| `go test ./internal/memory/ -count=1 -v` | after 树 | **0** | `RUN=78 PASS_TOP=41 FAIL_TOP=0 SKIP=2`（`ok … 21.148s`） |
| `go test ./internal/session/ -count=1 -v` | after 树 | **0** | `RUN=37 PASS_TOP=20 FAIL_TOP=0 SKIP=0`（`ok … 2.351s`） |

⛔ `go test ./cmd/wisp` 整包（票 314 死亡区）⛔ 跑；⛔ `go mod tidy`／`go get`。红句里⛔ `walk`／`build failed`／
`no tests to run`／`0xc0000135`／`RUN=0`（唯一一枚 `build failed` 在 M03a，已具名为台面故障并重做）。
原始件＝`logs/30-ac4-gates-exported.txt`、`logs/31-ac4-gates-repotree.txt`。

## 7. AC#5 生产可见性：说实话

**修完之后界面上仍⛔ 看得见。** 现跑（HEAD `f4b62861`，11:23–11:26 +08）：
`git grep -c 'session\.NewManager' -- '*.go'` ⇒ **0 枚、rc=1**（正控 `func NewManager`＝`manage.go:68` 1 枚）；
`grep -cE '"sessions"|"history"' cmd/wisp/main.go` ⇒ **0 枚、rc=1**；
`git grep -n 'memory\.Session' -- cmd/ internal/panel/ internal/tools/ internal/agent/` ⇒ **0 枚、rc=1**（正控 `internal/session/manage.go` 18 枚）；
`ListSessions`／`SearchSessions` 的非测试命中＝只有 DAO 自己的定义 ⇒ **产码调用点 0 枚**。
⇒ 差的那一跳＝装配根构造一枚 `session.NewManager` ＋一枚读者（子命令或面板快照），属 `cmd/wisp/**`＋`internal/panel/**`
＝本波唯一写腿 `220-r2` 的面，**本腿一个字节⛔ 动**（自证＝本笔提交名册只有四件、全在五处写面内）。

## 8. 本腿⛔ 治、具名上交的三形（全文＝`probes/326/r1/60-open-shapes-and-landing.txt`）

1. **缺母行那一形**：`UPDATE ... WHERE id=?` 影响 0 行时今天的净效果仍是"只多一行子表"。判成 `ErrNotFound`＝新增失败模式、
   票面⛔ 说、且本层注释明写"stores rows，生命周期策略归 `internal/session`" ⇒ 按**未定义即停**⛔ 自填。
   附尺：`git grep -n 'foreign_keys' -- '*.go'` ⇒ **0 命中、rc=1**；`dsn()` 只挂 `busy_timeout/journal_mode/synchronous/wal_autocheckpoint(＋_txlock)`
   ⇒ ddlV3 里那句 `REFERENCES session(id)` **今天根本不生效**，孤儿轮次在本票之前就可达。
2. **回填/迟到戳**：值取 `t.CreatedAt` ⇒ 若调用方给过去的时间戳，`updated_at` 会被写回那一秒（可能下沉）。
   今天唯一调用者 `manage.go:149 AppendTurn` 传 `m.now().Unix()` ⇒ 生产不可达；防它要 `MAX(updated_at,?)`＝⛔ 票面 ⓐ 形状，交裁。
3. **tie-break**：见 §3，两行逐字节未动，只交暴露面读数。

## 9. 交件自证三把尺

1. **提交号＋路径名册**：`86298e24`＝`git show --name-only --format=` ⇒ `internal/memory/dao_session.go`、
   `internal/memory/dao_session_append_326_test.go`、`internal/memory/dao_session_test.go`、`internal/session/manage_test.go`
   ＝**只含五处写面之内**（1／2 号面）。其后的读数与票尾 bullet 各自成笔、⛔ 改已入库字节。
   ⛔ `cmd/wisp/**`／⛔ `internal/panel/**`／⛔ `internal/winsec`／⛔ `internal/risk`／⛔ `scripts/**`／⛔ `tools/**`／
   ⛔ 三枚冻结件／⛔ golden／⛔ `thresholds.go`／⛔ `go.mod`／`go.sum`（尺＝名册逐枚比对，见 `logs/31-*` 第 5 节）。
2. **证据件绝对路径**：本件＝`D:\work\workspace\projects plans\Wisp\docs\evidence\s1\326-append-updated.md`；
   读数目录＝`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\326\r1\`
   （`00-setup-and-trees.txt`／`20-ac2-orderby-untouched.txt`／`40-de-tooth-and-teeth.txt`／`60-open-shapes-and-landing.txt`
   ＋`logs\` 十二枚原始输出＋`de-tooth.sh`／`de-tooth2.sh`／`zz326r1-tiebreak_bench_test.go`）。字节数现量＝票尾 bullet。
3. **未修码红句逐字＋改后成对读数＋去牙必须红与还原绿（含还原哈希）**：§2／§4／§5，
   还原哈希＝`ec6ee530b190795093829eca07e6003bf866db5da3e64c47e6d0b0b1bcd1d2b2`（＝`git show 86298e24:internal/memory/dao_session.go` 的 sha256，`cmp` rc=0）。

**没做完的格子**：本票五格（AC#1／AC#2／AC#3／AC#4／AC#5）⛔ 有任何一枚由本腿勾；AC#2 那格的"要动⛔ 动"⛔ 是缺陷、是票面明令射程外。
读数缺口＝§4 具名的那一枚：**"刷新挪到事务外"缺确定性判据**（只有真杀台件 56% 命中），要补需要一枚测试接缝，本腿⛔ 造。

**规程级自纠（入账）**：① 本腿踩了派单点名的 shell 陷阱一枚（双引号串里放反引号 ⇒ 路径被当命令替换，见 `logs/50-*` 末尾自纠）；
② 第一版真杀台件因"循环全落同一秒"在未修码上假绿，本腿自己发现并改成调用方自铸每秒一枚戳（两版读数都留在 `logs/`）；
③ `git archive` 六件套⛔ 覆盖 `go build ./...`（`internal/panel` 要 `frontend` 那枚包），建议抄进规程。
