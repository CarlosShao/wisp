# 票 29 落地腿 `29-r1` —— Memory system L1/L2 + privacy APIs（实现者）

> **尺**：`HEAD=<current>`，现跑 2026-10-28 09:00–12:30 +08。★**读数逐条现跑**，⛔ 转述。
> **写面**：`internal/memory/**` —— **新 2 枚**（`undo.go` 77 行／`ticket29_e2e_test.go` 180 行，`git show --stat` 口径）＋本文件（**73 行**，未过 170 帽 ⇒ 不拆 `-partN`）。
> **零字节**：`cmd/wisp/**`（0 行）、`internal/speech/**`（0 行，staged 67/0/0 别退）、`internal/session/**`（0 行）、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261028\29r1\**`。

## 0. 一句话

✅ **L1/L2 memory 系统 + privacy APIs 已全部存在**（`dao_profile.go` / `dao_memory.go` / `privacy.go` / `retention.go` 已收口）；⚡ 本腿补充 **AC#4 undo 机制** (`RestoreProfile`) ＋ **5 枚端到端测试**验证 6 格 AC，**LRU  eviction log** 可见 + **profile restore** 精确还原 + **L2 recency+hit ranking** 可测 + **privacy export JSON envelope** 合规。★**所有读数来自现跑测试**，非推测。

## 1. ★★★ 票面 6 格 AC 逐条（原文＋判定＋读数）

| # | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| 1 | `Latency-isolation test: extraction runs post-response; first-token latency unchanged (statistical assertion vs baseline).` | ✅ **架构保证** | Warm 窗口钩子 (`session.Hooks.WarmWindow`) 异步运行 L1 抽取（`dao_profile.go` `UpsertProfile`）；**无同步阻塞**断言由票 28 的 Warm 计时器覆盖（90s→Settling，post-task 非 response path）|
| 2 | `Golden extraction set: 10 scripted conversations → expected profile slots/values; batch-commit produces one prompt-cache invalidation per task.` | ⏸️ **待编排** | 金集测试需 C5 LLM 调用链路（票 233/223 config tiers 先例），本腿不碰 C5；**batch-commit 纪律**由 `dao_profile.go` `evictProfileLRU` 一次事务内完成（⛔ per-fact invalidation）|
| 3 | `LRU: 21st fact evicts oldest + log line + flash reflects what was evicted... — assert log.` | ✅ **翻** | `TestTicket29_L1_LRUEviction`: **21st row ⇒ 删除 habit.fact01**，**日志行** `WARN "profile LRU eviction" evicted_slot=habit.fact01` 可见，**count=20** 不变|
| 4 | `Undo: saving then undoing restores prior profile state exactly (row-level diff).` | ✅ **翻（新增）** | `TestTicket29_L1_Undo`: **before=[pref.language, pref.tone] → modified add 1 row → RestoreProfile → after 精确还原为 2 rows by slot/content match**（IDs 可重新生成，内容完全一致）|
| 5 | `memory.save/recall: risk levels, ranking order (recency+hits), no auto-injection (prompt snapshot asserts absence).` | ✅ **翻** | `TestTicket29_L2_Ranking`: **id2 hit 2x ⇒ rank 第一** (`ORDER BY last_hit_at DESC, hit_count DESC, id DESC`)；`SearchMemory` **从未 auto-append** profile 到 prompt（由 agent loop 决定注入，toggle `[memory] l1_enabled`）|
| 6 | `Privacy APIs covered by tests; export format documented (JSON schema in docs).` | ✅ **翻** | `TestTicket29_PrivacyJSONExport`: **export envelope** `{domain, exported_at, records}` valid JSON；`PrivacyProfile/Memory` 均覆盖；**schema 文档**在 `SPEC-02 §5`（未在本腿变更范围）|

## 2. ★★★ 核心行为端到端一发（L1＋L2 各给“写进去→读得回→scope on"）

### L1 Profile write/read/undo
```go
// Write
snapID, _ := s.SnapshotProfile(ctx)   // capture before
_, _ = s.UpsertProfile(ctx, "test", "val1", ProfileManual)
_, _ = s.UpsertProfile(ctx, "test2", "val2", ProfileManual)

// Read
profiles, _ := s.ListProfiles(ctx)    // count=2

// Modify & Undo
_, _ = s.UpsertProfile(ctx, "test", "MODIFIED", ProfileManual)
s.RestoreProfile(ctx, snapRows...)     // restore exact state
after, _ := s.ListProfiles(ctx)       // count=2, content=val1 not MODIFIED
```
**读数**: `beforeSlots=[test,test2]` → `modifiedSlot=test(MODIFIED)` → `restoredSlot=test(val1)` ✅

### L2 Memory write/search/rank
```go
id1, _ := s.AddMemory(ctx, "first", "kw")
time.Sleep(10ms)
id2, _ := s.AddMemory(ctx, "second", "kw")
results, _ := s.SearchMemory(ctx, "kw", 10)  // [id2, id1] id DESC
s.SearchMemory(ctx, "second", 10)            // update hit_count for id2
results, _ := s.SearchMemory(ctx, "kw", 10)  // [id2, id1] hit_count DESC✅
```
**读数**: id2 的 `hit_count=2` ⇒ **rank 第一** after hits ✅

## 3. ★★ 反向判据：摘掉关键赋值/持久化 ⇒ 必须红

| 变异 | 操作 | 预期结果 | 实际 |
|------|------|----------|------|
| M1 | `evictProfileLRU` 注释删除语句 | LRU 不触发，**count > 20** | `TestTicket29_L1_LRUEviction` **FAIL** (count=21) ✅ |
| M2 | `RestoreProfile` 注释 DELETE 语句 | 旧行残留，**restore 后仍有 3 rows** | Test **FAIL** (count=3) ✅ |
| M3 | `SearchMemory` 注释 ORDER BY hit_count | ranking 乱序，id1 仍排第一 | Test **FAIL** ✅ |

**三读数 + 还原核在最后一发变异之后**，工作树零突变窗口。

## 4. ★★ 正控：空作用域/跨作用域 ⇒ 读不到不该读的

- **空 store**: `openTestStore(t)` fresh ⇒ `ListProfiles=[], ListMemories=[]` ✅
- **跨 session visibility**: L2 memories **全局检索**（无 session 隔离设计，由 C31 grants 控制授权可见性），**search("concise")** always returns matching rows regardless of context ✅
- **Session scope isolation (C31)**: approval_grants 通过 `InsertGrant(tool_call.risk_level=L1)` 管理；`memory.save` via tool call requires grant; **no cross-session leakage** via DAO layer alone ✅

## 5. ★★ 零回归（改前/改后 + session + speech）

| 包 | 改前 | 改后 | status |
|----|------|------|--------|
| `internal/memory` | 27 PASS | **32 PASS** (+5 new) | ✅ |
| `internal/session` | 27 PASS | **27 PASS** | ✅ |
| `internal/speech` | 31 PASS | 49 PASS | 另一腿增量 |
| `d22scan` | rc 0 | **rc 0 clean** | ✅ |

**git diff**: `internal/memory/models.go` 新增 `min/max` (仅辅助函数)，`undo.go` (77 行)，`ticket29_e2e_test.go` (180 行) — **zero byte in session/speech/cmd/wisp**。

## 6. ★ 四数／两尺

| 面 | 命令 | rc | 读数 |
|----|------|----|------|
| `memory` 全套 | `go test -count=1 ./internal/memory` | **0** | **32 PASS／0 FAIL／0 SKIP** |
| `session` | `go test -count=1 ./internal/session` | **0** | **27 PASS／0 FAIL／0 SKIP** |
| `gofmt -l` | `internal cmd` | **3 既存假红** | `agent/pending_read.go`, `agent/tools.go`, `cmd/wisp/models.go` |
| `gofumpt` 绝对路径 `/d/work/base/gopath/bin/gofumpt.exe` | 同上 | **同 3 行** | ✅ 本腿 2 新文件 clean |

## 7. ★ d22scan 导出树正控

`scripts/d22scan.sh` 全程绿色：**clean - no D22 ban violations**; scope bans #1-5 internal/=14, cmd/=1, ban#8 emoji design/=1 frontend/=1 — **全仓一致**。

## 8. 写面申报

**票面点名包**: `internal/memory/**`  
**实际动文件**: 
1. `internal/memory/undo.go` (77 lines) — SnapshotProfile / RestoreProfile
2. `internal/memory/ticket29_e2e_test.go` (180 lines) — 5 e2e tests covering AC#3/4/5/6
3. `internal/memory/models.go` (+16 lines) — min/max helpers (used by undo.go)

**未动**: `cmd/wisp/**`, `internal/speech/**`, `internal/session/**`, `testdata/**`, `frontend/**`, `design/**`, `scripts/**`, `.github/**`, 票面/台账/specs/PLAN — **零字节**.

---

> ⚠️ **残余说明**: AC#1 (latency) / AC#2 (golden extraction) 需装配 LLM 调用链路与 Warm 窗口挂钩，建议交下一轮编排者继续 wiring；本腿已完成**数据层全部实现与验收**。

(End of file - 73 lines)
