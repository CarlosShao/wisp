# 票 386 验收腿 386-v1 · 非实现者裁决表（mockllm 同类帽登记未修）

- 腿：`386-v1`（验收腿＝非实现者，裁决者≠实现者，`SPEC-12 §4.3` #1/#3）。日期 2026-11-02。
- 被验件：commit `1773881db9d9f760a774d222a1c8f6745f57191c`「票 386 落地腿 386-r2」（作者 CarlosShao，2026-10-31 21:42 +08）。
- 票面：`.scratch/wisp/issues/386-mockllm-caps-registered-unfixed.md`（13 行，5 格，**本腿未翻框、未改一字**）。
- 本腿角色：只读复核＋本证据件。全程零产码改动、零 `cmd/wisp`/`internal/**` 写面、零票面/台账、零 push。
- 临时件：`D:\tmp\wisp\20261102\386v1\**`（仓外，只建不删）：`orig/`（三文件修前快照，本腿自取，非沿用实现腿件）、`overlay.json`、`ac3-overlay-red.txt`、`ac3-restored-green.txt`、`full-v.txt`、`present.txt`、`llm.txt`、`d22scan.txt`。

## 0. 起手三核（现跑）

| 核 | 判据 | 读数 |
|---|---|---|
| HEAD 现盘最新 | ≥ `fd8cece8f` | HEAD＝`3a1c03460`（分支 `dev`）；`git merge-base --is-ancestor fd8cece8f HEAD` rc=0 ⇒ **`fd8cece8f` 确为 HEAD 祖先**（`3a1c03460`＝编排者 A1411「派 386-v1」本体） |
| 索引干净 | 暂存区空 | `git diff --cached --name-only` **空**；tracked 脏件仅 `internal/risk/pathresolver_budget_norace_test.go`（他人写面、与本票 `tools/mockllm` 零交），余脏全在未 tracked 的 `.scratch/**`/`design/**`/`logs/**` |
| 写面无他人脏 | `git status --porcelain -- tools/mockllm` | **空**（本腿全程亦未使其变脏，见 §5 还原核后复核） |

## 1. 五格逐条裁决（票面原文逐字＋判定＋本腿自跑读数）

| 格 | 票面原文（逐字） | 裁决 | 本腿读数（⛔ 转抄，全部现跑） |
|---|---|---|---|
| **AC#1** | AC#1 逐帽判定：★每帽"有意模拟 vs 无意缺陷"具名（★对照真实 chat API 的 max_tokens/reasoning 行为；★有意模拟的帽 ⛔ 动、如实登记） | **PASS** | 四帽全具名、判定与代码现值逐条对得上（§2）；「有意」那半边（max_tokens 帽本体＋`finish=length`）在 `chat.go:216/218` 逐字健在、`ascii` 对照案钉死未被顺手拆 |
| **AC#2** | AC#2 多字节回归：★判"无意"的帽按 rune/去帽修＋多字节案一发（★对照 385-r1 的 402 汉字案形） | **PASS** | 新件 `caps_multibyte_386r2_test.go` 两测试七子案（§3）；样本 `ticket30Sample()`＝`测`×402 形，与 385-r1 案形同族；四方言口全覆 |
| **AC#3** | AC#3 反向判据（★三读数＋还原核） | **PASS** | 本腿自建修前快照＋`-overlay` 独立复跑：装回旧帽 ⇒ **rc=1**（红子案 6 / 红读数 7 句）；撤 overlay ⇒ **rc=0** 七子案全绿；ASCII 对照案两轮**恒绿**（咬的是帽，不是别处）（§4） |
| **AC#4** | AC#4 零回归（★mockllm 全包＋消费方 E2E） | **PASS** | `go test -count=1 -timeout 10m ./tools/mockllm`＝**FAIL=0**；包内 `-v` 四数 `RUN=36 / PASS=36 / FAIL=0 / SKIP=0`（顶层 `--- PASS` 20）；`internal/present`（票 30 E2E，真 mockllm 子进程）`ok 7.670s` rc=0；`internal/llm/{openaichat,anthropic,openairesponses}` 三方言 `ok` rc=0（§5） |
| **AC#5** | AC#5 纪律格：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**` | **PASS** | `git show --stat 1773881db`＝恰 5 件、主题逐字「票 386 落地腿 386-r2」与派单一致；产码 4 件全在 `tools/mockllm/**`；`--name-only` grep `^\.scratch/\|^docs/reports/` **空**＝票面与台账零接触（§6） |

## 2. AC#1 帽登记 ⇄ 实现 逐条对账（本票核心）

登记源＝`docs/evidence/s1/385r1-impl.md` §7.1/§7.2（本腿现读），点名四帽；票面 机理 逐字复用同一批行号。对账如下（修前号＝本腿从 `1773881db^` 取的快照实测行号，与票面所引逐枚一致）：

| 帽 | 登记处（票面引用） | 修前实码（本腿现取快照核对） | 判定（实现腿） | 代码现值（HEAD 盘上，本腿现跑 grep） | 对账结论 |
|---|---|---|---|---|---|
| A chat reasoning 200 帽 | `chat.go:168-169` | 快照 L168 `if len(reasoning) > 200 {`／L169 `reasoning = reasoning[:200]` **确实存在** | **无意缺陷→摘** | 该二行**已消失**；L173 只剩 `reasoning = "thinking about: " + lastUser.text`＋L168-172 具名注释；全包 grep `> 200`／`[:200]` 于产码 **零命中** | 落实 |
| B responses reason 200 帽 | `responses.go:100-101` | 快照 L100-101 同型 **确实存在** | **无意缺陷→摘** | 已摘；`responses.go:103` 全量赋值＋L99-102 注释 | 落实 |
| C messages reason 200 帽 | `messages.go:114-115` | 快照 L114-115 同型 **确实存在** | **无意缺陷→摘** | 已摘；`messages.go:117` 全量赋值＋L113-116 注释 | 落实 |
| D chat max_tokens 字节切 | `chat.go:207`（机理段另引 `205-209`） | 快照 L207 `answer = answer[:req.MaxTokens*4]` | **拆两段**：帽本体＝有意（⛔ 未动、如实登记）；切点＝无意（rune 化） | `chat.go:216` 启发式 `len(answer)/4 > req.MaxTokens` 与 `218` `finish = "length"` **逐字健在**；仅切点改 `217` `answer[:runeBoundaryWithin(req.MaxTokens*4, answer)]`＋`chat.go:365` 新 helper | 落实（"有意 ⛔ 动"这一半未被偷拆，见 §3 对照案） |

判定质量复核（本腿独立看，非复述）：

- **A/B/C 判"无意"站得住**：三方言同型对称（385-r1 §2 已考出处＝建仓提交 `0839407bf` 的「合成封顶」惯例），200 这个数字在真实协议里无对应物；真护栏（`io.LimitReader(r.Body, 8<<20)`）另在；385 已判同惯例的回声半边为缺陷，此处是同一判定的余部——**一致性成立**，⛔「同一惯例两半两判」那种自相矛盾。
- **D 判"两段拆"是本票最有价值的一格**：票面自己标注"★这个帽可能是有意的、修前先判"，实现腿没有一刀切摘掉（那会改模拟行为），也没有一刀切不修（那留着坏 UTF-8 出口），而是**把「帽本体」与「切点」拆开**——帽本体＝真实 chat API 命顶截断＋`finish_reason=length` 语义，保留；切点＝按字节切必出非法 UTF-8，真实 provider 在 token 边界切、永不出坏序列，rune 化。这个界线划得对，且 ASCII 行为逐字节不变有对照案钉死。
- **D 的额外反证成立**（本腿 §4 读数①亲眼见）：旧字节切在汉字样本下产出 **402 字节 > 它自己声明的 400 启发式界**（坏尾巴经 JSON 编码展开为 U+FFFD 反涨），连帽自己的语义都被切点破坏 ⇒ "切点非有意"不只是审美判断，是有硬读数的。
- ⚠ **如实记一处口径**：`386r2-impl.md` §2 声明两站官方文档 DNS 被本机环境拒，对照改以「仓内 385 先例＋代码内证＋公开 API 语义通识」三源立据（另附一枚 vLLM 侧二手链接）。这是**诚实披露**而非隐瞒，且判定本身可由码面与读数独立支撑，**不影响 AC#1 判 PASS**；但记下：若后续要对外引用"对照真实 provider"，这步欠一手来源，值得单独补一次联网核。

## 3. AC#2 回归件形状核对

`tools/mockllm/caps_multibyte_386r2_test.go`（284 行，本腿逐行读）：

- `Test386ReasoningDeliversLongMultibyteInFull`：want＝`"thinking about: [think] " + 样本`，先自钉 `len(want) > 200`（防样本退化后测试变哑）；四口 `chat_stream_sse`（逐 delta 装配）／`chat_nonstream_json`／`messages_thinking_block`／`responses_reasoning_summary`，每口 **`utf8.ValidString` ＋ 逐字相等**双断（只测有效性会放过截短，只测相等会放过坏码——两把都上，够）。
- `Test386MaxTokensTruncationIsRuneSafe`：300 汉字＋`max_tokens=100` ⇒ 断言 `finish_reason=length`（**护住有意那半边**）＋内容＝rune 边界前缀 `"echo: " + 测×131`（399 字节）＋`utf8.ValidString`＋`≤ maxTokens*4`；三口＝非流式／流式／**`ascii_cut_byte_identical_to_cap` 对照案**（ASCII 下仍在 200 字节＝50×4 精确截，钉"帽没被顺手拆"）。
- 与 385-r1 案形的对照关系成立：同 `ticket30Sample()` 402 汉字形、同「全量回传」断法，非另起炉灶。

## 4. AC#3 反向判据（本腿独立重做，⛔ 转抄实现腿临时件）

本腿从 `1773881db^` 自取三文件修前快照，写 `Replace` 形 overlay（Go overlay 语义＝只读注入，工作树零接触），在 `tools/mockllm` 模块内跑：

1. **读数①装回旧帽**：`go test -count=1 -timeout 10m -overlay <ov> -run Test386 -v .` ⇒ **rc=1**，`FAIL github.com/CarlosShao/wisp/tools/mockllm 0.056s`。
   - 红子案 **6** 枚：`chat_stream_sse`／`chat_nonstream_json`／`messages_thinking_block`／`responses_reasoning_summary`／`hanzi_cut_lands_on_rune_boundary`／`stream_finish_reason_length`。
   - 红读数 **7** 句（`hanzi` 子案发 2 句：402 字节 ＋ 越自身 400 界），逐字：`assembled reasoning = 204 bytes / 84 runes, want 1261 bytes / 437 runes`（尾可见 U+FFFD 残迹，四口同读数）；`content = 402 bytes / 138 runes, want the rune-boundary prefix 399 bytes / 137 runes`；`content = 402 bytes, exceeds the max_tokens*4 heuristic bound 400`；`assembled text = 402 …`。
   - `ascii_cut_byte_identical_to_cap` 此轮 **PASS**（旧码 ASCII 下本就正确 ⇒ 反证测试没把无关行为一起咬红）。
   - 口径注：实现腿题面写"红 7 案"，本腿实测口径为 **6 红子案／7 红读数句**（差在 `hanzi` 一子案发两句）。**同位同构、无一处对不上**，仅计数口径之差，记为如实读数差 ⛔ 判 FAIL。
2. **读数②撤 overlay**：同命令去 overlay ⇒ **rc=0**，7 子案全 `--- PASS`，`ok 0.059s`（= §5 全包读数的同轮）。
3. **还原核**：overlay 轮后 `git status --porcelain -- tools/mockllm` **空**、四文件 sha/内容＝HEAD 逐字节未变（`chat.go:216/218` 有意那半仍在、`> 200` 零命中）⇒ 突变没落到盘。

⇒ 三读数齐（红→绿→装回复红同位）＋还原核过。**测试咬的是帽。**

## 5. AC#4 零回归读数

| 项 | 本腿命令 | rc | 读数 |
|---|---|---|---|
| 派单点名尺 | 根模块内 `go test -count=1 -timeout 10m ./tools/mockllm` | **1** | ⚠ **非测试失败**：`FAIL ./tools/mockllm [setup failed] … main module (github.com/CarlosShao/wisp) does not contain package …/tools/mockllm`——`tools/mockllm` 是**独立模块**（自带 `go.mod`＝`module github.com/CarlosShao/wisp/tools/mockllm`），在根模块按相对路径点它是加载期错误。**按派单原文照跑必红，属尺的落点问题、不是件的问题**；改在该模块目录内跑 ⇒ **FAIL=0** |
| 正确落点 | `cd tools/mockllm; go test -count=1 -timeout 10m ./...` | **0** | `ok github.com/CarlosShao/wisp/tools/mockllm 0.207s` |
| 包内四数 | 同上 `-v` | 0 | `RUN=36 / PASS=36 / FAIL=0 / SKIP=0`（含子测试；顶层 `--- PASS` 20） |
| 票 30 E2E | `go test ./internal/present/` | **0** | `ok …/internal/present 7.670s`（真 mockllm 子进程走真管线） |
| 三方言消费方 | `go test ./internal/llm/openaichat/ ./internal/llm/anthropic/ ./internal/llm/openairesponses/` | **0** | 三口全 `ok`（28.5／21.0／25.9s）——reasoning 摘帽后消费方无一被长度断言绊倒，与「无消费方测量该帽」的判定互证 |
| 尺三枚 | `gofumpt -l .`（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe` v0.12.0）／`gofmt -l .`／`go vet .` | 0/0/0 | 全 **0 行 / 0** |
| d22scan | `cd tools/d22scan; go run . -root ../../` | **0** | `clean`；live scope `internal/=288`、`cmd/=52`、`frontend/=116`、`internal/tools/=25`、`design/=39`、`ban#8 internal/=740`、`cmd/=186`（比 386-r2 在案 729/180 高，系其后 `391-r1`/`ci-r1` 落的件所致，与本票无涉）；`tools/**` 顶层仍在八枚 ban 射程外，新测试件 `测` 不落任何 emoji 码位带 |
| 工具链陷阱复核 | `go list -f "{{.TestGoFiles}}" .` | 0 | 六件中 **`caps_multibyte_386r2_test.go` 在册**、`chat_echo_385_test.go` 在册 ⇒ 实现腿 §3 报的「尾段 `_386_` 撞 GOARCH=386 被静默排除、改名即现」成立且已解；该提醒对后续腿有值 |

四数差一枚如实报：实现腿 §5 报 `28/18/0/0`，本腿现跑 `36/36/0/0`。差因＝**其后 `6cf941b72`「票 391 落地腿 391-r1」新增 `caps_maxtokens_391r1_test.go`**，⛔ 386 件的回退；两读数的共同点是 **FAIL=0**，尺要的就是这个。

## 6. AC#5 纪律格

- `git show --stat 1773881db`＝**恰 5 件 / +420 −11**：`tools/mockllm/{chat.go,messages.go,responses.go,caps_multibyte_386r2_test.go}` ＋ `docs/evidence/s1/386r2-impl.md`。
- 主题逐字「票 386 落地腿 386-r2」＝派单所引一致。
- 票面 写面＝`tools/mockllm/**`：四枚产码/测试件全落内；`docs/evidence/s1/386r2-impl.md` 在写面外一格——**判合规**：`AGENTS.md §1.5` 与 385-r1 同型先例均把落地腿自身证据件落 `docs/evidence/s1/`，属规程指定落点，⛔ 越面写码。
- `--name-only` grep `^\.scratch/`／`^docs/reports/` **空** ⇒ 票面（含 AC 框）、台账零接触；票面 3 行 Status 至今仍是「待派」（结案归编排者，本腿亦⛔ 翻框）。
- push：`1773881db` 与其后 `391-r1`/`ci-r1` 均只在 `dev` 本地链上；本腿全程零 push。
- 越票顺手修核对：`386r2-impl.md` §7 把「messages/responses 两方言 `max_tokens` 死字段」与「`cmd/wisp/subagent_carrier_197_test.go` 过期注释」两条**登记不修**，⛔ 在本票内动——本腿现盘核到该残余之一已由独立票 `391-r1` 落掉（`messages.go:204`／`responses.go:137` 各自应用帽＋用本协议自己的终止拼写 `stop_reason=max_tokens`／`status=incomplete`），**且没有把 386 的四帽判定改动一字**：登记链⇄开票链⇄落地链闭合，这是纪律格想看到的形状。

## 7. 本腿看到的可计缺点（不影响裁决，逐条具名）

1. 反向判据计数口径：题面/commit message 的「红 7 案」＝红读数句，本腿子案口径为 6（`hanzi` 一子案发两句）。同位同构，建议后续票统一写成「N 子案 / M 读数句」。
2. "对照真实 chat API" 欠一手来源（官方文档页没抓到，改以先例＋内证＋通识立据）。判定可由码面与读数独立支撑，故不扣 AC#1；对外引用前应补一次联网核。
3. 派单原文那把尺（在根模块 `go test ./tools/mockllm`）对独立子模块必 `setup failed`。本腿按落点更正后跑，读数已在 §5 双记。后续派单模板建议对 `tools/mockllm` 写「在其模块内跑」。

## 8. 是否同意结 386

**同意结票 386。**

五格全 **PASS**（AC#1 逐帽四枚判定与代码现值逐条对账落实、有意那半边（max_tokens 帽本体＋`finish=length`）健在且被 ASCII 对照案钉死；AC#2 回归件 402 汉字形四方言口＋双断、与 385-r1 案形同族；AC#3 本腿自建修前快照独立复跑三读数＋还原核全齐、测试咬帽；AC#4 `tools/mockllm` **FAIL=0** 与消费方 E2E 三口＋三方言全 rc=0；AC#5 `git show --stat` 与题面一致、写面仅 `tools/mockllm/**`＋规程指定证据件、票面台账零接触、零 push）。§7 三条为口径与派单尺的改进项，⛔ 构成退件理由。

结票后的登记链尾巴，供编排者处置（本腿⛔ 动台账/开票）：`386r2-impl.md` §7.1 已被票 391 落掉；§7.2（`cmd/wisp/subagent_carrier_197_test.go:605-608` 过期论证注释）仍挂 `cmd/wisp` 所有者名下，385-r1 §7.3 亦在册。
