# d22scan gitignore 修复批 · 验收 r1 的三枚附条件就地闭合（close r1）

**我是谁**：实现方（写码代理），地界＝`tools/d22scan/**` ＋本件一枚路径。
**被闭的账**：`docs/evidence/s1/d22scan-gitignore-fix-accept-r1.md`（460 行，非实现者裁 ca84b75
＝**附条件成立**，总裁行点名三件、"全在一枚小 commit"）。本件只做那三件与它们的证明，
**验收已判成立的 anything 一律未动、未"顺手改善"**。

| 用途 | 锚点 | 现量方式 |
|---|---|---|
| 开工锚点（我读到的 HEAD） | `43a60432ac0c08a8f5f465660e7f30b885cd98ce` | `git rev-parse HEAD` @2026-09-25 10:00(+08) |
| 每节刷新 | §1=`888bbd5`，§2/§3=`888bbd5`，§4/§5/§6=`888bbd5` | 每节现跑 |
| 临时件（只建不删，全在仓库外） | `D:\tmp\d22scan-close-r1\` | `probe\`（真 git 前提台件）／`mut\`（首轮四发变异）／`mut2\`（终稿四发变异，末态已复原并按 hash 核过）／`gitignore.go.orig` `gitignore.go.final-orig`（变异前 pristine 拷贝）／`rowA-D.log` `final-rowA-D.log`／`baseline-gate.log` `final-gate.log`／`full.lst` `narrow.lst`／`*-code-only*.go`／`roster-{before,after}.txt` |

**地界现量**（开工前）：`git status --porcelain -- tools/d22scan docs/evidence/s1` 空 ⇒ 我接手时这两处无未提交改动；
`git ls-files tools/d22scan` ＝ 6 枚（`d22scan.exe` 未入库）。
**未碰清单**：`emojiRe`、`allowlist.txt`、任何阈值/golden、`internal/**`（含 `internal/panel/**`）、`frontend/**`、`design/**`
——见 §5 的逐枚读数。**未 push、未 `git add -A`/`.`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`。**

---

## 1. 条件①：补 §4 的"第三枚半"——全量 `git ls-files` 那一味从此有断言要求它

### 1.1 加了什么（`scan_test.go:1541-1640`，一枚新顶层测试 ＋101/−0）

`TestFullTrackedListCoversWhatTheNarrowListCannot`（`scan_test.go:1568`）。按验收 §4 末"最小修法"逐件落地：

- **台件**：`frontend/.gitignore` 追加 `  weird/*`（**行首两枚空格**）、种子 `frontend/weird/inside.tsx` 含 U+2264、
  **只跑普通 `git add -A`**（那一枚不被 force-add；同树的 `frontend/dist/tracked-or-not.tsx` 才 force-add，用途见下条）。
- **断言两半**（验收点名的那两半）：`git ls-files -i -c --exclude-standard` 在这枚树上**为空**（对那一枚种子件而言），
  ＋ **finding 点名** `frontend/weird/inside.tsx`。外加两半本文件风格的配套：分母 ＋1（两枚 frontend walk）、
  `runVerdict` 的 rc 必须是 1。
- **"0 命中"的正控写在同一枚断言里**：我把 `-i -c` 读成**整串等值**判断
  （`== "frontend/dist/tracked-or-not.tsx"`）而不是"不包含 weird"——同一棵树、同一条命令里
  既必须**打出** force-add 那一枚（仪器在看着这棵树），又必须**不打**出 weird 那一枚（前提仍然成立）。
  这一枚是 `t.Fatalf`：正控塌了本件就不再是证明，测试自己先红。
  若将来 git 或规则改动让 `-i -c` 也报出 weird，另有第三枚 `t.Fatal`（`holds` 块里）点名"这一枚台件不再钉全量腿，换一个形状"。
- **那一味本身的断言（打红 M1 的就是这一枚）**：`s.ign.index().holds("frontend/weird/inside.tsx", false)` 必须为真。
  它是只读的、同包内既有风格（本文件已经直接调 `g.decide(...)`、`ign.note()`），**不改生产码一行**。

**为什么必须有全量那一腿、窄清单结构上看不见这一类**：`-i -c` 只会举出"git 自己说撞了规则"的件，
而这一类的定义正是"匹配器比 git 多跳"（git 说没忽略 ⇒ 窄清单必然不说）。⇒ 窄清单举不出它。
这一类**不是假想**：同一批刚修过两枚实例（F3 的行首空白、F3 的尾随 `**`），所以它是活家族的保险。

### 1.2 夹具机制＝本模块既有那一套，未新造 helper

树全部落在 `t.TempDir()` 里，走包内**原有** helper：`liveFixture(t)`（内部就是 `t.TempDir()`）
＋ `seedFile` ＋ `gitIndexFixture`（`git init -q -b main` ＋ `git add -A`）＋ `gitForceAdd` ＋ `scanFixture` ＋
`gitCommand` ＋ `runVerdict`。**核销**：`grep -c '^func ' tools/d22scan/scan_test.go` 现量 **51（HEAD 版）→ 52（我的版）**，＋1 就是我新加的那枚测试函数本身
⇒ **包级 helper 一枚没新增**；测试内部的 `seed := func(t *testing.T, withWeird bool)` 是函数内闭包，
形状照抄上一枚测试 `TestTrackedPathsAreNeverSkippedByTheIgnoreFilter` 已有的同名局部 `seed`。**真树的 `frontend/`、`design/` 一字节未写。**

### 1.3 真 git 前提台件（仓库外 `probe\`，先问 git 再问被测码）

```
$ git ls-files                                   .gitignore frontend/.gitignore
                                                 frontend/dist/tracked-or-not.tsx frontend/weird/inside.tsx
$ git ls-files -i -c --exclude-standard          frontend/dist/tracked-or-not.tsx        <- 只此一枚
$ git check-ignore -v frontend/weird/inside.tsx  rc=1                                    <- git 说不忽略
$ comm -13 full.lst narrow.lst                   (空)   <- 窄清单 ⊆ 全量清单（条件③的前提）
$ comm -23 full.lst narrow.lst                   .gitignore / frontend/.gitignore / frontend/weird/inside.tsx  <- 反方向正控非空
```
⇒ "普通 `git add -A` 就追踪、而 `-i -c` 一辈子看不见"这一形在真 git 2.52.0 上成立；
subset 那一枚 0 命中是拿反方向非空正控一起量的（两把尺各一枚）。

### 1.4 双向变异读数（四行，全在仓库外副本 `mut2\`，被测码逐字＝交付版；**终稿又跑了一遍**）

首轮在条件②③的注释落地前跑（`rowA-D.log`），终稿落地后按同一串命令重跑（`final-rowA..D.log`，
`gitignore.go` 终稿 hash `6d07693`）⇒ **两遍四行读数逐枚同值**，下面引的是终稿那一遍：

| 行 | 树上代码 | 新测试 | 上一枚 F1 测试 | `TestGitIgnoreRuleSemantics` |
|---|---|---|---|---|
| **A** | 交付原样 | **PASS** | PASS | PASS |
| **B** | **M1**：`parseIndexPaths(all)` → `(withIndex)`（终稿 `gitignore.go:248`；首轮量在 `:245`，差三行是条件②③的注释加出来的；另加 `_ = all` 仅为让副本编译） | **FAIL** `scan_test.go:1633`（`holds()` 认不出索引明明持有的件） | PASS | PASS |
| **C** | **M4a**：`TrimRight(raw," \t\r")` → `TrimSpace(raw)`（终稿 `gitignore.go:580`，语义修回退，全量腿仍在） | **PASS**（正是全量腿把它救回来的） | PASS | **FAIL** 两支（`lead/inside.txt` / `"  lead/inside.txt"`，`scan_test.go:1787`） |
| **D** | **M4a ＋ M1** | **FAIL** 四枚：`scan_test.go:1606`（ban #6 分母 3 vs 3）、`:1609`（ban #8 同）、`:1619`（finding 没点名，findings 里只剩 dist 那一枚）、`:1633`（`holds()`） | PASS | FAIL |

⇒ **要求的两向都拿到了**：摘掉全量 `ls-files` 那一味 ⇒ 新测试**必红**（B 行）；它在 ⇒ 新测试**红不着**（A/C 行）。
C 行同时是这一味"防御不是装饰"的正证：**语义修一旦回退，今天只有这一腿挡着**（D 行的三枚行为红＝那一枚受追踪的交付字节对门永久隐身）。

⚠ **一处对验收文字的不同意（照单收会收错）**：验收 §4 末句写"从此摘掉 `all` 那一味必红"，
但它自己的读数表第 4 行是"只 M1（语义修完好）⇒ PASS"。盘上复算同它第 4 行一致：**纯行为的断言（分母／finding／rc）在 B 行全都不响**，
只有 `holds()` 那一枚白盒断言咬得住 M1 单独一发。⇒ 我在它点名的行为台件之外**多加了那一枚断言**，
否则这条闭不上；这不是扩大范围，是它那张表已经量到、但文字写歪了的那一位。

### 1.5 这一节的门读数（真树，`sh scripts/d22scan.sh` 整道，HEAD `888bbd5`）

step 1（正控）：**PASS=29 FAIL=0 SKIP=0 ＝ 29，`=== RUN`=69**（开工基线同树同命令＝28/0/0/68，见 §4），
名册差集 **＋1／−0**（唯一新增名 `TestFullTrackedListCoversWhatTheNarrowListCannot`，**没有一枚消失**）。
step 2 rc=0、八数未动（§4 逐枚）。整道门 rc=0。
日志：`D:\tmp\d22scan-close-r1\final-gate.log`（变异前的基线在 `baseline-gate.log`）。

---

## 2. 条件②：`gitignore.go` 那句"唯一形状"改成两支并列、并归对治理者（现号 `:103-119`）

**它点名的行号没漂**：`git show HEAD:tools/d22scan/gitignore.go | awk 'NR==103||NR==105'` 现量首尾正是
验收抄的那三行（"The one shape that could make it skip MORE than git …"）。**被换掉的原文（HEAD `:103-105`）**：

> The one shape that could make it skip MORE than git - not being able to reach the index at all -
> is handled by the rule above it, not by this list: no rule is applied then, so the tool over-scans and says why.

**换成的新文（工作树 `:103-119`）**：两支并列，**无最高级**，各自点名治理者与钉它的测试——

1. **问不到 index**（无仓库／无 git 二进制／超时／别人仓库的子目录）
   ⇒ 治理者＝`skip()` 的 `!ix.ok` 分支（一条规则都不应用）＋ `note()` 的第一行自陈。
   这一支**照旧成立**，验收 §3 四支各测过。
2. **匹配器把自己的规则读得比 git 宽**（＝它 §4 第四发 M4a 那一类，也是 F3 抓过的两枚实例：行首空白被 trim 成
   另一个模式、尾随 `**` 吃掉自己的父目录）⇒ 治理者**不是**"问不到就不套规则"那一支（这一支里 git 明明问得到、
   并且回答"没忽略"），治理者＝`holds()` 查**全量** `git ls-files` 集合而不只是"受追踪∩撞规则"那枚子集。
   钉它的是 `TestGitIgnoreRuleSemantics`（两枚读法本身）＋ §1 新加的那枚（挡下一枚的保险）。

**末尾明写这不是穷尽**：`// This is a list of the shapes seen so far, not a claim that there are no others.`
——理由：验收要的是"不要最高级"，而我若写"只有两支"仍然是同一枚未钉的 blanket（今天能举出两枚，
明天 F3 家族第三枚照样会长出来）。⇒ 措辞是"two shapes are known to do it"。
**真树证据＝§1.4 的 C/D 两行**：M4a 单独发时行为不变（全量腿挡着＝第二支的治理者真的是它），
M4a＋M1 时同一发当场隐身（第二支没有第一支可替）。

**未碰的同类句子（登记，不擅自扩大）**：`scan_test.go:1179` 与 `:1194` 里还有两句 "the one shape …" 式措辞，
属**上一批**（`3bb99aa`）写的注释、讲的是独立对照 walk 在 F1 树上的形状，**不在验收点名的三件之内**，本件一字未动。

---

## 3. 条件③：`-i -c` 那句"独立并集成员"换成真理由 ＋ `:11` 的引用改回 `:22`

### 3.1 `gitignore.go:180-194`（原 `:166-171`）

被换掉的原文末两行：
> … the set the self-report names, **and an independent union member so one command's parsing bug cannot blind the other**

我自己复算过它站不住（**不是只抄验收**，两把尺各一枚，`D:\tmp\d22scan-close-r1\probe\` 真 git 2.52.0）：

| 读数 | 命令 | 结果 |
|---|---|---|
| 子集关系 | `comm -13 full.lst narrow.lst` | **空** ⇒ `-i -c` 相对 `ls-files -z` 一枚新的都带不进 `holds()` |
| **正控**（防"0 命中＝尺子瞎"） | `comm -23 full.lst narrow.lst` | **非空 3 行**（`.gitignore`／`frontend/.gitignore`／`frontend/weird/inside.tsx`） |
| 枚数 | `wc -l` | full=4、narrow=1（同一棵树、同一次 `git add -A` ＋ 一枚 force-add） |
| 同一枚解析器 | `grep -n parseIndexPaths` 现量（终稿行号） | `parseIndexPaths` 定义在 `gitignore.go:301`，其第 302 行调的就是 `parseIndexPathsList`（定义 `:324`）；窄清单在 `:252` 直接叫同一枚函数 ⇒ 解析 bug 同打两枚 |

⇒ 新文改成两件事实：**它不是 `holds()` 里的第二意见**（严格子集 ＋ 同一枚 `parseIndexPathsList`），
**它承重的是 `note()` 点名的那三行里的第三行**（CI 日志里能看见"哪几枚交付字节正压在活规则下"，`gitignore.go:448`），
并指回 §1 那枚测试去钉"让跳过判定安全的是全量那一腿的射程，不是这份清单"。
原句里"the ticket's named instrument"与"the set note() names"两句是**真话**，保留。
`note()` 那一行本身**一字未改**，M6 那枚 pinner（`scan_test.go` 里 `note "" must name "TRACKED and matching an ignore rule"`）仍咬得住。

### 3.2 `gitignore.go:11` 的 `.gitignore:24` → `.gitignore:22`

现量（本仓根 `.gitignore`）：`22:frontend/dist/*`、`24:assets/web/dist/` ⇒ 旧引用指到了**另一条规则**上
（那是 `3bb99aa` 留下的，验收 N5 已登记"不是 ca84b75 写的"）。同段另一枚引用 `frontend/.gitignore:12`（`dist/*`）
现量**同值**，未动。

### 3.3 这一节的两枚条件都是纯文字：生产码零改动的证法

- `git diff --numstat -- tools/d22scan/gitignore.go` ＝ **＋30/−7**；
- 把两版的**非整行注释**都剥掉再 diff：
  `git show HEAD:…gitignore.go | grep -vE '^[[:space:]]*//' > A`、`grep -vE '^[[:space:]]*//' <work> > B`、`diff A B` ⇒ **空**（481 行 → 481 行）；
- `git diff -U0` 的改动行里**没有一枚非 `//` 开头的行**（`grep -vE '^[+-][[:space:]]*//'` 现量空）；
- `gofmt -l tools/d22scan` **空输出**、`go vet ./...` 与 `go build ./...` 在 `tools/d22scan` 内**全过**；
- 变异台件的 pristine 拷贝与真树文件做 `git hash-object` 三向同值（`baac05c3ae99ab775c2092759dada4a11a592f12`），
  **我从未把工作树文件改成变异态**——四发变异全在 `D:\tmp\d22scan-close-r1\mut2\` 里跑，因此也**没有用过 `checkout`/`reset`/`stash`**。

