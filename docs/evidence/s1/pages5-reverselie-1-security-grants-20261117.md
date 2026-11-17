# 反向谎言 #1 修复 · `security-settings.tsx` 会话授权「零写入」改真话

- **修复腿**：`agent-105`（写腿，frontend 写面；实现者，⛔ 非裁决者）
- **时间**：2026-11-17 +08
- **来路**：`docs/evidence/s1/pages5-adversarial-20261117-0953.md` §7.1（审计腿 agent-104 判定的反向谎言 #1）
- **修复提交**：见派单回话首行（本文件与代码同一枚 commit）

## 谎言与真相

- **原话（反向谎言）**：ABSENCES 行写「D45 会话授权生产零写入（grant_id 生产逐字传 nil，internal/agent/journal.go:101）」，头注同句。
- **真相（现核工作树，全部 file:line 可复跑）**：会话授权**生产在写**：
  - `internal/agent/approval/gate.go:859` — `g.grants.Record(...)`（路径授权行）；
  - `internal/agent/approval/gate.go:844` — `rt.RecordTarget(...)`（目标进程绑定行）；
  - `internal/session/grants.go:201` / `:274` — 两扇门都落到 `l.store.InsertGrant(...)`；
  - `internal/tools/bridge.go:1445` — `tc.GrantID = &g`（grant_id 带回调用面）；
  - 会话收口：`internal/session/grants.go:346-367` `RevokeAll` 盖 `revoked_at`、**行保留审计**（`internal/memory/dao_misc.go:48-67` `RevokeGrant` 是幂等盖章 UPDATE，不删行）。
- **仍然缺的（absence 行保留的理由）**：**面板侧没有 grants 列表读面**——快照无 grants 节；架构禁 Go 入站 grant 路由（`internal/session/grants.go:39-40` 自注「does not list or revoke grants for a UI」，`cmd/wisp/panel_inbound_fuzz_35g1_test.go:160` 把 `grants.revoke` 当非法入站名钉住）。
- **操作出口**：console 动词 `grants`（列本会话生效授权，`cmd/wisp/approval_reply.go:415`）与 `revoke-grant <编号>`（`cmd/wisp/approval_reply.go:474` 起，派发处 `:991-994`）。

## 改了哪里

- `frontend/src/components/security-settings.tsx:33-48` — 头注「zero production writers」段改真话（写面在、缺面板读面、出口＝console 动词，全部带 file:line）。
- `frontend/src/components/security-settings.tsx:70-75` — ABSENCES 该行 label 改「会话授权面板读面」，reason 改真话（生产在写 + 面板侧无读面 + 出口动词），保留「未接线」StateTag（缺的只是面板读面）。

## 硬约束自证

- 未动 `docs/PLAN.md`／`docs/specs/**`／thresholds／golden／allowlist／D43／design 契约文字。
- 未放宽任何断言；未 push；未提交 `frontend/dist` 字节（构建后删回 `.gitkeep`）；UI 零 emoji（本文件新增文案无 emoji，箭头用全角＝号与文字）。

## 串行读数（node = D:\work\server\node24，v24.21.0 / npm 11.19.0）

```
$ npm --prefix frontend run build          rc=0
$ npm --prefix frontend run typecheck      rc=0
$ npm --prefix frontend run lint           rc=0（38 warnings, 0 errors——存量告警，非本次引入）
$ npm --prefix frontend run render:history36  rc=0（6 XSS families neutralized，含新负控）
```
