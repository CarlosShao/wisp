# 139 — 压缩留痕 · 对抗验收件 r1（验收程，非实现者）

被验版本＝**`d949d9c5`**（票 139 四枚：`4d0866ad` / `1acdd032` / `e22da5a8` / `d949d9c5`）。
实现件：`docs/evidence/s1/139-compression-leaves-no-trace-r1.md`。票面：
`.scratch/wisp/issues/139-context-compression-leaves-no-trace-zero-log-calls-in-compress-go-zero-readers-of-res-compression.md`。

**本件是裁决，不是实现**。票面四框**本程不勾、也不替实现方勾**。零代码改动（见每格「本格改了哪些文件」）。

---

## 0. 锚点与工作树现量（开工第一发）

```
$ git rev-parse HEAD
2ebc119e781850d9a4cadfc37a4dc537a3b31595          # ≠ 被验版本 d949d9c5（其后另有 138/149 两枚验收程在交件）
$ git log --oneline -10                            # 摘要
2ebc119 evidence(149 验收 r1 第 2 格) …
39669af evidence(138 验收 r1 第 3 格) …
d949d9c evidence(139 AC#4 门禁 第 4 格 收口)       ← 被验版本
e22da5a evidence(139 AC#2+AC#3 第 2-3 格)
85c2cd7 evidence(149 验收 r1 第 0-1 格)
1acdd03 fix(139 AC#2+AC#3): 压缩成功那一跳留一条只记计数的痕
cce7510 evidence(138 验收 r1 第 2 格)
424ac84 feat(frontend beautiful-ui 转向)
$ git status --porcelain | head                    # 脏的全是别家的
 D design/assets/base.css … M design/doubao/demo/* … M docs/evidence/s1/138-cancel-kills-children-r1-accept-r1.md
?? .zcodeignore ?? design/doubao/01-ball-states.jpg ?? design/old/
```

⇒ 工作树确实脏（前端会话在写 `design/**`、138 验收件在飞）。**本程所有"某文件长什么样"一律按 sha 取**：
读源码用 `git show d949d9c5:<path>`，要跑就 `git archive d949d9c5` 到仓外快照。

**⚠ 派单 §0 那条快照坑：本程亲手撞上并复核成立。**
```
$ git -c core.autocrlf=false -c core.eol=lf archive d949d9c5 | tar -x -C /tmp/wisp139-snap-1   # rc=0
$ ls /tmp/wisp139-snap-1/third_party        → No such file or directory
$ git ls-files third_party | wc -l          → 0                  # 原生 dll 从来就不在树里
$ PATH=… go test ./cmd/wisp/ -run TestAccept139
   exit status 0xc0000135                                            # STATUS_DLL_NOT_FOUND
   FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.125s
```
`0xc0000135` ＝ 进程连加载都没完成，**不是被验版本的任何行为**。修法＝把宿主 `third_party/sherpa-onnx/*.dll`
（3 枚：onnxruntime / sherpa-onnx-c-api / sherpa-onnx-cxx-api）拷进快照并挂 PATH，之后同一条命令真的跑起来了。
⇒ 派单这句成立（且比它写的更具体：红形是 loader 级的 `0xc0000135`，不是 `no native DLLs` 那句人话）。

**本格改了哪些文件**：无（只读）。探针与快照在仓外 `/tmp/wisp139-snap-1`，副本落
`.scratch/wisp/probes/139/accept-r1/`（只建不删；未纳入 commit，见派单 §6「只 commit 自己那一枚件」）。

---

## 1. 第 2 节首要攻击点：这条痕到不到得了盘 —— **量成了读数，而读数把它推翻**

### 1.1 先复核派单 §2 那四行 grep（**四行全部按字面成立**）

```
$ git show d949d9c5:cmd/wisp/logsink.go  | grep -n 'slog.SetDefault'   → 153:	slog.SetDefault(slog.New(teeHandler{
$ git show d949d9c5:cmd/wisp/run.go      | grep -n 'installLogSink('   → 178:	sink, sinkErr := installLogSink(s.dataDir)
$ git show d949d9c5:cmd/wisp/run.go      | grep -n 'agent\.New('       → 547:	loop, err := agent.New(agent.Options{
$ git show d949d9c5:internal/agent/loop.go | grep -n -A2 'if opt.Logger == nil'
    207:	if opt.Logger == nil {   208:		opt.Logger = slog.Default()
```

派单 §2 那句「这只是一枚按行号的顺序推断」**也必须照字面收下**：`installLogSink`（:178）在 `runTextTask`（:137）里，
`agent.New`（:547）在 `(*agentRuntime).execute`（:542）里，**不同函数**，行号本身推不出先后。
真先后是调用图给的（同一条 goroutine，自上而下）：

```
main.go:129 → runTextTask → :178 installLogSink → :190 assembleRuntime → :203 rt.execute → run.go:547 agent.New
```

并且**全仓非测试代码里 `slog.SetDefault` 只有两枚调用者**（这把"中间有没有人换掉默认值"量掉了）：
`cmd/wisp/logsink.go:153`（run 腿装的 tee）· `internal/observe/logging.go:107 InstallAsDefault`
——后者全仓唯一调用者是 `cmd/wisp/slo_windows.go:279`（`wisp slo` 那条**另外的**腿，不与 run 同进程）。
`logging.go:465` 那枚在 `func init()` 里，**先于** :153，只可能被覆盖、不可能覆盖它。
⇒ 派单"取到的已是文件 tee"这句**在 run 腿上成立**，且现在是有读数的成立（下面 1.3）。

### 1.2 「有没有真测」——**没有，且实现件自报的这句成立**

`internal/agent/compress_trace_test.go` 四枚用例里最"端到端"的那枚 `TestCompressionTraceSurvivesTheLoopWiring`
（`:296`）走的是 `newHarness(..., withLogger(recs.logger()))`，而 `withLogger` 就是
`harness_test.go:88 o.Logger = lg` —— **它注入的是测试自己的 handler，压根不经过 `slog.Default()`、不经过 `installLogSink`、不落到盘**。
⇒ 它测的是"loop 把 opt.Logger 交给了压缩器"（M2 就是这一枚红），**不是**"这条记录到得了文件"。
实现件 §1.6 第 1 项 / §5.5 第 1 项自报"端到端一次没量"——**这句是真话，本程没有为它翻案**。

### 1.3 本程造的端到端读数（一）：**sink 那条链是通的**（把四段推理换成一发实测）

造法：仓外快照 `d949d9c5` ＋ 真 dll，在快照的 `cmd/wisp` 里加一枚**只存在于快照**的用例
（探针全文：`.scratch/wisp/probes/139/accept-r1/leg-e2e-probe_test.go.txt`），它做三件事：
① 用 `cmd/wisp/run_test.go:76 newRunFixture` 同一条**真组合根**跑 `runTextTask`（不是子进程、但同一条构造路径，
与票 131 那枚 nail 的"leg entry = runTextTask"同一形状）；② 把 provider 换成本程自己的脚本化 SSE 服务，
每轮吐 16 000 字 assistant 正文（＝len/4 下 4 000 token，且正文不受 spill 层管）；
③ 跑完打开 `<data>\logs\wisp-*.jsonl` 逐行解析。`context_window = 32768` ⇒ 阈值 12000×0.25＝**3072**（读数里就是这个数）。

为了能"一次读数"看到压缩器为什么不出声，本程在**快照**的 `loop.go` 里临时加了一行 `l.log().Info("accept139-probe-need", …)`
——**用的是 loop 自己那枚在 `agent.New` 时刻从 `slog.Default()` 抓来的 logger**，所以这一行本身就是那条链的试纸。
```
$ PATH="/tmp/wisp139-snap-1/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp/ -run TestAccept139 -count=1
    PROBE: {"conv_rounds":1,"level":"INFO","msg":"accept139-probe-need","msgs":1,"raw_rounds":1,"rounds":1,"threshold":3072,"tokens":8}
    PROBE: {"conv_rounds":1,…,"msgs":3, …,"raw_rounds":1,"rounds":2, …,"threshold":3072,"tokens":4048}
    …（每轮一条，直到）…
    PROBE: {"conv_rounds":1,…,"msgs":21,"raw_rounds":1,"rounds":11,"threshold":3072,"tokens":40408}
```
⇒ **11 条经 `agent.New` 抓去的 logger 出来的记录，全部落进了那个 JSONL 文件**（`sink records = 31`，逐枚可解析、键名未被红删改掉）。
这条读数同时结掉三件事：
- 派单 §2 的"`slog.Default()` 取到的已是文件 tee"——**不再需要按行号推，量到了**；
- 实现件 §2.1 支撑① "`installLogSink` 是唯一已经落到磁盘的那条路"——**方向成立**；
- 顺带把实现件 §2.6 第 2 项"`redactHandler` 会不会改我的字段名"也答了：**不改**（本程这条探针记录带 7 枚属性键，逐个原名在场）。

**快照还原自证**（探针不许污染被验版本）：
```
$ cp /tmp/loop139.orig.go /tmp/wisp139-snap-1/internal/agent/loop.go
$ grep -c 'accept139-probe' internal/agent/loop.go → 0
$ cmp internal/agent/loop.go <(git show d949d9c5:internal/agent/loop.go) → RESTORED IDENTICAL
```

### 1.4 本程造的端到端读数（二）：**票 139 那条痕在真腿上一发都没打过**（推翻点）

同一次跑里，本程找它交付的那条记录：
```
    zz_accept139_e2e_test.go:247: THE TRACE NEVER REACHED THE FILE:
        "agent: history compressed" records = 0 in the real `wisp run` log sink   （当时盘上共 20 条 / 加探针后 31 条）
$ … | grep -oE 'compression failed|history compress[a-z]*|probe-ran' | sort | uniq -c
      1 history compressed          # ← 这一枚是**本程失败句自己的引号文本**，不是盘上记录
                                    #   盘上 msg 名册里 `agent: history compression failed` = 0、
                                    #   `accept139-probe-ran` = 0（else-if 从未到达）
```
**为什么 0？不是链断，是那枚 `if rep.Ran` 在真腿上永远为假。** 上面读数里 `raw_rounds` 恒为 **1**：
`groupRounds`（`compress.go:291`）以 `m.Role == llm.RoleUser` 开一轮，而 `wisp run` 的历史
就是一个任务 ＝ **一枚 user 消息**（`loop.go:371 l.append(userMessage(input))`），后面全挂在同一轮里；
`Compress` 的折叠条件是 `if len(raw) <= c.b.KeepRawRounds { break }`（`compress.go:136-138`，`KeepRawRounds = 3`，
`budgets.go:42`）⇒ **1 轮历史里连一枚可折的"最老轮"都取不出来**，`rep.Ran` 保持 false，
成功边那条 `Info` 一次都不执行。token 明明早就过阈值（4048 → 40408 对 3072），**折叠就是不发生**。

**这个 `raw_rounds = 1` 不是本程构造出来的孤立形状，而是今天所有生产腿的共同形状**：
```
$ git grep -n '\.Steer(\|RunAsync(' d949d9c5 -- '*.go' ':!*_test.go'
    internal/agent/loop.go:321:func (l *Loop) RunAsync(…)      # 只有声明，零非测试调用者
    （.Steer( 非测试命中：0）
$ git grep -n 'agent\.New(' d949d9c5 -- '*.go' ':!*_test.go'
    cmd/wisp/run.go:547                                          # 全仓唯一非测试构造点
```
⇒ 没有任何一条在跑的腿会往同一枚 Loop 里塞第 2 枚 user 消息（`Steer` 零调用者、`RunAsync` 零调用者、
Loop 不跨任务复用）。实现件 §3 那枚"端到端"用例之所以能红，是因为它用
`compress_test.go:23 buildRoundHistory(4, 400)` **手工种了 4 枚 user 轮** ——
那是**只有测试会造的历史形状**。

### 1.5 另外两条腿（派单 §2 第 2 问）：**问题不在顺序，在根本不经过 `agent.New`**

`installLogSink` 共 4 处非测试调用：`run.go:178` · `resident_windows.go:57` · `models.go:284` · `secret.go:226`。
逐条读源码后的读数：
- `runResident`（`resident_windows.go:24-84`）＝ `proc.Boot` → `installLogSink` → `rt.RunEventLoop()`，
  **全程不构造 Loop**（`agent.New` 全仓唯一非测试调用者是 `run.go:547`）；它自己那行注释写着
  "the resident process is the leg owner actually uses"。
- `cmdModels`（`models.go:284`）/ `cmdSecret`（`secret.go:226`）同理，只装 sink、不跑 agent。
⇒ 派单猜的"哪条腿顺序反过来"——**没有一条反过来**；但更硬的是：**除了 run 这条腿，别的腿连这条痕的生产者都没有**，
而 run 这条腿本身（见 1.4）今天造不出这枚痕。所以"到得了盘"这句话覆盖的不是四条腿，也不是一条腿，
**是零条腿**：今天没有任何一条在跑的腿能产生这条记录。

### 1.6 档位与影响（这一格）

| 判据 | 档位 | 理由 |
|---|---|---|
| 实现件 §1.6#1 / §5.5#1 自报"端到端没量" | **成立**（真话，本程未翻案） | 四枚用例全在自己的 handler 上，无一过盘 |
| 派单 §2"取到的已是文件 tee" | **成立并升级为有读数**（§1.3） | 11 枚同链记录实落到 JSONL |
| 实现件 §2.1 支撑① "(a) 是唯一已落盘那条路" | **附条件入账**：机制成立、**这枚痕到不了那条路** | §1.4 |
| 票面 AC#2 的目标句"事后答得出压了多少" | **不成立（今天答不出）** | 生产侧零条记录；票面 §1 立的靶"根因是根本没有数据可看"**未被搬掉** |

**这不是 139 写坏了码**——`rep.Ran` 为假是 `KeepRawRounds=3` 对 1 轮历史的既定行为，
且票面 §3 明令"压缩算法本身、`Need()` 的阈值、D15 的预算分配：一律不改（改这些＝契约变更）"。
⇒ 本程**不提议改压缩器**（那要人工批准）；要报的是**这一格的正确判据形状缺一次读数，而那次读数量出来的是阴性**。

**本程没测什么（本格，按"漏了它谁会先被骗"排序）**
1. **没量"多轮会话"那条形状**：若哪天有腿复用同一枚 Loop 跨任务（或 `Steer` 真被接线），
   `raw_rounds` 才会 >1，本格的阴性结论就要重读。本程只证明"今天没有这种调用者"，
   没有证明"这种调用者一旦出现，痕就对"——那还需要一发盘上检。
2. **没起真子进程**：本程用 `runTextTask` 同函数入口（票 131 的 nail 也是这个形状），
   不是 `exec.Command` 起 `wisp.exe`。差异面＝`os.Args`/信号/stdout 重定向，本格结论不依赖它们，但"真进程"那一发仍未量。
3. **没量子进程版的 resident 腿**（票 127/131 有，本程未复用其形状）。
4. **`-race` 下这条链未量**（实现件 §4.6#1 同样未量）。
5. **没量"痕重复打"的盘上形状**（本程只数到 0 枚，无法在同一发里判 1 枚 vs 3 枚）。

---

## 2. 第 4 节承重与恒真：四发变异独立重跑 ＋ 本程自造四发

跑在**第二枚仓外快照** `/tmp/wisp139-snap-2`（`git archive d949d9c5`，与锚点逐枚 md5 相符：
`compress.go 9c5015dd… / loop.go 70f8fb10… / compress_trace_test.go 7bf2fbf7…`）。
驱动件 `mut.sh` / `drive.sh` / 逐发 patch `.py` 都在该快照内（快照非交付物，形状见 §0 落点规矩）。
每发都：先证落地（`grep` 原文改动在场）→ `go build ./internal/agent/` rc=0 → 读数 → 还原 → `cmp` 自证。

### 2.1 实现件自报的四发：**逐名复算全部成立**

```
[M0]  build rc=0 | whole-pkg-reds: NONE                                    （基线全绿）
[M1b] build rc=0 | whole-pkg-reds: BooksCountsOnSuccess SurvivesTheLoopWiring
[M2]  build rc=0 | whole-pkg-reds: SurvivesTheLoopWiring
[M3]  build rc=0 | whole-pkg-reds: SilentWhenNothingFolded
[M4]  build rc=0 | whole-pkg-reds: BooksCountsOnSuccess
```
（M1b ＝ 真删 `if rep.Ran { Info(…) }` 整块，比实现件的 `if false` 形状更狠：不留死码。）
⇒ 实现件 §3 那张表的**红名逐枚对得上**，且 M2 只端到端那一枚红、M3 被反向那枚独立咬住，两句都真。
名册口径按派单 §5：只数 `--- FAIL:` 行、`-count=1`、且这里跑的是**整包**（不只那四枚），
所以"别枚用例顺手也抓到了"这一支被排掉了——见 2.2 的 M6b/M7，整包里真出现了别枚红。

### 2.2 本程自造的四发（实现件没问过的问题）

| 变异 | 形状 | 整包红名（`-count=1`） | 判 |
|---|---|---|---|
| **M5 痕恒打（现实形）** | `if rep.Ran {` → `if c.Need(hist) {` | **NONE** | **逃逸。见下** |
| **M6b 痕顺手改写返回历史** | 记录打完 `out[len(out)-1] = …Text:"trace-altered"` | `TestCompressFoldsOldestKeepsLastThreePreservesIDs`、`TestCompressTriggerScalesWithWindow`、`BooksCountsOnSuccess` | 抓到，但**不是靠 T4**（另两枚既有用例也在管） |
| **M7 只在挂了 logger 那侧改写正文** | `if c.lg != nil { out[last]=… }`（长度/part 数不变、字节变了） | `BooksCountsOnSuccess`（靠 `:219-220` 那两行拿 `c.TotalTokens(out)` 复量） | 抓到；**T4 不红** ⇒ T4 只比"条数＋每条 part 数"，**不比正文 bytes**，"不改折叠"这句在它里面是结构级的、不是字节级 |
| **M9 只在挂了 logger 那侧改报告** | `if c.lg != nil { rep.CompressedMsgs = 0 }` | **只有 `DoesNotAlterTheFold`** | 抓到，且**唯一证人**就是 T4 |

**M5 是本格的点名缺陷。** 派单 §4 问："若 `if true` 只让'没压缩也记账'那一枚红，那把'痕恒打'的写法挡得住吗？"
⇒ 读数：**`if true` 挡得住（M3 被 T2 咬住），而现实里更可能写出来的那一形——"只要这次 `Need()` 过阈值就打痕"——今天零枚用例抓得到。**
根因读得出来：T2 那枚反例种的是 `buildRoundHistory(2, 20)`，测试自己先断言
`if c.Need(hist) { t.Fatalf("setup: …want a no-op pass") }`（`compress_trace_test.go:270-271`）——
它把"没过阈值"当**唯一**的"没折叠"形状。于是**"过了阈值但折不动"这一形在四枚用例里不存在**，
而它恰好就是第 1 格里量到的**生产形状**（`raw_rounds=1 ≤ KeepRawRounds=3` ⇒ `Need()` 真、`rep.Ran` 假）。
⇒ 两格是同一个洞的两面：**今天真机上唯一会发生的压缩结局（"想过阈值、折不动"）既没人打痕、也没人测。**
（M5 造出的痕还会带 `history_changed=false`、`compressed_msgs=0`，读的人把它当成"压了但没压动"，
比"没有痕"更误导——这条是本程按 M5 语义推的读数级判断，未再单独造发验它。）

### 2.3 「摘掉任意一味，是否存在一发变异从此打不红」——逐味答

先给"每味各自被谁抓着"（`X+M` ＝ 同时摘掉用例 X 并施变异 M，整包 `-count=1`）：
```
[T4+M9]  reds: NONE     ← 摘掉 DoesNotAlterTheFold，M9 逃逸
[T1+M4]  reds: NONE     ← 摘掉 BooksCountsOnSuccess，隐私探针那发（泄正文）逃逸
[T2+M3]  reds: NONE     ← 摘掉 SilentWhenNothingFolded，恒真化（if true）逃逸
[T3+M2]  reds: NONE     ← 摘掉 SurvivesTheLoopWiring，loop 那行接线断了也没人抓
[M9-only] reds: DoesNotAlterTheFold   （T4 是唯一证人）
```
⇒ **四枚用例枚枚都有专属证人，摘掉任意一枚都真有一发变异从此打不红**：
- 摘 **`loop.go:239` 那行接线** ＝ M2 → 只有 T3 红；**T3 一并摘掉才逃逸** ⇒ 那行承重。
- 摘 **`compress.go` 那条 Info** ＝ M1b → T1、T3 都红 ⇒ 承重（实现件说"两发红"复算一致）。
- 摘 **T1** → M4（把正文塞进痕）逃逸；摘 **T2** → M3 逃逸；摘 **T3** → M2 逃逸；摘 **T4** → M9 逃逸。
⇒ 这一格实现件**站得住**，四枚用例没有一枚是装饰。
**但**上面 M5 那一发说明：**"承重"是相对于"他们想到的那四形"而言的**；
本程自造的一形（`if c.Need(hist)`）在四枚用例全在场的情况下整包全绿。⇒ 牙有，牙口够不够是另一件事。

### 2.4 本格档位

| 项 | 档位 |
|---|---|
| 实现件 §3 的 M1–M4 与逐名红句 | **成立**（本程独立复算，含"跑整包"这一更严口径） |
| AC#3「反向判据：摘掉痕该用例必须转红」 | **成立**（M1b 两发红） |
| AC#3 的"恒真判据＝本仓登记过的一类新假绿"这一支 | **附条件入账**：`if true` 形被 T2 挡住，**`if Need` 形今天无人挡**（M5 全绿） ⇒ 建议补的那一枚用例是"过阈值但折不动"（＝1 枚 user 轮、`Need()` 真、`Ran` 假）——**这只治 M5 与 §2.2 那一形，治不到第 1 格的"生产侧根本不出痕"**（那要落点选在调用方或改 `KeepRawRounds` 语义＝契约变更，本程不开那张面） |

**本程没测什么（本格）**
1. **没量 M5 造出的痕会不会被下游误读**（本格只证"逃逸"，未证"误导"到盘面/诊断的那一跳）。
2. **没造"痕重复打"那一发**（把 Info 挪进 `for` 折叠循环里 ⇒ 一轮多条）：实现件 §3.4#1 自己说这条判据"排在功能断言之后"，本程未独立验它响不响。
3. **没量 `t.Try/TestBench` 之类的名册外入口**：本程红名一律取自 `go test -count=1` 整包输出，未跑 `-count=2`（§5 那格统一跑）。
4. **变异只施加在 `internal/agent`**：没量"改 `cmd/wisp` 侧"（如把 `Logger` 传成别的）会不会红——今天那里压根不传。

---

## 3. 第 3 节：它那几把"零命中"的尺，逐枚验正控是否同一把尺

口径先说死：**"同一把尺"＝同一个 pattern、同一组 flag、同一次进程内跑**，只换 pathspec；
换成正则形状不同的另一种检（如拿字面量键名去给一条 `grep -i` 背书）本程**不算同尺正控**。
本程全部在 `git grep <sha>` 形态下跑（按 sha 取树，不读脏工作树）。

### 3.1 R1「`compress.go` 日志调用 0 枚」——**同尺正控复算成立，且实现件那个 0 是真 0**

```
$ git grep -cE 'Info\(|Warn\(|Debug\(|Error\(' cce7510 -- internal/agent/compress.go
   (无输出，rc=1)                        ← 改前（1acdd03 的父）＝票面 §0 那枚 0，复算成立
$ git grep -cE 'Info\(|Warn\(|Debug\(|Error\(' d949d9c5 -- internal/agent/compress.go   → 1
$ git grep -cE 'Info\(|Warn\(|Debug\(|Error\(' d949d9c5 -- internal/agent/loop.go      → 14   ← 同尺正控，逐枚对上实现件的 14
```

### 3.2 R2「`Compression` 生产侧零读者」——**三枚读数逐名对上**

```
$ git grep -cE '\.Rounds\b'       d949d9c5 -- '*.go' ':!*_test.go' | awk sum → 3     （实现件：3）
$ git grep -cE '\.Usage\b'        d949d9c5 -- '*.go' ':!*_test.go' | awk sum → 54    （实现件：54）
$ git grep -cE '\.Compression\b'  d949d9c5 -- '*.go' ':!*_test.go' | awk sum → 1     （唯一命中＝loop.go:401 赋值，等号左边）
```
顺把 §1.2 那三处"逐名"验了（不是只给计数）：
```
$ git grep -n 'agent\.Result' d949d9c5 -- 'cmd/wisp/**' ':!*_test.go'
   run.go:647 func notifyBody(res agent.Result)   run.go:658 func costLine(res agent.Result)
$ git grep -n 'loop\.Run(' d949d9c5 -- 'cmd/wisp/**' ':!*_test.go'   → run.go:576
$ git grep -ln 'agent\.Result' d949d9c5 -- '*.go' ':!*_test.go'      → 只有 cmd/wisp/run.go（全仓）
```
⇒ 实现件 §1.2 点名的"三处"复算成立，且 `agent.Result` 在非测试代码里**只出现在 run.go 一个文件**。

### 3.3 R3「`internal/panel`＋`internal/tools`＋`diagnostics.go` 对 compression 0 命中」——**0 站得住，但它那发正控不是同一把尺**

实现件 §1.5 写"所有'0'都先打了正控（[P1][P3][P5-control]）"，而 §1.2/§2.1 给这一枚 0 配的
是 `diagnostics.go` 的字面量键名 `"id" "ticket" "summary" …` —— **那是另一种形状检，不是 `grep -i compression` 这把尺**。
本程补上同尺正控（同 pattern、同 `-ni` flag，只换 pathspec）：
```
$ git grep -cni 'compression' d949d9c5 -- internal/agent
   budgets.go:2  compress.go:11  compress_test.go:7  compress_trace_test.go:17  loop.go:5      （合计 42，尺活着）
$ git grep -ni  'compression' d949d9c5 -- internal/panel internal/tools internal/observe/diagnostics.go
   (无输出，rc=1)                                                            ← 同尺同 flag 的 0，真 0
```
另把它承重的那支撑（"没人整枚序列化 `Result`，所以字段级零读者＝整结构级也看不见"）也验了：
```
$ git grep -n '%+v\|%#v' d949d9c5 -- cmd/wisp internal/panel internal/tools ':!*_test.go'   → 无输出
$ git grep -c '%+v\|%#v' d949d9c5 -- cmd/wisp internal/panel internal/tools（含 _test.go）    → 115
```
⇒ 尺在场 115 次、非测试侧 0 次，**"全落 `_test.go`"这句成立**。
**本枚档位：读数成立；实现件 §1.5 那句"所有 0 都先打了正控"过 claim（R3 那发的正控换了尺）。**

### 3.4 AC#1② 的缩放链与 `12000` 名册——**读数成立，但"本票贡献 0 枚"这句在交付版是假的**

```
$ git grep -n 'refHistoryTokens\|HistoryCompressTokens = s(' d949d9c5 -- internal/agent/budgets.go
   :36  refHistoryTokens  = 12000 // D15(4) compression trigger
   :107 b.HistoryCompressTokens = s(refHistoryTokens)            ← 等比缩放链复算成立
$ git grep -n '12000' d949d9c5 -- 'internal/agent/*_test.go' | wc -l → 9   （实现件：8）
$ git grep -c '12000' d949d9c5 -- 'internal/agent/*_test.go'
   compress_test.go:8   compress_trace_test.go:1                  ← 多出的那 1 枚就在本票新增的文件里
$ git grep -c '12000' cce7510 -- 'internal/agent/*_test.go' → compress_test.go:8   （改前确是 8）
```
读那枚新增命中的上下文（`compress_trace_test.go:29-30`）：它写在文件头注释里，
原文是 **"(AC#1(2) forbids citing the 12000 reference value as an actual trigger)"** ——
即它正是"引用禁令"本身。⇒ **AC#1② 的实质合规成立**（本程复核四枚用例的阈值一律从
`BudgetsFor(...)`/`c.b.HistoryCompressTokens` 读，未写死 12000；第 1 格那次真跑里量到的
`threshold=3072` 也来自 `BudgetsFor(32768)`，不是字面量）。
但实现件 §1.3 那句"**本票新增用例对 `12000` 的贡献为 0 枚**"在 `d949d9c5` 上不成立：
**本票确实新增了 1 枚 `12000` 文本出现**（在新增的测试文件里）。这句在它写 AC#1 那格的时刻为真、
在交付版为假——是自家 `1acdd03` 造成的漂移（同 §3.5 那枚 `D28-1` 计数）。

### 3.5 `D28-1` 名册与 §5 登记表（AC#1③）——登记缺格复算成立；代码标记枚数两家只有一个对

```
$ git show d949d9c5:docs/specs/SPEC-12-roadmap-governance.md | grep -c 'D28-1' → 0 (rc=1)
$ git show d949d9c5:docs/specs/SPEC-12-roadmap-governance.md | grep -n 'D28'   → 无输出（整份文件 0 枚 D28）
$ git show d949d9c5:docs/specs/SPEC-12-roadmap-governance.md | grep -n 'D11'
   91:| REJECTED | 意图分类四级流水线 | D11 过度设计已废 | …            ← 只有 1 行、且正是 REJECTED 那行
```
⇒ 派单 §4 那两句（`D28` 命中 0、`D11` 只有 1 枚且是 REJECTED）**逐字成立**，
实现件"§5 里根本没有 D28-1 那一格"亦成立（且比它写的更宽：整份 SPEC-12 都不提 D28）。
代码侧标记数：
```
$ for c in cce7510 1acdd03 e22da5a d949d9c5; do git grep -c 'DEFERRED(D28-1)' $c -- '*.go'|awk sum; done
   2 → 3 → 3 → 3
$ git grep -n 'DEFERRED(D28-1)' d949d9c5 -- '*.go'
   compress.go:27  compress.go:49  loop.go:394      ← 交付版是 3 枚
$ git grep -c 'DEFERRED(D11-3)' d949d9c5 -- '*.go' → control.go:1（1 枚，与派单一致）
```
⇒ **实现件 §1.4 的"两枚"在 AC#1 那格时刻为真，在交付版为假**；第 3 枚是它自己 `1acdd03` 写
`CompressionReport` 的文档注释时加进去的。**本程派单里那句"代码里却有 `DEFERRED(D28-1)` 三枚"经复核成立。**
（按派单 §4 明令：本票不动 `docs/specs/**`、不为这族写今天不响的检；票 150／`Q-55` 归编排者。）

### 3.6 本格档位与「本程没测什么」

| 那把尺 | 档位 |
|---|---|
| R1（日志 0）、R2（读者 1、逐名三处、`%+v` 全在测试） | **成立**，读数与本程独立复算逐枚相符 |
| R3（panel/tools/diagnostics 0） | **读数成立；正控配对不成立**（换尺），本程已补同尺正控 |
| AC#1②（等比缩放、不拿 12000 当触发点） | **成立**；但"本票对 12000 贡献 0 枚"这句**过期**（实为 1 枚、在禁令注释里） |
| AC#1③（那格不存在 ⇒ 未被覆盖） | **成立**（本程另核：整份 SPEC-12 零枚 D28） |

1. **没验 §5 表体"47 行"那句**（实现件用它背书"0 是扫过表得到的"）：本程改用同尺正控那条路径否证了
   "文件没读到"这一支，就没去复量那 47；若那句话错，影响面只有它自己的措辞。
2. **没量 `refHistoryTokens` 之外那 15 枚 `ref*` 常量是否也等比**（只核了本票点名的这一条链 :36→:107→scaleInt）。
3. **没量 `internal/observe/**` 里除 `diagnostics.go` 之外还有谁可能读 `Compression`**：票面点名的就是 diagnostics，
   本程按票面口径跑；全仓那枚已由 R2 的 `*.go` 名册覆盖（计数 1）。
