# 201 — **到点该问的，今天没人问**：L1 一律自动执行、L2 一律自动拒绝；球／托盘／面板三处本来就能当"答复入口"，一个都没接

- Status: **部分结案（09-29 12:0x：七格里已勾 2＝AC#4／AC#7；未闭 5＝AC#1／AC#2／AC#3／AC#5／AC#6，前置具名见下面那节）**。⚠ **本票剩余的 AC#1／AC#2／AC#3a／AC#6b 全部归口同一枚前置——"宿主原生入口（球／托盘）在 `wisp run` 里没接"**，Go 侧的答复逻辑今天第一次全齐；⛔ **不翻 `-done`**（`-done` 是防重领唯一键，须 0 枚未勾才结）。早先那行"在派队列里"的口径已被上面这行取代，原文不再单独保留。
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

> ⚠ **09-29 12:0x 收 `201-r2`（五枚提交，逐枚 `git log -1` 核过存在）＋非实现者裁决 `201-v1`（`docs/evidence/s1/201-reply-listener-r2-verdict.md`，45,385 字节，提交 `21d9a310`）⇒ 逐格翻勾。先更正我自己：**上面 36／37 两行的"三枚钉"实为四枚**（多出 `cmd/wisp/run_mode101_test.go:389`，见台账 `A426`）**、37 行引的 `PLAN.md:1645-1646` 应为 `:1644-1645`**（`A423`）——**原句不抹，就在这里具名作废。**
> - **盘上多出来的**：`internal/agent/approval/replies.go`（宿主面：`Attach`／`Record`／`Look`／`Forget`／`Pending`／`AwaitingHuman`／`WaitingState`／`Allow`／`Reject`／`PanelReject`／`PanelAllow`，逐枚现读到）＋`replies_201_test.go`、`internal/config/allowdirs.go`（`[fs] allowed_dirs` 写路径：含 `..` 拒收、写失败回滚、写完采纳自己的 mtime）、`cmd/wisp/approval_always.go`（`always` 答复＝先印规则文本→再签一张真 L2 复确认卡→落 `allowed_dirs`→明写"下一次启动生效"）＋两枚 `cmd/wisp/*_201_test.go`。
> - **我自己跑的整包（不是腿的转述）**：`PATH` 带 sherpa、`-count=1`，起手 11:50:49／终态 11:52:27，全量存 `.scratch/wisp/probes/201/r2/gates-full.txt`＝**23 枚包 ok／3 枚包 FAIL 共 6 例**：`internal/ball` 1 例＋`internal/panel` 4 例（逐字原因都指向 `design/assets/tokens.css` 被另一队在工作树里删掉＋面板契约字段差＝**别人地界，不许顺手修**）＋**`internal/risk` 1 例 `TestResolvePerCallBudget`**（争用期 1.222ms／预算 1ms）⇒ **机器空下来 `-count=3` 复量全绿＝CPU 争用型假红、不是回归**，`thresholds.go` 一字节没动。⚠ 我一度把这里读成"两枚包 5 例"、并把 `config/perm/memory` 的耗时抄成了上一发被 `| tail -60` 截断的数——**两处都由裁决腿指出、我已复跑更正，记台账 `A428`**。
> - **翻勾结果**：✅ **AC#4**（正控真实存在：`approval_always_201_test.go:208-209` 断"要存的规则："与 `allowed_dirs += "` 逐字上屏、`:211-215` 断第二张 L2 卡真出现、`:218` 断它自己也印规则——我读到原文）；✅ **AC#7**（无面可点时兜底响亮＝这枚监听器本身）。⛔ **不翻**：**AC#1**（入口枚数仍＝1）／**AC#2**（改写形也不成立——**`replyVeto` 生产零赋值点**：`run.go:146` 只有声明、`:601` 只有读，全仓无赋值 ⇒ **"L1 等待期间能被否决"今天不可达**）／**AC#3**、**AC#6**（见下面那两格里我的理由）／**AC#5**（`approval_grant` 写手仍为 0 ⇒ "session 那一支不跨重启"被"根本没写"满足＝**恒真判据**；**且"可撤销"这一半今天零门：`SetAllowedDirs` 全仓零调用者、零用例** ⇒ 整格前置转**票 224**，撤销那半转**票 219 新增一格**）。
> - ⚠ **我与裁决腿不一致的两格（口径落档 `A428`）**：它裁 AC#3／AC#6"可翻（带注）"，我**按字面不翻**——两格文字都点名"**从球/托盘点**"／"**球进那一态**"，而**同一份表**在 AC#1 明说入口枚数＝1，⇒ 同表内自相矛盾；且 AC#6 那枚 `bookWaitingState` 函数体今天**只写审计日志、不驱动球**（我读了 `approval_always.go:171-186`）。⇒ **规矩：裁决腿的判语不自动等于翻勾；字面判据优先，分歧必须写死在票面。**
> - **要人批的两格拆法（按 `[H#]` 记档、不催 owner）**：AC#3 → **AC#3a 原生入口真执行**／**AC#3b 界面来源必被判红**（3b 今天已成立）；AC#6 → **AC#6a 等待态可计算可审计**（已成立）／**AC#6b 球真进那一态**（未成立）。**建议把本票剩余的 AC#1/2/3a/6b 全部归口到"宿主入口"那一枚前置**——它与票 219 的"三枚按钮"是同一条线，**这条线缺的是原生面，不是 Go 侧的答复逻辑**（Go 侧答复逻辑今天第一次全齐）。

## 验收判据（草，逐格要 `file:line` 与正控）
- [ ] **AC#1 有人能答**：三枚入口至少两枚有**生产调用者**（球或托盘算一枚，面板算一枚）。
- [ ] **AC#2 L1 不再自动执行**：造一发 L1 ⇒ **挂起等人**；正控＝超时那支必须落在"拒绝"，**不许落在"执行"**（这条要能判红）。
- [ ] **AC#3 L2 批得动**：从球/托盘点"允许" ⇒ 那一发真执行；**并且同样这一发从界面方法名直接送进来的"允许"要被判红**（`Q-49` 丙那批判据不许为便利开口子）。
  - ⚠ **本格今天不翻，且我与裁决腿判语不一致（`A428` 记档）**：非实现者裁"**可翻**"，但这一格**字面**要求"**从球/托盘点**允许 ⇒ 那一发真执行"，而**同一份裁决**在 AC#1 那格明说"入口枚数仍＝1"（球只被 `cmd/balldebug` 引、面板无 WebView2 宿主）⇒ **两格同时成立会自相矛盾**。⇒ 处置＝**按字面不翻**，把"界面方法名送进来的允许被判红"这半支单独记为**已成立**（现读：`Replies.PanelAllow`→`ErrPanelAllow`、缺 `grant`→`ErrBadGrant`，见裁决表），另**建议把本格拆成 AC#3a（原生入口真执行）／AC#3b（界面来源必被判红）**——**拆格属改判据＝要人批**，我按 `[H#]` 记档、不自己拆。
- [x] **AC#4 "一直"看得见存了什么**：选"一直"后**卡上出现规则文本**＋二次确认；正控＝去掉那行文本 ⇒ 红。
- [ ] **AC#5 授权真落库**：`approval_grant` 有生产写手，且重启后那条规则仍在（**票 181 `AC#7` 那枚教训：不许"字段/表在场、生产者从不填"**）。
- [ ] **AC#6 等待态可见**：等人答复时球进"等人"那一态；把这一态的驱动拿掉要能被判红。
  - ⚠ **本格今天不翻（我与裁决腿的"可翻带注"不一致，理由现读）**：Go 侧**确实第一次**产出了 D43 的等待态且有生产调用者（`cmd/wisp/run.go:1110`／`:1156` → `cmd/wisp/approval_always.go:171 bookWaitingState`，我读到了这两处调用），**但那枚函数体今天只做一件事——写审计日志**（`rt.auditf("approval: WAITING-STATE state=%q ...")`，逐行读于 `approval_always.go:171-186`）⇒ **没有驱动球的那一态**，而球本体根本不在 `wisp run` 里（`internal/ball` 只被 `cmd/balldebug` 引＝AC#1 那半同一根因）。⇒ 本格与 AC#1 是**同一枚前置**，前置不满足就不翻；**"态已可计算／已可审计"这一半单独记为已成立**，写在下面那节，免得下一位重造。
- [x] **AC#7 无界面兜底**：没有任何面可点时（纯 CLI／服务态）行为**明确且响亮**，不许静默执行（minimax/pi 两家的兜底都是"一律拒绝"，本仓按"拒绝并说明"走）。

## 禁区
`frontend/**`／`design/**` 零写面（那枚对不上的方法名**只登记、不改**，界面侧由 owner 带过去）；
C17 白名单**既有名字**不许动，**新增名字要先落一条 `A##`**（`A273`/`A389` 的口径：批准只到"新增字段"，**新增方法名不在其中**）；
`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；**面板侧来源的 L2 允许那条硬禁不许松**。

## 与其它票的关系（别在这票里顺手做）
- 界面那一侧"批准卡长什么样"＝界面那支的活；本票只交**Go 侧的答复通路、状态生产者与规则文本**。
- 唤醒词否决（KWS）与语音整条＝语音那族票；本票只保证"有面可点时必问、必能答"。

> ⚠ **09-29 14:3x 收只读普查腿 `201-c2`（`.scratch/wisp/probes/201/c2/census.md`，编排者代落盘＋逐条自己跑尺）⇒ 本票剩余四格的性质变了，三处就地更正（原句不抹）**
> - **①"缺宿主原生入口"这句不完整**：普查现读——**`wisp run` 跑完一发任务就退**（`cmd/wisp/main.go:85 os.Exit(cmdRun(...))` → `cmd/wisp/run.go:220 return rt.execute(task)` → `:855 bg.Wait()` → `:914 return`，编排者复跑行号复认），而**常驻那条腿 `cmd/wisp/resident_windows.go:24/:83` 没有 gate、没有 bridge、没有任务**。⇒ **两条腿各缺另一半**，本票 AC#1／AC#3a／AC#6b **整条前置转新立的票 228**（含"球挂 run 还是挂 resident"那枚未定义即停点）。
> - **②票面 §13／§55 那句"托盘'打开面板'那项是空 stub"口径错**：`internal/ball` 里**没有 stub**，它已经把选择发成回调了（`ball_windows.go:667/672` `b.fire(Events.OnTrayPanel)`）；那句 `fmt.Println("tray: open panel (stub, ticket 33)")` 在**消费者** `cmd/balldebug/main.go:199`。⇒ **对写腿的直接含义变了**：球那一面不缺"发"，缺"听"——落点不要在 ball 包里找空函数。
> - **③AC#2 那半格今天缺的是一枚赋值，不是功能**：`runSpec.replyVeto` 生产赋值点＝**0**（`run.go:146` 声明、`:601` 读，我复跑复认），而**现成半成品就在场**——`internal/agent/approval/approval.go:163 DefaultChannels()` 逐字 `return NewChannels(ChannelBall, ChannelEsc)`，**生产零调用者**；装配处 **`run.go:446` 传的是空参 `approval.NewChannels()`** ⇒ **今天是被"显式关掉"，不是"没建"**。三件都在（`approval.go:147`／`:163`、`gate.go:139 Channels()`、`approval_reply.go:415-417 SetLoaded`），**只差那一枚赋值** ⇒ **这一格可以从票 228 里 separable 地早做，但它撞 `run.go` ⇒ 仍须排在票 223 之后。**
> - **还有一条会影响票 228 排期的旧账**：`rt.close()`（`run.go:676-683`）**从不 join** 它 spawn 的 `replyHandle`（`:265`），注释 `:672-675` 自陈协程随进程死 ⇒ **D38(c)（`PLAN.md:2845`）"任务完成必须等所有派生 goroutine 退出"今天已经违约**（⚠ 这是**读到码**、不是读到尺；"有没有仪器在查这条"普查腿没查、我也没查）。**挂 GUI 之前先补这一格，否则新账叠旧账。**
> - **答复侧今天的真实状态（免得下一位以为还要重造）**：`Gate.DecideFromNative`（`gate.go:626`）／`DecideFromPanel`（`:638`）／`Veto`（`:387`）**生产调用者 5 枚、全在 `replies.go:325/372/377/401/427`**；`Replies` 那一层由 `approval_reply.go:210/249/251/278/310` 调；待答卡数据原生面**直接读得到**（`Replies.AwaitingHuman()` `replies.go:249`、`Pending()` `:224`、`ReplyCard` 带 `CorrelationID/Level/Tool/Grant/Paths`）。⛔ **但"L1 那一发"没有枚举口**（`Gate.windows` 全仓只有 `:321/:324/:330/:393` 四处命中，我复跑）⇒ 要枚举它就撞**票 220 AC#2 甲形**，先裁甲乙。
> - ⛔ **票面 §10 那把尺的读数（"三枚入口生产零调用者"）已过期**，上面 ③ 与 r1/r2 两节已各自解释过原因；**本票 AC#1 的"入口枚数＝1（控制台 stdin）"仍然成立**——球／托盘／面板 WebView **三处仍零**。
