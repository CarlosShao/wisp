# 386 — mockllm 同类帽登记不修三处：reasoning/reason 200 字节帽×3 方言＋max_tokens 字节切（385-r1 登记残余）

- **Status**：**待派**（编排者 2026-10-31 12:5x 开；来源＝票 385 落地腿登记不修三项）
- **Blocks:** nothing · **Blocked by:** nothing（★`tools/mockllm` 已被 385-r1 让出）
- **写面**：`tools/mockllm/**`
- **机理（385-r1 逐行点名）**：★`chat.go:168-169`／`responses.go:100-101`／`messages.go:114-115` **reasoning/reason 200 字节帽三方言**（同类多字节切风险、票面机理未点名）；★`chat.go:207` **max_tokens 字节切**（`answer[:MaxTokens*4]`、带 `finish=length` 模拟语义——★**这个帽可能是有意的（模拟真实 max_tokens 行为）、修前先判**）

## 判据
- [ ] **AC#1 逐帽判定**：★**每帽"有意模拟 vs 无意缺陷"具名**（★**对照真实 chat API 的 max_tokens/reasoning 行为**；★**有意模拟的帽 ⛔ 动、如实登记**）
- [ ] **AC#2 多字节回归**：★**判"无意"的帽按 rune/去帽修＋多字节案一发**（★**对照 385-r1 的 402 汉字案形**）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★mockllm 全包＋消费方 E2E）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`