# 票 385 落地腿 385-r1 · 实现件证据（mockllm 回声 400 字节截断修复）

- 腿：`385-r1`（写腿）。日期 2026-10-31。开工锚 HEAD＝`716874ae1`（分支 `dev`）。
- 写面：`tools/mockllm/**`（本腿实改 `chat.go`／`responses.go`／`messages.go`／新增 `chat_echo_385_test.go`）＋本证据件。未动 `cmd/wisp`、`internal/**`、票面、台账、testdata、specs、push。
- 环境：`GOFLAGS=`（空）、`GOTMPDIR` 未设、gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`（v0.12.0, go1.27.1）、突变件在仓外 `D:\tmp\wisp\20261031\385r1\**` 经 `-overlay` 注入、提交 `-F` 在 `--` 前、主题逐字 `票 385 落地腿 385-r1`。

## 0. 开工三核

1. HEAD 在 `dev` 且 `git merge-base --is-ancestor 0d240539d HEAD` rc=0（≥ 派单锚）。
2. `git diff --cached --name-only` 空（索引干净）。
3. `git status --porcelain -- tools/mockllm` 空（写面无他人脏）。
   共享树另有他人 WIP（`cmd/wisp/resident_voice_*.go` 等四枚在改）——与本腿写面无交；对本腿读数的影响在 §5 注记。

## 1. 票面 5 格逐条（原文＋判定＋读数）

| 格 | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| AC#1 | 按 rune 截断或去掉截断（对照该截断的原始意图——防超长回声；修法具名） | **成立（修法＝②去掉截断，具名理由 §2）** | 三方言（chat/responses/messages）回声帽全摘；防超长意图的两条真护栏保留（8MB LimitReader＋max_tokens 模拟），见 §2 |
| AC#2 | 多字节回归一发：402 汉字案走 mockllm 不再拦腰切（与 30-r1 的 C5 缝案同输出口径对照） | **成立** | 票 30 tier-2 同一样本（`测`×402＋`\n口播摘要：长文已归档`，echo 全文 1243 字节/419 rune）三模式全量回传；与 C5 缝案两路口径一致（§3） |
| AC#3 | 反向判据：摘修复 ⇒ 402 汉字案必红（三读数＋还原核） | **成立** | `-overlay` 装回字节帽必红（rc=1，红句 402 bytes/138 runes、尾部 U+FFFD）；三读数＋sha256 还原核见 §4 |
| AC#4 | 零回归：mockllm 现有消费方（票 12/30 的 E2E）全绿 | **成立** | 包内 19/16/0/0 rc=0；`internal/llm/...`＋`internal/present`（票 30 E2E）7 包 rc=0；`cmd/wisp` 定向消费方（票 12 E2E 等 37 枚顶层）46/36/0/0 rc=0＋ComposedGate 单发 PASS；`go build ./...` rc=0（§5） |
| AC#5 | 纪律格：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`（只 `tools/mockllm/**`） | **成立（自证）** | `git status --porcelain` 本腿名下恰 4 枚 `tools/mockllm` 文件＋本证据件；提交 pathspec 显式逐枚；零 push |

## 2. 原始意图对照与修法具名（AC#1）

**blame/出处**：`git blame -L 145,170 -- tools/mockllm/chat.go` 逐行同源＝**建仓提交 `0839407bf`**（2026-09-19 22:50 +0800，"feat(mockllm): standalone stdlib-only mock LLM server - chat/responses/messages SSE synthesis, golden replay, /v1/models, fault-injection control endpoints"）。该提交之后再无一笔碰过截断行；截断行**无专属注释**，commit message 也未专门解释。`git log -S 'answer[:400]'` 全仓仅此一笔。

**意图可辨的依据**（故未触发停手条件①）：同一提交里三个方言**同型对称**——`chat.go` `answer[:400]`＋`reasoning[:200]`、`responses.go` `answer[:400]`＋`reason[:200]`、`messages.go` `p.answer[:400]`＋`p.reason[:200]`。这是刻意的「合成回声封顶」惯例，不是孤行手误；意图＝防超长回声（与票面预判一致）。

**修法具名＝② 去掉截断**（非 ①按 rune 截断），理由四条：

1. **该帽与真护栏重复**：请求体已有 8MB 上界（`chat.go:21` `io.LimitReader(r.Body, 8<<20)`）——回声长度物理有界；按请求的长度控制已由 max_tokens 模拟承担（`chat.go:205-209`，含 `finish_reason=length` 语义）。400 字节帽是第二套更粗糙、不模拟任何真实 provider 行为的合成上限。
2. **无任何消费方测量该帽**：包内测试零断言（本腿逐文件核过，`mockllm_test.go`/`capability_test.go`/`thinking_capability_test.go` 全部短回声）；golden 模式整条绕过 `planChat`（`serveGoldenChat` 直放录制字节）；唯一被它挡死的正是票 30 的 402 汉字案。
3. **AC#2 的「两路输出一致」只有去掉截断能满足**：C5 缝假体回**全文**（`internal/present/e2e_present_30_test.go:144` 逐字 `&scriptedC5{text: strings.Repeat("测", 402) + "\n口播摘要：长文已归档"}`）；rune 截 400 仍把 420-rune 样本截短到 400 rune——「echo: 」前缀＋394 字后**摘要行整行丢失**，两路口径必然不一致。
4. **主流 harness 判断**：test double 的职责是忠实表达测试场景；不模拟真实行为的隐藏合成上限正是「mock 表达不了场景」这一类陷阱（票 30 被迫绕行就是实例）。长度控制归位到请求语义（max_tokens）而不是 mock 内的暗常数。

**协议形状零改动**：SSE chunk 结构、字段、finish/usage 语义全部未动；变的只是 content 可以有多长（正是「只动截断」的授权范围）。

## 3. AC#2 多字节回归一发（402 汉字案）

新回归件 `tools/mockllm/chat_echo_385_test.go`（`Test385EchoDeliversLongMultibyteInFull`）以**票 30 tier-2 逐字同一样本**钉三模式：

- 样本：`strings.Repeat("测", 402) + "\n口播摘要：长文已归档"`（样本本体 1206＋1＋30＝1237 字节；带 mock 的 `echo: ` 签名前缀共 **1243 字节 / 419 rune**；旧帽 400 字节会切在第 395 字节＝第 132 个汉字的 UTF-8 序列中间）。
- 修后读数：`chat_stream_sse`（SSE 逐 delta 装配，与 C5 缝案同走真 adapter 的流式口径）／`chat_nonstream_json`／`messages_and_responses_dialects` 三子测试 **全 PASS**；装配文本逐字节等于 `"echo: " + 样本`，`utf8.ValidString` 成立。
- **两路口径一致的读数**：C5 缝案侧＝`scriptedC5` 原样 emit 样本文本（一个 `EvTextDelta`，`e2e_present_30_test.go:52-61`）；mockllm 侧修后＝`"echo: " + 样本` 全量。剥掉 mock 的 `echo: ` 签名前缀后**两路逐字节相同**（签名前缀本身就是 mockllm 模式的既定形状，票 30 tier-1 的断言 `rc.Display == "echo: 生成结果如下：…"` 同款在案）。摘帽前该案物理回不来（§4 的红句即证）。

## 4. AC#3 反向判据（三读数＋还原核）

突变＝把原始字节帽逐字装回 `chat.go`（`if len(answer) > 400 { answer = answer[:400] }`）；突变件在仓外 `D:\tmp\wisp\20261031\385r1\chat.mut-clamp-385.go`，经 `go test -overlay` 注入，**工作树零接触**。

1. **修后（基线）**：`go test ./... -v` rc=0，四数 `RUN=19 / PASS=16 / FAIL=0 / SKIP=0`，Test385 三子测试全绿。
2. **装回帽（overlay）**：`go test -overlay ... -run Test385` **rc=1**，红句逐字：
   - `chat_stream_sse`：`assembled text = 402 bytes / 138 runes, want 1243 bytes / 419 runes (tail lost: "测…�")` —— 尾部可见 **U+FFFD**，即字节帽拦腰切出的坏序列经 JSON 化的残迹；402 bytes＝`echo: `(6)＋131 汉字(393)＋坏字节变 U+FFFD(3)。
   - `chat_nonstream_json`：`content = 402 bytes, want 1243 bytes with the sample intact`。
   - `messages_and_responses_dialects` 保持 **PASS**（突变只落 `chat.go`——红点外科级正中机理位点）。
3. **撤 overlay**：`go test -run Test385` rc=0，三子测试复绿。
4. **还原核**：overlay 前后三产码文件 sha256 前缀逐一相同（chat `6EAB45DBBE4F5B03`／responses `BD1D1107A092102A`／messages `D74E06418BA1573A`）；`git status --porcelain -- tools/mockllm` 恰本腿四件。

## 5. AC#4 零回归（消费方 E2E＋build）

| 消费方 | 命令 | rc | 读数 |
|---|---|---|---|
| mockllm 包内（含新回归） | `go test ./... -v -count=1`（tools/mockllm 模块内） | 0 | 四数 `RUN=19 / PASS=16 / FAIL=0 / SKIP=0`（RUN 含 3 子测试） |
| chat 方言族＋票 30 E2E | `go test ./internal/llm/... ./internal/present/ -count=1` | 0 | `internal/llm 50.2s`／`adaptertest 3.2s`／`anthropic 18.0s`／`golden 0.1s`／`openaichat 25.5s`／`openairesponses 23.4s`／`present 3.5s` 全 `ok`；present 即票 30 四案真管线（真 mockllm 子进程） |
| cmd/wisp 消费方（票 12 E2E 等） | `go test ./cmd/wisp/ -run '^TestRunTextTask|^TestHostDispatch|^TestMissingBlob|^TestTicket101|^Test365|^TestTaskOutputPointerNotice174r8|^TestProviders|^TestReplyListener|^TestPanelRouteRefuses|^TestUnansweredL2|^TestL1Veto|^Test375|^TestRunPacketCarriesTheSubagent|^TestSubagentStreamKey|^TestRunPacketReportsTheStreamLog|^Test197Subagent|^Test197NoAllowDoor' -v -count=1` | 0 | 四数 `RUN=46 / PASS=36 / FAIL=0 / SKIP=0`（RUN 含 10 子测试）；含票 12 `TestRunTextTaskTextPathEndToEnd` |
| cmd/wisp 慢枚单发 | `go test ./cmd/wisp/ -run '^TestComposedGateBlocksAWriteForTwoSeconds$' -v` | 0 | PASS（3.84s） |
| 编译 | 根模块 `go build ./...`；tools/mockllm 模块内 `go build ./...`＋`go vet ./...` | 0/0/0 | — |

注记（如实报）：
- cmd/wisp 首发裸跑 `exit status 0xc0000135`（DLL 未上 PATH）＝票 98/111 在册形状；按 `scripts/wisp-cli-tests.sh` 的台型前置 `third_party\sherpa-onnx` 后全绿。GOTMPDIR 未设（票 310/311 守卫天然过）。
- cmd/wisp 该包编译含他人四枚在改文件（resident_voice_*.go WIP）——共享树读数，非本腿可隔离；本腿字节不进该包（机制上 tools/mockllm 是独立模块，测试二进制对它的引用经 `adaptertest.StartMockllm` 子进程构建）。
- `cmd/wisp/subagent_carrier_197_test.go:605-608` 的「WHAT THIS FILE DID NOT MEASURE」第 2 条以本帽为不可达论证（echo 被帽在 400 字节 ⇒ 行 ~150 rune）。摘帽后该**论证前提失效**，但其全部**断言实测仍绿**（上面 46/36/0/0 含其三枚）。该注释归 cmd/wisp 所有者更新——本腿⛔ 越面，登记 §7。

## 6. 四数／两尺／d22scan（tools/ 射程如实报）

- **四数**（口径：`^=== RUN` 含子测试／`^--- PASS` 只顶层／FAIL／SKIP，⛔ 相减）：mockllm 包内 `19/16/0/0`；cmd/wisp 定向 `46/36/0/0`；d22scan 自检 `77/35/0/0`（RUN 含子测试）。
- **两尺**（甲＝gofumpt、乙＝gofmt，均为路径参数型）：`gofumpt -l tools\mockllm` **0 行**（rc=0）、`gofmt -l tools\mockllm` **0 行**（rc=0）。指向 `tools/mockllm` 即在射程内，本腿已跑；两尺没有 repo 级自动巡检面（它们只在被调用时工作）。
- **d22scan 真仓扫**：`cd tools/d22scan && go run . -root ../../` **rc=0 clean**，live scope 逐行：`bans#1-5 internal/=280`、`bans#1-5 cmd/=45`、`ban#6 frontend/=116`、`ban#7 internal/tools/=25`、`ban#8 design/=39`、`ban#8 frontend/=116`、`ban#8 internal/=714`、`ban#8 cmd/=162`。
- **tools/ 在 ban#8 射程吗——不在，如实报**：ban#8 的 scope 表只有 design/＋frontend/＋internal/＋cmd/（源码注释 `tools/d22scan/main.go:40-49` 与实跑 scope 表双向一致）；ban#7 的 `internal/tools/` 是仓内另一目录，与顶层 `tools/`（mockllm/d22scan 所在）无关。⇒ 顶层 `tools/**` 在八枚 ban 的任何一枚射程之外；本腿新增测试件里的 CJK 字符（`测`／`口播摘要`）也不在任何 emoji 码位带内。
- **正控**：仪器自带栽赃样本族自检 `go test ./tools/d22scan` rc=0（`TestScanDetectsAllSeededViolations` PASS＝ seeded 违例必被咬；`TestScannerSelfScanOfRealRepoIsGreen` PASS）。

## 7. 登记不修（残余，具名归 tools/mockllm 所有者后续开票）

1. **reasoning/reason 200 字节帽三方言**——与本次修复同类（字节切多字节文本拦腰切风险）：`chat.go:168-169`（`reasoning[:200]`）、`responses.go:100-101`（`reason[:200]`）、`messages.go:114-115`（`p.reason[:200]`）。票面 机理 只点名 400 字节回声帽，未点名 reasoning 帽 ⇒ 本腿不越 AC 顺手修。
2. **max_tokens 字节切**——`chat.go:205-209` `answer[:req.MaxTokens*4]`：~4 字节/token 启发式同样按字节切，长多字节答案＋小 max_tokens 时可拦腰切；且它带 `finish=length` 的模拟语义，动它是改模拟行为，超出「只动截断（回声帽）」授权。
3. **cmd/wisp/subagent_carrier_197_test.go:605-608 的过期论证注释**——引用本帽为不可达理由，帽摘后前提失效（断言实测仍绿）；归 cmd/wisp 所有者改注释，本腿⛔ 越面。

## 8. 交付形状

- 产码：`tools/mockllm/chat.go`（摘帽＋6 行注释）、`tools/mockllm/responses.go`（摘帽＋2 行注释）、`tools/mockllm/messages.go`（摘帽＋2 行注释）。
- 回归：`tools/mockllm/chat_echo_385_test.go`（新增，三子测试）。
- 证据：本件。临时件（突变件／overlay JSON／全部原始 log）在 `D:\tmp\wisp\20261031\385r1\**`（仓外，只建不删）。
- 提交：主题逐字 `票 385 落地腿 385-r1`，`-F` 在 `--` 前，pathspec 显式逐枚；只 commit 不 push。
