# 票 174 · AC#2c · 写码腿 `174-r2`（第二任 r2；本格＝`internal/agent/spill.go` 那枚桩的回执）

> 派单：编排者 2026-10-03 12:2x 口头派单（本格＝票面 `:36` 的 AC#2c，二选一 ⓐ/ⓑ，⛔ 不许两形都不落）。
> 写面＝`internal/agent/**`（产码＋`_test.go`）。AC 框一枚不勾。只 commit 不 push。

## 0. 起手锚与脏名册

- `date -Iseconds` 起手自取：`2026-10-03T12:32:49+08:00`
- 分支：`dev`
- `git log -1 --format=%H`＝`796b33eadf05bdb3ef9120506392983401d9c10e`（锚点自己取的，不引别人的号）
- **`git status --porcelain -- internal/agent`＝空**（逐字复跑，输出 0 行）⇒ 本格写面起手干净，无撞腿。
- 整仓 `git status --porcelain` 非空（`cmd/wisp` 4 枚 M＋若干未跟踪、`internal/llm`／`internal/config` 各 1 枚未跟踪、`design/**` 一片）——**全部在我禁改面内，我一字节不碰**。
- 起手门禁基线（时刻 `13:36 +08`，同机另有腿）：
  - `sh scripts/d22scan.sh` → `rc=0`
  - `go test -count=1 ./internal/agent/` → `ok github.com/CarlosShao/wisp/internal/agent 1.668s`
  - 顶层现绿名册已落 `logs/roster-start.txt`（**73 枚** `--- PASS/FAIL/SKIP` 顶层行，收尾 `comm` 逐名比）

## 1. 复跑编排者的现量（三句全是待验断言）

尺与逐字读数见 `logs/claims.txt`。**结论：两句成立、一句要更正。**

| # | 编排者的断言 | 我复跑 | 判定 |
|---|---|---|---|
| 1 | `internal/tools/task.go:562` 逐字 `"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token%s，全文见 %s…]\n%s"` ＝ 多一枚 `%s` 槽 | `grep -n '输出已落文件' internal/tools/task.go` → `562:` 该行逐字相符 | **成立** |
| 2 | `internal/agent/spill.go:141` 逐字同句但**没有那个从句** | `grep -n '输出已落文件' internal/agent/spill.go` → `141:` 该行逐字相符（`token，全文见`，无 `%s`） | **成立** |
| 3 | "`spill.go:18` 一带自陈是 stub" | `:18` 逐字是 `// D15(3) long-output spill (SPEC-05 §4.2): a single tool result over the` ——**那一带没有"stub"这个词**；自陈为 stub 的是 `:47`（`// Text is what goes into the context: either the original text or the` / `:48` `// head/tail stub.`）与 `:51`（`// Spilled is true when the full text went to a file and Text is a stub.`） | **推翻（行号错，实质对）**：`:18` 是 D15(3) 落盘契约的自述头，"自陈是 stub"在 `:47-48`／`:51`。本格不因此改变结论 |

> 另：`grep -rn "全文见"` 现量整仓**产码**只有那两枚（`task.go:562`／`spill.go:141`），其余命中全在 `_test.go` 的手写夹具与 `docs/evidence/**` 里 ⇒ "同一句在两条腿上"这一维成立。

## 2. 选形与理由：ⓐ（让 `spill.go` 走同一判据）

### 2.1 排程前提我先自己读了（AC#2c 的"等 177 裁完"）

**177 已裁完、且落点已定**，三条凭据（都现读，不引转述）：

1. `Q-61`＝**甲（精确豁免）** owner 已批 —— `docs/reports/pending-and-issues.md:1101` 末段逐字「**09-27 21:3x owner 已批＝甲（精确豁免：只免"本宿主在这一任务里亲手写下过的那几条具体路径"）**」；`Q-63`＝甲亦已批（同文件 `:1100` 末段、台账 `A353`）。
2. 甲形**已真落地**：`internal/tools/task.go:558` 逐字 `if box := hostPathBoxFromCtx(ctx); box != nil {` ＋ 其上注释逐字 `// Ticket 177 shape A (Q-61甲, precise exemption)`；常驻判据件 `internal/risk/shape_a_exemption_test.go` 在盘（票 177 `:94` 编排者自己数过八枚用例）。
3. **豁免落点＝"在宿主自产的那一枚 mark 里，那条路径所占的窗口不入片段索引"**，载体是**逐次调用的 box**（`box.set(rec.ArtifactPath)`），**不是**目录级、**不是**配置项（票 177 `:57` 第 4 条逐字）。

⇒ 排程前提成立，可以动文案。**而且正因为读了落点，我才敢说本格不碰它**：那枚豁免的输入是 `rec.ArtifactPath` 这一枚**路径值**，走的是 `internal/tools` 的盖戳链；`spill.go` 的桩文本从来不经那枚 box。⇒ 见 §2.3 的"零新路径"判据。

### 2.2 为什么是 ⓐ 不是 ⓑ

ⓑ 要求我写明"177 的豁免落点会改变它的文本形状"。现读否掉这句的因果：177 争论的**当事文本**是 `task.output` 那一支嵌路径的桩（票面 `:15` 与 `Q-61` 逐字点的都是 `task.go` 那枚），而 `spill.go` 的桩**今天没有第二枚已裁的形状在等**；ⓐ 不需要动 177 任何一个字节就能把话说全。**选 ⓑ ＝ 把一句我没证成的因果写进契约面**，且 AC#2c 明写"两形都不落不许交回来当完成" ⇒ 我落 ⓐ。

### 2.3 ⓐ 的形状（四句话）

1. **判定走 C26，不在 agent 里手工规范化**：新增**消费侧窄接口** `PointerJudge{ Canonicalize(string)(string,error); InAllowlist(string)bool }`，声明在 `internal/agent` 本包内。⇒ **零新依赖边**（`package agent` 现不引 `internal/risk`，也不为这一格去引；`*tools.PathCanonicalizer` 结构上直接满足，与 `internal/risk/assessor.go:139` 那枚同名同签的接口是同一族既有定式）。⛔ 无 `filepath.Clean|Abs`（AGENTS §1.2 逐字禁止项）。
2. **成句定式复制不抽公共**：从句插在 `约 %d token` 与 `，全文见 %s` 之间的同一位置（`task.go:562` 那一枚 `%s` 槽的位置），**不共享字符串常量、不跨包抽函数**（本仓既有裁定＝解析器复制不抽公共；另 `internal/tools/task.go` 我一字节不改）。
3. **两支说话、两支之间不粘连**：`judge == nil`＝fail-closed 明说"无法核实，按读不到处理"；`Canonicalize` 出错／`InAllowlist` 为假＝各说各的。**健康正例（根内）零说明**，并另钉一枚"逐字节等于改码前模板"的钉 ⇒ 防"把说实话写成噪音"。
4. **两处分歧是刻意的，都有具名理由**：
   - 本腿**不**拼 `err.Error()` 进模型可见文本 —— 那正是 AC#2d 在 `task.go` 上量出的无界形状（票面 `:62` 现量逐字 `（"+err.Error()+"）`，漏了 `tools: empty path` 六个字）。我这一支从第一天就不开这个口。
   - 本腿的"根外"那一支说出**两条**回来的路（加 `[fs] allowed_dirs`，**或**批那一张 L2 卡）—— 编排者 22:3x 第 3.2 条对 AC#2b 的落地腿就是这个要求；`task.go` 那两句受 `:299`／`:349` 两枚模板冻结钉系着（＝`Q-63` 边界，等谁裁），而 `spill.go` 的桩**没有**任何模板冻结钉（现量：`grep -rn "输出已落文件" --include=*_test.go internal/` 五枚命中全是手写夹具，agent 名下零枚逐字比对）。
5. **本腿不补 `os.Stat` 那一支**：`Prepare` 里桩文本是在**同一次调用写下那枚文件之后**才拼的，三条失败路（建目录／写／改名）全在那之前 `return` ⇒ 存在性在这一腿**由构造保证**，加一发 stat＝恒真尺（派单明禁）。AC#2 的形状 (b)（路径不存在／是枚目录）活在 `task.go` 那一腿，因为那枚路径来自**别人填的**名册记录。

## 3. 判据两向读数（⚠ 待填：种回静默那一形必红／补回从句必绿，逐字贴）

## 4. 门禁读数（⚠ 待填：build/vet/d22scan/test 终态三数＋名册 comm 对比）

## 5. 判不动的地方（四格，全部具名）

1. **生产接线（两枚 spill 构造点都在 `cmd/wisp`，本程禁改）**：`cmd/wisp/run.go` 现量造 `agent.NewSpiller(...)` 与 loop 内 `NewSpiller(opt.Config.ArtifactsDir, b)`（`internal/agent/loop.go:238`）。我把缝留成**各一行**：外部构造的走 `WithPointerJudge(...)`，loop 的走 `Config.PointerJudge`。⇒ **不接线时本腿的产码在生产上说的是"无法核实"而不是假话**（fail-closed，形状与 `task.go` 未接线那一支同款），但**它不说"读不到"**——那句要等 cmd 那一行。这与 AC#2b 的处境同形，去处也同类，**我不替它落 cmd**。
2. **`err.Error()` 无界那一格不属本程**：`internal/tools/task.go` 那支要不要改成不转述上游文本＝`Q-63` 边界（那两句受模板冻结钉系着），⛔ 我不判、也不动 `internal/tools` 一字。
3. **"批一张 L2 卡"这条路的精确适用条件**：票面 `:57` 第 4 条现量＝`[fs] allowed_dirs` 是**判级输入不是执行硬边界**（`fs.go:93-105` 只 `Canonicalize` 不查根），"会被拒"说的前提是**没人批**。我这一支的措辞按这个事实写（"落到 L2 要人批；没人批就是拒绝"），⛔ 不写成"根是硬边界"（那是 `Q-60` 的另一支，要 owner 说话）。
4. **AC 框**：`AC#2c` 与其余各框一枚不勾（勾它＝另一枚程）。本格交件套路＝"ⓐ 已落 + 两向读数 + 门禁"，**翻勾与否归非实现者**。

## 6. 推翻清单（编排者那几句现量的复算结果）

| 断言 | 结果 | 凭据 |
|---|---|---|
| `task.go:562` 那行字面（多一枚 `%s` 槽＝诚实从句） | **成立** | §1 表第 1 行，`grep -n` 逐字 |
| `spill.go:141` 那行字面（没有从句） | **成立** | §1 表第 2 行，`grep -n` 逐字 |
| "多一枚 `%s` 槽"＝诚实从句 | **成立** | `task.go:548` 逐字 `notice := t.d.pointerNotice(rec.ArtifactPath)`，该实参落进那个槽；`spill.go` 无对应物 |
| `spill.go:18` 一带自陈是 stub | **推翻行号** | `:18` 逐字是 `// D15(3) long-output spill (SPEC-05 §4.2): a single tool result over the`；"stub" 在 `:48`／`:51`。实质（这一支是 D15(3) 长输出桩）成立 |
| 票 177 已裁完这一条（AC#2c 的排程前提） | **成立** | §2.1 三条凭据（`Q-61`甲／`Q-63`甲已批、甲形已落进 `task.go:558`＋`internal/risk/shape_a_exemption_test.go`、落点＝逐次调用 box 的窗口排除） |
| "同一句'全文见'在两条腿上" | **成立** | 产码名下整仓只有那两枚命中，其余全在测试夹具/证据件 |

---

## 编排者补记（2026-10-03 15:5x，代跑不代判；§3/§4 两节「⚠ 待填」保持原样，我不代填）

本腿死于宿主"每日 Chat 额度用尽"（与 166-a2／255-r4／261-a2 同一波），死前已落三枚 commit：
`fad39b2c`（骨架＋现量台件）→ `5b14afe5`（产码，形ⓐ）→ `081019c0`（341 行测试）。以下三发读数全部以
**"编排者代跑"** 名义入账，判留给验收腿自己取数复认：

1. **包级绿尺**：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -run "Ticket174|Spill|PointerNotice" ./internal/agent/`
   ⇒ 21 枚 PASS／0 FAIL（其中 **7 枚是本腿新增的 `*174` 尺**：`TestSpillPointerUnwiredJudgeFailsClosed174`／
   `TestSpillPointerOutsideAllowlistNamesBothRoadsBack174`／`TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174`／
   `TestSpillHealthyPointerStillMatchesThePreFixSentence174`／`TestSpillPointerNoticeAddsNoSecondPath174`／
   `TestSpillPointerNoticeKeepsHeadTailAndBudget174`／`TestLoopSpillNoticeReachesHistoryUnwired174`）。
2. **变异 M1（种回静默形）**：把 `pointerNotice` 开头插 `if true { return "" }` ⇒ **5 枚具名用例当场红**
   （①④⑤⑦与 `TestSpillPointerCanonicalizeFailureQuotesNoUpstreamText174`）；还原用
   `git cat-file blob HEAD:internal/agent/spill.go >`，md5 还原前后一致（`f31147b4…`）。
3. **死腿残面**：它在跟踪树里留了一枚未还原的突变（`MUT-174R2-M1`，即上式）＝编排者已还原、`git diff` 归零。
   另有它自己的三份 commit 信息草稿与 `mut-m1-silent-shape.txt` 读数随本补记同批入库。**AC 框一枚未碰。**
