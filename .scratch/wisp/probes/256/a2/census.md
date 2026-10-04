# 票 256 `256-a2` 只读普查（成稿）——"ⓑ 不移动 `approval.New`"与"AC#1 补 `Options` 字段"是**同件事的两半还是互斥**

**⚠ 本件由编排者代落**：派单把交付件写给运行类型**没有写文件工具**的腿（`Explore` 型：有 Read/Grep/Glob，无 Write/Edit）⇒ 该腿**零 commit、件未建立**，全文以回报正文交回。这是同一小时内第二次（第一次＝`167-a2`，账 `A590` §1），定式收紧见 `A591`。
**代落口径**：读数全文保留；编排者对**决策关键的那几处**自己复跑并在此件内具名标注〔编排者复跑〕，其余为〔仅自述〕。
**现量锚点**：腿自报 HEAD `528bf7bd`（逐字 `528bf7bd 174-a4 只读普查终态…`）；零 Go 命令；`git status --porcelain -- cmd internal` 现量为空。

## §0 两个问句是不是同一件事（本票今天最值钱的一格）

判死：**票面 AC#1 那一句是两半，两半的命运不同。**
- **`[risk]` 那两项＝同件事的两半，合法。** 补 `Window`／`ApprovalTimeout` 两枚字段（`internal/agent/approval/gate.go:29`、`:32`）**不需要任何在 `:109` 那一刻还不存在的对象**：值来自 `config.toml`，`config.LoadFile` 的返回是纯数据，常驻腿此刻拿得到（链路见 §3）。⇒ **原地补字段＝不移动，ⓑ 裁定不被触碰。**〔编排者复跑：`cmd/wisp/resident_approval_windows.go:109` 逐字 `ra.gate = approval.New(approval.Options{`；`cmd/wisp/resident_windows.go:126` 逐字 `ra := newResidentApproval()`〕
- **`Grants` 那一枚＝互斥的一半。** 逐处现读，不是推理：
  - `Options.Grants`（`gate.go:59`）在 `New` 里只被赋值一次：`gate.go:160` `grants:   o.Grants,`〔编排者复跑命中 `:160`；⚠ 我第一发 `grep 'grants: o.Grants'` **零命中**＝**空格敏感的坏尺**，换成 `grants:[[:space:]]*o.Grants` 才中——本仓第 N 次撞上"词面尺因空白指错"，记我〕；字段声明 `gate.go:88`。
  - **全仓没有任何晚绑定入口**：`grep -rn "SetGrants|AttachGrants|WithGrants" --include=*.go cmd internal tools` ＝ **0 命中**〔编排者复跑＝0〕；`g.grants` 的写点只有 `gate.go:160` 一枚，读点两枚＝`gate.go:673`（判 nil）与 `gate.go:684`（`Record`）〔两枚均编排者复跑命中〕。
  - 账本**唯一**产码构造点在 `assembleRuntime` 之内：`cmd/wisp/run.go:473`〔编排者复跑：腿写 `:474`，**漂一行**，逐字 `sessID, err := session.Mint()`〕→ `run.go:478` `session.NewLedger(...)`〔复跑命中〕→ `run.go:487` `rt.session = ledger`〔复跑命中〕（`assembleRuntime` 函数体起于 `run.go:382`；`rt.session` 字段声明 `run.go:282`）。
  - 常驻腿那枚 `approval.New` 在 `resident_windows.go:126` → `resident_approval_windows.go:105`/`:109`；而 `assembleRuntime` 在常驻腿里要到 `resident_windows.go:255`（`src := startResidentTaskSource(rt, ra)`）→ `resident_task_source_windows.go:215` → `:265` `run, code := assembleRuntime(spec)` 才跑。⇒ **账本的构造点在门的构造点之后约 129 行**，且中间隔着一枚条件提前 return（`:218 console := interactiveStdin()` → `:230 return nil`）。
  - `runSpec` 今天**没有**注入账本的字段：字段全集现读 `run.go:102-175`＝`argv/stdout/stderr/dataDir/notify/probeSink/now/onRuntime/modeConfirm/sink/reply/replyVeto/gate/ui/cards/taskCtx`〔仅自述〕。
  - 任务管线那侧**确实已有一份可借**（`run.go:576 grantRead = rt.session` → `:758 Grants: grantRead` 给 bridge 读；`run.go:577 grantWrite = rt.session` → `:618 Grants: grantWrite` 给自建门写）〔仅自述〕，⚠ 但那份在常驻腿里是**读侧已接、写侧孤儿**：注入分支 `run.go:601-609` 走 `rt.gate = s.gate`，而那枚 gate 的 `g.grants` 早在 `:109` 定成 nil。
- ⇒ 要在 `:109` 原地补 `Grants: X`，`X` 只有两条来源：**(a) 新铸一枚第二 ledger**（与 `run.go:459-466` 注释那句 "One mint per process" 成双，且会造成 gate 写 id-A／bridge 读 id-B 的错配）；**(b) 新造一枚晚绑定 holder 类型**（仓里今天没有这种形状）。两条都**不是"补字段"**，是**新增接缝**。**具名写死：AC#1 的 `Grants` 半句与 ⓑ 的"不移动"不同命题，不可用"原地补字段"消解。**

## §1 `approval.Options` 完整字段名册（`internal/agent/approval/gate.go:17-60`，10 枚）

| # | 字段 | 声明 | 常驻腿 `:109-113` 实传 | 跑任务腿 `run.go:612-619` 实传 | `New` 里的兜底／钳位 |
|---|---|---|---|---|---|
| 1 | `UI` | `:20` | **有值** `ra.ui` | 有值 `rt.ui` | `:135-138` nil ⇒ `UIFuncs{}`（fail-closed） |
| 2 | `Clock` | `:22` | nil／零值 | nil／零值 | `:127-130` nil ⇒ `SystemClock{}` |
| 3 | `Channels` | `:25` | **有值** `approval.NewChannels()`（空册） | 有值（空册） | `:139-142` nil ⇒ `DefaultChannels()` |
| 4 | `Window` | `:29` | **零值 0** | 有值 `cfg.Risk.L1WindowSec`(`run.go:615`) | `:143-151` `<=0→3s`／`<2s→2s`／`>3s→3s` |
| 5 | `ApprovalTimeout` | `:32` | **零值 0** | 有值 `cfg.Risk.ConfirmTimeoutSec`(`run.go:616`) | 透传 `NewQueue`(`:158`)，`queue.go:85-87` `<=0→300s` |
| 6 | `WarningLead` | `:35` | 零值 | 零值 | `queue.go:88-90` ⇒ `30s` |
| 7 | `MaxPending` | `:37` | 零值 | 零值 | `queue.go:91-93` ⇒ `8` |
| 8 | `MaxTracked` | `:41` | 零值 | 零值 | `:152-155` ⇒ `64` |
| 9 | `Logf` | `:43` | **有值** `ra.residentAuditf`(`:112`) | 有值 `rt.auditf` | `:131-134` nil ⇒ noop |
| 10 | `Grants` | `:59` | **nil／零值** | 有值 `grantWrite`(`run.go:618`) | **无兜底**；`:160` 原样搬进 `g.grants` |

⇒ 常驻腿实传 **3／10**，零值 **7／10**（票面点名的三枚＝`Window`、`ApprovalTimeout`、`Grants`）。
`Grants` 的类型与持有：接口本体 `gate.go:66-68` `GrantRecorder{ Record(ctx, tool, pattern) (int64, error) }`（**只有 `Record`**；写侧专用，`:44-58` 注释明说 gate 不许读回）。边界是结构性的：`grep -rn "internal/session" --include=*.go internal/agent internal/tools` 非测试命中**全是注释**（`gate.go:52`/`:65`、`internal/tools/bridge.go:95`、`internal/tools/grant.go:38`），**没有一枚 import** ⇒ `internal/agent/approval` 见不到 `internal/session` 这个包。唯一实现者 `*session.Ledger`（`internal/session/grants.go:146`→`:157 InsertGrant`）。持有者＝`Gate`(`gate.go:88`) 与 `agentRuntime`(`run.go:282`)；建点唯一＝`run.go:478`（前置 `:473 Mint`）；生命周期＝本进程（`run.go:459-466`：id 不可重算、进程亡即失效）。

## §2 会话账本在常驻腿的可达性（调用链逐处）

造它的人（1 枚产码）：`run.go:473` `session.Mint()`（`internal/session/session.go:94`）→ `run.go:478` `session.NewLedger(...)`（`internal/session/grants.go:111`；`Store` 用 `mem`＝`run.go:452 memory.Open(s.dataDir)`，`:455 rt.store = mem`）→ `run.go:487 rt.session = ledger`。`internal/perm` **不造它**（`perm.New` 在 `run.go:654`，管 permission_mode，另一条通道）；`internal/agent/approval` **不造它**（只见接口）。
相对位置：门＝`resident_windows.go:126` → `resident_approval_windows.go:105` → `:109`；账本＝`resident_windows.go:255` → `resident_task_source_windows.go:215` → `:218` →（`:230 return nil`）→ `:265` → `run.go:382` 起体 → `:473/478/487`。⇒ **账本在门之后，且中间隔着提前 return；"之前还是之后"判死＝之后，无例外。**
任务管线那侧有一份可借，但借不到构造时刻：`rt.session` 在 `assembleRuntime` **返回之后**才对外可见（`resident_task_source_windows.go:290 src.run = run`），而 `g.grants` 在 `gate.go:160` 一次写定、无晚绑定入口（§0）。
**无控制台形状（双击／Explorer）**：`approval_reply_stdin_windows.go:41-52`（`GetStdHandle`/`GetConsoleMode` 失败 ⇒ nil）→ `resident_task_source_windows.go:218` → `:224-230` 打印 `taskEntryDisabledClaim`(`:85-89`) 后 `return nil` ⇒ **连 ledger 都不存在**，`Grants` 无从谈起（ⓑ 那条裁定的理由读数这条我逐处复认为成立）。

## §3 `[risk]` 两项从常驻腿取值的路径

字段与档位：节声明 `internal/config/schema.go:118` `Risk RiskSection`；结构体 `:447-472`；`ConfirmTimeoutSec` **`:450`**（`default:"300"`）、`L1WindowSec` **`:453`**（`default:"2"`）。
档位：`internal/config/tiers.go:41` 逐字 `"risk": "locked"`〔编排者复跑命中；另 `:83` 同键在第二张表里〕⇒ **既不是 hot 也不是 reload／restart**。应用表 `manager.go:248-251` `planLocked("risk", …, func(){ cur.Risk = fresh.Risk })`；方向审计 `manager.go:431` `riskDirection`（`l1_window_sec` 的 loosen/tighten 在 `:442-446`），**`confirm_timeout_sec` 没有方向分支**（`:431-455` 全文读）⇒ 改超时**不算 loosening**，走 `:335-341` 的 quiet-apply。
通知面：`manager.go:198-199` `OnReload` 只对 `rep.Reload`；`:201-202` `OnRestartPending` 只对 `rep.Restart`。`[risk]` 既不进 `rep.Hot`（`:293-304` 的 hotSections 逐字不含 risk，且 `:294`/`:299` 两道 panic 守卫锁死这张表只能装 TierRegistry 里 `"hot"` 的段）也不进 Reload／Restart ⇒ **没有任何 hook 会因 `[risk]` 改动叫醒任何人**；后果只落在 `rep.Locked[i]`。产码 hook 持有者只有两枚：`cmd/wisp/config_reload.go:115`、`cmd/balldebug/main.go:244`。
**常驻腿此刻手里有没有 config manager＝没有**（负向句三处读满）：产出者 `manager.go:101 NewManager(path, res)`；产码调用点三枚＝`run.go:409`（`assembleRuntime` 体内，门之后）、`panel_inbound.go:230`（`newComposerDispatchChain`，运行时刻在 `resident_windows.go:151 startResidentPanel` 之后、**不在 `:126` 之前**）、`cmd/balldebug/main.go:231`；持有者字段 `run.go:268 mgr *config.Manager`。`resident_approval_windows.go:43-56` 的 import 块**不含 `internal/config`**〔编排者复跑：该文件 import 逐枚＝`context/errors/fmt/log-slog/sync/time` ＋ `approval/ball/risk/statemachine/tools`，确无 `internal/config`〕⇒ **要用得加 import**。
**但值取得到，形状已有先例**：`resident_windows.go:177` 与 `:200` 两处 `config.LoadFile(filepath.Join(rt.Layout.DataDir, configFileName), nil)`（票 258 的 per-use 读形；`:168-171` 注释逐字说明"拥有 `rt.mgr` 的那条任务管线在此刻还不存在，所以这条链每次调用重读 `config.toml` 而不持有 Manager"）。`rt.Layout.DataDir` 在 `:126` 之前已可用（`:42` `proc.Boot`、`:66`、`:145`）；常量 `cmd/wisp/secret.go:63`。同类先例 `panel_resident_windows.go:201`、`providers.go:98`、`models.go:184`。
⚠ **接了之后的真实差异只有一枚**：`Window` 被 `gate.go:149-151`＋`queue.go:122`（`MaxL1Window=3s`）钳死，`schema.go:453` 的 2s 也出不了 3s ⇒ 与母票 248 AC#10 "只有超时是真差异"同形。
⚠ **接了也不会随改而动**：`g.window`(`:82`→`:157`) 与 `q.timeout`(`queue.go:64`→`:98`) 都是构造期定值，读面只有 `gate.go:176 Window()` 与 `queue.go:126 Timeout()` 两枚 getter，**全仓无任何 re-apply 路径** ⇒ **AC#2 的正控只能钉"带着种子值启动这一发"，钉不了"改一个活着的进程"**；这条对跑任务腿同样成立。

## §4 既有钉名册（穷举，不抽样）

### 4.1 `grep -rn "approval.New(" --include=*.go cmd internal tools` ＝ 10 行
产码 **2 枚**＝`resident_approval_windows.go:109`（实传 `UI`/`Channels`/`Logf`）与 `run.go:612`（实传 `UI`/`Channels`/`Window`/`ApprovalTimeout`/`Logf`/`Grants`）；测试 **7 枚**＝`approval/fakes_test.go:274`、`approval/ticket224_reply_grant_test.go:133`、`approval/ticket84_no_owner_test.go:107/230/247`、`tools/loop_approval_test.go:94`、`tools/wiring_test.go:83`；注释 1 处＝`approval/doc.go:28`。
⚠ 这把尺漏枚同包构造点，穷举须补 `New(Options{`：`approval/ticket220_l1_window_read_test.go:203` 与 `:383`（两枚都已传 `Window: 3*time.Second, Grants: grants`）⇒ **gate 构造点全仓合计＝11 枚（产码 2／测试 9）**。〔编排者只复跑了产码两枚的位置；枚数〔仅自述〕〕

### 4.2 会红的具名清单（`256-r1` 的雷区）
| # | 用例名／位置 | 断什么 | 落地后 |
|---|---|---|---|
| **P1** | `TestAC246ResidentPipelineAsksThroughTheOneGate`（`cmd/wisp/resident_task_source_246_windows_test.go:68`，钉在 `:89`） | **指针同一性**：`ar.gate != ra.gate` 即 `t.Fatalf`（"two gates in one process"）；`:71` 调 `newResidentApproval()` **无参** | **会红两形**：① 为带 grants 而立第二枚 gate ⇒ `:89` 必红；② 给 `newResidentApproval` 加参数 ⇒ **编译红**，冲击面＝**13 枚测试调用点**（`resident_approval_246_windows_test.go:63/94/137/313/330/342/359`＝7；`resident_task_source_246_windows_test.go:70/148/185`＝3；`resident_approval_live_246_windows_test.go:91/201/295`＝3，tag `windows && winlive`）＋产码 1 枚（`resident_windows.go:126`）〔**13 枚＝编排者复跑命中**〕 |
| P2 | `TestAC246ResidentGateInjectionIsRefusedHalfAssembled`（同文件 `:146`，拒句 `:168` 逐字 `注入的审批门必须连同它自己的界面与会话账本一起递进来`，真身 `run.go:602-605`） | `assembleRuntime` 必须**拒绝**半注入 | 不放宽 `run.go:601` guard 则**不红**；放宽即红 |
| P3 | `TestAC246ShippedResidentProcessOwnsItsCancelStep`（`resident_approval_246_windows_test.go:184`，`Contains` 两句常量 `:271/:272`，真身 `residentStatusLine()` `resident_approval_windows.go:359-369`） | 真起 `wisp.exe` 读状态句 | 改那两句字头 ⇒ **红**；补字段不动文案 ⇒ 绿 |
| P4 | `TestAC246StatusLineSaysWhatTheLegDoesNot`（同文件 `:358`，forbidden loop `"看得见"/"面板已就绪"/"已显示卡片"`） | 状态句不许承诺卡片可见 | 新句撞这三词之一 ⇒ 红（硬扫） |
| P5 | `…EscChannelStaysUnloadedWithoutABallWindow`／`…CardWithNoWindowFailsClosedThroughTheRealGate`／`…ChannelNeedsBothWindowAndExecutor`／`…VetoSentenceWithNoCard`（`:62/93/311/341`，读 `ra.gate` 在 `:64/73/77/82/321`） | 通道加载／fail-closed／双条件／无卡句 | 不读 `Window/Timeout/Grants` ⇒ 只要 P1 的签名冲击不炸就**绿** |
| P6 | `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`（`resident_approval_live_246_windows_test.go:72`，tag `windows && winlive`） | 真窗真借键 | 默认 CI 不跑；本机 winlive 受 P1 签名冲击 |
| **P7 ★** | `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`（`cmd/wisp/config_receipt_255_test.go:179`，表体 `config_readers_255.go:87` 起，尺在 `:209`＋要求 `citesChecked >= 8` `:219-224`） | 把 `hotRowClaims` 每条 `file.go:LINE [token]` 现读源文件逐行比对 | **本票最硬的间接钉**：被点的 `run.go` 行号四处＝`run.go:424`/`:435`/`:991`/`:1014` ⇒ **任何让 `run.go` 第 424 行以后整体位移的改动（含给 `runSpec` 加字段）＝四处 evidence drift ⇒ 红**；另 `models.go:163/183`、`resident_ball_windows.go:276`、`panel_config_store.go:92/96`、`panel_resident_windows.go:207`、`logsink.go:149`、`observe/logging.go:50`、`ball_windows.go:64` 同受这同一把尺管 |
| P8 | `TestTicket255RosterStillMatchesTheActualReadSites`（`:225`；表体 `config_readers_255.go:176-204`，射程 `:312-327`→`:332-…`） | 每段"谁读它"的文件集合与登记表逐枚相等 | **不在射程**：`hotRowClaims` 里没有 `risk` 行 ⇒ `rosterCoversSection("risk")`＝false ⇒ 新增 `cfg.Risk` 读点**不红**（三枚函数体逐行读完，非推理） |
| P9 | `internal/agent/approval` 全套（`fakes_test.go`／包内 25 个文件，例 `batch_test.go:90/123`、`window_test.go:22/59/87`、`queue_test.go:48/192`、`ticket87_veto_l2_test.go:60/129/159/203`、`ticket97_alias_direction_test.go:75/149`、`ticket220_l1_window_read_test.go:248/296/338/380`、`ticket242_binding_test.go:19/34/58/87`、`ticket242_panelface_test.go:30/53`、`ticket146_liveapprovals_backing_test.go:208/264`、`pending_read_test.go:111/190/224`） | 全部**包内自建 Options**，不吃 `cmd/wisp` | **不红**。⚠ 但 `fakes_test.go:250-279` 的 `newGate` **逐字段枚举拷贝**（`:261-269` 只搬 `Window/ApprovalTimeout/WarningLead/MaxPending/MaxTracked/UI/Clock/Channels/Logf`，**不搬 `Grants`，也无法搬任何未来新增字段**）⇒ 往 `Options` 加字段时包内仪器**静默丢弃**（不红、也测不到；要它搬就得改这枚既有钉） |
| P10 | `TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped`（`ticket224_reply_grant_test.go:273`，`:283` 查前缀 `approval: GRANT-DROPPED`） | nil recorder ⇒ 放行照旧＋必须说 dropped | **不红**；⚠ 它就是票面 §7-6 那枚"只查前缀、两形不加区分"的仪器，⛔ 不许为"前缀零命中"去动它 |
| P11 | `TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun`（`cmd/wisp/ticket224_assembly_test.go:80`，`:193-196` 反向钉不许 `GRANT-DROPPED`、`:190-192` 正钉 `GRANT-RECORDED … grant_id=%d`） | run 腿不许说 dropped | **不红**；这枚就是 AC#2"必须不再出现"那一形状在 run 腿的**现成同形尺**，可照其形而不必改它 |
| P12 | `TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline`／`TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds`／`TestDefaultDeadlineWallClockMeasurement`（`ticket84_no_owner_test.go:56/98/222`，`:108`/`:232` 读 `Queue().Timeout()`；`:224` 受 `WISP_84_MEASURE` 门并 `t.Skip`） | 默认 300s 有限上界 | 不红；⚠ `:224` 是 **SKIP 不是通过**（票面"不许把 SKIP 读成通过"对它成立） |
| P13 ★ | 三枚冻结件：`internal/panel/l2_grant_boundary_test.go`（`:1229`/`:1530`/`:1547`/`:1595`/`:1807`/`:1959`/`:2168`）、`internal/perm/ticket90_persist_test.go`（`:138`/`:181`/`:220`/`:291`/`:336`）、`internal/panel/tokens_fourway_test.go:439` | panel 报文／grant 词面／permission_mode 持久化／设计令牌 | 全不经 `cmd/wisp` 门构造点 ⇒ **不红**（前提：不往 `PanelAPI` 加任何 `Allow` 变体；`approval/ui.go:147`/`:158` 的 `AllowSession` 只在 `NativeAPI` 上） |
| P14 | `internal/config/unwired.go:119-122`（`lockedKeyDisposition` 两句）＋`TestEveryLockedSectionKeyIsAccountedFor`（守卫自述 `:113-117`） | 只要求"每个 locked 键有一行"，**不校验那句 prose 的完备性** | **不红**，但常驻腿一旦开始吃 `[risk]`，那句"consumed: cmd/wisp/run.go builds …"就变成**半谎**（产码里唯一的"谁吃这两项"登记表）；`unwired_test.go:66-67`/`:188-189` 默认档与非默认档都绿（不在 `unwiredKeys` 表里） |
| P15 | run 腿种 `[risk]` 的仪器（`run_mode101_test.go:110`、`approval_reply_201_test.go:145`、`panel_pump_test.go:136`、`unwired_test.go:66/188`、`migrate_test.go:27/99`、`run_test.go:331` 注释） | 写入／断言 `[risk]` 数值 | 不直接红；⚠ 但它们证明"这两项的正控今天**只在 run 腿有载体**" |
| P16 | `tools/d22scan/allowlist.txt`（7 行，全是 `pathresolver-bypass`／`mirror-hash`／`bare-goroutine`，**无 approval 条目**） | 禁令豁免表 | 不新增裸 `go` 语句、不绕 C26 ⇒ 不红；票面 AC#3 把它列为越界面 ⇒ 落地腿不许碰 |

**"常驻腿『会话档不落盘』这一格今天零用例守"**＝复认票面 §7-6 附记成立：`GRANT-DROPPED` 在 `cmd internal tools` 的 `*.go` 里总命中 **4 枚**＝产出两枚（`gate.go:674`/`:680`）＋断言两枚（`cmd/wisp/ticket224_assembly_test.go:194`、`approval/ticket224_reply_grant_test.go:283`）＋注释两处（`resident_task_source_windows.go:43`、`run.go:596`）；**两枚断言都不经 `resident_approval_windows.go`**。⇒ **钉住"三项形状／常驻没有 grants"的用例枚数＝0**（没有任何用例断言 `:109-113` 的字面字段集；`g.grants` 是小写私有，包外读不到）。

## §5 `GRANT-DROPPED`：在哪产出／什么条件／AC#2 今天能不能机读断言

- 产码点两枚（穷举）：`gate.go:673-677`（第一形，`if g.grants == nil`，句尾 `"(本机没有接入会话授权记账，本次按「仅本次」放行，没有落盘任何规则)"` 在 `:675`；**发生在 `allowScoped` 成功之后** `:669` ⇒ 放行照旧）与 `gate.go:679-681`（第二形，`if tool == ""`，句尾"卡片主题在答复前已离开队列，未落盘任何规则" `:680`；`tool` 来自 `:640 g.q.sessionSubject(corr)`，`:646` 置空）。函数体 `:658-694`。〔`gate.go:673`/`:684` 编排者复跑命中；其余行号〔仅自述〕〕
- 投递链：`NativeAPI.AllowSession`(`gate.go:653`) ← `replies.go:351→:363`（前置 `:358` 无 `Grant` ⇒ `ErrRouteHasNoAllow`；`:361` 无 gate ⇒ `ErrNoGateAttached`）← `approval_reply.go:257→:259` ← `:566-567 case "session"` ← `runReplyLoop` ← `attachReplyListener`(`:442`/`run.go:801-803`) 或常驻腿 `newReplySurface`(`resident_task_source_windows.go:288`，**被 `:287 if console != nil` 罩着**)。落盘者＝`grants.go:146 Record`→`:157 InsertGrant`，只在 `gate.go:684` 被调；`:673`/`:680` 两枚 return 都在它之前。
- **今天能不能机读断"第一形不再出现"＝能（有控制台形状那一支）**，可钉形状（本轮不写码）：
  1. **区分两形现成可做**：`本机没有接入会话授权记账` 这串在 `*.go` 里唯一产码点＝`gate.go:675`（另一命中 `run.go:596` 是注释）；第二形 `:680` 不含这串 ⇒ **尺应钉这句独有词组**，⛔ 不许钉前缀 `approval: GRANT-DROPPED`（前缀两形同吃＝§7-6 那个洞；`ticket224_reply_grant_test.go:283` 用的正是前缀，**不可复用**）。
  2. **钳位正控唯一可读面**＝`ra.gate.Queue().Timeout()`（导出：`gate.go:169 Queue()` ＋ `queue.go:126 Timeout()`；同形先例 `ticket84_no_owner_test.go:108`、`queue_test.go:48`）。`ra.gate.Window()` 对 `l1_window_sec` 敏感但被钳在 3s ⇒ **只能当"钳位仍在"的反控，不能当正控**。
  3. **载体先例已在同包**：`resident_task_source_246_windows_test.go:68-85` 已经会 `newResidentApproval()` ＋ `assembleRuntime(runSpec{gate/ui/cards/taskCtx})` 且 stderr 落进 buffer ⇒ 一枚"注入一行 `session <corr>` 后读 buffer"的用例形状接得上（不需跑进程、不需 winlive 档）。
- **量不到／判不动（具名）**：无控制台那一支（双击／Explorer）里 AC#2 的"不再出现"**机读不了**——不是尺不够，而是那句话在那一支**根本不产生**（`resident_task_source_windows.go:230 return nil` ⇒ 无 ledger、无 `runReplyLoop`、无 `session` 读者）。⇒ 那一支的"真落一行"＝〔待验，需跑一遍才知道；只读层面链路接不上〕。⛔ **没有**写成"做不到"。

## §6 `Confirming` 那一维：常驻腿持不持有状态机

- 常驻腿**不持有 Machine**（三处读满）：`grep -rn "statemachine.Machine|statemachine.New" --include=*.go cmd/wisp` 非测试命中**只有 `cmd/wisp/models.go:303`**（属 `wisp models` 链，不是常驻 boot 链）；`resident_approval_windows.go:54`(import)、`:492 SetState(StateSleeping)`、`:501-506 stateForCardLevel`（`"L1"→StateConfirming`，否则 `StateAwaitingApproval`，**输入是卡片自己的 `Level` 字符串，不查表**）；`resident_ball_windows.go:275` 的 `Initial: statemachine.StateSleeping` 是 `ball.Options.Initial`（`internal/ball/ball_windows.go:66`；`ball.New` 内只在 `:156-157` 填默认、`:263` 首帧渲染），**不是 `statemachine.New`**。
- 球不查表：`ball_windows.go:311-313 SetState` 只 `sta.PostTask(applyStateLocked)`——无 `Fire`、无转移合法性检查。
- ★ **票面 §7-4③ 那句"全仓产码只一枚调用者 `cmd/wisp/models.go:303`"字面不成立**：第二枚产码调用者＝**`cmd/balldebug/main.go:188`**（`//go:build windows`，`:1-8` 自述 ticket 07 调试台件）；测试档另 6 枚（`internal/ball/hotkey_live_test.go:319`、`interaction_live_test.go:54/134`、`internal/models/bridge_test.go:16/86`、`handoff_window_109_test.go:62`）。⇒ **实质读数（表那一半没有可破的东西）仍成立，被引范围（`cmd/wisp`）也成立，"全仓"二字是口径错。**
- **写码腿有没有先动过附近文件**：`git log --oneline 1d8106d8..HEAD -- cmd/wisp internal/agent/approval` ＝ 18 发；`git diff --stat` 的 36 枚文件里**不含 `resident_approval_windows.go`、不含 `run.go`、不含 `resident_task_source_windows.go`** ⇒ §0/§1/§3 引的那些行号**同形未漂**。
- **漂了的读数（引用前先重跑）**：`gate.go` 被 `ce693ce4`(220-r1) 插了 6 行 ⇒ 票面 §7-2 的 `g.grants` 两枚命中 `:667/:678` **今天＝`:673/:684`**〔编排者复跑对上新行号〕；§现量.2 的钳位 `gate.go:137-145` **今天＝`:143-151`**；`allowSession` 函数体今天 `:658-694`。常量未漂＝`queue.go:107`(300s)/`:116`(3s)/`:120`(2s)/`:122`(3s)。⚠ `ticket220_l1_window_read_test.go:4` 注释硬写的 `gate.go:302-316` 已被 220-r1 自己写漂（是注释不是尺 ⇒ 不红）。

## §7 排程约束（不裁，只报能不能／缺什么／撞什么）

- **只补 `[risk]` 那两枚字段（`Window`/`ApprovalTimeout`）＝能派。** 真实落点（按量到的，⛔ 不照抄票面）：`cmd/wisp/resident_approval_windows.go:109-113`（补两枚字段）＋ `:105` 的函数签名（要 `dataDir`；产码唯一调用点 `resident_windows.go:126`，`rt.Layout.DataDir` 在其之前已可用）＋ `:43-56` 的 import（加 `internal/config` 与 `path/filepath`）＋ 取值形状照 `resident_windows.go:177`/`:200` 的 per-use `config.LoadFile`。**写面＝1 枚文件（＋`resident_windows.go:126` 那一处调用）。**
  - **撞什么**：签名一改即撞 **13 枚测试调用点编译红**（P1 具名清单；枚数编排者复跑）。这些用例名**不许为落地去改**（票面禁区"不许为变绿放宽断言"）；但"补参数"不是放宽断言——**这 13 枚要不要跟着改调用形＝需要编排者裁的口径，枚数已给。**
  - **缺什么**：不缺读数。
- **补 `Grants` 那枚字段＝不能派（今天缺的不是勇气，是一枚值）。** 具名缺的接缝：构造时刻无 ledger 对象（`run.go:478` 在 `resident_windows.go:126` 之后约 129 行、且被 `resident_task_source_windows.go:230` 的条件 return 罩着）；`g.grants` 无晚绑定入口（`gate.go:160` 唯一写点，晚绑定三词全仓 0 命中）；`runSpec` 无 ledger 字段（`run.go:102-175` 逐枚）。要动它必须先答"第二枚 mint 还是新造 holder"——**两条都不是"补字段"，也不是"移动"**。
  - **另撞一枚最易忽略的硬钉**：任何让 `run.go` 第 424 行以后整体位移的改动（含给 `runSpec` 加字段）＝ P7 四处 evidence drift ⇒ 走 `run.go` 那侧的修法今天**必须连带重钉 255 的名册**，而那枚文件归票 255／248 AC#8，**不在本票写面**。
- **AC#2 的正控只有一把尺有效**＝`ra.gate.Queue().Timeout()`；`Window()` 被钳死 ⇒ 拿它当正控会读出恒绿假象。第一形的机读禁现＝可做，钉 `gate.go:675` 那句独有词组，⛔ 钉前缀。
- **与本腿无关但量到的排程事实**：ⓑ 那句话今天**还没落到任何产码里**——`grep -rn "只作用于跑任务的进程|常驻腿今天用常量"` 在 `cmd internal tools frontend` ＝ **0 命中**（只在 `.scratch/**` 与 `docs/**`）⇒ **`255-r2` 的 ⓑ 半格未交**；`config_readers_255.go` 的 `hotRowClaims` 里没有 `risk` 行 ⇒ ⓑ 若要"由登记表同源产出"，**那张表今天没有这一格**。
- **"我不替你裁"**：上面每一格都是盘上现读，无一处"应该没问题"。

## §8 我认为量不到的一切（穷举具名）

1. 无控制台形状里「真落一行」：链路接不上（`resident_task_source_windows.go:230`）⇒〔待验，需跑一遍才知道〕。
2. `GRANT-DROPPED` 第一形在常驻腿今天有没有真被印出来过：**无任何实跑凭据**；票 248-v1c §5-7 那句"没读到实跑一行"复认成立（盘上仍零用例）。
3. `Options` 那 7 枚零值字段"该不该补"：不是读数问题，没有量具。
4. winlive 档那 3 枚 `newResidentApproval()` 在本机此刻是否会被 CI 编到：tag 是 `windows && winlive`，"这台机器的 CI 带不带 winlive"在只读禁令内量不到。
5. `confirm_timeout_sec` 改了以后，常驻腿里**已经开着**的那张 L2 卡片会不会跟着变：读到的是"无 re-apply"，但"已挂起 item 用哪个值"要靠跑 ⇒ 〔待验〕。
6. 13 枚测试调用点的"编译红"是不是真红：按只读禁令没跑编译器；那是从"调用形无参"到"签名加参"的字面推断，红句措辞量不到。
7. 母票 248 AC#10 ⓑ 那句话最终落哪张表（`config_readers_255.go` 还是 `unwired.go` 的 `lockedKeyDisposition`）：两处都只是**没有**对应行，归口是裁定不是读数。
8. `docs/PLAN.md:1642` 原文：票面引它，本轮没核（属票面 AC#3 越界面，刻意不读不改）。

## §9 门禁读数（本件自量）

本轮**零 Go 命令**（腿自报，编排者复跑口径一致：只用 `grep`/`sed`/`Read`/`git`）。编排者代落本件的复跑清单见各节〔编排者复跑〕标注；复跑抓到两处**我自己的坏尺**：① `grep 'grants: o.Grants'` 因**空白对齐**零命中（真身在 `gate.go:160`，换 `[[:space:]]*` 才中）；② 腿写的 `run.go:474` 我量到 `:473`（漂一行）。
