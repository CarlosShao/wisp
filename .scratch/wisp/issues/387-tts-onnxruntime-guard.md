# 387 — `internal/speech/tts.go` NewMatchaTts 缺 onnxruntime 守卫：System32 1.17.1 遮蔽 ⇒ 裸跑 go test 崩 0xc0000005（59-r1 仪器发现①）

- **Status**：**实施中**（编排者 2026-10-31 13:3x 开；来源＝`59-r1` 仪器发现①；★**`NewDetector`/`NewStreamingAsr` 都有 `checkOnnxruntime` 守卫、独 `NewMatchaTts` 没有**）
- **Blocks:** nothing · **Blocked by:** nothing（★`internal/speech` 空闲）
- **写面**：`internal/speech/**`（★**只加守卫、⛔ 改任何引擎语义**）
- **机理（59-r1 实测）**：★**System32 驻 onnxruntime 1.17.1 遮蔽 PATH 的 1.28.2** ⇒ ★**裸 `go test ./internal/speech` 时 `NewMatchaTts` 崩 `0xc0000005`**——★**守卫家族里漏了它**（对照 127/131 先例：同族守卫＋具名 skip）

## 判据
- [ ] **AC#1 守卫补齐**：★`NewMatchaTts` 挂同族 `checkOnnxruntime`（★**与 NewDetector/NewStreamingAsr 同形**）
- [ ] **AC#2 裸跑不崩**：★**无 staging 直跑 `go test ./internal/speech` ⇒ 具名 skip、零崩溃**（★**对照守卫家族两先例**）
- [ ] **AC#3 反向判据**：★**摘守卫 ⇒ 裸跑崩 0xc0000005 复现**（★**仓外快照真改文件——A380 WalkDir 同款方法**；三读数＋还原核）
- [ ] **AC#4 零回归**：★**staging（三 DLL 同目录＋`WISP_*`）全绿 61/0/6 逐位同**（★**守卫 ⛔ 改 staging 路径下的行为**）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/{risk,session,config,memory,panel}/**`／thresholds/golden