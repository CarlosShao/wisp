# 293-a1 · 00 · 起手锚 + 尺清单 + 行号现量对账

腿＝`293-a1`（只读普查）。做的格＝票 293 的 **AC#0（ⓐⓑⓒ）** 与 **AC#1**。⛔ 未动任何产码，⛔ 未勾任何框。

---

## 1. 起手锚（`date` / `git log` / `git status` 的 stdout 逐字）

命令 1：`date`
```
Fri Oct  9 19:10:55 CST 2026
```

命令 2：`git log --oneline -1`
```
f019f45e 收 livewin-v1（A788＋§4.0as）：票 247 AC#4 按非实现者判语翻勾带两限、AC#2 判不成立改归口 297、AC#6 不可判＋设备被占那形欠我自己；findings 三处数具名更正；Q-82/Q-83 补登台账（记我"先宣布已上清单、当时台账零命中"）；票 296 行号 :55→:54 打旧
```

命令 3：`git status --porcelain -- .scratch/wisp/probes/293/`
```
warning: could not open directory '.scratch/wisp/probes/293/': No such file or directory
```
⇒ 起手时该目录**不存在**（零产出、零脏件）。落点目录由本腿 `mkdir -p .scratch/wisp/probes/293/a1` 新建，**只建不删**。

命令 4：`tasklist //FI "IMAGENAME eq wisp.exe"`
```
INFO: No tasks are running which match the specified criteria.
```
命令 5：`tasklist //FI "IMAGENAME eq balldebug.exe"`
```
INFO: No tasks are running which match the specified criteria.
```
⇒ 两枚计数都＝**0**，与票面要求一致。本腿**没有**杀过任何进程（起手即零，无需处置）。

---

## 2. 本腿跑过的全部命令（逐把具名，可复跑）

| # | 命令（逐字） | 用途 |
|---|---|---|
| C1 | `date` | 起手锚 |
| C2 | `git log --oneline -1` | 起手锚 HEAD |
| C3 | `git status --porcelain -- .scratch/wisp/probes/293/` | 起手锚落点 |
| C4 | `tasklist //FI "IMAGENAME eq wisp.exe"` | 进程计数 |
| C5 | `tasklist //FI "IMAGENAME eq balldebug.exe"` | 进程计数 |
| C6 | `git grep -n "SetTrayChecks" HEAD -- internal cmd` | 复跑票面 R1（AC#0 现量那把尺） |
| C7 | `git grep -n "SetTrayChecks" HEAD` | 含票池的全仓尺（发现被票池文本淹没，故有 C8） |
| C8 | `git grep -n "SetTrayChecks" HEAD -- '*.go'` | 含测试的全仓 **Go 源**尺（AC#0 ⓐ 要求同形） |
| C9 | `git grep -n "trayPauseWake" HEAD -- internal cmd` | AC#0 ⓐ 主尺 |
| C10 | `git grep -n "trayPauseWake" HEAD -- '*.go'` | AC#0 ⓐ 含测试全仓尺 |
| C11 | `git grep -n "OnTrayPauseWake\|OnTrayMute\|OnMuteHotkey" HEAD -- '*.go'` | AC#0 ⓐ 事件半 |
| C12 | `git grep -n "PauseWake\|pauseWake\|pause-wake\|pausedWake" HEAD -- '*.go' \| grep -v probes/` | AC#0 ⓐ 语义变体尺 |
| C13 | `git grep -n "attachMuteGate\|currentMuteGate\|muteGestureFunc" HEAD -- cmd internal` | AC#1 既有形状 |
| C14 | `sed -n '50,60p;295,320p' cmd/wisp/resident_windows.go` | AC#1 行号现量（只读） |
| C15 | `sed -n '14,30p' internal/audio/gate.go` | 票面 `gate.go:20-21` 行号现量 |
| C16 | `sed -n '272,282p' internal/config/schema.go` | 票面 `schema.go:277` 行号现量 |
| C17 | `sed -n '155,175p' cmd/wisp/resident_audio_windows.go` | 票面 `:162` 行号现量 |
| C18 | `grep -n "muted\|Muted" internal/audio/gate.go` | AC#0 ⓒ 门侧真相源名册 |
| C19 | `git grep -n "SetTrayTip" HEAD -- '*.go' \| grep -v probes` | AC#0 ⓒ 同族 setter |
| C20 | `grep -n "func (ra \*residentAudio) posture\|residentStatusLine" cmd/wisp/*.go` | AC#0 ⓒ 已有的状态出口 |
| C21 | `GOFLAGS= go list -deps ./internal/ball` | AC#1 硬约束尺（见 `40-ac1-landing-shape.md`） |
| C22 | `go env GOPATH GOFLAGS GOOS GOARCH` | 环境自报（未跑 build/vet/test） |

⛔ **未跑**（按派单禁令，让给 `292-v1` 的编译面）：`go build`／`go vet`／`go test`。⛔ 未跑 `gofumpt`／`d22scan`／`go test` 类门禁四数——那四数是票面 `AC#4`（落地腿）的格，不属本腿。

---

## 3. 票面行号快照 vs 本腿现量（⚠ 逐枚复跑，不符具名报漂移）

尺＝`git grep -n` / `sed -n` 在锚点 `f019f45e` 现跑。

| 票面主张 | 票面行号 | 本腿现量 | 判定 |
|---|---|---|---|
| `trayMuted bool` 字段声明 | `internal/ball/ball_windows.go:128` | `128: trayMuted     bool` | ✅ 未漂 |
| `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` | `:674` | `674` 逐字一致 | ✅ 未漂 |
| `appendItem(menuMute, "静音", muted)` | `internal/ball/tray_windows.go:88` | `88` 逐字一致 | ✅ 未漂 |
| `if checked {` / `flags \|= mfCheckd` | `:81-82` | `81`/`82` 逐字一致 | ✅ 未漂 |
| `SetTrayChecks` 注释＋定义两命中 | `ball_windows.go:952`＋`:953` | `952`/`953` | ✅ 未漂 |
| `OnTrayMute: func() { rb.muteGesture("tray-mute") }` | `resident_ball_windows.go:325` | `325` 逐字一致 | ✅ 未漂 |
| `outcome, executed := fn()` | `:527` | `527` | ✅ 未漂 |
| `OnMuteHotkey: func() { rb.muteGesture("mute-hotkey") }` | `:321` | `321` | ✅ 未漂 |
| `fmt.Printf("wisp: ball %s: %s\n", name, outcome)` | `:536` | `536` | ✅ 未漂 |
| `Ticket 117:` 两行注释开头 | `resident_windows.go:55-57` | `55` 起逐字一致 | ✅ 未漂 |
| `gate.SetMuted(!gate.Muted())` | `resident_audio_windows.go:162` | `162` | ✅ 未漂 |
| `MicMutedDefault bool` 那行 | `internal/config/schema.go:277` | `277` | ✅ 未漂 |
| 门侧 truth source "写死" | `internal/audio/gate.go:20-21` | ⚠ **票面这一枚指偏**：`:20-21` 是**包 doc 注释**里 "Muted (both paths, user intent wins)" 两行；**字段本体**在 `gate.go:79`（`muted bool`），读回在 `:202-205`（`Muted()`），写在 `:168-174`（`SetMuted`）。⇒ 语义主张（门是主）**成立**，⚠ 但引用的行号指的是注释不是字段。具名报此偏：doc 注释 vs 字段，**不是同一枚位置**。 |
| 后置 setter 既有形状 | `resident_windows.go:303` | `303: rb.attachMuteGate(raudio.toggleMute)` | ✅ 未漂 |

⇒ 结论：**票面 14 枚行号里 13 枚未漂**；唯一一枚偏了的是 `gate.go:20-21`（指到 doc 注释、不是字段本体），且偏的是**出处标注**不是**事实**——门侧为真相源在 `gate.go:79/168/202` 一样钉得住。

---

## 4. 名册差集的排他声明（防把探针副本算进名册）

全仓 `*.go` 尺扫到两枚 `.scratch/wisp/probes/**` 下的**源文件副本**，它们会污染"调用者枚数"这类计数，已从名册排除：

- `.scratch/wisp/probes/33/r8b/logs/pristine-ball_windows.go` —— `ball_windows.go` 的**整份历史拷贝**（含 `SetTrayChecks` 定义 `:918-919`、`trayPauseWake` `:129`/`:674`/`:922`）。排除理由＝探针目录里的 pristine 快照，非产码，不参与任何 build。
- `.scratch/wisp/probes/258/v1/mut1/cmd/wisp/resident_ball_windows.go` —— 票 258 验收腿的**变异拷贝**（含 `OnTrayMute: func() { recordBallGesture("tray-mute") }` 的**改前形状**）。排除理由＝变异夹具。

**排除枚数＝2 枚**（另有同名 pristine/变异件若再出现，一律按"探针目录＝非产码"处理）。⛔ 这 2 枚的命中**没有**被算进任何一节的"调用者枚数"。

---

## 5. 交付件名册（本腿写面全在此目录）

| 件 | 格的 |
|---|---|
| `00-anchor-and-rulers.md` | 起手锚 + 尺清单 + 行号对账 + 排他声明 |
| `10-ac0a-pause-wake-driver.md` | AC#0 ⓐ |
| `20-ac0b-contract-checkmark.md` | AC#0 ⓑ |
| `30-ac0c-second-surfacer-scan.md` | AC#0 ⓒ |
| `40-ac1-landing-shape.md` | AC#1 |
| `90-unrun-rulers.md` | 自报：没跑成的尺／判不了的格 |
