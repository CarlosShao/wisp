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
- §3 装配根接缝与最小改动面
- §4 可照抄的注入先例
- §5 不依赖真窗口的机读判据候选
- §6 我可能写错的条目（自我对抗）
- §7 量不到的地方（具名）
- §8 交件判语
