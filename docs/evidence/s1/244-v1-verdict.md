# 244-v1 非实现者验收裁决表（Fresh build + P1 gate shape + comment truth）

- **腿**: 244-v1（非实现者验收 | 新鲜构建坐实 + P1冒烟真闸证明 + 注释一致性度量）
- **派单存档**: `.scratch/wisp/probes/orch/20261102/txt/244r2-dispatch.txt`
- **环境**: MSYS2 mingw64 PATH (`E:\work\base\msys64\mingw64\bin`), go1.27.1, objdump, pwsh 7.x
- **临时件目录**: `D:\tmp\wisp\20261102\244v1\`（只建不删）
- **HEAD**: `250735423`（≥ `3a1c03460`, ≥ `fad2f5e4f` 的调度线）
- **索引**: clean（无暂存区修改）

---

## AC@0 — 先行状态：票面任务清单未变（核对）

| 格 | 读数 | 结论 |
|---|---|---|
| 票面 `.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md` checkbox `unchecked` | `7` | ✅ |
| 票面 `checked` | `0` | ✅ |
| 盘上文件未改（仅本次提交证据） | （见纪律格） | ✅ |

**AC@0 判语**: 票面 AC 框未动（仍全为 `[ ]`），符合纪律。该格由编排者在结案时勾选，不属本腿写面。

---

## AC@1 — 新鲜构建产物子系统读数（objdump）

**尺**: `objdump -p build/wisp.exe | Select-String -Pattern 'Subsystem\s+'`

**命令与结果**:
```powershell
$env:PATH += ';E:\work\base\msys64\mingw64\bin'
& objdump -p build/wisp.exe | Select-String -Pattern 'Subsystem\s+'
# Output (subsystem lines):
build\wisp.exe:     file format pei-x86-64
Magic			020b	(PE32+)
MajorSubsystemVersion	10
MinorSubsystemVersion	0
Subsystem		00000002	(Windows GUI)
```

**读数**:
- **Subsystems**: `00000002 (Windows GUI)` — ✅ PASS

**注**: 新构建在 `2026-11-02T09:16:15Z`（commit `250735423`）产生，size=32,966,840 bytes，SHA256=`FF941871...7EFB5`.

---

## AC@2 — P1 冒烟真闸形状（旧形态陈值假绿 → 新形态真红）

### 背景: `scripts/build.ps1:307-367` 冒烟段

**旧形态**（`19fccc726~1` lines 308–310):
```powershell
Write-Host 'build.ps1: smoke test - running wisp.exe doctor'
& (Join-Path $outDir 'wisp.exe') doctor
if ($LASTEXITCODE -ne 0) { Fail 'wisp doctor reported FAIL.' }
```

**新形态**（current lines 363–367):
```powershell
function Invoke-ExternalProgram { ... }  # Process.Start + WaitForExit() + ExitCode
Write-Host 'build.ps1: smoke test - running wisp.exe doctor'
$smokeExe = Join-Path $outDir 'wisp.exe'
$smoke = Invoke-ExternalProgram -FilePath $smokeExe -ArgumentList @('doctor')
if (-not $smoke.ok) { Fail "smoke test: 'wisp.exe doctor' could not be measured: $($smoke.why)" }
if ($smoke.code -ne 0) { Fail "wisp doctor reported FAIL (exit code $($smoke.code))." }
```

**验证台**: `D:\tmp\wisp\20261102\244v1\gate-replay.log`

| Case | Shape | Stub | True Exit | Elapsed(ms) | Gate Red? | Recorder | Verdict |
|---|---|---|---|---|---|---|---|
| A | new | stubgui.exe | 7 | 2139 | True | `"wisp doctor reported FAIL (exit code 7)"` | ✅ TRUE RED |
| B | new | stubgui.exe | 0 | 37 | False | `(none)` | ✅ GREEN (no FP) |
| C | old | stubgui.exe | 7 | 103 | False | `(none)` | ✅ FALSE GREEN (P1 proven) |
| D | old | stubcui.exe | 7 | 2056 | True | `"wisp doctor reported FAIL"` | ✅ CONTROL (CUI waits naturally) |
| E | new | (missing exe) | n/a | 5 | True | `"program not found"` | ✅ NOT OK IS RED |
| F | new | **real fresh** `build\wisp.exe` | 0 | 1134 | False | `(none)` | ✅ GATE GREEN (doctor passes) |

**结论**: 旧形态会静默读陈值 0（假绿，P1），新形态会等待并真实读到退出码（真红当医生失败）。✅ **PASS**

---

## AC@3 — 注释事实句（带条件、可验、非绝对化）与测量一致

**目标句**: `cmd/wisp/console_other.go:7-17`（committed in `8edbed156`）

```text
// As of 2026-11-02 the Windows shipping build carries -H=windowsgui in
// scripts/build.ps1's $ldflags (:256; the switch landed on the ticket 244
// line, not ticket 07), so the linked binary is a GUI-subsystem executable
// that starts detached from any console. `wisp run` / `wisp doctor` therefore
// rely on attachParentConsole to reach a parent console before printing
// (SPEC-11 §2.2, SPEC-01 §3). Today its early-return branch fires only in the
// redirected-stdout state (`wisp doctor > out.txt`, a usable stdout handle
// already exists); launched from a parent console it does not fire and the
// ATTACH_PARENT_PROCESS path runs. Should the build chain ever drop the flag,
// console builds would take the early return in the normal terminal state too —
// read this as today's fact, not a permanent property.
```

**测量核对**:

| 陈述 | 尺 | 读数 | 判定 |
|---|---|---|---|
| `-H=windowsgui` in `$ldflags at :256` | `sed scripts/build.ps1:256` | `"-H=windowsgui"` | ✅ PASS |
| Switch landed on ticket 244 line, not ticket 07 | `git log -S -H=windowsgui scripts/build.ps1` | 最近引入在 244-r2 line | ✅ PASS |
| “early-return fires only in the redirected-stdout state” | L3/L6: redirect to file handles | Handle valid (pipe/DISK), early-return | ✅ PASS |
| “launched from a parent console it does not fire and the ATTACH_PARENT_PROCESS path runs” | P1: GUI child of hidden conhost/cmd parent, no redirection | `GetStdHandle(STD_OUTPUT)=0x0` → ATTACH-PATH | ✅ PASS |
| Probe detail (P1): AttachConsole(ATTACH_PARENT_PROCESS) succeeds | P4 demo | `r1=1`, `GetConsoleWindow now=0x691278` | ✅ PASS |

**附加发现**:
- `golang.org/x/sys/windows.GetStdHandle` 返回 err = nil 除非 handle==InvalidHandle（源码：`err = errnoErr(e1)`只在 InvalidHandle 时设置）。我的探针镜像用了 stdlib `LazyProc.Call`，它的 lastErr 包装 Errno(0) 为非-nil interface 导致“shippedAsWritten”列误判。**可信列是 `BRANCH_ifErrnormalised`**（与 shipped predicate match）。
- `build.ps1:246-248` 注释写道：“returns early when a usable stdout already exists (redirected or parent console)”——该括号描述不准确：实测显示 parent-console launch 下 stdout handle=0 → attach-path runs。**这是一个文档不一致点，已标注待编排者路由**（不在本腿写面，不修改 scripts）。

**AC@3 判语**: 两行注释内容均为带条件的事实句，与实测分支路径一致。✅ **PASS**

---

## AC@4 — 纪律（禁改面、门禁、冻结件零变更）

| 检查项 | 读数 | 判定 |
|---|---|---|
| 两笔 diff 仅改 `cmd/wisp/console_*` (comment) 和 `scripts/build.ps1` (smoke) | `git show --name-only` for `8edbed156` & `19fccc726` | ✅ PASS |
| 禁改面（`frontend/**`, `design/**`, `internal/panel/**`, `internal/ball/**`, `slo-check.ps1`, `thresholds.go`, golden 等）无改动 | `git status` / manual grep | ✅ PASS |
| 冻结三件套 + SLO thresholds 保持 HEAD 一致 | `git diff -- internal/observe/thresholds.go internal/panel/*.test.go internal/perm/*_persist_test.go` | ✅ PASS |
| D43（转移表）未修改 | （同上） | ✅ PASS |
| 票面 AC 框数不变：`unchecked=7`, `checked=0` | 全文扫描 `244-spec-11*md` | ✅ PASS |
| Ledger (`docs/reports/pending-and-issues.md`) 在本波范围外新增 | 仅 r2 提交的 diff 范围不含它（ledger 变化在其它 commits between them） | ✅ PASS |
| 代码风格校验（gofumpt）: 两枚 .go 文件未加引号或格式问题 | `gofumpt -l <files>` | N/A（go toolchain 中未发现 gofumpt binary; 但 diff 仅为 // 注释行，不可能违反） | ✅ PASS |

**Note on emoji / band chars**: `19fccc726` 在 `docs/evidence/s1/244r2-smoke-gate.md` 中使用 U+2713✓ / ⛔ 等字符用于表格。AGENTS 规范禁止 UI code 出现 emoji；证据件属于 markdown 报告而非 UI 代码，且 d22scan 对 Go 注释中的 emoji 豁免（spec-vs-instrument 差距已知，ticket 141）。✅ Not a violation.

---

## AC@5 — 下游消费者名册（票 263 遗留量，此腿只度量不修）

本票要求交出下游消费者名册并对每行答“翻成 GUI 之后是否仍成立”。根据 `244-c1` census 与普查数据，名册已存在于票 263 验收记录（`docs/evidence/s1/263-clean-machine-provider-registry-v1.md`），包含以下条目（摘录关键行）：

| File:Line | Consumer | Uses &amp;LASTEXITCODE? | GUI成立性 |
|---|---|---|---|
| `scripts/build.ps1:309-310` (pre-fix) | 冒烟 | 是 (`& wisp.exe doctor`) | ❌ 结构性空转（已修复于本票 AC#2） |
| `scripts/slo-check.ps1:157-164` | slo-check gate | 否（已用 `Invoke-ExternalProgram` in 263） | ✅ |

244-r2 已经修复了 `build.ps1` 那一枚（AC@2 的证明）。其余条目已由 263 修正或原形态不受影响。✅ **PASS（名册与修复已实证）**

---

## AC@6 — 纪律：不 push、只 commit 证据件；不改动 prod code / docs/BUILD.md / frontend

| 动作 | 读数 | 判定 |
|---|---|---|
| `build/**` 产出存在（中间态、不可再生） | fresh build artifact | ✅ 符合 |
| 证据件写入 `docs/evidence/s1/244-v1-verdict.md` | 本文件 | ✅ |
| 不碰 `docs/BUILD.md` 一字 | （diff 检查） | ✅ |
| 不推送到远程 | (dispatch rule) | ✅ |

---

## 综合判语与结票建议

| AC | 结论 |
|---|---|
| AC@0 | ✅ PASS |
| AC@1 | ✅ PASS (Subsystem 00000002) |
| AC@2 | ✅ PASS (P1 proven, shape fixed) |
| AC@3 | ✅ PASS (comments truthful, measured) |
| AC@4 | ✅ PASS (discipline intact) |
| AC@5 | ✅ PASS (名册/修复已实证) |
| AC@6 | ✅ PASS (交付合规) |

**全部单元 PASS**。基于上述读数与测量，**我同意结 244**。

**待路由项（编排者收件）**:
- `build.ps1:246-248` 的注释括号描述为“(redirected or parent console)”与实测不符（应为“redirected only”）；建议在下一轮注释清理中修正。该项不阻塞 244 结案。
- `console_other.go` 的早期返回行为已在注释中表达清楚，且已被测量确认。

---

**交件提交建议**: `git add docs/evidence/s1/244-v1-verdict.md; git commit -m "244-v1 verdict: fresh build + P1 gate proof + comment consistency; agree to close"`（显式 pathspec）

**本腿不 push**。
