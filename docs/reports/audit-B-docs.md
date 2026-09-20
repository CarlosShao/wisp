# AUDIT-B — 文档与证据审计（13 张 done 票 01–09、13、14、17、18）

> 执行者：Audit-B（独立审计代理）。时间：2026-09-20 08:4x–08:5x（08:50 硬停，到点保存）。
> 方法：全量浅查（证据存在性、AC 勾选、进度日志标记、HANDOVER/registry 交叉核对）+
> 深查 4 张（04/07/09/13，每张随机抽 2 条证据引用到盘上核实）。
> 限制：受 08:50 硬停约束，未复跑全套件、未做 mutation 实验；单点核实以文件系统与现场单测为准。

## 一、证据存在性（全量）——PASS

13 张 done 票的验收报告全部在盘：
- s0：01/03/04/06-adversarial-acceptance.md + 02-spike-report.md + 02-adversarial-acceptance.md
- s1：05/07/08/09/17/18-adversarial-acceptance.md（另 c21-native-tokens.md、ball-states/ 20 张 PNG）
- s2：13/14-adversarial-acceptance.md
无缺失。

## 二、验收报告真实性抽查（深查 04/07/09/13）——所抽引用全部属实

- **票 04**：`TestAdvPanicWedge` 存在（internal/memory/dao_test.go:524）且本审计**现场实跑 PASS**（0.088s）；go.mod `modernc.org/sqlite v1.59.0` 直接依赖属实（复验声称）。✓
- **票 07**：docs/evidence/s1/ball-states/ 实有 20 张状态 PNG；internal/statemachine/table_test.go 实测 **314 行**、machine_test.go 195 行，与报告引用一致。✓
- **票 09**：internal/llm/testdata/golden/ 实有 10 个 .sse fixture；`TestMockllmGoldenByteIdenticalEvents` 存在（internal/llm/openaichat/mockllm_integ_test.go:176）；报告含自曝 MAJOR-1（sameEvent 比较失实）并给出 file:line 与最小修复——报告具真实对抗性。✓
- **票 13**：internal/audio 实有 **26 个**测试函数，与报告"26/26 race green"吻合（并纠正了交接声称的 21）；报告 §6 六条 AC 逐条裁决表引用具体测试名与 file:line。✓
- 唯一引用失效：票 07 报告指向不存在的 `docs/reports/pending-human-review.md`（见 MINOR-6）。

## 三、AC 覆盖核对（全量浅查）

13 张票票据内 Acceptance criteria checkbox **全部为 `- [ ]`（合计 71 条，勾选 0 条）**，票内无逐条裁决行；裁决在验收报告中，且粒度不一（04/09/13 有逐条 AC 表；05/06/07/17/18/14 为审计项表而非逐 AC 对照）。

## 四、HANDOVER.md 交叉核对

§3 完成度 13/61 ✓（盘上恰 13 个 `-done` 文件）；§4 无在途 ✓（git status 干净，HEAD 30adff5）；§5 行动清单（10 → 19∥20 → 12；15/16）与票据实际状态一致 ✓；R1/R2 在 registry ✓；H2 声称的 cmd/llmrecord 实际存在 ✓。

## 五、进度日志完整性

13 张票全部条目均带 `agent=` 标记与 ISO 时间戳；12/13 有显式 handoff 条目。发现乱序/缺项见下。

## 发现清单

### BLOCKER
（无）

### MAJOR
1. **[全部 13 张 done 票] AC checkbox 未勾选即标 done**——票据内 71 条验收标准全部 `- [ ]`、0 勾选，票内无逐条裁决行；最尖锐实例：**票 18** 的 AC"Resolver idempotence + perf ≤1ms 暖缓存"在报告"剩余"段自认 DEFERRED（未做 bench），票仍标 done。触发本审计规则"未勾选却标 done → MAJOR"。修复建议：验收时逐条把 `- [ ]` 翻成 `- [x]` 并附报告锚点；18 的 perf AC 补 bench 或在票据标注 DEFERRED+归属。

### MINOR
2. **[票 02/17/18] 进度日志时间戳乱序**——append-only 声称 newest-last，但 02 的 orchestrator VERDICT(09:31:25Z) 记录在 impl handoff(09:40:00Z) 之后；17 的 PASS(00:27:30Z) 在 impl commit(00:35:00Z) 之后；18 的 PASS(00:27:30Z) 在 impl 条目(00:31:00Z) 之后（疑为回填未排序）。
3. **[票 17/18] 报告路径乱码**——进度日志写 `docs/evidence/s1/.scratch/wisp/issues/17-*.md`/`18-*.md`，真实路径为 docs/evidence/s1/{17,18}-adversarial-acceptance.md。
4. **[票 04] 票头时间戳粘连**——"Last update: 2026-09-19T10:03:56ZT09:52:40Z"（双时间戳）。
5. **[票 18] 缺显式 handoff 条目**——T18-impl 末条为 `next=none-resolver-done`，无 `did=handoff-to-orchestrator`（13 张中唯一）。
6. **[HANDOVER §8 vs pending-and-issues.md] H 编号冲突**——HANDOVER 的 H5=minisign 密钥仪式、H6=票 16 真机三场景；registry 的 H5=TTS 选型（matcha P3）、**无 H6 条目**。下一会话按编号引用会张冠李戴；票 16 真机三场景是 S2 验收必需的人工件，registry 漏登。
7. **[票 07] 验收报告引用失效文件**——"已登记 docs/reports/pending-human-review.md"全盘无此文件；H1 内容实际在 pending-and-issues.md"待人工审核"节（内容未丢失，故 MINOR；若按"证据文件缺失"口径可升 MAJOR）。

## 已核实为准确（抽样通过项备案）

07 报告 line-count 引用、09 报告 file:line 引用与自曝缺陷、13 报告测试计数、04 的 TestAdvPanicWedge（现场实跑 PASS）、HANDOVER §3/§4/§5、H2→cmd/llmrecord、R1/R2 存在、ball-states 20 PNG、golden×10。

## 结论

文档-证据体系整体真实：13/13 验收报告在盘、深查 4 张所抽 8+ 条引用全部属实、HANDOVER 关键声称与盘面一致。主要缺口是**票据侧 AC 勾选纪律**（全量未勾选、18 一条 AC 未完成即 done）与**编号/路径类笔误**。

**AUDIT-B VERDICT: ISSUES: 7**
