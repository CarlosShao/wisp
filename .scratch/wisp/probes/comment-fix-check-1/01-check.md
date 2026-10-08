# 01-check — `comment-fix-prep-1/01-ready-to-apply.md` 独立复核（只读腿 comment-fix-check-1）

- 复核对象：`.scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md`（385 行、13 处〔只欠文案〕）。
- 现量：`2026-10-08 18:41+08` 起，HEAD `0a09a4f2`（`dev`）；起手锚见 `00-anchor.md`（commit `a1f0d2ff`）。
- **零 Go 命令**（181-v3 正种突变）；全部读数＝`git show/grep/ls-tree/cat-file`＋`sed/grep/awk/wc/diff`。
- 尺的射程逐条写明；范围尺一律现检 rc（下文 rc 均现跑）。
- 基线对账：prep 自报基点 `9e8477ed`；`git diff --name-only 9e8477ed 0a09a4f2` 非 `.scratch` 只动了
  `docs/reports/HANDOVER.md`／`docs/reports/pending-and-issues.md` ⇒ 13 处产码锚无中间漂移面。
- 工作树：`git status --porcelain -- <13 个目标件>` 输出空 ⇒ prep 在飞登记的 `panel_host_windows.go` 脏已不在，
  锚按 HEAD 取即按当时工作树取。
- 写面声明：本腿只新建本目录 `*.md`；未改任何产码/票面/台账/他人件；未 push。
- **射程判断，非内容引用**：范围尺命中过三枚冻结件之一 `internal/perm/ticket90_persist_test.go`（`ConfirmLocked` 命中
  :329/:342/:386）——只作射程命中登记，未读、未引其内容。

## 一、13 处四问表

| 处 | 锚在不在（逐行/逐字） | 原文逐字 | 替换文本的事实断言（逐句现验） | 撞禁词 | 判 |
|---|---|---|---|---|---|
| P01 | 在：`internal/panel/composer_dispatch.go:47-52`，`diff` 全等；:50 现量＝引语句（同发现取） | 对 | 对（两调用点现存；证据件存在且 §④ 在 :55） | 0 | **对上**（⚠ 见 §五① 行号互踩） |
| P02 | 在：`internal/agent/approval/doc.go:26-37`，`diff` 全等 | 对 | 前段全对；末条"Still owed"第 2 bullet 的"currently"**已过期** | 0 | **一处差**（§四①） |
| P03 | 在：`cmd/wisp/panel_inbound.go:11-14`，`diff` 全等 | 对 | 对（第二调用点＋仪器射程＝cmd/wisp） | 0 | **对上** |
| P04 | 在：`internal/panel/git.go:75-81`，`diff` 全等；:76 现量＝引语句 | 对 | 对（宿主/传输/router/调用点/Q-69/const :82 逐字描述均现验） | 0 | **对上**（git.go 无行号引用钉，rc=1） |
| P06 | 在：`cmd/wisp/config_readers_255.go:42-45`，`diff` 全等 | 对 | 对（两个 commit 均在；`TierOf(` 产码＝定义 tiers.go:93＋唯一调用 :235） | 0 | **对上** |
| P07 | 在：`cmd/wisp/config_reload.go:11-13`，`diff` 全等 | 对 | 对（**2715→2721 修正成立**；reloadOnce :152 内调 :153） | 0 | **对上** |
| P08 | 在：`cmd/wisp/config_reload.go:14-17`，`diff` 全等 | 对 | 对（**改引成立**：manager.go:172/:180 逐字；`ConfirmLocked != nil` 全仓只剩旧注释 :15） | 0 | **对上** |
| P09 | 在：`cmd/wisp/config_reload.go:18-21`，`diff` 全等 | 对 | 对（:115 赋值、manager.go:201-202 读侧）＋一处**归名差**（§四②） | 0 | **基本对上** |
| P10 | **不在**：自称块 8–15，引语实为 `firstrun.go:7-14` | 引语本身逐字对（对的是 7–14） | 对（loader.go **252 修正成立**、:238 确非 SaveFile；其余 7 项引用全对） | 0 | **对不上**（§四③） |
| P11 | 在：`cmd/wisp/resident_approval_windows.go:279-283`，`diff` 全等 | 对 | 对（唯一绑点 `resident_task_source_windows.go:309`；同文件 :322 确写 "called by startResidentTaskSource"） | 0 | **对上**（⚠ 见 §五①） |
| P12 | 在：`cmd/wisp/approval_reply.go:12-16`，`diff` 全等 | 对 | 替换句为真；risk 的 `.Veto(` 两处披露口径有一处瑕疵（§四④） | 0 | **基本对上** |
| P13 | 在：`internal/risk/assessor.go:28-32`，`diff` 全等 | 对 | 对（bridge.go:211-213/:924-927、panel_assets.go:66；横幅＝:9-36，待改句在横幅内且只动描述面） | 0 | **对上**（判不动已明标） |
| P39 | 在：`cmd/wisp/run.go:331-336`，`diff` 全等（含前导 tab）；:333 现量＝引语句 | 对 | 逐句对（`modeWrites` 读侧真 0；`ParseComposerRequest(` 调用真在 composer_dispatch.go:155） | 0 | **有条件对上**：事实句真，但照贴**撞红一枚现成判据**（§五② 头号） |

13 枚替换块行数（`sed -n` 现数，用于算行漂移）：P01 6→8(+2)｜P02 12→11(−1)｜P03 4→5(+1)｜P04 7→10(+3)｜
P06 4→5(+1)｜P07 3→4(+1)｜P08 4→6(+2)｜P09 4→5(+1)｜P10 8→10(+2)｜P11 5→6(+1)｜P12 5→7(+2)｜P13 5→7(+2)｜P39 6→9(+3)。

## 二、"判不动"三格复核（一行一条）

1. **P13（C19 冻结横幅边界）**：**明确标了**。P13 §4 原句「**"横幅内描述句是否算契约文本"的边界我判不动**（缺的尺＝人工/编排者裁定）」。
   我复跑：横幅真身 `assessor.go:9-36`（:10 「FROZEN CONTRACT — C19 rule set」），待改句 :28-32 在其内 ⇒ 位置判断对，
   且**没有**被写成"照做即可"。
2. **P04（`git.go:82` const 连带）**：**明确标了**。P04 §4「⚠ …**const 是否同步改＝产码面，本料不裁**」＋汇总表 §我判不动 #2。
   我复跑：`:82 const GitSwitchBlockedReason = "…还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主）…"` 逐字仍在
   ⇒ 替换文本对它的描述是准确转述，且未替作者裁。
3. **P04/P39（"hop 已落地"到票面验收的粒度）**：**明确标了**，但位置在汇总表 §我判不动 #3（「P04/P39 的"hop 已落地"粒度…判不动，
   缺的尺＝票 33/35 的验收口径」），**P39 自身 §4 未复述**。按代码事实层我逐句核过（下 §五②），不判票面粒度；此项属"标了但只在汇总处"。

## 三、三处引用修正的独立复验（各一把尺）

1. **P07 `PLAN.md` 2715→2721**：`docs/PLAN.md:2721` ＝ `### D36 — 配置模型（config.toml 全量 section 树 + 三档生效级别）`；
   `:2715` ＝ 「面板「数据与隐私」页…」别的内容；`:2726` ＝三档明细句。**修正成立**（同一对象＝D36 标题行）。
2. **P10 `loader.go` 238→252**（真身自定包＝`internal/config`）：`:252` ＝ `func SaveFile(path string, c *Config) error {`；
   `:238` ＝ `if ref := c.Voice.Realtime.APIKeyRef; ref != "" {`（确非 SaveFile）。**修正成立**。
3. **P08 引语改现行读法**：`internal/config/manager.go:172` ＝ `confirm := m.ConfirmLocked`（mu 下快照）、
   `:180` ＝ `approved := confirm != nil && confirm(c.section, c.keys)`；`git grep -F 'ConfirmLocked != nil'`（*.go 非 .scratch）
   **rc=0 唯一命中 `cmd/wisp/config_reload.go:15` 旧注释本身**。**修正成立**（替换文本引法＝:172＋:180 两步，逐字对）。

## 四、对不上的差（逐条具名；是它错 / 派单转述错）

1. **P02（它错，具名）——保留了一条已过期的"Still owed"分句**。替换文本末条 bullet 原文照旧：
   「delete loop.decideRisk's declared-L1/L2 refusal, **which currently kills fs.write before the bridge ever assesses it**」。
   现量：全仓产码 `agent.Options{` **唯一装配点**＝`cmd/wisp/run.go:995`（在 `func (rt *agentRuntime) execute` 内），
   `:1008 AdmitTask: rt.admitTask` 已把该钩子接上（:1004-1007 自述 D47 注册经此钩子）⇒ `loop.decideRisk` 的
   `l.opt.AdmitTask == nil` 分支（`internal/agent/loop.go:794`，在 :793-797 的 L1/L2 拒绝内）在产码里**不再被走到**：
   `AdmitTask:` 产码赋值点全仓只此一枚（范围尺：`*.go`、exclude `*_test.go`/`.scratch`，rc=0）。
   料的 §4 自报"仍欠"的现验只查了**分支存在**，没查**装订状态** ⇒ 该 bullet 属"该改未改"，照贴即把一句假话粘进 doc.go。
   建议落地腿改写该 bullet（或降级为"分支仍在代码里；`wisp run` 装订后已不可达"）。
2. **P09（它错·轻）——引短语归名差**。替换文本保留「票 223 AC#2 forbids（不许用「静默不生效」充当**这一档**）」。
   票面真身：`:35`（AC#2）写作「充当**第二档**」；「充当这一档」逐字出自 `:40`（**AC#7**）。票 223 存在、AC#2 存在、
   短语「静默不生效」存在（:35/:40 均有）⇒ 事实成立、逐字归名错一格。属既存句照抄（非本料新造），料的 risk 自己把归口记为
   `AC#7` ⇒ 替换文本与 risk 口径不一致，**两处都对不上**（文本挂 AC#2，risk 挂 AC#7）。
3. **P10（它错，具名）——块边界错位一行**。自称「原文逐字（块 8–15）／替换块＝8–15」，但引语 8 行逐字等于 `cmd/wisp/firstrun.go:7-14`
   （`:7` ＝ `// (NewDefaults, defaults.go:58 - …`；`:15` ＝ 单独的 `//` 段落分隔行）。`diff` 结果：`md[225-232]` 与
   `firstrun.go[7-14]` **IDENTICAL**；与 `[8-15]` 差 1 行（首行多余＋尾行 `//` 缺）。
   后果：**按自称行号粘贴会覆盖 :15 的分隔 `//` 并在 :7 留一行重复**；按引语实取（7–14）则自洽（替换块首行与 :7 逐字相同）。
   目标单行 `:11` 本身正确（`:11` 现量＝引语句）⇒ "13 枚锚零漂移（P07 也在 :11）"对 :11 这一句成立，错的是**块边界**。
   派单转述未提此行号（只提 loader.go 238→252），故**是料文错**。另：替换文本自称"8 行→10 行"与实取一致（7–14 也是 8 行）。
4. **P12（它错·轻）——risk 的 `.Veto(` 披露口径**。risk 写「`.Veto(` 产码两处＝approval_reply.go:363、resident_approval_windows.go:612」。
   现量（`git grep -F -e '.Veto('`，`*.go` excl `.scratch`）：`:363 s.live.h.Veto(corr)` 与 `:612 ra.cards.Veto(...)` 均非 `Gate.Veto`；
   `Gate.Veto` 的产码调用点是 **`internal/agent/approval/replies.go:463`（`g.Veto(Veto{…})`）——料文漏列**。
   替换句本身「.Veto, .DecideFromNative and .DecideFromPanel all have production call sites today」**仍为真**（replies.go:463/328/413/408/437，
   全为非测试产码），故只判 risk 披露瑕疵。另：票 201 `:10` 真身是三格行（`| 三种"人能答复"的入口 | **生产零调用者** | \`grep …\` ⇒ 只剩定义行 |`），
   替换文本只引前两格却写 "row verbatim" ⇒ 措辞略强于实际（两格逐字对，第三格未引）；`:73` 的"该尺读数已过期"登记存在、引用无误。
5. **派单转述的轻微不精确（我派单这句）**：派单与料文都把 `TestDispatcherSpellsNoRouteLiteralOfItsOwn` 标成
   `composer_dispatch_test.go:375`；真身 `func` 在 **:374**、regex `` `"panel\.[^"]*"` `` 在 **:380**（`:375` ＝ `root := …`）。
   按名取唯一、无歧义，属具名报回的转述误差（不影响结论）。另派单写的五禁词真身 `:714-716` 与料文 `:708-716` 两说均与盘上相符
   （:714 ＝禁词数组、:715 ＝Contains、:716 ＝红句含 "gated on Q-69"）。

## 五、料文未标、我另查到的（落地前必须处置）

① **跨处行号互踩（无仪器钉者仅列，有仪器钉者标红）**：
   - **P01 引 `composer_dispatch.go:155`（P39 的替换文本里）会被 P01 自己 +2 行推到 :157** ⇒ 两处不能各贴各的；
     同日另有一处 `cmd/wisp/panel_inbound_guards_35r3_test.go:40` 注释引 `composer_dispatch.go:153`（注释面，无仪器钉）。
   - **P02 引 `run.go:612/:748` 与 `resident_approval_windows.go:369`**：现量对上，但 P39(+3)/P11(+1) 落地后即失准（注释面引用，仪器不管）。
   - 各件 +N 行会让别处**注释/测试注释**里对它们的历史引用失真（我按 `file:line` 形状逐目标扫过 `*.go` 非 `.scratch`）：
     `git.go` 引用 rc=1（无）；`assessor.go` 有 3 处（`internal/agent/spill.go:56`、`internal/tools/bridge.go:918`、测试注释一处）；
     `config_reload.go` 有 7 处（`config_reload_223_test.go` 5 行注释、`config_sentences_223r2_test.go:23`、`resident_approval_risk_268_windows_test.go:281`）；
     `approval_reply.go` 有 6 处（含测试错误消息串）；`run.go`/`resident_approval_windows.go`/`composer_dispatch.go` 多处。
     **这些都不是判据，只有下面那一条是判据。**
② **P39 照贴会打红一枚现成判据（头号差，料文整个没提）**：
   - 判据真身 `cmd/wisp/config_receipt_255_test.go:179-211 TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`：
     从 `hotRowClaims`（`cmd/wisp/config_readers_255.go`）里把每条 `file.go:LINE [token]` 解析出来，
     `os.ReadFile` 那一行后要求 **该行仍含 token**（:207-211），另有 ≥8 条 checked cites 的地板（:220-223）。
   - 名册真身 `config_readers_255.go:109`（"agent" 行）：`"cmd/wisp/run.go:991 [cfg := rt.cfg] - …"`；
     我现取 `run.go:991` ＝ `cfg := rt.cfg`（今天该尺绿灯的基点）。
   - P39 换 331–336（6 行）为 9 行 ⇒ 之后行号 **+3**，`run.go:991` 漂到 **:994** ⇒ 该尺读 :991 不再含 token ⇒ **红**。
   - 同一 +3 还会推走 `run.go` 其余被引行（:423-424/:435/:424/:991/:1014，散在注释与图注），以及 `approval_reply_201_test.go:495-496`
     等注释引用（后者无仪器钉）。**处置建议**：落地腿要么把 P39 的替换块压回 6 行（不加行），要么同笔把名册里
     `run.go:991` 改成落位后的真行号（该行号属票 255 名册，改名册＝另一枚件的写面，须编排者裁）。
③ **P01 替换块的禁词体检是"整份文件字节"级**：真身尺 `composer_dispatch_test.go:708-718`（`os.ReadFile(composer_dispatch.go)`＋
   `strings.Contains`，注释也算）——P01 块 8 行现跑 `grep -F` 五禁词 **逐枚 0（rc=1）**，且"switching/switcher"不是禁词（真身数组只那五枚）。
④ **其余三把尺逐把点开＋对料文现跑**：
   - `:380`（`:375` 函数内）regex 尺只扫 `composer_dispatch.go` 且**跳过 `//` 行**（:381-388）⇒ 料文 0 命中（`"panel.` rc=1）、且替换块全为注释行；
   - `firstrun_257_test.go:430-438`：`"\nfunc "` 计数＝1＋签名逐字 ⇒ 料文命中 1/1 都在 §4 转述段，**替换块内 0**（`grep -F -c` 现跑），粘贴只增注释行 ⇒ 计数不动；
   - `resident_hotkey_258_test.go`（:73 落在读 `resident_windows.go` 的 err 分支；真钉是 :76 `config.LoadFile` 存在＋AST `startResidentBall` 参数数）⇒ 料文对这两串 **0 命中**（rc=1）；
   - `resident_grant_writer_265_windows_test.go:465` 起（`TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask`）＝AST 读
     `resident_task_source_windows.go` 的 `ra.bindResidentGrantLedger`／`src.submitTask`／`src.run`；另有 :450 起读 `Options.Grants` 选择子。
     AST 不吃注释；料文对 `ra.bindResidentGrantLedger` 1 命中（§4 转述）、`src.submitTask`/`ra.grants` 0（rc=1）——全在注释面，不撞。

## 六、13 处之外我另跑的（供编排者交叉核对，均现量）

- `.Handle(` 产码（`*.go` excl `*_test.go`/`.scratch`，rc=0）：`ComposerDispatch.Handle` 真调用**恰两处**＝
  `cmd/wisp/panel_host_windows.go:821`（`dispatchRaw` 内，`dispatchRaw` 定义 :814、调用 :803）与 `cmd/wisp/panel_inbound.go:163`；
  其余命中全是 `windows.Handle(...)`／http 类 Handler。函数值形状另查 `ComposerDispatch` 产码引用（rc=0）无 `= ….Handle` 赋值面。
- `main.go:115 case "panel-inbound":` 在；`panel_host_windows.go`：`:1 //go:build windows`、`:64` 引 `go-webview2`、`:386 webview2.NewWithOptions` 真建宿主；
  `resident_windows.go:151`→`panel_resident_windows.go:253 NewPanelManager(...)` 装配链在。
- 证据件 `docs/evidence/s1/33-minimal-inbound-hop-r1.md` 存在＝147 行（与料文同数），`§④` 在 `:55`。
- 票面存在性（`git ls-tree` 现查）：33／35／37／114／186／198(−done)／201／223／246／255／18(−done)／19(−done) 全在；
  Q-69 真身＝`docs/reports/pending-and-issues.md:1109`（另 `HANDOVER.md:1847`、`composer_dispatch_test.go:716` 红句）。
- `AppliedSteps` 在 `internal/tools/fs_write.go` 14 行；`ToolsCancelBus` 产码＝`run.go:749`；`CheckAndReload` 产码＝`config_reload.go:153`＋`cmd/balldebug/main.go:243`；
  `config_reload.go:105 func (rt *agentRuntime) startConfigReload()`（:114/:115 赋值、:119 起 tick ⇒ 料文"before the tick is spawned"对）。
- 五禁词对**整份料文** `grep -F` 逐枚 0（rc=1）；对 P01 替换块（md 30-37，8 行现数）逐枚 0（rc=1）。

**总判**：13 处里 **10 处clean对上**（P01/P03/P04/P06/P07/P08/P11/P13＋P09/P12"基本对上"各带一处具名轻差），
**1 处对不上（P10 块边界，必须改锚或改粘贴范围）**，**1 处带红（P02 过期 bullet，必须改文案）**，**1 处有条件（P39 事实句真但撞红 255 名册行号尺，必须改落地方式）**。
未跑 Go 测试（禁区），"红/绿"均按判据源码读法判定并给出 file:line 依据。
