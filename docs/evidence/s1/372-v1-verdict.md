# 372-v1 终裁：回退（件一）＋读侧降级（件二）**够格**，`AC#2` 能翻

- 腿 `372-v1`（非实现者＝裁腿）。量于 2026-10-24 11:37–12:0x +08，HEAD `9965e86a9`；被裁对象 `b4ee24cbe`（11 文件 `+746/−247`）。
- 判决依据（唯一权威）＝`docs/evidence/s1/m4-survey.md`（108 行，见 ① 的 blob 复量）。前腿 `372-r1-blocked.md`；件 `372-r2-impl.md`（109 行）。
- **写面**：本件一条路径；其余全部在仓外 `D:\tmp\wisp\20261024\372v1\**`（4 支 `-overlay` 探针 ＋ 4 枚变异件 ＋ 3 份旧 blob ＋ 日志）。**工作树 0 突变**：`git status --porcelain internal/perm internal/panel internal/tools cmd/wisp` 空、`git diff --cached` 空、`git diff HEAD -- <11 笔路径>` 空。
- 量具（可逐条重跑）：`logs/10-perm-probe`（六对/两发/panic/记账）、`11b-panel-probe`（handler）、`12-tools-probe`（A–F＋Z）、`13c-cli-probe`（CLI 管＋审计行）、`14-before-pairs`（改前 blob）、`21/22/23-step*`（三步法）、`31–38-m*`（变异）、`40–49-*`（门禁）。

## ① 票面前提新鲜度（对 HEAD blob 复量）＋ r1 那三条

| 前提 | 复量 | 判 |
|---|---|---|
| `m4-survey.md` 108 行 @`11ede6ef5` | `wc -l`=**108**；`git hash-object` 与 `git rev-parse 11ede6ef5:…` **同为 `1eebdf9a6804989241a2fa923e4ba898a8a8f45c`** | ★**新鲜** |
| `372-r2-impl.md` 109 行 | **109** | 新鲜 |
| `372-r1-blocked.md` **119 行** | **163 行**（唯一 blob `5f02c04bb` 即 163） | ★**派单那个数错 44**；内容与派单描述一致（§①–§⑥、甲/乙/丙），不影响裁 |
| 票面「`AE→AH` 今天被拒」 | 同一支探针挂 `255d42c75` 的 `store.go`+`direction.go`：**3 变宽全拦**（`AE→AH`/`AE→AA`/`AH→AA`，拒绝句点名 from/to、`saves=0`）；3 变窄全过 | ★**前提成立**（`14-before-pairs`） |
| 票面 `Status:` | 仍 `**ready-for-agent**`（未翻） | 合规 |

**核 r1 三条**：①§③「回退 ⇒ 2 枚 cmd 测试红在 `:106/:109/:160/:163`」⇒ 我三步法第 ② 步**逐行命中**；②§④「动作层缺口＝`sil.Silenced` 之后 `gate` 一次不被调」⇒ 我的 A/C/Z 三行复现（`MODE-SILENCE=1` 且 `L1_windows=0`）；③§④末段「读侧降级形：`perm.Store` 是全进程唯一同时知道『配的哪档』＋『有没有腿』的对象，且不动 `tools`/`cmd`；`Set` 的 from/to 必须读配置值」⇒ ACCOUNTING 行 ＋ M3 变红复现。**三条全立。**

## ② 核心判点一：六对表（★我自己造的，`Confirm=nil`）

| from→to | 变宽 | set_err | applied | config_after | **有效档读回** | saves | refusals |
|---|---|---|---|---|---|---|---|
| `AE→AH` | 1 | `<nil>` | 1 | `ask_high_risk` | **`ask_every_step`** | 1 | 0 |
| `AE→AA` | 1 | 拒（点名 AE/AA） | 0 | `ask_every_step` | `ask_every_step` | 0 | 1 |
| `AH→AE` | 0 | `<nil>` | 1 | `ask_every_step` | `ask_every_step` | 1 | 0 |
| `AH→AA` | 1 | 拒（点名 AH/AA） | 0 | `ask_high_risk` | `ask_every_step` | 0 | 1 |
| `AA→AE` | 0 | `<nil>` | 1 | `ask_every_step` | `ask_every_step` | 1 | 0 |
| `AA→AH` | 0 | `<nil>` | 1 | `ask_high_risk` | `ask_every_step` | 1 | 0 |

**总数核对**：3 变宽 → **2 拦 1 过**（只有目标 `auto_approve` 拦）；3 变窄 → 0 拦。**正控**（真 `Confirm`）：`AE→AH` ⇒ `applied=1`、**`confirm_calls=0`**；`AH→AA` ⇒ `applied=1`、**`confirm_calls=1`**。判语：**与 r2 报的六对表逐格一致**，`AE→AH` 不再是门，M4 的不对称未动。

## ③ 核心判点二：五发 A–E ＋ `M1` 对照（`ModeSource`＝真 `perm.Store`，同一发 `fs.write`）

| 发 | config／腿 | 有效档 | is_error | 落盘 | L1 窗 | L2 卡 | MODE-SILENCE |
|---|---|---|---|---|---|---|---|
| A | `auto_approve`／none | `ask_every_step` | **1** | **0** | 1 | 0 | **0** |
| B | `ask_every_step`／none | `ask_every_step` | 1 | 0 | 1 | 0 | 0 |
| C | `ask_high_risk`／none | `ask_every_step` | **1** | **0** | 1 | 0 | **0** |
| D | `auto_approve`／none，越白名单 | `ask_every_step` | 1 | 0 | 0 | 1 | 0 |
| E | `ask_every_step`／wired | `ask_every_step` | 0 | 1 | 1 | 0 | 0 |
| **M1 对照**（摘降级）：A／C | 同 | **`auto_approve`／`ask_high_risk`** | **0** | **1** | 0 | 0 | **1** |
| **Z 对照**（绕过读侧，直接给桥 config 值） | `auto_approve` | — | 0 | 1 | 0 | 0 | **1** |

★**A/C 两格从「静默落盘」变「拒」**（文本＝「L1 确认窗口尚未接入（票 21），已拒绝执行」），B/D 仍拒，E 仍问一发过一发 ⇒ **件二的牙口我复现了**；M1 与 Z 从反面把「静默」这一格摆出来。判语：**这是本腿的真价值，读数成立。**

## ④ 核心判点三：「只认 nil-ness」对不对（造 `F`＝`AA`＋**非 nil** 腿）

`F`：`config=auto_approve`、腿 wired、`NoGate` ⇒ **`effective=auto_approve`、`file_written=1`、`L1_windows=0`、`MODE-SILENCE=1`**（对照 M1 的 A 行同形，但这里腿**真在**）⇒ ★**降级只认 nil-ness、⛔ 按档位值乱降**：成立。

**判这对不对**：★**对本层是对的，而且就是主流判法**——perm 拿不到「可问性」，它手里只有「有没有腿」这一个信号；(a) dsh 的判据是 `approval === undefined`（**presence**），坏/抛的 answerer 由**动作层**降成 `unavailable`（`m4-survey:9`）；(b) openchamber 的 `displayedPermissionMode` 就是「这档赖以成立的东西不在 ⇒ 按更严的显示/跑」（`:88`）；(c) Step-Code `effectiveState` 的轴是 `cause.refused||unconfigured`（`:56-60`）。⇒ **「有腿就算有通道」在档位层是主流做法**，前提是**腿的行为在动作层 fail-closed**（本仓＝`internal/tools/gate.go:146-156` 的 `NoGate` ⇒ reject、`:126-140` 的 `GateFuncs` nil ⇒ `AnswerReject`，r1 §④ 已核）。
★**但这条代理在 `wisp run` 一处不成立**（r2 残余①）：那里的非 nil 值是**切档卡的路线**（`confirmModeSwitchVia` ⇒ `gate.PendingApproval`），其自述写明控制台跑没有原生面、卡只能超时成拒（`run.go:883-887`）⇒ **「有腿」≠「问得出来」**。⛔ 把 `F` 读成「`AA` 安全」；它是「`AA` 有一个**自称能问**的腿」。

## ⑤ 核心判点四：`PermissionMode()` 两发对照

| config.toml | 腿 | `PermissionMode()` | 配置被改写？ | `MODE-DOWNGRADE` |
|---|---|---|---|---|
| `ask_every_step` | none | `ask_every_step` | 否（`saves=0`） | 0 |
| `ask_high_risk` | none | **`ask_every_step`** | **否（`saves=0`，`config_file=ask_high_risk`）** | **1** |
| `auto_approve` | none | **`ask_every_step`** | **否（`saves=0`，`config_file=auto_approve`）** | **1** |
| `ask_high_risk`／`auto_approve` | wired | 原档 | 否 | 0 |

★配置**一字不改**、读回降级档、每格恰 1 行 `MODE-DOWNGRADE`（字段含 `from=/to=ask_every_step/origin="startup"/result=downgraded-no-channel`，且**不自称** `MODE-SWITCH`）；有腿 ⇒ 原档、0 行。判语：**成立**。附带：`Snapshot().Mode` 也随之报有效档，而 `run.go:724` 把面板的档位视图接在 `rt.modes.PermissionMode` 上 ⇒ **页面显示的是"真在跑"的档**（openchamber 形状），不是缺陷。

## ⑥ 核心判点五：那 2 枚 `cmd` 测试的期望改动

**三步法（我重跑）**：①`255d42c75` 码＋旧测 ⇒ **rc0**（5 枚 `TestTicket114r2*` 全绿）；②本笔码＋旧测（`-overlay` 挂回）⇒ **rc1**，逐行 **`:106`／`:109`／`:160`／`:163`**（与 r1/r2 报的字面一致）；③本笔码＋新期望 ⇒ **rc0**。
★**「`AA` 那半 0 改动」的核（`git show -U0` 的 hunk 为准，⛔ 只看绿）**：
- **resident 链**：hunk 只落在旧 `98`／`105-106`／`108-109`（AH 半）；**AA 那半（旧 `132-142`：一张卡、`leg.asked()==[auto_approve]`、disk `auto_approve`）不在任何 hunk 内 ⇒ 0 改动**；
- **CLI 管**：旧 `152-166`（含 `for _, to := range {ask_high_risk, auto_approve}` 整段）落在 `@@ -152,15 +171,39 @@` ⇒ **被展开成两个子测试**，`AA` 的两条断言（`errors.Is(ErrNoL2Confirm)`＋disk 不动）**在、行不在**（措辞微调，判据未松）。
- ★**新加的有效档读数是真牙**：M1（摘降级）⇒ `cmd` 红在 **`:207`**（「the pipe's disk write widened the EFFECTIVE档 to ask_high_risk」）。无 `Skip`、无改宽错误串。
判语：**改期望这一步做得够格**；但「`AA` 原断言 0 改动」这句**只对 resident 那半成立**，对 CLI 管是「语义原样、行重排」，⛔ 当 0-diff 读。

## ⑦ panic 那发（`recover()` 在侧）

from `AE` 与 from `AH`、`confirm=nil`、`Set(auto_approve)` ⇒ **`panicked=0`**（两形）、`last=refused-confirm`、`config_after` 原值、`saves=0`、有效档 `ask_every_step`。**M4（删整段目标门）⇒ 真 panic**：`invalid memory address or nil pointer dereference`，栈顶 `internal/perm/store.go:277`（变异件行号＝`s.confirm(ctx, sw)`）⇒ ★`store.go` 那条 `There is no nil test here on purpose` **仍承重**。

## ⑧ 核心判点：四变异（全仓外 `-overlay`，工作树零突变）

| 发 | 摘掉/翻回哪一跳 | rc | 逐名红行（我的） |
|---|---|---|---|
| `M1` | `PermissionMode` 去掉降级 | perm **1**／cmd **1** | `store_367_test.go:128`×3；`store_372_test.go:36`×2、`:63`、`:103`；**cmd `:207`** |
| `M2` | **两门**翻回 `IsWidening` | perm **1**／panel **1**／cmd **1** | `store_367_test.go:112/115/118/122`；`composer_gate_agreement_367_test.go:86/90/121/124`、`composer_handlers_test.go:180`；cmd `:119/:122`、`:193/:196` |
| `M3` | `from := s.PermissionMode()` | perm **1** | `store_367_test.go:103/115/118`；`store_372_test.go:113/116/119` |
| `M4` | 删整段目标门 | perm **1** | `TestSetGatesOnly…` 触发 **nil deref panic**（`store.go:277`）⇒ 包级 FAIL（panic 先于 `store_372_test.go` 那枚跑到） |

★四发**全红**；**还原核在最后一发（M2 cmd）之后**：见 ⑩（真树四数＋四包 `git status` 空）。★r2 报的 M1/M2/M3 红行与我逐行一致；M4 它报 `store_372_test.go:134`，我实测 panic 先生效在 `store_367_test.go` 的用例上——**同一件事、不同的先撞者**，⛔ 当它报错。

## ⑨ 边界（逐项）

- `b4ee24cbe` 恰好 **11 路径**（10 代码＋本笔件；`--numstat` 合计 `+746/−247`，与我逐文件相加一致）。
- **`internal/tools/**` = 0 路径**；`retention.go`／`schema.go`／`ddlV*`／任何 `*roster*` = **0 路径**（`git show --name-only | grep` 实测）。
- **`cmd/wisp` 只 2 个文件**：`panel_mode_confirm_114_windows_test.go`（那 2 枚期望）＋`panel_inbound.go`（hunk `@@ -341,4 +341,5 @@`，**改动行 100% 是 `//`**：非注释改动行 = **0**；即 `:341` 起那段注释 4 行改写成 5 行，⛔ 不是"一句"）。
- `b4ee24cbe^..HEAD` 多出的只有 `docs/reports/pending-and-issues.md`（编排者 A1220 的台账），与本笔无关。

## ⑩ 门禁复跑（真树、无 overlay；四数分开、既有红分开）

| 项 | 读数 |
|---|---|
| `go build ./...` | **rc0** |
| `go vet ./internal/perm ./internal/panel ./internal/tools ./cmd/wisp` | **rc0** |
| `./internal/perm` | **rc0，PASS 23／FAIL 0**（基线 19/0，+4 新） |
| `./internal/panel` | **rc1，PASS 166／FAIL 5**——既有 5 枚逐名：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree` ⇒ **零新增** |
| `./internal/tools` | **rc0，PASS 250／FAIL 0** |
| `./cmd/wisp` 定向（`PATH` 带 `third_party/sherpa-onnx`） | **rc0，PASS 29／FAIL 0**（整包会死在 `MARK299`，只跑定向） |
| `sh scripts/d22scan.sh`（**含正控**） | **rc0**，`d22scan: clean` |
| 格式甲（`gofumpt -l` 10 文件） | **空**；`gofumpt -d` **0 行** |
| 格式乙（`gofmt -l` 10 文件） | **空（0 行）** |

尺：PASS/FAIL＝`-v` 下 `^--- PASS:`／`^--- FAIL:` 计数（含子测试）；格式两把尺都量 10 个改动文件（`gofumpt -l` 空≠HEAD 行尾，⛔ 用 `gofumpt -l` 判 CRLF）。

## ⑪ 我那三条裁，逐条判

1. **`run.go:668-671` 非 nil 腿＝设计不是缺陷；真问题是「有腿≠问得出来」⇒ 374** —— ★**结论对**。补一刀：那里非 nil 的是**切档卡**的路线，与 L1 窗的表面本就不是一回事，所以「有通道所以不该弹卡」这半（M4）成立、"有通道"这半**不成立**；修法也只能在 seam（无可答表面⇒传 nil，或给 `perm` 一个显式「可问性」信号），⛔ 让 perm 猜。另：同一代理在**显示**侧是诚实的（`run.go:724`⇒页面显示有效档），无需动。
2. **CLI 管的 disk 可以改** —— ★**对**。我的 CLI 读数：`to=ask_high_risk ⇒ handle_err=<nil>、disk_after=ask_high_risk、effective_after=ask_every_step`；`to=auto_approve ⇒ ErrNoL2Confirm、disk 不动`；且放宽那一笔记了耐久账 **`perm: MODE-SWITCH from=ask_every_step to=ask_high_risk … result=applied`**（拒的那笔记 `panel: MODE-REFUSED`）⇒ 「档位照改 ＋ 动作层 fail-closed」正是主流（`m4-survey:96/104`）。**一个具名小口子**：降级记录只在**构造时**落一行，进程活着的期间把盘改宽**不会**当场再落一行 `MODE-DOWNGRADE`（要等下次启动）——安全不受影响（有效档立刻就是 AE），但「降级发生」这件事当时只有 `MODE-SWITCH` 一条侧证。
3. **`askOnce` 夹具补确认腿、断言语义未改** —— ★**成立且必要**：diff 只动夹具（`New(Options{Manager:mgr})` 加 `Confirm`），断言行**一字未动**；不补则「改配置键⇒链路行为变」在降级后不成立。代价：该用例不再覆盖无腿链路——但那件事另有牙（M1 在 perm 与 cmd `:207` 都红），⛔ 算放松。

## ⑫ `AC#2` 能翻／⛔翻 ＋ 缺什么

★**能翻**：四发必给读数齐（六对／动作层 A/C／正控／panic）、每处改动配了今天响的用例、反判据四变异全红并逐名报红行、还原核在最后一发之后。**翻框时必须在票面写清两条射程**：
1. **件二只在 `Confirm==nil` 的装配里生效**：`wisp run`（`run.go:668-671`）永远非 nil ⇒ 该 seam 的 A/C 两格**与改前同**（不是本笔弄坏的，但本笔也没盖住）⇒ **374 不落，⛔ 声称「全自动那条路已安全」**。
2. **「动作层那发」的实现层是 `perm.Store.PermissionMode()` 的读侧**（`tools/mode.go:27-31` 每发读一次），不是 `internal/tools` 里的一道门；效果等价（A/C 被 `NoGate` 拒），但 AC#2 第 ② 条的措辞（"动作层"）与本笔的落点差一层——⛔ 当它写了 `internal/tools` 的门。

## ⑬ 提交号 ＋ 逐笔越界

- 本件一笔：`docs/evidence/s1/372-v1-verdict.md`（唯一写面）。★落点前 `ls`：该路径与 `372-v1-verdict-part*` 均**不存在** ⇒ 用 `372-v1-verdict.md`，未改号。
- 越界：**0**。⛔ 未碰 `internal/**`、`cmd/**`、`scripts/**`、`.github/**`、`.scratch/**`（1263 行别人的脏一处未碰）、`docs/reports/**`、`docs/PLAN.md`、`docs/specs/**`；⛔ 未落任何 `.go`（探针/变异件全在 `D:\tmp\wisp\20261024\372v1\**`）；⛔ push、⛔ 翻框、⛔ 改票名/`Status:`/`-done`、⛔ 台账。
- 索引核对：`git diff --cached --name-status` 提交前**仅 1 行 `A`**（本件），无 `D`/多余 `A`。
