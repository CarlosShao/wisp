# 票 114 接线腿 `114-r2` — 把 L2 确认腿做成 `newComposerDispatchChain` 的**参数**，常驻面板链接上本进程审批门

**腿** `114-r2`（写腿／实现者）· **产码提交** `80376d58f`（6 件）· **起手锚** `127784149`
**写面**：`cmd/wisp/run.go`、`cmd/wisp/panel_inbound.go`、`cmd/wisp/panel_resident_windows.go`、`cmd/wisp/resident_windows.go`、`cmd/wisp/resident_approval_windows.go`、新件 `cmd/wisp/panel_mode_confirm_114_windows_test.go`；本件。
⛔ `internal/panel/**`、⛔ `internal/perm/**`、⛔ 票面／`Status:`／台账／`logsink.go`、零 push。

## §0 三句结论（先给判据）

1. ★**常驻面板那条产线**（`panel_host_windows.go:948` 的 `dispatchRaw`）★**从"变宽一律拒"变成"变宽过门进写入腿"**；切 `auto_approve` ★**恰好一张卡**，卡走**本进程审批门**（悬浮球 `ballCardUI`）。
2. ★★**`AE→AH`（`ask_every_step→ask_high_risk`）一张卡都不出**，因为 R20/M4 只对切进 `auto_approve` 出卡 ⇒ ★`run_mode101_test.go:250` 的 `want exactly [auto_approve]` **仍绿**＝**本腿没扩射程**的硬证明（§4）。
3. ★★**两处与派单字面不一致，本腿按硬边界处理并具名上报**：★①**`rt.modeWrites` 那条字面形状走不通**（`perm.New` 构造时抄走 `Confirm` ⇒ 晚绑不可行，§5），落的是**等价形状**（同一个适配器＋同一个卡体）；★②**CLI 那条产线仍拒**（stdin/stdout 管子上没有卡面，塞非 nil 腿＝把 AC#2 的门关掉，§6）⇒ ★**"一处改动、两条产线同时能走卡"只成立 1/2**。

## §1 `114-r1` 三条读数（本腿现跑，⛔ 抄）

**① `ParseComposerRequest` 产线调用者 ＝ 1**
`git grep -nE "ParseComposerRequest" -- "*.go" ":(exclude)*_test.go" ":(exclude).scratch/**" ":(exclude)docs/**"` ⇒ **9 行**：
`internal/panel/bridge.go:134` **定义** ／ `internal/panel/composer_dispatch.go:178` **调用** ／ 余 **7 行 `//` 注释** ⇒ **调用者 1**。
★尺的射程自查（`ruler-scope-checklist` 第 1/7 例）：同一模式**不加 `.scratch`/`docs` 排除**时是 **22 行**（13 行是 `.scratch/wisp/probes/**` 的探针副本）⇒ 数会变、**结论不变**；`114-r1` 记的"10 行/8 注释"是它的口径，本腿报自己的 9/7。

**② 绕过深度 1／深度 2 都拒、正控过**（`-overlay` 探针、仓外 `D:\tmp\wisp\20261024\114r2\probe\`、工作树零突变）：

| 发 | 读数（逐字） |
|---|---|
| 深度 2（走产线派发表） | `err="panel: 确认腿未接入，档位不得变宽: 请求 \"probe-depth2\" 要把档位从 ask_every_step 放宽到 ask_high_risk…"` `isErrNoL2Confirm=true` `disk="ask_every_step"` |
| 深度 1（跳派发表直调处理器） | 同句 `probe-depth1`、`isErrNoL2Confirm=true`、`store_reads=ask_every_step`、`disk="ask_every_step"` |
| **正控**（只把 `Confirm` 接上） | `err="<nil>"` `confirm_calls=0` `store_reads=ask_high_risk` `disk="ask_high_risk"` |
| **正控＋AA**（同一枚处理器） | `err="<nil>"` `confirm_calls=1` `asked=[auto_approve]` `disk="auto_approve"` |

**③ 逐字拒绝句**（深度 1／2／真二进制同句）：

```
panel: 确认腿未接入，档位不得变宽: 请求 "probe-depth1" 要把档位从 ask_every_step 放宽到 ask_high_risk，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）
```

真二进制复读（`bin\wisp.exe panel-inbound -data …`，DLL 在 PATH 上）：`AE→AH` **rc=1**、`AA` **rc=1**、两发后 `config.toml:4` 仍是 `permission_mode = 'ask_every_step'`。

## §2 改了什么（逐行）

| 件 | 改法 | 为什么 |
|---|---|---|
| `panel_inbound.go:263` | `newComposerDispatchChain(dir, auditf, actor, ws, confirm panel.ModeConfirm)` — 新增**第 5 个参数** | 腿只能在**构造时**给（§5），所以它是参数，不是事后可填的字段 |
| `panel_inbound.go:277` | `perm.New(perm.Options{Confirm: storeConfirmFromPanel(confirm)})`（原 `Confirm: nil`） | Store 的 `IsWidening && confirm==nil` 门拿**同一枚**腿 |
| `panel_inbound.go:288` | `ModeWriteHandler{Confirm: confirm}`（原 `Confirm: nil`） | handler 的 `ErrNoL2Confirm` 门拿**同一枚**腿；它只读 nil-ness、从不调用 |
| `panel_inbound.go:345` | 新 `storeConfirmFromPanel(panel.ModeConfirm) perm.ConfirmFunc`：**nil 进 nil 出**，否则包成 `perm.Switch` | 两层要同一条腿的两种形状；把 nil 换成"调用时才拒"的闭包会**把 Store 的门关掉** |
| `panel_inbound.go:236` | `newComposerDispatchChain(…, nil, nil)`＋注释 | CLI 腿**照旧 nil**，并写成具名理由（§6） |
| `run.go:699` | `rt.modeWrites` 的 `Confirm:` 由内联闭包改为 `panelModeConfirm(confirm)` | 适配器只有一枚 |
| `run.go:893/905` | 新 `confirmModeSwitchVia(gate *approval.Gate, ctx, sw)`；`rt.confirmModeSwitch` 改成一行转调 | **一个卡体**给两条腿用，句子不可能漂成两句 |
| `run.go:938` | 新 `panelModeConfirm(perm.ConfirmFunc) panel.ModeConfirm`（nil 进 nil 出） | 同上 |
| `run.go:331-345` | 字段注释重写（原文 "Nothing reads this field yet" 保留为事实，但把 114-r2 的另一端写清） | ⛔ 让注释继续误导 |
| `resident_approval_windows.go:111` | 新 `(*residentApproval).confirmModeSwitch` ＝ `confirmModeSwitchVia(ra.gate, …)` | 常驻腿的卡落在**本进程的门**（悬浮球）上 |
| `resident_windows.go:224` | `newResidentPanelManager(rt.Layout.DataDir, panelModeConfirm(ra.confirmModeSwitch))` | 装配根**交腿**的那一行 |
| `panel_resident_windows.go:228/263` | 两枚构造函数新增 `confirm panel.ModeConfirm` 并往下传；"Confirm stays nil on purpose" 那段注释改写 | 链是这条腿的落点 |

## §3 两条产线：`AE→AH` 改前／改后

| 产线 | 改前（`Confirm: nil`） | 改后（接上腿） |
|---|---|---|
| **常驻面板**（`newResidentComposerDispatch` ＋ `dispatchRaw`） | `ErrNoL2Confirm`（逐字同 §1③）、`disk=ask_every_step`、`confirm_calls=0` | **`err=nil`**、**`disk=ask_high_risk`**、**`confirm_calls=0`**（M4 不对它出卡） |
| **CLI 管**（`newPanelInboundDispatch`） | `ErrNoL2Confirm`（真二进制 rc=1） | ★**仍是 `ErrNoL2Confirm`**（真二进制 rc=1、盘未动）；理由与实测见 §6 |

★常驻那条的"改后"是**经产线门**读的：`PanelManager.dispatchRaw`（`panel_host_windows.go:941-949`），不是测试自搭的旁路。

## §4 M4 没被扩的硬证明（三处独立读数）

1. ★`go test ./cmd/wisp/ -run "TestTicket114r2|TestTicket101"` ⇒ **`TestTicket101ManualSwitchSurvivesRestart` PASS**（`run_mode101_test.go:250-252` 的 `want exactly [auto_approve] (R20/M4)` 未改一字）。
2. ★常驻门那一发 `AE→AH`：`confirm_calls=0`，且**AA 那一发** `asked=[auto_approve]`（len=1）⇒ ★**出来的卡是 `auto_approve` 那张，⛔ 不是 `AE→AH` 的**。
3. ★探针正控（§1②）与 run 腿（`TestTicket114r2RunSeamPanelHandlerStillCarriesTheSameLeg`）**同一个读数**：`asked=[auto_approve]`、`calls=1`。

## §5 为什么不是字面的"把 `rt.modeWrites` 交给链"（★裁决点）

`perm.New` 在**构造时**把 `Confirm` 抄进 `Store`（`store.go:139`），`ModeWriteHandler` 每次请求只读**自己那份**的 nil-ness（`composer_handlers.go:131`）⇒
- **晚绑**（holder／bind 之后再给）**不可行**：Store 永远留着 nil ⇒ 所有变宽被拒；
- **塞一个非 nil 但会拒的闭包**＝把 handler 的门**永久关掉**（§6 实测那一发就是它的后果）；
- 而 `rt.modeWrites` 由 `assembleRuntime` 装配，**晚于**常驻面板链（`resident_windows.go:224`），且在**无控制台启动形状**里根本不会装配（`startResidentTaskSource` 的 early return）。

⇒ ★腿只能是**构造参数**，常驻那条只能用**boot 就存在**的门：`ra.gate`。本腿用的是**同一个适配器**（`panelModeConfirm`）＋**同一个卡体**（`confirmModeSwitchVia`），所以"面板那张卡"和"`rt.modeWrites` 会出的那张卡"是同一段代码，⛔ 两个相似物。
★残留（具名）：`rt.modeWrites` 的**字段仍无读点**；本腿没删它、没让它说谎（注释已改准）。

## §6 CLI 那条线：为什么保持 nil（可执行的理由）

- 管子只有 stdin（封套流）／stdout（回执）／stderr（审计），**没有回答通道** ⇒ 任何"腿"都不是卡。
- ★`TestTicket114r2ANonNilLegAtTheCLISeamWouldUngateThePipe` **跑了那个替代**：给这条链一枚**非 nil**腿（内容＝立刻报错）⇒ handler 让路、`Set` 收下、`disk=ask_high_risk`、`confirm_calls=0` ⇒ ★**管子会一声不响地把档位放宽**。这就是 AC#2 那道门存在的理由。
⇒ ★`panel_inbound.go` 那个 nil **是门、不是桩**；把它"补上"＝把门关掉。

## §7 新用例逐名（★每枚今天响、不修不响）

| 用例 | 今天为什么响 |
|---|---|
| `TestTicket114r2ResidentPanelChainFromRefusalToTheWriteLeg` | 摘掉腿 ⇒ 红在 `:122`（§8 m1/m2/m3 三发各自点名） |
| `TestTicket114r2CLIPipeStillRefusesEveryWidening` | 给 CLI 那条链塞腿 ⇒ 红在 `:160`／`:163`（§8 **m6** 实测：管子把档位放宽了） |
| `TestTicket114r2ANonNilLegAtTheCLISeamWouldUngateThePipe` | 链忽略参数（＝腿被还原掉）⇒ 红在 `:187`（§8 **m1 打在 T3 上**实测：`the leg was refused instead`） |
| `TestTicket114r2TheResidentAssemblyRootHandsTheChainTheOrbLeg` | 装配根把腿换回 `nil` ⇒ 红在 `:226`（§8 m4，仓外树真种盘） |
| `TestTicket114r2RunSeamPanelHandlerStillCarriesTheSameLeg` | `rt.modeWrites.Confirm` 变 nil ⇒ 红在 `:307`（§8 m5） |

## §8 反向判据：六发变异（全部仓外／`-overlay`，工作树零突变）

| 发 | 突变处（仓外副本逐字） | rc | 红名 | 红行逐字 |
|---|---|---|---|---|
| 基线 | — | **0** | 5/5 绿 | — |
| **m1** | `panel_inbound.go` 两处 `Confirm:` 都退回 `nil`（＝改前形状） | **1** | 常驻门那枚／T3 | `:122 a widening that M4 does not card is still refused with the leg attached: panel: 确认腿未接入…`；打在 T3 上 ⇒ `:187 …the leg was refused instead: panel: 确认腿未接入…` |
| **m2** | **只**把 handler 的 `Confirm:` 退回 `nil`（Store 留腿） | **1** | 常驻门那枚 | `:122 … panel: 确认腿未接入…`（handler 的门先拒） |
| **m3** | **只**把 Store 的 `Confirm:` 退回 `nil`（handler 留腿） | **1** | 同上 | `:122 … perm: 把档位从 ask_every_step 放宽到 ask_high_risk 需要一条确认通道，但本机没有接入（fail-closed 未切）` ⇒ **两处 nil 都得接** |
| **m4** | 仓外树 `resident_windows.go:224` 的腿改成 `nil` | **1** | 装配根 AST 钉 | `:226 the resident boot hands the panel chain *ast.Ident instead of the L2 leg…` |
| **m5** | `run.go` 的 `Confirm: panelModeConfirm(confirm)` → `nil` | **1** | run 腿 | `:307 rt.modeWrites lost its L2 leg: panelModeConfirm dropped a usable confirm…` |
| **m6** | `panel_inbound.go:236` 的 CLI 调用改成塞一枚非 nil 腿 | **1** | CLI 管那枚 | `:160 …must be refused with ErrNoL2Confirm; got <nil>` ／ `:163 the pipe widened the档 to "ask_high_risk" with no card anywhere` |

★**三读数**：基线绿 **5/5** ／ 变异红（上表逐名）／ ★**还原绿**（仓外树 `mut/m4-restored.txt`：`go test -run TestTicket114r2` ⇒ **5/5 PASS, rc=0**）。
★**还原核（时机＝★全部六发变异之后量；m1–m3/m5/m6 走 `-overlay`、m4 走仓外树，两者都不写工作树；本行数值在 m6 之后又逐件复量一次、逐值同）**：`sha256sum` 六件 ＝ `e8f8ce64…`(panel_inbound.go) `da77ac93…`(run.go) `0b092eb1…`(resident_approval_windows.go) `2e0a541d…`(resident_windows.go) `b08ad82c…`(panel_resident_windows.go) `c622eb5d…`(新用例)；仓外树还原件与 pristine 逐字同值 `2e0a541d…`；`git status --porcelain -- cmd/wisp` ＝ **0 行**。

## §9 四数

| 包 | 命令 | 读数 |
|---|---|---|
| `internal/panel` | `go test -count=1 -v ./internal/panel/` | **rc=1** ⇒ `RUN 269 / PASS 264（顶层 166）/ FAIL 5 / SKIP 0` |
| `internal/perm` | `go test -count=1 ./internal/perm` | **rc=0** ⇒ `ok … 1.215s` |
| `cmd/wisp`（本腿定向） | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp/ -run "TestTicket114r2\|TestTicket101" -count=1 -v` | **rc=0** ⇒ `RUN 15（顶层 10、子测 5）/ PASS 15 / FAIL 0`（`ok … 171.687s`） |
| `cmd/wisp`（回归定向） | 同上 `-run "TestAC9\|TestAC1SessionDispose\|TestTicket255PanelHostBuildsItsWindowOptions"` | **rc=0** ⇒ **7 PASS / 0 FAIL / 0 SKIP** |

★`internal/panel` **5 枚红＝既有集**（与 `114-r1-impl.md` §8、`367-a1-census.md` §⑦ **逐名同集**，本腿零字节改动这些文件）：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`（frontend 耦合）· `TestComposerContractTypesMatchFrontend`（frontend 耦合）· `TestC21DesignTokensFourWayAgree`（`design/assets/tokens.css` 不在本树）· `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` · `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`。★**本腿新增红＝0**。

## §10 门禁 rc

| 门 | 命令 | rc |
|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** |
| vet（含测试件） | `go vet ./cmd/wisp` | **0** |
| 格式 | `/d/work/base/gopath/bin/gofumpt.exe -l`（六件逐名） | **0**，输出 **0 行** |
| D22 扫描 | `sh scripts/d22scan.sh` | **0**，`clean - no D22 ban violations`；`examined 291 production Go files`、`bans #1-5 internal/=249`、`#1-5 cmd/=42`、`ban #6 frontend/=116`、`ban #7 internal/tools/=25`、`ban #8 design/=39`／`frontend/=116`／`internal/=639`／`cmd/=150` —— 八 scope 逐数不降（`114-r1` 记 248/42/…/290，差额是他票自锚点以来长出来的件，非本腿） |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **0**，`VERDICT GREEN`（61 名册＝树） |

## §11 具名没做 ＋ 请裁决（⛔ 自行假设）

1. ★★**AC#3 字面 vs M4（最重要的开放点）**：票面 AC#3 逐字要"请求**自动档／更宽档**在真机上必须走 L2 原生卡片"；而 R20/M4 只对切进 `auto_approve` 出卡。本腿按派单**硬边界①**⛔ 扩 M4 ⇒ 现状是 ★**`AE→AH` 过门进写入腿但 0 张卡**。★"更宽档也要卡"＝**扩 M4 才做得到**＝派单停手条件②，故**不落、交裁决**。
2. ★**"两条产线同时能走卡"只成立 1/2**（§6 实测理由）。
3. ★`rt.modeWrites` 仍是"像路但不是路"（§5 残留）。
4. ★**票面三处过期句**（标题／`:21`／`:44`）＋`internal/panel/composer_handlers.go:41-45` 的过期注释：前者台账 `pending-and-issues.md` 记"归 `114-r2`"，但本腿派单**写面不含 `.scratch/wisp/issues/**`**且明令⛔改票名／`Status:`；后者在 `internal/panel/**`（本腿禁改）⇒ **两处都没动**，请编排者分派。
5. ★**未跑**：整包 `go test ./cmd/wisp`（`MARK299` 硬终止，派单禁）／真窗常驻读（WebView2，需真机手势，本腿无此授权）／POSIX 三合包。

## §12 纪律自报

- 起手 `git diff --cached --name-status` **空**；`cmd/wisp` 起手无他腿脏改动（工作树里的 `M`／`??` 怪名一律未碰）；`add`/`commit` 两处**字面点名**六件，⛔ `-A`、⛔ `.`、零 push、零 amend/reset。
- `unset` 只用 `GOTMPDIR`；`cmd/wisp` 一律带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`（首跑 `rc=127` 缺 DLL＝假失败，已带 PATH 重跑）。
- 临时件只建不删，全在 `D:\tmp\wisp\20261024\114r2\**`（含 `bin/`、`mut/`、`probe/`、`tree-m4/`、`*.txt` 读数）；仓内零 `.out`、零探针残留（§8 还原核）。
