# 367-a1 普查腿 —— `perm.Store.Set` 的「变宽 → 门在不在」表（只读 / AC#0）

- 腿 `367-a1`（**只读**）。HEAD `4ee57d318`，量于 2026-10-23 18:0x–18:1x +08。
- 写面：本文件 ＋ 仓外 `D:\tmp\wisp\20261023\367a1\**`。**仓内 `.go` 改动 0 字节**
  （`git status --porcelain internal/perm internal/panel cmd/wisp` **空**；`git diff --stat HEAD` 同三包 **空**）。
- 尺：`-overlay` 把**仓外**探针挂进包（`go test -overlay=…`），探针按产线原样装配**真** `perm.Store`
  ＋**真** `config.Manager`＋**真** `config.toml`（`D:\tmp\…\work\`）。逐字读数＝`log-perm2.txt` / `log-panel2.txt` 的 `A1|` / `A1P|` 行。
- 复跑：`GOFLAGS= go test -overlay=D:/tmp/wisp/20261023/367a1/overlay-perm.json -run TestA1Census -v ./internal/perm`（panel 同法，`-run TestA1ProdChain`）。

## ① 票面三条，逐条现跑复核（行号漂没漂）

| # | 票面断言 | 现跑复认 | 行号 |
|---|---|---|---|
| 1 | `Set` 在 `internal/perm/store.go:175` | HEAD blob 同位置、工作树同 | **未漂** 175 |
| 2 | `Set` 的**产线**调用者只有 1 枚＝`internal/panel/composer_handlers.go:141` | 全仓 `\.Set\(` 逐枚分流后仍只 1 枚 | **未漂** 141 |
| 3 | 那 1 枚在 `composer_handlers.go:133` 的 `Confirm == nil` 门**之后** | `:133` `if modeIsWidening(from, to) && h.Confirm == nil {`，`:141` 在其 `if` 之外、函数尾 | **未漂** 133 |
| 4 | `114-r1` 绕过实验：深度1 拒/0 写、深度2 拒/0 写、正控过/1 写 | 用**真 `perm.Store`**（不是假 writer）重做三发，读数全中（⑥） | **复现** |

★**一处票面前提已过期（不改结论、必须报）**：`composer_handlers.go:36-40` 仍写 "…**Nothing calls this handler yet**"，
但现在**有**产线调用者 `composer_dispatch.go:204`（`return "", d.Mode.HandleModeRequest(ctx, req)`），两条入口都到（⑥）。
⇒「今天不可从页面到达」的**理由**要改成「handler 层 `:133` 拒」，⛔ 再写「跳不存在」。`114-r1-impl.md` 用 `wc -l` 复量 **152** 行、`771bf451e`，票面数对。

## ② `Set` 的全部调用者（产线／测试分开，逐枚 `file:line`）

**产线 1 枚**
1. `internal/panel/composer_handlers.go:141` — `if err := h.Modes.Set(ctx, to, PanelModeOrigin, h.actor()); err != nil {`

**测试 12 枚 / 3 文件**
- `internal/perm/store_test.go`：`:116`（`TestTicket90ThreeModesAndNoFourth`，`risk.Mode(9)`）· `:136` `:142` `:152`
  （`TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited`）· `:185`（`…NoConfirmChannelMeansNoAuto`）·
  `:204`（`…ConfirmRefusalLeavesModeAlone`）· `:225`（`…PersistFailureKeepsMemory`）
- `internal/perm/ticket90_persist_test.go`：`:154`（`TestTicket90ManualSwitchSurvivesRestart`，循环 AH/AA）· `:404`（`TestTicket90HandEditLooseningGoesThroughD36`）
- `cmd/wisp/run_mode101_test.go`：`:233` `:236`（`TestTicket101ManualSwitchSurvivesRestart`）· `:546`（`TestTicket101ModeSwitchUsesTheRealL2Gate`）

**零其它**（同名不同物已逐枚剔）：全仓 `func (…) Set(` 只有 **3 个真实现** —— `perm.Store`（`store.go:175`）
＋两个**测试假腿**（`internal/panel/composer_handlers_test.go:48` `fakeModeWriter`、`composer_dispatch_test.go:69` `hop33ModeWriter`）；
另 3 个同名不同物（`memory.ProbeFlags`、`llm.MeasuredFlags`、`cmd/wisp.taintSourceFlag`）不是 `Set(ctx,to,origin,actor)` 这形状。

## ③ 「变宽 → 门在不在」完整表

**档位是自己现跑枚举的**（`A1|LADDER|` 逐 int −3…12）：`Valid()==true` 只有 `0=ask_every_step`、`1=ask_high_risk`、
`2=auto_approve`；`ModeNames_len=3`、`default=ask_every_step` ⇒ **恰好三档、没有第四档**。

`Set` 里**唯一**看方向的条件，逐字（`store.go:191-192`）：

```go
	if to == risk.ModeAutoApprove {
		if s.confirm == nil {
```

上游没有任何"变宽"判断，`:185` 只挡 `if from == to {`。⇒ **`Set` 判的是"目标档是不是 `auto_approve`"，⛔ "是不是变宽"。**

| 从 → 到 | 方向 | `Set` 拦？ | 凭什么（逐字） | 现跑读数（`A1|PAIR|`） |
|---|---|---|---|---|
| `ask_every_step → ask_high_risk` | **变宽** | **不拦** | `:191` 不命中 ⇒ 直落 `s.mgr.SetPermissionMode(to)`（`:208`） | `set_err="<nil>" applied=1 file="permission_mode = 'ask_high_risk'"` |
| `ask_every_step → auto_approve` | 变宽 | **拦** | `:191` ＋ `:192 if s.confirm == nil {` | `last_result=refused-confirm applied=0`，文件仍 `'ask_every_step'` |
| `ask_high_risk → ask_every_step` | 变窄 | 不拦（按 R20 不该拦） | 同两行不命中 | `set_err="<nil>" applied=1` |
| `ask_high_risk → auto_approve` | 变宽 | **拦** | `:191`＋`:192` | `refused-confirm applied=0` |
| `auto_approve → ask_every_step` | 变窄 | 不拦 | 同上 | `applied=1` |
| `auto_approve → ask_high_risk` | 变窄 | 不拦 | 同上 | `applied=1` |

**总数核对**：3 对变宽 → **2 对有门（target 全为 `auto_approve`）、1 对无门**；3 对变窄 → 0 有门（设计如此）。
**上一层（唯一产线调用者自己的门）** `composer_handlers.go:107/133`：`func modeIsWidening(from, to risk.Mode) bool { return to > from }` ＋ `if modeIsWidening(from, to) && h.Confirm == nil {`
⇒ **这一层是按方向的**，而今天产线装配 `Confirm == nil`（`panel_inbound.go:264`）⇒ 上表第 1 行**在产线上被这一层拦住**（⑥ 实测）。

## ④ 「非 `auto_approve` 变宽」那一发（量出来的，⛔ 读码推断）

装配：真 `config.Manager`（`D:\tmp\…\work\nonauto_widening\config.toml`）＋ `perm.New(Options{Confirm: nil})`，
直调 `Set(ctx, risk.ModeAskHighRisk, "census", "a1")`。逐字读数：

```
A1|NONAUTO|before_read=ask_every_step|before_file="permission_mode = 'ask_every_step'"
A1|NONAUTO|from=ask_every_step|to=ask_high_risk|set_err="<nil>"|read_after=ask_high_risk|applied=1|file="permission_mode = 'ask_high_risk'"|last_result=applied|last_detail=""
A1|AUDIT|tag=nonauto|line="perm: MODE-SWITCH from=ask_every_step to=ask_high_risk at=2026-10-23T10:12:54.8738476Z origin=\"census\" actor=\"a1\" result=applied detail=\"\""
```

⇒ **改成了 `ask_high_risk`**：`err=<nil>`、`result=applied`、`detail=""`（**没有一句"为什么变宽"**）、文件真落盘、`applied=1`。⇒ **没有门** —— 不是"读码看着没有"，是**零确认、零拒绝地真的改了**。

## ⑤ 正控（硬闸）：切进 `auto_approve` 必须有门

同装配同代码路径，只把 `to` 换成 `auto_approve`。**逐字拒绝句**（`store.go:194` 的 `sw.Detail`，经 `:196` 包 `perm: ` 前缀）：

```
perm: 切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）
```

配套读数：`last_result=refused-confirm`、`read_after=ask_every_step`、`applied=0`、`file="permission_mode = 'ask_every_step'"`（**没改**）。
**尺的双向**（证明它不是"一律拒"）：`Confirm` 接上答 yes ⇒ `confirm_calls=1 set_err="<nil>" read_after=auto_approve applied=1`；
答 no ⇒ `set_err="perm: 切到 auto_approve 被拒绝: operator clicked deny" applied=0`。
⇒ 同一把尺能读出"有门"（`refused-confirm`）与"无门"（`applied=1`）。

## ⑥ 今天有几条路能到 `Set`（决定"加固"还是"补洞"）

**产线：调用点 1 枚；入口 2 枚，装同一条链。**
- `cmd/wisp/panel_host_windows.go:948` — `return m.disp.Handle(ctx, raw)`（resident panel 页面，`dispatchRaw` 在 `:941`）
- `cmd/wisp/panel_inbound.go:164` — `reply, err := disp.Handle(ctx, raw)`（CLI seam，stdin 一行一封）
- 两入口 → 同一 `ComposerDispatch.Handle`（`composer_dispatch.go:176`）→ `dispatch` 的 `case MethodModeRequest`（`:200-204`）
  → `HandleModeRequest`（`composer_handlers.go:111`）→ `:133` 门 → `:141 Set`
- **两层 `Confirm` 都是 nil**（同一构造器 `newComposerDispatchChain`，`panel_inbound.go:240`）：
  `perm.Options{… Confirm: nil …}`（`:254`）、`ModeWriteHandler{… Confirm: nil …}`（`:264`）

**测试：12 枚调用点 / 3 文件**（② 逐枚列）—— 它们**绕过** handler 门直调 `Set`，测的是 `Set` 自己的语义。

**每条是否都已过门（实测）**：

```
A1P|DEPTH1|to=ask_high_risk|handler_err="panel: 确认腿未接入，档位不得变宽: …"|is_ErrNoL2Confirm=true|store_read=ask_every_step|set_writes=0
A1P|DEPTH2|to=ask_high_risk|handle_err=同上|is_ErrNoL2Confirm=true|store_read=ask_every_step|set_writes=0
A1P|POSCTL|to=ask_high_risk|handle_err="<nil>"|reply=""|handler_confirm_calls=0|store_read=ask_high_risk|set_writes=1
```

⇒ 产线路**已过门**（变宽在 `:133` 被拒、写腿 0 次）；正控把 `Confirm` 接上即过、写腿 1 次 ⇒ **`114-r1` 三条读数复现**（这次用的是真 `perm.Store`）。`handler_confirm_calls=0` ⇒ handler **只读 `Confirm` 的 nil-ness、从不调用它**（与 `composer_handlers.go:21-28` 一致）。

★**一条"像路但不是路"的，具名报上来**：`cmd/wisp/run.go:687` 建了一枚**带真 `Confirm`** 的 `ModeWriteHandler` 存进 `rt.modeWrites`
（声明 `run.go:340`）；现跑 grep 复认**无任何读点**（声明 1＋赋值 1、读 0），`run.go:333` 自陈 "Nothing reads this field yet"
⇒ **不算第二条产线路**；但它是"哪天有人把它接到 dispatch 上"的**现成入口** —— 那一刻 `ask_every_step→ask_high_risk` 会**过门进 `Set`**，而 `Set` 对它**没有门**。

**另两条"改档位但不经过 `Set`"的路（都不通向 `Set`）**：① `Manager.SetPermissionMode` 产线唯一调用者就是 `store.go:208`（`Set` 内部）；
② 设置页 `config.set` 按名字拒它 —— `config_handlers.go:131-132`：`case "permission_mode", "permissionmode", "mode": return "risk.", "权限档位的变更只由 panel.mode.request 发起（变宽要那张 L2 卡），设置页不写它"`。
手改 `config.toml` 走 `CheckAndReload`＋`ConfirmLocked`（`config_reload.go:118`），不经过 `Set`，不在本票射程。

**判**：**加固（latent seam），⛔ 补洞（live hole）**。判据＝**产线路 1 条、且它唯一的上层门今天实测拒**；
"活洞"需要"第二条产线路"或"产线门被绕"，两者今天都没有。★但 ④ 已证 **`Set` 自己这层对非 auto 变宽确实无门**，且 `Set` 是 exported ⇒
上面那枚现成入口一接上，**门就没了**。

## ⑦ 门禁 rc

| 命令 | rc |
|---|---|
| `GOFLAGS= go build ./...` | **0** |
| `go vet ./internal/perm ./internal/panel` | **0** |
| `go test -count=1 ./internal/perm` | **0** |
| `go test -count=1 ./internal/panel` | **1**（既有红 5 枚，见下） |
| `sh scripts/d22scan.sh` | **0**（clean） |
| `sh scripts/check-path-length-budget.sh` | **0**（VERDICT GREEN） |
| 两发 overlay 探针（`TestA1Census` / `TestA1ProdChain`） | **0 / 0** |

`internal/panel` 那 5 枚红是 **HEAD 在此工作树上的既有状态**（本腿 0 字节改动；与 `114-r1-impl.md` 记的 5 枚**同集**）：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`（`frontend/**` 不在本工作树）、
`TestC21DesignTokensFourWayAgree`（`design/assets/tokens.css` 不存在）、`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、
`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（票 35 的内存上界）。**⛔ 本腿写面，未动。**
格式尺：本腿**零 `.go` 改动** ⇒ 甲乙两尺都不适用（无入库 blob 变化）。

## ⑧ 提交

- 提交号：**`cc45a4b81`**（本文件入库那一笔；`git log -2 -- docs/evidence/s1/367-a1-census.md` 可见紧随其后的号回填笔，只改本节这一行）。
- 逐枚文件（每笔各只含这 1 个）：`docs/evidence/s1/367-a1-census.md`
- 逐枚 `.go`：**0**

## 附：尺与射程（落笔前五问的自答）

1. **量的是哪个对象**：`*perm.Store.Set` 的**运行行为**（真 `Manager`＋真 `config.toml`），⛔ 注释、⛔ 票面。
2. **输出形状**：一行一读数、`|` 分隔；错误一律 `<nil>` 或**原文**（`%q` 加引号），不转述。
3. **数的是"代码"还是"注释"**：`func (…) Set(` **3 个真实现**（1 产线＋2 测试假腿）＋3 个**同名不同物**分开数；
   `rt.modeWrites` 的"无读点"是**代码引用**计数（声明 1＋赋值 1、读 0），不是注释。
4. **和总数对上**：测试调用点 **12** ＝ `store_test` 7 ＋ `ticket90_persist` 2 ＋ `run_mode101` 3；全仓 `\.Set\(` 其余命中
   **全是 `Header.Set`/`reflect.Value.Set`/`jsValue.set` 等同名不同物**，已逐枚剔。
5. **谁算的/算法/几位**：行数用 `wc -l`（`114-r1-impl.md` = **152**，与票面一致）；⚠ PowerShell `Measure-Object -Line` 在同一文件给 **113**
   —— **那是错尺、已弃用**（`ruler-scope-checklist` 第 1/7 例同型）。
