# 65 — Jarvis 球「太糊 / 太粗」修复（高清 + 细粒度）

- 日期：2026-11-08
- 腿：65-jarvis-sharp-fine（写腿，只 commit 不 push）
- 机主原话（逐字，最高优先，刚刚下的令）：
  **「不行，太丑了这粒子，1. 太糊了，根本不像人家那个demo里那样高清 2. 粒子有点粗」**
  ⇒ 机主要两件：① 高清（不要糊）② 粒子细（不要粗）。⛔ 不是「再深」，是「清晰 + 细腻」。

## 1. 根因（三重截流，全部现跑定位，非猜）

上一轮深色腿留下的两个补丁方向错了，而且只碰到了表层：

1. **点径粗**＝`vite.config.mjs` 的 point-size 补丁把上游
   `gl_PointSize = uPixelRatio * uPointScale * (1.08 + twinkle * 1.38 + uPulse * 0.48)`
   的下限 **1.08 改成 7.00**（≈6.5 倍）⇒ 细颗粒变团块。这是「粗」的全部原因。
2. **「糊」的真根因不是 preset dpr**——读上游 dist（`jarvis-src/dist/index.js`）发现
   preset 的 `dpr: 3.5` 到达帧缓冲前有**两道隐藏截流**：
   - **quality 档降级**：`resolveQuality` 里 `window.innerWidth < 768` 判 narrow ⇒
     `quality="auto"` 在 100px 球窗里**永远降级 balanced**（avatar `maxDpr: 2.5`、
     `particleScale: 0.86` ⇒ 165 粒被砍成 142）。全屏 demo 不走这条分支——
     机主的参照系（demo）本来就渲染在嵌入窗达不到的分辨率上。
   - **devicePixelRatio 截流**：最终行
     `const renderDpr = Math.max(1, Math.min(window.devicePixelRatio || 1, renderPreset.dpr));`
     把渲染分辨率钉在 WebView2 的显示缩放（典型 1.0–1.5）⇒ 100px 窗只画 100–150px
     的缓冲 ⇒ 物理糊的下限。
   - 上一轮的 `dpr: 3.5→2.0` 补丁动的是**已经被这两道截流钳掉的值**——既不是糊的因，
   也不是糊的药（如实记）。
3. **密度不足**：细点之后白底覆盖塌（389 粒只有 6.52%），按派单正解「细而密」用密度补，
   ⛔ 不靠大点。

## 2. 改动（构建期补丁全部在 repo 侧，上游克隆零改动）

| # | 文件 | 改动 |
|---|---|---|
| 1 | `tools/jarvis-export/vite.config.mjs` | **删** point-size 补丁（1.08 上游字节原样出货）；**删** dpr 2.0 补丁（avatar 档恢复上游 `dpr: 3.5, minDpr: 2.75`）；**新增** `avatar render dpr supersample`（`renderDpr = renderPreset.dpr`，绕开 devicePixelRatio 截流＝demo 在 hi-dpi 屏上的同一条超采样路径，100px 窗画 350×350 再由合成器降采样）；**新增** `avatar particle density`（`particleCount: 165→1400`，×ultra 1.18 = **1652 粒**） |
| 2 | `tools/jarvis-export/entry.tsx` | `quality="auto"` → **`quality="ultra"`**（avatar ultra `maxDpr=4` ⇒ `min(3.5,4)=3.5` preset dpr 存活、粒子数不吃 0.86 折扣；auto 在窄窗必降 balanced 是确定性缺陷，钉死最高档＝确定性） |
| 3 | `internal/ball/assets/jarvis-orb.html` | 重建资产（760,132 B） |
| 4 | `internal/ball/orb_test.go` | sha 钉 `e9026f17…` → `33104399…`（siri 钉未动）＋**新增源码级清晰/细粒度钉**（见 §4） |
| 5 | `internal/ball/orb_jarvis_pure_particles_live_windows_test.go` | 细粒度双牙（新增）＋深度三阈值按新形状重定（见 §4） |

保留未动：alpha 下限 0.92、饱和 1.9、深青调色板、halo 深青站、NormalBlending（十站）、
几何/轨道半径、音频缝、阈值/golden/thresholds.go（⛔ 一字未动）。

## 3. 基线与资产（delta 可信链）

1. **基线重建先证逐字节＝已提交资产**：改前 `npm run build` → sha
   `e9026f17…` ＝ 已提交钉值 ⇒ 管线确定性成立，delta 可信。
2. 密度校准三档（静态白底实截）：389 粒覆盖 6.52% → 1003 粒 12.51% → **定档 1652 粒
   17.35%**（再密黏片：连通域均值面积 3.26→6.76→13.34 px 已在爬，点开始合并）。
3. 最终资产双跑确定性：两次构建 sha 均 `3310439910954d301dfc9c72c55eee17a997ead15ea900068ea17e1ca418a069`；
   落位 `internal/ball/assets/jarvis-orb.html`（LF，CR=0）。
4. 静态尺：`7.00 +` 0 命中、`particleCount:165` 0 命中、`dpr:3.5`×1（avatar 档，
   hero/panel 的 `dpr:2`/`dpr:2.25` 是上游原值）、外链 0、CR=0。

## 4. 判据（新牙 + 按新形状重定，每颗仍咬自己的缺陷）

**粗糊基线 vs 细清版**（静态白底 116×116 实截，同把青色尺，finegrain.py）：

| 指标 | 粗糊版（改前 e9026f17） | 细清版（改后 33104399） | 分离度 |
|---|---|---|---|
| cyan 行程条数 | 320 | 900 | ↑2.8x |
| **平均行程长** | **10.52 px** | **2.59 px** | ↓75%（「粗」被钉死） |
| **连通域个数** | **8** | **175** | ↑22x（「糊成大团」被钉死） |
| **连通域均值面积** | **420.9 px** | **13.3 px** | ↓97%（「细」被钉死） |
| cyan 占比 | 25.02% | 17.35% | 形状改变（见下） |
| cyan 均值亮度 | 114.5 | 122.9 | 形状改变（见下） |
| cyan 均值色度 | 48.7 | 39.7 | 形状改变（见下） |

- **★新增「细粒度」双牙（winlive，1px 步进 20 行采样）**：
  平均青行程长 **< 6.0 px**（粗糊 10.52；新 2.58——退回 7.00 点径即红）＋
  行程条数 **≥ 100**（空场/黏片即红）。
- **★新增「清晰度」源码级钉（非 winlive，普通 `go test` 就咬）**：
  资产必含 `particleCount:1400`、`dpr:3.5`、上游细点径整句
  `gl_PointSize = uPixelRatio * uPointScale * (1.08 + twinkle * 1.38 + uPulse * 0.48)`；
  必不含 `7.00 + twinkle`、`particleCount:165`（GLSL 是字符串字面量 minify 不动、
  对象字面量属性名不混淆 ⇒ 精确且确定）。
- **深度三阈值按新形状重定（如实，非放宽凑绿）**：细颗粒场**本质是边缘主导**——
  ~2px 的点几乎全是部分 alpha 与白的混色像素，同色下均值亮度必高于团块场（内部主导）、
  均值色度必低于团块场。旧数（lum<120、chroma>44.2、share>20%）是团块形状的性质，
  不是质量线。重定后仍各咬各的缺陷：
  - share > **14%**（新 17.35；pale 旧版 12.28–13.63% 在地板下；粒子丢失塌向 0）；
  - lum < **135**（新 121.2–122.9；pale 旧版 174.0；alpha 突变 0.92→0.24 把细点混色像素推到 ~230）；
  - chroma > **34**（新 39.7–40.6；饱和突变 1.9→1.0 使纯色 g-r 降 ~19% 落 ~32 即红）。
- **保留既有**：白底无盘（42% 环 >200）、四角双底 ±12、黑底可见（>6%）。

## 5. 读数（★winlive 全绿实测，TestBallLiveOrb 全家 194.700s PASS）

纯粒子牙 `-v` 单跑（attempt 1 一次过）：

- **白底**：四角 (255,255,255)×4 零偏差；42% 环 **8/8 非粒子、均值亮度 255.0**（>200 无盘）；
  cyan 占比 **26.01%**（存在 >2%、深 >14%）；粒子均值亮度 **121.2**（<135）、
  色度 **40.6**（>34）；**细粒度：206 条行程（≥100）、均值行程 2.58px（<6.0）**。
- **黑底**：四角 ±12 内；cyan 占比 **12.90%**（>6%）、lum 62.0、chroma 60.4。
- 门禁：`gofumpt -l internal/ball` 空；`go vet ./internal/ball` rc0；
  `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` rc0；`go test ./internal/ball -count=1` ok；
  ★winlive `TestBallLiveOrb` 全家 PASS 194.700s；`d22scan` clean（internal 964 + cmd 243 全量）；
  LF：五件触改文件 CR=0。

## 6. 实机验收（balldebug，真渲染真窗）

- `go build -o build/balldebug-orb.exe ./cmd/balldebug`；清干净残留；`-style=jarvis` 起；
  截 orb 窗 116×116 放大 4x → `D:\tmp\wisp\20261108\orblab\JARVIS-V4-SHARP.png`
  （对照粗糊版 `JARVIS-COARSE-BEFORE.png`，同机位同尺）。
- 球心 RGB **(122,147,162)**（粗糊版 (73,120,151)/改前深色版实测中心更糊灰）；
  画面中**独立细颗粒连通域 175 个**（均值 13.3px，粗糊版 8 个大团块/均值 420.9px）——
  目测即 demo 的细颗粒云。

## 7. 边界与如实

- ★**窗口 100px 是糊的物理下限之一，具名在此**：上游 demo 是**全屏渲染**（640px hero 档，
  devicePixelRatio 不被窄窗截流），本球的窗只有 100px。本轮已用「钉 ultra 档 + renderDpr
  强制 3.5 超采样」把**每物理像素拿到的渲染信息量**抬到 demo 同级（350×350 缓冲降采样到
  100px，等效 3.5x SSAA）；但**画幅本身**（粒子云的可辨细节总量）仍受 100px 窗约束——
  同样的粒子云在 200px 窗里会明显更像 demo。**若机主判断「必须放大球窗才真正像 demo」，
  那是契约面（球窗尺寸属 D 表/规格域），本腿不擅动，另裁。**
- `quality="ultra"` 的代价：渲染缓冲 350×350（~0.12MP）＋1652 粒子——对 WebView2 已在跑的
  GPU 是零头；60fps 不变。
- dpr 2.0 旧补丁的定性如实：它动的是被 quality/devicePixelRatio 双重截流钳掉的值，
  既非糊因亦非糊药；本轮删除。
- 深度三阈值重定的性质如实：是**形状再校准**（细颗粒边缘主导的物理），不是为绿放宽——
  每颗牙对各自的缺陷（pale 174.0、alpha 突变 ~230、饱和突变 ~32、粗糊 10.52px）
  保持红绿分离，出处读数全在 §4。
- winlive 采样在第二屏（orb 窗实位 (3316,1268)）；静态对照与 winlive 同窗位。
- 上游克隆零改动；`tools/jarvis-export/.verify.js` 为前腿临时件，只建不删。

## 8. 撤销口令

`vite.config.mjs` 删去 `avatar render dpr supersample`/`avatar particle density` 两补丁、
恢复 point-size 7.00 补丁与 dpr 2.0 补丁、`entry.tsx` 的 quality 改回 "auto"、重建资产、
sha 钉换回 `e9026f17…`、测试件深度阈值还原（120/44.2/0.20）并删细粒度牙与源码钉——
git 历史 `33104399` 之前的 `e9026f17` 提交有全量原文，一句话即可整体复原。
