# 390 — 球单击打断装配（381 残余①：381-r1 自认项、381-v3 判"建议独立票"）

- **Status**：**待派**（编排者 2026-10-31 19:0x 开；来源＝票 381 收口残余①）
- **写面**：`cmd/wisp/**`（球单击事件 → `rv.TriggerInterrupt()` 的挂接；★**Esc/veto 通路已落、只缺球单击这一跳**）
- **机理**：★381-r1 的 Esc→`vetoByEsc()`→`TriggerInterrupt` 通路已在（381-r2 增强＝speaker.Abort 真停 ≤400ms 实测 MET）；★**球单击＝第二个触发源**（381-r1 申报"球单击→预留钩子但未装配执行器（需 internal/ball API 扩展，不在本片范围）"）

## 判据
- [ ] **AC#1 球单击→TriggerInterrupt 真挂**（★**对照 Esc 通路的挂接形状**；★**internal/ball API 若需扩展 ⇒ 具名交回**——ball 是冻结面）
- [ ] **AC#2 端到端**：★**播报中球单击 ⇒ TTS 真停 ≤400ms ＋ 门重开**（★**复用 381-r2 的 staged 测量形状**）
- [ ] **AC#3 反向判据**：★**摘球单击挂接 ⇒ 球单击不再打断 ⇒ 红**（三读数＋还原核）
- [ ] **AC#4 零回归**（★**六族＋Esc 通路不动**）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `internal/{ball,risk,session,speech}/**`（★**ball API 扩展须具名交回**）