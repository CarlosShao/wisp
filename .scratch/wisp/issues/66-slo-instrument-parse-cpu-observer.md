# 66 — SLO 判据仪器返工：解析丢末项 + CPU 门的观测者自成本（票 12 桌面跑的连带发现）

**Status:** in-progress
**Claimed by:** agent-ticket66b
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

- [x] **AC#1 解析末项**：`parseSystemProcesses` 顺序修正；一条**回归用例**构造「目标 pid 恰为快照链最后一项」
  的缓冲并断言它在 map 里。⚠ 必须证明这条用例咬住缺陷：把函数退回旧实现，该用例**必须转红**
  （变异检验，`docs/evidence/s1/66-mutation-*.md` 留原始输出）。
- [x] **AC#2 无 keeper 可用**：`wisp slo -state Sleeping -seconds 10 -interval-ms 250`
  在**不启动任何 keeper 进程**的情况下 `exit≠2` 且九项门全出数；**连续 5 次**全部如此（一次成功不算，
  因为旧缺陷本来就是按进程创建顺序随机的）。把 A.6#1 里的 keeper 绕法从文档删掉。
- [x] **AC#3 CPU 树外口径**：`wisp slo` 获得树外测量能力（父读子 pid，复用 `cmd/balldebug/diff_windows.go`
  里已验证正确的遍历顺序；**不得复制粘贴出第二份实现**，须抽公共函数）。
  完成判据：同一 30s 窗口内，树外口径的 `Sleeping` CPU 均值 ≤0.5%，**且**树内口径在相同条件下仍能测出 >0.5%
  ——两条同时成立才证明"差的就是观测者"，而不只是换了个地方读。
- [x] **AC#4 CI 不再随机红**：`scripts/slo-check.ps1 -Subset smoke` **本地连续 5 次** exit=0；
  state 段不再引用树内 CPU 当门（引用位置须 diff 给复核者看）。
  ⚠ 复跑时挖出**第三起仪器缺陷**并同批修掉：`slo-check.ps1:129` 的聚合式在 `Set-StrictMode 2.0` 下
  恰好于「全过」路径抛 `.Count` 异常（⇒ `exit=1` 且不写报告，CI 连 artifact 都没有）；
  A14 修好前它被「所有 state 都 fail ⇒ 过滤器命中 ≥2 项」掩盖。数字与阳性对照见 `docs/SLO.md` 附录 C.4。
- [x] **AC#5 追溯票 08**：在票 08 的 ticket 文件 Progress log 追加一行，说明它归档的
  `build/slo/slo-report.json` 是 `exit=2`（state 口径从未出数），**不改它的 AC 勾选状态**
  （那是编排者的权，见纪律）。修好后重跑一次 full 子集，报告归档进 `docs/evidence/s1/66/`。
- [ ] **AC#6 文档同步**：`docs/SLO.md` 附录 A.1/A.4 与附录 B 的数字若因本票改动（例如树外口径新增），
  以**叠加新小节**的方式写，不得改写已录的原始样本表；A14/A15 在 `docs/reports/pending-and-issues.md`
  里移入「已解决」并附本票 commit 号。
- [ ] **AC#7 门禁**：`gofmt -l` 触及包为空、`go vet` 干净、`go test -count=2 ./internal/proc/ ./internal/observe/`、
  `-race` 同包跑一次；`tools/d22scan` 与静态禁用模式扫描 clean。

## Progress log（每次 commit 追加一行，格式 `- [ISO-Z] agent=... did=...`）
- [2026-09-20T15:26Z] agent=agent-ticket66b did=AC#2 复测收口 + AC#4 + AC#6 附录 C。**先自曝一次测量事故**：第一次跑 AC#4 五连时我把脚本调用接在 `| tee | grep` 后面，`$?` 取到的是 `grep` 的退出码 ⇒ 屏幕上五个「exit=0」全是假的（真退出码 1，聚合报告根本没写出来）；发现方式是 `ls` 找不到 `slo-report.json`。修正后逐字重跑，退出码直接取自脚本、不经管道。**新缺陷（第三起，A14 修好后才露头）**：`scripts/slo-check.ps1:129` 的 `($results | Where-Object { -not $_.pass }).Count` 在 `Set-StrictMode 2.0` 下**恰好在「全过」那条路上**抛 `PropertyNotFoundException`（零命中返回 `$null`、单命中返回标量，都没有 `.Count`）⇒ 四步检查全 pass 仍 `exit=1` 且不写报告，而 CI 那两步（`ci.yml:168/:198`，无 `continue-on-error` + `if-no-files-found: error`）连 artifact 都没有；它此前没露头是因为 A14 让每个 state 都 `exit=2` 全 fail ⇒ 过滤器命中 ≥2 项才有 `.Count`。⇒ **B.3 的预告被这次复跑反向加强：只修 A14 不是「随机红」而是「必红且无报告」**。修法 `@(...)` 强制成数组（同批 commit），阳性对照（StrictMode 下喂 0/1/2 项失败）= `allPass=True/False/False`，证明不是把门焊死成绿。修复后 `-Subset smoke` **5/5 exit=0**（15:14:32Z–15:18:46Z，五份日志归档 `docs/evidence/s1/66/66-smoke-run-{1..5}.log`，崩溃原文 `66-smoke-pre-fix-crash.log`），十个 state 样本的树外 CPU **全部 0.0000%**，树内记录 0.1612–1.0322%（再次印证那条门在树内口径无定义）。AC#2 侧：`-state Sleeping -seconds 10 -interval-ms 250` **零 keeper 连跑 5 次 exit=0/0/0/0/0**，九门每次 9 行，主体 CPU 0.0130/0.0000/0.0000/0.0000/0.0000%、观测者 CPU 0.6739/0.5963/0.6350/0.6872/0.9318%（全 >0.5%），run4 有 1 次 `sample_error`（`NtQuerySystemInformation: buffer never sufficient (last 1090464 bytes)`，200 次读里 1 次，**不是 A14**，登记为候选新项未修）。AC#5 的复跑腿：`-Subset full -SecondsPerState 6` **exit=0 / all_pass=true**，六态 24/24 样本全出数，报告归档 `docs/evidence/s1/66/66-full-subset-slo-report.json`；⚠ 六态 `posture` **全是 skeleton** ⇒ 这张表只证明「state 口径跑得通且 skeleton 达标」，**不等于** D32 六态表达标，边界已写进附录 C.5。AC#6：新增 `docs/SLO.md` **附录 C**（C.1-C.7，A/B 的原始样本表一字未改），A14/A15 移入 registry 已解决段在下一行 commit。测量六次前后 `tasklist|grep -iE "wisp|balldebug"` 全空。next=AC#7 门禁 + registry 搬迁。
- [2026-09-20T15:18Z] agent=agent-ticket66b did=AC#5 追溯票 08（只追加一行事实，未碰它的 AC 勾选与 Status 字段）：逐字段读它归档的 `build/slo/slo-report.json` 原件（`generated_at 2026-09-19T23:35:49Z`）确认 `Sleeping exit_code=2 / report=null`、`Warm exit_code=2 / report=null`、`settle exit_code=1 / pass=false`、`all_pass=false` —— 即 state 口径九门当时一门都没产出样本窗，红因是仪器（A14/A15）不是预算。票 08 文件里补的那行同时交代了三处修复归属（`86e868d` / `00bbb76` / 本次 ps1 崩溃修复）与重跑数字的位置（`docs/SLO.md` 附录 C）。同批把 `**Claimed by:**` 更新为 agent-ticket66b（承接被中断的 agent-ticket66），并按已完成证据勾 AC#3（同一 30s 窗：树外 CPU **0.0173%** vs 树内 **0.7629%**，两侧都按 `<=0.5%` 判、只有树外那条当门，原始 JSON `docs/evidence/s1/66/66-ac3-30s.json`）。next=AC#4 五连跑收尾。
- [2026-09-20T15:06Z] agent=agent-ticket66b did=AC#1 变异检验**独立复测**（编排者要求换人重做，不引用上一位代理的自述）：把 `WalkSystemProcesses` 退回「先测 `next==0` 再解码」，先 `grep -n MUTATION-66B` + `git diff --stat`（1 file changed, 6 insertions）证明变异真落盘，再 `go test -count=1 -v ./internal/proc/` → **rc=1，6 红 / 24 绿 / 1 SKIP**（比上一轮多 3 条：`TestSystemProcessSnapshotSeesSelf` 这次也红，`system snapshot does not contain the sampling process 12328 (513 entries)`；外加 AC#3 的两条树外采样用例 `system snapshot does not contain subject <pid>` —— 说明「换个地方读」并不天然免疫 A14，靠的是全仓唯一那份走链实现）。撤变异后 `grep -c MUTATION`=0、`git diff` 空、`go test ./internal/proc/ ./internal/observe/` rc=0 全绿。原始输出**追加**进 `docs/evidence/s1/66-mutation-parse-order.md`（未覆盖上一轮）。顺带完成 AC#2 复测的采样腿：`build/wisp.exe`（23:03 重建）`slo -state Sleeping -seconds 10 -interval-ms 250` **零 keeper 连跑 5 次 exit=0/0/0/0/0**，九项门每次各出 9 行、40 样本、run4 有 1 次 sample_error（39 样本，九门仍出数），主体（树外口径）CPU 0.0130 / 0.0000 / 0.0000 / 0.0000 / 0.0000%、观测者（树内，`gate:false`）0.6739 / 0.5963 / 0.6350 / 0.6872 / 0.9318%；测量前后 `tasklist|grep -iE "wisp|balldebug"` 均空。JSON 归档 `docs/evidence/s1/66/66b-nokeeper-{1..5}.json`。剩余：AC#4 slo-check 连跑 5 次、AC#5 票 08 追溯 + full 归档、AC#6 附录 C + A14/A15 移入已解决、AC#7 完整门禁。
- [2026-09-20T15:05Z] agent=agent-ticket66 did=AC#3 代码落地（同一 30s/同间隔的双子进程拓扑）：`wisp slo -state X` 现在在 Job 里起两个子进程——被测主体（boot+settle 后什么都不做，由 `proc.ExternalSampler` 树外读取，九项门全在它身上）与被测量化主体（跑遗留的树内自采样，报告塞进 `observer`，CPU 行 `observer_cost: true` 且 `gate: false`），观测者自己不再进被测集合；`cmd/balldebug/diff_windows.go::privateWorkingSetFor` 改调 `proc.SystemProcessSnapshot`，全仓只剩一份 SYSTEM_PROCESS_INFORMATION 走链（A14 就是第二份复制的顺序错）。阈值零改动：`cpu_percent_all_core` 两侧都仍按 `<=0.5%%` 判并写 `pass/fail`，只是树内那条不再当门（`internal/observe/observer_cost_test.go` 4 条用例钉住「只许 CPU 行动、limit 字符串两侧必须相等、非 CPU 行的 gate 不许随口径变化」）。冒烟（5s/250ms，同窗口同一对主体）：树外 CPU **0.0000%%** / 树内 CPU **0.5659%%**、`exit=0`；顺带发现树内私有工作集中位 8.38MB vs 树外 4.33MB，疑似采样器 1MiB 级扫描缓冲把内存行也污染了（SLO.md §7 早年记过的坑），待 30s 正式样本核实。另：先前 `go test ./internal/observe/` 因 `internal/llm/adaptertest/mockllm.go:68` 裸 `go func(` 转红，已确认早于本票（`git show 2825f28:...` 第 68 行即命中），且已被票 67 的 `7c9256b` 修掉，本票门禁不再受阻。剩余：AC#3 正式 30s 样本、AC#4 slo-check.ps1、AC#5 票 08 追溯 + full 归档、AC#6 附录 C/A14-A15 移已解决、AC#7 完整门禁。
- [2026-09-20T14:38Z] agent=agent-ticket66 did=AC#1 变异检验 + AC#2 无 keeper 连跑 5 次。变异：把 `WalkSystemProcesses` 退回「先测 `next==0` 再解码」，`go test -count=1 -v ./internal/proc/` → rc=1，3 条新用例全红（`parsed 4 of 5 snapshot entries: [1337 60060 4 4004]`／`walker delivered [100], want [100 200]`／`parsed zero entries`），包内其余 23 条仍绿；撤变异后 rc=0 全绿（grep -c MUTATION=0 已核）。原始输出 `docs/evidence/s1/66-mutation-parse-order.md`。AC#2：`./build/wisp.exe slo -state Sleeping -seconds 10 -interval-ms 250`（零 keeper）连跑 5 次 exit=1/1/1/1/1（≠2，九项门各出 9 行、40 样本、0 sample_errors），CPU 树内读数 0.6604/0.7625/0.7623/0.6356/0.6356% —— 即缺陷②本身，正是 B.3 预告的「只修① ⇒ 门随机红」，故 AC#3/#4 必须与本 commit 集同批。编排者那条 exit=2 的原 repro（`-seconds 5`）现在 exit=0（`docs/evidence/s1/66/66-repro-a1.json`）。A.6#1/#2 的 keeper 两行已从 docs/SLO.md 删掉并改为无 keeper 形态；无 keeper 的 `-settle` 复跑 exit=0（回 cap 262ms、FreeOSMemory=2、末值 8.5MB）。测量前后 `tasklist|grep -iE "wisp|balldebug"` 均为空。剩余：AC#3 树外口径（含 balldebug 改用公共快照函数）、AC#4 slo-check、AC#5 票 08 追溯 + full 子集归档、AC#6 附录 C + A14/A15 移入已解决、AC#7 门禁。

- [2026-09-20T14:24Z] agent=agent-ticket66 did=AC#1 主体：`parseSystemProcesses` 改为「先解析当前项并入 map、再测 `NextEntryOffset==0`」，并把这条链的唯一解码器抽成 `internal/proc/systemprocs_windows.go::WalkSystemProcesses`（`SysProcSample`/`SystemProcessSnapshot` 一并导出，供 AC#3 的树外读取与 `cmd/balldebug` 复用，不再养第二份实现）；新增回归用例 `internal/proc/systemprocs_windows_test.go` 7 条，其中 `TestParseSystemProcessesKeepsLastSnapshotEntry` 手工构造「目标 pid 恰为链末项」的缓冲并逐项断言 handle 数/私有工作集/线程数/user+kernel 时间（不是只断言 `len(out)`）。测量：`go test -count=1 -run 'TestParseSystemProcesses|TestWalkSystemProcesses|TestSystemProcessSnapshot' ./internal/proc/` → 7/7 PASS；`go build ./...` OK；`go test -count=1 ./internal/proc/` ok 0.464s；`gofmt -l internal/proc/` 空；`go vet ./internal/proc/` 干净。剩余：变异检验（退回旧实现 → 新用例必须转红，原始输出存 `docs/evidence/s1/66-mutation-*.md`）、AC#2-AC#7 全部未动。
