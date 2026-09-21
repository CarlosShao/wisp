# 90 — 面向用户的**权限模式开关**（三档：只问高危 / 每步都问 / 全自动，且"全自动"也吃不到不可逆操作）

**Status:** **ready-for-review（2026-09-21，agent-ticket90b 接续收尾）** —— 五框 AC#1–AC#5 现在**每一框都有可重跑的读数**，
**AC#4 那枚此前"涂了勾没跑"的框已由 5 轮变异补上证据**（见本票最后两条 Progress log）；代码 SHA：`1d4f289`+`582b9a9`+`d5564c2`。
**next=编排者验收**。前：**in-progress（2026-09-21，agent-ticket90）** —— 已按 owner 定案 M1–M5 + R20 选定施工形状（见本条下面第一条 Progress log），
开始落码。前：**unblocked**（2026-09-21 17:0x，owner 已把 M1–M5 全部拍完，见裁定 **R20** 与本票末尾"owner 定案"条）。UI 落点那半（输入框显示档位 + 附件 + 工作区选择）**拆给票 92**，本票只做权限语义、持久化、切档确认与审计。
**Type:** 安全能力 + 可用性（owner 2026-09-21 明确要求："有没有设计类似的权限切换的功能？这个也是很基本的，我希望都有"）
**Blocks:** 票 21 的 segment 2（原生确认卡）、票 48（审批队列）、票 49（会话授权）——**这四张是同一条链**
**Blocked by:** 票 80/83 已查清的前提：**`risk.Gate` 在生产里零调用点**（真实链走的是 `risk.Classify`），
所以"模式"这个东西今天**没有可挂的钩子**。

## 现状（我查过代码与契约，不是印象）

**已经有的（真在跑）**：
- **风险分级**：`internal/risk/`（`assessor.go` + `rules_gateway/shell/network/irreversible/scale`）
  决定一次工具调用是 L1（2–3 秒提示窗，可打断）还是 L2（原生强确认卡）。**这条链是活的**。
- **路径与来源防线**：C26 `PathResolver`（拒 reparse/junction、8.3 展开、句柄真实路径）、C25 污染/溯源、
  D31 原子改名、D34 规则分离（读写不同门）。
- **进程树治理**：C30 **Job Object**（所有子进程入 Job、`KILL_ON_JOB_CLOSE`、按进程树量私有内存）。
- **会话授权的数据层**：`internal/memory/models.go:105 GrantScopeSession = "session"`（SPEC-02 §3 的表），
  但**端到端还没接**（票 49 `ready-for-agent`）。

**没有的（今天确认的缺口）**：
1. **面向用户的模式开关**——你截图那种"询问审批 / 自动审批 / 完全访问"三档，我们**一档都没做**。
2. **OS 级隔离**（受限令牌 / AppContainer / 低权限账户跑子进程）——契约里**一个字都没有**，
   只有 Job Object（那是**生命周期与度量**，不是安全边界）。⇒ 另见票 91。
3. **静态豁免通道**：`blacklist_overrides` 这类键今天被证明是死的，已按票 83 改成**加载即报错**。

## 契约已经划死的红线（模式开关**不许**越过，逐条都有出处）

- `PLAN.md:1629`：不可逆操作**必须拒绝，且任何授权与会话授权都无法覆盖**。
- `PLAN.md:1640`：会话授权**不得覆盖 C25 污染升级**；会话结束授权**必须失效**。
- `PLAN.md:1588`（D33/F2）：面板 XSS 是"唯一能一次性绕过全部门控的路径" ⇒
  **「允许」只接受原生侧点击**，前端与任何 IPC 文本通道都不能代表用户批准（`ban #6` 就是这条的机器门）。
- `PLAN.md:1797`：曾有人提"面板可调用 `approval.decide`"，**被推翻**。
- `PLAN.md:3143`：**"L2 无限等待"被明文驳回**（永挂并持有 C20 路径锁）⇒ 超时一律判拒绝（300s，C18）。

## 待拍板（owner 一句话就能答完；括号里是我的推荐）

**M1 模式集合**：**推荐三档** = `每步都问`（默认，L1+L2 全开）/ `只问高危`（L1 静默、L2 仍问）/
`全自动`（L1+L2 静默放行，**但不可逆、写工作区外、外网三类永远仍问**）。
不推荐做"完全访问"这种**一键关掉一切**的档——契约里那三条红线不允许它存在。
**M2 默认档**：**推荐 `每步都问`**（自用期先建立信任，跑顺了再降）。
**M3 生效范围**：**推荐"会话级"**（改档只影响当前会话，重启回到默认），而不是持久写进配置文件——
理由是 `PLAN.md:1640` 那条"会话结束必须失效"的精神。
**M4 切档本身要不要确认**：**推荐要**（从"每步都问"往下调，是一次 L2 强确认 + 写审计日志）。
**M5 UI 落点**：**推荐球上右键菜单 + 面板顶部一条当前档位**（不做独立设置页，那是票 39 的事）。

## AC（等 M1–M5 点头后开工；先备着）

- [x] **AC#1** 模式作为**显式参数**进入决策链（`risk.Classify` 的调用方），**不许**做成包级全局变量：
      全局变量等于"任何代码路径都能改权限"，那是比没做更糟。
- [x] **AC#2** 三条红线各一条用例：不可逆操作 / 污染升级 / 前端代表用户批准 ⇒
      **在任何模式下都必须被拒或被拦**（`ban #6` 那条同时要有静态门）。
- [x] **AC#3** 切到"全自动"档要 L2 强确认 + 写审计（**审计三档都写**，只有全自动那一次要确认）。
      ⚠ **2026-09-21 按 R20/M3 更正本框**：原判据写的是"重启后回到默认档"，**owner 已推翻这一半** ⇒
      现在要的是**两条分开、不许合并**的用例：
      (a) **改过档位 ⇒ 重启后读回 = 那一次手动选的档**（不是默认档）；
      (b) **从未手动改过 ⇒ 重启后读回 = 默认第一档"每步都问"**。
      另见下面"owner 定案"条里新增的 **AC#3b**（会话授权**不能**跨重启，与 (a) 正好相反，两条都要在）。
- [x] **AC#4** 双向变异：把"模式"接进链但**去掉红线判定** ⇒ 必须有用例红；
      把红线改成"模式优先" ⇒ 必须红。锚点=承载行为的那一行，同链 grep 自证。
- [x] **AC#5** 台账：`risk.Gate` 从"生产零调用点"变成**有真实调用者**，或明写它仍归票 21 未建（不许假装它存在）。

## Rules（本仓固定）

冻结面禁改（`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、`pathresolver*.go`、`rules_gateway.go`）
⇒ **本票若要动这些，先停下来上报编排者**（D22）。其余按仓内固定规矩（引号 heredoc、显式路径、不 push、
不在仓内建 worktree、四种假绿逐跑点名、票面 append-only、渐进写）。

## Progress log（append-only）

- 2026-09-21（编排者）：建票。owner 提出"沙箱 + 权限切换都要有"。
  我查完的现状：**权限切换的"决策层"有、"用户可切的档"没有；OS 级沙箱没有（只有 Job Object）**。
  M1–M5 已按我的推荐预填，等 owner 一句话。

- 2026-09-21 17:0x（**owner 定案，编排者原样入账 ⇒ 施工以这条为准，别再翻上面 M1–M5 的推荐**）：
  - **M1 = 三档**（每步都问 / 只问高危 / 全自动）。不做"一键关掉一切"。
  - **M2 = 默认第一档"每步都问"**。
  - **M3 = 持久化（⚠ 这一条 owner 推翻了我的推荐）**：手动选过哪档，之后一直按那档走，
    **新会话与重启都不回默认**。
    ⇒ **我必须钉住一条边界，防止它渗到别处**：M3 说的是**"模式"这一个全局偏好持久化**，
    **不是把"会话授权"也变成持久**。`PLAN.md:1640` 那条"会话授权不得覆盖 C25 污染升级、
    会话结束必须失效"**一字不动**，票 49（`GrantScopeSession`）仍按"会话级、重启即失效"实现。
    ⇒ 新增判据 **AC#3b**：**改模式能跨重启，会话授权不能跨重启**，两条各一条用例，不许合并成一条。
  - **M4 = 只有切到"全自动"需要确认**（一次 L2 强确认）；**但审计日志三档切换全都写**
    （从哪档到哪档、时间、来源、谁触发的）。
  - **M5 = 档位显示在主面板的输入框里**（不另做设置页；参考 owner 给的截图形状：左下"添加附件"、
    输入框内可放模式选择器、右侧模型选择与发送）。**要**：附件、工作区选择。
    **不要**：git 分支切换那类（owner 原话：我们产品不要求 coding 能力）。
    ⇒ UI 那半**拆成票 92**（主面板输入框 v2）。本票交付：三档语义、持久化、切档确认与审计、
    红线不可越，以及把 `risk.Gate` 从"生产零调用点"变成有真实调用者（AC#5）。
  - **红线复述（owner 这次没改，照旧不可越）**：不可逆操作任何档都必须拒/问（`PLAN.md:1629`）；
    污染升级不能被任何授权覆盖（`:1640`）；"允许"只接受原生侧点击（`:1588`，`ban #6` 是它的机器门）；
    L2 不许无限等待，超时 300s 判拒（`:3143` + C18）。

- 2026-09-21（agent-ticket90，**checkpoint #1：形状选择，先落盘再写码**）：已逐字读完票面（含最后两条 Progress log）
  与 R20 的边界条款。施工形状定如下：
  1. **档位类型**：`internal/risk/mode.go`（**新文件、非冻结**；冻结的 `assessor.go` 一字不动）。
     `ModeAskEveryStep`（默认）/ `ModeAskHighRisk` / `ModeAuto` + `ParseMode` + `Mode.Screen(d Decision) (Decision, bool)`。
     类型放 `risk` 是为了让 `config` 与 `tools` 都能引它而不形成环（`risk` 不引 `config`）。
  2. **进链方式（AC#1）**：`internal/tools/bridge.go` 的 `route` 改成收**显式参数** `mode`，
     值由 `Options.Modes ModeSource`（`Mode() risk.Mode`）在每次 `Execute` 开头读一次；
     **没有任何包级可变全局**，未注入 = 永远最严档（fail-closed）。
  3. **持久化键（M3 / AC#3a）**：`config.toml` 的 `risk.permission_mode`（[risk] 是 locked section），
     默认 `"ask_every_step"`；未知值**加载即报错**（票 83 的规矩），并在 `unwired.go` 的
     `lockedKeyDisposition` 与 `manager.go` 的 `riskDirection` 里挂号（往全自动走 = loosen）。
  4. **M3 边界（编排者加的）**：只有**模式这一个偏好**进 config.toml 从而跨重启；
     **会话授权不碰这条通道**——`PLAN.md:1640` 一字不动，票 49 的 `GrantScopeSession` 保持"重启即失效"。
     ⇒ **AC#3(a)/(b) 与 AC#3b 是三条独立用例，不合并**。
  5. **M4 确认 + 审计**：新建 `internal/perm`（**新包，不踩别人的地界**）持有 `Store`：
     `Set(m, origin, actor)` ⇒ 只有目标档是 `ModeAuto` 时走一次 **L2 强确认**（复用 `tools.Gate.PendingApproval`，
     超时沿用 `risk.confirm_timeout_sec` 默认 300s，**这个数不动**）；**三档切换全部写审计**
     （from/to/时间/来源/谁触发），落在**现有 audit sink**（`Options.Logf` 那一族 `[audit]` 行）上，**不新建第二套**；
     另留 `Mode()` / `LastSwitch()` 给票 92 以**只读**形式显示。
  6. **红线判定（AC#2）**：`ModeScreen` 的静默白名单是**排除式**写的——
     `Deny`、R3、R4（`SessionOverrideBlocked`）、R8 不可逆、R2 工作区外、R5 外网、R9/未知等级 **一律不许被任何档静默**。
  7. **AC#5 台账**：`risk.Gate` 今天在生产里**仍是零调用点**；本票给它一个真实调用者
     （B 档的"单文件已确认"判定走 `Gate`，`internal/tools/paths.go`），
     若最终没落地就在此明写"仍归票 21 未建"，**不假装它存在**。
  8. **禁区自查**：不碰 `internal/winsec/**`、`frontend/**`、`internal/panel/`、`cmd/wisp/`、`internal/risk` 既有测试文件、
     `assessor.go`/`pathresolver*.go`/`rules_gateway.go`/`tools/d22scan/**`/`allowlist.txt`/`docs/PLAN.md`/`docs/specs/*.md`。
     ⇒ **已知代价**：`cmd/wisp/run.go` 那行 `Modes:` 注入由票 77 的文件承载，本票只把**需要的改动缩到一行**并在票面登记交接。

- 2026-09-21（**agent-ticket90b，接续说明 + AC#4 落地的账**）：**我从哪一步接**：前任 `agent-ticket90` 撞 150 轮上限，
  最后一句话是 "Now AC#4: bidirectional mutation testing in a fresh out-of-repo snapshot"，**然后就没了**。
  ⇒ **AC#4 那个框此前为什么是勾着的：前任把框涂满了却没跑变异**（全文检索该框，票面里 "AC#4" 只出现在框本身，
  log 里零次变异读数/红名/rc）。本条把那笔账补上（下面第二条），**不粉饰**：那一勾在此之前**没有证据支撑**。
  接到的工作树实况：`internal/risk/mode.go`(新)、`internal/config/permmode.go`(新)、`internal/perm/`(新包)、
  `internal/tools/mode.go`(新) 与 `bridge.go:285-292` 的 `mode := b.permissionMode(); sil := mode.Screen(verdict)` 已在，
  测试 `internal/tools/ticket90_test.go` / `internal/perm/store_test.go` / `internal/perm/ticket90_persist_test.go` 已在
  （已提交：`1d4f289` 落码 + `582b9a9` 落测试）。**未提交**的剩两处，我核完语义后代为落 `d5564c2`：
  ① `ticket90_persist_test.go` 里 AC#3b 那条**自带了一段"模式确实跨重启"的对照**——**这正是 R20/M3 边界条款禁止的合并**
  （合并 = 一条用例同时判两件事 = 会话授权那条哪天依赖上模式的机制就悄悄绿了），已删掉对照段，
  **三条用例现在各自独立**：`TestTicket90ManualSwitchSurvivesRestart`＝AC#3(a)、
  `TestTicket90UntouchedConfigStartsAtTheDefault`＝AC#3(b)、`TestTicket90SessionGrantDoesNotSurviveRestart`＝AC#3b，
  三条各自 `-run` 单跑各自 PASS（`=== RUN` 3 条、0 FAIL），AC#3b 那条函数体内**再无任何 `Mode` 引用**（grep rc=1）。
  ② `config/unwired.go` 里 `risk.blacklist_overrides` 的台账文案按实测改写（`Gate` 已有真实调用点，缺的是填 `bOverrides` 的确认记录）。
  **依赖/调用关系复验（票 94 的教训，不照抄）**：checkpoint #1 第 1 条说 "`risk` 不引 `config` 所以不成环"——
  `go list -deps ./internal/risk/ | grep -c wisp/internal/config` = **0**（rc=1），`internal/risk` 的仓内 import 只有
  `observe` 与 `winsec`，`config`/`tools`/`perm` 三包各命中 risk 1 次 ⇒ **这条断言为真**（与票 94 那条"winsec→risk 无环"不同，
  这条我实测过才留）。

- 2026-09-21（**agent-ticket90b，AC#4 双向变异原始读数；全程仓外纯净快照，工作树零改动**）：
  快照 `/tmp/wisp90b-sess90b`（`git archive HEAD | tar -x`，`git rev-parse` 确认**非 git 仓**=纯净），HEAD=`d5564c2`。
  **变异前**：`go build ./...` rc=0。基线四包 `go test -count=1`（risk/tools/perm/config）**rc=1**，唯一红是
  `TestResolvePerCallBudget`（C26 性能门，1.165534ms/op vs 1ms 预算，负载敏感）——单独 `-count=2` 两次 **ok**，
  随后**未变异的 `-v` 全量跑 rc=0**（PASS 222 / SKIP 1 / FAIL 0）⇒ 判定：**与票 90 无关的性能抖动，阈值一个字没动**，
  下面每轮都拿 222 当参照核对加减（四种假绿逐条点名：**唯一的 SKIP 是 `TestSyncRegistryProbeLive`，变异前后都是 1 条，
  没有任何 PASS 被它冒充**）。锚点=**承载行为那一行**，每轮**同一条 `&&` 链里 grep 自证变异真的落地**，每轮 `go build ./...` 先量到 rc=0。
  - **(i)-a 去掉"不可逆"红线**：`sed '/case R8:/{N;d}'`（`internal/risk/mode.go` `redLine()`）。
    PROOF：`grep -n "case R8:"` rc=1（空），同函数 `case R5:`/`case R9:`/`SessionOverrideBlocked` 仍在（外科手术）。build rc=0。**test rc=1**，红名 **2**：
    `TestTicket90IrreversibleStillAsksInEveryMode` → `ticket90_test.go:335: mode=auto_approve: irreversible call raised 0 L2 cards, want 1`；
    `TestTicket90ScreenTable` → `store_test.go:299: R8 不可逆 under auto_approve: level = L0, want L2`（220 PASS + 2 FAIL = 222 ✓）。
  - **(i)-b 去掉"污染升级"红线**：分两刀，因为这条红线**有两个证人**。
    b1 只删 `if d.SessionOverrideBlocked {` 分支（`{N;N;d}`）：build rc=0，**test rc=1**，红名 **1** =
    `TestTicket90TaintFlagAndDenyAreNeverSilenced` → `store_test.go:319: mode=auto_approve: a tainted verdict came out L0 (silenced=true); PLAN.md:1640 says no authorization - including a permission mode - covers a C25 escalation`；
    真实链那条 `TestTicket90TaintEscalationNeverSilenced` **仍绿**——因为 `case R4:` 还在（第二个证人兜住），**这条绿是有原因的绿，不是假绿**。
    b2 再把 `case R4:` 一起删（两证人同灭）：PROOF `grep -n "case R4:\|d.SessionOverrideBlocked {"` rc=1，build rc=0，
    **test rc=1**，红名 **3**（+ 那条已归因的性能抖动 = 4 FAIL，218+4=222 ✓）：
    `TestTicket90TaintEscalationNeverSilenced` → `ticket90_test.go:371: mode=auto_approve: tainted call raised 0 L2 cards, want 1`、
    `TestTicket90ScreenTable`、`TestTicket90TaintFlagAndDenyAreNeverSilenced`。
    ⚠ 过程诚实记账：b1/b2 **第一次尝试编译失败**（`sed {N;d}` 少删一个 `}` → `mode.go:199:2: syntax error`），
    **编译失败不算变异**，那一轮**没有测试读数**，已还原重做并先量到 `go build ./...` rc=0 才计入上面的数。
  - **(i)-c 去掉"前端不能代表用户批准"那条红线**（机器形状=**Deny 永不被静默**：`PLAN.md:1588`/`ban #6`；
    一个被静默的判定**不是"允许"，是"没人被问过"**）：把 `if d.Level == Deny {` 分支的返回改成
    `Silenced{Level: L0, Silenced: true}`。PROOF：`grep -n "Level: Deny, Kept"` rc=1，落地行 `mode.go:161` 已打印。build rc=0。
    **test rc=1**，红名 **4**（218 PASS + 4 FAIL = 222 ✓）：票 90 的 `TestTicket90TierADenySurvivesEveryMode` →
    `ticket90_test.go:417: mode=ask_every_step: A-tier path ran the tool 1 times, want 0`（三档全红）+
    `TestTicket90TaintFlagAndDenyAreNeverSilenced` → `store_test.go:328: mode=ask_every_step: Deny came out as L0`，
    **外加两条非本票的既有卫兵一起红**：`TestBridgeRefusesTheRealShortNameOfAnAListFile`、`TestSensitiveFileIsDeniedNotEscalated`
    ⇒ A 档路径**真的被写进去了**：这条红线不是票 90 独占，链上还有旧卫兵，**双保险**。
  - **(ii) 红线改成"模式优先"**：在 `Screen()` 的**第一句**前插 `if m == ModeAutoApprove { return Silenced{Level: L0, Silenced: true} }`
    （"全自动"档跳过一切红线）。PROOF：`sed -n '157,163p'` 打印出插入的三行确为函数体首句，`grep -c MUT-D` = 1。build rc=0。
    **test rc=1**，红名 **5**（217 PASS + 5 FAIL = 222 ✓）：`TestTicket90IrreversibleStillAsksInEveryMode`（`ticket90_test.go:335` 同上一句）、
    `TestTicket90TaintEscalationNeverSilenced`（`ticket90_test.go:371`）、`TestTicket90TierADenySurvivesEveryMode`
    （`ticket90_test.go:417: mode=auto_approve: A-tier path ran the tool 1 times, want 0`）、
    `TestTicket90ScreenTable`（`store_test.go:299`）、`TestTicket90TaintFlagAndDenyAreNeverSilenced`（`store_test.go:319`）
    ⇒ **三条红线各有一个以上会红的证人，两侧变异都红**，**AC#4 成立**。
  - **还原与"没在工作树里动过手"的证据**：每轮 `cp /tmp/mode90b.orig internal/risk/mode.go` 后 `diff -q` 与快照原件**一致**，
    末轮再 `grep -c "MUT-"` = **0** 且 `go build ./...` rc=0；工作树侧 `git diff --quiet -- internal/risk/mode.go` **rc=0**，
    `git status --porcelain` 只剩票面与本票无关的票 91 文件（**不是我的，没碰**）。**全程未在仓内建 worktree/checkout**（A38④）。
  - **收尾门禁（只跑本票碰的包）**：`gofmt -l` 空、`gofumpt -l` 空（`D:/work/base/gopath/bin/gofumpt.exe`）、
    `go vet ./internal/risk/ ./internal/tools/ ./internal/perm/ ./internal/config/` rc=0、
    `go test -count=2` 同四包 **rc=0 / 0 FAIL**、`sh scripts/d22scan.sh` **rc=0**（217 生产文件，ban #1-8 clean）。
    ⚠ 如实登记不追：`go test ./cmd/wisp/` 本机加载期 `0xc0000135`（缺 sherpa dll，票 98 的账），本票**没跑它**。
  - **AC#5 台账 + 交接（不扩写别人的地界）**：`risk.Gate` 的真实调用点=`internal/tools/mode.go:98`（`readBlacklist`，
    由 `bridge.go:291-293` 在 `len(rawPaths)>0 && verdict.Level>=L1` 时调），`TestTicket90BlacklistGateIsCalledOnLivePath` 盯住。
    **但生产组合根本还没注入模式**：`grep -rln "wisp/internal/perm\"" --include=*.go .` **零命中**、`cmd/wisp/run.go` 里**没有 `Modes:` 那一行**
    ⇒ 前任说的"缩到一行"实际**连那一行都没落**。代价被 fail-closed 兜住：`Modes` 为 nil ⇒ `permissionMode()` 返回
    `risk.DefaultMode()`＝最严档（`TestTicket90UnwiredModeSourceIsTheStrictestMode` 盯），所以今天的行为是"三档全问"，
    **不会**因为没接线而变松。**交接给票 92/77**：`cmd/wisp/run.go` 里把 `perm.Store` 注进 `tools.Options.Modes`
    （以及面板只读显示 `Mode()`/`LastSwitch()`），`run.go` 是票 77 地界，本票**不代写**。
