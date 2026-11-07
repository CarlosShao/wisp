# 65 — 球窗可点性回归修复（wave59 实现腿）

- 日期：2026-11-07
- 角色：实现写腿（裁决腿 e5ab45154 定性＋编排者派单；本腿开方实施）
- 基线：`bee208958`（A1600）；修复基点＝该 HEAD 未提交工作树
- 写面：`internal/ball/ball_windows.go`（2 处守卫）、`internal/ball/orb_windows.go`（样式＋注释＋删守卫函数）、本件、票 65 Progress log 一行

## 1. 红句复现（本腿亲跑，与派单读数同形）

```
go test -tags winlive ./internal/ball/ -run 'TestLiveNeverStealsFocus|TestLiveTransparentCornerFallsThrough' -count=1 -v
→ interaction_live_test.go:343: the pixel under the pointer is not the ball after moving there: hwnd=8ee06f2 want 17d0f80
→ interaction_live_test.go:371: WindowFromPoint(centre) = 8ee06f2, want the ball 17f0f80 - the orb must be clickable
→ 两枚 FAIL (1.333s)
```

## 2. 诊断：因＝球窗不画，果＝隐形＋穿透（不是 orb 挡球）

- 几何事实：orb 窗与球窗是**并排**的两个顶层窗（`orbFollowXY`：球右 12px gap，放不下翻左/夹紧），
  不是父子，也从不叠在球窗上 ⇒ "orb 挡住球"不成立。
- 真因果链：`orbLive()` 守卫（原 ball_windows.go :419 renderFrame、:554 present）⇒ orb 存活时球窗
  **从未 ULW 一帧** ⇒ WS_EX_LAYERED 窗 alpha 恒 0 ⇒ OS hit-test 对全透明分层窗跳过
  （WindowFromPoint 落到下层第三方窗，裁决腿探针同证）⇒ 球隐形＋点击穿透。
- 判别证据：**守卫一恢复（实验 B），两枚红测试即转绿**（见 §5），因果链闭环。

## 3. 技术选型：甲（球窗照常画＋orb 纯装饰穿透层），乙丙不选

- **乙否**：与已绿测试正面冲突——`orb_live_windows_test.go:152` 钉死 "the orb is a SECOND window,
  not a replacement"；出屏夹紧测试钉死并排几何。乙还要改两枚 winlive 断言对象，违背判据 3 的
  "断言对象不该改"。
- **丙否**：丙与甲在本几何下等价（orb 反正不叠球窗），甲比丙多走的只是把选型语义写明——
  零额外成本，按"写清楚"处置。
- **甲裁**：主流 Windows overlay 形状＝装饰层（LAYERED|TRANSPARENT|TOPMOST|NOACTIVATE，全穿透）
  ＋交互实窗（承载全部命中/拖拽/焦点语义）。原守卫的前提"painting the D2D frame into the layered
  window behind a transparent child would double-composite"**前提即错**——orb 不是球窗的 child，
  是独立顶层窗，两窗各自走 DWM 合成，不存在 double-composite。甲恢复球窗全部冻结契约
  （drag/dock/WM_NCHITTEST 圆内可点/透明角穿透/NOACTIVATE），orb 降为纯装饰，**测试断言零改动**。

## 4. 改动清单

1. `ball_windows.go renderFrame()`：删 orbLive 早退守卫＋重写注释（为什么不画是错的）。
2. `ball_windows.go present()`：同上删守卫。
3. `orb_windows.go`：`orbWindowExStyle` 增 `WS_EX_TRANSPARENT`（0x20，本文件定义常量）；
   `orbLive()` 函数删除（调用点已不存在）；文件头注释增 Clickability 段记录裁定。

## 5. 判别实验（pixel 枚摇摆取证，全程如实）

修复态下 `TestBallLiveOrbPageRendersPixels` 出现过 in-suite 三连红（0.000% 全黑、best 0.000%、
25s 超时），但同码 standalone 四连绿、其后全套四枚又连绿。判别链：

| 实验 | 内容 | 结果 |
|---|---|---|
| 修复态 standalone | Pixel 单跑 ×4 | 绿：8.50% / 38.27% / 35.37% / 36.51% |
| 修复态 in-suite（早期） | 4 枚整套 ×3 | Pixel 红：0.000% ×3 |
| 基线（仓外 clone bee208958）in-suite | 4 枚整套 | 全绿：28.74% |
| 实验A（修复态只撤 TRANSPARENT） | 4 枚整套 | Pixel 仍红 0.000% ⇒ TRANSPARENT 排除 |
| 实验B（守卫恢复＝基线画帧行为） | 4 枚整套 | Pixel 绿 31.97% ⇒ 红因锁定"球窗照画"与 in-suite 状态相互作用 |
| P1（Attach+Pixel） | 修复态两枚 | 绿 19.10% ⇒ Attach 单独前置不武装 |
| P2（Attach+Follow+Pixel） | 修复态三枚 | 绿 12.93% |
| P3（Placement+Pixel） | 修复态两枚 | 绿 22.05% ⇒ 单独前置也不武装 |
| 修复态 in-suite（后期） | 4 枚整套 ×2 | 全绿：19.95% / 12.47% |

判定：Pixel 枚对本机 WebView2 冷启动/首帧时机敏感，in-suite 早段出现过 0.000% 全黑但**同一二进制
后续多轮 in-suite 连绿**；基线与修复态都有通过样本。本腿未改该测试一字、未动 thresholds/golden。
两枚修复目标测试（:343/:371 断言对象＝球窗）在修复态下 5 轮全绿，未借任何断言改动。

## 6. 门禁读数（全部亲跑）

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空，rc0 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1 -timeout 20m` | ok 0.379s |
| ★`go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb'` | 修复态末两轮 4/4 PASS（19.95%、12.47%） |
| ★`go test -tags winlive ./internal/ball/ -run 'TestLiveNeverStealsFocus|TestLiveTransparentCornerFallsThrough'` | 2/2 PASS（1.47s/0.25s；修复后共 5 轮绿） |
| `sh scripts/d22scan.sh` | clean（internal/=941, cmd/=226 Go files） |

## 7. 如实段（未做到/边界）

- `orbdebug.go:36` 注释 "whether the WebView2 orb owns the pixels" 语义已过时，但该文件在授权写面
  之外，未动。
- 共享工作树里有他人未提交改动（`cmd/wisp/approval_reply.go` 等），本腿零触碰，commit 只点
  自己的 pathspec。
- Pixel 枚的摇摆未根除（时序敏感，与本修复无关的既有形状）；若编排者需要零摇摆，须另立票对
  该测试做启动确定性治理。
