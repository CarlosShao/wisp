# 70 — 让 CI 真的成为护栏：5 个 job 全红，逐因分诊（本仓的"门禁"从未生效过）

**Status:** in-progress
**Claimed by:** ~~agent-ticket70~~（174 次调用撞 turn 上限而死，**票面一条 log 都没写**）→ agent-ticket70-c（接续，只做 AC#2/AC#4/AC#6）
**Last update:** 2026-09-21 09:15（编排者：从 `git log` 重建断点，见下面那块）
**Blocked by:** ~~66、68~~ **两条都已解除**：票 66 已闭（`-done`），票 68 是 `blocked-on-owner` 且 `internal/ball`/`cmd/balldebug` 已让出

> ## 编排者重建的断点（09:15，我逐条亲自验过，**接续代理不要重做**）
> 死掉那位**5 个 commit 全在树上、票面账目为零**，所以这里按 commit 重建：
> - **AC#1 本地半** ✅：`f342413`（我入库前审过：`gofmt -l` 空、`go build ./...` rc=0、`go vet ./...` rc=0、
>   `git diff -w` 只剩 literal 展开）+ `8bfd47d`（它自己补上 gofumpt **版本漂移**漏掉的 1 个文件）。
>   **另一半"CI 的 lint job 真绿"仍要看一次真 run**，不许用"本地这步过了"替代。
> - **AC#5（R16 五条）** ✅ **全部落地且我核过**：`3539d47` 把 `cmd/balldebug` 三处具名 spawn 改走
>   `observe.Registry.Spawn`（HEAD 上裸 `go` 语句 **0** 处；名册名 `balldebug-hotkey-bridge`/
>   `balldebug-level-feeder` **不借产品名**）；`38b3715` 把 ban #1 从 `go func(` 扩到 `go <任意>`；
>   `allowlist.txt` 非注释行 4 → **5**，多出的那条**正是 R16#1 唯一授权的文件级豁免**
>   （`internal/observe/goroutine.go`，按路径不按调用形状）。
>   ⚠ **R16#3 的仪器没有回退**——我本来要查 `N==0` 致命是否还在，结果发现它**比裁定更强**：
>   `main.go:509` 是"生产 `.go` 文件数低于阈值 ⇒ 报『the ban scan would be a no-op』并失败"，
>   `:536` 仍自报 `examined N production Go files`，`go test ./tools/d22scan/` 8 条用例 `ok 0.315s`
>   （含 `TestScanAloneIsNotAFalsifier`、`TestCheckRootRejectsBlindRoots` 两条反空跑用例）。
> - **AC#3 → 移交票 72**（`a04d3e2` 定性为**安全分类失效**，不是占位；那张票在
>   `.scratch/wisp/issues/72-atble-classification-runner.md`）。**本票不要再碰 `internal/risk/`。**
> - **AC#4 配置已落**（`98fa8ae`：`slo-full` 缺 gcc 的根因是 runner 长驻进程拿的是 09-19 的旧环境，
>   msys64 在 E 盘且只在用户 PATH；修法 = 给该 job 设它自己文档化的 `MINGW64_ROOT` 旋钮），
>   **只剩"看一次真 run 的结果"这一步**。
> - 剩下真正没动的只有 **AC#2 `test-core` 分诊**与 **AC#6 逐 job 全绿证据**。

**Parallel slots:** ≤1 sub-agent；AC#1 落地期间**不得有第二个代理写码**
**Spec refs:** D22（"no job skippable"）、D42#9、C19、SPEC-08 §2（冻结）
**登记项:** A26（vacuous 静态门 + ban #1 覆盖洞）、A27（5/5 job 红）、A14/A15 的 CI 侧后果

## 为什么要这张票

我把 CI 当既有护栏写进了各票简报（"CI 会拦住裸 goroutine / 墙钟超时 / 越权 Clean"）。
2026-09-20 直接查 `gh run view 35517463335 --json jobs`（commit `0891907`）的结果是：
**`lint` / `test-core` / `test-windows` / `slo-smoke` / `slo-full` 五个 job 全部 `failure`**，
自 `f088ce3` 引入流水线以来**没有一次通过的证据**。⇒ **所有以"CI 会挡"为前提的论证都缺前提。**

逐 job 失败步骤（实测，非推断）：

| job | 失败步骤 | 我已核到的因 |
|---|---|---|
| `lint` | `gofmt (gofumpt)` | 用 CI 钉的 **v0.7.0** 本地复跑 CI 那条逐字命令 ⇒ **69 个文件**被标 |
| `test-core` | `Portable package tests` | `internal/observe` 的 `TestNoBareGoFuncInProductionCode`、`TestSampleStateCPUTotalDrivenMean`；另 `TestExistsAndDelete`、`TestBlobsListsMetadataOnly` |
| `test-windows` | `PathResolver junction placeholder` | 步骤名自带 **placeholder** ⇒ 可能是"设计上先红"的占位，**先定性再动** |
| `slo-smoke` | `SLO smoke gate` | 即 A14（解析丢末项常态 fail-closed），**票 66 修，不在本票重复** |
| `slo-full` | `Build wisp.exe` | 自托管 runner `wisp-slo` 的依赖/构建，未查 |

## 验收标准

- [ ] **AC#1 格式化 sweep（机械、零语义）**：`gofumpt -w` 过 CI 那条命令列出的全部路径，
  直到 `gofumpt -l . tools/d22scan tools/mockllm` **输出为空**。
  ⚠ **单独一个 commit**，message 里写明"纯格式化、无行为改动"；
  **不许**混进任何逻辑修改，也不许用"把 gofumpt 步骤改成 `continue-on-error`"糊过去
  （D22 明写 no job skippable，**放宽门不是修门**）。
  完成后复跑 `go test -count=1 ./...`（可跑的包）证明没改坏。
- [ ] **AC#2 `test-core` 分诊**：在 **Linux** 上复现那 4 条（CI runner 就是 ubuntu-latest），
  逐条定性为「平台差 / 票 66 或 67 带出的回归 / 本就在坏的断言」。
  **禁止**为了让它绿而改断言或加 skip；平台差异要落成**显式 build tag 或平台专属期望值**，并说明为什么。
- [ ] **AC#3 `test-windows` 的 placeholder 定性**：查清它是"待票 18/20 真实用例的占位"还是"真失败"。
  若确为占位，**改成显式的 TODO 步骤且**在 PLAN/HANDOVER 里留痕（改法由我定，代理先给证据）。
- [ ] **AC#4 `slo-full` runner**：查 `Build wisp.exe` 在自托管 runner 上为何失败（deps 缓存？PATH？CGO？），
  交付**可复跑的修复步骤**或"此 runner 今天不可用"的明确结论 + 需要 owner 做什么。
- [ ] **AC#5 ban #1 覆盖洞（A26 剩余半）—— 按裁定 R16 执行，不要再问**：
  匹配器 `tools/d22scan/main.go:251-252` 只对字面量 `go func(` 报警。我（编排者）逐条核过它说的"4 处具名协程"，
  **更正为 3 处要处理 + 1 处是机制本身**：`internal/observe/goroutine.go:281` 的 `go r.run(...)`
  **就是 Registry 的内部实现**，不是漏网；要改的是 `cmd/balldebug/main.go:231`（`runHotkeyBridge`）
  与 `:259`、`:440`（`feedLevels` ×2）。
  1. 新增 ban 匹配 `go <任意>`（**含具名调用**），**唯一按文件豁免的是 `internal/observe/goroutine.go`**；
     ⚠ **不得**以"这条是具名调用"为豁免理由——那正好把漏洞留在门上。
  2. `cmd/balldebug` 那 3 处**改走 `observe.Registry.Spawn`**（不是加豁免），让它们进 D38b 名册、被泄漏检查看见。
  3. 输出必须继续自报 `examined N production Go files`，**`N==0` 致命退出**（票 67 的 `checkRoot()` 已实现，**不许回退**）。
  4. **只许从严**：扩展后若又挖出具名 spawn，逐个登记进本票 log 由我判定，**不许默默加豁免凑绿**。
  5. 连带 loose end：`Spawn` 对不在名册的名字无条件 WARN，票 67 的 mockllm 测试因此每条多一行
     `WARN goroutine outside the D38 roster`。抑制它要在 `internal/observe/goroutine.go` 加 `TemporaryNames`
     ——**那是票 66 的文件，等票 66 收尾再做**；借用现有名册名来骗过泄漏检测**禁止**（代理拒得对）。
- [ ] **AC#6 全绿证据**：`gh run view <新 run> --json jobs` 的**逐 job 结论**贴进本票 Progress log。
  判据是**五个 job 全 pass**，不是我某一步本地过了。**`slo-smoke` 若仍红，先确认票 66 是否已合**，
  不要为了让它绿而调采样阈值（D32 的 0.5% 一字不动）。

## 编排者已裁定

1. **不放宽任何一步 CI**：不加 `continue-on-error`、不删步骤、不把阈值往下调。
2. **格式化与逻辑改动永不同 commit**（AC#1 单独成 commit）。
3. 本票 AC#1 落地期间**只有它一个写码代理**，避免全仓重写与别人 WIP 相撞。

## Progress log（append-only）
- [2026-09-20T15:39:24Z] agent=agent-ticket70 did=claim ticket70 (Status in-progress / Claimed by agent-ticket70); measured AC#1 baseline with the pinned version: gofumpt v0.7.0 -l . tools/d22scan tools/mockllm => 69 files (matches the orchestrator's count verbatim). Tree note: docs/SLO.md carries another agent's WIP edit (7+/1-) and is excluded from every commit I make (explicit paths only).
- [2026-09-21T01:33:20Z] agent=agent-ticket70-c did=**AC#2 分诊 `test-core`**（接手前任的两个死前发现，票面正文没写的那些）。复现方式＝**与 CI 那一行逐字相同的命令**跑在 **Linux** 上：`docker run golang:1.27`（GO 1.27.1 linux/amd64、`WISP_ENV=test`、挂载本仓 HEAD），包清单与 `.github/workflows/ci.yml:101-106` 一字不动。**先纠正 A27 的量级**：票面写"约 4 条测试红"，实测（run 35549416856 + 我本地各一次）**4 个包红、35 条以上 FAIL、其中 `internal/tools` 是 10 分钟超时 panic**（CI `FAIL internal/tools 600.027s` / 本地 `600.085s`，两边同一形状）。逐族定性与处置：
  (1) **`internal/observe` 2 条 = 本就在坏的断言（不是平台差，Windows 本地同样红）**。`TestRollingWriterDayRoll`：`newRollingWriter` 按文档**故意**在构造时 eager-open（logging.go:185「so a broken log dir fails at init」），用的是**真实时钟**，而 fixture 把 `w.now` 换成 2026-09-19/20 —— 于是在 19/20 号那两天真实 boot 文件与 fake 文件同名、`len(names)==2` 侥幸成立，**过了 21 号就永远多一个文件**（实测两侧都得 `[20260919,20260920,20260921]`）。`TestSampleStateCPUTotalDrivenMean`：fixture 写死"0.5 core-second per second"，而 `sampler.go:314` 是 `cores := runtime.NumCPU()` 除_all-core_——12 核本箱 4.17% 通过，CI 4 核 runner 12.5%（实测 14.846%）撞 `>10` 上界。**处置＝只动实现与 fixture**：给 `newRollingWriterClock(dir,sizeMB,days,now)` 开一条时钟注入缝（`newRollingWriter` 原样委托，生产行为零改动），fixture 改按核数换算成"恒为全核容量的 4.166%"。**两条断言一字未改**（`len(names)!=2`、`<=0 || >10` 原样）。Linux 与 Windows 两侧现均 ok。
  (2) **`internal/secret` 7 条 = 平台差**（`Store()/Resolve()` 走 DPAPI，`protect_other.go` 是按 C28 写死的 fail-closed stub，非 Windows 只能报 "DPAPI is only available on Windows"）。**处置＝AC#2 明给的"落成显式 build tag"**：`store_test.go`/`delete_test.go`/`migrate_test.go` 加 `//go:build windows` + 就地说明。**不是 skip、没删断言**：这三档文件里 15 条测试全部在 `test-windows` 那一步（`go test ./internal/secret/`）继续跑；portable 面（env: ref / redact / configrefs）留在**未加 tag 的 `refs_test.go`**，ubuntu 仍在测。**如实记代价**：8 条本就与 DPAPI 无关的测试因此退出 ubuntu 覆盖，要更细的拆法（把 portable 条单拎成 `*_other_test.go`）请编排者点头，我不在这轮擅自扩大改动面。
  (3) **`internal/risk` 17 条 = 冻结区，未碰，上报**。失败文字全部落在 P12 sync-root 探测与 provenance 写门上（`source "registry" must count as confirmed`、`env-grade roots are confirmed evidence`、`path outside the user profile must not be suspect`），而 `syncdirs_other.go` 自带 `DEFERRED(P12-macos)`、`registryProbe` 在非 Windows 直接返回 nil。判据在 `internal/risk/` 里 ⇒ 按简报"涉及改 `internal/risk/` 停下来上报"，**这条与票 72 同区**。
  (4) **`internal/tools` 19 条 + 超时 = 同一个平台差的不同下游，且是票 20 在飞的包，未碰，上报**。根因证据：CI 与我本地都打出 `open /home/runner/work/wisp/wisp/internal/tools/\tmp\TestFSRead…/note.txt` —— **C26 的 canonical 形状在 Linux 上仍是反斜杠**（`pathresolver.go:133` `strings.ReplaceAll(p, "/", "\\")`，`pathresolver_other.go` 顶部写明 `DEFERRED(macOS/Linux)`），tools 忠实照抄 canonical 去开文件，于是 404；超时那半（`wiring_test.go` 卡在 `approval.Gate.PendingApproval` 4m50s、整包 600s panic）**尚未定因**，可能是同一条路径故障的下游，也可能是审批交接在 Linux 上的第二个洞——这属票 20 的判据，我不替它下结论。
  next=AC#4 看一次真 run 的 `slo-full` 结果（前任已落 `98fa8ae` 的 `MINGW64_ROOT`），再 AC#6 贴逐 job 结论；`test-core` 本轮**闭不了**，卡在 (3)(4) 两族的冻结区判断。

