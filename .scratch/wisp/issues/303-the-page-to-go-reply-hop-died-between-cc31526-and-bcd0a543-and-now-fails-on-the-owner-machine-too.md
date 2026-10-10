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
- [ ] **`AC#2` 机制层（三形分开答，⛔ 合一格）**：**(甲) 页面 JS 压根⛔ 跑**／**(乙) JS 跑了但绑定名⛔ 接上**（`wispDispatch` 那枚注入名）／**(丙) 钩子接上了但门被再入帽挡住**。每形必答两问：① **能把这一形与另两形分开的读数长什么样**（现量给得出就给，给不出就具名写"缺哪枚读数＋谁来取"）；② 本票字节到不到得了那一步（`302-a1` 已判：这三形在现有 CI 字节里**分⛔ 开**，`-1.000` 与"零回执"**同时出现**）。⚠ 判"产品路径也断"必须⛔ 依赖这三形里任何一形的**假设**，要一发真机或一条真调用链读数。
- [ ] **`AC#3` 修复（只在 `AC#1`/`AC#2` 交完之后开工）**：**硬约束四条**——① ⛔ 放宽任何断言：`no report … within 15s (what DID arrive at the door: nothing at all)` 那三句判据本体一字⛔ 动；② ⛔ 把 `t.Fatalf` 换成 `t.Skip`、⛔ 造 `--- SKIP`（`tools/d22scan/runtests.sh:98`→`:102` 把任何 `^--- SKIP` 判红，那是**故意的**）；③ ⛔ 动 SLO 阈值／`thresholds.go`／D32 那一面（`TestPanelHostRealWindowHopAndLifecycle` 基线红在**预算**那一支属 SLO，⛔ 用"改预算"换绿）；④ 产码改动落在**门**那一侧（`cmd/wisp/**` 宿主传输边），⛔ 顺手改档／tag／ledger（那是票 302 的射程，本票⛔ 动 `scripts/portable-tests.sh`）。凭据＝成对两发＋三枚本机色（修后 `TestAC13`／两枚 `TestAC14*` 全绿）。
- [ ] **`AC#4` 门禁与越界**：`go vet ./cmd/wisp/` rc=0；`go vet -tags winlive ./cmd/wisp/ ./internal/ball/` rc=0；`sh scripts/d22scan.sh` rc=0；`gofmt -l cmd/wisp` 与 HEAD blob 那把**并排两把**都要报（先例＝票 298：加严残留具名、⛔ 顺手修）；`cmd/wisp` 整包改前／改后各两发取交集 ⇒ **新增红 0 枚**（⚠ 本机整包要铺 sherpa `PATH`，缺 DLL 是 `0xc0000135` 且**无 `--- FAIL`**＝用例根本没跑）；`git show --name-only --format=` **逐笔**名册＋与授权名册差集（⛔ 区间尺）；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec 且写在 `$( … )` **之外**。
- [ ] **`AC#5` 交回后由编排者办的两件**（⛔ 腿做）：① 票 302 的甲／乙裁语按 `AC#1` 结论回填（若"这一跳在本机都断"⇒ 修完就该回绿 ⇒ **甲／乙都⛔ 落**）；② 本票那一发 CI 色（修后推送产生的 `test-windows` 四数与红名册）。

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
