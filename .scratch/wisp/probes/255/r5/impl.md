# 票 255 AC#4 · 写腿 `255-r5` 交付件 — AC#4 收尾：把 `TestPanelHostRealWindowHopAndLifecycle` 修成确定

状态：**已完成**（§0–§6 全部写满；门禁四数齐；我的格子＝确定绿；整包 22 枚红全部在 `261-r2` 爆炸半径、逐字登记于 §4。交件 commit 号见 §0）。

---

## 0. 起手锚

| 项 | 读数 |
|---|---|
| 起手时刻 | `2026-10-03T15:56:38+08:00`（`date -Iseconds` 自取） |
| 起手 HEAD | `6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（`git log -1 --format=%H` 自取） |
| 分支 | `dev` |
| 起手写面态 | `git status --porcelain cmd/wisp internal` ＝ **空**（0 行）；会话中途 `261-r1`（`e2bfd7c8`）与 `212-a2`（`66fc537f`）落库，均在 `internal/**` 与 probes，不碰我写面；台账 `89c863f1` 明写「cmd/wisp 写面队列压在 255-r5 之后」⇒ 全程我是 `cmd/wisp` 唯一写手 |
| 我的第一 commit | 本节 commit 后补记（见 §6 末） |
| 票面 | `.scratch/wisp/issues/255-*.md`（AC#4 与「★ AC#4 的撞钉预检＋两格裁死」节读过；⛔ 未碰任何勾框） |
| 证据件 | 本文件＋`.scratch/wisp/probes/255/r5/` 下 8 份原始读数（只建不删） |

派单引的读数与我对库面的复核：

| 派单说 | 我复认 |
|---|---|
| r4 终态 1 红，红句在 `panel_host_windows_test.go:649` 一带 | 红名确是 `TestPanelHostRealWindowHopAndLifecycle`（`gate-final-verbose.txt:687`），但**逐字红行是 `:585` 的 `t.Errorf`**（"hide -> re-show changed the msedgewebview2 set ... 8 -> 6"）；`:649` 是 `t.Logf`（"ONE sample and no timing claim"），**它在 PASS 发里同样出现**（`r3/final-verbose-2.txt:686`，1.24s PASS） |
| `255-r3`（13:34 那发）同用例 PASS | 成立：`final-verbose-2.txt:687` `--- PASS (1.24s)` |
| 两发之间产码没变 | **对**：`git diff 7cb9d570..081019c0 -- cmd/wisp/panel_host_windows_test.go` ＝ 空（exit 0 无输出）。派单说"产码"指 `panel_host_windows.go`／`panel_resident_windows.go`——**这两个文件恰在同一对 commit 之间变了**（+108/−12、+51 行，＝`4658dbb6` 的 AC#4 取值闭包），但变的是建窗几何来源，与 hide→re-show 的子进程计数路无涉（详见 §5 推翻第 2 条） |

写面名册遵守：全程只改 `cmd/wisp/panel_host_windows_test.go` 一枚；`frontend/**`／`design/**` 未读未引；`docs/**`、`thresholds.go`、golden、`allowlist.txt` 未动；票面勾框未碰。

---

## 1. 我复跑的读数（修尺之前，`-count=3`）

命令（照派单逐字）：
`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=3 -run 'TestPanelHostRealWindowHopAndLifecycle' ./cmd/wisp/`
原始输出：`.scratch/wisp/probes/255/r5/rerun-count3.txt`（HEAD 读数 `1ba22d25`＝当时工作树顶，261-r1 在飞）。

| 发 | 终态 | 逐字红句 |
|---|---|---|
| 1 | FAIL 2.16s | `:549` cold 1689.5ms 超预算（连带）＋ **`:585` hide→re-show 树内计数 8→6** |
| 2 | PASS 1.09s | —（单发采样抓到 7=7） |
| 3 | FAIL 0.68s | **`:585` 树内计数 11→10** |

**＝ 2 红 1 绿。** 不是 r4 那种"3 发里 0 红"的侥幸形状——在我这台机器上它是**高频复现**（r4 终态 1 红、我 2 红），间歇面成立但概率远高于"偶发"。

### 1a. 判因四步（探针链，全部只加 `t.Logf`、断言未动，探针已撤）

1. **恢复曲线**（`probe-curve.txt`，`-count=2` 两发皆然）：re-show 后树内计数 `7→6` 后**整 4 秒稳定停在 6**，不回升。
2. **不 hide 对照**（`probe-nohide.txt`）：不 hide 连续采 6 秒＝`6,6,6,...` 纹丝不动 ⇒ "计数自己会掉"不是稳态背景；掉落与"建窗后数秒"这个时段绑定，与 hide 无关。
3. **名单转储**（`probe-names2.txt`，`-count=4`，2 红 2 绿）：FAIL 发里 pre/post 两侧名单**逐枚对照**——浏览器主体（ppid=测试进程那枚）与其 5 枚 helper **跨 transition 逐枚存活**（pid 全同）；唯一退掉的是那枚**每发换 pid 的短命 helper**（26636→16000→29996→31268）。**浏览器主体从未死过、HWND 全发同一** ⇒ "single-window reuse broken" 是误报：真相＝**建窗后 ~1s 内 WebView2 自行 spawn 然后退出的一枚短命 helper**，把"树内计数严格相等"这把单发尺打翻。
4. **等值沉降对照**（`post-fix-count3.txt`）：把读法改成"两侧各等两个一致样本再比"（等值沉降）后，红从间歇变**确定**（2/3 稳定 7 settled → 6 settled）⇒ 掉落是真实进程退出、不是抖动 ⇒ **等它回来这条路判本上就是死的**，判据必须换。

### 1b. 结论

派单的"初判＝间歇形（Destroy 后立刻 readWebviewTree，webview 子进程异步退出还没退干净）"**方向对（是异步退出）、但归因错（不是 Destroy 的子进程没退干净，而是 bring-up 的一枚短命 helper 自然退场；且此用例里 hide→re-show 根本不 Destroy）**。r3 PASS／r4 FAIL 两发的差异＝单发采样与那枚 helper 退场时刻的赛跑，不是产码差异。

---

## 2. 选形与理由

**选形甲（改读法），且判据换形（计数等值 → 浏览器宿主 pid 同一）。**

- **形乙先排除**：这发用例夹带 AC#3 端口尺、AC#2 延迟、IsCreated/HWND 检查、dispatchRaw 真答复流——摘整发把这些全带走，派单明令此形不可行。
- **形甲内部的两次换形**（都有探针读数撑腰）：
  1. **先试"等值沉降"**（两侧对称的有界等待＋两样本一致才判，本仓 `panel_host_windows_live_test.go:60` 定式）⇒ `post-fix-count3.txt` 显示它把间歇红变确定红——**等值判据本身错了**，不是读法不够耐心。
  2. **换成宿主 pid 同一**：`TestAC4...`（focus hop）那发里真实第二窗会出现**新 host pid**——等值计数抓不到"等数换浏览器"，pid 集合差抓得到；helper 自然退场则完全不碰 host pid。判据＝「本发起的浏览器 host 必须活着跨过 transition（`mine-died=[]`）＋re-show 不得新增 host（`added=[]`）＋HWND 同一」。这是**等强升形**（同向更强：新抓"等数换浏览器"一形），阈值未放宽、断言未删，红句全部带数据。
- **跨发复用的豁免**（`post-fix-count3-v3.txt` 的教训）：`-count>1` 共用同一测试进程 pid，WebView2 对同一 user-data-dir **复用上一发的浏览器**（host pid 30228 跨 3 发同一、cold 1689→328ms）⇒ pre-hide 采一枚 **baselineHosts**（bring-up 之前），判决只对「本发新增的 host」（`mineHosts`）或「transition 期间新增的 host」（`added`）响，遗留／复用浏览器不背这发的锅。
- **等待实现**：`settleTreeReading`＝循环计数封顶（`treeSettleMax = treeSettleWait/treeSettleStep`，3s/100ms），无 `time.Since`、无墙钟差判据，与 winlive 档既有定式同形。
- ⛔ 产码一字未动（`panel_host_windows.go`／`panel_resident_windows.go` 保持 `4658dbb6` 入库形）。

### 改动清单（唯一写面 `cmd/wisp/panel_host_windows_test.go`）

| # | 位置 | 动作 |
|---|---|---|
| 1 | `readWebviewTree` 后 | 新增 `settleTreeReading`（有界沉降读法）＋常量组 `treeSettleWait/Step/Max` |
| 2 | `pidSetOfTree` 前 | 新增 `browserHostPids`（ppid==测试进程的 msedgewebview2＝浏览器宿主集合）＋`pidSetToSortedSlice/pidSetDiff/pidSetSubtract/pidSetIntersect` |
| 3 | 用例开头 `baselineTree` 旁 | 新增 `baselineHosts`（bring-up 前的 host 集合基线，跨发复用豁免的锚） |
| 4 | `:647-731` 一带 | hide→re-show 那段：两侧换 `settleTreeReading`；计数等值断言换成 host-pid 同一两向断言（`added`／`mine-died`）；denominators Logf 换成带 host 集合的一枚（并删掉我中途多写的那枚重复 Logf） |
| 5 | 全文 | 判据注释落满四条读数依据（探针文件具名） |

红句数据形（新断言的红都带集合）：
- `re-show spawned a new browser host inside THIS process tree: added host pid(s) %v (pre %v -> post %v ...)`
- `hide -> re-show lost the browser this window started: host pid(s) %v gone after re-show (pre %v -> post %v ...)`

---

## 3. 修后 `-count=3` 三发全绿的逐字读数

命令：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=3 -v -run 'TestPanelHostRealWindowHopAndLifecycle' ./cmd/wisp/`
原始输出：`.scratch/wisp/probes/255/r5/post-fix-count3-v4.txt`。终态 **3 发全绿**（`ok github.com/CarlosShao/wisp/cmd/wisp 2.336s`，exit 0）。

| 发 | 终态 | 关键读数（逐字摘） |
|---|---|---|
| 1 | PASS 1.15s | `browser hosts pre=[27592] mine=[27592] post=[27592] added=[] mine-died=[]` · `same HWND 0x9870324` |
| 2 | PASS 0.49s | `browser hosts pre=[27592] mine=[] post=[27592] added=[] mine-died=[]` · `same HWND 0x3940f1e`（mine=[]＝复用上一发浏览器，host 仍同一） |
| 3 | PASS 0.55s | `browser hosts pre=[27592] mine=[] post=[27592] added=[] mine-died=[]` · `same HWND 0x4d307e6` |

（三发里 settle 读数 `7→6` 与 `7→7` 都出现过，判据对两者都判绿——因为 host pid 同一才是身份，计数只是报告。）

四用例组合（派单点名的整包门禁替身）：`post-fix-gate4.txt`，`TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2|TestAC3ListeningSocketRulerSeesItsOwnListener|TestAC4FocusReturnToPriorWindowGap33r5` × `-count=3` ＝ **12/12 全 PASS**。

---

## 4. 门禁读数

| # | 门禁 | 命令 | 读数 |
|---|---|---|---|
| 1 | `go vet` | `go vet ./cmd/wisp/` | **过**（exit 0） |
| 2 | d22scan | `sh scripts/d22scan.sh` | **clean**：`no D22 ban violations; ... bans #1-5 internal/=228, bans #1-5 cmd/=37 ...`（`d22scan.txt`） |
| 3 | gofumpt | `"$(go env GOPATH)/bin/gofumpt" -l cmd/wisp/panel_host_windows_test.go` | **空输出**（无格式偏差） |
| 4 | 四用例 ×3 | 见 §3 | **12/12 PASS** |
| 5 | 整包 | `PATH=... go test -count=1 ./cmd/wisp/` | 两发（见下） |

整包终态（**判不动格，如实登记**：编排者"修完应 0 红"的前提被会话中途落库的两枚别腿产物打破了）：

- **发 A（非 verbose，跑完）**`whole-package-final.txt`：包终态 `FAIL cmd/wisp 31.461s`，`--- FAIL` 计数 **22**，红名集合（逐名）＝ `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` / `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` / `TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused` / `TestReplyListenerAllowsAnL2CardFromTheNativeSide` / `TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel` / `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` / `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` / `TestL1VetoNeedsAChannelTheHostReallyWired` / `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` / `TestNativeHostSeamRefusesAPanelSourcedAllow` / `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence` / `TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere` / `TestTicket255ReceiptStillNamesTheLiveReadSection` / `TestTicket255ReceiptSentenceAssemblyIsFiltered` / `TestTicket223RunArmsTheReloadTick` / `TestTicket223HandEditedFsLooseningCostsAnL2Card` / `TestTicket223RefusedLooseningKeepsOldValues` / `TestTicket223TighteningRaisesNoCard` / `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` / `TestTicket223RestartTierSaysItWillNotApply` / `TestTicket223FailureSentencesAreDistinct` / `TestTicket223PermissionDeniedSitsInItsOwnSentence`。
  **红因逐字同一**（26 处命中）：`wisp run: 文本角色未配置（Unconfigured）：config: llm: model "m1" of provider "acme" is disabled (enabled=false); re-enable it or remove the entry` ＝ **`261-r1`（`e2bfd7c8`，16:19 落库）的 Enabled 真门**把 fixture 里的 `enabled=false` 变成 `wisp run` 硬拒（exit 2），而 223/255-receipt/246 家族的 fixture 还没跟上。⇒ **这 22 红全部在 `261-r2` 的爆炸半径里**（台账 `ef8266db` 已派 261-r2），无一枚在我的写面；**r4 那发唯一的红名 `TestPanelHostRealWindowHopAndLifecycle` 已从红名集合里消失**（它在包内跑过且未进 FAIL 名单）。SKIP＝0。
- **发 B（verbose，被炸断）**`whole-package-final-verbose.txt`：跑到 `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` 时整包 panic——`panic: Fail in goroutine after TestAC1AlwaysBranchDoesNotRevertAHandEditedKey has completed`（`approval_reply_201_test.go:672` `waitForSentence187`；fixture goroutine 在测试终结后仍 `t.Errorf`）。panic 前读数：**PASS 10 / FAIL 22 / SKIP 0**。这枚 panic 也是 261-r1 门＋226 fixture 泄漏 goroutine 的合谋，同归 `261-r2` 射程；我的文件不在这条调用链上。
- **我的四用例族**：在整包（发 A）中 0 红；单独 `-count=3`（§3/§4#4）12/12 绿。我修的那格＝**确定绿**，不是靠侥幸。
- ⛔ 本程没有跑全仓 `go test ./...`（派单 ⛔：`internal/tools`＋`internal/risk` 有 `252-v1` 在测量、`internal/llm` 有 `261-r1` 在飞）。

---

## 5. 判不动

1. **AC#4 勾框**：归编排者（派单 ⛔）。
2. **`[panel] height=0` 语义**：编排者已按最小诚实定案（0⇒260），本程未触碰、无新读数。
3. **"规格比仪器宽"的 emoji 缺口**（AGENTS §1.2）：非本程射程，未碰。
4. **`4658dbb6` 产码本身的验收**（取值闭包是否真随 `[panel] width` 变）：AC#4 的"改 width ⇒ 真窗宽随之变"那一发**只有本机可量**，且真窗几何尺在 `panel_geometry_255_winlive_test.go`（`//go:build windows && winlive`）——本程修的是**同一族另一枚尺的确定性**，没有跑那发 winlive 几何读数（真窗尺寸量测归实现者自证链，非本格任务）。若编排者要那发读数，须另派或由 r3 证据链补。
5. **跨发浏览器复用对 cold 延迟读数的影响**（cold 1689→328ms 的跨度）：AC#2 百分位尺在整包单发里只采 1 样本（`n=1`），这个跨发现象对 SLO 判据无影响（预算 1500ms 单样本都过），但**若未来把 `-count` 调大**，AC#2 的样本独立性会变——记在这里，不改尺。
6. **`winlive` 档"NO CI job"**（`:6` 注释 vs `ci.yml` 零命中）那条盘上不一致：编排者已裁"归 `winlive-census-1` 那一族，本程不据此派改造腿"——我遵守，没动。

---

## 6. 推翻清单（谁写的哪句、我的读数支不支持）

| # | 原话 | 谁写的 | 我的读数 | 判语 |
|---|---|---|---|---|
| 1 | "初判是**间歇形**（`Destroy` 后立刻 `readWebviewTree`，webview 子进程异步退出还没退干净）" | 派单（编排者引 r4 初判） | `probe-names2.txt`：红发里**浏览器主体与其 helper 跨 transition 逐枚存活**；退掉的是**每发换 pid 的短命 helper**（建窗后 ~1s 自然退场，`probe-nohide.txt` 不 hide 也退） | **部分推翻**。"异步退出"对、"间歇"对，但**归因错**：(a) 不是 Destroy 的子进程没退干净——此用例 hide→re-show **不 Destroy**；(b) 退的是 bring-up 的短命 helper，不是"刚被 Destroy 的窗"的遗留 |
| 2 | "同一用例 13:34 PASS（1.24s）… 两发之间**产码没变** ⇒ 间歇形" | 派单 | `git diff 7cb9d570..081019c0 -- panel_host_windows_test.go` 空（尺没变）；但 `panel_host_windows.go`/`panel_resident_windows.go` 同区间 +159 行（`4658dbb6`） | **修正**："产码没变"对那把尺成立，但字面上"产码"在那两发之间是变了的（AC#4 取值闭包入库）；好在那变化只动建窗几何来源，与 hide→re-show 子进程计数路无涉——所以结论碰巧仍立，前提要写准 |
| 3 | "3 发里若 0 红，仍要按第 2 条把读数造确定" | 派单 | 我 3 发 **2 红 1 绿** | 支持（未触发 0 红分支） |
| 4 | （隐含）"红句在 `:649` 一带" | 派单 | `:649` 是 `t.Logf`，真红在 `:585` 的 `t.Errorf` | **推翻**（定位行号；r4 gate 文件里 `:649` 行内容恰好同块出现，易误引） |
| 5 | （隐含）"等它回来"式的等待（形甲的字面读法：有界等待**直到树内 webview 归零或超时**） | 派单形甲原文 | `post-fix-count3.txt`：等值沉降后红更确定（2/3）；`probe-curve.txt`：6 不回升 | **推翻**（判据层）：那枚短命 helper 退了就不回来，"等归零"对这枚用例的判据是死路；须换 host-pid 同一（本程落法） |
| 6 | "前两程把产码与尺都落了，收尾只断在确定性" | 派单 | 修尺过程中三次暴露**新拓扑事实**（helper 自然退场／跨发浏览器复用），原"计数等值"判据本身错 aim | **补充**：收尾实际包含一次**判据换形**（等值→身份），不是纯粹的"读法加等待"；此换形在派单"甲乙两形"的射程内（甲＝改读法，未动产码、未删断言、未放宽阈值），具名登记 |

### 终态

- 我的第一 commit：本文件随产尺 commit 一起落（见交件消息）；commit 后本节不回填 hash（hash 即交件消息里那枚）。
- 工作树终态：`cmd/wisp/panel_host_windows_test.go` 已提交；probes/255/r5/ 十份读数随件入库。
- ⛔ 全程 0 push、0 amend、0 reset/rebase/stash；显式 pathspec 提交。
- ⚠ 共享树上的别人现场（未触碰、具名登记）：`internal/tools/paths.go` 脏件＝`252-v1` 工作现场；`design/**` 的 ` D`／` M`＝r3 证据件里登记过的既存形态；`internal/llm`+`cmd/wisp/run.go` 已由 `261-r1`/`261-v1` commit 入库。
