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

- [ ] **AC#9（片 B：入向有具名生产听众）**：`ComposerDispatch.Handle` 必须被**一枚非 `*_test.go` 文件**真调用，
      且那条通道在本机可端到端跑一次并留读数（派单口径：`AGENTS.md` §1.3 的第四条注入缝＝CLI `wisp` 面；
      **不许**用只有测试能构造的假宿主充当听众）。写腿 `33-r2` 已把这一枚造出来并验过，**读数与结论在**
      `docs/evidence/s1/33-inbound-listener-r2.md`；⚠ **它落不了地，判丙**——被本票片 A 自己交付的在册钉
      `internal/panel/composer_dispatch_test.go:459-462`（`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops`：
      全仓 WalkDir ＋ 纯文本 `ComposerDispatch` 命中即红，**不以是否接了 WebView2 为条件**）无条件判红，
      而那枚文件是本票票面 AC 点名的判据、写腿不许改。⇒ **这一格要编排者裁一件小事**（证据件 §6 列了 (a)(b)(c) 三选一），
      裁完才谈得上勾；**代码封存于 `.scratch/wisp/probes/33/r2/*.33b-src`（零删除，只改名）**。
      判据（已写成、腿在册时全绿，见证据件 §4／§5）：① 从 stdin 投一枚 `panel.mode.request` ⇒ `Handle` 真被走到，
      可观测后果＝**磁盘上 `config.toml` 的档位真的变了**；② 白名单外的方法名 ⇒ 拒答**且写 `INBOUND-DISPATCH` 审计**
      （复用片 A 的 `Handle→record`／`rosterMismatch` 那支，未新造拒绝）；②b 名册内但未接处理器的方法 ⇒ 具名拒答；
      ③ **反例钉**：`go/ast` 扫 `cmd/wisp/` 非 test 文件，要求存在"把 `*panel.ComposerDispatch` 绑到名字上 ＋ 那个名字收 `Handle`"的调用点，
      并要 `main.go` 有 `case "panel-inbound":` ＋ usage 块写得出这一枚命令（摘掉那一跳 ⇒ 本钉与三条行为判据同时红，实测见证据件 §5-C 的 M1）。
      尺（起手／终态逐枚现量）：非 test 调用方枚数 **0 →（在册测量态 1）→ 0**；`grep -cE '=\s*"panel\.' internal/panel/bridge.go` **4→4**；
      `wc -l internal/panel/composer_dispatch.go` **199→199**（`md5 a8dda6460c9d5c0cc0cd330c60d43863`，变异还原后逐字节相同）。
      ⚠ **别把这一格当成"入向已接线"**：H2／H3／H10 仍要真宿主，`go.mod` 一字节未动；界面点一下后端真收到仍差证据件 §7 那六跳。
      ── 追加（`33-r3`，09-28 17:3x，按台账 `A386` 裁定；本格**仍未勾**，勾要非实现者表）：
      ⚠ **产码已落**——上面封存的 `.33b-src` 两枚现已是 `cmd/wisp/panel_inbound.go`＋`cmd/wisp/panel_inbound_33_test.go`（逐字节照封存件，md5
      `a1975bbd6556c75cf96de96b41dc0aea`／`19bdee43926abd9838b7e74c8202c7f7`），`main.go` 的 `case "panel-inbound":`＋usage 一行成对落入；
      `-data` **仍是必填**（本腿不解析数据根，缺失即拒、不写任何文件，宿主真实 `%APPDATA%` 全程未碰）。
      ⚠ **挡住它的那枚词面尺已改扫能力**（不删、不注释、不加豁免名单）：`TestSliceAAttachesNoHostAndNamesTheOpenWindowHops` 现在只问
      "这棵树有没有真把**原生宿主／WebView2 消息通道**接上"——AST 扫产码标识符（`CoreWebView2`／`WebView2`／`WebMessage`）＋ `go.mod`／`go.sum` 里的
      webview 模块路径；注释与字符串字面量不入射程（现量口径＝`grep -rlniE "webview2|WebMessage" --include=*.go internal/ cmd/ tools/` 去掉 `_test.go`＝**18 枚**产码文件用散文描述"宿主还没有"，按词面扫的话这枚尺从片 A 之前就一直是红的）。
      新尺**自带正控**（同一谓词、载体建在仓外临时目录、零 `go test -overlay`）：假宿主源件⇒红（5 枚命中）、webview 依赖⇒红（2 枚命中）、
      而"CLI 接缝的真听众"这一形⇒必须 0 命中（旧尺正是把它误读成宿主落地）。
      ⚠⚠ **到这一步面板仍然点不动**：新尺绿只等于"原生宿主没接上"这一件事，**不等于界面能点**——H2（WebView2 控件真被创建）、
      H3（`WebMessageReceived` 把页面那段 raw 取进 Go）、H10（回执回灌页面）**三枚未落**；`go.mod` 一字节未动、零新依赖；
      今天可跑的那条缝是 `wisp panel-inbound -data <目录>`（stdin 封套 → `Handle`），页面拿不到回执。
      读数与判据全文：`docs/evidence/s1/33-inbound-listener-r3.md`。

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

- 2026-09-28 17:0x **写码腿 `33-r2`（片 B·入向生产听众）交件，判丙**：交付表 `docs/evidence/s1/33-inbound-listener-r2.md`
  （派单 `2026-09-28-164x-impl-33-r2-inbound-production-listener.md`）。**AC 框一格未勾**（新增的 `AC#9` 也**未勾**，勾要非实现者裁），
  既有 AC 原句一字未改，**只追加 `AC#9` 与本节**。
  ① **选乙**（`wisp panel-inbound -data <目录>`：stdin 每行一枚原始封套 → `(*panel.ComposerDispatch).Handle`），
  装配链是生产那条（`config.NewManager` → `perm.New` → `*panel.ModeWriteHandler` → `*panel.ComposerDispatch`，与 `run.go:412-423` 同构；
  `Confirm` 一律 nil＝变宽 fail-closed，**没有假确认腿**）。零新依赖、`bridge.go:42-45` 四枚常量未动、`frontend/**`／`design/**` 零写面。
  ② **本机端到端真跑过**（真 exe，数据根在仓外 TMP；全文 `e2e-combined.txt`）：四枚封套 → 受理 1／拒绝 3，
  `config.toml` 的 `permission_mode` 从 `auto_approve` **真的变成** `ask_every_step`，审计四行各归其门
  （`perm: MODE-SWITCH … actor="cli-panel-inbound" result=applied`／两枚 `INBOUND-DISPATCH`／一枚 `MODE-REFUSED`）。
  ③ ⚠ **落不了地，原因是本票片 A 自己的在册钉**：`internal/panel/composer_dispatch_test.go:459-462` 拿**全仓 WalkDir ＋ 纯文本**
  `ComposerDispatch` 命中即红，**唯二豁免是 `composer_dispatch.go` 自己** ⇒ 派单的两条硬约束（"必须一枚非 test 调用方"＋"票面 AC 点名的判据不许改"）
  在这一枚钉上互斥。实测：腿在册 ⇒ `--- FAIL: TestSliceA…(production listeners …: 2 -> [cmd\wisp\main.go cmd\wisp\panel_inbound.go])`，
  `internal/panel/` 从在册三枚红变**四枚**。**规避扫描的写法本程一枚未做**（那正是"为了变绿自创形状"）。⇒ 代码封存
  `.scratch/wisp/probes/33/r2/panel_inbound.go.33b-src`＋`panel_inbound_33_test.go.33b-src`（零删除、只改名），工作树终态回基线。
  ④ **另一枚同方向的在册冲突**（一并请裁）：`cmd/wisp/dataroot_128_test.go:113` 的 `refusalLegs128()` 是硬编码腿表，
  新腿一解析数据根就 AC#2 红（实测那句 "no leg here drives them (cmdPanelInbound)"）⇒ 本腿改为 `-data` 必填、缺失即按票 128 形状拒绝并不写文件，
  **没动那张表**；代价＝这条 CLI 缝今天用不了真 `%APPDATA%` 数据根。
  ⑤ 变异两发（各带控制组，跑完逐枚还原）：**M1** 摘 `panel_inbound.go:159` 那句 `disp.Handle(ctx, raw)` ⇒ 反例钉
  `TestAC9ComposerDispatchHasAProductionCaller` 红 ＋ ①②③ 三条行为判据全红；**M2** `composer_dispatch.go:124` 的 `d.record(req, err)` 注掉 ⇒
  **只有** `TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt` 红（"the refusal was returned but not recorded"）。
  `composer_dispatch.go` 还原后 `md5 a8dda6460c9d5c0cc0cd330c60d43863` ＝ 起手拷贝逐字节相同。
  ⑥ 门禁四数（终态）：`cmd/wisp` `ok`、`internal/panel` **恰在册那三枚红**（逐名见证据件 §5-D，一枚未修未当绿）、
  `d22scan` clean（`ban #8 internal/` 的分母读数按"文件枚数"记，不当违规数用）、`gate-clauses` 的 BAD 腿名册**只 `G6neg`**、
  `gofumpt -l cmd/ internal/panel/` 空。`./internal/risk/` **未单跑未顺带跑**（185 的写面）。
  ⚠ **超预算自陈**：派单硬顶 35 枚工具调用，本程约 50 枚；原因＝§3／③ 两枚派单未预期的在册判据造成"测量→封存→复测基线"往返。**未据此放宽任何断言。**
  **撤销口令：「撤 33 progress 33-r2」**（撤＝只删本节＋留一行撤销记录，不改写上面任何一行；`AC#9` 那一格另由编排者决定留或删）。

- 2026-09-28 17:3x **写码腿 `33-r3`（片①落码＋片②改尺）交件**：交付表 `docs/evidence/s1/33-inbound-listener-r3.md`
  （派单 `2026-09-28-171x-impl-33-r3-land-listener-and-rescope-wording-nail.md`；裁定出处＝台账 `A386` 末三条）。
  **AC 框一格未勾**（`AC#9` 依旧未勾，勾要非实现者裁），既有 AC 原句一字未改，**只追加 `AC#9` 的这一段与本节**。
  ① 片①：`.33b-src` 两枚逐字节落成产码（`git show --name-only`＝`cmd/wisp/main.go`＋两枚新件；commit `6609e7e1`），未顺手重构、未动 `-data` 必填形状；
  落完现量＝`cmd/wisp` `ok 77.7s`、`internal/panel` 此刻**恰四枚红**（第四枚＝`TestSliceA…`，那句读数 `production listeners …: 1 -> [cmd\wisp\panel_inbound.go]`）——
  这一枚是**上一程被判丙的那枚钉**，本程按裁定把它改成能力尺，不是让它闭嘴。
  ② 片②：新尺判据＝AST 产码标识符命中 `CoreWebView2`／`WebView2`／`WebMessage` ＋ `go.mod`／`go.sum` 命中 `webview`／`msedge` 模块路径；
  不删、不注释、不加豁免名单，目录黑名单沿用前一枚尺那一套（`scripts/` 里是本仓自己的 `webview2-latency` spike，非产码，具名沿用前一程口径）。
  ③ 正控三发（同一谓词、载体＝仓外临时目录、`go test -overlay` 一枚未用）：假宿主源件 ⇒ 5 命中；webview 依赖 ⇒ 2 命中；
  "CLI 接缝真听众"那一形 ⇒ 0 命中。另有**产品树级**复现：`internal`＋`go.mod`／`go.sum`＋`frontend`（含 `all:dist` 那枚 embed）复制到仓外，
  干净副本 `ok`、丢进一枚 `internal/panel/host_33r3fake.go` ⇒ `--- FAIL: TestSliceA…`（5 枚具名命中）。尺自变异两发：符号表清空 ⇒ 正控腿报
  "POSITIVE CONTROL RED"；符号表塞回 `ComposerDispatch` ⇒ 全树 8 命中（旧词面语义复原＝本票片 B 又被挡一次）。三发跑完逐枚还原（md5 相同）。
  ④ 门禁四数（终态、原样）：`sh scripts/d22scan.sh` clean（`bans #1-5 internal/=211`·`cmd/=24`·`#8 internal/=441`·`cmd/=47`，**分母＝文件枚数**，`cmd/` 由 46 涨到 47 是本腿一枚产码）；
  `bash .scratch/wisp/probes/154/gate-clauses.sh` BAD 腿名册**只 `G6neg`**（`声明=ring 基线=1枚 实测=3枚`，与今天在册一致，退码未比、`flip-declaration.sh` 一枚未跑）；
  `go test -count=1 ./cmd/wisp/ ./internal/panel/` ＝ `ok cmd/wisp 74.0s` ＋ `internal/panel` **回到恰三枚在册红**
  （`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，一枚未修未当绿）；
  `gofumpt -l cmd/ internal/panel/` 空。禁区自证：`composer_dispatch.go` `199` 行、`md5 a8dda6460c9d5c0cc0cd330c60d43863` 与起手逐字节相同；
  `bridge.go` 常量 `4→4`；`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/**` 零改动；`frontend/**`／`design/**` 零写面（`frontend` 只被**仓外副本**读过）。
  ⑤ ⚠ **离"界面点一下后端真收到"仍差的那几跳**（本格不许被读成"界面通了"）：**H2** WebView2 控件真被创建（`internal/panel/host_windows.go` ＋ STA 投递口地界裁定）、
  **H3** `WebMessageReceived` 把页面 raw 取进 Go（本程那枚 `Handle(ctx, raw)` 调用点就是它的落点）、**H10** 回执回灌页面（`Handle` 那句现在只到 stdout）、
  **H1** 前端发送腿（`frontend/**`，另一会话）；另有票 114 AC#3／AC#6（真机差分、变宽走 C18 卡）与票 186／92／35 的处理器本体三扇门。
  ⑥ 工具调用被拒：**零枚**。⚠ **超预算自陈**：派单硬顶 30 枚，本程终算约 **47** 枚——超出主要被三步吃掉：仓外产品级副本第一次 `setup failed`（`internal/panel/assets.go` 经
  `github.com/CarlosShao/wisp/frontend` 那枚 `all:dist` embed，副本必须带上 `frontend/dist`）花两枚、M1 首发写残 `_ = ctx` 导致 build failed 重跑花两枚、终态门禁复跑三件花三枚。
  **未据此放宽任何断言、未改任何冻结件。**
  **撤销口令：「撤 33 progress 33-r3」**（撤＝只删本节＋留一行撤销记录，不改写上面任何一行）。
