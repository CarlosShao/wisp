# 235 — 票 222 修完之后，那枚**守池钉子注释的理由已经 falsified**（`subagent_197_test.go:398-402` 还在写"一次 spawn 为孩子的整条生命占着一枚桥位"），而**真桥那枚新用例只有 1/3 能看见真桥、还得靠 30s 挂死护栏才红**

- **Status**：**待派，小活**（编排者 09-29 18:4x 立，起手锚点 `f23586fa`）。来源＝非实现者验收腿 `222-v1` 交件里"具名的两枚新缺陷"（`docs/evidence/s1/222-spawn-holds-bridge-permit-v1.md`，41,944 字节／354 行），台账 `A447`。
- ⚠ **两枚都不记在 `222-r1` 账上**（那枚腿具名说过，编排者复认）：缺陷① 是**票 197 批次留下的注释**（`git blame`＝`7ea14ce3e`），缺陷② 是**票 222 新用例自己的射程不足**、由验收腿在突变里量出来。本票只做**测试面**，⛔ **零行为变更**。
- **Type**：仪器／注释面。**这不是安全洞、也不是产品缺陷**，照三行读：① **现象在哪出现**＝仓里的注释文字与一枚用例的检测力，不在任何对外接口上；② **有没有本机被入侵的证据**＝**没有**（全程只读）；③ **最坏后果是什么形状**＝**下一位读注释的人会按已经塌掉的理由做判断**，以及**一处真回归可能悄悄不红**。

## 现量（起手逐条自己复跑，别信这里的行号）

| # | 事实 | 读数 | 尺（编排者 09-29 18:4x 现跑） |
|---|---|---|---|
| 1 | 注释还在写**已被证伪的理由** | `subagent_197_test.go:398-402` 逐字「one in-flight spawn holds **one bridge slot for its child's whole life**（`bridge.run` keeps the semaphore held across `entry.Tool.Execute`），so every admitted child above the ceiling is a row the roster prints as 「在跑」…」 | `sed -n '396,404p'` 现读 |
| 2 | ⚠ **这句话的前提正是票 222 修掉的东西** | 票 222 **AC#2 已于 09-29 判成立并翻勾**（父等待**不再**占许可；`222-v1` 的 M1 读数＝`:443`／`:477` 两枚 0.00s 看到"4 枚桥位被等待中的父占着"，修法落地后该形状消失） | 票 222 `:36` 那格＋表 AC#2 行 |
| 3 | 但**钉子本身不许拆**：它的**结论**仍对 | `Test197SubagentPoolNeverExceedsBridgeCeiling`（`:404` 起）判的是 `MaxConcurrentSubagents > MaxToolConcurrency` 这件事，与"为什么"无关；`222-v1` 的 M6 实测把池抬到 8 ⇒ **4 枚红**（三枚 197 钉＋一枚 222 用例）⇒ **牙还在** | 表 AC#5 行（M6） |
| 4 | 真桥那枚新用例的**载荷行**不是最显眼那行 | `subagent_222_test.go:278` ＝ `ParentTools: h.bridge`（**这行才是**）；`:288` 的 `Tools: h.bridge` 是**死行**——`subagent_197.go:284` 逐字 `opt.Tools = newSubagentToolChain(...) from base` 会覆写它 | `sed -n '276,289p'` ＋ `sed -n '282,286p'` 现读 |
| 5 | ⚠ **只有 1/3 用例能看见"绕过真桥"** | `222-v1` 的 M3（把孩子的工具面退回 `h.dir` 原形）⇒ **三枚用例里只有 1 枚红，且靠 30s 挂死护栏才红**；另两枚对同一缺陷**失明** | 表 AC#1 注① |

## 判据

- [ ] **AC#1 把 `:398-402` 那段理由改成实话，⛔ 结论与断言一个字都不许动**：新理由要按今天盘上成立的机制写——**"父等待已不占桥位（票 222 AC#2），但池帽仍必须 ≤ 桥帽，理由是名册会把排队的孩子印成『在跑』"**（M6 抬池 ⇒ 4 枚红＝这条还有效）。⛔ **禁法三条**：不许顺手删那枚钉子；不许放宽它的比较方向（`MaxConcurrentSubagents > MaxToolConcurrency`）；不许把注释里"7/8-red shape ticket 197 leg A reported"那段历史改掉——**那是出处，不是理由**。
- [ ] **AC#2 给缺陷② 补检测力：在 `await` 之前先读一次 `h.probe.runs`**（`222-v1` 具名的补法）。判据＝**同一枚 M3 形状（孩子的工具面绕开真桥、退回 `h.dir`）下，三枚用例必须都红**，⛔ 且红**不许来自 30s 挂死护栏**（要 0.00s 级读数，形状照它 M1 那两枚 0.00s）。⚠ 这等于给票 222 AC#1 那格**再加一层牙**——本格交完要在票 222 `:35` 下追加一行指向本票，**不改它已勾的原句**。
- [ ] **AC#3 两枚突变自证**：① 把 AC#2 新加的前置读数注释掉 ⇒ M3 形状必须重新变成"只有 1/3 看得见"；② 把 AC#1 改过的注释**再改回旧那句** ⇒ 必须有测试红（⚠ **若零枚红，就具名承认这枚注释没有任何仪器保护**，并把它转成一条常驻断言或写明"只能靠人读"——**不许假装改完了**）。
- [ ] **AC#4 门禁与逐名**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 -v` 四数只能从 `-v` 量；`$(go env GOPATH)/bin/gofumpt.exe -l` 零输出；`go vet`；`tools/d22scan`。⚠ 本仓已知坑：`internal/tools` 起手是 **161 顶层声明／161 PASS／0 FAIL**（`222-v1` 收尾读数），**改动后逐名对拉、只按用例名集合比、不按枚数**；红名册里历史在册的 `internal/panel` 4 枚＋`internal/ball` 1 枚属别人地界，不算本票新增。

## 禁区

零产码（`internal/tools/subagent_197.go`／`bridge.go`／`entry/**` 的判定分支**一条都不许动**，只动 `_test.go`）；不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；⛔ **不许为了让读数变绿而放宽任何断言**；`frontend/**`／`design/**` 零读零写零转述；⚠ **与票 222／221／234 同撞 `internal/tools`／`cmd/wisp` ⇒ 一律串行**，且**排在票 234 之后**（234 要在当前 HEAD 上复跑那 12 枚读数，本票若先动 `_test.go` 会改变名册基线）。

## 与其它票的关系（别在这票里顺手做）

- **票 222**：本票是它 **AC#1 的覆盖面补强**，不是返工——那格已按其字面判据翻勾（表 AC#1 行"成立（两注）"）。
- **票 221**：AC#4 那两句"可以单独停它"（`subagent_197.go:196`／`:389`）与 `task.cancel` **零枚生产注册**（编排者 09-29 18:3x 现跑 `grep -rn '"task.cancel"' --include='*.go' internal/ cmd/ | grep -v _test`＝**0**）——⛔ **不在本票**，那是票 221 写腿的活（批准已落 `A434`，只读普查 `221-c1` 已交=`A433`）。
- **票 211**：那段注释"得先动契约"式的推论属**票 211 已塌的前提**一族；本票只改注释文字，⛔ 不动并发帽的任何数值。

## Progress log

- [2026-09-29 18:4x +08] agent=orchestrator did=立票（零产码）：`222-v1` 具名两枚缺陷转本票，注释段与载荷行/覆写点四把尺我自己现跑对上（`:398-402`／`:278` vs `:288`／`197.go:284`）next=排在票 234 之后派，先 AC#1 改注释＋AC#2 补前置读数
- [2026-09-29 22:3x +08] agent=235-r1 did=**骨架先落**（本轮零判语，只把台件与这份 log 立起来）：起手锚点 `25556879`（＝HEAD，`ledger(A452 …)`）；票面四把尺已复跑两把（`sed -n '396,404p'` ⇒ 注释起于 **`:398`**、被证伪那句在 **`:399-401`**、函数名行 **`:403`**；`grep -n "runs" internal/tools/subagent_222_test.go` ⇒ **`:115`** 字段／**`:153`** 自增／**`:374-375`** 现有读法／**`:513`** 第二处读法），另两把（`grep -rn` 注释文字命中面＋`internal/tools` 基线名册）正在跑，读数一律落 `.scratch/wisp/probes/235/r1/readings.md`。⛔ 本腿不动票面任何 AC 框，勾由编排者翻 next=基线名册 → AC#1 改注释 → AC#2 三枚用例补前置读数 → AC#3 两枚突变 → AC#4 门禁
- [2026-09-29 23:0x +08] agent=235-r1 did=**四格交完（零行为变更，只动两枚 `_test.go`）**，起手锚点 `25556879`、本腿三枚 commit `795ed767`（骨架）／`1d2ad737`（AC#1）／`a71f0be9`（AC#2）。逐发读数全在 `.scratch/wisp/probes/235/r1/readings.md`（16,509 字节）。
  - **AC#1**：`subagent_197_test.go` 的 `:398-402`（5 行）换成 `:398-410`（13 行）——新理由按今日盘上机制写：父等待已不占桥位（`giveBackWhileWaiting`），池帽仍须 ≤ 桥帽的理由是**名册会把排队等许可的孩子印成「在跑」**；历史句 `7/8-red shape ticket 197 leg A reported` 逐字在位并写明它是出处不是理由。⛔ 结论与断言零字符变更（`git show 1d2ad737` 的 diff 只有那段注释）。**末句"抬池⇒这枚钉子仍红"本腿自己复跑过**（M6 overlay：`MaxConcurrentSubagents=8`，产码未动）⇒ **4 枚红**：`Test197SubagentPoolNeverExceedsBridgeCeiling`／`Test197SubagentPoolCapsAtBridgeCeiling`／`Test197FullPoolRefusesNextSpawnWithReadableReason` 各 0.00s ＋ `Test222SpawnConclusionArrivesThroughRealBridgeChildren` 3.00s（按 `222-v1` 注②＝测试构造的阻塞，本腿不当"池 8 有害"的证据）。
  - **AC#2**：三枚用例各加一次**前置读数** `h.awaitChildOnBridge222`（`subagent_222_test.go:378-407`，调用点 `:444`／`:536`／`:604`，反控那发先 `drainBridgeArrivals` `:602`）＝"桥上的到达必须早于任何一枚父任务返回"，因果判定、不读墙钟；配套 `probe222.arrived`（`:134`，Execute 里 `:172-176` 非阻塞投签，不承载判据）。**读数对拉**：M3（`ParentTools`＋`Tools` 同时改指 `&fake197Dir{}`）**改前 1/3 红**＝只有 leg1 且红在 `:374`"只等到 0/4 枚，护栏到点"**30.00s**；**改后 3/3 全红、全 0.00s** ⇒ 本格要的"红不许来自 30s 挂死护栏"成立。票 222 `:35` 下已追加一行指向本票（`git diff`＝1 insertion，已勾原句一字未动）。
  - **AC#3**：① 把三处前置读数与 drain 全注释掉＋同一枚 M3 ⇒ **退回 1/3、红在 30.00s 护栏**（`ac3-1-prereadoff-m3-final.txt`）⇒ 检测力确由本格新加的读数买的，**成立**。② 把 AC#1 的注释改回旧那句 ⇒ ⚠ **跑不出期望读数：`goexit=0`、218 名全 PASS、`--- FAIL` = 0、红名册为空**。⇒ **具名承认：这枚注释没有任何仪器保护**（`grep -rn "bridge slot for its child|whole life"` 只命中注释自身，行为尺看的是 `MaxConcurrentSubagents > MaxToolConcurrency`，不看文字）。按票面走第二条出路：写明**只能靠人读**；⛔ 本腿**没有**为变红新写"断言注释文字"的测试（那正是本票要修的毛病）。"要不要转成一条常驻断言"＝要人拍，本腿不自裁（见台件 §3/§7）。
  - **AC#4**：终态 `internal/tools` 全量 `-v`＝**`=== RUN` 218／全名（含缩进）218／`--- FAIL` 0／`^FAIL	` 0**、`ok 13.206s`；⚠ 只数顶层 `^--- PASS`＝169（缩进子测试 49 枚）——票面写的"起手 161"是 `222-v1` 的过期数，本腿以现跑 218 为分母，没改计数口径。**逐名对拉**（基线 `baseline-tools-v.txt` ↔ 终态 `head-tools-v.txt`，`comm` 双向）＝**差集为空**、零新增红。`gofumpt -l internal/tools/` 零输出；`go vet ./internal/tools/` clean；`scripts/d22scan.sh` clean（`internal/` 462 枚含注释与 `_test.go`）。外加稳定性：`-count=3` 那 6 枚池／桥钉 18 名全 PASS；`-race -count=2 -run Test222` **0 处 DATA RACE**。每发突变之后 `git status --porcelain -- internal cmd`＝**空**（六发全走 `-overlay`，工作树零写入）。
- [2026-09-29 23:0x +08] agent=235-r1 did=**没做完／留给编排者**（⛔ 不为空）next=① AC#3②"转成常驻断言"这一支要么动契约面（要 `A##`）要么写一枚断言注释文字的测试（本票禁形）⇒ **要人拍**，本腿只走了"写明只能靠人读"；② ⛔ **票面四枚 AC 框本腿一枚没翻**，勾由编排者核过之后翻；③ 台账两笔没落（`docs/reports/pending-and-issues.md` 属编排者的笔）：票 222 AC#1 覆盖面补强已交＋"这枚注释无人保护"的具名承认；④ 票面 AC#4 的基线数 161 已过期（现跑 169 顶层／218 全名），改不改那句＝编排者的笔；⑤ 整包 `./cmd/wisp ./internal/...` 终态本腿没跑（本票 AC#4 射程只到 `internal/tools`＋三件门禁），故 `internal/panel` 4 枚＋`internal/ball` 1 枚历史在册红**未经本腿复核**
