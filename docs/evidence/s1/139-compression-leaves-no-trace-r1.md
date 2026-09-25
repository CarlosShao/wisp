# 139 — 上下文压缩留痕（Go 侧那半）· 实现程证据件 r1

票面：`.scratch/wisp/issues/139-context-compression-leaves-no-trace-zero-log-calls-in-compress-go-zero-readers-of-res-compression.md`
件本落点规矩：`docs/evidence/s1/`（裁决表）＋ `.scratch/wisp/probes/139/`（原始读数）。
**本件出自实现者本人，不是裁决**——框不自勾，等非实现者程复算。

## 0. 锚点与工作树现量

```
$ git rev-parse HEAD
3f6322f1b6a5d77fac4e17dd0fa7ea33c19d3d53        # 与派单锚点 3f6322f 相同
$ git rev-parse --abbrev-ref HEAD
dev
$ git log --oneline -10                          # 全量贴进 probes/139/git-state.log
3f6322f evidence(138 落盘代理未提交的尾段) …
2f5c7f9 probes(149): 变异尺与逐发原始读数落盘
f8623c7 evidence(138 AC#4 门禁 第 6 格) …
b417d31 fix(149 AC#2+AC#3+AC#4) …
```

`git status --porcelain`（现量，`probes/139/git-state.log`）里 **`internal/agent/**` 与 `internal/observe/**` 各 0 枚改动**
⇒ 票面 §Packages「注意地界：`internal/observe/**` 此刻可能有别的程在跑」是当时写的，**本程所见现量为无人**；
本程因此**根本没碰** `internal/observe/**`（AC#2 的三选一里也没选它，理由见 §2）。
在飞的是 `cmd/wisp/**`（票 149）与 `frontend/**`＋`design/**`（owner/前端会话），本程一律零字节。

**地界自证**：本程写过的路径逐枚列在每格末行「本格改了哪些文件」，与 §派单「你的地界」逐名对齐。

---

## 1. AC#1 — 把"零留痕"量成读数（含票面 §0 两条未复算项）

**判据**（票面 AC#1 原文三段）：① `Compression` 读取者名册逐名（含 `cmd/wisp/**`、面板桥、诊断包）；
② `refHistoryTokens` 与 `context_window` 的关系，若确为等比缩放则本票例子一律不许拿 `12000` 当"实际触发点"；
③ `DEFERRED(D28-1)` 登记条目原文＋回答"压缩无留痕今天是否已被那格覆盖"。

### 1.1 量之前先把尺打一遍正控（本仓规矩：0 读数必须先证尺活着）

| 尺 | 正控（已知存在，必响） | 读数 | 待量目标 | 读数 |
|---|---|---|---|---|
| `git grep -n '\.X\b' -- '*.go' ':!*_test.go'` | `\.Rounds\b` | **3 枚**（含真读者 `cmd/wisp/run.go:608`） | `\.Compression\b` | **1 枚** |
| 同上 | `\.Usage\b` | **54 枚** | 同上 | 同上 |
| `git grep -nE 'Info\(\|Warn\(\|Debug\(\|Error\(' -- <file>` | `internal/agent/loop.go` | **14 枚** | `internal/agent/compress.go` | **0 枚（rc=1）** |

命令原文与输出：`.scratch/wisp/probes/139/ac1-readers.log` [P1]–[P3]。
⇒ 两把尺都会响，所以 1 枚与 0 枚不是"正则恒不匹配"造出来的漂亮 0。

### 1.2 ① 读取者名册（逐名，非计数）

`git grep -n '\.Compression\b' -- '*.go' ':!*_test.go'` 的**全部**命中：

1. `internal/agent/loop.go:401  res.Compression = rep` —— **写在等号左边，是生产者不是读取者**。

按票面点名的三处逐名核（`ac1-readers.log` [P4][P5]）：

2. `cmd/wisp/**`（非测试）拿到 `agent.Result` 的只有三处：`run.go:576 res := loop.Run(ctx, task)`、
   `run.go:647 notifyBody(res agent.Result)`、`run.go:658 costLine(res agent.Result)`。
   这三处逐名读的字段是：`TaskID · Status · Rounds · ToolCalls · Message · Text · Usage.InputTokens ·
   Usage.OutputTokens · CostMicros · Currency · Err`（`res.` 形状普查，见 [P4] 与件内 grep）——
   **`Compression` 不在其中**；`run.go:607-608` 那行用户看得见 summary 只印 任务号/状态/轮数/工具数/成本。
3. 面板桥：`internal/panel/**`、`internal/tools/**`（非测试）对 `compression` 的大小写不敏感 grep **0 命中**；
   全仓非测试文件里 `agent.Result` 只被 `cmd/wisp/run.go` 与一处测试签名
   （`internal/tools/loop_approval_test.go:252 dumpToolLog`，读的是 `ToolLog` 不是 `Compression`）引用。
4. 诊断包 `internal/observe/diagnostics.go`：`compression` **0 命中**；该文件字面量键名正控为
   `"id" "ticket" "summary" "registered" "false" "reason"`（尺活的）——它记的是检查项，不接 agent 侧观测。
5. 生产侧**也没有人整枚序列化 `Result`**：`%+v`/`%#v` 在 `cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`
   的命中**全部落在 `_test.go`**，无一处作用于 `agent.Result`（[P4] 同件）。
   ⇒ 这条堵住了"字段没人点名读、但整枚被 `%+v` 打出去所以其实看得见"这一支。

**读数**：票面 §0 那行「全仓唯一命中就是赋值那一行 `loop.go:401`，零读取者」**复算成立**，
且比"零读取者"更严：既无字段级读者、也无整结构级打印者。

### 1.3 ② `refHistoryTokens` 与 `context_window` 的关系 = 等比缩放（裁决表说法复算成立）

代码链（`ac1-scaling-deferred.log` [P7]）：`budgets.go:36 refHistoryTokens = 12000` →
`:91 Scale = float64(ctxWindow)/float64(ReferenceContextWindow)`（`ReferenceContextWindow = 128000`，`:18`）→
`:107 b.HistoryCompressTokens = s(refHistoryTokens)`，`s = scaleInt(ref, Scale)`（`:114` 只降不升、下限 1）。
实证：现跑既有用例 `go test ./internal/agent/ -run TestCompressTriggerScalesWithWindow -v` → **PASS**，
它自己算的是 `BudgetsFor(4096).HistoryCompressTokens == (4096*12000)/128000 == 384`。

⇒ **AC#1② 的禁令生效并已被本票遵守**：本票新增用例里 `12000` **只出现在"参考窗 128k 下的冻结参考值"这一格**
（那是 §4.3 原文数值，且既有 `compress_test.go:53-55` 已经这么用），
**触发点一律从 `Budgets.HistoryCompressTokens` 读出来，不写死**（见 §3 的 `th := c.b.HistoryCompressTokens` 形状）。
现量（`git grep -n '12000' -- 'internal/agent/*_test.go'`，本格时刻）：**8 命中，全部既有、全在 `compress_test.go`**
（`:53,54,55,171,172` 是"128k 参考窗下的冻结值"断言、`:174` 是缩放算式 `(4096*12000)/128000`、
`:168,180` 是注释），**无一枚把它当"实际触发点"** —— 本票新增用例对 `12000` 的贡献为 **0 枚**。

### 1.4 ③ `DEFERRED(D28-1)`：那格**不存在**，故"压缩无留痕"未被覆盖

代码侧标记（[P9]，逐名）：`internal/agent/compress.go:26`、`internal/agent/loop.go:394` 两枚，
内容都是**调用位置**这一件事（"Warm-window hook 才该拥有这次调用，同步回退是 S1 接受的形状"）。

登记侧现量：`grep -c 'D28-1' docs/specs/SPEC-12-roadmap-governance.md` → **0（rc=1）**，
而同文件 §5 表体有 47 行 `| ` 开头的条目 ⇒ 那枚 0 是扫过了表得到的，不是文件没读到。
仓内 D28-1 的全部文字痕迹只在：`docs/evidence/s1/10-adversarial-acceptance.md:106`
（"登记在 `compress.go:22-27`，不算违规"）、`docs/reports/2026-09-24-gap-analysis-audit-verdict.md:69`、
`docs/reports/pending-and-issues.md:4631,4708`。

**回答票面那一问**：**未被覆盖**。两层理由：
(a) `SPEC-12 §5` 里**根本没有 D28-1 那一格**（代码有标记、表里无条目 ⇒ 反向核对不成立，AGENTS.md §1.1 那条 1:1 双向）；
(b) 即便看代码注释的语义，D28-1 说的是**这次调用发生在哪儿**，与本票的**这次调用留下了什么**是两件事，
票面 §1 自己就把它们分开了（"本票不改这个决定，只给它补上'这一次到底压了多少'的读数"）。

⇒ 按 AC#1③ 的第二支「若未覆盖 ⇒ 去 `SPEC-12 §5` 按五字段补登记」应当补登记。
**但 `docs/specs/**` 在本程的零字节禁改面上**（派单 §2），台账 `pending-and-issues.md` 亦归编排者 ⇒
**本程不落笔，五字段备妥交编排者**（同 138 那程对 `SPEC-12 §5` 的处置形状）：

| 五字段 | 备妥文本 |
|---|---|
| 类型 | DEFERRED |
| 项 | `DEFERRED(D28-1)` 压缩调用位置（表里现在**缺这一格**，代码有标记） |
| 为什么现在不做（依据） | 票 10 契约明文豁免同步回退（`10-adversarial-acceptance.md:106`）；Warm-window hook 属票 28 |
| 完成判据 | 压缩由 Warm-window hook 拥有；同步路径只在无 hook 注册时兜底；留痕不随位置改变（`internal/agent` 用例钉住） |
| 当前残缺表现 | 压缩在响应路径上同步发生（时延落在用户那一跳）；**"这次压了多少"自本票起有 log 痕，但痕里不带 taskID** |

⚠ 同一把尺顺手扫到的、**不归本票修**的：`grep -c 'D11-3' SPEC-12` 亦为 **0**（代码 `internal/agent/control.go:15` 有标记）。
⇒ 表↔码双向核对的缺口**不止 D28-1 一枚**，本程只报不修（改 `docs/specs/**` ＝人工批准）。

### 1.5 放水两问自答

- **这一发在未修码上响不响？** 本格**不改码、不立判据**，只把票面 §0 的四条锚与两条未复算项落成读数；
  所有"0"都先打了正控（[P1][P3][P5-control]）。未修码上量到的 0 ⇒ 是真 0。
- **摘掉本票这一味，是否存在一发变异从此打不红？** 本格不适用（无判据）；判据在 §3。

### 1.6 本程没测什么（按"漏了它谁会先被骗"排序）

1. **没测"生产进程里这条痕真的落到 `<data>\logs\wisp-*.jsonl` 了吗"**。本格只读到
   `cmd/wisp/logsink.go:82 logSinkLevel = "info"`、`installLogSink` 在 `run.go:178`/`models.go:284`/
   `resident_windows.go:57`/`secret.go:226` 四处被调、且 `cmd/wisp/**` 非测试代码**从不设 `agent.Options.Logger`**
   （`git grep -n "Logger:" -- 'cmd/wisp/*.go' ':!*_test.go'` → rc=1）⇒ 环环推得"log 会落盘"，
   但**端到端量一次**需要真起进程＋真数据根，本程没量。若断在哪环，被先骗的是"事后答得出"这句话本身。
   ⚠ 另注：MEMORY 里"安全告警默认只到 stderr（持久 sink 只在 `wisp slo`）"那一条与 `logsink.go` 现码**不符**
   （现码 `installLogSink` 是 file+stderr 双路 tee），本程按**现码**走，并按派单 §0 的注入规矩记在这里、不采信记忆文本。
2. **没测 `wisp run` 之外的宿主入口**（`resident_windows.go` 那条常驻路径）会不会因为自己传了 Logger 而丢掉这条痕——
   现量 `Logger:` 在 `cmd/wisp/**` 非测试里 0 命中，所以今天不存在这一形，但那是"今天没人这么写"，不是"写了会被抓到"。
3. **没测 `Result.Compression` 出线在 149 那程手里会不会长出读者**——`cmd/wisp/**` 是别人的地界，本程零字节、
   也**不预判**。若它长出读者，本票 §2 里"选 (b) 无解"那句话的**时效**要重读。
4. **没测 D28-1 反向核对缺口的完整名册**：只顺手核了 D28-1、D11-3 两枚，其余 20+ 枚 `DEFERRED(...)` 代码标记
   是否各有 §5 条目没扫（那是编排者/验收面的活，本票只在 AC#1③ 点名范围内答）。

**本格改了哪些文件**：新增 `docs/evidence/s1/139-compression-leaves-no-trace-r1.md`（本件）、
`.scratch/wisp/probes/139/ac1-readers.log`、`.scratch/wisp/probes/139/ac1-scaling-deferred.log`、
`.scratch/wisp/probes/139/gate-pre-full.log`（§1.7 基线）。零生产码改动。

### 1.7 改前门禁基线（后续格拿它做名册差集的分母）

```
$ go test ./internal/agent/ -count=2 -v   > probes/139/gate-pre-full.log ; rc=0
RUN=152  PASS=152  FAIL=0  SKIP=0        # 四数（-count=2 ⇒ 76 枚 ×2）
$ grep -n '^FAIL\|^ok ' gate-pre-full.log
361:ok  	github.com/CarlosShao/wisp/internal/agent	5.041s
```
⚠ 按派单 §5 的坑：包级汇总行 `ok` 不参与 FAIL 计数，FAIL 只数 `--- FAIL`（本基线 0 枚）。

---

## 2. AC#2 — 最小落盘面：三选一选了"结构化日志字段"，且落在 `Compress` 的成功边上

**判据**（票面原文）：只要求"事后答得出"、不要求任何 UI；压缩**成功**时留下一条可核的痕，
内容至少含 **压前 token / 压后 token / 是否真的动了历史**；三者选一（结构化日志字段／`Result` 出线可读字段／诊断包一项）
**并在票面说明为什么选它**。硬约束两条：禁新增 Go→前端事件或面板方法（＝契约变更）、禁把原文内容落进日志。

### 2.1 三选一的裁决与代价（为什么不是另外两枚）

| 候选 | 裁决 | 依据（全部出自 §1 的现量，不是偏好） |
|---|---|---|
| (b) `Result` 出线上的可读字段 | **不选：这一味等于零改动** | 字段**已经存在**（`loop.go:100` 声明、`:401` 赋值，含 `TokensBefore/TokensAfter/Ran`）。§1.2 现量它在生产侧**既无字段级读者、也无整结构 `%+v` 打印者** ⇒ 再往这条线上加字段仍然是零读取者，"事后答得出"一句都不兑现。要让它可答就得去 `cmd/wisp/**` 加打印 ⇒ **那是别家的地界**（149 在飞），且改的是 CLI 输出面。 |
| (c) 诊断包里的一项 | **不选：越界且无落点** | `internal/observe/**` 是本程零字节禁改面（派单 §2；票面 §Packages 另加"注意地界"）。且 `diagnostics.go` 现量 `compression` **0 命中**、它的字面量键名只有 `"id" "ticket" "summary" "registered" "reason"` —— 那是一张**检查项**表，不是 agent 侧观测的收件处，硬塞要先动它的形状。 |
| (a) 结构化日志字段 | **选它** | 三条支撑读数：① 它是**唯一已经落到磁盘**的那条路——`installLogSink` 把 `slog.Default()` 换成"红删 JSONL 文件 + stderr"的 tee（`cmd/wisp/logsink.go:144`），门槛 `logSinkLevel = "info"`（`:82`）⇒ Info 级真的写进 `<data>\logs\wisp-*.jsonl`；② `cmd/wisp/**` 非测试代码**从不设 `agent.Options.Logger`**（`git grep -n "Logger:" -- 'cmd/wisp/*.go' ':!*_test.go'` rc=1）⇒ loop 用的就是那枚持久 sink，不需要为新痕开任何线；③ 落点在 Go 侧，不碰 SPEC-08 §5.2 的枚举白名单 ⇒ 不是契约变更。 |

### 2.2 为什么写在 `compress.go` 而不是 `loop.go` 的成功分支

票面 §0 自己把这一格挖深了一寸：那条失败 Warn **在调用方**，`compress.go` 从头到尾不产生日志。
痕落在 `Compress` 内部换来三件调用方给不了的事：
1. **未来的调用方自动继承**——D28-1 那句"Warm-window hook 才该拥有这次调用"落地那天，不必记得补打；
2. **与失败那一发同前缀**：`agent: history compressed` / `agent: history compression failed`
   共享搜索前缀 `agent: history compress`，一次 grep 拿到两种结局；
3. **"没折叠就不打"由数据决定**：`rep.Ran` 是本函数内的折叠事实，不是调用方事后推断的条件。

代价（照实记）：**痕里没有 taskID**。`Compress` 的签名不带任务号，加一枚参数要动 `Summarizer` 之外
的第二条构造线，且会牵动 `compress_test.go` 全部 8 枚调用点。⇒ 归因靠同一次 run 里相邻的
`task` 字段（失败 Warn 亦不带 task，二者对称）。这一条落进 §2.5「没测什么」第 1 项，并在 §1.4 的
五字段"当前残缺表现"里写了同一句话。

### 2.3 码的形状（`git show --stat 1acdd03`）

```
internal/agent/compress.go            |  +CompressorOpt/WithLogger（变参）+ log() + 成功边一条 Info
internal/agent/loop.go                |  1 行：New() 把 opt.Logger 交给 NewCompressor
internal/agent/compress_trace_test.go |  新增，4 枚用例
internal/agent/harness_test.go        |  新增 withLogger（纯追加，默认仍 discard）
```
构造器选**变参 opt** 而不是加第 3 枚必填参数，是为了让 8 枚既有 `NewCompressor(b, sum)` 调用点
（`compress_test.go` 7 枚 + `loop.go` 1 枚）**一字不改**：`git grep -n "NewCompressor(" -- '*.go'`
在改动后仍全部编译，`go build ./internal/agent/` rc=0。
"helper 是不是原有的那枚"自答：`harness_test.go` 里既有的 `withConfig/withTools/withRegistry` 等**语义未动**，
`o.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))` 那行默认值也**未动**，只是多了一枚 opt。

### 2.4 记录本身（真发一条出来长什么样，非人工排版）

`TestCompressionTraceDoesNotAlterTheFold` 里那枚未注入 logger 的 `NewCompressor(b, nil)` 走的是
`slog.Default()` 回退支，故 `-v` 输出里能看到原样一行（`.scratch/wisp/probes/139/mut/m0.log`，摘）：
```
time=2026-09-25T23:45:23.573+08:00 level=INFO msg="agent: history compressed" \
tokens_before=1236 tokens_after=780 threshold=384 msgs_before=18 msgs_after=10 \
compressed_msgs=9 kept_raw_rounds=3 history_changed=true
```
⚠ 这一行同时给出一个**读数级**的事实：`tokens_after=780` 仍**高于** `threshold=384`。
⇒ D15(4) 的退出条件是"预算够了 **或** 只剩 `KeepRawRounds` 枚原始轮"，两者都能收工，
所以"压后 token"合法地可以停在阈值之上。**这就是记录里必须带 `threshold` 的理由**：
没有它，读到 `784`（这里 780）的人无从判断"压够了"还是"撞到轮数地板"。
我第一版断言写的正是"`out` 必须不再 `Need()`"，它在真码上是**错判据**，跑出来第一发就红：
`trace booked 1284 -> 804 against threshold 384, but the output still needs compressing (804 tokens)`
—— 这条红是测试写错、不是码写错，改成了阈值感知的形状（`c.Need(out) && rep.KeptRawRounds > b.KeepRawRounds` 才红）。

### 2.5 放水两问自答（AC#2 这一格）

- **这一发在未修码上响不响？** 本格的"痕"本身没有判据（判据在 §3），但它的**不存在**被 §1.1 的尺量过：
  同一枚 `Info\(|Warn\(|Debug\(|Error\(` 正控在 `loop.go` 给 14、在 `compress.go` 给 0。
- **摘掉本票这一味，是否存在一发变异从此打不红？** 见 §3 的 M1：摘掉这条 `Info` 后有**两发**转红。不是装饰。

### 2.6 本程没测什么

1. **没测这条痕能回答"是哪个任务压的"**（无 taskID，见 §2.2 代价）。
2. **没测红删 handler 会不会改我的字段名**：`observe.redactHandler` 在写盘前会处理路径类内容，
   本程只读了它存在（`logging.go:132`），没量过一次真写盘后 JSONL 里的键名。
3. **没测 `LevelInfo` 在别的宿主入口被调低的可能**：`logSinkLevel` 是 `cmd/wisp` 的常量，不归本程核。

---

## 3. AC#3 — 牙：摘掉痕必须转红（四枚变异，逐名报红行）

**判据**（票面）：一枚用例钉住"压缩成功 ⇒ 痕必须存在"；反向判据是把痕摘掉，该用例必须转红并**逐名报出红行**；
先证变异落地（`grep -n` 原文 + `go build` rc=0）再读数；不许 `t.Skip`、不许放宽断言。

四枚变异全部用脚本施加，每枚跑完从 `probes/139/mut/*.orig` 还原并 `cmp` 自证（`RESTORED OK`）。
原文与逐发读数：`.scratch/wisp/probes/139/mutations.log`、`m4-only.log`、`mut/m0..m4.log`。

| 变异 | 施加与"落地证明" | build | 红名（逐名） | 判 |
|---|---|---|---|---|
| **M0** 基线（未变异） | — | rc=0 | 无 | 四枚全绿（`test_rc=0`） |
| **M1 摘痕**：删掉 `if rep.Ran { … }` 整块（＝票面点名的那一味） | `grep -q 'c.log().Info'` 无命中才算落地；diff 只减 16 行 | rc=0 | `TestCompressionTraceBooksCountsOnSuccess`<br>`TestCompressionTraceSurvivesTheLoopWiring` | **两发红** ⇒ 痕是承重的 |
| **M2 只断接线**：`loop.go` 去掉 `WithLogger(opt.Logger)` | diff 1 行 | rc=0 | `TestCompressionTraceSurvivesTheLoopWiring` | **只红这一发** ⇒ 端到端那发确实在测"接线"，不是重复测 M1 |
| **M3 恒真化**：`if rep.Ran {` → `if true {` | `grep -q 'if rep.Ran {'` 无命中 | rc=0 | `TestCompressionTraceSilentWhenNothingFolded` | 反向那一发**独立咬住**恒真判据 |
| **M4 泄内容**：记录里追加 `"peek", fmt.Sprint(out[0].Content[0])` | diff 1 行 | rc=0 | `TestCompressionTraceBooksCountsOnSuccess` | 隐私探针**咬得住** |

**红句原文**（这四句就是"改前那一发"的形状：M1 与 HEAD 版 `compress.go` 的唯一差异就是这味药）：
```
M1  compress_trace_test.go:206: "agent: history compressed" records = 0, want exactly 1 (all records: [trace-capture-ruler-control])
M1  compress_trace_test.go:328: a real Run that folded history left no "agent: history compressed" record; captured: [trace-capture-ruler-control]
M2  compress_trace_test.go:328: a real Run that folded history left no "agent: history compressed" record; captured: [trace-capture-ruler-control]
M3  compress_trace_test.go:286: no-op pass left 1 trace record(s), want none: [agent: history compressed]
```
⇒ **放水第一问**（未修码上响不响）：M1 = 未修码，**响**，两发同时红。
注意每条红句里都印着 `captured: [trace-capture-ruler-control]`——那是 §3.1 的正控记录，
它证明"0 枚痕"是**尺读到的 0**，不是尺本身没接上。

⚠ **诚实标注一处不等价**：M1 不是字面意义的"HEAD 版跑同一发用例"——HEAD 版里那四枚用例还不存在，
所以任何"新码+新用例"票的"改前红"都只能是**行为等价的模拟**。此处等价性可核：M1 的 diff 只删
`compress.go` 里我加的 16 行，`rep.Ran` 之外无别的行为改动，且 `harness_test.go` 的默认 discard 未动。

### 3.1 正控先行（本仓对"零读数"的规矩）

四枚用例都在断言痕之前先调 `recs.assertRulerLive(t)`：往同一枚 handler 打一条
`trace-capture-ruler-control`，打不进 capture 就 `t.Fatalf`。M1/M2 的红句里能看到这条控制记录仍在场。
另两枚已知正控在 §1.1（grep 尺：`.Rounds`=3 / `.Usage`=54 / `loop.go` 日志 14 枚）。

### 3.2 隐私探针为什么要放两枚（M4 教出来的）

M4 那次泄漏**不是被哨兵字符串抓住的**：`structuralTrim` 把每条正文截到 80 rune，
我埋在正文尾部的 `PRIVACY-SENTINEL-do-not-log-this` 被截掉了，所以第一枚探针没响；
响的是第二枚——查记录里是否出现正文里必然存在的 `问题`。红句原文（`m4-only.log`）：
```
compress_trace_test.go:256: trace leaked a user turn body: "agent: history compressed …
peek={[已压缩的历史] 以下轮次已被摘要压缩，仅供上下文参考：\nuser: 问题 0 xxxxx…
```
⇒ 记一条对本件自身的读数：**"往正文尾部埋哨兵"这一形对截断型摘要无效**，两枚探针不能省成一枚。

### 3.3 放水两问自答（AC#3 这一格）

- **断言方向动没动**：没动任何既有断言。`compress_test.go` 4 枚既有用例一字未改（`git show --stat 1acdd03` 里它不在改动清单）；
  新增用例只加不改。全仓 `t.Skip` 命中数：见 §4。
- **helper 是不是原有的那枚**：既有用例用的仍是原有 `buildRoundHistory` / `BudgetsFor` / `newHarness`（默认 discard 未动）；
  新加的 `traceCapture` 是本件自己的尺，并且每次先被正控驱动。

### 3.4 本程没测什么

1. **没测"痕重复打"**：`len(hits) != 1` 在端到端用例里是**排在功能断言之后的最后一条**，
   若某次改动让一轮里打了 3 条痕，会先被前面的断言放过、由这一条抓到；但**多轮压缩**（一轮一次）没测——
   那需要一条多轮的 golden fixture，而**任何 golden 在零字节禁改面上**，本程不新增。
2. **没测 `Compress` 失败路径的痕**：失败仍是调用方那枚 Warn（M1 之前也是），本票没给失败路径加码，
   所以"成功一条 Info、失败一条 Warn"的**对称性只由源码前缀保证，没有用例钉**。
3. **没测日志被降级到 `LevelWarn` 之后痕是否还在**：`Enabled` 由宿主 handler 决定，不在本包。

---

## 4. AC#4 — 门禁（逐包单跑，两形四数 + 名册差集 + 三道自核）

原文读数：`.scratch/wisp/probes/139/gate-ac4.log`（含命令），日志件 `gate-pre-full.log` /
`gate-post-count1.log` / `gate-post-full.log`。**口径承 138 第 6 格**（两形＝`-count=1` 与 `-count=2`）。

### 4.1 四数（`internal/agent` 单包）

| 形 | 采于 | RUN(all) | PASS_top | PASS_all | FAIL_any | SKIP_any | unique | panic: |
|---|---|---|---|---|---|---|---|---|
| 改前 `-count=2` | 进场时 HEAD `3f6322f`，本程未改码之前 | 152 | 118 | 152 | **0** | **0** | 76 | 0 |
| 改后 `-count=1` | `1acdd03`（码已入库） | 80 | 63 | 80 | **0** | **0** | 80 | 0 |
| 改后 `-count=2` | 同上 | 160 | 126 | 160 | **0** | **0** | 80 | 0 |

先说破两枚坑：**`PASS_top` 只数顶格**（子用例缩进，59/118 那类），所以顶格 118 ≠ 全量 152；
**FAIL 只认 `--- FAIL`**，包级汇总行 `ok  … 5.182s` 不进计数（本仓实测过汇总行会造出假 FAIL）。

### 4.2 三道自核（一枚 panic 会吞掉同包几十条读数，四数看不出来）

1. **算术自洽**：`unique(-count=1) = 80`、`RUN_all(-count=2) = 160 = 2 × 80`；`126 = 2 × 63` ⇒ 没有第二遍漏跑。
2. **名册两向 `comm`**：改前 76 → 改后 80。`comm -23`（新增）＝本票那 4 枚，逐名：
   `TestCompressionTraceBooksCountsOnSuccess` / `TestCompressionTraceDoesNotAlterTheFold` /
   `TestCompressionTraceSilentWhenNothingFolded` / `TestCompressionTraceSurvivesTheLoopWiring`；
   `comm -13`（丢失）**空**。⇒ 两向差集只有"我加的"，没有"被顶掉的"。
   两形名册互差亦为空（`diff n-post1.txt n-post.txt` 无输出）。
3. **同包"开跑 vs 出裁决"等集**：`comm -3` 空，`ran=80 / verdicted=80`，且三份日志 `grep -c 'panic:'` 全 **0**
   ⇒ 没有"开跑了没回来"的用例（`[0]` 直取那一族的病在这里没有发作）。

### 4.3 其余工具（版本现读，不背别人的读数）

```
$ go vet ./internal/agent/                       -> 无输出，rc=0
$ gofmt -l internal/agent                        -> 空（rc=0）
$ D:/work/base/gopath/bin/gofumpt.exe --version  -> v0.12.0 (go1.27.1)
$ gofumpt.exe -l . tools/d22scan tools/mockllm   -> 空（rc=0）
$ sh scripts/d22scan.sh                          -> rc=0，八枚分母全非零：
   bans #1-5 internal/=205  cmd/=23  ban #6 frontend/=66  ban #7 internal/tools/=18
   ban #8 design/=39  frontend/=66  internal/=413  cmd/=43
```
本轮 `gofumpt` 的清单里**没有** `cmd/wisp/**` 命中（派单预告过可能撞 149 的中间态；采数时它不在飞）。

**"各 scope 不降"怎么证的**（不拿别家锚点的数字冒充）：`internal/` 的 `.go` 枚数在
锚点 `3f6322f` = **412**、当前 HEAD 跟踪的仍是 **412**、工作树（含本程新增那枚未跟踪测试文件）= **413**
＝ ban #8 `internal/` 的分母 413。⇒ 分母只 +1，且那 +1 是本程自己的文件
（`probes/139/d22scan-denominators.log`）。bans #1-5 的 205 枚是"production Go files"，
本程新增的是 `_test.go`（不在其分母里），故不降。

### 4.4 放水自证

- **没有 `t.Skip`**：本包（`./internal/agent/`）零命中。全树 `git grep -n 't.Skip' -- 'internal/agent/**'`
  只有 **1 枚**，在子包 `internal/agent/approval/ticket84_no_owner_test.go:224`（票 84 既有的"有意慢"闸，
  与本票无关、本程未动）。
- **没有放宽断言**：`compress_test.go` 一字未改（`git diff --stat 3f6322f..HEAD -- internal/agent` 的清单里
  只有 `compress.go +60` / `compress_trace_test.go +393` / `harness_test.go +8` / `loop.go 1±`，
  **`compress_test.go` 不在列**）；既有的 `NewCompressor(b, sum)` 8 枚调用点一字未动。
- **阈值与 golden 零字节**：`git diff --name-only 3f6322f..HEAD -- internal/observe internal/agent/thresholds.go`
  → **空**；本程未新增/未修改任何 `testdata/golden/*.sse`（要造多轮压缩场景需要新 fixture，那在禁改面上，
  本程因此**没做**多轮痕的用例，见 §3.4 第 1 项）。

### 4.5 锚点漂移（本程所见如实记，不改写别人的话）

进场 `git rev-parse HEAD` = **`3f6322f`**（与派单锚点相同）。交件时 `3f6322f..HEAD` 之间共 **11 枚** commit，
其中本程 4 枚（`4d0866a` 格1 / `1acdd03` 码 / `e22da5a` 格2-3 / 本格随后），**别家 7 枚**
（138 验收程 3 枚、149 验收程 2 枚、台账 1 枚、`feat(frontend …)` 1 枚）。
逐名核过：**没有任何一枚别家 commit 触碰 `internal/` 或 `cmd/`**
（对 `3f6322f..HEAD` 的每枚非本程 commit 跑 `git show --name-only -- internal cmd` ⇒ 全部空，见 §4.5 命令）。
⇒ 本程的"改前"基线与"改后"读数额之间没有别人的代码插进来，两发的可比性成立。
采数过程中亲眼见到的漂移：23:47 `424ac84` → 23:49 `cce7510` → 23:53 `e22da5a`（`gate-*.log` 抬头各自带当时的 HEAD）。

### 4.6 本程没测什么（本格）

1. **没测 `-race`**：138/149 的 AC#4 口径里也没有它；本程新增的 `traceCapture` 自带 mutex，
   但**没有一次并发写它的实测**（四枚用例都在测试 goroutine 上串行调 `Compress`／同步 `h.run`）。
2. **没测跨包门禁**（`go test ./...` 全树）：派单要的是逐包单跑；全树那发由编排者与 CI 采，
   本程不冒充它的颜色。
3. **没测 `internal/observe/**` 与 `*_other_test.go` 那一族**：本程零改动那棵，故未在 linux 容器里量任何东西。
   ⚠ 顺带把派单 §5 那两句回核了一遍：本机 `docker` 可用性与 `golang:1.27` 是否在本地**本程未量**
   （用不上就没必要为它花轮次），Git Bash `docker -v C:\…` 假绿那一句是**编排者的断言**、本程既未采信也未推翻。

---

## 5. 收口

### 5.1 四格状态与 AC 对账

| 格 | AC | 落点 | 状态 |
|---|---|---|---|
| 1 | AC#1 | 本文 §1 | 读数齐（含两条未复算项落地）；`SPEC-12 §5` 补登记**未落笔**（禁改面），五字段备妥在 §1.4 |
| 2 | AC#2 | 本文 §2 + `1acdd03` | 选了结构化日志，落 `Compress` 成功边；三条内容要求齐（压前/压后/是否动了历史），无原文内容、无前端事件、无面板方法 |
| 3 | AC#3 | 本文 §3 + 4 枚用例 | 四枚变异 M1–M4 逐名红句，M1＝票面点名的"摘掉痕" |
| 4 | AC#4 | 本文 §4 | 两形四数 + 名册两向差集 + 三道自核 + 三件工具 |
| — | §3 明确不做 | 面板分解条 / 点开看原文 / 压缩算法与 `Need()` 阈值 / D28-1 hook | **一律未做**（票面 §3 与本程地界同向） |

**碰到即停的两条本票都没碰上**：未新增任何 Go→前端事件或 `C17` 面板方法；未改压缩算法/阈值/D15 预算分配。

### 5.2 票面与本程简报冲突处（按派单 §0 报回，不硬改）

**无实质冲突。**两处措辞差异照实记：
- 票面 §Packages 说 `internal/observe/**` "此刻可能有别的程在跑" ⇒ 本程现量**无人**（`git status --porcelain`
  那棵 0 枚，`probes/139/git-state.log`），且本程最终**一行都没碰它**。
- 票面 AC#1③ 说未覆盖就"去 `SPEC-12 §5` 按五字段补登记" ⇒ 与派单 §2 的"禁改面（零字节）`docs/specs/**`"相撞。
  本程按**更严的那条**（禁改面）执行：**登记文本备妥、不落笔**，等编排者落。

### 5.3 伪授权登记（两数分栏）

| 栏 | 数 | 逐条（带出处） |
|---|---|---|
| 真通知回显 | **2** | ① 后台命令完成通知（`gate-pre-full.log` 那发 `go test` 的 task-notification，id `b4l0dbn6v`）；② `MEMORY.md` 被外部改动的系统提示（自动事件，非指令，本程未据此行动） |
| 判为注入 | **0** | 本程未在任何工具输出里遇到"像编排者说的话"。相关的一条：`MEMORY.md` 里"安全告警默认只到 stderr（持久 sink 只在 `wisp slo`）"与 `cmd/wisp/logsink.go` 现码**不符**（现码 `installLogSink` 是 JSONL+stderr 双路 tee）⇒ 判为**过期记忆**、按现码走，已在 §1.6 第 1 项记名，不属伪授权 |

### 5.4 凭据/密钥

本票全程未接触凭据面。证据件与探针里出现的只有变量名与文件名
（`opt.Logger`、`logSinkLevel`、`<data>\logs\wisp-*.jsonl`），零字值。

### 5.5 本程整体没测什么（跨格，按"漏了它谁会先被骗"排序）

1. **"这条痕真的落到过一次磁盘"**：§1.6 第 1 项那串环环推导没有一发端到端实测（起进程→跑一次会压缩的任务→翻 JSONL）。
   最先被骗的是把"有 log"读成"事后答得出"的人，包括本件的读者。
2. **痕不带 taskID**（§2.2）：多任务并发时"哪次压缩属于哪个任务"答不出来。
3. **多轮/多次压缩**只测了"一轮恰好一条"（§3.4 第 1 项）；多轮的形状需要一枚新 golden，撞禁改面。
4. **失败侧对称性无钉**（§3.4 第 2 项）：`agent: history compression failed` 那枚 Warn 今天**没有**任何用例断言它存在，
   所以"两种结局共享前缀"这句话目前只有源码担保。
5. **`redactHandler` 对本记录键名的影响未量**（§2.6 第 2 项）。
6. **未复核派单 §5 的容器两句**（§4.6 第 3 项）：那两句是编排者的断言，本程按"未采信亦未推翻"处理。

*（本件出自实现者，不自勾。四格判据的每一枚读数都带原命令，可独立复算。）*
