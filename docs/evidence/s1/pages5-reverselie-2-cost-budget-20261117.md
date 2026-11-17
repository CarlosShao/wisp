# 反向谎言 #2 修复 · `cost-settings.tsx` 预算条文案改真话（不造数字）

- **修复腿**：`agent-105`（写腿，frontend 写面；实现者，⛔ 非裁决者）
- **时间**：2026-11-17 +08
- **来路**：`docs/evidence/s1/pages5-adversarial-20261117-0953.md` §7.2（审计腿 agent-104 判定的反向谎言 #2）
- **修复提交**：见派单回话首行（本文件与代码同一枚 commit）

## 谎言与真相

- **原话（反向谎言）**：ABSENCES 行写「预算字段不在 wire 上：没有预算读数就没有预算条」——把「字段不在 payload」说成了「没有预算读数」，而预算判定 Go 侧是活的。
- **真相（现核工作树，file:line 可复跑）**：
  - `internal/config/schema.go:608-613` — `CostSection`（`[cost]` 的 `daily_budget` / `monthly_budget` / `alert_threshold` 默认 0.8 / `over_budget` 默认 pause）；
  - `internal/observe/costbudget.go:103` — `BudgetParams.Judge(day, month)` 出判定（80% 警告 / 100% 超支）；
  - `internal/agent/costbook_44.go:147` — `l.publish(Event{Kind: EvWarn, ...})`（警告真上总线）。
- **缺口在哪**：`cost.summary` payload（`internal/panel/cost_handlers.go:124` costPayload）**尚未携带预算字段** ⇒ 判定结果到不了面板。等 `readface1` 单把字段加进 payload 后再接读面（另单）。
- **⛔ 造数字禁令**：字段落 wire 之前一格预算数字都不画——本修复只改文案，不添任何预算读数、不画零长度棒子。

## 改了哪里

- `frontend/src/components/cost-settings.tsx:29-35` — 头注「预算条」bullet 改真话（Go 侧判定活、payload 未携带、等字段）。
- `frontend/src/components/cost-settings.tsx:107-110` — ABSENCES「预算条（80% 警告 / 100% 暂停）」reason 改真话，保留「未接线」StateTag。

## 硬约束自证

- 未动 `docs/PLAN.md`／`docs/specs/**`／thresholds／golden／allowlist／D43／design 契约文字。
- 未放宽任何断言；未 push；未提交 `frontend/dist` 字节（构建后删回 `.gitkeep`）；UI 零 emoji。

## 串行读数（node = D:\work\server\node24，v24.21.0 / npm 11.19.0）

```
$ npm --prefix frontend run build          rc=0
$ npm --prefix frontend run typecheck      rc=0
$ npm --prefix frontend run lint           rc=0（38 warnings, 0 errors——存量告警，非本次引入）
$ npm --prefix frontend run render:history36  rc=0（6 XSS families neutralized，含新负控）
```
