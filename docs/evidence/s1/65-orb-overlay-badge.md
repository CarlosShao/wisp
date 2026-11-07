# 65 — orb 球面上的白色进度文字（overlay badge/text，2026-11-07）

**腿名**：wave65-overlay-no-text（写腿，只 commit 不 push）
**机主凭据**：编排者真机实测——build `build/balldebug-orb.exe`，EnumWindows 得「Wisp」与
「Wisp Orb」两窗 rect 完全重合 (3316,1268)-(3416,1368)（单球架构已落实，上一腿交付）；
区域截图放大 4 倍（`D:\tmp\wisp\20261107\orblab\orb-FINAL.png`）：球体、彩虹光带、边缘青色
高光弧全对，**但球正中央叠着一坨白色「69%」**。

## 一、修前读码（overlayOnly 分支今天画了什么）

`internal/ball/renderer_windows.go` `drawFrame` 的 `overlayOnly` 分支（orb 附加时球窗每一帧走的路径）：

| 绘制项 | 修前行为 | 判定 |
|---|---|---|
| hit-test 轮廓 | `pal.OrbBody.WithAlpha(1/255)` 椭圆（alpha byte 恰 1，全零＝OS 跳过＝不可点，上腿实测回归） | 合法，保留 |
| 步骤 9 进度环 | 无守卫直画：`BorderSoft` 底环 2px @ R+5px ＋ `Accent` 进度弧（-90° 起、360°×progress） | 合法（弧线不是文字，不压球面读数），保留 |
| 步骤 11 badge | 无守卫直画：`BGInset` 外圈＋`Danger` 红盘＋**白色数字**（`OnSolid` `drawText`，`BadgeFontPx` 字号） | **违判据**（白字压球面），数字须抑 |
| 步骤 12 queue dot | 无守卫直画：`BGInset`＋`Info` 小圆点（无文字） | 合法，保留 |
| 步骤 13 中心文本 | 无守卫直画：`OnSolid` 白色 `drawText`（`FontSizeMicroPx`）——「L1 倒计时秒数／下载百分比」 | **违判据＝机主截到的「69%」本体**，须抑 |

**wild 缺陷形状实测**（编排者截图逐像素机读，本腿复核）：150px 窗口区 ×4 放大 ⇒
白色近白像素 173 枚、x∈[257,336] y∈[284,312]（中心 ±11 窗 px＝文本块）；青色弧半径
≈21.7 窗 px（=D2D 进度弧 R+5=33 放大 4 倍，`Accent` 色）；**红色像素 0 枚**（badge 盘不在野）。
⇒ 野缺陷＝`cmd/wisp` presenter 经 `SetBadgeText("69%")` 挂载（`resident_present_windows.go:16,222`
的 `presenterBadger(rb.b.SetBadgeText)`），被步骤 13 无守卫直画上球面。

## 二、技术选型（铁律 ⓪⓪⓪：按「最有价值/主流 harness」拍，⛔ 按「动作最小」）

| 候选 | 内容 | 裁定 | 理由 |
|---|---|---|---|
| **甲** | orb 模式球窗**完全不画文字/数字**；进度改**球外圈细弧** | **采纳** | 主流圆形进度环（Echo/语音 harness 的标准进度语言）：进度信息不遮挡球体、可读、且 D2D 弧与页面玻璃球在 z 序上天然分层（球窗在 orb **之上**，弧贴 R+5px 落在页面球缘外的留白里）。真实价值：球面恢复纯净（机主判据），进度/审批信号仍宿主可达（页面 API 无这些通道） |
| 乙 | orb 模式只画 hit-test 轮廓，**什么都不显示**，进度交给 orb 页面自己渲染 | **排除** | 前提不成立：页面冻结 API 只有 `window.liquidOrb = { getState, setState, setAudioBands }`（orbasset.go 顶部钉死的唯一接缝，`orb_test.go` 钉死脚本体），**没有进度/文本/badge 通道**；要加就得改 `assets/orb-siri.html`——其 sha `d93da649…` 被票面钉死（上游逐字节嵌入是本票的票值核心），⛔ 动。且「什么都不显示」让 Downloading/AwaitingApproval 在 orb 模式下丢全部宿主侧信号 |
| 丙 | 进度数字移到球下方（不压球面） | **排除** | 球窗是方窗＋逐像素 hit-test，球下方就是窗外的桌面——「移到球下方」要给球窗加高度或另立文字窗：动冻结窗口几何（`WindowEdgePx`、dock 收缩、四角落空断言全挂在现几何上），复杂买不到价值；且一个悬浮的孤立百分数浮在桌面上不是任何主流 harness 的形状 |

**甲的落地边界（如实）**：
- 「细弧」**不用新画**——overlayOnly 分支的步骤 9 进度环本来就是球外圈细弧（底环 `BorderSoft`
  2px @ R+5px ＋ 进度弧 `Accent`），修前就在画、截图里的青弧就是它。本腿只做**减法**：
  把文字/数字从 overlay 分支摘掉（步骤 13 整体抑止；步骤 11 保留 `Danger` 盘作**无数字
  通知点**——主流 badge-dot 形状，审批「有货」的信号不丢，具体深度看面板）。
- D2D fallback（orb 未附加）**一字未动**：完整帧（含数字）原样保留——fallback 是
  attach 失败时的唯一可见面，摘文字会让它残废。

## 三、修法（改面与机制）

1. `internal/ball/renderer_windows.go`
   - 步骤 13：`if v.BadgeText != "" && !r.overlayOnly` —— orb 模式不画中心文本；
   - 步骤 11：数字 `drawText` 包进 `if !r.overlayOnly` —— orb 模式保留红盘、不画数字；
   - `overlayOnly` 字段注释与 `ball_windows.go renderFrame` 注释同步改写为
     「NO text of any kind」法则（防下一位照旧注释把文字加回来）。
2. `internal/ball/ball_windows.go` — 仅注释（机制零改：overlayOnly 切换、ULW 路径原样）。
3. `internal/ball/orb_live_windows_test.go` — 新回归用例（下）。
4. `internal/ball/renderer_windows.go` — DIB 审计访问器 `overlayGlyphPixels()`：
   `bits` 字段从 `uintptr` 改 `unsafe.Pointer`（创建处两处同步）——该指针今天在 Go 侧
   零解引用，vet-clean 的 `unsafe.Slice` 读法需要类型化指针；生命周期注释照旧
   （DIB 内存归 `hbm` 管，`release()` 的 `DeleteObject` 不变）。

## 四、判据（机主硬性判据的机械形式）＋ 新回归用例

**判据**：修后球面上不能有文字/数字 ⇒ 断言「球窗在 orb 模式下没画文字像素」，
落点＝**球窗 ULW 源 DIB 的逐像素审计**（球窗是叠加在 orb 页之上的那一层，机主截到的
「69%」就是这层画上去的——审计这层＝审计缺陷本体，不用截屏 OCR）：

`TestBallLiveOrbOverlayFrameHasNoGlyphPixels`（`orb_live_windows_test.go`，winlive）：
- 真球＋真 orb 附加（quiet-desktop 守卫＋泄漏 reap 环照旧）；
- 复刻**缺陷形状**：`SetState(Downloading)` → `SetProgress(0.69)` → `SetBadgeText("69%")`
  → `SetBadge(3)`（与 `cmd/balldebug` 缺陷复现路径、`cmd/wisp` presenter 调用形状一致）；
- 断言 `overlayGlyphPixels() == 0`：计数器只认**不透明近白像素**（alpha≥128 且 RGB 三通道
  ≥ alpha-40，预乘近白判据）——overlay 分支的合法调色板（轮廓 1/255、`BorderSoft`/`Accent`
  弧、`BGInset`/`Danger` 盘、`Info` 点）**没有任何近白元素**，`OnSolid` 白是唯一近白来源
  ＝`drawText` 专用色 ⇒ 计数非零即文字/数字像素压面；
- 再断言 Sleeping 静置帧同样 0（轮廓-only 法则）；
- 审计经 `sta.PostTask` 入队：FIFO 排在 setter 们的重绘任务之后，读到的是**已定格的
  present 帧**，无轮询无竞态。

**机主可自跑的验证命令**（PowerShell，仓根）：

```
$env:GOTMPDIR=""; & "D:\work\soft\Git\bin\bash.exe" -lc 'cd "/d/work/workspace/projects plans/Wisp" && unset GOTMPDIR TMPDIR TEMP TMP && export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH && go test -tags winlive ./internal/ball/ -run "TestBallLiveOrbOverlayFrameHasNoGlyphPixels|TestBallLiveOrbCornersReadTheDesktop" -count=1 -timeout 20m -v'
```

两枚都 PASS ＝ 球面无文字像素 ＋ 四角仍 (51,51,50) 精确读数。

## 五、门禁读数（全部本机实跑，2026-11-07）

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空（rc0） |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1` | ok 0.645s |
| ★`go test -tags winlive -run 'TestBallLiveOrb' -count=1 -timeout 20m` | **6/6 PASS**（含新增用例）：Attach 2.51s / FollowPlacement 0.00s / FollowStaysInside 0.27s / PageRendersPixels 29.77s / **CornersReadTheDesktop 30.84s PASS**（四角 (51,51,50) 法则未破）/ OverlayFrameHasNoGlyphPixels 1.27s |
| ★`TestLiveTransparentCornerFallsThrough` | PASS（0.26s） |
| ★`TestLiveNeverStealsFocus` | **PASS（1.65s）**——重试环第 2 发拿到判行（首几发物理指针被占用 SKIP-LOUD：SKIP 前样式位 `WS_EX_NOACTIVATE`/`MA_NOACTIVATE`/前台不变三断言已过；第 1 发无判行输出＝票面已记的"二进制静默死亡"同形状，第 2 发正常判） |
| `sh scripts/d22scan.sh` | **clean**（ban #8: internal/=944 Go files、cmd/=230、frontend/=116；首跑 FAIL＝tools/d22scan 自模块 1 枚，与本腿文件无关、系并发腿并行污染，隔离复跑 ok 73s，全量复跑 clean） |

**红→绿证据链**：修复未应用时新用例 **FAIL，计数 103**（`overlay frame carries 103 near-white
glyph pixels`）＝「69%」的文字像素确实在球窗帧里；修复后同用例 **PASS，计数 0**。判据有牙
（突变敏感：把抑止摘掉即红）。

## 六、改动面

- `internal/ball/renderer_windows.go` — overlay 分支去文字（步骤 11 数字/步骤 13 文本）＋
  `overlayGlyphPixels()` 审计 ＋ `bits` 类型化 ＋ 注释法则化。
- `internal/ball/ball_windows.go` — 注释同步（机制零改）。
- `internal/ball/orb_live_windows_test.go` — 新增 `TestBallLiveOrbOverlayFrameHasNoGlyphPixels`。
- 未碰：`assets/orb-siri.html`（sha `d93da649…` 钉测照旧绿——默认套件含 orb_test 的 sha 钉测）、
  `cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`、`frontend/**`、`design/**`、
  thresholds/golden、票面 AC 勾/Status。

## 七、未做到／如实注脚

1. **进度百分比在 orb 模式下当前无数字读数**（甲的代价，如实）：orb 模式进度＝外圈弧形
   （底环＋进度弧比例可目读），具体数字要看面板。页面 API 冻结（sha 钉死）是硬边界，
   数字要上玻璃面只能等「上游页加通道或换页」的票（那要动票面钉死的 sha，属改契约）。
2. **TestLiveNeverStealsFocus 收尾读数已落**：后台重试环（60s 间隔）第 2 发 **PASS 1.65s**
   （物理 SetCursorPos+点击真跑过、前台窗口未变）。首几发 SKIP-LOUD＝本机物理指针被占用
   （编排者/机主在用鼠标），该 SKIP 纪律拒绝伪绿、非测试缺陷；第 1 发曾无判行输出
   （票面在记的 winlive 二进制静默死亡同形状，未复现第二次）。
3. ~~d22scan 后台在飞~~ → **已落：clean**（首跑 FAIL＝tools/d22scan 自模块测试 1 枚，
   与本腿无关；当轮与后台焦点重试环并行，隔离复跑 ok 73s，全量复跑 clean 读数见 §五）。
4. 工作树上有一枚**并发 jarvis 腿的在飞未跟踪件** `internal/ball/assets/jarvis-orb.html`
   （commits 1652051d0/29b3501d5 的后续 WIP）：非本腿产物，pathspec 点名时不带它。
