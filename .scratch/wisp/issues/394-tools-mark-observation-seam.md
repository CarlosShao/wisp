# 394 — internal/tools mark 观察缝（喂 TTSGate.RecordTaintInheritance）——R382-1 任务级积累钩的真解（契约面待机主批）

- **Status**：**待机主批**（编排者 2026-11-01 09:3x 开；来源＝382-w-r2 停手②证据链：**生产唯一 mark 点＝internal/tools/bridge.go:687 MarkWithHostPath、tools.Options 无任何 mark 观察钩、ScopeTaints 按安全设计只回元数据** ⇒ admitTask 收尾闭包只能发空内容（假接线）或 Origin 冒充（改闸语义）——两形 ⛔）
- **写面**：`internal/tools/**`＋`cmd/wisp/**`（接线）
- **机理**：★382 形 B 的 RecordTaintInheritance 契约＝「调用方在 mark 时累积归一化内容」（tts_gate.go:245）⇒ **生产缺调用方**；★**真解两形**：①**internal/tools 加 mark 观察缝**（tools.Options 加 OnMark 观察钩喂闸——**扩 tools 契约**）②**闸/调用方下沉 internal/tools**（**更深的契约面**）

## 判据（机主批后生效）
- [ ] **AC#1 观察缝落地**（★**选型①或②、编排者按机主批执行**）
- [ ] **AC#2 端到端**：★**任务污点 → accumulate → 播报拒**（382-w-r2 的机制读数转生产真挂）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★六族＋tools 消费方）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ thresholds/golden