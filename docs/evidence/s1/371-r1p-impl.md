# 371-r1p —— `AC#2'` ＋ `AC#2'-b` 落地（形 `[bing-a]` ＋ `[bing-b]` 同波）

- 腿 **371-r1p**（写腿）。基线 **HEAD `8f03abd8b`**，量于 2026-11-04 22:0x–22:5x +08。
- 输入：票 `371-*.md` 的 `AC#1'`／`AC#2'`／`AC#2'-b`／`AC#5'` ＋ 只读腿 `371-a1p` 的
  `docs/evidence/s1/371-a1p-impl.md` ＋ 其探针 `internal/perm/store_371_a1p_test.go`。
- 写面：`internal/perm/**`、`internal/panel/composer_handlers.go`、`cmd/wisp/**`（继续用腿）、
  一枚 `internal/tools` 夹具。**产码 4 个文件 ＋ 2 支适配器所在的 2 个文件**，其余是夹具与新用例。
- ★**没有**碰：`M4`／`cmd/wisp/run_mode101_test.go:238/256-257` 的产品断言、票 367／372／375 原句、
  `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden。

---

## ⓪ 一屏判

1. **形**：`perm.Card` —— 一枚**不可零值伪造**的凭据（`AnsweredCard(corr)` 是唯一铸造点，零值＝“没答过”），
   `ConfirmFunc` 改成 `(Card, error)`；`Set` 只在 `card.Answered()` 时落盘，否则
   `ResultCardNotAnswered` ＋ `ErrCardNotAnswered`（逐字含 `card was not answered`）。
2. **两处一门**：新谓词 `perm.JudgeCard(to, channel)`（`GateNoCard`／`GateNoChannel`／`GateAsk`）
   被 `perm.Store.Set` 与 `panel.ModeWriteHandler` **同调**；文案从“各拼一遍”变“调用一次”。
3. **`[bing-b]`**：`Switch` 加 `Asked` ＋ `Card`，`MODE-SWITCH` 行加 `asked=0|1 card="<corr>"`；
   产线腿在调 gate **之前**自铸 corr（`host:mode-switch#N`）并塞进 `Decision.CorrelationID`，
   于是同一行可与 `approval: ANSWER-ALLOW corr=<同一个 id>` **join**（§④ 实测两行）。
4. **读数**：`AA` ＋ 非 `nil` 但不问的假腿 ⇒ **`applied=0`、`asked=0`、拒因逐字 `card was not answered`**；
   正控 `AA` ＋ 真腿 ＋ 有人答 allow ⇒ **`applied=1`、`asked=1`、`asked=1 card="host:mode-switch#1"`**。
5. **反向牙**三枚全红、且逐名报红行（§⑤）：摘掉 `Set` 里“答过”那一跳 ⇒ 假腿回到 `applied=1` 静默；
   摘掉 panel 的**同判/记账**那一跳 ⇒ 新的两处一致格红；把 panel 的门改回旧拼法（落盘变异）⇒ 源码钉格红
   （而**行为矩阵仍绿**——这条本身就是“为什么那枚钉是承重的”读数）。
6. **门禁**：`build`／`internal/perm`／panel 门禁子集／格式／`d22scan`／路径预算全绿；
   整包 `internal/panel` 与 `cmd/wisp -run '…Mode'` 各剩 **2 枚与本腿输入无关的红**（§⑦ 已证）。

---

## ① 形：怎么表示“答过允许”

### ①.1 选型（以及为什么它最像主流 harness、⛔ 最像“最小动作”）

```go
// internal/perm/store.go
type Card struct { answered bool; corr string }

func AnsweredCard(corr string) Card   // 唯一铸造点：答过 allow，且带卡片的 correlation id
func (c Card) Answered() bool         // 零值 = 没答过
func (c Card) Corr() string           // 供审计行 join

type ConfirmFunc func(ctx context.Context, req Switch) (Card, error)
```

- ★**为什么不是 `(Answer, error)` 三值枚举**：`deny`／`fault` 在今天的契约里已经由 **error** 承载
  （`RefusalResult` 分 `refused-confirm` / `confirm-failed`），再造一个枚举会让同一种“没答成”有两条通道。
  凭据只回答**一个**问题：“这张卡被答过 allow 了吗”，其余仍走 error。
- ★**为什么凭据带 `corr`**：`[bing-b]` 要求审计行能与 `ANSWER-ALLOW corr=` **按 corr** join。
  `PendingApproval` 只返回 `(Answer, why)`，所以**id 必须由腿自己铸**——见 §③ 的产线改法。
  凭据顺手带 id，等于把“答过”与“答的是哪张卡”合成一个事实。
- ★**零值语义**写在类型注释里（`THE ZERO VALUE MEANS "NOT ANSWERED"`），因为它是整个修的承重点：
  旧契约“`nil` error ＝ allowed”与新契约“零 `Card` ＝ not answered”**方向相反**，注释必须点名，
  否则下一个改腿的人会把 `return nil` 写回去（编译能过、洞会重开）。
- ★**信任边界**（照 `371-a1p` §⑧.3 的风险栏写进代码）：凭据是**同进程铸造**，挡的是**事故性假腿**，
  ⛔ 挡蓄意伪造 `AnsweredCard`。这和 `Options.Confirm` 原有的边界同层：`perm` 造不出“问”，
  它只能被交给一次问。

### ①.2 门：一个谓词、两处调用

```go
// internal/perm/direction.go
type CardGate uint8
const ( GateNoCard CardGate = iota; GateNoChannel; GateAsk )
func JudgeCard(to risk.Mode, channel bool) CardGate
```

| 层 | 逐字 | 作用 |
|---|---|---|
| `perm.Store.Set` | `switch JudgeCard(to, s.confirm != nil) { case GateNoChannel: … case GateAsk: 问并验凭据 }` | 权威判决（它才会问） |
| `panel.ModeWriteHandler` | `if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel { … }` | 面板侧**出卡前**的拒（页面拿到面板错误码），并把 Set 的拒按同一条规则记账 |

★`RequiresCard` **没有**退场：`JudgeCard` 的第一项就是它。逐字保留“目标规则只有一份”的
367／372 纪律，只是把**通道项**也收进同一个函数的签名，于是“两处一致”从**约定**变成**同一个调用**。

### ①.3 ★`AC#2'-b` 的诚实边界（写清，⛔ 藏）

**答案本身 panel 判不了，而且不该判**：它**不出卡**（`composer_handlers.go` 的头注释逐字解释：
在这里再问一次就是一次点击两张卡 = 两个 300s 窗口）。所以 panel 这一层能做、也必须做的是：

1. **同一谓词**：`JudgeCard` —— 目标 ＋ 通道项由一次调用决定（②.3 的源码钉把它钉死）；
2. **同一分类**：Set 拒了之后，panel 用 `perm.RefusalResult(err)` 给这一发**记自己的行**
   （`panel: MODE-REFUSED … result=card-not-answered`），与 store 的 `MODE-SWITCH … result=card-not-answered`
   **同名**；行里带 `requestId`，于是“页面这一发为什么没成”在日志里有落点。

⇒ 这就是 `AC#2'-b` 要的“两处同判”的**可判面**：不是“panel 也判答案”（不可能），
而是“**同一个谓词 ＋ 同一套 result 词面 ＋ 两层都留行**”。§⑤.2／§⑤.3 给出“只改 perm 不改 panel”的两种红。

---

## ② `AC#2'` 读数（`internal/perm/store_371_r1p_test.go`）

命令：

```bash
go test ./internal/perm -run 'Test371R1P|Test371A1P' -count=1 -v
```

### ②.1 靶心：`AA` ＋ 非 `nil` 但不问的假腿 ⇒ `applied=0` ＋ 拒因点名

```
--- PASS: Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt
371A1P|AC0p-2 AE->AA fake limb (the hole)|ask_every_step->auto_approve|invoked=1|asked=0
  |set_err=perm: 切到 auto_approve 被拒绝: perm: card was not answered (no L2 card was shown and answered allow)
  |last_result=card-not-answered|persisted=ask_every_step|saves=0
```

| 判据（票面 `AC#2'`） | 读数 | 判 |
|---|---|---|
| `applied=0` | `last_result=card-not-answered`、`persisted=ask_every_step`、`saves=0` | ✅ |
| 拒因**逐字**点名“这张卡没被答过” | `errors.Is(err, perm.ErrCardNotAnswered)` 真；`strings.Contains(err.Error(), "card was not answered")` 真 | ✅ |
| **⛔ 只说“拒了”** | 该拒**不含** `没有接入确认通道`（`nil` 通道那句），`errors.Is(err, ErrCardNotAnswered)` 对 `nil` 通道那一发为**假** | ✅ 两条故事分得开 |
| 审计行同一读数 | `result=card-not-answered asked=0 card=""` | ✅ |

### ②.2 正控：`AA` ＋ 真腿 ＋ 有人答 allow ⇒ `applied=1`、`asked=1`

```
--- PASS: Test371R1PTheAnsweredLimbAppliesAndTheRowCarriesTheCard
audit: perm: MODE-SWITCH from=ask_every_step to=auto_approve … asked=1 card="host:mode-switch#7" result=applied
```

- `invoked=1 asked=1`、`persisted=auto_approve saves=1`、`Last().Asked=true`、`Last().Card=host:mode-switch#7`。
- ★**⛔ 只测 `applied=0`**：这一发同时钉 `asked=1`，所以 ②.1 的 `asked=0` **不是假零**。

### ②.3 邻居没被移动（回归面）

| 用例 | 读数 | 判 |
|---|---|---|
| `…TheNilChannelRefusalIsStillItsOwnVerdict` | `ResultRefused` ＋ 逐字 `没有接入确认通道`；`errors.Is(ErrCardNotAnswered)=false`；`asked=0` | 372 的 fail-closed 原样 |
| `…AnOperatorsNoIsStillARefusalNotAnUnansweredCard` | `ResultRefused`，`errors.Is` 回操作者自己的拒因；**不**是 `card-not-answered`；`saves=0` | 三种故事各自可读 |
| `…TheMiddleTargetIsStillFreeOfCards` | `AE→AH` `invoked=0`、`asked=0`、`persisted=ask_high_risk` | **M4 未被扩** |
| `…RefusalResultIsTheOneVocabulary` | `nil→applied`／`ErrCardNotAnswered→card-not-answered`（含 `%w` 包裹）／`deny→refused-confirm`／`Deadline/Canceled→confirm-failed` | 一处词面，两层共用 |

### ②.4 `371-a1p` 探针：按它自己的注释翻面（⛔ 删、⛔ 弱化）

`store_371_a1p_test.go` 的 D／F 两格在 `786a13b22` 上记录的是**洞**，其注释逐字写了
“the day ticket 371's AC#2' lands … the D/F assertions below flip to applied=0 with a refusal naming
the unasked card. Flipping that assertion is part of that fix, not a regression here.”

⇒ 本腿只翻断言、**不动 limb 侧读数**（`invoked=1`／`asked=0`／`fake 从不 ask` 三条逐字保留）：

```
371A1P|AC0p-4 answered_limb     |refused=false|result=applied          |persisted=auto_approve |saves=1|asked=1
371A1P|AC0p-4 never-asking_limb |refused=true |result=card-not-answered|persisted=ask_every_step|saves=0|asked=0
```

★两枚用例改了名（`…IsSilentlyApplied` → `…IsRefusedNotSilentlyApplied`；`…HasTheSameHole` →
`…SplitsTheTwoLimbs`），因为旧名在新行为下会撒谎；旧名在新注释里点名，便于按 a1p 证据 grep。
A／B／②.1 的三条“稳定读数”逐字未动（A 仍 `invoked=0 asked=0 applied`，B 仍 `asked=1`）。

---

## ③ `[bing-b]`：审计行可 join（`cmd/wisp/perm_371_r1p_test.go`）

### ③.1 产线改法（三行 + 一枚计数器）

`confirmModeSwitchVia` 是**唯一**见过 `approval.Gate` 答案的产线函数，所以：

1. `modeSwitchCardCorr()` 铸 `host:mode-switch#N`（`atomic.Uint64`，进程内唯一）；
2. 塞进 `tools.Decision{CorrelationID: corr, …}` —— `approval/queue.go:166` 逐字
   `corr := d.CorrelationID; if corr == "" { corr = "approval-"+… }`，所以**队列就用这个 id**；
3. 答 allow 时 `return perm.AnsweredCard(corr), nil`，于是 `Set` 把它写进审计行的 `card=`。

★顺带把 D31 的 key 从**常量** `host:mode-switch` 改成**每卡唯一**（旧的常量在并发两发时会同键）。

### ③.2 实测：一发、一允许、两行同一个 id

```
371R1P|join|perm: MODE-SWITCH from=ask_every_step to=auto_approve at=… origin="panel" actor="owner" asked=1 card="host:mode-switch#1" result=applied detail=""
371R1P|join|approval: ANSWER-ALLOW corr=host:mode-switch#1 tool=permission.mode route=native decision=allow
```

- 这是**真 gate、真队列、真原生允许**（`gate.Native().Allow(ctx, corr, grant)`，grant 取自 `Prompt.Grant`）：
  ⛔ 由 fixture 代替。用例另外钉了 `p.Tool=permission.mode`、`p.Level=L2`、
  `p.CorrelationID` **以 `host:mode-switch#` 开头**（⇒ 不是队列自铸的 `approval-N`）。
- 反向格 `…NoClickMeansNoAnsweredCard`：同一腿、无人点、ctx 到期 ⇒ 拒、盘上仍 `ask_every_step`、
  行里 `asked=0 card=""`、**且**没有 `result=card-not-answered`（点了卡但超时 ≠ 没人答过）、也没有任何
  `ANSWER-ALLOW` 行。⇒ `asked=1` 不是常量。
- 适配器格 `…ThePanelAdaptersPreserveTheAnsweredCard`：`panelModeConfirm` / `storeConfirmFromPanel`
  两个方向都**原样带回**凭据（含 `corr`），且 `nil` 仍是 `nil`（否则两层门同时被关掉）。

---

## ④ `AC#2'-b` 读数（`internal/panel/composer_gate_agreement_371_test.go`）

### ④.1 两处一致的矩阵（真 `perm.Store` 在真 handler 后面，⛔ 假判决）

`6` 个有序对 × `3` 种腿（无通道／非 `nil` 但不问／答过 allow）＝ **18 行**；每行同时读两份**独立记录**：
handler 的返回、store 自己的 `MODE-SWITCH` 行。

| 行 | 面板 | store | 判据 |
|---|---|---|---|
| 无通道 ＋ 目标吃卡 | **自己拒**（`ErrNoL2Confirm`），`writeCalls=0`，腿**没被进过** | 未运行（无行） | 一行面板记录、带 `requestId` |
| 非 `nil` 但不问 ＋ 目标吃卡 | 放行给 store | **拒** `result=card-not-answered asked=0` | 面板也记一行，`result=card-not-answered`（**同一词面**）、带 `requestId` |
| 答过 allow ＋ 目标吃卡 | 放行 | 应用，`asked=1 card="t371/panel"` | `writeCalls=1`、`asks=1`、面板无第二判决 |
| 任何腿 ＋ 目标不吃卡 | 放行 | 应用，`asked=0` | 腿**一次都没被进过**（M4 未扩） |

计数断言：`pairs=6 / rows=18 / 面板拒=2 / store 拒=2 / 应用=14`（枚举被截短即红）。

### ④.2 源码钉：门调的是**同一个**谓词

`Test371R1PThePanelGateCallsTheSameCardPredicate` 解析 `composer_handlers.go`：

- `perm.JudgeCard(` 必须**恰好 1 次**（0 次＝面板自己拼规则；2 次＝两处判）；
- 第二实参必须是逐字 `h.Confirm != nil`（通道项）；
- 全文**不得**再出现 `perm.RequiresCard(to) &&`（旧拼法）；
- 文件必须用 `perm.RefusalResult(err)` 分类（同一 result 词面）。

★**为什么必须有一枚源码钉**：§⑤.3 实测——把这一行改回旧拼法，**行为矩阵仍然全绿**（两种拼法今天等价），
所以“只改 perm、不改 panel”这件事**只有**钉能看见。

### ④.3 契约格

`Test371R1PModeConfirmCarriesTheCredential`：`ModeConfirm` 的返回类型就是 `perm.Card`
（`var leg ModeConfirm = func(…) (perm.Card, error)` 在**编译期**成立）；零值必须读“没答过”。
若有人把 panel 的腿改回 `error`，本文件**编译不过**——这是跨包契约最强的红，且是**有意的**。

---

## ⑤ 反向牙（三枚，全部实测；时机在**最后一发变异之后**）

### ⑤.1 摘掉 `Set` 里“答过”那一跳（`-overlay`）

变异：`if !card.Answered() {` → `if false && !card.Answered() {`（`internal/perm/store.go`）。

```
--- FAIL: Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied
    371A1P|AC0p-2 AE->AA fake limb (the hole)|…|invoked=1|asked=0|set_err=<nil>|last_result=applied|persisted=auto_approve|saves=1
    the fake limb was silently accepted; ticket 371's gate is gone …
--- FAIL: Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt
    the switch applied although no card was ever answered: ticket 371's hole is open
--- FAIL: Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs
    the NEVER-ANSWERING limb was not refused as unanswered: got {refused:false result:applied …asked:0}
```

⇒ **正是**票面要求的红形：`applied=1`、静默。**恢复**：见 §⑥ 的 sha256 核（overlay 不落盘）。

### ⑤.2 摘掉 panel 的“同判/记账”那一跳（`-overlay`）

变异：`composer_handlers.go` 删掉 `h.recordStoreRefusal(req, err, …)`（把面板记回旧形）。

```
--- FAIL: TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371
    ask_every_step -> auto_approve / never-asking: the panel wrote 0 refusal line(s), want 1: ""
```

### ⑤.3 只改 perm、panel 的门留在旧拼法（**落盘变异**，因为源码钉读盘）

变异（落盘，随后恢复）：`if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel {` →
`if perm.RequiresCard(to) && h.Confirm == nil {`。

```
--- FAIL: Test371R1PThePanelGateCallsTheSameCardPredicate
    composer_handlers.go calls perm.JudgeCard 0 time(s), want exactly 1: the panel's half of the mode gate
    must be the SAME judgement perm.Store.Set applies (ticket 371-b), or the two gates are two rules
rc=1
```

★同时（`-overlay` 同形变异）**行为矩阵 `ok`**：这枚钉是**承重的**，不是装饰——两种拼法今天等价，
行为面永远看不见分叉。

**恢复核**：

```
d8eb3c3b7fae4efa084257a10d7b692c814be29eb716ec3a324929b0661dd1c4 *internal/panel/composer_handlers.go
d8eb3c3b7fae4efa084257a10d7b692c814be29eb716ec3a324929b0661dd1c4 *.scratch/371r1p/mut/composer_handlers.go.keep
RESTORED ok
```

---

## ⑥ 还原核（逐文件 sha256 对 `git show HEAD:`，在最后一发变异之后跑）

`DIFFERS` = 本腿有意改动；`SAME` = 未动；`NEW` = 新文件。

| 文件 | 工作树 | `HEAD:` | 判 |
|---|---|---|---|
| `internal/perm/store.go` | `059d4df6…` | `5d084db4…` | DIFFERS |
| `internal/perm/direction.go` | `08663161…` | `2a5ca71b…` | DIFFERS |
| `internal/panel/composer_handlers.go` | `d8eb3c3b…` | `70ec981d…` | DIFFERS（＝⑤.3 恢复后的值，与 save copy 同） |
| `cmd/wisp/run.go` | `584d61dc…` | `05de8110…` | DIFFERS |
| `cmd/wisp/panel_inbound.go` | `942ce7c8…` | `6135b9f4…` | DIFFERS |
| `cmd/wisp/resident_approval_windows.go` | `175cd510…` | `4412de3c…` | DIFFERS |
| `internal/perm/store_test.go` | `0f278256…` | `2d0b91c8…` | DIFFERS（夹具） |
| `internal/perm/ticket90_persist_test.go` | `1f9ce8f4…` | `36c8ebe9…` | DIFFERS（夹具） |
| `internal/perm/store_371_a1p_test.go` | `84badb13…` | `7815883b…` | DIFFERS（翻面，见 ②.4） |
| `internal/tools/booked_mode_377_test.go` | `66a1dbb8…` | `48ea5fcc…` | DIFFERS（夹具） |
| `cmd/wisp/run_mode101_test.go` | `81053c1d…` | `f798be26…` | DIFFERS（夹具；**产品断言逐字未动**） |
| `cmd/wisp/permission_mode_effective_377_test.go` | `74dddda2…` | `fcc69adc…` | DIFFERS（夹具） |
| `cmd/wisp/panel_mode_confirm_114_windows_test.go` | `2eecdb12…` | `4e7740a7…` | DIFFERS（夹具） |
| `internal/panel/composer_handlers_test.go` | `67e71f31…` | `f4c4df72…` | DIFFERS（夹具） |
| `internal/perm/store_371_r1p_test.go` | `db99b286…` | — | NEW |
| `internal/panel/composer_gate_agreement_371_test.go` | `bd3728bf…` | — | NEW |
| `cmd/wisp/perm_371_r1p_test.go` | `ffa36ac4…` | — | NEW |

★**直接调 `Set` 的断言一条没弱化**：`store_test.go`／`store_367_test.go`／`store_372_test.go`／
`ticket90_persist_test.go` 里 12 枚直调格只改**夹具的腿值**（`func(…) error { return nil }` →
`func(…) (Card, error) { return AnsweredCard(…), nil }`）；`store_367_test.go:86` 那枚“故意用字面
`p.to == risk.ModeAutoApprove`”的格子逐字未动。

---

## ⑦ 门禁（rc 逐条，⛔ 相减）

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `go build ./internal/perm ./internal/panel ./cmd/wisp` | **0** | 静默 |
| 2 | `go test ./internal/perm -count=1 -timeout 200s` | **0** | `ok 0.31s`（含 371 的 10 枚） |
| 3 | `go test ./internal/panel -count=1 -timeout 200s` | **1** | 整包 11.4s；**只有 2 枚与本腿无关的红**（见 ⑦.1） |
| 4 | `go test ./internal/panel -run 'Test371R1P\|TestPanelGateAndStoreGate\|TestOnlyTheAutoApproveTarget\|TestAWideningRequest\|TestAStricterOrEqual\|TestAModeRequestIsRefused\|TestTheWriteLegsOwnRefusal\|TestTheModeLadder'` | **0** | `ok 0.05s` |
| 5 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -run 'TestTicket101\|TestTicket367\|TestTicket371\|Mode' -count=1 -timeout 300s` | **1** | 139s；**只有 2 枚 `WISP_SPIKE_MODELS` 环境红**（见 ⑦.2） |
| 5b | 同上，`-run 'TestTicket101\|TestTicket114r2\|TestTicket367\|Test371'` | **0** | 135.6s `ok` |
| 6 | `gofumpt.exe -l`（17 枚改动/新增 .go） | **0** | **输出空** |
| 7 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（408 产线 .go；ban #8 含 `_test.go`） |
| 8 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（18622 条 tracked path，65 条全在名册） |

### ⑦.1 整包 panel 的两枚红：**与本腿输入无关（已证）**

```
--- FAIL: TestComposerContractTypesMatchFrontend      (Go ComposerState emits [files] that interface ComposerState does not declare)
--- FAIL: TestC21DesignTokensFourWayAgree             (Palette.GlassRing/GlassHi/GlassHiSoft 在 internal/ball/tokens.go 里没有颜色字面量)
```

两枚的输入逐字未改，现跑核对：

```
git diff --quiet HEAD -- internal/panel/composer.go frontend/src/lib/panel.ts internal/ball/tokens.go \
    docs/evidence/s1/c21-native-tokens.md frontend/src/styles/tokens.generated.css frontend/scripts/gen-tokens.mjs
inputs_unchanged_rc=0
PROVEN: every input of the two red panel cases is byte-identical to HEAD
```

⇒ 这两枚在 `8f03abd8b` 上就是红的（`frontend/**`／`internal/ball/tokens.go` 归机主外部 agent 的写面），
⛔ 本腿的读数，也⛔ 本腿去修。

### ⑦.2 `cmd/wisp` 的两枚红：**环境变量缺失**

```
--- FAIL: Test15eTheProductionHandOffIsTheModelStoreWalk          install the streaming ASR …: WISP_SPIKE_MODELS is not set
--- FAIL: Test16r1PremortemModelSha256CorruptRefusesTheHandOff    install the streaming ASR …: WISP_SPIKE_MODELS is not set
```

两枚是**语音模型**用例，只因 `-run '…Mode'` 里的 `Mode` 命中 `Model` 才被选中；缺 `WISP_SPIKE_MODELS`
即红，与 `perm`／panel 门无关（5b 用同包、同装配的 101／114r2／367／371 全绿即证）。

---

## ⑧ 边界（逐条对票面 `AC#5'`）

| 相邻票 | 它做了什么 | 本腿做了什么 | 重叠？ |
|---|---|---|---|
| **372**（读侧降级） | `nil` 通道 ⇒ 有效档塌到 `ask_every_step` | 未动：`PermissionMode()` 只读 `s.confirm == nil`，`ResultDowngraded` 逐字未改 | ⛔：372 治“不可问”，本腿治“可问而没问” |
| **375**（动作层 fail-closed） | “不可问 ⇒ 拒那一发” | 未动动作层；本腿在**档位切换**这一发上要求“答过” | ⛔：两层不互相替代 |
| **377**（记账面报有效档） | 报**有效档** | `[bing-b]` 的读数落在**同一条 `MODE-SWITCH`** 行（`asked=`／`card=`），⛔ 另造记账面 | ⛔：同一行加字段，⛔ 新账 |
| **114-r2**（面板腿） | 两条产线门同一条腿 | 两台适配器**保住凭据**（④.3／③.2 末格），门仍是同一个 `JudgeCard` | 复用，⛔ 重造 |
| **M4／R20**（机主 2026-09-21） | 只有切进 `auto_approve` 吃卡 | **未改产品行为**：`run_mode101_test.go:238/256-257` 逐字未动，`TestTicket101ManualSwitchSurvivesRestart` PASS；`AE→AH` 仍 0 卡（②.3） | ⛔ |

---

## ⑨ 剩余（本腿**没有**做，也⛔ 该由本腿做）

1. **R9／〔甲〕**：让**所有变宽**都出卡 ＝ 扩 `M4` ＝ **须摆机主**（机主 2026-09-21 定案；`A1210①`）。本腿未碰。
2. **凭据的同进程边界**（①.1 末段）：挡事故、⛔ 挡恶意；跨进程/跨用户证据不在本票射程。
3. **corr 撞车**：`modeSwitchCardCorr` 用进程计数器 ⇒ 唯一；若将来有人把 id 改成可外部指定，
   `queue.push` 的 clash 分支会追加后缀，`card=` 与 `ANSWER-ALLOW corr=` 的 join 会断（`asked=1` 仍在）。
4. **`store_371_a1p_test.go` 的两处改名**：a1p 证据文件里仍用旧名引述**历史读数**（保留不改，见 ②.4）；
   新名在 `-run 'Test371A1P'` 下仍被选中。★**复跑提示**：a1p 文件 §`⑤`／§`⑦` 里那两条
   `-run 'Test371A1PAC0PrimeFakeLimbIsSilentlyApplied'` / `…TheAutoApproveGateHasTheSameHole`
   是 `786a13b22` 上的**旧名**，在本提交之后**匹配 0 枚**（`go test` 会 `ok [no tests to run]`，
   ⛔ 当成本腿的绿）；要复跑同样的两格，用
   `go test ./internal/perm -run 'Test371A1PAC0PrimeFakeLimbIsRefusedNotSilentlyApplied|Test371A1PAC0PrimeTheAutoApproveGateSplitsTheTwoLimbs' -count=1 -v`。
