# 票 315 探针腿 315-a1 —— AC#0 生产面活体复量（结论：**复现**）

- 腿身份＝探针（只读复量＋一枚 zz 探针载具）；⛔ 不翻框、⛔ 不动票面 Status。
- 锚（现读 `git rev-parse --short HEAD`）＝`a5b2938c`（派单同锚）。载具＝新建 `cmd/wisp/zz315a1probe_test.go`（只加不改）；
  读数日志＝`logs/10-caseA-root-control.txt`（5,925 B）＋`logs/20-caseB-sibling-child.txt`（17,615 B，留盘）。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOTMPDIR=C:\tmp\wisp315a1gotmp`；`GOFLAGS=`；
  判据用 `-run` 小射程 `-count=1`；零 push、零 `add -A`。
- **结论先行：① 复现**。子代理的兄弟孩子（同一父任务的第二枚子代理）**自己**的首次走桥调用
  （一次白名单内 `fs.read`，与兄弟同路径）被真桥判 **L2 / R4 / unbound-scope**，落真审批门（`run.go:612` approval.New）
  → 真卡 → 生产默认路径 31s 超时自动拒；对照腿（同一枚子代理的第一次调用、尚无任何 mark 时）＝L0/allow/success。
  复现 run 的 `wisp run` 自身 rc=0（任务 completed），探针判据 rc=1（它断言「白名单内读取不应变假 L2」）。

## 1. 载具与形状（为什么是 golden，不是 mockllm）

- 一发**真跑组合根**：真 config.toml → 真 secret.Store → 真 approval.Gate → 真 tools.Bridge（真 C25 引擎）→
  真 agent.Loop；模型侧用 **golden SSE 接缝**（`AGENTS.md §1.3` 明列；`internal/llm/golden` 真 HTTP + 真适配器），
  所以**工具调用是模型发起的**（不是测试替孩子发）：根环路发 `task.spawn`，**孩子自己的环路**发 `fs.read`。
  孩子的工具面＝父桥（`internal/tools/subagent_197.go:296`），dispatch 带**孩子自己的 task id**（`internal/agent/loop.go:676`）。
- ⛔ mockllm 载具过不了这一形：它的 chat 规划器只在 `req.ToolChoice != nil` 时才发 tool_call
  （`tools/mockllm/chat.go` planChat），而 agent 环路从不设 ToolChoice（`internal/agent/*` 非测试零命中），
  openai-chat 适配器在 nil 时整字段不发（`internal/llm/openaichat/request.go:121`）⇒ 子环路永远不发工具调用。
  四枚 `subagent_*_197` 载具的 host-dispatch 接缝（如 `startChildWrite197`）用的是**测试自己选的 id**，
  不是孩子环路自己 dispatch 的那一枚（那是票 154 的「宿主自带 id」形）；`TestHostDispatchThroughTheAssembledBridge`
  又在 `rtHook`＝`execute()` 之前、任何 mark 之前。**均未跑**：不是量不到，是不属本形（此处点名，供裁）。
- 最小形（golden 七发，顺序确定性来自 `task.spawn` 阻塞等待子代理 join）：

  | # | 发信方 | 模型动作 | 时刻上的关键事实 |
  |---|---|---|---|
  | ① | 根 | `task.spawn(甲)` | 全树尚无任何 mark |
  | ② | 甲环路 | `fs.read(白名单内 note)` | 甲的首批走桥调用＝**对照腿** |
  | ③ | 甲环路 | 纯文本结论（19 runes） | join 时 `stampConclusion` 把 `task.output` mark 落**父名**（`subagent_197.go:487`） |
  | ④ | 根 | `task.spawn(乙)` | 根 scope 已有甲结论 mark；该 spawn 自身仍 L0/allow |
  | ⑤ | 乙环路 | `fs.read(同一 note)` | **被测腿**：乙 scope 未注册＋别处（根）有 mark ⇒ R4 unbound |
  | ⑥ | 乙环路 | 纯文本收尾 | 乙仍 completed（工具失败不杀任务） |
  | ⑦ | 根 | 纯文本收口 | run 结束 |

## 2. 场景 A：无别处 mark 的对照（根环路首次走桥）——rc=0

- 命令（起 2026-10-14 10:46:45 +08 → 止 10:46:59 +08）：
  `GOFLAGS= go test ./cmd/wisp/ -run 'TestZZ315A1RootLoopFirstBridgeCall' -count=1 -v` ⇒ **rc=0**，`ok … 3.076s`，`PASS (0.31s)`。
- 读数（`run rc=0 in 214.1ms`；golden requests=2，恰好脚本两发）：
  `tool=fs.read level=L0 decision=allow outcome=success corr="028305c3-…#call_zz315_a1_root_read"`。
- 意义：`scopeMarks` 的 fail-closed 是**有条件**的（“别的 scope 有 mark”才开火）；无 mark 时未注册 scope 首调用＝L0。
  这也证明 golden 载具确实驱动了真环路发起真桥调用（不是掏桥自演）。

## 3. 场景 B：兄弟孩子首批走桥调用——**rc=1（复现）**

- 命令（起 2026-10-14 10:47:20 +08 → 止 10:47:59 +08）：
  `WISP_315A1=1 GOFLAGS= go test ./cmd/wisp/ -run 'TestZZ315A1SiblingChildFirstBridgeCall' -count=1 -v`
  ⇒ **rc=1**，`FAIL TestZZ315A1SiblingChildFirstBridgeCall (31.20s)`，包 `FAIL … 33.519s`。
  （首过 10:44:15→10:44:59 同命令同形；键读数在两跑间逐字同族，下述逐字句取入库跑。）
- 该 run 自身 `wisp run` **rc=0**：`任务 b98af16f… 结束（completed，3 轮，2 次工具调用）`——假卡不杀任务，只吞掉那次调用。
- 逐字红句（stdout，审批卡）：
  ```
  [确认 L2 fs.read] R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容
    规则 R4
    影响路径：C:\tmp\wisp315a1gotmp\TestZZ315A1SiblingChildFirstBridgeCall3876597418\001\zz315a1-note.txt
    窗口 0.0s
    卡片编号：78892556-bebc-4e2c-9ddd-4f76712e9296#call_zz315_a1_yi_read
  [warning] 审批将在 30 秒后自动拒绝，请尽快确认
  [dismissed] 审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝
  ```
- 逐字红句（stderr，audit；`[audit] ` 前缀省略）：
  ```
  wisp run: 风险判定 tool=fs.read level=L2 rules=[R4] reason="R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容"
  approval: queued corr=78892556-…#call_zz315_a1_yi_read task=78892556-… tool=fs.read depth=1
  approval: WAITING-STATE state="AwaitingApproval" awaiting=true corr="78892556-…#call_zz315_a1_yi_read" tool=fs.read level=L2 paths=1 producer="ui-prompt"
  approval: ANSWER-EXPIRED corr=78892556-…#call_zz315_a1_yi_read tool=fs.read decision=timeout->reject after=31s
  tools: call kind=approval-timeout task=78892556-… corr=78892556-…#call_zz315_a1_yi_read tool=fs.read risk=L2 decision=timeout outcome=error rules_hit=[R4] in_allowlist_scope=true grant_id=0 reason="R4: 包含来自 unbound-scope …"
  ```
  ⚠ 同一行里 `in_allowlist_scope=true`：**调用在白名单内、仍被判 L2** —— 这条是「假阳性」最干净的自证。
- 落盘行（真 store，`ListToolCallsByTask`）：
  - 甲（对照）：`{ID:1 Tool:fs.read RiskLevel:L0 Decision:allow Outcome:success corr=05c26723-…#call_zz315_a1_jia_read}`
  - 乙（被测）：`{ID:3 Tool:fs.read RiskLevel:L2 Decision:timeout Outcome:error ErrorClass:user_rejected corr=78892556-…#call_zz315_a1_yi_read}`
  - 根两次 `task.spawn` 均 `level=L0 decision=allow outcome=success`（父自己的调用不中）。
- 因果链读数（卡挂起瞬间 5ms 级快照；`ScopeTaints` 只读视图，不落内容）：
  ```
  root  b98af16f-…: [tool=task.output origin="subagent:05c26723-…(甲)" runes=19]   ← 首枚子代理结论落父名
  run:project-instructions: marks=0（本 run 装载器 workspace=""，audit：没有找到任何项目说明文件）
  child 甲: marks=0（随其任务边界关闭：audit `C25 scope closed task=05c26723-… was_open=true dropped=1`）
  child 乙: marks=0（从未开：audit `C25 scope closed task=78892556-… was_open=false dropped=0`）← 未注册 scope 的硬证
  ```
  同刻 roster：根行 + 乙 `status="Thinking"`（**在飞**的兄弟孩子）+ 甲 `status="Settling"`。
  run 收尾处根 scope 关闭行 `dropped=2`（甲、乙两条结论 mark 都落了父名）——链 #4 双向坐实。
- 静态链在锚上复核（行号现读 `a5b2938c`）：`internal/tools/bridge.go:649`(mark 里懒开)/`:847`(OpenTask 定义)；
  `cmd/wisp/run.go:960`(admitTask 只 CloseTask)、`:1052`(项目说明 scope)；`internal/agent/loop.go:676`(dispatch 带自己 task id)；
  `internal/tools/subagent_197.go:487`(结论 mark 落 parentID)；`internal/risk/provenance.go:694`(scopeMarks fail-closed)。

## 4. 判据细节与量法可重跑

- 探针断言的**是应有行为**（白名单内 `fs.read`＝L0；不该出卡）⇒ 现树 rc=1＝复现；修好后同命令应翻绿。
  载具以 `WISP_315A1=1` 门控（无该变量按名 skip），所以包级默认跑不受一枚「以响为存在理由」的探针影响。
- 31s 来源＝`[risk] confirm_timeout_sec` 校验下限 31（`internal/config/validate.go:128`）；载具把
  `per_tool_timeout_ms=60000` 抬到卡之后，否则桥默认 30s 预算会先一步把调用以 ctx 截止拒掉（读法就变了）。
- 复跑：`WISP_315A1=1 GOFLAGS= go test ./cmd/wisp/ -run 'TestZZ315A1' -count=1 -v`
  （A 当发 ~0.3s；B 当发 ~31s＋编译；两发合计金请求数＝2／7，replayer 记 `Requests` 佐证无暗发）。

## 5. 未试的形与边界（喂 AC#1 与下一位）

- 未试：**并发**兄弟（两枚同时在飞、后跑那枚中招）——串行兄弟已复现，且被测那枚当时确实在飞（roster Thinking）；
  **project-instructions 作种子 mark**（数据目录放一份说明文件即可让首枚走桥任务也中——那会牵出根自己的首批调用，
  本腿未测；本 run 里该 scope marks=0，种子 mark 就是甲结论，正合票面链 #4）；**真人/原生允许**（只量了生产默认超时拒）。
- 未跑 mockllm 四载具与 `run_test.go:245`：见 §1（形不对，不是量不到）——此结论也直接喂「撤票判语」的反面：
  本形**不是**只有 harness 缺，生产组合根可达，且已活体复现。
- 下一位要点：修的三案（ⓐ`admitTask` 补 OpenTask／ⓑ子环路自开／ⓒ桥懒开扩到每次 dispatch）本腿不裁；
  只报读数：本案种子 mark＝**首枚子代理结论落父名**（`task.output`，origin＝兄弟 id，19 runes），
  被测 scope 全程 `was_open=false`；被测调用在白名单内。修后请用同门控命令复核翻绿。
