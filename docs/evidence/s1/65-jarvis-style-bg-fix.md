# 65 — Jarvis 粒子球页面底色/配色修复 ＋ winlive 判据加严

**Status:** in-progress（分步落盘，每步 commit）
**Leg:** 65-jarvis-style-bg-fix
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`
**前件:** `65-jarvis-style-integration.md`（jarvis 页集成腿，产物 sha `e95e1d96…` 钉在 `orb_test.go`）
**机主实测凭据:** 上游 demo（`D:\tmp\wisp\20261107\orblab\jarvis-4variants.png`）＝深黑底＋亮青粒子环＋轨道环线＋白色核；我们产物页浏览器打开（`jarvis-product.png`，node 静态服务 8925 ＋ Edge --app）＝白底＋灰白粒子环。

---

## 步1 现跑调查（2026-11-07 19:0x +08，全部亲读亲核）

### 1.1 上游 demo 的深黑底到底是谁给的

亲读上游 `examples/basic/index.html`（`D:\tmp\wisp\20261107\orblab\jarvis-src\examples\basic\index.html:8`）：

```css
html, body { margin: 0; padding: 0; background: #050913; ... }
```

⇒ **深黑底来自 demo 宿主页面的 `body{background:#050913}`**（机主从截图估的 `#0a0e14` 偏亮一档；以源码为准取 `#050913`）。
`examples/basic/src/App.tsx` 全文**不设任何背景**（只有卡片层的半透明 `rgba(255,255,255,0.04)`）——「正确形状的参照」就是：**组件零背景，底色由宿主页面给**。

### 1.2 上游组件本身：默认 palette 与底色依赖

亲读 `jarvis-src/src/JarvisOrb.tsx`（347 行）：

- **palette 默认就是 `"cyan"`**（`JarvisOrb.tsx:124` `palette = "cyan"`）；不传也不至于 undefined。
- 组件不依赖宿主底色：容器只有 `position:relative; width:100%; height:100%`，fallback 圆环自带发光，**没有任何背景声明**。
- 亲读 `jarvis-src/src/renderer.ts`（945 行）：WebGL context `alpha:true, premultipliedAlpha:false`（:532-535），`renderer.setClearColor(0x000000, 0)`（:539）⇒ **canvas 本身 orb 外像素 alpha=0，透明**；`outputColorSpace = THREE.SRGBColorSpace`（:540）；**全文无 toneMapping 设置**（默认 NoToneMapping）。材质全部 `THREE.AdditiveBlending`（10 处命中）。

### 1.3 我们的导出入口现在传了什么

亲读 `tools/jarvis-export/entry.tsx`（73 行）＋ `tools/jarvis-export/index.html`：

- **palette 显式传了**：`PALETTE_FOR_STATE`＝idle→cyan、thinking→aurora、success→cyan、alert→ember（entry.tsx:26-31，每帧按当前态传）。⇒「没传 palette」**不成立**。
- **three 的 color space / toneMapping 与 demo 一致**：我们与 demo 跑的是**同一个组件、同一段 renderer.ts**，入口与 vite 配置都不碰 three 初始化（`vite.config.mjs` 只有 dedupe＋jsx automatic）。⇒「renderer 设置不同」**不成立**。
- **底色**：`tools/jarvis-export/index.html:9` 是 `body { overflow: hidden; background: transparent; }`——**入口给的是透明底，不是黑底**。透明 body 在浏览器里＝透出浏览器默认白。⇒ 机主假设①方向对（入口没给深底），但准确根因是**给了 transparent 而不是 demo 的深底**。

### 1.4 变灰的机制（裁定：AdditiveBlending × 白底，机主假设③成立）

上游 10 处材质全 `AdditiveBlending`：粒子颜色**加**到已有底色上。底＝黑（#050913）时加出亮青；底＝白时青色分量加到接近饱和 ⇒ 环被洗成灰白。这不是我们的代码改了颜色，是**底色选错了**（白）导致加色混合的观感塌掉。palette/renderer 两条已排除（1.3）。

### 1.5 产品里 orb 窗到底该不该有底色（现读产码裁定）

亲读 `internal/ball/orb_windows.go` ＋ `internal/ball/orbasset.go`：

| 证据 | 读数 |
|---|---|
| 窗样式 | `orb_windows.go`：`WS_POPUP + WS_EX_TRANSPARENT|TOPMOST|TOOLWINDOW|NOACTIVATE`，**WS_EX_LAYERED 被刻意去掉**（黑方块缺陷，:29-39 注释）；透明链＝`DwmEnableBlurBehindWindow` ＋ `ICoreWebView2Controller2::PutDefaultBackgroundColor({A:0,...})`（:336） |
| served 注入 | `orbasset.go:133`：对**两页**注入 `<style>html,body{background:transparent !important}</style>`——`!important` 声明级压过页面自己的任何非 important body 底色（:125 注释逐字："transparent !important beats any opaque body rule"） |
| JS 引导 | `orbasset.go:170-182` `orbPageBootstrap` 再清一遍 body 背景（ExecuteScript 缝，防 race 的双保险） |
| siri 同构 | `assets/orb-siri.html:10` 自带 `body{background:#030409}` 深底——**siri 页本来就是「自带深底＋产品里被注入层清掉」**，产品透明从未依赖页面自己 transparent |
| 测试钉 | `orb_live_windows_test.go` `TestBallLiveOrbCornersReadTheDesktop`：orb 停在已知 (51,51,50) 底板上、四角 ±12 必须读出底板色 ⇒ **产品透明是测试钉死的契约** |

### 1.6 步1 结论（根因链）

1. **产品架构＝透明叠加窗**（1.5 五条证据），orb 页面在产品里**必须有透明底**——但这一点由 `orbasset.go` 的注入层保证，**不依赖页面自己写 transparent**。
2. 白底灰环的唯一根源＝**导出页 body 写了 `background:transparent`**：浏览器/截图场景透出白、加色混合洗白；这与「页面该不该透明」无关——页面自带什么底，产品里都会被注入层清成透明。
3. 判据缺陷确认：现 `orb_jarvis_live_windows_test.go` 唯一牙齿＝`orbTestNonDarkRatio > 1%`（orb_live_windows_test.go:385 定义，any channel > 30 计 non-dark）——**页面若真把白底画出来，non-dark 接近 100%，照样绿**；对「底色错了」完全无感，对「粒子是不是青的」也无感。
