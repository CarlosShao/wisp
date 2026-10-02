33-p1 (第二轮, 2026-10-02) 普查件：三问真跑读数

## §0 起手锚（同发取数）

- date：2026-10-02 19:29:31 CST（本轮第一次取钟）；主件落笔时刻见各节。
- git log -1 --format="%H %ci"：61c526e0ca0bbda05a7e1e9e978574d03327d379 2026-10-02 19:25:44 +0800（读数期中途推进到 cc6eaa65＝244-r1 的 build.ps1 提交，与本腿写面无交集，见 §0.2）。
- git branch --show-current：dev
- git rev-list --count origin/dev..HEAD ＝ cnb/dev..HEAD ＝ 22（19:37:23 现量，仅此刻有效）
- git status --porcelain ＝ 383 行（本仓常态：design/** 那批 owner 未提交件＋别家台件；⛔ 不还原不提交不删）
- git status --porcelain -- cmd internal ＝ 1 行：M internal/config/tiers.go（★ 归属在飞腿 255-v1 的 shot-3 突变体：.scratch/wisp/probes/255/v1/backup/tiers.go＋md5-baseline.txt＋shot3-red-full.txt 三件为凭，逐字注释 "255-v1 shot-3 mutation: was \"hot\", temporary"；⛔ 本腿不碰、不还原、不提交它）
- tasklist 起手：wisp.exe / wisp.test.exe / go.exe 全零枚；msedgewebview2.exe＝12（本机基线，与 33-r9 终态 12 同值）
- 票 33 框现读：13 未勾 / 1 已勾（19:55 复量同值）＝本腿一枚未碰

### §0.1 上一枚 33-p1 腿（10-01）与本次的关系

- 盘上现量：p1 目录上一枚腿的五枚提交（e7eee03e→c9c3e323→cf2d90ea→c8764a53→30ca36b1）全在 dev 历史里（git ls-tree HEAD .scratch/wisp/probes/33/p1/ 25 枚路径全在），本轮起手 status 里那行 "?? .scratch/wisp/probes/33/p1/" 是 git 对「目录里混有已提交件＋新的未跟踪件」的折叠显示，不是前腿被撤销（git log --diff-filter=D -- p1/ 零枚删除提交）。
- 本次派单是同一腿名的续发（三问已换），本件只答本轮三问；前腿读数（R25-R47、§A/§B）原样有效，不重跑、不改写。
- 前腿证据件是 .scratch/wisp/probes/33/p1/probe.md（编排者 A502 已收），不是五节 census；本次按派单落 census.md，两件并存。

### §0.2 读数期 HEAD 漂移

- 19:39:58 现量 git log -1 ＝ cc6eaa65（244-r1: build.ps1 追加 -H=windowsgui）。我的探针只 import internal/observe 与 go-webview2 依赖，零个 import 面、零个读数受它影响；写面名册核对按最终一次现量（§收尾）。

## §1 三问答（每问都带真跑读数，日志逐字在 logs/）

### 问① 再入形状：常驻进程里面板宿主起来之后，第二次 NewPanelManager（或等价入口）会怎样

先复认生产面（19:43 现读，全部 file:line 现取）：
- NewPanelManager 本体＝cmd/wisp/panel_host_windows.go:175-177：纯构造器，只填 mu/w/hwnd/dataPath/assets/disp 七个字段，零 Win32 调用、零 COM 调用 ⇒ 结构体层面多造几个都无副作用（〔量，读码〕；物理后果靠下面真跑）。
- 「等价入口」的建窗一步＝PanelManager.bringUp（:230），它内部就是 webview2.NewWithOptions（:298）。PanelManager 在 package main，探针导不进 ⇒ 真跑按依赖原语复刻：同一 STA 线程上、第一扇窗还活着时再调一次 webview2.NewWithOptions（这正是 bringUp 的建窗步）。复刻偏差具名：探针没有 PanelManager 的 disp/assets/dataPath 包装与 staleCloseQueued 前置检查——前置检查只认"队列里有 WM_CLOSE"，我的线程队列干净，不会触发它。
- 生产拓扑里"第二次"的实形＝panel_resident_windows.go:210 loop 的 handOverPump（:273）把控制交给 w.Run()，此后 RequestShow 全走 w.Dispatch 投递（post() :300-323 按 wRef 路由）——不是"再造一个 manager"。今天常驻装配只调一次 startResidentPanel（resident_windows.go 装配根），"进程内第二枚 manager"在今天的生产装配里不可达。

真跑读数（载体 .scratch/wisp/probes/33/p1/q1/main.go，两发逐字在 logs/q1-secondmgr-1.txt / -2.txt）：

发1（19:52:27，tid=30424）：
  HOST tid=30424 com=sta
  L1 two constructors done: m1=0x154463480b00 m2=0x154463480b40 (pure struct fill, no window)
  L2 first create ok hwnd=0x1C80F1E ms=667
  NOTE M1-LIVE
  L2 second create ok hwnd=0x30A0F32 ms=561 (both windows now on tid=30424)
  NOTE2 M2-LIVE
  TEARDOWN m2 destroyed / TEARDOWN m1 destroyed
  FINAL panics=0 sink_writes=0

发2（19:52:42，tid=33616）：同形，first ms=690 / second ms=536，两枚 hwnd 不同、两枚绑定都回话（M1-LIVE 与 M2-LIVE 都到），panics=0。

⇒ 判定：**不是 panic、不是静默失败、也不是"拿到第二实例"——是第二枚完整可用的实例真建出来了**（新 hwnd、绑定回话各自到达、两扇窗同线程并存）。615/679ms、536/561ms 两发同量级。销毁顺序后建先毁＋泵净，进程自退，msedgewebview2 台数回到 12。
⇒ 附问②侧的发现（dup-manager 形，q2 探针）：第二扇窗起来后它的页面照样能回话（REPORT_TEXT w2-page-said="W2-LIVE"），且第一扇不被打扰。物理上"一进程多面板窗"在依赖层是通的——今天挡住它的是产品形状（C27 单窗纪律＋装配只造一条 panel-sta），不是依赖能力。

### 问② Go→页面回执：宿主起窗后 ExecuteScript / PostWebMessage* 从 Go 侧发消息，页面有没有路回话；panel_resync 现在还是零吗

API 面现量（19:38-19:39 现取）：
- 依赖导出面（go-webview2@v0.0.0-20260205173254-56598839c808）：
  * 高层 WebView 接口（common.go:26-59）可导出的只有 Eval（＝ExecuteScript）与 Init/SetHtml/Navigate；接口名册里零枚 PostWebMessage*。
  * pkg/edge/chromium.go:138-145 Eval ⇒ e.webview.vtbl.ExecuteScript.Call（直 COM 调，结果忽略）。
  * PostWebMessageAsString 在依赖里只有一处调用＝chromium.go:242，是依赖自己的 MessageReceived 回声（收到什么原样发回），**不是给调用者的 API**。
- 产码面：grep -rn "PostWebMessage|ExecuteScript" --include=*.go cmd internal ⇒ **零命中**（19:38 现跑）。
- panel_resync：全仓（--include=*.go --include=*.md）grep ⇒ **代码侧仍然零命中**；只有规格四处：SPEC-08:150/:163、SPEC-10:18、PLAN.md:2972。⇒ 09-28 那句"规格有名、代码零命中"**今天仍然成立**（19:38 现量）。

真跑读数（载体 q2/main.go，Run() 泵形＝今天出货的 panel-sta 形；页面自己回话是唯一判据——前腿 R25 教训）：

- mode=eval（w.Eval 直推，UI 线程上发）三发（logs/q2-eval-2/-4/-5.txt）：
  PUSH issued via w1.Eval (ExecuteScript) on UI thread at_ms=600/600/630
  REPORT_TEXT page-said="EVAL_SEEN title=P1Q2-OK-33 polls=3" at_ms=745~787
  ⇒ **ExecuteScript 推送到达页面，页面回话也回到 Go**（往返闭环 145~187ms）。
- mode=eval-disp（推送从另一协程经 w.Dispatch 进 UI 线程再 Eval）一发（q2-evaldisp-1.txt）：
  PUSH issued via Dispatch->Eval at_ms=718 / REPORT_TEXT EVAL_SEEN ... at_ms=873
  ⇒ **Dispatch 这扇跨线程门真把闭包送上了 UI 线程并执行**（这正是常驻 post() 路由依赖的那一跳）。
- 判语纪律：q2-eval-1/-3 两发是仪器缺陷发（watchdog 没 post WM_QUIT 进程不退），读数本身（EVAL_SEEN）一致，日志保留不删；判语只采 -4/-5/-evaldisp-1 三发干净退出形。

⇒ 一句话：**今天出货形（panel-sta + Run() 泵）下，Go→页 的 Eval/ExecuteScript 通道与 页→Go 的绑定回话通道都是通的**；这是前腿 R25/R26 差异（自泵 vs Run）落成产码形状之后的复核读数。**PostWebMessageAsString 不在可调用面上**（依赖只在内部回声处用），所以"发一条消息"今天唯一可用形状就是 Eval 注入 JS——panel.resync 若要落地，要么走 Eval 推 JSON，要么等依赖加 PostWebMessage 暴露面（⛔ 本腿不做选型）。

### 问③ 起窗耗时（票 33 冷 ≤1500ms / 热 ≤200ms 那族，真机现量）

载体 q3/main.go，量法＝出货同义（cold＝建窗＋首轮 页→Go 往返；hot-approx＝通道已热的再往返，缺 ShowWindow/Hide 的 Win32 半边，**热读数是近似，具名降档**）。四发（19:54:27→19:54:57，logs/q3-latency-1..4.txt）：

  发1: COLD create alone: 810.626 ms → COLD(create+roundtrip) = 1129.557 ms；HOT-approx 209.610 / 128.335 ms
  发2: create alone 1123.567 → COLD 1348.502；HOT-approx 301.482 / 279.556
  发3: create alone 850.303 → COLD 1038.696；HOT-approx 157.227 / 203.220
  发4: create alone 648.771 → COLD 793.898；HOT-approx 142.528 / 142.941

- 冷：4 发 793.9~1348.5 ms，全部 < 1500ms 预算；两发贴近上缘（1348/1129）。建窗单独一步 648.8~1123.6 ms，占冷读数大头（首轮往返约 190~320ms）。
- 热（近似）：142.5~301.5 ms，8 枚读数里 6 枚 <200、2 枚超（301.5/279.6/209.6 三枚里两枚超）。⇒ 热族在争用下贴预算线，别拿单发当定论（与 33-r7 整包那两枚 D32 热重显红同族：quiet 档 37/36/40ms，负载档 292ms）。
- 【同机可能有争用】标注：读数期本机 go.exe＝0 枚、wisp.test.exe＝0 枚（每发前后 tasklist 现量），但 msedgewebview2 常驻 12 枚（别的会话）＋机器上可能有别的腿在编排；三问每发前都打了 date。争用无法排除，只能声明现量环境。

## §2 与 A487 / A504 判语的差集

- A487（00:03）说"frontend/dist 只有 .gitkeep ⇒ 页面产物一包都没有，整条面板链交白卷"。**今天的差集**：链的 Go 半边不白卷了——33-r5/r7/r9 落了常驻 panel-sta＋Run() 泵＋具名拒绝出口（我量到：dispatchRaw 生产链在 panel_inbound.go:228 五构造器全接好；我的探针页面穿过 wispDispatch 同族绑定成功回话）。"白卷"那半仍然真（AC#12 归口未动，frontend/** 本腿⛔未读未写）。
- A504（15:37）判"33-r6 没修好：毒线程被还给池，AC#14 15 秒超时"。**今天的差集**：那条红在 A505 四发名册里已修（泵到空＋具名拒绝＋r9 出口），我四发 q2/q3 真跑里零枚 15s 超时形、零 panic；第一次建窗（tid 各异）全部正常返回。
- A502 裁定 P1/P2/P3 我全部复认成立：P1 的 Run() 泵形＝今天出货形（我的 q2 三问发就是跑在这个形上）；P2 的"判据必须页面自己报"照做；P3 的 STA 0x2 守卫＝panel_resident_windows.go:214 逐字在产码（我的探针同款初始化，四发 r 非 0 即拒）。
- 新增（前人没判过的格子）：① 第二枚 webview 实例在同 STA 线程物理可行（q1 两发）；② PostWebMessageAsString 不在依赖可调用面＝"panel.resync 落地通道"这一格现在有了 API 面读数（要么 Eval 推、要么改依赖）；③ 热族争用敏感的现量补充（8 枚读数 2 枚超 200ms）。

## §3 我可能判错的条目

1. 【复刻偏差】q1/q2/q3 都不是 PanelManager 本体（package main 导不进）；量的是 webview2.NewWithOptions＋Bind＋SetHtml 原语与 panel-sta 的线程形状。如果 PanelManager 的 disp/assets/dataPath 包装层本身有副作用（读 assets.go：BuiltinAssets 只读 embed.FS，无窗口副作用），结论不翻；但"第二次 NewPanelManager"在产品装配里的真实可达性我⛔没量（今天装配只造一枚），别把"物理可行"读成"产品允许"。
2. 【热读数是近似】出货热路径＝Hide 后 HotShow（含 ShowWindow/SetForeground 的 Win32 半边＋focus 采样）；我的 HOT-approx 只量 SetHtml→回话。两族的共同部分是"通道已热的往返"，Win32 半边我的读数不覆盖。别把我的热数直接对 D32 200ms 判线。
3. 【争用不控】q3 四发在 24 秒窗口内连跑，互相可能污染（前一发销毁泵没排干净的 msedgewebview2 子进程 teardown）。发2 那对最慢（1348/301）可能吃了发1 的尾气； quiet 环境复测归验收腿。
4. 【单发 panel_resync 零命中的边界】我 grep 的是 --include=*.go 全仓＋md 四处规格引用；frontend/** 两层禁令没碰——如果那半边哪天有页面代码，我的"零命中"只对 Go/文档侧成立。
5. 【q2-eval-1/-3 两发进程不退】那是我 watchdog 的仪器缺陷（第一次实现忘了 post WM_QUIT），不是被测物缺陷；EVAL_SEEN 读数与干净发一致。日志保留，判语只采干净退出那三发。
6. 【"第二次 NewPanelManager"题面的两读】若题面要的是"同进程再走一遍 newResidentPanelManager＋startResidentPanel"（两条 panel-sta 线程），我量的是同线程两窗形（更贴近"第二个 manager 实例"的字面）。两条线程那一形物理上没量（q1 注释里 L3 具名没做）；若编排者要那一形，另派一发。

## §4 量不到的格子

1. 真常驻进程（cmd/wisp 装配）里第二次 NewPanelManager 的产品后果——生产装配根本不造第二次；复刻只能给原语层物理读数。
2. 热路径的完整 Win32 半边（ShowWindow/SetForeground/focus 采样）——那是 PanelManager 方法，探针复刻缺 Show/Hide 语义；D32 热预算的权威读数归 TestPanelHostLatencyPercentilesAC2（真桌面档）。
3. panel.resync 落地后的行为——规格有名、代码零实现；我只能量"推送通道今天的物理能力"（Eval 通/PostWebMessage 不在面上），落地形状是选型，⛔ 不归探针。
4. CI/无桌面环境的表现——winlive 族只在真桌面；本机读数不延伸到 runner（33-r9 §8 R8 同口径）。
5. MTA panic 的 HRESULT 真值——前腿 ⑦-结-3 已判"成因未定值"，本轮没造新仪器，沿用。
6. 两扇窗并存时的焦点/前台行为——q1 只证了并存与各自回话，没动 SetForeground；AC#4 那族的地界，⛔ 不越。

## 收尾（盘上现量，19:55-19:56）

- 写面名册：本腿新增全部在 .scratch/wisp/probes/33/p1/**（q1/q2/q3 三枚 main.go＋logs/ 14 份日志＋本件）；⛔ 产码/规格/冻结件/票面 AC 框零触碰（git status --porcelain -- docs/PLAN.md docs/specs internal/observe/thresholds.go tools/d22scan/allowlist.txt .github/workflows/ci.yml ＝ 空；票面框 13/1 同起手）。
- cmd internal 现量＝1 行（tiers.go＝255-v1 的，起手就在，不是我造的）。
- GOFLAGS= go build ./... rc=0；sh scripts/d22scan.sh rc=0 clean（logs/d22scan-final-q2.txt：ban #8 cmd/=87、internal/=484）；go list ./... 含 scratch ＝ 0 枚。
- msedgewebview2 终态 12（与起手同值，零孤儿）；33p1q1/q2/q3 进程终态零枚。
- ⛔ 本腿只 commit 不 push；pathspec 显式；⛔ add -A / --amend / reset / rebase / stash / checkout . / clean；临时件只建不删。
