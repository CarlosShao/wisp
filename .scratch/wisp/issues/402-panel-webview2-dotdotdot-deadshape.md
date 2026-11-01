# 402 — panel/WebView2 域：./... 整包扫 cmd/wisp 死形修复（401-v1 建议维持）

- **Status**：**待派**（编排者 2026-11-01 23:3x 开；来源＝401-r1 发现＋401-v1 双树归因闭合——**两形并记：工作库静默杀／导出树 33r5 断言红**）
- **写面**：`cmd/wisp/**`（TestAC14GoSideEvalPushReachesThePage／panel 33r9-nail 窗创建域）＋`internal/panel/**` 如需
- **两形**：①**工作库 `go test ./cmd/wisp/` 全包跑死于 TestAC14GoSideEvalPushReachesThePage 中途**（无 panic/--- FAIL/exit status——静默杀、401-r1 报＋401-v1 双树归因 HEAD 既有零翻转归属）②**导出树 TestAC4FocusReturnToPriorWindowGap33r5 断言红**（33r5 push 维度、摆动族——pre-flip 红/flipped 绿方向利好翻转但在册摆动）

## 判据（派时细化）
- [ ] **AC#1 静默杀根因定位＋修**（panel 33r9-nail 窗创建域）
- [ ] **AC#2 ./... 整包可跑完**（顶层完成集全、无静默杀）
- [ ] **AC#3 33r5 摆动稳定化或显式归因**
- [ ] **AC#4 零回归**（六族＋panel＋262/640/902）
- [ ] **AC#5 纪律格**