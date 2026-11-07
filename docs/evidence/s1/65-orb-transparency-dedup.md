# 65 — orb 真透明 ＋ 单球架构（transparency/dedup，2026-11-07）

**腿名**：wave65-transparency-dedup（写腿，只 commit 不 push）
**机主凭据**：编排者 EnumWindows+GetPixel 真机实测两枚缺陷（本腿不复测票面、只修）。

## 一、修前缺陷（编排者实测，本腿照单全收）

1. **orb 窗黑底不透明**：「Wisp Orb」132x132，exStyle `0x80800A8`，四角读数
   (30,24,29)/(38,33,40)/(26,21,32)/(38,31,38)，对照桌面 (51,51,50) ⇒ 黑方块盖桌面。
2. **球画了两遍**：orb 窗画一枚完整液态玻璃球，球窗又画一枚 D2D 玻璃球，两窗并排。

## 二、根因链（读码＋本机实测，逐层）

**A（透明）四层排查结论**：

| 层 | 结论 |
|---|---|
| 页面层 | 页面本身正确：`alphaMode:"premultiplied"`、clear `a:0`、injected bootstrap 清 body。**但 bootstrap 缝是竞态**（见下） |
| WebView2 层 | `PutDefaultBackgroundColor(a=0)`（controller2，IID c979903e…）确已在调；但 jchv 包 `AddScriptToExecuteOnDocumentCreated` 传 **NULL handler**（chromium.go:130），注册可落在首文档已加载之后 ⇒ bootstrap 整段丢失（真机实测：页面活、WebGPU 跑、四角整进程生命周期显示上游 body 底色）；另 API put 有文档明载的生效滞后 known issue ⇒ 部分实例整个生命周期白底 |
| 窗口层 | **首因**：orb 窗挂 `WS_EX_LAYERED` 却从不调 `SetLayeredWindowAttributes`/`UpdateLayeredWindow` ⇒ 顶层重定向表面是未初始化的黑层，页面 alpha-0 像素全合成到黑上 |
| DWM 层 | 修法走主流 Tauri/tao+wry 配方：`DwmEnableBlurBehindWindow` 空区域（`CreateRectRgn(0,0,-1,-1)`＋`DWM_BB_ENABLE|DWM_BB_BLURREGION`）——tao window.rs `on_create` 的原样做法，让重定向表面 per-pixel alpha 感知、无模糊；不需要 `DwmExtendFrameIntoClientArea`/`SYSTEMBACKDROP_TYPE`（同一 alpha 感知要求该调用已满足） |

**B（单球）选型＝甲**（orb＝球的可见面，球窗＝纯交互层），理由（铁律 ⓪⓪⓪）：
- 主流 overlay 形状就是「一个可见层＋一个交互层」；orb 页面就是票 65 要的玻璃质感本体。
- **乙排除**：把上游页裁成光晕＝毁掉「集成现成开源件」的票值核心。
- **丙排除**：把交互语义搬进 orb 窗要动冻结 hit-test/drag/dock 契约，且 D2D fallback 反正得把球窗立回来——复杂买不到价值。

## 三、修法（改面与机制）

1. `orb_windows.go`：去 `WS_EX_LAYERED`（保留 TRANSPARENT|TOPMOST|TOOLWINDOW|NOACTIVATE）；
   建窗后 `enableDWMAlpha`（blur-behind 空区域）；`WEBVIEW2_DEFAULT_BACKGROUND_COLOR=0x00FFFFFF`
   环境变量（文档明载的 API 滞后解法，进程出生即透明；面板 WebView2 同进程共用 env，其页面
   自带不透明 CSS 底，无暴露面）；`newOrbWindow` 边长改 `WindowEdgePx(sizePx,dpi)`＝与球窗同 rect；
   `follow()` 改中心对齐＋工作区夹紧（`orbOverlayXY`），**z 序＝球上、orb 正下**
   （`hWndInsertAfter=ballHWND`）；wndproc 加 `WM_NCHITTEST→HTTRANSPARENT` 兜底（穿透不再依赖
   LAYERED 的系统解释）。
   **为什么 orb 在球下面**：实测非 layered 的 `WS_EX_TRANSPARENT` 窗**不被 WindowFromPoint 跳过**
   （物理点击后球窗 LBUTTONUP 触发 followOrb 把 orb 叠上球心，WindowFromPoint(球心) 即答 orb——
   interaction_live_test.go:343 两连红锁定的机制）。反向后指针永远先落球窗，冻结语义零改。
2. `orbasset.go`：served 副本加注入 `<style>html,body{background:transparent !important}</style>`
   （与 CSP meta 同缝、同为内存变换，asset 字节零动，sha 钉测照旧）——声明式，无脚本竞态。
3. `renderer_windows.go`：`overlayOnly` 模式——orb 模式下球窗只画 badge/进度环/queueDot/文本
   （页面无这些 API）＋ **alpha=1/255 的 hit-test 轮廓**（`OrbBody.WithAlpha(1/255)`，byte 恰 1：
   OS 对分层窗按像素 alpha 决定跳不跳，全零＝不可点＝上一腿实测的回归形状）。
4. `ball_windows.go`：renderFrame/repaintAnimated 按 `b.orb.Load()` 切 overlayOnly；拖拽中
   `followOrb()` 逐帧跟随（原拖完才跟，会拖出「第二球」错觉）。
5. **写面偏离（如实）**：`orbasset.go` 与 `renderer_windows.go` 不在派单写面清单内，但
   ①CSS 注入是编排者问题 A 第 4 条点名之路，落点只能在 orbasset.go；②overlay 绘制分支落
   renderer 内聚处，避免在 ball_windows.go 复制 40 行绘制代码造漂移隐患。均在
   `internal/ball` 独占区内。

## 四、验收读数（全部本机实跑）

**门禁（最终轮）**：

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/balldebug` | rc0 |
| `go build ./cmd/wisp` | **BLOCKED——非本腿**：并发票 214/323 腿在飞 8 个文件（`panel_attachment_214.go` 等），盘上中间态不可编译（`l.rows` field/method 冲突等）；internal/ball 自身构建干净 |
| `go test ./internal/ball -count=1` | ok 0.218s |
| `go test -tags winlive -run 'TestBallLiveOrb'` | **5/5 PASS**（像素 41.94% 非黑；四角用例见下） |
| `go test -tags winlive -run 'TestLiveNeverStealsFocus|TestLiveTransparentCornerFallsThrough'` | **2/2 PASS** |
| `sh scripts/d22scan.sh clean` | clean（internal/=942 Go files 全量） |

**A 判据（编排者判据原样）**：`TestBallLiveOrbCornersReadTheDesktop`（新增，orb_live_windows_test.go）——
orb 移入工作区中心，脚下铺 (51,51,50) 已知底板窗（黑底板不可用：透黑与黑方块不可区分），四角
3px 内采样，**断言每通道与 (51,51,50) 差 ≤12**＋页面存活（与底板差 >25 的采样像素 >1%，死页不可
假绿）。多轮实测：**四角 = (51,51,50) 精确零差**（容忍 ±12），lit 25–93%。

**测试外真机复核（balldebug 二进制、非测试环境）**：PowerShell 建 (51,51,50) 底板窗插 orb 正下方，
四角实测 **(51,51,50)/(51,51,50)/(51,51,50)/(51,51,50)**，中心=orb 本体。截图：
`docs/evidence/s1/65-orb-over-known-floor.png`（底板上）、`65-orb-transparent-desktop.png`
（真桌面白窗背景上单球暗玻璃体，四周无黑方块）。

**B 判据（单球）**：EnumWindows 实测 `Wisp`（球窗）与 `Wisp Orb` rect **完全重合**
((3316,1268)-(3416,1368))，屏幕上只可能有一枚球；截图目视一枚暗玻璃球＋彩虹光带；
`TestBallLiveOrbFollowStaysInsideWorkArea` 钉死「orb rect == ball rect」的精确叠加。

**可复跑命令（编排者自跑）**：

```
go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 15m -v
```

看 `TestBallLiveOrbCornersReadTheDesktop` 的 log 行：`corners [[51 51 50] ×4] within +/-12`。

## 五、未做到／如实注脚

1. `go build ./cmd/wisp` 被并发腿在飞件压红（错误全文见门禁表），本腿零触碰 cmd/wisp/**。
2. 上腿的「像素非黑率摇摆」根因（泄漏堆积）由并发腿 79fdc5a8b 的 reap 环处理；本腿新增两枚
   winlive 用例都接了同一 reap 环。本腿轮次内未复现摇摆。
3. orb 模式下 dock 收缩成 tab 的视觉不可用（页面无 tab 状态）：docked 时 orb 以整球显示、
   位置夹进工作区；dock 交互语义（点击弹出/位置持久化）不变。D2D fallback 全保留。
4. `WEBVIEW2_DEFAULT_BACKGROUND_COLOR` 是进程级 env：面板 WebView2 同进程，其页面自带不透明
   CSS 底故无暴露面；若面板日后出现「未定义背景」形状需回看此处。
5. `CreateRectRgn`/`DwmEnableBlurBehindWindow` 失败仅告警不致命（缺它=黑角，球仍可用），实测未失败。

## 六、改动面

- `internal/ball/orb_windows.go` — 透明链重做（去 LAYERED＋DWM blur-behind＋env 透明底＋同 rect
  ＋中心叠加 z 序＋HTTRANSPARENT 兜底）。
- `internal/ball/orbasset.go` — served 副本注入透明背景 CSS（asset 零字节动）。
- `internal/ball/renderer_windows.go` — overlayOnly 模式（hit-test 轮廓＋宿主侧叠加件）。
- `internal/ball/ball_windows.go` — overlay 切换、拖拽逐帧跟随。
- `internal/ball/orb_live_windows_test.go` — exStyle 断言反转、follow 几何重钉、新增四角透明回归
  （彩色底板类）。
- `docs/evidence/s1/65-orb-over-known-floor.png`、`65-orb-transparent-desktop.png` — 实拍。
- 未碰：assets/orb-siri.html（sha d93da649… 钉测照旧绿）、cmd/wisp/**、internal/panel/**、
  frontend/**、design/**、thresholds/golden、票面 AC/Status。
