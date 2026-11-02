# 386 — mockllm 同类帽登记不修三处：reasoning/reason 200 字节帽×3 方言＋max_tokens 字节切（385-r1 登记残余）

- **Status**：**已结案 done**（2026-11-02 17:1x 编排者据非实现者 `386-v1` 五格全 PASS「同意结 386」结案；落地件＝commit `1773881db`（票386落地腿386-r2：reasoning/reason 200 字节帽×3 方言＋max_tokens 字节切逐帽判定；有意帽 ASCII 对照钉、无意帽 rune 修；`tools/mockllm` FAIL=0、消费方 E2E 全绿）；判语件 `docs/evidence/s1/386-v1-verdict.md`＝200fb51b6）｜原: 待派 2026-10-31
- **Blocks:** nothing · **Blocked by:** nothing（★`tools/mockllm` 已被 385-r1 让出）
- **写面**：`tools/mockllm/**`
- **机理（385-r1 逐行点名）**：★`chat.go:168-169`／`responses.go:100-101`／`messages.go:114-115` **reasoning/reason 200 字节帽三方言**（同类多字节切风险、票面机理未点名）；★`chat.go:207` **max_tokens 字节切**（`answer[:MaxTokens*4]`、带 `finish=length` 模拟语义——★**这个帽可能是有意的（模拟真实 max_tokens 行为）、修前先判**）

## 判据
- [x] **AC#1 逐帽判定**：★**每帽"有意模拟 vs 无意缺陷"具名**（★**对照真实 chat API 的 max_tokens/reasoning 行为**；★**有意模拟的帽 ⛔ 动、如实登记**）
- [x] **AC#2 多字节回归**：★**判"无意"的帽按 rune/去帽修＋多字节案一发**（★**对照 385-r1 的 402 汉字案形**）
- [x] **AC#3 反向判据**（★三读数＋还原核）
- [x] **AC#4 零回归**（★mockllm 全包＋消费方 E2E）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`
- [2026-11-02 17:1x +08] agent=编排者（结案）did=据非实现者 386-v1 五格全 PASS「同意结386」翻框结案；落地件=commit 1773881db（逐帽判定:有意帽ASCII对照钉/无意帽rune修;tools/mockllm FAIL=0）;判语=docs/evidence/s1/386-v1-verdict.md。
