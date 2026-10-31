# 388 — 测试播报音量衰减：语音测试真设备出声惊扰机主（2026-10-31 机主被吓一跳）——cgo 浮点音量衰减

- **Status**：**待派**（编排者 2026-10-31 14:2x 开；来源＝**机主亲耳听到测试台词「任务已经完成，明天下午三点的会议改到四点」受惊**＋59-r1 运行期干扰报告互相印证）
- **Blocks:** nothing · **Blocked by:** nothing
- **写面**：`cmd/wisp/**` 测试面（announce 测试的音量前置）或 `internal/audio/**` helper
- **机理**：★语音测试**真设备真 matcha**（≤400ms 打断验证必须真 waveOut）⇒ **测试台词从机主真音箱播出**；★台词实锚＝`resident_voice_interrupt_381r2_windows_test.go:458` `const sentence`／`:255` announce
- **★★⚠ 修法边界（59-r1 仪器发现②）**：★**`syscall.SyscallN` 整型路调 `ISimpleAudioVolume::SetMasterVolume` 会把会话音量写成≈0 且被 Windows 持久化** ⇒ ★**衰减必须走 cgo 真浮点调用**（59-r1 rig 已证可行）——**⛔ 整型路**
- **★★⚠ 别破坏被测物**：★≤400ms 打断验证量的是 abort 路径时延；★**衰减 ⛔ 改 waveOut 真路径语义**（⇒ 只动音量、⛔ 改成虚拟端点——那会改变被测对象）

## 判据
- [ ] **AC#1 衰减机制**：★测试 announce 前置 cgo 浮点音量衰减（★**值参数化、机主原音量测后复原**）
- [ ] **AC#2 出声可闻度对照**（★**衰减后播报→录回读数或机主确认**）
- [ ] **AC#3 反向判据**（★摘衰减 ⇒ 出声恢复；三读数＋还原核）
- [ ] **AC#4 零回归**（★≤400ms 打断读数分布不变——★**音量衰减 ⛔ 影响 abort 时延**）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `internal/{risk,session}/**`