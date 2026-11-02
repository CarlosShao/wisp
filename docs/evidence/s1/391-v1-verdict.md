# 票 391 验收腿 391-v1 · 非实现者裁决表（mockllm 两方言 max_tokens 死字段）

- 腿：`391-v1`（验收腿＝非实现者，裁决者≠实现者，`SPEC-12 §4.3` #1/#3）。日期 2026-11-02。
- 被验件：commit `6cf941b72265b688adc65ced142e7ef166fcb295`「票 391 落地腿 391-r1」（作者 CarlosShao，2026-10-31 22:01 +08）。
- 票面：`.scratch/wisp/issues/391-mockllm-maxtokens-dead-fields.md`（12 行，5 格，**本腿未翻框、未改一字**）。
- 本腿角色：只读复核＋本证据件。全程零产码改动、零 `cmd/wisp`/`internal/**` 写面、零票面/台账、零 push。
- 临时件（仓外，只建不删）：`D:\tmp\wisp\20261102\391v1\**` —— `orig/`（**本腿自 `6cf941b72^` 现取**的两文件修前快照）、`mut/`（本腿自制的四枚突变体＋两枚消费方端到端探针件）、`overlay-*.json`、`pre-overlay-test.txt`、`mut-A/B/C/D.txt`、`restored-green.txt`、`full-v.txt`、`run391.txt`、`http-*.json/.txt`、`http-readings.txt`、`wiring.txt`、`consumer.txt`、`mockllm-391v1.exe`。

## 0. 起手三核（现跑）

| 核 | 判据 | 读数 |
|---|---|---|
| HEAD 现盘最新 | ≥ `a4949255e` | HEAD＝`250735423`（分支 `dev`）；`git merge-base --is-ancestor a4949255e HEAD` **rc=0**；`git merge-base --is-ancestor 6cf941b72 HEAD` **rc=0**（被验件确在现盘链上） |
| 索引干净 | 暂存区空 | `git diff --cached --name-only` **空**；tracked 脏件**只有一枚**＝`internal/risk/pathresolver_budget_norace_test.go`（＝派单预告的 `86-r1` 在写面，本腿零接触、零判定，见 §5 说明）；余脏全在未 tracked 的 `.scratch/**` |
| 写面无他人脏 | `git status --porcelain -- tools/mockllm` | **空**（起手与还原核后各跑一次，两次皆空） |

## 1. 五格逐条裁决（票面原文逐字＋判定＋本腿自跑读数）

| 格 | 票面原文（逐字） | 裁决 | 本腿读数（⛔ 转抄，全部现跑） |
|---|---|---|---|
| **AC#1** | AC#1 两方言补齐 max_tokens 应用（★**对照 chat 方言 386-r2 的 rune 边界形状——同形**） | **PASS** | 两方言死字段**真消除**：`messages.go:202-207`／`responses.go:134-140` 各有应用点，五个语义位与 `chat.go:215-220` **逐位同形**（§2 对账表）；`chat.go` 本票零接触（`git show --name-only` 无 `chat.go`）；usage 按截断后文本计＝HTTP 现跑 `output_tokens=100` 实证（§6） |
| **AC#2** | AC#2 多字节回归：★**两方言各一发 max_tokens 汉字切案**（finish=length＋rune 边界前缀——对照 386-r2 案形） | **PASS** | 新件 `caps_maxtokens_391r1_test.go` 2 测试 6 子案：两方言各「非流式／流式／ASCII 对照」，汉字样本 300×`测` ＋ 帽 100 ⇒ `echo: `＋131 汉字＝399 字节，算术与 386-r2 chat 案逐位同款；`utf8.ValidString`＋逐字相等双断（§3）。⚠ 票面「finish=length」是 chat 词汇，落地按各协议自己的合法词汇（§4 突变体证明这条断言独立咬得住）——**判合规**，口径差记 §8-2 |
| **AC#3** | AC#3 反向判据（★三读数＋还原核） | **PASS** | 本腿**自建**修前快照（sha 与实现件 §4 逐位一致＝独立复现）＋四枚独立突变体，共 **6 轮读数**：修前快照 **rc=1／6 子案全红／14 句红读数**；错方言（chat 词 `length` 落进 Messages）**rc=1／3 红**；字节切（撤 rune 边界）**rc=1／2 红、ASCII 对照恒绿**；终止仍报 `completed` **rc=1／3 红**；`incomplete_details.reason` 错词 **rc=1／2 红**；撤全部 overlay **rc=0** 六子案复绿；还原核＝`tools/mockllm` 工作树空、两文件 sha256 前缀 `874B18E009700465`／`017F61A9763FE4E9` 全程未变（§4） |
| **AC#4** | AC#4 零回归（★mockllm 全包＋消费方 E2E） | **PASS** | 模块内全包 `-v`＝**`RUN=36 / 顶层 PASS=20 / FAIL=0 / SKIP=0` rc=0**（与实现件 §1 四数逐位吻合）；消费方 `go test -count=1 ./internal/llm/... ./internal/present/`＝**7 包全 `ok`／rc=0**；★本腿**加测**两枚端到端探针（`-overlay` 注入 `internal/llm/{anthropic,openairesponses}`，⛔ 落盘）：真适配器 ⇄ 真 mockllm 子进程，两方言截断在 C6 层各自落 `stop="max_tokens"`、文本＝同一 399 字节 rune 前缀、`usage.output_tokens=100`（§7）；派单原文那把尺在根模块必 `setup failed`＝独立子模块的落点问题，⛔ 件的问题（§5 双记） |
| **AC#5** | AC#5 纪律格：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**` | **PASS** | `git show --stat 6cf941b72`＝**恰 4 件 / +438 −10**，主题逐字「票 391 落地腿 391-r1」＝题面一致；`--name-only` 无 `cmd/wisp`、无 `internal/**`、grep `^\.scratch/\|^docs/reports/` **空**＝票面与台账零接触；票面 5 框至今全 `[ ]`、Status 仍「待派」（结案归编排者，本腿亦⛔ 翻框）；`chat.go`/`testdata` 零接触（§8） |

## 2. AC#1 死字段 ⇄ 应用 逐位对账（本票核心）

票面机理点名两枚死字段：`responses.go:67`、`messages.go:179`。本腿从 `6cf941b72^` 现取修前快照核对：**修前两文件除解析行外零应用点**（`Select-String MaxTokens` 于修前快照只命中结构体声明那两行，`truncated` 变量不存在）⇒ 机理成立。票面所引行号与修前快照实测差一枚（`messages` 声明在快照 L180，票面写 179；`responses` L67 与票面逐字一致）——同一枚字段，⛔ 影响判定，记 §8-3。

HEAD 盘上五语义位 ⇄ chat 方言（`chat.go:215-220`，386-r2 钉死的形状）逐位对账：

| 语义位 | chat（对照基准） | Messages | Responses | 同形？ |
|---|---|---|---|---|
| 生效条件 | `:215 if req.MaxTokens > 0 && toolName == ""` | `:202` 同（`p.toolName`） | `:135` 同 | 同形 |
| 帽启发式（有意模拟） | `:216 len(answer)/4 > req.MaxTokens` | `:203` 逐字同款 | `:136` 逐字同款 | 同形（帽本体 ⛔ 未被顺手拆） |
| 切点（rune 边界） | `:217 answer[:runeBoundaryWithin(req.MaxTokens*4, answer)]` | `:204` 复用同一 helper | `:137` 复用同一 helper | 同形（helper 在 `chat.go:365`，本票未动） |
| 只截正文 | `reasoning` 不入切点 | `p.reason` 不截（`p.answer` 单点） | `reason` 不截 | 同形 |
| usage 按截断后计 | `chatUsage(answer,…)` L343 用截断后的 `answer` | `:226 completionTokens = len(p.answer)/4+…`（在 `:204` 之后） | `:148 outTokens = len(answer)/4+…`（在 `:137` 之后） | 同形（HTTP 现跑 `output_tokens=100`＝399/4+1，坐实） |
| 终止信号 | `:218 finish="length"`（chat 合法词） | `:247`/`:300` `stop_reason="max_tokens"` | `:180-181` 非流式 `status=incomplete`＋`incomplete_details.reason=max_output_tokens`；`:284-291` 流式终结事件 `response.incomplete` | 形状同形、**词汇各归本协议**（见下） |

**「正确接线」的消费方实锤（本腿现读 HEAD，非转抄实现件）**：

- `internal/llm/anthropic/adapter.go:285` `case "max_tokens": → llm.StopMaxTokens`；同文件 `:302-307` default 分支对**未知 stop_reason 直接报错**（`anthropic: unknown stop_reason %q`）⇒ 若照抄 chat 的 `length` 进 Messages，消费方会被打崩。实现腿选 `max_tokens` **不是自行放宽，是被消费方契约钉住的唯一正确落法**。
- `internal/llm/openairesponses/adapter.go:235-244` `status=="incomplete"` 且 `incompleteReason=="max_output_tokens"`（或空）`→ llm.StopMaxTokens`；`stream.go:278-280` 三终结事件白名单本就含 `response.incomplete`，`:411-412` `statusFromType("response.incomplete")="incomplete"` ⇒ 选 `incomplete`＋`max_output_tokens` **落在消费方已解析的既有槽位**，零新字段。
- 三方言在 C6 归一到同一 `max_tokens`（`internal/llm/events.go:73` `StopMaxTokens`）＝票面机理「三方言不一致」病灶的正解；本腿 §7 端到端探针**读数级**证实两方言都真到了这一格。

## 3. AC#2 回归件形状核对（逐行读）

- `Test391MessagesMaxTokensTruncates`：非流式断 `stop_reason=="max_tokens"`＋文本＝`"echo: "+测×131`（399 字节）＋`utf8.ValidString`＋`≤ maxTokens*4`；流式逐帧装配 `content_block_delta.text` 后同断，且**额外断 `message_delta.delta.stop_reason`**（流/非流两槽分钉）；ASCII 对照 `max_tokens=50` ⇒ `echo: `＋194 个 `a`＝200 字节（rune 边界对 ASCII 即字节切）。
- `Test391ResponsesMaxTokensTruncates`：非流式断 `status=="incomplete"`＋`incomplete_details.reason=="max_output_tokens"`（且该字段用指针，缺失即红，不会静默放过）＋同款 399 字节前缀；流式断终结事件**必须是 `response.incomplete` 而非 `response.completed`**（`terminal` 三值白名单一起收，改错事件名同样红）；ASCII 同款对照。
- 与 386-r2 chat 案形对照关系成立：同样本长度、同「帽 100 ⇒ 131 汉字」算术、同 ASCII 对照钉「帽未被拆」。文件名尾段 `_391r1` 避开 `_386_` 撞 GOARCH 的静默排除陷阱（386-v1 §5 在册提醒）；本腿用 `go list -f "{{.TestGoFiles}}"` 口径确认该件**在测试文件名册内**（§4 六轮读数它确实在跑）。

## 4. AC#3 反向判据（本腿独立重做，⛔ 转抄实现腿临时件）

修前快照由本腿自 `6cf941b72^` 现取（sha256 前缀 `messages 4C03E0701A7EAB0D`／`responses 94E17C32CDDC5A88`——与实现件 §4 所记逐位一致，＝两条独立通道取到同一份修前件）。全部经 `go test -overlay` 注入（只读，工作树零接触）。

| # | 突变内容 | rc | 红子案 | 红读数（具名行） |
|---|---|---|---|---|
| ① | **恢复旧死字段**（两文件整体换回修前快照） | **1** | **6/6 全红** | 14 句：`messages.go:52 stop_reason = "end_turn", want "max_tokens"`／`:65 text = 906 bytes / 306 runes, want … 399 bytes / 137 runes`／`:69 906 bytes, exceeds the max_tokens*4 heuristic bound 400`／`:102 stream stop_reason = "end_turn"`／`:108` 流式装配同读数／`:143 ascii control: text = 1006 bytes (want 200), stop_reason = "end_turn"`；`responses :181 status = "completed", want "incomplete"`／`:185 incomplete_details = <nil>`／`:201 :205` 同款 906/400／`:250 terminal event = "response.completed", want "response.incomplete"`／`:254 terminal status = "completed" reason = ""`／`:260` 流式装配 906／`:299 ascii control: text = 1006 bytes (want 200)` |
| ② 撤 overlay（基线复绿） | — | **0** | 6/6 全绿 | `ok … 0.083s` |
| ③ 突变 A：**错方言**＝只把 Messages 终止词换成 chat 的 `length` | **1** | 3（messages 全组） | `:52 :102 :143` 三句皆 `stop_reason = "length", want "max_tokens"`；⚠ 关键＝**文本字节断言同轮为绿**（`:143` 那行前半 `text = 200 bytes (want 200)`）⇒ 词汇维度被**单独**咬住，⛔ 与截断维度混红 |
| ④ 突变 C：**撤 rune 边界**（`p.answer[:req.MaxTokens*4]` 字节切） | **1** | 2（非流式＋流式） | `:65 text = 402 bytes / 138 runes, want … 399/137（tail "…测测测…"）`、`:69 402 bytes, exceeds … 400`、`:108` 流式同读数；**ASCII 对照 PASS** ⇒ 多字节那一半单独咬住，且旧缺陷（越自身启发式界＋坏码尾）被复现 |
| ⑤ 突变 B：**截了但仍报 completed**（非流式 `status`＋流式终结事件都改回 completed） | **1** | 3（responses 全组） | `:181 status = "completed", want "incomplete"`、`:250 terminal event = "response.completed"`、`:254 terminal status = "completed" reason = "max_output_tokens"`、`:299 ascii control: text = 200 bytes (want 200), status = "completed"`；⚠ `:254` 的 `reason` 仍读到正确值＝证明红点落在**状态/事件**这一维，⛔ 蹭到 reason 维 |
| ⑥ 突变 D：**`incomplete_details.reason` 错词**（`content_filter`） | **1** | 2 | `:185 incomplete_details = &{Reason:content_filter}, want reason "max_output_tokens"`、`:254 … reason = "content_filter"`；**ASCII 对照 PASS**；★这枚突变在真 provider 语义里不是玩笑：消费方 `adapter.go:237-238` 会把 `content_filter` 映射成 `StopContentFilter`，＝C6 语义走歪，测试把它钉住了 |

**还原核**：六轮之后 `git status --porcelain -- tools/mockllm` **空**；`messages.go`／`responses.go` sha256 前缀仍＝起手读数（`874B18E009700465`／`017F61A9763FE4E9`）；`chat.go` 全程零接触（本腿未写、突变体只覆盖两文件）⇒ 突变没落到盘。

⇒ 三读数齐（绿→装回死字段复红同位→复绿）＋**四枚定向突变各自单独复红**＋还原核过。**测试咬的是死字段与方言词汇本身，不是别处。** 比票面最低要求（三读数）多四枚定位读数，记为加分项。

## 5. AC#4 零回归读数

| 项 | 本腿命令 | rc | 读数 |
|---|---|---|---|
| 派单原文尺（照跑） | 根模块 `go test ./tools/mockllm -count=1 -timeout 10m` | **1** | ⚠ **非测试失败**：`main module (github.com/CarlosShao/wisp) does not contain package …/tools/mockllm [setup failed]`——`tools/mockllm` 自带 `go.mod`＝独立模块，根模块按相对路径点它是**加载期错误**。尺的落点问题、不是件的问题；与 386-v1 §5 同一枚已在册口径（386-v1 §7-3 已建议派单模板改写本尺，⛔ 本腿改票面/模板） |
| 正确落点（全包 `-v`） | `cd tools/mockllm; go test -count=1 -timeout 10m -v .` | **0** | `ok github.com/CarlosShao/wisp/tools/mockllm 0.180s`；四数 **`RUN=36 / 顶层 PASS=20 / FAIL=0 / SKIP=0`＝FAIL=0**（与实现件 §1/§5 的 `36/20/0/0` 逐位吻合） |
| 点名用例 | 同上 `-run Test391 -v` | **0** | 2 顶层／6 子案全 `--- PASS`（§4 基线轮） |
| 消费方 E2E | 根模块 `go test -count=1 -timeout 10m ./internal/llm/... ./internal/present/` | **0** | 7 包全 `ok`：`llm 58.8s`／`adaptertest 4.4s`／`anthropic 22.6s`／`golden 0.12s`／`openaichat 32.6s`／`openairesponses 30.6s`／`present 4.7s`（present＝票 30 真管线，真 mockllm 子进程）——摘死字段变"有应用"后，消费方无一被绊倒 |
| ★本腿加测：端到端接线 | `go test -overlay wiring -run Test391V1 ./internal/llm/anthropic/ ./internal/llm/openairesponses/` | **0** | 两枚探针件由 `-overlay` **虚拟注入**（仓内零新文件，跑后 `git status -uno` 仍只剩 `internal/risk` 那枚他腿件）；读数：`anthropic: C6 stop="max_tokens" bytes=399 runes=137 incomplete=false usage={InputTokens:17 OutputTokens:100}`；`openairesponses: 同读数` ⇒ **C5 适配器 ⇄ 真 mockllm 子进程**双向坐实"两方言截断确实被消费方看见并归一到 C6 max_tokens"（＝票面机理的正解，⛔ 只停在协议字面） |
| 编译／静态 | `go vet .`／`go build ./internal/... ./cmd/...`（根模块） | **0** | vet rc=0；两路 build rc=0（`True/True`） |
| 两尺 | `gofumpt -l .`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）／`gofmt -l .` | 0/0 | 各 **0 行** |
| d22scan | `cd tools/d22scan; go run . -root ../../` | **0** | `clean`；live scope `bans#1-5 internal/=288`、`cmd/=52`、`ban#6 frontend/=116`、`ban#7 internal/tools/=25`、`ban#8 design/=39`、`internal/=740`、`cmd/=186`（与 386-v1 §5 读数同值，⛔ 新面）；`tools/**` 在八枚 ban 射程外，新测试件的 `测` 不落任何 emoji 码位带（386-v1 同口径复核） |
| `cmd/wisp` 定向名册 | **本腿未跑** | — | 如实记：实现件 §5 报 `46/36/0/0` rc=0，本腿**未独立复现**。理由（可核）＝`go list -deps ./cmd/wisp` 含 `github.com/CarlosShao/wisp/internal/risk`，而该包现有一枚他腿（`86-r1`）**未提交**的在写件——跑出来的红绿无法归因到本票。⛔ 误判他腿＝本票 AC#4 的两条已覆盖核心（mockllm 全包 FAIL=0＋消费方 E2E rc=0），故不据此扣格；若编排者要补这一读数，应在 `86-r1` 落盘后单独取 |

## 6. 线级读数（真子进程，非 httptest）

本腿另起 `mockllm-391v1.exe`（编译自现盘）于 `127.0.0.1:18097`，用 UTF-8 原始 body 打三方言同一算术（`300×测`＋帽 100），读数（`http-readings.txt`／`http-*.json|.txt` 在案）：

```text
MSGS nonstream stop_reason=max_tokens   bytes=399 runes=137 validUtf8=True  out_tokens=100
RESP nonstream  status=incomplete reason=max_output_tokens bytes=399 runes=137 out_tokens=100
CHAT same-cap   finish=length           bytes=399 runes=137
SAME_PREFIX chat == messages == responses: True
RESP control(no cap) status=completed bytes=906
MSGS stream  message_delta "stop_reason":"max_tokens"  ×1
RESP stream  "type":"response.incomplete" ×1 / "response.completed" ×0
```

⇒ 三方言在**同一请求算术**下给出**逐字节相同的截断文本**，终止信号各自落在本协议的合法词汇上，不设帽时行为不变（906 字节全文）。`RESP control` 那条同时反证：改动没把"帽 ≤ 文本"的正常路径碰歪。原始帧在案，含 `incomplete_details` 与 `usage.output_tokens=100` 的完整 JSON。

## 7. 与实现件的口径差／加分（本腿独立看，非复述）

- 实现件 §4 自报一处过程瑕疵：首次跑 responses 流式子案红在 `frame not JSON`（测试件漏跳 `[DONE]`），随即照 386-r2 chat 案形补上。本腿在修前快照轮**亲眼见该子案红在 `:250 terminal event` 断言上**（⛔ 再是 JSON 解析崩）⇒ **该瑕疵已闭合**，且闭合后的断言仍咬得住死字段。
- 实现件把票面 AC#2 的「finish=length」读成「截断终止信号」泛称、按协议词汇分落 `max_tokens`／`incomplete+max_output_tokens`：本腿独立复核消费方代码（§2 两处实锤）后判**该读法成立、且是唯一正确读法**——照字面抄 `length` 会踩 `adapter.go:307` 的未知 stop_reason 报错，那将是「修一个方言不齐、造出一个新方言不齐」，与票面机理相反。
- 加分：本腿的突变体 A/C/B/D 证明「词汇维」「切点维」「状态维」「reason 维」**四维各自独立可红**＝回归件不是靠一句复合断言蒙住的，粒度够。

## 8. 纪律格读数明细

- `git show --stat 6cf941b72`＝**恰 4 件 / +438 −10**：`tools/mockllm/{messages.go,responses.go,caps_maxtokens_391r1_test.go}`＋`docs/evidence/s1/391r1-impl.md`。题面（commit message）逐件自述：两产码＋一枚回归＋本证据件 ⇒ **与 --stat 逐件对得上，无未申报面**。
- 写面：三枚产码/测试件全在票面 `写面＝tools/mockllm/**` 内；证据件在写面外一格——**判合规**（`AGENTS.md §1.5` 与 385-r1/386-r1 同型先例：落地腿自身证据件落 `docs/evidence/s1/`，属规程指定落点，⛔ 越面写码）。
- 票面台账零接触：`--name-only` grep `^\.scratch/`／`^docs/reports/` **空**；`git log -- .scratch/wisp/issues/391-*.md` **只有开票那一枚** `1b9721519`（编排者），落地腿⛔ 动票面；票面五框至今 `[ ]`、Status 仍「待派」。
- 零越票顺手修：`chat.go` 不在本件写面（386-r2 刚修的对照件）——`--name-only` 零命中坐实；`testdata/**`、`internal/**`、`cmd/wisp/**` 零命中；本票也未把 386 的四帽判定改动一字（`chat.go:216/218` 现值健在）。
- push：`6cf941b72` 与其后 `79041f071`/`386` 结案件均已在 `origin/dev`——`git branch -r --contains 6cf941b72`＝`origin/dev`，该推送属**编排者权限动作**（`AGENTS.md §④` 推送只由编排者做），⛔ 可据此判落地腿违令；本腿全程零 push。
- 现盘在他腿之后：本件链上又落 `1773881db`（386-r2 已在 §2 基准内）与 `79041f071`（ci-r1，删 `chat.go` 死码 `estimateTokens`，与 391 语义无关）。本腿所有读数**取自现盘 HEAD `250735423`**（含 ci-r1 之后），⛔ 转抄实现腿当时的历史读数。

## 9. 本腿看到的可计缺点／残余（逐条具名，不影响裁决）

1. **三方言共同的一个角没截**（`fc=broken` ＋ forced `tool_choice` ＋ 小 `max_tokens`）：截断门槛看的是 fc 降级**之前**的 `toolName`，降级后正文照旧全文回传。★本腿判定＝**⛔ 本票回归**（chat 同形：`chat.go:215` 门槛亦在 `synthChat` 的 fcBroken 之前，`chat.go:238-244`），三方言齐；若要收紧应是一张覆盖三方言的新票。
2. **票面 AC#2「finish=length」与落地的协议词汇之差**：票面沿用了 386-r2 chat 案形的措辞，落地分两词汇——已在 §7 判合规并给出消费方实锤。建议结案时在台账里把这句口径钉成一枚注（本腿⛔ 动台账/票面），免得后续腿按字面再抄一遍 `length`。
3. 票面机理的 `messages.go:179` 与修前快照实测 `L180` 差一行（同一枚字段声明）；`responses.go:67` 逐字对得上。属薄引用行号腐烂，⛔ 影响判定。
4. **派单原文尺的落点**（根模块 `go test ./tools/mockllm`）对独立子模块必 `setup failed`——本腿双记（§5 前两行）。这是 386-v1 §7-3 已提过的同一枚模板缺陷，第二次复现，建议编排者统一改模板。
5. 同形状但**本票射程外**的余部：`messages.go` 的 `req.System`/`req.Thinking`、`responses.go` 的 `req.Instructions`/`req.Reasoning` 仍只解析不应用（现盘 grep `req\.[A-Z]` 逐枚核到）。本腿**判其非同一病灶**：这些字段的答案由 `/__control/last_request` 捕获后可被消费方断言（`protocol_test.go:50` 即用此道验 `instructions`），且 mock 的确定性回声按设计不以 system 为条件——与 `max_tokens`「请求方期待响应本身发生变化」不同类。建议作为**登记不修**看待，⛔ 据此判 391 未闭合（票面机理只点名 `max_tokens` 两枚，本腿现盘核到那两枚已真消除）。
6. `tools/mockllm/mockllm.exe`（约 10 MB）在写面目录里，已被 `.gitignore:2 *.exe` 挡住、**未入库**（`git ls-files tools/mockllm` 零命中）——形状合规，但建议后续腿把这类产物建到仓外（本票临时件已在 `D:\tmp\wisp\20261102\391v1\`）。

## 10. 是否同意结 391

**同意结票 391。**

五格全 **PASS**：

- **AC#1** 死字段真消除——两文件除解析行外各新增应用点，五个语义位与 `chat.go:215-220` 逐位同形、帽本体（`len/4` 启发式）与 helper 零接触，`chat.go` 本票一字未动；「正确接线」不由词汇表自证，而由消费方两处实锤（`anthropic/adapter.go:285`+`:302-307`、`openairesponses/adapter.go:235-244`+`stream.go:278-280`）与本腿端到端探针读数（C6 `stop="max_tokens"`、399 字节、`output_tokens=100`）三重落定。
- **AC#2** 两方言各一汉字切案（＋各一流式＋各一 ASCII 对照），算术与 386-r2 案形逐位同款，双断（`utf8.ValidString`＋逐字相等）。
- **AC#3** 本腿自建修前快照独立复现（sha 与实现件同值）⇒ 装回死字段 **6/6 子案全红／14 句读数**，再另出四枚定向突变证明**词汇／切点／状态／reason 四维各自单独可红**，撤后复绿，还原核工作树空、两文件 sha 未变。
- **AC#4** `tools/mockllm` 模块内全包 **FAIL=0**（`36/20/0/0` rc=0，与实现件四数吻合）；消费方 `internal/llm/...`＋`internal/present` **7 包全 ok rc=0**；本腿加测的端到端探针另给一层接线证据。派单原文尺的 `setup failed` 属独立子模块落点问题、已双记，⛔ 计为件的红。
- **AC#5** `--stat` 恰 4 件且与题面逐件一致；写面＝`tools/mockllm/**`＋规程指定证据件；票面五框未翻、Status 未动、台账零接触；`cmd/wisp`／`internal/**` 零命中；无未申报顺手修。

§9 六条为口径、模板与射程外观察，⛔ 构成退件理由。结票后的尾巴，供编排者处置（本腿⛔ 动台账/翻框/开票）：第 1 条（三方言共同的 `fc=broken`＋forced tool 角）与第 4 条（派单尺落点，386-v1 §7-3 二次复现）值得登记。
