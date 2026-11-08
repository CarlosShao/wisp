# 票 65 ＋ A1636 批 B：internal/ball 尺寸上限跟随 72 → 160

日期：2026-11-08　腿：Wisp 编队写腿（批 B，`internal/ball/**` 本波独占）
授权：台账 **A1636**，机主逐字「可放大，按你推荐的试试」；批 A（`internal/config/**`＋文档）已由编排者落账 commit `c24e1fa92`，本腿只做批 B，未重做批 A 一字。

## 一、编排者技术裁定的落实（逐条可核）

| 裁定 | 落实 | 证据 |
|---|---|---|
| 新范围 44..160，默认仍 56 | `BallSizeMaxPx = 160`；`BallSizeDefaultPx = 56` 一字未动 | `internal/ball/tokens.go:523-531` |
| `BallSizeLargePx` 必须留在 72（`--ball-lg` 档位语义，上限 ≠ 档位） | `BallSizeLargePx = 72` 一字未动；新牙 `TestBallSizeMaxWidenedTo160` 钉死 72，防有人把档位跟上限一起改 | `internal/ball/tokens.go`＋`internal/ball/tokens_test.go` |
| `BallSizeSmallPx=44` / `BallSizeMinPx=44` 不动 | 未动；新牙一并钉死 56/44/44 | 同上 |
| `RingMarginPx=22` 不动（几何契约） | 未动；`WindowEdgePx(160,96)=160+2×22=204` | `internal/ball/hit.go` |
| clamp 跟随点只读常量、不改逻辑 | `ball_windows.go:161-162`、`statevisual.go:133-134` 零字节动（读常量自动跟随）；逐枚现读确认**无硬编码 72** | git diff |

## 二、改动清单（git 现跑 `git diff --name-only` 全集＝6 枚）

1. `internal/ball/tokens.go` — `BallSizeMaxPx` 72⇒160（`[ball] size_max` 注释同步）；块顶注释 `config range 44-72` ⇒ `44-160`；新增注释一行写明机主授权 A1636 与「档位≠上限」的理由。
2. `internal/ball/ball_windows.go:64` — 仅注释：`SizePx` 字段文档 `44..72` ⇒ `44..160`。
3. `internal/ball/hit.go:30` — 仅注释：RingMarginPx 的波浪适配声明 `44..72` ⇒ `44..160`，并补一句理由（波列顶点 R+21.75px 是半径相对量，+22 余量在任意球径都成立，上限加宽不破坏该几何声明）。
4. `internal/ball/jarvis_windows_test.go:105` — 仅失败文案：硬编码 `the 72px budget` ⇒ `the %vpx budget`（用 `BallSizeMaxPx` 常量）；断言本体零动。
5. `internal/ball/tokens_test.go` — 新增 `TestBallSizeMaxWidenedTo160`：钉 `BallSizeMaxPx==160`、`BallSizeLargePx==72`（档位防连坐）、`WindowEdgePx(160,96)==204`、`stateSize(160,Warm)` 不被钳回 72。
6. `docs/evidence/s1/c21-native-tokens.md:230` — `--ball-size` 行 `可配 44–72` ⇒ `可配 44–160（上限 A1636 放宽）`。

### 第 6 条的越写面说明（如实，非静默）
本腿写面原为 `internal/ball/**`＋本证据件＋票 65 log。但 `TestC21GeometryRowsMatchCodeConstants`（`internal/ball/tokens_table_test.go:1007`）要求**常量值必须是表行内实际写出的数字**：行 230 认领 `BallSizeMaxPx`，行内只写了 56/44/72 ⇒ 常量改 160 后该牙必红。这行在 `docs/evidence/s1/`（与批 A 改 `SPEC-03:27`、`PLAN.md:3557` 完全同性质的机械跟随，非改契约文字），故按铁律⓪⓪⓪（最有价值、不为避风险降级）落了这一字范围跟随，特此具名上报。

## 三、★粒子预算裁定（编排者点名要算的）

公式（`internal/ball/jarvis.go:82-91`，现读）：
```
n = clamp(96 × size/56 + 0.5, 60, 120)   // jarvisParticleCountMin=60, Max=120
```
- `jarvisParticleCount(160)` = clamp(274.3, 60, **120**) = **120** ≤ 128 ⇒ `jarvis_windows_test.go:104-105` 的 128 上限**不挡路、牙不红**，**无需改预算公式、无需改 128**。
- 关键事实：该 120 钳制在 72px 时**早已生效**（96×72/56=123.4→钳 120）⇒ 上限改 160 后 max 档粒子数**一枚不变**，渲染成本零变化。
- 裁定：**两者都不改**（既不重算公式也不动 128），只在失败文案里去掉硬编码 "72px"。

## 四、★116 / 窗口假设排查（逐枚现跑）

- 全仓 `grep BallSizeMaxPx --include=*.go internal/ cmd/`＝19 处，全部落在 `internal/ball`（tokens 定义 1＋clamp 2＋渲染 1＋测试 5＋表 1），**无一处硬编码 72/116/204 以外的窗口假设**。
- `renderer_windows.go:939` `maxWindowEdge(dpi)=WindowEdgePx(BallSizeMaxPx,dpi)` 自动跟随；⚠ 如实报：`maxWindowEdge` **当前零调用者**（本次 grep 全包仅定义行）——历史遗留死代码，不红、未动、只记录。
- `WindowEdgePx` 全部调用点（`ball_windows.go:240/244/746`、`orb_windows.go:262`、测试）都拿实时 size/dpi 现算，无 ≤116 硬算。
- ⚠ 写面外的过期文案（未动，上报）：`cmd/balldebug/main.go:95` 帮助串仍写 `44..72`；`internal/ball/testdata/tokens_c21.css:222` 注释仍写 `可配 44–72`（testdata 参考件，无测试读它的该注释）。

## 五、★实机验收读数（2026-11-08，真桌面）

- 构建：`go build -o build/balldebug-orb.exe ./cmd/balldebug` rc0。
- 起法：`balldebug-orb.exe -style=jarvis -size=160 -state Warm -hold -x 300 -y 300`（起前清残留进程，读数后进程已杀）。
- **窗口矩形（现读 Win32 GetWindowRect）**：`L=300 T=300 R=504 B=504` ⇒ **204×204**，监视器 DPI=96，与 `WindowEdgePx(160,96)=204` 预测逐字节一致。
- 截图：`D:\tmp\wisp\20261108\orblab\JARVIS-160-SHARP.png`（208×208）。
- **球心 RGB＝(43,49,90)**；周点 above(104,84)=(45,48,94) / below(104,124)=(39,43,82) / left(84,104)=(44,47,93) / right(124,104)=(42,48,90)——暗蓝核＋密集粒子＋暖色光带，jarvis 形态肉眼可见（见截图）。
- 球未超出工作区（窗完全在屏内）、未撞任何 116 假设。

## 六、门禁读数（现跑）

| 门禁 | 结果 |
|---|---|
| `gofumpt -l internal/ball` | 空（rc0） |
| `go vet ./internal/ball` | rc0 |
| `go build ./...` | ⚠ **红，但非本腿造成**：`build/mut_rvw.go`（11-04 的历史未跟踪临时件，`git ls-files build/`＝0 枚）是 `package main` 残件，undefined 一串。⛔ 临时件只建不删 ⇒ 改跑等价集合 `go build ./cmd/... ./internal/... ./tools/...` ＝ **rc0** |
| `go test ./internal/ball -count=1` | ok 0.441s（含新牙 `TestBallSizeMaxWidenedTo160` PASS、`TestJarvisParticleCountBand` PASS） |
| `go test -tags winlive -run TestBallLiveOrb` | **10/10 全 PASS，203.981s**（StateAndAudioReachPixels 114.31s / JarvisStyleRendersParticles 7.81s / JarvisServedPageAndOverlay 0.72s / JarvisPureParticles 39.60s / Attach 1.08s / FollowPlacement 0.00s / FollowStaysInside 0.38s / PageRendersPixels 19.64s / CornersReadTheDesktop 19.31s / OverlayNoGlyph 0.80s） |
| `d22scan` | **clean**（从独立 module `tools/d22scan` 内跑 `go run .`；注：从仓库根 `go run ./tools/d22scan` 会报 not contain package——它是自带 go.mod 的独立 module，属跑法问题非违规） |
| LF | 全部改动件 CR=0（逐枚字节现扫） |

## 撤销口令
`git revert` 本腿 commit（单枚）即整批复原；`BallSizeLargePx=72` 与 `RingMarginPx=22` 全程未动。
