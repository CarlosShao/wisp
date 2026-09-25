# d22scan gitignore 修复批 · 验收 r1 的三枚附条件就地闭合（close r1）

**我是谁**：实现方（写码代理），地界＝`tools/d22scan/**` ＋本件一枚路径。
**被闭的账**：`docs/evidence/s1/d22scan-gitignore-fix-accept-r1.md`（460 行，非实现者裁 ca84b75
＝**附条件成立**，总裁行点名三件、"全在一枚小 commit"）。本件只做那三件与它们的证明，
**验收已判成立的 anything 一律未动、未"顺手改善"**。

| 用途 | 锚点 | 现量方式 |
|---|---|---|
| 开工锚点（我读到的 HEAD） | `43a60432ac0c08a8f5f465660e7f30b885cd98ce` | `git rev-parse HEAD` @2026-09-25 10:00(+08) |
| 每节刷新 | §1=`888bbd5`，§2/§3=`888bbd5`，§4/§5/§6=`888bbd5` | 每节现跑 |
| 临时件（只建不删，全在仓库外） | `D:\tmp\d22scan-close-r1\` | `probe\`（真 git 前提台件）／`mut\`（首轮四发变异）／`mut2\`（终稿四发变异）／`gitignore.go.final-orig`（变异前 pristine 拷贝）／`rowA-D.log`／`baseline-gate.log`／`final-gate.log`／`full.lst` `narrow.lst`／`*-code-only.go`／`roster-{before,after}.txt` |

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

### 1.4 双向变异读数（四行，全在仓库外副本 `mut2\`，被测码逐字＝交付版）

| 行 | 树上代码 | 新测试 | 上一枚 F1 测试 | `TestGitIgnoreRuleSemantics` |
|---|---|---|---|---|
| **A** | 交付原样 | **PASS** | PASS | PASS |
| **B** | **M1**：`parseIndexPaths(all)` → `(withIndex)`（`gitignore.go:245`，＝摘掉全量那一腿；另加 `_ = all` 仅为让副本编译） | **FAIL** `scan_test.go:1633`（`holds()` 认不出索引明明持有的件） | PASS | PASS |
| **C** | **M4a**：`TrimRight(raw," \t\r")` → `TrimSpace(raw)`（`gitignore.go:577`，语义修回退，全量腿仍在） | **PASS**（正是全量腿把它救回来的） | PASS | **FAIL** 两支（`lead/inside.txt` / `"  lead/inside.txt"`，`scan_test.go:1787`） |
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
