# 188 — 面板"任务监控"栏里**子代理状态**与**后台任务状态**这两堆今天没有 Go 侧来源（owner 09-28 放权：缺就补）

- Status: **ready-for-agent（先只读设计核，再落地）**。⚠ **功能票**。
- 来路：owner 09-28 15:2x 原话（逐字在台账 `A373`：「我现在给你放权，缺什么就给我补」）＋他圈了右上角那一栏的截图。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5／§7（`180-c1＋182-c1` 交件，编排者收）。
- 关联：票 182（那一栏的堆数与三档定性＝本票的上游）· 票 145（快照载体十六字段）· 票 163（常驻终端会话，**不同物**）· 票 33／35（面板宿主与桥，硬前置）

## 现量（`180-c1＋182-c1` 本程现跑，锚 `e9ef94d0`；编排者复跑同值）

| 维 | 今天有什么 | 缺什么 | 尺 |
|---|---|---|---|
| **子代理状态** | **规格真空**（不是"有规格没实现"） | 全仓没有"子代理"这一维的任何数据模型 | `grep -rniE "Subagent\|sub_agent\|SpawnAgent\|DelegateTask\|fanout" --include=*.go internal/ cmd/` ⇒ **0 命中** |
| **后台任务状态** | `internal/tools/task.go:110-152` 有 `Record`／`Look`／`Count`（任务名册） | **没有"状态"这一维**（成功／失败／进行中／被拒／取消都取不到） | 同上文件现读 |

⇒ 截图里那一栏今天画的"已完成 09:41／回复中 10:18／已拒绝 昨天 16:05"**在 Go 侧一枚都取不到**（前端是示意数据）。

## AC（每格都要答"这一发在未修码上响不响"）

- [ ] **AC#1 先答"子代理这一维到底要不要进产品"**：⚠ 这一格不是技术题——`180-c1＋182-c1` 现量它是**规格真空**，加它＝新需求（牵 D22 双角色／D31 并发审批／D38 线程模型三节）。⇒ 只读普查：要动哪几节规格、有没有更小的等价维（例："这一轮有几个后台任务在跑"而不是"有几个子代理"）。**结论摆给 owner 拍，不许自己定。**
- [ ] **AC#2 后台任务状态这一维先落地**（它不需要新需求，名册已在、缺字段）：给 `Record` 补状态维（枚举要**逐字对齐 D43 状态机权威转移表**，不许自造），并送进快照。⚠ 判据要钉"值来自真名册"，不许硬编码。
- [ ] **AC#3 反向判据**：状态维**不许**变成"面板侧可写"（面板只许显示＋发起请求，R20／票 92 口径）；补一枚常驻判据钉住"这一维没有面板写腿"。
- [ ] **AC#4 与票 145 分账**：快照载体扩张归票 145（`SPEC-08` 那一节的字段表），本票只交"源"；结题时逐条写明谁欠谁。
- [ ] **AC#5 契约轴**：`docs/PLAN.md`（D31／D38／**D43 转移表**）、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden 一字节不许动；**C17 白名单一枚不许加**。
- [ ] **AC#6 门禁**：逐包 `go test -count=1 ./internal/tools/ ./internal/agent/ ./internal/panel/`；`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**）；`gate-clauses.sh` **比红腿名册**（在册只 `G6neg`）；⚠ `./internal/panel/` 今天有 **2 枚已知红**（`design/assets/tokens.css` 被别家会话删了未 staged），照实记不许当绿、不许顺手修。

## 本票**不**解决
不画界面（`frontend/**` 归另一枚会话）；不做"未提交枚数／diff／文件树"（＝票 189／190）；不做面板终端与浏览器（＝票 191／192，两枚都要人先批准）。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权补票（`A373`）＋`180-c1＋182-c1` 名册里"没票认领"七处中的两堆。未派。
- 09-28 18:1x `188-r1`（写腿，只做 `AC#2`，派单 §E）：**`AC#2` 停手上报，零产码**；证据件 `docs/evidence/s1/188-task-state-r1.md`。三处前提冲突逐条带 `file:line`：
  ① **写面与票面 现量 不重合**——本票 现量 与 `AC#2` 点名的 `Record` 是 `*tools.TaskRoster` 的方法（`internal/tools/task.go:123`，票面行号段 `110-152` 对得上），缺字段的是 `TaskOutput`（`task.go:93-102`，只有 `Text:97`／`ArtifactPath:101`），而派单给的写面是 `internal/memory/**`，那包里**没有名为 `Record` 的载体**（只有 `dao_providerhealth.go:131 RecordProviderError`、`privacy.go:194 Records`）；
  ② **memory 侧不是"缺字段"**：`TaskLog.State` 已在（`internal/memory/models.go:62`；schema `state TEXT NOT NULL` 且刻意无 CHECK，`schema.go:57`／`:14-17`；DAO 只校验非空，`models.go:185-197`），生产者两枚都在写面外——`internal/agent/loop.go:983-998 taskLogState`（`done`/`stuck`/`cancelled`/`interrupted`/`error`，注释自称 free-text）与 `cmd/wisp/run.go:651-670`；这套词**一枚都不在 D43 的 20 个名字里**（那 20 枚的仓内逐字拷贝＝`internal/statemachine/states.go:11-31`，标题行 `docs/PLAN.md:3055`、40 条转移 `:3061-3100`，本程未动一字）；
  ③ **撞钉预检命中即红**：把 `task_log.state` 改成 D43 白名单会让既有断言变红——`internal/memory/dao_test.go:253,262,265,282,307`（用 `running`/`succeeded`/`interrupted`）、`internal/memory/privacy_test.go:28,34`、`internal/agent/forensics_test.go:128-132`（钉 `cancelled`）、`internal/agent/loop_golden_test.go:315-320`（钉 `done`）；两枚在写面外、两枚在写面内，一律不自放宽。
  需要的具名解冻＝编排者二选一：(i) 解冻 `internal/tools/task.go`＋`task_backfill.go`（只到"给 `TaskOutput` 新增状态字段"为止），或 (ii) 改判载体为 `internal/memory.TaskLog` 并同步解冻 `internal/agent/loop.go` 那枚映射与 `cmd/wisp/run.go` 调用点；另外"任务态 -> D43 名字"那张映射表是**规格真空**（D43 讲的是会话/球的 20 态，没有一行讲后台任务），要落 `AC#2` 先定它归谁。
  `AC#3` 反向判据也**未落**：现量 `internal/panel/` 非测试码里 `task_log`／`TaskLog`／`TaskOutput`／background task **零命中**，载体没定就只能交死钉（无正控）；等 (i)/(ii) 定案后同批落，并自带"种一发面板写腿必响"的正控。`AC#1` 子代理那一维**没碰**（`A389` 已摆 owner）。
  门禁读数：`d22scan` **clean**（`ban #8 internal/=441` 枚文件，票面基线 433 已过期＝别家加文件，涨不等于漂移）、`gate-clauses.sh` **BAD 名册只 `G6neg`**、`gofumpt -l internal/memory/` 空、`go test ./internal/memory/` ok 12.810s；`./internal/tools/`／`./internal/agent/`／`./internal/panel/`／`cmd/wisp` **未采读数**（在别人的飞件名下）。`flip-declaration.sh` 未跑。AC 框一枚未勾。
- 09-28 18:0x 编排者收 `188-r1`（五枚 commit，零产码、停手上报）：**载体选 (i)**＝解冻 `internal/tools/task.go`＋`task_backfill.go`，射程只到给 `TaskOutput` 加一枚状态字段，值**只能取自 `internal/statemachine/states.go:11-31`**（D43 逐字拷贝），不许新造、不许塞 `done/cancelled`；**不改判成 `TaskLog.State`**（那要动两处生产写入者＝改既有语义）。⚠ 它挖出的两套词表并存另立**票 196**；本票 `AC#2` 期间**不许与 196 互相顺手统一**。账 `A394`。
