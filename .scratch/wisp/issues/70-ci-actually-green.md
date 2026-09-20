# 70 — 让 CI 真的成为护栏：5 个 job 全红，逐因分诊（本仓的"门禁"从未生效过）

**Status:** blocked-on-tree（**必须等票 66 与票 68 收尾**：AC#1 是全仓格式化，会重写它们正在改的包）
**Claimed by:** —
**Last update:** 2026-09-20
**Blocked by:** 66（`internal/observe`/`internal/proc`/`cmd/wisp`）、68（`internal/ball`/`cmd/balldebug`）
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
- [ ] **AC#5 ban #1 覆盖洞（A26 剩余半）**：`go probeReader()` 这类**具名调用**今天扫不出来
  （生产范围实测 4 处：`internal/observe/goroutine.go:281` 受管 + `cmd/balldebug/main.go:231/259/440`）。
  要么把 ban 扩到具名调用并由我**书面**给 observe 豁免，要么把"clean"的语义在输出与文档里**写死成当前范围**。
  ⚠ 三处 `balldebug` 命中若被新覆盖面抓到，**照实报**，由我裁定归属，别自行豁免。
- [ ] **AC#6 全绿证据**：`gh run view <新 run> --json jobs` 的**逐 job 结论**贴进本票 Progress log。
  判据是**五个 job 全 pass**，不是我某一步本地过了。**`slo-smoke` 若仍红，先确认票 66 是否已合**，
  不要为了让它绿而调采样阈值（D32 的 0.5% 一字不动）。

## 编排者已裁定

1. **不放宽任何一步 CI**：不加 `continue-on-error`、不删步骤、不把阈值往下调。
2. **格式化与逻辑改动永不同 commit**（AC#1 单独成 commit）。
3. 本票 AC#1 落地期间**只有它一个写码代理**，避免全仓重写与别人 WIP 相撞。

## Progress log（append-only）
