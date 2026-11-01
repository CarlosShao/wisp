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
## 选型单更新（2026-11-01 10:4x 编排者，凭 382-w-v2 主新发现——原两形不穷尽）

- 生产带内容 mark 实测共 4 处跨 3 包：bridge.go:687／subagent_197.go:512（同任务作用域 mark——382-w-r2 的「唯一 mark 点」表述不周延）／panel_assets.go:245／projctx.go:349 ⇒ tools.Options.OnMark 只能盖 1 处
- 第三形＝risk 层 mark 观察缝（ProvOptions 钩或 TTSGate 借引擎自注册）：覆盖全部 4 处、契约面只在 risk、并收回调用方根本无法履行的 normalizeTaint 归一化义务（risk 未导出函数——382wr2 测试被迫手选已归一化常量即此痛）；不动 Scope.Close 语义、非 tri-v1 否决的形 A
- 代价如实：留存策略从显式 opt-in 变引擎侧姿势（建议只收任务作用域 mark）
- ⇒ 选型单＝三形（①tools.Options.OnMark ②闸下沉 ③risk 层观察缝）——批仍归机主（契约面）；382-w-v2 建议：批前把第三形并入选型单（已并）
