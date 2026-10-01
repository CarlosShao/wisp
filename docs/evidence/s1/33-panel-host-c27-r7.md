# 33-r7 面板宿主 C27 — 常驻建窗踩到别人留下的消息：层归属修复腿（写码腿）

票 33 `.scratch/wisp/issues/33-panel-host-c27.md`。本腿只修那一枚顺序依赖红：
整包跑时 `TestAC14AwaitedBindingReplyReachesThePage` 红、隔离单跑绿。
AC 框（`- [ ]`／`- [x]`）一枚没碰，勾与不勾归编排者／验收腿。

## 起手锚点（同发取，`date` + `git log -1 --format=%H` + `rev-parse --abbrev-ref HEAD`）

- 时刻：`2026-10-01 15:41:26 +08`
- HEAD：`416d9d56977a89a070e29e23391b1a2f02bdc9b5`（标题 `docs(33-r6 收尾两枚落盘)：A504＋它表尾的编排者标注……`）
- 分支：`dev`
- 台件目录：`.scratch/wisp/probes/33/r7/`（现读 `ls .scratch/wisp/probes/33/`＝a1 a2 h1 p1 r1 r2 r3 r4 r5 r6 ⇒ r7 未占）
- 证据件号：`ls docs/evidence/s1/ | grep -i r7` 零命中 ⇒ `33-panel-host-c27-r7.md` 未占
- 起手桌面枚数（`tasklist //FI "IMAGENAME eq msedgewebview2.exe"`，15:41）：6 枚（本机 owner 也可能自己开着 Edge ⇒ 只报数，不据此指控谁）

## 判据进度表

| 格 | 内容 | 状态 |
|---|---|---|
| ① | 最小确定性复现（谁先跑、留下什么、`bringUp` 读到什么）＋逐字读数＋命令 | 待填 |
| ② | 定性：产品缺陷 vs 测试卫生缺陷（两边说法都写＋我选哪支、凭什么） | 待填 |
| ③ | 修法＋反控（定向突变必须当场红，抄原始红句逐字） | 待填 |
| ④ | 整包逐名红册（PASS/FAIL/SKIP 三数＋逐名红，目标 rc=0 且红册空） | 待填 |
| ⑤ | 门禁四数（build / vet / d22scan / gofumpt） | 待填 |
| ⑥ | 我可能写错的条目（附"如果错了后果"） | 待填 |
| ⑦ | 判不动的地方（逐条 甲／乙／不做 ＋现量＋为什么判不了） | 待填 |

---

## ① 最小确定性复现（谁先跑、留下什么、`bringUp` 读到什么）

（待填）

## ② 定性：产品缺陷 vs 测试卫生缺陷（两边说法）

（待填）

## ③ 修法＋反控红句

（待填）

## ④ 整包逐名红册

（待填）

## ⑤ 门禁四数

（待填）

## ⑥ 我可能写错的条目（逐条，附"如果错了后果"）

- (a) **"整包那枚红的毒源＝`TestAC13BringUpSurvivesAReusedThreadQuit` 释放回池的那条线程"**：
  机制我用自己的台件复算过（同锁线程定植"活窗＋一枚没派的 `WM_CLOSE`"⇒ `bringUp` 3/3 panic，
  `.scratch/wisp/probes/33/r7/modes-grid.txt`），但"包内确实是这一枚用例把线程还了池、且 AC#14 落在它上面"
  是**归因**，凭的是具名 A/B（拿掉泄漏者⇒AC#14 绿 / 留着⇒红）＋时刻相邻，不是 TID 逐枚对上的。
  若错（真凶是别的还池点，或根本是 dataPath 复用/并发建窗），判据④仍会红，且我的修法（甲封点）不解决问题。
- (b) **"`WM_QUIT` 是闩锁（latched），过滤式 `PeekMessageW([WM_QUIT,WM_QUIT])` 在队列里还有别的消息时看不见它"**：
  读数支持（mode=dispatch 三发：drain removed=0、filtered found=false，紧接着的建窗照样 panic），
  但"闩锁"是**我对机制的命名**，不是 Win32 给我回的字段。若机制其实是别的
  （例如销毁过程在 drain 之后才异步投 quit），那"任何排干范围都不可靠"这句仍成立（因为读数是 removed=0 + panic），
  而我在 ② 里写的"闩锁"这个词要降格成假设。
- (c) **"整条队列排干会吃掉本次建窗自己的完成回调 ⇒ 挂在 `GetMessageW`"**：机制是假设。
  盘上事实只有两条：mode=wide 三发里 1 发 rc=2（`-test.timeout 120s` 到点，栈顶＝
  `pkg/edge/chromium.go:100` 的 GetMessageW，`modes/wide-2.txt`）；整包 `4be1d3f1` 那发
  `pre-create message purge hit its cap` removed=256 且 show 永不完成（`.scratch/wisp/probes/33/r6/fullpack2.log:640`，
  本腿未复跑＝引前人读数）。若"吃掉回调"这句错，两条读数仍然只支持"宽排干不可靠"，不支持我给的因。
- (d) **"出货拓扑今天不会中这一形（面板线程锁到底、队列建窗前必空）"**：我只读了 `panel_resident_windows.go:202`
  的锁而不解、`notify_windows.go:138-139`（message-only 窗、解锁前 `DestroyWindow`、全仓 `PostQuitMessage` 只有
  依赖的 `Terminate` 与我们自己的 `terminateOnThisThread` 两处）。**没**逐行读 `internal/ball` 的 ui-sta 与
  托盘那条腿。若还有别的"在面板线程上起过窗又投过 close"的路径，出货面也会中，且我的定性（测试卫生）要改。
- (e) **"产品侧 `bringUp` 拒绝建窗不会误杀合法状态"**：判据只认"队列里有一枚 `WM_CLOSE`"。
  `TestAC4PriorFocusSurvivesARefusedPanelSample` 在面板线程上先起了编辑器窗再 `showAndWait`，
  那一发合法地带着别的消息建窗——若过滤式 peek 在那一发里读到 `WM_CLOSE`，我就会把一枚绿用例改红。
  读数见 ⑤/④（整包若在那枚用例红，就是这支错了）。
- (f) **"`pumpAfterDestroy` 3/3 绿＝owner 侧排干是 sound 的修法"**：三发都是**同进程同线程连着建第二扇窗**。
  它没有覆盖"owner 排干后把线程还给池、再由**另一枚协程**拿去建窗"这一形（那是包内真实形状，
  受调度器支配，我只能靠整包名册验）。若差别要命，判据④会红。

## ⑦ 判不动的地方（逐条：甲＝按我读数／乙＝别的形／不做 ＋现量＋为什么判不了）

- **甲（我选的）**：层＝**线程的 owner**。谁的线程谁负责"还池／再建窗前把它泵到静"，
  出货侧已经是这一形（`loop` 锁到底）；包内唯一破口是测试把带活窗的线程还了池。
  现量＝`modes-grid.txt` 四形对照（none 3/3 panic、dispatch 3/3 panic、wide 2/3 建起但 +1 泄漏窗且 1/3 挂死、
  pumpAfterDestroy 3/3 建起且 windows 3→0）。
- **乙（复算过，不选）**：产码 `bringUp` 自己把线程弄干净。三种排干我都量过：窄形拦不住 close 引发的闩锁 quit
  （3/3 panic），宽形要么泄漏上一任的窗（`thread windows now=4`）要么吃掉本次建窗的消息挂死（1/3 rc=2 +
  fullpack2 的 cap=256），"派完再窄排"照样 3/3 panic。
  **我只保留了乙里可判定的那一半**：`bringUp` 把"队列里有没有 close"变成**能检查、能拒绝**的前提
  （检查＋具名拒绝＝读数支持；清理＝读数反对）。代价写在产码注释与 ③：脏线程上面板**不开**并留一句名，
  而不是半初始化的 controller（那枚能整机 segfault，见 ①）。
- **丙**：读数没指向别处，但**"为什么这条 M 被分给新协程"这一维我没判**——它属 Go 运行时调度，
  不是 cmd/wisp 能约束的；我只把"落到脏 M 上会怎样"变成可检查的。
- **不做**：⛔ 不碰依赖（`pkg/edge` 那三行没有 `if`）、⛔ 不 `go mod tidy`/`go get`/改 `go.mod`/`go.sum`、
  ⛔ 不动 `internal/**`、⛔ 不动 `shutdown.go`/`shutdown_hooks.go`、⛔ 不搬 `winlive`、⛔ 不加 `t.Skip`
  到新用例、⛔ 不放宽那 15 秒、⛔ 不把断言降级成 `t.Logf`、⛔ 不删 AC13/AC14 任何一枚、
  ⛔ 不给 `bringUp` 加"重试到绿"的循环（拒绝＝一次性具名错误，没有循环）、⛔ 不改票面 AC 框、
  ⛔ 不动台账／HANDOVER／票面（编排者独占）。
