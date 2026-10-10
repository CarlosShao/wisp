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

- [x] **`AC#0` 只读代价普查（这一格⛔ 任何产码）**：交回三张表——① **全名册**：当前三档（`cli`／`windows`／`core`）里，哪些用例的判据**必须**要真窗或页面回执（尺要写清是"红句里那句 `no report … from the page`"还是"`WISP_*` 环境门"还是 `winlive` tag，⛔ 三者混一把）；⚠ 枚数必须写清是**抽样还是整族**，射程目录逐枚带。② **死因两形分开答**：runner 上"没有 WebView2 Runtime" ↔ "有 Runtime 但拿不到回执"各自盘上证据到不到（接 `33-n1` 的读数，⛔ 自造一把）。③ **三形代价表**：甲＝把这 N 枚搬进 `winlive` tag 档；乙＝进 ledger 的 `-skip` 名单并写实话；丙＝**不动**，把"CI 上这 N 枚恒红"登记成具名已知红（要写清：恒红的灯会被后来人当成噪声，而本票的起因**正是**噪声让人分不清）。＋**必答**：丙形会不会让下一程把这三枚和 `internal/risk` 那 12 枚读成同一件事。判据＝每条带"哪把尺＋射程目录＋blob 还是工作树"，⛔ 裸数。
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

## 现量补充（2026-10-10 16:4x 编排者，来路＝`302-a1` 交件 `probes/302/a1/40-census.md` ＋ 我自己复跑的七把我没让腿跑的尺）

⚠⚠**本节⛔ 改写本票票面 `:1`／`:4`／`:14`／`:15` 四句的前提**（原句一字不改、留在上面当快照；本节＝现行读数）。读完这四句再派任何腿。

### ① 标题与 `:1`／`:4` 那句「那台机器永远给不出回执」「CI 上恒红」**⛔ 成立**

尺＝两发归档 CI 字节里同名用例的 `^--- [A-Z]*: <名>` 行逐枚打（件＝`probes/301/orch/logs/ci-baseline-cli-block.txt`／`ci-after-cli-block.txt`；两发 Runner Image 与 agent 版本我 16:0x 已逐字对拉＝`windows-2025-vs2026`／`20260925.250.1`，同⇒"镜像漂"那一支在这枚对拉里⛔ 有依据）：

| 用例 | 基线（`cc315261`，10-06） | 改后（`bcd0a543`，10-10） |
|---|---|---|
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | `--- SKIP`（`:1270` 具名理由逐字 `AC#13 has no subject in this tree: the embed resolves no entry`） | `--- FAIL (20.04s)` |
| `TestAC14AwaitedBindingReplyReachesThePage` | **`--- PASS (1.25s)`**，且逐字带页面自己的话（`:1293`：`page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`） | `--- FAIL (20.04s)` |
| `TestAC14GoSideEvalPushReachesThePage` | **`--- PASS (0.97s)`**（`:1300`：`page's own words: title="PUSHED-33R5-OK"`） | `--- FAIL (20.03s)` |
| `TestPanelHostRealWindowHopAndLifecycle` | `--- FAIL (3.01s)` 红于**预算**（`panel_host_windows_test.go:665`：`cold bring-up 2889.2 ms exceeds D32 panel cold budget 1500 ms`） | `--- FAIL (5.02s)` 红于**回执为零**（`:662`：`did not produce a browser round trip (got -1.000)`） |
| `TestPanelHostLatencyPercentilesAC2` | `--- FAIL (0.00s)` | `--- SKIP`（×2） |

- 同一基线发里还有 `:1201` 逐字 `our tree webview=7 … machine-wide msedgewebview2=7 … same HWND 0xc014c across hide->re-show=true` ⇒ **托管 runner 上 WebView2 起得来、窗建得起、页面回过话**。⇒ 腿的**形①（runner 没有 Runtime）＝已否证**，任何文案⛔ 许写那句。
- ⇒ 本票⛔ 是"一枚常红灯 + 一档放错"那一件事。**盘上是两件事**：归口／放置 ＋ **一枚未归因回归**。⚠ 我 16:0x 在证据件 §4 写的"三枚都是新可见、⛔ 是新坏"**只对 1 枚成立**（`TestAC13ColdStart…` 那一枚确实是从具名跳过走进 FAIL；另两枚基线是**绿**的）。**记我。**

### ② ★★**本机同码今天也红** ⇒ 「环境差」整条读法作废（尺＝编排者现跑，件＝`probes/302/orch/g4-local-three-cases.txt`）

`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -timeout 420s -v -run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|TestAC14AwaitedBindingReplyReachesThePage|TestAC14GoSideEvalPushReachesThePage' ./cmd/wisp/`
⇒ **`G4_rc=1`／三枚全 `--- FAIL (20.01s)`／红句逐字同 CI**（`no report "ac13-probe" … (what DID arrive at the door: nothing at all)`、`"ac14r-0"`、`"ac14-push"`）。
- 射程口径：跑的是**工作树**；`git status --porcelain -- cmd internal scripts .github docs` 现量 **0 行** ⇒ 被测码＝HEAD 等价（唯一脏件在我自家的 `.scratch/**` 文档面）。`tasklist` 起手 `wisp.exe=0`／`balldebug.exe=0`。
- ⇒ **票面 `:15` 那句"本机同码⛔ 红 ⇒ 这是环境差"作废**（它引的是 `33-p1`，那是**旧码**的读数，⛔ 覆盖当前 HEAD）。
- ⇒ **性质再往实现层挪一格**：这是**一枚可本机复现的页面→Go 回执断口**，⛔ CI 专属、⛔ 靠搬档能修。归因与修复＝**票 303**（本票只留"放置／归口"那一半）。
- 〔读码推的，⛔ 实测〕用户可见形状候选：这道门（票 35 的 `installPanelTransport`／消息钩子）是面板与宿主共用的，所以**"面板上点了没反应"是本票这一族最可能的产品形状**；把这句判死需要票 303 里那一发真机读数（⛔ 现在当已知）。

### ③ 名册＝**5 枚（＋1 枚登记）**，⛔ 3 枚

票面 `:9` 那三枚逐枚确认；同射程宽一把词面尺另抓到 `TestPanelHostRealWindowHopAndLifecycle`（要真 round trip，两发都红、**换因**，见①表）与派生的 `TestPanelHostLatencyPercentilesAC2`（⛔ 红而是 `--- SKIP` ⇒ 撞 `tools/d22scan/runtests.sh:98`→`:102` 那枚"SKIP ⛔ 不算过"的门）；相邻不同族一枚 `TestAC247LiveMicrophoneLevelsReachTheBallSeam`（⛔ 回执，⛔ 本票射程）。⇒ 凡报枚数处已写明"哪把尺＋射程"。

### ④ 三形代价里**由我现量补上的三枚读数**（腿⛔ go 编译面，这三把我替它取）

- **乙形（进 ledger `-skip`）今天⛔ 可行**——这条⛔ 是腿从 blob 推的，我量到了实物：`go test -list '.*' ./cmd/wisp/`（PATH **未**铺 `build/`，本机台面）⇒ **`exit status 0xc0000135`、`G3_rc=1`、列名 **0** 枚**（件 `probes/302/orch/g3-list-nopath.txt`／`…-body.txt`；`dll-visible=build/onnxruntime.dll …` 在目录里、只是⛔ 在 PATH 上）。而 ledger 的**包**会进 `scripts/portable-tests.sh:617` 那行全域 universe（⛔ 铺 PATH、⛔ 分档）⇒ 加 3 行＝把 `--scope=windows` 那一步今天的绿数拖成 rc=1。**要用乙，必须先动 `:617`（`scripts/**` 产码 ⇒ 触发 `AC#1`③"同笔连 pin 一起改"）。**
- **甲形的编译面凭据形状已验通**（票面 `AC#2` 要求的那一发）：`go vet -tags winlive ./cmd/wisp/ ./internal/ball/` ⇒ **`G1_rc=0`**（件 `g1-vet-winline.txt`）；`go test -tags winlive -list '.*' ./cmd/wisp/` ⇒ **`G2_rc=0`、名册 301 枚、含那三枚名**；**同尺正控**＝winlive-tagged 文件里的用例确实出现在这份名册（`cmd/wisp/panel_host_windows_live_test.go` 的 `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive` 命中）⇒ 搬进 `winlive` ⛔ 等于"既⛔ 跑也⛔ 编"。
- ★**没有任何一形能让 `--scope=cli` 变绿**（甲／乙只摘 3～5 枚 FAIL，剩下 4 枚别族 FAIL ＋ 2 枚 `--- SKIP` 各自撞 `runtests.sh:98`→`:102`；那 17 行我逐字复跑过）。⇒ "裁哪一形"裁的是**名册可读性**，⛔ 灯色。这句今后写进任何派单的第一段。
- 另两枚小腐烂（⛔ 门，但会说谎）：`ci.yml:625-626` 注释逐字 `7 under cmd/wisp, 5 under internal/ball`／`all twelve` ⇒ 加第 8 枚 winlive 文件即过期，而 `winlive` 在 `scripts/`＋`tools/` **0 命中** ⇒ 没有门钉这句注释；33-n1 引的 `ci.yml:388` 已腐烂（现量 `:534`）。
- ⚠ 甲与乙**互斥**：`:652` 的 `go test -list` ⛔ 带 `-tags`、⛔ 吃 `-skip` ⇒ 甲之后走乙 ⇒ 查⛔ 到名字 ⇒ `:667`→`:674 exit 1`（当场红）。腿那条"派单问反了方向"我认：我把陷阱写在了"改名⛔ 会被发现"那一处，**真静默面在"只验名⛔ 验判据"＋RUN 计数无标记蒸发**那一处。

### ⑤ 顶回我的七处，逐枚认（具名在 `40-census.md` §分歧）

1 档数（盘上 5 档 `tiers='core windows cli winsec census'`，我 X1 复跑对上，⛔ 我派单写的"三档"）／2 名册 3→5＋1／3 票面"永远拿不到回执"（见①）／4 `-skip` 静默方向问反（见④）／5 `go vet -tags winlive` **覆盖得到** `cmd/wisp`（我 V2 复跑 `ci.yml:672` 逐字含 `./cmd/wisp/ ./internal/ball/`）／6 winlive tag 文件＝**12** 枚（我 X2 独立尺复跑＝12；`33-n1` 那句"8 枚"作废）／7 派单预告"这五个目录此刻躺着别人脏件"⛔ 对（那五个路径此刻 0 行，脏件在全仓别处）。
另外腿自报的射程口径值得进模板：它的 blob 行号全部锚在**起手 `6414a4bb`**，交件时 HEAD 已被别的腿推到 `83cd66a8` ⇒ **引用它的行号前先重跑**（同族第四次）。

### ⑥ 裁形（裁决权＝编排者；本格⛔ 由落地腿自判）

- **`AC#0` 翻勾**（三张表齐、必答④ 有判读、七处顶回具名；缺的三枚读数里两枚我已代取（④ 第一、二条），bisect 与 ubuntu 台面那一枚⛔ 属本票射程）。
- **`AC#1` 按住、⛔ 派 `302-r1`**。具名前置条件两条：**(i)** 票 303 交回"页面→Go 回执断口"的归因；**(ii)** 只有当归因结论是**"这一跳在 CI 台面上永远拿不到回执"**时，搬档／进 ledger 才有意义——若结论是"能拿到，是某笔改动弄断了它"，那修完之后这三枚本该回绿，**本票甲／乙两形都⛔ 该落**（落了就是把一枚真回归搬进隐形，正撞票面 `AC#2` 自己警告的那句）。
- **`AC#1`② 那句"⛔ 动 `cmd/wisp/**` 的产码"的口径我现在定死**（⛔ 改原句）：**"产码"＝非测试文件**；在 `cmd/wisp/**` 内把测试函数／harness **逐字搬**到另一枚 `*_test.go` 属甲形本体，判据正文⛔ 改一字，凭据＝搬前搬后那几枚函数的 `cmp` 逐字节对拉。
- 本票状态＝**1 勾（`AC#0`）／4 未勾**，⛔ 改 `-done`。

## Progress log（本节）

- [2026-10-10 16:4x +0800] agent=编排者 did=收 `302-a1`（件 `probes/302/a1/40-census.md`，commit `8267be8a`，逐笔名册＝自家目录 2 枚，⛔ 越界）⇒ **翻 `AC#0`**＋**`AC#1` 按住**＋★票面四句前提就地改写（`:1`／`:4`／`:14`／`:15`）＋★新量到一枚**可本机复现的回归**（`G4`：三枚本机全 `--- FAIL (20.01s)`、红句逐字同 CI ⇒ "环境差"作废，⚠ 我此前那句是拿旧码探针当现行读数）＋★我 16:0x 那句"三枚都是新可见"⛔ 只对 1 枚（两枚基线**绿**）记我＋三枚读数我代腿取（`G1` rc=0／`G2` rc=0·301 枚含三枚名·带正控／`G3` `0xc0000135` rc=1 列名 0 枚 ⇒ **乙形今天⛔ 可行**）＋认七处顶回；next=**票 303**（页面→Go 回执断口的归因，本机可复现⇒⛔ 依赖 CI）→ 回来再裁本票甲／乙 → 票 300 `300-v3` → 票 111 `AC#12`

## 编排者裁（2026-10-10 19:1x，CI 改后那一发到手）⇒ §⑥ 的前置条件 (ii) **判死**：甲／乙**都⛔ 落**，`AC#1` 继续按住（⛔ 翻任何一枚 `- [ ]`）

件＝`.scratch/wisp/probes/303/orch/r9-ci-after-red-roster.txt`（票 303 那格 CI 色的原件，本票**借用**它，⛔ 复制一份读数到本票目录——先例＝同一份读数⛔ 两处重记）。

### ① 承重那一句（逐字，尺＝`gh run view --log --job 114194107792`，headSha `05db4bc6`＝票 303 修复那一笔）

`cmd/wisp/panel_resident_windows_test.go:825` 逐字：
`AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`
⇒ `--- PASS (0.30s)`。同一枚用例在**修复的父发**（`cf46c24a`，job `114187968428`）是 `--- FAIL (20.02s)`＋逐字 `:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)`。

- 这回答的是 §⑥ 那条**具名前置条件 (ii)**："只有当归因结论是『这一跳在 CI 台面上永远拿不到回执』时，搬档／进 ledger 才有意义"。盘上答＝**结论⛔ 是那一句**：CI 那台机器**拿得到回执**，而且判据是页面自己报回的（⛔ 是 Go 侧自说自话）。
- ⇒ **甲（搬 `winlive` 档）／乙（进 `portable-tests.sh` 的 `-skip` ledger）都⛔ 落**。落下去的后果正是 §⑥ 自己警告的那一句：**把一枚真回归搬进隐形**（还撞上 `AC#2` 那条"⛔ 让分母变成既⛔ 跑也⛔ 编"）。
- ⛔ 由此翻 `AC#1`：本格是**落地格**，而本程零产码；它的正确现态＝"按住，且本票⛔ 有可落的形"，⛔ 是"做完了"。

### ② 票面 `:11` 那句"另有 1 枚基线红→改后 `--- SKIP`＝`TestPanelHostLatencyPercentilesAC2`"——**这一发又翻回去了**（原句一字不改，留在上面当快照）

同一把尺（步日志 `portable-tests.sh: four numbers` 那行逐字）在 `cmd/wisp` 档：`RUN=398 PASS=278 FAIL=8 SKIP=2` → `RUN=398 PASS=279 FAIL=8 SKIP=1`。
- 名集合作差（⛔ 比枚数）＝**被修好 1 枚（nail1）／新增 1 枚（`TestPanelHostLatencyPercentilesAC2` 由具名 SKIP 走进 `--- FAIL`）**；`FAIL=8↔8` 同数而名不同 ⇒ ⚠ **两把尺⛔ 同物**，本票 `:11` 早就为同一枚坑写过一句同款警告，这次是它的第二次。
- 那枚新红的逐字＝`:1005: cold P95 2747.460 ms over 1 runs exceeds the D32 panel cold budget 1500 ms (the single-run assertion is not the only gate: the tail is what AC#2 asks for)`；它 SKIP 掉的前提（`:985` 逐字 `no cold/hot sample recorded in this process … Named skip - an empty aggregate is not a green latency gate`）随修复消失——上面的 lifecycle 用例**真的产出了样本**，判据才开始求值。
- ★这一枚**⛔ 属本票射程**（本票管"档／放置"，这枚的颜色变化来自"回执通道接上"＝票 303 的产码面），但它又**⛔ 是新坏**：§① 那张基线表里同一枚回归前就红过（`--- FAIL (0.00s)`）。⇒ 具名交下一程裁归口（候选＝票 33 的 D32 预算面／票 305 的文档交接面），⛔ 我盖章，⛔ 为它动 1500 那一枚数。

### ③ 本票现在的残差清单（谁仍红、归哪一票）——这是"丙＝登记成具名已知红"那一支能剩下的全部

- `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC14GoSideEvalPushReachesThePage`：**仍红，且已换形**（CI 上逐字＝`:327 AC#13 page answer … "0" - 0 of 1 probe id(s) present in the live document`；`:867 AC#14 nail 2 … page's own words: title=""`）⇒ 归**票 305**（产码面＝`bringUp` 的 `serveEntry`／探针文档交接次序＋`Eval` 推那一条），⛔ 本票的档面。
- `TestPanelHostRealWindowHopAndLifecycle`＋上面那枚百分位：红在**预算**那一支 ⇒ 归口⛔ 定（见 ②）。
- `internal/risk` 那 12 枚（A/List 表判定与 under-profile fallback 族）⇒ 仍按票面边界段留在台账 `A817`，先例＝票 115，⛔ 塞进本票。
- ⇒ 本票⛔ 改 `-done`（`AC#1`／`AC#2`／`AC#3`／`AC#4` 四枚 `- [ ]` 一枚没做完；先例＝票 62 那族"已收口而仍有未勾框"的教训，本票连收口动作都⛔ 触发）。
- ⚠ 一处台面口径要写进后续派单：本票 §① 的基线表用的是 10-06 那发（`cc315261`）与 10-10 改前那发（`cf46c24a`）对拉；本次新增的是**第三发**（`05db4bc6`）。三发同档同版本（逐字 `Image: windows-2025-vs2026`／`Version: 20260925.250.1`／agent `2.337.0`）⇒ "镜像漂"那一支在这三发之间⛔ 有依据。★而我 18:0x 在 `r7-ci-before-red-roster.txt` 里写的"`ImageVersion` 那把尺⛔ 命中"＝**我那把 grep 找的是 API 字段名，而日志里那行叫 `Version:`**，⛔ 是"盘上没这行"。**记我。**

- Progress log：
  - [2026-10-10 19:1x +0800] agent=编排者 did=取 CI 改后那一发（run `38045508579`／job `114194107792`／headSha `05db4bc6`）⇒ **§⑥ 前置条件 (ii) 判死**：托管 runner 上拿到了页面原话 `REPLIED,REPLIED,REPLIED`（`:825` 逐字）⇒ **甲／乙都⛔ 落**、`AC#1` **继续按住**且⛔ 翻框（本票⛔ 有可落的形，⛔ 是"做完了"）；★票面 `:11` 那句"1 枚基线红→改后 SKIP"这一发**又翻回去**（`TestPanelHostLatencyPercentilesAC2` 具名 SKIP→`--- FAIL`，逐字 `cold P95 2747.460 ms … exceeds the D32 panel cold budget 1500 ms`）＝回执接上的**次生色**，⛔ 本票射程、⛔ 新坏（基线发同枚 `--- FAIL (0.00s)`），归口交下一程；⚠ 四数尺 `FAIL=8↔8` 与名册作差⛔ 同物（同一枚坑本票第二次踩到，我按 `:11` 原句口径写清）；残差清单三条各自归口（票 305／预算归口未定／`internal/risk` 12 枚留 `A817`）；★`r7` 里我那格"ImageVersion ⛔ 命中"改正是**我的 grep 找错载体**（日志那行叫 `Version:`），定式并回＝换载体取读数先现量字面写法；本票⛔ 改 `-done`；next=票 303 `303-v1`（含本件 ② 那枚归口的必答题）→ `305-a1`
