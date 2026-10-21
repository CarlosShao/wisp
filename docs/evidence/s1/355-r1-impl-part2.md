# 355-r1 写腿（part 2／2）：新用例、三读数、门禁、未覆盖

接 part 1：`docs/evidence/s1/355-r1-impl.md`（头部、sec 0 选型、sec 1 闸门、sec 2 改了哪几行）。

---

## 3. 新用例逐名 ＋ 为什么今天响

新件＝`cmd/wisp/compression_trace_355_test.go`（本波唯一新 `.go`）。三枚：

| 用例 | 钉什么 | 为什么**今天响**（NOT 恒真） |
|---|---|---|
| `TestCompressionTraceThreeFlavoursReachTheClosingLine` | 一次**真压了历史**的 pass，三味都在收尾行上且与 report 一致 | 驱动**真** `agent.NewCompressor`＋**真** `Compress`（NOT 字面量）；fixture 先自证过了阈值（`c.Need`）再自证真折了（`rep.Ran`）。**摘掉 `:238` 赋值 => `HistoryChanged` 恒 false 而历史明明动了 => 红**（sec 4 逐字） |
| `TestCompressionTraceSilentWhenNothingCompressed` | 反向牙：**没压**就不许印（NOT 印一串 0） | 两形：零值 report；"过阈值但折叠不了"（1 轮 raw 打 `KeepRawRounds` 地板，票 153 点名的真实形状）。pad 由 `b.HistoryCompressTokens*4` 算出，NOT 手写数字。改成"没折也印"即转红 |
| `TestCompressionTraceReachesStdout` | 出口半：走**产码那枚** `printCompressionTrace` 写 stdout，恰好一行、换行结尾、含三味；静默形写 **0 字节** | 断言**生产 writer 路径**（NOT 自拼）。"0 字节"同时堵住"印个空行也算 print" |

**"今天最常见的那一形答什么"**（票面 AC#2 显式要求，NOT 只在退化形上存在却不说）：
今天 `wisp run` 的**常见形答"没有痕"**——压缩只在 (1) 历史过阈值 **且** (2) 可折 raw 轮数 > `KeepRawRounds`(3) 时发生；
CLI 一条任务的 user 轮只来自 `userMessage(input)` 与 C22 提醒（`reminderMessage`，同 `RoleUser`），
常见会话凑不到第 4 枚 user 轮 => `Need()` 或 fold 循环先 break。**这一形本腿如实答"没有痕"（不印行）。**
其预算契约＝`Q-312`、归机主，本腿**零触碰**。

---

## 4. 三读数逐字（牙：变异红 / 还原绿 / 原样绿）

### 4.1 先证变异落地（NOT 空口）
```
$ grep -n "MUTATION 355-r1 M1" internal/agent/compress.go
238:		// MUTATION 355-r1 M1: assignment removed (reverse judgement).
$ grep -n "rep.HistoryChanged = " internal/agent/compress.go
ABSENT (mutation landed)
$ go build ./...   ->  build rc=0
```

### 4.2 变异红（摘掉 `:238` 赋值）— `rc=1`，逐名报红行
```
--- FAIL: TestCompressionTraceThreeFlavoursReachTheClosingLine (0.00s)
    compression_trace_355_test.go:111: HistoryChanged = false for a pass that replaced the history: 1236 -> 780 tokens, 9 messages folded into the summary
    compression_trace_355_test.go:131: the closing line does not carry "history_changed=true": wisp run: 压缩痕：压前 1236 tokens => 压后 780 tokens（history_changed=false，9 条消息被折进摘要，保留 3 轮原文）
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.147s
```
注意：同一发那枚 slog 记录也印 `history_changed=false`（合法，它读同一处赋值），
而 `1236 -> 780 tokens` 说明**历史真的动了** => 这正是"读数与事实矛盾"被抓的形状。

### 4.3 还原绿（按裁决 (3) 核 sha256）
```
$ (Get-FileHash internal/agent/compress.go -Algorithm SHA256).Hash
A4254B1558996AC21E3793E058740677A0504AD9910A4797A728D67B36D599B8   # 与变异前保存值逐字相同
$ grep -n "MUTATION" internal/agent/compress.go   ->  (无输出)
--- PASS: TestCompressionTraceThreeFlavoursReachTheClosingLine (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingCompressed (0.00s)
--- PASS: TestCompressionTraceReachesStdout (0.00s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.135s
```

### 4.4 原样绿（在**未改的 HEAD blob** 上打同一发）
```
$ git show HEAD:internal/agent/compress.go > internal/agent/compress.go
$ grep -c "HistoryChanged" internal/agent/compress.go   ->  0
$ go build ./cmd/wisp/
cmd\wisp\run.go:1323:42: rep.HistoryChanged undefined (type agent.CompressionReport has no field or method HistoryChanged)
```
=> **原样根本编不过**：这一味今天**不存在**，故本用例在 HEAD 上不是"红"而是"没有该读数"——
这正是"今天不响、修了才响"的最强形状（NOT 我在未修码上把它跑绿过）。
另量：`internal/agent` 全包在改动后 `rc=0`、`--- FAIL:` **0 枚**。

---

## 5. 门禁（我自己现跑，贴末行；四数分开报，NOT 相减）

| 门 | 命令 | rc | 末行 |
|---|---|---|---|
| G1 构建 | `GOFLAGS= go build ./...` | **0** | （无输出） |
| G2 vet | `go vet ./cmd/wisp/` | **0** | （无输出；`./internal/agent/` 亦 **0**） |
| G3 整包 | `go test -count=1 -timeout 30m ./cmd/wisp/` | **1** | `FAIL github.com/CarlosShao/wisp/cmd/wisp 325.439s`（**既有红，见 5.1**） |
| G4 D22 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; ... ban #8 cmd/=144`（`cmd/` 143 涨到 **144**＝本腿新增那枚 `.go`；仍 clean） |
| G5 路径长度 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |
| G6 agent 包 | `go test -count=1 -timeout 15m ./internal/agent/` | **0** | `ok github.com/CarlosShao/wisp/internal/agent 3.145s` |

### 5.1 整包红名册（逐名；NOT 只报枚数）
**新增红＝0。既有红恰 1 枚**（起手那次未改码基线里同红，与本腿无关）：

| 红名 | 归属 |
|---|---|
| `TestAC14GoSideEvalPushReachesThePage` | **既有**：`panel_resident_windows_test.go:944/946`，Go 侧 Eval push 没到页面（`title=""`）。WebView2 实况面，与压缩痕无关 |

**被"吃掉"的枚名（派单要求登记）**：整包三发（起手基线／改后／未改的 HEAD 对照）**每次都恰有 1 枚 RUN 无判决行**，且每次落在不同面板用例上：

| 发 | RUN | PASS | FAIL | 无判决 | 死在 |
|---|---|---|---|---|---|
| 起手基线（未改码） | 213 | 211 | 1 | **1** | `TestBookedSnapshotReachesHostBitForBit329` |
| 改后（本腿） | 211 | 209 | 1 | **1** | `TestPanelThreadIsSTAAndExitsCleanly` |
| 对照（未改的 HEAD） | 203 | 202 | 0 | **1** | `TestAC13BringUpSurvivesAReusedThreadQuit` |

三发都死于 **~277-343s**、远早于 `-timeout 30m`，且被吃掉的是**不同**的面板/WebView2 用例 => 是**既有的进程级截断**
（`A1079` 同族"无判决行"形状；NOT 归因给本腿改动）。本腿只新增 3 枚用例，它们在**三次单跑**里都稳绿
（`ok 0.134s / 0.193s / 0.341s`）。

---

## 6. 提交账（逐枚 `git show --name-only`；含正控）

| 提交 | 主题 | 件数 | `.go` 枚 | 越界枚 |
|---|---|---|---|---|
| `b676092d` | `票 355 落地腿 355-r1 笔(1)=压缩痕三味出产码口（H3 收尾摘要行）+ CompressionReport 只读字段` | 5 | **3**（`cmd/wisp/run.go`、`internal/agent/compress.go`、`cmd/wisp/compression_trace_355_test.go`） | **0** |

- 逐枚 `git show --name-only b676092d` 原件：`cmd/wisp/compression_trace_355_test.go`、`cmd/wisp/run.go`、
  `docs/evidence/s1/355-r1-impl-part2.md`、`docs/evidence/s1/355-r1-impl.md`、`internal/agent/compress.go`。
- 越界判据＝白名单五枚之外的件数；读数 **0**。自门：`git ls-tree -r HEAD | grep compression_trace_355_test.go`
  => `100644 blob 980c634b690e3edc100d851e2243192da6ecf388`（件真在盘上、真在树里）。
- **正控（同尺打一笔含 `.go` 的真提交，证明尺不是恒 0）**：`d32b459c` => `.go` **4 枚**；
  同尺打 `abd09226`（纯交件）=> `.go` **0 枚**（能分辨）。故 `.go` 枚数这把尺**两向都可读**。

---

## 7. 未覆盖／未假设

- `Q-312`（`KeepRawRounds`／`Need()` 预算契约＝"常见路径走不到压缩"）**零触碰**，只按票面如实答"这一形没有痕"。
- H6（诊断包）**未接线**：今天零产码调用者，接线＝新增 CLI 动词/出口，超出"把已有的痕变成可答"。**已知残余**，NOT 假装已解。
- 面板/前端：零字节（`SPEC-08` sec 5.2 白名单零触碰）。
- `gofumpt`：本机无 v0.12.0 宿主 binary（`which gofumpt`＝127）=> **该项未跑**，只有 `gofmt -l`（两枚改动文件干净）。
  登记为**未覆盖**，NOT 报成"通过"。既有 `cmd/wisp/models.go` 命中 `gofmt -l`＝**既有**（HEAD blob 同样命中，非本腿引入）。
