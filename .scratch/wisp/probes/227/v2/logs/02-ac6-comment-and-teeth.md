# 227-v2 读数件②＝AC#6 两验（注释在场性⇄真轮询那发的去牙台⇄还原哈希）

台规＝★全部突变只打在仓外导出树 `D:\tmp\wisp\20261019\227v2\deliv-tree-ac5` 的物理字节上（⛔ 改跟踪文件、⛔ `-overlay`）；每发落刀**前**验旧锚命中＝1、落刀**后**验旧锚归零（MU-2/MU-3 另验新锚＝1）；每发跑完**当场还原**并比 sha256。pristine＝`72762bd5992479c7768d011b7126967556685952d777cc3fe2a5ff33896c602c` ⇄ `git cat-file blob HEAD:cmd/wisp/config_reload.go | sha256sum` **逐字同串**（blob `990d148b…`＝导出树⇄提交树同字节）。三发的 `RESTORED_SHA` 全＝pristine、`RESTORE_EQUALS_PRISTINE=yes`（驱动件 `mut-mu{1,2,3}.log` ⇄ `mu3-repeat-{2,3}.txt`；⛔ 事后补打）。

## 验① 「nothing polls」那句：在场、且仍是假话

- 在场性（量的）＝`HEAD:cmd/wisp/always_write_no_clobber_226_test.go` **第 43-44 行逐字**：
  `// The edit lands after assembly, while the process is running and nothing`
  `// polls: [app] theme, a key with nothing whatever to do with [fs].`
  ★跨行形（`nothing` 收在 `:43` 行尾、`polls:` 起在 `:44` 行首）＝单行 grep 打不中——`227-v1` 栽过的那把 `tr` 尺本腿没重犯，用 `git show | sed -n '43,44p'` 直读。
- blob＝`5fde793067116b14c2094b1b3deb4fefe9d80417` ⇄ 编排者收 bullet 记的同串⇄起手 `feac8f77` 同串＝**从票 223 接线至今没人动过它**。
- 假话半（读的＋量的）＝生产 1s tick 在 `cmd/wisp/config_reload.go:92`（`configReloadPollInterval = time.Second`）⇄ `:123`（`observe.Default.Spawn(… configReloadTick)`）⇄ `:861`（`run.go` 的 `assembleRuntime` 调 `startConfigReload`，`runTextTask:249` 走到）；★本腿控制发里该宿主真装配＝MU-1 红句把 `[audit] config: HOT-RELOAD state=armed tick=1s … goroutine=watchdog owner=config` 整行打了出来（见下）。⇒ **那句今天逐字在场、逐字为假＝要求①「可判且判未做」复认成立**。⛔ 本腿改它＝产码注释⛔ 我写面，具名交编排者（或续腿）落笔。

## 验② 常驻用例⇄226 用例的"牙"（四发一台）

对象两枚：`cmd/wisp/config_reload_223_test.go:333`（`TestTicket223RunArmsTheReloadTick`＝"真轮询下手改仍看得见"的常驻用例）⇄ `cmd/wisp/always_write_no_clobber_226_test.go:39`（`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`＝时间炸弹本体，`:102` 断 `rt.mgr.Config().App.Theme == "dark"`）。定向 `-run` 只装这两枚（`02-run.ps1 -Sel './cmd/wisp -run=…|…'`）。

| 发 | 落点（内容锚） | 锚尺 | 四数 RUN/PASS/FAIL/SKIP rc | 红名⇄红句逐字 |
|---|---|---|---|---|
| 控制（未突变） | ⛔ | ⛔ | `2/2/0/0 rc=0` | 两枚全绿＝⛔ 恒红正控过 |
| **MU-1** | `rt.reloadHandle = observe.Default.Spawn(… rt.configReloadTick)` 整行摘除（tick ⛔ 起） | 旧锚 1⇒0 | `2/1/1/0 rc=1` | 223 枚红：`config_reload_223_test.go:354: audit never carried "config: HOT-RELOAD state=applied" within 40s`（红 dump 里逐字带 `state=armed tick=1s`＝armed 过、applied ⛔ 来） |
| **MU-2** | `rep, err := rt.mgr.CheckAndReload()` 换成 `var rep/err`＋`MUB 227-v2 seam removed`（tick 醒着、从不重读） | 旧锚 1⇒0、新锚 1 | `2/1/1/0 rc=1` | 同 223 枚红（40.9s 预算耗尽那形） |
| **MU-3** | `configReloadPollInterval = time.Second` ⇒ `10 * time.Millisecond`（`MUC 227-v2`，同形状只加快） | 旧锚 1⇒0、新锚 1 | `1/0/1/0 rc=1` | ★**226 那枚红**：`always_write_no_clobber_226_test.go:103: the guarded write applied the hand-edited theme to the running process (got "light", want the dark this run loaded): reading a hand edit back is ticket 223`（复量×3 同句同 rc） |

**判读（逐句对应票面）：**
1. **常驻用例真有牙**＝MU-1（摘 spawn）⇄ MU-2（摘 CheckAndReload）两发落点不同、同一枚 223 用例都红＝票面要求的"真轮询下手改仍看得见"这一跳今天**有人钉、且摘了接缝就响**。成对两发落点不同（公共块突变条），少任一发不拿绿当凭据。
2. **226 那枚的绿仍站在时序那一侧**＝它在 1s tick 下控制发绿（本腿 `1.04s`，223-v2 那发 4.82s 也绿＝编排者更正句"⛔ 拿跑不满 1 秒当依据"在 1s 形下成立）；★但把同一条生产路径的 tick 只改快（10ms）、⛔ 动它任何断言⛔ 动任何接缝，它就红在 `:102` 那条方向断言上＝**它的绿依赖"轮询没在收口前赶到"这一运行学事实**＝时间炸弹**机制坐实**：宿主 1s 醒、手改与 `AddAllowedDir` 收口赛跑。票面 09-29 更正后的依据形（"没人钉"⇄`:102` 方向与 D36 热档相反）——本腿读数＝`app.theme` 确是热档（`internal/config/tiers.go:47 "app.theme": "hot"`、`manager.go:357-359 planApp` 热并），`:102` 期望 `dark`＝与热档相反的那条断言，1s 下侥幸绿、10ms 下必红 ⇒ "绿的原因是没人钉"这半⛔ 完全成立（223 那枚钉的是"看得见"，⛔ 钉"226 写回不 Adopt 内存"这一跳；两枚各钉各的），要钉死它得给 226 用例配一枚"tick 必赶得上"的确定性形状＝续腿的活（⛔ 本腿写面）。
3. 红句位置尺（⛔ 事后 grep 配错行）＝`go test` 把正文印在 `--- FAIL` **之前**：本腿全部按"FAIL 行往上取体"的切法取句（`07-attribute.sh` 那段就是这么写的）。

**还原自证**＝三发各自 `RESTORED_SHA=72762bd5…=pristine`；收尾再打一次全目录尺：导出树 `cmd/wisp/config_reload.go` ⇄ HEAD blob 同 sha256；跟踪文件工作树⛔ 被碰（`git status --porcelain -- cmd/wisp internal/config` 里本腿零枚，见件③尺）。
