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
- §2 配置侧名册＋`TierOf` 调用点枚数
- §3 装配根接缝与最小改动面
- §4 可照抄的注入先例
- §5 不依赖真窗口的机读判据候选
- §6 我可能写错的条目（自我对抗）
- §7 量不到的地方（具名）
- §8 交件判语
