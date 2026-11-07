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

## 步2 裁定：集成形态选甲（vite 构建独立单文件 HTML） 三条候选的裁定与代价（判据＝铁律 ⓪⓪⓪：最有价值＋最符合主流 harness；⛔ 动作最小不进判据）：

### 甲（选）：vite 构建 jarvis 单文件 HTML → `internal/ball/assets/jarvis-orb.html`

- **为什么最有价值**：
  1. **主流 harness 的形状**：上游自己给的标准用法就是「React 组件 + vite 页面」（examples/basic 原样）；把组件库变成可部署页面走前端工具链构建，是行业标准做法。产物与 orb-siri.html **同一个页面形状**（单文件自包含 HTML），宿主 Go 侧的 serving 缝、CSP 注入、透明链、状态推送、测试牙**全部复用已验证的那一条**，只加一个「按风格选 asset」的纯决策点。
  2. **上游生命周期零 fork**：quality=auto 降档矩阵、prefers-reduced-motion 降级、离屏暂停、ResizeObserver、fallback 圆环——这 347 行宿主逻辑原样随组件走，上游更新时重新构建即可（组件源码零改动，与 A1594「ship the upstream page unmodified」同一条纪律）。
  3. **体积代价如实入账**：react+react-dom+three+组件 minified 内联 ⇒ 产物 ~1 MB 量级，go:embed 后二进制增 ~1 MB（siri 页 88 KB 并存）。这是买「真粒子球＋全套自适应」的实价，不是浪费。
- **边界（诚实）**：单文件 HTML 是我们用上游组件＋vite 构建的**集成产物**，不是上游的原始发布物；asset 顶部注释与证据件都写明来源与构建方法（脚手架随腿入库，可复跑）。

### 乙（不选）：不用 React，vanilla 薄入口直调 renderer.ts（945 行）

- **具体代价**：`renderer.ts` 确实框架无关，但 quality 解析矩阵（cores/memory/dpr/saveData → 四档全参数缩放）、reduced-motion 检测与 CSS fallback、ResizeObserver/IntersectionObserver/visibilitychange 三套 observer 生命周期、fallback 圆环——全部长在 `JarvisOrb.tsx`（347 行 React）里。乙案要**把这 347 行 fork 成 repo 自有代码长期自养**，上游一改就得跟。
- **买到什么**：只省 minified ~140 KB 的 react+react-dom。按铁律「复杂必须买得到真价值」：这 347 行复杂度换 140 KB，买不到真价值，反而制造一个与上游组件并行的第二真相源。**不选。**

### 丙（不选）：把 React+three 直接打进产品 Go 二进制（放弃单文件页形状）

- **先纠正一个直觉**：甲也 embed 进 Go 二进制，体积代价与丙同量级——丙被拒**不是**因为体积（体积在甲里已如实入账），而是因为**形态**：丙意味着放弃「orb 页＝单文件 HTML」这条已验证 serving 缝，要为 JS 多文件建新的模块加载/资源服务面（Go 侧新 surface、新测试面），与 siri 页两条 serving 形状并存。
- **买到什么**：零。同样的字节，拆成多文件只多出一条文件管理缝。这是铁律点名的「复杂但毫无收益」。**不选。**

### 两条风格级设计裁定（随裁定一并定案）

1. **D32 待遇**：D32 静帧纪律管的是 D2D fallback 渲染器（`jarvis_windows.go`，本腿零触碰，继续静帧）。orb 窗是「真上游页面」面——siri 页 idle 本来就持续跑上游自己的 RAF（机主已签收的形状）；jarvis 页同待遇：跑上游自己的动画循环（其自带 reduced-motion 降级与离屏暂停）。这不是本腿发明的豁免，是既有 orb 面规则的对称适用。
2. **音频带**：jarvis 上游**没有** setAudioBands 类接缝（组件与渲染器都没有音频输入），本腿**不做发明**——jarvis 风格下语音家族不推音频带（pushAudio 跳过），如实记录为该风格的已知边界；audio 驱动仍是 glass 风格的 siri 页能力。

---

## 步3 落地：jarvis-orb.html 构建＋真渲染验证

**产物**：`internal/ball/assets/jarvis-orb.html`（758,355 B，LF，sha256
`e95e1d963139a08b50d0c6be1d460db137a005d602e05eb9451dc60a6d99ce62`，钉值进 `orb_test.go`）。
**脚手架**：`tools/jarvis-export/`（entry.tsx ＝ 本腿唯一自有集成代码；node_modules/dist 已 gitignore；构建方法见其 README）。

**真渲染验证（不是"文件生成了"）**：Tabbit 真浏览器打开产物页，`canvas:true`、
`window.jarvisOrb` 接缝 `setState("thinking")`→`getState()` 回读 `"thinking"`，
idle/thinking 两态截图：

- `docs/evidence/s1/65-jarvis-orb-idle.png`（青色球核＋多椭圆轨道环＋粒子散点＋轴线）
- `docs/evidence/s1/65-jarvis-orb-thinking.png`（teal 高能版，环更密、球核更实）

**构建过程踩掉的三颗雷（复跑必读）**：

1. **`String.replace` 字符串替换模式污染**：minified JS 里 `$'`/`$&` 序列被替换引擎解释，
   把原 index.html 尾部（含外链 `src="/assets/…"`）拼进内联代码 ⇒ SyntaxError＋5 个假 script。
   修＝函数式替换 `replace(m, () => s)`（inline.mjs 内注释钉住）。
2. **双 React 实例**：`file:` 依赖是符号链接，上游克隆自己的 react devDep 抢走解析
   ⇒ 组件 hooks 打到另一个 React 的 null dispatcher（`null (reading 'useRef')`）。
   修＝vite `resolve.dedupe: ["react","react-dom"]`（vite 文档对此的标准修法）。
3. **JSX classic runtime**：esbuild 默认 `React.createElement` 而入口未 import React
   ⇒ `React is not defined`。修＝`esbuild: { jsx: "automatic" }`。

**静态尺（与 siri 页同法）**：`src="http…|href="http…|fetch(` 全文仅 1 处命中＝vite 自带
modulepreload polyfill 的代码标识符 `fetch(a.href)`（文档内无 modulepreload link，路径永不触达；
serving 缝 CSP `connect-src 'none'` 运行时硬断）。无任何 http(s) URL 引用。bundle 内有
react-dom SSR 模板的 `<script src=…>` 字符串字面量、无 `</script` 序列 ⇒ 内联解析安全。

**机主自跑验证命令**（产物已入库后）：

```
python -m http.server 8923 --directory "D:\work\workspace\projects plans\Wisp\internal\ball\assets"
# 浏览器开 http://127.0.0.1:8923/jarvis-orb.html —— 截图就看这个 URL
```

（步4 Go 接线见下文追加节。）
