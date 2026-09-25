# d22scan 三条件收尾批（`304aeec`..`6a5b321`）· 非实现者对抗验收 r1

**我是谁**：裁决者（非实现者），本轮唯一可写路径＝本件。零代码改动、零 push。
**被验交付**：`docs/evidence/s1/d22scan-gitignore-fix-close-r1.md`（267 行，实现方自陈）＋它的两枚代码差
`tools/d22scan/scan_test.go` +101/-0（`304aeec`）／`tools/d22scan/gitignore.go` +30/-7（`b7c06d2`，自称注释级）／
三枚只动表的追加枚（`cec5e78`/`38f340b`/`6a5b321`）。
**它的上一环**：`docs/evidence/s1/d22scan-gitignore-fix-accept-r1.md` §4（工作单，含那行"只 M1 ⇒ PASS"与末句"从此摘掉 `all` 那一味必红"）。
**本件的靶心**：那枚白盒钉 `holds("frontend/weird/inside.tsx", false)`（`scan_test.go:1632`）是**承重**还是**镜子**，
以及实现程顶回上游验收措辞这件事的**程序**账。

| 用途 | 锚点 | 现量方式 |
|---|---|---|
| 我开始时的 HEAD | `310816f94216af4320b3a8ffa423eea5c4e1028d`（短 `310816f`） | `git rev-parse HEAD` |
| 我量 §0 时的 HEAD | `310816f`（§0 全程未变，末次复查同为 `310816f`） | 每节开头/结尾各查一次 |
| 被测批次 | `304aeec` → `b7c06d2` → `cec5e78` → `38f340b` → `6a5b321` | `git log --oneline 304aeec^..6a5b321` |
| 我读到的台账账 | `A207`/`A214`/`A218`/`A221②`/`A223`（`A223` 当声明读，不当证据读） | `docs/reports/pending-and-issues.md:5969,5989` |
| 我的仓库外工作区（只建不删，全在 `D:\tmp\d22scan-close-r1-accept\`；**现量 `find -maxdepth 1`＝15 个目录 ＋ 30 枚文件**） | 代码拷贝：`head\`（`tools/d22scan` 的 pristine 拷贝，`git hash-object` 三向核过，见 §1.1）／`pre\`（同一枚包在 `304aeec^` 的版）／`fullpre\`＝`git archive 304aeec^` 全量快照／`fullhead\`＝`git archive HEAD` 全量快照／`mut\M1` `mut\M4a` `mut\M4aM1` `mut\SN3-head` `mut\SN3-M1`（另 `mut\SEEDNARROW-*` 两枚是我第一次 sed 没生效的废件，留着）／二进制：`bin\d22-{post,M1,M4a,M4aM1}.exe`／我自己的三枚探针（**都不是交付件、都不进仓**）：`fuzz\zz_overread_probe_test.go` `fuzz\zz_combo_probe_test.go` `fuzz\zz_race_probe_test.go`／树：`probe\seed1`（§3 的真 git 台件）`seedtree\`（§1 的输出对照树）`dirprobe\`（目录规则那一发）`outer\`（§7.4 那支"临时根落在别人仓库里"）`fakehome\` ＋ `empty.gitconfig` ＋ `gitconfig-global` ＋ `mygitignore`（§7.2 的配置实验）／日志与名册：`step1-accept.log` `gate-accept.log` `gate-final.log` `vet.log` `combo-fuzz.log` `roster-{pre,head,fullpre,fullhead}.log` `status-{before,after}.txt` `names-{pre,real}.txt` `npre.txt` `nhead.txt` `r1.txt` `r2.txt` `out-*.txt` `dp-*.txt` `real-*.txt` `A.go` `B.go`（§6 那把"剥注释再 diff"的尺）`tmpwhere.go`（量 `os.TempDir()` 的三行小程序） | `find -maxdepth 1 -type d \| wc -l` ＝15、`-type f` ＝30 |

**并行读数的边界（先说清楚，免得我的数被人当别人的数用）**：**此刻有两路在同时作业**——
一位实现程在 `internal/panel/**` 上（`A222` 那三件补强），前端会话在 `frontend/**` 上写（`A227` 刚把它从 cold 切成 running）。
本件**没有**跑过整树 `go test ./...`；`go vet ./...` 与 `go build ./...` 我跑了（rc 均 0，见 §0.3），但那是编译不是行为读数。
§0.2 那八数是**我在 `310816f` 那一刻**的读数：`ban #8 internal/=407` 归 panel 那位（不是本批），
而**同一道门在我程末（`4fab445`）复跑时 `frontend/` 已经从 40 变 43**——那笔账单独记在 §0.5，别混进被验那批。

---

## 0. 闸门重跑（全部我现跑，锚点 `310816f`）

### 0.1 step 1（正控，`sh tools/d22scan/runtests.sh -C tools/d22scan ./...`）

```
$ git rev-parse --short HEAD
310816f
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./... ; echo "rc=$?"
rc=0
（末三行，逐字）
PASS
ok  	github.com/CarlosShao/wisp/tools/d22scan	32.050s
runtests.sh: OK - packages=[./...] top-level: PASS=29 FAIL=0 SKIP=0, === RUN=69, '[no tests to run]'=0
```

⇒ 派单里那对数 **PASS=29 / FAIL=0 / SKIP=0 / RUN=69 复现成立**（`28→29`、`68→69` 的"开工基线"那一半见 §6.1，我用另一把尺另量了一次）。
日志：`D:\tmp\d22scan-close-r1-accept\step1-accept.log`。

### 0.2 step 2（真树扫描，`sh scripts/d22scan.sh`）

```
$ sh scripts/d22scan.sh ; echo "rc=$?"
rc=0
（两处 d22scan.sh 标记都在，证明 step 2 确实跑在 step 1 之后）
d22scan.sh: positive control - runtests.sh -C tools/d22scan ./...
d22scan.sh: scan of /d/work/workspace/projects plans/Wisp
（真树读数，逐字 8 行）
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  32 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 407 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  39 Go files, comments and _test.go included
```

⇒ **八数 `203/22/40/18 · 32/40/407/39` 与派单给的、与交件件 §4.2 的第三列，三处同值；整道门 rc=0。**
日志：`D:\tmp\d22scan-close-r1-accept\gate-accept.log`（step 1 段与 §0.1 同值，`ok 32.050s` 那份）。

### 0.3 卫生三项

```
$ gofmt -l tools/d22scan        -> 空输出
$ gofmt -l .                    -> 空输出（整仓，不止那一枚目录）
$ go vet ./...                  -> rc=0，输出 0 行（D:\tmp\...\vet.log 行数=0）
$ go build ./...                -> rc=0
$ (cd tools/d22scan && go vet ./...) -> rc=0
```

### 0.4 测试跑完之后真树没有多出残渣（这一格同时属于 §2）

```
$ git status --porcelain > status-before.txt   # 20 行
$ go test ./tools/d22scan（经 §0.1 那把尺跑过一遍）
$ git status --porcelain > status-after.txt    # 20 行
$ diff status-before.txt status-after.txt      -> 空输出 ⇒ STATUS_IDENTICAL
```

那 20 行全是 owner 自己的东西：16 枚 `design/**` 未提交删除 ＋ 4 枚未追踪
（`design/doubao/01-ball-states.jpg`、`design/doubao/demo/lib/`、`design/doubao/demo/screenshots/`、`design/old/`）。
⇒ **`frontend/**` 与 `.gitignore` 在本轮任何一次测试运行里都没有被写过**（机制见 §2）。

### 0.5 ⚠ 我这一程跑完之后，真树的八数**变了**——变得与被验那批无关，但必须写在这里

交件前我又整道跑了一遍 `sh scripts/d22scan.sh`（锚点已从 `310816f` 走到 `4fab445`，
其间分支上落了 **21 枚** commit，其中 **4 枚是我自己这张表**、其余 17 枚是他人的）：

```
$ sh scripts/d22scan.sh ; echo "rc=$?"
rc=0
（末行逐字）
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=203, bans #1-5 cmd/=22, ban #6 frontend/=43, ban #7 internal/tools/=18, ban #8 design/=32, ban #8 frontend/=43, ban #8 internal/=407, ban #8 cmd/=39
$ git ls-files --others --exclude-standard -- frontend
frontend/scripts/render-stream.tsx
frontend/src/components/reveal-text.tsx
frontend/src/lib/reveal.ts
$ git diff --name-status 310816f..HEAD -- frontend        -> 空输出
$ git ls-files frontend | wc -l                            -> 40（与 310816f 上 ls-tree 的同值）
$ git status --porcelain | wc -l                           -> 30（我开始时是 20）
```

⇒ **两枚 `frontend/` 从 40 变 43，差的是 3 枚"未追踪也没被 ignore"的工作树文件**（第三方正在往 `frontend/**` 写；
`A227` 那格刚把前端会话从 cold 切成 running，我这程没去核它）。**归因清楚**：区间内没有一枚 commit 碰过 `frontend/`
（`git diff --name-status … -- frontend` 空），被验那批的两枚代码路径是 `tools/d22scan/**`（§9），
所以**这不是那批的回归、也不是我这程造的**。
⇒ 两条读法账：① §0.2 那八数只对**我量它的那一刻与那枚锚点**负责（本仓 `A220⑤` 已经为这件事改过做法，我照做）；
② 任何人今天重跑这道门会看到 `43/43`，**不要把它读成收尾批把门改坏了**。
`ban #8 internal/` 仍是 407（那枚并行的 panel 文件到我现在为止没有再长）。

---

## 1. 生死格：那枚白盒钉是承重还是镜子

**先声明我判它的方式**。题面给的三读法：(i) 承重、(ii) 镜子（⇒ 该删那味）、(iii) 更窄的一种。
这仓对承重的定义是"摘掉任一成分，若存在一发变异从此不再变红，那枚成分就是承重的"——题面也提醒了：
**只改一枚内部 struct 字段、永远不改任何 finding 的那种变异，正是这把尺会给出误导答案的地方**。
所以我不让公式替我想，先把盘上事实量出来，再回到公式。

### 1.1 被测码＝HEAD 逐字（先钉住"我量的确实是交付版"）

```
$ mkdir -p /d/tmp/d22scan-close-r1-accept/head
$ cp <repo>/tools/d22scan/{*.go,go.mod,allowlist.txt} /d/tmp/d22scan-close-r1-accept/head/
$ git hash-object gitignore.go                       （我的拷贝）
6d0769375aa5bcf4ac05b2179c62996912222692
$ git hash-object tools/d22scan/gitignore.go         （真工作树）
6d0769375aa5bcf4ac05b2179c62996912222692
$ git show HEAD:tools/d22scan/gitignore.go | git hash-object --stdin
6d0769375aa5bcf4ac05b2179c62996912222692
```

⇒ 三向同值，且**与交件件 §3.3 自陈的那枚 `6d07693…` 同值**。四枚变异全部种在仓库外副本，真树从未进入变异态。

### 1.2 复现交件件 §1.4 那张四行表（我自己跑，不抄它的数）

命令（四份目录各跑一次；`head`＝交付态，`mut/M1`＝`parseIndexPaths(all)`→`(withIndex)`＋`_ = all`，
`mut/M4a`＝`TrimRight(raw," \t\r")`→`TrimSpace(raw)`，`mut/M4aM1`＝两发同时）：

```
$ (cd <dir> && go test -count=1 -v -run 'TestFullTrackedListCoversWhatTheNarrowListCannot|TestTrackedPathsAreNeverSkippedByTheIgnoreFilter|TestGitIgnoreRuleSemantics' .)
```

| 行 | 我读到的数（逐字） | 交件件的数 | 判 |
|---|---|---|---|
| **A** 交付 | `--- PASS: TestTrackedPathsAreNeverSkippedByTheIgnoreFilter`／`--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot`／`--- PASS: TestGitIgnoreRuleSemantics` | 三枚 PASS | ✅ 同 |
| **B** 只 M1 | `--- PASS` F1、`--- PASS` 语义，新测试 **FAIL**，唯一一条 `scan_test.go:1633: holds() does not know a path the index plainly holds: ...` | FAIL 于 `:1633`，其余 PASS | ✅ 同（**只有那一枚白盒断言响**） |
| **C** 只 M4a | 新测试 **PASS**、F1 **PASS**、语义 **FAIL** 两支：`scan_test.go:1787: decide("lead/inside.txt", ...) = ignored true, want false` 与 `scan_test.go:1787: decide("  lead/inside.txt", ...) = ignored false, want true` | 同 | ✅ 同（**行为腿被全量那一味救回来**） |
| **D** M4a＋M1 | 新测试 **FAIL 四枚**：`scan_test.go:1606`（ban #6 counted 3 vs 3）、`:1609`（ban #8 同）、`:1619`（finding 没点名，`findings=[frontend/dist/tracked-or-not.tsx:1: [emoji] ...]`）、`:1633`（`holds()`）；语义测试 FAIL | 同四枚行号 | ✅ 同 |

⇒ **交件件那张表的四个读数、含全部行号，一枚不差地复现了。**这是本件第一次遇到"实现程的变异表我逐格重跑完全对得上"。

### 1.3 我把这一格往深里量了一层：变异改不改**门的可观察输出**（不只是改不改测试）

交件件只量到"测试红不红"。生死格的真正问法是题面那句：**"能不能造出一棵使 `sh scripts/d22scan.sh` 输出发生变化的树"**。
我于是拿四枚二进制（`bin/d22-{post,M1,M4a,M4aM1}.exe`，同一套 `scan_test.go` 之外的被测码，只差那一味）
去扫**同一棵树、同一串命令行**，做逐字节 `diff`。

真树（本仓工作树，`-root` 指过去）：

```
$ ./bin/d22-post.exe -root "D:/work/workspace/projects plans/Wisp" > real-post.txt   rc=0
$ ./bin/d22-M1.exe  -root "D:/work/workspace/projects plans/Wisp" > real-M1.txt     rc=0
$ diff real-post.txt real-M1.txt
IDENTICAL on the real repo tree (11 lines)
```

交件件那枚种子树（`seedtree\`：`.gitignore`＝`frontend/dist/*` + `!…gitkeep` + `build/`、
`frontend/.gitignore`＝`dist/*` / `!dist/.gitkeep` / `  weird/*`、`frontend/weird/inside.tsx` 含 U+2264、
普通 `git add -A` ＋ 只对 dist 那枚 `git add -f`；`git ls-files` 11 枚、`-i -c` 只有 dist 那枚）：

```
post   rc=1   M1 rc=1   M4a rc=1   M4aM1 rc=1
post == M1        （逐字节同）
post == M4a       （逐字节同）  <- 这一对交件件没量过，见下
M4a  != M4aM1     （逐字节不同）
M1   != M4aM1
$ diff out-M4a.txt out-M4aM1.txt
0a1
> d22scan: skipped as git-ignored: 1 file(s) under 0 ignored director(ies) [], decided by frontend/.gitignore (1 path(s))
5c6
< d22scan: scope ban #6 frontend/         examined   4 text files
> d22scan: scope ban #6 frontend/         examined   3 text files
8c9
< d22scan: scope ban #8 frontend/         examined   4 text files
> d22scan: scope ban #8 frontend/         examined   3 text files
12,13c13
< frontend/weird/inside.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned ...
< d22scan: 2 finding(s); D22 bans are not negotiable ...
---
> d22scan: 1 finding(s); D22 bans are not negotiable ...
```

**三条读数，一条比一条要紧**：
1. **交付态下摘掉那一味，门的可观察输出一字节都不动**（真树同、种子树同）。⇒ 读法 (ii) 的**前**提句，我复现了。
2. **`post == M4a`**：把 F3 的行首空白那枚语义修退回，门的输出仍然和交付态逐字相同——
   **原因就是那一味全量清单把被追踪字节接住了**（这正是交件件 C 行"新测试 PASS"的机理，我在门输出层看到了）。
3. **`M4a != M4aM1`**：同一发语义退回，再把那一味摘掉，**一枚受追踪的交付字节当场从分母与 finding 里消失**，
   并且 `note()` 从此自称"skipped as git-ignored"。⇒ 那一味的**行为后果是实的**，只是它显形的前提是"匹配器读宽"。

### 1.4 于是剩下的唯一问题：交付态的匹配器到底会不会读宽——我去造了，造不出来

M1 能改变门输出的**充要条件**是：存在一棵树、一个路径 P，使 **P 受追踪 ∧ git 的窄清单不含 P ∧ 我们的匹配器要跳过 P**
（即 `gitignore.go:380` 那一枚 `if ix.holds(rel,isDir) { return false }` 的答复在两版之间不同，且 `g.decide` 说 ignore）。
这就是"匹配器比 git 多跳"那一类。我不推理，直接搜：

- **手挑 47 形**（`fuzz\zz_overread_probe_test.go`）：F3 两枚实例的正反形、行首 1/2/3 空格与制表、行尾空格、
  行内空格、`**` 首中尾、`/`、`//`、`!/` 各种取反次序、`[a-z]`、`[!x]`、`?`、`{weird}`、`#` 首/中、
  反斜杠转义空格与结尾反斜杠、`\r`、大写形、真空白行、目录名真带两空格、路径带空格目录、`a/../weird/*`，
  文件位与目录位都问。**两把 git 尺同时问**（窄清单 `ls-files -i -c` ＋ `check-ignore -v` 并把 `!` 开头的模式判为"未忽略"）。

```
$ (cd fuzz && go test -count=1 -run TestProbeOverRead -v .)
CENSUS over=0 under=1 agree=54 total=55      -> PASS（我自己的探针，非交付件）
    under-read (safe) "upper" rule "WEIRD/*" path "frontend/weird/inside.tsx"
```

- **乘积 fuzz 384 形 × 4 条路径 = 1536 次比对**（`fuzz\zz_combo_probe_test.go`，
  行首空白 3 值 × 取反 2 × 词干 4 × 分隔 2 × 尾部 4 × 行尾空白 2，确定性枚举，不是随机）：

```
$ (cd fuzz && go test -count=1 -timeout 30m -run TestProbeOverReadCombos -v .)
COMBO CENSUS rules=384 over=0                 -> PASS
```

⇒ **合计 439 枚规则形（55 手挑 ＋ 384 乘积枚举；两批之间有几枚是同一形，如 `  weird/*` 在两处都放了，所以"互不相同的形"略少于 439），交付态下我造不出一枚"过度匹配"**。两件必须同时报出来的方法论账：
① 第一版我用 `git check-ignore`（不带 `--no-index`）当尺，被 3 形误判成 under-read——原因是**纯取反规则 `!weird/*` 也会被 check-ignore 以 rc=0 举出**，
而 `ls-files -i -c` 从不举目录；修正后的尺是"窄清单 ∧ 解析过 `!` 的 check-ignore"两问都答否。这一步不改变 over=0 的结论，但说明**这把尺本身需要正控**。
② 我的样本是 439 形，**不是全部形**；"今天没有过度匹配"是"我搜过 439 形没找到"，不是"没有"。

### 1.5 结构层：那一味买的是**一条不变式**，不是一次可观察输出

`skip()`（`gitignore.go:376-386`，现量）的顺序是：问不到 index ⇒ 一条规则都不应用；然后 **`holds()` 先答，答"是"就直接不跳**；
再轮到 `g.decide`。`holds()`（`:207-215`）文件位是 `tracked[rel] || ignored[rel]`，目录位只看 `dirs[rel]`。
于是全量那一味支撑的是 `gitignore.go:95-96` 那两句**自陈的不变式**：

> a path the index holds is never skipped, **whatever the rules say**, and a directory containing one is never pruned (F1).

摘掉 `all` 之后，这句话在盘上变成："**git 自己说撞了规则的**那批受追踪件永不被跳"。
两者的差集恰是 §1.4 那一类。而这一类**不是假想类**：M4a 一发就把它造出来（§1.3 第 2、3 条读数）。
⇒ 那一味的功能是**让"不隐身"这件事与规则形状无关**：不管将来长出第三枚什么形，只要 P 在索引里，它就扫得到。
窄清单那一味**结构上给不了**这个性质——`-i -c` 只能举 git 说"撞了"的件。

### 1.6 回到公式，并且说明为什么这里不能只听公式的

- 按公式：摘掉那一味 ⇒ 有没有一发变异从此不再红？**摘之前**（`304aeec` 之前）答案是"没有"——M1 全绿，
  这正是上一程验收与 `A221②` 记的那笔账；**摘之后**答案变成"有"——B 行那一枚 `:1633`。⇒ 在本仓这把尺下，**这一味现在是"被断言要求的成分"**。
- 但公式给的是"有没有断言要求它"，不是"它今天改变输出吗"。若只按后者判，
  §1.3 第 1 条读数（post≡M1）会把人引向 (ii)。**引向 (ii) 的那个推理链里藏着一个错步**：
  "今天没有可观察差异"≠"任何树上都没有"，它只等于"我手上这批树里没有过度匹配"。
  镜子（本仓用词＝装饰的钉子）的判据应当是"**任何可达输入上都不产生差别**"，
  而我量到的只是"交付态的 439 枚规则形上不产生差别"——**这两句话之间的距离就是那味保险的全部射程**。
  反证更硬：同一棵树上 `post == M4a`、`M4a != M4aM1`，即"这味在**存在过度匹配**的输入上确实决定门输出"是**量出来的**，不是推出来的。
- 另一点与"只改内部 struct 字段"那种误导性情形有关的关键差别：`holds()` **不是记账字段，是 `skip()` 的第一个分支**（`:380`）。
  它的答案今天在这两棵树上恰好通向同一个 return，是因为 `decide()` 恰好也说"不跳"；一旦 `decide()` 说"跳"，它就是唯一拦住那一句的东西（§1.3 第 2、3 条）。

### 1.7 能不能"因为错的理由为真"——量了，不能

- 第二参数是 `isDir`，**不是**交件件 §1.1 那句话没写但读者可能误会的别的什么：
  上一程 M5 摘的正是"`holds()` 的 `isDir` 支"，且本文件对文件位传 `false` 是对的。
- `holds(path,false)` 的真值有**两条来源**：`tracked`（全量）与 `ignored`（窄清单）。**只钉"必须为真"确实可能被窄清单蒙过去**。
  交件件在 `:1626` 钉 `ix.ok`、`:1629-1631` 钉 `ix.ignored[path]` 必须为否（`t.Fatal`）——
  即**先把窄清单那条来源关掉**，才让 `:1632` 的真值只能来自全量。这两枚守卫是这枚白盒钉成立的必要条件，不是装饰。
- 我又把种子里那两枚前导空格删掉（让 git 真的去匹配 `weird/*`，窄清单就会含它）来考这条守卫，两版二进制都跑：

```
$ (cd mut/SN3-head && go test -run '^TestFullTrackedListCoversWhatTheNarrowListCannot$' -v .)
    scan_test.go:1597: a plain `git add -A` must track frontend/weird/inside.tsx (git matches no rule on it), ls-files="...frontend/dist/tracked-or-not.tsx\nfrontend/src/panel.tsx\n..."   <- weird 那枚不见了
--- FAIL: TestFullTrackedListCoversWhatTheNarrowListCannot (0.52s)
$ (cd mut/SN3-M1 && go test ... 同一枚台件、M1 代码)
    scan_test.go:1597: 同一句
--- FAIL: TestFullTrackedListCoversWhatTheNarrowListCannot (0.58s)
```

⇒ 台件前提一塌，测试**先红在前提断言上并写明原因**（`:1591` 的等值与 `:1597` 的"必须被普通 add 追踪"），
**不会静默地把"错的理由"读成通过**。交付态与 M1 态同判，说明这条红不是来自被测码。
- 唯一还能说"钉得不够全"的一层：`:1632` 钉的是 `holds()` 的**文件位**；**目录位**（`dirs`，同样由全量喂）不在它的射程里。
  上一程的 M5（`:189` 把 `isDir` 支短路）已由 `:1485/:1488` 钉住 ⇒ 两半各自有主，**这一格不缺牙**，
  但如果有人要问"新测试单独钉住了全量的几分之几"，答案是**文件那一半**。

### 1.8 判（这一格）

**读法 (i) 成立，但要带一条明写的 (iii) 形限定；读法 (ii) 的"删掉那味"这一处置被我的读数正面否掉。**

- (i) 的**为真部分**：M1 是一枚真实回归类（同一棵树上 M4a 一发即可让它显形，量出来的是门输出的分母与 finding，
  不是 struct 字段）；`holds()` 是 `skip()` 的决策入口；这枚钉是**目前全树唯一**能把交付态与 M1 态分开的断言（B 行）。
- (iii) 的**为真部分（限定语）**：在**今天交付态的匹配器**下，439 枚规则形里我造不出一棵使门输出改变的树；
  所以"这味今天有没有行为后果"的答案是**没有**，它的价值全部在"下一枚匹配器读宽"那一方。
  **它不是"为一个不可能发生的类挡事"**——那一类已经发生过两次（F3 ①②），且一发单词元的退回就能再造（§1.3）。
- (ii) 的**为真部分**：它的前半句（可观察输出相同）我复现了。
  **它的结论（⇒ 窄清单那味无行为后果 ⇒ 该删）不成立**，理由＝§1.3 第 2 条 `post == M4a` 与 §1.5 那条不变式。
- **对 `A221②` 的账**：那格写"我前面批准保留那句现在有数了"。**按"防御、不是装饰"这个结论，数确实有了**
  （M4a／M4a＋M1 那对读数我逐枚复现）。但要说清它**是哪一种数**：它量的是"**存在**一个输入（带过度匹配的树）使这味决定门输出"，
  而**不是**"**当前**任何一棵真树使这味改变门的输出"（后者我今天专门去量了，答案是否）。
  ⇒ **留在判断里的那一块不是零，但也不是这张表**：它是"这一类会不会再长第三枚"的预报。
  这句要如实写回 `A221②`。

---

## 2. 那枚夹具到底是不是合成的（题面格 1）

`TestFullTrackedListCoversWhatTheNarrowListCannot` 的断言里点名了 `frontend/weird/inside.tsx` 与 `.gitignore`。
`frontend/**` 属另一位的地界，测试运行**绝不能**往那里写。我逐层查机制，再查盘。

**机制层**（每一层都在 HEAD 的文件上现读到）：

| 环节 | 现量 | 结论 |
|---|---|---|
| 树从哪来 | `liveFixture`（`scan_test.go:723-741`）第一行就是 `root := t.TempDir()` | 临时目录，不是仓库 |
| 谁写文件 | `seedFile(t, root, rel, content)`（`:15-24`）＝ `filepath.Join(root, rel)` ＋ `MkdirAll` ＋ `WriteFile` | 只能写在 root 之下 |
| 新测试内部 | `sed -n '1541,1641p' scan_test.go` 里 `grep -nE 'filepath.Abs｜os.Getwd｜"\\.\\./"｜os.Chdir｜WriteFile｜Remove'` ⇒ **空输出** | 没有一处绕过 root |
| git 台件 | `gitIndexFixture`（`:1411-1415`）＝ `git init -q -b main` ＋ `git add -A`，cwd＝root；`gitForceAdd`（`:1421`）＝ `git add -f --` | 建的是**临时目录自己的仓库** |
| 扫描入口 | `scanFixture(t, root)`（`:743`）＝ `scanWithStats(root)`，root 一直是那枚临时根 | 没有第二棵树 |
| `t.TempDir()` 落在哪 | 我用一枚极小的 Go 程序问 `os.TempDir()`：`C:\Users\swq\AppData\Local\Temp` | **在仓库外**，且不在任何 `.git` 之下 |

**helper 是否新造**（交件件 §1.2 的自陈）：`git show <ref>:tools/d22scan/scan_test.go | grep -c '^func '` ⇒
`304aeec^` **51** ／ `HEAD` **52**，＋1 恰是那枚新测试本身 ⇒ **包级 helper 一枚没新增**，复现成立。
"形状照抄上一枚测试的同名局部 `seed`"也对：`grep -n "seed := func"` ⇒ `:1440`（`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter` 内部）
＋ `:1569`（新测试内部），两枚都是函数内闭包。

**盘上层**（题面要的前后对照）：见 §0.4 —— 整道门（含 step 1 那一遍全量测试）跑完，
`git status --porcelain` 20 行 → 20 行、`diff` 空。我再单独按路径问一次：

```
$ git status --porcelain -- frontend .gitignore design
 （design 那 16 枚 D ＋ 4 枚 ?? ——全是 owner 自己的）
$ git ls-files --others --exclude-standard -- frontend
 （空）
```

⇒ **判：夹具是合成的，`frontend/**` 与根 `.gitignore` 一字节未写，真树零残渣。**
补一条我自己造的场景（交件件把它列在"没测"里）：见 §7.4，把 `TMP`/`TEMP` 指进"别人的仓库"里跑一遍——
结论是这一支**结构上打不开**（台件自己 `git init`，`rev-parse --show-prefix` 在任何临时根都返回空），
所以那条"没测"其实可以划掉。

---

## 3. `.gitignore` 那枚空白种子（题面格 2）

### 3.1 我自己在仓库外复现了 git 的三个读数（`probe\seed1`，照交件件 §1.1 那枚 seed 逐字重建，git 2.52.0.windows.1）

```
$ git ls-files
.gitignore cmd/wisp/main.go design/index.html frontend/.gitignore frontend/dist/tracked-or-not.tsx
frontend/src/panel.tsx frontend/weird/inside.tsx go.mod internal/ok/ok.go internal/tools/ok.go tools/d22scan/allowlist.txt
                                        <- 11 枚，含 frontend/weird/inside.tsx（普通 git add -A 就追踪）
$ git ls-files -i -c --exclude-standard
frontend/dist/tracked-or-not.tsx        <- 只此一枚
$ git check-ignore -v frontend/weird/inside.tsx ; echo rc=$?
rc=1                                    <- git 明说"我没忽略它"
$ cat -A frontend/.gitignore
dist/*$
!dist/.gitkeep$
  weird/*$                              <- 行首两枚空格，没有被 $ 前的任何 CR 混淆
```

⇒ 交件件 §1.3 那三行读数**逐字复现**（它那张表的枚数是一棵 4 枚的最小树，见 §5.2，关系同、枚数不同）。
"两枚前导空格的规则不被 git 匹配、但把空白削掉之后就成了我们匹配器的一条真规则"——**这一形在真 git 上成立**。

### 3.2 ubuntu 上的 git 会不会同意（题面点名要我推一下）

不是我猜：这台机上的 git 自带文档（`gitignore(5)`，PATTERN FORMAT 一节，
`C:\Users\swq\.qoder-cn\bin\git\mingw64\share\doc\git-doc\gitignore.html`）逐字是：

> A blank line matches no files, so it can serve as a separator for readability.
> A line starting with # serves as a comment. …
> **Trailing spaces are ignored unless they are quoted with backslash ("\").**
> An optional prefix "!" which negates the pattern; …
> It is not possible to re-include a file if a parent directory of that file is excluded.

⇒ 三点：①规范只把**行尾**空白列为被忽略，**行首**空白不在这条清单里 ⇒ 它是模式数据；
②这一节自 git 1.8.2 起就是这个形状（"trailing spaces are ignored"从未反向过）⇒ **跨平台同意是有据的，不是运气**。
⚠**runner 上那枚 git 的版本号我没有量过**（这台机不是那台机），所以这句话的力度止于"该规则在文档里存在且从未变向"，
不是"我在 ubuntu 上复现了"（另见 §7.2 与 §11-2）；
③同一节那句"parent directory 被排除时不能重新包含"倒是**我们主动偏离 git** 的一支——
我们的 `match()`（`gitignore.go:504-538`）从最深一层规则文件先答，`weird/` ＋ `!weird/inside.tsx` 会答"未忽略"，
git 答"已排除"。**方向＝多扫**（§1.4 那 439 形里 `neg-dir-then-file` 落在 agree，因为强制追踪后 `-i -c` 也没举它），
和 `gitignore.go:97-102` 那段自陈的"UNSUPPORTED ON PURPOSE … 只能让我们少跳"是同一族。

### 3.3 那枚"整串等值"断言稳不稳、塌的时候响不响

断言在 `scan_test.go:1591-1595`：`git ls-files -i -c --exclude-standard` 去空白后必须 **恰好等于**
`frontend/dist/tracked-or-not.tsx`，塌了是 `t.Fatalf` 且消息把两种塌法分开解释。

- **它比前提需要的更严**：前提只要求"weird 那枚不在里面"。所以任何一次 git 升级只要多报一枚**无关**的受追踪件，
  这条就会红——**代价是维护性（false alarm），不是牙齿**。它换来的是"0 读数自带正控"（dist 那枚必须在，否则就是仪器没在看这棵树），
  这个交换我认为是对的，符合本仓"0 命中必须配已知正控"的规矩。
- **塌得响不响，我量了**：把种子里那两枚前导空格删掉（＝git 开始真的匹配），交付态与 M1 态**都红**，
  且红在前提句上、消息点名原因（读数见 §1.7）。⇒ **fail-closed**，不会静默把"前提变了"读成"通过"。
- **另一支我也确认了**：`:1629-1631` 的第三枚 `t.Fatal` 在"窄清单也开始举这枚件"时点名"换形状"。
  两支合起来＝这条台件一旦不再钉住全量腿，**没有一条路能安静走过去**。

⇒ **判：那枚空白种子是真形、跨平台有文档依据、等值断言稳定度换成正控收益且塌时响亮。附条件：无。**

---

## 4. 行号与引用完整性（题面格 3）——**本件唯一判退回的一格，两枚都是注释级**

节内编号说明：题面把 8 枚格子列在 §2..§8、总裁列在 §9；我按"生死格＝§1、题面 8 格＝§2..§9"排，
所以题面的格 3 落在本件 §4、格 4 落在 §5，总裁在 §10、没测在 §11、给人看在 §12。

### 4.1 逐枚现量（全部在 HEAD `310816f`，另附漂移史）

| 引用（谁写的） | 声称指着什么 | 我量到的 | 判 |
|---|---|---|---|
| `gitignore.go:11` → `.gitignore:22` | 根 `.gitignore` 第 22 行＝`frontend/dist/*` | `sed -n '20,26p' .gitignore` ⇒ 22=`frontend/dist/*`、24=`assets/web/dist/` | ✅ 改对了（旧文 `:24` 指的确实是另一条规则） |
| `gitignore.go:11` → `frontend/.gitignore:12` | `dist/*` | 现量 12=`dist/*` | ✅ 未动且同值 |
| `gitignore.go:103-119`（新文） | 两支并列＋末行不声称穷尽 | 现读 103-119 正是新文，119 行＝`// This is a list of the shapes seen so far, not a claim that there are no others.` | ✅ |
| 旧 `:103-105` | 被换掉的那三行 | `git diff b7c06d2^..b7c06d2` 的 hunk 头 `@@ -100,9 +100,23 @@`，删的三行正是"The one shape that could make it skip MORE than git - not being able to reach the index at all - …" | ✅ 交件件"行号没漂"那句成立（在它落地那一刻） |
| `gitignore.go:180-194`（新文） | `-i -c` 那位的真理由 | 块起 180 ✓，但**被改的行是 183-194**，180-182 是三行未动的命令表 | ⚠ 松一格（块起点对、范围含未改行）；旧号同理：hunk `@@ -166,9 +180,18 @@`，实删的是旧 168-170，件里写"原 :166-171" |
| `gitignore.go:448` | `note()` 点名路径那一行 | 现量 447-450 是第三条 self-report，448 就是那句 `lines = append(... "TRACKED and matching an ignore rule" ...)` | ✅ |
| `gitignore.go:580`（M4a 落点） | `TrimRight(raw," \t\r")` | 现量 580 逐字是它 | ✅ |
| `gitignore.go:248`（M1 落点，§1.4 与台账 `A223③` 用的是这枚） | `parseIndexPaths(all)` | 现量 248 逐字是它 | ✅ |
| `gitignore.go:301 / :302 / :324 / :252` | `parseIndexPaths` 定义／它调 `parseIndexPathsList`／后者定义／窄清单调同一枚 | 301=`func parseIndexPaths(`、302=`list, err := parseIndexPathsList(z)`、324=`func parseIndexPathsList(`、252=`ignored, err := parseIndexPathsList(withIndex)` | ✅ 四枚全对 |
| 新测试内部 `:1568/:1591/:1596/:1606/:1609/:1619/:1626/:1629/:1632/:1633`、语义测试 `:1787` | 各自行上的那枚断言 | 逐枚现读同值；`:1787` 正是 `t.Errorf("decide(%q, isDir=%v) …")`，其 `lead/inside.txt` / `"  lead/inside.txt"` 两支在 1782-1783 | ✅ 11 枚全对 |
| 交件件 §1.1 说新测试＝`scan_test.go:1541-1640` | 那一枚测试的范围 | `git diff --unified=0 304aeec^..304aeec` ⇒ **只有一个 hunk `@@ -1540,0 +1541,101 @@`** ⇒ 实占 1541-1641（末行是分隔空行） | ⚠ 差一枚尾行，无实质影响；顺证：**101 行全在这一段里，没有任何夹带的第四处改动** |

### 4.2 退回枚 ①：新测试的注释里那枚 `gitignore.go:225` **在本批内部漂掉了**

`scan_test.go:1547` 写的是"removing it (parseIndexPaths(all) -> parseIndexPaths(withIndex), **gitignore.go:225**)"。
五锚现量：

```
ca84b75    call@225   line225=[	tracked, dirs, err := parseIndexPaths(all)]
43a6043    call@225   line225=[	tracked, dirs, err := parseIndexPaths(all)]
304aeec^   call@225   line225=[	tracked, dirs, err := parseIndexPaths(all)]
b7c06d2    call@248   line225=[func runGitIndex(root string) *gitIndexState {]
HEAD       call@248   line225=[func runGitIndex(root string) *gitIndexState {]
```

⇒ 这枚号**写下那一刻是对的**（`304aeec` 时 gitignore.go 还没动），
被**同一批的第二枚 commit `b7c06d2`**（＋23 行注释，全在该调用之上）漂到 248——
`225 + 14（②那一支）+ 9（③那一段）= 248` 我按 hunk 头算过，对得上。
而**同一批的证据件 §1.4 与台账 `A223③` 用的是 248** ⇒ 一批之内的两份文件对同一枚变异给了两个号，
码里那枚是过期的。今天 `:225` 落在 `runGitIndex` 的函数声明行（同一个函数、不是那句调用）。
**为什么这在本仓算一格而不是笔误**：条件②整枚条件就是"自我描述的注释不许过保"，
上一程 panel 那批也因 `:604` 指着一枚不存在的测试被判退回；本枚是同一族（号漂了、内容还找得回）。
**最小修法（点名，不动手）**：`scan_test.go:1547` 里把号改成 `gitignore.go:248`，
或按我这轮的教训**干脆不写裸行号**、改写"`runGitIndex` 里 `parseIndexPaths(all)` 那一句"——
这批自己已经演示了行号在一枚 commit 之内就会烂。

### 4.3 退回枚 ②（轻）：`3bb99aa` 那笔归属错了

交件件 §2 末段："`scan_test.go:1179` 与 `:1194` 里还有两句 'the one shape …' 式措辞，
属**上一批**（`3bb99aa`）写的注释"。现量：

```
$ git blame -L 1179,1179 -L 1194,1194 --porcelain tools/d22scan/scan_test.go
ca84b75  1179  1179  1     summary: fix(d22scan,A214/F1): skip() 改问 git index ...
ca84b75  1194  1194  1
```

⇒ 两句都是 **`ca84b75`**（＝被上一份验收裁的那一批）写的，不是 `3bb99aa`（再上一批），错了一代。
内容旁证也对：那两句讲的是"A214/F1 EXTENDS THAT DUTY TO THE INDEX, IN THE SAME COMMIT AS THE CODE"，
`ca84b75` 才是"同批改副本"那一枚。
**它真正的那句话（不在验收点名的三件之内、本件一字未动）我核过＝成立**（`:1179`/`:1194` 在 `git diff` 里没出现）。
⇒ 记**登记性缺陷**：把别人的账记到第三个人头上，会直接影响下一程"该谁修"的分派。

### 4.4 另两枚没到退回、但必须留名的数

1. **`§1.4` 的"首轮量在 :245"在任何一枚已提交锚点上都复现不出来**（提交态只有 225 与 248 两个号）。
   它必然是某一刻工作树的号 ⇒ 与本仓"`归因腐坏`"那族同形（一个数没带它是在哪棵/哪一刻量的）。
2. **`§3.1` 那张表的枚数（full=4、narrow=1、`comm -23` 三行）是一棵 4 枚的最小树**，
   不是交付测试那棵 `liveFixture` 树。我在**照 §1.1 逐字重建的 11 枚树**上重量（见 §5.1）：
   full=11、narrow=1、`comm -23`＝10 行。**关系两遍同向**（子集为空／反向非空），但**枚数不可复算**——
   下一位若想复算 §3.1 那三行，得先知道那是一棵最小树而不是台件那棵。

### 4.5 判

**这一枚格子＝退回（措辞级，两枚具名：`scan_test.go:1547` 的 `:225`、交件件 §2 末段的 `3bb99aa`）。**
其余 20 余枚 `file:line` 引用**逐枚现量同值**（含 §4.1 那 11 枚测试行号与 4 枚解析器行号），
两枚新写的散文段整体**没有虚高**（判据见 §5.4）；上面 4.4 两枚是登记性提醒，不绑修。
**本件没有一处是我"读着像对"就放过的**：每枚号我都落到行上看了一遍。

---

## 5. 子集／并集那次重导（题面格 4）

### 5.1 我自己重跑了一遍，两把尺各一枚（`probe\seed1`＝§3.1 那棵 11 枚树）

```
$ git ls-files -z | tr '\0' '\n' | sort > full.lst        # full=11
$ git ls-files -i -c --exclude-standard | sort > narrow.lst  # narrow=1
$ comm -13 full.lst narrow.lst
（空）                                    <- 窄 ⊆ 全，一枚新的都带不进来
$ comm -23 full.lst narrow.lst | tee c23.txt
.gitignore / cmd/wisp/main.go / design/index.html / frontend/.gitignore / frontend/src/panel.tsx
frontend/weird/inside.tsx / go.mod / internal/ok/ok.go / internal/tools/ok.go / tools/d22scan/allowlist.txt
count=10                                  <- 反方向正控非空（两清单不等价）
```

⇒ 交件件 §3.1 那对的**方向**复现成立（空／非空各一枚），枚数按 4.4-2 那笔账换成我这棵树的数。

### 5.2 "共用一枚解析器"这条理由，是把子集说成免费的还是把它钉住的

两问要拆开：
- **子集关系本身**是 **git 的语义**（`-i -c` 举的是"受追踪 ∧ 撞规则"），跟我们的解析器无关 ⇒
  在**静止的索引**上它确实近乎免费。我在注释里看到 `:302` 与 `:252` 都叫同一枚 `parseIndexPathsList`（现量在 §4.1），
  这条的作用**不是**证子集，而是**把原句那句"so one command's parsing bug cannot blind the other"打死**——
  同一个解析器 ⇒ 一次解析 bug 同打两枚，两枚从来不是彼此的后盾。**这一步是承重的**：
  被换掉的那句如果留着，就是一个假的安全属性（"两路互为备份"），下一程会照着它省掉一路。
- 所以我的判是：**"子集"这一腿近乎免费，"共用解析器 ⇒ 原来那句不成立"这一腿承重**，
  而它俩被写进同一段注释里，各自都对着。这段注释**没有**把免费的当昂贵的吹（它只说 "not a second opinion inside holds()"），
  也没把昂贵的当免费的略。

### 5.3 但我量出这句新注释**说过头了一格**：那两枚读不是原子的

`runGitIndex` 是**两次独立 spawn**：`:240` 读 `ls-files -z`，`:244` 读 `ls-files -z -i -c`。
中间若有一枚 commit 落地，窄清单就能含一枚全量清单没见过的路径 ⇒ **`its output is strictly a subset of ls-files -z`
这句无条件的话在 wire 层面不保证**。我在仓库外把这枚态**白盒造出来量了**（`fuzz\zz_race_probe_test.go`，
手写一枚 `gitIndexState{ok:true, tracked:{}, dirs:{}, ignored:{file:true}}` 塞进 matcher）：

```
$ go test -count=1 -v -run TestProbeRaceDirLeg .
    holds(file,false)=true holds(dir,true)=false
    raced: file NOT skipped (holds() consults ignored for files) - safe
    raced: DIRECTORY pruned -> holds(dir,true) reads only dirs, which the narrow list never feeds
--- PASS
```

⇒ 分两支：**文件支仍安全**（`holds()` 文件位读 `tracked || ignored`，窄清单能兜住）；
**目录支不兜**（`holds(:211-213)` 目录位只看 `dirs`，而 `dirs` 只由全量喂）⇒ 那一枚竞态下整目录被 `SkipDir`，
窄清单明明知道里面有受追踪件也救不回来，方向＝**少扫**（危险那一侧，`gitignore.go:103` 自己定义的那一侧）。
同一枚竞态还捎上下一句："`Nothing in a skip decision comes from this list that did not already come from the guard`"——
静止索引上对，竞态下那枚文件位就是"从这张清单多出来的"（且是好事）。

**这条不是本批造的行为**：两枚 spawn 是 `ca84b75` 的形状，交件件 §6-4 也明写了"并发 commit 中途读 index 的一致性"没测。
**本批写的是那句无条件的话**。按这批自己立的尺（`:119` "a list of the shapes seen so far, not a claim there are no others"），
`:186-188` 那两句缺同一枚限定。⇒ **最小修法（点名，不动手）**：把那两句加上"on a quiescent index"
／或改说"between the two reads nothing may commit"。**我没有**把这枚竞态在真 git 上跑出来（见 §11），
所以它是"码上可查＋白盒已造态"级，不是"端到端复现"级。

### 5.4 判

**附条件成立**：重导方向真、两把尺各一枚、"共用解析器"那一步确实在打死一句假担保（承重）；
两处 `strictly a subset` / `nothing … that did not already come from the guard` 的措辞过头（§5.3），
与本批被判的两枚行号账同属注释级，一起修一程即可。

---

## 6. 闸门与那几把尺（题面格 5）——每枚数旁边都写我量它的锚点

| 判据（交件件的声明） | 我的现量 | 锚点／时刻 | 判 |
|---|---|---|---|
| step 1 `PASS=28→29 FAIL=0 SKIP=0 RUN=68→69` | **三处各量一遍**：① 真树整道门＝`29/0/0/69 rc=0`（§0.1）；② `fullhead\`＝`git archive HEAD` 全量快照里 `go test -count=1 -v ./...`＝**29 / 0 / 0 / 69**；③ `fullpre\`＝`git archive 304aeec^` 快照同命令＝**28 / 0 / 0 / 68**（`SKIP=0`！） | `310816f`；快照＝HEAD 与 `22be55e` | ✅ **复现，而且我把"28 那一枚"从〔仅自述〕提到〔独立复现〕**：交件件是在活树上前后来各跑一次量的，我是拿两棵**已提交快照**在仓库外各跑一次，中间没有"我自己的改动"这个混杂量 |
| 名册差集 ＋1／−0（28 枚旧名逐枚仍在） | `diff` 两棵快照的 `--- PASS/SKIP` 名册 ⇒ 唯一一行 `> --- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot`，无删除行；再与真树名册对（29 枚）⇒ 差的 6 枚在快照里是 SKIP（它们在副本里找不到 `../..` 那棵真仓库），**一枚没消失** | `fullpre` → `fullhead` | ✅（这一条按台账教训必做：一枚 panic 会吞掉同包其余读数，四数看不出来） |
| "新测试内部没有 `t.Run`，所以两数各 ＋1" | `sed -n '1541,1641p' scan_test.go | grep -c 't\.Run('` ＝ **0** | HEAD | ✅ 解释了为什么不是 ＋2 |
| step 2 八数 `203/22/40/18 · 32/40/407/39`、`rc=0` | §0.2 逐字，**八枚同值**；`ban #8 internal/=407` 的归属＝`internal/panel/l2_grant_boundary_test.go`（盘上 10:09 那版），**不是本批**（本批对 `internal/**` 零动作，见 §9） | `310816f` | ✅（并注明：`design/=32` 是工作树形状，净快照 30——这一对我没法在活树上分开量，我引 `A221⑥`/验收 §9 那两笔已登记的读数，不当自己的证据用） |
| `gofmt -l` 空 | `tools/d22scan` 空；**整仓 `gofmt -l .` 也空** | HEAD | ✅（我把它扩到全仓量的） |
| `go vet ./...`／`go build ./...` 过 | 都 rc=0，`vet.log` 行数 0；`tools/d22scan` 内单独 `go vet ./...` 也 rc=0 | HEAD；⚠ 整树 vet 会编译 `internal/panel`，那是并行程的提交态，**我只据 rc=0 说"编译干净"，不拿它当行为读数** | ✅ |
| 断言 `164→173`（＋9，逐枚在新测试里）、删除列 0 | 我这把尺 `grep -cE 't\.(Error|Fatal)[A-Za-z]*\('` 在 `git show` 上量：`304aeec^`=**164**／`HEAD`=**173**；`sed -n '1541,1641p' | grep -c` 同一把尺＝**9**（⇒ ＋9 全在那一枚内，一枚没在别处加）；`git diff --numstat 304aeec^..b7c06d2 -- scan_test.go`＝`101  0`（删除列空） | HEAD | ✅ |
| `t.Skip[f]?(` 7→7，且**第一把尺恒不匹配已作废重做** | 正控我两问：① 字面 `t\.Skip(` 在整个 `tools/d22scan/*.go` 里 **0 命中** ⇒ 第一把尺不是"读到 0"，是**结构性不可能命中**（本包 7 枚全是 `t.Skipf(`）；② 修好的尺打印出那 7 枚行号 `241,261,677,1044,1189,1916,1989`，**与交件件列的七枚逐枚同值** | HEAD | ✅ 这一格我原本打算抓它一把，结果抓到的是它自己已经抓过的那把 |
| `emojiRe` 那一行四锚同值 | 我量了**六锚**（`3bb99aa`/`ca84b75`/`43a6043`/`304aeec^`/`b7c06d2`/`HEAD`）：`git show $r:tools/d22scan/main.go | grep '^var emojiRe' | git hash-object --stdin` ⇒ 六枚全是 `7355202a062eed144b52bdd81a386800b52643ae` | 六锚 | ✅ **一字未动**，且 hash 与件里写的那枚逐字符同 |
| 条件②③"生产码零字节改动" | `git diff b7c06d2^..b7c06d2 -- gitignore.go | grep -E '^[+-]' | grep -v '^+++|^---' | grep -vE '^[+-][[:space:]]*//'` ⇒ **空**；**正控**（同一把尺打在真代码差 `304aeec^..304aeec -- scan_test.go` 上）⇒ **63 行**，尺是活的；再把两版**整行注释剥光**做 diff ⇒ **481 行 vs 481 行、空输出** | HEAD | ✅ 三把尺各一枚，"注释级"这句话在盘上站得住 |

---

## 7. 那枚新 spawn-git 的测试稳不稳（题面格 6）

### 7.1 `-count=2`

```
$ go test -run '^TestFullTrackedListCoversWhatTheNarrowListCannot$' -count=2 -v ./...
=== RUN   TestFullTrackedListCoversWhatTheNarrowListCannot
--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot (1.02s)
=== RUN   TestFullTrackedListCoversWhatTheNarrowListCannot
--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot (0.97s)
ok  github.com/CarlosShao/wisp/tools/d22scan  2.032s   rc=0
```

⇒ 两遍同值、不依赖上一遍留下的东西（每遍自己 `t.TempDir()` ＋ 自己 `git init`）。多次复跑的耗时区间 0.51-1.44 s。

### 7.2 它到底吃不吃 git 的身份与全局配置

静态先看清它 spawn 了什么：`git init -q -b main`／`git add -A`／`git add -f --`／`git ls-files -i -c --exclude-standard`／`git ls-files`
——**没有 `commit`**，所以 `user.name`/`user.email` 在结构上就用不到。我照"CI 可能什么都没有"再实测一遍：

```
# 全局＋系统配置各指到一枚空文件、HOME/USERPROFILE/XDG_CONFIG_HOME 指到不存在的目录
$ GIT_CONFIG_GLOBAL=... GIT_CONFIG_SYSTEM=... HOME=/nonexistent-xyz ... go test -run '^TestFull...$' -v .
--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot (0.96s)      # 以及空配置文件的第二遍 (0.99s) PASS
```

⇒ **不吃身份、不吃全局配置**。`-b main` 需要 git ≥ 2.28（2020-06），本机 `git version 2.52.0.windows.1`。
CI 那一侧我**只量到"这道门整道跑在 ubuntu-latest"**（现读 `.github/workflows/ci.yml`：`:66` `runs-on: ubuntu-latest`、
`:68` `actions/checkout@v4`、`:81` 步 1＝`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`、`:109` 步 2＝`sh scripts/d22scan.sh`），
**那台 runner 上的 git 版本号我没有量过**（本机不是那台机，§11-1/§11-2）⇒ 只能说"任何现代 git 都支持 `-b`"，不能说"CI 上实测通过"。
⇒ 附带一条要紧的：**下一推这枚新测试就会在 ubuntu 上跑**，所以交件件 §6-2 那句"没在 Linux 上跑过"是**仍然开着的一格**，
不能由它自己或我这一程关掉。

**唯一真吃的全局项我自己造出来了**：`--exclude-standard` 的射程**包含 `core.excludesFile`**。
我在仓库外挂了一枚全局 excludes（内容一行 `inside.tsx`）再跑：

```
$ git check-ignore --no-index -v frontend/weird/inside.tsx        # 在带那枚全局 excludes 的环境里
D:/tmp/d22scan-close-r1-accept/mygitignore:1:inside.tsx	frontend/weird/inside.tsx
（上一条命令的 rc=0 ⇒ git 认为这枚件"被规则 bind 住了"，来源是那枚全局文件）
$ GIT_CONFIG_GLOBAL=<that config> go test -run '^TestFull...$' -v .
    scan_test.go:1597: a plain `git add -A` must track frontend/weird/inside.tsx (git matches no rule on it), ls-files="..."
--- FAIL: TestFullTrackedListCoversWhatTheNarrowListCannot (0.55s)
$ 撤掉那枚配置再跑同一枚测试 -> ok  github.com/CarlosShao/wisp/tools/d22scan 0.890s
```

⇒ **机器配置能把台件前提撬动，但撬动时它红在前提句上（`t.Fatalf`，还点名 ls-files 现场），不会静默通过。**
方向＝fail-closed，符合门的规矩；代价＝若哪天真给 runner 挂一枚撞名的全局 excludes，这枚测试要人来看一眼才知道是环境不是代码。
（另一枚我自己踩过的坑记在这里免得下位重踩：全局 excludes 写 `weird/*` 是**撬不动**的——
`--exclude-standard` 里的模式带 `/` 就是相对仓库顶锚定的，而临时仓库的顶下面没有 `weird/`；
换成不带斜杠的 `inside.tsx` 才命中。这是我第一遍"没测出来"的原因，不是测试的漏洞。）

### 7.3 它和 `A221③` 那支"没有 git ⇒ 全扫并响亮自陈"的关系

```
$ PATH="/d/work/base/go/bin:/c/Windows/System32:/c/Windows" go test -run '^TestFull...$' -v .   # PATH 里没有 git
    scan_test.go:1586: git init -q -b main in C:\Users\...\TestFullTrackedListCoversWhatTheNarrowListCannot203388652\001:
                       exec: "git": executable file not found in %PATH% ((no stderr)) -
                       the tracked-path guarantees of A214 cannot be built without a git binary
--- FAIL: TestFullTrackedListCoversWhatTheNarrowListCannot (0.01s)
```

⇒ 缺 git 时它**红**（`t.Fatalf`），不 Skip ⇒ `runtests.sh` 的 `SKIP != 0 即退出 1` 与 `PASS=0&&FAIL=0 即退出 1` 都还有牙。
**门的意思变了吗？没有换方向、只多了一枚实例**：本包里今天共 **3 枚**测试直接 spawn git（现量：按"函数体内出现过
`gitCommand`/`gitIndexFixture`/`gitForceAdd`"归并到所属 `func Test`）
⇒ `TestWalksSkipGitIgnoredPaths`／`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter`／本枚。前两枚的归属我用 blame 钉了：

```
$ git blame -L 1279,1279 tools/d22scan/scan_test.go   -> 3bb99aa   （TestWalksSkipGitIgnoredPaths 的行）
$ git blame -L 1439,1439 tools/d22scan/scan_test.go   -> ca84b75   （TestTrackedPaths... 的行）
```

⇒ **"正控需要 git"这条环境依赖是上一批登记的**（验收 §8 那格），本批**没有新造一类**，只把实例数从 2 加到 3。
（口径写清：还有一枚 `TestBuiltBinaryGoesRedEndToEnd` 会把扫描器**建成二进制**去跑，那枚二进制内部也会问 git——
那是**间接**依赖，不在上面这 3 枚里，我也没为它单独归账。）

### 7.4 顺手划掉交件件 §6-2 那条"没测"

件里写："`t.TempDir()` 落点（TMPDIR 在真仓库里时会走子目录那一支）未验"。我实测：

```
# 把 TMP/TEMP 指到一枚真仓库的子目录里（D:/tmp/.../outer 已 git init，outer/inner 是它子目录）
$ TMP=... TEMP=... go test -run '^TestFull...$' -v .
--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot (0.65s)
$ git -C outer/inner/probe rev-parse --show-prefix  -> inner/probe/     # 说明这一支的守卫本身是活的
```

⇒ **它不是"没验所以不知道"，而是结构上打不开**：台件自己在临时根里 `git init`，
所以从临时根问 `rev-parse --show-prefix` 永远是空串（自己是仓库顶），那枚"子目录"分支（`gitignore.go:233-239`）
不可能被 `t.TempDir()` 的位置触发。**这一支可以从"没测"清单里划掉。**
（另注：Go 在 Windows 上认 `TMP`/`TEMP`、**不认** `TMPDIR`——我先用 `TMPDIR` 试时它根本没挪窝，`os.TempDir()` 那枚小工具量出来的。）

### 7.5 判

**成立**：`-count=2` 稳、不吃身份与全局配置、缺 git 响亮红、不新增门的环境依赖；
外加两枚我给的信息：`core.excludesFile` 是它唯一真吃的机器配置（红得响亮）、§7.4 那条"没测"可以结案。

---

## 8. 程序格：五枚 commit 与"顶回上游措辞还自批不必停"这笔账（题面格 7）

### 8.1 逐枚 name-only（我现跑，`git show --name-only --pretty=format:`）

| commit | 它带的路径 | 枚数 | 判 |
|---|---|---|---|
| `304aeec` | `docs/evidence/s1/d22scan-gitignore-fix-close-r1.md` ＋ `tools/d22scan/scan_test.go` | 2 | ✅ 全是自己的 |
| `b7c06d2` | 同上，代码那枚换 `tools/d22scan/gitignore.go` | 2 | ✅ |
| `cec5e78` | 只有那枚证据件 | 1 | ✅ |
| `38f340b` | 只有那枚证据件 | 1 | ✅ |
| `6a5b321` | 只有那枚证据件 | 1 | ✅ |

另核：五枚**全是单父、线性**（`git rev-list --parents -n1` 逐枚＝2 个词）⇒ 没有 merge、
没有"改写已推送历史"的形状；整窗口（含并行的 panel 三枚）`git diff --numstat 304aeec^..6a5b321` **只有 4 枚路径**（§9.1）
⇒ **owner 那 16 枚 `design/**` 未提交删除、那 4 枚未追踪件，一枚都没被卷进这五枚 commit**。
这是"带 pathspec"最硬的一条旁证：`git commit -- <两枚路径>` 在形状上就吞不掉别家。

**我不能核的那一半，明写**：每枚 commit 之前那一刻 `git diff --cached --name-only` 里有没有出现过别家路径——
索引历史不落盘，事后无法复原。所以我**只核得到"零泄漏"，核不到"有没有触发 `A220④(a)` 那条停手条件"**；
交件件 §5 那行自称"枚枚只带我自己的路径"，与盘上一致，但它自述的是**索引态**，那一层 remains 〔仅自述〕。

### 8.2 题面那句"自批不报告"，盘上比它说的好——先把事实摆正

题面的表述是"contradict the previous acceptance's sentence, then self-approve **not reporting it**"。我查了发布物本体：

```
$ git log -1 --format=%B 304aeec | grep -n "不同意"
14:一处不同意验收的文字：它 §4 末"从此摘掉 all 那一味必红"与它自己读数表第 4 行矛盾，盘上复算同它表一致
$ git show 304aeec:docs/evidence/s1/d22scan-gitignore-fix-close-r1.md | sed -n '77,80p'
⚠ **一处对验收文字的不同意（照单收会收错）**：验收 §4 末句写"从此摘掉 `all` 那一味必红"，
但它自己的读数表第 4 行是"只 M1（语义修完好）⇒ PASS"。盘上复算同它第 4 行一致：**纯行为的断言（分母／finding／rc）在 B 行全都不响**，
只有 `holds()` 那一枚白盒断言咬得住 M1 单独一发。⇒ 我在它点名的行为台件之外**多加了那一枚断言**，
否则这条闭不上；这不是扩大范围，是它那张表已经量到、但文字写歪了的那一位。
```

⇒ 偏离与理由**写进了那枚 commit 的正文与它同时落盘的证据件 §1.4**（不是交件后补、不是只写在聊天里）。
所以"没报告"不成立；成立的是**"报告了但没停下来等批准"**。这两件事在该归的账上差一档，我要按盘上那一档记。

### 8.3 判：**方向对、程序欠一手，但真正该补的不是这一程的账，是一条没有规则的空白**

- **实质那一半我完全认**（并且我把它独立复现到了二进制输出层，§1.3）。
  照字面执行"必红"只有两条路：要么在测试里**顺手把匹配器改坏**再造红——那是本仓最忌讳的"为凑判制造形状"；
  要么承认打不红然后把那味删掉——§1.3 第 2 条读数证明删了就真丢牙。两害相比，**加一枚最小的真到达它的白盒断言是唯一正确答案**。
- **程序那一半确实欠一次报回**：`A220④(a)` 的升格句式管的是"别家暂存条目"，本批**没有触发它**（我也无法证否）；
  本批触发的是**另一件事**——上游工作单内部两句话打架，实现方自选一支继续走。
  而**现行规矩里没有这一条**：AGENTS.md §0.1/§2 的"未定义即停"讲的是"方案没覆盖的情况"，
  没说"方案自相矛盾时算不算"。⇒ 这一程不该被记成"违反了 `A220④(a)`"（那会把账记到错的格子里），
  该记的是**"这里有一类没写下来的情形，本程用一次个人判断填了"**。
- **我建议给编排者一条可复算的判据**（把这次的正当条件写成规则，而不是把这次当先例供人自由引用）：
  下游顶回上游措辞**只在同时满足三件时才允许不亦步亦趋地照做、且不必停手**——
  ① 顶回依据是**盘上可复算**的读数（本批：变异表；我这轮：全复现＋二进制 diff）；
  ② 偏离与理由**写进造成它的那枚 commit**（本批：满足，见 §8.2）；
  ③ **不放宽任何既有断言**（本批：删除列 0、`SKIP=0`、阈值/golden/allowlist 零字节，见 §9）。
  三条缺一即回到"停手报回"。反过来说：**这三条不写下来，下一程就只会引用这次的结论，不会引用这次的边界。**
- 一处小瑕疵（不判退回）：`git log --oneline` 只看**标题**是看不出这次顶回的（标题只说"从此有断言要求它"），
  要看正文才看得见。若编排者希望这类偏离在 `log --oneline` 一层就可见，那是**措辞习惯**问题，可提一句。

---

## 9. 契约轴（题面格 8）

区间＝本批全部五枚 `304aeec^..6a5b321`（＝`22be55e..6a5b321`）。**整窗口改了什么的完整清单**：

```
$ git diff --numstat 304aeec^..6a5b321
267	0	docs/evidence/s1/d22scan-gitignore-fix-close-r1.md
382	0	docs/evidence/s1/panel-l2-grant-nail-accept-r1.md   <- 并行那位的证据件，不算我的账
30	7	tools/d22scan/gitignore.go
101	0	tools/d22scan/scan_test.go
```

| 契约面 | 我的量法 | 读数 | 判 |
|---|---|---|---|
| `internal/observe/thresholds.go` | `git diff --numstat 304aeec^..6a5b321 -- <path>` ＋ blob 比对 `304aeec^` vs `6a5b321` vs `HEAD` | 空输出；blob **同值** | ✅ 零字节 |
| golden（`*.sse`） | `git ls-files | grep -c '\.sse$'` ＋ `git diff --numstat … -- '*.sse'` | 仓里 **52 枚**；diff **空** | ✅ 一枚未动 |
| `tools/d22scan/allowlist.txt` | numstat ＋ blob | 空；blob 同值 | ✅ 零新增豁免 |
| `internal/risk/rules_gateway.go` | numstat ＋ blob | 空；blob 同值 | ✅ |
| `internal/risk/**` | `git diff --numstat … -- internal/risk` | 空 | ✅ |
| `docs/PLAN.md` ／ `docs/specs/**` | 各一枚 numstat | 空／空 | ✅ 未碰契约与规格 |
| `frontend/**` | `git diff --numstat … -- frontend` | 空 | ✅（工作树里也没有：§2 那两问） |
| `design/**` | `git diff --numstat … -- design` | 空 ⇒ owner 那 16 枚删除**没被卷进任何一枚 commit** | ✅ |
| `main.go`（`emojiRe` 所在文件） | blob `304aeec^` vs `HEAD` | 同值 | ✅ 仪器未动 |
| `tools/d22scan/gitignore.go` | blob 比对 | `DIFF`（＋30/−7，且 §6 末行那三把尺证明**全在注释行**） | ✅ 生产**行为**零改动 |

⇒ **题面列的九面全零字节**（阈值／golden／allowlist／`rules_gateway.go`／`internal/risk/**`／`PLAN.md`／`specs/**`／`frontend/**`／`design/**`），
**另加一面**：仪器所在的 `main.go` 也 blob 同值。唯一 DIFF 的是 `gitignore.go`，而它每一行都是注释（§6 末行）。
`DEFERRED(D-xx)` 那一条我一并扫了：`git diff 304aeec^..6a5b321 | grep -c 'DEFERRED('` ＝ **0** ⇒ 本批没新增/删除任何推迟标记（登记表的 1:1 双向不受影响）。

---

## 10. 总裁

| 格 | 判 | 一句凭据（全是我自己现量的，锚点写在正文里） |
|---|---|---|
| §0 闸门重跑 | **成立** | step1 `29/0/0/69 rc=0`、step2 八数逐枚同值、`gofmt`(整仓)/`vet`/`build` 全清、测试前后 `git status` 逐字节同值 |
| §1 生死格（白盒钉＝承重还是镜子） | **成立**（读法 (i)，带一枚明写的 (iii) 形限定；**读法 (ii) 的"删掉那味"被我的读数正面否掉**） | 四行变异表连行号逐枚复现；我多量一层到二进制：`post≡M1`（真树＋种子树）但 `post≡M4a` 而 `M4a≠M4aM1`（分母 4→3、finding 消失、note 改口）；439 枚规则形里造不出过度匹配 ⇒ "今天无可观察差异"为真、"那味无行为后果"为假 |
| §2 夹具是否合成 | **成立** | `liveFixture`→`t.TempDir`→`os.TempDir()` 实测在仓库外；新测试区零 `Abs/Getwd/Chdir`；包级 helper 51→52；按路径问 status 只余 owner 那 20 枚 |
| §3 空白种子 | **成立** | `git ls-files` / `-i -c` / `check-ignore` 三读数逐字复现；`gitignore(5)` 原文只豁免**行尾**空白 ⇒ 跨平台有据；等值断言塌时红在 `:1591`/`:1597`（实测），fail-closed |
| §4 行号与引用完整性 | **退回（措辞级，两枚具名）** | `scan_test.go:1547` 的 `gitignore.go:225` 被**本批自己的** `b7c06d2` 漂成 `:248`（五锚现量），同批证据件用 248 ⇒ 两份文件两个号；交件件 §2 把 `scan_test.go:1179/1194` 归给 `3bb99aa`，`git blame` ＝ **`ca84b75`**。其余 20 余枚 `file:line` 逐枚同值 |
| §5 子集／并集重导 | **附条件成立** | 方向复现（`comm -13` 空、`comm -23` 非空 10 行）；"共用解析器"那一腿**承重**（它打死的是"两路互为备份"那句假担保）；但新写的 `strictly a subset` / `nothing … already came from the guard` 是无条件话，而 `:240`/`:244` 是**两枚非原子的 spawn**——白盒造竞态量得出来：`holds(file)=true` 但 `holds(dir)=false` ⇒ 目录支不兜、方向＝少扫。补一句限定即闭合 |
| §6 数与尺 | **成立** | 每一枚我另配一把尺＋正控；其中"28 那一枚"我从〔仅自述〕提到〔独立复现〕（两棵 `git archive` 快照在仓库外各跑一次） |
| §7 新测试稳度 | **成立** | `-count=2` 两遍同过；不吃身份/全局配置（空配置＋假 HOME 实测）；唯一真吃的 `core.excludesFile` 我造出来了 ⇒ 红在前提句；缺 git ⇒ 红在 `:1586` 不 Skip；本包 git 依赖测试 2→3 枚 ⇒ 门未新造一类依赖；§7.4 顺手把它一条"没测"结案 |
| §8 程序 | **成立（实质）＋程序那半按下面这条改写** | 五枚 commit 逐枚 name-only 只有自己的路径、单父线性、owner 那 16 枚删除零卷入；**"自批不报告"这句按盘上要降级为"报告了但没停下来等批准"**（偏离与理由写在 `304aeec` 的正文与同枚落盘的证据件 §1.4 里） |
| §9 契约轴 | **成立** | 九面零字节；`*.sse` 52 枚未动；`DEFERRED(` 本区间 0 命中；两枚允许动的面各有尺 |

**总体＝退回（措辞级），三件最小闭合集合，全部注释／文档级，零生产行为改动：**

1. `tools/d22scan/scan_test.go:1547`：`gitignore.go:225` → `:248`，**或**去掉裸行号、改写"`runGitIndex` 里 `parseIndexPaths(all)` 那一句"（这批自己已经演示了行号在一枚 commit 内就会烂）。
2. `docs/evidence/s1/d22scan-gitignore-fix-close-r1.md` §2 末段：`3bb99aa` → `ca84b75`（**只追加更正、不改写原句**，按台账老规矩）。
3. `tools/d22scan/gitignore.go:186-188`：给 `strictly a subset`／`Nothing in a skip decision …` 那两句加上"同一枚静止索引上"这一限定，
   并（可选、更大的一手）记一句"`:240` 与 `:244` 是两枚独立 spawn，中间不容 commit"。这一条**不是本批造的行为**，是 `ca84b75` 的两读形状；
   要真补牙就是生产码改动 ⇒ **该单开一格由编排者拍，不在这一程顺手做**。

**三件之外，本批交件质量我按"上一程那三枚条件"逐条对齐**：条件①（补第三枚半）**闭合并且比我预期的诚实**——
它没有照那句话去造一个假红；条件②（"唯一形状"改两支）**闭合**，末行那句不声称穷尽达到了这批自立的尺；
条件③（子集/联合换真理由）**实质闭合、措辞差一枚限定**（它确实自己重导了，不是抄验收的推论——我把两把尺都跑了一遍）。

**`A221②` 那格"批准保留"现在有没有数？**（题面点名要我判的一句）

- **"它是防御、不是装饰"这一句：有了，而且比我见过的任何一版都硬。**
  不是靠 M4a＋M1 那一行测试红，而是靠**门输出层**那对读数：同一棵树上 `post ≡ M4a`（那味把语义退回吞掉了）、
  `M4a ≠ M4aM1`（摘掉那味，受追踪的交付字节当场从分母与 finding 里消失）。
  外加一条**不需要量的**结构论证（§1.5）：那一味买的是"不隐身"与规则形状**无关**这件事本身。
- **但"今天有没有一棵真树能看出差别"这一句：没有，这是我量出来的否答（439 形）。**
  ⇒ `A221②` 里那句"现在有数了，不再是我觉得合理"如果被人读成"已经证明该留着"，**读过头了**；
  它真正的依据是三样东西叠起来：**一条码级不变式（最硬）＋ 一发已发生的同类缺陷家族（M4a/`post≡M4a`）＋ 一句预报（第三枚会不会长出来，这一枚仍是判断）**。
  ⇒ 建议把 `A221②` 就地追加这一句分层，别再让"数"和"判断"混在同一格里（**原话不改**）。
- **代价那一栏没变**：＋0.34 s/次、6 枚只读 spawn，仍只付在 step 2 与本地跑门（`A221②` 记的数，**我没有重测**——
  我这轮的读数全部在"改不改输出"这一层，没有去数 spawn）。我能给的担保是两条更基本的：
  `main.go` 的 blob 在本区间**同值**（仪器一行没动），`gitignore.go` 的 diff **每一行都是注释行**（§6 末行那三把尺），
  且我自己跑的整道门 step 2 rc=0 与八数同值（§0.2）⇒ 那笔"每次多花 0.34 秒"的账**没有理由在本批变化**，但它是引用不是复测。

**最后一层，也是我最想让编排者看到的一层**：这批真正值得入账的不是那枚测试，是**它拒绝把上游的一句话当成判据去凑**。
本仓反复记"审计给的修法也是未验证断言"，这次是同一规矩在**反方向**上兑现了一次（下游推翻上游）。
但它兑现得比规矩多走了一步（没停手）——**该补的是规矩，不是这一程的账**。§8.3 那三条可复算判据就是我要交的实物。

---

## 11. 我没有测什么（绝不让沉默被读成批准）

1. **GitHub CI 零次**：没跑过、没读过任何 run id／step 日志。§3.2 的"ubuntu 会同意"是从**这台机自带的 `gitignore(5)` 原文**推的，**不是 ubuntu 上实测的**。
2. **没在 Linux 上跑过任何一枚测试**：全部 win32 ＋ `git 2.52.0.windows.1`。已知两枚平台相关量我没验：
   `core.ignorecase`（本机 true ⇒ §1.4 那枚 `WEIRD/*` 落在 under-read；linux 上它变成 agree，谁都不吃亏）、
   以及 ubuntu 上 `/tmp` 的属主与 `git` 的 `safe.directory` 那一支（本机全局有 `safe.directory=*`，**我没有一台没有它的机器可测**）。
3. **§5.3 那枚竞态是我用白盒手造的状态，不是在真 git 上跑出来的**：我没做"`:240` 与 `:244` 之间塞一枚 commit"的端到端复现。
   所以那一格的力度止于"码上可查＋态可构造"，**不是"这仓现在正漏着什么"**。
4. **过度匹配只搜了 439 枚规则形**（55 手挑 ＋ 384 乘积枚举）。"造不出"只在**我搜过的形状**内成立；穷尽性我不敢说，这正是 §1.5 那条不变式存在的理由。
5. **没跑整树 `go test ./...`**（题面禁令：并行的 panel 程在编译，读数不可归因）。`go vet ./...`/`go build ./...` 我跑了（rc=0），但那只是编译。
6. **`ban #8 internal/=407` 与 `design/` 的 32/30 分家我没有自己重造**：前者我引 §0.2 的现量值并写明归属与锚点，
   后者我没有一枚"净快照"树可扫（`git archive` 出来的树没有 `.git`，扫描会走"问不到 index ⇒ 全扫"那一支，八数根本不可比）。
7. **每枚 commit 之前 `git diff --cached --name-only` 的索引态不可复原**（§8.1）：我只核到零泄漏。
8. **上一程的 M2/M3/M4b/M5/M6 我没重跑**：依据＝§6 末行那三把尺证的"生产行为零改动"＋`main.go`/`gitignore.go` 的 blob 面（行为路径未变）。
   这不等于"那五发今天仍红"，只等于"本批没有改变它们红不红的任何东西"。
9. **新测试对 `holds()` 的目录位（`dirs`）没有覆盖**（§1.7）——我量到了这一层缺口存在、也量到它由上一程的 `:1485/:1488` 钉着，
   但我**没有**独立重跑那一发来证明它今天仍咬得住。
10. **`runtests.sh`/`d22scan.sh` 这两枚脚本自身没被我改过也没被我逐行审**（不在本批射程）；我只用它们取数。
11. **owner 的 `NO CONCLUSION`／`slo-full` 那一套与本批无关**，我一次没碰。

---

## 12. 给人看的那一段（零术语）

这一批没有给用户做出任何新功能，也没有改动任何会跑起来的东西。它做的事用一句话说：
**项目里有一台"看代码有没有犯规"的检查器，它以前有一种可能被悄悄骗过去的走法——某份文件明明在交付清单里，
却被检查器当成"这是构建垃圾，不看"而整份跳过。上一批把这条堵上了，但堵住它的那道闩本身没有任何测试在看着，
谁都可以哪天顺手拆了而没有人报警；这一批就是去给那道闩装上一根会报警的针。**
针装好了：我把它拆下来重跑了一遍，检查器立刻尖叫（我复现了它的每一声）；针留着，检查器和拆掉闩的结果在你现在这棵树上
一个字都不差（这也量了）。

**还能不能被绕过？** 能，两条，都不新：① 如果哪天有人让检查器把某类文件"看漏"（这类失误这个项目已经犯过两次、
也都修了两次），闩仍然会拦住"漏看交付文件"这一种后果，而那根针保证下一次有人改坏时会响；
② 检查器读文件清单时问了两回 git，两回之间如果恰好有人往清单里加文件，理论上存在一个极小的窗口——
这一条我构造出来了但没有在真环境里跑出来，而且它属于**更早一批**留下的形状，这一批只是把话说得太满（我要求它补一句限定）。

**对你机器的安全性改变了什么？** 今天：**零**。这台检查器是开发期的一道门，不在你运行的产品里；
更重要的是，本项目这一带**今天什么都没有接到真的审批路径上**——防"面板替你按下同意"的那扇门压根还没装（见 `A217`/`A220`/`A222` 那三格），
所以这里没有任何东西在今天保护或危害你的机器。它保护的是**将来接上那一天**没人偷偷改坏判定逻辑。

**我要你做什么？** 不用你拍任何新的题。只有两枚一句话的更正需要走一遍实现程（改一个行号、改一处"这句是谁写的"的归属），
外加一处措辞补半个从句；三件都是文字，不碰行为。我建议顺手把 §10 里那三条"下游可以顶回上游"的条件写成规矩——
不写的话，下一程只会学到"顶回可以"，不会学到"在什么条件下才可以"。
