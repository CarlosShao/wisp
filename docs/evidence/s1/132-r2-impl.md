# 票 132 收尾腿 `132-r2` —— 两处"文字与事实不符"的收口（**⛔ 产码**）

- **腿**：`132-r2` ＝ **写腿（实现者）**。⛔ 翻框／⛔ 改票名／`Status:`／`-done`／⛔ 改台账／⛔ push／⛔ 动产码行为。
- **起手锚**：`21231ee05`；读数全部**本腿现跑**，尺逐条列在表里。
- **落笔 4 枚（都在允许面内）**：`cmd/wisp/logsink.go`（只那一段注释）、`docs/evidence/s1/132-r1-impl.md`／
  `132-r1-impl-part2.md`（只追加更正）、本件。`internal/**`、`scripts/**`、`.github/**`、`tools/d22scan/**`、
  `go.mod`/`go.sum`、`docs/PLAN.md`／`docs/specs/**`／`docs/reports/**`、`frontend/**`／`design/**`、
  `thresholds.go`／golden／`allowlist.txt` **一字未动**。

## §0 开工前两连：第一连**先红后转绿**，照实记

| # | 检查 | 读数 |
|---|---|---|
| 1 | `git diff --cached --name-status`（开工瞬间） | **1 枚不属于本笔的 `D`**：`130-…-silent.md`（票 130 旧名，`79 deletions`）⇒ 按派单**停手上报**、未落笔 |
| 1b | 同上（复读） | **空** —— 编排者以 `0eccd6435`（"改名删除半边"）补掉；索引干净后才落笔 |
| 2 | `git status --porcelain --` 写面三枚 | **空**（写面内无别人的未提交改动） |
| 3 | 行号／逐字引文 | **全部现跑**。派单给的 `:44-46` 实测 **`:44-50`**（§1）；`privateDirAll` 两行号**经核成立**（§2） |

停手那一步的原始读数留 `/d/tmp/wisp/20261023/132r2/01-baseline-index.txt`（留痕，不做判据）。

## §1 件一：`cmd/wisp/logsink.go` 的过期注释（真行号 **44–50**）

**改前逐字**（7 行；`:45-46` 正是 "nothing here seals or un-seals the log file"）：

```go
// INTERIM (ticket 117 AC#3): whether a log file counts as private data is the
// one question of ticket 95's that owner has not answered, so nothing here
// seals or un-seals the log file - it neither adds nor removes that decision.
// What this file does guarantee is *placement*: the sink is created only under
// the resolved data dir of the running env, so the question owner has to answer
// is about one directory, not about a temp file. See logSinkDir and its
// invariant test.
```

**它今天为什么站不住**（三条，全现读）：

| 它写的 | 今天的事实 | 尺 |
|---|---|---|
| "日志算不算私有数据是 owner 没答的那一问" | **已答**：`Q-31` ＝「**算、要锁**」⇒ 由此新建票 132 | `docs/reports/pending-and-issues.md:3463` |
| 弦外之音"没人封" | **管线真封目录**：`InitLog` → `InitLogWithRegistry` → `winsec.PrivateDirAll(cfg.Dir, 0o700)` | `internal/observe/logging.go:87`／`:99` |
| 同上 | **管线真封每个滚动文件**：`sealLogFile`(＝`winsec.SealFile`)，create 之后、首个内容字节之前 | `internal/observe/logging.go:223`／`:316` |

**改后逐字**（仍 7 行）：

```go
// SEALING, formerly INTERIM (ticket 117 AC#3): the question it deferred is
// answered - owner ruled under Q-31 that a log file counts as private data
// and is to be locked - and ticket 132 landed that lock one layer down, in
// internal/observe: InitLog creates the directory through winsec.PrivateDirAll
// and sealLogFile seals every rolled file. So nothing here seals or un-seals
// the log file: what this file decides is *placement*, under the resolved data
// dir, plus the mirrored stderr copy. See logSinkDir and its invariant test.
```

**为什么这版能站住**：①`INTERIM` **不成立就换** —— interim 期已结束（owner 答了、封落地了）⇒ 改 `SEALING, formerly INTERIM`；
旧词留痕有意为之：`grep -n INTERIM cmd/wisp/logsink.go` 仍命中 **1** 行，引"那段 INTERIM 注释"的旧件不悬空。
②"nothing here seals or un-seals the log file" 这半句**是真的、保留**：`grep -n "Seal\|seal" cmd/wisp/logsink.go` 的
6 处命中**全在注释里**（`:7 :29 :31 :41 :48 :169`）⇒ 本文件**零 seal 调用**；改的是**理由**（不再是"没人知道该不该锁"，
而是"封在下一层，本文件只管落点与 stderr 镜像"）。③**7 行→7 行**，`wc -l` 前后都是 **224** ⇒
`logsink.go:71/:82/:149/:153` 等被别处逐字引用的行号**一个都没漂**（本仓最贵的一类腐坏就是行号漂）。

## §2 件二：实现腿件 §2 那句假话（"封条在层出现之前就位"）

**改前逐字**（两处，同一个错法）：

- `132-r1-impl.md:42`：

  > | `:99` | `os.MkdirAll(cfg.Dir, 0o755)` | `winsec.PrivateDirAll(cfg.Dir, 0o700)` | 封条**在层出现之前**就位（父先于子），且既有树走修复分支 |

- `132-r1-impl-part2.md:53-54`：

  > 1. **`SealDir`（导出函数）生产调用者仍是 0**。落地用的是 `winsec.PrivateDirAll`——它在**创建每一层之前**就封
  >    （`MkdirAll`+`SealDir` 做不到，正是 AC#2「创建时就」的语义），且 `privateDirAll` 内部调的就是同一枚 `sealDir` 原语。

**错在哪**：没有层就没有对象可封 —— "层出现之前就封"断言的是**先封后建**，而 `privateDirAll` 是**先建后封**。现跑（`HEAD` 上现读）：

| 位置（现读） | 逐字 |
|---|---|
| `internal/winsec/winsec.go:283` | `if err := os.Mkdir(missing[i], perm); err != nil {` |
| `internal/winsec/winsec.go:292` | `if err := sealError(missing[i], sealDir(missing[i])); err != nil {` |
| `internal/winsec/winsec.go:282` | `for i := len(missing) - 1; i >= 0; i-- {` ⇒ 从最外层往里走，父先于子 |

**导出文档里的强保证，逐字**（`internal/winsec/winsec.go:223-229`）：

> PrivateDirAll creates path and any missing parents, sealing **every level it creates** before descending into it.
> Sealing afterwards is not enough: a child inherits the ACL it had at creation time, so a directory created wide and
> tightened later leaves everything written during that window wide. The order this guarantees is "sealed parent, then
> any child", which is what makes one seal cover files nobody seals individually - SQLite's -wal/-shm and a
> pre-rename retry temp.

**正确的说法**：**创建与封同一次调用、父先于子**（⛔ "层出现之前封"）。差别不是措辞洁癖：前者断言一件做不到的事
（对还不存在的对象设 ACL），后者才是覆盖 `-wal`／`-shm` 的真保证 —— 封落在**任何子层／子文件被写之前**。

**改后逐字**（两件都只**追加**、⛔ amend；各新增一节 `## 更正（票 132 收尾腿 132-r2，2026-10-23）`）：

- `132-r1-impl.md` 追加节的判决句：**创建与封在同一个循环体、同一次调用里完成，父先于子**。
- `132-r1-impl-part2.md` 追加节的判决句："与 `132-r1-impl.md` §2 那处是**同一个错法**。`privateDirAll` 是**先建后封**……"，
  并如实标注"§9 第 2 条所请的那一步已由本腿做完"（免得下一位以为 `logsink.go` 那段注释还旧）。
  两节都写明**保留**：`SealDir` 生产调用者仍 0、内部仍走同一枚 `sealDir` 原语、"按能力而非函数名读"仍成立
  —— 只改**能力的语义**，⛔ 借更正连带削弱已成立的部分。

## §3 门禁（rc 一律 `cmd > out 2>&1; rc=$?`，⛔ 走管道）

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 日志 0 字节 |
| `go vet ./cmd/wisp/` | **0** | — |
| `sh scripts/d22scan.sh` | **0** | `clean`；`bans #1-5 internal/=248 cmd/=42`；`ban #8 internal/=631 cmd/=148`（与 `132-v1-verdict.md:100` 逐字同 ⇒ 没动被计的东西） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；`tracked paths=17638  over-budget=61  covered by roster=61  not in roster=0` |
| `gofumpt -l cmd/wisp`（绝对路径） | 0 | **1 行 `cmd\wisp\models.go`** —— ⛔ 不是本笔，见 §5 |
| `gofmt -l cmd/wisp/logsink.go` | 0 | **0 行** |
| CRLF（`tr -cd "\r" \| wc -c`） | — | 落笔 4 枚**全 0** |

## §4 尺：怎么证"只改注释与文字、⛔ 动执行行"

```
git diff -U0 HEAD -- cmd/wisp/logsink.go | grep -v '^[+-][[:space:]]*//' | grep -v '^+++\|^---'
```

**读数**：只剩文件头／区块头五行（`diff --git`／`index`／`--- a/`／`+++ b/`／`@@ -44,7 +44,7 @@`）⇒ **执行行 0**。
同一把尺的计数式（尾接 `grep -c '^[+-]'`）⇒ **`0`**。这就是"没动执行行"的读数，不是"我看了觉得没问题"。

**尺的射程**（`ruler-scope-checklist.md` 第 3 条）：不过滤 `//` 时这把尺得 **14**（7 删＋7 增）—— 那 14 行是**注释**、
⛔ 执行行；两个数都写在这，免得下一位把 14 读成"改了 14 行代码"。

## §5 观察（⛔ 本笔地界，只记不修）

`cmd/wisp/models.go` 在 `HEAD` 上是 **CRLF**，故 `gofumpt -l cmd/wisp` 唯一命中它（反斜杠形状，判归属用 `[\\/]`）。
①`git status --porcelain -- cmd/wisp/models.go` **空**（与 `HEAD` 逐字节相同，不是我改的）；②`tr -cd "\r" | wc -c` **389**、
`wc -l` **389** ⇒ **每行都是 CRLF**，而 `.gitattributes:4` 写 `*.go text eol=lf`；③引入它的提交 `aeea4e5b7`（票 357 落地腿）。
⇒ `gofumpt -d` 要重写**整个 389 行**（全是 CRLF→LF）。**⛔ 本腿不修**：不在写面，且"顺手改别的"被明令禁止 ⇒ 交编排者另派。

## §6 没做什么

1. **"中途 roll 封失败 `slog` 吞错 ⇒ 产品侧无一行 seal 失败"** —— 票 132 的第三条残余，**要另开票**，本腿 ⛔ 办。
2. 上列写面外的一切**一字未动**（`internal/**`、`scripts/**`、`.github/**`、`docs/reports/**`、`thresholds.go` 等）。
3. ⛔ 翻框／改票名／`Status:`／`-done`／改台账／push；⛔ `add -A`／`.`；⛔ amend／reset／rebase／stash／`checkout .`／`clean`。

## §7 提交

主题逐字以 `票 132 收尾腿 132-r2` 开头；`add`／`commit` **两处都显式 pathspec、逐枚文件字面点名**（⛔ glob）。
本件与提交同笔，故提交 SHA 见 `git log`（⛔ 把 SHA 写进它自己那笔）。
