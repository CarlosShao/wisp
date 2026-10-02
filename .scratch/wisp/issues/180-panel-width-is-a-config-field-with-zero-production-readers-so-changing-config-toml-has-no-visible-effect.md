# 180 — `[panel] width` 是一枚**生产码零读者**的配置项：用户在 `config.toml` 里改它、或我们改 `default:"640"`，**界面上什么都不会变**（"提示存在但结果永不受影响"那一族在配置面上的复发）

- Status: **ready-for-agent（只读普查＋一支小落地）**。
- 来路：09-28 10:2x，前端那侧递给 owner 的"要拍板三件事"里第 ① 件（"`[panel] width` 默认 640 — 你要多少？改 `config.toml` 还是让 Go 改默认值？"）。我按盘上现量裁定：**这一件不该摆给他**——因为**两条路今天都没有可见效果**，先修效果再谈数值。台账 `A360`。
- 关联：票 92（composer 的 mode/附件/工作区那一格，同族"面板要真值不要自造"）· 票 145（快照扩成真载体）· 票 147／139／97（"字段存在、无人读"那一族）· `AGENTS §1.2` 与台账里那条"一条提示存在、但结果永远不受它影响的检查，比没有检查更坏"

## 现量（锚 `039efb47`，09-28 10:2x 本程现跑，每把尺可复制）

1. 字段在场：`internal/config/schema.go:527-528` 逐字
   `// Width is the panel width in px.` ＋ `Width int \`toml:"width" default:"640"\``（尺：`grep -rn "Width" --include=*.go internal/config/ | grep -v _test`）。
2. **读者＝零枚**：尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ **空输出**。
   ⇒ 今天**没有任何一条路径**把这枚数送到窗口尺寸上：改 `config.toml` 无效、改 `default` 也无效。
3. 它也不是"待接线的公开缺口"被别处登记过：尺 `grep -rn "width" internal/config/unwired_test.go` ⇒ 只有 `:217` 那行**测试用的 TOML 样本**（`width = 800`），**不在"未接线名册"的断言里**。
   ⚠ 这一条要写手复核：`unwired_test.go` 的语义到底是"登记未接线字段"还是别的；若它本应登记，那本票的第二格是"补登记"而不是"新发现"。

## 为什么值得做（不做会怎样）

最坏后果不是"宽度不对"，是**它假装可配**：用户（这里是 owner 自己）照文档/照配置树去设一个值，界面不动，而他**没有任何一处能读到"这项今天不生效"**——他会归因成"我设错了"或"这产品坏了"。同一族我们已经在提示音、审批窗口、SLO 步上各收过一次。

## AC（每格都要答"这一发在未修码上响不响"）

- [x] **AC#1 先判归属，再谈修法**（09-28 17:5x 编排者复跑对格后勾：三问答完且我今天自己重跑过——**没有任何人**决定面板尺寸，因为这棵树里**根本没有面板窗口**：`internal/panel/pump.go:15-17` 逐字"no WebView2 host"、`grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test`＝**0 命中**、那五枚 `[panel]` 字段的生产读取方各 **0**；规格要它生效＝`PLAN.md:2743` 逐字 `hot`；唯一真窗口是球（`internal/ball/ball_windows.go:233`，与 `[panel] width` 无关）＝账 `A391`）：逐枚答三问并落表——① 面板窗口尺寸今天由**谁**决定（现读 `internal/panel/` 里创建窗口/WebView 的那处，给出 file:line；尺：`grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test`）；② `cfg.Panel.Width` 到它之间**缺哪一跳**；③ D36 配置树里这一项的**规格出处**是哪一行（尺：`grep -n "\[panel\]" -A 6 docs/PLAN.md` 与 `grep -rn "width" docs/specs/SPEC-03*.md docs/specs/SPEC-08*.md`）——**规格有没有要求它生效**，决定本票是"补实现"还是"改规格文字（＝人工批准）"。
- [ ] **AC#2 正向落地（若 AC#1 判"规格要求生效"）**：把那一跳接上，并给一发**常驻判据**：设成两个不同值 ⇒ 面板侧读到的尺寸**随之变**（不许只断"字段非零"）。⚠ 判据要能区分"接上了"与"又抄了一遍默认值"。
- [ ] **AC#3 反向判据（与 AC#2 成对，缺一形＝装饰）**：**不许**用"把默认值改掉"来交差——把 `default:"640"` 改成别的数而**不接读者**，本票判不通过（那只是把一枚装饰换成另一枚）。若 AC#1 判"规格今天不要求生效"，则正解是**响亮地登记**：在配置项注释＋`docs/reports/pending-and-issues.md` 里写明"此项 RESERVED，设了不生效，缺口在票 N"，**不许静默留着**。
- [ ] **AC#4 顺带查同族**：`internal/config/schema.go` 里**还有几枚**"有 `default` 但生产零读者"的字段（尺：对每一枚字段名跑 `grep -rn "\.<字段名>" --include=*.go internal/ cmd/ | grep -v _test.go`，逐枚记 0/非 0）。⚠ 枚数只许现量、不许目测；本票**不修**它们，只出名单并逐枚归口（新票或既有票），⚠ 不许写"以后加固"。
- [ ] **AC#5 契约轴**：`docs/PLAN.md`（含 D36 那一节与 `:1531`／`:1532` 两行 DEFERRED）、`docs/specs/**`、`thresholds.go`／golden／`allowlist.txt` 一字节不许动；`internal/config/schema.go` 的**默认数值本身今天不许改**（要改数值是 owner 的一句话，见 Q-64 之外的另一枚待答项——若 AC#1 判"规格要求生效"，数值仍保持 640，**本票只接效果、不改数**）。
- [ ] **AC#6 门禁**：逐包 `go test -count=1 ./internal/config/ ./internal/panel/`（⚠ `./cmd/wisp/` 在本机测不到任何东西＝票 98，别把它算进你的红；CLI 那一面只能走 `-overlay` 台件）＋`sh scripts/d22scan.sh`＋名册两向 `comm` 差集；⚠ **最终读数取在最后一枚 commit 之后**。

## 本票**不**解决

- 不裁"面板应该多宽"——那是 owner 的一句话，且**在效果接上之前问他没有意义**（他会答一个不生效的数）。
- 不改任何前端文件（`frontend/**` 归别家，不读不写不引）。
- AC#4 只出名单不修。

## Progress log

- 09-28 10:2x 编排者立票：上面三把尺本程现跑（锚 `039efb47`）。未派。
- 09-28 14:4x `180-c1`（只读普查，与票 182 合派一程；表 `docs/evidence/s1/180-182-panel-fields-census-c1.md` §2–§3）：
  **AC 框一枚没勾、产码零字节未动**。现跑复核（起手 HEAD `e9ef94d0`，非票面锚 `039efb47`）——
  ① 字段仍在 `internal/config/schema.go:527-528`（`default:"640"`）；
  ② 尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ **仍然空输出（0 读者）**；
  `PanelSection` 生产里唯一被碰到的是 `internal/config/manager.go:211` 的**整块指针拷贝**，不读值不生效；
  ③ 三档结论＝**「规格要求生效」**：`PLAN.md:2726` 逐字「`hot`（立即生效）」＋`PLAN.md:2743` 把 `width` 列进 `[panel]` 的 `hot` 行
  ＋`SPEC-03-config-secrets-envs.md:39` 同行；SPEC-08 里 `width` 唯一命中是 `:218` 的 SVG `stroke-width`，与面板窗口无关。
  **缺的不是"一跳"是整条链**：这棵树今天**没有面板窗口**（`internal/panel/doc.go:16`／`pump.go:16`／`cmd/wisp/run.go:438`
  三处注释逐字指票 33/35 未开工），⇒ 票面 AC#1 的"谁决定尺寸"答案是**没有人**。
  AC#4 同族现量：带 `default:` 的行 **69**、去重字段名 **61**、**零生产读者 18 枚**（下界，粗尺同名歧义见表 §3），
  逐枚名与 struct 见表；⚠ **这 18 枚全部已在票 83（done）的 AC#1 全量表里登记并归口**
  （`issues/83-...md:155` 逐字把 `panel.{enabled,width,height,keep_alive_in_session,scale}` 归给 📋 票 33/34/36）。
  现量第 3 条的复核：`unwired_test.go` 的"名册"（`TestEveryLockedSectionKeyIsAccountedFor` `:311`）
  **射程只有 `risk`/`fs`/`net`/`plugins` 四枚 locked section** ⇒ `[panel] width` 不在册是设计如此，
  本票第二格既不是"补登记"也不是"新发现"；**落点本程不选，交编排者**。
- 09-28 17:5x `180-a1`（只读普查·派单 §B；本程锚 `38fc7c0e`，非票面锚 `039efb47`）：
  **AC 框一枚未勾、产码零字节未动**；只跑尺与读码，表在 `docs/evidence/s1/180-panel-width-ownership-a1.md`。
  ① **今天没有人决定面板窗口尺寸，因为本树没有面板窗口**——
  `internal/panel/doc.go:16`／`internal/panel/pump.go:15-17`／`cmd/wisp/run.go:438` 三处注释同指票 33/35 未落；
  尺 `grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test` **空输出 rc=1**；
  `go.mod` 里 webview 依赖**零枚**；仓内唯一真窗口是球，其尺寸走 `internal/ball/ball_windows.go:233` 的 `opts.SizePx`
  （默认与钳位 `:64`/`:144-151`），与 `[panel] width` 无关系——**别把球那一格当成面板的归属**。
  ② `default:"640"` 那一族＝`internal/config/schema.go:525-535` 的五枚字段
  （`enabled` `width` `height` `keep_alive_in_session` `scale`），生产读取方**逐枚现量各 0**
  （五把尺各空输出 rc=1，原文贴证据件 §2；尺射程只有 `internal/ cmd/ tools/`，`frontend/**` 与 `design/**` 不入任何零命中宣称）。
  生产里唯一碰到这枚 struct 的两处都不"读值去行动"：`internal/config/manager.go:211` 整块指针拷贝＋`:217` 整块比较；
  `internal/config/defaults.go:60` 的反射只把 `default:` 标签**写进**字段。
  连热加载的回执也停在字符串上：`rep.Hot` 全树命中三处全是 append（`manager.go:219`/`:278`/`:334`），无一处取值行动。
  ③ **规格要求它生效**（不是真空，故本票＝补实现那一支，不是改规格文字）：
  `docs/PLAN.md:2726` 逐字「**三档生效级别**：`hot`（立即生效）· `reload`（需重载子系统，如换 ASR 模型）· `restart`（需重启进程）。」
  ＋`docs/PLAN.md:2743` 逐字「| `[panel]` | `enabled` `width` `height` `keep_alive_in_session`(true) `scale` | `hot` |」
  ＋`docs/specs/SPEC-03-config-secrets-envs.md:39` 同行（`width(int)=640` … `hot`）。
  ⚠ 但要求**只由 `hot` 那一格承载**：`docs/specs/SPEC-08-ui-ball-panel.md:143-154`（§5.1 宿主）只写单例窗口/资源加载/resync，
  **没有一句把窗口尺寸指到 `cfg.Panel.Width`**——宿主侧文本是缺的，这条归编排者裁。
  ⚠ 风险栏（照派单 §B 末条预写；`AC#2`/`AC#3` 本程**未做**）：**不许用"把 `default:"640"` 改掉"交差**；
  票面问的"缺哪一跳"今天答案是**整条链**——宿主（票 33）与载体（票 35）都不在，
  且尺寸这一维在快照名册里也不存在（`internal/panel/composer.go:57-62` 只有四枚 key）。
- 09-28 17:5x 编排者收 `180-a1`（只读普查，两枚 commit `cb14a8b2`／`35abbc40` 只动票面与它自己证据件）：`AC#1` 勾（我复跑对格）。**排程结论（这一条比表更重要）**：面板尺寸这一维缺的**不是"接线那一行"，是整整一个窗口**——今天这棵树里没有面板窗口，所以 `AC#2`（正向落地）**排在票 33 的 WebView2 宿主之后**（H2／H3），本轮不排进队列；为免它再空一次，我已把"宿主必须读 `[panel] width/height/scale` 并真去设边界"追加成**票 33 的新格 `AC#10`**。另记一枚过期注释（`internal/panel/pump.go:15-17` 逐字写着"ticket 33 is unclaimed"——**33 早已被认领、片 A／片 B 都已落**）：归 `145-r2` 交件后的下一枚 `internal/panel` 写腿顺带更正，**现在不动那枚文件**（它在别人写面上）。

## 10-02 现量（追加节，正文与上面所有 AC 框一字未改）

- 10-02 10:2x `180-a1`（**第二程**·只读全名册普查，非 09-28 那一程；派单五问逐问落表；
  表 `.scratch/wisp/probes/180/a1/census.md`，588 行／53744 字节，占位符尺 **0** 命中）：
  **AC 框一枚未碰、产码零字节未动**。起手 HEAD `00e7efe`（10-02 09:57:47 +08）、
  终态 HEAD `f51cfef4`（10:27:45 +08）；
  尺 `git diff --stat 00e7efe..HEAD -- internal/config internal/panel cmd/wisp go.mod` ⇒ **空输出**
  ⇒ 本程引用的 12 处承重行号在两枚 HEAD 上**同立**（票面 `schema.go:527-528` 已漂到 `:532`，
  派单 `unwired.go:57-101` 已漂到 `:60-97`）。
  1. **全名册现量（问 1）**：段数 **18 枚 section ＋ 1 枚根表裸键 `schema_version`**
     （派单那句"18 个 section"**成立**，两把独立尺：`Config` struct 段字段枚数＝18；
     票 198 首建真件 33 个表头＝18 顶层＋15 子表）。叶子键 **150 枚**＝静态 115 ＋ 动态子表模板 35
     （`providers.<id>` 10 ＋ `models.<id>` 21 ＋ `plugins.<id>` 4）；
     静态 115 与首建真件的 115 行 `key =` 逐枚 `comm` 只差 `models.local_override`（空 map 不落盘）。
     带 `default:` 标签 **69** 枚（与 09-28 `180-c1` 独立对上）。
  2. **七档判读（150 枚全表见证据件 §1.5）**：
     R 有生产读者 **50** ／ B 读了只作旁注 **4** ／ W 类型在场·生产零赋值 **5**（`price` 五枚）／
     S 仅旁路调试程序 `cmd/balldebug` **4**（`[hotkey]` 全段）／
     G 无人读但 `unwiredKeys` 拦成加载错误 **6** ／ L 无人读仅登记 `lockedKeyDisposition` not built **5** ／
     **D 无人读且名册外 76**。⇒ **无人读或读了不做事 100／150**；**既无人读又不在任何名册 76／150**。
     整段无人读的段：`[voice]` 21／21、`[ball]` 7／7、`[audio]` 4／4、`[session]` 3／3、
     `[memory]` 4／4、`[cost]` 4／4、`[observe]` 4／4、`[privacy]` 5／5、`[app]` 4／4、`[panel]` 5／5。
  3. **`unwiredKeys` 覆盖率（问 2，本程最值钱那一格）**：名册 **6 枚**（`unwired.go:60-97`）。
     四个分母四个数——无人读的键里被任一名册登记者 **11／100＝11%**；
     会拦下来报错的 **6／100＝6%**；锁定四段之内 **11／11＝100%**；
     **锁定四段之外 0／89＝0%**。⇒ **"名册里安静"确实不等于"全仓没有哑键"，差的不是几枚是 89 枚。**
     完整性钉 `TestEveryLockedSectionKeyIsAccountedFor`（`unwired_test.go:311`，本程读到函数体为止）
     对四段各跑一次 `collect()`，我照它的逻辑复算 ⇒ 实际收集 **16 枚**路径，
     其中 `plugins` 只 1 枚：`PluginsSection.Entries` **没有 `toml` 标签**，而 `collect()` 首句
     `if tag == "-" || tag == "" { continue }` ⇒ 名册里那五条 `plugins.<id>.*` 行**这枚钉永远查不到**。
     且钉只问"有没有一句话"，`consumed:` 分支**不核实消费者是否在场**（本程逐枚读了 6 条 `consumed:`，
     都能在码上对上具名行 ⇒ 这一族今天没抓到假话）。
  4. **生效级别三档（问 3）**：作为**数据不存在**，作为**硬编码分支存在**。
     `type Tier string` ＋ `TierHot`/`TierReload`/`TierRestart`（`schema.go:30-44`）**零读者**
     （尺：`grep -rn "Tier\b|TierHot|TierReload|TierRestart" --include=*.go cmd internal tools | grep -v _test.go`
     ⇒ 命中全是别人的 `Tier`：`panel.EffectiveTier`、`risk.Tier`、`projctx.Tier`）；
     `schema.go` 里只有 `toml:` 与 `default:` 两种标签，**没有 `tier:`** ⇒ 级别不在字段上，
     只在 `plan()`／`planApp`（`manager.go:338-357`）／`planVoice`（`:361-417`）的分支里，**粒度是段不是键**。
     ⚠ 本程新量一枚：**reload 档在出货进程里没有听众**——`Manager.OnReload` 的非测试赋值点全仓只有
     `cmd/balldebug/main.go:244`，`cmd/wisp/config_reload.go` 只挂 `ConfirmLocked`（`:114`）与
     `OnRestartPending`（`:115`）⇒ 换 ASR/TTS 模型不会触发模型加载/卸载，
     而 SPEC-03:134 逐字要求"reload 段触发重载事件"。
     另：面板侧另有**四档**词汇 `EffectiveTier`（`config_handlers.go:196-203`：now/next_task/restart/not_applied），
     与 D36 三档**不同名也不同数**，且只走设置写入回执那一条路。
  5. **"此项今天不生效"有没有出口（问 4）**：**读侧盘上没有这句话；写侧已经有现成的一档。**
     现成＝`SettingWriteResult.Tier` ＝ `EffectiveNotApplied` ⇒ `tierSentence` 回"这一项没有被应用。"
     （`config_handlers.go:196-203`、`:210-217`、`:428-441`，票 248 落的）。
     但它只在 `config.set` 被**受理**之后存在，而白名单只有 **7 枚字段**（`:57-69`），
     `[panel] width` 过不了受理关（`:321-326` `unlisted-field`）⇒ 手改 `config.toml` 的人**连那句话都走不到**。
     快照侧也不能：`Snapshot`（`composer.go:57-92`）六个 `json:` 键里没有一栏承载配置键状态；
     `config.get` 走 `renderSettingsView`（`:445-462`）只说可读性／服务商数／凭据状态／聊天模型。
     最近的可扩展点具名到字段：`SettingsView`（`:142-156`）／`renderSettingsView`（`:445-462`）／
     `SettingWriteResult.Tier`（`:210-217`）／登记面 `lockedKeyDisposition`（`unwired.go:119-147`）。
     ⛔ 前三处属 **C17／契约面**（快照键集与方法白名单）——**本程不裁该不该加、加哪一栏**，交编排者落批准记录。
  6. **`[panel] width` 断在哪一环（问 5）**：十环逐点（证据件 §5.1）。
     值走完了"文件→`parse.go`/`loader.go`→`defaults.go:60`→`Manager.cur`（`manager.go:106`）→
     `Config()` 深拷贝（`:115`）"，**断在第 7 环：面板宿主从来拿不到配置对象**——
     `PanelManager`（`cmd/wisp/panel_host_windows.go:137-167`）没有 config 字段，
     构造函数 `NewPanelManager(disp, assets, dataPath)`（`:175-177`）不收配置，
     生产调用点 `cmd/wisp/panel_resident_windows.go:204` 也只传这三样；
     第 8 环尺寸在建窗那一行是**写死的字面量** `Width: 420, Height: 260`（`:304-305`），
     第 10 环没有第二条路（`SetBounds|MoveWindow|SetWindowPos|Resize(` 在本包外非测试码只命中
     `:45` 的注释）。
  7. **本程推翻票面／09-28 两程三句**（全表见证据件 §8）：
     ① 09-28 两程的"**这棵树今天没有面板窗口**"**已过期**——`go.mod:19` 已有
     `github.com/jchv/go-webview2`，`cmd/wisp/panel_host_windows.go` 与 `panel_resident_windows.go` 在场且真建窗。
     ⇒ 票面 AC#1 那问"谁决定尺寸"的答案从**没有人**变成**`:304-305` 那两枚字面量**；
     ⇒ 编排者 09-28 那句排程结论（AC#2 排在票 33 的宿主之后）**前置条件已满足**，这一格值得重排。
     ② **比票面"他没有任何一处能读到这项今天不生效"更坏一档**：手改 `[panel] width` 不是没回应，
     是收到一句**主动误报**——`manager.go:285` 把 panel 段列进 `rep.Hot`，
     `cmd/wisp/config_reload.go:170-172` 逐字打 stdout"配置热加载：这些段已立即生效（D36 立即档）：[panel]"。
     同一句误报覆盖 `[panel]` **全部五枚**（不只 `width`）。
     ③ 票面 AC#4 那句"还有几枚"：本程现量 **76 枚名册外哑键**（不是 09-28 的 18 枚）——
     两个数不矛盾，`180-c1` 那一遍只数了带 `default:` 标签的字段，本程按叶子键全量数。
     ④ 票面 09-28 记录"这 18 枚全部已在票 83 的 AC#1 表里登记并归口"：`issues/83-...md:155` 那一行
     本程**未读到**，按派单硬边界一律当待验断言引用，不据它追认任何格。
  8. **本程没做（明写，不写成"以后加固"）**：没跑任何编译／测试类尺（`198-v1` 在整包跑 `cmd/wisp`）
     ⇒ 全部枚数出自文本尺与 Python 走查，没有一把是编译器给的；
     没做 76 枚的逐枚归口（要读全工单池，超一条只读腿射程）；
     未读、未引 `frontend/**`／`design/**` 任何内容 ⇒ "读侧没有出口"这句的射程**只到 Go 侧为止**。
