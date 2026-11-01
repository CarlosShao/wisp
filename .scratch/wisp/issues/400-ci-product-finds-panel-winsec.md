# 400 — CI 修复腿 ci-r1 上报的产品级红三件（panel 契约漂移／35r8 流日志无界／winsec 342r1 竞态）

- **Status**：**待派**（编排者 2026-11-01 15:4x 开；来源＝ci-r1 上报不修清单——全为既有红/产品发现）
- **①panel 契约漂移**：★**Go Snapshot 发 [instructions tasks] 而 PanelSnapshot 未声明**（C17 桥面契约不一致——真产品发现）＋**panel 35r8 流日志无界**（4x flood→4.0x 驻留）
- **②winsec 342r1 竞态**：★**TestAC1..342r1 祖先换位窗 2/400 命中**（本机闸与 CI 同数＝真产品发现；winsec 封链语义未碰）
- **写面**：`internal/panel/**`＋`internal/winsec/**`（待派时按格细化）

## 判据（派时细化）
- [ ] **AC#1 panel 契约对齐**（★PanelSnapshot 声明补齐或 Go 侧收敛——对照 C17）
- [ ] **AC#2 35r8 流日志有界**
- [ ] **AC#3 winsec 竞态修**（★342r1 封链语义内）
- [ ] **AC#4 反向判据×3**（各对应红复现）
- [ ] **AC#5 纪律格**