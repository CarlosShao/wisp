# `194-batch1-cost-door-r1` — 票 194 分批落地第一枚（cost.summary）＋三枚受阻拼写现量停报

- **程**：`194-wisp`（写腿，写面＝`internal/panel/**`；只 commit 不 push）。
- **锚**：起手 HEAD＝`e90ea3cf`（现跑）；票面 11-07 代价表锚＝`a1c3f245b`，两者之间线上已推进 14 枚 commit ⇒ **本件全部读数现跑，不从票面抄**。
- **工单**：`.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md`（AC 框零勾，勾归非实现者）。
- **派单批次**：代价表裁定「乙堆先行（`tasks.list`/`task.detail`/`approval.current`/`cost.summary` 四枚最低价可同批）」⇒ 本腿现量后**只落成 1 枚**（§2 两把冻结仪器把另外三枚拦死，逐枚具名）。

## §1 现量（起手，全部可复跑）

| 读数 | 现量值 | 尺 |
|---|---|---|
| 规格名册（SPEC-08:163-174） | 18 invoke＋5 事件，拼写逐字未变 | `sed -n '155,178p' docs/specs/SPEC-08-ui-ball-panel.md` |
| 代码名册（bridge.go 常量） | 7 枚（原 4＋248 两枚＋320 一枚） | `sed -n '41,76p' internal/panel/bridge.go` |
| `disp.Handle(` 生产调用者 | 2 枚（`panel_host_windows.go:1191`＋`panel_inbound.go:174`） | `git --no-pager grep -n 'disp\.Handle(' HEAD -- cmd internal \| grep -v _test.go` |
| `panel.*` 字面量（git_test 抽取器） | 4 枚 | `grep -cE '=\s*"panel\.' internal/panel/bridge.go` |
| 前端 `panel.ts` 是否携带乙堆名 | 否（仅 `config.models.list`） | grep `frontend/src/lib/panel.ts` |
| 起手 `internal/panel` 测试 | ok 16.2s（票面 3 枚在册红已被 415/421 关闭） | `go test -count=1 ./internal/panel/` |

## §2 三枚受阻拼写（两把冻结仪器，均现撞；未落地、未起别名、未改仪器）

1. **`approval.current`／`approval.queue`**：`l2_grant_boundary_test.go` 的 `grantRouteWords`（:192-195 含 `"approval"`）＋`TestAnsweredPanelRoutesCarryNoApprovalDecision`（:1222-1239）——守卫回答任何带 approval 词的路由即红。**AC#6 明列该文件不动**。代价表的「approval.current 改动面=1、代价低」在今日树上**不成立**：改动面＝改冻结仪器或改规格表，两者都须机主批。**规格表 vs ban#6 仪器正面冲突＝机主拍板项。**
2. **`tasks.list`／`task.detail`**：`internal/tools/task_state_188_test.go:323-339`——任何值含 `"task"`/`"state"` 的 `Method*` 常量＝「panel-side write leg」（正控 `MethodTaskStateWrite`）。腿侧实测：落两枚 task 常量 ⇒ 该测试红 2 枚（`bridge.go:99/100` 逐行点名）；`internal/tools` 不在本腿写面 ⇒ 整体退回。**建议归编排者**：188 仪器按读写分形收窄（读门豁免＋正控保留），或机主另批。
3. **`cost.summary`**：不含 task/state/approval 词 ⇒ 干净落地（§3）。

## §3 落地内容（追加式，未回改任何既有行）

| 文件 | 内容 |
|---|---|
| `internal/panel/bridge.go` | `MethodCostSummary = "cost.summary"` 常量（注释钉三件命名事实＋两枚受阻拼写出处）；守卫 case 单行延伸 8 枚（l2/253 锚形状，行注释钉住「不许换行」） |
| `internal/panel/composer_dispatch.go` | 接口 `CostSummaryHandler`（单读方法）＋插座 `Costs`＋case（nil→按名拒绝）＋头注释门数 7→8 |
| `internal/panel/cost_handlers.go` | `CostReadHandler`＋`CostReader`（只读面声明）＋`CostDayRow`（整数微单位）＋`ErrNoCostReader`/`ErrSelectorRefused`＋`wisp-cost-summary:` payload 行（无点、无 `panel.` 前缀＝双仪器安全形状）＋空渲染 `[]` 不 null |
| `internal/panel/cost_roster_194_test.go` | AC#4①②③ 常驻判据测试（见 §5） |
| `internal/panel/inbound_roster_253_test.go` | 登记摩擦照程序：名册＋`cost.summary`、`wantFullInboundRosterSize253` 7→8、`wantNonPanelPrefixedInbound253` 3→4、两处 plant 串同步 |
| `internal/panel/bridge_test.go` | `rosterPageGapLedger194` 页面镜像缺口台账（**加牙不弱化**：未登记未镜像仍红＝plant 正控保留；双向审计＝幻影登记红＋过期登记红）；理由＝`frontend/src/lib/panel.ts` 归外部 agent（Q-51），按仓内 `composerTSGapLedger` 先例登记 |

## §4 门禁读数（全绿）

- gofumpt（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）`internal/panel/`＝**空**。
- `go vet ./internal/panel/`＝rc0；`go build ./cmd/... ./internal/... ./tools/...`＝rc0。
  ⚠ `go build ./...` 红于**未跟踪在飞件** `build/mut_rvw.go`（别人的变异模板，`git log` 零命中＝未提交；登记未动、未删）。
- `go test -count=1`：panel **ok 21.4s**｜tools **ok 36.2s**｜config **ok 4.0s**｜risk（单包）**ok 13.6s**。
- `sh scripts/d22scan.sh`＝**clean**（bans #1-5 internal/=390、cmd/=65；ban#8 internal/=976、cmd/=250 全扫 0 违规；runtests PASS=45 FAIL=0）。
- `sh scripts/check-path-length-budget.sh`＝**VERDICT GREEN**（tracked=19269、over-budget=62、roster=62、not-in-roster=0）。
- `sh .scratch/wisp/probes/154/gate-clauses.sh`＝**14/14 相符、聚合退码 0**（环腿读数逐枚＝各基线：G1b=4、G2=2、G5=2、G5neg=14、G6neg=1、G7=6、G7neg=7）。

## §5 AC 对照（勾归非实现者，本件只给读数）

- **AC#1**：逐枚四问＝本件 §1/§2（枚数名册全现跑）。
- **AC#2**：三条禁区零触碰——本批无 decide／无权限状态写／无安全节放宽；`config.set` 未动。
- **AC#3**：白名单扩张逐枚落账＝台账 **A1681**（文件/行/理由/边界/撤销口令齐全，只追加）。
- **AC#4**：① 名册外拒＋审计（`TestRoster194OffRosterNameIsRefusedAndAudited`）；② 无裁决通路（`TestCostDoorCarriesNoDecisionPath`：verdict 形字段被共享 parse 丢弃、接口无写法）；③ 归因不重写（`TestCostDoorNeverRewritesAttribution`：RequestID/Source 逐字进审计行）。
- **AC#5**：诚实边界——`cost.summary` 今天**无真听众**（装配在 `cmd/wisp/panel_inbound.go` 生产码，本腿禁改）；nil 插座按名拒绝。下一波 cmd/wisp 写腿接 `memory.Store.ListCostDaily`（`dao_misc.go:169`）适配器。
- **AC#6**：`docs/PLAN.md`／`docs/specs/**`／allowlist/thresholds/golden/审批超时/`go.mod` 零字节；三枚冻结测试文件零字节（diff 自证见 commit）。
- **AC#7**：§4 全绿；CLI 面未跑（未动 cmd/wisp，无对象）。

## §6 纪律自证

- 写面＝`internal/panel/**`＋本件＋台账＋票面 Progress log；`cmd/wisp` 零改动（两枚未跟踪 323 腿测试文件登记未动）。
- 自建自删两枚中间稿（`task_handlers.go`、`roster_history_194_test.go`，task 门退回时移除，未进 commit）。
- 在飞件具名登记未动：`build/mut_rvw.go`（未跟踪）、`cmd/wisp/panel_resident_323_windows_test.go`＋`panel_resident_assembly_323_test.go`（未跟踪）。
- 只 commit 不 push；add/commit 显式 pathspec。
