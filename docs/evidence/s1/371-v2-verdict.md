# 371-v2 裁决（独立验收·终版，**非实现者**）

- 腿 **`371-v2`**（独立验收者，未参与 `371-r1p`／`371-r2p` 任何实现）。
- 裁决对象：**`be74e2e7b`**（"371-r2p：R6 内层 `nil` 兜底 ＋ R7 边界表"），基线 **`e5a075586`**。
- 量于 **2026-11-05 08:4x–09:1x +08**；所有行号一律取 **`HEAD` blob**（`git show HEAD:<path>`），所有读数**本腿现跑**。
- 写面：**本文件**。产码／测试／票面**零改动**；§③ 的三发变异全走 `-overlay`（**不落盘**），还原核见 §③.4。
- 裁决范围：**R6**（`〔己〕` 内层 `nil` 兜底）、**R7**（`AC#5'` 边界表）、**`M4` 未动**、**无断言被弱化**、**R1–R9 收口矩阵**。

---

## ⓪ 一屏判

| 项 | 判 | 一句话 |
|---|---|---|
| **R6**（内层 `nil` 兜底） | **PASS** | 兜底是**产线码**（`store.go:369-376`），不是注释；摘掉它 ⇒ 假腿/`nil` 通道**静默落盘**；反向牙本腿自跑：pre-R6 blob 在 `store.go:364` **panic**、R6 后**拒**（`result=refused-confirm`、`persisted=ask_every_step`、`saves=0`）（§③）。 |
| **R6 用例的牙** | **PASS（有一枚形上的保留）** | 交付的用例**不修必红**（pre-R6 实测 panic）；但它钉的是"**判定与调用之间通道消失**"这一形，票面原话的"**谓词被改坏**"形只以**仓外 `-overlay` 读数**存在（§③.5 记，非阻断）。 |
| **R7**（边界表） | **PASS** | `docs/evidence/s1/371-r7-boundary.md` 的**每一枚锚点本腿逐条现读复核，全部逐字对得上**；372／375／377／371 **各治一层、判据不可互替**的判断成立；**R9 逐字写成"须摆机主、本笔未做"**（§④）。 |
| **`M4` 未动** | **YES（逐字同一）** | `be74e2e7b` 只碰 **4 枚文件**；`run_mode101_test.go` 的 blob 在 `e5a075586`／`8e0702fdb`／`HEAD` 上 **sha256 全同**＝`81053c1d…`（§⑤）。 |
| **无断言被弱化** | **YES** | 产线只加守卫（+20/−7）；panel 那 1 处是**把 panic 换成 `t.Errorf`**（严格更强）；`store_372_test.go:126-134` 那句注释因 R6 **过期**，但**不是本提交改的**、且**变异仍由别的格子抓住**（§③.5／§⑥）。 |
| **门禁** | **全绿** | build／perm 整包／R6 定向／panel 门禁子集（6 PASS）／`cmd/wisp` 定向（3 PASS）／`d22scan`（复跑 rc=0 clean）／路径预算 GREEN／`gofumpt` 甲乙两尺空／`vet` 静默（§②）。 |
| **剩余要做的** | **恰好 `{R9}`** | R1–R8 **全部有交付物且本腿复跑**；**R9（扩 `M4`）＝机主 2026-09-21 定案，`be74e2e7b` 一字未碰**（§⑤／§⑦）。 |
| **结票建议** | **可结（把 R9 摆机主）** | **ticket 371 implementable scope COMPLETE**，唯一残留是**owner decision**，见 §⑧。 |

---

## ① 复核对象：`be74e2e7b` 到底改了什么（`git show --name-status`）

```
A  docs/evidence/s1/371-r7-boundary.md                    (+142)
M  internal/panel/composer_gate_agreement_371_test.go      (+7 / −1)
M  internal/perm/store.go                                  (+20 / −7)
A  internal/perm/store_371_r6_test.go                      (+93)
```

- **`frontend/**`／`design/**`／`docs/specs/**`／`docs/PLAN.md`／`thresholds.go`／`golden`／`cmd/wisp/**`：命中 `0`**（现跑过滤，见 §⑤）。
- 产线改动只有 `store.go` 的 `GateAsk` 分支：`card, err := s.confirm(ctx, sw)` → `confirm := s.confirm; if confirm == nil { … ResultRefused … }` ＋ `card, err := confirm(ctx, sw)`（`store.go:369-377`）。
  ⇒ **捕获后再调用**：判的值与调的值是同一枚（文档 §0.4 的自述与代码一致）。
- `JudgeCard` 仍是唯一判据（`store.go:348`），守卫**不是第二道门**：只在"`JudgeCard` 保证到不了"的状态上拒绝。

---

## ② 门禁复跑（逐条 `rc` ＋ 逐字；⛔ 相减）

| # | 命令（仓根、Git Bash） | rc | 逐字读数 |
|---|---|---|---|
| 1 | `go build ./internal/perm ./internal/panel ./cmd/wisp` | **0** | 静默、**0 行**输出 |
| 2 | `go test ./internal/perm -count=1 -timeout 200s` | **0** | `ok  github.com/CarlosShao/wisp/internal/perm  0.316s` |
| 3 | `go test ./internal/perm -run 'Test371R6' -count=1 -v -timeout 200s` | **0** | `--- PASS: Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking (0.00s)`；`ok … 0.040s` |
| 4 | `go test ./internal/panel -run '371\|Gate\|Agreement' -count=1 -timeout 200s` | **0** | `ok … 0.045s`；`-v` **6 PASS**（逐名见 ②.1） |
| 5 | `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"; go test ./cmd/wisp -run 'Test371' -count=1 -timeout 220s` | **0** | `ok  github.com/CarlosShao/wisp/cmd/wisp  0.263s`；**3 PASS**（逐名见 ②.2） |
| 6 | `sh scripts/d22scan.sh`（**首跑**） | **1** | ⚠ **环境事故、非判据**：`TestBuiltBinaryGoesRedEndToEnd/armed_ban_6_…` 被 `*** Test killed: ran too long (11m0s)` 杀掉，整包墙钟 **33550.847s**（＝本机会话被挂起跨夜，与 08:4x 恢复之间**9h21m** 的空档吻合）。**未出现任何 ban violation**。 |
| 6b | `sh scripts/d22scan.sh`（**复跑**） | **0** | `d22scan: clean - no D22 ban violations; … ban #8 internal/=865 … cmd/=201`（全绿） |
| 7 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（over=65／roster=65／wall interval=0，字节与字符两尺同选 65 条） |
| 8 | `gofumpt -l` **甲＝`HEAD` blob**（3 枚抽到 `.scratch` 后跑） | **0** | **输出空** |
| 8b | `gofumpt -l` **乙＝工作树**（同 3 枚） | **0** | **输出空**（工作树与 `HEAD` blob sha256 相同 ⇒ 两尺同值，见 §③.4） |
| 9（加跑） | `go vet ./internal/perm ./internal/panel` | **0** | 静默 |
| 10（加跑） | `go test ./internal/panel -count=1 -timeout 200s` | **1** | 恰好两枚**既有红**：`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree` —— 与 `371-v1` §⑦ 记的同一对、输入面（`frontend/**`、`internal/ball/tokens.go`、`docs/evidence/s1/c21-native-tokens.md`）**本提交命中 0**，⛔ 本票担。 |
| — | ★`gofumpt` 不在默认 `PATH` | — | 用 `/d/work/base/gopath/bin/gofumpt.exe`（与 `371-v1` 同一枚）；**裸 `gofumpt` 会 rc=127，⛔ 当红读**。 |

### ②.1 `#4 -v` 逐名

```
--- PASS: TestRoutingDoesNotBypassTheNativeModeGate (0.00s)
--- PASS: TestPanelGateAndStoreGateAgreeOnEveryPair367 (0.00s)
--- PASS: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371 (0.00s)
--- PASS: Test371R1PModeConfirmCarriesTheCredential (0.00s)
--- PASS: Test371R1PThePanelGateCallsTheSameCardPredicate (0.00s)
--- PASS: TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree (0.00s)
```

### ②.2 `#5 -v` 逐名

```
--- PASS: Test371R1PTheAnsweredCardJoinsTheTwoAuditRows (0.02s)
--- PASS: Test371R1PNoClickMeansNoAnsweredCard (0.16s)
--- PASS: Test371R1PThePanelAdaptersPreserveTheAnsweredCard (0.00s)
```

### ②.3 R8 基线读数（本腿**自己复跑**，⛔ 抄 `371-a1p-impl.md`）

```
371A1P|AC0-2 AE->AH real limb|ask_every_step->ask_high_risk|invoked=0|asked=0|set_err=<nil>|last_result=applied|persisted=ask_high_risk|saves=1
371A1P|AC0-3 AE->AA real limb (positive control)|ask_every_step->auto_approve|invoked=1|asked=1|set_err=<nil>|last_result=applied|persisted=auto_approve|saves=1
```

⇒ `asked=1` 的正控证明"D 的 0 是真零"，⛔ 仪器读不到；`AE→AH` 仍 0 卡（`M4` 未被扩，符合预期）。

---

## ③ R6：**PASS** —— 兜底在产线码里，且**pre-R6 真 panic／R6 后真拒**

### ③.1 交付的用例不是空壳（`store_371_r6_test.go`）

它在**门已判 `GateAsk`（腿非 `nil`）之后**把腿摘掉：`s.now` 的第 2 次调用（= `GateAsk` 里的 `started := s.now()`）置 `s.confirm = nil`，随后断言：

- `calls != 2` ⇒ `t.Fatalf`（**防空转**：将来 `Set` 重排、跳过 `GateAsk` 计时器就**响亮地红**，⛔ 静默通过）；
- `err != nil`、`last.Result == ResultRefused`（且**⛔** `ResultCardNotAnswered`）、detail 含 **`内层 nil 兜底`**、**不含** `GateNoChannel` 的 `没有接入确认通道`；
- 腿被调用 **0** 次、`persisted=ask_every_step`／`saves=0`、审计**恰 1 行**。

### ③.2 反向牙（本腿自做）：pre-R6 blob ＋ **同一枚交付用例** ⇒ **panic**

`-overlay` 把 `store.go` 换成 `git show 8e0702fdb:internal/perm/store.go`（sha256 `059d4df6…`，仓内产码零变异）：

```
$ go test -overlay=".scratch/371v2/overlay-pre-r6.json" ./internal/perm \
    -run 'Test371R6' -count=1 -v -timeout 200s
rc=1

=== RUN   Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking
--- FAIL: Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal 0xc0000005 code=0x0 addr=0x0 pc=0x7ff77c4dca75]
github.com/CarlosShao/wisp/internal/perm.(*Store).Set(0xc8534b763c0, {…}, 0x2, …)
	D:/work/workspace/projects plans/Wisp/internal/perm/store.go:364 +0x1b5
github.com/CarlosShao/wisp/internal/perm.Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking(…)
	D:/work/workspace/projects plans/Wisp/internal/perm/store_371_r6_test.go:56 +0x19e
FAIL	github.com/CarlosShao/wisp/internal/perm	0.043s
```

**红行＝`store.go:364`**，pre-R6 blob 上逐字是：

```
		started := s.now()
		card, err := s.confirm(ctx, sw)
```

⇒ 与 `371-r7-boundary.md` §4.1 引的**同一行、同一 `pc`/偏移**（`+0x1b5`），★**独立复现、⛔ 抄**。

### ③.3 票面 `F4` 的原始动机（"**谓词被改坏** ⇒ 拒、⛔ panic"）：三发读数

造一枚**真·坏谓词**（`direction.go` 的 `JudgeCard` 把 `if !channel` 改成 `if false && !channel` ⇒ `auto_approve` + `nil` 通道也返回 `GateAsk`），配一枚**仓外探针**（`-overlay` 挂进 `internal/perm`，**不落任何仓内测试文件**）：

| 发 | 代码 | 读数 | rc |
|---|---|---|---|
| **①** 坏谓词 ＋ **pre-R6** `store.go` | `overlay-pre-r6-broken.json` | **panic**：`invalid memory address or nil pointer dereference` @ `store.go:364 +0x199`（`FAIL 0.044s`） | **1** |
| **②** 坏谓词 ＋ **R6 后** `store.go` | `overlay-after-broken.json` | **拒**（下引逐字） | **0** |
| **③** 原谓词 ＋ R6 后 ＋ 交付用例 | 无 overlay | `--- PASS`（§③.1） | **0** |

②的逐字：

```
371V2|broken-predicate|err=perm: 从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，但确认通道为 nil（内层 nil 兜底：门谓词被改坏时拒而不 panic，未切）
371V2|broken-predicate|result="refused-confirm"|detail="从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，但确认通道为 nil（内层 nil 兜底：门谓词被改坏时拒而不 panic，未切）"
371V2|broken-predicate|persisted="ask_every_step"|saves=0|last_ok=true|rows=[perm: MODE-SWITCH from=ask_every_step to=auto_approve … asked=0 card="" result=refused-confirm detail="…（内层 nil 兜底…）"]
--- PASS: Test371V2BrokenPredicateWithNilChannel (0.00s)
```

⇒ **判据成立**：谓词被改坏不再能打崩进程；`saves=0`＋一行审计 ⇒ 拒绝是**落账的**拒绝，⛔ 静默吞掉。

### ③.4 还原核（时机：**最后一发变异之后**；`-overlay` 从不落盘，故核"盘上本来就没动"）

```
$ git status --porcelain -- internal/perm internal/panel cmd/wisp
?? cmd/wisp/zz_probe_temp_test.go.txt          ← 2026-11-03 既有未跟踪件，非本腿

store.go                           work=206531b8…  HEAD=206531b8…  SAME
direction.go                       work=08663161…  HEAD=08663161…  SAME
store_371_r6_test.go               work=b1c6f143…  HEAD=b1c6f143…  SAME
composer_gate_agreement_371_test.go work=e3baefb2… HEAD=e3baefb2…  SAME
composer_handlers.go               work=d8eb3c3b…  HEAD=d8eb3c3b…  SAME
run_mode101_test.go                work=81053c1d…  HEAD=81053c1d…  SAME
```

### ③.5 两枚**形上的保留**（非阻断，⛔ 藏）

1. **交付用例钉的是"通道在判定与调用之间消失"，⛔ 票面字面的"谓词被改坏"。**
   票面 `〔己〕` 要求"一枚『谓词翻成 `to < from` ⇒ 拒而不 panic』的用例（`-overlay` 变异形、`371-a1` 已实证过红）"。
   交付的用例**不依赖 overlay**（这是对的，入库用例不该依赖变异），走的是同一枚守卫；**"坏谓词"那一形**只以 `371-r7-boundary.md` §4.3 的**仓外探针读数**存在。
   ★本腿 §③.3 已**独立复现**该形（`rc=1` panic → `rc=0` 拒）⇒ **功能诉求满足、判据可达**；保留点仅在于"入库用例的直接触发面"比票面字面窄一档。**建议**（非必须）：后续腿可在入库用例里再加一格"直接构造 `nil` 通道进 `GateAsk`"的变体，把 `-overlay` 那一发也钉进仓内。

2. **`store_372_test.go:126-134` 的注释因 R6 过期（不是本提交改的，⛔ 由 371 改）。**
   逐字：「A mutation that drops the nil half of the gate makes this test panic (a failure, not a pass), which is the property ticket 372 must not lose.」
   ★本腿现跑：**在坏谓词下这枚用例现在 PASS**（`rc=0`）—— 因为 R6 把 panic 换成了拒。**但这不等于覆盖丢失**：

   ```
   $ go test -overlay=".scratch/371v2/overlay-after-broken-panel.json" ./internal/perm -count=1
   --- FAIL: Test371R1PTheNilChannelRefusalIsStillItsOwnVerdict (0.00s)
       store_371_r1p_test.go:185: the nil-channel refusal does not name the channel: perm: …（内层 nil 兜底…）
   FAIL  github.com/CarlosShao/wisp/internal/perm  0.311s
   ```

   ⇒ 同一发变异在 `internal/perm` 仍**逐名报红**；`internal/panel` 侧同样红 5 处（`composer_gate_agreement_367_test.go:90/94/105/108`、`…371_test.go:246/250/254/283`、`composer_dispatch_test.go:185`）。
   ★**结论**：**没有断言被弱化**（一行断言文本都没改；`store_372_test.go` 根本不在本提交的改动面里），**且坏谓词仍被抓住**；留下的是一句**属于票 372 的过期注释**，按票面 `AC#4`（⛔ 碰 372 任何原句）**本票不该改它** ⇒ 记在这里、由编排者处置。

---

## ④ R7：**PASS** —— 边界表**每枚锚点本腿逐条现读**，全部对得上

文件：`docs/evidence/s1/371-r7-boundary.md`（142 行，`be74e2e7b` 新增）。它要答的两件事是"**谁治哪一层**"与"**R6 改前/改后读数**"。

### ④.1 归属（票面要求的四问，逐条对）

| 票 | 文档的判 | 本腿现读复核（`HEAD` blob） | 判 |
|---|---|---|---|
| **372 读侧降级** | "不可问 ⇒ 有效档塌到 `ask_every_step`" | `store.go:274` 逐字 `if s.confirm == nil && configured != risk.ModeAskEveryStep {`；`MODE-DOWNGRADE` 落 `:241-247`，`ResultDowngraded` | **对** |
| **375 动作层 fail-closed** | "`answer()==false` ⇒ L1 在开窗前拒" | `internal/agent/approval/reach.go:74-79` `func (g *Gate) answerable() bool { if g.answer == nil { return true }; return g.answer() }`；`gate.go:346-350` `L1-UNANSWERABLE … decision=reject`；`reach.go:132/150` `MODE-UNANSWERABLE` | **对** |
| **377 记账面** | "`tools: call` 行报**有效档**，是读数不是门" | `internal/tools/bridge.go:286` `mode := b.permissionMode()`、`:1215` `… reason=%q mode=%s`；真值来自 `cmd/wisp/approval_always.go:173` `return rt.modes.PermissionMode()` | **对** |
| **371 门要求"被答过"** | "门放行后要腿交回**被答过 allow 的凭据 `Card`**" | `store.go:384` `if !card.Answered() {`；`:389-393` `ResultCardNotAnswered`；`Card` 本体 `:138-160`（**零值＝没答过**）；`JudgeCard` 调用 `:348`、定义 `direction.go:57` | **对** |
| **R6** | "谓词被改坏 ⇒ 拒而不 panic（`:369-376`）" | `store.go:369-376` 逐字＝`confirm := s.confirm` ＋ `if confirm == nil { … 内层 nil 兜底 … }` | **对** |

### ④.2 "三条判据不可互替"成立（文档 §2 的推理，本腿验过判据面）

- 372 读的是**腿存不存在**（`s.confirm == nil`），**不触发**非 `nil` 假腿（文档 §2 第 2 行）；
- 375 读的是**此刻有没有人能答**（`Gate.answer()`），**不碰 L2**（`gate.go` 的 L2 路径未动，本腿 §② 复跑 375 相关包外的 `gate` 测试未涉及）；
- 371 读的是**这张卡有没有被答过允许**（`card.Answered()`）。
⇒ 三者值不同、落点不同（读侧／动作侧／档位切换），**没有一个能推出另一个**；文档"本票没有把 372／375／377 的收益算成自己的"这句**与代码一致**。

### ④.3 R9 是"明写剩余"（票面 `AC#5'`④ 的硬要求）

- 文档 §0.6 逐字：**"仍未做：R9（〔甲〕扩 M4，让所有变宽出卡）＝机主 2026-09-21 定案，须摆机主，本笔未碰"**；
- 文档 §3.1 逐字：**"这是机主 2026-09-21 17:1x 的定案 ⇒ 须摆机主，本笔未做、也未自决"**，并点名 `run_mode101_test.go:240/:259` **逐字未动**；
- 文档 §3 另记三条**诚实边界**：凭据的**同进程信任边界**、panel **刻意不判"答没答"**、372／375 行为**一字未改**。
⇒ **R7 = PASS**：它没有把"本票做完了"写成"整票做完了"，也**没有**把 R9 藏进已实现项。

---

## ⑤ `M4` 未动：**YES，逐字同一**（`be74e2e7b` 的裁决依据）

```
$ git show --name-only --format="" HEAD | grep -E '^(frontend/|design/|docs/specs/|docs/PLAN.md|cmd/wisp/run_mode101_test.go)|thresholds|golden'
（无输出 ⇒ NONE）

$ git show <c>:cmd/wisp/run_mode101_test.go | sha256sum
e5a075586  81053c1d0ce06ab645adbe7a3fa662d0db950387a6b4efa9f082c9b612bce402
8e0702fdb  81053c1d0ce06ab645adbe7a3fa662d0db950387a6b4efa9f082c9b612bce402
HEAD       81053c1d0ce06ab645adbe7a3fa662d0db950387a6b4efa9f082c9b612bce402
```

⇒ **`M4`（"只有切进 `auto_approve` 吃卡"）未被扩、未被弱化、未被注释漂移**；`cmd/wisp` 整包**不在本提交写面**（`00` 命中）⇒ `AE→AH` 仍 0 卡（§②.3 现跑 `invoked=0`）。

★**并且 R6 的守卫没有渗进中间档**：`GateAsk` 只由 `RequiresCard(to)==true`（＝`auto_approve`）到达。

---

## ⑥ "没有断言被弱化"：逐文件核（本提交只碰 4 枚）

| 文件 | 改动 | 判 |
|---|---|---|
| `internal/perm/store.go` | 只加守卫（+20/−7），**删的是那段"故意 panic"的注释**，换成"改成拒绝"的理由 | **不是弱化**：判据（`JudgeCard`）未动，守卫只拒"`JudgeCard` 保证到不了"的状态 |
| `internal/perm/store_371_r6_test.go` | 新增 93 行，**只加不减** | 新牙（§③） |
| `internal/panel/composer_gate_agreement_371_test.go` | +7/−1：`if !strings.Contains(handlerErr.Error(), …)` → **先判 `handlerErr == nil` ⇒ `t.Errorf`**，否则原断言逐字保留 | **严格更强**：原写法在"该拒却没拒"时**空指针 panic**（`371-v1` §⑧.4 记的同一枚工艺瑕疵），现在**明确报红且不阻断后续行**；**没有放宽任何条件** |
| `docs/evidence/s1/371-r7-boundary.md` | 新增 | 文档 |

★**反向对照**（R6 的"读侧"没被顺带改坏）：`PermissionMode()` 的降级（`store.go:274`）、`MODE-DOWNGRADE`、`answerable()`、`UnsilencedMode` **均未在本提交的 diff 里**；`internal/perm` 整包与 `internal/panel` 门禁子集**全绿**（§②）。

---

## ⑦ R1–R9 收口矩阵（票面"剩余价值 AC 清单"，逐条对）

| # | 诉求（票面原话摘要） | 交付物（盘上真在） | 本腿复核读数 | 判 |
|---|---|---|---|---|
| **R1** | `AC#0'` 现读所有"变宽／吃卡"谓词与调用点，逐处给"判什么属性" | `docs/evidence/s1/371-a1p-impl.md` §④ | **本腿自跑** `git grep 'JudgeCard(\|RequiresCard(\|IsWidening(' HEAD -- '*.go' ':!*_test.go'`：产线调用点**只有 2 处**（`store.go:348`、`composer_handlers.go:176`），**两处都传 `channel != nil`**；`RequiresCard` 只认 `AA`（`direction.go:23`）；`IsWidening` **产线 0 调用**（`direction.go:78`） ⇒ **全部只判 `nil-ness`，无例外** | **DONE** |
| **R2** | 门放行前要一枚"**这张卡被答过**"的凭据，⛔ 只认 `confirm != nil`（**本票主体**） | `internal/perm/store.go:138-160`（`Card`／`AnsweredCard`／`Answered()`）、`cmd/wisp/run.go:1163`（唯一产线铸造点，逐字 `return perm.AnsweredCard(corr), nil`）、判据 `store.go:384` | `Test371R1PTheAnsweredLimbAppliesAndTheRowCarriesTheCard` PASS（`asked=1`）；假腿 `card-not-answered`（`371-v1` PASS，本腿整包复跑绿） | **DONE** |
| **R3** | `MODE-SWITCH` 审计行带"问没问过"读数，让洞可被断言 | `Switch.Asked`／`Switch.Card`（`store.go:113-122`）、赋值 `:395`、打印 `:489-495`（`asked=%d card=%q`） | `Test371R1PTheAnsweredCardJoinsTheTwoAuditRows` **PASS**（`cmd/wisp`，⚠ 0.02s） | **DONE** |
| **R4** | 验收靶改到 `AA`：`AA` ＋ "非 `nil` 但不问" ⇒ `applied=0` ＋ 拒因点名；正控 `AA` ＋ 真腿 ⇒ `confirm_calls=1` | `store_371_r1p_test.go`；读数 `371-a1p` §⑤／`371-r1p-impl.md` §② | `371-v1` 判 **PASS**；本腿 `go test ./internal/perm` **rc=0**（整包）；正控 `asked=1`（§②.3） | **DONE** |
| **R5** | 两处门同判（`perm`／`panel`），"只改一处必须转红" | `composer_gate_agreement_371_test.go`（18 行矩阵）＋ `Test371R1PThePanelGateCallsTheSameCardPredicate` | `371-v1` 判 **PASS**；本腿 `-run '371\|Gate\|Agreement'` **6 PASS**（§②.1） | **DONE** |
| **R6** | 内层 `nil` 兜底：谓词被改坏 ⇒ **拒、⛔ panic** | `store.go:369-376` ＋ `store_371_r6_test.go` | **本次裁决 PASS**（§③，panic-before／refuse-after 三发读数） | **DONE** |
| **R7** | 与 372／375／377 的边界逐条写清"谁做了什么、本票还剩什么" | `docs/evidence/s1/371-r7-boundary.md` | **本次裁决 PASS**（§④，锚点逐条现读） | **DONE** |
| **R8** | 基线读数 ＋ 全门禁甲乙两尺 ＋ 纪律格（⛔ 改 `M4`、⛔ 碰 372／375 原句） | `371-a1p-impl.md` §②③ ＋ 各腿 §门禁 | **本腿自跑**：基线 §②.3、门禁 §②（含 **d22scan 复跑 clean**、路径 GREEN、`gofumpt` 甲乙空、`vet` 静默）；`M4` blob 三提交同 `81053c1d…`；票面／372／375 **零改动**（`git status` 空） | **DONE** |
| **R9** | **〔甲〕扩 `M4`：所有"变宽"都出卡** | **无（也不该有）** | 机主 **2026-09-21 17:1x** 定案（`371-a1`／`A1210①`；`PLAN.md`／`specs` **0 命中** ⇒ 非 `D` 级文字，但**是机主亲拍**）⇒ **须摆机主**；`run_mode101_test.go` blob **逐字未动**（§⑤） | **OWNER（唯一残留）** |

★**"恰好 `{R9}`"的判据**：
1. 票面九条**没有第十条**（`371-*.md:234-242` 的清单就是 R1–R9）；
2. R1–R8 的本腿复核**全部有盘上交付物且现跑绿**（上表"本腿复核读数"列）；
3. R9 是**唯一**被票面标"待人拍板／⛔ 本票自决"的一条，且**它要改的东西（`M4`）本提交一字未碰**；
4. 票面 `AC#4` 的纪律（⛔ 翻框／⛔ 改 `Status:`／⛔ `-done`／⛔ push／⛔ 碰 367／372／375 原句）**全部守住**。

---

## ⑧ 残留（⛔ 藏；均**不阻断** R6／R7／结票建议）

1. **R9（唯一真残留，owner）**：扩 `M4` ＝ 改机主定案。**必须摆机主**，且报票时**先讲"界面上哪处功能"**（把档位从"每步都问"拨到"只问高危"时，面板要不要也弹一次确认卡）。本票**不自决**。
2. **凭据的同进程信任边界**（`store.go:138-146` 逐字 `THE LIMIT`）：挡**事故性**假腿，⛔ 挡蓄意 `AnsweredCard`。跨进程/跨用户证据不在 371 射程；与 `Options.Confirm` 原有边界同层，本腿判**可接受**。
3. **`store_372_test.go:126-134` 的注释已过期**（§③.5.2）：坏谓词下那枚用例现在 PASS，覆盖由 `Test371R1PTheNilChannelRefusalIsStillItsOwnVerdict` ＋ panel 侧 5 处红承接。按 `AC#4` **本票不该改 372 的原句** ⇒ 请编排者决定是否另开一笔（属 372 写面）。
4. **`R6` 入库用例的触发面**（§③.5.1）：钉"通道在判定与调用之间消失"，票面字面的"谓词被改坏"只以仓外 `-overlay` 读数存在。功能与判据均达，仅建议后续补一格。
5. **`371-r1p-impl.md` 仍是未跟踪件**（`git status` 显示 `??`）—— 前一条腿的证据没入库；`371-a1p-impl.md`／`371-a1p` 的探针件同理是历史状态。**不属本提交**，但会让"盘上三把尺"里的"件真在"依赖工作树而非历史。
6. **整包 `internal/panel` 的两枚既有红**（`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`）：输入面 `frontend/**`／`internal/ball/tokens.go` 属机主外部 agent，本提交命中 0 ⇒ **⛔ 本票担、也⛔ 当本票的红**（§② #10）。

---

## ⑨ 结论（★给编排者的一句话）

> **R6 = PASS，R7 = PASS，`M4` 逐字未动，R1–R8 全部有交付且本腿现跑复绿 —— ticket 371 implementable scope COMPLETE，唯一残留是 R9（扩 `M4` 让所有变宽出卡）＝机主 2026-09-21 定案，属 owner decision。建议：把 R9 以"面板上把档位拨松要不要也弹卡"摆机主，同时把票 371 结掉（`-done` 由编排者翻，本腿 ⛔ 动票面）。唯一需要编排者另处置的两件事都不在 371 的可实现面：`store_372_test.go:126-134` 的过期注释（372 写面）与 `371-r1p-impl.md` 未入库。**

★**反向话（若不同意结票）**：只有一种情形——把 §③.5.1 的"入库用例触发面窄一档"判成 R6 **未达标**。本腿**不这么判**：票面 `〔己〕` 的判据是「**拒**、⛔ panic」，本腿已用**真·坏谓词**独立复现"pre-R6 panic（`rc=1`）→ R6 后拒（`rc=0`、`refused-confirm`、`saves=0`）"，**功能诉求与判据都在**；缺口只是"入库用例 vs 仓外探针"的**形**，属加牙建议、⛔ 结票阻断。
