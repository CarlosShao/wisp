# ruler-dedup-1 — 面板能力尺查重普查（纯只读，零 Go 命令）

- 腿：`ruler-dedup-1`／时刻 **2026-10-08 13:54:42+0800**（`date` 现跑）
- 锚点：分支 `dev`，HEAD `53d73db4`（`git rev-parse --short HEAD` 现跑）
- **CWD＝仓库根 `D:\work\workspace\projects plans\Wisp`**（⛔ 本程所有相对路径一律按仓库根解析；台账 `A715` §4 那条基准差异我按具名声明处理）
- 只跑文本尺：`grep` / `sed -n` / `wc -l` / `git show` / `git status --porcelain` / `git log` / `ls` / `date`。**零 `go` 命令**（build/vet/test/list/env 一枚没跑）。
- 本程**未改任何既有文件、未 commit、未 push**。工作树里 759 行 porcelain（别人的改动）一枚没碰。
  我引用的每一枚文件都单枚复量过 `git status --porcelain -- <path>`＝**空**（工作树＝HEAD），读数见 `readings.md` §0。
- ⛔ 本件不出判据、不提议改票面。只做重复面读数。

---

## 表一：重复判定表

`判定` 的口径：**"接下来要建的那一枚"对着"盘上现存的那一枚"**是否同一枚尺。
`方向`：入向＝网页/宿主 → Go（能改变卡状态的那一侧）；出向＝Go → 页面（只读展示面）。

| # | 候选尺（票面出处 `<file>:<line>`） | 断言对象 | 方向 | 载体文件（盘上现状） | 判定 |
|---|---|---|---|---|---|
| R1 | 259 `AC#3` ①＝`PanelAPI` 方法名封闭集（`259-grant-binding-identity-ruler-holes.md:22`） | **一枚 Go interface 的方法名封闭集**：`approval.PanelAPI` 现 3 名（`Reject`/`Head`/`View`）＋`Gate.Panel()` 返回值的**具体**方法名集＋4 枚形状 type-assert（`Allow`/`AllowSession`/`Approve`/`Settle`） | 入向（答复面在 Go 侧的形状） | **现存**：`internal/agent/approval/ticket259_panel_capability_rulers_test.go:44`（名册）`:66`（接口三面检查）`:93`（载体＋形状断言 `:108`/`:113`/`:118`/`:123`） | **重复**。再立任何"PanelAPI 方法名尺"＝同一枚封闭集。⚠ 现存尺看不见的两支：①**既有签名的形变**塞能力（`Reject(correlationID, reason string) error` 改成带 grant 也只比名字，看不见）②页面侧路由字面名（那是 `internal/panel` 面，见 R6） |
| R2 | 259 `AC#3` ②＝`PanelItem` 按能力判的字段封名单（同上 `:22`） | **一枚 struct 的字段"能携带什么"**：Kind 级禁 Bool/Func/Chan/Ptr/UnsafePointer/Interface/Map/Array/Struct＋非标量切片（`:138-164`）；再加**值级**（渲染出的串不得等于或含 nonce／`it.bind`，`:184`）；再加**不变量级**（撤销 grant 前后逐字段等，`:216`） | 出向（Go→页面读面） | **现存 5 枚尺在同一枚 `PanelItem` 上**：`ticket259_panel_capability_rulers_test.go:169`、`:184`、`:216` ＋票 242 那两把词面尺**一字未动仍在盘**（`ticket242_panelface_test.go:30` 名册双向、`:53` 禁词串`grant/allow/approve/nonce/token`） | **重复**，而且是第 6 把。⚠ 现存名盲三支合起来仍看不见这一形：**定义型串字段**（`type Ticket string` 的新字段 Kind()==String 过 `:169`，名字与类型串都不含禁词就过 `:53`，不是当前 nonce/digest 的字面或子串就过 `:184`）。这一格是**射程缺口**，不是"再立一把封名单"的理由 |
| R3 | 259 `AC#3` ③＝语义侧"出向读面拿到的东西不可回填成答复"（同上 `:22`） | **不可回填**这件事本身：取 `Head`＋`View` 每一枚字段每一枚元素的全部可观察串（`:331`），逐枚当凭据递两枚在册入向入口，全须被拒；控制段（真 nonce 仍能开卡 `:307`、`pendingCount()>0` `:323`）跑在最后 | 两者（出向取料 → 入向回填） | **现存**：`ticket259_panel_capability_rulers_test.go:283` | **重复**。⚠ 看不见的两支：取料器 `r1FieldStrings`（`:236`）只收 `String` 与 `[]string` 两种 Kind ⇒ 非标量材料的回填**从没被试过**；入口只有 `Native().Allow`（`:300`）与 `DecideFromPanel`（`:316`）⇒ 第三枚入向入口不在射程 |
| R4 | 242 `AC#2` 那半句"面板面向的读面零尺"（`242-the-grant-binding-layer…-zero-rulers.md:21`＋标题句 `:1`） | 同一枚 `PanelItem`／`PanelAPI` 出向读面（前半＝字段名册反射；后半＝"把出向读面扩到能改变卡的状态那一形也要红"＝能力形） | 出向（后半带一支入向后果） | **现存**：R2 的 5 枚＋R1 的 3 枚；`ticket242_panelface_test.go:9` 自称"Ticket 242 AC#2 (form A…)" | **重复**。票面自己已经把这一格推出去了，两处原文逐字：242:21「**AC#2 出向读面不得带 grant**：反射扫 `approval.PanelItem` 的字段名枚数与名字 ⇒ 出现 `Grant`／可答复类字段**就红**；再加一发**能力侧**判据（不是词面）：把 `PanelAPI` 的出向读面扩到"能改变卡的状态"那一形也要红。」＋242:56「归**票 259 AC#3** 三枚能力尺（`PanelAPI` 方法名封闭集／`PanelItem` 按能力判的字段封名单／语义侧"读面不可回填成答复"），⛔ 不为此解冻乙形」 |
| R5 | 253 `AC#1`＝可达性钉改问能力（`253-panel-inbound-three-ruler-holes.md:16`；§8 的窄义改判 `:42`） | ①**产码里那扇门存在**：`(*PanelManager).installPanelTransport` 里的 `w.Bind(panelDispatchBinding, …)`（`cmd/wisp/panel_host_windows.go:801`、bind 点 `:802`、名字常量 `:80`）②页面发起一次调用后 **Go 收到并回执** | 入向（网页→Go） | 能力形＝**现存**：`cmd/wisp/panel_transport_35r1_test.go:195`（调真产码 `mgr.installPanelTransport(fc, ctx)` `:203`，红点 `:217` `fc.doorRounds != 1`）、`cmd/wisp/panel_inbound_guards_35r3_test.go:97`（同一枚产码函数）、`cmd/wisp/panel_resident_windows_test.go:314`/`:812`（票面 `:42` 自己已列）。词面/AST 形＝**零现存**：`grep -rn "panelDispatchBinding" --include=*.go cmd internal` 的命中全在产码（`panel_host_windows.go:80/:792/:802`）与测试自己的假件里，没有任何尺断言产码含这枚 `Bind` | **部分重叠**。对"再造一台能力形可达性仪器"＝重复（且那族全带 `//go:build windows`：`panel_transport_35r1_test.go:1`、`panel_inbound_guards_35r3_test.go`、`panel_resident_windows_test.go` ⇒ POSIX 分母为零）。对"词面/AST 钉"＝票面 `:42` 那句今天仍为真，不重复 |
| R6 | 253 `AC#2`＝名册尺射程覆盖不带 `panel.` 前缀那批（`:17`；§8 裁定走形ⓐ `:38`） | **入向路由字面名的封闭集**（现 6 名＝4 枚 `panel.*`＋`config.get`＋`config.set`），对 5 个分母双向对账：`bridge.go` 包级串常量／运行中的 `knownComposerMethod`（真问一遍，不复刻）／`bridge.go:145-147` guard 的 case 标签（AST）／`composer_dispatch.go` router 的 case 标签（AST，且要求 `default:` 必调 `rosterMismatch`，尺 `:494`）／全包 ANSWERED 字面 | 入向（但测的是 **Go 侧字符串常量与 switch case**，不是 interface 方法名） | **现存**：`internal/panel/inbound_roster_253_test.go:419`/`:597`/`:696`/`:769`/`:809`（834 行，票面 §9 `:47` 记为 `253-r5` 交件 `87bc6aca`） | 对 R1＝**不重复**（两枚封闭集不在同一宇宙：一枚数 `approval.PanelAPI` 的 Go 方法名，一枚数 `internal/panel` 的路由字面量；包不同、形状不同、方向虽同为入向但被测物不同形）。对"另立一枚全量入向名册尺"＝**重复**。⚠ 253 票面 `AC#2` 框**未勾**（票 253 open=4/done=0），票面 `:52` 逐字「⛔ **AC#2 这格我不翻**（实现者自述不算凭据）：待非实现者验收腿 `253-v1`」⇒ 状态＝已交件未结案 |
| R7 | 253 `AC#3`＝旧封套静默丢弃那一支要有声（`:18`） | **Go 侧对"收到但形状不认识"可见**：①有可读现场（审计行或错误文本）②页面拿到的不是无声 `null` | 入向 | 相邻件**现存**（不是同一枚）：`internal/panel/bridge.go:131-133` 未知方法具名拒收（逐字 `方法 %q 不是面板 composer 通路的能力入口`）＋`:136-141` 来源／requestId 两支；`internal/panel/composer_dispatch.go:208` `return "", d.rosterMismatch(req)` ＋ `:215-222` 函数体（`d.record(req, err)` 落审计）；转发那支由 `cmd/wisp/panel_transport_35r1_test.go:210` 钉（红句逐字含 `installPanelTransport registered no page->wispDispatch forwarding hook`） | **说不准**（缺哪一发读数，具名）：缺的是**票面现量 3 那一形今天还在不在**——页面按旧形状直接 `postMessage` 时，帧在第三方库 `webview.msgcb` 里被当成 `id=0`／未绑定方法而消亡，**Go 根本没收到**，因此 `bridge.go` 的具名拒收与 `rosterMismatch` 都轮不到它。我这轮只读到假件复刻了那一支（`cmd/wisp/panel_transport_35r1_test.go:176-187` 的 `fc.rpcDeadSlot++`），**没有任何尺断言真库里的这一支**；而②那一半属界面侧，票面 `:23` 已写明本票只写 Go 侧。⇒ 在补到"真库 msgcb 今天是否已被 35 AC#6 的转发全部接住"这发读数之前，不能判它重复也不能判它不重复 |

**这一族今天的净重复面（只报读数，不开方案）**：票 259 `AC#3` 的 ①②③ **三枚全部已在盘**（同一次交件 `b6b1d6a4`，`git log --oneline -1` 逐枚＝"票 259 腿 259-r1 第 2 步：AC#2 拒因可指名＋AC#3 三枚能力尺的码与新尺"），而 259 的 `AC#3` 框**仍未勾**；242 的 3 框全未勾；票 242 `AC#2` 那一半在票面上已被逐字归口到 259 `AC#3`。⇒ **对着 R1/R2/R3/R4 再"新立一把面板能力尺"＝与 `b6b1d6a4` 那次交件重复**（形状＝票 262 那族被 README `:80` 明令禁止的"重复劳动"）。

### 259 那两枚尺文件逐枚（func Test* 名＋断的是什么＋看不见什么）

`ticket259_panel_capability_rulers_test.go`（349 行／6 枚 `func Test`）：

| 行 | 用例 | 断的是什么 | 看不见什么 |
|---|---|---|---|
| `:66` | `TestTicket259R1PanelAPIMethodSetIsClosed` | `reflect.TypeOf((*PanelAPI)(nil)).Elem()` 的三面：每个在册名必须在 3 名 allowed 集里（`:71-75`）／allowed 集里每名必须仍在（`:76-80`）／`NumMethod()==3` 计数等式（`:81-84`，红句自己写明"否则一枚别名静默通过"） | 只比**名字**：既有方法换参数/换返回形（`Reject` 收 grant）它读不到；具体载体 `panelAPI` 多出的方法（下一枚才管）；页面侧路由名 |
| `:93` | `TestTicket259R1PanelCarrierCarriesNoAnswerVerb` | `Gate.Panel()` 的**具体值**方法名也必须在同一天 3 名集里（`:99-103`）＋4 枚 Go 类型系统的形状断言：`Allow(corr,grant)`/`AllowSession(corr,grant)`/`Approve(corr)`/`Settle(corr,allow bool)` 任一被满足即红（`:108-127`） | 形状表之外的动词（`Confirm(corr string, decision uint8)` 那一形不在 4 枚 assert 里，只能靠上一段的名字表抓）；字段侧一切 |
| `:169` | `TestTicket259R1PanelItemFieldCapabilityFence` | `PanelItem{}` 逐字段走 `r1PanelFieldIsDisplayData`（`:138-164`），**只读 Kind 不读名**：Bool＝判决位、Func＝动词、Chan/Ptr/UnsafePointer＝活句柄、Interface＝不可见载荷、Map＝开口载荷、Array/Struct＝按值加宽（要递归）、非标量切片＝凭据形状，全部禁；String／各整型放行 | 定义型串（`type Ticket string`）过得了；语义是凭据而名字类型都不露痕迹的字段；方法集侧一切 |
| `:184` | `TestTicket259R1PanelItemCarriesNoCredentialValue` | 对 `Head`＋`View` 两枚投影的每字段每元素，逐串等于/含 `nonce` 或 `it.bind` ⇒ 红（`:198-206`）；红句只点字段名不点值（凭据不进日志，`:196-197` 注释写明） | 只认"同一枚 nonce/digest 的字面或子串"：**换形／再摘要后的凭据读不出**；非标量字段（`r1FieldStrings` `:236-252` 返回 nil）；未在册第三枚读面 |
| `:216` | `TestTicket259R1PanelItemProjectionIgnoresGrantState` | 撤销 grant（`q.revokeGrants(it.Corr)` `:219`）前后逐字段 `DeepEqual`，变了的字段逐枚指名红（`:227-231`） | 只试 `revokeGrants` 一条变异路；别的改变判决态的产码路（花牌／超时／并发答复）不在这枚的激励里；载具只有 1 枚卡 |
| `:283` | `TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer` | 取料器 `:331` 走 `Head`＋两枚 `View` 的全部可读串（2 枚卡，`:284-285`），逐枚试 `Native().Allow`（`:300`）与 `DecideFromPanel{Allow:true,Grant:v}`（`:316`）；控制段在最后：真 nonce 仍须能开（`:307`）且 `pendingCount()>0`（`:323`）；取料为空即 `t.Fatal`（`:294-296`，不许"读了零枚当过") | 两枚入口之外的第三枚答复入口；非标量材料的回填（取料器 Kind 限制）；`Source:"native"` 之外的来源形 |

`ticket259_denial_rulers_test.go`（288 行／6 枚 `func Test`）＝**票 259 `AC#2`（四种拒因可指名）那一格，不是 `AC#3`**；列在这里只为让下一腿不误把它算成能力尺的分母：

| 行 | 用例 | 断的是什么 | 看不见什么 |
|---|---|---|---|
| `:115` | `…DenialNamesMissingNonce` | 空 grant ⇒ `errors.Is(err, ErrBadGrant)`（对外仍合并）＋审计行逐字 `denial=missing-nonce`＋`FORGED-OR-STALE` 那句没丢 | 界面侧是否把这半句露出来（按裁定不许露） |
| `:134` | `…SpentOrNeverLiveNonce` | 路由上未发行 nonce ⇒ `denial=spent-or-never-live-nonce`；store 级二次 spend ⇒ `denialSpentNonce` | `misbound` 那一支（票面 §9 与 `:42` 的射程注＝路由级不可达，本文件 `:26-31` 逐字声明"不表态"） |
| `:165` | `…DenialNamesMisbound` | **只在 store 级**：`s.spend(nonce, minted+"-tampered")` 必回 `denialMisbound`；拒后仍烧牌（`s.live()!=0` 即红） | 任何经路由的调用（这枚尺自己不主张路由能造出它） |
| `:192` | `…DenialNamesCardNotPending` | 手写 `it.state = stateAnswered` 后路由 ⇒ `ErrNotPending`＋审计 `denial=card-not-pending`＋该支不得烧牌；再直调 `q.deliver` 钉姊妹点同因 | 生产那两枚 `it.state` 写点真跑出的形状（`:187-190` 注释写明是因为它们与 `dropLocked` 同临界区才手写的） |
| `:226` | `…FourDenialLabelsArePairwiseDistinct` | 4 枚 `grantDenial` 的 `label()` 两两不同、非空、不等于 `denialNone`；枚数等式 `len(seen)==len(cases)` | 四种原因**真的各自可达**（这只管名字不折叠，不管分支覆盖） |
| `:249` | `…OutwardAnswerStaysMergedAcrossDenials` | 对外句逐字钉 `approval: 原生令牌无效（缺失/已用/与本次请求不绑定）`；4 枚内部因不得出现在对外句里；两枚内部因仍须在审计里；且不得借用 5 枚既有拒因文案（`ErrPanelAllow`/`ErrNotAdmitted`/`ErrChannelUnavailable`/`ErrNotPending`/`ErrUnknownCorrelation`） | 除这 5 枚之外的文案重合；`ui.go` 之外别处是否把原因露出 |

---

## 表二：仓里现存的读盘型守卫名册

尺＝`grep -rln "go/parser" cmd internal`（22 条命中，其中 `internal/config/config.test` 是**未跟踪的二进制**：`git ls-files internal/config/config.test`＝空，不是尺）＋`grep -rn "os.ReadDir(\|os.ReadFile(\|parser.ParseFile\|parser.ParseDir\|filepath.Glob\|filepath.WalkDir" --include=*_test.go cmd internal`。

overlay 列的口径：**凡是走盘的（`os.ReadFile`／`parser.ParseFile(fset, path, nil, …)`／`os.ReadDir`／`filepath.Glob`／`filepath.WalkDir`），`go test -overlay` 喂不进它的眼**——它们读的是真目录里的字节，不是编译器看到的那份。仓内已定案的这句原话（提交 `3ee9b3ff`，票 197 口径）："**读盘与走 AST 的测试对 overlay 结构性失明**，活例 `TestSubagentStreamKeyHasOneMintSite` 与 `leg_sink_gate_131` 那族"。唯一反形：**解析测试自带源码串**的那一支是 overlay 可见的（overlay 改的就是那枚 `_test.go` 自己的字节）。

| `<文件:行>` | 射程（枚数／整包／全仓） | overlay 可见性 | 备注 |
|---|---|---|---|
| `cmd/wisp/panel_geometry_255_test.go:208`、`:316`、`:343`、`:376`（读取器 `:449`→`:451` `os.ReadFile`） | **具名 3 枚文件**（`panel_host_windows.go`、`panel_resident_windows.go`、`config_readers_255.go`；`:316` 与 `:208` 同一枚） | **失明** | 目录由 `runtime.Caller(0)` 导出（`:441-447`），不靠 cwd；`//go:build windows`（`:1`）；5 枚 `func Test` |
| `cmd/wisp/panel_locked_naming_33r11_windows_test.go:176`（`os.ReadDir(".")`）＋`:211` | **整包**（本包全部非 `_test.go`；`:186` 枚数为零即 `t.Fatal`＝空扫描不算过） | **失明** | `//go:build windows`；5 枚 `func Test`；★它的正控**不种在产码**、解析本文件自带的源码串（文件头 `:38-41` 逐字）⇒ **那几支是 overlay 可见的**（`:443` 的 `src` 是测试自己的串） |
| `cmd/wisp/panel_geometry_255r6_range_windows_test.go:42`＋`:78`＋`:82` | **整包**（dir 扫描） | 失明 | `//go:build windows`；`func Test` 枚数＝0（是别枚尺的取数件） |
| `cmd/wisp/panel_host_gate_test.go:33` | 具名单枚（`panel_host_windows.go`，`ImportsOnly`） | 失明 | 平台中性（本文件随 L1 core scope 一起跑，文件头 `:3-7` 逐字说明"读源文件而非开窗"） |
| `cmd/wisp/panel_host_gate_test.go:210`＋`:221`（用例 `:177`） | **整包**（`os.ReadDir(".")` 去 `_test.go`） | 失明 | ⚠ `:182` 有一支 `t.Skipf`（零构造点即跳过）＝按本仓规矩 SKIP 不算凭据 |
| `cmd/wisp/panel_inbound_33_test.go:207`／`:285`（`filepath.Glob("*.go")`）＋`:216`／`:296` | **整包**（cwd 相对 Glob） | 失明 | 依赖 `go test` 的 cwd＝包目录 |
| `cmd/wisp/leg_dispatch_gate_133_test.go:365`＋`:375` | **整包**（dir 枚举行列） | 失明 | 1 枚 `func Test`；`:1651` 自己在讲"本平台不编译的 func Test 读不到盘" |
| `cmd/wisp/leg_sink_gate_131_test.go:561`＋`:581`＋`:593` | **整包**（既走 AST 又走原文字节） | 失明（且被 `3ee9b3ff` 点名为结构性失明活例） | 1 枚 `func Test` |
| `cmd/wisp/panel_assets_143_test.go:267` | 具名单枚（`panel_assets.go`） | 失明 | 6 枚 `func Test`；`:98`/`:102` 另读子进程 stdout/stderr |
| `cmd/wisp/panel_reshow_255r1_windows_test.go:432` | 具名单枚（`hostFile255r1`，目录由 `runtime.Caller` 导出，见 `:427` 注释） | 失明 | `//go:build windows`；6 枚 `func Test` |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:340`／`:386`／`:521` | **具名 3 枚文件**（含跨包一枚 `internal/agent/approval/gate.go`） | 失明 | 目录由 `runtime.Caller`（`:78-80`）；`//go:build windows`；5 枚 `func Test` |
| `cmd/wisp/resident_approval_risk_268_windows_test.go:227` | 具名单枚（`resident_approval_windows.go`） | 失明 | `//go:build windows`；3 枚 `func Test` |
| `cmd/wisp/resident_ball_228_test.go:66`＋`:77` | **整包**（`os.ReadDir(".")`） | 失明 | 还扫 `//go:build` 注释本体（`:154`）；2 枚 `func Test` |
| `cmd/wisp/resident_grant_writer_265_windows_test.go:540`＋`:468`／`:644` | **整包 dir 扫描**＋具名单枚 | 失明 | 目录 `runtime.Caller`（`:630`）；`//go:build windows`；9 枚 `func Test` |
| `cmd/wisp/resident_hotkey_258_test.go:28` | 具名单枚（`resident_windows.go`） | 失明 | 2 枚 `func Test` |
| `cmd/wisp/subagent_carrier_197_test.go:466`（`filepath.WalkDir(root, sub)`） | **一个子树**（root＋子目录） | 失明 | 3 枚 `func Test`；同文件的 `:557` 是反射尺（见表二续） |
| `internal/agent/spill_path_invariant_test.go:436` | 具名单枚（`spill.go`，`ParseComments`） | 失明 | `:142` 的 `os.ReadDir(dir)` 是运行时目录扫描，不是源码尺；4 枚 `func Test` |
| `internal/ball/tokens_table_test.go:373`／`:1231` | **具名 2 枚文件**（路径由 `runtime.Caller` `:101`/`:114`） | 失明 | 5 枚 `func Test` |
| `internal/memory/artifacts_path_invariant_test.go:771`＋`:760` | **整包**（`os.ReadDir(".")`） | 失明 | 5 枚 `func Test` |
| `internal/panel/composer_dispatch_test.go:618` | **全仓生产 `.go`**（`WalkDir`，显式 SkipDir：`.git`/`.scratch`/`node_modules`/`dist`/`third_party`/`scripts`/`docs`/`frontend`/`design`/`build`，逐字在 `:604-606`） | 失明 | 12 枚 `func Test`；解析失败算发现不算跳过（`:620` 注释） |
| `internal/panel/inbound_roster_253_test.go:136`／`:250`／`:543`（`os.ReadFile`＋`parser.ParseFile`） | **具名 2 枚（`bridge.go`、`composer_dispatch.go`）＋整包**（分母⑤ `productionPanelFiles(t, dir)`，调用点 `:484`，循环 `:511-513`）；另有 TempDir 里的种形件 `:632`/`:658`/`:712`/`:738`/`:771` | 失明 | 仓库根由 `panelRepoRoot(t)` 导出＝**`os.Getwd()` 上溯找 `go.mod`**（定义在冻结件 `internal/panel/tokens_fourway_test.go:57-67`）⇒ 属 cwd 依赖那一族（`A715` §4 的基准风险面）；5 枚 `func Test` |
| `internal/panel/l2_grant_boundary_test.go:738` | **dir 扫描**（产码枚举，跳 `_test.go`；票 253 §8 `:38` 已引它这一条） | 失明 | ⚠ **三枚冻结件之一**（242:33、259:29 逐字列名）。同文件 `:2032` 解析的是测试自带的 `bridgeSrc` 串 ⇒ **那一支 overlay 可见**；7 枚 `func Test` |
| `internal/projctx/projctx_test.go:263`＋`:273`＋`:283` | dir 扫描＋测试自带串（`:283` 的 `src`） | 前者失明／后者可见 | 10 枚 `func Test` |
| `internal/tools/paths_twocontainments_252_r2_test.go:194`＋`:206` | **整包**（`os.ReadDir(".")`） | 失明 | 2 枚 `func Test` |
| `internal/tools/task_state_188_test.go:261` | **全仓生产 `.go`**（`WalkDir(root)`，去 `_test.go`；`:248-249` 注释逐字："零分母的扫描报告 clean＝一棵从没读过的树"） | 失明 | 4 枚 `func Test` |
| `tools/d22scan/main.go` | **全仓 `walkGo`**，按设计跳 `*_test.go`（票 33 §文件头 `:44-48` 逐字引它） | 不是测试尺（CI 门禁） | 列在这里只为让下一腿别把它算进尺的枚数 |

### 表二续：仓里现存的"按形状/能力判"的反射尺（⛔ 不读盘，编译器可见 ⇒ `-overlay` 喂得进它的眼）

尺＝`grep -rln "NumMethod()" --include=*_test.go cmd internal`＝6 枚文件（加 `ticket259_panel_capability_rulers_test.go` 自身）。

| `<文件:行>` | 被测类型／名册 |
|---|---|
| `internal/agent/approval/ticket259_panel_capability_rulers_test.go:66`＋`:93`＋`:169`＋`:184`＋`:216`＋`:283` | `approval.PanelAPI`（3 名）＋`Gate.Panel()` 载体＋`approval.PanelItem`；票 259 `AC#3` |
| `internal/agent/approval/ticket242_panelface_test.go:30`＋`:53` | `approval.PanelItem` 的 6 名串名册（`panelItemReadFace`，`:21-28`）＋禁词串表；票 242 `AC#2` 甲形 |
| `internal/agent/approval/ticket220_l1_window_read_test.go:338` | 读面类型 `NumMethod()==0`（`:343`）＋对 `*Gate` 方法集的一维扫描（`:361-362`） |
| `internal/panel/composer_test.go:171` | `ModeView` 方法集无写面（`:172-173`） |
| `internal/tools/grant_test.go:508` | `GrantSource` 恰 1 枚方法（`:522`）＋逐名检查（`:531`） |
| `internal/tools/ticket224_setter_scope_test.go:41` | 桥接收发器方法集扫描（`:81`）＋`GrantSource` 枚数等式（`:102-103`） |
| `internal/tools/ticket90_test.go:431` | `ModeSource` 恰 1 枚方法（`:434`）＝"读面不带 allow 权" |
| `cmd/wisp/subagent_selfapproval_197_test.go:557` | **4 枚根类型**递归扫 allow 门（`toolAssembledRoots197` `:547-554`＝`tools.SubagentDeps`/`tools.TaskDeps`/`tools.Options`/`agent.Options`），allow 门名册 `:541`＝`Allow`/`Native`/`DecideFromNative`/`DecideFromPanel`/`GrantNonce` |

⇒ **"按能力判"这一族今天不是 3 把、是 8 枚文件在册**（其中名盲的只有 259 那 4 枚：`:169`/`:184`/`:216`/`:283`；其余按**枚数等式或名字表**判）。197 那枚与 R1 **部分重叠**（同一"allow 门"概念、不同被测根）。

---

## 表三：归口句名册（哪张票把哪一格推给了谁，逐字）

| `<文件:行>` | 谁推给谁 | 原文逐字（截取归口那半句） |
|---|---|---|
| `.scratch/wisp/issues/README.md:70` | README 规矩 9 的两处不彻底 → **票 262** | `⚠ 两处已知的不彻底，都归**票 262**，别把这行当"已经装好牙"：① 这条帽**只有词面**——今天没有任何仪器` |
| `.scratch/wisp/issues/README.md:80` | 机主原话——**禁与票 262 重复**（本程那句"新立尺先查重"的定式出处） | `⛔ 不许为此新造脚本——那会与票 262 重复，机主原话"不要重复劳动哈"）。` |
| `.scratch/wisp/issues/README.md:79-80` | 缺口定性＝排程不是仪器 | `剩下的缺口只有一个，而且是**排程**缺口不是仪器缺口：**没任何东西要求"开票者在提第一笔之前跑它一发"**。` |
| `.scratch/wisp/issues/README.md:83` | "有没有牙"这一判 → `262-v1` | `归 `262-v1`（它此前被按住的原因＝会种超长文件名、把别人正在跑的门禁读数洗脏）。` |
| `242-…-zero-rulers.md:18` | 242 `AC#1` 的"路由级＋拒因指名"那一半 → **票 259**（编排者 10-08 就地加的射程注，`cade79b3`＝HEAD 的祖先，已用 `git merge-base --is-ancestor` 复认） | `⇒ 那一半**归口票 259**：它 `:12` 逐字记着"对外只有一个 `ErrBadGrant`、三种原因并成一句"，`internal/agent/approval/approval.go:477` 逐字写着 **"ticket 259 does not move it"**，而**选边（ⓐ具名降级／ⓑ补牙）是票 259 `AC#1` 明写了归编排者的那一格**。⛔ 任何腿不许为了闭合本格去自己改合并文案。` ⚠ 这半句在**射程注**里（带时刻的编排者追加），不是 `AC#1` 原判据正文 |
| `242-…-zero-rulers.md:56` | 242 `AC#2` 不成立那一格 → **票 259 `AC#3` 三枚能力尺**（本程最重的一条归口，它在 242 的 `AC#2` 与 259 的 `AC#3` 之间划了同一枚石头） | `归**票 259 AC#3** 三枚能力尺（`PanelAPI` 方法名封闭集／`PanelItem` 按能力判的字段封名单／语义侧"读面不可回填成答复"），⛔ 不为此解冻乙形（腿判"缺的零件可在甲内补"）。` |
| `242-…-zero-rulers.md:60` | 242 的载具前置 → **票 259 `AC#4`** | `该文件写面 `cmd/wisp` 此刻被 `255-r2` 占着 ⇒ 归票 259 AC#4 串行。` |
| `259-grant-binding-identity-ruler-holes.md:22` | 259 `AC#3` 自认＝**补 242 `AC#2`**，⛔ 不另开一票 | `**AC#3 三枚能力尺（补票 242 AC#2 不成立那一格，随本票一起做，⛔ 不另开一票）**` … 末句：`⛔ 这三枚只加尺，**不许为此解冻乙形**（`242-v1` 已判"缺的零件可在甲内补"）。` |
| `259-grant-binding-identity-ruler-holes.md:37` | 母票三格的凭据 → 本票与 `242-v1` 那份表 | `母票 242 的三格**保持未勾**、凭据归本票与 `242-v1` 那份表，⛔ 不许在 242 面`（该行到此被票面自己截断） |
| `253-panel-inbound-three-ruler-holes.md:6` | 253 **不接管**票 114 的 9 枚未勾格（那 9 枚归票 114） | `**⛔ 本票不接管票 114 的 9 枚未勾格**——那 9 枚是"差集"（档位／附件／工作区／消息），归票 114 自己的落地腿。` |
| `253-panel-inbound-three-ruler-holes.md:23` | 界面侧那半 → owner 带话（本编队永不转达） | `需要它配合的地方写进"归口"由 owner 带话。` |
| `253-panel-inbound-three-ruler-holes.md:31` | 与 114／242 **互不吞并** | `与票 114（差集九格）、票 242（`bindDigest` 绑定层与出向读面零仪器）**互不吞并**：那两票各自的名字与格子不许并进本票。` |
| `253-panel-inbound-three-ruler-holes.md:52` | 253 `AC#2` 的结案凭据 → 非实现者腿 `253-v1` | `⛔ **AC#2 这格我不翻**（实现者自述不算凭据）：待非实现者验收腿 `253-v1`，排在 `220-r1` 退出之后（那枚腿正写 `internal/panel`，同时跑测试＝互洗）。` |
| `275-…-exits-2.md:95` | 别枚票把 242／259 的**台件**认成在册件（禁顺手格式化） | `**(c) 19 枚里躺着两枚**别人的台件**，⛔ 不许当"格式回归"顺手格式化掉**：`probes/241/v1/posctl/badly_formatted.go`（票 241 的**正控**…）与 `probes/259/r1/gofumpt-negctl/probe.go`（票 259 的 **负控**）。` |

**票 259 `:22` 之外没有第二处把同一枚石头推给别的票的句子**（尺的归口面在 242/253/259 三张票里是**单向的**：242 → 259；259 不推给 253；253 只推给 114 与 owner）。

---

## 未勾框与已勾框现量（尺＝带 `[[:space:]]*` 的那把）

命令（CWD＝仓库根）：

```sh
for t in 242-…md 253-…md 259-…md; do grep -cE '^[[:space:]]*- \[ \]' "$t"; grep -cE '^[[:space:]]*- \[x\]' "$t"; done
```

| 票 | 未勾 | 已勾 | 未勾的是哪几枚（行号） |
|---|---|---|---|
| `242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md` | **3** | **0** | `:18` `AC#1`／`:21` `AC#2`／`:28` `AC#3`（票面 §5 `:63` 逐字："三格保持 `[ ]`、⛔ 不加 `-done`"） |
| `253-panel-inbound-three-ruler-holes.md` | **4** | **0** | `:16` `AC#1`／`:17` `AC#2`（★尺已在盘 `inbound_roster_253_test.go`，未翻勾）／`:18` `AC#3`／`:19` `AC#4` |
| `259-grant-binding-identity-ruler-holes.md` | **5** | **1** | 已勾＝`:19` `AC#0`（只读普查）；未勾＝`:20` `AC#1`／`:21` `AC#2`（★六枚尺已在盘 `ticket259_denial_rulers_test.go`）／`:22` `AC#3`（★三枚能力尺已在盘 `ticket259_panel_capability_rulers_test.go`）／`:23` `AC#4`／`:24` `AC#5` |

⚠ 状态读法（只报形状）：**票 259 的 `AC#2` 与 `AC#3` 是"尺已交件、框未翻"**（同一次交件 commit `b6b1d6a4`）；242 的 3 框按票面自己的裁定保持未勾。

---

## 我这条程自己不确定／看不见的那几格

1. **一枚绿/红读数都没有**。⛔ 本程禁跑 `go`，那 12 枚在册用例（259 两文件）今天是不是全绿、`M-D2`/`M-E` 两发是否仍被它们抓得住，我**看不见**；表一里"现存尺能不能抓 M-D2/M-E"全部是从代码形状推的射程陈述，不是运行读数。
2. **R7（253 `AC#3`）判不了**：缺的那发读数＝旧封套那一形今天到底停在哪儿。第三方库里那一支（票面现量 3 指的 `webview.go:162-168`〔待验〕）我没读，也不在上表任何一枚尺的分母里；`cmd/wisp/panel_transport_35r1_test.go:176-187` 的 `rpcDeadSlot` 是**假件对它的复刻**，不是对它的断言。
3. **读盘尺对 overlay 失明这一条我只在文档层复认**：依据是提交 `3ee9b3ff` 的原话与票 253 §8 `:38` 那批射程读数；**有没有哪一步 CI 真的用 `-overlay` 喂过上面任何一枚读盘尺，我看不见**（只在 `.scratch/wisp/probes/252/v1/overlay-inv/` 那类台件名里看到过这族的痕迹，没去核那一步的调用形状）。
4. **`frontend/**` 一枚没读**（242:33／259:30 两层禁令）。因此 `internal/panel/composer_test.go:502` `TestTheRendererHoldsExactlyOneDoorToTheHost` 与 `:394` `routeLiteralRe` 的**分母内容**我无法核，只核到尺文件自己那侧的正则与闭集（`:394` 逐字 `"panel\.[A-Za-z.]+"`、`composerRouteLiterals()` 5 名，`:409-419`）。
5. **票 242 `:56` 那句"今天仍只由五处注释守着"我复量不成 5 枚**：`grep -rn "PanelAPI" --include=*.go internal cmd \| grep -v _test` 给 ≥11 行注释命中（`approval.go:352`、`gate.go:638`、`pending_read.go:7`/`:16`、`queue.go:380`、`replies.go:97`/`:347`/`:471`、`ui.go:155`、`cmd/wisp/approval_reply.go:29`、`panel_resident_windows.go:224`、`resident_ball_windows.go:231`）。**哪五处是那句话所指，我判不动**，不当它错、也不照它算分母。
6. **`259-r1`／`253-r5` 的验收状态**：两枚腿都交了件（`b6b1d6a4`／`87bc6aca`），而 259 `AC#2`/`AC#3` 与 253 `AC#2` 的框一枚没翻；**在册的非实现者验收腿现在排到哪一格**，不在本程射程内（我没读 `docs/reports/HANDOVER.md`，本程要答的是重复面）。
7. **`cmd/wisp/panel_geometry_255_test.go:316` 与 `:208` 是同一枚文件**（`panel_host_windows.go`）⇒ 我给的"具名 3 枚"是**按文件去重**的枚数；按读取点算＝4 处。这两种数法都摆在件里，谁引用谁得自己点明用的是哪一种。
8. **`inbound_roster_253_test.go` 的仓库根靠 `os.Getwd()` 上溯找 `go.mod`**（定义在冻结件 `internal/panel/tokens_fourway_test.go:57`）⇒ 它与我这条腿的读数**不是同一个基准**（我用的是"仓库根 CWD"）；如果某条腿从别的目录跑它，它自己会跟着 cwd 走，这个偏差我这一程**没法量化**。
