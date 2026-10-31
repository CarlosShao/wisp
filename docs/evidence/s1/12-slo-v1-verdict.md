# 票 12 对抗验收腿 `12-slo-v1` 裁决表 — AC#2（Idle SLO gate 附录格）1:1 复核

**验收人**：12-slo-v1（非实现者对抗验收腿，SPEC-10 §8；⛔ 未落任何 `.go`、未翻任何框、未碰在飞写腿的工作树）
**时间**：2026-10-31 11:31–11:45 +08
**被验收物**：`docs/SLO.md` 附录 C（票 12 AC#2 的解除判据记录）
**判据来源**：票面 `12-cli-text-path-s1-gate.md` AC#2 原文（解除条件＝「票 66 的 AC#1–AC#4 全绿」，票面 :81）
＋ `docs/specs/SPEC-10-testing-acceptance.md` §8 七项 ＋ D32 阈值纪律
**复跑环境**：DESKTOP-LVS7839（与归档同机）· `git archive HEAD`（`2b581be`，2026-10-31 11:34 +0800）
导出至**仓外** `D:\tmp\wisp\20261031\12slov1\tree` 构建 `wisp.exe`（31,945,359B）；仓内写面仅本文件。
**登记一笔**：`.scratch/wisp/probes/orch/20261031/txt/12-slo-v1-dispatch.txt` 未在盘上（该目录只有 `12-r1-dispatch.txt`）；
派单全文随任务下发、判据未因此缺失，按「未定义即停」条款如实登记、不阻塞。

## 0. 结论

**AC#2 能翻。** 解除条件（票 66 AC#1–#4 全绿）已成立且经对抗复核无 blocker；附录 C 的数字锚全部对得上，
且 AC#2 的逐字命令于今日 HEAD 导出树上**实弹复跑 rc=0**。**编排者凭本表翻框（本腿 ⛔ 自翻）**；
翻框措辞请照抄 §5 的三条口径边界。

## 1. SPEC-10 §8 七项检查（逐项「过/不过＋读数」）

| 检查 | 判 | 读数 |
|---|---|---|
| ① stub 扫描 | **过** | `todo!(`/`unimplemented!(` 在 `internal/observe`、`internal/proc` 零命中；无 mock 替真实实现——归档五连跑与本腿复跑量的都是真进程树、真 exit code |
| ② DEFERRED 双向 | **过**（本格射程） | 仪器路径仅 `internal/observe/doc.go:22`（「implemented by ticket 08」，08 已 done、无欠账）与 `internal/proc/doc.go:22`（D41b，对应 SPEC-12 §3 S8 行「更新流程（D41b）+ N-1 回滚」）两处**包范围注记**，非未登记的开放推迟；无 `docs/DEFERRED.md#锚点` 形状的悬空标记 |
| ③ SLO 测量完整性 | **过** | 九门挂 `tree_private_bytes`（`WorkingSetPrivateSize` 树口径，`measurement_basis` 显式落盘）；阈值一字未动——CPU 两侧 limit 同由 `stateCPULimit(st)` 生成（`thresholds.go:133-138`），树内那条降为记录（`:140-152`）；采样未跳：run4 的 1 个 `sample_errors` 照录、39 样本仍出数，五跑全列未挑 |
| ④ 契约偏离 | **过** | 无契约文本改动；句柄 <600 是 `docs/SLO.md` 编排者裁定 1（2026-09-19）书面口径，25MB/0.5%/GDI<200 与 D32 原文逐字一致 |
| ⑤ D22 七禁令扫描 | **过** | `cd tools/d22scan && go run . -root ../..`（CI 同形）于 HEAD 导出树 → **rc=0 clean**，自报 internal 695＋cmd 161 Go files（含 ban#8） |
| ⑥ 场景↔切片双向 | **过** | 九门 ↔ D32 16.3.2 `Sleeping` 行（内存/CPU/GDI/句柄/goroutine/写盘/TCP）＋回落行（≤10s 回 cap＋FreeOSMemory>0）逐项 1:1；S1 done 判据的 idle gate 行即本格 |
| ⑦ 失败预演＋缺口审计 | **过** | 聚合门自带必须翻红的 leak 自检（归档 5 份全 `flipped_to_fail=True`）；本腿另撞出一次计划外阴性对照——子进程缺 DLL 时门按 INSTRUMENT BROKEN FATAL、零数字落盘（fail-closed 生效）；缺口审计＝票 66 裁决表（验收人非实现人）＋残余明写四条 |

## 2. 数字真实性（锚＋复跑读数）

- **C.1 五连跑**锚＝`docs/evidence/s1/66/66b-nokeeper-{1..5}.json`，逐字段核对：主体中位 4.39/4.35/4.31/4.27/**4.21**MB ✓、
  句柄 189/183/183/189/189 ✓、run4＝39 样本＋1 sample_error ✓、主体（树外）CPU 0.0130/0.0000×4% ✓、
  观测者（树内）CPU 0.6739/0.5963/0.6350/0.6872/0.9318% 全 `gate=false` ✓ —— 与附录 C.1 表**逐格相符**。
- **C.3** 锚＝`66-ac3-30s.json`：同窗树外 0.0173%（`gate=true`）＋树内 0.7629%（`gate=false`、`observer_cost:true`、note 在案）✓。
- **C.5** 锚＝`66-full-subset-slo-report.json`：六态 exit=0/pass、settle peak 132,636,672B=126.5MB → **back_within_cap_ms=261**、
  **free_os_memory_count=2**、final 8,916,992B=8.50MB ✓。
- **C.4** 锚＝`66-smoke-run-{1..5}.log`：五份尾部全 `all_pass=True`＋`leak exit=1 flipped_to_fail=True` ✓。
- ★**实弹复跑（本腿，2026-10-31 11:3x）**：AC#2 逐字命令 `wisp slo -state Sleeping -seconds 10 -interval-ms 250`（零 keeper、
  桌面前后 `wisp/balldebug` 均为 0）→ **rc=0、40/40 样本、0 error、basis=out-of-tree、CPU 门 0.026% ≤0.5%、内存 4.8MB、
  九门 9 行（7 门全过）**。JSON＝`D:\tmp\wisp\20261031\12slov1\out\12slov1-sleeping-1.json`。
- ★**A30 欠账变异复测**（票 66 裁决表内唯一〔仅代理自述〕行就此补上独立复现）：导出树上把
  `WalkSystemProcesses`（`internal/proc/systemprocs_windows.go:62`）退回「先测 `next==0` 再 `fn(e)`」旧序 →
  `go test ./internal/proc/` **rc=1**，4 红：`TestParseSystemProcessesKeepsLastSnapshotEntry`、
  `TestExternalSamplerTracksSubjectGrowth`、`TestWalkSystemProcessesStopsAtTerminator`、`TestParseSystemProcessesSkipsPIDZero`；
  还原后 MUTATION 计数=0、`./internal/proc/ ./internal/observe/` **rc=0**（ok 2.113s / ok 5.097s）。
- **复跑干扰如实记**：smoke 今日 3 发＝1 发我方环境错（新 pwsh 进程未带 sherpa DLL 的 PATH → 四个 `wisp.exe` 全 0xC0000135 →
  门按「leak 不翻红」FATAL——fail-closed 正确，反成阴性对照）；2 发被票 134 contention 预检拒（全机 CPU 96% ≥50%，
  rc=0、零数字落盘）——**无一次产出受扰数字**。AC#2 直跑那一发起测瞬间有 1 个 go 工具链进程在飞，读数仍 20 倍低于门，照实记录不抹。

## 3. 「树外口径当门」的合法性（照裁，不重裁）

票 66 书面裁定在案：裁决表（`docs/evidence/s1/66-adversarial-acceptance.md:29`）AC#3 行逐字要点——
「**D32 的 `Sleeping` CPU ≤0.5% 阈值一字未动**：两侧都仍按 `<=0.5%` 判并写 pass/fail，改的只是**哪一条有资格当门**——
这个区分由 `internal/observe/observer_cost_test.go` 4 条用例钉住」。机器可核三件：`thresholds.go:140-152` 注释逐字引
「Ticket 66 ruling 2」＋附录 B.2 的 ≈1.3ms/读；`observer_cost_test.go` 在盘；JSON 里树内行带 `observer_cost:true`＋note、
树外行 `gate=true`（本腿复跑同形状）。台账 A15 结案行（`pending-and-issues.md:3213`）＝「已解决（票 66，AC#3）」。
⇒ **已被裁过当门，本腿照裁**。C.4 smoke 主体 0.0000% 正是这条树外口径的读数，引用合法。

## 4. `docs/SLO.md:77` 的「待票 27」行 — 无耦合

该行＝§3 模型常驻表的「标点恢复（CT-Punc 类）未测 → 待票 27（S4）」，属 **Warm 态模型加载**口径；
九门不含任何模型行；票 12 AC#2 原文与其解除条件（票 66 AC#1–4）均不引票 27；`grep 票 27 docs/SLO.md` 全文仅此一行。
**12-r1 的「无耦合」核成立**，本格不受它牵制。

## 5. 残余具名（⛔ 本格射程）

1. **AC#8**（球态截图人工签收，D29）未勾，owner-only——本格翻转不使票 12 整体 done。
2. **AC#9**（2026-09-20 编排者新增，registry A25：回复后 3s 回落 ≤3000ms、单调钟、≥5 样本全报）未勾，归写腿排队。
3. 仪器残余：`NtQuerySystemInformation: buffer never sufficient`（run4 那发 sample_error，票 66 登记候选新项）未修——
   影响＝单样本丢弃、九门仍出数，不阻本格。
4. C.5 口径边界照旧有效：六态 posture=skeleton，**六态真测归票 15/28/33/36**——引用 C.5 时不得读成「六态已过」。

## 6. 判语

**AC#2 能翻**：解除条件成立＋§1 七项全过＋§2 数字锚逐数相符＋树外口径有书面裁定＋今日实弹复跑绿。
本格内无缺项。翻框归编排者；建议翻框措辞明写三条边界：①CPU 门读数是**树外口径**（票 66 裁定 2，`observer_cost_test.go` 钉住）；
②posture=skeleton，D32 六态真测未闭环（归 15/28/33/36）；③AC#8/AC#9 与仪器新项仍开着，票 12 整体未到 done。

## 附：本腿复跑产物（仓外，`D:\tmp\wisp\20261031\12slov1\`）

`head.tar`＋`tree\`（HEAD `2b581be` 导出树）· `out\12slov1-sleeping-1.json`（AC#2 复跑原件）·
`out\smoke{2,3}\slo-no-conclusion.json`（contention 拒采记录）· `logs\{build,ac2-rerun-1,smoke-rerun-1..3,d22scan-export,mut-proc,green-baseline}.*` ·
`systemprocs_windows.go.pristine`（变异复测的还原源）。
