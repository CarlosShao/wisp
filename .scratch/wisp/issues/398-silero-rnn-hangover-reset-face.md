# 398 — silero RNN hangover：internal/speech 无 session reset 面（392-v3 采纳路径②、非阻塞）

- **Status**：**待派**（编排者 2026-11-01 13:0x 开；来源＝392-v3 采纳路径②「RNN 无 reset 面、冷段是 owner 纪律非 session 清面、非阻塞」）
- **写面**：`internal/speech/**`（★387-r3 已让出）
- **机理（392-v3 三进程同读）**：★**silero RNN 语境跨空窗携带**（无输入不衰减）、**hangover 能在 post-speech 静音上开火**（探针 silence=1 talker=0）；**冷段吸收之**——但「更硬的清面」＝internal/speech 需 session reset 面（**新 API、契约相邻**）

## 判据（派时细化为逐格）
- [ ] **AC#1 reset 面设计**（★对照 detector Finalize 形状；★**⛔ 动 silero 语义**）
- [ ] **AC#2 hangover 清除读数**
- [ ] **AC#3 反向判据**
- [ ] **AC#4 零回归**（★internal/speech 全套 61/0/6＋裸跑 25/0/42）
- [ ] **AC#5 纪律格**