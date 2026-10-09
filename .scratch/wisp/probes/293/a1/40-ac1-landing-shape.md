# 293-a1 · 40 · AC#1 落点形状（⛔ 本格不改产码，只裁"谁跟着谁"要动哪几枚文件）

**锚点**＝`f019f45e`。硬约束尺＝`GOFLAGS= go list -deps ./internal/ball`（**本腿现跑**，⛔ 未引任何旧读数）。

---

## 1. ★硬约束现跑：`internal/ball` 的依赖名册里**不含** `internal/audio`（全称否定，已被尺钉住）

命令（逐字）：`GOFLAGS= go list -deps ./internal/ball` ⇒ `rc=0`，stderr 空。

```
总依赖枚数：145
其中本仓 internal/* 枚数：5
  github.com/CarlosShao/wisp/internal/winsec
  github.com/CarlosShao/wisp/internal/secret
  github.com/CarlosShao/wisp/internal/observe
  github.com/CarlosShao/wisp/internal/statemachine
  github.com/CarlosShao/wisp/internal/ball        ← 自身
internal/audio 命中数：0
```

**正控（这把尺有没有牙）**：同一把尺换射程＝`GOFLAGS= go list -deps ./cmd/wisp` ⇒ `rc=0`，总依赖 **287** 枚，其中 `internal/audio` 命中 **1** 枚、`internal/*` 共 **25** 枚。
⇒ 尺**能**看见 `internal/audio`（在 `cmd/wisp` 名册里就显形了），在 `internal/ball` 名册里显出 0 ⇒ **那句全称否定成立**，且⛔ 不是"尺瞎了"造成的假否定。

第二把尺（源面层交叉核对，不依赖 go list）：`git grep -n "internal/audio" HEAD -- internal/ball` ⇒ **命中 0 行**（整个 `internal/ball` 包没有一枚文件写这个 import 路径）。

⇒ 票面那句"那一跳**必须由 `cmd/wisp` 注入**、⛔ 不许新增 `internal/ball → internal/audio` 的包级依赖边"——**现跑证实**。

---

## 2. 既有形状（注入这一跳仓里已经走过一次，逐枚行号本腿现量）

| 角色 | 位置 | 逐字 |
|---|---|---|
| 后置 setter 的**调用点**（装配根，门存在之后那一行） | `cmd/wisp/resident_windows.go:303` | `rb.attachMuteGate(raudio.toggleMute)` |
| 后置 setter 的**定义** | `cmd/wisp/resident_ball_windows.go:472` | `func (rb *residentBall) attachMuteGate(fn muteGestureFunc) {` |
| 它的锁与"唯一写入者"自陈 | `cmd/wisp/resident_ball_windows.go:469-471`（注释）＋ `:476-477` `rb.muteMux.Lock()/defer Unlock()` | 注释逐字含「**it is the only writer**」 |
| 读回那枚注入物 | `cmd/wisp/resident_ball_windows.go:483` | `func (rb *residentBall) currentMuteGate() muteGestureFunc {`（`:481` 注释「reads the attached executor back under the same lock」） |
| 注入物的**类型**（这一跳今天长什么样） | `cmd/wisp/resident_ball_windows.go:89` | `type muteGestureFunc func() (outcome string, executed bool)` |
| 球句柄在宿主里 | `cmd/wisp/resident_ball_windows.go:134`（`b *ball.Ball`）· `:338`（`rb.b = b`）· 建球失败即早退 `:331-336` | ⛔ **没有** `rb.b == nil` 的构造保证 ⇒ 见 §5 风险 1 |
| 既有的 nil 护栏先例 | `cmd/wisp/resident_ball_windows.go:352`（`if hotReload != nil && rb.b != nil {`）· `:431`（`if rb.b == nil {`） | 同文件已有两枚"先问球在不在"的写法 |

⇒ **"装配根在门装好之后，用一枚后置 setter 把能力注给球宿主"这一形，在 `cmd/wisp` 里是既成事实、已被票 290 落地并有测试面**（`cmd/wisp/resident_mute_290_windows_test.go:93`/`:206`/`:235` 三处直调，票 294 正在补"走查"那一半）。⇒ 票 293 的那一跳**照这一形走即可**，⛔ 不需要新形状。

---

## 3. 那一跳要放哪一侧 ＋ 要动哪几枚文件

**方向（票面已裁、本腿复核凭据）**：门侧 `muted` 为**主**、托盘勾为**从**。
主的凭据（本腿现量，⛔ 不是照抄票面那个偏了的行号）：`internal/audio/gate.go:79` 逐字 `muted    bool`（字段本体）、写 `:168-174`（`func (g *HalfDuplexGate) SetMuted(muted bool)`）、读回 `:202-205`（`func (g *HalfDuplexGate) Muted() bool`）。⚠ 票面引的 `gate.go:20-21` 是**包 doc 注释**（逐字 `//   - Muted (both paths, user intent wins): capture closed; Muted/Unmuted` ＋ `//     events hook the mute hotkey path into the Muted state (ticket 07).`），主张对、位置指偏——已在 `00-...md` §3 具名上报。

### 形状甲（推＝门态一变就调 `SetTrayChecks`）——**动 3 枚文件，全在 `cmd/wisp`，`internal/ball` 0 枚**

| # | 文件 | 要动的地方（现量行号） | 动什么 | 为什么躲不开 |
|---|---|---|---|---|
| 1 | `cmd/wisp/resident_audio_windows.go` | 门句柄已在手里：`:157` `ra.gate`、拧门在 `:162` `gate.SetMuted(!gate.Muted())`、读回已发生在 `:171`/`:173` 那个 `switch` | 把"门此刻的 `muted`"**交出去**：两选一——(i) `toggleMute` 多回一枚 `muted bool`（改 `:156` 的签名 ⇒ 连带改 `resident_ball_windows.go:89` 的 `muteGestureFunc` 类型）；(ii) 新增一枚小读者 `func (ra *residentAudio) mutedNow() bool`（照 `:375 posture()` 的形） | 勾只能读回门，⛔ 不许自存（票面 `AC#2` ＋ 禁区第 2 条） |
| 2 | `cmd/wisp/resident_ball_windows.go` | `muteGesture` 体 `:515-538`；`executed==true` 那一支在 `:535-537` | 在门确实被拧了之后（`executed` 为真那一支）**推一次勾**：`rb.b.SetTrayChecks(<门读回的 muted>, <第二枚见 §5 风险 3>)`，⛔ 必须在 `rb.b != nil` 护栏内（先例 `:352`/`:431`） | 热键与托盘点击**都汇到这一枚函数**（`:321` 与 `:325` 两个闭包都调 `rb.muteGesture`）⇒ 一处落地即同盖"点击拧门"与"热键拧门"两种起点 |
| 3 | `cmd/wisp/resident_windows.go` | 装配根 `:303` 那一行**之后**（该处已在"门存在"之后） | 补**第三种起点＝boot 出厂静音位**那一次镜像：把门的初值推给勾（初值来源现量＝`resident_audio_windows.go:273` `audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`，默认值＝`internal/config/schema.go:277` `default:"true"`） | 球窗与热键**早于采集腿存在**（`resident_ball_windows.go:517-522` 注释逐字"the window and its hot key exist before the capture leg does"）⇒ boot 那一态不可能靠手势带到勾上 |

⇒ **新增包级依赖边＝0 枚**（`cmd/wisp` 早已依赖 `internal/ball` 与 `internal/audio`，尺＝§1 的 `cmd/wisp` 名册 25 枚 `internal/*` 含 audio）。
⇒ **新协程＝0 枚**（D38(b) 名册零膨胀）：`SetTrayChecks` 体只有 `b.sta.PostTask(...)`（`internal/ball/ball_windows.go:954-957`），而 `PostTask` 是往已存在的 STA 窗 `PostMessageW(WM_APPTASK)`（`internal/ball/sta_windows.go:233-249` 逐字注释「Safe from any goroutine」），⛔ 不起线程。
⇒ **无新数据竞争**：手势在 ui-sta 线程内联触发（`resident_ball_windows.go:542-544` 逐字「Callbacks run ON the ui-sta thread」），`showMenu` 也在这条线程读（`ball_windows.go:674`），写入经同一队列 ⇒ `trayMuted` 的读写同侧。⚠ 这一条与票 295（`ra.started` 跨线程裸 bool）是**两枚不同字段**，本形状⛔ 不修 295 那一格。

### 形状乙（拉＝`showMenu` 那一刻现读门态）——**要动 `internal/ball`，多出 1 枚导出 API**

凭据（本腿现量，这是乙形之所以成立、也是它之所以更贵的原因）：
- `internal/ball/ball_windows.go:673-674` ⇒ **每次右键都重建菜单**，不在右键之外缓存勾；
- `internal/ball/tray_windows.go:26-28` ⇒ `type tray struct { data notifyIconData }`，**托盘结构体里没有任何状态位**；
- 于是"勾正确"只在右键那一刻需要成立 ⇒ 只要球侧持有**一枚取门态的函数**，`showMenu` 现调现取，⛔ 不需要任何"门态一变就推"的接线。

代价（为什么本格不推荐它、但⛔ 不由本腿裁）：`internal/ball` 现在**没有**任何"取外部状态的 getter"面（尺＝`internal/ball/ball_windows.go:950-963` 那节只有两枚推式 setter：`SetTrayChecks`／`SetTrayTip`），乙形要**新增一枚导出 API**（形如 `SetTrayCheckSources(func() bool, func() bool)`）⇒ 动的文件是 `internal/ball/ball_windows.go` ＋ `cmd/wisp` 那两枚，且**扩大 `internal/ball` 的公开面**（该包自陈「无业务逻辑」，`docs/specs/SPEC-01-architecture.md:48` 逐字「ball/ # 悬浮球：… 托盘（**无业务逻辑**）」——一枚存着的"取外部静音态的函数"是否越那条线，**这是编排者的裁量，不是本腿的结论**）。

⇒ **两族形状的公共事实**：⛔ 都**不**新增 `internal/ball → internal/audio` 的包级边（乙形注入的是函数，不是包）；⛔ 都不新造协程；★ **甲形动 3 枚文件、全在 `cmd/wisp`、`internal/ball` 0 枚改动**（因为它复用既有的 `SetTrayChecks`）；乙形动 `internal/ball` 且要新导出 API。票面 `AC#2` 的判据文字（"门态一变，勾立刻镜像"＋三种起点各要对）**本身就是甲形的措辞**——这一句只作指认，⛔ 不代替编排者裁。

---

## 4. 落地腿的写面预告（供 `AC#3` 越界检查比对，⛔ 本腿没动它们）

按 §3 甲形，`git diff-tree -r --name-only --no-renames` 应逐笔只出现：
```
cmd/wisp/resident_audio_windows.go
cmd/wisp/resident_ball_windows.go
cmd/wisp/resident_windows.go
cmd/wisp/<新测试件>.go            ← AC#2 的用例面（票 294 与本票同写 cmd/wisp 测试面 ⇒ 串行）
（可选）internal/ball/ball_windows.go   ← ⛔ 只有走乙形才会出现；甲形下它出现即为越界信号
```
票面十枚禁列（`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`ticket90_persist_test.go`／`testdata/golden`）⛔ 任一命中即退回——**本腿已核：甲形三枚文件与十枚禁列零交集**（尺＝上面那三枚路径逐枚对 `git diff-tree` 的禁列前缀，无匹配）。

⛔ **禁区复核**：`internal/config/schema.go:197`（`wake_word.enabled`）／`:245`／`:277`（`mic_muted_default`）三枚默认值本形状**一枚都不必动**——boot 那一次镜像只是**读** `:277` 的既有默认，⛔ 不写、不改语义。

---

## 5. 三枚编排者下刀前必须知道的风险（本腿量到的，⛔ 不属本腿裁）

1. **`rb.b` 可能是 nil**。`resident_ball_windows.go:331-336`：建球失败时 `rb.verdict` 记"本进程没有浮球窗"后**早退**，`rb.b = b` 在 `:338` 才发生 ⇒ 装配根 `resident_windows.go:303` 之后任何一次镜像调用都**必须**先问 `rb.b != nil`（同文件先例 `:352`/`:431`）。⛔ 漏掉这一条＝无桌面会话（服务会话/CI runner）一 boot 就 panic，而那正是 `:53-57` 那段注释承诺要撑住的形状。
2. **那段"显式不镜像"的注释会被这一发改写**。`cmd/wisp/resident_windows.go:300-302` 逐字（本腿现读）：
   ```
   // and no mirror of the gate into the tray's own checkmark - the
   // ball-side trayMuted display flag stays undriven here rather than becoming a
   // second authority nobody reconciles.
   ```
   同族立场还在 `resident_ball_windows.go:40-51`（逐字 "this host keeps no copy of it … and no second 'am I muted' flag"）。⇒ 票 293 落地**必然**要改这两段注释文字。⛔ 这不是"顺手改注释"：它得同时说清**为什么现在不构成第二权威**（凭据现成＝勾是 `gate.Muted()` 的现读、不是新存的一份），否则改完的注释就是一句新的谎。
3. **`SetTrayChecks` 的第二枚实参今天没有真相源可喂**。见 `10-ac0a-pause-wake-driver.md` §4：`trayPauseWake` 该映的那个态（唤醒监听被暂停）在 `internal/speech` 里，而那枚包只有 `doc.go`；`schema.go:197` `wake_word.enabled` 出厂 `false`。⇒ 编排者要裁"只喂 `muted`、`pausedWake` 显式留 `false`（配一句指名理由）"还是别的。⛔ 本腿不挑。

---

## 6. 本节用过的尺（逐字，可复跑）

```
GOFLAGS= go list -deps ./internal/ball                    # rc=0；145 枚；internal/* 5 枚；internal/audio 命中 0
GOFLAGS= go list -deps ./cmd/wisp                          # 正控：rc=0；287 枚；internal/* 25 枚；internal/audio 命中 1
git grep -n "internal/audio" HEAD -- internal/ball         # 0 命中（源面层交叉核对）
go env GOOS GOARCH GOFLAGS GOVERSION                       # windows amd64 (GOFLAGS 空) go1.27.1
git grep -n "attachMuteGate\|currentMuteGate\|muteGestureFunc" HEAD -- cmd internal
sed -n '50,60p;295,320p' cmd/wisp/resident_windows.go
sed -n '150,205p;240,300p;370,400p' cmd/wisp/resident_audio_windows.go
sed -n '40,60p;130,150p;300,360p;445,490p;500,548p' cmd/wisp/resident_ball_windows.go
sed -n '118,135p;655,690p;940,965p' internal/ball/ball_windows.go
sed -n '1,60p;69,105p' internal/ball/tray_windows.go
sed -n '225,250p' internal/ball/sta_windows.go
sed -n '14,30p;55,80p;160,210p' internal/audio/gate.go
sed -n '193,206p;272,282p' internal/config/schema.go
```

⛔ 未跑：`go build`／`go vet`／`go test`（按派单禁令让给在飞的 `292-v1`）；`gofumpt`／`d22scan`／真机名册（属票面 `AC#4`，落地腿的格）。
