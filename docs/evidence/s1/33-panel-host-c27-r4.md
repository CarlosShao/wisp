# `33-r4` — 面板宿主仪器腿（⛔ 只改 `cmd/wisp/*_test.go`；六处"看起来在测、其实看不见"的尺装牙）

- 程：`33-r4`（**写腿·仪器**）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）
- 判据出处（逐格指回）：`docs/evidence/s1/33-panel-host-c27-v1.md` §AC#1／§AC#2／§AC#3／§AC#4／§AC#12 与 §A#24/#25/#26/#27/#30/#31
- ⛔ 本件**不勾票面任何一枚 `- [ ]`／`- [x]`**（勾框归编排者与下一枚验收腿）；⛔ 不为变绿放宽任何既有断言。
- ⛔ 产码那半边（`cmd/wisp/panel_host_windows.go`／`resident_*.go`／`internal/ball/**`／`internal/proc/**`）本程**一字未改**，归 `33-r2`。

---

## 起手锚点

| 尺 | 现量（同发取，`date` 与 `git log` 一条命令里跑） |
|---|---|
| `date "+%Y-%m-%d %H:%M:%S %z"` | `2026-10-01 11:30:24 +0800` |
| `git log -1 --format=%H` | `eed229e48cb9599c29f72fea6a1170739a5f2b41` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd internal` | **0 行**（起手名册＝空；终态必须等于这一枚名册） |

⚠ 起手锚的用途：收尾用 `git diff eed229e4..HEAD -- cmd internal` 复算⛔ 突变体零枚进提交；HEAD 期间会被别枚腿推进，所以本件**每条读数都带取数时刻＋当时 HEAD**。
⚠ 同刻在飞的只读腿（我不与它们抢写面，也不动它们在读的文件）：`33-a2`（读 `internal/ball`＋webview2 泵面）、`228-a3`（读 `internal/ball` 托盘面＋`internal/agent/approval`）。

---

## ① 六格逐格：改前读数／改后读数／反控读数

> 派单给的那几枚行号我**自己现取**过（HEAD `eed229e4`，起手态逐字）：`:115` ＝ `tcpTableOwnerPIDAll = 4 // TCP_TABLE_OWNER_PID_ALL`；`:116` ＝ `tcpStateListen = 10`；`:125` 只传 `windows.AF_INET`；`:145-146` ＝ 重试耗尽后 `t.Logf(...)` ＋ `return 0`；`:238` ＝ `prior := windows.GetForegroundWindow()`（取在 `:209` 的 `HotShow` 之后）；`:249-251` 那一跳整条只有 `t.Logf`；`:175`／`:237`／`:260` 只判 0／非 0。对照真值 `internal/proc/treemetrics_windows.go:57` ＝ `tcpTableOwnerPIDAll = 5`。⇒ 行号与派单一致，无过期指认。

### 格 1 — netstat 那把尺（AC#3 的 L2 半，改前恒真）

| 项 | 读数 |
|---|---|
| 改前 | 表类 `4`（真身 `TCP_TABLE_OWNER_PID_CONNECTIONS`，行 layout 不同）＋状态 `10`（真身 `MIB_TCP_STATE_ESTAB`）＋只读 `AF_INET`＋重试耗尽静默 `return 0` ⇒ 任何监听端口都数不到，断言不可能红。v1 §A#23 用真开端口的产码突变复认过：宿主确实持 `127.0.0.1:62971` 而用例整条 PASS |
| 改后 | 常量修正（`5`／`2`）；读表失败 ⇒ `(count, tableRows, err)` ⇒ 调用方 `t.Fatalf`（"量不到"与"没有监听"从此长得不一样）；射程从 `os.Getpid()` 扩到**这台件自己起的进程树**；两族行 layout 全部**实测**；`want`（API 回写的实际字节数）参与越界检查，读到截断表报错而不是少数几行 |
| 反控（必做那发） | 仓外副本（`git archive HEAD` → `%TEMP%/33r4-clean`）里把我新尺的两枚常量**改回旧值**（`5→4`、`2→10`）⇒ `--- FAIL: TestAC3ListeningSocketRulerSeesItsOwnListener (0.01s)`，红句逐字 `the AC#3 ruler did not see the loopback listener THIS test just opened on 127.0.0.1:62221: counted 0, baseline 0, want at least 1`（`11:50:38`）。同一枚正控在修正态＝**PASS**：`11:45:36` 那发 `IPv4 table rows=224, this pid owned 0 before / 1 while 127.0.0.1:64538 is listening`、`11:52:20` 那发 `0 -> 1`、Close 后回 `0` ⇒ **这格有牙，且牙是量过的** |
| 关掉那枚 socket ⇒ 绿 | 同上：Close 之后 `the AC#3 ruler ... stops after Close (0)`，且产品断言那一支（面板活着时）终跑逐字 `IPv4 rows=314 LISTEN-owned-by-tree=0`（`11:56:0x`）⇒ 不是"永远 0"，是"真没有" |
| 还原自证 | 副本测试文件用 `cp` 从仓内还原后 md5 逐字节相同：`0ef571f8ed176df830811ebe2ecb93b7`；仓内那两枚测试文件从未进过突变态（`git status --porcelain -- cmd internal` 见 §③ 收尾） |

**`AF_INET6` 那一族按派单要求的读数（"交读数、别自己扩范围"）**——仓外探针 `%TEMP%/33r4-clean/ipv6probe/main.go`（`11:50:2x`，同发真开 `127.0.0.1:62218` 与 `[::1]:62217` 两枚监听再读表）：

```
== family AF_INET6(23): entries=439 bytes_written=24596
   layout rowSize=48 state=40 pid=44 : state==LISTEN rows=3  owned-by-self=0   <-- 头文件那套＝本机数不到自己
   layout rowSize=52 state=44 pid=48 : state==LISTEN rows=1  owned-by-self=0
   layout rowSize=56 state=48 pid=52 : state==LISTEN rows=27 owned-by-self=1   <-- 本机实测能用那一行
   layout rowSize=48 state=44 pid=40 : state==LISTEN rows=0  owned-by-self=0
   layout rowSize=44 state=36 pid=40 : state==LISTEN rows=1  owned-by-self=0
== family AF_INET(2): entries=221 bytes_written=5316
   layout rowSize=24 state=0  pid=20 : state==LISTEN rows=42 owned-by-self=1
   port anchor hit at word 459 (value 2803 ...) : window ... [8]2 ... [13]888=SELF   <-- 888＝探针自己的 pid
```

⇒ 三条结论：① v1 §B#14 卡住的"自己开了 `[::1]` 却归不到自己"是**行大小**问题（本机 `AF_INET6` 该表 56 字节/行、state 在 48、pid 在 52），不是 Windows 归属口径问题；② 我把这套 layout 用进了尺 ⇒ IPv6 那枚读数是**活的**（正控在 IPv6 同样量通：`11:52:20` 逐字 `listening on [::1]:55872, table rows=1048, this pid owns 1 LISTEN rows (baseline 0)`）；③ **AC#3 的产品断言我只对 `AF_INET` 开**（终跑逐字 `IPv4 rows=314 LISTEN-owned-by-tree=0 ; IPv6 rows=44 LISTEN-owned-by-tree=0`），要不要把 IPv6 升成闸门＝编排者的射程决定（§⑥ 第 2 条），⛔ 我不自己扩。IPv6 那一支我在正控用例里**断言了尺本身能数到自己的 `[::1]` 监听**（那是仪器自检，不是产品闸门；两者在红句里分开措辞）。

### 格 2 — 焦点回还（AC#4）：**这一格现在会红，红因在产码，归 `33-r2`**

- 改前：`prior` 取在 `HotShow` 之后 ⇒ v1 那 11 发里 `prior` 逐发＝面板自己；`:249-251` 只有 `t.Logf`；全仓对 `prevFocus` 零枚断言（我起手自尺 `grep -rn "prevFocus" --include=*.go cmd internal`＝只命中产码 `panel_host_windows.go:115/269/303/311`，测试侧 **0 枚**）。
- 改后：新立 `TestAC4FocusReturnToPriorWindowGap33r2`（**用例名把归属腿写死**），取样时刻挪到 Show 之前（先 `Hide`、再读 `prior`），四枚断言：① 尺能区分窗口（`foregroundBefore != panelHwnd`，分不开 ⇒ `t.Fatalf`＝量不到、不算绿）；② 产码自己记下的 `mgr.prevFocus` 不许等于面板句柄；③ Show 之后前台必须＝面板（票面"only the panel takes focus when shown"）；④ Hide 之后前台不许还是刚被藏起来的那扇面板窗，且若 `prior` 是外来窗则必须回到它。
- 真机读数（同一条用例四发）：

| 时刻 | HEAD | 跑法 | before | prior | after Show | panel | prevFocus | after Hide | 结果 |
|---|---|---|---|---|---|---|---|---|---|
| `11:48:0x` | `0ecd725c` | 整包（lifecycle 之后） | `0x808c2` | `0x6bc0c56` | `0x6bc0c56` | `0x6bc0c56` | `0x6bc0c56` | `0x6bc0c56` | **FAIL**（②＋④） |
| `11:53:43` | `23627aba` | **单独跑** | `0x808c2` | `0x138105a0` | `0x138105a0` | `0x138105a0` | `0x138105a0` | `0x138105a0` | **FAIL**（②＋④） |
| `11:56:0x` | `2288265b` | 整包终跑 | `0x12084c` | `0x3a0eaa` | `0x3a0eaa` | `0x3a0eaa` | `0x3a0eaa` | `0x3a0eaa` | **FAIL**（②＋④） |
| `11:52:59` | 副本＋MUT-D | 单独跑＋改产码记值 | `0x808c2` | `0x808c2` | `0x808c2` | `0x12260638` | `0x1234567` | `0x808c2` | **FAIL**（③；②④转绿） |

- **红因具名（给 `33-r2` 的题面）**：面板在 `bringUp`（`webview2.NewWithOptions` 建窗）那一步就已经拿到前台，而第一次 `Hide()` 时 `prevFocus` 仍是 0（`panel_host_windows.go:311` 那句 `if prev != 0` 直接跳过）⇒ 没有任何一跳把焦点还给用户真正来处那扇窗；随后 `Show`（`:269`）把**当时的前台＝面板自己**记成 `prevFocus`，"回还"的目标从第一步就丢了。①那一枚在每发都成立（before 那枚是另一扇窗 `0x808c2`／`0x12084c`）⇒ **不是尺瞎，是那一跳没做**。
- ⛔ 没调成绿：没删断言、没加"量不到就算过"的分支、没为它放宽任何既有判据。这枚用例现在**逐发红**。
- ⚠ 顺序依赖必须说清：断言③（Show 后抢到前台）四发里红了一发（副本单独跑那发），根因＝Windows 前台锁——本进程不是前台进程时 `SetForegroundWindow` 会被拒。这不是"尺不稳"，是**产品路径在真机上本就不可靠**的形状；我原样保留断言（不降级为 `t.Logf`），`33-r2` 要连它一起解。
- 反控（MUT-D，只在仓外副本改产码 `:269`）：把记值换成外来句柄 `0x1234567` ⇒ 断言②**当场转绿**、④转绿、③转红 ⇒ 四条断言彼此独立、都可翻转＝不是恒真尺。还原：`git cat-file blob HEAD:cmd/wisp/panel_host_windows.go > 副本`，md5 `e75b9b2e36772ca4cae6043a6dbe7384` 与仓内逐字节相同（这枚 md5 与 v1 §E 收尾那枚同名＝同一份产码自 09-28 起未变，交叉复认）。

### 格 3 — HWND 身份（AC#1 的 "reuses window"）

- 改前：全文没比较过那枚 HWND（只判 0／非 0）。
- 改后：`hwndBeforeHide` vs `hwndAfterReShow` 逐名比较，红句写明"新 HWND＝第二扇浏览器，不是 re-show"。读数：终跑 `11:56` 逐字 `same HWND 0x1f20daa across hide->re-show=true`；另三发分别 `0x142e056e`／`0x109c0614`／`0xc8c09cc`（逐发不同＝不是硬编码来的绿）。
- dispose 那一半（票面 "session dispose destroys"）：新立 `TestAC1SessionDisposeHasNoProductionTriggerYet_AC1` ⇒ **具名 skip**，理由用尺取（AST 扫 `cmd/wisp` 非 test 文件）：终跑逐字 `AC#1 dispose scan: 0 production constructor(s) [] | 2 .Destroy() call site(s) [panel_host_windows.go:209 panel_host_windows.go:328]`。⚠ 那两枚 `.Destroy()` 是宿主内部对 `webview2.WebView` 的调用（不是会话级拆窗点），我**把这枚不精确原样写在注释里**；一旦有人接上生产构造点而没有任何 `Destroy` 调用点，这枚 skip 立刻转红。不可达的根＝常驻那两枚回调仍只打点（v1 §A#15），接它＝`33-r2`。
- ⛔ 这一格没留任何"全文只有 `t.Logf`"的形状：要么断言、要么具名 skip 并写出为什么不可达＋哪一程来填。

### 格 4 — process-tree 的分母（AC#1 的 child-count）

- 改前分母（尺＝全机按 exe 名数）：v1 §A#26 ＝ 基线 14／窗口期 20／销毁后回 14。我起手复跑取到同形：终跑窗口期 `machine-wide msedgewebview2=20`，而 `11:48` 那发销毁后是 **`machine-wide 14 -> 15`**＝**我们这棵树没动、全机那枚自己涨 1**（别人的窗）⇒ 旧尺会因别人起／退一扇窗而假红／假绿，这一枚是现抓的。
- 改后分母（尺＝从 `os.Getpid()` 递归走的**本机进程树**）：`our tree webview baseline 0 -> 窗口期 6 -> Destroy 后 0`；同一行还带 `direct-children=1`（**直接子级只有 1 枚**，其余是孙级 ⇒ 只数直接子级会读成"没起浏览器"＝又一枚假绿）与 `tree-pids=6`。三发逐字同形：`11:45:36`／`11:51:49`／`11:56:0x`（`AC#1 denominators ... our tree webview=6 direct-children=1 tree-pids=6 | machine-wide msedgewebview2=20`）。
- 断言形制：hide→re-show 比"我们树里的 msedgewebview2 枚数"（红句同时把全机那两枚当参照写出来）；Destroy 后 2 秒内回落到**本树基线**（`finalTree.TreeWebview > baselineTree.TreeWebview` ⇒ 红）。⛔ 全机枚数从此**不进任何判据**，只进日志。
- 反控：这格的"会响"由两枚形状证明——① `machine-wide 14 -> 15`（上面那发）证明旧分母会在我们没动的自己变＝旧判据今天可能因别人而红；② 换成树以后同一行两个口径并列，差值＝**别人的会话**（20－6＝14 枚，逐名对上 v1 的基线）。⛔ 没为了让它绿把树读成全机，也没反过来。
- ⛔ 没扩到 Job Object 那一维（v1 §A#25 指的现成口 `internal/proc/jobscope_windows.go:157 TreePIDs()`）：`internal/proc` 是本程禁写面，我用测试内 Toolhelp32 递归走树自给，读数具名如上。

### 格 5 — AC#12 那枚空尺

- 改前：`cmd/wisp/panel_host_gate_test.go:52-76` 除 `git ls-files` 报错的 `t.Fatalf` 外**零枚 `t.Errorf`**（v1 §A#31）。
- 现量起手态（HEAD `eed229e4`）：`git ls-files frontend/dist`＝1 枚 `.gitkeep`；`git status --porcelain -- frontend`＝**空**；`git status --porcelain --ignored -- frontend/dist`＝**2 行 `!!`**；`git check-ignore -v frontend/dist/index.html`＝`frontend/.gitignore:12:dist/*`。⇒ 工作树有产物、入库没有，且那些产物是 **ignored** 态（普通 `status` 看不见＝这条轴上一枚腿没写进表）。
- 改后：新用例 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`（旧名退役，见 §② 名册 diff）。判据**一律从 Go 侧能力问**：`panel.BuiltinAssets()` → `Built()`／`Resolve(EntryFile)`／`Check()`／`Manifest()`；⛔ 不问构建退出码、不问文案、⛔ 不开 `frontend/**` 一字节（provenance 轴只取 git 的**条目名**）。两条不变式＋一枚 provenance 轴。
- 两种口径的真实读数（同一枚尺，两棵树）：

| 口径 | 时刻／HEAD | 逐字读数 | 结果 |
|---|---|---|---|
| 工作树（ignored 产物在场） | `11:45:36`／`0d717917`；终跑 `11:56`／`2288265b` | `shape=page-bundle built=true entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil> manifest-entries=4 manifest-err=<nil> \| git-metadata=true tracked=1 tracked-beyond-anchor=0 tracked-has-entry=false ignored-or-untracked=2` | **PASS**（走"真有一包页面产物"那一支断言） |
| 清检出（`git archive HEAD` 副本，只有 `.gitkeep`） | `11:48:52` | `shape=anchor-only built=false entry-bytes=0 entry-err=panel: embedded assets are not built (run npm run build in frontend/) refs=0 check-err=同上 manifest-entries=0 manifest-err=同上 \| git-metadata=false` ＋ `AC#12 provenance axis not measurable in this tree` | **PASS**（走"占位 ⇒ 三门全 fail-closed"那一支断言） |

⇒ 这枚尺在**两种形状下都说话**（不是只会绿的那一枚），并把"此刻在哪个口径"写进同一行读数。
- 反控（MUT-B，只在仓外副本改产码）：`internal/panel/assets.go` 的 `newAssets` 改成"任何树都算 built"（＝把占位文件当一包页面）⇒ `--- FAIL: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 (0.05s)`，红句逐字 `the embed reports built=true yet Resolve("index.html") failed: panel: "index.html" is not in the embedded bundle: open index.html: file does not exist`（`11:52:59`，`rc=1`）。还原后 md5 `f9e3b3b078adeaa6c4d31f7cb2405d78` 与仓内逐字节相同；`git status --porcelain -- internal/panel/assets.go`＝**空**（仓内那枚文件从头到尾没进过突变态）。
- ⛔ 我没做的两件事（免得被读成做了）：① 没造任何页面产物、没写没读 `frontend/**`；② 没把 AC#12 的"归口未落"翻成红——那格按票面 `:37` 归编排者裁（§⑥ 第 5 条）。另：本尺把"半包"的完备性判定**委托给 `Assets.Check()`**，`Check()` 自己的牙在 `internal/panel/assets_test.go`（v1 §A#30 逐名复认过 `TestAnchorOnlyBundleIsNotBuiltAndFailsClosed`），本程禁写那枚包 ⇒ 这一层限制写进 §⑤ 第 8 条。

### 格 6 — latency 那格的两枚小账（AC#2）

- 改前：`:181` 把 `"HEAD 7a41db9b"` 写死在测量日志的字面量里（v1 §A#21：它 11 发都在别的 HEAD 上跑，日志仍写那枚）；票面要的 10-run P50/P95 在仓里不存在（我起手自尺 `grep -n "P50\|P95\|percentile" cmd/wisp/panel_host_windows_test.go cmd/wisp/panel_host_gate_test.go`＝**0 命中**）。
- 改后：① HEAD 运行时取数（`gitHeadShortForTest`；取不到 ⇒ 写 `HEAD-unknown`，⛔ 不编假锚——副本里实测就走了这一支：`head HEAD-unknown`）；② 同进程 `-count=N` 累积样本 ⇒ `TestPanelHostLatencyPercentilesAC2` 出 P50/P95（nearest-rank），按**同一对 D32 数**（1500／200）断言**尾部**，每发样本带 RFC3339 时刻逐行落日志；③ P11 那条 2000ms 只报最大冷值，不新设门槛。
- 真读数（本程自己那 10 发，`11:54:09`→`11:54:23`，HEAD `23627aba`）：cold `701.594／608.673／633.184／682.715／815.826／617.904／653.986／633.534／628.295／644.469` ⇒ **P50=633.534／P95=815.826**；hot `32.110／27.298／32.913／37.414／33.804／45.058／30.812／38.595／31.910／33.338` ⇒ **P50=32.913／P95=45.058**；P11 未触发（10 发最大 815.826 < 2000）。⛔ 阈值一字节未动（§③ 禁区自证），汇总只落 `docs/evidence/s1/`（本件 §⑦），⛔ `docs/SLO.md` 未碰。
- ⚠ 另两枚孤发读数（同一条用例、**非 10 发口径**，免得混算）：起跑基线那发（整包同跑、机器被我放并发压着）`cold=1366.298`／`hot=41.499`（`11:35:43`，距 1500 只差 134ms）；定向单跑 `cold=911.253／hot=46.601`（`11:51:49`）。⇒ "P95=815.826"是**同进程连续跑**口径，不是冷机首启口径（§⑤ 第 10 条）。
- 反控：这格的牙＝"分位真算、尾部真能过界"。我不能在不动阈值（禁区）的前提下造出越界样本，所以交**两枚替代证**：① n=10 时 nearest-rank 的 P95 取到的**正是那发最大尾值 815.826**（旧形态"每发只打一行"时这枚尾值只能靠人肉翻日志，现在它是判据本身）；② 断言写在 P95 而非 P50（红句带 `over %d runs`），一有过界尾值即红。⛔ 没为变绿动任何数，读数原样交。

---

## ② 起跑绿名册与终跑名册（逐名 diff）

### 起跑（`11:34:0x` 起钟→`11:37:54` 收钟，HEAD `eed229e4`）

尺：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp/ -count=1 -v`（日志 `.scratch/wisp/probes/33/r4/logs/baseline-start.txt`，名册 `.scratch/wisp/probes/33/r4/baseline-roster.txt`）

| 尺 | 现量 |
|---|---|
| 顶层 `--- PASS` 枚数 | **145** |
| 顶层 `--- FAIL`／`--- SKIP` | **0／0**（⇒ 起手名册＝全绿，且⛔ 不是 `0xc0000135` 那种"没跑"：日志里有 226 行 `=== RUN`、81 行子用例 `--- PASS`，包行 `ok github.com/CarlosShao/wisp/cmd/wisp 226.667s`） |
| 派单点名那枚时序 flake | `--- PASS: TestAC228ExitRequestDuringBootStillLeavesThroughD38E (4.72s)`＝**起跑不红**，无需 `-count=5` 隔离复认（终跑同样 PASS，4.06s） |
| 起跑日志顺手给的一枚**格 6 改前铁证** | `baseline-start.txt:528` 逐字：`cold bring-up measured on this box: 1366.298 ms (HEAD 7a41db9b, 2026-10-01T11:35:43+08:00)`——这一发跑在 HEAD `eed229e4` 之后、日志仍写 `7a41db9b`＝**硬编码假锚被我的起跑名册自己抓到了**（不是我引 v1 的话） |
| 起跑日志里另一枚格 2 铁证 | `baseline-start.txt:530` 逐字：`foreground after Show: 0x660dfe (panel hwnd 0x660dfe, foreground prior to test 0x660dfe)`＝三枚句柄同一枚＝旧尺"把焦点还给刚被藏起来的那扇窗"的形状，现抓 |

### 终跑（`11:55` 起钟→`11:57:57` 收钟，HEAD `2288265b`，名册 `.scratch/wisp/probes/33/r4/final-roster.txt`）

| 尺 | 现量 |
|---|---|
| 顶层枚数 | **149**＝`--- PASS` 147 ＋ `--- FAIL` 1 ＋ `--- SKIP` 1；包退码 `rc=1`（⚠ **这一发包是红的，红因与归属见下**） |
| 逐名 diff（起跑→终跑） | `+ PASS TestAC3ListeningSocketRulerSeesItsOwnListener`（格 1 正控）／`+ PASS TestPanelHostLatencyPercentilesAC2`（格 6 分位）／`+ SKIP TestAC1SessionDisposeHasNoProductionTriggerYet_AC1`（格 3 dispose 半，具名 skip）／`+ FAIL TestAC4FocusReturnToPriorWindowGap33r2`（格 2）／`- PASS TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12` → `+ PASS TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`（格 5 换名＋装牙） |
| **现在红的这一枚** | `--- FAIL: TestAC4FocusReturnToPriorWindowGap33r2 (1.04s)`。**红因＝产码缺那一跳**，⛔ 不是测试写坏：面板在 `bringUp` 建窗时就已拿到前台，第一次 `Hide()` 时 `prevFocus==0`（`panel_host_windows.go:311` 那句 `if prev != 0` 跳过），随后 `Show` 把"当时的前台＝面板自己"记成 `prevFocus`（`:269`）⇒ 没有一跳把焦点还给用户来处。逐发读数（5 发：`11:48`／`11:53:43`／`11:56`／`12:02:02`／副本 `11:52:59`）都在 §① 格 2 那张表里。**归属＝`33-r2`。** |
| 唯一一枚 SKIP | `TestAC1SessionDisposeHasNoProductionTriggerYet_AC1`＝尺自己量出来的不可达（生产构造点 0 枚），理由与"接上就转红"的机制在 §① 格 3。⛔ 不是"跳过以免红"，也没有任何 `t.Logf` 空尺混进来 |
| `internal/panel` 那四枚起手常红 | 我**没动它们**，也没为它们放宽任何东西：`12:02:45` 现跑逐名＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`＝**恰四枚，与 v1 §A#9 逐名相同**（含落在冻结件 `tokens_fourway_test.go:441` 那枚，红因＝工作树那 16 枚 ` D` 的 `design/**`）⇒ 不是本程造成，我也没读 `design/**` 内容 |
| ⛔ 我没有的权力 | 没删任何既有用例（唯一"消失"的名字是 AC#12 那枚换名，判据从 0 枚 `t.Errorf` 变成 5 枚）；没放宽任何既有断言（冷 1500／热 200／HWND／`IsCreated`／`h.count()!=1` 那几枚逐条仍在原位，只是**多了**断言）；没动 `TestAC228…` 一字 |

### 终态复跑（第二发整包，写面已含 `c8e715e3`，`12:08:5x`→`12:12:33`，HEAD `2dd28ace`）

尺：同一条 `PATH=… GOFLAGS= go test ./cmd/wisp/ -count=1 -v`（日志 `.scratch/wisp/probes/33/r4/logs/final-pkg-2.txt`，名册 `final-roster-2.txt`）。

| 尺 | 读数 |
|---|---|
| 顶层枚数 | **149**＝146 PASS／**2 FAIL**／1 SKIP |
| 与第一发终跑的逐名 diff | **只有 1 行翻转**：`PASS TestPanelHostRealWindowHopAndLifecycle` → `FAIL`；`TestAC4FocusReturnToPriorWindowGap33r2` 两发都 FAIL（＝本程故意的红）、SKIP 那枚两发都在 ⇒ 名册其余 146 枚逐名相同 |
| 这枚新红的逐字 | `after Destroy the WebView2 children this process started did not exit within 2s: our tree went baseline 0 -> now 5 (tree pids 5). The machine-wide count (14 -> 19) is reported only and is NOT the denominator`（`12:10:1x`） |
| **归因（不是我这枚改坏）** | 那一发**两个口径同时越界**：我们树 0→5、全机 14→19 ⇒ 旧尺（全机口径）在同一发也会红（`after 19 > baselineChildren 14`）＝**不是我换分母造出来的**；同一发冷值 977.297 ms 也高于同批 10 发的 608-816 ms ⇒ 那台机器当时在忙（这发跑在我连续跑 d22scan／vet／三枚副本突变之后） |
| `-count=5` 隔离复认（派单给 flake 定的姿势） | `12:13:04`→`12:13:1x`，`-run TestPanelHostRealWindowHopAndLifecycle -count=5` ⇒ **5 枚全 PASS**，逐发 `our tree webview baseline 0 -> now 0`；⚠ 其中第 2 发 `machine-wide 14 -> 15` 而我们树 0⇒**旧尺在这一发会假红**（`15 > 14`），这一枚是"全机分母会因别人起一扇窗而假红"的**活样本**（v1 §A#26 只是推断，本程抓到了实例） |
| 我的处置 | ⛔ 没放宽那 2 秒、没加重试、没把断言降级；登记为**负载敏感的已知抖动**：`WebView children exit ≤2s` 这一支在本机 6 发里红 1 发（整包并发那发），隔离 5 发全绿。⇒ 这条正好是票面 AC#1 那半句要抓的形状，留给 `33-r2`／验收腿判：是产码拆窗不够快，还是这台机在负载下就该报。⛔ 我不用"偶发"把它读成"不存在" |

---

## ③ 门禁四数（go build ./...／go vet ./...／d22scan／gofumpt，逐条真实读数）

时刻 `12:00:25`→`12:02:02`，HEAD `06acaca5`（终跑那一发在 `2288265b`）。

| 门 | 命令 | 真实读数（逐字） |
|---|---|---|
| build | `GOFLAGS= go build ./...` | **rc=0**（无输出） |
| vet | `GOFLAGS= go vet ./...` | **rc=0**（`PIPESTATUS[0]` 取法，无输出） |
| d22scan | `sh scripts/d22scan.sh` | **第一发 rc=1**：`scan_test.go:269: repo HEAD violates: cmd/wisp/panel_host_windows_test.go:571: [emoji] ban #8 glyph in scope cmd/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)` ⇒ 我把 `⛔`(U+26D4) 写进了 `t.Errorf` 的**字符串**里。⚠ **这一枚是我自己造的缺陷、被常驻门当场抓住**（我在 §⑤ 第 14 条事先写过这条雷，还是踩了）。改成 ASCII 后**rc=0**，逐字：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=224, bans #1-5 cmd/=33, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=476, ban #8 cmd/=78`（分母与 v1 §A#34 逐字相同＝我没把任何包的计数顶起来） |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/ internal/panel/` | **空列表，rc=0** |
| staticcheck | — | **〔未复认〕**：本机 `2025.1.1 (0.6.1)` 与 CI 钉的 `2026.2.1` 不同版（§⑥ 第 8 条），我不跑它拿假结论 |

### 禁区与写面自证（终态现量）

| 尺 | 读数 |
|---|---|
| `git status --porcelain -- cmd internal` | 起手 **0 行** → 本程写面＝**只两枚 `*_test.go`**（`panel_host_gate_test.go`／`panel_host_windows_test.go`），收尾提交后回 **0 行**；⛔ 零枚产码文件（编排者 `A497` 在 `11:46:28` 也独立量到同一形状并逐名核过"全是 `*_test.go`"） |
| `git diff --numstat eed229e4..HEAD -- cmd internal`（我的起手锚） | `cmd/wisp/panel_host_gate_test.go 250/25`、`cmd/wisp/panel_host_windows_test.go 516/66`＝**只有这两枚路径**（⛔ 别家腿在我锚后没碰过 `cmd`／`internal`）。删除列逐名解释见 §⑧ 第一枚 |
| `git diff --numstat eed229e4..HEAD -- internal/observe/thresholds.go docs/SLO.md docs/PLAN.md docs/specs docs/BUILD.md tools/d22scan/allowlist.txt` | **空**（一字节未动） |
| 三枚冻结件 | `git diff --numstat eed229e4..HEAD -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go`＝**空** |
| `go.mod`／`go.sum` | `git diff --numstat eed229e4..HEAD -- go.mod go.sum`＝**空**；⛔ 没跑 `go get`、⛔ 没跑 `go mod tidy`（它在 HEAD 上 exit 1，会造出不属本程的 diff） |
| `frontend/**`／`design/**` | `git diff --name-only eed229e4..HEAD -- frontend design`＝**空**；两层禁令自证：我只跑过 `git ls-files`／`git status --ignored`／`git check-ignore` 三枚**条目名**尺，⛔ 未打开一枚前端／设计文件 |
| 票面 | `git diff --numstat eed229e4..HEAD -- .scratch/wisp/issues/33-panel-host-c27.md`＝**1 增／0 删**＝只有我那一条 Progress log 追加行。⚠ 一枚**归因要说清**的数：现读 `^- [ ]`＝**13**（我锚点处是 12）＝**多的那一枚不是我翻的**，是编排者 `0ecd725c`（`11:47`，`A497`，新开票面 AC#14）；`^- [x]`＝1（AC#11，编排者 11:1x 翻的，早于我起手）。⛔ 我全程没碰任何复选框 |
| `docs/reports/pending-and-issues.md` | 未动（`git diff --numstat eed229e4..HEAD -- docs/reports/pending-and-issues.md`＝空） |
| git 纪律 | 只 commit、⛔ 未 push；每枚提交都带**显式 pathspec**；`git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean` **零次**；仓内**零删除**；突变全在仓外 `%TEMP%/33r4-clean` |

---

## ④ 我跑了哪些尺、每条真实读数

> 本节在**前 15 发工具调用内写满**（编排者硬要求），后续只追加、不删行。

1. **起手同发四读数**（`date`＋`git log -1 --format=%H`＋`git rev-parse --abbrev-ref HEAD`＋`git status --porcelain -- cmd internal`，一条命令里取）：`2026-10-01 11:30:24 +0800`／`eed229e48cb9599c29f72fea6a1170739a5f2b41`／`dev`／**0 行**。（时刻与锚点即 §起手锚点 那一发，同发。）
2. **验收表全文已读**：`docs/evidence/s1/33-panel-host-c27-v1.md` 219 行全文（§A#24/#25 常量错位与 AF_INET6 探针读数、§A#26 分母、§A#27 焦点、§A#30/#31 AC#12、§D 逐格判语、§E 收尾三把尺）。⛔ 我没有把它的任何读数抄成本程的交付读数——本件凡引用它的地方都标了"v1 §…"出处。
3. **被测仪器源码逐行现读**（不是抄表）：`cmd/wisp/panel_host_windows_test.go`（273 行全文）、`cmd/wisp/panel_host_gate_test.go`（180 行全文）、`cmd/wisp/panel_host_windows.go`（411 行全文，⛔ 只读）、`internal/proc/treemetrics_windows.go`（对照常量真值）。
4. **`cmd/wisp` 的 run 前提（派单点名那枚）**：`PATH` 形态照 `scripts/wisp-cli-tests.sh:100-106` 抄——脚本原话是 `pwd -W` 那种带盘符的形态会 **exit status 0xc0000135**（用例根本没跑），必须用 shell 自己的路径形态：`export PATH="$dll_dir:$PATH"`，`dll_dir="$root/third_party/sherpa-onnx"`。本程命令固定为
   `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp/ ...`（与 v1 §A#9 同形）。
   现量：`third_party/sherpa-onnx/` 里 `onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll` 三枚**在**（`ls` 现读，`11:3x`）。
5. **起跑基线（11:34:xx 起钟，后台跑，⛔ 不与其它腿并发）**：`GOFLAGS= go test ./cmd/wisp/ -count=1 -v`（带上面那枚 PATH）→ 日志 `.scratch/wisp/probes/33/r4/logs/baseline-start.txt`。名册与逐名绿／红待收（§②）。
6. **撞钉预检①（新符号名）**：`grep -n "func TestMain|func listenSocketsForPID|func countWebviewChildren|func treePIDs|panelHostCold|panelHostHot|func percentile|func gitHeadShort|func repoRootForTest" cmd/wisp/` ⇒ 现量 **只有两枚已存在的函数**（`panel_host_gate_test.go:79 repoRootForTest`、`panel_host_windows_test.go:89 countWebviewChildren`、`:110 listenSocketsForPID`）。⇒ `cmd/wisp` **没有 `TestMain`** ⇒ 本程若要用 TestMain 汇总是安全的（不会与既有 TestMain 撞名）。
7. **撞钉预检②（行为型钉，派单第 64 条那味）**：`scripts/portable-tests.sh:195` 的 cli scope 与 `scripts/wisp-cli-tests.sh:112` 只跑包、不按用例名取数 ⇒ 新增/重命名用例不会打断 CI 取数点；`internal/panel/composer_dispatch_test.go` 那枚反转钉扫的是**产码标识符**（v1 §A#12），`_test.go` 不入射程 ⇒ 我在测试文件里新增 AST 扫描不会打红它。⚠ 这两条是**读源码得出的**，不是跑出来的；跑出来的复认在 §③ 门禁那发（终态 `go test ./internal/panel/` 红名册必须仍等于起手的 4 枚）。
8. **AC#12 的三种口径现量**（见 §①格 5 那一发）：入库＝1 枚 `.gitkeep`；普通 `git status`＝**空**（因为 `frontend/.gitignore:12` 的 `dist/*` 把整批产物标成 ignored）；`--ignored` 口径＝2 行 `!!`。⇒ 这枚尺以后必须区分"tracked／ignored-dirty／embed 内容"三条轴，⛔ 只数 `git ls-files` 的那枚旧尺看不见第二条轴。
9. **禁区自证（起手态）**：`git status --porcelain -- cmd internal`＝**0 行**；`internal/observe/thresholds.go`／`docs/SLO.md`／`docs/PLAN.md`／`docs/specs/**`／`tools/d22scan/allowlist.txt`／三枚冻结件起手均未动（终态再取一发同尺）。⛔ 本程不跑 `go mod tidy`（v1 §A#29 现量它在 HEAD 上 `exit 1` 会造出不属本程的 diff）。
10. **`frontend/**`／`design/**` 两层禁令的自证**：本程只跑过 `git ls-files`／`git status --ignored`／`git check-ignore` 三枚**条目名**尺，⛔ 未打开过任何前端或设计文件，⛔ 结论不引到 `frontend/**` 身上（AC#12 那格的判据一律从 `panel.Assets` 的 Go 侧能力问）。

11. **起跑基线收讫**（`11:37:54`，HEAD `eed229e4` 那一发）：`ok github.com/CarlosShao/wisp/cmd/wisp 226.667s`，顶层 `--- PASS`＝**145**、`--- FAIL`＝**0**、`--- SKIP`＝**0** ⇒ 起跑名册逐名落 `.scratch/wisp/probes/33/r4/baseline-roster.txt`（§②）。这一发同时抓到两枚"改前铁证"（`baseline-start.txt:528` 的假锚 `HEAD 7a41db9b`、`:530` 的三枚句柄同一枚），逐字在 §②。
12. **10 发时延汇总（本程自己的交付读数）**：`11:54:09`→`11:54:23`，命令逐字
    `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp/ -count=10 -v -timeout 900s -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'`（HEAD `23627aba`，日志 `.scratch/wisp/probes/33/r4/logs/latency-count10.txt`，包行 `ok ... 11.530s`）⇒ appendix 在 §⑦。
13. **正控三枚真跑读数**（都是"种 X 必响／关掉必停"那一族）：① IPv4 `0 -> 1`（`127.0.0.1:64538`，`11:45:36`；`55871`，`11:52:20`；`59642`，`11:51:49`），Close 后回 `0`；② IPv6 `0 -> 1`（`[::1]:55872`，`11:52:20`，`table rows=1048`）＝**layout 56/48/52 是被量出来的**；③ 生命周期那发的产品断言态：`IPv4 rows=314 LISTEN-owned-by-tree=0`（终跑 `11:56`）＝面板真活着时确实零枚。
14. **仓外副本三枚反控**（载体全是 `%TEMP%/33r4-clean`，`git archive HEAD` 抽的，零枚 `go test -overlay`、仓内零突变）：MUT-C（旧常量塞回尺）⇒ `--- FAIL: TestAC3ListeningSocketRulerSeesItsOwnListener`，`11:50:38`；MUT-B（`newAssets` 任何树都算 built）⇒ `--- FAIL: TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12`，`11:52:59`，`rc=1`；MUT-D（`Show` 记值换外来句柄）⇒ 焦点用例三条断言独立翻转，`11:52:59`。逐枚还原后 md5：`assets.go f9e3b3b078adeaa6c4d31f7cb2405d78`、`panel_host_windows.go e75b9b2e36772ca4cae6043a6dbe7384`、副本测试文件 `0ef571f8ed176df830811ebe2ecb93b7` ＝与仓内逐字节相同；同发 `git status --porcelain -- cmd internal`＝只有我那两枚测试文件（`M`）。
15. **AC#12 两口径各一发**（同一枚尺、两棵树）：工作树 `shape=page-bundle built=true entry-bytes=1044 refs=2 manifest-entries=4`（`11:45:36`／终跑 `11:56` 同形）；清检出副本 `shape=anchor-only built=false`＋三门全拒 `panel: embedded assets are not built (run npm run build in frontend/)`＋`git-metadata=false`（`11:48:52`）。逐字在 §① 格 5。
16. **常驻门当场抓了我一枚**（`12:00:25` 第一发 `sh scripts/d22scan.sh` **rc=1**）：红句逐字 `repo HEAD violates: cmd/wisp/panel_host_windows_test.go:571: [emoji] ban #8 glyph in scope cmd/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)` ⇒ 我把 `⛔`(U+26D4) 写进了 `t.Errorf` 的字符串。改成 ASCII 后 rc=0（`12:01:5x`，分母逐字与 v1 §A#34 相同）。⚠ 这条我在 §⑤ 第 14 条**事先写明是雷**还是踩了，所以把"仪器射程 ≠ 规格射程"这一课记死：`AGENTS.md` 那句"注释豁免、字符串不豁免"对我自己同样成立。
17. **`internal/panel` 起手常红逐名复认**（`12:02:45`，⛔ 我没动它们）：恰四枚 `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，与 v1 §A#9 逐名相同。
18. **门禁四数**（§③ 那张表）：`go build ./...` rc=0、`go vet ./...` rc=0、`d22scan` rc=0（修掉第 16 条那枚之后）、`gofumpt -l cmd/ internal/panel/` 空。staticcheck〔未复认〕。
19. **`cmd/wisp` 单独包级终跑**（`11:55`→`11:57:57`，HEAD `2288265b`）：`rc=1`＝**147 PASS／1 FAIL／1 SKIP**（红与 skip 都是本程新装的尺，逐名在 §②）；⛔ 零枚既有用例由绿转红。
20. **写面闸门**：`git status --porcelain -- cmd internal` 起手 0 行 → 全程只两枚 `*_test.go` → 收尾提交后回 0 行（§⑧ 那一发现量）。

21. **终态复跑（第二发整包，HEAD `2dd28ace`）＋一枚负载敏感的红**：`146 PASS／2 FAIL／1 SKIP`，与第一发终跑逐名只差一行（`TestPanelHostRealWindowHopAndLifecycle` 由绿转红），红句与逐字归因在 §②；`-count=5` 隔离复认 5 枚全绿（其中一枚同时给出"全机 14→15、我们树 0→0"＝旧分母会假红的活样本）。⛔ 没放宽那 2 秒、没加重试、没降级成 `t.Logf`。

---

## ⑤ 我可能写错的条目（对抗我自己）

1. **修正常量不等于那把尺就有牙**。我把 class 改成 5、state 改成 2 之后，仍必须做**真开一枚 `127.0.0.1` 监听 socket** 的正控；如果正控不红（尺看不见自己进程刚开的端口），那说明我的行 layout 或归属口径还是错的，这一格判"仍无牙"并上报，⛔ 不许拿"常量改对了"交差。
2. **行 layout 是我自己算的，不是抄来的**：`MIB_TCPROW_OWNER_PID`＝6 枚 uint32＝24 字节、state 在偏移 0、owningPid 在偏移 20；`MIB_TCP6ROW_OWNER_PID`＝48 字节（localAddr16＋localPort4＋remoteAddr16＋remotePort4＋state40＋pid44）。v1 §B#14 在 `AF_INET6` 上取到过"表里有 LISTEN 行、own-pid 却对上 0 枚"，所以**IPv6 那一半我很可能又取错**；我只在本程自己的正控真跑通之后才主张任何 IPv6 结论，跑不通就具名写"这一维我没量到"。
3. **重试耗尽静默 `return 0` 的修法可能造出新假红**。我把它改成 `(count, error)` 后，调用方 `t.Fatalf` 会在**机器负载高**时红——那是"量不到"不是"违规"。这一枚红要有自证（错误句必须写清是工具失败），⛔ 不许反过来把它又降级成 `t.Logf`。
4. **焦点那一跳可能被我修成"两个都恒真"**。产品 `Show` 是在自己内部记 `prevFocus`（`panel_host_windows.go:269`，取在 `ShowWindow` 之前），所以只要测试的取样时刻也在 Show 之前，`afterHide == beforeShow` 就**几乎必然成立**——这正是"把断言写在产码已经做的那一步之后"的形状。⇒ 我必须补两枚才谈得上有牙：断言**面板 Show 后确实拿到前台**（否则"回还"这个问题根本不成立），以及**断言回还目标不是刚被藏起来的那扇面板窗**（today 的病形）。这两枚哪一枚在真机上是恒绿，我要具名说。
5. **焦点断言可能受 Windows 前台锁影响而 flake**。`SetForegroundWindow` 对非前台进程可能被拒（v1 §A#27 那句注释就是这意思）。⇒ 我不为它加 sleep 重试把它拧绿；如果它 flake，我把每一发的三个句柄逐名落表，交编排者判是不是产码那一跳真不稳。
6. **HWND 身份断言可能被 recreate 路径合法打断**。`Destroy` 之后 `Show` 会造**新**窗口（`:319-330` 明写），所以"同一枚 HWND"只在 **hide→re-show** 这一对上成立。⇒ 断言只放在 hide→HotShow 那一跳；dispose 之后另说，⛔ 不许把它写成"整个生命周期 HWND 不变"。
7. **process-tree 的分母我可能数窄了**。WebView2 的子进程可能不是测试进程的直接子级（loader/服务重挂父）；按 `ParentPID == 本 pid` 直接走会数到 0，那把"我们的树"读成"没人起窗"＝新的假绿。⇒ 我从 `os.Getpid()` 往下**递归**收集整棵子树，并具名报"直接子级枚数／全子树枚数／其中 msedgewebview2 枚数"三个数；如果全子树也是 0，这一格我判"改后仍量不到"，不硬说它有牙。
8. **AC#12 那枚尺的两个危险形状**：① 把"工作树有产物"读成"AC#12 已满足"（那是 33-v1 已判过的假绿，入库口径仍只有 `.gitkeep`）；② 把"入库只有占位"判成一枚**红**，替界面侧那枚 agent 背锅——归口未落是票面 `:37` 的原话，⛔ 本程不许拿红替编排者做那个决定。⇒ 我的判据只钉两条不变式（Built() 与 fail-closed 必须自洽；Built() 为真时入口引用必须全能解析），归口那一格写成**具名 skip 或具名读数**，不写成绿。
9. **反控可能需要临时改产码**：本程唯一允许的形态是**仓外副本**里改；若我不得不在树内做，改完立刻 `git cat-file blob HEAD:<path> > <path>` 并取还原前后 md5 相同。⛔ 绝不用 `checkout`／`restore`／`stash`，⛔ 突变体零枚进提交（编排者用 `git diff eed229e4..HEAD -- cmd internal` 复算）。
10. **10-run P50/P95 是自相关样本**：同进程 `-count=10` 的第 2 到第 10 发会吃到第一发留下的 WebView2 缓存/常驻子进程，把冷读低。⇒ 分位数落表时逐发带时刻与冷/热原文，并具名写"同进程连续 10 发"这一口径，⛔ 不写成"冷启真实分布"；P11 那道 2s 重评线我按同口径说明未触发。
11. **⛔ 不许把 `internal/panel` 起手那 4 枚红算成本程的**（`tokens_fourway_test.go:441` 落在**冻结件**上，红因＝工作树那 16 枚 ` D` 的 `design/**`）。⇒ 终跑名册按**逐名**比，不比退码；本程只要求 `cmd/wisp` 侧不新增红，新增的红只允许是"我给尺装牙装出来的、且红因具名"那几枚。
12. **时序 flake `TestAC228ExitRequestDuringBootStillLeavesThroughD38E` 可能在本程起跑/终跑里换态**。派单点名它起手即在。⇒ 偶发红用 `-count=5` 隔离复认并具名登记，⛔ 不改它、不放宽它、也不把它算成本程新增红（除非 `-count=5` 里它枚枚都红，那我如实上报）。
13. **硬编码 HEAD 的修法本身可能引入一枚新尺错**：`git rev-parse --short HEAD` 在 `go test` 的 cwd（`cmd/wisp`）里能取到仓根，但在**仓外副本**里 archive 出来的树**没有 `.git`** ⇒ 那枚 helper 若 `t.Fatalf` 会在副本里炸。⇒ 我把它写成"取不到就具名报 `HEAD-unknown`"，⛔ 不取用假值。
14. **d22scan 的 emoji 扫描（ban #8）可能打到我的新测试文件**：`cmd/` 在 ban #8 的走树里（v1 §A#34 现量 `ban #8 cmd/=78` 是分母），注释豁免、**字符串不豁免**。⇒ 新写的字符串里我不用 `✓`(U+2713)、`≤`(U+2264)、任何 U+1F000–U+1FAFF；终态 `sh scripts/d22scan.sh` 必须仍 exit 0。
15. **票面进度行**：往 `.scratch/wisp/issues/33-panel-host-c27.md` 的 Progress log 追加时**只用 Edit/Write**，⛔ 不带反引号的 heredoc（记忆里那族"静默吃字"事故），⛔ 不改上面任何一格、⛔ 不碰 AC 复选框。

16. **我确实踩了本条第 14 项那枚雷**（不是"差点"，是踩了）：`t.Errorf` 字符串里的 `⛔`(U+26D4，落在仪器实际扫的 `U+2600–U+27BF` 带里) 被常驻门 `sh scripts/d22scan.sh` 当场判红（`12:00:25`，红句在 §④ 第 16 项）。我事先在 §⑤ 写过"字符串不豁免"，还是把它写进了字面量。⇒ 这一条**不删、不改写**，留在这里证明"薄索引里那条仪器射程"对写手自己同样成立。
17. **我自己写正控时犯过一次基线取错的错**：IPv6 那一支我第一版把 `base6` 取在**已经开了 `[::1]` 监听之后** ⇒ `11:51:49` 那发当场误报"the ruler's IPv6 layout is blind"（`counted 1, baseline 1`）。把基线取回"种 X 之前"之后，同一把尺读到 `0 -> 1`（`11:52:20`）。⇒ 教训（值得进派单样板）：**正控的基线必须取在植入之前**，否则新装的牙会先咬自己一口，而"仪器自证"那一半就成了假案。这一条被我自己造的正控抓出来，说明那枚正控的写法是对的。
18. **焦点那枚断言读的是产码的非导出字段**（`mgr.prevFocus`，同包合法）。代价我说清楚：如果 `33-r2` 把那枚字段改名或搬走，这枚用例是**编译失败**而不是转红——我认这个代价（宁可编译失败，也不要一枚恒真的观察口），但这让本用例与产码内部形状耦合，验收腿该知道。
19. **dispose 那枚 AST 尺只认"名字叫 `Destroy` 的调用"**，所以它把宿主内部对 `webview2.WebView` 的 `w.Destroy()` 也算成调用点（现量正是那两枚：`panel_host_windows.go:209`／`:328`）。⚠ 今天不影响判定（构造点＝0 枚 ⇒ 走 skip 那一支），但等有人接上构造点、又只在库里销毁 webview 而不接会话拆窗，这枚尺会**误绿**。不精确我写进了用例注释与 §① 格 3；要精确到"receiver 是 `*PanelManager`"得引 `go/types`（测试变重），留裁给编排者。
20. **⚠ 本程让 `cmd/wisp` 整包变红，从而让 CI 的 windows cli 腿变红**：`TestAC4FocusReturnToPriorWindowGap33r2` 逐发红（§②），而 `scripts/wisp-cli-tests.sh` → `scripts/portable-tests.sh --scope=cli` 跑的就是 `./cmd/wisp/`（v1 §A#32 的接线，我现读同形）。⇒ 按派单"判据是会响、不是包全绿"我**没有**把它调成绿；但这一句属**门禁形状的账**，必须编排者本人知情：要么 `33-r2` 补那一跳（正解），要么临时把该用例名册写进"已知红"（我没那个权力，票面复选框与门禁文件我一枚未动）。
21. **P50/P95 的口径混用风险**：我这 10 发是**同进程连续**（第 2 发起 WebView2 环境/用户数据目录已热），所以 P95=815.826 ms 不是"冷机首启"分布；起跑那发整包同跑时冷读到 `1366.298 ms`（负载态、距 1500 只差 134 ms）——两枚我都逐名落了，⛔ 不合并成一个"冷启 P95"。票面 AC#2 那句"10-run P50/P95 into SLO appendix"我按"读数落 `docs/evidence/s1/`"交付（§⑦），⛔ 未写 `docs/SLO.md`。
22. **格 4 那枚 `direct-children=1` 说明我一开始的走树方式就可能是错的**：只按"直接子级"数会读成 1（甚至 0），所以递归是必须的；如果哪台机上 WebView2 把子进程挂到别的父进程下（服务重挂父），这棵树会读 0 枚＝又一枚假绿。⇒ 我在读数里同时打 `tree-pids` 与 `machine-wide` 两枚，任何一发出现"窗口期树内 0／全机涨"就能当场看出来。这条我判**本机未命中**（4 发逐发 6 枚），不是"已解决"。
23. **`nearestRankPercentile` 是我自己实现的**（仓里没有分位库件）。n=1 时 P50＝P95＝那一枚样本，所以单发跑不会"造出分位"——它只在 `-count>=2` 时才有意义，这一点我在用例注释里写了。如果编排者认为分位应当用别种定义（线性插值那一类），红/绿判定在 n=10、尾部只差一枚的情况下会不一样，请指名，我不自己换定义。
24. **⛔ 我没有做的事**：没动产码一字（`git diff` 可核）、没翻票面任何一枚框、没跑 `go mod tidy`、没读没写 `frontend/**`／`design/**`、没碰三枚冻结件与阈值／golden／allowlist、没把任何一枚旧断言放宽、没在仓内建 worktree、没删任何文件（副本里也是**只建不删**：MUT-* 三发都用覆盖还原，没 `rm`）。

25. **那枚负载敏感的红，我的归因只走了一半**：我说"不是我把分母换坏"的证据是**同一发里两个口径同时越界**（树 0→5、全机 14→19）——这是一发样本，不是分布。反过来我也**不能**把 `-count=5` 全绿读成"没问题"：那 5 发是隔离跑的，本来就不该期待它们红。⇒ 正确的口径是：**这一支（`WebView children exit ≤2s`）在本机 6 发里红 1 发，红那发是整包并发**；要定它是产码拆窗慢还是机器忙，需要"Destroy 调用→最后一个子进程退出"的**延迟分布**，那是 `33-r2`／验收腿该量的东西，⛔ 我不拿自己这一发行程去替它下结论，也⛔ 不为了少一枚红把 2 秒改宽。
26. **⚠ 本程可能让"包全绿"这一维从此不再成立**：除了 §⑤ 第 20 项那枚故意的红，AC#1 的退出那一支今天也开始偶发红。如果编排者要的是"`cmd/wisp` 名册必须等于起手名册"这一档闸门，本程**没做到**（终态名册比起手多了 2 枚 FAIL＋1 枚 SKIP，逐名可核），这是按派单"判据是会响，不是包全绿"的取舍，我把它写死在这里，不让它被读成疏忽。

---

## ⑥ 判不动的地方（逐条甲／乙／不做＋现量）

1. **格 2 的"回还那一跳在产码里到底做没做"——本程已把它量出来：没做**。起手我推不出（`panel_host_windows.go:269` 确实记值、`:311-313` 确实 `SetForegroundWindow(prev)`＋`SetFocus(prev)`，机制在场，所以真断言"不一定红"）。真机五发读数（§① 格 2 那张表）给的是**红**：面板在 `bringUp` 建窗时就已拿到前台 ⇒ 第一次 `Hide()` 时 `prev==0`、`if prev != 0` 那一支直接跳过，随后 `Show` 把"当时的前台＝面板自己"记成 `prevFocus` ⇒ 回还目标从第一步就丢了。⇒ 甲（我做的）＝逐名断言＋红因具名交编排者；乙＝先改产码——⛔ 本程写面禁；不做＝⛔ 不许把断言降级成 `t.Logf` 换绿。**这一格现在不再"判不动"，改成"已量、待 `33-r2` 补那一跳"**；留在本节只为说明我起手时它确实判不动（免得下一枚腿以为这是纸上推出来的）。
2. **`AF_INET6` 那一族的"能不能量"已定案，"要不要进门禁"仍归编排者**。本程现量：那台表在本机的行**不是**头文件那套 48 字节/行——探针实测 `56` 字节/行、state 在 48、pid 在 52，同一枚 `[::1]` 监听在这套 layout 下归到本 pid（`0 -> 1`，`11:52:20`），而 48/40/44 读 0 枚＝v1 §B#14 卡住的那一枚假象的根因（§① 格 1 那张矩阵）。⇒ 现在**尺是活的**，但 AC#3 的产品断言我只对 `AF_INET` 开（终跑逐字 `IPv4 rows=314 LISTEN-owned-by-tree=0 ; IPv6 rows=44 LISTEN-owned-by-tree=0`）。甲＝保持现状（IPv6 只进日志＋仪器自检断言）；乙＝编排者点头我就把 `v6 != 0` 也升成产品断言（一行的事，且 layout 已量过）；不做＝⛔ 不拿"今天 webview2 没开 IPv6 端口"（v1 §A#25 的读数）当"不需要覆盖"的理由，也⛔ 不自己扩。
3. **AC#13（冷启探测页盖掉真页面）**。现量：`bringUp` 的 `:221 serveEntry()` 之后 `:227 firstRoundTripLocked(...)`，而后者在 `:374-375` 又 `SetHtml` 一枚自造探针页。⇒ 这是**产码时序**，票面 `:39` 那格的判据（重排次序／过滤器＋一枚会响的"最终文档含真入口"断言）落点在 `panel_host_windows.go`，本程⛔ 不许改。甲＝本程只把"显示的是探针页还是入口页"这条**量出来**（读 embed 入口字节的特征 vs 探针页特征，从 Go 侧字符串比，不读前端文件），交 `33-r2`；乙＝要我现在钉成断言 ⇒ 会造出一枚"本程自己判自己产码缺跳"的红，归属混淆；不做＝⛔ 不留空尺。
4. **`PanelManager` 生产调用者枚数**（dispose 半、也是 owner 那句"我要能点"的那半）。现量：v1 §A#14 全仓 `grep -rn "PanelManager"`＝定义文件＋测试文件两枚 ⇒ **0 枚非 test 调用者**。这一格归 `33-r2`（接常驻）与票 248，本程⛔ 动不了；我能做的是把"dispose 在产品里不可达"从散文变成**一枚会随接线自动转红的具名 skip**。
5. **AC#12 的勾与归口**（谁把 `frontend/dist` 填上）。现量：入库 1 枚 `.gitkeep`、工作树 2 行 `!!`（§④#8）；票面 `:37` 逐字"这一格在'谁把 dist 填上'落定前勾不了"。⇒ 本程**只装牙、不判归口**，⛔ 不造页面产物、⛔ 不写 `frontend/**`。
6. **`winlive` 档在 CI 零岗位**。现量：本程所有真机读数（窗口、焦点、netstat、时延分位）只在开机这台机器成立；`scripts/portable-tests.sh` 的 cli scope 里没有 `-tags winlive` 这一档（v1 §A#19/#32 同判，编排者 `A489`/`Q-75` 第四次具名）。⇒ 每条读数带时刻＋HEAD，单开 workflow＝契约级，我不重开这一裁。
7. **`go mod tidy` 那味**（webview2 记成 `// indirect`）。现量：v1 §A#29 在仓外副本 `exit 1`，仓内 CI 无 tidy 门。⇒ 属 `go.mod` 禁区，本程⛔ 不跑不修；只在名册里说明" tidy 若被人跑一次会造出不属本程的 diff"。
8. **staticcheck 本机版与 CI 钉版不同**（本机 `2025.1.1`，CI `2026.2.1`，`ci.yml:181-200` 自陈旧版在 go1.27 export data 上导入期即崩）。⇒ 门禁四数里这一门**不复认**，交 CI；本件标〔未复认〕，⛔ 不跑它拿假结论。
9. **起跑名册的"没跑"与"绿"之分**。现量：`cmd/wisp` 缺 sherpa DLL 时是 `exit status 0xc0000135`＋`0.0xxs`＋**零枚 `--- FAIL`**（`scripts/wisp-cli-tests.sh:100-106` 逐字记录过这个坑）。⇒ §② 的名册只在日志里**真有 `--- PASS`／`--- FAIL` 行**时才算量到；若起跑那发是 `0xc0000135` 形，我具名写"这一维我没量到"，⛔ 不写"包全绿"。
10. **反控的载体归口**：格 1／2／5 的反控都可能需要"让产码或依赖变坏"。现量：派单允许只在仓外副本做，仓内做则须 md5 还原＋双尺复算。⇒ 我优先仓外（`git archive HEAD` 抽副本，与 v1 同一手）；副本没有 `.git`，凡我新写的 git 尺在副本里都走 §⑤#13 那枚 `HEAD-unknown` 分支，⛔ 不在仓内建 worktree／checkout。
11. **`AC#13`（冷启探测页盖掉真页面）本程量到了但没有判据可落**。现量（Go 侧能力问，⛔ 不读前端文件）：工作树口径下 `panel.Assets.Resolve("index.html")` 出来的入口字节是 **1044 bytes**（§① 格 5），而 `panel_host_windows.go:374-375` 那一发 `SetHtml` 塞进窗口的是它自己拼的探针文档（长度远小于 1044、内容不含入口特征）⇒ "最终显示的不是面板页"这一形在本机是可分辨的。⛔ 但"钉住最终文档"要动 `bringUp` 那两发的次序或上资源请求过滤器＝产码，票面 `:39` 已把它写成 `33-r2` 的格（判据①②③逐字在票面）。甲＝我把可分辨性量出来交过去（本条）；乙＝要我在测试里断言"最终文档含入口真内容"⇒ 今天的产码下必红，而那枚红的题面与 AC#13 重叠，会把 `33-r2` 的归属洗成"33-r4 自己造的红"；不做＝⛔ 不留一枚只有 `t.Logf` 的空尺（这一格因此**不做成用例**，只做成本条具名读数）。

---

## ⑦ AC#2 的 10-run appendix（票面要的那枚交付形状；⛔ 未写 `docs/SLO.md`）

尺（同进程连续 10 发，一条命令）：
```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= \
  go test ./cmd/wisp/ -count=10 -v -timeout 900s \
  -run 'TestPanelHostRealWindowHopAndLifecycle|TestPanelHostLatencyPercentilesAC2'
```
起钟 `11:54:09`／收钟 `11:54:23`，HEAD `23627aba`，包行 `ok github.com/CarlosShao/wisp/cmd/wisp 11.530s`，原始日志 `.scratch/wisp/probes/33/r4/logs/latency-count10.txt`。

| # | 时刻（+08） | cold ms | hot ms |
|---|---|---|---|
| 1 | `11:54:12` | 701.594 | 32.110 |
| 2 | `11:54:13` | 608.673 | 27.298 |
| 3 | `11:54:14` | 633.184 | 32.913 |
| 4 | `11:54:15` | 682.715 | 37.414 |
| 5 | `11:54:17` | **815.826** | 33.804 |
| 6 | `11:54:18` | 617.904 | **45.058** |
| 7 | `11:54:19` | 653.986 | 30.812 |
| 8 | `11:54:20` | 633.534 | 38.595 |
| 9 | `11:54:21` | 628.295 | 31.910 |
| 10 | `11:54:22` | 644.469 | 33.338 |

- 仪器自己出的汇总行（逐字，n=10 那一发）：`AC#2 percentiles over 10 runs (nearest-rank): cold P50=633.534 P95=815.826 (budget 1500) | hot P50=32.913 P95=45.058 (budget 200)` ＋ `AC#2 P11 line (cold > 2000 ms) - max cold observed over these runs: 815.826 ms`。
- 判定：**P95 冷 815.826 < 1500** ⇒ 冷预算未越界；**P95 热 45.058 < 200** ⇒ 热预算未越界；**P11（冷 > 2000 ms 才重评 L2 卡）未触发**（本批最大 815.826）。
- 口径具名（⛔ 不与上面合并）：① 起跑整包那一发 `cold=1366.298 / hot=41.499`（`11:35:43`，HEAD `eed229e4` 之后那一发，整包并发压着机器，距 1500 只差 134 ms）；② 定向单跑 `cold=911.253 / hot=46.601`（`11:51:49`，HEAD `23627aba`）；③ 终跑包内那一发 `cold=868.380 / hot=48.326`（`11:56:00`，HEAD `2288265b`，逐字见 `final-pkg.txt`）。⇒ 三枚都是**孤发**，不进 P50/P95；"10-run"那一批是**同进程连续**（第 2 发起 WebView2 用户数据目录与浏览器进程已热），所以它量的是"这台机上一次会话内的热重复"，⛔ 不是冷机首启分布——票面 AC#2 没写口径，我把口径钉在这里，`33-v2` 别拿它当冷启口径用。
- ⛔ 阈值与文档零改动：`internal/observe/thresholds.go`／`docs/SLO.md` 一字节未动（§③ 那枚 `git diff --numstat` 空读数）。1500／200 这两个数是从既有断言**沿用**的（同一对 D32 数），不是本程新设的门槛。

---

## ⑧ 收尾三把尺（逐字读数）

| 尺 | 读数 |
|---|---|
| `git status --porcelain -- cmd internal` | 起手（`11:30:24`，HEAD `eed229e4`）＝**0 行**；终态＝**0 行**＝等于起手名册 ✓（本程写过的两枚测试文件已全部提交，仓内零残留、零突变体） |
| `git diff --numstat eed229e4..HEAD -- cmd internal` | `cmd/wisp/panel_host_gate_test.go 250/25`＋`cmd/wisp/panel_host_windows_test.go 516/66`＝**只有这两枚路径**，⛔ 零枚产码文件。删除列逐名解释（一共 91 行，全是"被替换掉的旧仪器"，没有一行是别人的活）：`panel_host_gate_test.go` 那 25 行＝① 旧 `TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12` 的注释与函数体（0 枚 `t.Errorf` 的那枚空尺，被 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 取代）② 旧 import 块（新增 `go/ast`／`strconv`／`internal/panel` 三行，重写整块故计入删除）；`panel_host_windows_test.go` 那 66 行＝① 旧 `countWebviewChildren`（全机按名数分母）② 旧 `listenSocketsForPID`（`class=4`／`state=10`／静默 `return 0` 那枚恒真尺）③ 旧 AC#4 那一跳（只有 `t.Logf` 的 4 行）＋它头上那 8 行注释，判据整体迁进新用例 `TestAC4FocusReturnToPriorWindowGap33r2`（取样时刻修正＋4 枚真断言）。⇒ **零枚既有断言被删掉不补**，迁移逐名可对 |
| `wc -l -c docs/evidence/s1/33-panel-host-c27-r4.md` | 交付判据＝**307 行／60,593 字节**（取数时刻 `12:14:34`，本行已入列；若与本行自身字节数有 ±一位数误差，以编排者复尺为准——本程此后不再改本件） |

### 收尾其它读数（同批发）

- 本程提交清单（逐枚 `git log --format="%h %ad %s"` 现量）：`0a17c9fd`（`11:37` 证据件骨架，④⑤⑥ 起手写满）→ `2288265b`（`11:55` 六格装牙，只两枚 `*_test.go`，766 增／91 删）→ `c8e715e3`（`12:08` 把那枚 U+26D4 请出断言字符串，1 增／1 删）→ `2dd28ace`（`12:09` 证据件终态＋票面两条 Progress log）→ 本枚（复跑读数＋§⑤ 第 25/26 条＋这行清单）。⛔ 全程只 commit、**未 push**。
- 每一枚提交的 pathspec 都是显式单文件／双文件；`git add -A`／`.` 零次；`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean` 零次；仓内零删除（副本里也只用覆盖还原，没有 `rm`）。
- 票面：`git diff --numstat eed229e4..HEAD -- .scratch/wisp/issues/33-panel-host-c27.md`＝**1 增／0 删**（我那一条 Progress log 追加行），复选框**一枚未碰**；表里 §③ 已把"多出来的那枚 `- [ ]`＝AC#14"归给编排者自己的 `0ecd725c`。
- 台件目录：`.scratch/wisp/probes/33/r4/`（`baseline-roster.txt`／`final-roster.txt`／`logs/baseline-start.txt`／`logs/targeted-1..3.txt`／`logs/focus-alone.txt`／`logs/latency-count10.txt`／`logs/mut-b-ac12.txt`／`logs/mut-d-focus.txt`／`logs/final-pkg.txt`／`logs/d22scan.txt`／`logs/d22scan-after.txt`）；仓外副本 `%TEMP%/33r4-clean/`（含 `ipv6probe/main.go`）——**只建不删**，全部留在盘上供复尺。
- ⛔ 本程没做的：没翻任何 AC 框、没动台账、没读没写 `frontend/**`／`design/**`、没碰三枚冻结件／阈值／golden／allowlist／`PLAN.md`／`specs`／`BUILD.md`／`SLO.md`、没跑 `go mod tidy`／`go get`、没放宽任何既有断言、没在仓内建 worktree。
