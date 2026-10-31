# 票 386 落地腿 386-r2 · 实现件证据（mockllm 同类帽逐帽判定＋修）

- 腿：`386-r2`（写腿，386-r1 静默死亡重派）。日期 2026-10-31。开工锚 HEAD＝`46a80b579`（分支 `dev`，即派单锚本体）。
- 写面：`tools/mockllm/**`（实改 `chat.go`／`responses.go`／`messages.go`＋新增 `caps_multibyte_386r2_test.go`）＋本证据件。未动 `cmd/wisp`、`internal/**`、票面、台账、testdata、specs、frontend、push。
- 环境：`GOFLAGS=`（空）、gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`（v0.12.0, go1.27.1）、突变件在仓外 `D:\tmp\wisp\20261031\386r2\**` 经 `-overlay` 注入、提交 `-F` 在 `--` 前、主题逐字 `票 386 落地腿 386-r2`。

## 0. 开工三核

1. HEAD＝`46a80b579` 在 `dev`（`git log --oneline -1`＝编排者 A1332 重派提交，≥ 派单锚）。
2. 索引无本腿相关 staged 项（共享树他人脏件全在 `.scratch/**`、`design/**`、`docs/evidence/**`、`logs/**`——与写面零交）。
3. `git status --porcelain -- tools/mockllm` 空（写面无他人脏；目录里的 `mockllm.exe` 是 `.gitignore:2 *.exe` 在册的本地构建产物，未动）。

## 1. 票面 5 格逐条（原文＋判定＋读数）

| 格 | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| AC#1 | 逐帽判定：每帽"有意模拟 vs 无意缺陷"具名（对照真实 chat API 的 max_tokens/reasoning 行为；有意模拟的帽 ⛔ 动、如实登记） | **成立（四帽全具名，两顶判有意、五处切点判无意，§2）** | reasoning/reason 200 帽×3＝无意缺陷已摘；max_tokens 帽＝**帽本体有意模拟（⛔ 未动）＋字节切点无意缺陷（rune 化）** |
| AC#2 | 判"无意"的帽按 rune/去帽修＋多字节案一发（对照 385-r1 的 402 汉字案形） | **成立** | 新件 `caps_multibyte_386r2_test.go`：402 汉字 [think] 样本×4 方言口（chat 流式/非流式、messages、responses）＋max_tokens 汉字切两案＋ASCII 对照案，§3 |
| AC#3 | 反向判据（三读数＋还原核） | **成立** | 修前红 7 案→修后绿 7 案→`-overlay` 装回原始帽复红同位 7 案→撤 overlay 复绿；sha256 在案，§4 |
| AC#4 | 零回归（mockllm 全包＋消费方 E2E） | **成立** | 模块内全包 `28/18/0/0` rc=0；`internal/llm/...` 七包 rc=0；`internal/present`（票 30 E2E）rc=0；`cmd/wisp` 定向名册（票 12 E2E）`46/36/0/0` rc=0＋ComposedGate 单发 PASS；三 build＋vet rc=0，§5 |
| AC#5 | 纪律格：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`（只 `tools/mockllm/**`） | **成立（自证）** | 提交 pathspec 显式逐枚恰 5 件（4 产码/测试＋本证据件）；零 push；零翻框 |

## 2. AC#1 逐帽判定（本票核心交付）

判据声明：本腿开工时尝试抓取 DeepSeek/Anthropic 官方文档页作对照引用，两站 DNS 解析被本机网络环境拒（`api-docs.deepseek.com`／`docs.anthropic.com` 均解析到非公网地址）；对照改以三源立据——①仓内 385-r1 先例（`docs/evidence/s1/385r1-impl.md` §2/§7，其判定理由同样建在"对照真实 provider 行为"上）、②四个帽自身的代码内证（注释、启发式、协议语义）、③公开 API 语义通识（chat completions `max_tokens` 命顶即截断并报 `finish_reason="length"`；reasoning 类输出无任何 200 字节帽）。另据公开检索佐证 `finish_reason=length` 的截断语义（vLLM 侧同语义修复 [PR #58286](https://app.semanticdiff.com/gh/vllm-project/vllm/pull/58286/overview#1)）。

### 帽 A／B／C：reasoning/reason 200 字节帽三方言（`chat.go:168-169`／`responses.go:100-101`／`messages.go:114-115`）

**判定＝无意缺陷（非有意模拟），修法＝去帽（与 385 同法）。** 理由五条：

1. **不模拟任何真实 provider 行为**：真实 API 对 reasoning 类输出（chat 方言 `reasoning_content`、Responses 方言 reasoning summary item、Messages 方言 thinking block）**没有 200 字节响应帽**；请求侧的预算参数（如 Anthropic thinking 的 budget_tokens）是另一回事，不是响应截断。200 这个数字在任何方言的真实协议里都找不到对应物。
2. **出处与 385 已判死的回声帽同源同型**：385-r1 §2 已考 `git blame`——三方言的 `answer[:400]`＋`reasoning[:200]` 是建仓提交 `0839407bf` 里**同一笔的对称「合成封顶」惯例**；385-r1 已判该惯例"刻意但不模拟真实行为 ⇒ 缺陷"并摘掉回声半边，reasoning 半边是同一判定的未竟余部。
3. **与真护栏重复**：请求体已有 8MB 上界（`io.LimitReader(r.Body, 8<<20)`），reasoning 串＝`"thinking about: " + 请求回声`，长度物理有界。
4. **无任何消费方测量该帽**：包内测试全部短样本（逐文件核过）；golden 模式整条绕过合成路径；根模块侧 grep 无任何钉死 reasoning 长度的断言。
5. **两路口径一致性**：thinking 探针的确定性层（进程内假体）与 live mockllm 是同一测量的两口（`internal/llm/probe_thinking_test.go` 的双层形状）；假体给全文、mock 偷偷截 200，长样本下两口必分叉——这正是 385-r1 §2 理由 3 的同款陷阱。

**协议形状零改动**：SSE 事件结构、字段、finish/usage 语义全未动；变的只是 reasoning 内容可以有多长。

### 帽 D：max_tokens 字节切（`chat.go:207` `answer[:MaxTokens*4]`）

**判定＝两段拆开：帽本体＝有意模拟（⛔ 一字未动）；切点＝无意缺陷（rune 化修复）。**

- **帽本体有意模拟的具名证据**：①注释自证 `// max_tokens truncation: ~4 chars per token heuristic.`——启发式是被写明的设计选择；②语义对齐真实 chat API：`max_tokens` 命顶 ⇒ 截断＋`finish_reason="length"`，mock 的 `finish = "length"` 正是这个语义（票面自己标注的"带 finish=length 模拟语义"）；③385-r1 §7.2 当时就是因为它"带 finish=length 的模拟语义，动它是改模拟行为"而登记不修——先例已经把"帽本体＝有意模拟"钉过一次。
- **切点无意缺陷的具名证据**：真实 provider 在 token 边界截断，输出**永远是合法 UTF-8**（JSON 序列化层面也不可能携带坏序列）；mock 按字节切（`answer[:MaxTokens*4]`），多字节文本下拦腰切进 UTF-8 序列中间，JSON 编码器把坏尾巴变成 U+FFFD——与 385 已判死的回声字节帽同一缺陷类。读数①的额外反证：字节切出来的答案 402 字节**超过了它自己声明的 400 字节启发式界**（坏字节经 U+FFFD 展开后反而变长），连帽自己的语义都被切点破坏。
- **修法**：切点改 `runeBoundaryWithin(MaxTokens*4, answer)`——≤n 字节内最大的 rune 边界；**阈值启发式（len/4 > MaxTokens）与 finish=length 逐字未动**，ASCII 文本下行为与旧码逐字节相同（对照案钉死）。

### 未定义即停两条件的核对

①"所有帽都判有意模拟"——**不成立**（五处切点判无意），不触发登记即结；②"修帽必须动协议形状"——**不成立**（三处去帽只改内容长度、一处切点只改切进哪一字节，SSE/字段/finish/usage 全未动），不停手。

## 3. AC#2 多字节回归一发（402 汉字形）

新件 `tools/mockllm/caps_multibyte_386r2_test.go`，两枚测试七子案，样本沿用 385 的 `ticket30Sample()`（`测`×402＋摘要行）：

- `Test386ReasoningDeliversLongMultibyteInFull`（want＝`"thinking about: [think] " + 样本`＝1261 字节/437 rune；旧 200 帽切在第 201 字节＝第 59 个汉字序列中间）：`chat_stream_sse`（reasoning_content 逐 delta 装配）／`chat_nonstream_json`／`messages_thinking_block`／`responses_reasoning_summary` 四口全量回传＋`utf8.ValidString`。
- `Test386MaxTokensTruncationIsRuneSafe`：300 汉字＋`max_tokens=100`（全 echo 906 字节、启发式 226 tokens>100）⇒ 断言 `finish_reason=length`（有意模拟健在）＋内容＝rune 边界前缀 `"echo: " + 测×131`（399 字节）＋`utf8.ValidString`＋≤400 字节界；`hanzi_cut_lands_on_rune_boundary`（非流式）与 `stream_finish_reason_length`（流式）两口；`ascii_cut_byte_identical_to_cap` 对照案钉"帽仍按 `MaxTokens*4` 字节截断 ASCII"——**判有意的那半边没被顺手拆掉**。
- ⚠ **工具链陷阱存档（本腿开工实撞）**：测试文件初名 `multibyte_caps_386_test.go` 时 `go test` **静默不见**（`-list` 不列、`-run` 报 no tests to run，vet/build 亦不报错）——Go 文件名规则先剥 `_test` 再做 `_GOARCH` 后缀匹配，`386` 是 GOARCH 名（i386），文件被当 386 架构专属件在 amd64 上静默排除；`385` 不是架构名所以 385-r1 没踩。改名 `caps_multibyte_386r2_test.go`（尾段 `_386r2` 非架构名）即现。判据可复跑：`go list -f '{{.TestGoFiles}}' .` 对比改名前后。**后续任何腿起 `_386_`／`_amd64_`／`_windows_` 形文件名的测试件都会撞同坑。**

## 4. AC#3 反向判据（三读数＋还原核）

突变件＝修前原始三文件快照（`D:\tmp\wisp\20261031\386r2\orig\`），经 `go test -overlay` 注入，工作树零接触：

1. **读数①修前（帽在）**：`go test -run Test386 -v` **rc=1**，红句逐字（`caps_multibyte_386r2_test.go`）：reasoning 四案 `= 204 bytes / 84 runes, want 1261 bytes / 437 runes`，尾部可见 U+FFFD（200 字节切进的坏序列经 JSON 化的残迹）；max_tokens 两案 `content = 402 bytes / 138 runes, want 399 bytes / 137 runes`＋`402 bytes, exceeds the max_tokens*4 heuristic bound 400`；ASCII 对照案 **PASS**。
2. **读数②修后**：同命令 rc=0，七子案全 PASS；模块内全包 rc=0。
3. **还原核（overlay 装回原始帽）**：`go test -overlay <orig.json> -run Test386` **rc=1**，红点与读数①**同位同构**（reasoning 四案＋max_tokens 两案红、ASCII 对照仍绿）——测试咬的是帽，不是别处。
4. **撤 overlay**：rc=0 复绿；三产码文件 sha256（overlay 轮后现值）`chat 647B28BC4222104B`／`responses 94E17C32CDDC5A88`／`messages 4C03E0701A7EAB0D`（overlay 机制上只读注入，工作树零接触，此为轮后现值在案）。

## 5. AC#4 零回归（消费方 E2E＋build）

| 消费方 | 命令 | rc | 读数 |
|---|---|---|---|
| mockllm 包内（含新回归） | `go test -count=1 ./...`（tools/mockllm 模块内）＋`-v` 四数 | 0 | 四数 `RUN=28 / PASS=18 / FAIL=0 / SKIP=0`（RUN 含 10 子测试）；全包 0.16s |
| chat 方言族＋票 30 E2E | `go test -count=1 ./internal/llm/... ./internal/present/`（根模块） | 0 | `internal/llm 57.0s`／`adaptertest 4.0s`／`anthropic 20.6s`／`golden 0.1s`／`openaichat 29.0s`／`openairesponses 27.2s`／`present 4.1s` 全 `ok`；present 即票 30 真管线（真 mockllm 子进程） |
| cmd/wisp 消费方（票 12 E2E 等） | `go test -count=1 ./cmd/wisp/ -run '^TestRunTextTask|^TestHostDispatch|^TestMissingBlob|^TestTicket101|^Test365|^TestTaskOutputPointerNotice174r8|^TestProviders|^TestReplyListener|^TestPanelRouteRefuses|^TestUnansweredL2|^TestL1Veto|^Test375|^TestRunPacketCarriesTheSubagent|^TestSubagentStreamKey|^TestRunPacketReportsTheStreamLog|^Test197Subagent|^Test197NoAllowDoor' -v` | 0 | 四数 `RUN=46 / PASS=36 / FAIL=0 / SKIP=0`（与 385-r1 读数逐位一致）；含票 12 `TestRunTextTaskTextPathEndToEnd`；231.1s |
| cmd/wisp 慢枚单发 | `-run '^TestComposedGateBlocksAWriteForTwoSeconds$'` | 0 | PASS（3.66s） |
| 编译 | 根模块 `go build ./...`；tools/mockllm 模块内 `go build ./...`＋`go vet ./...` | 0/0/0 | — |

注记（如实报）：
- cmd/wisp 读数按 `scripts/wisp-cli-tests.sh` 台型前置 `third_party\sherpa-onnx` 上 PATH 后取（票 98/111 在册的 DLL 形状）；cmd/wisp 当前共享树零 WIP（本腿开工时唯一 tracked 脏件在 `.scratch/`）。
- 根模块 `go test ./cmd/wisp` 重包同一时刻只跑一发的规矩遵守：本腿全程单包串行。
- 三方言的 messages/responses 两路**解析了 `max_tokens`／`max_output_tokens` 但合成路径从不应用**（`responses.go:67`、`messages.go:179` 字段收下即弃）——票面四帽之外的独立观察，本腿⛔ 越票顺手补，登记 §7。

## 6. 四数／两尺／d22scan（tools/ 射程如实报）

- **四数**（口径：`^=== RUN` 含子测试／`^--- PASS` 只顶层／FAIL／SKIP，⛔ 相减）：mockllm 包内 `28/18/0/0`；cmd/wisp 定向 `46/36/0/0`；d22scan `clean`。
- **两尺**：`gofumpt -l tools\mockllm` **0 行**（rc=0）、`gofmt -l tools\mockllm` **0 行**（rc=0）。
- **d22scan 真仓扫**：`cd tools/d22scan && go run . -root ../../` **rc=0 clean**；live scope：`bans#1-5 internal/=287`、`bans#1-5 cmd/=52`、`ban#6 frontend/=116`、`ban#7 internal/tools/=25`、`ban#8 design/=39`、`ban#8 frontend/=116`、`ban#8 internal/=729`、`ban#8 cmd/=180`。顶层 `tools/**` 在八枚 ban 射程外（同 385-r1 §6 考）；新测试件 CJK 字符（`测`）不在任何 emoji 码位带。

## 7. 登记不修（残余，具名归 tools/mockllm 所有者后续开票）

1. **messages/responses 两方言的 max_tokens 死字段**——`responses.go:67`（`max_output_tokens`）与 `messages.go:179`（`max_tokens`）解析后合成路径从不应用；chat 方言有截断模拟而这两方言没有 ⇒ 两方言的 max_tokens 语义与真实 provider 不齐。补齐属"新增模拟行为"，超出本票"判+修已登记四帽"的授权，登记不动。
2. **cmd/wisp/subagent_carrier_197_test.go:605-608 的过期论证注释**——385-r1 §7.3 已登记（echo 帽摘除后其"不可达"前提失效，断言实测仍绿）；本腿动的是 reasoning 帽与 max_tokens 切点，对该注释的前提无新增破坏；归 cmd/wisp 所有者，本腿⛔ 越面。

## 8. 交付形状

- 产码：`tools/mockllm/chat.go`（摘 reasoning 帽＋max_tokens rune 边界切点＋`runeBoundaryWithin` helper＋注释）、`tools/mockllm/responses.go`（摘 reason 帽＋注释）、`tools/mockllm/messages.go`（摘 reason 帽＋注释）。
- 回归：`tools/mockllm/caps_multibyte_386r2_test.go`（新增，两测试七子案）。
- 证据：本件。临时件（原始帽快照／overlay JSON／全部原始 log）在 `D:\tmp\wisp\20261031\386r2\**`（仓外，只建不删）。
- 提交：主题逐字 `票 386 落地腿 386-r2`，`-F` 在 `--` 前，pathspec 显式逐枚；只 commit 不 push。
