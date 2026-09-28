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
