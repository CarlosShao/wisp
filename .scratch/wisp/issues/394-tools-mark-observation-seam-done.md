# 394 — internal/tools mark 观察缝（喂 TTSGate.RecordTaintInheritance）——R382-1 任务级积累钩的真解（契约面待机主批）

**Status:** **DONE-带残余**（2026-11-01 编排者收口：五格全勾（394-v2 七项自跑）；R382-1 缺口闭合；残余三条见末节）
- **写面**：`internal/tools/**`＋`cmd/wisp/**`（接线）
- **机理**：★382 形 B 的 RecordTaintInheritance 契约＝「调用方在 mark 时累积归一化内容」（tts_gate.go:245）⇒ **生产缺调用方**；★**真解两形**：①**internal/tools 加 mark 观察缝**（tools.Options 加 OnMark 观察钩喂闸——**扩 tools 契约**）②**闸/调用方下沉 internal/tools**（**更深的契约面**）

## 判据（机主批后生效）
- [x] **AC#1 观察缝落地**（★**选型①或②、编排者按机主批执行**）
- [x] **AC#2 端到端**：★**任务污点 → accumulate → 播报拒**（382-w-r2 的机制读数转生产真挂）
- [x] **AC#3 反向判据**（★三读数＋还原核）
- [x] **AC#4 零回归**（★六族＋tools 消费方）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ thresholds/golden
## 选型单更新（2026-11-01 10:4x 编排者，凭 382-w-v2 主新发现——原两形不穷尽）

- 生产带内容 mark 实测共 4 处跨 3 包：bridge.go:687／subagent_197.go:512（同任务作用域 mark——382-w-r2 的「唯一 mark 点」表述不周延）／panel_assets.go:245／projctx.go:349 ⇒ tools.Options.OnMark 只能盖 1 处
- 第三形＝risk 层 mark 观察缝（ProvOptions 钩或 TTSGate 借引擎自注册）：覆盖全部 4 处、契约面只在 risk、并收回调用方根本无法履行的 normalizeTaint 归一化义务（risk 未导出函数——382wr2 测试被迫手选已归一化常量即此痛）；不动 Scope.Close 语义、非 tri-v1 否决的形 A
- 代价如实：留存策略从显式 opt-in 变引擎侧姿势（建议只收任务作用域 mark）
- ⇒ 选型单＝三形（①tools.Options.OnMark ②闸下沉 ③risk 层观察缝）——批仍归机主（契约面）；382-w-v2 建议：批前把第三形并入选型单（已并）

## ★★★终裁落地（2026-11-01 18:1x 编排者，凭非实现者 `394-v2`，件 `docs/evidence/s1/394-v2-verdict.md` 64 行、`bc1417eb3`）

★**五格全勾 ⇒ 可收口**（七项自跑、纯导出树、三枚二进制自建）。
★**394-v2 七项**：①亚形 B 三条理由码上成立（ProvOptions 纯构造期定格／闸只在 attachProvenance 才存在／**全 diff 零新导出符号＝契约面只在 risk 逐字**／internal/tools 0 字节）②**汇流口主张坐实＝本腿关键核点**（**Mark＝MarkWithHostPath(...,"") 薄包装、四点全流经**——382-w-r2 的"唯一 mark 点"表述由此周延化）③留存姿势追认复核成立（MaxAncestorDepth 首-16 帽为真前提、RecordTaintInheritance 逐字节未动、AncestorInheritance 1→2＝已批姿势可见面）④生产真挂逐位复现（**全角 mark vs 半角播报＋白盒 normalizeTaint 断言＝双钉**）⑤反向两发全红有牙＋还原绿 ⑥零回归（risk rc=0＋六族两遍 33/2/0 红集全归因无 394 归因红＋**四分母 288/52/734/184 逐位一致、ls-tree 算术归因闭环**）⑦五格全勾
★**残余三条**：①注册前 mark 不入链（追溯需内容读出口、ScopeTaints 按安全设计只回元数据、不动）②缝日志走 logf（nil=drop）③MaxAncestorDepth(16) 帽＝382 既有防御参数（>16 作用域演化归 382 另开票）
★**394-v2 两点提示收档**：①staged 腿 runner 配方补 **manifest.json(+.minisig)＋third_party/model-fetch junction**（缺则 41a2r1 两红假报回归）②**394-v1 之死不在件上**（27bda97e0 本体独立复跑全部立住、直接进推送队列）
★**R382-1「生产 RecordTaintInheritance 调用方」缺口闭合**（382-w-r2 机制面＋394-r1 生产调用方＋394-v2 终裁＝三段闭环）
