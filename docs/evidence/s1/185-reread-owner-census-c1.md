# 185-c1 — 只读普查：「第二次读回该由谁给这枚声明」三支候选的现量与代价

- 派单：`.scratch/wisp/dispatches/2026-09-28-142x-readonly-185-c1-who-should-declare-the-read-back-path-and-what-does-each-route-cost.md`
- 工单：`.scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-because-the-fs-read-result-becomes-an-undeclared-c25-mark.md`
- 性质：**只读普查**。本表**不勾任何 AC 框**、**不派落地腿**、**不替编排者拍板选哪一支**。
- 本表所有行号皆可漂：凡引号内是**符号名＋`grep -n` 尺**，不是死号。
- ⚠ 三行定性（本票全程只在**本仓自己的判据与台件**里出现这些形状，**没有任何本机被入侵的证据**；最坏后果的形状是**产品正确性**的两张相反错——该拦的没拦／不该拦的被拦——不是攻击面）：出处＝票 185 现量第 5 条，本程未新增任何运行时外部面。

---

## ① 起手锚 / 写面闸门 / 三枚受保护文字 md5（起手向）

| 项 | 读数 |
|---|---|
| 起手 `date` | `Mon Sep 28 14:24:35 CST 2026` |
| 起手锚 `git log -1 --format='%h %cd'` | `72c76d42 Mon Sep 28 14:23:56 2026 +0800`（分支 `dev`） |
| **本票的"未修码"＝`72c76d42`（起手那一枚 HEAD，具名）** | 票 183 的修法已在 `0662a35a` 落地并在这枚 HEAD 上**在场**（现量：`grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/taintmatch.go \| grep -v _test.go` 命中 `attachDeclaredPath`（定义＋调用）与 `spellsDeclaredPath`（定义＋调用）；`provenance.go` 里 `declaredNorm`／`MarkWithHostPath` 的 183 更正段在场）⇒ **HEAD＝"183 修完之后、185 未修"**。本程凡说"今天红"，基线一律＝`72c76d42`，**不把 HEAD 与未修码当同义词**。 |
| 派单写的锚 `78ea1d19` | 现量：`git merge-base --is-ancestor 78ea1d19 HEAD` ⇒ **是 HEAD 的祖先**；`78ea1d19..HEAD` 之间另有 `72c76d42`（183 结案）等 commit。本程用**自己起手的 `72c76d42`** 作基线。 |
| HEAD 在程中前移（共享工作树，别家票的活） | 终量：`git log --oneline 78ea1d19..HEAD` 顶部＝`e9ef94d0`（票 181-c1 的 docs 骨架），其下 `03c01d71`（docs）、`c1c96008`（**本程骨架**）、`72c76d42`。⇒ 我看到的产码面未被这些 commit 改动（`git diff --name-only 72c76d42 HEAD` 见 §⑨ 终量，逐枚是 docs/`.scratch`）。 |
| 起手写面闸门 `git status --porcelain -- internal/ cmd/` | **空**（终态复量见 §⑨，仍为空） |
| 受保护文字 1／`internal/risk/provenance.go` 的 `Mark` 契约块（468-473 那一段） | md5 `858e45116383caa3e7c1dd4b0924fad1` ＝ 基线 ✓（尺：`sed -n '468,473p' internal/risk/provenance.go \| md5sum`） |
| 受保护文字 2／`internal/risk/taintmatch.go:11-15` | md5 `5680ddd18e2d2ec2a85e485b54f4c12e` ＝ 基线 ✓ |
| 受保护文字 3／`provenance.go` 的 `MarkWithHostPath` 文档块原 25 行（482-506） | md5 `f89e891e5eee3f3ea2b4f89d921072c4` ＝ 基线 ✓ |
| ⚠ `taintmatch.go:100-109` 一带 | **不在受保护名单**（"仅参考"格，随号推移）；本程**不据它报警**，全程只读。 |

**根因（本程自己从码上取的，不是转述）**：`Bridge.mark`（符号尺 `grep -n "func (b \*Bridge) mark" -A 22 internal/tools/bridge.go`）把每次执行后的 `res.Text` 交给 `MarkWithHostPath(..., hostPath)`，而 `hostPath` 取自 per-call 的 `hostPathBox.get()`；全仓今天**只有一枚 setter**＝`internal/tools/task.go` 的 `task.output` 指针分支（尺 `grep -n "hostPathBoxFromCtx" internal/tools/*.go \| grep -v _test.go` → 只有 task.go 一处调用 `box.set(...)`）。⇒ `fs.read` 成功读回的那枚 mark **没有声明** ⇒ 它正文里能拼出那条路径的 ≥8-rune 窗口照旧进索引 ⇒ 同一发路径的**第二次**续读被自己上一步的正文咬住（票 183 的 twin 机制：正文填充行 `-output-` 与产物文件名 `tool-output-agent-task-…` 同款拼写）。

---

## ② AC#2(a) 桥侧：`fs.read` 成功时也填 `hostPathBox`

**现量取号（本程 14:2x 现跑，号以这次为准）**：`grep -n "hostPathBox\|withHostPathBox" internal/tools/bridge.go internal/tools/task.go`

| 命中 | 内容 |
|---|---|
| `bridge.go:459` | `ectx, hostPaths := withHostPathBox(ectx)`——每发调用一枚新 box |
| `bridge.go:535` | `b.mark(dec, res, hostPaths.get())`——box 只在盖戳这一处被读 |
| `bridge.go:581-624` | `hostPathBox` 类型本体（`set`／`get`／`hostPathBoxFromCtx`；`set` 对空串与 nil box 静默返回） |
| `bridge.go:566-579` | `func (b *Bridge) mark(...)`：闸门是 `risk.IsSensitiveSource(dec.Tool)`，非名册工具**直接不盖戳** |
| `task.go:253-255` | **全仓唯一 setter**：`if box := hostPathBoxFromCtx(ctx); box != nil { box.set(rec.ArtifactPath) }`，值来自**名册记录**（宿主自己铸的路径） |

- **落点具名函数**：`fsRead.Execute`（`internal/tools/fs.go`，尺 `grep -n "func (t fsRead) Execute" internal/tools/fs.go`）。变异落点＝成功返回之前那句 `hostPathBoxFromCtx(ctx).set(canon)`；`canon` 来自 `t.d.open(a.Path)`（C26 规范化**模型给的 `path` 参数**；尺 `grep -n "canon, err := t.d.open" internal/tools/fs.go`）。
- **该函数今天有没有别的调用方**：无。`fsRead` 只由 `BuiltinFSEntries`（`grep -n "BuiltinFSEntries" internal/tools/fs.go`）注册、由桥的 `entry.Tool.Execute(ectx, ...)`（`grep -n "entry.Tool.Execute" internal/tools/bridge.go`）唯一调用 ⇒ 填 box 这条腿**只在过桥时可达**（不经桥的调用方拿到 nil box，声明自然缺席）。
- **变异现量（代价取数，副本在 `.scratch/wisp/probes/185/c1/mut/fs.go`，走 `-overlay`，跟踪件零改动）**：
  - 尺：`go test -count=1 -v -overlay=.scratch/wisp/probes/185/c1/overlay-muta.json ./internal/tools/`
  - 读数：`rc=0`、`ok github.com/CarlosShao/wisp/internal/tools 15.521s`（留档 `probes/185/c1/logs/tools-on-muta.txt`）
  - **会变红的既有判据＝逐枚点名：零枚。**（不是"枚数=0"的偷懒说法：点名射程＝`internal/tools` 全包 `-v` 全量跑，`--- FAIL` 行数为 0。）
  - **仪器自证（防"恒绿"）**：同一张 overlay 名册换成故意写坏的副本 `mut/fs_broken.go` → `go build -overlay=... ./internal/tools/` **rc=1**（`logs/overlay-proof-broken.txt`：`syntax error: imports must appear before other declarations`，报的就是我这枚副本），真实变异副本 rc=0（`logs/overlay-proof-muta.txt`）⇒ 映射确实生效，"零枚红"是**读数不是没跑**。
- **最贵的那枚代价**：它把"豁免的根据"从**宿主自己写下的那一串**换成**模型这次递进来的参数值**。票 177 W-2 那条承重边界（`外来正文携带同一条路径 ⇒ 仍命中 R4`）在 `fs.read` 这一族上就此失效——凡是模型递进来的路径，它读回来的正文里拼出该路径 ≥8 rune 的那段外来文字**不再算证据**（正是 `TestPointer183WorstCaseOfTheLandedExemptionIsPinned` 钉住的最坏形状，只是范围从"一枚宿主桩"扩到"每一次成功读取"）。而今天的判据面**一枚都抓不到**（上面 rc=0），＝零牙。
- **票 177 三条已裁边界碰到哪条**：
  - **碰到＝参数侧按值放行**。`provenance.go` 的 `MarkWithHostPath` 文档块里那句边界（尺 `grep -n "not a by-value pass on the parameter side" internal/risk/provenance.go`）明写"边界 (2)＝不是按参数值放行（票 183 AC#5 ①）——无名册、无 scope 级路径表、**不拿参数去和名字表比**"。(a) 的声明值**恰好就是参数**（`canon`＝C26 规范化后的 `path`），**本程判它等价于那一形 ⇒ 具名停手上报，不放行、不选它**。
  - 没碰 per-scope 名册（box 仍随发生死）、没碰 `allowlist`（全程不读 `InAllowlist`）。
  - 附带碰：**票 177 W-2**（同一枚 mark 正文里的外来同路径拼写要仍命中）——见上一格。

---

## ③ AC#2(b) D15 再落盘层：续读不产生第二枚 mark

- **那一层今天在哪**：`internal/agent/spill.go`（尺 `grep -n "func " internal/agent/spill.go`）：`NewSpiller:43` / `CapRaw:68` / `Prepare:84` / `nextSequence:150` / `artifactName:184` / `writeFileExclusive:249`。落盘桩文案在 `spill.go:141`（同款"全文见 %s"）。**唯一调用点**＝`internal/agent/loop.go:702` 的 `l.sp.Prepare(c.ID, log.Text)`（尺 `grep -n "sp.Prepare" internal/agent/loop.go`）。
- **⚠ 一条决定性的现量**：今天**根本不存在**"宿主替模型读"的通道。尺：`grep -rn "os.ReadFile\|os.Open" internal/agent/*.go \| grep -v _test.go` → **空**。`task.output` 的续读走的是进程内名册记录（`grep -n "rec.Text" internal/tools/task.go`），**不读文件**；文件那一侧唯一的读通道只有 `fs.read`（过桥 ⇒ 必盖戳）。⇒ **(b) 不是"关掉一个开关"，是要新铺一条数据入口**。
- **"不产生第二枚 mark"在码上等价于改哪一处**——两形，代价不同：
  - **(b-i) 新增宿主侧读回通道**（环路替模型把产物读回来、不进桥的 `mark()` 腿）。落点具名：`internal/agent/loop.go` 的 spill 调用点附近＋一枚不经 `Bridge.mark` 的入口。⇒ 这条新入口的回执**没有来源戳**，正面撞**票 175 的承重声明**（禁令面）。会红的既有判据（逐枚点名）：`internal/tools/ticket175r2_stamp_live_test.go` 的 `TestEveryRegisteredToolIsClassifiedForMarking175r2`（现量名册尺 `grep -n "^func Test" internal/tools/ticket175r2_stamp_live_test.go`＝344 行那枚）。risk 包那 8＋8 枚（`pointer_183_test.go`／`shape_a_exemption_test.go`）**不红**——它们直接调 risk API，不经桥。
  - **(b-ii) 把 `fs.read` 从盖戳家族摘掉**（改 `bridge.go` 的 `IsSensitiveSource` 闸门那一条腿）。会红的既有判据（逐枚点名，名册现跑尺 `grep -n "^func Test" internal/tools/fs_test.go internal/tools/ticket175r2_stamp_live_test.go`）：`internal/tools/fs_test.go` 的 `TestFSReadReturnsTheFileAndTaintsIt`、`TestFSReadTaintFeedsR4`，加 `ticket175r2_stamp_live_test.go` 的 `TestEveryRegisteredToolIsClassifiedForMarking175r2`。
  - ⚠ 这两形本程**都没做变异现量**（列在 §⑧3），上面是"落点＋点名预判"，不是读数。
- **会不会把票 175 的成功结果来源标记拆掉**：**会**。(b-ii) 直接拆；(b-i) 不拆 `fs.read` 但新增一条"带正文进上下文却没有戳"的路，正好落在 `TestEveryRegisteredToolIsClassifiedForMarking175r2` 与 `TestNonContentResultStaysUnmarked175r2` 之间那条线上（后者钉的是"非正文回执可以没戳"，新通道带的是正文）。⇒ **这一条是禁令面，具名报给编排者，不当成可选。**
- **最贵的那枚代价**：要么**新铺一条不受桥管的数据入口**（新的盖章责任真空＋新的 D34 表行），要么**摘掉一枚在册工具的来源戳**；两者都不修当前这一发的根因——第一次读回的那枚 mark 仍在，堵的是"内容巧合"而不是"通道"。
- **票 177 三条边界碰到哪条**：**没碰** per-scope 名册、没碰参数侧按值放行、没碰 `allowlist`。（碰的是**票 175**那条，见上。）

---

## ④ AC#2(c) 不修，把"同一条路径只许读一次"写成明说文案

- **文案落点具名文件**：唯一非冻结的宿主文案位＝`internal/tools/task.go` 的 `pointerNotice`（尺 `grep -n "func (d TaskDeps) pointerNotice" internal/tools/task.go`）。它今天已经在写"这条路径现在读不到"那一族真话（票 174 AC#2 的形状：不在授权范围／连规范化都没过／副本文件不存在／不是常规文件），且**是拼进桩正文里给模型读的那一句**，不是日志。
- **冻结判定**（现量）：尺 `grep -n "全文见" internal/tools/task.go docs/PLAN.md` → **只有 `task.go:257` 命中，`docs/PLAN.md` 里没有"全文见"这个字面**。同款字面另在 `internal/agent/spill.go:141`（D15 那一层的桩）。`PLAN.md:2564` 冻的是**那一行的定案形状**（尺 `sed -n '2564p' docs/PLAN.md`，现量逐字：`task.output` 行写"截断按 D15「单个工具结果」那一行的既有规矩（头 500 token＋尾 200 token＋总长＋文件路径），**只截不指＝不合格**（必须给可续读的路径）"）。
  ⇒ 结论：**文案落 `pointerNotice`＝不动任何冻结件**；**文案若要落桩格式串本体（`task.go:257`／`spill.go:141`）＝碰 `PLAN.md:2564` 那行的定案形状 ⇒ 需人工批准**（具名报回，本程一字未改）。
- **会红的既有判据**：不改行为、只加一句 note ⇒ **零枚行为判据转红**。但 note 是**逐条比对**的文案面：`internal/tools/task_output_pointer_notice_test.go` 那一族（该文件 `:331` 现量断 `"全文见 "` 在桩里）按 note 文本比——本程**没跑** `grep -n "^func Test" internal/tools/task_output_pointer_notice_test.go` 去逐枚点名 ⇒ 列在 §⑧4。
- **最贵的那枚代价（本程认为这一支真正的贵处，具名）**：这句话**今天不真**。第二次读回被不被拒，取决于第一次读回的**正文有没有恰好拼出那条路径的 ≥8-rune 窗口**（内容巧合），不是一条规则。票 183 那把频率尺的现量是**10 份正文里 3 份**会响（且语料＝仓内文件，真实文档/LLM 输出语料＝零读数，票 185 现量第 4 条自己具名）。⇒ 写"只许读一次"会把一个**内容依赖的偶然行为**说成**确定的产品承诺**，同时与 `PLAN.md:2564` 那行的**意图**（"可续读"）相撞（字面不撞、意图撞）。
- **要不要人工批准**：**要摆**（只报是哪一格、为什么，不替编排者拍板）——见 §⑫ 人工批准项 ②。
- **票 177 三条边界碰到哪条**：都不碰（零行为变更、零放宽）。

---

## ⑤ 有没有第四支（只作用于这一发、不放宽任何一类外来正文）

**码上确实有这条形，但它的判据本体碰已裁边界 ⇒ 具名上报，不当成免费路。**

- 现量：豁免只能**随一枚 mark** 存在——`taintmatch.go` 的 `attachDeclaredPath` 把声明挂在 `fragmentIndex` 上（尺 `grep -n "func (f \*fragmentIndex) attachDeclaredPath" internal/risk/taintmatch.go`），`provenance.go` 的 `MarkWithHostPath` 文档块明写"**The located span lives nowhere but this call** … Nothing per-scope is kept, because a scope-level path list would be a new roster (177-c1 §2, A353)"（尺 `grep -n "scope-level path list would be a new roster" internal/risk/provenance.go`），`taintmatch.go` 的 `newFragmentIndex` 文档同样写着"a per-scope path roster is explicitly NOT approved"。
- ⇒ 于是"只作用于同一条产物路径的续读"只有两种判据形状：
  1. **仍随 mark** ⇒ 就得让 `fs.read` 那枚 mark 带声明 ⇒ 回到 §②(a) 的参数侧；
  2. **按"这条路径是不是本 scope 里宿主自己落盘的那一枚"判** ⇒ 名册本体就是 **per-scope 名册＝票 177 已裁（未批）那一形**。宿主铸的路径今天**是有登记的**（`internal/tools/task_backfill.go:112` `rec.ArtifactPath = sp.Path`；尺 `grep -n "ArtifactPath" internal/tools/task_backfill.go`），所以这条路**在码上可写**，但一写就碰 177 的边界。
- **判据形状（若要它成立）**：`宿主铸过 ∧ 本 scope ∧ 只此一条 ∧ 随 scope 生死`，且必须**不**把该路径的窗口从**别的 mark** 的证据里去掉（否则又是 W-2 那面）。
- ⚠ 本程**没有**实现或测量这一支（无变异、无读数），只给"存在／碰哪条边界"的具名码量。⇒ 列 §⑧5。

**另有一形值得报（不是第四支，是它的反面）**：`fs.list` 的回执正文**必带宿主全路径**（`joinForListing(dir, name)` 拼全路径；尺 `grep -n "func joinForListing" internal/tools/fs.go`），但它**不在** `sensitiveSourceTools` 名册里 ⇒ 桥 `mark()` 的 `IsSensitiveSource` 闸门直接 return、**不盖戳** ⇒ 它今天**不造阻断者**，代价是"正文带宿主路径却无来源戳"。见 §⑥ 名册。

---

## ⑥ AC#3 覆盖面名册（逐枚带尺；名册与枚数一律现跑，未从票面抄）

**主尺（派单给的那把，现跑）**：`grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/ | grep -v _test.go`
命中名册（**9 个非测试位置**，逐枚）：

| # | 命中 | 是什么 |
|---|---|---|
| 1 | `internal/risk/provenance.go:479` | `Mark` 委派 `MarkWithHostPath(..., "")`（四参契约不变） |
| 2 | `provenance.go:482` | `MarkWithHostPath` 文档块（**受保护文字 3**） |
| 3 | `provenance.go:540` | `MarkWithHostPath` 函数本体 |
| 4 | `internal/risk/taintmatch.go:91` | 注释引用（183 段） |
| 5 | `taintmatch.go:100` | 注释引用（`declaredPath` 字段） |
| 6 | `taintmatch.go:133` | 注释引用（ordering (i)） |
| 7 | `internal/tools/bridge.go:575` | 全仓**唯一**调用 `MarkWithHostPath` 的地方 |
| 8 | `bridge.go:581-624` | `hostPathBox` 载具本体 |
| 9 | `internal/tools/task.go:253` | 全仓**唯一** `box.set` ＝ `task.output` |

⇒ **今天带声明的只有 1 枚工具（`task.output`）；其余全部不带。**

**谁会被盖戳**（名册尺 `grep -n -A 6 "var sensitiveSourceTools" internal/risk/provenance.go`＝`provenance.go:99-103`，9 枚：`fs.read`／`search.content`／`clipboard.read`／`system.get`／`web.fetch`／`doc.read`／`screen.capture`／`asr.transcript`／`task.output`）。

**今天真正存在的工具回执**（名册尺 `grep -rn "Name() string" internal/tools/*.go | grep -v _test.go` ＋ `grep -n "func (fsList) Name" internal/tools/fs.go`）：`fs.read`／`fs.list`／`fs.write`／`fs.edit`／`fs.trash`／`fs.move`／`fs.delete`／`task.output`。`clip.*`／`search.content`／`doc.read`／`screen.capture`／`transcript`／`system.get` 在 `internal/tools` **没有实现本体** ⇒ 今天不产生回执（D34 表里排 S3/S4），但**已在盖戳名册里**。

逐枚（工具 → 回执带不带宿主路径 → 会不会造同款阻断者 → 今天带不带声明）：

| 工具 | 回执带宿主路径？ | 造同款阻断者？ | 带声明？ | 本程尺 |
|---|---|---|---|---|
| `task.output` | **带**（桩里 `全文见 <path>`，值来自名册记录） | 不造（自己声明了自己那条）；但**只覆盖它那一枚路径** | **带**（唯一一枚） | `grep -n "box.set" internal/tools/task.go` |
| `fs.read` | 正文＝文件内容；路径**只在 `Origin` 字段**里（`Origin` 不进索引） | **会**——正文里拼出该路径 ≥8-rune 窗口就进索引（票 183 的 `-output-` twin；真机三行读数见 §⑧1 的具名出处） | **不带** | `grep -n "func (t fsRead) Execute" -A 45 internal/tools/fs.go`（返回 `Result{Text: string(buf), Origin: canon}`） |
| `fs.list` | **必带**（`joinForListing` 拼全路径） | **不造**——`fs.list` 不在盖戳名册 ⇒ 桥根本不盖戳（正文带路径却无戳＝另一面） | 不适用（无戳） | `grep -n "func joinForListing" internal/tools/fs.go`、`grep -n "IsSensitiveSource" internal/tools/bridge.go` |
| `fs.write`／`fs.edit` | 回执 `res.Origin = target`（尺 `grep -n "res.Origin" internal/tools/fs_write.go internal/tools/fs_edit.go`）；两者**都不在**盖戳名册 | 不造（无戳） | 不适用 | 同上 |
| 环路落盘桩（`spill.go:141`） | **带**（同款 `全文见 <path>`），进的是上下文不是桥 | **不造**——这层在桥**之后**，桩正文从不进 `MarkWithHostPath` ⇒ 该路径**在任何索引里都不存在** | 无声明（也无戳） | `grep -n "全文见" internal/agent/spill.go`、`grep -n "sp.Prepare" internal/agent/loop.go` |
| `clip.*`／`search.content`／`doc.read`／`screen.capture`／`transcript`／`system.get` | 未实现 ⇒ 零回执 | 未实现（**零读数**，不是"验过不会"） | 在册名册、无 setter | 尺见上一段 grep |

⇒ **同款阻断者今天的成因只有一枚**：`fs.read` 的正文。`fs.list` 那一族是"带路径但无戳"的**相反面**（少一张证据，不是多一枚阻断者），本票不裁它，具名给编排者。

---

## ⑦ AC#7「同名不同目录」那一形的现状（**验了**，具名两枚读数）

- **比对面现量**：尺 `grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/*.go | grep -v _test.go` →
  `taintmatch.go` 的 `attachDeclaredPath`（存 `[]rune(hostPathNorm)`，**整条**归一化路径）、`contains` 里那句 `spellsDeclaredPath(f.declaredPath, win)`、`spellsDeclaredPath` 本体（逐 rune 连续子串，**不查哈希**）；`provenance.go` 里 `m.idx.attachDeclaredPath(declaredNorm)`，`declaredNorm = hp`（＝`normalizeTaint(hostPath)` 整条）。⇒ 比对面是**整条归一化串**，**不是 basename**。
- **读数 1（未变异源＋我这枚探针，`-overlay`）**：
  `go test -count=1 -v -overlay=.scratch/wisp/probes/185/c1/overlay-ac7.json -run 'TestPointer185C1' ./internal/risk/` → **rc=0**；`TestPointer185C1SameNameDifferentDirectoryStillHits` **PASS**、`TestPointer185C1DeclaredPathItselfStaysClean` **PASS**（`probes/185/c1/logs/ac7-on-head.txt`）。⇒ **今天这一形不被放宽**（豁免吃不到 basename 级），与上面"整条串比对"一致。
- **读数 2（同一枚探针 ＋ basename 放宽变异）**：变异＝`.scratch/wisp/probes/185/c1/mut/taintmatch.go`（在 `contains` 顶部加"候选的 basename 出现在声明路径里 ⇒ 整支候选不算证据"＋`basenameTail` helper），名册 `overlay-ac7-mutant.json`：
  - `-run 'TestPointer185C1'` → **rc=1**，`--- FAIL: TestPointer185C1SameNameDifferentDirectoryStillHits`（`pointer_185_c1_ac7_probe_test.go:71`），控制腿仍 PASS（`logs/ac7-on-mutant.txt`）
  - **整包 risk 全量**（同 overlay，无 `-run`）→ 只有我这枚探针红，**已入库判据一枚都没红**（`logs/risk-suite-on-mutant.txt`；射程含 `pointer_183_test.go` 8 枚＋`shape_a_exemption_test.go` 8 枚＋其余）
  - 变异确实生效的证据：同一张名册里 `taintmatch.go` 换成我这枚副本后**读数的方向变了**（同一枚探针 PASS→FAIL），且该副本在 `overlay-ac7.json`（不含变异）下 PASS ⇒ 变量只有 `taintmatch.go`。
- ⇒ **本程的具名结论**：**"同名不同目录"那一形今天在已入库判据面上＝零牙**（basename 级放宽无人能看见）。它的牙必须由落地腿**新写**那一枚常驻判据；探针副本（形状、fixture、两腿配对）已落在 `.scratch/wisp/probes/185/c1/pointer_185_c1_ac7_probe_test.go`，可直接搬。落点建议具名：`internal/risk/pointer_183_test.go` 的第 9 枚（现量 `grep -c "^func Test" internal/risk/pointer_183_test.go`＝**8**，未抄票面）或 `internal/tools/pointer_183_cli_seam_test.go` 同批加一枚（现量名册：`TestPointer183BackfilledArtifactRereadStaysClean`／`TestPointer183SiblingArtifactInTheSameDirectoryStillHits`）。
- ⚠ 本程**没动跟踪件**（探针与变异只在 `-overlay` 名册里；`git status --porcelain -- internal/ cmd/` 见 §⑨ 仍为空）。

---

## ⑧ 本程没测什么（逐名）

1. **真机 CLI 那一发 `TestRereadHostPointerOnCLISeam183a1` 本程没重跑**。具名理由：(i) 它的写出 sink 会把 `.scratch/wisp/probes/183/r1/logs/e2e-readings.txt` **原地重写**，而那枚文件被票 185 逐行引用（第 55／60／64 行），我不覆盖别家已入库读数；(ii) 该腿**不在跟踪件里**——现量 `ls cmd/wisp/*_test.go` **没有** `zz183r1_e2e_test.go`，它只以 `.scratch/wisp/probes/183/r1/zz183r1_e2e_test.go` 的 `-overlay` 副本存在（名册尺 `cat .scratch/wisp/probes/183/r1/overlay-cli.json`）。HEAD 上"第一发绿／第二发红"那三行读数出处＝编排者 14:1x 的票面 Progress log ＋ `.scratch/wisp/probes/183/orch/logs/cli-with-prefix.txt`（本程引用、未复算）。⇒ **落地腿如要现量，先把 sink 改到自己的探针目录。**
2. **不带 PATH 前缀那一形（`exit status 0xc0000135`，票 98 那枚加载期坑）本程未亲取**，同 §⑧1 具名出处（在册尺＝`scripts/wisp-cli-tests.sh:20` 的 PATH 前缀）。
3. **(b) 的两支落法（b-i 新通道／b-ii 摘 `fs.read` 的戳）都没做变异现量**：§③ 是"落点＋点名预判"，红否未取数。
4. **(c) 文案落 `pointerNotice` 的既有判据名册未逐枚点名**：没跑 `grep -n "^func Test" internal/tools/task_output_pointer_notice_test.go`。
5. **第四支（§⑤）没有实现、没有变异、没有读数**：只给"这条形码上存在／一写就碰 177 per-scope 名册"的具名码量。
6. **票 185 AC#4 那把频率尺（3/10）本程未复量**：语料是别家探针的仓内文件，真实用户文档与 LLM 输出语料＝零读数（票面自己具名）；本程没造新语料。
7. **`clip.*`／`search.content`／`doc.read`／`screen.capture`／`transcript`／`system.get` 的回执形状**：只跑了"注册名册里有没有实现本体"一把尺（§⑥），未逐枚读形状（今天不存在 ⇒ 无从读）。
8. **`cmd/wisp` 的跟踪腿本程未跑**（含带/不带 PATH 前缀两形，同 §⑧1/2）；**全仓 `go test ./...` 未跑**（A359 并发假红；只读程按包单跑）。
9. **恒真判据的自证**：§② 的"零枚红"配了故意写坏的 overlay 证明腿（build rc=1）；§⑦ 的 PASS→FAIL 方向变化即仪器有牙的证明。**没有**拿"今天会绿的守卫"充当新功能有牙——`TestShapeAW3`／`TestPointer183AlmostPath…` 在 basename 变异下仍绿＝本程据此判它们**管不着**这一形（见 §⑦ 读数 2）。

---

## ⑨ 门禁终态

> ⚠ 按 `A363` 薄规矩：只读程的门禁读数**只作现状记录，不充当任何 AC 的结案凭据**。

| 尺 | 读数 |
|---|---|
| `sh scripts/d22scan.sh` | **rc=0**；`ban #8 internal/` **examined=433**＝基线（不多不少，无需解释漂移）；附 `ban #8 cmd/`=45、`ban #8 frontend/`=85、`ban #8 design/`=39、`ban #7 internal/tools/`=21、`bans #1-5 internal/`=208、`#1-5 cmd/`=23（留档 `probes/185/c1/logs/d22scan.txt`） |
| `bash .scratch/wisp/probes/154/gate-clauses.sh`（**比红腿名册、不比退码**） | 名册 14 腿：13 枚 `ok`，**BAD＝1 枚＝`G6neg`（声明 ring／基线 1 枚／实测 3 枚／因＝新增未成对）**；`G5neg` ring 8=8、`G7neg` ring 4=4、`G2` 2=2、`G6` 1=1 均 ok；基线过期枚数＝0；腿数断言相符（缺腿 0／空头声明 0）；聚合退码＝1。**与在册名册（今天唯一在册红腿＝`G6neg`＝票 178）差集＝空 ⇒ 不需要停手**（留档 `logs/gate-clauses.txt`；**没跑** `probes/161/r6/flip-declaration.sh`＝在册禁令） |
| `go test -count=1 ./internal/risk/`（**单跑**，A359） | `ok … 5.251s` |
| `go test -count=1 ./internal/tools/` | `ok … 16.050s` |
| `go test -count=1 ./internal/agent/` | `ok … 1.995s` |
| CLI 那一面（`-overlay`＋在册 PATH 前缀；两形） | **本程未跑**，具名理由与出处见 §⑧1／§⑧2 |
| 写面闸门（终态复量） | `git status --porcelain -- internal/ cmd/` ＝ **空**（与起手同为空；复量时刻在最后一枚 commit 之后，读数见交件消息） |
| 产码面是否被别家 commit 动过 | `git diff --name-only 72c76d42 HEAD` 交件时报，逐枚应为 `docs/**` 或 `.scratch/**`（本程自己的 commit 只含 `docs/evidence/s1/185-reread-owner-census-c1.md` 与 `.scratch/wisp/probes/185/c1/**` 与票 185 的 Progress log 追加） |

---

## ⑩ 被拒／没成功的调用（逐条）

- **零枚**：本程没有任何一次工具调用被权限系统拒绝，也没有一次意外的非零退出。
- 唯一的"非零 rc"是**设计如此**的两枚：§② 的 overlay 生效证明（故意写坏的 `mut/fs_broken.go` → `go build` rc=1）与 §⑦ 读数 2（basename 变异 → `go test` rc=1）。它们不是失败，是判据在响。

---

## ⑪ 有没有跑过删除命令 + 工具调用终值自报

- **删除命令：一枚都没跑。**全程无 `rm`／`mv`／`unlink`／`git clean`／`git restore`／`git checkout .`，无 `--amend`／`reset`／`rebase`／`stash`／`worktree`／`switch`，**没有 push**。临时件只建不删：`probes/185/c1/**` 下的 `mut/fs_broken.go`（证明件）、`mut/fs.go`、`mut/taintmatch.go`、三枚 overlay 名册、`logs/*.txt` 全部保留。
- **禁跑清单**：`probes/161/r6/flip-declaration.sh` **没跑过**。
- **工具调用终值**：见交件消息（硬顶 40；新探索在 28 枚之前停止，之后只用于填表、跑门禁、commit）。

---

## ⑫ next＝落地腿派之前还缺什么（含要不要摆人工批准项）

**还缺的读数（本程未取，逐名）**：

1. (b) 两支落法的红否现量（§⑧3）——尤其 `TestEveryRegisteredToolIsClassForMarking175r2` 那枚承重腿实际怎么响（名字为 `TestEveryRegisteredToolIsClassifiedForMarking175r2`）。
2. (c) 落 `pointerNotice` 时的既有文案判据名册（§⑧4）。
3. AC#1 那枚成对常驻判据的形状（本程只在 `probes/185/c1/pointer_185_c1_ac7_probe_test.go` 交了 **AC#7** 那一形的配对探针）。
4. 真机 CLI 两形（§⑧1/2）：如要现量，sink 必须落到落地腿自己的目录，**不再覆盖 `probes/183/r1/logs/e2e-readings.txt`**。
5. AC#4 那把频率尺若要报"这次往前推让哪一类外来正文不再算证据"，需要**真实语料**（今天 3/10 是仓内文件语料）。

**要不要摆人工批准项（只报是哪一格、为什么；本程不拍板）**：

- **①（碰已裁边界，具名上报）**：**§②(a) 判它等价于"按参数值放行"**——声明值就是模型递进来的 `path`（C26 规范化后）。出处边界文字＝`provenance.go` 的 `MarkWithHostPath` 文档块边界 (2)（尺 `grep -n "not a by-value pass on the parameter side" internal/risk/provenance.go`）＋票 183 AC#5 ①。同批碰 **票 177 W-2**。**⇒ 这一支要落地就必须先有人工批准**（本程不放行、不选它）。
- **②**：**§④(c)**——文案若落桩格式串本体（`task.go:257`／`spill.go:141` 的"全文见"那句）＝碰 `PLAN.md:2564` 那行的定案形状 ⇒ 人工批准；落 `pointerNotice` 则不碰冻结件，但"只许读一次"这句话与 2564 那行的**意图**（可续读）相抵 ⇒ 属产品承诺层决定。
- **③**：**§⑤ 第四支**——它的判据本体是 **per-scope 宿主产物名册**＝票 177 已裁（177-c1 §2／A353 明写"未批"）那一形 ⇒ 要写就得先有人批名册。
- **④（禁令面，具名）**：**§③(b)** 任一落法都动票 175「成功结果必须有来源标记」的承重声明 ⇒ 报给编排者裁，本程不裁。
- **一字节未动**（复量见 §①／§⑨）：`allowlist.txt`、审批超时常量、SLO 阈值、golden、`docs/PLAN.md`、`docs/specs/**`、三枚受保护文字。
