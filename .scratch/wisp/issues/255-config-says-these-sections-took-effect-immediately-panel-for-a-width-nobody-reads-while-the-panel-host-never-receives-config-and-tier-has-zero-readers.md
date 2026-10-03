# 票 255 — 改了 `[panel]` 那几项，程序**主动告诉你"这些段已立即生效：[panel]"**，而面板宿主根本拿不到配置对象：这句误报、三档生效级别零读者、段外 76 枚键无人登记

**立票时刻**：2026-10-02 10:3x +08，锚点 HEAD `fa575b49`（`dev`）
**来路**：只读腿 `180-a1`（`.scratch/wisp/probes/180/a1/census.md`，**588 行 / 53,744 字节**，九节齐、占位符 0）；台账 `A524`。票 180 原题为「`[panel] width` 生产零读者」——**现量比原题更糟**。
**性质**：⚠ **AC#1 算产品行为**（程序自己说了一句不成立的话）；AC#2／AC#3 是仪器与治理；AC#4 才是"让它真生效"那件功能。

## 现量（编排者本机自跑，2026-10-02 10:3x，逐字）

1. **误报那一行的产生处**：`internal/config/manager.go:280-288` 是一张**热应用段表**，`{"panel", &cur.Panel, &fresh.Panel, func(){...}}` 就在其中（我现读到 `:285` 那行确为 `panel`），其后 `for _, s := range rest {` 逐段应用并打那句"这些段已立即生效：…[panel]…"。⇒ **它报的是"内存里的值换掉了"，不是"有人按新值做了事"**。
2. **面板宿主确实拿不到配置**：`NewPanelManager(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string)`（我现跑 `grep -rn 'func NewPanelManager'`，全仓**只有这一枚定义**，⛔ 无配置参数）；窗口尺寸写死在 `cmd/wisp/panel_host_windows.go:304-306` —— 逐字 `Width: 420,` ／ `Height: 260,`。⇒ 断点**在"宿主拿不到配置对象"这一环**，不在解析、不在落盘。
3. **三档生效级别只是数据上的样子货**：`internal/config/schema.go:37-43` 定义了 `TierReload`／`TierRestart`（第三枚同处），**非测试读者：零**（我 `grep` 只命中定义与注释）；真正的档位判断硬编码在 `plan()`／`planApp`／`planVoice` 里，**粒度只到"段"**。⇒ 与 `AGENTS §4` 索引里 D36 那句"三档生效级别"对不上：**文档说三档按键，码上只有按段硬编码**。
4. **名册覆盖率（这条最要紧的账）**：`unwiredKeys` 实测 **6 枚**（`internal/config/unwired.go:60-97`；我派单写的 `57-101` 已漂）。叶子键总数 **150**＝静态 115＋动态模板 35；**没人读的键里被任一名册登记的只有 11/100＝11%**，锁定四段内 **100%**，**段外 0/89＝0%**。
5. 顺带一枚排程事实（它推翻票面 09-28 两程那句）：**"这棵树没有面板窗口"已过期**——`go.mod:19` 已带 `go-webview2` ⇒ 票 180 的 AC#2 排程前置**已满足**。

## 要建什么

- [ ] **AC#1 那句"已立即生效"只许说真话**：判据＝stdout 里"这些段已立即生效"那一段**只列真有生产读者的段**；**正控**＝手改 `[panel] width`  ⇒ 那句话**不许出现 `[panel]`**（要么不出现、要么换成诚实形如"值已换、但当前无组件应用它"）；**反控**＝真有读者的段（逐名给出证据）**仍要出现**。⛔ 不许把整句删掉——那是把热加载的可见性摘了，本仓定式：⛔ 摘尺不许当修尺。
- [x] **AC#2 三档生效级别要么落成可读数据，要么具名降级口径**：ⓐ 让 `Tier` 真被读者消费（谁在哪一档读它，具名到函数），或 ⓑ **具名登记"粒度只到段"**并给一把会响的尺（新增一枚键若其段的档位与实现不符 ⇒ 响）。⛔ **不许改 `docs/PLAN.md` 的文字**（改契约＝人工批准；本票只裁码与仪器）。**〔20:0x 编排者翻勾：ⓑ 落地＝`dd92bb92`（TierRegistry 16 段级＋app 4／voice 24 键级照第 47 行更正）；非实现者 `255-v1`（`45e33282`）第 3 发实证同源守卫——`"llm"` 改档 ⇒ panic 逐字 `config: section llm is hot-applied by plan() but not registered as "hot"`、热加载不再照常工作；判语＝成立（本写面内；`TierOf` 回执消费归 `255-r2`）。凭据＝`docs/evidence/s1/255-tier-registry-v1.md`〕**
- [x] **AC#3 段外那 76 枚键要一把**会响的尺**，⛔ 不要逐枚归口**：判据＝**新增一枚叶子键，既没有生产读者、也没有进任何名册 ⇒ 指名那一步必须红**（能力形，不是词表；正控必配"补登记后不响"）。⚠ 明写不做：把 76 枚逐枚判"该不该生效"是死活，交编排者按段抽判。**〔20:0x 编排者翻勾：`255-v1` 三处验真（种哑段红句逐字点名 `section "dumb"`／删行数钉双保险齐红／voice 假键活体红）＋正控齐＋两枚例外承重活体证明（清空例外名单红点名、非碰巧不红）。范围披露：`[app]` 新增键不在键级尺射程（固定名单非反射走查）——我裁＝记范围不算欠账，补尺（一行循环）归 `255-r2` 顺手做〕**
- [ ] **AC#4 让 `[panel] width` 真生效（这才是 owner 那句"改了要有反应"）**：落点＝**装配根把配置递给面板宿主**（既有裁定：装配根是唯一接缝），⛔ 不许新开 `panel→config` 依赖边、⛔ 不许在 `PanelManager` 里自己去读盘；写死的 `420×260` 变成"配置缺省时的值"而不是唯一来源。判据＝改 width ⇒ 真窗口宽度随之变（⚠ 这一格属**只有本机可量**那一族，别伪装成 CI 测过）。

## 禁区

- ⛔ 不许动 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt` 一字节；⛔ 不许放宽任何断言换绿；⛔ 不许 `t.Skip`；⛔ 不许把 SKIP 读成通过。
- ⛔ **凡要新增面板快照字段／新方法名＝C17 契约面**（`180-a1` 具名点出 `SettingsView`／`renderSettingsView`／`Snapshot` 三处都属这一面）⇒ **先停手上报**，由编排者落 `A##` 批准记录后才许动。
- ⛔ 草稿落库属 `docs/specs/SPEC-02` 镜像面；⛔ `voice` 那 21 枚键的终判归语音腿；`OnReload` 出货进程没听众（非测试赋值点全仓只有 `cmd/balldebug:244`）这一格**归票 42**，不许并进本票。
- `frontend/**`／`design/**` **既不读也不引**（⛔ 所以"读侧没有出口"这一句只到 Go 侧为止，界面侧那一半归 owner 带话）；grep/find 必写显式根（`cmd internal tools docs scripts .scratch`）。
- Git 纪律：只 commit 不 push；显式 pathspec；⛔ `git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 排程

- 写面＝`internal/config`＋`cmd/wisp`（AC#1／AC#2／AC#4 三格都碰）⇒ **此刻 `248-v1` 正在整包跑 `cmd/wisp`** ⇒ **按住**，起跑判据：`248-v1` 退出＋`git status --porcelain cmd/wisp internal/config` 为空。
- AC#3 那把尺属 `scripts/`＋`internal/config` 测试，可与 `254-r1`（只写 `scripts/`）串行不并发。
- ⛔ 拆腿建议两枚：**255-r1**＝AC#1＋AC#3（说实话＋会响的尺，纯 Go 测试与文案）；**255-r2**＝AC#4（真窗口尺寸，本机可量那一族）；AC#2 待我裁 ⓐ／ⓑ 之后再排。

## 编排者裁定（2026-10-02 11:00:02+0800，台账 `A527`；料＝只读腿 `255-a1` 的 `.scratch/wisp/probes/255/a1/census.md` 404 行／47,058 字节，我自己在 HEAD `59381f4a` 上逐条复认过下面五行）

**AC#2 裁＝ⓑ（具名登记＋一把会响的尺），ⓐ 不派。** 复认过的形状（**"生效级别"今天是四套词表，唯一真在跑的那套通篇不写 "tier" 这个词**）：

| | 词表 | 我复认到的 |
|---|---|---|
| T1 | `config.Tier` 三值（`internal/config/schema.go:31/36/39/43`） | 非测试写者 0、读者 0 |
| T2 | `restartTierKeys`（`cmd/wisp/config_reload.go:296-301`） | 段名/键路径硬编码；**`:297` 注释自称 "and by a test"，而 `grep -rn restartTierKeys --include=*_test.go cmd internal`＝0 命中＝那句没有 test** |
| T3 | `panel.EffectiveTier` 四值（`internal/panel/config_handlers.go:196-202`） | 写点全在 `cmd/wisp/panel_config_store.go`；`EffectiveNow`/`EffectiveNextTask` 的**非测试写者 0 枚**（实际只活两枚值） |
| **T4** | **`Report{Hot,Reload,Restart}`（`internal/config/manager.go:87-97`）** | **这才是真在跑的那套**；档名不在任何字符串里，在"段名被 append 进哪一枚切片"⇒ 按档名符号量的尺天生量不到它 |

- ⛔ **登记表逐"键"登记，不许逐字抄我票面那句"粒度只到段"**：`planApp:349/:355` 与 `planVoice:412/:415` 已是**一段双档的键级分档**（我 sed 复认），照票面原文登记＝**把诚实的码判红**——**这一处错记在我自己立的票面上，不改原话，就地具名更正。**
- ⛔ **`config.Tier` 那三枚常量不许"顺手接上线"了事**：它连的是死类型，买到的是编译期好看、运行时零改变。

**AC#5（本格由本裁定新增，⛔ 未勾，判语要非实现者裁）＝同一项设置在两套词表里答相反**：现量＝`cmd/wisp/panel_config_store.go:227` 逐字 `res.Tier = panel.EffectiveRestart`，而 `:222-229` 那一段**前面没有任何按键/按段的判断**；同时 `internal/config/manager.go:281` 把 `llm` 段列在**热应用段表**里。⇒ 面板对 `[llm]` 说"重启才生效"，热加载那条路把它当立即档，**两套之间零行码互相翻译**。判据＝**写入回执那句"什么时候生效"必须由 ⓑ 那张登记表产出**；正控＝改一枚表里标 hot 的键 ⇒ 回执**不许**出现"重启才生效"；反控＝`app.autostart` 那类 restart 键 ⇒ **必须**说重启。⛔ 不许为了让这格好看而把 `EffectiveNow`/`EffectiveNextTask` 那两枚零写者的值"用起来"——那是造出第五套；要么登记进同一张表，要么具名留着。

**排程（收 `255-a1` 给的两枚耦合点）**：ⓑ 与 AC#3 **共用同一枚反射枚举器**（"新增一枚哑键必红"那枚正控要真改 `schema.go` 的结构声明，纯测试侧造不出）⇒ **`255-r1`＝AC#2-ⓑ＋AC#3**（写面 `internal/config`）；**`255-r2`＝AC#1＋AC#5**（写面 `cmd/wisp`＋`internal/panel` 的回执与文案）；**`255-r3`＝AC#4**（真窗尺寸，〔只有本机可量〕那一族）。⛔ 上面"拆腿建议两枚"那行**被本段取代**，原话留着不删。AC#2 本格**不勾**：勾要非实现者裁，且 ⓑ 落地本身归 `255-r1`。

**⚠ 行号更正（2026-10-02 13:1x，台账 `A531`；原话不改，就地追加）**：上面 AC#5 那格写的 `cmd/wisp/panel_config_store.go:227` **差一行**——真身是 **`:228`** 逐字 `res.Tier = panel.EffectiveRestart`（`:227` 逐字 `res.Written = written`）。尺＝只读腿 `e2e-panel-1` 跑了 `grep` 与 `git cat-file blob HEAD:` **两把一致**；我自己在 `cmd/wisp` porcelain＝0 的工作树上 `sed -n '226,229p'` 复认（该包干净⇒工作树与 HEAD 同字节）——⛔ 我没有独立跑那枚 `cat-file` 尺，归属写清楚。**结论层不变**（"无条件硬填、不看键不看段"这条来自我 10:51 亲手 sed 过的 `:222-232` 整段），错的只是我引的行号。另补一枚现量（只读腿 `e2e-panel-1`）：`internal/config/schema.go:532` 的 `[panel] width` 默认值 **640**，而宿主写死 **420**（`cmd/wisp/panel_host_windows.go:304-305`）⇒ AC#4 那句"配置没递到宿主"又多一根凭据，**连默认值本身都对不上**。

**★ AC#4 的撞钉预检＋两格裁死（10-03 09:5x，只读腿 `255-c2` 交件＝`.scratch/wisp/probes/255/c2/census.md` 289 行／8 节全填／7 枚 commit 零 Go 命令；编排者复量四处，落账 `A561`）**
- **我现量坐实**：`cmd/wisp/panel_host_windows.go:304-305` 逐字 `Width:  420,`／`Height: 260,`，位置在 `bringUp`（`:230`）里 `webview2.NewWithOptions` 的实参块 ⇒ **建窗时用，不是显示时用**；全宿主无 `MoveWindow`／`SetWindowPos`／`Resize`，`Show`／`HotShow`／`Hide` 只走 `ShowWindow`，唯一"再读一次"的机会是销毁重建（`RequestDispose`）。`internal/panel` 一侧**零枚几何**（两把尺都空）。`internal/config/tiers.go:30 "llm":"hot"`／`:34 "panel":"hot"` 在。
- ⚠⚠ **AC#4 落地那一次必须同批改本票自己的尺**（这是第 64／70 条规矩现出的现场，⛔ 不许留给下一程踩）：`255-r2` 已把 `panel_host_windows.go:304` 写进机读判据——`cmd/wisp/config_readers_255.go:131/:134` 的判词 cite 该行＋token，`cmd/wisp/config_receipt_255_test.go:429` 要求台账含 `panel_host_windows.go:304`，`:389/:396` 那条正控**就是拿 `[panel] width = 641` 种的**（我现量，行号比腿报的 `:388` 漂一行）。⇒ 动 `:304` 这几处**必红**；落地腿要**同批重判 roster 行＋换新正控样本**（`[session]`／`[observe]` 可接手），⛔ 不许只改产码让 `255-r2` 变红，⛔ 不许删 `r2` 断言换绿。
- **我裁两格（⛔ 不摆 owner，低利害＋可逆＋零契约）**：**形＝取值闭包**（宿主每次建窗现读，dispose→show 重建即跟上新值），理由＝`tiers.go:34` 已把 `[panel]` 段登记为 `hot`，值快照形（＝重启才生效）与登记表矛盾；**回执文案随之必须诚实**＝"面板关窗再开即跟上新值"，⛔ 不许写成"拖动即变"（今天没有 resize 路，见上）。撤销口令**「255 AC#4 改取值快照」**。**`[panel] height`＝0 的语义**（schema 注释说"0=auto from content"而宿主侧无任何实现读过 "auto"）＝**未定义即停**那一类，我按最小诚实定：**0 ⇒ 沿用现常量 260，并把"由内容定高"具名登记为没做**（⛔ 谁也不许替它猜一个高度算法）；撤销口令**「255 height 那格重开」**。
- 一枚〔腿报，未复核〕的盘上不一致（不裁、只登记）：`cmd/wisp/panel_host_windows_test.go:1` 只带 `//go:build windows`，而 `:6` 注释自称属 `winlive` 档（"NO CI job"）；我现量 `ci.yml` 里 `tags`／`GOFLAGS` **零命中**⇒按标签读那枚真窗用例**会进 CI windows 档**。此事归已交件的 `winlive-census-1` 那一族，⛔ 本程不据此派改造腿。
- **可照抄的定式（腿量＋我复认形状）**：AC#4 不需要新造机制——最像的现成件＝热键桥（`internal/ball/hotkey_reload.go:11-13` 明写"消费包刻意不 import 配置，宿主供闭包"＋`type HotkeySource func() HotkeyConfig`），改动面最小＝变参 hook（`cmd/wisp/resident_ball_windows.go:63-81`，原文自述 hook 而非参数⇒不必改那 7 处测试构造点），同一条裁定还有审批门（`cmd/wisp/resident_windows.go:110-123` "BUILT HERE, by the assembly root"＋"Zero new package-level dependency edges"）。⇒ **"不开 `panel→config` 依赖边"这一格答＝不用开**：`PanelManager` 在 `cmd/wisp` 的 `package main` 里，而 `cmd/wisp` 本来就 import `internal/config`。⚠ 图级证明（`GOOS=windows go list -deps`）腿没跑（禁 Go 命令），**写腿自己补那一发**。
