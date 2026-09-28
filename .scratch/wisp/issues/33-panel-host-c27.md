# 33 — Panel host: C27 PanelManager singleton, WebView2 window, embed.FS serving

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 12-cli-text-path-s1-gate
**Parallel slots:** 1
**Spec refs:** SPEC-08 §5.1, D29, C27, D32 panel rows, D42#11, S5

## What to build
The `panel` module host side: singleton PanelManager owning at most ONE WebView2 window per
session (hide-don't-destroy), embed.FS resource serving via `AddWebResourceRequestedFilter`
(no localhost server), cold/hot show paths, focus return, and the no-runtime fallback
declaration consumed by 37's native card.

## Key constraints
- `jchv/go-webview2` (pure Go); one window per session; hide/show instead of destroy within a
  session; destroy at session dispose (31 hook). Cold show ≤1500ms / hot show ≤200ms — measured
  and recorded (D32).
- Resources: `//go:embed assets/web/dist` + `AddWebResourceRequestedFilter` mapping to embedded
  FS; NO localhost HTTP listener (D29; D21 anti-pattern).
- WebView2 created on the shared `ui-sta` STA thread (D38a); environment notes from spike (02)
  recorded — single-window strategy sidesteps Environment sharing entirely.
- Focus: only the panel takes focus when shown; on hide/close → focus returns to the recorded
  previous foreground window (D29).
- Multi-task partition hook: render context keyed by correlationId (surface built at 36/38).
- Runtime-missing detection: creation failure → `panel.unavailable` event + no-crash; L2 native
  fallback flag set (37 consumes); guide user to Evergreen install, never auto-download.
- CSP header/meta injection point for the embedded page (actual CSP value enforced at 35/36 —
  provide the injection hook + default strict CSP now: connect-src 'none' etc., SPEC-06 §9).

## Out of scope
- Frontend bundle contents (34); bridge protocol (35); panel pages (36–40).

## Acceptance criteria
- [ ] Show/hide/destroy lifecycle: second show within session reuses window (process-tree
      child-count stable); session dispose destroys + WebView children exit ≤2s.
- [ ] Latency: cold ≤1500ms / hot ≤200ms (10-run P50/P95 into SLO appendix).
- [ ] Embedded resource serving works offline; no listening socket exists (netstat assertion).
- [ ] Focus round-trip: editor → panel → hide → focus back in editor (automated + manual).
- [ ] Runtime-missing fixture (rename/mask loader) → fallback event, app alive, native L2 flag on.
- [ ] CSP present on served documents (response inspection) with connect-src 'none'.
- [ ] **AC#7（从票 35 快照泵 r1 验收转入，2026-09-25，来源 `35-panel-snapshot-pump-r1-accept-r1.md` F-PUMP-4）**：
      面板真的能看到快照之后，`SnapshotPump.Snapshot()` 那份包**跨两个瞬间**的形状要有一条会响的检——
      它不持泵锁地先后读 `Verdicts()` 与 `Results()`（`internal/panel/pump.go`），而今天四个读数点各自有锁、
      `SnapshotPump.mu` 只护计数器，所以 `-race` 两包全跑**是干净的**（实测 `ok 4.809s`／`ok 75.715s`、0 DATA RACE）。
      ⚠ **这一格今天不许在票 35 上开**：没有 hook 可测，硬开只会造出一枚**永远绿的假钉**（本仓那族叫"恒真判据"）。
      它要有 hook 的时刻，正是"有人在生产里读这份包的全量字节"那天——也就是本票把窗口接上的那天。
      判据：本票落地 `Marshal()` 的第一枚**生产**调用者之后，造一发"两个触发点之间状态变了"的用例，
      它必须能红（把 `Snapshot()` 改成持锁一次性快照 ⇒ 该发由红转绿）。
- [ ] **AC#8（从同一份验收转入，来源 §7 第 3 条 ＋ F-PUMP-6）**：本票接上出站通道时，**顺手量一次
      "哪些泵里的分支今天零执行"**并给它们用例或删掉。今天已知的两枚：
      ⓐ `cmd/wisp/panel_pump.go:186-192` 那条"摘要超过 440 字符就丢 ids"的分支，在 13 发变异＋工作树那发里
      **从未被走到**（落盘记录的 `pending=` 字段只出现过一枚值）；
      ⓑ `cmd/wisp/run.go:738` 那枚 `case agent.EvToolStart:`（`run.go:741` 里泵的新触发点 `changed = true`）
      **今天不可达**——`EvToolStart` 全仓**零枚发射者**（编排者 17:4x 独立复量：非测试引用恰 3 处＝定义
      `internal/agent/sink.go:25` ＋ 它的注释 `:24` ＋ 这一处消费；没有任何一处 `emit`）。
      判据：这两枚分支各要么有一发"拿掉它就红"的用例，要么被删；**不许留着让人以为它们在守什么**。

## Progress log (append-only, newest last)
- 2026-09-25 17:5x（编排者）：**加 AC#7／AC#8 两格，都是从票 35 快照泵 r1 的验收表转入的，不是本票新造的活。**
  转入口径写在这里：这两格**在票 35 上今天开不出来**——不是漏做，是**没有可测的 hook**，硬开会得到一枚
  永远绿的假钉（本仓那族缺陷有名字："恒真判据是一类新假绿"）。它们的 hook 恰好在**本票**把
  "Go → 面板"那根管子接上那天出现（AC#7 要 `Marshal()` 的第一枚生产调用者；AC#8 要真有人消费落盘摘要）。
  ⚠ 本票的 `Blocked by` 那行**没变**、`Out of scope` 那节**没动**（bridge protocol 35 仍在本票范围外）——
  这两格买的是"接上那天要顺手量"，不是把 35 的活搬过来。
  出处逐字在 `docs/evidence/s1/35-panel-snapshot-pump-r1-accept-r1.md`（F-PUMP-4／§7 第 3 条／F-PUMP-6）。
  **撤销口令：「撤 33 AC#7 AC#8」**（撤＝只删这两格＋在本节留一行撤销记录，不改写上面任何一行）。
- 2026-09-28 15:1x（只读设计核 `33-a1`，**零产码、AC 框一格未勾**）：交件全文
  `docs/evidence/s1/33-inbound-hop-design-a1.md`。本节只记四句结论，**不改写上面任何一行**：
  ① **"网页点一下 → Go 收到"这一跳被拆成 10 环（H1–H10），其中 H4–H9 六环一环都不引用 WebView2 符号**
  （现量：`grep -c webview go.mod` = 0；`ParseComposerRequest` 的签名就是 `raw string`，`bridge.go:84`）。
  ⇒ **票 186／187 不是"整体等地基"，是能先动 Go 侧那一半**（最小入向那片 = 1 枚新产码文件 + 1 枚同包测试 +
  零依赖 + 有 ubuntu CI 分母 + 零批准）。真宿主那片（H2/H3/H10 的手段）今天**在 windows scope 零用例分母**，
  要它就得先拍"改门禁形状"或"接受该腿分母＝本机人工"。
  ② **一枚规格级缺口，只报不填**：SPEC-08 §5.1 那五条宿主责任（`:143-154`）点名了生命周期／资源服务
  （`AddWebResourceRequestedFilter`，是**出向**）／无状态／Runtime 缺失／correlationId 分区，
  **没有一条写"消息接收"**；本票票面 `grep -nE "WebMessage|postMessage|inbound|receive"` = **0 命中**，
  而票 35 把 `bridge transport + dispatch`（`:7` slot A）划给自己、可 transport 要的那只手长在本票的窗口上，
  本票 `:33` Out of scope 又明写 bridge protocol 归 35。⇒ **入向接收器今天介于两票之间、两票都没逐字认领**；
  责任矩阵少一行是本跳断口的规格根因（补哪一行由编排者定，本程不动票面）。
  ③ **C17 白名单"定到哪一步"的现量比"待定稿"更难看**：规格侧 SPEC-08 `:163-174` 那 12 枚 invoke 方法名
  在**非测试 Go 代码里 0 命中**，代码侧 `bridge.go:42-45` 那 4 枚（`panel.mode.request` /
  `panel.workspace.request` / `panel.attachment.add` / `panel.message.send`）在 **`docs/specs/**` 里 0 命中**
  ⇒ **两条平行线，交集为 0**；节标题自带 `【SPEC 提案，S5 定稿走契约批准】`（`:156`）。`C24 GojaHostAPI`
  初始集一个名字都还没定（产码命中只有 `internal/config/manager.go`、`schema.go` 的配置字段）。
  ④ **P11（冷拉起 >2s 那条待定项）的现量其实已经在仓里**：`docs/evidence/s0/02-spike-report.md` §3.4
  两 run 的 cold P95 = 1041.6 / 1256.4 ms（**均未越过 2s，也均在 D32 冷 ≤1500ms 内**），唯一越过 1500ms 的是
  "进程内首次 create（竞争态）"那枚 1808.9 ms；而入向这一跳自己的传输成本有一枚可算的代理指标：
  **hot 浏览器往返 p50 = 3.5 ms**（量的就是 `w.Bind` JS→Go 那一趟，`scripts/spike/webview2-latency/main.go:118-134`）。
  ⚠ 两 run 同机、2026-09-19 未复跑，`grep -n "P11" docs/reports/pending-and-issues.md` = **0 命中**
  ⇒ 这条待定项在真相源台账里**没有对应 `A##`**，自偿要编排者先入册。
  ⚠ 本程另发现**一枚对本票不利的依赖现量**（利好写腿）：go-webview2 已被仓内那枚独立模块
  `scripts/spike/go.mod:9` 钉过（`v0.0.0-20260205173254-56598839c808`，含 indirect `go-winloader`）且
  `scripts/spike/bin/webview2-latency.exe` 真编译过 ⇒ "这枚依赖在本机能不能解析链接"**S0 已答**，
  本票落地时只剩"并进根 `go.mod` + `go.sum` + `deps.toml` 许可登记"三处动作（今天三处都 0 命中）。
  **撤销口令：「撤 33 progress 33-a1」**（撤＝只删本段＋在本节留一行撤销记录）。
- 2026-09-28 15:3x **编排者收 `33-a1`（只读设计核·零产码；派单 `2026-09-28-150x-readonly-33-a1-inbound-hop-design-core.md`）⇒ AC 框一枚没勾，但排程结论我照单采纳**（台账 `A374`）。
  **我自己复跑的**：`grep -c webview go.mod`＝**0**；`scripts/spike/go.mod:9/:18` 逐字＝`github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808`／`github.com/jchv/go-winloader // indirect`，`ls scripts/spike/bin/` 里**确有** `webview2-latency.exe`（我第一发用 `head -6` 截掉了、误以为没有＝**"空输出先怀疑仪器"又应验**）；`internal/panel/bridge.go:42-45` 逐字＝那四枚方法名；`internal/panel/composer_dispatch.go` **不存在**（落点确实是新文件）。
  ⚠ **更正它一枚读数（结论不变、尺偏宽）**：它 Q1 末句"SPEC-08 `:163-174` 那 12 枚方法名在非测试 Go 里 **0 命中**"——我逐枚 `grep -rlF` 复量**不成立**：`approval.queue` 在 `internal/statemachine/events.go:43` 逐字出现（`Event = "approval.queue-drained"`）、`approval.decide` 在 `tools/d22scan/main.go:23/361/827`（ban #6 的注释与拒绝文案）。**"两份名册没有交集"仍然成立**（没有一枚是入向方法的实现），但"0 命中"要说成"**零枚实现，两处同名噪声**"。⇒ **定式：报"零命中"必须同批给出扫了哪几个目录＋用没用 `-F`。**
  **我这一程定的三件事**：① **片 A（最小入向：`composer_dispatch.go`＋一枚具名诊断入口当生产听众，0 依赖、0 契约变更）判为可直接派**，但落点与在飞的 `181-r1` 同包 ⇒ **按住到 `181-r1` 交件**，不并行；② "要不要把 `./internal/panel/` 加进 windows CI scope"＝**不改契约、只改门禁分母、可逆**⇒ 我自决**暂不加**，真要加时单独落一枚 `A##`（因为那会把 `181-r1` 在飞的 4 枚红变成门红，责任要认）；③ "识别 git 要不要 spawn 外部二进制"＝我核过 `internal/panel/git.go` **今天零 `os/exec`**（`grep -n "exec.Command\|os/exec"` 零命中）⇒ **识别那一半不需要批准**，只有"切换"那一半要，归票 186 的 AC#2/AC#3 摆给 owner。
  ⚠ **规格级断口归我补（只追加、不改原句）**：入向接收器 **H3 在 SPEC-08 §5.1／本票／票 35 三处都没被逐字认领**（本票全文 `WebMessage|postMessage|inbound|receive` 零命中；票 35 把 `bridge transport + dispatch` 划给自己，可那只手长在本票的窗口上）⇒ **责任矩阵少一行**就是这枚断口的规格根因；下一份动本票规格面的程要把那一行补上，**别让它继续没人认**。
  ⚠ **它超顶**：工具调用 **42/35**（自己在 §11 具名登记，未美化）⇒ 我的派单坑：写了"到第 25 枚停"却没给"怎么数"的可执行尺，下一单补。它的 `internal/panel` 6 枚红（2 枚预告＋4 枚指向在飞 git 面）**它不判其对错＝正确处置**，那 4 枚要由 `181-r1` 自己交完再逐枚归因。
- 09-28 16:5x **写腿 `33-r1`（片 A·入向那一跳的六环）交件并被我核过**（末枚 `5ceffd15`，台账 `A380`；**AC 框一枚没勾**）。落了 `internal/panel/composer_dispatch.go`（199 行）＋判据件（481 行／12 枚）：入口 `Handle:120` → 唯一那道 `ParseComposerRequest:122` → 四条腿（mode `:143`／workspace `:148`／attachment `:153`／message `:158`）＋名册外走拒答（`:159`→`:167`）；真写腿**复用** `composer_handlers.go:141`、没复制。四发变异逐枚红（摘路由／默认放行／空插座放行／改 `RequestID`），`bridge.go` 白名单起手终态都＝4。
  - ⚠ **它具名更正我派单里的一发**：我写的"改 `source` ⇒ 必红"**打不红**——`bridge.go:93` 在路由前就钉死 `Source`，那条赋值走在不可能路径上；可观测的归因字段是 `RequestID`，它按可观测那枚重做。**它没为对上我的预测去拧判据。**
  - ⚠ **它答"生产调用者＝零枚"并停手**：我派单里"要一枚具名诊断入口当生产听众"和"⛔ 不碰 `cmd/wisp/**`"**本来就互斥**（矛盾在我这两句之间）⇒ **下一份入向派单必须二选一**：放开一枚 `cmd/wisp` 诊断命令，或明说"本片今天仍是零调用方、等真宿主"。**票 33 那六枚真窗口 AC 一格没结，这条不是回归、是我派单的边界写歪了。**
  - ✅ **它交付里那枚反向守卫值得抄**：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 既断"零宿主符号"又断"零监听者"，**将来谁接上真宿主，这枚判据会自己红并逼票面改口**——这是"不许拿恒真尺当新牙"的正解，我把它写进后续派单的样板。
  - ⚠ **它抓到一枚我没抓到的全仓红**（`internal/risk` 那枚"每条安全腿必须读改写账户"为 `181-r1` 的三行新消费者变红＝**票 102 的洞被新消费者重开**）⇒ 已派 `181-r2` 修，⛔ 不许动那枚判据；**我收腿必跑清单从此加 `./internal/risk/` 单包**（详见 `A380`）。
- 2026-09-28 16:3x **写码腿 `33-r1`（片 A·最小入向）交件**：交付表 `docs/evidence/s1/33-minimal-inbound-hop-r1.md`（派单 `2026-09-28-155x-impl-33-r1-minimal-inbound-hop.md`）。**AC 框一格未勾**（现量 `^- \[ \]`＝8／`^- \[x\]`＝0），原句一字未改，**只追加本节**。
  ① 新建 `internal/panel/composer_dispatch.go`（`Handle` `:120` → `ParseComposerRequest` `:122` → 四条腿 `:143/:148/:153/:158` → 默认拒绝 `:159`）＋判据件 `composer_dispatch_test.go`（12 枚）。**零新依赖、`bridge.go` 零字节改动、`cmd/wisp/**` 未动**。`bridge.go` 白名单枚数两向现量：常量尺 `= "panel.` **4→4**；⚠ 派单 §1 那把字面尺 `grep -c 'panel\.'` 起手就是 **5**（`bridge.go:35` 注释含 `"panel.*"`），两向差 0＝我没动过它，那把尺的"＝4"写法按"尺与描述不符"上报。
  ② 三发判据**先红后绿**（载具 `.scratch/wisp/probes/33/r1/mutate.sh`，复原后 hash `32e39e0a…` 逐字相同、重跑 `ok`）：摘路由 ⇒ `TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce` 红（**数被调次数**）；默认分支改放行 ⇒ `TestRosterMismatchBackstopRefusesInsteadOfAccepting` 红；nil 插座改放行 ⇒ `TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped` 红（4 枚）；改写 `RequestID` ⇒ `TestAttributionFieldsReachTheHandlerUnchanged` 红。⚠ **一枚载具级更正**：`33-a1` §7.4 第 3 发"改 `source` ⇒ 拒答"在绿色态已结，但**"入口改写 source"这一发变异造不出红**——`bridge.go:93` 在派发前就把 `Source` 钉死，路由内那句赋值落在不可能路径上；可测的归因字段是 `RequestID`（理由逐字在 `mutate.sh` 里）。
  ③ **AC#D 的诚实答案：生产调用者＝0 枚，本程没有那枚具名诊断入口。** 原因不是"要真 WebView2"，是**我的写面里没有 `cmd/wisp/**`**（`33-a1` §7.3 那形状要先新建一枚 `cmd/wisp/*.go`，同族先例 `panel_assets.go` ⇒ 用户敲的那条命令要等写面追加）。⇒ **停手上报，已交绊线**：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 同时钉住"本包零宿主符号"与"零生产听众"，谁接上真宿主它会主动报红并要求在票面具名——**这就是防"用假宿主充当已接线"那一枚**。`perm.Store.Set` 的生产调用者**今天仍是 0**，`HandleModeRequest` 非测试调用者也仍是 0。
  ④ 门禁终态：`internal/panel` **FAIL 恰 3 枚＝在册那 3 枚**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`），本程新文件带来红＝**0**；`internal/config` `ok 0.955s`、`internal/tools` `ok 15.428s`、`gofumpt` 名下空、`d22scan` rc=0 且 `ban #8 internal/` **436→438（+2＝本程两枚，逐名）**、`gate-clauses` 红腿名册仍**只 `G6neg`**（比名册不比退码，尺未改）。⚠ **一枚不在我名册里的红，具名上报不修**：`internal/risk`（**单包跑**，非 `A359` 那族并发假红）`--- FAIL: TestC26RewriteAccountIsConsumedAtEverySecurityLeg`——本程未动 `internal/risk/**` 一字（终态写面闸门只两枚 `internal/panel` 新文件），归因请编排者落 `A##`。
  ⑤ AC#E 词面尺：`composer_test.go:268` 那把尺**一字未动**；本程新产码不含任何"切换"词面，工作区那扇门只声明插座、处理器留给本票下游。自证＝`TestTheInboundHopAddsNoSwitchingCapability`（绿）。
  **撤销口令：「撤 33 progress 33-r1」**（撤＝只删本段＋在本节留一行撤销记录，不改写上面任何一行）。

