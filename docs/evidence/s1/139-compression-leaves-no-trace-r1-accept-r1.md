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
