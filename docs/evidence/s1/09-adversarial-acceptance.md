# T09 对抗验收报告（T09-adv 独立执行）

> 执行者：T09-adv 子代理（实现者为 T09-impl，相互独立）。
> 时间：2026-09-19T23:00Z 前后（本地 2026-09-19 深夜）。
> 环境：go1.27.1 windows/amd64，GOPROXY=goproxy.cn，Windows 11 x64，Git Bash。
> 并行保护：T14-impl 在同工作树开发 internal/models、internal/buildinfo、compose——本验收全程
> 未触碰；验收期间工作树出现 `M go.mod`/`M go.sum`（golang.org/x/sys 0.47.0→0.48.0、
> 新增 indirect golang.org/x/crypto 0.57.0），判定为 T14 在途改动，本验收未回滚、未提交；
> 该依赖 bump 与 internal/llm（stdlib+observe/config）无关，不影响本报告全部测试结论。
> 本验收 git 操作均使用显式路径。

## 1. 全量复跑（gates）

| 命令 | 结果 | 证据 |
|---|---|---|
| `go vet ./...` | PASS | VET_EXIT=0（T13 报告曾记录 adapter_test.go:190/238 self-assignment 在途 WIP 诊断，提交版已消除） |
| `go test ./... -count=1` | PASS | TEST_EXIT=0；llm 1.364s / llm/golden 0.084s / llm/openaichat 20.492s / memory 12.999s / audio 16.686s（T13）/ llmrecord 0.199s / 其余全 ok |
| `CGO_ENABLED=0 go test ./internal/llm/...` | PASS | ok x3（llm / golden / openaichat） |
| `CGO_ENABLED=0 go test ./cmd/llmrecord/...` | PASS | ok |
| `go test -race -count=1 ./internal/llm/...` | PASS | ok x3，无 race 报告（2.4s/1.1s/19.2s） |
| mockllm 模块（cd tools/mockllm）`go vet ./...` + `go test ./... -count=1` | PASS | MOCK_VET=0，MOCK_TEST=0，8 个测试函数与交接声称一致 |

注：清单原文命令 `CGO_ENABLED=0 go test ./internal/llm/... ./tools/mockllm/...` 在仓库根执行时
对 tools/mockllm 报 `directory prefix tools\mockllm does not contain main module`（嵌套独立
module，父模块不可见）——与交接日志"nested module excluded from parent ./. . ."一致，须在
tools/mockllm 目录内单独执行，非缺陷。

## 2. C6 契约对账（events.go vs SPEC-05 §3.2）— PASS

- 事件全集恰好 9 种：`EvTextDelta..EvDone`（internal/llm/events.go:31-41），与 SPEC-05 §3.2
  逐项一致；`StreamEvent` 结构体 union（注释明确契约级、agent 核心不见协议形状）。
- Stop 7 值 + `Valid()`：end_turn / max_tokens / tool_use / stop_sequence / content_filter /
  cancelled / error（events.go:71-88）；`TestStopReasonEnumIsExact` 钉死。
- Error 字段：class=D37（observe.ErrorClass）、provider_code、retry_after（Duration 字段）、
  retryable 以 `Error.Retryable()` 类策略方法 + no-retry sentinel 表达
  （internal/llm/errors.go:173、internal/observe/errors.go:214；TestRetryNoRetrySentinel、
  TestRetryNeverForAuthBudgetCancelled）。C6 四字段语义齐备，形态为方法而非布尔字段，判等价。
- Done 恰一次：单条 Stream 恰一对终态事件——成功 `Stop→Done`（adapter.go:230-233）、失败
  `Error→Stop{error}→Done`（fail，:273-284）、取消 `Stop{cancelled}→Done`（cancelled，
  :288-294），各路径单次发射后立即返回；TurnCollector 对 Done 后事件幂等忽略（events.go:216）。
- ctx 取消 → `Stop{cancelled}` 无 Error：cancelled() 路径结构性不发射 Error 且返回 nil；
  `TestGoldenCancellationMidStream`（adapter_test.go:389-427）独立复跑通过：首个 TextDelta 后
  cancel → 事件流含 Stop{cancelled}、末事件 done、Stream 返回 nil、流未跑完（<32 事件）。
- SPEC 硬要求 max_tokens→failToolCallsFromTruncatedMessage：seam 保证 reason 原样上抛、
  注释钉死票 10 义务（events.go:24-26、adapter.go:244），责任划分正确。

## 3. golden 字节一致独立复现 — PASS（测试强度不足，见 MAJOR-1）

`go test -count=1 -v -run TestMockllmGoldenByteIdenticalEvents ./internal/llm/openaichat/`
独立复跑：4/4 子测试 PASS（tool-call 0.01s / max-tokens 0.00s / usage-multichunk 0.00s /
long-text 1.34s）。同一 fixture 经两条跑器（internal/llm/golden httptest 回放器 vs mockllm
真子进程）喂同一 adapter。

**MAJOR-1**：比较函数 `sameEvent`（mockllm_integ_test.go:235-246）对 `EvTextDelta/
EvReasoningDelta` **无条件返回 true**——文本内容完全未比对；其注释声称 "assembled text
compared below"，但循环体内不存在任何累计文本比对；同文件 `rawGolden`（:154，mockllm 原始
字节保真探针）是**死代码**，从未被调用。后果：4 个 fixture 中 3 个（max-tokens、
usage-multichunk、long-text）若任一跑器的 golden 解析器在文本行上出错，现有测试无法失败
（事件计数、Usage/Stop/ToolCall 全等仍满足）；仅 tool-call 被 mockllm 模块测试字节级钉死
（TestChatGoldenReplayByteIdentity:287-308，served body == fixture 原始字节）。mockllm 侧
解析镜像测试（TestGoldenFormatParserMirrorsFixture）只对 backoff-429 做结构断言，不覆盖文本。
验收条 2 判 **PARTIAL**：属性大概率成立（两跑器对 LF 钉死 fixture 均逐字节供应，同 adapter
应产出全同事件流——回放器 parse 为 join 行 + \n 的规范化，对 LF 文件等价于原字节），但现有
证据链不成立。

**最小修复**（测试侧，无生产代码变更）：在 TestMockllmGoldenByteIdenticalEvents 中累计两路
TextDelta/ReasoningDelta 文本并在流末比对（稳健）；或直接全 JSON 逐事件比对（两跑器供应
字节应全同，TCP 分块不影响 SSE 行扫描，sameEvent 的"分块差异"前提不成立）；同时删除或真正
启用 rawGolden，并修正失实注释。

## 4. 故障注入独立复现 — PASS

对照真 mockllm 子进程 + /__control 端点，独立复跑 6 项全 PASS：

| 场景 | 断言（测试内） | 结果 |
|---|---|---|
| fail_next{429, times:2, retry_after:"1"} + Retry Max:3 | elapsed ≥ 1900ms（两次 1s 退避），最终 StopEndTurn + 非空文本 | PASS（3.20s） |
| fail_next{500, times:10} + Max:3 | err 分类=provider；/__control/state routes["chat"]==4（恰 4 请求） | PASS（1.41s） |
| fail_next{401, times:10} | err 分类=auth（Unconfigured 语义）；routes["chat"]==1（恰 1 请求，不重试） | PASS（1.36s） |
| latency{ms:120} | elapsed ≥ 300ms（3+ chunk × 120ms） | PASS（1.61s） |
| truncate{chunks:2} | ProviderCode=stream_disconnected + TurnResult.Incomplete=true | PASS（1.21s） |
| 429 单元级（TestRetryRateLimitHonorsHintVerbatim / TestGoldenBackoffLadder429） | 退避睡眠与 retry-after 逐字相等；录制的 429 经真实 HTTP 3 请求过梯 | PASS |

MINOR-2：429 集成断言阈值为 1900ms（2×1s 减调度容差），交接声称"实测≥2s"；数值口径不一致，
行为本身（retry-after 逐字遵守）成立。

## 5. schema v2 审计 — PASS（附 MINOR-1）

- **逐字节 diff（程序化）**：临时目录工具抽取 SPEC-02 §3 的 2 个 ```sql 块与 schema.go 的
  ddlV1/ddlV2 反引号原文比较——剥离围栏换行（提取产物，恰 1 字节）后 **ddlV1 ≡ block[0]、
  ddlV2 ≡ provider_health block[1]，逐字节一致**（含全部中文注释）。
- **provider_health 列**：provider/model 联合 PK、probe_json DEFAULT '{}'、last_probe_at、
  last_error、last_error_at、latency_ms_p50 DEFAULT 0、quota_state DEFAULT 'ok'（4 值枚举注释）
  ——与 SPEC-02 §3 v2 段文字描述（probe 六能力位 JSON、运行观测态 vs config 真相源分界）一致。
- **契约内省测试真实**：TestSchemaContractIntrospection 断言恰 9 表（contractTableOrder 含
  provider_health）、逐表 PRAGMA table_info 的 name/type/notnull/pk、schema_version='2'、
  ddlV1=15 语句、ddlV2=1 语句经 SQLite 真实执行往返。
- **迁移链**：0→1→2 生产步在 fresh open 时执行，备份 bak-0-1 + bak-1-2 恰各一份
  （TestMigrateFreshDatabaseCreatesCurrentVersion:348）；单事务原子性/更新版本拒绝/异库拒绝
  均有测试；TestMigrationChainWithBackup 以 test-only 2→3 步验证备份、水位、数据存活机制。
- **MINOR-1**：无独立"真实 v1 库（schema_version='1'、含种子数据、无 provider_health）重开
  迁移到 v2"场景测试——1→2 步仅在 fresh 0→1→2 链中对空 v1 库执行，2→3 机制测试用的是
  test-only 步骤。建议补一测：`Open(withSchemaTarget(1))` + 种子行 → 关闭 → 重开（目标 2）
  → 断言 provider_health 存在、数据存活、bak-1-2 落盘。鉴于 DDL 仅追加表、相邻步骤机制已测、
  DDL 与 SPEC 逐字节一致，判非阻塞。

## 6. 链 / 兼容 / 令牌桶 — PASS

- text_chain 步进：TestChainFailoverOrderAndCallback——401(auth)→下一家、配额(budget)→
  下一家、最终元素成功输出；OnFailover 回调携带 From/To/Cause 且恰 2 次。TestChainExhausted
  返回末元素分类错误；TestChainNoFailoverAfterPartialPayload——部分载荷已交付后链停止
  （已登记偏差：避免重复半截回复，合理）。
- SPEC-05 §3.3a"配额耗尽跳转必须产生用户可见事件"：seam 层以 OnFailover 回调满足（C6 事件集
  契约精确、不造假事件），悬浮球/面板接线归票 10/33——已登记偏差，判可接受。
- compat.loose / allow_missing_usage：missing-usage、missing-finish、no-[DONE]、missing tool
  index 四组 strict-vs-loose golden 测试齐备且通过。
- rpm/tpm 令牌桶**已实现**（无需核对票 11 登记）：TestBucketLimiterPacesRPM、
  TestBucketLimiterTightTPMWaitsAndCtxCancels、TestLimiterProviderReconcilesUsage 全 PASS。
- 预设 12 家 provider 钉死测试（catalog_test.go:139 起）+ /v1/models 自动发现导入
  capabilities-unknown（TestMockllmDiscoverModels 对真 mockllm）+ probe 原语契约
  （fc/vision/thinking 最小请求+检查；audio 诚实 not-implementable，归 15/60/61）——实现
  归票 11 已登记，符合票据范围。

## 7. 越界扫描 — PASS（附裁定）

5 个 T09 提交逐个 `git show --stat`：

| commit | 触碰 | 裁定 |
|---|---|---|
| fe87e12 | internal/llm（13 文件）+ 票 09 | 在范围内 |
| f27928e | internal/llm + internal/llm/openaichat + 票 09 | 在范围内 |
| 0839407 | tools/mockllm（11 文件 + 自有 go.mod + 2 fixtures）+ 票 09 | 在范围内 |
| 26f31b0 | .gitattributes（1 行，*.sse LF 钉死） | **裁定合理**：字节回放的跨平台前提，独立 chore 提交且已披露 |
| 3b7bf24 | internal/llm、internal/memory（schema v2 + DAO）、tools/mockllm/goldenfmt.go、cmd/llmrecord、docs/specs/SPEC-02、票 09 | **裁定合理**：schema v2 属票 09 原语范围（票据明文 "results → provider_health, SPEC-02 v2"）；SPEC-02 sql 块系逐字节契约的另一半，必须同步；cmd/llmrecord 为 golden 写格式消费端（票据"golden-SSE test harness"范围 + DEFER(a) 载体），自包含于 cmd/，进度日志与交接均已披露 |

- internal/audio、internal/models、internal/buildinfo、docker/、compose：**零卷入**（5 提交
  stat 逐个核对，无一路径）。
- go.mod/go.sum：5 提交均未触碰（mockllm 为 stdlib-only 独立 module）；工作树当前 M 状态为
  T14 在途改动（见头部说明）。
- 零 emoji：5 提交全量 diff 经 Unicode 范围扫描（U+1F300-1FAFF / U+2600-27BF / FE0F）无命中。
- 无密钥：internal/llm、tools/mockllm、cmd/llmrecord 代码与全部 golden fixture 扫描
  `sk-`/`api_key`/`Bearer` 高熵串——无真实 Key；catalog_test.go 仅 `sk-test` 假值。
- tools/mockllm/mockllm.exe 未入库（.gitignore *.exe 命中），git 跟踪清单干净。

## 8. 票据对照

验收 6 条（.scratch/wisp/issues/09-llm-provider-openai-mockllm.md:48-58）：

1. golden 回放测试（tool-call 组装 / max_tokens / mid-stream disconnect / cancellation /
   backoff ladder，全经录制 SSE 字节、无函数 mock）— **PASS**（adapter_test.go TestGolden*
   全经 golden 回放器真实 HTTP；断言非常量：End 顺序、args 分片按 index、usage{21,9,8}、
   disconnect 保留部分载荷 + incomplete、cancel 序、3 请求梯）。
2. mockllm golden + 单测回放器字节一致事件流 — **PARTIAL**（测试存在且通过，但文本轴未
   比对，见 MAJOR-1）。
3. 控制端点故障注入（429 retry-after / 5xx×3→provider / 401 不重试 / latency）— **PASS**
   （第 4 节独立复现，请求计数由服务端 state 钉死）。
4. 代理测试（HTTP(S)_PROXY；untrusted-cert ≠ unreachable）— **PASS**
   （TestHTTProxyEnvHonored 经录制代理 + TestTLSUntrustedCertDistinctFromUnreachable：
   tls_untrusted_cert vs connect_failed 两码不同，D42#4）。
5. Usage 多 chunk 聚合 — **PASS**（max-merge 设计对 Anthropic 累计式安全；
   TestGoldenUsageAggregation + usage-multichunk fixture + TestUsageMaxMergeAggregation）。
6. 12 家预设 + 发现导入 unknown 能力 + probe 原语契约 — **PASS**（第 6 节）。

DEFER 4 项登记核对：

- (a) 真录制归 H2：票 09 日志 ✓ + docs/reports/pending-and-issues.md:10 ✓（[H2] LLM API Key
  阻塞票 09 黄金 SSE 录制；llmrecord 就绪待 Key）。
- (b) probe 实现归票 11：票 09 日志 ✓ + 票 11:32 ✓（Capability probe implementation
  supplement，实测写入 provider_health）。
- (c) audio probe 请求定义归 15/60/61：票 09 日志 ✓ + 票 61:37-38 ✓（audio_in/audio_out
  probe 进 provider_health，票 11 框架）。
- (d) compose 接线（18080）：票 09 日志 ✓；**MINOR-4**：issues/ 中未钉到具体目标票
  （grep mockllm 于 14/16/12 号票无 compose 条目；mockllm 默认端口已是 18080，
  main.go:43，票 12 的 S1 门用 go run 不受影响）。

Progress log：5 条 append-only（claim → seam-core → openaichat → mockllm → 测试+收尾 →
handoff），完整、新者在后 ✓。Status: in-progress（未自行标 DONE，交 orchestrator 裁定）✓。

## 问题汇总

- **BLOCKER**：无
- **MAJOR**：
  1. TestMockllmGoldenByteIdenticalEvents 文本轴零比对（sameEvent 对 TextDelta 恒真、注释
     失实、rawGolden 死代码）——验收条 2 的证据链不完整，3/4 fixture 的双跑器文本分歧不可测。
- **MINOR**：
  1. 无独立"真实 v1 库（含数据）→v2"迁移场景测试（1→2 仅在 fresh 链中对空 v1 库执行）。
  2. 429 集成断言阈值 1900ms，与交接"实测≥2s"口径不一致（行为本身正确）。
  3. 取消测试未显式断言事件流中无 EvError（cancelled() 代码路径结构性保证，风险低）。
  4. DEFER(d) compose 接线未在 issues/ 钉到具体目标票。

## 裁决

12 项 checklist 中 11 项以真实行为断言独立复跑通过；全部 gates（vet / 全量 test / CGO0 /
race / mockllm 模块）绿；schema DDL 与 SPEC-02 §3 经程序化 diff 逐字节一致；越界、emoji、
密钥全部干净；票据与 DEFER 登记基本完备。唯一 MAJOR 是验收条 2 的测试强度：字节一致声称的
文本轴完全未被比对（注释与实现不符），属测试侧 5-10 行最小修复，无生产代码缺陷证据，但不修
则"两跑器字节一致"这一票 09 复用基石的证据不成立。

**VERDICT: FAIL (最小修复集)**
最小修复集（全部测试侧、归 T09-impl）：修复 MAJOR-1——TestMockllmGoldenByteIdenticalEvents
累计比对两路 TextDelta/ReasoningDelta 文本（或全 JSON 逐事件比对），删除或启用 rawGolden，
修正 sameEvent 失实注释。MINOR-1/2/4 建议顺手一并处理，不单独阻塞。
