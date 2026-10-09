# 290-v1 · 00 起手锚＋静态尺（非实现者验收腿，⛔ 零 AC 框翻动、零产码改动）

腿：`290-v1`（裁决者，≠ `290-r1` 实现者）。时钟＝`date` stdout，不手打。

## 0.1 锚与票面

- HEAD 起手：`5cff604e`（`票 290 改短名`）；被裁三笔＝`fdee536b`（产码）／`69d8c9ef`（cite 修复）／`75128ae5`（证据件）。
- 票面文件＝`.scratch/wisp/issues/290-gate-mounted-nobody-can-open-it.md`（25,768 字节，短名已生效）。
  尺＝`grep -c '^- \[ \]'` → **4**；`grep -c '^- \[x\]'` → **2**（rc=0）。逐枚现读：
  `:29 [x] AC#0`／`:30 [x] AC#1`／`:35 [ ] AC#2`／`:36 [ ] AC#3`／`:37 [ ] AC#4`／`:38 [ ] AC#5`。
  ⇒ **实现者未翻任何框**（票面 Progress log `:75` 那句"⛔ 零 AC 框自勾"与盘面一致）。
- 工作树脏件（别人的在飞改动，本腿一枚未碰）：`git status --porcelain` 行数＝**759**。

## 0.2 AC#2 那把尺（非测试调用者枚数）——本腿独立现跑

```
尺：grep -rn "SetMuted(\|SetSpeaking(" --include=*.go internal/ cmd/ | grep -v _test.go
命中 4 行（rc=0）：
  internal/audio/gate.go:97   // source is NOT started; SetMuted(false)/SetSpeaking(false) will start it   ← 注释
  internal/audio/gate.go:130  func (g *HalfDuplexGate) SetSpeaking(speaking bool) {                        ← 定义
  internal/audio/gate.go:168  func (g *HalfDuplexGate) SetMuted(muted bool) {                              ← 定义
  cmd/wisp/resident_audio_windows.go:162  gate.SetMuted(!gate.Muted())                                     ← ★唯一调用者
```
⇒ 起手锚（`fdee536b` 之前）＝**0 枚调用者**；落地后＝**1 枚**，且那枚在生产链路上、不是又一枚调试用 `cmd`。
票面 `AC#2` 写的判据是"由 2 枚变 ≥3 枚"（把注释行与两枚定义行都算进去了），实现者已具名报回这条枚数出入（`75128ae5` commit 消息末段"①AC#2 那条合尺在起手锚上现跑＝3 枚而非 2 枚"）——**本腿复核：出入成立，票面那句"2 枚"确实是把 `SetMuted(` 单尺与合尺混了**。真正有牙齿的是"调用者 0→1"，本腿认可这一把。

## 0.3 `attachMuteGate` 的生产调用者（尺＝票面指定的那条）

```
尺：grep -rn "attachMuteGate" --include=*.go cmd/ internal/   （rc=0，6 行）
  cmd/wisp/resident_ball_windows.go:469   ← 注释
  cmd/wisp/resident_ball_windows.go:472   ← 定义 func (rb *residentBall) attachMuteGate(fn muteGestureFunc)
  cmd/wisp/resident_windows.go:303        ← ★生产调用者 1 枚
  cmd/wisp/resident_mute_290_windows_test.go:93   /  :206 /  :235   ← 三枚测试
非测试行＝3（注释＋定义＋那一枚调用）；**真正的调用者＝1 枚**。
nilBall 那一处＝cmd/wisp/resident_mute_290_windows_test.go:235（派单写 `:234-236`，现读 `:234` 是 `var nilBall *residentBall`、
`:235` 是 `nilBall.attachMuteGate(...)`、`:236` 是 `nilBall.muteGesture("tray-mute")` ⇒ 派单那个区间准确）。
```

## 0.4 N1 的接线逐跳（每跳都是本腿现读，⛔ 未引用实现者的话）

| 跳 | 位置（现跑 `grep -n` 所得） | 逐字 |
|---|---|---|
| 1 热键注册 | `internal/ball/hotkey_windows.go:37` | `hkMute   = 2` |
| 1 热键注册 | `internal/ball/hotkey_windows.go:49` | `{hkMute, "mute"},` |
| 1 默认绑定 | `internal/ball/hotkey_windows.go:75` | `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}` |
| 1 注册动作 | `internal/ball/ball_windows.go:255` | `b.hotkeyReport = registerAll(b.hwnd, b.opts.Hotkeys)` |
| 1 球侧配置入口 | `cmd/wisp/resident_ball_windows.go:316` | `Hotkeys:  cfg,` |
| 2 Wndproc 分派 | `internal/ball/ball_windows.go:656` | `case wmHotkey:` |
| 2 mute 分支 | `internal/ball/ball_windows.go:660` | `case hkMute:` |
| 2 发球侧事件 | `internal/ball/ball_windows.go:661` | `b.fire(b.opts.Events.OnMuteHotkey)` |
| 3 球侧钩子 | `cmd/wisp/resident_ball_windows.go:321` | `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` |
| 3′ 托盘同支 | `internal/ball/ball_windows.go:678`＋`:679` | `case menuMute:` ＋ `b.fire(b.opts.Events.OnTrayMute)` |
| 3′ 托盘钩子 | `cmd/wisp/resident_ball_windows.go:325` | `OnTrayMute:      func() { rb.muteGesture("tray-mute") },` |
| 4 钩子本体 | `cmd/wisp/resident_ball_windows.go:515` | `func (rb *residentBall) muteGesture(name string) string {` |
| 4 取执行者 | `cmd/wisp/resident_ball_windows.go:516` | `fn := rb.currentMuteGate()` |
| 4 拧门 | `cmd/wisp/resident_ball_windows.go:527` | `outcome, executed := fn()` |
| 5 挂载动作 | `cmd/wisp/resident_windows.go:303` | `rb.attachMuteGate(raudio.toggleMute)` |
| 6 门侧执行者 | `cmd/wisp/resident_audio_windows.go:156` | `func (ra *residentAudio) toggleMute() (string, bool) {` |
| 6 真拧那一下 | `cmd/wisp/resident_audio_windows.go:162` | `gate.SetMuted(!gate.Muted())` |
| 7 门是生产装配的那枚 | `cmd/wisp/resident_audio_windows.go:273`＋`:277` | `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))` ＋ `ra.mic, ra.gate, ra.cancel = mic, gate, cancel` |

⇒ **每一跳都在生产装配里，没有一跳"只有测试里有"**。
⚠ 派单纠错一处（以盘上原文为准）：派单写 `case hkMuted:`——**盘上逐字是 `case hkMute:`（`internal/ball/ball_windows.go:660`）**，
票面 `:15` 那节本来就写 `case hkMute:`；派单还说行号"别信我写的"，本腿现跑取到 656/660/669/678，与票面 `:15` 一致。
派单另一处纠错：派单说 `OnMuteHotkey` 在球侧，现读 `Events` 名册 10 枚字段（`:318-329`），`OnMuteHotkey`/`OnTrayMute` 各 1 枚。

## 0.5 N5 的静态部分：托盘那枚勾

- 定义：`internal/ball/ball_windows.go:128` `trayMuted     bool`（注释同段是 display flag）。
- 被读进菜单的那一行：`internal/ball/ball_windows.go:674` `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)`。
- setter：`internal/ball/ball_windows.go:953` `func (b *Ball) SetTrayChecks(muted, pausedWake bool) {` → `:955` `b.trayMuted = muted`。
- **生产调用者尺**：
  ```
  grep -rn "SetTrayChecks" --include=*.go . | grep -v "^\./\.scratch" | grep -v "_test.go"   → rc=0，命中 2 行：
    ./internal/ball/ball_windows.go:952  ← 注释
    ./internal/ball/ball_windows.go:953  ← 定义
  含测试的全仓尺（internal/ cmd/）同样只有那 2 行 ⇒ **SetTrayChecks 的调用者＝0 枚（测试里也 0 枚）**。
  ```
⇒ `trayMuted` 今天**永远停在零值 `false`**，`showMenu` 拿到的静音勾恒为"未勾"；而 `menuMute` 那支（`:678-679`）自 `fdee536b` 起**真会拧门**（同一支走 `muteGesture`→`toggleMute`→`SetMuted`，见 §0.4 第 3′ 跳）。裁决正文见 `20-questions-n5-n6.md` §N5。
- 装配根自己把这件事写在了注释里（`cmd/wisp/resident_windows.go:300-302` 逐字）：
  `and no mirror of the gate into the tray's own checkmark - the` / `ball-side trayMuted display flag stays undriven here rather than becoming a` / `second authority nobody reconciles.`
  ⇒ **不是漏了，是具名不做的**；实现者自报"`SetTrayChecks` 产码调用者仍 0，归后续票"与本腿现跑一致。

## 0.6 N6 的两处新句（逐字）＋真值现跑

**(a) `cmd/wisp/resident_ball_windows.go:464-467`（`ballGestureWhy`）** 逐字：
`"no executor was handed to this ball host for this gesture; the gestures with one are the " +`
`"cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it " +`
`"injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg " +`
`"(ticket 290) - what is left is recorded by name, never invented"`

逐支核真（本腿现读，非引用）：
- cancel＝有执行者：`:322` `OnCancelHotkey:  func() { recordCancelHotkey(onCancelEsc) },`，`onCancelEsc` 来自装配根 `cmd/wisp/resident_windows.go:217` `rb := startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, withPanelHost(...))` ⇒ "when the assembly root injected an approval gate" **限定形属实**。
- panel＝有执行者：`:323` `OnPanelHotkey:   hs.requestPanelOpen("panel-hotkey"),`、`:324` `OnTrayPanel:     hs.requestPanelOpen("tray-open-panel"),`。
- mute＝有执行者（条件＝那枚腿装配了采集腿）：`:321`/`:325`，挂载在 `cmd/wisp/resident_windows.go:303`，条件性由 `toggleMute` 的 `ra == nil || ra.gate == nil || ra.mic == nil`（`resident_audio_windows.go:157`）兜住。
- 仍无执行者的三支＝`:319` click／`:320` summon-hotkey／`:326` tray-pause-wake ⇒ 与常量注释里点名的"the ball click and the summon key"／"the tray's pause-wake item"（`:451`、`:456`）一致。
⇒ **今天这句话是实话**，且它**没有**把任何还没做的事写成做了（电平不在句子里、托盘勾不在句子里、`Muted` 态/灰球不在句子里）。

**(b) `cmd/wisp/resident_ball_windows.go:390-393`（启动日志 `"gestures"` 串）** 逐字：
`"gestures", "recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) "+`
`"and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is "+`
`"assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four "+`
`"veto channels stay reduced here to the one the assembly root injected")`
- 时序自证为真：这行在 `startResidentBall` 内（`:387`），装配根的 attach 在 `cmd/wisp/resident_windows.go:303`，而建球在 `:217` ⇒ **"this line prints before the attach" 逐字属实**。
- 那句"turn this process's capture gate **once that leg is assembled**"是条件式，不是完成式 ⇒ 未声称已挂载。
⇒ **今天这句话是实话**。
- 附带核：`toggleMute` 成功那支的句子（`resident_audio_windows.go:174`）自己把越界面摘干净了，逐字
  `已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）` ⇒ **没有**把"电平已经在球上动"当成已验事实。

## 0.7 N8 逐笔禁区尺（⛔ 未用区间 diff）

```
git diff-tree --no-commit-id -r --name-only --no-renames <sha> | grep -cE "<十枚禁区＋frontend/design/testdata/golden/PLAN/specs>"
  fdee536b  hits=0  rc=1
  69d8c9ef  hits=0  rc=1
  75128ae5  hits=0  rc=1
```
- 三笔的文件名册全集：`fdee536b`＝4 枚全在 `cmd/wisp/`（`resident_audio_windows.go`/`resident_ball_windows.go`/`resident_mute_290_windows_test.go`/`resident_windows.go`）；
  `69d8c9ef`＝1 枚 `cmd/wisp/config_readers_255.go`；`75128ae5`＝18 枚全在 `.scratch/wisp/**`（含票面那一枚旧长名——见 §0.8）。⇒ **零枚落在禁区**。
- `git status --porcelain -- internal cmd` 现跑＝**空**（本腿起手；工作树里 `internal/`/`cmd/` 无未提交改动）。

## 0.8 与票面/派单不符项（本腿现跑，逐条具名）

1. **派单错**：`case hkMuted:` → 盘上 `case hkMute:`（`internal/ball/ball_windows.go:660`）。票面原文本来就对。
2. **`75128ae5` 的文件名册里有旧长名** `.scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md`——那是证据件提交**当时**的票名，`5cff604e` 才 `git mv` 成短名。工作树里旧名显示为 `D `（已 staged 的重命名左半边）。⇒ **不是实现者动了票面禁区**：AC#4 的禁区名册列的是 `PLAN.md`/`docs/specs/**`/`frontend/**`/`design/**` 等，票面本身不在禁区里，且只追加了 Progress log（现读 `:75` 一枚新行，`:73/:74` 是 a1 与编排者的旧行）。
3. **派单 vs 实现者自报的"三枚/七处"**：`69d8c9ef` 的 commit 消息写"三枚行号重指"，派单写"那七处重指"。现跑 `git show` ⇒ **改动行＝7 行、涉及的重指＝7 处 cite、不同的目标行号＝3 枚**（`:316`/`:256`/`:257`）。两种说法各自成立，本腿按"7 处逐处验"作答（见 `10-questions-n2-n4.md` §N4）。
4. **票面 `AC#5` 的命令名册比派单宽**：票面 `:38` 逐字要求 `go test ./cmd/wisp/ ./internal/audio/ -count=1`（**两枚包**），派单只写 `go test ./cmd/wisp/`。本腿按票面原文跑两枚包（以盘上原文为准）。
5. **实现者自报的三数口径**：`418/5/3` 与 `424/5/3` 用的是**含缩进子测试**的 `--- PASS` 计数（本腿在其 raw 件上现跑：`raw-pre-1`＝418/5/3、`raw-post-final-1`＝424/5/3、`pkgResultLines=2`），**只数顶层**则是 306/418→312/424 的另一套口径。⚠ 口径要在票面上写死，否则下一腿会拿顶层数当差集基准。这条不是缺陷，是"读数可复现但口径未声明"。
