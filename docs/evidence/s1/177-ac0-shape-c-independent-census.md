# 177 AC#0 交件表 — 候选丙（「每次续读都要一张人批的 L2 卡」）的独立普查

- 写手程：`177-a0`（派单只给 **AC#0**；票面明写「⚠ 不许并入甲/乙那一趟做」，本程独立起手）
- 起手时刻：`2026-11-05 20:44:07 +08`（`date` 现量）
- 分支/锚：`dev`、HEAD `34637ae6c`。派单写 `7eb64d5ce`；现量它是 HEAD 的祖先，**漂移 3 枚**（`b914db12e`／`6fef2a948`／`34637ae6c`，全是票 51/246 的收尾），且 `git diff --stat 7eb64d5ce..HEAD -- internal/ cmd/ tools/` **为空**、`git log 7eb64d5ce..HEAD -- internal/ cmd/ tools/` **零提交** ⇒ 本程读的产码面与派单锚一致。
- scoped 树（起手）：`git status --porcelain -- internal/ cmd/ tools/` 只有 `?? internal/tools/ticket23_image_channel_test.go`（**别家的未跟踪件**，不是本程写的；本程没有动它）。
- 结论一句话：**候选丙在真机上有四种不同形状，不是一句「会拒」**——没有审批层＝**永不咨询、当场拒**；`NoGate`＝**当场拒**；界面不可达＝**当场拒**；界面在、没人答＝**挂起到 C18 超时再拒**（测试用可配的 31s 量到墙钟 34.26s，生产默认 300s）。R4 的「不受 D45 覆盖」在代码里是**结构性的**：只有 L1 才会去查会话授权，L2 连查都不查 ⇒ 丙形下 **N 次续读＝N 张卡、N 段等待**，且票 164 AC#3 的两枚判据必须改形（＝AC 变更，须 owner）。
- 本程**零产码**：名下只新增本文件；没有改 `internal/**`／`cmd/**`／`tools/**` 一个字节。

## 0. 本格要答的三问（票面逐字）

> ① 真机上 `agent.Loop` 在**没有接审批通道**的界面里遇到 L2 是什么行为（CLI 那条路上命中 L2 之后到底是「拒绝」、「挂起」还是「永不咨询」）② R4 那一档「不受会话授权覆盖」（`SPEC-06:37` 逐字）在**每一次**续读上意味着什么 ③ 若接受丙，票 164 已勾的 AC#3 那两枚续读判据要按什么新形状改写（＝AC 变更，须 owner）。

## 1. ① 真机上 L2 的四种形状（四发都是〔现跑〕；命令逐字可复跑）

尺（本程所有跑法都带同一个环境前缀）：

```bash
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
```

| 形状 | 现跑命令 | 读数 | 定性 |
|---|---|---|---|
| **A. 没有审批层**（`agent.Options.AdmitTask == nil`），工具声明 L1/L2 | `go test -count=1 -run '^TestLoopStillRefusesL1WhenNoGateIsRegistered$' -v ./internal/tools/` | `--- PASS (0.17s)`；工具行 `Rejected=true`、`error_class=permission_denied`、正文含「审批通道尚未接入」、`decision=reject` 落账、目标文件不存在。⚠ 该用例喂的是**声明 L1** 的一发，但它走的就是同一条 `case RiskL1, RiskL2:`（`loop.go:919`），L2 没有第二条分支 | **永不咨询 ＋ 拒绝**（一张卡都不出，模型直接拿到错误） |
| **B. 桥 + `NoGate`**，工具声明 L2 | `go test -count=1 -run '^TestRiskDecisionRoutesEachLevelToItsBranch$/^NoGate_fails_closed_above_L0$' -v ./internal/tools/` | `--- PASS (0.00s)`；`Runs()==0`（工具一次都没执行）、正文含「票 21」 | **立即拒绝**（`NoGate.PendingApproval` 那一句） |
| **C. CLI 真装配（console 腿）＋ L2 ＋ 没有答复面** | `go test -count=1 -run '^TestUnansweredL2CardTimesOutIntoRejectNeverExecution$' -v ./cmd/wisp/` | `--- PASS (34.26s)`；卡片真的打印、C18 预警「审批将在 30 秒后自动拒绝」出现在 stdout、终止是 `ANSWER-EXPIRED … decision=timeout->reject`、`tool_call.decision=timeout`、目标文件**没有落盘** | **挂起**到 C18 截止（本发配置 31s）**再拒绝**；绝不在超时后执行 |
| **D. 有 gate、界面不可达** | `go test -count=1 -run '^TestAC246ResidentPipelineAsksThroughTheOneGate$' -v ./cmd/wisp/` | `--- PASS (2.84s)`；正文＝`审批界面不可达`（`gate.go:658-660` 的 fail-closed 句）、`displayedCards()==0`、文件未被改写 | **立即拒绝**（连挂起都不挂） |

每格对应的产码原文〔读码〕：

- A：`internal/agent/loop.go:916-924` —— `case RiskL1, RiskL2:` 先看 `l.opt.AdmitTask == nil`，是则 `j.decide(..., DecisionReject)` 并返回 `"该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"`。**整段没有调用任何 gate** ⇒ 「永不咨询」是代码结构，不是运气。
- B：`internal/tools/gate.go:147-161` —— `NoGate.PendingApproval` 固定返回 `AnswerReject, "L2 审批通道尚未接入（票 21），已拒绝执行"`。
- C：`internal/agent/approval/gate.go:617-681` —— `PendingApproval` 先 `q.push`／`grantNonce`／`ui.Prompt`，然后在 `it.answer`／`warn`／`deadline`／`ctx.Done()` 上选；匿名 0 秒不会自动放行，只有 `deadline` 那支 `q.expire`（`AnswerTimeout` → 调用方按拒绝处理）。`internal/agent/approval/queue.go:124-125`：`DefaultApprovalTimeout = 300 * time.Second`（「超时 300s 一律判拒绝」）。超时值由 `risk.confirm_timeout_sec` 走，`internal/config/validate.go` 把它夹在 `[31, 3600]`，所以本程的 31s 是**真的生产同一条路**、不是假时钟。
- D：`internal/agent/approval/gate.go:656-661` —— `ui.Prompt` 返回错误时 `q.abandon` 并直接 `AnswerReject`（`审批界面不可达，已 fail-closed 拒绝执行`）。

⚠ 一处口径更正（给下一位）：`docs/evidence/s1/179-cli-reread-r2.md` §4 那行 `窗口 0.0s` 是 **L2 卡片的显示值**，不是等待时长——现读 `cmd/wisp/run.go:2271-2275` 已经专门为 L2 改印 `p.Deadline`（注释逐字：「an L2 card's Window is 0 by construction; its wait is Deadline」）。所以那枚读数**不能**当作「L2 0 秒即拒」的证据；L2 的真实等待是 C18 截止（默认 300s）。本程没有拿它当读数用。

## 2. ② 「R4 不受 D45 会话授权覆盖」在每一次续读上是什么

- 规格逐字〔现读〕：`docs/specs/SPEC-06-security-gatekeeping.md:37` —— `| R4 | **C25 污染命中**（敏感读 × 外泄通道） | → L2，且**不受 D45 会话授权覆盖** |`。
- 产码是**结构性**排除，不是「查了再拒」〔读码〕：`internal/tools/bridge.go:383-403` 的会话授权查表整段包在 `if !sil.Silenced && sil.Level == risk.L1` 里 —— **L2 根本不进入查表分支**（注释逐字引 `SPEC-06 §8.3` 第一颗子弹与 D45-3：「an L2 verdict is never covered by a session grant」）。
- 二次确认〔读码〕：`internal/risk/mode.go:222-223` —— `if d.SessionOverrideBlocked { return "R4 污染升级（C25，任何授权都不可覆盖，含权限模式）", true }`；`internal/risk/assessor.go:96-98` 把 `SessionOverrideBlocked` 定为「must not be covered by a D45 scope session authorization (R4 taint hits, SPEC-06 §8.3)」。
- 也不许被「批量聚合」稀释〔读码〕：`internal/agent/approval/batch.go:27-40` —— `L2 is NEVER aggregated`，且 `d.SessionOverrideBlocked || hasRule(RulesHit, risk.R4)` 时 `Aggregate` 直接返回 nil。
- 规格/决策层的两处锚〔现读〕：`docs/PLAN.md:1598`（「**且 L2 永不聚合、永不进持久授权**」）、`docs/PLAN.md:1647-1648`（「批量聚合确认**不得**把 L2 聚进去；会话授权**不得**覆盖 C25 污染升级」）；`docs/PLAN.md:2034` 定案⑤「**C25 污染命中（R4）不受本授权覆盖**」。
- ⇒ **对「每一次续读」的含义**：丙形下模型每 `fs.read` 一次宿主指针，就是**一张新的 L2 卡**；卡上「本会话内允许」那一档对 R4 **永远不会解锁**（`bridge.go:395` 连查都不查），所以第 2 次、第 3 次……第 N 次都要人再答一遍。CLI 无人值守时每一次都是 §1 的 C 形（挂起到 300s 后拒），不是「首次问一次、以后放行」。
- 这条与 `177-c1` §4.1 的「名册形」无关：本程只回答「已有 L2 的语义」，没有再提任何豁免候选；也**没有引用** `177-c1` 那张表当证据（票面 09-27 23:0x 第 4 条已把它降为「线索」）。

## 3. ③ 若接受丙，票 164 已勾的 AC#3 两枚判据要改成什么形状（＝AC 变更，须 owner）

现状〔现读〕：两枚判据都活、都绿〔现跑〕：

```bash
go test -count=1 -run '^(TestLongOutputPointerRecoversEveryByte|TestPointerPast256KiBIsNotFullyReadable)$' -v ./internal/tools/
# --- PASS: TestLongOutputPointerRecoversEveryByte (0.02s)
# --- PASS: TestPointerPast256KiBIsNotFullyReadable (0.16s)
```

它们今天绿的前提是**甲形**（`177-r1` 的宿主自产路径豁免）＋ **`NoGate`**：两枚都用同一个 `taskBridge`，`Gate: NoGate{}` 写死在 `internal/tools/task_output_leg_test.go:48`。

丙形下这两枚**会当场红**，逐字形状不是推理——`docs/evidence/s1/175-mark-gate-r1.md:147-155` 记录过一次真实撞红（那是「盖戳已落、豁免未落」的树）：

```
--- FAIL: TestLongOutputPointerRecoversEveryByte
    task_output_leg_test.go:220: the pointer must be re-readable with fs.read, got: {Text:L2 审批通道尚未接入（票 21），已拒绝执行 IsError:true RiskLevel:L2 …}
--- FAIL: TestPointerPast256KiBIsNotFullyReadable
    task_output_leg_test.go:381: the first segment must read back, got {Text:L2 审批通道尚未接入（票 21），已拒绝执行 …}
```

⇒ 丙形要保住 AC#3，判据必须按下面**新形状**改写（三件缺一即装饰）：

1. **换成会答的 gate**：两枚判据的 bridge 不能再是 `NoGate`，要注入一个记录型 `Gate`（`PendingApproval` 返回 `AnswerAllow` 并记下 `Decision`）。
2. **把「每次一张卡」写成断言**：每次 `fs.read` 前后 `cards++` 必须 +1；卡上 `RulesHit` 含 `risk.R4`、`SessionOverrideBlocked == true`、`Reason` 指明来源是 `task.output`（而不是宿主自己写下的那枚来源）；**第二次**读同一路径再 +1（证明「本会话内允许」盖不住它）。
3. **把「读回全部字节」改成「读回 C1 能读回的那一段」**：`fs.read` 没有 offset、`max_bytes` 只能往小改不能往大改（`internal/tools/fs.go:55` 的 `defaultMaxReadBytes = 256 * 1024`、`:151-155` 的 `if a.MaxBytes > 0 && a.MaxBytes < limit` 夹法、`:115-121` 的 schema 只有 `path` + `max_bytes`）⇒ 丙形下 >256 KiB 的产物**再加多少张卡也读不到尾巴**（那是 `Q-59` 的地界）。AC#3 今天那句「recovers every byte」在丙形下只对 ≤256 KiB 成立；`TestPointerPast256KiBIsNotFullyReadable` 会被迫从「canary 写着今天读不全」变成「丙形下每段一张卡、尾巴仍读不到」。

⚠ 第 3 条是**改已勾的 AC 判据**，票面明写「那是 AC 变更＝owner 批准在前」；本程只把它摆清楚，**没有改任何一枚判据**。

## 4. 丙的代价（用本程/上游现量数字，不用形容词）

| 量 | 读数 | 出处 |
|---|---|---|
| 单次续读的等待 | 生产默认 **300s**；本程用可配下限 31s 量到墙钟 **34.26s**（含装配） | `queue.go:124-125`；§1 C 的现跑 |
| 一次 20 000 字节产物的「读到全文」 | 丙形下 **1 张卡**（`TestLongOutputPointerRecoversEveryByte` 只发一次 `fs.read`） | `task_output_leg_test.go:212-224` |
| 一次 300 KiB 产物的「读到全文」 | 丙形下至少 **2 张卡**（第一次读满 256 KiB，再试 `max_bytes` 一次），且**尾巴仍读不到** | `task_output_leg_test.go:346-420`；`fs.go:55/115/151` |
| 会话授权能否把 N 张卡压成 1 张 | **不能**：L2 不查会话授权（结构），L2 也不聚合 | `bridge.go:395`；`batch.go:39` |
| 无人值守 CLI 上读完一个 N 卡流程 | **N × 300s 的挂起**，每段都以拒结束 ⇒ 实际等于读不回来 | §1 C ＋ §2 |

## 5. 本程没测什么（别把上面的绿读成「丙形已量完」）

- **没有真跑 300s**：本程用的是 `risk.confirm_timeout_sec` 的可配下限 31s（生产同一条路，不是假时钟）。300s 那一段是〔读码〕不是〔现跑〕。
- **没有起独立 `wisp.exe` 进程**：C/D 两发是 `cmd/wisp` 的真装配（真 `run()`、真 `approval.Gate`、真 `consoleApprovalUI`、真 golden provider），不是外壳进程。
- **没有量丙形本身**（把豁免摘掉后的红集合）：那要临时改 `internal/risk/**`，是 AC#1b 的变异纪律；AC#0 只做普查，不重开变异。丙形的红面用 `175-mark-gate-r1.md:147-155` 的那次真实撞红作证据。
- **没有量面板/原生卡那条路**（丙形若走 resident 腿，L2 卡落在 `ballCardUI`）：本程只量了 console 腿与 resident「无球」fail-closed 那一支。
- 没有动票面任何 AC 框、没有动 AC#5（邻居分账）与 AC#6（排程）。

## 6. 门禁与自证（本程零产码，门禁按「本程引用到的包」取）

| 门 | 现跑读数 |
|---|---|
| `go build ./cmd/... ./internal/... ./tools/...` | `rc=0`，无输出 |
| `go test -count=1 -timeout 420s ./internal/tools/ ./internal/risk/ ./internal/agent/ ./internal/agent/approval/` | **四包全绿**（`21:05:31` 起跑）：`ok internal/tools 26.367s`／`ok internal/risk 14.516s`／`ok internal/agent 9.726s`／`ok internal/agent/approval 0.887s`，`RC=0` |
| `go test -count=1` 两发定向 `./cmd/wisp/`（§1 C/D） | 两发 **PASS**（34.26s／2.84s） |
| `gofumpt -l internal/tools internal/risk internal/agent cmd/wisp` | **空**（`v0.12.0 (go1.27.1)`） |
| `sh scripts/d22scan.sh` | `rc=0`；`clean - no D22 ban violations`（scopes 逐条打印，见日志） |
| 产码零改动 | `git status --porcelain -- internal/ cmd/ tools/` 只剩别家的 `?? internal/tools/ticket23_image_channel_test.go`；本程 commit 只带本文件 |

### 6.1 ⚠ `./cmd/wisp/` 整包在这台机上**不是可用的门**（具名红，不是本程造成的）

本程先按门禁把 `./cmd/wisp/` 整包也放了进去，连跑三发都不是绿——**三发的形状各不相同，必须分开记**：

| 发 | 命令 | 读数 | 归因 |
|---|---|---|---|
| 1 | `go test -count=1 -timeout 420s ./cmd/wisp/` | `FAIL 420.250s`（超时，tail 停在 `lifecycle_43_windows_test.go` 的常驻 owner goroutine）。当时机器上**另有一条腿**在跑 `go.exe test -v -count=1 ./...`（`Win32_Process` 现量：PID 57568、`CreationDate=2026/11/05 21:02:42`、命令行逐字 `D:/work/base/go/bin/go.exe test -v -count=1 ./...`） | 并发争用 |
| 2 | 同一条，另一条 `./...` 收工后复跑 | `FAIL 420.285s` | 仍超预算（本包历史读数 110–220s；这台机当时被多条腿占着） |
| 3 | `go test -count=1 -timeout 900s ./cmd/wisp/`（**不是超时**，跑完才判） | `FAIL 635.956s`；20+ 枚顶层红，逐类见下 | 大部分是环境；**至少一枚像是 HEAD 上的真漂移** |

第 3 发里逐类抽读到的原文（每条都在 `"$HOME/177-cmdwisp-full900.txt"` 里可复跑）：

- **语音链（占大头：`Test15e*`／`Test26P*`／`Test41a1*`／`Test395r1*` 等）**：`语音链未构造：VAD 没起来（speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)）` ⇒ **系统目录里的旧 onnxruntime.dll 抢先加载**，不是代码。
- **`TestAC4FocusReturnToPriorWindowGap33r5`**：`the panel did not take the foreground on Show: foreground 0xbec0936, panel hwnd 0xa450fc6` ⇒ 需要**交互桌面**的前台焦点，在本机（多腿共享、窗口被占）不成立。
- **`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`**：`the resident leg passes an unexpected Options field "Answer": the shipped field set is no longer the one ticket 256 adjudicated`／`passes 7 of Options' fields, want 6` ⇒ **这条不是环境的味道**：它是「常驻腿喂给 `approval.Options` 的字段集变了」。本程**不认领、不修**（写面之外），但必须具名交给它的归口（票 256／375 的接缝）。
- **`TestAC247HandingTheLevelToTheSeamIsNotVisibility`**：`prototypeVisuals is ON in this assembly: the shipped default is off` ⇒ 同类「装配默认值」红，也不是本程碰得到的东西。

处置：本程**不**拿整包超时/红去冒充绿，也**不**用整包红反过来否定 §1 的 C/D 两发——那两发是**定向**跑、各自 `--- PASS`（34.26s／2.84s），且与 §6 的四个内部包一起构成本程真正引用的面。**下一位若要整包 `./cmd/wisp/` 绿，请先修/先裁上面两条非环境红，并在树里没有别的 `go test ./...` 时复跑。**

## 7. 与票面其余格的关系

- AC#0 从「零读数」变成「四发现跑 ＋ 三问都有答案」；**翻不翻勾是编排者的裁定**，本程不自己翻。
- 本程**没有**碰 AC#1/AC#1b/AC#2/AC#3/AC#4 已勾的任何判据，也没有碰 AC#5（分账）与 AC#6（排程）。
- 现树里「盖章 × 续读」的互斥已由甲形（`177-r1`）＋ `175-r2` 化解（`a1` 普查同向）；本格的价值是**把没被选中的丙形价目表钉在盘上**，让下一次有人想走丙时不必重花一次预算。

next=编排者：① 核本表后决定 AC#0 的勾法；② 若要整包 `./cmd/wisp/` 绿，等 `go test ./...` 那条腿收工后复跑（§6.1）；③ 丙形仍未落地，本表不构成「丙已裁」；④ 排程不变：`177 → 175 → 176`。
