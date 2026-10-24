# 372-r2 交付件：回退到「只有目标档 `auto_approve` 要卡」＋ 读侧降级（有效档）

- 腿 `372-r2`（写腿＝实现者），量于 2026-10-24 10:5x–11:4x +08；开工时 HEAD `3efb53526`（`114-v2` 终裁），交件时 HEAD `255d42c75`（期间编排者／`368-r1` 只动 `docs/**`、`.scratch/**`、`internal/memory/**`；`git diff --name-only 3efb53526..HEAD` ⇒ **本笔 4 个包零他笔改动**；四数在 11:26 之后复跑，`internal/memory` 的改动已于 11:19 落盘）。
- 落盘：`internal/perm/**`（`store.go`／`direction.go` 改、`store_372_test.go` 新、`store_367_test.go`／`ticket90_persist_test.go` 必要测试）· `internal/panel/**`（`composer_handlers.go` 那道门＋3 枚必要测试）· `cmd/wisp/panel_mode_confirm_114_windows_test.go`（那 2 枚期望）· `cmd/wisp/panel_inbound.go:341` 一句注释 · 本件。**代码面 10 文件（9 改 1 新），加本件共 11 个 pathspec**。
- 纪律：⛔ push／⛔ 翻框／⛔ 改票名・`Status:`・`-done`／⛔ 台账；⛔ 碰 `internal/tools`、`cmd/wisp` 其他文件、`.scratch/**` 与别人的脏。
- ★开工前索引读数：`git diff --cached --name-status` ＝ `A .scratch/wisp/issues/114-…-done.md` ＋ `M .scratch/wisp/issues/114-….md`（**别人的、全在 `.scratch`**；`status` 里还有 `AD`／`M `）⇒ 按派单「工作树里有别人的脏」处理：**一处未碰、未清**；本笔 `add`／`commit` 全用显式 pathspec，这两项不会进我的提交。

## ① 判决依据（唯一权威）
- `m4-survey.md:96`（一句总判）：确认挂在**目标档位**与**这一发动作**上、⛔ 挂在「变宽」；**6 个有实现的 repo 没有一个是「因缺通道拒绝放宽档位」**。
- `m4-survey.md:104-106`：倾向＝**允许改档 ＋ 动作层 fail-closed**；`:106` 把 Step-Code `effectiveState`（`auto→confirm`，`permissions.ts:518-521`）列为**具名推荐**；openchamber `displayedPermissionMode`（`:88`）同形；`:105`「中间档不需要确认」有 dsh／minimax／openchamber 三票。
- `372-r1-blocked.md §③`：回退会让那 2 枚 `cmd/wisp` 测试转红（红行已具名）；`§④ 末段`：`perm.Store` 是**全进程唯一**同时知道「配的哪档」＋「有没有接确认腿」的对象、且已是 `tools.ModeSource`；读侧降级是本仓最省且同样主流的形；⚠ `Set` 的 from/to 必须改读配置值。**本腿接着做这两条，未重做它的量。**

## ② 件一：门回退到目标档（行号为改后）
- `internal/perm/direction.go:23` `+RequiresCard(to) = to == risk.ModeAutoApprove`（**唯一判据定义**）；`:36` `IsWidening` 保留，doc 改写成「**阶梯序见证、不再作门**，谁都不得拿它拒切换」⇒ `367` 做对的「一处定义」那半保留。
- `internal/perm/store.go:254` `from := s.configuredMode()`；`:268` `+if RequiresCard(to) && s.confirm == nil` ⇒ 只拦**目标 `auto`**（错误句点名 from/to）；`:281` `if RequiresCard(to)` 才走 M4 那张卡（仍无 nil 测试，靠 `:268` 保证非 nil）。
- `internal/panel/composer_handlers.go:153` 第二处调用 `+if perm.RequiresCard(to) && h.Confirm == nil`（同一函数，⛔ 各写一份）；`:140` `+if from == to` 短路交回 `Set`（no-change 不是要卡请求；也保住「配置 AA＋无腿、请求 AE」能落盘）。
- 效果：`AE→AA`／`AH→AA` 仍拦；`AE→AH` 不拦（＝R20/M4＋主流）。

## ③ 件二：读侧降级（有效档）——落在 `perm.Store.PermissionMode()`
- `internal/perm/store.go:196-205`：`configured := s.configuredMode(); if s.confirm == nil && configured != ask_every_step { return ask_every_step }`。
- `:215-224` `+configuredMode()` 是 config.toml **原值**，**不写配置**（两发对照里 `saves=0`、`config_file` 仍是原档）。
- 为什么在这层：`internal/tools/mode.go:27-31` 的 `ModeSource` **每发读一次**；`perm.Store` 全进程唯一同时知道「配的哪档」＋「有没有接确认腿」；形状＝Step-Code `effectiveState`（m4-survey:106）／openchamber `displayedPermissionMode`（:88）。
- 审计：`store.go:164-175` 在降级发生时落一枚 `MODE-DOWNGRADE from=<配置> to=ask_every_step result=downgraded-no-channel` 记录＋同字段日志行（**不是 MODE-SWITCH**：配置文件没动）；`auditAs` 在 `:363`。
- 记账：`Set` 的 `from` 读配置值（见 ②）——否则「配置 AA＋无腿＋请求 AE」会被当 no-change 而不落盘（r1 §④ ⚠ 已解，M3 牙口见 ⑪）。

## ④ 六对表（改后，`Confirm=nil`，仓外 probe rc=0）
| from→to | 目标 auto | set_err | 落盘(saves) | config 读回 | 有效档读回 |
|---|---|---|---|---|---|
| AE→AH | 0 | `<nil>` | applied(1) | `ask_high_risk` | `ask_every_step` |
| AE→AA | 1 | 拒（点名 AE/AA） | refused(0) | `ask_every_step` | `ask_every_step` |
| AH→AE | 0 | `<nil>` | applied(1) | `ask_every_step` | `ask_every_step` |
| AH→AA | 1 | 拒（点名 AH/AA） | refused(0) | `ask_high_risk` | `ask_every_step` |
| AA→AE | 0 | `<nil>` | applied(1) | `ask_every_step` | `ask_every_step` |
| AA→AH | 0 | `<nil>` | applied(1) | `ask_high_risk` | `ask_every_step` |
⇒ **3 变宽里 2 拦 1 过**（只有目标 `auto` 拦）；**3 变窄全过**。正控：`AE→AH`＋真 `Confirm` ⇒ `applied=1 confirm_calls=0`；`AH→AA`＋真 `Confirm` ⇒ `applied=1 confirm_calls=1`。

## ⑤ 五发 A–E（我自己重跑；`ModeSource` 用**真 `perm.Store`**，probe rc=0）
| 发 | config / 腿 | gate | 有效档 | is_error | 落盘 | L1 窗 | L2 卡 | MODE-SILENCE |
|---|---|---|---|---|---|---|---|---|
| A | `auto_approve` / none | NoGate | `ask_every_step` | **1** | **0** | 1 | 0 | 0 |
| B | `ask_every_step` / none | NoGate | `ask_every_step` | 1 | 0 | 1 | 0 | 0 |
| C | `ask_high_risk` / none | NoGate | `ask_every_step` | **1** | **0** | 1 | 0 | 0 |
| D | `auto_approve` / none，越授权目录 | NoGate | `ask_every_step` | 1 | 0 | 0 | 1 | 0 |
| E | `ask_every_step` / wired | 允许 | `ask_every_step` | 0 | 1 | 1 | 0 | 0 |
⇒ **A 与 C 从「静默落盘」变「拒」**（`file_written=0`，文本「L1 确认窗口尚未接入（票 21），已拒绝执行」）；B／D 仍拒；E 仍「问一发过一发」。
- **M1 对照**（摘掉读侧那一跳）：A `effective=auto_approve`、`is_error=0`、`file_written=1`、`MODE-SILENCE=1`；C 同；B／D 不动 ⇒ **A/C 两格就是件二的牙口**。
- **F 对照**（config `auto_approve`＋**非 nil** 腿）：`effective=auto_approve`、落盘、MODE-SILENCE=1 ⇒ 降级只认「腿的 nil-ness」，⛔ 按档位值乱降。

## ⑥ CLI 管那一发（`newPanelInboundDispatch` 生产装配，probe rc=0）
```
to=ask_high_risk  handle_err=<nil>                   disk_after=ask_high_risk  effective_after=ask_every_step
to=auto_approve   handle_err=ErrNoL2Confirm(面板句)   disk_after=ask_high_risk  effective_after=ask_every_step
```
⇒ `AE→AH` **确实把 disk 改了**（`ask_high_risk`）——这是 `114-v2` §4 点名的**已知代价**（管子无卡片表面、M4 不卡中间档），本腿按派单「明知地改期望」；但同一进程**有效档仍是 `ask_every_step`**（件二）⇒ 不构成静默放宽。`AA` 仍拒、disk 不动。

## ⑦ `PermissionMode()` 两发对照（语义变了的读数）
| config.toml | 腿 | `PermissionMode()` | 配置被改写？ | MODE-DOWNGRADE |
|---|---|---|---|---|
| `ask_high_risk` | none | **`ask_every_step`** | 否（saves=0） | 1 行 |
| `auto_approve` | none | **`ask_every_step`** | 否（saves=0） | 1 行 |
| `ask_high_risk` | wired | `ask_high_risk` | 否 | 0 |
| `auto_approve` | wired | `auto_approve` | 否 | 0 |

## ⑧ panic 那发（`recover()` 在侧）
- from `ask_every_step`、`confirm=nil`、`Set(auto_approve)`：`panicked=0`、`set_err="perm: 从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）"`、`config_after=ask_every_step`、`saves=0`。
- from `ask_high_risk` 同值 ⇒ **拒、⛔ panic**；`store.go` 那条 `There is no nil test here on purpose` 仍承重（门在它上面把 `AA`＋nil 拦掉）。M4（删掉整段目标门）⇒ `TestAutoApproveWithNoChannelIsRefusedNotPanicked` 真 panic（见 ⑪）⇒ 这条注释有牙。

## ⑨ 那 2 枚 `cmd/wisp` 测试改期望：三步法
1. **旧码＋旧测＝绿**：HEAD 起点跑定向 `-run 'TestTicket114r2|…'` ⇒ `rc=0`（`D:\tmp\…\base-cmdwisp.txt`）。
2. **本笔码＋旧测＝红**（旧测文件经 `-overlay` 挂回）：`rc=1`，逐行：
   `panel_mode_confirm_114_windows_test.go:106`（无腿必须拒，got `<nil>`）／`:109`（拒绝后 disk 却成 `ask_high_risk`）／`:160`（管子必须拒，got `<nil>`）／`:163`（disk 成 `ask_high_risk`）。
3. **本笔码＋新期望＝绿**：最终定向 `rc=0`（两枚在内）。
- **断言未松**：`AA` 那半原样保留（`errors.Is(ErrNoL2Confirm)`、`w.setCalls==0`、`audit` 恰 1 行、disk 不动）；`AH` 那半只把「拒＋disk 不动」改成「过＋disk 落 AH」，**另加**一枚有效档读数（⑥/⑦），⛔ `Skip`／⛔ 改宽错误串。

## ⑩ 新用例逐名 ＋ 为什么今天响（红＝摘掉那一跳就红）
| 用例 | 钉哪一跳 | 今天响（变异/读数） |
|---|---|---|
| `perm.TestSetGatesOnlyTheAutoApproveTargetAndEveryNarrowingApplies` | 件一·store 门 | M2 红 `store_367_test.go:112/115/118/122` |
| `perm.TestPermissionModeDowngradesWithoutAConfirmChannel` | 件二·读侧 | M1 红 `store_372_test.go:36`（AH/AA 读回原档） |
| `perm.TestPermissionModeDowngradeIsRecordedOnce` | 件二·审计 | M1 红 `store_372_test.go:63`（MODE-DOWNGRADE 行=0） |
| `perm.TestSetAccountsOnTheConfiguredValueNotTheDowngradedRead` | 件二·记账 | M3 红 `store_372_test.go:113/116/119` |
| `perm.TestAutoApproveWithNoChannelIsRefusedNotPanicked` | 件一·nil 安全 | M4 红 `store_372_test.go:134`（真 panic） |
| `panel.TestOnlyTheAutoApproveTargetIsRefusedWhenNoL2ConfirmLegIsAttached` | 件一·第二处调用 | M2 红 `composer_handlers_test.go:180` |
| `panel.TestPanelGateAndStoreGateAgreeOnEveryPair367` | 两处同判据 | M2 红 `composer_gate_agreement_367_test.go:90/94/124` |
| `cmd` 那 2 枚 | 端到端 | M2 红 `:119/122`、`:193/196`；旧测见 ⑨② |

## ⑪ 反向判据：4 发变异（全仓外 `-overlay`，工作树零突变）＋ 还原核
| 发 | 摘掉/翻回哪一跳 | rc | 逐名红行 |
|---|---|---|---|
| M1 | `PermissionMode` 去掉降级 | perm 1 | `store_367_test.go:128`×3；`store_372_test.go:36`×2、`:63`、`:103` |
| M2 | 两道门翻回 `IsWidening` | perm/panel/cmd 1 | `store_367_test.go:112/115/118/122`；`composer_gate_agreement_367_test.go:90/94/124`；`composer_handlers_test.go:180`；`cmd…:119/122`、`:193/196` |
| M3 | `from := s.PermissionMode()` | perm 1 | `store_367_test.go:103/115/118`；`store_372_test.go:113/116/119` |
| M4 | 删掉整段目标门 | perm 1 | `store_372_test.go:134`（nil pointer panic） |
- **还原核（在最后一发 M4 之后）**：真实工作树（无 overlay）四数见 ⑫；`git status --porcelain internal/perm internal/panel cmd/wisp` 仍只有本笔 10 个文件。

## ⑫ 四数与门禁（既有红分开报，⛔ 相减）
| 包 | 基线（HEAD） | 改后 | 备注 |
|---|---|---|---|
| `./internal/perm` | rc0 PASS 19／FAIL 0 | **rc0 PASS 23／FAIL 0** | +4 枚新用例 |
| `./internal/panel` | rc1 PASS 166／FAIL 5 | **rc1 PASS 166／FAIL 5** | 既有 5 枚**逐名**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`（本笔**零新增**） |
| `./internal/tools` | rc0（未计数） | **rc0 PASS 250／FAIL 0** | 同一把尺（`-v`，含子测试） |
| `./cmd/wisp` 定向（PATH 带 `third_party/sherpa-onnx`） | rc0 | **rc0 PASS 29／FAIL 0** | `-run 'TestTicket114r2|TestAC9|TestTicket101|TestTicket223|TestTicket90|TestPanelInbound|TestPanelMode|TestReplyListenerReject|TestTicket197|TestSubagent'`；整包会死在 `MARK299`，只跑定向 |
- 门禁：`go build ./...` rc0 · `go vet ./internal/perm ./internal/panel ./internal/tools ./cmd/wisp` rc0 · `tools/d22scan/d22scan.exe` rc0（clean）· `gofumpt -l` 对 10 个改动文件**空**、`gofumpt -d` **0 行**（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）。
- 尺：PASS/FAIL 数是 `-v` 下 `^--- PASS:`／`^--- FAIL:` 计数（**含子测试**）；`internal/tools` 基线未取（不在派单要求内），只报改后一列，⛔ 当"零变化"读。

## ⑬ 具名残余（⛔ 自行掩掉）
1. **`wisp run` 这条装配不会触发读侧降级**：`run.go:668-671` 在 `confirm==nil` 时替换成 `rt.confirmModeSwitch`（永远非 nil）⇒ `perm.Store` 认为「有腿」，`PermissionMode()` 不降级。⑤ 的 F 行就是这形状：非 nil 腿＋AA ⇒ L1 仍静默。⇒ 本腿按派单**未动 `cmd/wisp` 其他文件**；要堵这格需让 run seam 在「本机没有可答表面」时传 nil（或让 Store 收一个「可问性」而不是 nil-ness）——**下一票的事**，本件只把读数留在这里。
2. **CLI 管把 `ask_high_risk` 落盘**（⑥）：按 114-v2 §4 的「明知地改期望」处理；有效档不升（件二）是这条的兜底。
3. `ticket90_persist_test.go` 的 `askOnce` 夹具补了确认腿（否则「改配置键 ⇒ 链路行为变」在 372-r2 后不可能成立）；测试断言语义未改。
