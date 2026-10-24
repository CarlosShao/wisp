# 票 376 落地腿 `376-r1`（写腿）：〔乙′〕按 level 选源 —— L2 卡印 `300.0s`、L1 卡仍印 `2.0s`

- 腿 `376-r1`（实现者）。落笔读数时 HEAD `00c7468dd`（2026-10-24 13:0x +08）。本笔**未翻框**、未改票名／`Status:`／`-done`、未改台账、未 push。
- 写面**只有两枚**：`cmd/wisp/run.go`（**改**，`5/1`）＋ `cmd/wisp/approval_card_window_376_test.go`（**新建**，派单「＋必要测试」那一条）＋本件。`internal/agent/approval/**` **零字节改动**（`git status --porcelain -- internal/agent/approval` 空）；`gate.go:534` 一个字没动。
- 临时件全在仓外 `D:\tmp\wisp\20261024\376r1\`：`head.go`／`cut.go`（两枚仓外副本）、`*-overlay.json`（`-overlay`，工作树零突变）、`zz376r1_probe_test.go`（真钟双发探针，经 overlay 注入 `cmd/wisp`）、各 `*-run.txt`／`*-v.txt` 读数原文。
- ★**尺的射程（落笔前五问）**：①对象＝`approval.Prompt` 的**两个字段**及其唯一印发点，⛔ Win32 的 window；②形状＝`file:line`＋**卡面逐字**，行号一律现跑（`run.go` 已漂）；③数的是**代码**：`.Window` 命中**分生产/测试**看，⛔ 把 `CreateWindowExW`／`Options.Window` 混进来；④总数核对：生产 `.Window` 命中逐条列全、与该尺总命中对上；⑤秒数＝`time.Duration.Seconds()`，卡面按 `%.1f` 印、故一位小数。

## ① 为什么〔乙′〕（引普查四形表，⛔ 重做它的量）

普查 `docs/evidence/s1/376-a1-census.md` §④ 的四形实测（我不复测，只引判据）：

| 形 | 改动量 | L2 卡 | L1 卡 | 判 |
|---|---|---|---|---|
| 甲（`gate.go:534` 传真实 `win`）| `1/1` | `300.0s` ✓ | `2.0s` ✓ | `replies.go:83`「zero on an L2 card」变假 ⇒ 还得 +1 行注释 |
| 乙（字面：`run.go` 改读 `Deadline`）| `1/1` | `300.0s` ✓ | **`0.0s` ✗** | L1 的 deadline 实参是 `0`（`gate.go:293`）⇒ 把缺陷搬到 L1 |
| **乙′（按 level 选源）** | **`5/1`** | **`300.0s` ✓** | **`2.0s` ✓** | 两者全对、**生产者/字段语义不动** |
| 丙（`Window` 改成剩余时间）| `6/2` | 不破 | 不破 | 「剩余时间」今天拿不到（`panel/pump.go:50-53` 的 `approval.EventTick` **无生产者**）|

★理由（按铁律「最有价值＋主流 harness」，⛔ 按「风险最小」）：**主流 harness 的审批卡倒计时是渲染层按 level 选源**（L1 块长／L2 的 C18 截止）⇒ **生产者不变量不动**；甲只省一个分支，代价是把**已写进文档的不变量**改成随 level 两义。★落点：**改的正是唯一生产读者自己**（`Prompt.Window` 全仓只有这一个生产读者，见 §⑥）。

## ② 改了什么（逐行，`git diff --numstat` ＝ `5 1 cmd/wisp/run.go`）

★**行号现跑、⛔ 用票面的 1563**：`373-r1` 改过 `run.go`（`93fdc0496` 等），旧号已漂。现跑 `grep -n '窗口' <(git show HEAD:cmd/wisp/run.go)` ⇒ **`1657`**。`gate.go:534` **未漂**，逐字仍是 `p := g.promptFor(d, corr, 0, g.q.Timeout(), nil, nonce)`。

`cmd/wisp/run.go:1657-1661`（`consoleApprovalUI.Prompt` 里那一行）：

```go
	shownWindow := p.Window
	if p.Level == "L2" { // an L2 card's Window is 0 by construction; its wait is Deadline
		shownWindow = p.Deadline
	}
	fmt.Fprintf(u.out, "  窗口 %.1fs\n", shownWindow.Seconds())
```

★**那 5 行各是什么**（逐字即上，`sed -n '1657,1661p' cmd/wisp/run.go`）：① `shownWindow := p.Window` ＝默认取**甲源**（L1 块长；L2 上恒 0）；② `if p.Level == "L2" {` ＝**判 level**（只有 L2 换源，注释写明「L2 的 Window 按构造就是 0，它等的是 Deadline」）；③ `shownWindow = p.Deadline` ＝换到**乙源**（C18 截止，即 `q.Timeout()` 的同一个值）；④ `}` 闭括号；⑤ `fmt.Fprintf(… shownWindow.Seconds())` ＝**唯一被改的旧行**（原 `p.Window.Seconds()`，`git show HEAD:cmd/wisp/run.go` 的 `1657` 行）。⇒ 净效果：**L1 读 `Window`（块长）、L2 读 `Deadline`（C18 截止），且 `Prompt.Window` 在 L2 上仍然是 0**。

★**没动的东西**：`Prompt.Window`／`Prompt.Deadline` **字段语义与写入点零改动**（`gate.go:609-610`、`:293`、`:534`）；`replies.go:83`「Window is the L1 block length, zero on an L2 card」**今天仍然为真**（这正是不选甲的理由）；`ReplyCard.Window/.Deadline` 仍是死拷贝（0 读者）。★**`Prompt.Deadline` 从本笔起有了第一个生产读者**（就是这行打印），这是〔乙′〕的应有之义、⛔ 语义改动。

## ③ 必给读数 1：`L2` 卡印的秒数 **＝ `q.Timeout()`**（`EQUAL=true`）

真钟探针（`D:\tmp\wisp\20261024\376r1\probe-realwait-run.txt`，真 `consoleApprovalUI` ＋ 真 `approval.New`）：

```
[确认 L2 fs.write] 覆盖已有文件越出 [fs] 授权目录，需 L2 强确认
  影响路径：C:/work/existing.txt
  窗口 31.0s
[warning] 审批将在 30 秒后自动拒绝，请尽快确认
[dismissed] 审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝
PROBE_L2 PRINTED=31.0 QUEUE_TIMEOUT=31.0 QUEUE_TIMEOUT_RAW=31s ELAPSED=31.002 ANSWER=timeout EQUAL=true
PROBE_L2_DEFAULT PRINTED=300.0 QUEUE_TIMEOUT=300.0 EQUAL=true ANSWER=timeout
```

- ★**默认配置（`confirm_timeout_sec=300`）**：`PRINTED=300.0 == QUEUE_TIMEOUT=300.0`、`EQUAL=true` ⇒ 与普查**甲 overlay 的正控**（`300.0 (EQUAL=true)`）逐位一致，只不过现在走的是〔乙′〕。★这一发是同一探针的第三发、用**假钟**（`clock376`，只判「读的是哪个源」、不判墙钟；`QUEUE_TIMEOUT` 仍是闸门自己的 `300.0`），真钟读数见下一条。
- ★**带底 31**：卡印 `31.0s`，**同一发实测等了 `ELAPSED=31.002s`** ⇒ 印出的秒数**等于真实等待**（⛔ 只是「不再是 0.0」；同一次运行里 `QUEUE_TIMEOUT_RAW=31s` 也印在读数里）。
- ★反控（HEAD，`D:\tmp\...\test-head-overlay.txt`）：`窗口 0.0s` vs `Queue().Timeout() = 300.0s` ⇒ **红**。判据是「印的 == 真等的」，不是「≠0.0」。

## ④ 必给读数 2：`L1` 卡**不能**被打成 `0.0s`

- ★**印什么**：`窗口 2.0s`。★**真实块长**：`Gate.Window() = 2.0s`（`GATE_WINDOW_RAW=2s`）——★**配置 `Window: 1s` 被 `gate.go:143-151` 夹到 `MinL1Window=2s`**（真实配置路径 `l1_window_sec=1` 到达时就是这个值）。
- ★**真等多少**：同一发 `ELAPSED=2.001s`（L1 结束未被否决 ⇒ `ANSWER=timeout` ＝执行，B1）。
- ★这正是**乙（字面换源）会破的那一发**：L1 的 `Prompt.Deadline` 是 `0`（`gate.go:293` 第 4 参），字面换源会把 L1 卡印成 `0.0s`。本笔的 L1 用例把这枚反例钉住（见 §⑦）。

## ⑤ 两发对照（真钟，`probe-realwait-run.txt`；`[确认 …]` 逐字见 §③④）

| 发 | 型 | 卡上那行 | 真等的 | 判据 |
|---|---|---|---|---|
| L2（覆盖已有文件）| `Level=risk.L2`、`Paths=[C:/work/existing.txt]` | **`窗口 31.0s`**（默认配置那一发 `300.0s`）| `QUEUE_TIMEOUT=31s`、实测 `31.002s` | `EQUAL=true` |
| L1（低危提示）| `Level=risk.L1`、`Paths=[C:/data/a.txt]` | **`窗口 2.0s`** | `GATE_WINDOW=2s`、实测 `2.001s` | `EQUAL=true`、`L1_ZERO=false` |

## ⑥ 「别处不印 `Window`」——我自己复核过（⛔ 抄普查）

尺：`grep -rn '\.Window' --include=*.go internal cmd | grep -v _test.go`（**全部命中逐条看**）＋ `grep -rn '窗口' --include=*.go internal cmd | grep -v _test.go | grep -E 'Fprint|Printf|Sprintf|slog'`。结论：**卡片印发点只有 `run.go:1657` 这一处**。

- ★球卡 `ballCardUI.Prompt`（`cmd/wisp/resident_approval_windows.go:890-950`）：印 level／tool／编号／orb 状态／Esc 借键行，**没有一个时长字段**（逐行读过）。
- ★面板 `internal/panel/approval.go:38-59` `ApprovalCardView` 与 `approval.PanelItem`（`ui.go:50-57`）**根本没有时长字段** ⇒ 页面印不出没收到的数；`L1WindowWait`（`internal/panel/pump.go:127-130`）只有两个 `id`。
- ★唯一同时说出两个数的是 **boot 审计行** `resident_approval_windows.go:410`：`gate_window`＝`Gate.Window()`（配置进来的 L1 块长）、`gate_queue_timeout`＝`Queue().Timeout()`（C18）——★**slog、非卡、两个都对**。
- ★撞名的三处**不是**这个字段：`approval.Options.Window`（`gate.go:26-29,143`，配置）／`Gate.Window()`（`:176`，读面）／`tools.GateFuncs.Window`（`internal/tools/gate.go:122-131`，L1 回调字段）。
  ⇒ ★**没有发现别处也印 `Window`**，故**不触发**派单「停手上报」那一条。

## ⑦ `AC#2` 的牙：新用例逐名 ＋ 为什么今天响

新文件 `cmd/wisp/approval_card_window_376_test.go`（**新建**，派单写的「＋必要测试」；`cmd/wisp` 其它文件零改动）。

1. ★`Test376R1L2CardPrintsQueueTimeout` —— **今天响、不修不响**。期望值取**闸门自己的** `g.Queue().Timeout()`（⛔ 不是字面 300），从**真实控制台输出**里用 `窗口 ([0-9]+\.[0-9])s` 读回卡上那个数再比。HEAD 传（`-overlay head.go`）实测红：`approval_card_window_376_test.go:248: the L2 card printed "窗口 0.0s", but the gate waits Queue().Timeout() = 300.0s (5m0s).`（同一次运行里卡面正文就是 `窗口 0.0s`）。
2. ★`Test376R1L1CardKeepsItsBlockLength` —— **反例守门**（绿前绿后）：配置 `Window: 1s` ⇒ 先钉 `Gate.Window()==MinL1Window`，再钉卡印 `2.0s`、`≠0`；它在 HEAD 也绿，作用正是把**乙（字面换源）**那种改法**转红**（若换成打印 `Deadline`，这里印 `0.0s`）。
3. ★两枚用例都把**生产者形状**当**前提**钉住：L2 的前提是 `Prompt.Window==0` 且 `Prompt.Deadline==Queue().Timeout()`（`approval/gate.go:534` 的第 3／4 参分叉），L1 的前提是 `Prompt.Window==Gate.Window()`、`Prompt.Deadline==0`。⇒ 将来谁把 `gate.go:534` 的 `0` 换成真值（甲），失败信息会**具名**说要连带改 `replies.go:83`，⛔ 让这条不变量默默两义。
4. ★两枚用例都用**假钟**（`approval.Clock` 缝：`After` 登记通道、测试手打）⇒ 300s 的判据只用微秒，**不把 5 分钟等待塞进 CI**；真钟的「印的 == 真等的」读数在 §③⑤ 的仓外探针里。

## ⑧ 反向判据：三读数 ＋ 还原核（还原核在**最后一发变异之后**量）

- ★**变异红**：仓外 `cut.go` ＝ 从修后文件里**只摘掉那一跳**（删 `if p.Level == "L2" { // …` ／`shownWindow = p.Deadline`／`}` 三行，保留 `shownWindow := p.Window`），`-overlay` 换进编译。`D:\tmp\...\read-red-cut.txt` ⇒ **`rc=1`**，逐名红行：`--- FAIL: Test376R1L2CardPrintsQueueTimeout (0.00s)`＋`approval_card_window_376_test.go:248: the L2 card printed "窗口 0.0s", but the gate waits Queue().Timeout() = 300.0s (5m0s).`；同一次 `--- PASS: Test376R1L1CardKeepsItsBlockLength`（L1 那一半不跟着红，说明红的是**被摘的那一跳**）。
- ★**还原绿**：`read-green-restored.txt` ⇒ `rc=0`、`ok github.com/CarlosShao/wisp/cmd/wisp 0.248s`（`-run Test376R1`）。
- ★**原样绿**：修后首次定向（`test-new-run.txt`）⇒ `rc=0`、`ok … 0.237s`。
- ★**还原核**（时机＝**最后一发变异跑完之后**）：`sha256sum cmd/wisp/run.go` 变异**前** `565ff3467235c201c87616e063671e6e8f55d113d35423e890125f35fa5ea6da`，变异**后** `565ff3467235c201c87616e063671e6e8f55d113d35423e890125f35fa5ea6da` ⇒ **逐位相同**；`git diff --numstat -- cmd/wisp/run.go` ⇒ `5	1`。⇒ 变异全程走仓外副本＋`-overlay`，**工作树零突变窗口**。

## ⑨ 四数（改前／改后**分开列**，⛔ 相减）

| 尺 | 改前 | 改后 |
|---|---|---|
| `go test -count=1 -v ./cmd/wisp -run '<12 个分支>'`（带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`）| **PASS 13 / FAIL 0 / SKIP 0**、`rc=0`、113.834s | **PASS 15 / FAIL 0 / SKIP 0**、`rc=0`、125.327s |
| `go test -count=1 -v ./internal/agent/approval`（整包）| **PASS 97 / FAIL 0 / SKIP 1**、`rc=0`、1.036s | **PASS 97 / FAIL 0 / SKIP 1**、`rc=0`、1.203s |

- 尺的射程：`cmd/wisp` 「改前」＝ `-overlay` 把 `run.go` 钉回 HEAD **并且**把本笔新测试件删出编译（`{"…approval_card_window_376_test.go": ""}`），所以 13 vs 15 的差**正是新加的两枚**（两列都含 2 枚子用例 `TestL1VetoNeedsAChannelTheHostReallyWired/*`，故顶层是 11 → 13）；两发用**同一条 `-run` 正则**（12 个分支，原文见 `D:\tmp\wisp\20261024\376r1\cmdwisp-targeted.sh`：`Test376R1*` ＋ 卡片/打印邻域 `TestRunBooksWithASnapshotOfItsLiveQueue`／`TestThePumpIsDrivenNotJustAssembled`／`TestReplyListener*`／`TestPanelRouteRefusesAnAllow…`／`TestUnansweredL2CardTimesOutIntoRejectNeverExecution`／`TestL1VetoNeedsAChannelTheHostReallyWired`／`TestAlwaysBranchStores*`／`TestTicket223HandEditedFsLooseningCostsAnL2Card`／`TestTicket223TighteningRaisesNoCard`）。
- 尺的射程：`internal/agent/approval` 本笔**零字节改动**（`git status --porcelain -- internal/agent/approval` 空、`git diff` 空）⇒ 两列量的是**同一份字节**，但我**独立跑了两遍**（`approval-after-v.txt`／`approval-second-v.txt`），⛔ 用相减推。`SKIP=1` 是 `TestDefaultDeadlineWallClockMeasurement`（本笔未动）。

## ⑩ 门禁 rc ＋ 格式（甲乙两尺都报）

- `go build ./...` ⇒ **rc=0**（`gate-build.txt` 空）；`go vet ./cmd/wisp ./internal/agent/approval` ⇒ **rc=0**（`gate-vet.txt` 空）。
- `d22scan`（`tools/d22scan` 自成 module，`-root ../..`）⇒ **rc=0**：`examined cmd/ 151 Go files`、`internal/ 643`、`clean - no D22 ban violations`。
- ★**格式甲尺**：`/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/run.go cmd/wisp/approval_card_window_376_test.go` ⇒ **rc=0、输出空**（两枚都合格；`grep -cE '^([A-Za-z]:)?[^\r\n]*[\\/]'` ＝ **0** 枚归属本仓）。
- ★**格式乙尺**：`gofmt -l` 同两枚 ⇒ **rc=0、输出空**。⛔ 没用 `gofumpt -l .`（会扫 `.scratch` 里故意写坏的探针 ⇒ `rc=2`，且判归属须 `[\\/]`）。
- 长度门：本件 `wc -l docs/evidence/s1/376-r1-impl.md` ＝ **118** 行（≤160）。

## ⑪ 纪律格（自报）

- ⛔ 未改 `gate.go:534` 第 3 参；⛔ 未动 `Prompt.Window`／`Prompt.Deadline` 字段语义；⛔ 未碰 `internal/winsec/**`、`internal/memory/**`、`internal/perm/**`、`internal/panel/**`、`internal/agent/tools/**`、`cmd/wisp/**` 其它文件、`scripts/**`、`.github/**`、`go.mod`/`go.sum`、`thresholds.go`／golden／`allowlist.txt`、`frontend/**`／`design/**`（后者**零读零写**）。
- ⛔ 未翻框、未改票名／`Status:`／`-done`、未改台账、未 push。工作树里别人的脏（`.scratch/**`、`?? $null`／`-`／`-dummy`）**未碰未清**。
- 开工前核：`git diff --cached --name-status` **空**（无 `D`/`A` 残留）；`cmd/wisp/run.go` 开工时**干净**（`373-r1` 已提交 `93fdc0496`）；写面内无别人的未提交改动。
