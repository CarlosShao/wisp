# 165-v1 — 票 165（`plan.present`·甲路）非实现者终裁判语

- 裁决腿＝`165-v1`（只读·非实现者）。实现者＝`165-r1`（`internal/tools/plan_present.go`＋`plan_present_165_test.go`＝`4a49480da`）
  与 `165-r2`（死腿零提交，遗产由编排者代提＝`17ca21ee`＋补交证据件 `988ce413`，另 `a754a390` 是 r1 的票尾 bullet）。
- 起手 `git rev-parse HEAD`＝`71002d001af60feaf2bbd9b2a07863a6cb738a23`（`dev`），时刻 `2026-10-16 10:17:11 +0800`；交件时刻见文末。
- 票框尺（起手现跑）：未勾 **5** 枚／已勾 **0** 枚。⛔ 本件翻任何一枚框——五枚框的翻法归编排者。
- 台面＝三棵 `git archive` 仓外导出树（`t0`=起树锚 `b74b87d0`、`t1`=HEAD、`t2`=HEAD 留突变态），⛔ `-overlay`、⛔ 共享工作树、⛔ 仓内 worktree；
  逐把尺的原文、哈希与红句在 `.scratch/wisp/probes/165/v1/00-readings.md` 与 `01-mutations.md`（本件只写判语与凭据指路）。
- 读数纪律：下面每一枚"几枚"都带单位（枚路径／枚顶层用例／枚子测试／枚调用点／枚文件）；`RUN=`（`^=== RUN`）与 `PASS_TOP=`（`^--- PASS` 顶层）**分列**。

## 结论（一行）

**五格判语全为「成立」**；附**四处交付缺陷**（全是描述面与尺的射程，⛔ 任何一格因此不成立）与**六项本腿不裁、归编排者或人拍板**的清单。
建议编排者按 §1–§5 翻框，并把 §6 那四处按"只追加不抹平"的形落进台账。

## 1. AC#1 未修码读数 —— **成立（票前提坐实，且是跑出来的、⛔ 推的）**

- 判据原文问的是"今天只读地出一份方案能不能做到"。答：**做不到**。
- 凭据（四条，全我现跑）：
  1. 代码名册尺：`t0` 树里 `plan.present`／`BuiltinPlanEntries`／`planPresent` 三个词形在 `internal/`＋`cmd/` 的 `.go` 里＝**0 命中**；同一把尺在 `t1`＝**4 枚文件**命中（尺有牙）。
  2. 契约面尺：`plan.present` 那一行**只在冻结文本里**——`docs/PLAN.md:2567`＋`docs/specs/SPEC-07-tools-and-plugins.md:73` 两处现读都在（来源笔＝`824bb9c`，2026-09-27），码侧零枚 ⇒ 「契约已立、码里没有」这一形坐实。
  3. 未修码那发（跑）：`t0`＋我写的同形状尺（只读三个既有 builder＋registry 面）⇒ rc=0、`--- PASS`，逐字 `未修码名册 10 枚、零枚 plan/propose 形、registry.Lookup(plan.present)=false`。
     ⇒ 我这一发比 r1 交付的那枚（只数名册）**多一层 registry 面**；r1 那枚尺的 10 枚名册与我逐枚同名（`fs.read…task.spawn`）。
  4. 生产可达那发（跑）：把 r2 交付的调用点尺原样盖进 `t0`、靶向 `-run` ⇒ **rc=1**，阳性对照那行仍绿（先例 subagent 块在 `t0` 已接），
     红句逐字 `接线 RED：tools.BuiltinPlanEntries 在 internal/ 与 cmd/ 的产码里调用点 0 枚，要 1 枚（cmd/wisp/run.go 的装配那一处）；命中清单＝[]（0 枚＝写了没接线那一族，plan.present 在生产装配里根本够不到）`。
- 恒真三形：未修码上响不响＝**响**（第 4 发 rc=1）；反形＝第 3 发加一味即变 11 枚／1 枚 plan 形（`t1`，rc=0 且用例要求它响）；一把尺两种形状都放行否＝**不放行**（静态 grep 与跑测两条独立路径同读「零枚」）。
- 本格边界（具名）：这五条只裁"票前提"，⛔ 裁"接线之后有人真会调它"——后者见 §6 第 4 项。

## 2. AC#2 复用形状 —— **成立**（判据＝"没新增状态"；⛔ 本格读成"功能做全了"）

- 凭据（我现跑）：
  - 面尺：四枚 165 提交的全部 **26 枚路径**里，`internal/risk/**`＝**0 枚**、`internal/statemachine/**`＝**0 枚**、`internal/agent/approval/**`＝**0 枚**、`internal/panel/**`＝**0 枚**。
  - 通道尺（两把同读）：`Gate` 的回调＝**2 枚**（`internal/tools/gate.go:110`，`PendingWindow`＋`PendingApproval`）；文本尺与交付件的反射尺（`NumMethod()!=2 ⇒ t.Fatalf`）同读数。`GrantSource`＝**1 枚**方法（`Covering`，`internal/tools/grant.go:39`）。
  - 新形状尺：`internal/tools/plan_present.go` 里 `Answer[A-Z]|Mode[A-Z]|State[A-Z]`＝**0 命中**；`planPresent` 的方法名册＝**4 枚**（`Name`/`Description`/`Parameters`/`Execute`＝`Tool` 接口那一面，零枚额外出口）。
  - D43 引用点（我现读，⛔ 抄 r1 的号）：`#17`=`docs/PLAN.md:3077`、`#23`=`:3083`、`#24`=`:3084`（含超时 300s）、`#25`=`:3085`；
    码侧＝`internal/statemachine/table.go` 的 `D43: 17`（`:118`、`:123`）／`23`（`:158`、`:163`）／`24`（`:168`、`:172`）／`25`（`:176`、`:180`）＝**4 个号、8 条转移条目**，
    且这 8 条全是 165 之前就存在的条目（面尺已证 165 在其上零字节）⇒ r1 那句"复用 #17/#23/#24/#25"**名实相符**。
  - 反形（我自己造的世界，`t2`）：给 `Gate` 加第三条回调 `PendingPlan` ⇒ **rc=1、唯一红**＝`Test165PlanPresentAsksThroughTheExistingL2Card`，
    逐字 `Gate 有 3 个方法：本票复用现成的 PendingWindow/PendingApproval 两个回调，多出来的第三个就是一张新审批通道`；另四枚仍 PASS。⇒ 这枚尺有牙。
- 恒真三形：未修码（`t0`）上这格无从谈（工具不存在，红位由 §1 第 4 发接住）；反形响＝上面那发；一把尺两种形状都放行否＝"两枚回调"这形与"新增 State 常量"这形由**两枚不同的尺**分别盯（反射数枚＋grep 常量名册），两形都放行不了。
- 具名"缺世界"（⛔ 判失败、⛔ 由我裁）：`AwaitingApproval` 遇系统挂起那一格行为**至今未定义**（`docs/PLAN.md:1530` 的 DEFERRED 条目、`AGENTS.md §2` 的具名停手项，归 S7）——
  夹具里没造过这个世界，本票也没有任何一枚尺该为它负责。碰到即停手上报，归编排者拍板。
- 读法提醒（票面自己写的、本腿再钉一次）：`docs/PLAN.md:2567` 那一行自带一句"本票因此不减少确认次数"⇒ 「一批只问一次」**不是**本票的 AC，⛔ 在 AC#2/AC#3 成立之外再欠这一格。

## 3. AC#3 批准只批这一批 —— **成立**（两发我自己重放、复现「唯一红」，另补两形都有牙）

- 判据①（批完方案后来一发清单外的写必须响）：
  - 基线：`t1` 靶向 `Test165PlanApprovalRingsAgainForAnUnlistedWrite` PASS（rc=0 五枚全绿）。
  - 我重放的 M3（＝r1 的 M3，`internal/tools/bridge.go:451` 的 `PendingWindow` 调用换成强制 `AnswerAllow`）⇒ **rc=1、唯一红**，
    逐字 `方案批过之后，清单外的写没有响（L1 窗口被问 0 次，要 1 次）——批准变成了对后续写的放行`；还原后 `bridge.go`＝`a0c0736c639276b0835be7c0f85df9dff517eef73039ee882842759d63accfdd`（与 r1 基线逐位相等、与 HEAD blob MATCH）。
  - 装配那一半：我重放的 MA（删 `cmd/wisp/run.go:822-827` 那 6 行注册循环）⇒ **rc=1、RUN=3、PASS_TOP=0、FAIL_TOP=3**，三枚红句逐字（见读数件 01 §3），
    另两枚各红 11 秒＝真走到装配那一跳才红；还原 MATCH 后靶向重跑 rc=0。
- 判据②（摘掉「批过的范围」这一味，是否有一发越界从此打不红）：
  - 我重放的 M2（删 `planPresentText` 里那句 `planPresentScopeClause(len(steps))`）⇒ **rc=1、唯一红**＝`Test165ScopeClauseNamesTheBatch`（两行 Errorf 同一枚顶层用例）⇒ r1 的「唯一红」复现成立，
    **不存在**"摘掉之后越界打不红"那一形。
  - 我另造的第二形 M2b（边界句留在文本里、但数写死成 `1` 而批的是 3 步）⇒ **仍唯一红同一枚用例**（逐字带 `批准范围只限上面这 1 步` 那句被拍到）⇒ 这枚尺钉的是"句子的数跟批走"，⛔ 只钉"句子在不在"。
- 我补造的第三、四形（**r1／r2 都没造过这个世界**，本腿补上）：
  - M9（tools 级：清单外那一发之前往 `plan165Ledger` 种一行 `fs.write` 授权）⇒ rc=1，两条子判据同红：
    `方案批过之后，清单外的写没有响（L1 窗口被问 0 次，要 1 次）` ＋ `批准一份清单之后，会话授权里多出了 [fs.write]：批准只批那一发，不许铸出授权行`。
  - M8（装配级：用真 `memory.Store.InsertGrant` 往库里种一行覆盖清单外写）⇒ rc=1，红在授权行计数那一枚断言上，读数逐字
    `plan_assembled_165r2_test.go:267: AC#3① 装配读数：方案卡=plan.present（yes）→ 清单外卡=fs.write（no，IsError=true）→ 授权行 1 枚 → 目标文件未落盘`。
    ⚠ 形状差具名：这一发里**卡照样升起**（清单外那发在这套夹具里走 L2 支，`Covering` 只短路 L1 支）⇒ 装配用例对这个坏世界的凭据是**授权行计数**那一枚、⛔ 是"卡没升"那一枚。
  - 阳性对照（同树新写的一枚小尺）：同一会话身份 `ListGrantsBySession` 先读 **0 行** → 真 API 种 **1 行** → 再读 **1 行**（rc=0）
    ⇒ 交付件里那句"授权行 0 枚"**不是死尺**。
- 恒真三形：未修码上响不响＝§1 第 4 发已红；反形响＝M2／M2b／M3／M9／M8／MA 六发各红在其该红的格上；一把尺两种形状都放行否＝范围条款那枚尺**两种坏法都不放行**，
  授权那一形在 tools 级与装配级分别放行不了（但**只有装配级真落盘**，tools 级靠假账本，见上两形的区别）。

## 4. AC#4 面板不许代答 —— **成立**（按"能力类判据必问生产调用者"裁）

- 词面尺（我现读，两形都查）：
  - `type Decl struct` 字段名册＝**8 枚**（`Capabilities`/`Needs`/`Declared`/`PathParams`/`Facts`/`Timeout`/`Resident`/`Provider`），**零枚**含 `allow`/`approve`/`grant`；
    `PlanPresentDecl()` 函数体里这三个词形＝**0 命中**；`plan_present.go` **剥掉行注释后**全文对 `allow|approve|grant`＝**0 命中**（剩下的只在注释散文里）。
  - 参数面＝顶层 `properties` **2 枚键**（`summary`/`steps`），`planPresentArgs` 字段＝**2 枚**，且 schema 顶层写死 `"additionalProperties":false`。
- 生产调用者尺（⛔ 拿 `go build` 的 rc=0 当"接上了"）：
  - `git grep -n "BuiltinPlanEntries|plan.present" -- cmd/wisp internal/panel`（剥 `_test.go`）＝**1 枚调用点**（`cmd/wisp/run.go:822`，且 AST 尺现读它 `range=true register=true`＝与先例 subagent 块同形）；
  - `internal/panel/**` 对 `plan.present`／`BuiltinPlanEntries`＝**0 命中**（分母＝`git ls-files 'internal/panel/*.go'` 的 **41 枚跟踪文件**，含测试件）；
  - 面板侧那四枚既有边界钉我在 `t1` 自己复跑（靶向 `-run`，⛔ 整包）：**rc=0、RUN=6 枚（4 顶层＋2 子测试）、PASS_TOP=4、FAIL_TOP=0**
    ＝`TestAnsweredPanelRoutesCarryNoApprovalDecision`／`TestRealGuardRefusesEveryAssemblableApprovalRouteName`／`TestNoInboundEnvelopeCanBindAnApprovalVerdict`／`TestGrantWireShapesAreRefusedAtTheDoor`；
    本票在 `internal/panel/**` 与 `internal/agent/approval/**` 上零字节（§2 面尺）⇒ "发起请求的出口"加了、"给出允许的出口"一枚没加。
- 反形三发（有牙＋一枚射程边界）：
  - M6b（顶层加 `"approved":{"type":"boolean"}` 参数）⇒ 唯一红 `参数 "approved" 让模型自己说「这批已获准」：本票只许加发起请求的出口`（`:309`）。
  - M7（给 `Decl` 加 `AutoApprove bool`）⇒ 唯一红 `Decl 多了字段 AutoApprove：声明面能装下「已批准」就等于加了一个给出允许的出口`（`:292`），两发各跑两遍同句。
  - **M6（rc=0＝放行）具名为尺的射程边界**：把同一个键塞在 `steps` 属性**内部**时那枚尺读不到（它只扫顶层 `properties` 的键）。
    本腿判这一形**仍不构成出口**：`planPresentArgs` 没有字段接它、顶层 `additionalProperties:false` 拒未声明的顶层键；
    ⇒ 记为"尺的射程比声称的窄一格"，⛔ 记成漏洞，也⛔ 让 AC#4 因为这一格不成立。若要补，形＝把那枚尺扫到嵌套层（属产码改动，归实现者波次，⛔ 本腿）。
- 恒真三形：未修码上响不响＝`t0` 里调用点 0 枚已红（§1）；反形响＝M6b／M7；一把尺两种形状都放行否＝顶层形不放行、嵌套形放行（上面已具名＋两味兜住）。

## 5. AC#5 契约轴零字节 —— **成立**（越界尺命中 0 枚，带分母）

- 尺＝对**两腿全部提交**（`4a49480da`＋`a754a390`＋`17ca21ee`＋`988ce413`）逐枚 `git show --name-only` 取并集，再对禁区名册作差。
- 分母＝**26 枚路径**（9＋1＋15＋1，逐枚名册见读数件 00 §4）；禁区名册＝票面那一串＋`go.mod`/`go.sum`/`.gitignore`/`.github/`/`*ci.yml`/三枚冻结件/`tools/d22scan/allowlist.txt`/`internal/perm/ticket90_persist_test.go`。
- 结果＝**命中 0 枚**（`grep -c` 回 0 且退出码 1、stderr 空）。正控＝同一把尺喂 `824bb9c`（D34 冻结笔）的 2 枚路径 ⇒ **命中 2 枚**（尺有牙）。
- 三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）在 26 枚名册里＝**零枚出现**。
- 归属尺（⛔ 把别腿的路径算成 165 的账）：`b74b87d0..HEAD` 共 **93 枚提交**，其中动过禁区面的＝**8 枚**，逐枚＝`7c80ac032`／`f82b909e9`／`9b820f8e8`／`c7752e536`／`73eda0294`／`f3d447346`／`272405f53`（320-r2）／`4f2967e90`——**零枚署名 165-r1／165-r2**。
  另具名一笔早于本波的契约笔 `824bb9c`（动 `docs/PLAN.md`＋`docs/specs/SPEC-07`），那是 D34 那一行的来源、⛔ 记进两腿的账。
- 一处方法差具名：`165-r2` 的证据件 §7 用**当时工作树**的 `git status --porcelain` 当 AC#5 的尺（分母＝盘面），⛔ 等价于本票要求的"对全部提交逐枚 pathspec 作差"⇒ 那一格由本腿补跑（上面的 26/0）。
- 恒真三形：未修码上响不响＝本尺与"修没修"无关（它量字节归属），反形＝正控 2 枚命中；一把尺两种形状都放行否＝状态树尺会放行"已提交但盘面干净"那一形，正是 r2 §7 的形状，已由越界尺取代。

## 6. 交付缺陷（四处，全在描述面／尺射程；⛔ 影响任何一格成立）＋ 本腿不裁清单

1. `internal/tools/plan_present_165_test.go:14-15` 的文件头注释写「两发变异（**-overlay 注入**、树零改）」——与同腿 `06-readings-and-face-check.md` 那句「仓外备份＋sha256 还原证」自相矛盾，
   且 `-overlay` 冒充物理突变是本编队明令禁止的形状；真量是**三发**、非 overlay（我在同一份字节上重放 M2/M3/M1 全部复现）。⇒ 建议编排者按"只追加"更正该注释的账（⛔ 本腿改产码、⛔ amend）。
2. r1 票尾 bullet（`a754a390` 那笔对应的 bullet）写「M1 Declared L2→L1 ⇒ **AC#2 那条红**（附证）」，真量＝**两枚顶层红**（AC#2 用例＋AC#3① 用例；r1 自己的 `05-m1-level-downgrade.txt` 也记两枚）⇒ 少报一枚，属自陈与盘上读数不符。
3. `165-r2` 读数件**缺 10-* 编号**（编排者 A998 点名要本腿答"09 与 11 之间那格判据有没有落过盘"）：答＝**判据落过盘、⛔ 缺的是编号那一枚**。
   那一格的内容（AC#5 面尺＋ban #8 分母）落在件 `08-ac2-ac5-face-roster.txt` 与件 `01b`/`11`（568＝568、`cmd/` 136→138）里，证据件 §7 引用的正是这三枚；本腿⛔ 追认"10 号件曾存在"。
4. 「有没有人真的按门铃」＝**零读数**：全仓跟踪件对 `plan.present` 命中 **56 枚文件**，同一把尺喂夹具面（`testdata/**`、`**/*.golden`、`**/*.sse`、`internal/llm/**`）＝**0 枚** ⇒ 模型侧真调用今天没有任何凭据（`165-r2` 证据件 §10 第 1 条自己也具名了）。
   这**不是**票 165 五格里任何一格的判据（那一格归票 161），但它是"这一档刹车到底有没有人用"的欠账，建议编排者按 `Q##` 或票 161 的账挂住。

本腿不裁、往上交的六项：① 五枚 AC 框的**翻法**（归编排者）；② `AwaitingApproval` 遇系统挂起（DEFERRED，`docs/PLAN.md:1530`，S7）；③「方案条数上限」阈值（契约未说，两腿都没发明数字）；
④ 面板那一屏（`frontend/**`／`design/**` 归机主外部 agent）；⑤ 整包 `go test ./cmd/wisp` 四数（票 314 死亡区＋本波那一档归 `220-r2`，⛔ 本腿跑）；⑥ `-tags winlive` 那族（台面有别腿、归编排者车道）。
另具名一处**本波别腿的账**（⛔ 算 165）：共享工作树里 `go vet ./cmd/wisp` rc=1，红句 `cmd\wisp\panel_l1_blocked_220_test.go:236:84: (*testing.common).Errorf call has arguments but no formatting directives`，
该件 `git status` 读作 `??`（未跟踪、`220-r2` 在飞）；同一把尺在 HEAD 导出树里 rc=0。

## 7. 自证（三把尺）与时刻

- 交件时刻：`date '+%F %T %z'`＝见本腿最后一枚提交的时间戳（起手为 `2026-10-16 10:17:11 +0800`）。
- 写面（只三处，逐字照派单）：`.scratch/wisp/probes/165/v1/`（本件之外的 2 枚读数件）、`docs/evidence/s1/165-v1-acceptance.md`（本件）、票 165 末尾**一条** bullet。
- ⛔ 产码、⛔ 翻 AC 框、⛔ 动 `docs/PLAN.md:2567` 那行冻结文字（本件只**引用**它的行号与内容）；只 commit、⛔ push；两处 `git add`／`git commit` 逐枚显式 pathspec。
- 突变台件与原始日志留存（只建不删）：`D:\tmp\165v1\{t0,t1,t2,kept,logs}`；六枚被改文件全部对 HEAD blob 逐枚 MATCH（读数件 01 §0/§6）。
