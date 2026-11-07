# 65 — Jarvis 粒子球接入产品（第二种视觉风格，真 Three.js 页面腿）

**Status:** in-progress（分步落盘，每步 commit；本件先记步1调查）
**Leg:** 65-jarvis-style-integration（机主点名要的第二种视觉风格：真 cyber1443/jarvis-ai-orb-web-animation 粒子球，不是 D2D 仿品）
**票面:** `.scratch/wisp/issues/65-ball-glass-quality-rework.md`
**上游:** `cyber1443/jarvis-ai-orb-web-animation`（MIT，package.json `"license": "MIT"` 亲读；Three.js/WebGL，React 组件）
**本机克隆:** `D:\tmp\wisp\20261107\orblab\jarvis-src`（已 npm install + tsup 构建成功；demo `examples/basic` vite dev server 实跑通，四变体截图在机主侧）
**既有格局（本腿之前）:**
- 玻璃球＝`LerSent001/orb`（MIT）真 WebGPU 页面，机主已签收：`assets/orb-siri.html`（sha `d93da649…` 钉死）＋ `orb_windows.go`（WebView2 悬浮宿主）＋ `orbasset.go`（served 副本＋CSP/透明注入）＋ `orbstate.go`（20 冻结态 → idle|thinking）。
- jarvis 风格现状＝`95e5a5934`（jarvis leg 3）的**纯 D2D 仿品**（`jarvis.go`＋`jarvis_windows.go`，D2D 画的粒子点环，D32 静帧）。那只在 orb 宿主 attach 失败（D2D fallback）时才可见；orb 宿主当前不分风格、永远服务 siri 玻璃球页。
- 本腿任务：让 WebView2 orb 宿主**按风格选 asset**——glass ⇒ siri 玻璃球页（原样），jarvis ⇒ 真 jarvis 粒子球页（新）。D2D fallback 面不动（`renderer_windows.go` ⛔ 碰）。

---

## 步1 调查：两边页面形态实测（2026-11-07 现跑）

### 1.1 orb-siri.html（已集成那枚）的页面形态

亲读 `internal/ball/assets/orb-siri.html`（89660 B，sha 钉在 `orb_test.go:24`）：

| 维度 | 实测 |
|---|---|
| 文件形态 | **单文件自包含**：`<!doctype html>` 起一个文档；`grep -o "<script[^>]*>"` 全文只有 **1 个** `<script type="module">`（WGSL shader 全部以 JS 字符串内联）＋ 1 个内联 `<style>` ＋ `data:` favicon；零外链 |
| 外链扫描 | `orb_test.go:59` 钉死 `http://` / `https://` / `src="http` / `fetch(` 零命中 ⇒ 天然离线 |
| 渲染栈 | WebGPU/WGSL（`requestAdapter`），需要 secure context——这是宿主用 `https://wisp.orb` 虚拟主机映射而不走 `NavigateToString` 的原因（`orb_windows.go:334` 起的实测记录） |
| 透明 | canvas premultiplied、clear a=0，orb 外 alpha 恰 0 ⇒ 宿主窗口透明链（DWM 空区域 blur-behind ＋ controller2 默认背景 a=0 ＋ served 副本注入 `background:transparent !important`）成立 |
| 宿主接口 | `window.liquidOrb = { getState, setState, setAudioBands }`（frozen，`defineProperty`）；状态词表 **idle|thinking** 两个 |
| serving 缝 | embed 字节原样 + serve 时注入 CSP `<meta>` ＋ 透明背景 `<style>`（`orbPageDocument()`，注入在内存、asset 文件不动）；虚拟主机映射到每用户 site 目录的 `orb.html`，失败降级 `NavigateToString` |

### 1.2 jarvis 上游的产物形态

亲读 `jarvis-src`（tsup.config.ts / package.json / dist/index.d.ts / src/*）：

| 维度 | 实测 |
|---|---|
| 打包 | tsup `external: ["react","react-dom"]`（peerDeps `^18 \|\| ^19`）＋ `noExternal: ["three"]` ⇒ **three@0.184 已打进** `dist/index.js`（1,177,498 B，未压缩带 banner）；**React runtime 不在 bundle 里**，import 时需要页面提供 react/react-dom |
| 组件 | `<JarvisOrb size state palette quality intensity paused interactive draggableSpin breathing …>`（`JarvisOrb.tsx` 347 行） |
| 状态词表 | **idle / thinking / success / alert** 四态（`STATE_TARGETS`，每态是 energy/rotationSpeed/particleSpeed/shellRadius/ringSpread/filamentOpacity/coreScale/bloom 的目标向量，renderer 内部 lerp 过渡） |
| 调色板 | 内置 **cyan / aurora / ember** ＋ 自定义数值 palette（`PALETTES`，core/primary/secondary/tertiary/deep/fallback） |
| 自适应 | `quality="auto"`：按 cores/memory/dpr/saveData 降档 ultra/high/balanced/performance（粒子数/丝数/环段/DPR/targetFps 全套缩）；`prefers-reduced-motion` ⇒ 整组件降级成 CSS 圆环 fallback（不跑 WebGL）；IntersectionObserver **离屏暂停** ＋ document.hidden 暂停；ResizeObserver 跟随容器 |
| 尺寸 | 容器 100%×100%（`width:"100%",height:"100%"`），canvas 绝对定位铺满 ⇒ **任意宿主窗口尺寸都自适应**，无需按 px 选 preset |
| 透明 | `renderer.ts:532-539`：WebGL context `alpha: true, premultipliedAlpha: false`，`renderer.setClearColor(0x000000, 0)` ⇒ orb 外像素 alpha=0，与 siri 页同一条透明链成立 |
| 渲染器 | `renderer.ts` 945 行（createRenderer(canvas,…)），**框架无关**（不 import react）；但 quality 解析、reduced-motion、resize/visibility 生命周期全部长在 `JarvisOrb.tsx`（347 行 React）里 |
| 音频 | **无音频接缝**：组件/渲染器都没有 setAudioBands 类入口（只有 intensity 覆盖、pointer pulse、breathing）——诚实记下：jarvis 页接不了音频带，语音家族在该风格下无音频驱动，这不做发明 |
| demo | `examples/basic`（vite ^5.4.11 + react 19 + `jarvis-ai-web-animation: file:../..`），node_modules 已装齐，dev server 已在 127.0.0.1:5174 实跑出四变体 |

### 1.3 能不能做成与 orb-siri.html 同形状的单文件 HTML？（判定）

**能。** 依据：

1. vite（esbuild/rollup）可以把 react + react-dom + three + jarvis 组件打进**一个**自执行 ES module bundle；再把 bundle **内联回 index.html**（`<script type="module">…</script>`，与 siri 页同形状）——这就是 `vite-plugin-singlefile` 插件做的事，插件不可用时一次确定性的字符串后处理（读 dist/index.html，把 `<script src>`/`<link rel=stylesheet>` 换成内联体）等价可得，零新依赖。
2. 内联后的文档外链为零，与 `orb_test.go` 对 siri 页的同一把外链尺过得了。
3. size 估计（按 minified 常见值）：three ~650 KB ＋ react-dom ~130 KB ＋ react ~7 KB ＋ jarvis src ~40 KB ⇒ 产物 **~1 MB 量级**（go:embed 进二进制的体积代价，如实记录；siri 页 88 KB，两者并存后二进制多 ~1 MB）。
4. WebGL 不要求 secure context（区别于 WebGPU），所以 jarvis 页连 `NavigateToString` 降级路径都能真渲染；但 serving 缝统一走同一个 `https://wisp.orb` 虚拟主机映射，保持与 siri 页同一条 no-socket 形状。

**唯一形态差异**：siri 页是上游原样 export（上游本来就有单文件形态）；jarvis 上游发布形态是 React 组件库，单文件 HTML 是我们用上游自己的组件 + 主流工具链（vite）构建出来的**集成产物**——组件源码零改动（甲案的核心纪律，与 A1594「ship the upstream page unmodified」同法）。

---

（步2 裁定与后续步骤见下文追加节；本文件随各步 commit 增量追加。）
