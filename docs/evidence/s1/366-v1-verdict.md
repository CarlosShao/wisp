# 票 366 / 终裁腿 `366-v1`（非实现者）— `AC#2` **⛔ 不能翻**：形状与读数都对，缺一枚「坏尺」用例（＋措辞未改＋一枚 hash 不可复现）

- 腿＝`366-v1`（**非实现者裁腿**）· 裁的对象＝`366-r1`（`d4ce46f97`／`2abf150b2`／`c2e47f70c`，件 143 行）· 裁时 HEAD＝`434a1d4ab`（编排者收 r1 那一笔）· 2026-10-23 18:4x–19:1x +08
- 形已由编排者拍＝〔甲〕**基线化**（跑 `foreign_key_check`、**只记录**、⛔ 拒开库）⇒ 本腿**只裁"落成没有"**，不重开形。
- 写面＝本件 ＋ 仓外 `D:\tmp\wisp\20261023\366v1\**`：`headtree`＝HEAD 导出树 · `pretree`＝`4ee57d318`（改前）导出树 · `muttree-m1/m2/bad`＝三份变异副本 · `logs/`＝本件所有读数的原始输出 · `probes/`＝脚本。**仓内 `.go` 0 字节**；**工作树全程未当突变靶**（⑦ 末）。
- 尺的射程（先写清，⛔ 只说"跑过了"）：夹具＝**生产链建的真 v3 库**（`Open(dir, withSchemaTarget(3))`）；孤儿＝**生产调用**（`PurgeTaskLogs`／`DeleteTaskLog`／raw INSERT／`RunRetentionOnce`）；迁移＝生产 `Open(dir)`；`fkc`＝**迁移后另开一条 `dsn()` 业务连接**数 `PRAGMA foreign_key_check` 行数；计数＝`Store.ForeignKeyAudit()`；行数／`task_id`＝同一业务连接直接读。探针**只在导出树里**编译运行（`go test -run TestV1Probe366*`），⛔ 工作树。

---

## ① 票面前提新鲜度（对 HEAD blob 复量）

| 票面/普查件的话 | HEAD 复量（我自己跑） | 判 |
|---|---|---|
| 重建自称走"官方 12 步" | `schema.go:299-300` 逐字仍在；**普查件引的 `:283-284` 已被 r1 自己的 +122 行推到 `+16`** ⇒ 普查件行号**过期** | ✅ 事实在、行号漂 |
| 第 1 步（关 FK）`196-r3` 已补 | 实读：迁移连接 `PRAGMA foreign_keys=**0**`、DSN 串含 `foreign_keys(OFF)` 且**不含** `(1)`；业务 DSN **零** `foreign_keys`、实读也是 **0** | ✅ 仍在 |
| "第 10 步**从来没跑过**" | 已被本票自己的实现改写（`git grep -c 'foreign_key_check' HEAD -- '*.go'`＝schema 10＋open 5＋新用例 13＝**28 行**） | ⚠ **过期（正常）** |
| 孤儿**三条**路径 | 我造**四条**全部出孤儿（purge／delete／manual／retention） | ⚠ 票面"三条"过期 |
| `196-v2` 实测"迁移前后都＝1" | 方向真；数字是**夹具行数**不是常量（我量 3/2/1/1） | ✅ 与 `366-a1` 同判 |
| `dao_tasklog.go` 行号 | `:118`／`:134` 删除句、`:130` "own domain" 注释，**未漂** | ✅ |

## ② ★AC#2 判语：**⛔ 不能翻**（缺三项；行为层我已全部核过，缺的都不在行为）

- **缺 1（行为覆盖）**：`auditForeignKeys` 的"**审计自身读失败 ⇒ Warn＋`Checked=false`、照常开库**"那一支**没有用例**（r1 §5 已如实自报）。这是本腿**新增的一条判断**，也是唯一紧贴硬边界 #1 的分支 ⇒ AC#2 原文"**每处改动配一枚今天响、不修不响的用例**"未满足（判语见 ⑨）。
- **缺 2（措辞）**：新增注释**仍把"官方第 10 步"当已核事实**（8 处，见 ⑩）⇒ 编排者已拍的"改成**额外保险（可能与官方第 10 步重合）**"**没做到**。
- **缺 3（证据可复现）**：`还原核` 4 枚 sha256 里 **`open.go` 一枚与入库字节不符**（它报 `3d9f2655…`、入库 blob＝`d13000b1…`，见 ⑫ 附-1）⇒ 反向判据的还原核**不可复核**。
- ✅ **正面**：形状、位置、两处记录、四条路径读数、连开两次、边界零 diff、`M1`/`M2` 红行——**我全部独立复现**（③–⑦）。⇒ 缺的是**补牙＋改措辞＋重记一个 hash**，⛔ 返工设计。
- **落点核（票面最硬的那句）**：`open.go:509 step.apply` → `:517-520 auditForeignKeys` → `:521-527 写水位` → `:528 Commit` ⇒ 审计**恰在"apply 之后、写水位与 Commit 之前"**；`migrationStep.auditForeignKeys` 由步骤自报、链里**恰好 `3->4` 一枚**（我 `reflect` 逐枚量：`0->1`/`1->2`/`2->3` 全 false、`3->4` true）；**两处记录**都在（日志 INFO/WARN ＋ `Store.ForeignKeyAudit()`）。**验收判据用的是四条路径的读数**（不是"函数被调了"）。

## ③ 四条孤儿夹具：**我自己造的读数**（`logs/31-post-probe.txt`）

| 夹具 | `Open()` err | `Checked` | v | **`OrphanRows`** | `child_tables` | 迁移前 fkc | **迁移后独立 fkc** | task_log/tool_call（前→后） | `task_id` 逐字 |
|---|---|---|---|---|---|---|---|---|---|
| `purge-tasklogs`（`PurgeTaskLogs`） | **`<nil>`** | true | 4 | **3** | `[tool_call]` | 3 | **3** | 0/3 → 0/3 | `task-A,task-A,task-B` |
| `delete-tasklog`（`DeleteTaskLog`） | **`<nil>`** | true | 4 | **2** | `[tool_call]` | 2 | **2** | 1/3 → 1/3 | 同上 |
| `manual-insert`（raw，绕 DAO） | **`<nil>`** | true | 4 | **1** | `[tool_call]` | 1 | **1** | 2/4 → 2/4 | `…,ghost-missing-task` |
| `retention`（`RunRetentionOnce`） | **`<nil>`** | true | 4 | **1** | `[tool_call]` | 1 | **1** | 2/4 → 2/4 | `…,task-old` |
| `clean-v3`（无孤儿正控） | **`<nil>`** | true | 4 | **0** | `[]` | 0 | **0** | 2/3 → 2/3 | `task-A,task-A,task-B` |
| `clean-fresh`（0→4 全新） | **`<nil>`** | true | 4 | **0** | `[]` | — | **0** | 0/0 | — |

## ④ ★"永久锁死不复发"的分水岭读数（同一条孤儿库连开两次；`logs/31`＋`logs/44`）

- `purge`／`delete`／`manual`／`retention` **四形 + `clean-fresh`／`clean-v3` 两形，`OPEN#2` 全部 `err=<nil>`**（`logs/48` 单量了两种干净形）；第二次 `Checked=**false**`、`SchemaVersion=4`、`OrphanRows=0`、`ChildTables=[]`；**孤儿仍在**（独立 `fkc` 仍 **3/2/1/1**）、行数与 `task_id` 一字未动。
- **`Checked=false` 是如实的（不是"跑了得 0"）**：日志探针（`logs/44`）—— `OPEN#1`：`schema migration step applied`＝**true**、`found dangling references`＝true、`orphan_rows=3`、`child_tables=tool_call`、`opening anyway`、**无 `level=ERROR`**；`OPEN#2`：`step applied`＝**false**、**`migration baseline:` 一行都没打**（⇒ 既不谎报"干净"，也不谎报"跑过"）。机理＝`open.go:404 version == s.target` 早退。
- 与硬门槛的分水岭（`M2` 实测，见 ⑥）：门槛形下**三次连开三次被拒、水位停在 v3** ⇒ 那才是"永久打不开"。

## ⑤ ★计数如实＋行数未变（我自己核的）

- 计数 3/2/1/1 与**独立业务连接**的 `fkc` **逐条相等**（两把尺同值，⛔ 一把尺自证）；无孤儿库＝**0**（fresh＋clean-v3 两形）。
- **行数一行未变**：四条夹具迁移前→后 `task_log`/`tool_call` 与 `tool_call.task_id` 逐字相同（③）⇒ **⛔ 任何"删孤儿"动作**。
- **基线复跑（改前 `4ee57d318` 导出树，`logs/32-pre-probe.txt`）**：`0→0`／`3→3`／`2→2`／`1→1`／`1→1` ⇒ 与 `366-a1` 的 `3/2/1/1` 对得上，且**迁移既不清孤儿也不造孤儿**。

## ⑥ ★`M2`＝把记录改成门槛（**我自己造的那一发**，仓外 `muttree-m2`，+4 行）

- 变异体：`audit = auditForeignKeys(...)` 之后 `if audit.OrphanRows > 0 { tx.Rollback(); return audit, fmt.Errorf(...) }`。
- **红**：`rc=1`；366 的 **3 枚顶层用例**（`RecordsOrphans…`×4 子例、`NeverLocks…`、`IsLogged/orphans`）＋本腿探针共 **4 枚顶层／5 子例**；红行 `migration_fkcheck_366_test.go:271×4,329,392`；报错串 **`Open refused a database with 3/2/1/1 orphan row(s)`**（`logs/40-m2-test.txt`）⇒ **"⛔ 拒开库"这条边界有牙，且牙咬的是硬边界 #1**。
- **附带的锁死实证**（`logs/46`）：同一发 `purge` 库**连开三次、三次都 `err=…opening refused…`、`schema_version` 三次都停在 3**、行数原样 ⇒ 门槛形＝永久打不开，与交付形的"两次都 `<nil>`"正相反。
- 旁证 `M1`＝摘掉那一跳（−1 行）：`rc=1`、**4 枚顶层／8 子例**红，红行 `:207,232,278×4,333,380,383,403×4`（与 r1 报的行**逐枚相同**）。

## ⑦ 边界零 diff（逐项；sha256 都是我**在两棵树上各算一次**）

- **`ddlV1`–`ddlV4`**：`ddlV1 fda736e1…`（len 4463／15 句）· `ddlV2 3e9ece70…`（903／1）· `ddlV3 d46ee18e…`（1286／5）· **`ddlV4 c17241cd…`（1357／5）** —— `pretree` 与 `headtree` **四条逐字相同**。
- **12 枚名册**：两树同＝`running,done,succeeded,cancelled,interrupted,stuck,error,Thinking,Settling,Muted,Stuck,Error`。
- **四个边界函数 body 逐字同**（`probes/boundary-funcs.sh` 抽 `func` 到首个顶格 `}`）：`migrationDSN 8b40de05…`(3 行)／`migrationPragmas 97da5dd4…`(10)／`dsnPragmas 2e99efb5…`(8)／`requireForeignKeysOff 8050de10…`(25)；`git diff 6d47e1ff4..HEAD` 里**唯一**提到 `migrationDSN` 的新增行是**注释**（`schema.go:515`）。
- **⛔ DELETE**：`git diff -U0 6d47e1ff4 HEAD -- internal/memory` 的**新增行**里 `DELETE` **0 命中**；新用例文件 `grep -c DELETE`＝**0**。
- `numstat` 与它报的一致：`open.go +113/−20`·`schema.go +122/−1`·`taskstate_196r2_test.go +2/−2`·新用例 `+426`。
- **工作树未当靶**（还原核时机＝本腿最后一发变异之后量）：`git status --porcelain -- internal/memory` **空**；四枚文件 `sha256(工作树) = sha256(headtree) = sha256(git show HEAD:)`（`logs/47`）；三份变异副本各自**只差自己那一处**（M1 −1 行／M2 +4／bad +3），**其余产码文件全同**。

## ⑧ ★"坏尺不变成锁"：**行为上做到了**（注入实测，`muttree-bad`，`logs/42`）

- 注入：让 `foreignKeyCheck` 恒返回错误（＝PRAGMA 读不出来），其余交付字节不动。读数：`purge` 库 `err=**<nil>**`、`audit checked=**false**`、`version=4`、孤儿仍 3；日志含 `could not be read`＋`opening anyway`＋`level=WARN`、**无 `level=ERROR`**；**第二次开也 `<nil>`**。`clean` 库：`err=<nil>`，且**没有**打 `found no dangling references`（⛔ 把"没读到"伪装成"读到 0"）。
- ⇒ 代码路径＝`auditForeignKeys` 的 `if err != nil { Warn; return audit }`，**永不失败**，符合形的要求。**但**这一支**没有用例钉住**（⑨）。

## ⑨ ★那条"无牙"支：**⛔ 必须先补，不能翻**

- 理由（四条，逐条是尺）：①AC#2 要求"**每处改动**配一枚今天响、不修不响的用例"，而"坏尺⇒照常开库"是**本腿新增的判断**，不是既有行为；②它是本形里**唯一**紧贴硬边界 #1 的分支——把它写成 `return err` 就立刻变成 ⑧ 里那台**永久锁**；③本仓反复栽在"没有牙的护栏"这一族；④**现值在好尺下永不执行** ⇒ 任何回归在正常跑里**全绿**，只在现场坏尺那一发把用户锁死。
- **实测凭什么说"现有用例不算牙"**：把尺打坏后，**交付的 5 枚用例 4 红 1 绿**（`logs/43`）——但红的原因是它们要求 `Checked=true`，**信号饱和**：正确实现（照常开库）与"拒绝开库"都会红，区分不出来。
- 现成接缝（**⛔ 需改产码**，接缝已在）：`txExec` 是接口 ⇒ `type t366FailingTx struct{}`（`ExecContext`/`QueryContext` 各返回 error），调 `auditForeignKeys(ctx, t366FailingTx{}, SchemaVersionTarget, logger)`；建议名 `TestMigration366AuditNeverTurnsABrokenRulerIntoALock`；断言 5 条：不 panic · `Checked=false` · `OrphanRows=0`＋`ChildTables` 空 · 日志有 `WARN`＋`could not be read`＋`opening anyway`＋`schema_version` · 无 `level=ERROR`。

## ⑩ ★措辞那一格：**没做到**（8 处新增注释仍把"官方第 10 步"当已核事实）

- 处所（全是 r1 新增行）：`schema.go:197-198`（`this is its step 10`）· `:303`（`Step 10 of that procedure - PRAGMA foreign_key_check -`）· `:509-510`（`the step 10 of the documented 12-step ALTER TABLE procedure`）· `:518`（`the migration chain's step-10 record`）· `open.go:331-332`（`the 12-step ALTER procedure's step 10`）· `open.go:513-514`（`its step 10`）· 用例文件 `:3`、`:424`。件标题同口径；commit subject 也写"（12 步第 10 步）"，但受 ⛔ amend 约束、**永久留旧口径**（以本件为更正）。
- 它的 §1 第 24 行只在**件里**写了"未当已核事实"，**代码里没有这一句** ⇒ 编排者拍的改法**未落地**。
- **具名新措辞**（en，代码注释总句）：`An extra referential-integrity guard around the v4 rebuild. SQLite's 12-step ALTER TABLE procedure may run the same check at its step 10, but this repository has not verified that text, and that step is conditional on enforcement having been ON - which it never is here (migrationDSN turns it OFF). Extra insurance that may coincide with step 10; not a step this code claims to follow.`
  **逐处替换**：`:197-198` → `this is an extra check at the point where that procedure may check integrity (unverified here)`；`:303` → `The integrity check (possibly step 10 of that procedure; unverified here) - PRAGMA foreign_key_check -`；`:509-510` → `an extra guard for the v4 rebuild (possibly the same place as step 10 of the documented procedure; unverified here)`；`:518` → `the migration chain's post-rebuild integrity record`；`open.go:331-332` → `(extra insurance that may coincide with step 10 of the 12-step ALTER procedure; unverified here; ticket 366)`；`open.go:513-514` → `An extra integrity check after the rebuild and before the commit (where the 12-step procedure may put one; that step is conditional on enforcement having been ON, and is unverified here).`；用例 `:3`／`:424` 同口径。
  **中文（件/票面）**：`额外保险（可能与官方 12 步的第 10 步重合；官方原文本机取不回、未核；且官方那一步以"原本开着 FK"为前提，本仓迁移连接一直是 OFF）`。
- ⚠ **三处改不动的旧措辞**（不属 r1、且出本票边界）：`schema.go:324` **在 `ddlV4` 字符串里**（动它就破 ⑦ 的 sha256）· `open.go:184-185` 与 `schema.go:406-407/:442`（`196-r3` 的 step 1 措辞）⇒ 要统一**另立票**。

## ⑪ 门禁 rc（逐门分开、⛔ 相减、⛔ 走管道取 rc）＋格式两尺

| 门 | 命令 | 读数 |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **rc=0**（0 字节输出） |
| G2 | `GOFLAGS= go vet ./internal/memory/` | **rc=0**（0 字节输出） |
| G3 | `GOFLAGS= go test -count=1 ./internal/memory/` | **rc=0**；四数：`ok 25.401s`／`=== RUN`=**114**／`--- PASS`(顶层)=**64**／`--- FAIL`=**0**／`--- SKIP`=**2** |
| G4 | `sh scripts/d22scan.sh` | **rc=0**（`clean - no D22 ban violations`；ban #8 internal/ 639 枚） |
| G5 | `sh scripts/check-path-length-budget.sh` | **rc=0**（`VERDICT GREEN`；tracked 17670） |

- **格式甲（入库 blob＝判据）**：四枚 `git show HEAD:<p>` 导出到仓外再 `gofumpt -l` ⇒ **0 命中**。**乙（工作树＝只作诊断）**：`gofumpt -l internal/memory` 经 `[\\/]` 过滤 ⇒ **0 行**。
- **CRLF 尺**（⛔ 用 `gofumpt -l`）：`git show HEAD:<p> | tr -cd "\r" | wc -c` ⇒ 四枚全 **0**（LF）。
- 射程：G1 在**共享工作树**上跑，里面有**别的腿**在飞的 `internal/panel`／`internal/perm` 4 枚 `M` ⇒ G1 量的是"含别人改动的整模块"；`internal/memory` 本身与 HEAD 逐字同（⑦）。

## ⑫ 提交号＋逐笔越界

- 本件提交（`git log -1` 见提交后回填）；`add`/`commit` 两处**显式 pathspec**，只点 `docs/evidence/s1/366-v1-verdict.md` 一枚；提交前 `git diff --cached --name-status` 只见 **`A`** 这一行。⛔ push。
- 越界：仓内 `.go`＝**0 枚**；`internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`＝**0 字节改动**；临时件全在 `D:\tmp\wisp\20261023\366v1\**`（另有首次 `git archive` 误落到 `366v1\` 根的一份导出树，未参与任何读数）。

---

## 附：上报表（不占判语）

1. **还原核 1/4 枚不可复现**：r1 §6 报 `open.go 3d9f2655…`，入库 blob 的 sha256＝**`d13000b1…`**（工作树／`headtree`／`git show HEAD:` 三处同值）；其余三枚（`schema.go 3f2e62ea…`／`taskstate_196r2 a34a2591…`／新用例 `8730ae29…`）与它报的**逐字相同**。⇒ 它那句"posttree 与工作树逐字相同"至少对 `open.go` 不成立；也不是撞在变异产物上（M1＝`dde55d69…`、M2＝`977b6bdd…`）。**要求**：重量并给 **64 位全文**（它报的是 8 位截断，按尺清单第 9 问不可复核）。
2. **反向判据的子例数两处对不上**（行号与红项我全复现）：`M1` 它报 10 子例、我量 **8**；`M2` 它报 6 子例、我量 **5**（它把一枚**顶层**算进子例）⇒ 尺清单第 4 问"这个数有没有和总数对上"。
3. `366-a1` §10-1 的**官方第 10 步逐字条件句**本腿也取不回（`web_fetch` 同败）⇒ 本裁**不**建立在"官方条件句"上；⑩ 的新措辞正是为这个未核状态写的。
4. 未定案（⛔ 本腿自行扩）：`ForeignKeyAudit()` 返回结构体拷贝，但 `ChildTables` **与 `Store.fkAudit` 共用底层数组** ⇒ 调用方改元素会污染后续读数。今天零调用方，**不阻塞**；建议随"补牙"那笔在访问器里 `slices.Clone`。
