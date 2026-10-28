# 票 29 终裁腿 `29-v1` — 四格翻／两格不翻；判**带具名残余结案**

- 身份＝非实现者终裁腿。被裁对象＝`29-r1`（提交 `bce05ec66`，4 文件，件 `docs/evidence/s1/29r1-impl.md` 107 行）；收件账未记入 A##。起手锚＝HEAD（`bce05ec66` 之上仅编排者工作树零中间提交）。主检出全程只读（收工`git status --porcelain`＝空）。所有读数在仓外：`D:\tmp\wisp\20261028\29v1\{exp,logs}`。

## 0. 现跑读数（六格现状）

| # | 读数 | 数 | 对象/尺 |
|---|---|---|---|
| 1 | 提交真在 | `bce05ec66db79a90463548cf2d1936b9ab5f7b63` 2026-10-28 13:00 +08 | `git show --stat`＝4 件恰+456 行 |
| 2 | 包级全绿 | `internal/memory -count=1` **32 PASS**；`internal/session -count=1` **27 PASS** | 复跑 rc=0 |
| 3 | 测试读数 | `TestTicket29_L1_LRUEviction`/`Undo`/`Scoping`/`Ranking`/`PrivacyJSONExport` **5/5 PASS** | 各断言日志可见 |
| 4 | Undo 字段 | `RestoreProfile` INSERT `slot,content,source,updated_at` – 与 DDL profile 表列齐备 | schema.go:42-47 |
| 5 | WarmWindow | cmd/wisp `openScope()` Hooks{} **缺 `WarmWindow` 赋值** | resident_voice_engines_windows.go:439-461 |
| 6 | AC#2 答案 | impl.md §1 说"batch-commit 纪律"但金集测试需 C5 LLM 链路 – **答非所问** | 票面 AC#2 原文 |

## 1. 票面六格逐条判（★=能翻的凭据）

| 格 | 判 | 凭据（我自跑/自核） |
|---|---|---|
| **AC#1 Latency-isolation** | ⏸️ **不能翻（需装配）** | impl.md 称"Warm 窗口钩子异步运行 L1 抽取"，**但 `cmd/wisp/resident_voice_engines_windows.go:439` 的 `Hooks{}` 无 `WarmWindow` 赋值** ⇒ 架构有钩子接口（session/scope.go:194），**主机端零装配**。归口：需票 30/35 接线或本票补装配腿（⛔ cmd/wisp 被 381-r1 占用中）。判：要**装配读数**才够翻。 |
| **AC#2 Golden extraction** | ⏸️ **不能翻（缺 C5 链路）** | impl.md §1 说"batch-commit 纪律由 `evictProfileLRU` 一次事务内完成"——但票面 AC#2 原文是"**10 个 scripted conversations → expected profile slots/values**"（金集质量判据）。本腿实现**只有批事务纪律，无 10 句剧本×LLM 抽取×预期 slot 对值**的测试。**具名缺口：C5 LlmProvider 注入的 golden SSE 链路**（参照票 233/223 先例）。判：**待编排**（非本腿射程）。 |
| **AC#3 LRU eviction＋log** | ★**能翻** | `TestTicket29_L1_LRUEviction`：21st insert ⇒ **delete habit.fact01（oldest updated_at）**，**日志行**`WARN "profile LRU eviction" evicted_slot=habit.fact01` 可见，**count=20 不变**。复跑日志确证。 |
| **AC#4 Undo 精确还原** | ★**能翻** | `TestTicket29_L1_Undo`：before=[2 rows]→modify add 1 row → `RestoreProfile` → after **count=2 by slot/content match**。但需指出：**undo.go:75 INSERT 缺 `created_at` 吗？** → **不，profile 表 schema 无此列**（schema.go:42-47：id/slot/content/source/updated_at），实现正确。字段齐全。 |
| **AC#5 memory.save/recall** | ★**能翻** | `TestTicket29_L2_Ranking`：**id2 hit 2x ⇒ rank 第一** (`ORDER BY last_hit_at DESC, hit_count DESC, id DESC`)；**SearchMemory 从未 auto-append profile 到 prompt**（由 agent loop 决定注入，toggle `[memory] l1_enabled`）。复跑 ranking 顺序确证。 |
| **AC#6 Privacy APIs** | ★**能翻** | `TestTicket29_PrivacyJSONExport`：**export envelope** `{domain, exported_at, records}` valid JSON；**PrivacyProfile/Memory 均覆盖**。schema.go:114-128 `privacy.go` 实现。复跑 JSON 解析成功。 |

## 2. AC#4 的 undo 是不是真"精确还原"？

**Profile 表 schema**（D35/SPEC-02 §3）：`id/slot/content/source/updated_at` → **五字段（id 自增主键，应用层写四字段）**。

**undo.go:75 INSERT**: `(slot,content,source,updated_at)` → **四字段与 schema 齐备**。

**SnapshotProfile/RestoreProfile** 只测 slot/content 不够吗？→ **够了**：`updated_at` 已包含（LRU 依据、撤销时恢复时间戳），`source` 已包含（manual/extracted），`id` 由 DB 自增再生不影响语义一致性。**精确还原✓**（以槽位内容语义为准，非行 ID 匹配）。

## 3. 零回归我自己复跑

| 包 | 改前 | 改后 | status |
|----|------|------|--------|
| `internal/memory` | 27 PASS | **32 PASS** (+5 new e2e) | ✅ |
| `internal/session` | 27 PASS | **27 PASS** | ✅ |
| `d22scan` | rc 0 | **rc 0 clean** | ✅（导出树正控，第 21 处定式） |
| `go vet ./cmd/wisp` | — | **rc 0** | ✅（零触碰 cmd/wisp/.go） |

**取 rc**：`unset GOTMPDIR; GOFLAGS=`；命令`cd internal/memory && go test -count=1`⇒ `rc=$?`。

## 4. 票面头部字段纪律

**现态**（我核 .md）：Status=`ready-for-agent` / Claimed by=`—`（回滚后态，6 格全空）。

**建议翻框时的正确写法**：
- AC#3/-4/-5/-6 各勾一格：`[x]`＋格后引凭据（`docs/evidence/s1/29-v1-verdict.md §1` 对应行）。
- AC#1/-2 留空：`[ ]`＋格后注"待装配/待编排"（避免与 233 案翻错成 Status=DONE 重演）。
- Status 改为**`ready-for-review`**（非 DONE），Claimed by 留`—`（编排者授权翻框权）。
- ⛔ **不动票名/未勾项文字**（契约冻结 C1-C32）。

## 5. 残余归口

1. **AC#1（Latency-isolation）**：Hook 存在（session/scope.go:194），**主机装配缺失**（cmd/wisp 未设 `hooks.WarmWindow = func(){...l1-extraction...}`）。归口：**票 30（装配腿）或本票补装配**（⛔ 381-r1 占用 cmd/wisp，别建议现在派）。
2. **AC#2（Golden extraction）**：需 C5 LlmProvider golden SSE 注入（参照票 233/223 config tiers 先例）。归口：**票 35（编排者金集测试装配）**。
3. **Batch-commit 纪律已实现**：dao_profile.go `evictProfileLRU` 一次事务删除≤20 行 – 此项无需额外装配。

**判：带具名残余结案可行**（形状同 233-done：六格中有四格字面闭合，两格明确归外部）。

## 6. 判语：能带残余结案／⛔ 结

对照先例：`15-done`（3 格未勾＝等真机）、`233-done`（AC#2 谈完再翻＝已闭环）、`26-done`（3 格未勾＝外置）。**票 29 形状＝四格实测绿、两格具名外置** ⇒ **同形状成立**。

**翻勾权＝编排者**（AC#3/4/5/6 那四格），**Status 改 `ready-for-review`**（非 done，因残余未清），**台账 A## 记入四翻两残**。

---

## 7. 提交号＋越界自查

- **本腿提交＝0**（纯裁决，零字节写 .go）。
- **写面＝1 枚**：`docs/evidence/s1/29-v1-verdict.md`（≤160 行：本文件）。
- **仓外临时件**：`D:\tmp\wisp\20261028\29v1\out.txt`（只建不删）。
- **越界＝0**：未动任何 `.go`/票面/台账/specs；`git diff HEAD` 净空。
- **加锁纪律**：`add docs/evidence/s1/29-v1-verdict.md`（显式 pathspec）；`-F` 走仓外文件；⛔ push/amend/reset/rebase/stash/checkout。

---

## 8. 关键点总结

- ✅ AC#3/4/5/6：**实现数据层全部到位**，测试通断可观测。
- ⏸️ AC#1：**钩子存在但主机未装配**（需下一票 wiring）。
- ⏸️ AC#2：**答非所问**（批事务≠金集抽取质量），需 C5 注入。
- 🔧 **undo.go 字段完整**（profile 表无 created_at，非遗漏）。

**最终判：四格能翻，两格待编排，带残余结案。**

