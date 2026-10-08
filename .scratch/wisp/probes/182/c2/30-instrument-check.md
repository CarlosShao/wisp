# 182-c2 · 第 4 问查重：仓里有没有**已经存在**的同类判据/台件（按"这句话在说什么"扫）— 锚 HEAD `a7e9b2e7`

派单那句"⛔ 只搜符号名必错，要按'这句话在说什么'列变体再扫"⇒ 本件先列变体，再逐族扫，最后给"覆盖到哪儿、哪儿仍是空白"。

## 1. 变体表（六问里"有没有仪器钉得住"这句话在本项目可能长什么样）

| 这句话在说什么 | 盘上的说法变体（本腿实际扫过的词） | 扫到的仪器族 |
|---|---|---|
| 名册不许变宽 | `whitelist`／`roster`／`Door`／`knownComposerMethod`／`wantFullInboundRosterSize253`／`exactly one door` | `internal/panel/git_test.go`（13 枚 `^func Test`）／`inbound_roster_253_test.go`／`composer_test.go`／`l2_grant_boundary_test.go`／`config_route_248_test.go` |
| 不许把显示接成动作 | `no model-callable tool`／`action leg`／`off the table`／`never names` | `git_test.go:364 TestGitDimensionHasNoModelCallableTool`（＋`git_test.go:486 TestPlantedGitToolShapesGoRed` 正控）／`frontend_hygiene_test.go:281 TestFrontendNeverNamesAnApprovalDecision` |
| 声明面要与页面接口对账 | `MatchFrontendTypes`／`ContractTypes`／`pairs`／`FourwayAgree` | `composer_test.go:48 TestComposerContractTypesMatchFrontend`（**今天在册红**）、`TestApprovalCardViewJSONKeysMatchFrontendTypes`（红）、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（红）、`TestC21DesignTokensFourwayAgree`（红） |
| 载荷不许带凭据 | `canary`／`secretShape`／`leaked onto the wire`／`Sentinel` | `cmd/wisp/panel_config_248_test.go`（含恒真正控 `:218`）、`internal/panel/inbound_raw_leak_35r7_test.go`（今天新增族）、`cmd/wisp/instructions_200r2_test.go:73` |
| 由测试自己跑的"普查" | `census`／`普查`／`declaredRefSlots` | **`internal/agent/approval/ticket146_liveapprovals_backing_test.go`** 一族：`:62` 逐字"TYPE census (declaredRefSlots) no longer depends on the fixture to notice a…"、`:66` "it is the test of the census" ⇒ **本项目已有"把结构普查做成常驻尺"的先例**，但它普查的是**引用槽位形状**，⛔ 不普查"某一栏数据在 Go 侧有没有源" |
| 冻结件没被碰 | `git diff --name-only`／`--numstat`／`d22scan`／`allowlist` | **CI 里没有这一族**：`grep -n 'PLAN.md\|specs\|allowlist' tools/d22scan/main.go` ⇒ 只有 `allowlist.txt` 抑制机制（`:260 loadAllowlist`），**没有一枚 ban 管"文档/契约文件未被改动"**（rc=0）。人工台件广泛存在：`grep -rl 'git diff --name-only' docs/evidence/s1/ .scratch/wisp/probes/` ⇒ **182 个证据件在用这把尺**（rc=0）＝**项目定式的"人跑尺"，不是 CI 牙** |

## 2. 现量规模（尺与 rc 都在下面，⛔ 没挂管道取 rc）
- `ls internal/panel/*_test.go \| wc -l` ⇒ **21**
- `ls cmd/wisp/panel*_test.go cmd/wisp/subagent*_test.go \| wc -l` ⇒ **20**
- `grep -c '^func Test' internal/panel/git_test.go` ⇒ **13**
- `grep -rn 'func TestTheRendererHoldsExactlyOneDoorToTheHost\|func TestPlantedRendererDoorShapesGoRed' internal/panel` ⇒ rc=0，各 1 命中 ⇒ **`bridge.go:32-40` 注释里指名的两枚测试真身存在**（这一条按第 129 条纪律必须实跑，不能靠注释追认）
- `grep -rn 'census' --include=*_test.go internal cmd` ⇒ rc=0（首族＝票 146 那枚槽位普查）

## 3. 逐框覆盖结论（第 4 问的查重量）
| 框 | 已有同类判据？ | 空白具名 |
|---|---|---|
| AC#2 三档定性 | **无**（没有一枚尺读"哪一栏有没有 Go 侧源"；票 146 那族读的是引用槽位） | 空白＝"数据源普查"这一类判据在本项目**从来没被做成常驻尺**；能补的形状只有"每把尺可复跑＋正控" |
| AC#3 归属指认 | **无** | 同上；唯一凭据＝逐张票读框（本腿尺见 ①Q4，rc=0） |
| AC#4 与 181 分工 | **半有**：`git_test.go:383-391`＋Door 3 正控 `:505-521`＋`inbound_roster_253_test.go:59` | ⚠ **`composer_test.go:48` 的 `pairs` 名册不含 `GitView`** ⇒ 嵌子格不与 TS 对账；⚠ 那把尺**读 `frontend/src/lib/panel.ts` 工作树**，与本票 AC#7 撞面；⚠ 它**今天在册红** |
| AC#5 雷区 | **有，四枚独立钉＋两枚静态门**（`git_test`／`inbound_roster_253`／`composer_test:502`+`:533`／`frontend_hygiene_test:281`／d22scan ban #6／`ErrPanelAllow` 一族） | ⚠ **零枚尺扫 `design/doubao/demo/**`**（ban #6 只走 `frontend/`，`tools/d22scan/main.go:23` 定义逐字"in frontend/"）⇒ demo 里那些 `data-*-toggle`／`data-terminal-add` 动作面**根本没有仪器**，只有腿的尺 |
| AC#6 契约轴 | **无 CI 牙**；人工尺 182 个证据件在用 | 空白＝"冻结件未被碰"从未被做成常驻门（本腿扫过 `scripts/*.sh` 13 枚与 `tools/d22scan/main.go` 全部 ban） |
| AC#7 排程约束 | **无**（也没有任何机制能读"阻塞关系"） | 前半句是**恒真句**（179/176/177/175 从不引用票 182）；真关系在 188/189/190/191/192/181/186/180，那张反向依赖表**一块没写** |

## 4. 与"同类普查表"的查重（⛔ 不只看符号名）
`grep -rln '可开工\|能不能开工\|不可派\|死格' .scratch/wisp/probes/ docs/evidence/s1/` ⇒ rc=0，命中件里**成表的只有**：`probes/35/a6/verdict.md`（六格可开工性普查，本件形状的来源）、`probes/111/v1/verdict.md`、`probes/35/c2/summary.md`、`probes/35/census-resident-panel/summary.md`、`probes/274/a2/costs.md`、`probes/166/a1/census.md`、`probes/rate-census-4b/buckets.md`。
⇒ **票 182 名下从来没有这种表**（`probes/182/` 只有 `a1/census.md`＝逐堆普查、`c1/100-census.md`＝逐堆二轮）。本件＝**第一枚"每格能不能开工"**，⛔ 不重复 a1/c1 的活（它们答过的栏见 `verdict.md` §②）。
