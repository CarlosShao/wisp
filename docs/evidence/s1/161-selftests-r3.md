# 161-r3 证据件：把自检装成规矩（AC#2）＋ 止血（AC#7①）＋ 反向判据（AC#4）

- 程：票 161 的 **r3 写码位**（派单 `.scratch/wisp/dispatches/2026-09-26-221x-impl-161-r3-selftests-and-stop-the-bleeding.md`）
- 射程：**AC#7①（止血）＋ AC#2（自检装成规矩）＋ AC#4（反向判据）**。**AC#3 已归兄弟程 161-r2**（`docs/evidence/s1/161-doorbell-census-r1.md`，本件不重做、不复算）；**AC#6 归后一程**。
- 本程写面：`.scratch/wisp/probes/161/r3/**` ＋ 本件 ＋ `tools/d22scan/**` 的**测试面与自检入口** ＋ `.github/workflows/ci.yml`（**只加步骤**） ＋ `.scratch/wisp/probes/158/accept-r1/mut/guard.no{1,2,3}.go`（**仅空白格式化**） ＋ 票 161 Progress log（只追加）。
- 档位约定：本件默认档位＝**〔我本轮现跑过〕**。凡是抄派单或别人的数的，都单独标出。

---

## 0. step-0 五件（命令原文 ＋ 输出）

### 0.1 `date`

```
$ date
Sat Sep 26 22:43:51 CST 2026
```

（后续每次落时间前重新取过，见各节。）

### 0.2 锚点

```
$ git rev-parse HEAD
27f798ec1ebb5fa434c1ecc50b2b0cd815ee3242
```

⇒ 本程一切"改前/改后"的锚＝**`27f798e`**。派单明令不写死别人的号，所以 §3 的"入库件"归属是用 `git log --no-walk`／`git diff-tree` 现查的，不是抄的。

### 0.3 有没有别人未提交的东西躺在我要写的文件上

```
$ git status --porcelain -- .scratch/wisp/issues/160-*.md .scratch/wisp/issues/161-*.md docs/reports/pending-and-issues.md
(空)
```

⇒ **票 161 票面此刻是干净的**（编排者 21:2x 那 53 行在途编辑已经在库里了；161-r1 §5.7 请求的处置已经发生）。这是本程敢在最后落票面追加的前提。

全树现场（`git status --porcelain | wc -l` = **43**）里不是我写面的部分，逐条对上派单 §3：

```
 M .scratch/wisp/probes/152/my152.py                                    <- 停下来程序的半件，不动
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md     <- 同上
 D design/** （16 枚）+ ?? design/** 新件                                <- owner 自己的未提交删除，不动
?? .scratch/wisp/probes/161/r2/**                                        <- 兄弟程 161-r2 的写面，不碰
?? .scratch/wisp/probes/139/accept-r1/ .scratch/wisp/probes/156/... .zcodeignore 等
```

⇒ 每次提交都带**显式文件级 pathspec**（不是目录），`git add -A` / `git commit -a` 一次都没用。

### 0.4 票面 AC#7／AC#2／AC#4 三格（连行号抄；行号是读数）

票面 = `.scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md`

- **AC#2**（票面 **21–26 行**）
  - `21` 判据正文：「**AC#2 把自检装成规矩**：每道门自带 `--self-test`（干净样本必须不响、违规样本必须响，**两向都要过**），并有一枚 CI 入口把这些自检**全部跑一遍**；任何一枚自检挂了＝CI 红。⚠ 本票**允许**动 `tools/d22scan/**` 的**测试面与自检入口**，但**禁止改任何禁令的射程/波段**（那是 `Q-46`/票 141 已结案的契约面）。」
  - `22–24` CI 写面边界（09-26 补）：只许在 `lint` 作业里**新增**步骤（不带 `if:`、不带 `continue-on-error`、不取消任何一步）；不许动触发表／并发组／`slo-full`，不许挪位或删掉任何一步——要动＝契约级停手。
  - `25–26` 改前基线那两枚红的口径：「除那两枚之外零枚新增红」＋名册差集，不许放宽断言、不许 `t.Skip`、不许改 `frontend/**`。**⇒ 这两枚今天不红（§1 表末行现量），口径落地为"零枚红"。**
- **AC#4**（票面 **28 行**，一整行）：「**AC#4 反向判据（承重两问）**：把 AC#2 那对样本里"违规"那一枚摘掉，CI 会不会红？答不出＝装饰。」
- **AC#7**（票面 **35–43 行**）
  - `35` 标题：「**AC#7（09-26 22:0x 追加，来路＝161-r1 报回＋编排者独立复算）自家取证台件不许把全仓门点红，且"跑门"必须排在"落台件"之后**。」
  - `36–37` 现场：`gofumpt -l . tools/d22scan tools/mockllm`（版本现读 `v0.12.0 (go1.27.1)`）今天不空，**3 行是已入库件** `probes/158/accept-r1/mut/guard.no1.go`／`no2.go`／`no3.go`（`506cbae`，20:49，票 158 验收程自己加的）。
  - `38–39` 时序：验收表 **20:2x** 现量过 `rc=0 空` 并据此判 AC#5③ 成立，不到半小时后**它自己加的台件让那条读数过期**——「这不是谁读错，是"门先跑、台件后落"这个顺序本身没有门」。
  - `40` 本格两件事：① **止血**＝格式化那 3 枚 ＋ 复跑 `-overlay` 证明变异行为一字未变（"只改空白不改语义"这句话**要用读数证明，不许说**）；② **装顺序**＝全仓级那三把尺必须在台件全部入库之后跑最后一次，并写成**可复算的形状**，且必须**说得出删掉它哪一发会重新漏**。
  - `42` 判据红线：「不许写成"台件永远不许被这些门点红"——那是恒真检。」要的是**"门点红了有人能在推送之前知道是自己的台件造的"**这一发可复算。

### 0.5 前一程成果（当输入，不重做）

`docs/evidence/s1/161-gate-blindspot-r1.md`：§3 那 27 发成对样本 ＋ §4 六枚盲区（B1–B6）＋ §4.1 三枚脾气（A1–A3）。
⇒ §3 的 `--self-test` 样本表基本就是那张表＋那六枚盲区的**形状**；本程把它们搬进程序（见 §4）。

---

## 1. 前提核对：派单点名的每一条，逐枚现量

| 派单/票面写死的前提 | 盘上读数（我本轮现跑） | 结论 |
|---|---|---|
| `gofumpt -l . tools/d22scan tools/mockllm` **今天不空** | **不空**：rc=**2**，**8 行** | **成立** |
| "3 行是已入库件 `guard.no{1,2,3}.go`" | 成立，但**不完整**：另外 **5 行也是真的**，全部是 `probes/161/r1/runs/**`——161-r1 自己 `gitignore` 掉的台件。`git ls-files` 证明 5 行**未入库**、3 行**已入库** | **成立但只是半张表**，见 §1.1 |
| `506cbae`，20:49 加的这 3 枚 | `git log --no-walk` → `506cbae7 2026-09-26 20:49`；`git diff-tree --name-status` → 三枚都是 **A**（新增） | **成立** |
| 验收表**20:2x** 量的 `gofumpt=空` | 那三枚读数文件的**盘上 mtime 是 20:31**：`gofumpt-full.txt` 0 字节、`gofumpt-narrow.txt` 0 字节、`gofumpt-version.txt` = `v0.12.0 (go1.27.1)`。验收件 §1.3（`docs/evidence/s1/158-...accept-r1.md:101-103`）写的是 rc=0／0 行 | **成立，但"20:2x"应为"≤20:31"**；结论不变（20:49 的提交把它作废了） |
| gofumpt 本机 `v0.12.0` | `$(go env GOPATH)/bin/gofumpt.exe --version` → `v0.12.0 (go1.27.1)`；⚠ 裸名 `gofumpt` **不在 PATH**（`which` rc=1） | **成立** |
| `sh scripts/d22scan.sh` 此刻 rc=0 | **写码之前**现跑一次：rc=**0**、全仓 **0 枚 finding**、八枚作用域 examined 全非零（`logs/interim-d22scan-sh.txt`）。⚠ 这是**中途读数**，不作门禁交付——AC#7② 要求的正是"全仓级那三把尺排在台件全部入库之后跑最后一次"，所以本程的**判决性读数在最后一格里重新取** | **成立（当时）**，并按 AC#7② 的顺序规则重取 |
| `go vet ./tools/d22scan/` 跨模块 rc=1 | 原文形状现跑：`main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan`，**rc=1** | **成立**（本程一律用模块内 `go vet ./`） |
| 票面 25 行"改前基线有两枚红" | **不成立**（已被票 169 的 `611ae8b` 修掉）：锚 `27f798e` 上 `runtests.sh -C tools/d22scan ./...` = **PASS=30 FAIL=0 SKIP=0、=== RUN=70**，两枚点名用例都在且都绿（`logs/pre-runtests.log`，step-0 之后、写任何东西之前现跑） | **口径落地＝零枚红**，与 161-r1 §1.1 同一处置；判据文字一字未改，只往严格方向动 |

### 1.1 `gofumpt -l .` 那 8 行的归属（本程最该先讲清的一件）

改前（`logs/pre-gofumpt-l-full.txt` 之前的同形跑，逐字）：

```
.scratch\wisp\probes\158\accept-r1\mut\guard.no1.go          <- TRACKED  (506cbae)
.scratch\wisp\probes\158\accept-r1\mut\guard.no2.go          <- TRACKED  (506cbae)
.scratch\wisp\probes\158\accept-r1\mut\guard.no3.go          <- TRACKED  (506cbae)
.scratch\wisp\probes\161\r1\runs\b8-probe-bands\internal\a\p_b8.go                  <- UNTRACKED (161-r1 台件)
.scratch\wisp\probes\161\r1\runs\unp-bad\internal\a\p_unp.go:5:14: expected ')' ...  <- UNTRACKED (故意不可解析)
.scratch\wisp\probes\161\r1\runs\unp-bad\internal\a\p_unp.go:6:8: missing ',' ...    <- UNTRACKED
.scratch\wisp\probes\161\r1\runs\unp-hides-bans\internal\a\p_unp.go:15:14: ...       <- UNTRACKED
.scratch\wisp\probes\161\r1\runs\unp-hides-bans\internal\a\p_unp.go:16:8: ...        <- UNTRACKED
rc=2
```

**三条读数，每条都会骗不同的人**：

1. `gofumpt -l .` 只看**盘上后缀**，完全不知道 git 的存在 ⇒ 在共享工作树里它**永远无法**给出"CI 那一跑会看到什么"。派单 §4 要求它"必须空"，这一条**在别人的台件躺在盘上时不可能达成**，而且做不到的人有强烈动机去删别人的件——那是更坏的结局。
2. rc=**2** 不是"有几枚要格式化"，而是"**有几枚根本 parse 不了**"（161-r1 故意造的 `unp-*` 样本）。一个只读 rc 的人会把它当成"格式化门红了"，而它其实是"这台仪器被我自己的样本噎住了"。
3. 已入库的 3 行 vs 未入库的 5 行，**对 CI 是两个世界**：CI 的检出里只有前者。验收程 20:31 那条"0 行"读数在**它自己的检出语义下**并没有读错——错的是把"本地工作树的 8 行"和"CI 检出的 3 行"当同一个数（这正是 AC#7② 要钉的那件事）。

⇒ 本程的处置：**不碰任何不是我写面的行**（161-r1 的 `runs/**` 在派单 §3 特别名单里，禁改）；把自己那 3 枚格式化（写面内、且只改空白）；并把"谁的红色、入库了没有"做成一发可复算的读数（§5，AC#7②）。

---

## 2. 本程**没测**什么（第一格写作时刻的版本；最后一格会把这一节按"全部读数都在了"重写一遍）

派单 §6 明令"别在证据件里预先引用尚未产出的读数"，所以这一节现在只收**第一格已经成立**的缺口。第二、三格的缺口（以及本节的完整排序版）在最后一格里补，届时每条都给可复算的最小闭合动作。

**C1-N1 — 我没有跑过 CI。** 第一格的读数是**本地** `go test -overlay` 的两遍对照，与本程改没改 `ci.yml` 无关；GitHub 上的真实 step conclusion 一个字节都没量过（本机 gofumpt `v0.12.0` 与 CI 的 `gofumpt@latest` 会不会同判，161-r1 §5.5 也没量，我同样没量）。

**C1-N2 — "格式化后变异语义一字未变"只覆盖这三枚台件消费的那三发。** `guard.no{1,2,3}.go` 分别被 `ovl/r5/r6/r4` 消费（映射是我逐枚 `cat` 出来的，见 §3.2），我只跑了这三发，前后各一次。同目录的 `bridge.M1/M2/M3/as-is`、`guard.as-is.go`、`guard.stub.go` 与 `ovl/r1/r2/r3/r7` 我**没格式化、也没复跑**。r3/r7 用 `guard.as-is.go`（不在那 8 行里 ⇒ gofumpt 判它干净），所以"不受影响"对这两发是**推理，不是读数**。
最小闭合动作：`sh .scratch/wisp/probes/161/r3/run-mut.sh <label>` 把 ovl 名单加上 r3/r7 就是读数（约 2 分钟）。

**C1-N3 — 第一格没量"改前那 8 行里，每一行为什么脏"。** 我只对**我要动的 3 枚**取了 `gofumpt -d` 原文（§3.1，三枚的差异都只是一枚空行）。另外 5 行（161-r1 的 `runs/**`）我一行都没看内容——它们不是我写面，动它们＝踩别人的实验。

**C1-N4 — `.txt` 那一发只证了"能用"，没证"该换"。** §3.4 那一跑证明 `-overlay` 的 Replacement 指向 `.txt` 今天可用，且能复现改前红句。它**没有**证明"161-r1 应该把 `runs/unp-*/p_unp.go` 改成 `.txt`"——那要动别人写面；我只把这个形状记下来给下一位。

**C1-N5 — 我没量"票 158 验收表里那条 20:31 的读数当时是在哪棵树上取的"。** 盘上只有它的**产物**（`gofumpt-full.txt` 0 字节 ＋ §1.3 那句话）。它当时的工作树里还有什么别人的未入库件，我不知道，也无法从 git 反推（未入库的东西不留痕）。这不影响第一格的判据（我重取了自己的改前/改后两遍），但它决定了**为什么"0 行"这句话必须带树**——本程所有 `gofumpt` 读数都写了 rc 与行数，没有一枚只写"空"。

---

## 3. 第一格＝AC#7① 止血：格式化 3 枚入库台件 ＋ 证明变异一字未变

### 3.1 改前三枚文件的形状（`gofumpt -d` 原文，逐字）

三枚的差异**都是且只是一枚空行**（我从 `logs/pre-gofumpt-d-no{1,2,3}.txt` 抄，`-` 行全是空行，没有一枚 `+` 行带内容）：

```
guard.no1.go  @@ -107,7 +107,6 @@   - (blank line above "// 判据 2：...")
guard.no2.go  @@ -113,5 +113,4 @@   - (blank line before the closing "}")
guard.no3.go  @@ -82,7 +82,6 @@   - (blank line above "// 前置：...")
```

⇒ 能格式化，"非不可格式化"那一支（派单 §1 的停手条件）**没有触发**。

**blob 级的三枚对照（这是本格最硬的一发，逐枚现取）**：把 HEAD 的 blob 与入库后的 blob 各取一次字节数与 `\n` 枚数（`git show HEAD:<path>` / `git show :<path>`，Python 数字节，不是 `grep -c`）：

| 枚 | HEAD blob | 入库后 blob | 差 |
|---|---|---|---|
| `guard.no1.go` | 5413 字节 / 130 行 | 5412 字节 / 129 行 | **−1 字节 / −1 行** |
| `guard.no2.go` | 5008 / 117 | 5007 / 116 | **−1 / −1** |
| `guard.no3.go` | 4979 / 123 | 4978 / 122 | **−1 / −1** |

⇒ 一枚空行＝一个 `\n` 字节。**"只改空白"在 blob 层的算术形式就是这一列**：三枚各少 1 字节、少 1 行，其余字节全等（`git diff --cached --numstat` 同一件事的另一种说法：`0 1`、`0 1`、`0 1`）。

改前三枚内容 sha256（`git status --porcelain -- .../mut/` 当时为空 ⇒ 盘上＝HEAD）：

```
071522cd...  guard.no1.go   (130 行)
5eacc332...  guard.no2.go   (117 行)
3594b020...  guard.no3.go   (123 行)
```

改后（`gofumpt -w` 同一批三枚）：129 / 116 / 122 行，sha256 = `39ca85a4...` / `4e595dbd...` / `e7c41c2e...`。

**"只改空白"的第二支机器证明**（不依赖行尾假设）：把改前副本（`pre-format/guard.no1.go`，被 `.gitignore` 挡在库外的那三枚）与改后文件**剥掉全部空白**后取 sha256，三枚逐枚相等：

```
guard.no1.go pre=34f302831b3ec5... post=34f302831b3ec5...  IDENTICAL-MODULO-WHITESPACE
guard.no2.go pre=64cee74a750787... post=64cee74a750787...  IDENTICAL-MODULO-WHITESPACE
guard.no3.go pre=7a6207c2b1d9cd... post=7a6207c2b1d9cd...  IDENTICAL-MODULO-WHITESPACE
```

### 3.1b 本程自己犯的一支仪器错（如实报，因为它差点写进判据）

我一开始用 `grep -c $'\r' <file>` 量行尾，读数**恰好等于文件行数**（129/130/…）——那是这个 shell 把模式当成**空模式**用了（空模式匹配每一行），我差点据此写下"这三枚入库件是 CRLF，`.gitattributes:4` 的 `*.go text eol=lf` 被违反了"。换一支能数字节的仪器（Python 数 `b"\r"`）之后：**CR=0，LF=129/130**，`git ls-files --eol` 同判（`i/lf w/lf attr/text eol=lf`）。
⇒ 三条留给下一位的读数规矩：**数 CR 不许用 `grep -c`**（本仓已记过 `grep -c` 命中 0 → rc=1 吃掉 `&&` 链，这是同一族的另一半：命中"全行"时它不报错、只给错数）；`git ls-files --eol` 是行尾问题的**权威仪器**；结论换仪器重取一次再落笔。
⇒ 顺带一条对上派单 §4 的仪器坑清单：本仓 `* text=auto` ＋ `core.autocrlf=true` 这一族，今天**没有**在这三枚上产生 CRLF 假象（161-r1 §5.5 那句"盘上是 LF"与我的复算一致）。


### 3.2 变异台件消费关系的现读（不抄派单）

`ovl/*.json` 逐枚 cat 出来的映射（本程现读）：

| ovl | bridge 侧 | guard 侧 | 消费哪枚 |
|---|---|---|---|
| `r4-m2no3.json` | `mut/bridge.M2.go` | **`mut/guard.no3.go`** | 摘判据③ |
| `r5-m1no1.json` | `mut/bridge.M1.go` | **`mut/guard.no1.go`** | 摘判据① |
| `r6-m1no2.json` | `mut/bridge.M1.go` | **`mut/guard.no2.go`** | 摘判据② |
| `r1/r2` | M1 / — | `guard.stub.go` | 桩件 |
| `r3/r7` | M2 / M3 | `guard.as-is.go` | 原样三味 |

⇒ 我格式化的是**这三枚**，所以要证的就是 **r5／r6／r4** 三发（N5 讲清了我没证哪些）。

### 3.3 前后两遍读数（同一枚 `-overlay` 组合，各跑一次；全程未与 `-cover*` 同用）

命令原文（每发一枚日志，rc 单枚存档）：

```
$ go test -count=1 -v -overlay=.scratch/wisp/probes/158/accept-r1/ovl/r5-m1no1.json ./internal/tools/
$ go test -count=1 -v -overlay=.scratch/wisp/probes/158/accept-r1/ovl/r6-m1no2.json ./internal/tools/
$ go test -count=1 -v -overlay=.scratch/wisp/probes/158/accept-r1/ovl/r4-m2no3.json ./internal/tools/
```

| 发 | 改前四数 | 改前 rc | 改后四数 | 改后 rc |
|---|---|---|---|---|
| r5（M1＋摘判据①） | RUN=116 PASS=79 **FAIL=1** SKIP=0 | 1 | RUN=116 PASS=79 **FAIL=1** SKIP=0 | 1 |
| r6（M1＋摘判据②） | RUN=116 PASS=79 **FAIL=1** SKIP=0 | 1 | RUN=116 PASS=79 **FAIL=1** SKIP=0 | 1 |
| r4（M2＋摘判据③） | RUN=116 PASS=78 **FAIL=2** SKIP=0 | 1 | RUN=116 PASS=78 **FAIL=2** SKIP=0 | 1 |

红名（改前＝改后，逐枚）：r5/r6 = `TestSensitiveReadAcrossTheBridgeOpensItsTaskScope`；r4 = `TestFSListSummarizesADirectory` + `TestAtomicWriteKillsMidWrite`。
⇒ 与票 158 验收表 §2.3 那张表的 r5/r6/r4 三行**逐枚对得上**（档位〔验收程读数，我复算并同意〕）。

红句逐字节对照（`diff logs/pre-<ovl>.red.txt logs/post-<ovl>.red.txt`）：

```
r6-m1no2 : BYTE-IDENTICAL
r4-m2no3 : BYTE-IDENTICAL
r5-m1no1 : 两处差别 ——
  < --- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope (0.01s)
  > --- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope (0.02s)     <- 计时后缀
  < bridge_scope_open_ticket158_test.go:122: 收尾审计行没有承认这枚 scope 是开过的：...
  > bridge_scope_open_ticket158_test.go:121: 收尾审计行没有承认这枚 scope 是开过的：...  <- 行号
```

**明写：r5 那一发的红句变了，变的是行号，不是判据。** 我把空行摘掉了一枚，`t.Fatalf` 的调用点从 122 行上移到 121 行。两行并排（这就是"让下一位能判"的那一发）：

```
改前 guard.no1.go:122   ==   改后 guard.no1.go:121     （cmp 输出：SRC-LINE IDENTICAL）
		t.Fatalf("收尾审计行没有承认这枚 scope 是开过的：找不到含 %q 的行；"+
```

改后文件的 122 行现在是 `			"C25 scope closed 全部读数=%q\n"+`（同一枚 `t.Fatalf` 的第二段拼接），不是另一枚判据。

再把**只遮盖 (a) 计时后缀 (b) 文件名:行号**这两件事之后重新 `diff`，三发全部**逐字节相等**（`logs/{pre,post}-<ovl>.red.norm.txt` 三对，脚本 `extract-red.sh` 产出）。
⇒ 第一格的答话：**红名相同、四数相同、红句除行号位移外逐字节相同、且行号位移的位置就是被摘掉的那枚空行的下方**——"只改空白不改语义"这一句现在是**读数**，不是话。（计时后缀那一处是墙钟读数，与本判据无关，我照样列出来而不是抹掉。）

### 3.4 一个额外读数：改前的行为现在**只从库里的件**就能复跑

格式化之后，"改前状态"在盘上只剩我的副本。为免它变成"只有我这台机器能重跑"的读数，我做了一枚 `.txt` 镜像（`pre-format/guard.no1.go.txt`，sha256 与改前 `.go` **逐字节相等**）并用它当 `-overlay` 的 Replacement 再跑一次 r5：

```
$ go test -count=1 -v -overlay=.scratch/wisp/probes/161/r3/ovl/r5-m1no1-fromtxt.json ./internal/tools/
rc=1   RUN=116 PASS=79 FAIL=1 SKIP=0
--- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
    bridge_scope_open_ticket158_test.go:122: 收尾审计行没有承认这枚 scope 是开过的：...   <- 与 §3.3 的"改前"那一行逐字节相同
```

⇒ 两支结论：① 改前的红句**不依赖我这台机器**，从已提交件即可复算——而且镜像的**入库 blob 与 HEAD 的原 blob 逐字节相同**（三枚都是 `(bytes, LF) = (5413,130) / (5008,117) / (4979,123)`，`git show HEAD:<mut 件>` 与 `git show :<r3 .txt 镜像>` 两两相等，见 §3.1 表），所以我格式化那三枚**没有销毁任何取证能力**；② 派单 §1 提的"把 Replacement 指向 `.txt` 后缀、JSON 里映射到 `.go` 目标"那一形状**今天真的能用**（这是读数不是建议）。
⚠ 一条行尾账（不是问题，但下一位在 Linux 上复跑会撞）：`.txt` 走 `* text=auto`、`.go` 走 `text eol=lf`，两者入库都是 LF；在 `core.autocrlf=true` 的 Windows 检出上 `.txt` 会回到 CRLF、`.go` 保持 LF。行数相同（`\r` 不产生新行），所以 `-overlay` 的行号读数不变；真要逐字节复算 §3.3 的"红句对照"，请用**本程盘上那份**（LF，见 §3.1b）或在取用前统一行尾。
⇒ 顺带把本程自己的台件对齐到那条规矩：`pre-format/*.go` 三枚**只在盘上、被 `r3/.gitignore` 挡在库外**（内容＝改前字节的取证，格式化它们会毁掉证据），入库的是 `.go.txt` 镜像。⚠ 后果如实记：**本程让 `gofumpt -l .` 多了 3 行**（见 §3.5），因为 gofumpt 不看 gitignore。

### 3.5 第一格之后的 `gofumpt -l` 读数（不许藏的那一半）

```
$ $(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm      rc=2，8 行
  5 行 = probes/161/r1/runs/**          （UNTRACKED，161-r1 台件，不是我写面，没碰）
  3 行 = probes/161/r3/pre-format/*.go  （UNTRACKED，我自己的改前取证副本，本件 §3.4 说明为什么不能格式化）
  0 行 = TRACKED
```

⇒ **派单 §4 那句"`gofumpt -l .` 必须空"在共享工作树上今天做不到，而且不是本程做不到的**：只要 161-r1 的假根还在盘上，它就不空（那是它写面内的正当台件）。能达成、也确实达成的是**"已入库集合为空"**——也就是 CI 那一步的语义。把"这条红是谁的、入库了没有"做成一跑就能读的东西＝AC#7② 的活儿，本程在后面的格子里交，**写作此刻还没有那一发的读数**，所以这一段只记现状。

---
