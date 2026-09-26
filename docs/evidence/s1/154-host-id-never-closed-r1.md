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
