# 66 — SLO 判据仪器返工：解析丢末项 + CPU 门的观测者自成本（票 12 桌面跑的连带发现）

**Status:** in-progress
**Claimed by:** agent-ticket66
**Last update:** 2026-09-20
**Blocked by:** —（与票 12/20/21 无代码交集；本票只动 `internal/proc`、`internal/observe`、`cmd/wisp/slo*.go`、`scripts/slo-check.ps1`）
**Parallel slots:** ≤1 sub-agent（**测量类独占**：本票的验收要求真机跑采样，不能与任何其他跑测的代理并发）
**Spec refs:** D32 16.3.2（`Sleeping` 行 CPU ≤0.5% 全核，**阈值不得改动**）、C30（Job 进程树私有工作集口径）、D42#10、SPEC-10 §3
**登记项:** A14（解析丢末项）、A15（CPU 门自成本）；证据 `docs/SLO.md` 附录 A（代理那次）+ 附录 B（编排者复跑）

## 为什么要这张票

票 12 的桌面 SLO 跑挖出两起**判据仪器自身**的缺陷，我独立复跑后**两起都成立**，
但第二起的**解释**要换（更强的形式）。合起来的效果是：

> **D32 `Sleeping` 行的 CPU 门今天无法判定；而 state 口径的 SLO 门今天根本没产出过一份样本窗，
> CI 的两处 SLO 步骤（`.github/workflows/ci.yml:166` smoke / `:196` full）在无 `continue-on-error` 的情况下随机红。**

这不是功能缺陷，是**「我们用什么证明自己达标」的缺陷**。它的优先级高于票面上任何功能框：
仪器坏着的时候，后面每一张需要 SLO 数字的票（15/26/31/32 都要回填 `Sleeping`/`Warm` 行）
都会继承一个不可判的门。

## 缺陷① 解析器丢快照末项（A14，已复现）

`internal/proc/treemetrics_windows.go:236-262`：循环先
`next := ...; if next == 0 { break }`（:240-243），之后才解析当前项并入 map（:247-260）。
变长快照链的**最后一项** `NextEntryOffset` 恰为 0 ⇒ **永远进不了 map**。
新建进程排在 `ActiveProcessLinks` 末尾 ⇒ 「自己量自己」的采样进程常常正是被丢掉的那个 →
60 次重试（≈6s）耗尽 → fail-closed `system snapshot does not contain the sampling process`，`exit=2`。

**兄弟实现是对的**：`cmd/balldebug/diff_windows.go::privateWorkingSetFor` 先比 pid、再在 `next==0` 处 break。
**修法就是把「解析当前项 + 入 map」放到 `next == 0` 判断之前**（顺序调整，不是新增逻辑）。

## 缺陷② CPU 门测的是观测者自己（A15，已复现，解释形式已更正）

采样器运行在**被测树内部**。我用 8 个自跑样本把绝对成本钉住：

| 间隔 | 跳动样本 CPU 读数 | 折算绝对成本/次读 |
|---|---|---|
| 250ms | 0.5062 – 0.5276% | **≈1.28 – 1.32ms** |
| 2000ms | 0.0649 / 0.0651 / 0.0652% | **≈1.298 – 1.304ms** |

两个独立配置量出同一个绝对值 ⇒ 这就是每次 `ReadTree` 的 CPU。
⇒ **在 250ms 间隔下，光是"完成一次读"就是 0.52% > 0.5% 门**：门限与仪器最小可分辨量撞在同一数量级，
这条门在这个口径下**没有定义**（不是"噪声大"）。
同配置重复跑 PASS/FAIL 随机翻转（我的 5 个 10s/250ms 样本：0.0387 / 0.5838 / 0.5837 / 0.5053 / 0.2076%）。
**被测进程本身是清白的**：40 样本里 37 个 CPU 恰为 `0.0000`，与 `cmd/balldebug -diff` 树外实测球体
`Sleeping` **0.000%** 互证。

⚠ 票 12 记录里「CPU ∝ 读次数」那个解释形式**不成立**（40 次读反而比 5 次读更省），
不要照它设计修法（例如"降低采样频率"就不解决问题——真正要动的是**观测者所在的进程**）。

## ⚠ 两起缺陷必须一起修，只修①会让情况变坏

只修① ⇒ `slo-check.ps1` 从「今天必红」变成「**约一半概率随机红**」。随机红的门会训练所有人忽略它，
比确定性地坏更糟。所以本票的 AC 把「CPU 行离开树内口径」和「解析修复」绑在同一个 commit 集里。

## 编排者已作出的裁定（代理不得越界解释）

1. **D32 的 0.5% 阈值不改。** 要改的是量法，不是门。
2. **CPU 行的产品侧口径 = 树外测量**（父进程读子 pid，即 `cmd/balldebug -diff` 已有的做法）。
   把 `wisp slo` 的树内 CPU 数字**降级为记录项**需要书面批准——**本票就是这份批准**：
   修好①后，`slo-check.ps1` 的 state 段**不得**再拿树内 CPU 当门；树内 CPU 仍须写进 JSON 报告并标 `observer_cost: true`。
3. **票 12 AC#2 保持未勾**，不因本票完成而自动补勾。

## 验收标准（AC）

- [ ] **AC#1 解析末项**：`parseSystemProcesses` 顺序修正；一条**回归用例**构造「目标 pid 恰为快照链最后一项」
  的缓冲并断言它在 map 里。⚠ 必须证明这条用例咬住缺陷：把函数退回旧实现，该用例**必须转红**
  （变异检验，`docs/evidence/s1/66-mutation-*.md` 留原始输出）。
- [ ] **AC#2 无 keeper 可用**：`wisp slo -state Sleeping -seconds 10 -interval-ms 250`
  在**不启动任何 keeper 进程**的情况下 `exit≠2` 且九项门全出数；**连续 5 次**全部如此（一次成功不算，
  因为旧缺陷本来就是按进程创建顺序随机的）。把 A.6#1 里的 keeper 绕法从文档删掉。
- [ ] **AC#3 CPU 树外口径**：`wisp slo` 获得树外测量能力（父读子 pid，复用 `cmd/balldebug/diff_windows.go`
  里已验证正确的遍历顺序；**不得复制粘贴出第二份实现**，须抽公共函数）。
  完成判据：同一 30s 窗口内，树外口径的 `Sleeping` CPU 均值 ≤0.5%，**且**树内口径在相同条件下仍能测出 >0.5%
  ——两条同时成立才证明"差的就是观测者"，而不只是换了个地方读。
- [ ] **AC#4 CI 不再随机红**：`scripts/slo-check.ps1 -Subset smoke` **本地连续 5 次** exit=0；
  state 段不再引用树内 CPU 当门（引用位置须 diff 给复核者看）。
- [ ] **AC#5 追溯票 08**：在票 08 的 ticket 文件 Progress log 追加一行，说明它归档的
  `build/slo/slo-report.json` 是 `exit=2`（state 口径从未出数），**不改它的 AC 勾选状态**
  （那是编排者的权，见纪律）。修好后重跑一次 full 子集，报告归档进 `docs/evidence/s1/66/`。
- [ ] **AC#6 文档同步**：`docs/SLO.md` 附录 A.1/A.4 与附录 B 的数字若因本票改动（例如树外口径新增），
  以**叠加新小节**的方式写，不得改写已录的原始样本表；A14/A15 在 `docs/reports/pending-and-issues.md`
  里移入「已解决」并附本票 commit 号。
- [ ] **AC#7 门禁**：`gofmt -l` 触及包为空、`go vet` 干净、`go test -count=2 ./internal/proc/ ./internal/observe/`、
  `-race` 同包跑一次；`tools/d22scan` 与静态禁用模式扫描 clean。

## Progress log（每次 commit 追加一行，格式 `- [ISO-Z] agent=... did=...`）

- [2026-09-20T14:24Z] agent=agent-ticket66 did=AC#1 主体：`parseSystemProcesses` 改为「先解析当前项并入 map、再测 `NextEntryOffset==0`」，并把这条链的唯一解码器抽成 `internal/proc/systemprocs_windows.go::WalkSystemProcesses`（`SysProcSample`/`SystemProcessSnapshot` 一并导出，供 AC#3 的树外读取与 `cmd/balldebug` 复用，不再养第二份实现）；新增回归用例 `internal/proc/systemprocs_windows_test.go` 7 条，其中 `TestParseSystemProcessesKeepsLastSnapshotEntry` 手工构造「目标 pid 恰为链末项」的缓冲并逐项断言 handle 数/私有工作集/线程数/user+kernel 时间（不是只断言 `len(out)`）。测量：`go test -count=1 -run 'TestParseSystemProcesses|TestWalkSystemProcesses|TestSystemProcessSnapshot' ./internal/proc/` → 7/7 PASS；`go build ./...` OK；`go test -count=1 ./internal/proc/` ok 0.464s；`gofmt -l internal/proc/` 空；`go vet ./internal/proc/` 干净。剩余：变异检验（退回旧实现 → 新用例必须转红，原始输出存 `docs/evidence/s1/66-mutation-*.md`）、AC#2-AC#7 全部未动。
