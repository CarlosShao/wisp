# writeface-1 — 10 逐腿文件名册

锚＝HEAD `934b350a`（本腿锚件笔 `5f2289e3` 只加 `probes/orch/20261011/writeface-1/**`，⛔ 源码字节）。
三态声明：`git status --porcelain -- internal cmd`＝**0 行**（锚件现跑）⇒ 本件所有 HEAD blob 读数**与工作树逐字相同**；
每条读数尺＝`git grep -n <pat> HEAD -- <pathspec>` 或 `git show HEAD:<path>` 落 `logs/` 后 `sed -n`，⛔ 编译面命令 0 条。
行号只作"此刻在哪"，判据＝内容锚（逐字前缀，⛔ 不许把行号抄进派单）。

## A. `296-r2`（票 296 `AC#5`＝三处过头话改成说实话）

**必动 1 枚**（尺＝`git grep -n -e "only skip is the DeepEqual" -e 'an empty \[hotkey\] slot means' HEAD -- cmd internal` ⇒ 2 命中，逐枚读 blob 补第 3 处）：

| # | 文件 | 位置（函数/块） | 三处内容锚（现量行号） |
|---|---|---|---|
| 1 | `cmd/wisp/resident_windows.go` | `hotkeyReloadSource296` 的**文档块**（函数体⛔，函数起 `:43`） | ①`:38` 逐字 `// an empty answer: Check()'s only skip `<br>②`:40` 逐字 `// here would re-register the set and le`（该句跨 `:40-:41`）<br>③`:42` 逐字 `// owner's machine (an empty [hotkey] sl` |

★**顶回编排者的第 1 读数（雷#1 的答案＝现读，⛔ 推理）**：`AC#5` 的注释**⛔ 住在测试件里，住在产码文件 `resident_windows.go` 里**。
第二把尺（反证）＝`git show HEAD:cmd/wisp/resident_hotkey_296_windows_test.go | grep -nE "only skip|dead until|turned off|DeepEqual"` ⇒ **0 命中**（测试件没复述这三句，改不到它）。
连带后果：`296-r2` 与票 298 那两枚件（`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`）**同名册零交集**。

**可能被动（⛔ 票面要求，是本腿推的）**：无。`AC#5` 判据自写"⛔ 一行可执行代码都不许多变"＋`git diff -U0` 只含注释行。

## B. `295-r1`（票 295 `AC#1`＝`started` 换 `atomic.Bool`）

尺＝`git grep -n "\.started" HEAD -- cmd/wisp` ⇒ **11 行**（整族，含 `_test.go`，含注释；分布＝生产 3／测试 8）＋声明行 `	started bool`（⛔ 带点，不在那 11 里）。
本腿**自己复跑**了 `295-a1b` ⓓ 节那把承重尺，**11 枚逐点对上**（编排者 10-10 09:0x 复跑的 11 亦同）⇒ ⛔ 照抄，是复现。

**必动 4 枚**：

| # | 文件 | 被动的函数/块 | 内容锚（逐字）／枚数 |
|---|---|---|---|
| 1 | `cmd/wisp/resident_audio_windows.go` | struct 声明＋`toggleMute`＋`assembleCapture` 的 `default:`＋`posture()` | `	started bool`→`	started atomic.Bool`；`	ra.started = gate.Open()`（`toggleMute` 体内）；`		ra.started = true`（boot 支）；`	if !ra.started {`（`posture()` 体内）＝**4 处** |
| 2 | `cmd/wisp/resident_mute_290_windows_test.go` | `TestAC290BothMuteGesturesTurnTheGateAndBack` 等 | **4 枚读点**：3 枚条件位＋**1 枚 `t.Fatalf` 的 `%v` 参数位**（★参数位换形状后不红、只把日志打成 `{true}`；本腿在 blob 里逐字见到该形） |
| 3 | `cmd/wisp/resident_audio_247_windows_test.go` | 三枚 boot 侧用例 | **3 枚读点** |
| 4 | `cmd/wisp/resident_audio_247_live_windows_test.go` | 真机件 | **1 枚读点** |

**同批必改注释（票 293 定式"旧句不许与新代码并存"）**：`	// started mirrors what the gate now rep`（写点上方）＋`	// started records that the gate handed `（声明上方），都在文件 #1 内。
**⛔ 不在册（防顺手扩）**：`verdict`——尺＝`git grep -n "ra\.verdict" HEAD -- cmd/wisp/resident_audio_windows.go` ⇒ **10 命中，全部落在 `:263`..`:390`**（boot 侧 7 写含 `:321` 那条 Warn 读、`posture()` 侧 3 读），`toggleMute`（体起 `:156`）内 **0**；那几行里唯一一枚 `verdict` 字样是**字符串**（`:189` 的 `…（见启动时那句采集腿 verdict）`）⇒ 票面 `AC#0` 那句"也写 verdict"的前件不成立（与 `295-a1b` 一致，本腿自跑）。

## C. `294-r1`（票 294 `AC#1`/`AC#2`/`AC#3`）

**必动（本腿推：`AC#1` 那格⛔ 免新增文件不可）**：

| # | 文件 | 内容 | 依据（票面哪一格） |
|---|---|---|---|
| 1 | `cmd/wisp/resident_mute_294_test.go`（**新**；裁定节点名此路径；现量 HEAD **不存在**＝尺 `git ls-tree --name-only HEAD -- cmd/wisp/resident_mute_294_test.go` ⇒ 0 行） | `AC#1` 走查钉（读 AST）＋`AC#3` 的 `executed == true` | 裁定节 `AC#1` 落点裁 |
| 2 | `cmd/wisp/resident_hotkey_294_windows_test.go` 或并入 #1（★文件名票面与裁定都⛔ 指名＝本腿推的落点；⚠ `AC#2`/`AC#3` 要用 `residentAudio`/`toggleMute`/球窗，那些符号只在 windows 面存在 ⇒ 这一枚**必须**带 `_windows_test.go` 后缀，⛔ 与 #1 混成一枚非后缀件） | `AC#2` 真球窗 `IsLive(mute 镜像)` ＋**cmd/wisp 侧本地镜像常量 `hkMute… = 2`** | 裁定节 `AC#2`："缺的只是本地镜像常量"。现量：`git grep -n "hkMute" HEAD -- cmd/wisp` ⇒ **1 命中且是注释**（`resident_ball_windows.go:565`），镜像常量 **0 枚**；本包已有形状先例逐字 `	hkSummon258 uint32 = 1`（`resident_hotkey_258_windows_test.go:74`）；`IsLive` 在 `cmd/wisp` 整族只 **1** 枚读点（`:290`，读 summon）；被镜像的常量逐字 `	hkMute   = 2`（`internal/ball/hotkey_windows.go:37`，非导出）；仓内先例 `internal/ball/hotkey_live_test.go:404` 逐字 `	if !b.HotkeyReport().IsLive(hkMute) {` |

**只读面（被 parse、⛔ 被写）**：`cmd/wisp/resident_windows.go` 的 `runResident` 函数体（模板逐字在 `cmd/wisp/resident_hotkey_258_test.go`：`func runResidentBody258(t *testing.T) (*ast.BlockStmt, string) {`）；挂载那行现量 `	rb.attachMuteGate(raudio.toggleMute)`（现 `:342`）。
**`AC#3` 的反形正控＝瞬态变异、⛔ 入库写面**：目标是 `resident_audio_windows.go` 的 `	case gate.Muted():`／`	case gate.Open():` 两支末尾那两枚 `…, true`（blob 现读两枚俱在）。

**票面/裁定节读数已过期（本腿现量打旧，⛔ 改原句）**——四条全是票 293 的落地腿 `resident_tray_mute_293_windows_test.go` 进 HEAD 造成的位移：

| 原句（出处） | 当时 | 本腿现量 |
|---|---|---|
| "`attachMuteGate` 测试直调 **3**"（票 294 裁定节 21:1x／现量节 `:11`） | 3 | **5**（多出 `resident_tray_mute_293_windows_test.go:213`／`:293`） |
| 产码那枚 `:336` | :336 | **`:342`** |
| 定义 `resident_ball_windows.go:472` | :472 | **`:478`** |
| 手搓 `&residentBall{` **13** 枚（尺＝`git grep -n "&residentBall{" HEAD -- cmd/wisp`） | 13 | **16**（293 新件带 2 枚；`resident_ball_windows.go` 自身 1 枚在名单内，逐档分布见 `logs/r2-residentBall-literal.txt`） |
| `startResidentBall` 整族 12／真调用 10（裁定节 `:45`） | 12／10 | **射程差**：`-- cmd/wisp` 整族 **14**＝测试内 12（含注释 2）＋测试外 2（定义 `resident_ball_windows.go:302`＋生产调用 `resident_windows.go:250`）；**"测试真调用 10"复现一致** |

## D. `309-r1`（票 309 形 Ⓑ 注释＋`AC#1` 仪器＋`AC#2` 枚数钉）

**必动 2 枚**：

| # | 文件 | 内容 | 内容锚（现量行号） |
|---|---|---|---|
| 1 | `internal/tools/bridge.go` | `assessorFor` 文档块 `:920-923` **追加**＋守卫上方具名（⛔ 动 `:925` 布尔式） | `// assessorFor returns the assessor whose R4 input is bound to this task's`（`:920`）／`	if b.injected || b.prov == nil || taskID == "" {`（`:925`）／`		WithTaintDetector(b.prov.Detector(taskID))`（`:931`） |
| 2 | `internal/tools/*_test.go` **1 枚新增**（★文件名票面与裁定都⛔ 指名，本腿只报"枚数 1"；现量＝`git ls-tree` 该目录 `*_test.go` **53** 枚，⛔ 任一枚已含 `Test309`） | `AC#1` 恒真形仪器＋`AC#2` 枚数钉（同片写，票面 `AC#1` 逐字"⛔ 产码语义"） | 范本现量俱在：`internal/tools/paths_twocontainments_252_r2_test.go`／`task_state_188_test.go`／`failclosed_236_teeth_test.go`（三枚 blob 名已 `git ls-tree` 坐实） |

**枚数钉的分母（本腿复跑裁定节那句"60＝测试 59／产 1"，并补一枚它没写的维度）**：
尺＝`git grep -n "ToolRequest{" HEAD -- internal cmd` ⇒ **60**＝`_test.go` 内 **59**／外 **1**（`internal/agent/loop.go:675`，逐字 `		req := ToolRequest{`）。
★同一条尺去掉 pathspec 换成全仓 ⇒ **252** 枚（`.scratch/**` 变异拷贝＋`docs/**` 转述吃掉 192）⇒ 新钉必须把**射程目录**写进断言，否则分母随台面漂。
`AC#3`（`:869-883` 那段行号→内容锚）**同一枚文件 #1、另格**：裁定节点"排程更新"列的 `309-r1` 写面只含 Ⓑ注释＋`AC#1`＋`AC#2`，`AC#3` ⛔ 在列 ⇒ 派单要显式写"这格本腿取不取"。

## E. `298` 残余

**落地半边＝零枚文件（现量证成，⛔ 推理）**：
尺①＝`git merge-base --is-ancestor 18d43f8e HEAD` ⇒ **YES**（`298-r1` 的清洗笔已在 HEAD  ancestry 里）；
尺②＝两枚目标件 `git rev-parse HEAD:<path>` ⇄ `git rev-parse 18d43f8e:<path>` **逐字相同**（`9713dc37…`／`52d3c1b2…`，见 `logs/298-landed-check.txt`）⇒ 修后至今无人再动。
残余只有三件、**全是裁决/决策⛔ 代码**：`298-v1` 在飞的翻勾（它已交 `10-ac-readings.md`，`AC#0`..`AC#3` 判语俱"成立"）；`AC#2` 那句"要不要真换成本机 LF"＝编排者一句话（0 枚文件）；`gofumpt` 残留 1 枚（`panel_transport_35r2_test.go`）＝**⛔ 票 298 射程**，要走得另立一票。
⇒ **队列里⛔ 该有 `298-r2` 这样一枚落地腿。**
