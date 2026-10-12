# 票 295 · 落地腿 `295-r1` · 第 2 笔：`AC#2` 反形正控（删写那一发＝真跑，红句逐字）

时刻见各件头行时间戳（本笔落件＝2026-10-12 14:5x +08）。台面＝共享工作树（AC#1 已入库＝`d5ef9cb9`）。
凭据⛔"加了 atomic"这件事本身；凭据＝**这一发变异能让指名用例红，且红只在那一枚断言上**。

## 0. 三发是什么（一发都不⛔ 少）

| 发 | 内容（按内容锚，⛔ 绝对行号进件） | 目的 | 件 |
|---|---|---|---|
| pristine | 未变异（＝`d5ef9cb9` 的形状） | ⓐ 现有读 `started` 的用例全绿照跑 ＋ ⓒ 的正控制 | `logs/21-ac2-pristine-named.txt`、`logs/32-ac2-restore-and-pristine.txt` |
| **M1** | **删写**：把 `toggleMute` 体内 `	ra.started.Store(gate.Open())` 那一枚赋值删掉，其余一字不动（注释、`switch`、两枚 `Store/Load` 全留） | ⓑ 正控：指名用例必须红 | `logs/30-ac2-m1-deletewrite.txt` |
| M2（补充） | 把同一枚写换成 `	ra.started.Store(true)`（恒真） | ⛔ **不是 `AC#2` ⓑ 的正控**（ⓔ-3 具名：它测"恒真"那一形，不是"手势不写"那一形）；唯一用途＝现读 ⓓ-2 **#7 `t.Fatalf` 参数位**的输出形状 | `logs/31-ac2-m2-alwaystrue-shape.txt` |

变异内容确认尺（不是读码，是量的）：M1 那一发里 `git grep -nE '\.started\.(Load|Store)\(' -- cmd/wisp/resident_audio_windows.go` ⇒ **只剩 2 命中**（boot 的 `Store(true)` ＋ `posture()` 的 `Load()`），手势体内那枚写**不在名册里**＝删的就是它。

指名只许一枚＝`TestAC290BothMuteGesturesTurnTheGateAndBack`（`cmd/wisp/resident_mute_290_windows_test.go`）；`go test -run` 指名到用例级，逐字命令在各件头行。

## 1. M1 的红句（逐字，两处同句：命令 1 与命令 2 各打一枚）

```
    resident_mute_290_windows_test.go:125: mute-hotkey: the leg's derived mirror says nothing is running while the gate is open
```

- 包级行（命令 1）：`--- FAIL: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.079s` ＋ `M1_COMMAND1_RC=1`。
- 红**就红在 ⓔ-1 预判的那一枚断言**上：`t.Fatalf` 站在 `		if !ra.started.Load() {`（ⓓ-2 #5，本发现量 `:124`）的下一行 `:125`，句子里那枚 `resident_mute_290_windows_test.go:125` 就是行证。
- ⛔ 把功劳记给后面那枚：第一轮 `for` 里这枚断言**排在** `if ra.gate.Open() || ra.started.Load() {`（#6）之前；M1 的日志里 `:142`（#7 那枚 `t.Fatalf`）**一枚都没打**（尺＝`grep -c 'capture stayed handed over' logs/30-ac2-m1-deletewrite.txt` ⇒ 0）⇒ 红因是 #5，⛔ 是 #6/#7。

## 2. pristine 那一发（ⓐ 全绿照跑 ＋ ⓓ 那句"没放宽任何断言"）

| 发 | 逐字读数 |
|---|---|
| `logs/21` 命令 1（只跑指名用例） | `--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)`／`ok github.com/CarlosShao/wisp/cmd/wisp 0.087s`／`GO_TEST_RC=0` |
| `logs/21` 命令 2（`-run 'TestAC247\|TestAC290\|TestLive'`） | 顶层 **13 枚 PASS ＋ 1 枚 SKIP ＋ 0 枚 FAIL**（同发第二把口径含子项＝16 枚 PASS；两把⛔ 相加当一枚数） |
| `logs/32`（M1/M2 还原之后再跑一次指名用例） | `--- PASS: TestAC290BothMuteGesturesTurnTheGateAndBack (0.01s)`／`PRISTINE_AFTER_RESTORE_RC=0` |

SKIP 那枚＝`TestAC247LiveMicrophoneLevelsReachTheBallSeam`，日志逐字理由＝`resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0`。⛔ 我设 `WISP_LIVE_MIC`（真机件归机主的真机窗口），⛔ 把它当绿、⛔ 当"守住了"。
⛔ 放宽任何一条断言：本腿对测试件的改动只有"形状"（`ra.started` → `ra.started.Load()`）与 `%v` → `%t`，逐枚可见 `logs/22-ac1-diff.txt`。

## 3. pristine ⇄ M1 的名册作差（同行落盘，尺在 `logs/33-ac2-roster-diff.txt`）

同一把尺（抽 `--- PASS|FAIL|SKIP` 状态行、去时长、sort；顶层⇄子项同一把）：

```
名册枚数：pristine=17 行，M1=17 行
comm -12（两发逐字相同）=16 行
comm -3（只在一边出现）＝2 行：
  M1-only（变异那一发）:  FAIL TestAC290BothMuteGesturesTurnTheGateAndBack
  pristine-only（未变异那一发）: PASS TestAC290BothMuteGesturesTurnTheGateAndBack
```

⇒ 全场**只有一枚**用例的颜色被 M1 翻动，其余 16 行（含 3 枚子项）逐名逐色不变。

## 4. ⓔ-2 那 7 枚读点：逐枚标"对此变异不敏感"＋名册读数（⛔ 只写"读码推的"）

"不敏感"＝在 M1 那一发里，**含这枚读点的用例照样是它原来的颜色**（下面每一枚的状态都取自 `logs/30` 命令 2 的逐名册，与 `logs/21` 命令 2 逐字相同）。缺口审计请注意：`started` 看起来有 8 枚测试在守，**实际只有 1 枚有牙**。

| ⓓ-2 读点（内容锚，改后形状） | 件 | 所在用例（名册里的名） | M1 发里的状态（读数） | 为什么对它瞎 |
|---|---|---|---|---|
| #1 `	if ra.started.Load() {` | `resident_audio_247_windows_test.go` | `TestAC247VoiceDisabledBuildsNoCollector` | `PASS TestAC247VoiceDisabledBuildsNoCollector` | 该 rig `voice.enabled=false`，⛔ 构造出 `gate` ⇒ 用例⛔ 发手势；期望的是"boot 不许写 started" |
| #2 `	if ra.started.Load() {` | 同上件 | `TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice` | `PASS …OpensNoDevice` | 出厂默认 `mic_muted_default=true` ⇒ boot 走静音支，碰不到 `started`；不发手势 |
| #3 `			if ra.started.Load() {` | 同上件（表驱动子测试） | `TestAC247DeviceFailureShapesStillBootAndSayTheLoss` ＋ `/occupied`＋`/permission_denied`＋`/no_device` | 顶层与 3 枚子项**全 `PASS`** | 设备失败 ⇒ `gate.Open()` 永假；变异删的是手势写，这一路⛔ 发手势 |
| #4 `	if !ra.started.Load() {` | `resident_audio_247_live_windows_test.go` | `TestAC247LiveMicrophoneLevelsReachTheBallSeam` | **`SKIP …`＝本发零读数**（⛔ 当绿） | 读码依据：`WISP_LIVE_MIC` 未设 ⇒ 用例在 `:130` 早退；即便跑，它走 boot 的 `Store(true)` 支（未被 M1 碰）且不发手势。★这一枚**没有实测支撑**，具名为欠账 `Q-295-r1-live-skip` |
| #6 `		if ra.gate.Open() \|\| ra.started.Load() {` | `resident_mute_290_windows_test.go` | `TestAC290BothMuteGesturesTurnTheGateAndBack` | 用例整枚 `FAIL`，但**红不在这一枚**：`:141`/`:142` 零命中（`grep -c 'capture stayed handed over' logs/30…`=0） | 删写后 `started` 整轮停在 false ⇒ `false \|\| false` 不成立，分支根本不进。★它对"删写"瞎，对"恒 true"**有牙**（见下面 M2） |
| #7 `			t.Fatalf("… (open=%t started=%t)", gesture, ra.gate.Open(), ra.started.Load())` | 同上件（**格式串参数位**，⛔ 条件位） | 同上用例 | 同 #6：M1 发里一枚都没打 ⇒ 对此变异不敏感 | 它在 #6 的分支体里，#6 不进它就不打印 |
| #8 `	if ra.started.Load() {` | 同上件 | `TestAC290OutcomeIsReadOffTheGateNotOffTheRequest` | `PASS TestAC290OutcomeIsReadOffTheGateNotOffTheRequest` | 该 rig 的 source 开局就 `DeviceError(0x80070005…)` ⇒ `gate.Open()` 永假，写与不写都是 false；这枚用例的真牙在 `if ra.gate.Open()` 与 outcome 串上的 `设备未交接`／`Privacy & security` |

ⓔ-2 还有一行不占读点、但同一发里同样恒绿的：`TestAC290MutedDefaultStillDecidesTheBoot`（`PASS`，它读的是 `posture()` 的派生句子、⛔ 发手势直到最后一句 `muteGesture`）。

## 5. `AC#2` ⓓ：#7 那枚参数位的**输出形状**（验收腿现读的就是这两行）

M2（恒 `Store(true)`）那一发里，红落在 `:142`，日志逐字：

```
    resident_mute_290_windows_test.go:142: mute-hotkey: capture stayed handed over after muting (open=false started=true)
```

- 两枚布尔都打成裸 `true`／`false`；花括号尺＝`grep -c -e 'started={' -e 'started= {' logs/31-ac2-m2-alwaystrue-shape.txt` ⇒ **0 枚**。
- 反证如果只换类型不 `.Load()`：`atomic.Bool` 丢进 `%v` 会打 `{true}`——编译过、断言不红、**失败日志开始撒谎**；本发改了 `.Load()` 并把两个 `%v` 收成 `%t`（非布尔实参会被内建 vet 拦下＝更严）。这一枚**是** ⓔ-2 #6/#7 "对恒 true 有牙"的方向，⛔ 拿它当 `AC#2` ⓑ 的正控。

**两发的互斥读数（尺＝逐字 `grep -c`，本腿现跑）**：M1 那一发里 `resident_mute_290_windows_test.go:141` 与 `:142` 各命中 **0**（`capture stayed handed over` 也 0）；M2 那一发里 `:125`（`derived mirror says nothing is running…`）命中 **0**。⇒ 两枚断言各咬各的方向、互不顶名：#5 只在"手势不写"时红，#6/#7 只在"恒写 true"时红。缺口审计若看到"这一枚用例在两种变异下都红过"⛔ 就等于"两枚都在守"——名册只许按上面的方向记账。

## 6. ⓒ 的还原证据（⛔ 变异体留在树上＝当场作废）

`logs/32-ac2-restore-and-pristine.txt` 头两段：

```
worktree sha256   = 49dd9b13ba189a00b5b1022b2f84735e60029b86423f55672e07f9303d07edcf
HEAD-blob sha256  = 49dd9b13ba189a00b5b1022b2f84735e60029b86423f55672e07f9303d07edcf
worktree md5      = 244fb362e105fda6edfbc3b9227e8a4e
HEAD-blob md5     = 244fb362e105fda6edfbc3b9227e8a4e
git diff -- cmd/wisp/resident_audio_windows.go 的枚数 = 0
```

另三枚测试件同尺（工作树 sha256 ⇄ HEAD blob sha256 逐枚相等、`git diff --name-only` 各 0 行）。变异体自己的哈希也留档：M1＝`4a72c753…283c3cb`／md5 `b9ab5493…`；M2＝`ece60c4e…19b316`／md5 `4ad4e76b…`（⇒ 还原后的号与变异体的号**不**相等，这不是同一枚件）。还原后又复跑一次 pristine 指名用例＝`PASS`（见第 2 节第三行）。
