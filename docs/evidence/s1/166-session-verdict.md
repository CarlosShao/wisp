# ticket 166 · Go 侧会话管理层（列／搜／改名／归档带撤销／分叉）· 六格裁决表

- 腿＝**`166-v1`**（验收腿·裁决者·⛔ 实现者·只读产码）｜时刻＝2026-10-16 09:5x–10:4x +08
- 票＝`.scratch/wisp/issues/166-no-session-layer-list-search-rename-archive-fork-on-the-go-side.md`（六枚 AC 框**本腿一枚未翻**）
- 被裁的字节＝`427c62d0a`（166-r1，`internal/session/manage.go`＋`internal/memory/dao_session.go`＋`schema.go`＋四件测试件，9 路径 +975/−40）
  ·`24e36406`·`e1192509`（166-r1 注释与读数件）∥`c2b1ade5`＋`464f5549`（166-r2 仪器＋证据件＋读数件，编排者代提）
- 本腿 HEAD＝`feb7afae`；`git diff --stat 46007856 HEAD -- internal/session internal/memory`＝**0 行** ⇒ 派单锚与本腿读的字节同一枚
- 未修码基线＝`2f41a80cb`（＝`427c62d0a^`）
- ★这单的起点：**166-r1 在盘上没有证据件**（`ls docs/evidence/s1/ | grep 166` 起手只命中 `ticket166-ac6-166r2.md` 一枚，属另一条腿）
  ⇒ AC#1–AC#4 的读数全部由本腿自己跑，`probes/166/r1/` 那两枚过程件**只当索引、⛔ 当凭据**。
- 读数目录＝`.scratch/wisp/probes/166/v1/`（19 枚 logs ＋ 7 枚读数码 ＋ 3 份台件套件源码），本表每一格都指向那里的具体文件。

## 判语一览

| 格 | 判语 | 一句话凭据（详数在指向的件里） |
|---|---|---|
| **AC#1** 未修码读数 | **成立** | 命令档 `wisp sessions/history` 都 rc=2 `unknown command`；API 档基线五枚出口 `git grep` rc=1（正控 HEAD 5 命中）；`go doc Store` 基线唯一含 session 的是 `ListGrantsBySession` ⇒ 确实列不出 |
| **AC#2** 四件最小出口 | **不成立＝接线缺口** | 函数与用例逐枚对号，★生产调用者枚数＝**0/0/0/0/0**：`session.NewManager` 全仓 `.go` rc=1（正控 `session.NewLedger` 3 命中）、非测试文件里 Manager 方法调用点 rc=1、`memory.Session` 在 panel/tools/cmd 产码里 rc=1、CLI 十枚名册⛔ sessions、`PumpSources`/`Snapshot` ⛔ 会话读者 |
| **AC#3** 迁移崩溃原子性 | **成立** | 真杀进程（父 `cmd.Process.Kill`＝TerminateProcess）六枚定点＋十二枚随机：合法只有 `NOTHING(v2-untouched)` 与 `ALL(v3-complete)` 两形，**⛔ 一枚 `HALF-STATE`**；★其中一发在同一事务里塞出 `wisp.db-wal size=13962712` 的未提交帧，杀后重开仍是 v2 原形、`integrity_check=ok`；新老库名册逐字节相同（11 表/10 索引），老库升级前缺的是整张 `session`/`session_turn`＋三枚索引 |
| **AC#4** 分叉不许吞原件 | **成立** | 去牙两发（物理改盘上字节、⛔ `-overlay`）：M1"分叉时顺手清理"红在 `manage_test.go:213 original lost turns (2 -> 0): the fork swallowed history`；M2"不带走历史"红在 `:226 fork turns = [] (<nil>), want 2`；还原 sha 回到 `38bff760…`、`cmp` 0 差异、复跑 `--- PASS` |
| **AC#5** 契约轴零字节 | **成立** | 五枚笔并集 41 路径，15 枚命名 pattern 命中**全 0**（`docs/PLAN.md`/`docs/specs/`/`internal/risk/`/`internal/panel/`/`internal/agent/approval/`/`thresholds.go`/golden/`allowlist.txt`/`frontend/`/`design/`/`.github/`/`ci.yml`/`go.mod`/`go.sum`/`.gitignore`） |
| **AC#6** 同秒并列那一支 | **ⓐ成立（范围修正）／ⓑ成立** | 本腿复跑四发：改后 20/20（`RUN=20 TOPPASS=20 TOPFAIL=0`）；★摘掉步进时钟＝`TOPFAIL=10`（D 卷）／`TOPFAIL=9`（C 卷）＝牙响；改前字节 `-count=20`＝`TOPFAIL=12`（实现腿两发是 7 与 12，同一散布）；ⓑ 答"看不见"，成因⛔ 是排错而是**压根⛔ 链**（＝AC#2 那一格） |

**票的整体口径**：AC#2 判不成立 ⇒ **本票⛔ 已完成**，"用户今天仍没有一条路回到昨天那次对话"在产品面上照旧为真；
差的那一段是**装配根构造一枚 `Manager` ＋一枚命令或面板读者**（`cmd/wisp/**`、`internal/panel/**` 是别的写面的家，本腿只报位置）。

## 逐格读数指向

- **AC#1** → `01-ac1-unmodified-reading.txt` ＋ `01-ac1-cli-sessions-unknown-command.txt`（2498 字节，rc=2）
  ＋ `01-ac1-cli-history-unknown-command.txt`（2497 字节，rc=2）＋ `01-ac1-apidoc-store-baseline.txt`／`-head.txt`。
  台面插曲⛔ 读数：第一发 `rc=127 sherpa-onnx-c-api.dll: cannot open shared object file`（导出树里 `.dll`＝0 枚，三枚⛔ tracked）＝**台面故障**，补 DLL 后重跑才拿到数。
- **AC#2** → `02-ac2-outlets-and-wiring.txt`（§2.1 落点表、§2.2 五把尺）＋原始命中 `02-ac2-greps-outlets.txt`／`02-ac2-greps-restricted.txt`／`02-ac2-panel-nontest-session.txt`。
- **AC#3** → `03-ac3-crash-atomicity.txt` ＋台件源码 `zz166v1crash_bench_test.go`（`gofmt -l` 0 枚、`go vet` rc=0，跑在仓外 `tree-head`）
  ＋ `logs/ac3-crash-parent.log`（171 行）／`logs/ac3-crash-sweep.log`（161 行，seed `1792117039198721100`）／`logs/ac3-roster-bench.log`
  ＋两枚真二进制台件 `bench-mkv2-main.go`（基线树写 v2）／`bench-dumpdb-main.go`（HEAD 裸读与升级）＋`logs/ac3-old-pre.txt`／`ac3-old-upgraded.txt`／`ac3-new-fresh.txt`／`ac3-roster-diff-new-vs-old.txt`（`ROSTER_DIFF_RC=0`）。
- **AC#4** → `04-ac4-fork-de-tooth.txt` ＋字节账 `04-ac4-and-ac6-hashes.txt` ＋ `logs/ac4-M1-session.log`／`ac4-M2-session.log`／`ac4-restored-green.log`／`ac4-M1-memory.log`。
- **AC#5** → `05-ac5-contract-axis.txt` ＋ `05-ac5-commit-roster.txt`／`05-ac5-union-roster-41.txt`／`05-ac5-spec02-count.txt`（spec 正文 9 表 vs `contractSpec` 11 表）。
- **AC#6** → `06-ac6-instrument-audit.txt`（§6.1 字节账、§6.2 那份证据件十一枚抽取件的自洽审计、§6.3 本腿四发、§6.4 ⓑ）
  ＋ `logs/ac6-after-list-count20.log`／`ac6-teeth-flatclock-list.log`／`ac6-teeth-flatclock-search.log`／`ac6-restore-green-count20.log`／`ac6-before-166r2-count20.log`／`ac6-volC-after.log`／`ac6-volC-teeth.log`。

## ★本格必答复的那两句（⛔ 含糊）

1. **AC#6 那枚 `:84`（交付态 `:95` `TestSessionSearchTitleAndBody`）到底⛔ 是硬币。**
   本腿在**同一枚做平时钟的突变字节**上跑它：`-count=20` ⇒ `TOPFAIL=0`（逐字绿行 `--- PASS: TestSessionSearchTitleAndBody (0.18s)`）。
   读码对上：它对结果的断言只有"命中 1 枚／0 枚"（`:112`/`:119`/`:122`/`:125`），⛔ 比两枚之间的先后。
   ⇒ **本格真关掉的硬币只有 `:69`（`TestSessionListByWorkspace`）一枚；`:84` 那一处改步进时钟＝防未来。**
   166-r2 在证据件 §4 自报的这处"与派单/票面前提不符"，本腿复跑**支持它**、⛔ 替它抹平、也⛔ 因此判 ⓐ 红（ⓐ 的判据是"改前 FAIL>0 复现＋改后 20/20"，两发都在）。
2. **"写了没接线"这一形成立** ⇒ AC#2 判语＝**接线缺口**、⛔ 写成本票已完成（派单 §AC#2 明令）。
   ⛔ 用 `go build`/`go vet` 的 rc=0 冒充接线：Go⛔ 为未使用方法报错，两枚 rc 与"⛔ 调用者"完全可以同时成立。本腿认的只有 `git grep` 的调用点枚数（＋整包颜色只当台面证明）。
   ⚠ 口径分歧具名交回：若把 AC#2 的射程读成"Go 侧函数与用例齐"，那一读下它是**成立**的；
   判语的分量全在"**出口**"两个字——票面 §这是什么（"今天没有一条路能回到昨天那次对话"）与本腿的 0 命中读数取前者。

## 上交编排者裁的料（⛔ 本腿裁、⛔ 本腿动一字；详数在 `07-tiebreak-material.txt`）

1. ★`AppendSessionTurn`（`dao_session.go:253-279`）**⛔ 刷新 `session.updated_at`** ⇒ 连发 10 轮实测 `updated_at 1777593660 -> 1777593660 (delta=0)`、列表次序不变 ⇒
   **"继续聊不把会话顶到前面"＝行为缺陷**（写 `updated_at` 的只有 Create/Rename/Archive/Fork 四枚点）。
2. tie-break 本体（`dao_session.go:125` 与 `:144` 同一枚 `ORDER BY updated_at DESC, id DESC`，秒级 `updated_at` ＋随机 hex id）：
   生产时钟形状 200 对同秒会话实测**次序与创建次序相反 91/200＝45.5%** ⇒ 接线之前机主看不见（⛔ 链），接线之后就是一枚硬币。
   ⛔ 本腿动这一字：换排序键＝行为变更，且"加一枚单调列"⛔ 止一种形状（`seq`/`rowid`/毫秒写入时刻都还⛔ 定过）。
3. 对照组成立：`dao_profile.go:86`/`:111` 同形 SQL 的 `id` 是 `INTEGER PRIMARY KEY AUTOINCREMENT`（`schema.go:31`）⇒ 单调⇒⛔ 硬币；
   `session.id` 是 `TEXT PRIMARY KEY`＋`Mint` 的 `rand.Read`（`session.go:94-100`、`schema.go:221`）⇒ 随机⇒硬币。病灶是 id 形状、⛔ 是 ORDER BY 写法。
4. **owner 轴（契约）**：`schema.go:139` 抬 `SchemaVersionTarget` 2→3＋`ddlV3` 两枚新表、`schema_test.go` 的 `contractSpec` 从 9 枚扩到 11 枚，
   而 `docs/specs/SPEC-02-data-storage.md` 正文仍 `^CREATE TABLE`＝**9 枚**（`grep -c` 现量）⇒ 代码与 spec 正文枚数不一致；
   票面硬约束 1 逐字是"⛔ 新造存储…加**列**要走迁移"，交的是加**两枚表**（166-r1 自己在票尾两处具名报过：契约债 H1/H2＋"加列表述偏离"）。
   这一格落在"改契约＝人工批准"轴上（AGENTS.md §0 第 2 句／SPEC-12 §4.1）⇒ 本腿⛔ 判红也⛔ 放行，只交位置。
5. **DAO 层⛔ 分叉用例**（本腿顺手量到）：M1 突变打在 `dao_session.go`，`go test ./internal/memory/ -count=1` 仍 `ok … 12.429s`；
   全仓 `.go` 的 `ForkSession(` 调用点＝`manage.go:145` 一枚产码＋`manage_test.go` 一枚用例 ⇒ "不吞原件"今天**单层把守**，接线之前值得补一枚 DAO 级用例（⛔ 本腿动）。

## 未闭合与限制（具名，⛔ "应该没问题"）

- **⛔ 需要的那一发本腿做出来了**：AC#3 的"真把进程杀掉再重开"＝六枚定点杀＋十二枚随机杀，⛔ 缺件。
  唯一诚实边界写在 `03-ac3-crash-atomicity.txt` §3.1：注入的只有"死在哪一刻"（子进程在真 2→3 那步的 `apply` 里喊一声 READY 后被父进程打死），
  DDL 与事务/Commit 机器都是真的；`spill` 那一发专门堵"事务内的话压根没落盘"这句反驳。
- **`三枚冻结件` 这一条量不到**：AGENTS.md／`issues/README`／`HANDOVER.md` 里"三枚冻结件"命中 6 处、⛔ 一处给出名册 ⇒ 按"未定义即停"处理，
  AC#5 只交 15 枚**命名** pattern 的读数（全 0）。要量这一格请先给名册。
- **`forkTurns[i].ID == before[i].ID` 那条"复制⛔ 搬家"断言（`manage_test.go:232-234`）本腿⛔ 另做发去验牙**：
  造它要让 `session_turn.id` 显式复用父行 rowid，超出"顺手清理"那一形 ⇒ 本件只交"断言在盘上"，⛔ 交"验过"（`04-ac4-fork-de-tooth.txt` §4.2）。
- **⛔ 做的一件事**：`-count=20` 那两发在 D 卷与 C 卷各跑了一遍（同一颜色：D 卷 TF=10／C 卷 TF=9；改后两卷都 TF=0）；
  起因是票 311 的 GUARD 在 `scripts/d22scan.sh` 那一发以 **rc=2＝射程拒绝**挡下 `GOTMPDIR=D:/tmp`（⛔ 当红，按派单口径）。
  门禁两枚补跑在**授权卷/仓根**：`sh scripts/d22scan.sh`（`tree-head`，GOTMPDIR＝C 卷）⇒ **rc=0 clean**（bans #1-5 internal/=239、cmd/=41｜#6 frontend/=116｜#7 tools/=25｜#8 design/=30、frontend/=116、internal/=573、cmd/=138）；
  `sh scripts/check-path-length-budget.sh --with-self-test`（仓根，只读）⇒ **rc=0 VERDICT GREEN**（越帽 60、名册 60、worst 229）。
  ⚠ 这两枚 rc 覆盖的是**整棵树**（含别腿在飞字节），本腿只当台面读数，⛔ 替别的票收账。
- **别腿台面**：起手 `go build ./...` 三处都 rc=0（`tree-head`/`tree-base`/共享工作树），此刻⛔ 撞上 winsec/tools 的在飞红；
  工作树里 `go.mod`（webview2 从 indirect 提到 direct）与 `.gitignore` 是**别人的脏改动**，本腿⛔ 碰、⛔ 记进本票判语，只说明为什么一切读数走导出树。
- 本腿⛔ 勾任何 AC 框、⛔ 动 `Status:`/票名/`-done`、⛔ 改产码或测试件一字、⛔ push；临时件只建不删。

## 交件自证三把尺

1. **提交号**（本腿四笔，全部只 commit）：`28e37882`（台面＋AC#1/#2/#5）·`e396301d`（AC#3 台件与读数）·`fb157e95`（AC#4/#6/tie-break）·本表与票尾那一笔（见 `git show --stat HEAD`）。
2. **裁决表与读数件的绝对路径＋字节数**：
   - `D:\work\workspace\projects plans\Wisp\docs\evidence\s1\166-session-verdict.md`＝本文件（字节数见 `git show --stat` 那一行）
   - `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\166\v1\`＝7 枚读数件＋19 枚 logs＋3 份台件套件源码（`ls -la` 逐枚字节在交件正文里）
3. **本腿至少复跑的那两发 AC#6 逐字红句／绿句**：
   - 绿（改后，交付字节 `2bf320df…`）：`--- PASS: TestSessionListByWorkspace (0.06s)`，计数 `RUN=20 TOPPASS=20 TOPFAIL=0`（D 卷与 C 卷各一发，同色）
   - 红（摘掉步进时钟，字节 `7cfe34ee…`/C 卷；`b958d176…`/D 卷）：
     `    manage_test.go:87: workspace filter = [sess_91e70fe9130cfa7690030fc544e73b1d sess_3abcf8fc2cedfec1db78230d9d0b6e24], want [sess_3abcf8fc2cedfec1db78230d9d0b6e24 sess_91e70fe9130cfa7690030fc544e73b1d]`
     计数 `RUN=20 TOPPASS=11 TOPFAIL=9`（D 卷那一发是 `TOPPASS=10 TOPFAIL=10`）
   - 还原后绿：`manage_test.go` sha 回到 `2bf320df04e2819a…`、`cmp` 与 pristine **0 差异**、`-count=20` 两枚用例 `RUN=40 TOPPASS=40 TOPFAIL=0`
