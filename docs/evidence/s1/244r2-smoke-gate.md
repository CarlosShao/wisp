# 244-r2 — 冒烟形修「真读退出码」的证明（GUI 子系统 stub 两态读数）

- **腿**: `244-r2`（派单存档 `.scratch/wisp/probes/orch/20261102/txt/244r2-dispatch.txt`，2026-11-02 11:4x）
- **任务**: `scripts/build.ps1:309-310` 的冒烟从 `& (build\wisp.exe) doctor; if ($LASTEXITCODE -ne 0){Fail}` 改为**真正等到 GUI 进程退出 + 取真实 ExitCode** 的形态（镜像 `scripts/slo-check.ps1:143-200` 的 `Invoke-ExternalProgram`）。
- **边界**: ⛔只改 `scripts/build.ps1`；⛔不碰 `cmd/wisp/**` / `slo-check.ps1`（263 已修）；stub 临时件放 `D:\tmp\wisp\20261102\244r2\`。
- **收尾锚点**: HEAD=`c9ef940c853d417d528e0c47e1b0887fcac6cf14`；本文件是唯一交付物（commit message 行：`ticket 244-r2 smoke gate shape fix, scripts/build.ps1 only`，显式 pathspec `D:\work\workspace\projects plans\Wisp\scripts\build.ps1`, `D:\work\workspace\projects plans\Wisp\docs\evidence\s1\244r2-smoke-gate.md`）。

---

## 1. P1 描述与旧形态（`244-c1-census.md` §2 + §4 row 1）

`244-c1-census.md` §4 第 1 行明确列出 `build.ps1:309-310` 为 AC#6 头号残量：

```text
| 1 | `scripts/build.ps1` 冒烟 | `:309-310` [`& (build\wisp.exe) doctor`] → [`if ($LASTEXITCODE -ne 0)`] | `&` 起产物＋读自动变量 | **不成立⚠本票头号残量**. 〔263 量〕同形实测（`slo-check.ps1:157-164` 逐字记录：CUI 探针 2049ms 带 rc＝3 返回，GUI 探针 **70ms 放空、下一次取值在严格模式下死**）⇒ GUI 子系统下 `&` 不等待；本文件 `:40-41` 有 `$ErrorActionPreference='Stop'`＋`Set-StrictMode 2.0`，但 `$LASTEXITCODE` 已被 `:270` 的 `& go build`（rc=0）赋过值 ⇒ `:310` 读到的是**前一道命令的陈值 0 ⇒ 冒烟自 -H 落地起结构性空转（doctor 真 FAIL 也不会红）**。
```

**结论**: 该门自 `-H=windowsgui` 落地后，无论 `wisp doctor` 结果如何，都会以绿色"静默通过"。

---

## 2. 自建 GUI-subsystem stub 证明（现量）

### 2.1 Stub 来源与子系统设计

新建一枚 GUI 子系统测试程序（非仓库产物），`os.Exit(7)`，并支持延时参数以区分"等 vs 不等":

- **位置**: `D:\tmp\wisp\20261102\244r2\stub\main.go`
- **编译命令**:
  - GUI: `go build -ldflags '-H=windowsgui' -o stubgui.exe .`
  - CUI: `go build -o stubcui.exe .`
- **PE 子系统读数** (`objdump -p Subsystem`):
  - `stubgui.exe`: `Subsystem 00000002 (Windows GUI)` ✓
  - `stubcui.exe`: `Subsystem 00000003 (Windows CUI)` ✓
- **行为**: 启动后等待 `STUB_SLEEP_MS` 毫秒（默认 `0`），然后调用 `os.Exit(code)`（code 可从环境变量 `STUB_EXIT` 指定，默认 `7`）。

### 2.2 两态对比测试（`test-shape.ps1` 现量）

**TEST 1: OLD FORM (& + $LASTEXITCODE) on GUI SUBSYSTEM**

```powershell
& cmd /c exit 0                          # pre-seed LASTEXITCODE=0
$env:STUB_SLEEP_MS = '2000'              # process will sleep for 2 seconds
$t0 = Get-Date
& stubgui.exe doctor                     # old form: bare call without wait
$elapsedOld = ((Get-Date) - $t0).TotalMilliseconds
$guiOldResult = $LASTEXITCODE
Write-Host "old form result: $guiOldResult (true code is 7)"
```

**读数**: `elapsedOld ≈ 21 ms`，`$guiOldResult = 0`  
**[VERDICT]**: **STALE READ CONFIRMED** — `$LASTEXITCODE` 保留了 `& cmd /c exit 0` 的预置值 `0`，因为 `&` 没有等到 GUI 进程结束。

**TEST 2: NEW FORM (Process.Start + WaitForExit + .ExitCode) on GUI SUBSYSTEM**

```powershell
function Invoke-ExternalProgram { ... }   # defined below, mirrors slo-check.ps1
$r = Invoke-ExternalProgram -FilePath .\stubgui.exe -ArgumentList @('doctor')
Write-Host "new form result: $($r.code)"
```

其中函数实现（`Process.Start` + `WaitForExit()` + `$proc.ExitCode`）完全复现了 `slo-check.ps1:143-200` 的逻辑。

**读数**: `$r.code = 7`（耗时 > 2000 ms）  
**[VERDICT]**: **TRUE EXIT CODE CONFIRMED** — 真正读取到子进程的退出码 `7`。

**TEST 3–4: CUI SUBSYSTEM CONTROL**

作为对照，`stubcui.exe` 是控制台子系统（CUI），即使使用 OLD FORM 也正确等了约 2000ms 并读出 `7`，因为 CUI 二进制在 PowerShell 下的语义不同。

---

## 3. patched `build.ps1` 新形态（lines 307–370）

修改后的 `scripts/build.ps1` 冒烟段：

```powershell
# --- 6. smoke test ----------------------------------------------------------
# build\wisp.exe links with -H=windowsgui (`:256`), so it is a GUI-subsystem
# binary, and a bare '&' call is NOT waited for it: measured 2026-10-04 in
# ticket 263 (console probe returned after 2049ms with rc=3; the GUI probe
# fell through in 70ms), and re-measured 2026-11-02 by ticket 244-r2 with a
# self-built GUI stub that sleeps 2000ms then os.Exit(7): '&' returned after
# 21ms and $LASTEXITCODE read the stale 0 that '& $go build' (`:270`) had left
# there (docs/evidence/s1/244r2-smoke-gate.md). That made this gate
# structurally green no matter what doctor said - the exact defect the
# 244-c1 census §4 row 1 names. The shape below mirrors
# scripts/slo-check.ps1 Invoke-ExternalProgram (:143-200, ticket 263):
# Process.Start + WaitForExit() + the real ExitCode. doctor FAIL => red here.
function Invoke-ExternalProgram {
    param(
        [Parameter(Mandatory)][string]$FilePath,
        [string[]]$ArgumentList = @()
    )
    if (-not (Test-Path -LiteralPath $FilePath)) {
        return [pscustomobject]@{ ok = $false; code = $null; why = "program not found: $FilePath" }
    }
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = $FilePath
    $psi.Arguments = -join @(foreach ($v in $ArgumentList) {
        $s = [string]$v
        if ($s -match '[\s"]') { ' "' + ($s -replace '"', '\"') + '"' } else { " $s" }
    })
    # UseShellExecute=$false is the measured shape (ticket 263): the child's
    # own lines still reach this caller's streams, so the CI log keeps the
    # doctor output, and the exit code comes back the same.
    $psi.UseShellExecute = $false
    $proc = $null
    try {
        $proc = [System.Diagnostics.Process]::Start($psi)
    } catch {
        return [pscustomobject]@{ ok = $false; code = $null; why = ("cannot start '{0}': {1}" -f $FilePath, $_.Exception.Message) }
    }
    if ($null -eq $proc) {
        return [pscustomobject]@{ ok = $false; code = $null; why = ("no process object returned for '{0}'" -f $FilePath) }
    }
    try {
        $proc.WaitForExit()
        return [pscustomobject]@{ ok = $true; code = [int]$proc.ExitCode; why = '' }
    } catch {
        return [pscustomobject]@{ ok = $false; code = $null; why = ("'{0}' started but its exit code could not be read: {1}" -f $FilePath, $_.Exception.Message) }
    } finally {
        $proc.Dispose()
    }
}

Write-Host 'build.ps1: smoke test - running wisp.exe doctor'
$smokeExe = Join-Path $outDir 'wisp.exe'
$smoke = Invoke-ExternalProgram -FilePath $smokeExe -ArgumentList @('doctor')
if (-not $smoke.ok) { Fail "smoke test: 'wisp.exe doctor' could not be measured: $($smoke.why)" }
if ($smoke.code -ne 0) { Fail "wisp doctor reported FAIL (exit code $($smoke.code))." }
```

关键改动：
1. **定义 `Invoke-ExternalProgram`**：与 `slo-check.ps1:143-200` 等价的 `Process.Start + WaitForExit() + ExitCode` 形态。
2. **调用替代 `&`**：`$smoke = Invoke-ExternalProgram ...`。
3. **条件判断基于 `$smoke.code`**：失败时 `Fail` 报告具体退码。

**注释说明** (`# --- 6. smoke test`):
- `build\wisp.exe` 链接为 GUI 子系统 (`-H=windowsgui` on line `256`)。
- 旧形态 (`&`) 对 GUI 进程**不等待**。
- `Process.Start + WaitForExit + ExitCode` 保证读取真实退码。
- **doctor 真 FAIL ⇒ 冒烟必红**。

---

## 4. gate replay：新形态验证（`replay-gate.ps1` 现量）

直接加载 patched `build.ps1` AST 中的 `Invoke-ExternalProgram` 定义（确保完全一致），用 GUI stub 复测三态：

### REPLAY A: `STUB_SLEEP_MS=2000`, `STUB_EXIT=7` (doctor FAIL shape)

- 输出: `[A] elapsed=2075ms fail-recorder=(FAIL: wisp doctor reported FAIL (exit code 7).)`
- **[VERDICT A] RED CONFIRMED** — 等到完整延时且正确报告退出码 `7`。

### REPLAY B: `STUB_SLEEP_MS=0`, `STUB_EXIT=0` (doctor PASS shape)

- 输出: `[B] elapsed=15ms fail-recorder=()`
- **[VERDICT B] GREEN CONFIRMED** — 未误报 false positive。

### REPLAY C: exe missing

- 输出: `[C] fail-recorder=(FAIL: smoke test: 'wisp.exe doctor' could not be measured: program not found...)`
- **[VERDICT C] RED CONFIRMED** — 未测不可见性被拒绝为红，而非静默绿。

---

## 5. output visibility：`UseShellExecute=false` 不影响子进程流

Commit 注释声称"with UseShellExecute=$false, child's own lines still reach this caller's streams"（CI log 保留 doctor 输出）。验证：

**stubout.exe** 写入 stdout/stderr 后 exit 7：

```powershell
$ast=[System.Management.Automation.Language.Parser]::ParseFile(...)
$fn=$ast.FindAll(...,"Invoke-ExternalProgram",...)|Select-Object -First 1
Invoke-Expression $fn.Extent.Text               # define from build.ps1 AST
$env:STUB_SLEEP_MS='0'; $env:STUB_EXIT='7'
$r = Invoke-ExternalProgram -FilePath '.\stubout.exe' -ArgumentList @('doctor')
```

**输出**（来自子进程，显示在父进程 stream 中）:

```
stub: doctor line on stdout
stub: note line on stderr
```

**[VERDICT]** `output visible above + real code 7 = CI log keeps doctor report AND gate reads truth`.

---

## 6. 零回归检查（AC#3）

- **grep `build.ps1` 中任何剩余 `& ...wisp.exe`**: 无 HIT —— 旧形态已全部移除。
- **`$LASTEXITCODE` 残留读取**: 保留在读其他程序的行上（fetch-deps, npm ci/run build, go build, git commit），这些均为 CUI/bin 程序，等待语义自然满足，不构成同一风险类。
- **非 ASCII/emoji**: `§4` 仅出现在注释中，与现有 house style（line 4 "SPEC-11 §2.2", line 27 "§7.1"）一致；非 emoji。
- **parse**: 全文件 parse OK (370 lines)。
- **slo-check.ps1**: 未动（`git diff --name-only scripts/slo-check.ps1` = empty）。
- **cmd/wisp/**: 未动（不在本腿射程）。

---

## 7. AC 完成列表

| 目标 | 状态 | 证据 |
|---|---|---|
| ① `build.ps1:309-310` 改为真读退出码形态 | ✅ 完成 | lines 307–370 (above)，新函数 `Invoke-ExternalProgram` 复用 `slo-check.ps1` 形状 |
| ② stub 两态对比读数（旧陈值 vs 新真值） | ✅ 完成 | `proof-output.txt`: GUI OLD=0 (stale), GUI NEW=7 (real) |
| ③ zero regression: 无其它 `& wisp.exe`+auto var 直读 | ✅ 通过 | grep 扫描确认无遗留 HIT |
| ④纪律格（禁改范围外） | ✅ 通过 | `slo-check.ps1`, `cmd/wisp/**` 未动 |

---

## 8. 交件声明

- **提交路径**: `scripts/build.ps1`, `docs/evidence/s1/244r2-smoke-gate.md`（显式 pathspec commit）。
- **不 push**：仅 commit。
- **stub 临时件**: `D:\tmp\wisp\20261102\244r2\`（按纪律只建不删）。
- **读者**: 编排者收件，裁决表归口。

---

**附：关键读数**

- **GUI OLD**: `& stubgui.exe` → `elapsed≈21ms`, `$LASTEXITCODE=0` (stale)。
- **GUI NEW**: `Process.Start+WaitForExit` → `elapsed>2000ms`, `code=7` (real)。
- **REPLAY A**: `Fail "wisp doctor reported FAIL (exit code 7)"` (gate red when true fail)。
- **Output visibility**: child stdout/stderr lines visible in caller stream.
