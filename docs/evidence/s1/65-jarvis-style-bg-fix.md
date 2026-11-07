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

---

## 步2 裁定：导出页自带 demo 深底 #050913，产品透明零改动（铁律 ⓪⓪⓪）

**底色裁定：jarvis 导出页 body 自带 `background:#050913`（上游 demo 原值，非截图估值 #0a0e14）。**

- **为什么这不是"硬塞黑底破坏透明架构"**：产品的透明契约**从来不由页面自己提供**——`orbasset.go` 对两页统一注入 `html,body{background:transparent !important}`（声明级 `!important` 压过页面任何非 important 底色，orbasset.go:125 注释逐字），再加 JS 引导双保险。**siri 页就是这个形状**：自带 `#030409` 深底（其上游导出目标也是浏览器 tab），产品里被注入层清成透明，`TestBallLiveOrbCornersReadTheDesktop` 四角 (51,51,50) ±12 长期钉死透明没被破坏。jarvis 页自带深底＝对称适用同一条既有规则。
- **为什么有价值**：单文件 artifact 的直接观感场景（浏览器直开、静态服务截图、人眼验收）从「白底灰环＝看起来坏了」变成「上游 demo 原貌＝深黑底＋亮青环」；产品内行为字节级不变（注入层照清）。
- **被否的备选**：保持 transparent——产品没问题但任何人浏览器直开永远白底灰环，机主实测凭据里的「差得很远」原样保留，等于没修；在 orbasset.go 加条件注入/按风格区分——复杂且毫无收益（注入层本来就两页同治，siri 页先例已验证）。
- **配色裁定：产码零改动。** palette 已显式传 cyan（entry.tsx:26-31）；three 的 toneMapping/outputColorSpace 与 demo 是同一段代码（同一组件、入口与 vite 配置不碰 three 初始化）。变灰唯一根源是底色，底色修掉后加色混合自然显青。
- **撤销口令**：`git revert` 本腿的 tools/jarvis-export/index.html 那笔 commit 并重建（或把 body background 改回 transparent）。

## 步3 重建产物（2026-11-07 19:1x +08 实跑）

- 改 `tools/jarvis-export/index.html`（body `transparent`→`#050913`，注释写明底色来源与产品注入层的关系）。
- `npm run build`（vite build + inline.mjs）：`inline: wrote …dist\jarvis-orb.html (759031 bytes)`，6.78s；三颗历史雷防护（inline.mjs 函数式替换 / vite resolve.dedupe / esbuild jsx automatic）全部在位未触发。
- 拷贝 `dist/jarvis-orb.html` → `internal/ball/assets/jarvis-orb.html`：759,031 B（旧 758,355 B，差值=模板注释＋底色行），**CR=0 / LF=4332**（纯 LF，PowerShell 逐字节复核）。
- 新 sha256 `740c486827e918f0993e63ea1f213ca217067ec1e9e942dbc04777897c5d2017`，钉值 `orb_test.go jarvisUpstreamSHA256` 同 commit 换钉（README 钉值纪律）。
- 静态尺现跑：`grep -c 'src="https?://|href="https?://'` = **0 命中**（离线不破）。

---

## 步4 判据加严：`orb_jarvis_live_windows_test.go` 整牙重铸（本票核心价值）

**旧牙（废）**：`orbTestNonDarkRatio > 1%`（any channel > 30 计 non-dark，黑底板）——白底页面 non-dark≈100% 恒过，对底色错误完全无感（步1 §1.6 结论）。

**新牙（三条读数全要）**：

| # | 牙 | 判据 | 为什么 |
|---|---|---|---|
| ① | **背景必须＝桌面（透明形态）** | orb 停在已知 **(51,51,50)** 底板上（`orbTestCreateColoredFloor`；黑板不可用——透明压黑板与黑方块缺陷读数同为 (0,0,0)），窗四角内缩 3px（编排者探针形状）逐通道 **±12** 内必须读出底板色 | 产品 orb 窗是透明 WebView2 叠加（orb_windows.go DefaultBackgroundColor a=0／无 LAYERED；orbasset.go 注入 transparent !important）⇒ 按「与桌面底板一致」形态（TestBallLiveOrbCornersReadTheDesktop 同法），**不是**绝对亮度阈值——绝对「暗<40」会把透明架构改成黑底架构，正是任务令点名禁止的破坏 |
| ② | **粒子必须是青色** | 全窗采样（内缩 4px 步长 3px）中**青色主导像素**（`g-r ≥ 12 且 b-r ≥ 12 且 b ≥ 80`）占比 **> 2%** | 上游 cyan 色板 primary 0x38f4ff／secondary 0x19a8ff／tertiary 0x7fffee（states.ts:117-123）：每个粒子色 g、b 都远高于 r；加色混合压在灰底板上青色主导保持。白底洗白的灰白环 g-r≈0 → 归零；重上色页面（非青系）也归零。阈值 2% ＝ 实测健康读数 7.91-8.32% 的 ~4 倍余量 |
| ③ | **旧「non-dark > 1%」断言废除** | 本测试不再引用 orbTestNonDarkRatio；「页面活着」由 ② 蕴含（青像素存在＝真画了） | 任务令点名：那条在白底下恒真 |

**诚实边界（②的判别力范围）**：上游三个内置色板 **cyan/aurora/ember 全是青系**（states.ts:117-141，aurora/ember 只是青色微调，demo 里的紫色是 App.tsx 自定义色板）⇒ ②钉的是「页面自身青系词汇真的在画」，不区分三块内置色板（那是 entry.tsx 的映射决策，已由 orb_test.go 的映射钉子管住）；②真正杀的是**白底洗白／灰白环／没在画**三类缺陷。

**实测读数（本腿实跑，2026-11-07 19:15-19:16 +08，attempt 1 即过）**：

- ① `corners [[51 51 50] [51 51 50] [51 51 50] [51 51 50]] within +/-12 of (51,51,50)` —— **四角与底板精确零差**（透明未破）
- ② `cyan-dominant 80/961 sampled px (8.32% > 2.00% floor)`（首跑 1% 牙线读数 76/961＝7.91%，提线到 2% 后复跑 8.32%）
- ③ 旧断言已从本文件移除（`orbTestNonDarkRatio` 零引用）

**保留段**：`testJarvisRuntimeStyleSwitch`（ApplyStyle glass→jarvis 往返、窗柄不变）原样保留；3 发有界重建＋泄漏 WebView2 收割的 flake 形状同族保留。文件自含探针常量（jarvisProbeInset/jarvisCornerTol/jarvisCyanGate/jarvisCyanBright/jarvisCyanFloor），不耦合 wave68 在 `orb_live_windows_test.go` 的在飞改动；只复用 HEAD 已有的原语（orbTestCreateColoredFloor/orbTestFloorColor/orbTestPixelAt/orbTestReapLeakedWebView2）。
