# 393 — AEC3 绑定层 flush/reset API（392 残余：保温真解的前置）

- **Status**：**待派**（编排者 2026-11-01 00:1x 开；来源＝392-r1 残余「AEC3 无公开 flush/reset 面——保温真解属绑定层新 API」）
- **写面**：`internal/audio/aec/**`（★**绑定层＋shim；⛔ AEC3 源树**）
- **机理**：★**保温姿势判死的根因链**＝APM 跨窗携带失准＋无 flush/reset ⇒ **绑定层新增 flush/reset 面**（shim wap_flush/wap_reset＋Go 侧 Flush/Reset 方法）⇒ **后续保温修复票的前置**

## 判据
- [ ] **AC#1 绑定层 API 落地**：★**shim＋Go 方法（flush/reset）**（★**对照 wap_new_tuned 形状**）
- [ ] **AC#2 单测**：★**flush/reset 后状态清洁读数**（★**对照 392-r1 的 K=4 测量形状**）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★internal/audio 全套 94/83/0/11 基线＋392 的调参面不退化）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ AEC3 源树／`cmd/wisp`（26p 集成已收）／thresholds