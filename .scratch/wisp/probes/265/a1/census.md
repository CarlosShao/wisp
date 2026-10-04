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

### 4.A 本节的锚点与尺的射程（接续腿 `265-a1b` 补，⛔ 不改上面那句意图句一字）

- **锚点是散动的，本节逐把尺自带锚点**。本腿开工时 HEAD＝`c528035f`（15:34:03 +0800），交件时已推到 `862d1736`——在飞的写腿一枚接一枚上盘，所以**下面每一条读数都写着它自己的锚点与取数时刻**，原文读数存档在 `.scratch/wisp/probes/265/a1b/logs/00—17*.txt`。§0–§6 用的是 `35e852f9` ⇒ **凡本节与前人节行号不一致，差的是这 8 小时里 256/260/265 三批上盘的量，不是谁读错**（例：`resident_approval_windows.go` 的门字面量在前人节是 `:219-224`，在 HEAD `ba900551` 也是 `:219-224`，而**盘上工作副本已是 `:368-374`**）。
- ★ **本节写作的最大变数，先说白**：ⓐ-Ⅰ 的写腿 `265-r1` **在本腿普查期间已经上盘到工作副本**（`git status --porcelain -- cmd internal` 于 08:57Z 起＝`M cmd/wisp/resident_approval_windows.go`＋`M cmd/wisp/resident_task_source_windows.go`＋`M cmd/wisp/resident_approval_risk_256_windows_test.go`）。⇒ 所以本节 **"它今天断什么"一律按 HEAD 判**（HEAD 才是 ⓐ-Ⅰ 起跑前的世界），**"ⓐ-Ⅰ 会不会让它翻"两路都答**：一路是按 ⓐ-Ⅰ 设计推演（〔仅读码，未跑〕），一路是本腿**读到在飞工作副本里那一形真身**（`Grants: ra.grants` 已进字面量、`residentGrantHolder` 已存在、256 那枚测试已被改 ＋40/−20）。⛔ 本腿**不改写、不评价**那一枚在飞件，只在 §5.D 登记两格"我在它盘上看到的、没人钉的形状"。
- 裁形已定＝**ⓐ-Ⅰ**（票面「编排者裁定」节＋台账 `A601`／`A602`）⇒ "会不会翻"那一列**只对 ⓐ-Ⅰ 答**，ⓑ／ⓒ 只在它能反衬 ⓐ-Ⅰ 代价时附带一句（N5、N11 两枚就是这种）。
- 判定用词：**必翻／可能翻／不翻／量不到**；归类用词：**修 bug（本票写面自带）／改契约（须具名 `A##`）**。凡"必翻"而本腿没跑，逐条标〔仅读码，未跑〕。
- ⛔ **本腿一枚 go 命令都没跑**（派单硬约束），所以本节**没有任何一发"我跑过这枚钉，它红了/绿了"**。名册的枚数一律 `| wc -l` 收尾，⛔ 无 `head -N` 截断当枚数（唯一用了 `head` 的是为了限宽打印长文本，其计数在同一读数档里另发了 `wc -l` 一发）。


### 4.B 族⓪ 起手尺（票面点名的三词族）＝**5 命中／3 枚文件／承重钉 4 枚**

尺＝`git grep -n 'GRANT-DROPPED\|grants == nil\|Options.Grants' -- '*_test.go'` ⇒ **5 命中**（读数档 `.scratch/wisp/probes/265/a1b/logs/00-open-ruler.txt`，取数时刻 2026-10-04 08:08:20Z，锚点 `40aae961`）。
⚠ 那把尺的 glob 是全仓形，本腿补了一枚射程尺：`git ls-files '*_test.go'` 的一级目录分布＝`.scratch 59／cmd 62／internal 276／tools 5`（合计 402）⇒ **`.scratch/**` 里那 59 枚 `_test.go` 今天零命中**，且 go 工具不进点号目录，它们是**惰性件**，不在任何分母里。

| 枚号 | `file:line` | 它今天断什么 | ⓐ-Ⅰ 会不会让它翻 | 翻了算哪种 |
|---|---|---|---|---|
| N1 | `cmd/wisp/resident_approval_risk_256_windows_test.go:452` | `if got["Grants"] { t.Errorf(...) }`＝**读语法树**断常驻腿那枚 `approval.Options` 字面量里不许有 `Grants` | **必翻**（ⓐ-Ⅰ 的定义性动作就是加这一枚字段）〔仅读码，未跑〕 | **改契约**＝动票 256 判据面 ⇒ 已具名解冻 **`A602`**（票面「裁定」节第 3 条写的是 `A601`，**那一句指错了**，见 §5.C-3） |
| N2 | `internal/agent/approval/ticket224_reply_grant_test.go:283` | `if !f.log.has("approval: GRANT-DROPPED")`＝断"没有记账位时这句必须在" | **不翻**（ⓐ-Ⅰ 零改动 `internal/agent/approval`） | —（但它是 ⓐ-Ⅰ 不许把 drop 句子抹掉的护栏） |
| N3 | 同文件 `:287` | `if f.log.has("GRANT-RECORDED") \|\| f.log.has("grant_id=")`＝没有 recorder 就不许声称记了 | **不翻**（同上） | — |
| N4 | `cmd/wisp/ticket224_assembly_test.go:194` | `if strings.Contains(h.err.String(), "GRANT-DROPPED")`＝**`wisp run` 那条腿上这句不许出现** | **不翻**（ⓐ-Ⅰ 零改动 `run.go`；且它跑的是 else 支建的门，与常驻门互斥，见 §2.3） | — |
| （非钉） | `resident_approval_risk_256_windows_test.go:43`、`ticket224_reply_grant_test.go:128` | 两处是**注释**命中，不承重 | 不翻 | — |

★ 这一族**漏不掉的两枚**（同文件、同函数，只是不含那三个词，所以起手尺照不到）：
- N5 `internal/agent/approval/ticket224_reply_grant_test.go:277-279`＝`if err := f.g.Native().AllowSession(...); err != nil { t.Fatalf("a host with no ledger must still be answerable") }` ⇒ **它钉死了"没有记账位 ⇒ `AllowSession` 仍返回 nil"**。这一枚就是 §2.5 那句假话能活着印出来的**机器原因**：**ⓐ-Ⅰ 不碰它（不动 approval 包）；但任何"改用 error 路径来表达这条腿没记账"的形（＝ⓑ 的最省事走法）都会当场把它打红**，而那属**改契约**（D45 答复语义），不属修 bug。〔仅读码，未跑〕
- N6 `cmd/wisp/ticket224_assembly_test.go:189-190`＝逐字 `fmt.Sprintf("approval: GRANT-RECORDED corr=%s grant_id=%d tool=%s", ...)` ⇒ **钉的是 `gate.go:693` 那行日志的字面格式**。ⓐ-Ⅰ 不改 `gate.go` ⇒ 不翻；⚠ 但它意味着**落地腿不许顺手给 `GRANT-RECORDED` 加"来自哪条腿"的后缀**（那会同时打红这两枚）。〔仅读码，未跑〕

### 4.C 族(a) 语法树型 `Options` 字段集钉＝**全仓只有一枚实例（N1 所在函数），但它内部有 6 条断言臂**

尺（本腿自拉，读数档 `logs/01-family-a-ast.txt`＋`logs/01b-counts.txt`）：
- `git grep -ln 'parser\.ParseFile\|go/ast\|go/parser' -- cmd internal tools \| grep _test` ⇒ **19 枚文件**用语法树当尺（名册在该读数档，ⓐ-Ⅰ 射程内只有 `cmd/wisp/resident_approval_risk_256_windows_test.go` 与 `cmd/wisp/leg_dispatch_gate_133_test.go`／`leg_sink_gate_131_test.go` 三枚，后两枚见 §4.H）；
- `git grep -n '"Grants"\|"ApprovalTimeout"\|"Channels"\|"Logf"' -- '*_test.go'` ⇒ **2 命中，全在同一枚函数里**（`:445` 期望集＋`:452` 反向断言）⇒ **这一族确实只有一枚钉，编排者的"已知一枚"读数成立**，本腿不重复立案，只复认真身。

真身与射程复认（锚点现量＝`git show HEAD:cmd/wisp/resident_approval_risk_256_windows_test.go`，函数 `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims` 起 `:425`）：

| 臂 | `file:line` | 逐字断什么 | ⓐ-Ⅰ（6 枚字段）之下 |
|---|---|---|---|
| a-1 | `:428-433` | `if literals != 1` ⇒ 那枚文件里 `approval.Options` 字面量**必须恰好 1 枚**（它自己写明是 `TestAC246ResidentPipelineAsksThroughTheOneGate` 的"更早更便宜的版本"） | **不翻**——ⓐ-Ⅰ 是往既有那一枚里加字段；⚠ 但若落地腿"另建一枚门给 holder"＝立刻红，且红句会指到 §4.G 那枚指针同一性钉 |
| a-2 | `:445-451` | `want := []string{"UI","Channels","Window","ApprovalTimeout","Logf"}` 逐枚 `got[w]` | **不翻**（6 ⊇ 5，缺项才红） |
| a-3 | `:452-456` | `if got["Grants"]` | **必翻** ＝ N1 |
| a-4 | `:457-468` | 反向走查：`got` 里任何不在 `want` 的名字 ⇒ "unexpected Options field" | **必翻**（`Grants` 落进这支）——**与 a-3 是两枚独立红**，解冻时要一起改，只改 a-3 会留下一枚"a-4 红" |
| a-5 | `:469-471` | `if len(got) != len(want)` | **必翻**（6≠5）⇒ **同一枚函数今天会连红三支**〔仅读码，未跑〕 |
| a-6 | `:473-479` | `declared := declaredOptionsFields256(t)`（`:379` 起，读的是 **`gate.go` 里 `type Options struct` 的全字段**）＋ `if len(declared) < 10` | **不翻**（ⓐ-Ⅰ 不动 `gate.go`；本腿现数量到 `declared`＝**10 枚**：`UI/Clock/Channels/Window/ApprovalTimeout/WarningLead/MaxPending/MaxTracked/Logf/Grants`，**正好压在阈值上**⇒ 谁将来删一枚 Options 字段，这枚先红） |

⇒ **对 `A602` 的射程补一句**（编排者那条写的是"期望集 5 枚→6 枚"）：**要动的不止 `:445` 一行**——a-3／a-4／a-5 三支都得跟着搬家，加上函数名与文件头注释（`:43`、`:415-424` 那段"Grants is asserted ABSENT on purpose"）也变成假话。**这仍属同一枚 `A602` 的射程，不新立案。**〔仅读码，未跑〕

### 4.D 族(b) boot 打印／状态句**字面零漂移**钉＝**4 枚，ⓐ-Ⅰ 之下全部"不翻"但全是禁改面**

尺＝`git grep -n 'residentStatusLine\|taskPosture\|statusLine' -- cmd internal tools` ⇒ **42 命中**（`logs/02-families-bcde.txt`）。产码那几枚（定义处）不算钉，钉在测试里，逐枚：

| 枚号 | `file:line` | 它今天断什么 | ⓐ-Ⅰ | 归类 |
|---|---|---|---|---|
| N7 | `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go:107-109` | `if line := ra2.residentStatusLine(); line != old260r3StatusLine` ＝**逐字相等**，期望常量在 `:54`＝`"审批门已装配进本进程（取消通道：Esc 已加载；等待中的确认项：0）"` | **不翻**（本腿现量在飞写面：`git diff HEAD -- resident_approval_windows.go resident_task_source_windows.go` 里含 `residentStatusLine\|taskPosture` 的增删行＝**0 行**，读数档 `logs/15-queue-and-drift.txt`）；⚠ **但它是"ⓐ-Ⅰ 落地时不许把记账状态写进状态句"的唯一一枚钉**——写了就是零漂移红 | 翻了＝**改契约**（撞票 260 的判据面，那批有撤销口令「260 文案撤回」） |
| N8 | 同文件 `:103-105` | `want := "wisp: [audit] " + old260r3LoadedAudit`（常量 `:50-51`）＋ `Contains`＝**装配回执那条 audit 句不许动一个字** | **不翻**（ⓐ-Ⅰ 不印新句；⚠ holder 的 error 路径走的是 `gate.go:689` 那行，不是这条） | 同上 |
| N9 | `cmd/wisp/resident_cancel_key_label_260r4_windows_test.go:113-116` | `status := ra.residentStatusLine()` 必须含配置里那枚键、且**不含 "Esc"**＝"两处渲染同一枚读口"的单源钉 | **不翻**（ⓐ-Ⅰ 不加第二枚读口，也不碰 `cancelKeySpelling`） | 若被撞翻＝改契约（票 260 r4 判据） |
| N10 | `cmd/wisp/resident_task_source_246_windows_test.go:314`／`:385` | boot 那行必须含 `"任务来源："＋taskPostureAbsent／Refused`（四形常量 `resident_task_source_windows.go:107-110`） | **不翻**——⚠ 但**它是 `Contains` 不是相等**，所以"往 boot 那行**加**一段会话记账状态"这件事**它拦不住**；拦得住的只有 N7/N8 那两枚逐字钉。⇒ 本腿把这一格记为**"看加不加字"的射程差**，不许下一任把 N10 当"boot 零漂移"用 | —（不构成 ⓐ-Ⅰ 的阻挡面） |

### 4.E 族(c) `file:line` 型 evidence 钉＝**这一族没有一枚是机器红；且票面点名的"票 255 名册"在 ⓐ-Ⅰ 写面下不成立**

尺（读数档 `logs/03-family-c-mech.txt`＋`logs/05-family-c-narrow.txt`，取数 08:30:38Z／08:35:13Z）：
- `git grep -n 'resident_approval_windows\.go:[0-9]\|resident_task_source_windows\.go:[0-9]' -- docs` → 命中分布在 `docs/evidence/s1/246-*`、`248-*`、`docs/reports/pending-and-issues.md`；
- 同尺打 `-- .scratch` ⇒ **370 命中**（工单＋各家普查件，含本件 §0–§6 那 17 处）；
- 四枚模式（含 `run.go:[0-9]`／`resident_windows.go:[0-9]`）合起来打 `docs`＋`docs/reports` ⇒ **780 命中**＝**这不是能逐枚抄的名册，是一个"全仓注释级引用"的面**。

**三条判语，逐条带凭据：**

1. **没有任何一枚尺机器地核 `file:line`。** `tools/d22scan` 的 ban #9 用的正则在 `tools/d22scan/main.go:890-891`（`repoPathRe`）＝`(?:docs|\.scratch|internal|cmd|tools|scripts)/[A-Za-z0-9_./\-\x{4e00}-\x{9fff}]*[A-Za-z0-9_\-]`，它判的是**路径存在性**（`:872` 才 `add("phantom-citation", ...)`），**行号不是它的射程**；`git grep -n 'docs/evidence' -- cmd internal tools scripts`＝132 命中**全是注释**（`logs/04-family-c-mech.txt` §C5）。⇒ **这族"红"只红在人核／下一枚普查的眼里，卫生四门不会为它响。** 〔仅读码，未跑〕
2. **⛔ 更正票面／编排者起手尺的那句预期**：票 255 那两张名册（`docs/evidence/s1/255-tier-registry-r1.md`／`-v1.md`）里，指向我写面那两枚文件的带行号引用＝**0 命中**（尺＝`git grep -n 'resident_approval_windows\.go:[0-9]\|resident_task_source_windows\.go:[0-9]\|run\.go:[0-9]' -- docs/evidence/s1/255-tier-registry-{r1,v1}.md`⇒空）。255 名册钉的是 `tiers_255_test.go:44/:98/:169/:22/:30` 与 `manager.go:299/:346-367/:369-427`（`-v1.md:33/45/57/66/77/85/90/93/94`）——**全是 config 侧文件，ⓐ-Ⅰ 一枚都不动**。⇒ **"ⓐ-Ⅰ 会让票 255 名册行号漂移而红"这一格＝不成立**；真正会漂的是 **246 那两张**（下表）。〔仅读码，未跑〕
3. **已经漂了，不是 ⓐ-Ⅰ 造成的**（这点必须写，否则下一任会把旧账算到新腿上）。现量对照（HEAD `ba900551` 的产码 vs 名册落的行号）：

| 名册引用 | 落在哪 | 它说的 | HEAD 真身 | 判定 |
|---|---|---|---|---|
| `docs/evidence/s1/246-resident-task-source-v2.md:83` | `newResidentApproval` | `resident_approval_windows.go:105`，调用者 `resident_windows.go:123` | `:170`（定义）／`:213`（WithConfig）＋ `resident_windows.go:132` | **旧漂**（＋65／＋9 行） |
| 同上 `:84` | `startResidentTaskSource` | `resident_task_source_windows.go:215` | 现量 `:265` 那枚 `assembleRuntime` 调用还在同一函数内，入口行号已随 256/260 两批移动 | **旧漂** |
| 同上 `:87` | `ra.vetoByEsc` | `resident_approval_windows.go:171`↔`resident_windows.go:136` | 已随 260 那批重写 | **旧漂** |
| `246-resident-task-source-r2.md:105-106` | MUT-A／MUT-B 突变位 | `resident_task_source_windows.go:215`／`:182` | 突变体行号＝**当时**的行号 | **不可复现型引用**（⚠ 这一类最坏：正控重跑会种错地方） |
| `docs/reports/pending-and-issues.md` 内 10 处带行号引用 | 台账 | 逐条 | — | 台账是 append-only，**按纪律不回填更正** |

⇒ **本腿对这族的净判语＝ⓐ-Ⅰ 会让写面那两枚文件的**行号级引用**全部再漂一次（本腿量到 `resident_approval_windows.go` 在飞增 ＋176 行、`resident_task_source_windows.go` ＋31 行），但：① 没有尺为此响；② 票面点名的 255 名册不在射程；③ 真正该在落地件里补一句"行号已漂，旧引用按 `git show <锚>:` 读"的是 **246 那两张 ＋ 台账 ＋ 本件 §0–§6**。归 ⓐ-Ⅰ 落地腿的证据件，⛔ 本腿不改任何一件。**

### 4.F 族(d) "某词必须不出现"型反向钉＝**3 枚在场 ＋ 1 处本腿量到"零枚在场"**

尺＝`git grep -n 'forbidden\|must not\|不许\|仍在念 Esc\|不再重复提问\|记入本会话\|本会话内允许' -- cmd internal tools docs`＋`git grep -n 'PanelAllowSession' -- cmd internal tools`（读数档 `logs/02-families-bcde.txt`、`logs/06-session-reply.txt`、`logs/17-last-checks.txt`）。

| 枚号 | `file:line` | 它不许什么出现 | ⓐ-Ⅰ | 归类 |
|---|---|---|---|---|
| N11 | `cmd/wisp/resident_approval_246_windows_test.go:373-377` | `for _, forbidden := range []string{"看得见","面板已就绪","已显示卡片"}` 逐枚 `strings.Contains(ra.residentStatusLine(), forbidden)` ⇒ **状态句永远不许读成"有一张卡片看得见"** | **不翻**（ⓐ-Ⅰ 不往状态句加字）；⚠ 它同时是**票面 ⓑ 支的硬拦网**：ⓑ-Ⅱ 若把"这条路上只到本次为止"写进状态句并带上"看得见"这类词＝直接红 | 翻＝改契约（票 246 判据面） |
| N12 | `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go:132-134` ＋ `cmd/wisp/resident_cancel_key_label_260r4_windows_test.go:106-108` | seeded 档下"Esc"必须不在那几句里 | **不翻** | — |
| N13 | `internal/panel/l2_grant_boundary_test.go:1241`／`:1244` | 入向路由名字若是一枚审批决定（`AGENTS.md §1.2` ban #6／D33/F2/R20）⇒ 红；且声明的 guard 集必须与 `knownComposerMethod` 一致 | **不翻**（ⓐ-Ⅰ 不加面板路由）；⚠ 这一枚就是 §2.1 表 #3"面板侧连 Allow 都没有"的**机器护栏本体**——`git grep -n 'PanelAllowSession'` 产码／测试**零枚**（只有 `replies.go:347` 那句注释），**"不许有"这件事是由 N13 的路由名走查兜的，不是由一枚同名测试兜的** | 翻＝改契约（D33 必修五项面） |
| **N14（★零枚在场）** | `cmd/wisp/approval_reply.go:278`（审计 `记入本会话，paths=%d`）＋ `:280`（回执 `…不再重复提问`） | — | 尺＝`git grep -n '不再重复提问\|记入本会话' -- cmd internal tools docs` ⇒ 39 命中里**产码 2 枚、`*_test.go` 0 枚**（其余在 `docs/reports/pending-and-issues.md` 与本件／票件里）。⇒ **那两句"最坏后果"级假话今天没有任何一枚钉**，与 §2.5／票面 AC#1 ★ 条同读数，本腿是**独立复跑后复认**，不是照抄 | ⓐ-Ⅰ 之下它们自动变真话 ⇒ 不需翻；⚠ **但 AC#3 的正控（"种一次写侧没接上 ⇒ 这两句必须能被指出来"）今天无处可种**——要先有新钉。归 ⓐ-Ⅰ 落地腿 |

### 4.G 族(e) 单门指针同一性钉＝**4 枚，ⓐ-Ⅰ 全部不翻，但它们是 ⓐ-Ⅲ 的墓碑**

尺＝`git grep -n 'ThroughTheOneGate\|ar\.gate != ra\.gate\|ra\.gate != \|oneGate\|OneGate' -- cmd internal tools` ⇒ **4 命中**（`logs/02-families-bcde.txt`）。

| 枚号 | `file:line` | 它今天断什么 | ⓐ-Ⅰ |
|---|---|---|---|
| N15 | `cmd/wisp/resident_task_source_246_windows_test.go:89-92` | `if ar.gate != ra.gate { t.Fatalf("…runs on a DIFFERENT gate than the resident leg built…two gates in one process is the shape ruling 2.2 and ledger A481 refuse") }`＝**注入的门必须就是装配根用的门**（指针相等） | **不翻**（ⓐ-Ⅰ 往那枚唯一门里塞 holder，门还是同一枚对象） |
| N16 | 同文件 `:93-96` | `if ar.ui != nil` ⇒ 注入了门以后装配根**不许再造一枚控制台审批面** | **不翻**（`run.go:600-608` 那支未动；ⓐ-Ⅰ 零改动 `run.go`） |
| N17 | 同文件 `:97-100` | `ar.spec.taskCtx != ra.root` | **不翻** |
| N18 | `cmd/wisp/resident_approval_risk_256_windows_test.go:540-547` | `runResident` 体内 `newResidentApprovalWithConfig` 必须**恰好 1 次**、`newResidentApproval()` 必须 **0 次**（AST 走查，函数起 `:498`） | **不翻**（ⓐ-Ⅰ 不动 `resident_windows.go`——本腿现量在飞写面只有那两枚文件＋256 那枚测试，`logs/15-queue-and-drift.txt` §N3）；⚠ **这一枚＋ N15＋ a-1（`:428` 字面量恰好 1 枚）合起来＝ⓐ-Ⅲ"第二枚 mint"在三处各有一枚独立红**，与票面「明确不选 ⓐ-Ⅲ」那条一致：ⓐ-Ⅲ 不只是"形状不好"，它是**机器就拦着**的形状 |

⇒ 本腿补一句给落地腿的**同类不同名**护栏：a-1（`:428-433`）的错误文本**自己指回 N15**，说它是"更早更便宜的版本"。⇒ **ⓐ-Ⅰ 若为了拿 ledger 而新建第二枚 `approval.Options`／第二枚门，红两次**（一便宜一贵），这是好事，写进证据件时别只引一枚。

### 4.H 族(f) 本腿另拉出来、编排者起手尺点不到的两族＝**包级行为台账 2 枚 ＋ 卫生门 1 枚**

| 枚号 | `file:lang`→`file:line` | 它今天断什么（本腿读断言源码，⛔ 没跑） | ⓐ-Ⅰ |
|---|---|---|---|
| N19 | `cmd/wisp/leg_sink_gate_131_test.go:249`（`TestAC4EveryLegIsNailedOrRuled`；枚举器 `:814-817`＝"读 `main()` 的 argv 分派"，站点串形如 `main.go:%d`） | 每一枚 leg 要么被 `registerLegNail131` 钉住、要么在产码文件里带 `WISP-LEG-SINK-RULING:`（常量 `:161`）；`:296-301` 那支＝"记了 record、没装 listener、又没裁定句"⇒ 红；`:307-312`＝"裁定句与 install 并存"⇒ 红 | **可能翻＝本腿量不到**。读到的：leg 来自 `main()` 分派 ⇒ ⓐ-Ⅰ 不加 leg；全仓那枚裁定句只在 `cmd/wisp/slo_windows.go:184`（`logs/09-ruling-markers.txt`）⇒ 常驻那一枚 leg 今天的 state 是 nailed 还是 ruled，**只有跑这枚门读它打印的 ledger 行才知道**。**归口＝新增 U 尺 U13**（见下），⛔ 我不写"不翻" | 若翻＝**修 bug 面内**（本票写面自带的那枚 leg 该有钉），⛔ 不许用加裁定句的方式绕过 |
| N20 | `cmd/wisp/leg_dispatch_gate_133_test.go:178`（`TestAC1AC2DispatchHopGate133`）＋ `:1403`（"nail 必须是本目录源码里一枚 top-level `func TestXxx(t *testing.T)`"）＋ `:1651`（"本编译单元不采的文件里的 Test 函数＝那行是 no witness 的覆盖声明"）＋ `collectRulings133`（`:1744` 起，扫全部**非测试** `.go` 找 `WISP-LEG-COVERAGE-RULING:`，常量 `:138`；现量在场 5 枚：`main.go:74/79/83`、`panel_assets.go:16`、`slo_windows.go:199`） | **ⓐ-Ⅰ 的落地腿会被它咬的位置不是产码，是它新写的那枚测试**：新用例必须①落在本平台真编译的文件名／build tag 下（`_windows_test.go` 形），②若进 ledger 就必须真存在且可跑，③不许把覆盖声明挂在一枚不编译的 case 上 | **可能翻**（取决于落地腿怎么写测试），⛔ 本腿没跑＝**量不到具体那一支** | 翻＝**修 bug**（是新腿自己的判据没写对，不是别人的契约） |
| N21 | `tools/d22scan`（卫生门，非测试钉）：ban #1 裸 `go func(`、ban #9 phantom-citation（`main.go:872`/`:890-891`） | 本腿**读 diff 判**：在飞 holder 只加了一枚 `mu sync.Mutex`（`git diff HEAD -- resident_approval_windows.go` 里 `^\+.*(go func\|sync\.\|filepath\.)` ⇒ 唯一命中＝`+	mu     sync.Mutex`），**没有 `go func(`、没有 `filepath.Clean/Abs`、没有新增 `docs/` 引用** | **不翻**〔仅读码，未跑；真跑属 U6〕 | — |

### 4.I 名册汇总（ⓐ-Ⅰ 作用面）

- **必翻 1 枚（一处函数内 3 支独立红臂）**：N1（`:452`）＋a-4（`:457-468`）＋a-5（`:469`）。三支配一处解冻（`A602`）一起吃掉，⛔ 只改 `want` 一行会留下一支红。
- **零枚在场 1 处**：N14（那两句假话没有尺）⇒ AC#3 正控今天**无处可种**。
- **量不到 2 枚**：N19／N20 ⇒ **本件 §6.1 的 U 尺名册要加两枚**（本腿不重跑前人读数，只登记新格）：
  - **U13**（本腿新增，ⓐ-Ⅰ 专属）＝`go test -count=1 -run 'TestAC4EveryLegIsNailedOrRuled|TestAC1AC2DispatchHopGate133' ./cmd/wisp/` ⇒ 读那两枚台账打印的 leg 行，才知道 holder 那枚新函数进不进 `emits` 分母。归 ⓐ-Ⅰ 落地腿（它跑测试没有禁跑约束）。
  - **U14**（本腿新增）＝`go test -count=1 ./cmd/wisp/ ./internal/agent/approval/ ./internal/panel/` 全量基线，用来证伪本件 §4.D／§4.G／§4.F 里那 12 枚"不翻"判语。
- **不翻 12 枚**：N2／N3／N4／N5／N6（ⓐ-Ⅰ 不碰 `internal/agent/approval` 与 `run.go`）、N7／N8／N9／N10（状态句与 boot 句不写记账字）、N11／N12／N13（反向钉不撞）、N15／N16／N17／N18（同一性面）。逐条凭据在上四节，**全部〔仅读码，未跑〕**。
- **编排者起手尺没点到、本腿新立案的**：N5（`AllowSession` 必须返回 nil ⇒ ⓑ 的隐形代价）、N6（`GRANT-RECORDED` 字面格式钉）、N7／N8（260r3 两枚**逐字相等**零漂移钉）、N9（单源键名钉）、N11（三枚禁字）、N13（`l2_grant_boundary` 那枚路由名走查才是"面板无 allow"的真护栏）、N14（零枚在场）、N15—N18（同一性 4 枚＝ⓐ-Ⅲ 的三处墓碑）、N19／N20（包级行为台账）。**合计新增立案 16 枚**（N5—N21 减去与前人 §1.3 重复的 N2／N3／N4）。
- **本腿对票面的一处更正**：§4.E-2——"票 255 那张名册会因行号漂移而红"在 ⓐ-Ⅰ 的写面下**不成立**（255 名册零引用写面文件）；会漂的是 246 那两张＋台账＋本件。


## 5. 顺带上报（看到就报，⛔ 不许顺手改）

（本节答：（a）别的同类接缝——`Options` 里还有哪枚字段在某条生产腿上留空、后果是什么、有没有尺；
 （b）`GRANT-RECORD-FAILED`／`GRANT-RECORDED` 有没有尺在断（尺自拉）；
 （c）票 224／256／260／266 里与票 265 现量**冲突**的句子，指到行号，⛔ 不改别的票的票面。）

### 5.A （a）`Options` 里还有哪枚字段在某条生产腿上留空＝**四条腿同形留空 4 枚，但只有 `Grants` 这一枚没有兜底**

尺＝`git show HEAD:internal/agent/approval/gate.go` 取 `type Options struct` 全字段 ＋ `git show HEAD:cmd/wisp/run.go \| sed -n '605,625p'` ＋ 两枚产码字面量逐字段对表（读数档 `logs/13-sec5-ab.txt` §L1/L2、`logs/14-sec5a-fields.txt`、`logs/16-clock-status.txt`）。

**全集 10 枚**（`gate.go` 现量）：`UI / Clock / Channels / Window / ApprovalTimeout / WarningLead / MaxPending / MaxTracked / Logf / Grants`。
**两枚产码字面量各传几枚**：`run.go:612`（`wisp run` 腿）＝**6**（`UI/Channels/Window/ApprovalTimeout/Logf/Grants`）；`resident_approval_windows.go:219`（常驻腿，HEAD）＝**5**（少 `Grants`）。⇒ **两枚生产腿都留空的是同一批 4 枚：`Clock`、`WarningLead`、`MaxPending`、`MaxTracked`。**

| 字段 | 谁留空 | 后果（读码） | 有没有尺 |
|---|---|---|---|
| `Clock` | **两条生产腿都留空** | **无后果**——`gate.go:127-129` 有兜底：`if clock == nil { clock = SystemClock{} }`（那一行上面的注释逐字 "never a downgrade"）。⇒ 生产拿系统时钟，测试拿 fake | 只在包内有尺：`internal/agent/approval/ticket84_no_owner_test.go:59`、`:154`、`ticket87_veto_l2_test.go:49`、`ticket97_alias_direction_test.go:32`、`ticket220_l1_window_read_test.go:203/:383`、`ticket224_reply_grant_test.go:120` **全是包内自造 Options**。⛔ 没有一枚尺钉"生产腿该不该传 Clock" |
| `WarningLead` | 两条都留空 | **无后果**——`queue.go:88-90`：`warnBefore <= 0 \|\| warnBefore >= timeout` ⇒ `DefaultApprovalWarning` | 同上：只有包内 `batch_test.go:90/:123`、`queue_test.go:192`、`ticket84:61`、`ticket87:51` 显式传值。**产码侧零枚尺** |
| `MaxPending` | 两条都留空 | **无后果**——`queue.go:91-93` ⇒ `DefaultMaxPending` | 包内 `queue_test.go:335`（`Options{MaxPending: 1}`）。产码侧零枚 |
| `MaxTracked` | 两条都留空 | **无后果**——`gate.go:152-154` ⇒ 常量 `64`（**兜底值是枚字面量，不是 `Default*` 常量**，⚠ 这一形比另外三枚脆：将来谁把 64 换成常量，那枚换法没有任何产码尺在盯） | 零枚（`logs/14-sec5a-fields.txt` §M1 只命中 `fakes_test.go:269`，是仪器自己的拷贝表） |
| **`Grants`（对照面）** | **只有常驻腿留空** | **有后果，且是唯一一枚没兜底的**：`gate.go:155-165` 的 `New` 对上面 4 枚都有 `if`，**对 `Grants` 只有 `grants: o.Grants,` 一行直传**（`:160`）⇒ nil 不在构造期被换成任何东西，而是**推迟到答复那一刻**才由 `:673` 那一支兜——兜出来的是"放行但不记"＋一行日志。⇒ **本票的整条因果线就是"这一枚 Options 字段没有 default 分支"这件事**；如果 `New` 当时给 `Grants` 也写了一枚兜底（哪怕是"记不了就报错"），本票不会存在 | 见 §5.B |

⇒ **给编排者的一句话（不是问题，是形状）**：ⓐ-Ⅰ 落地之后，`approval.Options` 的 10 枚字段里，"两条生产腿都不传"的还剩 **4 枚**，其中 **3 枚有 `Default*` 兜底、1 枚（`MaxTracked`）有裸字面量 64 兜底**；**"一条腿传一条腿不传"的还剩 0 枚**。⇒ 这一格以后不会再长出本票这一形，**除非**有人给 `New` 加第 11 枚字段而忘了兜底分支——那件事的唯一哨兵是 §4.C a-6（`resident_approval_risk_256_windows_test.go:476` 的 `len(declared) < 10`），而**它的哨兵是"枚数"不是"有没有兜底"**，⛔ 它拦不住"加了字段没加兜底"。〔仅读码，未跑〕

### 5.B （b）`GRANT-RECORD-FAILED`／`GRANT-RECORDED` 有没有尺在断＝**RECORDED 12 命中／RECORD-FAILED 在 HEAD 只有 1 命中，且不在 `cmd/wisp`**

尺（本腿自拉，读数档 `logs/13-sec5-ab.txt` §L3、`logs/18-sec5c-conflicts.txt` §Q2）：
- `git grep -n 'GRANT-RECORDED' -- cmd internal tools` ⇒ **14 命中**，其中 `*_test.go` **12 命中**；
- `git grep -n 'GRANT-RECORD-FAILED' -- cmd internal tools` ⇒ **7 命中**，其中测试 **2 命中**——⚠ 但那 2 命中里**有 1 是在飞的 `265-r1` 工作副本**（`cmd/wisp/resident_approval_risk_256_windows_test.go:473` 的**错误文本**，不是断言）；
- ★ 决定性那一发＝`git grep -n 'GRANT-RECORD-FAILED' HEAD -- '*_test.go'` ⇒ **恰好 1 命中**＝`internal/agent/approval/ticket224_reply_grant_test.go:308`（`if !f.log.has("GRANT-RECORD-FAILED")`）。**枚数＝1，且它在 `internal/agent/approval` 包里、用的是包内 fixture。**

⇒ **三条上报，逐条给后果：**

1. **ⓐ-Ⅰ 的第一条硬约束（票面「落地三条硬约束」第 1 条："holder 在还没绑到 ledger 那段窗口里 `Record` 必须走 error 路径"）在 HEAD 上零枚产码侧尺**。包内那一枚 `:308` 断的是"recorder 的 `Record` 返回 error ⇒ 落 `GRANT-RECORD-FAILED`"，**它不知道常驻腿有没有 holder、也不知道 holder 有没有绑定**。⇒ **AC#3 的正控（"种一次 holder 未绑定 ⇒ 指名用例红"）今天无处可种**，必须随 ⓐ-Ⅰ 新建；这与 §4.F N14 那两句假话同格——**票面自己要求的那枚正控，正控对象今天没有仪器**。〔仅读码，未跑〕
2. **`GRANT-RECORDED` 的字面格式被两枚尺钉死**（`cmd/wisp/ticket224_assembly_test.go:189-190` 与 `internal/agent/approval/ticket224_reply_grant_test.go:219`/`:318`，后者还钉到 `grant_id=11`）⇒ ⓐ-Ⅰ **不许为了让审计"看得出是哪条腿"而改 `gate.go:693` 那行的字段顺序或加后缀**。这格是**给落地腿的红线**，不是缺陷。
3. **`GRANT-RECORD-FAILED` 与 `GRANT-RECORDED` 的极性在 `cmd/wisp` 侧只有一枚反形尺在场**（`ticket224_assembly_test.go:194` 只反钉 `GRANT-DROPPED`，**没有一枚反钉 `GRANT-RECORD-FAILED`**）⇒ 推论：**"装配正确"的那枚产码尺对"记账写失败"不敏感**。种一发"ledger 已接、但 `Record` 恒返回 error" ⇒ 本腿读码判 `cmd/wisp` 侧**零枚红**（`internal/agent/approval` 侧 `:308` 会绿，因为它要的正是这句）。⚠ 这一格是**量不到的那一种**：判它要么跑 §4.I 的 U14，要么由 ⓐ-Ⅰ 落地腿补一枚反形钉。归口：票 265 AC#3。〔仅读码，未跑〕

### 5.C （c）票 224／256／260／266 与票 265 现量冲突的句子（指到行号，⛔ 一枚票面未改）

| # | 别票的那一句（`file:line`） | 逐字在说什么 | 与票 265／本件的哪一句撞 | 本腿判 |
|---|---|---|---|---|
| **C-1** | `.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md:1` ＋ 同文件 `:19` | 标题句"`approval_grant` 没有生产写手" ＋ 表格行"⛔ **生产零写手／零读者**：那几枚 DAO 方法的非测试调用者＝0" | 票 265 `:4` 那句"它 09-29 那句'生产零写手'已被 224-r1/r2 部分推翻——今天 `run.go` 确实接了两枚" ＋ 本件 §2.3 现量（`run.go:618`／`:758`） | **冲突成立且是过期不是错判**：265 已具名说它"部分推翻"，但 **224 票面那两行没跟着加更正标注**（票面纪律＝append-only ⇒ 只有 224 自己或编排者能补）。⇒ 下一任读 224 票面会再得到一次"零写手"的错误先验。**归口：224**，本腿不动 |
| **C-2** | `.scratch/wisp/issues/256-resident-gate-built-before-session-grants.md:8` | "`resident_approval_windows.go:108-113` 逐字 `ra.gate = approval.New(approval.Options{` 里只有 `UI`／`Channels`／`Logf` 三项" | 本件 §1.2（现量＝5 枚 `UI/Channels/Window/ApprovalTimeout/Logf`，HEAD 行号 `:219-224`）＋ 票 265 §现量.2 | **256 票面自身已过期**（256-r1 落地把它从 3 枚变 5 枚）；256 自己在 `:58` 只更正了"区间起点偏一行"、**没更正"三项"**。⇒ 属"票面未回填"，不属矛盾。**引用 256 §现量.1 前必须按 `git show <锚>:` 读**，这一格与本件 §4.E-3 那批行号漂同族 |
| **C-3** | `.scratch/wisp/issues/265-resident-gate-has-no-session-grant-writer.md:53` | "AST 钉搬家需具名解冻＝**具名解冻已落 `A601`（五样齐＋撤销口令）**" | 台账：`docs/reports/pending-and-issues.md:11725`（A601 正文那句"AST 钉搬家需具名解冻＝**见 `A602`**"）＋ `:11738`（**`A602` 才是**"具名解冻·五样齐：期望集 5 枚→6 枚"，口令「265 撤 AST 钉解冻」在 `:11744`）＋ 另票旁证 `.scratch/wisp/issues/256-...md:24`（同一句写的是 `A602`） | ★ **票 265 面自己那一处是笔误（`A601`→应为 `A602`）**，三处独立凭据同向。⛔ 本腿不改票面（AC 框与票面归编排者）；**落地腿照 `A602` 的边界办，不要照票面那一行的 `A601`**——两枚的撤销口令不同，读错会撤错东西 |
| **C-4** | `.scratch/wisp/issues/256-...md:92` | 定式：**"禁现那格 ⛔ 不许钉前缀 `approval: GRANT-DROPPED`（前缀两形同吃＝票面 §7-6 已定的洞；`ticket224_reply_grant_test.go:283` 用的正是前缀，**不可复用、也不许为让前缀零命中去改它**）"**，并指定该钉的独有词组是 `gate.go:675` 那句「本机没有接入会话授权记账」 | 票 265 §现量.3（`:10`）把 `:283` 那枚当"既有尺两枚"之一引用 ＋ 本件 §4.B N2 与 §4.C 的起手尺族 | ★ **真冲突，且冲突在 265 的引用面**：265 用一把 256 已判"不可复用"的尺当现量凭据。本件照实复认了它的读数（它确实断"这句在场"），但 ⚠ **ⓐ-Ⅰ 落地时若要写"常驻腿上不许再出现 GRANT-DROPPED"这类禁现格，必须按 256 那句定式钉独有词组、不许钉前缀**——否则一票的判据面被另一票的定式否掉。**归口：编排者**（要不要在 265 面补一句指回 256:92） |
| **C-5** | `.scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md:63` | 撤销口令「**260 文案撤回**」的射程＝"给人看的那句话"，凭据＝`A595` 具名解冻五样齐（文件／行／理由／边界／口令） | 票 265 面 `:24` 那句"⚠ 文案射程刚被票 260 那批钉过" ＋ 本件 §4.D N7/N8（`resident_cancel_key_wording_260r3_windows_test.go:107`／`:103`）**与 N14（`approval_reply.go:278`/`:280` 那两句）** | **不冲突，是补边界**：260 那批钉的是**取消键那五句**，而 ⓐ-Ⅰ 之后可能被动的是 `approval_reply.go` 那两句**同族但不在 260 射程内**的人看句。⇒ 本腿量到：**那两句今天零枚钉（§4.F N14），所以它们既不在 260 的解冻面里、也不在任何禁改面里**——动它不需要 260 的口令，但 ⛔ **也不要有腿拿"没人钉"当"可以不改"的理由**（AC#1 ★ 条要求两句都处理） |
| **C-6** | `.scratch/wisp/issues/266-no-external-tooth-over-the-slo-check-self-locking-nails.md:52` | 禁改清单里"三枚冻结件"**具名**＝`internal/panel/tokens_fourway_test.go`／**`internal/panel/l2_grant_boundary_test.go`**／`internal/perm/ticket90_persist_test.go` | 票 265 `:35`（AC#2）只写"三枚冻结件"**没列名** ＋ 本件 §4.F N13（那枚 `l2_grant_boundary_test.go:1241/:1244` 恰是"面板侧不许有 allow 位"的真护栏） | **不是冲突，是一枚具名缺口**：265 AC#2 让腿"不动三枚冻结件"却没说哪三枚，而其中一枚正是本票 ⓑ／I6 那条"要不要把 drop 做成界面事件"路上**唯一拦得住面板侧 allow 的尺**。⇒ 建议（⛔ 本腿不动票面）：**在 265 面补一句指回 266:52 那三枚名**，否则 ⓐ-Ⅰ 之后再来一枚 ⓑ 的腿会去动那枚冻结件 |
| **C-7** | `.scratch/wisp/issues/256-...md:25` | "`GRANT-DROPPED 必须不再出现`那一半随 `Grants` 未接而**根本无法构造**⇒ 同归票 265" | 票 265 AC#3（`:36`）只写了三形各自的正控形状，**没接住 256 归过来的这一格**；本件 §5.B-1 量到该正控今天无处可种 | **归属撞车**：256 把一件事归给 265，265 面没有对应判据行。⇒ ⓐ-Ⅰ 落地腿若顺手把它做了，是"做了没被派的事"；若不做，那格在两张票之间悬空。**归口：编排者** |

### 5.D （d）三件不在 (a)(b)(c) 档里、本腿读码撞到的

1. **typed-nil 陷阱：`Options.Grants` 这一枚接法今天没有哨兵。** `gate.go:160` 是 `grants: o.Grants` **直传**，而 `:673` 判的是 `g.grants == nil`——**接口值 nil 与"装着 nil 指针的非空接口"不是一回事**：若 holder 以 `*residentGrantHolder` 形式接进 `Grants` 而**构造器里没 `new`**，那 `g.grants != nil` ⇒ 不走 drop 支 ⇒ 直接进 `:684` 对 nil receiver 调 `Record`。**本腿在飞工作副本里读到落地腿已经避开了它**（`resident_approval_windows.go` 工作副本 `:365` 先 `ra.grants = &residentGrantHolder{}`，再在 `:374` 传 `Grants: ra.grants`；⛔ 本腿不评价那枚件，只登记"这一格是靠读码撞到的、不是有尺的"）。**仓里零枚尺钉住"不许把 typed-nil 传进 Grants"**——它红不红只有 `go test` 知道。归口：票 265 AC#3。〔仅读码，未跑〕
2. **前人 §1.2 那句"129 已漂成 128"这一格，本腿量到"两个数都还在票面上"。** 现量（HEAD `f9dc152c`）：`resident_approval_windows.go` 的 HEAD 版本 `:188` 注释仍逐字写着 "gate is built roughly 129 lines BEFORE"；而 §4.E 的漂移表说明**门字面量已从 `:219` 挪到工作副本 `:368`** ⇒ 那句注释连同"129/128"之争**在 ⓐ-Ⅰ 之后两数皆错**。⛔ 不改别的件、不改票面；归 ⓐ-Ⅰ 落地件（它本来就要改写 `:185-192` 那一段，见 §3.1）。
3. **`.scratch/wisp/probes/265/r1/` 在盘上（未跟踪）且 `evidence.md` 已引用 `resident_approval_windows.go:<line>`**——尺＝`git grep -n 'resident_approval_windows\.go:[0-9]' -- .scratch/wisp/probes/265` ⇒ **1 命中**（`265/r1/evidence.md`）。⇒ **落地腿的凭据件自己也带行号引用，而它在 ⓐ-Ⅰ 自己那发之后必然再漂一次。**与本件 §4.E 同族，⛔ 不是缺陷，只是"下一任读它要带锚点"。


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

## 〔编排者更正 2026-10-04 14:31:09 +0800〕上一条标注里那句"§0–§6 正文已写满"**是我写错了**，按盘上现量更正（⛔ 上一条一字不改，只在此追加）

- **现量**（我 14:1x 自己跑 `grep -n "^## " census.md` ＋ 逐节读）：**§0／§1／§2／§3／§6 正文写满；§4 与 §5 只有标题＋那句"（本节答：…）"意图句，正文 0 行；§7 同形**。⇒ 上一条那句"§0–§6 正文已写满"**不成立**，按本条读。
- **漏掉的正是 AC#0 点名的两格交付物**：§4＝**撞钉名册**、§5＝**顺带上报**——票面 AC#0 逐字要求"会撞到的既有钉名册"，所以这一格**不是装饰性缺节，是 AC#0 本身没交完**。⇒ 编排者裁定：AC#0 **不翻勾**，缺口另派接续腿（见票面与 `A601`）。
- **为什么我会写错（记我，不当笔误藏）**：我代提之前只跑了两把尺——`wc -l`（315 行）与占位符尺（`待[填]|填写[中]`＝**0 命中**）——就把"写满"签了上去。⚠ **这两节并没有留"待填"两个字，它留的是"本节打算答什么"这一形 ⇒ 我那把占位尺对这一形全瞎，而行数照不出来（意图句本身就有三行）**。
- **由此补一条定式**：判"这一节写完没有"要跑 `grep -c "本节答" <件>`（命中＝只有意图句、没正文），⛔ 不许拿 `wc -l` ＋"占位 0"当完整性判据。
- **⚠ 这条对我今天已经签过的账有没有波及——我现跑了，没有**：同一把尺打过我今天据以翻勾的三份件：`.scratch/wisp/probes/260/v1/verdict.md`＝462 行／**本节答 0 处**、`263/v1/verdict.md`＝304 行／**0 处**、`212/v3/verdict.md`＝359 行／**0 处**（另 `266/a1/census.md`＝180 行／0 处）。⇒ **那三枚满稿的判语不受本次更正影响，我不做"连它们一起降档"的过度更正**；坏的是我那把尺，不是那三份件——**但今后任何一枚件的"满稿"我必须连这把尺一起跑**。
- 本件 §1／§2／§3／§6 的读数与结论**照收**（它按派单逐条标了〔仅读码，未跑〕）；§4／§5／§7 三格不作凭据用途。
