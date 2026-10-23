# 票 196 代码侧对齐腿 `196-c3` —— 实现件

**腿** `196-c3`（写腿）· 2026-10-23 · 仅 commit、⛔ push
**一件事**：`ddlV1` 里 `task_log.state` 的**行尾列注释**与契约侧 `SPEC-02:60` 对齐（「结束态」→「任务状态」），并把测试里那句"untouched"的**意图注释**改准。
**⛔ 未改**：任何 SQL 语句／标识符／列定义；`ddlV2`/`ddlV3`/`ddlV4`/`applyV4`/`buildDDLV4`；任何断言条件；契约文字。
**写面**：`internal/memory/schema.go`（只 `:68` 行尾 `--` 注释）· `internal/memory/taskstate_196_test.go`（只 `:447-448` 注释）· 本件。

## 1. 「改前逐字」

```go
schema.go:68                      (前导 2 空格；文件共 411 行)
  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）

taskstate_196_test.go:447-448     (前导 TAB；文件共 452 行)
	// ddlV1 must be untouched: this ticket adds a step, it does not edit the
	// frozen contract text.
```

## 2. 「改后逐字」

```go
schema.go:68
  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）

taskstate_196_test.go:447-448
	// ddlV1's structure and statements must be untouched: this ticket adds a step,
	// it does not edit the DDL; one wording-only comment was corrected (196-c2/196-c3).
```

`schema.go` **仍 411 行**（字符数 64→65）；测试**仍两行** ⇒ 断言仍在 `:449`，行号未漂。
意图注释改准三点：①「什么不许动」＝`structure and statements`；②「什么被动了」＝`one wording-only comment`；③依据＝`(196-c2/196-c3)`。原句 `frozen contract text` 在改注释后已成假话，故删。

## 3. 「三件已核」——我自己现跑复核（⛔ 抄派单）

**① 无测试钉 `ddlV1` 字面或哈希**：

```
$ grep -rn "sha256.*ddlV1\|ddlV1 ==\|== ddlV1\|ddlV1 :=\|ddlV1," --include="*_test.go" internal/
internal/memory/taskstate_196_test.go:307:	for _, script := range []string{ddlV1, ddlV2, ddlV3} {
```

⇒ 唯一命中是把 `ddlV1` 当**值**喂给 `seedLegacyV3`，**不是**文本/哈希比对；`sha256.*ddlV1`／`ddlV1 ==`／`== ddlV1` **各 0 命中**。成立。

**②③ 两处 `ddlV1` 断言都是语句数**（`taskstate_196_test.go:449`、`schema_test.go:264`，逐字同一行）：

```go
	if n := len(splitSQLStatements(ddlV1)); n != 15 {
```

**我另补两条射程核查**（决定"注释改动会不会红"）：
- `splitSQLStatements`（`schema.go:391-411`）**按 `;` 切分**再丢纯注释块；前后两版注释**都不含 `;`** ⇒ 切分点不变。
- `schema_test.go:59` 的 `contractSpec` 是**结构名册**（表/列/索引名、notNull、pk），**不含注释文字** ⇒ 措辞不进它的射程。

## 4. `splitSQLStatements(ddlV1) == 15` —— 现跑读数（改前＋改后）

**尺**：⛔ 不是我自己重写切分器。用 `go test -overlay=…` 把临时探针**注入包内**（探针在 `D:\tmp\…`，仓里不落盘），调**真实的** `splitSQLStatements` 与**真实的** `ddlV1`。

| | 读数 | rc |
|---|---|---|
| **改前** | `RESULT splitSQLStatements(ddlV1) len = 15`；纯注释开头语句 `= 8` | 0（PASS） |
| **改后** | `RESULT splitSQLStatements(ddlV1) len = 15`；纯注释开头语句 `= 8` | 0（PASS） |

探针打印**编译期看到**的那一行，证明 `ddlV1` 内部确实换了词：
- 改前 `RESULT ddlV1 state line = "  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）"`
- 改后 `RESULT ddlV1 state line = "  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）"`

⇒ **15 → 15 ＝ 结构未动的硬证据**；断言 `n != 15` **一字未动**（§5）。另单跑 `TestTaskState196DDLV4KeepsTheShape`、`TestSchemaContractIntrospection`：`--- PASS` ×2，`rc=0`。

## 5. `git diff -U0` 逐行核 ＋ `--numstat`

```
@@ -68 +68 @@ CREATE TABLE task_log (
-  state          TEXT NOT NULL,          -- 结束态（含 'interrupted'）
+  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）
@@ -447,2 +447,2 @@ func TestTaskState196DDLV4KeepsTheShape(t *testing.T) {
-	// ddlV1 must be untouched: this ticket adds a step, it does not edit the
-	// frozen contract text.
+	// ddlV1's structure and statements must be untouched: this ticket adds a step,
+	// it does not edit the DDL; one wording-only comment was corrected (196-c2/196-c3).

$ git diff --numstat
1	1	internal/memory/schema.go
2	2	internal/memory/taskstate_196_test.go
```

⇒ **2 hunk、3 删 3 增、无整块重写**。

**「全以 `--`／`//` 开头」逐行核 —— 一处尺的射程必须说清**：
- **测试文件 4/4 行**：剥掉 diff 标记后全是 `TAB`＋`//` ⇒ **字面满足**。
- **`schema.go` 那 1 行**是**行内注释**（SQL 列定义行尾的 `--`），剥标记后**以 `  state` 开头**、⛔ 不以 `--` 开头。⇒ 「整行以 `--` 开头」对行内注释是**结构上不可能满足**的尺，故换**更强**的判据（前缀/后缀逐字节）：

```
prefix_before_comment        = '  state          TEXT NOT NULL,          '
prefix_identical             = True    # 「-- 」之前逐字节相同
suffix_after_word_identical  = True    # 「（含 'interrupted'）」逐字节相同
old_only_diff_is_the_word    = True    # 换掉 结束态→任务状态 即得新行，别无他变
new_equals_spec_line_exactly = True    # 新行 == SPEC-02:60 整行，逐字节
old_equals_spec_line_exactly = False   # 改前 != 契约行
```

⇒ SQL 语句／标识符／列定义**一个字节没动**，改动**只在 `-- ` 之后**；目的证死：`new_equals_spec_line_exactly=True` ＝ 代码侧与契约侧**现在逐字节同一行**。

## 6. 门禁 rc（⛔ 未走管道：`cmd > out 2>&1; rc=$?`）

env：只 `unset GOTMPDIR`；**未** `unset TMPDIR/TEMP/TMP`（清了 `os.TempDir()` 变空串，在 `internal/memory` 造假红）。

| 门禁 | rc |
|---|---|
| `go build ./...` | **0**（输出 0 字节） |
| `go test -count=1 -v ./internal/memory/` | **0** |
| `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`） |
| `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`） |

**四数**（尺：`^=== RUN`、`^[[:space:]]*--- (PASS|FAIL|SKIP)`，**允许缩进**故顶层与子测试同尺，四数须自洽）：

| | RUN | PASS | FAIL | SKIP | PASS+FAIL+SKIP | 自洽 | rc |
|---|---|---|---|---|---|---|---|
| 改前 | 91 | 89 | **0** | 2 | 91 | ✅ | 0 |
| 改后 | 91 | 89 | **0** | 2 | 91 | ✅ | 0 |

⇒ 前后**四数完全相同**、`FAIL=0`（顶层 PASS 54、子测试 35）⇒ 停手条件①（改注释致红）**未触发**。

## 7. 格式：甲乙两尺

- **甲（判据）＝`git archive <待提交树>` 导出仓外**（`D:\tmp\…\export`，17640 文件、无 `.git`）：`gofumpt -l internal cmd tools scripts frontend` ⇒ **0 行，rc=0**。（`17640` 与 `check-path-length-budget.sh` 自报 `denominator: tracked paths=17640` **对上**。）
- **乙（只诊断）＝工作树**：同命令 ⇒ **3 行**：`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`、`cmd\wisp\models.go`。

**归属判定用 `[\\/]`**（Windows 上 gofumpt 吐**反斜杠**）：三行**都不含** `memory[\\/](schema\.go|taskstate_196_test\.go)` ⇒ **乙的 3 枚不是我这两个文件**；我这两个文件单独跑 gofumpt：**空列表，rc=0**。
这 3 枚的真身：`git status` 对它们**无改动**；`git ls-files --eol` ＝ `w/crlf`＋`attr/text eol=lf`；CR 计数 131／225／389，导出树同路径 **CR=0** ⇒ **工作树 CRLF 伪影**，非格式漂移、非本腿造成、且**不在我的写面**。

## 8. 前置核查 ＋ 尺的射程自述

**开工前**：`git diff --cached --name-status` ＝ **空**（无本笔不该有的 `D`/`A`）；`git status --porcelain -- internal/memory/ docs/evidence/s1/` ＝ **0 条**（写面内无别人的改动）。落笔后 staged 仅 `M schema.go`、`M taskstate_196_test.go`。

**落笔前五问**：①**量的是哪个对象**＝`15` 量的是 `ddlV1` 这一份 DDL 文本（编译期真值、经 `go test` 包内注入），非我另写的切分器；②**输出形状**＝gofumpt 在 Windows 吐**反斜杠**⇒ 归属用 `[\\/]`，diff 剥标记后**行内注释不以 `--` 开头**⇒ 换前缀/后缀逐字节判据（§5）；③**数的是代码还是注释**＝本腿**只动注释**，`prefix_identical=True` 就是"代码一字节没动"的证明；④**和总数对上没有**＝四数 `91=89+0+2` 自洽、甲尺 `17640 files` 与 `17640 tracked paths` 对上；⑤**谁算的、什么算法、几位**＝语句数由**仓内真实函数**算（非复刻）、格式由 `mvdan.cc/gofumpt`（`D:\work\base\gopath\bin\gofumpt.exe`）算、CR 由 `tr -cd "\r" | wc -c` 算。

**一处方法论错，已当场改正**：首次取 rc 写成 `bash -c "…; echo \"rc=$?\""` ⇒ PowerShell 先把自己的 `$?` 求值，打出 `RC=True`。⇒ rc 一律由**脚本文件内**的 `cmd > out 2>&1; rc=$?` 产出（`gate.sh`），⛔ 不经 PowerShell 插值、⛔ 不走管道。

## 9. 本笔提交号与提交后复跑

- **提交**：`7537627ef5f58d383c0a6176b1450f1cf7002849`（`dev`，⛔ 未 push）
- **提交内容**：3 文件、`+140 / -3`（`A docs/evidence/s1/196-c3-impl.md` 137 行 ＋ `M schema.go` 2 行 ＋ `M taskstate_196_test.go` 4 行）
- **提交后复跑（尺不变）**：
  - 提交里那一枚 blob 的 `:68`：`  state          TEXT NOT NULL,          -- 任务状态（含 'interrupted'）` ⇒ 与 `SPEC-02:60` 逐字节相同；
  - 提交里 `taskstate_196_test.go:449` 仍是 `if n := len(splitSQLStatements(ddlV1)); n != 15 {` ⇒ 断言一字未动、行号未漂；
  - 行数：`schema.go=411`、`taskstate_196_test.go=452`、本件 `=137`（与改前同，未膨胀）；
  - 格式**甲**改用 **`git archive HEAD`** 再导出一次到仓外：`gofumpt -l internal cmd tools scripts frontend` ⇒ **0 行，rc=0**；
  - 写面复跑：`git status --porcelain -- internal/memory/ docs/evidence/s1/` ⇒ **空**。
