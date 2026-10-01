# `33-r5` — 面板宿主功能腿：接进常驻＋AC#13 冷启真页面＋AC#4 焦点回还＋AC#14 回执两枚钉＋退净那一支搬 winlive＋两处过期注释

- 程：`33-r5`（**写腿·功能**）｜工单＝票 33 `.scratch/wisp/issues/33-panel-host-c27.md`
- 题面来源＝票面「编排者裁定（10-01 12:12）」九条 ＋「编排者补裁（二）（10-01 13:12）」P1-P5（P1 线程形状／P2 判据形状／P3 STA 显式初始化／P4 缺口另立票 249／P5 dispatchq 不当门）
- 读数来源三件：探针 `.scratch/wisp/probes/33/p1/probe.md` §A／§B；仪器腿 `docs/evidence/s1/33-panel-host-c27-r4.md` §① 格 2／格 3／格 5 与 §⑤ 第 18／19／20 项
- 写面：`cmd/wisp/**`（产码＋测试）＋ `.scratch/wisp/probes/33/r5/**` ＋ 本证据件。⛔ 未动 `internal/ball/**`／`internal/panel/**`／`internal/proc/**`／`internal/observe/**`
- ⛔ 本件**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾归编排者与非实现者验收腿 `33-v2`）；⛔ 不为变绿放宽任何断言、⛔ 不降级成 `t.Logf`

---

## 起手锚点（同发取数，一条命令里跑）

| 尺 | 现量 |
|---|---|
| `date "+%Y-%m-%d %H:%M:%S %z"` | `2026-10-01 13:22:12 +0800` |
| `git log -1 --format=%H` | `7a02b12190322cb2a0937940353e6bdfc123d766` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain` 起手名册 | **229 行**，逐字落 `.scratch/wisp/probes/33/r5/status-start.txt`（⛔ 闸门口径＝**终态名册等于起手名册＋只我这几枚路径**，不写"必须为空"） |
| `git status --porcelain -- cmd internal` | **0 行**（起手即净） |
| 代号复核 | `ls .scratch/wisp/probes/33/` ⇒ `a1 a2 h1 p1 r1 r2 r3 r4` ⇒ **`r5` 未被占用**（起手 `ls .scratch/wisp/probes/33/r5` ⇒ 不存在） |

> ⚠ 派单那句"起手脏项约 226 行"是编排者取数的时刻值；本腿 13:22:12 现量 **229** 行。按记忆里第 77 条那味（共享树里枚数会漂）这里记两枚各带时刻，不判谁错。

---

## 六件（正文在各节；本腿硬顺序＝先交 ⑤⑥⑦ 三节，再填这里）

- [ ] ① 面板接进常驻（`NewPanelManager` 非 test 调用者 0 → N）
- [ ] ② 线程形状＝P1 已裁（专用 STA 线程＋依赖库 `Run()`；⛔ 不投球的 `ui-sta`）
- [ ] ③ AC#13：冷启动那两发 `SetHtml` 不许盖掉真页面
- [ ] ④ AC#4：焦点回还那一跳（＋用例改名 `...Gap33r5`）
- [ ] ⑤ "Destroy 后 2s 退净"那一支移出默认档、进 `winlive`
- [ ] ⑥ 两处过期注释（`:34-38` 的 `PutBounds` 说辞／`:101` 的 `LockOSThread` 幻影指认）

（各件"改前读数—改后读数—反控读数"在下面的 §六件正文，本骨架先交 ⑤⑥⑦。）

---

## ⑤ 我可能写错的条目（对抗我自己）

1. **"接进常驻"可能被我自己写成"接进测试"**。判据是 `grep -rn "NewPanelManager" --include=*.go cmd/wisp | grep -v _test.go` 的**非 test 调用者枚数**。如果我只在新建的 `*_test.go` 里造它，那一枚尺仍然给 0＝本件未落地。⇒ 终态那一把尺逐字落表，⛔ 不引别人表里的数。
2. **`Run()` 形可能被我误当成"回执一定到"**。`33-p1` R26 是在**探针自己的线程**上跑 `Run()` 量到的；本腿把 `Run()` 搬进 `cmd/wisp` 之后是**新的一条线程、新的初始化路径**（`CoInitializeEx(0x2)` 显式、`Registry.Spawn` 的 owner/recover、`bringUp` 在 `Run()` 之前）。如果接完之后回执仍然不到，那是**新证据**，具名上报，⛔ 不许改判据（派单逐字）。
3. **两枚钉可能被我又并成一句**。P2 写死：`Eval` 主动推送与绑定回话是**两维**（R25 vs R27），各要各的凭据。我在代码与用例命名里必须能逐名指出哪一枚钉问哪一维；⛔ 一句 `t.Errorf` 里不许同时塞两维。
4. **AC#13 那枚断言可能是"问 SetHtml 被调用过"的伪牙**。判据是**最终文档含 embed 入口的真内容**，问能力。我的取数路径＝Go 侧 `panel.Assets.Resolve(panel.EntryFile)` 拿字节 → 从字节里抽稳定特征 → 由**页面自己**回答特征在不在（绑定报回 Go）。⛔ 不打开一枚 `frontend/**` 文件（两层禁令＝零读零写零转述；字节只经 `internal/panel` 的 Go API）。风险具名：**DOM 序列化会改写属性引号／补 `html/head/body`**，所以我抽的特征必须是**元素可寻址的东西（id）**，不是原始子串匹配；如果入口里连一枚 id 都没有，我的尺会退化——那一格我按"量不到＝失败测量"处理，⛔ 不当绿。
5. **反控可能在"另一枚口径"下不响**。`git archive HEAD` 抽的仓外副本里 `frontend/dist` 是 **gitignored**（`33-r4` §④#8 现量：入库只 1 枚 `.gitkeep`），那枚副本 `Assets.Built()=false` ⇒ AC#13 那一发只会 skip、不会红。⇒ 反控载体＝**带上工作树 dist 的仓外副本**（先例＝`33-r3` §③"产品树级复现：internal＋go.mod/go.sum＋frontend 复制到仓外"）；搬运是**字节搬运**，⛔ 不打开、不转述任何前端文件。若那一发仍不红，我就具名写"这枚反控在我的载体下量不到"，不假装量过。
6. **我把 `prevFocus` 采样挪到建窗之前，可能顺手改了那四枚断言**。派单逐字：⛔ **四枚断言一字不动**，只改标识符与注释里的归属代号。我对测试文件的改动必须逐枚能说出"这是 setup，不是断言"；断言块若被我的编辑工具碰到，我用 `git diff` 的**删除列**逐名核对并落表。
7. **⚠ 一枚必须让编排者看见的形状改动**：那枚焦点用例今天用**测试 harness**（`hostThreadHarness.bringUp`）先把窗建起来，再 `Hide`→`Show`。产品路径不是这个形（产品的 `Show` 自己建窗）。⇒ 若我把 setup 改成走**常驻线程＋产品的 `Show`**，那是让被测形状等于产品形状，**不是**放松断言；代价＝这条用例从此依赖本腿新建的产码线程。我把这一处单独记在 §六件④ 里，并保留 harness 那两枚既有用例的用法不动，免得 `33-v2` 以为我把"读产码非导出字段"那一格偷偷换成了别的观察口（裁定 5 的耦合我接受）。
8. **前台锁那一枚可能仍然红**（`33-r4` §⑤ 第 5 项：非前台进程时 `SetForegroundWindow` 可能被拒）。派单⛔ 不许加 sleep 重试、⛔ 不许降级。⇒ 若它红，我把**每一发的三枚句柄逐名落表**（before／recorded／afterShow／afterHide＋HEAD＋时刻＋整包或隔离两形）交编排者判，⛔ 不自己宣布"这是环境"。
9. **`winlive` 搬档可能被我搬成"默认档不再有退净判定"而没人知道**。那一支上界仍写 **2 秒**，⛔ 不放宽、⛔ 不删、⛔ 不加重试。代价逐字写进 §六件⑤：`winlive` 在 CI **零岗位** ⇒ 这一支从此〔仅本机可量、CI 永看不见〕。
10. **新线程的名字可能撞上契约面**。`observe.ResidentNames` 六枚名册／`ResidentBaseline` 我**一字未动**；不在名册里的名字落进 `rep.Unknown`（`goroutine.go:270-273` 那一支 `slog.Warn`＝只吵不红）。⇒ 如果哪枚既有门因为我多了一枚 unknown 而红，我**停下上报**，⛔ 不为了"看起来在册"去改名册或改线程名。
11. **退出十步是闭集**（`internal/proc/shutdown.go`／`shutdown_hooks.go`，注释逐字"冻结的是顺序"，可挂名册 8 枚）。我那枚线程的出口今天用 **defer**（`rb.stop()` 同形先例＝票 228 把球那一发挂在 `runResident` 的 defer 上，⛔ 没新增第 31 步）。⇒ 如果我发现"非挂进十步不可"，那就是要动 `internal/proc`＝禁写面 ⇒ **停手上报**，写清要动哪几行。这一格本腿判**不需要**越界，理由与读数在 §六件②。
12. **`Run()` 会占住线程 ⇒ 出货路径的出口可能不干净**。探针 R26 用的是 `w.Terminate()`（投 `WM_QUIT`，库的 `Run()` 在 `WMQuit` 分支 `return`）。我用同一形，并给"线程收干净"一枚**可观察读数**（`done` 通道＋有界等待，等待用的是 monotonic deadline，⛔ 不用墙钟时间差判超时——那是 ban #4）。⚠ 具名风险：`WM_QUIT` 是**线程消息**，如果哪天有人把面板窗挪进球的线程，那枚 `WM_QUIT` 会同时结束球的泵（`33-p1` ⑦-4 已具名这一格，我没测）。
13. **ban #1 的形状**：我新起的协程必须**有 owner ＋ 有 recover**。姿势＝`observe.Registry.Spawn(name, owner, root, fn)`（现成仪器，含 `recover`＋`debug.Stack()`），⛔ 不裸 `go func(`。测试 harness 那一枚既有 `go func()` 有注释 owner＋recover，我不动它。
14. **ban #8 会扫 `_test.go`，注释豁免、字符串字面量不豁免**。射程含 `cmd/`（`33-r4` §③ 现量它就是这样被抓的）。⇒ 本腿所有新增 Go 字符串**只用 ASCII**，⛔ 不写 `⛔`／`✓`／`≤`／任何 U+1F000-U+1FAFF。终态 `sh scripts/d22scan.sh` 必须 rc=0。
15. **`go mod tidy` 一字节不许跑**（它在 HEAD 上 exit 1，会造出不属本腿的 diff）。我用的是库里**已有**的 webview2 依赖，⛔ 零新增依赖、⛔ 不改 `go.mod`／`go.sum`。
16. **桌面**：真窗起完就收、一次一枚。⚠ 我这发会新增**会起真窗的用例**（AC#13／AC#14／线程那一族），默认档的 `cmd/wisp` 整包时间会涨；`msedgewebview2` 枚数起手／终态各量一次（`33-p1` R30 那把尺），⛔ 不给这台机器留孤儿子进程。
17. **⚠ 撞钉预检的漏计风险（记忆里第 64／70 条那一味）**。我起手跑过这些尺并把读数抄在下面 §⑥，但"今天绿的用例名册"要**整包跑一遍才算量到**；如果我在写完之前没跑起跑名册，那我可能在改语义时撞红别家的钉（本票已知的三枚：`internal/panel/composer_dispatch_test.go` 的反转钉（要求"必须有宿主"，本腿只会让它更成立）、`cmd/wisp/resident_ball_228_test.go` F1 的十枚回调名册（我只改 body、保留键名）、`cmd/wisp/panel_host_gate_test.go` 的 dispose AST 尺（裁定 4 要求本腿**同时收紧**它，见 §六件①）。
18. **⛔ 我没有做的事**（免得被读成做了）：没动 `internal/**` 一字；没翻票面任何一枚框；没跑 `go mod tidy`／`go get`；没读没写 `frontend/**`／`design/**`；没动 `PLAN.md`／`specs`／`SLO.md`／`BUILD.md`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件；没加 `Allow` 方法、没造第 5 枚 veto 通道、没造"页面点一下→Go 判成 L2 批准"的任何变体；没有把 artifacts 写入做成受门控的 Tool。
19. **可能被我写成"面板能用了"的两处夸大**。本腿落地之后仍然**不成立**的三件，我在结论里逐字保留：ⓐ 页面侧那一腿（`frontend/**`，票 248 与界面侧会话）今天仍没有人点；ⓑ 入向白名单只有四枚方法、**零枚 config/凭据方法** ⇒ owner 那句"点设置自己录 key"仍差票 248；ⓒ AC#12 的"有没有一包真页面"归口未落（工作树有产物、入库只 `.gitkeep`）。⇒ "用户能打开面板"这一句的成立条件是**上面三件之外都齐**，我在 §六件① 里按"宿主可用／用户可打开"两档分开写。
20. **一处归属卫生**：本件的行号一律本腿自己 `grep -n` 现取（⛔ 不照抄派单或前人表里的行号；引用别人读数时写明"出处＝谁的哪一节，本腿未复跑"）。

---

## ⑥ 我跑了哪些尺（逐条真实读数；起手与预检先登，真跑读数续编号）

| # | 命令（完整） | 逐字读数（摘要） |
|---|---|---|
| S1 | `ls .scratch/wisp/probes/33/` | `a1 a2 h1 p1 r1 r2 r3 r4` ⇒ **`r5` 不在名册**（撞名即停的预检，通过） |
| S2 | `date; git log -1 --format=%H; git rev-parse --abbrev-ref HEAD` | `2026-10-01 13:22:12 +0800`／`7a02b12190322cb2a0937940353e6bdfc123d766`／`dev`（同发，落 `logs/anchor-start.txt`） |
| S3 | `git status --porcelain \| wc -l`；`git status --porcelain -- cmd internal \| wc -l` | **229**；**0**（名册逐字＝`status-start.txt`） |
| S4 | `grep -rn "NewPanelManager" --include=*.go cmd/wisp \| grep -v _test.go`（起手那发） | 见 §六件① 的"改前读数"表（本腿自己现取，⛔ 不引票面 `:132` 那句） |
| S5 | `grep -n "CoInitialize" cmd/wisp/*.go`（排 test 另尺） | 见 §六件② 改前读数（派单那句"今天零枚"本腿复跑，⛔ 不采信转述） |
| S6 | 撞钉预检：`grep -rn "no panel host\|panel-hotkey\|tray-open-panel" --include=*.go cmd internal` | 命中 4 枚，全在 `cmd/wisp`（`resident_ball_windows.go:22/:127/:128/:203`）＋`approval_reply.go:428` 一枚注释 ⇒ **没有任何一枚断言钉死那句"no panel host"文案**；`ballGestureWhy` 那句改口不会撞词面钉 |
| S7 | 撞钉预检：`sed -n '40,60p;300,340p' cmd/wisp/resident_ball_228_test.go` | F1＝**词面型名册钉**：`ballEventCallbacks228` 十枚**回调键名**，判的是"宿主有没有把每一枚键设上"（`len(missing)`／`len(set) > len(list)`）。⇒ 本腿只改 `OnPanelHotkey`／`OnTrayPanel` 的 **body**，两枚键名保留 ⇒ 不红 |
| S8 | 撞钉预检：`grep -rn "runResident\|startResidentBall(" --include=*.go cmd/wisp` | `startResidentBall` 既有**三枚** live 用例＋一枚 `resident_approval_246_windows_test.go:314` 两枚参数调用 ⇒ 我改签名会撞编译；本腿**用变参**（`hooks ...ballHostHook`）保两枚参数形仍然合法，⛔ 不改任何既有用例的一行 |
| S9 | `grep -n "func (r \*Registry) Spawn" internal/observe/goroutine.go` | `:262 Spawn(name, owner string, root *Root, fn func(ctx context.Context)) *Handle` ⇒ 现成的 owner＋recover 形状，本腿用它是 ban #1 的正解（⛔ 不裸 `go func(`） |
| S10 | `grep -n "ResidentNames\|ResidentOverBaseline\|CategoryUnknown" internal/observe/*.go`（排 test） | `:35 CategoryUnknown`／`:43 ResidentNames`／`:65 slices.Contains`／`:81 return CategoryUnknown`／`:195 ResidentOverBaseline`／`:270 if cat == CategoryUnknown`／`:422 rep.ResidentOverBaseline = rep.Resident > ResidentBaseline` ⇒ 名册与基线我**一字未动**；新线程名落 `Unknown` |
| S11 | 必读三件全文已读 | 票面裁定九条＋补裁（二）P1-P5／`probes/33/p1/probe.md` 263 行（§A／§B）／`docs/evidence/s1/33-panel-host-c27-r4.md` 307 行。**引用它们的地方都标出处，⛔ 没把任何一枚前人读数当本腿读数** |
| S12 | `wc -l` 被测四件 | `cmd/wisp/panel_host_windows.go` **411 行**、`panel_host_windows_test.go` **723 行**、`resident_windows.go` 205 行、`resident_ball_windows.go` 271 行（全文逐行读过） |

（续编号 R1.. 留给起跑名册、终跑名册、门禁四数、反控读数。）

---

## ⑦ 判不动的地方（逐条 甲／乙／不做 ＋ 现量 ＋ 为什么判不了）

1. **球侧那一格：面板线程与球的 `ui-sta` 同时活着时，球的出帧会不会被插走**。现量：`33-p1` §⑦-结-2 逐字"我只证了 `wmAppTask` 这一类消息的顺序，D2D 的 `WM_TIMER` 出帧一枚都没测"；裁定 P1 也写明那一格〔没测〕。本腿**不往球的线程上投任何东西**（P1 的裁形），所以这一维在本腿的改动下**不新增暴露面**，但我⛔ 不能说它安全。甲＝具名交回；乙＝去动 `internal/ball/**`（禁写，归票 228）；**本腿做甲**。
2. **`Run()` 与别的泵在同一枚线程上共存**。现量：`33-p1` §⑦-结-4〔推〕。本腿的形状＝**面板线程只有 `Run()` 一枚泵**（`bringUp` 里那发 `pnlPumpOnce` 自泵在 `Run()` **之前**、同一条线程、栈不重叠），所以这一格在本腿**不需要**被裁；但"同一线程两枚泵"这一形我没有读数。⛔ 不许把它读成"已证共存无事"。
3. **MTA 那发的 `res` 真值**。现量：`33-p1` §⑦-结-3＝拿不到（要自己实现那枚未导出的 environment handler）。本腿只钉"显式 STA(`0x2`) ＋ 初始化返回值必须查"，⛔ 不写"STA 就一定成功"这种因果（P3 逐字：成因未定值、只定了现象）。
4. **`dispatchq` 只涨不取那一维**。现量：`33-p1` R25 间接读数（4 枚 `WM_APP` 被取走、0 枚闭包执行），直接长度未导出。裁定 **P5＝记欠账、不当门**，且 P1 选定 `Run()` 之后这一维由投递者本身消掉。⇒ 本腿⛔ 不为此单开闸门，只在表里说"今天的形状是 `Run()` 在取"。
5. **CI 那枚 windows cli 腿上有没有 WebView2 Runtime**。现量：未量（那是 `33-v2` 的格，票面 `:310` 逐字"判法＝读 `.github/workflows/ci.yml` 的 runner 标签＋那一步的命令"）。本腿新增的真窗用例**沿用既有默认档形状**（`//go:build windows`、起不来＝红不是 skip），所以这一格的暴露面**随本腿变大**（多几枚真窗用例）。⇒ 甲＝具名上报"我把默认档真窗用例从 N 枚涨到 M 枚"（读数在 §六件终态）；乙＝自己把它们搬进 `winlive`（⛔ 那是编排者的档级裁定，撤销口令都写好了的那一枚只有退净那一支归我）；本腿做**甲**。
6. **"焦点回还失败"该不该由产品上报**。现量：`Hide` 今天只有 `:311` 那一句 `if prev != 0`，`SetForegroundWindow` 的返回值**被丢掉**；前台锁（`33-r4` §⑤#5）会让那一跳真的失败。⇒ 修法要不要"检查返回值并响亮报出"属产品形状；本腿**加读数、不改判据**（把返回值取进日志，⛔ 不据此放宽断言④）。如果编排者要我把它升成断言，那是**新裁**，我不自己定。
7. **`Assets.Check()` 那一层的牙不在本腿射程**。现量：`33-r4` §① 格 5 逐字"本尺把半包的完备性判定委托给 `Assets.Check()`"，而 `internal/panel/**` 是本腿禁写面。⇒ AC#13 那枚"最终文档含真内容"的断言问的是**送进窗口的那包字节**，⛔ 不背书那包字节自洽。
8. **`go mod tidy -diff` 在 HEAD 上 exit 1／`staticcheck` 本机版与 CI 钉版不同**。现量：票面 `:267`／`:262`（本腿未复跑，⛔ 不跑）。⇒ 门禁四数里 staticcheck 标〔未复认〕，tidy 一字节不跑。
9. **"接进常驻"接到哪一枚触发口**。现量：`internal/ball` 的 `OnPanelHotkey`／`OnTrayPanel` 两枚回调今天只 `recordBallGesture`（`resident_ball_windows.go:127-128`），包外**没有**任何导出投递口（`33-a2` R20）。⇒ 本腿的接法＝在这两枚回调的 **body** 里投给面板线程（回调契约"快、非阻塞"＝一次 channel send／`Dispatch`，⛔ 不等窗建好）。⚠ 我没有真按过热键（那要桌面注入，`33-p1` §⑦-1 已判"做不到可复核前提"）⇒ 我这一格交的是**代码路径＋用例读数**，⛔ 不交"真按 Ctrl+Alt+P 看到了窗"。那一格留给 owner 手测或 `33-v2`。
