# 212-v3 非实现者验收件（攻 `212-r3` 的三笔改动：4ea6411e 骨架 → d03d166f 落地 → 0cb71a58 满稿）

- 本程：验收腿 **212-v3**（裁决者≠实现者；实现＝`212-r3`，ⓐ修法裁定＝编排者）。
- 起手锚：`dev` HEAD＝`57a33804`（本程第一条命令现跑，时刻 `2026-10-04 10:45:03 +08`）。
- 本程只读产码、只在 `tools/d22scan/**` 里做突变；AC 框一枚不碰；不放宽任何既有断言。
- 每一节的判语只认本程自己跑的尺；`212-r3` 自述里的每一句都按待验断言处理。

## ① 进场、基线与台件纪律自证

**起手（与本节同发取数，档 `logs/md5-start.txt`）**：`date`＝`2026-10-04 10:45:03+0800`；`git rev-parse --short HEAD`＝`57a33804`；分支 `dev`。

**被验三笔的写面名册（尺＝`git show --stat --oneline <sha>`，逐枚）**：

| commit | 写了什么 | 名册 |
|---|---|---|
| `4ea6411e` 骨架 | 只有本腿自己的证据件与四枚读数档 | `.scratch/wisp/probes/212/r3/{fix-and-readings.md,logs/*}`（**零产码**） |
| `d03d166f` 落地 | 四枚文件、**零新增文件** | `internal/agent/approval/pending_read.go`(M)·`tools/d22scan/main.go`·`tools/d22scan/scan_test.go`·`tools/d22scan/selftestsamples.go` |
| `0cb71a58` 满稿 | 证据件＋读数档 21 枚 | 全在 `.scratch/wisp/probes/212/r3/**` |

⇒ 三条推论当场可钉：① 本票三笔**没有碰票面一个字节**（`git log --oneline -4 -- .scratch/wisp/issues/212-*.md` 最新一枚是编排者的 `6a4db468`，212-r3 三枚 commit 的 name-only 里 `issues/212` **零命中**；票面现量 `- [x] **AC`＝**5 枚**、`- [ ] **AC`＝**0 枚**，文件名仍无 `-done`）；② `d03d166f` 在 `internal`／`cmd` 里**只 M 不 A**（§⑤ 归因用得上）；③ 写面与派单声明的五枚路径一致，⛔ 未见越界（禁区六枚路径＋`allowlist.txt` 零命中：`git show --name-only` 三枚逐枚核过）。

**四枚受检文件 md5 起手值**（`logs/md5-start.txt`）：`main.go 442d46cadef8c641c54f7ac00984a470`·`selftestsamples.go f9e10163199a89c2614daf1ab1ad7d52`·`scan_test.go 1bcf66c1bc7da3e7ff0601b66da13504`·`pending_read.go acef5b1739e9446cd7118bbffa6a0a3b`。
其中前三枚**工作树＝HEAD blob**（同 md5），第四枚不等——尺坑与处理写在 §⑧。

**"只在 `tools/d22scan` 一亩三分地里做突变"这一条怎么落地的（方法具名，因为这直接决定 §② 的种法）**：
ban #9 的扫描面是 `walkGo(internal)`＋`walkGo(cmd)`（`main.go:264`·`:267`，两枚调用，别处没有），
所以**任何"真树正控"都必然落在 `internal/` 或 `cmd/` 里**——那正是并发腿 260-r3 的写面，也是门禁分母本身。
两难之下本程选的形是：**不动仓内任何被跟踪文件，改在仓外搭一枚"真树形状"的假根，用同一枚产码 binary 以 `-root` 扫它**。
- 为什么这仍然算 AC#4 的正控而不是夹具复刻：跑的是 `go run . -root D:/tmp/212v3tree`＝**产码路径**（`checkRoot`→`walkGo`→`scanGoFile`→ban #9→`os.Stat`），
  不是 `selftest.go` 的夹具装配；假根里 `internal/probe/*.go` 是**真文件、真注释、真不存在的路径**，
  判"存在"的那次 `os.Stat` 打在盘上（`-root` 拼 `filepath.FromSlash(tok)`）。
- 搭台过程本身撞到两枚硬退出并留档（`logs/plant-00-baseline.txt` 前两版的尾行）：
  假根缺 `design/` 时报 `ban #8 scope design/ examined 0 files - ... never leave a scope pretending to scan (ticket 71 AC#4)`（rc 2），
  补了 design/frontend 又报 `scope ban #7 internal/tools/ examined 0 files but is declared live in declaredScopes()`（rc 2），
  补齐 `internal/tools/x.go` 后基线才 `rc=0` 全绿。⇒ 假根最终形态：12 枚产码 Go 文件＋design/＋frontend/＋`tools/d22scan/allowlist.txt`（从仓里**复制**、未改）。
- ⛔ 全程未在 `internal/`·`cmd/` 写过一个字节；`git status --porcelain -- tools internal cmd` 在每次突变后都取过（终态见 §⑧）。

**台件纪律自证**：突变一律先落文件（Edit 工具改 `main.go`）再判，输出全量重定向到 `logs/`，`finally` 语义由"还原→md5→复跑"三步代替；
取 rc 一律 `cmd > file 2>&1; echo rc=$?`（⛔ 没有 `cmd | head; echo $?`）；仓内台件只建不删（`logs/` 49 枚已随 `769eec1e` 入仓），
删掉的只有**盘外假根的载体文件**（那正是 AC#4 "删掉⇒必不响"要求的动作，且不在仓内）。

## ② 格1：AC#4 那把正控到底还有没有牙

**判语：成立——正控仍有牙；腿自述那句"`f.Comments` only、字符串字面量永远扫不到"为真，但它做瞎的不是 AC#4 这一格。**

### ②.1 我自己种的正控（两向，逐字）

种（`internal/probe/zz_plant.go` 的第 3 行注释＝`// See docs/evidence/s1/212-v3-plant-phantom.md for the readings.`，
该件在假根与真仓里都不存在）⇒ **响**，档 `logs/plant-p1-fullspelled-ring.txt:11`：

```
internal/probe/zz_plant.go:3: [phantom-citation] comment cites repo path "docs/evidence/s1/212-v3-plant-phantom.md" which does not exist on disk (ticket 212 ban #9): fix the citation or create the file; shorthand forms are a separate prescription, not a violation
```
同档 `:14`＝`rc=1`。

删掉载体（其余 12 枚文件一字未动）⇒ **不响**，档 `logs/plant-p10-plant-removed.txt:11-12`：

```
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=12, ... ban #8 internal/=12, ban #8 cmd/=1; ...
rc=0
```

第三向（把被引的那件**真造出来**再扫）⇒ 同一注释**不响**（档 `logs/plant-p7-page-exists-silent.txt`，`rc=0`、phantom 行 0 枚）
⇒ 证明响的是**存在性**而不是"注释里出现了路径形状"，这把 P2/P3/P4/P5/P9/P12 那六枚"静默"从"尺坏了"这一解释里排除掉。

### ②.2 喂给 ban #9 的节点集合（自己读码，不抄腿）

`main.go:831-832` 逐字 `for _, cg := range f.Comments { for _, c := range cg.List {` ＝**只走注释组**，
`f` 由 `parser.ParseFile(..., parser.ParseComments)` 得到（`:706`）；字符串字面量不是 `f.Comments` 的元素。
文件集合由 `walkGo` 决定，`:687` 逐字 `if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") { return nil }`，
调用面只有 `internal/`·`cmd/`（`:264`·`:267`）⇒ 名册行 `main.go:55-57` 那句"tools/** 与 _test.go 在射程外"**为真**。

四枚载体逐形现跑（同一条 binary，同一假根）：

| 载体 | 形状 | 结果 | 档 |
|---|---|---|---|
| S1 | 幻影路径**只在字符串字面量里**（`const beacon = "read docs/evidence/s1/212-v3-string-only.md first"`） | **rc=0，phantom 0 枚**＝看不见 | `logs/scope-s1-string-literal.txt` |
| S2 | 幻影路径在 `/* 块注释 */` 里 | **rc=1，响**（`internal/probe/zs2.go:3`） | `logs/scope-s2-block-comment.txt:11` |
| S3 | 幻影路径在 `internal/probe/zz*_test.go` 注释里 | **rc=0 不响**，且 `bans #1-5 internal/=12` 与 `ban #8 internal/=13` **同时出现**＝这枚文件只进 ban #8 的分母 | `logs/scope-s3-redo-testonly.txt:11` |
| S4 | 幻影路径在 `tools/probe/s4.go` 的产码注释里 | **rc=0 不响** | `logs/scope-s4-redo-toolsonly.txt:11` |

（⚠ 枚数口径具名：那两枚 12/13 是**假根**的数，不是真仓的 228/501；S3 那一发我在脚本里先落 `zs3_test.go` 才取数，
`zs2.go` 已在前一步删掉，所以 13＝12 枚产码＋1 枚 `_test.go`。第一次批量跑我把 S3/S4 与 S2 叠在同一发里，红句会被上一枚载体污染——
已重跑单发隔离，两发都留档：`scope-s3-underscore-test-go.txt` 是叠发的坏读数，`scope-s3-redo-testonly.txt` 是隔离后的**采信值**，`scope-s4-*` 同理。具名记账，不藏。）

### ②.3 那句话会不会同时把 AC#4 的正控做瞎——本程的裁法

**不会**，两格要分开：
- 票面 AC#4 写的是"造一件**注释**引用一枚盘上不存在的证据件 ⇒ 那一形必响"。我种的正是注释（`//` 与 `/* */` 两形都响，S2），
  所以 AC#4 的正控**在产码路径上是真的**，`f.Comments` 的限制**碰不到它**。
- 被 `f.Comments` 做瞎的是**另一格**：编排者那道待机的"要不要把 `tools/**` 纳入 ban #9 射程"。
  本程复认了腿 §⑥ 那句推论的形状——旧那枚 `docs/evidence/s1/212-citation-ruler.md` 在 `selftestsamples.go` 里**只出现在字符串字面量里**
  （尺＝`git show d03d166f^:tools/d22scan/selftestsamples.go | grep -n 212-citation-ruler` → 只有 `:375` 与 `:378` 两行，两行都是 Go 字符串），
  加上 S4 证明 `tools/**` 今天整体不在射程、S1 证明即使进射程也看不见字符串 ⇒ **"扩射程"这一支单独并不能把那一格变成有牙的判据，只有改措辞能**。
  ⇒ 这条是给编排者那道批准的**成色信息**，不是 `212-r3` 的伤；腿把它写在 §4 明处的做法本程判为**该记的账它记了**。


## ③ 格2：ⓐ 豁免的真实射程（本票最要紧的一格）

**判语：带条件成立。** 条件＝豁免面确实比"省略号在中段"宽，本程**在盘外真树把"谎报伪装成缩写"的真实写法造出来了**（五形），
但**真仓当前活体仍是 0 枚**（我自己写的同语义普查），所以编排者"不收窄"那道裁定在本程读数下不被推翻——
代价是：**这五形的免罪从此只能靠样本钉着，不能靠"缩写都是诚实的"这句假设**（腿把三形各钉一枚 silent 样本正是这条条件的落地，见 §⑦）。

### ③.1 双向都跑了（同一条产码 binary，档 `logs/plant-*.txt`，汇总 `logs/plant-summary.txt`）

| 发 | 种的形状（注释里的 token 逐字） | 期望 | 实测 |
|---|---|---|---|
| **P1** | `docs/evidence/s1/212-v3-plant-phantom.md`＝**全拼、不含任何标记、盘上没有** | **必红** | **rc=1 响**，token 逐字见 §②.1 |
| **P2** | `docs/evidence/s1/212-v3-...-phantom.md`＝**带 `...` 的全拼不存在**（票面说的"这是设计"那一形） | 静默 | **rc=0 不响**，phantom 行 0 枚 |
| P3 | `docs/evidence/s1/212-v3-tail-ellipsis.md`＋尾 `…`(U+2026) | 静默（裁ⓐ 的臂 2） | rc=0 |
| **P4** | `**docs/evidence/s1/212-v3-bold-phantom.md**`＝**markdown 加粗包一根完整幻影** | 静默（＝**谎报伪装成"强调"**） | **rc=0 不响** |
| P5 | `docs/evidence/s1/212-v3-glued.md`＋`见`(U+89C1)＋`docs/evidence/s1/...`＝CJK 粘连（臂 3） | 静默 | rc=0 |
| **P6** | `见docs/evidence/s1/212-v3-glued-nomark.md`＝**粘连但无标记** | 必红（边界） | **rc=1 响**，token `docs/evidence/s1/212-v3-glued-nomark.md` |
| P8 | 第一行全拼幻影、**第二行**才放 `...` | 必红（边界＝豁免按 `c.Text` 逐条注释算） | **rc=1 响**，token `docs/evidence/s1/212-v3-crossline.md` |
| **P9** | `docs/evidence/s1/212-v3-tail-dots.md...`＝**完整名＋ASCII 三点（"等等"式写法）** | 静默（＝**谎报伪装成"等等"**） | **rc=0 不响** |
| **P11** | `docs/evidence/s1/212-v3-space-then-star.md *`＝**标记前有空格** | 必红（边界＝region 被空格截断） | **rc=1 响** |
| P12 | `docs/evidence/s1/*.md`＝通指一枚目录 | 静默（③） | rc=0 |

⇒ 豁免的**真实射程**（代码给的，不是注释给的）＝"token 起点落在一个含 `...`/`…`/`*` 的 `shorthandRegionRe` region 起点"，
所以三臂免罪：**中段标记／尾随标记（任意一枚，包括 markdown 强调的两枚 `*`）／CJK 粘连尾段**；
边界三枚（P6 无标记、P8 换行、P11 空格）都照红——**腿写进 `main.go:925-927` 的 P6/P8/P9 三条边界本程逐枚复现，无一句过窄**。

### ③.2 "有没有一种真实写法把谎报伪装成缩写"——有，具名两枚新形（腿没点名的那一族）

腿与 v2 报的是"完整名＋尾 `…`／尾 `*`／CJK 粘连"。本程额外造出**两形更贴近日常书写**的：
- **P4 `**path**`（markdown 加粗/斜体）**：注释里 `// **docs/evidence/s1/xxx.md** is the reading.`
  是 Go 仓注释里真会出现的写法（有人拿 `**` 在注释里做强调），两枚 `*` 都落在 token 之后 ⇒ 整串进 region ⇒ **谎报静默**。
- **P9 `path...`（"等等/详见"式 ASCII 三点）**：`// See docs/evidence/s1/xxx.md... for details.`
  尾随的 `...` 是"等等"而不是"缩写"，但豁免只看 region 起点 ⇒ **谎报静默**。
我还试过但**没能溜过去**（当场响）的形：粘连无标记（P6）、标记换行（P8）、标记前加空格（P11）、
全角标点收尾 `docs/.../xxx.md。`（`。`U+3002 在两类字符集之外⇒ token 正常结束⇒响，这条与 P6 同机理，未单列档）。

⇒ **本程的定性**：这不是"实现者的缺陷"，而是**裁 ⓐ 的定义性代价**，而且它比票面那句话更硬——
**"缩写"与"谎报＋一个尾随符号"在词面上不可区分**。票面终裁节写的"今日零实伤"如果只指"中段省略号"，
本程把它扩到"尾随标记／粘连／加粗／等等式省略号"五形仍为零实伤（下条），所以裁定不改；
⛔ 但今后任何一形**出现活体**时，"不收窄"那句就不再自动成立——这条要由下一次现量说话，不由本程越权改 regex。

### ③.3 真仓现量（我自己写的普查，尺＝把 `main.go` 的两枚正则与 `symRefRe` 逐字搬到一枚盘外 Go 程序里跑真树）

档 `logs/census-real-repo.txt`（程序＝`D:/tmp/212v3census/main.go`，只读；文件名册＝`git ls-files internal cmd` 去 `_test.go`，与产码走的面同口径）：

```
CENSUS tracked production go files under internal/ cmd/
RING(exempt-free phantom, i.e. what the gate reports today)=0
EXEMPT(class 3 swallowed)=1
SYMREF(excluded api cite)=3
```

唯一那一枚被豁免吃掉的 token，逐字（`file:line` ＋整行注释）：

```
cmd/wisp/panel_host_windows.go:29  docs/evidence/s1  FULLLINE=//     EARLY by the nested one (about 0.55s per open) - docs/evidence/s1/... the 33-p1
```

判语＝**它是真缩写，不是被伪装的谎报**：那行写的是 `docs/evidence/s1/...` 加一句人话指向 33-p1 的探针，
不带任何"完整文件名"，把 `...` 去掉剩下的是 `docs/evidence/s1`（一枚**存在**的目录）。
⇒ **真仓"完整幻影被尾随标记/粘连救下"的活体＝0 枚**，与 `212-v2`（git archive 差集）和腿 §3"活体 0 枚"三把独立尺同结论。
另与门禁互证：`sh scripts/d22scan.sh` 终态 rc=0 零红（§⑧）＝ RING 那一路今天盘上确实没有东西。

## ④ 格3：新增那枚前缀同步钉是不是同义反复

**判语：成立（有牙）。** 两条独立理由：
1. **不是别名**：`repoPathRe`（`main.go:890-891`）与 `shorthandRegionRe`（`:904-905`）是**两份独立的手抄字面量**，
   钉 `scan_test.go:2399-2443` 用 `regexPrefixAlternation(repoPathRe.String())` 与 `...RegionRe.String()` **各读各的**，
   再比集合（`:2429-2442`）；`regexPrefixAlternation`（`:2449-2465`）在形状不是 `(?:a|b)/` 时**返回 error**（`t.Fatalf`），
   所以"把两枚表抽成一处再拼"这种重构若发生，钉会以"读不到形状"红而不是默默绿——**它没有把被钉对象定义成自己的输入**。
2. **我自己的刀**（与腿反向：腿动 `repoPathRe`，我动**另一枚**，正好打它红句的第二臂）：
   给 `shorthandRegionRe` 的前缀组加一枚 `plugins|`（Edit 逐字：`(?:docs|plugins|\.scratch|internal|cmd|tools|scripts)/[A-Za-z0-9_./\-\x{4e00}-\x{9fff}*…]*`），
   ⇒ `go test -count=1 -run 'TestBan9PrefixTablesAreTheSameList' ./...`（模块内）**指名 FAIL**，档 `logs/m2-test-prefix-drift.txt:1-2` 逐字：

```
    scan_test.go:2441: ban #9's two hand-copied prefix tables disagree: repoPathRe=[docs \.scratch internal cmd tools scripts], shorthandRegionRe=[docs plugins \.scratch internal cmd tools scripts]; only in repoPathRe=[] means a shorthand path there is NOT exempt and class 3 convicts again (ticket 212 AC#3), only in shorthandRegionRe=[plugins] means a COMPLETE phantom path there gets swallowed by the region, which widens the exemption without an owner's approval. Edit both lists in one commit, and read the ban #9 comment in main.go before deciding which shape is intended.
--- FAIL: TestBan9PrefixTablesAreTheSameList (0.00s)
```

   同一棵突变树上 `go run . -self-test` 仍 **`all 40 direction checks passed`／rc=0**（档 `logs/m2-selftest-same-tree.txt`）
   ⇒ 与腿 m2 的读数**同形**（腿报的是第一臂 `only in repoPathRe=[plugins]`，本程报的是第二臂），两臂都有牙。
   还原：`git cat-file blob HEAD:tools/d22scan/main.go > main.go` → md5 回到起手 `442d46cadef8c641c54f7ac00984a470`（`logs/m2-md5-restored.txt`），
   `git diff -- tools/d22scan` 空（§⑧）。

**刀法具名（腿踩过的那脚我换了刀）**：⛔ 不用"照行号 890/891 下 sed"——那会命中 `var` 行、突变静默不落地（腿 §5 自陈踩过一次 rc=0 才发现尺指错）。
本程两发突变都走**内容锚的 Edit 工具**（old_string 含整条正则字面量），改完**先 `md5sum` 证落盘变了**（m2 突变态 `e055e0f46caf6f694442534cf8e0f6e8` ≠ 起手）再跑判据。

**这枚钉的一处真实局限（不算伤，具名）**：它比的是**前缀组的字符串集合**，所以
① 列序变化不红（`:2412-2414` 自己说了理由：每个候选后面紧跟字面 `/`，序不改匹配——本程认这个推理）；
② 一枚表内部**拼法不同但集合相同**（如 `.scratch` vs `\.scratch`）会红（集合按原文比，✓ 有牙）；
③ 两枚表**同时**加同一枚新前缀时钉不会红——那种情况下豁免面与匹配面**一起**扩大，本票的"③ 在新前缀上复发"这一形确实不成立了，
   所以这不是洞。⇒ 结论：**不是同义反复，且它对"两张表被分别改动"这一真实漂移有牙**。


## ⑤ 格4：AC#5「八枚 ban 读数逐名不变」与两枚 +1 的归因

**判语：成立。** 八枚逐名不变（以本程自跑的基线终值对照腿与编排者的读数），两枚 +1 **都不是 `212-r3` 的伤**，
但**腿给 `ban #8 internal/` 那枚点名的文件是错的**，本程具名更正。

本程 10:53 自跑（`logs/gate-v3-final.txt`，rc=**0**，摘句 `logs/gate-v3-scoperows.txt`）八枚逐字：

```
d22scan: scope bans #1-5 internal/      examined 228 production Go files
d22scan: scope bans #1-5 cmd/           examined  38 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 501 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  98 Go files, comments and _test.go included
```

⇒ 与编排者 10:40 那发、与腿 §7 终值**逐字相同**；与腿起手（`internal/=500`·`cmd/=97`）差的两枚正是 §⑨ 里归因的两枚。

**归因尺（两把，都本程现跑）**：
1. `git log --name-status --diff-filter=A --format="COMMIT %h %ad %s" --date=format:"%H:%M" b1e59d63..HEAD -- internal cmd`
   （`b1e59d63`＝腿自述的起手锚，现读 `260-r1 补一句实话…`，它到 `0cb71a58` 之间隔 **20 枚** commit）⇒ 整个区间在 `internal/`＋`cmd/` **新增的 .go 只有两枚**（档 `logs/added-go-between-anchors.txt`）：

```
COMMIT e3e19e8e 10:36 test(ball)〔260-r2〕把"借到的必是裸 Esc"改成"借到的必是配置里那枚键"…
A	internal/ball/hotkey_cancel_borrow_expect_260r2_test.go
COMMIT 2a10134e 10:20 256-r1 落地：常驻那条腿的审批门改吃 [risk] 两枚配置（形ⓐ 签名不加参）
A	cmd/wisp/resident_approval_risk_256_windows_test.go
```

2. `git ls-tree -r --name-only <锚> -- internal|cmd | grep -c '\.go$'` 作差（同档尾部）：`internal` 500→501、`cmd` 98→99，
   两锚的**列表作差各只有一行**，就是上面那两枚新文件。⇒ 分母每边各 +1，**逐名对得上，没有第三枚**。

**结论与更正**：
- `ban #8 internal/` 500→501 ＝ `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go`，随 **`e3e19e8e`（260-r2，10:36）**。
  ⚠ **腿 §7 那一格点名的文件是 `internal/ball/hotkey_live_test.go`（mtime 10:25）——那枚在同一区间里是 `M` 不是 `A`（`git show --name-status e3e19e8e` 里两枚并列可见），改内容不改枚数**；
  腿用 mtime 猜，方向对（并发腿 260 那一族）、**具体那一枚指错**。本程判它"归因结论成立、点名需更正"，⛔ 不是伤。
- `ban #8 cmd/` 97→98 ＝ `cmd/wisp/resident_approval_risk_256_windows_test.go`，随 **`2a10134e`（256-r1，10:20）**——与腿一致，本程复认为真。
- 口径互证（⛔ 别把这两枚 +1 读成"八枚都漂了"）：两枚新文件都以 `_test.go` 结尾，而 `walkGo`（`main.go:687`）把 `_test.go` 挡在 bans #1-5 之外、ban #8 的 `walkEmoji` 把它们算进去
  ——这正是 `bans #1-5 internal/=228` 一动没动、只有 `ban #8 internal/` 动的机理；我在 §②.2 的 S3 那一发里独立看到了同一机理（假根 `bans #1-5 internal/=12` 与 `ban #8 internal/=13` 同发并存）。
- `212-r3` 自己**不可能**移动分母：它在 `internal`／`cmd` 里只有 `M pending_read.go`、**零新增文件**（§① 名册尺），且它改的全是注释行（§⑥ 的 0 行非注释尺）。
- ⚠ **口径差具名**（免得下一位把我的尺和门禁的尺当同一把）：`git ls-tree` 在 `cmd/` 给 **99**、门禁给 **98**，差的那一枚是 `cmd/wisp/testdata/esclistener/main.go`——`walkGo`/`walkEmoji` 都对目录名 `testdata` 走 `filepath.SkipDir`（`main.go:682`·`:1117`，`walkText` 同族在 `:1031`）。`internal/` 无 testdata 下的 .go，所以 501＝501。**两把尺都成立，只是口径不同**，我按门禁口径报数、按 ls-tree 口径作差，两句都在档。

## ⑥ 格5：(a) 那处"诚实化"引用的每一枚 `file:line` 真不真

**判语：成立**（四枚引用逐枚打开对全中，**没有把本票的缺陷再犯一遍**）。

新注释＝`internal/agent/approval/pending_read.go:40-54`（现读行号与派单给的范围一致）。它引用的每一处：

| 引用（注释里逐字） | 打开对 | 判 |
|---|---|---|
| `internal/tools/gate.go:12-13` (ticket 17) names exactly TWO fields ... `RulesHit and Reason` | `:12` `// what the caller must show the user. The confirmation card (tickets 21/37)`；`:13` `// renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note), so` | ✓ 行号对、"只两枚"对、票号 17 对 |
| `cmd/wisp/panel_pump.go:44-46`，且明写那一注是 `a copy, not a judgement` | `:44` `// The mapping is a copy, not a judgement: Level/RulesHit/Reason/`；`:45` `// SessionOverrideBlocked are the fields internal/risk put on the queue item`；`:46` `// when it was admitted, and panel renders them through the same` | ✓ 行号对、四枚逐字在该 3 行内、引号短语是原文 |
| `panel.CardViewFromDecision` `(internal/panel/approval.go:72-88)` is where all four are read onto the view | `:72` `func CardViewFromDecision(subject ApprovalSubject, decision risk.Decision) ApprovalCardView {`，`:77` 起结构字面量，四枚在 `:81 Level`／`:82 RulesHit`／`:83 Reason`／`:85 SessionOverrideBlocked`，`:88` 收到 `}` | ✓ 范围 72-88 正好框住函数与四枚赋值 |
| `ticket146_liveapprovals_backing_test.go`（同段**上一节**注释 `:36-38`，非本次改动） | `find internal cmd -name ticket146_liveapprovals_backing_test.go` → `internal/agent/approval/ticket146_liveapprovals_backing_test.go` | ✓ 存在（顺手核，防"改一处留一处"） |

改前的错处也逐字对上了账：旧句 `internal/tools/gate.go:13` 配的是"四枚字段"（`RulesHit, Reason, Level and SessionOverrideBlocked ... per ticket 17's frozen-contract note`），
而 `:13` 只点名两枚 ⇒ 编排者 `A592` §2 那条"指错出处"复认为真，**本程确认它已被改成实话**，
且新句把"两枚／四枚"分挂在两枚不同出处、并额外点名第三处实现（`approval.go`）作支撑——**归因层次比改前更细，没有新造断言**。
另外新句 `its phantom-citation ban checks only that a cited path exists` 本程实测为真（§②：它只看存在性，判不出"指错出处"）。

**行为零变化的独立尺**（不采信腿的 `--numstat` 叙述）：
`git show d03d166f -- tools/d22scan/main.go` 与 `-- internal/agent/approval/pending_read.go` 两枚文件，
把改动行滤成"非整行注释"后**各 0 行**（尺＝`grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | grep -vE '^[+-][[:space:]]*//' | wc -l`，档 `logs/diff-main-noncomment.txt`／`logs/diff-pendingread-noncomment.txt`，两档都是空文件）。
⇒ (a) 与"名册行／自述"两处的改动**确实是纯措辞**。

**顺手量到的一枚措辞账（⛔ 不构成退回，具名给编排者/续程）**：三枚新样本在 `selftestsamples.go:425/437/450` 自标 `SHAPE 1/2/3`＝
（尾部 `…`／尾部 `*`／CJK 粘连），而 `main.go:812-820` 把臂编号写成 (1) 中段标记／(2) 尾随标记／(3) CJK 粘连，
`main.go:825-826`·`:927-929` 与名册行 `main.go:67-68`（"All three are pinned expect-silent in selftestsamples.go (**212-r3**）"）都按**后者**的编号在指——
于是"臂 (1) 中段标记"的样本其实是 r2 那枚 `cites-shorthand.go`（`:396-406`，wantSilent，token `docs/evidence/s1/152-...-accept-r2.md` 等三枚），
而"臂 (2)"有两枚样本。**覆盖是齐的（4 枚 phantom-citation 的 expect-silent 用例：r2 一枚＋r3 三枚，我已用 m1 逐枚验红），错的是编号与"one case per arm"那句**。
判语＝**说得不准，但不是"比实际窄"**，与本票要防的"读起来比实际窄"不同形 ⇒ 记一枚措辞小账，改法＝把 `main.go` 的臂编号与样本标签对齐（纯注释，零行为）。

## ⑦ 格6：恒真判据自查——三枚新样本是不是装饰

**判语：成立**（三枚都有牙，没有一枚是装饰）。

**我自己的 m1**（刀＝`main.go:834` 逐字 `short := shorthandPathStarts(c.Text)` → `short := map[int]bool{} // 212-v3 MUTATION m1: exemption switched off`，
突变态 md5 `ac5f59968795a6347b45dfc7c0f198be`＝先证落盘再判，档 `logs/m1-md5-mutated.txt`）：

终行逐字（档 `logs/m1-selftest-exemption-off.txt:69`）：

```
d22scan -self-test: 4 direction(s) failed, 36/40 passed - the gate does not see what it claims
```

四枚 FAIL 逐枚（同档 `:61/:63/:65/:67`，我只摘 token 段）：

| FAIL 的用例 | 响在哪一形（token 逐字） | 是不是它声称钉的那一形 |
|---|---|---|
| `cites-shorthand.go`（r2 那枚，3 个 token：`docs/evidence/s1/152-...-accept-r2.md`／`docs/evidence/s1/212-`／`docs/evidence/s1/152`） | 中段 `...`／`*`／`…` | ✓ 臂(1) |
| `cites-tail-ellipsis.go:4` `docs/evidence/s1/probe-tail-ellipsis.md` | 完整名＋尾 `…` | ✓ 臂(2) 之一 |
| `cites-tail-star.go:4` `docs/evidence/s1/probe-tail-star.md` | 完整名＋尾 `*` | ✓ 臂(2) 之二 |
| `cites-cjk-attached.go:5` `docs/evidence/s1/probe-swallowed.md见docs/evidence/s1` | CJK 粘连把前面的完整名拖进 region | ✓ 臂(3)，且 `见` 在 token 里＝粘连形逐字复现 |

⇒ **摘掉豁免后没有任何一枚仍静默**（三枚新样本各响一次，token 互不重叠，且都响在自己声称的那一形上，不是"响在别的 token 上"的假钉）。
还原（`git cat-file blob HEAD:tools/d22scan/main.go > main.go`）：md5 回到起手 `442d46cadef8c641c54f7ac00984a470`（`logs/m1-md5-restored.txt`），
复跑 `all 40 direction checks passed (20 expect-ring, 20 expect-silent)`／rc=0（`logs/m1-selftest-restored.txt`）。

**顺带验的 (b) 那半格——ring 向没被削弱**（票面 AC#4 的凭据在 (c) 改址之后是否还成立）：
`-self-test` 的 ban #9 ring 向终行逐字（`logs/m1-selftest-restored.txt` 内）＝
`ban #9 phantom-citation    ring   OK     a comment citing docs/evidence/s1/212-comments-phantom-citation-v2.md, which the fixture does not seed | ranged on 1 finding(s): internal/probe/cites.go:3: ...`
⇒ 新指向的那枚件**真仓存在**（`ls -l`＝33,578 字节、mtime `Oct 4 09:57`，本程现跑，复认编排者③）而**夹具不 seed**，
所以 ring 向照旧响——**(c) 没有把正控改成"永远绿"**，也没有变成新幻影。


## ⑧ 门禁读数（终态）

本程自己跑的四把尺，全部重定向到 `logs/`，rc 用 `$?` 直取（⛔ 没有用 `cmd | head; echo $?` 那种取到 `head` 退码的形状）。

| 尺 | 命令（逐字） | rc | 终态逐字（关键行） | 档 |
|---|---|---|---|---|
| 门禁本体 | `sh scripts/d22scan.sh` | **0** | `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` ＋ `d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp` ＋ 八枚 scope 行 `bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=501, ban #8 cmd/=98` ＋ `d22scan: clean - no D22 ban violations` | `logs/gate-v3-final.txt`／八行摘出 `logs/gate-v3-scoperows.txt` |
| 仪器自检 | `cd tools/d22scan && go run . -self-test` | **0** | 终行 `d22scan -self-test: clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)`；名册行 `... = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation] + 1 finding type(s); 40 cases, 10 tag(s) covered, both directions required per tag` | `logs/m1-selftest-restored.txt` |
| 模块内 vet | `cd tools/d22scan && go vet ./...` | **0** | 零输出 | `logs/vet-tools.txt` |
| 新钉在门禁里真跑过 | `grep -n "PASS: TestBan9PrefixTablesAreTheSameList" logs/gate-v3-final.txt` | 命中 | `--- PASS: TestBan9PrefixTablesAreTheSameList (0.00s)`（第 224 行）＝这枚钉**在 `sh scripts/d22scan.sh` 的第一步里执行**，不在 `-self-test` 里 | `logs/gate-v3-final.txt:224` |

**八枚 scope 行与编排者 10:40 那发、与 `212-r3` §7 终值逐字相同**（我 this 程 10:53 自跑，非引用）。

**突变还原尺（起手 vs 终态）**：
- `md5sum -c logs/md5-start.txt` → 四枚全 `OK`：`tools/d22scan/main.go 442d46cadef8c641c54f7ac00984a470`·`selftestsamples.go f9e10163199a89c2614daf1ab1ad7d52`·`scan_test.go 1bcf66c1bc7da3e7ff0601b66da13504`·`internal/agent/approval/pending_read.go acef5b1739e9446cd7118bbffa6a0a3b`（档 `logs/md5-final-check.txt`）。
- `git diff -- tools/d22scan internal/agent/approval/pending_read.go | md5sum` = `d41d8cd98f00b204e9800998ecf8427e` ＝ **空 diff 的 md5**（档 `logs/diff-md5-end.txt`）⇒ 两发突变零残留。
- ⚠ **尺坑具名**（本程起手就撞到并改刀）：`pending_read.go` 的**工作树 md5 ≠ `git show HEAD:<path> | md5sum`**（工作树 `acef5b1739e9446cd7118bbffa6a0a3b` vs HEAD blob `219a0465f27c469217227ef7a48b907d`，档 `logs/md5-head-vs-work.txt`），`file` 报该文件 `with CRLF line terminators` 而 `.gitattributes:4` 写 `*.go text eol=lf` ⇒ 差异**只来自行尾**，`git status --porcelain` 对该文件为空。所以本程对**这一枚**文件的还原尺用 md5（对起手工作树）**＋ `git diff` 为空**两把，不用 md5(HEAD) 比 md5(工作树)。`tools/d22scan/*.go` 三枚是 LF，两把尺同值。
- ⚠ 编队对账：终态 `git status --porcelain -- tools internal cmd scripts docs .github` 有 **2 枚**（` M cmd/wisp/resident_approval_windows.go`·` M internal/agent/approval/approval.go`，档 `logs/porcelain-end.txt`），逐枚核对＝并发腿 **260-r3** 声明的写面，**本程一枚未碰**；我 10:53 那发门禁正是扫着这两枚在飞的修改跑的，它们是对**已存在文件**的修改（不增减文件枚数），故八枚分母不受影响——这一点与 §⑤ 的归因分开讲，不混。
- 卫生（非本票六格，具名不修）：`gofmt -l tools/d22scan` 报 `selftest.go`；本程把 **HEAD 版**该文件落到临时处再 `gofmt -l` **同样被报**（`logs/gofmt-head-selftest.txt`）＝起手既有，`212-r3` 未碰该文件，判定与腿 §7 一致。
- ⚠ **我这发读数的时刻性**（防下一位拿我的数当"永远如此"）：上表八枚是 **10:53** 那发。交件前复查时（`logs/porcelain-final.txt`）并发腿 260-r3 已在树里放了**两枚未跟踪的 `_test.go`**（`cmd/wisp/resident_cancel_key_wording_260r3_test.go`·`internal/agent/approval/cancel_key_wording_260r3_test.go`，另有 `M cmd/wisp/resident_task_source_live_246_windows_test.go`）——按 §⑤ 那把机理（ban #8 含 `_test.go`、bans #1-5 不含），**后来的运行会把 `ban #8 internal/`／`cmd/` 各读高 1 枚**，而 `228/38` 那两行仍不动。本程**没有**为这枚预测再跑一发门禁（不给在飞腿的半成品出颜色），只把机理与时刻写清楚。

## ⑨ 判不动的地方与六格三档判语

### ⑨.1 三档判语（六格一次裁，每格指回本程自己的尺）

| 格 | 被验对象 | 判语 | 一句话凭据（本程现跑） |
|---|---|---|---|
| 1 | **AC#4 正控还有没有牙** ＋ "ban #9 只读 `f.Comments`" 这句话成不成立、会不会把正控做瞎 | **成立**（两半都成立，且**不构成互相抵消**） | 种一发注释引用盘上不存在的件 ⇒ `rc=1` ＋逐字 `comment cites repo path "docs/evidence/s1/212-v3-plant-phantom.md" ...`；删掉载体 ⇒ `rc=0` 零 phantom（§② P1/P10）。`f.Comments` 那句**为真**（§② S1 字符串字面量 `rc=0`），但它**做不瞎 AC#4**——票面 AC#4 种的正是"注释"，块注释也响（S2）；被它做瞎的是**另一格**：`tools/**` 扩射程那一格（S4 已证 out of range，且 S1 证明扩了射程也看不见字符串里的旧幻影） |
| 2 | **ⓐ 豁免的真实射程**：全拼必红／`...` 全拼静默；有没有"谎报伪装成缩写"的真实写法 | **带条件成立**（条件＝豁免面比"省略号中段"宽，本程**在盘外真树造出两形真实写法并复现静默**，而**真仓当前活体 0 枚**） | 双向都有：P1 全拼 `rc=1`；P2 `docs/evidence/s1/212-v3-...-phantom.md` `rc=0`。溜过①的伪装形：P9 尾部 `...`（"等等"式）、P3 尾部 `…`、**P4 markdown 加粗 `**docs/evidence/s1/212-v3-bold-phantom.md**`**（＝把谎报写成"强调"）、P5 CJK 粘连、P12 `*.md` 通指——全部 `rc=0`；边界同时钉住：P6 粘连无标记 `rc=1`、P8 标记换行 `rc=1`、P11 标记前有空格 `rc=1`。真仓现量（我自己写的 ban#9 同语义普查，`logs/census-real-repo.txt`）：RING=**0**／EXEMPT=**1**（`cmd/wisp/panel_host_windows.go:29` 的 `docs/evidence/s1`，逐字是真缩写不是谎报）／SYMREF=3 ⇒ 三形的**活体仍是 0 枚**，编排者"不收窄"那道裁定在本程读数下**不被推翻** |
| 3 | **前缀同步钉是不是同义反复** | **成立（有牙，不是同义反复）** | 两枚正则是**两份独立字面量**（`main.go:891`·`:905`），钉用 `regexPrefixAlternation` 各读各的 `.String()`；**我自己的刀**＝给 `shorthandRegionRe` 加 `plugins\|`（腿加的是另一枚表），⇒ `scan_test.go:2441` 逐字 FAIL 且**指向第二臂**（`only in shorthandRegionRe=[plugins] ... widens the exemption without an owner's approval`），同棵树 `go run . -self-test` 仍 `all 40 direction checks passed`；还原 md5 全等 ⇒ 真实漂移看得见。**条件具名**：这枚钉只在 `go test` 那一侧有牙——我已复认它**确实在门禁第一步与 CI 的 `D22 scanner positive control` 步里执行**（`scripts/d22scan.sh` 调 `tools/d22scan/runtests.sh -C tools/d22scan ./...`；`ci.yml:74-81` 同一条；§⑧ 那发 `--- PASS` 在档），所以"只有本地跑测试才看得见"这一支**不成立** |
| 4 | **AC#5 八枚 ban 读数逐名不变**＋两枚 +1 归因 | **成立**（八枚逐名与编排者 10:40、与腿 §7 终值逐字相同；两枚 +1 **不是 `212-r3` 的伤**，且**腿对 `internal/` 那一枚的归因文件指错了，具名更正**） | `ban #8 internal/` 500→501 ＝ `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go`，随 **`e3e19e8e`（260-r2，10:36）** 新增（尺＝`git log --name-status --diff-filter=A b1e59d63..HEAD -- internal cmd`＋`git ls-tree -r --name-only` 两锚点作差，档 `logs/added-go-between-anchors.txt`）；⚠ 腿 §7 猜的是 `internal/ball/hotkey_live_test.go`（那枚是 **M 不是 A**，不改变分母）⇒ 结论方向对、**点名文件错一枚**。`ban #8 cmd/` 97→98 ＝ `cmd/wisp/resident_approval_risk_256_windows_test.go`，随 **`2a10134e`（256-r1，10:20）** 新增，与腿一致。两枚都是 `_test.go`＝只有 ban #8 含、bans #1-5 不含，所以 `228/38` 未动。另：`212-r3` 在 `internal\|cmd` 内**只 M 了 `pending_read.go` 一枚、零新增文件**（同一把尺），结构上就不可能移动分母 |
| 5 | **(a) 注释"诚实化"引用的每一枚 `file:line` 真不真** | **成立**（四枚逐枚打开对全中，一处未把本票的缺陷再犯一遍） | `internal/tools/gate.go:12-13` 逐字只点名 `RulesHit and Reason VERBATIM`＝**两枚**✓；`cmd/wisp/panel_pump.go:44-46` 逐字 `Level/RulesHit/Reason/SessionOverrideBlocked`＋`a copy, not a judgement`＝**四枚**✓；`internal/panel/approval.go:72-88`＝`CardViewFromDecision` 函数体（`:72` 起、`:88` 收到结构字面量右括号），四枚分别在 `:81/:82/:83/:85` 上视图✓；`internal/agent/approval/ticket146_liveapprovals_backing_test.go`（同段旧句）`find` 命中存在✓。新句另引 `tools/d22scan`／`risk.PathResolver`／`filepath.Clean/Abs` 均为符号不是路径✓。**没有一枚是幻影** |
| 6 | **三枚新样本是不是恒真／装饰** | **成立**（三枚都有牙，且各响在自己声称的那一形） | 我自己的 m1（`short := map[int]bool{}`）⇒ `4 direction(s) failed, 36/40 passed - the gate does not see what it claims`，三枚新样本**逐枚翻红且 token 各不重叠**：`cites-tail-ellipsis.go:4 ... "docs/evidence/s1/probe-tail-ellipsis.md"`／`cites-tail-star.go:4 ... "docs/evidence/s1/probe-tail-star.md"`／`cites-cjk-attached.go:5 ... "docs/evidence/s1/probe-swallowed.md见docs/evidence/s1"`，第四枚是 r2 那枚 `cites-shorthand.go`（三 token）；还原后 `all 40 ... (20 expect-ring, 20 expect-silent)` rc=0 ⇒ **没有一枚摘掉豁免仍不红** |

**其余两笔（不在六格里，本程顺手量到，具名入账）**：
- **(c) 没有变成新幻影**：`ls -l docs/evidence/s1/212-comments-phantom-citation-v2.md`＝**33,578 字节、mtime Oct 4 09:57**，复认编排者③；`grep -rn "212-citation-ruler" --include=*.go cmd internal tools` 只命中 `selftestsamples.go:381` 的 **`note:` 串**（那一句自己写明"该件从未存在"＝历史说明，不是路标），**产码注释零命中**⇒ 复认编排者④。
- **行为零变化这一条我另有一把更硬的尺**：`git show d03d166f -- tools/d22scan/main.go` 与 `-- internal/agent/approval/pending_read.go` 的**全部改动行都是注释行**（`grep -E '^[+-]' | grep -v '^[+-][[:space:]]*//'` 两枚文件各 **0 行**，档 `logs/diff-main-noncomment.txt`／`logs/diff-pendingread-noncomment.txt`）；`scan_test.go` 那一枚 hunk **86 增 0 删**（未放宽、未删除任何既有断言）；`selftestsamples.go` 删除行只有 3 枚，逐枚都是 (c) 那一句的 `src`／`summary`／`note`。⇒ 腿自述"行为零变化"**成立**，且比它自己给的 `--numstat` 版更强。

### ⑨.2 判不动／没跑到的地方（具名＋归口，⛔ 不写"应该没问题"）

| 枚 | 判不动的东西 | 为什么本程判不了／没跑到 | 归口 |
|---|---|---|---|
| 1 | 把 `tools/**` 纳入 ban #9 射程到底值不值 | 本程量到了**代价的形状**（S4 证 tools/ 现在完全不扫；S1 证扩射程也扫不到字符串字面量里那枚旧幻影；粗分母另算），但**扩射程＝人工批准面**，不是验收腿能裁的 | 编排者：随票 264 同批摆给机主（`A592` §4 已写死"不单摆一次"） |
| 2 | 把 `_test.go` 纳入射程 | 同上；唯一活体缩写形 `cmd/wisp/slo_report_144_windows_test.go:851` 本程**未复跑**（不在三格内，且它与谎报形不同＝它是缩写不是完整名） | 与 1 同批 |
| 3 | 豁免面会不会被"新前缀"绕过（除枚 4 之外） | 那需要真有人往两枚表里**只加一枚前缀**；我已用 m2 复现"只加一枚表"的形状并证钉响，**"两枚都加但拼法不同"（如 `.plugins` vs `plugins`）本程没做**——钉是**字符串集合**比较，那种拼法差异它会红（同一把尺），但红句会不会指对"哪一枚该改"我没验 | 若真发生漂移由那次改动自己的验收腿判 |
| 4 | 台账两条散文过期（`docs/reports/pending-and-issues.md:11408`·`:7832` 记 `37 direction checks`） | 台账**只追加不删**，且不在本程写面；本程终值 40＝20/20（§⑧） | 编排者下次记台账时带口径与锚点（腿 §6 同判） |
| 5 | NTFS 大小写洞（P10 静默／P11 响）与 `symRefRe` 三枚排除的**语义**是否都该排除 | 不在票 212 的三格内，本程**只复算了枚数**（SYMREF=3）没逐枚判该不该排除 | 已知残余，`A592` §4 ③ 已具名归口 |
| 6 | `internal/probe/*` 夹具样本的**真实产品后果**（夹具不是真树） | 三枚 silent 样本钉的是**豁免本身**，真树里那三形的活体是 0（§② census）；"豁免会不会某天盖住真谎报"只能等活体出现，本程**造得出形、造不出需求** | 若将来抓到活体 ⇒ 新票，⛔ 不许顺手改 regex |
| 7 | 名册行/自述里"扩面＝批准面"这类**治理句**是否与 `SPEC-12 §4.1` 逐字一致 | 本程核的是**射程与行为**，没有把每一句治理措辞回抄到 spec 逐字比 | 编排者；本程不判"措辞合规" |
| 8 | 票 212 AC 框与 `-done` 改名 | AC 框与改名归编排者（⛔ 本程一枚没碰，`git log -1 -- .scratch/wisp/issues/212-*.md` 见 §①） | 编排者按 §⑨.3 办 |

### ⑨.3 票 212 现在够不够格改名 `-done`

**够格**（本程判语）。判据逐条：
- ⛔ 编排者写死的那一枚否决条件＝"只要 AC#4 那格被你自己重跑推翻，就不许 `-done`"——**本程重跑没有推翻它**：种⇒响（`rc=1`＋逐字 token）、删⇒不响（`rc=0`），两向都在**盘外真树**用**同一枚产码 binary**复现（`logs/plant-p1-*.txt`／`logs/plant-p10-*.txt`）。
- 三笔欠账逐笔有下落：**(a)** 出处指对（§格5 四枚 `file:line` 全中）、**(b)** 三形各一枚 silent 样本且**摘掉豁免逐枚翻红**（§格6）、名册行与 `main.go` 自述写出真实射程（§格2/格1，本程用 10 发种删复现，未见"读起来比实际窄"的句子残留）、**(c)** 改指真实件且**没变成新幻影**（33,578 字节现读）。
- 加的那枚钉**有牙**（§格3 我自己的刀），且**没有放宽／删除任何既有断言**（`scan_test.go` 86 增 0 删）。
- 本程唯一**必须更正腿**的一处不是伤，而是归因点错文件：`ban #8 internal/` 那枚 +1 的真身是 `hotkey_cancel_borrow_expect_260r2_test.go`（随 `e3e19e8e`），不是腿猜的 `hotkey_live_test.go`——**结论不变（不是 212-r3 的伤）**，请编排者在翻勾节带上这一处更正。
- 另请带上**一枚措辞小账**（纯注释、零行为，改法写在 §⑥ 末）：`main.go:67-68`·`:825-826`·`:927-929` 的臂编号（1＝中段标记）与 `selftestsamples.go:425/437/450` 的 `SHAPE 1/2/3` 标签（1＝尾部 `…`）**错位一格**，"one case per arm" 那句也因此不准；**覆盖本身是齐的**（4 枚 expect-silent 用例，本程 m1 逐枚验红）。
- 本程自证没越权：AC 框 `- [x]`＝**5 枚**／`- [ ]`＝**0 枚**（与起手一致，本程一票未碰票面，尺见 §①）；票文件名仍无 `-done`；禁区六枚路径与 `tools/d22scan/allowlist.txt` 在 `git status --porcelain -- tools internal cmd scripts docs .github` 里终态只剩并发腿 260-r3 的两枚（§⑧）。
- 建议随 `-done` 一并带的两条**已知残余**（⛔ 不是退回理由）：`tools/**`／`_test.go` 是否在 ban #9 射程＝待人批准（§⑨.2 枚 1/2）；台账 `37 direction checks` 两条散文过期（枚 4）。

