# 票 114 接线终裁腿 `114-v2` — `114-r2` 的接线补上了哪几格、哪几格 ⛔

**腿** `114-v2`（**非实现者／裁决腿**）· **被裁对象** `80376d58f`（六件）＋`fda84b9ca`／`6ea57bbee`（件 144 行）
**起手锚** `6ea57bbee`（＝被裁对象尖端）· **收手时窗口已推进到** `11790a4ed`：窗口内他腿 `03179ea65`(A1213 派本腿)／`11ede6ef5`(m4-survey)／`26132d82e`(A1214 立票 372)／`e66c8eb65`+`b47638e58`+`11790a4ed`(114-r3)
**写面**：本件＋`D:\tmp\wisp\20261024\114v2\**`（仓外）。⛔ 产码／⛔ 翻框／⛔ 票名·`Status:`·台账／零 `push`。
**被测树**：`git archive HEAD`（＝`6ea57bbee`）→ `D:\tmp\wisp\20261024\114v2\tree`；六件 sha256 与 HEAD blob 逐值同（§7 末）。

## §0 三句判语（先给结论）

1. ★**接线做对了**，⛔「把空门接活了」：门在原生侧且不可绕、两处 `nil` 必须一起接、卡体与适配器各只有一枚、装配根 boot 交腿且**覆盖无控制台形状**；`M4` 一格没扩（§6）。
2. ★★但它**欠一句具名**：这是本仓**第一次让面板（渲染进程可达）拿到「静默放宽一档」的能力**——`AE→AH` 现在 `err=nil`／`disk=ask_high_risk`／`confirm_calls=0`。这**符合 `M4`**、⛔ 是 r2 的缺陷；但**能力增量本身要记账**，而它 §0 只写了「过门进写入腿」。
3. ★★`AC#3` **一半成立、一半 ⛔**（§11）。`CLI 仍拒` **今天是对的**，但它承重于票 367 的方向判据——**A1214 已判 367 反主流并立 372 回退** ⇒ 372 一落，r2 的 CLI 用例必红（§4，本件唯一一条前瞻风险）。

## §1 票面前提新鲜度（对 HEAD blob 复量）＋`114-r1` 三条复核

| 项 | 尺（逐字） | 读数 |
|---|---|---|
| 票面＝HEAD？ | `git rev-parse HEAD:<票面>` vs `git hash-object <工作树票面>` | 起手**同值** `6bf5fbbb…`（票面当时没被别人改） |
| `ParseComposerRequest` 产线调用者 | `git grep -nE "ParseComposerRequest" -- "*.go" ":(exclude)*_test.go" ":(exclude).scratch/**" ":(exclude)docs/**"` | **9 行**＝定义 `bridge.go:134`＋**调用 `composer_dispatch.go:178`**＋7 行 `//` ⇒ **1** |
| `perm.Store.Set` 产线调用者 | `git grep -nE "\.Set\(ctx" -- "*.go" ":(exclude)*_test.go" ":(exclude).scratch/**" ":(exclude)docs/**"` | **恰 1**：`composer_handlers.go:141` |
| 通路在不在 | 自建真二进制（仓外）`wisp.exe panel-inbound -data <seed>` | 非变宽封套 `rc=0` 逐字「第 1 行已受理，**处理器已被调用**」；`AE→AH` `rc=1` 逐字拒绝句、盘未动 |

⇒ 票面三处原句（标题／`:21`／`:44`）**在起手锚上都过期**（「零调用者」「通路仍不存在」均不成立）。
★**但窗口内已被别人的腿处理**：`114-r3`（`e66c8eb65`）在票面**追加**一节「★前提更正」（**原文一句未改**、框未翻），其表里两把尺与上表**逐值同**；并把 `composer_handlers.go:41-45` 注释与 `:132` 括注改准 ⇒ **r2 §11.4 那两处具名上报，事后由另一条腿按它的判法处理了**（⛔ 它漏做）。
★**注**：`internal/panel/composer_handlers.go` 在我起手 10:17 时是**他腿未提交的脏改**（mtime 10:17:17，6/6 行；收手前查明＝`114-r3`，已由 `e66c8eb65` 入库、diff 逐行同），我的 `-overlay` 探针因此带过一句新括注（「只有全自动要确认」）；**HEAD 原句是「只有变宽要确认」**，r2 引的那句对它当时的树**是准的**。

**`114-r1` 三条复核（等价形状，⛔ 逐字复跑它那三发）**

| r1 三条 | 我的读数 | 判 |
|---|---|---|
| ① 产线调用者＝1 | 上表 ⇒ 1 | **成立** |
| ② 深度 1／2 都拒、盘未动、正控过 | 深度1（直调 `HandleModeRequest`）／深度2（走派发表）⇒ 都 `ErrNoL2Confirm`、`disk=ask_every_step`；正控接腿⇒过 | **成立** |
| ③ `AA` 那发 `confirm_calls=1 asked=[auto_approve]` | 常驻门＋**真 gate**：`calls=1 asked=[auto_approve]` | **成立** |
| ★r1 §4「深度 3 不拦」 | ★**今天不成立**：直调 `Store.Set` 在 born-nil Store 上被 **367 的门**拒（`perm: …需要一条确认通道`）；接腿后放行且 `calls=0` | **被 367 变旧**（⛔ r1 判错） |

## §2 ★核心判点：常驻线 `AE→AH`（我自己造、真 gate、产线门）

**造法**（⛔ 抄 r2 的 stub）：`ra:=newResidentApprovalWithConfig(dir)`（**真 `ra.gate`**，无球＝fail-closed）→ `leg:=panelModeConfirm(计数包在 ra.confirmModeSwitch 外层)` → **`newResidentPanelManager(dir, leg)`**（＝`resident_windows.go:224` 那一行）→ `pm.dispatchRaw(ctx, 封套)`（＝**页面那一跳** `panel_host_windows.go:948`）。

| 发 | `err` | `disk` | `confirm_calls`／`asked` |
|---|---|---|---|
| 改前（`nil` 腿） | `ErrNoL2Confirm`（逐字同 r2 §1③） | `ask_every_step` | 0／`[]` |
| ★**改后（真腿）** | ★**`nil`** | ★**`ask_high_risk`** | ★**0／`[]`** |
| 正控 `AA`（同门同真腿） | 无球时 `perm: 切到 auto_approve 被拒绝: …审批界面不可达，已 fail-closed 拒绝` | 停在 `ask_high_risk` | **1／`[auto_approve]`** |
| `AA`＋无腿 | `ErrNoL2Confirm` | 未动 | 0／`[]` |

⇒ ★★r2 报的三个数**逐值复现**；我这一发比它硬：腿是**真 `ra.confirmModeSwitch`**，`AA` 证明它**真会去举卡**（且无球时真拒），⛔ 一个只记数的桩。

**判语**：★**「接线做对了」。**
- 门不是空的：同一枚腿在 `AA` 上真出卡（恰 1 张）、无球真拒；`nil` 时两层都真拒（§7 m1–m3）。
- `AE→AH` 的 0 张卡是 **`M4` 的决定**（owner 2026-09-21），⛔ 门失灵；把这算成 r2 的错＝把**规则账**记到**接线账**上。
- ★它欠的那一句（§0.2）是**记账**问题，⛔ 实现问题：改前面板**不能**放宽任何档，改后**能静默放宽一档**。

## §3 `1` 还是 `22`：那 13 行在哪、算不算产线

| 尺 | 行数 | 构成 |
|---|---|---|
| 带 `.scratch`/`docs` 排除（r2 口径） | **9** | 定义 1＋调用 **1**＋注释 7 |
| 不带排除（仍排 `_test.go`） | **22** | 9 ＋ **13 行全在 `.scratch/wisp/probes/**`**：`33/r11/mutations/{baseline,mut-m1…m3,restore}-panel_host_windows.go`、`319/r4/prechange-…`、`255/r6/mutations/baseline-…`、`198/v2/backup/run.go`、`197/r3b/pre/run.go`、`151*/run.{shipped,head,mut-noclose}.go` |

★**判：`1` 是产线调用者数，`22` 是「全仓含探针副本的文本行数」。** 判据三条：①13 枚全是**产码文件的副本**（`baseline-`/`mut-`/`prechange-`/`restore-` 前缀），无一枚是页面/宿主的新入口；②`go list ./...` 里 **`.scratch` 包数＝0**（`.` 开头目录 Go 工具直接忽略）⇒ `go build ./...` **不编译**它们；③把它们算进"调用者"，等于说"探针副本里有同一句"⇒ 自指。r2 的两个数都对、口径写在件里、**结论不变**＝可核。

## §4 「1/2 只成立」＋CLI 仍拒该不该（＋m6 复核）

| 读 | 我的读数 |
|---|---|
| CLI `AE→AH`／`AA` | 都 `ErrNoL2Confirm`、`disk=ask_every_step`；真二进制 `rc=1` |
| **m6**（CLI 塞非 `nil` 腿） | `err=nil`／`disk=ask_high_risk`／`calls=0` ＝**一声不响放宽**；红在 `:160`＋`:163`（两支各一对） |

★**今天"CLI 仍拒"是对的**：管子只有 stdin/stdout/stderr，**没有回答通道**；接一枚非 `nil` 腿只关掉 handler 的 `nil` 门（handler 只读 nil-ness），m6 已把后果跑出来；与调研的 `cline` 支（非交互会话明文拒）同向。
★★**但这是"借来的正确"**：它承重于 **367 的方向判据**（handler `perm.IsWidening`＋Store `IsWidening && s.confirm==nil`）。**A1214 已判 367 反主流、立 372 回退成「只对 `auto_approve` 判」** ⇒ 372 落地后 `AE→AH` 在 CLI 这条线上**两道门都不再拦**（Store 只对 AA 要通道；handler 换回共用谓词），`disk` 会被管子静默改掉，而 **`TestTicket114r2CLIPipeStillRefusesEveryWidening` 会红**。
⇒ ★这条**必须由 `372-r1`/`372-v1` 显式处理**（要么在 CLI 入口按名字拒档位写入，要么明知地改／退役 r2 那枚用例），⛔ 让它在 372 落地时悄悄烂掉。

## §5 `rt.modeWrites` 三条「走不通」逐条核 ＋ 「等价形状」判

| 理由 | 我的核法 | 读数 | 判 |
|---|---|---|---|
| ① `perm.New` 构造时抄 `Confirm`⇒晚绑不可行 | 读码（`store.go:139` `s.confirm=o.Confirm`；`internal/perm` 内**无 setter**）＋**造一发**：born-nil Store＋后置 `h.Confirm=腿` | handler 让路，**Store 拒**：`perm: 把档位 … 需要一条确认通道`；`calls=0`、盘未动 | ★**成立**（"holder"只救 handler、救不了 Store） |
| ② 塞"非 `nil` 但会拒"的闭包⇒永久关门 | **自造**：两层都拿"必拒"非 `nil` 闭包 | `err=nil`／`disk=ask_high_risk`／`calls=0` | ★**成立**（＝ m6 的机制） |
| ③ 无控制台形状里根本没装配 `rt.modeWrites` | 读码＋行号 | `resident_task_source_windows.go:243-249`（`console==nil && injected==""` ⇒ **return nil**）先于 `:284 assembleRuntime`；`rt.modeWrites` 只在 `run.go:692`（`assembleRuntime` 内）。顺序：`resident_windows.go:224`（面板链）< `:367`（task source） | ★**成立**；且非 r2 新说——`resident_approval_windows.go:361-367`（246/265 留注）早把"双击形状"写成同一条 |

**「等价形状」对不对 ⇒ ★对，而且比"等价"更强。** 常驻带控制台形状里 `runSpec.gate=ra.gate`（task source `:275`）→ `rt.gate=s.gate`（`run.go:624`）⇒ 真把 `rt.modeWrites` 交出去，它的 `Confirm`＝`panelModeConfirm(rt.confirmModeSwitch)`＝`confirmModeSwitchVia(rt.gate=ra.gate,…)`，**与 r2 交的 `panelModeConfirm(ra.confirmModeSwitch)` 是同一个 gate 对象、同一个卡体**，差别只在"谁持有它"；而**无控制台形状里 `rt.modeWrites` 不存在**⇒字面形状**在那里没有对象**。适配器 `panelModeConfirm`（`run.go:938` 一处定义）／卡体 `confirmModeSwitchVia`（`run.go:905` 一处定义）各一枚 ⇒「一个卡体、句子不会漂」**是可核的一处定义，⛔ 自陈**。

## §6 `M4` 没被扩的硬证明（三处独立）

1. ★被测树 `TestTicket101ManualSwitchSurvivesRestart` ⇒ **PASS**（`run_mode101_test.go:250` 逐字 `want exactly [auto_approve] (R20/M4)`；该件**不在六件里**）。
2. ★常驻门：`AE→AH` `asked=[]`（**0 张**）／`AA` `asked=[auto_approve]`（**恰 1 张**）⇒ 出来的卡是 AA 那张。
3. ★`AA`＋无腿 ⇒ `ErrNoL2Confirm`、0 张 ⇒ 接腿没把"AA 要卡"变成"AA 免卡"。

## §7 反向判据：六发（我自己造，全在仓外导出树）

| 发 | 突变（仓外树逐字） | `rc` | 红名 | 红行 |
|---|---|---|---|---|
| 基线 | — | **0** | 5/5 绿 | — |
| **m1** | `panel_inbound.go` 两处 `Confirm` 都退回 `nil` | **1** | 常驻门／T3 | `:122`「…still refused with the leg attached」＋`:187`「the leg was refused instead」 |
| **m2** | 只 handler 退回 `nil` | **1** | 常驻门／T3 | `:122`（`ErrNoL2Confirm`）＋`:187` |
| **m3** | 只 Store 退回 `nil` | **1** | 常驻门／T3 | `:122` ★**perm 句**「…需要一条确认通道…（fail-closed 未切）」＋`:187` ⇒ **两处 `nil` 都得接** |
| **m4** | `resident_windows.go:224` 腿改 `nil` | **1** | 装配根 AST 钉 | `:226`「hands the panel chain *ast.Ident instead of the L2 leg」 |
| **m5** | `run.go` `Confirm: panelModeConfirm(confirm)`→`nil` | **1** | run 腿 | `:307`「rt.modeWrites lost its L2 leg…」＋`:310`＋`:313` |
| **m6** | CLI 调用塞非 `nil` 腿 | **1** | CLI 管 | `:160`「must be refused with ErrNoL2Confirm; got <nil>」＋`:163`「the pipe widened the档…」 |

★**还原绿 5/5 PASS（`rc=0`，逐名 PASS）**；★**还原核（时机＝★m6 之后量）**：六件 `sha256`＝`e8f8ce64…`(panel_inbound)／`da77ac93…`(run)／`0b092eb1…`(resident_approval)／`2e0a541d…`(resident_windows)／`b08ad82c…`(panel_resident)／`c622eb5d…`(新用例)，**与 HEAD blob 及 r2 §8 报的六个值逐值同**；`cmp` 六对 **SAME**；`git status --porcelain -- cmd/wisp` ＝ **0 行**（量于 10:37，m6 之后）。⇒ 六发**无工作树突变窗口**。

## §8 四数 ＋ 格式（甲乙两尺）

| 包 | 命令 | 读数 |
|---|---|---|
| `internal/panel`（**live**） | `go test -count=1 -v ./internal/panel/` | `rc=1` ⇒ **RUN 269／PASS 264（顶层 166）／FAIL 5／SKIP 0** ＝ r2 逐值同 |
| `internal/panel`（**被测树**） | 同上，树内 | `rc=1` ⇒ 269／262（顶层 164）／**7**／0；多出的 2 枚＝**树工件**：`TestComposerRenderFixtureTellsTheTruth`、`TestReadGitOnThisRepositoryIsSelfConsistent`（无 `.git`／导出件 CRLF）⇒ 判 r2 用 live |
| `internal/perm` | `go test -count=1 ./internal/perm` | `rc=0` `ok … 0.818s` |
| `cmd/wisp`（本腿定向） | `-run "TestTicket114r2\|TestTicket101" -count=1 -v` | `rc=0` ⇒ **RUN 15／PASS 15／FAIL 0**（顶层 10＋子测 5） |
| `cmd/wisp`（回归定向） | `-run "TestAC9\|TestAC1SessionDispose\|TestTicket255PanelHostBuildsItsWindowOptions"` | `rc=0` ⇒ **7 PASS／0 FAIL／0 SKIP** |

★live 5 红**逐名＝既有集**（`114-r1` §8／`367-a1-census` 同集）：`TestApprovalCardViewJSONKeysMatchFrontendTypes`·`TestComposerContractTypesMatchFrontend`·`TestC21DesignTokensFourWayAgree`·`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`·`TestStreamLogDroppedNamingLedgerIsNotBounded35r8` ⇒ **新增红 0**。
**格式甲**：`gofumpt.exe -l <六件逐名>` ⇒ `rc=0`、**0 行**。
**格式乙**：`git ls-files "*.go"`（**1180 枚**）经 `xargs` 给 `gofumpt -l` ⇒ **34 行**；按 `[\\/]` 归属：`[\\/]cmd[\\/]wisp[\\/]` **0**、`[\\/]cmd[\\/]` **0**、`[\\/]internal[\\/]` **0**（34 行全在 `.scratch\wisp\probes\`）。
★**射程坑（记）**：`gofumpt -l .` 会去扫 `.scratch` 里**故意写坏的探针**⇒ `rc=2`＋54 行语法错（`expected ')', found '{'`）——**那不是格式读数**，⛔ 当"有未格式化文件"读。

## §9 门禁 rc（⛔ 走管道取 rc）

| 门 | 命令 | `rc` |
|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** |
| vet | `go vet ./cmd/wisp` ／ `go vet ./internal/panel` | **0** ／ **0** |
| D22 | `sh scripts/d22scan.sh` | **0**，`clean - no D22 ban violations`；`examined 291 production Go files`；八 scope **249/42/116/25/39/116/639/150**（逐数＝r2 报的、无一降） |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **0**，`VERDICT GREEN`，`roster entries=61`＝`over the hat=61` |

## §10 边界（越界＝0，逐项可核）

- r2 三笔提交的**足迹**：`80376d58f` **6 件**（全在 `cmd/wisp/`）·`fda84b9ca`/`6ea57bbee` **各 1 件**（`docs/evidence/s1/114-r2-impl.md`）⇒ **没有任何一枚在禁改面**。
- `internal/perm/**`：`git diff --numstat 80376d58f~1..6ea57bbee -- internal/perm` ＝ **0 行**；`internal/perm` 最后一次改动仍是 `367-r1`（`7fabccf1b`）⇒ `IsWidening` **没被碰**。
- `internal/memory/**`（`ddlV*`／`196-r2` 护栏／`196-r3`）＝ **0 行**；`cmd/wisp/logsink.go` **未动**（不在六件里、`git status` 对该路径 0 行）。
- `internal/**`／`cmd/wisp`／`scripts/**`／`.github/**`／`docs/PLAN.md`／`docs/specs/**`／`frontend/**`：`git status --porcelain` **0 行**（量于 10:37，m6 之后）。
  ★`.scratch/**` **有他腿的脏**（`git status --porcelain -- .scratch` = **1188 行**，其中 ` M` 4 枚在起手时就在）⇒ **本腿未碰**，也⛔ 把它算进任何读数。名册（长度门内嵌 roster）`61=61` GREEN。

## §11 `AC#3` 判语：哪一半成立、哪一半 ⛔ ＋ 为什么

| `AC#3` 的两半（票面 `:45-46` 逐字） | 判 | 为什么 |
|---|---|---|
| 「请求**自动档**…在真机上必须走 L2 原生卡片（C18），超时 300s＝拒」 | ★**成立（接线层）** | 常驻门 `AA` ⇒ **恰 1 张**、走 `confirmModeSwitchVia`→`ra.gate`→球；无球时 gate 自己 fail-closed 拒（`审批界面不可达`）⇒ 通道、卡体、fail-closed 三段都在；300s 是 `[risk] confirm_timeout_sec` 的构造期值（256 的既有面） |
| 「请求…**更宽档**…必须走 L2 原生卡片」 | ★★**⛔，且不该落** | 满足它**只能扩 `M4`**（owner 2026-09-21 定案）。而 `m4-survey`（108 行、8 对象）判：**6 个有实现的 repo 没有一个是「因缺通道就拒绝放宽档位」**，`M4` 与 dsh/minimax/openchamber 一致 ⇒ **A1214 已裁：`M4` ⛔ 动、367 的方向判据反主流、立 372 回退**。⇒ 这一半的正确处置是**作废**（由 372 把门改回目标档位＋补动作层），⛔ 实现它 |
| 「钉住**审计三档全写**（R20）不因接线而漏」 | ★**成立** | `AE→AH`／`AA`／无变宽三发各自落 `MODE-SWITCH` 行（`result=applied`／`refused-confirm`／`no-change`），盘面与审计同步 |

★**残余（⛔ 算成立）**：`AC#3` 的「**在真机上**」这一层没人证过——本腿与 r2 都没有真窗/真手势授权（`AC#6` 是 `skipped=frontend(owner-delegated)`）⇒ 成立的是**接线**，⛔ **玻璃上的那张卡**。

## §12 具名残留 ＋ 交下一腿

1. ★**r2 §11.1（`AC#3` 字面 vs `M4`）上报正确**：它按硬边界 ⛔ 自行扩 `M4`、把现状写成「`AE→AH` 过门进写入腿但 0 张卡」⇒ 与我的读数逐值同；**裁决见 §11**。
2. ★**r2 §11.4（票面/注释未动）判得对**：写面纪律；那两处已由 `114-r3` 按同一判法收口（票面**追加**前提更正、⛔ 改原文）。
3. ★**r2 §11.3（`rt.modeWrites` 仍无读点）＝真残留**：本腿复核仍无产线读者（只有 `run.go:692` 写点、`panel_mode_confirm_114_windows_test.go:284-293` 的测试读点）。
4. ★★**新残留（交 372）**：§4 那条前瞻——372 回退后 CLI 管对 `AE→AH` 失守、r2 的 CLI 用例会红。
5. ★**没做**：整包 `go test ./cmd/wisp`（`MARK299`，派单禁）／真窗读／POSIX 三合包／`internal/panel` 的 live 7 红中那 2 枚树工件的独立归因（只登记，未深挖）。

## §13 纪律自报

- 产码 **0 字节**；本件是**唯一**仓内写入。变异全在仓外导出树／`-overlay`（探针两枚在 `D:\tmp\...\probe\`），**工作树零突变窗口**；还原核在**最后一发（m6）之后**量（§7）。
- 起手 `git diff --cached --name-status` 空；`add`/`commit` 两处**显式 pathspec（单件字面点名）**，⛔ `-A`/`.`、⛔ amend/reset/rebase/stash/checkout、零 `push`。
- `unset` 只用 `GOTMPDIR`；`cmd/wisp` 全部命令带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`（导出树运行时指向主检出那份 DLL）；临时件只建不删，全在 `D:\tmp\wisp\20261024\114v2\**`，仓内零 `.out`。
- 窗口内他腿在飞（`372-r1`、`114-r3`），`HEAD` 由 `6ea57bbee` 推到 `11790a4ed`；**本件的被测对象是 `6ea57bbee` 那棵树**，⛔ 随 HEAD 漂移。
