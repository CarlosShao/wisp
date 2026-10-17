# 票 329 非实现者终裁 — 腿 `329-v1`（九格逐格判语，2026-10-17 10:3x +08）

裁决者＝`329-v1`（⛔ 实现者、⛔ 产码、⛔ 翻框）；被裁对象＝票 `.scratch/wisp/issues/329-panel-snapshot-last-mile-never-delivered.md` 九格
＋落地腿 `329-r1` 的产码七笔（`60e82f11`／`bed1c348`／`a0a2f440`／`8e1ea71d`／`2c4f829a`／`1d1261e0`／`a0d6ec85`）。
判语料只有我自己跑出来的读数；`329-r1` 件的读数我逐枚复量过，凡是它没跑的我一律自己跑（下面点名）。

**台面**：两枚仓外导出树＝`git archive` 八件套＋仓根 `docs/`＋手工三件；`D:/W329VA`＝起手锚 `feab45e3`（747 枚 `.go`）、
`D:/W329VB`＝被裁尖 `1d1261e0`（750 枚 `.go`，差＝本票三枚新件）；`GOTMPDIR=D:/wisp-gotmp-329v1`（先 `mkdir -p`）。
树 B 三枚被裁件的起手 sha256＝`panel_snapshot_deliver_329.go` `e7a11ceb0403…`／`internal/panel/pump.go` `728891008927…`／
`cmd/wisp/panel_pump.go` `3255fd96fff6…`，与 r1 件 `logs/pristine-hashes-B-resync.log` 三值逐枚相同 ⇒ **我这棵树与它那棵逐字节同形**（一枚交叉印证）。
`cmd/wisp` 一律带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` ＋ `-count=1 -timeout 30m`；⛔ `go test ./...`。
开工前后各一次现跑 `Get-CimInstance … CommandLine -match 'go test'`：本腿全程**没有第二枚 `cmd/wisp` 在飞**，故"重包排队"一格今天未触发（`228-r3` 此刻在共享树里只有 `internal/panel/display_scrub_327r2_test.go` 一枚未跟踪件，与本腿读数无关，⛔ 当凭据）。

## 1. 九格逐格

**AC#0 交件三把尺＝不成立（两处缺读数＋一枚指向不存在件的引用；缺陷全在"件"上，不在"字节"上）**
- 名册我自己复跑（`git show --numstat`，逐枚见 `logs/ac0-roster-and-hygiene.log`）：七笔**只**落 `cmd/wisp/**`（2 枚新件＋4 枚改）与 `.scratch/wisp/probes/329/**`（18 枚）；
  逐枚过滤 `issues/`⇄`evidence/`＝**七笔全 0**（票尾 bullet 与证据件由编排者 `7535d83b` 代提）。`cmd/wisp` 全量 `feab45e3..1d1261e0`＝364 插／18 删、6 枚文件。
- 缺陷甲＝r1 件 §1 指名 `probes/329/r1/logs/commits-roster.log` 作凭据，我现跑 `ls` **rc=2 不在盘上**（编排者 §8.2 已具名登记；本腿复量为"确不存在"）。名册本身可由 `git show` 现跑复得＝上面那条。
- 缺陷乙＝r1 件 §6 表第 68／69 行（改前·skip3 ⇄ 改后·skip3）与 §83"桌面卫生尺"三格**没有读数**（原词本件不复述）。⇒ 按 AC#0 自身那句"空读数词尺＝0"这一条**没满足**。
  两处我都另跑补齐：颜色见 AC#7、卫生见 §3。补齐≠替它勾框：AC#0 这一格判**不成立**，建议按"缺陷已具名、料已补全，翻框归编排者"处理。

**AC#1 起手三把尺＝成立**（我自己在树 A 重跑，⛔ 抄票面）
- 尺①：`grep -rn 'lastPanelSnapshot' --include='*.go' internal cmd | grep -v '_test.go' | grep -v '^cmd/wisp/panel_pump.go:4'` 树 A＝**0 命中**（管道 rc=1）。
  正控两发同尺：不滤 `panel_pump.go:4`＝2 命中（定义 `:424`＋注释 `:417`）；换一枚已知存在的名字 `bookPanelSnapshot`＝3 命中（`panel_pump.go:331`/`:360`＋`run.go:732`）。
- 尺②：`publishPanelSnapshot` 生产挂载在（`run.go:742` `rt.ui.publish`、`run.go:1063` `consoleSink.publish`）⇒ 票面"⛔ 把本票读成整条没接"这句成立。
- 尺③：`bridge.go` 入向名＝**7**（逐枚点名 `panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`/`config.get`/`config.set`/`config.models.list`）；
  `wantFullInboundRosterSize253 = 7` 两树同值（`inbound_roster_253_test.go:63`）。
- 顺带核正票面 §现量表一处数：票写"16 枚命中"，我树 A 现量＝**总 17**（非测试 2＋测试 15，逐枚 `3/3/3/2/2/2/2`… 实为 6 枚测试文件求和 15）；
  r1 §2 的登记句把求和写成"3+3+3+2+2+2+2＝17"（七项对六枚文件）＝**结论对、算式歪**，两处都只登记、⛔ 改原文。

**AC#2 落地那一跳＝成立**（"有非测试读者"⇄"字节逐位相等"两半都我自己跑到）
- 读者尺（同一把，树 B）＝3 命中：`panel_snapshot_deliver_329.go:9`/`:83` 两枚注释＋**`:95` 真读者**（`_, data, seen := rt.lastPanelSnapshot()`）。
- 逐位相等那半＝AC#3 三发突变背书（下格），⛔ 我这里复述结论。

**AC#3 逐位相等，不是"调用过"＝成立**（★本腿自己重放，四发＋一发正控，全落仓外树 B 物理字节，逐发还原、还原哈希逐发对回 §起手三值、`MUT329V1` marker 残留＝0）
- 绿基线（定向 `TestBookedSnapshotReachesHostBitForBit329`，`cmd/wisp/panel_snapshot_deliver_329_test.go:76`）＝`--- PASS (1.73s)`，四数 1/1/0/0。
- **去牙（空串形）**＝推 `panelSnapshotEvalJS(data[:0])` ⇒ `--- FAIL (1.33s)`，红句**只在** `:141`（push-side）＋`:164`（closing），`host="window.wispOnPanelSnapshot()"`；
  ★同发 `:134`（marshal-side）**保持绿** ⇒ 红的是"推出去的那串字节"，不是记账本身。
  ⚠ 具名差一笔：r1 §3 把这发写成"推 `panelSnapshotEvalJS(nil)`"，**那一形在 `1d1261e0` 字节上编不过**＝`cmd\wisp\panel_snapshot_deliver_329.go:95:5: declared and not used: data`＋`[build failed]`（逐字在 `logs/mutation-v1-defang-empty.log`）；
  build-failed 属台面故障、不是用例红 ⇒ 它的红句要成立只能是"字节换成空切片"那一形（我复现的那形红句与它盘上 `mutation-1-defang.log` 逐字同形状）。判：**实质成立、措辞不实**，AC#0 缺陷甲同批登记。
- **去牙（少最后一字节）**＝`data[:len(data)-1]` ⇒ `--- FAIL (1.47s)`，红同两行。
- **去牙（等长改内容）**＝把 payload 首字节翻一位、**长度一字不动** ⇒ `--- FAIL (1.45s)`，红 `:141`＋`:164` ⇒ 正面回答票面 AC#3 那句"证明它真在比字节、不是比长度"（r1 三发里没有这一发，本腿补）。
- **反控（marshal 侧少一条待批）**＝`internal/panel/pump.go:359` 改为对 trimmed 副本 marshal ⇒ `--- FAIL (1.47s)`，红句 **lead＝`:134`（marshal-side，另一枚尺、另一行）**，其后 `:141`/`:160`/`:164` 连带，`booked="{\"pending\":null,…"` ⇒ 与去牙两发**落点不同**（成对性成立）。
- **必答题（摘掉生产读者 `:95`）**＝`var data []byte; seen := false` ⇒ `--- FAIL (1.46s)`，红**逐枚落 `:141`／`:164`／`:167`**、`host=""`、"the stub host received 0 pushes" ⇒ 与盘上 r1 `mutation-3-reader-removal.log` 三行同形（我这次是自己跑的，编排者 §8.5 那一格读的是原文）。
- **还原正控**＝四发之后整树定向复发 `--- PASS (1.47s)` rc=0；逐发 `RESTORE deliver329=e7a11ceb04…match_pristine=YES`、`pump=72889100…YES`（`logs/mutation-hashes.log`）。

**AC#4 名册钉＝成立**
- `git diff --stat feab45e3 1d1261e0 -- internal/panel/bridge.go internal/panel/inbound_roster_253_test.go internal/panel/git_test.go`＝**空**（rc=0）。
- 入向名册现量仍＝7（两树同值，见 AC#1 尺③）；⛔ 为过关改名册数／加白名单：无。
- 出向那四枚 `window._rpc` 回执一字未动＝`git diff --numstat feab45e3 1d1261e0 -- cmd/wisp/panel_control_windows.go`＝**空**，
  `:737`/`:744`/`:748`/`:752` 两树逐字节同（我两边都 grep 过）。本票新用的 `Eval` 是**第五枚**使用者，挂在 `PanelManager.evalPanelJS`（`panel_host_windows.go:331-336`），不碰那四枚。

**AC#5 页面侧合同件＝成立（合同与产码逐字一致；一处行号漂移不实质）**
- 全局名＝合同写 `window.wispOnPanelSnapshot` ⇄ 产码 `panelSnapshotGlobal`（`panel_snapshot_deliver_329.go:57`）＋包裹 `:65-67`（`global + "(" + string(data) + ")"`，原文进参、不转义不重编码）＝**一致**。
- 参数形状＝合同四必选键 ⇄ `internal/panel/composer.go:57-61`（`pending`/`results`/`composer`/`generatedAt`）；两可选键＝`:74` `instructions,omitempty`、`:91` `tasks,omitempty`＝**一致**；
  卡面十字段 ⇄ `internal/panel/approval.go:40-58`（含 `:43` 注释的 `"L0"|"L1"|"L2"|"Deny"`、`:58` `decidedBy` 恒 "native"（`:87` 赋死）、`callChain` 今天恒空＝`internal/panel/pump.go:76` 具名"无生产者"）＝**一致**。
- 触发时机＝"每次 publish 恰好一发、⛔ 逐 token"⇄ `cmd/wisp/panel_pump.go:442`（`// Not on streamed text deltas.` 仍在）＋`:447-463`（推跳只在 publish 里）＝**一致**；
  窗口未开＝不推 ⇄ `panel_host_windows.go:332-335` nil 分支；控制台 `wisp run`＝零推 ⇄ 挂载只在 `resident_windows.go:348-350` 的常驻装配根。
- 一处漂移：合同 §4 引 `lib/panel.ts` 的 `hostBridge()` 在 `:161`，我现量＝**声明在 `:159`**（`:161` 是函数体里那行 `return window.wispBridge ?? …`）＝行号取自票面 10-16 读数、漂移 2 行，实质（"全是页面→Go 的门，没有 Go→页面的接收器"）我自己复量成立。
- ⚠ 定位一句（⛔ 任何人把它读成能力）：这枚件是**给机主转交外部前端 agent 的料**，不是本仓已实现能力；它 §4 自己写着"页面接收器今天缺位"。

**AC#6 ⛔ 发明合并／背压＝成立**
- 驱动点数与形**一枚没多**：`c.publish()`/`u.publish()` 两树同为 3 枚调用（A `:1363`/`:1452`/`:1481` ⇄ B `:1375`/`:1464`/`:1493`，只随本腿字段漂移行号），
  所在函数同名同形（`consoleSink.Publish`、`consoleApprovalUI.Prompt`、`consoleApprovalUI.Update`）；挂载两枚（`run.go:754`、`run.go:1075`）。
- 推跳只挂在 `publishPanelSnapshot` 内部一处（`panel_pump.go:463`）＝票面裁的"沿用既有 publish 真发生的那一刻"，⛔ 新驱动点。
- 逐 token delta 推＝无（`panel_pump.go:442` 那句还在；我四发红句里 `每个 delta`／`per delta` 零命中，见 `logs/mutation-*.log`）。
- 合并／背压规则＝一件没发明（`deliver_329.go:18-21` 与 `panel_pump.go:456-462` 两处都明文把 D38d 那条"多久推一次"留在票 35 名下）。⇒ 无需"停下来上报"那一格。

**AC#7 台面与颜色＝不可判（缺的正是"整包终态红名册"那一维；另路那一维我已跑到）**
- 命令形与台面我全部复现：八件套导出树＋`docs/`＋手工三件、`D:` 卷短名目录、`GOTMPDIR` 先建、DLL 路径＋`-timeout 30m`、⛔ `go test ./...`。这些**成立**。
- 整包那一维：r1 六发我逐枚自己 `grep -c` 复量（`logs/../r1/logs/full-*.log`）＝`166/115/2/0`、`203/148/0/0`（run2）、`203/148/0/0`（skipAC13）、
  `205/149/1/0`（skipAC13pair）、`203/147/1/0`（skipset2）、**改后 tip 树 `206/151/0/0`**；rc 全 1、`^ok` 全 0 枚；
  我另数了两发死亡发的形状：末枚 `=== RUN`＝A 侧 `TestAC13BringUpSurvivesAReusedThreadQuit`（log 第 1295 行）、B 侧 `TestBallPanelGesturesReachThePanelThread`（第 1361 行），两发尾行都停在 `MARK299 wndproc … known=true` 之后（第 1305-1306／1372-1373 行）
  ⇒ 形状＝**整进程在真窗区死亡**（票 314 现役症状、与 329 的字节无关）。⇒ ⛔ 读成"全绿"，⛔ 读成"329 把包打红"。
  我**没有**重跑整包赌一把（派单点名禁止；且共享树里 `cmd/wisp` 与 `internal/panel` 是别的腿的写面）。
- 另路取数（我自己跑的那一发）＝**具名 13 枚 `-run` 名册**（`lastPanelSnapshot`/`publishPanelSnapshot` 的全部测试读者所在文件，逐枚点名在 `logs/roster-named-*.log`）：
  起手树 A＝`RUN 12／PASS_TOP 12／FAIL_TOP 0／SKIP 0`、rc=0（第 13 枚＝新用例，A 侧不存在）；交件树 B＝`RUN 13／PASS_TOP 13／FAIL_TOP 0／SKIP 0`、rc=0。
  红名册两侧皆空 ⇒ **作差＝0**；名册差枚＝唯一新增 `TestBookedSnapshotReachesHostBitForBit329`，其余 12 枚两侧同名同绿。
- 判定＝**不可判**，缺的那一维写清楚：不是"329 的红不红"，而是"**整包终态红名册 ⇄ 起手红名册的作差**"这一维今天在这台机器上任何腿都取不到（死亡点逐发漂移、六发无一到达终点锚 `TestZZ315A1*`）。
  本票可触达面（那 13 枚）已由另路那一格判为作差 0；把它当成 AC#7 全格成立＝越界，我不做。

**AC#8 门禁四数＝成立（本腿自己在共享工作树跑的四数）**
- `sh scripts/d22scan.sh`＝**rc 0**、`d22scan: clean`（分母 internal=240／cmd=42／frontend=116／tools=25／design=39）。
- `sh scripts/check-path-length-budget.sh`＝**rc 0**、`VERDICT GREEN`（tracked=14579、over=60、roster=60、名册外 0）。
- `go vet ./cmd/wisp/`＝**rc 0**（共享树，量的是含别腿在飞字节的那棵树，⛔ 用作本票凭据之外的东西）。
- `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp`＝**rc 0**、列两枚＝`cmd\wisp\models.go`＋`cmd\wisp\panel_transport_35r2_test.go`——我这一跑**同时复核了 r1 §5 那一行**：两枚都不是 329 的写面＝既有欠账；本腿写面零枚 `.go`，故对我是 n/a（首轮 rc=2 只因判语件当时还没落盘）。

## 2. 能力必答（这一跳今天到底被谁读到）

- **① 被产码读到＝成立**，逐环 `file:line`（全在 `1d1261e0` 具名树）：
  `cmd/wisp/panel_snapshot_deliver_329.go:95`（第一枚生产读者）→ `:105`（`h.evalPanelJS(panelSnapshotEvalJS(data))`）← 调用者 `cmd/wisp/panel_pump.go:463`（`publishPanelSnapshot` 记账之后同发）
  ← 三个既有驱动点 `cmd/wisp/run.go:1375`（`consoleSink.Publish`）／`:1464`（`consoleApprovalUI.Prompt`）／`:1493`（`consoleApprovalUI.Update`）
  ← 挂载 `run.go:754`（`rt.ui.publish`）＋`run.go:1075`（`consoleSink{… publish:}`）← 装配根 `cmd/wisp/resident_windows.go:350`（`src.run.setPanelSnapshotHost(rp)`，`src==nil`／`src.run==nil`／`rp==nil` 三形＝不挂）
  ← 宿主半边 `cmd/wisp/panel_host_windows.go:331-336`（`currentWindow()` 非 nil 才 `w.Dispatch(func(){ w.Eval(js) })`）→ 泵 marshal `internal/panel/pump.go:354-373`。
- **② 被人眼读到＝今天不成立**，我自己现量的尺：`grep -rn 'wispOnPanelSnapshot' frontend/`＝**零命中**（rc=1，含 `frontend/dist`）；
  `frontend/src/main.tsx:97` 的 `product: {… node: <App />}` 不传 snapshot、`:110` 只 `createRoot(...).render(...)` 一次；`lib/panel.ts` 里只有页面→Go 的门（`:155` `wispDispatch`、`:159-161` `hostBridge`）。
  ⇒ 今天的事实是"Go 把一枚 JS 语句交给宿主，页面里没有那枚全局，调用被 WebView2 静默丢弃"。**⛔ 那一屏今天不跟着任务动**＝票面 AC#5 明文禁止合并的两半，本件一处也没合并。
  ⇒ 这一半**不是本票的失守**：票面 AC#5 与合同件 §4 明文禁止把它合并成端到端；本腿只是把"缺位"钉成我自己的现量而非转述。

## 3. 死腿三处读数的处置（⛔ 抹它原文，本件只另记）

1. §表 68／69 行的空读数 → 我用具名 `-run` 那一格补齐（AC#7），并把"整包那一维"判成不可判而非替它圆。
2. §83 桌面卫生尺的空读数 → 我现跑：`Get-Process -Id 21892`＝`ProcessName=wisp`、`StartTime=2026-10-15 19:03:14`、`Path=…\build\wisp.exe`（`logs/ac0-roster-and-hygiene.log`）
   ⇒ 活实例完好＝**本腿与 r1 都没杀它**（改前改后同一枚，跨 10-15→10-17 存活，可作两台面的同形凭据）；我全程⛔ 动它、⛔ 为它放宽任何断言。
3. §1 那枚不存在日志的引用 → 复量 rc=2 确不存在；名册我由 `git show` 现跑复得（§AC#0）。
4. 新抓一处：§3 的"nil 形"在 `1d1261e0` 字节上 build failed ⇒ 那发红句只能由空切片那一形产生（逐字见 `logs/mutation-v1-defang-empty.log` ⇄ `-r2.log`）。

## 4. 交件自证（三把尺，⛔ "已完成"）

- **① 逐枚提交＋逐枚 `git show --name-only`**（本腿写面只有三处：`docs/evidence/s1/329-v1-verdict.md`、`.scratch/wisp/probes/329/v1/**`、票 329 末尾一条 bullet）：
  `c27d5909`＝起手尺件×2＋回放脚本＋msg-1（4 枚）；`4b424368`＝突变四发读数＋哈希汇总＋还原正控＋进度件＋msg-2（9 枚）；`249786fd`＝具名 `-run` 名册两发＋门禁件＋AC#0/卫生件＋msg-3（5 枚）；
  其后两笔＝本判语件、票尾 bullet。⛔ push（推送只由编排者做）。
- **② 本件**：绝对路径 `D:\work\workspace\projects plans\Wisp\docs\evidence\s1\329-v1-verdict.md`；`wc -c`／行数／空读数词尺（词集＝`logs/verdict-selfcheck.log` 里那行 grep 式，本件不自带那些词）＝交件前现跑，逐值落同一枚 `logs/verdict-selfcheck.log`；0 字节件＝0 枚。
- **③ 决定性尺的逐字**：起手侧＝树 A `lastPanelSnapshot` 非测试读者 0 命中（管道 rc=1）＋正控 `bookPanelSnapshot` 3 命中；交件侧＝树 B 同尺 3 命中、真读者 `cmd/wisp/panel_snapshot_deliver_329.go:95`；
  成对红句逐字（`:134` ⇄ `:141`＋`:164`）、等长改内容发＝`:141`＋`:164`、必答题＝`:141`/`:164`/`:167`＋`host=""`＋"received 0 pushes"；
  还原哈希四发全部 `e7a11ceb0403…`／`728891008927…`（＝起手值）；另路名册四数＝A `12/12/0/0`、B `13/13/0/0`；门禁四数 rc＝0/0/0/0。
- **必答题（派单问的，不是我盖章的）**：摘掉生产读者 `:95` 后同一枚用例**响**，红落 `:141`/`:164`/`:167` ⇒ 该尺对这一跳敏感＝**是**，且我这一发是自己跑的。
- **给编排者的三句话**：AC#0 判不成立（料可补、字节无罪）；AC#7 判不可判（缺整包终态那一维，⛔ 任何腿今天取不到）；其余六格＋能力①成立、能力②今天不成立（属预期缺口）。
