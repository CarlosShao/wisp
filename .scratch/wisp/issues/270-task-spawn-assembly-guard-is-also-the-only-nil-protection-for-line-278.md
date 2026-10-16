# 票 270 — `task.spawn` 那道装配 guard **一票两形**：它同时是 `:278` 空函数值调用的唯一 nil 防护，拆之前必须先补

- Status: **已立，未派**（10-06 18:3x 编排者立；来源＝非实现者终裁腿 `236-v1` §5 第 2 行／§2 第 6 行，裁决表 `.scratch/wisp/probes/236/v1/verdict.md`；台账 `A643`）
- 关联：票 236 AC#1b（那三格已翻勾，本票是它交回来的**产码残余**，⛔ 不是 236 的欠账）／票 197（子代理三层）／票 164（名册只在本进程内）

## 这是什么（人话）

`task.spawn` 在派生子代理前先问一句"宿主给我装配了吗"，没给就**拒绝**。问题是这句拒绝顺带挡住了一件更糟的事：两行之后代码会**无条件**去调用那份装配，而"没给"的意思是那里**压根没有一个函数**——调用一个不存在的函数不是报错，是**整个测试二进制当场崩**。

所以今天这道 guard 承担了两件事：①给出诚实的拒绝文案；②**不让下一行崩**。任何把它拆成两支语句的"整理"，都会在只留一支的那一形上把②弄丢。

## 现量（⚠ 引用前先重跑；全部编排者 10-06 18:3x 自己跑，行号为真）

1. `internal/tools/subagent_197.go:256-258` ＝ `if t.d.BaseOptions == nil || t.d.ParentTools == nil { return Result{Text: "宿主没有给出派生用的装配…", IsError: true}, nil }` ——**一枚语句管两枚字段**。
2. `internal/tools/subagent_197.go:278` ＝ `base, wired := t.d.BaseOptions()`，**无条件调用**，就在 guard 之后两行。
3. `internal/tools/subagent_197.go:135` ＝ `BaseOptions func() (agent.Options, bool)`（函数值字段，nil 时调用＝panic）。
4. 生产侧唯一产出者＝`cmd/wisp/run.go:789`。
5. ★崩溃读数（`236-v1` 造，⛔ 本编队复述其原始读数，不是推断）：`A-4-without-recover`（guard 永假＋摘掉测试里那层 recover）＝**rc=1／PASS=25／FAIL=2／日志 `^panic` 2 行／`FAIL github.com/CarlosShao/wisp/internal/tools 1.444s` ⇒ 其后约 181 枚用例零读数**；正控 `ctrl-no-recover-only`（**只**换测试文件、产码一字不动）＝**rc=0／PASS=206／FAIL=0** ⇒ 失明只能归因于产码那一次空调用。

## 要建什么（⛔ 形不在票面写死，落地前先普查；三形都要量）

- **甲**：`:278` 改成"取之前先判空"（nil ⇒ 走同一句具名拒绝），guard 保持现状。⇒ 最小、零文案改动，但**拆语句仍要另批**。
- **乙**：把 `:256` 拆成两支独立 guard（各自文案），**同批**给 `:278` 装甲的判空。⇒ 形状最干净，代价＝两枚现有用例的"哪半边红"要重写，且它触碰票 236 AC#1b 刚钉的完整字面量（**必须同批把那些钉一起改**，不许先拆再补）。
- **丙**：不动结构，只在 `:256`／`:278` 各加一条**具名常驻判据**钉住"摘掉 guard ⇒ `:278` 必崩"这一事实（现由测试层 recover 转红，改由产码层拒收）。

判据要覆盖的形状：**拆掉任一半 guard 之后，`internal/tools` 整包不许出现"零读数"**（今天只有 recover 那一层替我们兜着）。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [ ] **AC#0 三形代价表**：甲／乙／丙各自的改动面（文件∶行）、会顶红谁现有的钉（⚠ 至少顶到 `internal/tools/failclosed_236_teeth_test.go` 那枚完整字面量钉与 `236-r2` 的两枚分支钉，须逐名读断言）、以及"拆完之后崩溃面还剩几枚零读数"的实测。⛔ 未出表之前任何一形都不许开工。
- [ ] **AC#1 `:278` 不再无条件调用**：`BaseOptions == nil` 时走到**具名拒绝**，⛔ 不许是 panic、不许是静默退化成"没装配也起一套运行时"。完成判据＝删除 `:256` 整支后的台件读数：`internal/tools` 整包**仍出全量名册**（PASS+FAIL ≈ 206，⛔ 不许出现 25／0 那种失明形）。
- [ ] **AC#2 拆分支要拆得动**（只在选乙时适用）：半边摘 ⇒ **恰 1 枚**具名红，两半各一形；这条直接接 `236-v1` 已给的读法（"各恰 1 枚"＝半边摘，"整支摘 ⇒ 2 枚红"＝正确反应、不是形状破）。
- [ ] **AC#3 文案不许被复用**：装配未接线那句与"池满"那句、"没有父任务"那句**三句互不充当彼此的替身**（`236-v1` 实测到"邻居分支理由是错的"那一形：nil 名册时池满门会说"你排满了"）。
- [ ] **AC#4 门禁四数＋还原自证**：`d22scan` rc=0 且正控先绿；`gofmt -l`／gofumpt `-l` 对本腿文件零输出；`go vet ./internal/tools/` rc=0；终态整包逐名红名册**与起手名册作差＝0**（⚠ 规矩＝整包红＋单包安静 `-count=3` 两发一起，见票 236 AC#5 尚未闭合那格）；`git status --porcelain -- internal cmd` 终态＝起手。
- [ ] **AC#5 后程硬约束落进注释**：`:256` 处写清"这支同时是 `:278` 的 nil 防护"，⛔ 只写现状，不许在注释里许诺未实现的能力。

## 禁区

⛔ 不许为了变绿放宽票 236 那批字面量钉（它们是刚翻勾的凭据）；⛔ 不许新增导出名／新造工具面（D34 无对应行）；⛔ 不许把 `:278` 的修法与本票 AC#2 的拆语句**分批**落地（先拆后补＝中间态可崩，那是本票立案的理由本身）；⛔ 不许顺手把 `ErrorClass` 那一维一起改了（那属〔契约邻接＝D37 错误分类〕，见票 236 交回来的第三件，需人工批准）。

## Progress log

- 10-06 18:3x 编排者立票。产码未动一字（`internal/tools/subagent_197.go` md5＝HEAD 同串）。
- 10-06 19:1x 腿 `270-a1`（只读代价普查）交件＝`.scratch/wisp/probes/270/a1/census.md`（五节全实，⛔ 不选形、不翻框、零 `go` 命令）。读数摘要：①派单五枚锚点（`:135`／`:256`／`:257`／`:278`／`:279`）现量**逐字对上，一枚未漂**；②顶红枚数＝甲 **0** 枚具名红／乙 **3** 枚具名红（`failclosed_236_teeth_test.go:95-96` 共用常量＋`:195`＋`:218` 两枚等值钉）另 **3** 枚次序敏感钉（`:176`／`:257`／`:269`）／丙 **0** 枚（越界加类型或字段则＋2＝`subagent_197_test.go:879`、`cmd/wisp/subagent_selfapproval_197_test.go:563`）；③★新发现＝`:293 opt.Tools = newSubagentToolProvider(t.d.ParentTools)` 同样只靠 `:256` 拦着 ⇒ `:256` 今天是**一票三形**，且 `internal/tools/task_cancel_221_test.go:51` 是派单未列的第 4 处 nil-`BaseOptions` 形状；④残余零读数面（静态）＝`internal/tools` 里 **25 处＋cmd 1 处无 recover 的 spawn 派发点**，三形都不动；⑤`d22scan` 侧唯一真会拦三形的是 ban 8（字符串不豁免）与 ban 9（丙那条注释若引路径必须盘上存在——本腿实测 `docs/DECISIONS.md`／`docs/DEFERRED.md`／`docs/TOOLS.md`／`docs/STATE_MACHINE.md`／`docs/contracts/` **全 MISSING**，`SPEC-05` 真名＝`docs/specs/SPEC-05-agent-core.md`）。⚠ **AC#0 里"拆完之后剩几枚零读数"那一格本腿未交实测**（禁跑 `go`），已作为 §3 第 3 条摆给编排者；本腿未动产码一字、未动任何 `- [ ]`（`git status --porcelain -- cmd internal` 起手＝终态＝0 行）。

---

### [10-07 09:3x] 编排者追加更正（来源＝只读腿 `270-a2` 交件 `aa6e8870`，其件 `.scratch/wisp/probes/270/a2/pairing.md` 103 行／23,759 字节／占位 0；下列三处**我逐条自己复量过**，原话不抹）

**更正一：本票 AC#1 那条完成判据今天恒真 ⇒ 不许当凭据（框未翻，仍待重写）。**
腿给的尺：`236-v1/verdict.md:153`（A-4 那发**带 recover＝名册全 ≈204/2**）与 `:164`（摘掉 recover 才 25/2）。⇒ 我原来那句"摘掉 `:256` 整支后仍出全量名册"正是**反形自己要读的那一发**——零改动也绿。按 [[记忆：判据换成反形仍全绿＝它对这事不敏感]]，这条判据对"有没有装配守卫"不敏感。
⇒ 处置：AC#1 框**保持未勾**；下一枚落地腿的判据改成腿建议的那枚定向形状——"同发突变里 `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired`（等值钉在 `failclosed…_test.go:218`）**红→绿**"，前提＝`:278` 复用 `:257` 那句串（该前提由腿给，落地腿必须自己再验一次）。

**更正二：我票面那句"就在 guard 之后两行"已过期。**
现量：`:258`→`:278`＝**20 行**（`:256`→`:278`＝22 行），不是两行。同一句英文仍在 `failclosed…_test.go:201-203`。⛔ 我只登记，不改写原句。

**更正三：★记我——"没有任何读盘／embed 的尺去读 `subagent_197.go`"这句是我凭空加的，不是 `270-a1` 的结论，且今天为假。**
- `270-a1` 的件里**没有这句话**（腿具名报回，我复按其件亦无）。
- 今天盘上有两把读源码形状的尺，且都在本包：`internal/tools/paths_twocontainments_252_r2_test.go:194` `os.ReadDir(".")` → `:203` 排除 `_test.go` → `:206` `parser.ParseFile(fset, name, nil, parser.ParseComments)`（＝逐枚解析本包全部非测试源，`subagent_197.go` 在其内）；`internal/tools/task_cancel_221_legs_test.go:224` `os.ReadFile("task.go")`。两把我都 `sed` 逐行看到原文。
- ⇒ 影响：**丙形（"只登记不装牙"）并非零具成本**；但腿同时给出关键限定——这类**源形状尺对 `-overlay` 突变失明**（`failclosed…_test.go:68-75`），所以它**不能顶替 AC#1**。两句要一起读，⛔ 不许只剩其中一句。

**裁定（本票不拆，拆票归下一轮）：**
腿建议"拆 2 枚＋1 枚待议"，理由是 `:293` 那一半**按静态读不是当场崩形**——我复量对上了：`internal/tools/subagent_197.go` 的 `Tools()`（`:552`）**有 `if p.inner == nil { return nil, nil }` 守卫**，nil `ParentTools` 不会在构造时炸，而是**延迟到 `Execute` 走 `p.inner.Execute(...)` 那一支**。⇒ 我票面"拆任一半不许零读数"只对 `BaseOptions` 半边成立，**对 `ParentTools` 半边说过头了**。
- 本票射程**收窄为 `BaseOptions` 那一半**（三形同落 `subagent_197.go`、互不打红：甲/丙 0 枚红、乙 3 枚红是它自己那批钉 ⇒ 文件层面装得住，这条腿已证）；
- `ParentTools` 那一半"延迟崩"要不要单立一票，**我按下不立**，等本票 AC#1 判据重写并落地之后再裁（⛔ 不在同一枚票上同时改判据与加射程）。

- 10-15 21:0x 落地腿 `270-r1` 交件（形＝甲，`A990` 裁③；写面只两枚：`internal/tools/subagent_197.go` 与新件 `internal/tools/nilbase_270_test.go` 191 行）。产码两处：`:290` 取之前先判 `t.d.BaseOptions == nil` ⇒ 先 `release()` 再走**与 `:257` 逐字节相同**的那句具名拒绝；`:256` 上方落 AC#5 注释（只写现状）。⛔ 未拆 `:256`、⛔ 未碰 `:296`／`ParentTools` 半边／`ErrorClass` 那一维／`failclosed_236_teeth_test.go` 一字。**AC#1 前提本腿自验**（⛔ 继承 `270-a2`）：两枚命中行 python 比对＝`True`，且该串与 `lit236SpawnNoAssembly` 相等已由新尺钉住。**AC#1 定向判据现场取到**（方法学承 `321-r1`：物理改写＋sha256 还原，⛔ `-overlay` 对读源码的尺失明；台件 `.scratch/wisp/probes/270/r1/mutate.sh`＋`guard_strip.py` 按条件文本定位、恰 1 枚命中否则拒绝突变）：同一"摘 `:256` 整支 3 行"的突变下，**旧形**（HEAD 产码）rc=1／PASS=235／FAIL=4／名册合计 **239**＝全量、`Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` 红且 got 以 `panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference` 开头；**甲形** rc=1／PASS=238／FAIL=1／名册合计 **239**、那枚 236 腿**红→绿**、`^panic`＝0 ⇒ ⛔ 25/2 那种失明形没有；唯一残余具名红＝`Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired`（`||` 的另一半，射程外，本腿没为它放宽任何钉）。AC#3 四枚输入四句理由逐枚等值＋两两不等＋反控两枚（nil 名册不许被池满门回答／装配齐全必须不拒）。门禁九把尺全 rc=0：`go build ./...`／`gofumpt -l` 空／`go vet ./internal/tools/`／`d22scan` clean 且射程非空（ban #7 internal/tools/=25、ban #8 internal/=564 含注释与 `_test.go`）／长度门 VERDICT GREEN 60==60／`-count=2` 478-0 与 `-count=1` 239-0 两发取交集／红名册与起手作差＝0（起手 0 枚、终态 0 枚）／`git status --porcelain -- internal cmd` 终态＝起手（0 行）。证据＝`docs/evidence/s1/270-nilbase.md`（七节，含突变三份 sha256 与还原自证）。**没做完／不适用逐枚具名**：AC#2 因选甲不适用（`:256` 未拆）；`ParentTools` 半边零改动（是否单立一票仍候裁）；AC#0/#1/#3/#4/#5 的 `- [ ]` 框一枚未翻、票名与 Status 未动，对抗验收归非实现者；一处弱牙已具名写进证据 §6 第 4 条（正跑树里 `:290` 那句的文案不可观察，其证明只在摘 guard 那一发成立）。⛔ 未 push。
- 10-16 11:1x **非实现者终裁腿 `270-v1` 交回**（角色＝裁决者·只读：⛔ 翻任何 AC 框——票框现量 **6 枚未勾／0 已勾**；⛔ 动产码/测试件一字；判语件 `docs/evidence/s1/270-assembly-guard.md`＝109 行／19,515 字节，读数件 `.scratch/wisp/probes/270/v1/**`（台件 `mutate270v1.py`＋`run-bench.sh`＋`logs/read-*.txt` 十一枚），台面＝`git archive HEAD(232f309e)` 的**仓外导出树**（`find … -name '*.go' | wc -l`＝1098 证非空），十发突变**全打物理字节**、⛔ `-overlay`。五格判语＝**AC#1 成立／AC#2 不适用／AC#3 成立／AC#4 成立／AC#5 成立**（⛔"附条件成立"）。承重读数逐枚我自己现跑：★`A990` 裁③ 那一发＝摘 `:265`〔票面 `:256`〕整支、甲留在 `:297`–`:299` ⇒ `RUN=325／PASS_TOP=241／FAIL_TOP=1／SKIP=0`、`^panic=0`、`Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` **红→绿**、唯一红＝`ParentTools` 半边（射程外）、突变串 `a9fdc98c…` 与落地腿**独立复现同串**；**去牙**＝在同一世界上再摘 `:297`–`:300` 那四行 ⇒ `325／238／4／0` **四枚具名红**，红句逐字 `nilbase_270_test.go:60: 调用点没自己判空（摘掉上游 guard 后的旧形）："panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference"`＋`nilbase_270_test.go:132`＋`failclosed_236_teeth_test.go:196`/`:219` ⇒ 名册仍全量（⛔ 25/2 那种失明形）；**反形三枚全响**（拒绝换成 panic＝4 红／换成"静默当成没孩子"＝4 红／`:299` 借名册门那句＝4 红，got 逐字 `"任务名册未接线（fail-closed：拒绝派生子代理，派生了也没有地方登记它）"`）⇒ ⛔panic、⛔静默、⛔替身三形都不放行；`release()` 真在被摘的那条路上＝**计数差**那一发坐实（只摘 `:298` 一行 ⇒ `nilbase_270_test.go:74: 被拒的这一枚占着池位：in-flight = 1, want 0`）；还原＝九发逐发 `MATCH=yes`（`86de44c4…`＝`git cat-file -p HEAD:internal/tools/subagent_197.go` 同串），还原后台面再跑一发＝`325／242／0／0`、rc=0。**具名不敏感的两形**（甲形固有，⛔ 算落地腿偷工——他在本条上面那句自报过同一处）：只摘甲那一支＝**0 枚红**（`strip-jia`）、摘 `||` 的 BaseOptions 半边＝**0 枚红**（`guard-keeps-parenttools-half`）；而摘 ParentTools 半边＝恰 1 枚红、那枚钉属票 236 ⇒ **AC#2 这一格在甲形下永远无人钉**，要"摘任一半必恰响一枚"只能另批走乙形。能力类必答：票面那把尺 `git grep -n '"task.spawn"' HEAD -- '*.go'`＝**126 命中**，其中 **104 枚落在 `.scratch/wisp/probes/**` 的存档突变副本**（＝票 324 那种同名假命中；⛔ 去掉它就会被读成"126 枚调用点"），真代码 **22** 枚＝产码 **2**（`internal/tools/subagent_197.go:187` `Name()`、`:583` `subagentHiddenTool`）＋测试 **20**，生产链上装配点 **1** 枚＝`cmd/wisp/run.go:787` 注册＋`:796` 那枚 `BaseOptions` 函数字段（⛔ 拿 `go build` 的 rc=0 当"接上了"）。五枚交付缺陷具名、归编排者处置（① `:257`–`:258` 那句 "the **only** site" 与同块 `:261`／`:294`–`:295` 的 "the only thing left" 互相顶＝过期文字；② 票面行号锚全漂：`:256`→**`:265`**、`:278`→**`:301`**、`run.go:789`→**`:796`**，"就在 guard 之后两行"今天实为 **35 行**；③ `nilbase_270_test.go:79`–`:81` 那枚"准入次数＝0"钉在我全部十发里一次都没响过 ⇒"⛔ 退化成另起一套运行时"这一半句只剩文案等值钉；④ `ParentTools` 半边（`:316`）零改动，`guard-keeps-baseoptions-half` 的红句证明它是延迟崩/空目录形，是否单立一票仍候裁；⑤ 本机 `command -v gofumpt`＝NOT_FOUND ⇒ 那一格 `gofmt -l` 只是**代用尺**、⛔ 冒充 gofumpt；另：长度门那一发在导出树里回 rc=2"not inside a git work tree"＝台面故障，已改在仓根复跑 rc=0 并具名）。`AC#0` ⛔ 本腿射程（派单只交 AC#1–AC#5＋能力类）。纪律：只 commit、⛔ push；写面只三处（判语件／`probes/270/v1/**`／本条 bullet）；`git add` 与 `git commit` 两处逐枚显式 pathspec；⛔ 动 `Status:`／票名／`-done`；临时件只建不删。
