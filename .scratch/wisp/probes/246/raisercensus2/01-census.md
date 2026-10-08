# 246-raisercensus-2 / 01 普查交件：问一逐跳链 ＋ 问二 `AC#7` 逐句

腿名：`246-raisercensus-2`（只读普查续程腿）。锚＝HEAD `cb9f1663`（取数钟 `2026-10-08 18:4x +0800`）。
本件**不裁任何东西**；只落静态尺读数与档位。⛔ 零 Go 命令（`181-v3` 种突变中）；⛔ 未改产码/票面/台账；⛔ 未种突变；⛔ 未翻框。
行号纪律：所有"文件:行号"与它所引的短语**同一发现取**（本仓今日已出三次"行号与短语错位"，见 `A729 §2`）。
调用形状尺＝`git grep -nE '[A-Za-z0-9_]\.<方法>\(' HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`；赋值/裸形状尺＝`git grep -nE '\.<方法>([^A-Za-z0-9_(]|$)'` 同射程。

---

## 1. 问一·原料：`approval.Gate` 导出方法名册（`internal/agent/approval/gate.go`＝**793 行**）

名册现取尺＝`git show HEAD:internal/agent/approval/gate.go | grep -nE '^func '`。导出方法 **14 枚**（`New` 是构造函数单列）：

| # | 方法 | 定义行 | 调用形状尺（产码、非测试）命中 | 档位判定 |
|---|---|---|---|---|
| 1 | `Queue()` | `:169` | 3 枚：`approval_reply.go:377`、`panel_pump.go:62`、`replies.go:311` | 〔已接〕（读面访问器） |
| 2 | `Channels()` | `:173` | 3 枚：`approval_reply.go:512`、`resident_approval_windows.go:567`/`:787` | 〔已接〕（Esc 通道登记：置 true/置 false） |
| 3 | `Window()` | `:176` | 1 枚产码：`resident_approval_windows.go:384`（同一行的 `Queue().Timeout()` 属 #1） | 〔已接〕 |
| 4 | `AdmitTextTask()` | `:194` | 5 枚：`approval_always.go:113`、`config_reload.go:272`、`resident_approval_windows.go:666`、`run.go:842`、`run.go:955` | 〔已接〕（D47 准入前置） |
| 5 | **`PendingWindow()`** | `:252` | **2 枚**：`internal/tools/bridge.go:443`（承载路径）、`resident_approval_windows.go:678`（只活在备用接缝内，见 §2 SeamX） | 〔已接〕（经 bridge；seam 那枚不可达） |
| 6 | `Complete()` | `:389` | 产码直呼 **0**；转发面＝`internal/agent/approval/report.go:91` `func (b bus) Complete(corr string) { b.g.Complete(corr) }`（⚠ `internal/tools/bridge.go:532` 的 `b.cancel.Complete` 是 cancel 总线，非本方法） | 〔建了但没接〕（产码直呼 0；经 bus 转发面存在，bus 的产码使用者本腿未逐枚展开） |
| 7 | `Veto()` | `:421` | 包内 1 枚：`replies.go:463` `g.Veto(Veto{...})`；cmd/wisp 两处是其他类型（`approval_reply.go:363` `s.live.h.Veto`、`resident_approval_windows.go:612` `ra.cards.Veto`）——它们经 Replies 的面回到包内 | 〔已接（包内）〕（Esc 否决链见 §2 H-Veto） |
| 8 | `LateVeto()` | `:482` | 2 枚 bus 转发：`report.go:78`、`report.go:185` | 〔经 bus〕 |
| 9 | **`PendingApproval()`** | `:502` | **6 枚**：`internal/tools/bridge.go:458`（承载路径）、`approval_always.go:115`、`config_reload.go:279`、`run.go:844`、`resident_approval_windows.go:680`（seam，不可达）、包内 `gate.go:270`（R7 升级：L1 窗口带 ≥50 路径时内部升级到 L2 路线） | 〔已接〕 |
| 10 | `Native()` | `:622` | 1 枚包内：`replies.go:363` `g.Native().AllowSession(...)`（native 答案路线） | 〔已接（包内）〕 |
| 11 | `Panel()` | `:626` | 2 枚包内：`replies.go:478` `g.Panel().Head()`、`:487` `g.Panel().View(corr)`（面板读面） | 〔已接（包内）〕 |
| 12 | `DecideFromNative()` | `:724` | 2 枚包内：`replies.go:328`、`:413` | 〔已接（包内）〕 |
| 13 | `DecideFromPanel()` | `:736` | 2 枚包内：`replies.go:408`、`:437` | 〔已接（包内）〕 |
| 14 | `Replay()` | `:755` | 调用形状尺 **rc=1 零命中**；裸形状尺 **rc=1 零命中**（均排 `.scratch` 与 `_test.go`） | 〔建了但没接〕（产码零调用者；测试侧未量） |

**赋值形状尺（不带括号那一族）逐枚结论**：14 枚的裸形状命中**全是注释或同名的其他字段**——`internal/tools/gate.go:117`/`:120` 的 `g.Window` 是 tools 自己的 `GateFuncs.Window` 函数值字段、`run.go:1368`/`:1374` 与 `resident_approval_windows.go:884` 是 doctor/报告结构体的 `Channels`/`Window` 字段、`gate.go:139`/`:143` 与 `replies.go:181` 同为字段拷贝。⇒ **`approval.Gate` 的方法作为函数值被赋走的产码点＝0 枚**。
**把两枚回调适配成 tools.Gate 的面＝`internal/tools/gate.go:110` `GateFuncs`**：使用者仅 `internal/tools/task_cancel_221_legs_test.go:63`（**测试件**）⇒ 产码使用者 0；产码走的是 `approval.Gate` 本体直插 bridge（见 §3.2）。

**判定入口的准确定义**：L1/L2 两条路线＝`PendingWindow`（`gate.go:252`）与 `PendingApproval`（`gate.go:502`）；两枚都收 `tools.Decision`、回 `(tools.Answer, string)`，与 `internal/tools/gate.go` 里 `Gate` 接口的两枚回调同名同形。`g.ui.Prompt` 在全包（排测试）**恰 2 处**：`gate.go:294`（PendingWindow 内）、`gate.go:536`（PendingApproval 内）——尺＝`git grep -nE 'ui\.Prompt\(' HEAD -- 'internal/agent/approval' ':(exclude)*_test.go'`，**只此两枚**。

---

## 2. 问一·逐跳链：从"任务真跑起来"到"会走到 `g.ui.Prompt`"

| 跳 | 锚（同一发现取） | 产码调用点 | 档位 | 条件/备注 |
|---|---|---|---|---|
| H1 双击进程 → 常驻腿 | `cmd/wisp/main.go:66` `runResident()`；`resident_windows.go:33` `func runResident()` | 1 枚 | 〔已接〕 | 无参启动＝常驻；`wisp run` 是另一条腿 |
| H2 造门（常驻那枚） | `resident_windows.go:132` `ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`；`resident_approval_windows.go:365` `ra.ui = &ballCardUI{ra: ra}`；`:369-376` `ra.gate = approval.New(approval.Options{`，`:370` `UI: ra.ui,` | 1 枚（另一枚是 `run.go:612` 的 CLI 自造门 posture） | 〔已接〕 | 门与它的展示面在同一处造；`ra.ui` 即举卡真身 |
| H3 起任务源 | `resident_windows.go:260` `src := startResidentTaskSource(rt, ra)`；func 在 `resident_task_source_windows.go:228` | **1 枚**（⚠ 另有 2 枚是 `.scratch/wisp/probes/258/v1/mut0|mut3` 的突变副本，非产码） | 〔已接〕 | `A722 §3ⓐ` 引 `:261`，本次现量落 **`:260`**（具名报回：1 行差） |
| H4 任务源装配管线 | `rts:278` `run, code := assembleRuntime(spec)`；`assembleRuntime` 定义＝`run.go:382`；`runSpec` 类型＝`run.go:102` | 1 枚 | 〔已接〕 | spec 字面量逐枚见 §3.3 |
| H5 装配根把门注入管线 | `run.go:604-606`（注入 posture：`rt.gate = s.gate` 在 `:606`）；`run.go:745-749` `rt.bridge = tools.New(tools.Options{` 其中 `:748` `Gate: rt.gate,`、`:749` `Cancel: rt.gate.ToolsCancelBus()` | 1 枚（bridge 的产码装配点） | 〔已接〕 | 「一个进程只许一枚 approval.Gate」（`run.go:583-600` 注释逐字）；`rt.bridge` 交 loop＝`run.go:997` `Tools: rt.bridge`（另有 `:782` `ParentTools: rt.bridge`） |
| H6 入口触发任务 | 控制台：`rts:231` `console := interactiveStdin()`；`:348` spawn `src.runConsoleLoop`；`rts:374` func；`:387` `case "task":` → `:392` `src.submitTask(strings.TrimSpace(rest))`。测试位：`:137` `const testTaskTextEnv = "WISP_TEST_TASK_TEXT"`；`:351` `src.submitTask(injectedText)`。公用消费者：`rts:430` `submitTask` → `src.run.execute(root.Ctx, text)`（loop 的单发执行） | 2 枚入口（控制台＝生产形；注入位＝测试形） | 〔已接〕 | 两条入口都汇入同一枚 `submitTask`；无入口时打 `taskEntryDisabledClaim`（`:98-101`）并 `return nil`（`:236-243`） |
| H7 任务→工具→问门 | `Bridge.Execute`（`internal/tools/bridge.go:268`）→ `:352` `ok2, why := b.route(ctx, &dec, sil, grantID)` → `route`（`:423`）内读 `Gate` 字段：`:443` `b.gate.PendingWindow(ctx, *dec)`（L1 支）、`:458` `b.gate.PendingApproval(ctx, *dec)`（L2 支） | 2 枚读点 | 〔已接〕 | 条件＝**有效级**（评估级经 permission mode 筛选后的 `sil.Level`）：L0 放行、Deny 直拒均**不问**；L1 且 `grantID==0` 问 `PendingWindow`（有 D45 会话授权则 `allow_grant` 不问）；L2 问 `PendingApproval`；其余 fail-closed |
| H8 gate→UI | `gate.go:294`（PendingWindow 内）与 `gate.go:536`（PendingApproval 内）`if err := g.ui.Prompt(ctx, p); err != nil` | 2 枚（全包仅此两枚） | 〔已接〕 | 两处都是"宿主不可达 ⇒ fail-closed 拒绝"；Prompt 结构由 `promptFor`（`:597`）装配 |
| H9 UI→举卡 | 常驻那枚 `g.ui`＝`ballCardUI`（H2）；`resident_approval_windows.go:864` `func (u *ballCardUI) Prompt(...)` 举卡真身；`:877` `b.TakeEscForCancel()`（仅 L1 借 Esc）；`:929` `Update` 回收；`:947` `settleOrb`；`:955` `b.ReleaseEscAfterSession()` | 1 枚 | 〔已接〕 | 无球时 `errResidentNoBall` fail-closed（`:866-871`） |
| SeamX 备用接缝（宿主自发起卡） | `askConfirmation`（`:641`，内 `:666` `AdmitTextTask`、`:678`/`:680` 两路线）← 唯一调用者 `AskOnTaskRoot`（`:697`，`:700` 调 `askConfirmation`） | **0 枚**（形状尺 `git grep -nE '[A-Za-z0-9_]\.AskOnTaskRoot\(' …`＝rc=1；全形状命中全为注释/定义：`resident_approval_windows.go:688`/`:697`、`rts:14`、`resident_windows.go:249`） | 〔建了但没接〕 | `K3` 判"保留不删"；`resident_windows.go:249` 注释逐字「This call is the caller」指 H3 那枚调用，**不是**指 `AskOnTaskRoot` 被接上 |
| H-Veto Esc 否决链（四步的第 3 步，供 §4 用） | L1 借键＝`:877`；登记通道 `:567` `ra.gate.Channels().SetLoaded(approval.ChannelEsc, true)`；归还＝`:787` 置 false；否决落点＝`resident_approval_windows.go:612` `ra.cards.Veto(card.CorrelationID)` → 包内 `replies.go:463` `g.Veto(...)` → `gate.go:421` `Veto` | 各 1 枚 | 〔已接〕 | D38(e) 第 3 步钩子＝`resident_windows.go:240` `rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)` |

---

## 3. 问一·三小问分栏答案

### 3.1 哪个包在什么条件下调用判定入口
- **`internal/tools`**：唯一承载路径。`Bridge.Execute`（`:268`）→ `route`（`:423`）→ L1 `PendingWindow`（`:443`）/ L2 `PendingApproval`（`:458`）。条件＝上表 H7 那条 `switch sil.Level`（源码逐字在 `bridge.go:423-478`）。
- **`cmd/wisp`**：4 枚产码直呼 + 1 枚不可达：
  - `run.go:836-859` `confirmModeSwitch`：切到"会移除询问"的档位 → L2 `:844`（先 `:842` 准入）；
  - `approval_always.go:110-125` `runWidening`：「一直允许」落盘前的第二重 L2 `:115`（先 `:113` 准入）；
  - `config_reload.go:262-284`：config.toml 被手改放宽 → 重确认 L2 `:279`（先 `:272` 准入）；
  - `run.go:954-961` `admitTask`：只做准入 `:955`（不直接问门），loop 每发任务走它；
  - `resident_approval_windows.go:641-681` `askConfirmation`（L1/L2 两路线）＝**不可达 seam**（SeamX）。
- **`internal/agent/approval` 包内**：`gate.go:270` R7 升级（L1→L2）；答案面 `replies.go:328`/`:413`（Native）、`:408`/`:437`（Panel）、`:463`（Veto）。
- `Replay`（`:755`）零调用者 ⇒ 〔建了但没接〕。

### 3.2 工具侧（`internal/tools/`）在什么条件下问门
- `Gate` 字段三处：`bridge.go:50` `Gate Gate`（Options 字段；同块 `:49` 逐字「nil is NoGate (fail closed)」）→ `:135` Bridge 结构体 `gate Gate` → 构造函数 `:176` `gate := o.Gate`（nil 填 `NoGate{}`）。
- 读点＝`:443`/`:458`（§2 H7），**读后走**两支 switch：L1 支把 `AnswerAllow/AnswerTimeout` 都当放行（窗口超时＝执行）、`AnswerVeto` 落 `DecisionReject`；L2 支 `AnswerAllow` 放行、`AnswerTimeout` 落 `DecisionTimeout`（自动拒绝）、其余拒绝。
- **有没有产码调用点＝有**：2 读点在生产函数 `route` 内、`route` 被 `Execute`（`:268`）在 `:352` 调用、`Execute` 经 `Tools: rt.bridge`（`run.go:997`）交给 loop；bridge 的产码装配点＝`run.go:745-749`（`Gate: rt.gate` `:748`）。

### 3.3 常驻任务源那一跳（`cmd/wisp/resident_task_source_windows.go`，557 行）
- **交给管线的 `runSpec` 字面量逐枚列**（`:265-276`，类型在 `run.go:102`）：
  ① `stdout:  os.Stdout`（`:266`）② `stderr:  os.Stderr`（`:267`）③ `dataDir: rt.Layout.DataDir`（`:268`）④ **`gate: ra.gate`（`:269`）** ⑤ **`ui: ra.ui`（`:270`）** ⑥ `cards: ra.cards`（`:271`）⑦ `taskCtx: ra.root`（`:272`，注释 `:52` 逐字「THE TASK ROOT (ruling 2.3)」）⑧ `notify: postSystemNotification`（`:275`，D10 通知路径）——共 **8 枚**；随后 `:278` `assembleRuntime(spec)`。
- **由谁触发**：`startResidentTaskSource` 的**产码调用者 1 枚**＝`resident_windows.go:260`（在 `runResident()` 内、`bindBallHost` 之后、boot 报告之前）。触发链＝`main.go:66`（无参启动）→ `runResident` → `:260`。任务入口两形：控制台（`interactiveStdin()`≠nil 时 spawn `runConsoleLoop`）或 `WISP_TEST_TASK_TEXT` 注入位；**两形皆无 ⇒ 打「任务入口未启用」claim 并 return nil**（`:236-243`）；装配失败 ⇒ 打 `pipelineAbsentClaim`（`:280-290`）。

---

## 4. 问二·`AC#7` 完成判据逐句核对（判据原文票面自取）

票面现量尺：`git show HEAD:<票面> | wc -l` ＝ **97**；`grep -c '^\s*- \[x\]'` ＝ **8**、`- \[ \]` ＝ **1**（与 `A488`/`A722`/`A723` 同数，枚数未动）。
判据原文（票面 `:60` 该格「完成判据＝」之后的**逐字**）：

> 完成判据＝常驻那条腿真起一条任务管线（或经装配根注入一条最小任务源），一发真机走通"起任务 → 举卡 → Esc 否决 → 任务被取消"四步；⛔ 不许用测试构造的任务源冒充。

（同行前置的问题陈述逐字留档：「门在、卡进得来，**还没有任务会去举它**」；末句「**这是票 246 标题那句"永久之家"真正落地的地方，属后端主链，排队列最前**」＝排程语，非判据。）

**逐句判定（每句配我自己跑的尺）：**

1. **「常驻那条腿真起一条任务管线（或经装配根注入一条最小任务源）」⇒ 静态＝成立〔已接〕。**
   尺：§2 H3（`resident_windows.go:260` 唯一产码调用者）＋ H4（`rts:278` `assembleRuntime(spec)`，`:269` 注入 `ra.gate`）＋ H5（`run.go:606` `rt.gate = s.gate`，即"经装配根注入"的乙形）。⚠ 条件与降级路径并记：`interactiveStdin()` 为 nil 且无测试注入时**这条管线不起**、启动只打 claim（`rts:236-243`）——今天那个二进制是 CUI、控制台形成立（`K5` 的三把记录，本腿未复量 `build/wisp.exe`）。
2. **「一发真机走通"起任务 → 举卡 → Esc 否决 → 任务被取消"四步」⇒ 静态四跳全接〔已接〕；"一发真机"的行为凭据＝本腿量不到新的那一天读数。**
   静态尺：起任务＝H6；举卡＝H7→H8→H9（`bridge:352→:443/:458 → gate.go:294/:536 → :864`）；Esc 否决＝H-Veto（`:877`/`:612`/`replies:463`/`gate:421`）；任务被取消＝D38(e) 第 3 步钩子已在位（`resident_windows.go:240`）。行为侧：盘上既有凭据＝`resident_task_source_live_246_windows_test.go`（**589 行**，`//go:build windows && winlive`）里 `TestLive246ResidentPipelineRaisesACardAndEscVetoesIt`（`:95`，头注逐字「AC#7's four steps on a shipped resident process: task -> card -> Esc veto -> the call cancelled」）＋ `TestLive246ExitCancelsARunningTask`（`:256`）；`A488` 记〔接缝注入〕6 发全过、**〔真模型〕栏本机为空**；`A725 ⓑ` 判「今天会不会真举出一张」**仍无凭据**。⇒ 本腿零 Go 命令不能复跑 winlive，"满足"的行为断言无新凭据；这正是 `A722 §3` 待人裁那一问的前提之一（⛔ 我不裁）。
3. **「⛔ 不许用测试构造的任务源冒充」⇒ 静态＝无冒充证据，成立。**
   尺：生产入口＝控制台形（`rts:231`/`:348`/`:374`，`:387` `case "task":`）；测试位＝`WISP_TEST_TASK_TEXT`（`:137`），且注入形**启动时必自报**（`:253-258` 逐字含「这是只在本机测试台件里受理的一发任务文本…」）；winlive 用例还专门断言"走了注入形却报告没自报＝红"（用例 `:139-143`）。⇒ 两形不自混；两形**汇入同一枚产码消费者** `submitTask`（`:430`，`A488` MUT-4 已验删消费者即两枚真机用例死＝消费者承重）。⚠ 具名：历史"四步"证据的任务文本走的正是注入位（用例 `:32` 逐字「The task text arrives through WISP_TEST_TASK_TEXT, the narrow injection ruling」）——它算不算"测试构造的任务源冒充"，＝待人裁那一问的另一半（⛔ 我不裁）。

### 4.9 「起管线 vs 真举卡」有没有人裁过 ⇒ **没有。只有登记与收窄，没有裁定。**
- `A722 §3`（`:14168` 起）：「这一格要的是**一枚能跑突变的非实现者裁决**…排在 253-r1 退出之后，与 `259-v1` 同批排队」＝**登记待人裁**；
- `A724 §3`（`:14221` 起）逐字：「`A722 §3` 是**登记**…**不是裁定** ⇒ …这一问**今天仍然未裁**…**原样挂着**」；
- `A725`（`:14250` 起）：编排者 18:1x 自取数只做**收窄**（「⛔ 不撤、⛔ 不翻框」）；
- `A729 §3`（`:14306` 起）：续派本腿＝「把"常驻今天会不会真举出一张卡"这条链逐跳跃量清」。
⇒ 本腿按派单写死的停手条件**继续跑完**（命中"已登记"不停手），现把逐跳读数交回；「起管线算不算满足 `AC#7`」仍待裁决腿。

---

## 5. 我没量到的（具名）与简报转述差异（具名报回）

**未量**：
- winlive 真机用例未跑（零 Go 命令约束），"真机四步"无今天的新行为读数；
- `interactiveStdin()` 定义行未复量（其出处引 `rts:20` 注释指向 `approval_reply_stdin_windows.go:41`）；
- `Replies.Veto`（`ra.cards.Veto` → 包内 `replies.go:463`）的中段函数归属未逐行读；
- `Replay` 的**测试侧**调用未量（§1 的 rc=1 是排测试尺）；
- `bus`（`report.go`）的上游产码使用者未逐枚展开（只量到转发面 :78/:91/:185）；
- `assembleRuntime` 在**非常驻腿**（`wisp run`）的调用枚数未展开（本腿射程＝常驻链）。

**与派单/前期转述的差异（以盘上现量为准）**：
- `startResidentTaskSource` 调用者：`A722 §3ⓐ` 写 `resident_windows.go:261`；本次现量 **`:260`**（同句同发现取）。
- 票面 `AC#7` 引「`newResidentApproval()` 有产码调用者（`resident_windows.go:122`）」：今天产码用的是 **`newResidentApprovalWithConfig(rt.Layout.DataDir)`**（`resident_windows.go:132`）；无参 `newResidentApproval()` 的**产码调用者＝0**（`git grep -nE 'newResidentApproval\(\)' HEAD -- cmd/wisp` 命中除定义 `:299-300` 外全在 `*_test.go`）⇒ 那枚旧引已随票 256 过期。
- `resident_windows.go` 那句 `This call is the caller` 落 **`:249`**（派单背景写 `:248`；`A729 §2` 已纠，我复量同 `:249`）。
