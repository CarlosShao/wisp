# 票 211 普查腿 c1 — 「要真 8 枚并发子代理」的代价表

现读口径：全部 file:line 取自**当前工作树**（2026-09-29）。票面行号已逐条复算，不一致处标出。
归档件（.scratch/wisp/probes/** 下的）一律标〔归档快照〕，不作为现状证据。
未触碰禁区：frontend/**、design/**、internal/panel/tokens_fourway_test.go。
未跑任何 go test/build/vet、未跑 wisp slo（有写腿在飞 internal/agent/approval、internal/tools、cmd/wisp）。

---

## ① 那枚「4」到底是谁定的 —— **字面写在冻结文本里，不是只在 Go 注释里**

**结论：4 是契约级。抬高它＝改契约，需要人工批准。** 而且不止一处，是**四处**：

### 尺一：docs/PLAN.md 的 D38（并发与线程模型）→ 其中（d）背压，即代码注释所说的「D38d」

docs/PLAN.md:2818 是小节标题：

    ### D38 — 并发与线程模型（**原文只有一行，而 D32 的可达性完全依赖它**）

docs/PLAN.md:2845 是它的第四个子项（＝「D38d」这个代号的来源）：

    **（d）背压**：

docs/PLAN.md:2849-2850 逐字（**这就是那枚 4 的冻结原文**）：

    - **工具并发上限 = 4**（防 LLM 一次发 50 个 tool call 打爆机器；
      也与 D45 的批量聚合配合 —— 4 路并发足以让批量场景快起来）。

⇒ 数字 **4 字面写在 PLAN 的冻结决策文本里**，并且**自带理由**（「防 LLM 一次发 50 个 tool call 打爆机器」）。

> ⚠ **这条理由的射程比实现的射程窄。** 冻结理由说的是「**LLM 一次发 50 个 tool call**」＝
> **单轮模型扇出**的限流；而代码把它实现成**全进程一枚 semaphore、所有任务的工具调用共撞**
> （internal/tools/bridge.go:19-23 明写宿主内部 caller 也必须撞同一枚天花板，否则「天花板只是一句注释」）。
> 丙 这条路能不能走，取决于 owner 认这条理由是**字面**（一个进程同一刻只跑 4 枚工具）
> 还是**射程**（防单轮扇出打爆机器）。我不裁这条，只把两处原文并排放这儿。
> 它**不等于**「进程内最多存在 4 个在跑的东西」——冻结文本对**任务并发**的态度恰恰是相反的，见 ④ 的丁。

### 尺二：docs/specs/SPEC-01-architecture.md（同一条规矩的第二、第三处字面）

docs/specs/SPEC-01-architecture.md:134 逐字：

    - 工具并发上限 **4**（防 LLM 一次 50 个 tool call 打爆机器；与 D45 批量聚合配合）。

docs/specs/SPEC-01-architecture.md:26 逐字（架构图那一行也带数）：

      ├ Agent 循环 / LLM adapter / 工具执行（工具并发 ≤4）

同文件 :136 的小节标题自称契约级：
「### 4.4 关停顺序（D38(e)，契约级，顺序错误即事故）」。
docs/PLAN.md:93 把 PLAN 定性为「47 条决策 + **32 条冻结契约**」；docs/PLAN.md:94 列出 SPEC-00…SPEC-12 十三份切片规格。

### 尺三：docs/specs/SPEC-05-agent-core.md（Agent 内核侧第四处）

docs/specs/SPEC-05-agent-core.md:70 逐字：

      与工具并发 ≤4、任务调度协同。

### 代码侧的拷贝（这些才是注释/常量，不是契约本体）

| 位置 | 内容 | 性质 |
|---|---|---|
| internal/tools/bridge.go:24 | const MaxToolConcurrency = 4（注释 :19-23 自称 D38d，并解释为何在 choke point 执行） | 常量拷贝 |
| internal/agent/budgets.go:47 | const MaxToolConcurrency = 4（注释 :44-46：并发计数，不按 context_window 缩放） | 第二份拷贝 |
| internal/tools/bridge.go:62 | 「Values above MaxToolConcurrency are ignored: **the D38d ceiling is not tunable upward**.」 | 注释。**票面把这句归到 budgets.go:62 —— 现树实际在 bridge.go:62** |
| internal/tools/bridge.go:153-154 | ceiling 越界即回落到 MaxToolConcurrency | clamp |
| internal/agent/guard.go:107-108 | 同一形状的第二道 clamp（cfg.ToolConcurrency，注释 :39-40 钉 D38d） | clamp |

### 关于「D38d」这个代号本身

代码注释里的 D38d 指 PLAN.md:2849 那一行，没问题。但 **「D38d」这个标签在冻结文本里并不专指工具并发**：
docs/specs/SPEC-04-voice-pipeline.md:43 把「有界 channel ≤200ms」也标成 D38d
（逐字：「- 有界 channel ≤200ms（D38d），满则丢帧**并计数**（计数进日志/诊断包）」）。
⇒ 拿 D38d 当唯一寻址键会撞名。**要批准就批准 PLAN.md:2849-2850 那一行 + SPEC-01:134（另 SPEC-01:26、SPEC-05:70 两处 ≤4）**，别批准一个代号。

### 顺带必须报的一件事：**冻结文本压根不承认「子代理」存在**

这不是「池 4 还是 8」的问题，但它是 8 的**前置契约面**，票面没写，我复算时撞上：

- docs/PLAN.md:1544：「| **REJECTED** | 多 Agent 协作 | 草案 6.6 明确不做 | — **不是债务** | — | 单 Agent，无编排（⚠ 注意与 D31 的**多任务并发**区分：并发是同一个 Agent 跑多个任务，不是多 Agent 协作） |」
- docs/specs/SPEC-05-agent-core.md:190：「- 不做意图分类前置层（D11 已废）；不做多 Agent/子代理（草案 6.6）。」
- docs/PLAN.md:1123（DSH 对照行）末尾：「✗ 我们不做子代理（草案 6.6 已拒多 Agent）」

⇒ 票 197 已经**在盘上把一行 REJECTED 变成了功能**（internal/tools/subagent_197.go 整件存在）。
而 PLAN.md:1544 那句 ⚠ 给了 丙 一条活的出路：冻结文本**认可的是「同一个 Agent 跑多个任务」的多任务并发**（D31），
被拒的是「多 Agent 协作／编排」。⇒ 「8 枚 child 各是一条并发任务、共用同一枚工具天花板」
比「8 枚子代理各占一枚工具许可」**离冻结文本更近**。这是本腿能给的最强一条**降低批准面**的读数。

**① 的答案（一句话）**：4 **字面写在冻结文本里**，逐字为 docs/PLAN.md:2849
「- **工具并发上限 = 4**（防 LLM 一次发 50 个 tool call 打爆机器；」，
另有 SPEC-01:134／SPEC-01:26／SPEC-05:70 三处 ≤4。**抬高它是契约改动，不是代码改动。**

---

## ② 那枚 30 秒是谁定的 —— **它不是契约，而且票面的语义读反了**

### 常量住在哪

- internal/tools/bridge.go:28 「const DefaultToolTimeout = 30 * time.Second」，
  注释 :26-27 自称「the C22 budget when neither the tool declaration nor the composition root names one」。
- 回落口：internal/tools/bridge.go:156-159（defTm 小于等于 0 时用 DefaultToolTimeout）。
- 单工具声明优先：internal/tools/bridge.go:550-555 timeoutFor（entry.Decl.Timeout 大于 0 就用声明值）。
- agent 侧同名旋钮：PerToolTimeout（见 internal/agent/loop_golden_test.go:167 的用法、internal/agent/forensics_test.go:95）。

### 它**不在**任何冻结件里

- **C22 的冻结原文不给数**：docs/PLAN.md:1372 逐字（节选）
  「| **C22** | **LoopGuard** | ⚠ **梯度式，非硬截断**：重复调用检测阈值 [3,5,8] 梯度注入提醒（学 DSH repeat-tool-reminder）· **每工具 timeoutMs** · 单任务 token 预算 · 轮数上限仅作**高兜底（50）**。」
  ⇒ 契约只要求「每工具要有 timeoutMs」，**30 这个数是代码自己定的默认值**。
- **PLAN 里所有 30s 都不是工具超时**（逐条排掉，防误引）：
  PLAN.md:1368 审批队列**超时 300s、超时前 30s 醒目提示**（C18，与工具执行无关）·
  PLAN.md:2204／:2215／:3072／:1320 是 Conversation 态**无语音 30s 回落 Warm**（隐私/麦克风，与工具无关）·
  PLAN.md:2664 是外部命令类工具默认 timeoutMs = 15000（**15 秒**，另一枚数）。
  grep -n 键词 30s／30 秒／timeoutMs 于 docs/PLAN.md 的全部命中即上面这些，**没有一条是「工具 30s」**。
- **thresholds.go 不涉**：internal/observe/thresholds.go 只有 D32 的内存/CPU 门禁
  （:19-24 memCapSleeping 25MB … memCapWorkPeak 700MB；:31 cpuLimitWorkPeak = 30.0 ← 这是**百分比 30.0%，不是 30 秒**，
  形似而实异，别被它骗）。里面**没有**工具并发数、**没有**工具超时。
  ⇒ 30 秒与 4 **都不在 thresholds.go 里**；AGENTS.md:27 那条「SLO 阈值 / golden / thresholds.go 一字节都不许动」
  **不覆盖这两枚数**（但它覆盖 D32 的 700MB 那一行，见 ④ 乙 的副作用）。
- **allowlist.txt 不涉**：tools/d22scan/allowlist.txt 全文 8 行，只有
  pathresolver-bypass／mirror-hash／bare-goroutine 三类豁免（:3-7），无并发项。
- **golden 不涉**：internal/agent/testdata/golden/ 的 11 个 .sse 全是 **LLM 线上报文**
  （six-tools.sse、parallel-tools.sse、slow-tool.sse、max-tokens-toolcall.sse 等），
  钉的是发给 provider 的请求体，**不编码并发数**。改天花板不会脏 golden。

### ⚠ 现读纠正票面：**排队时间不算在这 30 秒里**

票面 :23-24 断言「**排队时间算在这枚预算里**，所以多出来的那几枚最常见的结局是**超时失败**」。
现树**不是这个形状**。看 internal/tools/bridge.go:430-449 的次序：

    438    select {
    439    case b.sem <- struct{}{}:
    440        defer func() { <-b.sem }()
    441    case <-ctx.Done():
    ...
    449    ectx, cancel := b.execContext(ctx, timeout)

execContext（:542-547）就是 context.WithTimeout，而它在**第 449 行、即拿到许可之后**才建。
⇒ 30 秒**从执行开始算**，**不含等 sem 的时间**。排队那一枚**不会**因为排队而变 OutcomeKindTimeout。

那它会怎样？它**没有上界**：唯一出口是 :441 的 ctx.Done()。而子代理那条 ctx 是
internal/tools/subagent_197.go:327 的 WithCancel(WithoutCancel(ctx)) ——
**WithoutCancel 把父调用的 deadline 一起丢了**（注释 :309-312 自己写明 drops its cancellation/deadline，
这是票 197 §0 冻的「不许级联取消」）。
⇒ 现树形状＝**无限期挂住**，不是「30 秒红」。
⇒ 票面的「**7 of 8 红**」按现树**不可能由 30s 超时产生**；它要么来自旧形状，要么来自别的截止（父任务级／包测试截止）。
**这条按现读为准，并具名报出与票面的冲突，不替票面圆场。**

### 到期时模型看到什么（bridge.go:475-483，第二道同样的检查在 :495-503）

- outcome kind：OutcomeKindTimeout（internal/tools/bridge.go:482、:502）
- 文本：「工具 %s 超时（%dms），已协作式中止」（:478-479，两处一致）
- IsError: true、ErrorClass: string(observe.ClassTool)（:480-481）
- 顺序讲究：:471-474 注释说明**先查 deadline 再看 err**，否则会把 C22 的结构化措辞降级成 internal fault。
- 对照另三种：**取消** 「任务已取消，调用未执行」＋ OutcomeKindCancelled ＋ ClassCancelled（:441-447）；
  **工具内部故障** 「工具内部故障：」＋ err ＋ OutcomeKindError（:484-493）；
  **审批超时** 「审批超时未确认，已自动拒绝」＋ DecisionTimeout → OutcomeKindApprovalTimeout（:402-406、dispositionOf :424-425）。
  ⇒ 四者互不相同，模型能自纠的是前两类；**挂住那一枚什么都不产生**（这才是真症状）。

---

## ③ 谁钉着「4」这个数字（把 MaxToolConcurrency 改成 8，谁红？）

### 会**当场红**的（全部因为「必须正好摸到天花板」，而夹具只有 6 枚工具）

| file:line | 判据 | 为什么红 |
|---|---|---|
| internal/agent/loop_golden_test.go:178 | waitForInflight(MaxToolConcurrency, 5s) 不满足则 :181 t.Fatalf（only %d calls ran at once, want %d: the pool is serial） | blocking provider 加 six-tools 夹具最多同时 6 枚；等 8 枚在飞**永远等不到** |
| internal/agent/loop_golden_test.go:202 | tools.MaxConcurrent() **不等于** MaxToolConcurrency 即 t.Errorf（want exactly %d: the ceiling must be **reached, not merely respected**） | **值相等**断言：天花板 8、实测 6 ⇒ 红 |
| internal/agent/forensics_test.go:106 | waitForInflight(MaxToolConcurrency, 3s) 不满足则 :110 t.Fatalf（only %d/%d tool calls started; the cancelled path is untested） | 同 :178 形状（取证腿复用 blocking provider）⇒ 红 |

⇒ **只有这 3 处会红，而且没有一处是因为「4 这个数字被契约要求」**；它们红是因为
**夹具扇出是 6、断言要求「正好摸到天花板」**。修法是加大夹具（把 six-tools 变 ten-tools，
要动 internal/agent/testdata/golden/six-tools.sse 这类报文）＝**改写断言／夹具**，
按 AGENTS.md:27「不许为了变绿放宽任何断言」得谨慎，但严格说它动的不是契约数值本身。

### **不会**红（重要，别虚张声势）

| file:line | 断言形状 | 结论 |
|---|---|---|
| internal/tools/bridge_test.go:196 TestToolConcurrencyCeilingIsFour | 名字带 Four，判据体是 :221 大于则红、:225 小于 2 则红 | **相对断言，改 8 它照绿**（:210 注释 99 must clamp to 4 会变谎，判据不红）。**名字是一枚假钉子** |
| internal/agent/loop_golden_test.go:146 | MaxConcurrent 大于 ceiling 才红 | 不红 |
| internal/agent/loop_golden_test.go:187、:190 | 同为大于形状（:190 是 tools.Started()） | 不红（Started 为 6，大于 8 为假） |
| internal/tools/subagent_197_test.go:404 Test197SubagentPoolNeverExceedsBridgeCeiling | 只在「池 大于 天花板」时红 | 池 4、天花板 8 ⇒ **不红**。它是「池不许超天花板」的常驻钉，**不是 4 的钉** |
| internal/tools/subagent_197_test.go:443、:528 | 同一「池大于天花板」前置守卫（t.Fatalf） | 不红 |
| internal/tools/subagent_197_test.go:460／473／481／485／488／493／505／508／544／560／572／577／599／602 | 全部以 MaxConcurrentSubagents 为计数（名册行数／在跑行数／占位数） | 池不动全绿；池动到 8 时**这些数一起跟着走**，本身不钉 4 |
| internal/tools/subagent_197_test.go:421、:589 | 拒派生文本要含「上限 %d 枚」（常量拼） | 值无关，不红 |
| cmd/wisp/panel_pump.go:202 | state.PoolCap 赋值为 tools.MaxConcurrentSubagents | 载体透传，不钉数 |
| cmd/wisp/subagent_carrier_197_test.go:330 | PoolCap 不等于 tools.MaxConcurrentSubagents 才红 | **相对**：钉「透传不变」，不红 |
| internal/panel/subagent_roster_197_test.go:187 const poolCapGiven197 为 4 | 字面 4 | ⚠ 这是**测试自己喂给泵的 4**，注释 :180-186 明说它测的是载体透传、真正的池数归票 211。改天花板**不会**让它红。同文件 :193 TestAPumpWithoutARosterReaderSendsFourKeys 的 Four 是**四个 JSON 键**，与工具并发无关——**形似而实异，不许引成钉子** |

### 名册／差集尺（票面要的那枚）

- 现树对应物＝ Test197SubagentPoolNeverExceedsBridgeCeiling（internal/tools/subagent_197_test.go:404）
  加在跑计数 RunningSubagentIDs／InFlightSubagents（internal/tools/task.go:385-409）。
- **没有任何「差集」字样的判据**：在 internal/ 与 cmd/ 的 .go 文件里 grep 差集 ⇒ **零命中**。
- AC#6 在本仓被复用成很多别的东西（internal/llm/probe_health.go:156、internal/panel/composer_test.go:592、
  internal/panel/git.go:228、internal/proc/envfork.go:140 等），**没有一枚是并发差集尺**。

### CI / 全仓不变式扫不扫这两枚常量

- .github/workflows/ci.yml:55-56：lint job ＝ gofumpt/vet/staticcheck 加 **D22 七禁静态扫描 tools/d22scan** 加 ban#8 零表情。
  在 tools/d22scan 源码里 grep ToolConcurrency／MaxToolConcurrency／magic ⇒ **零命中**
  ⇒ **没有任何静态仪器扫并发常量**。d22scan 的豁免面（tools/d22scan/allowlist.txt:3-7）只有那三类。
- bare-goroutine 那一禁**间接管得到子代理**：完成 watcher 必须走 registry
  （internal/tools/subagent_197.go:345-346 由 observe.Default.Spawn 起，名为 subagent-finish- 加 id，:336 注释点名 D38b），
  豁免按**文件路径**给（allowlist.txt:7），不按调用形状给。
- scripts/slo-check.ps1、scripts/slo-fresh.sh、.github/workflows/slo-fresh.yml 跑的是 **D32 六态内存/CPU 门禁**，
  不读这 2 枚常量（internal/observe/thresholds.go 里没有它们）。
- **唯一的常驻钉＝那 3 处「必须正好摸到天花板」的判据，而它们钉的是「摸到」，不是「4」。**

---

## ④ 三条路（＋我在代码里翻出来的第四条）各自的代价表

先给一条**全表依赖的关键现读**，它改变甲／乙的性价比：

> **task.spawn 是一枚会阻塞的桥内调用，且它为整枚孩子的一生握着桥的一枚许可。**
> 证据链：internal/tools/bridge.go:439 拿 b.sem；:440 defer 归还；:470 调 entry.Tool.Execute。
> 而 spawn 的 Execute 一路走到 internal/tools/subagent_197.go:366-367（等 done 才 return t.answer），
> 那个 done 由 :348 的 bg.Wait() 才填。
> internal/tools/subagent_197.go:67-72 的注释自己承认这点（one in-flight spawn holds one bridge slot for
> the child whole life, because the bridge keeps the semaphore held across entry.Tool.Execute）。
> internal/tools/subagent_197_test.go:517-519 也承认（now that the pool EQUALS the ceiling, N real children
> also occupy all N bridge slots）。
>
> **推论（票面和这两段注释都没写出来的那一步）**：N 枚在飞子代理**吃掉 N 枚桥位**，
> 给孩子自己的工具调用只剩「4 减 N」枚。池＝4 ⇒ **剩 0 枚** ⇒ 孩子递第一枚工具调用时卡在
> bridge.go:438，而它那条 ctx 没有 deadline（subagent_197.go:327 的 WithoutCancel 把 deadline 丢了），
> **永远起不来** ⇒ **死锁，不是变慢**。
> 现树**没有一条用例覆盖这个形状**：subagent_197_test.go 里找不到孩子递工具的驱动；
> Test197SubagentPoolCapsAtBridgeCeiling（:441）的孩子用的是 fake197Provider（:452-455），
> **只吐 LLM 正文、不递工具**——它 4/4 全绿恰恰是因为它**没测**工具型孩子。

### 甲＝维持今天（池 4 ＝ 天花板 4）

- **文件**：无（现状）。
- **人工批准**：**不需要**。
- **红的移动**：无。
- **诚实的用户可见数**：**名册最多 4 行「在跑」，而「4 枚孩子同时真递工具」＝ 0。**
  1 枚孩子 ⇒ 余 3 枚桥位可用；3 枚 ⇒ 余 1 枚，工具调用被串成 1 路；4 枚 ⇒ **卡死**。
  ⇒ 甲**不是**「慢一点的 4」，而是「名册上 4、能用 3、满池时 0」。
  **这条比票面写的更糟，而且不需要动契约就能测出来。**

### 乙＝把桥的天花板调到 8

- **文件（契约侧）**：docs/PLAN.md:2849-2850、docs/specs/SPEC-01-architecture.md:134（另 :26）、
  docs/specs/SPEC-05-agent-core.md:70 —— 四处 ≤4 的文字。
  **文件（代码侧）**：internal/tools/bridge.go:24 与 :62 注释与 :153-154、internal/agent/budgets.go:47、
  internal/agent/guard.go:39-40 与 :107-108、池常量 internal/tools/subagent_197.go:78；
  另 subagent_197.go:60-77 与 subagent_197_test.go:390-400 那几段**注释会变谎**。
- **人工批准**：**要**。批准面**逐字就是**上面那四处 ≤4（主键 docs/PLAN.md:2849）。
  ⚠ 另外 docs/PLAN.md:2260 那一行把「推理 ＋ 面板 ＋ **工具并发**」的峰值钉在 **≤ 700MB**，
  即**工具并发已写进 D32 的验收行**；而 700MB 在 internal/observe/thresholds.go:24（memCapWorkPeak）
  是一字节不许动的阈值件 ⇒ **乙附带一次 wisp slo 复测义务**（票面 :37 自己点了这条，现读确认口径成立）。
- **红的移动**：loop_golden_test.go:178、loop_golden_test.go:202、forensics_test.go:106 **三枚红**（夹具只有 6 枚工具）；
  subagent_197_test.go:404 **不红**（4 大于 8 为假）。⇒ 要变绿就得扩夹具或放宽断言，后者撞 AGENTS.md:27。
- **诚实的用户可见数**：**还是 0。** 8 枚 spawn 握着 8 枚桥位，孩子的工具调用剩「8 减 8」＝ 0 枚。
  ⇒ **乙花的是契约的代价，买不到它宣传的东西**：它把死锁的门限从 4 挪到 8，
  **没有把「存在」和「执行」解耦**。这是本腿对乙最重要的一条读数：**乙单独走无用**，除非同时做丙。

### 丙＝把子代理槽位与「工具并发许可」**解耦**

**先答那句最关键的：是两枚许可，不是一枚。**

| 对象 | 住在哪 | 形状 | 生命周期 |
|---|---|---|---|
| 子代理池位 | internal/tools/task.go:358 的 TaskRoster.TryAcquireSubagentSlot | **互斥锁保护的整数计数** r.reservedSubagents（:365-370），release 由 sync.Once 幂等回收（:371-381） | 占位于 subagent_197.go:254；释于 :343 的 finish（:347、:351，以及各条拒绝路径 :264／:271／:305） |
| 工具执行许可 | internal/tools/bridge.go:439 对 b.sem 的写入 | **带缓冲 channel**（字段声明 :126；容量由 :152-154 的 ceiling 定，:171 建道） | 只活过 run 一帧（:440 defer 归还） |

⇒ **两个完全不相干的对象：不同包语义、不同数据结构、不同生命周期。**
把它们绑成「一枚」的**不是数据结构，而是一件实现事实**：task.spawn 自己以一枚桥内工具调用的身份，
在等孩子跑完的整段时间里**握着 b.sem 不放**（见本节开头的推论）。
⇒ **所以丙不需要新建第二座桥、不需要碰 MaxToolConcurrency：它只需要让 spawn 别替孩子握着那枚许可。**

- **文件**：internal/tools/subagent_197.go —— 把 :356-368 那段「等 done」改成派发即返回。
  **现成的形状就在同一个函数里**：:357-365 的返回文案已经写全了
  「父任务这一侧已经不等了…子代理没有被级联取消，它仍在名册里，可以单独停它，流键是 …」；
  配合 internal/tools/task.go 的 AttachCancel／DetachCancel（在 subagent_197.go:330、:391 调用）
  与 RunningSubagentIDs（task.go:397），**「父不等、子继续、槽仍占」这根管子已经通了**；
  结论回读走既有的 task.output（docs/PLAN.md:2565 已把 task.output 定成 L0 工具，并规定截断必须给可续读路径）。
  池常量 :78 由 4 改 **8**；subagent_197.go:60-77 与 subagent_197_test.go:390-400 那两段
  **关于「占一枚桥位」的注释必须重写**（否则注释变谎）。
- **人工批准**：**按现读，不需要动任何 ≤4 的冻结文字。** 桥仍 4、DefaultToolTimeout 仍 30s，
  thresholds.go／golden／allowlist 一个字节不碰。
  ⚠ 但有**一处该当面向 owner 说一句**，而它不是并发数：subagent_197_test.go:404 那枚常驻钉的**前提**
  （「池大于天花板＝谎」）在解耦之后**不再成立**；把池改成 8 会让那枚判据**当场红**（8 大于 4 为真）。
  红的是**判据的前提**，不是契约。而**那枚钉正是票 211 自己要求立的**（票面 :31-32：
  新增一枚常驻判据钉住「池不许大于 MaxToolConcurrency」）。
  **拆上周自己立的钉，我建议先要 owner 一句话再动**；这属于「改了验收的形状」，不属于「改了契约数值」。
- **红的移动**：subagent_197_test.go:404（前提失效）＋ :443、:528 两处守卫会在 8 大于 4 上 t.Fatalf；
  :460-511 那批以池为计数的名册断言**会从 4 变 8 并需重跑**（它们不钉值，但 :481／:485／:488／:505／:572／:577
  要真能起 8 枚）；cmd/wisp/subagent_carrier_197_test.go:330 是相对断言，不红；
  internal/panel/subagent_roster_197_test.go:187 是自喂常量，不红。
  **loop_golden_test.go／forensics_test.go／bridge_test.go 全绿（桥没动）。**
- **诚实的用户可见数**：**名册 8 行「在跑」，机器上同时最多 4 枚在真递工具，其余在等 LLM 或在等桥——但谁也不会饿死**：
  因为再没有「为孩子终身握着桥位」那枚占位，「4 减 N」那个减法消失了。
  ⇒ **丙能在桥＝4 的前提下交付「8 枚孩子可见且都能推进」。这就是那条不需要批准的答复。**
  ⚠ 代价说白：父模型的 inline 结论变成「派完就走、之后用 task.output 取」；
  answer()（:400-422）那套措辞要挪位置，票 197 §0 的「结论盖戳进父 scope」（stampConclusion，:424 起）
  触发点也要跟着挪。**这是语义形状变化，不是并发数变化。**

### 丁＝第四条形状：**孩子不建第二座桥，也不当工具调用——它是一条并发任务**

- 现读依据：internal/tools/subagent_197.go:277 把 opt.Tools 设成 newSubagentToolProvider(t.d.ParentTools)——
  孩子拿的是**父桥的包壳**（:509-521 只过滤掉 task.spawn，被藏的名字在 :513），
  **同一个 Bridge 实例、同一枚 b.sem**。⇒ 今天「孩子的工具执行」和「父的工具执行」**已经在抢同一份 4**。
- docs/PLAN.md:1544 那句 ⚠ 因此给出第四条路：「并发是**同一个 Agent 跑多个任务**，不是多 Agent 协作」——
  **多任务并发是 D31 认可的**。⇒ 把孩子做成宿主 TaskScheduler 直接调度的**一条平级任务**
  （而不是父的一块 spawn 工具调用），父用现成的 task.output 读它。
  这条路**同样不碰 ≤4 的冻结文字**，且比丙更彻底（连「占位」这层都不需要）。
- 代价：这是票 197 实体层之上的再一层重构，射程明显大于丙；且 MaxSubagentDepth 为 1
  （subagent_197.go:84，注释 :80-83 说明它靠「孩子目录里没有 task.spawn」结构性实现）要重新论证。
- 结论：**丙是丁的入口，丁是丙的尽头；两者都不需要动那枚 4。**

### 禁区对照（票面 :51-56）

- 调 MaxToolConcurrency／DefaultToolTimeout／C22 预算：**乙越界**。
  （按 ② 的现读，DefaultToolTimeout 其实**不是**契约件，但票面把它列进禁区，所以腿不能自己动。）
- 新建第二座桥／给孩子单开桥实例：**丙、丁都不做这件事**（两者都复用同一个 Bridge，只是不再为孩子终身握许可）。
- internal/panel/tokens_fourway_test.go、frontend/**、design/**：本腿未读未写；①②③④ 的尺没有一枚落在里面。

---

## ⑤ 别家怎么解这道题（D:/work/AI/open source，全程只读，未在其中写任何东西）

**总结论：四家里两家明确把「存在多少孩子」和「同时跑多少工具调用」分成两枚独立的数，
而且孩子那一枚**大于**工具那一枚——这正是 ④(丙) 要的形状。**

### Step-Code —— 两枚分开的数：**8 存在 / 4 在跑**，且孩子是独立进程

- Step-Code/packages/coding-agent/src/features/step-subagent.ts:72 「const MAX_PARALLEL_TASKS = 8;」
- :73 「const DEFAULT_CONCURRENCY = 4;」
- 两枚数的**语义就写在字段文档里**（StepSubagentExtensionOptions）：
  - :185-186 「/** Maximum number of tasks accepted in one parallel call. */ maxParallelTasks?: number;」
  - :187-188 「/** Number of child processes allowed to run at once. */ maxConcurrency?: number;」
  ⇒ **「被接受的任务数(8)」不等于「同时被允许跑的孩子进程数(4)」**——同一文件、同一个枚举、两枚独立的数。
- 落地：:504 maxParallelTasks 取 MAX_PARALLEL_TASKS 的钳制值；:505-507 maxConcurrency 取
  DEFAULT_CONCURRENCY 与 maxParallelTasks 的**较小者**（即并发孩子数可以小于被接受的任务数）。
- 例子层同形：packages/coding-agent/examples/extensions/subagent/index.ts:33 「MAX_PARALLEL_TASKS = 8」、
  :34 「MAX_CONCURRENCY = 4」、:219 mapWithConcurrencyLimit、:645 用它跑 8 里的 4。
- **为什么它不会撞上 Wisp 这个死锁**：孩子是**独立 RPC 会话／进程**
  （step-subagent.ts:404-411 的 runStepSubagentProcess → createSubagentRpcSession → session.runTurn）。
  孩子的工具调用**不消费父的任何工具许可**。
  ⇒ **丙的「派完就归还许可」在家外的等价物是「孩子根本不在父的桥里」。**

### deepseek-harness —— **三枚**分开的数：孩子（最多 16）**大于**工具（10）

- 工具那一枚：deepseek-harness/packages/core/agent-loop/src/constants.ts:6
  「export const DEFAULT_MAX_PARALLEL_TOOL_CALLS = 10」；
  读值 packages/core/agent-loop/src/tool-calls.ts:132；闸口 :200
  （while 条件里的 inFlight.size 小于 maxParallelToolCalls）；
  配置面 packages/core/agent-loop/src/index.ts:334-335（min(1)、default(DEFAULT_MAX_PARALLEL_TOOL_CALLS)、volatile）。
- 孩子那一枚（**独立数据结构、独立闸**）：
  packages/workflow/workflow-ptc/src/types.ts:16 「/** Concurrent agent() ceiling (already auto-resolved; 大于等于 1). */ maxConcurrentAgents: number」；
  :18 「/** Total agent() calls per run (the runaway-loop backstop). */ maxTotalAgents」；
  默认值 packages/workflow/workflow-ptc/src/index.ts:107（maxConcurrentAgents 的 zod 默认为 0）→
  :141-143（0 时取 Math.min(16, Math.max(1, availableParallelism() - 2))）；
  闸口 packages/workflow/workflow-ptc/src/runtime.ts:142-152（注释 Acquire one concurrency slot in FIFO order，
  activeSlots 与 limits.maxConcurrentAgents 比大小，满了就把 waiter 推进 slotWaiters）。
- **关键比例**：**工具并发 10、并发孩子最多 16** ⇒ **它家让孩子的数量超过工具许可数，而且是按机器核数自适应的。**
  并且 maxTotalAgents（总调用兜底）与 maxConcurrentAgents（同时在跑）又是分开的两枚 ⇒
  「存在多少／同时在跑多少／一共跑过多少」**三层**。
  ⇒ 对 ④(丙) 是**最强的外部支持**：**没人把「子代理池」钉在「工具并发许可」上。**
- 顺带：packages/client/ui-conversation/src/client/apply.ts:67 maxConcurrentFileUploads 默认 2，
  service.ts:412 用它做闸 ⇒ 它家连「上传」都是**每类工作单独一枚数**，不共用全局闸。

### minimax-code —— 有**工具**那一枚(5)，**孩子**那一枚**零命中**

- 工具：minimax-code/packages/agent-modules/mcp/src/runtime/types.ts:45 「maxConcurrency: 5,」（在 MCP_DEFAULTS，:41-46，
  注释 :44 为 Per-server semaphore size）；
  runtime/connection-pool.ts:9 「Per-server semaphore bounds concurrent tools/call to maxConcurrency」；
  :72 class Semaphore；:175 取 MCP_DEFAULTS.maxConcurrency 兜底；:401 new Semaphore(this.maxConcurrency)。
- **它是 per-server，不是全局** ⇒ 又一条「分层限流」证据（同一进程内不同服务器的工具各有 5，不是一枚全进程 4）。
- ⚠ **对我们最有用的一条**：connection-pool.ts:77-79 「aborted mid-queue we splice it out — the semaphore never
  grants a permit」；:106 对排队中被中止的 waiter 抛 Semaphore acquire aborted while queued；
  :250-259 「Early-abort: skip semaphore acquisition entirely when the caller has already aborted」。
  ⇒ **它家专门处理「排在队里的调用要能被摘出去」**。而 Wisp 现在的 bridge.go:441-447 只认 ctx.Done()，
  孩子那条 ctx 又**没有 deadline**（subagent_197.go:327）——minimax 这三行正对着我们那个「无限挂住」的洞。
- 孩子／子代理并发上限：**零命中**。只读 grep（模式 maxConcurrent(Subagent|Task|Agents|Children) 与
  concurren(t|cy).*subagent 与 subagent.*concurren 与 parallel.*(task|subagent)，覆盖 minimax-code/packages 的 .ts）
  ⇒ 唯一沾「并发」的是 packages/local-runtime/src/config/types.ts:126 的 channelBridge.lanes 下的 maxConcurrent，
  那是**聊天通道 lane** 的界，不是子代理池；subagent 的命中全在角色枚举
  （packages/local-runtime/src/agent/port.ts:2-4、feature-owned-skills.ts:7、subagent-telemetry.ts:1），**无并发数**。

### openchamber —— **两枚都是零命中**（它家根本没有这道题）

- grep（只读）模式一：maxConcurrent(Subagent|Task|Agents|Children)、MAX_CONCURRENT_(TASK|SUBAGENT|AGENT)、
  subagent.*limit、parallelTask、backgroundTask.*(limit|max)，覆盖 openchamber/packages 的 .js 与 .ts ⇒ **零命中**
- grep（只读）模式二：concurren(t|cy).*(tool call|toolCall)、toolCall.*concurren、parallelToolCalls ⇒ **零命中**
- 目录核对：openchamber/packages/web/server/lib/ 下有 agent-memory、agent-tool、message-queue、projects、quota 等，
  **没有 agent／task 并发闸**。它家仅有的并发旋钮都**不是 agent 的**：
  packages/web/server/lib/git/service.js:2470 createSerialRefresh 传 maxConcurrent（实现 lib/git/serial-refresh.js:11，
  默认 Number.POSITIVE_INFINITY）；lib/skills-catalog/scan.js:158 「const maxParallel = 10」（:207 起 worker 池）；
  renovate.json:7 prConcurrentLimit 为 5（那是 Dependabot 的，不是运行时的）。
  ⇒ **零命中就报零命中**：openchamber 不能给 ④(丙) 提供任何支持。

### ⑤ 的一句话

**Step-Code：8 存在 / 4 在跑，两枚独立数（step-subagent.ts:72-73 加 :185-188），孩子是独立进程所以不吃父的工具许可。
deepseek-harness：10 工具 / 最多 16 孩子 / 外加总数兜底，三枚独立数（constants.ts:6、workflow-ptc/types.ts:16-18、
workflow-ptc/index.ts:141-143），孩子数**故意大于**工具数。minimax：工具 5 且**按服务器分层**（types.ts:45），
孩子数零命中，但它家「排队中可摘除」（connection-pool.ts:77,106,250-259）正对我们「无限挂住」那个洞。
openchamber：两枚都零命中。**
⇒ **没有任何一家把「子代理池」钉在「工具并发许可」上；两家明确把它们分开，一家还让孩子多于工具。**

---

## 本腿具名报出的「没测／测不了」

1. **没跑任何 go test**（按派发令，有写腿在飞），所以 ③ 那张红表是**静态读断言**推出来的，不是实测；
   loop_golden_test.go:178／:202、forensics_test.go:106 三枚红的**因果**（6 枚夹具摸不到 8）未经运行确认。
2. **满池死锁那条推论是代码形状推论，没有实测**：现树没有任何一条用例让孩子真的递工具调用
   （Test197SubagentPoolCapsAtBridgeCeiling 的孩子只跑 fake197Provider）。
   要证它需要新写一条「4 枚孩子各递 1 枚工具」的用例——本腿是只读腿，没写。
3. **票面「7 of 8 红」未能复现，且按现树语义不可能由 30s 超时产生**
   （bridge.go:439 先拿许可，:449 才建 deadline）。这条与票面 :15、:23-24 直接冲突，**按现读为准**。
   要翻案得去读 docs/evidence/s1/197-subagent-entity-r1.md 的原始读数；本腿只引了票面的转述并标出冲突，
   **没引任何 probes 归档件作为现状**。
4. **「4 路并发足以让批量场景快起来」（docs/PLAN.md:2850）是不是 owner 心里的那个 4**——
   冻结文本给了理由，而理由的射程（单轮扇出）与实现的射程（全进程共闸）不一样。这条**不由我裁**，需要 owner 一句话。
5. **票 197 的功能与 docs/PLAN.md:1544 的 REJECTED 行、docs/specs/SPEC-05-agent-core.md:190 相冲**——
   本腿只报现读冲突，不裁它是否已被口头批准过。它决定的是「8」这句话要连着批准到哪一层。
6. **没读也没碰 frontend/**、design/**、任何 .scratch/wisp/probes 归档件**（c1 是本腿唯一落点）。
   internal/panel/subagent_roster_197_test.go:187 的字面 4 已核对为「自喂载体测试」，不是并发钉；
   同文件 :193 的 Four 是四个 JSON 键——两处**形似而实异**，已在 ③ 标明。
