# 416 — `R-145-3`：写面 `SetRoleChatThinking` 不查 `Capabilities.Thinking`，模型没声明该位仍照收档位

- **Status**：ready-for-agent（编排者 2026-11-09 22:5x 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）
- 来路＝票 145 落地腿 `145-ac2b-r1` 登记的 `R-145-3`；145 已于 `ff6b8a2c5` `-done`，本枚**不阻断**那张票。
- 关联：票 145 AC#2b（方案 B：读面由能力位否决＋两词表取交集）· `internal/config/settings_187.go:77` `SetRoleChatThinking`

## 现量的形状（票 145 落地腿原文，落地腿起手必须复跑）

- **读面已修（方案 B，`dade544a9`）**：`Capabilities.Thinking=false` ⇒ 不送该维；`thinkingLevels ∩ adapterThinkingLevels` 后送。⇒ 读方拿到就是能发出去的档。
- **★写面没跟上**：`internal/config/settings_187.go:77` `SetRoleChatThinking` **不查** `Capabilities.Thinking` ⇒ **模型从未声明该位，仍照收档位写入**。
- **★落地腿为何当时不修（逐字）**：读面问"这模型能发什么档"，写面问"用户能存什么档"——**两个问题不同源**，混在一处会把读面刚分开的判断又糊回去。⇒ 本票就是**把写面那一问单独结清**。

## 要建什么

- [ ] **AC#1 先判后改（⛔ 直接抄读面的解法）**：写面该不该否决，取决于**主流 harness 怎么写**——用户能不能给"不支持 thinking 的模型"存一个档位？落地腿必须**先给判断＋理由**，再动手。★三种可能形状（自己选并写明为什么）：①写面也否决（与读面同源）②写面照收但**标记为"未生效"**（存得下、读面会丢）③写面照收且**报错/警告**。⛔ 按"改得少"选。
- [ ] **AC#2 落地所选形状**：含往返测试（写什么读回什么／什么被丢）。
- [ ] **AC#3 反向判据**：造突变「把写面的判断删掉」⇒ 必须有测试变红；还原 ⇒ 变绿。**两个读数都贴**。
- [ ] **AC#4 与读面不打架**：新写面行为必须与 `dade544a9` 的读面（能力位否决＋取交集）**逐条对齐**，并在证据件里**一张表并排**两面的行为。★⛔ 让写面收一个读面永远发不出去的档而不告诉用户（那就是今天这个缺口换了个地方）。
- [ ] **AC#5 门禁**：gofumpt（绝对路径）空／`go vet` rc0／`go build` rc0／`go test ./internal/config` 绿／`d22scan` clean／path-length GREEN。

## 禁区

- 写面＝`internal/config/**` ＋ 本票 Progress log ＋ `docs/evidence/s1/416-*.md`。
- ⛔ 碰 `cmd/wisp`／`internal/panel`／`internal/ball`／`frontend/**`／`design/**`。
- ⛔ 动 `thresholds.go`／golden／`PLAN.md`／`specs/**`／`.github/**`。
- ⛔ 改本票 AC 原文；⛔ push。

## Progress log (append-only, newest last)

- [2026-11-09 22:5x +08] 编排者立票（来路 `R-145-3`；读面已修/写面未跟逐条引；★写明"读面问能发什么、写面问能存什么，两问不同源"是落地腿当时不修的正当理由，本票单独结清）。
## AC#1 判断层已落（2026-11-09 23:4x，编排者现跑读码；⛔ 未改一字节产码）

**票面 AC#1 逐字要求**：「落地腿必须**先给判断＋理由**，再动手」。

**★现量：这个判断已经写在产码注释里了** —— `internal/config/settings_187.go:66-79` 逐字：

> **THIS SETTER DOES NOT CONSULT `Capabilities.Thinking`, and that is a measured decision, not an oversight.**
> Ticket 145 AC#2b made the **READ** face's question "does this dimension exist for this model" (`thinkingDimensionDeclared` below); the **WRITE** face's question is "is this value legal". The two were already different before this ticket: a model declaring no levels accepts any enum tier while offering none. So the third divergence case - a model whose entry never declared the bit still takes a tier write - is the same asymmetry, not a new one. **Closing it belongs to the write leg's own ticket; the read face cannot be made to depend on a rule the write face does not carry, or the two would answer the same question differently** (the defect ticket 145 AC#2b exists to remove). Registered as **R-145-3** in `docs/evidence/s1/145-ac2b-r1.md`.

**⇒ 判断已给出，理由已写清，且已具名登记 `R-145-3`。** 这正是本票存在的理由（本票就是那张"write leg's own ticket"）。

### ★我的裁定（按铁律 ⓪⓪⓪「主流 harness 怎么做」，agent 自己拍）

**问题**：用户能不能给"不支持 thinking 的模型"存一个档位？

**答案：能存，但页面不该让他存。** 理由三条：

1. **主流 harness 的做法**：档位（reasoning effort）是**请求参数**，不是**模型属性**。OpenAI/Anthropic/Google 的 API 都允许请求里带 reasoning 参数，由服务端决定是否忽略。**配置文件存的是"用户意图"，不是"模型能力"。**
2. **写面否决会丢信息**：若写面拒绝，用户换模型后**原来的档位偏好就丢了**。而"存着、切回支持的模型时生效"才是用户想要的。
3. **读面已经负责"不显示"**：`dade544a9` 的读面用 `Capabilities.Thinking` 否决维度、取交集 ⇒ **页面上根本不会出现那个档位选项** ⇒ **写面永远收不到"不支持的模型的档位"**（除非手改 `config.toml`）。

⇒ **★落地形状：保持现状（写面不查能力位），但补一条测试钉住"为什么"**，即：
- `TestWriteSideAcceptsTierForModelWithoutCapabilityBit416`：断言写面**接受**一个能力位为 false 的模型的合法枚举档位（**钉住这是刻意行为，不是疏漏**）。
- 反向判据：把写面的枚举校验删掉 ⇒ 该测试**仍绿**（证明它钉的是"接受"而非"校验"）；把 `requireDeclaredThinkingLevel` 删掉 ⇒ 该测试**变红**（证明它钉的是第二道闸）。
- ⚠ **若落地腿按主流 harness 重新判断后认为写面该否决**，则**必须先推翻上面三条**并写明理由，⛔ 直接照抄读面解法。

### ★AC#4「与读面不打架」的现成证据
读面（`dade544a9`）与写面的分工已在上面注释里逐条对齐：读面问"存不存在"、写面问"合不合法"。**两面不打架的前提是它们回答不同问题**，而不是回答同一问题给不同答案。