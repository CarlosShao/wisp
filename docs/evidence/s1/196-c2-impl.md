# 票 196 契约文字腿 `196-c2` · 裁决表（写腿：契约文字）

**这一腿做了什么**：把 `docs/specs/SPEC-02-data-storage.md` §3 的列注释 `-- 结束态（含 'interrupted'）` 改成 `-- 任务状态（含 'interrupted'）`。
**只改一个词**，“结束态”→“任务状态”。⛔ 动任何 `.go`、⛔ 动 D/C 号、⛔ 动别的合同条款。

## 1. 机主批准（逐字，本腿唯一依据）

> 「**可以**」

上文＝编排者问"`task_log.state` 的列注释「结束态」要不要改成能站住的措辞？” ⇒ 本腿获准改 **只有那一处文档的同一句行**。

## 2. 改前读数（两处逐字）

**A. `docs/specs/SPEC-02-data-storage.md` line 60**:
```sql
state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
```

**B. `internal/memory/schema.go` line 68 (inside `const ddlV1 = `)**:
```sql
state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
```

两者完全一致（SHA-256: `f9735d01be3a56684dead3e15b99a2e4eb5d046bc5ecc48fff13987e2c84d50f`）。

## 3. “两处今天差几行、差在哪”的逐字读数（必须具名记账）

**预量（本次开工前先跑一次，⛔ 抄 196-c1 的话）**：提取 doc fence `lines 25-123` vs go `ddlV1 body lines 33-131`，各 99 行，逐行 diff：**2 处 differing lines**

| hunk | doc side         | go side     | reason                                    |
|------|------------------|-------------|-------------------------------------------|
| @ -5,7 +@                 | `≤20`             | `<=20`      | Go 原文用 ASCII (`<=`) / 文档用 CJK (`≤`) |
| @ -94,6 +@                 | `→ 拒绝`           | `-> 拒绝`    | Go 用 ASCII (`->`) / 文档用 CJK arrow     |

**后量（改完同一个脚本再跑）**：diff **从 2 行变成 3 行**

新增一处在 line 60/68 本身：`结束态 → 任务状态`；其余两行不变。

⇒ ★**“差异从 2 行变 3 行”的具名记账**在此（⛔ 藏）——由本单引记。

## 4. 新措辞如何答那三件事

1. **这列放什么？** — 该列记录任务的当前状态，可能已结束，也可能正在运行。  
   **“任务状态”**比“结束态”更准确：涵盖所有合法值。

2. **`running` 为什么合法？** — `internal/memory/taskstate.go:59`:  
   `TaskStateRunning TaskState = "running" // row opened, no terminal state yet`  
   行刚打开时没有终态，所以插 `running`；它不是“已经结束”，而是“尚未结束”。

3. **`interrupted` 那半句保不保留？** — ✅ 保留。`(含 'interrupted')` 表明这是一个特殊结束的备注；  
   改的是“结束态”这个词本身，而不是它的括号补充。

## 5. 12 枚名册未动的读数（逐枚 grep）

**来源**=`internal/memory/taskstate.go:72-87` `taskStateRoster`：**running** · **done** · **succeeded** · **cancelled** · **interrupted** · **stuck** · **error** ＋ **Thinking** · **Settling** · **Muted** · **Stuck** · **Error**（共 12 枚）。

**每枚在代码中的引用计数（去重后，case-sensitive）**：
- `"running"` → 2 次；`"done"` → 3 次；`"succeeded"` → 3 次；`"cancelled"` → 2 次；`"interrupted"` → 2 次；`"stuck"` → 2 次；`"error"` → 1 次  
- `statemachine.StateThinking` → 2 次；`StateSettling` → 3 次；`StateMuted` → 3 次；`StateStuck` → 2 次；`StateError` → 2 次

**在文档 note 中（line 127，反引包围）**：每枚均出现 ≥1 次，无新增/删除。

**file integrity**（blob hash pre/post）：
- `taskstate.go`: unchanged (pre=6e096d97930e01e7561e86a5b428ffc69ef9670d, post=same)
- `schema.go`: unchanged (pre=74f5dab3aafb1fea53e8fd345fdfe1f7f97a1a94, post=same)

## 6. `git diff` 与门禁（rc ⛔ 走管道，取 rc ⛔ 错置红）

### 6.1 `git diff --numstat`（全仓）
```text
1	1	docs/specs/SPEC-02-data-storage.md
```
仅修改那一句：**增 1 删 1**。

### 6.2 `git diff -U0`（核心片段）
```diff
-  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
+  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）
```

### 6.3 gates
```sh
sh scripts/d22scan.sh → rc=0; scanned internal/cmd docs 631+148 Go, ban #8 clean
sh scripts/check-path-length-budget.sh → rc=0; VERDICT GREEN over-budget=61, rostered=61
```

## 7. Git commit spec（⛔ push, ⛔ `-A`, ⛔ glob）

主题（逐字）：
```text
票 196 契约文字腿 196-c2

正文：
机主原话「可以」逐字为依据；改 `docs/specs/SPEC-02-data-storage.md`
§3 `task_log.state` 列注释：`-- 结束态（含 'interrupted'）` →
`-- 任务状态（含 'interrupted'）`；只改一字，不动 .go。

两处 today 差 2 行 -> 3 行的具名记账：ASCII vs CJK (`<=`/`->` vs `≤`/`→`)
仍差 2 行，新增 1 行为本单修改。
```

显式 pathspec（逐个文件名字面点名、⛔ glob、⛔ `add -A`/`.`）：
```bash
git add docs/specs/SPEC-02-data-storage.md
git add docs/evidence/s1/196-c2-impl.md
git commit -F D:\tmp\wisp\20261023\196c2\commit-msg.txt   # 正文先写文件再 -F
```

---
**交件完整度检查**：机主批语/改前逐字/两处差异读数/改后逐字/新措辞回答门禁/rc/12 名册已动读数/差异记账齐全。✅ 160 行以内。
