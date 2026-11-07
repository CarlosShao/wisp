# tools/jarvis-export — jarvis-orb.html 构建脚手架（票 65 第二视觉风格）

把上游 `cyber1443/jarvis-ai-orb-web-animation`（MIT）的 React 组件构建成
**一个自包含单文件 HTML**，落到 `internal/ball/assets/jarvis-orb.html`
（go:embed 进 `orbasset.go`，与 `orb-siri.html` 同一条 serving 缝）。

本目录是构建脚手架，⛔ 不是产品代码；`node_modules/` 与 `dist/` 已 gitignore，
入库的只有：入口源码（`entry.tsx` ＝ 本腿唯一自有集成代码）、构建配置、
内联器、本 README。

## 构建（本机已跑通的路径）

```bash
cd tools/jarvis-export
npm install            # 依赖：react/react-dom/jarvis-ai-web-animation/vite
npm run build          # = vite build && node inline.mjs
```

- `package.json` 里 `"jarvis-ai-web-animation": "file:D:/tmp/wisp/20261107/orblab/jarvis-src"`
  指向本机克隆（该克隆已 `npm install` + `tsup` 构建验证过）。
- 干净机器复跑：`git clone https://github.com/cyber1443/jarvis-ai-orb-web-animation`
  后把 file: 指过去即可（上游组件零改动）。
- 产物默认写 `dist/jarvis-orb.html`；拷贝到
  `internal/ball/assets/jarvis-orb.html`（LF；若 CRLF 先转 LF），
  然后**更新 `orb_test.go` 里的 `jarvisUpstreamSHA256` 钉值**
  （与 siri 页同一纪律：asset 变了钉值必须跟着换，测试强制）。

## 单文件是怎么来的

1. `vite build`：esbuild/rollup 把 react + react-dom + three + 上游组件
   minify 成一个自执行 ES module（`dist/assets/index-*.js`）。
2. `inline.mjs`：把 module script 与 stylesheet 内联回 `dist/index.html`
   （等价 vite-plugin-singlefile，零额外依赖），并**拒绝任何 http(s) 外链**
   ——离线是构建期保证，与 siri 页同一把外链尺。

## 入口（entry.tsx）的接缝决定

- 挂上游 `<JarvisOrb size="avatar" quality="auto" breathing>`（avatar 档＝
  112px 基准 165 粒子，为 44-72px 球窗选的档；容器 100%×100% 自适应任意窗）。
- 宿主缝 `window.jarvisOrb = { setState(name), getState() }`，状态词表＝
  上游自己的四态 `idle|thinking|success|alert`；未知名**忽略不抛**。
- palette 按态取上游内置三色板：idle→cyan、thinking→aurora、alert→ember
  （success 保持 cyan）。
- 音频：上游组件没有音频接缝，入口不发明；Go 侧对该风格跳过 setAudioBands。

## 验证命令（拿到产物后）

```bash
# 静态：外链零命中 + LF
grep -cE 'src="https?://|href="https?://' internal/ball/assets/jarvis-orb.html   # 期望 0
# 动态（真渲染粒子球，弹真窗）：
go test -tags winlive ./internal/ball/ -run TestBallLiveOrbJarvisStyleRendersParticles -count=1 -v
# 或浏览器直接看：任起静态服务，URL 路径指向 jarvis-orb.html
python -m http.server 8923 --directory internal/ball/assets   # http://127.0.0.1:8923/jarvis-orb.html
```
