# 票 293 · 写腿 `293-r1` · `10` `AC#2` 交付：勾镜像门态（ⓐ 行为断言 + ⓑ 反形正控）

落件时刻＝`2026-10-10 10:0x +08`（本腿 `date` 读数：起手 `09:43:32 +0800`）
HEAD 起点＝`6ef14788`；本腿产码＋测试笔＝`26289b9e`；本文件随该笔之后落。

---

## 1. 落点形状（甲形，票面已裁，未改形）

三处改动全在 `cmd/wisp`，**`internal/ball` 零改动**（尺＝`git diff 6ef14788..HEAD --stat -- internal/ball`＝**空**，
见 `30-ac3-scope.md`）；⛔ 新包级依赖边 0（`cmd/wisp → internal/ball` 与 `cmd/wisp → internal/audio` 两条本来就有，
本发没引入任何新 import）；新协程 0（没写一个 `go func(`，见 `20-ac4-gates.md` 的 d22scan 读数）。

| 文件 | 加了什么 | 角色 |
|---|---|---|
| `cmd/wisp/resident_audio_windows.go` | `func (ra *residentAudio) trayMuteState() (bool, bool)` | **读半**：`gate.Muted()` ＋「本进程有没有门」 |
| `cmd/wisp/resident_ball_windows.go` | `residentBall` 两枚字段（`trayMuteRead` / `trayCheckPush`，共用现成 `muteMux`）＋ `attachTrayMuteProjection(...)` ＋ `mirrorTrayMute() bool` ＋ `muteGesture` 末尾一跳 | **投影执行者**：读门 → 推给托盘显示面 |
| `cmd/wisp/resident_windows.go` | 挂载点**下一跳**（`if rb.b != nil { attachTrayMuteProjection(raudio.trayMuteState, rb.b.SetTrayChecks); rb.mirrorTrayMute() }`）＋ 同批改写那段注释 | **装配根**：boot 起点那一次投影 |

选「紧邻挂载点加一跳」而不是「包一层注入的执行者」的理由：`rb.attachMuteGate(raudio.toggleMute)` 那一行
**票面钉死不许动**（派单 §3「挂载点自己⛔ 不要动」），而包装执行者必须改写实参 ⇒ 只有独立一跳同时满足
「不动挂载行」与「三种起点都覆盖」。

## 2. 为什么这一形对三种起点都成立（逐起点调用链）

菜单是**每次右键现建**的（`internal/ball/ball_windows.go:674` `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)`），
`b.trayMuted` 的唯一写点在 `SetTrayChecks` 体内（`:955`）⇒ 正确形＝**门态每变一次就推一次投影**，
下一次右键读到的就是「截至那一刻」的门态；不需要、也没有第二处去「对账」。

- **① 点托盘那一项**：`ball_windows.go:678 case menuMute → b.fire(OnTrayMute)`
  → `resident_ball_windows.go:336 OnTrayMute: func() { rb.muteGesture("tray-mute") }`（＝`6ef14788` 的 `:325`，本腿在它上面加了 11 行字段注释 ⇒ 号已挪，认文本）
  → `muteGesture` 里 `fn()` ＝ `raudio.toggleMute` → `gate.SetMuted(!gate.Muted())`（`resident_audio_windows.go:162`，改前改后同号）
  → **返回后同一函数体内 `rb.mirrorTrayMute()`** → `trayMuteState()=gate.Muted()` → `push(muted,false)=rb.b.SetTrayChecks`。
- **② 全局静音热键**：`ball_windows.go:660-661 wmHotkey/hkMute → OnMuteHotkey`
  → `resident_ball_windows.go:332 rb.muteGesture("mute-hotkey")`（＝`6ef14788` 的 `:321`） → **与①同一条尾巴**（两枚起点共用 `muteGesture`
  的唯一 `executed==true` 出口，所以镜像只写一处就对两枚起点都成立；票面 `现量` 第 5 条「同罪的另一面」正是这条共用）。
- **③ boot 出厂静音位**：门在 `assembleCapture` 就以出厂默认挂载
  （`resident_audio_windows.go:291` `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`
  ＝`6ef14788` 的 `:273`，本腿在它前面插了 18 行的 `trayMuteState` ⇒ 号已挪，认文本；默认 `true`），
  这条路**不经过任何手势** ⇒ 靠 `resident_windows.go` 挂载点下面那一次 `rb.mirrorTrayMute()` 投影。
  ⚠ 这一枚也是唯一一枚「镜像不可能被①②顺带覆盖」的起点 ⇒ 突变体 B 的证据里它单独绿、单独被 boot 那一跳钉住。

三条禁区核对：① 勾读回 `gate.Muted()`，**没接 `ra.mutedAtBoot`**（该字段本腿一字未动，尺见 `30-ac3-scope.md`）；
② 第二实参**逐字 `push(muted, false)`**，理由写在 `mirrorTrayMute` 的注释里（静音不打勾＝谎；暂停唤醒今天无真相源可喂：
`internal/speech` 只有 `doc.go`、出厂 `wake_word.enabled=false`、那一项今天仍只 `recordBallGesture`）；
③ `rb.b` 的 nil 护栏**在取方法值之前**（`if rb.b != nil` 包住 `rb.b.SetTrayChecks`，注释写明取值的表达式自己会 panic）。

## 3. ⓐ 行为断言（七枚用例，同包夹具复用票 290 的 `muteRig`）

件＝`cmd/wisp/resident_tray_mute_293_windows_test.go`。尺＝`grep -c -- '--- PASS'`（含子测试，本仓沿用那把）。

定向读数（本腿 10:0x 现跑，命令逐字
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v -run TestAC293`）＝
**绿 7 枚 / 红 0 枚**，`ok github.com/CarlosShao/wisp/cmd/wisp 0.116s`，逐字：

```
--- PASS: TestAC293TrayItemCheckmarkFollowsTheGate (0.01s)
--- PASS: TestAC293MuteHotkeyCheckmarkFollowsTheGate (0.01s)
--- PASS: TestAC293BootProjectionShowsTheFactoryMutePosture (0.01s)
--- PASS: TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome (0.01s)
--- PASS: TestAC293NoGateWritesNoCheckmark (0.01s)
--- PASS: TestAC293ProjectionIsNilSafeWithoutABallWindow (0.01s)
--- PASS: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth (0.01s)
```

每枚断言的都是**同一个同值式**：送进显示面（＝生产里 `Ball.SetTrayChecks`）的那个 bool **必须等于那一刻的
`ra.gate.Muted()`**（助手 `assertMirrorsGate` 现读 `ra.gate.Muted()` 比对，⛔ 不与录制器自比），
外加每一次写都断言第二实参＝`false`。

| 用例 | 钉的是哪一枚起点 / 哪一条禁区 |
|---|---|
| `TestAC293TrayItemCheckmarkFollowsTheGate` | 起点①，开门与关门两向，且断言「一次被接受的翻转＝恰好一次写」（1→2 枚） |
| `TestAC293MuteHotkeyCheckmarkFollowsTheGate` | 起点②，两向 |
| `TestAC293BootProjectionShowsTheFactoryMutePosture` | 起点③；另断言出厂默认仍为 `true`（⛔ 不改默认值）＋ **投影没碰设备**（`src.calls()=0,0`，镜像只读不动门） |
| `TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome` | 设备拒绝交接时勾＝门自己的答案，⛔ 抄「请求的乐观结果」也⛔ 抄那句中文结果文本 |
| `TestAC293NoGateWritesNoCheckmark` | 无门（`voice.enabled=false`）与 nil 手把：`trayMuteState()`＝`(false,false)` ⇒ **一枚都不写**，不凭空造勾（禁区①的另一面） |
| `TestAC293ProjectionIsNilSafeWithoutABallWindow` | 禁区③：没挂投影／只挂读半／`rb.b==nil`／nil 宿主都不 panic 且不写；末尾断言门没被「问」动 |
| `TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth` | 三枚翻转＝三枚写，且 `executed==false` 的手势**一枚都不写**（勾永不跑在门前面） |

### 录制器那一面（诚实申报，⛔ 拿它冒充 AC#5）

`cmd/wisp` 今天**没有任何一枚测试造出真球窗口**（尺＝`grep -rln 'ball.New(' cmd/wisp/*_test.go`＝0 命中），
而 `Ball.SetTrayChecks` 要把任务投到 STA 线程（零值 `&ball.Ball{}` 没有那条线程），所以测试里 **push 半用录制器代替**，
读半仍是生产持有的那枚真 `HalfDuplexGate`。把录制器钉回真设面的是两件事：
- 件里的编译期断言 `var _ func(*ball.Ball, bool, bool) = (*ball.Ball).SetTrayChecks`（改名／变实参个数／变类型即不编译）；
- 生产挂载点写的是 `rb.b.SetTrayChecks` 本身（`resident_windows.go` 那一跳），不是任何适配层。
⇒ 球侧那三行（`:955` 写 / `:674` 读 / `tray_windows.go:88 appendItem(menuMute, "静音", muted)` → `:81-82 mfCheckd`）
**本腿一字未动**，勾最终**用户看得见没有**＝`AC#5`（机主在场的真机窗口），本腿不声称、也不翻框。

## 4. ⓑ 反形正控（两枚突变体，都是跟踪文件 ⇒ 附还原证明）

⛔ 本腿**没有**拿「`SetTrayChecks` 调用者枚数＝1」这类尺当凭据（那枚尺只在 `30-ac3-scope.md` 作现状描述用）。
两枚突变体各自：先 `cp` 到**仓外**（`/tmp/293_mut*_backup.go`）→ 改 → 跑定向用例 → 用仓外副本还原 → 比 `sha256sum`。

### 突变体 A＝镜像那行换成空操作（恒假，等价 `SetTrayChecks(false, …)`）

```
	push(muted, false)   →   _ = muted
                           push(false, false)
```
（第一发还试过 `push(false, false)` 不加 `_ = muted`：`muted` 变成已声明未使用 ⇒ `[build failed]`，
那一次不算数，已具名重跑；这条也顺带证明这枚变量确实被用在读回上。）

`mutantA_rc=1`；**红 4 枚 / 绿 3 枚**；红名与红句逐字（尺＝
`grep -E 'resident_tray_mute_293_windows_test\.go:[0-9]+:' logs/mutant-a-const-false.txt`）：

```
--- FAIL: TestAC293TrayItemCheckmarkFollowsTheGate (0.01s)
resident_tray_mute_293_windows_test.go:126: tray item, after muting again: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293MuteHotkeyCheckmarkFollowsTheGate (0.01s)
resident_tray_mute_293_windows_test.go:148: mute hot key, after muting: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293BootProjectionShowsTheFactoryMutePosture (0.01s)
resident_tray_mute_293_windows_test.go:167: boot, factory posture: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth (0.01s)
resident_tray_mute_293_windows_test.go:285: flip: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
```

★诚实申报**这枚突变体杀不到的那三枚**：`RefusedDevice`（那一刻 `gate.Muted()` 恰为 `false`，恒假与真值相等）、
`NoGate`（本就不该写）、`NilSafe`（不该写）⇒ 恒假形由上面 4 枚钉住，⛔ 全套 7 枚都能杀恒假。

### 突变体 B＝摘掉 `muteGesture` 里那一跳（只留 boot 投影）

`rb.mirrorTrayMute()` → 一行注释；`mutantB_rc=1`；**红 4 枚 / 绿 3 枚**：

```
--- FAIL: TestAC293TrayItemCheckmarkFollowsTheGate   (resident_tray_mute_293_windows_test.go:118: the tray checkmark was never written - the projection did not fire)
--- FAIL: TestAC293MuteHotkeyCheckmarkFollowsTheGate (…:143: the tray checkmark was never written - the projection did not fire)
--- FAIL: TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome (…:193: the tray checkmark was never written - the projection did not fire)
--- FAIL: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth (…:285: the tray checkmark was never written - the projection did not fire)
```
这枚**故意杀不到** `BootProjection`（boot 那一跳还在）⇒ 起点③不是靠手势那两枚顺带过的，
两枚突变体合起来把「boot 一跳」与「手势一跳」各自钉死；`NoGate` / `NilSafe` 仍绿是设计如此（无门不该写）。

### 还原证明（两枚各一发，读数相同）

```
sha_before=74a52c42871f830bc9377a462ee4a61030e7004fa8b27daa906d7534e70b9c32   (cmd/wisp/resident_ball_windows.go)
sha_after =74a52c42871f830bc9377a462ee4a61030e7004fa8b27daa906d7534e70b9c32   RESTORE_IDENTICAL
git status --porcelain -- cmd/wisp internal/ball  ⇒  空（两枚突变体之后各现跑一次，都干净）
```
⚠ 派单 §4 写的「`git diff --stat -- '*.go'` 为空」这一把在本腿**不可能为空**：本腿合法的产码改动就是 4 枚 `.go`
（`26289b9e`），那条尺的字面形是「零产码腿」的形状。本腿按**等价而更严**的形交：突变前后同一枚文件的 `sha256` 相同
＋ 跟踪树 `git status --porcelain` 复原为空 ⇒ 突变体在树上留下的净残留＝0 字节。原样读数与名册在
`30-ac3-scope.md`，突变原始 stdout 在 `logs/mutant-a-const-false.txt`／`logs/mutant-b-no-gesture-mirror.txt`。

## 5. 旧注释那一发（票面「编排者裁定」★点名要求的那件）

`cmd/wisp/resident_windows.go` 里那三行逐字（内容锚 `no mirror of the gate into the tray's own checkmark`）
**已不在树上**（尺＝`grep -n 'no mirror of the gate into the tray' cmd/wisp/resident_windows.go`＝**0 命中**，
见 `20-ac4-gates.md`），替换成主从关系的实话：`gate.Muted()` 是唯一真相源、`trayMuted` 是从它写出去的显示态、
没有第二座权威，并写明「自票 290 起托盘那一项真拧门，不再镜像就是隐私谎」。
挂载点那一行 `rb.attachMuteGate(raudio.toggleMute)` **逐字未动**。
`muteGesture` 的文档注释同批补了一段（说明投影是单向、且这里从不读回），⛔ 让下一位把「不镜像」当既有政策。
