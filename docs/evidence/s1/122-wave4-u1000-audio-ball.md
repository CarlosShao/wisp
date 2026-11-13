# 票 122 · 波 4b（重派）：`internal/audio` + `internal/ball` 的 U1000 清账

- 时间：2026-11-13 14:34–15:1x +08（写腿）
- 尺：`GOOS=windows CGO_ENABLED=0 staticcheck -checks=all,-SA9009 ./internal/audio/ ./internal/ball/`（不带 tag 与 `-tags=winlive` 各一组，读数文件在 `D:\tmp\wisp\20261113\122-wave4b\`）
- 派单：`.scratch/wisp/probes/orch/20261113/txt/122-wave4b-u1000-dispatch.txt`
- 铁律 0 全程执行：**U1000＝"没有 Go 代码读它"，不是"没有东西读它"**。四档处置逐条落账见下表。

## 1. 前/后读数

| 尺 | 前 | 后 | 残留 |
|---|---|---|---|
| 不带 tag | **40 条 U1000** | **2 条 U1000** | `orbHwndTopmost`＋`overlayGlyphPixels`（均为④，派单禁删禁 ignore，见 §3） |
| `-tags=winlive` | 43 条 U1000（40＋3 条**假阳性**，见 §4） | **2 条 U1000**（与不带 tag 同两条） | 同上 |

- ① ignore ＝27 条；② ignore ＝1 条；③ 删 ＝12 条目（10 符号＋capW/capH 2 枚字段＋winlive 修复后新显形 2 枚另计，见 §4）；④ 点名保留 ＝2 条。
- 其余残留读数为 ST1003×7（mmdevice `deviceId` 参数名、wasapi `pGetCurrentThreadId`、两枚测试 error var 名）——**存量命名类，非本票射程**（本票只清 U1000），如实记录不动。

## 2. 40 行逐条表

| # | 符号 | 声明位置 | 档 | 处置与理由 |
|---|---|---|---|---|
| 1 | `iidIUnknown` | `internal/audio/wasapi_windows.go:49` | ③删 | 全仓 grep 零引用；它是纯 Go 值变量（GUID 常量），从不作为内存块被 API 按偏移读——COM 调用传的是 `&iidIMMDeviceEnumerator`/`&iidIAudioClient` 等**别的** GUID 的地址 |
| 2 | `fadeStart` | `internal/ball/ball_windows.go:118` | ③删 | 全仓零读零写（批 1/批 2 判例一致：写侧已不存在，非 OS 镜像字段） |
| 3 | `startedErr` | `internal/ball/ball_windows.go:160` | ③删 | 同上，全仓零读零写 |
| 4 | `(*Ball).recenterAt` | `internal/ball/ball_windows.go:1113` | ③删 | 派单点名③；全仓零调用；函数体＝`pGetWindowRect`＋`moveWindow` 纯调用、无需要长期存活的副作用资源；`pGetWindowRect` 在 ball_windows/dock_windows 另有 7+ 活调用者，不级联 |
| 5 | `d2d1RenderTargetProps.typ` | `internal/ball/d2d_windows.go:139` | ①ignore | **D2D1_RENDER_TARGET_PROPERTIES ABI**：`CreateDCRenderTarget` 按内存偏移读整个结构体 |
| 6 | `d2d1RenderTargetProps.usage` | `internal/ball/d2d_windows.go:144` | ①ignore | 同上 |
| 7 | `d2d1RenderTargetProps.minLevel` | `internal/ball/d2d_windows.go:145` | ①ignore | 同上 |
| 8 | `dwmBlurBehind.fTransitionOnMaximized` | `internal/ball/orb_windows.go:165` | ①ignore | **DWM_BLURBEHIND ABI**：`DwmEnableBlurBehindWindow` 按偏移读 |
| 9 | `orbHwndTopmost` | `internal/ball/orb_windows.go:193` | **④留** | winlive 测试真调用：`orb_live_windows_test.go:559`/`:577`、`orb_glass62_live_windows_test.go:145`/`:160`（四行文件头 build tag 均 `//go:build windows && winlive`），作 `SetWindowPos` 的 hWndInsertAfter；⛔删（必红 winlive 编译）⛔ignore（会把真调用者藏起来）；不带 tag 读数里的残留属预期，已在源码处留 Ticket-122 注释点名调用者 |
| 10 | `orbDataLock` | `internal/ball/orb_windows.go:466` | ②ignore | 写-only 即设计：`:508` 的 `orbDataLock = h` 本身是副作用（句柄寄存），"held OPEN from before the folder exists until process exit"，sweep 以句柄存在为活性信号（票 408r3 `da5d559b7` 判例）；删赋值＝句柄失去 Go 侧锚点 |
| 11 | `renderer.capW` | `internal/ball/renderer_windows.go:29` | ③删 | 唯一读写都在 `ensureDCAndDIB` 内（`:952` 读 / `:975` 写）；宿主函数零调用者 ⇒ 同生共死一起删 |
| 12 | `renderer.capH` | `internal/ball/renderer_windows.go:29` | ③删 | 同上 |
| 13 | `(*renderer).resize` | `internal/ball/renderer_windows.go:218` | ③删 | 全仓零调用；函数体一条 `comCall(BindDC)`，无副作用资源 |
| 14 | `(*renderer).overlayGlyphPixels` | `internal/ball/renderer_windows.go:838` | **④留** | winlive 测试真调用：`orb_live_windows_test.go:972`（文件头 `//go:build windows && winlive`，`TestBallLiveOrbOverlayFrameHasNoGlyphPixels` 把它 post 到 STA 线程审计 DIB）；⛔删⛔ignore，同 #9 处置 |
| 15 | `maxWindowEdge` | `internal/ball/renderer_windows.go:938` | ③删 | 全仓零调用；纯算术（`WindowEdgePx(BallSizeMaxPx, dpi)` 包装），无副作用 |
| 16 | `(*renderer).ensureDCAndDIB` | `internal/ball/renderer_windows.go:945` | ③删 | 全仓零调用；函数体副作用只落在 `r` 自身字段（hdcMem/hbm/bits/capW/capH），宿主死了副作用无处泄露；删后核过 `bitmapInfo`/`dibRGBColors`/`slotDCRTBindDC` 在 `newRenderer`（`:129/:138/:163`）另有活读者，**不级联出新死码** |
| 17 | `wndClassEx.style` | `internal/ball/win32_windows.go:175` | ①ignore | **WNDCLASSEX ABI**：`RegisterClassExW` 按偏移读 |
| 18 | `wndClassEx.clsExtra` | `internal/ball/win32_windows.go:177` | ①ignore | 同上 |
| 19 | `wndClassEx.winExtra` | `internal/ball/win32_windows.go:178` | ①ignore | 同上 |
| 20 | `wndClassEx.icon` | `internal/ball/win32_windows.go:180` | ①ignore | 同上 |
| 21 | `wndClassEx.bkgnd` | `internal/ball/win32_windows.go:182` | ①ignore | 同上 |
| 22 | `wndClassEx.menuName` | `internal/ball/win32_windows.go:183` | ①ignore | 同上 |
| 23 | `wndClassEx.iconSm` | `internal/ball/win32_windows.go:185` | ①ignore | 同上（`size`/`wndProc`/`instance`/`cursor`/`className` 有 Go 侧赋值不在 40 条内，未动） |
| 24 | `trackMouseInfo.hoverTime` | `internal/ball/win32_windows.go:247` | ①ignore | **TRACKMOUSEEVENT ABI**：`TrackMouseEvent` 按偏移读 |
| 25 | `bitmapInfoHeader.compression` | `internal/ball/win32_windows.go:220` | ①ignore | **BITMAPINFOHEADER ABI**：`CreateDIBSection` 按偏移读整个 40 字节头 |
| 26 | `bitmapInfoHeader.sizeImage` | `internal/ball/win32_windows.go:221` | ①ignore | 同上 |
| 27 | `bitmapInfoHeader.xPelsPerMeter` | `internal/ball/win32_windows.go:222` | ①ignore | 同上 |
| 28 | `bitmapInfoHeader.yPelsPerMeter` | `internal/ball/win32_windows.go:223` | ①ignore | 同上 |
| 29 | `bitmapInfoHeader.clrUsed` | `internal/ball/win32_windows.go:224` | ①ignore | 同上 |
| 30 | `bitmapInfoHeader.clrImportant` | `internal/ball/win32_windows.go:225` | ①ignore | 同上 |
| 31 | `bitmapInfo.bmiColors` | `internal/ball/win32_windows.go:230` | ①ignore | **BITMAPINFO ABI**：紧跟 40 字节头之后，`CreateDIBSection` 按偏移读 |
| 32 | `notifyIconData.dwState` | `internal/ball/win32_windows.go:262` | ①ignore | **NOTIFYICONDATA ABI**：`Shell_NotifyIconW` 按偏移读写 |
| 33 | `notifyIconData.dwStateMask` | `internal/ball/win32_windows.go:263` | ①ignore | 同上 |
| 34 | `notifyIconData.szInfo` | `internal/ball/win32_windows.go:264` | ①ignore | 同上 |
| 35 | `notifyIconData.uVersion` | `internal/ball/win32_windows.go:265` | ①ignore | 同上 |
| 36 | `notifyIconData.szInfoTitle` | `internal/ball/win32_windows.go:267` | ①ignore | 同上 |
| 37 | `notifyIconData.dwInfoFlags` | `internal/ball/win32_windows.go:268` | ①ignore | 同上 |
| 38 | `notifyIconData.guidItem` | `internal/ball/win32_windows.go:270` | ①ignore | 同上 |
| 39 | `notifyIconData.hBalloonIcon` | `internal/ball/win32_windows.go:271` | ①ignore | 同上 |
| 40 | `winHeapPtr` | `internal/ball/win32_windows.go:332` | ③删 | 全仓零调用（`cmd/balldebug` 用的是自己的 `winHeapPtrLocal` 副本，不是本符号的读者）；纯 uintptr→pointer 转换，无副作用；注释指向的 `renderer.dumpDIB` 根本不调用它（批 1 同判，文档已烂） |

（声明行号以改动前的基线读数为准；①档的行号在加 ignore 注释后自然下移。）

## 3. ★ABI 反向判据：删了编译照样过（现场实测）

本腿在收尾前做过一次现场演示：把 `bitmapInfoHeader.clrUsed`（①档 ABI 字段）真删掉，跑 `GOOS=windows go build ./internal/ball/` ⇒ **exit 0，一个错都不报**；随后立即原样还原并复跑 build＝绿。结论逐字落档：**删 ABI 字段，go build 不报错——编译器不是护栏**。这就是①档 ⛔删 的机器依据：错位只发生在运行期（窗口类注册／DIB 头／托盘数据）。

## 4. winlive 尺的如实报（存量坏编译的修复＋新显形 2 条）

- **修复前**：`-tags=winlive` 下 ball 测试包**编译本身是坏的**（存量，早于本腿）：
  `interaction_live_test.go:27` 本地重声明 GWL_EXSTYLE 索引常量（生产版在 `win32_windows.go:87`）；`:291` 还把 `^uintptr(19)` 直接喂 `pGetWindowLongPtrW.Call`（溢出）；`live_guard_windows_test.go:39-40` 本地重声明 `pSetWindowLongPtrW`/`pGetWindowLongPtrW`。
  后果①：派单"差集"判据被污染——winlive 多出的 3 条（`hotkey_windows.go` 的 `takeEscWith`/`cancelBorrowedLine`/`cancelFailedLine`）是**分析退化假阳性**：测试包被编译错误丢弃，`hotkey_status_test.go:216/:225/:265/:276/:280`（tag 仅 `//go:build windows`）那批默认档引用从分析里消失。它们在两把尺下的真实状态都是**有主**，不是 40 条之一，也不归四档处置。
  后果②：winlive 独有的死码被编译错误遮住看不见。
- **修复**（本腿写面内，internal/ball/**）：删两处本地重声明，调用点改用生产符号（生产版常量经 int32 变量转换，遵循 `win32_windows.go:83-87` 注释的既有规矩）。修复后 `-run TestNothingMatchesThis -tags=winlive` 编译冒烟 exit 0。
- **修复后新显形 2 条 U1000**：`live_guard_windows_test.go:44/:45` 的 `pGetClientRect`/`pClientToScreen`——全仓 grep 零引用的测试本地死 proc，按③删。
- 过程如实记：本腿一度误删④档的 `overlayGlyphPixels`（grep 命中时未先核 tag），发现后**立即原样恢复**，最终状态＝保留；此事不改变最终盘面，留档防重犯。

## 5. 门禁（逐项真实读数）

| 门禁 | 读数 |
|---|---|
| `go build ./internal/audio/ ./internal/ball/`（GOOS=windows CGO_ENABLED=0） | exit 0 |
| `go vet ./internal/audio/ ./internal/ball/` | exit 0 |
| `go test -count=1 ./internal/audio/ ./internal/ball/` | `ok .../internal/audio 18.374s`＋`ok .../internal/ball 0.854s`，exit 0 |
| `gofumpt -l internal/audio internal/ball` | 空（改动中曾点名 renderer_windows/win32_windows/interaction_live_test 三处对齐，已 `-w` 收口后复验为空） |
| `sh scripts/d22scan.sh` | **clean，0 finding**（`PASS=45 FAIL=0`； examined 481 production Go files） |
| SKIP | `internal/audio/aecspike`（sherpa cgo DLL 依赖，本机跑不了，如实 SKIP，未当绿）；`go test ./cmd/wisp/` 按派单票 98 豁免**未跑** |

## 6. 写面声明

改动只落在：`internal/audio/wasapi_windows.go`、`internal/ball/{ball,d2d,orb,renderer,win32}_windows.go`、`internal/ball/{interaction_live_test,live_guard_windows_test}.go`、本证据件、票 122 面 Progress log。`git diff --name-status` 与 commit hash 见回话（显式 pathspec 提交，未 push）。
