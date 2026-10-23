# 票 130 收尾腿 `130-r3` —— 只换 `:183` 那句 `Fatalf` 文案：把"132 reads the file seal exactly there"换成站得住的读数

- **写腿（实现者）**。**写面**：`internal/observe/earlylog_takeover_130r1_test.go`（**只那段 `Fatalf` 文案字符串**）＋本件。
  零码路径：`git diff -U0` 里 **`if`／断言条件 0 处**、**非字符串字面量改动 0 行**（`+5 / -2`，全部是字符串片段行）。
- **起手锚** `87036ca95`；本笔期间编排者落地 `0eccd6435`（`A1182` 补丁·改名删除半边）＋`21231ee05`（台账 `A1183`）⇒ **终锚 `21231ee05`**。
  两锚之间我的文件 blob **逐字相同**（`4787b3b3c5eac23b22d435105d5272273c90142e`）⇒ 改前读数不受影响。
- **尺的声明**：`RUN/PASS/FAIL/SKIP` 四数＝`go test -v` 里 `^=== RUN   Test`／`^--- PASS: Test`／`^--- FAIL: Test`／`^--- SKIP: Test` 的**行数**（`tr -d '\r'` 后数）；rc 一律 `cmd > log 2>&1; echo $?`（**⛔ 走管道**，避免把红读成绿）。
- **环境**：`GOFLAGS=`、**只 `unset GOTMPDIR`**（`TMPDIR`/`TEMP`/`TMP` **未动**）、未设 `CGO_ENABLED`；Go `go1.27.1 windows/amd64`。
- 全部临时件在 `D:\tmp\wisp\20261023\130r3\**`（**只建不删**）；日志 `…\logs\01…15-*.txt`；突变副本 `…\mut\*`；overlay `…\overlay-m2*.json`。

## ① 开工前第 1 道闸：索引里有一条不该有的 `D`（我先停手上报，编排者已清）

| 项 | 我现跑的读数 |
|---|---|
| 索引 | `git diff --cached --name-status` ⇒ `D  .scratch/wisp/issues/130-…-silent.md`（**唯一一条**，79 行全删） |
| 它在 HEAD 里？ | `git cat-file -e HEAD:<原票>` ⇒ **IN-HEAD**；`test -e <原票>` ⇒ **NOT-ON-DISK**；`git ls-files -s` ⇒ **索引无条目** |
| 来源 | `HEAD`=`87036ca95 A1182`，该笔 stat **只有两枚**（`…-done.md 104+`／台账 `49+`）⇒ 原票的删除**被 stage 了但没 commit** |
| 判 | **不在我的写面**、不是我产生的 ⇒ 按派单「出现本笔不该有的 `D`/`A` ⇒ 停手上报」**停手交读数**（未动任何仓内文件） |
| 落定 | 编排者 `0eccd6435`（pathspec 逐字点名旧名＋新名）⇒ 我**现跑复核** `git diff --cached --name-status` **空**、`HEAD` 里 `issues/130-` 只剩 `-done` 一枚 ⇒ 闸 2（写面内有别人的未提交改动）也**空** |

## ② 改前逐字（`git show 87036ca95:internal/observe/earlylog_takeover_130r1_test.go`，行 `182–186`）

```go
	if got := fileBytes130r1(t, dir); got != 0 {
		t.Fatalf("after InitLog (and a flush) the boot file holds %d bytes, want 0: construction is not the "+
			"handover instant - ticket 132 reads the file seal exactly there, and a replay this early makes "+
			"that reading unreadable", got)
	}
```

**为什么是假话**：票 132 的用例**不 flush**、量的是 `bufio`（`logprivate_132_test.go:64-70`：`os.Stat` ＋ `info.Size() != 0`）⇒ 把补发挪进 `InitLogWithRegistry`（M-2）后，132 六枚**仍 6 PASS / 0 FAIL rc=0**（§④ 我自己造过），所以"132 reads the file seal exactly there … makes that reading unreadable"**不成立**。

## ③ 改后逐字（工作树 `182–189`）

```go
	if got := fileBytes130r1(t, dir); got != 0 {
		t.Fatalf("after InitLog (and a flush) the boot file holds %d bytes, want 0: construction is not the "+
			"handover instant - a replay moved here is caught by this case's own flushed reading, not by "+
			"ticket 132's six, which stay green under the same move (6 PASS / 0 FAIL, "+
			"docs/evidence/s1/130-v1-verdict.md section 6): the one case there that reads the boot file's "+
			"size at construction stats it without flushing, so it reads the writer's bufio rather than the "+
			"handover, and the other five never take that reading", got)
	}
```

新文案里每一句都**可复跑**：真代价＝本笔 `:182`（带 flush）变红；132 六枚 M-2 下仍绿＝§④ 的读数；"只有一枚读构造瞬间的文件大小"＝`grep -n "Size()\|os.Stat" internal/observe/logprivate_132*_test.go internal/observe/sealcalls_132_test.go` 只命中 `:34`（目录）、`:64/:68`（构造瞬间，**唯一**）、`:83`（**先 FlushNow** 的另一个读数）、`:134-139`（封印失败支）。

## ④ `130-v1` §6 我自己复核：132 六枚在 M-2 下仍绿

- **M-2 造法**（仓外 overlay，仓内**零 `logging.go` 字节**）：`internal/observe/logging.go` 的 `p.flushDone = h.Done()` 后插一行 `_, _ = earlyLogs.drain(handler)`（即"把补发挪进 `InitLogWithRegistry`"）。
  `diff -u` 只多这一行；`go test -overlay=…/overlay-m2.json`。
- **六枚名单**：`TestLogDirIsSealedAtCreation132r1`／`TestLogFileIsSealedBeforeTheFirstContentByte132r1`／`TestExistingWideLogTreeIsRepaired132r1`／`TestSealFailureRefusesTheWrite132r1`／`TestUnsealableLogDirIsRefusedLoudly132r1`／`TestSealEntryPointsHaveProductionCallers132r1`。

| 树 | RUN | PASS | FAIL | SKIP | rc | 日志 |
|---|---|---|---|---|---|---|
| 原样（改前） | 6 | 6 | 0 | 0 | 0 | `02-before-six132.txt` |
| **M-2** | **6** | **6** | **0** | **0** | **0** | `03-M2-six132.txt` |
| 原样（改后，文案已换） | 6 | 6 | 0 | 0 | 0 | `07-after-six132.txt` |

⇒ **`130-v1` §6 那条我复现成立**：M-2 下 132 六枚**仍绿**，"塞进 `InitLog` 会花掉 132 的封印判据"是假话。机制＝`info.Size()` 读的是未 flush 的 `bufio`（132 那枚**不调用** `FlushNow`）。

## ⑤ 改前／改后四名逐数相同（同一 `-run` 正则、同一 `-count=1`）

- 正则：`^(TestInstallAsDefaultTakesOverTheBootRecord130R1|TestHandoverBelongsToTheFirstListener130R1|TestInstallThenFlushDeliversOnce130R1|TestHandoverKeepsTheRedaction130R1)$`

| 时点 | RUN | PASS | FAIL | SKIP | rc | 日志 |
|---|---|---|---|---|---|---|
| **改前** | **4** | **4** | **0** | **0** | **0** | `01-before-four.txt` |
| **改后** | **4** | **4** | **0** | **0** | **0** | `06-after-four.txt` |

四名逐枚（改后）：`TestInstallAsDefaultTakesOverTheBootRecord130R1` PASS／`TestHandoverBelongsToTheFirstListener130R1` PASS／`TestInstallThenFlushDeliversOnce130R1` PASS／`TestHandoverKeepsTheRedaction130R1` PASS —— **逐数与改前相同**（只是"红的时候说的话"变了，见 §⑥ 的红行）。

## ⑥ `:182` 的 `FlushNow()` 换成 `error(nil)` ⇒ 头号用例 **PASS**（我自己复现，逐字读数）

突变（只有一处，verdict 已验过这条路）：`if err := p.FlushNow(); err != nil {` → `if err := error(nil); err != nil {`（`diff -u` 只这一行）。

| 组合 | RUN | PASS | FAIL | SKIP | rc | 日志 |
|---|---|---|---|---|---|---|
| M-2 ＋ **去 flush**（改前文案） | 4 | 4 | 0 | 0 | **0** | `05-M2-noflush-four.txt` |
| M-2 ＋ **去 flush**（改后文案） | 4 | 4 | 0 | 0 | **0** | `09-after-M2-noflush.txt` |
| **对照**：M-2 ＋ 保 flush（改后文案） | 4 | 3 | 1 | 0 | **1** | `08-after-M2-flush.txt` |

对照那一发的红行**逐字**（证明"真代价是本笔自己那条带 flush 的钉子"，且 158 字节是**现跑**的）：

```
    earlylog_takeover_130r1_test.go:183: after InitLog (and a flush) the boot file holds 158 bytes, want 0: construction is not the handover instant - a replay moved here is caught by this case's own flushed reading, not by ticket 132's six, which stay green under the same move (6 PASS / 0 FAIL, docs/evidence/s1/130-v1-verdict.md section 6): the one case there that reads the boot file's size at construction stats it without flushing, so it reads the writer's bufio rather than the handover, and the other five never take that reading
```

⇒ 带 flush ⇒ **红**（158 字节）；去掉 flush ⇒ **绿**（4/4 rc=0）⇒ `:182` 才是那条钉子，与 132 无关。

## ⑦ `git diff -U0` 的形状（⛔ 改断言的证据）

- hunk 头：`@@ -184,2 +184,5 @@ func TestInstallAsDefaultTakesOverTheBootRecord130R1(t *testing.T) {`（**只有这一处**）；`--numstat` ＝ `5  2`。
- 变更行里**非字符串字面量**行数：`grep -E '^[+-]' … | grep -vE '^[+-][[:space:]]*"' | wc -l` ⇒ **0**。
- 变更行里含 `if`／`for`／`!=`／`==`／`Fatalf(`／`Errorf(` 的行数 ⇒ **0**。
- `if got := fileBytes130r1(t, dir); got != 0 {` **未在 diff 里出现**；`FlushNow` 调用、`error(nil)` 都**没有**进仓（突变全在 `D:\tmp`）。

## ⑧ 门禁 rc（现跑，⛔ 走管道）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| `go build` | `go build ./...` | **0** | 空输出（`10-build.txt`） |
| 包测试 | `go test -count=2 ./internal/observe/` | **0** | `ok … 9.274s`（`11-observe-count2.txt`） |
| ban #8 | `sh scripts/d22scan.sh` | **0** | `clean`；`internal/` **631** Go 文件（含 `_test.go`、注释也扫）／`cmd/` **148**；计数 `248/42/116/25/39/116/631/148`（`12-d22scan.txt`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0** | **VERDICT GREEN**；longest=185、worst=229、hat=165、over=61、wall interval=**0**、roster=61（`13-pathlen.txt`） |
| gofumpt（**绝对路径**） | `/d/work/base/gopath/bin/gofumpt.exe -l internal/observe/earlylog_takeover_130r1_test.go` | **0** | **0 行**（`14-gofumpt-mine.txt`） |
| gofumpt 全树（诊断） | 同上 `-l . tools/d22scan tools/mockllm` | 2 | **54** 行＝**51 `.scratch\`**（含故意写坏的 `p_unp.go`）＋**3 eol 伪红**（`internal\agent\approval\pending_read.go` CR=131／`internal\agent\tools.go` CR=225／`cmd\wisp\models.go` CR=389，`git ls-files --eol` 三枚都是 `i/lf w/crlf`）；**本笔文件 0 行**（`15-gofumpt-full.txt`） |
| CRLF（本笔文件） | `tr -cd "\r" < <file> \| wc -c` | — | **0**；`git ls-files --eol`＝`i/lf w/lf` |

## ⑨ 声明与残余

1. **只换了字符串**：checkout 里 `git show HEAD:<file>` 与工作树的差异**逐字**＝§②→§③；`if`、量什么、`Fatalf` 的条件**一字未动**（§⑦ 三读数）。
2. **`ban #8` 对本件的射程**：新文案**全 ASCII**（无 `1F000–1FAFF`／`2200–22FF`／`2600–27BF`／`2B00–2BFF`／`FE0F`／`1F1E6–1F1FF` 带内码位）；`d22scan` rc=0 覆盖 `internal/` 631 个 Go 文件（**字符串不豁免**）。
3. **未碰**：`logging.go`（`130-r2` 的注释）、`cmd/**`、其余 `internal/**`、`scripts/**`、`.github/**`、`tools/d22scan/**`、`go.mod`/`go.sum`、`docs/PLAN.md`／`docs/specs/**`／`docs/reports/**`、`frontend/**`／`design/**`、`thresholds.go`／golden／`allowlist.txt`／三枚冻结件。**⛔ push**、⛔ amend/reset/rebase/stash/checkout。
4. **一条射程声明**：§④ 的 6/6 与 §⑥ 的 4/4 都是**本机 Windows（`go1.27.1 windows/amd64`）单发**读数，与 `130-v1` §6 同向；`cmd/wisp` 整包**未跑**（既有环境上限，⛔ 当判据）。
5. **派单未定义即停项①**（改文案会让别的用例的 `strings.Contains` 失效）**不成立**：全仓 grep `reads the file seal exactly there`／`makes that reading unreadable`／`construction is not the handover instant` ⇒ `.go` 里**只有本文件这一处**（其余命中全在 `docs/`、`issues/` 的**叙述**里，不是断言）。
