# 391 — mockllm max_tokens 死字段：messages/responses 两方言解析后合成路径从不应用（386-r2 登记不修）

- **Status**：**已结案 done**（2026-11-02 18:1x 编排者据非实现者 `391-v1` 五格全 PASS「同意结 391」结案；落地件＝commit `6cf941b72`（票391落地腿391-r1：messages/responses 两方言补齐 max_tokens 应用，对照 chat 386-r2 的 rune 边界同形；两方言各一发汉字切案；`tools/mockllm` FAIL=0、消费方 E2E 全绿）；判语件 `docs/evidence/s1/391-v1-verdict.md`＝dcb2a28a3）｜原: 待派 2026-10-31
- **写面**：`tools/mockllm/**`（★**386-r2 已让出**）
- **机理（386-r2 逐行）**：★**messages/responses 两方言解析 max_tokens/max_output_tokens 后、合成路径从不应用**（responses.go:67、messages.go:179 死字段）——**chat 有截断模拟（386-r2 修为 rune 边界）而这两方言没有** ⇒ max_tokens 语义与真 provider 不齐（三方言不一致）

## 判据
- [x] **AC#1 两方言补齐 max_tokens 应用**（★**对照 chat 方言 386-r2 的 rune 边界形状——同形**）
- [x] **AC#2 多字节回归**：★**两方言各一发 max_tokens 汉字切案**（finish=length＋rune 边界前缀——对照 386-r2 案形）
- [x] **AC#3 反向判据**（★三读数＋还原核）
- [x] **AC#4 零回归**（★mockllm 全包＋消费方 E2E）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`
- [2026-11-02 18:1x +08] agent=编排者（结案）did=据非实现者 391-v1 五格全 PASS「同意结391」翻框结案；落地件=commit 6cf941b72（两方言补齐max_tokens应用,对照chat 386-r2 rune边界同形）;判语=docs/evidence/s1/391-v1-verdict.md。
