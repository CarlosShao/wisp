# 78 — Fix the Linux-only `go vet` errors that have been silently skipping the D22 gate in CI since `fd8f838`

**Status:** ready-for-agent (**优先级最高：它是 A44① 的因——D22 门禁从未在 CI 上产出过一次结论**)
**Type:** build/portability defect (tiny diff, large governance consequence)
**Blocks:** 票 71 的 AC 收尾、票 77（新 CI job 不能建立在一个哑步骤上）、`test-core`/`lint` 的可信度
**Blocked by:** nothing — `internal/ball`（票 74 已 done）与 `cmd/wisp`/`internal/proc`（票 66 已 done）现在都空
**Evidence:** registry **A44②**；run `35551819606` / job `106188167868`；票 71 报告

## 已测事实（不要重做定位）
```
internal/ball/statevisual.go:287:17: undefined: mulA
cmd/wisp/slo.go:324:49: undefined: proc.Runtime
```
- **只在 Linux 红**，本机 windows/amd64 不复现 ⇒ 所以它活了很久。今日纯净树
  `GOOS=linux go vet ./internal/ball/` 报**逐字节相同**的错误。
- 根因：`mulA` 只存在于 `internal/ball/renderer_windows.go:368`（该文件 `//go:build windows`），
  而**引用方 `statevisual.go` 没有 build tag**；`proc.Runtime` 同理只在 `internal/proc/boot_windows.go:28`。
- 引入提交：`fd8f838`（**编排者为一个被杀的原型代理做的检查点提交**）与 `00bbb76`。

## 为什么这张票比它的 diff 大小重要得多
CI 的 lint job 里 `go vet` **排在本项目最重要的门禁（D22 七禁令静态扫描）之前** ⇒ 自 `fd8f838` 起
vet 一失败，**扫描步骤被 `skipped`，D22 门从未在 CI 上给出过一次结论**。
票 71 已经把步骤顺序修好（没删步骤、没加 `continue-on-error`），
但**只要这两个 undefined 还在，lint 仍然到不了扫描那一步的结论**。

## What to build
每个未定义符号选**一种**正确形状，并说明为什么（不要为了让 Linux 闭嘴而造空实现骗过 vet 又留下静默行为）：
- **要么**给引用方补上与实现对齐的 build tag / 把该文件拆成 `_windows.go` + `_other.go`；
- **要么**把符号移到无 tag 文件里（如果它本来就不平台相关）；
- **禁止**：写一个返回零值的 Linux stub 让 `cmd/wisp/slo.go` 在 Linux 上"跑起来但什么都不量"——
  那是把"编译错误"换成"假绿"，比原来更糟。若某条路径在非 Windows 上**本就不该存在**，
  就用 tag 把它挡掉，让 Linux 上根本不存在这条路径。

## AC（1:1 裁决表）
- [ ] **AC#1** `GOOS=linux go vet ./...`（主模块）干净——**这就是本票的判据仪器**，
      必须给出真实命令与 exit code。⚠ 变异检验：把任一处修复退回旧形状 ⇒ 该命令必须转红
      （**先 grep 证变异落盘再跑**，还原后证明确实还原）。
- [ ] **AC#2** Windows 侧**零行为变化**：`go test ./internal/ball/ ./internal/proc/ ./cmd/wisp/ -count=2` 全绿，
      且 `gofmt -l` 触及包为空。
- [ ] **AC#3** 一条防回归的机械检查：新增用例或 CI 步骤，使"无 tag 文件引用 windows-only 符号"这类形状
      **在 CI 上必红**（可选实现：`GOOS=linux go vet ./...` 作为 lint job 的一步）。
      判据：它必须自己有一次真红（种子一个故意的引用 ⇒ 步骤红），否则等于没装（A15/A44①）。
- [ ] **AC#4** D22 扫描步骤**第一次产出真实 CI 结论**：交回 push 之后的 run id + job id + 该步骤的结论。
      ⚠ 区分 **failure / cancelled(0 job) / 未跑完**（A44③）：排队被取代的 run **不算样本**。
      你不需要 push（编排者推），把这条留着不勾并写 `next=`，我推完回填。

## 硬规矩
- 只动 `internal/ball/**`、`internal/proc/**`、`cmd/wisp/**`（以及为 AC#3 需要时的 `.github/workflows/ci.yml` 一行）。
  ⚠ `ci.yml` 若与票 71/77 冲突，**先报告再动**（同文件两人同改是假并行）。
- 不许动 `docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`、`internal/tools/**`（他人领地/冻结区）。
- 提交：`git commit -q -F - -- <显式路径> <<'MSGEOF'`（**引号**）；禁 `git add -A`；提交前核
  `git diff --cached --name-only`；禁 `--amend`/`reset`/`rebase`/`stash`（A34）。
- 共享树里**不要跑全仓 `go test ./...`**（会吃到别人未提交的 WIP）；要看 HEAD 的真实状态用
  `git archive HEAD | tar -x -C /tmp/<dir>`，**且绝不在仓库内建 worktree**（A38④）。
- 前 15 次工具调用内必须有一个 checkpoint commit；每次 commit 同步 Status + 勾框 + `next=`。
