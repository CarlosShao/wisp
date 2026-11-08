# 65 — Jarvis 纯粒子态（删盘 + NormalBlending）

- 日期：2026-11-08
- 腿：65-jarvis-pure-particles（写腿，只 commit 不 push）
- 机主原话（逐字，最高优先，推翻 2026-11-07 深盘裁定）：
  **「这个球有个问题，就是不要有这种球体，就直接像Jarvis那样，那没有明显的边框，背景，直接就是粒子态，你能懂我意思吗」**
  ⇒ 纯粒子态：不要深色圆盘、不要可见边框、不要背景；粒子自己悬空。

## 1. 裁定与理由（技术选型自行拍板，依据铁律 ⓪⓪⓪）

- **病根**：上游 `cyber1443/jarvis-ai-orb-web-animation` 的 renderer 把全部 **10 处**
  粒子材质钉成 `THREE.AdditiveBlending`（`src/renderer.ts` :95/:140/:190/:246/:562/:577/:600/:619/:637/:687）。
  Additive 把颜色**加到背景上**：白面上低 alpha 粒子边缘被 clamp 向白（机主 JARVIS-PRODUCT-REAL.png 的灰白环缺陷）。
- **上一轮的圆盘方案**（`#jarvis-disc` 深色底盘）被机主本人否掉 —— 盘本身就是他不要的「球体」。
- **正解＝换混合模式，不加任何背景**：`THREE.NormalBlending`（three.js 透明材质的**默认**混合模式）——
  每颗粒子按**自身 alpha** 与背后任何内容做正确 over-合成：白底上是青、暗底上也是青，无需底盘。
  这就是主流 three.js 粒子系统在**未知背景**上保证可见性的标准形状（normal blending + 高饱和色板），
  而非 additive + 深底盘。
- **落点＝构建期补丁（repo 自有），上游克隆零改动**：`tools/jarvis-export/vite.config.mjs` 新增
  `wisp-jarvis-normal-blending` vite 插件，对上游 dist 做 `blending: AdditiveBlending` → `blending: NormalBlending`
  精确重写，**硬护栏＝恰好 10 处，数目不符即 fail build**（上游更新静默漏补丁整类风险封死）。
  干净机器复跑路径（README：fresh clone + npm install + npm run build）**不变、可复现**；
  产物 sha 钉钉的是补丁后的构建。补丁包（patch-package）同思路，零新增依赖。

## 2. 基线与资产（delta 可信链）

1. **基线重建先证逐字节＝已提交资产**：改前 `npm run build` → `cmp` 逐字节相等，sha
   `2cb15fc2aaa2972181079a01ff889d832b52e8ab4daa288de1ede317bc0c72ff` ＝ 上一轮钉值 ⇒ 管线确定性成立，delta 可信。
2. 模板 `tools/jarvis-export/index.html`：删 `#jarvis-disc` 整段（CSS 规则 + `<div>` 元素 + 旧注释），改为纯粒子裁定注释；
   `body { background: transparent; }` 构造即透明**保持不动**。
3. 重建 + 双跑确定性验证：两次构建 sha 均为
   `754037b9a4c393de03800483cba81d718b91dd6b5075ff89b1d421e6d07a5163`；
   落位 `internal/ball/assets/jarvis-orb.html`（759,964 B，`cmp` 一致）。
4. 静态尺：`jarvis-disc` 0 命中、`#050913` 0 命中、CR=0（纯 LF）、外链 `src/href="http…"` 0 命中。
5. `orb_test.go` 的 `jarvisUpstreamSHA256` 同 commit 换为 `754037b9…`；**siri 钉 `d93da649…` 未动**。

## 3. 判据（判据翻转：上一轮证明「盘在」，本轮证明「盘没了 + 粒子任何底都在」）

测试件重命名（如实报）：
- `orb_jarvis_disc_test.go` → **`orb_jarvis_pure_particles_test.go`**（git rm + 新建）
- `orb_jarvis_disc_live_windows_test.go` → **`orb_jarvis_pure_particles_live_windows_test.go`**（git rm + 新建）

新牙（winlive `TestBallLiveOrbJarvisPureParticlesOverWhiteAndBlackFloors`，白/黑双底一球跑完）：

| # | 牙 | 判据 | 实测读数 |
|---|---|---|---|
| ① | **无盘牙**（翻转上一轮） | 42% 边长 8 点环的非粒子样本 ≥4/8，平均亮度 **>200** | **8/8，均值 255.0**（盘在时同环读 ~17.9；分离度干净） |
| ② | 四角透明（双底） | 四角 3px 内逐通道 ±12 读底色 | 白底 4×(255,255,255) 零偏差；黑底同过 |
| ③ | **白底青粒子**（本腿核心） | 青主导占比 (g-r≥12 ∧ b-r≥12 ∧ b≥80) > 2% | **12.59%**（此前盘上 7.08–7.49%；白底上青粒子可见性成立） |
| ④ | **黑底青粒子**（双底判据） | 同③，换黑底再测 | **3.33%** > 2% |

★**如实（突变验证 ④ 的真实结果）**：把 `jarvisBlendingPatch` 改回 additive、重建（sha
`6c099e37…`）、重跑 ⇒ **牙没有变红**：白底青 13.63%、黑底 4.89%，与 normal 构建（12.28%/3.54%）同区间，
四牙全绿。**机理**：非重叠粒子像素在透明画布上两种混合模式的缓冲区**逐位相同**（rgb=cyan·a, alpha=a，
清屏为透明黑），重叠区占比小 ⇒ 平底色相占比尺**原理上分不开**两种混合模式（黑底 3.54 vs 4.89 的 1.35pp 差
与历史批间方差 ~1.2pp 同量级，不可作阈）。⇒ 现有四牙咬的是**盘回归**（牙①：盘在=17.9 vs 无盘=255.0）、
**整窗背景回归**（牙②）、**粒子丢失/洗白**（牙③④），**不咬混合模式常量**；NormalBlending 的落位依据是
机主裁定（盘被否）＋ three.js 默认实践，已在测试件头注与 vite.config 注释里如实记下本轮突变实测数。

## 4. 门禁（全绿）

- `gofumpt -l internal/ball`：空（新件一处对齐由 gofumpt -w 归位后复检空）
- `go vet ./internal/ball` rc0；`go build ./internal/ball ./cmd/wisp ./cmd/balldebug` rc0
- `go test ./internal/ball -count=1` ok
- ★winlive `TestBallLiveOrb` 全家 **10/10 PASS 165.313s**（旧 8 枚 + jarvis 风格 + 新纯粒子牙；四角 (51,51,50) 钉未破）
- `d22scan` clean（internal 963 + cmd 243 等全量）
- LF：六件触改文件 CR=0

## 5. 边界与如实

- 桌面上残留 `wisp.exe`（pid 41396）挡测量，按牙自带指令 taskkill 后重跑（非本腿进程，遗留 balldebug/wisp）。
- 上游克隆（`D:/tmp/wisp/20261107/orblab/jarvis-src`，node_modules 以 symlink 指向）**零改动**；
  补丁全部在 repo 侧 `vite.config.mjs`，干净机器路径不变。
- 粒子数量/轨道几何零改动（既有判据未动）；音频缝零改动。
- 白底青占比 12.59% 的机理说明：白底上青色相门限 a≥6.4% 即过（g-r=188a），淡晕也计入 ⇒ 数值高于深底盘时代属预期。
- 混合模式无平底像素尺可钉（见 §3 突变段）——若机主要「常量级」的钉，需要重叠区专项仪器（当前判为低价值高脆性，未上）。

## 6. 撤销口令

`vite.config.mjs` 里 `jarvisBlendingPatch` 改回 `"blending: AdditiveBlending"`、模板还原 `#jarvis-disc` 段、重建资产、
sha 钉换回 `2cb15fc2…`、恢复两枚 disc 测试件（git 历史 `24167a4cc` 有全量原文）——一句话即可整体复原。
