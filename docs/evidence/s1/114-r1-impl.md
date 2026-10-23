# 票 114 落地腿 `114-r1` —— **停手上报**（零产码改动）＋ 三组读数

**腿** `114-r1`（写腿／实现者）· **锚** `dd16b005c`（起手＝收手，本腿未提交产码）
**起手树态**（他腿，非本腿）：` M cmd/wisp/run.go`（365-r1）· ` M internal/agent/loop.go` · 收手时另见 `?? cmd/wisp/tasklogstate_365_test.go`
**写面**：本腿只落本文件；`internal/panel/**`、`cmd/**`、`frontend/**`、`internal/perm/**` **零字节改动**。

## §0 停手判语（先结论）

派单两条硬停手条件**同时命中** ⇒ **不硬改**：

| 派单条件 | 本腿读数 | 判 |
|---|---|---|
| ①「若**绕过也拦得住**（门已在原生层）⇒ 停手上报」 | 绕开 composer 直调**下面那一层**（两道深度）**都被拒**，逐字见 §3；只有再往下两层的 `ModeWriter.Set` 不拦（§4 边界） | **命中** |
| ②「若**必须动 `cmd/**`** ⇒ 停手上报」 | 票面唯一未做格 **AC#3**（真机上变宽必须走 C18 卡、超时 300s＝拒）的落点全在 `cmd/wisp/panel_inbound.go:240-309`（两条产线共用的装配函数） | **命中** |

## §1 票面那句逐字

> # 114 — `ParseComposerRequest` 在**生产里零调用者**：面板那条"档位/附件/工作区"的输入通路今天还没接上；接上的那一票**必须把门放在原生侧**（票 92 验收 `R-92-1` / `R-92-2`）

票面对"**原生**"**给了**，且给到判据级（票面 `:38-40` 逐字）：

> **AC#2** 原生侧门（`R-92-2` 的正身）：**mode 写入处理器在注入的 `Confirm`（L2 确认腿）为 nil 时必须拒**，
> 并且**加一条用例钉住"没有 L2 确认腿就不许写档位"**。判据要能被变异证伪：
> 把 nil-检查删掉 ⇒ 该用例必须红；把它改成"只记日志不拒" ⇒ 也必须红。

⇒ **"原生" ＝ Go 侧那一枚 mode 写入处理器**（票面没说要落到"副作用最底层"）。这一格**今天在盘上**：
`internal/panel/composer_handlers.go:133` `if modeIsWidening(from, to) && h.Confirm == nil {` ⇒ 拒，且在任何 `Modes.Set` 之前。

## §2 「没有生产调用者」的读数（现跑）

| 尺（逐字） | 命中 | 非注释·非测试调用者 |
|---|---|---|
| `git grep -nE "ParseComposerRequest" -- '*.go' ':(exclude).scratch' ':(exclude)*_test.go'` | **10 行** | — |
| └ 只取**调用形状** `ParseComposerRequest(`（同过滤） | **2 行** | — |
| └ 逐行判：`internal/panel/bridge.go:134` 定义 ／ `internal/panel/composer_dispatch.go:178` **调用** ／ 余 **8 行全 `//` 注释** | | ★**1（不是 0）** |

**那枚调用者逐字**：`internal/panel/composer_dispatch.go:178` `	req, err := ParseComposerRequest(raw)`（在 `(*ComposerDispatch).Handle` 体内）；
`Handle` 的**产码听众 2 枚**：`cmd/wisp/panel_host_windows.go:948` `	return m.disp.Handle(ctx, raw)`（常驻面板宿主）·
`cmd/wisp/panel_inbound.go:164` `		reply, err := disp.Handle(ctx, raw)`（`wisp panel-inbound`）。
**端到端复现**（不是 grep）：真二进制跑一封**非变宽**封套 ⇒ `rc=0`、逐字 `wisp panel-inbound: 第 1 行已受理，处理器已被调用`
⇒ ★**票面标题前半句「生产里零调用者」＝已过期**（过期时刻＝票 33 落 `composer_dispatch.go` 那笔；台账 `A##` 记 `114-a2` 同一读数）。
本仓 `probes/stale-claim-1/02-ticket-and-ledger-roster.md` 的 `Q01` 早把它列为"害度最大的一枚"；**本程独立复跑，结论一致**。

## §3 「绕过它」那一发（改前形状；改后＝同一读，见 §5）

三发在 `-overlay` 注入的探针里（仓外 `D:\tmp\wisp\20261023\114r1\ov\`，**工作树零突变**）；被测装配＝产线原样
`ModeWriteHandler{Modes: 写腿, Confirm: nil, Audit: 审计, Actor: …}`（`cmd/wisp/panel_inbound.go:260-267`）。

| 深度 | 跳过了什么 | 结果 | 写腿被调用 |
|---|---|---|---|
| 1 | composer **＋**派发表，直调 `HandleModeRequest` | **拒** | **0 次** |
| 2 | 只跳 composer，走产线派发表 `dispatch` | **拒** | **0 次** |
| 3 | 跳到底：直调副作用层 `ModeWriter.Set` | ⚠ **不拒**（`err=<nil>`） | 1 次 |
| 正控 | 同请求同处理器，**只把 `Confirm` 接上** | 过 | 1 次 |

★**「绕过也拦得住」逐字拒绝句（深度 1）**：

```
panel: 确认腿未接入，档位不得变宽: 请求 "bypass-depth1" 要把档位从 ask_every_step 放宽到 ask_high_risk，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）
```

深度 2 同句、只换 `requestId="bypass-depth2"`；**正控**（`Confirm` 接上）`err=<nil>` / `calls=1`
⇒ 拦住它的**正是** `Confirm == nil` 那一跳，不是别的拒绝。
**真二进制那一发（同一句、走完整产线，非白盒）**：

```
wisp panel-inbound: 第 1 行被拒绝：面板请求被拒绝 [panel.mode.request 114r1-a]：panel: 确认腿未接入，档位不得变宽: 请求 "114r1-a" 要把档位从 ask_every_step 放宽到 auto_approve，而本机没有接入 L2 确认腿（R20/M4：只有变宽要确认）
[audit] panel: MODE-REFUSED request="114r1-a" from="ask_every_step" to="auto_approve" origin="panel" actor="cli-panel-inbound" err=… detail="档位未变化：写入腿没有被调用"
```

`rc=1`；`to=ask_high_risk` 同句（`from="ask_every_step" to="ask_high_risk"`）；`to=ask_every_step`（不变宽）`rc=0` 且写 `MODE-SWITCH … result=no-change`；
伪造 `source` 被**更早**一道门拒（`来源 "renderer" 不是 "panel-composer"，按伪造/串台拒绝`）。
**盘面复核**：三发之后 `data/config.toml:124` 仍是 `permission_mode = 'ask_every_step'` ⇒ **没有任何一发真的写进去**。

## §4 「真正的门」该在哪一层 ＋ 为什么

**票面点名那一层＝对的原生层，且已在盘上。** 按主流 harness 口径给三条可核理由（⛔ 只抄票面）：

1. **门要挂在"能改档"那一层的调用侧，不在 UI 侧。** 本仓能改档的只有 `(*perm.Store).Set`（`internal/perm/store.go:175`），
   它的**产线调用者只有 1 枚**：`internal/panel/composer_handlers.go:141`，**就在** `:133` 那道 nil 门之后
   ⇒ 门与副作用之间没有第二条产线岔路。这是"门位置正确"的可核依据，不是自陈。
2. **`perm.Store` 自己那道门只盖一档**（`to == risk.ModeAutoApprove && s.confirm == nil` ⇒ 拒），
   而 R20 的方向不对称要求"任何变宽"都拒 ⇒ 这一层补不了 `ask_every_step → ask_high_risk`；门落**调用侧**是对的，
   也**不该**挪进 `perm`：票面 `Packages` 逐字 `internal/perm/`（存储语义**别改**）。
3. **JS 侧的门是装饰品，本仓已证**（票 92 验收：被攻陷的渲染进程运行期 `postMessage` 不需要源码出现任何被禁字样）
   ⇒ 门必须在 Go 侧、且在副作用之前——今天就在。

★**一条诚实的层位边界（不修，只登记）**：深度 3（直调 `ModeWriter.Set`）**不拦**——`perm.Store.Set` 只对"切进 `auto_approve`"要确认。
今天**不可从页面到达**（`Set` 的产线调用者只有 §4.1 那一枚、且在其后）⇒ **潜伏缝不是活洞**；要堵得落 `internal/perm`（票面禁改）或包 wrapper，本腿不动。

## §5 改了什么（逐行）

**无。** 产码零改动（`git status --porcelain -- internal/panel internal/perm` 起手／收手同值，见 §6 还原核）；不新增用例、不新增产码注释。
理由＝§0 两条同时命中，派单写明此情形"⛔ 硬改"。
★**顺手登记一条产码真伪问题（不修，归派单方裁）**：`internal/panel/composer_handlers.go:36-40` 文件头仍逐字写着
「Nothing calls this handler yet, because the WebView2 "event -> ParseComposerRequest" hop does not exist in this tree」，
而本程已用真二进制证明**它被调用**（§2 那发 `处理器已被调用`）。同族的 `cmd/wisp/run.go:333` 已改准（"hop now exists"），**只有这一处没跟上**。
改它只需删那三句＋换指 `panel_inbound.go:164` / `panel_host_windows.go:948`；但它属"文档句"不属"门"，按停手令不动。

## §6 反向判据（三读数＋还原核）—— 用在**既有**门钉上，证明它不是装饰

票面 AC#2 要求"判据要能被变异证伪"；本腿无新钉，于是把这条证伪**真跑一遍**（改动全在 `-overlay`，仓外副本）：

| 突变 | 突变处（仓外副本逐字） | `rc` | 红名 | 红行逐字 |
|---|---|---|---|---|
| 基线 | — | **0** | 6 枚全绿 | — |
| **M-1** 删掉 nil 检查 | `if modeIsWidening(from, to) && false {` | **1** | `TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached` | `composer_handlers_test.go:137: ask_every_step -> ask_high_risk with no confirm leg: got nil error, the档 would have been written` |
| **M-2** 记日志但照样往下走（删 `return refused`） | `h.record(…)` 之后 | **1** | 同上 | 同上 `:137` |
| **M-2b** 既返回拒绝**又**照样写 | 在 `return refused` 前插 `_ = h.Modes.Set(ctx, to, PanelModeOrigin, h.actor())` | **1** | 同上 | ★`:145: … the injected ModeWriter was called 1 time(s) despite the refusal - the gate did not stop the write` |
| **M-3** 梯子反过来 | `func modeIsWidening(from, to risk.Mode) bool { return to < from }` | **1** | 3 枚 | `:109` `… must read as widening; risk.Mode's constants are ordered strict -> loose on purpose` ／ `:137` 同上 ／ `:223` `ask_high_risk -> ask_every_step (stricter) was refused: …` |

★**M-2b 是要害**：它让处理器**照样返回那句拒绝**（"返回了 error"这条判据**照样绿**），唯一红的是 `:145`「**注入的 `ModeWriter` 被调用了 1 次**」
⇒ 票面 AC#2 写的"判据＝写腿未被调用，不是返回了 error"**今天成立**。
**还原核**：`sha256sum internal/panel/composer_handlers.go` 四发突变前后逐字同值 = `a8620ca4f9d7c5470f958e76692fbb99513c45462661b0dc8ee22a589525ef88`，`RESTORE-CORE: identical`。

## §7 新用例逐名（＋为什么今天响）

**无新增**（§5）。理由：AC#2 的钉子**已存在且有牙**（§6 四发证），再加一枚只会是同一判据的第二种写法。

## §8 门禁 rc（现跑，rc 不走管道）

| 门 | 命令 | `rc` |
|---|---|---|
| 编译 | `go build ./...` | **0** |
| vet | `go vet ./internal/panel/` | **0** |
| 格式 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/panel` | **0**，输出 **0 行** |
| D22 扫描 | `sh scripts/d22scan.sh` | **0**，`clean - no D22 ban violations` |
| 包测试（**甲**：`-v` 顶层口径） | `go test -count=1 -v ./internal/panel/` | **1** ⇒ `RUN 268 / PASS 165 / FAIL 5 / SKIP 0`（`PASS` 含子测 263） |
| 包测试（**乙**：包级） | 同上末行 | `FAIL github.com/CarlosShao/wisp/internal/panel 17.535s` |

d22scan scope 逐数（自陈行逐字）：`bans #1-5 internal/=248`·`bans #1-5 cmd/=42`·`ban #6 frontend/=116`·`ban #7 internal/tools/=25`·
`ban #8 design/=39`·`ban #8 frontend/=116`·`ban #8 internal/=635`·`ban #8 cmd/=149`·`examined 290 production Go files`。
**5 枚红逐名**（与票面历史 `3 枚前端耦合 + 2 枚 35-r8` **同集**、本波零新增）：`TestApprovalCardViewJSONKeysMatchFrontendTypes` ·
`TestComposerContractTypesMatchFrontend` · `TestC21DesignTokensFourWayAgree` · `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` ·
`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`。本腿 `internal/panel` 零改动 ⇒ 这 5 枚的红绿均非本腿造成（判据＝§5 porcelain 同值）。

## §9 未定义即停：请派单方拍三件（⛔ 本腿自行假设）

1. **AC#3 要不要开 `cmd/**`**：`cmd/wisp/panel_inbound.go:249-267`（`newComposerDispatchChain`，**两条产线共用**）两处 `Confirm: nil`；
   而 `cmd/wisp/run.go:687-694` **已有**一枚接好 L2 腿的 `rt.modeWrites`，`run.go:333` 自陈"Nothing reads this field yet"。
   ⇒ 把那一枚交给这条链＝"一处改动、两条产线同时变宽可走卡"的最小形状；但全在 `cmd/**`，且 `365-r1` 正在改 `cmd/wisp`。
2. **票面 `:44` / 标题 / `:21` 的过期句要不要改**：本腿**未动票面**（⛔ 改票名／`Status:`），但读票人会读到"通路不存在"（§2 已证过期）。
3. **深度 3 的潜伏缝要不要单独立票**（§4 边界）：`perm.Store.Set` 对非 `auto_approve` 的变宽不设门。

## §10 纪律自报（⛔ 藏）

- 为取"一份全默认 `config.toml`"跑过一次 `wisp run`（它不认 `-data`，落到 `%APPDATA%\wisp-dev\` 并起过 WebView2 面板）；
  **只杀了本程起的进程（PID 81876）**，他腿的 `wisp`（PID 21892、`build\wisp.exe`、10-15 起）**未动**。此后改用 `-overlay` 种子＋`panel-inbound`，未再起窗。
- 起手 `git diff --cached --name-status` **空**；写面（`internal/panel/`、`cmd/`、`internal/perm/`）本腿**零写入**；两枚 ` M` 是他腿的，本腿未碰其行为。
- `unset` 只用 `GOTMPDIR`；`cmd/wisp` 相关命令一律带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`
  （首跑 `rc=127` `sherpa-onnx-c-api.dll` ＝**假失败**，已重跑取真读数）。
- 临时件只在 `D:\tmp\wisp\20261023\114r1\` 下建、未删；读数 `.out` 全在该仓外目录，仓内无 `.out`。零 `push`；本腿只 commit 本文件（pathspec 字面点名）。
