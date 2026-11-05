# R9 阻断报告：扩 M4（371-〔甲〕）的真实写面越过本波禁区（cmd/wisp ＋ internal/command）

- 腿：`R9`（实现腿）· `HEAD 557096c82ac40b23a769b93a79e572215a3ccab5` · 2026-11-05 14:4x–15:0x +08
- 结论：**STOP，不写产码。** 一句话理由：**R9 一落地，`cmd/wisp` 里 7 个断言必红**
  （它们逐字钉着旧 M4「只有切进 `auto_approve` 才出卡」），且同一谓词在
  `internal/command/roster.go` 还有一个生产调用点（371 写面外）⇒ 本波「no cmd/wisp budget」的禁区先于实现。
- 本腿唯一写入＝本文件。零产码、零测试改动、零 `cmd/wisp` 字节。

---

## 0 输入（STEP0，现跑）

| 项 | 读数 |
|---|---|
| `.scratch/wisp/issues/R9-*.md` | **不存在**（`Get-ChildItem .scratch/wisp/issues -Filter 'R9*'` 空读数）。R9 不是独立票面，是 **371 的残余项**。 |
| R9 的登记处 | `371-iswidening-...-done.md:242`：`R9`（待人拍板，⛔ 本票自决）＝〔甲〕扩 M4 让所有变宽出卡 ⇒ 须摆机主；另见 `:7` Status 行「★R9 扩 M4＝待机主、已列入摆机主清单」 |
| 机主授权 | 台账 `A1541`（`docs/reports/pending-and-issues.md:31174`）逐字「按推荐 → **R9 全弹卡（所有变松动作都出确认卡）批准**」 |
| 派单 | 台账 `A1542`（`:31187`）「机主上轮已批 51-AC#3/347/R9 ⇒ 本波起派 51-B/339/R9」 |
| wave9 派单原文 | `.scratch/wisp/probes/orch/20261105/txt/wave9-182-dispatch.txt:4` 逐字「R9 (自定面)：机主批"所有变松动作都出确认卡"，Go 侧把 loosening 都引到卡/审批队列；**STOP 若需 frontend/risk冻结/cmd-wisp**」；`:5`「本波不分配 cmd/wisp 写腿（**若面落 cmd/wisp ⇒ STOP**）」 |
| 371 票状态 | `Status: done`（2026-11-05，A1529：R1–R8 DONE、**R9 surfaced**）。按「withdrawn/done ⇒ STOP」的字面本可整腿停工；因 R9 已被机主批、本波点名派出，本腿按「371 之外的独立工作项」处理，**停因不是 done，是写面越界（见 §2）**。 |

## 1 三条核（THREE-CORE）

### ① 产品定案（R9 要什么）

| 出处 | 逐字/要点 |
|---|---|
| `docs/PLAN.md:132`（D4） | 「用户侧只能**调严**；调松需显式勾选风险告知（不可静默降级）」 |
| 机主 2026-09-21 | 「只有切进全自动要确认」（M4／R20 定案，`371:37`、`store_367_test.go:141-143` 引用） |
| 台账 `A1541` 第 4 条 | R9 全弹卡**批准**（2026-11-05）⇒ 把「只有 `auto_approve` 要卡」扩成「**每一次调松都要卡**」 |
| `docs/evidence/s1/m4-survey.md:96` | 主流做法不按「方向」门控（6 个 repo 无先例）⇒ 这是**机主自选的产品决定**，不是主流对齐；票面已如实记（`371:242`） |

### ② 现状代码：门判的是「目标档」，不是「方向」

| 位置 | 现读 |
|---|---|
| `internal/perm/direction.go:23` | `func RequiresCard(to risk.Mode) bool { return to == risk.ModeAutoApprove }` |
| `internal/perm/direction.go:57-65` | `JudgeCard(to, channel bool)`：`GateNoCard / GateNoChannel / GateAsk`（唯一判决；生产调用点实为三处：`perm`／`panel`／`roster`） |
| `internal/perm/store.go:348` | `switch JudgeCard(to, s.confirm != nil)`；`GateAsk` 分支要「**被回答过的** `Card`」（371）＋ R6 内层 nil 兜底（`:370-376`） |
| `internal/panel/composer_handlers.go:176` | `if perm.JudgeCard(to, h.Confirm != nil) == perm.GateNoChannel`（面板半） |
| `internal/command/roster.go:287` | `switch perm.JudgeCard(to, env.CardChannel)`（**第三处**生产调用点，回答「这条命令要不要卡」） |
| `internal/risk/mode.go:50-55` | `ModeAskEveryStep < ModeAskHighRisk < ModeAutoApprove`（「Do not reorder」）⇒ `IsWidening`（`direction.go:36`）仍可作方向判据 |

### ③ 跨通道不一致（R9 真正要治的东西）

同一个「每步都问 → 只问高危」，**两条通道今天判决相反**：

| 通道 | 现读 | 出卡？ |
|---|---|---|
| `config.toml` 手改（D36 rule 1） | `internal/config/manager.go:472-485`：`newMode > oldMode` ⇒ `risk.permission_mode` 记 loosening；`:203-215` 逐 section 调 `ConfirmLocked`，没答＝拒绝（`不得静默生效`） | **出**（每一次调松） |
| 运行时切档（`perm.Store.Set`，panel / CLI / ball 共用） | `direction.go:23`：只有 `auto_approve`；`AE→AH` 零卡、直接落盘 | **不出** |

⇒ R9 ＝ 把运行时通道对齐到配置通道与 D4:132，即「**所有调松动作都出确认卡**」。
本腿确认这条需求本身清楚、可实现，**卡在写面**（§2）。

### ④ Blocked-by

- **无外部依赖**：不需要前端（卡是 C18/`approval.Gate` 既有卡片，`cmd/wisp/run.go:1148-1156` 已为 `permission.mode` 造卡，不需要新渲染）· 不需要 `internal/risk`（只读 `Mode` 序）· 不需要真机（测试用 stub leg）。
- **唯一阻塞＝写面预算**：`cmd/wisp/**` 与 `internal/command`（见 §2）。

## 2 真实写面（blast radius，现读）

R9 不可能只在 `internal/perm` 落地：谓词要从「只要目标」变成「要 from＋to」，四个生产调用点同判，否则 371 建立的
「一处定义、多处同判」当场破功，且 `roster` 会对用户说反话。

| 层 | 文件:行 | 性质 | 本波可否写 |
|---|---|---|---|
| 生产 | `internal/perm/direction.go:23,36,57-65` | 谓词签名＋判据 | 可（371 面） |
| 生产 | `internal/perm/store.go:325-406`（门/文案/注释）、`:181-197,299-324` | `Set` 门 | 可（371 面） |
| 生产 | `internal/panel/composer_handlers.go:176-181` | 面板半＋拒绝文案 | 可（371 面） |
| 生产 | **`internal/command/roster.go:270-298`** | 同一谓词第三调用点＋三条解释文案（`:290` 逐字「R20/M4 只有切进全自动要确认…写腿直接落盘」） | **面外**（371 未含 `internal/command`） |
| 生产注释 | `cmd/wisp/panel_inbound.go:275,345-349` | 逐字「refuse every auto_approve switch」「gate is perm.JudgeCard(to, …)」 | **禁写** |
| 测试 | `internal/perm/store_367_test.go:134-136,157-159`（「only the two auto_approve targets are gated」／「non-auto widening cost 0 confirmations」）、`:144-176`；`store_371_a1p_test.go:105-122`；`store_371_r1p_test.go:236-245`；`store_372_test.go`；`store_test.go:143-155,240` | 旧 M4 的正/反钉 | 可（聚焦 `_test.go`），但**必须随行为翻转**（不是放宽） |
| 测试 | `internal/panel/composer_gate_agreement_367_test.go:61,87-91`、`composer_gate_agreement_371_test.go:155,340-364`、`composer_handlers_test.go:158` | `RequiresCard(to)` 编译面＋两处一致格 | 可，同上 |
| 测试 | `internal/command/roster_test.go:314` 及文案断言 | 第三处同判 | 面外 |
| **测试（禁写、必红）** | **`cmd/wisp/run_mode101_test.go:238-244,253-257`** | `want exactly [auto_approve]` | **禁写** |
| **测试（禁写、必红）** | **`cmd/wisp/panel_mode_confirm_114_windows_test.go:121-126、137-146、150-158、189-199、232-240、342-357`** | 6 格（无腿落地／有腿零卡／`[auto_approve]`／CLI 落地／假腿静默／run 缝 `[auto_approve]`）＋ 4 处注释 | **禁写** |

### 这 7 格为什么必红（不靠猜，靠调用链）

`perm.Store.Set` 现在的 `asked` 只在 `GateAsk` 分支里 +1（`store.go:368-395` 调 `s.confirm`）。
R9 把 `AE→AH` 也判成 `GateAsk` ⇒ 同一串 `Set` 调用：

- `run_mode101_test.go` boot1：`Set(ask_high_risk)` 先问 1 张、`Set(auto_approve)` 再问 1 张 ⇒
  `gotAsked=[ask_high_risk, auto_approve]`，`:256-257` 的 `want exactly [auto_approve]` **红**。
- `panel_mode_confirm_114`：`:140-143` 的 `leg.calls()==0` **红**；`:153-155` 的 `[auto_approve]` **红**；
  `:121-126`／`:189-199`／`:232-240` 的「无腿/假腿也落地」变成 `GateNoChannel` 拒绝 **红**；`:355-357` **红**。

**今天它们是绿的**（现跑，非推断）：

```
$ export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
$ go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart|TestTicket114r2' -count=1 -timeout 300s
ok  github.com/CarlosShao/wisp/cmd/wisp  45.077s
```

⇒ 本腿若落 R9，会在 `cmd/wisp` 留下 7 个**我无权修**的红格；这不是「已有红」，是**我造成的红**。
票面自己也是这么写的：`371:280` 逐字「必须改 `run_mode101_test.go:251` 那条断言（★那是**改产品决定**、⛔ 改注释）」。
⇒ 命中派单 STOP 条件 **`face=cmd-wisp`**（外加 `internal/command` 生产面越界）。

## 3 为什么不交「半片」

1. **roster 会说反话**：`internal/command/roster.go:287-290` 拿同一谓词回答「这条命令要不要卡」，
   只改 `perm`/`panel` 而不改它 ⇒ `wisp commands` 继续打印「只有切进全自动要确认…写腿直接落盘」，
   而实际 `AE→AH` 会弹卡 ⇒ 用户可见的说明书与行为相反。它是生产文件、在 371 写面外。
2. **371 的价值是「一处定义、多处同判」**：`371:174-176`（AC#2'-b）要求 `perm` 与 `panel` 同判且「只改一处必须红」。
   R9 把调用点变三处（perm／panel／roster），**必须同批**，否则那条「两处一致」的牙从设计上就漏了第三处。
3. **半片的绿是假的**：只在 `perm`/`panel` 内绿、`cmd/wisp` 与 `roster` 红/说谎，等于把「产品决定已落地」写成假读数。

## 4 下一波落地草案（需 1 枚 `cmd/wisp` 预算＋`internal/command` 面）

> 只是草案，**本腿未改一字节**；给有预算的腿照着做。

```go
// internal/perm/direction.go
func RequiresCard(from, to risk.Mode) bool { return IsWidening(from, to) } // 调松（含中间档）都要卡
func JudgeCard(from, to risk.Mode, channel bool) CardGate { /* 同三态，判据换 RequiresCard(from,to) */ }
```

- `store.go:348` → `switch JudgeCard(from, to, s.confirm != nil)`；`GateAsk` 内文案把「切到 auto_approve」改成「从 X 切到 Y」；
  **读侧 372-r2 不动**（无腿时 `PermissionMode` 仍塌到 `ask_every_step`）；**371 的 answered-card 与 R6 nil 兜底原样保留**。
- `composer_handlers.go:176` → 同传 `from`；拒绝文案去掉「只有全自动要确认」。
- `roster.go:287` → `perm.JudgeCard(env.Mode, to, env.CardChannel)`；三条文案同步。
- 要翻的格：`perm` 5 文件（`store_367_test.go:134-176` 是主格）、`panel` 3 文件、`roster_test.go`、`cmd/wisp` 2 文件 7 格。
- 新增牙（红→绿）：`AE→AH` ＋ 无腿 ⇒ `ResultRefused`；`AE→AH` ＋ 真腿 ⇒ `Asked=true`、`confirm_calls=1`、`applied=1`；
  `AH→AA` ＋ 真腿 ⇒ 1 张；`AA→AH`（调严）⇒ **0 张**；反向变异＝把 `RequiresCard` 退回 `to==auto_approve` ⇒ 新格全红。
- 已知代价（写进票面）：**CLI 无卡面 ⇒ 从此不能调松**（`AE→AH` 也拒，fail-closed）；配置手改通道仍可调松且出卡（D36）。
  D4:132 支持这个方向，但这是用户可见 UX 变化，须具名。

## 5 门禁读数（HEAD，本腿零产码）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `GOFLAGS= go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| test（面内） | `go test ./internal/perm ./internal/panel ./internal/command -count=1 -timeout 240s` | `internal/perm` **ok 0.639s** · `internal/command` **ok 0.338s** · `internal/panel` **FAIL**：2 枚**既有**红（`TestComposerContractTypesMatchFrontend`、`TestC21DesignTokensFourWayAgree`，前端/C21 侧，HEAD 干净树同上，非本腿） |
| test（禁区基线） | `go test ./cmd/wisp -run 'TestTicket101ManualSwitchSurvivesRestart|TestTicket114r2' -count=1` | **ok 45.077s**（⇒ 旧 M4 的 7 格今天是绿的，§2 的「必红」由此成立） |
| gofumpt | `gofumpt -l cmd internal tools` | 仅 `internal/agent/approval/pending_read.go`（既有 CRLF 131 对、`git status` 干净、非本腿；A1540/A1542 已具名披露） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**（ban#8 internal/=884、cmd/=207、frontend/=116、design/=39；clean） |

## 6 本腿写入

- `docs/evidence/s1/R9-stop-cmdwisp-face.md`（本文件）。零产码、零 `_test.go`、零 `cmd/wisp`、零 `internal/command`。
- 未翻任何票面 `[ ]`／`Status:`／`-done`；未改 `M4` 任何断言；未 push。
- 共享工作树的并行读数（**非本腿**）：本腿运行期间（14:51 前后）同浪 `339-r2` 腿在改
  `internal/risk/anchor_miss_reason_334_test.go`、`internal/risk/anchor_miss_reason_339_other_test.go`（去掉 windows build tag／加 339 格）；
  我的提交用了显式 pathspec，只含本文件。

## 7 给编排者的一句话

R9 的需求与授权都成立、实现形状也已定（§4），**唯一缺的是写面**：请给一条 `cmd/wisp` 预算的写腿
（顺带 `internal/command/roster.go`），我可按 §4 一次落全；若坚持本波零 `cmd/wisp`，则 R9 只能继续挂起，
不能以「perm 内绿」冒充落地。
