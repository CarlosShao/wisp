# 票 65 · Jarvis 球深色圆盘底（dark disc）裁定表

- 腿名：65-jarvis-dark-disc（Wisp 编队写腿）
- 日期：2026-11-07 22:1x–23:2x +08
- 派单：编排者（A 腿＝`internal/ball/**` ＋ `tools/jarvis-export/**`；B 腿＝`frontend/**`，互斥）
- 机主原话（逐字）：「这个你说的球风格切换，我不知道怎么弄，要不你直接把 Jarvis 风格搞出来给我看算了，我就不自己搞了」

## 0. 起手三核（现跑，全过）

1. **HEAD 即机主授权 commit**：`git log --oneline -1` ＝ `7873c8328 A1622: OWNER DIRECTIVE x2. (1) PERMANENT AUTHORIZATION: frontend/** and ALL UI now belong to this fleet's write surface…(2) OWNER DOES NOT WANT TO SWITCH MANUALLY - he wants Jarvis simply SHOWN to him…My ruling: give the ball its OWN dark circular backdrop (not a full-window black bg…)…`——派单前提（frontend 归本编队、深色圆盘裁定）与 HEAD 逐字一致，「frontend 归外部 agent Q-51」旧记载确实已废。
2. **引文出处核**：派单点名的四件全部现读——`orb_windows.go`（DefaultBackgroundColor a=0，:336）、`orbasset.go`（注入层 `html,body{background:transparent !important}`，:133；embed 两页）、`renderer_windows.go`（overlayOnly 分支 :473 起）、`tools/jarvis-export/index.html`（body transparent 构造）。
3. **在飞件核查**：起跑时跟踪面零改动（`git status --porcelain` 非 `??` 行＝空）；收尾时出现他腿新在飞件（frontend 5 件＝B 腿、`internal/panel/composer.go`/`pump.go`、票 215）——与 A 腿写面零重叠，本腿 commit pathspec 严格五件点名，未吞任何他腿 hunk。

## 1. 现状复核（机主实测读数全部复现确认）

- 球 100px（`-size` 上限 72＋窗口 edge 100）⇒ 粒子细节糊；
- 透明叠加在白窗上 ⇒ AdditiveBlending 青粒子洗成灰白（`JARVIS-PRODUCT-REAL.png` 缺陷形状）；
- 上游 demo 好看的前提＝大尺寸＋深底（`jarvis-darkbg.png`，深底 `#050913`）。

## 2. 方案裁定：甲（页面侧圆盘）——判据不是感觉

**判据（派单指定）＝哪种能让粒子对比度真的上来且四角仍透明。**

### 乙（D2D 侧画盘）被 z 律结构性否决

`orb_windows.go` follow() 文档 + `newOrbWindow` :398-403 的 z 律：**ball 窗（D2D）在顶层，orb 窗（WebView2 页面）在其正下方**（Interaction 语义全在 ball 窗，画面全在 orb 窗）。renderer 的 overlayOnly 分支画在 **ball 窗 ULW**＝画在粒子**上面**——D2D 圆盘会盖住粒子而不是垫底。乙要成立必须倒 z 序或另开窗口，两样都动冻结交互契约（WindowFromPoint(球心)必须答 ball、四角钉、dock/hit 全挂现几何）。**乙结构性不成立，非代价问题。**

### 甲（页面内 #jarvis-disc）成立

页面是唯一在粒子**下面**的层。圆盘＝圆形、只在圆内不透明、圆外保持页面透明——盘内给粒子对比度，盘外四角仍读桌面。这正是主流 orb UI（Jarvis/Siri/Cortana 类）的「深盘＋亮粒子」形状。

### 几何设计（全部可推导）

| 参数 | 值 | 推导 |
|---|---|---|
| 盘直径 | **窗口边 96%**（半径 48%） | 四角探针在 3px 内缩的对角线上：100px 窗＝66.5% 边长、渐近 70.7%——都在 48% 半径之外，四角读桌面 |
| 粒子最大投影 | 74.2% 半边（＝37.1% 边长） | 相机 PerspectiveCamera(35°) z=4.1、sceneScale 0.82、最宽元素轨道半径 1.17 环；tan(17.5°)×4.1≈1.29 半边可视高，(1.17×0.82)/1.29≈74.2% 半边 |
| 盘能否盖住粒子 | **能**（48% 边长半径 > 37.1% 边长投影） | AdditiveBlending 青粒子从此不再与桌面合成 |
| 盘色 | `radial-gradient(#101a2e 中心 42% 高度 → #0a1322 → #050913)` | 上游 demo 深底族 `#050913`（jarvis-src examples/basic），非纯黑（浅色桌面突兀，派单明令） |

## 3. 落地

- `tools/jarvis-export/index.html`：模板加 `<div id="jarvis-disc">` ＋ CSS（fixed 居中、96%、border-radius 50%、radial-gradient 深底族、z-index:0）；`#root` 提 `position:relative;z-index:1`（React 粒子画布在盘上方）。注释完整记载几何推导。
- 重建 `internal/ball/assets/jarvis-orb.html`（761,765B，纯 LF，外链 0 命中，`jarvis-disc` 2 命中＝CSS＋元素）。
- **基线确定性先证**：改模板前先跑了一次重建，产物 sha 逐字节等于已提交资产（`66eb4b80…`）——构建链可复现，资产 delta 恰为圆盘一项。
- `orb_test.go`：换 sha 钉 `66eb4b80…` → **`2cb15fc2aaa2972181079a01ff889d832b52e8ab4daa288de1ede317bc0c72ff`**（同 commit 换钉纪律）。siri 钉 `d93da649…` 未动。

## 4. 判据（牙）与读数

新静态件 `orb_jarvis_disc_test.go`（2 枚）＋新 winlive 件 `orb_jarvis_disc_live_windows_test.go`（1 枚，白地板三牙）：

### 牙① 盘内亮度必须低于阈值（盘真的存在）

winlive 白地板（0xFFFFFF）＋ 8 点采样环（窗口边 42% 半径——盘内 48% 半径留 ~6px 边距，粒子投影 37.1% 边长之外）：

- **绿**：8/8 非粒子采样，平均亮度 **17.9**（< 60 地板）——首发 PASS 13.91s；
- **突变红**（见 §5）：摘盘后同一采样环 **255.0**。

### 牙② 四角仍在地板色 ±12（透明架构未破）

白地板是 (51,51,50) 灰地板的**相反极端**（机主缺陷形状＝白）：四角 3px 内缩读数 **(255,255,255) 逐通道零偏差**。盘圆够不到对角线探针（66.5% 边长 > 48% 半径）。既有 `TestBallLiveOrbCornersReadTheDesktop`（灰地板 ±12）全套复跑未破。

### 牙③ 粒子对比度：青色主导占比高于加盘前的白底形态

白底（无盘）时青粒子被洗白：`g-r`/`b-r` 收敛 → cyan 主导判据（g-r≥12 ∧ b-r≥12 ∧ b≥80）读数崩；盘上粒子保持青色。读数 **7.08% / 7.49%**（> 2% 地板；与既有 jarvis 各腿 7.08–8.32% 同族——盘没有吃掉粒子亮度）。

### 牙④ 突变（摘盘 ⇒ 牙①必须红）

见 §5，红→还原→复绿闭环。

### 静态牙

`TestJarvisAssetShipsDarkDisc`：embed 内 `id="jarvis-disc"` 存在、`border-radius: 50%`、`#050913` 在盘内、**无平底 `background: #050913`**（平底会破透明架构，防御性负牙）、body 仍 transparent 构造、盘直径 96% 正则钉、文档序盘在 #root 前（垫底层）。`TestJarvisServedDocumentKeepsDiscAndBelt`：注入 belt 与盘在 serving 后共存（belt 管盘外 html,body，盘只管圆内——两层各司其职）。

## 5. 突变④记录（红→逐字节还原→复绿）

1. 备份：`.scratch/wisp/probes/65-disc-mut/{index.html,jarvis-orb.html}.bak`（sha `0f497d04…` / `2cb15fc2…`）。
2. 突变 M1＝模板摘 `#jarvis-disc` CSS 块＋元素 → 重建 → 资产 sha **`66eb4b80…`**＝加盘前已提交资产（delta 恰为圆盘，突变干净）。
3. 红：`TestBallLiveOrbJarvisDarkDiscOverWhiteFloor` **FAIL 94.41s**，3 发重建重试全红，盘环亮度 **255.0**（≥60），非粒子采样 4/8；四角仍 (255,255,255)（盘的缺席不破透明——牙①就是「盘存在」的专属探测器，符合派单牙④预期）；cyan 16.02%（白底上环更宽但亮度判据不背书它好看）。
4. 还原：备份逐字节 `cmp` 双对齐（模板 `0f497d04…`、资产 `2cb15fc2…`，BYTE-EXACT）。
5. 复绿：**PASS 13.08s**（盘环 17.9、cyan 7.49%、四角零偏差）。

## 6. 门禁读数（全绿）

| 门禁 | 读数 |
|---|---|
| gofumpt -l internal/ball（全目录） | 空（rc0），含新两件＋钉件 |
| go vet ./internal/ball | rc0 |
| go build ./internal/ball ./cmd/wisp ./cmd/balldebug | rc0 |
| go test ./internal/ball -count=1 | ok 0.223s |
| **winlive TestBallLiveOrb 全套** | **10/10 PASS 160.176s**（旧 8 枚＋新 1 枚；四角 (51,51,50) ±12 钉未破；新枚 12.85s） |
| d22scan | clean rc0（internal 963 / cmd 239 / frontend 116） |
| LF | 五个触碰件 CR=0 全过 |

## 7. 球放大（-size 上限 72）：**有契约限制，如实报，未硬改**

- 契约＝**44–72px 冻结**：`docs/PLAN.md:3557`（「格子内球体 56px（配置范围 44–72px…）」）、`docs/specs/SPEC-08` §2（配置范围 44–72）、`internal/config/validate.go:62`（`ball.size %d out of range [44, 72]`）、`internal/config/schema.go:179`、`internal/ball/tokens.go:524-528`（`BallSizeMinPx=44`/`MaxPx=72`）＋ `stateSize`/`New` 双处 clamp、tokens_table_test.go 对账。
- 上限不在 `cmd/balldebug`（其 flag 文案引用 ball 常量，:95「44..72」），改它＝动契约面 ＋ `internal/config/**`、`cmd/balldebug/**` 均不在本腿写面。
- **如实结论**：放大到 120–160 需要改契约文字（44–72 → 更宽），按铁律属「摆机主批」的三类之一，本腿不动。当前 100px 窗口内深盘已把对比度缺陷修掉（盘上粒子不再洗白），尺寸是下一个独立改进项。

## 8. 写面与边界

- 触碰件（5 件＋2 文档）：`tools/jarvis-export/index.html`、`internal/ball/assets/jarvis-orb.html`、`internal/ball/orb_test.go`（仅 sha 钉一行＋注释）、`internal/ball/orb_jarvis_disc_test.go`（新）、`internal/ball/orb_jarvis_disc_live_windows_test.go`（新）、本证据件、票 65 Progress log 一行。
- ⛔ 未动：`assets/orb-siri.html`（sha `d93da649…` 钉死未触）、`cmd/wisp/**`、`internal/panel/**`、`frontend/**`（B 腿独占）、`internal/tools/**`、`internal/config/**`、thresholds/golden、票面 AC 勾与 Status。
- 他腿在飞件（frontend 5＋panel 2＋票 215）现读存在，本腿 commit pathspec 严格点名，零吞并。

## 9. 遗留（如实）

1. **UI 代码里零 emoji**：新代码零 emoji，d22scan ban#8 全量 clean 复核。
2. **PrintWindow 截图**：本腿判据全部走 winlive 读屏（GetPixel 合成像素），未产截图——派单明令 ⛔ 靠屏幕截图下结论；机主要看效果直接跑 `build/balldebug-orb.exe -style=jarvis`（盘已在产物里）。
3. `-size` 上限提升＝契约变更，待机主拍板（见 §7）。
