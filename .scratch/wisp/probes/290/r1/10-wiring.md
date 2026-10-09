# 票 290 / 290-r1 — 10 接线（甲-1：只把门拧开，这一发不动状态机）

派单转述与票面裁定逐字读过；落点照「编排者裁定（2026-10-09 14:2x）」那一版执行，⛔ 未重新设计。
本腿起手锚＝`6c976aa6`（见 `00-anchor-and-baseline.md` §1）。
**锚漂具名**：本腿跑到一半，另一枚腿 `livewin-1` 落了 `fcc60648`（真机窗口逐形读数仪器与命令备好），
尺＝`git show --name-only --format= fcc60648 | grep -E "\.go$"` ⇒ **零枚 Go 文件**（只 1 枚
`.scratch/wisp/probes/orch/2026-10-09-live-window-runsheet.md`），⇒ 本腿的改前/改后 Go 面名册
仍然同一棵树可比；本腿产码笔落其下＝`fdee536b`。

## 1. 动了几枚文件（枚数＋行数，尺＝`git show --stat fdee536b`）

| 文件 | 行数 | 这一枚改了什么 |
|---|---|---|
| `cmd/wisp/resident_audio_windows.go` | +60 / -0 | `(*residentAudio).toggleMute()`——本发唯一那枚门侧写口调用者；文件头「Shape」清单多一条 mute（票 290 甲-1） |
| `cmd/wisp/resident_ball_windows.go` | +142 / -14 | `muteGestureFunc` 类型、`residentBall` 的 `muteMux`/`muteGate` 两枚字段、`attachMuteGate`/`currentMuteGate`/`muteGesture` 三枚方法、`Events` 里两支手势改走门、`ballGestureWhy` 与启动日志 `"gestures"` 那句按事实重推、文件头多一节「What changed since (ticket 290…)」 |
| `cmd/wisp/resident_windows.go` | +25 / -0 | 装配根在 `startResidentAudio` 之后 `rb.attachMuteGate(raudio.toggleMute)`（后置 setter）＋该跳边界说明（不 dispatch、不第二枚真相源、零新协程、托盘勾号不镜像） |
| `cmd/wisp/resident_mute_290_windows_test.go` | +287 / -0（新文件） | 六枚读数（见 §5） |

`git show --stat fdee536b` 逐字：`4 files changed, 514 insertions(+), 14 deletions(-)`。

## 2. 接线形状（每处改前后逐字）

### 2.1 `resident_ball_windows.go` 的 Events 两支（票面现量第 3 条那两枚空落地）

改前（逐字，原 `:281`／`:285`）：
```go
			OnMuteHotkey:    func() { recordBallGesture("mute-hotkey") },
			OnTrayMute:      func() { recordBallGesture("tray-mute") },
```
改后：
```go
			OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },
			OnTrayMute:      func() { rb.muteGesture("tray-mute") },
```
`rb` 在 `startResidentBall` 第一句就存在（`rb := &residentBall{cancelHosted: …}`），闭包捕获的是
指针，故 `ball.New` 时刻不需要知道门长什么样——这是后置 setter 能成立的前提。

### 2.2 门侧执行者（`resident_audio_windows.go`，紧接 `setAudioLevel` 之后）

```go
func (ra *residentAudio) toggleMute() (string, bool) {
	if ra == nil || ra.gate == nil || ra.mic == nil {
		return ra.unrunningClaim(), false
	}
	gate := ra.gate
	gate.SetMuted(!gate.Muted())          // ← AC#2 尺命中的那枚（HEAD:162）
	ra.started = gate.Open()             // 派生镜像，不是第二枚真相源
	switch { case gate.Muted(): … case gate.Open(): … default: 分类＋指引 }
}
```
要点：
- **写下去再读回来**：三条形句子全部取自己 `gate.Muted()` / `gate.Open()` 的现读，⛔ 不取"我刚才请求了什么"。
  设备拒绝交接那一形把 `mic.Err()` 的 `observe.Error.Class` 与 `mic.Stats().LastError` 打进去，
  与装配段（原 `:241-256`）同一把读法，⛔ 不谎报成功。
- **无门那两形**（`voice.enabled=false` ⇒ `ra.gate==nil`；或整枚 handle 为 nil）返回
  `ra.unrunningClaim(), false`——那句就是启动时本来要说的那句（`采集腿未构造：[voice] enabled=false（配置来源 %s）`
  或 `采集腿未装配：本进程没有任何麦克风路径`），⇒ 手势说出原因，而不是演一次静音。

### 2.3 球侧入口（`resident_ball_windows.go`）

`muteGesture(name string) string` 三条形：
1. 执行者尚未 attach（只可能在开机窗口里，`:217` 建窗、`:278` 建门之间）⇒ 逐字说明是哪一段窗口，
   ⛔ 不假装键生效；
2. `executed=false` ⇒ 把门侧给的原因原样说出（测试 `TestAC290NoGateSaysWhichShapeThisProcessIsIn`
   钉住"主机不许改写它收到的句子"）；
3. 正常 ⇒ Info 级日志＋同一句 console。
返回句子是给调用者/测试读的；`ball.Events` 要的是 `func()`，两支闭包按形状丢弃返回值（非放宽，是形状）。

### 2.4 装配根后置 setter（`resident_windows.go`）

```go
	raudio := startResidentAudio(rt, rb.setAudioLevel)
	…
	rb.attachMuteGate(raudio.toggleMute)
```
`attachMuteGate` 走 `rb.muteMux` ⇒ 与 ui-sta 线程的 `currentMuteGate()` 互为定序；
这一枚锁同时把 `assembleCapture` 写进 `residentAudio` 的一切（含 `gate`/`mic`/`started`）
排在任何一次手势读它们之前（boot 写 → unlock → ui 线程 lock → 读），⇒ 后置注入不是 data race。

## 3. 真相源那句（引文，⛔ 未倒过来）

`internal/audio/gate.go:20-21` 逐字：
```
//   - Muted (both paths, user intent wins): capture closed; Muted/Unmuted
//     events hook the mute hotkey path into the Muted state (ticket 07).
```
⇒ **门侧 `muted` 为主**。本发落地后 `residentBall` 里与"谁在静音"有关的字段只有 `muteGate`（一枚函数值，
不是状态），⛔ 零新增布尔；`ra.started` 是 `gate.Open()` 的镜像（派生，注释逐字写明），
⛔ 不参与"谁被静音"的判定。

**未做的两件事（裁定的射程，具名留档）**：
- ⛔ 不 dispatch `EvMuteKey`、不动 `internal/statemachine/table.go`、不做灰球＋斜杠——
  裁定逐字给了三条硬读数（表里无 `Sleeping→Muted` 边；`EvMuteKey` 产码发射者只有 `cmd/balldebug/main.go:603`；
  `WithGateEvents` 产码零调用者），那一格今天落不了，要人工批准 `A##`，不归写腿。
- ⛔ **托盘勾号 `internal/ball/ball_windows.go:128 trayMuted` 本发不镜像**（`SetTrayChecks` `:953` 产码调用者仍为 0）。
  它是 290-a1 点名的"第三枚真相源"位置：勾号今天仍是球侧一枚没人驱动的显示位，
  本腿若把它写成受门驱动就会新增一条跨面写口、并把 AC#2 那把尺的射程换成"镜像了几个字段"。
  ⇒ 具名残余：静音键真机按下去之后托盘菜单的勾号不会变，那一格归后续票（裁定里 乙-托盘半 那一支）。

## 4. 同批改掉的那句假话（票面裁定倒数第 2 条）

### 4.1 `ballGestureWhy`

改前逐字（原 `:408-410`）：
```go
const ballGestureWhy = "this process has no task pipeline and no microphone, so the gesture has no executor here; " +
	"the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel " +
	"gestures that reach the resident panel thread (ticket 33)"
```
改后逐字：
```go
const ballGestureWhy = "no executor was handed to this ball host for this gesture; the gestures with one are the " +
	"cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it " +
	"injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg " +
	"(ticket 290) - what is left is recorded by name, never invented"
```
推导逐名写在紧邻的注释块里（哪些有执行者、哪些没有、为什么）：
- **有**：cancel key（须装配根注入了审批门，票 246）、两支面板手势（须注入了面板 host，票 33）、
  两支静音（须本进程装配了采集腿，票 290）。三枚条件形是**故意的**：`requestPanelOpen` 在没有
  showPanel 时仍落回 `recordBallGesture`，写成一刀切的"这三类always work"就会变成新谎。
- **没有**：球点击与 summon 热键（本 host 没被交过"从球手势到任务"的任何口子；本 leg 确实装配了
  任务管线 `resident_task_source_windows.go`，但它读控制台，把球接进去不是本票的射程），
  托盘 pause-wake（要暂停的 KWS 住在没写的 `internal/speech`，且 `wake_word.enabled` 出厂 false）。
- **各自说自己的**：drag end（`recordBallDragEnd`，无 `Options.Store` ⇒ 位置不持久）、tray Exit
  （`recordTrayExit`，退出请求今天只有控制台信号这一条路）。
- 旧句两处假：`no microphone` 自票 247 起在本进程不成立（本进程有采集腿）；`no task pipeline`
  自票 246 AC#7（`startResidentTaskSource`，`resident_windows.go:260`）起也不成立 ⇒ 一并不再由本句代答。
- ⛔ 未写现在时计数句（注释里没有"全仓有 N 处调用"那种会被 d22scan 自己扫的形状；d22scan rc=0）。

### 4.2 启动日志里同形的第二句（本腿另找到的一处，票面没点名）

改前逐字（原 `:350-351`）：
```go
		"gestures", "recorded only except the cancel key: this leg has no task pipeline and no microphone, "+
			"and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)")
```
改后：同一条假话的两处 claim（`no task pipeline`/`no microphone`）与"只有 cancel key 有执行者"（自票 33 起就已偏低）
一并按事实重写，并逐字点明"这一行印在 attach 之前，因为球先建"。
⇒ 具名上报：**票面/派单只点了 `ballGestureWhy` 一处，同一枚假话在同文件那枚 slog record 里还有第二处**；
本腿按"不许留一句已被自己推翻的假话"把它同批改了，仍在同一枚文件内，未扩到别处。

### 4.3 文件头历史块

票 228 的原句（`no task pipeline and no microphone`）在文件头是**带标签的历史陈述**
（`:30` 逐字「this file's wording … is kept because it is still most of the truth」），
本腿⛔ 未改写历史句一字，只按其既有惯例追加一节「What changed since (ticket 290, form 甲-1)」，
内文逐字点明：门侧为主、本 host 不留副本、不动状态机、不做第二枚真相源。

## 5. 新增读数（`cmd/wisp/resident_mute_290_windows_test.go`，六枚，全绿）

| 尺名 | 问的是什么（能力，⛔ 不是词面） |
|---|---|
| `TestAC290BothMuteGesturesTurnTheGateAndBack` | 两枚手势名各自：静音位起步（源码零 start）→ 按一次 ⇒ `gate.Muted()==false` 且 `gate.Open()==true` 且源码 start 数＝开启数；再按一次 ⇒ 回到静音且 stop 数＝开启数；`reg.CountByName("audio-capture")==0` ⇒ 本发零新协程 |
| `TestAC290OutcomeIsReadOffTheGateNotOffTheRequest` | 设备拒绝交接（0x80070005 那枚 canonical `audio.DeviceError`）时句子必须打"设备未交接"＋错误分类＋指引文本，⛔ 不许报成功；`ra.started` 不许为真 |
| `TestAC290NoGateSaysWhichShapeThisProcessIsIn` | `voice.enabled=false` ⇒ 无门可拧，返回 `(原因, false)` 且原因含"未构造"；球 host 原样转述⛔ 不改写；nil handle 不 panic |
| `TestAC290GestureBeforeTheAttachSaysSo` | 未 attach 时那一发说的是"采集腿尚未装配"，⛔ 不谎称改了状态；nil ball 也安全 |
| `TestAC290MutedDefaultStillDecidesTheBoot` | 三枚默认值现读一字未动（`mic_muted_default=true`/`voice.enabled=true`/`wake_word.enabled=false`）＋boot 仍打"设备未打开"＋源码 boot 零 start；离开静音只能由那一发手势做到 |
| `TestAC290UnhostedGestureWordingStillSaysTrueThing` | `ballGestureWhy` 不再自称 `no microphone`，且点名票 290 的采集腿那一支 |

注入面：只走票 247 已有的那枚 `captureSource` 装配接缝（`assembleCapture` 的第 4 参数，
`resident_audio_247_windows_test.go:251` 同形），**门是真的门（`audio.NewHalfDuplexGate`）、装配是真的装配、
D38(e) 注册是真的注册**；⛔ 未用 mock 代替"真的门"来报完成，⛔ 未拿测试里的 `MicMutedDefault=false`
临时配置冒充用户路径（本批测试一律用出厂表 `writeAudioConfig(t, nil)`）。
设备真开那一发＝AC#3 第②形，本腿具名欠（见 `20-ac-readings.md`）。

## 6. 落点后的行锚（⛔ 用票面行号引的东西一律两态都给）

| 锚（票面/裁定给的旧号） | 旧号逐字 | 本腿落地后的新号（现读） |
|---|---|---|
| 球的静音手势那支 | `resident_ball_windows.go:281` | `:321` `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` |
| 托盘静音那支 | `resident_ball_windows.go:285`（票面现量写 `:284`，裁定写 `:285`） | `:325` `OnTrayMute:      func() { rb.muteGesture("tray-mute") },` |
| 常驻球初始态那句 | `resident_ball_windows.go:275` `Initial:  statemachine.StateSleeping,` | `:315`（一字未动，本腿没碰状态） |
| `ballGestureWhy` 那句 | `resident_ball_windows.go:408` | `:464` |
| 门的构造那枚 | `resident_audio_windows.go:213` | `:273`（一字未动） |
| 默认档那句 verdict | `resident_audio_windows.go:238` | `:298` |
| 本腿新增那枚调用者 | — | `resident_audio_windows.go:162` `gate.SetMuted(!gate.Muted())`（函数体自 `:156`） |
| 装配根两枚顺序 | `resident_windows.go:217`／`:278` | `:217`／`:278`（未漂；新 setter 在 `:303`） |

⚠ **具名上报（票面/裁定自己的一处枚数不一致）**：静音托盘那支，票面裁定段写 `:285`、
现量段与派单写 `:284`——起手锚 `6c976aa6` 上现读逐字＝`:285`（`:284` 是 `OnTrayPanel`），以原文为准。

## 7. 本腿造成的一处漂移与修法（⛔ 不是放宽断言）

改后整包第一发把票 255 的
`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 翻红（改前两发都绿）⇒ **本腿的账**：
`cmd/wisp/config_readers_255.go` 的 roster 用 `file.go:LINE [token]` 形钉住读点，
本腿在那两枚文件里的插入把行号推走了。修法＝**只更新三枚 cite 的行号**（`:196→:256`、`:197→:257`×4 row、
`:276→:316`×2 处），token 与 verdict 文字一字未动，⛔ 未动断言、⛔ 未把本票措辞塞进 255 的句子
（那句"[audio] 只在 boot 读一次、热重载要重启"在票 290 之后语义仍然成立——用户那一跳拧的是门、不是配置）。
先例＝`ada5c563`（247-r1 因自己的读点重裁过同一张 roster）。
笔＝`69d8c9ef`（`1 file changed, 7 insertions(+), 7 deletions(-)`）。

