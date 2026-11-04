# 371-v1 裁决（独立验收，**非实现者**）

- 腿 **371-v1**（独立验收者，未参与 `371-r1p` 实现）。裁决对象：修复提交 **`8e0702fdb`**，基线 **`8f03abd8b`**，量于 2026-11-04 22:56–23:10 +08。
- 裁的 AC：票 `371-*.md` 的 **`AC#2'`** 与 **`AC#2'-b`**（并复核 `AC#0'`／`AC#1'`／`AC#5'` 的相邻边界）。
- 写面：**本文件**。产码／测试／票面**零改动**；两枚反向变异腿各自还原（见 §④／§⑤.4，含 sha256 核）。
- ★所有行号一律取 **`HEAD` blob**（`git show HEAD:<path>`）；所有读数**本腿现跑**，⛔ 抄 `371-r1p-impl.md`。

---

## ⓪ 一屏判

| 项 | 判 | 一句话 |
|---|---|---|
| **`AC#2'`** | **PASS** | `AA` ＋ 非 `nil` 但不问的假腿 ⇒ `applied=0`（`persisted=ask_every_step`／`saves=0`）且拒因逐字含 **`card was not answered`**；该句与 `nil` 通道那句**逐字不同**（§③）。 |
| **`AC#2'-b`** | **PASS** | 两处门**同调** `perm.JudgeCard`（`store.go:348`／`composer_handlers.go:176`）；never-asking 行在**两份独立记录**上都读到 `result=card-not-answered`（§⑤）。 |
| **`M4` 未动** | **YES（逐字同一）** | `run_mode101_test.go` 全文件 diff **只有 1 个 hunk**、纯夹具；四行产品断言在基线与修复 blob 上 **sha256 相同**（§②）。 |
| **无断言被弱化** | **YES** | 6 枚夹具文件只改腿的返回值形状；`store_371_a1p_test.go` 的 D/F 翻面是**基线文件自己预先授权的**，A/B 断言未动，新增断言更严（§⑥）。 |
| **反向牙** | **RED（正确方向）** | 摘掉"答过"那一跳 ⇒ `applied=1` 静默，`store_371_r1p_test.go:97` 逐名报红；还原后 sha256 回 `059d4df6…`（§④）。 |
| **门禁** | **全绿** | build／perm 整包／panel 门禁子集／`cmd/wisp` 定向／`d22scan`／路径预算／`gofumpt` 全 `rc=0`；整包 `internal/panel` 仅剩 2 枚**与本票无关**的既有红（§⑦）。 |
| **票能否结** | **本腿 AC 可结；整票⛔** | `AC#2'`／`AC#2'-b` 已满足 ⇒ 可进结案；但 **R9（〔甲〕扩 `M4`）仍待人拍板**、**R6（〔己〕内层 `nil` 兜底）仍未实现**，两者都不属本腿 AC，须由编排者具名处置（§⑧）。 |

---

## ① 复跑读数（逐条 `rc` ＋ 逐字，⛔ 相减）

| # | 命令（`cd` 到仓根、Git Bash） | rc | 逐字读数 |
|---|---|---|---|
| 1 | `go build ./internal/perm ./internal/panel ./cmd/wisp` | **0** | 静默（无输出） |
| 2 | `go test ./internal/perm -count=1 -timeout 200s` | **0** | `ok  github.com/CarlosShao/wisp/internal/perm  0.370s` |
| 3 | `go test ./internal/panel -run '371\|Gate\|Agreement\|ComposerGate' -count=1 -timeout 200s` | **0** | `ok  github.com/CarlosShao/wisp/internal/panel  0.053s`（首跑）／`0.048s`（反向牙还原后复跑）；`-v` 逐枚 **6 PASS**（见下） |
| 4 | `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"; go test ./cmd/wisp -run 'Test371\|TestTicket367\|TestTicket101' -count=1 -timeout 360s` | **0** | `ok  github.com/CarlosShao/wisp/cmd/wisp  135.460s` |
| 4b | 同上 `-v`（复跑取逐名） | **0** | `ok … 134.344s`；**7 枚 PASS** ＋ 3 子格 PASS（见下） |
| 5 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; … ban #8 internal/=864, cmd/=201` |
| 6 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（tracked=18625、over=65、covered=65、not in roster=0） |
| 7 | `/d/work/base/gopath/bin/gofumpt.exe -l <6 枚产线 .go>` | **0** | **输出空**（6 枚：`cmd/wisp/panel_inbound.go`／`resident_approval_windows.go`／`run.go`／`internal/panel/composer_handlers.go`／`internal/perm/direction.go`／`store.go`） |
| 8（加跑） | `go vet ./internal/perm ./internal/panel` | **0** | 静默 |

### ①.1 `#3` 的 `-v` 逐名（还原后复跑）

```
--- PASS: TestRoutingDoesNotBypassTheNativeModeGate (0.00s)
--- PASS: TestPanelGateAndStoreGateAgreeOnEveryPair367 (0.00s)
--- PASS: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371 (0.00s)
--- PASS: Test371R1PModeConfirmCarriesTheCredential (0.00s)
--- PASS: Test371R1PThePanelGateCallsTheSameCardPredicate (0.00s)
--- PASS: TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree (0.00s)
ok  github.com/CarlosShao/wisp/internal/panel 0.048s
```

### ①.2 `#4b` 的 `-v` 逐名（`cmd/wisp`）

```
--- PASS: Test371R1PTheAnsweredCardJoinsTheTwoAuditRows (0.02s)
--- PASS: Test371R1PNoClickMeansNoAnsweredCard (0.16s)
--- PASS: Test371R1PThePanelAdaptersPreserveTheAnsweredCard (0.00s)
--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.48s)
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.53s)
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (81.62s)
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.73s)
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.72s)   ← 3 子格：存储损坏／版本不认识／权限读不到
```

★**一处要说清的读数**：`-run` 里的 **`TestTicket367` 臂在 `cmd/wisp` 命中 0 枚**（`grep -rn '^func Test.*367' cmd/wisp/*_test.go` 无输出；`-v` 里也 0 行）⇒ 367 的一致性格**不在** `cmd/wisp`，而在 `internal/panel` 的 `TestPanelGateAndStoreGateAgreeOnEveryPair367`，由 `#3` 跑到并 **PASS**。⛔ 把"该臂没红"读成"367 有读数"。

---

## ② `M4` 未动：**YES，逐字同一**（裁决依据）

### ②.1 整文件 diff = 1 个 hunk，且只碰夹具

```
$ git diff 8f03abd8b 8e0702fdb -- cmd/wisp/run_mode101_test.go | grep -c '^@@'
1
$ git diff … | grep '^@@'
@@ -221,13 +221,15 @@ func TestTicket101ManualSwitchSurvivesRestart(t *testing.T) {
```

该 hunk 的全部内容（逐字）：`confirm := func(_ context.Context, sw perm.Switch) error {` → `(perm.Card, error)`；
`return nil` → `return perm.AnsweredCard("t101"), nil`；外加 2 行注释。**没有一行断言被碰。**

### ②.2 产品断言四行：两枚 blob 上 **sha256 相同**

抽出的四行（注释行 ＋ `want exactly [auto_approve]` 那条 ＋ `if len(gotAsked)…` ＋ 第三条 `cards shown` 断言）：

```
		// R20/M4 is asymmetric: only the switch INTO auto_approve costs a card.
	if len(gotAsked) != 1 || gotAsked[0] != risk.ModeAutoApprove {
		t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", gotAsked)
		t.Errorf("cards shown = %d, want exactly one L2 card per switch (R20/M4)", cards)
```

| blob | 抽取行 | sha256 |
|---|---|---|
| `8f03abd8b:cmd/wisp/run_mode101_test.go`（`:238`／`:256-257`／`:564`） | 同上 | `92026504b0f3ce46b211575fd159192e2648a013b411d758371f10bc26cf0731` |
| `8e0702fdb:cmd/wisp/run_mode101_test.go`（`:240`／`:258-259`／`:566`） | 同上 | `92026504b0f3ce46b211575fd159192e2648a013b411d758371f10bc26cf0731` |

`diff` 两枚抽取文件 `rc=0`（无差异）。行号 `+2` 的位移**只**来自那两行新注释。

### ②.3 工作树与 `HEAD` 同 blob

```
81053c1d0ce06ab645adbe7a3fa662d0db950387a6b4efa9f082c9b612bce402 *cmd/wisp/run_mode101_test.go
81053c1d0ce06ab645adbe7a3fa662d0db950387a6b4efa9f082c9b612bce402 *-
```

⇒ **`M4`（"只有切进 `auto_approve` 吃卡"）未被扩、未被弱化**；`TestTicket101ManualSwitchSurvivesRestart` 本腿现跑 **PASS（41.48s）**，它在产线装配上仍断言 `gotAsked == [auto_approve]`。

---

## ③ `AC#2'`：**PASS** —— 洞确实闭合，且**两种拒因分得开**

### ③.1 靶心读数（`go test ./internal/perm -v`，本腿现跑）

```
371A1P|AC0p-2 AE->AA fake limb (the hole)|ask_every_step->auto_approve|invoked=1|asked=0|set_err=perm: 切到 auto_approve 被拒绝: perm: card was not answered (no L2 card was shown and answered allow)|last_result=card-not-answered|persisted=ask_every_step|saves=0
```

| 票面判据 | 读数 | 判 |
|---|---|---|
| `applied=0` | `persisted=ask_every_step`、`saves=0`、`last_result=card-not-answered` | ✅ |
| 拒因**逐字**点名"这张卡没被答过" | `card was not answered`（`strings.Contains` 真、`errors.Is(err, ErrCardNotAnswered)` 真） | ✅ |
| **⛔ 只说"拒了"** | 与 `nil` 通道拒因**逐字不同**（见 ③.2） | ✅ |

### ③.2 两枚拒因**同时**取到（`-overlay` 聚焦探针，落盘零改动）

探针放在 `.scratch/371v1/probe_371v1_test.go`，经 `-overlay` 挂进 `internal/perm`，**不落任何仓内测试文件**：

```
371V1|nil-channel|err="perm: 从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）"|result=refused-confirm|recorded=true|persisted_saves=0
371V1|never-asking|err="perm: 切到 auto_approve 被拒绝: perm: card was not answered (no L2 card was shown and answered allow)"|result=card-not-answered|recorded=true|persisted=ask_every_step|saves=0|invoked=1|asked=0
371V1|distinct=true|contains_card_was_not_answered(fake)=true|contains_card_was_not_answered(nil)=false
371V1|fake_applied=false|nil_applied=false
```

**两枚拒因逐字**：

1. `nil` 通道（`store.go:351-352`，`ResultRefused`）：
   `perm: 从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）`
2. 非 `nil` 但不问（`store.go:371-380`，`ResultCardNotAnswered`）：
   `perm: 切到 auto_approve 被拒绝: perm: card was not answered (no L2 card was shown and answered allow)`

⇒ **两句话不同、两个 `Result` 不同、`errors.Is` 交叉为假** ⇒ 票面担心的"与 `nil` 通道的拒分不开"**没有发生**。

### ③.3 正控（⛔ 省）

```
371A1P|AC0p-4 answered_limb|refused=false|result=applied|persisted=auto_approve|saves=1|asked=1
371A1P|AC0p-4 never-asking_limb|refused=true|result=card-not-answered|persisted=ask_every_step|saves=0|asked=0
```

⇒ `asked=0` **不是假零**：同一尺对"真问过的腿"读出 `asked=1` 且 `applied=1`；两枚腿在 `asked` 上分叉、在判决上分叉。
另：`Test371R1PTheAnsweredLimbAppliesAndTheRowCarriesTheCard` PASS，审计行带 `asked=1 card="host:mode-switch#7"`。

### ③.4 `M4` 隔壁没被扩（`AC#2'` 的反向边界）

`Test371R1PTheMiddleTargetIsStillFreeOfCards` PASS：`AE→AH` 下 `invoked=0`、`persisted=ask_high_risk`、行里 `asked=0`。
⇒ 新的"答过"判据**只**长在 `GateAsk`（`auto_approve`）分支上，没有渗到中间档。

---

## ④ 反向牙（本腿自己做，⛔ 信实现者）：**RED，且逐名报红行**

变异：`internal/perm/store.go:371` `if !card.Answered() {` → `if false && !card.Answered() {`
（用 `-overlay`：把 `HEAD` blob 复制到 `.scratch/371v1/mut/store.go` 再改，**工作树一字节未动**）

```
$ go test -overlay .scratch/371v1/overlay-reverse.json ./internal/perm \
    -run 'Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt|Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied|Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs' -count=1 -v
rc=1

    store_371_a1p_test.go:177: 371A1P|AC0p-2 …|invoked=1|asked=0|set_err=<nil>|last_result=applied|persisted=auto_approve|saves=1
    store_371_a1p_test.go:191: the fake limb was silently accepted; ticket 371's gate is gone - the never-answering limb must be refused (<nil>)
--- FAIL: Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied
    store_371_a1p_test.go:267: the NEVER-ANSWERING limb was not refused as unanswered: got {refused:false result:applied persisted:auto_approve saves:1 asked:0}, want refused with "card-not-answered". If the two limbs read the same verdict again, ticket 371's gate is gone
    store_371_a1p_test.go:272: both limbs left the same persisted state ("auto_approve"/1): the gate did not split them
--- FAIL: Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs
    store_371_r1p_test.go:97: the switch applied although no card was ever answered: ticket 371's hole is open
--- FAIL: Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt
FAIL	github.com/CarlosShao/wisp/internal/perm	0.044s
```

★**红形与票面要求逐字一致**：摘掉那一跳 ⇒ 假腿回到 `applied=1` / `persisted=auto_approve` / `saves=1` **静默**，且三条用例**逐名报红行**（`store_371_r1p_test.go:97` 是最直接的一枚）。

**同一变异打到 panel 一致性格也红**（两份记录同时被读坏）：

```
    composer_gate_agreement_371_test.go:160: ask_every_step -> auto_approve / never-asking: panel refused=false (<nil>), want refused=true
    composer_gate_agreement_371_test.go:166: … the refusal still moved the persisted档 to auto_approve
    composer_gate_agreement_371_test.go:197: … the store's own row says refused=false, want true: perm: MODE-SWITCH … asked=1 card="" result=applied detail=""
--- FAIL: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371
rc=1
```

### ④.1 还原核（时机：**最后一发变异之后**）

| 对象 | 读法 | 值 |
|---|---|---|
| `internal/perm/store.go`（变异前后各一次） | 工作树 | `059d4df618468c77738da527a1c23b3360614cb115fbd8a847792f8a616bb871` |
| `git show HEAD:internal/perm/store.go` | blob | `059d4df618468c77738da527a1c23b3360614cb115fbd8a847792f8a616bb871` |
| `git status --porcelain -- internal/perm internal/panel cmd/wisp` | — | 只有 `?? cmd/wisp/zz_probe_temp_test.go.txt`（**2026-11-03 23:46 的既有未跟踪件，非本腿**） |

⇒ `-overlay` 变异**从不落盘**，还原＝核"盘上本来就没动"，比 `git checkout` 更强。

---

## ⑤ `AC#2'-b`：**PASS** —— 两处一门 ＋ 两条记录

### ⑤.1 两个调用点**同调** `perm.JudgeCard`（逐字，`HEAD` blob）

| 层 | 文件:行 | 逐字 |
|---|---|---|
| `perm`（会问的那层） | `internal/perm/store.go:348` | `switch JudgeCard(to, s.confirm != nil) {` |
| `panel`（出卡前拒的那层） | `internal/panel/composer_handlers.go:176` | `if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel {` |

谓词本体一份（`internal/perm/direction.go:57-65`）：`JudgeCard(to, channel)` ⇒ `GateNoCard`／`GateNoChannel`／`GateAsk`，其第一项就是 `RequiresCard`（`:58`）。
**产线里没有第二套拼法**：`git grep 'JudgeCard\|RequiresCard(\|IsWidening(' HEAD -- '*.go' ':!*_test.go'` 的生产命中只有上面两处调用 ＋ `direction.go` 的定义；`IsWidening` 仍是"⛔ 门"（`direction.go:78`，产线调用点 0）。

### ⑤.2 never-asking 行在**两份独立记录**上（本腿探针，`-overlay`，真 `perm.Store` 在真 handler 后面）

```
371V1|store-record|perm: MODE-SWITCH from=ask_every_step to=auto_approve at=… origin="panel" actor="v1" asked=0 card="" result=card-not-answered detail="切到 auto_approve 需要一张被回答过的 L2 强确认卡，但确认通道没有交回任何“有人答过允许”的凭据（card was not answered）：未切"
371V1|panel-record|panel: MODE-REFUSED request="v1-never-asking" from="ask_every_step" to="auto_approve" origin="panel" actor="v1" result=card-not-answered err=perm: 切到 auto_approve 被拒绝: perm: card was not answered (no L2 card was shown and answered allow) detail="档位未变化：写入腿按同一条规则拒绝了这一发"
371V1|handlerErr=perm: 切到 auto_approve 被拒绝: perm: card was not answered (…)|writeCalls=1|limbEnters=1|persisted=ask_every_step
```

⇒ 两条记录**同名** `result=card-not-answered`、同一发请求（`requestId` 可对）、盘上没动；面板**没有**第二个词面。

### ⑤.3 一致性矩阵的覆盖（`composer_gate_agreement_371_test.go`）

- 矩阵＝ `6` 个有序对 × `3` 种腿 ＝ **18 行**，计数断言 `pairs==6 && kinds==18`（枚举被截短即 `t.Fatalf`）。
- never-asking 那一支（`carded && kind == "never-asking"`）**同时**断言：面板放行给 store（`writer.calls==1`、`limb.asked()==1`）、surfaced 拒因含 `card was not answered`、**store 行**含 `result=card-not-answered` ＋ `asked=0`、**panel 行**恰 1 条且含同一 `result=` 与 `requestId`。
- `no channel` 行断言面板**自己**拒（`没有接入 L2 确认腿`）、写入腿 `0` 次、腿没被进过、panel 恰 1 行。
- `answered allow` 行断言 `writeCalls=1`／`asks=1`／**没有**第二条判决行、store 行 `asked=1 card="t371/panel"`。
- 非吃卡目标（`default` 支）断言 `limb.asked()==0` ⇒ **R20/M4 未被扩**。

### ⑤.4 `AC#2'-b` 的反向判据："只改 `perm`、不改 `panel` ⇒ 那一格必须红"

本腿**自己落盘变异**并当场还原（这一枚必须读盘，因为钉子用 `os.ReadFile`）：

变异：`composer_handlers.go:176` → `if perm.RequiresCard(to) && h.Confirm == nil {`

```
    composer_gate_agreement_371_test.go:343: composer_handlers.go calls perm.JudgeCard 0 time(s), want exactly 1: the panel's half of the mode gate must be the SAME judgement perm.Store.Set applies (ticket 371-b), or the two gates are two rules
--- FAIL: Test371R1PThePanelGateCallsTheSameCardPredicate
rc=1
```

★**同一次变异里**：`--- PASS: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371`（行为矩阵**仍绿**）＋ `--- PASS: Test371R1PModeConfirmCarriesTheCredential`。
⇒ 源码钉是**承重的**：两种拼法今天行为等价，行为面永远看不见分叉 —— 这正是 `AC#2'-b` 要防的"某一处碰巧对"。

**还原**：

```
$ git checkout -- internal/panel/composer_handlers.go          # rc=0
d8eb3c3b7fae4efa084257a10d7b692c814be29eb716ec3a324929b0661dd1c4 *internal/panel/composer_handlers.go
d8eb3c3b7fae4efa084257a10d7b692c814be29eb716ec3a324929b0661dd1c4 *.scratch/371v1/composer_handlers.go.keep
d8eb3c3b7fae4efa084257a10d7b692c814be29eb716ec3a324929b0661dd1c4 *-            ← git show HEAD:…
```

还原后复跑 `#3`：**6 PASS、`ok 0.048s`**（见 ①.1）⇒ 还原是行为级干净的。

---

## ⑥ 没有断言被弱化（逐文件核）

**被 `8e0702fdb` 碰到的 17 枚文件**（`git diff --name-only`）：产线 6 枚（`internal/perm/store.go`／`direction.go`、`internal/panel/composer_handlers.go`、`cmd/wisp/run.go`／`panel_inbound.go`／`resident_approval_windows.go`）＋ 夹具 8 枚 ＋ 新用例 3 枚。**`docs/**`／`specs/**`／`thresholds.go`／`golden`／`frontend/**`／`design/**` 命中 0**。

### ⑥.1 六枚夹具：只改腿的**返回值形状**，断言一行未动

`store_test.go`／`ticket90_persist_test.go`／`booked_mode_377_test.go`／`permission_mode_effective_377_test.go`／`panel_mode_confirm_114_windows_test.go`／`composer_handlers_test.go` 的 diff 共 153 行，形态统一：

```diff
-	return func(_ context.Context, sw perm.Switch) error { … return c.fail }
+	return func(_ context.Context, sw perm.Switch) (perm.Card, error) { … return perm.AnsweredCard("…"), nil }
```

`store_test.go` 的 `confirmSpy` 还保留了**三分支语义**（`answer != nil` ⇒ `Card{}, c.answer`；否则 `AnsweredCard("spy")`），即"错误那条腿不带凭据"仍被测到。**没有一条 `if`／`Errorf` 被删或放宽。**

### ⑥.2 `store_371_a1p_test.go` 的翻面：**基线文件自己预先授权**

基线 blob 逐字：

```
// the D/F assertions below flip to applied=0 with a refusal naming the unasked
// card. Flipping that assertion is part of that fix, not a regression here.
```

本腿核：A／B 两枚用例的**断言体未动**（diff hunk 只落在文件头注释、limb 类型、D/F 两枚用例）；D/F 的 limb 侧读数（`invoked=1`／`asked=0`）**逐字保留**，而门侧断言从"applied 就是洞"翻成"必须拒且点名"，并**新增** `errors.Is(ErrCardNotAnswered)`、`asked=0 card=""`、`persisted=ask_every_step/0` 三条 ⇒ **严格更严**。

### ⑥.3 直接调 `Set` 的既有断言一条没弱化

`store_367_test.go`／`store_372_test.go` **根本不在本提交的改动面里**（`store_372_test.go:134` 那句 `Set panicked on auto_approve with no channel` 仍在）；`store_367_test.go:86` 那枚"故意用字面 `p.to == risk.ModeAutoApprove`"的格子未动。

---

## ⑦ 整包 `internal/panel` 的两枚红：**与本票无关**（已证）

```
$ go test ./internal/panel -count=1 -timeout 200s
rc=1
--- FAIL: TestComposerContractTypesMatchFrontend (0.00s)
    composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)
    tokens_fourway_test.go:638: parties: frontend/scripts/gen-tokens.mjs=43 tokens, docs/evidence/s1/c21-native-tokens.md=78 colour rows, internal/ball/tokens.go=37/35 palette fields, … 11 mapped + 23 waived
```

**输入逐枚 `git diff --quiet 8f03abd8b 8e0702fdb -- <path>` 全 `rc=0`**（＝无差异）：

| 输入 | rc |
|---|---|
| `frontend/src/lib/panel.ts` | 0 |
| `internal/panel/composer.go` | 0 |
| `internal/ball/tokens.go` | 0 |
| `docs/evidence/s1/c21-native-tokens.md` | 0 |
| `frontend/src/styles/tokens.generated.css` | 0 |
| `frontend/scripts/gen-tokens.mjs` | 0 |
| 两枚用例自身 ＋ `workspace_test.go` | 0 |

⇒ 两枚红的**输入与判据**都与 `8e0702fdb` 无关；`frontend/**`／`internal/ball/tokens.go` 属机主外部 agent 的写面，⛔ 本票修，⛔ 本票担。

---

## ⑧ 残留与边界（⛔ 藏）

1. **R9／〔甲〕（待人拍板，⛔ 本票自决）**：让**所有变宽**都吃卡 ＝ 扩 **`M4`** ＝ 改机主 2026-09-21 定案（台账 `:3019`／`:3031`）。本提交**没碰**（§② 已证），**仍须摆机主**。这枚不解决，票面"〔甲〕"一行就仍挂着。
2. **R6／〔己〕（仍未实现，非本腿 AC）**：`store.go:355-362` 逐字写着"There is no separate nil test here on purpose"，即**内层 `nil` 兜底仍缺**；票面要求的"谓词被改坏 ⇒ 拒而⛔ panic"那枚用例也没在本提交里出现（`store_372_test.go:134` 只覆盖 `nil` 通道那条老路径）。⇒ **须由编排者具名处置**（实现，或在台账里显式划出本票射程）。
3. **凭据的同进程边界**（`store.go:138-145` 逐字 `THE LIMIT, stated where a reader can see it`）：挡**事故性假腿**，⛔ 挡蓄意 `AnsweredCard`。与 `Options.Confirm` 原有信任边界同层，本腿判**可接受**。
4. **一枚工艺瑕疵（不阻断）**：`composer_gate_agreement_371_test.go:214` 在"store 该拒却没拒"时会 `handlerErr.Error()` 空指针 panic（本腿 §④ 的变异里实测到）。它仍然 `rc=1` 且前三条行已报红，但把 `handlerErr` 先判空再取值会更干净。**建议**（非必须）由后续腿加一枚 `t.Fatalf` 守门。
5. **`AC#5'` 边界（复核结论）**：372 的读侧降级（`PermissionMode` 只看 `s.confirm == nil`）与 375 的动作层 fail-closed **都**只治"不可问"；本提交治"可问而没问"，**三层不重叠**。377 的记账面未另造 —— `[bing-b]` 的 `asked=`／`card=` 加在**同一条 `MODE-SWITCH`** 行（字段赋值 `store.go:382`，打印 `store.go:471-484` 的 `auditAs`），且实测能与 `ANSWER-ALLOW corr=` join。
6. **`-run 'TestTicket367'` 在 `cmd/wisp` 命中 0 枚**（①.2 末段）：这不是红，但也⛔ 当绿读；367 的格在 `internal/panel`，本腿已复跑 PASS。

---

## ⑨ 结论

- **`AC#2'` = PASS**：`AA` ＋ 非 `nil` 但不问的假腿 ⇒ `applied=0`（`persisted=ask_every_step`、`saves=0`）＋ 拒因逐字 `card was not answered`，且与 `nil` 通道那句**逐字不同**；正控 `asked=1` 证明零不是假零；反向牙摘跳即红并逐名报行。
- **`AC#2'-b` = PASS**：`perm.Store.Set` 与 `panel.ModeWriteHandler` **同调** `perm.JudgeCard`（两处逐字已引）；never-asking 行在**store 行与 panel 行两份独立记录**上都是 `result=card-not-answered`；"只改 `perm` 不改 `panel`"由源码钉转红（`:343`），而行为矩阵仍绿 ⇒ 该钉承重。
- **`M4` = 未动（YES，逐字同一）**：`run_mode101_test.go` 单 hunk、纯夹具，四行产品断言 sha256 相同，工作树＝`HEAD`。
- **门禁 = 全绿**（两枚既有 panel 红与本票无关，输入逐枚无差异）。
- **票能否结**：**本腿 AC（`AC#2'`／`AC#2'-b`）可结、`8e0702fdb` 可作数**；但**整票 `371` ⛔ 就此标 `-done`** —— **R9 待人拍板**、**R6 未实现**，两枚都要编排者具名处置后才谈结票。
