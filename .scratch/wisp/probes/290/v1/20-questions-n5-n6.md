# 290-v1 · 20 N5 托盘勾号／N6 重写句／N8 禁区与越界（含 0 字节证据件那一格）

## N5 · 托盘那枚勾（`trayMuted`）现在是不是在说谎——是，而且这次它有后果了

现读三件（⛔ 未引用实现者自报，尺都是本腿现跑）：
- 定义：`internal/ball/ball_windows.go:128` `trayMuted     bool`。
- 被读进菜单的那一行：`internal/ball/ball_windows.go:674` `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)`
  → `internal/ball/tray_windows.go:88` `appendItem(menuMute, "静音", muted)` → `:81-82` `if checked { flags |= mfCheckd }`。
  ⇒ 这枚 bool 就是**菜单上那个勾**，且它是**用户看得见静音状态的唯一一面**（另一面是 `fmt.Printf` 那句 console）。
- setter 的生产调用者尺：
  ```
  grep -rn "SetTrayChecks" --include=*.go . | grep -v "^\./\.scratch" | grep -v "_test.go"   → rc=0，2 命中（:952 注释／:953 定义）
  含测试的全仓尺（internal/ cmd/）同样 2 命中 ⇒ 调用者＝0 枚（改前 0、改后仍 0）
  ```
  ⇒ `trayMuted` 恒为零值 `false` ⇒ "**静音**"这一项**永远不带勾**。

**这一发改变的是"这一项会不会动手"，不是"勾准不准"**：
- 改前：`case menuMute:`（`:678`）→ `OnTrayMute: func() { recordBallGesture("tray-mute") }`（旧逐字见票面 `:13-15`）
  ⇒ 点它**什么都不做**，勾也不动。用户得到的信号是"点了没反应"＋"没勾"——两者一致地指向"麦克风没开"，而默认档**确实**没开。
- 改后：`:325` `OnTrayMute: func() { rb.muteGesture("tray-mute") }` → `:527` `outcome, executed := fn()` → `toggleMute` → `SetMuted`
  ⇒ 点它**真的拧门**，勾**仍不动**。

⇒ **裁决（N5 问的那一判）：这是这次新造的、用户看得见的错状态——但只造了"有后果"那一半；"勾本身不准"是既有缺口。**
拆开来写清，免得下一腿把两件事混成一件事：
1. **既有（不归本票）**：`trayMuted` 自始未驱动，而默认档 `mic_muted_default=true`（`internal/config/schema.go:277` 现读 `default:"true"`）
   ⇒ **boot 时就已经在说谎**：门关着、菜单不打勾。这一格 `fdee536b` 之前就存在，`SetTrayChecks` 的调用者改前改后都是 0。
2. **新造（归本票这一发）**：这一发让"静音"项第一次**能改变设备状态**，于是那枚永远不动的勾从"惰性的错"变成"**有后果的错**"——
   它现在与一件**已经发生的事**并排显示，而不是与一件从未发生的事并排显示。

**最坏形状（用户会误以为什么）**：
- 形状 A（隐私面，最坏）：用户双击启动 `wisp.exe`（boot＝静音位、设备未打开、菜单"静音"无勾）→
  右键托盘点"静音" → 门开、`WASAPIMicrophone` 交接设备、采集线程在跑；
  他再右键打开菜单想确认状态 → **"静音"仍然无勾** ⇒ 他**看不出自己刚把麦克风开了**；
  过一会儿再点一次（本意"看看会不会怎样"）→ 门被拧回静音。他把一件真的改变了采集状态的动作当成一个没反应的装饰。
  ⛔ 这里没有 console 兜底：`cmd/wisp/resident_windows.go:55-57` 逐字
  `// Ticket 117: the resident process is the leg owner actually uses - double` / `// click the icon, no terminal attached, stderr going nowhere. Everything`
  ⇒ `muteGesture` 那句 `fmt.Printf("wisp: ball tray-mute: %s\n", outcome)`（`resident_ball_windows.go:536`）**在双击启动的形状里根本没有读者**。
- 形状 B（误以为还在听）：他按到静音位后菜单仍无勾 ⇒ 若他把"无勾"读成"没静音"，会以为 Wisp 还在采集——这一向危害小些（多疑不少信），但同样是错的。
- 形状 C（热键那支同罪）：`Ctrl+Alt+M` 走的是同一支 `muteGesture`，同样只剩 console 那句（无常驻终端时进日志文件），
  托盘勾不镜像 ⇒ **两面（菜单勾 / 日志句）给出的确定度不一样**，勾是零信息的那一面。

**⛔ 本腿未顺手修**（未动 `internal/ball` 一字、未动 `SetTrayChecks`、未加镜像）。落成新票的建议名册（交编排者立）：
1. `29x-tray-checkmark-never-follows-the-gate-mute-item-acts-while-its-own-check-stays-unchecked`
   内容＝`SetTrayChecks` 的产码调用者 0 枚；判据建议：挂载之后 `trayMuted` 必须与 `gate.Muted()` 一致（或按 `gate.go:20-21` 的主从关系，
   明确写"门为主、勾为从"并给一枚真机/台件读数）。⚠ 这一枚要回答的正是本票裁定段留给后续的那格（`290-a1` 已命名"第三枚真相源"）。
2. 名册另外两枚见 `10-mutations-and-cites.md` §7（`executed` 契约的一半／装配根与球侧闭包无生产形状用例）。

## N6 · 那句被重写的假话：两处同形假话，逐字读回来判"今天是不是实话"

（凭据全在 `00-anchor-and-static-scales.md` §0.6，这里只给判语与实跑尺；逐字引文不重复抄。）

| 处 | 位置 | 判语 | 凭据（本腿现跑） |
|---|---|---|---|
| (a) `ballGestureWhy` | `cmd/wisp/resident_ball_windows.go:464-467` | **【成立】今天这句话是实话** | 三支有执行者逐支对上：`:322` cancel（`onCancelEsc` 来自 `resident_windows.go:217` 的 `ra.vetoByEsc`）、`:323`/`:324` panel（`hs.requestPanelOpen`）、`:321`/`:325` mute（挂载在 `resident_windows.go:303`）；仍无执行者的三支＝`:319` click／`:320` summon／`:326` pause-wake，与注释里点名的三支（`:451`、`:456`）一一对上。**没有**任何"已做"形：电平、托盘勾、`Muted` 态/灰球一个都没进这句 |
| (b) 启动日志 `"gestures"` | `cmd/wisp/resident_ball_windows.go:390-393` | **【成立】今天这句话是实话** | 它自己那句限定 `this line prints before the attach, because the ball is built first` 可核为真：本行在 `startResidentBall` 内（`:387-393`），装配根建球在 `resident_windows.go:217`、建门在 `:278`、挂载在 `:303` ⇒ **顺序逐字属实**；"turn … **once that leg is assembled**" 是条件式不是完成式 |
| (c) 顺带核的那句（不算重写、但同一族） | `cmd/wisp/resident_audio_windows.go:174` | **【成立】没有越界声称** | 逐字 `…电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）` ⇒ 电平交接与"球会呼吸"被划成两件事，后者明说未声称 |

⛔ 未发现"把还没做的事写成做了"的形：
- "电平已经在球上动"——两处新句都没写；
- "托盘勾会跟随"——没写，且装配根 `resident_windows.go:300-302` **明说没有镜像**；
- "陪聊全双工／`Muted` 态／灰球＋斜杠"——没写；`EvMuteKey` 零发射者本腿复跑（尺：`grep -rn "EvMuteKey" --include=*.go cmd/ internal/ | grep -v _test` → 只命中 `cmd/balldebug/main.go:603` 与 `internal/statemachine` 表自身，见下表），与裁定段"产码发射者只有 balldebug:603"一致。

```
本腿现跑 EvMuteKey 尺（非测试）：
  cmd/balldebug/main.go:603            ← 唯一产码发射者（票面裁定段 :60 那句）
  internal/statemachine/*.go           ← 表自身与事件定义
  cmd/wisp/resident_audio_windows.go、resident_ball_windows.go、resident_windows.go 三枚件内 = 0 命中
rc_evemutekey=0；三枚产码件内 EvMuteKey 计数 = 0（⇒ 本票没碰状态机，与"⛔ 不动状态机"自报一致）
```

## N8 · 禁区与越界（逐笔尺；含 0 字节那枚证据件）

### 8.1 逐笔（⛔ 未用区间 diff 作判据；区间只用来做归因旁证）

```
git diff-tree --no-commit-id -r --name-only --no-renames fdee536b
  cmd/wisp/resident_audio_windows.go / resident_ball_windows.go / resident_mute_290_windows_test.go / resident_windows.go
git diff-tree --no-commit-id -r --name-only --no-renames 69d8c9ef
  cmd/wisp/config_readers_255.go
git diff-tree --no-commit-id -r --name-only --no-renames 75128ae5
  18 枚：票面那一枚（旧长名，见 8.4）＋ probes/290/r1/** 17 枚
禁列 grep（十枚禁区＋frontend/design/testdata/golden/PLAN/specs）三笔各自 hits=0，rc=1（读的是 grep 无命中）
```
逐枚禁区现读（相对 `6c976aa6` 的**区间**尺，只作旁证，因区间内除本票外无人动 Go 面——见下面那条）：
```
git diff --no-renames --name-only 6c976aa6..HEAD -- \
  internal/config/schema.go internal/statemachine/table.go internal/observe/thresholds.go \
  internal/audio/gate.go internal/ball frontend design testdata \
  internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go
→ 空（零字节）
git diff --no-renames --name-only 6c976aa6..HEAD -- internal cmd → **恰好 5 枚**，＝上面两笔产码笔的那 5 枚文件
⇒ 区间里没有别腿的 Go 笔，所以这一发"逐笔 vs 区间"两种尺给同一个答案；实现者 §3 那句"中间两笔只动 .scratch、零 .go"复现为真
```
- 三枚默认值一字未动：`internal/config/schema.go:197` `Enabled bool \`toml:"enabled" default:"false"\``（wake word）／
  `:245` `default:"true"`（voice）／`:277` `MicMutedDefault bool \`toml:"mic_muted_default" default:"true"\``。
  读取处亦未动（`resident_audio_windows.go:256` 仍是 `ra.mutedAtBoot = c.Audio.MicMutedDefault`；票面现量 `:8-9` 引的那枚
  `WithStartMuted(c.Audio.MicMutedDefault)` 今天住在 `:273`——票面写 `:213`，是插入造成的行号漂移，见 §8.5）。
- `frontend/**`／`design/**` 零读零写：三笔文件名册零命中（8.1 第一行尺）；工作树里那 16 枚 ` D design/**` 是别人的在飞件，本腿一枚未碰。
- 票面 AC 框：**零枚被翻**（`^- [ ]`＝4、`^- [x]`＝2，现读见 `00-anchor…` §0.1）；
  `git diff-tree --numstat -r 75128ae5 -- <票面旧长名>` ＝ **`1 0`**（＋1 行／−0 行），
  那一行逐字是 `- [2026-10-09 15:55:59 +08] agent=290-r1 did=照裁定那一版（甲-1）落地 AC#2/AC#4/AC#5…`
  ⇒ **只追加 Progress log 一行，票面原句零改**。这是本票"实现者只许追加 Progress log 行"的最硬一发凭据。

### 8.2 `raw-post-2.md` 是 0 字节——那一格到底交没交

盘上事实：
```
.scratch/wisp/probes/290/r1/raw-post-2.md      = 0 字节
.scratch/wisp/probes/290/r1/rosters/ 里**没有** raw-post-2.red.txt（五枚名册＝pre-1/pre-2/post-1/post-final-1/post-final-2）
```
实现者自己的登记（`30-gates.md:35`）逐字：`| 改后·中间树第二发（被我主动中止，故 0 字节） | raw-post-2.md | — | — | — | 中止 | — | — | — | — |`
⇒ **那一格有同格替代凭据**：票面 `AC#5` 要的是"改前改后**各一次**"，盘上交的是改前 2 发（`raw-pre-1/2.md`＋两枚名册）与
改后终态 2 发（`raw-post-final-1/2.md`＋两枚名册）；`raw-post-2.md` 是**中间树**（`fdee536b` 那一态，cite 还没修）那一发的**中止残件**，
它的同格（中间树改后第一发）由 `raw-post-1.md`（321,015 字节）＋`rosters/raw-post-1.red.txt`（6 枚红，含 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`）**已经交付**。
**判语：那一格不算"没交"**——0 字节件的语义是"我自己叫停的那一发"，不是某个 AC 格的唯一凭据；
本票的判据格（AC#5 的改前改后各一次）另有四件非 0 字节凭据＋五枚名册。
⛔ 本腿未删该件（仓规：临时件只建不删）。
⚠ 具名一条形状欠账（不阻塞裁决）：**0 字节件本身不带任何自证**（"中止"只在叙述里，raw 里没有）。按"证据件须能独立复现"的标准，
这一发应当写成 `raw-post-2.aborted.md` 之类**带一句中止原因的非 0 字节件**。留给编排者当台件纪律的一条，⛔ 不是实现者的 AC 账。

### 8.3 实现者 `30-gates.md` §3 那条尺的口径提醒

它写"禁区命中"时用的正则名册与票面 `AC#4` 的十枚路径**不完全同构**（例如它写了 `docs/` 一整条，票面写的是 `PLAN.md` 与 `docs/specs/**`；
票面还列了 `golden`，它的尺里 `golden` 也在）。本腿**重跑的是票面名册**（见 §8.1），两把尺都零命中 ⇒ 结论一致，口径差异只记录不追。

### 8.4 一处路径事实（不是缺陷，具名防误读）

`75128ae5` 的文件名册里那一枚是**旧长名**票面路径；`5cff604e` 才 `git mv` 成短名。
⇒ 任何拿"短名路径"去 `diff-tree 75128ae5` 的尺都会读到空，误以为"实现者没动票面"或"票面在禁区里"。
本腿两枚名册都跑过（`00-anchor…` §0.8 第 2 条）。

### 8.5 本腿现跑发现的与票面原文不符（票面行号漂移，非实现者造成）

票面"现量"节 `:8-9` 引 `cmd/wisp/resident_audio_windows.go:213` 的那枚 `WithStartMuted(c.Audio.MicMutedDefault)`，
本腿现读＝**`:273`**（`gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`）。
同理票面 `:13-14` 引的 `resident_ball_windows.go:281` `OnMuteHotkey: func() { recordBallGesture("mute-hotkey") }` 今天已被改写、住 `:321`。
⇒ 这些是**票面写于起手锚（`6c976aa6`）之前的行号快照**，落地腿按规矩⛔ 不改票面原句（`AC#4` 也禁止动 `docs/`），所以漂移留在原处——
**归编排者**（要不要在票面追加一节"行号已漂移，现读名册"是票面动作，本腿不动一字）。
⚠ 这条对裁决无影响：`WithStartMuted` 那一枚读取点今天仍是 `:273`，且 `TestAC290MutedDefaultStillDecidesTheBoot` 钉住了 boot 仍在静音位。
