# 391 — mockllm max_tokens 死字段：messages/responses 两方言解析后合成路径从不应用（386-r2 登记不修）

- **Status**：**待派**（编排者 2026-10-31 22:1x 开；来源＝386-r2 登记不修①）
- **写面**：`tools/mockllm/**`（★**386-r2 已让出**）
- **机理（386-r2 逐行）**：★**messages/responses 两方言解析 max_tokens/max_output_tokens 后、合成路径从不应用**（responses.go:67、messages.go:179 死字段）——**chat 有截断模拟（386-r2 修为 rune 边界）而这两方言没有** ⇒ max_tokens 语义与真 provider 不齐（三方言不一致）

## 判据
- [ ] **AC#1 两方言补齐 max_tokens 应用**（★**对照 chat 方言 386-r2 的 rune 边界形状——同形**）
- [ ] **AC#2 多字节回归**：★**两方言各一发 max_tokens 汉字切案**（finish=length＋rune 边界前缀——对照 386-r2 案形）
- [ ] **AC#3 反向判据**（★三读数＋还原核）
- [ ] **AC#4 零回归**（★mockllm 全包＋消费方 E2E）
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`