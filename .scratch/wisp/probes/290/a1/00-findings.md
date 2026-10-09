# 票 290 AC#0 + AC#1 — 只读普查件（腿 `290-a1`）

**生成腿**：`290-a1`（只读普查＝AC#0 契约普查 + AC#1 三形代价表）
**锚点**：分支 `dev`，起点 HEAD `34a012599118516ae41cc56ff9453216df875779`（`2026-10-09T14:03:23+08:00`）
**性质**：⛔ **本格零产码改动、零默认值改动、零 AC 翻框**。三形全部只写代价、不落地。
**纪律自报**：
- ⛔ **我没有跑 `go build` / `go vet` / `go test` 任何一枚**（Go 编译/测试面此刻由写腿 `289-r1` 独占）。跑过的取数只有：`git grep -n <HEAD>`、`git show HEAD:<path>`、`git ls-tree`、`git cat-file -e`、`git status`、`GOFLAGS= go list -deps`（只出不入）、`date`。
- ⛔ 零读零写 `frontend/**` 与 `design/**`；本件不引它们任何结论。
- 文件存在性一律只认 git 对象层；行号锚全部现跑（每把尺带命令原文，可重跑）。
- 写点唯一＝本件；票面只在 `## Progress log` 末尾追加一行，⛔ 未改票面任何原句、⛔ 未勾任何 AC。

**脏度前置尺（决定本件行号可信度的那一把）**：
`git status --porcelain -- internal/audio/gate.go cmd/wisp/resident_audio_windows.go cmd/wisp/resident_ball_windows.go internal/ball/ball_windows.go internal/statemachine/table.go internal/statemachine/events.go internal/config/schema.go internal/config/manager.go docs/PLAN.md docs/specs` → **输出为空** ⇒ 这些文件工作树 == HEAD，本件引用的行号既是工作树行号也是 HEAD 行号。

**名册尺（对象层）**：`git ls-tree --name-only HEAD docs/specs/` → 13 份 SPEC 全在，票面点名的三枚实名＝`SPEC-04-voice-pipeline.md` / `SPEC-06-security-gatekeeping.md` / `SPEC-08-ui-ball-panel.md`（票面写的就是这三个名，无笔误）。
`git ls-tree -r --name-only HEAD internal/speech` → 只有 `internal/speech/doc.go` 一枚（票面 AC#1 丙的这块料自证成立）。

---

## 0. 我对票面现量的复量（尺与读数，先自证再引用）

尺＝`git grep -n "SetMuted(\|SetSpeaking(" HEAD -- internal cmd`（14:0x 现跑），逐枚点名：

| # | 命中 | 属类 |
|---|---|---|
| 1 | `internal/audio/gate.go:97` `// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it` | 注释 |
| 2 | `internal/audio/gate.go:130` `func (g *HalfDuplexGate) SetSpeaking(speaking bool) {` | 定义 |
| 3 | `internal/audio/gate.go:168` `func (g *HalfDuplexGate) SetMuted(muted bool) {` | 定义 |
| 4-10 | `internal/audio/gate_test.go:58/:68/:96/:99/:140/:143/:144/:147/:184/:204` | 测试（`SetSpeaking` 七枚：58/68/96/99/143/144/204；`SetMuted` 三枚：140/147/184） |

- 分尺读数：`git grep -n "SetMuted(" HEAD -- internal cmd \| grep -v _test` → **2 枚**（gate.go:168 定义 + :97 注释）；`grep -c _test` → **3 枚**（gate_test.go:140/:147/:184）。
- ⇒ **产码调用者＝0 枚**，票面"全仓没有任何一条生产路径能把门拧开"**成立**。
- ⚠ **与票面的一处数字出入见 §4.1**（票面 `:12` 说 `SetMuted` 的测试命中"四枚"却只列了三枚行号；现量＝三枚）。

其余票面现量逐枚复核**全部成立**（行号一字不差）：`resident_audio_windows.go:213` 门装配、`:238` 默认档文案、`:260` "gate 未 open"、`internal/audio/gate.go:97` 注释、`internal/config/schema.go:277` `default:"true"`、`resident_ball_windows.go:281` `OnMuteHotkey: func() { recordBallGesture("mute-hotkey") },`、`internal/ball/ball_windows.go` 的 `Events` 十枚字段（`:49-60`，字段本体 `:50-59` 恰好十枚）、`case wmHotkey:` `:656`、`case hkMute:` `:660`、`case menuMute:` `:678`、`internal/statemachine/table.go:65/:75/:79`。

### 0.1 我这把尺新量出来的、票面没说的三枚事实（都影响 AC#1）

1. **门早就备好了"把 Muted 态接过去"的那枚钩子，而且它自己也没被接**：`internal/audio/gate.go:61-65` 逐字
   `// WithGateEvents registers the event sink (mute hotkey path -> Muted state;` ／ `// the ball/state machine subscribes at wiring time, ticket 07/16).`
   尺＝`git grep -n "WithGateEvents" HEAD -- internal cmd` → 命中只有定义 `gate.go:61/:63` 与测试 `gate_test.go:86/:127` ⇒ **产码调用者＝0**。
   同族：`gate.go:12` 逐字 `// Half-duplex gate (D16): the seam the TTS side and the mute hotkey drive.`、`gate.go:165-167` 逐字 `// SetMuted is driven by the mute hotkey. Muted closes capture on both paths` / `// (user intent outranks the duplex policy); the Muted/Unmuted events are the` / `// hook into the Muted ball state (07).`
   ⇒ **"谁是谁的真相源"这件事，产码注释已经替票面答了**：门侧 `muted` 是主、球侧 `Muted` 态是订阅者（投影）。见 §2.1。
2. **`EvMuteKey` 在产码里零发射者，而且它的入态 `Armed` 也零发射者 ⇒ 球侧今天连进 `Muted` 这条路都不存在**：
   尺＝`git grep -n "EvMuteKey" HEAD`（去测试/票面/台账）→ 产码命中只有 `cmd/balldebug/main.go:603`（dev 调试件）与定义 `internal/statemachine/events.go:21`。
   尺＝`git grep -n "EvKwsEnabled" HEAD -- internal cmd | grep -v _test` → 只有 `events.go:18` + `table.go:55` 两枚定义行，**零发射者**。
   尺＝`git grep -n "Dispatch(" HEAD -- cmd internal | grep -v _test | grep -v balldebug` → 产码状态机发射只有 `internal/models/bridge.go:40/:46/:59/:64` 那四枚（`EvModelMissing`/`EvDownloadFailed`/`EvDownloadCompleted`），**没有一枚是 mute/kws**。
   ⇒ 后果：`table.go:65`（#8 `Armed`→`Muted`）是 `Muted` 的唯一入边，而 `Armed` 的入边 #5（`table.go:55`，`EvKwsEnabled`）零生产者 ⇒ **甲形若要求"球先进 `Muted` 态"，今天必撞 D43 无合法边**（`machine.go` 判非法转移），那正是票面 AC#4 说的"要先落一枚人工批准的 `A##`"那一格。
3. **球侧还有第三枚静音旗，今天没人写它**：`internal/ball/ball_windows.go:128` `trayMuted bool`（托盘勾号位），写点只有 `SetTrayChecks`（`:952-958`）——尺＝`git grep -n "SetTrayChecks" HEAD -- cmd internal | grep -v _test` → 命中只有注释 `:952` 与定义 `:953` ⇒ **产码调用者＝0**。而 `:674` `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` 每开一次托盘都读它。
   ⇒ 真相源风险不止票面说的那两枚（门侧布尔 / D43 态），托盘勾号是**第三枚**，代价表里必须点名（见 §2.2）。

### 0.2 `recordBallGesture` 到底做了什么（票面说"记一笔就断"，现量到句）

`cmd/wisp/resident_ball_windows.go:412-420` 逐字：
```
412: // recordBallGesture books one ball gesture that arrived with nowhere to go.
417: func recordBallGesture(name string) {
418: 	slog.Warn("ball gesture arrived with no executor", "gesture", name, "why", ballGestureWhy)
419: 	fmt.Printf("wisp: ball %s: %s\n", name, ballGestureWhy)
420: }
```
十枚回调的填法（`:278-289`）：`OnClickBall`/`OnSummonHotkey`/`OnMuteHotkey`/`OnTrayMute`/`OnTrayPauseWake` 五枚走 `recordBallGesture`；`OnCancelHotkey` 走 `recordCancelHotkey(onCancelEsc)`（`:282`，同族里唯一有执行者那枚，`onCancelEsc` 是装配根注入的 `escVetoFunc`，定义 `:63`、注释 `:237`）；`OnPanelHotkey`/`OnTrayPanel` 走 `hs.requestPanelOpen(...)`；`OnTrayExit` 走 `recordTrayExit`；`OnDragEnd` 走 `recordBallDragEnd()`。
⚠ 那枚 `ballGestureWhy` 常量（`:408-410`）逐字：`"this process has no task pipeline and no microphone, so the gesture has no executor here; ..."` —— **"no microphone" 这半句自票 247 起在本进程里是假的**（同一进程 `cmd/wisp/resident_windows.go:278` 就装配了采集腿）。甲形必须改这句，见 §4.3。

---

## 1. AC#0 契约普查：「麦克风的开／关由谁发起」

### 1.1 搜过的字样名册（尺与真实命中数，不是"我搜过了"）

尺＝`git grep -n -i -e "静音" -e "mute" -e "开关" -e "toggle" -e "按住" -e "push-to-talk" -e "ptt" -e "授权" -e "privacy" -e "隐私" -e "唤醒" -e "wake" -e "麦克" -e "microphone" -e "mic" HEAD -- docs/PLAN.md docs/specs/SPEC-04-voice-pipeline.md docs/specs/SPEC-06-security-gatekeeping.md docs/specs/SPEC-08-ui-ball-panel.md`
→ 原始命中 259 行（落 `/tmp/p290/all.hits`，仓库外）。逐字样计数（大小写不敏感 ci / 敏感 cs）：

| 字样 | ci | cs | 字样 | ci | cs |
|---|---|---|---|---|---|
| 静音 | 12 | 12 | privacy | 5 | 5 |
| mute | 11 | 2 | 隐私 | 44 | 44 |
| Muted | 9 | 8 | 唤醒 | 37 | 37 |
| 开关 | 10 | 10 | wake | 10 | 1 |
| toggle | **0** | **0** | 麦克 | 44 | 44 |
| 按住 | **0** | **0** | microphone | 2 | 1 |
| push-to-talk / PTT | **0** | **0** | mic | 27 | 21 |
| 授权 | 82 | 82 | | | |

分文件命中：`PLAN.md` 203 ／ `SPEC-04` 18 ／ `SPEC-06` 16 ／ `SPEC-08` 22。
⇒ **字样名册里"按住／push-to-talk／PTT／toggle"四支在全套契约文本里零命中**——契约从没许过 push-to-talk，也几乎没用"toggle"这个词（中文侧写的是"一键…静音"与"再按静音键"）。

### 1.2 逐字引文（每条带 `文件:行`）

**（一）"由谁发起"＝静音热键，契约写了，而且写了两次（进与出）**

1. `docs/PLAN.md:462`（D16② 隐私三件套之内）逐字
   `  ② KWS 激活时悬浮球有常亮可见指示 + 一键静音全局快捷键`
2. `docs/PLAN.md:1312`（D16 那节的态表）逐字
   `| \`Muted\` | 灰球 + 斜杠 | 一键静音（D16） | 再按静音键 |`
   ⇒ **这一枚是 AC#0 最硬的一条**：它同时定了"发起者＝那枚一键静音热键"和"出去也靠同一枚热键再按一次"。
3. `docs/PLAN.md:3068` 逐字 `| 8 | \`Armed\` | 静音键 | \`Muted\` | 停 KWS 推理（模型可保留） |`
   `docs/PLAN.md:3070` 逐字 `| 10 | \`Muted\` | 静音键 | \`Armed\`（KWS 开）/ \`Sleeping\` | — |`
4. `docs/specs/SPEC-08-ui-ball-panel.md:96` 逐字 `| 8 | \`Armed\` | 静音键 | \`Muted\` | 停 KWS 推理（模型可保留） |`
   `docs/specs/SPEC-08-ui-ball-panel.md:98` 逐字 `| 10 | \`Muted\` | 静音键 | \`Armed\`(KWS 开)/\`Sleeping\` | — |`
   产码那三枚对应行（现量逐字）：`internal/statemachine/table.go:65` `D43: 8, From: StateArmed, Event: EvMuteKey, To: StateMuted,`／`:66` `SideEffects: []string{"kws.stop-inference"},`／`:75` 与 `:79` 两枚 `D43: 10, From: StateMuted, Event: EvMuteKey, To: StateArmed,` / `To: StateSleeping,`（守卫分别是 `kws-loaded` / `kws-not-loaded`）。
5. `docs/specs/SPEC-04-voice-pipeline.md:47` 逐字 `- 静音键：一键全局静音 → \`Muted\` 态；KWS 激活时悬浮球常亮指示（D16②）。`
   `docs/specs/SPEC-04-voice-pipeline.md:80` 逐字 `2. KWS 激活常亮指示 + 一键静音快捷键。`（票面 :22-:24 已引，我复量行号一字不差）
6. `docs/specs/SPEC-08-ui-ball-panel.md:229` 逐字 `- 托盘：右键菜单（打开面板/静音/暂停唤醒/退出）；左键 = 打开面板。`
   ⇒ **第二个发起者（托盘同名菜单项）契约也写了**，且 `internal/ball/ball_windows.go:678-679` 的 `case menuMute: → b.fire(...OnTrayMute)` 已把这条路由到产码（只是终点仍是 `recordBallGesture`）。

**（二）"默认不开麦"＝隐私侧的硬默认，契约写了三处**

7. `docs/PLAN.md:460` 逐字 `  ① 默认不打开麦克风（D2 已定 KWS 为 opt-in）`
8. `docs/specs/SPEC-04-voice-pipeline.md:79` 逐字 `1. 默认不打开麦克风（KWS opt-in）——可写进 README 第一行的强论证。`
9. `docs/PLAN.md:1999` 逐字片段 `② 默认不碰麦克风（D16）③ 本地优先、零云依赖可跑 ④ 插件不常驻**。`

**（三）另有一条"用户开麦"的形状被定死了，且它自带确认门（本票 AC#0 的第三发现）**

10. `docs/specs/SPEC-08-ui-ball-panel.md:118` 逐字 `| 30 | \`Warm\` | 开启 Conversation | \`Conversation\`→\`Listening\` | **首次开启需 L2 级隐私确认** |`
    `docs/specs/SPEC-08-ui-ball-panel.md:116` 逐字 `| 28 | \`Speaking\` | 播报完（Conversation 模式） | \`Listening\` | 麦克风开启 + 红色常亮环 |`
    `docs/specs/SPEC-08-ui-ball-panel.md:71` 逐字 `| \`Conversation\` | **danger 常亮环，不得渐隐、不得弱于 Confirming**（麦克风开着的唯一交代） |`
    与 `docs/PLAN.md:2026`（§16 定案表第 5 行，逐字片段）`**\`Conversation\` 保留为显式开启项**（\`[voice] conversation_mode = false\` 为默认值）… **首次开启需 L2 级隐私确认并记日志**`
    ⇒ **这是全仓唯一一处"契约明写了一个把麦克风打开的用户动作"并且给它的代价定了形**：显式开启项落在 config（`conversation_mode=false`）＋**首次开启一枚 L2 级隐私确认**＋常亮环交代。本票甲/乙两形与它同族，编排者裁形时应把这一枚当参照。
11. 唤醒词那支：`docs/specs/SPEC-08-ui-ball-panel.md:95` 逐字 `| 7 | \`Armed\` | 唤醒词命中 | \`Listening\` | 同 #4；KWS 暂停 |`；`docs/specs/SPEC-04-voice-pipeline.md:70` 逐字 `- 唤醒词命中 → \`Listening\`，**KWS 暂停**（半双工内不需要同时跑）`。
    ⚠ 注意射程：契约的唤醒词命中是 `Armed`→`Listening`（**前提是 KWS 已加载**），它**不是**"开麦"那一跳的发起者；"开麦（KWS）"那一跳的发起者是 #5 `Sleeping`→`Armed`（`EvKwsEnabled`），而**契约文本里我找不到一枚"用户怎么把 KWS 打开"的动作句**（只有 `PLAN.md:1884` 把"KWS 开关"列进不可热加载段）。这一格空白见 §1.3 ⓒ。

**（四）32 条冻结契约（C 表）里与音频／权限相关的行（现读 `docs/PLAN.md:1351-1382`）**

12. `docs/PLAN.md:1358` 逐字 `| C8 | \`AudioSource\` | 真实麦克风 / **wav 注入（测试钩子）** / \`[Aec]\` 仅接口位 | D16, 测试地基 |`
13. `docs/PLAN.md:1359` 逐字（C9）`| C9 | \`AsrEngine\` \`TtsEngine\` \`WakeWordEngine\` | 本地 sherpa 与云端 provider 可互换；**D47 增 \`RealtimeEngine\`（见 C32）**：云端 S2S realtime，opt-in、零工具权限 | D5, **D47** |`
14. `docs/PLAN.md:1362` 逐字 `| C12 | \`BallState\` | **20 态** + **完整转移表（D43，第四轮才真正写出）** + 每态超时 | §2 + D43（⚠ 原标「来源 D24」为误引） |`
    ⇒ **C12 就是 D43 那张表**：甲形若新增态或改转移边，动的就是冻结契约（票面 AC#4 同一判语）。
15. `docs/PLAN.md:1367`（C17 `PanelBridge`：前端↔Go 双向通道，回复按 correlationId 路由，前端必须无状态）、`docs/PLAN.md:1368`（C18 `ApprovalQueue`，逐字末句 `**「允许」决策只接受原生侧来源（D33/F2，面板来源直接拒绝）**`）⇒ 乙形的两条硬边界。
16. `docs/PLAN.md:1382`（C32 `RealtimeEngine`，逐字片段 `opt-in + 显式「音频将上传」告知 + 首次 L2 隐私确认。`）⇒ 与第 10 条同形参照。
17. `docs/PLAN.md:1390` 逐字 `**契约变更流程**：任何修改必须人批准，且同步更新 \`docs/DECISIONS.md\` 与受影响切片卡。`（⚠ 该句指向的 `docs/DECISIONS.md` 按 AGENTS.md 索引段与对象层名册目前不存在，不属本票射程，仅记一句）

**（五）`SPEC-06` 里与"开麦是不是授权"有关的行（结论：SPEC-06 全篇讲的是工具风险/路径/外泄，**没有任何一枚句子管麦克风**）**

18. 尺＝`grep -i -E "静音|mute|麦|mic" /tmp/p290/all.hits` 过滤 `:docs/specs/SPEC-06` → **零命中**。SPEC-06 的 16 枚命中全部落在"授权目录 / 会话授权 / L2 / 面板来源拒绝"那族（`SPEC-06:35`、`:37`、`:113`-`:122`、`:131-132`）。
    ⇒ 票面把 SPEC-06 列入 AC#0 射程是对的（要证"它没写"），**SPEC-06 对"谁发起开麦"一句都没写**。

### 1.3 三问结论

**ⓐ 契约有没有已经定死的某一形？——有，而且是两枚，不是零枚。**
定死的是：**(主发起者) 一键全局静音热键**（最硬一枚＝`docs/PLAN.md:1312` 那行同时写了进与出："一键静音（D16）／再按静音键"；佐证 `PLAN.md:462`、`PLAN.md:3068`/`:3070`、`SPEC-04:47`/`:80`）**＋ (次发起者) 托盘右键菜单的"静音"项**（`SPEC-08:229`）。
定死的**不是**面板开关、**不是**按住说话（`按住`/`push-to-talk`/`PTT` 全套契约零命中）、**不是**唤醒词（唤醒词命中在 D43 里是 `Armed`→`Listening`，前提是麦已经开）。
⚠ 但**"那一跳的语义射程"是空白**：契约给 #8 挂的副作用逐字只有 `停 KWS 推理（模型可保留）`（`PLAN.md:3068`/`SPEC-08:96`/`table.go:66` 的 `kws.stop-inference`）——**契约从没写"静音键要关/开采集设备"**。门侧那枚 `muted` 恰恰管的是设备交接（`gate.go:96-98` 注释＋`gate.go:234-242` `effectiveOpen()`：`muted` 一票否决）。
⇒ **结论必须说全**：发起者＝已定死（热键，托盘同名项次之）；**"热键是否顺手拧门（开/关设备）"＝契约没写**。后半句不是实现决策，是功能语义决策 ⇒ **摆机主**（⛔ 本写码腿不替它挑）。票面 AC#0 ⓒ 那一支的判据在此**部分触发**，不是我漏搜。

**ⓑ 有没有关于持久化的许诺？——一半有一半没有；且契约的 `hot` 档与产码自述互相打脸。**
有的那一半：
- `docs/PLAN.md:2735`（D36 三档生效级表）逐字 `| \`[audio]\` | \`input_device\` \`sample_rate\` \`half_duplex\`(硬编码 true，只读显示) \`mic_muted_default\` | \`hot\`（须重开设备） |` ⇒ 配置侧那枚键归 **`hot` 档（须重开设备）**，不是 `reload`、不是 `restart`。
- `docs/PLAN.md:2757` 逐字 `3. **默认值的唯一真相源是 Go 结构体 tag**，TOML 只是它的序列化。` ＋ `docs/PLAN.md:2759-2760` 逐字 `4. **GUI 与 TOML 的双向一致由"读写同一个结构体"保证**…  GUI 不得持有任何 TOML 里没有的状态。` ＋ `docs/PLAN.md:1886-1887` 逐字 `- **单一真相源始终是文件**（D6），不得出现 GUI 内存态与文件态不一致` ⇒ 若"开关"做成面板/配置态，它**必须落 `config.toml`**，这是硬规矩。
- `internal/config/schema.go:276-277` 逐字 `// MicMutedDefault starts every session muted.` / `MicMutedDefault bool \`toml:"mic_muted_default" default:"true"\`` ⇒ 产码把语义写死成"**每次会话启动都从配置回到静音**"，运行时选择不进 TOML。
- `docs/PLAN.md:1884` 逐字 `（provider、插件开关、阈值、快捷键）；**不可热加载的段**（KWS 开关、平台相关）` ＋ `:1885` `必须明确提示「需重启生效」，**不得静默忽略**` ⇒ **KWS 开关在契约正文里被归为"不可热加载"**，而 `internal/config/tiers.go:29` 却把整段 `"audio": "hot"` 登记为热档（两枚不是同一个键，但两句话读起来会打脸，落地腿必须点名分清）。
没有的那一半（尺＝`git grep -n -E "重启|持久|落库|热加载|热重载|restart|persist" HEAD -- <四份契约>` 再过滤 `静音|mute|麦|mic|KWS|唤醒|wake|状态机|ball|悬浮球`，14:0x 现跑）：
> 命中只有四行且**没有一枚管静音/开麦的持久化**：`PLAN.md:549`（一键重启提示）、`PLAN.md:1884`（上面那条）、`PLAN.md:1972`（子进程崩溃只带走语音能力）、`SPEC-08:125`（#35 看门狗一键重启）。
⇒ **"用户主动开门之后重启，是回到静音还是记住开着"——契约一句都没写**；`"每次显式"`这个形状也没有许诺过。这一问＝**摆机主**（功能决策）。
再一条现量硬事实（落地腿必须知道）：即便把 `mic_muted_default` 改成 `false` 写回 TOML，**本次运行也不会生效**——`cmd/wisp/config_readers_255.go:150` 逐字 `… it reads the file once at boot, so a hot reload of this section waits for a restart instead of acting in this run`，而 `cmd/wisp/resident_audio_windows.go:196` 那次读取是 `config.LoadFile` 新读、不是活 Manager 的重读。

**ⓒ 契约一句都没写吗？——不是"一句都没写"，但确实有一格空白，且那一格是本票的要害。**
写了：发起者（热键＋托盘）、默认不开麦（三处）、D43 #8/#10 两条边、`mic_muted_default` 的 `hot` 档与"每次会话从配置起"。
没写（⇒ 必须摆机主，⛔ 本腿不挑）：
1. **静音键是否兼作"设备开关"**（契约只写 `停 KWS 推理`）；
2. **开门后的持久化语义**（重启回静音 vs 记住开着）；
3. **"把 KWS/麦打开"这个动作由谁发起**（`EvKwsEnabled` 那枚入态事件在四份契约里找不到用户动作句；唯一有用户动作句的开麦形状是 Conversation 那支，而它要 L2 隐私确认）。
⇒ 依票面 AC#0 ⓒ 与"排程"段那半句（"若 AC#0 判'契约没写'则先摆机主"），本票**不是整格空白**而是**三处局部空白**：甲形（热键接门）有契约文本撑腰（发起者已定死）但**语义射程不足**，⛔ 不许由写码腿自行把"停 KWS 推理"读成"关设备"。

---

## 2. AC#1 三形代价表（⛔ 不落地、⛔ 不改默认值）

先给共同底座（三形都建立在它上面，全部现量）：

| 事实 | 出处（逐字/尺） |
|---|---|
| 门的两枚唯一开门手段 | `gate.go:130 SetSpeaking` / `gate.go:168 SetMuted`；产码调用者各 **0**（§0 尺） |
| `muted` 一票否决，`speaking` 只在 Path T 关门 | `gate.go:234-242` `effectiveOpen()`：`if g.muted { return false }` / `if g.speaking && g.path == PathT { return false }` |
| 静音位上设备**根本不打开** | `gate.go:96-98` `// Start implements AudioSource. When the gate is (currently) closed the inner / source is NOT started; SetMuted(false)/SetSpeaking(false) will start it / later with the stored ctx/buf, so Start still returns nil.` |
| 门自己声明的钩子 | `gate.go:61-65 WithGateEvents`（产码调用者 0）；`gate.go:165-167`（`SetMuted is driven by the mute hotkey… the Muted/Unmuted events are the hook into the Muted ball state (07)`） |
| 球侧态由注入写，不由状态机写 | `cmd/wisp/resident_approval_windows.go:957` `b.SetState(statemachine.StateSleeping)`；`cmd/wisp/resident_ball_windows.go:275` `Initial: statemachine.StateSleeping`；产码 `Dispatch` 只有 `internal/models/bridge.go` 那四枚 |
| 电平出腿（门的下游） | `resident_audio_windows.go:112 levelOut` → `:132 rb.b.SetAudioLevel(level)`；`resident_windows.go:278` `raudio := startResidentAudio(rt, rb.setAudioLevel)` |
| 装配根顺序（关键代价） | `resident_windows.go:217` 先 `rb := startResidentBall(rt.Registry, ra.vetoByEsc, …)`，`:278` 才 `raudio := startResidentAudio(...)` ⇒ **注入是单向的（音频腿吃球的 setter）**，反向（球腿吃门的 setter）今天没有槽 |
| 停机 | `resident_audio_windows.go:289-296` `func (ra *residentAudio) stop() … err := ra.gate.Stop()`；门侧 `Stop` 会把 `running/open` 归零（`gate.go:115-126`） |
| 三枚默认值（⛔ 一字不动） | `schema.go:245 VoiceSection.Enabled default:"true"`、`schema.go:197 WakeWord.Enabled default:"false"`、`schema.go:277 MicMutedDefault default:"true"`；`HalfDuplex` `:275` 硬编码 true |

### 2.1 甲形：静音热键（＋托盘同名项）那支接到门上

- **要动的文件＝3 枚，全在 `cmd/wisp`（同包）**：
  1. `cmd/wisp/resident_audio_windows.go` — 给 `*residentAudio` 加一枚翻转方法（内部调 `ra.gate.SetMuted(!ra.gate.Muted())`）；同时 `:234`/`:238`/`:259` 那三段 verdict 分支与 `:315 posture()`、`:94 mutedAtBoot` 的口径必须在翻转后重打，否则 `:238` 那句"设备未打开"会在开门后继续说谎。
  2. `cmd/wisp/resident_ball_windows.go` — `startResidentBall` 签名（`:256`）多一枚注入（先例逐字＝同位的 `onCancelEsc escVetoFunc`，`:63` 定义、`:237` 注释、`:282` 使用）；`:281 OnMuteHotkey` 与 `:285 OnTrayMute` 两支从 `recordBallGesture` 改指向注入的执行者；`:408-410 ballGestureWhy` 那句必须重写（现含"no microphone"，票 247 起为假，见 §4.3）。
  3. `cmd/wisp/resident_windows.go` — 装配根解顺序（`:217` 早于 `:278`）：需要一个后置 setter 或一个指向"稍后才出现的 ra"的闭包壳（`rb` 已有 `cancelHosted` 那类"有没有执行者"的自我申报形状，`:257`）。
  - ⛔ **`internal/audio/gate.go` 不必改**：`SetMuted`／`Muted`／`Open`／`WithGateEvents` 全已存在（`:168/:202/:209/:63`）。若采用"门→球态投影"，改的是 `resident_audio_windows.go:213` 那枚 `NewHalfDuplexGate(...)` 加一个 `WithGateEvents(...)` option，仍是 `cmd/wisp` 内。
- **谁跟谁（真相源）**：**门侧 `muted` 布尔为主，球侧 D43 `Muted` 态为从（投影）**。不是我推的，是产码注释已写死的：`gate.go:20-21`（`Muted (both paths, user intent wins): capture closed; Muted/Unmuted events hook the mute hotkey path into the Muted state`）＋ `gate.go:61-62`（`the ball/state machine subscribes at wiring time`）＋ `gate.go:166-167`。
  ⚠ **两枚必须点名的次生真相源**：(i) 托盘勾号 `ball_windows.go:128 trayMuted`（写者 `SetTrayChecks` 现零调用者）——若门为主，勾号必须由门事件驱动，⛔ 不许再来一枚布尔；(ii) `internal/tools/subagent_197.go:109` 把 `subagentStateStopped = statemachine.StateMuted` 复用了一枚 `StateMuted` 常量（尺 7），与球态无关，落地腿别把它读成第二真相源，但它会让 `StateMuted` 的读者混淆。
  ⚠ **硬代价（票面 AC#1 那一问的正答案）**：若坚持"球也要真进 `Muted` 态"，今天**无路可走**——`Muted` 唯一入边 `table.go:65`（#8 需先 `Armed`），`Armed` 入边 `table.go:55`（#5 `EvKwsEnabled`）**零发射者**（§0.1 第 2 条）。所以甲形只有两种诚实形状：**(甲-1)** 球态仍由注入直接 `b.SetState(StateMuted)`（绕过转移表＝本仓裁过两次的形状，票面 246 同族，⛔ 我不推荐、只点名）；**(甲-2)** 只动门、球态保持 `Sleeping`、视觉上另有交代（则 D43 #8/#10 那两枚边继续无人走，"Muted 态"契约项继续悬空）；**(甲-3)** 真要态跟随 ⇒ 需 KWS/`Armed` 链路或新增转移边 ⇒ **那是人工批准的 `A##`，不是本票**（票面 AC#4 逐字同判）。
- **新增包级依赖边＝0 枚（有尺）**。尺＝`GOFLAGS= go list -deps ./cmd/wisp | grep github.com/CarlosShao/wisp/internal` → 名册（14:0x）：`internal/agent internal/agent/approval internal/audio internal/ball internal/buildinfo internal/config internal/llm internal/llm/anthropic internal/llm/openaichat internal/llm/openairesponses internal/memory internal/models internal/observe internal/panel internal/perm internal/plugin internal/proc internal/projctx internal/risk internal/secret internal/session internal/statemachine internal/streamkey internal/tools internal/winsec`。`internal/audio`／`internal/ball`／`internal/statemachine` **三枚已在册** ⇒ 甲形接线在 `package main` 内完成 ⇒ 作差为空。
  ⛔ 反例也量了：`GOFLAGS= go list -deps ./internal/ball` → `internal/ball internal/observe internal/secret internal/statemachine internal/winsec`（**无 `internal/audio`**）⇒ 若把开关做进 `internal/ball` 让它直接吃门，会新增 1 枚 `internal/ball → internal/audio` 边。**该形状要避免。**
- **会不会把隐私默认改掉？——不。** 默认值那三枚一字不动（§2 表末行），`WithStartMuted(c.Audio.MicMutedDefault)`（`:213`）读的是同一枚 true；甲形加的是**运行期用户主动那一跳**。⚠ 但有两处说谎风险要照实写：`:238` 与 `:260` 两句是 boot 时刻的一次性 verdict，开门后不改就是过期读数；`ballGestureWhy`（`:408`）那句现在是假话。
- **额外一条必须写进代价的**：热键那支今天**已在注册表里**（`resident_ball_windows.go:267 residentBallHotkeyChain258`＋`[hotkey] mute` 档＝`PLAN.md:2732`），⛔ 不需要重注册，所以甲形不涉及 D36 的"须重注册全局热键"代价。

### 2.2 乙形：面板／托盘开关

拆成两半量，两半代价差一个数量级：

- **乙-托盘半（便宜，几乎等于甲）**：路由已经通到头（`ball_windows.go:673-684` 的 `case wmRButtonUp:` → `showMenu(...)` → `case menuMute:` `:678` → `OnTrayMute`），终点是 `resident_ball_windows.go:285 recordBallGesture("tray-mute")`。⇒ **文件数与甲同一枚集合（3 枚，`cmd/wisp` 内），新增依赖边 0 枚（同尺）**；只多一件：勾号需要人写 ⇒ 调用 `ball_windows.go:953 SetTrayChecks`（现零产码调用者），把门事件当唯一写者，⛔ 不许再攒第二枚布尔。契约撑腰句＝`SPEC-08:229`（托盘"静音"项已定名）。
- **乙-面板半（贵，且契约没给这一跳留门）**：
  - **要动的文件＝甲的 3 枚 ＋ 3～4 枚**：`internal/panel/bridge.go`（方法名册，现量逐字 `:66 MethodConfigGet = "config.get"` / `:67 MethodConfigSet = "config.set"`）、`cmd/wisp/panel_inbound.go`（入站白名单，`:55` 逐字注释 `it answers two of the six whitelisted methods`）、`internal/panel/config_handlers.go`（`config.set` 的落笔处），外加一枚处理支路。合计 **6～7 枚文件**。
  - **依赖边**：`GOFLAGS= go list -deps ./internal/panel` → `internal/observe internal/panel internal/projctx internal/risk internal/secret internal/streamkey internal/winsec`（**无 `internal/audio`／`internal/ball`／`internal/statemachine`**）。⇒ 走 `cmd/wisp` 内接线 ⇒ **0 枚新增**；若让 `internal/panel` 直接吃门 ⇒ **新增 1 枚 `internal/panel → internal/audio` 边**（并可能连带 `→ internal/ball`）。尺与名册都在上，⛔ 不是感觉。
  - **契约侧的硬代价（照实写）**：面板改配置那一跳的合法形状只有 `config.set [audio] mic_muted_default`（D36 GUI 规矩 `PLAN.md:2759-2760`＋`PLAN.md:1886`），而 `config_readers_255.go:150` 自述本次运行不读 ⇒ **乙-面板天然落在"写了文件但麦不会开"的坑里**，除非再加一枚活重读器；那条重读路径又与 `PLAN.md:1884`"KWS 开关不可热加载／需重启"读起来冲突（虽非同一枚键），必须在裁决里说清。
  - **历史开过的洞**：`Q-49`/票 114 那族"网页事件→Go"——契约现量边界＝`PLAN.md:1368`（C18 末句：面板来源的"允许"直接拒绝）＋ `SPEC-06:131-132`（`approval.decide` 的「允许」拒绝一切面板来源；PanelBridge 方法白名单＋每方法标注 capability 与是否需要原生侧授权）＋ `SPEC-08:159`（`correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志`）。⇒ **⚠ "面板来源不许拧门"这句契约里没有**，但它是 C18/SPEC-06:131 那一族的直接类比 ⇒ 乙形必须自己写死"这一跳的原生侧来源是什么"，否则就是把隐私开关交给网页事件。三形里**只有乙形会动到面板那三枚冻结件所在包**（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go` 是票面 AC#4 点名的冻结件；⛔ 我不主张会改它们，只主张落点同包、越界检查必红）。
  - **默认值**：乙-面板形最容易顺手做出"把默认改成开着"的错（用户认知：面板上是枚开关 ⇒ 关掉就等于默认开麦）。票面禁区第 1 条要逐字带着：`mic_muted_default` 的 true 一格不许动。

### 2.3 丙形：说话起始／唤醒词

- **今天不可能落地（现量三条，票面 AC#1 那一问的自证）**：
  1. `git ls-tree -r --name-only HEAD internal/speech` → **只有 `doc.go`**（一块没写）；
  2. 尺＝`git grep -n "wisp/internal/speech" HEAD -- cmd internal | grep -v _test` → **零命中（rc=1）**，全仓无产码 import 它；
  3. `wake_word.enabled` 默认 false（`schema.go:197`）＋ `EvKwsEnabled` 零发射者（`events.go:18`/`table.go:55` 两枚定义行，§0.1）＋ KWS 引擎本身在产码里不存在（C9 `WakeWordEngine` 只有契约行 `PLAN.md:1359`）。
- **要动的文件（若真做）**：`internal/speech/`（KWS 实现，新包多枚件）＋ `cmd/wisp/resident_audio_windows.go`（把 KWS 挂进采集腿）＋ `cmd/wisp/resident_ball_windows.go` 或装配根（发射 `EvKwsEnabled`/`EvWakeWord`）＋ `internal/config` 读取面（`wake_word.*` 现零产码读取者）⇒ **≥4 枚落点、且是跨包新功能**，不是一跳。
- **新增依赖边＝≥1 枚（有尺）**：`cmd/wisp → internal/speech` 是**史上第一枚**（尺 27 现证零 importers；名册见 §2.1 与 `go list -deps ./internal/speech` → 只有它自己）。
- **隐私默认：这一形最危险**。它把"开麦"从**一次显式动作**换成**设备常驻监听**；契约对此的交代是硬要求而非可选项——`SPEC-04:79`（默认不打开麦克风、KWS opt-in）、`PLAN.md:462`/`SPEC-04:80`（KWS 激活**常亮指示**）、`SPEC-08:71`（"麦克风开着的唯一交代"是 danger 常亮环，且 Conversation 那支还要 `#30 首次开启需 L2 级隐私确认`）。⇒ **丙形归"唤醒词/KWS 那一格"（P14、S4 那族），本票不许装得能走**（票面 AC#1 逐字同判）。

### 2.4 三形汇总（给编排者裁的那张表）

| 形 | 动的文件枚数 | 新增包级依赖边 | 真相源关系 | 隐私默认 | 今天能否落地 |
|---|---|---|---|---|---|
| 甲（热键→门；托盘同支） | **3**（全 `cmd/wisp`） | **0**（尺：`go list -deps ./cmd/wisp` 已含 audio/ball/statemachine） | 门 `muted` 为主，球 `Muted` 态与托盘勾号为从（`gate.go:20-21/:61-62/:166-167` 已写死主从） | 不变（默认 true 一字不动） | **能**——但"球真进 `Muted` 态"那一半**不能**（#8 需 `Armed`，`Armed` 零生产者 ⇒ 要么绕表 SetState、要么人工批准新边） |
| 乙-托盘半 | 同甲 **3**（+ 调 `SetTrayChecks`） | **0** | 同甲 | 不变 | 能 |
| 乙-面板半 | **6～7**（+`internal/panel/bridge.go`、`cmd/wisp/panel_inbound.go`、`internal/panel/config_handlers.go`） | 0（走 `cmd/wisp`）／**1**（若 `internal/panel→internal/audio`） | 若走 `config.set`，真相源变成**文件**，与门侧布尔成三源（文件/门/球） | 不变但**易被读成"关掉即默认开"** | 能，但"写了不生效"（`config_readers_255.go:150` boot-only）要另接一枚活重读 |
| 丙（唤醒词/说话起始） | **≥4**，跨包新功能 | **≥1**（`cmd/wisp→internal/speech`，史上第一枚） | 需新定义（KWS 与门谁主谁从契约没说） | **会变**：显式动作 → 常驻监听 | **不能**（`internal/speech` 只有 `doc.go`、零 importer、默认 false、`EvKwsEnabled` 零发射者） |

---

## 3. 够不到的地方（具名，不装）

1. **没跑门禁四数**（`go build`/`go vet`/`go test`/`d22scan`）——按指令让给 `289-r1`；本件所有结论都只靠 `git grep`/`git show`/`go list -deps`/`wc`。⇒ 若某枚形需要"能编译的接缝"，本件不能替它担保。
2. **没有真机读数**：`resident_windows.go:217/:278` 的装配顺序是从代码读的，`ra.vetoByEsc` 那条注入链（票 246 那族）我没有现跑到进程里；AC#3 那两形真机不在本腿射程。
3. **`go list -deps` 的口径**：它给的是**传递依赖名册**，能证"这枚包已在册／不在册"，⛔ 不能证"某一枚文件级 import 在哪一行"。我在 §2 每处新增边都写成了"名册作差"而不是"import 行级"，请读者按此口径读。
4. **`SPEC-08` §6 方法白名册的产码对应物**我只量到 `internal/panel/bridge.go:66-67` 两枚常量与 `cmd/wisp/panel_inbound.go:55` 那句注释，**没找到六枚白名册的完整产码名册落点** ⇒ 乙-面板那"3～4 枚文件"是下界，不是定数。
5. **`docs/specs/SPEC-09-platform-windows.md` / `SPEC-03` / `SPEC-05` 不在票面点名的四份里**，我没有对它们做名册尺（票面 AC#0 只列了 PLAN + SPEC-04/06/08）。若编排者要"全套 spec 都证一遍"，那一格我没做。

## 4. 与票面/代码原文的出入（本仓铁律：以原文为准，具名上报）

1. **票面 `:12` 的枚数与它自己的行号列表不符**：原文"其余**四枚**全在 `gate_test.go`（`:140/:147/:184`）"——列了三枚行号，现量也是**三枚**（尺见 §0 表：`SetMuted(` 命中 = 定义 `gate.go:168` + 注释 `:97` + 测试 `:140/:147/:184`）。⇒ 我按现量报"三枚"，不改票面一字。
2. **票面 `:18` 说"全仓唯一'把门拧开'的写法在测试里"、引 `resident_audio_247_live_windows_test.go:139`**：同一枚尺（尺 4 `MicMutedDefault`）现量到**第二枚同形写法**在同文件 `:244`（`dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })`）。⇒ "唯一"应读成"唯一那一族/唯一那个文件"，不是唯一一枚。结论方向不变（⛔ 都不是用户能走到的路）。
3. **票面 `:15` 的 `case wmAppTray:` 段范围写作 `:669-682`**：现量该 `switch` 块实际起于 `ball_windows.go:669`、`case menuExit:` 在 `:682`、整支 `return 0` 在 `:686`。⇒ 我引用时写 `:669-686`（含 `menuMute` 的 `:678` 与票面一致）。
4. **票面 AC#1 甲说"门侧另有一枚 `muted` 布尔"是两枚真相源**：现量**还有第三枚**——托盘勾号 `ball_windows.go:128 trayMuted`（写者 `SetTrayChecks:953` 零产码调用者）；另 `internal/tools/subagent_197.go:109` 复用了 `StateMuted` 常量名（异义同符）。⇒ 甲形的主从声明必须把第三枚一起钉住。
5. **票面把"契约许过那一跳"落在 `SPEC-04:47/:80`**（"一键全局静音 → `Muted` 态"）：我复量成立，但**同一支 D43 #8 的副作用只有 `kws.stop-inference`**（`table.go:66`/`PLAN.md:3068`/`SPEC-08:96`）——**契约从没写"静音键拧设备"**。⇒ "契约已经许过"这句要收窄成"已经许过发起者"，语义射程那一半是空白（详见 §1.3 ⓐ/ⓒ）。这是本件对票面性质判定 (`:4` 那枚 ★) 的一处**实质补充**，不是指它写错。
6. **`internal/config/tiers.go:29 "audio": "hot"`** 与 **`PLAN.md:1884`"不可热加载的段（KWS 开关…）"**、**`config_readers_255.go:150`"hot reload … waits for a restart"** 三句并读会互相打脸（键不同、档名同）。⇒ 不是票面错误，但落地腿必须在裁决里点名分清，⛔ 不许按其中一句静默实现。

## 5. 追加给票面的 Progress log 那一行（第 3 笔 commit 的内容，此处只留副本）

见票 `.scratch/wisp/issues/290-…mute-hotkey…md` 末尾；本腿未勾任何 AC、未改任何原句。
