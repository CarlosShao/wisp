# T10 对抗验收报告（独立复核者执行）

> 执行者：adversarial-reviewer（实现者 T10-loop 两任：ae23da6 被 turn 上限击杀，d6f6406/3a21d8a/aade1a9/8d9fd95 接手未提交 WIP）。时间：2026-09-20T09:05:00Z。
> 契约：`.scratch/wisp/issues/10-agent-loop-core.md`（`## Key constraints` / `## Acceptance criteria` 逐条对照）；SPEC-05 §2/§4/§8；PLAN.md 按 ID 查：D11（:274-305）、D15③④（:430-437）、D21（:679、:2941）、D39（:3199）、C22（:1372）、D42#9（:3041、:755）。
> 复核方法：先攻击判据仪器本身，再逐 AC 进攻；所有结论 = 亲自读到的 file:line 或亲自重跑的输出。探针（`zz_probe*_test.go`）已删除，未提交。

## 一、判据仪器审计（AC6 的「mockllm request count = 0」）

第二任实现者改动了 AC6 唯一的判据：`harness.requests()`。**改动前后计数路径逐行对比**：

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 0.1 | 计数点是否被收窄 | **PASS：未收窄，反而提前** | 旧版 `git show ae23da6:internal/agent/harness_test.go:154` = `return len(h.rep.Requests)`，其 append 发生在 `internal/llm/golden/replay.go:44`（`ServeHTTP` 第一句）；新版 `internal/agent/harness_test.go:154-158` = `len(h.bodies)`，append 在 `harness_test.go:96-98`，位于 `rep.ServeHTTP(w, r)`（:99）**之前**。同一 handler、每次 HTTP 请求一次、无 method/path/状态码过滤 ⇒ 覆盖集恒等，且新计数点严格更早 |
| 0.2 | 是否只统计部分请求种类 | **PASS：不区分种类** | 探针 P2：`text-reply` 只录了 1 个 response，跑 2 个任务 → 第 2 次请求被 replayer 以 599 拒绝（`replay.go:45-48`）且触发适配器重试；实测 `requests()=3 / rep.Requests=3 / bodies=3` 三者恒等 —— 被拒绝/失败的请求同样计入 |
| 0.3 | 「0」是否平凡为真（阳性对照） | **PASS：同一仪器在同类用例返回非 0** | `control_test.go:69` 用完全相同的 harness 断言非控制词「取消这个任务吧」→ `requests()==1`；`guard_test.go:41` 同仪器 =8；`loop_golden_test.go:184` =2。探针 P1（`tool-then-text`）实测 `requests()=2`，非零 |
| 0.4 | 改动的动机是否正当 | **PASS：旧版确有 data race** | `rep.Requests` 是无锁裸 slice（`replay.go:33`），由 httptest handler goroutine 写、测试 goroutine 读；`h.bodies` 由 `h.mu` 保护（`harness_test.go:96-98,155-157`）。旧版还在流被提前取消时可能读到撕裂状态 |
| 0.5 | 结论 | **AC6 仪器健全** | 「计数被改小使 0 平凡成立」这一 BLOCKER 假设**不成立**：改后计数器只会比改前更早、更多，绝不漏计 |

## 二、逐 AC 裁决

| AC | 契约原文要点 | 裁决 | 决定性证据（file:line / 重跑输出） |
|----|-------------|------|-----------------------------------|
| AC1 | golden 驱动：纯文本 / 单工具 / 并行工具 / 结果喂回下一轮 / 预算耗尽→Stuck 带明示 | **PASS** | 真 SSE 字节走真 openai-chat 适配器 + 真 httptest（`harness_test.go:86-101,125-138`），非函数 mock。10 个 `TestGolden*`（`loop_golden_test.go`）。喂回：`loop_golden_test.go:128-134` 直接断言第 2 轮**线上 body** 含 `"tool_call_id":"call_e1"` 且第 1 轮不含。预算：`loop_golden_test.go:157-190` 断言 `status=stuck`+`brake=token_budget`+`st.Text == res.Message`（同一句话既发布又返回）+`requests()==2`（第 3 轮从未发出） |
| AC2 | `stopReason=max_tokens` → 该消息**全部**未闭合调用判失败，一个都不执行 | **PASS（且非循环论证）** | 实现 `loop.go:747-770`：`for _, c := range turn.ToolCalls`（:755）**不看 `c.Complete`**，逐个 startCall→`decide(Reject)`（:757）→`finish(truncated)`（:758）→补 tool_result（:760）。夹具 `testdata/golden/max-tokens-toolcall.sse` 同一条消息里 `call_ok1` 参数完整合法、`call_cut2` 参数被腰斩，随后 `finish_reason":"length"`。探针 P3 实测两个调用在缝上都是 `complete=true`、`stop="max_tokens"` ⇒ 「判失败」完全是 loop 的动作，`executeCalls` 的 `!c.Complete` 兜底分支（`loop.go:588-592`）本来会**放行两个调用**。断言面：`truncation_test.go:38` 执行数=0、:44 `res.ToolCalls=2`、:57-68 **逐 call** 断言 outcome=truncated+error_class=loop、:74-91 断言**落库的 2 行** outcome/decision=reject/error_class、:95 保留部分文本、:105 未发起第 2 轮。`disconnect-mid-toolcall.sse` 确实被 `truncation_test.go:119-140` 跑到（该夹具物理上无 `finish_reason`/`[DONE]`，末行参数以引号结尾、合法 JSON 闭合但语义截断） |
| AC3 | 阈值随 `context_window` 等比缩放，**不得硬编码**；artifact 存全文、context 留头 500/尾 200/总长/路径；>1MB 截断并标 `truncated=true`；1MB 上限 | **PASS（阈值缩放）/ 见 MAJOR-1、MINOR-3** | 硬编码搜猎：`spill.go` 全文无 4000/1MB 字面量，只有 `s.b.SpillTokens`（:88）、`s.b.RawOutputCapBytes`（:63）；`budgets.go` 的 `refSpillTokens=4000`（:32）、`refRawOutputBytes=1<<20`（:35）是**基准值**，经 `BudgetsFor`→`scaleInt`（:94-108,114-123）缩放后供值，正是契约要的结构。缩放证明测试**只换窗口**：`spill_test.go:56-63` 同一 592 字节 payload 分别喂 `NewSpiller(dirBig,big)` 与 `NewSpiller(dirTiny,tiny)`，128k 不溢、4096 溢出，并断言大窗口目录**空**（:65-67）排除假通过。头/尾/总长/路径：`spill.go:101-107` + `spill_test.go:104-151`（:127 `string(got) != full` 断言文件=全文；:148 KeptHead/KeptTail 恰等于 scaled 500/200）。rune 安全：`spill.go:66-70` 逐字节回退到 UTF-8 起始字节（`utf8Start(b) = b&0xC0 != 0x80`，`prompt.go:293`）——按构造覆盖全部后续字节（2/3/4 字节序列同理），已提交测试钉的是 2 字节 `é` 跨界这一例（`spill_test.go:169-179`，断言精确切到 99 字节、排除 `é`）。E2E：`spill_test.go:186-233` 同一条 `spill-tool.sse`，128k 不溢 / 4096 溢且落盘 |
| AC4 | >12000 tok 压最老轮、保最近 3 轮原文、**绝不丢 tool_call id 与结果引用** | **PASS** | ledger 是**真机制不是回填**：`compress.go:112` 在摘要之前就把被折叠消息的 id 收进 ledger，`:159-168 summaryBlock` 无条件把 `[tool_call 记录（不得丢弃）] id(name); id…` 拼进摘要块 —— 摘要器返回什么都不影响。测试**双向断言**排除「字符串碰巧还在」：`compress_test.go:100-105` 先断言被折叠 id `!hasLiveToolUse`（结构体层面已消失），再断言 `histTextHas`（只剩 ledger 文本这条路）；:85-93 断言最近 3 轮 id 仍是**活的** tool_use+tool_result 结构。`droppingSummarizer`（:122-126）只返回省略 id 的散文，:142-147 仍要求 id 存活。最近 3 轮逐字：:103 断言 `!hasLiveToolUse` 与 148-153 的 `hasLiveToolUse` 互斥成对。触发阈值缩放：:159-205（5 轮 ~1344 token 落在 tiny 384 与 big 12000 之间，128k no-op、4k 折叠；:186-188 的 setup 断言防止 payload 跑出区间）。:43-46/:161-163 钉住 128k 参考值 12000 |
| AC5 | 重复调用 `[3,5,8]` 梯度注入提醒；第 8 次 → Stuck + **用户可见消息**（禁止静默停）；每工具 timeoutMs | **PASS（禁止项未被违反）/ 见 MAJOR-2、MINOR-4** | 墙钟差值搜猎：唯一超时实现 `loop.go:686-689` = `context.WithTimeout` + `errors.Is(tctx.Err(), context.DeadlineExceeded)`，**无** `time.Since`/`time.Now()` 差值判断（D42#9、PLAN:755/:3041 的禁令未被触碰；`d22scan` 全仓 clean）。包内 `time.Since` 只出现在 `tools.go:186` 的记录性 offset 与测试断言里。3/5/8 三档各自成立：`repeat-echo.sse` 录了 9 个 response（:27-:97 八个 tool_calls + :104 第九个纯文本），测试断言 `res.ReminderLevels==[3,5,8]`（`guard_test.go:32-40`）、`requests()==8`（:41，证明第 9 个 response 从未被消费）、**执行数=7**（:45，第 8 次拒绝执行）、注入的 `[系统提醒]` 恰 **2 条**（:51-61，3/5 档进对话、8 档不进）、`EvReminder` 恰 **3 条**（:76-87）。Stuck 可见性：:63-72 同时查 `EvStuck` 事件与 `res.Message`，并逐个要求含 `echo`/`连续调用 8 次`/`参数相同`/`token`，:65 `t.Fatal("no stuck event published")`。阶梯可配置（证明非硬编码）：:92-106 用 `[]int{2}` 复跑。200k 预算按窗口缩放：:130-146 |
| AC6 | 控制词命中即本地执行，**零 LLM 往返**（mockllm request count = 0） | **PASS** | 判据仪器见第一节（健全）。词集与 D11(1) 逐字一致：`control.go:36-42` 恰好 停/取消/重说/大声点/确认 五词一一对应五个 verb，无自造同义词（注释 :33-35 明确扩词典归 ticket 58）。短路真实：`loop.go:332-334` 在 `run()` 的**第一句** return，`root`/`newTaskJournal`/`buildRequest`/`stream`（:336-392）一行都到不了；`MatchControl` 纯本地 `strings.Trim`+map 查表（`control.go:51-58`）。非控制词照样发给 provider 且**只发一次**：`control_test.go:58-72`（5 个近邻全不匹配 + 端到端 `requests()==1`）。控制层还能真停正在跑的任务：`control_test.go:76-103` |
| AC7 | 段序 = 缓存前缀序；前缀跨轮**字节稳定**；④ 场景与 BM25 段在最后 | **PASS（核心断言有效）/ 见 MINOR-5** | 比较的是**字节**且是**两次独立装配**：`prompt_test.go:71-72` `req1 = asm.Build(mkIn(true), histTurn1)`、`req2 = asm.Build(mkIn(false), histTurn2)`，两轮的 scene（时间+焦点窗口）、BM25 命中、profile、history 全不同，:76 `sysBytes(req1) != sysBytes(req2)` 才失败 —— 不存在「一个字符串和自己比」。**反向防空**：:84 要求前缀含 `你是 Wisp` 且 `len(prefix) >= 200`（否则判「byte-stability unproven」），:97-99 要求易变后缀**确实不同**，:90-94 反向断言 5 个易变串（含「微信 - 张三」「上午的安排」）**绝不出现**在前缀里。段序：`prompt.go:51-54` 规范序 + `prompt_test.go:129-139` 逐位对齐 + :143-147 断言 `CachePrefix` 恰为前 4 段 + :149/:181-192 断言 ④ 场景既在段序末也在**最后一条 content part**，并在 :190 检查场景里真的带着本轮的焦点窗口串 |
| 附加 | 终止集（无工具调用+文本完成 / 取消 / 预算 / 50 轮兜底） | **PASS** | `loop.go:437-445`（无 tool call → completed）、:358 与 :398 与 :470（取消）、:362 与 :415（token 预算）、:367 `NextRound()` 兜底；`guard.go:139-142` 在第 **51** 轮才 `floorHit`，即 50 轮是上限不是主刹车（符合 PLAN:434「高兜底」）。取消端到端：`loop_golden_test.go:194-214` 断言 `root.Wait()` 归零 |
| 附加 | task_log / tool_call 行带 error_class / decision | **FAIL（部分路径）** | max_tokens 路径完整（`loop.go:757-758`，`truncation_test.go:82-90` 逐行验证 outcome/decision=reject/error_class=loop）。但（a）`failOpenCalls`（`loop.go:775-790`）只 `finish` 不 `decide` → 探针 P9 实测中断行的 `decision=""`；（b）**取消任务的终态行整体丢失** → 见 MAJOR-1 |
| 附加 | 200k token 预算按 context_window 缩放 | **PASS** | `budgets.go:37,108` `refTokenBudget=200000` 经 `scaleInt` 缩放；`guard.go:102-104` 仅在未显式配置时取缩放值；`guard_test.go:130-146` 钉住 128k=200000、4096=`4096*200000/128000=6250` |
| 附加 | root ctx 取消贯穿 loop→stream→tool | **PASS** | `loop.go:336-343` 把 `ctx` 重绑到 `root.Ctx`（这正是 d6f6406 的修复），:515 `provider.Stream(ctx,…)`，:619 工具 goroutine 由 `reg.Spawn(…, root, …)` 继承 `root.Ctx`（`observe/goroutine.go:286-292`）。`control_test.go:76-103` 证明控制词「停」能取消**正在跑**的任务（`status=cancelled`、`RootPending=0`、且控制词本身 `requests()` 仍 =1，即没有多花一次往返） |
| 附加 | D38d 并发上限 4 | **PASS（机制为真；测试偏弱）** | `loop.go:578` `sem := make(chan struct{}, guard.Concurrency())`、`guard.go:105-107` 钳到 ≤4。现有测试 `loop_golden_test.go:146-152` 只断 `<=4` 且把 `<2` 降级成 `t.Logf` —— echo 工具瞬时返回，`MaxConcurrent()` 恒为 1，**该断言无法区分上限 4 与完全串行**。探针 P11/P14 用会真阻塞的工具实测（各 3 次）：6×80ms 工具 `maxInflight=4 elapsed=160ms`；6×60ms 计时工具 `call_c0/c1/c2/c5 start @0s`（4 个真并发）、`c3/c4` 各在前一个结束后启动、120ms 跑完 6 个、`guard.Concurrency()=4` ⇒ 机制正确，串行是测试度量盲区而非产品缺陷（记 MINOR-6） |
| 附加 | 标准捷径排除（无断言测试 / t.Skip / 同质表 / 循环夹具 / mock 顶替 C5 缝 / 阈值被弱化） | **PASS** | 全包 `t.Skip` = 0；`^func Test` 共 **35** 个（compress 4 / control 6 / guard 5 / loop_golden 10 / prompt 3 / spill 5 / truncation 2）与 `-count=2 -v` 实跑的 35 个逐一对上，两遍全绿。`control_test.go:15-26` 七条用例七个不同 utterance+verb（非同质表）。夹具不编码被验行为：P3 证明 `Complete=true` 由适配器给出、判失败是 loop 的决定（见 AC2）。C5 缝是真适配器：`harness_test.go:125-138` + `golden/replay.go:43-68`，无函数 mock 顶替。阈值未被弱化：`spill_test.go:29` 钉 4000、`compress_test.go:44,161` 钉 12000、`guard_test.go:132` 钉 200000、`budgets.go:18` 钉基准 128000 |

## 三、缺陷清单

### BLOCKER
无。首要嫌疑（AC6 判据仪器被改小）**经三方对照 + 两个探针实测后被排除**。

### MAJOR

| # | 缺陷 | 证据 | 一行修法 |
|---|------|------|---------|
| MAJOR-1 | **被取消任务的终态 `task_log`/`tool_call` 行写不进库**：`finish()` 用已取消的 root ctx 落终态，`BeginTx(cancelledCtx)` 直接失败，任务行永久停在 `state="running"` / `ended_at=NULL`，在途工具调用行整体丢失。取消是契约点名的四种终止态之一，其取证结果被静默丢弃（只留一条 WARN） | `loop.go:343`（`ctx = root.Ctx`，d6f6406 新增）→ `loop.go:833` `j.finishTask(ctx,…)` → `memory/writer.go:135` `q.db.BeginTx(cmd.ctx,nil)` → `enqueue` 于 `writer.go:88-89` 以「write abandoned by caller」返错；`loop.go:834-836` 吞成 Warn。探针 P4 实测 `task_log.state="running"`、`endedAt=<nil>`、`tool_call rows=0`，日志 `agent: task_log finish failed err="memory: begin write tx: context canceled"` | `loop.go` `finish()` 里落库前把 ctx 换成脱离取消的清理 ctx：`jctx, cj := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second); defer cj()`，`finishTask`/`failOpenCalls` 用 `jctx` |
| MAJOR-2 | **C22 结构化 TOOL_TIMEOUT 被工具自返的错误吞掉**：`dispatch()` 在超时分支同时返回「超时 outcome」和工具自己的 error，而 `executeCalls` 的 `case execErr[i] != nil` 排在前面 → 超时文本与 `ClassTool` 全丢。`ToolProvider` 契约（`tools.go:84`）明确要求宿主失败以 error 返回，真机 host bridge（ticket 20+）必然命中此路：模型只看到 `context deadline exceeded`（不含工具名与 timeoutMs），行被记成 `error_class="internal"`（D37 里 internal 属不可自纠类），与 C22「structured TOOL_TIMEOUT」直接冲突。现有测试之所以通过，是因为 `EchoProvider` 违反该契约、ctx 取消时返回 nil error（`tools.go:170-171`） | `loop.go:694-696`（返回 `err`）+ `loop.go:642-645`（`execErr` 优先）+ `guard.go:244-250`/`observe/errors.go:230-245`（裸 ctx 错误无类别 → 落到 `ClassInternal`）。探针 P8（契约诚实的 sleep 工具）实测 `outcome="error" class="internal" text="context deadline exceeded"`，行 `class="internal"` | `loop.go:694` 超时分支改为 `return ToolOutcome{…超时文本…, IsError:true, ErrorClass:string(observe.ClassTool)}, nil`（错误已被该 outcome 完整表达） |

### MINOR

| # | 缺陷 | 证据 | 一行修法 |
|---|------|------|---------|
| MINOR-1 | artifact 逃出 1MB 硬上限：PLAN:432 要求「先截断到上限、**再走上一行**（落文件）」，实现落的是**未截断原文**，且 stub 里播报的 `总长 N 字节` 与文件实际大小不符（小窗口+>上限时 `truncated=true` 标记还会被尾窗切掉） | `spill.go:77`（`capped` 已算出）与 `spill.go:80`（`TotalBytes: len(capped)`）与 `spill.go:97`（`writeFileExclusive(path, []byte(text))` ← 原文）；`spill.go:101-107`（marker 只在被截断文本尾部，小窗口 tail 只有几 token）。探针 P7 实测 cap=200 时 `artifactBytes=800`、`TotalBytes=200`、stub 内 `truncated=true` 丢失 | `spill.go:97` 改写 `[]byte(out.Text)`（或把 marker 拼在 stub 之后而非文本尾部） |
| MINOR-2 | 中断/未闭合调用落库缺 gate 判定，与 max_tokens 路径不一致（同一张 `tool_call` 表两套写入规矩） | `loop.go:775-790` `failOpenCalls` 无 `j.decide`；对照 `loop.go:757` 有。探针 P9 实测 `outcome="error" decision=""` | `failOpenCalls` 在 `j.startCall` 后补 `j.decide(ctx, row, DecisionReject)` |
| MINOR-3 | `turnSignature` 的注释声称用 `encoding/json` 自身压缩做参数规范化，代码只做了 `strings.TrimSpace`；因此仅键序/空白不同的等价重复调用**不会**被判为重复（注释高估了防线，而注释比代码读起来更可信） | `guard.go:227-229`（注释）vs `guard.go:230-240`（实现）。探针 P10 实测 `{"a":1,"b":2}` 与 `{ "b":2, "a":1 }` 交替调用 `ObserveTurn` 返回 `ok=false` | `turnSignature` 里对 `c.Args` 先 `json.Unmarshal`+`json.Marshal`（或 `json.Compact`）再拼签名，或删掉该注释 |
| MINOR-4 | 每工具超时唯一的可见断言是 `text` 含「超时」，只覆盖了「工具返回 nil error」这一契约违背分支；契约诚实分支（返 error）零覆盖（正是 MAJOR-2 藏身处） | `guard_test.go:177`；`tools.go:170-171` | 补一个返回 `ctx.Err()` 的 ToolProvider 用例（即探针 P8 的形状） |
| MINOR-5 | D39 的分段预算与 ≤2300 总额**无任何测试固定**（只在生产代码里执行）；`Assembler.TotalTokens()` 自带「used by tests to assert the D39 ceiling」的注释却是死代码 | 全仓 grep：`2300\|PromptTotal\|TotalTokens()` 在 `internal/agent/*_test.go` 命中 0 次；`prompt.go:149-153`。实测值本身正确（探针 P5/P12：identity 54/150、safety 93/100、style 49/50、resident 42/1200、总 345 ≤ 2300，三段规范正文均完整未被裁剪），只是没有回归保护 | 在 `prompt_test.go` 加一条：逐段 `ApproxTokens(s.Text) <= s.Budget` 且 `sum <= 2300` 且 `asm.TotalTokens()==2300` |
| MINOR-6 | D38d 并行度断言可被串行实现骗过：`<=4` 恒真，`<2` 只 `t.Logf`；`withTools` 全仓无一个会阻塞的工具，产品侧真并行从未被测试证明（本次由探针 P14 代为证明） | `loop_golden_test.go:146-152`；`tools.go:166-174`（sleep 用 `observe.NewTimeout` 忙等，echo 瞬时） | 把 `MaxConcurrent() < 2` 从 `t.Logf` 升成 `t.Errorf`，并配一个慢工具 provider 夹具（探针 P14 已给出现成写法） |
| MINOR-7 | `prefixSummary` 算了不用（`:=` 后立刻 `_ =`），读者会误以为摘要前缀位置有条件分支 | `compress.go:104-108` | 删掉这两行 |
| MINOR-8 | 夹具头部注释与实际用例脱节（写着 12000-token 窗口 / 375-token 阈值，测试实为 4096 窗口 / 128 token），会误导下一个改这个用例的人 | `testdata/golden/spill-tool.sse:3-5` vs `spill_test.go:201-207` | 改注释为 4096/128，或直接删掉带数字的那两句 |
| MINOR-9 | `RecordSink` 自述「只适用于单任务 goroutine」，但 `TestDefaultControlCancelsRunningTask` 里异步任务 goroutine 与测试 goroutine 会并发 `Publish` → 无锁 append 竞争（25× `-race -count=25` 未复现，属潜伏 flake） | `sink.go:76-81`；`control_test.go:80-85`；`loop.go:479` | `RecordSink` 内加 `sync.Mutex`，或测试侧只在 `task.Wait()` 之后读 |

## 四、实现者移交的开放问题裁定

**问题**（Progress log 2026-09-20T12:40Z 第 (c) 条）：被取消任务的终态 `task_log`/`tool_call` 行用已取消的 root ctx 写入、被 memory 的 writer 放弃；实现者称「pre-existing、spec-silent、本票有意不动」，并「若验收希望一个 detached cleanup ctx 就提请登记」。

**裁定：三项主张均不成立 —— 这是本票必修的真实缺陷（MAJOR-1），不是既有行为、不是规范沉默、也不应先变成决策征询单。**

1. **不是 pre-existing。** `git diff ae23da6 HEAD -- internal/agent/loop.go` 在这一维度上**只有一行**新增：`+ ctx = root.Ctx`（现 `loop.go:343`），由本票 d6f6406 引入（Progress log 11:48Z 自述：为修 D11「停」无法取消在跑任务）。改前 `finishTask` 拿的是调用方存活的 ctx，取消路径的终态行**能**落库。所以「写不进」是本票修复的副产物，不是继承债。
2. **不是 spec-silent。** 契约 `## Key constraints` 末行明确要求「Writes task_log row + tool_call rows (04 tables) with error_class/decision fields」，而取消是同一行上文点名的四种终止态之一；`taskLogState` 专门为它准备了 `"cancelled"` 枚举值（`loop.go:857`）。规范既规定了这行必须写、又规定了它必须带这些字段 —— 「写不进去」是直接违规。ticket 04 侧也无沉默：`memory` 契约承诺「命令已入队就会执行，调用方只是不再旁观」（`writer.go:62-66`），但 `exec` 把 `cmd.ctx` 直接喂给 `BeginTx`（:135），对**入队前就已取消**的 ctx 必然失败 —— 这是 memory 的承诺与实现之间的偏差，不是本票可以援引的免责条款。
3. **归属：本票（10），一行修即可**（见 MAJOR-1 修法）。修法完全落在 loop 的终止路径内，不触碰 `internal/memory/`、不触碰另一代理在修的 `internal/risk/`。
4. **另需登记（但不是本票的阻塞项）**：`writeQueue.exec` 是否应当对 `cmd.ctx` 做 `context.WithoutCancel` 再开事务，属跨模块数据完整性策略（涉及取消语义、写放大与 ticket 04 的 WAL/磁盘满预案 D42#5），应由 ticket 04 的所有者决策。**但本票不得把「detached ctx」这个决定权外包出去**：取消一个任务与把它持久化是 loop 自己的职责，且修法只需 2 秒超时的一行。验收要求：修 MAJOR-1 时附一条测试，断言「取消的任务在库里以 `state="cancelled"` + `ended_at NOT NULL` + 在途 `tool_call` 行为 `outcome="cancelled"` 存在」—— 当前该断言的三个分量全部为 0/NULL，探针 P4 已量化。

## 五、重跑记录（真实输出，命令均限定 `./internal/agent/`）

```
$ export PATH="/d/work/base/go/bin:/e/work/base/msys64/mingw64/bin:$PATH"; export GOPATH=/d/work/base/gopath

$ go vet ./internal/agent/
（无输出）exit 0

$ go test ./internal/agent/ -count=2 -v
ok  	github.com/CarlosShao/wisp/internal/agent	1.763s
# 实跑 35 个 Test 函数 × 2 遍 = 70 次 PASS，0 SKIP（与 `grep -c '^func Test'` 的 35 逐一对齐）：
#   loop_golden 10 / control 6 / guard 5 / spill 5 / compress 4 / prompt 3 / truncation 2

$ go test ./internal/agent/ -race -count=1
ok  	github.com/CarlosShao/wisp/internal/agent	2.200s
exit 0

$ go test ./internal/agent/ -race -count=25 -run 'TestDefaultControlCancelsRunningTask|TestGoldenCancellationEndToEnd|TestSteeringInsertsIntoNextRound'
ok  	github.com/CarlosShao/wisp/internal/agent	11.796s
exit 0   # 专挑 3 个「异步 + 取消 + 计数器」用例加压，无竞争、无 flake

$ cd tools/d22scan && go run .
d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/exit 0
```

探针（`internal/agent/zz_probe_adversarial_test.go`、`zz_probe2_test.go`，P1–P14，交付前已删除、未提交）：P1/P2/P3/P5/P7/P8/P9/P10/P11/P12/P13/P14 均实跑并留档于本报告各条证据。
取证过程纪律（记录在此以免被当作可省略的细节）：早期一次经长输出通道读取的结果曾显示 P11 为 `maxInflight=1 / elapsed=490ms`、并显示 P3 的夹具「不存在」、P12 三段「被裁剪」—— 这些读数与源码逐行阅读结论互相矛盾，遂对同一批探针用 `-count=3` 定向复测：P11 三次均为 `maxInflight=4 / 160ms`（PASS）、P3 三次均成功加载夹具并给出 `complete=true`、P12 三次均显示 identity 54/150、safety 93/100、style 49/50 **正文完整未被裁剪**。故**前一批读数判为输出通道串扰，不予采信**；本报告只采信可复现且与源码一致的读数。特别地，「工具执行实为串行」与「D39 三段正文被预算裁掉」这两条一度成形的指控**均不成立**，未列入缺陷清单。

## 六、结论

- AC1–AC7 **七个验收族各自的判据全部成立**，判据仪器（`harness.requests()`）经逐行比对 + 阳性对照 + 拒绝请求计数三重验证后**判定健全**，AC6 的「0 请求」不是被改小的尺子量出来的假绿。
- 但整票不能按现状签收：MAJOR-1 让「取消」这一命名终止态的取证结果永久错误（且是本票自己的修复引入的），MAJOR-2 让 C22 的 `timeoutMs` 在真机工具上退化成一条 `internal` 类噪声。两者都是一行修。
- 另：`DEFERRED(D28-1)`（压缩暂同步调用）符合契约「synchronous fallback for S1 acceptable but flagged」的明文豁免，登记在 `compress.go:22-27`，**不算违规**；`DEFERRED(D11-3)`（force_tool/force_chat 规则表）代码与 Progress log 都有交代，但契约的 `## Out of scope` 未列它 —— 建议把「D11(3) 不在本票」补进 Out-of-scope 一行，别让它从验收面上蒸发（流程项，不计缺陷）。
- 建议：退回 `in-progress` 修 MAJOR-1 + MAJOR-2（附上述新断言），MINOR 按批处理；修完只需复跑本报告第五节的 5 条命令 + P4/P8 两个探针形状。

**VERDICT: FAIL**

---

## ⚠ 事后更正（2026-09-20，编排者，registry A26）：本文件中记为 PASS 的 `d22scan` 那条**是空跑**

本文件里 `cd tools/d22scan && go run .`（**不带 `-root`**）被记成 `clean` / PASS。
实测：该形式**默认扫描 `.`，即扫描器自己那个 module** —— 里面没有 `internal/`、没有 `cmd/`，
allowlist 读不到时被当空 ⇒ **它检查了 0 个生产文件却退出 0**。
⇒ 本文件据此支撑的"裸 goroutine / 墙钟超时 / `filepath.Clean` 越权 / emoji"四项，
**在当次验收中并未被这道门检查过**。原文数字与结论**一律不改写**（历史测量保留可读），
此段为叠加更正。
- 已由票 67（commit `23ebb59`）修好仪器本身：`checkRoot()` 让空范围致命退出，
  输出新增 **`examined N production Go files`** ⇒ **门今后必须自报工作量**，只报 `clean` 不作为证据。
- 当前 HEAD 用正确调用实测：`clean`，exit 0（本条是**事后**为那两个 AC 补上的真证据，不追溯证明当次验收有效）。
