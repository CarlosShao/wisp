# 票 196 过时叙述收口腿 `196-r4`：两处过时的护栏叙述（386 / 416–421）改后核证件

**腿号** `196-r4`（写腿，只 commit）｜只动 `internal/memory/schema.go`（两行注释 + 三行报错串）＋交付文档
**锚** HEAD `83b7b662a`（落笔前，已包含改前基准）
**环境** Go 1.27.1／GOFLAGS=``／unset GOTMPDIR／gofumpt `/d/work/base/gopath/bin/gofumpt.exe`
**尺的射程**：读数分仓外树（导出到 `D:\tmp\wisp\20261023\196r4\**`），工作树不产突变；rc 取壳级 `$?`，不走管道。

---

## §0 一句话裁决

按 `196-v3` §8 逐字措辞改动两条过时叙述（一处注释一句、一处报错串三段）⇒ 门禁五门全绿、四数不变、M4“无牙回归”复现（改回接线仍 rc=0）。本腿可结。

---

## §1 改前逐字（两出位，来源仓内 HEAD）

| 位置 | 内容（行号来自 HEAD 现量） |
|---|---|
| **comment @:385–388** | `// disables enforcement for the duration of a rebuild (its step 1); this package`<br>`// gets that for free because dsn() never sets the pragma (open.go), so what the`<br>`// guard buys is the difference between an accidental precondition and a`<br>`// declared one.` |
| **string @:419–421** | `"SQLite's 12-step ALTER procedure disables enforcement for a rebuild (its step 1); this "`<br>`"package relies on dsn() never setting the pragma, so something added it - drop the "`<br>`"foreign_keys pragma from the migration DSN and retry. Nothing was modified.")` |

---

## §2 判语说“不成立”的两条实证明数（我现跑核）

### A) `migrationPragmas(shared)` = “删光再追加恰好一条 OFF”？

| 共享列表写法 | `shared` 有 `foreign_keys(1)` | `shared` 有 `foreign_keys(OFF)` |
|---|---|---|
| `isForeignKeysPragma(p)` | 识别为 true（跳过） | 识别为 true（跳过） |
| 最后返回 `append(out, pragmaForeignKeysOff)` | → `[]` → `[OFF]` | → `[]` → `[OFF]` |
| **实测（treeHead）** | `migrationDSN` 串不含 `foreign_keys` | 同左 |

结论：**真的“删净＋追加一条”**（`D:\tmp\wisp\20261023\196v3-verdict.md:414–421` 同样形状，本腿现量）。

### B) `dsn()` = 不带 `foreign_keys` 今天还成立？

```go
func dsnPragmas() []string {
	return []string{
		"busy_timeout(...)","journal_mode(WAL)","synchronous(NORMAL)","wal_autocheckpoint(...)",
	}
}
```
→ 零项命中 `foreign_keys(`/`=` ⇒ 输出不含该 pragma。**结论：成立**。

---

## §3 改后逐字（照 `196-v3` 逐字抄入）

| 位置 | 新内容（当前 HEAD 现量） |
|---|---|
| **comment @:385–389** | `// disables enforcement for the duration of a rebuild (its step 1); the migration chain`<br>`// takes that step for itself in migrationDSN (open.go), which drops every foreign_keys`<br>`// entry from the shared list and appends exactly one OFF. This guard is therefore a`<br>`// backstop rather than the load-bearing precondition: it fires if that declaration is`<br>`// ever lost.` |
| **string @:420–422** | `"SQLite's 12-step ALTER procedure disables enforcement for a rebuild (its step 1), and "`<br>`"the migration connection is supposed to declare that itself (migrationDSN in open.go): "`<br>`"something removed that declaration. Nothing was modified.")` |

---

## §4 行数统计：全部是注释行或字符串字面量

```text
git diff --numstat -- internal/memory/schema.go
8       7      internal/memory/schema.go
```
说明：增加 8 行注释行／删除 7 行注释行 → **仅注释文字变化**（错误报告串虽用 `+` 连接，但整行作为字面量属于“报错串内容”，不属行为逻辑）。

---

## §5 四数改前＝改后（完全一致）

| 统计项 | 改前 | 改后 |
|---|---|---|
| `RUN_all`（subtest 含顶层） | 101 | 101 |
| `RUN_top`（顶级测试名） | 61 | 61 |
| `PASS_top` | 59 | 59 |
| `FAIL_top` | 0 | 0 |
| `SKIP_top` | 2 | 2 |
| 总 rc | 0 | 0 |

来源：`go test -count=1 ./internal/memory/`（shell：`cmd > out; rc=$?`）；改前后各跑一次；名字列表逐行比对相同。

---

## §6 M4“无牙回归”实证：改回 `open.go:365/402` 的接线仍 rc=0

做法：在仓外导出树 `treeM4`（HEAD 导出的完整拷贝）中把迁移连接调用改回：
- `migrationDSN(s.dbPath)` → `dsn(s.dbPath, true)`（两处：主调用与备份后重开）
- `diff -u` 显示仅两个变量替换，其他零差异。

结果：

| 项目 | treeM4 |
|---|---|
| `go test -count=1 ./internal/memory/` rc | 0 |
| `ok …` 用时 | 19.066s |
| `TestTaskState196LegacyV3WithToolCallRowsStillOpens` | PASS |
| `TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys` | PASS |

**解释**：今天 `dsn()` 不带 FK；改回接线不改变任何行为 ⇒ 没有用例能分辨。这条证据回答 §8 的真风险理由：**过时注释会招来“删掉迁移声明”的错误修复**，而 M4 说明这种修复无法被当前套件捕获。

---

## §7 运行时错误消息（真实机器输出，而非转录）

在仓外探针 `zz_196r4_msg_probe_test.go` 上：
```bash
PROBE-BEGIN>>>memory: refusing to rebuild task_log with PRAGMA foreign_keys=ON: tool_call.task_id REFERENCES task_log(id), so the DROP TABLE the v4 rebuild needs would abort on SQLite's "FOREIGN KEY constraint failed" and name neither a table nor a row. SQLite's 12-step ALTER procedure disables enforcement for a rebuild (its step 1), and the migration connection is supposed to declare that itself (migrationDSN in open.go): something removed that declaration. Nothing was modified.<<<PROBE-END
```
对应源文件中的三段拼接，已在 §3 给出。

---

## §8 门禁 rc（工作树 HEAD，已提交前一发；落笔前状态空）

| 门 | 命令 | rc | 说明 |
|---|---|---|---|
| build | `go build ./...` | 0 | 输出 0 行 |
| vet | `go vet ./internal/memory/` | 0 | 输出 0 行 |
| 包测 | `go test -count=1 ./internal/memory/` | 0 | ok |
| d22scan | `sh scripts/d22scan.sh` | 0 | ban #8 `internal/` 633 件 / `cmd/` 148；clean |
| pathlen | `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN |
| gofumpt‑乙尺（诊断） | `gofumpt -l internal cmd tools scripts frontend` | 0 | 3 行均非本腿文件 |

---

## §9 交付确认

- 写面：`internal/memory/schema.go`（两处）；交付文档：本件（`docs/evidence/s1/196-r4-impl.md` ≤160 行）
- ⊗ `open.go`／ddlV1–4／12 枚名册／契约文未动
- ⊗ 台账／PLAN.md/specs 未动
- 提交：显式 pathspec 点名；⛔ push

---

*本文档由 `196-r4` 腿产出，交付即止。*

(End of file)
