# 票 65 · Jarvis 纯粒子深色腿（65-jarvis-pure-deep-color）

- 日期：2026-11-08（同日第三腿，接 pure-particles 腿被 kill 留下的半成品）
- 机主原话（逐字，最高优先）：**「太淡了，颜色深一点，粒子颜色」**
- 判读：⽠主要的是**更深、更浓、更饱和**（浅色桌面上不要淡成灰白），⛔ 不是更亮更白。
- 角色：实现腿（本档为实现腿自证档；缺口审计/对抗验收归另一 agent）。

## 0. 接手状态（前腿半成品，未丢弃）

- `tools/jarvis-export/vite.config.mjs` 有未提交改动：前腿加了构建期 GLSL 补丁
  （alpha 下限 0.24→0.55、能量因子 0.58→0.72、点径下限 1.08→1.35、twinkle 增益 1.38→1.50、
  饱和 luma-mix 因子 1.4），注释分析扎实；资产本体仍是已推的 754037b9…（759,964B）——前腿没走到重建。
- 本腿在其上继续：保留其全部补丁，叠加本腿改动（见 §2）。

## 1. 基线与「前值」读数（白底 = 机主缺陷形状的极端）

基线先证：`git show HEAD:tools/jarvis-export/vite.config.mjs` 复原配置后 `npm run build`
逐字节重建 ＝ 已提交资产 754037b9a4c393de03800483cba81d718b91dd6b5075ff89b1d421e6d07a5163
（759,964B，`cmp` BYTE_IDENTICAL）⇒ 之后所有 delta 可信。

前值读数（旧资产 + 白底，本机实截 `D:\tmp\wisp\20261108\orblab\leg65-baseline.png`，orb 窗口 100×100 实测）：

| 量 | 前值（754037b9 资产） | 本腿终值（e9026f17 资产） | 判据 |
|---|---|---|---|
| 白底青占样本比 | 13.22% | **39.13%** | >20%（派单①）✓ |
| 白底粒子像素平均亮度 | 174.0 | **115.1** | <120（派单②）✓ |
| 白底粒子色度 mean min(g−r,b−r) | 44.2 | **47.1** | >前值 44.2（派单③）✓ |
| 球心 RGB（5×5 均值） | (134.7,183.0,198.8) lum 172.2（机主截图时代 217） | **(50.0,108.9,156.8) lum 105.2** | 明显低于 217 ✓ |
| 42% 环均值亮度 | 255.0 | 254.7 | >200 无背板 ✓（一点 252.7＝粒子掠过，牙容忍） |
| 黑底青占样本比 | 3.33–3.54%（上腿读数） | **23.41%** | >6%（派单·黑底）✓ |
| 黑底粒子亮度/色度 | — | 64.2 / 62.8，rgb (0,62.8,129.9) | 黑底不发黑看不见：**两全成立** ✓ |
| 黑底四角 | ±0 | **±0（全 0,0,0）** | ±12 保持 ✓ |

白底截图：`D:\tmp\wisp\20261108\orblab\leg65-r13.png`（260×260，含 80px 边距）；
放大交付件 **`D:\tmp\wisp\20261108\orblab\JARVIS-PURE-V3.png`**（orb 实测 100×100（机主机 116 档；本机 DPI 实测 edge=100），5× 最近邻放大到 500×500）。
黑底截图：`leg65-r13-black4.png`。

## 2. 落了什么（全部在 `tools/jarvis-export/vite.config.mjs` 的构建期插件里，上游克隆零改动）

`wisp-jarvis-bright-particles` 插件补丁面（每处站点恰好 1 次、buildEnd 硬护栏，上游变了会 fail build 而不是带病出货）：

1. **粒子 alpha 下限**（前腿 0.55→本腿 0.92+twinkle·0.08；能量因子 0.78 保持前腿抬升）——idle 实际 α≈1。
2. **粒子点径下限**（前腿 1.35→本腿 **6.00**，twinkle 增益 1.90）——见 §3 机理④。
3. **饱和 luma-mix**（前腿 1.4→本腿 **1.9**），且从「只粒子」扩到**全部四个出色站点**：
   粒子 `mix(uPrimary,uAccent,…)`、核心球 `mix(uPrimary,uCore,0.58+…)`、平环 `mix(sideColor,capColor,cap)`、
   数据包 `mix(uCore,uPrimary,…)`。
4. **调色板本体变深（机主令②，前腿没做到、本腿补上）**：调色板源头＝上游克隆
   `D:/tmp/wisp/20261107/orblab/jarvis-src` 的 `src/states.ts`（cyan/aurora/ember 三套），但其发布 dist
   里 hex 已被 esbuild 转成十进制字面量（`0x38f4ff→3732735`）。补丁站点＝**整个 PALETTES 对象字面量**
   （一条精确串、计数护栏）：

   | 字段 | 原值（cyan/aurora/ember） | 新值 | 说明 |
   |---|---|---|---|
   | primary | 0x38f4ff / 0x58f7ff / 0x4df6ff | **0x0088b7**（三套同值） | 粒子主色 |
   | secondary | 0x19a8ff / 0x2cc7ff / 0x0ea5e9 | **0x007d91** | 环线用色 |
   | tertiary | 0x7fffee / 0x9afff4 / 0x99fff0 | **0x0093a8** | 粒子 accent |
   | core | 0xeaffff / 0xf2ffff / 0xffffff | **0x3a8cab / 0x4a9cba / 0x5aa8c6** | fresnel 核心球＝机主看到的「淡」的主体 |
   | deep、fallback | — | 不动 | 已深/渐变串 |

   选值在**合成域**设计：实测场色＝`f·lin(hex)+(1−f)·255`（f≈0.6，见 §3），hex 按线性域目标反推 sRGB。
5. **halo 光晕三站**：上游 halo 渐变**中心是字面白 `rgba(255,255,255,0.92)`**（淡源之一），且把线性值当
   sRGB 画进 SRGB 贴图＝双重伽马压暗。三站全改字面深青 `rgba(0,155,203,…)`（0.92/0.85/0.55）。
   ★如实：r5/r7 两轮读数证明 halo 对被采样场**无可观测贡献**（可能被核心球遮挡或不在采样环带），
   补丁保留（无害、方向正确），但「halo 是主淡源」这一假设**被实测否定**，机理结论以 §3 为准。
6. **avatar 档 dpr 3.5→2.0（minDpr 2.75→1.5）**：渲染质量旋钮（超采样降糊化），轨道几何/粒子数/圈数零改动。

资产重建：双跑确定性 sha **e9026f175a6f5fd6508d537603af7942afb619a732688d6e4d5a03f2aac156f9**（760,176B，纯 LF、外链 0），
同 commit 换钉 `orb_test.go jarvisUpstreamSHA256`（**siri 钉 d93da649… 未动**）。

## 3. 机理（本轮用地面真值钉死的，写下来防下腿重走弯路）

方法：Tabbit 浏览器加载 dist/jarvis-orb.html → `addInitScript` 钩 `shaderSource`（看 GPU 实际编译的源）
＋钩 rAF 第 300 帧 `readPixels` 全画布（看帧缓冲真值），另用 PowerShell 截实机合成屏读数互相印证。

1. **uniform 是线性域**：`new THREE.Color(hex)`（ColorManagement）把 sRGB hex 转 linear-sRGB 再进 shader；
   ShaderMaterial 无输出端 sRGB 重编码 ⇒ **画布上写的就是线性值**（比 hex 观感深）。
2. **合成模型**（白底）：`out = f·lin(场色) + (1−f)·255`，f≈0.6＝该采样像素的覆盖×α；
   黑底：`out = f·lin(场色)`（无白地板）。
3. **饱和补丁是生效的**：帧缓冲满 α 像素实测 (0,88,158)＝恰好 lin(0x0088b7) 过 1.9 饱和的预期值
   （此前几轮「看起来没生效」是覆盖混白造成的读数假象）。突变复证见 §4。
4. **f≈0.6 的墙**：alpha 地板（0.55→0.92）与点径（1.35→3.2）都推不动它；真正推开它的是
   **点径 6.0/7.0**（点在 100px 窗里的纯色核心变大，门内样本从「点边缘混白」变成「点实体」）——
   色度 29.9→48.3、占比 9→39% 的那一步就是它。
5. **色相被黑底门槛钉在蓝青侧（如实取舍）**：黑底牙 `b≥jarvisCyanBright(80)`（与 siri 族共享常量，⛔动）
   在黑底合成（无白地板）里要求 `f·b_lin≥80` ⇒ b_lin≥~130 ⇒ b/g≈1.6-1.8，偏蓝；纯青绿（b_lin≈76）
   会让黑底占比牙变红。当前取点是「白底深度＋黑底可见＋色度牙」三约束的可行点，
   双底读数见 §1；撤销口令＝改 hex 重 build 换钉（一行 sed）。

## 4. 突变证据（派单：饱和因子改回 1.0 ⇒ 亮度/色度判据必须红）

- 突变：四个饱和站点 1.9→1.0，重建（BUILD_RC=0），白底实截读数：
  色度 **40.8 < 44.2 ⇒ 色度牙红 ✓**（帧缓冲纯色 (0,88,158)→(0,71,105)，正如注释预告）。
- ★如实：**亮度牙不红**（s=1.0 实测 lum 100.5，反而更暗）——突变是「更暗但更灰」：
  深来自调色板 hex＋alpha/点径，浓来自饱和；「亮度/色度」这对牙**合起来**才咬得住突变，
  单亮度牙对饱和突变原理上不设防。此如实结论已写进测试件常量注释与 vite.config 注释。
- 还原 1.9 重建：sha 回到 e9026f17…（逐字节复原，`cmp` 过），embed 同步复核 ASSET_IN_SYNC。

## 5. 测试牙（改既有 `orb_jarvis_pure_particles_live_windows_test.go`，新增本地常量，⛔未动共享 jarvisCyan* 与 thresholds/golden）

四牙保持＋新增深色牙（同一批合成屏样本、同一 3px 步进栅格）：
- 白底：占比 > jarvisDeepWhiteRatioFloor=**0.20**；粒子均值亮度 < jarvisDeepLumCeiling=**120**；
  色度均值 > jarvisDeepChromaFloor=**44.2**（＝基线实测前值，出处 leg65-baseline.png）。
- 黑底：占比 > jarvisDeepBlackRatioFloor=**0.06**（亮palette build 实测 3.3-3.5%）。
- 四角 ±12、42% 环 >200、存在性 2% 牙全部保持。`orbPureCyanRatio`→`orbPureParticleStats`
  （多返亮度/色度均值），失败句带全部读数。

## 6. 门禁与状态

- gofumpt `-l internal/ball` 空；`go vet ./internal/ball` rc0；`go build ./internal/ball ./cmd/wisp ./cmd/balldebug` rc0；
  `go test ./internal/ball -count=1` ok。
- winlive `go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -count=1 -timeout 30m`：见 Progress log 行（本档写就时在跑，读数以 log/票行为准）。
- d22scan：见 Progress log 行。
- 改动件全部 LF（`file` 六件无 CRLF 标记）；临时件只建不删（本腿截屏/脚本在 `D:\tmp\leg65\` 与 orblab，仓库外）。
- ⛔ 未碰：`cmd/wisp/**`、`internal/panel/**`、`frontend/**`、`internal/config/**`、`design/**`、
  thresholds/golden、SLO、siri 钉、票面 AC 勾与 Status。只 commit 不 push。
