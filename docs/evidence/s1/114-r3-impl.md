# 114-r3 写腿交件 —— 过期措辞收口（只四处，⛔ 零行为变更）

**身份**：写腿（实现者）。**起手锚**：`03179ea65`；**收手时 HEAD**＝`11ede6ef5`（窗口内他腿推进，只落 `docs/evidence/s1/m4-survey.md`，与本笔零交集）。
**判据**：四处全是注释／票面文字 ⇒ `go test` 四数必须与改前**逐数相同**（§5）。**写面（只有三处）**：`internal/panel/composer_handlers.go`（只那两处）＋`.scratch/wisp/issues/114-…-must-be-native.md`（只追加一节）＋本件。

## §0 起手核（⛔ 跳步）

| 核什么 | 现跑读数 |
|---|---|
| `git diff --cached --name-status` | 只有 `M docs/reports/pending-and-issues.md` —— **无 `D`/`A`** ⇒ 不触停手。该 `M` 是**别人**的暂存（台账），本笔 ⛔ add、⛔ commit |
| 写面三路径 `git status --porcelain` | **改前全空**（写面内无别人的未提交改动） |
| `git status --porcelain -- cmd/wisp` | **空**（`114-r2` 残余已不在；本笔 ⛔ 动 `cmd/**`） |
| 行号／逐字 | 全部**现跑重取**：派单给的 `:41-45`／`:132` 与现读一致；旧措辞在 HEAD 的 `*.go` 里**恰 1 处**＝`composer_handlers.go:132` |
| 别人的脏（`issues/115-*`、`probes/161/r6/logs/*` 的 `M`、几个 `??` 怪名） | **看见、⛔ 碰、⛔ 清** |

## §1 处① `internal/panel/composer_handlers.go:41-45` —— "Nothing calls this handler yet"

**改前（逐字，5 行）**

```
// WHAT THIS IS NOT: a path. Nothing calls this handler yet, because the
// WebView2 "event -> ParseComposerRequest" hop does not exist in this tree
// (status table ①.1; tickets 33/35 are still ready-for-agent). AC#2 is the
// precondition for that hop, not a substitute for it - the panel still cannot
// change the档, and this file must not be read as if it could.
```

**被否证的三处（本腿现跑，⛔ 转述）**

| 旧句的主张 | 现跑读数 |
|---|---|
| "Nothing calls this handler yet" | **有产线调用者**：`internal/panel/composer_dispatch.go:204` `return "", d.Mode.HandleModeRequest(ctx, req)`（接口 `:73`、实现 `composer_handlers.go:109`） |
| "the WebView2 hop does not exist in this tree" | **那一跳在**：`cmd/wisp/panel_host_windows.go:929` `w.Bind(panelDispatchBinding, …)` → `:930 m.dispatchRaw` → `:948 m.disp.Handle(ctx, raw)` |
| "the panel still cannot change the档" | **CLI 管也到同一枚 router**：`cmd/wisp/panel_inbound.go:173` `reply, err := disp.Handle(ctx, raw)`；`ComposerDispatch.Handle`＝`composer_dispatch.go:176` |

⇒ 两条产线入口（页面门＋`wisp panel-inbound` 管）都到同一枚 `ComposerDispatch.Handle`，`dispatch` 在 `:204` 分流。

**改后（逐字，5 行）**

```
// WHO CALLS THIS NOW: the earlier "nothing calls this handler yet" was true
// while the inbound hop did not exist, and it is not true now. Both production
// doors reach this handler through one router - the page's door at
// panel_host_windows.go:929/948, `wisp panel-inbound` at panel_inbound.go:173 -
// and ComposerDispatch.dispatch routes the mode case here (composer_dispatch.go:204).
```

**为什么恰好仍是 5 行**：165 行 → **165 行**（§7）。零行位移是刻意的——`HandleModeRequest` 仍 `:109`、门仍 `:131`、
另一处 `:132`、`Set` 仍 `:141` ⇒ 别处写着的跨文件行号引用**没被本笔推歪**（`panel_inbound.go:37` 引 `:126`、
`resident_approval_windows.go:98` 引 `:131`、`composer_dispatch.go:10` 引 `:141`，改前改后同号；那三条是否过期⛔ 本笔射程）。

## §2 处② `internal/panel/composer_handlers.go:132` —— 括注词面漂移（★字符串，非注释）

**改前（逐字）**

```go
		refused := fmt.Errorf("%w: 请求 %q 要把档位从 %s 放宽到 %s，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）",
```

**改后（逐字）**

```go
		refused := fmt.Errorf("%w: 请求 %q 要把档位从 %s 放宽到 %s，而本机没有接入 L2 确认腿（R20/M4：只有全自动要确认）",
```

**为什么那四字是错的**：机主原话是「**只有全自动要确认**」（`pending-and-issues.md:3019` 裁定 `R20` 逐字）；「只有变宽要确认」在
`PLAN.md`／`specs`／台账 **0 命中** ⇒ 那枚括注把**方向判据**（本文件的变宽闸）写成了**产品规则**（R20/M4 只对 `auto_approve` 收一次 L2）。

**★先查"会不会破某枚测试的红行断言"（派单指定的停手条件）——现跑读数**

| 尺 | 读数 |
|---|---|
| `git grep -n "只有变宽" HEAD -- '*.go'` | **恰 1 处**（`composer_handlers.go:132`，就是被改这句） |
| `git grep -n "只有变宽" -- '*.go'`（改后） | **0 处**；同尺 `只有全自动` ⇒ 2 处（本句 ＋ `internal/perm/store_test.go:16` 的注释，**同侧**） |
| `git grep -n "只有变宽\|只有全自动\|R20/M4：\|而本机没有接入 L2 确认腿" -- '*_test.go'` | **1 行**且**不是断言**：`store_test.go:16` 注释「M4 只有全自动吃一次 L2」 |
| `git grep -n "strings.Contains" -- 'internal/panel/*_test.go' 'cmd/wisp/*_test.go'` | **795 行**（本仓测试大量读盘断字面）⇒ 单看它无判定力，须逐条落回这枚句子：**无一处**把该句当被断言的串 |
| 依赖该拒绝句的三处，逐条落点 | `composer_handlers_test.go:156` 断 `audit.joined()` **含 requestId**；`cmd/wisp/panel_inbound_33_test.go:147` 是**否定**断言（不含 `panel: MODE-REFUSED`）；`errors.Is(ErrNoL2Confirm)`（`composer_handlers_test.go:144`／`composer_dispatch_test.go:184`／`cmd/wisp/panel_mode_confirm_114_windows_test.go:105/:159`）断的是**哨兵**，与文案无关 |
| 有没有测试按路径读这枚文件断字面 | `git grep -n "ReadFile.*composer_handlers\|\"composer_handlers.go\"" -- '*.go'` ⇒ **0 命中** |

⇒ **结论：没有红行断言依赖该串 ⇒ 可改，不触停手上报**（派单指定的停手条件：查过，没有）。

## §3 处③ 票 114 票面三处过期句（⛔ 抹原文）

**改法**：**只追加一节**（`## ★前提更正（2026-10-24，凭 114-r1/114-r2 现跑）`，46 增／**0 删**），原 118 行逐字未动。
⛔ 未动 `Status:`、⛔ 未动票名、⛔ 未加 `-done`、⛔ 未勾任何框。

| 位置 | 原文（逐字，仍在原地） | 现跑读数 |
|---|---|---|
| 标题 `:1` | 「`ParseComposerRequest` 在**生产里零调用者**」 | **≠ 0**：`internal/panel/composer_dispatch.go:178` 是调用（定义 `bridge.go:134`） |
| `:21` | 「`ParseComposerRequest` 的**生产调用者 = 0**，`perm.Store.Set` 的**生产调用者 = 0**」 | **两半都不是 0**：①`:178`；②`git grep -n "\.Set(ctx" -- '*.go' ':!*_test.go'` ⇒ 全仓**恰 1 处** `composer_handlers.go:141`（定义 `internal/perm/store.go:189`；`internal/tools/mode.go:20-22` 自陈写侧不在 tools） |
| `:44` | 「⚠ 门已立，**通路仍不存在**（P0 地界：宿主那一跳留票 33/35）」 | **通路已在**：`panel_host_windows.go:929/930/948` ＋ `panel_inbound.go:173` ⇒ 同一枚 `Handle`（`:176`）⇒ `dispatch` `:204` |

**原文当时为什么为真**（写进更正节，留史）：2026-09-21 立案时入向那一跳整层不存在，三个计数确实都是 0、票 33/35 确实仍是
ready-for-agent（`114-ac1-status-table.md` ①.1 量的正是那个状态）；**那一跳后来真落地了**。更正节同时写死"⛔ 读成票 114 变小了"：
①`Confirm` 的 nil 语义一字未变（`:131`），`wisp panel-inbound` **仍递 nil**（`panel_inbound.go:236`，理由 `:33-40`）⇒ 管这一路变宽仍被拒；②AC#3 仍未做。

## §4 处④ `cmd/wisp/logsink.go` 那条 `INTERIM` —— **判语：无需改**

派单："只现跑核它还准不准（`132-r2` 已改过），若已准 ⇒ 如实写无需改、⛔ 硬改"。逐条现跑核：

| 现文的主张（`:44-50`） | 现跑读数 | 判 |
|---|---|---|
| 已不叫 `INTERIM`：标「SEALING, formerly INTERIM (ticket 117 AC#3)」，正文说"deferred question is answered" | `logsink.go:44` 逐字如此 | 准 |
| `Q-31` 已批（日志＝私有数据、要锁） | 台账 `:3463`「`Q-31`＝算、要锁 ⇒ 新建票 132」；`:23530` 记 `132-r2` 改后标签与正文 | 准 |
| `InitLog` 经 `winsec.PrivateDirAll` 建目录 | `internal/observe/logging.go:87` `func InitLog` → `:99 winsec.PrivateDirAll(cfg.Dir, 0o700)` | 准 |
| `sealLogFile` 封每一枚 rolled 文件 | `logging.go:223` `var sealLogFile = winsec.SealFile`；`:316` `sealLogFile(name)` | 准 |
| 「nothing here seals or un-seals the log file」（本件只管落点） | `grep -n "PrivateDirAll\|SealFile\|sealLogFile\|SealDir" cmd/wisp/logsink.go` ⇒ **只命中注释** `:29/:47/:48`，**零调用** | 准 |

⇒ **无需改**；本笔 **`cmd/**` 一字节未动**（派单第三条停手条件因此不触发）。

## §5 四数（改前／改后逐数相同）—— `internal/panel` 的 5 枚既有红分开报

**尺**：`cd <repo>; unset GOTMPDIR; export GOFLAGS=; go test -count=1 -v ./internal/panel/ > out 2>&1; rc=$?`

| 时点 | RUN | PASS | FAIL | SKIP | rc |
|---|---|---|---|---|---|
| **改前**（原树） | **269** | **166** | **5** | **0** | 1 |
| **改后**（本笔） | **269** | **166** | **5** | **0** | 1 |

**5 枚既有红逐名（两跑同集；去掉耗时后 `diff` rc=0＝SAME SET）**——与本笔无关，⛔ 未修、⛔ 未放宽：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestC21DesignTokensFourWayAgree`（3 枚前端耦合：
共享工作树里 `design/assets/*` 被别家会话删且未 staged）、`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（2 枚 35-r8）。

**`cmd/wisp` 定向**（派单禁整包：`MARK299` 死区 ⇒ 只跑定向、必须带 sherpa PATH）：`export PATH="$PWD/third_party/sherpa-onnx:$PATH";
go test -count=1 -run "TestTicket114r2" ./cmd/wisp/` ⇒ **rc=0**、`ok … 9.581s`（5 枚：resident 链／CLI 管仍拒每一发变宽／
非 nil 腿会解锁管／常驻根递 orb 腿／run 缝同一条腿）。**⛔ 未跑整包**。

## §6 门禁 rc（逐条现跑）

| 门 | 命令要点 | rc | 输出 |
|---|---|---|---|
| build | `go build ./...` | **0** | 空 |
| vet | `go vet ./internal/panel/` | **0** | 空 |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/panel/composer_handlers.go` | **0** | **0 行**（绝对路径；⛔ 拿 `gofumpt -l` 判 HEAD 是不是 CRLF） |
| d22scan | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；scope：`ban #8 internal/=639`（含注释与 `_test.go`）／`cmd/=150`／`design/=39`／`frontend/=116`；`bans #1-5 internal/=249`／`cmd/=42`；`ban #6 frontend/=116`；`ban #7 tools/=25` |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；pre-stage `tracked=17680  over-budget=61  rostered=61  not in roster=0`；**staged 后**（本件入册）`tracked=17681`，over-budget 仍 61／全 rostered |

★本笔在 `.go` 里**零新增符号**：没有 `★`／`⛔`／`⇒`／带圈数字进 `.go`（改后注释全 ASCII；字符串只换 4 个字）——d22scan `ban #8` rc=0 即仪器读数。

## §7 `git diff --numstat` 读数

| 文件 | +/− | 说明 |
|---|---|---|
| `internal/panel/composer_handlers.go` | **6 / 6** | 只两处：`:41-45` 注释 5 行＋`:132` 字符串 1 行；**165 行 → 165 行** |
| `.scratch/wisp/issues/114-…-must-be-native.md` | **46 / 0** | 只追加"★前提更正"一节；原 118 行逐字未动（`:1`／`:21`／`:44` 与 HEAD 逐行比对 **IDENTICAL**） |
| `docs/evidence/s1/114-r3-impl.md` | **160 / 0**（新增） | 本件；`git diff --cached --numstat` 逐字：`160	0	docs/evidence/s1/114-r3-impl.md` |

**sha256（收手回执，⛔ 跨形比同一性）**：`git show HEAD:… | sha256sum` ⇒ `3558ee2b943c9df0ccd04115873ba708e757c31cc407073fdcbdaa8cb2e43575`（git blob 形，LF）；
改后工作树文件 ⇒ `3565cce4424f97230ab0b5d324ea23025ee2fab3f91d243a8e08081e2a1ce8b5`。两者**必然不同**（内容改了）。

## §8 边界自报（越界＝作废）

- **⛔ 改任何产码行为**：一段注释＋一枚错误句的括注；无控制流／新符号／导入／API 变化；四数逐数同（§5）。**⛔ 动 `internal/perm/**`**（`367-r1` 地界）／**⛔ 动 `cmd/**`**（含 `logsink.go`，判语见 §4）：**各零字节**。
- **⛔ 顺手改别的过期句**：只这四处。具名"看见没动"（留下一波）：`cmd/wisp/panel_inbound.go:37`（引 `:126`）、
  `cmd/wisp/resident_approval_windows.go:98`（引 `:131`）、`internal/panel/composer_dispatch.go:10`（引 `:141`）、
  `internal/panel/pump.go:15-16`（`ticket 33 is unclaimed`）。
- **⛔ 别的面**：`scripts/**`、`.github/**`、`tools/d22scan/**`、`go.mod`/`go.sum`、`docs/PLAN.md`、`docs/specs/**`、
  `docs/reports/**`（含别人已 staged 的台账 `M`——⛔ add、⛔ commit）、`frontend/**`、`design/**`（**零读零写零转述**）、
  `thresholds.go`／golden／`allowlist.txt`／三枚冻结件——**全部零字节**。
- **⛔ 翻框／改 `Status:`／改票名／加 `-done`／改台账／push**：**全部未做**。本腿**只 commit**（显式 pathspec、逐个文件具名）。
