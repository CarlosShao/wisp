# 票 171 · 171-r2 写码位（仪器程）证据件＝AC#7（`run()` 那一族腿带上枚数）＋ AC#8（减法侧装牙三处）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-173x-impl-171-r2-run-legs-count-and-subtraction-teeth.md`
- 票面＝`.scratch/wisp/issues/171-…-count-blind-pair-legs.md`（本格引的行号：**AC#7＝`:50`、AC#8＝`:52`（三处最小闭合集合逐字照抄，未改写）、AC#8 的「顺带一字」＝`:53`**；已结案两格的 `>`＝AC#2 的 `:28`–`:29`、AC#6 的 `:44`–`:48`）
- 本程 step-0 锚点（现量）＝**`f1b99a70c3dfdad04f34094306b5942a06355f52`**，分支逐次现量＝`dev`（起手 `git rev-parse --abbrev-ref HEAD` ⇒ `dev`；结案时再量一次仍 `dev`）
- 本程性质：**仪器程**。写面只有 `probes/154/gate-clauses.sh`＋`probes/161/r6/flip-declaration.sh`（只新增一格）＋`probes/171/r2/**`＋本证据件＋票面 Progress log 一行。
- ⚠ **共享工作树里另一程 `162-r4` 正在动 `internal/tools/**`**：本程取数期间它落了 `7a41af51`／`27f0cb8a`／`9c00bb33`／`b542bc1c`／`fd8201e7`／`bb840a4f` 等枚，`git rev-parse HEAD` 因此在程中前移过。所以本程**所有"改前 vs 改后"的比对一律把两把尺打同一枚锚**（§5 那台 `r2-roster-diff.sh`），不拿两次不同 HEAD 的读数比。

---

## 1. step-0 起手五件（逐件带命令）

| 件 | 命令 | 读数 |
|---|---|---|
| 时刻 | `date -Iseconds` | `2026-09-27T17:29:10+08:00` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev`（＝派单要求，不是＝停手；全程未用 `switch`/`checkout`/`merge`/`rebase`/`reset`/`stash`/`worktree`/`clean`） |
| 锚点 | `git rev-parse HEAD` | `f1b99a70c3dfdad04f34094306b5942a06355f52` |
| 现场 | `git status --porcelain -- .scratch/wisp/probes/154/ .scratch/wisp/probes/161/r6/flip-declaration.sh internal/ cmd/` | **空**（写面起点干净）。工作树里另有 `.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*.txt` 8 枚、`design/**` 16 枚删除＋若干未跟踪件、`docs/evidence/s1/152-…-accept-r1.md` ⇒ **都带"不是我的"状态，本程未提交、未还原、未补完、未评论** |
| 改前基线两份 | `sh .scratch/wisp/probes/154/gate-clauses.sh` ／ `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` | **gate rc=0**、`# 腿数＝14 声明与实测不符＝0`（`logs/gate-pre-fix.txt`）；**flip rc=0 GREEN**、九格全 ok、`腿数＝14 声明与实测不符＝0`（`logs/flip-pre-fix.txt`） |

派单里"flip-declaration 现在九格全绿"、"腿数＝14"、"`SHR` 今天不响"三条前提**都量过了，都成立**；为它们跑的改前读数就是本节这两份，结案时"由红转绿／由绿转红"拿它们比。

## 2. 本程没测什么（不假装覆盖）

- **没跑全仓 `go test ./...`**（派单明令禁止，同树另一程正在写 `internal/tools/**`）。Go 侧只跑 `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` 那一把。
- **没写、也没改任何 `.go`**：本程写面 `find .scratch/wisp/probes/171/r2 -name '*.go' | wc -l` ⇒ **0**；因此"台件必须 gofumpt 干净"这一格的真读数是 **`gofumpt -l .scratch/wisp/probes/171/r2` 输出 0 行**（版本另见 §8，非引用他人读数）。AC#7 的 Go 样本只活在 `$TMPDIR` 的合成树里（不落仓、不求编译——尺吃的是 `git grep` 的行，不是 `go build` 的图；这一条在样本文件注释里逐字写着）。
- **没核 AC#1／AC#3**（尾随注释被当调用点／自检接 CI）：不在本派单。⚠ 顺带一条与 AC#1 同族的**读数口径**：run 腿的枚数＝**命中行数**，而 run() 的名册**吃整张 git grep 输出**（注释行也算）——`cmd/wisp/run.go` 那一发 CloseTask 与它上面那行注释就是 G2 的 2 枚。**这不是本程新造的粗糙**，是本程之前既有形状（pair 腿只吃调用点、run 腿吃全名册）；本程把它的**行数**接进判据，没有把它改成调用点口径——改成调用点口径＝射程收窄，归 AC#1 那一格处理。
- **`runtests` 的改前名册没取**：`git diff --name-only f1b99a70 HEAD -- tools/d22scan scripts | wc -l` ⇒ **0 枚**（本程没碰那把尺的 module）⇒ 名册不可能因本程位移，所以本程只交改后名册 34 枚（`logs/runtests-roster.txt`）＋四数，没假装做过两向 diff。
- **"摘掉一条腿的样本"那一发没单独造**：票面 AC#8 的现量症状是"腿或它的样本被摘"，本程交的是**整腿摘除**（§4 ①）＋**基线高于实测＝样本消失的形状**（§4 ②，真树读数）两发；没有第三发"删样本不删腿"的合成树台件（`next=` 第 2 条）。

## 3. AC#7＝`run()` 那一族腿带上枚数（形制逐字跟 AC#2：`want_n <基线枚数>` ＋"实测 > 基线才 BAD"，不选"逐枚相等"）

**登记用的现量命令**（锚 `f1b99a70`，逐腿一条，与本腿 pathspec 逐字同形）：

```
git grep -nE <pattern> f1b99a70 -- <该腿 pathspec 原样> | grep -c .
```

| 腿 | 声明 | 现量命中行数 | 登记成 |
|---|---|---|---|
| G1 `ToolRequest\{` | quiet | 0 | `want_n 0` |
| G1b `\.Execute\(` | quiet | 0 | `want_n 0` |
| G2 `OpenTask|CloseTask`（`-w`） | ring | **2** | `want_n 2` |
| G3 `^func \(l \*Loop\) [A-Z]…taskID` | quiet | 0 | `want_n 0` |
| G4 `PassThroughUnclassifiedRisk:` | quiet | 0 | `want_n 0` |

**那一发"同一条 `run()` 腿内命中从 2 枚变 3 枚"的两向读数**（台件＝`r2-run-legs-and-subtraction.sh` PART 1；改前的尺＝`git show f1b99a70:…gate-clauses.sh` 的字节，两把尺打同一棵合成树、同一射程）：

| 状态 | 该腿真命中（`git grep` 独立复算） | 改前的尺：聚合退码 | 改后的尺：聚合退码 |
|---|---|---|---|
| state1 | 2 行 | **6** | **6** |
| state2（同一条腿加第 3 枚） | 3 行 | **6**（一动不动＝洞复现） | **7**（动了） |

改后那枚新 BAD 逐字：`# BAD  腿=G2 声明=ring 基线=2枚 实测=3枚 因=新增命中（票 171 AC#7：实测 > 基线）`
改前的尺**根本没有** `# 命中行数＝` 这一行（`leg_field … '# 命中行数＝'` ⇒ `NO-READING`），改后逐状态打 `命中行数＝2`／`＝3`。
合成树里那 6 枚既有 BAD 全是"声明 ring、合成树里未成对 0 枚"那六腿（G5/G5neg/G6/G6neg/G7/G7neg），两状态**逐枚相同**⇒ 退码 6→7 那一枚只能归 G2 自己（原样输出：`logs/harness-pre-fix.txt`／`logs/harness-post-fix.txt`，逐腿名册在 `logs/part1-state{1,2}-{pristine,fixed}.txt`）。

**为什么这一格不是装饰**：判据问的那一问（"今天不响吗"）在**未修码**上先跑了一遍，读数就是上表 state1=state2=6 那一格。

**九腿（连同五枚 run 腿＝十四枚）"腿 → 声明 → 实测"全表，改前／改后各一张**（命令 `sh .scratch/wisp/probes/154/gate-clauses.sh`，锚不同 HEAD、尺不同字节；两向同锚版在 `logs/gate-{pre,post}-at-anchor.txt`）：

| 腿 | 改前（`logs/legs-table-pre.txt`） | 改后（`logs/legs-table-post.txt`） |
|---|---|---|
| G1 | `# ok   腿=G1 声明=quiet 实测=quiet` | `# ok   腿=G1 声明=quiet 基线=0枚 实测=0枚` |
| G1b | `# ok   腿=G1b 声明=quiet 实测=quiet` | `# ok   腿=G1b 声明=quiet 基线=0枚 实测=0枚` |
| G2 | `# ok   腿=G2 声明=ring 实测=ring` | `# ok   腿=G2 声明=ring 基线=2枚 实测=2枚` |
| G3 | `# ok   腿=G3 声明=quiet 实测=quiet` | `# ok   腿=G3 声明=quiet 基线=0枚 实测=0枚` |
| G4 | `# ok   腿=G4 声明=quiet 实测=quiet` | `# ok   腿=G4 声明=quiet 基线=0枚 实测=0枚` |
| G5 | `# ok   腿=G5 声明=ring 基线=1枚 实测=1枚` | 逐字相同 |
| G5pos | `# ok   腿=G5pos 声明=quiet 基线=0枚 实测=0枚` | 逐字相同 |
| G5neg | `# ok   腿=G5neg 声明=ring 基线=8枚 实测=8枚` | 逐字相同 |
| G6 | `# ok   腿=G6 声明=ring 基线=1枚 实测=1枚` | 逐字相同 |
| G6pos | `# ok   腿=G6pos 声明=quiet 基线=0枚 实测=0枚` | 逐字相同 |
| G6neg | `# ok   腿=G6neg 声明=ring 基线=1枚 实测=1枚` | 逐字相同 |
| G7 | `# ok   腿=G7 声明=ring 基线=3枚 实测=3枚` | 逐字相同 |
| G7pos | `# ok   腿=G7pos 声明=quiet 基线=0枚 实测=0枚` | 逐字相同 |
| G7neg | `# ok   腿=G7neg 声明=ring 基线=4枚 实测=4枚` | 逐字相同 |
| 文末 | `# 腿数＝14 声明与实测不符＝0` | `# 腿数＝14 声明与实测不符＝0`＋`# 腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`＋`# 腿数断言＝相符`＋`# 基线过期枚数＝0` |

⇒ **九枚 pair 腿的声明一枚都没悄悄过期**（实测与基线仍逐枚全等），`pair` 腿红句格式一字未动（AC#2 结案记录里那句逐字红句仍然对得上：kind 只改"新增"后面那两个词，算术不多不少）。

## 4. AC#8＝减法侧装牙（三处，票面 `:52` 原话逐字实现）

### ① `legs()`／聚合段加"腿数＝登记期望值"的断言，期望值与九腿 `want` 同源

做法（**只有一处写腿号**）：`LEGS_EXPECT='G1 G1b G2 G3 G4 G5 G5pos G5neg G6 G6pos G6neg G7 G7pos G7neg'`，
14 那枚数不在别处再抄——`LEG_EXPECT_COUNT` 由这张表 `for` 一圈现算（命令：`sh …gate-clauses.sh | grep 腿数断言` ⇒ `名册=14`）。
`want()` 起跑线核三件事：腿号**在不在册**（不在＝rc=3 死）、**有没有第二次出现**（重复＝rc=3 死）；
聚合段核两件事：**在册的每一腿都记过账**（缺＝BAD）、**登记过 want 的腿都真的跑了**（空头声明＝BAD）。

两向读数（台件 PART 2；从**尺的副本**里整条摘掉 G6 的 4 行＝want/want_n/pair/参数行，仓里字节未动；两把尺都打同一枚锚 `fd8201e7`）：

| 尺 | 摘掉 G6 之后 | 聚合退码 |
|---|---|---|
| 改前 | `腿数＝13 声明与实测不符＝0` | **0**（打印 13 然后安静＝洞复现） |
| 改后 | `腿数＝13 声明与实测不符＝1` | **1**，新 BAD 逐字 `# BAD  腿=G6 声明=在册应跑 实测=腿不在场 因=整条腿没跑（票 171 AC#8①：腿数断言）` |

### ② `SHR` 改判＝允许变好，但必须单列一张『基线过期』表，并让 `flip-declaration.sh` 那把核它

两向读数（台件 PART 3；把 `want_n 8` 抬成 `want_n 9`、实测仍是 8）：

| 尺 | 退码 | 单列的表 |
|---|---|---|
| 改前 | **0**（正确：变好了不算言行不一） | **0 枚**（只剩聚合行里那句 `注=读数变好了…`＝只有人眼看得见） |
| 改后 | **0**（仍然不计进退码——没换成"为好消息响"的恒红形状，台账 A320／A321） | **1 枚**，表逐字如下 |

『基线过期』那张表的原样输出（`logs/part3-stale-fixed.txt`，awk 取 `## 基线过期` 到 `# 基线过期枚数`）：

```
## 基线过期（票 171 AC#8②）：实测低于基线的腿，逐枚列成一张表
# 这一味**不**计进退码：变好了不算言行不一（算了就是把门换成「为好消息响」的恒红形状，台账 A320／A321）。
# 改前它只是聚合行里那句「注=读数变好了」，现量＝退码 0、单列表 0 枚（只有人眼看得见）。
# 表头：腿 | 声明 | 基线（want_n） | 实测 | 差 | 处置
# STALE 腿=G5neg 声明=ring 基线=9枚 实测=8枚 差=1枚 处置=把 want_n 核下来；样本一枚没变好而是被摘掉的，先回答谁摘的
# 基线过期枚数＝1   （核它的那一把＝.scratch/wisp/probes/161/r6/flip-declaration.sh 新增的那一格）
```

`flip-declaration.sh` **只新增这一格**（`== (stale) …`，既有九格一字节语义未动），它的**全量读数**（`logs/flip-post-fix.txt`，脚本 `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` ⇒ **rc=0 GREEN**）：

| # | 格 | 读数 |
|---|---|---|
| 0 | baseline（未翻的声明） | `ok … -> rc=0 (want zero)`，`腿数＝14`／`声明与实测不符＝0` |
| 1 | G2 ring->quiet | `ok … rc=1 (want nonzero)`，BAD＝`腿=G2 声明=quiet 基线=2枚 实测=2枚 因=空/非空那一味与声明不符（票 161 AC#6①）` |
| 2 | G5neg ring->quiet | `ok … rc=1`，BAD＝`腿=G5neg … 基线=8枚 实测=8枚` |
| 3 | G7 ring->quiet rev | `ok … rc=1`，BAD＝`腿=G7 … 基线=3枚 实测=3枚` |
| 4 | G1 quiet->ring | `ok … rc=1`，BAD＝`腿=G1 声明=ring 基线=0枚 实测=0枚` |
| 5 | G7pos quiet->ring rev | `ok … rc=1`，BAD＝`腿=G7pos … 基线=0枚 实测=0枚` |
| 6 | G1＋G4 两枚同翻 | `ok … rc=2 (want 2)`，两行 BAD 都在 ⇒ 退码仍是**说谎腿的枚数**而不是一面旗 |
| 新 | stale：`want_n 8`→`9`（实测 8） | `ok … still exits 0`＋`ok the stale baseline is its own readable table row: # STALE 腿=G5neg 声明=ring 基线=9枚 实测=8枚 差=1枚 …` |
| restore | 未翻字节重跑 | `ok … rc=0`＋`ok restore matches the baseline (rc=0)` |

⇒ 这把门禁**不是**被本程改成"永远说到做到"的空尺：它现在**多**核一件事（SHR 必须成表），**没有少**核任何一件（0–6＋restore 七格全部照旧逐格 ok、退码数值照旧）。改前那份（九格、无 stale 格）在 `logs/flip-pre-fix.txt`，同样 rc=0。
跑它的副作用照旧：脚本自己往 `probes/161/r6/logs/flip-*.txt`（8 枚 tracked 件）写 log ⇒ 那 8 枚带 ` M`，**本程没提交、也没还原**。

### ③ 与 AC#7 同程落

同一程、同一次改 `want`／`want_n`／记账三处（本件 §3＋§4 就是那一程）；落地后**九腿基线与 14 枚腿数一并复算**＝§3 那张全表＋§4① 的 `名册=14 声明=14 记账=14`。

## 5. 我有没有减射程（单独一节，归非实现者专判）

**一枚未减、一枚未动。** 证据不是句子，是集合差（台件 `r2-roster-diff.sh`：把改前的尺（`git show f1b99a70:…gate-clauses.sh`）与工作副本的尺**打同一枚锚**，十四枚腿的名册两向 `comm`）：

```
r2-roster-diff.sh: anchor bb840a4f  pre-fix ruler (git show f1b99a70c3dfdad04f34094306b5942a06355f52) rc=0  working-copy ruler rc=0
  腿=G1     改前名册=0   改后名册=0   新增点名=0 退场=0
  腿=G1b    改前名册=0   改后名册=0   新增点名=0 退场=0
  腿=G2     改前名册=2   改后名册=2   新增点名=0 退场=0
  腿=G3     改前名册=0   改后名册=0   新增点名=0 退场=0
  腿=G4     改前名册=0   改后名册=0   新增点名=0 退场=0
  腿=G5     改前名册=2   改后名册=2   新增点名=0 退场=0
  腿=G5pos  改前名册=88  改后名册=88  新增点名=0 退场=0
  腿=G5neg  改前名册=16  改后名册=16  新增点名=0 退场=0
  腿=G6     改前名册=1   改后名册=1   新增点名=0 退场=0
  腿=G6pos  改前名册=6   改后名册=6   新增点名=0 退场=0
  腿=G6neg  改前名册=1   改后名册=1   新增点名=0 退场=0
  腿=G7     改前名册=3   改后名册=3   新增点名=0 退场=0
  腿=G7pos  改前名册=28  改后名册=28  新增点名=0 退场=0
  腿=G7neg  改前名册=4   改后名册=4   新增点名=0 退场=0
  合计名册行数：改前=151 改后=151   （两向 comm 全空＝射程一枚未动）
r2-roster-diff.sh: GREEN - 十四枚腿的名册两向差集全空，一条射程行未动。
```

（名册行数与腿的"枚数"不是一回事：run 腿这里取的是**逐行名册本体**，G5pos 那 88 行＝该腿 `git grep` 全部命中行，它的**未成对**枚数仍是 0。）
另外两笔：① 一条 pathspec 都没增删——`git diff f1b99a70 HEAD -- …gate-clauses.sh | grep -E '^[+-].*(:!|\$GO|\$GO2|pathspec)'` 只命中**三行注释**（本程写的说明文字），零枚射程行；② 腿没删——`腿数` 改前改后都 14（§3 文末那两行）。
⚠ 本程**自己犯过一次恒真比较并已在台件里改掉**：`r2-roster-diff.sh` 第一版给六枚对照腿取的段名前缀是 `echo "## G5 正控（…）"` 那一行，而 pair() 随后自己打 `## G5-正控 …` ⇒ awk 的段界立刻结束、六枚腿全报 `0/0`（＝什么都没比的绿）。改成取 pair() 自己那行标题之后才有上面这张表（G5pos=88／G5neg=16／G7pos=28 是真读数）。
⚠ **第一版那份空读数不在盘上**：两次运行 `tee` 进的是同一个 `logs/roster-diff.txt`，后一次把前一次**覆掉了**（本程没删任何东西，但也没留两版并存）⇒ 那一版只有本件这段自述＋台件注释里的记录，**没有可复算的原件**；要重看恒真那一形，把 `r2-roster-diff.sh` 的六枚段名前缀改回 `## Gx 正控`/`## Gx 负一负` 即可复现（验收方若要这条凭据，请自己跑，别按本件措辞去找文件）。

**没有用过任何排除**：本程没新增 `':!'`、没新增 `grep -v` 过滤、没把任何文件或腿排出射程。

## 6. 那枚 40→42 的现量（票面 `:53` 的「顺带一字」）

```
$ git grep -lEw 'Close\(\)' $(git rev-parse HEAD) -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l
42
$ git grep -lEw 'Close\(\)' f1b99a70 -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l
42
```

⇒ 尺里那句注释的 **40 改成 42**，两枚锚同数（与 `171-v1` 那两枚锚的读数一致）；命令一并写进了注释本体，下一程不必再猜它是怎么数出来的。

## 7. 契约轴（本程自量）

```
git -c core.quotePath=false show --name-only <本程四枚提交> | grep -E '^(internal|cmd|tools|frontend|design|docs/(PLAN|specs|reports))'  ⇒ 0 枚
git diff --name-only f1b99a70 HEAD -- tools/d22scan scripts | wc -l ⇒ 0
```

本程写面逐枚＝`probes/154/gate-clauses.sh`、`probes/161/r6/flip-declaration.sh`、`probes/171/r2/{r2-run-legs-and-subtraction.sh,r2-roster-diff.sh,logs/*}`、`docs/evidence/s1/171-run-legs-and-subtraction-r2.md`、票面 171 的 Progress log 一行。`internal/**`／`cmd/**` 里另一程 `162-r4` 的字节**与本程无关**（本程未提交、未还原它们）。

## 8. 门禁四数＋名册＋gofumpt（逐枚现跑）

| 门禁 | 命令 | 读数 |
|---|---|---|
| 尺本体 | `sh .scratch/wisp/probes/154/gate-clauses.sh` | 改前 rc=0（`腿数＝14 声明与实测不符＝0`）／**改后 rc=0**（同数＋`腿数断言＝相符`＋`基线过期枚数＝0`）⇒ 本程没有把自己的门改成恒红，也没有为变绿放宽任何断言 |
| 门禁聚合 | `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` | 改前 rc=0（九格）／**改后 rc=0（十格）**，restore＝baseline |
| D22 静态扫 | `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；ban #8 覆盖 `internal/` 422／`cmd/` 45 枚 Go 文件（`logs/d22scan.txt`） |
| 尺的自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**，**PASS=34 FAIL=0 SKIP=0**，`=== RUN=76`，`'[no tests to run]'=0`；名册 34 枚顶层测试名＝`logs/runtests-roster.txt`（复算命令：`grep -E '^=== RUN' logs/runtests.txt \| sed 's/^=== RUN  *//' \| grep -v '/' \| sort`）；改前名册未取，理由＝§2 第五栏 |
| gofumpt | `export PATH="$PATH:$(go env GOPATH)/bin"; gofumpt --version` | **`v0.12.0 (go1.27.1)`**（现跑，非引用）；本程 0 枚 `.go` ⇒ `gofumpt -l` 0 行 |

## 9. 提交的原样输出（每格一件；共享工作树里 HEAD 会被另一程前移，故逐枚点名本程自己的提交号）

原样贴在 `logs/commit-attest.txt`（同一节里还留着另一程 `fd8201e7` 夹在中间的现场，未改写）。摘要：

| 件 | 提交 | 内容 |
|---|---|---|
| ① | `4cee47f2` | step-0 基线读数 `logs/gate-pre-fix.txt` |
| ② | `7f78a163` | 台件 `r2-run-legs-and-subtraction.sh`（先落未修码读数） |
| ③ | `290aa63e` | 尺本体 AC#7＋AC#8①②＋全部改后 logs |
| ④ | `822e2001` | `flip-declaration.sh` 新增那一格＋`r2-roster-diff.sh`＋门禁 logs |
| ⑤ | 本件之后一枚 | 证据件＋票面 Progress log 一行＋`logs/commit-attest.txt` 收尾 |

## 10. 被拒／没成功的调用

- **权限系统拒绝：0 次**（没有一次工具调用被拒；没有"假设成功往下写"的段落）。
- 没成功但非权限的一次：一枚 `Edit` 的 `old_string` 不匹配（我把 `book()` 上方的注释行记错了），随后按文件里**逐字原文**重做成功。发生在**取数之后**（改前基线两份读数已落盘、未修码三发"今天不响"也已落盘），所以没有任何一份证据件里的读数依赖那次失败的编辑。

## 11. 删除命令次数＝**0**

本程未执行、也未创建任何删除命令（无 `rm`／`rmdir`／`Remove-Item`／`git clean`）。复算命令：`history` 不在本程写面里，但台件里刻意**没有** `rm`——两状态是**相加**的（state2 往同一棵合成树多落一枚样本文件，从不删 state1 那枚），`r2-run-legs-and-subtraction.sh` 与 `r2-roster-diff.sh` 全文可查：`grep -c 'rm -' .scratch/wisp/probes/171/r2/*.sh` ⇒ 两枚脚本各 **0**。临时件（`/tmp/wisp-171-r2.*`、`/tmp/wisp-171-r2-tree.*`）留在盘上未清，路径逐条打在两份 harness log 的末尾。

## 12. 伪授权两栏＋凭据

| 我引用过的"授权／同意"文字 | 出处（工具名＋命令前 40 字） |
|---|---|
| **0 栏**——本程没有引用过任何一条看起来像"批准／同意／可以推送／可以改契约"的文字 | 无 |

凭据值零抄录：本程未读过、也未抄录任何密钥／token／`config.toml` 内容（写面只有上述几枚文件）。

## 13. `next=`

1. **归编排者**：AC#7／AC#8 两格翻勾与缺口审计必须由**另一个** agent 做（`SPEC-12 §4.3` #1/#3）。派单里我引的每一枚"几枚"都带命令，专判建议＝①"我有没有减射程"那一节（§5，含我自犯的那次恒真比较）；② flip 新增那一格是不是真在核表而不是核打印（§4②）；③ run 腿的枚数口径＝**命中行数含注释行**（§2 第三栏）是不是可接受的判据口径——它与票面 AC#1 那枚尾随注释尺洞同源，本程**没有**顺手把它修掉（修它＝改射程）。
2. **减法侧还缺一发**：票面症状"腿**或它的样本**被摘"里的"删样本不删腿"那一发，本程只用 SHR 那张表侧面覆盖（§4②）。要正面钉，需要一枚"从合成树里少给一枚样本"的两向台件——本程未做（§2 末栏）。
3. **AC#1／AC#3 仍在票面未勾**（尾随注释被当调用点／归因尺自测接 CI），本程一字节未碰。
4. **台账可入账的三条过期文字**（本程未替别人改）：① `gate-clauses.sh` 里票 161 写下的行号组（`bridge.go:559/:633/:684` 等）在本程锚上**仍未复算**（票 171 r1 已在旁边并排过现量，没人核完）；② `flip-declaration.sh` 文件头 `:19`–`:20` 那句注释写的是 `the restore run (expect rc=0 with 6 of 14 legs ringing by design)`，而它自己的 GREEN 那行现量打印的是 **7 响／7 静**（`G2 G5 G5neg G6 G6neg G7 G7neg` 响）⇒ **那句注释是过期文字**（不是判据、也不是本程改坏的：本程只加一格、未动那九格的语义），归编排者入账；③ 票面 `:36` 那句 `:211/:642/:691` 已由 r1 记录为位移过。
5. 若下一程给 G5neg 那一腿加"排测试"，必须在同一格把 `want_n 8` 核下来（r1 在票面 Progress log `:75` 留的那条账），别留 `SHR` 当默认——本程之后那张『基线过期』表会让这件事**看得见**，但仍然**不替人做**。
