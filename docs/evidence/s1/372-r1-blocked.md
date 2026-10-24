# 372-r1 **停手上报**（⛔ 这不是交件：本腿**一行代码都没落盘**）

- 腿 `372-r1`（写腿）。HEAD `e66c8eb65`（`114-r3` 于 10:3x 刚落的那枚），量于 2026-10-24 10:2x–10:5x +08。
- **落盘 0 字节**：未改任何 `.go`、未 commit、未 push、未翻框、未改票名/`Status:`、未碰台账。
- 停手原因＝派单里**具名的两条**停手条件**同时命中**：③（要动 `cmd/wisp`）与「动作层若在 `internal/tools` ⇒ 停手上报」。
- 本文件只放**读数与具名代价**，供编排者串行安排；**所有探针/候选件都在仓外** `D:\tmp\wisp\20261024\372r1\**`，工作树 0 突变。

## ① 开工前核对（逐条现跑）

| 尺 | 读数 |
|---|---|
| `git diff --cached --name-status` | **空**（无 `D`/`A`） |
| `git status --porcelain internal/perm internal/risk` | **空** |
| `git status --porcelain internal/panel` | 10:30 那一跑是 ` M internal/panel/composer_handlers.go` ⇒ **本应停**；**10:35 复跑已空**，因 `114-r3` 在 10:3x 提交了 `e66c8eb65`（`114-r3` 的四处措辞改准，零行为变更）⇒ **该条件已自行解除** |
| `git status --porcelain`（全树） | 别人的脏仍在（`.scratch/**` 的 `M`、`?? $null`/`-`/`-dummy` 等）⇒ **一处未碰** |

★所以**开工前的停手条件（③脏文件）不成立**了；**真正挡住本腿的是下面的 ② 和 ③**。

## ② AC#0：我自己复跑的核心读数（⛔ 抄 `367-r1-impl.md`）

尺＝同一支仓外探针经 `-overlay` 挂进 `internal/perm`（`perm_probe_test.go`，**只调行为、不点名任何谓词**，故 367 码与回退码都能编译）；`before`＝**当前树**（＝367 码），`after`＝仓外候选 `cand/store.go`（**只把门换成 `RequiresCard(to)`＋还原拒绝句**）。两跑 rc 均 0。

`R1|LADDER|names=[ask_every_step ask_high_risk auto_approve]|names_len=3|default=ask_every_step`

### 六对「变宽→门」表（`Confirm` 全 nil）

| 从 → 到 | 方向 | 改前（367，今天） | 改后（候选＝回退） |
|---|---|---|---|
| `AE → AH` | 变宽 | **拦** `applied=0 refused-confirm` | ★**不拦** `set_err=<nil> applied=1 file_after="ask_high_risk" saves=1 refusal_lines=0` |
| `AE → AA` | 变宽 | 拦 | **拦**（`切到 auto_approve 需要一次 L2 强确认…`） |
| `AH → AA` | 变宽 | 拦 | **拦**（同） |
| `AH → AE` | 变窄 | 不拦 `applied=1` | 不拦 `applied=1` |
| `AA → AE` | 变窄 | 不拦 | 不拦 |
| `AA → AH` | 变窄 | 不拦 | 不拦 |

**总数核对**：3 变宽 → 改前 **3 拦**、改后 **2 拦 1 过**；3 变窄 → 两跑**都 0 拦**（硬边界 1「⛔ 改机主 2026-09-21 定案」未破：只有目标 `auto_approve` 要确认）。
★**`AE→AH` 这一行就是本票的全部收益**，也是主流那句「允许改档」。

### 其余三发（两跑同值）

| 发 | 读数（改前＝改后） |
|---|---|
| `AE→AH` ＋真 `Confirm` | `set_err=<nil> applied=1 confirm_calls=0` ⇒ **M4 不对称未动** |
| `AH→AA` ＋真 `Confirm` | `applied=1 confirm_calls=1`（正控：该出卡的那档仍出卡） |
| **panic 那发**（`AA`＋`confirm==nil`，`recover()` 在侧） | `panicked=0 set_err="perm: 切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）" applied=0 mode_after=ask_high_risk` ⇒ **拒、⛔ panic**；`store.go` 那条 `There is no nil test here on purpose` **仍承重**（门在它上面把 `AA`+nil 拦掉） |

## ③ 件一为什么落不了：**回退会让 2 枚 `cmd/wisp` 测试转红**（该包＝本腿禁区）

尺＝仓外候选经 `-overlay` 换掉 `internal/perm/store.go` ＋ `internal/panel/composer_handlers.go`（**测试全用树里的 HEAD 版**），`PATH` 带 `third_party/sherpa-onnx`，`-run "TestTicket114r2|TestTicket101|TestTicket223|TestReplyListenerReject|TestPanelInbound|TestPanelMode"`：

| | rc | PASS / FAIL |
|---|---|---|
| 基线（不动） | **0** | 21 / 0 |
| 候选（回退） | **1** | 19 / **2** |

红的两枚**逐名＋逐行**（`panel_mode_confirm_114_windows_test.go`）：

```
--- FAIL: TestTicket114r2ResidentPanelChainFromRefusalToTheWriteLeg
 :106  with no leg attached the resident door must refuse with ErrNoL2Confirm; got <nil>
 :109  a refused widening moved the档 on disk: got "ask_high_risk", want ask_every_step
--- FAIL: TestTicket114r2CLIPipeStillRefusesEveryWidening
 :160  the pipe has no card surface, so a widening to ask_high_risk must be refused with ErrNoL2Confirm; got <nil>
 :163  the pipe widened the档 to "ask_high_risk" with no card anywhere
```

⇒ 这两枚**逐字钉的是 367 的判据**（`AE→AH` 无腿必须拒），而 `114-r2` 的「读数 1」与「读数 3」**正是它们**。回退＝**把 `114-r2` 刚交的那两格读数作废**，而这**必须在 `cmd/wisp/**` 里改测试** ⇒ **本腿写面外**（硬边界 3）。
★**另有一处过期注释**（同样是禁区）：`cmd/wisp/panel_inbound.go:341` 逐字 `gate is \`IsWidening(from, to) && s.confirm == nil\` (store.go:209)`，回退后即不成立。**一处未碰。**

### 本腿写面内会红的三枚（这三枚是**我的**、可改，已备好改法）

| 包 | rc（候选） | 红的 | 归属 |
|---|---|---|---|
| `./internal/perm` | 1 | `TestSetGatesEveryWideningAndNoNarrowing` | 我的（改成「只有 `AA` 目标被拒」表＋`AE→AH` 点名行） |
| `./internal/panel` | 1 | `TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached`、`TestPanelGateAndStoreGateAgreeOnEveryPair367` | 我的（同上，改 9 对全表） |

**既有红分开报、⛔ 相减**：基线 `./internal/perm` **rc 0 / PASS 19 / FAIL 0**；`./internal/panel` **rc 1 / PASS 166 / FAIL 5**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`、`TestComposerContractTypesMatchFrontend`、`TestC21DesignTokensFourWayAgree`）。候选跑面板 **PASS 164 / FAIL 7**＝既有 5 ＋ 我的 2，**本笔未新增任何「不可归我」的红**。

### 12 枚直调 `Set` 的测试：**一枚都不用改**（逐枚量过）

尺＝`-run` 点名覆盖它们的那 9 个函数，两跑（`overlay-before`＝367 码 / `overlay-cand`＝回退码）**逐名并排**：

| 函数（覆盖的调用点） | 改前 | 改后 |
|---|---|---|
| `TestTicket90ThreeModesAndNoFourth`（`store_test.go:116`） | PASS | PASS |
| `TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited`（`:136 :142 :152`） | PASS | PASS（`spy.n()` 断言未动） |
| `TestTicket90NoConfirmChannelMeansNoAuto`（`:185`） | PASS | PASS |
| `TestTicket90ConfirmRefusalLeavesModeAlone`（`:204`） | PASS | PASS |
| `TestTicket90PersistFailureKeepsMemory`（`:233`） | PASS | PASS |
| `TestTicket90ManualSwitchSurvivesRestart`（`ticket90_persist_test.go:154`） | PASS | PASS |
| `TestTicket90ConfigKeyChangesWhatTheChainAsks`（`:404`） | PASS | PASS |
| `TestTicket90HandEditLooseningGoesThroughD36` | PASS | PASS |
| `TestTicket90SnapshotIsReadOnly` | PASS | PASS |
| `TestTicket101ManualSwitchSurvivesRestart`（`run_mode101_test.go:233 :236 :546`，cmd 包） | PASS | PASS（在③那次 21 发的跑里） |

★**367-r1 改过的那枚夹具（`store_test.go:227`，给 `TestTicket90PersistFailureKeepsMemory` 加的 `Confirm`）现在应当**恢复**（摘掉）**，理由不是变绿（它**本来就是绿**），而是两件**变强**的事**：①它是 367 为「任何变宽都拒」加的，判据回退后**理由消失**；②**留着它反倒丢牙**——有通道在场时，**若有人把方向门悄悄加回来，这枚测试依然是绿的**；摘掉之后同一发会红在门那一跳。**断言一字不动**（`ResultPersistError` 计数＋「内存保持原值」两条都在）。

## ④ 件二为什么落不了：**动作层存在，但它不在我的写面里**

### 「动作层」在哪（现跑找的，具名到行）

| 环节 | 位置 | 它今天怎么 fail-closed |
|---|---|---|
| 每次调用的判级→路由 | `internal/tools/bridge.go:310-314`（读档＋`mode.Screen`）→ `:318-339`（`MODE-SILENCE` / `MODE-REDLINE`）→ `:364-380` → `route` `:444-499` | 只在**真去问**的时候 |
| 问的两条腿 | `internal/tools/gate.go:126-140`（`GateFuncs`：回调 nil ⇒ `AnswerReject`）、`:146-156`（`NoGate` ⇒ reject） | 是，但**得先走到它** |
| 档位把问题**抽掉** | `internal/risk/mode.go:192-212`（`ask_high_risk` 静默 L1；`auto_approve` 静默非红线 L2） | ⛔ 不判通道 |
| 产线装配 | `cmd/wisp/run.go:775`（`Gate: rt.gate`）、`cmd/wisp/panel_inbound.go:272-291`（`perm.New{Confirm: storeConfirmFromPanel(confirm)}` ＋ 同一枚 `confirm` 给 handler） | 两处 nil 写死 |

⇒ ★**缺口精确到一行**：`bridge.go:311 sil := mode.Screen(verdict)` 之后，**只要 `sil.Silenced` 为真，`route` 就永远走 `case risk.L0`，`gate` 一次都不被调** ⇒ 「档位要求确认而本进程没有通道」这件事**没有任何地方读得到**。

### 五发实测（仓外探针 `-overlay` 挂进 `internal/tools`，业务调用同一发 `fs.write`）

```
A auto_approve  channel=none   in-scope  is_error=0 file_written=1 L1_windows=0 L2_approvals=0 MODE-SILENCE=1
B ask_every_step channel=none  in-scope  is_error=1 file_written=0 L1_windows=1 L2_approvals=0 MODE-SILENCE=0  text="L1 确认窗口尚未接入（票 21），已拒绝执行"
C ask_high_risk  channel=none  in-scope  is_error=0 file_written=1 L1_windows=0 L2_approvals=0 MODE-SILENCE=1
D auto_approve   channel=none  outside   is_error=1 file_written=0 L1_windows=0 L2_approvals=1            text="L2 审批通道尚未接入（票 21），已拒绝执行"
E ask_every_step channel=wired in-scope  is_error=0 file_written=1 L1_windows=1 L2_approvals=0
```

- ★**A＝件二那条「改前：静默通过」的读数**：档位 `auto_approve`、**零通道**、可逆写**真落盘**、**零张卡**、只留一行 `MODE-SILENCE`。
- ★**C 同病**，且在**中间档**上 ⇒ 回退后 `AE→AH` 一旦被允许，**中间档也把 L1 静默掉**（而 `367` 的拒恰好遮住了这一格）。
- **B/D 说明今天的 fail-closed 是真的**，但只在**问得出来**的时候（D 的红线仍拒）。
- **E＝正控**：有通道 ⇒ 问一发、过一发。

### ★可达性（现跑核过，这一条**改变结论**）

「放宽档位 ＋ 无通道」这个组合**今天一个门都进不去**：

| 入口 | 今天怎么被挡 | 出处 |
|---|---|---|
| `Set` | 367 的方向门：任何变宽无腿 ⇒ 拒 | `internal/perm/store.go:209` |
| 手改 `config.toml` 放宽 `risk.permission_mode` | `[risk]` 是 **locked** 段 ⇒ 走 D36 的 `ConfirmLocked`，**hook 为 nil 即拒** | `internal/config/manager.go:248-251`、`:314-343`、`:431-460`、**:180 `approved := confirm != nil && confirm(...)`** |

⇒ ★★**件一自己会开出本仓第一条「放宽了却没有通道」的路**：`Set` 放行 `AE→AH` 之后，上面探针那格 **C**（`ask_high_risk` 静默 L1）就从「不可达」变成「面板/CLI 一点就到」。**所以件一 ⛔ 单独落**——这正是派单说「两件都要做、⛔ 只做一件」的硬理由，不是排期偏好。

### 要补这一跳，必须动哪些包（具名）

1. **`internal/tools/**`**（判据落点：`bridge.go:311` 之后、「`sil.Silenced` 生效之前」加一跳：**通道不可问 ⇒ 就按没被静默处理**，落到 `route` 的 L1/L2 分支由 `gate` 拒）。
2. **`cmd/wisp/**`**（唯一知道「本进程有没有一条真能出卡的通道」的装配点：`run.go:775` / `panel_inbound.go:272-291`）。
3. **⛔ 不能给 `tools.Gate` 加方法**：`cmd/wisp/subagent_selfapproval_197_test.go:575-579` 把该方法集**逐字钉成 2 枚**（`PendingWindow`/`PendingApproval`，197「子代理永不自批」的护栏）⇒ 形只能走 **`tools.Options` 的新字段**（如 `Askable func() bool`，197 的禁名表 `Allow/Native/DecideFromNative/DecideFromPanel/GrantNonce` 不含它）。
4. 代价具名：该字段**产线 nil ⇒ 一律按「问不出来」** ⇒ 在 `cmd/wisp` 接上真答案之前，`auto_approve` 与 `ask_high_risk` 会比今天**更严**（每一发本来会被静默的调用直接拒）。**这正是 fail-closed 的方向**，但**必须先接上再声称「全自动能用」**。

### ★另一条**在写面内**的等价形（供编排者拍，我未擅自选）

`perm.Store` 是**全进程唯一同时知道「配的哪档」和「有没有接确认腿」**的对象，且它**已经是 `tools.ModeSource`**（`internal/tools/mode.go:27-31` 每发读一次）⇒ 把「无通道 ⇒ 有效档降到 `ask_every_step`」做在**读侧**（Step-Code `permissions.ts:518-521` 的 `effectiveState` 形状，m4-survey.md:106 具名推荐），**一行判据＋一枚 `MODE-DOWNGRADE` 审计行**，就能让 A/C 两格变成「问 ⇒ `NoGate` 拒 ⇒ 拒」，**完全落在 `internal/perm` 内、⛔ 动 `internal/tools`/`cmd/wisp`**。
- ⚠ 代价：`Set` 的 `from/to` 记账必须改读**配置值**（否则「配置 `AA`、无腿、请求 `AE`」会被当成 `no-change` 而不落盘）；`PermissionMode()` 的语义变成**有效档**（openchamber 的 `displayedPermissionMode` 就是这么做的，m4-survey.md:88）。
- ★**为什么这条路真的通到产线**（现跑核过）：产线桥的 `tools.Options.Modes` 就是 `rt.modes`，而 `rt.modes` 的类型是 **`*perm.Store`**（`cmd/wisp/run.go:780` ＋ `:330` ＋ `:682`）⇒ 降到读侧**不需要动 `cmd/wisp` 一个字节**，也不需要动 `internal/tools`。
- ⚠ 但它**不消掉③那两枚 `cmd/wisp` 红**：`newComposerDispatchChain` 里 CLI 那条路的 `perm.Store` 是 `Confirm: storeConfirmFromPanel(nil)`＝nil，`Set` 仍会按目标档放行 `AE→AH` ⇒ 114-r2 的读数 1/3 仍是红的。**两件事互相独立。**
- 我的判断：**它是本仓最省、且同样主流的形**；但派单把件二写成「动作层」，**形由编排者拍**，我不擅自换层。

## ⑤ 我需要什么（三选一，★我推荐甲）

- **(甲) 放 `cmd/wisp`＋`internal/tools`**（**推荐**）：`372-r1` 一口气把件一（回退）＋件二（调用层 fail-closed）**同波**落完，并改那 2 枚 114-r2 测试。★只有这一形能**同时**满足「件一＋件二」和「⛔ 单独落件一」。⚠ 那两枚是 `114-v2` 正在裁的读数 ⇒ **建议先有 114-v2 的裁决**。
- **(乙) 同波派一条 `cmd/wisp` 腿**（只改 `panel_mode_confirm_114_windows_test.go` 的两处期望＋`panel_inbound.go:341` 一句注释）＋我落件一。**⛔ 不建议**：即使门禁转绿，**件二缺位 ⇒ 上面那格 C 的路被件一打开**（见 ④ 可达性）。
- **(丙) 件二走「档位层降级」**（**全在我写面内**、产线真生效，见 ④ 末段）：A/C 两格变成「问 ⇒ 拒」，**但仍需 (乙) 那条 `cmd` 腿**去改那 2 枚测试的期望（它们钉的是 `ErrNoL2Confirm`，而降级只改「谁被问」，不改 `AE→AH` 本身是否放行）。
- ⇒ **无论哪条，`cmd/wisp` 那 2 枚测试的期望都要改**；差别只在**件二的判据落在哪一层**。**这必须你拍，我不擅自换层。**

## ⑥ 纪律格

- ⛔ commit、⛔ push、⛔ 翻框、⛔ 改票名/`Status:`/`-done`、⛔ 台账；`.scratch/**` 与票 **367/371/114** 一字未碰。
- 突变/候选/探针**全在仓外**：`cand/store.go`、`cand/composer_handlers.go`、`perm_probe_test.go`、`tools_probe_test.go`、`overlay-{before,cand,tools}.json`、`logs/*.txt`（`D:\tmp\wisp\20261024\372r1\**`，**只建不删**）。
- 探针**没有**写进工作树：`git status --porcelain internal/perm internal/panel internal/tools internal/risk` 在全部跑完之后仍**空**。
