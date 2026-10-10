# 票 305 — 面板**冷启动之后"活文档"到底是谁**：通道已经修好（页面能答了），但内嵌入口没交接上来，而 `Eval` 推不进文档那一枚**早在回归之前就红、从来没绿过**

**立票**：2026-10-10 18:0x 编排者（来路＝票 303 修复腿 `303-r1` 交件时**顶回我票面那句"三枚回绿"**；我随后自己跑了三发把它证成）
**性质**：★**缺陷票**，两枚症状、**⛔ 一条因果**（`AC#0` 必须先分开，混一票必然顺手改错那半）。
**为什么要紧**：这就是用户真正会碰到的那一屏——**面板打开了，但里面不是给他看的那页**。通道修好之后红句从"页面没答话"换成了"页面答了话、可它答的是探针那张文档"，⇒ **症状从"链路断"升级成"链路通而内容错"**，这对后面每一张要往面板里放东西的票（186／187／167／145）都是前提。

## 现量（每条带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**母仓 @ `807497c1`（修复已在树里）**，尺＝`go test ./cmd/wisp/ -count=1 -timeout 420s -v -run '<三枚名>'`，件 `probes/303/orch/r2-orch-targeted-after.txt`（我 18:0x 现跑，`rc=1`）：
  `TestAC14AwaitedBindingReplyReachesThePage` ⇒ **PASS(0.55s)**，页面自己的话逐字 `REPLIED,REPLIED,REPLIED`（`Go's handler was reached by 3 of the 3 real requests`）。
  `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` ⇒ **FAIL(1.30s)**，红句逐字 `after a real cold start the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0")`（出处 `cmd/wisp/panel_resident_windows_test.go:330`）。
  `TestAC14GoSideEvalPushReachesThePage` ⇒ **FAIL(0.54s)**，红句逐字 `Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK"`（`:869`）。
- ★**同一枚仓外 clone 上的前后对照**（这是本票最关键的一把尺：台面相同、只有 commit 不同）：
  `f718e9b6`（＝回归笔 `fb2fb802` 的父）⇒ 件 `probes/303/orch/r3-f718e9b6-three-nails.txt`：nail1 **PASS**、**nail2 FAIL 且红句与今天逐字同形**（`title as "", want "PUSHED-33R5-OK"`）、AC13 **SKIP(0.00s)**。
  `807497c1`（修复后）⇒ 件 `probes/303/orch/r4-clone-at-HEAD-nails.txt`：nail1 **PASS**、nail2 **同一句 FAIL**、AC13 **SKIP**。
  ⇒ **判语①**：nail2 的那枚红**⛔ 是 `fb2fb802` 造成的**——它在回归之前的同一台面上就是同一句话。**"Eval 推不进文档"是一枚独立的、更早的缺陷**（本票症状②）。
  ⇒ **判语②**：AC13 在这枚 clone 里两发**都 SKIP** ⇒ 它的颜色**只有母仓台面能量**（先例＝`frontend/dist` 只跟踪 `.gitkeep`，本机有旧产物、干净 clone 没有）⇒ 本票⛔ 拿 clone 的 SKIP 当"它本来就好"，也⛔ 拿它当"它坏了"。
- ★**母仓改前那一发**（腿件 `probes/303/r1/20-targeted-before-red.txt`，`rc=1`）：三枚全 FAIL 且都是 `no report … from the page within 15s (what DID arrive at the door: nothing at all)` ⇒ 通道那一层确实被 `fb2fb802` 断过、也确实被 `303-r1` 接回来了（我这发的 PASS 是凭据）。⚠ **AC13 在母仓⛔ 已知"改前绿"过**（改前那一发它红的是**另一句**，因为通道断在前面挡住了它）⇒ 归因欠着，见 `AC#0` ②。
- **腿具名的嫌疑笔（〔仅腿报〕，我⛔ 复跑过归因）**：`70b00885`（`bringUp` 把探针文档当最后活文档那一族的引入处）。⛔ 把它当结论引。

## 要建什么（`AC#0` 之前⛔ 任何产码）

- [ ] **`AC#0` 只读归因＋分层普查**：交回四张表——
  ① **"活文档"名册**：`bringUp` → `firstRoundTrip` → `serveEntry` 这条装配序里，每一次 `SetHtml`/导航各自把哪份文档变成活的、最后一次是谁（尺＝逐枚调用点文件:行＋调用序；⛔ 按注释读）。
  ② **两枚症状分开归因**（各用一枚**同台面**的前后对照，⛔ 拿 clone 的 SKIP 充数、⛔ 拿不同台面相减）：症状①＝AC13（入口没交接）；症状②＝nail2（`Eval` 推不进）。⚠ 若某一枚在**能判死的台面上从来没有绿过**，就写"它可能⛔ 是一枚回归、而是一枚从没做通过的功能"，并具名说凭据是哪一发。
  ③ **判据是不是同一条**：nail1 绿／nail2 红 是同一枚用例文件里的**两个维度**（awaited reply vs Go-side push，票 33 `A475` 那批裁过"两维⛔ 混一枚"）⇒ 本票必答"push 那一维今天缺的是哪一跳"。
  ④ **代价表**：修①要动 `bringUp` 的文档交接次序（**第二枚未归因的产码改动**）；修②要动 `Eval` 那一路（今天是不是根本没人调用它？先跑调用点尺）。＋**必答**：动②会不会撞上票 33 已裁的那条"面板线程＝专用 STA ＋ 库 `Run()` 泵"（`Eval` 只能在泵线程上跑）——撞就具名报回、由我裁，⛔ 自己绕。
- [ ] **`AC#1` 症状①（入口交接）落地**：只在 `AC#0` ② 归因之后开工。**硬约束**＝⛔ 放宽 `TestAC13…` 那句判据（"live document 必须含内嵌入口声明的那些 element id"是票 33 立的行为判据，⛔ 改成"含任一 id"或改成 `t.Skip`）；⛔ 为了让 AC13 绿而动 `firstRoundTrip` 探针文档的存在性（探针那发⛔ 是"最后活文档"这件事本身就是要修的对象）；⛔ 动 `frontend/src/**`（页面源码不在本仓射程，读一律 `git show <ref>:<path>`）。
- [ ] **`AC#2` 症状②（`Eval` 推不进）落地**：同上，⛔ 放宽"页面自己报 title"那句判据（判据必须**由页面报回**，先例＝票 33 `A502` 形）。若 `AC#0` 判出"这一维从来没有绿过 ⇒ 是缺功能⛔ 是回归"，本格**自动降级为⛔ 开工**，由我另立落地票。
- [ ] **`AC#3` 反恒真成对两发**：每枚症状各一对＝(a) 撤掉修复 ⇒ 目标用例当场按预期红（**具名红句逐字引**，⛔ 只报枚数）；(b) `cmp` 逐字节还原 ⇒ 复绿。⚠ 外加一发**跨维度负控**：修①的那笔⛔ 可能让 nail2 变绿——若变绿，必须具名说明为什么（否则就是"两枚症状其实同源"，本票的分层假设作废）。
- [ ] **`AC#4` 门禁与越界**：`go vet ./cmd/wisp/` 与 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/` 各 rc=0；`sh scripts/d22scan.sh` rc=0；格式名册**两把并排**（工作树／HEAD blob，各自写射程）；**整包 `go test ./cmd/wisp/ -v` 改前改后各一发取交集＝新增红 0 枚**（⚠ 两发必须**同一台面**，见下面"规矩"里那条）；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；每把门禁件自落一行 `rc=N`；⛔ 零 push；commit 显式 pathspec 写在 `$( … )` **之外**、只落自家 `probes/305/<腿名>/**`（票 301 `AC#4b`）。
- [ ] **`AC#5` 与票 33 的对账**：本票收掉的两枚症状对应票 33 的 `AC#13`／`AC#14` 哪几格，逐枚写"⛔ 翻票 33 那格／翻本票本格"，⛔ 两本账同时记同一件事（先例＝`A817`/票 302 那处"同一物理缺陷链上只记一次"）。

## 边界（⛔ 塞进本票）

- ⛔ **`fb2fb802` 那一跳的回执通道**＝票 303 已闭（`AC#3` 的修复），本票只在它之上判"内容对不对"。
- ⛔ **`--scope=cli` 里这三枚该不该搬档**＝票 302（现在**仍未裁**，且它的甲／乙裁语本来就按在本票之后）。⚠ 顺序⛔ 倒：先修内容，再谈搬档——搬档会把"链路通而内容错"这枚进步藏起来。
- ⛔ **`msedgewebview2.exe` 那 18/24 枚**＝机主自己的应用（`A821`＋`A823`），⛔ 残留、⛔ 本票射程。
- ⛔ **`internal/risk` 那 12 枚 CI 红**＝`A817` §3〔待归因〕。

## 规矩（本票全程）

- 子代理**只 commit、不 push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ 在仓内建 worktree 或 checkout；临时件**只建不删**；证据件⛔ 叫 `.out`。
- ★**台面规则（本票立的一条新规矩，来自我这次自己抓到的坑）**：**前后对照必须在同一枚台面上跑**（母仓↔母仓、clone↔clone）。母仓带 `frontend/dist` 旧产物、干净 clone 没有 ⇒ 同一枚用例会给出**不同颜色**（这次实测：AC13 在 clone 两发都 SKIP、在母仓是 FAIL）。⇒ 凡跨台面比较，读数⛔ 可比，必须具名写成"两把台面"。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件；⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP`。
- 跑 `cmd/wisp` 测试必带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（★shell 形 `/d/…`，`D:/…` 那形 `0xc0000135` 且零 `--- FAIL`）；跑真机/整包前 `tasklist` 现量 `wisp.exe`／`balldebug.exe`＝0。
- 裁决者≠实现者（`SPEC-12 §4.3`）：`AC#1`/`AC#2`/`AC#3` 的判语⛔ 由落地腿自勾。
- ⚠ 腿给的嫌疑笔 `70b00885` 是〔仅腿报〕，本票 `AC#0` ② 未跑之前，任何文案⛔ 把它当归因结论写。

next=`303-v1`（非实现者裁票 303 `AC#2`/`AC#3`，含"这修复⛔ 放宽判据"那道必答题）→ `305-a1`（`AC#0` 只读，⛔ Go 编译面按补位腿规矩自报；若它要"改前必红"的夹具色，需先与编排者抢台面）→ 我裁 → 视 `AC#0` ② 的结论决定 `305-r1` 开⛔ 开
