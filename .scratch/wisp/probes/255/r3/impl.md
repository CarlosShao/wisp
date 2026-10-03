# 票 255 AC#4 · 写腿 `255-r3` 交付件 — `[panel] width/height` 到达建窗实参块

状态：**进行中**（本文件是第一枚骨架 commit，§0／§1／§2 已写满并复核，§3－§8 在后续 commit 填实；
本枚**不是交件**，交件判据见文末「交件自查」。）

---

## 0. 起手锚

| 项 | 读数 |
|---|---|
| 起手时刻 | `2026-10-03T10:38:31+08:00`（`date -Iseconds` 自取） |
| 起手 HEAD | `67ab595d3ad107d3198252bb01f08a922dcfd548`（`git log -1 --format=%H` 自取） |
| 分支 | `dev` |
| 写面起手态 | `git status --porcelain -- cmd/wisp internal/panel internal/config` ＝ **空**（0 行）⇒ 无别人的脏件，我是这两面唯一写手 |
| 票面全名 | `.scratch/wisp/issues/255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md`（61 行，含末尾「★ AC#4 的撞钉预检＋两格裁死」节，已全文读过） |

同时刻在飞／已交的同票腿：`255-r2`（AC#1／AC#5，最后写盘 10:33／10:35，产码**已 commit**＝
`67ab595d` 与 `13dd60b4` 两枚之内）⇒ 它与我在同一写面**串行**，不并发。

### 基线现绿的用例名册（收尾要逐名 comm 对比、不许丢名）

跑法＝`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./cmd/wisp/ ./internal/panel/ ./internal/config/`，
原始输出已落 `.scratch/wisp/probes/255/r3/baseline-verbose.txt`（`tools/d22scan/runtests.sh` 那一族的 `-v` 形，
名字才出得来）。

| 读数 | 值 |
|---|---|
| `--- PASS` 枚数 | **386** |
| `--- FAIL` 枚数 | **5** |
| `--- SKIP` 枚数 | **0** |
| 包终态 | `FAIL cmd/wisp 190.506s` · `FAIL internal/panel 1.622s` · `ok internal/config 1.052s` |

**起手就红的 5 枚（全部不是我造成的，逐名登记＋具名原因）**：

1. `internal/panel` · `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`
2. `internal/panel` · `TestC21DesignTokensFourWayAgree`
   ⇒ 这 2 枚＝编排者预告的**已知常红**，起因＝`design/**` 在工作树里被删未 staged
   （我现跑 `git status --porcelain -- design` 复认：` D design/assets/base.css` 等 8+ 行在）。
   ⛔ 我不动 `design/**`、不动 `internal/panel/tokens_fourway_test.go` 一字。
3. `cmd/wisp` · `TestApprovalCardViewJSONKeysMatchFrontendTypes` — 现量红句逐字
   `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`（`approval_test.go:129`）。
4. `cmd/wisp` · `TestComposerContractTypesMatchFrontend` — 同族（前端契约面，属 `frontend/**` 两层禁令射程，我没读它）。
5. `cmd/wisp` · `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` — 真进程＋真答复流那一族，
   失败点是「卡片没点名 `risk.permission_mode`」。**不是 `[setup failed]`、不是 `0xc0000135`**：
   本次跑了 386 枚 PASS ⇒ 用例真在跑，PATH 那一格是对的。

⇒ 这三枚 `cmd/wisp` 红**在我落码之前就在**，我的判据是**名字集合只增不减、且不新增第 6 枚红**，
不是「整包绿」。收尾 comm 表在 §6。

---

## 1. 复跑编排者的现量（逐条，我自己跑的尺）

| # | 编排者的待验断言 | 我的复跑尺 | 结果 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_host_windows.go:304-305` 逐字 `Width:  420,`／`Height: 260,` | `grep -n 'Width:\|Height:' cmd/wisp/panel_host_windows.go` | **成立**：`:304 Width:  420,`／`:305 Height: 260,`（另 `:51` 注释里有一句 "the 420x260 size this host asks WindowOptions for"） |
| 2 | 位置在 `bringUp`（`:230`）里 `webview2.NewWithOptions` 的实参块 | `grep -n 'func (m \*PanelManager) bringUp'` ＝ **`:230`**；实参块起于 `:299` `w := webview2.NewWithOptions(webview2.WebViewOptions{` | **成立**，且 `:302 WindowOptions: webview2.WindowOptions{` ⇒ 建窗时用，不是显示时用 |
| 3 | 全宿主没有 `MoveWindow`／`SetWindowPos`／`Resize` | `grep -n 'MoveWindow\|SetWindowPos\|Resize\|SetBounds' cmd/wisp/*.go` | **成立（一处措辞要更正）**：零调用点；唯一命中是 `panel_host_windows.go:45` 的**注释**，它引用库侧的 `(*edge.Chromium).Resize()` 并明说那不构成阻塞。⇒ "没有 resize 路" 对，"全宿主没有 Resize 这个词" 不对 |
| 4 | `Show`／`HotShow`／`Hide` 只走 `ShowWindow` | `grep -n 'ShowWindow' cmd/wisp/*.go`：宿主侧 `:92` 取 proc、`:413` Show、`:500` Hide，别无他用 | **成立** |
| 5 | 唯一"再读一次"的机会＝销毁重建（`RequestDispose`） | `grep -rn 'RequestDispose' cmd internal tools` | **成立**，且它是 `residentPanel` 的方法不是 `PanelManager` 的：`cmd/wisp/panel_resident_windows.go:393`，体逐字 `return rp.post(func() { rp.mgr.Destroy() })`；`PanelManager.Destroy` 在 `panel_host_windows.go:529`，`:534` 把 `m.created = false` 复位 ⇒ 下一次 `Show` 走第二次 `bringUp` |
| 6 | `internal/panel` 一侧零枚几何 | 我在 §3 另跑（`internal/panel` 全库 `Width\|Height\|Geometry` 非测试命中） | 待 §3 落 |
| 7 | `internal/config/tiers.go:34` 把 `[panel]` 段登记为 `hot` | `grep -n '"panel"\|"llm"' internal/config/tiers.go` | **成立**：`:30 "llm": "hot"`／`:34 "panel": "hot"`，行号一字不差 |
| 8 | `NewPanelManager(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string)` 是全仓唯一定义、**无配置参数** | `grep -rn 'func NewPanelManager' cmd internal tools` | **成立且唯一定义在 `cmd/wisp/panel_host_windows.go:175`＝`package main`**，⛔ 不在 `internal/panel` ⇒ 编排者那一格「在后者就停手上报」**不触发**，加参数／加变参 hook 属合法装配根动作 |
| 9 | `cmd/wisp/config_readers_255.go:131/:134` 的判词 cite 了 `panel_host_windows.go:304` 那行＋token | 我 sed 复认 115-145 段 | **成立**：`:131` 注释 cite `(cmd/wisp/panel_host_windows.go:304 [Width:  420,])`，`:134` 判词逐字 `... cmd/wisp/panel_host_windows.go:304 [Width:  420,] is hard-coded and NewPanelManager receives no config ...` |
| 10 | `config_receipt_255_test.go:429` 要求台账含 `panel_host_windows.go:304` | sed 复认 420-435 | **成立**：`:429 if !strings.Contains(line, "panel_host_windows.go:304") {`；且 `:425` 还要求同一行含 `no-reader`（编排者没点名这半枚，**它同样必红**，见 §5） |
| 11 | `:389/:396` 那枚正控**就是拿 `[panel] width = 641` 种的** | sed 复认 380-400 | **成立**：`:388 func TestTicket255ReceiptOmitsPanelFromTheImmediateSentence`（编排者给的 `:389` 是函数体行，函数名在 `:388`）、`:389 const plantedWidth = 641`、`:396 r.plant(t, "[fs]", "[panel]\nwidth = 641\n\n[fs]")` |
| 12 | `[panel] width` 的 schema 默认值是 **640**，与宿主写死的 420 不一致 | 我在 §2 现读 schema.go | 待 §2 落（编排者引的是 `A531` 那句，行号 `schema.go:532` 未复核） |

---

## 2. 取值闭包落点（编排者裁死的形，我不再选边）

**裁形复述**：宿主**每次建窗现读**一次配置，`dispose → 再 show` 重建时跟上新值。
理由＝`tiers.go:34` 已登记 `[panel]` 为 `hot` ⇒ 值快照形（重启才生效）与登记表矛盾。
**`height == 0` ⇒ 沿用现常量 260**，"由内容定高"具名登记为**没做**。

### 落点（三处产码，全在 `cmd/wisp`，零新依赖边）

| # | 文件 | 动作 |
|---|---|---|
| P1 | `cmd/wisp/panel_host_windows.go` | `NewPanelManager` 增**变参 hook**（`opts ...panelHostOption`，照 `cmd/wisp/resident_ball_windows.go:73 withPanelHost` 的自述定式："It is a hook and not a parameter so the three existing two-argument call sites ... keep compiling untouched"）；新增 `windowOptions()` 方法**在每次 `bringUp` 现取一次几何**并返回整个 `webview2.WindowOptions` 块；`:304-305` 两行字面量撤出建窗实参块 |
| P2 | `cmd/wisp/panel_resident_windows.go:183` `newResidentPanelManager(dataDir)` | 装配根供闭包：`config.LoadFile(filepath.Join(dataDir, configFileName), nil)` **在每次建窗时**现读，取 `cfg.Panel.Width`／`cfg.Panel.Height`；读不到 ⇒ 回落 420／260 并**出声**（`slog.Warn`，不静默） |
| P3 | `cmd/wisp/panel_host_windows.go` 顶部常量 | `defaultPanelWidth = 420`／`defaultPanelHeight = 260`＝"配置缺省时的值"（票面 AC#4 原话），不再是唯一来源 |

### 拓扑事实（这条决定闭包为什么只能长成 LoadFile 形，编排者未提，具名补上）

**面板宿主只在 resident 进程里存在，而 resident 进程在建窗主机的时刻手里没有 config Manager。** 现量：

- `grep -rn 'config.NewManager' cmd internal | grep -v _test` ＝ 4 处：`cmd/balldebug/main.go:231`、
  `cmd/wisp/panel_inbound.go:230`、`cmd/wisp/run.go:409`（＋ `internal/ball/hotkey_reload.go:15` 注释）。
  `resident_windows.go` **一处都没有**（`grep -n 'config\.' cmd/wisp/resident_windows.go` 只命中 `:200` 的注释）。
- 建窗主机在 `cmd/wisp/resident_windows.go:142` `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)`，
  它**早于** `:208` `src := startResidentTaskSource(rt, ra)`，而 run.go 那枚 Manager（`:409`）在管线里、
  由 task source 触发才存在 ⇒ 把 Manager 当构造参数递给面板宿主**在时间上不可能**。
- 结论：闭包只能自己解析 config.toml，而 `config.LoadFile(path, nil)` 是本仓既有的"一次性现读"定式
  （`cmd/wisp/models.go:184` 就是这么用的）。⛔ 我**不**在 resident 腿新起一枚 Manager——
  那会把 `residentPanelHotReloadNote`（`panel_resident_windows.go:162` 逐字
  "panel host (resident): this leg does not tick config.toml either; "）那句判词弄假，属越界改别的腿的裁定。

⇒ 取到的形状＝**每次建窗现读盘**，比内存快照更新鲜；`PanelManager` 本身**不认识** config 包
（`panel_host_windows.go` 的 import 块里没有 `internal/config`，落码后我再用 `go list -deps` 出图级证明）。

### `height` 那一格的诚实登记

`internal/config/schema.go` 的 `[panel] height` 注释写了 "0=auto from content"（§8 里我自己 grep 过再落原文）。
宿主侧**没有任何实现读过 "auto"**。按裁死的最小诚实：`0 ⇒ 260`。"由内容定高"＝**没做**，登记在本节与 §7。
