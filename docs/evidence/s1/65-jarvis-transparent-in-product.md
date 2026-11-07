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

（步2 完成后补：模板 transparent-by-construction、badge 跳过、重建与钉值。）

## 3. 步3 判据加严

（步3 完成后补。）

## 4. 步4 门禁

（步4 完成后补。）
