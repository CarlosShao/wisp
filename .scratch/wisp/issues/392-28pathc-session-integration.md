# 392 — 28-PathC 会话集成：消费 26p-r1 钉死的设计输入（DEFERRED 重开②）

- **Status**：**实施中**（编排者 2026-10-31 22:0x 开；来源＝A1325 DEFERRED 重开②＋**59-v3 字面重开条件已满足**：「26 门模块落地〔389 收口〕＋首格验收过〔26p-r1 Feed 等价 6 窗零漂移〕」均在盘）
- **写面**：`cmd/wisp/**`＋`internal/audio/**`（★**26p-r1 已让出**）
- **设计输入（26p-r1 读数钉死、本腿消费）**：
  ①★**AEC 双讲抑制损伤摆动 0-8dB（-29.4/-38.2）是主导项 ⇒ 28 集成的调参活**（引擎默认参数未调）
  ②★**monitor 跨公告保温＝28 的集成深度**（公告窗 -30 工作点是 monitor 构造面选择）
  ③★26p-r1 的接线（voiceBargeinMonitor／HalfDuplexGate.SetConversation／drain goroutine tts-infer）＝**28 会话侧要对齐的形状**

## 判据
- [ ] **AC#1 会话侧集成**：★**跨公告保温/会话窗与 barge-in monitor 的深度对齐**（★**消费设计输入②**；★**对照 28-r1 的 Scope 形状**）
- [ ] **AC#2 AEC 调参读数**：★**双讲损伤摆动收窄尝试**（★**引擎参数面**；★**⛔ AEC3 源树改**）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★**六族＋Path-C 三真数不退化**——★**若 26p-v1 判语与本腿前提冲突 ⇒ 停手重对齐**）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `internal/{risk,session,speech}/**` 语义（★**SetConversation 公开面可用**）／⛔ D43 表／thresholds/golden