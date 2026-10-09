# 293-a1 · 30 · AC#0 ⓒ — 全仓有没有**第二处**已经在替托盘维护静音状态的东西？

**锚点**＝`f019f45e`。⛔ 结论先说：**没有。** 但这一节的实际交付不是那句"没有"，而是**四枚"看起来像、其实不是"的东西**具名列出（编排者要防的正是"另造一枚其实已经躺在仓里"）。

---

## 1. 票面指定的语义变体名册（逐枚现跑，尺＝`git grep -n <pat> HEAD -- '*.go' '*.ts' '*.tsx' '*.js' '*.html' '*.css'`，已滤 `.scratch/`）

| 变体 | 命中（逐字） | 是不是"维护静音状态的第二处" |
|---|---|---|
| `trayMuted` | `internal/ball/ball_windows.go:128`（声明）· `:674`（读进菜单）· `:955`（`SetTrayChecks` 体内写）· `cmd/wisp/resident_windows.go:301`（**注释**，见 §3） | ⛔ 否。字段只有这一枚，写入者只有那座 0 调用者的桥 |
| `menuMute` | `internal/ball/tray_windows.go:21`（`menuMute = 2` 常量）· `:88`（`appendItem(menuMute, "静音", muted)`）· `internal/ball/ball_windows.go:678`（`case menuMute:` → `b.fire(...OnTrayMute)`） | ⛔ 否。命令 id，不存状态 |
| `mfCheckd` | `internal/ball/tray_windows.go:82`（`flags \|= mfCheckd`）· `internal/ball/win32_windows.go:136`（`mfCheckd = 0x00000008` 常量） | ⛔ 否。**全仓唯一一处往菜单里塞 checked 的地方**就是 `:82`，而它的实参只可能来自 `showMenu` 的两枚入参 |
| `SetTrayChecks` | `internal/ball/ball_windows.go:952`（注释）· `:953`（定义） | ⛔ 否。定义体＝两行赋值，调用者 0（见 `10-...md` §1.3） |
| `showMenu` | `internal/ball/tray_windows.go:69`（注释）· `:72`（定义）· `internal/ball/ball_windows.go:674`（**全仓唯一调用点**） | ⛔ 否。★但它给出本节最重要的一形，见 §2 |
| `pauseWake` 一族（本腿自补）`PauseWake\|pauseWake\|pause-wake\|pausedWake` | 见 `10-...md`：`tray_windows.go:22/:89`、`ball_windows.go:43/:57/:129/:674/:680/:681/:956`、`cmd/wisp/resident_ball_windows.go:326/:456`、`cmd/balldebug/main.go:201-202`、`cmd/wisp/resident_ball_228_test.go:52` | ⛔ 否。全是从 `SetTrayChecks` 那条死桥或从事件路由走过的路，⛔ 无一处另存状态 |
| `SetTrayTip` / `setTip`（本腿自补，同族第三面） | `internal/ball/ball_windows.go:960-963`（定义，调用者 0）· `internal/ball/tray_windows.go:54-60`（`tray.setTip` 写 `szTip` 后 `NIM_MODIFY`） | ⛔ 否。tooltip 面今天**也没人驱动**，与两枚勾同病；它不存静音态 |

⇒ 名册扫完的硬结论：**托盘静音态在全仓只有一个存储位置（`ball_windows.go:128`），一个写入者（`:955`，隶属 0 调用者的 `SetTrayChecks`），一个读取者（`:674`）**。**没有第二处。**

---

## 2. ★本节最有用的一形：菜单是**每次右键现建**的（不是缓存态）

尺＝`sed -n '669,686p' internal/ball/ball_windows.go`（`10-...md` 已录逐字）＋ `sed -n '1,60p' internal/ball/tray_windows.go`。

- `internal/ball/ball_windows.go:673-674`：`case wmRButtonUp:` → **每次**都新调一次 `showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)`。
- `internal/ball/tray_windows.go:26-28`：`type tray struct { data notifyIconData }` ⇒ **`tray` 结构体里没有任何静音字段**，它只装 Shell 的 `NOTIFYICONDATA`。

⇒ 事实交给编排者（⛔ 不是本腿的建议）：**勾的正确性只在"右键那一刻"需要成立**，不是连续需要。⇒ 于是 `AC#1` 那一跳逻辑上有两族形状：**推**（门态一变就调 `SetTrayChecks`）与**拉**（`showMenu` 那一刻现读门态）。票面 `AC#2` 判据写的是"门态一变，勾立刻镜像"＋点名三种起点，走的是"推"这一族；而 setter 的形状（`SetTrayChecks(muted, pausedWake bool)`）也是"推"。**"拉"这一族今天不需要改 setter 签名就能成立**——这一条本腿只写事实、⛔ 不裁该走哪族。

---

## 3. `cmd/wisp` 里已经写死的**一处立场**（票面没提，必须让编排者知道）

`cmd/wisp/resident_windows.go:293-303`（`sed -n '295,320p'` 现读，逐字）：
```go
// (D43's table has no Sleeping -> Muted edge, and inventing one is a
// ...
// names that follow-up as outside this installment), no second "who is muted"
// flag anywhere
// (the gate's own flag is the truth source, internal/audio/gate.go:20-21), no
// new goroutine (the gestures fire on the ui-sta thread and the device open is
// synchronous inside internal/audio, on the roster name ticket 247 already
// booked), and no mirror of the gate into the tray's own checkmark - the
// ball-side trayMuted display flag stays undriven here rather than becoming a
// second authority nobody reconciles.
rb.attachMuteGate(raudio.toggleMute)
```
⇒ 判读：**票 290 的落地腿当时是"知道这枚勾、并显式决定不去镜像它"**，理由写的是"免得变成第二座没人调和的权威"。
⇒ ⚠ 这对票 293 是**直接承重**的：本票要的正是"让勾镜像门"，⛔ 编排者要么在 `AC#1` 裁定里**具名推翻这段注释的理由**（并说清为什么现在不构成第二权威——凭据现成：`AC#2` 要求勾**读回 `gate.Muted()`** 而不是另存一枚），要么承认这段注释与新要求冲突、把它一并改写。⛔ 本腿不替他裁，但**这段注释不指出就等于没说**。
⇒ 同族第二处立场写在 `cmd/wisp/resident_ball_windows.go:40-51`（逐字）：
```go
// The gate's own muted flag is the truth source for "who is muted"
// (internal/audio/gate.go:20-21); this host keeps no copy of it - see
// muteGesture, whose sentence is read back off the gate after the write.
// What did NOT change here: no state machine event is dispatched (D43's
// Sleeping -> Muted edge does not exist, and inventing one is a human-approved
// contract change, not a wiring detail), and no second "am I muted" flag.
```

---

## 4. 四枚"看起来像第二处、其实不是"的东西（编排者下刀前先认识它们）

| # | 东西 | 位置 | 为什么不是"第二处维护者" | ⚠ 风险备注 |
|---|---|---|---|---|
| 1 | **`ra.mutedAtBoot`** | `cmd/wisp/resident_audio_windows.go:101`（声明）· `:256`（唯一写入＝`c.Audio.MicMutedDefault`） | 尺＝`git grep -n "mutedAtBoot" HEAD` ⇒ 全仓（去探针/去台账）**只有这 2 命中，读＝0 枚** ⇒ **write-only 死字段**，今天没替任何东西维护状态 | ★**这才是本票真正该防的那一枚**：它长得就像"我以为的静音"，而它的注释（`:96-99`「recorded so the boot line and any later reading of this handle state the same privacy posture」）**声称**自己在为可读性服务、实际从没被读过。⛔ 落地腿**绝不能**把勾接到 `mutedAtBoot`——它在 `toggleMute` 之后**不再更新**（`:156-190` 全文只回字符串、不动 `mutedAtBoot`），接上去＝用户刚开麦、勾还按出厂值说"静音"＝**把错状态从"永不动"换成"永远错"**。这正好是票面 `AC#2` 禁的"另存一枚我以为的静音"。⚠ 兄弟字段 `voiceEnabled`（`:100`/`:257`）同样是 0 读者（尺＝`git grep -n "voiceEnabled" HEAD -- cmd`）⇒ 这一对是同一发留下的，⛔ 不属本票射程，具名给编排者记账即可 |
| 2 | **`ra.verdict` 的 boot 静音句**＋`posture()` | `resident_audio_windows.go:298`（`采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）`）→ `:375 posture()` → `resident_windows.go:320`（尺＝`grep -n "采集腿：%s" cmd/wisp/resident_windows.go`）`fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())` | 它是**已在维护静音态的那一处**，但**面的不是托盘**、且**只在 boot 那一刻取一次**（`toggleMute` 不改 `ra.verdict`）⇒ 不是"第二处托盘维护者"，是"第一处 boot 维护者" | ⚠ 与本票**共存不冲突**；落地腿⛔ 不许把它当勾的来源（boot 快照 ≠ 当前门态）。它同时是票面"console 那句不是兜底"（`resident_windows.go:55-57` 双击无终端）的同一面 |
| 3 | **`WithGateEvents` 那枚现成接缝** | `internal/audio/gate.go:61-64`（定义）· `:176-178`（发 `EventMuted`/`EventUnmuted`） | 尺＝`git grep -n "WithGateEvents" HEAD` ⇒ 产码调用者＝**0**，只有 `gate_test.go:86`/`:127` 两枚测试 | ⛔ 不是"已有人在维护"，而是"**门侧本来就预留了把静音事件推出去的口子，没人接**"。这一条**编排者早已记账**（`docs/reports/pending-and-issues.md:14829`＋`docs/reports/HANDOVER.md:2018` 逐字"现成接缝 `WithGateEvents` 产码调用者＝0"）⇒ 本腿**复跑并确认**，⛔ 不新报 |
| 4 | **球态 `Muted`（灰球＋斜杠）那套视觉机器** | `internal/ball/statevisual.go:15`（`IconSlash`）· `:160`/`:290`（`case statemachine.StateMuted:`）· `internal/ball/tokens.go:448`（`SlashStrokePx = 1.5`） | 尺＝`git grep -n "EvMuteKey" HEAD` ⇒ 产码发射者只有 `cmd/balldebug/main.go:603`；`internal/statemachine/table.go:65/:75/:79` 三枚边里**没有 `Sleeping→Muted`**；常驻球初始态 `resident_ball_windows.go:315` `Initial: statemachine.StateSleeping,` | ⇒ 球那一面**机器齐备、没人开火**，⛔ 不构成第二处托盘维护者。且**票面 `AC#3` 已明写本票不许把它打进 D43**（需人工批准的 `A##`）。同三条读数在 `pending-and-issues.md:14829` 已由编排者亲跑记录 ⇒ 本腿是**复现**不是新量 |

---

## 5. ⛔ 没有停在这些目录之外——`scripts/`、`tools/`、`frontend/`、票池也都扫了

| 射程 | 尺 | 读数 |
|---|---|---|
| `scripts/`＋`tools/` | `grep -rn "trayMuted\|SetTrayChecks\|mfCheckd\|showMenu\|menuMute\|trayPauseWake" scripts tools` | **0 命中**（逐字 `wc -l` ＝ `0`） |
| `scripts/` 里凡提托盘的 | `grep -rln "托盘\|tray\|Tray" scripts` → `scripts/dev/ball-cycle.ps1:46`（只讲"把球放在托盘角上方的净空区"，⛔ 无菜单无勾）·`scripts/spike/common/winshell.go:367-421`（`AddTrayIcon`/`Remove`，只有 `NIM_ADD`/`NIM_DELETE`，⛔ **无 `TrackPopupMenu`/`AppendMenu`/`MF_CHECKED`**，尺＝`grep -n "TrackPopupMenu\|AppendMenu\|MF_CHECKED" scripts/spike/...` ＝ 0 命中）· `scripts/spike/{shell-baseline,xy-verdict}/main.go`（同一族 spike 件）· `scripts/spike/bin/*.exe`（编译产物，⛔ 源码面） | ⛔ 无一维护静音态 |
| `frontend/`（面板那一侧会不会也显示静音） | `grep -rn "mute\|Mute\|静音" frontend/src frontend/scripts` | 命中**全是 Tailwind 颜色 token**（`text-muted-foreground`／`--tooltip-muted`／`--muted`）与 `frontend/src/fixtures/harness.ts:228` 一枚**状态视觉表夹具**（`{ state: "Muted", visual: "灰球 + 斜杠", enter: "一键静音（D16）", exit: "再按静音键" }` ＝ `PLAN.md:1312` 的抄件，走球那一面）。⇒ 面板/前端**没有任何一处**显示或维护麦克风静音态 |
| `internal/panel`（C17 桥有没有静音方法） | `grep -rn "mute\|Muted" internal/panel/*.go \| grep -v _test` | **0 命中** ⇒ 面板桥没有静音面，⛔ 不构成第二处 |
| 票池（谁还可能来动这枚勾） | `grep -rn "SetTrayChecks\|trayMuted\|那枚勾\|打勾" .scratch/wisp/issues/*.md` 去掉 290/293/294/295 | 唯一实质命中＝`228-...md:139`（**既有事实记录**：「`SetTrayChecks`／`SetTrayTip` 全仓零枚生产调用者 … 今天那枚勾在产码里恒为 `false`」，本腿复跑成立）。⇒ ⚠ **票 228 只是量到了它、没有领它**——那枚票的射程是"球和托盘不在跑任务的那个进程里"，与本票的"勾跟不跟门"是两件事。⇒ 残余两枚在飞票：`294`（静音挂载那一跳的生产形状测试）、`295`（`ra.started` 裸 bool 跨线程写）——**都不主张驱动这枚勾**，但 `294` 与本票**同写 `cmd/wisp` 测试面**、`295` 与本票**同碰 `resident_audio_windows.go`**（票 293 排程段已写 294 串行；⚠ 295 与本票的写面重叠**票面没写**，见 `90-...md` 自报）。 |

---

## 6. 一句话回票面

**没有第二处。** 但要说清"为什么没有"才算数：仓里躺着四枚**形状相近**的东西（`mutedAtBoot` 死字段 / boot 的 `ra.verdict` 静音句 / 没人接的 `WithGateEvents` / 没人开火的球态 `Muted`），⛔ 落地腿若从这四枚里任一枚取勾，得到的都不是当前门态（分别是：出厂值快照、boot 快照、无人注册的回调、进不去的态）⇒ **`AC#2` 那句"勾必须读回 `gate.Muted()`"是这四枚的唯一排他凭据，建议落地腿逐枚具名排除，别只写"没找到"。**
