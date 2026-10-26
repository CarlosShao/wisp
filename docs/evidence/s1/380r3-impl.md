# 380r3 收尾小修腿 实现件 — R2：表 `:207` 的自述被本笔自己证伪

**写面**：`docs/evidence/s1/c21-native-tokens.md` 只改 `:207` 一句 ＋ 本件。
**未碰**：`docs/PLAN.md` / `internal/**` / `cmd/**` / `frontend/**` / `design/**` / 票面 / 台账 / `Status:` / `-done`。
**索引**：开工前 `git diff --cached --name-status` 空（无预存暂存）；工作树别人的脏（`.scratch/**`、`design/old/`、`design/doubao/` 等未跟踪件）一个未碰。

## R2 判语复核（不采信，自己跑）

- 终裁原文 `docs/evidence/s1/380-v1-verdict.md:71`：表 `:207` 写「`PLAN.md` 3611 行里 `c21-native` **0 处命中**」——同一笔把它证伪，`HEAD` 现 **1 处**（`:3604`，本笔自己写的）。
- 同件 `:106` 给的两个写法：「改成『本笔前 0 处、本笔 `:3604` 起 1 处』或删括号」⇒ 本腿取前者，且**只改这一句**。
- 出处核对：`git show --numstat 0a8920a15 -- docs/PLAN.md` ⇒ `2 2`；该笔两对 `+/-` 正是 `:3308` 与 `:3604` ⇒ `:3604` 那个 `c21-native` 确由 `380r2` 写入，不是更早的腿。

## 八项读数

### ① 改前现量 ⇒ **1 处**

`grep -c 'c21-native' docs/PLAN.md` ⇒ `1`（rc 0）；命中行 `:3604`；`wc -l docs/PLAN.md` ⇒ `3611`。
★ 现量 ＝ 1，与终裁的量一致 ⇒ 不触发「不是 1 处就停手上报」。
句子两半也各自核过：`git show 0a8920a15^:docs/PLAN.md | grep -c 'c21-native'` ⇒ `0`（rc 1，零命中）；`git show 0a8920a15:docs/PLAN.md | grep -c` ⇒ `1`。

### ② `:207` 改前原文（逐字）

> **这两节不是 `C1–C32` 文本**（`PLAN.md` 3611 行里 `c21-native` 0 处命中）：它们是 `C21` 契约的**书面镜像**，

### ③ 改后原文（逐字）

> **这两节不是 `C1–C32` 文本**（`PLAN.md` 3611 行里 `c21-native`：票 380 之前 0 处，票 380 在 `:3604` 起 **1 处**）：它们是 `C21` 契约的**书面镜像**，

★ 三个数对上：`grep -c` 改后现值 **1** ＝ 句中写的「`:3604` 起 **1 处**」；`wc -l` 现值 **3611** ＝ 句中写的「3611 行」；「票 380 之前 0 处」＝ `0a8920a15^` 实测 `0`。

### ④ 78 行数据零改动（`-` 行数恰等于我改的那一句）

- `git diff --numstat -- docs/evidence/s1/c21-native-tokens.md` ⇒ `1` 删 / `1` 增（提交后 `git show --numstat` 同值：同一对 blob 的同一个 diff）。
- 唯一 hunk `@@ -204,7 +204,7 @@`：`-` 恰 1 条（`:207` 那一句）、`+` 恰 1 条。
- `wc -l` 改前 `371` / 改后 `371` ⇒ 只换句、不插删行。
- ⇒ `-` 行数 **1** ＝ 我改的那一句的行数；表里任何一行都没进 diff。

### ⑤ 78 行 ＋ 两张对照表（★已排表头）

| 块 | 行范围 | `^\| ` 计数 | 减表头 1 | 期望 |
|---|---|---|---|---|
| `## 明色系` | 46–90 | 41 | **40** | 40 ✅ |
| `## 亮色系` | 91–136 | 39 | **38** | 38 ✅ |
| 第四代槽位对照 | 137–173 | 12 | **11** | 11 ✅ |
| 第四代未覆盖 | 174–209 | 24 | **23** | 23 ✅ |

- 78 ＝ 40 ＋ 38；34 ＝ 11 ＋ 23。改前改后四个块计数**一模一样**。
- ★表头陷阱（本会话已踩）：分隔行 `|---|---|---|---|` 以 `|-` 开头、**不进** `^\| ` 计数；但**表头**行形如 `| CSS 变量 | … |`、会被数进去 ⇒ 每块必须 **−1**。本件所有数都按「排表头」口径报。

### ⑥ 别处同类「0 处」自述清单（★只登记，⛔ 未改）

| 位置 | 原文（节选） | 现在假不假 | 处置 |
|---|---|---|---|
| `c21-native-tokens.md:325` | `Palette.TintSuccessHi` …「`--tint-success-hi` 全仓无引用（`design/` 实测 0 处 `var()`）」 | **不假**：现跑 `grep -rn "tint-success-hi"` 只命中**声明**与 Go 注释/测试注释，`var(--tint-success-hi)` 消费 **0 处** | 同形状但成立 ⇒ 按硬边界③ **不动** |

- 全文件 `c21-native` 只出现在 `:207` 一处；`PLAN.md` 字样只在 `:207`／`:208`；`3611` 只在 `:207`。
- 除 `:207`（已修）外，本件内**无**第二处「0 处命中／0 命中／零命中／无命中」类计数自述。

### ⑦ 门禁 rc

- `go test -run TestC21 ./internal/ball` ⇒ **rc 0**（`ok github.com/CarlosShao/wisp/internal/ball 0.169s`）
- `go test -run TestC21DesignTokensFourWayAgree ./internal/panel` ⇒ **rc 0**（`ok github.com/CarlosShao/wisp/internal/panel 0.047s`）
- `./tools/d22scan/d22scan.exe`（仓根跑）⇒ **rc 0** clean（`examined 294 production Go files under internal/ and cmd/`）。
  ★射程说明：ban#8 扫的是 `design/ 39` ＋ `frontend/ 116` ＋ `internal/ 652` ＋ `cmd/ 153`，**不含 `docs/`** ⇒ 本腿这处 docs 改动**不被它覆盖**，⛔ 不得拿它的 rc 0 当本文件的证据；也按派单**没有**拿 ban#8 去扫 `docs/`。

### ⑧ 格式甲乙两尺（★分开报）

- **甲尺** `gofmt -l internal cmd` ⇒ 3 条：`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`cmd/wisp/models.go`。
- **乙尺** `/d/work/base/gopath/bin/gofumpt.exe -l internal cmd` ⇒ 同 3 条（两尺点名集合相同）。
- ★这 3 条**是既存脏、与本题无关**，两条证法：
  1. 归属：`git diff HEAD -- <这 3 个>` **空** ⇒ 工作树字节 ＝ `HEAD` 字节，不是本腿改的（本腿提交 **0 个 `.go` 文件**）。
  2. 真因＝**CRLF 不是代码形**：三文件工作树 `LF == CRLF`（225／389／131 行**全 CRLF**）；而 `git cat-file -p HEAD:<f>` 的 blob 是 LF、拷到临时件再跑 `gofmt -l` **不被点名**。`core.autocrlf=true`（`.gitattributes` 反而写 `eol: lf`）⇒ Windows checkout 产物。
- ⇒ 只看本腿改动面（0 个 `.go`），甲乙两尺**皆空**；上面 3 条建议另开票（不在本腿写面）。

## 停手判据复核

- 现量 **1 处** ⇒ 未触发「不是 1 处就停手上报」。
- 改动未超过那一句（`numstat` `-` 恰 **1**）⇒ 未触发「必须超过那一句就停手上报」。
- 未改：`PLAN.md`、`internal/**`、`cmd/**`、`frontend/**`、`design/**`、票名／`Status:`／`-done`／台账；⛔ push。
