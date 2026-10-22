# 票 166 落地腿 `166-r1` · `AC#2` 四件最小出口（**Go 侧**）· 2026-10-22

腿＝**写**（实现者）。只 commit ⛔ push。本件是**读数件**；判语归非实现者——`166-v1` 那枚"AC#2＝接线缺口"的判语
本腿**不翻**（残余逐字见 §9）。台面：主检出 `D:\work\workspace\projects plans\Wisp`；读数件 `D:\tmp\wisp\20261022\166r1\logs\`。

## 0. 一句话

四件的 Go 侧出口（列／搜／改名／归档带撤销）**已在盘上**（10-13 落、10-18 装配半边）。本腿做两件事：
① 现量"每一发在**未修码**上响不响"（⛔ 转述旧日志）；② 量到并修掉**一件真缺陷**——`SetSessionArchived` 的两个分支都写
`updated_at=now` ⇒ 用户按"撤销"拿回来的是**改了排序**的列表，那不是撤销。⛔ 抹平：生产调用者仍 **0/0/0/0/0**（§9）。

## 1. 票面前提现跑（腐坏 1 处，如实记）

| 票面逐字 | 现跑 | 判 |
|---|---|---|
| `Status:` **"ready-for-agent，但立而不派**（排在票 158 验收与 161 之后…）" | 票面 `:3` 逐字仍在 | 前提两条已满足：`158-…-done.md`、`161-…-done.md` 在盘 ⇒ "立而不派"的理由消失、编排者已定派 |
| `AC#6` 地基 `manage.go:93` 写 `m.now().Unix()` | 今天 **`manage.go:106`**，**偏 +13**（10-18 腿在 `:49-60` 插了 `var _ ManageStore` 守卫块） | 腐坏 1 处；⛔ 抹平 |
| `AC#6` 地基 `dao_session.go:125` 排 `ORDER BY updated_at DESC, id DESC` | `:125` 逐字命中 | 对 |
| `AC#6` 地基 `session.go:96` 由 `rand.Read` 铸 id | `:96` 逐字命中 | 对 |
| `AC#2` 未勾、其余五格已勾 | 票面 `:25` `- [ ]` | 本腿只做这一格 |

## 2. 起手基线（未动代码，主检出）

`go test -count=1 -v` × 两包，**四数分开**：`internal/session` RUN=24 PASS=24 FAIL=0 SKIP=0 ／
`internal/memory` RUN=45 PASS=43 FAIL=0 SKIP=2（`logs/01-baseline*.txt`）。两包 rc=0。

## 3. 四件逐件："这一发在未修码上响不响"（★"不响"也如实报）

台面＝`git archive 2f41a80cb` 的**仓外导出树**（1003 枚 `.go`；⛔ 工作树）。
★**先说尺**：Go 1.27 上 `go test -run <不存在的名字>` 的 **rc 仍是 0**，只有输出里的字面 `[no tests to run]` 认得出
（量到：`logs/03-unfixed-readings.txt` §C）。⇒ ⛔ 拿 rc 判"响不响"；本件一律贴**字面**。

| # | 件 | 未修码上（`2f41a80cb`） | 修后（`17cfe8cf`） |
|---|---|---|---|
| 1 | 列（按时间／按工作区） | **不响**：`internal/session` 只有 `doc.go`/`grants.go`/`grants_test.go`/`session.go`/`ticket224_*`；无 `manage.go`、无 `ListSessions` | `TestAC2OutletListAssembled` PASS；本腿 `TestAC2Undo…:113` 顺序断言绿 |
| 2 | 搜（标题与正文） | **不响**：`go doc ./internal/memory Store \| grep -i session` 只有 `ListGrantsBySession`（D45 授权台账，不是会话行）；`SearchSessions` 全仓 0 命中 | `TestAC2OutletSearchAssembled`、`TestAC2SearchHonours…` PASS |
| 3 | 改名 | **不响**：`RenameSession` 在 `2f41a80cb` 的 `internal/` 下 0 命中（`git grep -c` rc=1）；schema 连 `session` 表都没有（`CREATE TABLE … session` 在 `internal/memory` 下 rc=1） | `TestAC2OutletRenameAssembled`、`TestAC2OutletsRefuseUnknownSession` PASS |
| 4 | 归档带撤销 | **不响**：`SetSessionArchived` 全仓 0 命中 | 见 §4／§5：修前**红**，修后绿 |

★强化读数（同一台面，`-overlay` 把四件探针送进未修树）：**编不过**——
`undefined: NewManager`、`undefined: ManagerOptions`、`undefined: memory.SessionListOptions`，rc=1；
同一发探针在 HEAD 树上 rc=0 ⇒ 未修码**连这四个出口的名字都没有**（`logs/04-overlay-probe.txt` §G）。
★正控（同尺同树）：`-run '^TestTicket224MintIsRandomAndValid$'` ⇒ `--- PASS`（尺灵，`logs/35-controls.txt` §a）。

## 4. 本腿落地（两枚文件，`17cfe8cf`）

- `internal/memory/dao_session.go`：`SetSessionArchived` 两分支**只写 `archived_at`**（把 `, updated_at=?` 与 `now` 摘掉，
  今天 `:205`），并把 326 AC#1 注释里"archived"从 touch 清单里摘掉（**注释**，⛔ 契约文字）。
  **行为变更具名**：归档不再算"活动"；撤销把行**放回原位**。`ORDER BY` 逐字未动 ⇒ ⛔ 碰 AC#6ⓑ 那格。
- `internal/session/ticket166r1_ac2_gaps_test.go`（**新增**，224 行）：三枚用例，全部经 `ManageStore` 接口装配真 `wisp.db`。

| 件 | 新用例（逐名） | 为什么今天响 |
|---|---|---|
| 4＋1 | `TestAC2UndoPutsTheConversationBackWhereItWas` | 修前**红三行**（`:86` 归档改写了 last-activity、`:105` 撤销搬了行、`:113` 撤销没把顺序放回）；修后绿 |
| 2 | `TestAC2SearchHonoursArchiveMaskAndWorkspace` | 今天绿＝**钉住**（此前任何一层都没有"归档要从**搜**里也藏起来"和"搜按工作区收窄"的用例）；牙 T2/T6/T7 逐枚让它红 |
| 3＋4 | `TestAC2OutletsRefuseUnknownSession` | 今天绿＝钉住"不存在的 id 必须报错、⛔ 静默成功"（含正控：真 id 三发都返回 nil）；牙 T3/T4/T8 让它红 |

⇒ 四件**每件都有今天响的用例**：件 1／2／3 的"响"由 10-18 的装配用例＋本腿两枚钉住＋八发牙口证明（§5），
件 4 的"响"是**修前就红**（§5 T5 逐字）。

## 5. 反向判据：八发（把"这一跳"摘掉 ⇒ 用例必须转红）

台面＝`git archive 17cfe8cf` 导出树（1157 枚 `.go`），**物理改字节**、串行、每发之后 `cp` 还原。
逐发 rc：**变异 rc=1**（红）、**还原 rc=0**（绿）、**原样 rc=0**（绿）。

| 发 | 摘掉哪一跳（`grep -n` 证明） | 红名＋红行（逐字） |
|---|---|---|
| T1 | `manage.go:117` `return nil, nil // TEETH T1`（item 1 `Manager.List` 的 store 跳） | `TestAC2Undo…:72` premise broken: list = []；`OutletListAssembled:87/:90/:98/:101`；`OutletSearchAssembled:139`；`OutletArchiveUndoAssembled:212` |
| T2 | `manage.go:126` `return nil, nil // TEETH T2`（item 2 `Search`） | `TestAC2Undo…:120`；`SearchHonours…:151/:155/:167/:171/:177`；`OutletSearchAssembled:123/:130`；`OutletRenameAssembled:172` |
| T3 | `manage.go:136` `return nil // TEETH T3`（item 3 `Rename`） | `OutletsRefuseUnknownSession:197` Rename(unknown id) reported success；`OutletRenameAssembled:162` title = "before name", want "after name"／`:172` |
| T4 | `manage.go:148` `return nil // TEETH T4`（item 4 `Unarchive`） | `TestAC2Undo…:102/:113/:120`；`SearchHonours…:177`；`OutletsRefuseUnknownSession:203`；`OutletArchiveUndoAssembled:209/:212` |
| **T5** | `dao_session.go:200` `now := observe.NowWallUTC().Unix() // TEETH T5` ＋ SQL 里装回 `updated_at=?` ＝ **本腿的修复被摘掉** | `TestAC2Undo…:86` archive rewrote the last-activity stamp: updated_at 1792659720 -> 1792678029（**冒烟的那一行**）；`:105` undo moved the row instead of putting it back；`:113` undo did not restore the position |
| T6 | `dao_session.go:138` `sqlText += "" // TEETH T6`（搜的归档掩码） | `SearchHonours…:167` archived conversation still surfaces in search |
| T7 | `dao_session.go:140` `if false && opts.Workspace != "" // TEETH T7`（搜的工作区过滤） | `SearchHonours…:155` workspace-scoped title hit = [两枚]／`:159` body hit leaked across the workspace filter |
| T8 | `dao_session.go:181`＋`:209` `n < 0 // TEETH T8`（两处 `ErrNotFound` 守卫） | `OutletsRefuseUnknownSession:197/:200/:203`（改名／归档／撤销三发全部"reported success"） |

**还原核 `sha256` 三处逐字相同**（提交版＝`git show 17cfe8cf:<path>`／还原台＝导出树／工作树；八发前后各量一次，`SAME=yes`）：

- `internal/session/manage.go` → `0577b500c1fc705ee0246ab6b87d74782dfbb37831930550774cbd645abc7e6d`
- `internal/memory/dao_session.go` → `4bb05526ce8651e94190125505b4f87ec359d50826e9aebdccf75a8b1aec8268`
- `internal/session/ticket166r1_ac2_gaps_test.go` → `9729adefb98cc2a2e112446d06164af73a772fc39e5c12c0f6fcefcdca58f0be`

工作树**从未**携带变异（`git status --porcelain -- internal/session internal/memory` 空）。逐字读数：`logs/10-teeth.txt`。

## 6. "归档带撤销"的撤销路径（⛔ 纸上）

- 撤销是**一等出口**：`manage.go:147` `func (m *Manager) Unarchive`（`grep -n` 落地，§4 已贴）。
- 撤销**被真的触发**：`TestAC2UndoPutsTheConversationBackWhereItWas` 调 `Unarchive` 并从存储读回三样——
  `archived_at` 清空（`:102`）、`updated_at` 逐字回到归档前的值（`:105`）、默认列里回到**原位置**（`:113`）；
  外加**第二条读路**：撤销后同一会话重新可被**搜**到（`:120`）⇒ 撤销不是"翻一个列"。
- 牙口两发证明它**能被摘掉**：T4（摘 `Unarchive` 的 store 跳）红 3 行；T5（把 `updated_at=?` 装回去）红 3 行。
- 行与正文都还在：`Turns` 逐字断言 `kept through the undo`（`:123`），归档＋撤销不吞历史（AC#4 不重开）。

## 7. 门禁（自己现跑，逐发 rc）

| 门 | 命令 | rc／读数 |
|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **rc=0** |
| G2 | `go vet ./internal/session/ ./internal/memory/` | **rc=0** |
| G3 | `go test -count=1 ./internal/session/ ./internal/memory/` | `session` **RUN=27 PASS=27 FAIL=0 SKIP=0** rc=0 ／ `memory` **RUN=45 PASS=43 FAIL=0 SKIP=2** rc=0（⛔ 相减） |
| G4 | `sh scripts/d22scan.sh` | **rc=0** `clean - no D22 ban violations`（`internal/` 619 ＋ `cmd/` 148 枚 `.go` 扫过） |
| G5 | `sh scripts/check-path-length-budget.sh` | **rc=0** `VERDICT GREEN`（band：over-the-hat=61、wall interval=0） |
| G6甲 | `git archive 17cfe8cf` 树 + `gofumpt -l`（**判据**） | 20 枚未格式化，**本腿两枚 0 枚**；同尺打父提交 `5b6338f6` 树也是 20 枚，`sort` 后 **`diff` rc=0**（逐行同） |
| G6乙 | 工作树 `gofumpt -l`（**只作诊断**） | 54 行、**本腿两枚 0 行**；54≠20＝工作树里有别的腿的未跟踪件（`gofumpt -l` 走文件系统、⛔ 只走索引）⇒ 甲乙量的是两个对象 |

★G4 第一发 **rc=2**，**不是缺陷**：`d22scan.sh` 里 ticket 311 AC#1 shape B 的守卫**拒跑**——我把测试临时根设在 D: 卷
（`GOTMPDIR=D:/tmp/…`），它的参照卷是 C。把临时根摘回 C ⇒ rc=0。★如实记：⛔ 把"守卫拒跑"读成"门禁红"。
★G3 的两组四数是在**参照卷 C**上重取的（D: 卷上那一组同数，见 `logs/11-gate-*.txt`；memory 开库会走 winsec 封印，
色彩与卷相关＝ticket 311 的教训，故以 C 卷为准）。

## 8. 红名册（★主检出，⛔ 导出树）

尺：`go test -count=1 ./internal/...`（**30 包**），两侧**同一棵树**，只把本腿两枚文件用 `-overlay` 拉回 `5b6338f6` 的形。
两侧各 **10 枚**，**9 枚同名**；1 对**一进一出**：before 独有 `TestFloodRulerFiresOnAnUnboundedSink35r8`（`internal/panel`），
after 独有 `TestGateHalfDuplexClosedZeroFrames`（`internal/audio`）。定点复跑（两侧同尺）：panel 那枚两侧 PASS、
audio 那枚两侧 PASS ⇒ **满包并发下的负载抖动、⛔ 本腿因果**（`logs/40..43-swap*`）。
⇒ **新增红（归因本腿）＝0**。十枚逐名：`logs/31-roster-after.names`（如 `TestApprovalCardViewJSONKeysMatchFrontendTypes`、
`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 等，均与本腿两包无关）。

## 9. 残余（⛔ 本腿写面，如实报）

四件的**生产调用者仍＝0/0/0/0/0**：`git grep -n "session\.NewManager" -- '*.go'` **rc=1**（0 命中）；`_test.go` 之外的
Manager 方法调用点 rc=1；`internal/panel`／`cmd/wisp` 对 `memory.Session`/`ListSessions`/`SessionByID` **0 命中**。
接线要 `cmd/wisp`（**本波归别的腿**）＋ `internal/panel`（**归别的腿**）⇒ 本票"出口"在产品面上**仍不成立**（`166-v1` 判语照旧）。
⛔ 拿 `go build` rc=0 或"用例绿"冒充接线；⛔ 造第二个 Store／假调用者凑数。

## 10. 提交（⛔ push）

| 提交 | `git show --name-only` | `.go` 枚 | 越界 |
|---|---|---|---|
| `17cfe8cf` | `internal/memory/dao_session.go`、`internal/session/ticket166r1_ac2_gaps_test.go` | **2** | 0 |
| **本笔**（本件） | `docs/evidence/s1/166-r1-impl.md` | 0 | 0 |

越界尺：把路径逐个打 `^(internal/session/\|internal/memory/\|docs/evidence/s1/166-r1-[^/]*\.md$)` ⇒ 越界行 **0**（rc=1）。
正控（同尺）：`17cfe8cf` 的 go-files=**2**（尺看得见 `.go`）／docs-only 提交 `f895fe5c` go-files=**0**（尺不会乱数）。
