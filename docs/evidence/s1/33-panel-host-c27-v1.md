# `33-v1` — 面板宿主真起来：非实现者对抗验收表（射程＝票 33 的 AC#1／AC#2／AC#3／AC#4／AC#11／AC#12）

- 程：`33-v1`（**对抗验收腿，非实现者**）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）
- 被验腿：`33-r1`（写腿），五枚提交 `e4ed8304`→`697b4fae`→`b1e3d94a`→`fafe0445`→`f0ad3b51`
- ⛔ 它的证据件 `docs/evidence/s1/33-panel-host-c27-r1.md` 属**〔实现方自述＋台件〕**，**不构成任何一格 AC 的翻勾依据**；本表每一格只用我自己现跑的尺。
- ⛔ 本表**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾框归编排者）；⛔ 不写"附条件通过"式补判语。
- 写作顺序（硬要求）：**§A／§B／§C 三节自对抗内容先写满并 commit**，之后才回填 §D 逐格判语。本表判据交付＝`wc -l -c` 这个文件。

---

## 0. 起手同发读数（`date`＋HEAD＋脏项，一次取）

| 尺 | 现量 |
|---|---|
| `date` | `Thu Oct  1 10:35:27 CST 2026` |
| `git rev-parse --short HEAD` | `94bd133e`（**起手锚**；HEAD 会动——只读普查腿 `244-a3` 在往 `.scratch/wisp/probes/244/a3/census.md` 写文档件，所以收尾时 `git diff 94bd133e..HEAD -- cmd internal` 必须为空） |
| `git status --porcelain -- cmd internal` | **0 行**（起手名册＝空；终态必须等于这一枚名册） |

---

## A. 我跑了哪些尺、每条真实读数（逐条现跑，无一条抄编排者或抄实现腿）

1. **起手同发三读数**：见 §0，时刻 `10:35:27`、HEAD `94bd133e`、`cmd internal` 脏项 **0 行**。
2. **被验的五枚提交逐枚复认**（`git show --stat --format="%h %ad %s" --date=format:"%H:%M"`）：
   - `e4ed8304`（docs，`1 file changed 53 insertions`）＝证据件 ⑥⑦⑧ 起手件；
   - `697b4fae`（`10:10`）＝`cmd/wisp/panel_host_windows.go` **+456**／`cmd/wisp/panel_host_gate_test.go` **+126**／`cmd/wisp/panel_host_windows_test.go` **+233**／`go.mod` **+2**／`go.sum` **+6**，共 823 增 0 删；
   - `b1e3d94a`（`10:26`）＝同三枚测试/产码文件的 **112 增／63 删**（专用线程 bring-up 从产码移进 `_test.go`）；
   - `fafe0445`（`10:26`）＝**只碰** `internal/panel/composer_dispatch_test.go`（**97 增／32 删**）＝那枚反转钉单独一枚提交；
   - `f0ad3b51`（docs）＝证据件 69 增／2 删（表里删除列非 0 的是它自己那份文档的 2 行，不是台账追加件 ⇒ 与"追加式文件删除列须为 0"那条纪律不同族，记下来备查）。
3. **实现腿的写面范围**（`git log --oneline 7a41db9b..HEAD -- cmd internal`，`7a41db9b`＝实现腿自报的真实起手 HEAD）＝**恰三枚**（`697b4fae`／`b1e3d94a`／`fafe0445`），与它自述一致；`git rev-parse --short 7a41db9b` 复认存在。
4. **票面 AC 框归零自证**（`git diff --numstat 7a41db9b..HEAD -- .scratch/wisp/issues/33-panel-host-c27.md`）＝**该路径零行输出**＝实现腿对票面**一字未动**；现量票面 `grep -c '^- \[ \]'`＝**12**、`grep -c '^- \[x\]'`＝**0**。⇒ "AC 框一枚没勾"这一条我复认成立（不是它自述）。
5. **裁定与预检的出处我已读**：票面「编排者裁定 J1–J8」全表（`:20-31`）＋新判据 AC#11（`:35`）／AC#12（`:37`）原文；台账 `A489`（`:10120` 起，反转裁与三处会被打假的指认）／`A490`（`:10153` 前，HEAD `7297e521` 上"设计内的红"）／`A491`（`:10153`，票 248 预检，与本程只共享"33 独占 `cmd/wisp`"这一句）。**两处编排者自认的错（`A490` 的文件名写错／`A491` 的措辞缺口）我不重复立案。**
6. **写面形状起手现量**：`cmd/wisp/panel_host_windows.go`＝**411 行**（`697b4fae` 落 456 行、`b1e3d94a` 净 −45）；`cmd/wisp/panel_host_gate_test.go`＝**180 行**；`cmd/wisp/panel_host_windows_test.go`＝**273 行**；`internal/panel/composer_dispatch_test.go`＝**718 行**。⇒ "宿主产码 456 行"那枚自述已过期，现读 411。
7. **工具链前提**（起手 `go env GOFLAGS` 空、`GOMODCACHE`、`GOPROXY` 见 §A#10 同发读数）＋ 基线命令必须带 `PATH=$PWD/third_party/sherpa-onnx:$PWD/build:$PATH`，否则 `0xc0000135`＝用例根本没跑（这条不是我的假设：票面与台账多处记录同形坑，我起手第一发就带上，并在红名册里核对有没有 `--- FAIL` 行）。
8. **三枚冻结件起手未动自证**（`git log --oneline 7a41db9b..HEAD -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go`）＝**零行**＝实现腿没碰那三枚（它的 `composer_dispatch_test.go` 不属冻结名册，反转属编排者在 `A489` 已裁）。

（本节继续追加：基线红名册、反转钉正控/反向正控复跑、定向突变、SLO 两发、无端口尺射程、AC#11 仓外清检出、AC#12 两口径、`AGENTS.md` §1.2 两条顺手核。）

---

## B. 我可能写错的条目（对抗我自己）

1. **"go test 里起得来的窗口" ≠ "产品里点得开"**。本表最可能的错就是把宿主原语的真机读数（测试在专用锁线程上建的窗）直接读成"用户点设置看得见"。⇒ 判 AC#1/AC#2/AC#4 时每一格都要具名写清**调用者是谁**（`cmd/wisp` 产码里枚数几枚、非用例调用者逐枚点名），⛔ 两件事不许合并成一句。
2. **我会不会把 `internal/panel` 起手那 4 枚红算成这程造成的**。归因尺＝`git diff 7a41db9b..HEAD -- internal/panel`（只有一枚测试文件 97/32）＋红句源码点名；⛔ 我不读 `frontend/**`／`design/**` 内容（两层禁令），只用红句、`git status` 的 `D` 行与断言源码定归因。判语只写"是不是这程造成"，不为它们放宽任何断言，也不动三枚冻结件让它们变绿。
3. **定向突变的还原风险**。我只改产码一枚、读完红名**立刻** `git cat-file blob HEAD:<path> > <path>` 还原（⛔ 禁 `checkout`/`restore`/`stash`），还原后取 `md5sum` 与起手拷贝逐字节比对，再复跑同发改绿；突变体一枚都不进提交。若还原后 hash 不同 ⇒ 这一格我判"判不动"并具名上报，不留"应该没问题"。
4. **SLO 读数的自相关**：`-count=2` 的第二发可能吃到第一发的 WebView2 缓存/常驻子进程，把"冷"读成"不冷"或把"热"读得更热。⇒ 两发各带时刻＋HEAD，并具名写这族在 CI **无 `-tags winlive` 档**（`A489`／`Q-75`）⇒ 只有本机有效；⛔ 阈值一字节不动（`internal/observe/thresholds.go` 起手只读不改）。
5. **AC#12 的两种口径混用**：`find frontend/dist`（磁盘，含他人**未跟踪**构建物）与 `git ls-files frontend/dist`（入库）差得很远。判 AC#12 只用**入库/干净检出**口径；⛔ 我不造页面产物填 `frontend/dist`（那是界面侧那枚 agent 的活），也⛔ 不把 AC#11 的 `rc=0` 读成"有内容"。
6. **"无监听端口"那把尺可能恒真**：如果 L2 那半按**测试进程 PID** 过滤，而真开端口的是 webview2 的**子进程**，或它取 PID 的方式取错，那"本 pid 零条 LISTEN"是**必然成立**的空尺。⇒ 我必须读它的射程（它数的是哪个 PID、有没有把行数当命中），并做一发**真开端口的仓外载体/产码突变**证明它会响；光看"零条"绝不算凭据。
7. **我会不会把它的自述当成它跑出来的**：那句泵再入 panic（`edge.Chromium.Init` nil `webview`）我只有复认"证据件里有没有真栈原文＋可复现命令"的义务；若只有栈的形状没有复现尺 ⇒ 标〔仅自述〕并登记为下一程前置，⛔ 不替它追认也不替它推翻。
8. **载体必须在仓外**：正控/反向正控/清检出全部建在 `C:/Users/swq/AppData/Local/Temp/`（或 `t.TempDir()`），⛔ 零枚 `go test -overlay`，⛔ 绝不在仓内建 worktree／checkout。
9. **HEAD 会动**：`244-a3` 正在往 `.scratch/wisp/probes/244/a3/census.md` 写并 commit 文档件。⇒ 每条读数都带取数时刻＋当时的 `git rev-parse --short HEAD`；收尾必用**起手锚** `94bd133e` 跑 `git diff 94bd133e..HEAD -- cmd internal`（应空），而不是用"最新 HEAD"自证。
10. **反转钉的射程我可能读歪**：新语义是"产码标识符族 ≥1 命中 **且** 主模块 `go.mod` 有 webview 模块"。我要分清**谓词**（`hostChannelCapabilityHits` 那类）与**判定形**（`==0` 改成 `!=0`）分别在哪一行，并逐枚数正控的命中数；只看函数名改名就算"已反转"是空判。
11. **`AGENTS.md` §1.2 的两条顺手核不是我该扩大射程的口子**：读 `panel_host_windows.go` 时只核"明文 API 密钥"与"`risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策"两形；若发现别的越界形状，**具名登记**交编排者，⛔ 不自己判它的票、⛔ 不动票 248／票 228 的活。
12. **我会不会把"这格勾不了"写成"这格红"**：AC#12 与 AC#3 的过滤器半、AC#1 的 dispose 半是**归口未落定／依赖边界**，判语要写"证据形状成立否＋归口"，勾不勾归编排者。

---

## C. 判不动的地方（逐条甲／乙／不做＋现量）

1. **AC#1 的"session dispose  destroys ＋ WebView children exit ≤2s"那一半**：票面 `:65-66` 要求 dispose 触发。现量（起手 `grep -rn "PanelManager" cmd/wisp/main.go cmd/wisp/resident_windows.go cmd/wisp/run.go`）**0 命中**（编排者 10:33 同尺，我自己会再跑一发）⇒ 若真没有会话级 dispose 的生产触发者，那"会话结束拆窗"这半**在产品里不可达**，我只能判"显式 `Destroy` 机制＋子进程回落"这一半。**请裁**：甲＝按机制判这半、那半登记未落地；乙＝要求这程先接 dispose 跳（超出它自述射程）；不做＝不放宽、不宣称。
2. **J1 的第一枚用例（票面 `:24`：落地腿"第一枚用例就证泵期间球仍能出帧、消息仍被派发"）**：实现腿自述这一枚**未提交**（本机 panic，改落进 ②/⑧）。⇒ 本程能判的是"证没证／停手上报合不合规"（它自述证不出就停手、未造第二线程绕＝符合票面），但**"证不出来"这一句是它跑出来的还是推出来的**要靠复认（见 §C#3）。真裁"要不要授权动 `internal/ball` 交泵权"＝契约/别模块写面，**归编排者摆 owner**，我不选。
3. **那句 panic 的可复现性**：我要在它的表里找**真栈原文＋复现命令**；现读 `33-panel-host-c27-r1.md` ②#2 给的是栈的**形状串**（`ballWndProc→…→webview2.NewWithOptions→CreateWithOptions(webview.go:340)→Embed→Init(chromium.go:131)`）与时刻口径，**没有贴原始 panic 输出、没有复现命令行**。⇒ 判〔仅自述〕；下一程前置：要么复跑一发留下栈文，要么改成别的架构裁。⛔ 我不自己复现那发（要动 `internal/ball` 写面，超出我只读＋突变射程）。
4. **AC#3 的"多文件页面经 `AddWebResourceRequestedFilter` 服务"那一半**：它交回的是 `pkg/edge` 控制器 `PutBounds` 吃模块内私有 `w32.Rect`（`pkg/edge/ICoreWebView2Controller.go:59`）⇒ 主模块给不了尺寸 ⇒ 只能退 `SetHtml`。这枚我要独立复认（读依赖模块真实签名），并明确写**它是依赖边界问题还是产品形状问题**——这句决定编排者要不要摆 owner，我不替他选，但**必须给"如果成立，属哪一类"**。
5. **AC#4 的 manual 那一半**（票面 `:69`「automated ＋ manual」）：我没有真编辑器在场上，`GetForegroundWindow` 取到的是控制台前台。⇒ 机制可判（记录前窗＋回还调用），"手感"判不动，标〔仅本机机制层〕；⛔ 不许写成"焦点回还在真编辑器里成立"。
6. **AC#12 的勾**：产物归界面侧那枚 agent，本编队⛔ 不写 `frontend/**` ⇒ 判据能否区分"只匹配到占位"与"真有一包页面产物"我可以复尺，但**这一格在"谁把 dist 填上"落定前勾不了**（票面 `:37` 逐字），我把这件事写成"判据形状是否成立＋现量归口未落"，不当红、不当绿。
7. **`frontend/dist` 磁盘态那三枚未跟踪产物**（`git status` 里属他人地界）：我不动、不提交、不删、不读内容；它们存在这一事实只用来解释"本机能 build 出带页面的 exe"为什么**不是** AC#12 的凭据。⛔ 我不打开那些文件、⛔ 结论不引到 `frontend/**` 身上。
8. **winlive 档在 CI 零岗位**（`A489`／`Q-75`）：本表所有真机读数只在开机这台机器成立，每条带时刻＋HEAD；要不要单开 workflow＝契约级，已由编排者登记为默认"不做"，我不重开这一裁。

（本节会继续追加：基线红名册里判不动的、突变跑不出红的、以及任何票面没覆盖的形状。）

---

## D. 逐格判语（六格，起手全部"未判"）

| 格 | 判据（我自己定的尺） | 判语 |
|---|---|---|
| **AC#1** 生命周期（二次 show 复用同窗、进程树子进程数稳定；dispose 拆窗＋子进程 ≤2s 退出） | 未判 | **未判** |
| **AC#2** 冷 ≤1500ms／热 ≤200ms | 未判 | **未判** |
| **AC#3** 离线供给＋无监听端口（＋embed 过滤器那一半） | 未判 | **未判** |
| **AC#4** 焦点回还 | 未判 | **未判** |
| **AC#11** 清检出能建（＋反向证：embed 模式改成不存在路径 ⇒ 必失败） | 未判 | **未判** |
| **AC#12** 能建 ≠ 有页面（embed 条目数／关键入口，不问退出码） | 未判 | **未判** |

### J1–J8 复核（裁定落地形）

| 裁定 | 判语 |
|---|---|
| J1 宿主投 `ui-sta`（甲） | 未判 |
| J2 缺 runtime 降级＋`log.Fatalf` 不可达 | 未判（本腿射程外的 AC#5，只核"有没有新增可达 fatal"） |
| J3 三口径分开报／P11 未触发／`thresholds.go` 一字不动 | 未判 |
| J4 落点甲（宿主进 `cmd/wisp`，`internal/panel` 零平台分叉） | 未判 |
| J5 依赖只走 `go mod` 解析、⛔ 手抄哈希（ban #5） | 未判 |
| J6 embed 路径按活接缝／"清检出能建"升硬判据（AC#11） | 未判 |
| J7 无监听端口 L1 进仓／L2 本机证据＋具名登记 CI 无这一半 | 未判 |
| J8 33 先于 248、不合批、不同 `244` 同批改 | 未判 |

---

## E. 收尾三把尺（终态填）

| 尺 | 起手 | 终态 |
|---|---|---|
| `git status --porcelain -- cmd internal` | 0 行（`10:35:27`／`94bd133e`） | 未判（须**等于起手名册**） |
| `git diff 94bd133e..HEAD -- cmd internal` | （锚＝起手 HEAD） | 未判（须**为空**） |
| 三节自对抗非空＋占位符 `grep -c` 为 0 | — | 未判 |

