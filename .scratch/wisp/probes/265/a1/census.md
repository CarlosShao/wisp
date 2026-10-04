# 普查件 265-a1 — 常驻（GUI）腿的审批门没有会话授权记账位：现量复认与三形代价

**腿**：`265-a1`（只读普查，AC#0）
**取数时刻**：2026-10-04，锚点 HEAD `35e852f9f44b6763f211a714ffe8c4ac303b9fa8`（`dev`）
**对象票**：`.scratch/wisp/issues/265-resident-gate-has-no-session-grant-writer.md`（提交号 `35e852f9`）

## 0. 本腿的运行约束（先声明，免得后面读数被误读）

- ⛔ **本腿一枚 go 命令都没跑**（禁 `go build`/`vet`/`test`/`run`，禁 `scripts/d22scan.sh`，禁 `scripts/slo-check.ps1`）——在飞的 `260-v1` 独占 `cmd/wisp`＋`internal/agent/approval`＋`internal/ball` 的取数期。
  ⇒ 本件里所有"接得上吗／会不会红"的答案**都是读码级**，凡需真跑才知道的一律标〔仅读码，未跑〕。
- ⛔ 本腿**不裁形**：§3 只摆料，ⓐ／ⓑ／ⓒ 的选择归编排者。
- ⛔ 本腿**没读也没引用** `.scratch/wisp/probes/260/v1/**`。
- ⚠ **盘脏读数（本腿第一枚读数，影响全部行号引用）**：`git status --porcelain -- cmd internal` 取数时＝
  `M cmd/wisp/resident_approval_windows.go`（唯一一枚脏文件，`260-v1` 的突变体：`residentCancelKeySpelling` 里读热键报告的那个 `for` 环被替换成 `_ = b`，产码净行数 **−6**）。
  ⇒ **本件所有 `file:line` 一律以 `git show HEAD:<path>` 的行号为准**（把该文件导出为临时件后逐行编号）。
  `cmd/wisp/run.go` 与 `internal/agent/approval/gate.go` **盘上与 HEAD 逐字节相同**（已核：`git diff --stat HEAD -- cmd/wisp/run.go` ＝空；`gate.go` 盘上 793 行＝HEAD 793 行，diff 空）。
  ⛔ 本腿不动那枚脏文件（只建不删的临时件都在 `/tmp`，见 §7）。

## 1. 现量复认（对编排者 13:2x 那四条逐条判：对／错＋真读数）

**判语只用三词：对／错／量不到。** 尺全部是 `git grep` / `git show HEAD:` / `sed`，一枚 Go 命令都没用。
除注明外，**行号一律是 `git show HEAD:<path>` 的行号**（盘上 `resident_approval_windows.go` 此刻被 `260-v1` 的突变体污染，见 §0）。

### 1.1 §现量 1（接缝本体）＝**对**，八处行号逐处复认命中

尺＝`grep -n "grants" internal/agent/approval/gate.go`（该文件盘上＝HEAD，`git diff --stat HEAD -- ...gate.go` 空）＋ `sed -n '640,710p'`。

| 编排者说 | 我读到的 | 判定 |
|---|---|---|
| `:160` 是 `grants:   o.Grants,` ＝全仓唯一写入点 | `gate.go:160` 逐字 `grants:   o.Grants,`；全仓 `g.grants` 写点尺＝`git grep -n "\.grants =" HEAD -- cmd internal tools` → **命中 `queue.go:563`（`fresh.grants = newGrantStore()`，那是 `Queue.grants *grantStore`，一次性 nonce 表，**与 `Gate.grants` 同名不同物**）**＋测试两枚；`Gate.grants` 的写入**只有 `gate.go:160` 那一枚构造期赋值** | **对**（⚠ 补一条：写"唯一"时必须带射程，`Queue.grants` 会被下一枚程读成第二枚写手） |
| `:673` `if g.grants == nil` | 逐字命中 | 对 |
| `:674-676` 落 `GRANT-DROPPED corr=%s tool=%s paths=%d (本机没有接入会话授权记账，本次按「仅本次」放行，没有落盘任何规则)` | 逐字命中（句子跨 `:674`＋`:675` 两行字面量，`)` 与参数在 `:676`） | 对 |
| `:679` `if tool == ""` → `:680` 第二句 | 逐字命中 | 对 |
| `:684` 才 `g.grants.Record(ctx, tool, p)` | 逐字命中 | 对 |
| `:689-690` `GRANT-RECORD-FAILED` / `:693` `GRANT-RECORDED` | 逐字命中 | 对 |

函数体边界：`allowSession` ＝ `gate.go:658-696`；`nativeAPI.AllowSession` ＝ `:653-655`（**纯转发，不是接缝**——引用接缝时别指 `:653`）。

### 1.2 §现量 2（生产建门点枚数＋那句注释）＝**对，但"两枚"这句要带射程才不被读错**

票面尺我原样重跑（`git grep -n "approval.New(approval.Options{" HEAD -- cmd internal | grep -v _test`）⇒ **2 枚**：
`cmd/wisp/resident_approval_windows.go:219`、`cmd/wisp/run.go:612`。枚数复认。

**本腿自拉的两把更宽的尺**（票面那句"不许照抄两枚"的兑现）：
- 尺 A＝`git grep -n "approval\.New(" HEAD -- cmd internal tools` 去掉 `_test.go` ⇒ 产码命中 **3 行**：上述两枚 ＋ `internal/agent/approval/doc.go:28`（**注释，非建门**）。⇒ 产码建门仍 **2 枚**。
- 尺 B＝`git ls-files cmd` 逐目录核（`cmd/balldebug`／`cmd/llmrecord` 全查）⇒ `git grep -n "approval\." HEAD -- cmd/balldebug` ＝ **0 命中** ⇒ 调试台件不建门。
- 尺 C＝补同包构造形（`256-a2 census §4.1` 具名说票面那把尺"漏同包构造点"）：`New(Options{` 在本仓**产码零额外命中**（同包自建的只有 `approval` 与 `tools` 的**测试**）。⇒ 票面"两枚"成立。

注释那句：`resident_approval_windows.go:186` 逐字 `// Options.Grants stays unset here. The session ledger's only production`——**仍是真话**（现读 `:219-225` 字面量实传 5 枚字段：`UI`/`Channels`/`Window`/`ApprovalTimeout`/`Logf`，无 `Grants`）。⚠ 同段 `:188` 那句 "this file's gate is built roughly 129 lines BEFORE" 的 **129 已漂**：真身＝`resident_windows.go:132`（建门）↔ `resident_windows.go:260`（`src := startResidentTaskSource(rt, ra)`）＝**128 行**，且两行之间还隔着 `resident_task_source_windows.go:224-230` 那枚条件 return（见 §2.2）。

### 1.3 §现量 3（既有尺两枚，极性相反）＝**对**，并补第三枚

票面尺 `git grep -c "GRANT-DROPPED" HEAD -- '*_test.go'` ⇒ **两枚文件**：`cmd/wisp/ticket224_assembly_test.go`＝1、`internal/agent/approval/ticket224_reply_grant_test.go`＝2（合计 3 处命中）。逐枚读断言后：

| 枚 | `file:line` | 它今天**逐字**断什么 | 极性 |
|---|---|---|---|
| 尺① | `internal/agent/approval/ticket224_reply_grant_test.go:283` | `if !f.log.has("approval: GRANT-DROPPED") { t.Errorf(...) }` ⇒ **断这句在场** | 正向（要它出现） |
| 尺①补 | 同文件 `:287` | `if f.log.has("GRANT-RECORDED") \|\| f.log.has("grant_id=")` ⇒ 断"没记账位就不许声称记了" | 反向 |
| 尺② | `cmd/wisp/ticket224_assembly_test.go:194` | `if strings.Contains(h.err.String(), "GRANT-DROPPED") { t.Errorf("audit says the scope was DROPPED on a host that does have a ledger") }` ⇒ **断这句不在场** | 反向（不许出现） |
| **尺③（票面没列）** | `cmd/wisp/resident_approval_risk_256_windows_test.go:452-456` | `if got["Grants"] { t.Errorf("the resident leg now passes Options.Grants. …") }` ⇒ **断常驻腿的 `Options` 字面量里没有 `Grants`**，读的是**语法树不是日志** | 反向，且**唯一一枚直接钉住本票那一格** |

⇒ 编排者那句"这两枚极性相反"读数成立；**本腿的增量是第三枚**：它不读日志、读 AST，所以 ⓐ 一落地**它必红**，而那两枚日志尺反而都不红（见 §4）。

### 1.4 §现量 4 那两处待复认：① **错（我票面那句不成立）**／② **对**

**① "resident_windows.go prints the same limit at boot"＝不成立。** 尺＝把 `resident_windows.go`（HEAD，288 行）**全部 16 处打印/记录语句逐枚列出**（`git grep -nE "fmt\.(Print|Fprint)|slog\."`）＋ 在该文件里 `grep -niE "grant|会话|记账|ledger"` ⇒ 只命中 **2 行注释**（`:113`/`:137`，都是 `ledger A481` 那种台账引用，**不是会话账本、也不是打印**）。boot 那一句（`resident_windows.go:269-270`）逐字拼的是：
`"wisp: %s; 任务来源：%s; 面板：%s; D38(e) steps with an owner in this process: %s\n"` ＋ `ra.residentStatusLine()`（`resident_approval_windows.go:589-598`，逐字只有"审批门已装配进本进程（取消通道：X 已加载；等待中的确认项：N）"两形）＋ `src.taskPosture()`（四形常量在 `resident_task_source_windows.go:107-110`）＋ `panel.statusLine()`。
**这四形没有一枚含"会话授权／记账／Grants"字样。**
⇒ **更正：`cmd/wisp/run.go:597` 那句话（逐字 "resident_windows.go prints the same" ＋ `:598` "limit at boot;"）在今天的盘上指不到任何产码句子**，属**注释里编造了一条不存在的可见性**，不是"改了以后不印了"。
⚠ 归因（`git log -S"会话" -- cmd/wisp/resident_windows.go`＝**0 命中**，`git log -G"prints the same limit" -- cmd/wisp/run.go`＝**0 命中**，两句合起来读＝**这句话从来没有对应的产码打印**）：它自诞生就是假话，不是过期。

**② `resident_approval_windows.go:186` "Options.Grants stays unset here"＝仍是真话。** 见 §1.2。

## 2. 四问（票面 AC#0 原文）

### 2.1 ① 常驻腿的门被谁持有、答复从哪几处进来

**持有者**：门赋给 `residentApproval.gate`（字段声明 `resident_approval_windows.go:84`，`gate *approval.Gate`），
赋值那一枚在 `:219`（`ra.gate = approval.New(approval.Options{…})`，字面量到 `:225` 闭合）。
`residentApproval` 这枚对象今天**唯一**的产码构造者＝`newResidentApprovalWithConfig`（`:213`），
它唯一的产码调用点＝`resident_windows.go:132`（`ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`，票 256 §8.2 那形ⓐ留下的那枚）。
姊妹签名 `newResidentApproval()`（`:170`）**产码调用者 0 枚**，只被本包测试领（尺＝`git grep -n "newResidentApproval()" HEAD -- cmd internal | grep -v _test`＝空）。

**生命周期谁管**（逐段，全在 `resident_windows.go` 的 `runResident` 体内）：
| 时刻 | `file:line` | 做了什么 |
|---|---|---|
| 建 | `resident_windows.go:132` → `resident_approval_windows.go:219` | 门在这一刻存在（早于面板 `:151`、早于球 `:217`、早于任务源 `:260`） |
| 交给球当否决执行者 | `resident_windows.go:217`（`startResidentBall(rt.Registry, ra.vetoByEsc, …)`） | 门不进门，交的是**方法值** |
| bind | `resident_windows.go:228`（`ra.bindBallHost(rb)`）→ `:360-361` | 挂球 UI ＋ `Channels().SetLoaded(ChannelEsc, true)` |
| 注入装配根 | `resident_windows.go:260`（`startResidentTaskSource(rt, ra)`）→ `resident_task_source_windows.go:252-264`（`runSpec{gate: ra.gate, ui: ra.ui, cards: ra.cards, taskCtx: ra.root}`）→ `:265`（`assembleRuntime(spec)`） | 门以**指针值**递进 `cmd/wisp/run.go` |
| 封闭 | `resident_windows.go:227`（`defer ra.detachBall()`）＋ `:240`（`RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)`） | 退出时卸通道 ＋ 拒待批卡 |

**答复入口逐枚**（六枚候选全查；每枚给 `file:line` 并答"是不是最终走 `gate.go:673` 那一支"）：

| # | 入口 | `file:line`（产码可达链） | 能不能触到 `gate.go:673` |
|---|---|---|---|
| **1** | **控制台动词 `session <编号>`**（唯一一枚真能触到的） | `resident_task_source_windows.go:345 runConsoleLoop` → `:376 default` → `:381 src.surface.handle(...)` → `approval_reply.go:562 handle` → `:566-567 case "session": return s.session(corr)` → `approval_reply.go:257-259 s.live.h.AllowSession(s.ctx, corr)` → `internal/agent/approval/replies.go:351-363 g.Native().AllowSession(...)` → `gate.go:653-654` → `gate.go:658 allowSession` → **`:673 if g.grants == nil` 命中** | **是**，且**全仓只有这一条路**（尺＝`git grep -n "AllowSession" HEAD -- cmd internal | grep -v _test`：产码调用者只有 `approval_reply.go:259` 一枚） |
| **2** | 球上的取消键（Esc／配置的组合键） | `resident_windows.go:217` 交 `ra.vetoByEsc` → `resident_approval_windows.go:398-413 vetoByEsc` → `:406 ra.cards.Veto(...)` → `approval/replies.go:463 g.Veto(...)` | **否**——veto 一支不进 `allowSession`，只否决 L1 窗口 |
| **3** | 面板入向（`approval.decide`） | `cmd/wisp/panel_inbound.go:277 Message: nil`（产码路由没有审批消息位）；`gate.go:736-746 DecideFromPanel` 对任何 allow **按路线直接拒绝** | **否**——面板侧连 `Allow` 都没有（`approval/ui.go` PanelAPI 无 Allow；`replies.go:347-350` 逐字"There is deliberately no PanelAllowSession"）。⇒ **票面上"用户在那张卡片上**点**「本会话内允许」"这个动作今天不存在**，见 §2.5 |
| **4** | 托盘 | `internal/ball/tray_windows.go:86-91`：菜单只有 打开面板／静音／暂停唤醒／退出，`Allow`／`允许` **零命中** | **否**（入口不存在） |
| **5** | 原生路由 `allowSession` 本身 | `gate.go:658`——它是**被到达的**，不是入口；上游只有 #1 | 它是 `:673` 的**宿主函数** |
| **6** | 宿主自发的卡（模式切换 L2） | `run.go:826-844 confirmModeSwitch`（走 `rt.gate`，注入形状下＝同一枚 `ra.gate`）；`resident_approval_windows.go:491 AskOnTaskRoot` | **间接触到**：卡能建、能显示，但**答复仍只能从 #1 进来**。⚠ `AskOnTaskRoot`／`askConfirmation` 的产码调用者 **0 枚**（尺＝`git grep -n "AskOnTaskRoot\|askConfirmation" HEAD -- cmd internal | grep -v _test`＝只有定义处与三处注释），票 256 §7-3 记过这一格，本腿复认**仍成立** |

⇒ **#1 的前提是"这一发常驻进程有交互控制台"**：`resident_task_source_windows.go:218 console := interactiveStdin()` → `:224-230` 没有控制台就 `return nil`（连 `assembleRuntime` 都不跑）。
⇒ 另一条更硬的读码结论：**L2 卡才走得通这条路**——`replies.go:356-357` 在 `card.Grant == ""` 时返回 `ErrRouteHasNoAllow`，而 L1 窗口没有 allow（SPEC-06 §2 B1），所以 `:673` 那一支**只在常驻腿对一张 L2 卡打 `session` 时才被踩到**。〔仅读码，未跑〕

### 2.2 ② 那本 ledger 在常驻腿建门的时刻存不存在

**两件事分开答，各自给 `file:line`。**

- **存不存在：建门那一刻（`resident_windows.go:132` → `resident_approval_windows.go:219`）不存在。**
  它在全仓**唯一的产码构造点**＝`cmd/wisp/run.go:478`（`ledger, lerr := session.NewLedger(session.LedgerOptions{ID: sessID, Store: mem, Logf: rt.auditf})`），
  前置是 `run.go:473`（`sessID, err := session.Mint()`），落到对象是 `run.go:487`（`rt.session = ledger`）。
  三行都在 `assembleRuntime` **函数体之内**（函数起于 `run.go:382`），而 `assembleRuntime` 在常驻腿里要到 `resident_windows.go:260` 才第一次被叫。
  尺＝`git grep -n "session\.Mint\|session\.NewLedger" HEAD -- cmd internal tools | grep -v _test` ⇒ **产码命中只有 `run.go:473`／`run.go:478` 两枚**。
  ⛔ 别把 `internal/session/grants.go:166`、`queue.go:563` 读成第二枚构造点：前者是 `Record` 的失败包装，后者是 `Queue.grants *grantStore`（一次性 nonce 表，同名不同物）。
- **被没被递进去：没有，而且"那一刻"根本没有值可递。**
  递的入口只有一枚＝`gate.go:160`（`grants:   o.Grants,`，`New` 的构造期赋值）。晚绑定三词现跑＝0 命中：
  尺 `git grep -n "SetGrants\|AttachGrants\|WithGrants" HEAD -- cmd internal tools` ⇒ **0**（本腿复跑，与票 256 §8.1 那句〔编排者复跑〕同读数）。
  `runSpec` 也**没有**账本位：字段全集现读 `run.go:102-174`＝`argv/stdout/stderr/dataDir/notify/probeSink/now/onRuntime/modeConfirm/sink/reply/replyVeto/gate/ui/cards/taskCtx`（**16 枚**），**无 ledger／无 Grants**。
  ⇒ 所以即使把建门挪到 `run.go:478` 之后，也还差一枚"把 ledger 交回给 `residentApproval`"的形状。

- **"把建门挪到它之后"这一形在时序上可行不可行：不可行（读码级，且这不是本腿新造的结论）。**
  现读链条与**票 256 §7-5／§8.1 说的是同形**，逐处对得上：
  `resident_task_source_windows.go:218`（`interactiveStdin()`）→ `:224-230`（无控制台 ⇒ 打印 `taskEntryDisabledClaim` 后 `return nil`）→ `:265`（`assembleRuntime(spec)`）。
  ⇒ 门一旦挪到 `assembleRuntime` 之内，**双击／Explorer 拉起那一支（正是 D2／票 228 那个"用户真正启动的进程"形状）永远走不到装配**＝没有门、没有 `PendingWindow`／`PendingApproval`、球不进 `Confirming`、取消键不借、Esc 通道不加载。
  这就是票 256 §7-5 那句"**⛔ 这不是'那一维变弱'，是载体不存在**"，也是母票 248 AC#10 定 ⓑ（不移动）的理由；**本腿复认：那条提前 return 今天仍在建门之后、装配之前，形状未变。**
  ⚠ 我与那两票的**唯一差别是枚数与行号**，不是结论：票 256 §现量.1 引的是 `resident_approval_windows.go:108-113`（256 立票时）、§7-2 更正为 `:109-113`（256-a1 复认时），**今天真身是 `:219-225`**；
  票 256 §8.1 说"账本构造点在门构造点之后约 129 行"（`resident_windows.go:126`↔`:255`），**今天是 `:132`↔`:260`＝128 行**。两处都是注释里的自述数字，没有一枚尺钉着它（见 §4 的 N#7）。

### 2.3 ③ `run.go` 那两枚（`:618 Grants: grantWrite`／`:758 Grants: grantRead`）与常驻腿的关系

**具名答：同一枚函数、同一次进程，但两枚的"可达性"完全不同——一枚常驻腿永远走不到，另一枚常驻腿天天走。**

- 两枚都在 `assembleRuntime`（`run.go:382` 起）里，而**常驻腿确实调用它**（`resident_task_source_windows.go:265`，把 `ra.gate`/`ra.ui`/`ra.cards`/`ra.root` 注进去）。⇒ **不是"根本不同的入口"**：`wisp run` 与常驻腿**共用同一枚装配根**，区别只在 `runSpec.gate` 有没有被注入。
- **`grantWrite`（`run.go:618`）与常驻腿完全无关。** 它在 `if s.gate != nil {…} else {…}` 的 **else 支里**（`:600` 判、`:609` 起 else、`:612-619` 建门、`:628` 闭 else）。常驻腿 `s.gate != nil` ⇒ 走 `:600-608` 那支（`rt.gate = s.gate`），**`:612` 那枚字面量在常驻进程里一次都不执行**。
  ⇒ 推论（本腿自己拉的尺，不照抄票面）：**"生产建门只有两枚"这句在枚数上对，在"两枚各归一条腿"这句上也对；但它会让人误读成"run.go 接了、常驻没接"＝同一条码路径上的两个兄弟。真相是常驻腿与 `grantWrite` 分属互斥的两支，同一次进程里不可能都跑到。**
- **`grantRead`（`run.go:758`）常驻腿是跑得到的**：它在 if/else **之外**（`tools.New(tools.Options{…})` 起 `run.go:745`），两支都过。
  ⇒ **常驻腿今天是"读侧已接、写侧孤儿"**：`run.go:573-578` 现铸的 `grantRead`/`grantWrite` 里，读侧进了桥（`internal/tools/grant.go:65-68`，nil 则 `return 0` 按未授权处理）、写侧进了那枚**常驻进程根本没建**的门。
  净结果：常驻进程里那份 ledger **被铸出来、被挂到桥上、然后没有任何东西往里写** ⇒ `Covering` 永远查不到行 ⇒ **每次都问**。
  ⚠ 这一形**不是**"悄悄放宽"，也**不是**"记到了别的会话"——它比票面写的更保守一档（fail-closed），但代价是 ledger 那半条链在常驻进程里是**活体死代码**。〔仅读码，未跑〕
- **"常驻腿今天有没有任何一处把 grants 补上"＝没有。** 尺（本腿自拉，⛔ 不照抄票面）：
  `git grep -n "Grants" HEAD -- cmd ':!*_test.go'` ⇒ 8 命中里**赋值只有 2 枚**（`run.go:618`／`run.go:758`，都在装配根内、都归 `wisp run` 那一支或两支共用的读侧），
  其余 6 命中是**注释**（`run.go:595`、`run.go:643`、`approval_reply.go:327`、`run.go:1398`、`resident_approval_windows.go:186`、`resident_task_source_windows.go:42`）＋ 1 枚 `run.go` 的行内说明。
  再补一枚窄尺：`git grep -in "grants" HEAD -- cmd/wisp/resident_windows.go cmd/wisp/resident_task_source_windows.go cmd/wisp/resident_approval_windows.go` ⇒ **产码赋值 0 枚，三枚命中全是注释**。
  ⇒ **常驻腿三条腿文件（resident_windows／resident_task_source／resident_approval）里补 grants 的地方一处都没有，也没有别处（`approval_always.go`／`panel_*` 全查零命中）。**

### 2.4 ④ 票面 §现量 4 那两处待复认（答复浓缩在此，全读数在 §1.4）

- **"resident_windows.go prints the same limit at boot"＝不印。** 真印的那句限制不在 boot，在**答复那一刻**：
  逐字句子（`gate.go:674-676`）＝`approval: GRANT-DROPPED corr=%s tool=%s paths=%d (本机没有接入会话授权记账，本次按「仅本次」放行，没有落盘任何规则)`。
  印到哪儿去（常驻腿的 sink 是 `resident_approval_windows.go:329-333 residentAuditf`，两路同时走）：
  ① `slog.Info("audit: " + line)` → 进程默认 logger ＝ `logsink.go:153-161` 那枚 `teeHandler`，**primary＝`<dataDir>\logs` 下的 redact JSONL 滚动件**（`logsink.go:144-149`＋`:76 logDirName = "logs"`），**mirror＝`os.Stderr`**；
  ② `fmt.Printf("wisp: [audit] %s\n", line)` → **stdout**（不是票面说的 stderr）。
  ⇒ **票面"唯一的痕迹是一行 stderr"这句在本腿读下来不准确**：落点其实是 **stdout ＋ stderr ＋ `<dataDir>\logs\*.jsonl`** 三处，其中只有 JSONL 那枚在 GUI 形状里真存在。
- **一个没有控制台的 GUI 用户看不看得到：看不到。**
  尺：`scripts/build.ps1:103-115` 逐字给 `wisp.exe` 传 `-H=windowsgui`（票 244），⇒ 产物是 GUI 子系统；
  `resident_windows.go:55-60` 注释逐字"**double click the icon, no terminal attached, stderr going nowhere**"（票 117 立此段）；
  `main.go:65 attachParentConsole()`（`console_windows.go:33-43`）只在**有父控制台可挂**时救得回来，Explorer 拉起没有父控制台 ⇒ stdout/stderr 两路**无处可去**，只剩 JSONL 文件那一路。
  ⇒ 净答案：**界面上看不见（面板没有审批入口、球只换状态、托盘没有 allow），日志文件里有，终端里两条都有——而"终端里两条都有"那一形已经不是双击的 GUI 用户了。**〔仅读码，未跑〕
- **`resident_approval_windows.go:186` 那句注释仍是真话**（见 §1.4②）。

### 2.5 ★ 本腿量到的一枚票面之外读数（不在四问射程，但直接压 AC#1 与 ⓑ）

`approval_reply.go:279-281` 是 `session` 动词的**成功回执**，逐字：
`已按「本会话内允许」答复 <corr>（<tool>）：这一发已放行，卡片上那些路径对本会话后续的 L1 询问不再重复提问（L2 永不被会话授权覆盖，SPEC-06 §8.3 第 1 条）`。
它的返回路径**只看 `AllowSession` 有没有报错**（`:260-276` 三个 `case errors.Is(...)` ＋ `:276 default err != nil`），
而 `gate.allowSession` 在 `g.grants == nil` 那一支**先放行、再落 `GRANT-DROPPED`、然后 `return nil`（`gate.go:673-677`）** ⇒ 报错路径不触发 ⇒
**常驻腿上这句话今天照样逐字印出来**，而它承诺的"不再重复提问"在常驻腿上**不发生**（写侧孤儿，见 §2.3）。
⇒ 落点＝`resident_task_source_windows.go:319 src.runConsoleLoop(ctx, console, os.Stdout)` → `:387 fmt.Fprintf(out, "wisp: %s\n", text)`，**stdout**。
⇒ 这一枚是"可见性缺"的**反向形**：不是"没话说"，是**说了一句做不到的话**，且它由 `err == nil` 那一路自然产出、没有一枚尺在钉它（`ticket224_assembly_test.go:198` 只断"这句里含 corr"，不断含"不再重复提问"）。
⚠ **它归 ⓑ 的射程（票面 ⓑ＝"说到明处"），但它把 ⓑ 的语义翻了一面：ⓑ 不是"补一句真话"，是"改一句已经在说的假话"。** 本腿不裁形，只把这枚摆上台。〔仅读码，未跑〕

## 3. 三形代价表（⛔ 本腿不裁形，只把料摆齐）

### 3.0 三形共同的前置：两条腿的形状差别（先摆这一格，不然三形的"作用面"读不准）

常驻进程今天有**两形**，三形各自的代价在这两形里完全不同：

| 形状 | 触发条件（产码 `file:line`） | 门在不在 | ledger 在不在 | 有没有答复入口 | 有没有卡片 |
|---|---|---|---|---|---|
| **带控制台**（终端里跑 `wisp.exe` 无参） | `resident_windows.go:132` 建门 → `:260` → `resident_task_source_windows.go:218` 拿到 stdin → `:265 assembleRuntime` | 在（`ra.gate`） | **在**（`run.go:473→478→487` 在这一发里真跑了） | **在**（`runConsoleLoop` 的 `session` 动词） | 有（L1 窗口／L2 队列卡） |
| **双击／Explorer**（D2／票 228 那个"用户真正启动的进程"形状） | `resident_task_source_windows.go:224-230 return nil`（`interactiveStdin()`＝nil，`approval_reply_stdin_windows.go:41-52`） | 在（`ra.gate`） | **根本不存在** | **不存在**（`src.surface == nil`，`:287-290` 那条件不成立） | **不存在**——`assembleRuntime` 没跑，没有桥、没有任务；`AskOnTaskRoot` 产码调用者 0 枚（尺 R38） |

⇒ **三形都必须对着这两形各读一遍**：ⓐ 只能修上面那一行，ⓑ 只能让上面那一行的用户看见，ⓒ 两行一起登记。下面每形都按这两形分档写。

### 3.1 ⓐ 接上（常驻腿建门时也把会话记账递进去）

**最小写面（逐枚具名；本腿按"能不改 `internal/agent/approval` 一字"那一支摆，另一支单列）**

形 ⓐ-Ⅰ（**晚绑定 holder，`GrantRecorder` 由 `cmd/wisp` 自己实现**；尺＝`gate.go:66-68` 那个接口只有 `Record` 一个方法，任何包都能实现）：
1. `cmd/wisp/resident_approval_windows.go` —— 新增一枚 holder 类型（或就近一个小 struct）＋ `:219-225` 字面量里加 `Grants:`（实传 5 枚 → 6 枚）；`:185-192` 那段"WHAT THIS DELIBERATELY DOES NOT CLOSE"必须同步改写。
2. `cmd/wisp/resident_task_source_windows.go` —— 在 `:281 src.run = run` 之后、`:321 submitTask` 之前，把 `run.session` 绑进那枚 holder。
3. `cmd/wisp/resident_approval_risk_256_windows_test.go` —— ④ 那枚 AST 钉的期望集（`:445 want` 5 枚→6 枚、`:452` 那句"Grants 不许在场"必须搬家）。**⚠ 这一枚不是本票写面，动它＝改票 256 的判据面。**
⇒ 共 **3 枚文件**（2 枚产码＋1 枚测试），`internal/agent/approval` **零改动**，`cmd/wisp/run.go` **零改动**（这是它比 ⓐ-Ⅱ 值钱的地方，见 N#12）。

形 ⓐ-Ⅱ（**晚绑定入口开在 approval 包里**，即 256-a2 §0 说的"新造一枚 holder 类型"落在 `internal/agent/approval`）：
1. `internal/agent/approval/gate.go`（新增 setter／holder 类型 ⇒ `g.grants` 从"唯一写点＝构造期"变成两处写点，⚠ 那要重读 `gate.go:86-88` 那句"No method on Gate reads it back"是否仍成立）
2. `internal/agent/approval/fakes_test.go:250-279`（`newGate` 逐字段拷贝表，见 N#15）
3. `cmd/wisp/resident_approval_windows.go` ＋ 4. `cmd/wisp/resident_task_source_windows.go` ＋ 5. 256 那枚 AST 钉
⇒ **5 枚**，且新增两枚包内仪器（setter 的正／反控）。

形 ⓐ-Ⅲ（**第二枚 mint**，256-a2 §0 的 (a)）：
1. `cmd/wisp/resident_windows.go`（在 `:132` 之前开 `memory.Open` ＋ `session.Mint` ＋ `NewLedger`）
2. `cmd/wisp/resident_approval_windows.go`
3. `cmd/wisp/run.go` —— **要么让它复用那枚 ledger，要么任它再铸一枚**
⇒ 枚数看着少，**但它带着两枚硬伤**（本腿只摆不裁）：
- **gate 写 id-A／bridge 读 id-B 的错配**：`grantWrite` 走 `resident_windows.go:132` 那枚 mint，`grantRead` 走 `run.go:478` 那枚——两张 `approval_grant` 行键在 A、`tools/grant.go:65` 那侧查的是 B。⇒ 净结果不是 fail-closed 而是**"记了，但没人查得到"**，审计还会逐字打印 `GRANT-RECORDED ... grant_id=N`（`gate.go:693`）——**那是比 `GRANT-DROPPED` 更响的一句假话**。这一条直接对上票 265 AC#1 三档里的**"记到了别的会话"**。〔仅读码，未跑〕
- 同进程第二枚 store 句柄（`run.go:450 memory.Open` 之外再开一枚），D38(e) 第 7 步的关闭归属（`resident_task_source_windows.go:478-497` 只关 `src.run`）要重答。

**必须先解冻什么（三形共用现量）**
- `git status --porcelain -- cmd internal` 本腿取数＝**`M cmd/wisp/resident_approval_windows.go`（唯一一枚脏文件，`260-v1` 的突变体）**。⇒ 三形落地前都要等它变干净。
- ⛔ **`cmd/wisp/**` 与 `internal/agent/approval/**` 此刻在 `260-v1` 的取数射程内**（它要在这两个包跑定向 `go test` 与突变）；票 265 §排程逐字"同包一律串行，不接受'改的是不同文件'这种推理"。⇒ **ⓐ 三形今天一枚都不能开工**，且 ⓐ-Ⅱ 还多撞 `internal/agent/approval`（同在其射程）。
- ⓐ-Ⅲ 额外要**具名解冻 `cmd/wisp/run.go`**——票 256 §8.2 硬边界②那句"本轮不许碰 run.go"虽是为 `256-r1` 写的，但 P7 那枚钉（N#12）是事实约束、不因换票而失效：**动 `run.go:424` 以后的行 = 四处 evidence drift 必红**，而那枚名册归票 255／248 AC#8，不在本票写面。
- 256 那枚 AST 钉（N#4）的期望集搬家＝**改票 256 的判据**，须编排者落一枚具名 `A##`。

**用户看得见什么（ⓐ-Ⅰ／Ⅱ）**：带控制台那一形——他答一次「本会话内允许」，同一发进程里那条路径后续的 L1 询问**不再重复问**（`GRANT-DROPPED` 那行也不再出现）。双击那一形——**他什么都看不见，也什么都没变**（那一形根本没有卡片和答复入口，见 §3.0）。

### 3.2 ⓑ 不接，但说到明处

**★ 本腿量到的前置事实（决定 ⓑ 到底能落在哪儿）：常驻腿上今天没有任何"文字面"可写。** 逐枚尺：
- 球画的是状态与图标，**没有一行字**：`internal/ball/statevisual.go:177-184` 的 `Confirming`／`AwaitingApproval` 只填 `RingColor`／`RingPulse`／`BadgeCount`；唯一会画字的 `renderer_windows.go:493`（`BadgeCount` 的数字）与 `:508-510`（`BadgeText`）——`SetBadge`／`SetBadgeText` 的**产码调用者只有 `cmd/balldebug/main.go`**（尺 R91/R92），常驻腿零调用。
- 托盘 tooltip：`internal/ball/ball_windows.go:960-963 SetTrayTip` **产码调用者 0 枚**（尺 R83）；初值是常量 `"Wisp"`（`ball_windows.go:248`）。
- 托盘菜单：`tray_windows.go:86-91` 四枚项，无审批位（尺 R15）。
- 面板：常驻进程**确实链了 WebView2 宿主**（`panel_host_windows.go:64` import go-webview2；`resident_windows.go:151/157` 真的建了并起了），**但审批卡的数据过不去**——`run.go:734-736` 的 `rt.ui.publish = rt.publishPanelSnapshot` 只在 `rt.ui != nil` 时装，而注入支 `run.go:607` 明确 `rt.ui = nil`，`:727-733` 那段注释逐字写着"its own UI owns its publishing - which this leg does NOT give it"；快照的 `Out` 是 `rt.bookPanelSnapshot`（`run.go:725`）＝**落 ledger，不是落页面**（`panel_pump.go:12-21` 逐字："There is no Go -> page channel in this tree … the last mile is still open"）。
  ⚠ **顺带更正三处过期声称**（票 265 AC#0 没问，但本腿读到了就具名报，见 §5.1）：`resident_approval_windows.go:32-33`、`resident_task_source_windows.go:52-53`、`approval_always.go:165` 都写着"this process links no WebView2 host"——**这一句在今天的常驻进程里不成立**（票 33 已落地）；真正成立的是"**没有到页面的通道**"，两件事别混。
⇒ **净读数：ⓑ 的"说到明处"今天能落的可见面只有 stdout／stderr／`<dataDir>\logs\*.jsonl`**（＋带控制台那一发的终端回执本身）。**双击的 GUI 形状里 ⓑ 落不了任何用户看得见的位置**——那一形没有卡片、没有文字面、也没有能收消息的页面。这一条是本件对 ⓑ 最硬的代价读数：**ⓑ 不是"便宜的降级形"，它在机主真正用的那一形里等于零可见性。**

**最小写面（按能落地的两支）**
- ⓑ-Ⅰ（改回执，**唯一在带控制台那一形真能看见的一支**）：`cmd/wisp/approval_reply.go:257-282`（`session()` 的成功回执 `:279-281`）＋ 让 `replySurface` 能知道"我这枚门有没有记账位"。⚠ 这一步今天**做不到**：`approval.Gate` 导出的 18 枚方法里没有一枚把 `g.grants` 的有无报出去（尺 R43 名册），`gate.go:86-88` 还逐字钉着"No method on Gate reads it back" ⇒ **ⓑ-Ⅰ 必须新开一枚读面**（I2），或把事实沿构造链传下去（`resident_task_source_windows.go:252-264` 的 `runSpec` 加位 → 撞 N#12 P7）。枚数：**枚不齐**——最小也得起码 `approval_reply.go` ＋ 一枚传递形状（`run.go` 或 `gate.go`），加 1 枚新判据。
  ⚠ 另有一枚**本腿读码判红的既有形状**：`replySurface` 是 `wisp run` 与常驻腿**共用**的（`approval_reply.go:476 newReplySurface`，两支都调），所以在 `:279` 直接改字＝**同时改掉了 `wisp run` 腿那句真话**（那一腿 `Grants` 有值、承诺是兑现的）。⇒ ⓑ 不是"加一句"，是"给一句分成两形"，而分形需要输入 ⇒ 回到 I2。
- ⓑ-Ⅱ（改 boot／audit 句子，只对读日志的人可见）：`cmd/wisp/resident_approval_windows.go:229-234`（建门回执那条 `slog.Info`）或 `:589-598 residentStatusLine`。枚数：**1 枚产码＋1 枚判据**，最小。**但它撞两枚零漂移钉**（N#9 逐字比较 `residentStatusLine`、N#10 禁止那几行出现 `Esc`），且**对用户仍然不可见**。
**必须先解冻**：同 3.1（`cmd/wisp/**` 在 `260-v1` 射程内）；ⓑ-Ⅰ 若走 `gate.go` 读面则连 `internal/agent/approval/**`；ⓑ-Ⅰ 若走 `runSpec` 传位则**必须具名解冻 `cmd/wisp/run.go`＋票 255 名册**。
**用户看得见什么**：ⓑ-Ⅰ＝带控制台的人会在答复那行看见"这条路上『本会话内允许』只到本次为止"；双击的人**看不见**。ⓑ-Ⅱ＝**没有用户看得见**，只有读日志的人看得见（⛔ 我不许把这一支写成"说到明处"，它写的是"说到日志"）。

### 3.3 ⓒ 登记成 DEFERRED（五字段）

**最小写面（逐枚）**
1. `docs/specs/SPEC-12-roadmap-governance.md` §5 表加一行（五字段现读表头＝`docs/specs/SPEC-12-roadmap-governance.md:64`：`类型 | 项 | 为什么现在不做（依据） | 完成判据 | 前置 | 当前残缺表现`）。⚠ **`docs/specs/**` 在本腿与本票 AC#2／AC#3 的禁改清单里**（票 256 §禁区、票 265 AC#2 逐字），⇒ ⓒ 落地需要**编排者自己动那一行或具名解冻**，任何产码腿都不许碰。
2. 产码里的具名标记一枚（最自然的位置＝`cmd/wisp/resident_approval_windows.go:185-192` 那一段旁边）。
3. 台账 `docs/reports/pending-and-issues.md` 一枚 `A##`（归编排者写）。
⇒ 产码枚数：**1 枚**（最省的一形）。
**★ 但本腿量到 ⓒ 的一枚硬撞（不在票面上，必须摆出来）**：`SPEC-12:95-96` 逐字规定代码内标记的形状是
`// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点`，
而 **`docs/DEFERRED.md` 在盘上不存在**（尺＝`git ls-files docs | grep -E "DEFERRED|DECISIONS"`＝**空**；AGENTS.md §4 末段也具名说过这批交付物"截至锚点 `4e66817` 在仓里不存在"）。
`tools/d22scan` 的 **ban #9 phantom-citation**（`main.go:799-873`，尺路径正则 `:890-891` 覆盖 `docs/…` 全拼路径）判的正是"产码注释引一条盘上不存在的仓内路径"⇒
**照 SPEC-12 那一句写的 ⓒ 标记会被 d22scan 直接判红**（今天仓里 28 处 `DEFERRED(...)` 标记**没有一处**带 `docs/` 引用，尺 R111＝0 命中，所以这枚撞是**新造的**、不是既有的）。
⇒ ⓒ 的落地形状必须**要么不写 `docs/DEFERRED.md` 那枚锚点、要么先由人工批准建那一枚文件**——这是 D22 闸门③意义上的"未定义即停"，本腿只摆不裁。〔仅读码，未跑——真跑 d22scan 属 U6〕
**必须先解冻**：`docs/specs/**`（编排者面）＋ `cmd/wisp/resident_approval_windows.go`（`260-v1` 射程，等它变干净）。
**用户看得见什么**：**没有。** ⓒ 是唯一一形"三形里对用户零变化"的——它的价值全在盘上那行登记与下一任读者能不能查到。⚠ 写这一形给机主交代时不许说"以后会看得见"。

## 4. 撞钉名册（逐枚：它今天断什么／哪一形会让它翻／翻了算修 bug 还是算改契约）

（本节答：尺＝`git grep -n "GRANT-DROPPED\|grants == nil\|Options.Grants" -- '*_test.go'` 起手，**但名册不只这一族**：
 行为型钉只写产物不写常量名，所以另拉（a）语法树型 Options 字段集钉、（b）boot 打印／状态句**字面零漂移**钉、
 （c）`file:line` 型 evidence 钉（票 255 名册）、（d）"某词必须不出现"型反向钉、（e）单门指针同一性钉。
 每枚给 `file:line` ＋ 今天断什么 ＋ ⓐ／ⓑ／ⓒ 各形会不会翻 ＋ 翻了算哪种。）

## 5. 顺带上报（看到就报，⛔ 不许顺手改）

（本节答：（a）别的同类接缝——`Options` 里还有哪枚字段在某条生产腿上留空、后果是什么、有没有尺；
 （b）`GRANT-RECORD-FAILED`／`GRANT-RECORDED` 有没有尺在断（尺自拉）；
 （c）票 224／256／260／266 里与票 265 现量**冲突**的句子，指到行号，⛔ 不改别的票的票面。）

## 6. 判不动／量不到（具名＋归口）

### 6.1 本腿**没跑**的尺，逐枚具名（每枚都写了"为什么跑不了＋谁能跑＋跑了会读到什么"）

⛔ 下面每一枚都是**该票判据要用、而本腿只能给读码预测**的。凡在本件其它节里被引用为结论的，都已就地标〔仅读码，未跑〕。

| # | 没跑的尺（原样命令） | 为什么本腿跑不了 | 归口／跑了会读到什么 |
|---|---|---|---|
| U1 | `go build ./...` | §1 禁令（禁任何 go 命令）；`260-v1` 独占取数期 | 门禁读数由编排者事后补（票 265 AC#4）。本腿只保证**没让盘多脏一枚**（见 §6.3） |
| U2 | `go vet ./cmd/wisp/ ./internal/agent/approval/` | 同 U1 | AC#4 四门之一 |
| U3 | `go test -count=1 ./cmd/wisp/ ./internal/agent/approval/ ./internal/tools/ ./internal/session/` | 同 U1，且**这正是会洗掉 `260-v1` 读数的动作** | **本件 §4 全部"会不会红"的判语，今天只有读码凭据**。尤其 N#3（256 的 AST 钉）与 N#9（255 的行号钉）在 ⓐ 下**本腿断言必红但没跑** |
| U4 | `go test -count=1 -run 'TestTicket224' ./cmd/wisp/ ./internal/agent/approval/`（两枚极性相反钉的基线） | 同 U1 | 编排者若要核"两枚今天各自绿在哪一支"，只有这一发能答；本腿 §1.3／§4 的极性是**读断言源码**读出来的 |
| U5 | 定向正控（票 265 AC#3 那一发）：ⓐ 形＝摘掉那处接线必须让指名用例红 | 需要产码＋测试同时存在，本腿不写产码不写测试 | 归 ⓐ 的落地腿；⚠ 本件已具名"今天常驻腿那一格**零枚用例守**"（§5.2），所以**这一发今天无处可摘**——先要有钉 |
| U6 | `sh scripts/d22scan.sh`（卫生门） | 该脚本第一步就编 `tools/d22scan` ⇒ 属禁跑的 go 命令族 | AC#4。⚠ 票 265 AC#4 自己写了"那两枚『CI 红／本机绿』属票 212 射程，看到具名上报"——本腿**连看都没看到**（没跑），故 §5.4 报的是"量不到"而不是"没发现" |
| U7 | `powershell scripts/slo-check.ps1` | §1 明令禁跑 | AC#4 无关本票，但票面排程写了；未跑 |
| U8 | `gofumpt -l cmd internal` ＋ 票 212 ban #9 的实跑 | 前者要 go tool；后者是 d22scan 的一条 ban ⇒ 同 U6 | AC#4。**注意**：ban #9（phantom-citation）对本票特别相关——§1.4① 那句假话正是"引用了一条不存在的打印"，而 `resident_windows.go` **存在**，所以 ban #9 结构上抓不到它（它只判路径存在性，不判句子真伪）。见 §5.4 |
| U9 | 真跑一发常驻进程（`wisp.exe` 无参，带控制台／不带控制台两形），看 GRANT-DROPPED 到底出现在哪几路 | 需要 build ＋ 起进程 | **这是 §2.4／§2.5 那两条"看得见／看不见"判语唯一能证伪的尺**。归 ⓐ／ⓑ 落地腿的凭据档；本腿全部标〔仅读码，未跑〕 |
| U10 | 突变：删掉 `gate.go:673-677` 那一支 ⇒ 哪一枚用例红 | 需要 go test | 用来验"第一形有没有牙"。本腿只报 §5.2 的**尺在场／缺席**读数 |
| U11 | 突变：把 `resident_approval_windows.go:219` 那枚字面量加上 `Grants:` ⇒ N#3 红、别的钉动不动 | 需要 go test | **这是 ⓐ 的作用面读数**，本件 §3.1 的"必须先解冻／会撞的钉"全是读码预演，⛔ 不许被读成实测 |
| U12 | 突变：删掉 `approval_reply.go:279-281` 那句成功回执 ⇒ 有没有一枚尺红 | 需要 go test | 本腿读码判**零枚**（§2.5／§5.2）。这一发是 ⓑ 的正控原型（票 265 AC#3 ⓑ 支要求的正是这个形状） |

### 6.2 判不动的地方（⛔ 不用"应该没问题"填空）

| # | 判不动的那一格 | 本腿读到哪一步就停了 | 归口 |
|---|---|---|---|
| I1 | **"双击 GUI 那一支今天到底有没有任何一张卡片"** | 读满的部分：`resident_task_source_windows.go:218→:224-230` 无控制台即 `return nil`（管线不装配）；`AskOnTaskRoot`／`askConfirmation` 产码调用者 **0 枚**（尺 R38）；面板那侧 `Message: nil`（`panel_inbound.go:277`）＋ `panel.approval.request` **在 Go 的受理表里不存在**（`internal/panel/bridge.go:146-152 knownComposerMethod` 六枚：mode/workspace/attachment/message/config.get/config.set）。⛔ 判不动的部分：**我没有把整棵 `cmd`＋`internal` 的每一枚 `PendingApproval`／`PendingWindow` 产码调用者穷举**，所以"这一支今天一张卡都举不出"这句**全称负向我不写**，只写"本腿穷举到的三条入口都不通" | ⓐ／ⓑ 定形前必须补 U9（真跑一发无控制台的常驻进程）；归口票 265 AC#1 |
| I2 | **⭐ ⓑ 那枚"已按『本会话内允许』答复…不再重复提问"要不要按支拆句** | 本腿量到了形状（§2.5），**判不了形**：拆句要 `replySurface` 知道"我这枚 gate 有没有记账位"，而 `approval.Gate` **今天没有任何一枚导出方法把 `g.grants` 报出去**（导出名册尺 R43：`Queue`/`Channels`/`Window`/`AdmitTextTask`/`PendingWindow`/`Complete`/`Veto`/`LateVeto`/`PendingApproval`/`Native`/`Panel`/`DecideFromNative`/`DecideFromPanel`/`Replay`/`Bus`/`Report`/`ToolsCancelBus`/`LiveL1Windows`——**没有 `HasGrantWriter` 之类**）。⇒ ⓑ 的最小写面里**藏着一枚新接缝（新增导出读面）**，而那要撞 `gate.go:86-88` 那句"No method on Gate reads it back"——这枚边界是不是 D45 的契约面，**归人工批准**（`SPEC-12 §4.1`），不是腿能定的 | 编排者裁；裁之前 ⓑ 的"最小写面"只能写"待定" |
| I3 | **ⓐ 走 (a) 第二枚 mint 会不会破 A435 第 1 条** | `run.go:459-467` 注释逐字 "One mint per process, from crypto/rand, and nothing derived: A435 clause 1 forbids recomputing it…"。**这句是注释还是契约条文，本腿判不了**：`PLAN.md:1642` 我读了射程（ledger A591 引它"A435 第 1 条"），但 **A435 原文里"一枚 mint"是硬约束还是描述**——本件不裁 | ⛔ 不许任何腿据此选 (a)。要选 ⓐ 先由编排者把 A435 原文那一格贴出来 |
| I4 | **常驻腿上"会话"这个词对用户的含义** | 票面禁口令「224 会话时长改跨重启」＝"本次会话内＝一次进程存活期"是既有裁定。本腿读到：**常驻进程"一次进程存活期"对用户来说＝"开机到关机"**（`resident_windows.go` 无重启路径），而 `run` 腿的＝"一次命令"。**同两个汉字，两种长度**。判不动：这是不是要算 D45 的语义漂移，属产品裁量 | 摆给编排者＋机主，⛔ 本腿不改票面那句 |
| I5 | **N#3 那枚 AST 钉翻了以后算"修 bug"还是"改契约"** | 本腿能说的是：它的错误文本**自己写了**"the half was filed as pending ticket 265"（`resident_approval_risk_256_windows_test.go:453-455`），**即它是票 256 主动留在盘上、让票 265 去翻的**。但"票 256 的 AC#1 期望集是 5 枚字段"这件事被写成了钉死形态 ⇒ **翻它需要编排者落一枚具名 `A##`**，⛔ 腿不许自己判成"修 bug" | 归编排者；本件 §4 N#3 已给结论 |
| I6 | **`GRANT-DROPPED` 该不该由"日志行"升级成"界面事件"** | 这是 ⓐ／ⓑ 之外的第四形（把 drop 做成 D43 状态或面板一行）。本腿能判的是：`bookWaitingState`（`approval_always.go:158-180`）已经给了"名字从冻结表里出来"的先例，而 D43 转移表＝C12 冻结件 ⇒ **任何新状态名一律是人工批准面** | ⛔ 本腿不列进三形表（票面只给三形），只在此登记 |
| I7 | **本腿能不能证实 `run.go:597` 那句"prints"是"从来没印过"还是"曾经印过后来丢了"** | `git log -S`／`-G` 三发都跑了：`-S"prints the same limit"` 只命中**票 265 立票那一发**（`35e852f9`，因为票面正文含这句），`-G"prints the same limit" -- cmd/wisp/run.go` **0 命中**，`-S"会话" -- cmd/wisp/resident_windows.go` **0 命中**。⇒ 读起来支持"从未印过"，但 `-G` 在重命名／折叠行上是**出了名的钝尺**，⛔ 我不把"从来是假话"写成判死 | 若 ⓑ 定形需要归因（是谁写的假话），需要一发 `git log -L` 型尺，归编排者或产码腿 |

### 6.3 本腿对盘的污染面（自证）

- `git status --porcelain -- cmd internal` 交件时刻＝**仍只有 `M cmd/wisp/resident_approval_windows.go`** 一枚（`260-v1` 的突变体，⛔ 不是本腿写的，本腿没动它）。
- 本腿写面＝`.scratch/wisp/probes/265/a1/**` 一枚文件（本件）。临时件全在 `/tmp`（`gate_head.go`／`raw_head.go`／`run_head.go`／`rw_head.go`／`rts_head.go`／`ar_head.go`／`t256.go`／`t260r3.go`），**都在仓外**；仓内临时件只建不删这条对本腿不适用（本腿没在仓内建临时件）。
- `sh scripts/check-path-length-budget.sh --with-self-test` 本腿**跑过一发**（纯 shell，不碰 Go，§1.2 那把尺的分母要引用它）：VERDICT GREEN；`tracked paths=5536／over-budget=57／roster entries=57／not in roster=0`；三发正控全 ok（正控台件 `bench=/tmp/tmp.nG43jn8EhK`，脚本逐字写着"kept on disk; this project never deletes temp artifacts"⇒ 该台件在 `/tmp` 下，不在仓内）。
- ⚠ 本件的**路径长度**自量：`.scratch/wisp/probes/265/a1/census.md` ＝ 37 字符，远低于 issues 目录那枚 100／121 的帽子（帽子读数见上一行 `hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121`）。

## 7. 本腿自己写错的尺与读数（全数留下，不当笔误藏）

（本节答：本腿跑过的每一枚尺的原样命令与读数存档；写错／重跑／更正的行数逐条登记，⛔ 不删。）

---

## 〔编排者标注 2026-10-04 14:07:13 +0800（＝06:07:13 UTC）〕本腿死在 §7，两件事分开记

- **本腿撞 150 轮上限**（通知逐字：`Reached the maximum turn limit (150)`，168 次工具调用／2,409,128 ms）。盘上终态＝**307 行／47,045 字节、占位尺 0 命中**；§0–§6 正文已写满（`afbeb2c7` 落了 §6，其后 §2.5／§3／§4／§5／§6.3 的增量**未提交**，由我以显式 pathspec 代提保住，⛔ 不改它一字）。
- ⛔ **§7「本腿自己写错的尺与读数」只有标题与那句"本节打算答什么"，正文我没代填**。理由（既有定式）：那一节是实现腿对**它自己的尺**的自证，我填了就把"谁做的判"洗混；它这一枚是只读普查腿，§7 的内容只有它自己知道跑过哪些原样命令。⇒ **§1–§6 的读数我照收（它按派单逐条标了〔仅读码，未跑〕）**，§7 这一格归下一枚接手 `265-r*`／验收腿自己去量它写错的尺时补，本件不作 §7 的凭据用途。
- ⚠ 本件里 §6.3 那句"交件时刻 `cmd/wisp/resident_approval_windows.go` 仍是 `M`"是**时序读数**（那正是 `260-v1` 的突变窗口）。编排者 13:52:00 与 14:06:29 两次现跑：`git status --porcelain -- cmd internal`＝**空**，且该文件 md5 ＝ `git cat-file blob HEAD:` 的 md5（`7a26c7a990dbd2351bdf9898b5bdc192`）⇒ **无残留突变**，那一行不必被下一任读成事故。
