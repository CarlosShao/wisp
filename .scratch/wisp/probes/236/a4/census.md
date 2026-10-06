# 236-a4 · AC#3 第一半静态普查：今天有几枚尺会在「task id 换成模型给的值」时变红

> 本腿＝票 236 AC#3 **第一半**（"先量清"）的只读普查腿 `236-a4`。⛔ 零 `go` 命令、零产码改动、零票面框改动、零 `docs/**` 改动。
> 原始读数全部落 `D:/tmp/wisp236a4/logs/`（绝对前缀临时件，只建不删），本件只抄结论与逐字断言。

---

## §0 起手锚

| 项 | 现量读数（本腿自己跑的） |
|---|---|
| `git log -1 --date=iso-strict --pretty='%h %ad %s'` | `bd39f17f 2026-10-06T19:21:34+08:00 probes(票 111 追加编排者记 19:2x)：ci-census-read-1 死腿的 logs 我自己复跑后读回三问` |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal` | **0 行**（起手即干净） |
| `git status --porcelain` 全仓 | 751 行（全在 `.scratch/**` 与 `design/**` 等工作树残留，与本腿射程无关；本腿一枚未动） |

★**在飞腿声明**：起手那一刻 `cmd/wisp` 有别人的在飞腿 **`232-v1`** 正在跑突变。本腿对它**零读、零写、零 `go` 命令**——
一条 `go test`／`go build`／`go vet`／`go list` 都没发过；`cmd/wisp` 的读数全部是**读源码字节**得到的。
⛔ 本件里"会红／不会红"全部是**静态判读**，不是突变读数（能静态判动的都附了尺原文，判不动的写进 §4）。

禁读名单（本腿一次没打开）：`.scratch/wisp/probes/232/**`、`.../270/**`、`.../271/**`、
`.../pool-validity/5g2/**`、`.../expired-premises/7b2/**`、`frontend/**`、`design/**`。

---

## §1 铸造点与 model-supplied 那条路（行号＝本腿现取）

### 1.1 `agent.newTaskID()` 定义与熵源失败那一支的退化形

`internal/agent/loop.go:1128-1143` 逐字：

```go
1128  // newTaskID returns a UUIDv4-shaped task id (the task_log.id column is
1129  // documented as a UUID), built from crypto/rand so no new dependency appears.
1130  func newTaskID() string {
1131      var b [16]byte
1132      if _, err := rand.Read(b[:]); err != nil {
1133          // Entropy source unavailable: fall back to a clock-derived unique id.
1134          n := time.Now().UnixNano()
1135          for i := range b {
1136              b[i] = byte(n >> (8 * uint(i%8)))
1137          }
1138      }
1139      b[6] = (b[6] & 0x0f) | 0x40
1140      b[8] = (b[8] & 0x3f) | 0x80
1141      h := hex.EncodeToString(b[:])
1142      return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
1143  }
```

票面那句"熵源失败那一支退化成时钟派生"＝**成立，逐字对上**（`:1132-1137`，注释自己写的就是 `clock-derived`）。
本腿再加两条票面没写、但直接决定"这枚钉该钉什么"的现量：

1. **退化形长得和正常形一模一样**——`:1141-1142` 的格式化在 `if` 之外，输出仍是 36 字符 uuid 形
   ⇒ 全仓那枚唯一的形状尺（§2.1 A-1）**分不出这两支**。
2. 退化形既不是 128 bit，也不是"干净的一个时钟值"：`i%8` 把 64 bit 的 `UnixNano()` 循环铺满 16 字节，
   `:1139-1140` 再抹掉 12 bit ⇒ 真实熵＝**纳秒计数器在一个进程寿命内的取值个数**
   ⇒ "唯一"那一半站得住（单调计数器），"**不可猜**"那一半站不住。
   ⚠ 这一句是**读码推出来的**，本腿禁 `go`、**没有真跑**；要坐实得由带 go 权限的腿喂一次失败支（§4 第 2 条）。

### 1.2 全部调用点（锚在带括号的调用形状上，不裸 grep 符号名）

尺原文：`git grep -n "newTaskID(" -- internal cmd` → tracked 全量 6 行（`logs/01-newtaskid-hits.txt`、`logs/65-verify2.txt`）：

| file:line | 形状 | 性质 |
|---|---|---|
| `internal/agent/loop.go:329` | `id := newTaskID()`（在 `RunAsync`，函数头在 `:328`） | **铸造点 1** |
| `internal/agent/loop.go:340` | `return l.run(ctx, newTaskID(), input)`（在 `Run`） | **铸造点 2** |
| `internal/agent/loop.go:1130` | `func newTaskID() string {` | 定义 |
| `internal/agent/compress_trace_test.go:557` | `idA, idB := newTaskID(), newTaskID()` | **测试直调（同一行两次）** |
| `internal/tools/pointer_183_cli_seam_test.go:59` | `// agent.newTaskID() gives in production …` | **注释，不是调用** |
| `internal/tools/task_backfill.go:38` | `// model-supplied. agent.newTaskID() is …` | **注释，不是调用** |

⇒ 编排者给的三处参考读数（`:329`／`:340`／测试 `:557`）**逐枚对上**。现场教训一条：
裸 `grep newTaskID` 会把两行**注释**吸进来（本腿靠"带括号"剔掉的）；
`internal/agent/loop.go:333` `t.result = l.run(c, id, input)` 是**传**不是**铸** ⇒ **铸造点确实只有两枚，没有第三枚**。
（另核：`git grep "\.run(" -- internal cmd ':!*_test.go'` 里 `internal/tools/bridge.go:362` 是 `Bridge.run` 同名不同物，
`cmd/wisp` 那几枚 `.Run(` 是 ball/panel 的，都与环路无关。）

### 1.3 `internal/tools/task_backfill.go:35-46` 注释逐字（票面引用的那一段）

```
35  // The reason the artifact key cannot be the bare task id: the spill file name
36  // is derived from whatever key it is handed (agent.spill.go artifactName), and
37  // the loop's ordinary spill path hands it the tool-CALL id, which is
38  // model-supplied. agent.newTaskID() is a plain lowercase hex uuid, so a bare
39  // task id is a perfectly reachable call id - and two keys that encode to the
40  // same name land on the same file, where spill.go:117-131 documents
41  // last-writer-wins. A collision there is not an error, it is a silent byte
42  // swap: the model re-reads a pointer the host printed and gets somebody else's
43  // bytes, with C25 provenance broken and nothing raised. Prefixing the key (not
44  // changing the overwrite semantics - that is contract-adjacent) is what keeps
45  // the two namespaces apart.
46  const taskArtifactPrefix = "agent-task-"
```

票面 `:28` 那两处引用（`model-supplied`、"裸 task id 是完全可以被够到的 call id"）**逐字成立**。

### 1.4 谁把模型给的值当 task id 用（追到调用点）

| 通道 | 现量（file:line） | 结论 |
|---|---|---|
| task.spawn 的父 id | `internal/tools/subagent_197.go:259` `parentID := CorrelationID(ctx)`，`:260-263` 空即拒 | 只认宿主侧 |
| task.spawn 的参数面 | `subagentSpawnArgs` 只有 `description`／`prompt`；schema 逐字 `"additionalProperties":false` | ★**模型今天递不进一枚 id** |
| task.cancel 的调用者 id | `internal/tools/task.go:707` `caller := CorrelationID(ctx)`；`:730` `if rec.ParentTaskID != caller` | 只认宿主侧 |
| task.cancel 的目标 id | `internal/tools/task.go:614` `TaskID string `json:"task_id"`` | **模型自由递**（但要比对父） |
| task.output 的目标 id | `internal/tools/task.go:473`；`taskOutput.Execute`（`:503-560` 全文读过）**不比对 caller**，直接 `Roster.Look(a.TaskID)` | ★静态上"同进程内任何任务可点名读走任何别的任务的输出"**可达**。现有尺里唯一碰这带的是 `internal/tools/pointer_185_cli_seam_test.go:327` T5（逐字 `// T5 — AC#2 on the seam: another task cannot borrow this task's roster.`），它钉的是 **C25 作用域借道**，不是读取权限 ⇒ ⛔ 不许把 T5 读成这道门有牙 |
| 跨层兜底 | `internal/tools/bridge.go:269-271` 逐字 `if req.CorrelationID == "" { req.CorrelationID = req.TaskID }` | ★**只兜底到 TaskID，永不兜底到模型给的 CallID** |
| 三个 id 并排那一处 | `internal/agent/loop.go:654` `TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID` | 模型给的 `call.ID` **单占第三格**，从不进前两格 |
| `task.output` 的 schema 说法 | `internal/tools/task.go:481` 逐字 `"后台任务的 id（由宿主记录，不是路径）"` | 词面写了"由宿主记录" ⇒ **描述那一侧有句子，断言那一侧没尺**（按票 236 自己的判据：词面不构成牙） |

**今天真的喂进接缝的 model-supplied id 长啥样**
（尺：`grep -rhoE "\"id\"[: ]*\"[^\"]{1,48}\""` 全量去重计数，`logs/10-callid-shape.txt`）：
`resp_mock1`×20、`msg_out1`×12、`msg_mock1`×11、`fc_1`×3、`call_b2`×3、`call_e1`×3、`rs_1`×2、`r`×2、`m`×1、`chatcmpl-fake`×1 ……
⇒ **最短 1 字符，没有一枚是 36 字符 4 破折号的 uuid 形**（这一条直接决定 §2.1 的判语与 §3 的结论）。

### 1.5 D34 权威表的射程（现取行号；票面禁区＝不许新造工具面）

- D34 节头＝`docs/PLAN.md:2516`（逐字 `## 16.5 D34 — 内置工具权威表（**取代 D14 的六件套；docs/TOOLS.md 的唯一来源**）`），
  权威表本体＝`PLAN.md:2529`「### 16.5.2 权威表」起（表头在 `:2531`）。
- task 那一族在表里的行＝`PLAN.md:2564`（`task.list` / `task.cancel`，L0 / L1，S7）与 `PLAN.md:2565`（`task.output`，L0，S7）。
  ★**表里没有 `task.spawn` 行**：`grep -nE "task\.spawn" docs/PLAN.md` ＝ **0 命中**（本腿现量）。
  那一枚工具仓里已实现（`internal/tools/subagent_197.go`），名册来源是票 197 批次。本腿只登记事实、不判它算不算缺陷（超射程）。
- **表里没有任何一枚工具的动作是"铸造／替换／领取 id"**（表段内 `grep -nE "铸造|mints|由宿主"` ＝ 0 命中；
  逐字命中的只有 `task.go:481` 那句 schema 描述）。
  ⇒ 只到事实层的结论：**要钉"id 由谁铸造"这件事，不需要新增工具面**——它不是"模型能调的一个动作"，而是宿主内部一个值的出处。
  可落的**现有接缝**（三条，本腿都不裁）：`C5 LlmProvider`（golden SSE；§1.4 末那批 id 就从这儿进来）、
  装配根（`cmd/wisp/run.go:1099` `bg := loop.RunAsync(ctx, task)`、`internal/tools/subagent_197.go:344` `bg := child.RunAsync(childCtx, prompt)`）、
  以及 `internal/agent` 包内（铸造点在哪儿，尺就在哪儿）。
  参照形状：仓里**已有**一枚同族钉落在"接缝/包内"而不是工具面上——`internal/session` 的 minted-shape 锁（§2.4）。
  ⇒ **这一格能派产码腿**；⛔ 落点本腿不裁（§4 第 5 条）。

---

## §2 尺分类表

### 2.0 分母尺（三把，原文读数）

```
ls internal/tools/*_test.go cmd/wisp/*_test.go | wc -l                             → 116
grep -rilE "task ?id|taskID" internal/tools/*_test.go cmd/wisp/*_test.go | wc -l   → 36
grep -rinE "task ?id|taskID" internal/tools/*_test.go cmd/wisp/*_test.go | wc -l   → 203
```

- 票面"粗尺命中约 36 枚文件"＝**逐字对上**。
- ★**203 行里只有 18 行带断言动词**（`… | grep -iE "errorf|fatalf|assert|equal|want"` → 18，全量落 `logs/07-assert-lines.txt`）
  ⇒ 其余 **185 行是注释／wire struct tag／取值器／`t.Logf`**（§2.4 给了形状与样例）。
  **"含 taskID 这个词"与"它会红"之间差着一个数量级**，这就是本件主要防的事。
- 第三把（全仓形状词面尺，86 行，`logs/12-shape-vocab.txt`）：
  `grep -rnE "uuid|UUID|MustCompile" --include=*_test.go internal/ cmd/`
  ⇒ **真·task id 形状断言全仓只有 1 行**＝`internal/agent/compress_trace_test.go:526`；
  其余是 CSS token（`internal/ball/tokens_table_test.go`）、panel 路由正则、emoji 段、spill 文件名一族的正则。

### 2.1 类别 A｜形状检查（验 id 长啥样，不验来源）＝ **1 枚**（票面那两目录内＝0 枚）

**A-1 `internal/agent/compress_trace_test.go:524-528`**（属 `TestCompressionTraceCarriesTheOwningTaskID`，函数在 `:478`）逐字：

```go
524  // The shape check is not decoration: newTaskID mints a 36-char uuid-shaped
525  // id, and this is what tells the reading apart from a hand-typed constant.
526  if len(task) != 36 || strings.Count(task, "-") != 4 {
527      t.Errorf("trace task = %q, want the uuid-shaped id newTaskID mints (36 chars, 4 dashes)", task)
528  }
```

- **"换成模型给的值会不会红"＝会红**（机制写死）：这一枚的 `task` 是从**真 `h.run()` 的 `Result.TaskID`** 读回来的
  （`:481-488` 走真环路＋golden provider；`:513` 先 `if task != res.TaskID`）。
  一旦铸造点改成"用模型给的那枚 id"，喂进来的就是 §1.4 末那批形状（`msg_mock1`／`call_e1`／`r`）
  ⇒ `len(task) != 36` 成立 ⇒ `t.Errorf` ⇒ **该用例 FAIL**。
- ★**但它的红买不到"来源钉"**：喂一枚 36 字符 4 破折号的模型串（`00000000-0000-4000-8000-000000000000` 这种，
  C-16 已在仓里写死过）它**照绿**；且它比的是"长度＋破折号数"，**不是出处**。
- ⛔ 看不见：熵源失败那一支（§1.1 第 1 条，退化形同形状）；同一枚用例里 `:513` 那条自洽比对恒真，红只可能来自长度。

**A-2 相邻带（判：不是 task id 的形状尺，剔出类别 A）**：
`internal/agent/spill_name_injectivity_test.go:252` `if !regexp.MustCompile("-[0-9a-f]{16}\.txt$").MatchString(name)`、
`internal/tools/pointer_183_cli_seam_test.go:176` `!strings.HasPrefix(name, "tool-output-agent-task-")`
⇒ 钉的是**产物文件名那一族的前缀/摘要尾巴**（call-id 命名空间），换 task id 来源不红（前缀由 `taskArtifactPrefix` 给）。
`cmd/wisp/run_test.go:56` `taskIDRe = regexp.MustCompile(`任务 (\S+) 结束`)` 只**取回**不约束形状 ⇒ 不红（见 C-15）。

### 2.2 类别 B｜来源钉（断言 id 由 Go 侧铸造／断言它不等于模型给的值）＝ **0 枚**

四把负向尺原文（⛔ 零命中在这台仓要读成"没有"，不是"命令失败"）：

| 尺 | 命令 | 读数 |
|---|---|---|
| B-尺1 | `grep -rnE "TaskID *[!=]= *[a-zA-Z_.]*CallID" --include=*_test.go internal cmd`（正反两个方向） | **0 命中**（`logs/47-final.txt`）⇒ **没有任何一枚断言把两族 id 互相比对** |
| B-尺2 | `grep -rniE "entrop\|guessab\|unguess\|unpredict\|clock-derived" --include=*_test.go internal cmd` | **0 命中**（`logs/43-entropy-221.txt`）⇒ 熵支零射程 |
| B-尺3 | `grep -rn "UnixNano" --include=*_test.go internal cmd` | 4 命中，**全在 `internal/proc/` 的 mutex 命名**（`boot_windows_test.go:23/24`、`envfork_mutex_windows_test.go:28`、`singleinstance_windows_test.go:52`）⇒ **没有一枚测过退化形** |
| B-尺4 | `grep -rnE "rand\." --include=*_test.go internal/agent` | **0 命中** ⇒ `rand.Read` 用的是包级 crypto/rand、**没有注入接缝**，退化支在今天测试面上不可达（这是"尺看不见"，不是"支不存在"） |

**名字像来源钉、其实不是的三枚（逐枚读体，具名剔除）**：

- **B-x1 `cmd/wisp/subagent_carrier_197_test.go:430` `TestSubagentStreamKeyHasOneMintSite`**——名字带 mint。
  体（`:431-438`）用 `go/ast` 扫 `internal`＋`cmd` 里 stream key **字符串字面量**的非测试 mint 点，
  逐字 `t.Fatalf("the stream key literal %q is minted in %d non-test sources %v, want exactly one: %s", …)`（`:436`）。
  ⇒ 钉的是**前缀常量只有一处写**，从头到尾不碰 task id 的出处。**判：不是来源钉。**
  ★这一枚是本腿抓到的"名字像尺、其实不是"的现场样本（与派单里那条"别把含 taskID 当会红"同一族）。
- **B-x2 `internal/agent/compress_trace_test.go:557-560`**（属 `TestCompressionTraceNeverInventsATaskID`，`:534`）逐字：

  ```go
  557  idA, idB := newTaskID(), newTaskID()
  558  if idA == idB {
  559      t.Fatalf("setup: newTaskID returned the same id twice (%q)", idA)
  560  }
  ```

  ⇒ **全仓唯一一处测试直调铸造函数**，但它钉的是"同一进程内两次调用不撞"，且红形是 `t.Fatalf("setup: …")`＝**台件塌了**，不是"来源被换"。
  换成模型给的值：模型值不参与这两枚直调 ⇒ **不红**。
  ★结构上看不见跨进程／跨时间碰撞（同仓 session 那枚同类钉的注释 `internal/session/grants_test.go:106`
  逐字写了这个天花板：`which is why this function asserts 200 distinct mints rather than trusting …`；task 这边连这个都没有）。
- **B-x3 `internal/tools/task_cancel_221_legs_test.go:9`** 注释逐字
  `agent.Loop, real memory.Store journal) with a REAL host-minted task id as the`
  ＋体（`:153 return r.TaskID`、`:333`、`:380-381` 比的是 `parent197` 那枚宿主常量）
  ⇒ 它保证"载体用真铸的 id"，但**没钉"铸的 id 不许来自模型"**：两端都不与铸造点比出处 ⇒ **不红**。

### 2.3 类别 C｜字面量钉 ＝ **27 行 / 48 枚去重串 / 19 枚文件；判"会红"的 0 枚**

分母尺原文（`logs/51-literals-count.txt`）：

```
grep -rnE "(TaskID|taskID|CorrelationID|task) *= *\"[a-zA-Z0-9_:.-]{3,24}\"|\"task_id\":\"[a-zA-Z0-9_:.-]{3,24}\"" \
  --include=*_test.go internal/tools cmd/wisp | wc -l                        → 27 行
grep -rinE "task ?id|taskID" internal/tools/*_test.go cmd/wisp/*_test.go \
  | grep -oE "\"[A-Za-z0-9_:.-]{3,32}\"" | sort -u | wc -l                   → 48 枚不同串
```

逐组判语（★"钉的是 task id 还是别的字段"；"会红"口径＝把 `newTaskID()` 换成模型给的那枚串）：

| 组 | file:line（现取） | 逐字 | 钉的是谁 | 会红＋为什么 |
|---|---|---|---|---|
| C-1 | `internal/tools/failclosed_236_teeth_test.go:91`／`:98` | `lit236CancelNoHostTaskID = "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）"`／`lit236SpawnNoHostTaskID = "…不知道父任务是谁…"` | **拒因整句**（票 236 AC#1 的字面量钉） | **不红**：`:123`／`:261` 比的是 `Text` 整句，与 id 值无关 |
| C-2 | 同上 `:147`／`:287` | `Args: json.RawMessage(`{"task_id":"target-236r2"}`)`／`CallID: "call-236r2-spawn-no-parent"` | 目标/调用的 **fixture 串**（`:115 target := x.h.childRow(t).TaskID` 才是真铸的） | **不红** |
| C-3 | `internal/tools/ticket176r1_start_port_test.go:131` | `const taskID = "e3f1c0a97b2d4865"` | ★**唯一一枚"最像 uuid"的字面量**，在 L-2 命名空间钉里 | **不红**：三处比对（`:144 Contains("agent-task-"+taskID)`、`:155 sp.Path == pointer.ArtifactPath`、`:162 string(body) != mine`）**两端同源于这个 const** ⇒ 恒真 |
| C-4 | `internal/tools/pointer_183_cli_seam_test.go:63-64` | `r183BackgroundID = "c66c0634183r2a1b"`／`r183SiblingTaskID = "c66c0635183r2b2c"` | 背景任务 fixture id；注释 `:59` 自称"这就是 production 给的形状"＝**拿字面量冒充铸造产物** | **不红** |
| C-5 | `internal/tools/subagent_197_test.go:29` | `parent197 = "parent-task-197"` | 宿主盖给 harness 的父身份 | **不红**（`:308 if row.Out.ParentTaskID != parent197` 两端都是它） |
| C-6 | `internal/tools/subagent_197_test.go:327-329` | `key := SubagentStreamKey(row.TaskID)` / `if key != "subagent:"+row.TaskID` | **键形状**，右端复用铸来的值 | **不红**（换 id 两边一起换） |
| C-7 | `internal/tools/task_cancel_221_legs_test.go:333`／`:380-381` | `if rec.Out.ParentTaskID != parent197 \|\| rec.Out.Kind != TaskKindSubagent \|\| rec.Out.Label != label197`／`t.Errorf("tool_call 行的调用者是 %q, want %q", r.TaskID, parent197)` | **父子关系＋账上调用者** | **不红** |
| C-8 | `internal/tools/bridge_test.go:149`／`:671`／`:686` | `TaskID: "task-1", CorrelationID: "corr-1", CallID: "call-1"`／`rq.TaskID, rq.CorrelationID = "task-1", fmt.Sprintf("corr-%d", i)` | 手戳请求的 fixture 身份 | **不红** |
| C-9 | `internal/tools/bridge_scope_open_ticket158_test.go:92`／`:104`／`:122` | `const plainTask = "task-158-plain"`／`const taintedTask = "task-158-taint"`／`want := "task=" + taintedTask + " was_open=true"` | **C25 作用域开没开**，id 只当索引 | **不红** |
| C-10 | `internal/tools/task_output_leg_test.go:64/93/182/253/299/323/360` | `"bg-1"`／`"ghost-9"`／`"bg-7"…"bg-11"`；`:99 t.Fatalf("an unknown task id must be an error outcome, got: %+v", out)`；`:102 …!strings.Contains(out.Text, "ghost-9")` | **查不到要响亮报错＋错误分类** | **不红** |
| C-10b | `internal/tools/task_output_leg_test.go`（`roster.Record(taskID, …)` 处，注释约在 `:83`） | 注释逐字 `The task id is a parameter, not a hard-coded constant, because the caller's id and the target's id must differ in some cases` | ★唯一一处"铸造点会路过的"fixture 用法，但比的是**输出字节** | **不红** |
| C-11 | `internal/tools/task_output_ac2_before_test.go:43-45`＋`:65` | `TaskID: "164-ac2-before"`、`CallID: "call-164-ac2-before"`、`{"task_id":"bg-1"}` ＋ `t.Errorf("the refusal leaked the task id it never looked up: %q", out.Text)` | **未知工具那一支不许复述 id**（泄漏尺） | **不红**（比的是"文本里有没有那个串"） |
| C-12 | `internal/tools/ticket175r2_stamp_live_test.go:124`／`:148-149`／`:222`／`:280`／`:290-292` | 注释逐字 `The host mints the name; no leg below lets the model invent it.`；`const task = "175r2-j3"`＋`{"task_id":"bg-175r2-j3"}`；`func TestHostMintedPointerRereadStaysCleanOnRealBridge175r2` | **指针再读不重触 R4**（豁免落没落） | **不红**。★`:124` 那句"no leg lets the model invent it"是**注释不是断言** ⇒ 本腿按判调用形状同样的严谨判它：**词面不构成牙** |
| C-13 | `internal/tools/pointer_185_cli_seam_test.go:195/225/256/294/374` | `const task = "185r1-second"…"185r1-sibling"` | C25 作用域／roster 借道 | **不红** |
| C-14 | `cmd/wisp/subagent_carrier_197_test.go:203-204`／`:273`／`:283`／`:294` | `if row.Kind == tools.TaskKindRoot && row.TaskID != "" { return row.TaskID, nil }`；`if rootID == "" \|\| len(childIDs) != 1 \|\| childIDs[0] == ""`；`if sect.Rows[i-1].TaskID > sect.Rows[i].TaskID`；`if child.ParentTaskID != rootID` | ★**真·从生产铸造流读回来的 id**，但四处全是**自洽比对／非空／排序** | **不红**（宿主侧只把 id 铸成**空串**才会红 `:273`；换成模型给的**非空**串不红。`:283` 是排序尺，uuid 形与非 uuid 形都保序） |
| C-15 | `cmd/wisp/run_test.go:54-56`／`:138-144`／`:231` | `taskIDRe = regexp.MustCompile(`任务 (\S+) 结束`)`；`// taskID reads the id this run reported`；`tl, err := f.openStore().TaskLogByID(ctx, f.taskID())` | ★**最该成为来源钉的位置**：真装配根跑一遍、从状态行读回 id、再拿它去 `task_log` 找行 | **不红**（`(\S+)` 不约束形状，打印与落账同一个值即恒真） |
| C-16 | `cmd/wisp/subagent_stream_key_197_test.go:52-56`＋`internal/tools/subagent_197_test.go:886` | `"task-197-a"`, `"00000000-0000-4000-8000-000000000000"`, `"带中文的id"`／`"aB-1", "00000000-0000-4000-8000-000000000000", "带中文的id"` | **uuid 形字面量出现了**，但钉的是"两侧拼出的键相同" | **不红，且方向相反**：这两枚今天**主动在证明"uuid 形的 id（模型也能给）拼出的键一样合法"** ⇒ 它们是在为"模型给一枚 uuid 形 task id"**开绿灯**的尺 |
| C-17 | `internal/tools/tasklist_deferred_236r3_teeth_test.go:150`／`:202` | `func Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment`／`func Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt` | 名册在册 vs 注册表实存；overlay 仪器自证尺 | **不红**（本腿只核到"它不比 id 来源"这一层，逐行 taskID 断言未展开，`logs/05-all-lines.txt` 在案） |

**类别 C 的口径写死（防"枚数"被读成"牙数"）**：字面量行 **27**、去重串 **48**、涉及文件 **19**。
判"换铸造点会红"＝**0 枚**；判"钉的其实是别的字段"＝表内全部 17 组。
其中**唯二两处值来自生产铸造流**＝C-14（`subagent_carrier_197_test.go:203`）与 C-15（`run_test.go:138`），
**两枚都只做自洽比对 ⇒ 都不红**。

### 2.4 那 185 行"命中但不算尺"的形状（本腿为什么不列它们）＋一枚现成参照

| 形状 | 样例（file:line 逐字片段） | 为什么不是尺 |
|---|---|---|
| wire struct tag | `cmd/wisp/subagent_carrier_197_test.go:74` `TaskID            string `json:"taskId"`` | **类型声明**，与 AC#3 无关（正是派单警告的"类型声明混进命中"） |
| 文档注释 | `internal/tools/subagent_197_test.go:213` `// parent task id reaches the tool the way production delivers it.` | **词面**；按票 236 自己的判据（`tasklist_deferred_236r3` 那枚"不靠注释"的钉）注释不构成牙 |
| 取值器 | `cmd/wisp/subagent_carrier_197_test.go:127` `if r.TaskID == taskID {` | **查表索引**（配 `:131 t.Fatalf("the packet carries no roster row for %q: %+v", taskID, sect.Rows)`＝"这一行在不在"，不比来源） |
| 日志/诊断 | `cmd/wisp/subagent_carrier_197_test.go:381`／`:584`、`internal/tools/task_output_leg_test.go` 的 `t.Logf` 行 | 只打印 |

★**仓里已有一枚同族"来源钉"的成熟形状**（不是本腿造的新概念，给后程落点用的现成参照，射程在 `internal/session`，超出票面两目录）：
- `internal/session/grants_test.go:152` `TestTicket224TestLiteralsAreOutsideTheMintedShape`，逐字断言 `:161`
  `t.Errorf("the hand-written literal %q validates as a minted session: the shape "+"lock is no longer excluding fixtures", lit)`；
- `internal/session/grants_test.go:498` `TestTicket224NewLedgerRefusesHandWrittenIDs`；
- `internal/session/session.go:67 func (id ID) Valid() bool` ＋ `:94 func Mint() (ID, error)`；
- `cmd/wisp/ticket224_assembly_test.go:433` 逐字
  `t.Fatalf("boot 1 minted %q, which the shape lock rejects: the minting point and the …")`。
⇒ **task id 这一族今天没有这套东西**（B-尺1～4 零命中）。这一条是"钉该长什么样"的仓内先例，不是本腿的设计建议。

### 2.5 三类合计

| 类别 | 票面射程（`internal/tools`＋`cmd/wisp` 的 `_test.go`） | 全仓 |
|---|---|---|
| A 形状检查 | **0 枚** | **1 枚**（`internal/agent/compress_trace_test.go:526`） |
| B 来源钉 | **0 枚** | **0 枚** |
| C 字面量钉 | **27 行 / 48 串 / 19 枚文件**（判会红 0 枚） | 同（其余包的 id 字面量是 session/grant/corr 族，不属 task id） |
| **合计会在"换成模型给的值"时变红** | **0 枚** | **1 枚**（A-1；且红因＝长度，不是来源） |

---

## §3 对抗"我现读＝0 枚"这一句

**判决：不成立——但要分两个口径下刀，整体收下和整体推翻都是错的。**

★**编排者那句『0 枚』不成立，真值是 1 枚，尺是 `internal/agent/compress_trace_test.go:526`**
（用例名 `TestCompressionTraceCarriesTheOwningTaskID`，函数在 `:478`）。断言逐字：

```go
526  if len(task) != 36 || strings.Count(task, "-") != 4 {
527      t.Errorf("trace task = %q, want the uuid-shaped id newTaskID mints (36 chars, 4 dashes)", task)
528  }
```

**为什么会红（机制，三句，每句都指回本件现量）**：

1. `task` 这个值不是手写的，是从**真环路跑一遍**读回来的（`:481-488` `h.run("收个尾")` → `:507` 取 `rec.attr["task"]` → `:513` 先与 `res.TaskID` 自洽比对）；
2. 铸造点一旦改成"用模型给的那枚 id"，这里的值就变成 §1.4 末现量到的那批形状
   （`msg_mock1` 9 字符／`call_e1` 7 字符／最短一枚是 `"id":"r"` 的 1 字符），`len != 36` 当场成立；
3. ⇒ `t.Errorf` ⇒ 该用例 FAIL ⇒ **会红**。

**这一枚今天被漏掉的原因（本腿判这就是那句负向句的成因，形状与"读盘尺对 `-overlay` 结构性瞎"同族）**：
它的射程在 **`internal/agent`**，而 AC#3 让筛的是 `internal/tools/*_test.go` 与 `cmd/wisp/*_test.go` 两个目录
（`logs/04-census-raw.txt` 的分母尺就是从这两目录开的）。**按包开尺就看不见跨包的那一枚**——
这正是记忆里那条"按包门禁看不见全仓不变式"的现场版。

**但"0 枚"在另一个口径下成立，这一半必须一并写清，不许被上一半抹平**：
**在票面点名的两目录（`internal/tools/*_test.go` ＋ `cmd/wisp/*_test.go`）内＝0 枚**
（§2.5 那张表的第一列；类别 B 四把负向尺零命中、类别 C 27 行字面量全数判不红）。

**★比枚数更要紧的那一句：这 1 枚红，买不到 AC#3 想要的东西。**

- 它的红因是**长度巧合**，不是出处：模型今天恰好没有给 uuid 形（`logs/10-callid-shape.txt` 现量），所以它红；
  喂一枚 36 字符 4 破折号的模型串它**照绿**——而这种串**仓里已经写死在 fixture 里**了
  （`cmd/wisp/subagent_stream_key_197_test.go:55` 与 `internal/tools/subagent_197_test.go:886` 的
  `"00000000-0000-4000-8000-000000000000"`，两枚都在**主动证明"uuid 形的 id 拼出的键一样合法"**）。
- 熵源失败那一支的退化形**同形状**（§1.1 第 1 条），这一枚也分不出。
⇒ **结论（给编排者翻勾用的那一句）**：**"task id 由谁铸造"这一格今天没有来源钉（B＝0 枚，四把负向尺零命中）；
唯一的红来自一枚形状尺，而它对"模型给一枚 uuid 形串"和"退化支"两头都是瞎的。**

---

## §4 判不动的地方（这一节是保命节，写满）

### 4.1 本这把静态普查尺结构性看不见的五类形状

1. **看不见"哪个包在 CI 真跑"**。尺：`grep -rn "internal/tools" scripts/portable-tests.sh` 命中 core tier（`:240`），
   而 windows tier（`:249-253`）**不含 `internal/tools`**、cli tier 只有 `./cmd/wisp/`。
   ⇒ 静态尺能列出"哪些用例存在"，列不出"哪些用例今天真的被执行"。这一格必须由 `ci.yml`＋archived CI logs 那条腿答，本件不替它答。
2. **看不见 `-overlay` 突变后的字节**。票 236 自己在 `internal/tools/tasklist_deferred_236r3_teeth_test.go:202`
   把这件事钉成了仪器事实（`Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt`）。
   ⇒ 本件的"不会红"全部是**谓词层判读**，不是**突变读数**；真值要以 `go test -overlay` 那把腿为准。
3. **看不见注释与断言之间的落差被写成"看起来有牙"**。本腿现场抓到三例：
   `internal/tools/ticket175r2_stamp_live_test.go:124`（`The host mints the name; no leg below lets the model invent it.`——是注释，不是断言）、
   `internal/tools/task_cancel_221_legs_test.go:9`（`with a REAL host-minted task id`——是载体自述，不是钉）、
   `internal/tools/pointer_183_cli_seam_test.go:59`（把 `newTaskID()` 的产物形状冒充进 fixture）。
   ⇒ **静态尺把这些行计进"36 枚文件"，但它们一行牙都没有**；反方向也成立：删掉一句注释不会让任何用例红。
4. **看不见跨进程／跨时间的 id 碰撞**。唯一一枚不撞尺
   （`internal/agent/compress_trace_test.go:557-560`）只在**同一测试进程内**调两次
   （同仓 session 那枚同类钉的注释逐字承认了这个天花板：`internal/session/grants_test.go:106`）。
5. **看不见"名字像尺"与"其实是尺"的差**。B-x1 `TestSubagentStreamKeyHasOneMintSite` 带 mint 字、体里没有 task id 出处；
   反过来 `cmd/wisp/run_test.go:138` 那个 `taskID()` 取值器**才是唯一真从生产铸造流读回来的位置**，
   而它今天只做自洽比对——**静态尺若按名字筛，会把这两枚读反**。

### 4.2 要真跑才知道的三件事（本腿禁 `go`，一条都没跑）

1. **A-1 到底会不会红**（§3 的机制是静态推的）。定向突变最小形：
   只把 `internal/agent/loop.go:329`＋`:340` 两枚 `newTaskID()` 换成 golden provider 给的 id（不动其它任何行），
   指名用例必须是 `TestCompressionTraceCarriesTheOwningTaskID` 一枚红；本腿**没有这个读数**，不许被引用成有。
2. **熵源失败那一支能不能被够到**。§1.1 第 2 条（"退化形真实熵＝纳秒计数取值个数"）是**读码推的**；
   而 B-尺4 现量 `rand.` 在 `internal/agent` 测试面 **0 命中** ⇒ `rand.Read` 是包级直调，**今天没有注入接缝**，
   这一支在测试里**不可达**（"尺看不见"≠"支不存在"）。要量它必须先开接缝，属产码射程、不属本腿。
3. **task.output 那条"不比对 caller"的读取路径今天到底能不能读到别人的输出**（§1.4 那一行）。
   静态读到的只是"没有那道检查"；真跑一次跨任务读取才知道是**读到了**还是**被别处拦了**。
   ⛔ 本腿**没有**把它写成漏洞结论，只写成"静态可达、无尺"。

### 4.3 一枚尺"今天有没有牙"要分两档答（票面点名的那一枚）

`internal/tools/task_output_pointer_notice_test.go:420` `TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`：

| 档 | 读数 | 尺原文 |
|---|---|---|
| **本机（Windows）跑不跑得到** | **跑得到**：`:404` `exec.Command("cmd", "/c", "mklink", "/J", …)` 成功即建真 junction，不 skip ⇒ 有牙 | `sed -n '395,432p'` 逐字（`logs/31-ac5-skips.txt`） |
| **CI 跑不跑得到** | **ubuntu 档：进得了作用域、执行不了**——`internal/tools/...` 在 core tier（`scripts/portable-tests.sh:240`），但 ubuntu 造不出 NTFS junction ⇒ `:408 t.Skipf`。而且它**不在 ledger 名册里**（HEAD 现量：`git show HEAD:scripts/portable-tests.sh \| grep -c "NotRelayTheFiledPath"` ＝ **0**），`runtests.sh:98` 与 `portable-tests.sh:704` 两档都把未记账 SKIP 判红 ⇒ archived 日志 `run-37166458550-failed.log` 里它 **`--- SKIP`×2 ＋ `portable-tests.sh: unaccounted SKIP lines`×1**（本腿 grep 计数）。**windows 档：不在作用域**（win tier 不含 `internal/tools`）⇒ 今天这一枚在 CI **跑不到执行路径**，只跑得到"它 skip 了"这一声 | `logs/32-skip-policy.txt`、`logs/57-skipnames.txt`、`logs/62-core-summary.txt`、`logs/72-leg-date.txt` |

⇒ 判语：**这一枚今天的牙只在**本机会咬，CI 侧它既咬不到、还**倒过来把 core step 顶成红**（未记账 skip）。
⛔ 本腿**不裁**它该进 ledger 还是该改平台 tag（那是票 111 AC#10 与票 236 AC#4 那一格的地界，且票面 `:28` 写明了本票不许顺手扩边界）。

### 4.4 为什么"钉该落在哪一层"本腿不裁

1. **本腿只有静态尺**：落点的选择标准是"哪一层能因定向突变而红"，这**必须是真读数**才能定，静态判不了（§4.2 第 1 条）。
2. **三个候选层各自的天花板不同**，静态只能到"它在不在作用域"这一层：
   `internal/agent` 包内（铸造点旁边，A-1 就在这儿，天花板＝只看形状）／
   `C5 LlmProvider`（golden SSE，能造"模型给一枚 uuid 形串"那一形，天花板＝测不到装配根怎么用它）／
   装配根（`cmd/wisp/run.go:1099`，唯一真把两族 id 串起来的地方，天花板＝windows 作用域问题，见 4.3）。
   ⇒ 这三条**互不覆盖**，选哪条是**设计裁决**，不是普查读数。
3. **票面 `:28` 自己把顺序定死了**："先量清…再定钉的落点"。本腿交的是**前半**。
4. ★**给编排者的一条硬事实，不是建议**：仓里**已有一枚成熟同族形状可抄**（§2.4 session minted-shape 锁：
   `Valid()` 谓词 ＋ "字面量必须在 minted 形状之外" ＋ "真装配铸的 id 永不等于任何字面量" 三段双向），
   而且它**不靠新工具面**（§1.5 D34 现量：表里没有一枚工具的动作是铸造 id）⇒ **这一格不触发票面禁区，能派产码腿**。
5. ⛔ **未定义即停**：本腿不写落点、不写判据措辞、不勾任何框。AC#3 那枚框要由**另一个 agent** 凭真读数翻（AGENTS §0 第 3 句）。


