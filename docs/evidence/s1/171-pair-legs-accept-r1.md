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

### ④ 判据③ 那句前提（"未修码点不到它"）的独立复算——两向
- **B 形（接住句柄、全文件不关）**：改前的尺**也点名**（§4② 的 control 行：旧尺三条腿全点 `bridge.go`；M2 另证：新加的 `internal/fake171v1/leak.go` 在旧尺上同样被点名）。
  ⇒ 票面那句"未修码上主尺点不到它"**对 B 形不成立**＝实现者报回的话**复算一致**，且它的数学解释站得住（同射程同开方下，把合方词根集合变大只可能让正向名册**变小**，点不出改前点不到的文件）。
- **C 形（`_ =p.OpenScope(id)` 丢句柄＋同文件 `l.CloseScope("task-other")`）**：台件 `v1-ruler-mutations.sh` S1（树 `internal/fake171v1c/c_discard.go`）——
  **改前的尺 G5 名册只有 `cmd/wisp/panel_assets.go`＋`internal/tools/bridge.go` 两枚、未成对枚数＝2、不点这一枚（rc=1）**；
  **改后的尺点名它，三条腿（G5／G5-正控／G5-负一负）＋`DISCARD-HANDLE` 批注，rc 1→3**。⇒ 这才是 AC#6③ 要的那发"改前不响→改后点名"，**两向都真跑出来了**。
- **这样结案算不算把洞说清楚**：算。票面 16:3x 追加的那段（原句未抹、按"两形分别成立"结案、并给出"红因换轨"这一口径）与我的独立读数逐字对得上，**不当缺陷报**。

**⇒ AC#6 判〔成立〕（①②③④ 四问全过）。**

<!-- §5（AC#2 三发）与 §6 之后各节由本程第二枚 commit 追加：先测后写、按格提交。 -->

