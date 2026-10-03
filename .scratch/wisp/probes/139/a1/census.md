# 139-a1 普查证据件 — 「上下文压缩不留痕」这一格在现树上的真实密度

派单代号 `139-a1`（只读普查腿）。票面：
`.scratch/wisp/issues/139-context-compression-leaves-no-trace-zero-log-calls-in-compress-go-zero-readers-of-res-compression.md`。
本篇**只做一件事**：把这一格在**当前工作树**上的读数取到编排者能当场决定派哪一枚落地腿的密度。

---

## 0. 锚点与口径

- 取锚时刻 `2026-10-03 10:01:25 +0800`；`git log -1 --format=%H` = `ccd9543b3d86f3f627d331b3c399cb3d4ee99eec`；分支 `dev`。
- 工作树**非干净**（`git status --porcelain` 有他程在飞的改动，含 `cmd/wisp/config_receipt_255_test.go` 与
  `.scratch/wisp/probes/**` 若干），故本篇所有**产码读数一律取自工作树**、所有**归属读数一律用 `git log` 指到具体 commit**，
  两者不混用。工作树里 `internal/agent/{compress.go,loop.go,budgets.go}` 与 `cmd/wisp/run.go` 当前 **0 枚改动**
  （`git status --porcelain` 全文里无 `internal/` 也无 `cmd/` 的 M 行，除上面那枚 255 的测试文件），所以本格的产码读数可直接当锚点树读。
- 尺：只读工具（`Read` / `Grep` / `find` / `grep` / `sed -n`）。⛔ 本篇**未跑任何一枚 `go` 子命令**
  （`build` / `vet` / `test` / `run` / `list` 全零），原因＝同机三枚写腿在飞，编译会互洗读数。
  因此本篇**不报任何 rc、不报任何用例通过与否**；需要 rc 的判据一律落到 §6。
- **grep 分母口径**（本仓踩过：变异副本长得跟产码一样）：
  `cmd/**` + `internal/**` + `docs/**`，**排除 `.scratch/**`**。这条不是洁癖——
  `./.scratch/wisp/probes/183/a1/mut/loop.go` 是一枚**完整的 loop.go 副本**，
  它在 `res.Compression = rep`（该行副本 `:404`）和 `comp *Compressor`（`:182`）上都命中，
  不排除就会把「生产侧读者」数出虚高的一枚。本篇每条全称否定都按这个分母跑，并带正控。
- 路径写法：报到**包／目录级**＋行号（例 `internal/agent/compress.go:174`），不停在文件名。
- 禁令执行：`frontend/**` 与 `design/**` 未读、未引用；`internal/observe/thresholds.go`、golden、
  `tools/d22scan/allowlist.txt`、`.github/workflows/ci.yml`、`docs/PLAN.md`、`docs/specs/**`、票面 **一字未动**；
  产码／测试 **0 字节改动**（本程唯一写面＝本文件与同目录的提交说明 txt）。

---

## 1. 压缩到底在哪跑

### 1.1 生产侧唯一的入口链（逐枚行号，工作树锚 `ccd9543`）

| 跳 | 位置 | 是什么 |
|---|---|---|
| ① 构造（每枚 Loop 一次） | `internal/agent/loop.go:239` `comp: NewCompressor(b, opt.Summarizer, WithLogger(opt.Logger))` | 压缩机**不是按任务造的**，是 `agent.New` 里按 Loop 造的；`opt.Logger` 在 `:207-209` 已被兜成 `slog.Default()` |
| ② 字段落位 | `internal/agent/loop.go:182` `comp *Compressor` | Loop 的共享字段，跨任务复用 |
| ③ 判该不该压 | `internal/agent/loop.go:393` `if l.comp.Need(hist)` → 实现 `internal/agent/compress.go:119`（阈值比较在 `:120`） | 阈值＝ `internal/agent/budgets.go:107` 的 `HistoryCompressTokens`，由 `:36` `refHistoryTokens = 12000` 按 `:91` `Scale` 等比缩放（`scaleInt` 在 `:114`，只收不放） |
| ④ 真压 | `internal/agent/loop.go:399` `nh, rep, err := l.comp.Compress(withTraceTask(ctx, taskID), hist)` → 实现 `internal/agent/compress.go:174` | `withTraceTask` 定义在 `internal/agent/compress.go:152`，读侧 `:160` |
| ⑤ 失败边 | `internal/agent/loop.go:401` `l.log().Warn("agent: history compression failed", "err", err)` | 唯一一条**在调用方**的痕 |
| ⑥ 成功边 | `internal/agent/loop.go:402` `else if rep.Ran` → `:403` `l.replaceHistory(nh)` → `:404` `res.Compression = rep` | `replaceHistory` 见 `internal/agent/loop.go`（历史换手）；痕本身在压缩机里（下条） |
| ⑦ 痕的产出点 | `internal/agent/compress.go:240` `c.log().Info("agent: history compressed", attrs...)`，由 `:216` `if rep.Ran` 把门 | 属性表在 `:222-231`，task 键在 `:237-239` |

### 1.2 它在哪个循环里

`internal/agent/loop.go:376` `for {` —— **①~⑦ 整条链在 run 的轮次循环体内**（`hist := l.History()` 在 `:392`，
压缩紧跟其后，`buildRequest` 在 `:408`）。⇒ **一枚任务的一次 Run 里这段代码可以被走 N 次**（每轮模型调用前各一次），
不是一个任务一次的形状。这条直接决定 §2.3 那格「算了但没存」的读数。

### 1.3 「活在 run 局部」那一格（本仓旧账，不是新发现，但必须量全）

`res` 是 `run` 的局部变量：`internal/agent/loop.go:369` `res := Result{TaskID: taskID, Status: StatusCompleted}`。
它出 run 只有三条口，全部经 `Result` 这一个结构：

1. `RunAsync` 把它塞进 `t.result`（`internal/agent/loop.go:326`），由 `Wait()` 返回（`:308-314`）；
2. `Run` 直接返回（`:332-334`）；
3. 各收尾分支 `l.finish(...)` / `brakeStuck` / `failWith` 收的也是这枚 `&res`。

生产侧**真的拿到 `Result` 的腿**只有两枚：

- `cmd/wisp/run.go:1099` `bg := loop.RunAsync(ctx, task)` → `:1106` `res := bg.Wait()`，
  消费面是 `cmd/wisp/run.go:1184` `notifyBody(res agent.Result)` 与 `:1195` `costLine(res agent.Result)`；
- `internal/tools/subagent_197.go:344` `bg := child.RunAsync(...)` → `:383` `res := bg.Wait()`，
  消费面是 `:412` `finalize(...)` 与 `:437` `answer(...)`。

⇒ 「泵／快照读不到 run 局部」这一格在 `Compression` 上是**同一形状的第二次命中**：
`res` 到得了宿主、到得了子代理工具，**但它身上那枚 `Compression` 字段在两枚消费面里都没被看一眼**（逐跳见 §4）。

### 1.4 装配面（谁把这格接到盘上）

- `cmd/wisp/run.go:995` `opts := agent.Options{...}` —— 具名读一遍：`Provider` / `Tools` / `Sink` / `Journal: nil` /
  `AdmitTask` / `Registry` / 内嵌 `Config` 全在，**没有 `Logger` 这一项**。⇒ 生产 Loop 的压缩机走的是
  `internal/agent/loop.go:207-209` 的兜底 = `slog.Default()`。
- 常驻腿的面板入口**不装 sink**：`cmd/wisp/panel_inbound.go:49` 明写「it installs no persistent log sink (installLogSink)」，
  且 `grep -rn 'agent\.New(' cmd internal`（排除测试）在生产侧只有 **2 枚命中**：
  `cmd/wisp/run.go:1023`（`wisp run` 腿）与 `internal/tools/subagent_197.go:319`（子代理腿）。
  ⇒ **今天这台机器上会走到压缩的只有 `wisp run` 那条 CLI 腿和它的子代理**，常驻面板腿根本不构造 Loop（尺见 §3.4）。
- 落盘那条管子（痕**能**到盘的证据）：`cmd/wisp/run.go:222` `installLogSink(s.dataDir)` 在
  `:249` `assembleRuntime(s)` **之前**；`assembleRuntime` 内 `:995` 才造 `agent.Options`、`:1023` 才 `agent.New`。
  `installLogSink` 实现在 `cmd/wisp/logsink.go:144`，`:153` `slog.SetDefault(slog.New(teeHandler{...}))`，
  `teeHandler.Handle` 在 `:208-214`（primary＝`p.Handler()`＝带 redact 的 JSONL 文件，mirror＝stderr，`min_level` 写死在 `:82`）。
  ⚠ 一个**时序口径**必须记：`agent.New` 在 `:207-209` 把 `opt.Logger` 兜成的是**那一刻**的 `slog.Default()`，
  压缩机字段 `lg` 于是**持有一个 Logger 指针**（`internal/agent/compress.go:77` / `:87-89` / `:102-107`），
  不是每次调用现取默认值 ⇒ 若哪天真出现「先 `agent.New` 后 `installLogSink`」的装配序，痕会进旧默认（stderr），
  今天的顺序不会（`run.go:222` 早于 `:1023`）。

---

## 2. 它算了什么、扔了什么

### 2.1 一次 `Compress` 调用里算出来的每一枚量（逐枚行号 × 存到哪 × 能不能被答出来）

| 量 | 算在哪 | 存到哪 | 出了这次调用还有人能答吗 |
|---|---|---|---|
| 压前 token | `internal/agent/compress.go:175` `rep := CompressionReport{TokensBefore: c.TotalTokens(hist)}` | 字段 `:69` ＋ 日志 `:223` | 日志能；结构体能（但结构体没人读，见 §4） |
| 阈值 | `:120` / `:180`（两处同一 `c.b.HistoryCompressTokens`） | **只进日志 `:225`**，`CompressionReport` 里**没有这个字段**（结构体定义 `:61-71`） | 只有日志能答 ⇒ 拿不到 `res.Compression` 的一方**无法算占比**（`tokens_after/threshold`） |
| 压后 token | `:212` `rep.TokensAfter = c.TotalTokens(out)` | 字段 `:70` ＋ 日志 `:224` | 同上两路 |
| 压前/压后条数 | `len(hist)` `:226`、`len(out)` `:227` | **只进日志**，字段表里没有 | 结构体侧答不了 |
| 被折掉的条数 | `:208` `rep.CompressedMsgs += len(victims)`（跨内层迭代累加） | 字段 `:63` ＋ 日志 `:228` | 两路都有 |
| 留下的原始轮数 | `:213` `rep.KeptRawRounds = len(rawRoundIndexes(rounds))` | 字段 `:64` ＋ 日志 `:229` | 两路都有 |
| 历史真的动了没 | `:230` 现算 `rep.TokensAfter != rep.TokensBefore \|\| len(out) != len(hist)` | **只进日志**，字段表里没有 | 结构体侧答不了；`:402` 那个 `else if rep.Ran` 用的是另一枚量 |
| 这次是哪枚任务的 | `:160` `traceTaskID(ctx)` 读、`:237-239` 追加 | **只进日志**（键 `task`） | 结构体侧没有 taskID；`Result.TaskID`（`internal/agent/loop.go:83`）在场，但两者之间没有代码把它们对上 |
| id 台账（硬规则那枚） | `:178` `var ledger []string` 累加、`:196` 追加、`:214` `rep.PreservedIDs = dedupe(...)`、`:215` `rep.RawIDsInsideSummary = dedupe(ledger)` | **只到内存字段 `:65` / `:66-68`**，不进日志（隐私口径 `:59-60`） | **答不出**：全树非测试读者 0 枚（尺见 §4.1），`Compress` 一返回就没人持有 |
| 折了几轮（次数） | 内层循环 `:180-209` 可迭代多次，每轮 `fold := raw[:len(raw)-c.b.KeepRawRounds]`（`:185`） | **没有任何计数器**：`rep` 里只有 `Ran bool`（`:207`） | 答不出「这次折了几轮」 |
| 折之前一共有几枚原始轮 | `:181` `raw := rawRoundIndexes(rounds)` | **`len(raw)` 从没进过任何输出** | 答不出「从 N 枚原始轮压到 3 枚」里的 N |
| 总轮数 | `:176` `rounds := groupRounds(hist)` | `len(rounds)` 从没算过 | 答不出 |
| 是不是**在旧摘要上再摘要**（二次信息复损） | `:190-192` `if idx := summaryRoundIndex(rounds)` ＋ `:206` `summaryMarkerIn(victims)` 传给 `rebuildRounds` | **算了、用了、没记**（既不在字段也不在日志） | 答不出「这段历史今天被复损过第几次」 |
| 折进去的正文本身 | `:198` `c.summarize(ctx, victims)` → `:202-205` `summaryBlock` | 进 `out` 的新消息，即**新历史本身** | 在；日志侧**故意不进**（`[privacy] keep_transcript` 硬编码 false，`internal/config/schema.go:514-515`、`internal/config/validate.go:82-84`） |

### 2.2 「算了但没存」这格的确切形状（本票核心）

上面表里加粗的五行是同一件事的五个面：**阈值、条数、动没动、taskID、以及折前规模（`len(raw)` / `len(rounds)` / 折叠次数 / 是否叠摘要复损）
这组量，今天有一多半只活在日志那一路，另一多半连日志都没有、算完即弃。**
具体到落地腿要挑的话，`internal/agent/compress.go:214` 与 `:215` 那两枚（id 台账）是**唯一已经落在 `CompressionReport` 字段上、
却既不出声也没人读**的量——它们正是 D15(4) 那条「tool_call id 永不丢弃」硬规则（`:18-21` 注释）的**证据本体**。

⚠ 一句防止落地腿跑歪的话：**把 id 台账本身写进日志是禁的**——`:59-60` 与 `keep_transcript` 那条硬约束覆盖的是历史内容；
要出声只需要**计数**（`len(rep.PreservedIDs)` / `len(rep.RawIDsInsideSummary)`），
本仓已有的隐私探针会盯正文（`internal/agent/compress_trace_test.go` 里那两枚哨兵，见 §5.3）。

### 2.3 一枚任务里这段被走 N 次，只有最后一次活下来

- `internal/agent/loop.go:376` 是轮次循环，`:392-406` 在循环体内 ⇒ 一枚任务可以多次触发压缩。
- 成功边是**覆盖写**：`:404` `res.Compression = rep`；字段自己的注释（`:99`）就写着
  "the last history compression pass (zero if none)"。⇒ **「这个任务一共压过几次、累计省下多少」在 `Result` 这一路今天无解**。
- 但**日志那一路是能答的**：每次折叠各打一条 `:240`，且 `:237-239` 带 `task` 键 ⇒ 同一条任务的多枚痕能靠 task id 聚合。
  ⇒ 落地腿若要加「累计」，先看清这不是从零起：它要么在 `rep` 上加累计字段（改结构体＝改 `internal/agent/**` 的出线），
  要么在日志侧靠计数聚合（**今天已经可行，零改动**）。

### 2.4 被丢的部分还能找回吗（三态，逐条给尺）

1. **内存里的原文**：`rebuildRounds`（`internal/agent/compress.go:406-424`）把被折轮次整批 `drop`，
   调用方 `internal/agent/loop.go:403` → `:1016-1020` `l.history = hist` 直接换手 ⇒ 换手后**没有任何一方再持有旧切片**
   （`Compress` 声明它不改入参内容 `:172-173`，但没人留着那个入参）。
2. **磁盘上的原文**：找不回。三处都读了——
   环路侧 `Journal: nil`（`cmd/wisp/run.go:1004`，理由写在 `:1000-1003` 注释里：桥已写权威 tool_call 行）；
   `internal/memory/schema.go:23-136` 的**七张表里没有任何 message/conversation 表**（`task_log` 在 `:53`，只到任务粒度）；
   `internal/session/grants.go:67` `Ledger` 存的是授权 pattern（`:146` `Record(ctx, tool, pattern)`），不是正文。
   且 `keep_transcript` 是**硬编码 false、写 true 直接报错**（`internal/config/schema.go:514-515`、`validate.go:82-84`）
   ⇒ 「留原文副本」这条路**根本没有合法落点**，任何落地腿别往这里想。
3. **还能保住的**（结构性残留，不是原文）：`structuralTrim`（`:283-306`）留 tool_call id / 工具名 / 参数字节数 / 结果状态 /
   图片 mime，散文只留**头 80 rune**（`:290` `firstRunes(t, 80)`，实现 `:308-314`）；`summaryBlock`（`:269-279`）再追加整条 id 台账。
   ⇒ 所以「压缩=信息全丢」是**错的**，准确说法是**散文变 80 rune、id 与状态全留、且留痕只记计数**。
   ⚠ 这条 80 rune 的形状有个副作用要留给落地腿：**「往正文尾部埋哨兵」这类隐私探针单用会漏**
   （尾部被截掉），票 139 实现腿在 `:91` 那段自打里就是撞在这上面——`internal/agent/compress_trace_test.go` 里那两枚探针是这么来的。
   现量确认这两枚都在、且各管一头：尾部哨兵 `internal/agent/compress_trace_test.go:186` 定义、`:193` 埋进每一条 TextPart 的**末尾**、
   `:257-258` 检；头部探针 `:261-263`（检 `"问题"`，正是夹具用户轮的第一个词，`internal/agent/compress_test.go:31`
   `fmt.Sprintf("问题 %d %s", i, pad)`）⇒ **`firstRunes(t, 80)` 砍得掉尾哨兵、砍不掉头探针**，两枚一起才闭环。

---

## 3. 「零 log 调用」复量（现跑尺＋正控＋历史对照）

票面 §0 第 5 行那句「`compress.go` 全文零条日志调用」**在现树已不成立**。不许沿用，以下是我这轮现跑的。

### 3.1 尺的字面命令与读数

```
grep -n 'Info(\|Warn(\|Debug(\|Error(\|Trace(\|log()\|slog\.' internal/agent/compress.go
```

7 行命中，逐行分型（**别把定义当成调用**）：

| 行 | 是什么 | 算不算一次日志出口调用 |
|---|---|---|
| `internal/agent/compress.go:45` | 注释里出现的 `slog.Info` 字样 | 否 |
| `:77` | 字段声明 `lg *slog.Logger` | 否 |
| `:87` | 形参声明 `func WithLogger(lg *slog.Logger) CompressorOpt` | 否 |
| `:102` / `:106` | 包内 logger 取用器 `func (c *Compressor) log() *slog.Logger` 及其 `return slog.Default()` | 否（它自己不落痕） |
| `:240` | `c.log().Info("agent: history compressed", attrs...)` | **是，唯一一枚** |

⇒ **读数：`internal/agent/compress.go` 今天有 1 枚日志出口调用**（`:240`，`Info` 级），
外加 1 枚把门 `:216` `if rep.Ran`、8＋1 枚属性（`:222-231` 的 8 个键 ＋ `:237-239` 条件追加的 `task`）。

### 3.2 正控（同一把尺在别处必须命中，否则尺是瞎的）

```
for f in internal/agent/loop.go internal/agent/spill.go internal/agent/tools.go \
         internal/agent/journal.go internal/tools/subagent_197.go cmd/wisp/run.go; do \
  printf '%s : ' "$f"; grep -c 'Info(\|Warn(\|Debug(\|Error(' "$f"; done
```

| 文件 | 命中行数 |
|---|---|
| `internal/agent/loop.go` | **14** |
| `internal/agent/spill.go` | 0 |
| `internal/agent/tools.go` | 1 |
| `internal/agent/journal.go` | 0 |
| `internal/tools/subagent_197.go` | 2 |
| `cmd/wisp/run.go` | 2 |

⇒ 尺看得见（同包同形状的文件给 14／2／1），**但 0 命中不等于缺陷**：
`internal/agent/spill.go` 与 `internal/agent/journal.go` 今天也是 0 条，它们的可观测面在别处
（spill 写 artifacts 文件、journal 写 `memory` 的 `task_log`/`tool_call` 行，`internal/agent/journal.go:19-23`）。
⇒ **这条对落地腿重要**：本格的判据不该是「`compress.go` 有日志」，而是「owner 那两问能答」；
`compress.go:240` 那枚日志存在，是否兑现成「答得出」，看 §4。

### 3.3 历史对照（同一把尺、同一个文件，修码前后各跑一次）

```
for r in d25f139b 1acdd032 ff550f38 b23c7f77 HEAD; do \
  printf '%s compress.go log-calls: ' "$r"; \
  git show $r:internal/agent/compress.go | grep -c 'Info(\|Warn(\|Debug(\|Error('; done
```

读数：`d25f139b` → **0**，`1acdd032` / `ff550f38` / `b23c7f77` / `HEAD` → **1**。

- `1acdd032 fix(139 AC#2+AC#3)` 是那枚把 0 变成 1 的 commit；`ff550f38 test(153 AC#1)` 与 `b23c7f77 feat(153 AC#2)` 各补判据/补 `task` 键。
- ⇒ 票面 §0 那格、以及票面「§0 五条锚」里的第 5 条，**从 `1acdd032` 起就过期了**（见 §7 第 1 条）。

### 3.4 「谁产出／谁投递／谁落盘」三处都读了才写的结论（本仓规矩：数枚数答不了能力问题）

- **产出**：`internal/agent/compress.go:240`，把门 `:216` `if rep.Ran` ⇒ **只有真折了才出声**；
  `:180` 那个内层循环折不折由 `:182` `len(raw) <= c.b.KeepRawRounds` 定。
- **投递**（三级，逐枚）：`internal/agent/compress.go:102-107` `log()` → 字段 `:77` `lg` →
  由 `:87-89` `WithLogger` 在 `internal/agent/loop.go:239` 注入 `opt.Logger` →
  `opt.Logger` 生产侧**从没被赋值**：`cmd/wisp/run.go:995-1021` 那段 `agent.Options` 里没有 `Logger` 项，
  且 `grep -rn 'Logger:' --include='*.go' cmd internal \| grep -v '_test\.go'` = **0 命中**；
  兜底在 `internal/agent/loop.go:207-209`（`opt.Logger = slog.Default()`）。
  ⇒ 生产侧唯一的投递路径**就是进程默认 Logger**；`internal/agent/harness_test.go:87` 那枚 `withLogger` 是**测试专用**注入口。
- **落盘**：`cmd/wisp/logsink.go:144` `installLogSink` → `:149` `observe.InitLog` →
  `:153` `slog.SetDefault(slog.New(teeHandler{primary: p.Handler(), mirror: stderr}))` →
  `:208-214` `teeHandler.Handle`（primary 先、mirror 后）。
  primary 内部是 `internal/observe/logging.go:87-90`（`redactHandler` 套 `slog.NewJSONHandler`，级别 `:90`），
  目录是 `cmd/wisp/logsink.go:76`/`:87-89` 的 `<dataDir>\logs`，级别 `:82` 写死 `"info"` ⇒ **`Info` 级过得去**。
  ⚠ 这条管子**不设「压缩」这一维的开关**：痕到不到盘只取决于 sink 装没装、级别过不过。
- **两条腿的装配顺序都成立**（不是只成立一条）：
  `wisp run` 腿 `cmd/wisp/run.go:222`（装 sink）早于 `:249` `assembleRuntime` → `:1023` `agent.New`；
  常驻腿 `cmd/wisp/resident_windows.go:63`（装 sink）早于 `:206` `startResidentTaskSource` →
  `cmd/wisp/resident_task_source_windows.go:265` `assembleRuntime` → 同一个 `:1023`。
- ⇒ **结论（能力，不是枚数）**：这枚痕**今天有两条真实的到会落盘的腿**（`wisp run` 与常驻任务源），
  不依赖任何测试注入口；`cmd/wisp/panel_inbound.go:49-54` 那一腿**明写不装 sink 且不造 Loop**
  （`grep -n 'assembleRuntime\|agent\.New' cmd/wisp/panel_inbound.go` = 0 命中）⇒ 那条腿连压缩都不会发生，不是「发生了没记」。

### 3.5 会不会响（触发可达性，纯结构读数，未跑码）

`Compress` 折得动需要**两件事同时成立**：

1. `internal/agent/compress.go:119-120` `Need()`：总 token **>** `b.HistoryCompressTokens`（缩放值，`internal/agent/budgets.go:107`）；
2. `internal/agent/compress.go:182-184`：`len(raw) > b.KeepRawRounds(3)`，`raw` 由 `:376-384` `rawRoundIndexes` 数**非摘要轮**，
   轮的分界在 `:347-362` `groupRounds` 的 `:350`（**只有 `RoleUser` 才起新轮**）。

⇒ **一轮 = 一枚 user 角色消息 ＋ 它后面跟着的一切**。单次任务里 `RoleUser` 的产生点只有三处：
`internal/agent/loop.go:371`（任务输入本身，实现 `:1045-1047`）、
`:281`→`:1028` `drainSteering` 里的 steer（`userMessage`， steering 通道）、
`:493` `reminderMessage`（实现 `:1080` `userMessage("[系统提醒] " + text)`）——
**C22 提醒阶梯与 steering 都会各起一枚新「轮」**，工具结果不是（`RoleTool`，`:715`）。
并且 `internal/agent/loop.go:232-249` 造 Loop 时**不给 history 播种**（`l.history` 零值），
而 `cmd/wisp/run.go:990` `execute` 每枚任务都现造一枚 Loop（`:1023`）⇒ **历史不跨任务累计**。

⇒ 所以：**今天要让这枚痕响一次，得在一枚任务里凑出 ≥4 枚 user 角色消息且总量过阈值**——
正好是票面 §1 那条 09-26 更正（普通收尾形 0 枚 / 退化形 1 枚）的**结构原因**，
我这轮是从码里读出来的、没跑任何用例。⇒ 落地腿别把「加日志」当成「owner 那两问就有答案」：
**常见路径压根不折，因此压根不打**；能答的只有「这次如果折了，折了多少」。
