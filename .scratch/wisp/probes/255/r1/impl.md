# 255-r1 交件件（甲形落地＋同批改漂＋无窗仪器）

时刻 `2026-10-08 11:3x +08`。起手 HEAD `e6c3be17`（见 `00-anchor.md`），交件 HEAD `8b32060b`。
⛔ 本件不宣称量过"真浏览器里的宽度"——那一族只有本机可量，且本程⛔ 没开真窗。

---

## 1. 三件各一笔

### 1.1 甲形（`cmd/wisp/panel_host_windows.go`）

新增 `(*PanelManager).requestGeometryOnReshow`，只在 `Show` 的**已建窗分支**（`if !created` 的 `else`）跑一次：
`m.geometry != nil` ⇒ 先取 `m.w`（ nil ⇒ 立刻返回，⛔ 不为一条发不出去的请求读盘）⇒ 调 `m.windowOptions()`
（＝建窗用的同一枚解析函数，height 0→260 与"读不到配置回常量"两条规则只活在一处）⇒
`w.Dispatch(func(){ w.SetSize(width, height, webview2.HintNone) })`。

三条硬约束逐条守：

| 约束 | 这一程怎么落的 |
|---|---|
| ① 外框 vs 客户区 | **具名分开写**，⛔ 不宣称"与配置一致"。建窗那份进 `CreateWindowExW`（`webview.go:296-320`，`0xCF0000 WS_OVERLAPPEDWINDOW`）＝**外框**；`SetSize` 那份过 `AdjustWindowRect`＝**客户区**。同一个 420 是两个宽度。边框差多少＝〔仅本机可量〕，本程⛔ 没造减法（造减法＝拿常识冒充实测）。**具名后果**：一次不动配置的重新显示也可能把外框挪动那个差值；配置读不到时重新显示会把窗口移到宿主常量 420×260——与新建一扇窗今天的结果一致，写在方法注释里而不是让操作者自己推。 |
| ② ⛔ HintFixed | 传 `webview2.HintNone`，两处钉（记录值＋AST 读调用点实参）。`HintFixed` 会清掉 `WSThickFrame｜WSMaximizeBox`（`webview.go:408-413`）＝关掉用户拖窗口边，不在本票射程。 |
| ③ 必须走 Dispatch | 走 `common.go:39 Dispatch`，⛔ 不新造跨线程通道、⛔ 不起 goroutine。仓里对"跨线程直接 `SetWindowPos` 会怎样"**零凭据**（`255-a2` 具名写了没有），所以这一发只说"挂上了库自己的回线程通道"，⛔ 不说"实测过跨线程会崩"。具名后果：`Dispatch` 的队列只有 `Run()` 会排，所以泵交接前的重新显示会等泵起来才改尺寸，不是当场。 |

⛔ 没有热加载 hook：`tiers.go:34` 写着 `"panel":"hot"` 而 `OnReload` 只对 Reload 档开火，`startConfigReload` 唯一非测试调用者在 `run.go`，常驻腿 0 命中 ⇒ 文案与注释一律写"下一次显示请求跟上"，⛔ 没有一处写"立即生效／改完自己变"。

落点：`panel_host_windows.go` `Show` 的 else 支＋一枚新方法；装配根⛔ 未改（`withGeometrySource` 早已存在），
⛔ 未新开 `panel→config` 依赖边（`TestTicket255HostStillDoesNotParseConfigItself` 绿），⛔ 未动 C27。

### 1.2 同批改漂

- `config_readers_255.go` 起手 `:17`：`cmd/wisp/panel_host_windows.go:304-305 at HEAD 67ab595d` ⇒ 换成**内容锚**
  （"the create call inside bringUp spelled its geometry as the literals `Width: 420` / `Height: 260`，去
  `git show 67ab595d:cmd/wisp/panel_host_windows.go` 里 grep 这两枚"），并具名写下为什么⛔ 不留行号
  （那段冷启动尾巴自 33-r5 搬进 `coldStartPageHandover`，33-v4 10-08 点了这次腐坏）。
- `config_readers_255.go` 起手 `:141` 同族措辞一并换成内容锚。
- `config_readers_255.go` `[panel]` 行的两个 bullet ⇒ 三个：读者在常驻进程／关窗再开仍算＋**新增**"活窗由
  Show 的重新显示支发一次 `SetSize`（客户区语义≠外框）"／**新增**"什么还没动：没有那一下按钮，尺寸要等下一次
  显示请求"。verdict 串保留 `关窗再开` 与两枚既有证据锚（`panel_resident_windows.go:207 [cfg.Panel.Width]`、
  `panel_host_windows.go:262 [Width:  uint(width)]`），⛔ 没碰 `AC#1` 那条判据本体。
- `panel_resident_windows.go` 的 `residentPanelGeometryNote`：`state=per-create` ⇒ `state=per-create-and-per-reshow`，
  ⛔ 删掉那句已不成立的"已建好的窗口也不会自己改大小（今天没有 resize 路）"，换成说实话（会发一次／客户区语义／
  没人会在保存那一刻去按）。**该 const 的三行结构保持三行**，所以名册里 `:207`／`:253` 两枚行号锚不位移。
- `panel_geometry_255_test.go` 同票自己的尺里两处**过期陈述**（"there is no resize path today"、
  "this host has no resize path at all"、以及把 `panel_host_windows.go:304` 当现行位置引的那句注释）⇒ 改成
  "只在一次显示请求时发、不在保存那一刻发"＋内容锚。**逐行等价替换，断言逻辑零改动**（枚数 8 增 8 删）。

行号守恒尺（改完复跑）：`panel_resident_windows.go:207` = `return cfg.Panel.Width, cfg.Panel.Height`、
`:253` = `return NewPanelManager(...withGeometrySource...)`、`panel_host_windows.go:262` = `Width:  uint(width),`、
`:392` = `WindowOptions: m.windowOptions(),`——四枚全部仍命中，`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 绿。

### 1.3 无窗仪器（`cmd/wisp/panel_reshow_255r1_windows_test.go`，563 行）

查重（写之前跑的）：`git grep -n "SetSize" -- cmd internal tools scripts` ⇒ 生产码**零枚** `SetSize` 调用点，
唯一命中是 33-r10 假控件的**空体**；`git grep -n "HintNone\|HintFixed\|webview2.Hint"` ⇒ 除那枚空体外**零命中**。
⇒ 没有任何既有钉覆盖"发没发／发的是哪一对／hint 是哪枚／走没走 Dispatch"，本程没有新造同义尺。

⚠ 具名一处：派单说 `panel_pageover_33r10_windows_test.go` 那枚假控件可以直接用——它的 `SetSize` 是**空体、不记录**，
`attachSink33r10` 的形参又是 `*docSink33r10` 具体类型，所以**复用不了**（要它记录就得改 33 族在飞的件）。
本程改为自带一枚 `geoSink255r1`，⛔ 没动 33-r10 的文件。

尺问的是能力，⛔ 不是词面：
- (a) 发出去的是哪一对 ⇒ 记录型 sink 收 `(w,h,hint)`，六档表逐一钉（517×331／640×0→260／0×400→420×400／
  0×0→420×260／-1×-1→420×260／无源⇒一枚都不发），并钉"一次重新显示只读源 1 次"。
- (b) 走没走 `Dispatch` ⇒ sink 把"在 Dispatch 闭包内收到的 SetSize"与"在外面收到的"**分两个名单**，
  外加一枚 `TestTicket255r1SinkTellsDispatchFromABareCall` 自证这枚区分不是恒真。
- (c) hint 是哪枚 ⇒ 记录值 == `HintNone`，**并且**用 AST 读调用点第三实参必须是 `webview2.HintNone`；
  同时钉"整个宿主只允许 1 枚 `SetSize` 调用点"、"它必须裹在 `Dispatch` 的闭包里"。
- (d) 出货入口真能走到吗 ⇒ `TestTicket255r1ShowOnAnExistingWindowSendsTheResize` 直接开 `Show`（无窗）；
  冷启动支跑不了（`bringUp` 要真运行时），所以那一半用 AST 钉"调用只在 `else` 支、且是无条件语句"。

---

## 2. `AC#1` 正控那一发的真实颜色：**绿，且前提没翻**

`go test ./cmd/wisp/ -run 'TestTicket255' -v -count=1` ⇒ `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence`
**PASS**（`logs/test-255-all-rulers:54-55` 抓到进程自己打出的两行，逐字）：

```
IMMEDIATE "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：[llm]"
HONEST    "wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，本次运行不会因此改变行为
           （票 255 AC#1：这一半不许说成「已立即生效」；逐段的读者判定见 HOT-RELOAD-READER 行）：
           [ball] [session] [audio] [agent] [privacy] [memory] [panel] [cost] [models] [observe] [hotkey] [app] [voice]"
```

判读：`AC#1` 那条判据（"那句已立即生效里不许出现 `[panel]`"）**落地后照字面跑仍然成立**——因为 `AC#1` 量的是
`wisp run` 那条 reload 回执，而我这枚新读者活在**常驻 `wisp` 进程**，`wisp run` 根本不建面板宿主；`[panel]` 名册行
仍是 `other-process`。那一格最强的"哪儿都没读者"形态，`255-r3` 早就搬去 `[session]` 了。⇒ **⛔ 这一格不需要为
"前提翻转"重写**；我没有改它、没有放宽它、也没有为好看去动名册文案。

## 3. 门禁 rc 名册（每件自落 rc，⛔ 无管道取退码）

| 门禁 | 命令 | rc | 件 |
|---|---|---|---|
| vet | `go vet ./cmd/wisp/` | **0**（零输出＝干净，rc 自己落进行里） | `logs/vet-cmd-wisp`、`logs/vet-cmd-wisp-final` |
| build | `go build ./...` | **0** | `logs/build-all` |
| 新立名册 | `go test ./cmd/wisp/ -run 'TestTicket255r1' -v -count=1`（PATH 带 `third_party/sherpa-onnx`） | **0**，6 枚函数 / 12 条 PASS 记录全绿 | `logs/test-255r1-newcases` |
| 同票全部尺 | `go test ./cmd/wisp/ -run 'TestTicket255' -v -count=1` | **0**，20 枚顶层全 PASS（含 `ReceiptOmitsPanel...`、`RosterEvidenceLines...`、`RosterStillMatchesTheActualReadSites`、`HostStillDoesNotParseConfigItself`、`PanelRosterVerdictIsTheHonestShape`） | `logs/test-255-all-rulers`、`logs/test-255-final-confirm` |
| d22scan | `sh scripts/d22scan.sh` | **0**，`clean - no D22 ban violations`；正控 `PASS=35 FAIL=0 SKIP=0` | `logs/d22scan` |
| gofmt | `gofmt -l cmd internal tools scripts` | **2**（既有红，见下），⛔ 本程 5 枚文件**一枚都不在名册里** | `logs/gofmt-repo`、`logs/gofmt-targetfiles`、`logs/gofmt-fix` |

⚠ 本机坑第一发就撞上：不带 DLL 直接跑 `go test ./cmd/wisp/` ⇒ `exit status 0xc0000135`、**零 `--- FAIL`**
（＝用例根本没跑）。按 buildrc 分诊后把 `third_party/sherpa-onnx` 挂 PATH 重跑才有颜色（`logs/test-255r1-newcases`
第一版留的就是那枚 0xc0000135 读数，没删）。

### gofmt 那格带着既有红交件（⛔ 不记本程账，逐枚具名）

`gofmt -l` 点名 3 枚 + 1 处解析报错，逐枚出处：

| 名册里的文件 | 是不是本程的 | 现量定性 |
|---|---|---|
| `cmd\wisp\models.go` | 否 | **CRLF 幻影**：`git show d37e75bd:...` 取出的 blob 跑 `gofmt -l` ⇒ 未格式化计数 **0**，工作树那份是行尾符差 |
| `cmd\wisp\panel_inbound_guards_35r3_test.go` | 否 | blob 层就未格式化（anchor 计数 **1**）＝**真红**，属 35r3/35-r5 族 |
| `cmd\wisp\panel_inbound_guards_35r5_test.go` | 否 | `git cat-file -e HEAD:` ⇒ **不在 HEAD**、`git log --` 零命中＝**另一枚在飞腿的未跟踪新件**，本程⛔ 不碰 |
| `internal\panel\inbound_raw_leak_35r7_test.go:270 raw string literal not terminated` | 否 | 35-r7 那枚在飞的件（A705 刚收的族），gofmt 解析报错 |

本程 5 枚文件（`panel_host_windows.go`／`panel_resident_windows.go`／`config_readers_255.go`／
`panel_geometry_255_test.go`／`panel_reshow_255r1_windows_test.go`）单独跑 `gofmt -l` ⇒ 新件先红后
`gofmt -w` 修平、复跑空名册（`logs/gofmt-fix`），另 4 枚始终 0 命中。⛔ 一枚既有红都没顺手改。

## 4. 整包 `go test ./cmd/wisp/`（窗口依赖那族既有红＝不记本程账，附对照实验）

`go test ./cmd/wisp/ -count=1 -timeout 480s` ⇒ **rc=1**，`logs/test-full-package-postchange`。
名册：**5 枚 `--- FAIL`、0 枚 SKIP**，且那一发在 8 分钟帽上 `panic: test timed out`（未跑完，⛔ 别把这份当整包名册读）：

```
TestPanelHostRealWindowHopAndLifecycle              cold bring-up = -1.000 ms（哨兵未替换＝没拿到浏览器往返）
TestAC4FocusReturnToPriorWindowGap33r5              面板没抢到前台（注意：这枚 AC#4 是**票 33 的焦点那条**，与本票 AC#4 无关）
TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe   15s 内页面零回报（A702 具名那枚"字节交对而 #root 恒空"的最坏形状）
TestAC14AwaitedBindingReplyReachesThePage            同上，door 处 "nothing at all"
TestAC14GoSideEvalPushReachesThePage                 同上
```

双向作差本程**没有改前整包基线**（跑一发要另建检出，共享树里⛔ 不许），所以改用一枚更能指向问题的对照：
把那 5 枚单独拎出来跑两次——

| 跑法 | 结果 | 件 |
|---|---|---|
| 带着我的改动，5 枚隔离 | **5 枚全 FAIL**（rc=1，70.4s） | `logs/window5-isolated-withchange` |
| **把 `Show` 里那枚 `m.requestGeometryOnReshow()` 摘掉**，同样 5 枚隔离 | **同样 5 枚 FAIL、同名同耗时**（rc=1，70.5s） | `logs/window5-isolated-withoutmyhop` |

⇒ 摘掉我那发**颜色不变**＝这 5 枚红与本程改动无关，全部属"要真窗＋要 bundle＋要前台权利"那一族，
⛔ 不记本程账、⛔ 不为变绿放宽任何断言。对照做完立刻原样装回，`git diff -- cmd/wisp/` 枚数 **0**
（工作树与 `8b32060b` 逐字节相同），并把 `TestTicket255` 全名册复跑回绿（`logs/test-255-final-confirm` rc=0）。

## 5. 突变表（三发各咬一枚不同的尺，⛔ 没有一枚恒真）

| 突变 | 红了谁 | 件 | rc |
|---|---|---|---|
| A：`w.Dispatch(func(){SetSize})` 改成裸 `w.SetSize(...)` | (b) 两格 + (d) `ShowOnAnExistingWindow` + AST 的"裹在 Dispatch 里"，**共 3 枚函数红**；hint/几何那两格照绿（＝它们问的不是这一格） | `logs/mutation-A-bare-setsize` | 1 |
| B：hint 换 `HintFixed` | (c) 两格：记录值 == HintNone **与** AST 读调用点实参，独立咬；(a) 的 pair 断言照绿 | `logs/mutation-B-hintfixed` | 1 |
| C：把那枚调用从 `else` 支搬到 `if !created` 支 | 分支钉三连红（`!created` 支 1 次≠0／`else` 支 0 次≠1／不是无条件语句）+ `ShowOnAnExistingWindow` 红 | `logs/mutation-C-branch-swap` | 1 |

三发全部当场还原；`grep -n "MUTATION" cmd/wisp/panel_host_windows.go` 与 `git diff` 双尺确认树上没留突变。

## 6. 具名顶回（派单转述与盘上原文冲突处，逐枚）

1. **`config_readers_255.go:141` 并没有引用 `panel_host_windows.go:304`。** 现量
   `grep -n "304" cmd/wisp/config_readers_255.go` ⇒ 起手**只有 `:17` 一枚**命中。起手 `:141` 的原文是
   "`Width: 420 / Height: 260` in cmd/wisp/panel_host_windows.go's create block at HEAD 67ab595d"——它引的是
   **文件＋那段**，没有行号。所以"两处引用已漂行号"实为**一处**。台账 `A700` §5 与票面 `:99` 都写作 `:17/:141`；
   本程按"两处都换成内容锚"执行了，⛔ 没因此改动别的事实，但那一格"枚数"要按原文更正。
2. **`AC#1` 的正控前提**没有**随本票 `AC#4` 翻转。** 现量＝那枚正控用例 `TestTicket255ReceiptOmitsPanelFromTheImmediateSentence`
   在我改动之后 **PASS**，且它抓到的进程自打句子里 `已立即生效` 只列 `[llm]`、`[panel]` 落在那句诚实形里（§2 逐字）。
   根因是 `AC#1` 量的是 `wisp run` 的 reload 回执，而我这枚读者活在**常驻 `wisp` 进程**（`wisp run` 根本不建面板宿主），
   `[panel]` 名册行仍是 `other-process`；另外"哪儿都没读者"那一半 `255-r3` 已搬到 `[session]`。
   ⇒ 那一格的重写**不是**被本票 `AC#4` 翻掉的，按"翻转"去改写 `AC#1` 会把一条仍成立的判据改坏。⛔ 本程一字未动它。
3. **`:304` 在 `67ab595d` 上不是假号。** `git show 67ab595d:cmd/wisp/panel_host_windows.go` 的 304-305 逐字就是
   `Width:  420,`／`Height: 260,`，而原句本来就带 "at HEAD 67ab595d" 限定。所以漂的是**形态**（裸行号会被读成现行位置），
   不是事实。本程仍按令换成内容锚，并把这一点写进注释，⛔ 没把它当"引用了个不存在的行"来修。
4. **33-r10 那枚假控件桩复用不了**（派单说"`:130` 已实现 `SetSize` 空体／`:123` 已实现 `Dispatch`／`:154` 是注入点"，
   三处都真，但空体**不记录**、`attachSink33r10` 形参是具体类型 `*docSink33r10`）⇒ 想拿它问"发的是哪一对"必须改 33 族
   在飞的件。本程自带 sink，⛔ 未碰 `panel_pageover_33r10_windows_test.go`。
5. 一枚同族过期锚，本程**只登记未改**（不在派单三件里）：`panel_resident_windows.go:171` 的注释把装配点写成
   `cmd/wisp/resident_windows.go:142`，现量 `:142` 是注释行、真正的 `newResidentPanelManager(...)` 调用在 `:151`，
   而 `resident_windows.go` 里 `NewPanelManager(` **零命中**（非测试装配点在 `panel_resident_windows.go:253`）。

## 7. 改动枚数

| commit | 文件 | 增/删 |
|---|---|---|
| `d37e75bd`（起手锚，先于任何 go 命令） | `.scratch/wisp/probes/255/r1/00-anchor.md` | 1 file, 46 insertions / 0 deletions |
| `8b32060b`（实现＋仪器） | `cmd/wisp/config_readers_255.go` **35/15**；`cmd/wisp/panel_geometry_255_test.go` **8/8**；`cmd/wisp/panel_host_windows.go` **93/4**；`cmd/wisp/panel_reshow_255r1_windows_test.go`（新，563 行）**563/0**；`cmd/wisp/panel_resident_windows.go` **3/3** | 5 files, 702 insertions / 30 deletions |
| （本件与 logs、票 255 追加节） | `.scratch/wisp/probes/255/r1/**`、票 255 末节 | 见下一笔 commit |

（上表数字＝`git show --numstat 8b32060b` 现量，⛔ 不是 `--stat` 的圆整读数。）

⛔ 一枚 0 字节副产品的具名：`logs/gofmt-anchor-check/panel_inbound_guards_35r5_test.go.anchor` 是空的，
因为 `git show d37e75bd:cmd/wisp/panel_inbound_guards_35r5_test.go` 在起手锚上**没有这个路径**（它是另一枚在飞腿的
未跟踪新件，见 §4 的 gofmt 名册），取回的字节自然是 0——那是**读数本身**而不是交付缺陷，所以本程⛔ 不把它当门禁件交、
也⛔ 不删（临时件只建不删）。同一格 `internal\panel\inbound_raw_leak_35r7_test.go` 的 gofmt 解析报错同理属 35-r7。

票 255 **纯追加一节**：追加前后 `grep -cE '^[[:space:]]*- \[ \]'` 都是 **2**、`- [x]` 都是 **2**（AC 框一字未勾未动）。
TICKET_BOXCOUNT before=2/2 after=2/2
