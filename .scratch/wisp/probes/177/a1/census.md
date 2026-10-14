# 177-a1 — 只读设计核查 census：「给外来内容盖章」×「超长产物要能被模型续读」，**今天还互斥吗**

- 腿：`177-a1`（只读设计核查；零产码、零写生产面；⛔ 未翻框／未动 `Status:`／未碰台账／未 push）
- 起手锚 `af983cc4`（19:42 现读）→ 读数锚 **`a84fec49`**（19:46:57；漂移 1 枚＝编排者派本条腿的提交，`git show --name-only` 只含 `probes/orch/…/177a1-*.txt`＋台账，`internal/`＋`cmd/` 零命中——逐报）。
- 读数窗口 19:52–20:0x +08；交件＝本件＋`logs/risk-narrow.txt`｜`logs/tools-narrow.txt`（原始读数逐字）＋票 177 尾 1 条 bullet。

## 结论先行

1. **不互斥了——这枚排他今天不成立。** 而且化解它的**不是** `315`／`179` 两条修复（那两枚是同一根管子上的**另两扇门**，也都已落地）；真正拆掉排他的链＝**`177-r1`（甲形载具）→ `175-r2`（`task.output` 进 C25 名册＝盖章真生效）→ `183`（按值规则，第一发续读过）→ `185`（per-scope 派生名册，`A381` 批准，第二发也过）**。`175-r1` §7.1 当初撞红的两枚（`TestLongOutputPointerRecoversEveryByte`／`TestPointerPast256KiBIsNotFullyReadable`）**今天在真桥具上逐枚绿**（本腿现跑，逐字见 §2）。
2. **票 175 没有可派的写腿了**：`175-r2` 已在 09-28 落地（名册＋五枚常驻判据），票面 4/5 勾，唯一未勾＝AC#5「与 164/174 分账」——那是**结题动作**，不是产码活。派单里"175 可派"的前提**已过期**；它的"最小形"今天就是树上的实体（§3②）。
3. **真正可派的下一条（等 `cmd/wisp` 让位）＝`179` AC#8 的 CLI 端到端**：零产码、一条现成 overlay 命令（`probes/179/r3/00-anchor-and-gates.md` §4 已备），一发同时喂 **177 AC#3 端到端条件／175 AC#5／176 AC#3、#4**。**本腿没跑它**——`cmd/wisp` 面本波被 `133-r1` 占（窗口具名，见 §2）。

## 1. 复核排他前提（逐条现读；引内容锚，行号会烂）

- **`docs/evidence/s1/175-mark-gate-r1.md` §7.1（撞红原文）**：两枚红逐字＝`--- FAIL: TestLongOutputPointerRecoversEveryByte … the pointer must be re-readable with fs.read, got: {Text:L2 审批通道尚未接入（票 21），已拒绝执行 …}`／`--- FAIL: TestPointerPast256KiBIsNotFullyReadable … the first segment must read back, got {Text:L2 审批通道尚未接入 …}`。机制（它现量）：D15(3) 续读桩逐字嵌产物路径（`…全文见 %s…`）⇒ 盖章后那段路径进污染索引；续读的 `fs.read` 参数只有 `{"path": <它>}` ⇒ `Inspect` 对 path 形参数**不豁免** ⇒ 命中 ⇒ R4 升 L2 ⇒ 测试 `NoGate{}` ⇒「已拒绝执行」。四支候选出路（甲具名豁免／乙桩不嵌路径／丙接受续读要卡／丁切半句盖戳判不成立），**甲/乙/丙三形同撞**——它的原话：不是选形问题。
- **票 164 AC#3 两枚续读腿**：名称＝上句两枚（`internal/tools/task_output_leg_test.go`）。票面逐字「截断策略必须带可找回的指针…按那条做，不新造」；翻勾判语带三条件，其一逐字「指针能不能被模型**自己**走通，取决于**票 174**」。164 未勾＝AC#2b（署名）／AC#4（取消尾巴，对象在 176 AC#3）／AC#5（契约轴）。
- **票 176 的 e2e 格**：AC#1 [x]；AC#2–#5 [ ]。其 Progress log 的 e2e 读数逐字：前半截＝「跑后台任务→超长落盘→名册真有记录→路径能读回全文＝绿（第一次真机读数）」；后半截＝模型调 `task.output` 被 `internal/agent/loop.go` 的 default 支拒（逐字 `风险未分级且直通开关关闭，已拒绝执行`）⇒「票 177 那条『续读不被自家门拒』的端到端**仍没有读数**，而且现在挡在豁免之前的另一枚门上」。
- **315 收口件**：票名带 `-done`；`docs/evidence/s1/ticket315-315-v1-accept-v1.md` 六格全「成立」；修法＝装配根补一行（我现读在场：`cmd/wisp/run.go` 的 `rt.bridge.OpenTask(taskID)`）；复现读数 31s 拒 → 修后 0.14–0.20s 放行。
- **179 收口件**：`docs/evidence/s1/179-declared-l0-refused-accept-v1.md`＋票尾 2026-10-14 17:30 编排者亲验；**AC#2/#3/#4/#5/#7＋#9/#10 已勾**，**AC#1／AC#8 未勾**；修法我现读在场：`internal/agent/loop.go` 的 `case RiskL0:` 支；AC#8 前置逐字备在 `probes/179/r3/00-anchor-and-gates.md` §4（命令／overlay 件／预期读法／上一发 R4 基线）。

## 2. 量今天（仓外快照重放；读数逐字）

**环境与争用窗口（自证）**：
- 快照＝`git archive HEAD` → `/d/tmp/wisp177a1snap`（HEAD `a84fec49`；`git status --porcelain -- internal/ cmd/` 起手与全程为空）；`GOTMPDIR='C:\tmp\wisp177a1gotmp'`；`GOFLAGS=`；PATH MSYS 形；`nice -n 10`；原始读数列目录 `logs/`。
- **争用窗口（`tasklist` 自证）**：133-r1 的 `go.exe test -count=2 -v -timeout 40m ./cmd/wisp/`（PID 33936，起 **19:45:19**，`GOTMPDIR=C:\tmp\wisp133r1gotmp`）＋其 `wisp.test.exe`（PID 44924，19:45:27 起）＋其内嵌 `go.exe build -o …\wisp.exe ./cmd/wisp`（PID 26924，19:51:26 起）；**20:01 复读仍在**。⇒ 按派单：**`cmd/wisp` 台件（179-r2 overlay e2e）⛔ 未跑**；改在快照上重放同一形的**桥级／包级**常驻件（不占该面）。
- 另两枚旧进程（不属本波）：`wisp.exe`（15:31 起，机主试用件）＋`cargo.exe`×2（10-10 起）——与编译面无关，未动。`A952` 曾记 133-r1 在飞变异件（`M main.go` 等）；我 20:01 现读 `git status --porcelain -- cmd/ internal/`＝**空**（可能已还原），但它持续在跑 `./cmd/wisp/`，面仍被占。

**读数 A（`internal/risk`）**：尺＝`go test -count=1 -run 'TestPointer183|TestPointer185|TestShapeA' ./internal/risk/` ⇒ `rc=0`｜`--- PASS`＝**26／FAIL 0**（183 族 8＋185 族 10＋shape_a 族 8）｜逐字 `ok github.com/CarlosShao/wisp/internal/risk 1.792s`。点名含 `TestPointer185SecondRereadOfAHostArtifactIsClean`（第二发续读干净）、`TestPointer185TheReadBackBodyIsWhatBlocksTheRereadWhenNothingDeclaredIt`（**反向控制**：没有名册时读回正文确实阻断）、`TestShapeAW2ForeignMarkCarryingSamePathStillHits`（外来携带同路径仍命中）。
**读数 B（`internal/tools`）**：尺＝`go test -count=1 -run '<下列七枚>' ./internal/tools/` ⇒ `rc=0`｜**7 PASS／FAIL 0**｜逐字 `ok github.com/CarlosShao/wisp/internal/tools 6.206s`。逐名：
- `TestLongOutputPointerRecoversEveryByte`（164 AC#3 腿 1）`PASS (0.03s)`
- `TestPointerPast256KiBIsNotFullyReadable`（164 AC#3 腿 2）`PASS (0.30s)`
- `TestTaskOutputStampedOnRealBridge175r2`（J-1 盖章正向）`PASS (0.01s)`
- `TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2`（J-2 反向）`PASS (0.04s)`
- `TestHostMintedPointerRereadStaysCleanOnRealBridge175r2`（J-3 宿主指针续读干净）`PASS (0.03s)`
- `TestPointer183BackfilledArtifactRereadStaysClean`（CLI 接缝形、回填路径）`PASS (0.12s)`
- `TestPointer183SiblingArtifactInTheSameDirectoryStillHits`（同目录兄弟仍命中）`PASS (0.05s)`
⇒ **"模型续读宿主指针"这一形今天过门**（桥级）；两向都在场（正向干净＋反向仍命中），不是把墙拆掉换来的安静。

**对照：盘上旧读数（本腿复量的对象曾经长什么样）**——`probes/179/r2/logs/e2e-readings.txt`（09-28 11:3x）逐字 `tool=fs.read risk=L2 decision=reject … rules_hit=[R4] … reason="R4: 包含来自 task.output 的内容"`（`in_allowlist_scope=true`＝不是越界升档）；`probes/183/r1/logs/e2e-readings.txt`（09-28 13:0x）第一发 `kind=success risk=L0 … rules_hit=[]`／第二发 `kind=refused risk=L2 rules_hit=[R4] reason="R4: 包含来自 fs.read …"`（＝票 185 的面）。**post-185 的 CLI 读数：盘上不存在**（`185-r1` 交付表 §⑦ 待填；179 AC#8 等 `cmd/wisp` 空出）——**这一格是今天唯一量不到的**。

**静态链复核（现读，与读数同向）**：`bridge.go` `mark()` 早退＝`!risk.IsSensitiveSource(dec.Tool)`，而名册已含 `SrcTaskOutput`（九枚）；`task.go` 在 `rec.ArtifactPath != ""` 时 `box.set(rec.ArtifactPath)` 并写指针句；`provenance.go` `MarkWithHostPath` 三段＝(i) 177 位置跨度＋(ii) 183 `attachDeclaredPath` 按值扣掉声明路径自己的窗口＋(iii) 185 `hostMintedPathFor(scopeID, origin)`（**派生不存储**：由本 scope 自己 marks 的 `declaredPath` 投影，随 `Scope.Close` 生灭；只对"origin 恰为宿主亲手落盘那一条"的 mark 生效，**参数侧永不查表**）；`taintmatch.go` `contains()` 对拼出声明路径的候选窗口跳过。三枚边界都在测试面：外来同路径仍命中（W-2／J-2／185 L4）、兄弟／一字符之差／同 basename 仍命中（W-3／183 AC#4／185 AC#7）、名册不跨 scope 且随 marks 死（185 L5/L6）。

## 3. 裁定建议（编排者裁；本腿只给推荐）

**① 两件还排不排他**：**不排**。判据三向：(a) 排他机制（盖章→路径入索引→R4 拒续读）已被 183/185 按值规则拆掉，桥级两枚原位腿今天绿；(b) 反向面仍在（外来携带同路径仍命中）⇒ 不是靠放宽换的；(c) 盖章本身今天真生效（J-1 绿＋名册九枚）。⚠ 两点保留：(i) 化解者是 `177→175-r2→183→185` 链，**不是** `315/179`（后者是同一管子上更后面的两道门——"子代理 scope 未注册⇒白名单内读被假 L2 拒"与"声明 L0 当未分类⇒只读内置被拒"——都已落地，是 e2e 通路的**前置**而非排他机制本身）；(ii) **真机 CLI 端到端今天仍以 09-28 的读数为准**（post-185 未复量）⇒"不排他"目前是**桥级/包级事实＋静态链**，不是"CLI 已复绿"。

**② 若不排他 ⇒ 票 175 是否可直接派**：**175 的写腿不用再派——它已经落了。** 若编排者要的是一句"最小形"的实体清单，今天树上的就是：
- **判据**：`internal/tools/ticket175r2_stamp_live_test.go` 五枚——J-1 正向／J-1b 非内容不盖戳／**J-2 洗戳警报器（外来同路径仍命中 R4、`SessionOverrideBlocked=true`）**／J-3 宿主指针续读干净／普查腿（每枚注册工具必须被分类）。
- **落点**：`internal/risk/provenance.go`（名册九枚，`+SrcTaskOutput`，`+4−0`）＋`internal/tools/task.go`（既有 `box.set`，本批未改）＝`962ba9d1`；判据件＝`58db7400`。
- **face**：`internal/risk`＋`internal/tools`（**不含 `cmd/wisp`**）。
- **要碰的既有钉：零枚**——`task_output_leg_test.go` 两枚、`provenance_test.go`、`thresholds.go` 一字节未动（票面自证＋本腿读数复认）。
- 风险一句话：豁免今天只对 `task.output` 一枚（全仓唯一 `box.set`）生效，下一枚带宿主路径的工具没登记就靠普查腿红——这正是 175 AC#3 买的东西。
- **175 唯一未勾＝AC#5 分账**，前置＝端到端读数（见③）。

**③ 若仍排他（本条件未触发）：不适用。** 但把"e2e 分叉怎么办"备好（probe 内建分岔清单）：`风险未分级`（指回 179 面）／`L2 级…审批通道尚未接入`（指回票 21/162 族）／`R4`（指回 177/183/185 面）／`查不到这个任务`（指回 176 写者）。**下一枚可派腿（推荐）＝`179` AC#8 的 CLI e2e**：零产码，一把尺＝`go test -count=1 -v -run 'TestRereadHostPointerOnCLISeam179r2' -overlay=.scratch/wisp/probes/179/r2/overlay-e2e.json ./cmd/wisp/`（带在册 PATH 前缀）；**face＝`cmd/wisp`**（等 133-r1 落地）；一发同时结：**177 AC#3 端到端条件／175 AC#5／176 AC#3、#4**。⚠ 台件最后实跑＝09-28，**对现树的兼容性必须先编译一验**（`cmd/wisp` 此后动过 313／255-r8／133 等多轮）。

## 4. 依赖与风险（派 175 后继／179 AC#8 所需全部前置）

- **`cmd/wisp` 面**：133-r1 独占中（窗口见 §2）⇒ e2e 腿的硬前置＝该面空出＋台件编译复验。
- **174 家族**（"读回来"同一根管子）：AC#1/#2/#3/#2d/#6 已勾；**未勾＝AC#2b**（把 `Paths:` 接进生产——179 票尾现读仍写 `TaskDeps{Roster: rt.tasks}` 未接，面 `cmd/wisp`）、**AC#2c**（`internal/agent/spill.go` 那枚桩的说法；票面逐字「**排程：等票 177 裁完再动**」——**今天可判它解锁了**，面 `internal/agent`）、AC#4/AC#5（结题动作）。`Q-59`（`fs.read` 偏移）＝仍 DEFERRED（`PLAN.md:1532` 那行的报警器＝`TestPointerPast256KiBIsNotFullyReadable`，今天绿＝"没人偷偷加偏移"，属正常）。
- **167／267**：与本题**不相关**——167＝面板四枚小出口（立而不派，面在面板侧）；267＝审批窗口配置值域，已 `-done`。
- **今天不可达的格（逐名）**：① **179 AC#8／177 AC#3 端到端条件／175 AC#5／176 AC#3、#4**——同一发 CLI e2e，本波量不到（`cmd/wisp` 占中，具名）；② **177 AC#0（丙形）**——**零读数，且不许被本腿任何结论冒充**（本腿只量已批的甲链）；③ **185 证据表 §②–⑪ 待填＋AC 全未勾**（决定性修法 `ee1e2118` 代码在树，交付表是残件）——建议派 e2e 腿时顺手记一笔，别让"最后的牙"只有代码没有交付面。
- **风险**：豁免代价面（正文恰拼出路径前缀的段不再算证据；183 频率尺 3/10、真实语料零读数）＝已由 `TestPointer185CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan` 钉住，属既知既批；`A952` 记的 133-r1 在飞变异件提醒：**它落地前任何人不许拿工作树当读数面**（本腿用 `git archive HEAD` 快照正是为此）。

## 具名没做（别把上面的绿读成"都量过了"）

- ⛔ **未跑 `cmd/wisp` 任何东西**（含 179-r2 overlay e2e）——争用窗口见 §2；"今天过门还是被拒"的 **CLI 一格无本腿读数**，只有桥级/包级＋静态链。
- ⛔ 未跑全包（只窄 `-run`）；未跑 `d22scan`／`gate-clauses`／`runtests`／`flip-declaration.sh`（本腿零产码改动；那些尺面属别人）。
- ⛔ 未读 `frontend/**`、`design/**`；⛔ 未翻框、未动 `Status:`／票名、未碰台账；⛔ 未 push。
- ⛔ 未复算 `185-r1` 的变异读数与 §⑦（其表残件，留它的主子）。
