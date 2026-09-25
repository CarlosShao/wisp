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
