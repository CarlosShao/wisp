# 票 391 落地腿 391-r1 · 实现件证据（mockllm 两方言 max_tokens 死字段补齐）

- 腿：`391-r1`（写腿）。日期 2026-10-31。开工锚 HEAD＝`1b9721519`（分支 `dev`，即派单锚本体；提交时 HEAD 已被编排者推进至 `e2ce24d2b`——收 387-r2 的台账件，与本腿写面零交）。
- 写面：`tools/mockllm/**`（实改 `messages.go`／`responses.go`＋新增 `caps_maxtokens_391r1_test.go`）＋本证据件。未动 `chat.go`（386-r2 刚修的对照件，零接触）、`cmd/wisp`、`internal/**`、票面、台账、testdata、specs、frontend、push。

## 0. 开工三核

1. HEAD＝`1b9721519` 在 `dev`（＝派单锚 `1b9721519` 本体）。
2. 索引干净（`git diff --cached --stat` 0 行）；共享树他人脏件在 `.scratch/**` 与 `internal/speech/**`（他腿在飞写面，与本腿零交）。
3. `git status --porcelain -- tools/mockllm` 空（写面无他人脏）。

## 1. 票面 5 格逐条（判定＋读数）

| 格 | 判定 | 读数 |
|---|---|---|
| AC#1 两方言补齐 max_tokens 应用（对照 chat 的 rune 边界形状——同形） | **成立**（§2） | messages/responses 合成路径各落截断块：帽 `len(answer)/4 > MaxTokens` 逐字同款（有意模拟 ⛔ 未动）、切点复用 chat 的 `runeBoundaryWithin(MaxTokens*4, answer)`（helper 未动）、`toolName == ""` 才截、reason/thinking 不截、usage 按截断后文本计——五个语义位与 chat 逐一同形 |
| AC#2 多字节回归两案（finish=length＋rune 边界前缀——对照 386-r2 案形） | **成立**（§3，⚠ 终止信号词汇落法见 §2 选型） | 新件 `caps_maxtokens_391r1_test.go`：2 测试 6 子案——messages 非流式/流式/ASCII 对照＋responses 非流式/流式/ASCII 对照；汉字切案算术与 386-r2 的 chat 案逐位同款（300 汉字＋帽 100 ⇒ `echo: `＋131 汉字＝399 字节） |
| AC#3 反向判据（三读数＋还原核） | **成立**（§4） | 修前红 → 修后绿 → `-overlay` 装回修前快照复红同位 → 撤 overlay 复绿；sha256 在案 |
| AC#4 零回归（mockllm 全包＋消费方 E2E） | **成立**（§5） | 模块内全包 `36/20/0/0` rc=0；`internal/llm/...` 七包＋`internal/present`（票 30 真管线）rc=0；cmd/wisp 定向名册（票 12 E2E）`46/36/0/0` rc=0；三 build＋vet rc=0 |
| AC#5 纪律格：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**` | **成立（自证）** | 提交 pathspec 显式逐枚恰 4 件（2 产码＋1 测试＋本证据件）；零 push；零翻框 |

## 2. AC#1 修法与关键选型（本票核心交付）

**帽本体与切点（三方言同形，照 386-r2 钉死的形状）**：两方言各在合成路径插入截断块——`if req.MaxTokens > 0 && toolName == "" { if len(answer)/4 > req.MaxTokens { answer = answer[:runeBoundaryWithin(req.MaxTokens*4, answer)]; truncated = true } }`。启发式（~4 bytes/token）、`toolName==""` 才截（forced tool 不截）、reason/thinking 不截、usage 按截断后文本计——全部与 `chat.go` planChat 的应用点同形；`runeBoundaryWithin` 直接复用 chat 的 helper（chat.go 一字未动）。落位：messages 在 `planMessagesFrom` 之后、capability 段之前（＝chat 的 plan 层同位，fcBroken 清 tool 前的 toolName 判定同序）；responses 在 tool 决策后、fcBroken 段前（同位）。

**★终止信号词汇的选型（技术判断，具名理由）**：票面 AC#2 的 "finish=length" 按 386-r2 案形理解为「截断终止信号」的泛称——386-r2 案形是 chat 方言的，`length` 是 OpenAI Chat 词汇。两方言各落**本协议已有的合法词汇槽位**：

- **Messages**：`stop_reason: "max_tokens"`（非流式顶层字段＋流式 `message_delta.delta.stop_reason`）。理由：①真 Anthropic provider 的 stop_reason 词汇表没有 `length` 这个值，照抄 chat 词汇恰好制造新的方言不齐——与本票机理相逆；②消费方实锤：`internal/llm/anthropic/adapter.go` 的 mapStop 把 wire `max_tokens` 映射到 C6 `max_tokens`，且对**未知 stop_reason 直接报错**（`anthropic: unknown stop_reason %q`）——报 `length` 会把消费方打崩。
- **Responses**：非流式 `status: "incomplete"`＋`incomplete_details.reason: "max_output_tokens"`；流式终结事件 `response.incomplete`（payload response 内带同两字段）。理由：①真 Responses API 无 finish_reason，命顶即 status=incomplete＋incomplete_details；②消费方实锤：`internal/llm/openairesponses/adapter.go` 映射表 `incomplete + max_output_tokens → C6 max_tokens`，流解析本就认 `response.completed/response.incomplete/response.failed` 三终结事件（`stream.go:278`）。

三方言终止信号在 C6 层殊途同归到同一个 `max_tokens`（chat `length→max_tokens`、Messages `max_tokens→max_tokens`、Responses `incomplete/max_output_tokens→max_tokens`）——三方言不一致的病灶由此消除。

**未定义即停两条件的核对**：①"两方言的 max_tokens 语义与 chat 不同"——**不成立**：`max_tokens`（Anthropic）与 `max_output_tokens`（Responses）与 chat 的 `max_tokens` 同义（输出 token 上限、命顶截断并报告），只是字段名与终止信号拼法各协议不同（方言保真的本职），不触发停手；②"修齐必须动协议形状"——**不成立**：`stop_reason`/`status` 都是协议**已有**字段、新值都是协议**合法词汇**，`response.incomplete` 是协议标准事件（消费方本就解析），SSE 事件结构、字段集、finish/usage 语义零改动，不停手。

## 3. AC#2 多字节回归两案（对照 386-r2 案形）

新件 `tools/mockllm/caps_maxtokens_391r1_test.go`（文件名尾段 `_391r1` 非 GOARCH 名——386-r2 §3 存档的 `_386_` 静默排除陷阱照避），2 测试 6 子案，样本算术与 386-r2 的 chat max_tokens 案逐位同款：

- `Test391MessagesMaxTokensTruncates`：300 汉字＋`max_tokens=100`（echo 906 字节＝启发式 226 tokens>100）⇒ 非流式断言 `stop_reason="max_tokens"`＋文本＝`"echo: "+测×131`（399 字节）＋`utf8.ValidString`＋≤400 字节界；流式断言 `message_delta` 的 `stop_reason="max_tokens"`＋装配文本同上；`ascii_cut_byte_identical_to_cap` 钉"帽仍按 `max_tokens*4` 字节截 ASCII（`echo: `＋194 a＝200 字节）"——**判有意的那半边没被顺手拆掉**。
- `Test391ResponsesMaxTokensTruncates`：同算术＋`max_output_tokens=100` ⇒ 非流式断言 `status="incomplete"`＋`incomplete_details.reason="max_output_tokens"`＋文本＝同款 rune 边界前缀；流式断言终结事件**是 `response.incomplete` 而非 `response.completed`**＋status/reason 两读数＋装配文本同上；`ascii_cut_byte_identical_to_cap` 同款。

## 4. AC#3 反向判据（三读数＋还原核）

突变件＝修前两文件快照（`D:\tmp\wisp\20261031\391r1\orig\`，sha256 `messages 4C03E0701A7EAB0D`／`responses 94E17C32CDDC5A88`——与 386-r2 §4 在案现值逐位一致，即自 386-r2 后两文件零改动），经 `go test -overlay` 注入，工作树零接触：

1. **读数①修前（死字段在）**：`go test -run Test391 -v` **rc=1**。红句逐字（`caps_maxtokens_391r1_test.go`）：messages `stop_reason = "end_turn", want "max_tokens"`、`text = 906 bytes / 306 runes, want the rune-boundary prefix 399 bytes / 137 runes`、`906 bytes, exceeds the max_tokens*4 heuristic bound 400`、ascii 对照 `text = 1006 bytes (want 200)`；responses `status = "completed", want "incomplete"`、`incomplete_details = <nil>`、同款 906/400 红句、ascii 对照同款。⚠ 如实报：首次跑 responses 流式子案红在 `frame not JSON`（测试件自身漏跳 `[DONE]` 哨兵，随即按 386-r2 chat 案形补上跳过）；死字段的红以其余五子案为准，回填后修前形态已由读数②③钉死。
2. **读数②修后**：同命令 rc=0，六子案全 PASS；模块内全包 `36/20/0/0` rc=0。
3. **还原核（overlay 装回修前快照）**：`go test -overlay D:\tmp\wisp\20261031\391r1\overlay.json -run Test391` **rc=1**，红点与读数①**同位同构**（六子案全红）——测试咬的是死字段，不是别处。
4. **撤 overlay**：rc=0 复绿；修后产码现值 sha256 `messages 874B18E009700465`／`responses 017F61A9763FE4E9`（overlay 机制上只读注入、工作树零接触，与快照值不同可证）。

## 5. AC#4 零回归（消费方 E2E＋build）

| 消费方 | 命令 | rc | 读数 |
|---|---|---|---|
| mockllm 包内（含新回归） | `go test -count=1 -v ./...`（tools/mockllm 模块内） | 0 | 四数 `RUN=36 / PASS=20 / FAIL=0 / SKIP=0`（386-r2 基线 28/18＋本腿 2 顶层 6 子案，增量吻合）；全包 0.05s |
| 两方言适配器＋chat 族＋票 30 E2E | `go test -count=1 ./internal/llm/... ./internal/present/`（根模块） | 0 | `llm 53.7s`／`adaptertest 2.9s`／`anthropic 18.0s`／`golden 0.05s`／`openaichat 27.1s`／`openairesponses 23.4s`／`present 3.1s` 全 `ok`；present 即票 30 真管线（真 mockllm 子进程） |
| cmd/wisp 消费方（票 12 E2E 等） | 386-r2 同款定向名册 `-run '^TestRunTextTask|^TestHostDispatch|…|^Test197NoAllowDoor' -v` | 0 | 四数 `RUN=46 / PASS=36 / FAIL=0 / SKIP=0`（与 386-r2 读数逐位一致）；前置 `third_party\sherpa-onnx` 上 PATH（`scripts/wisp-cli-tests.sh` 台型） |
| 编译 | 根模块 `go build ./...`；tools/mockllm 模块内 `go build ./...`＋`go vet ./...` | 0/0/0 | — |

注记（如实报）：
- cmd/wisp 读数按 `scripts/wisp-cli-tests.sh` 台型前置 `third_party\sherpa-onnx` 上 PATH 后取（票 98/111 在册的 DLL 形状）。
- 根模块 `go test ./cmd/wisp` 重包同一时刻只跑一发的规矩遵守：本腿全程单包串行。

## 6. 四数／两尺／d22scan

- **四数**（口径：`^=== RUN` 含子测试／`^--- PASS` 只顶层／FAIL／SKIP，⛔ 相减）：mockllm 包内 `36/20/0/0`；cmd/wisp 定向 `46/36/0/0`；d22scan `clean`。
- **两尺**：`gofumpt -l tools\mockllm` **0 行**（gofumpt v0.12.0, go1.27.1 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）、`gofmt -l tools\mockllm` **0 行**。
- **d22scan 真仓扫**：`cd tools/d22scan && go run . -root ../../` **rc=0 clean**；live scope 与 386-r2 读数一致：`bans#1-5 internal/=287`、`bans#1-5 cmd/=52`、`ban#6 frontend/=116`、`ban#7 internal/tools/=25`、`ban#8 design/=39`、`ban#8 frontend/=116`、`ban#8 internal/=729`、`ban#8 cmd/=180`。新测试件 CJK 字符（`测`）不在任何 emoji 码位带。

## 7. 登记不修

无新增。本票即 386-r2 §7 登记项①的落地；两方言截断落位后三方言 max_tokens 语义在 C6 层齐平，无残余死字段、无新缺陷面。

## 8. 交付形状

- 产码：`tools/mockllm/messages.go`（截断块＋非流式/流式 stop 词汇）、`tools/mockllm/responses.go`（截断块＋非流式 status/incomplete_details＋流式终结事件）。
- 回归：`tools/mockllm/caps_maxtokens_391r1_test.go`（新增，2 测试 6 子案）。
- 证据：本件。临时件（修前快照／overlay JSON）在 `D:\tmp\wisp\20261031\391r1\**`（仓外，只建不删）。
- 提交：主题逐字 `票 391 落地腿 391-r1`，`-F` 在 `--` 前，pathspec 显式逐枚恰 4 件；只 commit 不 push。
