# 212-r3 修码证据件（票 212 第 3 轮小写码腿；三格欠账＋一枚前缀同步钉）

- 本程：写码腿 **212-r3**。派单来源＝非实现者验收腿 **212-v2** 件
  `docs/evidence/s1/212-comments-phantom-citation-v2.md`（33578 字节，本程现读）§8 表里点给 212-r3 的三笔具名欠账
  ＋编排者派单追加的一枚"两处前缀集合相等"钉。票面＝`.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md`
  （含 2026-10-03 终裁节与 2026-10-04 翻勾节；只读，本程零写）。
- 起手锚：`dev` HEAD＝`b1e59d63`。交件时 HEAD 已被并发腿推到 `2708289f`（256-r1／260-r2 落了四枚），
  本程的写面只有那五枚路径：`tools/d22scan/main.go`·`tools/d22scan/selftestsamples.go`·`tools/d22scan/scan_test.go`·
  `internal/agent/approval/pending_read.go`·`.scratch/wisp/probes/212/r3/**`。
- 一句话交付：**行为零变化**——(a)(c) 是纯措辞，(b) 是加样本＋把自述追平行为，钉是加测试；
  两道突变（摘豁免 m1／改前缀表 m2）都成对跑过并还原，md5 与 diff 哈希双双相等（§2／§5）。

## §0 起手锚与本程自跑的基线尺（先落档再动手）

起手 md5（本程现跑，档 `logs/md5-baseline.txt`）：

| 文件 | 起手 md5 |
|---|---|
| `tools/d22scan/main.go` | `e6d1745996e8e02c82ab691ac2cbadf8` |
| `tools/d22scan/selftestsamples.go` | `c61abbeba9b55a6ea452cec8478c2c93` |
| `internal/agent/approval/pending_read.go` | `e0c9dc3d5798514f7acba07984605680` |

基线读数（本程现跑，档 `logs/selftest-pristine.txt`／`logs/gate-pristine.txt`，rc 全 0）：

| 尺 | 起手逐字 |
|---|---|
| `cd tools/d22scan && go run . -self-test` 末行 | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` |
| 同上·名册行 | `d22scan -self-test: roster read from main.go = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation] + 1 finding type(s); 37 cases, 10 tag(s) covered, both directions required per tag` |
| `sh scripts/d22scan.sh` | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ＋ `d22scan: clean - no D22 ban violations; ... ban #8 internal/=500, ban #8 cmd/=97` |

起手 `git status` 里 `design/**` 的十六枚删除／`design/doubao/*` 的修改／`.gitignore` 的修改**在本程第一条命令之前就在盘上**
（本程第一条 git 命令的原始输出留档 `logs/working-tree-dirty-at-end.txt`，同一批路径全程未进本程任何 commit），
所以 §7 的越界检查按**单腿 commit 文件名册**算，不按全树 `git diff` 算——后者在共享树里必然带着别人的料。

## §1 (a) `pending_read.go` 的出处指错——改成实话

三处现读逐字（本程现跑，行号是本程取数时刻的，非派单给的）：

| 处 | 逐字 | 覆盖哪几枚字段 |
|---|---|---|
| `internal/tools/gate.go:12-13` | `// what the caller must show the user. The confirmation card (tickets 21/37)`／`// renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note), so` | **只有两枚**：`RulesHit`·`Reason` |
| `cmd/wisp/panel_pump.go:44-46` | `// The mapping is a copy, not a judgement: Level/RulesHit/Reason/`／`// SessionOverrideBlocked are the fields internal/risk put on the queue item`／`// when it was admitted, and panel renders them through the same` | **四枚**：`Level`·`RulesHit`·`Reason`·`SessionOverrideBlocked`，且措辞是"映射"不是"冻结契约" |
| `internal/panel/approval.go:81-85`（`CardViewFromDecision`，`:72` 起） | `Level: decision.Level.String(),`／`RulesHit: rules,`／`Reason: decision.Reason,`／`ReasonKnown: ...`／`SessionOverrideBlocked: decision.SessionOverrideBlocked,` | 四枚**确实都上视图**——这一枚本程现读，用来支撑"四枚之说为真、只是出处不是 `gate.go`" |

改前（`pending_read.go:41-46`，逐字）：

```
// observer sees it. Decision is the verdict as it was admitted (RulesHit,
// Reason, Level and SessionOverrideBlocked are the fields the L2 card shows
// verbatim, per ticket 17's frozen-contract note, which lives in
// internal/tools/gate.go:13 - the note is a display contract on Decision, and
// tools/d22scan does not enforce it (its pathresolver-bypass ban is about
// filepath.Clean/Abs outside risk.PathResolver), copied per
```

改后（`pending_read.go:40-54`，逐字首六行）：

```
// LiveApproval is one approval still waiting for an answer, as an in-process
// observer sees it. Decision is the verdict as it was admitted, and the two
// sentences about it are not the same sentence: the frozen-contract note in
// internal/tools/gate.go:12-13 (ticket 17) names exactly TWO fields the card
// renders VERBATIM - RulesHit and Reason - while the four the L2 card actually
// carries (Level/RulesHit/Reason/SessionOverrideBlocked) come from the panel's
```

归属写法：**两枚**（`RulesHit`/`Reason`）挂 `internal/tools/gate.go:12-13` 的票 17 冻结注；
**四枚一起**挂 `cmd/wisp/panel_pump.go:44-46` 的映射注，并明写那一注"不是第四份冻结契约"；
四枚上视图的那处实现 `internal/panel/approval.go:72-88` 一并点名；
`tools/d22scan does not enforce it` 那句保留（212-v2 §4 已判真），并补了 ban #9 也不执法（它只核路径存在性）。
路标一枚没拆（票面禁区），改的是"哪句话由哪一行支撑"。

**行为零变化的尺**（档 `logs/comment-only-proof.txt`）：
`git diff -U0` 里**非整行注释**的改动行数——`internal/agent/approval/pending_read.go`＝**0**、`tools/d22scan/main.go`＝**0**；
`--numstat` 分别是 `12 6` 与 `59 11`（删除列非零的每行都以 `//` 开头，是改写注释而非删逻辑）。

## §2 (b) 三形 silent 样本＋成对突变

三枚新样本（`tools/d22scan/selftestsamples.go:416-452`），一枚压一形；每形对应 `212-v2` §3 的探针：

| 样本 file（夹具内路径） | 压的形 | 注释里的 token（夹具内逐字） | 对应探针 |
|---|---|---|---|
| `internal/probe/cites-tail-ellipsis.go`（`:417-427`） | 完整名＋尾部 `…` | `docs/evidence/s1/probe-tail-ellipsis.md` + U+2026 | P3 |
| `internal/probe/cites-tail-star.go`（`:429-439`） | 完整名＋尾部 `*` | `docs/evidence/s1/probe-tail-star.md*` | P4 |
| `internal/probe/cites-cjk-attached.go`（`:441-452`） | CJK 粘连尾段 | `docs/evidence/s1/probe-swallowed.md` + U+89C1 + `docs/evidence/s1/...` | P5 |

粘连字符按本文件既有惯例用码点构串（`glyphCJKGlue = string(rune(0x89c1))`，`:32`），样本 token 起名用 `probe-` 前缀
而非 `174-`／`212-`——**故意不让它看起来像某枚真票的证据件**（那正是本票立的形）；三枚都不被夹具 seed，所以三枚都是"真②形状、被③豁免压住"。

成对突变 m1（档 `logs/m1b-*.txt`；改名之前的第一发 `logs/m1-*.txt` 只建不删、一并留档，那发的红句用的是旧 token 名）
＝把 `tools/d22scan/main.go:834` 的
`short := shorthandPathStarts(c.Text)` 逐字换成 `short := map[int]bool{} // MUTATION m1`，跑 `-self-test`，再逐字还原。

| 发 | 尺 | 逐字 | rc |
|---|---|---|---|
| 改前 | `md5sum main.go` | `442d46cadef8c641c54f7ac00984a470`（含本腿正常改动，非起手值） | — |
| 红发 | `go run . -self-test` | `d22scan -self-test: 4 direction(s) failed, 36/40 passed - the gate does not see what it claims` | **1** |
| 还原后 | `md5sum main.go` | `442d46cadef8c641c54f7ac00984a470`＝改前**全等** | — |
| 还原后 | `go run . -self-test` | `d22scan -self-test: clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)` | **0** |

红句逐字（摘豁免时四枚 silent FAIL 报出的 token；档 `logs/m1b-selftest-exemption-removed.txt:61,63,65,67`）：

```
internal/probe/cites-shorthand.go:4: [phantom-citation] comment cites repo path "docs/evidence/s1/152-...-accept-r2.md"
internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/212-"
internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/152"
internal/probe/cites-tail-ellipsis.go:4: [phantom-citation] comment cites repo path "docs/evidence/s1/probe-tail-ellipsis.md"
internal/probe/cites-tail-star.go:4: [phantom-citation] comment cites repo path "docs/evidence/s1/probe-tail-star.md"
internal/probe/cites-cjk-attached.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/probe-swallowed.md见docs/evidence/s1"
```

读法：三枚**新**样本各响一次（第 4/5/6 行），r2 那枚样本照旧响三次——`4 direction(s) failed`＝四个**用例**（其中三枚是本腿新加的），
⇒ 三形确实**只因豁免而静默**，样本压的是豁免本身，不是别的什么。第 6 行那一枚报出的 token 里带着 `见`，
逐字复现了 `212-v2` §3 P5 的"粘连把前面的完整幻影一起吃掉"。

`git diff -- tools/d22scan` 在还原后**不等于空**——本腿的正常交付就改着那三枚文件；
所以"突变零残留"的尺用两把：`md5sum main.go` 前后全等（上表），以及 §5 那发的 `git diff -- tools/d22scan internal/agent/approval/pending_read.go | md5sum` 前后全等。

## §3 (b) 名册行与 `main.go` 自述改成真实射程

三处，每处给改前／改后要点（全文见 `git show` 本腿产码枚）：

| 处 | 改前 | 改后 |
|---|---|---|
| `main.go:50-69`（`// Bans` 第 9 枚名册行） | 末两句只有 `Judged against the scanned root; shorthand forms are a separate prescription (AC#3), not a violation.` | 明写射程两刀：①**扫描面**＝`internal/` 与 `cmd/` 的产码 Go 文件，`tools/**` 与 `_test.go` 在射程外（仪器自己就能躺幻影路标，212-v2 §8 #4，扩面＝批准面）；②**豁免面**＝按 region **起点**判，所以三形免罪＝标记**在路径中段**／标记** trailing 一枚完整路径**／**CJK 粘连尾段**；三形各有一枚 expect-silent 样本钉着；并保留反句"A fully-spelled path with no mark in it still rings" |
| `main.go:805-830`（`scanGoFile` 里 ban #9 的自述段） | `③ has two measured shapes, and both are exempt here: prose with no path token at all, and a path spelled WITH an abbreviation mark INSIDE it` | 改成 `What that actually exempts, stated as the range the code has`＋**三形编号列全**（(1) 中段标记／(2) 尾部 `…` 或 `*`／(3) CJK 粘连），并写明 (2)(3) 是裁 ⓐ 的代价、今日活体 0 枚、编排者因此不收窄（票面翻勾节，口令「212 收窄豁免」）、样本落点与"再收窄就会以红方向检查出现而不是靠重读这段注释" |
| `main.go:908-931`（`shorthandPathStarts` 自述） | 只写了"尾部标记也是③"，且说豁免由"the expect-silent sample"（单数）钉 | 补第三臂（CJK 粘连＝一枚完整名＋一个 CJK 字母＋一段缩写、中间无缝，整串进 region）＋`212-v2` P6/P8/P9 划出的三条边界（无标记仍响／换行不外溢／空格截断）＋"每一臂各一枚样本，窄化会响在方向检查上"＋活体 0 枚 |

派单里那句 ⛔ 的落实：全腿**没有**留下任何读起来比仪器实际更窄的自述；也没有把"三形"写成"省略号中段一种"。

## §4 (c) 仪器自己注释里的那枚幻影

现读确认（`ls docs/evidence/s1/ | grep 212`）：盘上只有 `212-comments-phantom-citation-v2.md` 一枚，
`212-citation-ruler.md` **从未存在**。旧句在 `tools/d22scan/selftestsamples.go` 的 ring 样本（`:373-382`）里，本程改成：

```
"// Reading the readings first is mandatory; see docs/evidence/s1/212-comments-phantom-citation-v2.md\n" +
```

`summary` 与 `note` 同批改齐（note 里逐字记下：旧名在真仓与夹具里都不存在＝本票的形报在仪器身上，212-v2 §8 #4；
改指一枚**真仓存在、夹具不 seed** 的件，所以 ring 向照旧响；⛔ 修法是改措辞，不是把 `tools/**` 拉进射程）。

ring 向未被削弱（本程现跑）：`logs/selftest-final.txt` 那一行逐字＝
`d22scan -self-test: ban #9 phantom-citation    ring   OK     a comment citing docs/evidence/s1/212-comments-phantom-citation-v2.md, which the fixture does not seed | ranged on 1 finding(s): internal/probe/cites.go:3: [phantom-citation] comment cites repo path "docs/evidence/s1/212-comments-phantom-citation-v2.md" ...`

**本程顺手量到的一格，具名报给编排者（不属本腿改动面）**：ban #9 读的是 `f.Comments`（`main.go:831-832`），
**只有注释**——字符串字面量里的路径它永远看不见。所以：
① 旧那枚 `212-citation-ruler.md` 即便把 `tools/**` 纳入射程**也不会被报红**（它在 `src:` 串里），
"扩射程"并不能自动修掉这一格，只有改措辞能——这直接影响编排者那道待机的成色；
② 本腿新样本里的 `docs/evidence/s1/probe-*.md` 也全在字符串字面量内，同样不会被扫到（现跑复认：见下条尺）。

粗分母（档 `logs/phantom-head-side.txt`／`logs/phantom-now-side.txt`，尺＝`tools-cited.sh`：只取注释行、抽
`docs|.scratch|internal|cmd|tools|scripts/` 形状、跳过含 `...`/`…`/`*` 的③形、对 `[ -e ]` 比）：
起手锚的 `tools/` 注释侧 7 枚不存在 token，本腿之后**逐名相同（diff 为空）**⇒ **本程没往仪器源码里种下任何新幻影路标**；
那 7 枚是 `internal/build`·`internal/build/leak.go`·`internal/proc.WithRegistry`·`internal/tool`·
`internal/tools.Result.AppliedSteps`·`internal/tools/tmp`·`tools/ok.go`，其中两枚会被 `symRefRe` 排除（带 `.大写` 的 API 引用），
⇒ 若真扩射程，今天会红的大约 5 枚。**这只是粗分母，不是裁决**：扩射程＝人工批准，本程一字未动射程。

⛔ 未做且不许做：`tools/**` 未进 ban #9 射程；`tools/d22scan/allowlist.txt` 未碰；任何 ban 的扫描范围未改。

## §5 前缀同步钉（`repoPathRe` ↔ `shorthandRegionRe`）

落点：`tools/d22scan/scan_test.go:2384-2465`，`TestBan9PrefixTablesAreTheSameList` ＋ helper `regexPrefixAlternation`。
两枚正则的注释各加一句指回这枚钉（`main.go:886-889` 属 `repoPathRe`、`main.go:901-903` 属 `shorthandRegionRe`；
正则字面量各在 `:891` 与 `:905`）。

选"测试向"而不是"抽公共表"的理由（写清行为影响）：抽成一处再编译**行为等价但形状搬家**——两枚判据的正则字面量会从源码里消失，
而名册行／自述／样本都按"源码里读得到形状"来写；本腿的约束是行为与形状都不许动，所以装**断言**不装重构。
比较用**集合**不用字符串序：两枚正则里每个候选后面都紧跟一个字面 `/`，列序改不了匹配结果，把重排当红是噪声。

正控 m2（成对；档 `logs/m2b-*.txt`；同形的第一发 `logs/m2-*.txt` 只建不删、一并留档）＝只给 `main.go:891` 的
`repoPathRe` 前缀组加一枚 `plugins|`。⚠ 这一发的行号是**终态树的 891**；第一次跑时它还在 890（本腿后来的措辞修正把行推下去了），
照 890 复跑会命中 `var` 行、** mutation 静默不落地**——本程真踩过一次（`m2b` 第一发 rc=0 才发现尺指错行），具名记在这里。

| 发 | 尺 | 逐字 | rc |
|---|---|---|---|
| 改前 | `md5sum main.go` | `442d46cadef8c641c54f7ac00984a470` | — |
| 改前 | `git diff -- tools/d22scan internal/agent/approval/pending_read.go \| md5sum` | `6748f19f3239df647828c1cdfb053260` | — |
| 红发 | `go test -count=1 -v -run 'TestBan9PrefixTablesAreTheSameList\|TestScannerSelfScanOfRealRepoIsGreen' ./...` | 见下 | **1** |
| 红发·同棵树 | `go run . -self-test` | `d22scan -self-test: clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)` | **0** |
| 还原后 | `md5sum main.go` | `442d46cadef8c641c54f7ac00984a470`＝改前**全等** | — |
| 还原后 | 同上 `git diff \| md5sum` | `6748f19f3239df647828c1cdfb053260`＝改前**全等** | — |
| 还原后 | 同一枚测试单跑 | `--- PASS: TestBan9PrefixTablesAreTheSameList (0.00s)` | **0** |

红句逐字（`logs/m2b-test-prefix-divergence.txt`，指向 `scan_test.go:2441`）：

```
    scan_test.go:2441: ban #9's two hand-copied prefix tables disagree: repoPathRe=[plugins docs \.scratch internal cmd tools scripts], shorthandRegionRe=[docs \.scratch internal cmd tools scripts]; only in repoPathRe=[plugins] means a shorthand path there is NOT exempt and class 3 convicts again (ticket 212 AC#3), only in shorthandRegionRe=[] means a COMPLETE phantom path there gets swallowed by the region, which widens the exemption without an owner's approval. Edit both lists in one commit, and read the ban #9 comment in main.go before deciding which shape is intended.
--- FAIL: TestBan9PrefixTablesAreTheSameList (0.00s)
```

同一发里 `TestScannerSelfScanOfRealRepoIsGreen`＝`--- PASS`、`-self-test` 仍 40/40 全绿
⇒ **没有这枚钉，前缀漂移在盘上是完全不可见的**（这正是 212-v2 §8 #3 说的"AC#3 那一形会在新前缀上悄悄复发"）。

## §6 撞钉预检：哪些断言里硬写了数字

起手那发（逐字命令 `grep -rn '37\|expect-silent\|direction checks' tools/d22scan --include=*.go`）命中 13 行，逐枚分类：

| 命中 | 是什么 | 是断言里的数字吗 |
|---|---|---|
| `selftest.go:606` | 末行的格式串 `all %d direction checks passed (%d expect-ring, %d expect-silent)`，参数是 `len(selfCases)`／`countWant(...)` | **不是**，分母现算 |
| `selftest.go:568`（未命中 37 但同族） | 名册行 `%d cases`＝`len(selfCases)` | 不是 |
| `selftest.go:332` | 单向 tag 的 HOLE 文案（字符串） | 不是数字 |
| `selftest.go:363` | 散文注释 | 不是 |
| `selftest_test.go:101` | `strings.Contains(joined, "only an expect-silent sample")` | 文案子串，无数字 |
| `main.go:35`·`:66`·`:142` | `37 / 40 / 43` 是 ban #6 frontend/ 的历史分母叙述（账 `A207`） | 散文，非断言 |
| `main.go:280`·`:551` | "the whitelist drops 5 of 37 tracked files" 的散文 | 散文，非断言 |
| `gitignore.go:21` | 同族 37/40/43 散文 | 散文 |
| `scan_test.go:1036`·`:1098`·`:1272` | ban #6/#8 分母的散文注释 | 散文，非断言 |
| `selftestsamples.go:112`·`:304`·`:337` | 样本 `note` 文案 | 文案 |

⇒ **结论：`tools/d22scan` 里没有任何断言硬写用例分母**（37 只出现在散文与文案里），
所以本腿**既没有需要同步的既有断言、也没有动过任何既有断言**——这条不是"没找到"，是 13 枚逐枚读完的判语。
分母从 37 变 40 只落在两处**现算输出**上（末行与名册行），是算术结果不是被改的数。

另外两处带数字的**外部散文**本程判不动、具名归口：`docs/reports/pending-and-issues.md:11408`
（记着 `all 37 direction checks passed (20 expect-ring, 17 expect-silent)`）与 `:7832`（更早的 34＝19/15）——
台账是**只追加不删**的真相源，且**不在本程写面**（派单⛔）；分母现已是 40＝20/20，
**归编排者在台账新记一条时带上口径与本锚点**，不是回头改那两行。

## §7 门禁读数终态（AC#5：八枚 scope 行逐名对比）

本程交件前的最后一发（档 `logs/gate-final.txt`，rc=**0**）：

```
runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 228 production Go files
d22scan: scope bans #1-5 cmd/           examined  38 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 501 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  98 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=501, ban #8 cmd/=98; ban #8 emoji coverage: design/ 39 text files; frontend/ 85 text files; internal/ 501 Go files, comments and _test.go included; cmd/ 98 Go files, comments and _test.go included
```

与起手 `logs/gate-pristine.txt` 的八枚 scope 行对比（`diff logs/scope-pristine.txt logs/scope-final.txt`，档留两枚文件）：

| scope | 起手 | 终值 | 判 |
|---|---|---|---|
| #1-5 internal/ | 228 | 228 | 逐字不变 |
| #1-5 cmd/ | 38 | 38 | 逐字不变 |
| #6 frontend/ | 85 | 85 | 逐字不变 |
| #7 internal/tools/ | 23 | 23 | 逐字不变 |
| #8 design/ | 39 | 39 | 逐字不变 |
| #8 frontend/ | 85 | 85 | 逐字不变 |
| #8 internal/ | 500 | **501** | **疑似争用**：`internal/ball/hotkey_live_test.go` mtime 10:25（本程起手发 10:05 之后），写面属并发腿 **260-r1/r2**；ban #1-5 的产码分母 228 未动＝它是 `_test.go`，只有 ban #8 含 |
| #8 cmd/ | 97 | **98** | **疑似争用**：`cmd/wisp/resident_approval_risk_256_windows_test.go` 随并发腿 **256-r1 的 `2a10134e`**（10:20）进树；尺＝`git diff --name-status b1e59d63..HEAD -- internal cmd`（档 `logs/concurrent-leg-file-changes.txt`）只有那三行，本程一枚未碰 |

⇒ **八枚里 6 枚逐字不变，2 枚各 +1 且都指得出"哪枚并发腿的哪一行文件"**；本程产码面只碰
`tools/d22scan/**`（不在任何 ban 的扫描范围）与 `internal/agent/approval/pending_read.go` 的**注释行**，
门禁 rc=0、零新增红，AC#5 的"逐名不变"以本程起手为基线成立。
`go vet ./...`（模块内）无输出 rc=0；`go build ./...`（根模块）rc=0；
`gofmt -l` 在 `tools/d22scan` 只报 `selftest.go`——本程未碰那枚文件，`git show HEAD:tools/d22scan/selftest.go` 落盘同样被报
（起手既有状态，具名不顺手修）；`gofumpt -l` 对三枚 tools 文件全部零输出。
`gofumpt -l internal/agent/approval/pending_read.go` 报该文件：同包**未碰过的** `queue.go` 同样被报，
尺＝`gofmt -d` 显示整文件逐行替换＝工作树 CRLF 状态（`.gitattributes` 对 `*.go` 写 `eol=lf`，
所以进树的 blob 仍是 LF：`git show HEAD:internal/agent/approval/pending_read.go` 落盘后 `gofmt -l` 零输出）。
本程**未**去"修行尾"——那会造出一枚 131 行的整文件 diff，不属本票三格。

## §8 判不动的地方（具名，不用"应该没问题"填空）

| 枚 | 内容 | 为什么本程动不了 |
|---|---|---|
| 1 | 把 `tools/**` 纳入 ban #9 射程 | 扩射程＝人工批准面（派单⛔⛔ 逐字；票面翻勾节记着"与票 264 同批摆"）。本程只把 §4 那句"即便扩了射程也扫不到字符串字面量"量出来供裁 |
| 2 | 把 `_test.go` 纳入 ban #9 射程 | 同上；唯一活体缩写形 `cmd/wisp/slo_report_144_windows_test.go:851` 仍在射程外（212-v2 §8 #6 归口不变） |
| 3 | 收窄豁免（让尾部标记／CJK 粘连真红） | 编排者 2026-10-04 已裁**不收窄**，口令「212 收窄豁免」；本程做的正是派单指定的替代动作＝让真实射程可见而不改行为 |
| 4 | 行号型路标有没有人管（212-v2 §8 #8 两枚错锚） | ban #9 只钉存在性；要不要立第十枚归编排者开票，本程不裁 |
| 5 | NTFS 大小写不敏感的 stat 残洞（P10 静默／P11 响） | 已知残余，两任判语与 v2 同判"量得到、无猎物"；本程未复跑该探针（不在三格内） |
| 6 | 台账里 `37 direction checks` 两条散文过期（`:11408`／`:7832`） | 台账不在写面（派单⛔），且规矩是只追加不删；§6 已具名归口编排者 |
| 7 | `tools/d22scan/selftest.go` 起手即被 `gofmt -l` 报 | 本程未碰那枚文件；"顺手格式化"会把别人的历史混进本腿 diff（票面 AC#3 越界检查的精神） |
| 8 | 工作树里 `internal/agent/approval/*.go` 的 CRLF 状态 | 不是本票的格；进树 blob 已是 LF，改了只会造整文件噪声。判语＝现状无害，具名登记 |
| 9 | 票面 AC 框／`-done` 改名 | AC 框只归编排者；票面翻勾节写死"212-r3 交完后须再来一枚非实现者小程 212-v3 核三笔"，本程不追认自己的判据 |
| 10 | `pending_read.go` 那句"四枚字段"要不要连带把 `panel_pump.go` 的注也升级成契约 | 那是**契约变更**（D43／C 族口径），本程只把出处写对，一字未动被引用方 |
