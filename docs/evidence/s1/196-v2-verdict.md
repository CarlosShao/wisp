# 196-v2 终裁件：`196-r2` 对「`applyV4` 的 `DROP TABLE` 靠 FK 关着才安全」的收口够不够格

**票** 196（隐患终裁腿 `196-v2`，**非实现者**）｜**只 commit、⛔ push**
**被裁** `196-r2`：`0c08eca14`（码＋测试＋件）＋`7178f7f12`（件两处尺误自纠）
**锚** 开工 HEAD `234ff1997`；落笔时 HEAD `78a4cacdf`（别的腿在飞）。★`git log ef3ec49bd..HEAD -- internal/memory/schema.go internal/memory/taskstate_196r2_test.go` **只有 `0c08eca14`** ⇒ r2 的字节没被别人动过；下列读数不因 HEAD 前移而变。
**环境** Go 1.27.1 windows/amd64｜`GOFLAGS=`｜只 `unset GOTMPDIR`（⛔ 动 `TMPDIR/TEMP/TMP`）｜⛔ `CGO_ENABLED=0`｜gofumpt `/d/work/base/gopath/bin/gofumpt.exe`
**尺的射程**：每条读数都标「哪棵树／哪把尺／量的是代码还是注释」。探针 `.go` **只存在于仓外** `/d/tmp/wisp/20261023/196v2/tree*`，⛔ 进仓。

## §0 一句话裁决

★★**够格翻 `196-r2` 这一格；⛔ 够格把隐患整条注销。**
⇒ 护栏把「将来在第 3 句炸 `(787)`」变成「第 1 句可读地拒、**零语句执行**」；★**但 FK 开着时 `Open()` 仍然打不开**——隐患是**被声明并前移**，⛔ **被消除**。真修法（`open.go`）＋孤儿处置仍在别的腿，见 §2。

## §1 九项逐项判语（读数全部我自己现跑，⛔ 采信件）

### 1) 票面前提新鲜度：**基本真，两处尺过宽**
对 HEAD blob 逐点复量：`schema.go:292/294/303/316/321/322/323/361/368/398/401/407/411/416/421/438/440/443/453`、`ddlV1:80/82`、`open.go:146/148-151/152/153/394`、`196-v1-verdict.md:142`、`pending-and-issues.md:11196/23130`、`196` 票 `:59`、`326-append-updated.md:147`、`326-v1-verdict.md:78` —— **逐点命中**。
两处尺过宽（不影响结论，但正是「量错对象」那族）：
- ★① §3「`git grep 'foreign_keys' -- '*.go' …` ⇒ 本腿之后只剩我自己新增的 **7 行**」——HEAD 全是 **14 行**：`schema.go` **7**（361/368/401/407/411/416/421）＋新测试 **7**（`:4/72/134/152/162/167/177`）。正确说法是「schema.go 7 行／全仓 14 行」。
- ★② §3「**全仓**唯一挂 pragma 的地方是 `open.go:146`」——**生产**成立（`git grep '_pragma=' ef3ec49bd -- internal/memory` 只有 4 行），但 `.scratch/wisp/probes/166/v1/bench-dumpdb-main.go:18,94` 也有 `_pragma=busy_timeout(3000)`（探针）。应写「生产 `internal/memory`」。
- ★另：§2 P4a 写 `task_log=3`，但交付测试夹具 `seedLegacyV3(…, theFourLegacyWords)` 是 **4 行**；我独立跑得 **4**。3 属它自己的探针夹具（它明说塞 3 枚），⛔ 矛盾；但复跑者用交付夹具会得 4。
- ★「仓里没有『将来要开 FK』的票」：`grep -rn 'foreign_keys\|外键\|FK' .scratch/wisp/issues` 只命中 196／326 两条欠账，无开启计划 ⇒ **成立**；三条欠账行号全真。

### 2) 「今天不炸」：**真**（我自造夹具复现）
夹具＝生产 `ddlV1+ddlV2+ddlV3` 建的 v3 库（`schema_version=3`）＋**4** 枚 `task_log` 行 ＋ 1 枚 `tool_call` 行（`task_id` 从 `task_log` `SELECT` 出来）。
- 生产 `dsn()`（FK 关）：`PRAGMA foreign_keys`＝**0**；`Open()` ⇒ **`<nil>`**、`schema_version=4`、`task_log=4`、`tool_call=1`、`matched_parent=1`、`leftover_v4=0`。
- FK 开（`dsn()` 加 `&_pragma=foreign_keys(1)`，**唯一差别**；treeD＝改前形状）：生产 `Open()` 逐字
  `memory: migration 3->4 failed (database restored from backup is at …\backup\wisp.db.bak-3-4): apply: constraint failed: FOREIGN KEY constraint failed (787)`
- ★**第 3 句才炸**（逐句执行 `splitSQLStatements(ddlV4)`，5 句）：stmt 1 `<nil>`、stmt 2 `<nil>`、stmt 3 **`constraint failed: FOREIGN KEY constraint failed (787)`**；失败事务内读回 **`task_log_v4 objects=1 task_log=4 tool_call=1`** ⇒ 前两句**已**建表搬行，第 3 句 `DROP TABLE task_log` 才炸。★与件逐字一致。

### 3) ★P2（整个选型的支点）：**成立**（同连接、事务内/外两读）
- P1：生产 `dsn()` 上 `PRAGMA foreign_keys`＝**0**。
- **事务内**：`PRAGMA foreign_keys=ON` 后读回 **0**。
- **事务外**（同一 `*sql.Conn`、同一条语句）：读回 **1**（正控）。
- ★交叉验证：在 FK-开的 DSN 树上，护栏在**事务内**读到 **1** 并拒 ⇒ 事务内**读**反映连接真值，**写**才是 no-op。
⇒ ★`applyV4` 不可能自己关 FK；真正「关」只能落在 `open.go` 开事务之前。**件据此选〔甲〕的理由站得住。**

### 4) ★护栏在破坏性语句**之前**开火：**是**（比交付测试更强的一发）
FK 开＋1 枚子行的连接上，**直接**在事务内调 `applyV4(tx)`：
- 返回**可读**那句（含 `foreign_keys`／`tool_call.task_id REFERENCES task_log(id)`／`Nothing was modified.`；**不含 `(787)`**）。
- ★**同一事务、尚未回滚时**读回：`task_log_v4 objects=0`、`task_log=4`、`tool_call=1`、`schema_version=3` ⇒ **5 句一句都没执行**（⛔ 只是「回滚后残留 0」）。
- 改前形状（同夹具、删那 3 行）：`(787)`，且事务内 `task_log_v4 objects=1` ⇒ 前两句**确实跑了**。
- 生产 `Open()` FK 开后：可读拒；事后 `schema_version=3 task_log=4 tool_call=1 leftover_v4=0`。★「Nothing was modified」指**库内容**：`migrate()` 在步骤前仍会写 pre-migration backup 文件（错误串里可见），⛔ 缺陷。

### 5) ★反向判据（我自己造，⛔ 看件）
| 读数 | 树／唯一差别 | rc | 逐字 |
|---|---|---|---|
| **R1** | treeA/treeE＝交付形状 | **0** | 两枚新用例 `PASS`；`TestTaskState196DDLV4KeepsTheShape` **PASS** |
| **R2** | treeC＝treeA＋`dsn()` FK(1)（`diff -rq` 只有 `open.go` 一行） | **1** | 两枚全 `FAIL`；`Open()` **可读拒**；事后 `v3/4/1/0` |
| **R3** | treeD＝treeC **删掉那 3 行**（`diff -rq` 只有 `schema.go`，`diff -u` 正是 440–442） | **1** | 两枚全 `FAIL`；`Open()` 退回 **`apply: constraint failed: FOREIGN KEY constraint failed (787)`**；`applyV4` 直调也 `(787)` 且事务内 `task_log_v4=1` |

★**单变量**：`diff -rq treeA treeB` ⇒ 只有 `schema.go`；`treeA treeC` ⇒ 只有 `open.go`；`treeC treeD` ⇒ 只有 `schema.go`（3 行）⇒ **护栏就是那 3 行换来的**，件这句成立。
★诚实一条：R2 里 `TestTaskState196RefusesRebuildWithForeignKeysOn` 是**控制半**先 `Fatalf`（它假定 `dsn()` 不设 pragma）⇒ 主体半在 R2 没被跑到。件自己披露了；这是第二条绊线，但也是**测试形状的弱点**（见 §2(e)）。

### 6) ★边界零 diff：**逐项核，全真**
- `0c08eca14` `--name-status`：`A docs/…/196-r2-impl.md`／`M internal/memory/schema.go`／`A internal/memory/taskstate_196r2_test.go` **恰 3 枚**；`7178f7f12`：**只 M 件**；`234ff1997`：**只 M 台账**（不属 r2 写面）。
- `git diff --numstat ef3ec49bd..HEAD -- internal/memory/schema.go`＝**73/0**。⚠ 派单的 `grep -cE '^[+-][^+-]'` 在此得 **71**：有 **2 行是空白新增行（单独的 `+`）**被该模式漏掉 ⇒ **numstat 才是 73**（这是「算法不同」第 9 例的邻居）。
- `ddlV4`：`schema.go` 279–324 行 `cmp`（ef3ec49bd vs HEAD）**rc=0**；`splitSQLStatements(ddlV4)`＝**5**；形状测试 **PASS**。
- `ddlV1/V2/V3`：`schema.go` 1–355 行 `cmp` **rc=0**；全 diff **0 删除行**；`git log ef3ec49bd..HEAD -- internal/memory/taskstate.go internal/memory/open.go docs/PLAN.md docs/specs` ⇒ **空**。
- 12 枚名册：`taskStateRoster` 字面 **12** 条；`taskstate.go` **零 diff**。
- `docs/PLAN.md`＋`docs/specs`：`git diff … | wc -c`＝**0**。`cmd/**`／`scripts/**`／`.github/**`：零写（两笔名册已证）。

### 7) 门禁复跑（四数分开／rc 直取）
- `go build ./...` **0**；`go vet ./internal/memory/` **0**；`go test -count=1 ./internal/memory/` **0**（`ok … 24.448s`）。
- `sh scripts/d22scan.sh` **0** `clean`（production 290＝internal 248＋cmd 42；ban #8 internal 632／cmd 148 Go 文件）。
- `sh scripts/check-path-length-budget.sh` **0** `VERDICT GREEN`（`tracked paths=17644  over-budget=61  covered=61  not in roster=0`；比 `196-v1` 件的 17628 大是共享树前移，判据同）。
- **gofumpt 甲尺**（`git archive` 仓外 treeE，**判据**）：`gofumpt -l internal cmd tools scripts frontend` ⇒ **0 行**；两件单跑 ⇒ **0 行**。
- **gofumpt 乙尺**（工作树，**只诊断**）：**3 行**，全是反斜杠输出（归属按 `[\\/]` 判）：`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`cmd/wisp/models.go`。三枚 `git show HEAD:<p> | tr -cd "\r" | wc -c`＝**0**、`git ls-files --eol`＝`i/lf w/crlf`、工作树 CR＝行数（131/131、225/225、389/389）⇒ **工作树字节、⛔ 提交内容**，且都不是 r2 的文件。正控：故意写坏的 `bad.go` **被列出 1 枚**。
- **四数**：改后 `rc=0 RUN=93 PASS_TOP=56 FAIL_TOP=0 SKIP_TOP=2`；改前（treeBase=`ef3ec49bd`）`rc=0 91/54/0/2` ⇒ **Δ=+2/+2 恰两枚新用例**，与件**逐数复现**。⚠ 口径注：`RUN` 数的是**全部** `=== RUN`（含 **35** 个 subtest），PASS/FAIL/SKIP 是**顶层**；顶层 RUN 是 **58/56**。Δ 结论不变，但「RUN 与 PASS 同一分母」是错觉（`grep -cE '^=== RUN'`＝93 可复现）。

### 8) ★它上交的两件：**第 1 件对；第 2 件前提真、拒得对，但检查仍欠**
1. **真修法在 `open.go`**：**判对**。官方 12 步第 1 步就是「FK 开着先关」，第 10 步是重建后、提交前 `foreign_key_check`；P2 已证事务内关不掉 ⇒ 只能落在 `migrate()` 开 `BeginTx` **之前**。⚠ 补一句：`dsn()` 是 reader/writer/migration **共用**的，最稳是**迁移专用 DSN**（或迁移连接上显式 `PRAGMA foreign_keys=OFF`），⛔ 把 `_pragma=foreign_keys(0)` 一把加进共用 `dsn()`——那会在将来开 FK 时把运行时也钉死。
2. **不加结尾 `foreign_key_check`**：**前提真**。我造了一枚带**孤儿** `tool_call` 行的 v3 库：`PRAGMA foreign_key_check` 迁移**前**＝**1 条**（`table=tool_call rowid=1 parent=task_log fkid=0`），迁移**后**仍＝**1 条**；生产 `Open()`（FK 关）今天照开、`schema_version=4`。且孤儿今天**可达**：`PurgeTaskLogs`（`dao_tasklog.go:134`）删光 `task_log` 而注释明写 `tool_call` 行归自己的域，`DeleteTaskLog`（`:118`）删父不删子，FK 又关着 ⇒ 只要有 `tool_call` 行，一发 purge 就造孤儿。⇒ **硬门槛**真会让这类库打不开（派单停手条件③）⇒ **这一腿拒加是对的**。★但「不加」只能是**暂缓**：第 10 步仍欠，落 `open.go` 时应做成**基线化**检查（重建前先 `foreign_key_check` 记基线，重建后只对**新增**违规报错）或先配「孤儿清理」票；⛔ 把检查永久免掉。

### 9) 它自报的尺误：**两处都真，但同类的新一处没自报**
- 测试件 `wc -l`＝**213**（件已改；`224` 是估的）✓；`7178f7f12` `--name-status` **只 M 件** ✓。
- `sha256sum`（SHA‑256、64 hex）四读同值：`schema.go 6b985a280092becbfc0c72bcb7e932339a0e08f3a62272ea591a8c3f61084094`、测试 `9ac9ef48a7e590c78560f0864732c331673445529b0e1e12dbc6a499e8bb627f`（`0c08eca14`＝`7178f7f12`＝HEAD＝工作树）⇒ **码与测试零字节改动** ✓。
- `open.go:152/153` 两行拆分已在件里 ✓；件 **118 行** ✓。
- ★新一处：§3 的「7 行 vs 14 行」（见 §1.1①）——同一族「尺的射程」，⛔ 在它自报的两处里。

## §2 ★够不够格翻 ＋ 缺什么

**够格翻这一格**：护栏**在破坏性语句之前**、**零语句执行**、**可读**、**单变量**（3 行）、两侧都有测试，边界零 diff，门禁全 0，且它自报并改正了两处尺误。
**⛔ 够格把隐患整条注销**：FK 开着时 `Open()` 仍然**打不开**（只是从 `(787)` 变成人话）。所以台账/票面应写「**本腿交付成立；隐患＝已声明＋第 1 句可读拒；真修法另立 `open.go` 腿**」，⛔ 写「已修」。

**缺什么（按优先级）**
- (a) ★`open.go` 迁移连接：`BeginTx` **之前**关 FK（迁移专用 DSN），重建后按官方第 10 步 `foreign_key_check`——这是**消除**隐患的那一半。
- (b) ★第 10 步必须**基线化**（或先清孤儿）：我一发孤儿库实测「检查前后都 1 条」⇒ 硬门槛会拒老库；基线比较/先清理才能既查新破坏又不拒老库。
- (c) ★孤儿处置票：`PurgeTaskLogs` 今天保证「purge 后必有孤儿」；开 FK 前要先清/先修，否则 FK 一开，删父就炸。
- (d) 件里 §3「7 行」改 14（或写明只数 `schema.go`）；「全仓唯一挂 pragma」改「生产 `internal/memory`」。
- (e) 可选：把 `TestTaskState196RefusesRebuildWithForeignKeysOn` 的**控制半**显式钉 `foreign_keys(0)`（别依赖 `dsn()` 默认），这样 R2 仍能跑到主体半；否则「前提一变」时两半一起红，诊断力下降。
- (f) 门禁/四数引用时注明 `RUN` 含 35 个 subtest（顶层 58/56），⛔ 让读者以为 93 与 56 同分母。

## §3 提交号 ＋ 逐笔越界

- 本件提交号：`<待填>`（提交后回填；只 `A docs/evidence/s1/196-v2-verdict.md` 一枚，`git diff --cached --name-status` 提交前已核）。
- 逐笔越界：`0c08eca14`＝3 枚（件 A／`schema.go` M +73/−0／测试 A 213 行），无 `cmd`／`scripts`／`.github`／`PLAN`／`specs`；`7178f7f12`＝1 枚（只件）；`234ff1997`＝1 枚（只台账，非 r2 写面）。★**零越界**。
- 本腿写面：本件 ＋ 仓外 `/d/tmp/wisp/20261023/196v2/**`（探针 `treeA..treeE`、`treeBase`、日志）。⛔ 改仓内任何 `.go`／`internal/**`；工作树 `git status --porcelain -- internal/memory docs/evidence/s1` 在落笔前为空。
