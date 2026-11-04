# 371-a1p 只读腿 —— `AC#0`②③ ＋ **全部 `AC#0'`**（读数、谓词普查、洞、跳数、`[bing-a]` 的契约面与爆炸半径）

- 腿 `371-a1p`（**只读**，⛔ 实现者）。HEAD **`786a13b22`**，量于 2026-11-04 21:4x–22:0x +08。
- 写面：**本文件** ＋ 一枚聚焦用例 `internal/perm/store_371_a1p_test.go`（新建，只加断言、⛔ 动产码字节）。
- 产码字节：**0**。`git status --porcelain` 里没有 `internal/**`／`cmd/wisp/**` 的 ` M`；本腿只新增上述两枚文件。
- ★所有行号一律取自 **HEAD blob**（`git show HEAD:<path>`），⛔ 工作树（理由＝`371-a1-census` ⑥.3 的同一条纪律）。
- 本腿**不**改票、不翻框、不碰 367／372／375 原句、不碰 `docs/PLAN.md`／`docs/specs/**`／`frontend/**`／`design/**`／`thresholds.go`／任何 golden。

## ⓪ 一屏判（六句，逐条落点见后）

1. **A（`AC#0`②）复现成立**：`AE→AH` ＋ **真腿** ⇒ `invoked=0`、`asked=0`、`set_err=<nil>`、`last_result=applied`、`persisted=ask_high_risk`、`saves=1` ⇒ 票面的 `confirm_calls=0, applied=1` **成立，且比 367 时代更强**（连"调腿"都没有，因为 `:268`／`:281` 都在 `RequiresCard(ask_high_risk)==false` 上短路）。
2. **B（`AC#0`③ 正控）成立**：同一真腿、只把目标换成 `auto_approve` ⇒ `invoked=1`、`asked=1`、`applied` ⇒ **尺读得出"问过"**，所以 A/D 的 0 不是假零。另跑产线装配 `TestTicket101ManualSwitchSurvivesRestart` **PASS**（它逐字断言 `asked == [auto_approve]`）。
3. **C（`AC#0'`①）结论**：`HEAD` 上判「变宽／吃卡」的 5 处谓词与调用点，**没有一处**判「问过／被答过」；门的两个合取项只判 ①**目标档身份**（`RequiresCard`）②**`Confirm` 是否为 `nil`**。严格点名例外：`RequiresCard` 判的是**目标档**、`IsWidening` 判的是**方向**（且已**不是门**，产线调用点 **0**）。
4. **D（`AC#0'`② 靶心）洞是真的**：`AA` ＋ **非 `nil` 但从不问**的假腿 ⇒ `invoked=1`、`asked=0`、`set_err=<nil>`、`applied`、`persisted=auto_approve`、`saves=1`。
5. **E（`AC#0'`③）**：门放行 → 卡被答 共 **6 跳**（另加审计行＝第 7 跳）。**不作判据的那一跳＝`Set` 里 `s.confirm(ctx, sw)` 的返回值**（`store.go:283`）：`nil` 只被用来把结果记成 `applied`，`err` 只被 `classifyConfirmFailure` 分类；**没有任何一跳要求"有人看过这张卡"**。
6. **F（`AC#0'`④）同洞成立**：真腿与假腿在 `AA` 上得到**同一组判决**（`refused=false`／`result=applied`／`persisted=auto_approve`／`saves=1`），只有 `asked` 不同（1 vs 0）⇒ **今天的 `AA` 门与 D 是同一个洞**。
7. **G（契约面）**：`ConfirmFunc`／`Options`／`Switch` **不是** `C1–C32`／`D1–D47` 契约文字（`docs/PLAN.md`＋14 份 `specs` 里 `ConfirmFunc`／`RequiresCard`／`ModeConfirm`／`MODE-SWITCH`／`permission_mode`／`auto_approve` 全部 **0 命中**）⇒ `[bing-a]` 是**内部 API 决定**、⛔ 触发"改契约文字"的人工批准。**爆炸半径**见 `§G.2`（产线 2 支适配器／2 处装配；测试 6 枚文件构造 `ConfirmFunc`，另 2 枚碰 `panel.ModeConfirm`）。**推荐 `[bing-a]`（并把 `[bing-b]` 的审计读数并进同一波）**。

---

## ① 尺：命令、仪器、写面

### ①.1 仪器（本腿新建、已入库）

`internal/perm/store_371_a1p_test.go` —— 4 枚用例，**只断言、只打印**，不碰产码。尺的关键是**两个计数器**：

| 计数 | 含义 | 为什么必须分开 |
|---|---|---|
| `invoked` | `Set` **调用了**这枚 `ConfirmFunc` | "接线被用到"≠"有人被问" |
| `asked` | 这枚腿**真的出了卡**并拿回 `allow` | `confirm_calls` 在票面上指这个 |

- **真腿**：`invoked++` **且** `asked++`，返回 `nil`（＝答了允许）。这正是产线腿的契约：`cmd/wisp/run.go:1126-1132` 的 `confirmModeSwitchVia` 只有 `ans == tools.AnswerAllow` 才 `return nil`。
- **假腿**：`invoked++` **只此一句**，`return nil` —— **非 `nil`、从不问**（`AC#0'`② 点名的新形状）。
- ★**保真度声明（⛔ 冒充"我点了真卡"）**：本机没有可点的原生卡面；"真腿"按**本仓自己的约定**建模 —— `cmd/wisp/run_mode101_test.go:225-228` 逐字说那枚 stub「stands for "the native side clicked allow once"」，`panel_mode_confirm_114_windows_test.go` 的 `countingPermConfirm` 同形。⇒ 本腿量的是**`perm` 这一层"有没有被要求问／有没有真问"**，⛔ 量 C18 队列的 300s 行为（那归 375／C18 用例）。

### ①.2 命令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"

/d/work/base/gopath/bin/gofumpt.exe -l internal/perm/store_371_a1p_test.go   # 输出空 ⇒ 已格式化
go vet ./internal/perm                                                       # rc=0
go test ./internal/perm -run 'Test371A1P' -count=1 -v                        # 4 枚 PASS，读数见 §A–§D
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart' -count=1 # PASS (41.7s)
```

### ①.3 写面核对

```bash
git status --porcelain | grep -E '^ ?M'        # 空 ⇒ 无产码改动
git status --porcelain | grep -E '371_a1p|371-a1p'   # 只有本文件＋用例
```

---

## ② A —— `AC#0`②：`AE→AH` ＋ **真 `Confirm`**（`confirm_calls` 应为 0）

**命令**：`go test ./internal/perm -run 'Test371A1PAC0RealLimbOnTheMiddleTier' -count=1 -v`

**逐字读数**（`store_371_a1p_test.go` 的 `t.Logf`）：

```
371A1P|AC0-2 AE->AH real limb|ask_every_step->ask_high_risk|invoked=0|asked=0|set_err=<nil>|last_result=applied|persisted=ask_high_risk|saves=1
```

**判**：

| 指标 | 读数 | 票面期望 | 判 |
|---|---|---|---|
| `confirm_calls`（＝`asked`） | **0** | 0 | ✅ 复现 |
| `applied` | **1**（`last_result=applied`、`persisted=ask_high_risk`、`saves=1`） | 1 | ✅ 复现 |
| `invoked`（新增的更强读数） | **0** | —（票面没要求） | ★ 372 之后**连腿都没被碰**：`:268` 与 `:281` 的 `RequiresCard(ask_high_risk)==false` 直接短路 |

★与 `371-a1-census` ③ 的 `leg1.out` 同值（`confirm_calls=0`／`applied`），**独立复跑、不抄**。

---

## ③ B —— `AC#0`③：**正控** `AA` ＋ 真 `Confirm`（必须 `confirm_calls=1`）

**命令**：`go test ./internal/perm -run 'Test371A1PAC0PositiveControlOneCardOnAutoApprove' -count=1 -v`

**逐字读数**：

```
371A1P|AC0-3 AE->AA real limb (positive control)|ask_every_step->auto_approve|invoked=1|asked=1|set_err=<nil>|last_result=applied|persisted=auto_approve|saves=1
```

**判**：`asked=1` ⇒ **尺能读出"问过"**。⇒ A/D 的 `0` 是**真零**，⛔ 仪器读不到造成的假零。

### ③.1 产线装配的旁证（端到端、真 `perm.Store` ＋ 真适配器）

**命令**：

```bash
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart' -count=1
```

**读数**：`--- PASS: TestTicket101ManualSwitchSurvivesRestart (41.67s)`；该用例在 `run_mode101_test.go:249-252` 逐字断言：

```go
if len(gotAsked) != 1 || gotAsked[0] != risk.ModeAutoApprove {
    t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", gotAsked)
}
```

而它的 hook（`:236-245`）先做 `AE→AH`、再做 `AH→AA`；两条都成功、退出码 0 ⇒ **`gotAsked` 恰为 `[auto_approve]`** ⇒ 在**产线装配**上：`AE→AH` 的 `asked=0`、`AH→AA` 的 `asked=1`。★这是 B 的产线旁证，⛔ 替代 §②/§③ 的直读。

---

## ④ C —— `AC#0'`①：所有判「变宽／吃卡」的谓词与调用点

**命令**：

```bash
git grep -n "RequiresCard\|IsWidening" HEAD -- "*.go" ":!*_test.go"
git show HEAD:internal/perm/store.go | grep -n "RequiresCard\|s.confirm"
git show HEAD:internal/panel/composer_handlers.go | grep -n "RequiresCard\|Confirm == nil"
git grep -n "IsWidening(" HEAD -- "*.go"
```

| # | `file:line`（HEAD blob） | 逐字 | **判什么属性** | 能证明"问过"吗 |
|---|---|---|---|---|
| 1 | `internal/perm/store.go:268` | `if RequiresCard(to) && s.confirm == nil {` | ①目标档是 `auto_approve`（经 `direction.go:23`）②`s.confirm` 字段**是 `nil`** | ❌ 两个合取项都与"问过"无关 |
| 2 | `internal/perm/store.go:281` | `if RequiresCard(to) {` | **只**判目标档；决定"要不要调那一发" | ❌ |
| 3 | `internal/perm/store.go:283` | `if err := s.confirm(ctx, sw); err != nil {` | 判**调用的返回 err**；`nil` 落到 `:297 sw.Result = ResultApplied` | ❌ `nil` 也可能是"没问"（§D 实测） |
| 4 | `internal/panel/composer_handlers.go:153` | `if perm.RequiresCard(to) && h.Confirm == nil {` | 同 #1 的两项，对象换成 handler 的 `ModeConfirm` 字段；**它从不调用这枚腿**（`:157-160` 的注释逐字「Set judges the same target with the same predicate」） | ❌ 连调用都没有 |
| 5 | `internal/perm/direction.go:23` | `func RequiresCard(to risk.Mode) bool { return to == risk.ModeAutoApprove }` | **目标档身份**（唯一一枚返回 `true` 的档） | ❌ |
| 6 | `internal/perm/direction.go:36` | `func IsWidening(from, to risk.Mode) bool { return to > from }` | **阶梯序**（`risk.Mode` 按 strict→loose 排）；`direction.go:30` 逐字「It is NOT a gate any more (ticket 372)」 | ❌ 且**已不是门**：产线调用点 **0**，只有 2 枚测试文件 5 处（`composer_handlers_test.go:114/118`、`store_367_test.go:209/212/217`） |

### ④.1 结论（按票面必须二选一的那句）

- ★**"全部只判 `nil-ness`"要加一个精确注**：门的**两个合取项**里，判**通道属性**的那一项确实是**只判 `nil-ness`**（`s.confirm == nil`／`h.Confirm == nil`）；另一项 `RequiresCard(to)` 判的是**目标档身份**（不是 nil-ness）。
- ★**例外点名**：`RequiresCard`（`direction.go:23`）判**目标**；`IsWidening`（`direction.go:36`）判**方向**但**已不作门**。
- ★★**本条的承重结论**：**6 处里没有任何一处判"这张卡有没有被答过"** —— 这个"缺席"就是 `AC#1'` 的靶。★措辞上请按"**判目标 ＋ 判 nil-ness；⛔ 判 answered**"引用，别写成"全都在判 nil-ness"（那会让 `RequiresCard` 被误读成通道谓词）。

---

## ⑤ D —— `AC#0'`②：**活洞**（`AA` ＋ 非 `nil` 但从不问的假腿）

**命令**：`go test ./internal/perm -run 'Test371A1PAC0PrimeFakeLimbIsSilentlyApplied' -count=1 -v`

**逐字读数**：

```
371A1P|AC0p-2 AE->AA fake limb (the hole)|ask_every_step->auto_approve|invoked=1|asked=0|set_err=<nil>|last_result=applied|persisted=auto_approve|saves=1
```

**判**：

| 指标 | 读数 | 含义 |
|---|---|---|
| `invoked` | **1** | `Set` 确实调了这枚腿（`:283`） |
| `asked`（＝`confirm_calls`） | **0** | 一张卡都没出 |
| `set_err` | `<nil>` | 腿返回 `nil` ⇒ 被当成"允许" |
| `applied` | **1** | `persisted=auto_approve`、`saves=1` ⇒ **最宽档落盘** |

⇒ ★★**这就是 `AC#1'`〔丙〕要治的活洞**：`perm` 只看到"非 `nil`"。★它**不属于** 372 的读侧降级（372 只在 `s.confirm == nil` 时把有效档塌到 `ask_every_step`；本形状非 `nil` ⇒ 不触发），★也**不属于** 375 的动作层 fail-closed（375 治"不可问"⇒拒那一发；本形状**可问而不问**）。⇒ **`AC#5'`①/② 的边界成立**：本洞是第三层。

---

## ⑥ E —— `AC#0'`③：从「门放行」到「卡被答」有几跳，哪一跳不作判据

**命令 / 尺**：读 HEAD blob 的 `store.go:268-299`、`run.go:1104-1157`、`panel_inbound.go:335-356`、`approval/gate.go:582`、`approval/queue.go:443-453`。下表的"跳"从**门放行之后**算起（门本身＝`store.go:268`，另计）。

| 跳 | 落点（HEAD） | 逐字/形状 | 在 HEAD 上是判据吗 |
|---|---|---|---|
| H1 | `internal/perm/store.go:281` | `if RequiresCard(to) {` | 判"目标是 AA"，**不判问了没** |
| H2 | `internal/perm/store.go:283` | `if err := s.confirm(ctx, sw); err != nil {` | **❌ 只判"调用没返回 err"**；调用≠问 |
| H3 | `cmd/wisp/run.go:1143` / `cmd/wisp/panel_inbound.go:346` | `panelModeConfirm`／`storeConfirmFromPanel`：nil in nil out，纯类型翻译 | ❌ 不携带任何证据 |
| H4 | `cmd/wisp/run.go:1119-1125` | C18 入队并等：`gate.PendingApproval(ctx, tools.Decision{Level: risk.L2, …})` | 这一步**真的出卡**（可问性由 C18 决定） |
| H5 | `internal/agent/approval/queue.go:443-453` | 原生点击 → `answer{a: tools.AnswerAllow}` → 审计行 `ANSWER-ALLOW corr=… decision=allow` | ✅ **证据在这一跳产生**（`ans`、`corr`） |
| H6 | `cmd/wisp/run.go:1126-1132` | `if ans != tools.AnswerAllow { … } return nil` | ✅ **腿内部判**了"答的是允许" |
| H7 | `internal/perm/store.go:283-289` ＋ `:297` | `err == nil` ⇒ `sw.Result = ResultApplied` | ❌ **这一跳不作判据**：返回 `nil` 只被当成"没被拒"，⛔ 没被当成"有人答过" |
| H8 | `internal/perm/store.go:358-370` | `MODE-SWITCH … result=applied detail=""`（`Secs` 是秒数） | ❌ 行里**没有**"问过/答过"字段 |

**判**：

- ★**跳数**：门放行 → 卡被答共 **6 跳（H2–H7）**；把审计行算上是 **7 跳**。
- ★★**不作判据的那一跳**：**H7（`Set` 里 `s.confirm` 的返回值）** —— `err` 只喂 `:284 classifyConfirmFailure`（`store.go:305-313`：只分"被拒／出错"），`nil` 只落到 `:297 ResultApplied`。**"有人看过这张卡"在 `perm` 侧没有任何判据**。H2/H3 同理（调用与翻译都不带证据）。
- ★**要紧的细节**：证据在 **H5/H6 已经存在**（`ans == tools.AnswerAllow`、`ANSWER-ALLOW corr=…`），**只是过了 H7 就被丢成 `nil`**。⇒ `[bing-a]` 的楔子就在这里：让 H6 的那枚证据**穿过 H7**。
- ★**一处要写清的边界**：`perm.ConfirmFunc` 的文档（`store.go:95-103`）已经把"`nil` = allowed"写成语义；今天缺的是**可判性**，⛔ 缺定义 —— 这也是为什么 `[bing-a]` 是"收紧证据"而不是"发明语义"。

---

## ⑦ F —— `AC#0'`④：今天的 `AA` 门有没有**同样的洞**？

**命令**：`go test ./internal/perm -run 'Test371A1PAC0PrimeTheAutoApproveGateHasTheSameHole' -count=1 -v`

**逐字读数**：

```
371A1P|AC0p-4 answered_limb|refused=false|result=applied|persisted=auto_approve|saves=1|asked=1
371A1P|AC0p-4 never-asking_limb|refused=false|result=applied|persisted=auto_approve|saves=1|asked=0
```

**判**：★**有，同一个洞。** 判决四元组 `(refused, result, persisted, saves)` **逐项相同**，唯一分叉的是 `asked`（1 vs 0）⇒ 门的判决是 **`Confirm != nil` 的函数**，而两枚腿都满足它。★这正是 `AC#2'` 验收判据要打的那一格：假腿 ⇒ 必须 `applied=0` **且** 拒因逐字点名"这张卡没被答过"。

---

## ⑧ G —— 供你裁 `AC#1'`：`[bing-a]` 的契约面 ＋ 爆炸半径 ＋ 与 `[bing-b]` 的代价对照

### ⑧.1 契约面：`ConfirmFunc`／`Options`／`Switch` **不是** C 级（也不是 D 级）契约

**尺（逐条现跑，HEAD blob）**：

| 命令 | 读数 | 判 |
|---|---|---|
| `git grep -c "ConfirmFunc" HEAD -- docs/PLAN.md docs/specs/*.md` | **0 行、rc=1** | 类型名不在契约文档 |
| `git grep -c "RequiresCard\|ModeConfirm\|MODE-SWITCH" HEAD -- docs/PLAN.md docs/specs/*.md` | **0 行、rc=1** | 谓词／审计行也不在 |
| `git grep -n "permission_mode" HEAD -- docs/PLAN.md docs/specs/*.md` | **0 行** | 配置键不在 |
| `git grep -n "auto_approve\|ask_high_risk" HEAD -- docs/PLAN.md docs/specs/*.md` | **0 行** | 三档词面不在 |
| `git grep -n "权限模式" HEAD -- docs/PLAN.md docs/specs/*.md` | **1 行**：`PLAN.md:2138`（**Claude Code 借鉴表**的一行，⛔ 决策条目） | 只是提到别家的功能 |
| `git grep -n "三档" HEAD -- docs/PLAN.md docs/specs/` | 命中全是**别的三档**（D36 配置生效级别 `:2731/2736`、动效 `:3371`、模型镜像三档 `SPEC-10:141`/`SPEC-11:93`、`input.type` 授权三档 `:2034`） | **没有一处**是权限三档 |

**`C1–C32` 名单逐枚核过**（`PLAN.md:1357-1388`）：`C1 Tool`…`C18 ApprovalQueue`…`C32 RealtimeEngine`。**没有** `perm.Store`／`ConfirmFunc`／`Options`／`Switch` 任何一枚。最近的是 **`C18 ApprovalQueue`**，而它冻结的是**队列行为**：

> `PLAN.md:1374`（C18 行）逐字：`全局 FIFO；每项含 correlationId / 所属任务标识 / 工具名 / 完整参数 / 风险级 / 请求理由。队头单显 + 深度计数；超时 = 300s，一律判拒绝…；「允许」决策只接受原生侧来源（D33/F2，面板来源直接拒绝）`

⇒ **C18 冻结的是"卡怎么排队、怎么判、谁能答"，⛔ 冻结 Go 函数签名。**

**`D` 面**：`D4`（`PLAN.md:122`）是 **L0/L1/L2 风险级**的确认门控（另一根轴）；`D31`（并发/审批）、`D33`（安全必修五项）、`D45`（批量/会话授权）都**没有** `ConfirmFunc` 的形状。★另：本腿**不碰** `M4`／`R20` 的**产品行为**（`AE→AH` 仍不出卡）—— 那是机主 2026-09-21 定案（见 `371-a1-census` ①.2）。

**答**：

1. ★**`ConfirmFunc` 签名／`Options`／`Switch` ⛔ `C` 级契约、也 ⛔ `D` 级契约**（0 命中）⇒ `[bing-a]` **不触发 `SPEC-12 §4.1` 的"改契约文字＝人工批准"**，属编排者可拍的**内部 API 决定**。
2. ★**但有一处必须写进裁形记录**：`[bing-a]` 会让 `perm` 的"允许"**带上证据**，这**加强**（⛔ 违背）C18 的「允许只接受原生侧来源」与 D4 的「L2 必须显式确认」。
3. ★★**票面〔丙-a〕的第二形（"`Switch` 上落 `Answered bool` ＋ 时间戳"）要更正**：今天的 `type ConfirmFunc func(ctx context.Context, req Switch) error` **按值传入 `Switch`、只返回 `error`** ⇒ 腿**无法**把 `Answered` 塞回 `Switch`。要用 `Switch` 当证据载体，必须让 `ConfirmFunc` **返回 `Switch`（或收 `*Switch`）** —— 那**同样是签名变更**，⛔ "比返回凭据更小的动作"。⇒ 裁形时**两个子形其实是同一量级**，别把第二形当零成本。

### ⑧.2 爆炸半径（`git grep` 计数 ＋ 逐文件，HEAD blob）

| 面 | 计数 | 逐文件（HEAD） |
|---|---|---|
| `ConfirmFunc` 文本提及 | **16 行 / 8 文件**（其中带 `perm.` 限定 **8 行 / 6 文件**；`store.go`／`store_test.go` 用包内未限定写法） | `internal/perm/store.go`（:95/:104/:117/:129）· `cmd/wisp/run.go`（:126/:889/:1134/:1143）· `cmd/wisp/panel_inbound.go`（:336/:346）· `cmd/wisp/run_mode101_test.go`（:145）· `cmd/wisp/permission_mode_effective_377_test.go`（:42）· `cmd/wisp/panel_mode_confirm_114_windows_test.go`（:67）· `internal/perm/store_test.go`（:99）· `internal/tools/booked_mode_377_test.go`（:119/:264） |
| **产线** `perm.New(perm.Options{…})` | **2 处** | `cmd/wisp/run.go:867`（`Confirm: confirm`，:869）· `cmd/wisp/panel_inbound.go:272`（`Confirm: storeConfirmFromPanel(confirm)`，:277） |
| **产线**适配器（两个方向） | **2 支** | `cmd/wisp/run.go:1143 panelModeConfirm`（perm→panel）· `cmd/wisp/panel_inbound.go:346 storeConfirmFromPanel`（panel→perm） |
| `panel.ModeConfirm`（若 `AC#2'-b` 要求两处同判） | **1 类型 ＋ 1 字段 ＋ 3 构造点** | `internal/panel/composer_handlers.go:98`（type）/`:106`（字段）；构造：`cmd/wisp/run.go:894`、`cmd/wisp/panel_inbound.go:288`、`cmd/wisp/resident_windows.go:301`（`panelModeConfirm(ra.confirmModeSwitch)`） |
| **测试**构造/接受 `ConfirmFunc` 的文件 | **6 枚**（＋本腿 1 枚） | `internal/perm/store_test.go`（`newTestStore` 参数 ＋ 1 枚字面）· `internal/perm/ticket90_persist_test.go`（3 枚字面）· `internal/tools/booked_mode_377_test.go`（参数 ＋ 1 枚）· `cmd/wisp/run_mode101_test.go`（`start` 参数 ＋ stub）· `cmd/wisp/permission_mode_effective_377_test.go`（参数 ＋ 1 枚）· `cmd/wisp/panel_mode_confirm_114_windows_test.go`（`countingPermConfirm.fn`） |
| **测试**构造 `panel.ModeConfirm` 的文件 | **2 枚** | `internal/panel/composer_handlers_test.go`（`countingConfirm.fn`，:59）· `cmd/wisp/panel_mode_confirm_114_windows_test.go` |
| `perm.Store.Set` 调用点（**签名不变则不受影响**） | 产线 **2**（interface 调用）／测试 **22 文本点 / 7 文件** | 产线：`composer_handlers.go:147`、`:163`（`h.Modes.Set`）。测试：`store_test.go` 7 · `store_367_test.go` 5 · `store_372_test.go` 3 · `ticket90_persist_test.go` 2 · `booked_mode_377_test.go` 1 · `run_mode101_test.go` 3 · `composer_gate_agreement_367_test.go` 1 |
| `Switch{…}` 字面构造（**仅当 `Switch` 加字段才动**，且加字段不改签名） | **5 处** | `store.go:159/168/248/255` · `run.go:1148` |

⇒ **`[bing-a]` 的真实半径**：**产线 4–5 个文件**（`internal/perm/store.go`、`internal/panel/composer_handlers.go`、`cmd/wisp/run.go`、`cmd/wisp/panel_inbound.go`，若连带 panel 同判再加 `resident_windows.go` 的调用点）＋ **测试约 8–9 枚文件**（上表 6＋本腿＋2 枚 panel；`store_367/372` 与 `ticket90_persist` 里的 `Set` 调用点**不必改**，只改**夹具的腿值**）。★**`Set` 的签名不必动**（⇒ 22 个 `Set` 调用点原地不动），这是 `[bing-a]` 比"给 `Set` 加参数"便宜的地方。

### ⑧.3 代价对照表（**你裁**）

| 维度 | **`[bing-a]`**：`ConfirmFunc` 返回非零凭据，`Set` 判"答过" | **`[bing-b]`**：不动门/契约，只把"问过"变成可观测读数 |
|---|---|---|
| 改契约文字（`C1–C32`／`D1–D47`） | **0**（§G.1 已量） | **0** |
| 产线文件 | **4**（perm 门＋类型、panel 门＋类型、2 支适配器所在的 `run.go`／`panel_inbound.go`） | **1**（`cmd/wisp/run.go` 的产线腿加一行审计；`perm` 侧可选加一字段） |
| 测试文件 | **≈8–9**（§G.2；`Set` 调用点不改） | **1–2**（新增断言＋1 枚既有审计格式用例） |
| **洞是否闭合** | ✅ **闭合**：假腿 ⇒ `err==nil` 但凭据为零 ⇒ 拒，且拒因可点名"卡没被答过"（＝`AC#2'` 的验收判据） | ❌ **不闭合**：只让洞**有读数** |
| 用户可见价值 | ✅ 切到"全自动"时，**要么真弹卡并等到答案、要么不许放行** | ⚠ 只让事后从日志看出"这次拨全自动没弹卡"（〔丁〕的诉求） |
| 两处门同判（`AC#2'-b`） | ⚠ **必须**同改 `panel.ModeConfirm`（否则 `perm` 判、panel 不判 ⇒ 你点名的"只改一处必须转红"那格会红） | ✅ 不涉及（两处都仍是 nil-ness） |
| 反主流核对（`m4-survey.md` ④/⑤） | ✅ 主流把确认挂在"这一发真的被确认过"上（Step-Code `effectiveState`、dsh `permission/preset` 事件、openchamber `displayedPermissionMode`）⇒ 门的判据应是**答案**，⛔ **接线存在** | ⚠ 主流**没有**"只记读数、不判"这一形；它是本仓为"让洞可断言"加的辅助面 |
| 主要风险 | ①凭据是**同进程**铸造 ⇒ 挡的是**事故性的假腿**（本票的形状），⛔ 挡恶意代码；②`panel` 侧必须同波改，否则两层判据分叉（正是 `AC#2'-b` 要防的） | ①**票面字面"`Set` 落 `asked=1/0`" 不可实现**：`Set` 只能观测"调用过"，⛔ 观测"问过"；要让读数**诚实**，`asked` 必须由**腿**发（产线腿 `run.go` 有 `rt.auditf`，可发 `MODE-SWITCH-CARD corr=…`，与既有 `ANSWER-ALLOW corr=…`（`queue.go:452`）**天然可 join**）⇒ 这是**对票面措辞的一处更正**，⛔ 改票；②洞仍在，用户照旧可能被切进"全自动"而没人看过卡 |
| 与 `AC#1'` 的关系 | 就是**〔丙-a〕本体** | 是**〔丙-b〕**，票面自己说"推荐与〔丙-a〕同批" |

### ⑧.4 我的推荐（**一行为由**）

> **选 `[bing-a]`，并把 `[bing-b]` 的审计读数并进同一波**：`[bing-a]` 的值是**把"允许"从"接线存在"改成"有人答过"** —— 这是主流 harness 判"改档"的判据（`m4-survey.md` ④），而 `[bing-b]` 单独交付时**用户仍会在没人看过卡的情况下被切进全自动**，只买到"事后能看见"，买不到"洞被堵"。

- ★**成本护栏（⛔ 被误读成"大动作"）**：`[bing-a]` **不动 `Set` 签名**（22 个 `Set` 调用点原地不动）、**不动 `C`／`D` 文字**、**不动 `M4` 产品行为**；产线 4 个文件、测试 8–9 枚，**其中一半是夹具的腿值**。
- ★**必须写进裁形记录的一处**：`[bing-b]` 的"`Set` 落 `asked=1/0`"要落成"**腿发审计行 ＋ `corr` join**"，否则那个 `asked` 是假的（`Set` 只能看见 invocation）。
- ★**若你只批 `[bing-b]`**：那也**必须**同时把 `§⑤` 的读数固化成一条"红线用例"（今天的 `applied` 必须被记录成"未答已放行"），否则洞会继续静默 —— 这**仍然是 `[bing-b]` 的下限**。

---

## ⑨ 门禁 rc（本腿）

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/perm/store_371_a1p_test.go` | **0** | **输出空**（无待格式化） |
| 2 | `go vet ./internal/perm` | **0** | 静默 |
| 3 | `go test ./internal/perm -run 'Test371A1P' -count=1 -v` | **0** | 4 枚 **PASS**；读数为 `§②/§③/§⑤/§⑦` 的逐字行 |
| 4 | `go test ./internal/perm -count=1` | **0** | 整包 PASS（本腿用例未扰动既有格） |
| 5 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart' -count=1` | **0** | `--- PASS`（41.67s） |
| 6 | `git status --porcelain`（产码面） | **0** | 无 `internal/**`／`cmd/wisp/**` 的 ` M` |

★**未跑**：`./cmd/wisp` 整包（会死在 `MARK299`，按纪律只跑定向 `-run`）；`go build ./...`（本腿 0 产码字节，⛔ 用别人的红冒充本腿读数）。

---

## ⑩ 附：本件的边界与可复核点

1. **量的是什么**：`*perm.Store.Set` 的**运行行为**（真 `Options` ＋ 真 `memConfig`／真时钟桩，同 `store_test.go` 的 `newTestStore` 形状）；并且**另**用 `cmd/wisp` 的产线装配做了一次旁证（§③.1）。
2. **"真腿"的射程**：见 §①.1 的保真度声明。本腿**⛔**声称"原生卡被点过"；它声称的是"**腿被要求问时会不会真的记一次问**"。
3. **`asked` 是自建计数器**，⛔ 读 C18 的队列深度；C18 侧的"卡真的进了队列"由既有 `ANSWER-ALLOW corr=…` 审计行（`queue.go:452`）与 375／C18 用例承重。
4. **行号一律 HEAD blob**：`store.go:268/281/283/297`、`composer_handlers.go:153`、`direction.go:23/30/36`、`run.go:1143/1104-1132/894/867`、`panel_inbound.go:272/277/288/346`、`queue.go:443-453`、`PLAN.md:1374/122/1357-1388/2138`。
5. **本腿新增的判据**：只新增一枚**用例**（`internal/perm/store_371_a1p_test.go`，4 枚 case），⛔ 新增任何产线谓词、⛔ 改任何既有断言、⛔ 弱化任何尺。
6. **票面待更正两处（只记，⛔ 改票）**：①〔丙-a〕第二形"`Switch` 加 `Answered`"在今天的按值签名下**不可实现**（§G.1 第 3 条）；②〔丙-b〕"`Set` 落 `asked=1/0`"要改成"腿发审计行 ＋ `corr` join"（§G.3 风险栏①）。
