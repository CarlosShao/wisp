# 票 358 终裁腿 358-v1 - 判语件 part 2/3（非实现者，裁决者 NOT 实现者）

前件＝`docs/evidence/s1/358-v1-verdict.md`（part 1/3：sec 0 档位表、sec 1 前提新鲜度、sec 2 `AC#0`）。

## 3 `AC#2`（落地＋牙）- 成立，可翻（本腿自造三发变异）

NOT 采信实现腿报的 M1-M5。本腿**自己造**三发；每发**先证落地**（`grep -n` ＋`go build ./...` rc=0）
再读数，**每发读完立刻单独还原**并核 `sha256`（快照＝`D:\tmp\wisp\20261021\358v1\m0\`）。
M0 快照 sha256 与实现腿报的**逐字相等**（`run.go` `58e8d8d6...ab3bfb`、355 test `4aa4fe44...114262`、
358 test `939eae29...845939b`），故“还原到哪”两造同尺。

| 变 | 突变（落地读数） | build | 变异读数 | 还原 sha256 |
|---|---|---|---|---|
| **A** 摘掉 attempted 那一行 | `grep -n "超了阈值但一枚都没折"` rc=1；`grep -n "case agent.CompressionAttempted"` rc=1 | **0** | **rc=1**，2 名转红 | `58e8d8d6...ab3bfb` 逐字相等 |
| **B** 让 skipped（含零值）也印 | `:1355 case agent.CompressionSkipped`、`:1356`/`:1358` 两条 print | **0** | **rc=1**，3 名转红 | 同上逐字相等 |
| **C** compacted 那一行改一字 | `:1343` `折入摘要`（`折进`=>`折入`），`grep -c 折进摘要`=0 | **0** | **rc=1**，1 名转红 | 同上逐字相等 |

**变异 A 红句逐字**（摘掉 attempted 那一行）：
```
--- FAIL: TestCompressionTraceSilentWhenNothingCompressed (0.00s)
    compression_trace_355_test.go:214: a pass that cleared the threshold but folded nothing answered "" (ok=false):
      the state is "attempted", so this run DID try and the reading must say it folded nothing -
      silence here is the hole ticket 358 exists for
--- FAIL: TestCompressionTraceAttemptedSaysItFoldedNothing (0.00s)
    compression_trace_358_test.go:172: an attempt that folded nothing answered "" (ok=false)    [两子用例各报一次]
    compression_trace_358_test.go:212: compressionLine answered nothing for state "attempted"
    --- FAIL: .../keep_raw_rounds_floor       --- FAIL: .../provider_triggered_under_fold_budget
```
**变异 B 红句逐字**（让 skipped 也印）：
```
--- FAIL: TestCompressionTraceSilentWhenNothingCompressed (0.00s)
    compression_trace_355_test.go:183: a run that never reached the compressor printed
      "wisp run: compression trace: no pass ran (tokens_before=0)" (ok=true); a line of zeros would claim a measurement nobody took
    compression_trace_355_test.go:252: a report that names the skipped state printed
      "wisp run: compression trace: never triggered (tokens_before=0)" (ok=true); that shape never attempted a fold, so a line there is noise
--- FAIL: TestCompressionTraceReachesStdout (0.00s)
    compression_trace_355_test.go:285: printCompressionTrace claimed a print for a run that never compressed
    compression_trace_355_test.go:288: a run that never compressed wrote "..." ; want zero bytes
--- FAIL: TestCompressionTraceSkippedStaysSilent (0.00s)
    compression_trace_358_test.go:248 / :253 / :256    [两子用例各报]
```
**变异 C 红句逐字**（compacted 改一字）：
```
--- FAIL: TestCompressionTraceCompactedLineIsTheOne355Shipped (0.00s)
    compression_trace_358_test.go:77: the compacted line changed (ticket 355's wording is pinned byte for byte):
    compression_trace_358_test.go:86: what reached stdout is not the pinned line plus one newline:
```
=> 三发都**只剩本票那一条腿在红**（C 只红 358 那枚），说明“改一处、响一处”，NOT 靠别处连带变红。

**特别核（派单点名的独立副本尺）**：`compacted` 那一行是否**逐字等于票 355 原来那一行**。
本腿用**第三把尺**（既不读产码格式串、也不读 358 测试里的副本）：从 `b676092d:cmd/wisp/run.go`
用 `sed` 抽出该字面、`cat -A` 看行尾、`sha256sum` 比对三份：

```
b676092d 原文   936b2aa53c19b4f6b501665578eda90960c4618f631ec61bd7dd80784392166d   135 bytes
HEAD 产码       936b2aa53c19b4f6b501665578eda90960c4618f631ec61bd7dd80784392166d   135 bytes
HEAD 358 测试   936b2aa53c19b4f6b501665578eda90960c4618f631ec61bd7dd80784392166d   135 bytes
```
=> **逐字节相等**（含全角标点与空格）。故票 355 那一行**一个字节都没被 358 动过**，且 358 测试里
的期望值是**独立副本**（若它从产码反推，改产码就不会红；变异 C 证明它会红）。

## 4 `AC#5`（触发基数）- 成立，可翻（附一枚残余面）

### 4.1 真值 / 估算两支各一枚用例，且今天都响

`TestCompressionTriggerPrefersProviderUsage`（真值优先，**两向**都测：估算低/真值高、估算高/真值低）、
`TestCompressionTriggerFallsBackToEstimateOnFirstTurn`（`observed=0` 必落估算，且与 `Need()` **逐值相等**）、
`TestLoopTriggerUsesProviderReportedInputTokens`（**真 `Loop.run`** 打 `tool-then-text.sse`）、
`TestLoopTriggerFallsBackToEstimateWhenNothingObserved`、`TestLoopTriggerBasisCanChangeBetweenTurns`
（同一历史两基数必须给不同读数，防“算一次钉死”）。

**本腿自造变异证这两支都有牙**（NOT 抄它报的 M3）：
- `MT-A`＝`compress.go:269` `if false && observed > 0`（真值支短路）=> **rc=1**，4 名转红：
  `TriggerPrefersProviderUsage`（`:270` Source=estimate、`:273` Tokens=32 want 768、`:276` Need=false）、
  `AttemptedWhenProviderSaysOverButFoldHasNothing`（`:363`）、
  `LoopTriggerUsesProviderReportedInputTokens`（`:73`/`:80`）、
  `LoopTriggerBasisCanChangeBetweenTurns`（`:196`/`:199`/`:204`）。
- `MT-B`＝`:278` `est := 0`（估算支谎报）=> **rc=1**，1 名转红：
  `TriggerFallsBackToEstimateOnFirstTurn`（`:328` Tokens=0 want 32/1236、`:331` Need=false、`:341` 与 `Need()` 不等）。
=> **两支各有独立的牙**，NOT 只有一支有牙、另一支恒真。两发还原后 `compress.go`
`sha256=a23d4287...80a0e1`、`git hash-object=b23910d9...6d37`（＝HEAD），`grep -c MUTATION`=0。

### 4.2 `Source` 一路带到报告与两条痕

现读：`Trigger()` 两支各自 `Source: TriggerProviderUsage` / `TriggerEstimate`（`:274`/`:283`）；
`Compress()` 回填 `rep.TriggerSource`（`:385`）；`compacted` 痕 `:477` 与 `attempted` 痕 `:510` **都**带
`"trigger_source"`＋`"trigger_tokens"`；`run.go:1352-1354` 的 attempted 行印 `source=%s`。
用例按名钉住：`compress_state_358_test.go:171`（12 键含两 trigger 键）、`:209`/`:399-403`（痕里逐值）、
`compress_trigger_358_test.go:72`/`:150`（`Result.Compression` 上逐值）、`compression_trace_355_test.go:233`、
`compression_trace_358_test.go:188`（CLI 行上逐值）。
=> **静默混用没有发生**：拿不到真值时痕里是 `trigger_source=estimate`，NOT 空着也 NOT 假装。

### 4.3 残余面一枚（本腿独立发现，NOT 两条腿报过）

**`source` 为空时保持空、NOT 填 `estimate` - 这属实现选择；本腿判它合判决书，但留一枚残余面。**

`run.go:1347-1350` 注释逐字：“source is printed as the report carries it. An empty one stays empty
rather than being filled with `estimate`: absence is a fact about the report”。
**判决书 sec 4.2 判据 (1) 要求**“真值不可用时回落到估算并**在痕里标明用了哪个**” - 这条**已被走通
路径满足**（两支都填了 `Source`）。本腿自造探针 `zz358v1b_probe_test.go` 量那唯一能造出空 `Source`
的面（调用方传 `TriggerReading{Source:""}`）：

```
State="attempted" TriggerSource="" TriggerTokens=1152 TokensBefore=20
trace: agent: history compress skipped || state=attempted reason=under_fold_budget ... trigger_source= trigger_tokens=1152
FINDING: report carries an INVALID/unnamed TriggerSource "" while State="attempted"
```
=> 报告会出现 `TriggerSource` **不满足自己的 `Valid()`**，且痕里 `trigger_source=` 是空值。
**但**：`loop.go:418` 是**唯一产码调用点**，传的 `trig` 永远来自 `Trigger()`（两支都取名），
故**产码路径不可达**；可达面只有“将来某调用方自己拼 `TriggerReading`”。
`TriggerSource.Valid()` 现无产码调用者（`grep` 只有 `compress.go:130` 定义与 `String()` 内部用）。
=> 判：**今天不构成 `AC#5` 的缺陷**（判据要求“不掉真值、不静默混用”，两者都成立），但**属真实残余**：
`Compress` 对“未命名的基数”既不拒也不补，报告可自称 attempted 却答不出基数。
**建议**（NOT 本腿改码）：写进 `AC#5` 后续收尾，或由 `TriggerSource.Valid()` 在 `Compress` 入口校验一次
=> 不满足即回落 `TriggerEstimate`。**NOT 判 `AC#5` 不可翻** - 三格判据逐条满足，这枚是超出字面的加固面。
# 票 358 终裁腿 358-v1 - 判语件 part 2/2（非实现者，裁决者 NOT 实现者）

前件＝`docs/evidence/s1/358-v1-verdict.md`（part 1/2：sec 0 档位表至 sec 4 `AC#5`）。

## 5 票 355 那颗钉子 - 判：改对了，且确比原来强

**改前全文**＝`git show 5ca2472e:cmd/wisp/compression_trace_355_test.go`（`05f76ad6^` 同 blob
`980c634b...f388`）。那颗钉子改前 `:143-174`，两半：(a) 零值报告必须静默（`:146-149`）；
(b) 过阈值但只 1 轮 =>`Ran=false` 且**必须静默**（`:171-173` 逐字
`if line, ok := compressionLine(rep); ok || line != "" { t.Errorf("a pass that cleared the threshold but folded nothing printed %q (ok=%v)", line, ok) }`）。

**逐条比对四项读数**（改后＝HEAD `20a8777a...`）：

1. **(a) 原句一字未删**：改后 `:182-185` 仍是
   `if line, ok := compressionLine(agent.CompressionReport{}); ok || line != ""` ＋`t.Errorf(...)`，
   只把消息里 `never compressed` 改成 `never reached the compressor`（更准，non-weakening）。
2. **(b) 从“必须静默”改成“必须明说不折”**：原断言**被替换**而非删除 - 换成 `:212-217` `if !ok || line == ""`
   => `t.Fatalf`，加 `:220` 必须含 `超了阈值但一枚都没折`、`:223-227` 必须**不**含
   `没触发`/`未触发`/`跳过`、`:229-238` 四读数逐名带真值。=> 这一半**方向翻转且加强**
   （原只禁“出声”，现要求“说对话”）。
3. **(c) 新增第三半**：`:250-254` 填好的 `skipped` 报告也必须静默（原版没有；防只特判 `{}`）。
4. **另两枚钉子零字节**：`TestCompressionTraceThreeFlavoursReachTheClosingLine`、
   `TestCompressionTraceReachesStdout` 的 diff 里只有它们**上方注释段**被改准（`@@ -20,13 +20,15 @@`
   那一段），函数体零改动。

**本腿自己造反例证“新两半真的钉住了”**（NOT 只看它报的读数）：
- **反例1（让 attempted 静默）**＝变异 A => **必红**，`:214` 逐字转红（sec 3）。=> (b) 真钉住。
- **反例2（让 skipped 印）**＝变异 B => **必红**，`:183`＋`:252` 两条同时转红（sec 3）。=> (a)(c) 真钉住。
- **反例3（比原来更强的反面）**：旧判据只禁“出声”，**一行只说数字不说人话**的 zeros-only 行它**会放行**。
  本腿现量该形：`"...tokens_before=774，tokens_after=774，trigger=774 tokens，source=estimate"`
  => 不含 `超了阈值但一枚都没折` => **今天必被 `:220` 判红**，旧判据（只查 `ok||line!=""`）**会放它过**。
  => 这是“更强”的可核证据：**新判据拒掉了一类旧判据接受的实现**。

**=> 终判：改对了。** 理由三条：(1) 判据方向由“禁”变“禁＋必须说对话”，是**单调加强**；
(2) 原 (a) 一字未删、另两枚钉子零字节 => NOT 删断言换绿；
(3) 改的理由（前提被 358 打掉）**与盘上事实相符** - 本腿已独立确认改前两形在可读面上完全同形（sec 2），
358 后 `CompressionState` 让两形在**类型上**分开（`Valid()` 三值、空串不算态）。
=> 裁决“改既有测试判据”**合铁律**（NOT 因“改测试有风险”而绕路），且该改**属测试判据、NOT 契约文字**
（D 表无此条）。

