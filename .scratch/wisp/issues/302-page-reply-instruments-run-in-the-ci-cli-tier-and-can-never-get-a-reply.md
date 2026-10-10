# 票 302 — 三枚**要页面自己回话**的判据躺在 `--scope=cli` 档里被托管 runner 求值：那台机器永远给不出回执，于是它们在 CI 上恒红，而"没有回执"恰恰就是那三枚用例的**判据本身**

**立票**：2026-10-10 16:0x 编排者（来路＝票 301 `AC#2` 的 CI 色回填，正文＝`.scratch/wisp/probes/301/orch/2026-10-10-ci-color-backfill.md` §4；**这 3 枚⛔ 是票 301 造成的**，只是随 742 枚一次推送第一次进了 CI 才被看见）
**性质**：★**归口／放置票（instrument placement），⛔ 缺陷票**。那三枚用例写得**诚实**——它们的红句逐字写着"没有回执就⛔ 判"，正是本仓 10-08 以来反复立的那条规矩（"判据必须页面自己报回"，票 33 `A502` 形）。坏的是**它们被放进了一个永远拿不到回执的档**，于是 CI 上有一枚常红的灯，而常红的灯会让人开始忽略红色。
**为什么要紧**：`test-windows` 今天六枚 job 里就是红的（改前基线也红）。里面混着三种完全不同的红：① 门的守卫红（⛔ 该出现）、② `internal/risk` 那 12 枚环境敏感红（另一族，见 `AC#4` 边界）、③ 本票这 3 枚"真窗判据被 CI 求值"红。**枚数混在一起＝下一程分不清哪一枚该救**。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**三枚名册**（尺＝两发同一 job 的 `--- FAIL` 名册作差 `comm -13`，原始字节＝`probes/301/orch/logs/ci-baseline-cli-fail.txt`／`ci-after-cli-fail.txt`）：
  `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`。
- ★**同档基线↔改后的四数**（尺＝步日志里 `portable-tests.sh: four numbers` 那一行逐字引）：`RUN=335 PASS=231 FAIL=6 SKIP=1` → `RUN=398 PASS=278 FAIL=8 SKIP=2`。另有 1 枚基线红→改后 `--- SKIP`＝`TestPanelHostLatencyPercentilesAC2`。⚠ **⛔ 把这 +2 读成"本批退步"**：`FAIL` 6→8 是名集合差（＋3／−1）后的净值，两把尺⛔ 同物。
- ★**红句形状三条同形**（逐字，尺＝每枚取 `=== RUN` ↔ `--- FAIL` 之间那段；⚠ Go 的 `-v` 把明细写在 `--- FAIL:` **之前**，取反方向的那把尺会读到空）：
  `no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all)`（`cmd/wisp/panel_resident_windows_test.go:325`）／同形两枚在 `:821`（`"ac14r-0"`）、`:866`（`"ac14-push"`），第三枚出处件另含 `cmd/wisp/panel_pageover_33r10_windows_test.go`。
- ★**窗体那一侧在同发里是**建成**的**（同一段日志逐字：`AC#13 33-r10 world: embed resolves 1068 entry byte(s), 1 element id(s) [root]`、`entry bytes carried=true`、`probe document is 136 byte(s)`）⇒ **"CI 开不出面板"那句在本发⛔ 成立**；缺的是"页面 → Go"那一跳的回执。⚠ 那到底是因为 runner 没有 WebView2 Runtime、还是有 Runtime 而没人在泵／没桌面会话——**两形本票⛔ 混、⛔ 现在判**，那一枚读数归在飞的 `33-n1`（票 33），本票 `AC#0` ② 只负责把它**接过来答**。
- ★**本机同码不红**（尺＝票 33 探针腿 `33-p1` 的交件 `probes/33/p1/**`：`Eval` 推＋页面回话双向闭环 145~187 ms）⇒ 这是**环境差**，⛔ 是断言写坏。
- ★**既有机制两枚可用**（⛔ 新造）：① `winlive` tag 档＋`ci.yml` 里那枚 `go vet -tags winlive` 编译门（票 111 `AC#11` 落的那道）；② `scripts/portable-tests.sh` 的夹具 ledger（`-skip` 名单＋逐行理由）。⚠ 引 ledger 前必读 `A811`/`A812` 那两条已证事实：**`-skip` 是静默过滤器、⛔ 产 `--- SKIP`**（尺＝两发归档 CI 字节），而 `runtests.sh:98`→`:102` 把任何 `^--- SKIP` 判红 ⇒ 走②买到的是"这枚⛔ 再被求值"，而**是⛔ 是**"它被记成跳过"，那句话要写进文案。
- ⚠**边界一枚**：本票⛔ 管 `internal/risk` 那 12 枚（5 枚红句逐字点名 runner TEMP 的 8.3 别名 `C:\Users\RUNNER~1`、2 枚同文件 A/B 表判定但红句不含该串、5 枚 `syncdirs_test.go` 的 under-profile fallback 家族）。那一族的先例在**票 115**（`AC#2`／`AC#3` 做过"复现 runner 8.3 形状、改前红／改后绿"的 tree attribution），**归口记在台账 `A817`**、⛔ 塞进本票——两族的修法一个在"档"、一个在"判据读的是哪棵树"，混一票必然顺手改错那半。

## 要建什么（`AC#0` 之前⛔ 任何产码；三形都必带代价，⛔ 只摆两形）

- [ ] **`AC#0` 只读代价普查（这一格⛔ 任何产码）**：交回三张表——① **全名册**：当前三档（`cli`／`windows`／`core`）里，哪些用例的判据**必须**要真窗或页面回执（尺要写清是"红句里那句 `no report … from the page`"还是"`WISP_*` 环境门"还是 `winlive` tag，⛔ 三者混一把）；⚠ 枚数必须写清是**抽样还是整族**，射程目录逐枚带。② **死因两形分开答**：runner 上"没有 WebView2 Runtime" ↔ "有 Runtime 但拿不到回执"各自盘上证据到不到（接 `33-n1` 的读数，⛔ 自造一把）。③ **三形代价表**：甲＝把这 N 枚搬进 `winlive` tag 档；乙＝进 ledger 的 `-skip` 名单并写实话；丙＝**不动**，把"CI 上这 N 枚恒红"登记成具名已知红（要写清：恒红的灯会被后来人当成噪声，而本票的起因**正是**噪声让人分不清）。＋**必答**：丙形会不会让下一程把这三枚和 `internal/risk` 那 12 枚读成同一件事。判据＝每条带"哪把尺＋射程目录＋blob 还是工作树"，⛔ 裸数。
- [ ] **`AC#1` 落地形（只在 `AC#0` 交出代价、并经编排者裁"甲／乙／丙"之后才开工）**：**硬约束三条**——① ⛔ 放宽任何断言（"没有回执就⛔ 判"那句必须原样保住：任何一形都⛔ 许把"无回执"改成 pass 或改成 `t.Skip` 之前先 pass）；② 只许动**档／tag／ledger 文案**，⛔ 动 `internal/**`、⛔ 动 `cmd/wisp/**` 的产码（判据本体⛔ 改）；③ 若动 `scripts/portable-tests.sh` 的清单，则**同一笔 commit** 里连 pin 一起改（GUARD D 的退码文本逐字要求，先例＝票 301 `AC#1`；引用那两句时注意它**跨两行**）。⛔ 顺手改别的档、⛔ 顺手把 `internal/risk` 那 12 枚一起处理。
- [ ] **`AC#2` 反恒真那一格**：交付含**成对两发**＝(a) 把那枚豁免／换档撤回去 ⇒ 目标档当场按预期红（具名红句逐字引，⛔ 只报枚数）；(b) `cmp` 逐字节还原 ⇒ 复绿。⚠ 若裁的是甲形（搬进 `winlive`），必须另附一发证明**那三枚仍进编译面**（尺＝`go vet -tags winlive ./cmd/wisp/` rc=0 且 `-list` 名册含那三枚）——⛔ 让"进 CI 的分母"变成"既⛔ 跑也⛔ 编"，那是把一枚假绿换成一枚隐形。
- [ ] **`AC#3` 门禁与越界**：`sh scripts/d22scan.sh` rc=0；`bash -n scripts/portable-tests.sh` rc=0（若动它）；`bash scripts/portable-tests.sh --scope=census` 的 totals 行**逐字未变**（现量基线＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`）；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺——`301-v1` §3-⑥ 已证区间尺会把别人的笔算到腿头上）；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec 且写在 `$( … )` **之外**。
- [ ] **`AC#4` 边界声明（本格⛔ 产码、⛔ 改判据）**：在交付里**逐枚写明**本票收掉了哪几枚、⛔ 收哪些，并把 `internal/risk` 那 12 枚**逐枚具名**留回台账那条 `A817`。判据＝名册差集（尺＝`comm` 两边各自的红名册，⛔ 比枚数）＋一句"哪几枚从此⛔ 在 CI 恒红、哪几枚仍红"。⚠ 若腿把那一族"顺手修了"＝越界，停手上报（本票与它是两枚射程）。

## 规矩（本票全程）

- 子代理**只 commit、不 push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree；临时件**只建不删**；证据件⛔ 叫 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`）。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件；⛔ 为变绿放宽任何断言（本票的第一风险恰恰是"为了 CI 变绿而把判据改软"）。
- 裁决者≠实现者（`SPEC-12 §4.3`）。本票的 `AC#1`/`AC#2` 判语⛔ 由落地腿自勾。
- ⚠ 三形里"丙＝不做"是**合法选项**，⛔ 把它当失败；它的代价是具名的（一枚常红灯继续被忽略），选它要把那句写进台账。

next=票 301 `AC#4`/`AC#5` 的验收腿 `301-v3`（在飞）→ `302-a1`（`AC#0` 只读）→ 编排者裁形 → `302-r1` → `302-v1`
