# 票 65 · jarvis 产品态透明与 badge 缺陷（真机实测两缺陷收口）

> 起因（机主 2026-11-07 真机实测凭据）：`build/balldebug-orb.exe -style jarvis`，
> 两窗 rect 完全重合 (3316,1268)-(3416,1368)，青色粒子环＋轨道环线在画
> （球心 RGB=(170,203,204)），但 ①同一截图左上角 RGB=(255,255,255)（疑「页面没清透明」）
> ②球窗上有一个红色圆点 badge 叠在粒子上。
> 本件 = 步1 现跑取证 → 步2 落地 → 步3 判据加严 → 步4 门禁，逐节 append。

## 1. 步1 现跑取证（2026-11-07 20:0x，HEAD 007f1ee21）

### 1.1 注入层对两个 style 是否一视同仁 —— 是（同一个函数、同一个注入串）

`internal/ball/orbasset.go`：

- `orbServeDocument(page []byte)`（orbasset.go:91-107）是**唯一**注入函数，
  siri 与 jarvis 共用：`orbPageDocument()`（:62）与 `orbJarvisDocument()`（:70）
  都直接调它。注入串（:93-95）：
  `<meta http-equiv="Content-Security-Policy" ...>` ＋
  **`<style>html,body{background:transparent !important}</style>`**，
  打在文档第一个 `<head>` 之后；找不到 `<head>` 时退化为 doctype 后注入
  （:102-105），策略永不静默丢弃。
- `orbPageForStyle`（:203-208）：StyleJarvis → `orbJarvisDocument()`；否则 siri 页。
  **两个 style 走同一条注入路径，不存在「jarvis 分支没注入」。**
- 落盘路径：`orbWriteSite`（:110-127）每次 attach 把**两个** entry
  （orb.html / jarvis.html，均为注入后字节）写进
  `%LOCALAPPDATA%\wisp\orb\site`（orb_windows.go:363,411-417），WebView2 经
  虚拟主机映射加载的就是这两份文件。

### 1.2 jarvis 页的底色规则 —— 一处、无 !important

`internal/ball/assets/jarvis-orb.html`（759031B，sha
740c486827e918f0993e63ea1f213ca217067ec1e9e942dbc04777897c5d2017，d8da4e83e 落）：

- grep background：正文 CSS 只有一处规则 `body { overflow: hidden; background: #050913; }`
  （:18，上游 demo 原值），**无 !important**；其余 background 命中全部是
  three.js 着色器源码字符串（`backgroundIntensity` 等 uniform 名），不是 CSS。
- CSS 裁决：注入串 `transparent !important`（作者层 !important）对页面
  `#050913`（作者层普通声明）——**!important 必胜，与出现顺序无关**。
  页面没有任何能与注入串打平或压过它的规则。
- canvas 本身透明（排除「WebGL canvas 自带不透明底」假设）：构建产物内
  three.js 渲染器以 `alpha:!0, premultipliedAlpha:!1` 创建、
  `setClearColor(0,0)`（黑 α0）——canvas 四角像素 α=0，不是不透明面。

### 1.3 D2D overlay 的 badge/圆盘绘制点

`internal/ball/renderer_windows.go` `drawFrame`：

- overlay 分支（`r.overlayOnly`，:473-492）只画 alpha=1/255 命中剪影——
  **overlay 模式下角落不可能由 D2D 画出白底**。
- badge 圆盘＝**step 11（:575-593）**：`v.BadgeCount > 0` 时无条件画
  BGInset 边环 ＋ Danger 红盘（数字在 overlay 模式被 :588 `!r.overlayOnly`
  压掉）——overlay 模式保留「无数字通知点」是 wave 早前的裁定
  （「mainstream badge-dot shape」）；**它不分 style**，jarvis 粒子球上照样画，
  这就是机主看到的红点。跳过判据落点＝此处按 `activeStyle` 分流。
- 参考：queue dot＝step 12（:595-603）、进度环＝step 9（:556-564）
  （页面的冻结 API 画不了进度，step 9 保留不动）。

### 1.4 ★真机复现（本机＝机主同一台机器，2026-11-07 20:02）

现行 winlive 判据 `TestBallLiveOrbJarvisStyleRendersParticles` 现跑：

```
attempt 1: jarvis orb style proof: corners [[51 51 50] [51 51 50] [51 51 50] [51 51 50]]
  within +/-12 of (51,51,50); cyan-dominant 68/961 sampled px (7.08% > 2.00% floor)
--- PASS: TestBallLiveOrbJarvisStyleRendersParticles (8.17s)
```

四角读数与已知地板 (51,51,50) **逐通道零偏差**（容差 ±12）——真实产品路径
（`ball.New` → 真实 WebView2 attach → orbWriteSite 落盘 → 虚拟主机加载注入后
jarvis.html）在像素层成立：页面四角合成结果就是窗后景物，页面没有自绘任何底色。

### 1.5 ★裁定：四角白不是「页面没清透明」，是透明球后面那块白色的桌面

机主自己录的关键旁证：NBA2K 全屏盖桌面 ⇒ CopyFromScreen 连**远处桌面**都读
(255,255,255)。三段推理：

1. 页面自己的底色是 #050913（**暗色**）：若注入失败、页面底色漏出来，
   四角应读**暗色** (≈5,9,19)，不可能是白色。白 ⇒ 漏出来的不是页面底色。
2. 球心 (170,203,204)＝青粒子合成在白色游戏面上的读数——粒子在画说明页面
   活着且参与合成；若 WebView2 整窗是不透明白底（DefaultBackgroundColor 未
   生效），粒子会被白底洗掉且四角与球心同为纯白，与实测读数不符。
3. 本机 20:02 现跑（1.4）：同一台机器、同一份代码、真实产品路径，四角
   **精确等于**窗后已知色。transparent 成立。
4. 补充对账：机主 19:46 构建的 `build/balldebug-orb.exe` 与终收口
   007f1ee21（19:56）之间代码零差异（该 commit 只动 2 个文档件，
   `git show --stat` 复核）⇒ 排除「旧 exe 缺修复」。

**结论：jarvis 页在产品里已是真透明；机主截图的四角白＝透明球窗后面的
白色游戏面。** 残余的真实缺陷只有一枚：overlay 模式的 D2D badge 圆盘（1.3）。

## 2. 步2 落地

### 2.1 jarvis 页透明改为「构造即透明」（手法裁定：模板自带 transparent，注入层保留为 belt）

选型理由（铁律 ⓪⓪⓪）：步1 已证注入层在产品里**赢**（1.4 四角零偏差），但「透明」目前
**依赖**注入层在级联里压过页面自身的 `#050913`——这是一条还存在的依赖链。主流 harness
对 overlay 页的做法是页面**自身**不画任何底色（透明 by construction），注入/宿主兜底只是
belt。因此把 `tools/jarvis-export/index.html` 的 body 规则改为
`body { overflow: hidden; background: transparent; }`（注释写明动机与代价），重建单文件
产物。这样即便未来 serving 缝丢了注入 belt，产品页也不会再画出底色——整类「注入优先级
被顶掉」的风险被消除；siri 页（钉死 sha 的上游页）继续由注入层负责，不受影响。

- 重建：`npm run build`（vite build + inline.mjs），产物
  `internal/ball/assets/jarvis-orb.html` **759648B**（旧 759031B），纯 LF（CR=0），
  属性形外链 0 命中（构建期离线尺保持）。
- 新 sha256：`66eb4b800e2f37c77ce2b3e9a7d3d4a6a11dd102b85097bc06d7732f2c9c856b`，
  同 commit 换钉 `orb_test.go jarvisUpstreamSHA256`（README 既有纪律：asset 变钉值必须
  同 commit 换）。产物内仅注释文本残留一处 `#050913` 字样（无 CSS 效果）。
- 代价（已知并接受）：单独在浏览器标签页打开该 artifact 时，粒子叠在标签页自己的底色上
  （白底上青粒子会偏淡）——demo 深底观感仍在 `jarvis-src/examples`；该 artifact 的唯一
  消费方是产品嵌入资产。
- 撤销口令：把 `tools/jarvis-export/index.html` 的 body 规则改回
  `background: #050913`，`npm run build` 重建、拷贝、把钉值换回
  `740c4868...2017` 即复原。

### 2.2 jarvis 风格下 overlay 不画 D2D badge 圆盘

`internal/ball/renderer_windows.go` step 11（原 :575-593）条件改为：

```go
if v.BadgeCount > 0 && !(r.overlayOnly && activeStyle == StyleJarvis) {
```

- 只收窄 overlay 态（页面在脸上）＋ jarvis 家族：glass 家族的 badge-dot 行为不变；
  orb 未附着（D2D fallback 面）时所有风格照画——审批深度仍可达用户（面板＋fallback 面）。
- 进度环（step 9）、queue dot（step 12）、命中剪影都不动。
- 撤销口令：把该条件改回 `if v.BadgeCount > 0 {` 即复原。

## 3. 步3 判据加严

### 3.1 为什么测试能过而真机「白底」——测试没有底板假设错误，白是窗后景物

现行四角牙的底板＝测试自建 (51,51,50) 已知色窗（orbTestCreateColoredFloor，topmost、
插在 orb 正下方）——这不是「把页面背景改成已知色」，页面与 serving 缝零特设；它模拟的
是「窗后有一块已知桌面」。真机上四角读白的原因＝**窗后桌面被 NBA2K 盖成白色**（机主
自记：连远处桌面都读 255,255,255），透明球叠在白面上四角就是白。三段证明见 §1.5
（漏底只会漏暗色 #050913；球心粒子读数与透明合成一致；同机同码 20:02 复跑四角零偏差）。

### 3.2 为什么不加「四角==环境桌面色」的产品态牙（如实说明差异）

「断言四角读到的就是桌面色 ±12」在真机上**不可证伪**：透明球后面的游戏画面是活动
图像，两次采样之间像素自己会变，即便球完全不透明，比较也会随机红/随机绿——这条牙
在真机上测不出任何东西。测试环境（静止桌面）能做这种比较，但那正是机主指出的
「测试环境与产品不同」之处。**替代判据（覆盖真机的两条）**：

1. **served-page 牙**（新，TestBallLiveOrbJarvisServedPageAndOverlay 相位 1）：真 attach
   之后读产品 WebView2 实际加载的那份文件 `%LOCALAPPDATA%\wisp\orb\site\jarvis.html`
   （orbWriteSite 每次 attach 重写，非测试特设路径），断言
   (a) 注入 belt `html,body{background:transparent !important}` 在，
   (b) 页面自身 body 规则已是 `background: transparent`（2.1 构造即透明），
   (c) 无 `background: #050913` 不透明底规则。这条牙在任何机器、任何桌面背景下都成立，
   把「页面自己画底色」这类缺陷整类钉死。
2. **角点像素牙保持**（TestBallLiveOrbJarvisStyleRendersParticles 现行四角 ±12 对已知
   地板）：这是合成像素层的真值判据——若 DefaultBackgroundColor 或页面合成出问题，
   四角不可能与地板逐通道零偏差。真机白角疑云的像素层排除靠它＋1.5 的暗色反证。

### 3.3 无 badge 牙（新，同文件相位 1/2）

- 相位 1：jarvis 家族＋orb 已附挂，`SetBadge(3)`（真产品 setter，STA 队列保序）后，
  在 STA 线程走球窗 ULW 源 DIB，数「红主导」像素（α≥100 且 R≥1.5G 且 R≥1.5B——
  Danger 两行主题 #E07A70/#BC544C 都满足，且对预乘不透明度缩放不变；overlay 分支
  允许画的其余元素——BorderSoft 近白环、Accent 青弧、Info 蓝点、α=1/255 剪影——
  全部不进这个色带）。断言 **==0**。
- 相位 2（阳性对照）：同一 badge 深度切到 glass 家族，断言红主导像素 **≥50**——
  证明审计看得见圆盘，jarvis 的 0 是风格域跳过而不是审计瞎。

## 4. 步4 门禁（2026-11-07 20:0x-20:1x，全部现跑）

| 门 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空（GOFUMPT_CLEAN） |
| `go vet ./internal/ball` | rc 0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc 0 |
| `go test ./internal/ball -count=1` | ok 0.300s |
| ★`go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 25m` | **8/8 PASS 58.605s**（旧 7 枚＋本腿新 1 枚；重建后页面上重跑） |
| ├ TestBallLiveOrbJarvisStyleRendersParticles | PASS 8.32s：四角 [[51 51 50]×4] 对地板 ±12 零偏差；青主导 71/961＝7.39%（构造即透明的新产物上成立） |
| ├ TestBallLiveOrbJarvisServedPageAndOverlay（新） | PASS 0.91s：served jarvis.html 759931B＝注入 belt 在＋body `background: transparent` 构造即透明＋无 `#050913` 底规则；badge 深度 3 下 jarvis overlay 红主导像素 **0**；glass 阳性对照同深度 **245** px（floor 50）——审计看得见圆盘 |
| ├ TestBallLiveOrbAttachStateAndTeardown / FollowPlacement / FollowStaysInsideWorkArea | PASS |
| ├ TestBallLiveOrbPageRendersPixels | PASS 26.70s |
| ├ TestBallLiveOrbCornersReadTheDesktop | PASS 20.24s（glass 四角零差未破） |
| └ TestBallLiveOrbOverlayFrameHasNoGlyphPixels | PASS 0.74s（glass 徽章点行为未动，近白 0＋青弧在） |
| ★`go test -tags winlive ./internal/ball/ -run 'TestLiveNeverStealsFocus\|TestLiveTransparentCornerFallsThrough' -count=1` | **2/2 齐绿**：CornerFallsThrough PASS 0.53s；NeverStealsFocus 第 1/2 次 **SKIP-LOUD**（SetCursorPos 被环境拒绝＝全屏游戏在桌面上的现场条件，跳过前 NOACTIVATE/前台不变断言已全绿），第 3 次同命令 **PASS 1.21s**——如实记录非首跑绿 |
| `sh scripts/d22scan.sh` | clean（internal 956＋cmd 234＋design 39＋frontend 116＋risk 21） |

环境差异如实：winlive 现场桌面被 NBA2K 全屏占用中（机主在飞），已知 CopyFromScreen 远处
桌面读白——本腿全部判据按机主指示用像素读数（合成屏 GetPixel 的四角牙读的是 orb 自己
窗后已知地板；DIB 牙走 ULW 源位图，与屏幕抓取无关），未受影响。

## 5. 收件对账

- commits（本腿，只 commit 不 push）：3104ac37a（步1 证据）／步2（模板透明重建＋钉值＋badge 跳过）／步3（新判据）／步4（本节＋Progress log）。
- 写面：tools/jarvis-export/index.html、internal/ball/assets/jarvis-orb.html（重建产物）、internal/ball/orb_test.go（sha 钉）、internal/ball/renderer_windows.go（badge 跳过）、internal/ball/orb_jarvis_live_windows_test.go（新牙）、本证据件、票面 Progress log 一行。cmd/wisp、internal/panel、internal/config、internal/tools、frontend、design、assets/orb-siri.html、thresholds/golden、票面 AC 勾与 Status：零触碰。
- 如实两项：①「四角白」判定为窗后白色游戏面而非页面缺陷（§1.5/§3.1），jarvis 页在产品里此前已是真透明，本腿把透明升级为构造保证；②NeverStealsFocus 非首跑绿（前两次环境 SKIP-LOUD）。


