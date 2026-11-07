# 68 — parity closure：既有红 `TestBallLiveEdgeDock`（boundary −4）收件

- **日期**：2026-11-07 20:3x–20:5x +08
- **腿**：wave75 / 68-parity（`internal/ball/**` 独占写面）
- **锚点**：HEAD `7202103f3`（A1617，wave74 已收）开工；本件交付于其上
- **票**：`.scratch/wisp/issues/68-ball-default-visuals-parity.md`

## 1. 起手三核（现跑，非采信票面）

| 核 | 读数 |
|---|---|
| ①票面现量 | 票面正文停在 2026-09-20（AC#2 未翻默认值）；**码里已过期**：`statevisual.go:110` `var prototypeVisuals = true`、`tokens_test.go:159` 断言默认 ON ⇒ AC#2 翻转与 W68b（settling 0.95 / dock 42% / 断言迁移）已在早前波次落地（台账 `pending-and-issues.md:31432`）。票面**未过期**且仍欠的现量＝winlive 既有红。 |
| ②引文出处 | 台账 `:31496` 与 `:31786` 逐字核过：红句 `live_windows_test.go:384: edge 1: docked orb boundary -4, want the work edge 0`，归「68-r2 改泊靠几何后未跑 winlive 的欠账」。现跑复现一致（见 §2）。 |
| ③红绿归属 | `git status --porcelain` 过滤 `^??` 后**零已跟踪改动** ⇒ 无在飞写腿，红不是别人造成的。 |

## 2. 既有红今日现跑（开工时）

```
=== RUN   TestBallLiveEdgeDock
    live_windows_test.go:384: edge 1: docked orb boundary -4, want the work edge 0
--- FAIL: TestBallLiveEdgeDock (1.26s)
=== RUN   TestBallLiveEdgeDockHover
--- PASS: TestBallLiveEdgeDockHover (0.47s)
FAIL  github.com/CarlosShao/wisp/internal/ball  1.962s
```
（`TestBallLiveEdgeDockHover` 今天本来就绿；红的只有 `TestBallLiveEdgeDock` 本体。）

## 3. 根因（码读，非转述）

| 侧 | 半径 | 泊靠 p=1 时沿轴出画量 |
|---|---|---|
| 位置法（码，已对） | `dock_windows.go:69` `orbR = dockOrbRadius(b.curVisual.SizePx)`＝态半径 **17**（Sleeping 34.72/2）；`dock.go:153` `off = int(orbR × DockSquash(1))`＝`int(17×0.42)`＝**7** | 7px |
| live 测试（旧，过期） | `live_windows_test.go:369`（旧）`R = edgePx/2 − RingMarginPx`＝窗口半径 **28**；`drawn = int(28×0.42)`＝**11** | 11px |

`11 − 7 = 4` ⇒ 测试按旧模型算出的球界比法线所在**凹进 4px**；主屏 work 边在 0 时 `got` 字面就是 **−4**。位置法本体（R15#4/Q-4 的「喂画出的本体半径」）是对的，欠的只是 live 测试的几何模型还停在 68-r2 之前的窗口半径口径。纯面牙早已在场：`dock_test.go:189`、`tokens_test.go:319` 都钉了 `dockOrbRadius(34.72)=17`。

## 4. 修法（测试模型迁移，零断言放宽、零行为改动）

`internal/ball/live_windows_test.go` `TestBallLiveEdgeDock`：
- 期望出画量改为与位置法同源：`orbR := dockOrbRadius(frameVisual().SizePx)`（UI 线程读，Sleeping 全程不变，读一次），`drawn = int(orbR × DockSquash(p))`——与 `dock.go:153` 逐表达式相同（连截断都一致，不是只靠 ±2 容差兜住）。
- 弹出后「全 orb 在工作区内」检查的半径同步改用 `orbR`（p=0 时窗口按 `DockPos` 本法横骑边界只露透明边距，orb 本体必须完整可达——注释写明）。
- 断言强度不变：四边相切判据、`±2` 容差、底部任务栏门、ramp 中间帧 ≥2、零定时器门全部原样保留。
- 删掉的 `marginPx` 局部变量只服务于旧 R，别无消费者。

## 5. 牙（突变后必须红）

突变：`dock_windows.go:69` → `orbR = edgePx/2 - marginPx`（窗口半径，pre-R15#4 旧错法）：

```
--- FAIL: TestBallLiveEdgeDock (1.09s)
    live_windows_test.go:396: edge 1: docked orb boundary 4, want the work edge 0
```

镜像方向 +4，同一枚 11−7=4 ⇒ **牙有**。还原后 `sha256(dock_windows.go) = be0e093f9a10ab55110707ba3c2dc0fe0157381cb86fcdbc7d621127a5c62e38`＝HEAD 版逐字节相同（开工时留底 `.scratch/wisp/probes/68parity/dock_windows.pristine.sha`），`git status` 对该文件零输出。日志：`.scratch/wisp/probes/68parity/mut-window-radius-red.txt`。

## 6. 收件读数表（全部现跑）

| 门禁 | 读数 |
|---|---|
| `gofumpt -l internal/ball` | 空，rc0 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1` | `ok 0.384s` |
| `go test -tags winlive ./internal/ball/ -run TestBallLiveOrb -count=1 -timeout 25m` | **8/8 PASS**（JarvisStyleRendersParticles 9.08s / JarvisServedPageAndOverlay 0.82s / AttachStateAndTeardown 0.99s / FollowPlacementDecision 0.00s / FollowStaysInsideWorkArea 0.25s / PageRendersPixels 20.73s / CornersReadTheDesktop 21.22s / OverlayFrameHasNoGlyphPixels 0.98s；`ok 54.353s`） |
| `go test -tags winlive ./internal/ball/ -run TestBallLiveEdgeDock -count=1 -v` | **PASS 1.47s** ＋ Hover PASS 0.55s（收件前红 1.26s） |
| `tools/d22scan`（own module，`-root ../../`） | `clean - no D22 ban violations`（internal/=956、cmd/=236、frontend/=116、design/=39） |
| LF | `git ls-files --eol`：`live_windows_test.go` 与 `dock_windows.go` 均 `i/lf w/lf` |

## 7. 边界如实

- **零生产行为改动**：`dock.go`/`dock_windows.go`/`statevisual.go` 等**一字未改**（dock_windows.go 经突变-还原后 sha 与 HEAD 相同）；本件只改一个测试文件的几何模型。
- 票 68 的 **AC#2/AC#3 勾与 Status 未动**（需真桌面复测/机主签收，归裁决腿与机主；本腿 ⛔ 翻勾）。
- `assets/orb-siri.html`、`assets/jarvis-orb.html`、`cmd/wisp/**`、`internal/panel/**`、`internal/tools/**`、`internal/config/**`、`frontend/**`、`design/**` 零触碰。
- 票面 Progress log 追加一行（append-only）。
