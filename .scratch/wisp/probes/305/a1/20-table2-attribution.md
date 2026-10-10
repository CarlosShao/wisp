# `305-a1` · `20` 表②：两枚症状分开归因（逐字回读四件＋CI 那件，再作差）

尺口径先写死：本表**⛔ 新取颜色**——所有颜色一律逐字读回编排者已有的五件读数件；本腿只加了**对象层源码尺**与**`docs/evidence/s1/` 既有读数尺**。
台面按票面规则具名（母仓＝带 `frontend/dist` 旧产物的共享工作树；clone＝`C:/Users/swq/tmp/…` 那份 ⛔ bundle 的干净检出；CI＝托管 `windows-2025-vs2026` runner）。

## A. 五发逐字回读（对拉表）

| 发（件） | 台面 / head | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | `TestAC14AwaitedBindingReplyReachesThePage` | `TestAC14GoSideEvalPushReachesThePage` |
|---|---|---|---|---|
| `probes/303/r1/20-targeted-before-red.txt`（腿，改前） | **母仓** @ `cf46c24` 系 | **FAIL(20.01s)** `:325` 逐字 `no report "ac13-probe" from the page within 15s (what DID arrive at the door: nothing at all)`；`:315` 逐字 `AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]`；`cold -1.0 ms` | **FAIL(20.05s)** `:821` 同形（`ac14r-0`，nothing at all） | **FAIL(20.02s)** `:866` 同形（`ac14-push`，nothing at all） |
| `probes/303/orch/r2-orch-targeted-after.txt`（编排者，改后） | **母仓** @ `807497c1` | **FAIL(1.30s)** `:327` 逐字 `AC#13 page answer (head 807497c1): "0" - 0 of 1 probe id(s) present in the live document` ＋ `:330` 逐字 `… the live document contains NONE of the 1 element ids the embedded entry declares (the page itself answered "0")`；`cold 1005.8 ms` | **PASS(0.55s)** `:825` 逐字 `… page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)` | **FAIL(0.54s)** `:867` 逐字 `AC#14 nail 2 (Eval push hop), page's own words: title=""` ＋ `:869` 逐字 `… the page reports its title as "", want "PUSHED-33R5-OK" …`；`cold 368.5 ms` |
| `probes/303/orch/r3-f718e9b6-three-nails.txt`（编排者，**回归之前**） | **干净 clone** @ `f718e9b6`（＝`fb2fb802` 的父） | **SKIP(0.00s)** `:315` 逐字 `Resolve(entry): panel: embedded assets are not built (run npm run build in frontend/)` ＋ `:317` 具名跳过句 | **PASS(1.26s)** 同句 `REPLIED,REPLIED,REPLIED` | **FAIL(0.54s)** `:867`／`:869` **与 r2 逐字同形**（`title=""`, want `"PUSHED-33R5-OK"`） |
| `probes/303/orch/r4-clone-at-HEAD-nails.txt`（编排者，**修复之后**） | **干净 clone** @ `807497c1` | **SKIP(0.00s)** 同上两行逐字 | **PASS(1.44s)** | **FAIL(0.57s)** 同一句逐字 |
| `probes/303/orch/r9-ci-after-red-roster.txt`（编排者，托管 CI） | **托管 runner** headSha `05db4bc6`，job `114194107792`，`test-windows` | **FAIL(0.62s)** `:327`／`:330` 逐字（`"0"` − `0 of 1 probe id(s)`；`cold 470.2 ms`） | **PASS(0.30s)** 逐字 `REPLIED,REPLIED,REPLIED` | **FAIL(0.37s)** `:867`／`:869` 逐字 |

四数/名册尺（`r9` 自带）：`cmd/wisp` 档改前 `RUN=398 PASS=278 FAIL=8 SKIP=2` → 改后 `RUN=398 PASS=279 FAIL=8 SKIP=1`；两把尺⛔ 同物，本表⛔ 相加、⛔ 互减。
⚠ 时长那一列按派单写法处理：改前 `20.01s`/`20.02s` 是 15s 死线走满（回执那一层断），改后 `0.62s`/`0.37s`/`0.54s` 是**回执到了而内容不对**——**时长⛔ 当分母用**，只用作"形状换了"的旁证。

## B. 作差结论（逐条带尺）

**B-1 症状②（nail2）⛔ 由 `fb2fb802` 造成——票面 判语① 的这一半，盘上支撑。**
尺＝同台面（干净 clone）两发对拉：`r3`（`f718e9b6`＝回归笔的父）与 `r4`（`807497c1`＝修复后）**同一句逐字红**（`title=""`，want `"PUSHED-33R5-OK"`），且 `nail1` 在两发上都绿。
⇒ 回归笔之前 nail2 就红；修复笔之后 nail2 仍红。**这一枚⛔ 动在回执通道那一跳上。**

**B-2 症状①（AC13）的颜色只有带 bundle 的台面量得到。**
尺＝`r3`/`r4` 逐字 `:315` `panel: embedded assets are not built` ＋ `:317` 具名跳过 ⇒ 干净 clone 上 AC13 **恒 SKIP**（⛔ 颜色）；母仓两发（`r1` 改前／`r2` 改后）**都 FAIL**，但红的是**两句话**（`no report … nothing at all` vs `the page itself answered "0"`）。
⇒ 票面第 18 行那句"AC13 在母仓⛔ 已知'改前绿'过"——**在票 303 那一对改前/改后的范围内为真**（改前那一发它没给出内容判语的机会，通道挡在前面）。

**B-3 ★但是：两枚症状都在**另一枚**带 bundle 台面上绿过——这与票面标题/判语①的读法冲突，具名报回。**
尺＝`grep -rn -E "PASS: TestAC13ColdStart|PASS: TestAC14GoSideEvalPush" docs/evidence/s1/`（件 `logs/l4`）：
- `docs/evidence/s1/33-panel-host-c27-r7.md:201` 逐字 `--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.97s)`（同件 `:241`、`:277` 另两发亦 PASS）
- `docs/evidence/s1/33-panel-host-c27-r7.md:207` 逐字 `--- PASS: TestAC14GoSideEvalPushReachesThePage (0.50s)`
- `docs/evidence/s1/33-panel-host-c27-r6.md:36` 逐字 `14:24:22  AC13 alone, count=1:  --- PASS: … (1.47s)`
- 那一发的身份（同件 `:10`／`:23`）：起手 HEAD＝**`416d9d56`**（2026-10-01），整包 **rc=0／PASS=240／FAIL=0／SKIP=0**，且名册里**逐枚真跑了面板家族**。
  ⇒ `SKIP=0` 里 AC13 在场 ⇒ **那一发必然在带 bundle 的台面上**（clone 上它只能 SKIP）。
⇒ **两枚症状都⛔ 是"从没绿过的功能"；两枚都落在"某时绿过、后来红"那一族。**

## C. 票面那条必答（"若某一枚在能判死的台面上从来没有绿过……"）——**显式答**

- **"能判死的台面"这把尺先定义**：对 AC13＝带 bundle 的台面（母仓／或 clone＋拷入产物）；对 nail2＝任何开出真窗的台面（它在带 bundle 与 ⛔ bundle 两把台面上都给判语）。
- **答：⛔ 一枚命中那条件。** 凭据逐枚：AC13＝`33-panel-host-c27-r7.md:201/241/277`＋`r6.md:36`（HEAD `416d9d56`，带 bundle）；nail2＝`33-panel-host-c27-r7.md:207`（同一发、同一台面）。
- ⇒ 票面标题里"从来没绿过"那一读法⛔ 成立；票面 判语① 只否掉了 **`fb2fb802` 那一支**，⛔ 否掉"回归"这一整个类。
- ⇒ **本腿⛔ 据此下任何 `AC#2` 判语**（那是编排者裁权）。本腿只把这一条摆明：**"降级为⛔ 开工（缺功能）"的前提在盘上⛔ 成立**；若按"缺功能"裁，`r7` 那三枚 PASS 就成了没人解释的反例。
- ⚠ **一票欠账要具名**：`r7` 的绿是 **10-01 那枚未跟踪 `frontend/dist` 字节**给的；今天的 1044 字节那份 mtime＝`2026-10-10 08:51`（`logs/l6`）。⇒ 同一把"母仓台面"尺的**第二个自变量已经漂过**，票面台面规则没覆盖这一轴。

## D. 归因窗口与〔仅腿报〕嫌疑笔 `70b00885` 的现量

**归因窗口（尺＝`git log --oneline 416d9d56..f718e9b6 -- cmd/wisp/panel_host_windows.go`，件 `logs/l5`）：**
`416d9d56`（三枚全绿）→ `f718e9b6`（nail2 已红、AC13 在该台 SKIP）区间里动过那枚文件的＝**三枚**：
`4658dbb6`（10-03，255-r3 尺寸取值闭包·编排者代提死腿半成品）、`7a0236b3`（10-01，33-r9 具名拒绝的出口）、`f7d28ef0`（10-01，33-r7 建窗线程前提）。
⚠ 这把尺射程**只有那枚 `.go`**；两枚症状还吃**未跟踪 dist 字节**与测试 harness 的窗数（`r6` 的组 A/组 B 就是窗数依赖红）⇒ **"哪枚 commit 动了哪枚文件"判不死**，本腿⛔ 指名。
⚠ 另具名一处台面纪律：`r7` 的**绿**在带 bundle 台面上、`r3` 的**红**在 ⛔ bundle 台面上 ⇒ 这俩**⛔ 同一枚台面**。⇒ 严格意义的"同台面绿↔红对照"**欠一发**（＝下表 R2）。

**嫌疑笔 `70b00885` 现量（尺＝`git show --stat --format` ＋ `git show 70b00885 -- cmd/wisp/panel_host_windows.go` 的逐字 `+/-/@` hunk，件 `logs/l5`）：**
- 身份逐字：`70b00885 2026-10-08 10:16:41 +0800 test(33-r10 AC#13 item 2): give the cold-start page handover a window-free content ruler`
- stat：`cmd/wisp/panel_host_windows.go | 41 +++-`、`cmd/wisp/panel_pageover_33r10_windows_test.go | 293 ++++++++++++++++++++++++++`（新档，纯测试）
- hunk 逐字（唯一的行为面变化）：从 `bringUp` 里**删**掉那 9 行（`rtMs := m.firstRoundTripLocked(ctx, t0)` / `if err := m.serveEntry(); err != nil { m.serveNotBuiltNoticeLocked() }` ＋搬走的那段注释），在**新函数 `coldStartPageHandover`** 里以**同样语句、同样顺序**加回，`bringUp` 的原位改成一句 `rtMs := m.coldStartPageHandover(ctx, t0)`。
⇒ **它⛔ 改任何一处"文档交接次序"**：`firstRoundTrip` 仍在前、`serveEntry` 仍在后；三枚 `SetHtml` 一枚没挪（`10-` 表 §A 第 9-11 行那三枚调用点在这发前后同一序列）。
⇒ **腿的嫌疑笔⛔ 归因**：`bringUp` 把探针文档当"最后活文档"那一族的引入处**⛔ 是 `70b00885`**（那一族在 33-r10 之前就已被修掉；本 commit 只把修好的形态搬了个家）。
⇒ 它**确实留下一处污染**：搬进 `coldStartPageHandover` 的注释与 `…test.go:330` 的判据串都还逐字写着"finished on the probe page / bringUp must hand the entry over AFTER the probe"——在 `a81c2980` 的源码序上⛔ 成立（`10-` 表 §D 第 1、2 条）。**那是文字过期，⛔ 行为回归。**

## E. 判不死的东西＋各欠哪一发（⛔ 本腿跑，交编排者）

| 缺的那一发 | 尺面/台面 | 它能判死什么 |
|---|---|---|
| **R1** | 母仓 targeted：`go test ./cmd/wisp/ -count=3 -v -run 'TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe|TestAC14AwaitedBindingReplyReachesThePage|TestAC14GoSideEvalPushReachesThePage'`（带 `PATH=third_party/sherpa-onnx:build`） | 两枚症状是**确定红**还是**窗数/顺序依赖红**（先例＝`r6.md` 组 A 全绿／组 B 命中目标红）。缺它 ⇒ "哪一枚先坏、会不会是同一枚 harness 干扰"判不死 |
| **R2** ★ | **同一枚台面**的绿↔红对照：`git clone` 到仓外（`C:/Users/swq/tmp/…`）检 `416d9d56`，再把母仓今天那枚 `frontend/dist`（1044 字节）整份拷进那份 clone，跑三枚 | 一次答两事：① 带 bundle 台面上今天还绿不绿得回来（坐实"回归"类）；② 自变量是**代码**还是**那枚未跟踪 dist 的字节**。⛔ 在共享工作树里 checkout |
| **R3** | （登记为做不到）取 10-01 时代那份 dist 字节 | 未跟踪产物⛔ 任何 commit ⇒ 台面第二轴**⛔ 复现**，只能由 R2 侧证 ⇒ 具名欠账，⛔ 假装能跑 |
| **R4** | `go test -tags winlive ./cmd/wisp/`（`panel_transport_live_35v2_windows_test.go:159` 是本仓**唯一另一枚真窗 `Dispatch`+`Eval` 面**，今天⛔ CI 档、⛔ 本机默认档） | 给"Eval 推不推得进"第二次**独立**测量；⛔ 与票面 `AC#4` 要求的 `go vet -tags winlive` 混成一发 |

## F. 一句话结论（表②）
**症状②⛔ 是 `fb2fb802` 造成的（clone 同台面 `r3`↔`r4` 逐字同句），但票面"从来没绿过"那一读法⛔ 成立——`33-panel-host-c27-r7.md`（HEAD `416d9d56`，带 bundle，整包 240-0-0）里 AC13 与 nail2 两枚逐字 PASS ⇒ 两枚都在"绿过、后来红"那一族；嫌疑笔 `70b00885` 现量＝纯 extract-function（同语句同顺序），⛔ 改任何文档交接，只留下一句过期诊断文字。**
