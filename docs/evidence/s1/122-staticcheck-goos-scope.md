# 票 122 · staticcheck 量尺口径修正：U1000 按 GOOS 分庭（A1773 裁定的落地腿）

- 日期：2026-11-13 · HEAD（开工时）`8a8673fe8`（dev）· 收工 HEAD `de64fb589`
- 依据：台账 **A1773**（`docs/reports/pending-and-issues.md:35686`）——「47 U1000 里 44 条是
  linux 构建下真 unused、windows 构建下有主的 GOOS 盲区，删了会把 windows 功能搞坏；
  真正主流做法＝staticcheck 按 GOOS 分分母」；派单
  `.scratch/wisp/probes/orch/20261113/txt/122-goos-scope-dispatch.txt`。
- 写面：`.github/workflows/ci.yml`（staticcheck 步 run 块＋注释）＋本证据件＋票 122 面 Progress log。
  **零 .go 改动进 commit**（判据用的临时种子已拆净，`git diff` 对 `internal/ball/jarvis.go` 为空）。

## 1. 口径裁定（一句话版）

**U1000（未使用符号）是平台相关的判断**：一个符号"是否被使用"取决于哪个 GOOS 在编译它的调用者。
本仓是 Windows 优先桌面应用（PLAN.md D7），大量调用者住在 `//go:build windows` 文件里
（典型：`internal/ball` 的 jarvis 全族由 `jarvis_windows.go` 消费），linux runner 上判它们死＝量尺错，
不是代码错。故：

| 遍 | GOOS | checks | 判什么 |
|---|---|---|---|
| pass 1 | linux（=runner 原生） | `all,-U1000` | 跨平台通用缺陷（SA/ST/S/QF），两个 GOOS 下都成立 |
| pass 2 | windows（runner 交叉求值） | `all`（含 U1000） | 死代码判定，在"看得到全部调用者"的平台上 |

两遍**都必须 exit 0** 这一步才绿；windows 遍无 `continue-on-error`、无 `if:`（D22 mode-6 原样）。
**⛔ 未用 `-checks=-U1000` 全局关门**——那是放宽；U1000 仍有门，只是搬到了正确的 GOOS 上。

## 2. 落地形状（ci.yml staticcheck 步）

85a 立的三模块循环（`.` / `tools/d22scan` / `tools/mockllm`）、版本钉 2026.2.2、
finding 计数正则、`export data version` 崩溃行计数、"rc 永远不单独当读数"原则全部保留；
唯一改动＝循环体从单遍变两遍，自报行按遍分打分母：

```
staticcheck self-report: version=... modules=3 [pass1 GOOS=linux checks=all,-U1000] packages=P1 findings=F1 toolchain-crash-lines=C1 [pass2 GOOS=windows checks=all] packages=P2 findings=F2 toolchain-crash-lines=C2 (step exit is R; both passes must exit 0)
```

run 块 `bash -n` 语法验 rc=0（`D:\tmp\wisp\20261113\122-goos\extracted-run.sh`，临时区）。

## 3. 本机复跑读数（两遍官方口径 + 对照）

环境：Windows 11 宿主机（`go1.27.1`，`staticcheck.exe 2026.2.2 (0.8.2)`），复刻 CI 步的循环形状
（脚本：`D:\tmp\wisp\20261113\122-goos\runchk.sh`，与 ci.yml run 块逐行同形，只多 `GOOS/CGO_ENABLED`
环境注出）。工作树当时 `.scratch` 外零脏 .go ⇒ 读数即 HEAD `8a8673fe8` 的读数。

### pass 1 — GOOS=linux，`-checks=all,-U1000`

`packages=42 findings=32 toolchain-crash-lines=0 (step exit is 1)`，逐模块：

| module | exit | packages | findings |
|---|---|---|---|
| `.`（root） | 1 | 40 | 31 |
| `tools/d22scan` | 0 | 1 | 0 |
| `tools/mockllm` | 1 | 1 | 1 |

32 条按 checker：ST1003×10、ST1021×10、ST1020×3、ST1000×2、ST1022×2、ST1023×2、
SA5000×1（`internal/plugin/panic_isolation_51_test.go:62`，票面已登记"夹具本身、有意不修"）、
SA9003×1（`tools/mockllm/messages.go:145`，两遍同报）、SA9009×1（`frontend/embed.go:4`，已知假阳性，票 122 AC#3 登记格）。逐条全文见附录 A。

### pass 2 — GOOS=windows，`-checks=all`

`packages=44 findings=124 toolchain-crash-lines=0 (step exit is 1)`，逐模块：

| module | exit | packages | findings |
|---|---|---|---|
| `.`（root） | 1 | 42 | 123 |
| `tools/d22scan` | 0 | 1 | 0 |
| `tools/mockllm` | 1 | 1 | 1 |

124 条按 checker：**U1000×79**、ST1003×17、ST1021×11、ST1020×3、ST1000×2、ST1005×2、
ST1012×2、ST1022×2、ST1023×2、SA4010×1、SA5000×1、SA9003×1、SA9009×1。逐条全文见附录 B。

### 关键交叉读数

1. **CGO 等价性**：pass 2 用 `CGO_ENABLED=0`（runner 交叉分析时 cgo 必然关）。宿主上以
   `CGO_ENABLED=1`（mingw gcc 在场）同形复跑一遍：**finding 集逐字节一致（124=124）**，
   日志唯一差异＝自报行标签。仓内两处 cgo（`internal/audio/aec`、`internal/audio/aecspike`
   的 `import "C"` 文件）都带 `windows && cgo` tag 且有 `!cgo`/stub 配对文件 ⇒ GOOS=windows
   包集在 CGO=0 下完整可加载，44 包零加载损失。
2. **包数差异是真实形状**：`go list ./...` root 模块 GOOS=linux 40 包 / GOOS=windows 42 包，
   linux 侧整体缺席的 2 个：`build`、`cmd/balldebug`（其文件全被 windows 构建约束覆盖）。
3. **windows-only 加载文件上的非 U1000**（旧 CI 形状从未看过的 14 条，windows 遍新增视野）：
   `internal/audio/mmdevice_windows.go` ST1003×4、`internal/audio/wasapi_windows.go` ST1003×1、
   `internal/ball/hotkey_borrow_refused_260r5_test.go` ST1012×1、
   `internal/ball/hotkey_refused_borrow_release_245_test.go` ST1012×1、
   `internal/proc/boot_windows.go` ST1021×1、`internal/proc/jobscope_windows.go` ST1003×1、
   `internal/proc/treemetrics_windows.go` ST1003×1、`internal/secret/delete_test.go` SA4010×1、
   `internal/tools/recycle_windows.go` ST1005×1、`internal/tools/search_apps_windows.go` ST1005×1、
   `internal/winsec/pinnedleaf_342_windows.go` ST1000×1。
4. **仪器行为注记（U1000 集合随 checks 组合有两处反向移动）**：windows 默认 checks 口径
   U1000=82，全集口径 U1000=79。差集具名：`internal/tools/web_search.go` 的
   `unlandedProvider` 族 4 条只在**默认**口径报（全集口径多开 stylecheck 后该类型被判定有引用路径）；
   `internal/tools/web_bridge_test.go:49` `(*scriptGate).levels` 1 条只在**全集**口径报
   （此符号是票 122 AC#4 判表的"真死码候选"之一，全集口径把它捞回来了）。
   两口径各咬各的，谁也不是对方的子集——这正是"自报行必须带 checks 字样"的原因。

## 4. 与既有读数的对账（为什么数字对不上，逐条归因）

| 口径 | GOOS | checks | findings | U1000 | 出处 |
|---|---|---|---|---|---|
| 编排者派单现量 | windows | 默认 | 95 | 86 | 派单 :9（本机，2026-11-12 树） |
| CI run 43362781437 | linux | 默认 | 64 | 62 | 派单 :10（2026-11-12 的 CI） |
| 本腿复跑 | windows | 默认 | 89 | 82 | `D:\tmp\...\pass-windows-default.log` |
| 本腿复跑 | linux | 默认 | 54 | 52 | `D:\tmp\...\pass-linux-default.log` |
| **新口径 pass 1** | linux | `all,-U1000` | **32** | 0（口径排除） | `pass-linux.log` |
| **新口径 pass 2** | windows | `all` | **124** | 79 | `pass-windows-cgo0.log` |

归因（全部是**真实修复方向**的减少，无形状异常）：

- 编排者 95 → 本腿 89：之间 HEAD 推进了真实修码（票 439/440 波、`e918b4952` 已修的
  SA4004/S1016 在 122sa1/sb2 波里清掉）＋ 4 条 U1000 被清账波删除 ⇒ −6。
- CI 64 → 本腿 54：同样的 10 条 U1000 真删（62−10=52）＋ 2 条 linux 侧非 U1000 同期被真修 ⇒ −10。
- 默认 89 → 全集 124：多出 stylecheck 全族（ST1000/1003/1020/1021/1022/1023 共 43 条，
  默认不开）−3 条 U1000 集合移动（见 §3.4）＝ +35。

## 5. 成对判据（四发，逐字读数；日志 `D:\tmp\wisp\20261113\122-goos\probe.log`）

对象：`internal/ball/jarvis.go`（判据指定的文件；该族 38 条 U1000 由
`jarvis_windows.go` 消费＝A1773 盲区的标本）。种子＝临时函数 `jarvisProbe122Unused`。

| 发 | 形状 | 命令形状 | 读数 | 判 |
|---|---|---|---|---|
| A 正向·搬家 | 种子在场，linux 遍 `-checks=all,-U1000` | 同 pass 1，单包 `./internal/ball/` | `findings=1`（基线 ST1003），**probe-lines=0** | 判庭搬家，linux 不再对 GOOS 盲区误报 ✓ |
| B 正向·咬合 | 同一枚种子，windows 遍 `-checks=all` | 同 pass 2，单包 | `findings=43`（基线 42+1），**probe-lines=1**：`internal\ball\jarvis.go:239:6: func jarvisProbe122Unused is unused (U1000)` | windows 遍对真死码照咬、点名到行 ✓ |
| C 反向·对照 | 同一枚种子＋`//lint:ignore U1000`，windows 遍 | 同上 | `findings=42`（回基线），**probe-lines=0** | 指令级抑制才是"真闭嘴"（连正确平台都不看）——本修一条指令都没加，门没有关小 ✓ |
| D 拆除还原 | 种子删净，两遍复跑 | 同 pass 1/2 单包 | windows `findings=42`、probe-lines=0；linux `findings=1`＝基线 | 还原逐字节（`git diff` 空）✓ |

B 发即"反向判据"的实体：**同一枚真死码种子**，在 linux 遍沉默（A）、在 windows 遍被咬（B）——
口径改动没有把门关小，只是把判庭搬到了看得到证据的法庭；C 发证明如果哪天真要"不报"，
唯一手段是加抑制指令（可被 grep 审计），而不是改 GOOS 口径。

## 6. 边界（这把量尺看不到什么，具名）

1. **linux-only 加载文件的 U1000 不设门**：`//go:build !windows` 文件（如
   `internal/audio/wasapi_other.go`、`internal/winsec/pinnedleaf_342_unix.go`）只在
   GOOS=linux 下加载，其 U1000 两遍都看不到（pass 1 排除 U1000，pass 2 不加载这些文件）。
   本树该面读数＝0 条（linux 默认口径 U1000 52 条全部来自共用文件），故今天无新增盲区；
   但这是结构性边界，后续若长出 linux-only 死码，此步不报。
2. **pass 2 是交叉求值不是本机执行**：runner 上 GOOS=windows 意味着不跑 windows 二进制、
   只做静态分析；SA 族里依赖运行时行为的检查不受影响（staticcheck 本就是静态的），
   但 cgo 内部（C 侧）不可见——与 CGO 等价性读数一致（§3.1，集合逐字节同）。
3. **U1000 清理本身不在本腿**：A1773 明令"不是靠删符号"；真死码候选（3 条）与 44 条
   "留"的归属表已在 `docs/evidence/s1/122-u1000-judgment.md`，后续波按**新分母**清。

## 7. CI 能不能绿（如实）

**不能，还差得远——但这次是读得懂的红，不是口径错的红。** 新口径下 lint 的 staticcheck 步
需要两遍都净：pass 1 还差 **32 条**（31 root＋1 mockllm，分布见 §3；其中 30 条是
stylecheck 默认关的文风族，2 条是已登记的夹具/假阳性）；pass 2 还差 **124 条**
（79 U1000＋45 其他，逐条附录 B）。两份逐条清单就是后续清账波的工单面。

## 8. next

1. 编排者推送本 commit 后取 CI 步级读数：两遍分母应各为
   `[pass1 ... packages=? findings=?]` / `[pass2 ... findings=?]`（ubuntu runner 与本机
   GOOS/CGO 形状一致性以 CI 实测为准）；**期望 findings 接近 32/124**（非逐字节——runner
   的 go patch 版本与 .scratch 外并行腿的新 .go 会动分母，判读以两遍分开的数字为准）。
2. 票 122 后续清账波改按新分母领工单：pass 1 的 32 条（GOOS 无关）＋ pass 2 的
   windows-only 14 条（§3.3）＋ pass 2 的 U1000 79 条按 `122-u1000-judgment.md` 归属表清。
3. linux-only 文件 U1000 盲区（§6.1）今日读数 0，留此登记，不另开票。

## 附录 A — pass 1（GOOS=linux，`-checks=all,-U1000`）32 条逐条

- `frontend/embed.go:4:1: ineffectual compiler directive due to extraneous space: "// go:embed below is the ONLY way those bytes reach wisp.exe. There is no" (SA9009)`
- `internal/agent/compress.go:42:1: comment on exported type CompressionState should be of the form "CompressionState ..." (with optional leading article) (ST1021)`
- `internal/agent/images.go:1:1: package comment should be of the form "Package agent..." (ST1000)`
- `internal/agent/loop_golden_test.go:222:21: should omit type int from declaration; it will be inferred from the right-hand side (ST1023)`
- `internal/audio/aecspike/longrun.go:180:2: struct field dropDb should be dropDB (ST1003)`
- `internal/audio/aecspike/longrun.go:181:2: struct field lowAbsDb should be lowAbsDB (ST1003)`
- `internal/audio/aecspike/longrun.go:203:2: const collapseDropDb should be collapseDropDB (ST1003)`
- `internal/audio/aecspike/longrun.go:204:2: const collapseLowAbsDb should be collapseLowAbsDB (ST1003)`
- `internal/audio/audio.go:41:1: comment on exported const TargetRate should be of the form "TargetRate ..." (ST1022)`
- `internal/audio/bargein.go:64:2: const BargeinCollapseDropDb should be BargeinCollapseDropDB (ST1003)`
- `internal/audio/bargein.go:65:2: const BargeinCollapseLowAbsDb should be BargeinCollapseLowAbsDB (ST1003)`
- `internal/audio/bargein.go:197:2: struct field dropDb should be dropDB (ST1003)`
- `internal/audio/bargein.go:198:2: struct field lowAbsDb should be lowAbsDB (ST1003)`
- `internal/audio/bargein.go:534:1: comment on exported method UntrustedWindows should be of the form "UntrustedWindows ..." (ST1020)`
- `internal/ball/tokens.go:447:5: var wensityStopTs should be wensityStopTS (ST1003)`
- `internal/llm/content.go:19:1: comment on exported type MessageRole should be of the form "MessageRole ..." (with optional leading article) (ST1021)`
- `internal/llm/llm_test.go:28:9: should omit type Content from declaration; it will be inferred from the right-hand side (ST1023)`
- `internal/llm/resolver.go:36:1: comment on exported type RoleSettings should be of the form "RoleSettings ..." (with optional leading article) (ST1021)`
- `internal/memory/undo.go:49:6: var lastTs should be lastTS (ST1003)`
- `internal/plugin/manifest.go:11:1: comment on exported const CurrentSchemaVersion should be of the form "CurrentSchemaVersion ..." (ST1022)`
- `internal/plugin/panic_isolation_51_test.go:62:3: assignment to nil map (SA5000)`
- `internal/present/dispatch.go:13:1: comment on exported type Speaker should be of the form "Speaker ..." (with optional leading article) (ST1021)`
- `internal/risk/syncdirs.go:454:1: comment on exported method SyncRoots should be of the form "SyncRoots ..." (ST1020)`
- `internal/session/scope.go:60:1: comment on exported type EngineKind should be of the form "EngineKind ..." (with optional leading article) (ST1021)`
- `internal/statemachine/timeouts.go:18:1: comment on exported type TimeoutSpec should be of the form "TimeoutSpec ..." (with optional leading article) (ST1021)`
- `internal/tools/capability.go:8:1: comment on exported type Capability should be of the form "Capability ..." (with optional leading article) (ST1021)`
- `internal/tools/refusal.go:7:1: comment on exported type RefusalKind should be of the form "RefusalKind ..." (with optional leading article) (ST1021)`
- `internal/tools/registry.go:13:1: comment on exported type Kind should be of the form "Kind ..." (with optional leading article) (ST1021)`
- `internal/tools/system.go:395:1: comment on exported function SystemGetDecl should be of the form "SystemGetDecl ..." (ST1020)`
- `internal/tools/tool.go:11:1: comment on exported type Tool should be of the form "Tool ..." (with optional leading article) (ST1021)`
- `internal/winsec/pinnedleaf_342_unix.go:3:1: package comment should be of the form "Package winsec..." (ST1000)`
- `tools/mockllm/messages.go:145:3: empty branch (SA9003)`

## 附录 B — pass 2（GOOS=windows，`-checks=all`）124 条逐条

- `cmd/balldebug/shot_windows.go:36:2: var pGetModuleHandle is unused (U1000)`
- `frontend/embed.go:4:1: ineffectual compiler directive due to extraneous space: "// go:embed below is the ONLY way those bytes reach wisp.exe. There is no" (SA9009)`
- `internal/agent/compress.go:42:1: comment on exported type CompressionState should be of the form "CompressionState ..." (with optional leading article) (ST1021)`
- `internal/agent/images.go:1:1: package comment should be of the form "Package agent..." (ST1000)`
- `internal/agent/loop_golden_test.go:222:21: should omit type int from declaration; it will be inferred from the right-hand side (ST1023)`
- `internal/agent/reminder_timer_test.go:109:29: func (*fakeReminderClock).releasedCount is unused (U1000)`
- `internal/audio/aecspike/longrun.go:180:2: struct field dropDb should be dropDB (ST1003)`
- `internal/audio/aecspike/longrun.go:181:2: struct field lowAbsDb should be lowAbsDB (ST1003)`
- `internal/audio/aecspike/longrun.go:203:2: const collapseDropDb should be collapseDropDB (ST1003)`
- `internal/audio/aecspike/longrun.go:204:2: const collapseLowAbsDb should be collapseLowAbsDB (ST1003)`
- `internal/audio/aecspike/rig_windows.go:49:2: var iidIUnknown is unused (U1000)`
- `internal/audio/audio.go:41:1: comment on exported const TargetRate should be of the form "TargetRate ..." (ST1022)`
- `internal/audio/bargein.go:64:2: const BargeinCollapseDropDb should be BargeinCollapseDropDB (ST1003)`
- `internal/audio/bargein.go:65:2: const BargeinCollapseLowAbsDb should be BargeinCollapseLowAbsDB (ST1003)`
- `internal/audio/bargein.go:197:2: struct field dropDb should be dropDB (ST1003)`
- `internal/audio/bargein.go:198:2: struct field lowAbsDb should be lowAbsDB (ST1003)`
- `internal/audio/bargein.go:534:1: comment on exported method UntrustedWindows should be of the form "UntrustedWindows ..." (ST1020)`
- `internal/audio/mmdevice_windows.go:284:40: func parameter deviceId should be deviceID (ST1003)`
- `internal/audio/mmdevice_windows.go:291:33: func parameter deviceId should be deviceID (ST1003)`
- `internal/audio/mmdevice_windows.go:293:35: func parameter deviceId should be deviceID (ST1003)`
- `internal/audio/mmdevice_windows.go:295:42: func parameter deviceId should be deviceID (ST1003)`
- `internal/audio/wasapi_windows.go:36:2: var pGetCurrentThreadId should be pGetCurrentThreadID (ST1003)`
- `internal/audio/wasapi_windows.go:49:2: var iidIUnknown is unused (U1000)`
- `internal/ball/ball_windows.go:118:2: field fadeStart is unused (U1000)`
- `internal/ball/ball_windows.go:160:2: field startedErr is unused (U1000)`
- `internal/ball/ball_windows.go:1113:16: func (*Ball).recenterAt is unused (U1000)`
- `internal/ball/d2d_windows.go:139:2: field typ is unused (U1000)`
- `internal/ball/d2d_windows.go:144:2: field usage is unused (U1000)`
- `internal/ball/d2d_windows.go:145:2: field minLevel is unused (U1000)`
- `internal/ball/hotkey_borrow_refused_260r5_test.go:101:5: error var t260r5ErrNoWindowDoor should have name of the form errFoo (ST1012)`
- `internal/ball/hotkey_refused_borrow_release_245_test.go:58:5: error var t245ErrNoWindowDoor should have name of the form errFoo (ST1012)`
- `internal/ball/orb_windows.go:165:2: field fTransitionOnMaximized is unused (U1000)`
- `internal/ball/orb_windows.go:193:7: const orbHwndTopmost is unused (U1000)`
- `internal/ball/orb_windows.go:466:2: var orbDataLock is unused (U1000)`
- `internal/ball/renderer_windows.go:29:2: field capW is unused (U1000)`
- `internal/ball/renderer_windows.go:29:8: field capH is unused (U1000)`
- `internal/ball/renderer_windows.go:218:20: func (*renderer).resize is unused (U1000)`
- `internal/ball/renderer_windows.go:838:20: func (*renderer).overlayGlyphPixels is unused (U1000)`
- `internal/ball/renderer_windows.go:938:6: func maxWindowEdge is unused (U1000)`
- `internal/ball/renderer_windows.go:945:20: func (*renderer).ensureDCAndDIB is unused (U1000)`
- `internal/ball/tokens.go:447:5: var wensityStopTs should be wensityStopTS (ST1003)`
- `internal/ball/win32_windows.go:175:2: field style is unused (U1000)`
- `internal/ball/win32_windows.go:177:2: field clsExtra is unused (U1000)`
- `internal/ball/win32_windows.go:178:2: field winExtra is unused (U1000)`
- `internal/ball/win32_windows.go:180:2: field icon is unused (U1000)`
- `internal/ball/win32_windows.go:182:2: field bkgnd is unused (U1000)`
- `internal/ball/win32_windows.go:183:2: field menuName is unused (U1000)`
- `internal/ball/win32_windows.go:185:2: field iconSm is unused (U1000)`
- `internal/ball/win32_windows.go:220:2: field compression is unused (U1000)`
- `internal/ball/win32_windows.go:221:2: field sizeImage is unused (U1000)`
- `internal/ball/win32_windows.go:222:2: field xPelsPerMeter is unused (U1000)`
- `internal/ball/win32_windows.go:223:2: field yPelsPerMeter is unused (U1000)`
- `internal/ball/win32_windows.go:224:2: field clrUsed is unused (U1000)`
- `internal/ball/win32_windows.go:225:2: field clrImportant is unused (U1000)`
- `internal/ball/win32_windows.go:230:2: field bmiColors is unused (U1000)`
- `internal/ball/win32_windows.go:247:2: field hoverTime is unused (U1000)`
- `internal/ball/win32_windows.go:262:2: field dwState is unused (U1000)`
- `internal/ball/win32_windows.go:263:2: field dwStateMask is unused (U1000)`
- `internal/ball/win32_windows.go:264:2: field szInfo is unused (U1000)`
- `internal/ball/win32_windows.go:265:2: field uVersion is unused (U1000)`
- `internal/ball/win32_windows.go:267:2: field szInfoTitle is unused (U1000)`
- `internal/ball/win32_windows.go:268:2: field dwInfoFlags is unused (U1000)`
- `internal/ball/win32_windows.go:270:2: field guidItem is unused (U1000)`
- `internal/ball/win32_windows.go:271:2: field hBalloonIcon is unused (U1000)`
- `internal/ball/win32_windows.go:332:6: func winHeapPtr is unused (U1000)`
- `internal/llm/content.go:19:1: comment on exported type MessageRole should be of the form "MessageRole ..." (with optional leading article) (ST1021)`
- `internal/llm/llm_test.go:28:9: should omit type Content from declaration; it will be inferred from the right-hand side (ST1023)`
- `internal/llm/resolver.go:36:1: comment on exported type RoleSettings should be of the form "RoleSettings ..." (with optional leading article) (ST1021)`
- `internal/memory/l1_warm_29_test.go:151:6: func l1WaitFor is unused (U1000)`
- `internal/memory/undo.go:49:6: var lastTs should be lastTS (ST1003)`
- `internal/models/acl_sid_121_test.go:53:6: type aceLine is unused (U1000)`
- `internal/observe/logprivate_132_windows_test.go:99:6: func dirIsSealed132 is unused (U1000)`
- `internal/plugin/manifest.go:11:1: comment on exported const CurrentSchemaVersion should be of the form "CurrentSchemaVersion ..." (ST1022)`
- `internal/plugin/panic_isolation_51_test.go:62:3: assignment to nil map (SA5000)`
- `internal/present/dispatch.go:13:1: comment on exported type Speaker should be of the form "Speaker ..." (with optional leading article) (ST1021)`
- `internal/proc/boot_windows.go:22:1: comment on exported type Runtime should be of the form "Runtime ..." (with optional leading article) (ST1021)`
- `internal/proc/jobscope_windows.go:31:7: const jobObjectBasicProcessIdList should be jobObjectBasicProcessIDList (ST1003)`
- `internal/proc/treemetrics_windows.go:47:2: var procGetExtendedTcpTable should be procGetExtendedTCPTable (ST1003)`
- `internal/risk/anchor_world_334_test.go:55:2: const hostilePOSIX347 is unused (U1000)`
- `internal/risk/anchor_world_334_test.go:60:2: const danglingTarget347 is unused (U1000)`
- `internal/risk/syncdirs.go:454:1: comment on exported method SyncRoots should be of the form "SyncRoots ..." (ST1020)`
- `internal/secret/delete_test.go:131:9: this result of append is never used, except maybe in other appends (SA4010)`
- `internal/session/scope.go:60:1: comment on exported type EngineKind should be of the form "EngineKind ..." (with optional leading article) (ST1021)`
- `internal/statemachine/timeouts.go:18:1: comment on exported type TimeoutSpec should be of the form "TimeoutSpec ..." (with optional leading article) (ST1021)`
- `internal/tools/capability.go:8:1: comment on exported type Capability should be of the form "Capability ..." (with optional leading article) (ST1021)`
- `internal/tools/recycle_windows.go:151:13: error strings should not be capitalized (ST1005)`
- `internal/tools/refusal.go:7:1: comment on exported type RefusalKind should be of the form "RefusalKind ..." (with optional leading article) (ST1021)`
- `internal/tools/registry.go:13:1: comment on exported type Kind should be of the form "Kind ..." (with optional leading article) (ST1021)`
- `internal/tools/screen_capture_windows.go:42:2: var procGetDeviceCaps is unused (U1000)`
- `internal/tools/screen_capture_windows.go:43:2: var procSetCaptureDpiAwareness is unused (U1000)`
- `internal/tools/search_apps_windows.go:94:25: error strings should not be capitalized (ST1005)`
- `internal/tools/system.go:395:1: comment on exported function SystemGetDecl should be of the form "SystemGetDecl ..." (ST1020)`
- `internal/tools/system_backend_windows.go:35:2: var ole32 is unused (U1000)`
- `internal/tools/system_backend_windows.go:36:2: var shlwapi is unused (U1000)`
- `internal/tools/system_backend_windows.go:60:2: var procOpenInputDesktop is unused (U1000)`
- `internal/tools/system_backend_windows.go:63:2: var procMonitorFromWindow is unused (U1000)`
- `internal/tools/system_backend_windows.go:72:2: var procGetDpiForWindow is unused (U1000)`
- `internal/tools/system_backend_windows.go:76:2: var procGetWindowIcon is unused (U1000)`
- `internal/tools/system_backend_windows.go:77:2: var procOleInitialize is unused (U1000)`
- `internal/tools/system_backend_windows.go:78:2: var procOleUninitialize is unused (U1000)`
- `internal/tools/system_backend_windows.go:79:2: var procDwmGetWindowAttr is unused (U1000)`
- `internal/tools/system_backend_windows.go:80:2: var procMapVirtualKeyExW is unused (U1000)`
- `internal/tools/system_backend_windows.go:81:2: var procGetAsyncKeyState is unused (U1000)`
- `internal/tools/system_backend_windows.go:82:2: var procIsWindowEnabled is unused (U1000)`
- `internal/tools/system_backend_windows.go:83:2: var procGetWindowLongPtrW is unused (U1000)`
- `internal/tools/system_backend_windows.go:84:2: var procGetAncestor is unused (U1000)`
- `internal/tools/system_backend_windows.go:85:2: var procRegisterClipFmt is unused (U1000)`
- `internal/tools/system_backend_windows.go:86:2: var procGetClipboardSeqNum is unused (U1000)`
- `internal/tools/system_backend_windows.go:87:2: var procGetKeyState is unused (U1000)`
- `internal/tools/system_backend_windows.go:105:2: const gaRoot is unused (U1000)`
- `internal/tools/system_backend_windows.go:145:2: const dwmExtendedFrameBounds is unused (U1000)`
- `internal/tools/system_backend_windows.go:808:2: var procGetWindowIconSafe is unused (U1000)`
- `internal/tools/system_backend_windows.go:809:2: var dwmGetWindowAttrSafe is unused (U1000)`
- `internal/tools/system_dpi_windows.go:16:6: func setDPIAwarenessSafe is unused (U1000)`
- `internal/tools/tool.go:11:1: comment on exported type Tool should be of the form "Tool ..." (with optional leading article) (ST1021)`
- `internal/tools/web_bridge_test.go:49:22: func (*scriptGate).levels is unused (U1000)`
- `internal/winsec/acl_windows_test.go:32:2: const wellKnown is unused (U1000)`
- `internal/winsec/leftstanding_354_windows_test.go:68:2: const leftStandingInjected354r1 is unused (U1000)`
- `internal/winsec/pinnedleaf_342_windows.go:3:1: package comment should be of the form "Package winsec..." (ST1000)`
- `internal/winsec/seam_guard_windows_test.go:47:6: type narrowOnlyResolver is unused (U1000)`
- `internal/winsec/seam_guard_windows_test.go:49:27: func narrowOnlyResolver.Resolve is unused (U1000)`
- `internal/winsec/seam_guard_windows_test.go:53:27: func narrowOnlyResolver.ResolveAccounted is unused (U1000)`
- `internal/winsec/seam_guard_windows_test.go:70:31: field installs is unused (U1000)`
- `tools/mockllm/messages.go:145:3: empty branch (SA9003)`
