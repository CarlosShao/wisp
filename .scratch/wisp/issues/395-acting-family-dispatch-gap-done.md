# 395 — D43#15 工具臂与 Acting 家族（#17-20）零生产派发：语音回合带 tool_calls 却恒走纯文本臂（16-v2 升级发现）

- **Status**：**已结案 done**（2026-11-02 14:0x 编排者据非实现者 `395-v1` 五格＋SPEC-10 七项全 PASS「同意结 395」结案；落地腿 `395-r1`=commit `2216f0d3d`（D43#15 工具臂 Thinking→Acting＋#18 Acting→Speaking 派发真接）；判语件 `docs/evidence/s1/395-v1-verdict.md`=`059011810`）｜原: 待派 2026-11-01
- **写面**：`cmd/wisp/**`（resident_voice_engines_windows.go 的派发位）
- **机理（16-v2 逐字实锤）**：★demo 回合 SSE 首响应带 tool_calls（tool_call 行落真库）、**产线 resident_voice_engines_windows.go:979-981 恒以 HasToolCall:false 走 #15 纯文本臂**；★**D43#15 工具臂与 Acting 家族（#17-20：EvApprovalNeeded/EvToolsCompleted/EvRepeatThreshold/EvPathConflict）全仓零生产派发位**（只存在于冻结表＋表测试＋ball 交互测试）⇒ **语音回合的工具结果永远不走 Acting 呈现**
- ★**demo 走位仍每行合法**（无非法转移）＝缺口非破坏

## 判据
- [x] **AC#1 工具臂派发真接**：★**HasToolCall:true ⇒ 派 #15 工具臂（Acting 家族可达）**（★**对照 D43 表内行、⛔ 表一字动**）
- [x] **AC#2 端到端**：★**语音回合带 tool_calls ⇒ Acting 呈现读数**（★**对照 16-r1 demo 形状**）
- [x] **AC#3 反向判据**（★摘派发 ⇒ 回纯文本臂 ⇒ 红；三读数＋还原核）
- [x] **AC#4 零回归**（★六族＋Path-C 套件＋262/640/902）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ D43 表／statemachine／internal/{risk,session,speech} 语义

- [2026-11-02 14:0x +08] agent=编排者（结案）did=据非实现者 395-v1 五格＋SPEC-10 七项全 PASS「同意结 395」翻框结案：落地腿 395-r1=commit 2216f0d3d（D43#15 工具臂 Thinking→Acting＋#18 Acting→Speaking 派发真接）；判语=docs/evidence/s1/395-v1-verdict.md(059011810)。
