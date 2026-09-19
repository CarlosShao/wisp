# T08 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T08-impl 独立）。时间：2026-09-19T23:48:48Z（截止警戒窗口内内联验收）。

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 复跑 | PASS | observe/models `go test -count=1` 全绿；gofmt 漂移已修（internal/models，本报告同步提交） |
| 2 | 采样器口径 | PASS（关键裁定核实） | 门禁单位=私有工作集（NtQuerySystemInformation WorkingSetPrivateSize），commit 并行记录——实测本机 8.7MB 全门禁 PASS（commit 口径会假红 48MB，裁定正确）；fail-closed（零值=error、无样本=FAIL） |
| 3 | 泄漏 fixture | PASS | 100MB 页触摸持有必须翻红否则 FATAL；回落 CheckSettle 要求 10s 内回上限且 FreeOSMemory 计数>0（实测抓到 buffer 未释放真 bug 并已修） |
| 4 | CI 矩阵 | PASS | ci.yml YAML 校验过；5 job 无可跳过；d22scan 实跑全仓 clean（自测旗标名为文档级 MINOR）；lint 含七禁令+emoji 扫描 |
| 5 | compose | PASS | `docker compose config` OK；T14 model-mirror 三服务原样保留 + mock-llm(18080) 接线 |
| 6 | 越界 | PASS | internal/models 仅 gofmt 修复；T14 其他文件零触碰 |

## 登记的风险/移交
1. 本地 sandbox 间歇从 SystemProcessInformation 隐藏自身进程 → 本地 state 采样段不稳；**首次 CI 运行（slo-full on wisp-selfhosted-01）是最终验证点**。
2. Armed/Warm 含模型态现为骨架足迹（JSON posture:"skeleton" 如实标注），票 15/26/33 落地后阈值自动生效。
3. tools/d22scan 的自测触发方式文档化（MINOR，转票 12 gate 检查项）。

**VERDICT: PASS**
