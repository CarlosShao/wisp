# 票 65 · orb follow 出屏修复（wave55 `orb-fix-offscreen`）

- 日期：2026-11-07 13:1x +08 · 腿：orb-fix-offscreen · 触发：台账 **A1596**（编排者亲跑实测）
- 实测凭据（编排者，`build/balldebug-orb.exe` + EnumWindows）：「Wisp Orb」rect=(3428,1268)-(3560,1400)、虚拟屏宽 3440 ⇒ **orb 右边缘出屏 120px**；根因＝`follow()`（orb_windows.go :311-324）无条件放球右侧 `wr.r+gap`。

## 修了什么

`internal/ball/orb_windows.go`：`follow()` 拆出纯函数 `orbFollowXY`（球 rect、orb 边长、gap、工作区 `MonitorRect.Work`、是否有显示器信息）：

1. `MonitorFromWindow(MONITOR_DEFAULTTONEAREST)` + `GetMonitorInfoW` 取 **rcWork**（复用既有 `monitorRectsForWindow`，避任务栏；零新 win32 绑定）；
2. 右侧放得下（`wr.r+gap+edge <= rcWork.right`）→ 保持右侧（原行为）；
3. 放不下且左侧放得下（`wr.l-gap-edge >= rcWork.left`）→ 翻到球左侧；
4. 两侧都放不下 → 水平夹紧进工作区；垂直同夹紧到 `rcWork.top..bottom-edge`；
5. 取不到显示器信息 → 原样回退旧形状（右侧），不猜边界。

WS_EX（LAYERED|TOPMOST|TOOLWINDOW|NOACTIVATE）、`HWND_TOPMOST`、`SWP_NOACTIVATE` 语义逐字未动。

## 有牙的回归（`internal/ball/orb_live_windows_test.go`）

- `TestBallLiveOrbFollowPlacementDecision`：`orbFollowXY` 六形状表测（右有位保持/右缘停靠翻左/底缘垂直夹紧/球越顶夹下/窄屏两侧都放不下夹紧/无显示器信息回退）。
- `TestBallLiveOrbFollowStaysInsideWorkArea`（真窗）：球 `SetWindowPos` 到工作区**右下角贴边**，调 `followOrb()`，断言 orb 真实 rect **完全落在 rcWork 内**（left≥work.L 且 right≤work.R、top≥work.T 且 bottom≤work.B），并断言精确坐标（`orb.r == ball.l-12` 翻左、`orb.b == work.B-edge` 垂直夹紧）。修前形状（`wr.r+12`）在此摆位下必红。

## 门禁读数（本腿自跑）

| 门禁 | 读数 |
|---|---|
| `gofumpt -l`（orb_windows.go + 测试） | 空 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1` | ok |
| `go test -tags winlive ./internal/ball/ -run TestBallLiveOrb` | ok |
| `sh scripts/d22scan.sh` | clean |

（最终读数以腿收件报为准，此件随 commit 落盘。）

## 边界（如实）

- 窄屏"两侧都放不下"分支在真显示器上无法物理摆出，由纯函数表测覆盖，winlive 只测真窗可达形状。
- D2D fallback、上游 `assets/orb-siri.html`、阈值/golden、禁区（cmd/wisp、panel、winsec、agent、frontend/design）零改动。
- 裁决（审计者≠实现者）未做，归编排者。
