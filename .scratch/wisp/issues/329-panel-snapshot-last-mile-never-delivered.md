# 票 329 — 快照的"最后一公里"：`lastPanelSnapshot()` 生产零读者，泵 marshal 出的字节今天走不到页面

- Status: **已立，派 `329-r1`（写腿·形甲）**（2026-10-16 16:2x 编排者立；来路＝只读普查腿 `33-a3`＋我自己逐字复认的三把尺；台账 `A1009`）
- 关联：票 33 `AC#14`（投递者那一格，本票交完才谈得上翻）／票 35 `:93`·`:391`（出向整跳，10-07 读数、行号已漂）／票 145 ⓑ（记账读者 0）／C17（`docs/PLAN.md:1367`）／`SPEC-08:150/158/163/174`
- 写面：`cmd/wisp/**`（本波该面**只许这一枚写腿**）；⛔ `internal/panel/**`（`327-r1` 在写）

## 这是什么（人话）

面板要显示"现在在干什么、有几条待批"，靠的是一串叫**快照**的字节。今天这串字节**已经做好了、也记进了账**，但**没有任何一条正式代码把它交给网页**——只有测试代码会去读它。所以用户在屏幕上看到的那一屏，今天不会跟着任务动。

这一跳缺的不是 API，是"**有人去调它**"：`internal/panel/pump.go` 里那句错误文案自己就写着"最后一公里在票 33/35 名下"，本票就是去把那一跳接上。

## 现量（⚠ 引用前先重跑；全部编排者 10-16 16:1x 自己在**工作树**跑，行号为真）

| 尺 | 命令 | 读数 |
|---|---|---|
| 记账点 | `grep -rn 'bookPanelSnapshot' --include='*.go' internal cmd` | 定义 `cmd/wisp/panel_pump.go:360`；生产唯一挂载 `cmd/wisp/run.go:732`（`Out: rt.bookPanelSnapshot`） |
| 那串字节的读者 | `grep -rn 'lastPanelSnapshot' --include='*.go' internal cmd tools` | 16 枚命中＝**定义/注释 2**（`panel_pump.go:417`/`:424`）＋**其余全在 `_test.go`**（`instructions_200r2_test.go` 3／`panel_l1_blocked_220_test.go` 3／`panel_pump_test.go` 3／`subagent_blocked_197_test.go` 2／`subagent_carrier_197_test.go` 2／`subagent_selfapproval_197_test.go` 2）⇒ **非测试调用者 0 枚** |
| 泵确实被驱动 | `grep -rn 'publishPanelSnapshot' --include='*.go' internal cmd \| grep -v _test` | 生产挂载 `run.go:742`（`rt.ui.publish`）＋`run.go:1063`（consoleSink），调用点 `run.go:1363/1452/1481` ⇒ ⛔ "整条没接"这句**说过头**，死的只有出向那一跳 |
| 宿主能送字节的原语 | `grep -n 'Eval(' cmd/wisp/panel_control_windows.go` | 定义 `:667`；生产 4 枚调用 `:737/:744/:748/:752` **全在 `msgcb` 里回 `window._rpc` 的 RPC 回执**，没一枚读快照 |
| 出向名册 | `grep -nE 'Method[A-Za-z]+ +=' internal/panel/bridge.go` | **7 枚**（`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`／`config.get`／`config.set`／`config.models.list`），**没有快照这一族**；两枚会拦住新增名的钉＝`internal/panel/git_test.go:387`、`internal/panel/inbound_roster_253_test.go:63`（`wantFullInboundRosterSize253 = 7`） |

**页面侧第二处未接（不在本票写面，`frontend/**` 归机主的外部前端 agent）**：`frontend/src/main.tsx:97` 渲染 `<App/>` **不传 snapshot**、`:110` 只挂一次 ⇒ 落 `EMPTY`；`frontend/src/lib/panel.ts` 声明的全是**页面→Go 出向门**（`window.wispDispatch` `:155`、`hostBridge()` `:161`），**没有任何 Go→页面接收器**。⇒ 本票只负责把 **Go 这半**接通并**把页面那半的合同写成文件**；⛔ 任何腿不许把这两半合并成"端到端已通"。

## 裁定（编排者 10-16 16:2x，落台账 `A1009`）

**选形＝甲**：Go 侧在"publish 真发生"的那一刻，经宿主把**已 marshal 的那串字节**推给页面（沿用 `panelControl.Eval`，推一枚具名全局调用），**⛔ 新增入向 `Method*` 常量**。
- 甲⛔ 碰契约：`SPEC-08:150/158/163/174` 已把 Go→前端推送标成**出向**行（请求列＝—），C17（`PLAN.md:1367`）本就要求"Go→前端事件推送"⇒ 本票是**兑现冻结契约**，不是改它。
- 乙（新增一枚入向名、页面主动拉）⛔ 不走：要同时撞 `git_test.go:387` 与 `inbound_roster_253_test.go:63` 那两枚名册钉（＝改名册数＝骗过自己的尺），且页面同样得改（本编队零写面）⇒ 若将来要乙，先落一枚**人批记录**（`A##`）再动。
- 丙（真机手测当凭据）⛔：属〔仅本机可量〕。
- **背压/合并规则⛔ 本票发明**：D38d 那条"多久推一次"的合并规则仍在票 35 名下；本票只沿用既有调用点（状态真动时 publish），**⛔ 逐 token delta 推**。

撤销口令＝**「329 改形乙」**（届时先量那两枚名册钉的解冻代价）。

## 验收（⛔ 实现者自勾；本腿只交料，翻框归编排者＋非实现者）

- [ ] **AC#0 交件三把尺**：逐枚提交号＋逐枚 `git show --name-only`（名册只许含 `cmd/wisp/**`／`.scratch/wisp/probes/329/r1/**`／`docs/evidence/s1/**` 一枚短名件／本票末尾**一条** bullet）；证据件 `wc -c`＋行数＋占位符 0；⛔ 0 字节件。
- [ ] **AC#1 起手先自己复认那三把尺**（⛔ 抄本票读数）：`lastPanelSnapshot` 非测试调用者＝**0**、`publishPanelSnapshot` 生产调用点在（⇒ ⛔ 把本票读成"整条没接"）、`bridge.go` 入向名＝**7**。三把各自带正控（同尺去一个已知存在的名字必须出数）。
- [ ] **AC#2 落地那一跳**：新增**一条生产代码**，在泵真的 publish 时把快照字节交给宿主；★判据尺＝`grep -rn 'lastPanelSnapshot' --include='*.go' internal cmd | grep -v '_test.go' | grep -v '^cmd/wisp/panel_pump.go:4'` **非空**（＝至少一枚非测试读者，且不是定义行自己）。
- [ ] **AC#3 逐位相等，不是"调用过"**：用例必须断**桩宿主收到的那串字节 ⇄ 泵 marshal 出的那串字节逐位相等**（⛔ 只断 `Eval` 被调过一次）；★去牙一发＝把推出去的字节换成空串或日志摘要，**指名那枚用例必须红**；★反控一发＝把 marshal 侧换一个字段（例如少一条待批），用例也必须红（证明它真在比字节、不是比长度）。
- [ ] **AC#4 ⛔ 撞那两枚名册钉**：交件时 `git diff --stat <起手锚> -- internal/panel/bridge.go internal/panel/inbound_roster_253_test.go internal/panel/git_test.go` 必须**为空**；⛔ 为了让新名字过关去改名册数／加白名单。
- [ ] **AC#5 页面侧合同写成文件**：新增一枚 `.scratch/wisp/probes/329/r1/page-contract.md`，逐字写清**具名 JS 全局／参数形状（就是那串快照 JSON）／触发时机**，供机主转给他的前端 agent；★票面必须**另起一行**具名"页面接收器今天缺位＝`main.tsx:97/110`＋`lib/panel.ts` 无 onmessage/订阅"⇒ ⛔ 本票任何一句不许写成"屏幕上已经会动了"。⛔ 动 `frontend/**`／`design/**` 一字。
- [ ] **AC#6 ⛔ 发明合并/背压**：只允许挂在既有 publish 调用点（`run.go:1363/1452/1481` 那三形）；红句不许出现"每个 delta 推一次"。若确实需要节流规则，**停下来上报**，⛔ 自己定。
- [ ] **AC#7 台面与颜色**：`cmd/wisp` 整包那一格＝★必须带 DLL 路径＋显式超时（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -timeout 30m`），四数逐枚点名（`^=== RUN` 含子测试 ⇄ `^--- PASS` 只数顶层），终态红名册与起手**作差＝0**；⛔ `go test ./...`；⛔ 用默认 10 分钟帽判死活。起手那一发⛔ 用共享工作树跑 ⇒ `git archive` 仓外导出树**八件套** `go.mod go.sum internal cmd tools scripts frontend models`＋手工补 `build/*.dll`／`third_party/sherpa-onnx/*.dll`／`frontend/dist`，树放 `D:` 卷且目录名要产生 8.3 短名别名（`/tmp` 那形会把 `paths_shortname_252_*` 那族打红＝台面故障，⛔ 计入判据）。
- [ ] **AC#8 门禁四数只在工作树跑**：`gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`／`go vet ./cmd/wisp/`，逐枚落 `rc=`；⛔ 把门禁指向部分导出树（`design/`／`.git` 缺 ⇒ `exit 2` 的 GUARD 射程拒绝，长得像红）。

## 已知会一起变旧的句子（⛔ 本票动它一字，只登记）

- `internal/panel/pump.go` 那句"快照已构造但没有出口（最后一公里在票 33/35 名下）"——`Out` 今天**已挂载**，那句错误的可达性本票交完要重看；但那包此刻有 `327-r1` 在写，⛔ 同波两枚写腿碰同一面 ⇒ 归**下一波**一枚独立裁决腿。
- 票 33 `AC#14` 里"本仓宿主对 `Run()`／`Dispatch()` 的调用者＝0 枚"（10-01 读数）已过期：现量 `cmd/wisp/panel_resident_windows.go:320 w.Run()`＋`:381 w.Dispatch(fn)`。⛔ 改原文，只在票尾记现量。
- 票 35／145 引的 `run.go:725`、`panel_pump.go:321/383`、`wantFullInboundRosterSize253 = 6` 全是 10-07 读数，现量＝`:732`／`:360`／`:424`／`= 7`。
