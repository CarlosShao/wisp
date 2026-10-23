# 367-v1 终裁腿 —— `367-r1` 的 `AC#2`（非 `auto` 变宽的门：落地＋牙）能不能翻

- 腿 `367-v1`（**非实现者裁腿**，⛔ 实现者）。HEAD **`7fabccf1b`**，量于 2026-10-23 18:5x–19:1x +08。
- 写面：本文件 ＋ 仓外 `D:\tmp\wisp\20261023\367v1\**`。**仓内 `.go` 改动 0 字节**（`git status --porcelain internal/perm internal/panel cmd/wisp internal/risk` 空、`git diff --stat` 同四包空）。
- 尺：**我自己另写**的三支探针（`perm`／`panel`／`cmd/wisp` 各一支），真 `config.Manager` ＋ 真 `config.toml` ＋ 真 `perm.Store`。
  `before` ＝ **只换** `git show 7fabccf1b^:internal/perm/store.go` 这一枚 blob（`cmd/wisp` 那一发再换旧 `composer_handlers.go`）⇒ **单变量 A/B**，⛔ 删 `direction.go`（删了会被已入库的 `store_367_test.go` 引住、编译不过，见 ⑨ 附注）。
- 复跑：`GOFLAGS= go test -count=1 -overlay=D:/tmp/wisp/20261023/367v1/overlay-perm-{after,before}.json -run '^TestV1' -v ./internal/perm`（panel／cmd 换各自 JSON；`cmd` 必须带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`）。

## ⓪ 判（一句话）

**`AC#2` 能翻。** 形＝〔甲′〕按判落地、牙真咬、六对与两条 arm 都经我独立复现。**翻的条件不是再改代码**，而是把 r1 证据件里**两处不成立的句子**更正（⑥「收益证明」、①「复跑」），另有一处**过期注释**与一处**稳健性建议**随件记；四条都不碰产码字节。

## ① 票面前提新鲜度 ＋ 复核 `367-a1` 的三条

| 断言 | 立案时 blob `4ee57d318` | 本笔 HEAD `7fabccf1b` |
|---|---|---|
| `func (s *Store) Set` 在 `store.go:175` | **175** ✔ | **189**（+14：本笔插入门与注释） |
| 旧判据 `if to == risk.ModeAutoApprove` 在 `:191-192` | **191** ✔ | 已被 `:209 if IsWidening(from, to) && s.confirm == nil {` 取代 |
| 产线唯一调用点 `composer_handlers.go:141 h.Modes.Set` | **141** ✔ | **141 未动** ✔ |
| 那 1 枚在 `Confirm == nil` 门之后；门在 `:133` | **133** ✔ | **131**（-2：`modeIsWidening` 那一段被删） |
| 产线调用者 **1 枚** | 1 ✔ | **1**：全仓非测试 `.go` 里 `perm` 形状的 `Set` 只有 `composer_handlers.go:141`（其余是 `Header.Set`／`reflect.Value.Set` 同名不同物） |
| 「加固、⛔ 补洞」 | 成立 | **仍成立**：两入口实测都拒（③ arm=none），`Set` 新门在产线上**走不到** |

★**行号会说话**：后续票再引 `:175`／`:133` 都是**过期数**，要引 `:189`／`:131`。
★**census 那条「过期处」没被 r1 收拾干净（F3，纯注释）**：`367-a1 ①` 点名 `composer_handlers.go:36-40` 的 “Nothing calls this handler yet” 已假；r1 **重写了该文件头三段**，但 `:41-45` 仍逐字写 “WHAT THIS IS NOT: a path. **Nothing calls this handler yet** … tickets 33/35 are still ready-for-agent” ⇒ 与 `composer_dispatch.go:204` ＋ 两入口实测矛盾。

## ② 核心判点一：六对我自己核（`Confirm=nil`，真 `config.toml`）

| 从 → 到 | 方向 | 改前（旧判据） | 改后（本笔） |
|---|---|---|---|
| `ask_every_step → ask_high_risk` | 变宽 | **`applied=1`、零问答、文件真落 `ask_high_risk`、`detail=""`** | **`refused-confirm applied=0`**，文件仍 `'ask_every_step'` |
| `ask_every_step → auto_approve` | 变宽 | `refused-confirm applied=0`（旧句） | `refused-confirm applied=0`（新句） |
| `ask_high_risk → auto_approve` | 变宽 | `refused-confirm applied=0`（旧句） | `refused-confirm applied=0`（新句） |
| `ask_high_risk → ask_every_step` | 变窄 | `applied=1` | **`applied=1`（未变）** |
| `auto_approve → ask_every_step` | 变窄 | `applied=1` | **`applied=1`（未变）** |
| `auto_approve → ask_high_risk` | 变窄 | `applied=1` | **`applied=1`（未变）** |

**总数核对**：变宽 3 对 → 改前 **2 有门**（两对目标都是 `auto_approve`）、改后 **3 全有门**；变窄 3 对 → 改前/改后 **0 拦**。⇒ **硬边界 1（⛔ 改变窄语义）未破。**
逐字新句（store 的）：`perm: 把档位从 ask_every_step 放宽到 ask_high_risk 需要一条确认通道，但本机没有接入（fail-closed 未切）`。

## ③ 核心判点二：`arm=half`（我自己造的那一发）

尺＝真 `panel.ModeWriteHandler` 在真 `*perm.Store` 前，只换确认腿接没接；同一请求 `AE→AH`。

| arm | handler 腿 | store 腿 | 改前 | 改后 |
|---|---|---|---|---|
| `none`（产线装配） | nil | nil | handler 拒，`applied=0`，文件 `'ask_every_step'` | handler 拒，同 |
| **`half`** | **接上** | **nil** | **`handle_err=<nil>`、`applied=1`、文件 `'ask_high_risk'`（静默变宽）** | **`applied=0`**，句是 **store 的** `perm: 把档位从 … 放宽到 … 需要一条确认通道…`，`last_result=refused-confirm` |
| `both` | 接上 | 接上（yes） | `applied=1` | **`applied=1`** |

**判：这「另做一半」是必要的，也是充分的 —— 但射程要写清。**
- **必要**：产线是 `none`，两入口实测都在 handler 就拒（`reply`/`handle_err` 同一句、文件不变）⇒ **产线那一发永远走不到 `Set` 的新门**，只有 `half` 能把「store 自己那层门」单独打亮。
- **充分**：`half` 正是「上层有腿、下层没腿」这一形状，改前 `applied=1`／改后 `applied=0` 且句源从 panel 换成 perm ⇒ 门确实在这一层咬住。
- **射程（⛔ 当活洞）**：`half` 这个装配**在生产里不存在** —— 今天两层要么都 `nil`（`panel_inbound.go:254/264`）、要么都非 `nil`（`run.go:667/687` 同一枚 `confirm`）。所以这一发读的是**病的形状**，⛔ 读「今天有一条活路」。它不改变①的「加固」判。

## ④ 核心判点三：`arm=both` ⇒ 门不是一律拒

`AE→AH` ＋ 真 `Confirm` 答 yes ⇒ **`applied=1`**、文件真落 `ask_high_risk`、`MODE-SWITCH … result=applied`。⇒ **门是 fail-closed，⛔ 一律拒。**
★同一发里 **`store_confirm_calls=0`、`handler_confirm_calls=0`** ⇒ **M4 的不对称没动**（只有 `auto_approve` 吃问答）。这条同时是 ⑩ 票 371 的读数。

## ⑤ 「谓词只有一份」＋ 依赖方向（我自己核）

- `git grep modeIsWidening -- '*.go'` ⇒ **0 命中**（`internal/panel` 里那枚**真删了**；命中只剩 docs/探针快照与新测试的一句沿革注释）。
- `git grep 'func IsWidening' -- '*.go'` ⇒ **1 命中**：`internal/perm/direction.go:26`。调用点 **2**：`store.go:209`、`composer_handlers.go:131`；钉子 **2**：`store_367_test.go:187`（perm）、`composer_handlers_test.go`（panel，驱动 `perm.IsWidening`）。
- **反证（我加的一发 M3：把谓词翻成 `to < from`）** ⇒ 两包 `rc=1/1`，**两枚梯子钉同时红**：`store_367_test.go:195/198`（“must read as widening/narrowing”）＋ `composer_handlers_test.go:113`（同句）⇒ **两枚钉子确实驱动同一个函数**，⛔ 各钉各的。
- **依赖**：`go list` ⇒ perm 只 import `context fmt config risk sync time`；`git grep wisp/internal/perm -- internal/panel` ⇒ **3 命中（全是本笔新加）**；反向 ⇒ **0**。⇒ `panel → perm` 下行，`perm → panel` 反向 ⛔。
- ★**执行缺口（如实报）**：没有任何用例**禁止**「再抄一份同 body 的谓词」；能钉的是「两层同判」（agreement 测试），⛔ 文本唯一。⇒ 「只有一份」今天靠**评审**守，⛔ 靠门禁守。

## ⑥ 反向判据 M1／M2（我自己造，突变全 `-overlay`）

- **M1（`:209` 回成旧判据 `to == risk.ModeAutoApprove`）** ⇒ **perm rc=1 / panel rc=1**。红名与红行：
  `TestSetGatesEveryWideningAndNoNarrowing` → `store_367_test.go:78/82/85/88/98`；
  `TestPanelGateAndStoreGateAgreeOnEveryPair367` → `composer_gate_agreement_367_test.go:76/91`。
- **M2（变宽＝任何变化 `from != to`）** ⇒ **perm rc=1 / panel rc=1**，**红在 3 个变窄行**：
  `store_367_test.go:104/107/110/113/116`（`AH→AE`、`AA→AE`、`AA→AH` 各五行逐条报）＋ `composer_gate_agreement_367_test.go:76/97`（“is a narrowing and neither layer may refuse it”）。
- ⇒ **硬边界 1 有牙**：M2 只打变窄、M1 只打那一个新变宽，两发互不代偿。
- **M3 顺带暴露一处稳健性（F4，⛔ 阻塞）**：谓词一翻，`AE→AA` 就不再满足 `IsWidening` ⇒ 新门放行、直落 `if to == risk.ModeAutoApprove` 里 **`s.confirm` 为 nil 被调用 ⇒ panic**（`store.go:224`，两包测试二进制都崩）。今天的谓词正确所以不可达；但旧码那枚内层 `if s.confirm == nil` 是**免费的兜底**，建议后续（371 或加固笔）在 auto 分支补回一枚 nil 判，别把 panic 留给重构。

## ⑦ `TestTicket90PersistFailureKeepsMemory` 三步法（我自己核）

| # | 装配 | rc | 读数 |
|---|---|---|---|
| ① | `HEAD` 码 ＋ `HEAD` 测试（`-overlay` 换回两枚旧 blob） | **0** | PASS |
| ② | **本笔码 ＋ `HEAD` 测试** | **1** | FAIL；红在旧 `store_test.go:232` `persist failure not audited: [… result=refused-confirm …]` ⇒ 撞的是新门的 `ResultRefused`，⛔ `ResultPersistError` |
| ③ | 本笔码 ＋ 改后夹具 | **0** | PASS |

**「断言一字未松」我的核**：`git show 7fabccf1b^:internal/perm/store_test.go` 与 `HEAD:` 那一枚的**整文 diff** 只有一处代码改动 —— `New(Options{Manager: mc, Logf: rec.logf})` → 同调用**多加 `Confirm: func(…) error { return nil }`**，其余 `+8` 行全是注释；三条断言（`err != nil`／内存档位不变／`ResultPersistError` 恰好 1 条）**逐字未动**。⇒ **只给夹具加腿、⛔ 放宽门**。

## ⑧ 边界零 diff（逐项）

`git diff --numstat 7fabccf1b^ 7fabccf1b` 全量只有 **8 个文件**（7 个 `.go` ＋ 本笔 impl 件），故凡不在表内者皆 0；再逐项点名复核（各得 **0 文件**）：
`internal/memory/schema.go`（`ddlV*`）· `internal/memory/open.go`（`migrationDSN`）· `internal/memory/taskstate_196r2_test.go`（196-r2 护栏）· `internal/memory/open_196r3_test.go`（196-r3 的 DSN 逻辑）· `internal/panel/inbound_roster_253_test.go`（名册）· `cmd/wisp/**`（**0 字节**，含 `panel_dispatch_binding_roster_253r1_windows_test.go`）· `internal/risk/**`（**0 字节**）。

## ⑨ 门禁四数（分开报）＋ 格式甲乙两尺

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `GOFLAGS= go test -count=1 ./internal/perm` | **0** | **19 PASS / 0 FAIL** |
| 2 | `go test -count=1 ./internal/panel` | **1** | **166 PASS / 5 FAIL**，**逐名同集既有红**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8` ⇒ **本笔新增红 0** |
| 3 | `PATH=…/third_party/sherpa-onnx:… go test -count=1 -run 'TestReplyListenerReject…|TestTicket223ModeLoosening…|TestTicket101' ./cmd/wisp` | **0** | **7 PASS / 0 FAIL**（含 `TestTicket101ManualSwitchSurvivesRestart`，`asked==[auto_approve]` 那条断言在 `run_mode101_test.go:251`、本笔 0 改动） |
| 4 | `go build ./...` ／ `go vet ./internal/perm ./internal/panel` | **0 / 0** | — |
| 5 | `sh scripts/d22scan.sh` ／ `sh scripts/check-path-length-budget.sh` | **0 / 0** | `clean - no D22 ban violations`（ban #8：internal 639、cmd 149）／`VERDICT GREEN`（over-budget 61＝roster 61、`not in roster=0`） |

**格式甲（入库 blob，判据）**：7 个 `.go` 逐枚 `git show HEAD:<path> | gofumpt.exe -l` ⇒ **全 CLEAN（空）**。
**格式乙（工作树，只诊断）**：同 7 枚 `gofumpt -l` **空**、`gofmt -l internal/perm internal/panel` **空**。
★**尺本身双向验过**：坏输入被点名（`<standard input>`）、好输入静默 —— **但只在 Git Bash 的管道里成立**；同一发在 PowerShell 管道里**好输入也被点名**（PS 改了字节）⇒ 本表全部读数取自 bash，⛔ 取 PS。
★`git ls-files --eol` 这 7 枚 **`i/lf w/lf`**；⛔ 拿 `gofumpt -l` 输出反推换行（那读的是工作树）。
**还原核的时机**：在**最后一发变异（M3）之后**量的 ⇒ perm `rc=0`、`19/0`；panel `rc=1`、`166/5`（同上 5 枚）；工作树四包 `git status` 空。
★**一处证据缺陷（F2，⛔ 代码）**：impl 件 ① 那句「复跑 `-overlay=…overlay-perm-{before,after}.json -run TestR1`（rc 0 / 0）」**今天跑不出来**：`before` 因 `direction.go` 被 `""` 删掉、而已入库的 `store_367_test.go` 引用 `IsWidening` ⇒ `[build failed] undefined: IsWidening`；`after` 缺 `R1ROOT` 环境变量 ⇒ 3 枚 `FAIL: R1ROOT unset`。读数本身我用自建探针复现了，⛔ 抄它的命令。

## ⑩ `AA` 拒绝句那处 ＋ 我判票 371 立得对不对

**(a) 换句对不对 ⇒ 对。** 旧句 `切到 auto_approve 需要一次 L2 强确认…` 对 `AE→AH` **说假话**（R20/M4 对非 auto 变宽**不问**），方向句 `把档位从 %s 放宽到 %s 需要一条确认通道…` 对三种变宽**都成立**，且**具名两端**（我用例核过：`store_367_test.go:93-96` 要求句里含 from 与 to）。
**(b) 旧句真的没钉吗 ⇒ 真的没钉。** `git grep '需要一次 L2 强确认|没有接入确认通道' -- '*_test.go'` 在 **`7fabccf1b^`** 与 **HEAD** 各得 **0 命中**（rc 1）；旧句今天只活在 docs／探针快照里。⇒ 换句没有绕过任何钉。
★**一处措辞隐患（记给 371）**：新句「需要一条**确认通道**」在 `AE→AH` 上加真腿时**并不会出卡**（④ `confirm_calls=0`），操作者可能读成「会问我一次」。建议 371 顺手把句子改成「…需要一条确认通道（`auto_approve` 才会真的出卡）」。

**(c) 票 371 立得对不对 ⇒ 立得对，但票面有两处要更正。**
1. **真**：`Set` 对 `AE→AH` 要求通道、而真通道对它 `confirm_calls=0` ⇒ **「门装上了、但卡不出来」**，我复现（④）。`cmd/wisp` 侧 `run.go:885-900` 的 L2 适配器只在 `to == auto_approve` 被调。
2. **要更正一（归属）**：这个形状**不是 367 引入的** —— `panel` 那层自票 114 起就是「只看 `Confirm` 的 nil-ness」，367 只是把**同一个形状**搬到 `store`，并把两层判据统一。票面应写「367 让它**可读**」，⛔「367 造成」。
3. **要更正二（「门是空的」）**：门**不是一律空** —— 没通道时它真拦（②③）；准确说法是「**门是通道存在性检查，⛔ 问答要求**；有通道的装配上，非 auto 变宽**静默通过**」。
4. **处置方向**：按铁律（最有价值＋主流 harness），应选**让 M4 对所有变宽出卡**（变宽＝人确认，主流 harness 的方向），⛔ 把 `IsWidening` 退回只判 `AA`（那是把 367 的类级收益扔掉）。这一条要动 `cmd/wisp` ＋ 改产品决定 ⇒ 摆机主批。

## ⑪ r1 上报的两处判断（我逐条判）

1. **`AA` 句换了** ⇒ **判对**（见 ⑩(a)(b)）。
2. **R20/M4「只有 auto 吃问答」它没动** ⇒ **判对**：`run_mode101_test.go:251` 在 `cmd/wisp/**`，是**本笔写面外**且属**产品决定**；它「另立一笔」的判断我同意。★但它由这条推出的 **⑥「收益证明」不成立（F1）**：`run.go:687` 那枚 handler 的 `Modes` 是 `run.go:667` 建的同一条 `modeStore`，`Confirm` 走 `rt.confirmModeSwitch`（**非 nil**）⇒ 哪天把它接到 dispatch 上，是 **`arm=both`** 那一行：`applied=1`、`confirm_calls=0`，**⛔「接上也得先把确认腿递给 Set」**。新门的真收益只在「将来有人建一个**没有腿**的 store 去变宽」这一形；这句要更正。

## ⑫ 提交 ＋ 逐笔越界

- 本件入库那一笔：**`6576470ed`**（本行由紧随其后的回填笔改成这个号，只改这一行）。
- 逐枚文件（每笔只含这 1 个）：`docs/evidence/s1/367-v1-verdict.md`。逐枚 `.go`：**0**。
- 越界核对：⛔ 改 `internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`；⛔ 落 `.go`；⛔ 翻框／改票名／`Status:`／`-done`／台账；⛔ push。探针／突变全在仓外，**工作树 0 突变**。

## 附：落笔前五问（尺的射程）

1. **量的是哪个对象**：`*perm.Store.Set` 与 `*panel.ModeWriteHandler` 的**运行行为**（真 `config.Manager`＋真 `config.toml`），⛔ 注释、⛔ 票面。
2. **输出形状**：`V1|…|` 一行一读数、`|` 分隔；错误一律 `<nil>` 或**原文**（先 render 再引号）。
3. **数的是代码还是注释**：调用点数**代码引用**（`git grep` 非测试 `.go`）；「旧句没钉」数的是 `*_test.go` 命中，⛔ docs。
4. **和总数对上**：六对＝3 变宽＋3 变窄；panel 5 枚既有红逐名同集；perm 19＝`^--- PASS` 计数；cmd/wisp 7＝`^--- PASS`。
5. **谁算的/算法/几位**：行数 `wc -l`（census=159、impl=160、票面=61）；diff 行数用 `git diff --numstat`（⛔ `grep '^[+-][^+-]'`）；哈希是 `git rev-parse`（SHA-1、40 位），⛔ `sha256sum`。
