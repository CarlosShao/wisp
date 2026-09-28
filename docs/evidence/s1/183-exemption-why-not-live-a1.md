# 183-a1 — 只读设计核：票 177 那枚"宿主自己写下的路径精确豁免"为什么在真机 CLI 上没护住

- 程：`183-a1`（只读核，零跟踪改动）｜派单 `.scratch/wisp/dispatches/2026-09-28-114x-readonly-183-a1-…md`
- 起手时刻 `2026-09-28 11:53:41 +0800`｜分支 `dev`｜HEAD `e2dde3f0d9351539e773e30519ed6c5bfcc45727`
- ⚠ 派单写"我这边 `95218c4c`"：现量 HEAD 比它**多一枚 commit**（`e2dde3f0`＝两枚派单本身）。尺：`git merge-base --is-ancestor 95218c4c HEAD && echo yes` → yes。**不是分岔，是编排者自己把那两枚派单提交了**，本程按 HEAD 做。
- 写面：本文件＋`.scratch/wisp/probes/183/a1/**`＋票 183 Progress log（只追加，一枚框没勾）。

## 1. step-0 五件（全部现跑）

| 件 | 读数 |
|---|---|
| `date` | `2026-09-28 11:53:41 +0800` |
| `git rev-parse --abbrev-ref HEAD` | `dev` ✅ |
| `git rev-parse HEAD` | `e2dde3f0…`（派单号 `95218c4c` 是其父，见上） |
| `git status --porcelain -- internal/ cmd/` | **空**（起手） |
| `go test -count=1 -run 'ShapeA\|177' ./internal/risk/` | `ok github.com/CarlosShao/wisp/internal/risk 0.086s`（全绿，与票面一致） |

起手之后**立刻**发生的一件事（决定本程全部读数的可信度）：

```
$ git status --porcelain -- internal/ cmd/
 M internal/agent/loop.go        ← 在飞的 179-v1 的 M3 变异（把 PassThroughUnclassifiedRisk 判据反了）
```

它把我**第一次**跑的真机 CLI 变成另一发（`task.output → tool error: 风险未分级且直通开关关闭，已拒绝执行`，见 `probes/183/a1/logs/e2e-readings.txt`）。我从 `git show HEAD:internal/agent/loop.go` 抽了一份钉进自己的 overlay（`probes/183/a1/mut/loop.go`）后重跑，之后所有读数都在钉死的 HEAD 码上。尺：`git diff -- internal/agent/loop.go | grep -c "MUTATION 179-v1"`。
⚠ **通用坑（比 `runtime.Caller` 那枚更阴）**：`-overlay` 只保护你不脏别人的文件，**不保护你不吃别人脏的文件**——同一棵工作树里，任何一枚走 `go test` 的程都可能把另一枚的临时变异编进自己的二进制。要防这一手，把自己依赖的每枚跟踪件都钉进 overlay（我钉了 `loop.go` 一枚就够，因为它是唯一被脏的）。

## 2. 本程**没**测什么

- 没跑整包 `go test`、没跑 `gate-clauses.sh`、没跑 d22scan（派单 §0 禁令：179-v1 在飞 ⇒ 数不可归因）。所有 `go test` 都带 `-overlay` 且带 `-run` 定点过滤器。
- M-A（跨度用原始字节下标）只取到**常驻判据侧**的红腿名册，**没取到 CLI 侧读数**——那条命令漏了 `WISP183DBG=` 前缀，读数没落盘（`cut: …/m-A.txt: No such file` 就是证据），预算到顶未重跑。所以 M-A 的"两向"只有向。
- 没测 `dropped=1` 那枚线索的跨任务面：票 179 台件把 run A／run B 跑在**同一进程、同一 prov、不同 scope**，跨进程／重启后的名册与 mark 对应关系本程零读数。
- 没测"真背景输出里到底多常见路径同款 8-gram"这一维（§3(d) 只证了机制，没证发生率的分布）。
- 没动 AC#6 的 `窗口 0.0s`：本程读数里连卡都没下（续读那一发在 R4 就被拒了），归口仍是票 162 那一族。
- 没验"修法会不会洗戳"（那是 AC#3 反向判据的活），只把风险点具名写进 §6。

## 3. 四支候选逐支判定

判据的**共同现量**：把 `internal/risk/provenance.go` 换成一份**只加日志**的副本（`probes/183/a1/mut/provenance.go`，逐字节等价＋`dbg183()`），在真 CLI 台件（`probes/183/a1/zz183a1_e2e_test.go`，＝179-r2 那枚台件改名）上跑，读数落 `probes/183/a1/logs/`。`dbg183` 的路径取自环境变量 `WISP183DBG`，**没有**用 `runtime.Caller`。

复现（一次跑完三枚读数）：
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
L="$PWD/.scratch/wisp/probes/183/a1/logs"; L=$(cygpath -m "$L")
WISP183DBG="$L/mark-hit.txt" go test -count=1 -overlay=.scratch/wisp/probes/183/a1/overlay-instrumented.json \
  -run 'TestRereadHostPointerOnCLISeam183a1' ./cmd/wisp/
```

### (a) 跨度定位偏 —— **不成立**（但这一支确实"被钉着"，钉它的判据在真机形状上从没跑过）

现量（`logs/mark-hit.txt` 第 1 行）：
```
MARK scope=f9ce0ecd-… tool=task.output origin="" contentBytes=3204 normRunes=3172 truncated=false
 minChars=8 maxSrc=262144 hostPathBytes=155 hostPathNormRunes=155 lo=2101 skip=[[2101 2256]]
 windows=3051 excluded=162 hashed=136 rawByteIndex=2244
```
⇒ 归一化坐标 `lo=2101`，原始字节下标 `rawByteIndex=2244`，**两者差 143**；码用的是前者（`provenance.go:530` 的 `runeIndexOf`）。跨度长度 `2256-2101=155` ＝ `hostPathNormRunes`，`excluded=162 = 155+8-1` 正是"重叠即排除"的语义。**窗口坐标没有偏。**
变异兑现（M-A：`runeIndexOf(rp, nhp)` → `strings.Index(content, hostPath)`，即派单问的"改成原始下标"）：常驻判据里 **5 枚**响，逐名：
```
--- FAIL: TestShapeAHostReReadOfDeclaredPathIsClean
--- FAIL: TestShapeAW2ForeignMarkCarryingSamePathStillHits
--- FAIL: TestShapeAW3NeverMintedSiblingUnderSameDirStillHits
--- FAIL: TestShapeAFExclusionDiesWithScopeClose
--- FAIL: TestShapeAOrderOneSpanLocatedAfterNormalization
```
尺：`go test -count=1 -overlay=… -run 'ShapeA|177' ./internal/risk/ 2>&1 | grep -aE -- "--- FAIL"`。
⇒ 结论：**如果**这一支是病，票 177 的判据会红；今天它们全绿且 `lo` 读数自洽 ⇒ 这一支清白。（`maxSrc=262144`、`truncated=false` 也顺手排掉了"路径落在截断尾之外 ⇒ 排除不到"这一形——它在真机形状上远够不着。）

### (b) 载具没送到 —— **不成立**（我的派单人最怀疑这一支，它错了）

现量：上面那行 `hostPathBytes=155`（**非空**）＋ `lo=2101`＋`skip=[[2101,2256]]`＋`excluded=162` ⇒ `mark` 那一刻 `hostPathBox` 里有值、值就是名册登记那条产物路径。
名册回填与当场 `box.set` **不是两条对立的路、也没有先后问题**：`Backfill`（`internal/tools/task_backfill.go:112` `rec.ArtifactPath = sp.Path` → `:115` `Roster.Record`）只写**数据**；`box.set`（`internal/tools/task.go:253-255`，号现取：`grep -n "box.set" internal/tools/task.go`）在 `task.output` 被模型调起的那一刻从**同一条记录**读出来。`Bridge.run` 先 `withHostPathBox`（`internal/tools/bridge.go:459`）、后 `Execute`（`:470`）、再 `b.mark(dec, res, hostPaths.get())`（`:535`）——同一次调用、同一枚 box、顺序正确。
变异兑现（M-B：把 `task.go` 那句 `box.set` 摘掉，走 `probes/183/a1/mut/task.go`＋overlay）：
```
MARK … hostPathBytes=0 hostPathNormRunes=0 lo=-1 skip=[] …
MATCH … paramRunes=155 … frag="c:\\users" hit=true
```
且桥级那五枚里 **2 枚红**（尺：`go test -count=1 -overlay=… -v -run '175r2' ./internal/tools/ 2>&1 | grep -aE -- "--- (FAIL|PASS)"`）：
```
--- FAIL: TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2
--- FAIL: TestHostMintedPointerRereadStaysCleanOnRealBridge175r2
```
⇒ 这一支**有**判据钉着（175-r2 那对桥级腿），而且摘掉载具时命中的是路径**开头** `c:\users`；今天真机上命中的**不是**开头（见 (c)/(d)）⇒ 载具送到了。派单的"回填与当场可能不同一条"这一前提**证伪**。

### (c) 命中的不是路径窗口而是那句中文说明 —— **不成立**

现量：base 那发命中的片段逐字是 `frag="-output-"`（`logs/mark-hit.txt` 第 2 行）。它是产物路径 `…\artifacts\tool-output-agent-task-<uuid>.txt` 里的一段，**不是** `全文见 ` 也不是 `省略 17200 字符`（那几句归一化后是 CJK＋数字，8-rune 窗口与 `-output-` 无交集）。M-B 那发的 `frag="c:\\users"` 同样是路径本体。⇒ 拆"路径子串／非路径子串"各造变异这一步省了：**命中一直落在路径本体上**，中文说明无嫌疑。

### (d) 索引之外还有第二条比对路 —— **成立，且这一支就是根因**

严格说不是"第二条路"，是**第一条路的两处位置盲**，票 177 §4.1 那句"`f.src` 复核面不动 ⇒ 不会产生假命中"在真机形状上**不成立**（本程现量，未照抄）：

1. **索引侧按位置排除，但哈希集是去重的、位置无关的**（`taintmatch.go:134-146` 排序＋去重，只留 `uint64`）。被豁免窗口 `"-output-"` 的哈希**因为同一个 8-rune 串在标记正文的另一处（非豁免区）出现过**而照样躺在集合里。base 读数就是它的量：`windows=3051 excluded=162 hashed=136` —— 3051-162＝2889 个窗口压成 136 枚哈希，压缩比正是"同一串多处出现"的直接体现（台件的背景正文是 `strings.Repeat("WISP179R2-background-output-line.", 700)`，其中 `-output-` 与路径里的 `tool-output-agent` 同款）。
2. **复核侧整串比对，也按位置无关**（`taintmatch.go:207` `strings.Contains(f.src, w)`，`f.src` 是**整段归一化正文**，豁免窗口没从里面挖掉）。于是"哈希来自别处 × Contains 命中别处"＝一枚本该被豁免的窗口被**重新捞回**，R4 升 L2、当场拒。

**决定性变异 M-D1**（唯一变量＝标记正文，把背景正文换成纯 CJK、与路径零同款 8-gram；改的是我那份台件副本里的 `long183a1`）：
```
MARK … hostPathBytes=154 lo=767 skip=[[767 921]] windows=1182 excluded=161 hashed=127
MATCH … paramRunes=154 … frag="" hit=false
→ [audit] … 无 R4；整份 fs.read 成功；「整份读回逐字节相等：20040 == 20040」
```
尺：`grep -a "hit=false\|逐字节" .scratch/wisp/probes/183/a1/logs/mark-hit-D1.txt`；台件断言读数见同目录 `gotest-D1.txt`。
⇒ **同一枚豁免、同一份生产码、同一个真 CLI 接缝**，只把"正文里有没有路径同款片段"这一维换掉，那一发就从"R4 拒"变成"逐字节读回"。四支里只有 (d) 有这种一翻两效的形状，**判定：根因＝(d)**。
（`gotest-D1.txt` 里剩下的 `读回 10002 vs 10000` 是我这枚 CJK 台件自己踩在 rune 中间被 `max_bytes` 切开的边界伪影——顺带说明 `fs.read` 不切半枚 rune，与本题无关。）

## 4. 还原证明（每次变异之后各跑一次，全部为空）

```
$ git status --porcelain -- internal/ cmd/        # 起手
$ git status --porcelain -- internal/ cmd/        # M-D1 之后
$ git status --porcelain -- internal/ cmd/        # M-A＋MUT-B 之后（两次）
$ git status --porcelain -- internal/ cmd/        # 本程最后（§1 表格同一条尺）
（每条命令输出均为空 —— 上面四个时点的实测都是空，逐次贴在本程会话里）
```
唯一"脏"过 `internal/` 的是 **179-v1 的 `M internal/agent/loop.go`**，不是我；我在 §1 具名报了，并用 overlay 钉回 HEAD 版本自保。本程结束时 `internal/ cmd/` 仍为空。

## 5. 派单 §2 三格

**① 为什么八枚常驻判据全绿却挡不住这一发——逐枚答接缝。** 尺（枚数与函数名，现跑）：`grep -n "^func Test" internal/risk/shape_a_exemption_test.go`。八枚**全部长在 `internal/risk` 包级、直接调 `p.MarkWithHostPath(...)` 那一层**（`shape_a_exemption_test.go:61/74/90/131/154/172/185/206/213/222/240` 逐枚可见），桥级一枚都没有、CLI 级零枚：

| 判据（名字取自上面那把尺） | 它测的接缝 | 为什么拦不住 (d) |
|---|---|---|
| `TestShapeAHostReReadOfDeclaredPathIsClean` | 包级：stub→mark→Inspect | `shapeAStub` 正文里那条路径**只出现一次**、且正文没有 `tool-output-` 同款片段 ⇒ 豁免窗口没有"别处的孪生哈希" |
| `…W2ForeignMarkCarryingSamePathStillHits` | 包级：两枚 mark | 它只测"外来**那枚**命中"，外来那枚自带全文索引，与豁免窗口的去重回灌无关 |
| `…W3NeverMintedSiblingUnderSameDirStillHits` | 包级：同目录兄弟路径 | 同上：测的是"不同路径不误豁免"，不是"同一路径的同款子串" |
| `…F…ExclusionDiesWithScopeClose` | 包级：scope 生命周期 | 只跟 `Close` 有关 |
| `…OrderOneSpanLocatedAfterNormalization` | 包级：ordering (i) | 钉的是坐标（M-A 时它红，说明它干活），与去重后回灌无关 |
| 其余三枚（absent / empty / 全串即路径 那族，`:185/:206/:213/:222`） | 包级：豁免失败形与不吞豁免 | 都不喂"正文重复"这一形 |

⇒ 这不是"测试覆盖不足"这种废话，而是一句具体的：**豁免的强度是按"位置"定义的，索引与复核的强度是按"字符串集合"定义的，两者之间缺一条"同款片段的第二处出现"判据。** 要长在哪一枚接缝上才拦得住：**最省的是 `internal/risk` 包级再加一枚形状**（同一枚 mark 的正文里，路径同款片段在豁免窗口之外**再出现一次** ⇒ 断言今天会命中、修法之后不命中）——它今天就能红（M-D1 的反向即它的现量），不必等 CLI。CLI 级那枚（真机 `-overlay` 台件，本程已经现成可抄：`probes/183/a1/**` 三件套）该长成 **AC#2 的正向判据**，因为票 183 AC#2 明写"必须长在 CLI 接缝形状上"；但它跑一次要 60 秒且要吃别人脏件，**不能替代**包级那枚便宜红腿。桥级（175-r2 那对腿）已被证明钉得住载具（M-B 两枚红）、钉不住 (d)，别指望它。

**② `dropped=1` 那一格。** 尺：`grep -an "dropped=" .scratch/wisp/probes/179/r2/logs/e2e-readings.txt` → run A `dropped=0`、run B `dropped=1`。语义现取：`dropped := len(b.prov.ScopeTaints(taskID))`（`internal/tools/bridge.go:754`，尺 `grep -n "dropped :=" internal/tools/bridge.go`）＝**关闭时该 scope 里有几枚 mark**。⇒ run B 的 scope 里只有它**自己**那次 `task.output` 的 mark（本程直接看到：`MATCH … marks=1 markTool=task.output`）。裁一句：**豁免要作用的是 run B 当场刚产出的那一枚 mark**，run A 的那枚既不在场也不是判据（run A 的 scope 是另一枚 id）。跨任务／跨 scope 时名册与 mark 的对应关系：**名册（`ArtifactPath`）是进程级共享的，mark 是 per-task-scope 的**——A 任务落的产物路径被 B 任务的 `task.output` 打印出来 ⇒ 豁免只落在 B 的那枚 mark 上；外来内容若出现在**另一个 scope 的**某枚 mark 里，两枚 mark 各自独立索引，W-2 那一族仍然各判各的。本程**没有**跨进程／重启后的读数（§2 已具名）。

**③ 行号重锚（本程现跑，全部 `grep -n`）**：`internal/tools/task.go:253-255`＝`hostPathBoxFromCtx`＋`box.set`；`internal/tools/task_backfill.go:112`＝`rec.ArtifactPath = sp.Path`、`:115`＝`Roster.Record`；`internal/tools/bridge.go:459/470/535/575`＝`withHostPathBox`／`Execute`／`b.mark`／`MarkWithHostPath`；`internal/risk/provenance.go:530`＝`runeIndexOf`、`:541`＝`newFragmentIndex(string(rp), p.minChars, skip...)`；`internal/risk/taintmatch.go:134-146`＝排序去重、`:207`＝`strings.Contains(f.src, w)`；`internal/agent/loop.go:801`＝"风险未分级"那一句。尺：`grep -n "TaskDeps{" cmd/wisp/run.go`（派单说 362→365 那一处，本程没引它）。⚠ 派单 §1(d) 引的 `provenance.go:662` 在 HEAD 上是 `var cands []cand` 那一行——`662` 这一枚**号已漂**，现读"路径参数从不豁免"的语义在 `Inspect` 的 `if !gateOpen && !isPathKey(k) { continue }`（HEAD 上 `provenance.go:709`，尺 `grep -n "isPathKey(k)" internal/risk/provenance.go`）。

## 6. 推荐落地点（行，不是码）＋要不要动冻结件

**根因一句话**：豁免只作用于"建索引时哪些窗口进哈希集"，而"命中判定"用的是**去重后的哈希集**＋**对整段正文的 `strings.Contains`**，两者都不带位置；于是豁免窗口内的那 8-rune 串只要在正文别处出现过一次，豁免就当场失效。

**推荐落地点（一处，最小）**：`internal/risk/provenance.go:541` 那一句——**喂给 `newFragmentIndex` 的那段文本，先把被豁免的跨度物理挖掉**（`rp[:lo] + rp[hi:]`），使**哈希集与复核语料 `f.src` 同时**不含豁免窗口，位置性自然一致。它同时满足：不需要新字段（不会把对抗验收 M-5 那笔"写而不读的窗口拼写表"请回来）、不引入按值放行（挖的还是"这一枚 mark 的跨度"）、`taintmatch.go` 一字不动。
- 已知代价（必须写进落地判据）：挖掉跨度会在切口处**新造**跨界的 8-rune 窗口（例：路径前的 `…总长` 与路径后的 `…]\n` 被缝到一起）。方向是**多判**（fail-closed）＝安全方向，但**要一枚常驻判据钉住它有多严重**，别让它变成莫名多出一张 L2 卡。
- 副作用面：`marks=2` 那一格（`fs.read` 的正文随后也被盖戳，见 `logs/mark-hit-D1.txt` 第 3 行 `MARK … tool=fs.read origin="C:\…txt"`）——同一份路径文本在**另一枚 mark 里是正文不是跨度**，所以 AC#3 反向那一发"外来内容出现同一条路径 ⇒ 仍要命中 R4"必须现量复测，别假设挖跨度会误放它。

**要不要动冻结件（具名报回，本程一字未动）**：
- **不必**改 `provenance.go:468-473`（`Mark` 的契约文字）与 `taintmatch.go:11-15`（"逐 token 追踪已被否决"）这两句**文字**。
- ⚠ 但**要人工拍一枚**：`MarkWithHostPath` 的文档块（`provenance.go:482-506`，含 ordering (i)/(ii) 那段）里那句"only the normalized window that hostPath occupies inside THIS mark is kept out of the fragment index"在挖跨度修法下**语义会变强**（从"窗口不进索引"变成"这段文本从索引与复核语料里一起消失"）。它是**契约性描述文字**、按 AGENTS §0.2 属"改契约＝人工批准"的面 ⇒ **`183-r1` 派单前请编排者先定：允许追加新注释、是否允许改这段描述文字**（改 `PLAN.md`／`thresholds.go`／golden／审批常量／`allowlist.txt` 一律不碰）。
- 三支毒修法（AC#5①②③）本程一枚都不需要：不用放宽 R4、不设 `PassThroughUnclassifiedRisk`、不把 `fs.read` 摘出名册。

## 7. 被拒／没成功的调用

**取数前**：无（没有一次工具调用被权限拒）。
**取数后（读数作废／没拿全，逐条）**：
1. 第一次真机 CLI 跑（`logs/e2e-readings.txt`）被 179-v1 的脏 `loop.go` 污染 ⇒ 那发不是本题那一发，已重跑并具名（§1）。
2. M-A 的 CLI 读数**没落盘**：命令漏 `WISP183DBG=` 前缀（`cut: …m-A.txt: No such file`）⇒ (a) 支只有判据侧变异。
3. 第一次 M-B **整枚作废**：`overlay-instrumented.json` 里没有 `internal/tools/task.go` 这一枚映射，所以那次"五枚桥级腿全绿"是**变异没生效**造成的假读数（我一度据此写"载具没判据钉着"——**那句话错，以 §3(b) 修正版为准**）。补上映射后重跑：两枚红、`hostPathBytes=0`。
4. `grep` 两次因未转义括号／`--- FAIL` 被当选项而失败（`grep: Unmatched (`、`unknown option`），换 `grep -aE` 后拿到数。

## 8. 有没有跑过删除命令

跑过 `rm -f` **两枚、作用面只有我自己的探针读数文件**（`probes/183/a1/logs/mark-hit.txt`、`…/m-A.txt`、`…/m-B.txt`，为的是让 append 型 sink 从空开始）。`probes/**` 既有台件、`internal/`、`cmd/`、`docs/` 一律未删；未跑 `git clean`／`restore`／`checkout .`／`reset`；临时件只建不删（作废的两份读数仍在盘上）。

## 9. 工具调用：用了 **43** 次 vs 硬顶 **45**（自报枚数：第 43 枚＝把这份表与探针提交进来的那一步；到顶前停手，未追加探索）。

## 10. 伪授权两栏

- **本程收到的、判为"真授权"**：派单文件＋票 183 票面（写面清单、硬顶、禁令）——出处 `.scratch/wisp/dispatches/2026-09-28-114x-readonly-183-a1-…md` §3、`issues/183-….md` AC#1/AC#7。
- **本程遇到的、判为"不是授权"（没照它做，具名）**：① 工作树里 `internal/agent/loop.go` 的 `MUTATION 179-v1 M3 inverted` 注释——那是**另一枚程的临时变异**，不是"HEAD 已改成这样"的授权，也没人叫我在本题里接受它 ⇒ 我没跟它交互、直接 overlay 钉回 HEAD。② 派单里"我最怀疑第二支（载具没送到）"与"回填与当场可能不同一条"——编排者的怀疑**不是**结论授权；现量把它证伪了（§3(b)），我照实报回。③ 票 183 面上 `AC#1` 之外七格全空 ⇒ 一枚没勾。④ `docs/PLAN.md`／`provenance.go:468-473`／`taintmatch.go:11-15` 的冻结文字：谁都没授权我动，我也没动。

## 11. 凭据值零抄录

本程未接触任何 API 密钥／DPAPI 明文；读数里的路径只有 `$TEMP` 下的测试产物与仓内探针目录；`sha256` 值是**测试正文**的哈希（台件自己打的），非凭据。

## 12. next=

**`183-r1` 落地腿派之前还缺什么**——三件，缺一件就会再产一枚"做了但没生效"：

1. **缺一枚包级红腿的形状定义**：派单必须把"豁免窗口同款片段在正文别处再出现一次"写成 AC#2 的**反向孪生判据**（今天它红、修法后它绿；现成可抄的现量就是 `logs/mark-hit.txt` 的 `frag="-output-" hit=true` 与 `logs/mark-hit-D1.txt` 的 `frag="" hit=false`）。不写它，AC#2 那枚 CLII 级判据一旦因"背景正文刚好不含路径同款片段"而假绿，本票等于没修。
2. **缺编排者一枚拍板**：`MarkWithHostPath` 文档块（`provenance.go:482-506`）那句"kept out of the fragment index"的语义要不要随"挖跨度"改写（§6）；以及切口新造窗口（多判方向）算不算可接受残留。这属"改契约＝人工批准"面，agent 不该自己定。
3. **缺一条排程护栏**（本程实测踩到）：`183-r1` 若与任何一枚改 `internal/agent/**`、`internal/tools/task.go`、`internal/risk/**` 的程同时在飞，它必须把自己依赖的**每枚跟踪件**钉进 overlay（`git show HEAD:path` 抽副本），或者两枚程错开；否则"绿"与"红"都可能不是它的码。**建议：`183-r1` 之前先收掉在飞的 `179-v1`。**

其余：`177 AC#3` 条件那一格、`175 AC#5`、`176 AC#3-5` **维持不翻**（要等 `183-r1` 端到端）；票 183 的 AC#2..AC#8 本程一枚没做、一枚没勾。
