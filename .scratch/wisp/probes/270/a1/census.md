# 270-a1 — 票 270 三形代价普查（AC#0 的料）

腿＝`270-a1`（**只读代价普查腿**）。射程＝票 270 **AC#0**：三形（甲／乙／丙）各自的改动面、顶红谁现有的钉、动不动文案、要不要新导出名，＋ panic→零读数那一形的现有兜底枚数。
⛔ 本腿**不选形**、⛔ 不翻任何 `- [ ]`（AC#0 由编排者凭本表翻）、⛔ 不改产码一字、⛔ 不写裁决表到 `docs/evidence/s1/`、⛔ 不动台账、⛔ **零 `go` 命令**（另有一枚腿在 `cmd/wisp` 跑突变，并发 go 进程会洗它的读数）。
§3 的判语格子一律留空＝给编排者拍的，本腿一个不填。

---

## §0 起手锚

| 项 | 命令（逐字） | 读数 |
|---|---|---|
| HEAD／分支 | `git rev-parse --short HEAD` ／ `git branch --show-current` | `2822ee04` ／ `dev` |
| 产码面工作树（起手） | `git status --porcelain -- cmd internal \| wc -l` | **0 行** |
| 全仓工作树（起手） | `git status --porcelain > logs/gitstatus-full.txt` | 原样存盘（⛔ 本腿不清他腿的脏，不动 `cmd/wisp` 那枚在飞的腿的面） |
| 终态对拉（交件时） | `git status --porcelain -- cmd internal \| wc -l` | **0 行＝起手**（票面 AC#4 里"起手＝终态"这半条本腿可自证；另一半"整包红名册作差"本腿不量，见 §4 第 1 条） |

### §0-补 本腿零 go 命令的自证（跑过的命令全文，逐字）

`git rev-parse --short HEAD` · `git branch --show-current` · `git status --porcelain -- cmd internal` · `git diff --numstat` · `git show -s --format=…` · `git add <点名 pathspec>` · `git commit -F <txt> -- <点名 pathspec>` · `mkdir -p .scratch/wisp/probes/270/a1/logs` · `ls` · `wc -l` · `wc -c` · `sed -n '…p'` · `cat` · `cat -n` · `head` · `tail` · `cut` · `sort` · `uniq -c` · `awk` · `grep -rn`（含 `-c`／`-o`／`-P`／`-B`／`-A`／`--include`）· `md5sum` · `for f in …; do [ -e … ]; done`

⇒ **清单里没有任何一条以 `go ` 开头，也没有任何一条以 `go.exe` 开头。**

落盘日志（大输出先落文件，⛔ 未把上百行原文读进对话）＝`.scratch/wisp/probes/270/a1/logs/`：
`baseoptions-all.txt`（194 行，全仓含 `.scratch` 副本）· `baseoptions-tracked.txt`（22 行，只 `cmd internal`）· `callshape.txt` · `nilcheck.txt` · `prod-shape.txt` · `sentences.txt` · `recovers.txt` · `spawn-sites.txt` · `d22scan.txt` · `verdict-s2.txt` · `gitstatus-full.txt` · `236-verdict-head.txt`。

### §0-补2 读尺纪律（每把尺锚在带括号的调用形状上）

```
grep -rn "\.BaseOptions(" --include=*.go .
  → 19 命中：产码真调用 1 处 ＝ internal/tools/subagent_197.go:278
    （17 处在 .scratch 的历史／突变副本，1 处是测试注释 failclosed_236_teeth_test.go:203）
grep -rn "BaseOptions == nil" --include=*.go .
  → 19 命中：产码判空 1 处 ＝ internal/tools/subagent_197.go:256
    （17 处同上＝.scratch 副本，1 处＝failclosed_236_teeth_test.go:12 的注释）
grep -rn "BaseOptions != nil" --include=*.go .        → 0 命中（无第二支正向判空）
```

⚠ 对照：裸 `grep -rn "BaseOptions" --include=*.go .` ＝ **194 命中**，其中 **172 命中落在 `.scratch` 下不可编译的副本里**（含一枚嵌套重复目录 `.scratch/.scratch/wisp/probes/222/v1/mutations/asis/subagent_197.go`）。⇒ 派单红线成立：**裸符号名会把历史突变副本／注释／类型声明／构造点全数当"调用者"**。本件所有"有没有调用者／有没有防护"的结论只用上面三把锚定尺，且**双向都追到 main**。

**正向（生产唯一产出者，追到 main）**：`cmd/wisp/run.go:789` ＝ `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },`，位于 `tools.BuiltinSubagentEntries(tools.SubagentDeps{…})`（`:780`）实参内；该语句在 `func assembleRuntime(s runSpec) (*agentRuntime, int)`（`cmd/wisp/run.go:382`）里；链路 `main`（`cmd/wisp/main.go:58`）→ `os.Exit(cmdRun(args[1:]))`（`:91`）→ `cmdRun`（`:151`）→ `runTextTask`（`:181`）→ `assembleRuntime` ⇒ **闭合；生产产出者恰 1 枚，与派单给的 `run.go:789` 对上。**
配套（"装了函数但没快照"那一维）：`loopOptSet` 全仓只有三处＝`cmd/wisp/run.go:293` 声明・`:789` 读・`:1022` `rt.loopOpt, rt.loopOptSet = opts, true` ⇒ **`:789` 那枚闭包永不为 nil，它只是可能在 `:1022` 之前回 `(_, false)`**——那是 `:279 if !wired` 管的事，与 `:256` 管的"nil 函数值"**不是同一件事**（＝AC#1 的语义边界）。

**负向（别处判空／别处喂它，也追到 main）**：`grep -rn "SubagentDeps{" --include=*.go cmd internal` ＝ **7 处**；`grep -rn "BuiltinSubagentEntries" --include=*.go cmd internal` ＝ 8 命中（产码 2＝定义 `subagent_197.go:233`＋注释 `:229`；调用＝产码 1＋测试 4）。分布＝产码 1（唯一接 `BaseOptions` 的）／测试 6：接**非 nil** 产出者 3（`internal/tools/subagent_197_test.go:180`+`:185`、`internal/tools/subagent_222_test.go:299`+`:303`、`internal/tools/task_cancel_221_legs_test.go:86`+`:91`）、留 **nil** 的 2（`internal/tools/task_cancel_221_test.go:51`＝`SubagentDeps{Roster: roster}` 只枚名；`cmd/wisp/subagent_selfapproval_197_test.go:551`＝`reflect.TypeOf(tools.SubagentDeps{})` 只做类型遍历）。⇒ 两枚 nil‑deps 形状都**不派发**，今天摸不到 `:278`。
★**新发现（派单未列）**：`internal/tools/task_cancel_221_test.go:51` 是仓里把 `BaseOptions` 留 nil 的**第 4 处**（第 5 处＝`subagent_selfapproval_197_test.go:551` 的类型遍历形）；它只读 `Tool.Name()`、不调 `Execute` ⇒ 今天不崩，但属"未兜面"。

---

## §1 形状与逐字原文（行号＝本腿 `grep -n`／`sed -n` 现取）

### 1.1 五枚锚点全部核实：**编排者派单给的锚点一枚未漂**

| 派单锚 | 现量 | 逐字原文 |
|---|---|---|
| `:135` 声明 | 对上 | `BaseOptions func() (agent.Options, bool)`（函数值字段 ⇒ nil 调用＝panic） |
| `:256` guard | 对上 | `if t.d.BaseOptions == nil \|\| t.d.ParentTools == nil {` |
| `:257` 拒绝体（派单写作"…拒绝…"） | 对上 | `return Result{Text: "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）", IsError: true}, nil` |
| `:278` 无条件调用 | 对上 | `base, wired := t.d.BaseOptions()` |
| `:279` 未接线支 | 对上 | `if !wired {` → `:280 release()` → `:281 return Result{Text: "宿主现在给不出装配快照（fail-closed：不在没有父装配的情况下起子环路）", IsError: true}, nil` |

同一 `Execute` 体内的次序（`sed -n '240,300p' internal/tools/subagent_197.go` 现读）：`:253` `if t.d.Roster == nil {`（体 `:254`）→ **`:256` 装配 guard** → `:259 parentID := CorrelationID(ctx)`／`:260 if parentID == "" {`（体 `:261-262`）→ 深度上限支（`MaxSubagentDepth`）→ 池满支（`TryAcquireSubagentSlot(MaxConcurrentSubagents)`＋那句"同时在跑的子代理已达上限 %d 枚"）→ **`:278` 调用** → `:284 if base.AdmitTask == nil {`（"子代理永不自批"支）。
⇒ **`:278` 之前没有任何一处对 `BaseOptions` 的判空**；`:279` 的 `!wired` 管的是"装了函数但取不到快照"，**管不到 nil 函数值**（nil 函数值在 `:278` 当场崩，走不到 `:279`）。⇒ **票面"这道 guard 一票两形"的成立性：本腿独立复认＝成立。**
⇒ 另记一枚同类事实（同一把尺现取）：`t.d.ParentTools` 在 `:293` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` 被**无条件使用**，同样只靠 `:256` 那一支拦着 ⇒ **`:256` 今天其实是"一票三形"（`:278` 的 nil 防护＋`:293` 的 nil 防护＋那句诚实拒绝文案）**；差别只在 `newSubagentToolProvider(nil)` 会不会自己崩——本腿读它需要实跑才能定性，⛔ 不许推测（见 §4 第 6 条）。

### 1.2 `BaseOptions` 命中分类（＝改动面分母；只算 tracked 的 `cmd/`＋`internal/`）

`grep -rn "BaseOptions" --include=*.go cmd internal` ＝ **22 命中／7 文件**：

| 类 | 枚数 | 逐枚性质 |
|---|---|---|
| **产码非测试** | **8 枚／2 文件** | `internal/tools/subagent_197.go` 6＝`:26`・`:132`・`:134` 注释・`:135` 声明・`:256` guard・`:278` 调用；`cmd/wisp/run.go` 2＝`:291` 注释・`:789` 构造（生产唯一）。另 `:257` 是含该字样的**文案行** |
| **`_test.go`** | **14 枚／5 文件** | `internal/tools/failclosed_236_teeth_test.go` 10＝`:12`・`:30`・`:201`・`:203`・`:204` 注释・`:95` 完整字面量常量・`:208` 用例名・`:212` 变异（`d.BaseOptions = nil`）・`:216`・`:220` 红句；三枚测试产出者 `subagent_197_test.go:185`・`subagent_222_test.go:303`・`task_cancel_221_legs_test.go:91`（**派单给的三枚行号逐字对上，一枚未漂**）；`cmd/wisp/subagent_selfapproval_197_test.go:544` 注释 1 枚 |
| **注释文档（md，全仓）** | 票 270 面 4・`probes/236/v1/verdict.md` 11・`probes/236/r2/evidence.md` 14・`docs/reports/pending-and-issues.md` 3・`docs/evidence/s1/*` 4・`probes/236/a1/census.md`・`probes/235/**`・`probes/221/c1/census.md`・`probes/167/a3/census.md` 等 | ⛔ 本腿不判其内容，只按 §2.4 引用 `236-v1` 的读数 |
| **`.scratch` 下不可编译的 `.go` 副本** | 调用尺 17・判空尺 17・裸符号尺 172 | 不参与编译，但**会被裸 grep 数进去**（§0-补2）；⛔ 本腿未在 `.scratch` 落任何 `.go`（票 236 AC#6 登记的正是"在 `.scratch` 放故意编译不过的样本"这种雷） |

**改动面分母结论**：三形的**产码射程恰＝1 枚文件里的 4 枚行**（`internal/tools/subagent_197.go` 的 `:135`・`:256`・`:257`・`:278`）；测试射程由 §2.2 逐枚点名（`failclosed_236_teeth_test.go`＝唯一顶红载体，`subagent_197_test.go:879`＝字段数红线，`cmd/wisp/subagent_selfapproval_197_test.go:551`＝类型遍历红线）。

### 1.3 panic→零读数那一形的现有兜底落在哪几枚文件（本腿实测枚数）

`grep -rn "recover()" --include=*.go internal cmd` ＝ **17 命中／14 文件**。按"能不能兜住 `:278` 那一次空函数值调用"分三类：

| 类 | 位置∶枚数 | 读数 |
|---|---|---|
| **在 spawn 派发路径上＝今天唯一的兜底** | `internal/tools/failclosed_236_teeth_test.go` ∶ **2 枚**（`spawnGuarded` 的 `:232-236`、`spawnWithoutHostTaskID` 的 `:282-286`） | 经它们发出的 `task.spawn` 派发点＝`:172`／`:190`／`:214`／`:250`（**4 处**，含定义处引用）。panic 被转成 `Result{Text: fmt.Sprintf("panic（guard 摘掉后不是拒绝而是崩）: %v", rec), IsError: true}` ⇒ 具名红 1 枚；该串⛔ 永不等于任何 want ⇒ 造不出绿 |
| **包内但兜不住** | `internal/tools/bridge.go:823`（`safeFacts` 的 recover）∶ 1 枚 | 只包 risk 事实钩子那一次调用，⛔ 不包 `Execute` 本体 ⇒ 结构上拦不到 `:278`（原文：`f = risk.Facts{Irreversible: []string{fmt.Sprintf("宿主事实钩子 panic: %v", rec)}}`） |
| **不在本票射程** | `cmd/balldebug/main.go` 1・`cmd/wisp/panel_host_windows_test.go` 1・`cmd/wisp/panel_resident_windows_test.go` 2・`internal/agent/approval/fakes_test.go` 1・`internal/agent/approval/ticket220_l1_window_read_test.go` 1・`internal/ball/hotkey_reload.go` 1・`internal/ball/sta_release_windows_test.go` 3・`internal/memory/writer.go` 1・`internal/observe/goroutine.go` 1・`internal/observe/sampler_faketree_guard_136_test.go` 1・`internal/panel/subagent_stream_197_test.go` 1・`internal/risk/assessor.go` 2・`internal/tools/grant.go` 1・`internal/tools/mode.go` 1 | 与 spawn 调用栈无交 |

**未兜的暴露面（本腿数的是"派发点枚数"，⛔ 不是用例枚数）**：`grep -rn "h.spawn(\|\.spawn(t, \|spawnResult(\|\"task.spawn\"" --include=*_test.go internal/tools cmd/wisp` ＝ **28 处／5 文件**＝`internal/tools/subagent_197_test.go` 20・`internal/tools/task_cancel_221_legs_test.go` 3・`internal/tools/failclosed_236_teeth_test.go` 3・`internal/tools/subagent_222_test.go` 1・`cmd/wisp/subagent_carrier_197_test.go:174` 1（走真 `bridge.Execute`）。⇒ **只有 `failclosed_236_teeth_test.go` 那 3 处包了 recover，其余 25 处裸奔**；今天不崩只因三处产出者（`:185`／`:303`／`:91`）恒接**非 nil** `BaseOptions`。
⚠ 点名一枚：`internal/tools/subagent_197_test.go:852` ∶ `res := h.spawn(t, t.Context(), label197, "没有审批层就想派一枚")` 属那 20 枚未兜用例——任何"让它拿到 nil `BaseOptions` 而 `:278` 未装甲"的改动落到它头上＝**不报红，打死测试二进制**。

---

## §2 三形代价表

> 口径：⛔ 本腿禁跑 `go` ⇒ "顶红"＝**读断言原文的静态推演**，不是实跑红名册；每格点名**被顶到的那一句**。"拆完之后剩几枚零读数"的**实测**本腿交不出（§2.3＋§4 第 1 条）。

### 2.1 改动面（文件∶行）／是否动文案／是否需新导出名

| 形 | 动哪几枚文件 | 射程细节 | 动文案？ | 新导出名？ |
|---|---|---|---|---|
| **甲**＝只给 `:278` 装甲判空（nil ⇒ 具名拒绝），`:256` 原样 | 产码 **1 文件**：`internal/tools/subagent_197.go`（`:278` 前插判空＋体内引用 `:257` 那句）。测试 **0 文件** | ⚠ 新分支**今天不可达**（`:256` 已在它上游把所有 nil 拦走）⇒ 对现有用例**零新读数**；它的牙只在"overlay 摘 `:256`"的台件里看得见（＝票面 AC#1 那条台件）。另：甲**不解决** `:293` 那一形（见 §1.1 末段） | **否**（nil 支复用 `:257` 现串 ⇒ 零文案改动；另写一句则顶 AC#3"三句互不充当彼此替身"，见 §2.5） | **否**（同包控制流；票面禁区亦禁新增导出名） |
| **乙**＝`:256` 拆两支＋**同批**给 `:278` 装甲 | 产码 **1 文件**（`:256-258` → 两支＋`:278` 判空）＋测试 **≥1 文件**＝`internal/tools/failclosed_236_teeth_test.go`：常量 `:95-96` 一拆为二、两枚用例的 want（`:195`、`:218`）与红句（`:196-197`、`:219-220`）逐枚改、文件头注释 `:12`・`:30`・`:201-207` 同步；`subagent_197_test.go`（有无红取决于是否只拆文案不动字段，现读＝不动字段则不红） | 票面禁区第 3 条＝拆语句与装甲**不许分批**（先拆后补＝中间态可崩，正是本票立案理由）⇒ 乙**只能一发落地**，落地发内改动面＝产码 1＋测试 1（＋可能的注释若干） | **是**——两半各有自己的文案才谈得上"拆得动"（AC#2），于是顶红 §2.2 第 1/2/3 行；那批钉是**票 236 刚翻勾的凭据** ⇒ 改＝同批改钉＋登记，⛔ 不许为变绿放宽（票面禁区第 1 条） | **否**；⛔ 也不许用"加字段／加类型"表达装配状态（§2.2 第 6/7 行两枚红线钉） |
| **丙**＝不动结构，只在 `:256`／`:278` 各加一条**具名常驻判据** | 产码 **1 文件**：`:256` 处注释（票面 AC#5 那句"这支同时是 `:278` 的 nil 防护"）＋`:278` 处判空／断言性拒绝。测试侧：只加注释 ⇒ **0 文件**；判据要"常驻有牙"（不靠 overlay 台件）⇒ **≥1 文件**（新增一枚自带 recover 的用例，或改 `:212` 那枚变异载体） | ⚠ **丙与甲的差只剩那条注释**：AC#5 的注释是丙的必要非充分件；"改由产码层拒收"落地＝甲的 `:278` 判空 ⇒ 选丙时改动面＝甲＋注释，**不额外省任何事**。★"只加注释"那一支**没有仪器钉**（§2.2 第 8 行） | **否**（＝甲的全部代价＋注释） | **否**（做成新函数／新类型／新字段 ⇒ 撞 §2.2 第 6/7 行两枚钉＋票面禁区"不许新增导出名"） |

**三形共同不覆盖的两格（选形时必须看得见，否则是在猜）**：

1. 票面 §26 判据＝"拆掉任一半 guard 之后 `internal/tools` 整包不许出现零读数"。三形解决的都是 **BaseOptions 半边摘 ⇒ 崩**这一维；**ParentTools 半边摘后不是拒绝而是"真派生"**（`236-v1` 的 `v-m1d1` 读数 got 逐字＝`子代理 83c0c55a-…（父工具面没接线）已结束，状态 Settling。`），这一维**三形都不动**（乙拆支后才谈得上给它独立文案＝票面 AC#3）。⇒ **甲/丙解"失明形"、不新增 AC#2 的牙；乙同时动 AC#2，代价＝改票 236 那批字面量钉。**
2. **`:293 opt.Tools = newSubagentToolProvider(t.d.ParentTools)` 也是只靠 `:256` 拦着的无条件使用**（§1.1 末段）。若选乙并拆支，`ParentTools` 那一支摘掉后 `:293` 的行为（崩／还是静默给空目录）**不在三形任何一形的射程里**——票面 §26 那句"整包不许出现零读数"若也覆盖 `ParentTools` 半边，那甲乙丙**都还不够**。

### 2.2 ★顶红谁现有的钉（逐枚：文件∶行∶断言逐字∶甲／乙／丙会不会红＋为什么）

| # | 文件∶行∶断言逐字 | 甲 | 乙 | 丙 | 被顶到的是哪一句＋为什么 |
|---|---|---|---|---|---|
| 1 | `internal/tools/failclosed_236_teeth_test.go:95-96` ∶ `lit236SpawnNoAssembly = "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，" + "fail-closed：不起第二套运行时）"` | 不红 | **红**（两半文案分开时） | 不红 | 该常量是第 2、3 两枚分支钉**共用的 want**：本身不断言，一改它＝同时顶到第 2 与第 3 枚。乙必须同批改＝触碰票 236 刚翻勾的字面量凭据 |
| 2 | `…:184` `Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired`（变异 `:188 d.ParentTools = nil`）∶ `:192 if !res.IsError { t.Errorf("ParentTools 没接线时 task.spawn 回了非拒绝：%q（guard 不在时子代理照样派生，只是拿到一个空目录）", res.Text) }` ＋ `:195 if res.Text != lit236SpawnNoAssembly {`（红句 `:196-197` ∶ `"subagent_197.go:256 那一支的完整字面量没被原样说出（got %q, want %q）——这一支是 \|\| 的 ParentTools 半边"`） | **不红**（`:256` 未动 ⇒ 仍由整支拦，got 逐字等 want） | **红** | 不红 | 被顶到的是 **`:195` 那枚等值比**（`!=`，不是 `Contains`）：乙给 ParentTools 半边换文案 ⇒ got≠旧 want 立红；`:192` 还同时要求"仍必须是拒绝"，不许改成成功回执形 |
| 3 | `…:208` `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired`（变异 `:212 d.BaseOptions = nil`）∶ `:215 if !res.IsError { t.Errorf("BaseOptions 没接线时 task.spawn 回了非拒绝：%q", res.Text) }` ＋ `:218 if res.Text != lit236SpawnNoAssembly {`（红句 `:219-220` ∶ `"subagent_197.go:256 那一支的完整字面量没被原样说出（got %q, want %q）——这一支是 \|\| 的 BaseOptions 半边；got 以「panic」开头就是那道 guard 已经不在了"`） | **不红**（同上：`:256` 仍先手拦） | **红**（半边文案一变） | 不红（只加注释／加今天不可达的判空时）；若丙顺手改写 `:256` 文案则红 | 被顶到的是 **`:218` 的等值比**。★另有**非断言级代价**：`:220` 把"got 以 panic 开头"直接读成"guard 已经不在了"；甲/丙装了 `:278` 判空后，这枚用例上 **panic 那一形永远不再出现** ⇒ 该诊断句退化成"只在 overlay 台件里成立"的**过期叙述**（注释层漂移，不是断言层红） |
| 4 | `…:230-238` ∶ `func spawnGuarded(…) (out Result) { defer func() { if rec := recover(); rec != nil { out = Result{Text: fmt.Sprintf("panic（guard 摘掉后不是拒绝而是崩）: %v", rec), IsError: true} } }() return h.spawn(t, t.Context(), label, prompt) }`（另 `:280-296` `spawnWithoutHostTaskID` 同形） | 不红 | 不红 | 不红 | **没有任何断言钉"recover 必须被触发"** ⇒ 三形都不因"panic 消失"而红。⚠ 代价在注释层：`:225-229`・`:201-207`・文件头 `:30-36` 都把"这一支摘掉＝崩"写成事实，甲/丙落地后即过期；**改票 236 的凭据注释属 AC#0 之外的射程，须编排者定**（本腿只点名，不代改） |
| 5 | `…:166-181`（roster 支）∶ `:176 if res.Text != lit236SpawnRosterUnwired {`（红句含 `"若 got 是「同时在跑的子代理已达上限」，说明 Roster == nil 的拒绝分支已经不在了"`）；`…:244-273`（parentID 支）∶ `:257 if res.Text != lit236SpawnNoHostTaskID {` ＋ `:269 if n := h.roster.Count(); n != 1 {`（红句 ∶ `"名册行数 = %d, want 1（只有 harness 那枚根任务）——宿主没给父任务 id 时不许派生、更不许登记：%v"`） | 不红 | 不红（**前提＝guard 次序不动**） | 不红 | 三形都不碰 `:253`／`:259-263`。⚠ 但乙把一枚 `\|\|` 拆成两支时**先后次序决定"哪一支先说话"**：这三枚各钉一句完整字面量，次序挪了就得重跑它们的台件（票面 AC#3"三句互不充当替身"正是这一族） |
| 6 | `internal/tools/subagent_197_test.go:879` ∶ `if got := subagentDepsFieldNames(); len(got) != 5 {`（红句 ∶ `"SubagentDeps 字段集合变了：%v（多出来的每一枚都要问一句是不是新的允许出口）"`；定义 `:834-841` 用 `reflect.TypeOf(SubagentDeps{}).NumField()`） | **不红**（判空＝局部控制流） | **不红** | **不红** | ★**红线一枚**：任何"用 `SubagentDeps` 新字段表达装配状态"（例 `AssemblyWired bool`）**直接顶红 `:879`**，并把讨论推进到票 197 的"子代理不许有第六枚旋钮"。⚠ 同处注释 `:877-878` 写的是"字段只有装配、目录、溯源、流**四枚**"，现量已是 **5** ⇒ 注释与钉**已漂**（本腿只报现象，不定性）。⇒ **三形都不许加字段**，否则本表重算 |
| 7 | `cmd/wisp/subagent_selfapproval_197_test.go:551` ∶ `{"tools.SubagentDeps", reflect.TypeOf(tools.SubagentDeps{})},` ＋ `:563` 红句 ∶ `"从子代理的装配里能静态走到声明了允许出口的类型：%v。孩子只要拿得到这些名字里的任何一个，「子代理永不自批」就退化成一句约定"`（禁止方法名册 `:541` ∶ `{"Allow","Native","DecideFromNative","DecideFromPanel","GrantNonce"}`） | 不红 | 不红 | **不红（只加注释时）**；⚠ 丙若把判据做成**新增导出类型**（包一层拒绝器）⇒ 这枚与 `:879` **同红**＋撞票面"不许新增导出名" | 它走 `findAllowDoors197` 的字段＋方法集遍历：**只有新增字段／类型／方法会被扫到**；改函数体、加注释扫不到 |
| 8 | `internal/tools/task_cancel_221_legs_test.go:181` ∶ `t.Fatalf("BuiltinTaskEntries 没注册 task.cancel（现名册：%v）—— 说明书又在许诺一枚不存在的工具", names)` ＋ `:184 if len(names) != 2 {` ＋ `:224 src, err := os.ReadFile("task.go")` | 不红 | 不红 | 不红 | 名册尺只吃 `task.go` 的注册行与 `task.cancel` 文案；`:224` 这把**读盘尺读的是 `task.go`，不是 `subagent_197.go`** ⇒ 本票三形在它盲区外。★对丙的直接后果：`grep -rn "subagent_197.go\"" --include=*_test.go internal cmd` ＝ **0 命中**；`grep -rn "go:embed" --include=*_test.go internal/tools` ＝ **1 命中**（`tasklist_deferred_236r3_teeth_test.go:109` ∶ `//go:embed task.go`）⇒ **仓里没有任何读盘／embed 尺吃 `subagent_197.go` 的字节** ⇒ **丙"只加注释"那一支今天没有任何仪器能钉住它是否真存在**（写没写、漂没漂都扫不出） |
| 9 | `internal/tools/subagent_197_test.go:846-854`（`Test197SubagentHasNoSelfApprovalOutlet`）∶ `res := h.spawn(t, t.Context(), label197, "没有审批层就想派一枚")` ＋ `if !res.IsError \|\| !strings.Contains(res.Text, "永不自批") { t.Fatalf("宿主没有准入钩子时派生成功了（子代理就有了自己的出口）：%q", res.Text) }` | 不红 | 不红 | 不红 | 走 `:284 base.AdmitTask == nil` 支，harness `:185` 恒回非 nil ⇒ 三形都在它下游。**★但它是 §1.3 那 20 枚"未包 recover"的 `h.spawn(` 之一**：`:278` 若仍无判空，任何让它拿到 nil `BaseOptions` 的改动＝打死二进制而非报红 |
| 10 | `internal/tools/subagent_222_test.go:299-317`（`wireSpawn`）与 `internal/tools/task_cancel_221_legs_test.go:86-101` ∶ 各自 `BaseOptions: func() (agent.Options, bool) { … }` 恒回 `true` | 不红 | 不红 | 不红 | 两枚是**产出者**不是断言 ⇒ 三形都摸不到它们。⚠ `subagent_222_test.go:320` 的裸 `go func(i int)`（该文件注释自称 d22scan ban 1 只施于非测试文件）与三形无关，只提醒别顺手改（§2.6） |
| 11 | `internal/tools/task_cancel_221_test.go:51` ∶ `BuiltinSubagentEntries(SubagentDeps{Roster: roster}),`（名册尺，只读 `e.Tool.Name()`） | 不红 | 不红 | 不红 | **派单未列的第 4 枚 nil‑BaseOptions 形状**：构造 deps 但**不派发** ⇒ 即使 `:278` 无判空也不崩。⚠ 哪天有人让它真 `Execute` 一次 spawn，它就落在无 recover 的面上 |

**枚数小结**（本表点名到的现有钉共 **11 行**；因形而红的枚数）：

- **甲＝0 枚红**。代价不体现在红，体现在两处：①新分支今天不可达 ⇒ 必须靠一枚 overlay 台件才有读数（AC#1 的完成判据）；②`failclosed_236_teeth_test.go` 里 4 处注释（`:30-36`・`:201-207`・`:219-220` 的后半句・`:225-229`）成为过期叙述。
- **乙＝3 枚具名红**＝`failclosed_236_teeth_test.go:95-96`（共用常量）＋`:195`（ParentTools 等值钉）＋`:218`（BaseOptions 等值钉）；另 **3 枚次序敏感钉**（`:176`／`:257`／`:269`）需重跑台件。乙是唯一"能同时把 AC#2 的牙建出来"的形，也是唯一**必须动票 236 凭据**的形。
- **丙＝0 枚**（只加注释／加今天不可达判空的形）；若越界加类型或字段 ⇒ **＋2 枚**（`subagent_197_test.go:879`、`cmd/wisp/subagent_selfapproval_197_test.go:563`）。丙"只加注释"那一支的判据**无仪器可钉**（第 8 行）。

### 2.3 "拆完之后崩溃面还剩几枚零读数"——**本腿交不出实测**（禁跑 `go`），只给静态结论

- 三形任一（`:278` 装了判空）⇒ **摘 `:256` 整支**（`A-4` 那一形）时 `BaseOptions==nil` 走具名拒绝 ⇒ 不再产生 `^panic` ⇒ "整包零读数"形消失；两枚分支钉各报具名红，名册仍全（＝票面 AC#1 要的"仍出全量名册，PASS+FAIL ≈ 206"）。
- 乙另加"半边摘 ⇒ 恰 1 枚具名红"的形状（AC#2）；甲/丙不新增 AC#2 的牙。
- **残余零读数面（静态）**＝§1.3 末段那 **25＋1 处无 recover 的派发点**＋`task_cancel_221_test.go:51` 那枚 nil‑deps 名册尺。三形**都不动它们**：那一片今天靠"三处 harness 恒接非 nil `BaseOptions`"这条**测试侧不变式**兜着，⛔ 不靠产码。⇒ 若编排者要"零读数的最后一层挪进产码"，那是**第四形**（票面没有），本腿不自造，只在 §3 摆成一道要人拍的题。
- ⛔ 上面全是**推**，不是**量**。AC#0 原文里"以及『拆完之后崩溃面还剩几枚零读数』的实测"这一件，**本件没交**，因为本腿被禁跑 `go`（派单红线）。⇒ 这一格必须由**另一枚能跑 `go` 的腿**（或编排者自己在核过后）补实测；补的时候参照 §2.4 的配对法（正控＝只换测试文件、产码一字不动）。

### 2.4〔引用件，非本腿实测〕`236-v1` 的崩溃读数（⛔ 只取读数，不取判语）

出处＝`.scratch/wisp/probes/236/v1/verdict.md`（§2 表内 `:119`／`:132` 行、§3 名册 `:153`／`:161`／`:163`／`:164` 行、§5 第 2 条 `:255-258` 行）。

- `A-4-without-recover`（`:256` 改 `if false` ＋ overlay 换掉 `failclosed_236_teeth_test.go:232-236` 那五行 recover）＝**rc=1／PASS=25／FAIL=2／日志 `^panic`＝2 行／`FAIL github.com/CarlosShao/wisp/internal/tools 1.444s`** ⇒ 其后约 **181 枚用例零读数**（206−25）。
- `ctrl-no-recover-only`（★**只**换测试文件、产码一字不动）＝**rc=0／PASS=206／FAIL=0** ⇒ 摘 recover 这个动作本身不产红、不产崩 ⇒ 上面那份失明只能归因于产码那一次空函数值调用。
- `A-4`（同一产码突变＋**保留** recover）＝**rc=1／PASS=204／FAIL=2**（两枚具名红，其一 got 逐字以 `panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference` 开头）。
- `v-m1d1`（`:256` 只留 BaseOptions 半边）＝FAIL=1＝`…WhenParentToolsIsUnwired`（got＝成功回执形）；`v-m1d2`（只留 ParentTools 半边）＝FAIL=1＝`…WhenBaseOptionsIsUnwired`（got 以 panic 串开头）。
- 该件 §0-补2 记的基线：`internal/tools` 顶层 **PASS=206／FAIL=0／SKIP=0**，且它声明了口径——`206` 是 `grep -c '^--- PASS'` 的**顶层枚数**，行首锚定看不见缩进子测试。
- ⛔ 本腿**不复用**该件 §2 第 6 行与 §5 第 2 条里的任何**定性句**（"算不算真红／算不算掩盖"属编排者射程）。

### 2.5 文案面（票面 AC#3 只列事实，⛔ 不判对错）

三句现量（`grep -rn` 逐字，均在 `internal/tools/subagent_197.go`）：

- 装配未接线（`:256` 支）＝`:257` "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）"
- 快照取不到（`:279` 支）＝`:281` "宿主现在给不出装配快照（fail-closed：不在没有父装配的情况下起子环路）"
- 池满（`:270` 支）＝"拒绝派生：同时在跑的子代理已达上限 %d 枚（在跑的：%s）。这是硬拒，不是排队——…"
- 无父任务 id（`:260` 支）＝`:261-262` "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，就没法登记父子关系…）"

★三形对文案的处置差异（AC#3 的射程）：甲/丙 **零文案改动**（nil ⇒ 复用 `:257` 现串，AC#3 天然不违）；乙 **必须新写一句**（否则拆支无意义）⇒ 乙落地时"装配未接线"会从**一句**变**两句**，而票面 AC#3 只要求"装配／池满／无父任务"三句互不充当替身，**没规定拆开后那两句之间算不算替身**——这是乙留给编排者的一道未定义项（§3 第 4 条）。
另记：`:257` 那句文案里**同时写着两个字段名**（`BaseOptions/ParentTools`）；乙拆支后这句要么改名要么留一句作废——`failclosed_236_teeth_test.go:12` 与 `:95-96` 都逐字含这串，属第 1 行同源。

### 2.6 `d22scan` 会不会拦三形里的某种写法（⛔ 未跑它，只读 `tools/d22scan/main.go` 原文）

本票相关 9 枚 ban（原文抄自 `tools/d22scan/main.go:7-70`，行号＝现取）：

| ban | 原文射程要点 | 对三形的判定 |
|---|---|---|
| **1 bare-goroutine** | `every \`go <anything>\` - closure literal OR named call (widened by R16#1)`；scope＝`internal/`＋`cmd/` **非测试、非 testdata**；唯一文件级豁免＝`internal/observe/goroutine.go`（见 `allowlist.txt` 逐字 ∶ `bare-goroutine	internal/observe/goroutine.go	R16#1: this file implements Registry itself…`） | ⛔ **三形都不新增 goroutine** ⇒ 不拦。（旁证：`subagent_222_test.go:320` 的裸 `go func` 因"non-test"豁免而不在射程，本腿只读到注释自述，⛔ 不实跑验证） |
| **2 pathresolver-bypass** | `filepath.Clean / filepath.Abs outside the allowlist` | ⛔ 三形都不碰路径（改动面只读 `internal/tools/subagent_197.go` 一个函数体） ⇒ 不拦 |
| **3 plaintext-key** ／ **4 wallclock-timeout** ／ **5 mirror-hash** ／ **6 panel-approval（scope frontend/）** ／ **7 internal-artifact-tool（scope internal/tools/，指 gated tool 名）** | 形状＝密钥字面量／墙钟差／镜像哈希／前端 `approval.decide`／把宿主内部 artifacts 写入实现成**受门控的 Tool 名** | ⛔ 都不沾。**特别提醒 ban 7**：乙/丙若"顺手"把那句拒绝做成一枚新工具面或新 entry 名，正好撞它（票面禁区第 2 条同源＝D34 无对应行） |
| **8 emoji** | scope＝`emojiScopes()`＝`design/`＋`frontend/`＋`internal/`＋`cmd/`（`goOnly` 对 `internal/`＋`cmd/`），**`_test.go` 与注释都在射程内**；`emojiRe`＝`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`（`main.go:186` 逐字）；匹配前 `blankCommentRanges()` **把注释抹白**（`:1150-1161`，注释里的 `≤` 不判、串里 `<=` 的 `=` 本就不在射程） | ⚠ **唯一真的会拦三形的 ban**：注释豁免，**字符串常量不豁免**。⇒ 甲／乙／丙的**新文案里若出现 `≤`（U+2264）、`✓`（U+2713）、`→`（U+2192 不在 ban 的带里、仪器抓不到）、`①`（U+2460 亦不在带里）**这类字符，落在 `Text:` 字符串上＝**CI 直接红**。乙最危险（它要新写一句文案）。⇒ 派单"文案零改动"这一条在甲/丙天然满足；乙必须由编排者先定文案写法，⛔ 别让写腿临场造句 |
| **9 phantom-citation** | 产码文件（`internal/`＋`cmd/`，**非 `_test.go`、非 `tools/**`**）注释里写了 `docs|\.scratch|internal|cmd|tools|scripts/` 开头的**全拼路径**，`os.Stat` 判它存不存在（`main.go:880-891`＋`:860` 逐字 ∶ `s.add("phantom-citation", path, pos.Line, "comment cites repo path \""+tok+"\" which does not exist on disk (ticket 212 ban #9)…")`）；带 `…`／`*` 的简写形被 `shorthandRegionRe` 豁免（不判、也不许当判据） | ⚠ **直接打在丙的 AC#5 注释上**：那条"这支同时是 `:278` 的 nil 防护"的注释**一旦引用仓库路径**就必须逐字存在。本腿实测盘上：`docs/PLAN.md` **EXISTS**、`docs/SLO.md` **EXISTS**、`docs/reports/HANDOVER.md` **EXISTS**；`docs/DECISIONS.md`・`docs/DEFERRED.md`・`docs/TOOLS.md`・`docs/STATE_MACHINE.md`・`docs/contracts/` **全部 MISSING**；`docs/specs/` 真名册＝`SPEC-00-product-overview.md`…`SPEC-05-agent-core.md`（⚠ **不是** `SPEC-05-agent-kernel.md`）、`SPEC-11-build-deploy-containerization.md`。⇒ 写注释时**照上面这张真名册拼**，或直接不引路径（只写 `internal/tools/subagent_197.go:256` 这种"包内行号"形状，`repoPathRe` 不吃） |

`tools/d22scan/allowlist.txt` 现读**共 5 行豁免**（逐字全文见 `logs/`；含 `pathresolver-bypass` 3 枚、`mirror-hash` 1 枚、`bare-goroutine` 1 枚），**没有一条与本票相关** ⇒ 别指望白名单替三形开_any_口子。
CI 侧接线：`.github/workflows/ci.yml:74` `D22 scanner positive control (tools/d22scan tests, seeded red)`＋`:81` `run: bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（**本腿只 grep 到这两行原文，未读全文、未跑**）。

---

## §3 要编排者裁的清单（⛔ 判语一律留空，本腿不预选任何形状）

> 用法：每条只给"题＋本腿量到的料＋空白判语格"。⛔ 本腿不填、不推荐、不排序。

| # | 待裁的题（只问编排者能答的那一层） | 本腿给的料 | 判语（留空） |
|---|---|---|---|
| 1 | **选哪一形**：甲（`:278` 装甲、零文案改动、0 枚红）／乙（拆两支＋同批装甲、3 枚具名红＋3 枚次序敏感钉）／丙（甲＋一条 AC#5 注释、0 枚红但注释无仪器可钉）？ | §2.1 改动面三行＋§2.2 枚数小结 | |
| 2 | **要不要"第四形"**：票面 §26 那句"拆掉任一半 guard 后整包不许零读数"里，**ParentTools 半边**那一维三形都不解（`:293` 也是无条件使用）；而**残余 25＋1 处无 recover 的派发点**今天靠测试侧不变式兜着。要不要把它挪进产码？ | §1.3 末段＋§2.1"共同不覆盖的两格"第 2 条＋§2.3 | |
| 3 | **AC#0 的实测那一格谁跑**：本腿被禁 `go`，"拆完之后剩几枚零读数"只有静态结论＋`236-v1` 的引用读数。是否另派一枚**能跑 `go` 的读数腿**（照 §2.4 的配对法：产码突变一发＋"只换测试文件"正控一发）？ | §2.3 末段＋§2.4 | |
| 4 | **乙的新文案怎么写**：票面 AC#3 只禁"装配／池满／无父任务三句互为替身"，**没规定拆开后那两句装配文案之间算不算替身**；且 `:257` 现串里同时写着两个字段名（`BaseOptions/ParentTools`），`failclosed_236_teeth_test.go:12`＋`:95-96` 逐字含这串。拆开后那两句是编排者定，还是写腿定？ | §2.5 全段 | |
| 5 | **票 236 那批凭据（注释＋常量＋两枚等值钉）允许被乙同批改吗**：票面禁区第 1 条写的是"⛔ 不许为了变绿放宽"；乙的改**不是放宽而是换串**——这条边界要编排者点名批准（改的是刚翻勾的凭据，属"改凭据"级别，⛔ 不是本腿能定的） | §2.2 第 1/2/3/4 行＋票面禁区第 1 条 | |
| 6 | **过期注释归谁处置**：若选甲或丙，`failclosed_236_teeth_test.go` 的 `:30-36`・`:201-207`・`:219-220` 后半句・`:225-229` 会全部变成"描述已不成立的崩溃"。要不要同批改？改测试文件的注释属不属于"动票 236 凭据"？ | §2.2 第 3/4 行＋§2.3 | |
| 7 | **`subagent_197_test.go:877-878` 那条注释与 `:879` 那枚钉已漂**（注释写"四枚字段"、钉要求 `!= 5`、现量是 5）——本票三形都不碰它，但它是本腿扫到的既有缺陷。登不登台账（`A##`／`Q##` 由编排者写，⛔ 本腿不许动台账） | §2.2 第 6 行 | |
| 8 | **AC#4 的门禁谁执行**：本腿只能自证"`git status --porcelain -- cmd internal` 终态＝起手（0 行）"这一半；`d22scan`／`gofmt -l`／gofumpt `-l`／`go vet ./internal/tools/`／整包红名册作差 四件**都要求跑 go**，本腿一枚没跑。是否由写腿落地时同批发，还是编排者核过后代跑？ | §0 表末行＋§2.6＋§4 第 1 条 | |

---

## §4 判不动的地方（本腿这把尺结构性看不见什么——⛔ 不许推测填空）

1. **一切"会红几枚"的实际名册**。⛔ 零 `go` 命令 ⇒ §2.2 的"红／不红"全部是**读断言原文的静态推演**，不是实跑结果；票面 AC#0 要的"崩溃面还剩几枚零读数的**实测**"本腿**未交**（见 §3 第 3 条）。任何"甲落地后名册仍是 206"这类句子都不能从本件得出。
2. **`go:build` 标签／构建约束看不见的文件**。`grep --include=*.go` 读的是盘上字节，看不见构建标签、看不见 `GOOS=windows` 才编译的文件到底参不参与本次测试二进制。⇒ 本件所有"文件参与/不参与"的判断只到"它在 `cmd/`＋`internal/` 且是 `.go`"这一层。
3. **模块与包边界看不见全貌**。本腿没跑 `go list` ⇒ 不知道 `internal/tools` 测试二进制的**真实文件集**（含 `//go:build` 过滤后的），也不知道 `tools/d22scan` 这枚独立 module（`main.go` 文档注释自述 ∶ `this package is its OWN Go module (tools/d22scan/go.mod)`）在 CI 里被怎样调用（只 grep 到 `ci.yml:74`／`:81` 两行）。
4. **`newSubagentToolProvider(nil)` 会不会崩判不动**（`:293` 那一形）。要判它需要读那枚函数体＋实跑；本腿读了它的调用点，⛔ 未定性。⇒ §2.1"共同不覆盖的第 2 格"只陈述"这一维三形都不动"，不陈述"这一维会不会崩"。
5. **`d22scan` 的真实输出判不动**。§2.6 是**读源码文本**得出的形状判断，不是扫描结果：⛔ 没跑 `scripts/d22scan.sh`、没跑 `runtests.sh`、没读 `allowlist.txt` 全文（只读了 grep 命中的行）。ban 8 的 emojiRe／ban 9 的 `repoPathRe` 逐字串是原文，**但"三形的具体文案会不会命中"依赖最终写下的字节**，本腿无法预知。
6. **overlay 语义看不见**。票面 AC#1／AC#2 的台件靠 `go test -overlay`；`failclosed_236_teeth_test.go:68-75` 与 `tasklist_deferred_236r3_teeth_test.go:33-45` 都记着同一条仪器事实——**读盘型尺（`os.ReadFile`）对 overlay 换进去的字节结构性失明**。⇒ 本件 §2.2 第 8 行只敢判"`subagent_197.go` 没有读盘尺"，⛔ 不敢判"有读盘尺时谁红"。
7. **用例枚数口径看不见缩进子测试**。`236-v1` §0-补2 自己声明了 ∶ `206` 是 `grep -c '^--- PASS'` 的**顶层枚数**，行首锚定看不见子测试。⇒ 本件引用"206／181"时继承同一口径；**"派发点枚数 28"与"用例枚数"不是一回事**（§1.3 已具名说明，本件一律不换算）。
8. **他腿在飞的面看不见**。另有一枚腿在 `cmd/wisp` 跑突变 ⇒ 本件只锚 `git status --porcelain -- cmd internal`＝0 行这一把尺；⛔ 看不见那枚腿的中间态、也不引用它任何读数。§0 那份 `logs/gitstatus-full.txt` 是起手时刻的快照，**不代表交件时刻的全仓状态**。
9. **票 236／235／222／221 的历史判语不在射程**。本腿只按派单取了 `236-v1` 的**读数**（§2.4），⛔ 未采纳其 §2 第 6 行与 §5 第 2 条的定性句；其余件（含 `.scratch/wisp/probes/236/a1/census.md`、`235/**`、`docs/evidence/s1/*`）只统计了命中枚数，⛔ 未读内容。禁读名单（`probes/232/**`、`probes/pool-validity/5g2/**`、`probes/expired-premises/7b2/**`、`frontend/**`、`design/**`）**一条未碰**。
10. **"该不该拆语句"本身不在本腿射程**。§2.1/§2.2 只回答"拆了要花什么代价"，⛔ 不回答"值不值"；AC#2 要的形状（半边摘 ⇒ 恰 1 枚红）只有乙给得出，但那是**需求优先级**判断，属 §3 第 1 条的留空格。

