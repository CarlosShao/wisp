# 180-c1 格 1 · 票面 AC#1 前提的重验判语

起手锚 `601c2b18`；本件读数取于 **`70b00885`**（本程跑尺期间工作树被别的在飞腿推进过，见 §5）。
只读：产码零字节改动。

---

## 1. 三把尺的现跑读数

### 尺 (P3) —— 票面 AC#1 自己那把窗口尺，**原样复跑**

```
grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test
=> 空输出，rc=1
```

**票面那句"＝0 命中"今天仍然一字不差地成立。** 但这把尺证不了票面以为它证的东西 —— 见 §3。

### 尺 (P2) —— `internal/panel/pump.go` 顶部那句原文（整行读，未截断）

```
grep -n "" internal/panel/pump.go | sed -n '1,30p'   rc=0
```

逐字（行号＝pump.go 行号）：

```
15: // WHAT THIS FILE IS NOT: the transport. There is no Go -> page channel in this
16: // tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer,
17: // no local HTTP/SSE/websocket server, and assets.go:8 quotes D29 saying there
```

⇒ 票面引用的 `"no WebView2 host"` **确在 `pump.go:16`**，`15-17` 的格位也对。
⇒ 同一行还挂着 **`(ticket 33 is unclaimed)`** —— 这句**也过期**：`.scratch/wisp/issues/33-panel-host-c27.md` 在场且正在被
`33-r10` 那枚腿续做（本程现量 HEAD 就是它的 commit `70b00885`）；票面第 75 行编排者自己已记过这一笔过期，
并明确"现在不动那枚文件（它在别人写面上）"。**本程照办，不动。**

### 尺 (P1) —— 面板宿主/真窗的存在性

```
grep -n "webview" go.mod
=> 19:	github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect        rc=0
```

```
grep -rn "NewPanelManager(" --include=*.go cmd/wisp/*.go | grep -v _test.go
=> cmd/wisp/panel_host_windows.go:213:func NewPanelManager(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string, opts ...panelHostOption) *PanelManager {
=> cmd/wisp/panel_resident_windows.go:253:	return NewPanelManager(disp, assets, dataPath, withGeometrySource(panelGeometrySource(dataDir))), nil
   rc=0     ⇒ 定义 1 枚 ＋ 非测试调用点 **恰 1 枚**
```

```
grep -n "windowOptions()\|WindowOptions: m.windowOptions()" cmd/wisp/panel_host_windows.go
=> 236:func (m *PanelManager) windowOptions() webview2.WindowOptions {
=> 260:	return webview2.WindowOptions{
=> 262:		Width:  uint(width),
=> 263:		Height: uint(height),
=> 392:		WindowOptions: m.windowOptions(),        rc=0
```

---

## 2. 判语（两格，按派单要求分开）

### 格 ①「面板窗口今天到底存不存在」—— **存在。票面 AC#1 那句「这棵树里根本没有面板窗口」已过期。**

不是"部分过期"，是**整句作废**：`go.mod:19` 有 `github.com/jchv/go-webview2`；
`cmd/wisp/panel_host_windows.go` 定义 `NewPanelManager`（`:213`）并在 `windowOptions()`（`:236`）里
返回 `webview2.WindowOptions{Title, Width, Height}`（`:260-264`），该 options 在建窗那一行被交出去（`:392`）。
本程**没有开真窗**（硬约束 4），存在性全部出自文本尺与码上调用链。

票面那一格**勾错了前提**，但它勾的那三把尺**各自都还能复现**（P3 空输出、P2 引文准确）。
失效的不是读数，是**从读数到"没有窗口"这一步推理** —— 见 §3。

### 格 ②「存在的话它由哪条腿装配、什么时候才真开」

装配链（逐跳 file:line，全部非测试码）：

| 跳 | 位置 | 内容 |
|---|---|---|
| 1 | `cmd/wisp/resident_windows.go:151` | `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)` —— 常驻腿启动时装配 |
| 2 | `cmd/wisp/panel_resident_windows.go:228` | `func newResidentPanelManager(dataDir string) (*PanelManager, error)` |
| 3 | `cmd/wisp/panel_resident_windows.go:253` | `NewPanelManager(disp, assets, dataPath, withGeometrySource(panelGeometrySource(dataDir)))` |
| 4 | `cmd/wisp/panel_resident_windows.go:198-209` | `panelGeometrySource` 每次建窗 `config.LoadFile` 现读，`:207` `return cfg.Panel.Width, cfg.Panel.Height` |
| 5 | `cmd/wisp/panel_host_windows.go:199-200` | `withGeometrySource(src func() (width, height int))` → `m.geometry = src` |
| 6 | `cmd/wisp/panel_host_windows.go:236-264` | `windowOptions()`：`gw > 0` 才覆盖常量 `panelWidthPx = 420`（`:89`），`return webview2.WindowOptions{Width: uint(width)}` |
| 7 | `cmd/wisp/panel_host_windows.go:392` | 建窗调用点 |

**什么时候才真开**：**不是开机即开**。常驻腿的 panel 线程在 `loop` 的 select 上等 show 请求，
`show(via)` → `rp.post(...)` → `showOnThread`（`panel_resident_windows.go:408`）→ `rp.mgr.Show(...)`（`:409`），
成功才打 `"wisp: panel window is up (...)"`（`:411`）。也就是**按需拉起**（热键召唤／派发路径），
且 `bringUp` 可能拒绝（`:386-407` 的注释记着 33-r9 那一格：拒绝会把原因存进 `startUp` 并让线程退役）。

**并且：只有常驻 `wisp` 进程开。`wisp run` 那一侧今天不建 panel 宿主**
（`config_readers_255.go:151-153` 自己写着 NewPanelManager 的唯一非测试调用点就是 `:253`，
并把 `[panel]` 那一行判成 `other-process`）。

---

## 3. 票面那把尺为什么会判反（这是本程的头一格，写透）

`P3` 的射程是 **`internal/panel/`**。而面板**宿主**（建窗、取尺寸、调 WebView2）在 **`cmd/wisp/`**：
`internal/panel/` 装的是 snapshot 组装与 pump（数据半），从来就不含窗口代码。
⇒ 在 `internal/panel/` 里搜 `NewWindow|SetBounds|Rect{|Width:` 得 0，
**对"这棵树有没有面板窗口"这个问题没有回答能力**——它只回答"`internal/panel/` 这个包里有没有建窗码"。
`0 命中` ＋ `pump.go:16 那句注释` 两条各自都真，拼起来推"全树无面板窗口"是**射程错配**，不是读数错。

⚠ 同一处 `pump.go:16` 那句 "no WebView2 host" 也是**同一射程错配的产物**：
它写的是"这棵树没有 Go→page 的通道"，被票面读成了"这棵树没有 WebView2 宿主"。
前者到今天仍成立（`config_readers_255.go` 与 `pump.go` 都没声称有反向通道），后者已经不成立
（`panel_host_windows.go:392` 就是宿主在建窗）。**这两件事被同一行注释混在一起说了。**

---

## 4. 对 AC#2 / AC#3 方向的影响（头一格的后果，⛔ 本程不落地）

派单担心的是"AC#1 判错 ⇒ AC#2/AC#3 修法方向跟着错"。现量结果：

- **宽度那一跳今天已经接上了**，接它的是 **票 255 AC#4**（不是本票）。
  所以 AC#2「把那一跳接上」在 `[panel] width` 这一枚上**已经是完成态**，
  再按票面"缺整条链"去排一次队 = 重复施工。
- 票面 AC#3「不许用改默认值交差」**仍然有效但对象变了**：`default:"640"` 一字节没被改过，
  而被改的是**读者**，正是 AC#3 想要的那一支（接效果、不改数）。
- **仍然没接上的两枚**（同段、同一族，本程现量）：
  `PanelSection.Enabled`（`schema.go:540`，`default:"true"`）与
  `PanelSection.KeepAliveInSession`（`schema.go:546`，`default:"true"`）—— 见 `census.tsv` 与本件 §6。
- **仍然缺的不是"跳"而是"档"**：已建好的窗口不会跟上新值（没有 resize 路，`config_readers_255.go:155-160` 自己钉着；
  本程复跑 `grep -rn "MoveWindow\|SetWindowPos\|Resize(\|SetBounds" --include=*.go cmd/wisp/ internal/panel/ | grep -v _test.go`
  ⇒ 非注释命中只有 `cmd/wisp/testdata/esclistener/main.go`（testdata，非产码），`rc=0` 但**产码零命中**）；
  效果是"关窗再开"，不是"拖窗即变"。

---

## 5. 与本程派单的冲突（具名报回）

1. **派单说**：「`cmd/wisp/resident_windows.go` 里 `NewPanelManager` 已有非测试调用者」。
   **现量**：`resident_windows.go` 里没有 `NewPanelManager(` 的调用；该文件 `:141` 是**一句提到这个名字的注释**，
   `:151` 调的是包装函数 `newResidentPanelManager`。真正的 `NewPanelManager` 非测试调用点是
   **`cmd/wisp/panel_resident_windows.go:253`**（文件差一个 `panel_` 前缀）。
   ⇒ 派单的**方向对、文件名错**；本程一律以现量的 `panel_resident_windows.go:253` 为准。
2. **派单给的入口尺 70 与票内两程的 69 冲突**：两边都对，量的是两回事。
   `grep -c 'default:' schema.go` ＝ **70 行**；其中 `schema.go:8` 是**文档注释**在引用 `default:"..."` 这个写法本身
   （`// D36 rule 3: defaults live ONLY in the `default:"..."` struct tags. The TOML`），不是字段标签。
   ⇒ **真实带 `default:` 标签的字段＝69 枚**，去重字段名 62 枚。本程 TSV 的母体是 **69**。
3. **HEAD 漂移**：派单说起手 HEAD 不早于 `601c2b18`。起手＝`601c2b18` 整；
   本程跑到中途工作树被 `93147902`（编排者停车点）/`0ec19368`（`182-c1` 锚件）/`70b00885`（`33-r10` 测试腿）推进。
   `70b00885` 动过 `cmd/wisp/panel_host_windows.go`（+41/-11），**本程引用的该文件行号（`:199/:213/:236/:260-264/:392`）
   全部按 `70b00885` 复核过**；`internal/config/schema.go` 在 `601c2b18..70b00885` 区间 diff 为 **0 行**
   （`git diff 601c2b18..HEAD -- internal/config/schema.go | wc -l` ⇒ `0`，`rc=0`），所以票面 `schema.go:527-528` 那两行仍立（现量 `:541-542`）。

---

## 6. `[panel]` 五枚的逐枚现量（票面"各 0"那一说的复验）

尺：`grep -rn "\.<Field>\b" --include=*.go internal/ cmd/ | grep -v _test.go`（票面原尺）

| 字段 | schema.go | default | 非测试命中 | 归属 |
|---|---|---|---|---|
| `Width` | `:542` | `"640"` | **2** | 1 枚真读者 `panel_resident_windows.go:207`；1 枚是 `config_readers_255.go:161` 的**字符串引文**（不是读者）⇒ **真读者 1** |
| `Height` | `:544` | **无 `default:` 标签** | ≥1 | 同一行 `:207` 一起返回 ⇒ 有读者；但**不在本程 69 枚母体内**（无 default 标签） |
| `Enabled` | `:540` | `"true"` | 21（全部同名撞车） | ⛔ **无一枚属 `PanelSection.Enabled`**：`grep -rn "Panel\.Enabled\|panel\.Enabled" --include=*.go internal/ cmd/ \| grep -v _test.go` ⇒ **rc=1 空输出** ⇒ **读者 0** |
| `KeepAliveInSession` | `:546` | `"true"` | **0** | 全树（含测试、含 `.scratch` 外的产码）只有 `schema.go:545-546` 的定义与注释 ⇒ **读者 0** |
| `Scale` | — | **无 `default:` 标签** | 1（`internal/agent/budgets.go:94` 的 `b.Scale`，**另一枚字段**） | 不在 69 枚母体内；那一枚命中与 `[panel] scale` **无关**（撞名） |

⇒ 票面「那五枚 `[panel]` 字段的生产读取方各 **0**」**今天只对 3 枚成立**（`Enabled` / `KeepAliveInSession` / — ），
对 `Width`、`Height` **已过期**。
