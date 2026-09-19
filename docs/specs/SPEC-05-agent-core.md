# SPEC-05 · Agent 核心：循环、LLM 抽象、上下文与记忆

> 追溯：D8/D10/D11/D15/D20/D21/§14.2/§14.3/D39/D32/C5/C6/C7/C13/C22/C23/C31；工具面见 SPEC-07。

## 1. 背景与问题

Agent 循环参考 Pi（MIT，10.6 万 star）的骨架（D21）：**循环薄、难题外移**——流式与 tool-call
解析的真正工程量在 D8 的三协议 adapter（Pi 的 `StreamFn`），权限层（D4/D19/D30）是 Pi 完全没有、
必须自建的部分。本 spec 冻结循环结构、事件全集、上下文预算数值、失败语义与记忆机制。

## 2. 循环结构（每任务）

```
输入（语音转写 / 文字 / 面板） → 会话控制层（本地正则，D11①）
    ├ 命中控制词（停/取消/重说/大声点/确认）→ 直接执行控制语义，不经 LLM（几十毫秒级）
    └ 业务意图 → Agent 循环：
        组装上下文（§4） → LlmProvider 流式调用 → StreamEvent 消费 →
        ├ 文本 → 结果分流（D10，§7）
        └ tool_calls → 每调用：RiskAssessor(C19) 判级 → 门控（L0 直通/L1 阻止窗口/L2 审批队列）
                        → 执行（并发 ≤4，per-tool timeoutMs）→ 结果（超长 spill，D15③）→ 回灌循环
        stopReason=length → failToolCallsFromTruncatedMessage：该消息内所有未闭合 tool call 全部判失败
                            （D21 必偷设计，不得拿截断参数执行）
        终止：无 tool call 且文本输出完毕 / 用户取消 / C22 刹车 / token 预算 / 50 轮兜底
```

- force_tool / force_chat 规则表（TOML，D11③）：用户显式偏好覆盖，注入为指令而非前置分类。
- Steering（Pi 双层队列）：运行中插入的指令进内层队列；结构性重启（换指令重来）走外层。
- 取消：root ctx 贯穿 loop → stream → tool.execute（对应 Pi 的 AbortSignal 贯穿，D21#5）。

## 3. LLM 抽象（D8/C5/C6/C7）

### 3.1 三协议 adapter

OpenAI Chat Completions（`/v1/chat/completions`，含 tool_calls）· OpenAI Responses
（`/v1/responses`，reasoning item）· Anthropic Messages（`/v1/messages`，原生 tool use /
prompt caching / extended thinking）。**Agent 核心只见归一化表示**，协议差异止步于 adapter。

### 3.2 C6 StreamEvent（归一化事件全集，契约级）

```
TextDelta{text} · ReasoningDelta{text}（provider 有才发）
ToolCallStart{id,name} · ToolCallArgsDelta{id,partial} · ToolCallEnd{id}
Usage{in,out,cached}
Stop{reason: end_turn | max_tokens | tool_use | stop_sequence | content_filter | cancelled | error}
Error{class(D37 枚举), provider_code, retryable, retry_after}
Done
```

硬要求：`Stop.reason == max_tokens` 必须触发 `failToolCallsFromTruncatedMessage`。

### 3.3 C7 内部消息表示（必须含多模态部件，否则 `screen.capture` 是死的）

```go
type Content interface{ /* sealed */ }
type TextPart struct{ Text string }
type ImagePart struct{ MimeType string; BytesRef string; Alt string } // bytes_ref 指内存/临时文件，
                                       // 不落 DB、不进日志、不进诊断包（隐私）
type ToolUsePart struct{ ID, Name string; Input json.RawMessage }
type ToolResultPart struct{ ID string; Content []Content; IsError bool }
```

### 3.3a Provider 目录、角色与兜底链（2026-09-19 用户批准，SPEC-03 §3.1 的运行侧）

- **角色系统**：`roles.{chat, memory_extract, handoff, summarize}` 各自独立 model + temperature +
  thinking_intensity；记忆提取/摘要允许挂便宜小模型（省钱且隔离缓存影响）。
- **兜底链（取代单条 fallback_provider）**：`text_chain` 有序遍历——401/403/配额尽/5xx 重试耗尽
  → 跳下一家；**配额耗尽的跳转必须产生用户可见事件**（悬浮球 + 面板），不得静默降级。
- **voice 链**：`realtime → 云级联(C9) → 本地 sherpa 级联`；本地链零成本、永远最后（票 61/60）。
- **限流自约束**：per-provider `rpm/tpm` → 本地令牌桶，在 provider 429 之前先自守规矩；
  与工具并发 ≤4、任务调度协同。
- **兼容开关**：`compat.loose` 等（SPEC-03 §3.1）由 adapter 消费——缺失字段按开关容忍或报错。
- **探测（probe）**：tool_calls/vision/thinking/audio 逐项实测（票 11 实现）；结果写
  `provider_health`（SPEC-02 schema v2）。「声明 ✓ 实测 ✗」必须可见。
- `context_window` 挂在**模型条目**上（per-model，非 per-provider）——D15 全部阈值按它缩放。

### 3.4 重试与失败语义（§14.2，失败必须可见且可区分）

| 情况 | 行为 | 用户看到 |
|---|---|---|
| 429/5xx/网络抖动 | 指数退避重试 ≤3 次（重试在循环外，不污染核心） | 保持 `Thinking`，>5s 显示「重试中 (2/3)」 |
| 重试耗尽 | 明确报错，**不得静默降级** | `Error` + 面板给原因与建议（换 provider/查 Key/查额度） |
| 401/403 | 不重试 | `Unconfigured` + 引导配置面板 |
| 额度耗尽 | 不重试 | `Error` + 区分「网络问题」与「账户问题」 |
| 流式中途断开 | 保留已收部分并标注不完整；未闭合 tool call 全判失败 | 面板显示部分内容 + 「响应中断」 |

`text_chain`（有序）：`provider` 类错误重试耗尽后切下一家，链尽 → `Error(provider)`，
任务 ctx 保留、可从断点重试（D40#4）。
网络细节：尊重系统代理（WinHTTP 默认）+ 支持 `HTTP(S)_PROXY`；企业 TLS 拦截时错误必须区分
「证书不受信」与「连不上」（D42#4）。

## 4. 上下文预算（D15/D39，数值已定案）

### 4.1 系统提示词分段（顺序 = 缓存前缀排序，硬约束）

| 段 | 内容 | 预算 | 位置 |
|---|---|---|---|
| ① | 身份与边界（含「不写代码」） | 150 | **缓存前缀** |
| ⑤ | 安全与确认规则 | 100 | **缓存前缀** |
| ⑥ | 输出与播报风格 | 50 | **缓存前缀** |
| ② | 工具索引 · 常驻部分（内置全量 ~1200 tok） | 1200 | **缓存前缀** |
| ③ | L1 用户画像（≤20 条，~400 tok） | 400 | 后缀（偶尔变） |
| ② | 工具索引 · BM25 top-K（≤5） | ≤300 | 后缀（每轮变） |
| ④ | 当前时间/焦点窗口/场景 | ≤100 | **必须最后** |
| | 合计 | **≤2300**（D39 修正：原 2000 算不平） | |

- 排序是成本问题：前缀逐字节稳定才能命中 prompt caching；④ 与 BM25 段每轮都变，
  必须排在缓存断点之后。**C5 须把「缓存断点位置」作为 provider 能力暴露**
  （Anthropic 显式 `cache_control`，OpenAI 隐式）。
- L1 画像变更使缓存失效 → 画像提取**批量提交**，不逐条更新提示词（§14.3）。

### 4.2 工具注入两级（D15①②）+ 长输出 spill（③）

- 内置工具全量常驻注入；第三方工具本地 BM25/关键词检索 top-K ≤5，零 LLM 往返；
  检索未命中 → `list_tools` 元工具自查完整目录（循环内普通工具调用，非额外 LLM 往返）。
- 单工具结果 >4000 token（约 6000 汉字/16KB）→ 全文落 `artifacts\tool-output-<id>.txt`，
  上下文只留头 500 + 尾 200 token + 总长度 + 路径；模型可 `fs.read` 按需再读（宿主内部 spill，
  不经门控，D34 注②）。
- 原始输出硬上限 1MB（`shell.exec`/D46 command 的 stdout 同为 1MB）→ 截断标 `truncated=true`。
- **所有阈值必须按所配模型的 `context_window` 等比缩小，不得硬编码**（否则小上下文模型直接 400）。

### 4.3 会话历史压缩（D15④）

- 历史 >12000 token → 对最老轮次摘要压缩，保留最近 3 轮原文；仍超 → 继续压最老直到达标。
- 压缩在**响应路径外**（Warm 窗口内）做；**不得丢弃 tool_call 的 id 与结果引用**
  （否则截断失败判定与 C25 溯源断链）。

## 5. 会话与保活（B4/C31）

- 会话定义：唤起 → 用户显式结束，或空闲 90s（`[session] warm_timeout_sec`）。
- `Warm` 态：VAD/ASR/TTS/标点模型保持加载、**麦克风关闭**、面板可隐藏保活
  （`[panel] keep_alive_in_session`）；单击球/快捷键 → `Listening` **零模型加载**。
- `Conversation`（显式开启，默认 false）：麦克风持续开启、红色常亮环；30s 无语音 → `Warm`；
  首次开启 L2 级隐私确认并记日志。
- 省电模式：Warm 缩至 30s、Conversation 至 15s；关闭 KWS 常驻（D42#6）。
- 会话结束 → `Settling` → Dispose SessionScope（C11 全程，含 `FreeOSMemory`）→ 卸模型/
  隐藏面板/失效会话授权（D45-2 的授权随 session_id 过期）。
- 模型引用计数归 `session` 模块（唯一「保活」所有者）；语音侧串行加载策略见 SPEC-04 §4。

## 6. 记忆系统（D20/C13）

- **L1 画像**：`profile` 表 slot 有限枚举，≤20 行；**全量注入** system prompt（~400 token）；
  提取在任务结束后的 Warm 窗口**异步**跑（不占感知延迟）；反馈 = 悬浮球短闪「记住了：X」+
  一键撤销（记错比不记更伤人，纠错成本必须≈0）；满 20 条按 `updated_at` LRU 淘汰**并写日志**
  （淘汰必须可见）。禁同步提取（+300–1500ms）、禁 embedding、禁对话摘要堆积。
- **L2 显式记忆**：「记住…」→ `memory.save`（L1 风险级）；`memory.recall`（L0）关键词 + 最近
  优先；**不自动注入** system prompt。
- **L3 任务日志**：`task_log` 表，30 天；与运行日志分离（产品数据 vs 诊断数据）。
- 面板全量可见、逐条可删、可导出；存储位置写进文档。

## 7. 结果分流（D10，判定从上到下第一条命中，结构化标记优先于长度）

| # | 条件 | 通道 | 拉面板 |
|---|---|---|---|
| 1 | 含代码块/表格/≥3 项列表/≥2 链接/文件路径清单 | 落 `artifacts\` + 面板 + 播报摘要与路径 | ✓ |
| 2 | >400 汉字（或 >800 字符） | 落文件 + 面板 + 播报摘要（≤60 字）与路径 | ✓ |
| 3 | 61–400 汉字 | 系统通知 + 面板（播报首句 ≤60 字） | ✓ |
| 4 | ≤60 汉字且无结构化标记 | TTS 播报 + 悬浮球 + **写剪贴板** | ✗ |

- 摘要由 LLM **同一次响应**给出（系统提示词⑥段要求长结果附 ≤60 字口播摘要），不得为此再发请求。
- 落文件 = 宿主内部工件写入，不经门控；但必须面板可见可删可导出、路径写日志、计入 500MB 配额。
- 剪贴板只在第 4 档写（长结果覆盖用户剪贴板是负体验）。
- 第 3 档依赖 `notify`，Focus Assist 可能抑制 → **必须同时更新悬浮球角标**（D42#8）。

## 8. 防死循环与预算（C22 LoopGuard + C23 CostMeter）

- 梯度刹车：重复调用检测阈值 **[3,5,8]** 梯度注入提醒（学 DSH repeat-tool-reminder）；
  达 8 → `Stuck` 态并**明示重复了什么**，不得静默停。
- 每工具 `timeoutMs`（结构化 TOOL_TIMEOUT，协作式中止）；单任务 token 预算默认 200k；
  轮数上限 50 仅作高兜底。任一超限 → `Stuck`，明示「本次任务已用 X token / 约 ¥Y」。
- CostMeter（C23）：每任务 token/费用/工具次数/耗时；日/月累计（`cost_daily`）；
  预算达 80% 提示、100% 默认暂停新任务（`[cost] over_budget` 可配仅告警）；
  单任务 200k 预算（防死循环）与费用预算（钱）是两个层级，都要。
- 与 D31 的关系：轮数兜底、梯度提醒、路径冲突（C20）见 SPEC-06 §6。

## 9. 测试决策

- **黄金流回放**：录制真实 SSE 流存 `testdata/golden/`，回放驱动 Agent 循环——
  验 tool-call 解析、`stopReason=length` 全失败、取消、重试退避、fallback 切换（不是 mock 函数，
  是真实字节流）。mock-llm（SPEC-11）同时提供故障注入端点（429/5xx/断流/慢流）。
- 接缝：C5 `LlmProvider` 是 LLM 唯一接缝；C8 是音频唯一接缝；本 spec 的测试全部打到这两个接缝上，
  不为内部函数开测试后门。
- 上下文预算：构造超长工具输出/超长历史断言 spill 与压缩行为、tool_call id 保留、阈值随
  `context_window` 缩放。
- 记忆：提取异步性（断言首字延迟不受影响）、LRU 淘汰写日志、撤销行为、批量提交。
- 分流：D10 四档各构造样本（含「40 字表格」这种结构化优先用例）断言通道选择。
- 成本：用量归集、日/月聚合、预算暂停。

## 10. 不做什么

- 不做意图分类前置层（D11 已废）；不做多 Agent/子代理（草案 6.6）。
- 不做 MCP client（D13，接口位已留）。
- 不做对话摘要自动堆积/embedding 检索（D20/16.9#3）。
- 不为「未来本地模型」预留 provider 之外的东西（ollama 预设已覆盖本地需求）。
