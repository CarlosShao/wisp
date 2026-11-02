# 385 — mockllm 回声 400 字节字节级硬截断：多字节文本拦腰切（归 tools/mockllm 所有者）

- **Status**：**已结案 done**（2026-11-02 16:2x 编排者据非实现者 `385-v1` 五格全 PASS「是否同意结 385: 是」结案；落地件＝commit `fad2f5e4f`（2026-10-31 原实现：三处 400-byte hard clamp 全摘、`runeBoundaryWithin` 回退 UTF-8 边界、streaming 按 `[]rune` 分块）；判语件 `docs/evidence/s1/385-v1-verdict.md`）｜原: 待派 2026-10-31
- **Blocks:** nothing · **Blocked by:** nothing
- **写面**：`tools/mockllm/**`
- **机理（`30-v1` 逐字核）**：★`tools/mockllm/chat.go:159-162` 逐字核到 **`len(answer) > 400` 字节截断** ⇒ ★**402 汉字（多字节）物理回不来、拦腰切** ⇒ 票 30 `AC#2` 的 402 汉字案被迫改走 C5 进程内缝假体
- ★**影响**：★**golden/契约测试凡走 mockllm 的长响应全部被截**（⛔ 影响判定、只影响测试基建）

## 判据
- [x] **AC#1 按 rune 截断或去掉截断**（★**对照该截断的原始意图——防超长回声；★修法具名**）
- [x] **AC#2 多字节回归一发**：★**402 汉字案走 mockllm 不再拦腰切**（★**与 `30-r1` 的 C5 缝案同输出口径对照**）
- [x] **AC#3 反向判据**：★**摘修复 ⇒ 402 汉字案必红**（三读数＋还原核）
- [x] **AC#4 零回归**：★**mockllm 现有消费方（票 12/30 的 E2E）全绿**
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ `cmd/wisp`／`internal/**`（★**只 `tools/mockllm/**`**）
- [2026-11-02 16:2x +08] agent=编排者（结案）did=据非实现者 385-v1 五格全 PASS「同意结385」翻框结案；落地件=commit fad2f5e4f（三处400byte clamp全摘+runeBoundaryWithin UTF-8边界+streaming []rune分块）；判语=docs/evidence/s1/385-v1-verdict.md。
