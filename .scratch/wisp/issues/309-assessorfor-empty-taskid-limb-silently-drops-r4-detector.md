# 票 309 — `Bridge.assessorFor` 的三枚前提零仪器：`taskID == ""` 那一肢把 R4 探测器**静默摘掉**，而同一种"宿主没给任务 id"在别处是**拒绝**或**自己命名**

**立票**：2026-10-11 11:5x 编排者（来路＝只读腿 `taskid-a1`（笔 `93643a29`，件 `probes/orch/taskid-a1/**`）答掉我自己登记在 `A838` §3 的那枚欠读数；八问读数我自己逐枚复跑，见台账 `A842`）
**性质**：★**今天不可达**（生产构造点收敛到 `newTaskID()`，恒非空）⇒ 本票⛔ 缺陷、本票＝**缺结构保证＋缺仪器＋同形三种处置不一致**。低利害内务，⛔ 不上机主清单。
**Status:** 未开工（六格全未勾）。⛔ 零翻框（翻框归非实现者验收程）、⛔ 零 push。

## 现量（每条带尺名＋射程；⛔ 引用前先重跑，行号是 10-11 11:5x 快照，一律配内容锚）

- 那枚守卫逐字（尺＝`sed -n '924,932p' internal/tools/bridge.go`，工作树＝HEAD）＝
  `func (b *Bridge) assessorFor(taskID string) *risk.RiskAssessor {` / `if b.injected || b.prov == nil || taskID == "" {` / `return b.assess`。
  ⇒ 三枚前提任一成立就**整枚摘掉 R4 探测器**（`:931 WithTaintDetector(b.prov.Detector(taskID))` 那一支不走）。
- 它的文档块（`:920-923`，锚 `assessorFor returns the assessor whose R4 input is bound to this task's`）**只讲按 scope 绑定**，
  ★对 `taskID == ""` 那一肢**零理由注释**（尺＝`sed -n '918,926p' internal/tools/bridge.go`）。
- 消费者面（尺＝`grep -rn 'assessorFor' internal cmd --include=*.go`，滤 `_test.go`）＝**真调用点 1 枚**（`bridge.go:297` 锚 `verdict := b.assessorFor(req.TaskID)`）；
  其余 4 枚命中全是注释（`assessor.go:30`、`panel_assets.go:60`/`:227`、定义自己）。
- 空串可达性（尺＝`grep -rn 'ToolRequest{' internal cmd --include=*.go`，滤测试）＝生产构造点**全仓 1 枚**＝`internal/agent/loop.go:675`，实参收敛到 `newTaskID()`（`:1153`，返回 `8-4-4-4-12` 十六进制＝恒 36 字符；熵失败那支退化成时钟派生**仍非空**）⇒ **今天传不进空串**。
- ★同一种"宿主没给任务 id"，三处三种处置（逐字）：
  ① `bridge.go:925` ＝**静默摘掉探测器**、零注释；
  ② `cmd/wisp/panel_assets.go:164` 注释逐字 `a CLI probe has no task, so it names one`＝**自己命名一枚并写明故意**；
  ③ `internal/tools/subagent_197.go:263`（锚 `parentID := TaskID(ctx)` 后 `if parentID == "" {`）＝**fail-closed 拒绝**，文案逐字"不知道父任务是谁，就没法登记父子关系…"。
- 仪器面（尺＝`grep -rn --include=*_test.go 'assessorFor' internal cmd`）＝点名 `assessorFor` 的用例 **0 枚**（唯一命中 `panel_assets_143_test.go:7` 是注释行）。
  ⚠ 但⛔ 等于"从没被经过"：`internal/tools/failclosed_236_teeth_test.go:244`（锚 `func Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID`）在 `:288` 逐字 `TaskID: "", CorrelationID: ""` 真把空 id 喂进 `Execute`；
  它的断言全是 `res.IsError`／字面量／`h.roster.Count()`，⛔ 一行讲检测器挂⛔ 挂 ⇒ 准确口径＝**被经过、未被观测**（⛔ 写"零仪器经过"，那一说过头）。
- 摘掉守卫的反面地板（⛔ 本票走法）：`internal/risk/provenance.go:727`（锚 `if unbound {`）对未开 scope 无条件回 `Hit{… SrcTool: SrcUnboundScope}, true`
  ⇒ 少了 `taskID == ""` 那一肢，空 id 会撞成**每一次写都报**＝误报风暴。所以"删掉那一肢"⛔ 是更安全的方向，⛔ 有人顺手这么"从严"。
- ⚠ 注释里的行号已腐烂（尺＝`sed -n '869,872p' internal/tools/bridge.go` ⇄ `grep -n 'func (l \*Loop) Run\|func (l \*Loop) RunAsync\|newTaskID()' internal/agent/loop.go`）：那段逐字写 `internal/agent/loop.go:332/:321`，本轮现量真身＝`RunAsync` 起 **`:328`**（铸造点 `:329`）／`Run` 起 **`:339`**（`:340`）⇒ 两枚号全漂。
- ⛔ 属本票那一问（另记，⛔ 顺手）：`b.injected` 那一肢同样不挂探测器，但生产面 `Assessor:` 赋值点＝**0 枚**（尺＝`grep -rn 'Assessor:' internal cmd --include=*.go`，滤测试＝1 命中且那是注释）⇒ `injected` 今天恒 false，归装配那族另裁。

## 验收标准（六格；★次序硬门＝`AC#0` 未裁形之前 ⛔ 派落地腿）

- [x] **`AC#0` 三处同形处置的归口与选形（只读，判语交编排者裁）**：把上面那三处逐字名册核实或更正，并给出"要不要对齐成同一种处置"的两形代价：ⓐ**入口拒绝**（照 `subagent_197.go` 那支：空 id 就报错、不猜）ⓑ**写明故意**（照 `panel_assets.go:164` 那支：静默降级留在原地＋注释具名"没有任务 ⇒ 不挂探测器 ⇒ 这一支少报"）。每形要写：会动哪几枚文件／会不会打红既有钉（点名）／用户侧表现。★两形都⛔ 是"删掉那一肢"（见上面那枚地板）。
- [ ] **`AC#1` 给三枚前提装上仪器**：判据＝**恒真形必须响**——把 `taskID == ""` 那一肢改成"空 id 也挂探测器"或整枚摘掉守卫时，指名用例必须红；反过来探测器没挂上时也要红。⛔ 只交"改⇒红"不交"换反形也红"＝那格没交。写面＝`internal/tools/*_test.go`；⛔ 产码语义（那归 `AC#0` 裁完的落地程）。
- [ ] **`AC#2` 结构保证那一问**：今天靠"生产构造点只有 1 枚"保证非空＝**接线事实、⛔ 结构保证**（`bridge.go:872-875` 已点名三类未来候选者：自带 id 的面板/球宿主、重试包装器、跨轮次带 id 的定时任务）。本格只回答并落地一件最小事：**将来任一 caller 手搭 `ToolRequest` 递进空 `TaskID`，要让它在盘上可核地响**（ⓐ 入口断言／拒绝，或 ⓑ `AC#1` 那枚仪器覆盖到）。⛔ 本票新造一条生产可达路径。
- [ ] **`AC#3` 文案一致性**：`bridge.go:869-876` 那段引的 `loop.go:332/:321` 行号已漂 ⇒ 换成**内容锚**；同段那句"Today no production code dispatches on the bridge except loop.go"要带尺重述（现量＝构造点 1 枚）。⛔ 改这段里的任何**契约声称**文字（只改行号与枚数）。
- [ ] **`AC#4` 门禁与越界**：`GOFLAGS= go build ./...`／`go vet ./internal/tools/`／`sh scripts/d22scan.sh` 各 rc=0 且**门禁件自己落一行 `rc=N`**（⛔ 0 字节件＝那格没交）；格式名册两把并排（工作树⇄HEAD blob）新增 0 枚并写射程；每笔 commit 带显式 pathspec，commit 前 `git diff --cached --name-only` 数一遍；写面⛔ 出 `internal/tools/**`＋`probes/309/**`＋本票末尾一条 append-only bullet（1 加／0 删）。
- [ ] **`AC#5` 三本账对账（⛔ 同一件事记两次）**：具名写清本票与 **票 307**（已 `-done`，那枚是折法不对称、成因⛔ 同）、**票 40**（面板 Security 页零消费者，归口见 `A838` §3）、**票 197**（`subagent_197.go` 那支拒绝形的出处）、**票 242**（correlation id 改每调用一枚，正是 `:920-923` 那句"per scope, hence per task"的前因）四家的分工；台账 `A838` §3 那枚"新抓具名欠读数"⇒ **由 `A842` 结掉**，⛔ 在本票再记一遍。

## 边界与已知禁区（⛔ 派单要逐字带上）

- ⛔ 动风险层判级（`internal/risk/**`）：探测器没挂是**装配/入口**问题，⛔ 靠改判级"从严"来修（见那枚 `SrcUnboundScope` 误报风暴地板）。
- ⛔ 动 `internal/observe/thresholds.go`／`testdata/golden`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件（`internal/panel/tokens_fourway_test.go`·`internal/panel/l2_grant_boundary_test.go`·`internal/perm/ticket90_persist_test.go`）。
- ⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP` 换绿；⛔ 把注释里的"今天没有生产 caller"当成永久保证（那是会被别人一行 commit 缩短的现读句）。
- ⚠ 台面三根轴照旧：ⓐ母仓⇄干净 clone 颜色⛔ 可比；ⓑ母仓那枚未跟踪 `frontend/dist/index.html` 自己会漂；⏂导出树⛔ 带 sherpa/build 两枚 DLL 目录 ⇒ 跑 `cmd/wisp` 一律 PATH 写母仓绝对路径，且第一发先断言具名用例真发出 `--- PASS/FAIL/SKIP`。
- ⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree/checkout；⛔ push（子代理只 commit）；临时件只建不删、从一开始就落自家 `logs/`；证据件⛔ 叫 `.out`。

## 排程

- 按在 **`72-r1`**（写腿，`internal/risk/*_test.go`）交完之后：本票 `AC#1` 的写面是 `internal/tools/*_test.go`，⛔ 与它抢同一枚 Go 编译/测试面的整包名册。
- next＝`AC#0` 只读普查（`309-a1`）⇒ 我裁形 ⓐ／ⓑ ⇒ `309-r1` 落地 ⇒ `309-v1` 非实现者终裁。

## `AC#0` 裁语（2026-10-11 12:2x +08，编排者；普查腿＝`309-a1` 三笔 `2a37d612`／`27059960`／`ff301daf`，件 `probes/309/a1/**`；判语者⛔ 实现者）

★**裁＝形 ⓑ（写明故意）；`AC#2` 的保证落"形 B 枚数钉"；⛔ 形 ⓐ（桥入口拒绝）。撤销口令＝「309 改形 ⓐ」。**

- **⛔ ⓐ 的三条理由**（每条指着我自己复跑过的尺，正文＝台账 `A844` §1／§2）：①**射程比病灶宽**＝ⓐ 让**所有**工具调用在空 id 时整枚失败，而要修的只有 R4 少报那一支（机械短路点现量＝`internal/risk/rules_gateway.go:102` 逐字 `if ctx.taint == nil {`）；②**必红两枚既有仪器**＝`internal/tools/failclosed_236_teeth_test.go:122` 与 `:257`（载体逐字＝`:118 x.cancel(t, "", target)` → `task_cancel_221_legs_test.go:116-119` 的 `TaskID: callerID`，与 `:288 TaskID: "", CorrelationID: ""` 直喂 `Execute`）——换期望值＝动票 236 的仪器本体，改载体＝③④ 两支必须同批补仪器否则覆盖搬家式消失；③**遮蔽教学文案**＝③④ 那两句逐字在教"派生要有父／停人要知道是谁"，被桥上游一句"没有任务 id"盖掉。
- **ⓑ 的账（不许写成"更严所以更好"）**：行为逐字节不变 ⇒ **对空 id 的实际防御＝零**（用户／模型／审计行零变化）。这形里"响"只有两枚来源＝`AC#1` 那枚"恒真形必须响"的仪器＋下面那枚枚数钉；腿在 `21-shape-b-cost.md:20` 具名过：**注释先写、仪器后补＝`AC#1` 那格没交** ⇒ `309-r1` 必须同片写、口径一致。
- **`AC#2` 落点定形**＝一枚 `internal/tools/*_test.go` 里的**枚数钉**：`go/parser` 读生产源码，断言"`agent.ToolRequest{` 在 `internal/agent` 之外的生产构造点枚数＝**0**"。现量基线两把独立尺：枚数尺＝全仓 `ToolRequest{` **60** 枚＝`_test.go` 内 59／外 **1**（那枚＝`internal/agent/loop.go:675`，被排除项吃掉）；门尺＝票 154 `gate-clauses.sh` 的 **G1** 那格登记 `want_n 0` 今日对得上。★**范本⛔ 去 `cmd/wisp` 找，本包就有两枚现成的**：`paths_twocontainments_252_r2_test.go:195-206`（`os.ReadDir(".")` 排除 `_test.go` 后 `parser.ParseFile`）与 `task_state_188_test.go:243-265`（`WalkDir`＋把 `examined` 当返回值，注释逐字 *"a scan whose denominator is zero reports clean for a tree it never read"*）⇒ **新钉必须自带分母**，⛔ 交一枚"扫了 0 个文件然后报绿"。
- **★两形都没防住的形状（下一程不许以为裁完就闭）**：ⓐ／ⓑ 都⛔ 回答"未来那条路真递进空 id 之后 R4 响不响"。ⓐ 把那次调用整个掐掉，ⓑ 只保证"有人写它时会被钉叫住"。本票真正的"少报有人管"＝`AC#1` 那一格，`AC#0` 裁形**⛔ 闭合 `AC#1`／`AC#2` 之外的任何格**。
- **未计价的第三形（登记，本轮⛔ 采纳）**：ⓐ 还有一枚窄变体＝"只对会用到探测器的写类工具在空 id 时拒"。票面与腿都⛔ 给它算过价（必动文件／必红枚数／与 `internal/risk` 判级面的接缝），要真走得先再来一枚普查。

### 对本票「现量」两节的更正（★只追加，⛔ 改上面任何一句原话；两处都是**我写错的**，被抓出来记在 `A844` §1）

1. **"三处三种处置"＝数目不足**。现量名册是**四处**：①`bridge.go:925` 静默摘探测（⛔ 注释）／②`cmd/wisp/panel_assets.go:164` 自己命名并写明故意／③`internal/tools/subagent_197.go:263-266` 拒／**④`internal/tools/task.go:708-711` 拒**（锚 `caller := TaskID(ctx)` 后 `if caller == "" {`，文案逐字"这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）"，与 ③ 共享同一枚字面量前缀）。⇒ 准确口径＝**四处同形、两拒一自命名一静默**，"静默那枚⛔ 仪器"这一句⛔ 变。
2. **那枚反面地板我写宽了**。原文写"对未开 scope **无条件**回 `Hit{… SrcTool: SrcUnboundScope}, true"。现量真身＝`internal/risk/provenance.go:694-710` 的 `scopeMarks`：`reg != nil` 时一律返回 `(nil,false)`（`:701-702`），只有**该 id 没登记ⓔ另有 scope 手里有 marks** 才置 `unbound=true`（`:704-709`）⇒ 单任务、全仓无 marks 的那一发⛔ 报。**方向结论⛔ 变**（"删掉 `taskID == ""` 那一肢"仍⛔ 是更安全的方向，因为它把"漏报"换成"看情形的误报"），但措辞该按上面这句读。⚠ 另补一枚我漏的消费者＝`:828`（`CheckText` 同形），本票原只指了 `:727`。
3. **口径小账**：`assessorFor` 名册尺现量 **7 命中**＝真调用点 **1**（`bridge.go:297`）＋函数声明 **1**（`:924`）＋**注释行 5**（`:920`／`assessor.go:30`／`panel_assets.go:60`／`panel_assets.go:227`／`panel_assets_143_test.go:7`）；原文那句"其余 4 枚命中全是注释"把文档块与声明混成了一枚，按 **5** 读。"点名它的测试只有注释那一枚"⛔ 变。

### 排程更新（追加，⛔ 改上面「排程」那两行）

- `AC#0` 既裁 ⇒ **`309-r1`（落地腿）现在仍按住**，具名理由＝它要动 `internal/tools/**` 并跑 `go test ./internal/tools/`，而 `298-v1`（票 298 非实现者验收）正在取验收色、`308-a1`（票 308 只读普查）正在读 `internal/risk` 的台面颜色；两枚交完才派。⚠ 别把"裁完形"读成"可以开工"。
- `309-r1` 的写面与判据按本节定：Ⓑ 的注释（`bridge.go:920-923` 追加＋`:925` 上方具名）＋`AC#1` 那枚恒真形仪器＋`AC#2` 那枚枚数钉（自带分母）；⛔ 动 `:925` 的布尔式、⛔ 删那一肢、⛔ 动 `internal/risk/**`、⛔ 改 `lit236*` 两枚常数的值。
- 腿具名的 7 枚欠读数归口逐枚写在 `A844` §3（其中 **U-3 已由我这发结掉**、**U-4 结掉一半**＝G1b 那把门尺今日会响，成因⛔ 本票，归票 154 那一族，⛔ 在 309 里重登记 `want_n`）。

### 翻框那一问（★本轮翻 `AC#0` 一枚，理由具名）

- ★翻＝**`AC#0` 一格**（尺后＝1 勾／5 未勾，尺＝`grep -cE '^[[:space:]]*- \[x\]'` ⇄ `- [ ]`）。理由三条逐字：①本格票面自写"只读，判语交编排者裁"，它的交付物＝**核实或更正名册＋两形代价**，两半都已交（普查腿交、我逐枚复跑）；②本格⛔ 产码＝**没有"实现者自证"可攻**，`SPEC-12 §4.3`「裁决者≠实现者」在这格的作用对象是"谁来判定实现"，⛔ 是"谁来签收一句裁语"；③本仓同一形制已有三枚先例（票 228／票 257／票 258 的 `AC#0` 都是"收只读普查腿 ⇒ 编排者裁形并翻那一格"）。
- ⚠ **其余五格一枚都没翻**＝`AC#1`／`AC#2`／`AC#3`／`AC#4`／`AC#5` 全归 **`309-v1`（非实现者终裁）**。⛔ 把"⛔ 动框"这条规矩读成"已裁完也留成过期框"——那正是台账抱怨的"Status 行整批过期"那一族病。
- ⚠ 也别把这句读成"等验收腿才能开工"：`309-r1` 的开工门已开（形 Ⓑ 已定，口令「309 改形 ⓐ」）；它现在按住的原因⛔ 是缺裁语＝纯粹是排面冲突（见上节那条）。
- ⚠ 票头那行「**Status:** 未开工（六格全未勾）。⛔ 零翻框」随本笔**过期**＝现量 1 勾／5 未勾（尺＝`grep -cE '^[[:space:]]*- \[x\]'` ⇄ `- [ ]`）。⛔ 改原句（工单池 `.scratch/wisp/issues/README.md:201` 原话逐字"别因为看到 `-done` 就以为它零残余"那族病，⛔ 由我再造一枚过期 Status 行），以本节为准；`-done` 后缀仍⛔ 加（五格未闭合）。
- **落地腿 `309-r1` 交件（2026-10-11 13:1x +08，实现者；⛔ 翻框＝五格全归 `309-v1`）**：按 `AC#0` 裁定的**形 Ⓑ** 落地——`bridge.go` `assessorFor` 文档块追加"写明故意＋UNDER-REPORTS"注释、`:925` 肢上方具名一行（⛔ 改布尔式／⛔ 删肢／⛔ 动 `internal/risk`，行为逐字节不变）；`AC#1`＝`internal/tools/assessorfor_taint_instrument_309_test.go::Test309AssessorForEmptyTaskIDDropsTaintDetector`（恒真形两向都有牙：突变"摘 `|| taskID == ""` 肢"⇒空 id 腿红、"探测器改没挂"⇒非空腿红，改前绿/突变红/还原逐字节一致三件齐，读数 `probes/309/r1/logs/mutation-rigs.log`）；`AC#2`＝`internal/tools/toolrequest_nail_309_test.go::Test309ToolRequestProductionConstructionNail`（枚数钉自带分母，现量 `internal/agent` 之外 `agent.ToolRequest{` 生产构造点＝0、分母 245，内置正控证尺⛔ 瞎，Rig 证"加一枚即红"读数 `logs/nail-rig.log`）；`AC#3`＝`bridge.go:861`/`:871` 两处绝对号换内容锚＋`877` 那句带尺现量重述（U-5 现读：`loop.go:366` 漂⇒真身＝`defer revokeAdmission()`@`:379`）；`AC#4`＝`go build`/`go vet ./internal/tools/`/`sh scripts/d22scan.sh` 各 rc=0（d22scan 一度红＝枚数钉 `t.Fatalf` 字面量里的 `≈`U+2248 落禁带、非注释豁免命中，改"接近零"后绿，⛔ 动 allowlist）＋格式名册 `internal/tools` 两尺 0↔0＋逐名作差纯新增（209→211 顶层 PASS、0 FAIL/0 SKIP、两枚 236 空 id 载体照绿）＋本轮自抓自纠：`bridge.go` 被我编辑成 CRLF（`CR=1313`，⊄ `.gitattributes *.go eol=lf`）已归一 LF；`AC#5` 对账 307/40/197/242 见 `probes/309/r1/90-final.md` 末节。件＝`probes/309/r1/{00,10,20,30,40,90}.md`＋`logs/`。commit `1cd57ec1`（锚）／`16e86bae`（产码+测试）／本笔（证据件+logs＋此 bullet）。
- **非实现者终裁 `309-v1` 判语（2026-10-11 13:5x +08，裁决者⛔ 实现者；★本 bullet ⛔ 翻任何框，五格翻不翻归编排者复跑后定）**：`AC#1`＝**附条件成立**（四形突变本腿逐形 LANDS 自证＋逐字红：摘 `|| taskID == ""`⇒`:64` 红、`WithTaintDetector(nil)`⇒`:54` 红（`level=L0 rules=[]`＝⛔ 恒真比较的反证）、整枚守卫 off⇒`:64` 红、肢取反⇒双红，每形 `RESTORE cmp: byte-identical OK`＋`sha256sum cb86eaa7…`；⛔ 空集换绿；**条件**＝标题"三枚前提"只装了一肢——摘 `b.prov == nil`／摘 `b.injected` 两发**各 rc=0**＝仪器对那两肢瞎，且 `bridge.go:947` "Who rings if this limb's shape changes"⛔ 全称）；`AC#2`＝**附条件成立**（分母真算：钉自打印 245＝第三方尺 `269−24`；"加一枚即红"在钉**真扫的** `bridge.go` 上复现，constructions 0→1 而分母⛔ 动；**条件**＝格文"任一 caller"⛔ 闭合——声明＋赋值、别名两形实测 rc=0，`internal/agent` 整包 24 枚源⛔ 解析、`tools/` 13 枚源在 roots 外，而这族形状洞**⛔ 新发现**＝票 154 验收程合成树已量过五形（`probes/154-accept/ac1-g1-blindspot-synthetic.txt`）⇒ ⛔ 另立新账）；`AC#3`＝**附条件成立**（最硬两把＝± 行 35 枚里**非 `//` 注释行 0** ＋ 两版 blob 剥注释后 `cmp` **IDENTICAL（代码行 768⇄768）**⇒"行为逐字节不变"由尺坐实；6 枚新锚今日逐枚在盘上（`defer revokeAdmission()`=`loop.go:379`、`ruleTaint`+`ctx.taint == nil`=`rules_gateway.go:101/102`、`scopeMarks`=`provenance.go:694`＝正是追加更正②的口径）；**条件**＝"带尺重述"那半句没过——`bridge.go:879-881` 自称 1、同一把尺实测 **3**，多出的两行＝**该注释自己**（剥行首 `//` 后＝1，结论⛔ 变）；另★我派单那句"`git diff -w --numstat` expect 0"**是错前提**，实测 `32 3`，已按实测量报、⛔ 换尺凑 0）；`AC#4`＝**成立**（三枚门禁本腿各一发即 `rc_build=0`/`rc_vet=0`/`rc_d22scan=0` 且各自落行、件⛔ 0 字节；eol 三把尺复量那次自抓自纠＝`CR=0`／`i/lf w/lf attr/text eol=lf`／`.gitattributes:4`；格式名册四档**归 309-r1 新增＝0**，`gofumpt` 在 `$(go env GOPATH)/bin` 在场、⛔ 免检，blob 侧那 1 枚＝票 298 `AC#5` 已承接；越写面＝0）；`AC#5`＝**成立**（四家按**它们自己的正文**分开，`A838` §3⛔ 重记✅、`want_n`⛔ 揽✅；措辞更正一枚⛔ 阻格＝③ 拒绝形的文本出处是**产码** `subagent_197.go:263-266`＋文件名对票号命名，票 197 正文 `grep "没有宿主给的任务 id"`＝**0 命中**）。**★承重那句**：Ⓑ 之后**对空 id 的实际防御＝零**，三把尺坐实（代码行 `cmp` 同一／两枚真喂空 id 的 236 载体在**三张台面**全 PASS／整包名册**纯新增 2 枚、删 0、SKIP 0**，两棵导出树各自余下那枚红＝`TestTicket252R1AlignmentAddsNoAuthorization`，改前改后同一枚＝台面缺件⛔ 归 309）；盘上只有**两枚跳线**，⛔ 一格⛔ 读成"以后有人递空 id 会有人管"。**具名欠账 10 枚**（`Q-309-v1-prov-nil-limb-blind`／`-who-rings-overclaim`／`-other-gates-uninstrumented`／`-nail-misses-decl-and-alias-shapes`／`-nail-excludes-whole-agent-package`／`-nail-roots-miss-tools-dir`／`-nail-control-not-exercisable`／`-selfref-ruler`／`-stale-count-comment`／`-anchor-debt-154`，最后一枚⛔ 本票账＝归票 154/151 族；逐枚内容与归口见件）。**撤销口令**＝逐枚「**309-v1 撤 `AC#n`**」（编排者复跑对不上即撤，⛔ 一把撤整表）；改判形 ⓐ 仍用既有口令「309 改形 ⓐ」，届时本表五格判语随裁形整批作废。件＝`probes/309/v1/{00,10,20,30,40,50,60,90}.md`＋`logs/`（含可逐字重跑的 Rig 脚本）＋裁决表 `docs/evidence/s1/309-v1-verdict.md`；台面 HEAD `0d2260eb`，收尾 `git status --porcelain -- internal cmd`＝**0 行**；⛔ push／⛔ 翻框／⛔ 动过一枚断言，瞬态突变只打 `bridge.go` 且每形当场还原。
