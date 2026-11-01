# 396 — 保温修复（resampler/detector 半隔离＋393 Reset 统一器接入）——392/393 链收口腿

- **Status**：**实施中**（编排者 2026-11-01 11:3x 开；来源＝392-r1 残余「保温修复票收窄为 resampler/detector 半」＋393-v1「Reset 即统一器」＋26p-v2 接力裁定）
- **写面**：`internal/audio/**`（★**392-r1/393-r2 已让出**）
- **根因链（392-r1 §6.1＋393-v1）**：保温判死根因＝resampler/detector 半（毒源）＋APM 跨窗携带（**393 Reset 已可原地清＝统一器**）⇒ **本腿＝毒源半隔离＋Reset 接入保温路径 ⇒ 保温姿势复活读数**

## 判据
- [ ] **AC#1 毒源半隔离**：★**resampler/detector 逐件隔离定位**（392-r1 §6.1 具名残余的执行）
- [ ] **AC#2 Reset 接入保温**：★**窗间 Reset → 第二窗重新会判**（393-r2 的「第二窗重新会判」形状在保温场景复现）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★internal/audio 全套＋**keepalive=false 默认 ⛔ 动——392-r1 裁定、修复经显式开启读数**）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ AEC3 源树／`cmd/wisp`／thresholds