# 183 — **票 177 那枚"宿主自己写下的路径精确豁免"在真机 CLI 上没护住宿主自己的指针**：模型照 `task.output` 给的那条路径去 `fs.read`，被自家 R4 判成"包含来自 task.output 的内容"⇒ 升 L2 ⇒ 当场拒 ⇒ `读回 0 字节`——**"超长输出读回来"这件产品承诺到今天仍未达成**

- Status: **ready-for-agent（先只读核，再落地）**。⚠ 本票是**缺陷票**，不是新功能票：判据在票 177 里已经写死并翻过勾（AC#2 反向／AC#3 正向〔成立·带条件〕），**本票要查的是"那枚豁免为什么在真机接缝上没生效"**。
- 来路：`179-r2`（提交 `391822a2`／`4fbda6a5`，表 `docs/evidence/s1/179-cli-reread-r2.md` 20739 字节）第一次把"模型拿宿主指针去 `fs.read`"这一发在真 CLI 上跑通到**门那一跳**，然后被拒。它按票 179 AC#8 的分岔规矩**具名报回、没写成本票结案、一枚框没勾**。台账 `A362`。
- 关联：**票 177**（豁免本体，`Q-61` 甲＝精确豁免已批）· **票 175-r2**（把 `task.output` 送进 C25 名册那一跳）· 票 164（那两枚续读腿＝同一件事的桥级版本）· 票 162（`窗口 0.0s` 即拒那一形）· 票 174 AC#2b（回执里"C26 未接线"那句，**正交、不是本票的阻断者**）

## 现量（锚 `4fbda6a5`，09-28 11:3x 编排者自己从读数文件里逐字取的，尺可复制）

读数文件＝`.scratch/wisp/probes/179/r2/logs/e2e-readings.txt`（28496 字节）。

1. **阻断那一行，逐字**（尺：`grep -an "rules_hit=\[R4\]" .scratch/wisp/probes/179/r2/logs/e2e-readings.txt`）：
   ```
   51:[audit] tools: call kind=refused task=c3d0f0f7-… corr=c3d0f0f7-… tool=fs.read risk=L2 decision=reject
      outcome=error rules_hit=[R4] in_allowlist_scope=true reason="R4: 包含来自 task.output 的内容"
   ```
   ⚠ **`in_allowlist_scope=true` 与 `risk=L2` 同时成立**＝路径在授权根内、却被外来内容规则升了级。⇒ 阻断者**不是**票 174 那条"artifacts 不在授权根"，是 **R4 命中了宿主自己盖的那枚戳里的路径文本**。
2. **豁免没生效的第二条线索**（尺：`grep -an "dropped=" 同一枚文件`）：
   ```
   44: C25 scope closed task=c66c0634-…（run A，产出那一发） was_open=false dropped=0
   55: C25 scope closed task=c3d0f0f7-…（run B，续读那一发） was_open=true  dropped=1
   ```
   ⇒ run B 自己那枚 `task.output` 的 mark **在场**（close 时落了 1 枚），而续读的 `fs.read` 参数**正是被这枚 mark 判中的**。
3. **到手的前半（说明起跑口与盖戳都活着，别误判成"没接线"）**：`task.output` 执行成功、层级 `risk=L0 decision=allow outcome=success`，回执里逐字带着那条产物路径与省略数（`…省略 17200 字符，总长 20000 字节…全文见 C:\Users\swq\…\artifacts\tool-output-agent-task-c66c0634-….txt`），且模型送进 `fs.read` 的 `path` **与名册登记的那条逐字相同**（台件两向比对通过）。
4. **最坏后果的形状（数字）**：产物 **20000** 字节 / 读回 **0** 字节；整份与 `max_bytes=10000` 两发的工具行都没回到上下文。随后 `[确认 L2 fs.read] … 窗口 0.0s` → `[工具 fs.read -> error]` → run B 退出 1、任务 `cancelled`。
5. **⚠ 一枚"尺不够宽"的自污要一起读**（别拿文件级 `grep -c` 当命中数）：`风险未分级`／`查不到这个任务`／`L2 级…审批通道尚未接入` 三枚在**三区**（console／rig-turns／tool-rows）计数都是 **0/0/0**，而"升 L2"真实存在（`level=L2` 2 处、`risk=L2` 1 处、`[确认 L2 fs.read]` 1 处）。⇒ **判"哪一形在响"要分区分字段，不能一把 `grep -c` 糊过去**。

## 为什么值得做（不做会怎样）

这是 owner 从第一天就在意的那句话——**"超长的输出，我能不能让它自己读回来"**。今天的答案仍然是**不能**，而且**不是"没做"，是"做了但没生效"**：票 177 的豁免有八枚常驻判据（`internal/risk/shape_a_exemption_test.go`）、票 175-r2 的桥级判据也全绿，**可它们全都在桥级／包级，真机接缝上第一次跑就红**。⇒ 不做这一枚，前面那两批的"成立"就永远停在〔带条件〕，而条件是**一个今天不成立的事实**。

## AC（每格都要答"这一发在未修码上响不响"；先测→再写→再提交）

- [ ] **AC#1 先答"豁免为什么没护住"——要读数不要推理**。至少排掉这四支候选，每支给一条现量：
  **(a) 跨度定位偏**——`MarkWithHostPath` 那条路径窗口在**归一化之后**的坐标对不对（`normalizeTaint` 会丢空白与零宽字符；产物路径里有 `\`、盘符、`agent-task-<uuid>` 长串）；
  **(b) 载具没送到**——真 CLI 上 `hostPathBox` 在 `mark` 那一刻是不是空的（`task.go:253-255` 的 `box.set` 与桥的 `mark` 之间的**先后**，尤其名册是 `176-r1` 的**回填**路径填的、不是 `task.output` 当场填的）；
  **(c) 命中的不是路径窗口**——R4 匹配到的片段是不是那句**中文说明**（`全文见 `／`省略 17200 字符`）而不是路径本体；
  **(d) 索引外的第二条路**——`provenance.go:662` 那一支（路径参数从不豁免）与本次命中的关系。
  ⚠ 判"是哪一支"必须**造变异**证明（改一支只有那一支红），不许只读码下结论。
- [ ] **AC#2 正向判据（常驻，落 `internal/tools/`，走真桥具）**：模型照宿主指针续读 ⇒ **不命中 R4**、`fs.read` 成功、**逐字节读回**。⚠ 这条判据必须长在**CLI 接缝形状**上（含回填路径产生的那条 `ArtifactPath`），**不能只是桥级复用**——票 164 那两枚桥级腿今天全绿却挡住了这一发，这就是"包级全绿≠端到端通"第 N 次应验。
- [ ] **AC#3 反向判据（与 AC#2 成对，缺一形＝装饰）**：**外来内容里出现同一条路径、模型随后去读它 ⇒ 仍然要命中 R4**（票 177 AC#2 原句）。⚠ 任何修法若让这一发也变安静＝**洗戳**，本票判不通过。
- [ ] **AC#4 第三形：同名不同目录／改一个字符的路径**必须仍命中（票 177 的 W-3 那一族在 CLI 形状上的对应物）。
- [ ] **AC#5 ⚠ 三支毒修法明禁**：① 放宽 R4 或把"路径类参数一律豁免"（那是把 `provenance.go:662` 那支反着改，票 177 §4.2 已裁过＝匹配侧按值放行与 AC#2 结构冲突）；② 设 `Config.PassThroughUnclassifiedRisk`（票 179 AC#5 同一枚禁令）；③ 把 `fs.read` 从 C25 名册里摘出去"让它别盖戳"（那是拆票 175 的承重声明）。⇒ 判"必须动这三枚之一才修得好"**停手上报**。
- [ ] **AC#6 顺带答一句 `窗口 0.0s`**：真 CLI 上那张 L2 卡的窗口为什么是 0.0 秒（是审批通道未接、还是超时常量在此形状下取零）？⚠ **本票不修它**，只具名归口（候选＝票 162 那一族），并把读数留在表里。⚠ 不许顺手改任何审批超时常量（冻结件）。
- [ ] **AC#7 契约轴**：`docs/PLAN.md`（含 `:1531`／`:1532` 两行 DEFERRED）、`docs/specs/**`、`thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节不许动；`provenance.go:468-473` 那句 `Mark` 契约文字与 `taintmatch.go:11-15` 那句"逐 token 追踪已被否决"**保持不动**（`Q-63` 甲的边界，撤销口令「撤 Q-63 甲」）。
- [ ] **AC#8 门禁**：逐包 `go test -count=1 ./internal/risk/`（**单跑**，四包并发会假红＝`A359`）＋`./internal/tools/`＋`./internal/agent/`；CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，出处 `scripts/wisp-cli-tests.sh:20`；不带就是 `0xc0000135`＝票 98 那枚加载期坑，两形都要贴）；`sh scripts/d22scan.sh`；`gate-clauses.sh` **比红腿名册**不比退码（今天唯一在册红腿 `G6neg`＝票 178）；⚠ 最终读数取在最后一枚 commit 之后。

## 本票**不**解决

- 不裁"该不该给续读加一张 L2 卡"（票 177 的丙形，AC#0 至今零读数，另议）。
- 不修票 174 AC#2b（`TaskDeps.Paths` 在 `cmd/wisp/run.go:365` 仍未接——尺：`grep -n "TaskDeps{" cmd/wisp/run.go`；⚠ **号会从 362 漂到 365 是因为 `179-r1` 改了同一枚文件的注释**，引 `run.go` 行号一律要重锚；⇒ 回执带"C26 未接线"那句 fail-closed 文案；**它是文案缺陷、不是阻断者**，AC#1 的 (b) 支若把它牵进来要分开报）。
- 不修 `窗口 0.0s`（AC#6 只归口）。
- 不翻任何别人的勾（票 177 AC#3、票 175 AC#5、票 176 AC#3/4/5 都**维持现状**，等本票落地那一发端到端）。

## Progress log

- 09-28 11:3x 编排者立票：上面五把尺本程现跑（锚 `4fbda6a5`，读数我自己从文件里逐字取）。未派。
- 09-28 12:1x 只读核 `183-a1` 交件（表 `docs/evidence/s1/183-exemption-why-not-live-a1.md`，AC#1 那一格）：四支候选逐支给现量＋变异兑现，判定**根因＝(d)**——豁免按「位置」排除，命中判定用的是去重后的哈希集＋对整段归一化正文的 `strings.Contains`，两处都位置无关 ⇒ 被豁免窗口里那 8-rune 串只要在正文别处出现过一次，豁免当场失效（真机读数 `frag="-output-" hit=true`，出处 `probes/183/a1/logs/mark-hit.txt`）。**(a)(b)(c) 全部不成立**：坐标 `lo=2101` vs 原始字节 `2244` 两者差 143 且用的是前者；载具送到了（`hostPathBytes=155 lo=767 skip=[[767,921]] excluded=162`，⇒ 派单"回填与当场不同一条"这一前提**证伪**）；命中片段是路径本体不是中文说明。决定性变异 M-D1：只把标记正文换成与路径零同款 8-gram 的正文，同一份生产码同一个 CLI 接缝上那一发就 `hit=false` 并「整份读回逐字节相等」。
- 09-28 顺手两格：八枚常驻判据全在 `internal/risk` 包级、且都不喂「路径同款片段在正文重复」这一形（这就是"全绿挡不住"的具体答案，不是"覆盖不足"）；`dropped=1` 裁的是 run B 自己当场那枚 mark（`marks=1`，名册进程级共享／mark per-scope）。推荐落地点＝`provenance.go:541` 把豁免跨度从喂给 `newFragmentIndex` 的文本里物理挖掉（`taintmatch.go` 一字不动、不必动 `provenance.go:468-473`），但 `MarkWithHostPath` 文档块那句语义要变强＝**人工批准面**，已在 §6 具名报回。**AC#2..AC#8 一枚没做、一枚没勾。**⚠ 另报一枚排程坑：起手第一发被在飞的 `179-v1` 脏件 `internal/agent/loop.go` 污染（读出"风险未分级"），已用 `git show HEAD:` 钉进 overlay 重跑——`183-r1` 派之前建议先收掉 179-v1。跟踪文件零改动（`git status --porcelain -- internal/ cmd/` 每次变异后均为空）。
- 09-28 12:34-12:56 写码腿 `183-r1` 交件（表 `docs/evidence/s1/183-pointer-exemption-r1.md`，台件与读数 `.scratch/wisp/probes/183/r1/**`；**AC 框一枚没勾**）。先按派单立裁判再逐族现量：**AC#0 裁判在未修码上红**（`logs/referee-on-HEAD.txt`，`--- FAIL: TestPointer183RefereeTwinFragmentOutsideSpanStaysExempted … got hit … source=fs.read task.output fragment_len=8`），归因控制腿（同一 stub、正文不含同款片段）在未修码上就绿 => 红确实来自"别处同款片段"那一维。逐族：**族 F1（`183-a1` §6 推荐的挖跨度）没翻裁判**（`logs/f1-referee.txt`：裁判与"路径出现两次"两枚仍红、票 177 八枚 ShapeA 全绿）=> 派单 §2 〔编排者预测〕**坐实，不是推翻**；**族 F2 翻了**（六枚全绿＋八枚 ShapeA 一枚没红）。落的是 F2 的**复核侧形态**：`taintmatch.go` 给 `fragmentIndex` 加一个被读的 `declaredPath`、`contains()` 先逐字问"这枚候选窗口是否拼出这条声明串"（不查哈希 => 碰撞不可能多扣），`provenance.go` 在 span 定位成功后把声明串挂到这枚 mark 的索引上；两枚跟踪件**删除列 0**（48/59 行纯插入），文档块 `provenance.go:482-506` 原句一字未动（md5 仍 `f89e891e…`，只在其后追加更正段），两支冻结文字两向 md5 不变。**真机 AC#2 半翻**：模型照宿主指针的第一发 `fs.read` 从 `kind=refused rules_hit=[R4] source=task.output / 读回 0 字节` 变成 `kind=success risk=L0 rules_hit=[]`、整份 20000 字节回全（台件背景正文钉回 179-r2 原版、一字未换）；**第二发（有界重读）仍被拒**，判中它的是 `fs.read` 读回来那 20000 字节所盖的**没有声明的** mark => 按票 183 AC#3/177 W-2 的边界它本就必须命中，这不是回归而是新缺陷面，修法要动 `internal/tools/bridge.go` 或 D15 再落盘层（都在本程写面之外，已具名停手上报，建议另立票）。AC#3 三处现量绿（包级归因换成 `web.fetch`、桥级 175-r2 那对腿 `ok 15.339s`、真机 `marks=2` 那一格）；AC#4 改一个 rune 的近邻路径仍命中；代价钉在常驻腿 `TestPointer183WorstCaseOfTheLandedExemptionIsPinned`（少判一枚卡，方向 fail-open：同 mark 正文里拼出声明串 >=8 rune 的片段不再是续读证据）。**真机频率那道必答题（现量，未修码 `git show HEAD:` 钉进 overlay 跑的）**：`3 of 10 bodies make the host's own pointer unreadable`，且响的片段不是夹具那枚 `-output-` 而是**用户 profile 前缀**（`\appdata\`、`\roaming`）=> 真内容也会响，不是只夹具响；真实 LLM 输出语料的发生率本程零读数。门禁：`./internal/risk/` 单包 `ok 4.078s`、tools/agent 绿、d22scan rc=0（`ban #8 internal/` 432＝多扫我那一枚 `_test.go`）、gate-clauses rc=1 且红腿名册**只有 `G6neg`**（与在册同名册）。AC#6 `窗口 0.0s` 只具名归口票 162，超时常量一字节没碰。
- 09-28 13:2x-13:5x 非实现者对抗验收 `183-v2`（**只裁不改**，表 `docs/evidence/s1/183-pointer-exemption-accept-v2.md`，读数 `.scratch/wisp/probes/183/v2/**`，**AC 框一枚没勾**，工具调用 38/40）。四格结论：① **恒真性矩阵**＝七枚常驻判据（现量 7 枚，不是派单写的 6）× 六发自取变异 M0/M1/M2/M3/M4/M5——**没有一枚是装饰**（每枚至少被一发打红；M1≡M2 红名册逐字相同＝值规则是一枚行为）；**但 M3（不要求 `runeIndexOf` 命中就挂 `declaredNorm`）零枚红** ⇒ `provenance.go:568-572` 那支 fail-closed 今天只靠注释活着，**判据集这一格不通过**，要补的腿形状＝正文不含声明路径时那一发必须仍命中。② **进攻**＝自造四发：与声明串共享 >8 rune 前缀的**另一条资源**路径仍命中（L1）、可逃集合精确＝"整条候选参数都是声明串的逐字子串"（L2 `hit=false`）、跨 scope 借不到（L3 归因 `web.fetch`）、同一路径第二发仍命中且判中它的是没声明的读回 mark（L4）⇒ **影响面只在同一枚 mark 内**，我这发没结出可用外带。频率口径复跑（同一 overlay、我自己 rc=0 那份 `freq-recount.txt`）：`3 of 10` 逐字复现，但分母＝**8 枚真仓内文件＋2 枚作者自造夹具**、另 2 枚 SPEC 语料因文件名取错未计入 ⇒ 该拆成"真文件 2/8＋自造 1/2"说，且那 8 枚是**仓内文件不是用户文档**，"自用第一天就会撞上"是〔推断〕不是读数（⚠ 段已披露组成，结论行与它口径不一致）。③ **三条边界**＝(a) 不是 AC#5① 禁的那一支（键不在参数名／通道上），但它**通道无关**地按值扣证据，可逃值只等于宿主自己公开过的那段字；(b) `declaredPath` **不是名册**——存储侧 `grep -rn "declaredPath\|attachDeclaredPath" internal/risk/*.go | grep -v _test.go` 现量 11 处全在 mark 侧、`Scope.Close` 删 mark 即删豁免、桥上 `hostPathBox` 每调用新建，跨 mark／跨 scope 不可见由 L3/L4 **实测**；(c) 逐 token 那一票＝**坐实**编排者"不算复活"（无 token 级污点状态、不传播、只逐字比一条常量串），并给了下次判分界线。④ **AC 判读（不勾）**＝票 183 **AC#2 判部分达成**：措辞上第一发满足，两处缺口＝(硬) AC#2 自己要求的"常驻、落 `internal/tools/`、长在 CLI 接缝形状上"**没落**——`grep -rln "183" internal/tools/*_test.go internal/agent/*_test.go` 零枚，tools 里只有 175-r2 那对桥级腿（本程现跑 `ok 0.078s`），CLI 接缝台件只在 `probes/183/**`＝非常驻；(归 185) 同一次任务第二发读回 0 字节。票 177 AC#3／票 175 AC#5／票 176 AC#3-5 **今天都不翻且理由各不相同**（177＝端到端条件只满足首读那一半；175 AC#5 是分账格、本票不构成为它翻勾；176 三枚等的是"起跑口"，与本票无依赖）。**票 185 判独立**（写面在桥／D15、性质是"越读越堵"不是 183 的回归），建议票 183 面追加限定句"逐字节读回指照宿主指针那一发"（只追加、不改原句）。没裁的十格逐名登记在表 §6（AC#1 四支、AC#4 不同目录同名形、AC#6 秒数、AC#7 完整轴、`-race`、真机 CLI 两形、门禁复跑、F2 两版等价性、M3 真机可利用性、`notify` 通道那一发）。护栏自证：零删除命令、零 push、六发变异每发 `git cat-file blob HEAD:<path> > <path>` 还原、终态 `internal/ cmd/` 空＋三向 md5 两向同值、`git log --oneline d61aee1b..HEAD -- internal/risk/` 现量无人碰过；M2 第一次跑是 `build failed`（`declaredNorm declared and not used`）＝**不是判据红**，改 `_ = declaredNorm` 重跑才取到真名册。
