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