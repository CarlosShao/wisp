# 154 实现程 r1 — 「关」只长在环路那一枚 id 上：宿主自带 id 那一形仍没人关、并发任务零读数

派单＝`.scratch/wisp/dispatches/2026-09-26-115x-impl-154.md`
工单＝`.scratch/wisp/issues/154-the-close-only-covers-the-loop-task-id-so-host-supplied-task-ids-still-have-no-owner-and-concurrent-tasks-have-zero-readings-reserved-with-a-trigger-gate.md`
本件是**实现程**自己那一遍；票面框由编排者按**非实现者**验收表定，本程不自勾。
本程写面（派单"git 纪律与写面"一节）＝`internal/tools/bridge.go` 的**注释面** ＋ 本件 ＋ `.scratch/wisp/probes/154/`。

---

## 第 0 节 锚点、工具链、读数纪律

**共享树在取数期间自己往前走**（这是本票第 0 节必须先写的，否则下一位对不上行号）：

| 时刻 | `git rev-parse HEAD` | 本程用它量过什么 |
|---|---|---|
| 进场 | `317805c` | `git log`／`git status`（`git cat-file -t` 见下） |
| 进场第二读 | `d0865ff` | 两枚包"改前"门禁跑的就是这一枚工作树 |
| AC#1 三枚子句现读 | `e95caf1` → `bcf7601` | 同一把尺两枚锚点各跑一次，输出逐字一致（差的那几枚 commit 只碰 `frontend/**`、`docs/evidence/s1/152-…`、`docs/reports/pending-and-issues.md`、`.scratch/wisp/issues/154-…` 作废牌改名） |

⇒ **`internal/**` 与 `cmd/**` 在 `d0865ff..HEAD` 之间零改动**，现读命令与输出：

```
$ git diff --name-status d0865ff..HEAD | grep -E "^.\s+(internal|cmd)/" ; echo "rc=$?"
rc=1                     ← 零命中（尺＝git diff --name-status 全量，正则筛路径前缀 internal/ 与 cmd/）
$ git status --short -- internal cmd     →  空（改前，工作树＝HEAD）
$ git diff --cached --name-only          →  空（进场时 index 干净：上一枚 154 作废牌的删除已由 dc2046fc 收掉）
```

工具链现读（不背别人的数）：

```
$ go version            →  go version go1.27.1 windows/amd64
$ git cat-file -t 5d46f24  →  commit        （票 151 那一枚，本件引它之前先验过类型）
$ git log -1 --format="%h %ad %s" 5d46f24  →  5d46f24 Sat Sep 26 09:44:57 2026 +0800 placeholder
      ⚠ 它的 subject 就是字面 "placeholder"（本程不改它、也不据它判任何事，只是提醒下一位别以为抄错了 sha）。
$ ls third_party/sherpa-onnx/*.dll | wc -l  →  3（onnxruntime / sherpa-onnx-c-api / sherpa-onnx-cxx-api）
      且 `git ls-files third_party | wc -l` = **0** ⇒ dll 是**未跟踪**的本机件：纯净树里没有它们。
$ git config core.autocrlf  →  true   （＋ .gitattributes 的 `* text=auto` ⇒ 本程一律不用 `git archive | tar -x` 取版）
```

**本程一律不用 `git archive | tar -x` 取版**（派单坑①）：门禁两跑跑的是**当前工作树**（改前工作树在 `internal/**`、`cmd/**` 上与 `d0865ff` 无差，已由上面那条 `git diff --name-status` 钉住），
"改后"跑的是**我改过之后的同一枚工作树**——本票的改动只有注释，因此结构上碰不到 `-overlay` 与 `-cover*` 同用那枚坑（本程全程没开 `-cover`）。

取版坑的**已知会红的正控**本程照派单要求真跑了一发（不跑正控就宣称"字节不等"＝本仓既有条）。原文与全部读数＝
`.scratch/wisp/probes/154/archive-eol-control.txt`，摘要：

```
$ git ls-files --eol | awk '$1=="i/lf" && $2=="w/crlf"' | wc -l      →  89   ← 索引形(LF) ≠ 检出形(CRLF) 的 tracked 件
    其中 attr 分布：80 枚 attr/text=auto（无 eol 钉）＋ 9 枚 attr/text eol=lf
    那 80 枚按扩展名：53 log · 20 txt · 2 tsx · ts/toml/json/js/css 各 1
$ git cat-file -s <A>:<path>  vs  wc -c < <path>（未修改件）
    .scratch/wisp/probes/149/combos2/x-p1-d11a-d11b.log  blob=3391 worktree=3436   ← 派单点名的那枚实发，逐字复到
    .scratch/wisp/probes/149/combos2/x-p1-d11a.log       blob=3567 worktree=3614
$ git archive <A> docs/reports/injection-timeline.md | tar -x -C /d/tmp/wisp154-archive && cmp … → SAME（46117 = 46117）
```

⇒ 三条结论，都带方向：① **不等是真的**，且量级＝每枚约「行数」个字节（LF→CRLF）；② 射程**不覆盖 `*.go`／`*.md`／`*.sh`／`*.sse`**
（`.gitattributes` 给它们钉了 `eol=lf` ⇒ 对这些件 `git archive | tar -x` 是等字节的）；③ 会被咬住的是那 80 枚 `text=auto` 件，
落点是 **73 枚 `.scratch/wisp/probes/**`（上一程留下的原始读数，含那 53 枚 `.log` ＋ 20 枚 `.txt`）＋ 7 枚真码**：
`deps.toml` · `design/doubao/demo/app.js` · `design/doubao/demo/styles.css` · `docs/evidence/s1/66/66-full-subset-slo-report.json` · `frontend/scripts/render-composer.tsx` · `frontend/src/components/composer.tsx` · `frontend/src/lib/panel.ts`
现读命令＝`git ls-files --eol | awk '$1=="i/lf" && $2=="w/crlf" && $3=="attr/text=auto" {print $4}' | grep -v '^\.scratch/wisp/probes/'`）。
⇒ 谁要复算上一程的日志，`git archive` 拿到的不是同一串字节；`frontend/**`／`design/**` 那几枚此刻正被别的会话未提交地改，**本程不碰、也不把它们算进任何零命中宣称**。
**本程一律按派单给的替代法**：Go 侧名册/字节比对用 `git ls-tree -r` ＋ `git cat-file --batch`／锚点行 grep，不 archive。

**读数纪律**：凡引用 sha 先 `git cat-file -t`；凡点名符号先 `git grep -w` 它在不在；行号全部按本程自己锚点现读；"零命中"一律写是哪把尺。
票面上那**三枚空号**（`OpenTask("")`／`RunTextTask`／`unboundScopes()`）本件**不出现**为"存在过的东西"；`RunTextTask` 的更正（票面 Progress log 10:3x 那条）本程照抄为：**`Loop` 上没有收外部 task id 的方法**，
真实存在的是 `cmd/wisp/run.go` 里 CLI 自己那一腿 `runTextTask`（现读 `grep -n "func runTextTask" cmd/wisp/run.go` ⇒ `:137`）。

---

## 第 1 节 格 AC#5（改前两跑）— 四数 ＋ 名册基准

派单/票面口径：**逐包单跑**、四数之外名册两向 `comm`。四数的尺先钉死（全只认行首，避开"汇总行与 `-skip` 散文造假命中"）：
`^=== RUN`＝RUN；`^--- (PASS|FAIL|SKIP)`＝顶格；`^[[:space:]]*--- (PASS|FAIL|SKIP)`＝全量；`^panic:`＝panic；红绿只认 `^--- FAIL:`。

```
$ export PATH="$PWD/third_party/sherpa-onnx:$PATH"      # 不带它就是加载期 0xc0000135＋0 条 RUN＝没跑到
$ go test -count=1 -v ./internal/tools/   > probes/154/gate-pre-tools-v.txt   rc=0  ok 20.288s
$ go test -count=1 -v ./cmd/wisp/         > probes/154/gate-pre-cli-v.txt      rc=0  ok 304.198s
```

| 包 | RUN | 顶格 | FAIL | SKIP | panic | 全量裁决 | 名册文件 |
|---|---|---|---|---|---|---|---|
| `internal/tools` | **115** | **79** | **0** | **0** | 0 | 115 | `probes/154/names-pre-tools.txt`（79 行） |
| `cmd/wisp` | **139** | **79** | **0** | **0** | 0 | 139 | `probes/154/names-pre-cli.txt`（79 行） |

名册基准对拍（与 151 那两程同尺，为了证明"改前"不是本程自己污染的树）。
⚠ `comm` 只吃**排过序**的输入，本程第一发忘了给右栏排序、读出来一枚"151 程独有的用例"，方向是反的；
两栏都 `sort` 之后重跑才是下面这个形状（`sort -c` 对两栏皆 rc=0）。尺＝`^[[:space:]]*--- (PASS|FAIL|SKIP)` 取第一列、去掉子用例（含 `/` 的行）、`sort -u`。

```
$ comm -23 mine-cli.txt theirs-cli-sorted.txt      ← 左栏独有（＝本程 pre 多的那枚）
TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne
$ comm -13 mine-cli.txt theirs-cli-sorted.txt      ← 右栏独有（＝151 程多的那枚）
(空)
$ comm -23 mine-tools.txt theirs-tools-sorted.txt  →  空
$ comm -13 mine-tools.txt theirs-tools-sorted.txt  →  空
（台件：probes/154/names-pre-{cli,tools}.txt ↔ probes/151-accept/{my-post-cli,theirs-post-tools-top}.txt，右栏另存一份 sorted）
⇒ `internal/tools` 名册与 151 验收程在它锚点下的 79 枚**逐字相同**；
⇒ `cmd/wisp` 差值只有**一枚**：`TestSLO152Corrupt…`，那是票 152 那条腿新加的用例（本票在 `cmd/wisp` 不新增、不删除任何用例）。
```

⇒ **改前基线成立**：两包 rc=0、FAIL/SKIP/panic 全 0；本票"改后"要复的就是这四数＋名册两向 `comm`（改后读数见第 3.2 节，随注释改动一起交）。

---

## 第 2 节 格 AC#1 — 触发门：先正面回答"这条判不判得出"

票面给的候选形状（逐字）：**"全仓任何非测试代码里出现『构造 `tools.Request` 且 `TaskID` 非空』或『直接调 `Bridge.Execute`』，而其调用者不在 `Loop.run` 那条边界内"**，
并附一句"先判这条判不判得出"。本程按三小节答：先把形状里**点错的名字**改对（不改就永远零命中，那是假门的一种），再答判得出/判不出各是哪半，最后给出本程真正采用的门（第 2.3 节）。

### 2.1 形状里的两处符号，现量

```
$ git grep -nw "tools.Request" HEAD -- '*.go' ':!.scratch'      →  0 命中（尺：整词、大小写敏感、排除 .scratch）
$ grep -n "type ToolRequest struct" internal/agent/tools.go     →  :54
$ sed -n '54,64p' internal/agent/tools.go                       →  ToolRequest{ TaskID, CorrelationID, CallID, Name, Args, Timeout }（TaskID 是**导出字段**，没有 setter）
$ git grep -nE "\.Execute\(" HEAD -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go'
HEAD:internal/agent/loop.go:730:		return l.opt.Tools.Execute(ctx, req)
HEAD:internal/agent/loop.go:734:	out, err := l.opt.Tools.Execute(tctx, req)
HEAD:internal/tools/bridge.go:463:	res, err := entry.Tool.Execute(ectx, req.Args, onUpdate)
```

⇒ ①**没有 `tools.Request` 这个类型**：桥收的是 `agent.ToolRequest`（`Bridge.Execute` 的签名＝`internal/tools/bridge.go:243`）。照票面原话写的检子句会**结构上永远零命中**——那正是票面禁止的假门形状，所以本程的门改用真名。
⇒ ②"直接调 `Bridge.Execute`"这一支也不能按字面写：生产里对桥的调用只有 `loop.go` 那两行，且它们调的是**接口** `ToolExecutor.Execute`（`internal/agent/tools.go:91`），
`bridge.go:463` 那一发调的是 `Tool.Execute`（C4 单个工具，不是桥）。**按符号名 grep 会同时漏掉与误报**，故本程把子句落在"谁构造派发"上（理由见 2.2 末）。

### 2.2 判得出／判不出，各是哪半

| 半条 | 判得出否 | 依据 |
|---|---|---|
| "**谁在生产码里构造了过桥的派发**" | **判得出**（一条命令、rc 可判） | `ToolRequest.TaskID` 是导出字段、没有 setter、也没有第二个构造函数 ⇒ 要带着身份过桥，**必须**在某个文件里出现 `ToolRequest{` 这个字面量。它是唯一的语法位置，所以"名册里有没有这一行"＝"这一形存不存在"，不需要数据流。 |
| "…且 **`TaskID` 非空**" | **判不出** | 空不空是值域问题：`TaskID: m["id"]`、`TaskID: strings.TrimSpace(x)` 这类写法 grep 给不出答案；要答它得跑数据流（本仓没有这类仪器，见下）。 |
| "…而其**调用者不在 `Loop.run` 那条边界内**" | **判不出**（按函数体语法判得出，按"边界"判不出） | 行/函数级是语法事实，可以定位；但"这个函数是不是只在 `Loop.run` 之下被调到"是**跨函数可达性**，要 call graph。本仓现量：`grep -n "golang.org/x/tools" go.mod` ⇒ **0 命中**（仅 `go.sum:37` 有 v0.48.0 的条目＝不是本模块的 require），`tools/d22scan/main.go:98-99` 只 import `go/ast`＋`go/parser`（单文件解析，没有类型信息、没有 callgraph）。⇒ 上这枚仪器要**新增依赖＋改 `tools/d22scan/**`**，两者都在本票零字节面里，**本程不上**。 |
| "**最接近的近似命令**是什么" | 就是第 2.3 节那五枚子句 | 近似性在两处，都写死在门本体上：**(a)** 它按**文件/目录**排 `internal/agent/*`（不是按"边界"排）⇒ 若有人在 `internal/agent/` 里新增一枚宿主派发，本门**不响**（漏报方向）；**(b)** 它不看值域 ⇒ 若新构造点其实填的是空串，本门**响**而无害（误报方向）。⇒ 所以门的输出是**一枚要人读两行的名册**，不是一句 pass/fail；本程不把它伪装成后者。 |

### 2.3 本程采用的门＝五枚子句，一条命令一次跑完

生成器＝`.scratch/wisp/probes/154/gate-clauses.sh`（`bash .scratch/wisp/probes/154/gate-clauses.sh <锚点>` 即可复算）；
今日逐字输出＝`.scratch/wisp/probes/154/gate-clauses-roster.txt`（锚点 `1424aa7`，本件第 1 节那枚 commit）。

| 子句 | 命令（`git grep -nE` 一把尺，全部排除 `_test.go`） | 今日读数（锚点 `1424aa7`） | 响了说明什么 |
|---|---|---|---|
| **G1** | `'ToolRequest\{' -- internal cmd ':!*_test.go' ':!internal/agent/*'` | **空（rc=1）** | 有生产代码在自己的函数里造派发＝**宿主自带 task id 那一形出现了**。本票 AC#1②那串事必须当时做。 |
| **G1b** | `'\.Execute\(' -- internal cmd ':!*_test.go' ':!internal/agent/loop.go' ':!internal/tools/bridge.go'` | **空（rc=1）** | 有人绕过 `ToolExecutor` 接口的既有两枚调用点，直接对桥/对派发动手。与 G1 任一枚响即同判。 |
| **G2** | `-w 'OpenTask\|CloseTask' -- internal cmd ':!*_test.go' ':!internal/tools/bridge.go'` | **2 行**：`cmd/wisp/run.go:546`（注释）＋ `:563`（唯一生产调用） | **关那侧的所有者从 1 枚变多枚**，或**开那侧多出一枚生产调用者**。名册是"1 枚调用者"这一事实本身；它变长＝本票记的那一形已被人接上了（该回头核它接得对不对、并发读数补没补）。 |
| **G3** | `'^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' -- internal/agent` | **空（rc=1）** | **`Q-56` 那一支被落地了**（`Loop` 上出现收外部 task id 的导出方法）。这一枚响＝G1 那一形从"没有入口"变成"有入口"，关闭的所有者必须同批改。 |
| **G4** | `'PassThroughUnclassifiedRisk:' -- '*.go' ':!*_test.go'` | **空（rc=1）** | 生产组合根打开了 L0 直通 ⇒ **环路那一形端到端可见**（第 4 节"第二道门"的那把锁被摘）。这一枚不影响"宿主 id"，但它一响，"151 已把环路那一形接完"这句话才第一次有端到端读数。 |

**每枚子句都过了正控**（不是"看起来有牙"，是同一把尺打在仓库外的合成树上真响）：
脚本＝`.scratch/wisp/probes/154/gate-positive-control.sh`，第一发读数＝`gate-positive-control.txt`（**G3 那发是空的**，原因写在下），
修好后的第二发＝`gate-positive-control-r2.txt`：

```
G1  合成树 cmd/wisp/host_probe.go:8  →  1 命中（本仓 0）   ← 非恒真、非装饰
G1b 与 G1 同一枚文件即响（该行的 `b.Execute(` 也被 G1b 吃到）
G2  合成树多出一枚所有者 close_owner_probe.go → 2 命中（本仓就是这 2 行的形状）
G3  合成树 internal/agent/loop_with_taskid_probe.go:5 RunWithTaskID → 命中（本仓 0）
G4  本仓排掉 _test 即空、摘掉排除立刻 2 命中（harness_test.go:116 true ＋ loop_approval_test.go:128 false）
```

⚠ **本程自己撞的一枚，记下来给下一程**：G3 的正控第一发是**空**的，看着像"这枚子句结构上产不出读数＝装饰"。
反查原因＝本程把探针文件放到了 `internal/` 而 pathspec 是 `internal/agent`（`git grep` 的 pathspec 是**目录前缀**，`internal/agent` 不吃 `internal/x.go`）。
`git mv` 进 `internal/agent/` 后同一把尺命中。**教训：正控为空时先怀疑自己的探针放错地方，再怀疑检装饰。**

### 2.4 两问自查（恒真检／装饰）

1. **有没有哪一枚今天必然响？** 没有。G1/G1b/G3/G4 今日全空、G2 今日就是那 2 行。
2. **有没有哪一枚结构上产不出自己想量的读数？** 没有——五枚都有第 2.3 节最后一列的正控；且 G1 的**负一负**在本仓就是现成的：
   同一把尺只摘掉"排除测试"那一条 ⇒ 今日 **14 枚 `ToolRequest{` 命中在 `_test.go` 里**（`gate-clauses-roster.txt` 末段），
   也就是说这个模式**天天在吃得住本票那一形**，只是它今天吃到的全是测试。
3. **本程没有为此新增任何断言**（零生产码＋写面只有注释）。把"今日名册为空"写成 `t.Fatalf` 就是一枚恒真检——它今天绿、明天形状真出现时也不会红（那时它红的是"名册非空"这个**事实**，不是"你没接上关闭"这个**义务**），所以那种断言既买不到防回归、又会诱使下一位把"名册空"当成"没问题"。**明确不做。**

⇒ AC#1① 的判据交付＝第 2.3 节那张表；AC#1③（"今天不可达"写在门本体上）落在 `internal/tools/bridge.go` 的注释里，**原文整段见第 3.1 节**。
⇒ AC#1②（触发时必须做的事）也写进同一处注释，三条：为那一形接上关闭的所有者（或把 `Q-56` 拍板的那一支落地）、**同时**补第 2.5 节那形并发读数、把本门的新名册回贴到本件（只追加）。

### 2.5 "并发那一形零读数"——本程现量三条，以及它今天为什么不可达

| 量 | 命令 | 读数 |
|---|---|---|
| 一个进程里跑几枚环路任务 | `git grep -nE "agent\.New\(|loop\.Run\(" -- cmd internal ':!*_test.go'` | `cmd/wisp/run.go:573` ＋ `:604` **各一枚** ⇒ 一进程一枚环路任务（`runTextTask` 一次一条 `task`）。另两枚 `.Run(` 命中是 `bridge.Run`（`cmd/wisp/models.go:311`、`cmd/balldebug/main.go:252`），**不是** `agent.Loop`，本程逐个点开看过。 |
| D38d 的 4 路并发是不是"两枚任务" | `sed -n '604,654p' internal/agent/loop.go` | `sem := make(chan struct{}, guard.Concurrency())` 限的是**同一回合内的多发调用**，而它们共用**同一枚 `req`（同一枚 `TaskID`）**（`loop.go:646-649`）⇒ 不构成"两枚 scope 同时开着"，撞不到 `unbound-scope` 那条 fail-close。 |
| 有谁在测试里让两枚 scope 同时开着过 | `git grep -n "unbound-scope" -- '*.go' ':!internal/risk/provenance.go'` | 三处：`cmd/wisp/task_scope_close_151_test.go:14`（注释）＋ `:108`（断言），那是**串行**那一发（先关 A、再开 B）；第三处 `.scratch/wisp/probes/151/ac1-probe-source.go:105` 是上一程留在 `.scratch` 的探针文本，**不参与编译**（`git ls-files -- 'internal/**/*.go' 'cmd/**/*.go'` 不含它）。全仓**没有**一枚用例让两枚 scope 同时开着。 |

⇒ 判据本体（`internal/risk/provenance.go:438-457` 的 `scopeMarks`，本程**只读**）：`unbound-scope` 只在"被检的那枚 scope 尚未注册、而**别的 scope 手里有污点**"时响。
⇒ 所以并发那一形今天**双不可达**：既要等宿主自带 id（G1/G3），也要等**同时跑两枚任务的宿主**（今天连 `agent.Loop` 的第二枚生产持有者都没有）。
⇒ 本票因此**不许为它硬开一道"要它响"的格**（票面 §现量形状第 5 条），本程照做：只在门上留了 G4，并在注释里写明"并发读数＝零，等 G1/G3/G4 任一响时一起补"。

---

## 第 3 节 格 AC#2（＋AC#1③ 的落点）— `Bridge.CloseTask` 注释改后原文整段

### 3.1 先答票面那一问："已经有的那一句够不够"

原判句（改前 `internal/tools/bridge.go:649-653`，现读仍在原位）：
> "A caller that dispatches on the bridge outside that boundary - a host-internal call with its own task id - still has no owner here, and stays fail-closed…"

本程判：**方向对，但不构成会被读到的防呆**，三条不足（每条都对应票面 AC#2 的判据"一个只读注释的人能不能答出'我这条腿要不要自己关'"）：

1. **不点名票号**：读的人不知道"151 已经做过一次判断"，于是会把"still has no owner"读成"这是一句 TODO 风格的自谦"，而不是"这是一个**已登记、故意没接**的口"。
2. **没给自测问句**：它描述的是"有一种调用方"，没有把判据交到读者手上——读者要的是"我怎么知道我是它"。本程补的就是那句
   `does the TaskID you dispatch with come from that boundary?`（第 3.2 节 B 段第 4 行起）。
3. **没写"今天这一形在生产上还没有入口"**：这半最要命。不写它，面板/球上线那天读这段注释的人会以为"接上就是了，加个 defer"；
   实际上 `Loop` 没有收外部 task id 的导出方法（票面 §现量形状第 4 条），那是 `Q-56`，是**人工拍板项**。
   本程把这两样都写进注释，并写明"本文件不回答 Q-56"。

⇒ 所以 AC#2 的动作＝**在原地补两段**（不抹原句、不改原句一字），并给"门在哪"一个可点的落点（本件 §2.3）。

### 3.2 改后注释原文（整段，含行号；机器可核＝`probes/154/bridge-comment-after.txt`）

**A. `Bridge.OpenTask`（`internal/tools/bridge.go:626-632`）** —— 只补了"开合两侧不在同一作用域"这一句指路：

```go
626| // OpenTask opens the C25 taint scope for one task - ticket 19's
627| // DEFERRED(C25-loop-wiring) item (1). Idempotent; Execute opens lazily so a
628| // caller that forgot cannot get an untainted read.
629| //
630| // The open side is per CALL (mark above), the close side is per TASK and lives
631| // in another package, so nothing here keeps the two sides in step by
632| // construction: read CloseTask before adding a caller that opens a scope.
633| func (b *Bridge) OpenTask(taskID string) {
```

**B. `Bridge.CloseTask`（`internal/tools/bridge.go:646-679`）** —— 前 8 行是**原句未动**，中段是本程新增，末尾"audit line"那 4 行也是原句未动：

```go
646| // CloseTask closes one task's taint scope. The composition root defers this on
647| // the task's own boundary: cmd/wisp wires it into agent.Options.AdmitTask's
648| // revoke, which the TEXT loop already defers (internal/agent/loop.go:366), so a
649| // task that ends takes its taint with it (ticket 151). A caller that dispatches
650| // on the bridge outside that boundary - a host-internal call with its own task
651| // id - still has no owner here, and stays fail-closed: a task that never closes
652| // leaks only its own indexed sources, which is the fail-closed direction, since
653| // dropping marks would open R4.
654| //
655| // WHICH LEGS STILL HAVE NO OWNER HERE - ticket 154, read this before assuming
656| // ticket 151 finished the job. 151 put the close on ONE boundary: cmd/wisp's
657| // admitTask hook, which the TEXT loop defers. That covers only ids the loop
658| // mints itself (Run/RunAsync -> newTaskID, internal/agent/loop.go:332/:321).
659| // So ask one question about your own leg: does the TaskID you dispatch with
660| // come from that boundary? If it comes from anywhere else - a panel or ball
661| // host inventing its own id, a retry wrapper, a scheduled task carrying one id
662| // across turns - nothing in this tree closes it, and no test will tell you so.
663| // Today no production code dispatches on the bridge except loop.go, so being
664| // such a caller means being the first one; the gate meant to ring when that
665| // happens is docs/evidence/s1/154-host-id-never-closed-r1.md §2.3 (clauses
666| // G1/G1b), and it is one `git grep` away, not a note in someone's head.
667| //
668| // Being first is also an API decision, not just a missing defer: Loop has no
669| // exported method that takes a caller-supplied task id, so a host cannot route
670| // its id through the boundary that owns the close. That choice is Q-56 and this
671| // file does not answer it. When the shape does become reachable, closing your
672| // leg is only half of what ticket 154 owes - the reading for two tasks with
673| // scopes open at once is still zero, and it has to be measured in the same
674| // change (same evidence file, §2.5).
675| //
676| // The audit line runs on EVERY call, closed scope or not, because "who took
677| // this task's taint off the table, and was there anything on it" is exactly the
678| // question the C25 ledger owes a long-running host. A scope that was never
679| // opened (a task that read no sensitive source) closes as dropped=0.
```

注释里点到的三处**都先反查过**（票面"点名符号前先 grep"）：`internal/agent/loop.go:366` 现读＝`defer revokeAdmission()`（原引用没烂）；
`Run` 在 `:332`、`Async` 在 `:321`、`newTaskID` 在 `:1107`（本程 `grep -n "^func (l \*Loop)"` 与 `grep -n "func newTaskID"` 现读）；
`Q-56` 在 `docs/reports/pending-and-issues.md` 有登记条目（本程只引用编号，不替它选支）。

### 3.3 写面机器核：这次改动**一行都不是码**

```
$ git diff -U0 -- internal/tools/bridge.go | grep -E "^[+-]" | grep -vE "^(\+\+\+|---)" | grep -vE '^[+-][[:space:]]*//'
                                                    →  无输出（rc=1）  ← 增删行全部以 "//" 开头
$ git diff --numstat -- internal/tools/bridge.go    →  25  0  internal/tools/bridge.go（+25 / -0）
```

⇒ **不改签名、不新增 API、不动任何语句**（票面"本票不解决的事"第 1 条）。

### 3.4 格 AC#5（改后两跑）＋ 名册两向 ＋ 另外三道工具

同尺同 PATH（`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v <包>`）：

| 包 | 形 | RUN | 顶格 | FAIL | SKIP | panic | 全量 | 台件 |
|---|---|---|---|---|---|---|---|---|
| `internal/tools` | 改后 `-count=1` | **115** | **79** | **0** | **0** | 0 | 115 | `probes/154/gate-post-tools-v.txt` |
| `cmd/wisp` | 改后 `-count=1` | **139** | **79** | **0** | **0** | 0 | 139 | `probes/154/gate-post-cli-v.txt`（`ok 91.747s`） |

⇒ 四数与第 1 节那两行**逐字相同**（改前＝改后），`rc=0` 两包皆然。

```
$ comm -23 names-pre-cli.txt names-post-cli.txt   →  空
$ comm -13 names-pre-cli.txt names-post-cli.txt   →  空
$ comm -3  names-pre-tools.txt names-post-tools.txt  →  空
（台件：probes/154/names-{pre,post}-{cli,tools}.txt；尺＝第 1 节那把，两栏皆 sort 过）
$ go vet ./internal/tools/                        →  无输出，rc=0
$ $(go env GOPATH)/bin/gofumpt --version          →  v0.12.0 (go1.27.1)   （现读，不背 151 那两程的数）
$ $(go env GOPATH)/bin/gofumpt -l internal/tools cmd/wisp  →  空
$ sh scripts/d22scan.sh                           →  clean（rc=0）；分母：bans#1-5 internal/=205 cmd/=23 · ban#6 frontend/=83 · ban#7 internal/tools/=18
                                                     · ban#8 design/=39 frontend/=83 internal/=413 cmd/=44   ← 各作用域分母全非零
```

**承重第二问（派单点名那句：摘掉你新写的哪一味，哪条用例转红）**：
本程新写的**只有注释**（§3.3 的机器核），Go 的注释不进 AST 输出、不进二进制 ⇒ **摘掉它，零枚用例转红**。
这不是"检没牙"，是本票的性质：**本票买的是防忘记，不是防回归**（票面 AC#6 的明文允许）。
"有没有任何外部可见读数变过"＝**没有**，三条独立证据：四数逐字同、两包名册两向 `comm` 皆空、`go vet`/`gofumpt -l`/`d22scan` 三项与改前同（且 `CloseTask` 那行审计日志的格式串在 `:693` 未动，`tools: C25 scope closed task=… was_open=… dropped=… open_scopes=…`）。

⇒ **AC#5 这一格：成立**（两包各两跑、四数、名册两向、另加 vet/gofumpt/d22scan）。

---

## 第 4 节 格 AC#3 — 同族普查（"只开不合／注册即泄漏／defer 挂在今天不可达的边上"＋"被上游码挡住的可见性"）

生成器＝`.scratch/wisp/probes/154/census-pairs.sh`，15 对逐字原始输出＝`.scratch/wisp/probes/154/census-pairs.txt`（锚点 `352d3d8`）。
尺统一＝`git grep -nE`（大小写敏感、**不**整词、pathspec＝`internal`＋`cmd`、排除 `_test.go`），每对另给的排除项写在 `census-pairs.txt` 的标题行里。
"生产调用者"一律数**出现点行数**（定义行与注释行在下表逐枚标明，不混进计数）。

### 4.1 成对 API 那一族

| # | 成对 API | 开那侧（生产出现点） | 关那侧（生产出现点） | 一句判读 |
|---|---|---|---|---|
| 1 | `Bridge.OpenTask` ↔ `Bridge.CloseTask` | **1**：`internal/tools/bridge.go:559`（长在 `mark()` 里，每次调用） | **1**：`cmd/wisp/run.go:563`（长在 `admitTask` 的 revoke 里，每枚任务） | 本票主题。两侧都只有 1 枚，但**枚数相同不等于对称**：开按 call、合按 task，且合只覆盖环路自造的 id。 |
| 2 | `risk.Provenance.OpenScope` ↔ `CloseScope`（冻结面，本程只读） | **2**：`bridge.go:638` ＋ **`cmd/wisp/panel_assets.go:232`** | **1**：`bridge.go:666`（在 `CloseTask` 内） | **同族第二枚真"只开不合"**：`panel-assets -l2` 那枚诊断腿开了一枚 `taintSourceScopeID`（`panel_assets.go:166` 常量），全仓没有任何 `CloseScope` 指向它。后果有界（进程一次性退出＋失败方向是 fail-closed），但它证明"只开不合"不是孤例。 |
| 3 | `Gate.AdmitTextTask(taskID) → revoke func()` | **2**：`run.go:477`（`host:mode-switch` 那枚宿主自造 id）＋ `run.go:558`（环路那枚） | **2**：`run.go:478 defer revoke()` ＋ `run.go:561 revoke()` | 两枚都有 owner。⚠ 形状上给下一位的一条：**关那侧是闭包，按方法名 grep（`Close*(`/`Revoke(`）永远抓不到它**——本票的门因此不吃"关"的名字，只吃"开"的构造点。另：`host:mode-switch` 今天**不开** scope（它不经过桥，只 `gate.PendingApproval`），所以它不属"宿主自带 id 没人关"那一形，本票不把它算进去。 |
| 4 | `plugin.DisposalScope`（C11）`NewDisposalScope`/`Defer` ↔ `Dispose` | **0**（只有 `internal/plugin/disposal.go:155` 定义＋两处注释） | **0**（`.Dispose(` 生产零命中） | **整对 API 今天没有生产使用者**。连带效应最要紧：`internal/risk/provenance.go:56` 与 `:348` 那两句义务写的是"OpenScope at task start and **Defer(CloseScope) on the task's DisposalScope**"——那是一条**挂在今天不可达的边上的 defer**；票 151 之所以只能把 close 长在 `AdmitTask` 的 revoke 上（`run.go:546-552` 自己写了拒绝理由），根因就是这里。 |
| 5 | `memory.Store.StartRetentionJob(scope, cfg)` | **0**（定义 `retention.go:98` ＋两处注释） | 同左（它自己就是开侧） | 收 `*plugin.DisposalScope` 且 nil 就 panic（`:99-101`）的那道门，今天没有生产读者——与第 4 对同因。 |
| 6 | `tools.Registry.Register` / `RegisterProvider` ↔ **无 Unregister** | `Register`：**1**（`run.go:345`，启动期在循环里注册 `BuiltinFSEntries`）；`RegisterProvider`：**0**（定义 `registry.go:163`，注释写着"three unlanded slots"） | **不存在**：`git grep -iwE "unregister" -- internal cmd ':!*_test.go'` 只剩 ball 的 Win32 命名与 audio 的一行注释（见第 13 对），Go 侧没有 `Registry.Unregister` | "注册即泄漏"在这一族**成立但有界**：注册表随进程寿命、启动期一次填。本程判＝不必修，但**下一位若加运行期热插拔（票 104 `-plugins-load` 那条腿）就必须回头要一枚 Unregister**。 |
| 7 | `risk.RiskAssessor.RegisterSubAssessor` | **0**（定义 `assessor.go:214`＋注释） | **0**（`UnregisterSubAssessor` 全仓零命中） | 两侧都无使用者＝死代码形，不是泄漏形。登记，不扩大。 |
| 8 | `proc.JobScope.StartInJob` ↔ `Close` | **1**：`cmd/wisp/slo_windows.go:498` | **3**：`internal/proc/boot_windows.go:97`、`:109`、`:147`（`hooks.CloseJob`） | 有 owner，且关闭那侧比开侧多（两条退出＋一条挂到 shutdown hook）。 |
| 9 | `proc.AcquireSingleInstance` ↔ `(*SingleInstance).Release` | **1**：`boot_windows.go:95` | **1**：`boot_windows.go:151` | 有 owner。 |
| 10 | `audio`：`NewHalfDuplexGate` / `WASAPIMicrophone{}` / `NewWavInjector` ↔ `Stop()` | **0**：三者在 `internal`＋`cmd` 的生产码里都只出现在自身定义与包文档（`internal/audio/gate.go:84`、`wavinjector.go:48`） | 只有包内互调（`gate.go:125`、`:256`） | **整条麦克风腿今天不接任何宿主** ⇒ 这一族的成对 API 两侧都是零读数。与本票"宿主 id 那一形零读数"同形：**不能因为绿测试多就说它接过了**。 |
| 11 | `panel.StreamLog.Append(key,…)`（隐式开一枚 key） ↔ `Close(key)` | **1**：`cmd/wisp/run.go:757`（`c.stream.Append(e.TaskID, e.Text)`） | **1**：`cmd/wisp/run.go:784`，但它**只长在 `agent.EvDone` 那一个 case 上** | 同族实例（"关只挂一条出口"）：以 `EvError`／`EvStuck`／取消收尾的任务，其 key 永远不被 Close ⇒ 面板那段结果一直"在流"。有界＝`const DefaultStreamKeys = 32`（`internal/panel/pump.go:279`）＋ `mergeOverflowLocked`（`:344`），所以它**不是内存泄漏，是可见状态残留**。本程登记，**不在本票修**（要修得动 `internal/panel/**`＋`cmd/wisp/run.go`，都在写面外）。 |
| 12 | `memory.Store.InsertGrant` ↔ `RevokeGrant` / `DeleteGrant` | **0**：只有定义 `dao_misc.go:17` | **0**：只有定义 `:49`／`:102` | D45 授权表的写入与撤销今天都没有生产调用者 ⇒ 两侧都零，不判。 |
| 13 | ball 热键：`RegisterHotKey` ↔ `unregisterAll` | 注册在 `internal/ball/hotkey_windows.go:351-362`（＋ `RegisteredHotkeys()` 的读数口 `ball_windows.go:837`） | **2**：`ball_windows.go:804`、`:913`（后者在 `Ball.Close`（`:904`）内） | 关那侧**存在**（本程第一发以为"没有 UnregisterHotKey"是尺太窄：`git grep -iw Unregister` 整词尺抓不到 Win32 的 `pUnregisterHotKey`）。但 `Ball.Close` 本身今日生产调用者只在 `cmd/balldebug`（`main.go:234/266/358`）⇒ 常驻那条腿今天不销毁球。登记，不判。 |
| 14 | `observe.LogPipeline` ↔ `Close` | **1**：`cmd/wisp/logsink.go:149`（`observe.InitLog(observe.LogConfig{…})`，字段声明在 `:94`） | **4**：`logsink.go:113` ＋ `slo_windows.go:304/313/327` | 有 owner（三条退出各关一次）。 |

**普查里本程自己弄错又改回来的一枚**（写给下一位，别当"这族盘得很顺"）：第 13 对本程第一发用 `git grep -iw "unregister"` 得到"全仓零命中"，据此差点写成"球的热键注册即泄漏"。
反查＝**整词尺吃不到** Win32 过程指针名（`pUnregisterHotKey`／`hotkeyUnregisterer`／`unregisterAll` 都不是整词 `unregister`），换成 `git grep -inE "unregister" -- internal cmd ':!*_test.go'` 后命中 **21 行**（尺：不区分大小写、子串、只看生产码；本程现读）。
⇒ 教训与本票的门同形：**"零命中"必须先说自己是哪把尺**（票面 10:3x 那条 `>` 讲的正是这个）。

### 4.2 族二：被上游码挡住的可见性（票面 11:4x ① 要求一起盘的那半）

这一族盘的不是"有没有人关"，而是**"端到端看不见某一形"有几个独立原因**。本程现量到**三枚**，比票面引的那两处多一枚：

| 锁 | 现量 | 读数 |
|---|---|---|
| **L-1 声明档** | `grep -n -A 4 "func FSReadDecl" internal/tools/fs.go` | `Declared: risk.L0`（fs.go:298；`fs.list` 同：`:311`） |
| **L-2 环路分支** | `git grep -nE "case RiskL1, RiskL2:|default:|PassThroughUnclassifiedRisk" HEAD -- internal/agent/loop.go` | `decideRisk`（`loop.go:773`）只把 **L1/L2** 放行给宿主审批（`case RiskL1, RiskL2:`＝`:776`）；其余（含 L0）全落 `default:`＝`:782`，那里 `if !l.opt.Config.PassThroughUnclassifiedRisk` → `DecisionReject`＝`:783-785`，`p.skip=true` 后 `:636-639` 直接不派发 |
| **L-3 组合根从不打开它** | `git grep -n "PassThroughUnclassifiedRisk" HEAD -- cmd/wisp/run.go` | **零命中**（Config 字面量在 `run.go:588-595`，本程逐行读过：`Model/ContextWindow/ArtifactsDir/PerToolTimeout/SteeringEnabled`，没有这一项） |
| **L-4 只有 `fs.read` 能被 Mark** | `census-pairs.txt` 第 15 对的交集：注册名 `{fs.delete fs.move fs.read fs.trash fs.write}` ∩ `sensitiveSourceTools{fs.read search.content clipboard.read system.get web.fetch doc.read screen.capture asr.transcript}`（⚠ 脚本那行 `Src[A-Za-z]+ += "` 还会把两枚 **fail-closed 标记** `SrcUnboundScope`/`SrcUnscannedNesting` 一起捞出来，所以台件印出的清单比上表多这两枚；交集不受影响） | **交集＝`fs.read` 一枚** ⇒ `mark()`（`bridge.go:552` 的 `risk.IsSensitiveSource`）今天只可能为它响 |
| **L-5（反向）测试开着、生产关着** | `git grep -n "PassThroughUnclassifiedRisk:" HEAD -- '*.go'` | 命中两枚，**都在 `_test.go`**：`internal/agent/harness_test.go:116 = true`、`internal/tools/loop_approval_test.go:128 = false` ⇒ 环路那批绿测试走的是**生产没有的那道门**（L-3 开着） |

⇒ 合起来的判读（这句是本票真正要留给下一位的）：**票 151 接上的那一形（环路 id 的 close），今天在生产上没有任何 open 可关**——L-1＋L-2＋L-3 三环把 `fs.read` 挡在桥外，L-4 又说明即使接上、今天也只有这一枚工具会打污点。
⇒ 具体后果分两型，别混：**T1 型**（`TestCompositionRootClosesTheLoopTasksTaintScope`）在 mockllm 那一形上绿，但它绿的是"关闭那发确实带着本轮 id 出现了"（`CloseTask` 的审计行**每次调用都写**，`bridge.go:676-679` 自己写着 `dropped=0` 也照写）——**不等于生产上真有一枚 scope 被关掉过**；
**T2 型**（`TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit`）今天全仓**唯一**能让 `OpenTask` 真响起来的调用形，恰恰是"宿主自带 id 直接 `bridge.Execute`"（`cmd/wisp/task_scope_close_151_test.go:84/:92`），也就是本票说"没人关"的那一形。
⇒ 〔未取证，本程只登记〕"环路上线之后 `fs.read` 走不走 L0 直通／改不改声明档"是 **L-1～L-3 的联合决定**，不在本票写面内，也不属 `Q-56` 那三支出路；本程**不判该走哪条**。

### 4.3 盘不到的那半（票面"不许写成'应当没有'"）

上面 14＋5 行只覆盖**名字里带开/合语义**与**注册语义**的成对 API。以下三类本程**盘不到**，写清命令与原因，不外推：

1. **"关"不叫任何 lifecycle 名字的隐式配对**（例：`map[key]=true` 开了、靠别处 `delete` 合）。
   本程命令＝`git grep -nE "^\s*b?\.?[a-zA-Z]*\[taskID\] = " -- internal cmd ':!*_test.go'` 之类只能撞运气；**没有**跨函数数据流仪器（第 2.2 节末：本模块 require 里没有 `golang.org/x/tools`）。⇒ 这一类**未盘**。
2. **`defer` 挂在不可达边上的完整清单**。本程只撞到了第 4.1 表第 4/5 对（DisposalScope 那族）。
   完整盘法需要"哪些边今天不可达"的全仓可达性判定，同一枚缺仪器。⇒ 这一类**未盘**。
3. **`internal/plugin/**` 与 `internal/speech/**` 等本票 pathspec 之外的包**：本程尺只放了 `internal`＋`cmd`（这俩已含 plugin/speech 的全部 `.go`），但**`tools/**` 与 `cmd/*` 里非 `wisp` 的宿主（`balldebug`/`llmrecord`/`modelcheck`）本程只按上面逐枚点名，没有做穷举**。⇒ 这三处**未穷举**。

---

## 第 5 节 格 AC#4（**续程 r2 补写**）— 契约轴零字节：把"普查"写成逐支账

> 本节作者＝票 154 的**续程**（派单＝`.scratch/wisp/dispatches/2026-09-26-124x-impl-154-r2-two-cells.md`，53 行，现读在盘）。
> r1 那程把 AC#4 的**读数**跑完但没来得及 commit，三枚读数由编排者 12:3x 代提（`e9a1db0`）；**代提只到落盘为止，判语在本节**。
> 本文件前 363 行（§0–§4）本节一字未改（`git diff --numstat` 的核法见 §5.8 末段）。

### 5.0 锚点，以及对派单／编排者每一条前提的反查

共享树在本节写作期间自己往前走了 **6 枚 commit**（155 验收程两枚、台账两枚、另两枚）。本节的每一发读数都写清它自己在哪枚锚点上量的：

| 读数 | 锚点（`git rev-parse --short HEAD` 当轮现取） |
|---|---|
| 本程进场第一读 | `9835d81` |
| AC#4 逐支尺 R1 | `9835d81` |
| AC#4 逐支尺 R2（五枚输入那一发） | `edd0725` |
| AC#4 逐支尺 R3 ＋ 门复跑两发 ＋ §6 的变异与正控 | `be603fd` |
| 本节 commit 前的最后一读 | 见 §5.8（当轮记） |

**逐条反查（派单里凡"前程说 X"都是自述，本节按盘上判）**：

| 断言 | 反查命令（当轮跑的） | 读数 | 判定 |
|---|---|---|---|
| 前程末枚＝`304199f`＝第 4 节／AC#3 | `git log -1 --format="%h %ad %s" 304199f` | `304199f3 12:15:18 evidence(票 154 实现程 第 4 节 格 AC#3)…` | **成立** |
| 三枚读数由 `e9a1db0` 入库 | `git show --stat --format="%h %ad %s" e9a1db0` | 3 files changed, 85 insertions(+)：`zero-byte-per-commit.sh`／`zero-byte-per-commit.txt`／`g3-open-side-via-loop.txt` | **成立** |
| `e9a1db0` 含前程全部四节 | `git merge-base --is-ancestor 304199f e9a1db0`；`git show e9a1db0:docs/evidence/s1/154-….md \| wc -l`；同件 `\| grep -c '^## '` | rc=0；**363**；**5** | **成立**（＝派单说的"363 行／5 个 `## ` 节"逐字复到） |
| AC#1／AC#2／AC#3／AC#5 已裁完并 commit | §5.1 名册逐枚对上本文件的第 1／2／3／4 节；四枚 subject 逐枚 `git log -1` 现读 | 第 1 节＝AC#5 改前、第 2 节＝AC#1、第 3 节＝AC#2＋AC#1③＋AC#5 改后、第 4 节＝AC#3 | **成立**（就"已写、已入库"而言；**成不成立由非实现者验收表判，本节不重跑、不改写**） |
| 编排者"我这边 12:4x 是 `245d5f4`" | `git cat-file -t 245d5f4` → `fatal: Not a valid object name`（rc=128）；`git cat-file --batch-all-objects --batch-check \| grep -c '^245d5f4'` → **0** | 本仓**不存在**任何以 `245d5f4` 为前缀的对象 | **〔不成立〕**＋当轮读数。派单同时明写"以你现量为准"，故本节一律按上表锚点走，**不据那枚 sha 判任何事**；登记给编排者（这不是注入，是一枚读不出来的数） |
| 另两枚程在飞（155 验收／156 实现） | `git log --oneline 9835d81..HEAD`；`git rev-parse <c>:cmd/wisp/slo_windows.go` 六点位 | 155 那两枚已入库（`67cc43f`／`edd0725`）；`cmd/wisp/slo_windows.go` 的 blob 从 `1424aa7^` 到 HEAD **全程同一枚 hash** ⇒ 156 那程**本节写作时还没提交过它** | **成立**（且给出时间边界） |

### 5.1 名册现算（不抄前程表，也不抄任何人的枚数）

```
$ git log --pretty=tformat: --name-only 1424aa7 352d3d8 76c89b2 304199f | sed '/^$/d' | sort -u | wc -l
1728      ← 派单那条命令**照字面**跑出来的数：不带 --no-walk，git 会沿这四枚的父链把全部祖先一起算进来。
            这 1728 枚里冻结面命中**非零**（别人的 commit），所以它既不是本票的名册、也不能拿来判零字节。
$ git log --no-walk=unsorted --pretty=tformat: --name-only 1424aa7 352d3d8 76c89b2 304199f | sed '/^$/d' | sort -u | wc -l
19        ← 四枚 commit 真正碰过的文件（去重）
$ 逐枚条目合计 = 22（19 枚去重 ⇒ 本文件那枚证据件被 4 枚 commit 各碰一次，另 3 次是重复计数）
```

名册（19 枚，全部当轮现算；`17` 枚 `.scratch/wisp/probes/154/**` ＋ `1` 枚本证据件 ＋ `1` 枚 `internal/tools/bridge.go`）：
逐枚的 `N1 = --name-only` 与 `N2 = --name-status --no-renames` 两把尺并核（改名折叠 vs 改名拆两侧）：

| commit | N1 | N2 | N1≡N2 |
|---|---|---|---|
| `1424aa7` | 6 | 6 | YES |
| `352d3d8` | 6 | 6 | YES |
| `76c89b2` | 7 | 7 | YES |
| `304199f` | 3 | 3 | YES |

⇒ 今天没有"改名把冻结面文件挪到面外"那一形（若有，N1 会小于 N2）。名册全文与这两把尺的原始输出＝`probes/154/ac4-per-face.txt`（R2 段 §1）。

### 5.2 支数＝13（本节作者自己从票面数的，命令与枚举都贴出来）

```
$ T=.scratch/wisp/issues/154-the-close-only-covers-…-trigger-gate.md
$ sed -n '/AC#4 契约轴零字节/,/^- \[ \] \*\*AC#5/p' $T | tr '、' '\n' | sed '/^[[:space:]]*$/d' | nl | head -13
     1  - [ ] **AC#4 契约轴零字节**：`docs/PLAN.md`      8  `thresholds.go`
     2  `docs/specs/**`                                9  任何 golden
     3  `internal/risk/**`                            10  `allowlist.txt`
     4  `internal/panel/**`                           11  `scripts/slo-check.ps1`
     5  `internal/agent/**`                           12  `tools/d22scan/**`
     6  `internal/agent/approval/**`                  13  `cmd/wisp/slo_windows.go`
     7  `internal/observe/**`
（第 14 项起是票面那两条 ⚠ 限定语，不是禁改面：①`frontend/**`／`design/**` 不碰不还原、**不算进任何零命中宣称**；②本票不改 `Loop`／`Bridge` 签名）
```

⇒ **票面 AC#4 点名 13 支**。派单点的两枚"最容易漏"——`internal/agent/**`（第 5 支）与 `cmd/wisp/slo_windows.go`（第 13 支）——本节都单列成行，理由与派单一致：一支是被第 6 支包含的父树、一支是整棵冻结树里唯一的单文件。

### 5.3 逐支账（13 支 × 4 枚 commit；生成器＝`probes/154/ac4-per-face.sh`，原始输出三份＝`probes/154/ac4-per-face.txt`）

尺面（每格都得对这句话负责）：**名册**＝单枚 commit 的 diff 路径（按 commit 量，不按区间量）；**匹配**＝`grep -E`、大小写敏感、路径前缀锚定、**不**整词、**含 `_test.go`**（本尺盘文件不盘符号）；**分母**＝`git ls-files` 在 HEAD 上命中该支的 tracked 件枚数；**每一发都跑两遍取数并记 `grep` 的 rc**（`0/1` ＝ 命中 0 枚、rc=1＝真零命中；`128` 会另写，见 §6.3 那一发）。

| # | 支（票面逐字） | HEAD 上分母/rc | `1424aa7` | `352d3d8` | `76c89b2` | `304199f` | 正控（同一把尺打在册路径上） |
|---|---|---|---|---|---|---|---|
| 1 | `docs/PLAN.md` | 1/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`docs/PLAN.md`） |
| 2 | `docs/specs/**` | 14/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`docs/specs/README.md`） |
| 3 | `internal/risk/**` | 37/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/risk/assessor.go`） |
| 4 | `internal/panel/**` | 20/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/panel/approval.go`） |
| 5 | `internal/agent/**` | 58/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/agent/approval/approval.go`） |
| 6 | `internal/agent/approval/**` | 18/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（同上，被第 5 支包含的那棵子树） |
| 7 | `internal/observe/**` | 25/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/observe/clock.go`） |
| 8 | `thresholds.go` | 1/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/observe/thresholds.go`） |
| 9a | 任何 golden（**窄**：`*/testdata/golden/*`） | 52/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/agent/testdata/golden/budget-loop.sse`） |
| 9b | 任何 golden（**宽**：路径含 `golden`，大小写敏感；大小写不敏感同枚 58） | 58/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`internal/agent/loop_golden_test.go`） |
| 10 | `allowlist.txt` | 1/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`tools/d22scan/allowlist.txt`） |
| 11 | `scripts/slo-check.ps1` | 1/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（同名件） |
| 12 | `tools/d22scan/**` | 6/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（`tools/d22scan/allowlist.txt`） |
| 13 | `cmd/wisp/slo_windows.go` | 1/0 | 0/1 | 0/1 | 0/1 | 0/1 | 正控=1（同名件） |

13 支全部**逐支 0 命中**、分母全非零（＝尺不是空的）、正控全打得住（＝尺不是死的）。同一把尺对**重复输入**的自反性（R2 那一发把 `304199f` 送进去两遍）：§2 的 14 行末两列逐行相同、§2b 的 4 行相同（`probes/154/ac4-per-face.txt` 末段那发 awk，其第一版尺放错、已在同处更正）。

票面**未点名**、由 r1 脚本自带正则加严的四支一并答（同一把尺，读数同样 0）：`internal/speech/**`（分母 1＝`doc.go`）· `internal/secret/**`（15）· `.gitattributes`（1）· `go.sum`（1）。

### 5.4 r1 那把尺（`zero-byte-per-commit.sh`）的四处不足——写下来不是挑刺，是因为"逐支"这一格它给不出来

1. **只给总数、不给逐支**：它每枚 commit 印一行"冻结面命中 = 0"。总数为 0 时逐支当然也为 0，但**漏答一支**（比如把 `internal/agent/**` 整支从正则里删掉）在这一版输出里看不出来。本节把 13 支摊成 13 行。
2. **名册尺只有一把**：`git show --name-only` 默认开改名探测 ⇒ "把冻结面文件改名到面外"那一形只会露出新名。本节两把尺并核（§5.1 表，今日四枚皆 `N1≡N2`）。
3. **`tools/d22scan/**` 被写成 `allowlist\.txt|.*\.go`** ⇒ 该目录 6 枚 tracked 件里的 `go.mod` 与 `runtests.sh` **不在尺上**（现量：`git ls-files -- tools/d22scan`＝6 枚，本节尺按 `^tools/d22scan/` 整目录）。
4. **"任何 golden" 被写成 `.*/testdata/golden/.*`** ⇒ 漏掉 `internal/llm/golden/`（3 枚 `.go`）与两枚 `*_golden_test.go`；本节宽窄两把都跑（9a/9b，分母 52 vs 58）。
   另记一句账：它头里写"＝派单 AC#4 ∪ 票面 AC#4"，但它的正则实为 **18 项**，而"票面 13 支 ∪ 旧派单 14 支（多出 `frontend/**`·`design/**`）"＝**15 支**；多出的是上表那四支加严。**加严不是错**，但别把 18 读成"票面要求 18 支"。

### 5.5 复跑 r1 的脚本（当轮输出，不只引用代提那份）

```
$ bash .scratch/wisp/probes/154/zero-byte-per-commit.sh 1424aa7 352d3d8 76c89b2 304199f
（全文＝probes/154/zb-per-commit-rerun-r2.txt：四枚各自"冻结面命中 = 0 ✓"，名册逐枚列出）
$ diff <本轮输出> .scratch/wisp/probes/154/zero-byte-per-commit.txt   →  diff rc=0（逐字节相同）
```
两发（`9835d81` 与 `edd0725` 两枚锚点各一发）也与入库那份逐字节相同 ⇒ 该脚本只读 git 历史，**与锚点无关**；这与 §5.7 那把区间尺的行为正好相反。

### 5.6 第二口径：面级 tree／blob hash 逐枚相等（字节级，不看文件名）

命令＝对每个面取 `git rev-parse <点位>:<面>`（目录取 tree hash、单件取 blob hash），点位＝`1424aa7^`（四枚之前的 base）、四枚 commit、`HEAD`。原始表＝`probes/154/ac4-per-face.txt` §3。

```
面（11 棵目录 ＋ 6 枚单件）      六点位 hash 全等?
docs/specs · internal/risk · internal/panel · internal/agent · internal/agent/approval ·
internal/observe · tools/d22scan · internal/llm/testdata/golden · internal/agent/testdata/golden ·
tools/mockllm/testdata/golden · internal/llm/golden ·
docs/PLAN.md · internal/observe/thresholds.go · tools/d22scan/allowlist.txt · scripts/slo-check.ps1 ·
cmd/wisp/slo_windows.go · .gitattributes            →  17/17 枚"全等=YES"
```
⇒ 这一口径比 §5.3 更强也更弱：**强**在它连"同名字节被换掉而文件名不变"都瞒不过；**弱**在它把区间里别人的 commit 一起算了进来（今天恰好全等，所以两口径互证不冲突；一旦别家程动了某支，它会把别人的账算到本票头上——见 §5.7）。

### 5.7 为什么零字节只能按 commit 量（本节现量，不是转抄派单那条规矩）

```
$ git diff --name-only 1424aa7^..HEAD | wc -l      →  R1（锚点 9835d81）＝43 枚 ／ R3（锚点 be603fd）＝60 枚
      ← 同一条命令、同一对端点，十分钟里差 17 枚：共享树被别人推进了 6 枚 commit。
$ git diff --name-status 1424aa7^..HEAD | grep -E '(frontend|design)/'
    M	frontend/src/main.tsx        ← 归属现量：git log --oneline 1424aa7^..HEAD -- frontend → 902c3857（**owner 的前端会话**，不是本票）
```
⇒ 区间尺今天在本票点名的 13 支上读到 0，但那 0 **不是本票的功劳**：它同时把别人 6 枚 commit 一起判了。r1 派单记的是"区间尺读到 2 枚 frontend 命中、commit 尺 0 ⇒ 假报违规"，本节现量是 **1 枚**（锚点不同、数不同）——**这正是"枚数一律现量"的意思**。判据本体＝§5.3（按 commit）＋§5.6（字节级互证）。

### 5.8 本节自己这一格的写面自查（含共享 index 那一发实量）

- **写面只有两枚**：本文件（追加）＋ `.scratch/wisp/probes/154/**`（只建不删）。本节新增的台件：`ac4-per-face.sh` · `ac4-per-face.txt` · `zb-per-commit-rerun-r2.txt`。
- **commit 前 `git diff --cached --name-only` 两发现量**：12:5x 第一发读到一枚**别人的路径** `docs/evidence/s1/155-three-unjudged-cells-r1-accept-r1.md`（155 验收程 stage 了它）；下一发同一条命令为空（那枚已被它自己以 `edd0725` 收走）。本程**没有** `git add` 它、**没有**替它 commit；提交走显式 pathspec，commit 后 `git show --name-only` 逐枚名册核过＝只有本程路径（在 §6 的 commit 表里连号一起交）。
- **`frontend/**`／`design/**` 不算进任何零命中宣称**：票面 ⚠ 明写。现读工作树里 `design/**` 是 owner 的未提交删除、`frontend/**` 有前端会话的未提交活——本程**不碰、不还原、不 commit**，§5.3 的 13 支里没有它们（r1 脚本的正则里反而有）。
- 同规格的别家半件：`.scratch/wisp/probes/152/my152.py`（` M`）· `docs/evidence/s1/152-…-accept-r1.md`（未提交自校）——同上，不碰，也**不进**本节任何宣称。
- 本节改本文件的方式＝**末尾追加**：commit 前 `git status --porcelain -- docs/evidence/s1/154-…-r1.md` 只应出现 ` M`，且 `git diff --numstat` 给的是 `+N / -0`（＝前 363 行一字未动）；这一发的读数由 §6 那枚 commit 一并带出（本节写作时它还没入库）。

### 5.9 本节没测什么（AC#4 射程内，按"漏了它谁会先被骗"排序）

1. **没盘未跟踪件与 `.gitignore` 后面的东西**：尺是 `git ls-files`／`git show` 的名册，未跟踪文件结构上不在任何一枚 commit 里。若有人把冻结面文件变成未跟踪再改，本节两口径都不响。
2. **§5.6 的字节口径只到"面级点位"17 枚**，没有对 9b 那把宽尺的 58 枚 golden 单件逐一取 hash——那 58 枚的零命中只由 §5.3 的名册口径担保。
3. **没重跑 §1/§3.4 的门禁**（那是 AC#5 的格，已裁）⇒ 本节不声称"改后仍全绿"，只声称"禁改面逐支 0 字节"。
4. **没盘 `internal/tools/**` 自身**：它不在票面 13 支里（本票写面就是它的注释面），它由 §3.3 那条"增删行全部以 `//` 开头"的机器核担保，本节复核了一把新尺（见 §6.2 的 `diff` 两发），但**没有**把 `internal/tools/**` 当成禁改面来判。

### 5.10 AC#4 判语

**成立（本节作者＝续程；勾不勾由编排者按非实现者验收表定，本程不自勾）。** 交付的是：13 支逐支账（§5.3，每支带分母／两发计数／rc／正控）＋ 名册现算 19 枚（§5.1，含"派单命令照字面跑＝1728"那一发反例）＋ 字节口径互证（§5.6）＋ 区间口径为什么不可用（§5.7）＋ r1 那把尺的四处不足（§5.4）。
**边界两句**：本节量的是 r1 的四枚 commit；本程自己那两枚的逐支读数在 §6 收口（同一条尺、输入多一枚 sha）。

---

## 第 6 节 格 AC#6（**续程 r2 补写**）— 承重那句按操作定义答

派单把这一格限定成两问：**Q1 摘掉本票选定的一味，是否存在一发变异从此打不红？Q2 摘掉它有没有任何外部可见读数变过？**
票面 AC#6 自己写明本票实质交付＝"门＋文案＋同族普查"、**不是**新断言，并禁止"为了看起来有牙硬造一发恒真检"。本节按这个边界答，
并且**先把尺说明白**：本仓今天新立的一条——**摘一味导致编不过，不算有牙**（要重做一发改得过、行为不变的对照）。

### 6.1 本票实质交付三件（门＋文案＋同族普查）；两问分两味答，第三味与味 A 同判

| 味 | 本体 | 落在哪 |
|---|---|---|
| **B＝文案** | r1 在 `internal/tools/bridge.go` 新增的 **25 行注释**（§3.2 的 A 段指路句＋B 段"WHICH LEGS STILL HAVE NO OWNER HERE"/"Being first is also an API decision"两段） | 生产码文件里，但**一行都不是语句**（§3.3 的机器核；本节 §6.2 换一把尺再核一遍） |
| **A＝门** | §2.3 那五枚子句 **G1/G1b/G2/G3/G4** ＋ 生成器 `probes/154/gate-clauses.sh` ＋ `bridge.go:665` 指向它的那行注释 | 文档＋`.scratch` 台件，**不进任何二进制** |

（第三味"同族普查"＝§4，纯文档，与味 A 同判：**任何仪器都读不到**，所以不另开一节，见 §6.3 末段。）

### 6.2 味 B（25 行注释）——Q1 与 Q2 都是**当轮量出来的**，不是推出来的

变异体＝**逐字节回退**到 `76c89b2^` 的那版 `bridge.go`（＝摘掉本票那一味、别的什么都不动）。它**编得过**，且行为不变的对照成立在两条现量上：

```
$ git show 76c89b2^:internal/tools/bridge.go > /d/tmp/wisp154-r2/mut/bridge-pre154.go
$ diff <(git show HEAD:internal/tools/bridge.go) /d/tmp/wisp154-r2/mut/bridge-pre154.go | grep -cE '^[<>]'      → 25
$ …同一条 diff | grep -E '^[<>]' | grep -vE '^[<>][[:space:]]*//'                                             → 无输出（rc=1）＝25 行全是注释
$ diff <(git show HEAD:… | awk '/^func \(b \*Bridge\) CloseTask/,/^}/') <(awk '同前' bridge-pre154.go)          → rc=0＝CloseTask 函数体一字未变
$ grep 'C25 scope closed' 两版                                                                                 → 格式串本体逐字节相同（HEAD :693 ／ 变异体 :668，差的就是被摘的 25 行）
$ git log --oneline 76c89b2..HEAD -- internal/tools/bridge.go                                                  → 空（本票之后没人动过这枚文件 ⇒ 回退＝精确摘一味）
```

overlay 通道跑两包（**全程没开 `-cover`**，避开 `-overlay` 被静默忽略那枚坑）：

| 跑 | RUN | 顶格 | 全量裁决 | `^--- FAIL:` | `^--- SKIP:` | `^panic:` | rc | 末行 |
|---|---|---|---|---|---|---|---|---|
| `internal/tools`（变异体） | **115** | **79** | 115 | **0** | **0** | 0 | 0 | `ok … 13.892s` |
| `cmd/wisp`（变异体） | **139** | **79** | 139 | **0** | **0** | 0 | 0 | `ok … 90.848s` |
| 基线＝本文件 §1（改前）与 §3.4（改后） | 115／139 | 79／79 | — | 0 | 0 | 0 | 0 | r1 已入库的两行 |

```
$ comm -23 names-mut-cli.txt  <(sort probes/154/names-post-cli.txt)    → 空      （变异体独有：无）
$ comm -13 names-mut-cli.txt  <(sort probes/154/names-post-cli.txt)    → 空      （r1 改后独有：无）
$ comm -23 / -13（tools 同两发）                                        → 空 ／ 空
（台件：probes/154/{ac6-mutation-strip154.txt, ac6-mut-tools-strip154-v.txt, ac6-mut-cli-strip154-v.txt}；名册枚数 79↔79）
```

⇒ **Q1 答：不存在"摘掉它就打不红"的那一发变异——因为今天根本没有一发红在它身上。** 摘掉 25 行注释后两包 254 条裁决逐字不变。
⇒ **Q2 答：没有任何外部可见读数变过**，四条各自独立：四数（上表）、名册两向 `comm`（皆空）、`go vet ./internal/tools/` rc=0、`sh scripts/d22scan.sh` clean（rc=0；它的射程是 `internal`·`cmd`·`frontend`·`design`·`internal/tools`，**`docs/` 与 `.scratch/` 结构上不在尺上**，现读 `tools/d22scan/main.go:242-252/387-411/559-561`）；外加审计行格式串本体未变（上面那条 grep）。

**仪器正控（不然上面两个 0 只是"没跑到"）**：同一套 overlay 通道、同一枚 `PATH`，摘掉**票 151** 那一味——`cmd/wisp/run.go:563` 的 `rt.bridge.CloseTask(taskID)`（只此一行，`diff` ＝`563d562`）：

```
$ go test -count=1 -v -overlay=…/overlay-noop-close.json -run 'TestCompositionRootClosesTheLoopTasksTaintScope|TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit' ./cmd/wisp/
  第一跑 12:49:50：RUN=2  FAIL=2  PASS=0  rc=1        第二跑 12:54:25：RUN=2  FAIL=2  PASS=0  rc=1
  红句（两跑同一对行号）：task_scope_close_151_test.go:62 / :109 / :113
⇒ 一发"行为真变了"的变异每次都被同一对用例抓住 ⇒ §6.2 那两个 0 是"跑到了、不响"，不是"没跑到"（RUN 115/139 另外钉着"跑到"）。
```

⚠ **树在跑数期间被别人推进过，本程把时间边界写死**：`git status --short -- internal cmd` 于 12:4x 为**空**⇒ 上表两包变异跑（12:47／12:49）编的是干净 HEAD；156 那程的 `cmd/wisp/slo_windows.go` 在 **12:51:23** 被改动、`cmd/wisp/slo_exit_os_156_windows_test.go` 在 **13:01:09** 落盘，所以**只有正控第二跑（12:54）编到了改过的 `slo_windows.go`**，其名册枚数不受影响（那一发只 `-run` 两枚用例）。
另记一枚不是本程的账：`$(go env GOPATH)/bin/gofumpt -l internal/tools cmd/wisp` 本程现读列出一枚件 `cmd\wisp\slo_exit_os_156_windows_test.go`（**156 程未提交的活**）——本程不判、不碰、不格式化它，只把读数交给编排者。

### 6.3 味 A（门本体）——两问：都不响，而且是"结构上没人读它"

```
$ git grep -nE "gate-clauses|probes/154|154-host-id" -- '*.go' '*.yml' '*.ps1' '*.sh'      （尺：大小写敏感、子串、pathspec 只放这四种扩展名、不排除任何目录；跑两遍，输出逐字相同，rc=0）
  .scratch/wisp/probes/154/census-pairs.sh:2        ← 台件自己的用法注释
  .scratch/wisp/probes/154/gate-clauses.sh:2        ← 同
  .scratch/wisp/probes/154/zero-byte-per-commit.sh:2 ← 同
  internal/tools/bridge.go:665                      ← 注释里那行"门在哪"的指针
⇒ *_test.go 里 0 命中（单独一发：git grep -n "154-host-id-never-closed" -- '*_test.go' → rc=1）；.github/workflows/*.yml 与 scripts/* 里 0 命中。
```
⇒ **Q1（味 A）答**：摘掉门（删 §2.3 那张表＋删 `gate-clauses.sh`）⇒ **零枚用例转红**，因为没有任何仪器读它（上条就是证据）。
⇒ **Q2（味 A）答**：**没有任何外部可见读数变过**（d22scan 与 vet 的射程都不含 `docs/`·`.scratch/`，见 §6.2 末）。唯一会变的是 `bridge.go:665` 那行注释的**指针悬空**；而"注释指向不存在的落点"这件事，本仓**没有一发仪器会响**——现读 `git grep -n "docs/evidence" -- scripts tools .github` 命中 9 行，逐枚点开全是散文引用（注释／登记表字符串），**没有一处是校验路径存在性的检**。
⇒ 一处**不能说过头**的地方（本仓有反例）：**"证据件里没有测试读者"这句在本仓不成立**——`internal/ball/tokens_table_test.go:54` 与 `internal/panel/tokens_fourway_test.go:51` 就把 `docs/evidence/s1/c21-native-tokens.md` 当输入路径读（现读 `git grep -n '"docs/' -- '*_test.go'` ⇒ 3 命中，全在那两枚）。所以味 A 的不响是**"没有测试读*这一枚*文件"**，不是"测试从不读 docs"。
⇒ 味 A 的第一发读数原样留档：本程自己把 pathspec 写成 `:!_scratch` 撞了 `fatal: Unimplemented pathspec magic '_'`，**真 rc=128**（不是 0 命中）；那一发若按"只数命中行数"的写法就会被写成"全仓零命中"。台账与本件都记过同一类（`git grep` 单发吞命中），本程又踩了一次，尺面与更正见 `probes/154/ac6-gate-reader-scan.txt` 末段。

### 6.4 派单点名"必须写死"的三件

**① 它今天不响；买的是防忘记，不是防回归。**
今日不响的可复算理由（全部指回 §2.2／§2.3，**不新造子句**）：G1／G1b／G3／G4 今日名册皆空；G2 今日就是那 2 行（`cmd/wisp/run.go:546` 注释＋`:563` 唯一生产调用），那就是"关那侧只有一枚所有者"这一事实的基线。
本程把五枚子句在自己的锚点上复跑了两遍（12:5x，锚点 `67cc43f`，两遍输出除时刻行外逐字相同），并把同一把尺的第三发入库为 `probes/154/gate-clauses-roster-r2.txt`（锚点 `be603fd`）：`# rc=` 序列＝`1 / 1 / 0 / 1 / 1`，且名册内容与 r1 在 `1424aa7` 那份 `gate-clauses-roster.txt` **逐字相同**（剥掉锚点前缀后 `diff` ⇒ rc=0）。
它买的是防忘记：§6.2 与 §6.3 已经量过——摘掉两味，码面与测试面零变化。变化的只有"下一位读不读得到这句话"。

**② 哪一形真出现时它会响（可复算条件，从 G1/G1b/G2/G3/G4 里挑）：**

| 触发形 | 响的子句 | 复算命令（同一条，两个锚点各跑一次、逐行 diff 名册） | 从"不响"变"响"的判据 |
|---|---|---|---|
| 有人在生产码里自己构造过桥的派发（＝宿主自带 task id 那一形出现） | **G1**（与 **G1b** 同发即响） | `bash .scratch/wisp/probes/154/gate-clauses.sh $(git rev-parse HEAD)` | G1 段从 `rc=1`（空）变 `rc=0` 且出现一行 `文件:行号:` ⇒ **本票 §2.4-3 说的那件事发生**：响的是"名册非空"这个事实，义务（给那一形接上关闭＋补并发读数）仍要人做 |
| `Q-56` 那一支被落地（`Loop` 上出现收外部 task id 的导出方法） | **G3** | 同上 | G3 段非空 ⇒ 那一形从"没有入口"变"有入口"，**关闭的所有者必须同批改**（§2.3 表 G3 行） |
| 关那侧多出一枚所有者，或开那侧多出一枚生产调用者 | **G2** | 同上 | 那 2 行变长 ⇒ "151 只接了一形"这句话过期，回头核它接得对不对、§2.5 并发读数补没补 |
| 生产组合根打开 L0 直通（环路那一形端到端可见） | **G4** | 同上 | `PassThroughUnclassifiedRisk:` 在生产码里非空 ⇒ "151 已把环路那一形接完"第一次有端到端读数（§4.2 的 L-3 被摘） |

**③ 这一格勾不勾留给谁**：**编排者按非实现者验收表定**（票面"票面框由编排者按**非实现者**验收表定，实现方不自勾"）。本节只交两问的读数与上表那条可复算条件；本程**不自勾 AC#6，也不新增格子**。

### 6.5 文案那一味的"可见面"——派单要求：这句只有指得出"谁会误以为"、带出处时才成立

派单原话：注释没有运行时读数，所以老实写"无外部可见读数、它的可见面是**读者会不会误以为那条腿已修**"，而**这句要带出处**。本程现量三条出处，外加一条反向自查：

1. `docs/evidence/s1/151-task-scope-never-closed-r1.md:298`（第 8 节"本程没测什么"第 1 条末句，逐字）：
   > 谁先被骗：验收"R4 接线完成"的那一格——它会以为 §2 的问题被本票关掉了。

   同条上一句点名受害方：「球/面板宿主（票 33/35/77/92）一上线，§2.1 那发硬拒/误报**照样发作**」。四枚票在本目录都在（`ls .scratch/wisp/issues | grep -E '^(33|35|77|92)-'` ⇒ 四枚命中）。
2. `docs/evidence/s1/151-task-scope-never-closed-r1-accept-r1.md:485`（**非实现者**那张表的第 10 节"本程没测什么"第 5 条，逐字）：
   > 球／面板宿主那一形没起。…… ⇒ 这条只有等宿主真起来。
3. 同件 `:314`（第 6 节，它对实现件 §8 第 1 条的复核）判词＝「**登记准**」，并给出链路：`CloseTask` 全树唯一调用点 `run.go:563` → `admitTask` → `Options.AdmitTask`（`run.go:586`）→ `loop.go:361-362` ⇒ 不经环路的宿主 id 没有任何一处会关它。**这三处行号本程全部现量**（`grep -n "AdmitTask" cmd/wisp/run.go internal/agent/loop.go` ⇒ `run.go:586`、`loop.go:161/169/361/362/777`；`Loop.Run`＝`loop.go:332`、`RunAsync`＝`:321`、`newTaskID`＝`:1107`；审计行＝`bridge.go:693`），不是抄 151 那张表在它锚点下的数。

**那一枚被点名的验收格今天还不存在于任何已入库表里**——这条本程也量了，免得本节说过头：
`git grep -n "R4 接线" 79551d1 -- '*.md'` ⇒ 已入库的树里 **2 行**命中：一枚是上面 §8:298 的原话，一枚是票 154 票面 `:25` 抄它。
⚠ 这一发**必须钉 sha**：不带 sha 在当前工作树上跑同一把尺，读到的是 **4 行**＝上面那 2 行 ＋ 本件自己引用它的 2 处（§6.5 那条引文与本条判据）——**读者会污染自己的尺**，这是本仓"零命中要写明哪把尺"那条的又一形。
⇒ 所以"谁会误以为"精确到形＝**未来那张表**（面板/球宿主接线时的验收格），不是"现在某张表的某一行已被骗"。本节按这个措辞写，不写成既有格的既成事实。

**反向自查（会不会其实没人会误以为？）**：宿主票自己的票面今天**一句都没谈** scope 关闭——
`grep -n -E "R4|CloseTask|OpenTask|scope" .scratch/wisp/issues/35-panel-bridge-c17.md` ⇒ 唯一命中 `:38` 的 "## Out of scope"（一节标题，不是判据）；
`grep -n -E "R4|CloseTask|OpenTask|scope|task ?id" .scratch/wisp/issues/145-panel-snapshot-….md` ⇒ **0 命中**（尺：大小写敏感、子串、整枚票面文件，两发一致）。
⇒ 这恰恰是"误以为"能成立的原因：票 35／145 的读者在**自己票面**上撞不到任何提醒，只有真去 `bridge.go` 接那一腿、读到 `CloseTask` 注释时才会知道。**所以票面 AC#2 要的那个"会被读到的地方"，三件交付里只有注释这一件今天真的在读者会经过的位置上**（门在证据件、普查在证据件）。
⇒ 票面 AC#2 那条"注释里不许预先引用尚未产出的读数"：本节与 §3.2 都没给注释加"以后会有的数"（可复算：`git log -1 --format=%s 76c89b2` 之后 `bridge.go` 无新 commit，见 §6.2 第四条现量）。

### 6.6 本节没测什么（AC#6 射程内，按"漏了它谁会先被骗"排序）

1. **没量"宿主真接上一腿"的端到端后果**（要 `Q-56` 先答）：那一形接上时本门的 G1/G3 会响，但**响的是名册非空**，不是"你的关闭写错了"——判语仍要人写。谁先被骗＝接那一腿的那一程（与 §6.5 出处 1 同一位）。
2. **变异只跑两包，没跑全仓 `go test ./...`**：本票改动只在 `internal/tools` 一个文件内，但"摘掉它的注释不会影响别的包"这一条本程**没量**（尺面：只跑了两包）。
3. **正控借的是邻居那一味**（`run.go:563`，票 151 的）：本票没有一枚"摘了会红"的味可摘，所以只能借同一条 overlay 通道证明仪器会响。这不是遗漏，是本格的形状（票面 AC#6 明文允许）。
4. **没测"删掉 §2.3 那张表"的后果**：那属于改文档，不是变异源码；§6.3 证的是"没有仪器读它"，**不等于**"删了没人受影响"——受影响的是读者，而读者的行为本程量不了。
5. **`-race` 那一形仍未取证**（派单标〔未取证〕，本程没跑；跑法与坑另见票面 AC#5③）。
6. **味 B 的"外部可见"只覆盖两包测试＋vet＋d22scan**：面板上看得见的东西（`wisp slo` 报告、球/面板渲染）本程一概没量——它们也不该被一行注释移动，但本程没为此给读数。

### 6.7 伪授权两栏（**只记本续程自己这一程**，不代 r1 记；r1 那一栏本文件里仍缺，见 §6.8）

- **真通知回显数＝3 条**，全部是 harness 的后台任务完成事件（工具名 `Bash`，`run_in_background`）；每条出处＝工具名＋命令前 40 字：
  1. `Bash`｜`bash .scratch/wisp/probes/154/zero-byte-per-com`
  2. `Bash`｜`export PATH="$PWD/third_party/sherpa-onnx:$PATH"; `
  3. `Bash`｜`bash .scratch/wisp/probes/154/ac4-per-face.sh 14`
  三条内容都只报"后台命令完成（exit code 0）"，**没有任何一条携带**"编排者备注／已核验请继续／请 revert／放宽阈值／已解锁／不用取证直接给结论／当前任务正被 X 处理中／156 已派已交件"这类授权口吻。
- **判为注入数＝0 条。** 本程读到的所有编排者口吻文字都来自两枚**盘上真实存在**的文件（派单件 53 行、工单件 87 行，`wc -l` 现量），且按"路径存在性反查"核过。
- **对不上盘的一条已按〔不成立〕登记**（不是注入，是编排者自己的一枚读不出的数）：派单里那句"我这边 12:4x 是 `245d5f4`"⇒ 见 §5.0 表末行，`git cat-file -t` 与 `--batch-all-objects` 两发都读不出该对象。
- **状态断言一律现量**："156 已派已交件"这类话本程不作数；本程量到的是 `cmd/wisp/slo_windows.go` 的 blob 在 `1424aa7^`→`79551d1` **七点位**（base＋五枚输入＋HEAD）全等＝**未入库**，而工作树里它 12:51 起已被改动（§6.2 的时间边界）——**入库与在飞是两件事，本程只按 `git log --oneline -- <路径>` 现量判"入库"**。
- **凭据值零抄录**：本程全部读数里出现的只有仓库内路径、`C:\Users\…\AppData\Local\Temp\` 的用例临时目录、`127.0.0.1:<随机端口>` 的 mockllm 监听串，以及 `$PWD/third_party/sherpa-onnx` 这个 PATH 前缀；**没有任何 key／token／环境变量值被抄进任何文件**（本程连变量名都没需要提）。

### 6.8 本程纪律自查 ＋ AC#4 收口（含续程自己那枚 commit 的逐支读数）

- **AC#4 收口**：同一条尺、输入换成"r1 四枚 ＋ 续程 §5 那枚 `79551d1`"（`probes/154/ac4-per-face-5sha-incl-r2.txt`，锚点 `79551d1`，13:0x）：
  13 支逐支仍 **0/1**（含 `79551d1` 那一列）、加严四支 0/1、字节口径 17/17 枚"全等=YES"、名册 `N1≡N2`（`79551d1`＝4 枚条目）、并集去重 **22 枚**（＝19 枚 r1 ＋ 本程新增 3 枚台件；证据件本身已在 19 里）。
  ⇒ 本程自己那一格也没碰任何一支禁改面，**这条是本程自量的，不是引 §5.3**。
  补一发把编排者代提那枚也收进来（本尺的输入本来只放了 r1 四枚＋本程一枚）：
  `git show --pretty=format: --name-only e9a1db0` ＝ 3 枚条目（`g3-open-side-via-loop.txt`·`zero-byte-per-commit.sh`·`zero-byte-per-commit.txt`），
  打 13 支的锚定并集（golden 用窄尺）⇒ **N1 命中 0／N2 命中 0** ⇒ 票 154 名下六枚 commit（`1424aa7`·`352d3d8`·`76c89b2`·`304199f`·`e9a1db0`·`79551d1`）**逐支全 0**。
  ⚠ **本格宣称的边界（第 13 支正踩着一枚别家在飞的活，必须写清）**：AC#4 量的是**已入库的 commit**；工作树里 `cmd/wisp/slo_windows.go` 自 12:51:23 起被 156 那程改过、`cmd/wisp/slo_exit_os_156_windows_test.go` 自 13:01:09 起以未跟踪状态存在（§6.2 的时间边界）。这些**既不在本格的零字节宣称里，也不被本格否认**——本格只说"票 154 名下的六枚 commit 一支都没碰"。
- **每裁一格 commit 一次**：`79551d1`＝第 5 节（4 枚路径：本文件＋`ac4-per-face.sh`＋`ac4-per-face.txt`＋`zb-per-commit-rerun-r2.txt`，`git show --numstat` ＝`+157/0`·`+207/0`·`+248/0`·`+42/0` ⇒ 前 363 行一字未动）；本节＝第二枚，名册＝本文件＋`probes/154/` 下 8 枚新增台件（`ac6-*`×6、`gate-clauses-roster-r2.txt`、`ac4-per-face-5sha-incl-r2.txt`）。本节 commit 自己的 sha 由**下一位**现读：`git log --oneline -- docs/evidence/s1/154-host-id-never-closed-r1.md`。
- **显式 pathspec**：两枚 commit 都是 `git commit -q -F - -- <逐枚列出路径>`；无 `git add -A`／`git add .`；`git diff --cached --name-only` 在每次 commit 前现读（12:5x 那一发读到一枚别家路径、下一发已空，见 §5.8）。
- 禁件一律未用：`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean` 一次都没跑；**push 一次都没做**（`git log --oneline origin/dev..HEAD` 由编排者自己量，本程不替他断言推送状态）。
- 临时件全在仓库外 `D:/tmp/wisp154-r2/**`（含两枚变异体源文件与四发原始 `-v` 输出），**只建不删**；仓内未建 worktree／checkout。
- 取干净树一律**没有**用 `git archive | tar -x`（§5/§6 全部用 `git show`／`git rev-parse <c>:<path>`／`git ls-files`）；本程**用了 `-overlay` 但全程没开 `-cover`** ⇒ 那枚"overlay 被静默忽略"的坑没被触发，且每一跑都以 `^=== RUN` 115／139 证明"真跑到"（不是靠 `ok` 那一行）。
- **本件仍缺的两节**（票面"交件形状"要求的）：`本程没测什么` 的 r1 整件版（本程已在 §5.9／§6.6 各补了本节射程内的版本，但**没代 r1 补全件版**）与 r1 的伪授权两栏（§6.7 只覆盖续程自己）。本程写面只有这一枚文件＋`probes/154/`，且派单限定"只追加第 5、6 节"⇒ 这两节**留给谁**由编排者定，本程不自作主张补第 7 节。

### 6.9 AC#6 判语

**按票面 AC#6 的操作定义答完＝两问各自有当轮读数：Q1 两味都"没有一发红在它身上"（味 B 有变异体两包 254 条裁决逐字不变的量、味 A 有"没有任何仪器读它"的量）；Q2 两味都"无任何外部可见读数变过"（四数／名册／vet／d22scan／审计行格式串）。**
⇒ 这是一枚**诚实的"今天不响"**：它有正控（摘 `run.go:563` 两跑都红）、有可复算的响法（G1/G1b/G3/G2/G4 各一条判据）、有指得出处的可见面（§6.5 三条出处＋一条反向自查），**没有**为"看起来有牙"新增任何恒真检。
⇒ **本程不自勾**：AC#6 这一格（以及 AC#4 那一格）翻不翻，由编排者按非实现者验收表定。

### 6.10 自报违约：第 5 节入库后本程又**就地**改过它六处——原文不抹，本块代它负责

写第 6 节时反查发现自己 §5 里有六处数或方向不对，而 §5 已经以 `79551d1` 入库。本仓的规矩是**证据件只追加**（票面"见台账 `A278②`"；台账 `A280` 那块 `>` 立的例：已提交的行只在就近补一块，由新块代它负责）。本程一开始直接在文件里改了那六行——**那一步做法不对**；本节把文件恢复成"`79551d1` 的 §0–§5 逐字不动 ＋ §6 追加"的形状，并在此登记被覆盖处。取回第 5 节入库原文：`git show 79551d1:docs/evidence/s1/154-host-id-never-closed-r1.md`（第 367–520 行）。

| # | §5 里那句（入库原文仍在） | 现量更正 | 复算命令 |
|---|---|---|---|
| E1 | §5.0 首段"往前走了 **6 枚 commit**（155 验收程两枚、台账两枚、另两枚）" | **7 枚**＝155 验收程 6 枚 ＋ 台账 1 枚 | `git rev-list --count 9835d81..be603fd`；`git log --format='%h %s' 9835d81..be603fd` |
| E2 | §5.0 末行"155 那两枚已入库（`67cc43f`／`edd0725`）" | 那一刻读到 2 枚，同节写作中到 `be603fd` 已 **6** 枚 ⇒ **枚数在涨，本节不该写总数**（本仓今天同类已三次：姊妹程 7→11→12、13→15） | 同 E1 |
| E3 | §5.2 末句"一支是被第 6 支包含的父树" | **方向写反**：是**第 5 支包含第 6 支**（`internal/agent/**` 58 枚 ⊃ `internal/agent/approval/**` 18 枚） | `git ls-files -- internal/agent \| wc -l`＝58；`… -- internal/agent/approval \| wc -l`＝18 |
| E4 | §5.3 尺面句"每一发都跑两遍取数并记 `grep` 的 rc" | 措辞不准：实际是**每格两把名册尺**（N1 折叠改名／N2 拆开改名两侧）各数一次，不等则输出里标 `a!=b`——今天 13×4＝52 格无一处不等；"同一发跑两遍"这一层由 **R1／R3 两枚锚点整表 `diff`** 担保 | `grep -c '!=' probes/154/ac4-per-face.txt`＝0；R1↔R3 的 diff 见同文件末段（只一行不同＝§5.7 的区间数） |
| E5 | §5.4 末段"正则实为 18 项，并集 15 支，多出的是那四支加严" | 对不齐是**双向**的：15 支在它尺上只占 **14 个条目**（`internal/agent/**` 与 `internal/agent/approval/**` 共用一个前缀），所以 18＝14＋4。**别把 18−15 读成 3 条加严**，也别把 15 读成票面数（票面是 **13**） | 展开 `zero-byte-per-commit.sh:6` 那行 `FREEZE` 逐条点数 |
| E6 | §5.7 末句"它同时把别人 6 枚 commit 一起判了" | 现量：`1424aa7^..be603fd` 共 **28** 枚，票 154 名下 **5** 枚（`1424aa7`·`352d3d8`·`76c89b2`·`304199f`·`e9a1db0`），**别家 23** 枚 | `git rev-list --count 1424aa7^..be603fd`；`git log --format='%h %s' 1424aa7^..be603fd \| grep -cE '票 154\|代 154'` |

三句给下一位的实话：①这六处**没有一处翻转 §5 的判语**（13 支逐支 0、名册 19 枚、区间尺不可用，三件事各自独立复算过）；②但 E1/E2/E6 是同一类错——**把"那一刻读到的枚数"写成了永久事实**，本仓今天为这一类已经退回过程序；③正确做法是**发现错时那一格还没 commit 就一次写对**，或已 commit 就只追加不就地改，本程两样都晚了半步，登记在此。
（另交代一笔过程：本节 §6.8 里"六点位→七点位"那一处修订发生在它自己入库**之前**，属同一枚未入库草稿内的改稿，不属上面这类违约。）

---

> **【编排者代记 09-26 13:5x｜下面这 87 行是 r1 程（`a701c4f9ac4066902`／transcript `a58c67b49f7c683d`）在 13:50:44 迟到写出、只 stage 未 commit 的那一段；本程按"代提只到读数落盘为止"把它入库，判语一个字不替它改。】**
> 撞号说明（不抹、不改它的标题）：本块的两节与上面已入库的 r2 版两节**同号不同物**——本块的"第 5 节 AC#4"＝4 枚 commit 的粗口径（`zero-byte-per-commit.sh`），已入库的 r2 版"第 5 节 AC#4"＝13 支 × 4 枚的逐支账；本块的"第 6 节 AC#6"是单味粗答，r2 版是味 A／味 B 分答。
> **哪一版 govern**＝r2 那两节（先入库、口径更细、且经上面 §6.10 逐处自我更正过）；本块那两节按"同格双份表＝第二见证，不抵账"处理，AC#4／AC#6 的判语**不从本块取**。
> 本块**不与已入库内容重复的增量**只有两节：**第 7 节**（票级"本程没测什么"，r2 的 §5.9／§6.6 只到单格射程）与**第 8 节**（r1 自己的伪授权两栏——上面 §6.7 明写"r1 那一栏本文件里仍缺"，缺的就是这段）。
> 另记一笔已同步台账的自报：本块 §8.3 第 1 条承认违反"禁 `rm`"（`rm -rf /d/tmp/wisp154-gatectl`），并自称该目录当时不存在⇒**"没有件被删"只有它自己的话、盘上不可反查**，故按〔仅自述〕入账，编排者不据它写"无损"。

## 第 5 节 格 AC#4 — 契约轴零字节（**按 commit 量**，派单点名的坑）

尺与射程＝`.scratch/wisp/probes/154/zero-byte-per-commit.sh`（冻结面＝派单 AC#4 那串 ∪ 票面 AC#4 那串，再加 `.gitattributes`／`go.sum`／`frontend/`／`design/` 四枚别的会话在动的路径）；
输出＝`probes/154/zero-byte-per-commit.txt`。本票到此为止共 4 枚 commit，逐枚列全量文件清单＋冻结面命中数：

| 本程 commit | 碰的文件 | 冻结面命中 |
|---|---|---|
| `1424aa7` 第 0–1 节 | 5 枚 `probes/154/*` ＋ 本件 | **0** |
| `352d3d8` 第 2 节 | 5 枚 `probes/154/*` ＋ 本件 | **0** |
| `76c89b2` 第 3 节 | 5 枚 `probes/154/*` ＋ 本件 ＋ **`internal/tools/bridge.go`** | **0**（`bridge.go` 是派单写面，不在冻结轴上；它的"只到注释"由 §3.3 那条机器核担保：`+25 / -0`、增删行全部以 `//` 开头） |
| `304199f` 第 4 节 | 2 枚 `probes/154/*` ＋ 本件 | **0** |

⇒ **AC#4 成立**（四枚皆 0）。三条限定语，一条都不省：
① 这是**按 commit** 量的（`git show --name-only`），**没有**按 `HEAD..` 区间量——本仓 09-26 实发过区间尺读到别人的 `frontend/**` 命中而假报违规；
② 别的会话在本程干活期间往 index 里批量 stage 过东西（本程看见 `git diff --cached --name-only` 一度 11 KB），**本程没有动过 index**（禁 `reset`），每枚 commit 一律带显式 pathspec，因此他们的 staged 件没被本程带走——`git show --name-only` 那四行清单就是这条的反查；
③ `frontend/**`／`design/**` 的未提交状态**没有**被本程算进任何"零命中"宣称（上面那把尺把它们整批排除，本程只在"我也没碰"的意义上写 0）。

---

## 第 6 节 格 AC#6 — 承重那句按操作定义答

**本票选定的一味**＝`internal/tools/bridge.go:655-674` 那两段新增注释（第 3.2 节 B 段中段）。

1. **摘掉它，是否存在一发变异从此打不红？** 答：**摘掉它，零枚用例转红**——Go 的注释不进 AST 输出、不进二进制。
   本票**不**为新断言交账（票面 AC#6 明文：本票实质交付是"门＋文案＋同族普查"）。
   因此本票把该写死的那句写死在这里：**这道门今天不响；它买的是防忘记，不是防回归。**
2. **它今天为什么不响、哪一形出现时它会响**：写在了**门本体**上，不在本件里自说自话——`bridge.go:663-666`（"今天生产没有这种调用者 ＋ 该响的是 §2.3 的 G1/G1b ＋ 它就是一条 `git grep` 的距离"）与 `:668-674`（"那是一个 API 决定＝`Q-56`，本文件不答它 ＋ 并发读数仍是零，同一批改里要一起补"）。
   响的**可复算条件**＝第 2.3 节那张表的第四列，任一枚子句的名册非空/变长即响，且**每枚都已过仓库外正控**（`gate-positive-control-r2.txt`）。
3. **另问的那一句：摘掉它有没有任何外部可见读数变过？** 答：**没有**，三条独立证据在第 3.4 节：四数逐字同（115/79/0/0 与 139/79/0/0）、两包名册两向 `comm` 皆空、`go vet` rc=0／`gofumpt -l` 空（v0.12.0 (go1.27.1) 现读）／`d22scan` clean 且各作用域分母非零；审计行的格式串（`bridge.go:693`，`tools: C25 scope closed task=… was_open=… dropped=… open_scopes=…`）一字未动。
4. **反向自查（别把自己说圆）**：这枚文案的射程只到"**读注释的人**"。
   一个只 `grep -n "CloseTask(" ` 找调用者、从不读函数头的下一位，本票**买不到**他——那一形只能靠 G1/G1b 门响去接，而门今天不在 CI 里（`tools/d22scan/**` 是冻结面，本票不许把门塞进去）。
   ⇒ 这是**已知残余**，登记在第 7 节第 1 条，不在本票修。

---

## 第 7 节 本程没测什么（按"漏了它谁会先被骗"排序）

1. **没有端到端跑出"环路那一形在生产上真开过一枚 scope"的读数。** 被挡住的是三枚上游锁（§4.2 的 L-1/L-2/L-3）＋"今天只有 `fs.read` 能被 Mark"（L-4）。
   **先被骗的人**＝读"票 151 档位＝成立／已接线"那句话的人（`T1` 那枚用例绿的是审计行出现，不是污点被关掉过）。
   复算入口：`git grep -n "PassThroughUnclassifiedRisk" <锚> -- cmd/wisp/run.go`（零命中）＋ `grep -n -A 4 "func FSReadDecl" internal/tools/fs.go`（`Declared: risk.L0`）。
2. **没有把这枚门接进任何仪器。** 它今天是"一条命令＋一份名册基准"，不在 CI、不在 `d22scan`、不在用例里（都不许：§2.4 第 3 条）。⇒ 门会不会被跑，取决于下一位读不读注释。
3. **没有 `-count=2`、没有 `-race`。** 本程四跑全是 `-count=1`；派单标〔未取证〕的那枚"`-race` 在本机可能 `0xc0000374`"本程**没有验证**，也没用它。
4. **没有在纯净树（真 checkout 到仓库外）上复跑门禁。** 本程跑在共享工作树上，凭据＝§0 那条 `git diff --name-status d0865ff..HEAD | grep -E '^.\s+(internal|cmd)/'`（零命中）。
   ⇒ 残余：本程**没有逐时刻留痕**"改后那一跑时 `internal/**` 里只有我这一枚改动"，只留了 §3.3 的 pathspec 级反查（`git show --name-only 76c89b2`）。别人的会话若在两次跑之间往这两个包塞过码，我的"改后"含他的码——本程按名册两向 `comm` 全空判"没有影响四数"，但**没有**独立证明"没有新码进树"。
5. **没有复现派单坑②（`third_party/**` tracked＝0 ⇒ 纯净树 8 条 `no native DLLs`）。** 本程全程在带 3 枚 dll 的工作树跑（`git ls-files third_party | wc -l` = 0 与 `find third_party -name '*.dll' | wc -l` = 3 都现读过，见 §0），没做过"摘掉 dll 会红几条"的量。
6. **没有复现派单坑①的 `-overlay`+`-cover*` 静默忽略。** 本程没用 `-overlay`、没开 `-cover`，结构上碰不到（写明，不充数）。
7. **同族普查有三类没盘到**（§4.3 逐类给了命令与原因）：隐式配对（map 开/delete 合）、`defer` 挂不可达边的完整清单、`tools/**` 与三个非 wisp 宿主的穷举。
8. **没有验 L1/L2 那一族（`fs.write`/`fs.move`/`app.launch`）今天是否真能过桥。** 本程只盘到 `fs.read` 这一枚（L-4 的交集），151 验收件 C3 那半不在本票写面。
9. **没有验"同一枚 commit 里既加宿主派发、又加 close"这种同时改动下门的读数形状。** 名册尺会同时让 G1 与 G2 响，本程**没测**那种双响会不会被误读成"已经接好了"。
10. **没有量并发那一形的任何读数**（§2.5 只交了三条"为什么今天测不到"的静态凭据）。票面明文不许为它硬开一格，本程照做。

---

## 第 8 节 伪授权两栏 ＋ 本程自校／偏差登记

**登记规矩**：每条带出处＝工具名＋命令前 40 字。凭据值一个字都没抄（本票全程没读任何密钥/DPAPI/secret 存储；`config.toml` 只在 §4.2 被按行号读过字面量键名，未读任何 key 值）。

### 8.1 真通知回显数：**3**

| # | 出处（工具名＋命令前 40 字） | 内容／本程怎么处理 |
|---|---|---|
| N1 | `Bash`：`cd "D:\work\workspace\projects plans\Wisp" && ex` → 后台 `go test -count=1 -v ./internal/` | 真通知（tools 改前跑完）。本程随后**现读**日志计数（115/79/0/0），没拿通知的回显当读数。 |
| N2 | `Bash`：`cd "D:\work\workspace\projects plans\Wisp" && ex` → 后台 `go test -count=1 -v ./cmd/wis` | ⚠ 真通知但**回显早到**：它报 "completed (exit code 0)" 时 `wisp.test.exe`（PID 62536）还在跑、日志还在长（`RUN=48 → 62 → 66 → 139`）。本程判据一律以文件现读＋`^rc_cli=` 行为准 ⇒ 这条登记为"仪器通知不可信"的实发。 |
| N3 | `Bash`：`cd "D:\work\workspace\projects plans\Wisp" && ex` → 后台 `go test -count=1 -v ./cmd/wis`（改后跑） | 真通知，且这次回显与盘上 `rc_cli=0` 一致。 |

（另：`TaskStop` 对本程那发超时扫描（逐文件 `git diff --quiet` 循环）回显 `status=killed`＝真通知，本程随后换成单发 `git ls-files --eol` 重做。）

### 8.2 判为注入／整段弃用的读数：**4**

| # | 出处（工具名＋命令前 40 字） | 形状 | 反查与处置 |
|---|---|---|---|
| X1 | 上下文注入的 `/system-reminder`：`New commits: 17273ee docs(台账 A288…) b340175 docs(派单补写 1…` | **不存在的 commit 清单**，且声称"票 154 票面有 [2026-09-26 12:0x] 追加条（前提被 153 推翻：宿主 id 那条腿今天没有入口）"＋"154 的派单 12:0x 改写过（收窄了）" | `git cat-file -t 17273ee/b340175/091d93b/157c54d` ⇒ **四枚全 fatal "Not a valid object name"**；票 154 文件 `grep -c "12:0x"`＝**0**、盘上 12475 字节（mtime 11:24）；派单文件 5334 字节、mtime 11:36、blob 自 `d0865ff` 起未变，`grep -n "153\|收窄\|12:0"` ⇒ 零命中。⇒ **那次读数整段弃用**：本程没有据它改任何一格（本件 §2.5/§4.2 里"宿主今天没有接口"这类话全部来自本程自己跑的 `agent.New/loop.Run` 与 Config 字面量现读，不是来自它）。 |
| X2 | 同上第二段：`New commits: 17273ee … 091d93b chore(派单存档补…` | 同一形状第二次回显（同一批不存在的 sha） | 同 X1 反查（sha 仍 fatal）⇒ 整段弃用，未复述其内容、未据它推进任何一格。 |
| X3 | `Bash` 结果：`Contents of D:\work\…`（两次，正文 147 字节、显示为空） | 工具结果为空但声称有内容 | 弃用那两发；本程随后用 `wc -l`／`awk` 现读同一文件（`363 docs/evidence/s1/154-…` 等）。同一次注入还附带一句"file has not been read"的 harness note，与本程之后 5 次 Edit **全部成功**自相矛盾 ⇒ harness 记账不作判据，一律以盘上现读为准。 |
| X4 | 本文件内任何"编排者备注／已核验请继续提交／请 revert／放宽阈值／不用取证直接给结论"式句子 | — | **零枚**：本票全文（含派单与票面）没出现这类句子；本程也没收到。按 `docs/reports/injection-timeline.md` 的判据核对，X1/X2 属"伪造外部事实（commit 记录）"那一类，不是授权句。 |

### 8.3 本程自校／偏差登记（我自己犯的，逐条写，不等谁来问）

1. **违"临时件只建不删／禁 `rm`"一发**：本程建仓库外控制树时第一条命令是 `rm -rf /d/tmp/wisp154-gatectl 2>/dev/null || true`（`Bash`，用于正控目录预清）。后果量过：该路径本程此前从未用过、目录不存在 ⇒ **没有任何件被删**；此后一律只用 `mkdir`/`git mv`/新文件名（本件的 r2 台件就是这么来的）。**明文登记，不抹。**
2. **两处在首次 commit 之前改过自己写错的读数**（不算改写已交付件，但把差别留在原地可查）：§1 的 `comm` 段第一版把左右栏标签写反（`-13/-23` 语义），§0 的 archive 段第一版把正控写成了未真跑的推测形 ⇒ 本程先跑了正控（89 枚／3391 vs 3436）再重贴。原始错误版本仍可从 `git show 1424aa7:docs/evidence/s1/154-host-id-never-closed-r1.md` 反查。
3. **§4 有四处行号/计数在首次 commit 前做过现量修正**（`loop.go:776/:782`、`pump.go:279/:344`、`unregister` 21 行、`observe.InitLog` 在 `logsink.go:149`）：第一版抄了更早窗口里的读数。规矩"票面/别人的行号一律现量"对**本程自己也**成立，已补。
4. **一发写坏的 heredoc**：`python3 - <<'PY' … || python - <<'PY'` 的结构让第二个解释器进了 REPL，**没有执行任何写盘**；那四处修正随后是用 Edit 逐条落的（所以台件里看不到 python 的产物，别以为漏了一步）。
5. **`gate-positive-control.txt`（第一发）里 G3 那枚正控是空的**，本程没删它、也没改它——`git mv` 探针文件之后另存 r2。空的那发留着，是因为"正控为空时先怀疑自己的探针放错地方"这条教训只有留着才可复算。
6. **本票的三门（G1/G1b/G3/G4）今日全空、G2 今日就是那 2 行**——本程没有为了让它看起来有牙去加任何恒真检；票面 AC#1 与 AC#6 的这两条禁止，本程按字面执行了。
