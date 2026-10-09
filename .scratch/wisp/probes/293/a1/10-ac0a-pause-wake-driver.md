# 293-a1 · 10 · AC#0 ⓐ — 另一枚勾 `trayPauseWake` 今天有没有驱动者？

**锚点**＝`f019f45e`。尺＝`git grep -n ... HEAD`（对 HEAD，不对工作树）。⛔ 全部只读。

---

## 结论（一句）

**没有。`trayPauseWake` 与 `trayMuted` 同病：全仓（含测试）驱动者＝0 枚**——它唯一的写入点在 `SetTrayChecks` 的函数体里，而 `SetTrayChecks` 自己 0 枚调用者，所以那枚勾今天**恒为 `false`**、和静音那枚一样永不动。

⚠ 但**两枚勾的"错状态"性质不一样**，这一条是本节要交给编排者裁的实质（详见 §4）：`pause-wake` 的"不打勾"目前**说的是实话**（那一项点了什么也不会发生），`mute` 的"不打勾"**说的是谎**（票 290 之后点了真开门）。

---

## 1. 尺一形（与票面 `SetTrayChecks` 那把同形）

### 1.1 主尺：`git grep -n "trayPauseWake" HEAD -- internal cmd`

逐字 stdout（去掉票池噪音前已确认 internal/cmd 只含产码）：
```
HEAD:internal/ball/ball_windows.go:129:	trayPauseWake bool
HEAD:internal/ball/ball_windows.go:674:			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)
HEAD:internal/ball/ball_windows.go:956:		b.trayPauseWake = pausedWake
```
⇒ **3 命中，全是球侧内部**：1 枚字段声明（`:129`）、1 枚**读**（`:674` 送进菜单）、1 枚**写**（`:956`）。

### 1.2 含测试的全仓尺：`git grep -n "trayPauseWake" HEAD -- '*.go'`
```
HEAD:.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go:129:	trayPauseWake bool
HEAD:.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go:674:			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)
HEAD:.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go:922:		b.trayPauseWake = pausedWake
HEAD:internal/ball/ball_windows.go:129:	trayPauseWake bool
HEAD:internal/ball/ball_windows.go:674:			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)
HEAD:internal/ball/ball_windows.go:956:		b.trayPauseWake = pausedWake
```
⇒ 全仓 `*.go` 共 6 命中，**3 枚是探针副本**（`.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go`，历史快照，非产码，已在 `00` §4 具名排除），**3 枚＝1.1 那三枚**。
⇒ **测试文件命中＝0 枚**。没有一录用例读它、也没有一录用例写它。

### 1.3 那个唯一写入点的宿主：`git grep -n "SetTrayChecks" HEAD -- '*.go'`
```
HEAD:.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go:918:// SetTrayChecks updates the mute / pause-wake checkmarks.
HEAD:.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go:919:func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
HEAD:internal/ball/ball_windows.go:952:// SetTrayChecks updates the mute / pause-wake checkmarks.
HEAD:internal/ball/ball_windows.go:953:func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
```
⇒ 票面那把尺（`HEAD -- internal cmd`）我复跑＝只 2 行（注释 `:952` ＋定义 `:953`）；**放宽到全仓含测试仍＝同 2 行**（另 2 行是探针副本）。
⇒ **`SetTrayChecks` 调用者枚数＝0（改前 0、含测试仍 0）**。两枚勾共享这唯一一座桥，桥两头都没人走。

`SetTrayChecks` 函数体逐字（`internal/ball/ball_windows.go:952-958`，`sed -n` 现读）：
```go
// SetTrayChecks updates the mute / pause-wake checkmarks.
func (b *Ball) SetTrayChecks(muted, pausedWake bool) {
	b.sta.PostTask(func() {
		b.trayMuted = muted
		b.trayPauseWake = pausedWake
	})
}
```
⇒ 两枚 bool 在**同一个 `PostTask` 闭包**里一起被写——**形状上就是一次调用同时驱动两枚勾**；票面"同一枚 setter 的第二个实参"的说法在产码里得到逐字确认。

### 1.4 菜单侧的读（`internal/ball/tray_windows.go:72-89` 现读逐字）
```go
func showMenu(hwnd windows.HWND, muted, pausedWake bool) uint32 {
	...
	appendItem := func(id uintptr, label string, checked bool) {
		flags := uintptr(mfString)
		if checked {
			flags |= mfCheckd
		}
		pAppendMenuW.Call(menu, flags, id, unsafePtr(utf16(label)))
	}
	appendItem(menuOpenPanel, "打开面板", false)
	pAppendMenuW.Call(menu, mfSepart, 0, 0)
	appendItem(menuMute, "静音", muted)
	appendItem(menuPauseWake, "暂停唤醒", pausedWake)
```
⇒ **两枚勾走的是同一条 `appendItem(..., checked)` 通路**，`menuPauseWake` 的标签逐字＝`暂停唤醒`（票面没写这枚标签，本腿现量补上）。⇒ "恒 `false`" 的实际用户可见后果＝**"暂停唤醒"这一项也永远不带勾**。

---

## 2. 点击那一头（`ⓑ` 的对照面：勾没驱动，那"动手"驱动了吗？）

尺＝`git grep -n "OnTrayPauseWake\|OnTrayMute\|OnMuteHotkey" HEAD -- '*.go'`，产码命中（探针副本已滤）：

| 文件:行 | 逐字 | 性质 |
|---|---|---|
| `cmd/wisp/resident_ball_windows.go:321` | `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` | **有 executor**（票 290） |
| `cmd/wisp/resident_ball_windows.go:325` | `OnTrayMute:      func() { rb.muteGesture("tray-mute") },` | **有 executor**（票 290） |
| `cmd/wisp/resident_ball_windows.go:326` | `OnTrayPauseWake: func() { recordBallGesture("tray-pause-wake") },` | **无 executor**，落 `recordBallGesture` 空档 |
| `cmd/balldebug/main.go:201-202` | `OnTrayPauseWake: func() { fmt.Println("tray: pause-wake toggled (stub, ticket 41)")` | debug 腿的显式 stub |

⇒ `pause-wake` 那一头**也只有空落地**：`recordBallGesture`（`resident_ball_windows.go:545-548`）逐字只做 `slog.Warn` + `fmt.Printf`，**不动任何设备、不动任何状态**。

---

## 3. `pause-wake` 为什么是空落地——产码里已有一段**写死的理由**（逐字）

`cmd/wisp/resident_ball_windows.go:456-457`（现读）：
```go
//   - the tray's pause-wake item: the wake-word listener it would pause lives in
//     internal/speech, which is unwritten, and wake_word.enabled ships false.
```

本腿独立复核那两句事实（⛔ 没照抄注释）：
- `ls internal/speech` → 只有 **`doc.go` 一枚文件** ⇒ 该包**已声明、无实现**。✅ 注释那句 "unwritten" 与盘上一致。
- `internal/config/schema.go:197` → 逐字 `Enabled bool \`toml:"enabled" default:"false"\``（属 `WakeWord` 结构体，`schema.go:196` 起）⇒ 出厂 `wake_word.enabled` **确实**是 `false`。✅

（⚠ 顺带：`schema.go:197` 正是票 293「禁区」一节列的三枚不许动的默认值之一。本节只是**读**它，⛔ 未改一字。）

---

## 4. 交给编排者的判语（⛔ 本腿不挑支、不替他写决定）

把 §1–§3 合起来，两枚勾今天的**状态**与**错的状态**是两回事：

| | `trayMuted`（静音） | `trayPauseWake`（暂停唤醒） |
|---|---|---|
| 勾的驱动者 | **0 枚** | **0 枚** |
| 点击会不会真动手 | **会**（`:325` → `muteGesture` → `resident_audio_windows.go:162` 真拧门） | **不会**（`:326` → `recordBallGesture`，纯记账） |
| 界面上"不打勾"这句话 | **假的**（用户看得见：刚开了麦，勾还说没静音）⇒ 票 293 的靶 | **暂时真的**（确实没暂停任何东西 ⇒ 不打勾是**对的**） |

⇒ 本腿量到的**可裁事实**：`SetTrayChecks(muted, pausedWake)` 一次调用同时写两枚，⚠ 但**只有第一枚有真相源可镜像**（`gate.Muted()`）；**第二枚今天在整个仓里没有任何可镜像的真相源**——它该映的那个态（唤醒监听被暂停）在 `internal/speech` 里，而那枚包只有 `doc.go`。

⇒ 于是编排者要裁的那一格是：**落地腿该只镜像 `muted` 一枚、把 `pausedWake` 实参显式留 `false`？还是要求同一发里两枚都给？** 前者的凭据＝"给 `pausedWake` 造一枚状态"＝**禁区第 2 条禁止的'新造真相源'**；后者的凭据＝setter 的形状是两枚一起的、只喂一枚会让"暂停唤醒"这一项继续零信息（但零信息 ≠ 说谎）。⚠ **本腿不在这两者里挑，也不在票面里替他勾任何框**——票 293 `AC#2` 那句"不许新造第三枚真相源"已把这题的一半钉住，剩下那一半（显式留 `false` 要不要写成一条注释/断言）留给编排者。

---

## 5. 顺带量到的同族第三枚（供 `AC#0 ⓒ` 交叉引用）

尺＝`git grep -n "SetTrayTip" HEAD -- '*.go'`（已滤探针）：
```
HEAD:internal/ball/ball_windows.go:960:// SetTrayTip updates the tray tooltip (state surfacing).
HEAD:internal/ball/ball_windows.go:961:func (b *Ball) SetTrayTip(tip string) {
```
⇒ **托盘的第三面（tooltip）也是 0 枚调用者**，与两枚勾同形、同一个"有名字、没使用者"的病。⚠ 这一枚**不在票 293 的射程内**（本票只裁勾），但票 228 的现量段（`.scratch/wisp/issues/228-...:139`）早就把 `SetTrayChecks`／`SetTrayTip` 并列为"全仓零枚生产调用者"——**这条既有主张本腿复跑后成立**，具名记此以说明"不是本腿新发现"。
