# 票 171 r1 验收表（非实现者）＝AC#6（合方词根认句柄那一形）＋ AC#2（枚数进判据）——主攻「它到底有没有减射程」

- 验收位｜程 `171-v1`｜时刻 `2026-09-27 17:0x +08`｜派单＝`.scratch/wisp/dispatches/2026-09-27-164x-accept-171-r1-v1-did-it-shrink-the-scope.md`
- 被验对象＝`3df749607a0da1e70579f5c3abb71003350a45d2`（尺＋台件＋表）与 `0cb6e8f2fe54aef082f0a5f4a81aa932efee0e79`（票面进度）
- 台件＝`.scratch/wisp/probes/171/v1/**`（本程唯一写面之一，另一写面＝本件）｜尺的字节**一律 `git show` 到仓外 `$TMPDIR`**，`probes/154/**` 与 `probes/171/r1/**` 一字节未动（现量：`git status --porcelain -- .scratch/wisp/probes/154 .scratch/wisp/probes/171/r1 internal cmd tools .github docs/PLAN.md docs/specs docs/reports` ⇒ **空**）

## 0. step-0（分支＋每枚号复核＋锚点策略）

```
$ git rev-parse --abbrev-ref HEAD   ->  dev
$ git rev-parse HEAD                ->  b1cbdad43f4deee7009c89fb6c251455e9277985   （取数途中兄弟程又推进到 4193f3a4，见 §8）
$ git cat-file -t 3df74960          ->  commit   3df749607a0da1e70579f5c3abb71003350a45d2
$ git cat-file -t 0cb6e8f2          ->  commit   0cb6e8f2fe54aef082f0a5f4a81aa932efee0e79
$ git cat-file -t 19513cce          ->  commit   19513ccefa45d51206871f00b3efa6ab19dd620a   （派单点名的"改前那版尺"）
$ git cat-file -t ff000784          ->  commit   ff0007840ac9489c6c5ef017ca5092cdd2ed6094   （171-r1 的父节点，旁证）
```

**锚点策略（派单头两条＝本格成不成立的先决条件）**：本程所有名册／腿读数**全部钉在 commit 上取**，没有一枚读数来自工作树——
读尺的输出走 `git grep <pat> <sha>`，跑尺的工作目录要么是真仓（锚作参数传进去），要么是仓外的丢弃树。
本程用的三个锚：**`3df74960…`（被验对象，主锚）**、**`19513cce…`（改前那版尺的字节出处＋第二锚）**、
**丢弃树自己的 HEAD（`b6cdb7b`／`e888913`／`169f9b3`／S1 那枚）**。兄弟程正在写 `internal/tools/**` ⇒ 那些提交**不进本表任何一格**。
尺的字节交叉核对：`git diff 19513cce 3df74960 -- .scratch/wisp/probes/154/gate-clauses.sh` 的 `index 80c22482..f60792c2`
＝本程 `git hash-object` 抽出的两版字节 `80c2248274244e544e28e7dbc3fba672a4c0b82e`／`f60792c2348740bf3ff70e822264166905007f8c`（同两枚 ⇒ 抽的就是那两版）。

## 1. 本程没测什么（先把没做的说清楚）

1. **票面任何一格都没勾**，`0cb6e8f2`／`3df74960` 的字面也没改（写面只有本件＋`probes/171/v1/**`）。
2. **AC#7（`run()` 门上同一枚洞）本程不判**（派单明令）。**只写一句给读者**：本表 §5 判的是 `pair` 腿；
   **同一枚"减法侧无牙／枚数只往增加那侧装牙"的洞在 `run()` 那一族（G1–G4，仍是空／非空两态）还开着，已由票面 AC#7（`issues/171-…:46`）接单** ⇒ 别把本表读成"判据族全修好了"。
3. **AC#6①／②（尺的聚合退码、不改射程凑 0）我复算的是读数，不是它内部实现的对错**；AC#1（尾随注释）、AC#3（接 CI）、AC#4／AC#5（契约轴）不在派单两格里，本程零判定。
4. **`CLOSE-BY-GENERIC-CLOSE` 那一格我没扩**（派单明令"别顺手把它扩成能修"）：它要数据流，已记进本票"不解决"。本程只回答"今天这一版比改前看得见更多还是持平"（§4③）。
5. 没跑全量 `go test ./...`（派单只点名 `tools/d22scan` 那把 runtests），没跑 `d22scan` 的 CI 触发面，没跑 golden／SLO。
6. 丢弃树里的样本**没有编译过**（尺只吃 `git grep`，样本只为造读数；`gofmt` 过了，见 §3 台件里的 `gofmt -w`）。这条是"没测"，不是"测过没问题"。
7. 票 162 兄弟程的读数（`b1cbdad4`／`4193f3a4` 上的 `internal/tools/**`）**一律不作本表判据**，只用于 §8 那两把门禁的"门还能跑"。

## 2. 票面两格（连行号逐字抄，出处＝`.scratch/wisp/issues/171-three-ruler-holes-…-count-blind-pair-legs.md`）

```
27  - [ ] **AC#2（尺洞③＝枚数）**：把成对腿的判据从"空／非空"改成**带枚数**（声明"该响几枚"或"枚数只增不减"，两形选一形并写明为什么）。**硬判据**：造一发"同一条腿内未成对从 1 枚变 2 枚"的台件 ⇒ 聚合退码**必须变**；这发在**未修码上必须不响**（否则本格是装饰）。
33  - [ ] **AC#6（09-27 11:5x 追加，来路＝票 160 的 r1 落地后编排者现量复现）合方词根写死成 `CloseScope` ⇒ …G5 的正控那腿言行不一、聚合退码从 0 变 1，而主尺多假点一枚。**
37    ① 修完 `gate-clauses.sh` 的**聚合退码回到 0**（14 腿零枚言行不一），且这一发必须是你真跑出来的、不是推的；
38    ② **不许靠把 `internal/tools/bridge.go` 从 G5 主尺射程里排除掉来凑 0**——那是放水…改射程必须与"认得句柄那一形"**同时**发生；
39    ③ **反向判据（这一发在未修码上必须不响）**：造一枚"用句柄形状**开而不关**"的样本…⇒ 未修码上主尺**点不到它**、修完之后**必须点名它**。
46  - [ ] **AC#7（…）**：`pair` 腿今天带了 `want_n` 枚数，而 **`run()` 那一族腿（G1–G4）仍是"空／非空"两态** ⇒ AC#2 收掉的那个洞在隔壁那扇门上还开着。
```

## 3. 主攻那一节的差集读数（"我一条 pathspec 都没减"是机器证的，不是读它的话）

台件＝`v1-run-rulers.sh`（两版尺 × 两枚锚＝四份读数 `logs/gate-{old,new}-{19513cce,3df74960}*.txt`）
＋`v1-leg-diff.sh`（逐腿差集 `logs/leg-diff-{19513cce,3df74960}.txt`、`legs/*.tsv`）。
两版尺都从 commit 抽字节（`git show <sha>:.scratch/wisp/probes/154/gate-clauses.sh`），跑在**同一条射程、同一枚锚**上。

### 3.1 射程本体：逐腿 `git grep` 命令行对照（九腿，两枚锚各一遍，读数相同）

`-- G5 / G5-正控 / G5-负一负 / G6 / G6-正控 / G6-负一负 / G7 / G7-正控 / G7-负一负` 九腿全部 **[SAME-pathspec]**，
唯一变化在 pattern 段（合方词根 `'CloseScope'` → `'CloseScope|Close\(\)'`）。逐字样本（锚 `3df74960…`）：

```
G5        old: … 'OpenScope|CloseScope'              … -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/risk/*
G5        new: … 'OpenScope|CloseScope|Close\(\)'    … -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/risk/*
G5-正控   old/new 同上，只多一条 ':!cmd/wisp/panel_assets.go'（两版都有，一字未动）
G5-负一负 old/new … -- internal/**/*.go cmd/**/*.go            （两版都无排除）
G6／G6-正控／G6-负一负／G7／G7-正控／G7-负一负：old 与 new 的整行**逐字相同**（含 ':!internal/tools/bridge.go'、':!internal/plugin/disposal.go'、正控那两腿的收窄 pathspec）
```

⇒ **"一条 pathspec、一枚 `:!` 都没动"这句话在九腿上全部成立**，G6／G7 那六腿连 pattern 段都没变。

### 3.2 UNPAIRED 名册差集（comm 两向；两枚锚读数逐字相同）

| 腿 | 改前 | 改后 | 退场 | 新增 |
|---|---|---|---|---|
| G5 | 2 枚 | 1 枚 | `internal/tools/bridge.go` | （空） |
| G5-正控 | 1 枚 | 0 枚 | `internal/tools/bridge.go` | （空） |
| G5-负一负 | 9 枚 | 8 枚 | `internal/tools/bridge.go` | （空） |
| G6／G6-正控／G6-负一负 | 1／0／1 枚 | 1／0／1 枚 | （空） | （空） |
| G7／G7-正控／G7-负一负 | 3／0／4 枚 | 3／0／4 枚 | （空） | （空） |

九腿并集：**新增点名＝0 枚；退场＝恰 1 枚路径 `internal/tools/bridge.go`，跨 3 条腿（G5／G5-正控／G5-负一负）**；
名册行合计 21 → 18（＝同一枚文件 × 那 3 条腿）。锚 `19513cce…` 与锚 `3df74960…` 两枚锚**逐字同读**（⇒ 两锚之间那 5 枚提交没有一只落进这九腿射程，见 §8 末）。

⇒ **实现者声明的那三个数（新增 0／退场恰 1／同一枚路径）复算一致。多一枚就是减了射程——没有多。**

可复制命令（读数已落盘，重跑同一枚台件即可）：

```
sh .scratch/wisp/probes/171/v1/v1-run-rulers.sh 19513cce… 3df74960…
sh .scratch/wisp/probes/171/v1/v1-leg-diff.sh   3df749607a0da1e70579f5c3abb71003350a45d2   # 全文＝logs/leg-diff-3df74960.txt
```

## 4. AC#6 四问，逐问读数

### ① 减没减射程 → **没减**（§3 两向机器证）。聚合退码也复算到：
新尺在锚 `3df74960…` 与 `19513cce…` 上都是 **rc=0／`腿数＝14 声明与实测不符＝0`**；改前那版尺在同两枚锚上都是 **rc=1／不符＝1**（`BAD 腿=G5pos`）。⇒ AC#6① 那句"真跑出来的 rc 回到 0"＝**成立**。

### ② 那一发"只开不关"＝本格的判死发（我自己造的，不是抄它的）
台件＝`v1-mutation.sh`。做法：`git archive 3df74960 internal cmd` 到仓外 `$TMPDIR` 的**丢弃 git 树**（`git init`＋`git add internal cmd`＋commit），
变异 M1＝`sed '/closeErr = scope\.Close()/d'` 摘掉 `internal/tools/bridge.go:706` 那一行（`git show --stat` 证落点＝`internal/tools/bridge.go | 1 -`）；
本仓 `internal/**` **零字节落件**（§0 的 status 现量＝空）。丢弃树 sha：control `b6cdb7b`／M1 `e888913`。

| 树 | 改前的尺点名 `bridge.go` 的腿 | **改后的尺**点名的腿 | 新尺聚合退码 |
|---|---|---|---|
| control（逐字节＝锚点） | G5／G5-正控／G5-负一负（旧尺盲认 Close，按设计必点） | **（无）** | **rc=0** |
| **M1 unpaired-bridge（只开不关）** | G5／G5-正控／G5-负一负 | **G5／G5-正控／G5-负一负（重新点名）** | **rc=3**（三腿各比基线多 1 枚） |

⇒ **"认句柄那一形"是真认，不是假认**：新尺只在"接住的句柄确实被关掉"时安静，Close 那一行一摘就立刻重新点名三条腿。
**AC#6② 判〔成立〕。**

### ③ 它摊出来的那句 `CLOSE-BY-GENERIC-CLOSE` 够不够"给人判"
- 真尺（锚 `3df74960…`，新尺）在 G5／G5-正控／G5-负一负三条腿上**逐腿都把退场那枚摊回来**：`logs/gate-new-3df74960….txt` 第 127／241／530 行 `#   CLOSE-BY-GENERIC-CLOSE internal/tools/bridge.go`。
- 派单那一问「句柄被丢掉＋同文件另有一句无关的 `Close()` ⇒ 今天会不会又静默」我自己造了样本 D（台件 `v1-sample-d.sh`，树 `internal/fake171v1d/d_discard_plus_generic_close.go`：`_ = p.OpenScope(id)` ＋ `_ = f.Close()`）：
  **两把尺都点名**（改前 rc=1 点 G5／G5pos／G5neg；改后 rc=3 点同样三条腿，并额外摊出 `DISCARD-HANDLE` ＋ `CLOSE-BY-GENERIC-CLOSE` 两行批注）。
  ⇒ **今天这一版比改前【看得见更多】，不是持平、更不是更少**：唯一换来的安静就是"接住句柄又用句柄关掉"那一形（§3.2 恰 1 枚），而那枚一被摘就重新响（§4②）。
- `CLOSE-BY-GENERIC-CLOSE` 本身仍然只是"给人判"的一句话（文件级、要数据流才能判成尺）。这一格**没扩、也不该在本票扩**——它已在票 171 的"不解决"里；票面 §12 next=① 说的是同一件事。
- 现量一枚计数偏差（**不影响判定、只是文字账**）：尺内新注释写"非测试生产码 **40** 枚文件里有 `Close()`"，
  我现量 `git grep -lEw 'Close\(\)' <锚> -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l` ⇒ **42 枚（两枚锚都是 42）**＝被验表 §4 写的数。⇒ **表的数对、尺里那句注释的数错（40 应为 42）**，下一程顺手核一下即可，本程没动它一字。

**⇒ AC#6 判〔成立〕（①②③ 三问过，第④问＝票面那句前提的复算在 §6）。**

## 5. AC#2（枚数进判据）三发——它选的形、有没有牙、藏不藏东西

**它选的形复算确认**（字节层，`git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh`）：`want_n <基线枚数>` 紧跟 `want`、只对 `pair` 腿生效；
`pair` 腿的记账换给 `book_count()`，判据是**两味并集**——① 票 161 的空／非空（`bin != EXPECT ⇒ BAD`）逐字保留、② 新增`实测 > 基线 ⇒ BAD`；
`实测 < 基线 ⇒ SHR`、**不计进退码**；`want_n` 缺失 ⇒ `exit 4`。现量：新尺里 `^want_n ` 登记行数＝**9**＝九枚 pair 腿全登记，
锚 `3df74960…` 上九腿`基线=N枚 实测=N枚`逐枚全等（`G5 1／G5pos 0／G5neg 8／G6 1／G6pos 0／G6neg 1／G7 3／G7pos 0／G7neg 4`，原文在 `logs/gate-new-3df74960….txt`）。
**它拒绝"逐枚相等"的理由我认**：相等形会让"把 `panel_assets.go` 补上一发收尾"这种**改好了**变成红——那正是票 161 AC#7／台账 A320、A321 付过钱的恒红形状。

### 5.1 有牙那一发（我自己塞的，不是它的台件）——两向都动／都不动
台件 `v1-mutation.sh` M2：丢弃树 control（逐字节＝锚点）之上加一枚 `internal/fake171v1/leak.go`（接住句柄、全文件不关）⇒ **同一条腿 G5 内未成对 1 枚 → 2 枚**。

| 尺 | control（1 枚） | extra-leak（2 枚） | 这一腿的判定 |
|---|---|---|---|
| 改前的尺 | **rc=1** | **rc=1（不动）** | `# ok 腿=G5 声明=ring 实测=ring`（旧尺只比空／非空） |
| 改后的尺 | **rc=0** | **rc=3（动了）** | `# BAD 腿=G5 声明=ring 基线=1枚 实测=2枚 因=新增未成对（票 171 AC#2：实测 > 基线）` |

⇒ **未修码上这发不响、修完才响**＝票面 AC#2 那句硬判据在**我自己的台件**上复算成立；"新增的第二枚躲在非空里"这件事在改前的尺上确实躲得进（旧尺 rc 一字不动），在新尺上进不来了。**AC#2〔成立〕**。

### 5.2 `SHR` 那一味藏不藏东西：**藏，但只藏"减法"这一侧、且有边界**（这一格得点名）
台件 `v1-ruler-mutations.sh` S2＋S2b（尺的副本在仓外，`probes/154/**` 未动）：

```
S2  把 G7 主尺射程悄悄调小（追加 ':!internal/observe/*' ':!internal/proc/*'，want_n 仍写 3）
    ⇒ # SHR  腿=G7 声明=ring 基线=3枚 实测=1枚 注=读数变好了，基线过期，下一程把 want_n 核下来
    ⇒ # 腿数＝14 声明与实测不符＝0 ／ 聚合退码＝0        ← 退码【不动】
S2b 把同一腿射程整条清空（读数 3→0）
    ⇒ # BAD  腿=G7 声明=ring 基线=3枚 实测=0枚 因=空/非空那一味与声明不符 ／ 聚合退码＝1   ← 响
减法侧的第二形（台件 v1-declaration-flips.sh 末发）：整腿删掉（去掉 G7-正控那一发 pair 调用）
    ⇒ # 腿数＝13 声明与实测不符＝0 ／ 聚合退码＝0        ← 退码【不动】，且尺【没有任何一处断言腿数必须是 14】
```

可复制命令（S2／S2b 的尺副本一律 `git show` 到仓外，射程那一行只在这份副本里被改）：

```
T=$(mktemp -d); git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh > $T/gate.sh
awk -v add="':!internal/observe/*' ':!internal/proc/*'" 'index($0,":!internal/plugin/disposal.go")      && index($0,"GO2"){print $0 " " add; next}{print}' $T/gate.sh > $T/gate-shrunk.sh
sh $T/gate-shrunk.sh 3df749607a0da1e70579f5c3abb71003350a45d2      # ⇒ logs/s2-shrunk-range.txt：SHR、rc=0
# S2b＝同一形再加 ':!internal/memory/*'（读数 3→0）⇒ logs/s2b-range-to-zero.txt：BAD、rc=1
```

⇒ **判这是不是一条静默通道：是。** 只要一条腿的读数还 ≥1，把它的射程往小调（或整腿删掉）都**不会动聚合退码**；
唯一的兜底是"掉到 0 且声明 ring"时票 161 那一味会响。**它不是装饰**（`SHR`／`腿数` 两行都摊在 stdout 里、`基线` 是写在尺旁边的可读数，谁调大调小都在 diff 里露头），
但**"名册缩小"这一侧没有机器牙**。这条我**不判成本格的失败**——票面 AC#2 的原话就是"新增的第二枚躲进了非空"，靶子在**增加**那一侧，实现者按票面造牙、
并且明写了"变好不算言行不一"的理由；把减法侧也装牙＝"逐枚相等"那一形，正是它论证过不要的恒红。
⇒ **要装的牙属于另一格**：本程把它记进 §14 next=②，交编排者定（并给 AC#7 还是另立一格），**不要求实现程在本票修**。

### 5.3 空心那一发（缺 `want_n` 必须死）
台件 `v1-ruler-mutations.sh` S3：只摘掉 G5 那一行 `want_n 1` ⇒ **rc=4**，末行逐字：

```
想死在起跑线上：pair 腿 G5 没登记 want_n 基线——枚数不进判据＝票 171 AC#2 那枚洞还开着
```

且它**没有出结论**：`# 腿数＝`／`# 聚合退码＝` 两行 grep 零命中（尺死在登记处，不会给出"干净"）。⇒ 下一程新加一枚 pair 腿**不可能**悄悄退回空／非空那一形。

> 再钉一句给读者（派单 §1 末）：**同一枚洞在 `run()` 那一族（G1–G4）的门上还开着**——那四腿仍是空／非空两态、没有枚数，
> **已由票面 AC#7（`issues/171-…:46`）接单**，本程不判。**别把 §5 读成"判据族全修好了"。**

## 6. AC#6 判据③ 那句前提的独立复算（两向都是自己的读数）
**复算对象＝票面 `:39` 那句「未修码上主尺点不到它」**
- **B 形（接住句柄、全文件不关）**：改前的尺**也点名**（§4② 的 control 行：旧尺三条腿全点 `bridge.go`；M2 另证：新加的 `internal/fake171v1/leak.go` 在旧尺上同样被点名）。
  ⇒ 票面那句"未修码上主尺点不到它"**对 B 形不成立**＝实现者报回的话**复算一致**，且它的数学解释站得住（同射程同开方下，把合方词根集合变大只可能让正向名册**变小**，点不出改前点不到的文件）。
- **C 形（`_ =p.OpenScope(id)` 丢句柄＋同文件 `l.CloseScope("task-other")`）**：台件 `v1-ruler-mutations.sh` S1（树 `internal/fake171v1c/c_discard.go`）——
  **改前的尺 G5 名册只有 `cmd/wisp/panel_assets.go`＋`internal/tools/bridge.go` 两枚、未成对枚数＝2、不点这一枚（rc=1）**；
  **改后的尺点名它，三条腿（G5／G5-正控／G5-负一负）＋`DISCARD-HANDLE` 批注，rc 1→3**。⇒ 这才是 AC#6③ 要的那发"改前不响→改后点名"，**两向都真跑出来了**。
- **这样结案算不算把洞说清楚**：算。票面 16:3x 追加的那段（原句未抹、按"两形分别成立"结案、并给出"红因换轨"这一口径）与我的独立读数逐字对得上，**不当缺陷报**。


## 7. 声明机制那两发（翻一枚／翻两枚／还原）＋ flip-declaration.sh

我自己翻（尺的副本在仓外 `$TMPDIR`，**`probes/154/gate-clauses.sh` 一字节未动**；台件 `v1-declaration-flips.sh`，锚 `3df74960…`）：

```
f0 未翻（＝基线）          rc=0   # 腿数＝14 声明与实测不符＝0
f1 翻一枚 G7pos quiet→ring rc=1   BAD 腿=G7pos … 因=空/非空那一味与声明不符（票 161 AC#6①）
f2 再翻一枚 G6neg ring→quiet rc=2  两枚 BAD 逐行都在 ⇒ 退码＝2、不是"任何一腿错⇒1"
f3 把 f2 翻回去             rc=0（cmp 证还原后的尺与 f0 **逐字节相同**）
```

门禁那把（派单 §2 点名）：`sh .scratch/wisp/probes/161/r6/flip-declaration.sh` ⇒ **rc=0**，末两行逐字：
`flip-declaration.sh: GREEN - the aggregate exit code tracks the declarations in both directions and returns to 0 (baseline rc=0, 14 legs booked).`
＋`…the baseline books 7 ringing legs (G2 G5 G5neg G6 G6neg G7 G7neg) and 7 silent ones … and still exits 0.`（`FAIL` 行数＝**0**）
⇒ **`want_n` 没把这套声明换成"永远说到做到"的空尺**：声明仍逐枚连着退码，两枚同翻仍是 2，翻回仍是 0。
⚠ 副作用同实现者报的那条：跑这枚门禁会往 `probes/161/r6/logs/flip-*.txt`（**8 枚 tracked 件**）写 log ⇒ 本程那 8 枚**不提交、不还原**。

## 8. 门禁读数（派单 §2 四把，逐把跑、没跑全量）

| 门禁 | 我的现量 | 与派单给的断言 |
|---|---|---|
| `sh .scratch/wisp/probes/154/gate-clauses.sh`（工作树里的尺，默认锚＝HEAD `4193f3a4`） | **rc=0**／`腿数＝14 声明与实测不符＝0` | **同色**（编排者 16:2x 现量一致） |
| 同一把尺钉在锚 `3df74960…` 与 `19513cce…`（§0 的策略） | 两枚锚皆 **rc=0／不符＝0** | 名册判据只用这一行，不用上面那行 |
| `sh .scratch/wisp/probes/161/r6/flip-declaration.sh` | **rc=0**／GREEN／FAIL 行数 0 | 同色 |
| `sh scripts/d22scan.sh` | **rc=0**（`d22scan: clean - no D22 ban violations`） | 只当"门还能跑"（`internal/=420`、`cmd/=45` 是 HEAD 已被兄弟程推进后的数，**不作名册判据**；被验表里那枚 `internal/=418` 是它自己锚上的数） |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | rc=0；`packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` | 四数逐字同它报的基线 |
| runtests 名册两向 `comm`（66 枚唯一名 vs `probes/171/r1/baseline-runtests-pre-fix.txt`） | 66 枚／`comm -23`＝空、`comm -13`＝空 ⇒ **名册逐字相同** | 一致 |
| gofumpt | 本程**仓内零枚新 `.go`**（台件全是 `sh`；样本只在 `$TMPDIR`，`gofmt -w` 过） | 不适用（无新 Go 件） |

另两枚"名册没漂移"的旁证：① 两枚锚 `19513cce`↔`3df74960` 上九腿读数逐字相同 ⇒ 被验那 5 枚提交没动这九腿射程；
② §3.2 的差集只由**尺**解释得通（pathspec 九腿全同、只换合方词根）。

## 9. 被拒／没成功的调用（各标取数之前／之后）

- **权限层拒绝：0 次**（没有一次工具调用被门禁挡掉，也没有改过任何权限设置）。
- **我自己写坏、全部副作用只在台件／仓外**（六次，逐条）：
  1. **取数之中（§3 之前）**：`v1-leg-diff.sh` 首跑把 `UNPAIRED` 行的字段取成 `$2` ⇒ 差集读数不可用（"退场=[UNPAIRED ]"）。修成 `$2":"$3`、连 `-REV` 一起收 ⇒ 才是 §3.2 那张表。**没产生仓内副作用。**
  2. **§4 取数之前**：`v1-mutation.sh` 首跑 `git -C <我的台件目录> archive … internal cmd` ⇒ pathspec 相对 cwd、`fatal: pathspec 'internal' did not match any files` ＋ tar 报错；改成从仓根跑 `git archive`。
  3. **§4 取数之前（同一台件第二版）**：跑尺时 `cd` 进丢弃树、log 路径却是相对路径 ⇒ `No such file or directory` 与 `rc: unbound variable`；改成绝对 `S=$(cd "$D/logs" && pwd)` ＋ `GIT_DIR`/`GIT_WORK_TREE` 显式指向丢弃树（否则尺会打到真仓）。
  4. **§5 取数之前**：`v1-ruler-mutations.sh` 的 S2 一处 shell 引号写坏（`add="…"'`）⇒ `syntax error near unexpected token '('`，那一轮 S2／S3 没跑；修后重跑。
  5. **§5 取数之前**：同台的 S1 树当时只 `git add internal`（漏 `cmd`）⇒ 旧尺 G5 名册少点 `panel_assets.go`、G6 两腿读成 quiet。修后 S1 重跑（§4④ 引的是修后的读数；错的那轮读数是"旧尺名册只有 bridge.go 一枚"，我没用它下判）。
  6. **§6 取数之前**：`v1-declaration-flips.sh` 首跑 flip 规格写成空格分隔（应写 `LEG|旧|新`）⇒ 报 `!! 翻 G7pos 的声明没命中` 并 `exit 9`；修后重跑（§6 是修后的读数）。
- **兄弟程干扰：0 次读数被污染**——我所有名册／腿读数取自 commit（§0）；HEAD 在取数途中从 `b1cbdad4` 推进到 `4193f3a4`（162 系那几枚），受影响的只有 §7 那把全仓级门禁的计数（`internal/=420`），已明写"不作名册判据"。

## 10. 有没有跑过删除命令

- **仓内：零次** `rm`／`git rm`／`clean`／`checkout`／`restore`／`reset`／`stash`／`--amend`／`rebase`。现量：`git status --porcelain -- .scratch/wisp/probes/154 .scratch/wisp/probes/171/r1 internal cmd tools .github docs/PLAN.md docs/specs docs/reports` ⇒ **空**。
- **仓外 `$TMPDIR` 丢弃树内**：`rm -rf "$T/base"`（台件④／⑥各一次，清自己的临时目录）、`sed` 生成副本而非就地改、`mv` 只在 `$TMPDIR` 内挪尺的副本。样本文件都建在丢弃树里，**`internal/**` 一字节未落件**。
- 别人的未提交件（`design/**` 的 16 枚删除、`probes/152/my152.py`、`docs/evidence/s1/152-…-accept-r1.md`、`probes/162/**`、`probes/161/r2| r5` 的既有未跟踪件、以及我自己跑门禁脏掉的 `probes/161/r6/logs/flip-*.txt` 8 枚）：**不提交、不还原、不补完**。

## 11. 伪授权两栏

- **我没做过、派单也没授权我做的事**：改 `D1–D47`／`C1–C32`／`R1–R9`／D43 转移表；动 SLO／golden／`thresholds.go`；动禁令射程或 `allowlist.txt`；
  勾票 171（或任何票面）的框、改票面 Status／Progress log；改 `probes/154/**`、`probes/171/r1/**`、`internal/**`、`cmd/**`、`tools/d22scan/**`、`docs/reports/**` 的字节；
  把 `CLOSE-BY-GENERIC-CLOSE` 扩成可判（派单明令"别顺手修"）；切分支／建 worktree／push／改写已推送历史。
- **看起来像授权、我按"不是"处理的事**：① 派单里那四条断言（`rc=0／不符 0 枚`、`退场恰 1 枚`、`flip 现在 rc=0`、`缺 want_n 会死 rc=4`）**全按未验证断言逐枚现量**，
  其中尺内注释"非测试生产码 **40** 枚文件里有 `Close()`"我量到 **42**（两枚锚同数）⇒ 按派单"不符就报回、不许为对上一句话改判据"处理：**报回、没改尺、也没改判据**；
  ② 工作树正被兄弟程写 ⇒ 我没拿任何工作树读数当判据（连 `d22scan` 的计数也只当"门还能跑"）；
  ③ "票面 AC#6 判据③ 那句前提已被实现程顶回一半、已按两形分别成立结案" —— 我**独立复算了**（§4④，两向都是自己的读数），
  复算支持它那句更正，于是**按派单说的"别把它当缺陷"**处理；但我也**没有**顺手把票面那句话改写成与实测一致（票面禁改）。

## 12. 凭据值零抄录

本程读过的东西只有：尺的 shell 字节、`git` 元数据（sha／名册／status）、日志文本与文件路径。
没有读过、也没有抄录任何 API 密钥／DPAPI 密文／token／`.env`；本件里出现的字符串只有路径、commit 号、函数名与腿号。

## 13. 判定汇总

| 格 | 判 | 依据（本格自己的读数） |
|---|---|---|
| **AC#6** | **成立** | ① 新尺两锚皆 rc=0／不符＝0（旧尺皆 rc=1／G5pos BAD）；② 九腿 pathspec 逐字相同、名册新增 0 枚／退场恰 1 枚同路径（§3）；③ 我造的 M1（摘掉 `scope.Close()` 那一行）让新尺立刻重新点名三条腿＝句柄那一形真认；④ C 形改前静默／改后点名、B 形改前也点名（票面那句对 B 不成立＝实现者报回的话复算一致） |
| **AC#2** | **成立**（附一条不判死它的名册账） | 九腿全登记 `want_n`、基线==实测逐枚；M2：1 枚→2 枚 ⇒ 旧尺 rc 不动、新尺 rc 0→3（硬判据两向）；S3：缺 `want_n` ⇒ rc=4 且不出结论。**账＝§5.2 那条"减法侧无牙"（射程调小／整腿删掉都不动退码），它不在票面 AC#2 的靶子上，落 §14 next=② 交编排者定** |

## 14. next=

1. **交回编排者**：AC#6／AC#2 两格**都不勾**（勾归编排者）。若按本表结案，两格可勾。
2. **给"减法侧无牙"立案**（本程现量：`SHR` 不响退码、整腿删掉不响退码、`腿数` 只打印无断言）——
   最小闭合集合是三处、都很小：① `legs()`／聚合段加一枚"腿数＝登记期望值"的断言（期望值写进尺旁边、与九腿 `want` 同源）；
   ② `SHR` 那一味改成"允许变好，但必须在 stdout 单列一张『基线过期』表"并让 `flip-declaration.sh` 那把核它（现在只有人眼看得见）；
   ③ 与 AC#7 同程落（`run()` 腿一起带枚数），否则两扇门的牙还是不齐。**这条不是本派单的判据，我没要求实现程在本票修。**
3. 一枚文字账：尺内新注释"非测试生产码 40 枚文件里有 `Close()`" ⇒ 现量 **42**（`git grep -lEw 'Close\(\)' <锚> -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l`，两枚锚同数）。下一程顺手核一字。
4. `CLOSE-BY-GENERIC-CLOSE` 维持"给人判"：要接成可判的东西需要数据流，已记在票 171 的"不解决"；本程未扩、也不建议在本票扩。
5. AC#7（`run()` 门上的同一枚洞）与本票 AC#1（尾随注释）**都还没做**；做 AC#1 那一程要一并把 `DISCARD-HANDLE`／`CLOSE-BY-GENERIC-CLOSE` 两种行形状算进名册 diff（`pair()` 现在每腿会多打这几行，逐枚点名行的形状变了）。


