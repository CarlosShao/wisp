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
| 我的仓库外工作区（只建不删，全在 `D:\tmp\d22scan-close-r1-accept\`） | `head\`（`tools/d22scan` 的 pristine 拷贝，`git hash-object` 见 §1.1）／`bin\d22-post.exe`／`bin\d22-M1.exe`／`bin\d22-M4a.exe`／`bin\d22-M4aM1.exe`／`mut\M1` `mut\M4a` `mut\M4aM1`／`fuzz\`（两枚我自己写的探针，`zz_overread_probe_test.go` `zz_combo_probe_test.go`）／`probe\seed1`（§2 的真 git 台件）／`seedtree\` `dirprobe\`（§1 的输出对照树）／`step1-accept.log` `gate-accept.log` `vet.log` `status-before.txt` `status-after.txt` `out-*.txt` `dp-*.txt` `real-*.txt` `combo-fuzz.log` | `ls` 现量 |

**并行读数的边界（先说清楚，免得我的数被人当别人的数用）**：另一位实现程此刻正在 `internal/panel/**` 上作业。
本件**没有**跑过整树 `go test ./...`；`go vet ./...` 与 `go build ./...` 我跑了（rc 均 0，见 §0.3），
但那是编译，不是行为读数。§0.2 的 `ban #8 internal/=407` 是**我在 `310816f` 上量的那枚数**，
它的归属是 `internal/panel/l2_grant_boundary_test.go`（盘上 10:09 的那版），**不是本批造的**，我不追它的漂移。

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

⇒ **合计 431 枚规则形，交付态下我造不出一枚"过度匹配"**。两件必须同时报出来的方法论账：
① 第一版我用 `git check-ignore`（不带 `--no-index`）当尺，被 3 形误判成 under-read——原因是**纯取反规则 `!weird/*` 也会被 check-ignore 以 rc=0 举出**，
而 `ls-files -i -c` 从不举目录；修正后的尺是"窄清单 ∧ 解析过 `!` 的 check-ignore"两问都答否。这一步不改变 over=0 的结论，但说明**这把尺本身需要正控**。
② 我的样本是 431 形，**不是全部形**；"今天没有过度匹配"是"我搜过 431 形没找到"，不是"没有"。

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
  而我量到的只是"交付态的 431 枚规则形上不产生差别"——**这两句话之间的距离就是那味保险的全部射程**。
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
- (iii) 的**为真部分（限定语）**：在**今天交付态的匹配器**下，431 枚规则形里我造不出一棵使门输出改变的树；
  所以"这味今天有没有行为后果"的答案是**没有**，它的价值全部在"下一枚匹配器读宽"那一方。
  **它不是"为一个不可能发生的类挡事"**——那一类已经发生过两次（F3 ①②），且一发单词元的退回就能再造（§1.3）。
- (ii) 的**为真部分**：它的前半句（可观察输出相同）我复现了。
  **它的结论（⇒ 窄清单那味无行为后果 ⇒ 该删）不成立**，理由＝§1.3 第 2 条 `post == M4a` 与 §1.5 那条不变式。
- **对 `A221②` 的账**：那格写"我前面批准保留那句现在有数了"。**按"防御、不是装饰"这个结论，数确实有了**
  （M4a／M4a＋M1 那对读数我逐枚复现）。但要说清它**是哪一种数**：它量的是"**存在**一个输入（带过度匹配的树）使这味决定门输出"，
  而**不是**"**当前**任何一棵真树使这味改变门的输出"（后者我今天专门去量了，答案是否）。
  ⇒ **留在判断里的那一块不是零，但也不是这张表**：它是"这一类会不会再长第三枚"的预报。
  这句要如实写回 `A221②`。
