# 票 92（返工后）独立对抗验收裁决表 —— acceptor-ticket92b

**验收开始**：`date` 实测 `Mon Sep 21 21:15:03 CST 2026`（会话第一条命令）。
**锚定 sha**：`91b5fc4`（代码）+ `a8f9459`（票面，= 当前 HEAD）。快照 = `git archive 91b5fc4 | tar -x -C /tmp/ac92b-sfx`（A38④，仓内无 worktree/checkout）。
**只读**：本轮**未改一行生产码**、未 commit、未 push、未改票面、未动 `docs/reports/pending-and-issues.md`；**未打开任何窗口、未截屏**（`R-92-5` 不许升格）。
**返工代理交件自述**：gofumpt 5⇒0；POSIX `278/175/0/1→4`；门二选结构面；`48/47/0/1` 撤回。

## 本轮要复核的四条被打回项 + 一条矛盾读数

| 编号 | 上一轮的账 | 本轮裁决 | 证据档 |
|---|---|---|---|
| R-92-3 | `gofumpt -l` 命中 5 文件；且交件写"本机无二进制（未跑）"不实 | **已清**（快照/工作树/`-l .` 三处皆空；不实登记已撤回） | 〔独立复现〕 |
| R-92-4 | POSIX 四数 `48/47/0/1` 不可复现 | **已清**（我逐位复算出 `278/175/0/4`，并把它的对账算式 `273+3+2=278` 用逐包分解证平）；⚠ 一条口径账：公布的命令块**逐字跑是 `rc=1/174/FAIL=1`** | 〔独立复现〕 |
| R-92-1 | 门二只覆盖 `.ts/.tsx/.css` 字面量 | **未清**：扩展名那一半真修好了（F3a/F4 我做出红），但**计算式路由 F5 全绿**、**通道放在 `frontend/src` 之外的 F6 全绿** | 〔独立复现〕 |
| R-92-2 | 真边界在原生侧（运行时 postMessage 不需要源码字样） | **未清**（本票未越界，`cmd/wisp` 零改动我核过）；另补一条它没点名的缺口：宿主有两种装法，钉只认一种 | 〔独立复现（"未越界"）〕 |
| 矛盾读数 | 容器整步 `portable-tests.sh` rc=1（`TestComposerRenderFixtureTellsTheTruth` FAIL）vs 92b 三包 FAIL=0 | **两边都不假**：那条 FAIL 是 `git archive` 的 **CRLF 伪形**（我把两形各量一遍），整步在 `91b5fc4`/`a8f9459` **仍 rc=1 但红因换人**，`e563a61` 已被票 111 转绿。92b 的形态**对 SKIP 那一半结构性看不见、对 fixture 那一半看得见** | 〔独立复现〕 |

**总判见下面的 `## 总判`（**附条件通过**），逐格裁决见 `## 逐格裁决表`，新账见 `## 新账 R-92b-x`。**

## 格 R-92-3 · `gofumpt` / `gofmt` 与"上一轮的谎"有没有再犯

**命令与读数（两处、逐条亲跑）**：

| 位置 | `gofmt -l internal cmd` | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal cmd` | `gofumpt -l .` |
|---|---|---|---|
| `91b5fc4` 纯净快照 `/tmp/ac92b-sfx` | **空**（rc=0） | **空**（rc=0） | **空** |
| 当前工作树（HEAD `a8f9459` + 邻居 ` M`：`internal/winsec/winsec_windows.go`、`scripts/portable-tests.sh`、`internal/models/*` 等） | **空** | **空** | **空** |

`gofumpt --version` 我本机实测 **v0.7.0 (go1.27.1)**，路径 `D:\work\base\gopath\bin\gofumpt.exe` ⇒ 与 92b 自述同版本同位置。
⚠ 一条 Git Bash 坑（不是它的错，但影响复算）：`export PATH="$PATH:$(go env GOPATH)/bin"` 之后 `gofumpt` 会解析成
`\work\base\gopath/bin/gofumpt`（MSYS 把盘符 `D:` 吃掉）⇒ `command not found`（rc=127）。
**必须**写 `"$(go env GOPATH)/bin/gofumpt.exe"`（上一轮验收方用的就是这个形式）。若有人照 92b 票面那句
`export PATH=... && gofumpt --version` 逐字跑，会得到"二进制不存在"的**假读数** —— 本轮它的命令块里给了 `which gofumpt` 与 `.exe` 两条退路，未构成不实登记。

**结论：`R-92-3` 已清。** 上一轮那 5 枚红（`attachments_test.go`/`bridge_test.go`/`composer.go`/`composer_test.go`/`workspace.go`）
在 `91b5fc4` 与工作树两头都复算为空；那句"本机没有 gofumpt 二进制（未跑）"已在票面 `91b5fc4` 正文与 `a8f9459` 两处撤回，
本轮交件**没有**新增任何"未跑/没有二进制"式的免检登记。〔独立复现〕

## 格 R-92-4 · POSIX 四数我自己量一遍

**同一棵树（`git archive 91b5fc4 | tar -x -C /tmp/ac92b-sfx`）、同一条命令、两种 fixture 字节形态**，
容器 `golang:1.27`（`uname -s`=Linux、`go version go1.27.1 linux/amd64`、`CGO_ENABLED=1`、`GOMODCACHE` 挂载、
挂载自证 `ls -l /wisp/go.mod` 有字节 + `ls /wisp/internal/tools | wc -l` = **35** + `MOUNT_OK` + `PANEL_OK`）：

| fixture 字节 | rc | `=== RUN` | `^--- PASS` | FAIL | SKIP | 含子测试的 `--- PASS` |
|---|---|---|---|---|---|---|
| **纯归档快照（CRLF，容器内 `FIXTURE_CR=8`）** | **1** | **278** | **174** | **1** | **4** | — |
| **工作树字节 cp 回（LF，`FIXTURE_CR=0`）** | **0** | **278** | **175** | **0** | **4** | **274** |

第二行与 92b 自述 **逐位相等**（`278/175/0/4`，含它另给的"含子测试 =274"）。4 枚 SKIP **名单也逐条相等**：
`TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`、
`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`、`TestWorkspaceSwitchRefusesAJunctionToOutside`。
唯一 FAIL 行的原文在 `--- FAIL: TestComposerRenderFixtureTellsTheTruth`，包级 `FAIL internal/panel / ok internal/tools / ok internal/risk`。

**对账（把它的"票 105 的 3 枚 + 本轮 2 枚"算式化）**：
- 上一轮验收在 `internal/panel` 量到 `count=2 → RUN=126 / ^--- PASS=76` ⇒ 单趟 **63/38**；我在 `91b5fc4` 的 Windows 单趟是 **65/40** ⇒ **panel +2**，
  且正是它点名的 `TestTheRendererHoldsExactlyOneDoorToTheHost` + `TestPlantedRendererDoorShapesGoRed` 两枚新用例 ⇒ **+2 成立**。
- 上一轮验收在票 105 落地后的 `internal/tools` 量到 `count=2 → 230/158` ⇒ 单趟 **115/79**（其 `b878b30` 时刻是 `224/152` ⇒ **112/76**，票 105 净 **+3**）。
- ⇒ `273 + 3(tools，票 105) + 2(panel，本轮) = 278`，`170 + 3 + 2 = 175`。**算式平。**
  ⚠ 但这条算式的**分母**要钉清楚：上一轮的 273 是在 `21da8f3`（票 105 之后）量的，`internal/tools` 目录当时 `ls | wc -l` = **34**、
  现在是 **35** ⇒ 那 +3 属 `internal/tools`，与 92 无关，它没有把它算成自己的覆盖，这一点它是诚实的。

**结论：`R-92-4` 已清，但带一条口径账**。它的四数**可重跑、逐位复现**，`48/47/0/1` 已正式撤回且未再引用；
⚠ 唯一残留：**它公布的命令块本身跑不出它报的数** —— 照块里逐字跑是 `rc=1 / 174 / FAIL=1`，要凑到 `rc=0/175/0` 必须**额外**做
它只在正文另一段用文字提过的 `cp 工作树 fixture 回快照`。这不是造假（它如实披露了成因、成因我复算为真：
`git archive` 在 `* text=auto` 下把该文件 8 个 LF 全变 CRLF，`-c core.autocrlf=false` 挡不住，我两条都验过），
但"可重跑"的最低标准是**命令块逐字可跑**，缺一步就该写进步骤里。〔独立复现〕

## 格 · 覆盖面对账（`ban #6` / `ban #8`，在**当前工作树**复算，未照抄）

`sh scripts/d22scan.sh`（工作树 = HEAD `a8f9459` + 邻居在飞的 ` M`），**rc=0**，末行 `clean - no D22 ban violations`：

| scope | 92b 报（`91b5fc4` 真树） | **我实测（当前工作树）** | 差 |
|---|---|---|---|
| `bans #1-5 internal/` | 202 | **202** | 0 |
| `bans #1-5 cmd/` | 20 | **20** | 0 |
| **`ban #6 frontend/`** | 43 | **43** | 0（≥ 票面基线 40 ⇒ 覆盖证明成立） |
| `ban #7 internal/tools/` | 18 | **18** | 0 |
| `ban #8 design/` | 16 | **16** | 0 |
| **`ban #8 frontend/`** | 43 | **43** | 0 |
| **`ban #8 internal/`** | 371 | **373** | **+2**（邻居在飞的 `internal/models`/`winsec` 文件，只增不减成立，与我上一条票 113/115 的观察同向） |
| `ban #8 cmd/` | 26 | **26** | 0 |

⚠ **快照 vs 工作树那条读数差我量出来了，别当覆盖下降**：同一枚扫描器跑 `91b5fc4` 的**纯净归档快照**给的是
**`ban #6 frontend/=40`、`ban #8 internal/=368`** —— 差额 3 恰好是 `frontend/dist/` 里那 3 个**未进 git 的构建产物**
（`index.html`、`assets/index-CEH-Pz8P.css`、`assets/index-UL9kYvYl.js`；`dist/` 里只有 `.gitkeep` 被 tracked）。
⇒ 43 与 40 之差**不是**本票覆盖面变小，是归档快照天生不含产物。92b 报的 43 是工作树读数，成立。〔独立复现〕

## 格 · 第 5 条：那条矛盾读数的裁决（`portable-tests.sh` rc=1 vs 92b 三包 FAIL=0）

我按派单要求在容器里跑了整步 `bash scripts/portable-tests.sh`（`golang:1.27`、`CGO_ENABLED=1`、挂载自证 `ls -l /wisp/go.mod` + `MOUNT_OK`），
外加一次**逐包分解**（这一步才是把 `R-92-4` 真正结掉的东西）：

| 树 | 命令 | rc | 四数 | 红因（点名） |
|---|---|---|---|---|
| `91b5fc4` **纯归档**（fixture CRLF） | `go test -count=1 -v ./internal/panel/ ./internal/tools/ ./internal/risk/` | **1** | 278/174/**1**/4 | `TestComposerRenderFixtureTellsTheTruth` |
| `91b5fc4` **fixture 还原为 LF** | 同上 | **0** | **278/175/0/4** | 无 |
| `b878b30`（上一轮验收那棵树，LF） | 同上，逐包 | 0 | panel **63/38** + tools **76/59/SKIP3** + risk **134/73/SKIP1** = **273/170/0/4** | 无（**与上一轮验收报的数逐位相等**） |
| `91b5fc4` 逐包 | 同上 | 0 | panel **65/40/0** + tools **79/62/SKIP3** + risk **134/73/SKIP1** = **278/175/0/4** | 无 |
| `91b5fc4` **整步**（LF） | `bash scripts/portable-tests.sh` | **1** | RUN=931 PASS=575 **FAIL=1** SKIP=1 | `internal/models` 的 `TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified`（**票 109 的修前红，不是票 92 的**） |
| `a8f9459`（会话开始时的 HEAD，LF） | `bash scripts/portable-tests.sh` | **1** | RUN=1002 PASS=641 **FAIL=0** SKIP=1 | **1 条未记账 SKIP：`TestWorkspaceSwitchRefusesAJunctionToOutside`**（`internal/tools/paths_workspace_test.go:198`，24 包全部 `ok (ran)`） |
| `e563a61`（交回时的 HEAD，邻居又推进了 11 枚） | `bash scripts/portable-tests.sh` | **0** | RUN=1036 PASS=661 **FAIL=0 SKIP=0**，`internal/panel` 打出 `ok (own line)`、`--- PASS: TestComposerRenderFixtureTellsTheTruth` | **无** —— 那枚未记账 SKIP 已由 **票 111 的 `7699ec3`** 记进 ledger（`… |linux|fixture|ticket111 AC#10 …`），整步转绿 |

⚠ 三条读数的口径都是**容器内 `PORTABLE_RC=$?` 取自被测命令本身**（外层 `docker` 的 rc 我另记，两者不同形）。

## 总判

**附条件通过（conditional pass）** —— **不是退回返工代理，也不判总 FAIL**，因为按派单口径"探针达成 AC 声称要防的结局才算 FAIL"，
我的三形**都没达成结局**（下面逐条给了实测理由）；但 **`R-92-1` 那一格我判"不通过"**（造形全绿的判据在此生效），
条件是 **`R-92b-1 / R-92b-2 / R-92b-3` 三格不许静默结案，`R-92-5` 不许升格**。

一句话理由：**`R-92-3` 与 `R-92-4` 我逐位复算成立（gofumpt 三处空、POSIX `278/175/0/4` 与它的对账算式 `273+3+2` 我用**逐包分解**证平），
但 `R-92-1` 按派单口径**未清**——我造的第三形里**计算式路由 `sendRequest("panel"+"."+noun+"."+verb,…)`（运行时 = `panel.approval.decide`）与
"第二条通道放在 `frontend/src` 之外"两形全绿**，四枚结构判据各留一字面前提，`AC#1` 那句"词汇是闭集、加路由必须两边同改"是**没被钉住的承诺**；
同时 `R-92-5` 那一格我量出了它的可见性证据有**一行深度的说谎空间**（组件退化 / 产物手写两形全绿，`render:composer` 不在 CI 任何一步）。**

**为什么不走"退回补齐"**：按派单总判规则，只有探针**达成 AC 声称要防的结局**才 FAIL/退回。三形我都没达成——
原生侧确实无处可发（我在 `91b5fc4` 重核：`ParseComposerRequest` 生产调用者 **0**、`perm.Store.Set` 生产调用者 **0**、
`cmd/wisp` 与 `frontend/**` 在 `91b5fc4..a8f9459` **零改动**、且 `internal/panel/bridge.go:99 knownComposerMethod` 是 Go 侧**闭集 switch**，
未列出的方法**运行时响亮拒**）。⇒ **附条件通过**成立，条件就是下面 `R-92b-1/2/3` 三格**不许静默结案**。

**我造了哪几形、各什么结局（派单要的对照本）**：F3a 红、F4 红、**F5 全绿**、**F6 全绿**、F6b 红（证明 F6 的漏不在扫描器而在结构钉的根）、
F1 红、**F2 全绿**、**G 全绿**、MUT-H 红（11 枚）⇒ **8 发里 4 绿 4 红**，四条绿各对应一条新账。

## 逐格裁决表（含三档标签；`92-adversarial-acceptance.md` 的 AC#1-#7 本轮只重判被返工碰到的格）

| 格 | 本轮裁决 | 证据档 | 原因（没亲手跑的一律第三档并写原因） |
|---|---|---|---|
| **R-92-3** gofumpt + "上一轮的谎" | **已清** | 〔独立复现〕 | 快照 / 工作树 / `gofumpt -l .` 三处我都跑，全空；v0.7.0 版本与路径我本机同验；票面那句不实登记已撤回 |
| **R-92-4** POSIX 四数可重跑 | **已清（带一条口径账）** | 〔独立复现〕 | 我逐位复算出 `278/175/0/4 rc=0` 与 SKIP 名单；对账算式我用逐包分解证平（`b878b30` = 63/76/134 = 273，`91b5fc4` = 65/79/134 = 278）。⚠ 它公布的命令块**逐字跑是 `rc=1/174/FAIL=1`**，少写了命令块里的 `cp` LF fixture 那一步（正文披露了，块里没有） |
| **R-92-1** 门二覆盖面 | **未清** | 〔独立复现〕 | F3a/F4 我做出红（它修好的那一半是真的），**F5 全绿**（它 `next=` 里自己点名的第三形的**加强版**），F6 另证扫描根比渲染器加载面窄 |
| **R-92-2** 真边界在原生侧 | **未清（本票不越界，诚实）** | 〔独立复现（"未越界"这一半）〕 | `cmd/wisp` 零改动我核过；我另补一条它没点名的同源缺口：宿主有两种装法，钉只认 `postMessage` 一种 |
| **R-92-5** 真机差分截屏 | **不升格、未清** | 〔仅自述，不背书〕 | **本轮我未开任何窗口**（owner 签收窗口未定）；我只裁了"fixture 有没有说谎空间"= **有**（F2/G 两形全绿） |
| AC#1 档位只读 | **不通过**（按派单口径：门的覆盖面残留） | 〔独立复现〕 | 反半边我钉住了（2 处合法调用点照旧绿）；正半边 F5/F6 绿 |
| AC#2 附件响亮失败 | **通过** | 〔独立复现〕 | MUT-H：统一出口改 `nil` ⇒ 11 枚红，点名到 7 个子测试 |
| AC#3 / AC#4 | **沿用上一轮"通过"** | 〔日志＋归档，我抽验〕 | 本轮返工范围不含；我抽验的是"承载这两格的文件 `91b5fc4..a8f9459` 零改动"（`git show --stat` + `git diff --stat` 亲跑）⇒ 上一轮读数对本轮成立 |
| AC#5 变异三向 | **(i) 本轮我重做为"红+绿各半"**（F3a 红 / F5 绿）；(ii) MUT-H 红；(iii) 未重做 | (i)(ii) 〔独立复现〕／(iii) 〔仅自述，不背书〕 | `internal/risk/**` 在禁改清单，我不去动它的生产文件；签收时补跑 `go test -count=1 ./internal/risk/` 即可 |
| AC#6 台账与门禁 | **通过（数字全复算）** | 〔独立复现〕 | `gofmt`/`gofumpt` 空、`go vet ./internal/panel/ ./internal/tools/ ./internal/memory/` rc=0、`go test -count=2 -v ./internal/panel/ ./internal/tools/` = **360/238/0/0 rc=0**（与它报的逐位相等）、`sh scripts/d22scan.sh` 工作树 rc=0 且 `ban #6 frontend/=43`、`ban #8 internal/=373`（它报 371，+2 是邻居在飞文件）；⚠ **整步 `portable-tests.sh` 它未跑**（`R-92b-4`） |
| AC#7 不做 git 切换 | **通过** | 〔独立复现〕 | 包内用例 PASS + 我把 `fixtures/`、`dist/` 也 grep 了一遍 0 命中 |

## 亲手命令清单（按可重跑形状给出）

见上面各格的表格内命令原文。汇总：
1. `date` ⇒ `Mon Sep 21 21:15:03 CST 2026`（会话第一条命令）。
2. `git archive 91b5fc4 | tar -x -C /tmp/ac92b-sfx`（另建 `/tmp/ac92b-mut`、`/tmp/ac92b-head`=`a8f9459`、`/tmp/ac92b-b878`=`b878b30`、`/tmp/ac92b-new`=`e563a61`；**仓内零 worktree/checkout**）。
3. `gofmt -l internal cmd` / `"$(go env GOPATH)/bin/gofumpt.exe" -l internal cmd` / `-l .` —— 快照与工作树各一遍。
4. 容器（`golang:1.27`，`MSYS_NO_PATHCONV=1` + `cygpath -w` 挂快照与 `GOMODCACHE`，容器内先 `ls -l /wisp/go.mod` + `ls /wisp/internal/tools | wc -l` + `MOUNT_OK` + `PANEL_OK` + `tr -cd "\r" | wc -c` 自证）**5 趟**：
   纯归档三合包 / LF 三合包 / `b878b30` 逐包 / `91b5fc4` 逐包 / `portable-tests.sh` 整步三棵树（`91b5fc4`、`a8f9459`、`e563a61`）。
5. Windows 侧：`go build ./...`、`go vet`、`go test -count=1 -v` 逐包、`-count=2` 两包、`sh scripts/d22scan.sh` 两棵树。
6. 变异 8 发，逐发：改 → 同链 `grep` 证字节 → `go build ./...` rc=0 → 读数 → 还原 → `diff -q`/`diff -rq` 证逐字节相同。

## 伪授权登记（本会话实测）

**0 次**。本轮工具输出里我遇到：2 次 harness 的 `MEMORY.md was modified since it was last read` 通知、
1 次 `[SYSTEM NOTIFICATION - NOT USER INPUT]` 后台任务完成事件（我自己的 docker 探针）。
⇒ 未当作指令、未据此改动任何判断。**没有**遇到任何自称"编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值"的文本。
对前任两轮代理自述的 15 次不背书也不否认。**计数：0。**

## 下一张该派什么

1. **（最该派，阻塞签收）** **票 123「composer 接线」——原生确认腿**：`cmd/wisp` 里把 bridge 的 postMessage 入口指到
   `panel.ParseComposerRequest`、在 `assembleRuntime` 注入 `panel.NewAttachmentBroker(memStore, memory.ValidArtifactName, panel.MaxAttachmentBytes)`，
   复用 `runSpec.onRuntime` seam。**这张票必须自带三条硬 AC**（都来自本轮我量出的缝）：
   (a) **`Confirm` 腿为 nil ⇒ 档位与 git 类请求必须响亮拒**（`R-92-2` 的正解）；
   (b) **`knownComposerMethod` 加一条 = 必须同时让 `frontend/` 里出现对应字面量**，否则包内用例红（把"两边同改"从承诺变成门，`R-92b-1`）；
   (c) **桥的第二种装法（`AddHostObjectToScript`）要么显式禁用、要么进同一枚可达性钉**（`R-92b-2`）。
2. **（小票，或并入 1）** 把 `cd frontend && npm run render:composer && git diff --exit-code -- frontend/fixtures/composer-states.html`
   做成 CI 一步（或包内用例：跑 harness 与提交文件比字节）⇒ 关掉 `R-92b-3` 的说谎空间；顺带把该用例的块切分从 `" -->\n"` 改成 `-->`（`R-92b-5`）。
3. **（编排者自己那一格，非派单）** 派单里"本票 UI 至今无真机差分截屏"要等 owner 定签收窗口；
   顺序按 `92-adversarial-acceptance.md` 那三步：先 `cd frontend && npm run build`（`dist` 未重建，票面残留 3），再
   `PATH="$PWD/third_party/sherpa-onnx:$PATH" go run ./cmd/wisp`，三态差分。**本轮我未开窗。**

**裁决（三句话）**：

1. **`audit-runner-readings` 那条 `TestComposerRenderFixtureTellsTheTruth` FAIL 是真的，但它是仪器伪形，不是生产红。**
   该证据自述树 = `git archive HEAD | tar -x`，我在**同一形状**上把它做出来了：纯归档快照里该文件 **8 个 CR**（`git archive` 在 `.gitattributes: * text=auto`
   下按宿主换行导出，`-c core.autocrlf=false` **挡不住**，两条我都量了），而该用例的块切分写死 `" -->\n"` ⇒ CRLF 下三个块全部解析不出来 ⇒
   读数正是那份证据里的原文 `the render fixture has no block for … ×3` + `verified across 0 painted states`。
   **工作树与该文件的 git blob 都是 CR=0**（`git ls-files --eol` = `i/lf w/lf`），`ubuntu` 上 `actions/checkout` 出来也是 LF
   ⇒ **真 CI 那一步不会因此红**；红的是"Windows 上用 `git archive` 造快照"这一族所有人（含我、含 92b、含上一轮验收）。
   **`R-92-5`/`R-92-4` 都不该为它记账；该记的是一条仪器账**：`TestComposerRenderFixtureTellsTheTruth` 对**换行字节敏感**，
   一条"证据文件"的用例把自己钉在了 checkout 方式上 ⇒ 新账 `R-92b-5`（不阻塞本票，属 `internal/panel` 地界，改法是把分隔符切成 `-->` 而不是 `" -->\n"`，或对该文件加 `*.html -text`）。
2. **这条 FAIL 现在还在不在**：**用例层面不在了**（LF 树里它 PASS，两棵树我都跑过）；**整步 rc=1 还在**，但红因已经换人：
   `91b5fc4` 时刻是 `internal/models` 的票 109 修前红（现由 `76fc5b0` 记为"109 因『链不在二进制里』退回"），
   `a8f9459` 时刻 models 已绿，**唯一剩的是那条未记账 SKIP**。
   ⚠ 那条 SKIP 的**主体用例是票 92 的**（AC#3(ii) 的真 junction 用例），**入账动作是票 111 做的**（`7699ec3 ci(111,AC#9/AC#10)` 把
   `TestWorkspaceSwitchRefusesAJunctionToOutside|./internal/tools/|linux|fixture|…` 写进 `scripts/portable-tests.sh` 的 ledger，
   我用 `git log -S` 定位到该 commit，并核到 `a8f9459` 那份副本里**还查不到这行**）⇒ **不是 92 的错，也不需要 92 修**；
   但账要记对：**票 92 在 Linux 上留了一条结构性 SKIP，而它的交件自述里把 4 枚 SKIP"逐条点名"了却没认出其中一枚会让整步红**
   —— 因为它压根没跑整步（见下条）。
3. **92b 的读数形态是否结构性看不见它：一半是，一半不是。**
   - **看不见的那一半（是）**：它跑的是裸 `go test`。**裸 `go test` 把 SKIP 记成 `ok`**（这正是 `scripts/portable-tests.sh` 头注写的立步理由），
     也没有 `runtests.sh` 的"任何顶层 `--- SKIP` 即致命"规则、没有 GUARD A/B/C、没有另外 21 个包。
     实测：它的命令在 `a8f9459` 那棵树上给出 `FAIL=0/SKIP=4 rc=0`，同一条树整步跑是 **`PORTABLE_RC=1`**。⇒ **这是门禁形状账，不是谁的错**。
   - **看得见的那一半（不是）**：`TestComposerRenderFixtureTellsTheTruth` 就在 `./internal/panel/` 里，**它那条命令覆盖得到**；
     它没报红，只因为它按自己披露的方法把 fixture 换成了 LF。**所以"92b 看不见 fixture 那条红"这个假设不成立**——它看得见，只是它选的树是绿的。
   - ⚠ 一条**未答的问**：票面 `20:5x` 的编排者注要求"你的 AC#6/AC#7 门禁读数**必须包含这条**（整步 portable）并说清是谁的账"。
     我在 92b 的交件段里 grep `portable` ⇒ **0 次命中**（它只提过 `TestComposerRenderFixtureTellsTheTruth` 的 CRLF 成因一次）。
     时序上那条约在它 `91b5fc4`（20:18）之后、`a8f9459`（21:12）之前插入，**它可能没读到** ⇒ 我记 `R-92b-4` 为**未答问**，不记为违规。

## 格 · 变异清单（锚点 / 落地证明 / 红的测试名）—— 共 8 发，全在 `/tmp`，仓内零 worktree

| # | 锚点与改法 | 落地证明（同一条 `&&` 链） | `go build ./...` | 结局 | 红的测试名 + 断言行 |
|---|---|---|---|---|---|
| F3a | `frontend/src/lib/panel.ts:237` `sendRequest("panel.mode.request", { to })` → `const routePrefix = "panel."; sendRequest(routePrefix + "mode.request", { to })` | `grep -n "routePrefix"` 打出 237/238 两行原文 | rc=0 | **红** | `TestTheRendererHoldsExactlyOneDoorToTheHost`@`:518`；`TestFrontendComposerRequestsMatchTheEnvelope` |
| F4 | 新增 `frontend/src/components/ac92b-fourth-door.tsx`，内 `window.chrome?.webview?.postMessage(...panel.mode.set...)` | `grep -n "postMessage"` 打出第 2 行原文 | rc=0 | **红 ×4** | `:506` outside / `:510` count=3 / `:522` unknown literal / 门二 `TestPlantedComposerModeWriteGoesRed`@`:140` |
| **F5** | `panel.ts` **纯增量**导出 `mut92bAskHostToDecide`：`sendRequest("panel" + "." + noun + "." + verb, { correlationId, outcome: "grant" })`，运行时 = **`panel.approval.decide`** | `grep -n "mut92bAskHostToDecide\|sendRequest(\"panel\""` 打出 240/243 行 + `grep -c` = 1 | rc=0 | **全绿（本票最硬的一条未清账）** | 无。`internal/panel` 全包 rc=0；`sh scripts/d22scan.sh` 同树 rc=0 `clean - no D22 ban violations` |
| **F6** | 第二条通道**放在 `frontend/src` 之外**：新增 `frontend/public/ac92b-door.js`（计算式路由）+ `frontend/index.html` 加 `<script src="/ac92b-door.js">` | `grep -n "ac92b-door" frontend/index.html` = 22 行；`grep -n postMessage` = 3 行 | rc=0 | **全绿** | 无。`composer_test.go:524` 仍打印 `21 files scanned, 2 host call sites` —— **一个数都没动** |
| F6b | 同一位置把路由换成**字面** `approval.decide`（测 ban #6 的**文件面**到底多宽） | `grep -n "approval.decide"` = 第 3 行 | — | **红** | d22scan rc=1：`frontend/public/ac92b-door.js:3: [panel-approval] ...` + 两枚自扫描用例红 ⇒ **ban #6 覆盖 `frontend/` 全树，结构钉只覆盖 `frontend/src`** |
| F1 | `frontend/fixtures/composer-states.html`：`no host attached` 块里 `档位未知` → `全自动` | `grep -c` 前后计数 | rc=0 | **红 ×2** | `TestComposerRenderFixtureTellsTheTruth`@`:706` + `:710` |
| **F2** | 同文件：把 `workspace chosen, video stored` 的 **1870 字节真 React 输出**换成 **108 字节手写** `<div>手工写下 …</div>`（四个必需子串都在） | Python 打印替换字节数 + 容器外 `repr()` 证子串在位 | rc=0 | **全绿** | 无；三包 `-count=1` 全 `ok` |
| **G** | **fixture 一字未动**，改真组件 `composer.tsx:123`：`{state.mode.l2ConfirmNames.includes(m) ? "（需原生 L2 强确认）" : ""}` → `{""}` | `grep -c "需原生 L2 强确认" frontend/src/components/composer.tsx` = **0**（同链打印） | rc=0 | **全绿** | 无，`^--- PASS=40`、`FAIL=0` |

**还原证明**：每发还原后 `diff -q` 对纯净快照逐字节相同；`diff -rq /tmp/ac92b-sfx/frontend /tmp/ac92b-mut/frontend` **空输出**；
末态 `go test -count=1 ./internal/panel/` **rc=0**。仓内 `git status --short` 与开工时同构（我只新增了这一份证据文件）。

**四种假绿各自过了一遍**：
- **跳过冒充通过** ⇒ 命中（见第 5 格：`TestWorkspaceSwitchRefusesAJunctionToOutside` 在 Linux 记 `ok`，整步靠 ledger 才看得见）；
- **断言恒真** ⇒ 未命中：`TestTheRendererHoldsExactlyOneDoorToTheHost` 有 5 条独立计数、`TestPlantedRendererDoorShapesGoRed` 走 `t.TempDir()` 副本正向控制，我 F3a/F4 独立做出红；
  但 **`composer_test.go:524` 那句 `t.Logf("… 2 host call sites (all in src/lib/panel.ts)")` 是无条件打印的字面**，F4 里它和上一行"outside call site"同时出现 ⇒ **日志自相矛盾**（一条小账，不判假绿，登记 `R-92b-6`）；
- **跑错对象/错文件** ⇒ 命中两发：`TestComposerRenderFixtureTellsTheTruth` 读的是**提交进 git 的产物 HTML**而不是组件渲染结果（F2/G）；
  结构钉的扫描根是 `frontend/src` 而不是渲染器实际加载的集合（F6/F6b）；
- **门禁压根没跑** ⇒ 命中一发（`R-92b-4` 整步 portable 未跑）；`gofumpt` 那一枚**本轮真跑了**（三处读数我全部复算为空）。

## 格 · AC#7 负判据（"不做 git 切换"）

`TestNoGitSwitchCapabilityInThePanelSurface` 在 `91b5fc4` 快照 **PASS**，日志原文
`AC#7 negative criterion held: no git switch entry point in frontend/ or internal/panel/ (production files)`。
我另把它**跳过**的两棵树（`fixtures`、`dist`）在工作树里用同一条 11-pattern 正则逐字 grep ⇒ **0 命中**（含 `frontend/dist/assets/index-*.js` 产物）。
⇒ **AC#7 通过**，本轮没有把 git 切换能力加回面板面（`git show --stat 91b5fc4` 只有 5 个 `internal/panel/` 文件，`cmd/wisp`、`frontend/**` 一字节未动 ⇒ "顺手加回来"这一族为零）。〔独立复现〕

## 格 · `R-92-5` 不许被升格：fixture 与真实渲染之间有没有"说谎空间"

**先说纪律**：本轮**未开任何窗口、未截屏**；我不给 `R-92-5` 任何升格，也不判它清。下面只裁"证据够不够格当可见性证据"。

**链路读码（一条决定性事实）**：`frontend/fixtures/composer-states.html` 是 `npm run render:composer`
（`frontend/scripts/render-composer.tsx`，内含 `renderToStaticMarkup(<Composer state=… />)` + 三态 `expect`/`forbid` 断言 + `failures>0 ⇒ return 1`）
**写出来的产物文件**；而 `TestComposerRenderFixtureTellsTheTruth` 读的是**提交进 git 的那份 HTML 文本**。
我把 `.github/workflows/ci.yml` 逐字 grep 过：CI 跑 `npm run build`（381 行）与 `npm run render:l2`（389 行），**没有任何一步跑 `render:composer`**。
⇒ 组件与那份 HTML 之间的唯一纽带**不在任何门禁里**。这不是推断，下面两发变异各自把它做成读数。

| # | 变异（`/tmp/ac92b-mut`，同链 `grep` 证字节 + `go build ./...` rc=0） | 结局 | 说明 |
|---|---|---|---|
| **F1** | 改 fixture：把 `no host attached` 块里的 `档位未知` 换成 `全自动` | **红**（派单点名的那一发，达成了） | `composer_test.go:706 "with no snapshot at all the row does not say so (\"档位未知\" missing)"` + `:710 "an unknown mode rendered as a named档"`；同跑 `:728` 仍打印 `verified across 3 painted states` |
| **F2** | 把 fixture 里 `workspace chosen, video stored` 那一整块**真实 React 输出 1870 字节**删掉，换成一行**手写** `<div>手工写下 D:\work\Wisp\notes 与 clip.mp4 video/mp4 需原生 L2 强确认</div>`（108 字节） | **全绿** | `--- PASS: TestComposerRenderFixtureTellsTheTruth`，并且我另跑**三包** `-count=1` ⇒ `ok panel / ok tools / ok risk`（rc=0）：整包里没有任何一枚用例能分辨"这是组件画出来的"还是"这是人手写的字" |
| **G** | fixture **一字未动**，改**真组件** `frontend/src/components/composer.tsx`：把 `{state.mode.l2ConfirmNames.includes(m) ? "（需原生 L2 强确认）" : ""}` 换成 `{""}`（`grep -c "需原生 L2 强确认"` 在该文件 = **0**，同链证明） | **全绿** | `go build ./...` rc=0，`go test -count=1 -v ./internal/panel/` rc=0、`^--- PASS=40`、`FAIL=0` ⇒ **"HTML 里写了 L2 提示、实际渲染不出 L2 提示"这一形本票的门根本看不见**，而这正是派单第 4 条问的那句话 |

**裁决**：**有说谎空间，而且是一行深度的**，两个方向各证成一发（F2：产物可手写；G：组件可退化而产物不重生成）。
F1 只证明它抓得住"有人改了 HTML 里的字"，抓不住"HTML 不再对应组件"。
⇒ 该用例能发现**措辞漂移**，**不能**发现**渲染漂移**；票面残留 1 写的"它证明值被画出来了"这句**说过头了**，
应改成"它证明**那份被提交的 HTML 里含这些字**，画没画要等 `render:composer` 进门禁或等差分截屏"。
这不阻塞结案（它本来就不是像素证据，`R-92-5` 仍欠真机差分），但**属门禁形状账**：一条名叫 `TellsTheTruth` 的用例给出一条它测不到的承诺 ⇒ 登记 `R-92b-3`，修法唯一且便宜：
把 `cd frontend && npm run render:composer && git diff --exit-code -- frontend/fixtures/composer-states.html` 做成 CI 一步（或包内一条用例：跑 harness 并把 stdout 与提交文件比字节）。
**本轮我没有权限动 CI，也明确不动。**〔独立复现：F1/F2/G 三发都是我亲手种、亲手还原，还原后 `diff -q` 证与快照逐字节相同、包回 rc=0〕

## 格 · AC#2 的"响亮失败"这一族（本轮唯一一处我自己重做的能力变异）

派单把"响亮失败被静默"列为总判 FAIL 的触发条件之一，所以我没有只复述上一轮的读数，重做了一发：
**MUT-H** = 把 `internal/panel/attachments.go:282` 的**统一出口** `refuse()` 的 `return ref, fmt.Errorf(ErrAttachmentRejected…)`
改成 `return ref, nil`（同链 `grep -n "return ref, nil"` 打出 276/291 两行、`go build ./...` rc=0）⇒
`go test -count=1 -v ./internal/panel/` **rc=1，FAIL=11**，点名到
`TestRefusesUnsupportedAndMasqueradingInputsLoudly` 的 **7 个子测试**（空文件 / 伪装成 .png 的 exe / 未知类型 / 谎报 mime /
路径注入名 / 2GB / nil sink）+ `TestMessageCarriesAttachmentRefsAndRefusals` + `TestAttachmentPayloadCarriesTheBytes`（含 `an empty payload is refused`）。
还原：`diff -q` 逐字节相同、包回 rc=0。
⇒ **AC#2 的"不许静默丢弃"是活的门**（不是恒真断言、不是跑错对象）。〔独立复现〕
AC#3（工作区）与 AC#4（关面板不丢）本轮**未重做**：返工范围明写只补 ①②+门二，两格上一轮由 `acceptor-ticket92` 以自有探针复现过
⇒ 我这一轮**沿用其结论但标〔日志＋归档，我抽验〕**：我抽验的是"这两格的文件本轮一字节未动"（`git show --stat 91b5fc4` 只有 5 个 `internal/panel/` 文件，
`internal/tools/paths_workspace.go`、`internal/tools/paths.go`、`internal/risk/**` 在 `91b5fc4..a8f9459` 零改动 ⇒ 上一轮读数对本轮成立）。

## 格 · 一条"未跑"登记的文字复核（派单第 2 条的后半）

92b 本轮**没有**再写"本机没有某二进制"。它写了另一条免检：
`frontend/**` 本轮一行未改 ⇒ 未重跑 `npm run typecheck`/`oxlint`，登记为"本轮不适用"。
- **前提一半真一半假**：`git show --stat 91b5fc4` 确实只带 `internal/panel/` 五个文件（**它自己的 commit 一行未改前端 = 真**）；
  但"本轮"到 `a8f9459` 为止，`frontend/` 相对上一次 typecheck 读数所在的 `8e10095` **有 1 行改动**
  （`git diff --stat 8e10095..HEAD -- frontend/` ⇒ `frontend/src/lib/panel.ts | 2 +-`，来自 `afd47cb fix(92,AC#2)`）。
- 我读了那一行的**原文**：只有 JSDoc 注释 `/** The sniffed type, never the declared one. */` → 末尾多了 `; empty unless stored.`
  ⇒ **注释文本，`tsc -b` / `oxlint` 的输出不可能因此变** ⇒ 结论仍然成立，但**理由要给对**：不是"一行未改"，是"改的那一行是注释"。
  这条登记为 `R-92b-7`（**不阻塞**、纯精度账；上一轮那枚硬谎是"二进制不存在"，性质不同，我不升格它）。
- ⚠ 另记一条**票面残留 3 的连带事实**：`frontend/dist/` 未重建（`git ls-files frontend/dist` 只有 `.gitkeep`），
  所以"归档快照 ban #6 = 40 / 工作树 = 43"这 3 枚之差**就是产物**，别当覆盖下降（我在覆盖格里量死了）。

## 新账 `R-92b-x`（本轮新增，逐条标阻塞性与地界）

| 编号 | 内容 | 阻塞结案？ | 地界 |
|---|---|---|---|
| **R-92b-1** | 结构钉的三枚新判据各自留了一缝：`.join(".")` 只堵一种拼法、`computedSendRequestRe` 把**以引号起首的拼接**当字面量、`routeLiteralRe` 只认 `"panel.*"` 这个字面形状。⇒ **F5 全绿**：`sendRequest("panel"+"."+noun+"."+verb, …)` 运行时就是 `panel.approval.decide`，门一（ban #6）+门二+四枚结构判据**六把锁一起静默** | **不阻塞**（原生侧今天无处可发：`ParseComposerRequest` 生产调用者 0、`perm.Store.Set` 生产调用者 0，且 `knownComposerMethod` 是 Go 侧闭集 switch，未列方法**运行时响亮拒** ⇒ 探针**没达成**结局；但**"加路由必须两边同改"这句承诺不成立**） | `internal/panel`（`composer_test.go`） |
| **R-92b-2** | 结构钉的扫描根是 `frontend/src`，比 ban #6 自己的文件面（整个 `frontend/`）**和**渲染器实际加载的集合（`index.html` + 同源脚本）都窄。⇒ **F6 全绿**：第二条通道放进 `frontend/public/ac92b-door.js` + `index.html` 引用，`21 files scanned, 2 host call sites` 一个数没动。**F6b 证明这不是设计而是漏**：同一位置写字面 `approval.decide` 时 ban #6 **能**看见并 rc=1 | **不阻塞**（同上：无听众；且 `index.html` 的 CSP 是 `script-src 'self'`，内联脚本本来就跑不了，能跑的是同源文件） | `internal/panel`（扫描根）＋一句前端事实 |
| **R-92b-3** | **fixture 的可信度只到"字符串在 HTML 里"**：`npm run render:composer`（唯一把断言打在**真组件渲染结果**上的东西）**不在 CI 任何一步**（`ci.yml` 只跑 `npm run build` 与 `npm run render:l2`）。⇒ **F2 全绿**（1870 字节真 React 输出换成 108 字节手写文字）、**G 全绿**（真组件不再画"需原生 L2 强确认"，HTML 照旧声称画了）。它**能**抓住 F1（改 HTML 的字），**不能**抓住渲染漂移 | **不阻塞结案**，但**阻塞"票面残留 1 的那句话"**：`它证明值被画出来了` 要改成 `它证明提交进 git 的 HTML 含这些字` | `internal/panel` + `.github/workflows/ci.yml`（编排者的门，我只登记） |
| **R-92b-4** | 票面 `20:5x` 编排者注要求"AC#6/AC#7 读数必须包含整步 `portable-tests.sh` 并说清是谁的账"⇒ 92b 交件 grep `portable` **0 次命中**，未答。（时序上可能没读到，故不判违规） | 不阻塞（其内容已由票 111 在 `7699ec3` 落地） | 本票票面 AC#6 |
| **R-92b-5** | `TestComposerRenderFixtureTellsTheTruth` 用 `" -->\n"` 切块 ⇒ **对换行字节敏感**，任何 `git archive` 快照（Windows 宿主，`* text=auto`，8 个 CR）都会让它**假红**；本轮我与上一轮验收、`audit-runner-readings` 三方各自踩到的是同一枚 | 不阻塞（LF 树/真 CI 都绿），但**是仪器账**：改法是切 `-->` 而不是 `" -->\n"`，或 `frontend/fixtures/*.html -text` | `internal/panel` |
| **R-92b-6** | `composer_test.go:524` 的 `t.Logf("renderer door held: … (all in src/lib/panel.ts)")` 是**无条件字面**，F4 里它与上一行"outside call site"同屏出现 ⇒ 日志会自相矛盾（断言没错，读数会骗人） | 不阻塞 | `internal/panel` |
| **R-92b-7** | "本轮 `frontend/**` 一行未改"这句免检前提**在树级别为假**（差 1 行，且那一行是注释 ⇒ 结论仍成立）；理由应写成"改的是注释" | 不阻塞 | 本票票面 |
| **R-92b-8** | 它为"不扩 AC#7 扩展名"给的理由**字面为真**（`frontend/scripts/vendor-shadcn.mjs:2` 的注释里确实有 `git checkout`，我读了原文），但推论**未穷尽**：本包早为解决同一问题写过 `codeOnly()`，"**扩面 + 剥注释**"这一组合它没测过就宣布不可维护 | 不阻塞（改进账） | `internal/panel`（`composer_test.go` AC#7 那枚谓词） |

**没有一条达成"面板改宽档位 / git 切换回到面板面 / 响亮失败被静默"**：
第一条被我逐层核到生产面（零调用者 + Go 侧闭集 `knownComposerMethod` + `cmd/wisp` 零改动）；
第二条我另把 AC#7 的 11 条 pattern 扩到 `fixtures/` 与 `dist/` grep ⇒ 0 命中；
第三条我用 MUT-H 主动把它做成 11 枚红 ⇒ 它是活的。

## 格 R-92-1 / R-92-2 · 门二选"结构面"是不是一个更弱的门

先钉反半边（它必须照旧通过）：`91b5fc4` 快照未改一字时 `go test -count=1 ./internal/panel/` **rc=0**（Windows，`RUN=65 / ^--- PASS=40 / FAIL=0 / SKIP=0`），
`TestTheRendererHoldsExactlyOneDoorToTheHost` 的日志原文 `renderer door held: 21 files scanned, 2 host call sites (all in src/lib/panel.ts), 5 panel.* route literals`
⇒ **`panel.ts` 里那 2 处合法调用点照旧绿，门不是"永远红"**。

造形清单（每一形：改 → 同一条 `&&` 链里 `grep` 证字节在文件 → `go build ./...` rc=0 → 才读数）：

| # | 形状（全部落在 `/tmp/ac92b-mut`，仓内未建 worktree） | 结局 | 红的用例 / 断言行 |
|---|---|---|---|
| **F3a** | `panel.ts` 内：`const routePrefix = "panel."; sendRequest(routePrefix + "mode.request", { to })`（它自己在 `next=` 里点的第三形，标识符起首） | **红** | `TestTheRendererHoldsExactlyOneDoorToTheHost` @`composer_test.go:518`（computed route）+ `TestFrontendComposerRequestsMatchTheEnvelope` |
| **F4** | `panel.ts` 之外新增 `src/components/ac92b-fourth-door.tsx`，内含 `window.chrome?.webview?.postMessage(JSON.stringify({ method: "panel.mode.set", to }))` | **红（四重）** | `:506` outside call site、`:510` count=3、`:522` unknown literal、门二 `:140` `TestPlantedComposerModeWriteGoesRed` |
| **F5** | `panel.ts` 内**纯增量**（不动任何既有字面量）：`sendRequest("panel" + "." + noun + "." + verb, { correlationId, outcome: "grant" })`，`noun="approval"`、`verb="decide"` ⇒ 运行时路由 = **`panel.approval.decide`**（PLAN.md:1588 那条红线本身） | **全绿** | 无。`internal/panel` 全包 rc=0；`sh scripts/d22scan.sh` 同一棵树上 rc=0 `clean - no D22 ban violations`；容器 ban #6 也看不见它 |

**⇒ `R-92-1` 未清。** 结构面确实把它上一轮那两形堵住了（F3a/F4 都红，扩展名那一半是真修好了），
但它的四枚结构判据里，**唯一能挡住"拼名"的两枚各自有一个字面前提**：
(ii) 只认 `.join(".")` 一种拼法、(iii) 的正则 `sendRequest\(\s*(?:[A-Za-z_$][\w$]*|`|\[)` **把"以引号起首的拼接"当字面量放过**、
(iv) 只扫 `"panel.*"` 这个**字面量形状**（`"panel"` 与 `".mode.request"` 分开引就不再是 `panel.` 字面量）。
F5 同时命中这三条缝，得到的是**比上一轮更靠近本票红线的那一形**：它不是"面板改宽档位"（`panel.mode.request` 本来就是合法请求），
而是**面板向宿主发出"批准这条操作"的路由** —— 也就是 `AC#1` 引的 `PLAN.md:1588` 与 ban #6 立规的原话。

**但按总判规则我仍不据此判本票 FAIL**，理由与上一轮同一条且本轮复核仍然成立（不是我推断，是 grep 亲测的形状）：
原生侧今天**没有任何听众** —— `ParseComposerRequest` 与 `perm.Store.Set` 的生产调用者各 0（上一轮实测），
`cmd/wisp/run.go` 本轮一行未动（票面残留 2 自述，`git show --stat 91b5fc4` 只有 5 个 `internal/panel/` 文件 ⇒ 成立），
所以 F5 **达不成**"面板改宽档位 / 响亮失败被静默"这个 AC 声称要防的结局。⇒ 记**门的覆盖面**账（`R-92b-1`），AC#1 该格**不通过**（按派单口径），总判走"附条件"那条。

**`R-92-2`（真边界在原生侧）**：本轮**未清也未推进**，且它自己在票面把这条腿明写"归接线票"，我核到它**没有偷装那条腿**（`cmd/wisp` 零改动）⇒ 是诚实的越界自限。
另记一条它没点名的同源缺口：`panel.ts:140` 的注释自己写着宿主有两种装法
（`WebView2's AddHostObjectToScript / postMessage pipe`），而结构钉**只认 `.postMessage(`** 一种；
`hostObjects.*` 那一形源码里不出现 `.postMessage(`、不出现 `panel.` 字面量 ⇒ 结构面天生扫不到它。这是 `R-92-2` 的另一半，归接线票。

**它自报的"扩面会让 `vendor-shadcn.mjs:2` 的一句注释红" —— 复核：字面为真。** 我读到该行原文：
`// scripts/vendor.mjs but against a registry instead of a git checkout:` ⇒ 含 `git checkout`，
而 `gitSwitchCapabilityRe` 走的是 `reMatchesInFile`（**不剥注释**），把它换成 `rendererSourceFile` 的扩展名集（含 `.mjs`）必红。
⚠ 但**理由成立 ≠ 推论成立**："因此 AC#7 不可扩面"这条推论被本包自己已有的一件工具反驳：门二就是因为同样问题写过 `codeOnly()`
（`composer_test.go:650`，剥注释不剥"穿着注释的代码"）。**扩面 + 复用 `codeOnly()`** 这条组合它没测过。⇒ 记 `R-92b-8`（不阻塞）。
〔独立复现：F3a/F4/F5/d22scan/vendor-shadcn 行；F3a/F4/F5 的落地证明是同一链里的 `grep -n` 原文 + `BUILD_OK`〕

**共树时间漂移（写这份表时的状态）**：我 21:15 开工时 HEAD=`a8f9459`，21:46 已推进到 `e563a61`（邻居 11 枚），
交回时工作树带着别人的 ` M cmd/wisp/{resident_windows.go,run.go}`（`run.go` +33/−1，票 117/121 在飞）与 ` M docs/reports/pending-and-issues.md`，
**这些都不是我的**，我也一行没碰（我只新增了本文件）。
⚠ 我总判的那条前提**在交回这一刻仍然成立、但会被接线票本身作废**：
`grep -rn "ParseComposerRequest(\|NewAttachmentBroker(" --include='*.go' cmd/`（**含脏工作树**）⇒ 仍然 **0 命中**。
⇒ **下一位验收方必须在接线票落地后重测这一条**，不能引用我这里的 0；本表的"无处可发"只对 `91b5fc4..e563a61` 负责。

**验收时间戳**：会话开始 `date` 实测 `Mon Sep 21 21:15:03 CST 2026`；本文件写完时 `date` 实测 `Mon Sep 21 21:53:48 CST 2026`。
**未 commit、未 push、未改票面、未动 `docs/reports/pending-and-issues.md`、未开任何窗口、未改一行生产码。**
