# 票 290 — 门装上了、**没有任何一条生产路径能把它拧开**：默认档下球永远收不到电平（`SetMuted` 产码调用者＝0）

**立票**：2026-10-09 13:5x 编排者（机主令「及时补票」；来路＝票 247 非实现者验收腿 `247-v1` 的"哑键被吃掉之后下一枚哑键"这一问，件 `.scratch/wisp/probes/247/v1/50-seams-and-defaults.md`、`20-ac-verdicts.md`；根因由本编排者 13:5x 现跑两把尺自证，见下面现量）
**性质**：★**链路只接了一半，不是缺陷暴露**。票 247 把「真麦 PCM → 一枚 `float32` → 球的液态」这半条接完了，另半条——**"谁允许麦克风打开"**——今天在本仓**根本不存在**。⚠ 这一格与票 247 `AC#4`（默认档真机读数）是同一枚欠账的两面：那格量的是"默认档会不会偷开麦"，本票量的是"默认档能不能开麦"。**答案今天都是"不开"，一个是好事、一个是坏事。**

## 现量（每条都带尺；⛔ 引用前先重跑，枚数与行号都是快照）

- **门在装配时开在静音位（正确的保守形，票 247 P1 甲裁的那一支）**：`cmd/wisp/resident_audio_windows.go:213` 逐字
  `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`
  而 `mic_muted_default` 的默认值＝**true**（`internal/config/schema.go:277`）。静音位上**设备根本不打开**：`internal/audio/gate.go:97` 逐字
  `// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it`
- ★**`SetMuted` 的非测试调用者＝0 枚**（尺＝`git grep -n "SetMuted(" HEAD -- internal cmd`，13:5x 现跑）：命中只有**它自己的定义行** `internal/audio/gate.go:168`、`:97` 那句注释，其余**四枚全在 `gate_test.go`**（`:140/:147/:184`）。**`SetSpeaking` 同形**（尺＝`git grep -n "SetSpeaking(" HEAD -- internal cmd`）：定义 `gate.go:130` ＋ 注释 `:97` ＋ 测试 `:58/:68/:96/:99/:143/:144/:204` ⇒ **产码里没有任何一处调用那两枚唯一的开门手段**。
- **球上那枚静音热键今天只"记一笔"，不动门**：`cmd/wisp/resident_ball_windows.go:281` 逐字
  `OnMuteHotkey:    func() { recordBallGesture("mute-hotkey") },`
  对照同族里**唯一有执行者**的那枚（票 244 普查已量）：`:126`/今天的 `OnCancelHotkey → recordCancelHotkey → ra.vetoByEsc`。⇒ 手势名册（`internal/ball/ball_windows.go:49-59` 十枚字段；`case wmHotkey:` 在 `:656`、`case hkMute:` 在 `:660`）里 `mute` 那一支走到 `recordBallGesture` 就**落地为空**。托盘那支同名目（`case menuMute:` `:678`，`case wmAppTray:` 段 `:669-682`）同形。
- **今天全仓唯一"把门拧开"的写法在测试里**：`cmd/wisp/resident_audio_247_live_windows_test.go:139` 逐字
  `dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })`
  ⇒ 票 247 `AC#2` 那发真机读数是在**专门为它写的一枚临时配置**下取的，**不是用户能走到的那条路**。⛔ 这条不构成"接好了"，本票就是它剩下的那一半。
- **装配自己的文案已经承认这件事**：`resident_audio_windows.go:238` 逐字（默认档落点）
  `采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）`；`:260` 逐字 `采集腿在跑但设备未交接（gate 未 open）：球不会收到电平`。⇒ **本票不是"发现了一个没人说过的洞"，是"它说过了、然后没有票去做"**。
- **别把它读成隐私缺陷**：默认档不偷开麦＝**票 247 P1 甲裁的正确行为**，⛔ 本票不许把 `mic_muted_default` 的默认值改哪怕一格。要的是**一条用户能主动发起的开门路**，不是把门拆了。
- ★**契约其实已经许过那一跳**（13:5x 现读，⛔ 这条只是"不是无中生有"，AC#0 仍要整族搜）：`docs/specs/SPEC-04-voice-pipeline.md:47` 逐字
  `- 静音键：一键全局静音 → `Muted` 态；KWS 激活时悬浮球常亮指示（D16②）。`
  与 `:80` 逐字 `2. KWS 激活常亮指示 + 一键静音快捷键。`
  ⇒ 契约许诺的是**"一键切换"**，而今天 `hkMute` 那支走到 `recordBallGesture` 就断了；**球侧 `Muted` 态与门侧那枚 `muted` 布尔谁跟谁**这件事契约里没写＝本票 AC#1 甲形必须回答的那问。

## 要建什么（⛔ 先普查再落地；两格都不许改默认值）

- [ ] **AC#0 先答"契约里有没有许诺过谁来开麦"（只读，⛔ 不许只答"没找到"）**：在 `docs/PLAN.md`（`D16` 语音双工与隐私边界、`D5`、`D36` 配置模型那几节 + C 表）与 `docs/specs/SPEC-04-voice-pipeline.md`、`SPEC-06-security-gatekeeping.md`、`SPEC-08-ui-ball-panel.md` 里逐字搜并摘录"麦克风的开／关由谁发起"的原文句子（字样名册：`静音`/`mute`/`Muted`/`授权`/`privacy`/`唤醒`/`wake`/`push-to-talk`/`按住`/`开关`/`toggle`）。⇒ 三问：ⓐ 契约有没有**已经定死**的那一支（热键切静音／面板开关／按住说话／唤醒词）？ⓑ 有没有"必须每次显式、重启回到静音"这类**持久化许诺**（若有 ⇒ 归 `D36` 的生效级别，本票照它做）？ⓒ 如果契约一句都没写 ⇒ 本票升机主清单，⛔ 不许写码腿自己挑一支（那是功能决策，不是实现决策）。完成判据＝每条带 `文件:行` 逐字引文＋搜过的字样名册。
- [ ] **AC#1 三形代价表（⛔ 本格不改产码）**：把三条候选各写成"要动哪几枚文件＋门的状态从谁的真相源来＋会不会开出第二个真相源"：
  - **甲：静音热键那支接到门上**（`OnMuteHotkey` 现在走到 `recordBallGesture` 就断）。⚠ 硬约束：球侧 `Muted` 是 D43 状态（`internal/statemachine/table.go:65` 逐字 `D43: 8, From: StateArmed, Event: EvMuteKey, To: StateMuted,`；`:75`/`:79` 两枚 #10 `StateMuted` 的出去边），**门侧另有一枚 `muted` 布尔** ⇒ 两枚真相源就是本仓裁过两次的形状（票 246 的"第二个谁在等批准"同族）。必须写清**谁跟着谁**。
  - **乙：面板／托盘开关**（托盘名册里已有 `menuMute`，`internal/ball/ball_windows.go:669-682`，`case menuMute:` 在 `:678`）。⚠ 面板那条"网页事件→Go"那一跳历史上开过洞（`Q-49`/票 114），代价要照实写。
  - **丙：说话起始／唤醒词**（`wake_word.enabled` 默认 **false**，且 `internal/speech` 今天**只有 `doc.go`**＝一块没写）⇒ 这一支今天**不可能落地**，写清楚它是"归唤醒词票"那一格，不装得能走。
  ⇒ 判据＝三形各带现读凭据（`文件:行`＋尺读数），由编排者裁（若 AC#0 判"契约没写"则先摆机主）。
- [ ] **AC#2 落地那一发（⛔ 按在 AC#0/AC#1 裁之后）**：选定的那一形接上之后，**尺＝非测试调用者枚数**（`git grep -n "SetMuted(\|SetSpeaking(" HEAD -- internal cmd \| grep -v _test` 由 2 枚变 ≥3 枚，且新命中的那枚**在生产链路上**、⛔ 不许是又一枚调试用 `cmd`）。
- [ ] **AC#3 默认档一字不改，且两形真机读数都要有**：① **默认双击启动**：`mic_muted_default` 仍为 true、装配走 `WithStartMuted(true)`、球的 verdict 仍打 `:238` 那句"设备未打开"；② **用户主动开门之后**：`SetAudioLevel` 真的收到随声音变化的数（此形与票 247 `AC#2`/`AC#4` 的欠账**同批发**，⛔ 不许拿测试里那枚临时配置（现量最后一条）充当"用户能走到"）。两形各给采样数与出处；跑真机前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 必须为 **0**。
- [ ] **AC#4 越界检查**：`git diff` 名册里出现 `frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）／D43 转移表任一路径 ⇒ 直接退回。⛔ 若 AC#1 甲形需要新增状态或改转移边，**那不是本票**——那要先落一枚人工批准的 `A##`。
- [ ] **AC#5 门禁四数**：`GOFLAGS= go build ./...` rc=0；`gofumpt -l <自己动过的目录>` 空；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ -count=1` 改前改后各一次（**PASS/FAIL/SKIP 三数并排＋红名册差集**，⚠ 必须带 sherpa DLL harness，不带时那枚 `0xc0000135` 是环境红不是被验物）；`sh scripts/d22scan.sh` 纯净树 rc=0（⛔ `tools/d22scan` 是独立模块，`go run ./tools/d22scan` 必失败）；逐名照抄终态、每枚门禁件自落一行 `rc=N`。

## 禁区

- ⛔ **不改任何默认值**：`voice.enabled`（`schema.go:245` true）／`wake_word.enabled`（`:197` false）／`mic_muted_default`（`:277` true）三枚一字不动，读取处也不动。本票加的是"用户主动发起的那一跳"，不是把闸门挪位置。
- ⛔ **不许新造第二枚真相源**：门侧 `muted` 与球侧 `Muted` 必须有具名的主从关系，写进票面而不是留给读者猜。
- ⛔ 不许把"拒绝启动"扩大（票 128 只定了没有数据根那一种）；不许为开门新造协程（D38(b) 名册六枚零膨胀，票 247 P8 已裁电平在既有 `audio-capture` 协程内算）；⛔ 裸 `go func(`（ban #1）。
- ⛔ 不许顺手把票 247 的三格未勾（`AC#2`/`AC#4`/`AC#6`）勾掉，也不许顺手改票 291 那把电平尺的定义。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree。

## 排程与串行

- **两拍**：`290-a1`（只读普查＝AC#0＋AC#1）→ 编排者裁（若 AC#0 判"契约没写"⇒ 先摆机主一句话）→ `290-r1`（落地＝AC#2..AC#5）。
- 与票 247 的残余三格**同属地界**（都要那发"机主在场的真机两形"）⇒ **合批发、别各跑一次**；与票 291 互不阻塞（一枚是尺的定义、一枚是谁开门），但 ⛔ 不许混进同一枚写腿：票 291 若判"契约面"要等人的话，会把整批一起卡住。
- ⛔ `290-r1` 按在 `255-r2` 之后（两枚都碰 `cmd/wisp`，同包两枚写腿会把红的归因搅浑——本仓 09-30 实测过两次代价）。

**Status:** **未开工**（2026-10-09 13:5x 立票；`290-a1` 只读腿可先派，此刻 Go 编译/测试面空）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-09 14:15:16 +08] agent=290-a1 did=AC#0＋AC#1 只读普查交付件 `.scratch/wisp/probes/290/a1/00-findings.md`（38,744 字节；第 1 笔骨架 `fd692f7`、第 2 笔正文 `012ca47`，⛔ 零产码、零默认值改动、本行以外未动票面一字、零翻框、未跑 go build/vet/test 让给 289-r1）。AC#0 三问：ⓐ **发起者契约已定死**＝一键全局静音热键（`docs/PLAN.md:1312` 逐字「| `Muted` | 灰球 + 斜杠 | 一键静音（D16） | 再按静音键 |」同时定了进与出、`:462`、D43 #8/#10 `PLAN.md:3068/:3070`=`SPEC-08:96/:98`=`table.go:65/:75/:79`、次发起者托盘同名菜单项 `SPEC-08:229`）；⚠ 但 #8 副作用逐字只有 `kws.stop-inference`（`table.go:66`）⇒**「热键是否兼作设备开关」契约没写**；ⓑ 持久化只写了一半（`mic_muted_default` 归 `hot`（须重开设备）`PLAN.md:2735`＋`tiers.go:29`；GUI 态必须落文件 `PLAN.md:2759-2760`/`:1886`；产码语义「starts every session muted」`schema.go:276`），**「开门后重启回不回静音」四份契约零命中**（尺＝重启/持久/热加载/restart/persist 交叉过滤，只命中 `PLAN.md:549/:1972`、`SPEC-08:125` 三枚无关句）；ⓒ 故＝两处局部空白（语义射程＋持久化＋「谁发起把 KWS/麦打开」）⇒**摆机主**，⛔ 不由写码腿挑。字样名册 16 枚带命中数（`按住`/`push-to-talk`/`PTT`/`toggle` 全套契约 **0 命中**；`SPEC-06` 对麦克风 **0 命中**）。AC#1 三形：甲＝**3 枚文件**全在 `cmd/wisp`（`resident_audio_windows.go`/`resident_ball_windows.go:256/:281/:285/:408`/`resident_windows.go` 装配根顺序 `:217` 早于 `:278` 需后置 setter），**新增包级依赖边 0 枚**（尺＝`GOFLAGS= go list -deps ./cmd/wisp` 名册已含 audio/ball/statemachine；反例已量：`go list -deps ./internal/ball` 无 audio ⇒ 做进 ball 会加 1 枚边）；真相源＝**门 `muted` 为主、球 `Muted` 态与托盘勾号为从**，此主从是产码注释早写死的（`gate.go:20-21`/`:61-62`/`:166-167`），⚠ 但球今天**无路进 `Muted`**（`EvMuteKey`/`EvKwsEnabled` 产码零发射者 ⇒ 真要态跟随须人工批准 `A##`，即本票 AC#4 那一格）；乙-托盘半＝与甲同集合（另需 `SetTrayChecks` `ball_windows.go:953`，现零产码调用者）；乙-面板半＝**6～7 枚文件**（＋`internal/panel/bridge.go:66-67`、`cmd/wisp/panel_inbound.go:55`、`internal/panel/config_handlers.go`），边 0 枚若走 `cmd/wisp`／**1 枚**若 `internal/panel→internal/audio`（已量名册），且天然落「`config.set` 写了文件但本次运行不开麦」坑（`config_readers_255.go:150` 产码自述 boot-only）＋面板来源的隐私那一跳契约无授权形状（参照 `PLAN.md:1368` C18、`SPEC-06:131-132`、`SPEC-08:159`）；丙＝**今天不可能落地**（`internal/speech` 对象层只有 `doc.go`、`git grep wisp/internal/speech` 产码零命中、`wake_word.enabled` 默认 `false` `schema.go:197`）⇒ ≥4 枚跨包文件＋史上第一枚 `cmd/wisp→internal/speech` 边，且把显式动作换成常驻监听＝三形里唯一真会碰隐私默认的一支，归唤醒词票。另新量三枚票面没说的料：`WithGateEvents` 产码调用者＝0（门自建的钩子也没接）、托盘勾号 `trayMuted` 是**第三枚**真相源、`ballGestureWhy`（`resident_ball_windows.go:408`）那句「no microphone」自票 247 起在本进程为假。票面四处枚数/范围出入已具名记在件内 §4（含 `:12` 的「四枚」现量三枚、`:18` 的「唯一」在同文件 `:244` 还有第二枚同形）。next=编排者裁三形；若采纳我 §1.3 的读法（发起者已定死、语义射程未定死），**先摆机主一句「静音键要不要兼作设备开关＋开门后重启回不回静音」**再派 `290-r1`（按票面排程须在 `255-r2` 之后）；AC#2..AC#5 本腿一概未做、未勾。
