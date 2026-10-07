# 273-a2 — 机器实例普查（只读腿，零产码改动）

票面：`.scratch/wisp/issues/273-shipping-process-builds-the-state-machine-without-a-sink-so-every-d43-side-effect-falls-into-a-no-op.md`（整份已读，60 行）。
本腿射程＝§1 机器实例点名 / §2 逐枚副作用名追生产者 / §3 三层条数原文 / §4 落点三事实。**不裁决、不选形、不写产码、不跑任何 `go test`。**

---

## §0 起手锚

```
$ git rev-parse --short HEAD
15699a2f

$ git status --porcelain -- cmd internal docs
(空)

$ date
Wed Oct  7 10:24:58 CST 2026
```

---

## §1 全仓"机器实例"逐枚点名

### 1.1 先把尺定下来（三把，互相补漏）

**尺 A（构造形状全量，含测试与 `.scratch`）**
```
$ grep -rn 'statemachine\.New(' --include=*.go .
```
原始末 3 行（不截断，总 8 枚命中）：
```
./internal/models/bridge_test.go:16:	machine := statemachine.New(statemachine.Options{
./internal/models/bridge_test.go:86:	machine := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
./internal/models/handoff_window_109_test.go:62:	return statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
RC=0
```
命中 8 枚＝产码 2 ＋ `_test.go` 6（`.scratch` 里那些 statemachine 变异拷贝**贡献 0 枚 New**，它们只 import）。

**尺 B（同名噪声防漏：包别名写法穷举）**——先确认 import 有没有别名，再确认 `New` 只有限定名一种写法：
```
$ grep -rhn 'wisp/internal/statemachine"' --include=*.go . | sort | uniq -c
```
原始末 3 行：
```
      1 8:	"github.com/CarlosShao/wisp/internal/statemachine"
      1 9:	"github.com/CarlosShao/wisp/internal/statemachine"
      3 3:import "github.com/CarlosShao/wisp/internal/statemachine"
```
⇒ 全仓该包 import **一律无限定别名**，所以 `statemachine.New(` 就是唯一的构造调用形状。

**尺 C（绕过 `New` 的直接复合字面量——同族写法的第二形）**
```
$ grep -rn 'statemachine\.Machine{' --include=*.go .
RC=1        (零命中)
```
⇒ 没有任何"不经 `New` 手搓 Machine"的形状；`sink` 字段是 `Machine` 的私有字段（`machine.go:52`），包外也写不到。

**尺 D（`Sink` 那枚类型的名字，防"同名噪声"给假读数）**
```
$ grep -rn 'statemachine\.Sink' --include=*.go .
RC=1        (零命中)
```
⇒ 全仓（含测试）**没有任何一处写出 `statemachine.Sink` 这个类型名**。判"有没有人接"因此**不能**靠 `grep Sink`：仓里叫 Sink 的三枚都不是它——
`grep -rn 'Sink:' --include=*.go cmd internal | grep -v _test` 的产码命中是 `cmd/wisp/providers.go:194`（`storeHealthSink`＝`llm.HealthSink`）与 `cmd/wisp/run.go:998`（`consoleSink`＝`agent.Sink`），**两枚都不是状态机的 Sink**。

**尺 E（"哪枚构造点真带了 Sink 字段"＝票 AC#1 要的尺）**——New 点向后看 4 行取 `Sink:`：
```
$ grep -rn -A4 'statemachine\.New(' --include=*.go cmd internal | grep 'Sink:'
internal/models/bridge_test.go-18-		Sink:    func(e statemachine.Effect) { effects = append(effects, e) },
RC=0
```
⇒ **8 枚构造点里只有 1 枚带 Sink，且那 1 枚在测试里**。产码 2/2 不带。这条与票面 [10-07 09:5x] 更正一/更正二**同数**（测试 6 枚、其中 1 枚真传），我复跑独立得到。

### 1.2 逐枚点名

| # | 位置（`file:line`） | ⓐ 落在哪个函数／入口 | 那条链今天真有人调吗（调用点枚数） | ⓑ 填了 `Options` 哪几个字段 | ⓒ 活多久 | ⓓ 有没有 Sink |
|---|---|---|---|---|---|---|
| 1 | `cmd/wisp/models.go:303` | `(*modelStore).handOffModel`（`cmd/wisp`，出货二进制里那台） | 生产链 3 跳，逐跳 **1 枚**真调用点：`handOffModel` ← `models.go:293`（`modelsEnsure`）；`modelsEnsure` ← `models.go:127`（`cmdModels` 的 `case "ensure"`）；`cmdModels` ← `cmd/wisp/main.go:109`（`case "models"`，`func main`）。⇒ 只有跑 `wisp models ensure <id>` 这一条 CLI 腿会造它 | 只 `Initial: StateFirstRun`（`Sink`/`Timeouts` 未填） | **一次 CLI 调用**：`:304` `defer machine.Close()`，函数返回即死；不是进程生命周期 | **无** ⇒ `machine.go:64→:67` 缺省 no-op |
| 2 | `cmd/balldebug/main.go:188` | `func main()`（`cmd/balldebug`，开发期调试器） | `main` 即入口，无上游；下游机器被 `gesture`(`:591`)／`dispatch`(`:617`)／`syncMachineToBall`(`:630`) 使用，这三个由 `ball.Events` 回调（`:194-208`）驱动。⚠ **`balldebug` 不是出货件**：`scripts/build.ps1:129` 只 `go build ... ./cmd/wisp`；`balldebug` 仅 `scripts/dev/ball-cycle.ps1:30` 构建 | 只 `Initial: StateSleeping` | **一次 balldebug 进程**（dev 手跑，随进程退出） | **无** ⇒ 同样 no-op |
| 3 | `internal/ball/hotkey_live_test.go:382` | `TestLiveMuteHotkeyEndToEnd` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 4 | `internal/ball/interaction_live_test.go:54` | `TestLiveClickSummonsAndDragDoesNot` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 5 | `internal/ball/interaction_live_test.go:134` | `TestLiveConfirmingCancelAndEscReturned` | 测试；`winlive` 门控用例 | 只 `Initial` | 一次用例 | **无** |
| 6 | `internal/models/bridge_test.go:16`（跨 `:16-19`） | `newWalkRig` 测试夹具 | 测试夹具（`grep -c newWalkRig` 见下） | **`Initial` + `Sink`**（唯一一枚） | 一次用例 | **有** ⇒ 全仓唯一真收件人，把 `Effect` append 进切片 |
| 7 | `internal/models/bridge_test.go:86` | 同文件另一枚（`Initial: StateSleeping` 的负控形状） | 测试 | 只 `Initial` | 一次用例 | **无** |
| 8 | `internal/models/handoff_window_109_test.go:62` | `newWalkMachine()` 夹具 | 测试夹具 | 只 `Initial` | 一次用例 | **无** |

**小结（我这把尺的数，不是引 a1 的）**：机器实例 **8 枚**＝产码 **2** ＋测试 **6**；带 Sink 的 **1** 枚（测试 #6）；产码 **0** 枚带 Sink。
⇒ 出货进程（`wisp.exe`）里**只有 `wisp models ensure` 这一条腿会造出一台机器**；`wisp run`（`cmd/wisp/run.go`）与常驻腿（`runResident`，`resident_windows.go:33`／`resident_other.go:22`）**一台都没有**——常驻腿只把 `statemachine.State*` 当常量用（`resident_ball_windows.go:275` 是 `ball.Options{Initial: ...}`，不是 Machine；`resident_approval_windows.go:956` 是 `b.SetState(...)`）。

**"某符号零调用者"锚在调用形状**：上表的调用点枚数一律锚成 `标识符(` 或 `x.method(` 形状（`grep -rn 'cmdModels'`／`'modelsEnsure'`／`'handOffModel'` 的命中里，注释行与声明行我都逐条分辨后剔除；例如 `handOffModel` 的 7 枚文本命中里只有 `models.go:293` 是调用，其余 6 枚是声明与注释）。

---

## §2 那 50 枚副作用名，逐枚问"谁会让它响"

### 2.1 副作用名清单的尺（我自己现跑，不引 a1）

**尺 F（名字清单＝`table.go` 权威表本体）**
```
$ grep -o 'SideEffects: \[\]string{[^}]*}' internal/statemachine/table.go | grep -o '"[a-z0-9._-]*"' | tr -d '"' | sort
62 .scratch/wisp/probes/273/a2/logs/table-effect-names.txt      (出现次数，含重名)
$ ... | sort -u | wc -l
50                                                                (逐名去重)
```
原始末 3 行：
```
ui.one-click-restart
ui.voice-cancel-availability
RC=0（清单见 `logs/table-effect-names.txt`，50 枚去重名＝票面③那段引用的同一靶面）
```

**尺 F-防漏（证明"按行抽 SideEffects"不漏名）**——若某枚 `SideEffects` 跨行，抽行就会漏：
```
$ grep -c 'SideEffects:' internal/statemachine/table.go          -> 39
$ grep -c '\]string{' internal/statemachine/table.go             -> 39
$ grep -n '^\s*"[a-z]' internal/statemachine/table.go            -> RC=1（零命中）
```
⇒ 39 枚 `SideEffects:` ＝ 39 枚 `[]string{`，且表里没有"裸字符串续行"，所以每枚 SideEffects 都在单行上，尺 F 不漏。
（重名 12 枚：`approval.cancel-call`／`approval.dequeue`／`approval.replay-offer`／`asr.exclude-playback`／`audio.release-output`／`badge.decrement`／`mem.rss-verify-10s`／`panel.stream-push`／`session.scope-create`／`speech.load-vad-asr`／`tts.stop-bargein-400ms` ＝ `sort | uniq -d` 现量 11 行，其中 `audio.release-output` 出现 3 次；`62 - 50 = 12` 枚重复出现。）

**尺 G（"事件有没有生产者"＝锚在 `Dispatch(` 调用形状，⛔ 不锚在 `Sink` 字符串）**
```
$ grep -rn '\.Dispatch(' --include=*.go cmd internal | grep -v _test
```
原始末 3 行：
```
internal/models/bridge.go:59:		if _, derr := b.machine.Dispatch(statemachine.EvDownloadFailed, nil); derr != nil {
internal/models/bridge.go:64:	if _, derr := b.machine.Dispatch(statemachine.EvDownloadCompleted, nil); derr != nil {
internal/statemachine/machine.go:202:	_, _ = m.Dispatch(ev, nil)
RC=0
```
⚠ 同名噪声一枚：`cmd/wisp/panel_resident_windows.go:360` 的 `w.Dispatch(fn)` 是面板工作队列，**不是**状态机 Dispatch（它没有 `statemachine.Machine` 接收者；尺"全仓 `statemachine.Machine` 类型出现处"＝只有 `cmd/balldebug/main.go` 4 处 + `internal/models/bridge.go` 2 处，见下）。
```
$ grep -rn '\*statemachine\.Machine\|statemachine\.Machine\b' --include=*.go cmd internal | grep -v _test
cmd/balldebug/main.go:178 / :591 / :617 / :630
internal/models/bridge.go:24 / :33
```

**尺 H（逐事件穷举生产者，⛔ 不用扩展名猜分隔符；退码取自身）**——完整命令与逐名计数落 `logs/event-producer-scan.txt`（132 行）：
```
$ for ev in $(grep -o 'Ev[A-Za-z0-9]*' internal/statemachine/events.go | sort -u); do
    hits=$(grep -rn "\b$ev\b" --include=*.go cmd internal | grep -v '_test\.go' | grep -v '^internal/statemachine/')
    n=$(printf '%s' "$hits" | grep -c .)      # 计数取 grep -c，不取管道退码
    printf '### %s prod_nonstatemachine_hits=%s\n%s\n' "$ev" "$n" "$hits" >> logs/event-producer-scan.txt
  done
$ grep '^###' logs/event-producer-scan.txt
```
41 枚 `Ev*` 常量里**只有 7 枚**在生产码（`cmd`＋`internal`，排除 `_test.go` 与包内）有命中：
```
EvModelMissing=1      internal/models/bridge.go:40
EvDownloadFailed=2    internal/models/bridge.go:46 / :59
EvDownloadCompleted=1 internal/models/bridge.go:64
EvSummon=1            cmd/balldebug/main.go:600
EvVeto=3              cmd/balldebug/main.go:596 / :606 / :608
EvMuteKey=1           cmd/balldebug/main.go:603
EvInterrupt=1         cmd/balldebug/main.go:598
其余 34 枚 = 0        （零生产者）
```
⚠ 尺 H 的 `Ev[A-Za-z0-9]*` 还会抽出类型名 `Event`（命中 44），那是**词形噪声不是事件常量**，已从"41 枚常量"里剔除。
⚠ 防"只扫 cmd/internal"的漏尺：
```
$ grep -rn 'statemachine\.Ev' --include=*.go . | grep -v '^\./cmd/' | grep -v '^\./internal/'
RC=1（零命中）
$ find . -name '*.go' -not -path './.scratch/*' | grep -v '^\./cmd/' | grep -v '^\./internal/'
-> frontend/embed.go、scripts/spike/common/machine.go 等 11 枚 spike、tools/* 18 枚（`scripts/spike/common/machine.go` 是 Win32 消息循环，不是状态机；尺 A 已证它不 `New`）
```

### 2.2 三台"可能响的机器"的边界（本腿只做静态推演，⛔ 没跑过任何一发）

- **M-ship**＝`cmd/wisp/models.go:303`（出货 `wisp.exe`，无 Sink）。灌进去的事件只有尺 G 的 `bridge.go:40/:46/:59/:64` 三个名。状态闭包＝`FirstRun -> Downloading -> FirstRun | Error`（#2 无副作用、#37 完成行有 1 枚、#37 失败行无副作用）。⇒ 出货腿今天真投的副作用名＝**1 枚**。
- **M-debug**＝`cmd/balldebug/main.go:188`（dev 件，`build.ps1:129` 不构建它，无 Sink）。灌进去的事件＝`gesture()`(`:591-612`) 的 4 枚名 ＋ 机器自带定时器（`machine.go:202`←`timeouts.go:45-52`）。起始 `Sleeping`。静态闭包（按 `Table` 逐行匹配 From/Event/Guard）＝
  `Sleeping -EvSummon-> Listening(#4) -EvVeto-> Warm(#13) -EvSummon-> Listening(#29) | -90s EvWarmIdle-> Settling(#31) -3s EvSettleExpired-> Sleeping(#32，GuardFn 取 `Facts{}` 零值 ⇒ `!KwsLoaded` 成立) | Settling -EvSummon-> Listening(#33)`。
  到不了的态：`Armed/Muted`（要 #5 `EvKwsEnabled`／#8，全仓零生产者）、`Confirming`/`AwaitingApproval`（#17 零生产者）、`Thinking/Acting/Speaking`（#11/#15/#16/#18 零生产者）、`Downloading/FirstRun/Unconfigured/NoNetwork/WatchdogAlert/Queued/Stuck/Conversation`（无生产者或无入边）。
- **测试机**＝`bridge_test.go:16`（唯一带 Sink）。

### 2.3 逐名表（50 枚，⛔ 不是总数）

档位：**A**＝出货腿今天真投进 no-op；**B**＝dev 件 `balldebug` 静态可达、真投进 no-op（本腿未跑，属代码路径推演）；**C**＝事件有生产者、但今天到不了触发点；**D**＝触发事件零生产者（键盘/托盘/语音/面板/定时器/CLI 里没有任何一条真路径）。

| 副作用名 | 所在 D43 行 | 触发事件 | 事件入口分类 | 今天真路径（具名行／"零生产者"） | 档 |
|---|---|---|---|---|---|
| `model.verify-sha256-signature` | #37(`table.go:265`) | `EvDownloadCompleted` | CLI（`wisp models ensure`） | 有：`internal/models/bridge.go:64` ← `cmd/wisp/models.go:305 WireDownloading`，机器＝M-ship | **A** |
| `session.scope-create` | #4(`:52`)、#7(`:62`) | `EvSummon`／`EvWakeWord` | 单击球·快捷键／唤醒词 | `EvSummon` 仅 dev 件 `cmd/balldebug/main.go:600`；`EvWakeWord` 零生产者 | **B** |
| `speech.load-vad-asr` | #4、#7 | 同上 | 同上 | 同上 | **B** |
| `audio.discard-buffer` | #13(`:94`) | `EvVeto` | 否决词·`Esc`·单击球 | 仅 `cmd/balldebug/main.go:596`（Listening 分支） | **B** |
| `session.zero-load-resume` | #29(`:216`) | `EvSummon` | 单击球·快捷键 | 仅 `cmd/balldebug/main.go:600`，需已在 `Warm`（可由 #13 到达） | **B** |
| `settling.cancel-fallback` | #33(`:243`) | `EvSummon` | 单击球·快捷键 | 仅 dev 件，需已在 `Settling`（#31 到达） | **B** |
| `session.dispose-scope` | #31(`:228`) | `EvWarmIdle` | **机器自带定时器**（`timeouts.go:49` Warm 90s） | 定时器在 M-debug 的 `Warm` 会武装；出货腿走不到 `Warm` | **B** |
| `speech.unload-asr-tts` | #31 | `EvWarmIdle` | 定时器 | 同上 | **B** |
| `mem.free-os-memory` | #31 | `EvWarmIdle` | 定时器 | 同上 | **B** |
| `panel.destroy-or-hide` | #31 | `EvWarmIdle` | 定时器 | 同上 | **B** |
| `mem.rss-verify-10s` | #32(`:233`、`:239` 两枚行共用 #32) | `EvSettleExpired` | 定时器（`timeouts.go:50` Settling 3s） | 同上（M-debug 的 `Settling`） | **B** |
| `error.ack` | #38(`:270`) | `EvErrorAck` | 定时器（`timeouts.go:52` Error 10s） | 出货腿**真进 `Error`**（`bridge.go:46/:59` 投 `EvDownloadFailed`，#37 失败行 → `StateError`），但 `cmd/wisp/models.go:304` `defer machine.Close()` 在 `:321 return 1` 时拆表（`machine.go:169-174`→`stopTimerLocked`），10s 不可能到期 | **C** |
| `tools.execute` | #21(`:151`) | `EvConfirmExpired` | 定时器（`timeouts.go:48` Confirming 3s） | 定时器存在，但 `Confirming` 唯一入边是 #17 `EvApprovalNeeded`＝零生产者 | **D** |
| `approval.cancel-call` | #24(`:169`、`:173` 两枚行) | `EvApprovalDenied`／`EvApprovalTimeout` | 原生卡（面板·托盘）／定时器 300s | `EvApprovalDenied` 零生产者；`EvApprovalTimeout` 需已在 `AwaitingApproval`（#17 零生产者） | **D** |
| `approval.replay-offer` | #24 两枚行 | 同上 | 同上 | 同上 | **D** |
| `config.persist` | #1(`:44`) | `EvOnboardingCompleted` | 引导界面（面板）／CLI | 零生产者 | **D** |
| `db.create-schema` | #1 | 同上 | 同上 | 零生产者 | **D** |
| `kws.load` | #5(`:56`) | `EvKwsEnabled` | 配置·托盘开关 | 零生产者 | **D** |
| `kws.pause` | #7(`:62`) | `EvWakeWord` | 语音（KWS 命中） | 零生产者 | **D** |
| `kws.stop-inference` | #8(`:66`) | `EvMuteKey` | 静音键·托盘 | `EvMuteKey` 仅 dev 件 `:603`，但需已在 `Armed`；`Armed` 入边是 #5（零生产者），dev 件从 `Sleeping` 起 ⇒ `Sleeping + EvMuteKey` 无行＝非法转移（`machine.go:136-141`） | **D** |
| `kws.keep-alive-alert` | #9(`:71`) | `EvWatchdogOverrun` | 看门狗 | 零生产者 | **D** |
| `audio.stop-capture` | #11(`:87`) | `EvVadStop` | 语音（VAD 判停） | 零生产者 | **D** |
| `asr.punctuate` | #11 | 同上 | 语音（标点模型） | 零生产者 | **D** |
| `input.taint-mark` | #11 | 同上 | 语音·输入通道 | 零生产者 | **D** |
| `error.device-name` | #14(`:98`) | `EvAudioDeviceLost` | 音频设备回调 | 零生产者 | **D** |
| `panel.stream-push` | #15(`:103`、`:109` 两枚行共用 #15) | `EvFirstToken` | LLM 首 token | 零生产者 | **D** |
| `error.classify-retry` | #16(`:113`) | `EvBrainFailed` | LLM／网络失败 | 零生产者 | **D** |
| `confirm.countdown-start` | #17(`:120`) | `EvApprovalNeeded`（L1 guard） | 工具门控 | 零生产者 | **D** |
| `approval.enqueue-c18` | #17(`:125`) | `EvApprovalNeeded`（L2 guard） | 工具门控 | 零生产者 | **D** |
| `agent.inject-reminder` | #19(`:138`) | `EvRepeatThreshold` | Agent 循环 | 零生产者 | **D** |
| `panel.queue-waiting` | #20(`:147`) | `EvPathConflict` | 路径锁（C20） | 零生产者 | **D** |
| `agent.cancel-tool-call` | #22(`:155`) | `EvVeto`（From `Confirming`） | 否决 | `EvVeto` 有 dev 件生产者，但 `Confirming` 不可达 | **D** |
| `ui.voice-cancel-availability` | #22 | 同上 | 同上 | 同上 | **D** |
| `approval.dequeue` | #23(`:160`、`:165` 两枚行) | `EvApprovalGranted` | 原生侧允许（面板·托盘） | 零生产者 | **D** |
| `badge.decrement` | #23 两枚行 | 同上 | 同上 | 零生产者 | **D** |
| `mic.keep-off` | #26(`:188`) | `EvSpeakDone`（not-conversation） | TTS 播报完 | 零生产者 | **D** |
| `panel.keep-alive` | #26 | 同上 | 同上 | 零生产者 | **D** |
| `tts.stop` | #27(`:192`) | `EvInterrupt` | 快捷键·单击·否决词 | `EvInterrupt` 有 dev 件生产者 `:598`，但需已在 `Speaking`（入边 #15/#18/#25 全零生产者） | **D** |
| `audio.release-output` | #27、#41(`:201`、`:207` 两枚行) | `EvInterrupt`／`EvBargeIn` | 打断·AEC 检出 | `EvInterrupt` 同上不可达；`EvBargeIn` 零生产者 | **D** |
| `tts.stop-bargein-400ms` | #41 两枚行 | `EvBargeIn` | 语音（AEC） | 零生产者 | **D** |
| `asr.exclude-playback` | #41 两枚行 | `EvBargeIn` | 语音 | 零生产者 | **D** |
| `mic.enable` | #28(`:212`) | `EvSpeakDone`（conversation） | TTS＋会话模式 | 零生产者（且 `Conversation` 无入边，见 §3ⓓ） | **D** |
| `conversation.ring-solid` | #28 | 同上 | 同上 | 同上 | **D** |
| `conversation.privacy-confirm-l2` | #30(`:224`) | `EvConversationOn` | 配置·面板开关 | 零生产者 | **D** |
| `conversation.entered` | #30 | 同上 | 同上 | 零生产者 | **D** |
| `conversation.suspend` | #42(`:247`) | `EvTaskIntent` | 语音／文本循环（D47） | 零生产者；且 `From StateConversation` 在表里**没有任何一枚 `To: StateConversation`** ⇒ 该行的 From 今天进不去（尺：`grep -n 'To: *StateConversation' internal/statemachine/table.go` → RC=1） | **D** |
| `ctx.carry-c7-d47` | #42 | 同上 | 同上 | 同上 | **D** |
| `task.ctx-keep` | #34(`:253`，`From AnyState`) | `EvNetworkDown` | 网络探测 | 零生产者（`AnyState` 通配也救不了"没有事件"） | **D** |
| `ui.one-click-restart` | #35(`:257`，`AnyState`) | `EvWatchdogFatal` | 看门狗 | 零生产者 | **D** |
| `diag.record-stack` | #36(`:261`，`AnyState`) | `EvPanicRecovered` | panic recover | 零生产者 | **D** |

**逐名档位计数（我这几把尺的数）**：A＝**1**、B＝**10**、C＝**1**、D＝**38**，合计 **50**＝尺 F 的去重枚数。
⇒ 票面③那段"1 真投 / 1 结构可达 / 48 没被投"的**形状我复核为一枚不差地同意方向，但档位划分不同**：我这边出货腿真投 1 枚（同），`error.ack` 归 C（同"结构可达"），其余 **49** 枚里我另分出 **10** 枚属 dev 件 `balldebug` 可达（票面把这一段并进了"48 枚根本没被投"）。**"48"这枚数在我这三把尺上不成立**（出货腿口径下是 49，含 dev 件口径下没投的是 39）——按票面 AC#0 的要求，我以现跑的尺为准并在此具名顶回。

---

## §3 条数：三层原文摆齐（本腿不裁决、不改任何文件）

ⓐ **`docs/PLAN.md`（D43 权威表本体）关于条数的原话**：
- `docs/PLAN.md:3057`：`20 态（§2）· 40 条转移。**未列出的转移一律非法**（未定义即停，D22 闸门③）。`
- `docs/PLAN.md:1748`：`| **`docs/STATE_MACHINE.md`** | **第四轮新增**：D43 的 20 态 + 40 条转移表。**必须单独成文** —— C12 冻结的就是它，散在方案里 agent 会自行简化 |`
- `docs/PLAN.md:3151`（§16 差异总览里的"驳回"行）：`… **态数不是复杂度指标，转移数才是**（40 条转移，20 态，比例正常） …`
- 表本体行数尺：`awk 'NR>3055 && NR<3102 && /^\| [0-9]+ \|/ {c++} END{print c}' docs/PLAN.md` → **40**（编号 1–40 逐枚连续，末三行原文＝`docs/PLAN.md:3098/3099/3100` 的 `| 38 |`／`| 39 |`／`| 40 |`）。**PLAN 里没有 #41/#42。**

ⓑ **`docs/specs/SPEC-08-ui-ball-panel.md` §3 的 `#41`/`#42` 逐字**：
- `docs/specs/SPEC-08-ui-ball-panel.md:85`：`## 3. 状态机权威转移表（C12/D43——20 态 42 条转移；未列出的一律非法）`
- `:119`：`| 41 | `Speaking` | 播报中检出用户语音（Path C，AEC） | `Listening` | 停 TTS 并释放输出，≤400ms；播报音频不得进 ASR（D47） |`
- `:120`：`| 42 | `Conversation` | 任务意图（如「帮我做 X」） | 交回文本循环（Path T） | realtime 大脑会话挂起或结束；上下文经 C7 携带（D47） |`
- 该节表本体尺：`awk 'NR>=87 && NR<=140 && /^\| [0-9]+ \|/ {…print 编号…}' docs/specs/SPEC-08-ui-ball-panel.md` → **42 行**，编号序列原文＝
  `1..30, 41, 42, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40`（**#41/#42 插在 #30 与 #31 之间，非数值序**）。

ⓒ **代码里写着 `D43: 41`／`D43: 42` 的行，逐字**：
- `internal/statemachine/table.go:195`：`		D43: 41, From: StateSpeaking, Event: EvBargeIn, To: StateListening, Listen: PhaseInSession,`
- `internal/statemachine/table.go:204`：`		D43: 41, From: StateSpeaking, Event: EvBargeIn, To: StateListening, Listen: PhaseConversation,`
- `internal/statemachine/table.go:246`：`		D43: 42, From: StateConversation, Event: EvTaskIntent, To: StateThinking,`
- ★**如实报**：`D43: 41` **确实出现在两枚不同的行上**（`:195` 与 `:204`；同 `From`／同 `Event`／同 `To`，差别只在 `Listen` 相位与 guard）。尺＝
  `grep -o 'D43: [0-9]*' internal/statemachine/table.go | sort -t' ' -k2 -n | uniq -c` → 末三行原文：
  ```
        2 D43: 41
        1 D43: 42
        2 D43: 40
  ```
  全仓 `D43:` 命中：**56 枚 slice 行 / 42 枚去重编号（1–42）**；其中共用同一编号的行还有 #10(×2)、#12(×3)、#15(×2)、#17(×2)、#18(×2)、#19(×2)、#23(×2)、#24(×2)、#25(×2)、#32(×2)、#37(×2)、#40(×2)——即"一枚编号多行"是本表的既有形状，不是 #41 独有。
- 实现表的自述文字（同一文件里，⚠ 两处口径并不一致）：
  - `table.go:3`：`// The authoritative D43 transition table (SPEC-08 §3 / PLAN.md D43): 20` / `:4` `// states, "anything not listed is illegal" …`
  - `table.go:6-9`：`// Numbering reconciliation: the frozen D43 table carries 40 rows; SPEC-08 §3` / `// (the ticket's line-by-line authority) appends two numbered rows to it -` / `// #41 (Path C barge-in) and #42 (Conversation task intent) - for 42 rows` / `// total. All 42 are implemented and pinned by the row-presence test.`
  - `table.go:11-13`：`// Table-driven by contract … A few rows list two destinations in their To column (guard-branched); those are encoded as adjacent entries sharing one row number.`（**这是文件自己对"同编号多行"的既有约定，不是我这腿的新解释**）
  - `table.go:29`（字段注释，⚠ 与表体共存）：`D43         int   // row number in SPEC-08 §3 (1..40)`
  - `internal/statemachine/doc.go:6-7`：`//   - the authoritative D43 transition table (table.go: the 40 frozen rows` / `//     plus SPEC-08 §3's appended #41/#42) and the per-state timeout table`
  - 测试钉：`internal/statemachine/table_test.go:15-16` `if r.D43 < 1 || r.D43 > 42`；`:26` `t.Errorf("got %d distinct rows, want 42 (40 D43 + #41/#42)", len(seen))`

ⓓ **顺手量到、与条数同层但形状不同的一枚（只报形，不裁决）**：尺
`grep -n 'To: *StateConversation' internal/statemachine/table.go` → **RC=1（零命中）**，而 `D43: 42`（`table.go:246`）的 `From` 正是 `StateConversation`。⇒ 实现表里 #42 那一行的 `From` 没有任何一枚行的 `To` 能到达。SPEC-08:118 的 `#30` 把 `To` 列写成 ``Conversation`→`Listening`，而 `table.go:219` 的 `To` 只有 `StateListening`（`table.go:220-223` 的注释对这一点是明说的）。三层原文到此摆齐，**是否算缺陷由编排者裁**。

---

## §4 落点形状的三个事实（本腿不选形）

### ⓐ 装配根在哪

| 腿 | 入口 | 装配函数（具名行） | 机器造点 | 该装配函数今天被调几处 |
|---|---|---|---|---|
| `wisp models ensure`（出货 CLI 子腿） | `cmd/wisp/main.go:109`（`case "models"`） | `cmdModels` `cmd/wisp/models.go:104` → `modelsEnsure` `:262` → `(*modelStore).handOffModel` `:302` | `:303`（唯一造点，逐跳 1 枚调用点） | `cmdModels` 1 枚产码调用；`modelsEnsure` 1 枚；`handOffModel` 1 枚 |
| `wisp run`（干活 CLI 腿） | `cmd/wisp/main.go:91` → `cmdRun` | `cmd/wisp/run.go:990 (*agentRuntime).execute`（`:998` 装 `agent.Options.Sink`） | **0 枚**（尺 A 全仓无） | — |
| 常驻腿（无参数 GUI，`wisp.exe` 双击） | `cmd/wisp/main.go:66 runResident()` | `cmd/wisp/resident_windows.go:33 func runResident()`；同文件装配行＝`:66 installLogSink`、`:132 newResidentApprovalWithConfig`、`:217 startResidentBall(...)`；非 Windows 是同名的 stub `cmd/wisp/resident_other.go:22` | **0 枚** | `runResident` 产码调用点 1 枚（`main.go:66`） |
| dev 件（⛔ 非出货件） | `cmd/balldebug/main.go:87 main()` | 同函数内装配（`:188` New、`:189 ball.New`、`:194-208 Events`） | `:188` | 仅 `scripts/dev/ball-cycle.ps1:30` 构建；`scripts/build.ps1:129` 只 `go build ... ./cmd/wisp` |

⇒ 事实：**出货进程里"有机器可装"的位点只有 `cmd/wisp/models.go:302-305` 一处**；常驻腿与 `wisp run` 里连机器都没有，"往哪儿注 Sink"这两条腿目前是空集（要先有机器）。

### ⓑ 既成的"由装配根注入"先例（票 246 的真身在哪几行）

- 票面自述：`cmd/wisp/resident_approval_windows.go:5`——`// The approval gate held by the resident process (ticket 246, 乙形：装配根注入).`
- 被注入方的形状（参数就是"函数值＋nil 合法"）：
  - `cmd/wisp/resident_ball_windows.go:63` `type escVetoFunc func() string`
  - `cmd/wisp/resident_ball_windows.go:256` `func startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc, hotCfg ..., hotReload ..., hooks ...ballHostHook) *residentBall`
  - `:257` `rb := &residentBall{cancelHosted: onCancelEsc != nil}`；消费点 `:430-437 recordCancelHotkey(onCancelEsc)`（nil ⇒ 只登记不可用）
  - 注释钉：`:34` "It is handed in as a function value by startResidentBall's parameter, so…"、`:237` "onCancelEsc is the injected cancel executor … nil is legal"、`:351` "D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)"
- 装配根递进去的那一行：`cmd/wisp/resident_windows.go:217 rb := startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, …)`（`:132` 先 `ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`）。
- 同族第二先例（**同一字段名 `Sink`，装的都是别的包的 Sink**，即"由装配根往 Options 里塞一枚 Sink"这条形状在本仓已有两枚产码真身）：
  - `cmd/wisp/run.go:998 Sink: consoleSink{out:…, stream:…, publish: rt.publishPanelSnapshot}`（`internal/agent/sink.go:63 type Sink interface`；缺省不塞＝`internal/agent/sink.go:67-68 NopSink … used when no UI is wired`）
  - `cmd/wisp/providers.go:194 Sink: storeHealthSink{store: store}`（`internal/llm/probe_health.go:117 HealthSink`、`:208 Sink HealthSink`）
- ⚠ 反面事实：状态机这枚 `statemachine.Sink`（`machine.go:34`）**没有任何一枚先例被写出过类型名**（尺 D 全仓 `grep -rn 'statemachine\.Sink'` → RC=1，含测试），唯一的真收件人是函数字面量 `internal/models/bridge_test.go:18`。

### ⓒ 新增注入会不会新开一条依赖边（`go list` 现量，⛔ 不凭印象）

```
$ go list -f '{{.ImportPath}}: {{.Imports}}' ./internal/statemachine
github.com/CarlosShao/wisp/internal/statemachine: [fmt log/slog sync time]
```
⇒ 状态机包**零 wisp 内部依赖**。"包内自带默认 Sink"（票 AC#2 的乙形）若要在包里真发出去，`ball`/`panel`/`agent` 全都要新开边，且 `internal/ball` 与 `internal/models` 已经各自 import `statemachine`（下量），包内自带真收件人会**成环**。
```
$ for p in ./internal/models ./internal/ball ./internal/panel ./internal/agent; do go list -f '{{.ImportPath}} {{range .Imports}}{{.}} {{end}}' $p | grep CarlosShao; done
```
原始末 3 行（逐包 wisp 内部依赖，已剥前缀）：
```
./internal/models -> internal/observe internal/statemachine
./internal/ball   -> internal/observe internal/statemachine
./internal/panel  -> internal/projctx internal/risk internal/streamkey
./internal/agent  -> internal/config internal/llm internal/memory internal/observe internal/winsec
```
⇒ `internal/models`（`WireDownloading` 造桥侧）**不 import** `ball`/`panel`/`agent`；在这里注一枚真 Sink 就要新开 `models -> ball`/`models -> panel` 之类的边。
```
$ go list -f '{{range .Imports}}{{.}}{{"\n"}}{{end}}' ./cmd/wisp | grep wisp
（21 枚，含 internal/agent、internal/ball、internal/models、internal/panel、internal/session、internal/statemachine、internal/observe …）
```
原始末 3 行：
```
github.com/CarlosShao/wisp/internal/secret
github.com/CarlosShao/wisp/internal/session
github.com/CarlosShao/wisp/internal/statemachine
github.com/CarlosShao/wisp/internal/tools
```
⇒ 在 `cmd/wisp` 装配根里注一枚真 Sink（票 AC#2 的甲形落点），被调的 `ball`/`panel`/`agent` **已经在 `package main` 的直接 import 集里**，这条注入**不需要新增 import 边**（尺已跑）。另有一条必须记的事实：`cmd/wisp` 装机器的那条腿（`models.go`）今天**没有** `ball`/`panel` 的任何窗口或事件循环在场，而常驻腿（`resident_windows.go`）有窗口却**没有机器**——"零新边"只说明 import 图上不加箭头，不说明那两条腿能直接对接。

---

## §5 未做完／没量的格子（照实报，⛔ 不写成做完了）

1. **§2 的 B 档 10 枚＝静态可达推演，没有任何一发运行时读数**：`winlive` 未批、本腿一枚 `go test` 都没跑。若编排者要"balldebug 真响过这 10 枚"的凭据，本腿给不了。
2. **§2 没有把"面板/托盘/键盘/语音"入口逐枚追到 OS 事件源**：我只量了"事件常量的生产者枚数＋具名行"。托盘 `OnTray*` 我读到的是 `cmd/balldebug/main.go:199-208` 的 `fmt.Println` stub（不投事件），出货常驻腿的托盘回调链**没逐枚追**。
3. **§1 的 8 枚里没有查"机器被别处经接口值持有"的形状**：`internal/models/bridge.go` 用的是具体类型 `*statemachine.Machine`，尺"全仓 `statemachine.Machine` 出现处（非测试）"＝ balldebug 4 枚 + bridge 2 枚；但我**没有**扫"把机器塞进某个 `any`/结构体字段再传"的形状。
4. **§3 只做三层原文摆齐**：`#41/#42` 有没有获人工批准（票面末段那格）＝不在本腿射程，我**没去台账查 `A##`**。
5. **`error.ack` 归 C 档的依据**是"`Close()` 拆定时器"这一读码推理（`machine.go:169-174`＋`models.go:304/:321`），**没有**时间测量；如果哪天 `handOffModel` 在 return 前做够 10s 的事，这一格会变。
6. **§2 的 12 枚重名去重后＝50**，但"每枚名字分布在几枚行上"我只在表里逐名标了行号，**没有**单独给"每枚名字出现的次数分布尺输出"（`uniq -c` 我跑了 `uniq -d`，见 2.1 那句）。

