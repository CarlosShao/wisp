# 151 验收 r1（非实现者对抗验收程）— `Bridge.CloseTask` 只开不关这一案的第二枚验收表

> 本件是**验收表**，不是实现件。派单＝`.scratch/wisp/dispatches/2026-09-26-103x-accept-151-r2.md`。
> 被验版本＝**`4e16976`**（票 151 自己链的最后一枚）。实现件＝`docs/evidence/s1/151-task-scope-never-closed-r1.md`（只读，未改一字）。
> 票面 151 与票面 154 **本程一字未动**（派单与编排者都点名了这条，见末节的现量）。

## 第 0 节 锚点与环境（本程现读，不抄派单）

```
$ git rev-parse --short HEAD          ->  b0f68b2      （本程开工那一刻的工作树头）
$ git branch --show-current           ->  dev
$ git cat-file -t 4e16976             ->  commit
$ git cat-file -t 5d46f24             ->  commit
$ git cat-file -t 45c920e             ->  commit
$ git log --oneline -1 4e16976        ->  4e16976 evidence(151 r1): 判定＋修法交件……
```

- **取版**：`git archive 4e16976 | tar -x -C D:/tmp/wisp151accept`（仓库外）。**未在仓库内建 worktree、未 checkout、未 stash/reset**（AGENTS §1.4）。
- 变异与复算用三枚仓库外树：
  | 树 | 内容 | 用途 |
  |---|---|---|
  | `D:/tmp/wisp151accept` | `4e16976` 纯净树 | 摘一味／摘另一味的变异尺 |
  | `D:/tmp/wisp151accept-base` | 同一枚树 ＋ `third_party/sherpa-onnx`（从工作树复制的三枚 dll） | 改后门禁全跑、`go vet`、`gofumpt`、`d22scan` |
  | `D:/tmp/wisp151accept-pre` | `5d46f24^` 的 `run.go`＋`bridge.go`（用 `probes/151/run.head.go`/`bridge.head.go`，本程 `sha256sum` 现核＝与 `git show 5d46f24^:<file>` **逐字节相同**），新用例文件摘成 `package main` | 改前门禁全跑（不走 `-overlay`，见第 5 节） |
- **PATH 用 shell 自己的路径形**（`/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx`），未用 `pwd -W` 盘符形。
  ⚠ 本程第一发全包门禁**就是这么死的**：纯净树里没有 `third_party/`（`git archive` 不带未跟踪件），
  8 枚"要起真子进程"的用例报 `no native DLLs in ..\..\third_party\sherpa-onnx` 并转红
  （原始读数＝`probes/151-accept/gate-base-cli-v.txt`：RUN=132 顶格 70 红 8）。
  **按派单第 55 行的次序处理＝那是"仪器没跑到"，不是红**；补上 dll 目录后同一枚树重跑＝RUN 138／顶格 78／红 0（第 5 节）。
- 本程写面只用过两枚路径：`docs/evidence/s1/151-task-scope-never-closed-r1-accept-r1.md`（本件）＋ `.scratch/wisp/probes/151-accept/**`。
  其余全在仓库外。临时件**只建未删**。

## 第 1 节 票面 5 枚 AC ↔ 本程 5 格 双向对账（派单攻击点②）

现量（锚点 `4e16976` 的那枚票面，两把尺都给）：

```
$ git show 4e16976:.scratch/wisp/issues/151-*.md | grep -c '^- \[ \]'   ->  5
$ git show 4e16976:.scratch/wisp/issues/151-*.md | grep -c '^- \['      ->  6
```

⇒ **票面 AC 框＝5 枚**（`:27 :29 :34 :36 :39` 五行，逐行现读＝AC#1..AC#5）。实现件 §10 写"票面 **6** 枚 AC 框"。
那枚"6"最省事的复现尺就是上面第二行：`^- \[` 会把 `:56` 那行 Progress log 的 `- [2026-09-26 00:5x +08]` 一起数进来。
⇒ 判：**数错的是实现件，且是一枚"尺"错、不是"格"错**——它自陈未勾的框实际是 5 枚，5 枚它都交了读数，
所以**没有哪一格因此没人裁**（下表第五列＝本程那一格的号，五格全在本件里）。

| 票面 AC（原文标题，现抽） | 实现件那一节 | 本程那一格 | 本程裁决 |
|---|---|---|---|
| AC#1 先量"不关"今天到底改变了什么 | §2 | 第 2 节 | **成立，但一句结论要收窄**（可达性那一问，见第 2 节末） |
| AC#2 会响的检（硬核心） | §3 | 第 3 节 | **成立，两味都有牙**；但 §3 那句"两枚在未修码上结构上跑不起来"**本程现量证伪其一半** |
| AC#3 修法只在"谁拥有关闭"这一层（三选一） | §4 | 第 4 节 | **成立**，且"不许顺手改 `OpenTask`"那半句逐字节核过 |
| AC#4 契约轴（零字节） | §5 | 第 5 节前半 | **成立**（三枚 commit 的名册并集＝5 个路径，一枚不在禁写面里） |
| AC#5 门禁（逐包单跑） | §6 | 第 5 节 | **六数全中、名册两向对得上**（本程独立跑，未走 `-overlay`） |
| （票面另四条"现量的形状"＝AC 之前的第 1 节，不是 AC 框） | §1 | 第 2 节开头 | 四条本程都重走过，一条尺写清了 |

派单那七条攻击点落在哪：① →第 2 节；② →本节；③④⑤ →第 3 节；⑥ →第 5 节；⑦ →第 6 节。总判＝第 7 节。

## 第 2 节 格 AC#1 —「不关」今天到底改变了什么（派单攻击点①落这一格）

### 2.1 票面那四条"现量的形状"，本程在 `4e16976` 上重跑（尺写清楚）

| 那一条 | 本程用的尺（大小写敏感／整词／范围） | 现读 |
|---|---|---|
| `CloseTask` 调用者 | `grep -rn "CloseTask" --include='*.go' cmd internal`（大小写敏感、子串，范围＝**纯净树 `4e16976`**） | 6 行命中，其中**调用点只有 1 枚**＝`cmd/wisp/run.go:563`；其余是 `run.go:546` 注释、`task_scope_close_151_test.go:7,16` 注释、`bridge.go:642` 注释、`bridge.go:655` 声明 ⇒ **修法接上线之后**确实长在组合根上 |
| 未修码上的同一问 | `git grep -n 'CloseTask' 5d46f24^ -- '*.go'` | 只有注释＋声明两枚，**调用者 0 枚**，测试里也 0 枚 ⇒ 票面那条"连测试都没有"在开票那一刻成立 |
| `OpenTask` 唯一非测试调用者 | `grep -rn "OpenTask" --include='*.go' cmd internal \| grep -v _test` | `bridge.go:559`（`mark` 里，实参 `dec.TaskID`）＋声明＋注释 ⇒ 与实现件 §1 第 2 条同一枚位置（**行号按 `4e16976` 取，不是按它们的 `86b0161`**） |
| 一个进程几枚环路任务 | `grep -rn "agent.New(" / "loop.Run(" --include='*.go' cmd internal \| grep -v _test` | 各 1 枚（`cmd/wisp/run.go`），与实现件 §1 补的那条一致 |

⇒ 这四条**全部对得上**，实现件 §1 没有过量、也没有少量。

### 2.2 它的 AC#1 读数：本程**独立重跑**了那枚探针

尺：把 `probes/151/ac1-probe-source.go` 原样复制进纯净树的 `internal/tools/`（仓库外那枚树，包名就是 `tools`，
与它们用 `-overlay` 注入是同一枚编译单元），`go test -count=1 -v -run 'TestProbe151' ./internal/tools/`。
原始读数＝`probes/151-accept/ac1-probe-rerun.txt`，rc=0（PASS 的是"量到了"）。逐枚对数：

| 腿 | 实现件 §2 记的 | 本程现量 | 对得上？ |
|---|---|---|---|
| `no-close_today` 第二发裁决 | `L2 rules=[R4] approvals=1`，reason＝`R4: 包含来自 unbound-scope …` | **逐字相同**（`scopes_left=1 marks_left=1`） | 是 |
| `close-first_control` | `L0 rules=[] approvals=0 reason=""`，关闭即时 `scopes 1->0 marks 1->0` | **逐字相同** | 是 |
| 内容级三腿（同 id 复用／换新 id／换 id 且已关） | 命中 `fs.read <路径>`／命中 `unbound-scope`／`L0` 零卡 | **逐字相同**（本程的命中路径是这台机器的临时目录，形状一致） | 是 |
| 无审批通道那一形（rounds=8） | `executed=1 refused=7 L2cards=7` | **`executed=1 refused=1`（rounds=2）／`executed=1 refused=7 L2cards=7`（rounds=8）** | 是 |
| 一路点允许那一形（rounds=64） | `executed=64 refused=0 L2cards=63`、scopes 64、marks 64 | **相同**，heap `594,656 -> 3,015,152`＝`delta=2,420,496`＝**37,820 B/轮** | 是（斜率与它的 37,652 差 0.45％，同一次跑内 host 噪声） |
| 再多一枚什么都不读的干净任务 | `64 -> 65` | **`scopes 64 -> 65`、`L2cards 63 -> 64`** | 是 |
| 逐一 `CloseTask` 之后 | 表归零、堆回到基线附近 | **`scopes=0 marks=0 heap=587,008`（基线 594,656）** | 是 |
| 小 N 的堆读数不可用（它自己认的坑） | rounds=2 两腿 delta 为负 | **本程也是负的**（-183,960 / -100,056） | 是（坑是真的，不是它推的） |

⇒ **AC#1 的量法与量到的东西成立**："渗"＝会（且是变严不是漏内容，这个方向它也钉了）；"长"＝会、无上限、可归零；
内容级渗＝不会。票面 AC#2 给的作废前提（"既不渗也不长"）**不成立**这一判语，本程跟着成立。

### 2.3 派单攻击点①那一问：这条链到底是"生产可达"还是"真桥可达、CLI 那一轮不可达"

**本程自己造了那一发，并把它造到哪儿为止量了下来。** 台件＝`probes/151-accept/accept151_e2e_test.go`
（原始输出 `e2e-shot-shipped-code-v3.txt`、`…-v5.txt`）与对照件 `accept151_control_test.go`（`e2e-control-l1-write.txt`）。

造法（不动仓内任何一枚文件，全部在仓库外那枚纯净树里）：mockllm 除"确定性合成"之外还有**金样本回放腿**
（`tools/mockllm/main.go:19-21`：model 写作 `golden/<name>` 就逐段回放 `internal/llm/testdata/golden/<name>.sse`，
`server.go:145 nextGolden` 一次请求吃一段）。⇒ 实现件 §8 第 5 条那句"mockllm 只在 `tool_choice` 非空时才吐
`tool_calls`"只描述了**合成腿**（`tools/mockllm/chat.go:172` 那一行确实是它的尺），**金样本腿不看 tool_choice**。
于是本程现写一篇两段的金样本（第 1 段＝模型发起一发 `fs.read`，参数是本轮 allowlist 内那份身份证号文件；
第 2 段＝纯文本答复），config.toml 把 `text_chain` 指到 `acme/golden/…`，真跑 `runTextTask`：

```
现读（本程，4e16976 树）：exit=0
  stdout: [工具 fs.read -> error] … 任务 08281a77-… 结束（completed，2 轮，1 次工具调用 …）
  stderr: [audit] tools: C25 scope closed task=08281a77-… was_open=false dropped=0 open_scopes=0
```

⇒ **模型确实发起了那一发（2 轮 1 次工具调用），但它到不了桥**：`was_open=false dropped=0`
说明这一轮的 scope 根本没被开过（`mark` 没跑到）。**为什么到不了，本程做了对照**（同一枚腿、同一套金样本机制，
只把工具换成声明档位 L1 的 `fs.write`）：

```
对照现读：exit=0
  stdout: [工具 fs.write -> success]        且 out.txt 真落盘（19 字节）
  stderr: [audit] tools: C25 scope closed task=b27e0404-… was_open=false dropped=0 open_scopes=0
```

⇒ 档位之外别无变量，所以卡住 `fs.read` 的就是**声明档位**，三处原文（全在 `4e16976` 上现读）：
`internal/tools/fs.go:298` `FSReadDecl → Declared: risk.L0`；
`internal/tools/bridge.go:224` `ToolInfo.RiskLevel = levelString(e.Decl.Declared)`；
`internal/agent/loop.go:776-789` `decideRisk` 只让 L1/L2 过，default 分支在
`Config.PassThroughUnclassifiedRisk == false` 时**直接拒绝**，而 `cmd/wisp/run.go:588-595` 的
`agent.Config{…}` 字面量里**没有这一枚键**（本程尺：`grep -rn "PassThroughUnclassifiedRisk" --include='*.go' internal cmd | grep -v _test`
⇒ 只有 `loop.go:125/128/767/783`、`agent/tools.go:23`、`tools/gate.go:133` 五处，**组合根一处都没有**）。

**裁决（两问分开答，档位不同）**：
1. **"这一类缺陷在生产上可达"＝成立**，但要按 `bridge.go:559` 那一发的字面读：**任何在桥上派发的宿主**（面板／内部调用方／已装配好的 `rt.bridge`，
   `cmd/wisp/run_test.go` 的 `TestHostDispatchThroughTheAssembledBridge` 就是这一形的现成形状）一发敏感读就挂一枚 scope，
   且今天没人关。实现件 §2.1 那张"两形代价"表里的每一枚数本程都独立复到了（§2.2），所以"渗/长"不是推的。
2. **"今天的 `wisp run` 控制台腿会天天吃到这一发"＝不成立**，且原因**比它自陈的那一枚更深一层**：
   它 §8 第 5 条只归给 mockllm 的 `tool_choice` 门；本程现量到**第二道门是码给的**——
   环路在 L0/未分级这一档就把它拒了，所以**换真 provider 也到不了桥**。
   ⇒ 实现件 §2.1 那句"（今天的 `wisp run` 控制台腿、`NoGate`）"是把**桥的配置形**（NoGate）写成了**这条腿的实际后果**；
   按本程的读数，这条腿今天连"第一枚任务读完敏感内容"都发生不了。
   **这是登记不准，不是假绿**：它反而把自己的后果面**说大了**一格，方向与放水相反，但下一程若照它这句去验收"控制台腿已修"，会被它骗。
   ⇒ 本程把这一条作为**对本票结论的更正值**记进总判的附条件里（第 7 节），并点名它 §8#5 那句"本程没测什么"应当补的正是这第二道门。

⚠ 本程**没有**造出"端到端一轮里关掉一枚带污点的 scope"那一发（`dropped>=1` 那一形）。上表两枚读数就是本程造到为止的位置：
金样本能让环路发起 `fs.read`，但环路自己的风险政策不给它到桥。要造出那一发得动 `internal/agent/**` 或
`cmd/wisp/run.go` 的 Config 字面量——**那是别人的地界（票 153／组合根），本程没那个写面，也没自作**。
所以"实现件的 T1 钉的是 `dropped=0` 那一形"这句话**本程跟着成立**（它 §8#5 自陈的那一半没被本程推翻）。

## 第 3 节 格 AC#2 — 会响的检（派单攻击点③④⑤落这一格）

### 3.1 那三枚 commit 与它自陈的形状（攻击点④）

```
$ git cat-file -t 45c920e / 5d46f24 / 4e16976   ->  commit / commit / commit
$ git show --name-only --format='%h %s' 45c920e  ->  23 枚，全在 .scratch/wisp/probes/151/  下，别枚＝0
$ git show --name-only --format='%h %s' 5d46f24  ->  恰 3 枚：cmd/wisp/run.go、cmd/wisp/task_scope_close_151_test.go、internal/tools/bridge.go
$ git show --name-only --format='%h %s' 4e16976  ->  恰 1 枚：docs/evidence/s1/151-task-scope-never-closed-r1.md
$ git show -s --format='%s' 5d46f24              ->  placeholder
```

⇒ 攻击点④那一问的前半**对得上**：名册恰那三枚、无夹带；message 确实是 `placeholder`，实现件 §10 自打一枪登记过，
本程跟着核过"内容正确、消息坏"这一句为真。**这枚不是验收表能替他补的**（补记＝改历史，AGENTS §1.4 禁 `--amend`/`reset`），
只能由编排者以追加的方式处理 ⇒ 记进第 7 节的"要编排者动手"那一栏，不算它的实现缺陷。

### 3.2 摘一味（攻击点③）：两味**都有牙**，本程各造一枚

尺：变异全在仓库外的 `4e16976` 纯净树上做，逐枚先 `diff` 证明"其余一字不动"，再 `go test -count=1 -v -run '<两枚新用例>' ./cmd/wisp/`。

**味 1＝`run.go` 里那一行 `rt.bridge.CloseTask(taskID)`**（`probes/151-accept/mutA-drop-closetask-call.txt`）

```
$ diff run.shipped.go cmd/wisp/run.go
563d562
<               rt.bridge.CloseTask(taskID)        ← 只少这一行，别的一字未动
$ go test … -run 'TestCompositionRoot…|TestAdmitTask…' ./cmd/wisp/     rc=1
    task_scope_close_151_test.go:62: 任务结束时没有关闭它自己的 C25 污点 scope：stderr 里找不到以
      "[audit] tools: C25 scope closed task=51a044eb-… " 开头的审计行
--- FAIL: TestCompositionRootClosesTheLoopTasksTaintScope (1.67s)
    task_scope_close_151_test.go:109: … 第二发 judged L2, want L0；text="实时语音（Path C）…未经文本循环登记（AdmitTextTask），已 fail-closed 拒绝"
    task_scope_close_151_test.go:113: 第二轮没有读到文件内容，text="…"
--- FAIL: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.97s)
```

⇒ **两枚全红，红因不同**（一枚丢审计行、一枚丢裁决），行号 `:62/:109/:113` 与实现件 §3 引的那三句**逐字对得上**（只有 uuid 是本程现读的）。
"哪条用例变得不响？"＝没有一条不响。⇒ 这一味**不是装饰**。

**味 2＝`CloseTask` 里那行 `b.log(...)` 审计**（本程量到一件实现件没说的事：**这一味摘不干净**）

```
B1 只删那两行 b.log(...)：
$ go test … ./cmd/wisp/
# github.com/CarlosShao/wisp/internal/tools
internal\tools\bridge.go:659:2: declared and not used: dropped
internal\tools\bridge.go:663:2: declared and not used: left
FAIL    github.com/CarlosShao/wisp/cmd/wisp [build failed]        （probes/151-accept/mutB1-drop-audit-log-compile.txt）
```

⇒ 那枚审计行是 `dropped`/`left` 两枚局部量的**唯一消费者**，按"其余一字不动"的字面去摘会得到**编译失败**，
不是一枚可判的读数（编译失败既不能算红也不能算绿＝仪器没跑到）。所以本程补一发**行为中立的 B2**
（删掉 `b.log` ＋ 补 `_, _ = dropped, left` 让编译过去，除这两处外一字不动）：

```
B2：probes/151-accept/mutB2-drop-audit-log-behavior-neutral.txt
--- FAIL: TestCompositionRootClosesTheLoopTasksTaintScope (2.01s)     ← 只这一枚红
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.73s)      ← 另一枚仍绿
```

⇒ **审计行有牙，且牙口正好只长在 T1 上**——这正是实现件 §4 末段那句"摘掉它 `TestCompositionRootCloses...` 直接红"
与 §3 那张"一枚丢审计行、一枚丢裁决"的分工。**两味分开钉两半，本程跟着成立。**
（但"摘一味"这一动作在味 2 上必须付一发额外改动才成立 ⇒ 实现件 §4 那句话该补一句"这枚审计行与两枚局部量绑死，
不是可独立删除的一行"，本程把它列进"要更正措辞"那一栏，不改判定。）

### 3.3 恒真这一问（攻击点⑤）：它那句推理**只对两枚里的一枚成立**

实现件 §3 原文："这两枚在'没有本票改动'的码上**结构上跑不起来**（`rt.admitTask` 这枚方法就是本票引入的，
未修码上连编译都过不去）。所以本件不用'未修码'当红绿尺，用的是**变异尺**。"

- `rt.admitTask` 确是本票新引入（本程尺：`git show 5d46f24 -- cmd/wisp/run.go` 里那一枚 `+func (rt *agentRuntime) admitTask`；
  `grep -c admitTask` 在 `5d46f24^` 的三枚相关文件上＝**0**）⇒ **T2 那一半成立**（它函数体里直接调 `rt.admitTask`）。
- **T1 那一半不成立**：`TestCompositionRootClosesTheLoopTasksTaintScope` 函数体只用到 `newRunFixture` / `f.run` / `f.taskID()` /
  `strings.Contains`，这四件在 `5d46f24^` 上**全都有**（本程尺：`git show 5d46f24^:cmd/wisp/run_test.go | grep -n 'func (f \*runFixture) taskID\|taskIDRe\|rtHook'`
  ⇒ `:56 :67 :123 :134` 四行命中）。本程真去跑了：把 `probes/151/run.head.go`＋`bridge.head.go`（先 `sha256sum` 现核＝`5d46f24^` 逐字节相同）
  铺进树，把那枚测试文件裁成 **T1-only**（T2 连同它独有的 `context`/`time` 两枚 import 一起摘，否则编译单元假死），

```
$ go test -count=1 -v -run 'TestCompositionRootClosesTheLoopTasksTaintScope' ./cmd/wisp/      rc=1
    task_scope_close_151_test.go:17: 任务结束时没有关闭它自己的 C25 污点 scope：stderr 里找不到以
      "[audit] tools: C25 scope closed task=d51b8723-… " 开头的审计行
--- FAIL: TestCompositionRootClosesTheLoopTasksTaintScope (1.80s)
```

⇒ **T1 在未修码上编得过、且直接红**（`probes/151-accept/mutC-t1-on-prefix-code.txt`）。
票面 AC#2 点名要的那一发"这一发在**未修码**上响不响"，对 T1 本可以**拿未修码当尺直接答**，
实现件却用"两枚都跑不起来"把它整体推给了变异尺＋§2.1 的等价读数。
**这不改变判定**：变异尺本程复到了、T1 的未修码直读也确实是红的，钉的牙比它自陈的更硬。
但**它这句话是一句会被下一程当尺用的过度陈述**（例如"票 154 那两枚用例也不用拿未修码量了"），
⇒ 列进"要更正措辞"那一栏，与本程总判的档位挂勾（第 7 节）。
⚠ 同一枚 §6 里它把改前那一跑做成 `overlay-pre.json`（两枚用例**都**摘掉）是**另一件事、做法没错**：
那一跑要的是"改前 rc=0"，留着 T1 就会红。本程这句只针对 §3 那句推理的射程，不针对它 §6 的改前尺。

### 3.4 顺带一条同义反复检查（它没写、本程自己加的）

T2 是**测试自己手动**调 `rt.admitTask("host-151-taint")()` 走那一发边界，不是让环路自己 defer。
这种形状最容易长成"把被测函数写进断言里"的假绿 ⇒ 本程的判据照旧是摘一味：
味 1 拿下后 T2 仍红（judged L2）⇒ 它钉的是**效果**，不是自己调自己。这一条**过**。

## 第 4 节 格 AC#3 —「谁拥有关闭」三选一（裁的是"组合根 defer"）

票面这一格给了两条硬约束：**(a) 三选一先裁再写**；**(b) 不许顺手把 `OpenTask` 也改了来"让两头对称"**。两条本程都核了。

### 4.1 它选①的那两条"现量理由"，本程各换一把尺复到

| 它的理由 | 本程的尺（现读于 `4e16976`） | 读数 |
|---|---|---|
| (a) 这枚 hook 是组合根里唯一拿得到环路 task id 的边界 | `grep -n "newTaskID" internal/agent/loop.go` ＋ 全包 `grep -rn "newTaskID" --include='*.go' internal cmd \| grep -v _test` | 调用点只有 `:322`（`RunAsync`）与 `:333`（`Run`），声明在 `:1107`——**小写、包内可见**，包外拿不到 ⇒ 这一条不是"它推的"，是**符号可见性给的硬界** |
| (b) 这条边界本来就 `defer`，所以 brake/cancel/error/panic 四条形都走它 | `sed -n '355,370p' internal/agent/loop.go` | `:361` `if l.opt.AdmitTask != nil {` → `:366` `defer revokeAdmission()`——**admit 之后立刻 defer**。非测试码里 `AdmitTask` 只有 `cmd/wisp/run.go:586` 一枚赋值（尺：`grep -rn "AdmitTask" --include='*.go' internal cmd \| grep -v _test` ⇒ `loop.go:161/169/361/362/761/777` ＋ `run.go:21/586` ＋ `bridge.go:643` 注释），且 `loop.go:777` 那一枚只是 **nil 检查**、不产生第二条 revoke 路 ⇒ "恰好一枚、长在 defer 上"成立。panic 那一形本程按 Go 语义读（defer 在 unwinding 里跑），**没有真造一发 panic 去量**（见第 8 节） |
| ②"桥内部自管"被否的理由（`Execute` 返回≠任务结束） | `grep -n "sem := make(chan struct{}, guard.Concurrency())" internal/agent/loop.go` | `:604` 命中 ⇒ 同一轮的工具跑在并发闸上（D38d 多路），同一任务多轮 ⇒ 桥确实不知道"任务什么时候完"，②要成立得新造契约。**否得成立** |
| ③"任务生命周期事件"被否的理由（今天没有这个面） | 本程 §2.1 那两枚 `agent.New(`/`loop.Run(` 唯一调用者读数 ＋ `run.go` 的 `onRuntime` 只被测试填 | 桥不在事件那条边上 ⇒ 选③要先造一条总线。**否得成立** |

⇒ **"三选一先裁再写"：本程判它裁对了，且裁的依据是真读数不是偏好。**
"零新生命周期、零新接口"这句本程跟着核过：`admitTask` 只是把已有的 `rt.gate.AdmitTextTask` 包一层、
revoke 里多一发 `rt.bridge.CloseTask(taskID)`，`agent.Options` 的形状一字未动
（尺＝`git diff 5d46f24^ 5d46f24 -- cmd/wisp/run.go` 全文逐行读过：唯一的行为改动就是 `AdmitTask: rt.admitTask` 那一行换掉 `rt.gate.AdmitTextTask`）。

### 4.2 那句"不许顺手改 `OpenTask`"：逐字节核过，**没犯**

```
$ git diff 5d46f24^ 5d46f24 -- internal/tools/bridge.go | grep -E '^[+-]' | grep -n "OpenTask\|OpenScope"   ->  空
$ git diff 5d46f24^ 5d46f24 -- internal/tools/bridge.go | grep -c '^+'                                     ->  17（含 `+++ b/…` 那一枚头 ⇒ 真实新增 16 行）
$ git diff 5d46f24^ 5d46f24 -- internal/tools/bridge.go | grep -c '^-'                                     ->  4 （含 `--- a/…` 那一枚头 ⇒ 真实删 3 行，全在 CloseTask 的旧注释里）
```

⇒ diff 里那枚 `@@ … func (b *Bridge) OpenTask(taskID string) {` 是**上下文行、不是被改行**（本程专门把它挑出来看了）。
`OpenTask` 函数体一字未动 ⇒ **射程没扩大，这一条过**。

### 4.3 它"刻意不做"的那三件事，本程判哪几条真是射程外

1. **不在任务开头 `OpenScope`**：它给的理由是"那会把 unbound 的 fail-closed 换成 fail-open，且正主在
   `internal/risk/**` 那面零字节墙上（`DEFERRED(C25-loop-wiring)` 第 (1) 项）"。本程现量：
   `CloseTask` 新用的那枚 `prov.ScopeTaints` **在 `5d46f24^` 上就存在**
   （`git grep -n "func (p \*Provenance) ScopeTaints" 5d46f24^ -- internal/risk` ⇒ `internal/risk/provenance.go:412`）
   ⇒ 它没有为了加一行审计去动 `internal/risk` 一字。这一条"留给编排者判"**判得对**，且**票 154 已收下**。
2. **不给桥加"列出当前开着哪些 scope"的公开 API**：本程跟着成立——AC#1 那枚探针是经 `-overlay`／包内文件量的
   （本程 §2.2 照同一形状复跑），不需要新面。
3. **不在 `CloseTask` 里顺手删 `b.seqs`**：本程现读函数体（`bridge.go:655-670`）只有
   `ScopeTaints` → `delete(b.scopes, taskID)` → 条件 `CloseScope` → `b.log` ⇒ 没夹带。

⇒ **AC#3 这一格：成立，无条件。**
