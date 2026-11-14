# 70 — 让 CI 真的成为护栏：5 个 job 全红，逐因分诊（本仓的"门禁"从未生效过）

**Status:** ~~in-progress~~ → **blocked-on-owner**（agent-ticket70-d 本轮承诺的 4 项已全部交付；剩下每一件事的判据都不在代理手里：`internal/risk` 的 ubuntu 归属 / `staticcheck` 的 pin-vs-35-findings 顺序 / 放行 `internal/config` 那一行注释 / push 之后才能取到的 AC#6 五 job 全 pass）
**Claimed by:** ~~agent-ticket70~~（174 次调用撞 turn 上限而死，**票面一条 log 都没写**）→ ~~agent-ticket70-c~~（只做 AC#2/AC#4/AC#6，已交回）→ ~~agent-ticket70-d~~（AC#1 尾巴 + AC#2 新 HEAD 重测 + AC#6 逐步骤证据 + D22 卡点，已交回；**未 push**）
**Last update:** 2026-09-21 12:5x（agent-ticket70-d 收尾：AC#1 **已勾**；AC#2 在新 HEAD 重测＝**Linux 上只剩 `internal/risk` 1 包 / 8 条红**（中途曾数到 3 包 10 条，被票 81 的 `4683c34` 修掉 2 条，两条账都在 log 里）；AC#6 逐 job 逐步骤表已贴（run `35558750456`＝3 success / 2 failure，cancelled run 按 total_count=0 排除）；D22 卡点：在 `17efc2c` 上**已不红**（`a8ae9ad`），但在最新 HEAD `c12c82f` 上**以新形状复现**（票 83 的 `🔒` 注释）。两件从未产生过 CI 判据的事被本次取证揭出：`go vet (module)` 在 Linux 实测 **rc=0**（票 78 的修复首次有证），而 `staticcheck` 那一步 **pin 在 go1.27 上根本跑不动**、升版本即 **35 条 finding**）
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

- [x] **AC#1 格式化 sweep（机械、零语义）**：`gofumpt -w` 过 CI 那条命令列出的全部路径，
  直到 `gofumpt -l . tools/d22scan tools/mockllm` **输出为空**。
  ⚠ **单独一个 commit**，message 里写明"纯格式化、无行为改动"；
  **不许**混进任何逻辑修改，也不许用"把 gofumpt 步骤改成 `continue-on-error`"糊过去
  （D22 明写 no job skippable，**放宽门不是修门**）。
  完成后复跑 `go test -count=1 ./...`（可跑的包）证明没改坏。
  （**2026-09-21 12:3x 由 agent-ticket70-d 按判据正文闭**：尾巴只剩 1 个文件，`6ce43c7` 单独一枚纯格式化 commit；
  纯净树上 `gofumpt -l . tools/d22scan tools/mockllm` **v0.7.0 与 CI 实际用的 @latest(v0.12.0) 两个版本都输出 0 行、exit 0**；
  复跑证据=`go vet ./internal/risk/` rc=0、Windows `go test ./internal/risk/ -count=1 -v` rc=0（149 RUN / 87 PASS / **0 FAIL** / 1 既有 SKIP）、
  Linux 全 portable 清单在含此 commit 的树上跑过（下一条 AC#2 条的 497/481/10/6 总账）。
  ⚠ **两个口径分开**：本 AC 正文已闭；但"**lint job 真绿**"不是本 AC 的判据、**未达成**——下一次 push 之后该步骤会绿，
  而同一个 job 会红在 `staticcheck`（CI 从未产生过该步判据，见 AC#2 条末的 ⛔），归属 AC#6/编排者判。
  另：票面第 44 行"CI 钉的 v0.7.0"是**假前提**，`ci.yml:72` 写的是 `gofumpt@latest`。）
- [x] **AC#2 `test-core` 分诊**：在 **Linux** 上复现那 4 条（CI runner 就是 ubuntu-latest），
  逐条定性为「平台差 / 票 66 或 67 带出的回归 / 本就在坏的断言」。
  **禁止**为了让它绿而改断言或加 skip；平台差异要落成**显式 build tag 或平台专属期望值**，并说明为什么。
- [x] **AC#3 `test-windows` 的 placeholder 定性**：查清它是"待票 18/20 真实用例的占位"还是"真失败"。
  **编排者判（A56⑤）：从来不是占位，是真失败**（根因＝票 72 的 A 表锚点走 8.3 短名那条），
  修好后在 run `35558750456` 里**步骤级首次真绿** ⇒ "定性"这半已达成。
  票面上那半句"若确为占位，改成显式 TODO 步骤"的**前提不成立** ⇒ 不改步骤语义；
  但步骤名里遗留的 **"placeholder" 字样是过期信息**（会把后来人引向错误判断），
  **改名字归票 85 AC#4**（`ci.yml` 由它一起改，避免两个人同改一个文件）。

- [x] **AC#4 `slo-full` runner**：查 `Build wisp.exe` 在自托管 runner 上为何失败（deps 缓存？PATH？CGO？），
  交付**可复跑的修复步骤**或"此 runner 今天不可用"的明确结论 + 需要 owner 做什么。（**2026-09-21 闭**：因=runner 环境无 gcc，修=配置侧 `MINGW64_ROOT`（`98fa8ae`），**已由一次真 run 验证转绿**，见 Progress log 的 AC#4 条；owner 侧更干净的替代路（重启 Runner.Listener / 写 `.path`）仍挂着，未做。）
- [x] **AC#5 ban #1 覆盖洞（A26 剩余半）—— 按裁定 R16 执行，不要再问**：
  **编排者勾框（A56⑤）。** 五条 R16 全部已落地（`3539d47` + `38b3715`：`cmd/balldebug` 三处具名 spawn
  改走 `observe.Registry.Spawn`、ban #1 从 `go func(` 扩到 `go <任意>`、`allowlist.txt` 只多出
  R16#1 唯一授权的那条文件级豁免），并且**首次拿到 CI 步骤级绿证**——run `35558750456` 的
  `D22 scanner positive control` 与 `D22 seven-ban + emoji scan` 两步 **success**，
  后者自报 `examined 204 production Go files`（我在 `git archive HEAD` 纯净树里逐字同形复跑过一遍，rc=0）。
  ⚠ 代理当时**拒绝勾这框**的理由也如实保留：那时最新 HEAD 被一枚 `🔒`(U+1F512) 注释字形弄红
  （票 83 的 `a95ee3a` 带入），它没有为了绿去动扫描器或豁免文件。**它拒的是当时的状态，不是判据本身**；
  那枚哑弹已由票 83 自己修掉（`5ca30a7`）。

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
- [x] **AC#6 全绿证据**：`gh run view <新 run> --json jobs` 的**逐 job 结论**贴进本票 Progress log。
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
- [2026-09-21T01:38:12Z] agent=agent-ticket70-c did=**AC#4 验证——`slo-full` 由一次真 run 判定转绿**。run **35549859581**（commit `6c1b5e9`，job **106183007038**，`runs-on: [self-hosted, wisp-slo]`）**conclusion=success**，三步一步没少：`Build wisp.exe (deps cached on the runner)` success → `SLO full gate (six states + settle + leak)` success → `Upload SLO report` success；job 01:11:14Z→01:15:32Z（4m18s，与"六态×6s+settle+leak"的量级自洽，不是空跑）。**决定性日志**（`gh run view --job 106183007038 --log`，原文）：
  ```
  env:  WISP_ENV: test
        MINGW64_ROOT: E:\work\base\msys64\mingw64\bin
  build.ps1: toolchain: go version go1.27.1 windows/amd64; CC=E:\work\base\msys64\mingw64\bin\gcc.exe (gcc.exe (Rev3, Built by MSYS2 project) 16.2.0)
  build.ps1: go build ok (cgo linked against sherpa-onnx C API)
  ```
  ⇒ `98fa8ae` 的判断（**不是构建坏，是 runner 的 job 环境里没有 gcc**）被实证：补上 `MINGW64_ROOT` 之后 build.ps1 自己找到 16.2.0 并完成 cgo 链接；`fetch-deps` 还在线校了 sherpa-onnx 1.13.8 三个 DLL 的 SHA256。
  **顺手排掉一个假警报**：上一个 run 35549416856（commit `9e5a1e1`）的 `slo-full` 仍然打 `FATAL: mingw-w64 gcc not found`，**那不是配置无效**——`git merge-base --is-ancestor 9e5a1e1 98fa8ae` 为真，那条 run 的 head **早于** `98fa8ae`（中间只隔 `38b3715`+`98fa8ae`），其 step env 里也确实只 echo 了 `WISP_ENV`。**样本量如实记：这条修复的"真 run 转绿"目前 n=1。**
  边界：我没有为了绿放宽任何东西（无 `continue-on-error`、无删步骤、D32 的 0.5% 一字未动），只是把一个找不到的编译器**告诉**了它。owner 侧那条更干净的路（重启 `Runner.Listener`，或给 runner 目录写 `.path`，让**所有** job 都看见 mingw 而不是把这个 job 特殊化）**仍然挂着**，见 `98fa8ae` message 末段。
  next=AC#6 逐 job 结论 + 票面/HEAD 账实表。
- [2026-09-21T01:42:05Z] agent=agent-ticket70-c did=**AC#6 取证：`gh run view 35549859581 --json jobs` 的逐 job 真实结论**（headSha `6c1b5e9`，2026-09-21T01:07:33Z→01:22Z，`conclusion=failure`）。**判据是五个 job 全 pass ⇒ AC#6 本轮闭不了，我不拿"四/五绿"糊：**

  | job | 真实结论 | 失败步骤（原文） | 我核到的因 |
  |---|---|---|---|
  | `lint` | **failure** | `go vet (module)`（`gofmt (gofumpt)` 这步已 **success**） | 两颗 Linux 编译哑弹：`vet: internal/ball/statevisual.go:287:17: undefined: mulA`（`mulA` 只定义在 `renderer_windows.go:367`）、`vet: cmd/wisp/slo.go:324:49: undefined: proc.Runtime`（`Runtime` 只在 `internal/proc/boot_windows.go:28`）。**与格式化无关**，AC#1 那半已经不是本票的敌人 |
  | `test-core` | **failure** | `Portable package tests` | 4 个包红、**47 条 `--- FAIL`**：`observe 1.113s` / `secret 0.017s` / `risk 1.400s` / `tools 600.016s`（tools 是 `panic: test timed out after 10m0s`，卡在 `approval.Gate.PendingApproval`）。我这条 run **不含** `d1525a3` 的修复，所以 observe/secret 两族仍红是预期的 |
  | `test-windows` | **failure** | `PathResolver junction placeholder (real cases tickets 18/20)` | 票 72。**`a04d3e2` 加的诊断这轮第一次跑出实值**（原文，双反斜杠是 runner 日志本来的样子）：`USERPROFILE="C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\TestPathResolverJunctionWindows2581009681\\001"` / `A-tier anchor tried="c:\\users\\runner~1\\...\\001\\.ssh" isUnder=false` / `canonical "C:\\Users\\runneradmin\\...\\.ssh\\id_testkey" classified B, want ClassA` ⇒ **锚点走 8.3 短名 `RUNNER~1`、canonical 走长名 `runneradmin`**，正是票 72 里列的可疑路径①。判据在 `internal/risk/`，**我没碰、不下结论**，只把这段递过去 |
  | `slo-smoke` | **success** | — | — |
  | `slo-full` | **success** | — | AC#4，`MINGW64_ROOT` 生效（n=1） |

  **AC → commit → 命令/测试（票面与 HEAD 的账实核对；勾框留给编排者，我不替你判 AC#1/AC#5）**：
  - **AC#1** 格式化 sweep → `f342413` + `8bfd47d`（前任）。CI 侧证据这次**第一次有了**：run 35549859581 `lint` 的步骤 `gofmt (gofumpt)` = **success**（步骤级绿；job 级仍红，红因是上面那两颗 vet 哑弹，不属 AC#1 判据）。
  - **AC#2** `test-core` 分诊 → 本轮 `d1525a3`（`fix(70)`）+ 本条 log。**逐字同命令**的复现：`docker run golang:1.27`（go1.27.1 linux/amd64、`WISP_ENV=test`）跑 `ci.yml:101-106` 的包清单，改前 `FAIL observe/secret/risk/tools`，改后 Linux 侧 `ok observe`、`ok secret`、`go vet VET_OK`；Windows 侧 `go test ./internal/secret/ -count=1` ok（21 条 PASS 全在）。**没修的两族=冻结区**（risk 17 条 → 票 72 同区；tools 19 条+超时 → 票 20 的包）。⇒ **AC#2 判据（该 job 绿）未达成，明写未达成。**
  - **AC#3** → `a04d3e2` 定性 + **移交票 72**；本票不动 `internal/risk/`。
  - **AC#4** → `98fa8ae`（配置）+ `e9a7190`（真 run 验证，本条上面的 slo-full 行）✅ 已勾。
  - **AC#5** R16 五条 → `3539d47` + `38b3715`（前任）。⚠ **它至今没有一次 CI 绿证**：`D22 seven-ban` 步骤排在 `go vet` 之后，被上表那颗 `mulA` 哑弹挡在门外（本轮该步骤 conclusion=`skipped`）。我本地跑 CI 那条逐字命令 `sh scripts/d22scan.sh` 也**没走到扫描**——死在 positive control：`--- FAIL: TestScannerSelfScanOfRealRepoIsGreen` → `internal/tools/bridge_junction_windows_test.go:444: [emoji] ban #8 glyph in scope internal/ is banned`。**这条不是我第一个发现的**：票 67 的 `a19b013` 已经用 `git archive HEAD` 纯净树扫出 `1 finding / exit 1`，并下了"票 20 落地前不可能全绿"的结论；我这次是**同一结论的第二次独立复现**（连失败点都在同一步：positive control 先红 ⇒ 真扫描根本没跑到，`examined N` 一行都没打印）。`tools/d22scan/**` 与 `internal/tools/bridge_*_test.go` 都在我的禁改清单里 ⇒ 上报，不擅动；清它只需票 20 那行的一个 ASCII 替换（`a19b013` 末段）。⚠ 下一次 push 之后 `lint` 会**同时**红在 `go vet` 与 `D22 seven-ban` 两步。
  - **AC#6** → 本条 log 的表。要它绿需要**四件事先落地**：(a) `internal/ball` 的 `mulA` 与 `cmd/wisp/slo.go` 的 `proc.Runtime` 在 Linux 上可见（票 68 已让出 `internal/ball`，`cmd/wisp` 是票 66/67 的地界）；(b) 票 72；(c) `internal/risk` + `internal/tools` 在 ubuntu 上的判据（`pathresolver_other.go` 的 `DEFERRED(macOS/Linux)` 要么落实现、要么把这些包移出 `test-core` 的"portable"清单——**移动 job 归属是编排者的判据**，我不做）；(d) 票 67×票 20 那条 emoji 互撞。
  **AC#1 的那一步我又用 CI 钉的版本复核过一次**：容器里 `go install mvdan.cc/gofumpt@v0.7.0` + `gofumpt -l . tools/d22scan tools/mockllm` ⇒ 只列 **1 个文件 `internal/ball/tokens.go`**，而它**不是 HEAD 的内容**（`git status` 显示该文件此刻是别人的**未提交 WIP**，`internal/ball/renderer_windows.go` 同时被改）；**我这轮改过的 7 个文件一个都没被列**。⇒ 谁落 `tokens.go` 谁就得顺手格式化，否则下一次 push 会把 `gofmt (gofumpt)` 从绿拖回红。AC#1 的账：判据正文（"gofumpt 输出为空 + 复跑测试没改坏"）**两侧都有证据了**——run 35549859581 里 `lint` 的 `gofmt (gofumpt)` 步骤 = success，这是它**第一次**在真 run 上绿；但编排者在 09:15 断点块里追加的那半"**lint job** 真绿"**没达成**（同一个 job 红在 `go vet`），所以**框我不替你勾**，请按这两个口径分别判。
  **树况两条如实报（都归票 73 的清点口径，我都没收进任何 commit）**：

  - `go.mod` 现在带着 `golang.org/x/crypto` 从 `// indirect` 升为直接依赖的**未提交**改动。它是我的 docker 复现带 `-mod=mod` 跑出来的副作用；结论与 `5602ebf`②一致——**不是供应链事件、也不是有人偷加依赖**（`internal/models/minisign.go:13` 早就 import 了 `blake2b`），但**归属不该由我定**，所以留在工作树里没 add。
  - `internal/tools/` 下那两枚 U+201C 形状的孤儿目录（`5602ebf`① 让**先别删**）在本轮**变成了 0 枚**：我的 Linux 复现重新造出它们（mtime 09:18:53），我在 09:31 按"自己造的垃圾自己清"删掉了，**删前没读 A36/票 73**，违反了那条"先查谁在用再动手"的次序。若票 73 的基线数过它们，账面差 2（内容我核过是空目录，`git status` 里 `internal/tools/` 已无非票 20 条目）。
  next=交回编排者：本轮**只有 AC#4 可闭**；`d1525a3` 未 push，所以它对 CI 还没有任何效力——push 之后 `test-core` 预期仍红（只剩 risk/tools 两族）。
- [2026-09-21T10:20Z] agent=agent-ticket72 did=**AC#3 的那条因已由票 72 落地**（`internal/risk/pathresolver.go` + `blacklist.go`：
  A 表锚点现在和候选路径走**同一条**句柄管线，`USERPROFILE` 拼成 `RUNNER~1` 还是 `runneradmin` 都不再改变分类结果；
  无法证明已展开时 fail-closed 判 A）。你留在 `a04d3e2` 的那段诊断输出就是本票的起点证据，一字未改地留着。
  ⇒ **本票 AC#3 的判据（那条 `test-windows` 步骤）的根因已闭，剩下的只是"哪一次真 run 证明它"**：
  票 72 的 AC#4 需要编排者 push 之后看新 run 的 `test-windows`；**步骤名 `PathResolver junction placeholder (real cases tickets 18/20)` 仍然误导**
  （那一步跑的是一条真红队用例，不是 placeholder），票 72 AC#6 想改成
  `PathResolver red-team cases (junction/8.3/UNC/\\?\ - tickets 18/20/72)`，但 `ci.yml` 是票 71 的在途文件 ⇒ 交回你排。




- [2026-09-21T11:5xZ] agent=agent-ticket70-d did=**认领 + 起点实测**（本轮承诺范围＝编排者指定的 4 项：AC#1 尾巴、AC#2 在新 HEAD 重测、AC#6 逐 job 逐步骤证据、D22 positive control 卡点）。起点 HEAD `17efc2c`，AC#1 尾巴先用 **CI 钉的 v0.7.0** 复核：`go install mvdan.cc/gofumpt@v0.7.0` ⇒ `gofumpt --version` = `v0.7.0 (go1.27.1)`，`gofumpt -l . tools/d22scan tools/mockllm` **只列 1 个文件** `internal/risk/provenance_syncdirs_windows_test.go`（票 75 代理新加、未跑 gofumpt）。⇒ 与简报一致，**旧"69 个文件"已过期**。
  next=一枚纯格式化 commit（只 `-w` 那个文件）→ 用 `docker run golang:1.27` 逐字复刻 `ci.yml:127-134` 重测 AC#2 → 把 run 35558750456 的逐步骤表贴进本 log → 自己复现 D22 positive control。
- [2026-09-21T11:58Z] agent=agent-ticket70-d did=**AC#1 的尾巴落地**（`6ce43c7`，**单独一枚纯格式化 commit**，只含 1 个 `.go` 文件）。判据逐字用 CI 钉的版本：`go install mvdan.cc/gofumpt@v0.7.0` ⇒ `gofumpt --version` = `v0.7.0 (go1.27.1)`；`gofumpt -l . tools/d22scan tools/mockllm` 改前 **1 行**（`internal\risk\provenance_syncdirs_windows_test.go`）→ 改后 **0 行、exit 0** ⇒ 本票 AC#1 正文的判据（"输出为空"）**第一次真的成立**。**行为改动＝0**：`git diff -w` 只剩那一条跨行 composite literal 的括号形状（字面量字符串、字段值、真值全部一字未动）。改后复跑：`go vet ./internal/risk/` rc=0；`go test ./internal/risk/ -count=1 -v` rc=0，`=== RUN` **149** / `--- PASS` **87** / `--- FAIL` **0** / `--- SKIP` **1**（SKIP 是既有的 `TestSyncRegistryProbeLive`："本机 HKCU 无 registry-grade 同步记录"，不是我造的，也没被算成 ok）。⚠ 这只闭 AC#1 的**本地半**；CI 侧 `gofmt (gofumpt)` 步骤要在**下一次 push 之后的 run** 里才是绿——本轮 HEAD `17efc2c` 的那次 run 该步骤就是 **failure**（红因正是这个文件），见下面 AC#6 条。
  next=AC#2（`docker run golang:1.27` 逐字复刻 `ci.yml:127-134` 的包清单，在新 HEAD+格式化树上重测）。
- [2026-09-21T12:1xZ] agent=agent-ticket70-d did=**AC#6 在新 HEAD `17efc2c` 上的逐 job 逐步骤取证** + **D22 positive control 卡点：我自己复现的结论是"已经不红了"**。

  **AC#6（run `35558750456`，headSha `17efc2c1d950975f712c8c159221be096de91c1a`，branch dev，event push，2026-09-21T03:48:12Z，`gh run view 35558750456 --json jobs` 原文结论）**：判据＝五个 job 全 pass ⇒ **本轮未达成，3 success / 2 failure**。逐步骤（不是逐 job）：

  | job | job 结论 | 逐步骤（真实 conclusion） |
  |---|---|---|
  | `lint` | **failure** | `D22 scanner positive control` **success** → `D22 seven-ban + emoji scan` **success** → `gofmt (gofumpt)` **failure** → `go vet (module)` `skipped`、`go vet (tools/d22scan module)` `skipped`、`staticcheck` `skipped`、`mockllm module vet` `skipped`。**硬规矩 (c)：那 4 个 skipped 不是"过了"，是 GitHub 在 `gofmt` 失败后不再执行**——所以 `17efc2c` 上"Linux 的 `go vet ./...` 到底绿不绿"至今**没有一次 CI 证据**（票 78 的两颗哑弹修没修好，CI 层面仍是未证，见下） |
  | `test-core` | **failure** | `Start compose test services` success → `Probe mock-llm` success → `Portable package tests` **failure** → `Environment fork assertion (WISP_ENV=test data dir)` `skipped`（同前：前一步失败，不代表它过）→ `Stop compose services` success |
  | `test-windows` | **success** | 7 个真实步骤全 success，含 `cgo build smoke`、`Portable windows tests (proc/secret/config)`、**`PathResolver junction placeholder (real cases tickets 18/20)`**（票 72 的因落地后这步第一次转绿；步骤名仍误导，归票 71/72 改） |
  | `slo-smoke` | **success** | `Build wisp.exe` + `SLO smoke gate` + `Upload SLO report` 全 success |
  | `slo-full` | **success** | `Build wisp.exe (deps cached on the runner)` + `SLO full gate (six states + settle + leak)` + `Upload SLO report` 全 success（AC#4 的 `MINGW64_ROOT` 现在是 **n=2**：`35549859581` + 本次） |

  **硬规矩 (a) 的排除**：同一时段 run `35554681609`（headSha `131f722`）`conclusion=cancelled`，`gh run view --json jobs` 回 `jobs: []`（**total_count=0**）⇒ **不是样本，直接排除**，没进上表任何一格。
  **`gofmt (gofumpt)` 那一步的红因（从 CI 日志原文取的，不是推断）**：`gofumpt would reformat:` / `internal/risk/provenance_syncdirs_windows_test.go` / `##[error]Process completed with exit code 1.` ⇒ **正是我 `6ce43c7` 修掉的那个文件**；顺带纠正一处本票记了很久的假前提：**CI 这条命令是 `go install mvdan.cc/gofumpt@latest`，没有钉版本**（日志：`go: downloading mvdan.cc/gofumpt v0.12.0`）。"CI 钉的 v0.7.0"是**票面早期写的**、与实际不符 ⇒ 我用**两个版本**都复核了：同一棵纯净树上 `v0.7.0 -l` 与 `v0.12.0 -l` **均 0 行、exit 0**，所以 `6ce43c7` 之后这一步在两个版本下都该绿。⚠ **`@latest` 本身是个门随机性来源**（哪天 gofumpt 改版就全体历史 commit 一起红），但 `ci.yml` 是票 71 的在途文件 ⇒ **只上报、不改**。
  另一条**账实纠正**：票 70-c 记过"gofmt 步骤第一次真 run 绿"（run `35549859581`，head `6c1b5e9`）——那是真的，但**只在 `6c1b5e9` 那个点上为真**；票 75 在 `ea01b9f` 新加的 `_windows_test.go` 又把它拖回红。⇒ **"AC#1 已闭"这件事没有永久属性，只有"某次 push 的某一步"属性**，这也是为什么这一步的判据必须每次 push 重看。

  **D22 positive control 卡点（简报第 4 步）——我自己复现：已经不红了，前任的锚点确实漂了**：
  - 复现环境与命令：`git archive HEAD` 解到仓外纯净树 ⇒ `sh scripts/d22scan.sh`（脚本自己派生 repo root）⇒ **exit 0**：`positive control - go test ./...` → `ok github.com/CarlosShao/wisp/tools/d22scan 2.631s`，随后真扫描自报 `examined 203 production Go files`、ban #8 `internal/ 303 Go files` + `cmd/ 24` + `design/ 16 text files`，末行 `clean - no D22 ban violations`；唯一 `NOT COVERED` 是 ban #6 `frontend/`（该树在此 HEAD 不存在，脚本文本明写这不是豁免）。
  - **谁修的 / 修的是什么**：`a8ae9ad`（**编排者本人**，09:40）把 `internal/tools/bridge_junction_windows_test.go:444` 注释里的 **U+26A0 `⚠`** 换成 ASCII 标记 `NOTE(d22 ban #8):`，1 文件 1 行、纯注释。⇒ 简报里"裸 typographic-quote 字形换成 Go 字符串转义（`\u201c` 形态）"这个修法**在本例上既不可能也不需要**：字形在**注释**里（注释内 Go 转义不生效），且它是 `⚠` 不是弯引号。**我没有为此动任何文件**。
  - 独立仪器复核（不信扫描器自己的嘴）：我用 ban #8 的**同一字符类**（`[\U0001F000-\U0001FAFF\u2600-\u27BF\u2B00-\u2BFF\uFE0F\U0001F1E6-\U0001F1FF]`）自己扫 `internal/` + `cmd/` 的 `.go`：**扫到 327 个文件（与扫描器自报的 303+24 逐字相符）、命中 0 处**。第 444 行的非 ASCII 码点逐字打出来只剩 CJK 与全角标点。
  - **`tools/d22scan/**` 与 `allowlist.txt` 一字未动**：`grep -vc '^#\|^$' tools/d22scan/allowlist.txt` = **5**（与票 70-c/票 67 记录的 4→5 之后的数字相同，既没变长也没变短）。CI 侧同样两份证据：上表 `lint` 的 D22 两步在 `17efc2c` 均为 **success**。⇒ 票面 AC#5 那条"⚠ 下一次 push 之后 lint 会同时红在 go vet 与 D22 seven-ban 两步"的**后半句作废**（D22 已被 `a8ae9ad` 清掉），前半句（`go vet`）仍未证。
  next=AC#2 在新 HEAD 的 Linux 逐包重测（`docker run golang:1.27` 逐字复刻 `ci.yml:127-134`）+ `go vet ./...`/`staticcheck` 在 Linux 上的真结果，一起写下一条。

- [2026-09-21T12:3xZ] agent=agent-ticket70-d did=**AC#2 在新 HEAD 上重新分诊**（票面"四族红"那张表按编排者的话**作废**，下面是替换它的新表）。判据＝**与 CI 那一行逐字相同的命令**跑在 **Linux** 上：`git archive HEAD` 纯净树 → `docker run --rm -i -e WISP_ENV=test golang:1.27`（容器内 `go version go1.27.1 linux/amd64`），包清单照抄 **`ci.yml:127-134`（编排者更正过的行号，不是 101-106）**的 16 个 `./internal/...` 树，只多带一个 `-v`（为了让仪器数得出条数，不改判据）。**exit code：`DOCKER_RC=1`（红，如实报）。**

  **总账（全量输出仪器数出来的，不是"默认只打第一条"的那种数）**：`=== RUN` 顶层 **497** = `--- PASS` **481** + `--- FAIL` **10** + `--- SKIP` **6**（**逐条对得上，0 条来路不明**；含子测试的 `=== RUN` 782、子测试 `--- PASS` 285、子测试 FAIL **0**）；包级 **16 ok / 3 FAIL / 4 no test files**；假绿四项逐个查过：`no tests to run` **0** 次；6 条 SKIP 逐条点名且**没一条被算成 ok**（`TestSubprocessCrashWriter`=helper 只被子进程用、`TestSyncRegistryProbeLive`=本机无 registry-grade 记录、`TestD34WriteMatrix`+`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`=Linux 上无卷标识、`TestRealDownload*`×2=要 `WISP_IT_REAL_MIRROR=1`）。

  | 族 | 票 70-c 在 `6c1b5e9` 的旧账 | **新 HEAD（`17efc2c`+我的 `6ce43c7`）实测** | 定性 |
  |---|---|---|---|
  | `internal/observe` | 2 FAIL | **ok 2.548s** ✅ | `d1525a3` 的时钟注入缝 + 按核数换算仍成立 |
  | `internal/secret` | 7 FAIL | **ok 0.095s** ✅ | `//go:build windows` 方案仍成立；portable 面（`refs_test.go`）仍在 ubuntu 跑 |
  | `internal/risk` | 17 FAIL | **8 FAIL / 2.271s** ❌ 仍红 | **同族但少了一半**（票 72 + 票 75 的路径形状修掉了 9 条）。剩下 8 条**同一个因**：`syncdirs_test.go:118 source "registry"/"config"/"env" must count as confirmed`、`:168 env-grade roots are confirmed evidence`、`provenance_test.go:469/541/608/670` 与 `syncdirs_test.go:356` 的 `precondition: the injected registry-grade root must confirm detection` —— Linux 上 `registryProbe` 返回 nil、`syncdirs_other.go` 带 `DEFERRED(P12-macos)`，于是 P12 那一层的前置条件在 ubuntu 上永远不成立。**全部在 `internal/risk/**`＝我的禁改区，与票 72 同区 ⇒ 未动、上报** |
  | `internal/tools` | 19 FAIL + **600s 超时 panic** | **ok 13.003s** ✅ | **票 75 的 C26 canonical-on-Linux 落地被这次复现证实**：不仅 19 条转绿，那条 `panic: test timed out after 10m0s`（旧 `wiring_test.go` 卡在 `approval.Gate.PendingApproval` 4m50s）**一并消失**——前任猜的"超时可能是路径故障的下游"这条**被证实**，没有第二个 Linux 审批洞 |
  | `internal/agent` | （旧账 ok） | **1 FAIL / 7.885s** ❌ **新增** | `TestSpillContainmentByDirectoryListing`（`spill_path_invariant_test.go:280/286`，票 76 的 `39420cf` 新加，`6c1b5e9` 时该文件**不存在**）：夹具拿 `..\..\..\..\CONTROL-escape` 当"控制逃逸"样本，Linux 上反斜杠是**合法文件名字符**⇒ 文件老老实实落在 artifacts 目录里（名字含反斜杠），夹具自己的前置"listing 必须看得见那个逃逸"不成立，于是它**主动拒绝跑成空断言**（打印 "AC#2 would be vacuous"）。**同一命令在 Windows 上 `--- PASS`（实测，`internal/agent 0.071s`）** |
  | `internal/memory` | （旧账 ok） | **1 FAIL / 19.364s** ❌ **新增** | `TestArtifactsContainmentByDirectoryListing`（`artifacts_path_invariant_test.go:464`，票 76 的 `d9224af` 新加）：purge 报 `removed [data/artifacts/nested\canary-separator.txt ...]`，want 却是 `data/artifacts/nested/canary-separator.txt` —— 同一条"反斜杠在 Linux 上是文件名而不是目录分隔符"的形状差。**Windows 上 `--- PASS`（实测，`internal/memory 0.102s`）** |

  **我的复现与 CI 的账实核对（这条最重要，因为它排除"是我机器/我容器特殊"）**：`gh run view 35558750456 --log` 里 CI 自己打的包级结论是 **`FAIL internal/agent 0.851s` / `FAIL internal/memory 10.822s` / `FAIL internal/risk 1.968s`，其余 16 包 `ok`（含 `ok internal/tools 9.285s`）**，`--- FAIL` 计数 **10**，与我在容器里数出的 **同一批 10 个测试名逐字相同** ⇒ 族定性不是单机偶发。
  **⇒ AC#2 判据（`test-core` 绿）未达成：3 包 / 10 条红。本轮一条断言都没改、一条 `t.Skip` 都没加、一个包都没移出 portable 清单。** 三族红的处置全是**归属判断**，见下面"交回编排者"。

  **顺手把 CI 从未产生过判据的三步补了本地证据（它们在 `17efc2c` 全是 `skipped`，不能记成"过了"）**：
  - `go vet ./...`（Linux，逐字同命令）⇒ **`VET_RC=0`**，`go vet -gcflags=-e ./...` 全量仪器**一行输出都没有** ⇒ 票 78 的两颗哑弹（`internal/ball/statevisual.go` 的 `undefined: mulA`、`cmd/wisp/slo.go` 的 `undefined: proc.Runtime`）**在 Linux 上确实修好了**（此前**从未**有过一次 CI 绿证）。
  - `go vet ./...` 在 `tools/d22scan`、`tools/mockllm` 两个嵌套模块 ⇒ **`D22SCAN_VET_RC=0` / `MOCKLLM_VET_RC=0`**。
  - `test-core` 的 `Environment fork assertion`（`bash tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`）⇒ **`FORK_RC=0`**，`top-level: PASS=1 FAIL=0 SKIP=0, === RUN=3, '[no tests to run]'=0`。
  - ⛔ **新炸点（本条唯一"下次 push 会多一步红"的预警）**：`staticcheck ./...`。CI 钉的是 **`staticcheck@2025.1.1`**，它在 **go1.27 上根本跑不起来**：6 行 `-: internal error in importing "internal/byteorder"/"internal/cpu"/"internal/goarch"/"math/bits"/"unicode/utf8" (cannot decode ..., **export data version 4 is greater than maximum supported version 2**)` ⇒ `SC_2025_RC=1`，**0 条真实 finding**（也就是说：**这一步的 CI 判据至今一次都没产生过**，它一直被前面的失败挡着）。把工具升到 `@latest` 后能跑（`staticcheck 2026.2.1 (0.8.1)`），但会**立刻挖出 35 条真 finding、`SC_LATEST_RC=1`**：`27×U1000(未用的函数/类型) + 3×SA1019(runtime.GOROOT 已废弃) + 3×S1011 + 1×SA4006 + 1×SA4000`，散在 **15 个包**：`cmd/wisp`×5、`internal/agent`×3 + `internal/agent/approval`×4、`internal/llm`(+adaptertest/anthropic/openaichat)×6、`internal/audio`×4、`internal/ball/tokens_table_test.go`×2、`internal/tools`×4、`internal/memory`×2、`internal/observe`×2、`internal/config/parse.go:223`×1、`internal/models/archive.go`×1、`internal/proc/crossvet_test.go`×1（**`internal/risk` 一条都没有**）。**我没碰其中任何一个文件**（`cmd/wisp/**`、`internal/config`（票 80 在途）、`internal/tools`、`internal/memory/*_test` 全在禁改或别人在途的清单里），也没动 `ci.yml`（票 71 的在途文件）⇒ 只上报：**"pin 太老跑不动"与"升上去就有 35 条账"是同一件事的两面，得编排者定顺序**（先清 35 条再升 pin，还是给 staticcheck 单开一张票）。
  next=把 AC#1 的框勾上（判据正文已满足，含 v0.7.0 与 CI 实际用的 @latest v0.12.0 双版本复核），然后交回编排者：AC#2 三族红的归属、staticcheck 的 35 条、以及"lint 下一次 push 会红在 staticcheck 而不是 gofmt/vet"。

- [2026-09-21T12:4xZ] agent=agent-ticket70-d did=**AC#2 的账在 15 分钟里被别人的落地改写了 ⇒ 追一条更正的数（HEAD 现在是我自己的 `c12c82f`，含票 81 的 `4683c34` 与票 83 的 `a95ee3a`）**。同一判据（`docker run golang:1.27` + `WISP_ENV=test` + `git archive HEAD` 纯净树）只重跑上面那三个红包：
  ```
  ok    github.com/CarlosShao/wisp/internal/agent            4.372s
  ok    github.com/CarlosShao/wisp/internal/agent/approval   0.037s
  ok    github.com/CarlosShao/wisp/internal/memory          39.089s
  --- FAIL: TestWriteGateNotSelectedByPayloadKey / ...PlainLocalWriteNotFlagged / ...EveryPathTargetJudged /
      ...AllPathsNonSyncStaysExempt / TestSyncFallbackNotDisarmableByWeakRoot / TestSyncEnvConfiguredRoots /
      TestSyncNormalNewFileWriteNotFlagged / TestSyncDotDotTailFailsClosed      (8 条)
  FAIL  github.com/CarlosShao/wisp/internal/risk  1.895s     RC=1
  ```
  ⇒ **`test-core` 在 Linux 上现在只剩 1 个包 / 8 条红，且 8 条全在 `internal/risk/**`（我的禁改区、票 72 同区、P12 的 registry/config/env 三档在 Linux 上永不 confirmed）**。上一条里我数的 `internal/agent` 与 `internal/memory` 各 1 条**已被票 81 的 `4683c34`（"derive the purge expectation from where the OS resolved each canary, not from a name list"）修掉**——修法正是 AC#2 允许的那条"平台专属期望值"，**不是** skip 也不是删断言，值得记一笔：**它证明了那两条红不需要我来动手**。
  **AC#2 结论不变：判据（`test-core` 绿）未达成 ⇒ 框不勾**，但"未达成"的原因从"四族"收敛成"一族半"（risk 8 条 + 下面两件与 test-core 无关的 lint 侧炸点）。
  next=交回编排者三件事：(1) `internal/risk` 那 8 条在 ubuntu 上该不该算 portable（**移出清单是编排者的判据**）；(2) `staticcheck` 的 pin 跑不动 vs 升上去 35 条 finding 的先后顺序；(3) push 之后按 AC#6 的口径重看 `lint`/`test-core` 两步——**`gofmt (gofumpt)` 这次会绿**（AC#1 的 CI 侧最后一块证据），`test-core` 预期仍红 8 条。

- [2026-09-21T12:5xZ] agent=agent-ticket70-d did=⛔ **简报第 4 步的 D22 positive control 卡点：以另一种形状复现了——而且这次的制造者是我自己动不了的人**。我上面那条"D22 已不红"在**它当时的锚点上仍然为真**（run `35558750456`/head `17efc2c` 两步 success、`a8ae9ad` 清掉了 :444 的 `⚠`），但我交完 AC#2 之后按"最新 HEAD 再核一遍"的规矩在 `git archive HEAD`（HEAD=`c12c82f`，已含票 83 的 `a95ee3a`）的纯净树上重跑 `sh scripts/d22scan.sh`：**`D22_RC=1`，又死了，死在同一条 positive control 上**：
  ```
  --- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.82s)
      scan_test.go:269: repo HEAD violates: internal/config/unwired.go:13:
          [emoji] ban #8 glyph in scope internal/ is banned (D23):
          covers comments and _test.go, not only string literals
  --- FAIL: TestRealRepoLedgerIsHonest (0.59s)
      scan_test.go:667: HEAD must be green, rc=1 out=d22scan: examined 99 production Go files ...
  FAIL github.com/CarlosShao/wisp/tools/d22scan 12.810s
  ```
  **命中点逐字**：`internal/config/unwired.go:13` 注释里的 **U+1F512 `🔒`**（`//  1. it sits in a locked (🔒) section ...`），由**票 83 的 `a95ee3a`（"six locked-section keys that lie now fail the load loudly"）**带进来，落库时间就在我这一轮的中段。我的独立仪器（ban #8 同一字符类，自己扫 `internal/`+`cmd/` 的 `.go`）给出全仓真数：**扫 329 文件、命中 1 处**，就是这一处 ⇒ 不是"2 处错误"那类采样假象。⚠ 顺带一条对仪器的更正：`TestRealRepoLedgerIsHonest` 那次自报的是 `examined 99`（不是上次的 203）——它扫的是 `tools/d22scan` 自己的**测试夹具树**而非全仓，别把 99 读成"覆盖面缩了"；全仓那一步的 203/184/19/303/24 在 `17efc2c` 上仍是 CI 日志里的原值。
  **为什么我这次也不修**：(a) 简报的禁改区写着"票 80 的代理正在 `internal/config` 里查 `bOverrides`，别去那里"，`a95ee3a` 又是几分钟前的活人产物；(b) `tools/d22scan/**` 与 `allowlist.txt` 是硬禁（allowlist 此刻仍 **5 条非注释行**，我没动）。**修法与 `a8ae9ad` 同形、且比它更省事**：把注释里的 `🔒` 换成 ASCII（例如 `locked (L2-locked) section`），1 文件 1 行、纯注释、零逻辑；**简报第 4 步说的"换成 Go 字符串转义 `\u201c`"在这个案子上做不到**——字形在注释里，注释内 Go 转义不生效（两次卡点都是这种形状：`⚠`/`🔒` 都在注释，都不是弯引号）。**若票 83 的代理已经收工，这条我随时可以补一刀，但需要编排者点头放行 `internal/config`。**
  ⇒ **对下一次 push 的预测（先写下来，别事后归因）**：`lint` 会红在 `D22 scanner positive control`（第 1 步，`tools/d22scan` 的 `go test ./...` 因真仓违规而 FAIL）+ `D22 seven-ban + emoji scan`，**而 `gofmt (gofumpt)` 会绿**（我在 HEAD `c12c82f` 的纯净树上双版本复核：`v0.7.0 -l` **0 行 rc=0**、`@latest=v0.12.0 -l` **0 行 rc=0**）；`go vet` 三步入 CI 之后**大概率绿**（Linux 容器实测 rc=0），再往下就会第一次真正跑到 `staticcheck` ⇒ **红**（pin 跑不动）。也就是说：**AC#1 绿了、AC#5 的门反而被别人的新文件弄红了，而 AC#6 的判据（五 job 全 pass）一次比一次看得清——它红得越来越"诚实"了。**
  next=交回编排者拍板：`internal/config/unwired.go:13` 的 1 行注释我改不改（放行与否）；`internal/risk` 8 条在 ubuntu 的归属；`staticcheck` 的 pin-vs-35-findings 顺序。票面 6 枚 commit（`d1455d1`/`6ce43c7`/`7f050b7`/`7b85b09`/`c12c82f`/`e4082a0` + 本条），**未 push**。

- [2026-11-12 11:4x +08] agent=70-ac2-triage 只读分诊 test-core（ubuntu，run `43146479546` step7）现取 **FAIL=15**（tools 9 + plugin 6；`=== RUN=3037 PASS=2037 SKIP=4`，日志 ：12321-:12322）：**A 类 11 条**（`joinForListing` 硬编码 `\` 一处致 search 全家 7 条＋`launcher_other.go` 缺 Job 失败措辞/AUMID 臂 2 条＋plugin 纯函数 fixture/错误形状 2 条）／**B 类 4 条**（Test52 pin 家族：D46 盘符判据 Windows 形状，建议 ledger 具名出 core scope，windows 腿覆盖不丢）。另记两条给编排者：票面曾写的 "4" 是归档 run-37166458550 的数非本次；**SKIP=4 全不在 ledger，15 条修完后 step 仍会因未记账 SKIP 红**。表见 docs/evidence/s1/70-ac2-testcore-triage.md。AC#2 是否算满足留编排者裁，本腿未翻框。

- [2026-11-13 14:1x +08] agent=70-ac2 写腿分诊+修复（接上条，run `43362781437` 两 job 现取）：**"6 枚"横跨两个 job**——test-core FAIL=4（Test52 pin 四枚，`resolves to a drive-less path`）＋ test-windows FAIL=2（`Test52AbsoluteTempPath` 的 `contains "~"` ＋ config `TestReadFailureDoesNotAdoptStat` 的 `got report &{Hot:[ball] ...}`）；config 那枚不在 test-core（`_windows` 文件 linux 不编译）。**基线更正**：`255224a56` 的"最后一次绿"是 slo-fresh；255224a56 上 11-12 15:00Z 的 ci run `43243081268` test-core 已 FAIL=4 恰是四枚 pin。**裁定**：四枚 pin＝**A 类**（前腿判 B 出 scope 的结论被本腿推翻——根因是 `command.go` 把 Windows 的 `filepath.VolumeName==""` 语义无条件化，POSIX 绝对路径完全确定，不属"可移动"；修=isAbsSpelling/resolveExe 分平台收窄，windows 判据一字未变）；`Test52AbsoluteTempPath`＝**A 类**（`~` marker 把 runner 8.3 别名 `RUNNER~1` 当 home-相对拼写抢跑，理由错；`~` 收窄为行首语义＋测试改规则探测器，ubuntu 从 SKIP 变真跑）；config 一枚＝**C 类**（hosted windows-2025 会话不兑现 share-0 冲突，本机 ×21 绿含 runner 同款 go1.27.4；ci-r1 §5#3 独立记载过同一枚）——新件 `manager_readfail_42r1_envprobe_windows_test.go` 落 PROBE-UNDO 探针闸＋具名 skip，**ledger 五字段行已备好在证据件 §4，scripts/portable-tests.sh 在本腿写面外，需编排者登记，否则 test-windows 该步从 vacuous-pass 红变 honest-skip 红**。**成对判据**：linux 容器（golang:1.27 真跑）修前恰红四枚 pin → 修后 `ok`×2 且 SKIP census=0；本机模拟 `~`-TEMP 环境红句逐字复现 → 修后 PASS。**门禁**：build/vet/-count=2/gofumpt（我写面 0 行；models 1 行是 64a97a610 他腿今日件）/d22scan clean 全过。证据：docs/evidence/s1/70-testcore-triage-6.md。commit 显式 pathspec 未 push。AC#2 框不翻（编排者裁）。

- [2026-11-13 15:5x +08] agent=70-core-green 写腿（test-core 在 FAIL=0 之后只剩 2 枚 SKIP → 走到真绿）。**逐枚裁定**：①`TestAppendSessionTurnCrashChild`——**实测在 `internal/memory/dao_session_append_326_test.go:44`（派单写 internal/session 是笔误；全仓 grep 仅此一处）**，票 326 AC#3 的 reexec **子进程角色**（首行即 Skip，仅被父测试 `exec.Command(os.Args[0])` 重启 + TerminateProcess 硬杀时干活）＝**顶层永不可跑**（它判据的一半就是被真杀；把它变"能跑"等于复刻父测试的机制而非执行一个用例）⇒ **ledger 具名登记** `any|reexec`（scripts/portable-tests.sh:819，与既有 `TestSubprocessCrashWriter` 同形先例，票 93 机制）；真断言持有者（父测试）ubuntu 容器实测照跑真过，断言零丢失。②`Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm`＝**能在 ubuntu 真跑 ⇒ 删过期守卫让它真跑**（判据链：形状门 `windowsAUMIDPattern` 在平台中立 `launcher_aumid.go`（原注释 "EVERY GOOS classifies ... the same way"），`launcher_other.go:62` 在 Linux **同样**分类 Kind=aumid；skip 理由 "the AUMID shape gate is the Windows launcher" 已过期；断言一字未动；Docker go1.27.1 实测 `--- PASS`、SKIP=0）。③mockllm SA9003（messages.go:145 空分支）＝有意为之的注释性 no-op，删空分支、说明并入真实分支注释，staticcheck 复跑 `SC_RC=0`、模块测试 ok。**ubuntu 读数**（Docker `golang:1.27`＝CI 同款镜像）：archive 树全量 core scope `RUN=3038/PASS=2055/FAIL=2/SKIP=0`——FAIL=2 是**导出法假象**（无 .git ⇒ panel 的 git 自洽测试前提不成立；core.autocrlf 把 composer-states.html 导成 CRLF ⇒ render 夹具切 0 块），改容器内 git clone 复跑 `RUN=2993/PASS=2023/FAIL=0/SKIP=0`（panel 转绿坐实；**顶层断言 FAIL=0、SKIP=0 双跑实测**）；llm 3 包在 clone 跑进程级死（无 --- FAIL 行无 panic，同树同 overlay 第一跑全绿 ⇒ 宿主容器资源 flake，非本腿改动）。**门禁**：build/vet/test(-count=1)/d22scan(clean)/bash -n×2 全过；`scripts/runtests.sh` 不存在（rc=127），规矩持有者是 `tools/d22scan/runtests.sh`（未在 scripts/ 下新造）。证据：`docs/evidence/s1/70-testcore-green.md` + `D:\tmp\wisp\20261113\70-core-green\repro-log.md`。commit 显式 pathspec 未 push。
- [2026-11-14 12:0x +08] agent=70-triage 分诊腿（只读，⛔ 零产品/测试/脚本改动，未 push；证据 docs/evidence/s1/70-ci-red-triage.md ＋ 探针件 .scratch/wisp/probes/70/）。run `43549905387`（`096f247ed`）逐 job 定因，四步结论：**step 4＝连坐**（`RUN=0` 是真值、不是日志假象：`scripts/portable-tests.sh:851` 的"名册宇宙"把 `./cmd/wisp/` 拉进同一发 `go test -list`，在 DLL 暂存之前（`ci.yml:954` vs `:1089`）必然 `exit status 0xc0000135`＝STATUS_DLL_NOT_FOUND，`:856 exit 1` 让整个 winsec tier 在跑第一枚测试之前退出；`winsec-tests.sh:235` 的 GUARD 2 反被 `-list` 的名册尾巴 `ok …internal/winsec 0.069s` 满足＝真判据假输入；点着它的是 2026-11-13 ledger 新增的 `./cmd/wisp/` 五行（`02939b8a2`/`7386f221c`），老机制 × 新行）；**step 10＝真红**（派单指的 `leg_sink_nail_131_windows_test.go:358` 是 `readLegSink131` 的夹具调用（`t.Helper()` 把日志记到调用者行）、`:396` 是 `t.Logf`，两处都不是断言且用例 PASS；真正的红是 `--- FAIL: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands (11.16s)` → `subagent_selfapproval_197_test.go:232` 的 `t.Fatalf("exit %d…")`＝根任务 exit 1，尸检是 `QuotaLedger.roll`（`internal/observe/quota_44.go:256`）nil 指针 panic，而 `QuotaLedger`（`:221-228`）无锁、`QuotaGate.Decide`（`internal/agent/quota_44r4.go:155`）在每次派发前 `Spend()→roll()`；END ANCHOR 成立（`619=PASS 599+FAIL 1+SKIP 19`，末名 `TestZZ315A1…` 打印 PASS）⇒ 跑完了、不是连坐；19 枚 SKIP 不是红因（同 20 发里 17 发 step 10 success 且 SKIP 同样 19））；**step 13＝连坐**（与 step 4 同一道 `:851` 闸：`--- FAIL` 0 枚、四数一行未打、END ANCHOR **没响且不可能响**（`portable-tests.sh:1029` 只对 cli tier 生效），17 行 `ok …0.0xx s` 是名册尾巴（本机实测 `go test -list` 末行就长这样）；它排在被点名的 `Fetch voice fixtures` 之后仍红 ⇒ 补夹具治不了）；**step 7＝时序性真缺陷（窗口型）**（`internal/tools/task_cancel_221_legs_test.go:452` 没有自己的上界，30.01s 是 `internal/tools/bridge.go:28 DefaultToolTimeout=30s` 到点后两枚挂着的 `task.spawn` 收口才让 park 循环继续；根因是 `internal/tools/subagent_197.go:373-374` 先 `PublishSubagent` 再 `AttachCancel`，而子侧 `AdmitTask`（`:337-345`）在第一次模型调用前**只发布行**，测试放行点 `spawnHeld`（`:132-158`）恰以第一次模型调用为准 ⇒ 行可见/句柄未挂，产品自己的答复自证矛盾（`state=Thinking` 却说"没有在跑"）；分母：最近 20 发 step 7 18 success/2 failure＝10%，两发命中是同一句拒绝文本的两枚不同用例（`:452` 与 `:322`），同枚 1/20；本地 `-count=3` 3/3 PASS（0.01/0.00/0.00s，go1.27.1）⇒ 窗口窄但真）。**追加第二发** run `43551571039`（`1ce39a651`＝当前 HEAD）同结构红：step 4/13 仍是 `0xc0000135`（19/19 确定性红），test-core 换了一枚本族用例红（`:322`）。**修法候选（都没做）**：①名册复验改逐包隔离 / 名册阶段也走 DLL 暂存（owner 编排者，`scripts/**` 的账）；②`QuotaLedger` 加锁或 `Spend()` 只读快照（owner `internal/observe`，票 44 家族）；③发布与挂句柄合成一次发布、在 `AdmitTask` 内完成（owner `internal/tools`，票 221/197 家族）。三处的"文件:行"与逐字读数都在证据件里。
- [2026-11-14 12:2x +08] agent=44-quota-race 写腿（接 70-triage 的 ②：`internal/observe`，票 44 家族）did=给 `QuotaLedger` 加 `sync.Mutex`（`roll` 拆成持锁的 `rollLocked` ＋ 上锁薄壳 `roll`，`Post`/`Spend` 持锁；stamp 与三张 map 进同一临界区；⛔ 行为语义零改动、⛔ 不加第二把尺）。复现件 `internal/observe/quota_44_race_test.go`（16 goroutine×200 轮纯 `Spend`＋每读前进一天的时钟；8 读×8 写 `Post`/`Spend` 固定日）：**修前 HEAD `-race` = 39 条 `WARNING: DATA RACE` ＋ 1 条 `fatal error: concurrent map iteration and map write`，rc=1**，命中行覆盖 `quota_44.go:251/252/253/256/257/258`＋`:284/286/287/295`（**含 CI 点的 `:256`**）⇒ 复现成立；修后同一件 `-race -count=5` **10/10 PASS**、5 个独立进程 **5/5 rc=0 races=0**、`go test -race -count=1 ./internal/observe/` = `ok … 33.079s`。定向 `go test -count=3 -run Test197 ./internal/agent/` 四数 **RUN=0 PASS=0 FAIL=0 SKIP=0**（**this package has no Test197 family**，逐字 `[no tests to run]`；CI 真红那枚在 `cmd/wisp`，要 DLL 暂存＋夹具，本腿未跑、不冒充；同族 `./internal/tools/` **33/33 绿**）。门禁：build rc=0／vet rc=0／`gofumpt -l internal/observe` 空／`sh scripts/d22scan.sh` **rc=0 clean**（正控 PASS=45）。**遗留具名（不越界）**：`SeedFrom`（`internal/observe/costmodel_44b.go:186-201`）仍在锁外直写三张 map，不在本腿写面（它只在 boot seed 一次、不在本红路径上；修法＝该函数外层加锁＋首行改 `rollLocked`）。证据：`docs/evidence/s1/44-quota-ledger-race.md`。commit 显式 pathspec、**未 push**。
- [2026-11-14 12:3x +08] agent=44-quota-seedfrom-lock 写腿（**收**上一腿 `6f23ce632` 具名留下的那一枚尾巴，写面仅 `internal/observe/costmodel_44b.go`＋扩既有仪器＋证据件＋本行）did=`SeedFrom` 原先锁外直写三张 map（`costmodel_44b.go:193/195/198`），现**整段在 `l.mu` 里**且首行改调 `rollLocked`（持锁版；若先 `Lock` 再调 `roll()` 会重入同一把 `sync.Mutex` ⇒ 当场死锁，故持锁段内只许 `rollLocked`）。**正控（修前 HEAD 现跑）**：新仪器 `TestQuotaLedgerSeedFromOverlapsSpend`（4 读 `Spend`×200 轮 ＋ 1 写 `SeedFrom`×40 轮，时钟固定同一天，两行 seed 分别打 day/month/plan 三桶）**4/4 发 rc=1、6/6/7/6 条 `WARNING: DATA RACE`、fatal=0**，逐字头一条＝`Write … costmodel_44b.go:193`（`runtime.mapassign_faststr`）↔ `Previous read … quota_44.go:333`（`Spend.func1` 的 map 拷贝），另有 `:195`↔`:334` 一对 ⇒ **重叠造得出来，不是"造不出"**。**修后**：同一枚 5 个独立进程 **5/5 rc=0 races=0**、同进程 `-race -count=5` **5/5 PASS**、`go test -race -count=1 ./internal/observe/` = `ok … 25.134s`。**行为语义零改动**：入桶三等式、负行拒绝的同一句错误文本与同一时机（roll 之后/入桶之前）、返回值、`Spend` 可见字段全未动，既有两枚语义钉子（`TestQuotaLedgerSeedFromBucketsDayMonthPlan`/`TestQuotaLedgerSeedRolloverIsModelCorrect`）照旧绿；仪器尾部另加**安全半边**（丢任一微单位即红），不是"检测器安静就算过"。门禁：build rc=0／vet rc=0／`gofumpt -l internal/observe` 空／`sh scripts/d22scan.sh` **rc=0 clean**（正控 PASS=45 FAIL=0，真扫 481 枚生产 Go 文件 0 发现；**第一次红在本人刚写的那一行 phantom-citation 上**——注释引了本件而本件还没落盘，建件后重跑即 clean，非他人文件、非放宽）。**具名后果（不越界）**：`SeedFrom` 一改，`quota_44.go:297` 的薄壳 `roll()` 已无生产调用点，`staticcheck` 现读 `internal\observe\quota_44.go:297:23: func (*QuotaLedger).roll is unused (U1000)`（CI staticcheck **pass 2 `GOOS=windows -checks=all`** 会读到的新增一条；pass 1 `-U1000` 读不到）；修法＝删 `quota_44.go:289-301` 的 `roll`＋其上过期注释共 13 行，**不在本腿写面**（`A1773` 那条"不许删 U1000 符号"针对 windows 侧活着/linux 侧假死的符号，这一枚两个 GOOS 都死）。证据：`docs/evidence/s1/44-quota-seedfrom-lock.md`。commit 显式 pathspec、**未 push**。
- [2026-11-14 12:4x +08] agent=44-quota-u1000 写腿（**摘** `44-quota-seedfrom-lock` §5.1 具名留下的那一枚 U1000，因为它挡 CI `lint` 的 staticcheck pass 2）did=删掉 `internal/observe/quota_44.go` 里**零调用者**的 `roll()` 薄壳及其过期 doc（`756ebbcfa` 把唯一调用者 `SeedFrom` 改成 `rollLocked` 之后它就死了），并把 `rollLocked` doc 里两处对 `roll` 的指涉改成现状（`Post`/`Spend`/`SeedFrom` 各自持 `l.mu`；`SeedFrom` 自己持锁、先调 `rollLocked` 再在同一临界区入桶 ⇒ 没有取锁变体可调），⛔ 没为"让它被用到"加任何调用（加调用是把 U1000 藏起来而不是修它）、⛔ `rollLocked`/`Post`/`Spend`/`SeedFrom` 函数体一字未动。**删前数调用点**：全仓 `.go` 的 `\.roll()` 命中 6 处＝1 定义＋5 注释，**生产/测试调用点 0** ⇒ 两个 GOOS 都死（唯一真调用点在平台中立的 `costmodel_44b.go`），与 `A1773` 那条"windows 活/linux 假死不许删"不同类。**读数**：`staticcheck -checks=all ./internal/observe/`（`GOOS=windows`，2026.2.2）修前 `quota_44.go:297:23: func (*QuotaLedger).roll is unused (U1000)` **rc=1** → 修后 **rc=0 逐字无输出**；三包（observe+agent+tools）同口径 **rc=0 无输出**；`go test -race -count=1 ./internal/observe/` **rc=0 `ok 29.560s`**；build 三路 rc=0／vet rc=0／`gofumpt -l internal/observe` 空（0 行）／`sh scripts/d22scan.sh` **rc=0 clean**（481 枚生产 Go 文件）。`staticcheck ./...` 仍 rc=1 但三处**全在 HEAD 就存在且文件与 HEAD 逐字节相同**：`build/mut_rvw.go` 11 处 undefined（**未跟踪**的变异实验遗留件，CI 的 build 那发不含 `build/`）、`frontend/embed.go:4` SA9009、`internal/ball` 两条 U1000 恰是 `429387f46` 票 122 具名的 Tier-4 残渣。**具名遗留（写面外、纯注释）**：`costmodel_44b.go:191` 仍写 "calling `roll()` would … deadlock"（已指空，建议改 "no lock-taking variant to call"）、`quota_44_race_test.go:19/36/71` 三处以 `roll()` 指代换日动作（非调用）。证据：`docs/evidence/s1/44-quota-seedfrom-lock.md` **§6**。commit 显式 pathspec、**未 push**。
- [2026-11-14 13:0x +08] agent=70-roster-gate 写腿（接分诊的 ①：名册闸连坐 ＋ GUARD 2 假输入；写面 `scripts/portable-tests.sh`、`scripts/winsec-tests.sh`，未 push）did=**名册阶段改逐包隔离**：合并 `go test -list` 先走快路，失败才逐包 `-list`；列不出来的包具名 ＋ rc ＋ 加载器签名、计数打印；只在 ledger 里的包（本 tier 不跑）其 ledger 行仍留在 `-skip` 模式里但**具名 DEFERRED**（由跑它的 tier 复验，`./cmd/wisp/` → `mode=cli`），scope 内的包则由 GUARD B 的每包结果行 ＋ 末尾 ticket-70 判决兜红；**没有**在名册阶段暂存 DLL（理由写进证据件 §2：DLL 载具只对一包有效，且会把载具知识复制成两份）。**GUARD 2 收紧**：包级结果行之外还必须看到执行读数（`--- PASS|FAIL|SKIP:` 行或 `portable-tests.sh: four numbers` 行），否则红并说明"抓到的是 `-list` 名册尾巴"。读数：本机 HEAD 字节 `bash scripts/winsec-tests.sh` = `:1247 0xc0000135` → `:1252 RUN=0 PASS=0 FAIL=0 SKIP=0` → rc=1（与 CI step 4 逐字同形，确定性不是 flake）；改后同一命令 = `1 of 10 roster package(s) could NOT be listed`、`./cmd/wisp/ [rc=1 0xc0000135]`、5 条 DEFERRED 逐名、`RUN=129 PASS=78 FAIL=0 SKIP=0`、真结果行 `ok … 57.463s`、rc=0。step 13 镜像 `--scope=windows` 改后 `RUN=1274 PASS=832 FAIL=0 SKIP=0`、13 包逐个 ok、rc=0。正控：GUARD 2 面对"只有名册尾巴"（桩 exit 0）必红（逐字红句在证据件 §4.2），回放真实 129 条执行行则绿（rc=0）。END ANCHOR 票 314 两段按标记抽取 sha256 前后相同（LF `3dac299b…`／`bceb9103…`；CRLF worktree 字节 `89aff5ba…`／`c60726a0…`），`git diff -U0` 无一 hunk 落入；真跑 `--scope=cli` 打印 `END ANCHOR … roster=619 … TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（随后为省时中止，未主张该发绿）。门禁：`bash -n` rc=0×2；d22scan 第一发 rc=1 是别人在飞的 `costmodel_44b.go` phantom-citation，等其落盘后重跑 **rc=0 clean**；slo-nail-census 52/0；path budget GREEN；既有载体 `portable-tests-selftest.sh` 在 LF 字节上 **32/0 GREEN**（本机 CRLF 检出会让 3 条变异用例假红，已给 LF/CRLF 对照读数）。证据 `docs/evidence/s1/70-roster-gate-fix.md` ＋ `.scratch/wisp/probes/70/fix/**`。commit 显式 pathspec、未 push。
- [2026-11-14 13:2x +08] agent=today-four-fixes-audit（非实现者，审计腿）did=审计今日四笔修复对 CI 红的闭合：A 类突变逐笔现跑（44 锁 25 races+fatal、SeedFrom 锁 8 races、197 形状尺红、314 造不出红），B 类零放宽，D 类仍有 2 处固定名 profile 与形状尺字符串逃逸；四笔结论=成立/成立/成立/部分；AC#6 五 job 全 pass 仍不可翻（需 push 后真 run）。证据 docs/evidence/s1/today-four-fixes-audit.md

## AC#6 全绿证据（编排者现取，2026-11-14 13:4x）：**CI 6/6 全绿，归因单变量**

**run `43559883923`（head `b52365b7120c8128f0b2a9ca2c3f297ba361bd62`，起于 13:06:18）逐 job：**

| job | 结论 | 失败步骤数 |
|---|---|---|
| `lint` | **success** | 0 |
| `slo-smoke` | **success** | 0 |
| `slo-full` | **success** | 0 |
| `test-windows` | **success** | 0 |
| `test-core` | **success** | 0 |
| `lint-frontend` | **success** | 0 |

### 归因（三发对照，单变量）
| commit | 配额锁 `6f23ce632` | 197 窗口 `d8e2f4dc8` | 名册闸 `dada2ed9e` | CI |
|---|---|---|---|---|
| `0c9ff7c0c` | ✓ | ✓ | ✗ | **failure**，仅 step 4＋step 13 |
| `b70ea907f` | ✓ | ✓ | ✗ | **failure**，仅 step 4＋step 13 |
| **`b52365b71`** | ✓ | ✓ | **✓** | **success（6/6）** |

⇒ 两发失败的**唯一**差异是 `dada2ed9e`（名册逐包隔离＋GUARD 2 假输入），而失败步骤恰好是它治的 step 4／13 ⇒ **因果闭环**。
另：`test-core` 在 `0c9ff7c0c` 已转绿 ⇒ 之前那枚 `Test221SubagentCannotStopSiblingOrItself`（10% 时序窗口）在 197 修复后不再出现（⚠ 一发样本不算证明，真证据是 6400 枚在飞 0 命中）。

**AC#2 一并说明**：票面写的「在 Linux 上复现那 4 条」前提已被分诊推翻 —— `test-core` 实际只有 **1 枚**红（`Test221…`，30.01s＝撞 `DefaultToolTimeout=30s` 上界），且它的根因是 `subagent_197.go` 的发布/挂句柄窗口，已修并在 CI 上验绿 ⇒ **AC#2 以「CI 上不再出现」为准判过**，不再另建 Linux 复现。
---

## 结案（编排者，2026-11-14 14:0x）：**6/6 全勾，两发连续绿**

### 全绿证据（`gh run view <id> --json jobs` 逐 job，两发）
| run | head | 结论 | 逐 job |
|---|---|---|---|
| `43559883923` | `b52365b71` | **success** | lint ✅ · slo-smoke ✅ · slo-full ✅ · test-windows ✅ · test-core ✅ · lint-frontend ✅ |
| `43561054718` | `df6a8e035` | **success** | slo-full ✅ · lint ✅ · test-windows ✅ · test-core ✅ · slo-smoke ✅ · lint-frontend ✅ |

**为什么要求两发**：这套测试有 10% 的 flake 前科（`Test221SubagentCannotStopSiblingOrItself` 20 发里 2 发命中），**一发绿不足以结案**。

### 归因（三发对照，单变量）
| commit | 配额锁 `6f23ce632` | 197 窗口 `d8e2f4dc8` | 名册闸 `dada2ed9e` | CI |
|---|---|---|---|---|
| `0c9ff7c0c` | ✓ | ✓ | ✗ | failure，**仅** step 4＋13 |
| `b70ea907f` | ✓ | ✓ | ✗ | failure，**仅** step 4＋13 |
| `b52365b71` | ✓ | ✓ | **✓** | **success 6/6** |

### AC#2 的处置（**前提被推翻，按「已解决」勾**）
票面写「在 Linux 上复现那 4 条」——分诊实测 `test-core` **只有 1 枚**红：`Test221SubagentCannotStopSiblingOrItself (30.01s)`，`30.01s` 是撞 `internal/tools/bridge.go:28 DefaultToolTimeout=30s` 上界，**不是墙钟猜想**。根因在 `internal/tools/subagent_197.go` 的「行先发布、取消句柄后挂」窗口，已由 `d8e2f4dc8` 修好，并在两发 CI 的 `test-core` 上验绿 ⇒ **本格以「CI 上不再出现」判过**，不再另建 Linux 复现。撤销口令：「70 AC2 复议」。

### 本轮落地的四笔（全部已推 `origin/dev`）
1. `dada2ed9e` 名册阶段逐包隔离 ＋ `winsec-tests.sh:235` GUARD 2 假输入（**END ANCHOR 两段 sha256 逐字节同串**）
2. `6f23ce632`＋`756ebbcfa`＋`f8e94bbd8` `QuotaLedger` 无锁竞争（修前 39 条 DATA RACE ＋ 1 fatAL，修后 `-race` 5/5 绿）
3. `d8e2f4dc8`＋`0c9ff7c0c` 子代理发布/挂句柄窗口（6400 枚在飞 0 命中）
4. `b52365b71` 起 CI 转绿

### 更正记（2026-11-14 15:5x，编排者补，只追加不改上文）

本票结案判据是「两发连续绿」，那两发**确实**达到了，判据没被降低。但今天后续 10 发里 **3 发红**，两处**互不相干**的 flake 源现量如下 ⇒ **「CI 已经稳稳绿了」这句话不成立，本票不作此声明**：

1. `test-core`（Linux）`TestPOSIX120SealWindowLeafSwapIsClosed` —— run `43565131755`(head `05fb21317`) 红，同套在 `069f97cac`／`779d6c6fa`／`86e5655ba` 三发绿 ⇒ **已立票 441，已修**（拆两载体、弃权具名带数字）。
2. `test-windows` step 13 `Test52TimeoutKillsTheProcessTree` —— run `43567331249`(head `e1cc304e0`) 与 `43568581591`(head `1f0d04714`) 红，`a46352fd3` 绿 ⇒ **3 发红 2**，**已立票 443**。红句 `command_52_test.go:748: the fixture never reported a grandchild pid` ＝仪器没观测到就判红（与 441 同机理）。

⇒ 本票修的是**名单闸那道**（roster gate），那部分证据依然成立；**整条 CI 的可信度**由 441／443 负责，不等于本票已结。
