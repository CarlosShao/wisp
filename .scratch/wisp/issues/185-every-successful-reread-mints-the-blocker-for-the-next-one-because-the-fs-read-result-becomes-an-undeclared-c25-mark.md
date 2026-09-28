# 185 — **每一次成功的续读，都会亲手造出下一次续读的阻断者**：`fs.read` 读回来的那份正文被桥当场盖成一枚**没有声明路径**的 C25 mark ⇒ 模型对**同一条宿主指针**再发一发（有界重读／续读后半段）就又被 R4 拒（真机逐字 `reason="R4: 包含来自 fs.read C:\Users\…\001\a…"`、close 时 `dropped=2`）——**票 183 修好的是"第一发能读回"，不是"读得回来这件事成立"**

- Status: **已立、按住不派**（等 `183-v1` 交完并由我裁归属：本票与票 183 AC#2 的措辞是同一场官司）。
- 来源：`183-r1`（提交 **`0662a35a`**，表 `docs/evidence/s1/183-pointer-exemption-r1.md` 24874 字节）。它把票 183 那枚豁免落成"按值扣掉声明路径自身的窗口"，**第一发续读真机翻绿**，然后**照票面停手上报**了这第二发——**没有顺手去动桥**（桥在它的禁改面上），这一格是正面样本。台账 `A364`。
- 关联：**票 183**（豁免本体，同一接缝）· **票 177**（W-2：外来内容携带同一条路径**仍要命中**＝本票的边界，不是本票要拆的墙）· **票 175-r2**（桥级盖章那一跳）· **票 164**（截断带可找回指针）· **票 176 AC#3-5／票 177 AC#3 端到端条件／票 175 AC#5**（这几枚条件格等的都是"读回来"，本票不解决它们就翻不了完整版）· **D15**（上下文预算／再落盘层）· **C25**（盖章契约）。

## 这是什么（人话）

上一发修好的是："模型拿着宿主自己写下的那条路径去读文件，不再被自家门挡住。" 真机读数（`probes/183/r1/logs/e2e-readings.txt` 第 **55** 行，逐字）：

```
[audit] tools: call kind=success tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true reason="无规则命中（L0 直接执行）"
```

可是同一次运行里，模型对**同一条路径**再发的一发（有界重读）仍然被拒（同一枚文件第 **60** 行，逐字，路径已按本仓规矩截短）：

```
[audit] tools: call kind=refused tool=fs.read risk=L2 decision=reject outcome=error rules_hit=[R4] in_allowlist_scope=true
          reason="R4: 包含来自 fs.read C:\Users\swq\AppData\Local\Temp\TestRereadHostPointerOnCLISeam183a1929971568\001\a…"
```

⇒ 判中它的**不是** `task.output` 那枚戳，是**上一步 `fs.read` 成功读回来的那份正文自己**被盖成的那枚戳（尺：`grep -an "kind=refused" .scratch/wisp/probes/183/r1/logs/e2e-readings.txt`；旁证第 **64** 行 `C25 scope closed … dropped=2`＝收尾时该 scope 里**两枚** mark，比票 183 那一发多一枚）。

**形状**：豁免是**逐枚 mark 声明**的（`MarkWithHostPath` 只把宿主当场声明的那条路径的窗口从**那一枚** mark 的证据里去掉，随该 mark 生死）。`task.output` 那一枚声明了 ⇒ 第一发通。`fs.read` 的正文那一枚**没人替它声明** ⇒ 它的正文（里面就有那条路径的同款片段）成了下一发的证据。**于是"读回来"这件事今天只对第一次成立。**

## 现量（锚 `0662a35a`，09-28 13:0x 编排者本程现跑）

1. 上面两行读数（第 55／60 行）＋ `dropped=2`（第 64 行）。尺：`sed -n '55p;60p' .scratch/wisp/probes/183/r1/logs/e2e-readings.txt`、`grep -an "dropped=" 同一枚文件`。
2. **谁给 `fs.read` 盖的戳**（尺：`grep -n "func (b \*Bridge) mark" -A 22 internal/tools/bridge.go`）：桥在每次执行后拿回执正文盖一枚 mark，**只有** `MarkWithHostPath` 那一支会带声明路径，而带不带由 `hostPathBox` 里有没有值（尺：`grep -n "hostPathBox\|withHostPathBox" internal/tools/bridge.go internal/tools/task.go`）⇒ **`fs.read` 这一路今天不填那枚 box**（填它的是 `internal/tools/task.go` 里 `task.output` 那支）。
3. **为什么这不是票 183 的回退**：同一发在未修码上第 55 行也是红的（`183-r1` 的 HEAD 读数 `probes/183/r1/logs/referee-on-HEAD.txt`：包级裁判与"路径出现两次"两枚都 FAIL）⇒ 第一发是**它翻过来的**，第二发**在票 183 之前也在**，只是当时被第一发的红挡住了看不见。
4. **⚠ 一枚要摊开的次生面（本票与票 183 共用，别只报一半）**：`183-r1` 量的真机频率读数（尺：`cat .scratch/wisp/probes/183/r1/logs/freq-head.txt`）＝**10 份正文里 3 份**会让宿主自己的指针读不回来，而响的片段全部是**用户 profile 前缀**那一类（`\swq\app`／`\appdata\`／`\roaming`）。⚠ **这 10 份"正文"是拿仓内文件当语料**（`AGENTS.md`、`docs/PLAN.md`、`internal/agent/loop.go`、`go.mod`、两枚 synthetic，另两枚 SPEC 文件名它猜错＝未计入），**真实用户文档与 LLM 输出语料的命中率＝零读数**——写手自己在"本程没测什么"里具名了。⇒ 同一把尺的另一面：**落地那枚按值豁免，也会让同一枚 mark 正文里恰好拼出路径前缀的那段外来正文不再算证据**（写手自陈"少一张 L2 卡，索引从不变宽"，并把最坏形状钉成常驻判据 `TestPointer183WorstCaseOfTheLandedExemptionIsPinned`）。⇒ **"读得回来"每往前提一步，"外来内容算不算证据"就往后退一寸**——这两面要一起报，不许分开报成"修好了"或"有洞"。
5. **三行定性**（安全字样自带）：**①现象出现在哪**——只出现在本仓自己那枚 CLI 接缝台件与 `internal/risk` 的包级判据里；**②有没有本机被入侵的证据**——没有，一条都没有；**③最坏后果是什么形状**——不是"有人进来了"，而是**两张相反方向的错**：该拦的没拦（外来正文里拼出路径前缀那一段漏判），或不该拦的被拦（模型正常续读被升 L2 拒掉，用户看到的是"产品读不回自己的输出"）。两者都是**产品正确性**问题，不是攻击面。

## 为什么值得做（不做会怎样）

owner 从第一天在意的那句"超长的输出，我能不能让它自己读回来"——今天的答案从**不能**变成了**能读一次**。真实用法里"读一次"通常不够（截断后分片读、读一半再要后半段、模型重试同一发），而每一次成功读回都会**新增**一枚没声明的 mark ⇒ **越读越堵**。不做这一枚，票 176 AC#3-5／票 177 AC#3／票 175 AC#5 那几枚"端到端条件格"就只能停在"第一发成立"这一半。

## AC（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）

- [ ] **AC#1 先把"第二发为什么被拒"钉成一枚常驻判据（红在 HEAD）**：落 `internal/tools/`（走真桥具）或 `internal/risk/`（包级两枚 mark 那一形），要求"同一任务里 `task.output`→`fs.read` 成功→对**同一条路径**再发一发 ⇒ 不命中 R4"。⚠ 判据要**成对**：同批补一枚"外来正文携带同一条路径 ⇒ 仍命中"（票 177 W-2 原句），缺一形＝把豁免又扩成装饰或把墙又拆掉。
- [ ] **AC#2 答"这枚声明该由谁给"**，逐支给现量与代价，不许凭偏好选：**(a)** 桥侧——`fs.read` 成功时也填 `hostPathBox`（它读的就是那条路径，宿主当场知道）；**(b)** D15 再落盘层——续读走"宿主替模型读"的那条通道，根本不产生第二枚 mark；**(c)** **不修**，把"同一路径只许读一次"写成明说的事实与文案。⇒ 判"(a) 或 (b) 会碰到票 177 已裁过的边界（per-scope 名册／参数侧按值放行／`allowlist`）"＝**停手上报**。
- [ ] **AC#3 覆盖面要量，不要宣布**：落地后**逐枚**答"哪些工具回执里会带宿主路径、它们各自会不会造出同款阻断者"（已知候选：`fs.read`／`task.output`／`clip.*`／落盘回执）。尺要现跑（`grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/`），**名册与枚数不许从本票抄**。
- [ ] **AC#4 与票 183 共用那枚代价面**：AC#2 无论选哪支，都要复量第 4 把尺那一面（真实正文命中率 3/10 那把尺）并具名说"这次往前推一步，让哪一类外来正文不再算证据"。**不许**只报收益不报退让。
- [ ] **AC#5 门禁**：逐包 `go test -count=1 ./internal/risk/`（**单跑**）＋`./internal/tools/`＋`./internal/agent/`；CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，两形都要贴）；`sh scripts/d22scan.sh`；`gate-clauses.sh` **比红腿名册不比退码**；⚠ 终态读数取在最后一枚 commit 之后。⚠ **排程护栏**：`-overlay` 不保护你不吃别人的在飞脏件——起手复量 `git status --porcelain -- internal/ cmd/`，不为空就把你要依赖的每枚跟踪件用 `git show HEAD:<path>` 钉进自己的 overlay；**不许跑** `probes/161/r6/flip-declaration.sh`；注释里凡自指本文件行号，写"符号名＋`grep -n` 尺"，**不写死号**。
- [ ] **AC#6 契约轴**：`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节不许动；`internal/risk/provenance.go:468-473`（`Mark` 契约文字）与 `internal/risk/taintmatch.go:11-15`（逐 token 追踪已被否决）保持不动（基线 md5 在 `A363`／`A364`）；**票 183 的文档块原 25 行同样不许改写**，要更正只许追加。
- [ ] **AC#7（09-28 14:1x 由票 183 的 AC#4 归口过来·纯测试面·零产码）**：`同名不同目录` 那一形——同一个**文件名**挂在另一枚目录下、宿主正文里逐字拼出那条兄弟路径、而宿主只声明了自己那一条 ⇒ 续读兄弟**必须仍命中 R4**。⚠ 它和 `183-r2` 已落的那枚"同目录兄弟"是**两枚不同的洞**：那枚防"放宽到目录级"，这枚防"放宽到**文件名／basename**级"。今天修法是按整条归一化串逐字比、看着到不了这一形——**但"看着到不了"不是凭据**（票 183 整张票的教训就是这句）。⇒ 本格今天**应当是绿的**，所以它是**"不许弄坏"的守卫，不许充当本票新增的牙**；它的牙必须现量：对"豁免改成只比 basename"那一发变异（走 `-overlay`，别动跟踪件）它要红。落点＝`internal/tools/pointer_183_cli_seam_test.go` 同批加一枚，或 `internal/risk/pointer_183_test.go` 第 9 枚（尺 `grep -c "^func Test" internal/risk/pointer_183_test.go`＝**8**，别抄票面）。

## 本票**不**解决

- 不裁"该不该给续读加一张 L2 卡"（票 177 丙形，AC#0 至今零读数）。
- 不动 `窗口 0.0s` 那一族（归票 162，票 183 AC#6 已具名）。
- 不动票 174 AC#2b（回执里"C26 未接线"那句文案，正交）。
- 不翻任何别人的勾（票 176 AC#3-5、票 177 AC#3 端到端条件、票 175 AC#5 **维持现状**），不 push。

## Progress log

- 09-28 13:0x 编排者立票：上面 5 把尺本程现跑（锚 `0662a35a`，`git status --porcelain -- internal/ cmd/` 起手为空）。第 55／60 行两枚读数我自己从文件里逐字取。⚠ 本票与票 183 AC#2 的"一次成功算不算达标"是同一场官司，**归属等 `183-v1` 交完再裁**；裁完之前**不派**。未派。
- 09-28 14:1x 编排者：**归属已裁完、本票解除按住**（`183-v2` 判"独立"，我照它翻勾并把票 183 改名 `-done`）。本程真机重跑再确认一次阻断者仍是这枚（HEAD=`78ea1d19`，带在册 PATH 前缀）：第一发续读 `kind=success risk=L0 rules_hit=[]`，第二发 `kind=refused risk=L2 rules_hit=[R4] reason="R4: 包含来自 fs.read …"`，台件 `TestRereadHostPointerOnCLISeam183a1` 在 HEAD 上 `--- FAIL (60.10s)` 红在 `zz183r1_e2e_test.go:373`。⚠ 同批**新加 AC#7**（票 183 AC#4 那支"同名不同目录"零读数、原话不改、归口到本票）。⇒ 下一枚＝**`185-c1`（只读普查，先答 AC#2 那三支的现量与代价，不许直接派落地腿）**。未派之前不改任何勾。
- 09-28 14:25-14:43 只读普查 `185-c1` 交件（表 `docs/evidence/s1/185-reread-owner-census-c1.md` 31185 字节，探针与四份 overlay 在 `probes/185/c1/**`，**AC 框一枚没勾**，工具调用 35/40）。**本票的"未修码"它具名为起手 HEAD `72c76d42`**（并自证票 183 的修法已在场：`spellsDeclaredPath`／`attachDeclaredPath` 在 `taintmatch.go` 有定义有调用）⇒ 这一格按我上一轮的教训做对了。**四支结论**：**(a) 桥侧填 box**＝落点 `fs.go` 的 `fsRead.Execute`、值取 `canon`（＝模型递进来的那个 `path`）⇒ **它判上去就是票 183 AC#5① 明禁的"参数侧按值放行"**，且我本程**自己复跑**它的变异（`-overlay=probes/185/c1/overlay-muta.json`，`./internal/tools/` 全包 `-v`）⇒ **`ok 17.264s`、`--- FAIL` 计数＝0** ⇒ **这一支今天落地没有任何一枚已入库判据会响＝零牙**（不是它安全，是没人看着）；**(b) D15 再落盘层**＝今天不存在"宿主替模型读"的通道（`internal/agent/*.go` 里 `os.ReadFile|os.Open` 零命中），任一落法都要动**票 175"成功结果必须有来源标记"那枚承重声明**（(b-ii) 摘戳会红三枚：`TestFSReadReturnsTheFileAndTaintsIt`／`TestFSReadTaintFeedsR4`／`TestEveryRegisteredToolIsClassifiedForMarking175r2`）；**(c) 不修、写成"同一条路径只许读一次"**＝非冻结落点存在（`task.go` 的 `pointerNotice`），但**这句话今天不真**——第二发响不响取决于正文巧合（票 183 频率尺 3/10、真实语料零读数）⇒ **写成规则＝拿偶然当承诺**（票 97/179/183 同族：注释说的动作没有判据钉着）；**第四支（它自己找到的）**＝判据换成"这条路径是不是本 scope 里宿主自己落盘过的那一枚"（宿主确实登记过：`task_backfill.go:112`），⚠ **但它的本体＝per-scope 路径名册＝票 177 已裁"未批"那一形**（`taintmatch.go` 逐字写着 "a per-scope path roster is explicitly NOT approved"）。**AC#7 那一形它验了不是推理**：我本程复跑 basename 放宽变异（`overlay-ac7-mutant.json`）⇒ **只有它那枚探针红、`internal/risk` 已入库的 16 枚（`pointer_183` 8＋`shape_a` 8）一枚都没红** ⇒ **"同名不同目录"今天在判据面上零牙**，落地腿必须自己补那枚常驻判据。**名册现量**：全仓只有 `task.output` 一枚 setter（`box.set` 唯一命中 `task.go:254`）、`MarkWithHostPath` 生产调用唯一一处（`bridge.go:575`）⇒ **今天带声明的工具＝1 枚**；顺带量出一枚反向面：`fs.list` 回执必带宿主全路径却**不在**盖戳名册里（少一张证据、不是多一枚阻断者，本票不裁、具名转下一轮）。纪律自证我复量同值：`internal/ cmd/` 起手与终态都空、三枚受保护文字 md5 全等基线、`ban #8 internal/` 433、红腿名册仍只 `G6neg`、**AC 框 7 枚未勾 0 枚已勾**。**⇒ 本票四支全部要么撞已批边界、要么需要一枚新裁定，落地腿我不派**；已把"要不要为分页续读推翻 09-27 那条名册裁定"摆给 owner（人话后果＋我的推荐，见台账 `A370`）。**只 commit 两枚（`c1c96008` 骨架、`4813567e` 正件），没 push。**

- 09-28 14:2x–14:4x `185-c1`（只读普查程）交件：**产码面零字节改动**（起手与终态 `git status --porcelain -- internal/ cmd/` 都为空；三枚受保护文字 md5 起手与终态均＝基线 `858e4511…`／`5680ddd1…`／`f89e891e…`）。**本票的"未修码"具名＝本程起手的 HEAD `72c76d42`**（票 183 的修法已在场：`grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/taintmatch.go | grep -v _test.go` 命中；派单写的 `78ea1d19` 经 `git merge-base --is-ancestor` 验为 HEAD 的祖先）。裁决表＝`docs/evidence/s1/185-reread-owner-census-c1.md`（十二节齐；探针与变异全在 `.scratch/wisp/probes/185/c1/**`，走 `-overlay`）。**AC 框一枚没勾。** 三格读数摘要：
  - **AC#2(a) 桥侧填 box**＝落点 `fsRead.Execute`（`internal/tools/fs.go`）；全仓今天只有 `task.output` 一枚 setter（尺 `grep -n "hostPathBoxFromCtx" internal/tools/*.go | grep -v _test.go`）。变异现量（`overlay-muta.json`）：`go test -count=1 -v -overlay=… ./internal/tools/` **rc=0／零枚红**，且用"故意写坏的副本 build rc=1"证了 overlay 确实生效 ⇒ **这一支今天没有任何判据管得着**；本程判它**等价于"按参数值放行"（声明值就是模型递进来的 `path`）＝碰票 183 AC#5①＋票 177 W-2**，**具名停手上报，未放行、未选它**。
  - **AC#2(b) D15 再落盘层**＝那层在 `internal/agent/spill.go`（`Prepare:84`）＋唯一调用点 `loop.go:702`；**今天不存在"宿主替模型读"的通道**（尺 `grep -rn "os.ReadFile\|os.Open" internal/agent/*.go | grep -v _test.go`＝空）。"不产生第二枚 mark"两形：新通道／摘 `fs.read` 的戳——**都动票 175「成功结果必须有来源标记」的承重声明＝禁令面**，点名会红的既有判据：`TestFSReadReturnsTheFileAndTaintsIt`／`TestFSReadTaintFeedsR4`／`TestEveryRegisteredToolIsClassifiedForMarking175r2`（两形射程不同，见表 §③；本程**未做变异现量**）。
  - **AC#2(c) 不修只写文案**＝非冻结落点存在＝`internal/tools/task.go` 的 `pointerNotice`；桩格式串本体（`task.go:257`、同款字面另在 `spill.go:141`）碰 `PLAN.md:2564` 那行的定案形状（现量：`grep -n "全文见" internal/tools/task.go docs/PLAN.md` → PLAN.md 里**没有**"全文见"这个字面）。⇒ 落桩本体＝**要人工批准**；另报一处真代价："只许读一次"**今天不真**（第二发响不响＝正文内容巧合，票 183 频率尺现量 3/10、真实语料零读数）。
  - **第四支**：码上有形（判据＝"这条路径是不是本 scope 里宿主自己落盘的那一枚"，宿主确实登记过：`task_backfill.go:112 rec.ArtifactPath = sp.Path`），**但它的本体就是 per-scope 名册＝票 177 已裁（未批）那一形**（`provenance.go`／`taintmatch.go` 各有一句"a per-scope path roster is explicitly NOT approved"）⇒ 具名上报，未实现未测量。
  - **AC#3 覆盖面**：非测试命中 9 处，全仓只有 `bridge.go` 一处调 `MarkWithHostPath`、只有 `task.go:253` 一处 `box.set` ⇒ **带声明的只有 1 枚工具**；造同款阻断者的今天只有 `fs.read` 的正文。⚠ 顺量出一枚反向面：`fs.list` 回执**必带宿主全路径**（`joinForListing`）但**不在** `sensitiveSourceTools` ⇒ 桥根本不盖戳（少一张证据，不是多一枚阻断者），本票不裁它、具名给编排者。
  - **AC#7（验了，不是"看着到不了"）**：比对面是**整条归一化串**（`spellsDeclaredPath` 逐 rune、不查哈希）。配对读数＝`probes/185/c1/pointer_185_c1_ac7_probe_test.go`（只走 `-overlay`，跟踪件没动）：未变异 **rc=0 两枚 PASS**；加"basename 级放宽"变异 → **同一枚探针 FAIL、而整包 risk 已入库判据（`pointer_183` 8 枚＋`shape_a` 8 枚）一枚都没红** ⇒ **这一形今天在已入库判据面上＝零牙**，AC#7 的常驻判据必须由落地腿新写（现量 `grep -c "^func Test" internal/risk/pointer_183_test.go`＝8）。
  - **门禁（只作现状记录，按 `A363` 不充当 AC 结案凭据）**：`d22scan` rc=0／`ban #8 internal/` examined=**433**＝基线；`gate-clauses.sh` 名册 14 腿、BAD＝**1 枚＝`G6neg`**（与在册红腿名册差集＝空，**没跑** `flip-declaration.sh`）；`./internal/risk/`（单跑）／`./internal/tools/`／`./internal/agent/` 三枚全 ok。**真机 CLI 那一发本程没重跑**——它的 sink 会原地重写 `probes/183/r1/logs/e2e-readings.txt`（别家票逐行引用第 55／60／64 行），且该腿只以 `probes/183/r1` 的 `-overlay` 副本存在（`ls cmd/wisp/*_test.go` 没这枚文件）；不带 PATH 前缀那一形同具名未取。两形的 HEAD 读数出处＝本票上面那条 14:1x 的编排者亲取记录。
  - **零删除命令／零 push／未派落地腿**。人工批准项本程只摆不拍：①(a) 是否算"按参数值放行"＋破 W-2、②(c) 文案落桩本体是否动 2564 的形状、③第四支要不要 per-scope 名册、④(b) 要不要拆 175 的承重声明。**待编排者裁：选哪一支、以及要不要先落批准记录。**
