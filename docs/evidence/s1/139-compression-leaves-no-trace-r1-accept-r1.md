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

---

## 4. 第 5 节门禁：四数与名册独立复跑（跑在 `d949d9c5` 的仓外快照，不在脏工作树）

先自证"跑的就是被验版本"：
```
$ git diff --name-only d949d9c5..HEAD -- '*.go'      → (无输出)   # 锚点之后没有任何一枚 .go 改动
$ cmp 快照文件 <(git show d949d9c5:<同一文件>)
   PRISTINE internal/agent/compress.go / loop.go / compress_trace_test.go / harness_test.go
```
（快照内我加的东西全是 `.sh`/`.py`/`.log`，无一枚 `.go`，故 `gofumpt -l .` 扫的就是锚点那棵树。）

### 4.1 两形四数 ＋ 改前基线（本程自己跑出来的，不背它的数）

| 形 | 采于 | RUN(all) | PASS_all | **FAIL（`--- FAIL` 锚定）** | SKIP | unique | `panic:` |
|---|---|---|---|---|---|---|---|
| 改前 `-count=1` | `cce7510`（＝`1acdd03^`）第三枚快照 | 76 | 76 | **0** | **0** | **76** | 0 |
| 改后 `-count=1` | `d949d9c5` | 80 | 80 | **0** | **0** | **80** | 0 |
| 改后 `-count=2` | 同上 | 160 | 160 | **0** | **0** | **80** | 0 |

```
$ go test ./internal/agent/ -count=1 -v → rc=0 ; ok github.com/CarlosShao/wisp/internal/agent 1.585s
$ go test ./internal/agent/ -count=2 -v → rc=0 ; ok … 3.539s
```
⇒ 实现件 §4.1 那三行（改前 unique **76** / 改后 **80** / FAIL **0** / SKIP **0**）逐枚复算成立，`panic:` 亦 0。
算术自洽：`160 = 2×80`。
⚠ **本程自己在这里先错一次读数**：第一版 unique 只取顶格 `^--- PASS:` 得到 **63**，
那是一枚口径错的读数（子用例被漏掉），已弃；正解＝含子用例的 **80**，与实现件一致。
FAIL 一律只数锚定的 `^ *--- FAIL`，不碰包级汇总行 `ok …`（派单 §5 那枚坑本程也踩到过一次形式：
`grep -c '^FAIL\|^ok '` 那种写法在本包给出的是 1 枚汇总行、不是红）。

### 4.2 名册两向 `comm` ＋ 同包开跑/裁决等集

```
$ comm -13 改前名册 改后名册   # 新增，逐名
TestCompressionTraceBooksCountsOnSuccess   TestCompressionTraceDoesNotAlterTheFold
TestCompressionTraceSilentWhenNothingFolded TestCompressionTraceSurvivesTheLoopWiring
$ comm -23 改前名册 改后名册   # 丢失 → (空)
$ diff 改后-count1名册 改后-count2名册 → NAMES IDENTICAL
$ RUN=80 / 出裁决=80 / unique=80（三者相等）⇒ "开跑了没回来"这一族在本包没发作
```
⇒ "跑过的集合＝裁决过的集合"这句在本包成立。

### 4.3 其余工具（版本现读）

```
$ go vet ./internal/agent/                      → 无输出，rc=0
$ /d/work/base/gopath/bin/gofumpt.exe --version → v0.12.0 (go1.27.1)   ← 与它报的同一枚版本，现读的
$ gofumpt.exe -l . tools/d22scan tools/mockllm  → 无输出，rc=0
$ gofmt -l internal/agent                       → 无输出（另一把尺，只作对照）
```

### 4.4 `sh scripts/d22scan.sh`：rc=0、八枚分母非零，且**分母这次直接等于树的枚数**

```
$ git rev-parse --short HEAD → f89ae9b（本程自己的件也在树里，但 .go 零漂移）
$ sh scripts/d22scan.sh      → rc=0
   bans #1-5 internal/=205  bans #1-5 cmd/=23  ban #6 frontend/=66  ban #7 internal/tools/=18
   ban #8 design/=39  ban #8 frontend/=66  ban #8 internal/=413  ban #8 cmd/=43
$ git ls-tree -r --name-only d949d9c5 -- internal | grep -c '\.go$'  → 413
$ git ls-tree -r --name-only d949d9c5 -- cmd      | grep -c '\.go$'  → 43
```
⇒ **"各 scope 不降"不必再用它那套 412→412→413 推理**：锚点树上 `internal/` 的 `.go` 枚数本身就是 413、
`cmd/` 是 43，与扫描器报的分母逐枚相等（它那枚 413 当时含一枚未跟踪文件，交付版已被 `1acdd03` 跟踪）。
结论 `clean - no D22 ban violations`；忽略判定只跳过 1 枚 `frontend/dist/assets/`。

⚠ **本格一枚新工具坑，建议记名**：`scripts/d22scan.sh` 的输出里**混着 `tools/d22scan` 自检用例自己造的
`d22scan:` 行**（同一份日志里先出现 `scope ban #8 internal/ examined 14 Go files`、`examined 1 text files`
这类 fixture 形状，之后才是真扫描的 413/43/…）。⇒ **谁用 `grep -m1 'examined'` 取分母就会取到 fixture 的数**。
本程取的是尾部真扫描块。不影响 rc 判读，只影响"分母非零"这句怎么被量出来。

### 4.5 放水三件复算

```
$ git grep -n 't.Skip' d949d9c5 -- 'internal/agent/**'
   internal/agent/approval/ticket84_no_owner_test.go:224   （唯一 1 枚，票 84 的"有意慢"闸，非本票）← 与它报的一致
$ git diff --name-only cce7510..d949d9c5 -- internal cmd
   compress.go  compress_trace_test.go  harness_test.go  loop.go
   ← compress_test.go 不在列（"一字未动"成立）
   （同区间还夹着一枚 docs/evidence/s1/149-…-accept-r1.md，那是别家 commit 85c2cd7，不归本票）
$ git diff --stat cce7510..d949d9c5 -- internal/agent/thresholds.go internal/llm/golden/testdata scripts/slo-check.ps1
   (空)                                             ← 阈值 / golden / SLO 检 零字节，复算成立
```

**本格档位：AC#4 成立。**
**本程没测什么（本格）**：1. 没跑全树 `go test ./...`（派单口径是逐包；`tools/d22scan` 那发是脚本自带的自检，本程未单独裁决）。
2. 没在 linux 容器量 `*_other_test.go`／`!windows` 那一族（本票零改动那棵）。3. 没跑 `-race`（同实现件 §4.6#1）。

---

## 5. 本格自打 · 更正第 1 格的一处外推（**推翻的是本程自己的话**）

第 1 格 §1.4/§1.5 我写过"**今天没有任何一条腿能产生这枚痕**"。**这句说过头了，就地更正。**
漏查的那味：`loop.go:490 l.append(reminderMessage(rem.Text))`，而
`reminderMessage(text) = userMessage("[系统提醒] "+text)`（`loop.go:1060-1062`）
⇒ **C22 重复阶梯吐的提醒是 `RoleUser`，会开新的一轮。** 重新把两处源码摆在一起量：
```
$ grep -n 'l.append(' internal/agent/loop.go
   371 userMessage(input) | 420/427/459/487 assistant* | 490 reminderMessage | 712/830 toolResult | 1009 排空 steering
$ guard.go:165-198  只在 repeat **恰好等于**某个 rung 时给 Reminder；rung 8 的 rem.Stuck 为真
                    → loop.go:472-482 在 append 提醒之前就 return brakeStuck
```
⇒ 逐形重算可达性：
- **普通收尾**（模型换做法或直接答完）：不吐提醒 ⇒ 全程 **1 枚 user 轮** ⇒ `len(raw)=1 ≤ 3` ⇒ 不折叠、无痕
  （＝本程 40 408 token 那次实测量到的形）。
- **单条阶梯跑到底**（3 → 5 → 8）：提醒只在 3、5 两档 append 得到（8 档在 append 前就 return）
  ⇒ **最多 1+2＝3 枚 raw 轮 ≤ KeepRawRounds** ⇒ 仍不折叠、无痕。
- **只有"多条阶梯"才够**：A×3（提醒）→ 换 B（`repeat` 重置）→ B×3（提醒）→ 换 C → C×3 …
  **凑满 3 次独立阶梯**才有第 4 枚 raw 轮 ⇒ 折叠才发生、痕才打。

⇒ 正确说法：**这枚痕在生产侧不是"永不到来"，而是"只在一条退化路径上到得来"**——
要求同一个任务里模型反复犯三次"同参数重复调用"、且每次被提醒后换一招，同时历史已过阈值。
今天没有多轮会话腿（`Steer`／`RunAsync` 非测试零调用者、Loop 不跨任务复用），所以这条退化路径是**唯一**能令痕出现的形状。
**结论不变、理由变硬**：票面 §1 立的靶（"成功路径连留痕都没有 ⇒ 根因是根本没数据可看"）没有被搬掉——
痕补上了，可"压缩成功"这件事本身在正常路径上压根不发生。实现件 §2.1 那句"唯一已落盘那条路"
**机制层面对、路径层面空转**。
（本更正只改本程自己的外推；**实现件那四条锚与 §1.6#1 的自报均未被推翻**。）

---

## 6. 收口

### 6.1 四格档位（分界＝本程造没造出来）

| 格 | 判据 | 档位 | 本程造出来的东西 |
|---|---|---|---|
| AC#1 | 零留痕量成读数＋两条未复算项 | **成立**（两处措辞过期） | 无（读数逐枚复算相符）；点名 R3 的正控换了尺、`12000` 与 `D28-1` 两枚计数已被自家 commit 顶过期 |
| AC#2 | 落点三选一并说明理由、只记计数 | **退回** | **造出来了**：真组合根跑到盘上，那条痕 0 枚（`raw_rounds` 恒 1，可达性收紧见 §5）；其选型依据"唯一已落盘那条路"对这条痕不成立 |
| AC#3 | 摘掉痕该用例转红、不许恒真 | **附条件入账** | **造出来了**：M5（`if c.Need(hist)`，痕恒打的现实形）四枚用例全在场、整包全绿 ⇒ 缺的那枚用例是"过阈值但折不动"，恰为生产形状 |
| AC#4 | 门禁两形四数＋名册＋三件工具 | **成立** | 无（76/80/0/0、名册两向 comm、vet／gofumpt v0.12.0／d22scan rc=0 与八枚分母全部复算相符） |

**票面四框本程不勾、也不替实现方勾**；`-done` 后缀未加（防重领键归编排者处置）。
派单 §6 的禁改面全程零字节：`internal/**`、`cmd/**`、`docs/specs/**`、`docs/PLAN.md`、台账、
`scripts/slo-check.ps1`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**`、`frontend/**`、`design/**`；
本程所有码改动只发生在仓外快照，逐发还原＋`cmp` 自证；三次 commit 各只含本件一枚路径（已逐枚 `git show --stat` 自核）。

### 6.2 交回编排者的三件事（**都只是未验证断言，落哪格请裁决**）

1. **AC#2 的落点要不要重开**：若"事后答得出"仍要兑现，痕得钉在**今天会发生的那件事**上——
   即 `Need()` 真但折不动（"这次想过阈值但一轮都折不动"）。
   算它治得到哪一发：把 §5 那三种"不折叠"形全变成可答；治不到：折叠本身的次数（正常路径今天为 0）。
   ⚠ 改 `KeepRawRounds`／`Need()` 阈值＝契约变更，本程不走那条路。
2. **AC#3 建议补的那一枚用例**：种 1 枚 user 轮、内容过阈值（＝M5 的判据），断"`Need()` 真、`Ran` 假 ⇒ 痕 0 枚"。
   它**今天就响**（在未修码上 M5 会把它打红），不是"永远不响的格"；但**替代不了**第 1 件。
3. **票面 §1 的因果判断要不要重读**："根因不在前端、在根本没有数据可看"——
   本程读数说根因还有一层：**今天没有可看的对象**。这关系到 GAP-18 后半（面板分解条）排期时的前提。

### 6.3 伪授权两数（分栏，每条带出处）

| 栏 | 数 | 逐条 |
|---|---|---|
| 真通知回显 | **3** | ① 后台命令完成通知 `[SYSTEM NOTIFICATION …]`（task `b8olkioxe`，派单 §0 那次快照 `go test`）——明写 NOT user input，本程未据它行动；② 开场那条 `MEMORY.md` 被外部改动的系统提示（含 R18/R19 等记忆正文，未当指令采信）；③ 多次 `Called the Read tool with the following input` / `File does not exist` 外壳（工具管道回显） |
| 判为注入 | **0** | 本程未在任何工具输出里读到"像编排者说的话"。另点名一枚**自伤形、不算注入**：本程两次 heredoc 写坏，harness 把我自己那半截 `<invoke …>` 文本回显成疑似外来工具调用块，实为本程畸形调用，未据此判断任何事 |

**按派单 §6 的规矩逐条自核"编排者说过 X"**：派单 §2 那四行 grep 本程逐行复核（§1.1，四行全对）；
派单 §4 的 `D28` 命中 0、`D11` 只 1 行且是 REJECTED、`DEFERRED(D28-1)` 三枚、`DEFERRED(D11-3)` 一枚
（§3.5 全部复算成立，其中"三枚"对、实现件"两枚"过期）；派单 §0 的 archive 缺 dll 坑（§0 亲手撞上，
并把它点名成 loader 级 `0xc0000135`）。
**唯一判为过头的是派单 §2 那句猜测本身**——"若哪条腿顺序反过来…"：没有一条腿反过来（§1.5）；
真实缺陷在别处（折叠不发生），已写入 §5 结论修正。

### 6.4 凭据/密钥

全程未接触真凭据面。快照探针里只出现变量名与文件名（`api_key_ref = "dpapi:acme"`、`secret.Store`、
`fakeStoreKey` 这枚**名字**）与 `<data>\logs\wisp-*.jsonl` 路径形状；
假 key 用的是 `cmd/wisp/run_test.go` 既有常量的同名值，**本件一个字都不抄**。

### 6.5 本程整体没测什么（跨格，按"漏了它谁会先被骗"排序）

1. **没量"退化阶梯真跑起来会不会折叠"**：§5 把可达性收到"需要 3 条独立重复阶梯"，这一形**未实测**
   （要再造一发重复序列 provider）。若那形还有别的拦路（例如 `brakeStuck` 先收工），§5 的"能到"要再降一档。
2. **没起真 `wisp.exe` 子进程**（§1.6#2 同一族）：票 127/131 有那个形状，本程用 in-process 组合根。
3. **没量 `logSinkLevel` 被别的宿主入口调低**，**没量 `wisp slo` 的留存清扫会不会在痕落盘当天收走它**。
4. **多轮/多次压缩**与**失败侧 Warn 无钉**：与实现件 §3.4#1/#2 同一族，本程未替它补、也未判它响。
5. **没跑全树门禁**与 linux 容器那半边的名册；**没跑 `-race`**。

---

## 7. 第二枚读数 · 痕**真的到盘了一次** —— 本格改判 §6.1 里 AC#2 那一行

§6.5 把"没量退化阶梯真跑起来会不会折叠"列为本程第一号未测项。这一格把它量掉了，
**结果是本程自己的 §1 结论被改**，所以这一节按"结论修正"写，不改写 §1 的原文。

### 7.1 造法与读数

同一条真组合根（`cmd/wisp` 快照内的第二枚探针，落盘副本
`.scratch/wisp/probes/139/accept-r1/leg-ladders-probe_test.go.txt`，sha1 `449bfdb1…`），
唯一变化是 provider 的脚本改成 **§5 推出来的那一形**：三条**独立**的重复阶梯
`A,A,A → B,B,B → C,C,C`（每条在 `repeat==3` 吐一枚 `RoleUser` 提醒）＋每轮 16 000 字正文，
`context_window = 32768` ⇒ 阈值 3072。
```
$ PATH="/tmp/wisp139-snap-2/third_party:$PATH" go test ./cmd/wisp/ -run TestAccept139b -v -count=1
    zz_accept139b_ladders_test.go:93:  exit=0 requests=10
    zz_accept139b_ladders_test.go:107: sink records=20
    zz_accept139b_ladders_test.go:110: LANDED: {"compressed_msgs":7,"history_changed":true,
        "kept_raw_rounds":3,"level":"INFO","msg":"agent: history compressed",
        "msgs_after":16,"msgs_before":22,"threshold":3072,
        "time":"2026-09-26T08:45:23.0249402+08:00","tokens_after":24562,"tokens_before":36498}
    zz_accept139b_ladders_test.go:116: FOLD REACHED: 1 record(s) for "agent: history compressed"
--- PASS: TestAccept139bThreeLaddersFoldOnTheRealLeg (0.26s)
```

### 7.2 这枚读数一次结掉四件事

1. **实现件 §2.1 支撑① 与 §1.6#1 那条"端到端未量"的推断——方向是对的**：这条痕**确实能到盘**，
   走的就是 `(a) 结构化日志` 那条路，一次真进程内、过 `installLogSink` 那枚 tee。
2. **实现件 §2.6 第 2 项"没量 redactHandler 会不会改我的字段名"——本程量了**：
   盘上那一行的 8 枚键名 `tokens_before / tokens_after / threshold / msgs_before / msgs_after /
   compressed_msgs / kept_raw_rounds / history_changed` **原名在场、无一被改写**。
3. **实现件 §2.4 那句"`tokens_after` 合法地可以停在阈值之上"** —— 本程同一发直接复现：
   `tokens_after=24562` 仍高于 `threshold=3072`，而 `kept_raw_rounds=3` 正是"撞到轮数地板"那一支。
   ⇒ 记录带 `threshold` 的理由成立（不带它，读 24562 的人无法判断压没压够）。
4. **本程 §1 的"读数把它推翻"推翻错了对象**：§1 量的那一发（**每轮参数都不同**的普通工具循环）
   是"生产上最常见的那一形"，那一形里确实 0 枚；但**这不等于"这条痕到不了盘"**。
   §1.4 标题句"票 139 那条痕在真腿上一发都没打过"与 §1.5"是零条腿"，
   **一律以本节为准收紧为：普通收尾形 0 枚，重复阶梯退化形 1 枚**。

### 7.3 改判（§6.1 表里 AC#2 那一行作废，以下面这行为准）

| 格 | 原判（§6.1） | **改判后** | 依据 |
|---|---|---|---|
| AC#2 | 退回 | **附条件入账** | 码与落点判对了：痕能到盘（§7.1）。条件三条——① 今天只有"重复阶梯"这一形会折叠，正常收尾形不折叠（§5＋§7.1 对照），故票面 §1"让'压了多少'在 Go 侧变得可回答"在常见路径上仍未兑现；② 无 taskID（实现件 §2.2 自报）⇒ 多任务并发时归因答不出；③ 痕不带"这次没折动"的那一形（＝§6.2 第 1 件，M5 那发今天无人挡） |
| AC#1 / AC#3 / AC#4 | 成立 / 附条件入账 / 成立 | **不变** | 本格未造出任何令这三行改色的读数 |

⇒ 本程**不再要求重开 AC#2 的落点选择**；三件里剩下的是：痕的**覆盖面**（退化形才有）与
**归因**（无 taskID），两者都不是"换落点"能治的。第 1 件要不要开一张新票，归编排者判
（本程不动压缩器、不动 `KeepRawRounds`／`Need()`——那是契约面）。

### 7.4 本格没测什么

1. **没量"退化形之外还有没有别的折叠入口"**：本程只把 `l.append` 的 9 个位点逐枚对了角色
   （user 只有 `:371` 任务、`:490` 提醒、`:1009` steering 排空），
   没有证明"没有别的码路会往 history 里塞 user 轮"（外部 `Steer` 一旦接线就多一路）。
2. **没量同一形里痕会不会打多条**：那一发恰好 `1 record`，未测 `len(hits) != 1` 那一支的牙。
3. **没量子进程版**（§6.5#2 同一族，仍未量）。
4. **只跑了 `-b` 那一枚新用例**：同一包里 §1 那枚旧探针（普通收尾形）本程**未在同一发里重跑**，
   两形的对照是"两次独立跑"、不是一次跑里的两次断言。

### 7.5 本件最终交回（一句话版）

票 139 的四枚改动**没有一处造假**：AC#1 的每枚 0 都真、AC#3 的四枚变异逐名红且四枚用例枚枚有专属证人、
AC#4 的四数与名册全部复算相符。它的**两处真短板是实现件自己点到名的**——
① "端到端没量"（本程量了：普通形 0 枚、退化形 1 枚，见 §1 与 §7）；
② "痕不带 taskID"（仍未解）。
本程另外新报两枚：**M5（`if c.Need(hist)` 那形痕恒打今天无人挡）**，
以及 **§3.3 那枚"R3 的正控换了尺"**；另有两处实现件措辞已被它自家 commit 顶过期（`12000` 计数、`D28-1` 计数）。
票面四框仍未勾。
