# 票 309 — `309-v2` 非实现者终裁表（`AC#1`／`AC#2`／`AC#3`）

裁决者＝`309-v2`（⛔ 实现者；`309-r1`/`309-r2` 的码与本表零关系，本表只读盘＋跑尺）。
锚点＝`f6067e5d`（`git rev-parse --short HEAD` 起手现跑＝`f6067e5d`，与派单 §1-1 要求的号**相等**）；起手
`git -c core.quotePath=false diff --cached --name-only | wc -l`＝**0**（派单时刻也是 0；我没看到 staged 别家路径）。
本轮**⛔ 翻任何框**（起手尺＝`grep -cE '^[[:space:]]*- \[x\]'`＝**3** 勾／`- [ ]`＝**3** 未勾，收尾同一对数＝3／3）。
台面＝工作树＋两棵仓外导出树（`git archive f6067e5d | tar -x`，⛔ `cp` 别的锚 ⇒ **单面**，`Q-309-r2-tree-rigs-mixed-face` 对我这些读数已结）。

> ★**承重那句先写在最前面**（⛔ 让下一位以为三格成立＝这件事闭了）：形 Ⓑ 之后**对空 `taskID` 的实际防御仍然＝零**。
> 本表能判的只是"有人改那枚肢／有人新写一条构造点时会不会响"。盘上仍然只有两枚跳线：一枚定向仪器
> （只认"挂没挂探测器"这一件事的两向）＋一枚枚数钉（只认三形写法、且 roots 不被谁钉着）。

## 判格（每格：判语／最强那一行凭据／尺名＋射程＋量的还是算的＋provenance）

| 格 | 判语 | 最强凭据（可逐字重跑） | 尺（名／射程／量或算／provenance） |
|---|---|---|---|
| **`AC#1`** 三枚前提装仪器（判据＝恒真形必须响） | **成立**（`309-v1` 那枚"只装了一肢"的条件按编排者 §4.4 的口径＝**已闭合**；残余另记 `Q-309-v2-1`，⛔ 本格的欠） | 四形突变各**红**、pristine **绿**，逐形当场还原：摘 `\|\| taskID == ""`⇒`--- FAIL: Test309AssessorForEmptyTaskIDDropsTaintDetector`；`WithTaintDetector(b.prov.Detector(taskID))`⇒`WithTaintDetector(nil)`⇒同一枚 FAIL；整枚守卫 off（`if false {`）⇒FAIL；肢取反（`taskID != ""`）⇒FAIL。件＝`probes/309/v2/logs/rigA-instrument.log`＋`rigA-{0-pristine,a1,a2,a3,a4}.log`；还原证明＝每形 `RESTORE cmp: byte-identical OK sha=995c23d2ef6ea16c… match=YES`（⇄ `git show f6067e5d:internal/tools/bridge.go \| sha256sum` 同一枚 sha，见 `rigD-callsite.log` 末行） | 尺＝`go test ./internal/tools/ -count=1 -run Test309AssessorForEmptyTaskIDDropsTaintDetector -v`／射程＝包内点名那一枚用例／**量的**／provenance＝**仓外导出树 A**（`git archive f6067e5d`）。"世界"在不在＝读码尺 `sed -n '32,55p' internal/tools/helpers_test.go`（`fsBridgeWith` 恒 `Provenance: prov`、恒不写 `Assessor:`）⇒ `b.prov != nil && !b.injected` 的夹具真存在，两肢只差 `taskID` 的取值 ⇒ **非缺世界**（`rigA-0-pristine.log` 里 ① 腿 `--- PASS` 就是"挂上了"那一侧真发生过） |
| **`AC#2`** 结构保证（枚数钉） | **成立**（`309-v1` 的三条条件全部由**我这一发**复现：整包 `internal/agent` 不再被跳／声明＋赋值形／别名形／roots 加 `tools`＋`scripts`；新抓的四处残余是**新欠**不是旧条件，见 §欠账） | 名册 1→2 且具名新件 ⇒ 红：`internal/agent/rig309v2bare.go:3(same-package-bare-literal)`（B2，`rc=1`）、`cmd/llmrecord/rig309v2e.go:5(var-decl-no-initializer)`（E1，`rc=1`）、`internal/agent/rig309v2e.go:5(alias-literal)`（E2，`rc=1`）、`tools/rig309v2root/rig.go:5`（B3）、`scripts/rig309v2root/rig.go:5`（B4）；**移除**也响：把 `loop.go:675` 那一枚改写成 `req := *new(ToolRequest)`＋逐字段赋值（编译得过、语义等价）⇒ 钉红在 `名册=(空名册) … 数到 0 枚`（B5，`rc=1`）＝判据是 `==1` ⛔ `<=1`；**往宽里挪一档**：`wantTotal 1→2` 打在⛔ 变的树上⇒**红**（B6，`rc=1`）＝谁悄悄放宽期望值会被 CI 当场叫住。pristine 正控＝`--- PASS … 构造点=1 … 分母=293 枚非测试源（每 root=cmd=40 internal=229 scripts=11 tools=13）`（B0/C0） | 尺＝`go test ./internal/tools/ -count=1 -run Test309ToolRequestProductionConstructionNail -v`（钉自己＝`go/parser` 读源码，射程＝roots `cmd internal tools scripts`、滤 `_test.go` 按**路径后缀**）／**量的**（分母是解析计数，293 另有一把**算的**第三方尺：`find <root> -name '*.go' -not -name '*_test.go' \| wc -l`＝40/229/13/11，件 `probes/309/v2/logs/tree-path.txt` 旁 `s14-outside-roots.txt`）／provenance＝**仓外导出树 B** |
| **`AC#3`** 文案一致性 | **附条件成立**——**一条**条件，且是本轮新抓的（⛔ 在 `309-v1` 的账上）：`bridge.go:959-963` 那句全称"counts **every** production construction point … so a second one rings **wherever** it is written"**量上不成立**（`Q-309-v2-3`）；同段其余句句实测对得上 | 反证＝**B14**：往 `internal/tools/` 放一枚真的第二构造点 `return new(agent.ToolRequest)`（钉头 `:56` 自己列的漏形之一）⇒ 顶层 **`--- PASS`、构造点仍＝1、分母漂到 294**＝有第二枚而⛔ 无人响（件 `rigB-b14.log`，判读 `rigB-nail.log` 第 "B14" 节）。同一把钉在**同树**里对 B1 的限定名字面量是红的 ⇒ 差的不是射程是**写法**。格文其余部分**过**：①内容锚全在盘上（`defer revokeAdmission()`／`id := newTaskID()`／`return l.run(ctx, newTaskID(), input)`／`if ctx.taint == nil {`／`func (p *Provenance) scopeMarks` 各 `grep -c`＝**1/1/1/1/1**，件 `s16-anchors-and-contract.txt`；点名段内 `loop.go:<数字>` 绝对号＝**0**）；②"带尺重述"那半句这次真带尺：注释写的那把逐字跑＝**71 行／文件名以 `_test.go` 结尾 69／非测试按文件名字段 2／再剥行首 `//`＝1**，且注释自己那句"Skip that last strip … yields 2 lines today"＝**2**（`internal/agent/loop.go:675` ＋ `bridge.go:882` 那一行注释本体），简写形 `\| grep -v _test.go` 也确是 **2**＝它警告的那枚假数成立（件 `r4-grep-rulers.txt`）；③"⛔ 改契约声称文字"＝`git diff bd124b2a f6067e5d -- internal/tools/bridge.go` 里带 `C19`/`C25`/`C26` 的增删行 **0**；④"行为逐字节不变"由尺坐实：该区间 `bridge.go` 改动行 **70 枚全是注释**、非注释改动行 **0**，两版 blob 剥整行注释后 `cmp` **IDENTICAL**（854⇄854，sha `496ed6d683e8e070…`，正是守卫那行逐字 1⇄1 的同一把尺） | 尺①＝`grep -rn 'ToolRequest{' internal cmd --include=*.go` ＋两支过滤（射程＝`internal`＋`cmd`，⛔ 四根）／**量的**／provenance＝**工作树**，且 `git status --porcelain -- internal cmd tools scripts`＝**0 行**⇒工作树⇄HEAD blob 同一枚（另用 `git show f6067e5d:` 复量守卫行）。尺②＝`grep -n 'loop\.go:[0-9]' internal/tools/bridge.go`＝**5 枚**，逐枚都在**票 154/151 那族段落**（`:113/:114/:230/:252/:1119`）⛔ `AC#3` 点名段内＝既有账 `Q-309-v1-anchor-debt-154`，编排者已裁"⛔ 本票账"，我照收不改判。尺③＝`go/parser` 那枚钉（射程四根）／量的／导出树 B |

## 四面进攻的正面答复（派单 §2）

**① 那枚钉有没有牙**——有，而且不是只有"改⇒红"：`==1` 使**加**（B1/B2/B3/B4/E1/E2）与**删**（B5）两向都红，**期望值放宽**当场红（B6）。
`Q-309-v1-nail-control-not-exercisable`（"内置正控⛔ 可测"）我这发用**扫描面**代证区分度并签收：把 matcher 打瞎（限定名 typo）⇒ 红在 `control "shape 1 qualified literal": matcher counted 0, want 1`（B11）；把 matcher 打成 catch-all ⇒ 红在 `control "guard: bare literal in a NON-agent package": matcher counted 1, want 0`（C1，`rigC-controls.log`）＝**两向都咬**，正控⛔ 恒真。

**② 恒真三面自查**
- (i) 永红形（`==1` 打在永不匹配的名册上）：**不会静默**——B11 实测红在控制那一支。
- (ii) 空循环形（遍历空集）：`roots=[]`⇒红在分母那一支（B9：`只解析了 0 枚非测试生产源（roots=[]）`）；root 改名成不存在的路径⇒红（B10）；`分母 floor 200→1`⇒**绿**（B7）＝floor 可被无声放宽，但 per-root"名册为空即 `t.Fatalf`"与 `==1` 两支仍在（这就是我对"往宽里挪一档为什么还守得住"的答复：`==1` 守得住，floor 守不住，见 `Q-309-v2-2`）。
- (iii) **缺世界形**——逐肢答：
  `taskID == ""`：世界**在**（夹具恒 `prov != nil`、恒 `injected == false`，两腿只差 id），四形红＝非恒真的实证。
  `b.injected`：整包摘掉那一支⇒顶层 FAIL **1 枚＝pristine 基线同一枚**（`TestTicket252R1AlignmentAddsNoAuthorization`，票 298 `AC#5` 已承接＝**净新增 0**），且 309 仪器自己照 `--- PASS`（A5）＝**没有**任何断言讲这一肢。世界**造得出**（`ticket90_test.go:366` 就写了 `Assessor: assessor`，尺＝`grep -rn 'Assessor:' internal cmd --include=*.go`＝2 行＝1 注释 `bridge.go:132`＋1 测试；生产写点 **0** ⇒ r2 那句"生产写点 0 枚"实话）＝缺的是**断言**⛔ 输入。
  `b.prov == nil`：整包摘掉⇒顶层 FAIL **17 枚＝基线＋16**（A6，正是 r2 自报"16 枚红"那枚数，我复现），FAIL 名册里**没有一枚**讲这一肢（`grep -cE '^--- FAIL' | grep -icE 'prov|assessorfor|taint|309'`＝**0**），309 仪器照 `--- PASS`＝噪声红。⇒ `⛔ 具名该肢的仪器` 两肢**同色**（r2 顶回派单 §5 的那处更正＝**成立**，我这一发独立量到）。
  台件怎么造（⛔ 语言层手段，全用现成夹具）：照 `assessorfor_taint_instrument_309_test.go` 的 ② 腿改两枚前提——(a) `New(Options{Registry, Paths, Gate, Logf})`（**不写** `Provenance:`）＋非空 taskID，断"同一句哨兵文本判不出 R4，而挂了 prov 的那枚 bridge 判得出"；(b) `New(Options{… Assessor: risk.NewRiskAssessor().WithCanonicalizer(paths).WithSensitiveClassifier(cls)})`（无探测器）＋非空 taskID，断同一件事。零新导出 API、零新字段。

**③ roots 覆盖（现跑答复，件 `s13-roots-outside.txt`／`s14-outside-roots.txt`／`s15-frontend-root.txt`）**
`cmd` **在**射程（钉的 roots 第 1 根，`cmd=40` 枚非测试源真解析）；`tools` **在**（B3 实测红）；`scripts` **在**（B4 实测红）；`internal` **在**且**不再整包跳**（B2/E2 红）。
`_test.go`＝按**路径后缀**滤（`strings.HasSuffix(path, "_test.go")`），⛔ 按行内容滤——这正好是注释里警告的那枚简写假数的反面，两处口径我各量到（工作树简写＝2；导出树 B13 那一发简写＝**4**，AST 钉仍＝1＝词面尺把注释里的字面量数成**假枚**、解析器尺⛔ 数）。
注释行＝AST 不认注释，B13/E4 实测绿（一枚只把 `agent.ToolRequest{…}` 写在注释里的文件⛔ 进分子）＝**⛔ 假枚**。
★新发现的**根外**：`frontend/embed.go` 是 `go list ./...`（**35 枚包**）里唯一落在那四根之外的模块包（`github.com/CarlosShao/wisp/frontend`，1 枚非测试 `.go`，今日 `ToolRequest` 命中 **0**）⇒ 见 `Q-309-v2-4`。`.scratch/**` 下 615 枚非测试 `.go`（其中 314 枚 tracked）全在**点目录**里＝Go 工具链自己看不见 ⇒ ⛔ 该进 roots，⛔ 算洞。
**这算不算 `AC#1` 的欠？＝⛔ 算。** `AC#1` 问的是"三枚前提装上仪器"，roots 是 `AC#2` 那一格的载荷物；`Q-309-r2-scripts-roots-admitted-by-prefix-only` 里"scripts 到底进没进射程"这一问由 B4 直接答：**进了**。它另外那半（"前缀放行"是现读句⛔ 是保证）我另跑一把：nested `go.mod` 三枚 module 前缀逐字含 `github.com/CarlosShao/wisp/`（`tools/d22scan`／`tools/mockllm`／`scripts/spike`），活例＝`tools/signmodels/main.go:34` ⇒ r2 的判据今天**如实**，但那句"⛔ 是保证"的自陈⛔ 变。

**④ `bridge.go` 那三句注释，逐枚自己现跑（件 `r4-grep-rulers.txt`／`s11-universal-claims.txt`／`rigA-*.log`）**
1. `:982` "Who rings if this limb's shape changes is Test309AssessorFor…"＋"that test is the only instrument, and this comment states the debt, not a guarantee"＝**实话**：该用例在场（`go test -run` 点名真跑）、四形各红（A1–A4）、同段自陈"⛔ guarantee"。⚠ 一句要收窄：`⛔` 这一肢⛔ 是唯一一枚——真正"有人不传 id 时"响的是**别家**测（D 组实测：把调用点 `b.assessorFor(req.TaskID)` 改成恒空 id ⇒ 整包净新增 **7** 枚红，含 `TestFSReadTaintFeedsR4`、`TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2` 与三枚 Pointer185；改成 `req.CorrelationID` ⇒ 净新增 **3** 枚＝`TestToolCallRowsAreComplete`／`TestFSReadTaintFeedsR4`／`Test283IdentityChainThroughTheRealLoop`；**两发里 309 那枚仪器都照 `--- PASS`**，件 `rigD-callsite.log`）⇒ "谁响"这句对**肢的形状**为真、对**调用点**为⛔（`Q-309-v2-5`）。
2. `Every other gate still runs` **已不在盘上**（`grep -n 'Every other gate\|still runs'`＝0 命中）；换成的那句 `:949-952`（"the other rules are STRUCTURALLY still reached on that path - and nothing in ticket 309 asserts it"）＝**实话且带自限**：`internal/risk/rules_gateway.go` 里 `func rule*` **4** 枚，读 `ctx.taint` 的只有 `ruleTaint`（`:101/:102/:105`），退路 `b.assess` 在 `bridge.go:211-213` 确带 canonicalizer＋sensitive classifier ⇒ "少报只少 R4 那一支"结构成立；那句"nothing in ticket 309 asserts it"⛔ 是全称声称，正合编排者 §4.3 要的射程句。
3. `:878-890` 那把 `grep -rn 'ToolRequest{' … yields the single point`＝**逐字现跑对得上**（71／69／2／1，简写形＝2 也如注释所说）；注释里"measured at HEAD 6f389389 … whose go files are identical to the dispatch anchor 202ed4a5"＝**真**（`git diff --name-only 202ed4a5 6f389389 -- '*.go'`＝**0** 枚）。⚠ 同段下面 `:959-963` 那句"every/wherever"**⛔ 实话**（B14/E3 实测绿）——这就是 `AC#3` 唯一那枚条件。
4. 全称词清点（我穷举过这一族写法，⛔ 只引一枚）：`bridge.go` 里命中 `every\|wherever\|only\|no production caller` 的四行＝`:951`（自限，真）／`:957`（带尺＋自限，真）／`:959-963`（**⛔ 真**）／`:985`（"the only instrument"，射程＝这一肢，真）。

## 欠账（具名，⛔ 我给任何格盖章翻勾）

- `Q-309-v2-1-assembly-limbs-stay-comment-only`＝`b.injected`／`b.prov == nil` 两肢⛔ 具名仪器（编排者 §4.4 裁⛔ 造第二枚，注释那半句我实测如实：A5 净新增 0／A6 净新增 16 且⛔ 一枚讲这一肢）。归装配那一族；台件形状见上面 ②(iii)。
- `Q-309-v2-2-nail-roots-and-floor-unpinned`＝roots 列表与分母 floor 都⛔ 被任何东西钉着：**B8** 把 roots 收成 `cmd internal` 两根、同时 `tools/` 里真有一枚第二构造点⇒**绿**（净少扫 24 枚源）；**B7** 把 `examined < 200` 改成 `< 1`⇒**绿**。修法便宜：把"每 root 分母下限"或总分母（现量 293）当第二枚钉，或把 roots 名册本身钉住。
- `Q-309-v2-3-blind-shapes-are-production-reachable`＝**B14**（`new(agent.ToolRequest)` 那枚真第二构造点令钉绿）＋**E3**（`[]agent.ToolRequest{{TaskID: ""}}` 内层省类型，钉绿）＝把 `r2` 具名⛔ 实测的那两条漏形里的两条补上实测读数（另：`Q-309-r2-nail-blind-shapes-unrigged` 剩三条仍**具名⛔ 实测**＝函数体内别名／改名 import 与 dot-import／定义型＋转换）。⇒ 直接判死 `bridge.go:962` 那句 "wherever"（`AC#3` 条件）。
- `Q-309-v2-4-frontend-is-a-fifth-module-root`＝`frontend/embed.go` 在 `go list ./...` 的 35 枚包里、在四根之外；今日 0 命中。裁：要不要进 roots（进＝词面尺分母也要重述；⛔ 进＝"四根全覆盖"那句注释要收窄成"四根"）。
- `Q-309-v2-5-who-rings-at-the-callsite-not-in-309`＝调用点解耦由**票 124/175/183/185** 那族测的牙接住（D1 净＋7／D2 净＋3），309 自家仪器两发全绿 ⇒ 注释那句"Who rings…"该配一句" rings，但⛔ 由本票那枚仪器"。
- `Q-309-v2-6-legacy-go-snapshots-under-probes`（⛔ 本票三格账，只登记）＝`.scratch/**` 下 tracked `.go` **314** 枚（含 `probes/309/a1/logs/bridge.go`、`probes/309/a1/logs/loop.go` 等），与"⛔ 在 `probes/**` 下建 `.go`"那条规矩相反；对我这四把门禁零影响（d22scan clean／pathlen GREEN 实测）。

## 我这发⛔ 做了／⛔ 复跑（分层）

- `Q-309-v2-A-scope-not-covered`＝⛔ 跑整仓测试名册（只 `./internal/tools/`）、⛔ 跑格式名册与 eol 三把（`AC#4` 已由编排者翻勾，不重判）、⛔ 翻框／⛔ 动 `Status:`／⛔ 改名 `-done`／⛔ push（`A813` 一次未用）。
- `Q-309-v2-B-a6-cause-read-not-per-test`＝A6 那 17 枚红的**成因**我只逐字读了 2 枚（`TestToolConcurrencyCeilingIsFour`、`TestPerToolTimeoutHonoredViaContext`＝"L2 审批通道尚未接入（票 21）"文案不符那一族），其余按名册＋"⛔ 一枚名带 prov/assessorFor/taint/309"这把尺推＝"全是管道噪声"这一句是**读码推**，⛔ 16/16 逐枚量过。
- `Q-309-v2-C-r359-numbers-i-did-not-repeat`＝R359-1（module 前缀）／R359-2（`^package agent$`＝34 枚全在 `internal/agent/`）／R359-4（`Assessor:` 2 行）／R359-6（词面尺两支）／R359-7（分母四根）我各现跑过；R359-3／R359-5（`filterdemo`、`bridge.go` 逐枚对账 7→10、2→1）**⛔ 重跑**，只引用。
- 真机量过 vs 读码推：上面每一行标了 provenance；"缺世界"那三肢的**存在性**是读 `helpers_test.go`/`ticket90_test.go` 推的，**颜色**（净新增红数）是 `go test` 量的。

## 撤销口令（逐枚，编排者复跑对不上即撤）

「309-v2 撤 `AC#1`」／「309-v2 撤 `AC#2`」／「309-v2 撤 `AC#3`」；改判形 ⓐ 仍用既有口令「309 改形 ⓐ」，届时本表与 `309-v1` 表整批作废。
