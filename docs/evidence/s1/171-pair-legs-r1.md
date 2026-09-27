# 票 171 r1（写码位·仪器程）＝AC#2（枚数进判据）＋ AC#6（合方词根认句柄那一形）——同一枚 `pair` helper，一程做完

- 锚点（step 0 现量）：分支 `dev`｜`git rev-parse HEAD` ＝ `19513ccefa45d51206871f00b3efa6ab19dd620a`｜时刻 `2026-09-27 15:43:09 CST +08`
- 写面：`.scratch/wisp/probes/154/gate-clauses.sh`（唯一一枚既有文件）· `.scratch/wisp/probes/171/r1/**`·本件·票 171 的 Progress log（最后一步、只追加）
- **生产码零字节**：`internal/**`、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、`probes/160/**`、`probes/161/**`、`probes/162/**` 全 0 改动（逐枚现量命令与读数见 §7 末）
- step-0 第 4 件的一条偏差（先报回）：`git status --porcelain -- .scratch/wisp/probes/154/ .scratch/wisp/probes/161/ internal/ cmd/` **不是全空**，三行是 `probes/161/` 下的既有未跟踪件（`r2/__pycache__/`、`r2/ctl/`、`r5/negative-control/`），**不是 tracked 改动、不是本程造的**；`probes/154/`、`internal/`、`cmd/` 三处逐字为空。本程没碰它们。
- ⚠ 取数之后 HEAD 被**另一程（票 162 的验收程）推进过**：`git log --oneline 19513cce..HEAD` 今天实量 6 枚，`git diff --stat 19513cce..HEAD -- internal cmd` **＝空** ⇒ 两把尺的射程里一枚 Go 文件都没变，本程所有读数仍然对得上；改前/改后两读都明写锚是 `19513cce`。


## 1. 票面两格（连行号逐字抄自 `.scratch/wisp/issues/171-…-pair-legs.md`）

```
27	- [ ] **AC#2（尺洞③＝枚数）**：把成对腿的判据从"空／非空"改成**带枚数**（声明"该响几枚"或"枚数只增不减"，两形选一形并写明为什么）。**硬判据**：造一发"同一条腿内未成对从 1 枚变 2 枚"的台件 ⇒ 聚合退码**必须变**；这发在**未修码上必须不响**（否则本格是装饰）。
33	- [ ] **AC#6（09-27 11:5x 追加，来路＝票 160 的 r1 落地后编排者现量复现）合方词根写死成 `CloseScope` ⇒ 生产码换成"句柄自带 `Close()`"那一形之后，G5 的正控那腿言行不一、聚合退码从 0 变 1，而主尺多假点一枚。**
37	  ① 修完 `gate-clauses.sh` 的**聚合退码回到 0**（14 腿零枚言行不一），且这一发必须是你真跑出来的、不是推的；
38	  ② **不许靠把 `internal/tools/bridge.go` 从 G5 主尺射程里排除掉来凑 0**——那是放水，且是把票 161 AC#7 刚收掉的恒红形状从另一扇门放回来（同一枚错，台账 `A320`／`A321`）；改射程必须与"认得句柄那一形"**同时**发生，两件事一起做才有意义；
39	  ③ **反向判据（这一发在未修码上必须不响）**：造一枚"用句柄形状**开而不关**"的样本放 `.scratch/wisp/probes/171/**`（不进树）⇒ 未修码上主尺**点不到它**（因为词根不匹配，今天不响）、修完之后**必须点名它**。答不出"哪条用例因此变红"＝这一格是装饰。
```

## 2. AC#6①：改前／改后逐字（同一条命令、同一个锚）

改前（台件 `.scratch/wisp/probes/171/r1/baseline-gate-clauses-pre-fix.txt`，`sh .scratch/wisp/probes/154/gate-clauses.sh; echo rc=$?`）：

```
# BAD  腿=G5pos 声明=quiet 实测=ring
# 腿数＝14 声明与实测不符＝1
# 聚合退码＝1
```

改后（台件 `.scratch/wisp/probes/171/r1/readings-gate-post-fix.txt`，同一条命令）：

```
# ok   腿=G5 声明=ring 基线=1枚 实测=1枚
# ok   腿=G5pos 声明=quiet 基线=0枚 实测=0枚
# ok   腿=G5neg 声明=ring 基线=8枚 实测=8枚
# 腿数＝14 声明与实测不符＝0
# 聚合退码＝0
```

⇒ **由红转绿**：rc 1→0，14 腿零枚言行不一，两行都是真跑出来的。

改因的现量（不是推的）：**`CloseScope` 这个词根在生产码里已经没有任何调用点**——
`git grep -n CloseScope 19513cce -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go'` ⇒ 两枚命中，**两枚都是 `internal/risk/provenance.go` 的注释行**
（`:361`、`:562`，那两行讲的正是"旧的两口子里 by-id 那一口已经拿掉了"），调用行 0 枚；而 `pair()` 本来就剔整行注释 ⇒ 合方名册结构性为空。
（同一条命令的原文输出与 rc 留在 `.scratch/wisp/probes/171/r1/readings-write-surface.txt`。）
票 160 之后唯一的收尾是句柄自带的那枚：`git grep -nE 'OpenScope|scope\.Close\(\)' 19513cce -- internal/tools/bridge.go`
⇒ `:648 b.scopes[taskID] = b.prov.OpenScope(taskID)` 与 `:706 closeErr = scope.Close()`（改前 `:211` 那句注释写的 `:642/:691` 已因这次落地位移，同一格注释本程已改成现量数字）。
所以一把合方只写 `CloseScope` 的尺，今天会把**任何**用句柄正常关掉的文件读成漏。

## 3. AC#6③：反向判据（三枚样本、两把尺、两向读数）

台件＝`.scratch/wisp/probes/171/r1/r1-pair-shapes.sh`（全份 stdout：`readings-pair-shapes-r1.txt`，逐腿 log：`logs/part1-*.txt`、`logs/part2-*.txt`）。
样本**不进主树**：脚本在 `$TMPDIR` 里 `git init` 一枚丢弃树（`internal/fake171/*.go`，落在尺的 pathspec 射程 `internal/**/*.go` 内），
改前的尺＝`git show 19513cce:.scratch/wisp/probes/154/gate-clauses.sh` 写进 `$TMPDIR`，两把尺打**同一枚 commit**，
且两把都跑在**同一条射程**上（一条 `:!` 都没减）。这枚丢弃树在本仓目录之外，不是 worktree、不是 checkout、没切分支。

| 样本 | 形状 | 改前的尺 | 改后的尺 |
|---|---|---|---|
| A `a_paired.go` | `k.scopes["a"] = …OpenScope("a")` ＋ 同文件 `…Close()`（＝bridge.go 那一形，真的关了） | **点名**（rc=1，假红） | **不点名**（rc=0），且打印 `CLOSE-BY-GENERIC-CLOSE internal/fake171/a_paired.go`＝"这文件的关是谁关的"摊给人判 |
| B `b_leak.go` | 接住句柄、**全文件一次都没关** | 点名 | **点名**（rc=1） |
| C `c_discarded.go` | `_ = p.OpenScope("task-C")`（句柄被丢）＋ 同文件 `l.CloseScope("task-other")`（关的是别人的 id） | **不点名**（rc=0） | **点名**（rc=1），并打印 `DISCARD-HANDLE` ＋ `note 丢句柄那一味点的名` |

- **③ 要的那一发（未修码不响 → 修完响）＝样本 C，两向都真跑出来了。**
- **③ 的字面前半句对样本 B 不成立，报回不改判据**：票面写"未修码上主尺点不到它（因为词根不匹配）"，实测 B 在改前的尺上**就被点名**（rc=1）。理由是结构性的：AC#6② 禁止减射程，而"认句柄那一形"是把**合方**的词根集合变大；同一射程、同一开方之下，合方集合变大只能让**正向名册变小**，数学上点不出改前点不到的文件。真正"改前点不到"的形状是**改前的尺会当成成对**的那一形＝C（开方与合方的词根都在场、句柄却已经被丢）。本程按实测写判据，没有为了对上门面那句话去动尺。
- 附带一条自纠（本程实测撞出来的，不是设想的）：`discarded_files()` 第一版把**定义行** `func (p *Provenance) OpenScope(...)` 也算成"裸调用丢句柄"，于是样本 A 被误点一枚；定义行天生没有赋值号，拿声明当调用点就是噪音。现版本先 `grep -vE ':[0-9]+:[[:space:]]*func[[:space:]]'` 再判，log `logs/part1-A-fixed.txt` 与上表 A 行是修好之后的读数。

## 4. 我有没有减射程（逐枚集合差，不是一句话）

台件＝`.scratch/wisp/probes/171/r1/roster-diff.sh`（读数 `readings-roster-diff.txt`，九腿的 pre/post 名册落在 `logs/roster-*-pre.txt` / `*-post.txt`），命令：

```
sh .scratch/wisp/probes/171/r1/roster-diff.sh \
  .scratch/wisp/probes/171/r1/baseline-gate-clauses-pre-fix.txt \
  .scratch/wisp/probes/171/r1/readings-gate-post-fix.txt
```

逐字读数：

```
== ## G5 OpenScope    改前 2 枚／改后 1 枚   退场: - #   UNPAIRED internal/tools/bridge.go   新增: (空)
== ## G5-正控         改前 1 枚／改后 0 枚   退场: - #   UNPAIRED internal/tools/bridge.go   新增: (空)
== ## G5-负一负       改前 9 枚／改后 8 枚   退场: - #   UNPAIRED internal/tools/bridge.go   新增: (空)
== ## G6 主尺／G6-正控／G6-负一负／G7 主尺／G7-正控／G7-负一负：退场＝(空)  新增＝(空)
```

⇒ **一条 pathspec 没减、一枚 `:!` 没动**（写面自证：`git diff -- .scratch/wisp/probes/154/gate-clauses.sh | grep -E '^[-+].*pair "` 只有 `'CloseScope'`→`'CloseScope|Close\(\)'` 与多出的 `handle`/`-` 一味，行尾的 pathspec 逐字未变）。
全门**新增点名＝0 枚**、**退场＝1 枚文件 × 3 条腿＝同一枚 `internal/tools/bridge.go`**，而它退场的理由不是"被排除了"，是"它确实把接住的句柄关了"（`:648` 接、`:706` 关），并且改后的尺**仍把这枚文件摊出来**：`#   CLOSE-BY-GENERIC-CLOSE internal/tools/bridge.go`（谁关的？泛用 `Close()`——人来判是不是这族句柄）。

这条批注不是装饰：`Close()` 在本仓是泛用动词（现量：非测试生产码 42 枚文件里出现 `Close()`，命令 `git grep -lEw 'Close\(\)' 19513cce -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l` ＝ 42，原文在 `readings-write-surface.txt`），
一把只按文件求差集的尺分辨不出**那一只** `Close()` 是不是 scope 的句柄。所以合方认两形换来的每一分安静，都同时打印成一句可判的话；
真正的兜底是 §3 的样本 C 那一味（丢句柄 ⇒ 合方写什么都不算关），它让"用句柄开而不关"比改前**更**看得见，不是更看不见。

## 5. AC#2：枚数进判据——选了哪一形、为什么、两向读数

**选的形**：`want_n <基线枚数>`（紧跟 `want`，只对 pair 腿有效）＝"这一腿今天实测几枚"，判据是**两味的并集**、一枚腿最多计一进退码：
① 票 161 AC#6① 的空／非空那一味**逐字保留**；② 新增 `实测枚数 > 基线 ⇒ BAD`。`实测 < 基线` 打 `SHR`（读数变好＝不算言行不一，但基线过期，逐行摊出来给下一程核）**不计进退码**。

**为什么不选"该响几枚＝实测必须逐枚相等"那一形**：相等形式会让**改好了**也变红（把 `panel_assets.go` 补上一发收尾 ⇒ G5 从 1 枚变 0 枚 ⇒ 恒红）。本仓已经付过这笔钱：票 161 AC#7／台账 A320、A321 收掉的正是"门为不可修的东西常红 ⇒ 下一程就去把它放宽"。而 AC#2 要补的洞只在**增加**那一侧（票面原话："新增的第二枚躲在非空里看不见"），所以响的方向就只往增加那一侧装。基线本身是**可读的数**、写在尺旁边，谁把它调大都在 diff 里露头。

**硬判据的两向读数**（同一枚 `r1-pair-shapes.sh` PART 2，同一条腿 G5、同一射程、同一枚 `want_n 1` 基线）：

```
phase1（1 枚未成对）：G5 实测 改前=1枚 改后=1枚 | 聚合退码 改前=6 改后=6
phase2（2 枚未成对）：G5 实测 改前=2枚 改后=2枚 | 聚合退码 改前=6 改后=7
# BAD  腿=G5 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）
  ok    UNMODIFIED CODE: 1 枚 -> 2 枚 moves the aggregate exit code NOT AT ALL (rc=6 both)
  ok    FIXED CODE: the same jump moves the aggregate exit code (6 -> 7)
```

⇒ **未修码不响、修完才响**，两向都是现量。洞的成因也在同一份读数里看得见：改前的尺**也数得出** 1 枚/2 枚（它把 `# rc=` 打出来了），只是判据不读那个数。

**G6／G7 那几腿有没有被我改过期**：没有——同一枚 helper 的改动会连带过期，所以九枚 pair 腿**全部**登记了 `want_n`，改后逐枚现量：`G5=1／G5pos=0／G5neg=8／G6=1／G6pos=0／G6neg=1／G7=3／G7pos=0／G7neg=4`（§2 那份 post 读数的九行 `基线=N枚 实测=N枚` 全等）。
G6／G7 的**算术一字未动**：§4 的 comm 对它们逐腿空差集就是这条的证。它们只多了两列（`基线=`、`因=`／`注=`）。
另外 `pair()` 现在**拿不到基线就直接死**（rc=4，"枚数不进判据＝票 171 AC#2 那枚洞还开着"）——下一程新加一枚腿不可能悄悄退回空／非空那一形。
`--diffsets` 文本喂入那一味（票 161 r6 的 `diffsets-text.sh` 吃的就是它）**逐字未改**：`diffsets()` 的第四味是可选参数，不传＝行为与改前一致。

## 6. flip-declaration.sh（声明机制还活着没有）

- 改前（`baseline-flip-declaration-pre-fix.txt`）：**RED**，rc=1，三格 FAIL，全是 AC#6 那枚常红的下游：
  `FAIL 14 legs each say what it does -> rc=1 (want zero)`／`FAIL G1 quiet->ring AND G4 quiet->ring together -> rc=3 (want 2)`（多出那一枚正是 G5pos）／`FAIL nothing flipped any more -> rc=1 (want zero)`。
- 改后（`readings-flip-declaration-post-fix.txt`）：**GREEN，rc=0**，九格 ok、零枚 FAIL，逐字读数：

```
  ok    nothing flipped any more -> rc=0 (want zero)
  ok    restore matches the baseline (rc=0)
flip-declaration.sh: GREEN - the aggregate exit code tracks the declarations in both directions and returns to 0 (baseline rc=0, 14 legs booked).
```

⇒ **由红转绿**，而且它证的正是本程最容易改坏的那件事：翻一枚声明退码跟着变（`G2 ring->quiet`、`G5neg ring->quiet`、`G7 ring->quiet`、`G1 quiet->ring`、`G7pos quiet->ring`、两枚同翻＝rc=2、翻回＝rc=0）。
`want_n` 之所以没把这套声明改成"永远说到做到"的空尺，是因为 `book_count()` 里票 161 那一味（空／非空 vs `ring`／`quiet`）**逐字保留**，枚数只是**多出来的一味**——flip 6 那发（两枚同翻＝rc=2 不是 1）就是这条的证。
⚠ 一枚副作用要报：`flip-declaration.sh` 自己把逐腿 log 写在 **`.scratch/wisp/probes/161/r6/logs/flip-*.txt`（8 枚 tracked 件）**，
所以跑这枚门禁会让 `probes/161/` 出现 ` M`——那是门禁脚本的落点，不是本程改了 161 的尺；本程**没提交、没还原**那 8 枚（见 §9）。

## 7. 门禁四数＋名册差集＋gofumpt

| 门禁 | 改前 | 改后 |
|---|---|---|
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | rc=1（`腿数＝14 声明与实测不符＝1`） | **rc=0**（`腿数＝14 声明与实测不符＝0`） |
| `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` | RED（3 格 FAIL，全部下游于上面那一枚） | §6 下格 |
| `sh scripts/d22scan.sh` | rc=0（clean，live scope work 九项计数见 `baseline-d22scan.txt`） | rc=0，`internal/=418`、`cmd/=45`、`frontend/=85`、`design/=39` 逐字同改前 |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | `packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` rc=0 | 四数逐字相同，rc=0 |
| runtests 名册两向 comm（34 枚顶层＋76 枚 RUN 的名册，`grep -oE '^\s*--- (PASS\|FAIL): [A-Za-z0-9_/]+'` 去重排序） | 66 枚 | 66 枚；`comm -23`＝空、`comm -13`＝空 ⇒ **名册逐字相同** |
| gofumpt | — | `gofumpt --version` ＝ `v0.12.0 (go1.27.1)`（不在 PATH，按派单 `export PATH="$PATH:$(go env GOPATH)/bin"` 后现跑）；本程**没有在仓里落任何 `.go`**（样本是脚本在 `$TMPDIR` 生成的），生成物逐枚过 `gofumpt -l` ⇒ `ok generated samples are gofumpt-clean` |

写面自证（禁改名单零字节；命令与输出留在 `readings-write-surface.txt`）：`git status --porcelain -- internal cmd tools .github frontend design docs/PLAN.md docs/specs docs/reports .scratch/wisp/probes/160 .scratch/wisp/probes/161 .scratch/wisp/probes/162`

## 8. 被拒／没成功的调用

- 发生在**取数之后**（step 0 五件与改前基线三份读数都已落盘）：
  1. 一次 `Bash` 因单行 shell 语法写坏（`for`/`awk` 混排漏了 `done`）exit 2，未产生副作用 ⇒ 改成台件脚本 `roster-diff.sh` 重做（同一格因此多一枚文件，是台件不是残渣）。
  2. 一次 `Edit` 因 tab／缩进不匹配 0 命中，重试一次即中；期间在 `pair()` 留过一段重复的旧代码（同一次 Edit 的 old_string 止于函数中部），已在下一次 Edit 里删净，`bash -n` 与两轮实测读数都在删净之后跑的。
- 取数之前：无。

## 9. 有没有跑过删除命令

- 仓内：零次 `rm`／`git rm`／`clean`／`checkout`／`restore`。
- 仅在 `$TMPDIR` 的丢弃树里有 3 次 `rm -f`（`r1-pair-shapes.sh` 用来在样本之间切换状态，落在本仓之外；派单/161  precedent："Deleting inside this synthetic tree is the ONE deletion this dispatch allows"）。
- **本程唯一一处"写进了别人的写面"是门禁自己干的**：跑派单点名的 `probes/161/r6/flip-declaration.sh` 会往 `probes/161/r6/logs/flip-*.txt`（8 枚 tracked 件）写 log。
  那 8 枚的 ` M` 因此出现在 `git status` 里；本程**没提交、没还原、没改内容**，处置与下面那些未提交件同一条。
- 别人的未提交件（`design/**` 的 16 枚删除与几枚未跟踪新件、`probes/152/my152.py`、`docs/evidence/s1/152-…-accept-r1.md`、`probes/162/r3/`、`docs/evidence/s1/162-fs-edit-r3.md`、`probes/161/r2| r5` 的既有未跟踪件）**不提交、不还原、不补完、不评论**，逐字保持现状。

## 10. 伪授权两栏

- **我没有做过、也没被授权去做**的事：改契约（D/C/R/D43 转移表）、动 SLO／golden／`thresholds.go`、动禁令射程与 `allowlist.txt`、把 G5 接进 CI、切分支／建 worktree、push、改写已推送历史。派单原文里没有任何一句把这些写成本程范围。
- **看起来像授权、我按"不是"处理**的事：① 票面 AC#6③ 那句"未修码上主尺点不到它"与实测不符（§3 表 B 行）——按派单 §5 末段"不符就报回、不许为了对我那句话去改判据或改测试"处理：**没有**为了让 C 之外的形状"对上"而放宽任何一形；② 编排者说"flip-declaration 现在还活着"——改前现量它是 RED，我按读数把它当**必须转绿的那一发**跑，而不是当参考（改后读数在 §6）；③ `probes/161/**` 下三枚未跟踪件——按"只读、别动"处理，没把它们当谁丢的垃圾清掉。

## 11. 凭据

本程没有读过、也没有抄录任何 API 密钥／DPAPI 密文／token；台件与读数里出现的字符串只有文件路径、commit 号与函数名。

## 12. next=

1. AC#6 的"合方认两形"只做到**文件级**：`CLOSE-BY-GENERIC-CLOSE` 那一批注把"哪一只 Close()"的问题交给人，没交给尺——把它接成可判的东西要做数据流，超出 git grep 那把尺的地界（要不要另立尺，归票 154/158 的地界裁决）。
2. AC#2 只补了 pair 腿；`run()` 那一形（G1–G4）的"名册枚数"仍然只有 0／非 0——同一枚洞的另一半，本程**没动**，因为它不在派单的两格里。
3. 票 171 的 AC#1（尾随注释被当调用点）与 AC#3（自检接 CI）两格本程**没做**；§3 那条"泛用 Close()"的批注会让 AC#1 的射程变长一点，做 AC#1 的那一程要一并把 `CLOSE-BY-GENERIC-CLOSE` 的行形状算进名册 diff。
4. 基线枚数是**会过期**的数：G5neg=8 里有 5 枚是 `internal/risk/**` 的测试与定义本体（`want G5neg ring` 的射程本来就不排它们）；下一程若给 G5-负一负加排测试，记得同一格里把 `want_n` 核下来，别留 `SHR`。
5. 本件的 §6 若被下一程读到还没有读数＝那一发没跑成，不要当成绿。
