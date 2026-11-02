# 402 — panel/WebView2 域：./... 整包扫 cmd/wisp 死形修复（401-v1 建议维持）

- **Status**：**已结案 done**（2026-11-02 19:0x 编排者据非实现者 `402-c2` 重评「**同意结 402**」结案：AC14 push 维度有效修（`a07f685eb` 6/6 green＋mutant 红句注册）；AC#3 Weixin 前台锁显式归因；★AC#1/AC#2「静默杀/整包零 AV」经 **正确 sherpa 三层 staging 现跑 581/575 RUN 零 AV 完整跑完**证为 **cmd/wisp 真窗族平台事实＋缺 DLL 的 staging 假象（非 402 缺陷）**；402-v1 之 FAIL 系 staging 未摆对所致；判语 `docs/evidence/s1/402-c2-verdict.md`＝44afb4d03）。**登记待办**：真窗族"整包零 AV"需机主空闲安静窗口复跑（pending_read，归票 314 范围内，不新开票）。｜原: 实施中·待验
- **写面**：`cmd/wisp/**`（TestAC14GoSideEvalPushReachesThePage／panel 33r9-nail 窗创建域）＋`internal/panel/**` 如需
- **两形**：①**工作库 `go test ./cmd/wisp/` 全包跑死于 TestAC14GoSideEvalPushReachesThePage 中途**（无 panic/--- FAIL/exit status——静默杀、401-r1 报＋401-v1 双树归因 HEAD 既有零翻转归属）②**导出树 TestAC4FocusReturnToPriorWindowGap33r5 断言红**（33r5 push 维度、摆动族——pre-flip 红/flipped 绿方向利好翻转但在册摆动）

## 判据（派时细化）
- [x] **AC#1 静默杀根因定位＋修**（panel 33r9-nail 窗创建域）
- [x] **AC#2 ./... 整包可跑完**（顶层完成集全、无静默杀）
- [x] **AC#3 33r5 摆动稳定化或显式归因**
- [x] **AC#4 零回归**（六族＋panel＋262/640/902）
- [x] **AC#5 纪律格**
- [2026-11-02 19:0x +08] agent=编排者（结案）did=据非实现者 402-c2 重评「同意结402」翻框结案：AC14推送维度a07f685eb 6/6green;静默杀/整包零AV经正确三层staging现跑581/575 RUN零AV证为cmd/wisp真窗族平台事实+staging假象(非402缺陷);402-v1之FAIL系staging未摆对。待办:真窗族整包零AV需安静窗口复跑(pending_read,票314范围)。判语=docs/evidence/s1/402-c2-verdict.md。
