# 72-a2 — 必答 A：票 307 落地之后，这 7 枚还红的机制

## A.1 结论一句话

**这 7 枚与 `internal/risk/syncdirs.go` 因果＝0。** 它们的红全部落在**两枚测试文件自己体内**（`pathresolver_anchor_spelling_windows_test.go` ＋ `pathresolver_junction_windows_test.go`），机制＝**测试把 `t.TempDir()` 的返回拼写当成了"唯一拼写"或"OS 担保拼写"来用**；产码侧一行都没错（三处被误读的产码行为逐枚列在 A.3）。

## A.2 为什么因果＝0（四条独立尺，件名随附）

| # | 尺 | 读数 | 件 |
|---|---|---|---|
| 1 | `git show --numstat --format= 51a0a0f2`（票 307 的产码笔） | `internal/` 面只有两行：`17 3 internal/risk/syncdirs.go` ＋ `141 0 internal/risk/syncdirs_alias_307_windows_test.go` ⇒ **`pathresolver*.go`／`blacklist.go` 一字节⛔ 动**；其余 27 行全在 `probes/307/r1/**` | 本程 stdout（正文抄入本文件 A.5 末） |
| 2 | `git grep -c -E 'IsSyncPath\|syncSet\|NewProvenance' HEAD -- internal/risk/pathresolver_anchor_spelling_windows_test.go internal/risk/pathresolver_junction_windows_test.go` | **0 命中（`rc=1`）** ⇒ 这 7 枚的调用图⛔ 进 `syncdirs.go` | 同上 |
| 3 | `git grep -n -E 'syncSet\|IsSyncPath\|resolveTarget\|syncSuspect' HEAD -- internal/risk/pathresolver.go internal/risk/blacklist.go` | **0 命中（`rc=1`）** ⇒ 反向也成立：`Resolve`/`formsOf`/`Classify`/`Gate` 这条链⛔ 调 `syncdirs` 的任何东西 | 同上 |
| 4 | 两台面 × 同一份码、只差一枚 env：别名台面 `-run '<7 枚>'` ＝ `rc=1`/FAIL=**7**/PASS=0/SKIP=0；默认台面 ＝ `rc=0`/FAIL=0/PASS=**7**/SKIP=0 | 甲族的红**⛔ 经过**任何 307 改动即可复现，且 307 落地后**一枚不少一枚不多**（与编排者 `A834` 那四发的"改后剩 7 枚＝甲族原样"同形，但那是他们两把整包读数；本腿只跑甲族 7 枚这一把） | `logs/aliased-7named.txt`／`logs/default-7named.txt` |

⇒ 派单里那句"乙族被 307 治好了、甲族⛔ 动过一根手指"＝**本腿独立复认成立**（尺 1 直接给出"⛔ 动过"的对象层凭据，比颜色对拉强）。

## A.3 那它们的红落在哪一枚产码文件？——⛔ 落在任何产码文件

三个**产码事实**被测试误读；每一处都是**设计如此**，本腿⛔ 主张任何一处该改：

1. **`pathresolver.go:105 Resolve` 对"尚不存在的叶子"交回词汇形**（`#2` 的红件里那枚 `canonical` 自身还带着 `RISKAT~1\RUNNER~1`）⇒ 于是"Resolve 交回的东西恒一种比较形"这句**只对存在的路径成立**。
   ⚠ 派生结论（本腿⛔ 跑仪器，⛔ 判它是不是缺陷，只登记）：这条与 `blacklist.go:115-126 overrideApplies` 的窄规则合起来，意味着**在别名台面上，一条尚未创建的 B 级路径即便已有人批准过也⛔ 会放行**（`cand.raw != cand.real`）。方向＝**多问一次**＝票 72 的"拒绝可以宽、放行必须窄"那一侧，⛔ 是 fail-open。
2. **`pathresolver.go:345 formsOf` 产两种比较形**（`spellings()` 的注释 `:329-331` 逐字＝"The list is a union on purpose: matching an A anchor on any form can only ADD a deny, never remove one"）⇒ `#2` 拿"1"去断言它，本身就是把"这台机器只有一种拼法"写进了判据。
3. **`blacklist.go:115 overrideApplies` 只认句柄真实拼法**（`:122-124` `if cand.raw != cand.real { return false }`）⇒ `#3`/`#7` 拿短拼同时当入参与当 override key，是把这条**票 72 自己要求的窄规则**当成旧行为用。
   它的来⛔ 是臆造：票面 `:86-92`（编排者插单 10:08 的落地说明）＋票面 `:216-230` 插单原文＋票面 `:236-243` 逐条回执。
   ★同一枚 API 的**文档契约**在 `blacklist.go:58-60` 逐字＝`// Classify classifies an already-resolved canonical path. Callers MUST pass the path through Resolve first; lexically-spelled inputs are matched on a best-effort normalized form only (defence in depth, not a bypass route).` ⇒ **#3/#7 递进去的是"词汇形"，正是文档里那句 best-effort 那一侧**；拿它去领放行票＝按⛔ 契约的方式用这个 API。

### ★形状侧的正证据（本腿自己量的，⛔ 抄任何人的读数）

第三发尺＝把两枚 windows 文件里**全部 16 枚**用例一起放上别名台面（名册由 `git grep -h -oE '^func (Test[A-Za-z0-9]+)\(' HEAD -- <两文件>` 现抽，逐行打印对拉＝**16 枚**）：
件＝`logs/aliased-all16-bothfiles.txt`，读数＝**`rc=1`／PASS=9／FAIL=7／SKIP=0**。
FAIL 的 7 枚逐字＝本名册那 7 枚；PASS 的 9 枚里含**票 72 三条不变量用例本身**＝`TestClassifySpellingInvariance`／`TestOverrideOnlyAcceptsTheResolvedForm`／`TestClassifyFailClosedWhenSpellingUnprovable`。

⇒ **别名台面上那 7 枚就是全部红，⛔ 多一枚⛔ 少一枚**；而"比较写法⛔ 依赖 `t.TempDir()` 拼法"的 9 枚在同一个世界、同一份码上是**绿的**。
⇒ 这就是"判级⛔ 被证伪"最硬的一枚形状证据，⛔ 需要任何变异。（⛔ 拿它当"本格成立"：它只证"红在写法侧"，⛔ 证"改法有牙"。）

### 绿着的 9 枚有一个共同形状

`pathresolver_junction_windows_test.go` 十条红队用例（票面 `:13-15` 说它"Case 1..10"）：**红 4 枚（Case 2/3/4/8＝我的 #4/#5/#6/#7）、绿 6 枚（Case 1/5/6/7/9/10）**；`pathresolver_anchor_spelling_windows_test.go` 六枚里**红 3（#1/#2/#3）、绿 3**。
⚠ 这一列颜色是本腿**自己那一发 16 枚同跑**的读数（件 `logs/aliased-all16-bothfiles.txt`，别名台面），⛔ 是转述；⛔ 整包 200 枚的颜色⛔ 由本腿复现（那是编排者 `A834` 与 307 腿的两把整包读数）。

⇒ 落到"最小诚实修法"上，两枚文件里⛔ 需新造形状：绿着的这 9 枚就是现成本子。它们的共同形状（`sed` 逐枚读 `pathresolver_junction_windows_test.go:88-338` ＋ `anchor_spelling:74-117`/`:157-198`/`:250-292`）：**从不把 `res.Canonical` 与测试自己 Join 出来的字符串比**，只把 `res.Canonical` 喂给 `Classify`/`Gate`，比的是**类别**（或 canonical 比 canonical）：
- Case 1 `:103` `if got := Classify(res.Canonical); got != ClassA`（那两句 `t.Errorf` 打在 `:110`/`:112` 的是 `a04d3e2` 加的诊断输出，⛔ 判据）；
- Case 5 `:206`、Case 7 `:260`/`:263`、Case 10 `:334` 同形；
- Case 6 `:221`/`:224`/`:231` 只断 `Resolve` 的 `err`（`ErrReparseDenied` 的正反）；
- Case 9 `:314` `if r1.Canonical != r2.Canonical` —— **canonical vs canonical**，两侧同源，别名台面上⛔ 会错开。

⇒ 这一列就是"⛔ 放宽断言的改法长什么样"的**现成本子**：把比较的两端换成"同一把折法产出的东西"，⛔ 换掉判据本身。

## A.4 与票面不符的一处，具名报回（★本腿认为最值钱的一条）

票面 `:245-247`（AC#5 回执第 4 条）逐字＝
> `实测 grep -rn "\bGate(" --include=*.go internal/ cmd/ 的非测试命中只有 blacklist.go:76 它自己的定义 ⇒ 今天还没有生产 caller 传 override 表`

**这一句在 HEAD 上⛔ 成立。** 本腿现跑（尺＝`git grep -n -E '\bGate\(|\bClassify\(' HEAD -- internal cmd tools | grep -v _test.go`）＝非测试命中 **6 行**，其中一条**正是拿着 override 表的生产调用**：

```
HEAD:internal/tools/mode.go:91   overrides := b.confirmations()
HEAD:internal/tools/mode.go:93   c, err := b.paths.Canonicalize(p)
HEAD:internal/tools/mode.go:98   d := risk.Gate(c, overrides)
HEAD:internal/tools/bridge.go:103-104  Confirmations returns the B-tier single-file overrides on record (the set risk.Gate's bOverrides argument was designed to receive)
HEAD:internal/risk/rules_gateway.go:77 switch ctx.classifier.Classify(canonical)
HEAD:internal/tools/paths.go:381       switch risk.Classify(canonical)
```

⇒ 三处后果，逐条给下一程：
1. 票 72 当年**放弃反驳通道**（`:244` "不采用 — …'恒为 1'⛔ 是可依赖的前提"）的理由之一（"今天还⛔ 有生产 caller"）已过期；但它的**结论仍然对**——因为生产 caller 交进去的是 `Canonicalize` 的产物（`mode.go:93`），而 #3/#7 交进去的是 `filepath.Join(t.TempDir(), …)`。**"⛔ 恒为 1"与"生产按契约传已解析形"是两件事，票面把前者写成了后者的理由。**
2. ⇒ 甲族 #3/#7 的修法**⛔ 是零生产风险的纯测试侧改动**这一句⛔ 能再说：override 那条窄规则现在**有**一个生产消费者（`mode.go:98`，`Bridge.readBlacklist`，走 `internal/tools` 的桥面）。⇒ 任何去动 `overrideApplies` 的"修法"现在会**改变生产行为**（`BlacklistNote.AlreadyUnlocked` vs `Unlockable` 的分档，见 `mode.go` 的 `:34-38` 相对段），⛔ 再是"只动测试"。
3. ⇒ 台账要记一句（⛔ 本腿动票面原句）：票 72 `AC#5` 的回执第 4 条**取证已过期**，出处＝`mode.go:98` 这一枚 caller 是 `a04d3e2` 之后长出来的（本腿⛔ 量它长在哪一枚 commit，那一发要 `git log -S 'risk.Gate(' -- internal/tools/mode.go`，属下一程）。

## A.5 尺与件（本程 Go 面发数＝**4**，全部是具名 `-run ./internal/risk/`；其中**一枚作废**，具名如下）

```
1) -run '<7 枚名册>'  别名台面 → logs/aliased-7named.txt        rc=1  FAIL=7 PASS=0 SKIP=0
2) -run '<7 枚名册>'  默认台面 → logs/default-7named.txt        rc=0  FAIL=0 PASS=7 SKIP=0
3) -run '<16 枚名册>' 别名台面 → logs/aliased-all16-bothfiles.txt  ← 第一次那发⛔ 成立，见下
4) -run '<16 枚名册>' 别名台面（重搭）→ 同一枚件（覆盖前那发的正文已在本节记下）
```
**第 3 发作废并具名自报（本腿自己漏搭台面）**：那一发**根本⛔ 设任何 temp env**（本工具的 shell 状态跨命令⛔ 保留，上一发的 `export` 到这一发已经没了），于是 Go 从 Git Bash 继承来的 `TMPDIR` 那一形（POSIX 风格、无卷标）被 `os.TempDir()` 直接采信 ⇒ `TestClassifyAnchorSpellingIsNotVerdict` 报 `Resolve("\\TestClassify…")`、`TestClassifySpellingInvariance` **panic**：`slice bounds out of range [:1] with length 0`，崩点在 `pathresolver_anchor_spelling_windows_test.go:78` `share := strings.ToLower(vol[:1]) + "$"`（`vol := filepath.VolumeName(long)` 返回 `""` 时零守卫）。第 4 发改成三行独立 `export`（`env | grep -E '^(TMPDIR|TEMP|TMP)='` 现量三枚都是 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`）后 panic 消失、名册如 A.3。
★⚠ **顺手量到一枚⛔ 在甲族 7 枚射程内的既有脆面（本腿⛔ 修、⛔ 判它该⛔ 该修，只登记）**：`TestClassifySpellingInvariance`（⛔ 属那 7 枚，它在正确台面上是绿的）对 `filepath.VolumeName` 的返回值**零守卫**，任何让 `os.TempDir()` 交回无卷标路径的环境都能让它 **panic 掉整个测试二进制**、把同场其余用例一起带走（本腿第 3 发就是实锤：那一发 `--- PASS`=0、16 枚只跑到第 2 枚就崩）。⇒ 这条对本票族的用处＝**任何"拿这 7 枚当凭据"的发数都必须先断言名册行数**，否则会读到"0 枚过"这种假读数而⛔ 报错。

每发前一把 `git status --porcelain -- internal cmd`＝**0 行**（四发各量一次；第 2/3/4 发的读数直接打在 stdout 里：`0`／`0`／`0`）。⛔ `./...`、⛔ `cmd/wisp`、⛔ `d22scan`、⛔ 建树做变异。⛔ 翻任何框、⛔ push。

`51a0a0f2` numstat 原文（本腿逐行看过，共 30 行，`internal/` 面 2 行）：
```
17	3	internal/risk/syncdirs.go
141	0	internal/risk/syncdirs_alias_307_windows_test.go
```
