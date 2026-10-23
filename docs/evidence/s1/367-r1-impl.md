# 367-r1 落地腿 —— `perm.Store.Set` 的门改判**方向**（写腿 / AC#2）

- 腿 `367-r1`（**写腿**）。HEAD `f3a4055d2`（叠在普查腿 `cc45a4b81` 之上），量于 2026-10-23 18:2x–18:5x +08。
- **形＝〔甲′〕**（编排者拍）：`Set` 的门从「目标是不是 `auto_approve`」改成「**是不是变宽**」，判据**复用** `modeIsWidening` 那一族谓词（下沉后＝`perm.IsWidening`），**⛔ 另写一份**。
- 写面：`internal/perm/**`、`internal/panel/**`、两包测试、本文件。**`cmd/wisp/**` 改动 0 字节**（两条产线路用 `-overlay` 在**仓外**驱动真产线装配去量）。
- 突变/探针全在仓外 `D:\tmp\wisp\20261023\367r1\**` ＋ `-overlay` ⇒ **无工作树突变窗口**。
- 开工前核过：`git diff --cached --name-status` **空**（无 `D`/`A`）；`git status --porcelain internal/perm internal/panel cmd/wisp` **空**（写面内无别人的未提交改动）。

## ① 普查核心读数：我自己复跑的（⛔ 抄 `367-a1-census.md`）

尺＝**同一支探针的 A/B**：`D:\tmp\…\probe\perm_probe_test.go` 经 `-overlay` 挂进 `internal/perm`；`before` 那一跑**同时**把 `internal/perm/store.go` 换成 `git show HEAD:internal/perm/store.go`、把 `direction.go` 置空（overlay 的 `""` ＝ 删该文件）⇒ **同一支探针、同一批 `R1|` 行、只换被测代码**。复跑 `go test -overlay=…/overlay-perm-{before,after}.json -run TestR1 -v ./internal/perm`（rc 0 / 0）。
**梯子是我现跑枚举的**：`R1|LADDER|valid=[0=ask_every_step 1=ask_high_risk 2=auto_approve]|names_len=3|default=ask_every_step` ⇒ **恰好三档、无第四档**。

## ② 依赖方向（先核、⛔ 假设）

| 尺（现跑） | 读数 |
|---|---|
| `go list -f '{{.ImportPath}} -> {{join .Imports " "}}' ./internal/perm` | `context fmt config risk sync time` ⇒ **perm 不 import panel、不 import 任何 UI 包** |
| `git grep "wisp/internal/perm" -- internal/panel` | **0 命中** ⇒ panel **今天不** import perm |
| `git grep "wisp/internal/panel" -- internal/perm` | **0 命中** |
| 全仓 import `internal/perm` 的产线包 | 只有 `cmd/wisp`（`panel_inbound.go`/`run.go`）＋ `internal/projctx` 的**测试** |

⇒ **perm 是下层叶子、panel 是上层**：`perm → panel` 是反向（⛔ 用），`panel → perm` 是下行（本次新加的唯一一条边）。

## ③ 谓词下沉到哪、为什么

- **下沉到 `internal/perm/direction.go`**：`func IsWidening(from, to risk.Mode) bool { return to > from }`，**一份定义、两个调用点**（`perm.Store.Set` 与 `panel.ModeWriteHandler.HandleModeRequest`）。
- **为什么不在 `panel`**：`perm` 不能 import `panel`（病根就是"门装在一层、另一层看不见"）；放 panel 就只能一层看得见。
- **为什么不新造第三包／不放 `internal/risk`**：`risk` 不在本腿写面（⛔ 动）；且谓词是**权限政策**（"放宽安全姿态"），语义归属就是 `perm`——`risk.Mode` 只提供"常数有序"这个前提，不提供"宽/严"这个政策。
- ⇒ `modeIsWidening` 从 `composer_handlers.go` **删除**；`:133` 改叫 `perm.IsWidening(from, to)`；panel 那枚梯子钉子改成驱动**同一个**函数（两包各一枚钉子、同一组 `risk.ModeNames()`）。
- ★代价具名：`panel` 的依赖面**动了**。`composer_handlers.go` 里 "`ModeWriter` … with no import of internal/perm, so this package's dependency face does not move" 说的是**接口**（仍是结构化声明、仍不 import perm），它对；但**谓词这一处必须共享**，否则又是两个判据。该处头注释已同步改写。
## ④ 6 对「变宽 → 门」表（改前 → 改后，`Confirm` 全 nil）

| 从 → 到 | 方向 | 改前 `:191-192` 逐字 `if to == risk.ModeAutoApprove {` ＋ `if s.confirm == nil {` | 改后 `:209` 逐字 `if IsWidening(from, to) && s.confirm == nil {` |
|---|---|---|---|
| `ask_every_step → ask_high_risk` | **变宽** | **不拦** `applied=1`、文件真落 `ask_high_risk`、`detail=""` | **拦** `applied=0` `last_result=refused-confirm` |
| `ask_every_step → auto_approve` | 变宽 | 拦 `refused-confirm applied=0` | 拦（同） |
| `ask_high_risk → auto_approve` | 变宽 | 拦 `refused-confirm applied=0` | 拦（同） |
| `ask_high_risk → ask_every_step` | 变窄 | 不拦 `applied=1` | **不拦** `applied=1`（未变） |
| `auto_approve → ask_every_step` | 变窄 | 不拦 `applied=1` | **不拦** `applied=1`（未变） |
| `auto_approve → ask_high_risk` | 变窄 | 不拦 `applied=1` | **不拦** `applied=1`（未变） |

**总数核对**：3 对变宽 → 改前**只有 2 对有门**、改后 **3 对全有门**；3 对变窄 → 改前/改后**都 0 拦**（硬边界 1 未破）。逐行读数 `read-before.txt`/`read-after.txt` 的 `R1|PAIR|`。
**非 auto 变宽那一发**（`R1|NONAUTO|`）：改前 `set_err=<nil> read_after=ask_high_risk applied=1 last_result=applied last_detail=""`（**零确认、零拒绝、真落盘**）；改后 `set_err="perm: 把档位从 ask_every_step 放宽到 ask_high_risk 需要一条确认通道，但本机没有接入（fail-closed 未切）" applied=0`。★**一处逐字句变了，具名报**：`auto_approve` 的拒绝句从 `切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）` 改成方向句 `把档位从 %s 放宽到 %s 需要一条确认通道，但本机没有接入（fail-closed 未切）`。理由：`ask_high_risk` 的变宽**不欠一次问答**（R20/M4 只对 auto 出卡），照抄旧句会说出不成立的理由。旧句**无测试钉它**（`git grep '需要一次 L2 强确认\|没有接入确认通道' -- '*_test.go'` = 0 命中）；`367-a1-census.md ⑤` 与台账里那两处引用是**改前读数**，保留不动。

## ⑤ 两条产线路的实测（`AE → AH` 现在必须拒）

尺：`-overlay` 把 `D:\tmp\…\probe\cmd_probe_test.go` 挂进 `cmd/wisp` 包，经 `newComposerDispatchChain`（真 `config.Manager`＋真 `config.toml`＋真 `perm.Store`，两层 `Confirm` 都 nil）驱动**两个入口各一发** `panel.mode.request to=ask_high_risk`。`PATH` 必带 `third_party/sherpa-onnx`（不加就 `0xc0000135`＝假失败）。

| 入口（票面行号） | 走法 | 改前 | 改后 |
|---|---|---|---|
| `cmd/wisp/panel_host_windows.go:948` | `m.disp.Handle(ctx, raw)`（经 `dispatchRaw`） | 拒，`file_after="ask_every_step"` | 拒（同） |
| `cmd/wisp/panel_inbound.go:164` | `reply, err := disp.Handle(ctx, raw)`（经 `newPanelInboundDispatch`） | 拒，`file_after="ask_every_step"` | 拒（同） |

**逐字拒绝句**（两入口同一句，`handle_err` 与 `reply` 同源；文件仍 `ask_every_step`）：

```
panel: 确认腿未接入，档位不得变宽: 请求 "r1-host" 要把档位从 ask_every_step 放宽到 ask_high_risk，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）
```

★**诚实说明**：产线上这一发是**handler 层 `:133` 先拒**的，`Set` 的新门**没被走到**（两层 nil ⇒ 上层先拦）⇒ 这张表**不能**证明"新门在产线上会拦"；那是下表 `arm=half` 的读数。

## ⑥ 正控（尺双向）＋「有真 Confirm 就能过」那一发

手装三条 arm（产线装配把两处 nil 写死，⛔ 改 `cmd/wisp`）：同一请求 `AE → AH`，只换确认腿接没接。

| arm | handler 腿 | store 腿 | 改前 | 改后 |
|---|---|---|---|---|
| `none` | nil | nil | handler 拒、`applied=0` | handler 拒、`applied=0` |
| **`half`** | **接上** | **nil** | **`handle_err=<nil>`、`applied=1`（静默变宽）** | **`applied=0`**，句＝`perm: 把档位从 ask_every_step 放宽到 ask_high_risk 需要一条确认通道，但本机没有接入（fail-closed 未切）` |
| `both` | 接上 | 接上（真 `Confirm` 答 yes） | `applied=1` | **`applied=1`**（`MODE-SWITCH … result=applied`） |

★**`arm=half` 就是病根的直接读数**：改前只要 handler 有腿而 store 没腿，`AE→AH` **过 handler 门进 `Set` 并被 `Set` 放过**；改后同一装配**被 `Set` 自己拒**，且句是 store 的（`perm: …`），不是 handler 的。
**正控双向**（`R1|POSCTL|`，`to=auto_approve`）：`arm=nil` 拒／`arm=yes` `confirm_calls=1 set_err=<nil> applied=1`／`arm=no` `set_err="perm: 切到 auto_approve 被拒绝: operator clicked deny" applied=0` ⇒ 尺两侧都能读；**R20/M4 的不对称没动**（`R1|POSCTL|arm=yes-nonauto`）：`AE→AH` ＋ 真 `Confirm` ⇒ `confirm_calls=0`、`applied=1`。
★★**收益证明（具名）**：今天 `cmd/wisp/run.go:687` 那枚**带真 `Confirm`** 的 `ModeWriteHandler` 躺在 `rt.modeWrites` 里零读点。**哪天它被接到 dispatch 上**，`AE→AH` 会**过 handler 的 `:133` 门**（那层只看 `Confirm` 的 nil-ness，接上就放行）**进 `Set`** —— 改前 `Set` 对它无门（上表 `arm=half` 改前＝那一刻的读数：`applied=1`、零问答）；**现在 `Set` 也有门了**（`arm=half` 改后 `applied=0`）。⇒ 那枚现成入口从"一接上就漏"变成"接上也得先把确认腿递给 `Set`"。**本轮没有去接它**（另一件事、⛔ 本笔写面）。

## ⑦ 12 枚直调 `Set` 的测试：逐枚改前改后

尺＝`git grep '\.Set(' -- internal/perm/store_test.go internal/perm/ticket90_persist_test.go cmd/wisp/run_mode101_test.go` ⇒ **12 命中**，由 **9 个函数**覆盖。改前＝HEAD 的码＋HEAD 的测试（`-overlay` 换回 `git show HEAD:` 的 blob）。

| 调用点 | 那一发是什么 | 改前 | 改后 |
|---|---|---|---|
| `store_test.go:116` | `Mode(9)` 非法档 | 绿 | 绿 |
| `store_test.go:136` `:142` `:152` | 有 spy：`AE→AH`、`AH→AA`、`AH→AE` | 绿 | 绿（`spy.n()==1` 未变 ⇒ **M4 不对称没动**） |
| `store_test.go:185` | `AE→AA` 无通道 | 绿（拒） | 绿（拒） |
| `store_test.go:204` | `AE→AA` 答 no | 绿（拒） | 绿（拒） |
| `store_test.go:233` | `AE→AH`、**无 Confirm** | 绿（**应用**） | **不修就红** ⇒ 已改夹具（见下） |
| `ticket90_persist_test.go:154` | 有 Confirm，循环 AH/AA | 绿 | 绿（`confirmCalls==1` 未变） |
| `ticket90_persist_test.go:404` | 有 Confirm，`AH→AA` | 绿 | 绿 |
| `run_mode101_test.go:233` `:236` | 有 Confirm，AH 再 AA | 绿 | 绿（`asked==[auto_approve]` **未变**） |
| `run_mode101_test.go:546` | 真 L2 适配器，`AE→AA` 300ms 超时 | 绿（拒） | 绿（拒） |

**唯一变红的＝`TestTicket90PersistFailureKeepsMemory`**（`store_test.go:233`）。三读数：①HEAD 码＋HEAD 测试 **绿**；②**本笔码＋HEAD 测试 红**（`--- FAIL: TestTicket90PersistFailureKeepsMemory`，撞的是新门的 `ResultRefused` 而不是 `ResultPersistError`）；③本笔码＋改后夹具 **绿**。**判＝正确行为变化**：它测"写盘失败不得当成已切换"，`AE→AH` 只是夹具前提；按新不变量"没通道不得变宽"，夹具得先有一条通道才走得到写盘。**只加 `Confirm`，断言一字未松**（`ResultPersistError` 计数与"内存保持原值"两条都在）。**⛔ 为让它绿而放宽新门**——没有。

## ⑧ 新用例逐名 ＋ 为什么今天响

| 用例 | 在哪 | 把哪一跳摘掉它就红 |
|---|---|---|
| `TestSetGatesEveryWideningAndNoNarrowing` | `internal/perm/store_367_test.go` | 6 对一张表：3 变宽必拒／3 变窄必过。摘新门 ⇒ `AE→AH` 行红；变宽改成"任何变化" ⇒ 3 个变窄行红 |
| `TestSetWideningPassesOnceAChannelIsAttached` | 同上 | 通道接上后 `AE→AH` 必过且 `confirm_calls==0`（M4 不对称）；门改成"一律拒"就红 |
| `TestSetAsksAutoApproveEvenAfterAnUnreadableStoredMode` | 同上 | 钉"`from` 恒为三档之一"（`PermissionMode` fail-closed）⇒ 进 `auto_approve` 分支时不会 nil 解引用 |
| `TestIsWideningOnTheDocumentedLadder` | 同上 | 谓词本体的梯子钉子（perm 侧）；`risk.Mode` 一改序就红 |
| `TestPanelGateAndStoreGateAgreeOnEveryPair367` | `internal/panel/composer_gate_agreement_367_test.go` | **跨层一致性**：6 对两层必须同判。任一层单独改回旧判据就红 ⇒"两个转换点"的类级钉子 |
| `TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree`（改） | `internal/panel/composer_handlers_test.go` | 原有钉子改成驱动 `perm.IsWidening`：同一函数、同一组名字 |

## ⑨ 反向判据：三读数 ＋ 逐名报红行 ＋ 还原核

突变全在仓外 `-overlay`（两次突变后 `git status --porcelain internal/perm internal/panel` 仍只有本笔 7 个文件 ⇒ **工作树 0 突变**）。

**M1＝把新加的那一跳摘掉**（`:209` 回成旧判据 `if to == risk.ModeAutoApprove && s.confirm == nil {`）⇒ rc **1 / 1**（perm / panel）：
```
--- FAIL: TestSetGatesEveryWideningAndNoNarrowing
store_367_test.go:78: ask_every_step -> ask_high_risk was APPLIED with no confirmation channel; a widening must be refused (ticket 367)
store_367_test.go:82: … audit result = "applied", want "refused-confirm"
store_367_test.go:85: … moved the mode to ask_high_risk despite the refusal
store_367_test.go:88: … persisted anyway: mode="ask_high_risk" saves=1, want "ask_every_step"/0
store_367_test.go:98: … 0 refusal lines in the audit, want 1
--- FAIL: TestPanelGateAndStoreGateAgreeOnEveryPair367
composer_gate_agreement_367_test.go:76: ask_every_step -> ask_high_risk: the two gates disagree - perm.Store refused=false (<nil>), panel handler refused=true (…)
composer_gate_agreement_367_test.go:91: … is a widening and both layers must refuse it without a channel: store=<nil> handler=…
```
**M2＝把"变宽"定义改成"任何方向的变化"**（`:209` → `if from != to && …`）⇒ rc **1 / 1**，红的是**变窄**那三行：
```
--- FAIL: TestSetGatesEveryWideningAndNoNarrowing
store_367_test.go:104: ask_high_risk -> ask_every_step was REFUSED (…); getting stricter must never be blocked (R20)
store_367_test.go:107/110/113/116: … = "refused-confirm", want "applied" / did not apply / persisted saves=0, want 1 / wrote a refusal line for a narrowing
（`auto_approve → ask_every_step`、`auto_approve → ask_high_risk` 同五行，逐条各报）
--- FAIL: TestPanelGateAndStoreGateAgreeOnEveryPair367
composer_gate_agreement_367_test.go:76/97: … is a narrowing and neither layer may refuse it: store=… handler=<nil>
```
**还原核（时机＝最后一发突变之后量的）**：去掉 overlay 原样再跑 ⇒ `./internal/perm` **rc 0、PASS 19 / FAIL 0**；`./internal/panel` **rc 1、PASS 166 / FAIL 5**（既有红，见 ⑩）。三读数＝**变异红（M1/M2 各 1）／还原绿（perm 19/0）／原样绿（同上）**。

## ⑩ 四数（分开报、⛔ 相减）

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `go test -count=1 ./internal/perm` | **0** | PASS **19** / FAIL **0** |
| 2 | `go test -count=1 ./internal/panel` | **1** | PASS **166** / FAIL **5** ⇒ **5 枚全是既有红**，与普查**同集**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestC21DesignTokensFourWayAgree`、`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`。**本笔新增红 0 枚**（判据：`--- FAIL` 名单逐名等于普查那 5 枚） |
| 3 | `PATH="…/third_party/sherpa-onnx;$PATH" go test -count=1 -run 'TestTicket101\|TestTicket223ModeLoosening…\|TestReplyListenerReject…' ./cmd/wisp` | **0** | PASS **7** / FAIL **0**（含 `TestTicket101ManualSwitchSurvivesRestart`：`asked==[auto_approve]` **未变**） |
| 4 | `GOFLAGS= go build ./...` ＋ `go vet ./internal/perm ./internal/panel` | **0 / 0** | — |

★`cmd/wisp` **整包不跑**（会死在 `MARK299`），只跑上面定向 `-run`；★不带 sherpa `PATH` 是 `0xc0000135` 假失败，本表全部带。
## ⑪ 门禁 rc ＋ 格式甲乙两尺

| 门 | rc | 读数 |
|---|---|---|
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；ban #8 扫 `internal/` 639 ＋ `cmd/` 149 个 Go 文件（含注释与 `_test.go`） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；over-budget 61 ＝ roster 61、`not in roster=0` |
| 格式**甲（入库 blob＝`git show :path`｜`gofumpt.exe -l`，判据）** | — | 本笔 7 个 `.go` **逐个**：全 **CLEAN**；尺本身双向验过（未格式化的 stdin 被点名、已格式化的静默） |
| 格式**乙（工作树，只作诊断）** | — | 同上 7 个文件 `gofumpt -l` 与 `gofmt -l` **都空** |

## ⑫ 提交

- 逐字点名的文件（⛔ glob）：`internal/perm/direction.go`、`internal/perm/store.go`、`internal/perm/store_test.go`、`internal/perm/store_367_test.go`、`internal/panel/composer_handlers.go`、`internal/panel/composer_handlers_test.go`、`internal/panel/composer_gate_agreement_367_test.go`、本文件。
- `git diff --numstat`（改动的 4 个）：`store.go 39/14`、`composer_handlers.go 20/20`、`composer_handlers_test.go 11/7`、`store_test.go 9/1`。
- **只 commit、⛔ push**；⛔ 翻框、⛔ 改票名/`Status:`/`-done`、⛔ 改台账、⛔ 碰票 114 原句。

## 附：落笔前五问（尺的射程）

量的是**运行行为**（真 `config.Manager`＋真 `config.toml`／真 `perm.Store`），⛔ 注释、⛔ 票面；输出一行一读数、`|` 分隔，错误一律 `<nil>` 或**原文带引号**（`%q` 打在 `error` 上出 `%!q(<nil>)`＝尺的产物，已改成先 render 成字符串）；数的是**代码引用**（12 枚测试调用点＝7＋2＋3；`Set` 的真实现仍只有 3 个：`perm.Store` ＋ 2 个测试假腿）；**数和总数对上了**（变宽/变窄 3/3；panel 5 枚既有红逐名同集；perm 19 枚 PASS 由 `-v` 数 `^--- PASS`）；**谁算的**＝`Select-String '^--- (PASS|FAIL)' | Measure-Object`、`git diff --numstat`、`wc -l`（⛔ `grep '^[+-][^+-]'` 会漏"加号后为空"；`Measure-Object -Line` 对本文件少报，见普查第 5 问那个坑）。
