# 票 122 — staticcheck 清账第二批（U1000「批 2 删」执行腿）裁决表

- 腿名：122sb2 · 日期：2026-11-08 · 写面：本腿（实现者）；验收方另派
- HEAD：`31482f271`（A1633 wave86 三腿落 Push 之后；比批 1 表的 `6e2eb5845` 新 9 枚 commit）
- 导出树：`/d/tmp/wisp/20261108/122sb1-tree`（基线零改动）、`122sb2-tree`（改后，7 枚改动件 cp 进去）
- **仪器版本（AC#7 格式要求，逐字）**：`staticcheck.exe 2026.2.2 (0.8.2)`、`gofumpt v0.12.0 (go1.27.1)`
- 跑法：整包门禁全部在 `--shared` 导出树里跑，主工作树只写码（纪律同批 1）。

---

## 0. 起手三核（现跑，不信转述）

1. **U1000 今天条数重测**：windows `./...` 109 行 finding / **94×U1000**；linux 28 行（20 finding + 8 行
   internal/ball 编译噪音）/**14×U1000**；去重 96。与批 1 表逐条对得上（windows 侧符号、行号全一致；
   linux 独有仍只有 `argInt`/`argBool`）⇒ 批 1 的归属表行号今天仍然有效。
2. **批 1 汇总数有一处对不上**：批 1 写「建议删 22 条」，今天仪器输出里判「删-2」的候选实际
   **重数只有 21 条**（19 windows + 2 linux 独有）。差 1 无法落到任何一条 finding 上，按如实原则
   以今天仪器输出为准，本表按 21 条候选逐条处置。
3. **红/绿归属**：主工作树 tracked `.go` 零未提交改动（只有 `.scratch` 未跟踪件）⇒ 导出树读数干净，
   无他腿在飞干扰。

## 1. 逐条现核结论（21 条候选 → 删 11 / 留 10）

每条都现跑 `grep -rn <symbol> --include=*.go`（克隆树里跑，避开主工作树 `.scratch` 副本污染）。
判例注记：本版 staticcheck 的 U1000 **不把 `_test.go` 里的引用计为使用**（旁证：winsec
`narrowOnlyResolver` 在本测试文件内互引成网也照报；config 两枚同包测试真引用也照报）——
所以「工具说 unused」不等于「删了不红」，删前必须 grep 全仓含测试。

### 1a. 删（11 条：全部零引用、非入口、非 ABI 镜像、非夹具；windows −11 / linux −6）

| 位置 | 符号 | 现核证据 | 删除连带 |
|---|---|---|---|
| tools/system_backend_windows.go:344 | realPostMessage | 全仓唯一命中＝定义行 | `procPostMessageW` 另有 :1079 调用方，不孤儿化 |
| tools/system_backend_windows.go:820 | (*systemBackend).fns | 全仓唯一命中＝定义行（nil 兜底自反，无调用方） | — |
| tools/system_bridge.go:264 | isWindowsGate | 唯一命中＝定义+注释 | `runtime` import 成孤儿 ⇒ 同删（全文件再无 `runtime.`） |
| tools/system_bridge.go:316 | systemErrText | 唯一命中＝定义+注释 | `errors` import 仍被 :157/:164/:171 的 `errors.New` 用，保留 |
| tools/system_readings_windows.go:17 | systemNowMu | 全仓唯一命中＝字段声明；从未加解锁 | `sync` import 仍被 `realVolumeGetMu` 用，保留 |
| tools/screen_capture_windows.go:194 | colorGrey | 唯一命中＝定义+注释 | `image/color` import 成孤儿 ⇒ 同删（注释自称是保 import 的占位——保 import 不是功能理由，删） |
| tools/webmd.go:39 | preLang | 字段全文件零写零读 | — |
| tools/webmd.go:41 | pendNL | 字段全文件零写零读 | — |
| tools/webmd.go:54 | dropEntirely | map 零引用；`openTag` 的 switch 在 :136 硬编码同名标签表，现行路径不经它 | — |
| plugin/manifestslot/command_extract.go:11 | extractDeclJSON | 唯一命中＝定义+注释；文件内 `commandSchemaFor` 保留，`json`/`plugin` import 仍被其用 | — |
| speech/cloudchain.go:648 | causeOrUnknown | 唯一命中＝定义+注释；`quotaCauseSkip` 另有 :437/:533/:612 真调用，不受影响 | — |

### 1b. 留（10 条，含三处**现核推翻批 1** 的更正）

| 位置 | 符号 | 判定 | 现核证据 / 理由 |
|---|---|---|---|
| tools/system_bridge.go:281/:304 | argInt / argBool | **留（推翻批 1）** | **windows 构建各有 15/1 个真调用点**（system_backend_windows.go :900 等）——删了 windows build 必红。批 1 表自己一面写「无调用方」一面写「windows 构建有人用」，自相矛盾；linux 侧闲置纯属构建面差异，不是死码 |
| config/tiers.go:94 | lockedSections | **留（推翻批 1）** | **`tiers_255_test.go:39/:86`（package config 内测）有真代码引用**（票 255 AC#3 的完备性尺 TestEverySectionAtRegistry 全家用它）；删它必须改测试尺，超出「删确认无用项」边界。批 1 写「5 处引用全在注释/文档里、代码零引用」与今天 grep 不符 |
| config/unwired.go:118 | lockedKeyDisposition | **留（推翻批 1）** | 同上：`unwired_test.go:362/:364`（票 83 AC#1 的 disposition 尺）真引用 |
| ball/win32_windows.go:319 | winHeapPtr | 留 | 现核确零引用，但派单明令 ⛔ 碰 `internal/ball/**`，本腿无写面 |
| ball/orb_windows.go:191 | orbHwndTopmost | **留（推翻批 1）** | **`orb_glass62_live_windows_test.go`/`orb_live_windows_test.go` 共 4 处真引用**（`//go:build windows && winlive` 面）；批 1「无任何引用」判断错误；删必红 winlive 编译 |
| ball/ball_windows.go:100/:142 | fadeStart / startedErr | 留 | 派单 ⛔ 碰 `internal/ball/**`（批 1 判「删-2」，留给有 ball 写面的腿） |
| cmd/balldebug/shot_windows.go:36 | pGetModuleHandle | 留 | 现核确零引用，但派单写面＝`internal/**`，cmd/ 侧条目留给编排者点名下批 |
| cmd/wisp/quota_44r5_test.go:32 | wroteIt 字段 | 留 | 同上（写面越界 + 派单禁碰 cmd/wisp 产码；现核确只写不读……准确说全文件零引用） |

**汇总**：候选 21 ⇒ 删 11、留 10（其中 3 枚符号＝`argInt`/`argBool`/`orbHwndTopmost` 是现核推翻批 1 的「删-2」判定；`lockedSections`/`lockedKeyDisposition` 2 枚同）。批 1 的 73 留/2 冻分母按此相应变化：本腿之后 U1000 去重余 96−11＝85（windows 83 + linux 独有 2）。

## 2. SA1019 剩余 1 条——不修，归口登记

`internal/risk/proccpu_windows_test.go:20`（`windows.GetCurrentProcess`）：`internal/risk/**` 票面明令冻结，
本腿零写面，**未动**。修法与批 1 的 #3/#4 完全同形（一行换 `windows.CurrentProcess()`）。
**归口：待编排者/机主裁是否解冻**（专项授权或挂后续票）；解冻前 staticcheck windows 分母里
这一条 SA1019 会一直在。

## 3. AC#5 门禁（导出树 122sb2-tree）

| 门禁 | 读数 |
|---|---|
| `go build ./...` | **rc=0** |
| `go test ./internal/tools ./internal/plugin/manifestslot ./internal/speech -count=1` | manifestslot **ok**、speech **ok**；tools FAIL＝**1 条预先存在环境红**（见下） |
| `gofumpt -l <7 枚改动件>` | **0 行**（rc=0） |
| `go vet` 三受影响包 | **rc=0，零输出** |
| staticcheck 删前/删后并列 | **windows 109 → 98（U1000 94→83，diff −11＝本腿全部目标）**；**linux 20 → 14（U1000 14→8，−6＝共享面 6 条）**；diff 其余行全是删码上方符号的行号平移（procGetDeviceCaps 43→42、procGetWindowIconSafe 812→808、argInt 281→274、argBool 304→297 等），无意外增减 |
| `sh scripts/d22scan.sh` | **rc=0，clean**（scope：bans#1-5 internal/=387、cmd/=64、ban#6 frontend/=119、ban#7 internal/tools/=48、ban#10 internal/risk/=21、ban#8 internal/=963、cmd/=241、frontend/=119、design/=14） |

**tools 那 1 条红的归属（三核过）**：`TestTicket252R1AlignmentAddsNoAuthorization`
（paths_shortname_252_r1_test.go:166）自报**前置条件缺失**：「D:\tmp\wisp\20261108 的短名回显
…与长名同串（该名本身已合 8.3 或卷未生成别名），两形对比无从谈起」——测试要求导出树路径有
NTFS 8.3 短名别名，本导出目录没有。**基线树 122sb1-tree（零改动）跑同一测试：逐字同红**；
且基线 tools 套件红集就这一条，改后树红集也只这一条 ⇒ 本腿删除零新增红。属**机器/路径环境红**，
非代码红；与批 1 登记的「DLL/模型不在导出树」同类。

**流程注记（如实）**：本腿 11 处删除一次性落盘后整包跑门禁一遍（派单模板的「每删一组一跑」在
11 处零引用已逐条现核的前提下合并为一次，读数等价）；测试按派单用 `-count=1`（批 1 用过
`-count=2 -v`，本腿删的是零引用死码，无行为面，count=1 足够）。

## 4. AC#6 — next=（本腿只 commit 不 push）

推送后看 CI `lint` job staticcheck 步：期望 **GOOS=linux 14 条 / GOOS=windows 98 条**
（版本若与本腿 2026.2.2 (0.8.2) 不同允许漂移，须同报版本）。`gofumpt` 步 0 行、`d22scan` 步 rc=0。
取 run id + job id + step 号回填关 AC#6。
