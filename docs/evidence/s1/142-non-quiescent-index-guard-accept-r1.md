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
