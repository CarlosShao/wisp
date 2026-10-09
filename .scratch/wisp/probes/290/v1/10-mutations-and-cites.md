# 290-v1 · 10 变异台件（N2 六枚用例有没有牙／N3 恒真面／N4 七处重指是否只改了数）

⛔ 三发变异都**当场还原**，还原后与 HEAD 逐字节等值（§4 给 md5 与 blob 双尺）。
台件纪律：输出先落 `/tmp/290v1/*.md` 再读；退码取**被测那条命令**的 `$?`（`go test ... > file 2>&1; echo rc=$?`），⛔ 无 `cmd | head; echo $?` 形状。
计数尺逐字＝`grep -cE "^[[:space:]]*--- (PASS|FAIL|SKIP)"`（含缩进子测试，与实现者 `30-gates.md` §2 同口径）。

## 1. 基线（HEAD=`5cff604e`，未突变）

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -run 'TestAC290' -count=1 -v
rc_targeted_baseline=0    --- PASS x6（六枚 TestAC290* 全绿），ok cmd/wisp 0.102s
```

## 2. N2 · M1＝把"真拧门那一下"改成拧了等于没拧（保留 grep 命中）

突变体（`cmd/wisp/resident_audio_windows.go:162`）：
```
改前  gate.SetMuted(!gate.Muted())
改后  gate.SetMuted(gate.Muted())   // 290-v1 MUTATION-M1
```
⇒ **AC#2 那把尺在这个突变体下仍然命中**（现跑 `grep -rn "SetMuted(\|SetSpeaking(" --include=*.go internal/ cmd/ | grep -v _test.go | grep -c "resident_audio_windows.go:162"` → **1**，rc=0）。
读数：
```
rc_m1=1
--- FAIL: TestAC290BothMuteGesturesTurnTheGateAndBack
--- FAIL: TestAC290OutcomeIsReadOffTheGateNotOffTheRequest
--- PASS: TestAC290NoGateSaysWhichShapeThisProcessIsIn
--- PASS: TestAC290GestureBeforeTheAttachSaysSo
--- FAIL: TestAC290MutedDefaultStillDecidesTheBoot
--- PASS: TestAC290UnhostedGestureWordingStillSaysTrueThing
```
**必须红的三枚＝预测＝实测（三枚全中，无一枚多余）**，且红在各自的哪一条断言（实跑原文）：
```
resident_mute_290_windows_test.go:119: mute-hotkey: the gate still reports muted after the gesture: "已静音：采集已关闭，设备未打开（再按一次取消静音）"
resident_mute_290_windows_test.go:178: a failed open was reported without saying it failed: "已静音：…"
resident_mute_290_windows_test.go:265: the user's gesture cannot leave the boot posture: "已静音：…"
```
⇒ **这六枚有牙**：拧门那一行被拔掉时三枚立刻红，红的都是"门自己的读数"（`ra.gate.Muted()`/`ra.gate.Open()`）而不是刚写进去的布尔。
另两枚绿的是**该绿的两枚**——`NoGate…`／`GestureBeforeTheAttach…`／`UnhostedGestureWording…` 判的是"没有门／没挂载／文案"三个面，与拧不拧无关（这正说明它们不是凑数的，见 §3 的 M2）。

### ★票面 AC#2 那把尺单独不够（M1 就是反证）

票面 `AC#2`（`:35`）的判据**只有**"非测试调用者枚数 由 2 枚变 ≥3 枚，且新命中的那枚在生产链路上"。M1 证明：
**一枚 `gate.SetMuted(gate.Muted())` 也满足这把尺（命中位不变、枚数不变），但门一格都不拧。**
⇒ AC#2 那条尺是**必要不充分**；本票真正交付的是 `resident_mute_290_windows_test.go` 那三枚行为断言。
裁 `AC#2` 时本腿把这条记在判语里（见 `30-ac-verdicts.md`），⛔ 不据此判失败——票面上写的"完成判据"由测试补住了。

## 3. N2／N3 · M2＝反形正控：门真拧了，但 `executed` 谎报为 `false`

突变体（`resident_audio_windows.go:172` 与 `:174` 两枚 `true`→`false`；`:162` 保持 `SetMuted(!gate.Muted())` 原样）：
```
rc_m2=0    --- PASS x6  ← 六枚全绿
```
⇒ **实测发现一处覆盖缺口**：六枚用例里**没有一枚**在"门真的被拧了"这条happy path 上断言 `executed==true`。
`TestAC290NoGateSaysWhichShapeThisProcessIsIn`（`:196-199`）只在**没有门**那一支断言 `executed` 为假；
`muteGesture`（`resident_ball_windows.go:527-537`）拿到 `executed=false` 时只是改打 `slog.Warn("mute gesture found no gate to turn")`
并把 `toggleMute` 给的那句原样打印——而那句**仍然是从门读回来的真话**，所以用户可见输出不错、只有**日志记录**会变成"找不到门可拧"。
⇒ 严重性：低（可观测性面，非用户可见错状态），但具名交回：这不是"门存在但不可调用"那类，也不是"仪器里没造出那个世界"那类，
是**那一枚函数值契约的一半没被钉**。落成新票建议名册的一项（§5）。

## 4. N3 · 恒真面：被测的那枚门是生产装配的吗

逐条现读（⛔ 不采信实现者 `resident_mute_290_windows_test.go:11-15` 那句自述）：
- **门＝生产装配的那枚**：台件走 `assembleCapture(rt, dir, tap.call, <source 工厂>)`
  （`resident_mute_290_windows_test.go:88`），而生产入口 `startResidentAudio`→`buildResidentAudio`
  （`resident_audio_windows.go:199-200`、`:225-226`）＝`assembleCapture(rt, dataDir, out, newRealCaptureSource)`——
  **同一个函数**，只差 source 工厂那一枚注入位（`:229-230` 的形参）。门在 `:273` 由
  `audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))` 造、`:277` 挂到 `ra.gate` ⇒ **不是手搓的门，是生产那个构造点造的门**。
  注入的只是**设备那一侧**（票 247 已有的 `captureSource` 接缝），且台件用的是**真门**（未 mock `HalfDuplexGate`）。
- **但装配根那一段没有被走**：生产挂载点是 `cmd/wisp/resident_windows.go:303` `rb.attachMuteGate(raudio.toggleMute)`，
  台件是在 `resident_mute_290_windows_test.go:93`／`:206` **手抄了同一句话**（`rb := &residentBall{}` + `rb.attachMuteGate(ra.toggleMute)`）。
  `runResident`（装配根本体）**在任何测试里都没被调用过**（尺：`grep -rn "runResident\b" --include=*_test.go cmd/` → 零命中，rc=1）。
- **球侧两支闭包也没被走**：台件直接调 `rb.muteGesture("mute-hotkey")`（`:117`），
  跳过 `resident_ball_windows.go:321`/`:325` 那两枚 `Events` 闭包；`internal/ball` 的 `case hkMute:`（`ball_windows.go:660`）与
  `case menuMute:`（`:678`）更要真窗口＋真消息泵才会走到，台件里一枚都没走。
- ⚠ **本腿一度把这条写重了，现跑后更正**（以盘上原文为准）：仓里**今天已经有**能造"生产形状"的两件仪器，
  290 只是没用它们——
  ① 真窗口形状：`cmd/wisp/resident_hotkey_258_windows_test.go:266` `Test258BridgeRebindsLiveKeysFromConfigEdit`
     在 `:278` 逐字 `rb := startResidentBall(observe.NewRegistry(), nil, src, src)` 起**真球窗**、`:289-291` 读 `rb.b.HotkeyReport()` 的注册表，
     而且它的 `writeConfig258` 那一段（`:272-275`）**逐字已经写了 `mute = "Ctrl+Alt+M"`**；
     `t.Skipf("SKIP-LOUD…")` 那一支在本机**没有触发**（整包红/跳名册里既无它也无不跳＝本机真拉起过球窗）。
     ⇒ "**mute 那枚热键到底注册成功了没**"这一跳今天**有现成的尺可读**，只是没人为 `hkMute` 读它。
  ② 装配根形状：仓里有成熟的 **对 `runResident` 函数体做 AST 走查**的先例三枚
     （`resident_ball_258` 族：`cmd/wisp/resident_hotkey_258_test.go:25/:33/:45-68`、`resident_ball_228_test.go:166-247`、
     `resident_approval_risk_256_windows_test.go:527-569`），专治"`runResident()` 永不返回所以不能直调"这件事
     （`resident_ball_228_test.go:17` 逐字：`Why a source walk and not a runtime call: runResident() never returns (it`）。
- ⇒ **具名结论**："装配根把执行者挂上（`resident_windows.go:303`）"与"球侧那两枚闭包真被路由到 `muteGesture`"这两条判据，
  **在生产的形状上今天没人验过**：`grep -rn "attachMuteGate" --include=*_test.go cmd/wisp/` → 3 命中全在
  `resident_mute_290_windows_test.go`（`:93`/`:206`/`:235`，全是**直调**），**无一枚走查 `runResident` 体**；
  `runResident` 本身在任何测试里都**没有被调用过**（只被 AST 走查与子进程名引用）。
  缺口命名：它是①"门存在但不可调用"的**续集**（今天可拧了，但"拧门那只手是装配根挂上去的"这件事只由 grep 尺与一句注释保证），
  而不是"守卫分岔所需的世界仪器里没造出来"那一类——**仪器在仓里躺着，没被指向这一跳**。

**要让它在生产也可见，台件该怎么造**（⛔ 本腿没造、⛔ 不要 `frontend/**`、⛔ 不加新 seam——两件现成仪器都够用）：
```go
// (a) 装配根那一跳：照 cmd/wisp/resident_hotkey_258_test.go:25-68 的 runResidentBody258 走查形状加一枚，
//     断言三件事，全部读 AST、不需要跑 runResident：
//       1) runResident 体内 callExpr 里 id.Name == "attachMuteGate" 的次数 == 1（不是 0、不是 2）；
//       2) 它接收者是 startResidentBall 那枚返回值 rb，实参是 raudio.toggleMute（点名 selector 链，不接受"任意函数值"）；
//       3) 它在体内的位置**晚于** startResidentAudio 那一次调用（后置 setter 的时序是这条链成立的前提，
//          resident_windows.go:217 建球早于 :278 建门早于 :303 挂载）。
// (b) 球侧注册那一跳：复用 resident_hotkey_258_windows_test.go:266-294 的真球窗形状（含它那句
//     SKIP-LOUD 的 headless 退路），在 :289 那把 HotkeyReport 尺上**对 mute 那枚 id 现读 IsLive**，
//     红/跳都具名——今天 hkMute 的注册成功与否没有任何用例读过。
// (c) 若还要那一支闭包本身：同一枚 live-ball 用例里从 rb 侧把 Events.OnMuteHotkey/OnTrayMute 取出来
//     直接 fire（⛔ 不 PostMessage 真热键、⛔ 不 press 键盘），断言读 ra.gate.Muted()/Open() 与 source 的 start 计数，
//     ⛔ 不读返回值——理由见本件 §2 的 M1：读返回值会被"句子里写了成功"骗过去。
```

## 5. N4 · `69d8c9ef` 那七处重指是不是只改了数

`git show 69d8c9ef` 实跑＝**7 行改动**（`+7/-7`）、**7 处 cite 重指**、指向 **3 枚不同的目标行号**。⛔ 无任何 token 文字、verdict 前缀、断言改动。
逐处（⛔ 未抽样）核"那一行确实含方括号里的 token"：

| # | 引用所在行（`cmd/wisp/config_readers_255.go`） | 重指后 cite 的目标 | 目标行现读（逐字） | token 命中 | 这把尺是否检查它 |
|---|---|---|---|---|---|
| 1 | `:121`（Go 注释） | `resident_ball_windows.go:316 [Hotkeys:  cfg,]` | `:316` = `		Hotkeys:  cfg,` | 是 | **否**（注释不在 `hotRowClaims` 里） |
| 2 | `:129`（`"hotkey"` 行） | 同上 `:316 [Hotkeys:  cfg,]` | 同上 | 是 | 是 |
| 3 | `:150`（`"audio"` 行） | `resident_audio_windows.go:256 [c.Audio.MicMutedDefault]` | `:256` = `	ra.mutedAtBoot = c.Audio.MicMutedDefault` | 是 | 是 |
| 4 | `:199`（`"voice.tts.speed"`） | `resident_audio_windows.go:257 [c.Voice.Enabled]` | `:257` = `	ra.voiceEnabled = c.Voice.Enabled` | 是 | 是 |
| 5 | `:200`（`"voice.punctuation"`） | 同上 `:257` | 同上 | 是 | 是 |
| 6 | `:201`（`"voice.wake_word.thresholds"`） | 同上 `:257` | 同上 | 是 | 是 |
| 7 | `:202`（`"voice.wake_word.veto_words"`） | 同上 `:257` | 同上 | 是 | 是 |

- 旧行号残留尺：`grep -n "resident_ball_windows.go:276\|resident_audio_windows.go:19[67]" cmd/wisp/config_readers_255.go` → **零命中**（rc=1）。
- **M3 正控（证"这把被编辑过的尺仍然有牙"）**：把 `:150` 那枚 `:256` 故意写成 `:255`（`:255` 是空行）。
  ```
  PATH=… go test ./cmd/wisp/ -run 'TestTicket255' -count=1 -v
  rc_m3=1
  --- FAIL: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim     ← 点名的那一把，红
  （同一次 -run 里其余 19 枚 TestTicket255* 全绿，含三枚 receipt 文案用例 ⇒ 红只归 cite 那一枚）
      config_receipt_255_test.go:209: row "audio" cites cmd/wisp/resident_audio_windows.go:255 for
      "c.Audio.MicMutedDefault", that line now reads "" - the evidence drifted…
  ```
  尺的实现现读：`config_receipt_255_test.go:182` 遍历 `hotRowClaims`、`:191` 取 `(rel, lineNo, token)`、`:207` `got := strings.TrimSpace(lines[n-1])`、
  `:208` `if !strings.Contains(got, token)` 才报错 ⇒ **行号写歪一位必红**，实测证实。
  ⚠ 同时现读到这把尺的两个**已知软面**（都不是本票账，具名给编排者）：
  ① `:121` 那一处重指是**注释**，`evidenceCite` 只扫 `hotRowClaims` 的值 ⇒ 注释里那枚 `:316` **没有任何尺盯着**，写歪不会红；
  ② `:220` 的 `citedRowsFloor = 8` 是**总数下限**，个别 cite 被改成 `扫描零命中` 标记不会触到它（要 ≥9 处都退化才会红）。
- ⇒ **"只改了数"这一问：成立**（7/7 处只动行号，token/句式/断言零改动，且被编辑过的那把尺 M3 实测仍红得下来）。
  派单那句"若这发正控红不了就判 `AC#4`/`AC#5` 相关格不成立"——**正控红了**，故不触发。

## 6. 还原证明（⛔ 三发突变都没留在跟踪文件里）

```
md5sum（突变前留存 vs 还原后现读，逐字节等值）：
  cmd/wisp/resident_audio_windows.go  fe196316985a9b739cfa4006e965f067  ==  fe196316985a9b739cfa4006e965f067   (M1/M2 还原)
  cmd/wisp/config_readers_255.go      2b86241559822ea34dc60bda02bdab8a  ==  2b86241559822ea34dc60bda02bdab8a   (M3 还原)
另一把独立尺（对 HEAD 的 blob 比，而非只对自己比）：
  git show HEAD:<file> | md5sum  ==  md5sum <file>   → 两枚文件都 YES
工作树尺：git status --porcelain -- cmd internal   → 空（0 行）；git diff --stat HEAD -- cmd internal → 0 行
```

## 7. 建议落成新票的名册（⛔ 本腿未动任何产码，交编排者立）

1. `SetTrayChecks` 的生产调用者＝0 枚 ⇒ 托盘静音勾号与门脱钩（N5 正文在 `20-questions-n5-n6.md`）——**用户看得见的那一枚**。
2. `executed` 那一枚函数值契约的"真拧了门"一半无人断言（N2·M2 实测）。
3. 装配根 `runResident`＋`ball.Events` 两支闭包无生产形状用例（N3，含台件写法）。
