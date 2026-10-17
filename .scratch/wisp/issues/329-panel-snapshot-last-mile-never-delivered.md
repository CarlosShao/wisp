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
- [2026-10-17 10:0x +08] agent=`329-r1`（落地腿·写·形甲）did=★把"记进账→交给宿主"那一跳接上：`cmd/wisp/panel_snapshot_deliver_329.go`（`panelSnapshotHost` 窄接口＝一枚 `evalPanelJS`；`window.wispOnPanelSnapshot(泵 marshal 字节原文)`＝具名全局出向调用，⛔ 新增入向 `Method*` 常量、⛔ 动那四枚 `_rpc` 回执）＋`publishPanelSnapshot` 记账后同发 `rt.pushPanelSnapshotToHost()`（＝`lastPanelSnapshot` 第一枚生产读者，现量 `panel_snapshot_deliver_329.go:95`；只挂既有三个驱动点，⛔ 逐 token delta、⛔ 合并/背压规则——那条仍在票 35 名下）＋`run.go` `snapHost/snapHostMu` 字段＋`panel_host_windows.go` `PanelManager.evalPanelJS`（Dispatch 回面板线程，沿用 `residentPanel.post` 形状；窗口未开＝不推，字节留在账上）＋`resident_windows.go` 装配根 `src.run.setPanelSnapshotHost(rp)`（装配失败／宿主未建两形＝不挂）。`AC#3` 尺＝`TestBookedSnapshotReachesHostBitForBit329` 三枚子尺（marshal 侧⇄推侧⇄状态钉）＋账上摘要行 sha256 同链；★成对突变两发各自红在不同枚尺（去牙红推侧 `:141`、反控红 marshal 侧 `:134`，均仓外导出树物理字节、逐发还原哈希对回起手值）；★必答题＝摘掉那枚生产读者后同一用例**响**（`:141`+`:164`）⇒ 尺对该跳敏感。`AC#2` 尺交件现量三命中（真读者一枚另两枚新件注释）；⚠ 登记票面两处过期数：现量表"16 枚命中"其自身列举求和＝**17**；`publishPanelSnapshot` 生产挂载行号随本腿漂移（`run.go:742→754`、`:1063→1075`）。★`AC#5` 合同件＝`.scratch/wisp/probes/329/r1/page-contract.md`（全局名/参数形状/触发时机逐字）；**页面接收器今天仍缺位（`main.tsx:97/110`＋`lib/panel.ts` 无 Go→页面接收器）＝票面原文继续成立，本票⛔ 任何一句不构成"屏幕上已经会动了"**。`AC#4` 名册钉空 diff rc=0、入向名册仍＝7。门禁四数 rc＝0/0/0/0（d22scan 第一发自触 ban #9 幽灵引用红，`1d1261e0` 追加更正复跑转净，⛔ amend）。★整包颜色如实交缺口：本机今天处票 314 窗口区早停实例活跃中——起手树四发（标准/复跑/skip1/skip2/skipset2）**没有任何一发跑到终点锚 `TestZZ315A1*`**，死点逐具名且逐发漂移（`TestAC4FocusReturnToPriorWindowGap33r5`→`TestAC13BringUpSurvivesAReusedThreadQuit`→`TestAC13BringUpRefusesAThreadWithAQueuedClose`→`TestAC14GoSideEvalPushReachesThePage`→`TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`），且同一枚起手树里机器态红名册本身逐发不同（`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`／`TestRunPacketLightsTheRosterRowOfItsLiveL1Window` 各只在个别发红、其余发绿）⇒ "起手 6⇄终态 6 作差 0"那格今天在这台机器上取不到（先于本腿成立），逐发四数与名册见 `docs/evidence/s1/329-snapshot-delivery-r1.md` §6，该格请编排者裁。本腿新用例在两种台面各自定向发均绿（工作树 1.65s／导出树 1.84s，且文件排序在窗杀族之后、整包早停发轮不到它＝具名）。⛔ 翻任何 AC 框（九格现量未动）、⛔ push、⛔ 动 `frontend/**`/`design/**`/三枚冻结件、⛔ 杀 PID 21892（交件时点复量完好）。
- [2026-10-17 10:3x +08] agent=`329-v1`（非实现者终裁腿·`cmd/wisp/**` 只读射程）did=★九格逐格判语落 `docs/evidence/s1/329-v1-verdict.md`（119 行／17869 字节）。**我自己重放的成对突变四发**（全打仓外导出树 `D:/W329VB`＝`1d1261e0` 字节，逐发当场还原、还原哈希对回起手值 `e7a11ceb0403…`／`728891008927…`、marker 残留 0、还原后正控 `--- PASS 1.47s`）＝① 去牙·空切片形 `panelSnapshotEvalJS(data[:0])` 红 `:141`＋`:164` 而 `:134` 绿；② ★补一发等长改内容（payload 首字节翻一位、长度一字不动）红同两行 ⇒ 正面回答票面"真在比字节、不是比长度"；③ 反控 marshal 侧少一条待批 红 lead＝`:134`（另一枚尺另一行，与①落点不同）；④ 必答题摘生产读者 `:95` 红 `:141`/`:164`/`:167`＋`host=""`＋"received 0 pushes"。★新抓一处不实：落地腿 §3 把①写成 `panelSnapshotEvalJS(nil)`，那一形在 `1d1261e0` 字节上 **build failed**（`cmd\wisp\panel_snapshot_deliver_329.go:95:5: declared and not used: data`），其盘上红句只能由空切片那一形产出（逐字两发见 `probes/329/v1/logs/mutation-v1-defang-empty.log` ⇄ `-r2.log`）。格判：`AC#0`＝**不成立**（r1 件两处没有读数＋§1 指的 `commits-roster.log` 我复量 `ls` rc=2 确不存在；缺陷在件不在字节，名册我已由 `git show --numstat` 逐枚复得＝七笔只落 `cmd/wisp/**`＋`probes/329/**`，issues/evidence 逐枚 0）；`AC#1`/`#2`/`#3`/`#4`/`#5`/`#6`/`#8`＝**成立**（#4：名册钉 diff 空 rc=0、入向名册两树皆 7、`panel_control_windows.go` A..B numstat 空＝那四枚 `_rpc` 回执一字未动；#5：合同全局名⇄`deliver_329.go:57`+`:65-67`、四必选键⇄`composer.go:57-61`、两可选 `omitempty`⇄`:74`/`:91`、卡面十字段⇄`approval.go:40-58`、时机⇄`panel_pump.go:442`/`:447-463` 逐字一致，一处行号漂移＝合同写 `lib/panel.ts` `hostBridge()` `:161` 实为 `:159`（`:161` 是其函数体那行），⛔ 影响实质；#6：驱动点 A⇄B 同数同形（`run.go:1375`/`:1464`/`:1493`）＋红句零枚 delta 推；#8：我自己在共享树跑的四数 rc＝gofumpt 0（列 `models.go`＋`panel_transport_35r2_test.go`＝非本票写面的既有欠账，同复核 r1 §5）／d22scan 0 clean／path-length 0 GREEN／`go vet ./cmd/wisp/` 0）；`AC#7`＝**不可判**，缺的那一维写清＝"整包终态红名册 ⇄ 起手红名册作差"今天在这台机器上取不到（r1 六发我逐枚 `grep -c` 复量 `166/115/2/0`…`206/151/0/0`，rc 全 1、`^ok` 全 0、改前末枚 `=== RUN`＝`TestAC13BringUpSurvivesAReusedThreadQuit`、改后＝`TestBallPanelGesturesReachThePanelThread`、两发尾皆停 `MARK299 wndproc` 之后＝整进程真窗区死亡，⛔ 读成全绿 ⛔ 读成 329 打红；我也⛔ 重跑整包赌一把）⇒ **另路取数我自己跑**＝具名 13 枚 `-run` 名册：起手树 A `12/12/0/0` rc=0 ⇄ 交件树 B `13/13/0/0` rc=0，红名册两侧皆空＝作差 0、唯一差枚＝新用例 `TestBookedSnapshotReachesHostBitForBit329`。能力必答：①**被产码读到＝成立**（逐环 `deliver_329.go:95`→`:105`←`panel_pump.go:463`←`run.go:1375/1464/1493`←挂载 `:754`/`:1075`←装配根 `resident_windows.go:350`←宿主半边 `panel_host_windows.go:331-336`←`internal/panel/pump.go:354-373`）；②**被人眼读到＝今天不成立**（`grep -rn wispOnPanelSnapshot frontend/` 我现量零命中 rc=1，含 `dist`；`main.tsx:97` 不传、`:110` 只挂一次）⇒ 票面 AC#5 那句"两半不许合并"我这腿同样不合并。另把 r1 §83 那格没有的**桌面卫生尺**现跑补齐＝`Get-Process -Id 21892` 活着、`StartTime 2026-10-15 19:03:14`、`Path …\build\wisp.exe`（改前改后同一枚存活，本腿全程⛔ 触碰）。⛔ 翻任何 AC 框（九格现量我交件前后各复量一次皆＝9）、⛔ 改 `Status:`／票名／`-done`、⛔ push、⛔ 写面外零字节（本腿只写判语件＋`probes/329/v1/**`＋这一条 bullet）。
