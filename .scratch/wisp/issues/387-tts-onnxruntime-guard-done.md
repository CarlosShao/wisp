# 387 — `internal/speech/tts.go` NewMatchaTts 缺 onnxruntime 守卫：System32 1.17.1 遮蔽 ⇒ 裸跑 go test 崩 0xc0000005（59-r1 仪器发现①）

**Status:** **DONE-带残余**（2026-10-31 编排者收口：五格全部站得住零翻框（387-v1 七项自跑）；残余三条见末节）
- **Blocks:** nothing · **Blocked by:** nothing（★`internal/speech` 空闲）
- **写面**：`internal/speech/**`（★**只加守卫、⛔ 改任何引擎语义**）
- **机理（59-r1 实测）**：★**System32 驻 onnxruntime 1.17.1 遮蔽 PATH 的 1.28.2** ⇒ ★**裸 `go test ./internal/speech` 时 `NewMatchaTts` 崩 `0xc0000005`**——★**守卫家族里漏了它**（对照 127/131 先例：同族守卫＋具名 skip）

## 判据
- [x] **AC#1 守卫补齐**：★`NewMatchaTts` 挂同族 `checkOnnxruntime`（★**与 NewDetector/NewStreamingAsr 同形**）
- [x] **AC#2 裸跑不崩**：★**无 staging 直跑 `go test ./internal/speech` ⇒ 具名 skip、零崩溃**（★**对照守卫家族两先例**）
- [x] **AC#3 反向判据**：★**摘守卫 ⇒ 裸跑崩 0xc0000005 复现**（★**仓外快照真改文件——A380 WalkDir 同款方法**；三读数＋还原核）
- [x] **AC#4 零回归**：★**staging（三 DLL 同目录＋`WISP_*`）全绿 61/0/6 逐位同**（★**守卫 ⛔ 改 staging 路径下的行为**）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/{risk,session,config,memory,panel}/**`／thresholds/golden
## ★★★终裁落地（2026-10-31 22:2x 编排者，凭非实现者 `387-v1`，件 `docs/evidence/s1/387-v1-verdict.md` 61 行、`76600cdd4`）

★**五格全部站得住（零翻框）⇒ 可带残余结案**。
★**387-v1 七项自跑**：①守卫同位同形核＝**四先例逐一对核＋全包 sherpa 直调点 11 处自数、生产面 8 处全在守卫后＝100% 覆盖成立** ②裸跑复跑 **25/6/36 零 Exception、8 处 skip 句逐字含守卫具名错误** ③反向**真改文件**摘守卫 ⇒ 0xc0000005 栈逐字（**mutant SHA 与 387-r2 记录逐字相同＝同一编辑独立复现、反证其真用 A380 真改法**）④staging 改前后各 61/0/6＋名册 diff=0＋matcha 真合成在跑（⚠ 复跑陷阱：两 WISP_* env 必须指二进制所出自的树）⑤**归因校正核过＝两处都真**（全包首发崩点 TestKwsProbeSherpaDebug 直调＋NewMatchaTts 定向 tts.go:158 算术闭合）⑥d22scan clean＋正控（★**正控须用生产面 .go——禁令豁免 _test.go、第一发放 test 文件未被咬**）⑦五格判。
★**残余三条**：①**6 KWS Fatalf＝测试断言形状、387-r3 在飞修** ②台账 :28956 修订＝A1335 append-only 已办（387-v1 现读核实）③轻微备查＝387r2-impl 的 07-d22scan.log 实为 go run 失败残留、真读数来自预编译 exe（判语不受影响）。
