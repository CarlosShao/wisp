# 358-r1b 落地件（part 2/2）：attempted 变成 CLI 可读

身份＝写腿（实现者，续跑实例，业务身份同 `358-r1`）。不出终裁、不翻框、不改票名/`Status:`/`-done`/台账、**不 push**。
派单＝`.scratch/wisp/probes/orch/20261021/txt/358-ruling-nail-355.txt`（`5fe1d3be`）｜锚＝`e1de7bab`｜钟点 `2026-10-21 17:2x-18:0x +08`。
前作＝`docs/evidence/s1/358-r1-impl.md`（160 行＝帽上，故本笔拆 `-part2`）。

## 0 开工前

- `git status --porcelain -- cmd/wisp/ internal/agent/ docs/evidence/s1/` 起手＝**空**（写面内零残留）。
- 三枚涉改件**改前 sha256**（`D:\tmp\wisp\20261021\358r1b\pre-change-sha256.txt`）：`d9dfee30...794be3 run.go`｜`f85a9ceb...9233b9 compression_trace_355_test.go`｜`a23d4287...b80a0e1 compress.go`。
- 上一枚结论**复认不重做**：358 三态已在盘（`internal/agent/**` 本笔**零改动**，`git hash-object` = `b23910d9...6d37` = `git rev-parse HEAD:internal/agent/compress.go`，逐字相等）。

## 1 改了哪几行（逐行）

**`cmd/wisp/run.go`（唯一产码面）**

| 行 | 改动 |
|---|---|
| `:1306-1338` | `compressionLine` 文档整段重写：按 **State 三态各说各话**立论，逐条写明 skipped 与未填两形为何**仍静默**（"没试过就不该有行；印一串 0 等于声称一次没人做过的测量"） |
| `:1341-1345` | `case agent.CompressionCompacted:` => **原格式串一字不动**（`压前 %d tokens => 压后 %d tokens（history_changed=%t，%d 条消息被折进摘要，保留 %d 轮原文）`） |
| `:1346-1354` | **新增** `case agent.CompressionAttempted:` => `wisp run: 压缩痕：超了阈值但一枚都没折（tokens_before=%d，tokens_after=%d，trigger=%d tokens，source=%s，保留 %d 轮原文）`；四项读数全带；`source` 按 report 原值印（空就空，NOT 填 `estimate`，那是替报告声称一个没人记过的基数） |
| `:1355-1356` | 其余 => `return "", false`（`CompressionSkipped` 与"压根没跑"的零值报告都落这里） |
| `:1228-1229`／`:1289-1296` | 调用点与 `printCompressionTrace` 注释改成"一态一行、只在**没试过**的形上静默"，行为零变化 |

**`cmd/wisp/compression_trace_355_test.go`**：只按裁决第 5 条改**那一枚钉子** `TestCompressionTraceSilentWhenNothingCompressed` 的判据（`:142-253`），文件头 `:23-30` 那句"silent unless a pass ran"随之改准。**另两枚钉子（`:106` `ThreeFlavoursReachTheClosingLine`／`:261` `ReachesStdout`）零字节**。

**新件 `cmd/wisp/compression_trace_358_test.go`**（260 行）：三形各一枚用例（`:65`／`:132`／`:221`）＋助手 `attemptedAtKeepRawFloor`（`:96`）／`mustLine358`（`:53`）＋ `compactedLine355`（`:43`，355 原文**逐字抄**当期望值，NOT 从产码格式串反推）。

**没动的（逐条对齐派单边界）**：`internal/agent/**` **零字节**；`budgets.go`／`docs/PLAN.md`／`docs/specs/**` 零字节；`internal/{tools,config,winsec,risk}/**`、`tools/d22scan/**`、golden、`thresholds.go`、`frontend/**`、`design/**` 零字节；未造前端事件／`panel.*`／JSON 键。`Ran` 镜像语义**未动**（`rep.Ran = State==Compacted` 是 358 已交的产码，本笔只读 `State`，11 处旧读者含义一个没改）。

## 2 三形用例逐名 ＋ 为什么今天响

| 用例 | 钉什么 | 为什么**今天响**（NOT 恒真／NOT 只印零） |
|---|---|---|
| `...CompactedLineIsTheOne355Shipped` | `compacted` 那一行**逐字不变**：真压缩机真折叠的 report 过**产码** formatter，整行 `!=` 期望即红；再走 `printCompressionTrace` 比"整行＋恰好一个换行" | 期望串是 355 产码原文**抄进测试**的独立副本 => 改动/换序/换字段都红。变异 M3（"折进"=>"折入"）即红（sec 3.1） |
| `...AttemptedSaysItFoldedNothing` | `attempted` **必须明说不折**：(1)必须出话 (2)必须含"超了阈值但一枚都没折" (3)**不许**含"没触发/未触发/跳过" (4)四项读数逐名带真值 (5)隐私哨兵 (6)出口字节＝行＋换行。**两形都驱动**：`keep_raw_rounds_floor`（过阈值打地板）＋`provider_triggered_under_fold_budget`（provider 说超、估算没超，`tokens_before==tokens_after<阈值`） | 改前对非 `Ran` 一律 `return "", false` => (1) 直接 `t.Fatalf`，红。**两形必须给出不同的行**（末尾比对）=> "回一句固定话／只印零"都过不了：第二形两数相等，只印零会被读成"没触发" |
| `...SkippedStaysSilent` | `skipped` **仍静默**：(1)`ok=false` 且 `line==""` (2)`printCompressionTrace` 返 false (3)**写 0 字节**（NOT 空行）。**两形都驱动**：`filled_skipped_report`（真压缩机跑过、`State=skipped`）＋`zero_value_no_pass_ran`（零值报告，`State` 是空串） | 这一半是**反向牙**：把 skipped 也印出来即红（变异 M2，sec 3.1）。零值那格同时堵"只特判空 struct"的假绿 |
| （钉子，改判据）`TestCompressionTraceSilentWhenNothingCompressed` | **两半都钉**：(a) 零值报告静默；(b) `attempted` **必须明说不折**＋四项读数＋不许说"没触发"；(c) **新增第三半**：填好的 `skipped` 报告也必须静默（防只特判 `{}`） | 比 355 原判据**更强**：原来只有 silent 一侧，现在 silent 与"必须出声"两侧都在。NOT 删断言换绿，`(a)` 原句一字未删，另加 `(b)(c)` |

**改理由（逐字进测试注释，引裁决）**：原钉子论据是"`ok=false` 同时覆盖'从没压缩'与'不是压缩的那轮'，两形不可分 => 印一行就是在声称一次没人做过的测量"（`run.go` 原注释逐字）。**那论据当时是对的**，但它依赖的前提（两形在报告上**完全同形**）已被 358 的 `CompressionState` 打掉（`skipped`／`attempted`／`compacted` 是一个字段的三个值）。继续守它就等于**用一颗过时的钉子挡住新读数**（台账 `A201` 第二点那族"规格与仪器分叉"）。裁决全文＝`.scratch/wisp/probes/orch/20261021/txt/358-ruling-nail-355.txt`，测试注释按名引了。

## 3 三读数逐字（变异红／还原绿／原样绿）＋ `sha256`

### 3.1 变异红（逐名报红行）
每变异**先证落地**（`grep -n` 原文＋`go build ./...` rc=0），三条落地读数：
```
M1  grep -n "MUTATION 358-r1b M1" cmd/wisp/run.go  =>  1347:  // MUTATION 358-r1b M1: the attempted line removed (reverse judgement).
M2  grep -n "MUTATION 358-r1b M2" cmd/wisp/run.go  =>  1356:  return fmt.Sprintf("wisp run: 压缩痕：没触发（tokens_before=%d）", rep.TokensBefore), true // MUTATION 358-r1b M2
M3  grep -n "压缩痕：压前" cmd/wisp/run.go          =>  1343:  （"折进摘要"=>"折入摘要"，格式串其余逐字不动）
三次 go build ./... 均 rc=0
```
**M1＝把 attempted 那一行摘掉（派单第 2 条那颗反向判据）**：`rc=1`
```
--- FAIL: TestCompressionTraceSilentWhenNothingCompressed (0.00s)
    compression_trace_355_test.go:214: a pass that cleared the threshold but folded nothing answered "" (ok=false): the state is "attempted", so this run DID try and the reading must say it folded nothing - silence here is the hole ticket 358 exists for
--- FAIL: TestCompressionTraceAttemptedSaysItFoldedNothing (0.00s)
    compression_trace_358_test.go:172: an attempt that folded nothing answered "" (ok=false); state="attempted" means the run DID try, so silence here is the hole ticket 358 exists for
    （同句在 keep_raw_rounds_floor 与 provider_triggered_under_fold_budget 两子用例各报一次）
    compression_trace_358_test.go:212: compressionLine answered nothing for state "attempted"
```
**M2＝让 skipped 也印（反向：把静默那半去掉）**：`rc=1`，三枚同时红
```
--- FAIL: TestCompressionTraceSilentWhenNothingCompressed
    compression_trace_355_test.go:183: a run that never reached the compressor printed "wisp run: 压缩痕：没触发（tokens_before=0）" (ok=true); a line of zeros would claim a measurement nobody took
--- FAIL: TestCompressionTraceReachesStdout
    compression_trace_355_test.go:285: printCompressionTrace claimed a print for a run that never compressed
--- FAIL: TestCompressionTraceSkippedStaysSilent
    compression_trace_358_test.go:248: a report (state="skipped") that attempted no fold printed ... (ok=true); that shape never tried, so a line there is noise
```

**M3＝把 compacted 那一行改一个字**：`rc=1`
```
--- FAIL: TestCompressionTraceCompactedLineIsTheOne355Shipped (0.00s)
    compression_trace_358_test.go:77: the compacted line changed (ticket 355's wording is pinned byte for byte):
    compression_trace_358_test.go:86: what reached stdout is not the pinned line plus one newline:
```

### 3.2 还原绿 ＋ 还原核 `sha256`（逐字）
```
还原后 run.go    58e8d8d6d173a7cf99d73b005de95f3826ed142387a25a5f1e9ceaf346ab3bfb == M0 快照逐字相等
还原后 355 test  4aa4fe448450ec1e381be0f4feac2c7a44fc8aaedcc38be584f9650e9f114262 == M0 快照逐字相等
还原后 358 test  939eae292b9d37e9900328dcae517c42de1473a54119b8b6ab0cbddc0845939b == M0 快照逐字相等
grep -c MUTATION cmd/wisp/run.go => 0；还原绿 family_rc=0，六枚顶层 PASS、零 FAIL（逐字见 sec 4）
```
**"回 HEAD blob"为何换成快照 sha**：本票要改的 `run.go` 与 `compression_trace_355_test.go` **与 HEAD 必然不等**（差＝本票改动），故按上一枚已认的口径测**另一枚本票未改的文件**：`compress.go` `git hash-object` = `b23910d9b9390d73675461b112060d4c0a7e6d37` = `git rev-parse HEAD:internal/agent/compress.go`，**逐字相等**。

### 3.3 原样绿（未改的 HEAD blob 上打同一发，`-overlay`，工作树零触碰）
```
overlay：cmd/wisp/run.go => git show HEAD:cmd/wisp/run.go（blob 4865cfde...d3ae，与 HEAD 逐字相等）
grep -c "超了阈值但一枚都没折" <HEAD blob> => 0（那一行今天不在）；同一发 rc=1
--- FAIL: TestCompressionTraceSilentWhenNothingCompressed
    compression_trace_355_test.go:214: ...answered "" (ok=false): the state is "attempted"...
--- FAIL: TestCompressionTraceAttemptedSaysItFoldedNothing
    compression_trace_358_test.go:172: ...answered "" (ok=false); state="attempted"...
```
=> "**今天不响、改完才响**"：同一断言在 HEAD 上红、本树上绿；NOT 我在未改码上把它跑绿过。（`CompactedLineIsTheOne355Shipped` 与 `SkippedStaysSilent` 在 HEAD 上也**绿**，符合预期：本笔**只**改 attempted 那一跳 => 红只落在 attempted 相关两枚，正是"改一处、响一处"。）

## 4 同发跑全族（四数分开，NOT 相减）

`go test -count=1 ./cmd/wisp/ -run 'CompressionTrace' -v` => **rc=0**；顶层 RUN **6** / 子测试 RUN **4** / 顶层 PASS **6** / 子测试 PASS **4** / 顶层 FAIL **0** / 顶层 SKIP **0**
逐名（355 三枚与 358 三枚**同发全绿**，证明 355 其余钉子没被碰坏）：
```
--- PASS: TestCompressionTraceThreeFlavoursReachTheClosingLine (0.00s)      <- 355
--- PASS: TestCompressionTraceSilentWhenNothingCompressed (0.00s)           <- 355（本笔按裁决改判据）
--- PASS: TestCompressionTraceReachesStdout (0.00s)                        <- 355（零字节）
--- PASS: TestCompressionTraceCompactedLineIsTheOne355Shipped (0.00s)      <- 358
--- PASS: TestCompressionTraceAttemptedSaysItFoldedNothing (0.00s)         <- 358
    --- PASS: .../keep_raw_rounds_floor (0.00s)
    --- PASS: .../provider_triggered_under_fold_budget (0.00s)
--- PASS: TestCompressionTraceSkippedStaysSilent (0.00s)                   <- 358
    --- PASS: .../filled_skipped_report (0.00s)
    --- PASS: .../zero_value_no_pass_ran (0.00s)
ok  github.com/CarlosShao/wisp/cmd/wisp  0.285s
```

## 5 门禁 rc（自己现跑，贴末行）

| 门 | 命令 | rc | 末行/读数 |
|---|---|---|---|
| G1 build | `GOFLAGS= go build ./...` | **0** | （无输出） |
| G2 vet | `go vet ./cmd/wisp/ ./internal/agent/` | **0** | （无输出） |
| G3 agent | `go test -count=1 ./internal/agent/` | **0** | `ok github.com/CarlosShao/wisp/internal/agent 4.068s` |
| G4 cmd/wisp | `go test -count=1 -timeout 30m -v ./cmd/wisp/` | **1** | `FAIL github.com/CarlosShao/wisp/cmd/wisp 333.647s`；顶层 **RUN 159 / PASS 158 / FAIL 0 / SKIP 0**，**整进程死在 `TestAC14GoSideEvalPushReachesThePage`（无终态行）** |
| G5 d22scan | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; ... ban #8 cmd/=145`（144=>**145**＝本笔新增那枚 `.go`，仍 clean） |
| G6 路径长度 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name ...`；`over the hat=60 ... roster entries=60` |
| G7 gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l <3 件>` | **0** | 列表**空** |

gofumpt 用**绝对路径**（裸名 rc=127，NOT 据此判"本机没有"）；未用 `CGO_ENABLED=0`；`cmd/wisp` 整包**抬 `-timeout 30m`**。

**G4 的 rc=1 NOT 本笔的伤（按名作差，两向皆空）**：同一把尺、带 `-v`、**按名**比对本笔起手基线（`baseline-cmdwisp-v.txt`）与改后（`post-cmdwisp-v.txt`）：
```
NEW RED (本笔后红、基线不红):        (空)
RECOVERED (基线红、本笔后绿):        TestAC4FocusReturnToPriorWindowGap33r5   <- 台账 :14768 在册真窗焦点仪，颜色随会话漂
ABSENT ON POST (基线有名、本笔后无终态): (空)
NEW ON POST (本笔后才有名):          本笔三枚 TestCompressionTrace*(Attempted/CompactedLine/Skipped)（全 PASS）
  ＋ TestAC13BringUpRefusesAThreadWithAQueuedClose / TestAC13BringUpSurvivesAReusedThreadQuit /
    TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever / TestAC14AwaitedBindingReplyReachesThePage（全 PASS）
```
=> **新增红 0 枚**。四个 `TestAC13/AC14*` 是**上一发死在它们之前、这一发跑到了它们**，列在"new on post"是因基线那发被死点吃掉、**没有它们的终态行**，不是新增红。
**两发死点不同名**（基线＝`TestAC13BringUpSurvivesAReusedThreadQuit`；本笔后＝`TestAC14GoSideEvalPushReachesThePage`）=> 印证"该族死点随会话漂"（`A1079`／`A1100`／票 314），NOT 拿一发读数指控另一发。
整包 **362 枚编译的顶层用例里只有 159 枚跑到**（`grep -cE '^func Test'` = 362 vs RUN = 159），未跑到的那批**本笔无法用整包这一发判色**，**属票 314 射程**。

## 6 提交账（逐枚 `git show --name-only`）

| 提交 | 主题 | 件数 | `.go` 枚 | 越界枚 |
|---|---|---|---|---|
| `43f44f59a99f658a60ecad12c3574589c7cba7f0`（`43f44f59`） | `票 358 落地腿 358-r1b：attempted 变成 CLI 可读（三态各说各话）` | 3 | **3**（`cmd/wisp/run.go`、`cmd/wisp/compression_trace_355_test.go`、`cmd/wisp/compression_trace_358_test.go`） | **0** |
| `a8ae09cebe2d5781258fe0c5d0783e7541ae0c2a`（`a8ae09ce`） | `票 358 落地腿 358-r1b：交付件 358-r1-impl-part2` | 1（本件首次落盘） | **0** | **0** |
| 本件其后各笔（压到帽内、补齐提交号） | 主题均以 `票 358 落地腿 358-r1b：` 开头 | 各 1（本件） | **0** | **0** |

- 本表不写自己那一笔的 sha（写了就立刻过期）。现跑尺：`git log --format='%h %s' e1de7bab..HEAD -- docs/evidence/s1/358-r1-impl-part2.md` 得本件全链；`.go` 逐枚数 `git show --name-only --format='' <sha> | grep -c '\.go$'`。
- 越界判据＝派单白名单（`cmd/wisp/run.go` 的 `compressionLine`／`printCompressionTrace` 一带、`cmd/wisp/compression_trace_355_test.go`、`cmd/wisp/compression_trace_358_*.go`、`docs/evidence/s1/358-r1-impl.md`／新 `-part2.md`）之外的件数；逐笔合并跑 `grep -vE` 逐行输出**为空** => 读数 **0**（`43f44f59` 那笔实测：3 件全在白名单内）。
- 新件唯一性：`git ls-tree -r --name-only HEAD | grep -c '<新件名>'` => `compression_trace_358_test.go` = **1**、`358-r1-impl-part2.md` = **1**。
- **正控**（同尺打一笔含 `.go` 的真提交，证明尺不是恒 0）＝`b676092d`（票 355 落地腿）=> `.go` **3** 枚；**负控**＝`3ed3fdce`（纯交件）=> `.go` **0** 枚。两向都可读，故"`.go` 3 / 越界 0"不是瞎尺读数。
- 未 push（各笔都在本地 `dev` 上）；只有 code 那一笔带 `.go`，交付件各笔零 `.go`。
## 7 未做／未假设

- NOT 翻任何票面框、NOT 改票名/`Status:`/`-done`/台账、NOT push。
- `Q-312`（`KeepRawRounds`／`Need()` 预算契约）**零触碰**；本笔只让 attempted 这一形**可读**，NOT 改它何时发生。
- 前端/面板零字节（`SPEC-08` 5.2 白名单零触碰）；未造事件、未加 `panel.*`、未加 JSON 键。`internal/agent/**` 零字节（358 产码面已交，本笔只读 `State`）。**未定义即停**：未遇到需要自行假设的分叉 => 无停手项。
