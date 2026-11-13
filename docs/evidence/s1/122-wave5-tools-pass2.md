# 票 122 — pass2 清账波 5（`internal/tools` ＋ `cmd/balldebug`）证据件

- 腿：122-wave5-tools-pass2（写腿）
- 开工 HEAD：`b39126832`（commit 前；`git log --oneline -1` 见文末）
- 派单：`.scratch/wisp/probes/orch/20261113/txt/122-wave5-tools-pass2-dispatch.txt`
- 工具版本（AC#7 格式）：staticcheck **2026.2.2 (0.8.2)** · gofumpt **v0.12.0 (go1.27.1)** · go **1.27.1 windows/amd64**
- 口径：`GOOS=windows CGO_ENABLED=0 staticcheck -tags=winlive -checks=all,-SA9009 ./internal/tools/ ./cmd/balldebug/`（自跑，未照抄派单数）

## 1. 复跑前读数（开工现测，28 条）

```
cmd\balldebug\shot_windows.go:36:2 U1000 pGetModuleHandle
internal\tools\recycle_windows.go:151:13 ST1005
internal\tools\screen_capture_windows.go:42:2 U1000 procGetDeviceCaps
internal\tools\screen_capture_windows.go:43:2 U1000 procSetCaptureDpiAwareness
internal\tools\search_apps_windows.go:94:25 ST1005
internal\tools\system_backend_windows.go:35:2 U1000 ole32
internal\tools\system_backend_windows.go:36:2 U1000 shlwapi
internal\tools\system_backend_windows.go:60:2 U1000 procOpenInputDesktop
internal\tools\system_backend_windows.go:63:2 U1000 procMonitorFromWindow
internal\tools\system_backend_windows.go:72:2 U1000 procGetDpiForWindow
internal\tools\system_backend_windows.go:76:2 U1000 procGetWindowIcon
internal\tools\system_backend_windows.go:77:2 U1000 procOleInitialize
internal\tools\system_backend_windows.go:78:2 U1000 procOleUninitialize
internal\tools\system_backend_windows.go:79:2 U1000 procDwmGetWindowAttr
internal\tools\system_backend_windows.go:80:2 U1000 procMapVirtualKeyExW
internal\tools\system_backend_windows.go:81:2 U1000 procGetAsyncKeyState
internal\tools\system_backend_windows.go:82:2 U1000 procIsWindowEnabled
internal\tools\system_backend_windows.go:83:2 U1000 procGetWindowLongPtrW
internal\tools\system_backend_windows.go:84:2 U1000 procGetAncestor
internal\tools\system_backend_windows.go:85:2 U1000 procRegisterClipFmt
internal\tools\system_backend_windows.go:86:2 U1000 procGetClipboardSeqNum
internal\tools\system_backend_windows.go:87:2 U1000 procGetKeyState
internal\tools\system_backend_windows.go:105:2 U1000 const gaRoot
internal\tools\system_backend_windows.go:145:2 U1000 const dwmExtendedFrameBounds
internal\tools\system_backend_windows.go:808:2 U1000 procGetWindowIconSafe
internal\tools\system_backend_windows.go:809:2 U1000 procDwmGetWindowAttrSafe
internal\tools\system_dpi_windows.go:16:6 U1000 func setDPIAwarenessSafe
internal\tools\web_bridge_test.go:49:22 U1000 (*scriptGate).levels
```

28 条，与派单清单逐条一致（0 新增 0 缺失）。

## 2. 复跑后读数

同口径命令 → **0 行输出，exit 0**。
不带 tag 复核（默认档 `-checks=all,-SA9009` 同两包）→ **0 行输出，exit 0**（本波无④类，两档同空）。

## 3. 逐条表（符号 / 位置 / 四档 / 理由）

| # | 位置 | 符号 | 档 | 处置 | 理由（现核证据） |
|---|---|---|---|---|---|
| 1 | system_backend_windows.go:35 | `ole32` | ③ | 删 | `NewLazySystemDLL("ole32.dll")` 句柄；全仓唯一下游 = 本文件 `procOleInitialize`/`procOleUninitialize`（同批死）；剪贴板/回收站另有各自的 ole32 入口（recycle_windows.go:59-60 内联取）不依赖它。整组同进退 |
| 2 | :36 | `shlwapi` | ③ | 删 | `NewLazySystemDLL("shlwapi.dll")` 句柄；全仓零 proc 从它取（forks/go-webview2 与 tools 无关）。彻死 |
| 3 | :60 | `procOpenInputDesktop` | ③ | 删 | `user32.NewProc("OpenInputDesktop")`；本文件 1468 行 + 全仓零引用。桌面枚举动词后端无实现 |
| 4 | :63 | `procMonitorFromWindow` | ③ | 删 | 同上，零引用；窗口贴靠走 `snapGeometry`（纯几何，读 `windowRect`），不需要 monitor API |
| 5 | :72 | `procGetDpiForWindow` | ③ | 删 | 零引用；DPI 声明实际入口 = `ensureDPIAwareness`（system_ops_windows.go:100，`platformSystemBackend:85` 活调用），用的是 `procSetProcessDpiAwareness` |
| 6 | :76 | `procGetWindowIcon` | ③ | 删 | `= procGetWindowIconSafe` 纯别名；零引用 |
| 7 | :77 | `procOleInitialize` | ③ | 删 | `ole32.NewProc("OleInitialize")`；零引用（与 #1 同组） |
| 8 | :78 | `procOleUninitialize` | ③ | 删 | 同上 |
| 9 | :79 | `procDwmGetWindowAttr` | ③ | 删 | `= dwmGetWindowAttrSafe` 纯别名；零引用。本文件无任何 DWM 调用 |
| 10 | :80 | `procMapVirtualKeyExW` | ③ | 删 | 零引用；键序列走 `parseKeySequence`/`namedKeys` 表 + SendInput scan 码 |
| 11 | :81 | `procGetAsyncKeyState` | ③ | 删 | 零引用；按键注入是 SendInput 单向，不读异步键态 |
| 12 | :82 | `procIsWindowEnabled` | ③ | 删 | 零引用；可用性检查走 `IsIconic`/`IsWindowVisible`（另 proc，活的） |
| 13 | :83 | `procGetWindowLongPtrW` | ③ | 删 | 零引用 |
| 14 | :84 | `procGetAncestor` | ③ | 删 | 零引用（`gaRoot` 常量就是给它配的，见 #22） |
| 15 | :85 | `procRegisterClipFmt` | ③ | 删 | 零引用；剪贴板读写只用 `CF_UNICODETEXT=13`（`cfUnicodeText`，活） |
| 16 | :86 | `procGetClipboardSeqNum` | ③ | 删 | 零引用 |
| 17 | :87 | `procGetKeyState` | ③ | 删 | 零引用（同 #11） |
| 18 | :105 | `gaRoot` | ③ | 删 | `GA_ROOT` 常量唯一消费者是 `procGetAncestor`（#14）——两者一起死，真删不违反「删一半留一半」 |
| 19 | :145 | `dwmExtendedFrameBounds` | ③ | 删 | `DWMWA_EXTENDED_FRAME_BOUNDS` 常量唯一消费者是 DWM 属性读（#9 族）——一起死 |
| 20 | :808 | `procGetWindowIconSafe` | ③ | 删 | 声明为 `*windows.LazyProc` 且**从不赋值**，恒 nil，零调用——比 proc 更死；随 #6 别名一起整链删 |
| 21 | :809 | `dwmGetWindowAttrSafe` | ③ | 删 | 同 #20 |
| 22 | system_dpi_windows.go:16 | `setDPIAwarenessSafe` | ③ | 删（整文件） | `procSetCaptureDpiAwareness = setDPIAwarenessSafe`（screen_capture_windows.go:43）是它唯一引用，而该 var 自身零调用 ⇒ 死链。且它**与活函数重复**：`ensureDPIAwareness`（system_ops_windows.go:100，被 `platformSystemBackend:85` 每次后端构建调用）用同一 `procSetProcessDpiAwareness` 干同一件事。活逻辑不受影响（详见 §5） |
| 23 | screen_capture_windows.go:42 | `procGetDeviceCaps` | ③ | 删 | `gdi32.NewProc("GetDeviceCaps")`；全文件零调用。BitBlt 通路（`realCaptureScreen`/`realCaptureWindow`/`copyToPNG`）不用它 |
| 24 | :43 | `procSetCaptureDpiAwareness` | ③ | 删 | `= setDPIAwarenessSafe` 死链头部；截屏前 DPI 已由 `ensureDPIAwareness` 进程级声明过，截屏文件内无需再调（见 §5） |
| 25 | search_apps_windows.go:94 | `errNoAppPathsHive` 文案 | ST1005 | 真改 | `"App Paths 注册表不可读"` → `"读不到 App Paths 注册表"`（首字符从大写拉丁字母改为 CJK，语义不变）；改前 grep 全仓：该错误文本无任何测试断言/解析依赖；错误变量本身有主（:87 返回），非死码 |
| 26 | recycle_windows.go:151 | 报错文案 | ST1005 | 真改 | `"Shell 操作被中止：…"` → `"shell 操作被中止：…"`（仅首字母大小写，语义不变）；改前 grep 全仓：该文本无测试断言/解析依赖 |
| 27 | web_bridge_test.go:49 | `(*scriptGate).levels` | ③ | 删 | 测试 helper；全包 `grep '\.levels('` 零调用（AC#4 表 2026-11-12 判真死码候选，本腿现核坐实）；`find`/`seen` 等姊妹方法活 |
| 28 | cmd/balldebug/shot_windows.go:36 | `pGetModuleHandle` | ③ | 删 | `kernel32.NewProc("GetModuleHandleW")`；全包零引用（cmd/wisp/notify_windows.go 的 `pGetModuleHandleW` 是另一符号、自取自用，不受影响）。批 2 曾判「删-2」，当时 cmd 侧越界，本波写面含 cmd/balldebug，落删 |

**四档落点汇总**：③真无主删 24 项（含 2 条纯别名、1 个整文件、1 对「常量+唯一消费者」）；ST1005 真改 2（派单明令文风不许 ignore 蒙）；**①ABI ignore 0 条 ②故意保留 ignore 0 条 ④点名保留 0 条**。全波 **零 lint:ignore**。

## 4. ★procX 整链核对（派单风险 1 的回答）

- 15 枚 user32/kernel32 死 proc + 2 枚 ole32 proc + 2 枚 DLL 句柄（ole32/shlwapi）+ 2 枚死常量（`gaRoot`/`dwmExtendedFrameBounds`，唯一消费者就是上面这批 proc）+ 2 枚恒 nil 的 Safe 变量 + 2 枚别名 var：**23 个声明在同一文件整组同进退**，删后 `go build`/`go vet` rc0、test ok——无悬空调用（该文件剩余的 proc 全部有活调用：`procGetForegroundWindow`/`procIsIconic`/`procBitBlt` 等，逐一未被本波触碰）。
- `setDPIAwarenessSafe` 死链 3 节点（`procSetCaptureDpiAwareness` → `setDPIAwarenessSafe` → `procSetProcessDpiAwareness`/`shcore`）删了前两节；第三节 `procSetProcessDpiAwareness`/`shcore` **留**——`ensureDPIAwareness`（system_ops_windows.go:100-114）在用，删它会砸真功能。
- 全仓 `git grep` 核对过 27 个符号名（.go+.md），命中只剩：docs/evidence 历史档（批 1「留-口」表，历史不改）、`.scratch` probe 夹具（非仓码）、他包同名不同物（`modOle32`/`pnlModOle32`/`pGetModuleHandleW` 均独立声明）。**零处引用被本波删码打断。**

## 5. 为什么删 setDPIAwarenessSafe 不是撕掉截屏 DPI 修复

- `ensureDPIAwareness`（system_ops_windows.go:94-114，`dpiOnce` 幂等）在 `platformSystemBackend`（:85）里于**每次构建 Windows 后端时**执行——即任何 system 动词（含 `screen.capture`）到达后端之前。它是全包唯一的 DPI 声明入口，用的是与死副本同一个 `procSetProcessDpiAwareness`。
- 死副本 `setDPIAwarenessSafe` 与它三处行为差异（无 sync.Once、错误映射不同、`ERROR_ACCESS_DENIED` 分支）全部无消费者，删整文件不改变任何运行时行为。test `Test23ScreenCapture*` 四条 + `Test23CaptureAltNamesDRMBlackFrame` 全 PASS 佐证。

## 6. 门禁读数（逐项，含版本）

| 门禁 | 读数 |
|---|---|
| `go build ./internal/tools/ ./cmd/balldebug/` | rc0，无输出 |
| `go vet` 同两包 | rc0，无输出 |
| `go test -count=1 ./internal/tools/` | **ok 50.955s**；`-v` 复跑 ok 39.629s，**SKIP 恰 1 条**：`Test23LiveClipboardRoundTripThroughRealOS`（真 OS 剪贴板占用夹具，环境性，改动前同形态；其余全 PASS 无 FAIL） |
| staticcheck `-tags=winlive`（派单口径） | **28 → 0 行，exit 0**；不带 tag 复核也 0 行（无④类，两档同空） |
| gofumpt -l | 初跑点名 2 件（本波删行后的 var 块对齐）→ `-w` 收口 → 复跑 **0 行**；diff 复核 = 纯对齐（shot_windows.go / screen_capture_windows.go），无格式化引入的语义改动 |
| `bash scripts/d22scan.sh` | **clean rc0**：PASS=45 FAIL=0 SKIP=0（=== RUN 93）；扫面 `d22scan: clean - no D22 ban violations`，internal/ 403 + cmd/ 77 产码文件照常入镜（tools 51） |
| `go test ./cmd/wisp/` | 按派单明令（票 98 豁免）**未跑** |

## 7. 写面核对

`git diff --name-status`（本波全部改动）：

```
D         internal/tools/system_dpi_windows.go
M         cmd/balldebug/shot_windows.go
M         internal/tools/recycle_windows.go
M         internal/tools/screen_capture_windows.go
M         internal/tools/search_apps_windows.go
M         internal/tools/system_backend_windows.go
M         internal/tools/web_bridge_test.go
```

7 个 .go，全部在派单写面（`internal/tools/**` ＋ `cmd/balldebug/**`）。共享树另有的他腿脏件（`.scratch/wisp/probes/291/...`、`346-a2-dispatch.txt`、`frontend/fixtures/composer-states.html`）**零触碰零提交**。统计：19 insertions / 98 deletions（.go 部分；另有证据件＋票面 log）。

## 8. 一句话收口

**internal/tools ＋ cmd/balldebug 的 pass2 本波归零**（28→0，两档口径同空）；全仓 pass2 剩余面归 audio/ball 之外的后续波与 AC#6 的 CI 步级读数。
