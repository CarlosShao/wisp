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

## Progress log (append-only, newest last)

- [2026-11-10 22:3x +08] wave121（`internal/audio/aecspike`，7 件，台账 A1688）＝AC#1 的机制面先行交付：cgo 浮点守卫 `ApplyAnnounceAttenuation`＋`Restore`＋真机读数 original=1.000000→attenuated=0.150000（回读一致）＋loopback 对照 750.3→46.5 peakrms（AC#2 的录回读数）。未做：cmd/wisp 接线（当波 423 占槽）。

- [2026-11-11 14:3x +08] agent=388-closeout-N（写腿 cmd/wisp/** 测试面·只 commit 未 push）did=announce 测试接线＋≤400ms 分布断言 verdict=**两枚真设备播报测试已带衰减前置，重测全绿，分布读数在册**：①新件 `cmd/wisp/voice_announce_attenuation_388_windows_test.go`：`attenuateAnnounce388(t, attenuated)`（值参数化＝AC#1；调用点传 `aecspike.DefaultAnnounceAttenuatedVolume` 0.15，数字不抄进本包）——播报前 `ApplyAnnounceAttenuation`、`t.Cleanup` 里 `Restore()` 并以 `ATTENUATE-388 original/attenuated/restored` 三读数落日志（AC#3 的还原核＝回读核；守卫自身契约：一 Apply 恰一 Restore）；守卫在本机不可用时 **skip 不响跑**（满音量跑＝388 要拦的惊扰本身，假绿此路不通）；＋`assertInterruptDistribution388`（n/min/median/max 分布日志＋每样本同 400ms 界断言，不松任何既有界）；＋`Test388AnnounceAttenuationWired` 接线钉（runtime.Caller 落盘读两枚测试件，逐枚验「attenuateAnnounce388 在第一发 rv.announce 之前」＋helper 含 Apply/t.Cleanup/Restore 配对；253r1 名册钉同形）。②接线面＝两枚真设备 ≤400ms 测试各一行前置：`Test381r2RealDevicePlayingPhaseAbortWithin400ms`（Esc 路，票面台词锚 381r2:458）＋`Test390r1BallClickRealDeviceWithin400ms`（球点路，同句）——控制轮与全部中断轮同罩；**被测物语义零触碰**（waveOut 真路径、abort 计时、逐轮 ≤400ms Fatalf 全部原样；衰减只走会话音量 cgo 浮点路＝59-r1 发现②的边界）。③真机重测读数（staged 二进制，票 409 载体形，System32 onnxruntime 遮蔽须 DLL 贴测试二进制）：两测试各 **PASS**；381r2＝ATTENUATE original=1.000000→attenuated=0.150000→restored=1.000000，unstick 分布 n=5 min=29ms median=39ms max=49ms、waveOutReset 同 29/39/49ms（界 400ms）；390r1＝同三读数 33/39/49ms——**衰减后 abort 时延仍在 29-49ms 同量级，AC#4 的"衰减 ⛔ 影响 abort 时延"读数成立**；AC#3 反向＝restored=1.000000 回读核。④门禁：gofumpt -l cmd/wisp/ 空／`go vet ./cmd/wisp` rc0／cmd/wisp 整包 staged 载体干净一发 FAIL=3＝HEAD 存量三枚（198/313/33r11）逐名同色；`Test388AnnounceAttenuationWired` PASS；d22scan clean／path-budget GREEN。⑤没做：AC#2 的"机主亲耳确认"那半（loopback 录回读数 wave121 已交；听感确认需机主在场，如实留格给收件）；`internal/{risk,session}/**`、thresholds/golden/PLAN/specs 零字节。