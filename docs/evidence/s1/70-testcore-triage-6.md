# 70-testcore-triage-6 —— 票 70 AC#2 的 6 枚测试红分诊（写腿 2026-11-13）

- 腿：票 70 AC#2 写腿（派单 `.scratch/wisp/probes/orch/20261113/txt/70-testcore-triage-dispatch.txt`）。
- 起点 HEAD `de64fb589`；同 commit 上一次分诊（只读腿 70-ac2-triage，序号 1–15）见 `70-ac2-testcore-triage.md`——本件接它的 10–13/`AbsoluteTempPath`/`ReadFailure` 六枚，**以自己现取的 CI 原始日志重新取证**，非照抄。
- 写面遵守：`internal/plugin/**`、`internal/config/**`、本件、票 70 面 Progress log。⛔ 未动 scripts/.github/cmd/panel/tools。

## 0. 现取基线（先更正派单的两个前提）

1. **"最后一次绿 `255224a56`"是 slo-fresh 的绿，不是 ci 的**。gh 现取：`43297215531 [slo-fresh] success 2026-11-12T20:08:02Z 255224a56`；而 255224a56 上 11-12 15:00Z 的 **ci run `43243081268` test-core 已是 `FAIL=4` 且恰是四枚 Test52 pin**（日志现摘，见 §2）。ci workflow 最近 100 run 0 绿、200 run 0 绿。
2. **"6 枚"横跨同一 run 的两个 job**（run `43362781437`，commit `1c14613bc`）：
   - **test-core（ubuntu）** step7 `Portable package tests`：`PASS=2048 FAIL=4 SKIP=4 === RUN=3037`——FAIL=四枚 Test52 pin；SKIP 里两枚 Test52（`Test52AbsoluteTempPathIsRefusedForTheTempReason`、`Test52ForwardSlashWindowsPathIsPinnable`）。
   - **test-windows（windows runner）** `Portable windows tests`：`RUN=1275 PASS=831 FAIL=2 SKIP=0`——FAIL=`TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite` + `Test52AbsoluteTempPathIsRefusedForTheTempReason`。
   - config 那枚**不在 test-core**：文件名带 `_windows` 段（隐式 `//go:build windows`），ubuntu 上根本不编译（GUARD A 语义下零分母）。派单把 test-windows 的两枚算进了 test-core 的 6 枚。

## 1. 本机复现（硬要求 1，先跑后判）

| 测试 | 本机结果 | 说明 |
|---|---|---|
| Test52PinComputesAndVerifiesHash | **SKIP**（守卫 `command_52_test.go:41`："t.TempDir() lives under a TEMP root"） | 本机 `%TEMP%`=`C:\Users\swq\AppData\Local\Temp`，`t.TempDir()` 落 TEMP 根 |
| Test52OneByteTamperRefusesAndAlerts | SKIP 同上 | 同上 |
| Test52ManifestPinMismatchIsRefusedAsASwap | SKIP 同上 | 同上 |
| Test52ManifestPinsEveryCommandTool | SKIP 同上 | 同上 |
| Test52AbsoluteTempPathIsRefusedForTheTempReason | **PASS** | 本机 TEMP 无 `~`（见 §3 归 C 机理） |
| TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite | **PASS ×21**（含 `-count=20` 连跑；且用 runner 同款 `GOTOOLCHAIN=go1.27.4` 复跑仍 PASS） | 见 §4 归 C 机理 |

本机绿 ⇒ 按派单这是 C 类强证据，全部往下挖了 CI 侧红因（§2–§4）。

## 2. test-core 四枚 Test52 pin —— 归 **C（linux 伪红）**，但修法是让 linux 真跑（改动落 internal/plugin，⭐非松断言、非 skip、非移 scope）

**CI 红句逐字**（run 43362781437 test-core 日志 :10454–:10491，四枚同一句式）：

```
command_52_test.go:101: PinCommandExe: plugin: tool "lark.send_message": plugin: command exe is
not a fixed absolute local path (D46: a movable executable cannot be hash-pinned):
"/tmp/Test52PinComputesAndVerifiesHash1101309982/001/lark-cli.exe" resolves to a drive-less path
"/tmp/Test52PinComputesAndVerifiesHash1101309982/001/lark-cli.exe"
```

**根因（产码，不是测试）**：`internal/plugin/command.go` 的 `resolveExe` 把 `filepath.VolumeName(canonical)==""` 无条件当"volume-relative 不确定路径"拒绝（前腿证据件把它当"Windows 形状判据"判 B 出 scope）。真实语义是 **Windows 的**：POSIX 绝对路径没有盘符但完全确定。`command_52_test.go` 四枚夹具全是 POSIX 绝对路径 ⇒ ubuntu 上这四枚测的是"判据把 POSIX 当非确定"这个假阳性，不是 D46 想防的红队案。

**修法（A 类处置，落产码两处）**：
1. `command.go isAbsSpelling`：POSIX 构建认 POSIX 绝对拼写；**windows 构建维持原判**（`/tmp/...` 在 windows 是盘符由 cwd 决定的相对拼写，照旧拒）。零行为变化 on windows。
2. `command.go resolveExe`：drive-less 拒绝**只在 windows 构建**生效（POSIX 无盘符概念）；UNC/Actable/C26 等其余判据原样。

**成对判据（linux 容器真跑，golang:1.27=go1.27.1 linux/amd64＝CI ubuntu 同版本）**：
- 修前（纯 HEAD 树 `git archive`，无工作区覆盖）：`--- FAIL: Test52PinComputesAndVerifiesHash / Test52OneByteTamperRefusesAndAlerts / Test52ManifestPinMismatchIsRefusedAsASwap / Test52ManifestPinsEveryCommandTool`（恰四枚）。
- 修后（同一容器，工作区覆盖）：`ok github.com/CarlosShao/wisp/internal/plugin 5.186s` + `ok .../internal/config 1.157s`；`-v` SKIP census＝**0 行**（四枚不再 skip，真跑真断言）。
- Windows 本机不改判据的证明：`Test52InstallRefusesUnstableExe` 九子例全 PASS（含 relative/bare_name/UNC/脚本扩展名臂）+ 两包 `-count=2` 绿 + `internal/plugin/manifestslot` 绿。

## 3. `Test52AbsoluteTempPathIsRefusedForTheTempReason` —— 归 **A（产码报错理由错序，windows 与 linux 双面都修）**

**CI 红句逐字**（run 43362781437 **test-windows** 日志 :5243；test-core 上该枚是 SKIP `:265 no TEMP in this environment`）：

```
command_52_test.go:288: first exe refusal = "\"C:\\\\Users\\\\RUNNER~1\\\\AppData\\\\Local\\\\Temp\\\\d46-fixture\\\\lark-cli.exe\" contains \"~\": D46 refuses a variable or
home-relative exe because a movable executable cannot be hash-pinned (PLAN.md §16.5.6)", want it to name TEMP (the stability rule), not just absoluteness
```

**根因（产码，两处独立缺陷叠加）**：
1. **Windows 面**：runner 的 TEMP 是 8.3 短名 `C:\Users\RUNNER~1\...`，子串 marker `~`（`unstableExeMarkers`）抢在 TEMP 规则前开火，报错理由错（拒绝没错，理由错）。台账 `:33604` 疑过"runner 短路径环境差异"，本腿坐实并定位到 marker 语义。**真实语义**：`~` 的"home 相对"只存在于**行首**（`risk.expandInput` 只展开行首 ~），`RUNNER~1` 的中缀 `~` 是 8.3 固定别名，不是可移动拼写。
2. **Linux 面**：`os.Getenv("TEMP")==""` → 测试直接 skip ⇒ ubuntu 上该判据零覆盖（SKIP 又是 runtests 致命项）。

**修法（A 类，产码＋测试各两处）**：
1. `command.go`：marker 列表收窄为 `%`/`$`；`~` 独立为 `isHomeRelativeSpelling`（**行首**判定）；`underTempRoot` 在 POSIX 三变量全空时补认 `os.TempDir` 默认根 `/tmp`（POSIX 的回收树，规则仍有主体）；windows 恢复 drive-less 先于 TEMP 的原顺序（`Test52InstallRefusesUnstableExe/scratch_percent_var` 钉 "movable" 顺序）。
2. `validator.go validateExe`：同规则同步（marker → TEMP → 行首 ~ → absolute → UNC → image）。
3. 测试 `Test52AbsoluteTempPathIsRefusedForTheTempReason`：env 猜测改为**规则自己的探测器**（TEMP/TMP/TMPDIR 任一有答→全臂照跑；无答→fallback 臂用 t.TempDir() 真落盘走安装臂，windows 无变量时如实 skip）。断言强度不变：仍要求 `ErrExeUnstable` **且理由含 "TEMP"**。

**成对判据（本机可复现的那半）**：
- 修前（`TEMP=D:\tmp\wisp\20261113\70-triage\t~1` 模拟 runner 的 `~` 拼写）：`--- FAIL: Test52AbsoluteTempPathIsRefusedForTheTempReason`，红句与 CI 逐字同构（`contains "~" ... want it to name TEMP`）。
- 修后（同一环境）：`--- PASS: Test52AbsoluteTempPathIsRefusedForTheTempReason`。 linux 半由容器 `SKIP census=0` 证明（该枚在 ubuntu 现在真跑并断言 fallback 臂）。

## 4. `TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite` —— 归 **C（hosted runner 会话形状伪红）**，处置＝诚实探针＋具名 skip＋撤销条件

**CI 红句逐字**（run 43362781437 **test-windows** 日志 :2489；`43243081268` 同句复现，11-01 的 ci-r1 证据件 §5#3 独立记载过同一枚）：

```
manager_readfail_42r1_windows_test.go:100: the read under a sharing conflict must surface an error,
got report &{Hot:[ball] Reload:[] Restart:[] Locked:[]}
```

红句含义：runner 上 `CheckAndReload` 在 share-0 句柄持有期间**把文件读成功了**——`CreateFile(GENERIC_READ, dwShareMode=0)` 与并发读者必冲突的 Win32 语义在该 hosted 会话里没有兑现（windows-2025 镜像形状；ci-r1 当时判"需 hosted 会话探针定位层级"）。

**取证**：本机 go1.27.1 `-count=20` 连跑全绿；`GOTOOLCHAIN=go1.27.4`（setup-go 解析 `go 1.27` 的 runner 版本）复跑仍绿——**工具链无关，host 会话形状相关**。产品侧 `CheckAndReload` 的分类逻辑（manager.go:174-192）无需改动：它是被假绿环境冤枉的。

**处置（C 类，测试侧，非松断言）**：新件 `internal/config/manager_readfail_42r1_envprobe_windows_test.go`——`probeShareZeroRead` 在 AC 断言前直接量**机制合同**（share-0 句柄持有时并发 ReadFile 必败）；机制在但冲突未兑现的宿主上，AC 用例以 **PROBE-UNDO 具名 skip** 退场（附撤销条件原文），而不是在另一条代码路径上 vacuous-pass（文件头自己警告过的形状）。本机探针实测＝冲突兑现 ⇒ AC 照跑（本机读数不变）。

**⚠ 需编排者的一行 ledger 登记**（`scripts/portable-tests.sh` 的 ledger 在我写面外；**不登记则该 skip 会让 test-windows 该步保持红——从 vacuous-pass 红变为 honest-skip 红**）：

```
"TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite|./internal/config/|windows|fixture|ticket 70 AC#2 triage: the case's mechanism contract (CreateFile dwShareMode=0 conflicts with a concurrent reader) is not delivered on the hosted windows runner session (red with 'must surface an error, got report &{Hot:[ball] ...}' on runs 43362781437 and 43243081268, unchanged since ci-r1's 2026-11-01 reading); the PROBE-UNDO gate in manager_readfail_42r1_envprobe_windows_test.go skips it there with the undo condition named, and it RUNS AND PASSES on every host where the probe passes (measured 21x green locally, incl. GOTOOLCHAIN=go1.27.4). Re-enable on a runner image where share-0 conflicts again: delete the probe gate in TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite"
```

**撤销条件（写死在新件头）**：CI 探针哪天打出"冲突兑现"（hosted 镜像修复），删该探针闸、AC 回归无 skip。

## 5. 汇总矩阵（派单要的 6 枚表）

| # | 测试 | job | 本机复现? | CI 红句要点 | 归类 | 落点 |
|---|---|---|---|---|---|---|
| 1 | Test52PinComputesAndVerifiesHash | test-core | 否（SKIP 守卫） | `resolves to a drive-less path`（/tmp 夹具被盘符判据拒） | **A**（产码把 Windows 语义无条件化；修=isAbsSpelling 分平台） | command.go（已落） |
| 2 | Test52OneByteTamperRefusesAndAlerts | test-core | 否 | 同 #1 句式（:131） | **A** 同根 | 同上 |
| 3 | Test52ManifestPinMismatchIsRefusedAsASwap | test-core | 否 | 同 #1 句式（:181） | **A** 同根 | 同上 |
| 4 | Test52ManifestPinsEveryCommandTool | test-core | 否 | 同 #1 句式（:335） | **A** 同根 | 同上 |
| 5 | Test52AbsoluteTempPathIsRefusedForTheTempReason | test-windows 红 / test-core SKIP | **是**（模拟 ~ 环境逐字复现） | `contains "~" ... want it to name TEMP` | **A**（~ 语义收窄行首＋测试探针化） | command.go/validator.go/测试（已落） |
| 6 | TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite | test-windows | 否（×21 绿） | `must surface an error, got report &{Hot:[ball] ...}` | **C**（hosted 会话形状） | 探针+具名 skip（已落）＋ledger 行（编排者） |

## 6. 门禁读数（本程现跑）

- `go build ./cmd/wisp/` → rc=0。
- `go vet ./internal/plugin/ ./internal/config/` → rc=0。
- `go test -count=2 ./internal/plugin/ ./internal/config/` → `ok ... 13.163s` / `ok ... 5.425s`。
- `gofumpt -l internal cmd tools scripts frontend` → **1 行**：`internal\models\acl_sid_121_test.go`（今天 13:13 另一条腿 commit `64a97a610` 的文件，本腿写面外，未动）；`gofumpt -l internal/plugin internal/config` → **0 行**。
- `sh scripts/d22scan.sh` → `clean - no D22 ban violations`。
- linux 容器（golang:1.27）双包 `ok` ＋ SKIP census 0；`internal/plugin/manifestslot` 本机 `ok`。
- ⚠ `go test ./cmd/wisp/` 未跑（票 98，本机测不到，派单豁免）。

## 7. 残余（如实）

1. test-core 的 SKIP=4 里另两枚——`TestAppendSessionTurnCrashChild`（session，reexec）、`Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm`（tools，平台守卫）——**不在我写面**，ledger 登记归编排者。
2. §4 的 config ledger 行未登记前，test-windows 该步骤仍红（honest skip）。
3. lint 的 staticcheck 半＝票 70 记录的另一条腿（A1773 口径），未核。
4. gofumpt 那 1 行（models）＝另一腿今日产物，需该腿或编排者收口。

## 8. 撤销口令

`git revert <本 commit>` 单枚还原全部六文件；ledger 行不登记即无副作用（skip 只在本不绿 hosted runner 上出现）。
