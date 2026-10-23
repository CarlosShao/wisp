# 票 196 终裁腿 `196-v1` —— `196-r1` 的〔乙〕落地（甲＝权威映射 ＋ 乙＝schema v4 CHECK）够不够格结票

- **腿**：`196-v1`（**非实现者／裁决者**，裁决者≠实现者）。被裁对象＝`196-r1`（`6e19f4a43`，8 文件 `959+/40−`，件 `docs/evidence/s1/196-r1-impl.md` 139 行；我逐笔核过 `--numstat` 与件行数，均对得上）。
- **本票是"资料/裁决记录型"、无判据框** ⇒ 我判的是**这次落地的质量**，⛔ 逐格勾。
- **落点先 `ls`**：`docs/evidence/s1/196-v1-verdict.md` **不存在**（`docs/evidence/s1/196*` 只有 `196-r1-impl.md`）⇒ **未改号**。
- **起手 HEAD `96f1af9c9`、收手 HEAD `1a02c3fca`**：本腿期间共享树**落了 15 笔**（`git rev-list --count 96f1af9c9..HEAD`），我先后看到 96f1af9c9→47987d323→b1d1edd3c→b22225513→50f97c1b7→1a02c3fca 六枚快照。★**但 `git log 6e19f4a43..HEAD -- internal/memory internal/statemachine` 收手时仍为 0 笔** ⇒ 我量的是 `196-r1` 的字节，⛔ 别人的活。
- **环境**：`go1.27.1 windows/amd64`、`GOFLAGS=`、`unset GOTMPDIR`；`CGO_ENABLED` 未动。★**`TMP/TEMP` 我留着**——理由见 §7 末（照派单 `unset TMPDIR TEMP TMP` 会造一枚**假红**）。
- **仪器**：①仓外导出树三枚（`before`＝`6e19f4a43^`、`landed`＝`6e19f4a43`、`pristine`＝`landed` 去我探针的再导出）；②`-overlay` 自造突变体三枚；③**我自己写的探针** `zz_196_v1_probe_test.go`（落在**仓外**导出树、⛔ 抄它的 fixture、⛔ 落任何仓内 `.go`）；④`git archive` 甲尺。
- ★**仓内 `.go` 零字节、工作树零突变**：本腿唯一写入＝本件（`git status --porcelain -- internal cmd scripts` 首尾各一次，只有别人的件）。
- ★**本件 §1–§9 与必判清单 ①–⑨ 逐项一一对应**（§10＝裁决）。
- **裁语**：**够格结票（能结）**，附 §10 五条具名残余（**无一阻断**）。★★**最要紧的一条（老库能打开）过**。

## §1 票面前提新鲜度 —— 它报「票面自己腐坏 3/4」，我逐条复核

| # | 票面那条 | 我在 **`6e19f4a43`** 上量到（`git show <sha>:<p> \| sed -n`） | 判 |
|---|---|---|---|
| F1 | `states.go:8` ／ `:11-31` | `:8` ＝ `// The 20 states, names exactly as in D43 / SPEC-08…`；`:11-31` ＝ 20 枚 `StateXxx State = "…"` ＋收尾 `)` | **不漂** |
| F2 | `PLAN.md:3055` ＝ D43 标题 | `:3055` ＝ **D42 表第 11 行（WebView2 缺失）**；D43 标题在 **`:3061`** | **漂 +6，它说得对** |
| F3 | `loop.go:983-998`／`run.go:651-670` | `loop.go:983-998` ＝ `func (l *Loop) finish(` 的签名与赋值；`run.go:651-670` ＝ `tools.Options` 注释块；真位置＝`loop.go:1046-1063`／`run.go:1361-1372`（两处都是 `taskLogState`） | **漂，它说得对** |
| F4 | `models.go:57` 邻域 | `models.go:57` ＝ `// complete.` | **漂，它说得对** |

★★**但它给的替代锚点，两条在交付件落定时已经不准**（我判这是本件唯一实质缺陷）：
- ★**F2 的 `:3061` 今天再漂 +4 → `:3065`**：`196-c1`（`47987d323`，机主批的契约文字腿）在 D35 插了 4 行。**⛔ 它的错**（它量时是对的），但**件里的号现在会扑空**。
- ★★**F4 的 `schema.go:53-66`／“`state TEXT NOT NULL` 在 `:57`” 是自指腐坏**：落地前 `schema.go:53-66` 恰是 `task_log` 块，`:57` 恰是那行 `state TEXT NOT NULL`（我现跑核过：`6e19f4a43^` 的 `:57` 逐字就是它）；**但它自己那一笔把 `ddlV1` 推后 11 行** ⇒ 落地后 `state` 在 **`:68`**、块在 **`:64-77`**。⇒ ★**它用"引用会腐"证明票面腐坏，同时自己写了一枚会腐的引用**。

## §2 ★★核心判点一：改前「三扇门全开」—— 我自己造那三发（`before` 导出树 ＝ `6e19f4a43^`）

探针＝`zz_196_hole_probe_test.go`（我写、仓外、跑完留档）。三扇门逐字：

```
READING door1 StartTaskLog(state="definitely-not-a-task-state") err = <nil>
READING door2 FinishTaskLog(state="definitely-not-a-task-state", errorClass="") err = <nil>
READING door3 raw INSERT(state="definitely-not-a-task-state") err = <nil>
READING rows now holding "definitely-not-a-task-state" = 3
READING schema_version = 3
READING task_log DDL:  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
```

★**三发全 `<nil>`、3 行垃圾真落库、列上只有 `NOT NULL`** ⇒ **成立**。★门 2 比票面写的更宽这条也成立：`FinishTaskLog` 在 `errorClass == ""` 时**一次都没校验过**（原实现 `if errorClass != "" { … }` 包着校验）。

## §3 ★★核心判点二：CHECK 是不是真拦住（`landed` 导出树，我自己的探针）

```
door1 StartTaskLog(garbage)            err = memory: invalid task_log.state "…" (want one of: running|done|succeeded|cancelled|interrupted|stuck|error|Thinking|Settling|Muted|Stuck|Error)
door2 FinishTaskLog(garbage, "")       err = 同上（干净退出那条路已关）
door2b FinishTaskLog(garbage, class)   err = 同上（带 error class 那条路仍关）
door3 raw INSERT(garbage)              err = constraint failed: CHECK constraint failed: state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error') (275)
door3b raw UPDATE(garbage)             err = 同上
rows holding garbage = 0
legacy word "running"/"done"/"cancelled"/"succeeded"   DAO err=<nil> raw err=<nil>   （四词 × 两嘴，8/8 过）
整册 12 枚 raw SQL                      逐枚过（0 报错）
DDL 列逐字: state TEXT NOT NULL CHECK (state IN ('running', 'done', 'succeeded', 'cancelled', 'interrupted', 'stuck', 'error', 'Thinking', 'Settling', 'Muted', 'Stuck', 'Error')),
```

⇒ **两向都成立**：非法被拒（DAO 与 raw SQL 两条路、INSERT 与 UPDATE 两向）、四词各一过。★它报的 `(275)` 与逐字 CHECK 文本，我**逐字节**复现。

## §4 ★★核心判点三：老库能不能打开 —— **本票最要紧的一条**

夹具我自己造：用**生产 `ddlV1`+`ddlV2`+`ddlV3` 常量**建真 v3 库（`schema_version=3`、**落地前 `sqlite_master` 里无 CHECK**），塞 **12 行**＝旧词 7 枚 ＋ D43 名 5 枚，**再塞 1 枚 `tool_call` 行**（外键指向 `task_log`，它那份夹具没塞）。

```
fixture schema_version = 3        fixture has CHECK before Open = false
after Open: schema_version = 4    after Open: rows = 12 (want 12)
states read back = [Error Muted Settling Stuck Thinking cancelled done error interrupted running stuck succeeded]
idx_task_log_started present = 1  leftover task_log_v4 objects = 0
raw INSERT illegal AFTER migration = constraint failed: CHECK constraint failed: … (275)
tool_call rows after migration = 1
```

★**开得开、一枚不少、逐词读回、索引重建、旧表无残留、迁移后垃圾被拒、连带 `tool_call` 行活着**。
**反向**（本来就坏的库）：塞 `bogus-state-written-by-a-fork` ⇒ `Open()` 失败并**点名**：

```
memory: migration 3->4 failed (database restored from backup is at …\backup\wisp.db.bak-3-4):
apply: task_log.state holds 1 value(s) outside the v4 roster (bogus-state-written-by-a-fork); refusing to rebuild …
after refused migration: schema_version = 3   rows = 2
```

⇒ 「**CHECK 若把老库锁死 ⇒ 用户升级即炸**」这条**没发生**：v3 库（含四词）打开成功。「打不开」只发生在**库里已有名册外词**时，且**先点名、后回滚、有备份**——比让 `INSERT…SELECT` 撞一句不带行身份的 `CHECK constraint failed` 强。★**我判这是〔乙〕能做到的最优形态**，非阻断。

## §5 名册覆盖度 —— 我核的是"生产真写几词 ＋ 旧库可能几词"

- **今天生产写 6 词**（我普查 `INSERT/UPDATE task_log` 的非测试调用点只有 `dao_tasklog.go` 两处，喂它的只有三处）：`running`（`journal.go:131`）＋ `done/stuck/cancelled/interrupted/error`（`loop.go:1046-1063`）＋ `done/cancelled/error`（`run.go:1361-1372`）。
- ★**旧库可能的词**——我跑了**历史普查**（⛔ 只看今天）：`git log --all` 里碰过 `loop.go`/`run.go`/`journal.go` 的 **53 笔**提交，逐笔抽 `func taskLogState` 的 `return "…"`，**去重只有 5 枚** `done`/`error`/`cancelled`/`stuck`/`interrupted`；`journal.go` 历史上只写过 `State: "running"`。⇒ ★**除 `succeeded` 外没有第二枚"旧词"**，而 `succeeded` 今天只在测试（`dao_test.go`、`privacy_test.go`）与旧行里。
- **名册 12 枚＝旧词 7（6 生产 ＋ `succeeded`）＋ D43 名 5**，**无 `Sleeping`/`Listening`/`Acting`/`Speaking`** 等球态（我逐枚对 `internal/statemachine/states.go:11-31` 的 20 枚核过，5 枚全在其中，另 15 枚一枚没进）。
- **判语**：**覆盖够、没多收**。★另核：票面"今天的四词一枚都不在 D43"**成立**——`running`/`done`/`succeeded`/`cancelled` 对 `states.go:11-31` 的 20 枚**逐枚不在**（大小写敏感、⛔ 近似匹配）。★两条残余（非阻断）：①5 枚 D43 词**今天零生产者**（为将来迁移预留，它自己 §9.3 已点名）；②反向映射**多对一有损**（`succeeded`→`Settling`→`done`），`TaskStateFromD43` 取名册靠前者 ⇒ 往返确定但**不保原拼写**（它写了、也是刻意的）。

## §6 反向判据 —— 我自己造（⛔ 只看它的 M1/M2/M3）

载体＝`-overlay` 把 `landed` 树里的单文件换成突变体（仓内零改动）。三发全部 `rc=1`：

| 我的突变 | 摘掉哪一跳 | 红数 | 逐名报红（件:行） |
|---|---|---|---|
| **A** | CHECK 改成恒真 `CHECK (1 = 1 OR state IN (…))` | **4** | `taskstate_196_test.go:196`（`raw INSERT … succeeded: the column still has no CHECK`）、`:374`（`the migrated database still accepts an illegal state`）、我探针 `zz…:106/115`（door3/door3b 全 `<nil>`） |
| **B** | DAO 名册判定 `if false && !ValidTaskState(...)` | **2** | `taskstate_196_test.go:147`／`:169`——★逐字是 `refusal "…CHECK constraint failed…" does not quote the offending value` ⇒ ★**关掉甲，乙先开口**（纵深成立，它的 M2 说法我复现） |
| **C** | 名册摘掉 `TaskStateSucceeded` | **9** | 既有 3 枚：`dao_test.go:266`、`privacy_test.go:32`、`retention_test.go:123`；196 的 4 枚：`taskstate_196_test.go:53`／`:114`／`:232`／`:335`；我探针 2 枚 |

★**它的三个读数（2／2／7）与我（4／2／9）逐项自洽**：差额正是我多写的两枚探针（A +2、C +2、B +0）⇒ **它的 M1/M2/M3 可信**。
★**一条它没说的**：A 那发里 `TestTaskState196CheckIsDeclaredOnTheColumn`（只验 DDL 文本含 `CHECK`）**仍绿** ⇒ **"DDL 里有 CHECK"这枚测试单独挡不住恒真 CHECK**，承重的是行为那枚；它两枚都有，所以没漏。
★**还原核**：`sha256sum -c` 三枚 `OK`（`sha256sum`＝SHA-256／64 位；⛔ 与 `git hash-object` 的 SHA-1／40 位混比——我先问了口径再比）。

## §7 门禁 rc（四数分开报、⛔ 相减）＋ gofumpt 口径

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 空 |
| `go vet ./internal/memory/ ./internal/statemachine/` | **0** | 空（`go vet ./...` 我也跑了一发 **0**） |
| `go test -count=1 ./internal/memory/ ./internal/statemachine/` | **0** | `ok …/memory 24.109s`、`ok …/statemachine 0.351s` |
| `sh scripts/d22scan.sh` | **0** | `clean`；ban #8 扫 internal 630／cmd 148 Go 文件 |
| `sh scripts/check-path-length-budget.sh` | **0** | `tracked paths=17625  over-budget=61  covered by roster=61  not in roster=0` ⇒ `VERDICT GREEN`。★与 `196-c1` 那天看到的 `62／61` 不同，因为票 365 那枚超长件已在 `a957dfdc0` 改名（我这一发在改名后） |
| **gofumpt 甲尺**（`git archive 6e19f4a43` → 仓外 `pristine`，**判据**） | **0**，**输出 0 行** | ★第一次跑出 1 行，是我自己留在 `landed` 树里的探针 ⇒ **重导 pristine 再跑才 0 行**（"谁的字节"这一问）。当班 HEAD `50f97c1b7` 的导出树同样 **0 行** |
| **gofumpt 乙尺**（工作树，**只作诊断**） | **0**，**输出 4 行** | 3 枚整文件 CRLF（见 §9④）＋1 枚**真脏** `internal/proc/envfork_test.go`（CR=0、`M`＝别人未提交的在飞件；同文件在甲尺导出树里干净） |
| `gofumpt -l internal/memory`（工作树） | **0**，**0 行** | `196-r1` 的 8 枚文件 CR 全 0 |

★**归属判定口径**：`gofumpt -l` 在 Windows 上吐**反斜杠**（`internal\memory\…`），我一律先把 `\` 归一成 `/` 再分类，⛔ 按字面前缀 `grep '^internal/'` 直接判（那会把 4 枚全判成"不在 internal"）。

★★**它更正的那个口径，我对过了、它是对的**：`gofumpt -l` **列文件时自身 rc 仍 0**——我用正控实证（一枚故意写坏的 `dirty.go`：`rc=0`、输出 1 行）；`.github/workflows/ci.yml:217-226` 的判据就是 `OUT="$(gofumpt -l …)"; if [ -n "$OUT" ]; then … exit 1; fi` ⇒ **CI 判"输出非空"、⛔ 判 gofumpt 的 rc**。★这是"谁算的、什么口径"那一族，它踩对了。
★★**一条台面陷阱（我记、⛔ 算它的）**：照派单字面 `unset GOTMPDIR TMPDIR TEMP TMP` 跑，`os.TempDir()` 变**空串**（我实测），`internal/memory` 的 `TestArtifactsContainmentByDirectoryListing` **假红**（`legit delete removed [data/artifacts/real-1.txt], want exactly [D:\…]`）；我在 **`6e19f4a43^`** 上复跑**同样红** ⇒ **与 196 无关、纯环境**。故本腿用的口径是 **只 `unset GOTMPDIR`**（`GOTMPDIR` 本就空，d22scan 未触 311 守卫）。

## §8 新增红 0 ／ 写面纪律

- **逐笔 `git show --name-status`**：`6e19f4a43` ＝ 8 枚（`docs/evidence/s1/196-r1-impl.md` ＋ `internal/memory/` 的 7 枚）；补两笔 `e89d1fe98`／`2c3628679` 均**只动件**。★`grep '\.go$' | grep -v '^internal/memory/'` ⇒ **空**；`grep -v '^internal/memory/' | grep -v 196-r1-impl.md` ⇒ **空** ⇒ **零越界**。
- **正控（同尺打一笔含 `.go` 的真提交）**：`5ffa553e3`（348-r2）→ 尺列出 `internal/winsec/deleteseam_348r2_windows_test.go`、`internal/winsec/winsec.go` 两枚 `.go` ⇒ ★**这把尺不是瞎的**。
- **新增红 0**：`internal/memory`＋`internal/statemachine` 全绿（§7）；另跑受 DAO 改动影响的三个包（HEAD `1a02c3fca`，`rc=0`）：`ok …/internal/agent 4.666s`、`ok …/internal/session 4.346s`、`ok …/internal/tools 19.148s`。★**我跑过的六个包里唯一的红**＝`cmd/wisp` 的 `0xc0000135`，**归因见 §9③：预存、与本次无关**（`internal/observe` 等我没跑，⛔ 越报）。
- **测试没被放宽**：`schema_test.go`／`migrate_v1seed_test.go` 的 `-` 行只有"3→4 水位／备份名加 `wisp.db.bak-3-4`"，**且新增** `splitSQLStatements(ddlV4) != 5` 与 `bak-3-4` 两枚断言 ⇒ **改严了、⛔ 改松**。

## §9 写面外 4 条 —— 我的判语（⛔ 算它的）

① **第二个 `taskLogState`**：**真**。`loop.go:1046-1063` 把 `StatusStuck`→`"stuck"`；`run.go:1361-1372` 无该 case ⇒ 落 `default` → **`"error"`**。⇒ 同一状态两枚函数各说各话，且 `cmd/**` 不在它写面 ⇒ **它不动是对的**（编排者已立票 365）。★补一句：DAO 收口后**库里**只会剩名册词，但**两个函数仍各写各的**，收敛属 365。
② **契约文字欠账**：★**"已摆机主"这个处置对**——当场改 `SPEC-02`/`D35` 是"改契约文字"，须另走批准（`AGENTS.md` 铁律第 4 条三类之一），它只**具名**、⛔ 顺手改。★且**已闭环**：机主批「1. 可以」⇒ `196-c1`（`47987d323`）把 12 枚写进 `SPEC-02 §3` **围栏外** ＋ `D35`。★我核了它那句陈述**今天仍逐字真**：`docs/specs/SPEC-02-data-storage.md:60` 仍是 `state TEXT NOT NULL,` ⇒ **围栏内确实仍无 CHECK**。
③ **`go test ./cmd/wisp/` 预存红 `0xc0000135`**：**归因对**。我在 **`6e19f4a43^`** 导出树跑 ⇒ `rc=1`、`exit status 0xc0000135`（`0.069s`）；当班工作树同码同错；★再打一发 `-run TestNothingMatchesThis`（**一枚测试都不匹配**）⇒ **同样 `0xc0000135`** ⇒ **进程在加载期就没起来**（DLL 缺失），与 196 无关。
④ **3 枚工作树 CRLF**：**判"⛔ 归一"对**。逐枚 `tr -cd "\r" | wc -c`＝**行数**（`pending_read.go` 131/131、`tools.go` 225/225、`cmd/wisp/models.go` 389/389）；`git status --porcelain` 对三枚**无输出**；`git ls-files --eol` ＝ `i/lf w/crlf`（`.gitattributes` 写死 `*.go text eol=lf`）；仓外导出树与 `git show HEAD:` 的 blob **CR=0** ⇒ ★**这是工作树字节、⛔ 提交内容**，写它只会给共享树添一笔三方文件的无主 diff。★**新增一条**：乙尺今天多出第 4 枚 `internal/proc/envfork_test.go`（`CR=0`、`M` 未提交、别人的在飞件）⇒ 乙尺今天＝**3 假 ＋ 1 真**，更说明它只配当诊断。

## §10 裁决：够不够格结票

★★**能结（够格）**。判据四条，逐条都有我自己的读数撑着：①**改前的洞是真的**（§2 三发 `<nil>`、3 行垃圾落库）；②**改后两侧都关**（§3 DAO＋raw SQL、INSERT＋UPDATE，四词 8/8 过）；③★★**老库开得开**（§4：12 行 v3 库 → v4、一枚不少、逐词读回、带 `tool_call` 行也活）；④**名册覆盖对**（§5：生产 6 词 ＋ 历史无第 7 枚 ＋ 5 枚 D43 名，没多收球态）。加：反向判据三发全红且**它的数我复现**；门禁五门全 **0**；写面零越界。

**五条具名残余（⛔ 阻断结票，但请编排者知悉）**
1. ★**F4 的替代锚点在交付件里自腐**：`:57`→落地后 `:68`、块 `:53-66`→`:64-77`（它自己那一笔 +11 行）。按件复核会扑空——**这是本件唯一的实质瑕疵**，属"引用会腐"的同一种病。
2. **F2 的 `:3061` 已被 `196-c1` 再推 +4 到 `:3065`**（⛔ 它的错，但件里的号今天不指向 D43）。
3. **读侧 API 今天零调用者**：生产代码里 `TaskStateFor`/`D43State`/`TaskStateFromD43` **0 处**（只有 `taskstate.go` 自己与 `schema.go` 的 `TaskStateRoster`）⇒ "读走它"目前是**约定**、⛔ 既成事实；5 枚 D43 词也**零生产者**。
4. **`applyV4` 的 `DROP TABLE` 今天靠 `foreign_keys` 关着才安全**（`dsn()` 不设该 pragma）：若将来打开 FK 且库里有 `tool_call` 行，重建会在 DROP 处失败（我这一发塞了 `tool_call` 行是绿的，正是因为这个 pragma 关着）。属**将来开关**的射程，非本笔缺陷。
5. **台面**：照派单字面 `unset TMPDIR TEMP TMP` 会在 `internal/memory` 造一枚**假红**（§7 末），与本笔无关；后续腿照抄会把假红读成新红。

★一句话：**〔乙〕的形状是对的（名册一处权威、SQL 由 Go 生成、v4 重建走同一条 3→4、先点名再动手），读数是真的，反向判据是真的，老库是真能打开 —— 这票可以结。** 它交的 4 条写面外事项，处置**全部正确**，其中 ② 已由 `196-c1` 闭环、① 已由票 365 接住。
