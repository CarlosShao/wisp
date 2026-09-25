# 142-r1 — 两读不原子的守卫（实现方自证；裁决表须出自非实现者）

工单：`.scratch/wisp/issues/142-the-two-git-spawns-are-not-atomic-so-a-tracked-directory-can-be-skipped-during-the-window.md`
地界：只动 `tools/d22scan/gitignore.go` · `tools/d22scan/main.go` · `tools/d22scan/scan_test.go` + 本文件。
本程**只 commit 不 push**；未勾票面任何一格；未写台账；未自派 `A##`。

## 0. 锚点、被验的东西、以及票面前提里我独立核过的那几枚

共享树在动（别的程在提交），所以每一枚读数旁边单独引它自己的锚点，不写裸数。
本文件成文期间 HEAD 依次经过 `0fe7629` → `1b98c7b` → `6e316a0` → `9d06544`（我第一枚）→ `21617e5` → `64b404d`（我第二枚）。

**票说什么。** `runGitIndex` 问 git 三次，其中两次是列目录（`ls-files -z` 与
`ls-files -z -i -c --exclude-standard`），两枚独立 spawn，中间落一枚 commit ⇒ 全量那一读缺那枚新变的
路径、也从没为它记过父目录 ⇒ `holds(dir)=false` ⇒ 一条 ignore 规则把**受追踪的目录**整支剪掉 ⇒
门少扫且不报错。出处＝`d22scan-close-r1-accept-r1.md` §5.3（锚点 `310816f`，非实现者验收程）。

**三条前提我各自量过（不是照抄票面）。**

1. **"两次 spawn"为真。** 改前 HEAD（`0fe7629`，与 `1b98c7b` 之间 `tools/d22scan` 零提交，
   `git log --oneline 0fe7629..1b98c7b -- tools/d22scan` 输出为空）里三枚 `runGitAt` 顺次发出，
   中间无任何原子性措施。读源码确认。
2. **"窄集溢出错集"在 git 那一层就真，且方向是少扫。** 仓外副本
   `D:\tmp\ticket-142\probe2\tree`（`git init` 自造、两枚 commit、git version 2.50.0 系，实量
   `git version 2.52.0.windows.1`）现量：
   - 读 1（commit 之前 `git ls-files -z`）= `.gitignore` / `frontend/src/panel.tsx` / `internal/ok/ok.go`；
   - 读 2（commit 之后 `git ls-files -z -i -c --exclude-standard`）= `frontend/dist/assets/bundle.js`；
   - `comm -13 full1 narrow2` ⇒ **`frontend/dist/assets/bundle.js`**（narrow ⊄ full，窗口确实打开）；
   - 同一棵树静止态（`comm -13 full2 narrow2`）⇒ **0 枚**（检查不误伤正常树）。
   文件那一支救不到：`holds(file)` 读 `ignored` 会答 true，但走不到那一层——父目录先被剪。
   ⇒ 与验收量到的 `holds(file)=true / holds(dir)=false` 同形、单向。
3. **"今天没有任何断言要求它"为真。** 改前 HEAD 上 `grep -n` 命中的只有注释与那句限定语，
   没有一处判定；`ab9d5f4` 那一批动的正是注释。⇒ 本票把限定语换成断言。

**真树底数（新守卫"不许误伤"的那一面）。** 锚点 `0fe7629` 的仓库根现量：
`git ls-files -z` = **1156** 枚、`git ls-files -z -i -c --exclude-standard` = **0** 枚、`comm -13` = **0** 枚
⇒ 本仓库此刻窄集是空集，子集关系平凡成立，新分支在干净签出上零次触发。
⚠ 这枚 0 不是"尺没看"：同一把尺的非空正控有两发——既有 `scan_test.go` 的
`TestFullTrackedListCoversWhatTheNarrowListCannot` 断言 `-i -c` 恰好等于 `frontend/dist/tracked-or-not.tsx`
一枚；本票 §1 块 1 断言它恰好命中 `frontend/dist/assets/shipped.tsx`。

**票面前提里我没复算的一条（写明）。** §5.3 那句"白盒造竞态量得出来"是验收程自述；我没找到它的副本，
也没重跑它那一发。我独立造的是上面 probe2 与 §1 的用例，同形同向，故本程结论按独立复现入账，
**不**按"复用了它的读数"入账。

## 1. AC#1 — 那一形能不能在测试里打开：**能**（本票存活）

### 1.1 用哪一面打开它（先说清哪一味是新的）

时间竞态那一支我**没有**走：窗口只有一枚 spawn 那么宽（十几毫秒量级），按它开测就是
"闲时绿、忙时红"的用例，本仓对这类东西的判法是先划掉（`d22scan-close-r1-accept-r1.md` §7.4
把"TMPDIR"那格判成结构上不可达就是先例）。走的是白盒那一支，而且只在**一处**换了装配口：

- `runGitIndex` 尾部原本内联的"两串 → 状态"装配，抽成 `buildGitIndexState(root, all, withIndex)`；
  抽的时候逐字搬，本票新增的判断只有那一段（见 §2）。这是本文件已有的风格：
  `verdict`/`fixtureVerdict` 早就是同一个体位的两扇门。
- 于是测试能做而生产不做的事只剩一件：**把两枚 spawn 之间插一枚真 commit**。
  两串字节都出自真 git（`runGitAt`），装配走生产的 `buildGitIndexState`，
  第 5 块的走查走生产的四条 walk（`scanWithIgnore`，见 §2.3）。

### 1.2 红的读数确实出现过（可重跑，命令在 §3.3）

副本目录 `D:\tmp\ticket-142\battery\M-A-no-guard`：`gitignore.go` 与交付版**只差被摘掉的那 6 行**，
`main.go`/`scan_test.go` 与交付版同 hash（`git hash-object` 现量：
`scan_test.go` = `e64b9631b600cd39a7abb5604df424babf7cd26e`、`main.go` = `a05c53b3c84cc525b524eceecde52868986e71ef`）。
守卫那 6 行摘掉后，同一发用例 `TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`
**红 14 条**，其中三句是要听的：

```
scan_test.go:1746: an index that moved between the two reads was ACCEPTED (holds("frontend/dist/assets", dir)=false, holds(file)=true): the walks prune a directory the index holds ...
scan_test.go:1818: ban #6 counted 3 frontend/ files on the torn index and 5 on a clean one, want +1 - the torn read must fail toward MORE scanning, never fewer
scan_test.go:1844: verdict on the torn index must be 1 (it scanned more, and more is what it found), got 0 err=
```

`got 0` 那一句就是票说的结局：**门是绿的，而它少读了两枚文件**（3 对 5）。
同一发用例在交付码上 `--- PASS (1.59s)`。⇒ 窗口打得开、用例打得响，AC#1 成立。
M-A 其余 11 条红行与逐枚归因在 §3。

### 1.3 顺带把"结构上打不开"的那一支划掉（免得下一位再试）

`rev-parse --show-prefix`（第三读）**不参与这条不变式**：它答的是"root 是不是仓库顶"，
只取决于 cwd 与 GIT_DIR，不取决于树里有几枚 commit，所以两读之间它翻不了。
它现有的一票否决（子目录 ⇒ 整棵树按"问不到"处理）原样保留，本票没动它。

## 2. AC#2 — 改动本体：只往"多扫"那一侧掉

### 2.1 牙（`tools/d22scan/gitignore.go`）

`buildGitIndexState` 在装配两串之前先问一句：**窄集里有没有全量集没有的路径？** 有 ⇒
"这两串描述的是两棵树"，返回 `ok=false` 并带一行原因；于是 `skip()` 落到它已有的那一支
（git 问不到 ⇒ **一条 ignore 规则都不应用**），`note()` 落回它已有的第一行机制。
新增判断本体 6 行：

```go
	if strays := straysOutsideTheFullRead(tracked, ignored); len(strays) > 0 {
		return gitIndexUnavailable("the git index changed between the two `git ls-files` reads of %s: %d path(s) "+
			"the second read reports as tracked are missing from the first read, so no single answer describes "+
			"this tree (a commit landed mid-scan) - strays: %s",
			filepath.ToSlash(root), len(strays), quotePathList(strays, maxStraysNamed))
	}
```

配套的三味小的（每一味都有"摘掉它谁红"的答案，见 §3）：
`straysOutsideTheFullRead`（差集，排序）、`quotePathList` + `maxStraysNamed = 8`
（点名前 8 枚、其余只报枚数；`strconv.Quote` 保证路径里带换行也不把那一行劈成两行）、
`gitIndexUnavailable`（把原来 `runGitIndex` 里的 `fail` 闭包提成包级函数，两处共用一个失败形状）。

**没做的（票面明令禁止，逐条对）**：没加重试、没加锁、没"比较后任选一枚读数"（M-C 那一发变异
就是它，测试专门为它红一遍）、没加豁免名单、没放宽任何 scope 定义、没把方向反过来、
没动任何既有断言、没加 `t.Skip`。也没动 `emojiRe` / 任何 scope 的文件清单 / `allowlist.txt` /
`thresholds.go` / golden / `docs/PLAN.md` / `docs/specs/**` / `internal/risk/**` / `rules_gateway.go`。

### 2.2 注释那一半：把"被承认"改成"被执行"，但**不**收回限定语

`ab9d5f4` bounded 的那两句仍是 bounded 话，本票没把它改回无条件话。写进去的是：
窗口仍然存在（还是两枚 spawn），变化在于**分歧现在会被查出来**；并且明写残余——
一枚落在窗口里的 commit 若只加了"git 自己不匹配任何规则"的路径，窄集仍是全量集的子集，
检查看不见它，所以 `holds()` 仍必须从**全量** `ls-files` 装载（那一支是 `A214/F3` 两形的守卫）。
顺手把原文里 `:240`/`:244` 两枚裸行号换成不会腐烂的指法（"runGitIndex 里那两枚 `ls-files`"）：
本票恰好就在移动那些行，留行号等于埋雷（同批 `TestFullTrackedList...` 的注释写明了这条规矩）。
文件头那段"已知会多扫的形"从两枚改成三枚，末句"不是穷尽声明"原样保留、没写最高级。

### 2.3 门（`tools/d22scan/main.go`）：`scanWithIgnore`

`scanWithStats` 拆成"建 matcher"与"带着 matcher 走查"两半，生产只走前一枚（它自己建新 matcher）；
后一枚多的那个参数只被本票的第 5 块用。理由与 `fixtureVerdict` 同一条：测试要走的必须是
**生产那四条 walk**，不是在测试里另抄一份走查。删的那一行只有 `ign: newGitIgnore(root),`
（改成 `ign: ign,`），逐枚点名在 §5。

### 2.4 顺带答票面那一问：子集关系今天有没有被执行？

改前：**没有**，任何一处都没有。今天执行它的是本票 §2.1 那 6 行。
改前两串都只被各自解析、从不互相比对，所以 `ab9d5f4` 那句限定语是纯措辞。

## 3. AC#3 — 变异自证，两向 + 三枚单点回退

全部在**仓外副本**做：`D:\tmp\ticket-142\battery\<mutant>`，脚本
`D:\tmp\ticket-142\battery.py`（逐枚复制四件文件、按唯一锚行改一刀、只跑这一发用例）。
仓库里那三枚文件在变异期间一字未动（跑完 `git status --porcelain -- tools/d22scan` 只列出我自己的
改动，`git diff --numstat` 与提交前一致）。每一枚都取 `--- FAIL/PASS` 与全部 `scan_test.go:` 行。

### 3.1 正向：造出 `narrow ⊄ full` 的种子 ⇒ 必须走全扫那一支并自陈

种子就是 §1 那一发用例的前 5 块（`shipped.tsx` 在两次读之间被 `git add -f` + commit）。
交付码上的读数（`--- PASS (1.59s)`，用例本体是 `scan_test.go` 里的
`TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`，按本文件自己的规矩不给裸行号——同批已经改坏过一次）：

| 拿什么量 | 交付码读到 | 说明 |
|---|---|---|
| 块 1（前提） | 窄集命中 `frontend/dist/assets/shipped.tsx`；读 1 既无此路径也无 `frontend/dist/assets` 目录 | 窗口真打开了，不是推的 |
| 块 2（守卫） | `ok=false`、`why` 同时含 `changed between the two` 与该路径 | 拒绝装配这两串 |
| 块 3（反"永远拒"） | `buildGitIndexState(root, 静止的全量, 同一窄集)` ⇒ `ok=true` 且 `holds(dir)=holds(file)=true` | 检查只在分歧时响 |
| 块 4（决策 + 已知正控） | `decide(dir).ignored=true`（规则活的），`skip(dir)=false`、`skip(file)=false`，`note()` 第一行 `d22scan: gitignore rules NOT APPLIED` | "什么都没跳"不是因为没东西可跳 |
| 块 5（端到端） | `examined[ban#6]`、`emojiSeen[frontend/]` 都比静止态 **+1**；两条 finding 都在；`verdict=1`；静止态那一遍 `note()` 里没有 `NOT APPLIED`、仍有 `skipped as git-ignored` | 全扫 + 响亮 + 正常树不受影响 |

### 3.2 反向：摘掉一味 ⇒ 谁变红（四枚）

| 变异 | 摘掉/改成的东西（生产侧，一行不差） | 读数 |
|---|---|---|
| **M-A** | `buildGitIndexState` 里 §2.1 那 6 行整体删除 | **FAIL，14 条红行**，含 `:1746` ACCEPTED、`:1818` `3 ... and 5 ... want +1`、`:1844` `got 0`。⇒ 守卫承重，且它一摘门立刻回到"绿的但少读两枚文件" |
| **M-B** | `main.go` 里 `ign: ign,` 改成 `ign: newGitIgnore(root),`（门还开着，但塞进去的状态被忽略） | **FAIL，3 条红行**：`:1818`、`:1822`（撕裂索引那遍读到 4 对 5）、`:1841`（只剩 `shipped.tsx` 一条 finding）。⇒ seam 不是装饰，第 5 块量的确实是那个状态 |
| **M-C** | 那 6 行换成"把窄集多出来的路径并进全量集"（= 票面禁的"任选一枚读数"那一支） | **FAIL，11 条红行**，`:1746` 那句变成 `holds(dir)=true`（事实是它自己造的），`:1818`/`:1822`/`:1841` 仍红（未追踪的构建产物又被跳掉了）。⇒ 测试分得开"拒绝"与"缝合"，不是只要多扫就算过 |
| **M-D** | `quotePathList` 的函数体换成 `return strings.Join(paths, ", ")`（同时去掉 8 枚上限与 `strconv.Quote` 转义） | **FAIL，3 条红行**：`:1879`（原因里出现真换行）、`:1882`（未点名的 stray 不再可见）、`:1890`（`(+1 more)` 没了）。⇒ 那两味小的也各自有人守着 |

M-A 的另外 12 条红行（同一发用例，摘守卫之后，14 条 = 上面引的 2 条 + 这 12 条；12 个不同站点）：
`:1752` 两枚（`why` 为空，一枚该含 `changed between the two`、一枚该含那枚 stray 路径）、
`:1785`（`skip(dir)` 真把目录剪了）、`:1792` + `:1796` 两枚（自陈首行不是那句机制，内容见下）、
`:1800`（多行）、`:1822`（ban #8 与 `:1818` 同因）、`:1841`（两条 finding 全空 `got []`）、
`:1876` + `:1882`（块 6 那枚换行种子无人拒）。

⚠ 一条值得单独交出来的读数：M-A 里 `note()` 长这样——

```
d22scan: skipped as git-ignored: 0 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (1 path(s))
d22scan: 1 path(s) git reports as TRACKED and matching an ignore rule (git ls-files -i -c --exclude-standard) - a tracked path is part of the delivery, so none of them was skipped: frontend/dist/assets/shipped.tsx
```

第二行是**假话**：`shipped.tsx` 确实"没被跳"，但它所在的目录被剪了，所以它一行也没被读。
改前那种状态下自陈会主动误导读者，这是本票比"措辞"走得远的一点。

### 3.3 复现这三枚/四枚读数

```
python D:\tmp\ticket-142\battery.py        # 四枚副本变异 + 各跑一发用例
# 或只看某枚的完整输出：
cat D:\tmp\ticket-142\battery\M-A-no-guard.log
```
交付码那一遍：`sh tools/d22scan/runtests.sh -C tools/d22scan ./...`（§4）。
脚本只读 `tools/d22scan` 并往 `D:\tmp\ticket-142` 里写；仓库目录内没有建过副本、worktree 或 checkout。

### 3.4 两枚我按定义标成"白盒/纵深防御"、没吹成抓到东西的

- 块 6（换行路径）与块 7（8 枚上限）：今天真实 git 输出里不会长出带换行的 tracked 路径，
  这两枚防的是"我把这份原因写进 note() 的第一行"这件事本身不破。摘掉它们红的是 M-D。
- 块 3/块 4 的反"永远拒"与"规则活的"两枚控制：它们不抓今天的任何伤，它们让块 2/块 4 的
  零读数有意义。按本仓"承重"的口径：摘掉它们，块 2 与块 4 就退化成"永远为真"，
  存在一种变异（永远拒绝装配）能打不红任何其它东西——所以它们承的是"这一格可证伪"的重。


## 4. AC#4 — 门禁与四枚数：改动前后同一把尺的读数

### 4.1 扫描器自家测试（`sh tools/d22scan/runtests.sh -C tools/d22scan ./...`）

两枚读数都取在 **`git archive HEAD | tar -x` 的纯净快照**里（同一形态、同一台机器），
因为工作树在动（前端会话正在往 `frontend/` 落未追踪文件）：

| 锚点 | PASS | FAIL | SKIP | === RUN | `^panic:` | 名册枚数 |
|---|---|---|---|---|---|---|
| `1b98c7b`（改前，我的第一枚 commit 之前） | 29 | 0 | 0 | 69 | 0 | 29 |
| `640cff1`（改后，含本票两枚代码 commit） | 30 | 0 | 0 | 70 | 0 | 30 |

**名册差集**（`grep -oE '^--- (PASS|FAIL|SKIP): NAME' | sort -u` 之后 `diff`）：

```
15a16
> TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules
```

⇒ **+1 枚、-0 枚**。四枚数之外还给名册，是因为一枚 panic 会吞掉同包其余读数：
两枚快照的 `^panic:` 都是 0，且没有任何既有名字从名册里消失。
工作树形态（非快照）在锚点 `21617e5` 上是同一组数：`PASS=30 FAIL=0 SKIP=0 === RUN=70`。
`go vet ./...`（在 `tools/d22scan` 目录内，它是自己的 module）干净，`gofmt -l .` 空。

### 4.2 门的读数（八枚数 + rc，每枚都带锚点）

| 取数形态 | 锚点 / 树 | bans #1-5 internal/ | cmd/ | ban #6 frontend/ | ban #7 | ban #8 design/ | frontend/ | internal/ | cmd/ | rc |
|---|---|---|---|---|---|---|---|---|---|---|
| 工作树（改前代码） | `0fe7629` | 203 | 22 | 43 | 18 | 32 | 43 | 407 | 39 | 0 |
| 纯净快照（改前代码） | `1b98c7b` | 203 | 22 | 43 | 18 | 30 | 43 | 407 | 39 | 0 |
| 纯净快照（改后代码） | `640cff1` | 203 | 22 | 43 | 18 | 30 | 43 | 407 | 39 | 0 |
| 工作树 A/B（同一枚树，两版码） | `d61281c` | 203 | 22 | 46 | 18 | 32 | 46 | 407 | 39 | 0 |

- **纯净快照那一行是本票 AC#4 指名要的形态**：`git archive` 出来的树没有 `.git`，
  改前改后都必须"响亮自陈"。两版都自陈，逐字同形：
  `d22scan: gitignore rules NOT APPLIED - git cannot be consulted in ...: git rev-parse --show-prefix: fatal: not a git repository ...`
  且它是扫描输出的**第一行**（A218(7) 那条性质没被本票动到）。八枚数两版一字不差。
  ⚠ 快照那行 `design/=30`（工作树是 32）：owner 手里那 16 枚未提交的 `design/` 挪动不进 archive，
  不是回归，与本票无关。
- **A/B 那一行是"正常树上的零行为改动"的正面证据**：同一枚工作树、同一条命令、
  跑两遍——一遍用 `1b98c7b` 的 `gitignore.go`+`main.go`（副本 `D:	mp	icket-142b-before-1b98c7b`），
  一遍用交付码——`diff` 八枚数与那行 provenance **完全为空**（`EIGHT_IDENTICAL`）。
  原因就是 §0 那枚底数：静止索引上子集关系成立，新分支零次触发。
- `frontend/` 今天一天里从 43 长到 46（`git ls-files frontend` 恒为 43，多的全是未追踪的工作树文件），
  是前端会话在落文件，本票的 A/B 已经把这条与码分开。⛔ 没有为变绿动任何断言、阈值或豁免。

## 5. 契约轴：断言删改枚数、Skip、以及"哪些生产行被删了"

### 5.1 `git diff --numstat`，逐枚点名（本票是生产行为改动，所以全列）

| commit | 文件 | + | - |
|---|---|---|---|
| `9d06544` | `tools/d22scan/gitignore.go` | 108 | 15 |
| `9d06544` | `tools/d22scan/main.go` | 13 | 1 |
| `9d06544` | `tools/d22scan/scan_test.go` | 214 | **0** |
| `64b404d` | `tools/d22scan/scan_test.go` | 37 | **0** |
| `640cff1` | `docs/evidence/s1/142-non-quiescent-index-guard-r1.md` | 197 | 0 |

三枚 commit 的全部路径（`git show --name-only --format="" 9d06544 64b404d 640cff1 | sort -u`）
只有上面这四枚文件，没碰 `allowlist.txt`、任何 `thresholds.go`、任何 golden、
`docs/PLAN.md`、`docs/specs/**`、`internal/**`（含 `internal/risk/**`、`internal/panel/**`）、
`frontend/**`、`design/**`、`rules_gateway.go`、`emojiRe` 那一行（`git show | grep -c emojiRe` = 0）。

**测试文件删除列 = 0 ⇒ 断言删改枚数 = 0。** 这是可从numstat 直接读出的强结论：
删除列为 0 意味着没有任何一行既有测试被改过或删过，我加的 251 行全部是插在既有行之间的新行。

**那 16 枚生产删除逐枚归因**（`git diff -U0 1b98c7b 64b404d -- tools/d22scan/{gitignore,main}.go | grep '^-'`）：

| 哪几枚 | 被删的原句 | 去哪了 |
|---|---|---|
| 1-2 | `// that list. Two shapes are known to do it, ... somewhere` / `// else:` | 措辞改成 "Three shapes ..."，同一句注释的续写，无语义 |
| 3-10 | `ab9d5f4` 那 8 行 bounded 话（含 `:240 and :244` 两枚裸行号） | 换成 §2.2 那段：限定语保留、"被承认"改"被执行"、补残余、裸行号换成不腐烂的指法 |
| 11-13 | `fail := func(format string, args ...any) *gitIndexState { return ... }` | 提成包级 `gitIndexUnavailable`，`fail := gitIndexUnavailable` 一行接上，`runGitIndex` 里 5 处 `fail(...)` 调用一字未改 |
| 14-15 | `return fail("cannot parse git ls-files ...")` 两枚 | 随装配一起搬进 `buildGitIndexState`，只把 `fail(` 换成 `gitIndexUnavailable(`，消息文本逐字保留 |
| 16 | `ign: newGitIgnore(root),`（main.go） | 搬到 `scanWithStats` 里那一行 `return scanWithIgnore(root, newGitIgnore(root))`，struct 字面量改成 `ign: ign` |

⇒ **没有一枚被删的行携带语义判断**：十枚是注释、六枚是搬家。删除列里没有任何一条断言、
任何一条 `return`/`fail` 的判定条件被去掉。

### 5.2 `t.Skip` 枚数前后，以及"这把尺看得见它"的正控

| 形态 | `grep -c 't\.Skip'` |
|---|---|
| `1b98c7b` 的 `scan_test.go` | 8 |
| 交付（`64b404d`）的 `scan_test.go` | 8 |

同 8 枚：7 枚既有 `t.Skipf(...)`（"not inside the wisp repo" 五枚 + 另两处）+ 1 枚注释里提到这个词；
本票新增 0 枚。⚠ 这枚 0 用同一把尺在**副本**里种了一枚真的验过：
`D:	mp	icket-142\skipcontrol\planted.go` 在 `TestScanCleanRepoIsGreen` 头里插一行
`t.Skip("planted known-positive for the grep ruler")`，同一条 grep 立刻报出
`113:	t.Skip("planted known-positive for the grep ruler")`（行号在此点名）。
交付文件里那 7 枚 `t.Skipf` 的行号（改后）：`241 / 261 / 677 / 1044 / 1189 / 2169 / 2242`
——前 5 枚与改前同行号，后 2 枚只是被我插在上面的新用例推下去（改前是 `1918 / 1991`）。
并且两枚快照的 `--- SKIP` 计数都是 **0**，即这些分支今天都不响。
