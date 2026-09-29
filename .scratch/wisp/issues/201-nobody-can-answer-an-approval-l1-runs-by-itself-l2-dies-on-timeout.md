# 201 — **到点该问的，今天没人问**：L1 一律自动执行、L2 一律自动拒绝；球／托盘／面板三处本来就能当"答复入口"，一个都没接

- Status: **在派队列里（等 `cmd/wisp/**` 与 `internal/agent/**` 的写面空出来）**。
- 来源：09-28 五家对标（`A403` 我方现状＋`A407`/`A408` 界面普查）。**这是差集里唯一一枚"安全＋可用性双失"的洞**，
  也是我复算过、不是转述的（判据在下面）。

## 现量（起手逐条复算，别信这里的行号）
| 事实 | 读数 | 尺 |
|---|---|---|
| 三种"人能答复"的入口 | **生产零调用者** | `grep -rn "DecideFromPanel\|DecideFromNative\|\.Veto(" --include=*.go internal/ cmd/ \| grep -v _test.go` ⇒ 只剩定义行 |
| 四条否决通道（球／Esc／面板／唤醒词） | 装配时**传空** | `cmd/wisp/run.go` 里 `NewChannels()` 的空参处（现读定位） |
| 后果 | **L1 到点即执行**（写文件类，无人能否决）；**L2 到点即拒绝**（`fs.edit`／`fs.delete`／切到全自动 全都办不成，超时退码 4） | `internal/agent/approval/gate.go` 的 `Veto`／`DecideFromNative`／`DecideFromPanel` 附近；`run.go` 里 L2 控制台只"打印卡片"那一段 |
| 载体其实**都在** | 球（`internal/ball`，20 态视觉齐）、托盘（`internal/ball/tray_windows.go`，但"打开面板"那项是空 stub）、面板审批卡（`pending` 键每轮真造） | `internal/ball/statevisual.go`；`internal/panel/composer.go` 的 `Pending` |
| ⚠ 界面白名单对不上 | 页面发 `panel.approval.request`，**Go 的入向名册里没有这一枚**（只有 4 枚），会被直接拒 | `internal/panel/bridge.go` 那 4 枚 vs `frontend/src/lib/panel.ts`（**只读，不许改**） |

## 别人怎么做（三家现抄，都带证据）
- **系统托盘里直接出批准菜单**：`Allow once` / `Always` / `Deny` / 去应用里看（openchamber `packages/electron/tray.mjs:207`，图标还会"呼吸"表示在忙）。
- **通知深链直达那条会话**（手机上收到通知就能批）：openchamber `apps/deepLinks.ts:26-33`。
- **批准卡上把"会存成哪条规则"印出来**：openchamber i18n 逐字 `chat.permissionCard.alwaysAllowPatterns: 'Always: {patterns}'`；
  minimax 的"一直允许"**还要二次确认并写明存了哪条规则**（`permission-picker.ts:118-120`）。
- ⚠ **"面板侧来源的 L2 允许"是本仓硬禁形状**（`AGENTS.md` §1.2＋`Q-49` 丙那批判据）⇒ **上面三形不许原样搬**：
  能批的只有**用户在这个界面上亲手点的那一发**，且**必须走原生确认那一条路**（`SPEC-06` 的分层），**不能做成"界面发个方法名就放行"**。

## 要建什么（四段，一段都不许砍）
1. **一条真的"答复"通路**：把 `DecideFromNative`／`DecideFromPanel`／`Gate.Veto` 接上真实入口，**先接球与托盘这两枚**（它们是本机的原生面，不碰界面那支的地盘）。
2. **到点必须问**：L1 那批发起动作**要等人**（有超时，但超时是**拒绝**而不是执行——这条是现在最危险的反差）；
   L2 那一发**要能批得动**，不再"必然超时退码 4"。
3. **答复卡要看得见内容**：卡片上有**这次/一直/拒绝**三枚语义，且选"一直"时**把会存成哪条规则印出来＋二次确认**；
   "一直"落到 D45 已设计的那张 `approval_grant` 表（现量：DAO 齐全、**生产零写手**）。
4. **等答复时球要变样子**：二十态里"等人"那一态要真被驱动（现在产品路径没有状态生产者）。

> ⚠ **09-29 12:3x 收 `201-r1` 的半成品（腿死于轮次上限，不是做完；产码我代落 commit，读数我自己跑）**
> - **盘上有**：新件 `cmd/wisp/approval_reply.go`（23,711 字节）＋`cmd/wisp/approval_reply_211_test.go`（26,269）＋`approval_reply_stdin_{windows,other}.go`；改 `cmd/wisp/main.go`（+21/−1）／`run.go`（+104/−4）／`gate.go`（+16/0）／`queue.go`（+32/0）。⚠ 文件头注释与测试文件名写着 **"Ticket 211"**——**那是我派单写错的编号**（211 是"池 8 vs 桥天花板 4"那枚契约票），同一枚文件 `:15`／`:44` 又正确引用本票 ⇒ **续腿顺手把注释与文件名改回 201**。
> - **我自己的读数**（编队已安静才跑）：`go build ./...` 净；`go test ./cmd/wisp` **ok 76.919s**／`./internal/agent/approval` **ok 0.362s**／`./internal/tools` **ok 12.861s**／`./internal/panel` **4 FAIL**＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`——**这四枚是历史在册红（`HANDOVER` §4.0v 那张表），本半成品没新增红、也没把它们修绿**。
> - **逐格现判**：**AC#1 半**——只有**控制台/stdin 那一枚入口**（`main.go:147` `reply := interactiveStdin()` ⇒ `run.go:593` `rt.attachReplyListener(...)`），球／托盘／面板 WebView **三处仍零**；它真调 `DecideFromNative`（`approval_reply.go:204`）与 `DecideFromPanel`（`:265`，`:271` 注释自陈"这条路今天被路由单独拒掉、留着当反控"）⇒ **"答"这一侧第一次有了生产调用者**，但"至少两枚入口"不满足。**AC#7 成立**（无界面时响亮兜底＝这枚监听器本身）。**AC#3 方向对**（面板侧 `allow` 被**路由**拒、不是被便利豁免）。**AC#4／AC#5／AC#6 未见实现**。
> - ⛔ **AC#2 是我自己写反的一格**：原句要"L1 挂起等人、超时落拒绝"——**与冻结语义冲突**：`internal/tools/bridge.go:384-386` 注释逐字「**The L1 window running out unopposed MEANS EXECUTE (SPEC-06 §2)**」，L1 的既定语义就是"不打断就干"，到点自动拒绝的是 **L2**。**我的裁定＝不翻极性**（翻它＝动 SPEC-06/D31，要 owner 一句话，已按 `[H#]` 记档、不催）；**AC#2 改写为**："L1 等待期间**能被否决**、**等待态在名册上可见**（票 220 那一格）；到点无答复**仍按 SPEC-06 §2 执行**"。
> - ⚠ **AC#5 我写错了对象，与本仓三枚钉互斥，就地改**：原句要 `approval_grant` 的授权"重启后仍在"——**那正是 `internal/perm/ticket90_persist_test.go:220` 判死的形状**（"session grant does NOT survive a restart"，文件头逐字"They are three test functions on purpose"），也撞 `PLAN.md:1642`「会话结束后授权**必须**失效」。⇒ **改写**："落库只落**长期**那一支，其真落点＝往 `[fs] allowed_dirs` 加一行（持久、**可撤销**），并按 `PLAN.md:1645-1646` **必须触发一次 L2 重新确认**；`approval_grant` 的 **session** 那一支**必须不**活过重启"。
> - ⇒ **续腿 `201-r2` 的射程**：AC#4／AC#6 要么落地、要么**具名报"这格要宿主面（球/托盘），今天做不到的原因是 X"**；AC#1 的入口数现跑点数报我；⛔ **不许动 L1 极性**、不许动 `ticket90_persist_test.go` 那三枚钉、命名避开 `internal/panel/l2_grant_boundary_test.go:1882-1886` 那十二个词根（**只叫 `reason` 安全**，且那枚是冻结件）。

## 验收判据（草，逐格要 `file:line` 与正控）
- [ ] **AC#1 有人能答**：三枚入口至少两枚有**生产调用者**（球或托盘算一枚，面板算一枚）。
- [ ] **AC#2 L1 不再自动执行**：造一发 L1 ⇒ **挂起等人**；正控＝超时那支必须落在"拒绝"，**不许落在"执行"**（这条要能判红）。
- [ ] **AC#3 L2 批得动**：从球/托盘点"允许" ⇒ 那一发真执行；**并且同样这一发从界面方法名直接送进来的"允许"要被判红**（`Q-49` 丙那批判据不许为便利开口子）。
- [ ] **AC#4 "一直"看得见存了什么**：选"一直"后**卡上出现规则文本**＋二次确认；正控＝去掉那行文本 ⇒ 红。
- [ ] **AC#5 授权真落库**：`approval_grant` 有生产写手，且重启后那条规则仍在（**票 181 `AC#7` 那枚教训：不许"字段/表在场、生产者从不填"**）。
- [ ] **AC#6 等待态可见**：等人答复时球进"等人"那一态；把这一态的驱动拿掉要能被判红。
- [ ] **AC#7 无界面兜底**：没有任何面可点时（纯 CLI／服务态）行为**明确且响亮**，不许静默执行（minimax/pi 两家的兜底都是"一律拒绝"，本仓按"拒绝并说明"走）。

## 禁区
`frontend/**`／`design/**` 零写面（那枚对不上的方法名**只登记、不改**，界面侧由 owner 带过去）；
C17 白名单**既有名字**不许动，**新增名字要先落一条 `A##`**（`A273`/`A389` 的口径：批准只到"新增字段"，**新增方法名不在其中**）；
`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；**面板侧来源的 L2 允许那条硬禁不许松**。

## 与其它票的关系（别在这票里顺手做）
- 界面那一侧"批准卡长什么样"＝界面那支的活；本票只交**Go 侧的答复通路、状态生产者与规则文本**。
- 唤醒词否决（KWS）与语音整条＝语音那族票；本票只保证"有面可点时必问、必能答"。
