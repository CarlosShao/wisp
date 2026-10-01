# 33 — Panel host: C27 PanelManager singleton, WebView2 window, embed.FS serving

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 12-cli-text-path-s1-gate
**Parallel slots:** 1
**Spec refs:** SPEC-08 §5.1, D29, C27, D32 panel rows, D42#11, S5

## 编排者收件（2026-09-30 23:1x，台账 `A484`）：Blocked by 07 已消 ⇒ 宿主这段提到队列头

owner 当场提了一条功能要求，逐字：「问题我早就说过了，要录入key，起码我要能看到主面板，**我要点击设置，自己配置模型这些参数**，直接给你就太不合适了」。⇒ 两件事：

1. **`Blocked by: 07-ball-state-machine-core` 这一行今天过期了**：球与托盘已进常驻进程（票 228，非实现者验收表 `228-resident-ball-v1.md` 裁成立，台账 `A477`）。本票顶部这行不改写（保留出处），但**排程按已解除走**。
2. **宿主真起来那段（AC#1..AC#4：窗口生命周期／冷启 ≤1500ms 与热显 ≤200ms 的实测／embed.FS 离线供给且无监听端口／焦点回还）提到队列头**，排在写腿 `246-r2` 交件之后——它要新增 WebView2 依赖、动构建链，所以⛔ **不与票 244（GUI subsystem 构建）同批改**，也⛔ 不与 `246-r2` 并发（同占 `cmd/wisp`）。
3. ⚠ **一处要提前说清、免得宿主起来了却交不出他要的东西**：入向白名单今天只有四枚方法（`internal/panel/bridge.go:42-45`），**零枚 config/凭据方法** ⇒ "点设置、自己录 key"这条路由属**新立的票 248**，不在本票射程内。本票的验收腿**不许**因为窗口开出来了就宣称"设置可用"。

## 编排者裁定（09-30 23:4x，台账 `A486`）：`33-h1` 的八问逐条判完，落地腿 `33-r1` 可以派

普查件＝`.scratch/wisp/probes/33/h1/census.md`（**367 行／74,999 字节**，占位符现量 **0**，70 枚编号尺，三枚 commit 写面只有这一路径；我按三把尺复认过行数与零占位符）。八问裁定：

| 问 | 裁 | 依据与边界（写死，落地腿不许自行改形） |
|---|---|---|
| **J1 宿主投 `ui-sta` 还是独立 STA 线程** | **甲：投现成的 `ui-sta`** | 独立那一支要**新增 D38(b) 名册外的一枚协程名**＝改名册属契约面，要人工批准，⛔ 本票不许顺手要。代价要**当场证**：`pkg/edge/chromium.go:96-111` 的 `Embed` 自带嵌套 `GetMessageW` 泵直到 `inited`，四类回调直接落在 `ui-sta` 上 ⇒ **落地腿的第一枚用例就证"泵期间球仍能出帧、消息仍被派发"**（普查 ⑨ 第 3 条明写它只有行号与 Win32 语义、**没跑过**）。证不出来就**停手上报**，⛔ 不许改成独立线程绕过。 |
| **J2 缺 runtime 走哪支** | **降级＋响亮，⛔ 任何路径都不许让 `log.Fatalf` 可达** | 普查在库的错误两支里翻出 **2 枚 `log.Fatalf`**，与票面 AC#5「no-crash」正面对撞。⇒ 判据＝用票面已有的 rename/mask fixture（本机装了运行时**不是**跳过这一格的理由），并具名写出我们**不走**那两枚会 fatal 的 API。 |
| **J3 `1808.9ms` 那一形算不算本票冷启口径＋P11 要不要入册** | **算另一枚口径名，不并入判据；P11 未触发，不改契约** | 三口径分开报：`cold`（沿用 S0 口径，判绿用它，P95 现量 1256.4）／`cold-embed`（embed 供给＋CSP 那两件新成本，1808.9 那一形归这里）／`recreate`。⇒ **P11（冷拉起 >2s 才重评 L2 卡回原生）两形都 <2000，未触发**；⛔ `docs/SLO.md` 与 `internal/observe/thresholds.go` 一字节不动（现量：那文件里**零枚延迟毫秒字段**，只有 `memCapPanel`），读数只进证据表。`cold-embed 1808.9 > 1500` 这一条我登记为**已知代价**（INTERIM），⛔ 不摆给 owner 做选择题。 |
| **J4 落点甲／乙** | **甲：宿主进装配根 `cmd/wisp`，`internal/panel` 保持零平台分叉** | 两形的依赖增量是**同一批 9 枚包名**，差别只在"挂哪枚包、那枚包有几份 CI 分母"，而 `internal/panel` 是仓里**唯一零平台分叉包**且在 ubuntu 有分母 ⇒ 把 Windows-only 宿主塞进去会毁掉那一格。现成旁证：`cmd/wisp/panel_assets.go:13` 逐字写着这枚宿主"expected to call the same panel.Assets API"。⚠ 票 197／238 那两条推广到本票，**普查已具名标成"这是我的推广"**，我沿用同一姿势。 |
| **J5 `GOPROXY=goproxy.cn` 取 go.sum hash 撞不撞 ban #5** | **撞，禁止** | ban #5＝从镜像站取哈希。⇒ 依赖只能走 `go mod` 正常解析、哈希由 Go 自己写 `go.sum`；⛔ 不许手工抄任何 hash 进文件，⛔ 不许从镜像目录取。**拉不到就停手上报**，不许为变绿放宽。普查另登记一条：**上游"最新 tag"两发网络尺都失败**（`WebFetch github.com/jchv/go-webview2` fetch failed），所以"最新 tag"这一格是〔未核到〕，落地时按本机 module cache 里已有的 zip/mod 走（它现量：纯 Go 零 cgo、离线可拉、要加的是 **2 枚 module／4 行 `go.sum`／9 枚包名**，`x/sys` 不抬）。 |
| **J6 票面 embed 路径 vs 活的 `all:dist`** | **按现存活接缝（`panel.Assets` 那枚 `all:dist`）做；票面那行路径不作为判据；但"清检出必须能建"升成本票一枚硬判据** | 三枚现量把这件事钉住了：票面写的 `assets/web/dist` **两种口径都不存在**；`.gitignore:24` 整条忽略且**无锚豁免**；Go 官方原话＝空目录匹配被忽略、模式不匹配则 **the build will fail**。⇒ 我**不改 `docs/specs/**` 一字**，改在本票面具名说明"规格那行的路径不作为判据"。**新增 AC#11（见下）**：在**干净检出**（临时目录另拷一份、不带未跟踪产物）里 `go build ./...` 必须通过；⚠ 我 23:4x 现量 `git ls-files frontend`＝**85 枚已跟踪文件**（`frontend/.gitignore`／`VENDORED.md` 等），**dist 那批是不是在内我没查**——这一格正是 AC#11 要证的：如果产物不入库，那"能构建"要靠构建链（票 244／CI 那侧）而不是靠谁手工生成，**不许把 `frontend/**` 提交进来糊它**（两层禁令）。 |
| **J7 "无监听端口"那半要不要进 windows CI scope** | **L1 进仓当常驻用例，L2 只留本机证据并具名登记"CI 里没有这一半"** | L1＝`go/ast` 扫宿主包 import（禁 `net`／`net/http`），在 ubuntu core scope **有真分母**（`scripts/portable-tests.sh:179`）；L2＝`GetExtendedTcpTable` 且**必须过滤 `dwState==LISTEN`**（现成那枚只数行数不读 state，直接拿来会判错，`internal/proc/treemetrics_windows.go:269`）。⛔ 词面 `netstat` 门一律不做（全仓 **0 命中**＝那种门必恒真）。`winlive` 在 CI 里**零岗位**＝这是第四枚"只在开机那台机器上成立"的同形坑，**登记而不假装覆盖**。 |
| **J8 33 与 248 谁先** | **33 先（owner 要先看见能点的窗口）；两枚可并行开发、不可合批** | 关键读数＝本票可做到**零 `internal/panel` 写面** ⇒ 与票 248（要动 `internal/panel/bridge.go`）不互斥；⛔ 但**不合并 commit**，也⛔ 不与票 244（GUI subsystem，同动构建链）同批改。串行事实：**`246-r2` 已死于 150 轮帽**（`A486`），它的残留我已收尾、`cmd/wisp` 现为干净 HEAD ⇒ `33-r1` 可以开工，前提是先复跑一次整包并逐名比红名册。 |

## 新增判据（本票自加，勾仍归非实现者）

- [x] **AC#11（09-30 23:4x 编排者追加，来路＝`33-h1` ④ 节 J6）**：**清检出能建**——把仓拷进临时目录（不带未跟踪产物）后 `GOFLAGS= go build ./...` 必须通过；若今天不通过，本票**不许**用"提交 `frontend/**` 产物"或"改 `.gitignore` 加豁免"来让它通过（那两条分别撞两层禁令与 `A207` 那味 index-aware 过滤），要**停手上报**由我改派构建链那侧。判据自带反向证：把 embed 模式改成一个不存在的路径 ⇒ 构建**必须**失败（证的是"这枚判据有牙齿"，不是"我改了就红"）。

- [ ] **AC#12（10-01 00:0x 编排者追加，来路＝我自己量的那枚 `.gitkeep`；勾仍归非实现者）**：**"能构建"不等于"有页面可发"**。现量：库里 `frontend/dist` 只跟踪**一枚 `.gitkeep`**（尺：`git ls-files frontend` ⇒ `frontend/dist/.gitkeep`），而活模式是 `frontend/embed.go:19` 的 `//go:embed all:dist`——`all:` 连点文件一起收，所以**AC#11 今天很可能是绿的，而绿的原因是目录里躺着一枚占位文件**。⇒ 本票不许把 AC#11 的绿当成"面板有内容"：判据要能区分**"embed 只匹配到占位文件"**与**"真有一包页面产物"**（问能力：embed 后的文件系统条目数／关键入口是否真存在，⛔ 不问构建退出码、不问文案）。⚠ 产物由界面侧那枚 agent 产出，本编队⛔ 不写 `frontend/**` ⇒ 这一格在"谁把 dist 填上"落定前**勾不了**，与票 248 的 AC#9 是同一件事的两面（同一枚归口）。

- [ ] **AC#13（10-01 11:2x 编排者追加，来路＝`33-v1` 验收表 §AC#3"供给那半"＋我自己现读；账 `A494`）：冷启动那段"往返探测"不许把真页面盖掉。** 现读时序（我 11:1x 亲自复认过行号）：`cmd/wisp/panel_host_windows.go:221` 调 `serveEntry()` → `:250` 把 embed 的入口字节 `SetHtml` 进控件，紧接着 `:227` 调 `firstRoundTripLocked(...)` → **`:372-375` 又一发 `SetHtml`，内容是一枚自造的探测页**（`<!doctype html>…<script>if(window.wispProbeRT)window.wispProbeRT();</script>…`）⇒ **每次冷启动最终显示的是那枚探测页，不是面板**。⚠ 这条不是"注释过期"：**面板真接进常驻之后，用户看到的仍然是枚空壳页**，而台件只会绿（它断的是"往返成立"，不是"页面是内容页"）。**完成判据**＝① 探测与供页面**次序重排**或改用资源请求过滤器（`AddWebResourceRequestedFilter` 在 `pkg/edge` 是**导出**的，同文件 `:34` 的注释自己就承认了这一点）；② 一枚**会响的**断言钉"最终文档里含 embed 入口的真内容"（问能力：解析出的条目／文档正文特征，⛔ 不问 `SetHtml` 被调用过）；③ 反控＝把两发 `SetHtml` 的次序调换，该用例**必须**红。⛔ 不许用"删掉探测"来糊——探测是 AC#4/冷启动"可用"判定的来源，删它会把另一格打空。
  ⚠ **连带更正实现腿两枚停手上报里的一枚**（判语＝`33-v1`，我追认）：「多文件过滤器接不上，因为 `PutBounds` 吃模块私有类型 `w32.Rect`」**理由不成立**——不需要构造那枚类型：`(*edge.Chromium).Resize()`（`chromium_amd64.go:12-22`）是导出的、内部自己 `GetClientRect` 再调 `PutBounds`，高层 `SetSize`／创建链路用的正是它（`webview.go:429`、`:340-344`）。⇒ **这一形属"产品形状决策"，不属"依赖边界"**，任何下一程不许再拿它当阻塞。

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

## 编排者收件（10-01 10:3x，实现腿 `33-r1` 交件；⛔ 六格判语一律等验收腿 `33-v1`，本节只记"我盘上核到了什么"）

**盘上卫生（我 10:33:59 现量，HEAD＝`94bd133e`）**：五枚提交逐枚 `git log -1` 复认——`e4ed8304`(09:52 证据件⑥⑦⑧先落)→`697b4fae`(10:10 宿主＋依赖＋测试)→`b1e3d94a`(10:26 把专用线程 harness 从产码移走，理由＝它自己撞了 `AGENTS.md` §1.2 那条 ban #1「裸 `go func(` 无 owner/recover」)→`fafe0445`(10:26 **单独一枚**按我 `A489` 的裁定反转那枚负向钉)→`f0ad3b51`(10:32 结论表)。⛔ `git status --porcelain -- cmd internal`＝**0 行**；⛔ 产码里 `MUT-` 残留＝**0 处**；票面 AC 框**它一枚没碰**（`git diff f1e7a3e2..HEAD -- .scratch/wisp/issues/` 里出现的复选框行全是**我自己**那两枚提交的内容）。

**反转后的钉我复认落地**：`internal/panel/composer_dispatch_test.go:433` 现为 `TestPanelHostIsAttachedAndNamesTheWindowHops`，旧名只剩两枚注释（`:17`／`:396`）⇒ 归档日志里那个名字不用去追改。**这枚钉从此变成"必须有宿主"的锁**——它反转是**我**的裁定（`A489` 甲形），有没有牙由 `33-v1` 用定向突变判，不由我签字。

**⛔ 一条今天最要紧的未接上事实（我在派单里就防着"import 了≠装起来了"这一形，现量确认）**：**面板今天仍然开不出来**。宿主代码进了仓（`cmd/wisp/panel_host_windows.go`），但 `grep -rn "PanelManager" cmd/wisp/main.go cmd/wisp/resident_windows.go cmd/wisp/run.go`＝**0 命中**——实现腿**没有**把常驻那条腿的 `OnPanelHotkey` 接上，也**没有**为绕开阻塞而偷加第二枚线程；它交回两枚依赖边界（`pkg/edge` 控制器 `PutBounds` 吃模块私有类型 ⇒ 多文件资源过滤器接不上；`go-webview2` 开窗是阻塞嵌套 `GetMessageW` 泵，从球的 `ui-sta` 泵再入 ⇒ panic）。⇒ 这两条真伪与边界由 `33-v1` 复认，**成立与否决定我要不要把"泵所有权"摆给 owner**（另一条形＝加第二枚常驻线程受 `observe.ResidentNames` 六枚冻结名册限制＝契约变更＝要他一句话）。**在它裁完之前，本票不许 `-done`、AC#1..AC#4 不许翻勾。**

**记在我身上的一枚派单疏漏（小、可补，具名）**：我给 `33-r1` 的派单写死了"⛔ AC 框一枚不许碰"，**但没写"必须在下面这节 Progress log 里逐枚提交留一条 `agent=… did=… next=…`"**（模板里那句"只许在进度记录里写我做了什么"这次漏了）⇒ 它五枚提交在票面上没有对应进度行，历史只剩证据件那份自述。处置＝**本节就是编排者代记**，⛔ 不冒充实现方自述、不改写它任何一行；今后派单模板补上这一句（同条已进我的记忆文件）。

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
- [ ] **AC#10（09-28 17:5x 编排者追加，来路＝`180-a1` 普查＋账 `A391`；不勾，勾要非实现者裁）**：宿主真起来那天，**`[panel]` 那五枚字段必须有生产者**——`width`／`height`／`scale` 要真去设窗口边界（`internal/config/schema.go:525-535` 是定义处，生产读取方现量＝**各 0**），且规格逐字要求 `hot` 生效（`PLAN.md:2743`＋`SPEC-03:39`）。⚠ 判据要能区分"读值并真设边界"与"只在测试里读一下"；**不许用"把 `default:"640"` 改掉"交差**（票 180 `AC#3` 同一条雷）。

---

## 编排者收件（10-01 **11:18**，收 `33-v1` 验收表 `docs/evidence/s1/33-panel-host-c27-v1.md`＝**219 行／60,796 字节**，四枚提交 `a75387cc`→`345b058b`→`e0ed6ef4`→`d6064656`；⛔ 票面 AC 框它一枚未碰，我核过＝现读 13 框（其中 AC#11 那一枚勾是**我**这轮翻的，其余 12 枚未勾））

**判语收表（勾哪几格、各归哪一发腿）**
- **AC#11 成立 ⇒ 已勾**。凭据我认：清检出 `go build ./...` 两发 exit 0，**反向证由验收腿自己跑**（只在仓外副本里把 embed 目标改名、零删除、未动 `frontend/**` 一字 ⇒ exit 1，红句逐字 `frontend\embed.go:19:12: pattern all:dist: no matching files found`）。⚠ 本格范围**只有"能建"两字**，不许任何程把它读成"有页面可发"（那是 AC#12）。
- **AC#3 不成立 ⇒ 不勾，归 `33-r3`**：「无监听端口」那一维**今天没有任何会响的检**。验收腿把产码改成真开一枚 `127.0.0.1` 监听 ⇒ 断言照绿（八发全 PASS），同发仓外 PowerShell 取到 ground truth（`wisp.test` pid 28572 确有 LISTEN `127.0.0.1:62971`）。根因定位到行，**我 11:1x 自己复认**：`cmd/wisp/panel_host_windows_test.go:115` 写 `tcpTableOwnerPIDAll = 4`，而**同一仓**的 `internal/proc/treemetrics_windows.go:57` 写 `= 5`（真值 5）；`:116` 写 `tcpStateListen = 10`，而 LISTEN 真值是 2（10 不是 LISTEN）。⇒ 两枚常量错位 ⇒ 查询取错表＋状态永不匹配 ⇒ **恒真**。附带：只读 `AF_INET`、重试耗尽静默 `return 0`。
- **AC#4 不成立 ⇒ 不勾，归 `33-r4`**（⚠ 10-01 更正：这里原先写 `33-r3`，**那枚代号 09-28 已被"入向听众"那一发用过**、它的台件目录 `.scratch/wisp/probes/33/r3/` 还在——我把新腿错派成同名，那枚新腿随即因连接报错死掉且**盘上零残留**（我核过：`git status --porcelain` 里 `cmd`／`internal` 零行、无 `33-panel-host-c27-r3.md`、票面无它的进度行）。重派代号为 `33-r4`，账 `A496`）：焦点回还那一跳**只有 `t.Logf`**（`panel_host_windows_test.go:249-251`），全仓 `prevFocus` **零枚断言**；且 `prior` 取在 `HotShow` 之后 ⇒ 11 发里"Show 之后的前台窗"**逐发就是面板自己的句柄**（等于把焦点还给刚被藏起来的那扇窗）。
- **AC#1／AC#2／AC#12 ⇒ 一律不勾**（读数真，判据形不齐或归口未落）：AC#1 有牙的是 `IsCreated` 状态断言，票面要的 **process-tree** 数成了**全机按名数**（基线 14 枚属别人会话、我们贡献约 6），且**从未比较 hide→re-show 的 HWND 身份**，dispose 那一半在产品里不可达；AC#2 的 11 枚读数（冷 max 937.093／中位 760.377，热 max 45.426／中位 37.825）未越界、P11 未触发、阈值一字未动，但票面要的 10-run P50/P95＋appendix 在仓里**不存在**（`grep P50|P95`＝0），日志里那句 "HEAD 7a41db9b" 是**硬编码字面量**；AC#12 那枚用例全文只有 `t.Logf`、零 `t.Errorf`＝**空尺**（牙在 `internal/panel/assets.go:56`＋`assets_test.go:22`，我逐名复认存在）。
- **本轮新开一格 AC#13**（写在上面「新增判据」那节）：冷启动那两发 `SetHtml` 把真页面**盖掉**了。

**交我裁的五件，我此刻的答**
1. **AC#2 按哪种形状结** ⇒ **不放宽、不改判据**：把"10-run 汇总"做成**真台件读数**，appendix 落 `docs/evidence/s1/`（⛔ 不进 `docs/SLO.md`，它在禁改清单），硬编码 HEAD 改成运行时取数。⚠ 这是修仪器，不是重定门槛。
2. **J2／AC#5 现在开还是留** ⇒ **留给下一程**，但先具名更正实现腿那句：**"不走会 fatal 的 API"不成立**——产码入口 `webview2.NewWithOptions` 自带**三枚 `log.Fatal`**（模块 `webview.go:115/120/125`），edge 层 `:173`（运行库缺失那一支）／`:295` 也在。⇒ 任何"运行库缺失"的修法**不许声称绕开了 fatal**，要真探测再创建。
3. **泵所有权三形** ⇒ **不摆 owner，先量**（理由具名：那句"从球的 ui-sta 再入会 panic"目前只是〔仅自述＋栈形状〕——无原始栈、无复现命令、用例自陈未提交、`probes/33/r1/` 里没有该发工件；而它**也没说**"投递原语已在、只差一枚导出包装"，`staThread.PostTask` 之所以只在包内可见是因为 `staThread`／`b.sta` 未导出，这直接改变代价）。⇒ 已派 `33-a2`（只读普查）。⛔ 任何腿不许为此动 `internal/observe` 那六枚名册或基线常量（那一档＝人工批准）。
4. **`staticcheck` 本机版（2025.1.1）与 CI 钉的（2026.2.1）不同、`ci.yml:181-200` 自陈那一步导入期即崩** ⇒ 这一门**我不复认**，交 CI，本表结论标〔未复认〕。
5. **AC#12 的勾** ⇒ 同意"等谁把 dist 填上"（票面 `:37`），与票 248 AC#9 同归口；但**空尺这格不等**：`33-r4` 要给它装牙（该红时红，或明确 skip 并说一句为什么），⛔ 不许留一枚"看起来在测"的用例。

**另三枚登记（只登记，本票不改文件）**
- **CI 分母错位**：`cmd/wisp/panel_host_gate_test.go:3-6` 自述"ubuntu core scope too"，实际 L1／AC#11／AC#12 三枚只在 **windows cli 腿**（`scripts/portable-tests.sh:195`＋`.github/workflows/ci.yml:475`/`388`＋`scripts/wisp-cli-tests.sh:65`），ubuntu core 那份清单（`ci.yml:341`）只有 `./internal/panel/...` ⇒ 过期指认，未改文件。⚠ 反转钉本身在 core 有分母（这条成立）。
- **`go mod tidy -diff` 在 HEAD 上 exit 1**（webview2 记成 `// indirect`；另两枚被 tidy 会抹掉的 `x/sys v0.47.0` 行起手就在）。仓内 CI 无 tidy 门 ⇒ 今天不红，但**任何腿跑 tidy 就会造出一枚不属于它的 diff**，派单里要写死。
- **起手 4 枚 `internal/panel` 红逐名判归因＝非本程**（`approval_test.go:129`／`composer_test.go:74`／`frontend_hygiene_test.go:216`／`tokens_fourway_test.go:441`；第四枚落在**冻结件**上，红因＝`git status` 里那 16 枚 ` D` 的 `design/**`）⇒ 与任务 #109"推送前逐名比红名集合"同源，我按**已知常红**读，不据此判任何格。

**排程（本轮定的序）**：`33-r4`（仪器腿，**只改 `cmd/wisp/*_test.go`**；原写 `33-r3`，代号与 09-28 那一发撞名，见上面 AC#4 那行的更正）→ `33-r2`（功能腿：接进常驻＋AC#13 那两发 `SetHtml`）→ 票 248 落地腿。⚠ `228-a2`（只读）此刻在飞、读面含 `cmd/wisp`／`internal/ball` ⇒ **要动 `resident_*.go` 的那一发必须等它交完**。

- `agent=33-a2 did=只读普查（零产码／零跑）：①库侧线程要求逐处 file:line（pkg/edge/chromium.go:95-111 阻塞嵌套泵、:130-136 无 nil 判定；库根无 chromium.go、真身 pkg/edge/）②.scratch/wisp/probes/33/ 零该发 panic 工件（唯一栈＝33-panel-host-c27-r1.md:73 散文）③甲形射程＝internal/ball 现零导出投递面（staThread/PostTask/b.sta 全未导出），会被叫红的三族逐枚点名（observe 名册钉只乙形打红；hostThreadHarness 与球侧 PostTask 取数＝行为型钉）④第四形：仓内无第二枚独立泵线程（notify_windows.go 尚未读，已挂欠账） next=编排者裁 ⑦-A（panic／阻塞／乱序唯有真跑可定）与 ⑦-E/⑦-C。文件=.scratch/wisp/probes/33/a2/census.md`
- `agent=33-a2 did=只读普查骨架先落 ⑤⑥⑦（提交 eed229e4）next=补 ①–④ 终态后交件`
- `agent=33-a2 did=终态交件：①库侧线程要求 13 处 file:line（pkg/edge/chromium.go:95-111 阻塞嵌套泵／:130-136 无 nil 判定／webview.go:108 记死 mainthread／**webview.go:445-447 Dispatch 只投线程消息 WM_APP、dispatchq 唯一读者＝webview.go:362-363 的 Run()**，本仓宿主零枚 Run/Dispatch 调用者＝R17/R18）；②.scratch/wisp/probes/33/ 全 84 枚文件＝零该发 panic 工件（散文栈只在 33-panel-host-c27-r1.md:73；且 harness 的 recover 只产 %v 无栈＝panel_host_windows_test.go:45，故"没栈"不等于"没跑"）；③甲形射程＝全目录已证 internal/ball 零导出投递面（R20）＋11 枚测试全在包内（R21），会被叫红逐枚点名：**resident_ball_228_test.go F1 TestAC228BallHostAnswersEveryGesture 是词面型在册名册钉（:50-53 十字面名册＋:319-321 枚数断言），只有往 ball.Events 加回调那一形会当场打红**；F2 行为+形状钉禁动 ball.New/defer stop；F3 能力型正向钉（搬家要留载体）；**F4 panel_host_gate_test.go:27 按文件名寻址 ⇒ 甲形不打红它但会把 L1 掏成假绿**；F6/F7 行为型待真跑；F8/F9 只有乙形打红（ResidentBaseline=6，只登记不建议）；④第四形排除 notify_windows.go（R40 零泵）与 tray（R41 复用球窗句柄），未列形状浮出＝"宿主自持 webview.Run() 那一形"且与 ui-sta 同线程互斥；⑤R1–R47 全带逐字读数；⑥14 条自攻（最重＝⑥-5 Win32 线程消息语义是通识非读数）；⑦A–H：已结清 C/D/E/F，仍判不动＝⑦-A（panic／阻塞／乱序）与新增 **⑦-H"Go→页那一跳今天没有投递者"（PLAN.md:1367 C17 逐字要求"Go→前端事件推送＋按 correlationId 路由"⇒ 属契约面未兑现）**。next=交编排者裁；⛔ 本腿未给选型 next=无`
