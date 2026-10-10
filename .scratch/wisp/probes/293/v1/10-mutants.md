# 票 293 · `293-v1` · 10 突变台件（本腿自造三枚，⛔ 复述实现者的）

锚＝`b3593507`；时刻＝`2026-10-10 11:05–11:2x +08`。
突变面＝`cmd/wisp/resident_ball_windows.go`（`mirrorTrayMute` 的 `push` 那行 / `muteGesture` 的挂载那一行）。
定向尺（三枚共用，逐字）＝
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v -run 'TestAC293'`
原始 stdout＝`logs/targeted-head.txt`／`logs/mutantA.txt`／`logs/mutantB.txt`／`logs/mutantC.txt`（⛔ 0 字节件）。

## 0. 正控（先证明"没突变时全绿"，⛔ 缺一枚突变读数不作凭据）

`HEAD` 未改动那一发：`rc=0`、`--- PASS` **7** 枚、`--- FAIL` **0** 枚（7 枚逐名全在 `logs/targeted-head.txt`）。

## 1. 突变体 ›1＝镜像值钉成恒假

改：`push(muted, false)` → `_ = muted` ＋ `push(false, false)`。
⚠ **第一发我踩到同一枚坑**：只写 `push(false, false)` 不加 `_ = muted` ⇒
`cmd\wisp\resident_ball_windows.go:554:2: declared and not used: muted` ＋ `FAIL ... [build failed]`、
`rc=1` 且 `--- PASS`/`--- FAIL` 计数**都是 0** ⇒ **那是"没跑成"不是"跑绿"**（`logs/mutantA.txt` 头部逐字留有这一发的形状，
本腿具名作废这一发、按可编译形重跑）。

`rc=1`；**红 4 枚／绿 3 枚**。红句逐字（含文件行号，尺＝
`grep -E 'resident_tray_mute_293_windows_test\.go:[0-9]+:' logs/mutantA.txt`）：

```
--- FAIL: TestAC293TrayItemCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:126: tray item, after muting again: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293MuteHotkeyCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:148: mute hot key, after muting: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293BootProjectionShowsTheFactoryMutePosture
resident_tray_mute_293_windows_test.go:167: boot, factory posture: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth
resident_tray_mute_293_windows_test.go:285: flip: the tray's mute checkmark is false while the gate reports true - the checkmark must mirror gate.Muted()
```

恒假**杀不到**的三枚（本腿自己数出来的，不是抄它的）：
`RefusedDevice`（那一发门最终报 `false`，恒假与真值撞同 ⇒ 这一枚对恒假不敏感）、`NoGate`／`NilSafe`（本就不该写）。

## 2. 突变体 ›2＝摘掉"手势执行完那一跳"的挂载

改：删 `muteGesture` 里 `rb.mirrorTrayMute()`（`:617`），留注释占位。
`rc=1`；**红 4 枚／绿 3 枚**：

```
--- FAIL: TestAC293TrayItemCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:118: the tray checkmark was never written - the projection did not fire
--- FAIL: TestAC293MuteHotkeyCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:143: the tray checkmark was never written - the projection did not fire
--- FAIL: TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome
resident_tray_mute_293_windows_test.go:193: the tray checkmark was never written - the projection did not fire
--- FAIL: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth
resident_tray_mute_293_windows_test.go:285: the tray checkmark was never written - the projection did not fire
```

摘跳**杀不到** `BootProjection`（boot 那一跳独立存在于 `resident_windows.go:367`，不经手势）⇒ 起点③由 boot 那一枚用例单独钉住，
这一枚恰好证明"三种起点不是靠一枚顺带覆盖"。

## 3. 突变体 ›3＝本腿加造的第三形：假造一枚"永远说静音"的真相源（恒真）

派单只要求两枚；本腿多加一枚，因为它能把**禁区①**那形（另存一枚"我以为的静音"）单独问一遍。
改：`push(muted, false)` → `_ = muted` ＋ `push(true, false)`。
`rc=1`；**红 4 枚／绿 3 枚**：

```
--- FAIL: TestAC293TrayItemCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:118: tray item, after unmuting: the tray's mute checkmark is true while the gate reports false - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293MuteHotkeyCheckmarkFollowsTheGate
resident_tray_mute_293_windows_test.go:143: mute hot key, after unmuting: the tray's mute checkmark is true while the gate reports false - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome
resident_tray_mute_293_windows_test.go:193: refused device open: the tray's mute checkmark is true while the gate reports false - the checkmark must mirror gate.Muted()
--- FAIL: TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth
resident_tray_mute_293_windows_test.go:285: flip: the tray's mute checkmark is true while the gate reports false - the checkmark must mirror gate.Muted()
```

⇒ **恒真与恒假的红名册不重合**（恒假含 `Boot`、恒真含 `RefusedDevice`）⇒ 断言在**两个方向**上都对门态敏感，
按 `A710` 那把尺（"判据换成反形它也绿＝不敏感"）**不成立**：这套用例对"勾是不是抄的门"是有牙的。
禁区①那形（`mutedAtBoot` 出厂恒 `true`）在形状上＝恒真 ⇒ 被 ›3 这一族钉住。

## 4. 覆盖矩阵（本腿现量，行＝用例，列＝突变体）

| 用例 | ›1 恒假 | ›2 摘手势跳 | ›3 恒真 |
|---|---|---|---|
| TrayItem | 红 | 红 | 红 |
| MuteHotkey | 红 | 红 | 红 |
| BootProjection | **红** | 绿（设计如此） | 绿（恒真与出厂态撞同） |
| RefusedDevice | 绿（恒假与门态撞同） | **红** | **红** |
| NoGateWritesNoCheckmark | 绿 | 绿 | 绿 |
| NilSafe | 绿 | 绿 | 绿 |
| UndrivenTrayFlag | 红 | 红 | 红 |

`NoGate`／`NilSafe` 三枚全绿**是设计如此**（它们断的是"不该写的时候一枚都不写"，本就该对镜像**值**不敏感）；
⚠ 但本腿⛔ 没造出能杀它俩的突变体（形如"`!ownsGate` 也照写"）⇒ 那两枚的牙＝**只由代码读形支撑、未由突变证明**，具名记欠账。

## 5. 还原证明（⛔ `git checkout`／`stash`，用本腿仓外备份写回）

- 备份＝`$TEMP/rbw.293v1.bak`（仓外），动手前 `sha256sum` 两行同值：
  `e67b7031afc79a242387cb838de52de80f0319c35b96121b0f99a34a03778d00`（＝`cmd/wisp/resident_ball_windows.go`＝备份）。
- 三枚**每一枚都在同一条命令里 `cp` 写回**（防中途炸），写回后现量：
  `sha256sum cmd/wisp/resident_ball_windows.go` ⇒ **`e67b7031…78d00`，与动手前逐字同值**（三枚各一次，见三份 logs 末行）。
- `git diff --stat -- cmd/wisp` ⇒ **空**；`git status --porcelain cmd/wisp | wc -l` ⇒ **0**。
- 附：`sha256sum` 同一枚文件的 `git show HEAD:` 副本 ⇒ 亦＝`e67b7031…78d00`（树上工作文件与 blob 逐字节同，LF）。
- ⚠ **本腿⛔ 复现不了实现者件的还原值**：`.scratch/wisp/probes/293/r1/10-ac2-behaviour-and-mutation.md:137`
  写 `sha_before=sha_after=74a52c42871f830bc9377a462ee4a61030e7004fa8b27daa906d7534e70b9c32`，
  而本腿在同一枚文件、同一枚 HEAD 上现量＝`e67b7031…`。树本身可证没被留脏（`git diff` 空＝更硬的尺），
  但**这一枚 sha 值今天不可复算**＝证据件缺陷，具名记（口径与 `A799` 已抓到的那两枚同类：号写错、枚数写错，这里是**值写错**）。
