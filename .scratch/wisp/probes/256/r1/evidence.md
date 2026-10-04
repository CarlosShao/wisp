# 票 256 `256-r1` 落地写码腿证据件 —— 常驻那条腿的审批门改吃 `[risk]` 两枚配置（形ⓐ：签名不加参）

**腿**＝`256-r1`（只做编排者已裁的 `[risk]` 那一半，`Grants` 那一半归口待立票 265，本轮不碰）
**任务书**＝`.scratch/wisp/issues/256-resident-gate-built-before-session-grants.md` §8（编排者 10-04 09:5x 裁定，账 `A591`）＋只读普查件 `.scratch/wisp/probes/256/a2/census.md`
**写面**（逐字照任务书）＝`cmd/wisp/resident_approval_windows.go` ＋ `cmd/wisp/resident_windows.go` 的 `:126` 一处调用 ＋ 新增测试件 ＋ 本证据件目录
**仓库**＝`D:\work\workspace\projects plans\Wisp`，分支 `dev`
**起手锚点**＝HEAD `d0847aa1`（`264-a1 先写满尾部两节…`）
**成稿时刻**＝2026-10-04 10:2x +08

---

## §0 这一轮到底把什么换掉了（一句话，带读数）

常驻那条腿（悬浮球＋托盘那一支进程）建审批门时，`approval.Options` 里**实传枚数从 3/10 变成 5/10**：新增的正是
`Window` 与 `ApprovalTimeout` 两枚 ⇒ 改 `confirm_timeout_sec` 现在**对这一发常驻进程建出来的那枚门生效**。
`Grants` 那一枚仍然不传，且判据 ④ 把它**钉成"不许出现"**（见 §5）。
枚数不是注释里数出来的，是从 `go/ast` 的语法树上数的，见 §3④。

---

## §1 起手撞钉预检读数

### 1.1 环境读数（本轮新量，两条具名记着，后来者别再撞一遍）

**① `cmd/wisp` 的测试二进制在本机需要先挂 DLL**——这是票 98 的洞、票 111 的 `scripts/wisp-cli-tests.sh` 存在的原因，
本轮**现场复现两发**：

| 命令形状 | 读数 |
|---|---|
| `go test -count=1 -run '...' ./cmd/wisp`（裸跑） | `exit status 0xc0000135` ＋ `FAIL ... 0.051s`，**0 行测试输出** |
| 直接跑 `go test -c` 出来的 `wisp.test.exe` | `rc=127`，`error while loading shared libraries: sherpa-onnx-c-api.dll` |
| `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ... ./cmd/wisp` | **绿**（起手一发 `rc=0`，25.–30.5s） |

⇒ 本腿所有 `go test ./cmd/wisp` 读数**都是带这枚 PATH 前缀跑的**，逐字命令在 §6。
`go build ./...` 与 `go vet ./cmd/wisp/` **不需要**这枚前缀（它们不启动二进制）。

**② 任务书里那条 `go test ... cmd/wisp` 的包路径形状在 Git Bash 下不成立**：
裸 `cmd/wisp` 会被当 import path 解析 ⇒ `package cmd/wisp is not in std (D:\work\base\go\src\cmd\wisp)` ⇒ `[setup failed]`。
本腿一律改成 `./cmd/wisp`。**这不是放宽，是把命令写对**；票面原句不改，就地记这一处形状差。

**③ 起手时工作树已经有别人的活**（不是本腿造成的，逐字记录以免被当成越界）：
起手第一发 `git status --porcelain` 在 HEAD `d0847aa1` 上就列出了 `D design/assets/*`、`D design/screens/*`、
`M design/doubao/**`、`M .gitignore`，以及后来 `git diff --stat` 里的 `M internal/agent/approval/pending_read.go`、
`M tools/d22scan/{main.go,scan_test.go,selftestsamples.go}`。这些是 §任务书里点名的在飞腿（`260-r1`／`212-r3`）
与编排者的暂存面，**本腿一枚都没碰、一枚都没提交**（尺＝§8 的 `git show --numstat`，逐发枚数＝1／3）。

### 1.2 今天绿的用例名册（起手一发，`-run 'TestAC246|TestTicket224|TestTicket255Roster'`）

原始输出＝`.scratch/wisp/probes/256/r1/start-targeted.txt`。
**终值：top-level PASS 20／FAIL 0／SKIP 0；子用例 PASS 7**，`ok github.com/CarlosShao/wisp/cmd/wisp 26.170s`。

```
TestAC246CancelGestureUsesTheInjectedExecutor
TestAC246CancelStepHookRunsOnTheRealShutdownSequence
TestAC246CardWithNoWindowFailsClosedThroughTheRealGate
TestAC246ChannelNeedsBothWindowAndExecutor
TestAC246DevLegIgnoresTheTestTaskInjection
TestAC246EscChannelStaysUnloadedWithoutABallWindow
TestAC246ResidentGateInjectionIsRefusedHalfAssembled
TestAC246ResidentPipelineAsksThroughTheOneGate
TestAC246ResidentTaskRootCancelStopsTheModelCall
TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry
TestAC246ShippedResidentProcessOwnsItsCancelStep
TestAC246StatusLineSaysWhatTheLegDoesNot
TestAC246TestTaskInjectionPredicate
TestAC246VetoSentenceWithNoCard
TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking
TestTicket224ProductionSessionDoesNotSurviveRestart
TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun
TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim
TestTicket255RosterStillMatchesTheActualReadSites
```

⚠ `TestTicket255Roster*` 两枚**起手绿**；收尾一发见 §6（P7 自锁，任务书要求的两发）。

### 1.3 P1／P3／P4 会不会被本轮顶到（逐枚现读断言，不给"应该没问题"）

**P1 — `TestAC246ResidentPipelineAsksThroughTheOneGate`（`cmd/wisp/resident_task_source_246_windows_test.go:68`，钉在 `:89`）⇒ 不顶到。**
真身逐字：`if ar.gate != ra.gate { t.Fatalf("the assembled pipeline runs on a DIFFERENT gate than the resident leg built: %p vs %p (two gates in one process is the shape ruling 2.2 and ledger A481 refuse)"`。
⇒ 本轮**一个 `approval.Options{}` 字面量、一枚门**：新函数是**替换**旧构造体、不是并列加一枚门；
旧签名只做一行委托。判据 ④ 另加一枚更早更便宜的同形尺（同一文件里 `approval.Options` 字面量枚数必须＝1）。
收尾读数＝**PASS**（§6）。

**P3 — `TestAC246ShippedResidentProcessOwnsItsCancelStep`（`resident_approval_246_windows_test.go:184`）⇒ 不顶到。**
它 `Contains` 的三句常量逐字在 `:269-276`：`gateAssembledClaim = "审批门已装配进本进程"`、
`gateAbsentClaim = "审批门未装配"`、`cancelStepRosterClaim = "3:cancel-task-roots"`，
真身是 `residentStatusLine()`。⇒ **本轮一个字没改 `residentStatusLine()`**（它在 `resident_approval_windows.go:488`，
本轮写面没有碰它）。新增那句 provenance 措辞**不进状态句**（见下面 P4 的理由），只进 `slog`。
收尾读数＝**PASS**（真起 `wisp.exe` 那一支，30.5s 那发里跑的）。

**P4 — `TestAC246StatusLineSaysWhatTheLegDoesNot`（同文件 `:358`，forbidden 在 `:373`）⇒ 不顶到。**
禁词三枚逐字 `"看得见"/"面板已就绪"/"已显示卡片"`，**尺面只作用于 `ra.residentStatusLine()` 的返回值**。
写之前先 grep 了一遍全仓（`grep -rn '看得见\|面板已就绪\|已显示卡片' cmd internal --include=*.go`）⇒
命中两枚：这枚尺自己，和 `internal/tools/ticket175r2_stamp_live_test.go:168` 一句无关文案。
本腿新增的措辞是 `"config"`／`"defaults (config.toml unreadable)"`／`"defaults (no host config view)"`
与一条 `slog.Info("resident gate: [risk] tier taken at construction", ...)`——**三枚禁词一枚不撞**，
而且它们**根本不进状态句**（状态句是 P3/P4 共同的尺面，往那儿加字是同时惹两枚钉的形状，本轮刻意避开）。
收尾读数＝**PASS**。

**顺带复核的两枚（任务书 §起手必做.2 之外，本轮判定与 P1 同源）：**
- P2 `TestAC246ResidentGateInjectionIsRefusedHalfAssembled` ⇒ 本轮**没碰 `run.go:601` 那枚 guard**，也不放宽任何注入条件 ⇒ 不顶到；收尾 **PASS**。
- P12 那三枚默认 300s 的钉在 `internal/agent/approval` 包内，**不在本腿定向尺的射程里**；
  其中 `TestDefaultDeadlineWallClockMeasurement` 受 `WISP_84_MEASURE` 门并 `t.Skip` ⇒ **那一枚是 SKIP 不是通过**，
  本腿没有把它读成通过，也没有跑它（票面"不许把 SKIP 读成通过"对它成立）。

**13 枚测试调用点编译面（形ⓐ 的成本核对，编排者复跑过枚数＝13）**：
`grep -rn 'newResidentApproval' cmd/wisp --include=*.go` 本腿现跑，命中 17 行＝
测试 13 枚（`resident_approval_246_windows_test.go:63/94/137/313/330/342/359`＝7、
`resident_task_source_246_windows_test.go:70/148/185`＝3、`resident_approval_live_246_windows_test.go:91/201/295`＝3）
＋产码 1 枚（`resident_windows.go`，见 §2）＋声明与注释 3 枚。
⇒ **本轮一枚测试文件没改**，旧签名保住了 ⇒ 13 枚调用点零编译冲击。
⚠ 那 3 枚 winlive 档在 `windows && winlive` tag 后面，**不带 tag 的 `go test` 根本不编它们**，
所以"没红"不能当证据。本腿**单独跑了一发带 tag 的语法检查**：
`go vet -tags winlive ./cmd/wisp/` ⇒ **`rc=0`、空输出**（件 `.scratch/wisp/probes/256/r1/winelive-compile.txt`）
⇒ 13 枚调用点连 winlive 档一起，**编得动**，这一格不是〔待验〕。

---

## §2 落地形状（file:line 逐处，全部是改完之后现读的号）

### 2.1 `cmd/wisp/resident_approval_windows.go`

| 位置 | 是什么 |
|---|---|
| `:43-58` | import 块新增两枚：`"path/filepath"` 与 `github.com/CarlosShao/wisp/internal/config`（照普查件 §7 的点名；`resident_approval_windows.go` 此前**不** import `internal/config`， census §3 复认） |
| `:88-94` | `residentApproval` 结构体新增三枚字段 `riskWindow time.Duration`（`:92`）/ `riskTimeout time.Duration`（`:93`）/ `riskProvenance string`（`:94`），注释 `:88-91`（这张回执留着，判据靠 `ra.riskProvenance` 读） |
| `:117-119` | **旧签名 `newResidentApproval()` 保留，函数体只有一行**：`return newResidentApprovalWithConfig("")` ⇒ "没有宿主配置视图"那一形 |
| `:160` | **新函数 `func newResidentApprovalWithConfig(dataDir string) *residentApproval`** |
| `:164` | `window, timeout, provenance := residentRiskGateValues(dataDir)` |
| `:165` | 三枚回执字段落进 `ra` |
| `:166-171` | `approval.New(approval.Options{...})`：`:167 UI` / `:168 Channels` / **`:169 Window: window`** / **`:170 ApprovalTimeout: timeout`** / `:171 Logf` ⇒ 实传 5 枚 |
| `:176-181`（注释在 `:173-175`） | `slog.Info("resident gate: [risk] tier taken at construction", ...)`——**`window_sec_read`/`confirm_timeout_sec_read` 之外还把 `ra.gate.Window()` 与 `ra.gate.Queue().Timeout()` 现读出来打进同一行**，那句"用的是常量档"是从**门自己身上**读的不是从本文件的算式读的，所以日志不可能和自己的值不一致 |
| `:204-218` | 三枚 provenance 具名常量：`:209 riskProvenanceRead = "config"`、`:213 riskProvenanceUnreadable = "defaults (config.toml unreadable)"`、`:217 riskProvenanceNoView = "defaults (no host config view)"`（同族先例＝`resident_ball_windows.go:154-158` 的 `hotkeyProvenance*` 三枚；**为什么不照抄 258 的值比较口径，见下面 §2.3**） |
| `:226-241` | `residentRiskGateValues(dataDir)`：空 dataDir ⇒ `(0, 0, noView)`；`config.LoadFile(filepath.Join(dataDir, configFileName), nil)` 失败或 nil ⇒ Warn 一句 + `(0, 0, unreadable)`；否则 `(L1WindowSec 秒, ConfirmTimeoutSec 秒, config)`。**两枚兜底形返回零值是有意的**：让 `approval.New` 自己文档写明的零值兜底（`Options.Window → DefaultL1Window`、`NewQueue timeout<=0 → DefaultApprovalTimeout`）成为机制，本文件因此**不持有这两枚常量的私有副本**，将来 approval 侧改默认值不会在这里留下第二套真值 |
| `:246-251` | `residentRiskConfigPathForLog`：没有宿主视图时打 `(no host config view)` 而不是空串——空串会被读成"读失败了"而不是"没去读" |

### 2.2 产码调用点，与"形ⓐ 走得通"这格的答案

`cmd/wisp/resident_windows.go:132`＝`ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`
（原来那一行在 `:126`；本腿在它**上面加了 6 行注释**（`:126-131`），所以调用本身落到 `:132`——
**位移只发生在本文件内部**，`run.go` 一行没动，见 §6 的 P7 自锁两发）。

**任务书那句"若 `dataDir` 在 `:126` 拿不到就停手上报"的撤销口令没有用到——形ⓐ 走得通，读数是跑出来的不是推出来的：**
- `rt` 在 `:42` 由 `proc.Boot(env)` 产出；`rt.Layout.DataDir` 在**同一函数体更早的两处已经在用**：
  `:66 installLogSink(rt.Layout.DataDir)`（在 `:132` **之前**，现读复认）与
  `:151 newResidentPanelManager(rt.Layout.DataDir)`（在之后）。
  ⚠ **普查件 §2 给的复认三处是 `:42`/`:66`/`:145`；那个 `:145` 本腿现读是 `:151`——漂了 6 行，漂的原因就是本腿自己加的那 6 行注释，不是别人改的。** 已在 §2.2 顶部具名记录。
- 闭合性另有两枚实跑凭据：`go build ./... rc=0`，以及判据 ⑤
  `TestTicket256ResidentBootPassesTheDataDirToTheGate`（从 AST 上要求 `runResident` 里那枚调用的**实参必须是 `*.DataDir`**）
  ⇒ **传进去的确实是 dataDir，且枚数＝1、裸签名调用枚数＝0**。
⇒ **⛔ 那 13 枚测试调用点一枚没改；⛔ `run.go`／`runSpec` 一枚字段没加；⛔ 没有任何断言被放宽。**

### 2.3 ★ 现场学到一枚没有任何上游件写过的读数（这条会影响 ⓑ 那句文案的字面）

**一旦常驻腿开始吃 `[risk]`，一份"没有 `[risk]` 段"的可读 `config.toml` 也会把 L1 窗口从恒 3s 改成 2s。**
机制：`internal/config/schema.go:453` 的 `L1WindowSec` 带 `default:"2"`，`config.LoadFile` 会把这枚 tag 填进返回值；
而 `approval.DefaultL1Window` 是 **3s**（`queue.go:116`）。⇒ 两枚"默认"**不是同一个数**，
而 2s 落在 `[MinL1Window=2s, MaxL1Window=3s]` 带内、钳位不动它。
读数凭据＝判据 ② 的用例名 `timeout only - the window then comes from the schema tag, not from the gate's compiled constant`
（种 `confirm_timeout_sec = 45` 且不种 window ⇒ 实测 `Window() = 2s`，**本腿第一版把它预期成 3s 时被自己的尺判红过**，
红句逐字留在 §4 的 M0 那发）。

⇒ 这条是本腿对**票面 §7-5／§8-6 那句 ⓑ 文案**的实质修正：那句话今天还没落到任何产码里
（本腿自己复跑 `grep -rn '只作用于跑任务的进程\|常驻腿今天用常量' cmd internal tools` ＝ **0 命中**，复认普查件 §7 末条），
而 `255-r2`／票 265 将来落它的时候，**"常驻腿今天用常量 300s／3s"这半句在 256-r1 之后已经不等于是**：
超时那一枚仍是常量档（没有宿主视图时），窗口那一枚在"文件可读但没说 window"时是 **2s 不是 3s**。
⛔ **本腿不改票面原话、不改 `docs/**` 一字**（那是 AC#3 越界面）；只把这格读数交给编排者归口。

---

## §3 四枚常驻判据（＋一枚落点副尺）：用例名／尺面／终值

全部在新增件 `cmd/wisp/resident_approval_risk_256_windows_test.go`（552 行，`//go:build windows`）。
原始输出＝`.scratch/wisp/probes/256/r1/new-rulers.txt` 与 `final-targeted.txt`。

### ① 默认档钉

**`TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants`**
- 尺面＝**唯一有效读面 `ra.gate.Queue().Timeout()`**（`gate.go:169 Queue()` ＋ `queue.go:126 Timeout()`；
  同形先例 `ticket84_no_owner_test.go:108`、`queue_test.go:48`），辅以 `ra.gate.Window()` 与 `ra.riskProvenance`。
- 两个 sub-case 都是"常量档"那一形：`no host view at all`（空 dataDir＝旧签名）与 `config.toml missing`（真目录、没文件）。
  断 `Timeout() == approval.DefaultApprovalTimeout`、`Window() == approval.DefaultL1Window`、
  provenance 必须分别等于 `riskProvenanceNoView` / `riskProvenanceUnreadable`（**那句"用的是常量档"要说得出，就是这一枚断言**）。
- **配对正控（同一 case 内）**：在 `config.toml missing` 那个目录里种一发 `confirm_timeout_sec = 45` 再建一版门 ⇒
  要求 `Timeout() == 45s` 且 provenance 变成 `config`。**没有这半，"兜底读到 300s"与"这枚门根本不读配置"无法区分。**
- 另钉一枚委托形状：`newResidentApproval()` 与 `newResidentApprovalWithConfig("")` 的 provenance 必须相等（＝"委托、不是第二套实现"的可执行定义）。
- **终值：PASS（top 1 ＋ sub 2）**。M2 突变下**它的正控那半会红**（§4）。

### ② 吃配置的正控

**`TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow`**，5 枚 sub-case：
- **正控面只用 `Queue().Timeout()`**。种 `confirm_timeout_sec = 45` ⇒ 断 `Timeout() == 45s`。
  `45` 是刻意选的：常量档 300s、schema tag 300s，**没有任何兜底路径能自己产出 45s**，所以这一枚绿不可能是恒绿。
- **`Window()` 明确只当"钳位仍在"的反控**，用例的 `winIsWhat` 字段把那句话说在**失败文案里**（不是只写在注释里）：
  `l1_window_sec = 99` ⇒ 必须读出 `MaxL1Window`（3s）、`l1_window_sec = 1` ⇒ 必须读出 `MinL1Window`（2s）。
  两枚反控都在断"钳位还在"，**都不当作"配置被吃进去了"的证据**——`gate.go:143-151` 与 `queue.go:122` 会把它们压回带内，
  拿 `Window()` 当正控就是普查件 §5／任务书判据②点名的那枚恒绿假象。
- 唯一一枚有意义的 window 正控＝种 `l1_window_sec = 2` ⇒ 断 `Window() == 2s`，因为 2s **既在带内、又不等于**这枚门在 256 之前的恒 3s。
- 另钉一枚带界：任何 sub-case 读出 `Window()` 落在 `[MinL1Window, MaxL1Window]` 之外 ⇒ 直接说"钳位没了、本轮所有 `Window()` 读数要重判"。
- **终值：PASS（top 1 ＋ sub 5）**。

### ③ 构造期定值的限制（不许写成"改配置活进程立刻跟着变"）

**`TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly`**——这条限制**被跑出来，不是写在散文里**：
1. 种 45s 建门，先断 `Timeout() == 45s`（拿不到这一枚就 fail-fast 说"本 case 无法在一个它从未有过的值上量限制"）；
2. **把同一个目录的 `config.toml 重写成 90s`**；
3. **正控**：自己用 `config.LoadFile` 再读一遍那个文件，断 `c.Risk.ConfirmTimeoutSec == 90`——
   没有这一枚，下面那句"还是 45"就和"我的重写根本没落盘"完全不可区分；
4. 限制本体：那枚**已经建好的**门必须**还是 45s**。红句写得是二义的、并且把后果说清楚：
   要么有人给 `[risk]` 加了 re-apply 路径（⇒ 本票那句限制作废、255/265 的 ⓑ 文案一起重判），
   要么超时不再从构造期取值（更糟）；
5. 同一条 case 里补上反向半枚：**重写之后新建**的那枚门必须读到 90s——
   有了"旧门守 45／新门吃 90"这一对，限制才是被钉住的，而不是被猜出来的。
6. 另钉 provenance 不随文件改动而漂（它是**构造期回执**，不是活探针）。
- **终值：PASS（top 1）**。
- **这条限制在证据里的明话**（任务书判据③要求）：**本票接上的是"带着种子值启动的这一发"，不是"改配置活进程立刻跟着变"**。
  `g.window`（`gate.go:82→:157`）与 `q.timeout`（`queue.go:64→:98`）都是构造期复制进对象的定值，
  全仓读面只有 `Gate.Window()` 与 `Queue.Timeout()` 两枚 getter，**没有任何 re-apply 路径**；
  `[risk]` 在 `internal/config/tiers.go:41` 是 `locked` 档，**既不进 `rep.Hot` 也不进 Reload／Restart**，
  所以没有任何 hook 会因为改它而叫醒任何人（普查件 §3，本腿按上面第 4、5 步实跑复认了它的后果面）。

### ④ 落点自证（不是词面尺）

**`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`**——尺面＝**`go/ast` 语法树**：
- `optionsFieldSet256(t, path)` 解析 `cmd/wisp/resident_approval_windows.go`，收集**每个 `approval.Options{...}` 复合字面量里
  真实传出去的字段名**。注释不在语法树上 ⇒ **在文件里写"risk"这个词、写十遍，也变不绿**（任务书"不许做成词面尺"的硬要求）。
- 先钉**枚数**：该文件里 `approval.Options` 字面量**必须恰好 1 枚**（两枚＝常驻腿自己铸了第二枚门，
  那就是 P1 的 `ar.gate != ra.gate` 指针同一性钉要红的方向；这枚尺是它的更早、更便宜的版本）。
- 再钉**集合**：必须恰好是 `{UI, Channels, Window, ApprovalTimeout, Logf}`——缺 `Window` 或 `ApprovalTimeout` ⇒ 指名说
  "这就是票 256 [risk] 那一半的全部"；多出任何一枚（例如有人顺手把 `Grants` 或 `WarningLead` 接上）⇒ 也红，
  理由逐字写在红句里："shipped 的字段集不再是票 256 裁过的那一份"。
- **枚数双向都响**：`len(got) != 5` 单独红一句 ⇒ 3/10（回退）与 6/10（夹带）都会红。
- **分母是从 `internal/agent/approval/gate.go` 现读的**（`declaredOptionsFields256` 走 AST 找 `type Options struct` 数名字），
  **不是硬编码 10** ⇒ `Options` 加字段时这枚尺会先说"分母动了，普查件 §1 那张表要重判"。
- **读数（本轮终值，`t.Logf` 逐字）：`landing site reading: resident leg passes 5 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)`**。
- **终值：PASS（top 1）**。

### ⑤ 落点副尺（另一半"到底接上了没有"）

**`TestTicket256ResidentBootPassesTheDataDirToTheGate`**——解析 `cmd/wisp/resident_windows.go` 的 `runResident` 函数体：
- `newResidentApprovalWithConfig` 的调用枚数**必须＝1**；
- `newResidentApproval`（**精确名比对，不是前缀**，否则会把新函数误吞进来）的调用枚数**必须＝0**，
  红句逐字带着后果："the shipped resident leg would be back to the compiled 300s / 3s, which is the exact defect ticket 256 was filed for"；
- 并且实参**必须是 `*.DataDir` 形状的 selector**（`sel.Sel.Name == "DataDir"`）⇒ 钉住"交给它的是这个进程自己的目录"。
- 为什么要有这枚：①②③④ 全都**直接构造门**，它们看不出 boot 那条线到底调了哪一枚 ⇒ 若调用点被退回，
  四枚判据会**全部照常绿**。这枚尺补的就是那个洞。
- **终值：PASS（top 1）**。

---

## §4 突变记录（红句逐字；⚠ 本腿自己也被自己的尺判红过一次，那一发留着没抹）

**还原凭据（三发共用）**：突变前把两枚产码文件复制到 `.scratch/wisp/probes/256/r1/*.pristine`，
还原后逐枚 `md5sum` 相等——
`6b33f0df963e0eed3a2e35a29aedcb15`（`resident_approval_windows.go` 与它的 pristine 同一枚）、
`a5ad69a4f8300a14ee7fd596b5c3a2fb`（`resident_windows.go` 与它的 pristine 同一枚），
并 `diff -q` 两枚均空输出；`git diff --stat -- cmd/wisp/resident_windows.go` 落回 `7 insertions(+), 1 deletion(-)`＝本腿应有的形状。
⚠ 线 endings 也核了：两枚文件都是 0×CRLF／纯 LF，突变用 python 改写时**没有把行尾换掉**（`LF=651` vs pristine `LF=653`，差的 2 行正是摘掉的两枚字段）。

### M0 —— 不是突变，是本腿第一版的预期被自己的尺判红（具名保留）

第一版判据 ② 里我把"只种 `confirm_timeout_sec`、不种 window"那一枚 case 的期望写成 `approval.DefaultL1Window`（3s）。红了：
```
resident_approval_risk_256_windows_test.go:247: Window() = 2s, want 3s for seed "confirm_timeout_sec = 45\n"; what this reading proves: nothing: no [risk] window was seeded, so this is the compiled 3s and reading it back green would not show config was consulted
```
⇒ **这不是尺写坏了，是产码告诉我一件没被任何上游件量过的事**（schema tag 的 2s ≠ 门自己的 3s）。
处理＝**改期望、不改断言强度**，并把那枚 case 改名成 `"timeout only - the window then comes from the schema tag, not from the gate's compiled constant"`，
`winIsWhat` 里把那件事实写进失败文案（§2.3 那条读数就是这么来的，交给 `255-r2`／票 265 归口）。
⛔ 本腿**没有**为了变绿删掉任何断言、没有放宽任何阈值。

### M1 —— 摘"调用点"：`resident_windows.go:132` 退回 `ra := newResidentApproval()`

件＝`.scratch/wisp/probes/256/r1/mutation-M1.txt`。指名用例＝判据 ⑤。**红，三条全出**，逐字：
```
--- FAIL: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)
    resident_approval_risk_256_windows_test.go:541: runResident calls newResidentApprovalWithConfig 0 times, want exactly 1
    resident_approval_risk_256_windows_test.go:544: runResident still calls the no-host-view newResidentApproval() 1 time(s): the shipped resident leg would be back to the compiled 300s / 3s, which is the exact defect ticket 256 was filed for. The bare signature exists for the 246 roster only (orchestrator ruling §8.2).
    resident_approval_risk_256_windows_test.go:549: runResident's call does not pass a *.DataDir value, so the gate is not being handed the host's own directory: the read would go to the wrong config.toml or to none
```
⇒ 这枚形状**只有 ⑤ 看得见**（①②③④ 直接造门，M1 之下它们全绿）——这正是任务书判据④"落点自证"要的"变了就响"。

### M2 —— 摘"新逻辑本体"：从 `approval.Options{...}` 字面量里删掉 `Window` 与 `ApprovalTimeout` 两行

件＝`.scratch/wisp/probes/256/r1/mutation-M2-final.txt`（gofumpt 之后重跑的那一发；行号与格式化前一致，
因为 gofumpt 只动了同一行内的空白对齐）。**判据 ①正控、②、③、④ 四枚红**，逐字（摘最长的文案中间的 `…` 不吞行）：
```
--- FAIL: TestTicket256ResidentGateWithoutAHostViewRunsOnTheCompiledConstants (0.00s)
    resident_approval_risk_256_windows_test.go:160: positive control failed: seeded confirm_timeout_sec = 45 but Queue().Timeout() = 5m0s, so ruler ① above was measuring a face that ignores config at all
    --- PASS: .../no_host_view_at_all          ← 兜底那一半照旧绿，符合预期：M2 没有动兜底形
    --- FAIL: .../config.toml_missing         ← 红的是它的配对正控，不是兜底断言本身
--- FAIL: TestTicket256ResidentGateTakesTheSeededRiskTimeoutAndWindow (0.02s)
    resident_approval_risk_256_windows_test.go:250: Queue().Timeout() = 5m0s, want 45s for seed "confirm_timeout_sec = 45\n" - this is AC#2's positive control, a miss means the resident leg still ignores confirm_timeout_sec
    resident_approval_risk_256_windows_test.go:254: Window() = 3s, want 2s for seed "confirm_timeout_sec = 45\n"; …
    resident_approval_risk_256_windows_test.go:254: Window() = 3s, want 2s for seed "l1_window_sec = 2\n"; …
    resident_approval_risk_256_windows_test.go:250: Queue().Timeout() = 5m0s, want 45s for seed "confirm_timeout_sec = 45\nl1_window_sec = 2\n" - …
    resident_approval_risk_256_windows_test.go:254: Window() = 3s, want 2s for seed "confirm_timeout_sec = 45\nl1_window_sec = 2\n"; …
    resident_approval_risk_256_windows_test.go:254: Window() = 3s, want 2s for seed "l1_window_sec = 1\n"; …
    --- PASS: .../window_way_too_large_-_reverse_control,_the_clamp_must_survive   ← 反控在 M2 之下照旧绿：它断的是钳位，不是配置
--- FAIL: TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly (0.00s)
    resident_approval_risk_256_windows_test.go:287: setup: the gate was not built with the seeded 45s, got 5m0s; ruler ② owns that failure, but this case cannot measure a limitation on a value it never had
--- FAIL: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
    resident_approval_risk_256_windows_test.go:448: the resident leg does NOT pass Options.Window to approval.New. That is the whole of ticket 256's [risk] half: Window and ApprovalTimeout are the two fields this leg exists to feed.
    resident_approval_risk_256_windows_test.go:448: the resident leg does NOT pass Options.ApprovalTimeout to approval.New. That is the whole of ticket 256's [risk] half: …
    resident_approval_risk_256_windows_test.go:470: resident leg passes 3 of Options' fields, want 5
    resident_approval_risk_256_windows_test.go:474: landing site reading: resident leg passes 3 of 10 declared approval.Options fields (was 3 of 10 before ticket 256)
--- PASS: TestTicket256ResidentBootPassesTheDataDirToTheGate (0.00s)   ← 它只看调用点，M2 没动调用点，绿是对的
```
⇒ **枚数 3/10 那一发是跑出来的红句，不是推断**；两枚突变合起来把 ①②③④⑤ 五枚尺的作用面**逐一区分开**了。

### M1＋M2 之后

还原 ⇒ 五枚判据 **全绿**（`final-targeted.txt`），`md5sum` 与 pristine 相等（上面），`git diff` 只剩本腿该有的形状。

---

## §5 `GRANT-DROPPED` 那一格：**没钉，具名说为什么、归口给谁**

**结论：这一格本轮钉不了，也**不该**由本腿钉。** 三个理由，逐条带出处：

1. **它断的是"第一形不再出现"，而第一形的成因是 `g.grants == nil`——`Grants` 本轮没接、也接不了**。
   本票 AC#1 被普查件 §0 量成两半，`Grants` 那一半卡在：账本唯一产码构造点在 `cmd/wisp/run.go:473→:478→:487`，
   在常驻腿门构造点**之后约 129 行**、且被 `resident_task_source_windows.go:218→:230` 那枚条件提前 return 罩着；
   `g.grants` 全仓只有 `gate.go:160` 一枚写点、`SetGrants|AttachGrants|WithGrants` 三词 0 命中。
   ⇒ 常驻腿今天仍然会在有人答 `session` 时产出 `gate.go:673-677` 那句第一形。**"不再出现"是假话，钉它就是把红写成绿。**
2. **票面 §7-6 已经定形过那枚洞**：既有仪器 `internal/agent/approval/ticket224_reply_grant_test.go:283` 只查前缀
   `approval: GRANT-DROPPED`，而前缀**两形同吃**；第二形（`gate.go:679-681`，"卡片主题在答复前已离开队列"）
   与 `Grants` 无关、挪门挪不掉它 ⇒ 若本腿去钉前缀零命中，就会顺带把第二形也禁掉，**那是票面 §7-6 明令禁止的动作**。
3. **本票 §8-5 已经把这半归口成待立票 265**（`Grants` 那半的两形代价：(a) 第二枚 mint 撞 `run.go:459-466` 的 "One mint per process"
   并造成 gate 写 id-A／bridge 读 id-B 错配；(b) 晚绑定 holder 新类型，仓里今天没这形状）。**接不上接缝之前，"落一行"无可判。**

**本腿实际钉住的（不是没钉，是钉在能钉的那一层）**：判据 ④ 用 AST **正向断言 `Options.Grants` 不出现在常驻腿的字段集里**，
红句逐字："the resident leg now passes Options.Grants. That is NOT ticket 256-r1's scope … Either a seam was minted outside that ticket, or this leg grew."
⇒ 效果是**双向闸**：本腿不许悄悄接上（越界就红），票 265 落地那天这一枚会红 ⇒ **逼着那天的腿回来改这份名册、
顺手把 `GRANT-DROPPED` 禁现那格一起接走**，不会让它变成一个没主的空格。

**归口**：`GRANT-DROPPED` 第一形禁现＋第二形仍须能出现＋"会话档不落盘"那格（票面 §7-6 附记与普查件 §3⑤ 都说今天**零用例守**）
⇒ **随 `Grants` 那半一起归待立票 265**；形状照普查件 §5 第 1 条：钉 `gate.go:675` 那句独有词组
`本机没有接入会话授权记账`，⛔ 不钉前缀。本腿**没有**动 `ticket224_reply_grant_test.go`、没有动 `gate.go` 一字。

---

## §6 门禁读数（终；逐字抄，⛔ 整包 `./...` 没跑）

所有 `go test ./cmd/wisp` 都带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`（原因见 §1.1①）。

| 尺 | 命令 | 读数 | 件 |
|---|---|---|---|
| 起手定向尺 | `go test -count=1 -v -run 'TestAC246\|TestTicket224\|TestTicket255Roster' ./cmd/wisp` | `ok ... 26.170s`，**PASS 20／FAIL 0／SKIP 0**（top-level），子用例 PASS 7 | `start-targeted.txt` |
| 收尾定向尺（含新增五枚） | `go test -count=1 -v -run 'TestAC246\|TestTicket224\|TestTicket255Roster\|TestTicket256Resident' ./cmd/wisp` | `FINAL_RC=0`，`ok github.com/CarlosShao/wisp/cmd/wisp 30.497s`，**top-level PASS 25／FAIL 0／SKIP 0** | `final-targeted.txt` |
| **P7 自锁·起手** | `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` | **PASS**（起手一发，见上） | `start-targeted.txt` |
| **P7 自锁·收尾** | 同上，收尾一发 | **`--- PASS: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)`**；同发另一枚 `--- PASS: TestTicket255RosterStillMatchesTheActualReadSites (0.54s)` | `final-targeted.txt` |
| 构建 | `go build ./...` | `BUILD_RC=0`，空输出 | `gate-build.txt` |
| vet | `go vet ./cmd/wisp/` | `VET_RC=0`，空输出 | `gate-vet.txt` |
| winlive 档编译（13 枚调用点的成本） | `go vet -tags winlive ./cmd/wisp/` | `rc=0`，空输出 | `winelive-compile.txt` |
| D22 禁令 | `sh scripts/d22scan.sh` | `D22_RC=0`；`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；`d22scan: examined 266 production Go files ...`；**末句 `d22scan: clean - no D22 ban violations`**（ban #8 覆盖 `cmd/` 98 枚 Go 文件含注释与 `_test.go` ⇒ 本腿新增件被扫过且干净） | `gate-d22scan.txt` |
| 路径长度预算（票 262 刚落地的门禁） | `bash scripts/check-path-length-budget.sh` | `PATHLEN_RC=0`；`denominator: tracked paths=5310  over-budget=57  covered by roster=57  not in roster=0`；**`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`**；同发 unit cross-check 一句：字节视图与字符视图选出**同样 57 枚** | `gate-pathlen.txt` |
| gofumpt | `$(go env GOPATH)/bin/gofumpt -l <本腿三枚文件>` | 第一发**列出** `cmd/wisp/resident_approval_risk_256_windows_test.go`（三处结构体字面量对齐＋两处 `+` 换行位置）⇒ `gofumpt -w`  applied ⇒ 第二发**空输出、rc=0**（三枚全干净） | `gofumpt.txt` |

⚠ `gofumpt` 不在裸 PATH 上，在 `$(go env GOPATH)/bin/gofumpt.exe`——记着，别把它读成"尺没跑"。
⚠ 本腿**没有跑** `go test ./...` 整包（任务书逐字禁：`cmd/wisp` 单包 240 s）。上面 d22scan 那一发里的 `PASS=35`
是 `tools/d22scan` 自己那套 `runtests.sh` 的读数，**不是 `cmd/wisp` 的整包读数**，别混着引。

### AC#3 越界面（逐字核过）

本腿两发 commit 的实际触及面（尺＝`git show --numstat --format= <sha>`，逐字）：
```
5ad8ec27  65  0  .scratch/wisp/probes/256/r1/evidence.md
2a10134e  552 0  cmd/wisp/resident_approval_risk_256_windows_test.go
          140 11 cmd/wisp/resident_approval_windows.go
            7 1  cmd/wisp/resident_windows.go
```
⇒ `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／
`frontend/**`／`design/**` **一枚都不在本腿的 commit 里**；三枚冻结件
（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）
**没读没改**；`.scratch/wisp/issues/**`（含 256 自己）与台账 `docs/reports/pending-and-issues.md` **没碰**；
`cmd/wisp/run.go` **一行没动**（形ⓐ 不需要，`runSpec` 一枚字段没加）。

⚠ 但**工作树**里今天确实躺着 `design/**` 与 `tools/d22scan/**`、`internal/agent/approval/pending_read.go` 的未提交改动。
**那不是本腿的**——起手第一发 `git status --porcelain` 在锚点 `d0847aa1` 上就已经列出它们（§1.1③），
归属是在飞的 `260-r1`／`212-r3` 与编排者暂存面。**本腿用 `git add` 显式 pathspec ＋ `git commit -- <pathspec>` 同发一条命令，
把它们一枚都留在门外**（这正是票面 §7-1 那枚 `A532` 归属事故定式的形状）。⇒ **AC#3 判定＝通过，且判定尺是 `git show --numstat`，不是 `git diff`。**

---

## §7 判不动的地方／没做成的（逐条具名，⛔ 没有一条用"应该没问题"填）

1. **`Grants` 那一半没做**（票面 AC#1 的另外半句）。⛔ 不是"暂缓"，是**判死**：
   构造时刻没有账本对象、`g.grants` 无晚绑定入口、`runSpec` 无账本字段（普查件 §0/§2 逐处 file:line）。
   **归口＝待立票 265**（票面 §8-5），它的 AC#0 要先把 (a) 第二枚 mint 与 (b) 晚绑定 holder 两形的代价与线程／时序约束量齐。
   ⇒ **本票 AC#1 整格不许由本腿翻勾**：本腿只交付了它被裁的那半，`Grants` 那半今天仍然不落一行。
2. **`GRANT-DROPPED` 禁现那格钉不了**，见 §5 三条理由。**归口＝票 265**。
3. **无控制台那一支（双击／Explorer 拉起）里 AC#2 的"不再出现"机读不读得了——本腿〔待验〕**。
   普查件 §2/§5 说链路接不上（`resident_task_source_windows.go:218→:230` 提前 return ⇒ 无 ledger、无 `runReplyLoop`、无 `session` 读者），
   但那是**只读层面的结论，本腿没有跑那一支**。按任务书"任何'这条路接不上'的结论属【否证】，必须自己把那条路的作用面跑一遍"
   的口径，**本腿没跑 ⇒ 只能记〔待验〕，不许记成"做不到"，也不许记成"成立"**。
   本轮改的是 `[risk]` 两枚，**与那一支的门是否存在无关**（门仍然在 `resident_windows.go:132` 无条件建，ⓑ 那条裁定的理由读数照旧成立），
   所以这一格不因本腿而变好也不因本腿而变坏。
4. **`internal/config/unwired.go:119-120` 那两行今天变成半谎**（普查件 §4.2 P14 预告的那一格，本腿落地后**成事实**）。现读逐字：
   ```
   "risk.confirm_timeout_sec": "consumed: cmd/wisp/run.go builds the approval timeout from it",
   "risk.l1_window_sec":       "consumed: cmd/wisp/run.go builds the L1 auto-approve window from it",
   ```
   现在这两枚值**也被 `cmd/wisp/resident_approval_windows.go` 在常驻腿里建进门**，那两行是产码里唯一的"谁吃这两项"登记表，
   于是它们**少说了一个消费者**。⛔ `internal/config/**` **不在本腿写面**（写面逐字只有 §2 那三枚文件＋证据目录）⇒
   **本腿不改它一字，具名上报归口**：`255-r2` 的 ⓑ 那半格（或票 265）落那句同源文案时，把这两行的射程一起改掉。
   今天**不会因此判红**——`TestEveryLockedSectionKeyIsAccountedFor` 只要求"每个 locked 键有一行"，
   不校验那句 prose 的完备性（P14 的读数，本腿复认了那两行文字与守卫自述 `:112-117`）。
5. **`hotRowClaims` 里没有 `risk` 这一行，本腿复跑了**：`grep -n '"risk' cmd/wisp/config_readers_255.go` ⇒ **0 命中**，
   `sectionReadSites` 的键全集也没有 `Risk`。⇒ 票面 §8-6 那句"ⓑ 若要由登记表同源产出，那张表今天没有这一格"**复认成立**，
   归 `255-r2`／票 265；本票不代它落。
   ⚠ 顺带一枚对本腿有利的读数：正因为 `rosterCoversSection("risk")` 为假，新增的 `cfg.Risk` 读点**不在**
   `TestTicket255RosterStillMatchesTheActualReadSites` 的射程里 ⇒ 收尾那发 P8 同形尺**实测 PASS**（不是推理）。
6. **`cmd/wisp/resident_task_source_windows.go:13` 那句注释引用的行号被本腿写漂了**（本腿自己造成，具名不抹）：
   原文 `246-a1 … found AskOnTaskRoot (resident_approval_windows.go:261)`，
   而 `AskOnTaskRoot` 现在在 **`resident_approval_windows.go:390`**（本腿在文件上半部加了 129 行：委托＋新函数＋三枚常量＋两枚 helper）。
   ⛔ 那枚文件**不在本腿写面**，而且它是**注释不是尺**（普查件 §6 对 `ticket220_l1_window_read_test.go:4` 那处同形事故定的口径＝不红）。
   ⇒ 本腿**没有顺手改它**，只在这里记账：谁下一次动那枚文件，把 `:261` 一起改成现读号。
7. **`ⓑ` 那句文案本腿没落**（`grep` 复跑 0 命中，§2.3）——**并且本腿给它添了一枚必须改的内容**：
   "常驻腿今天用常量 300s／3s"这半句在 256-r1 之后**不再等于是**（见 §2.3 那枚 2s 读数）。
   ⛔ 本腿不改 `docs/**`、不改票面一字；**这一格是交给编排者的读数，不是本腿自行修的文案**。
8. **`WarningLead`／`MaxPending`／`MaxTracked`／`Clock` 四枚零值字段没碰**。票面没要，普查件 §8-3 说"该不该补"不是读数问题、没有量具。
   本腿**没有**顺手接（接了会撞上判据 ④ 那枚"多出任何一枚都红"的闸，正是要它拦的行为）。
9. **判据 ③ 只钉了 timeout 那一侧的"构造期定值"**。`Window()` 那一侧的"不随改而动"本腿**没单独跑**——
   因为 `Window()` 被钳在 3s（`l1_window_sec=99` 与不种值都可能读出 3s），拿它做这枚限制**读出恒绿的可能性高于读出事实**；
   带界那枚断言（`[Min,Max]`）与两枚反控 case 是本腿对 window 侧做到的全部。**这条限制本身对两枚字段同形成立**
   （`g.window`/`q.timeout` 都在构造期复制、全仓无 re-apply，普查件 §3；本腿对 timeout 实跑复认）。
10. **winlive 那 3 枚（`resident_approval_live_246_windows_test.go:91/201/295`）本腿只证明"编得动"，没证明"跑得绿"**：
    `go vet -tags winlive` rc=0（真窗真借键那一支要桌面，任务书也写了"默认 CI 不跑"）。⇒ **那一发的运行时读数〔待验〕**。
11. **本腿没做的两件交付节奏外的事，明写以免被当成做了**：没跑 `go test ./...` 整包（禁）；
    没跑 `internal/agent/approval` 包内那套 P9 尺（普查件 §4.2 P9 判"不红"，前提是包内 `fakes_test.go:250-279` 的 `newGate`
    逐字段拷贝——本腿**没往 `Options` 加新字段**，那枚拷贝面没被动到；⚠ 若将来加字段，那台仪器会**静默丢弃**它，
    这是 P9 记着的形状，不是本腿新造的）。

---

## §8 Git 流水（只 commit 不 push；add 与 commit 同发一条命令、显式 pathspec）

| # | sha | 内容 | `git show --numstat --format=` 枚数 | pathspec |
|---|---|---|---|---|
| 1 | `5ad8ec27` | 骨架：证据件九节标题 | **1 枚** | `.scratch/wisp/probes/256/r1/evidence.md` |
| 2 | `2a10134e` | 落地：两枚 `[risk]` 字段进常驻门＋五枚判据 | **3 枚** | `cmd/wisp/resident_approval_windows.go` `cmd/wisp/resident_windows.go` `cmd/wisp/resident_approval_risk_256_windows_test.go` |
| 3 | `2708289f` | 证据件写满（§0–§9 全填，占位符 0） | **1 枚** | `.scratch/wisp/probes/256/r1/evidence.md` |
| 4 | 本发 | §8 补第三发的 sha ＋ **已提交树上的一发复认**：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -run 'TestAC246\|TestTicket224\|TestTicket255Roster\|TestTicket256Resident' ./cmd/wisp` ⇒ **`rc=0`，`ok github.com/CarlosShao/wisp/cmd/wisp 27.500s`**（件 `committed-tree-check.txt`）；同一发的 `git status --porcelain -- cmd internal tools docs frontend design` 里仍然躺着 `design/**`、`tools/d22scan/**`、`internal/agent/approval/pending_read.go`——**逐条与起手 §1.1③ 那一份对照过，全是别人的在飞面，本腿一枚没碰** | 1 枚 | `.scratch/wisp/probes/256/r1/evidence.md` |

**没做的事，逐条对票面**：⛔ 没 `git add -A`／`.`；⛔ 没 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；
⛔ 没 push；⛔ 没碰 `.scratch/wisp/issues/**` 一枚复选框（含 256 自己那三格——**AC#1 整格不归本腿翻勾，理由见 §7-1**）；
⛔ 没碰台账；临时件只建不删（本目录里 `*.pristine`、`wisp.test.exe`、各 `.txt` 读数全留着）。

**共享工作树里的时间线读数**：起手锚点 `d0847aa1`；本腿第 2 发之前，HEAD 已被并发腿推进到
`b831ef02`（`ledger(A595)`）／`84dbda52`（`260-r2 骨架`）／`263658e0`（`264-a1`）等，
本腿第 2 发落在它们之后＝`2a10134e`。**没有一枚别人的 commit 被本腿吞掉或改写**（尺＝上面每发的 numstat）。



---

## §9 本件自量（成稿判据＝盘上，不是工具回执）

| 尺 | 读数 |
|---|---|
| 占位符词面扫描（四枚词的字面见票面 256 §7-1 那句成稿判据，本件不重复其字面，否则这一行会把自己的扫描命中） | **0 命中** |
| `wc -c` | **44337 bytes** |
| `wc -l` | **465 lines** |
| 节次 | §0–§9 十节全填（§0 结论句／§1 起手预检／§2 落地形状／§3 四枚判据／§4 突变／§5 GRANT-DROPPED／§6 门禁／§7 判不动／§8 Git 流水／§9 自量） |
