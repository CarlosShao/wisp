# 142-accept-r1 — 两读不原子守卫的对抗验收（非实现者程）

被验物：票 142 交付 —— `tools/d22scan/gitignore.go` · `main.go` · `scan_test.go`
实现方自证表：`docs/evidence/s1/142-non-quiescent-index-guard-r1.md`（§0–§8，六枚 commit 落档）
票面：`.scratch/wisp/issues/142-the-two-git-spawns-are-not-atomic-so-a-tracked-directory-can-be-skipped-during-the-window.md`（AC#1–AC#4，四格全未勾）
本件作者＝**非实现者**；本程**一字未改代码**，只取证与裁决；不勾票面、不写台账、不 push。

总裁与八格的落点在本件 §9／§10／§11。**每格现量现写，每格一枚 commit**，
所以下面每一节的读数都只对它那一节写明的锚点负责。

---

## 0. 锚点、脏树隔离、以及简报前提的逐条复算

### 0.1 锚点（自量，不抄简报）

```
$ git rev-parse HEAD
cd87354ffbe1f0529159f38872dfcdd70c7298cb          # 进场锚点，与简报给的 cd87354 一致
```

成文期间 HEAD 已被别的程推到 `b293784e9bced93290fb672722c5d74f8371fea3`，但：

```
$ git log --oneline cd87354..HEAD -- tools/d22scan
（空）
```

⇒ 被验的那枚 module 自我进场锚点之后**无人再动**，本件全部 `tools/d22scan` 读数对 `cd87354` 负责。

### 0.2 禁读脏工作树：逐枚证明

进场第一动作与每一发变异之后都跑同一条（**四枚路径分开列**）：

| 时刻 | 命令 | 读数 |
|---|---|---|
| 进场（12:2x） | `git status --porcelain -- tools/d22scan/` | 空 |
| 九枚变异跑完 | 同上 | 空 |
| 快照／A/B／容器跑完 | 同上 | 空 |
| 交件前 | 同上 | 空 |

工作树别处**确实脏**（`D design/assets/*`＋`D design/screens/*` 共 16 枚、`M internal/panel/l2_grant_boundary_test.go`、
`?? design/doubao/**`、`?? design/old/`）——那是另两程未提交的活，**不在本票地界，也一并不归本票的账**。
所以本程取数一律走 `git show <commit>:<path>`、`git archive <commit> \| tar -x`、以及**仓外副本**，
没有一次拿工作树当被验版本。

### 0.3 我的仓外工作区（只建不删，全在 `D:\tmp\accept142\`）

| 路径 | 是什么 | 现量 |
|---|---|---|
| `tools/d22scan/` | `git archive cd87354 tools/d22scan \| tar -x` 的 pristine 基座；三枚 .go 与工作树 `cmp` 逐字节相同 | `cmp` 三枚全 IDENTICAL |
| `mut-accept/{M-A-no-guard,M-B-seam-ignored,M-C-merge-strays,M-D-plain-join}` | 复算实现方那四发变异（**我自己的脚本** `battery-accept.py`，不复用它的 `D:\tmp\ticket-142`） | 每枚一份 `.log` |
| `mut-accept/{S-1-cap-raised,S-2-no-quote,S-2b-no-quote,S-3-no-sort,S-4-runGitIndex-bypasses-guard,S-5-reversed-subset-check}` | 本程自己加的单点回退／方向变异 | 同上 |
| `mut-accept/{AR-1,AR-3,AR-4,AR-5}` | 本程自己加的"控制块到底有没有牙"变异 | 同上 |
| `snap-head/` `snap-old/` | `git archive cd87354`／`1b98c7b` 全量纯净快照（无 `.git`），各跑 `runtests.sh`＋`scripts/d22scan.sh` | `gate.log` `scan.log` `roster.txt` |
| `snap-skip/` | `snap-head` 的拷贝，真植一枚 `t.Skip` 当尺子的正控 | `skip-gate.log` |
| `ab/{head-code,old-code}` | 仓外两版码，`go run . -root <仓库>` 打同一枚工作树（A/B 与稳定性） | `*.digest` `ab.diff` |
| `shim/` | 本程自造的 `git` PATH shim（`bin/git.exe`）＋一发**端到端探针**（§1.3），跑在 `cd87354` 与 S-4 两版码上 | `delivered.log` `s4.log` |

`git archive` 出来的快照里 `.git` 不存在（现量 `NO .git (pure snapshot)`），这是 AC#4 指名要的形态。

### 0.4 简报给我的前提，逐条复算（**有一条只对了半条**）

| 前提 | 我的量法 | 读数 | 判 |
|---|---|---|---|
| 锚点 `cd87354` | `git rev-parse HEAD` | 同值 | ✅ |
| `git diff --numstat 1b98c7b..cd87354 -- tools/d22scan/` ＝ `108/15`、`13/1`、`262/0` | 同命令现量 | `108 15 gitignore.go`／`13 1 main.go`／`262 0 scan_test.go` | ✅ |
| `^func Test` 29→30、名册差集恰 `+ TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`／`−0` | 源码级 `grep -E '^func Test'`＋`diff`；**再拿运行名册量一遍**（`--- PASS/FAIL/SKIP` 名字集合，两枚快照各一份） | 源码级 `15a16 > func TestIndexChanging…(t *testing.T) {`；运行名册 29→30、差集同一行、`−` 0 枚 | ✅（两把尺同值） |
| 它是独立 go module ⇒ 根目录 `go vet ./tools/d22scan/` 必失败 | 现跑 | `rootvet rc=1`，报文 `main module (github.com/CarlosShao/wisp) does not contain package …/tools/d22scan` | ✅ 是设计不是伤 |
| 无 `.git` 的纯净快照那一支必须响亮自陈（A218⑦） | `snap-head`＋`snap-old` 各跑 `scripts/d22scan.sh` | 两版都有 `d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …`，且**是 scan 那一步的第一行**（`awk` 取 `d22scan.sh: scan of` 之后第一行＝该句） | ✅ |
| `t.Skip` 自报 8→8，且"正控是靠副本里真植一枚量的" | 我自己植（§7.3） | 改前 8／交付 8／植后 9，且严格 runner `rc=1` | ✅ 尺子看得见 |
| 这条守卫在 linux 下"有没有分母" | 不推，**真跑容器**（§8.3） | `golang:1.27`（linux，git 2.47.3）上该用例 `--- PASS`、整包 `ok`；同一形态下 M-A `--- FAIL` 14 条红 | ✅ 有分母，且能红 |
| 实现方 §1.2 的"与交付版同 hash：`scan_test.go` = `e64b9631…`" | `git rev-parse <c>:tools/d22scan/scan_test.go` 逐枚 | `9d06544`=**e64b9631…** · `64b404d`=51b35ee7… · `1ed231e`/`cd87354`=3fa6c24b… | ⚠ **只对了半条**：那枚哈希是 `9d06544` 那一版的，不是最终交付版（差 48 行新用例）。同段的 `main.go = a05c53b3…` 到 `cd87354` 仍未变 ✅。⇒ 属读数记账精度，不推翻任何判据（我用 `cmp` 直接证明基座＝`cd87354`），但要按 §8 那条"引读数必须同引它取于哪一版"的自家规矩更正 |

### 0.5 本程不做计时测量（编队约束）

另有一程正在 `internal/panel` 跑 `go test`，所以本件**不引用任何秒数**、不做任何"快/慢"判据。
颜色、枚数、名册、字节、diff 才是本件的凭据。（`go test` 输出里带的 `(0.20s)` 一类我只当它是那行文本的一部分，
从不据它判任何东西。）

---

## 1. AC#1 —— 这一形能不能被钉成能响的用例：成立（本程另造了一发端到端凭据）

### 1.1 先答那一问：**这条守卫在生产里走不走得到？谁调它？**

`grep -rn 'runGitIndex\|buildGitIndexState\|scanWithIgnore\|scanWithStats\|newGitIgnore(' tools/d22scan/ --include=*.go`
现量，把调用者逐枚点名（生产＝`gitignore.go`/`main.go`，测试＝`scan_test.go`）：

| 被调 | 调用点 | 属性 |
|---|---|---|
| `buildGitIndexState`（定义 `gitignore.go:308`） | `gitignore.go:282`（`runGitIndex` 尾部） | **生产唯一入口** |
| | `scan_test.go:1744`／`1767`／`1885` | 测试（块 2／块 3／块 6） |
| `runGitIndex`（定义 `:261`） | `gitignore.go:510`（`gitIgnore.index()`） | **生产** |
| | `scan_test.go:1860` | 测试（块 5a） |
| `index()` | `gitignore.go:476`（`skip()` 里） | **生产** |
| `skip()` | `main.go:660/674`（walkGo）、`853/865`、`931/964`（walkText/walkEmoji 两支）、`1180/1185`（`checkRoot` 自家那枚 matcher） | **生产，八处** |
| `scanWithIgnore`（定义 `main.go:232`） | `main.go:221`（`scanWithStats`，塞一枚新 matcher） | **生产唯一入口** |
| | `scan_test.go:1808`／`1870` | 测试（块 5 撕裂态／块 5b） |
| `scanWithStats` | `main.go:210`（`Scan`）、`main.go:1331`（`main` 里 `checkRoot` 之后真扫描） | **生产** |
| | `scan_test.go:619/639/647/679/745`＋`scanFixture` | 既有测试 |

⇒ 链条今天是真的：`main` → `scanWithStats` → `scanWithIgnore` → 四条 walk → `skip()` → `index()` →
`runGitIndex` → **`buildGitIndexState`**（守卫就在那一枚的 `:317`）。**"测试构造撕裂索引"与"真实运行中撕裂"
不是同一件事**：前者只保证"装配器看见撕裂会拒"，后者还要保证"生产确实把那两枚 spawn 的结果交给这枚装配器"。
自证表 §6 第 1 项登记的正是这条缝（它写明"凭据是 `runGitIndex` 尾部那一行 `return buildGitIndexState(root, all, withIndex)`（人眼读过，未测）"）。
**本程不接受"人眼读过"当作一格的凭据，所以自己下去量了两发**（§1.2、§1.3）。

### 1.2 第一发：把装配口从生产里挪走，用例响不响？（**不响**，登记成立）

变异 `S-4-runGitIndex-bypasses-guard`：`gitignore.go:282` 那一行改成调用一枚我新写的
`assembleWithoutSubsetCheck`（＝把 `1b98c7b` 里内联的装配逐字搬回来，**没有**子集检查），
`buildGitIndexState` 本体一字未动。

```
S-4-runGitIndex-bypasses-guard   PASS     redlines=0        # Windows，go test -run 该用例
--- PASS: TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules (0.15s)   # golang:1.27 容器里同一枚变异
```

⇒ **自证表 §6 第 1 项那条残余是真的，我复现了**：今天没有任何断言要求守卫待在生产的调用链上。
（同一枚变异下实现方那 5 块全绿，包括块 5a —— 因为 5a 钉的是"静止态给回可信状态"，与守卫在不在无关，
它自己的 §8 也是这么写的。）

### 1.3 第二发：**端到端**，真 spawn、真 commit 插在两读之间、生产唯一那扇门（**响**）

自证表 §6 第 1 项说"没造 PATH-shim 端到端那腿"。本程把那腿造出来了，造在**仓外探针**里
（`D:\tmp\accept142\shim\`，不进仓）：一枚 `git` shim 挂在 PATH 上，所有调用原样转发真 git，
只有一件额外的事——**看见第二枚 `ls-files -z -i -c --exclude-standard` 时，先落地一枚真 commit 再回答**。
被测包一字未改（用的就是 `cd87354` 那三枚文件），入口是生产唯一的 `scanWithStats`。

```
== delivered ==   --- PASS: TestProbeRealSpawnsTearThroughProduction
```

探针里被断言的读数（交付码）：`s.ign.idx.ok == false`；`note()` 第一行
`d22scan: gitignore rules NOT APPLIED … changed between the two`；
`examined["panel-approval"] == 5`（撕裂那一遍把 `frontend/dist/assets/` 里两枚 `.tsx` 都读了）；
第二遍（同一棵树、同一个二进制、marker 已耗尽 ⇒ 两读一致）`ok==true`、`examined==4`、自陈里没有 `NOT APPLIED`。
⇒ **守卫在生产里走得到，而且它今天就在响**：这条不是推的，是真 spawn 之间落了一枚真 commit 量出来的。

同一枚探针打 §1.2 那枚变异：

```
== S-4 ==         --- FAIL: TestProbeRealSpawnsTearThroughProduction
    zz_shim_probe_test.go:64: production runGitIndex ACCEPTED a pair whose second read names frontend/dist/assets/shipped.tsx and whose first read does not
    zz_shim_probe_test.go:74: torn scan examined[panel-approval]=3, want 5 … findings=[]
```

⇒ 那腿**造得出来、也不贵**（shim 现量 77 行＋探针 102 行 `_test.go`，全在仓外、不进仓），
并且它抓得到 §1.2 那枚变异。

### 1.4 判：**AC#1 成立**；残余＝登记即可，但要带一条入账口径

三条理由摆在一起看：

1. **票面自己点的就是白盒那一手**（票 §AC#1："白盒构造两次读之间索引发生变化的状态……本仓已有同风格先例：
   直接问 `holds()`"）。那枚先例 `TestFullTrackedListCoversWhatTheNarrowListCannot`（现量 `scan_test.go:1570`
   起）做的是"用 `gitCommand` 直接问 git 两遍、再问 `ix.holds(...)`"——它**不动装配器**；
   更早那手"手搓 `gitIndexState{ok:true, tracked:{}, dirs:{}, ignored:{file:true}}` 塞进 matcher"是**验收程自己在仓外造的探针**
   （`d22scan-close-r1-accept-r1.md` §5.3 的 `fuzz\zz_race_probe_test.go`）。
   交付那发用例喂进装配器的两串字节**都出自真 git、跨一枚真 commit**，比这两手都强一档。
2. **放弃计时竞态那条腿的理由站得住**，而且它没把判据换成能通过的那版：原句不抹，
   按 `d22scan-close-r1-accept-r1.md` §7.4 的先例写明为什么不给读数（窗口≈一枚 spawn ⇒ flaky 守卫比没有更糟）。
   我把那条先例本身读了一遍（它 §7.4 是把"TMPDIR 落点"那格判成**结构上打不开**结案），
   形状同类：**用它当先例不越界**。
3. §1.3 已经把"生产到不到得达"这问**用真 spawn 答掉了**，所以这条残余不再是"不知道生产响不响"，
   只剩"防不住将来有人把装配口搬走"——而那是**回归防护**问题，不是本票 AC#1 的成立条件。

⇒ **入账口径建议（编排者写，本程不写台账）**：AC#1 记"成立（白盒钉住＋验收程端到端复核过一次）"；
同时把"**没有任何断言钉住 `runGitIndex → buildGitIndexState` 那一跳**"作为**开放残余**登记，
归口见 §9 的 F-142-1（建议下一枚 `tools/d22scan` 票顺手收：shim 腿或一条源码形状断言，二选一，别为它单开票）。

⚠ 顺带一句诚实话：§1.3 那腿在我这里是**取证**，不是**交付**——它跑在仓外副本，仓里没有它，
所以它对"今天的生产"有效、对"明天被人改坏"无效。这一点写进 §10 的"没测"里，别让读的人以为门已经钉上了。

---

## 2. 承重复算：那四发变异我自己跑了一遍（**不照收它的读数**）

### 2.1 我的台件与它的台件是两套

脚本 `D:\tmp\accept142\battery-accept.py`（**不是**它的 `D:\tmp\ticket-142\battery.py`；那目录我一次都没读），
基座是 `git archive cd87354 tools/d22scan` 的 pristine 拷贝，每枚变异**只改一枚文件的一刀**、
按唯一锚串定位（锚串出现次数≠1 就 `sys.exit` 拒绝跑），只跑 `TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`。
红行计数口径：`grep -cE "^    scan_test.go:[0-9]+:"`，并且先排除本仓那个坑——`t.Logf` 也带 `file:line:` 前缀：
现量该用例体内 `t.Logf` 命中 **0** 枚，所以"红行数＝断言数"这层等式在这一发上是成立的
（M-A 站点分布：1746×1／1752×2／1785／1792／1796×2／1800／1818／1822／1841／1844／1887／1893＝14）。

| 变异 | 我摘/改的东西（生产侧） | 我的读数 | 它的读数（取自它的表 §3.2／§8） | 判 |
|---|---|---|---|---|
| **M-A** | `buildGitIndexState` 里守卫那 6 行整体删除 | **FAIL，14 条红**，站点与 §8 更正后的名册**逐枚同值**（含 `:1887`/`:1893`）；关键三句：`ACCEPTED (holds(dir)=false, holds(file)=true)`、`counted 3 … and 5 … want +1`、`verdict … got 0` | 14（同站点） | ✅ 复算一致 |
| **M-B** | `main.go` 里 `ign: ign,` → `ign: newGitIgnore(root),` | **FAIL，3 条红**：`:1818`、`:1822`（撕裂那遍读到 4 对 5）、`:1841`（只剩 `shipped.tsx` 一条 finding） | 3（同站点） | ✅ |
| **M-C** | 守卫换成"把窄集多出来的路径并进全量集＋补 `dirs`"（＝票面禁的"任选一枚读数"那一支） | **FAIL，11 条红**；`:1746` 那句真变成 `holds(dir)=true`（它自己造的假事实），`:1818`/`:1822`/`:1841` 仍红 | 11（同） | ✅ |
| **M-D** | `quotePathList` 函数体换成 `strings.Join(paths, ", ")` | **FAIL，3 条红**：`:1890`（原因里出现真换行）、`:1893`、`:1901`（`(+1 more)` 没了） | 3（同） | ✅ |
| **linux 半边** | 同一枚 M-A、同一份源码，在 `golang:1.27`（linux，git 2.47.3）容器里再跑 | `--- FAIL`，`linux-M-A-redlines=14` | 〔它没跑过，自报 §6 第 2 项〕 | ✅ 红在 CI 用的那一侧也成立 |

⇒ **四发的"摘掉它就没人红"我独立复现了，枚数与站点全同**；它给的因果链（守卫承重／seam 承重／"缝合"那一支被专测拒掉／两味小的各有红）没有一处需要我打折。

### 2.2 按本仓"承重"的可操作定义逐味判

定义（简报与 `issues/README` 同口径）：**摘掉任意一味，都存在一发变异从此打不红 ⇒ 那味承重。**
上一表就是这条定义的直接执行：M-A/M-B/M-C/M-D 各有一发变异、摘掉对应那一味就**不再红**。
本程再加两味它没列进这张表的（§3 详解）：`quotePathList` 与 `strconv.Quote` 拆成两枚独立回退后**各自仍红**，
而 `straysOutsideTheFullRead` 里那枚 `sort.Strings` **摘掉无人红** ⇒ 只有那一味是装饰。

### 2.3 另问一句它可能给出误导答案的：**摘掉它，外部可见读数变过没有？**

这是这类守卫最容易被读错的地方，两半都要给：

- **静止树上：一字没变，而且这是设计而非缺陷。** 我的 A/B（仓外两版码、同一条 `go run . -root <仓库>`、同一枚工作树，
  三遍跑）：`A_B_IDENTICAL`（八枚 scope 行＋两行自陈＋`clean` 那行全同）、
  `HEAD_RUNS_IDENTICAL`（同一版码跑两遍也全同 ⇒ 树没在两遍之间动）。
  门禁四数在**纯净快照**形态：`1b98c7b`＝`PASS=29 FAIL=0 SKIP=0 === RUN=69`、
  `cd87354`＝`PASS=30 FAIL=0 SKIP=0 === RUN=70`，`^panic:` 两版均 0；八枚数只差 `frontend/ 43→46`，
  归因是别人那一枚 `d61281c`（§7.4）。
  ⇒ **谁要是拿"基线四数没变"来判这味不承重，他就读错了**：一条只在两读分歧时响的守卫，
  在静止树上的**正确**读数就是零变化。
- **撕裂树上：外部可见差别很大，我量到三处。** §1.3 的端到端探针在交付码上给
  `note()` 第一行 `NOT APPLIED … changed between the two`、`examined[panel-approval]=5`、两条 finding；
  同一探针在 M-A／S-4 上给 `examined=3`、`findings=[]`、第一行换成一句
  `skipped as git-ignored … none of them was skipped: …`（§5 那枚假话）、而 M-A 里 verdict 从 1 掉到 0。
  ⇒ **门"绿着少扫文件"这件事在本票之后变成"红着多扫文件"**，方向确实是票要的。

---

## 3. 单点回退进攻：它超出最小修法的那三味，各自"撤掉哪条用例变不响"

它的 §3.4 把这三味标成"纵深防御、没吹成抓到东西"。我按简报要的口径**只撤其中一处**、逐味问那一句：

| 撤掉的这一处（只此一处，其余一字不动） | 我的变异名 | 哪条用例变得不响（现量站点） | 判 |
|---|---|---|---|
| `quotePathList` 整枚（换成 `strings.Join(paths, ", ")`，即上限与转义一起没） | `M-D-plain-join` | **3 条**：`:1890` 原因里长出真换行、`:1893` 未点名的 stray 不再可见、`:1901` `(+1 more)` 消失 | 承重（但红只出现在块 6/块 7 两枚合成串上，见下面"它的诚实处"） |
| `maxStraysNamed = 8` 的上限（抬到 64，函数体一字不动） | `S-1-cap-raised` | **1 条**：`:1901`，报文直接跟着变成 `must name the first 64` | 承重，但**只有一枚**、且那一枚喂的是合成列表 `p1..p9` |
| `strconv.Quote` 的转义（上限留着，只把逐枚转义拿掉） | `S-2b-no-quote` | **3 条**：`:1890`、`:1893`、`:1901`（`p1` 不再带引号） | 承重 |
| ⚠ 第一次尝试的同名变异 | `S-2-no-quote` | **NORESULT**：`.\gitignore.go:139:2: "strconv" imported and not used` ⇒ **编译都没过，那枚读数不作数** | 台账式提醒：这类"摘一味"的变异如果留下未用 import／未用变量，得到的是 NORESULT 不是红，必须换形（`S-2b` 用 `_ = strconv.Quote` 把 import 钉住，才真正隔离出"转义"这一味） |

**三味都答得出"撤掉谁红"⇒ 按定义不是装饰，本程不退回。** 但两句实话要留下：

1. **红的都是块 6/块 7 那两枚合成断言**（`buildGitIndexState(root, "", "we\nird.tsx")` 与
   `quotePathList([]string{"p1",…,"p9"}, 8)`），不是真 git 产物。它自己在 §3.4 第一句就是这么写的
   （"今天真实 git 输出里不会长出带换行的 tracked 路径，这两枚防的是'我把这份原因写进 note() 的第一行'
   这件事本身不破"），§6 第 4 项又登记了"9 枚只是合成字符串，不是真 git 输出"。
   ⇒ **措辞诚实，未越界**，我不按"它吹了"退它。
2. `maxStraysNamed` 抬到 64 仍然只红 1 条，说明这味的全部重量就挂在一枚 `HasSuffix("(+1 more)")` 上；
   它防的是"自陈第一行被撑成文件清单"，而**没有任何断言钉住"第一行不许超过多长"**。这不算错（票没要求），
   但记账时要按"一枚弱断言"记，别按"两向都响的守卫"记。

### 3.1 本程额外撤的两处（它没列进自己那张表）

| 撤掉的东西 | 变异名 | 读数 | 判 |
|---|---|---|---|
| `straysOutsideTheFullRead` 里那枚 `sort.Strings(strays)`（差集排序） | `S-3-no-sort` | **PASS，0 条红** | ⚠ **这一味才是真装饰**：今天没有任何断言要求"被点名的 strays 是有序的"（`note()` 第三行的 `ignoreSet` 排序是另一处，被既有 `:1525` 那类断言看着）。窄集本身出自 git、今天已是字典序 ⇒ 摘掉排序不可见。**不必退回**（它没声称过这味有牙），但自证表 §2.1 把它列进"配套的小的"时应补一句"排序无断言" |
| `gitIndexUnavailable`（把原 `runGitIndex` 里的 `fail` 闭包提成包级函数） | 未造变异 | —— | 它是**搬家不是新造**：`1b98c7b` 里那个闭包返回的是同一个 `&gitIndexState{why: …}`。它也没声称这枚有自己的 red mutant（§5.1 把 11-13 那三枚删除归因为"提成包级函数"），⇒ 按搬家入账，不要求牙 |

### 3.2 撤干净核对口令（它 §7 末段给的那句）复核

它写："撤回这批＝把 `gitignore.go` 里那 6 行判断删掉，撤完 §3 的 M-A 那发变异就会从红变绿"。
⇒ 我的 M-A 就是"只删那 6 行"那一枚，它红 14 条。这句话**方向讲反了一半**：删掉那 6 行之所以**仍然红**，
是因为用例本身要的是"守卫存在"；"撤干净"的判据不该是"用例变绿"，而应是"用例连同 seam 一起撤"。
本程量到的更准确的撤净口令是：**撤那 6 行 ⇒ 用例红 14 条（M-A）；只撤 seam ⇒ 红 3 条（M-B）；
两处都撤干净（回到 `1b98c7b` 的码）⇒ `scan_test.go` 里 `buildGitIndexState` 三处引用直接编译不过**。
它给的那句留作读者口令够用，但**别把它当判据抄进下一张票**。

---

## 4. AC#2 —— "只许一个方向"：实现对了，但"别的形状谁挡"要分四支写清

它实现的是：`narrow ⊄ full` ⇒ `gitIndexUnavailable(...)` ⇒ 一条 ignore 规则都不应用、全扫、走 `note()` 第一行。
这一支**票面钦定的四禁全部没碰**（没重试、没锁、没"比较后任选一枚读数"、没豁免名单、没把方向反过来）：
`buildGitIndexState` 里那 6 行只做一件事——返回 `ok=false`；`git diff 1b98c7b..cd87354 -- tools/d22scan/gitignore.go`
里没有任何 `for` 重试、没有 `sync.`、没有新表。它举的"问到了但两个答案对不上＝问不到"这条类比
（`ca84b75` 的既定失败方向）我对着 `note()` 的现文核过：撕裂态与"没有 .git"走的是**同一枚**首行机制
（`gitignore.go:531-534` 那一支），§1.3 端到端探针量到的第一行与快照里那一行同前缀，只差 `why` 内容。

### 4.1 反方向（全量里有、窄查里没有）：今天有人挡，且挡的是"不许响"

窄查本来就是全量上的一个过滤器，所以这一方向**绝大多数时候都在发生**（每一枚未被规则盖住的 tracked 文件都是）。
挡它不许响的是三处断言：块 3（`:1768` 静止对必须 `ok=true`，`t.Fatalf`）、块 5 的前提（`:1812` 静止扫描必须 `examined==4`）、
块 5a（`:1860` 生产那扇门在静止树必须给回可信状态）。**我把反方向真造出来量了一遍**：
`S-5-reversed-subset-check`（把 `straysOutsideTheFullRead` 的差集方向掉个头＝"全量里窄查没有的"也算 stray）⇒

```
S-5 reversed check: FAIL redlines=2
    scan_test.go:1752: the refusal must name "frontend/dist/assets/shipped.tsx", got "the git index changed between the two …"
    scan_test.go:1769: a quiescent pair was refused (…) - the subset check fires on more than the race and every real scan would lose its ignore filter
```

⇒ 反方向**有人挡**，且 `:1769` 那一枚正是为它立的牌位（它在静止树上就红，不用等到撕裂树）。

### 4.2 第三种形状（本票没碰、票面也没要求）：三枚我自己造出来，逐支说"该不该挡"

| 形状 | 我的变异 | 读数 | 该不该挡 |
|---|---|---|---|
| 两枚 spawn **先后顺序**被换（窄读在前、全量在后 ⇒ 窄是旧树、全量是新树） | `S-6a-spawn-order-swapped`（整包跑） | `PASS=24 FAILset=[]`，与基座**一模一样**（基座同形态 `PASS=24 FAILset=[] SKIP=6`）⇒ **今天没人挡** | **不该为它开 AC**：这一支伤不到少扫那一侧——`dirs` 与 `tracked` 由**更新的**那枚全量喂，能持住的路径只多不少 ⇒ 方向仍是多扫。要写的是注释，不是断言：`gitignore.go:119-129` 那第三枚 bullet 现在只登记了"窗口里那枚 commit 只加了不匹配任何规则的路径"这一族残余，"两读顺序反过来也是同一族（子集不破 ⇒ 检查看不见）"没写。⇒ 记进 §9 的 F-142-3（措辞级、不阻塞本票） |
| 全读失败被吞（`all=""` 继续装配） | `S-8-full-failure-swallowed` | `PASS=24 FAILset=[]` ⇒ **没人挡** | 伤害面：全量为空 ⇒ 窄集里任何一枚都成 stray ⇒ **本票这枚守卫反而把它兜住了**（除非两串都空）。所以这一支现在比改前更安全，只是"更安全"这件事没人测。该不该挡：**该**，但注入面只有 shim（见 §1.3），与 F-142-1 是同一枚小腿，别单开票 |
| 窄读失败被吞（`withIndex=""` 继续装配） | `S-7-narrow-failure-swallowed` | `PASS=24 FAILset=[]` ⇒ **没人挡** | 伤害面**不是少扫**：`holds()` 的文件位读 `tracked \|\| ignored`、目录位只由全量喂，所以吞掉窄读只让 `note()` 少掉第三行（那枚"哪些受追踪件坐在忽略规则下"的出处），覆盖面一点没丢。⇒ 与本票无关的**先存在**缺口：那句"three read-only questions are all-or-nothing, a half-loaded index would be a new blindfold"在 `1b98c7b:176` 就写着，本票一字未碰（`git diff … \| grep -c all-or-nothing` = **0**）。记进 §9 的 F-142-2，别算在本票账上 |

### 4.3 判

**AC#2 成立**：方向只有一处、落点是既有的 `ok=false` 机制、四禁未碰、反方向有专断言挡着且我量到它红。
"别的形状"我按四支分开写了——**两支该补的是注释和一枚小注入腿，一支与本票无关**，
没有一支值得我为它硬开一道今天不响的 AC（票面 §AC#1 的 ⚠ 就是禁这件事）。

---

## 5. 它自报的那条副作用（"自陈会说谎"）：**独立复现成立，是一枚先于本票的缺陷**

### 5.1 我自己复现的两发（不复用它的读数）

**第一发（白盒）**：`mut-accept/M-A-no-guard.log`，摘掉守卫那 6 行、其余一字不动，用例 `:1792` 把当场的
`note()` 原文打出来：

```
d22scan: skipped as git-ignored: 0 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (1 path(s))
d22scan: 1 path(s) git reports as TRACKED and matching an ignore rule (git ls-files -i -c --exclude-standard) - a tracked path is part of the delivery, so none of them was skipped: frontend/dist/assets/shipped.tsx
```

**第二发（端到端，更硬）**：§1.3 那枚 PATH-shim 探针打在 `S-4`（＝生产装配回到 `1b98c7b` 的内联形状）上，
`note()` **同一句假话**照样出现，而且这一次我手里有独立证据说明它是假的：
`examined["panel-approval"] = 3`、`findings = []` —— 也就是说 `shipped.tsx` **一个字节都没被读过**，
它的父目录 `frontend/dist/assets/` 在第一行里刚被点名剪掉。
⇒ 与它 §3.2 那段引文逐字同形，**这一发我认**。

### 5.2 归因：先于本票，且本票没有把它改回无条件话

```
$ git show 1b98c7b:tools/d22scan/gitignore.go | grep -n "none of them was skipped"
459:			"- a tracked path is part of the delivery, so none of them was skipped: %s",
$ git diff 1b98c7b..cd87354 -- tools/d22scan/gitignore.go | grep -c "none of them was skipped"
0
```

⇒ 那句话说的是"**这一枚清单里**的路径没被 skip() 剔掉"，但在撕裂索引上它被读成"这些文件被看过了"，
而目录已经被剪。它是 `A214/F1` 那一批（`ca84b75` 系）立第三行自陈时带来的，**不是本票造的**，
本票也没动那一行一个字。

### 5.3 本票到底把它怎么样了：**遮住了一半，但遮住的方式是一条不变式，不是一条断言**

我把"守卫开着 ⇒ 第三行不可能是假话"这条推导走完（并且只用文件里的事实）：
`pruned` 只在 `skip()` 里、且只有 `holds(rel,dir)==false` 时才写入（`gitignore.go:483-500`）⇒ 一枚目录被剪，
意味着它不是任何 tracked 路径的祖先；而守卫保证 `ignoreSet ⊆ tracked`（`:317` 那 6 行，不然直接 `ok=false`）⇒
**第三行点名的路径不可能坐在被剪掉的目录下**。守卫关着的两行第一行（`ok=false`）不输出剪枝与第三行的
撕裂组合同时无效。⇒ 结论：这一支**今天结构上不可达**。

但要说清三件事，免得下一位把它当"已经钉住"：

1. 全仓**没有任何断言**写"第三行点名的路径必须落在未被剪枝的目录里"。`ignoreSet` 的引用只有三处
   （`gitignore.go:237/328/550/553`），测试侧命中的是 `scan_test.go:1525` 那句"note 必须点名这一枚路径"——
   它证明这行**会说话**，不证明说出来的**是真话**。
2. 遮住的只有"窄集冒到全量之外"这一支。§4.2 那两支（两读顺序掉包、吞掉某枚读失败）**不在这枚守卫的射程里**；
   我逐支推过：吞掉窄读只让第三行消失（不会说谎），两读顺序掉包时全量是新的、`dirs` 反而更全 ⇒ 也造不出这枚假话。
   ⇒ 目前**没有**已知的活着的通路；这句是推导，不是量到的零（我没能造出任何一枚今天还能让它响的变异）。
3. 它 §3.2 把这条主动交出来了（"改前那种状态下自陈会主动误导读者，这是本票比'措辞'走得远的一点"），
   **没有把它吹成"我修了自陈诚实性"**；措辞与它实际做的相符。

### 5.4 建议归口（**本程不动手修**）

- 不新开票。这一条与 §9 的 **F-142-1**（"没有任何断言钉住 `runGitIndex → buildGitIndexState` 那一跳"）
  是同一枚小腿的活：那腿（shim 注入或源码形状断言）一旦补上，顺手加**一行**断言
  "第三行点名的每一枚路径，都出现在本次 `examined`/非 `pruned` 目录之下"就同时把这条也钉住。
- 台账口径建议写成：**先于票 142 存在的自陈诚实性缺陷（`1b98c7b:459` 那一句），在"窄集溢到全量之外"这一支上
  被票 142 的守卫顺带堵住（结构推导＋验收程两发复现），但无断言持有**；
  归口 F-142-1 那枚小腿，不单开、不阻塞本票验收。

---

## 6. 恒真判据自查：本票每一条"必须红/必须绿"，先在别处试响一次

口径：一条判据如果**在任何变异下都响不了**，它就是装饰，不是被验物的功劳。
我在 §2–§5 之外又自己造了 8 枚变异（合计 **18 枚**，全部仓外副本，`mut-accept/` 与 `shim/`），
逐枚点名"谁把它打红"。表里"唯一命中"＝摘掉这条断言后**再没有任何断言**能抓住同一枚变异
（我用"把那枚断言删掉再跑同一枚变异"验的，见 AR-3／AR-5）。

| 判据（交付物里的那一句） | 我在未修旧码／反例上试响的结果 | 判 |
|---|---|---|
| AC#3(ii)「摘掉守卫 ⇒ 用例从绿变红」 | 真旧码（`1b98c7b`）上这枚用例**根本编译不过**（`buildGitIndexState` 不存在）⇒ 拿它当"旧码试响"是**平凡真**，不作数；行为等价旧码的那一发是我的 **M-A** ⇒ **红 14 条**（Windows 与 linux 各一遍） | ✅ 有牙，且判据选得对（它没拿"旧码编不过"糊弄这一格） |
| AC#3(i)「造出 `narrow ⊄ full` 的种子 ⇒ 走全扫且自陈行出现」 | 同一枚种子在 M-A 下走的是**剪枝**那一支（`examined 3 对 5`、`verdict got 0`）⇒ 判据可反 | ✅ |
| 块 1 前提第一句（`:1733` 窄读必须点名 stray） | `S-9b2-no-rule-in-the-seed`（把两枚 `.gitignore` 种子换成不盖住 `dist/` 的规则）⇒ **红 1 条，正是 `:1733`**（报文 `got []`） | ✅ 有牙 |
| 块 1 前提第二句（`:1737` 读 1 不许已经含新路径） | `S-9c-read1-after-commit`（把 read 1 的取数挪到 commit 之后）⇒ **红 1 条，正是 `:1737`** | ✅ 有牙 |
| 块 1 的另一发（撤掉 `gitForceAdd` 让种子失效） | `S-9`：先在 `:1707` 的 `gitCommitAll` 上炸（空索引 ⇒ `git commit` 失败）⇒"种子失效就 fail OUT"这句成立，但**站点不是块 1**，这一枚只算旁证 | ⚠ 记账精度：块 1 的两枚前提**不是**唯一防线，台件更早响 |
| 块 3 反"永远拒"（`:1769`） | `AR-1-always-refuse` ⇒ 唯一命中 `:1769`；**但**把块 3 整段删掉再跑同一枚变异（`AR-3`）⇒ 块 5 的静止前提在 `:1799` 照样红 | ✅ 有牙，❌ **不唯一** ⇒ 它 §3.4 那句"存在一种变异（永远拒绝装配）能打不红任何其它东西"**实测说重了** |
| 块 4 反"规则是死的"（`:1782`） | `AR-4-rules-dead` ⇒ 唯一命中 `:1782`；删掉这枚断言后同形变异（`AR-5`）在 `:1810` 仍红 | 同上：有牙、不唯一 |
| 块 4 的 `skip(dir)=false`（`:1785`）与三句 `note()`（`:1792/1796/1800`） | M-A ⇒ 全部命中 | ✅ |
| 块 5 的 delta（`:1818/1822`）、finding 集合（`:1835/1838/1841`）、verdict（`:1844`） | M-A／M-B／M-C 三枚各自命中不同子集；`S-10` 命中它的静止前提 `:1813` | ✅ |
| **块 5a（`:1861/1863`，"生产那扇门静止态必须可信"）** | **18 枚变异里没有一枚命中它**。最接近的是 `S-10-runGitIndex-injected-refusal`（让 `runGitIndex` 无条件拒），但它在 `:1813`（块 5 的静止前提）就 `t.Fatalf` 中止，**永远走不到 5a** | ⚠ **今天不可能响的那一枚**。它 §8 已自报"5a 在四枚变异下都不红"，措辞诚实；但按票面 ⚠"不许硬开一道永不响的"那把尺，5a 属**正向对照**、不是抓东西的腿 ⇒ 台账里不许写成"5a 钉住了生产那一跳"（真缺的那条腿＝§9 F-142-1） |
| **块 5b（`:1875/1878`，"seam 不是装饰"）** | 同样 0 命中，而且是**构造性**的：`main.go:221` 的 `scanWithStats` 正文就是 `return scanWithIgnore(root, newGitIgnore(root))`，所以 5b 比的是 `X` 与 `X`——今天除非 `scanWithStats` 多出步骤，它不可能不响 | ⚠ 同上：防将来分叉的引线，不是今天的牙。"seam 不是装饰"这句话的**实际凭据是块 5 在 M-B 下红的那 3 条**（它 §3.2 也正是这么记的，没记到 5b 头上）⇒ 措辞过关，账面别串线 |
| AC#4「门禁四数＋八枚数不降」 | 两枚纯净快照（`1b98c7b` vs `cd87354`）各自真跑（§7）⇒ 差值只有 `+1 枚用例`；八枚数只差 `frontend 43→46`，归因到别人那一枚 commit | ✅ 可反证，不是恒真 |

**这一格要不要动被验物？不。** 5a/5b 都是"多出来的正向对照"，摘掉它们不会让任何今天响的判据变哑
（块 3/块 4 的同位情形我用 AR-3/AR-5 那种"删断言再跑"的手法验过，5a/5b 的 0 命中是 18 枚变异的实测集合）。它们的毛病是**账面归因**容易被下一位读成"生产那一跳已被钉住"。⇒ 归口 §9 的 F-142-1（补真腿），
外加一句给编排者的更正建议：自证表 §3.4 那句"存在一种变异能打不红任何其它东西"**实测不成立**（AR-3 仍红），
按 append-only 补一行即可，不必退回重做。

⚠ **行号口径**（本仓那条"引行号必须同引它取于哪一版"在这里同样要守）：表里 `AR-3` 的 `:1799` 与 `AR-5` 的
`:1810` 是**那两枚变异各自删掉一段断言之后的**行号；同一枚因在**交付版**上的站点都是 `:1813`
（块 5 的静止前提 `t.Fatalf`）。`S-10` 没删测试文件的行，所以它的 `:1813` 就是交付版行号。

---

## 7. 门禁与契约轴：我自己跑的那一遍（四数／名册／Skip／冻结件）

### 7.1 `sh tools/d22scan/runtests.sh -C tools/d22scan ./...`（两枚纯净快照，同一台机器）

命令与形态：`git archive <锚> | tar -x` 到 `D:\tmp\accept142\snap-{old,head}`，进快照目录跑同一条脚本
（快照里没有 `.git`，现量 `NO .git (pure snapshot)`；这是 AC#4 指名要的形态）。

| 锚点 | PASS | FAIL | SKIP | === RUN | `^panic:` | `[no tests to run]` | runner 结论 |
|---|---|---|---|---|---|---|---|
| `1b98c7b`（改前） | 29 | 0 | 0 | 69 | 0 | 0 | `runtests.sh: OK` rc=0 |
| `cd87354`（交付） | **30** | 0 | 0 | **70** | 0 | 0 | `runtests.sh: OK` rc=0 |

**名册差集**取的是运行输出的 `--- PASS/FAIL/SKIP: NAME` 集合（不是我 grep 源码猜的），两枚快照各一份：
`diff` ⇒ `15a16 > TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`，**`+1 枚、`−` 0 枚**，
roster 枚数 29→30。⇒ 与它 §4.1 同值；源码级 `^func Test` 那把尺我也跑了（同 29→30、同名册），两把尺一致。

根目录 `go vet ./tools/d22scan/` 现量 **rc=1**，报文 `main module (github.com/CarlosShao/wisp) does not contain
package …/tools/d22scan`；`go.work` 不存在（`test -e go.work` ⇒ NO）。
⇒ 这是设计不是伤（简报这条前提我复算过，`.github/workflows/ci.yml:124-126` 那一步就是
`working-directory: tools/d22scan` 里跑 `go vet ./...`）。module 内 `go vet ./...` **rc=0**、`gofmt -l .` **空**。

### 7.2 `sh scripts/d22scan.sh` 的八枚数（四形态，每枚带锚点）

| 形态 | 锚点/树 | #1-5 internal/ | #1-5 cmd/ | #6 frontend/ | #7 internal/tools/ | #8 design/ | #8 frontend/ | #8 internal/ | #8 cmd/ | rc |
|---|---|---|---|---|---|---|---|---|---|---|
| 纯净快照 改前 | `1b98c7b` | 203 | 22 | 43 | 18 | 30 | 43 | 407 | 39 | 0 |
| 纯净快照 交付 | `cd87354` | 203 | 22 | **46** | 18 | 30 | **46** | 407 | 39 | 0 |
| 工作树（交付码，仓外副本 `go run . -root <仓库>`） | 工作树 | 203 | 22 | 46 | 18 | **32** | 46 | 407 | 39 | 0 |
| 工作树（同一枚树、改前码） | 工作树 | 203 | 22 | 46 | 18 | 32 | 46 | 407 | 39 | 0 |

- **A/B 独立复算**：两版码各跑一遍、同一枚工作树，`diff` 出的 scope 行＋自陈行＋`clean` 行
  **完全为空**（我打的标 `A_B_IDENTICAL`）；同一版码再跑第三遍仍 `HEAD_RUNS_IDENTICAL`
  ⇒ 它 §4.2 那行 `EIGHT_IDENTICAL` 我复现了，而且"静止索引上新分支零次触发"这句话在**当下这枚脏工作树**上
  也是真的（现量自陈第一行是 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`，
  **没有** `NOT APPLIED`）。
- **快照 `design/=30` 对工作树 `32`**：owner 那 16 枚未提交的 `design/` 挪动不进 archive（`git status` 现量在列），
  与本票无关，它 §4.2 的 ⚠ 已自报，我核过口径成立。
- **`43→46` 归因**：`git ls-tree -r --name-only <c> -- frontend | wc -l` 逐枚
  `1b98c7b`=43、`640cff1`=43、`d61281c`=46、`173a77d`=46、`cd87354`=46；
  `git log --oneline 1b98c7b..cd87354 -- frontend` 只列出 **`d61281c`** 一枚（前端会话的换屏层，
  新增 `render-nav.tsx`／`nav-rail.tsx`／`panel-views.ts` 三枚受追踪文件＝43+3）。
  ⇒ **这条差值不是本票造成的**，引用它时必须带这层归因（我核过，它 §4.2／§8 两处都写了）。

### 7.3 `t.Skip`：它自报 8→8，我把"这把尺看得见它"那一发自己复算

| 量法 | 读数 |
|---|---|
| `grep -c 't\.Skip'` 改前／交付 | **8 / 8**（本票新增 0 枚） |
| 我在副本里真植一枚（`snap-skip`，往该用例头里插一行 `t.Skip("planted known-positive for the acceptance ruler")`） | 同一条 grep ⇒ **9**，行号 `1683` |
| 同一份植入后的严格 runner | **rc=1**，`--- SKIP: TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules`（`skip-gate.log:133`），verdict 行 `1 test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3)`（`:226`），四数 `PASS=29 FAIL=0 SKIP=1 RUN=70` |

⇒ 简报那条 ⚠ 我照办：**没拿 `--- SKIP`＝0 当"没有跳过"**，而是把正控种进副本量。
两把尺（grep 与 runner）都看得见一枚真的 Skip。既有那 7 枚 `t.Skipf` 今天都不响（两枚快照 `--- SKIP` 均 0）。

### 7.4 断言只增不减（我自己数，不抄它的表）

`scan_test.go` 的删除列在三枚 commit 上都是 **0**（`git show --numstat` 现量：`9d06544` 214/0、
`64b404d` 37/0、`1ed231e` 11/0；区间累计 `262/0`）⇒ **没有任何一行既有测试被改或被删**，
这是能从 numstat 直接读出的强结论，我复算成立。计数字典（`git show <c>:… | grep -oE … | wc -l`）：

| 尺 | `1b98c7b` | `cd87354` | 差 |
|---|---|---|---|
| `Fatalf` | 43 | 52 | +9 |
| `Errorf` | 90 | 115 | +25 |
| `t.Fatal(` | 36 | 39 | +3 |
| `t.Skip` | 8 | 8 | 0 |
| `^func Test` | 29 | 30 | +1 |

⇒ **只增不减**，三枚计数差（`+9/+25/+3`）相加＝**37 枚断言调用点**，正好等于新用例体内
（`sed -n '1682,1903p' | grep -oE 't\.(Fatal|Fatalf|Errorf)\('`）现量的 `Fatalf 9／Fatal 3／Errorf 25`
⇒ 新增的每一枚断言都在这枚新用例里，旧用例一枚未动。
另：既有 7 枚 `t.Skipf` 我单独数过（`grep -c 't\.Skipf'` = **7**，加注释里那一枚词＝8），与本票无关。

### 7.5 契约轴：先证明东西存在，再说零命中，再拿同一把尺跑正控

**存在性（本程现量）**：`tools/d22scan/allowlist.txt` EXISTS、`docs/PLAN.md` EXISTS、`docs/specs` 目录 EXISTS、
`internal/observe/thresholds.go`（`git ls-files` 命中的唯一一枚 `thresholds.go`）EXISTS、
`grep -c 'emojiRe' tools/d22scan/main.go` = **4**、golden 一族 `git ls-files | grep -icE 'golden|testdata/.*\.(sse|json|txt|wav)'` = **52**。

**零命中（对 `1b98c7b..cd87354` 整区间，不只是对本票三枚 commit）**：

```
$ git diff --numstat 1b98c7b..cd87354 -- tools/d22scan/allowlist.txt internal/observe/thresholds.go docs/PLAN.md docs/specs '*/testdata/golden/*'
（空）
```

`emojiRe` 的定义块（`sed -n '/emojiRe = regexp/,/^$/p'`）两枚锚点**同一枚 sha1**（`61c0b60c…`）
⇒ 字符类一字节未动。`internal/**` 在区间内**零文件**（`git log --oneline 1b98c7b..cd87354 -- internal` 空；
工作树里那枚 `M internal/panel/l2_grant_boundary_test.go` 是**另一程未提交**的活，不记在本票头上）。

**本票自己三枚 commit 的路径集合**（`git show --name-only --format='' 9d06544 64b404d 1ed231e | sort -u`）：
`tools/d22scan/gitignore.go`／`main.go`／`scan_test.go` —— 只有这三枚。

**同一把尺的正控（防"尺根本没看"）**：对隔壁前端会话那一枚 `d61281c` 跑同一条
`git show --name-only` ⇒ 打出 8 枚路径（`frontend/package.json`、`frontend/src/App.tsx`…）。
⇒ 这枚"路径集合差集"尺子看得见越界，本票的零命中不是瞎读。
