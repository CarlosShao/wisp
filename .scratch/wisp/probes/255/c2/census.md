# 255-c2 只读普查：票 255 AC#4「面板宽度改了要真生效」到底要动哪几行

代号 `255-c2`。角色：**只读普查腿**。今天 2026-10-03。

## §0 起手锚

- 取数命令（同一发）：`date "+%Y-%m-%d %H:%M:%S%z"` ／ `git log -1 --format="%h %ci"` ／ `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`
- 本地时刻：`2026-10-03 09:26:05+0800`
- HEAD：`5f9ff9d4 2026-10-03 09:25:59 +0800`
- 脏文件计数（显式根 `cmd internal tools docs scripts .scratch`）：`362`
- ⚠ 与本腿相关的**在飞写面**（同一次 `git status --porcelain -- cmd/wisp internal/config internal/panel` 现量）：
  - `M cmd/wisp/config_reload.go`（他人未提交改动）
  - `?? cmd/wisp/config_readers_255.go`（他人新增、未跟踪）
  - `?? internal/config/tiers_app_255r2_test.go`（他人新增、未跟踪，`[app]` 键级补尺归 `255-r2`）
  ⇒ 这三枚的路径本腿**只读不写**。凡本腿引用到 `config_reload.go` / `config_readers_255.go` 的行号，都是**工作树现量**（可能随对方落盘再漂），会在具名处标注；`git show HEAD:` 对照过的会写明。

## §0.1 硬规自检

- ⛔ 本腿未跑任何 Go 命令（`go test`／`go build`／`go vet`／`go run` 全部零次）。需要跑才拿到的读数一律进 §7。
- 只读手段：`Read` / `grep` / `find` / `git log` / `git show` / `git diff` / `git status`。
- ⛔ `frontend/**`、`design/**` 未读、未引、未转述。
- ⛔ 冻结件一字未动（本腿全程只读）。
- ⛔ AC 框未碰、票面未改。

## 章节占位

## §1 那枚写死的尺寸现在在哪

### 1.1 `cmd/wisp` 里的逐枚名册（现量：工作树＝HEAD，该文件 `git diff --stat` 为空）

| # | file:line | 逐字 | 给谁用 | 时机 |
|---|---|---|---|---|
| 1 | `cmd/wisp/panel_host_windows.go:304` | `Width:  420,` | `webview2.WebViewOptions.WindowOptions.Width`（字面量在 `webview2.NewWithOptions(...)` 的实参里，调用起于 `:298`、整块 `:298-307`） | **创建窗口时**（`bringUp`，定义 `:230`），⛔ 不是显示时 |
| 2 | `cmd/wisp/panel_host_windows.go:305` | `Height: 260,` | 同上 `.Height` | 同上 |

这两枚就是全仓**唯一**的面板几何字面量。其余佐证：

- **结构体里没有第二处**：`PanelManager` 的字段名册（`cmd/wisp/panel_host_windows.go:137-167`）里 `dataPath string:141`／`assets *panel.Assets:142`／`disp *panel.ComposerDispatch:143`，**没有任何宽度/高度字段**；构造函数 `NewPanelManager(disp, assets, dataPath)`（`:175-177`）签名里也没有（票面现量第 2 条那句「全仓只有这一枚定义」本腿复核成立：`grep -rn "func NewPanelManager" cmd internal tools` 只命中 `:175`）。
- **显示时不重设几何**：`grep -n "MoveWindow\|SetWindowPos\|Resize" cmd/wisp/panel_host_windows.go` 只命中 `:45` 的**注释**（那句指的是依赖库 `pkg/edge` 的 `Chromium.Resize`，本腿只判断其射程＝不是本仓代码，未引用其内容作证据）。方法名册（同文件 `grep -n "func (m \*PanelManager)"`）里 `Show:395`／`HotShow:461`／`Hide:488` 都只走 `ShowWindow`（`:92` 那枚 proc）。⇒ **窗口一旦建成，420×260 之外没有任何一条码再往里写尺寸**；唯一"再读一次"的机会是销毁重建路径 `RequestDispose`（`cmd/wisp/panel_resident_windows.go:393-399`，注释 `:390-392` 明写 "the recreate path ... a later RequestShow builds a fresh window"）。这一条对 §5 的判据形状是承重的。
- **同名单但无关**（免得下一枚腿重复找）：`cmd/wisp/testdata/esclistener/main.go:219` 的 `40, 40, 420, 220` 是测试数据里另一枚 `CreateWindowEx` 窗口（x,y,w,h），⛔ 不是面板；`internal/ball/liquid.go:33` 的 `420` 是音频包络毫秒、`internal/ball/tokens.go:427` 与 `:436` 的 `260` 是动画毫秒，均与几何无关；`internal/observe/thresholds.go:11/:12/:103` 的 `420` 是句柄数读数——**该文件是冻结件，本腿只判断"这些数字不是像素"这一射程，未取内容作规矩**。

### 1.2 `internal/panel` 一侧有没有第二处硬编码：**没有，零枚**

- `grep -rn "Width\|Height" internal/panel --include=*.go`（含测试文件）⇒ **0 命中**。
- 放宽到 `grep -rni "width\|height\|geometry\|bounds" internal/panel --include=*.go` ⇒ 只有两类**子串假阳性**：`internal/panel/l2_grant_boundary_test.go` 的标识符 `inboundSeeds`（含 "in-bound"）与 `internal/panel/pump.go:431` 的一句散文里的 "it bounds"。前者属**冻结测试文件**，本腿只说明"这枚命中是字符串巧合、不是几何"这一射程判断，⛔ 未引用其内容。
- 结论：**面板这一族的几何字面量全部集中在 `cmd/wisp/panel_host_windows.go:304-305` 两行**，`internal/panel` 根本不知道"窗口尺寸"这件事存在。

### 1.3 缺省值口径不一致（票面 `A531` 那格的现量复核：成立）

- 配置侧缺省：`internal/config/schema.go:532` 逐字 `Width int \`toml:"width" default:"640"\`` ⇒ **640**。
- 宿主侧写死：**420**。⇒ 差 220 px，"配置缺省时的值"这句话今天**连缺省值自己都对不上**（票面 `A531` 那句）。
- 高度更糟：`internal/config/schema.go:534` `Height int \`toml:"height"\`` **没有 default tag**，`schema.go:533` 的注释写着 "0 = auto from content" ⇒ 配置缺省是 `0`，而宿主写死 `260`。⇒ AC#4 若把 `260` 直接换成"配置值"，一枚没写 `[panel]` 的默认配置会得出 **height=0 的窗口**，这是本腿能给 AC#4 写腿的最要紧预警（不是"做不到"，是要在装配根把 0 解释成 260 还是解释成 auto，二者都得具名选）。
## §2 配置那一侧有什么

### 2.1 `[panel]` 段字段名册（全量，五枚）

段挂在 `internal/config/schema.go:123` 逐字 `Panel   PanelSection   \`toml:"panel"\``。段体 `internal/config/schema.go:528-539`：

| 字段 | file:line | toml 键 | default tag | 缺省语义（取自同行注释） |
|---|---|---|---|---|
| `Enabled` | `internal/config/schema.go:530` | `enabled` | `"true"` | — |
| `Width` | `internal/config/schema.go:532` | `width` | `"640"` | `:531` 注释：面板宽度，px |
| `Height` | `internal/config/schema.go:534` | `height` | **无**（缺省即零值 0） | `:533` 注释：px；0 = auto from content |
| `KeepAliveInSession` | `internal/config/schema.go:536` | `keep_alive_in_session` | `"true"` | `:535`：会话期内保持 WebView2 温热 |
| `Scale` | `internal/config/schema.go:538` | `scale` | **无**（零值 0） | `:537`：UI 缩放；0 = 跟随系统 DPI |

⇒ 除宽高外段里还有三枚；**AC#4 的题面只圈了宽度**（票面第 20 行把 `420×260` 整体说成"配置缺省时的值"），所以 `height` 要不要一起递、按 §1.3 那个 0 的语义怎么裁，是写腿必须具名选的一格，本腿不代裁。

### 2.2 `TierRegistry` 给 `panel` 登记的档位

- `internal/config/tiers.go:34` 逐字 `"panel":   "hot",` —— **段级一枚，panel 没有键级行**。对照：`[app]` 走键级（`internal/config/tiers.go:47-50`，`app.theme`=hot、其余三枚=restart），`[voice]` 走键级（`:53-76`）。⇒ 档位词表里 `[panel]` 是**整段 hot**。
- 与 `plan()` 的同源守卫：`internal/config/manager.go:288` 逐字 `{"panel", &cur.Panel, &fresh.Panel, func() { cur.Panel = fresh.Panel }},` 在热应用段表里；`:294` 若 `TierRegistry[s.name] != "hot"` 就 `:299` panic。⇒ **"panel 是热档"这件事在 `internal/config` 内部是真的、且被一把会 panic 的尺钉着；缺的从来不是档位登记，是段外的读者。**

### 2.3 `cfg.Panel` 的非测试读者普查（回答"没人读"到底几行证据）

`grep -rn "\.Panel\b" cmd internal tools scripts --include=*.go | grep -v _test.go` 的**全部**命中归类：

| 命中 | 是不是 `config.Config.Panel` 的读者 |
|---|---|
| `internal/config/manager.go:288` | **是**，但在 `internal/config` 内部（热应用表本身），票面早就算过 |
| `cmd/wisp/config_readers_255.go:130` | **不是**：那是一枚**字符串**（在飞腿 `255-r2` 给 panel 写的 no-reader 判词），不是读值 |
| `cmd/balldebug/main.go:240`、`internal/ball/hotkey_reload.go:19/:102`、`internal/ball/hotkey_windows.go:98/:99/:455` | **不是**：这些是 `ball.HotkeyConfig.Panel`（热键字符串字段）与 panel 钩子，跟 `[panel]` 段无关 —— 同名陷阱，本仓有多枚同名文件/字段，这五处逐枚看过定义才敢归类 |

⇒ **`PanelSection` 的五枚字段在生产码全部零读者**：`grep -rn "Panel\.Width\|Panel\.Height\|Panel\.Enabled\|Panel\.Scale\|Panel\.KeepAlive" cmd internal tools scripts --include=*.go` 的非测试命中＝只剩上面那枚字符串。
- 加载期也不会响：`internal/config/unwired.go:32-37` 明写"unlocked 段（ball, **panel**, cost, voice, memory…）的零消费者键**故意不进** `unwiredKeys` 那张报错表"。⇒ 用户手写 `[panel] width = 900` 今天**既没读者也没报错**，与票面那句"这些段已立即生效"合起来就是误报的全部材料。

### 2.4 ⚠ `TierOf` 今天的生产调用者枚数（判据＝调用点枚数）

- 定义：`internal/config/tiers.go:93-96`。
- 尺：`grep -rn "TierOf(" cmd internal tools scripts --include=*.go`（**含测试**，本尺无 mock 面）。
- **HEAD `5f9ff9d4`：生产调用者 0 枚**。依据：`git show HEAD:cmd/wisp/config_reload.go` 里 `TierOf`／`TierRegistry` **零命中**（本腿实跑，grep rc 非 0＝无匹配），HEAD 上其余命中只有定义行与注释。
- **工作树（在飞）：生产调用者 1 枚**＝`cmd/wisp/config_readers_255.go:178`（同行注释自称 "TierOf's first production caller"）。**该文件未跟踪（`??`，属正在写 `cmd/wisp`＋`internal/config` 的 `255-r2`）**，所以这一枚**不属 HEAD、也随时可能被对方改写**。
- ⚠ 本腿实测到的**同一次会话内漂移**（记下来防后来人误读我的行号）：同一枚调用我在 `09:2x` 读到的是 `:174`，`09:3x` 再读已是 `:178` —— 对方在写这个文件。
- 另有**不经 `TierOf` 直读 map** 的产码点两枚：`internal/config/manager.go:294`（同源守卫）与 `cmd/wisp/config_readers_255.go:187`（在飞）；测试侧直读点在 `internal/config/tiers_255_test.go` 与 `cmd/wisp/config_receipt_255_test.go`（后者亦未跟踪）。
- ⛔ 本腿**没跑过任何** go 命令，也没打算用 `go build` 的 rc=0 当"接上了"的证明；上面每一枚都是 grep 计数。
## §3 装配根接缝与最小改动面

### 3.1 装配根是哪一段（具名函数＋行号）

面板宿主**只有一条生产构造路径**，在**常驻进程**那一族（不是 `wisp run`）：

| 层 | 函数 | file:line | 今天递给面板宿主什么 |
|---|---|---|---|
| 进程根 | `runResident()` | `cmd/wisp/resident_windows.go:30`；面板那一行 `:142` `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)` | **只有一个字符串** `rt.Layout.DataDir`（`rt` 是 `proc.Boot(env)` 的产物，`:39`），⛔ 没有任何配置对象 |
| 装配助手 | `newResidentPanelManager(dataDir string)` | `cmd/wisp/panel_resident_windows.go:183-205` | 组装三样：`disp`（`:189`）、`assets`（`:193` `panel.BuiltinAssets()`）、`dataPath`（`:203` `filepath.Join(dataDir, "panel-webview2")`），最后一行 `:204` `return NewPanelManager(disp, assets, dataPath), nil` |
| 入向链 | `newResidentComposerDispatch` → `newComposerDispatchChain` | `cmd/wisp/panel_resident_windows.go:166-171` → `cmd/wisp/panel_inbound.go:228` | **配置对象就在这儿被造出来又丢掉**：`cmd/wisp/panel_inbound.go:230` `mgr, err := config.NewManager(cfgPath, nil)`，它只喂给 `perm.New`（`:237-244`）与 `newConfigStore(mgr, ...)`（`:267`），**函数签名 `:228` 只返回 `(*panel.ComposerDispatch, error)`，mgr 出不去** |

⇒ **断点的精确形状**：装配根手里其实**已经有**这一进程唯一的那枚 `*config.Manager`（`cmd/wisp/panel_inbound.go:230` 造的），只是它被**关在 `newComposerDispatchChain` 的返回值之外**。AC#4 缺的不是"去读盘"，是"把已经在手上的那枚对象交出去"。
⇒ 常驻进程的球腿与审批腿**完全不读配置**（`grep -rn "config\.|internal/config" cmd/wisp/resident_windows.go cmd/wisp/resident_approval_windows.go cmd/wisp/resident_ball_windows.go` 零命中），所以"再造一枚 Manager"这条路等于**同一文件第二个加载者**，正是票里第 25 行、`cmd/wisp/panel_inbound.go:204-208` 写着的"one truth / ticket 101 的分裂状态"禁令。

### 3.2 ⛔ 会不会被迫新开 `internal/panel → internal/config` 那条边：**不会，一条都不必开**

- 决定性事实：**面板宿主根本不在 `internal/panel` 里**。`PanelManager` 与 `NewPanelManager` 都在 `cmd/wisp/panel_host_windows.go`，**`package main`**（同文件 `:3`），而 `cmd/wisp` 的产码**本来就 import `internal/config`**（`cmd/wisp/panel_inbound.go:73`、`cmd/wisp/panel_config_store.go:39`）。
- `internal/panel` 今天**不** import `internal/config`：`grep -rn "internal/config" internal/panel --include=*.go` 只命中两处**注释**（`internal/panel/config_handlers.go:30`、`:208`，本腿只判断其射程＝文字提及而非 import，未取内容作规矩）。
- ⇒ 只要那枚"宽度"落点是 **`PanelManager` 的字段（package main）**，AC#4 的改动**一行 import 都不新增**。
- ⚠ **唯一会踩禁区的走法**（点名防下一枚腿顺手做）：把几何塞进 `internal/panel` 的某个类型（例如给 `panel.ComposerDispatch`／`SettingsView`／某个 handler 加个 width 字段）——那条路要么新开被禁的边、要么撞 C17 契约面（票面禁区第 25 行"新增面板快照字段／新方法名＝先停手上报"）。**本腿量到的结论是：AC#4 不需要碰 C17，也不需要碰那条边。**

### 3.3 最小改动面（逐枚 file:line，按"改动点数"排序，两形都给）

**共同必动（无论哪形，4 点）**
1. `cmd/wisp/panel_host_windows.go:137-167` —— `PanelManager` 加落点字段（今天**没有**任何几何字段，见 §1.1）。
2. `cmd/wisp/panel_host_windows.go:175-177` —— `NewPanelManager` 的签名/构造体。
3. `cmd/wisp/panel_host_windows.go:302-306` —— `WindowOptions{Title/Width/Height}` 那三行里，`Width: 420` / `Height: 260` 改为"读宿主字段，字段为缺省时才落到 420/260"（票面 AC#4 原话：写死值变成"配置缺省时的值"而⛔ 不是唯一来源）。
4. `cmd/wisp/panel_resident_windows.go:204` —— 唯一的产码构造点，把值/闭包递进去。

**为了"值从哪来"必动（1–2 点）**
5. `cmd/wisp/panel_inbound.go:228` —— `newComposerDispatchChain` 的签名多返回一枚 `*config.Manager`（或返回已解析好的几何），并在 `:230-233` 之后把它带出去；连带它的两枚调用者 `cmd/wisp/panel_inbound.go:219` 与 `cmd/wisp/panel_resident_windows.go:170` 各改一行。**这是"不开第二条 import 边、又不再造第二个 Manager"两条约束逼出来的唯一走法。**

**形 A：值快照（装配时定型）**——就是上面 5 点。
- 若用**第 4 枚位置参数**：还得改 7 处测试构造点 `cmd/wisp/panel_host_windows_test.go:516`、`:689`、`cmd/wisp/panel_host_windows_live_test.go:39`、`cmd/wisp/panel_resident_windows_test.go:63`、`:449`、`:526`、`:527`（现量 `grep -rn "NewPanelManager(" cmd internal tools --include=*.go`）。
- 若用**变参 hook**（照抄 §4 的 `ballHostHook` 定式）：这 7 处**一行都不用改**。⇒ 差异纯粹在改动面大小。
- 生效语义：值在 `runResident` 起跑那刻定型；**改配置要重启才看得见**（因为 `bringUp` 用的是构造时的快照）。

**形 B：provider 闭包（用时才读）**——1–5 点同上，只是字段类型是 `func() (int,int)`，在 `bringUp`（`cmd/wisp/panel_host_windows.go:298` 那块字面量处）调用。
- 生效语义：多一条真生效路径——`RequestDispose`（`cmd/wisp/panel_resident_windows.go:393-399`）销毁后 `RequestShow` 重建窗口，会读到**热加载后**的新值。**不需要任何 resize 机制**（§1.1 已量：本宿主没有 `MoveWindow`/`SetWindowPos`/`Resize`，`Show` 只 `ShowWindow`）。
- ⛔ 本腿不裁 A／B：这一格直接决定票 255 AC#1 那句"已立即生效"对 `[panel]` 能不能说真话（见 §3.4），属编排者裁。

### 3.4 ⚠ 与在飞腿 `255-r2` 的**硬耦合**（本腿认为是最要紧的一条排程事实）

在飞的 `cmd/wisp` 写面已经把 `panel_host_windows.go:304` 那行**写进了机读判据**：
- `cmd/wisp/config_readers_255.go:130`（未跟踪，在飞）给 panel 的判词逐字含 `cmd/wisp/panel_host_windows.go:304 [Width:  420,] is hard-coded and NewPanelManager receives no config (票 255 AC#4 owns that break)`。
- 尺子＝`cmd/wisp/config_receipt_255_test.go:172` `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`：`:190` `os.ReadFile` 把 cite 指向的文件读回来、`:200-202` 要求**那一行仍然含那个 token**，否则红（"the evidence drifted, so re-adjudicate this row"）。
- 同一枚未跟踪文件里 `:406-412` 还要求热加载审计行含 `config: HOT-RELOAD-READER section=panel` + `no-reader` + `panel_host_windows.go:304`；`:369` `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence` 的整条正控就是拿 `[panel] width = 641` 种的。

⇒ **AC#4 一动 `:304`（不动就谈不上"变成缺省值"），这三处必红。** 这不是"断言被放宽"的问题，而是**判据的证据行moved ⇒ 那一行必须重判**：写腿必须在同一次改动里把 panel 的判词从 `hotClaimNoReader` 换形（形 A ⇒ `hotClaimSnapshotOnly`，与 `cmd/wisp/config_readers_255.go:99` 那枚 `[agent]` 行同族；形 B ⇒ 可争 `hotClaimConsumed`），并给一把新的正控替掉 `[panel]` 那枚"读者为零"的样本（⛔ 摘尺不许当修尺：票面第 17 行那句仍然管着 AC#1 那把尺，`[session]`／`[observe]` 这类真无读者段在 `cmd/wisp/config_readers_255.go:118`、`:123` 已具名，可以接手当样本）。
⇒ 排程后果：AC#4 **不能**在 `255-r2` 落盘之前动 `cmd/wisp/panel_host_windows.go:304`，否则两边互相把对方的测试改红；票面第 32-34 行那条起跑判据（`git status --porcelain cmd/wisp internal/config` 为空）对 AC#4 同样适用。

### 3.5 ⛔ 本腿**不**停手上报"非开那条边不可"

量到的事实是反的：那条边**根本不需要开**。本腿唯一想标为"未定义即停"的是 §1.3 那格——`[panel] height` 的配置缺省是 **0 且 schema 注释说 0＝auto from content**（`internal/config/schema.go:533-534`），与票面"写死的 260 变成配置缺省时的值"在字面上互相拉扯（0 到底翻成 260 还是翻成让内容定高？票面与裁定都没写）。这枚**留给编排者具名裁**，本腿不按自己的判断填。
## §4 仓里现成的同类定式（AC#4 能不能照抄、不新造机制）

**能照抄，而且至少有三枚现成形**。按"与 AC#4 的相似度"排：

### 定式 ①（最像）：消费侧⛔ 不 import 配置，装配根递一枚**取值闭包**进去 —— 热键桥

- 规矩原文：`internal/ball/hotkey_reload.go:11-13` —— "The bridge is deliberately free of any config import - internal/ball owns no business/config knowledge (SPEC-01 §3) - so **the host supplies two closures and one optional refresh hook**"（`:16-21` 就是那段用法示例）。
- 类型：`internal/ball/hotkey_reload.go:50` 逐字 `type HotkeySource func() HotkeyConfig`；构造器 `internal/ball/hotkey_reload.go:73` `func NewHotkeyReloader(binder HotkeyBinder, current HotkeyConfig, src HotkeySource) *HotkeyReloader`。
- 装配根那一头的实配点：`cmd/balldebug/main.go:231` `mgr, errC := config.NewManager(*configPath, nil)` → `:237-242` `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig { h := mgr.Config().Hotkey; ... })` → `:243` `bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }` → `:244` `mgr.OnReload = bridge.OnReload()`。
- **⇒ 这一枚的形状与 AC#4 的形 B 逐字同构**：`internal/ball` 不开 `→internal/config` 的边、自己绝不读盘；配置由宿主以闭包供给；消费包只看得到 `HotkeyConfig` 这种自己的类型。AC#4 只要把 `HotkeyConfig` 换成"面板几何"那一枚 package-main 里的小类型即可，⛔ 不需要新机制。
- ⚠ 同一枚定式在**出货进程里其实没接上**：全仓产码里 `NewHotkeyReloader` 只有 `cmd/balldebug/main.go:237` 一处调用，常驻球腿用的是编译期缺省（`cmd/wisp/resident_ball_windows.go:171` 逐字 `Hotkeys:  ball.DefaultHotkeys(),`）。这格已被在飞腿写进判词（`cmd/wisp/config_readers_255.go:113` 的 `hotClaimDebugHostOnly` 行）。⇒ **别把 balldebug 当"已经接好了"的先例来省掉自己那一步**，它的价值只在形状。

### 定式 ②（改动面最小）：装配根用**变参 hook**递能力，旧调用点一行不改

- `cmd/wisp/resident_ball_windows.go:63-68`：`// ballHostHook is how the assembly root hands the ball host one capability it does not know how to build ... the resident ball file receives a function value and stays ignorant of what is behind it`；类型 `:68` `type ballHostHook func(*panelHostHooks)`；载体 `:70-74` `panelHostHooks{showPanel func(via string) bool}`。
- **原文替本腿说话**：`cmd/wisp/resident_ball_windows.go:76-78` —— "It is a hook and not a parameter **so the three existing two-argument call sites of startResidentBall keep compiling untouched**."
- 实配点：`cmd/wisp/resident_windows.go:163-165` `startResidentBall(rt.Registry, ra.vetoByEsc, withPanelHost(func(via string) bool { return panel.RequestToggle(via) }))`；被注入方签名 `cmd/wisp/resident_ball_windows.go:158` `func startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc, hooks ...ballHostHook)`，收 hook 的循环 `:160-163`。
- **⇒ AC#4 若给 `NewPanelManager` 加第 4 枚位置参数，要连带改 7 处测试构造点（§3.3 列了）；照这一枚定式走变参 hook，那 7 处一行不动。** 这正是"不新造机制"最省的那条路。

### 定式 ③（同一条裁定）：审批门"根上造、往下递"，且明写**一条新依赖边都不开**

- `cmd/wisp/resident_windows.go:110-115`：the approval gate is **BUILT HERE, by the assembly root**, and handed down as a function value and a binding. The resident file does not compose a gate of its own ... 同段还引了裁定原文（`:114` "两枚正向依赖边一律不开、改注入"）。
- `cmd/wisp/resident_windows.go:117-122`："Zero new package-level dependency edges"，并把尺子写成 `go list -deps ./cmd/wisp` 改前改后对照（⛔ 本腿不跑，见 §7）。
- 实配点：`:123` `ra := newResidentApproval()` → `:174` `ra.bindBallHost(rb)` → `:163` 把 `ra.vetoByEsc` 作为函数值递给球腿。
- ⇒ **AC#4 的"落点＝装配根把配置递给面板宿主"不是新裁定**，`:125-128` 那段注释已经把面板宿主纳进同一个形："It is BUILT here, by the assembly root, and handed to the ball host as a function value - the same shape ticket 246's ruling (ledger A481) set"。

### 定式 ④（配置**值**→组件的两枚老样本，供对照语义）

- 快照形：`cmd/wisp/run.go:423` `cfg := mgr.Config()` → `:424` `rt.cfg = cfg`（`:420-422` 注释自称 "A snapshot copy ... stay frozen at what this boot read"），消费点 `cmd/wisp/run.go:763` `DefaultTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) * time.Millisecond`、`cmd/wisp/run.go:1014`、`:1016`。⇒ **这就是"值递到了、但热加载追不上"的那一族**，在飞腿给 `[agent]` 的判词正是它（`cmd/wisp/config_readers_255.go:99` `hotClaimSnapshotOnly`）。
- 活读形：`cmd/wisp/panel_config_store.go:88` `cfg := s.mgr.Config()`（在 `ReadSettings` 里，每次调用现读），配置对象持有在 `cmd/wisp/panel_config_store.go:66-71` `type configStore struct{ mgr *config.Manager ... }`，构造 `:75` `newConfigStore(mgr, secrets, auditf, actor)`，socket 侧只有接口 `internal/panel/config_handlers.go:232` `type ConfigStore interface`（三枚方法 `:234/:238/:242`）。⇒ **这是"面板这一族今天唯一真在读配置的那条路"，也是"接口声明在 `internal/panel`、配置对象留在 `package main`"这一分工的原样先例**——AC#4 要的正是同一分工，只是把"读设置"换成"读几何"。

### 定式 ⑤（顺手一枚：几何类选项的缺省/夹取模板）

- `internal/ball/ball_windows.go:64` `SizePx int // configured orb size 44..72 (0 = default 56)`，`internal/ball/ball_windows.go:144-151` 给了"0→缺省、越界夹住"的三步处理。⇒ §1.3 那枚"`[panel] height` 缺省是 0 还是 260"的纠结，**仓里有现成的表达模板**（一枚具名缺省常量＋显式夹取），但**选哪个语义仍要编排者裁**（§3.5），本腿只指出模板存在。
- 对照：常驻球腿根本没设这枚字段（`cmd/wisp/resident_ball_windows.go:165-185` 的 `ball.Options{}` 里无 `SizePx`）⇒ 面板之外的 `[ball]` 段犯的是同一枚病（在飞判词 `cmd/wisp/config_readers_255.go:117` 亦如此记）。
## §5 不依赖真窗口的机读判据：有没有、长什么样、落在哪

**直接答案：有，而且不止一枚。票面 AC#4 那句"改 width ⇒ 真窗口宽度随之变"确实只有本机可量，但"配置值真被面板宿主用上了"这件事可以拆成两枚机读断言：(i) 创建窗口那一步**收到的实参**取自配置；(ii) 装配根到宿主那一段真跑过一次。二者都不需要 WebView2 窗口。**

先例（本仓已有同族尺，⛔ 不是本腿新造的规矩）：
- `internal/ball/hotkey_reload.go:42-43` —— `HotkeyBinder` 的注释原文 "the seam exists so the diff logic is testable with no window"（类型体 `:44-46`）；假件实现在 `internal/ball/hotkey_status_test.go:444`（`func (r *recorder) RebindHotkeys(...)`），配合 `NewHotkeyReloader(..., func() HotkeyConfig { return src })` 用（`internal/ball/hotkey_status_test.go:458`、`:511`）。⇒ **"窗口那一侧不可测 ⇒ 把决定窗口参数的那一步做成无窗口可测的 seam"是本仓已经付过钱的走法。**

### 候选 A｜字段被赋成配置值（最便宜，但只钉一半）
- 形状：default tier 测试里 `NewPanelManager(disp, nil, dataPath, <planted>)` ⇒ 断言 `mgr` 的那枚新字段 == planted。
- 落点候选：`cmd/wisp/panel_host_windows_test.go`（已有同族构造点 `:516`、`:689`，⛔ 不必新建文件、不必进 winlive 档）。
- 钉得住：配置值进了宿主。**钉不住**：`bringUp` 到底有没有拿它去建窗（§1.1 已量：那是全宿主唯一一处用几何的地方）。⇒ **A 不许单独当 AC#4 的判据。**

### 候选 B｜把 `WindowOptions` 的组装抽成纯函数，断言返回值（推荐的主判据）
- 形状：`bringUp` 里 `cmd/wisp/panel_host_windows.go:298-307` 那块字面量改成一枚可单测的构造（`func (m *PanelManager) windowOptions() webview2.WebViewOptions` 或包级 `panelWindowOptions(geom)`），缺省值以**具名常量**留在同文件；测试不调 `bringUp`、不建窗，直接叫这枚函数，断言 `Width` 随注入值变、⛔ 注入为零时才等于那枚缺省常量。
- 落点候选：被测函数留 `cmd/wisp/panel_host_windows.go`；尺放 `cmd/wisp/panel_host_windows_test.go`。
- 为什么它比 A 强：它钉的是**建窗那一步的实参**，正是票面 AC#4 "写死的 420×260 变成配置缺省时的值"的可读投影。
- ⛔ 这条⛔ 不许顶掉真窗那一格：它证的是"参数从哪来"，不是"Win32 真把窗口开成多宽"（§1.1 那条注释 `panel_host_windows.go:50-52` 自己就写着 `WindowOptions` 与实际客户区尺寸的关系**尚未被证明**）。

### 候选 C｜把"造窗口"本身做成可注入函数值（最接近端到端，成本最高）
- 形状：`PanelManager` 加一枚默认＝`webview2.NewWithOptions` 的函数字段；测试注入 recorder 捕获实参并返回 nil/假件，`bringUp` 走到那一步即被断言"确实调了、且实参来自配置"。
- 先例同族：`Refresh` 那枚可注入函数字段（`cmd/balldebug/main.go:243` 装配）。
- ⚠ 具名风险：这一枚引入了"生产永远走默认值"的注入点 ⇒ 按 AGENTS §1.3 与票面禁区"⛔ 不许用 mock 代替真的来假报完成"，**C 必须与真窗那一格同发**，并配一把"默认值没被偷偷换掉"的尺（形如：产码构造路径里那枚字段仍为 nil 的断言）。本腿把 C 列为候选，不是列为推荐。

### 候选 D｜装配根那一段真跑一次（种一枚临时 config.toml，不建窗）
- 形状：测试直接叫 `newResidentPanelManager(dir)`（`cmd/wisp/panel_resident_windows.go:183`），事先在 `dir` 里写好 `[panel] width = <planted>`，断言返回的 manager 携带 planted。
- 这条真正测的是 §3.3 第 5 点那格改动（`newComposerDispatchChain` 把 `cmd/wisp/panel_inbound.go:230` 造出的那枚 `mgr` 交出来）：**只要它变红，就说明"配置→装配根→宿主"这条路又断了**，而它既不需要 WebView2 也不需要真窗口。
- 落点候选：`cmd/wisp/panel_resident_windows_test.go`（同包、default tier、已有 `NewPanelManager` 构造先例 `:63`、`:449`、`:526`、`:527`）。
- ⚠ 排程预警：`cmd/wisp` 的测试在这台机与 CI 上都要先把 sherpa DLL 摆进 PATH（`scripts/wisp-cli-tests.sh:7-16` 记的就是这个坑；`docs/reports/HANDOVER.md:143` 第②条又记了一次）。本腿⛔ 没跑过，只转述那两处的**射程**。

### 候选 E｜复用**已存在**的那把证据行尺，把"宽度来自配置"挂上去长期盯
- 现量：`cmd/wisp/config_receipt_255_test.go:172` `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 会把 roster 里每枚 `file.go:LINE [token]` cite 的**那一行读回来**要 token 还在（`:190`、`:200-202`）。
- 形状：AC#4 落地时，把 `cmd/wisp/config_readers_255.go:130` 那行 panel 判词的 cite 从"硬编码那行"改指**新的取值处**（宿主字段/`windowOptions` 那一步）。此后任何人把宽度改回写死，**这把已经在跑的尺自己会红**。
- ⇒ 这不是新机制，是**把新事实挂到旧尺上**；也正是 §3.4 那格冲突的解法（冲突不是靠删断言消，是靠重判那一行消）。

### 关于"做不到"这句
本腿**不**提交"AC#4 没有无窗口判据"这种否证。上面 A/B/D/E 四枚都可从盘上读出来、都不新建机制、都不放宽任何既有断言（B/D/E 只做加法；C 需要额外一把防 mock 的尺）。真窗那一格照票面留在"只有本机可量"族，⛔ 不伪装成 CI 测过。
- §6 我可能写错的条目（自我对抗）
- §7 量不到的地方（具名）
- §8 交件判语
