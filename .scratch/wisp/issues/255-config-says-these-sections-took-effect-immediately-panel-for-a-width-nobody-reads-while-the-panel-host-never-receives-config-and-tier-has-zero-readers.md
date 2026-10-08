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

## 10-08 10:3x 编排者补格（来源＝只读普查腿 `180-c1`，件 `.scratch/wisp/probes/180/c1/`；AC 框一枚未动）

⚠ 这一格**归进已有的 `AC#1`**（那句"已立即生效只许说真话"），⛔ 不新开 AC 号——本票 `AC#4` 也已经把"`[panel] width` 真生效"登记在案，`180-c1` 独立量到同一件事 ⇒ **票 180 的 AC#2/AC#3 从此归口本票，不在 180 里重复施工**（记我：180 立票在前、255 落地在后，我没查重到这一步，第 123 条同族）。

**现量到的过强宣称**（`180-c1` 具名上报，本程**一字未改**，它在别人写面上）：`cmd/wisp/config_readers_255.go:161` 那句写着 `[panel] width` 属"每次 `wisp run` 重读、关窗再开即跟上新值"这一档。三把尺对拉后：
1. `cmd/wisp/panel_host_windows.go` 的 `Show` 只在 `!created` 时才走 `bringUp`（`:499-509`，现量 `created := m.created` 与 `if !created {` 同段）；
2. C27 定的是**隐藏而非销毁**，而 `RequestDispose` 的**非测试调用者＝0 枚** ⇒ 面板一旦建过，关窗再开**不会**重建，尺寸也就不会重读；
3. ⇒ 真实生效条件＝**重启常驻进程**，不是"关窗再开"。

**给下一位碰这一格的判据（⛔ 不是本程的活）**：把 `:161` 那句改成说实话（"重启进程后生效"），或者把"重建"那一跳真接上并配一发**会响**的仪器（两种都行，选哪种归 `AC#4` 的落地腿一起裁——那枚腿本来就在处理宽度生效）。⚠ 无论哪种，⛔ 不许只把文案里的时间口径改宽来交差（那正是 `AC#1` 存在的理由）。

## 10-08 10:4x 机主当场拍板（原话入台账 `A698`）：`AC#4` 的口径从"重启后生效"**升为"不用重启就见效"**

问的是只有他能答的那件事：**"改面板宽度，要不要做到改完立刻见效、不用重启？"** 他回：**"要"**。
⇒ 本票 `AC#4` 从此不是"把话说实话"那一格（那是 `AC#1` 的活），而是**一枚功能**：已存在的面板窗口要在配置变更后拿到新尺寸。
⛔ **边界（我拍的，写死给下一枚腿）**：他批的是**结果**，⛔ 不是"改 C27（关窗＝隐藏不是销毁）"——**不许把这条读成契约变更的许可**。因此"关窗时销毁窗口"那形**默认不选**，除非前置普查证明保留隐藏语义做不到、且我把两形代价摆完再单独要他那一句。
前置普查＝只读腿 `255-a2`（零 go 命令，落点表六问：句柄与库侧出口／像素换算与几何源／STA 线程归属／热加载挂不挂得上／无窗仪器能钉哪一格、哪一半只能〔仅本机可量〕／雷区与契约）。⇒ **写码腿按在它之后**（此刻 `cmd/wisp` 的测试面与哈希对拉归验收腿 `33-v4`，并行写入会造出假"被篡改"）。

## 10-08 只读腿 `255-a2` 交件登记（六问落点表，⛔ 本段不改写上面任何一行、不勾任何框）

件＝`.scratch/wisp/probes/255/a2/10-landing-table.md`（起手 HEAD 现量 `5738661`，零 go build/vet/test，零产码写点）。
**只登记读数，不裁形、不选 hint、不排程**（编排者的活）。六问各一句：

1. **句柄与库侧出口**：HWND 在 `cmd/wisp/panel_host_windows.go:398` 取（`windows.HWND(w.Window())`）、`:415` 落 `m.hwnd`、`:637` 清零，读口 `:290 windowHandle()`。**库侧有现成出口＝有**：`go-webview2` 的 Go 包装层接口 `common.go:54 SetSize(w, h int, hint Hint)`，实现 `webview.go:405` 内部就是 `AdjustWindowRect`＋`SetWindowPos`（`:427-430`）并尾随 `browser.Resize()`（`:431`）；同接口 `common.go:39 Dispatch(f func())` 是现成的回线程通道。⚠ **两条路数值语义不等**：建窗 `webview.go:296-320` 把 `WindowOptions.Width` 当**外框**用，`SetSize` 把入参当**客户区**（过 `AdjustWindowRect`）⇒ 同一个 420 不等同一件事。
2. **像素换算**：`panelWidthPx/panelHeightPx`＝`panel_host_windows.go:89-90`，唯一消费点 `:244`，注入口 `:199 withGeometrySource`，生产装配点 `panel_resident_windows.go:253`，几何源 `:200-209`（`config.LoadFile` 现读，失败回 `0,0`），schema `internal/config/schema.go:539-549`。⚠ **顶回派单转述一处：DPI 逻辑→物理那层换算今天根本不存在**（两把尺：非测试码里 DPI 族符号唯一命中是 `schema.go:547` 那行注释；`PanelSection.Scale` 非测试读者 0 枚），现量口径只在 `cmd/wisp/panel_geometry_255_winlive_test.go:169-176` 由真窗**反推** scale。
3. **线程归属**：改尺寸必须落那条专用 STA 线程（锁点唯一在 `panel_resident_windows.go:262`，泵 `:299 w.Run()`）。**现成机制＝`residentPanel.post` `:347`**（已发布走 `:358 w.Dispatch`、未发布走 `:362 rp.tasks`），既有同形用户 `:450`／`:482` ⇒ ⛔ 不需要新造通道。**跨线程调 `SetWindowPos` 的后果＝仓里没有凭据**（5 处 md 命中全是"宿主没有 SetWindowPos"那一类断言），本腿未拿常识冒充实测。
4. **热加载挂不挂得上＝挂不上现成的**：`OnReload`（`internal/config/manager.go:54`）只在 `:198 rep.Reload` 非空时开火，而 `[panel]` 登记在 **hot**（`internal/config/tiers.go:34`）⇒ 天生不开火；且 `startConfigReload` 的唯一非测试调用者是 `cmd/wisp/run.go:813`，**常驻腿 0 命中**。最少落点表给到 3 枚（宿主方法／`rp.post` 请求／同批文案），**第 4 处"谁去按"才是难点**，两枚候选都在别人写面上（`config_reload.go:105/:153` 与 `resident_windows.go:152-158`/`:205` 的 `hotReload258`）。⚠ **撞钉预检**：`cmd/wisp/config_readers_255.go:161` 今天逐字写「面板关窗再开即跟上新值」并 cite 已漂的 `panel_host_windows.go:262/:392`，`config_readers_255.go:17/:141` 与 `config_receipt_255_test.go:440` 还在 cite `:304` ⇒ **落地必须同批改这三处**，否则 AC#1 那格被踩红。
5. **无窗仪器钉得住**：`cmd/wisp/panel_pageover_33r10_windows_test.go:130` 那枚假控件**已经实现了 `SetSize(w,h,hint)` 的空体**、`:123` 已实现 `Dispatch`，注入点 `:154 m.w = c` ⇒ 改成记录型 sink 就能无窗问"发没发、发的是不是几何源此刻那一对数、hint 是哪一枚"，能力尺写法照该文件头 `:19-33`（问能力、⛔ 不问调用枚数）；`panel_geometry_255_test.go:227/:239/:266` 那枚 AST 钉（恰好 1 枚 `NewWithOptions`、`WindowOptions` 必须是 `windowOptions()` 调用）会自己拦住"顺手加第二条建窗路"。**只能〔仅本机可量〕**＝真窗改完的实际宽度（同类先例仅 `panel_geometry_255_winlive_test.go:48/:57/:173` 一枚）、DPI factor 真值、`AdjustWindowRect` 边框差几像素。winlive 家族现量 **12 枚文件**、默认 `go test ./cmd/wisp/` **不编译**（该文件头 `:9-11` 自述），⛔ 不许读成 CI 有载体。
6. **雷区与契约**：入站白名单六枚无几何（`internal/panel/bridge.go:42-45`＋`:66-67`）⇒ **这条链上面板侧零写、零批准**；若日后做"页面拖滑块改宽"＝新增入站方法名＝**C17 契约面**（本票禁区 `:25`），本腿只指认。**甲形（活窗直接改尺寸，最少 3 枚落点）不碰 C27**——不销毁、不新建、隐藏语义一字不动，`PLAN.md:1377/:2352/:3113` 的"单窗口复用绕开 Environment"原样成立；其风险是 `SetSize` 客户区语义、hint 选型（`HintFixed` 会把窗改成拖不动，库 `webview.go:409`）、跨线程须走 `rp.post`。**乙形（保留隐藏语义、Show 时重建）2 枚落点但要复活 `RequestDispose`（现量非测试调用者 0 枚）**，代价＝每次改宽后第一次 Show 从热 ≤200ms 掉回冷 ≤1500ms（烧 D32 那两行，阈值本体 `internal/observe/thresholds.go` 一字节不许动）＋每窗新建一枚 WebView2 Environment（`PLAN.md:1568` P11 整条理由）＋把 `bringUp:381-384` 那条 33-r7/r9 已量到危险的"同线程 dispose→recreate"从零调用者变成每次走。**⇒ "不动 C27 也能做到"的那条路就是甲形本身，乙形不比它省契约、反而更贵。**
