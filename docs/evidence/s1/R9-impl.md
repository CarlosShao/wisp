# R9 落地读数：每一次调松都出确认卡（一处定义、四处同判、cmd/wisp 七格翻转）

- 落地腿：`R9`（实现者；**判定归非实现者**，本件只给读数与形状）
- 起手锚：`a2b3b80c9`（HEAD，`git log --oneline -1`）
- 授权：台账 `A1541` 第 4 条逐字「**R9 全弹卡（所有变松动作都出确认卡）批准**」（机主 2026-11-05）；派单 `A1542`／wave9
- 前一腿：`R9` STOP（`docs/evidence/s1/R9-stop-cmdwisp-face.md`）——它给的 §4 草案就是本腿照做的形状；本波机主给了**唯一一枚 `cmd/wisp` 写腿**＋`internal/perm`／`internal/panel`／`internal/command` 三面
- 写面（`git status --porcelain` 逐件，见 §7）：`internal/perm/**`、`internal/panel/**`、`internal/command/{roster.go,roster_test.go}`、`cmd/wisp/**`（仅注释＋两枚测试＋一枚新测试件＋`run.go` 一句卡片文案）、本件
- 零 `frontend/**`、零 `design/**`、零 push、零 amend／reset／rebase／stash／checkout

---

## 0 一句话

门从「**目标是不是全自动**」改成「**这一步是不是调松**」，判据仍然**只有一份**（`perm.RequiresCard`／`perm.JudgeCard`），
四处生产调用点（`perm.Store.Set`／`panel.ModeWriteHandler`／`command.Roster.Judge`／`perm` 自己的定义）同批改；
带腿时**每次调松一张卡**、无腿时**每次调松一律拒**（fail-closed），**调严永远不出卡**。

## 1 判据：一处定义（改前 → 改后）

| 处 | 改前（`a2b3b80c9` blob） | 改后 |
|---|---|---|
| `internal/perm/direction.go:23` | `func RequiresCard(to risk.Mode) bool { return to == risk.ModeAutoApprove }` | `func RequiresCard(from, to risk.Mode) bool { return IsWidening(from, to) }` |
| `internal/perm/direction.go:57` | `func JudgeCard(to risk.Mode, channel bool) CardGate` | `func JudgeCard(from, to risk.Mode, channel bool) CardGate`（三态 `GateNoCard`／`GateNoChannel`／`GateAsk` 不变） |
| `internal/perm/store.go:348` | `switch JudgeCard(to, s.confirm != nil)` | `switch JudgeCard(from, to, s.confirm != nil)`；`from` 仍 `= s.configuredMode()`（372 的记账口径未动） |
| `internal/panel/composer_handlers.go:176` | `if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel` | `if perm.JudgeCard(from, to, h.Confirm != nil) == perm.GateNoChannel` |
| `internal/command/roster.go:287` | `switch perm.JudgeCard(to, env.CardChannel)` | `switch perm.JudgeCard(env.Mode, to, env.CardChannel)`（**第三处生产调用点**，上一腿 STOP 点名的越界面） |

**没有第二套拼法**（现跑）：

```
$ git grep -n "JudgeCard(\|RequiresCard(" -- '*.go' ':!*_test.go'
internal/command/roster.go:294:	switch perm.JudgeCard(env.Mode, to, env.CardChannel) {
internal/panel/composer_handlers.go:184:	if perm.JudgeCard(from, to, h.Confirm != nil) == perm.GateNoChannel {
internal/perm/direction.go:23:func RequiresCard(from, to risk.Mode) bool { return IsWidening(from, to) }
internal/perm/direction.go:75:	if !RequiresCard(from, to) {
internal/perm/store.go:373:	switch JudgeCard(from, to, s.confirm != nil) {
```

### 1.1 一处**新增**的接口，以及为什么非加不可（给验收者：这是本腿唯一一处设计决定）

`ModeWriter` 加了一枚 `ConfiguredMode()`（`perm.Store` 侧由 `ConfiguredMode()` 导出，内部仍 `configuredMode()`）。

理由（**不是**为了让测试变绿，是因为不加就有一条真·假拒）：

- `perm.Store.Set` 的门拿 **`configuredMode()`**（落盘值）当 `from`——372 的记账牙 `TestSetAccountsOnTheConfiguredValueNotTheDowngradedRead` 要求它必须这样；
- 而 `panel.ModeWriteHandler` 原来拿的是 **`PermissionMode()`**（有效档）。**无腿的机器**上有效档一律塌到 `ask_every_step`（372-r2）。
  两者在「配置 `auto_approve` ＋ 无腿 ＋ 请求 `ask_high_risk`」这一格**不是同一对数**：
  - 面板按有效档看是 `AE→AH`＝调松 ⇒ 拒；
  - 写入腿按落盘值看是 `AA→AH`＝调严 ⇒ 本来会放行。
  ⇒ 用户「**把文件调严**」会被面板拒掉，直接违反 `PLAN.md:132`「调严永不设卡」。
- 若不加这枚接口，唯一让门禁变绿的办法是**把 367／371 的「两层同判」矩阵放宽成「面板可以更严」**——那正是本仓禁止的「为变绿放宽断言」。
- 加完之后两层判**同一对数**，`composer_gate_agreement_367/371` 两条矩阵**一格没放宽**（读数见 §3）。

同批的基线口径同步：`command.Env.Mode` 的文档改成「**门将来判的那一端**＝落盘值」，`roster_test.go:314` 从 `store.PermissionMode()` 改成 `store.ConfiguredMode()`
（不改则 roster 会对用户说反话：矩阵里 3 格红，见 §2 第 3 条）。`Env` 的字段、`judgeMode` 的三条判决文案结构未变。

## 2 RED-before（先写新牙、在旧码上跑，逐字）

新牙两件在**不动产码**的 `a2b3b80c9` 上先落地：

- `internal/perm/store_r9_test.go`（3 格）
- `cmd/wisp/r9_mode_card_test.go`（2 格：run 缝 1＋CLI 管 1）

```text
$ export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
$ go test ./internal/perm -run "TestR9" -count=1
--- FAIL: TestR9EveryLooseningCostsOneAnsweredCard (0.00s)
    store_r9_test.go:61: ask_every_step -> ask_high_risk cost 0 L2 card(s), want exactly 1 (R9: every loosening shows a card)
    store_r9_test.go:65: ask_every_step -> ask_high_risk: result="applied" asked=false card="", want applied/true/<id> ...
--- FAIL: TestR9ANeverAskingLimbCannotLoosenTheMiddleTier (0.00s)
    store_r9_test.go:117: ask_every_step -> ask_high_risk applied although no card was ever answered: R9's middle tier inherited ticket 371's hole
FAIL	github.com/CarlosShao/wisp/internal/perm	0.067s

$ go test ./cmd/wisp -run "TestR9" -count=1 -timeout 420s
--- FAIL: TestR9TheRunSeamCardsEveryLooseningAndNoTightening (2.88s)
    r9_mode_card_test.go:102: L2 cards raised = [auto_approve], want [ask_high_risk auto_approve] ...
    r9_mode_card_test.go:121: ask_every_step -> ask_high_risk row does not read asked=1/applied: [audit] perm: MODE-SWITCH from=ask_every_step to=ask_high_risk ... asked=0 card="" result=applied ...
--- FAIL: TestR9TheCLIPipeRefusesALooseningItCannotCard (0.04s)
    --- FAIL: .../a_loosening_is_refused_and_nothing_lands (0.02s)
        r9_mode_card_test.go:148: ... must be refused with ErrNoL2Confirm ...; got <nil>
        r9_mode_card_test.go:152: the refused loosening still moved the file: reads "ask_high_risk", want ask_every_step
FAIL	github.com/CarlosShao/wisp/cmd/wisp	3.029s
```

⇒ 四条读数就是票面那句「每一次调松都要一张卡」在旧码上**不成立**的原样：中档调松 0 卡、无卡面的管子照样调松。

**第三处红（不在新牙里，是矩阵在旧基线上自己报的）**：把产码改完后、`command/roster_test.go` 还没换基线时，18 行矩阵当场 3 格红，
逐字点名 `ask_high_risk-to-ask_high_risk-card-false`／`auto_approve-to-ask_high_risk-card-false`／`auto_approve-to-auto_approve-card-false`
「roster says Refused ... but the door accepted the switch」——那正是 §1.1 那条「roster 与门不是同一对数」，不是本腿猜的。

## 3 GREEN-after（同一批牙，改后）

```
$ go test ./internal/perm -count=1                       ok  github.com/CarlosShao/wisp/internal/perm   0.426s
$ go test ./internal/command -count=1                    ok  github.com/CarlosShao/wisp/internal/command 0.260s
$ go test ./cmd/wisp -run "TestTicket101ManualSwitchSurvivesRestart|TestTicket114r2|TestR9" -count=1 -timeout 420s
ok  github.com/CarlosShao/wisp/cmd/wisp   45.899s
```

矩阵口径（`perm/store_367_test.go`，三档 6 对，逐对枚举并在尾部断言 3/3）：

| 通道 | 3 对调松（AE→AH／AE→AA／AH→AA） | 3 对调严（AH→AE／AA→AE／AA→AH） |
|---|---|---|
| 无确认腿 | **拒**、`result=refused-confirm`、文件不动、审计一行、拒绝句点名两端 | **放行**、`result=applied`、文件落盘 |
| 有真腿（应答） | **一张卡**、`asked=1`、`card=<id>`、`applied=1` | **0 张**、腿一次都不进、`applied=1` |
| 有腿但从不问（假腿） | `ResultCardNotAnswered`（371 的牙，**现在覆盖中档**）、文件不动 | 0 张、放行 |

## 4 翻转的格（哪些断言改了、为什么不是放宽）

| 文件 | 格 | 旧期望 → 新期望 | 为什么不是放宽 |
|---|---|---|---|
| `internal/perm/store_367_test.go` | 首表 | `refused=2/applied=4`（只判目标）→ `refused=3/applied=3`（判方向） | 同一张六对表，**多拒的一格正是中档调松**；调严三格逐字仍在，表尾 3/3 断言保留 |
| 同上 | 次表 | 「中档不出卡」→「中档 1 张、再调松第 2 张、调严 0 张」 | 多出的两格是**新增牙齿**；`spy.calls[i]` 的 from/to 逐格钉住 |
| `internal/perm/store_371_a1p_test.go:97` | AC#0② | `invoked=0/asked=0` → `invoked=1/asked=1` | 这是**历史读数格**（为票 371 量的），R9 改的就是它；limb 侧事实（进入次数）反而更严 |
| `internal/perm/store_371_r1p_test.go:233` | 中档免卡 | 名字与期望一起翻成「中档 1 张」（＋调严 0 张的控制臂） | 旧名 `...StillFreeOfCards` 留着就是谎；新格把两个方向都钉住 |
| `internal/perm/store_test.go:135` | `TestTicket90EveryLooseningCostsAConfirm...`（改名） | `AH` 0 张 / `AA` 1 张 → `AH` 1 张 / `AA` 2 张 | 审计三档全写的断言**逐字未动**；新增 `spy.calls[0]`／`calls[1]` 两格 |
| `internal/perm/ticket90_persist_test.go:177` | 重启持久 | `confirmCalls=1` → `2` | 两次切换都是调松，读数随规则 |
| `internal/panel/composer_handlers_test.go:134` | 无腿矩阵（改名为 `TestEveryLooseningIsRefused...`） | 拒 2 → 拒 3；放行 4 → 3；尾部三条计数公式同步 | 9 格逐格枚举不变，`unchanged=3` 断言不变；多拒的是中档调松 |
| `internal/panel/composer_gate_agreement_367_test.go` | 两门同判 | `RequiresCard(to)` → `RequiresCard(from,to)`；`refused 2→3` | **矩阵本身一格没放宽**（仍是「两层必须同时拒／同时放」） |
| `internal/panel/composer_gate_agreement_371_test.go` | 18 行矩阵 | `refusedByPanel/Store 2/2` → `3/3`；源码钉 `len(args)==2` → `3`，并**新增** `Args[0]=="from"`、`Args[1]=="to"` 两枚 | 18 行（6 对 × 3 种腿）仍是全枚举；`result=card-not-answered` 两处同词断言未动；源码钉**更严** |
| `internal/command/roster_test.go:314` | 18 行对撞 | `Env.Mode` 取 `PermissionMode()` → `ConfiguredMode()` | 判据/断言结构未动，只把基线换成门实际判的那一端（§1.1） |
| `cmd/wisp/run_mode101_test.go:256` | 重启缝 | `want [auto_approve]` → `[ask_high_risk, auto_approve]` | 两次切换都是调松；`MODE-SWITCH from=ask_high_risk to=auto_approve` 与 `result=applied` 逐字未动 |
| `cmd/wisp/panel_mode_confirm_114_windows_test.go` | 六格（见 §4.1） | 见下 | 见下 |

### 4.1 114 那六格，逐格

| 格 | 旧（372 口径） | 新（R9） |
|---|---|---|
| ① 无腿 + `AE→AH`（resident 门） | 放行、盘上 `ask_high_risk` | **`ErrNoL2Confirm` 拒**、盘上仍 `ask_every_step` |
| ② 有腿 + `AE→AH` | 放行、**0 张卡** | 放行、**1 张卡** |
| ③ 同门再 `AH→AA` | 1 张（`[auto_approve]`） | 2 张（`[ask_high_risk, auto_approve]`） |
| ④ CLI 管 `ask_high_risk` | 落盘、有效档仍 `ask_every_step` | **拒**（`ErrNoL2Confirm`）；同格改量「管子只能调严」＋有效档仍 `ask_every_step` |
| ⑤ CLI 管 + 假腿 `AE→AH`（改名 `...StillOwesAnAnsweredCard`） | 落盘（那是当时的洞） | **拒**、腿被问 1 次、文件不动 |
| ⑥ run 缝 handler 两次请求 | 1 张（`[auto_approve]`） | 2 张（`[ask_high_risk, auto_approve]`） |

④ 那格原来的**两条**读数（拒绝 ＋ 有效档降级）都保住了：拒绝改成对调松量，降级改成用一份手写的 `ask_high_risk` 配置让管子把它调严到 `ask_every_step`、再用 `perm.New` 冷读有效档。**没有删掉任何一条读数。**

## 5 门禁读数（本腿，现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（无输出） |
| test（面内） | `go test ./internal/perm ./internal/command` | `perm` **ok 0.426s** · `command` **ok 0.260s** |
| test（panel） | `go test ./internal/panel` | 本腿射程全绿；仍**2 枚既有红** `TestComposerContractTypesMatchFrontend`（前端 `ComposerState.files`）、`TestC21DesignTokensFourWayAgree`（C21 令牌表）——与 `R9-stop-cmdwisp-face.md` §5 同读，HEAD 干净树同样红，非本腿 |
| test（cmd/wisp 门） | `go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart\|TestTicket114r2\|TestR9' -count=1 -timeout 420s` | **ok 45.899s** |
| gofumpt | `gofumpt -l cmd internal tools` | 仅 `internal/agent/approval/pending_read.go`（**既有** CRLF 件，`git status` 干净、本腿零触碰；A1540/A1542 已具名）。本腿改过的每个文件 **0 命中** |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**（读数见 §5.1；本腿第一次跑被 ban#9 抓到一处 phantom citation，正是本件缺席所致，补上后转绿） |

## 6 反转变异（红，逐字）

| 变异 | 改法 | 结果 |
|---|---|---|
| `M1` 谓词退回目标 | `direction.go` 的 `RequiresCard` 改回 `return to == risk.ModeAutoApprove` | `perm` 的 R9/367/371 六格红（中档 0 卡）；`cmd/wisp` 的 `TestR9...` 两格红（`[auto_approve]` / 管子又放行） |
| `M2` 只改一处 | `store.go` 传 `to`（面板/roster 不动） | 两门同判矩阵红：面板拒、store 放 |
| `M3` 去掉 channel 项 | `JudgeCard` 里 `if !channel` → `if false` | 无腿矩阵红（`refused=0`）、`Store` panic 守卫红 |

## 7 接受的代价（写进票面，具名）

1. **CLI 管（`wisp panel-inbound`）从此不能调松**：它没有卡面 ⇒ 每次调松都 `GateNoChannel` 拒。配置手改通道（D36）仍可调松且**本来就出卡**，所以「所有调松都有卡」在两条通道上都成立，代价是命令行少了一条调松的路。
2. **`AE→AH` 的运行时通道从「0 卡」变「1 卡」**：这是机主裁定的本体；面板/球/CLI 的调松动作从此每次一张 C18 卡。
3. **`PermissionMode()` 的读侧降级（372-r2）不变**：无腿进程仍按 `ask_every_step` 执行，`config.toml` 原值不被改写。

## 8 复核（任意人可重跑）

```sh
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
git grep -n "JudgeCard(\|RequiresCard(" -- '*.go' ':!*_test.go'   # 定义 2 处 + 调用 3 处
go build ./cmd/... ./internal/... ./tools/...
go test ./internal/perm ./internal/command -count=1
go test ./internal/panel -run 'TestEveryLoosening|TestPanelGateAndStoreGateAgree|TestAWideningRequest|TestAStricterOrEqual|TestTheModeLadder|Test371R1P' -count=1
go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart|TestTicket114r2|TestR9' -count=1 -timeout 420s
/d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools
sh scripts/d22scan.sh
```
