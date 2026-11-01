# 400 — CI 修复腿 ci-r1 上报的产品级红三件（panel 契约漂移／35r8 流日志无界／winsec 342r1 竞态）

**Status:** **DONE-带残余**（2026-11-01 编排者收口：五格全过（400-v1 三方实测）＋残余①执行链闭环（400-r2 停手→400-r3 Go 侧→TS 键落地机主预批）；残余=342 两格归自己线）
- **①panel 契约漂移**：★**Go Snapshot 发 [instructions tasks] 而 PanelSnapshot 未声明**（C17 桥面契约不一致——真产品发现）＋**panel 35r8 流日志无界**（4x flood→4.0x 驻留）
- **②winsec 342r1 竞态**：★**TestAC1..342r1 祖先换位窗 2/400 命中**（本机闸与 CI 同数＝真产品发现；winsec 封链语义未碰）
- **写面**：`internal/panel/**`＋`internal/winsec/**`（待派时按格细化）

## 判据（派时细化）
- [x] **AC#1 panel 契约对齐**（★PanelSnapshot 声明补齐或 Go 侧收敛——对照 C17）
- [x] **AC#2 35r8 流日志有界**
- [x] **AC#3 winsec 竞态修**（★342r1 封链语义内）
- [x] **AC#4 反向判据×3**（各对应红复现）
- [x] **AC#5 纪律格**
## ★★★收口落地（2026-11-01 22:4x 编排者，凭 400-v1 五格全过＋残余①执行链）

★**机主三连批复（A1378 逐字「追认、TS 预批、R15 全按推荐」）**：①**panel.ts 追认**＝收口条件满足（400-v1 附条件①外部 agent 知悉措辞在 400-v1 裁决表内、由机主转）②**TS 预批**＝droppedCount 键落地（frontend/src/lib/panel.ts TaskRosterSection＋字段＋注释、**「Go 发 TS 未跟」窗口闭**）③R15＝68 解锁（另票）
★**残余①执行链**：400-r2 停手具名（wire 类型面错位）→400-r3 Go 侧落地（23f69b8b9）→**TS 键落地（本笔）＝闭环**
★残余：票 342 AC#1/AC#4 两格归自己终裁线／wire 后续演化按需
