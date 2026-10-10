# 票 303 — 「页面 → Go」那一跳的回执在 `cc315261`（10-06 那发 CI）还拿得到、在 `bcd0a543`（10-10 这批）拿不到了，而**今天连本机也拿不到**：这是票 302 那三枚红里被"环境差"那句话盖住的**一枚真回归**

**立票**：2026-10-10 16:4x 编排者（来路＝`302-a1` 的表② ＋ 我自己现跑的 `G4` 那发本机复跑；件＝`probes/302/a1/40-census.md`、`probes/302/orch/g4-local-three-cases.txt`；账＝`A819`。**⛔ 塞进票 302**——那一枚修的是"判据长在哪个档"，这一枚修的是"那道门断了"，两件事的判据形状与代价面完全不同）
**性质**：★**缺陷票**（regression），⛔ 仪器票、⛔ 归口票。
**为什么要紧（用户可见形状）**：面板与宿主之间那道门（票 35 的 `installPanelTransport`／消息钩子那条边）是**产品路径共用的**〔读码推的，⛔ 实测——把这句判死＝本票 `AC#2` 的一格〕。盘上现在的形状是：窗**建得起**（`window_opened=true`、`panel thread exited cleanly shows=1`）、页面**收得到 Go 的推**、但 Go **收不到页面的话**（`what DID arrive at the door: nothing at all`，15 s 超时）。最坏后果的形＝**"球/热键把面板叫出来之后，点什么都没反应"**。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**同一条判据的三色**（尺＝`^--- [A-Z]*: <名>` 在两发归档切片里逐枚打；件＝`probes/301/orch/logs/ci-baseline-cli-block.txt`／`ci-after-cli-block.txt`；两发 Runner Image 与 agent 版本逐字相同＝`windows-2025-vs2026`／`20260925.250.1` ⇒ "镜像漂"那一支⛔ 依据）：
  `TestAC14AwaitedBindingReplyReachesThePage` **PASS(1.25s，带页面的话 `REPLIED,REPLIED,REPLIED`) → FAIL(20.04s，零回执)**；`TestAC14GoSideEvalPushReachesThePage` **PASS(0.97s，`title="PUSHED-33R5-OK"`) → FAIL(20.03s)**；`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` `--- SKIP`（具名理由逐字 `the embed resolves no entry`）→ FAIL（这一枚是**新可见**，⛔ 新坏——它 10-06 那发根本没有对象可判）。
- ★**同发里"预算那一枚"换因**（尺＝每枚取 `=== RUN` ↔ `--- FAIL` 之间那段；⚠ Go 的 `-v` 把明细写在 `--- FAIL:` **之前**）：`TestPanelHostRealWindowHopAndLifecycle` 基线红句逐字 `cold bring-up 2889.2 ms exceeds D32 panel cold budget 1500 ms`（`panel_host_windows_test.go:665`）→ 改后逐字 `did not produce a browser round trip (got -1.000)`（`:662`）。⇒ **枚数没变而被名集合差读成"没动过"**，这一枚的"换因"是本票与票 302 的接缝。
- ★**本机 HEAD 复现**（尺＝编排者 16:4x 现跑，`PATH` 铺了 `third_party/sherpa-onnx`＋`build`，`-count=1 -timeout 420s -v -run '<那三枚>'  ./cmd/wisp/`）⇒ **`rc=1`、三枚全 `--- FAIL (20.01s)`、红句逐字同 CI**。射程口径：被测＝工作树，而 `git status --porcelain -- cmd internal scripts .github docs` 现量 **0 行** ⇒ 工作树≡HEAD；`tasklist` 起手 `wisp.exe=0`／`balldebug.exe=0`。⚠ **单发**，复现率⛔ 知（`AC#0` 的活）。
- **票面那句旧话作废**（写在这里防下一个人都踩）：票 302 `:15` 我 16:0x 写的"本机同码⛔ 红 ⇒ 这是环境差"引的是 `33-p1`，那是**旧码**的读数；当前 HEAD 本机**红**。⛔ 再拿"环境差"替这三枚说话。
- **嫌疑面（尺＝`git log --oneline cc315261..bcd0a543 -- <4 files>`；⚠ 枚数随文件集变，两把都具名）**：
  腿那把（4 枚文件）＝ **5 枚**；我这把（`cmd/wisp/panel_host_windows.go`／`cmd/wisp/panel_resident_windows.go`／`cmd/wisp/panel_pageover_33r10_windows_test.go`／`cmd/wisp/panel_transport_35r2_test.go`）＝ **7 枚**＝`9995f9b1`（33-r11 改名 `firstRoundTripLocked`→`firstRoundTrip`，自述 `Zero lines added or removed`）／`8b32060b`（255-r1 面板宽度到真窗）／`70b00885`（33-r10 冷启交接补 window-free 尺）／`3a343bc7`＋`2fc5f5c9`（票 35 `AC#8` 夹具面，⛔ 产码改动）／`286a7f30`＋`fb2fb802`（票 35 `installPanelTransport`／`wispDispatch` 绑定＋postMessage 转发 Init）。
  ⇒ **只有 `fb2fb802`／`286a7f30` 逐字落在那条"页面回话"的边上**（消息钩子的再入形状、绑定名注入形状）。⚠ 本票**⛔ 归因**——这是 `AC#1` 的产物，⛔ 由票面替它填数。
- ⚠**bisect 的一枚陷阱（现量，别浪费一发）**：`TestAC13ColdStart…` **依赖 `frontend/dist` 的产物字节**（CI 那发逐字 `embed resolves 1068 entry byte(s)`；本机也有旧产物），而**仓里只跟踪 `.gitkeep`** ⇒ 干净的仓外 clone 里它会走"no subject"具名跳过、**复现⛔ 出来**。⇒ **最小复现集要挑 `TestAC14AwaitedBindingReplyReachesThePage`**：基线那一发 `embedded assets are not built` 时它**照旧 PASS 并拿到回执** ⇒ ⛔ 依赖 dist。

## 要建什么（`AC#0` 之前⛔ 任何产码）

- [x] **`AC#0` 复现钉死（⛔ 归因）**：① 起点／终点具名（起点＝`cc315261` 那台 CI 台面能拿到回执；终点＝`bcd0a543` 或**当前 HEAD**，两者要分清并各留一发）；② 最小复现集＝**一枚用例名**＋跑法逐字（含 `PATH` 那两枚目录、`-count`、`-timeout`），并**具名回答它依赖⛔ 依赖 `frontend/dist`**；③ **复现率**＝同一条命令 `-count=3` 的三色（⛔ 单发当恒红）；④ 一枚"这⛔ 是我这台机器的毛病"的对照＝**同码在 CI 那发也红**（引用现量第 3 条，⛔ 重跑 CI 来证这句）。判据＝每条带"哪把尺＋射程目录＋工作树还是 HEAD"。
- [x] **`AC#1` 归因到**一笔**（⛔ 一块批）**：在**仓外** clone 里做（⛔ 在仓内建 worktree 或 checkout；AGENTS §1.4），交回＝`git bisect` 的**那笔 SHA**＋该笔的**定向复现两发**＝(a) 在该笔上跑最小复现集 ⇒ 红，(b) 只把那⼀笔单独 revert（或 `git checkout <父号>` 同尺）⇒ 绿。**两发缺任一发＝本格⛔ 交**。若两枚以上同犯（例如"绑定名"与"再入帽"各断一半）必须具名写"单 revert ⛔ 复绿"并给组合读数。⛔ 共享工作树上做任何 checkout。
- [x] **`AC#2` 机制层（三形分开答，⛔ 合一格）**：**(甲) 页面 JS 压根⛔ 跑**／**(乙) JS 跑了但绑定名⛔ 接上**（`wispDispatch` 那枚注入名）／**(丙) 钩子接上了但门被再入帽挡住**。每形必答两问：① **能把这一形与另两形分开的读数长什么样**（现量给得出就给，给不出就具名写"缺哪枚读数＋谁来取"）；② 本票字节到不到得了那一步（`302-a1` 已判：这三形在现有 CI 字节里**分⛔ 开**，`-1.000` 与"零回执"**同时出现**）。⚠ 判"产品路径也断"必须⛔ 依赖这三形里任何一形的**假设**，要一发真机或一条真调用链读数。
- [ ] **`AC#3` 修复（只在 `AC#1`/`AC#2` 交完之后开工）**：**硬约束四条**——① ⛔ 放宽任何断言：`no report … within 15s (what DID arrive at the door: nothing at all)` 那三句判据本体一字⛔ 动；② ⛔ 把 `t.Fatalf` 换成 `t.Skip`、⛔ 造 `--- SKIP`（`tools/d22scan/runtests.sh:98`→`:102` 把任何 `^--- SKIP` 判红，那是**故意的**）；③ ⛔ 动 SLO 阈值／`thresholds.go`／D32 那一面（`TestPanelHostRealWindowHopAndLifecycle` 基线红在**预算**那一支属 SLO，⛔ 用"改预算"换绿）；④ 产码改动落在**门**那一侧（`cmd/wisp/**` 宿主传输边），⛔ 顺手改档／tag／ledger（那是票 302 的射程，本票⛔ 动 `scripts/portable-tests.sh`）。凭据＝成对两发＋三枚本机色（修后 `TestAC13`／两枚 `TestAC14*` 全绿）。⚠ **2026-10-10 21:4x 编排者判：这半句凭据在本票里从来⛔ 可能成立**（那两枚属票 305 的症状，逐字凭据见下面"收 `303-v1`"那一节第 2 节）⇒ **本格⛔ 翻、原句⛔ 改**，可满足的部分搬到下面新格 `AC#3b`。
> **[2026-10-10 22:2x +0800] 编排者·`AC#3` 那枚"不可满足判据"的正式登记（形照票 141 的先例＝`141-…-done.md:50-54` 那段登记＋`:59` 那枚撤销口令；本登记⛔ 翻这一格的勾，`-done` 那一步从此⛔ 需要它先闭）**：本格凭据末句"修后 `TestAC13`／两枚 `TestAC14*` 全绿"**结构性⛔ 可满足**＝它把两枚**本票射程之外、且属还没做出来的功能**算进了分母。盘上三边凭据（各带件名）：① `nail2` 的红句在**回归之前**那一发就是逐字同一句（我 `probes/303/orch/r3-f718e9b6-three-nails.txt` ＋ 腿件 `probes/303/v1/logs/`，`f718e9b6`＝`fb2fb802` 的父发、同一枚干净 clone 台面）；② `TestAC13…` 在干净 clone **恒 `--- SKIP`**（母仓跟踪的 `frontend/dist` 只有 `.gitkeep`，本机有旧产物、clone 没有＝票 252 那族"产物缺失"形），在母仓**恒红**且红的是"活文档还是探针那张页"＝**票 305 症状①**；③ 非实现者 `303-v1` 对这格的判语（件 `probes/303/v1/10-verdicts-ac3-ac4.md:48-51`）＝**ⓣ 判"凭据句本身有缺陷"（⛔ 判"修复不成立"）／ⓤ 判"更正形状（nail1 绿＋两枚残余另立票 305）下闭合，附 ⓓⓖ 两条件"**，ⓓⓖ 我 21:4x 各自现跑核闭（ⓓ＝整包红名册逐名对拉、ⓖ＝票 305 已在盘上）。⇒ **处置＝原句⛔ 删、本格⛔ 翻、可满足那半句落在下面 `AC#3b`（已勾）**；⛔ 任何人把这一格读成"漏勾"，也⛔ 拿"AC#3 未勾"说票 303 的修复没做——**修复做没做看 `AC#3b`＋`AC#4`＋本票"收 `303-r1`／收 `303-v1`"两节**；⛔ 为凑那半句放宽任何断言／改 D32 那枚 `1500`／造 `--- SKIP`（照票 141 那句"不要为此把这格退回，也不要为了凑它去放宽箭头段"）。**撤销本登记＝口令「撤 303 AC#3 登记」**（射程＝删本段；`AC#3b` 的勾与 `AC#4` 的勾⛔ 随之撤，撤它们要另外的口令）。
- [x] **`AC#3b`（2026-10-10 21:4x 编排者追加，来路＝`303-v1` 的 ⓤ 判语＋我自己三条对拉；它承担 `AC#3` 那半句⛔ 可满足的凭据，⛔ 取代 `AC#3` 上面任何一行原文）**：修复那一笔在本票射程内**成立**的判据＝以下三件**同时**真——① `AC#3` 那四条硬约束逐条有编排者自己现跑的尺与读数（⛔ 任何一条只凭腿自述）；② 修后 `TestAC14AwaitedBindingReplyReachesThePage` **绿**且判据**由页面自己报回**（逐字 `REPLIED,REPLIED,REPLIED`，本机＋CI 两侧各有）；③ `cmd/wisp` 整包改前／改后**各两发、同一台面取交集＝稳定新增红 0 枚**，且**单发红过的那一枚⛔ 被这把尺藏住**——必具名报出并写清归因（本程＝`TestTicket223HandEditedFsLooseningCostsAnL2Card`，⛔ 归本票修复）。⇒ ⛔ 拿本格去说"票 303 修完三枚都绿了"：剩下两枚的账在**票 305**。判语正文与三条读数＝票面"编排者收 `303-v1`"那一节（⛔ 在此重复，免成第二份真相源）。
- [x] **`AC#4` 门禁与越界**：`go vet ./cmd/wisp/` rc=0；`go vet -tags winlive ./cmd/wisp/ ./internal/ball/` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l cmd/wisp` 与 HEAD blob 那把**并排两把**都要报（先例＝票 298：加严残留具名、⛔ 顺手修）；`cmd/wisp` 整包改前／改后各两发取交集 ⇒ **新增红 0 枚**（⚠ 本机整包要铺 sherpa `PATH`，缺 DLL 是 `0xc0000135` 且**无 `--- FAIL`**＝用例根本没跑）；`git show --name-only --format=` **逐笔**名册＋与授权名册差集（⛔ 区间尺）；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec 且写在 `$( … )` **之外**。
- [x] **`AC#5` 交回后由编排者办的两件**（⛔ 腿做）：① 票 302 的甲／乙裁语按 `AC#1` 结论回填（若"这一跳在本机都断"⇒ 修完就该回绿 ⇒ **甲／乙都⛔ 落**）；② 本票那一发 CI 色（修后推送产生的 `test-windows` 四数与红名册）。

## 边界（⛔ 本票射程，混进来＝顺手改错那半）

- `internal/risk` 那 **12** 枚 CI 红（名单＝`A817` §3；一族读的是**哪棵树**，先例＝票 115 `AC#2`/`AC#3`）。
- 票 302 那三枚的**放置**（档／tag／ledger 文案）。
- `TestPanelHostRealWindowHopAndLifecycle` 基线那一支**预算红**（D32／SLO 面，`AGENTS §1.1` 一字节⛔ 动）——本票只许把它从"回执为零"移回"红于预算"，⛔ 更多。
- 那 4 枚 `mockllm`／config／`goroutine outside the D38 roster (leak symptom)` 一族的 cli 红（⛔ 窗⛔ 回执）。

## 规矩（本票全程）

- 子代理**只 commit、不 push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree 或 checkout（bisect 一律在**仓外** clone）；临时件**只建不删**；证据件⛔ 叫 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`）。
- 跑任何真机／整包前先 `tasklist` 现量 `wisp.exe`／`balldebug.exe`＝0；杀进程只停本腿自己起的那枚。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件；⛔ 为变绿放宽任何断言。
- 裁决者≠实现者：`AC#1`/`AC#2`/`AC#3` 判语⛔ 由落地腿自勾。
- 证据件一律只新建 `.md`／`.txt`（⛔ `.go` 进 `.scratch`——会挪分母；其余文件类型实测⛔ 动门，见 `A818`）。

next=`303-a1`（`AC#0`＋`AC#1`：复现钉死＋仓外 bisect 到一笔）→ 编排者复跑 → `302` 甲／乙回填裁语 → `303-r1`（修复）→ `303-v1`（非实现者）

## 编排者收 `303-a1`（2026-10-10 17:3x，六笔 `b3f9c6d7`→`0edce3b6`→`f2ac16dd`→`48ee655e`→`86ed3592`→`9ebc8b35`，件 `.scratch/wisp/probes/303/a1/**`，逐笔名册与授权名册差集＝∅，⛔ push）⇒ **翻 `AC#0`＋`AC#1`**，归因**收敛到一笔**＝`fb2fb802`，★定向两发**由我在同一枚仓外 clone 上逐字复跑对上**

### ★★结论一句话（这条是本票全部的价值）

「页面 → Go 的回执」那一跳**不是环境问题、不是档放错、不是 flake**，是 **10-07 18:05 那笔 `fb2fb802`（票 35 的 35-r1：`installPanelTransport` 绑上 `wispDispatch` 门＋把 postMessage 转给 Init）** 弄断的。**单撤这一笔就复绿** ⇒ ⛔"多枚同犯"那一支没触发。

### `AC#0` 凭据（腿交 ＋ 我复跑）

- **最小复现集＝一枚**：`TestAC14AwaitedBindingReplyReachesThePage`。跑法逐字（腿件 `10-repro-head-count1.txt`／`11-repro-head-count3.txt`）＝
  `export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH"` ＋ `go test -count=1 -timeout 420s -v -run 'TestAC14AwaitedBindingReplyReachesThePage' ./cmd/wisp/`
  ⚠★一枚**形状级**新知（我没料到的）**PATH 必须是 shell 形 `/d/…`**，写成 `D:/…` 就直接 `0xc0000135`（腿的尺＝`scripts/wisp-cli-tests.sh:101-109`；它起手把它写成 `113-118` 并自己更正，笔 `f2ac16dd`）。
- **复现率**＝同命令 `-count=3` ⇒ **FAIL 3／PASS 0／SKIP 0**，判据句命中 3 次（`rc=1`）⇒ 我此前那句"单发、复现率⛔ 知"的欠账**由这一发销掉**。
- **⛔ 依赖 `frontend/dist`**：干净 clone（`frontend/dist` 只有 `.gitkeep`）里起点照 PASS、终点照 FAIL（腿件 `16-`／`17-`）⇒ 票面写死的"最小集挑这枚"是对的。
- **"同码在 CI 也红"那一枚**＝引用归档字节 `ci-after-cli-block.txt`（腿⛔ 为它推 CI）✓。

### `AC#1` 凭据（★这一格我只认自己复跑过的那一发）

- bisect 台面＝腿建的**仓外** clone `~/wisp-303-bisect`（`cc315261`→`bcd0a543`＝742 枚、0 merge）；**12 步全有效（GOOD 4／BAD 8／INVALID 0）**，判红**只认逐字判据句**（`21-bisect-log.txt`／`22-bisect-index.tsv`）。⚠ 这正是我派单里点名的那枚陷阱：`0xc0000135`／`[build failed]`／`no tests to run` 都算**无效步**⛔ 算 bad——它 0 枚无效步＝判据守住了。
- first bad＝**`fb2fb802f75a0e3eeacad488f1adc6f064e29f85`**，在那把 7 枚嫌疑名册里（＝我 `A819` 记的那 7 枚之一，⛔ 当时判死，现在判死了）。
- 那一笔动了几枚文件（尺＝`git show --name-only --format= fb2fb802`，我 16:4x 就跑过这把尺）＝**2 枚**：`cmd/wisp/panel_host_windows.go`（**产码** +57/−5）＋`cmd/wisp/panel_transport_35r1_test.go`（新增测试码 +243）⇒ ⛔ 纯夹具，那把"是产码断还是尺断"的刀**落在产码侧**。
- ★**定向两发由我逐字复跑**（件 `.scratch/wisp/probes/303/orch/r1-orch-pair-summary.txt`＋`r1-pair-a-culprit.txt`／`r1-pair-b-parent.txt`，17:2x，同一枚 clone，复跑前 clone 工作树脏枚数＝0）：
  (a) `fb2fb802` ⇒ `rc_pair_a=1`、`--- FAIL: TestAC14AwaitedBindingReplyReachesThePage`，红句逐字 `no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)`；
  (b) 其父 `f718e9b6` ⇒ `rc_pair_b=0`、`--- PASS`，绿句逐字 `page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`。
  ⇒ 两发都**与腿的读数一字不差** ⇒ `AC#1` 我判**成立**。
- 腿另交的两发加料（⛔ 本格所要求，⚠ 我**没**复跑，读数按〔仅腿报〕记）：删掉那笔新增的测试文件⇒**仍红**（＝断口⛔ 靠那枚测试存在）；把那笔的产码半边换成父发⇒`installPanelTransport undefined`＝**INVALID**（＝那一笔的产码与测试互相引用，⛔ 能半边撤）。

### 顶回我一处，理由在盘上，我认

票面 `:10`（现量段第二行）那句"`TestAC13ColdStart…` 是 SKIP→FAIL 那一枚"——**只对母仓／CI 台面成立**（那里的 `frontend/dist` 有产物字节）；**干净 clone 今天仍 SKIP**（腿件 `18-`／`19-`）。⇒ 派单／票面凡引用"某枚用例今天是什么色"，**必须同时写台面**（母仓工作树／仓外干净 clone／CI），这是"射程口径"那一族第五次同名事故。

### 本票现态与 next

- **`AC#0` 翻、`AC#1` 翻**（凭据如上）；**`AC#2`（机制三形）一枚没量**＝腿具名欠，⛔ 它的射程（派单里我禁了）；**`AC#3` 修复⛔ 开工**（票面写死要先有 `AC#2`）；`AC#4` 门禁随修复；`AC#5` 归我。
- ⚠ **残体一枚具名记下（机主台面上看得见的事）**：本腿终态 `msedgewebview2.exe` ＝**24 枚**（起手⛔ 量，它自己报了），`wisp.exe`／`balldebug.exe`＝0，宿主 `go test` 全部已退 ⇒ 那 24 枚是测试遗留的**孤儿渲染进程**。⛔ 我按"只杀自己起的"的规矩处理＝**⛔ 杀别人的、⛔ 假装没事**：由我在下一程具名量一次父链，再决定停哪些，并把"为什么它们没随窗退出"这件事**单独归口**（窗退出净空那一支本来就有 `*_WinLive` 那枚 2s 判据在管，那是真窗档、⛔ 进 CI）。
- Progress log：
  - [2026-10-10 17:3x +0800] agent=编排者 did=收 `303-a1`（六笔，件 `probes/303/a1/**`，越界 0）⇒ **翻 `AC#0`／`AC#1`**；归因**收敛到一笔 `fb2fb802`（35-r1）**，定向两发**我在同一枚 clone 上逐字复跑对上**（(a) rc=1＋`nothing at all`／(b) rc=0＋`REPLIED,REPLIED,REPLIED`）；★销我自己那枚"复现率⛔ 知"的欠账（`-count=3`＝FAIL 3/3）；★认一处顶回（`TestAC13…` 的颜色**随台面变**：母仓 FAIL／干净 clone SKIP／CI 基线 SKIP）；★新知一条：**PATH 的 shell 形 vs `D:/` 形**决定 `0xc0000135`；⚠ 残体 24 枚 `msedgewebview2.exe` 孤儿进程记在案、由我下一程具名处置；next=`303-r1`＝先答 `AC#2` 机制三形（页面 JS⛔ 跑／绑定名⛔ 接上／再入帽挡门），再动那一笔的产码半边，**⛔ 放宽任何断言、⛔ 造 `--- SKIP`、⛔ 碰 D32 预算**


---

## 编排者收 `303-r1`（2026-10-10 18:3x，两笔 `e3341368`→`807497c1`，件 `.scratch/wisp/probes/303/r1/**`，逐笔名册越界＝∅，腿⛔ push＝我这一程按 `A813` 推）⇒ **只翻 `AC#2`**，`AC#3`／`AC#4`／`AC#5` ⛔ 翻，★**票面那句"修后三枚全绿"是我造的，作废并具名更正**

### 一句话

「页面 → Go」那一跳的回执**被接上了，而且我自己复跑对上了**；票面当初顺手算进本票射程的**另外两枚症状**与本票无关，已单独立成**票 305**。

### 1. 我现跑的三发（每发都写台面，⛔ 拿腿的话当凭据）

- **母仓 HEAD 定向三枚**（件 `probes/303/orch/r2-orch-targeted-after.txt`，`rc=1`）：nail1 `--- PASS`，绿句逐字 `page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`；`TestAC13ColdStart…` `--- FAIL`（页面自答 `"0"`）；`TestAC14GoSideEvalPush…` `--- FAIL`（`the page reports its title as "", want "PUSHED-33R5-OK"`）⇒ 与腿的 `51-restore-green.txt` **逐字同**。
- **clone 成对两发**（件 `probes/303/orch/r5-r6-orch-pair-verdict.txt`，★还的就是腿具名欠的那格 `30-package-before.txt`＝5 字节只有 `rc=1`）：改前 `detached=cf46c24a`＝顶层 268 PASS／17 FAIL／3 SKIP，改后 `detached=807497c1`＝270／16／2。**被修好 2 枚**＝nail1 ＋ `TestPanelHostRealWindowHopAndLifecycle`（★后者在**母仓**腿那发 `31-` 里改后仍 `--- FAIL`、红句仍是 `did not produce a browser round trip (got -1.000)` 那一支 ⇒ **颜色随台面**，这条进票 305 症状①那一族）。**逐名作差给出新增红 1 枚**＝`TestAC4FocusReturnToPriorWindowGap33r5`。
- **交替换头隔离发**（件 `probes/303/orch/r8-ac4focus-verdict.txt`）：同一枚 `807497c1` 二进制在 **47 秒内** 5 连绿 → 3 连红，而 `cf46c24a` 那枚 5 连红、可它在整包那一发里是 `--- PASS (5.28s)` ⇒ **两枚 SHA 各自都出现过两色** ⇒ 判语＝那枚"新增红"**⛔ 能归因给修复**（红句在有／无修复两支下一字不差，而它测的是 Show 时前台焦点归谁，与那 +27 行 JS 转发钩子不在同一层）。⚠ **⛔ 因此把它当"flake 就没账了"**：本机现量 13 发＝5 绿／8 红（尺＝三发定向 13 枚计数）⇒ 单独立成缺口归票 304 那一族（"这台面上一枚真窗焦点判据⛔ 有可重跑的确定性"），⛔ 混进本票 `AC#4` 的门禁。

### 2. ★顶回我那一句，理由在盘上，我认（票面原句⛔ 改，更正写在这里）

票面 `AC#3` 我写的凭据末句"修后 `TestAC13`／两枚 `TestAC14*` 全绿"**⛔ 成立**——而且它当时就是把两枚**本票管不到**的症状算进了射程。尺＝我在同一枚 clone 里跑的 `f718e9b6`（＝`fb2fb802` 的**父发**，回归之前）定向三枚（件 `r3-f718e9b6-three-nails.txt`）：nail2 红句**逐字同一句** `the page reports its title as "", want "PUSHED-33R5-OK"`；`r4-clone-at-HEAD-nails.txt` 在 `807497c1` 复现同形。
⇒ **判语①**＝nail2 的那枚红**⛔ 是 `fb2fb802` 造成的**，它在回归之前的同一枚台面上就是同一句话。**判语②**＝本票 `AC#1` 那笔归因（收敛到 `fb2fb802`）**不受影响**，因为它钉的是 nail1 那一跳，而那一跳有成对两发＋我这发母仓复跑。两枚残余症状＝**票 305**（新立的**台面规则**：前后对照必须在同一枚台面上跑，母仓↔母仓、clone↔clone）。

### 3. 四条硬约束我核到的是**名册级**（⛔ 腿自勾；"有没有变松"那道必答归 `303-v1`）

- ① ⛔ 放宽断言＝**结构性成立**：`git show --name-only --format= e3341368` 十枚路径里 `cmd/wisp/` 只**一枚**（`panel_host_windows.go`，`+27/−0`），**⛔ 一枚测试件都没进名册** ⇒ 判据那三句所在的文件根本没被碰，⛔ 需要读 diff 才判得出。
- ② ⛔ 造 `--- SKIP`：clone 两发 SKIP＝3 → 2，且改后的 SKIP 名册是改前的**子集**（消失那枚＝`TestPanelHostLatencyPercentilesAC2`，负载敏感那族）。
- ③ ⛔ 碰 SLO／`thresholds.go`／D32：那 +27 行全落在一枚 JS 常量字符串内（`git diff --numstat cf46c24a e3341368 -- cmd/wisp/panel_host_windows.go`＝`27 0`）。
- ④ 落点＝门那一侧 ✓；⛔ 动 `scripts/**`／档／tag／ledger ✓（名册里根本没有那几枚路径）。
- ⚠ **两处我读了但明确⛔ 盖章**：**(a) ⛔ 扩权限**——我的判断是"门全在 Go 侧（`dispatchRaw` → `ComposerDispatch.Handle`、C17 名册、`ErrNoL2Confirm`、`perm.Options` nil 即 fail-closed），而 `window.external.invoke` 在修复前**就已可达** ⇒ 没新增能力"；★这句**⛔ 算凭据**，`303-v1` 的必答题。**(b)** 腿的门禁件 `31-package-after.txt` **末尾⛔ 落 `rc=N` 那一行**（尺＝`grep -n '^rc=' 0 命中`）＝踩我那条"每把门禁件自落一行 `rc=N`"，记档、⛔ 我顺手修（归口＝票 304 `AC#4` 那条编队模板族）。

### 4. 翻勾与现态

- **`AC#2` ✔**：三形分开答齐（甲＝因果对否证＋**具名欠**那枚门外 `Eval` 回读；乙＝代码＋因果＋HEAD 自洽三重否证；丙＝**活动形，有硬读数** `re-entered 9 levels / native exit 0 / door fired 0`，再叠加修复那发的成对反证）。本格文本自己写了"给不出就具名说缺哪枚＋谁来取"，腿照办 ⇒ 我判成立；⚠ 欠的那枚（门外 `Eval` 读 `typeof window.chrome.webview`／`document.title`）**具名交 `303-v1`** 裁"甲这一形在现有凭据下能不能判死"。
- **`AC#3` ⛔ 不翻**（★这一格⛔ 由我改定义后自己盖章）：我核到的是**事实层**——成对两发齐（撤修复⇒三枚逐字红 `nothing at all`／`cmp` 逐字节还原⇒nail1 绿）＋我母仓复跑同色＋上面那四条名册级约束；可本格**原文**要求的凭据是"三枚本机色全绿"，那句我已判定⛔ 成立（且⛔ 该由本票满足）。⇒ **我不把格子翻成"按更正后的目标成立"**，那枚更正后的目标交 `303-v1` 一次裁：ⓐ 这修复有没有把判据变松（腿自己声明⛔ 盖章）；ⓑ ⛔ 扩权限（我读过、给了判断，⛔ 算凭据）；ⓒ 按更正形状（nail1 绿＋两枚残余归票 305）本格算不算闭合。
- **`AC#4` ⛔ 不翻**：票面要"整包改前／改后**各两发**取交集"，我手上＝各**一发**＋一枚定向隔离发；而逐名那把尺给的是**＋1 枚新增红**，我用**另一把尺**（交替换头）判它⛔ 归因给修复 ⇒ 形状不满足，欠的那一发具名交 `303-v1`（⛔ 我用换尺去补主尺的缺）。
- **`AC#5` ⛔ 不翻**：① 票 302 甲／乙回填现在有料了——CI `test-windows` 在**修复的父发**那一发（run 38043392486）红名册 **20 枚**，含 nail1／nail2／`TestAC13…`／`TestPanelHostRealWindowHopAndLifecycle`（件 `probes/303/orch/r7-ci-before-red-roster.txt`；台面 `windows-2025-vs2026`／agent `2.337.0`，⚠ `ImageVersion` 那把尺⛔ 命中＝没取到，具名记档）⇒ 推送后那发作差才能答"这一跳在 CI 也回绿了吗"；② CI 色＝推送之后才有，本程就推。
- **现态＝3 勾／3 未勾**（`AC#0`／`AC#1`／`AC#2` 勾；`AC#3`／`AC#4`／`AC#5` 未勾），⛔ 改 `-done`。next＝推送取 CI 色 → `303-v1`（非实现者一次裁 `AC#3`：ⓐ 有⛔ 变松／ⓑ 有⛔ 扩权限／ⓒ 更正后的目标闭不闭合 ＋ 裁 `AC#4` 欠的那一发 ＋ 裁甲那枚欠读）→ `305-a1`。
- Progress log：
  - [2026-10-10 18:3x +0800] agent=编排者 did=收 `303-r1`（两笔，件 `probes/303/r1/**`，越界 0）⇒ 翻 `AC#2`、⛔ 翻 `AC#3`／`AC#4`／`AC#5`（★`AC#3` 那格我⛔ 按"更正后的目标"自盖章，交 `303-v1` 一次裁ⓐⓑⓒ）；★**销我自己票面那句"修后三枚全绿"**（clone 里回归之前那发拿到逐字同一句红 ⇒ nail2⛔ 是本票那一笔造成的）；★还腿具名欠的那格（clone 成对两发：268/17/3 ↔ 270/16/2，被修好 2 枚／逐名新增红 1 枚）；★那 1 枚新增红我用**交替换头**隔离发判成⛔ 归因给修复（同一枚二进制 47 秒内 5 绿→3 红；两枚 SHA 都出过两色），⛔ 顺手盖"新增红 0"；⚠ 另立缺口＝真窗焦点类判据本机 13 发 5 绿／8 红；★腿的门禁件 `31-` 末尾⛔ 落 `rc=N` 那行，记档⛔ 顺手修；next=推送（`A813`）→ `303-v1` → `305-a1`

## 编排者收 CI 改后那一发（2026-10-10 19:1x，件 `probes/303/orch/r9-ci-after-red-roster.txt`）⇒ **翻 `AC#5`**（两件都有盘上凭据了），`AC#3`／`AC#4` ⛔ 碰

### 0. 这一格原文要的两件，逐件答

- **② 本票那一发 CI 色**＝已取。改前发＝run `38043392486`／job `test-windows` `114187968428`／headSha `cf46c24a`（＝修复的**父发**）；改后发＝run `38045508579`／job `test-windows` `114194107792`／headSha `05db4bc6`。改后那发六枚 job 结论逐字（尺＝`gh run view --json jobs --jq`）＝`lint-frontend`=success／`slo-full`=success／`slo-smoke`=success／`test-windows`=failure／`lint`=failure／`test-core`=failure。
- **① 票 302 的甲／乙回填**＝**甲／乙都⛔ 落**，凭据从两枚来：本票 `AC#1` 已证"这一跳在**本机**都断"（⇒ 修完就该回绿，⛔ 是环境差），而这枚 CI 对拉把它**反过来**也证了——修复那一发在**托管 runner** 上拿到了页面自己的话（逐字：`AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`，`cmd/wisp/panel_resident_windows_test.go:825`）。票 302 票面那句"那台机器永远给不出回执"⛔ 成立 ⇒ 搬档（甲）／进 `-skip` ledger（乙）都⛔ 落；丙那一支剩下的射程只是"仍红的哪几枚、各归哪一票"，按下面第 3 节那份清单办。

### 1. 逐名作差（尺＝两份 `--- FAIL:` 顶层名 `sort -u` 后 `comm`，⛔ 比枚数；件＝`r9-ci-after-red-roster.txt`）

- 改前 20 枚／改后 20 枚 ⇒ **仍红 19 枚（逐字同名）／被修好 1 枚／新增 1 枚**。
- **FIXED**＝`TestAC14AwaitedBindingReplyReachesThePage`：改前 `--- FAIL (20.02s)`＋逐字 `panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)` ⇒ 改后 `--- PASS (0.30s)`＋上面那句页面原话。**判据是页面自己报回的**，⛔ 是 Go 侧自说自话（票 33 `A502` 形）。
- **NEW**＝`TestPanelHostLatencyPercentilesAC2`：改前 `--- SKIP (0.00s)`（`:985` 逐字 `no cold/hot sample recorded in this process … Named skip - an empty aggregate is not a green latency gate`）⇒ 改后 `--- FAIL (0.00s)`（`:1005` 逐字 `cold P95 2747.460 ms over 1 runs exceeds the D32 panel cold budget 1500 ms …`）。
- ⚠ **四数那把尺⛔ 能佐证"没变化"**：`cmd/wisp` 档逐字 `RUN=398 PASS=278 FAIL=8 SKIP=2` → `RUN=398 PASS=279 FAIL=8 SKIP=1`——`FAIL` 枚数**同数**而名集合作了差（＋1／−1），两把尺⛔ 同物（先例＝票 302 `:11` 那句同款坑）。

### 2. 本票边界段授权的那一枚"形状迁移"，CI 上如实拿到了（⛔ 更多）

- `TestPanelHostRealWindowHopAndLifecycle` 颜色**没变**（两发皆红），变的是形：改前 `:660` 逐字 `cold bring-up measured on this box: -1.000 ms`＋`:662` `did not produce a browser round trip (got -1.000); the message channel did not come up on the real window` ⇒ 改后 `:660` 逐字 `cold bring-up measured on this box: 2747.460 ms (HEAD 05db4bc at read time, 2026-10-10T10:41:20Z)`＋`:665` `cold bring-up 2747.5 ms exceeds D32 panel cold budget 1500 ms`。
- ⇒ 从**回执为零**移回**红于预算**，＝票面边界段写死允许的那一支，且⛔ 越界：⛔ 动 SLO 阈值／`thresholds.go`／D32 的 1500 那一枚数（一字未动，本程零产码）。
- ★⛔ 把这枚红读成"本票修复造成的退步"：票 302 现量补充节那张表里，**回归之前**的基线发（`cc315261`，10-06）同一枚就是红于预算（逐字 `cold bring-up 2889.2 ms exceeds D32 panel cold budget 1500 ms`）⇒ 2747.460 那一形是**旧形的回来**，⛔ 新坏。
- ⚠ 那枚 SKIP→FAIL 是上面这枚迁移的**次生色**（SKIP 的前提"本二进制里没有冷启样本"随修复消失，判据才开始求值）。它该归哪一票＝**⛔ 我裁**（本票只授权"移回红于预算"这一迁移，⛔ 顺手把预算面收进本票）；具名交 `303-v1`／票 302 的归口那格。

### 3. 票 305 那两枚症状在托管 runner 上一比一复现（换形而⛔ 换因）＝票 305 的台面⛔ 限于本机

- `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`：改前 `:325` 逐字 `no report "ac13-probe" … nothing at all`（`20.01s`）⇒ 改后 `:327` 逐字 `AC#13 page answer (head 05db4bc): "0" - 0 of 1 probe id(s) present in the live document`＋`:330` `… the round-trip probe page is the last document shown, so the user sees a stub instead of the panel.`（`0.62s`，同发逐字 `wisp: panel window is up (test-harness, cold 470.2 ms, hot path 0.0 ms)`）
- `TestAC14GoSideEvalPushReachesThePage`：改前 `:866` 逐字 `no report "ac14-push" … nothing at all`（`20.02s`）⇒ 改后 `:867` 逐字 `AC#14 nail 2 (Eval push hop), page's own words: title=""`＋`:869` `Go's Eval push did not reach the document … This is the push dimension, separate from the awaited-reply dimension …`（`0.37s`）
- ⇒ 与我本机 clone 台件拿到的是**同一句逐字**（件 `r5-r6-orch-pair-verdict.txt`／`52-restore-analysis-and-topback.md` 那两枚形状）。⚠ 一处台面差要具名：改前两枚各走满 15s 超时（20.01s／20.02s），改后 0.62s／0.37s——时长差本身是"回执到了而内容不对"的证据，⛔ 当分母用。**派票 305 的腿时必须把这两条逐字带进派单。**

### 4. ★台面那格：`r7` 里我具名欠的那枚读数，这次取到了，而⛔ 取不到是我那把尺写坏的

两发同一 job 的 `Set up job` 段逐字：`Image: windows-2025-vs2026`／`Version: 20260925.250.1`／`Current runner version: '2.337.0'` ⇒ **同档同版本**，"镜像漂"那一支在这枚对拉里⛔ 有依据。
`r7-ci-before-red-roster.txt` 里我写的"`ImageVersion` 这一把尺⛔ 命中"＝**我的尺的缺陷**：我在日志里找的是 GitHub API 的字段名 `ImageVersion`，而日志里那行写在 `##[group]Runner Image` 下面、叫 `Version:`。**记我**；定式并回＝**同一枚事实在两种载体上名字⛔ 相同（API 字段 ↔ 步日志文本），换载体取读数要先现量那一行的字面写法，⛔ 把"我这把 grep 没命中"写成"盘上没这行"**。

### 5. 翻勾与现态

- **`AC#5` ✔**：本格文本自己写死了"⛔ 腿做、交回后由编排者办的两件"⇒ 两件都办了且各带盘上凭据（①＝票 302 现量补充节＋本件第 0 节那条页面原话；②＝本件第 1 节）。⛔ 实现者＝我（修复那一笔的作者是 `303-r1`），⛔ 角色撞。
- **`AC#3` ⛔ 动**：仍交 `303-v1` 三道必答（ⓐ 有⛔ 变松／ⓑ 有⛔ 扩权限／ⓒ 按更正形状本格闭不闭合）。CI 那一发的逐字红句**正是**ⓐⓑ 要的对照料，已写进派单。
- **`AC#4` ⛔ 动**：票面要"整包改前／改后**各两发**取交集"，欠的那一发仍然欠（我这程做的是 CI 档的另一把尺，⛔ 用换尺去补主尺的缺，先例＝我上一条定式）。
- **现态＝4 勾／2 未勾**（勾＝`AC#0`／`AC#1`／`AC#2`／`AC#5`；未勾＝`AC#3`／`AC#4`），⛔ 改 `-done`。
- next＝**`303-v1`**（非实现者，锚 `9cef4589` 之后我这一笔）一次裁 `AC#3` 三道＋`AC#4` 欠的那一发＋甲那枚门外 `Eval` 欠读 → **`305-a1`**（票 305 落点只读，带上面第 3 节两条逐字）→ 票 302 的归口那格（含本件第 2 节那枚次生色）。
- Progress log：
  - [2026-10-10 19:1x +0800] agent=编排者 did=取 CI 改后那一发并与 `r7` 改前名册逐名作差 ⇒ **翻 `AC#5`**（① 票 302 甲／乙**都⛔ 落**，凭据＝托管 runner 上页面原话 `REPLIED,REPLIED,REPLIED`；② 四数＋名册作差已交，件 `r9-ci-after-red-roster.txt`）；★FIXED 1 枚（nail1，红句 `nothing at all`→页面原话）／NEW 1 枚（`TestPanelHostLatencyPercentilesAC2` 由具名 SKIP 走进 FAIL，`cold P95 2747.460 ms` 超 D32 的 1500）／仍红 19 枚逐字同名；★本票边界授权的那枚"回执为零→红于预算"迁移**如实出现**且⛔ 是新坏（基线发同枚逐字 2889.2 ms 同形），⛔ 动预算；★票 305 两枚症状在 CI 一比一复现（逐字已在件里）⇒ 票 305 台面⛔ 限本机；⚠ 四数尺 `FAIL=8↔8` 与名册作差⛔ 同物；★销掉 `r7` 里我那格"ImageVersion ⛔ 命中"＝**我那把 grep 找的是 API 字段名而日志那行叫 `Version:`**（两发同 `20260925.250.1`，镜像漂那一支⛔ 成立），定式并回＝换载体取读数先现量那一行字面写法；`AC#3`／`AC#4` ⛔ 碰；next=`303-v1` → `305-a1`

## 编排者收 `303-v1`（2026-10-10 21:4x，四笔 `8fc42230`→`d9864baf`→`2d68e28d`→`49f7e371`，件 `probes/303/v1/**`；越界＝我自己逐笔取名册重算＝**0 枚**、`git diff --name-only origin/dev..HEAD -- '*.go'`＝**0 枚**）⇒ **翻 `AC#4`**＋**补一枚可满足的 `AC#3b` 并翻它**、★**`AC#3` 原句⛔ 翻（那枚凭据句在本票里从来⛔ 可能成立，原句⛔ 删，只追加本节）**＋★**它交回的整包名册与我 18:1x 自己那对逐名对拉上了：改后那一侧⛔ 一字之差、改前那一侧只差一枚已具名的 flake**

### 0. ★两把独立名册对拉（本程最要紧的一把尺：⛔ 信任何一方的转述）

- 我 18:1x 那对（件 `probes/303/orch/r5-orch-package-before-names.txt`／`r6-orch-package-after-names.txt`，台面＝仓外 clone，各**一发**）vs 它这对（`logs/names-{before,after}-{1,2}.txt`，台面＝另一枚仓外干净 clone `~/wisp-303-v1`，各**两发**）。尺＝剥 `--- FAIL:` 前缀与时长后 `sort`，`diff`／`comm` 比**集合**、⛔ 比枚数：
  - **改后那一侧＝`IDENTICAL_SETS`**：它的"两发改后取交集"16 枚与**我那一发改后**的 16 枚**逐名同集合**（无一枚只在其中一边）。
  - **改前那一侧只差一枚**＝`TestAC4FocusReturnToPriorWindowGap33r5`（它两发**都红**、我那一发**恰好绿**）。⇒ 这正是我 A823 里写过的那枚**非确定性真窗焦点仪**（台账既有先例：本机 13 跑 5 绿／8 红）。
  - ⇒ ★**我 A823 那句"逐名新增红 1 枚"在这里被我这把尺自己收回一半**：那一枚当时是 flake 撞在我那一发上；把它按"各两发取交集"这把尺放平＝**新增红 0 枚**。**枚数没变、结论换了＝换尺的效果**，⛔ 写成"它算对了而 I 算错了"，两把尺各带各自射程（我＝各一发、它＝各两发取交集）。
- ⚠ **这把"交集尺"的盲区必须具名（⛔ 让本格被引成"从来⛔ 新增红"）**：取两发交集＝**只红过一次的枚会被排除**。它的 `after-1` 里就多红过一枚 `TestTicket223HandEditedFsLooseningCostsAnL2Card`（`after-2` 没有）。⇒ 本格成论的射程＝**"稳定新增红 0 枚"**；那一枚单发红我⛔ 归因给本票那笔修复（它属 FS／L2 时序一族，台账 `:10553` 有过同一枚 SHA 两发一红一绿的决定性读数），登记成 flake 形状残余。

### 1. 四条硬约束——本程**我自己**在 `git diff cf46c24a 807497c1` 上现跑（⛔ 抄腿的 ⓐ 判语）

| 硬约束（票面原句） | 我这把尺 | 读数 |
|---|---|---|
| ① ⛔ 放宽任何断言（那三句 `no report … nothing at all` 本体一字⛔ 动） | `git diff cf46c24a 807497c1 \| grep -c '^-[^-]'` | **0**（整个区间**零删行**，唯一产码文件＝`cmd/wisp/panel_host_windows.go`，＋27 行全在新增的注释块＋`external.invoke` 包装里） |
| ② ⛔ `t.Fatalf`→`t.Skip`、⛔ 造 `--- SKIP` | `git diff … -- cmd/wisp/panel_host_windows.go \| grep -E '^+' \| grep -cE 't\.Skip\|SKIP\|budget\|1500\|thresholds'` | **0**（diff 里那几处 `SKIP`／`no report` 命中全在 `probes/303/r1/**` 的**日志文本**里，⛔ 产码行） |
| ③ ⛔ 动 SLO 阈值／`thresholds.go`／D32 那一面 | 同上那一把＋名册 | 产码只碰一枚文件、零 `budget`/`1500` 新增；`internal/observe/thresholds.go` **0 字节** |
| ④ 改动落在**门**那一侧、⛔ 顺手改档／tag／ledger | `git diff --name-only cf46c24a 807497c1` | 只 `cmd/wisp/panel_host_windows.go`＋自家 `probes/303/r1/**`；`scripts/portable-tests.sh` **零字节**＝⛔ 碰票 302 射程 |

### 2. ★`AC#3` ⛔ 翻，⛔ 是因为修复没落地——是因为**那半句凭据在本票里从来⛔ 可能成立**

- 原句要"成对两发＋**三枚本机色**（修后 `TestAC13`／两枚 `TestAC14*` 全绿）"。现量（三边独立）：修后 `TestAC14AwaitedBindingReplyReachesThePage` **绿**（页面原话 `REPLIED,REPLIED,REPLIED`，CI 与本机两侧都有）；`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 在干净 clone 里**恒 SKIP**、在母仓**恒红**（红句＝页面自己答 `"0"`）；`TestAC14GoSideEvalPushReachesThePage` **同一句红在它回归之前那发就在**（`f718e9b6`，腿件 `logs/nails-at-f718e9b6.txt` ＋**我 18:0x 自己在同一枚 clone 上复跑过那一发**，件 `probes/303/orch/r3-f718e9b6-three-nails.txt`）。
- ⇒ 后两枚**⛔ 是本票那一笔造成的**，它们是**票 305** 的两枚症状（我 18:4x 立票时就是按这个判的）。所以"三枚全绿"这半句⛔ 满足得了一个还没被做过的功能——**本格⛔ 翻、⛔ 删原句**，可满足的那部分搬到下面新格。
- ⚠ 腿自己那条判语（ⓣ＝凭据句本身有缺陷／⛔＝修复⛔ 成立／ⓤ＝更正形状下闭合，挂 ⓓⓖ 两条件）我认；ⓓ（各两发补齐）＝本程第 0 节已独立对拉过、ⓖ（另立票 305）＝票已在盘上 ⇒ **两个条件都闭了**，闭合形状落成 `AC#3b`。

- ⇒ 可满足的那部分搬到票面 AC 清单里的新格 **`AC#3b`**（紧接 `AC#3` 之后，勾已翻；判据正文与三条射程＝本节，⛔ 重复一遍造成第二份真相源）。

### 3. `AC#4` 翻勾的凭据与⛔ 它覆盖的三样

- 门禁逐把我自己核过（件 `probes/303/v1/logs/`）：`go vet ./cmd/wisp/` rc=0、`go vet -tags winlive ./cmd/wisp/ ./internal/ball/` rc=0、`sh scripts/d22scan.sh` rc=0、**`gofmt` 并排两把**＝工作树 `cmd/wisp` **1 枚**（`models.go`，既有加严残留、⛔ 顺手修，先例＝票 298）／HEAD blob 导出树 **0 枚**。四发同一枚 clone、只差一次 `git checkout --detach`（⛔ 母仓工作树）。
- ⚠ 一枚**我派单里写错的尺**要更正：我给的"全仓 `gofmt` 大约 5 枚"那把尺对不上它现量的**38 枚**。它**报差⛔ 自行调和＝对**。⇒ 差额属**射程／台面**那一族（工作树 vs HEAD blob／CRLF），**归我下一程现量核对并具名改正**（记我：本票与票 296 那两处派单都引过"5 枚"那句）。
- ⛔ 本格⛔ 覆盖的三样，逐枚留在票面与台账、⛔ 塞进翻勾句：① 腿的 **M1–M4 变异取证**＝〔仅腿报，我⛔ 复跑〕；★其中 M3 那一形的台件在 `logs/phase2-out.txt` 里**中途撞上一行 Python `SyntaxError: unmatched ')'`**（它随后自陈"M3: fix removed as well"＝机制那一发⛔ 干净），那部分判语我⛔ 盖章。② 它自陈的 **M4 网眼**＝拆 Go 侧帧那一形十枚钉全绿＝今天⛔ 任何仪器看得见，欠补钉设计。③ 它自陈**母仓台面 lifecycle 色未复测**。

### 4. ★那枚超预算红的归口（`AC#5` ②那格⛔ 答的一半，本程裁掉）

- 形状：CI 托管 runner 上 `cold 2747.460 ms` 超 D32 面板冷启预算 `1500` ⇒ `TestPanelHostRealWindowHopAndLifecycle` 与 `TestPanelHostLatencyPercentilesAC2` 各红（A824 §"本票有关那四枚"逐字）；而**改后本机同一枚码⛔ 红**——腿报的三发 `cold 980.666／1002.355／1341.042 ms` 我⛔ 复跑，但**我自己那对名册（第 0 节 `r6`）里 lifecycle 压根不在红侧** ⇒ "本机不红"这一支我有凭据，"具体那三个毫秒数"＝〔仅腿报〕。
- **裁语**＝⛔ 动那枚 `1500`（D32 面，人工批准才动）、⛔ 为它另立票、**归口＝票 302 的"runner 时序 vs 我们的读数"那一面**（先例＝`pending-and-issues.md:10553` 那枚同 SHA 两发一红一绿的决定性读数）。⇒ 这条从 A824 挂着的"未裁"里销掉；票 302 那格保持"按住＝⛔ 做完了"。

### 5. 现态与 next

- 票 303 现态＝**6 勾／1 未勾**（勾＝`AC#0`/`AC#1`/`AC#2`/`AC#4`/`AC#5`/`AC#3b`；未勾＝⛔ 满足的 `AC#3` 原句那一格，理由见本节第 2 节）⇒ **⛔ 改 `-done`**：收口那一步（全勾→`Status: done`→改名→更新索引）等我把 `AC#3` 这枚"死凭据句"按票 141 那枚不可满足判据的先例正式登记完再办，⛔ 顺手盖。
- next＝**派 `306-a1`**（票 306 的 `AC#0` 三张表，只读）→ **`305-a1`**（票 305 的归因，派单必带票 305 现量那两句逐字）→ `296-r2` → `295-r1` → `294-r1` → `294-v1` → `298-v1`；我这边两笔小账＝**现量核对 `gofmt` 那 38 枚的射程并具名更正派单口径**、**M4 网眼那枚欠账具名入票面**。
- Progress log：
  - [2026-10-10 21:4x +0800] agent=编排者 did=收 `303-v1`（四笔，越界逐笔重算 0、`-- '*.go'` 0 枚；它误名 `logs/phase2.out` 那件我 `cmp` 补件＝逐字节 SAME 292 字节）＋★**两把独立名册对拉**（它的"各两发取交集"vs 我 18:1x 的"各一发"：改后那一侧 `IDENTICAL_SETS`＝16 枚逐名同集合；改前那一侧只差 `TestAC4FocusReturnToPriorWindowGap33r5`＝已具名 flake，它两发都红而我那一发恰好绿）⇒ **我 A823 那句"逐名新增红 1 枚"换尺后＝0 枚**，同时把"交集尺会藏单发红"这枚盲区写死（`after-1` 多红过一枚 `TestTicket223HandEditedFsLooseningCostsAnL2Card`，⛔ 归本票修复，台账 `:10553` 有同 SHA 一红一绿先例）；★**四条硬约束我在 `git diff cf46c24a 807497c1` 上自己现跑**（整个区间**零删行**、产码唯一文件 `cmd/wisp/panel_host_windows.go` 的＋27 行里 `t.Skip|SKIP|budget|1500|thresholds` 命中 **0**、`scripts/portable-tests.sh` 零字节）；⇒ **翻 `AC#4`**＋**补 `AC#3b`（可满足形状）并翻它**；**`AC#3` ⛔ 翻**＝它那半句"三枚本机色全绿"要的是票 305 那两枚症状（`AC13` 干净 clone 恒 SKIP／母仓恒红、nail2 回归前同句已红＝我自己 18:0x 在 clone 上复跑过的 `r3-f718e9b6`），⛔ 删原句、只追加本节；★**裁掉一枚挂着未裁的归因**＝CI `cold 2747.460 ms` 超 D32 `1500` ⛔ 动预算、⛔ 另立票、归口票 302 的 runner 时序那一面（"本机不红"我有我自己 `r6` 的凭据，三个毫秒数＝〔仅腿报〕）；⛔ 给腿的 M1–M4 变异判语盖章（M3 台件 `logs/phase2-out.txt` 中途撞 Python `SyntaxError`、M4 十枚全绿那枚网眼欠补钉、母仓 lifecycle 色未复测，三样各自留档）；⚠ 我派单"全仓 gofmt 约 5 枚"对不上它现量 38 枚＝**报差不自调、对**，差额归我下一程核并更正两处派单口径；现态 6 勾／1 未勾（`AC#3` 那枚死凭据句），⛔ `-done`（收口等按票 141 先例正式登记不可满足判据）；next=派 `306-a1` → `305-a1` → `296-r2`
