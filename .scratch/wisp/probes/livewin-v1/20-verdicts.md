# livewin-v1 — 逐格判语（五格；票 247 AC#2/AC#4/AC#6 · 票 290 AC#3 · 票 291 AC#2）

票面原句一律 `sed -n 'Np'` 现抽（HEAD `d432d728`）。凭据行号＝盘上原始件现量行数（`wc -l` 见 00 号件）。
尺的编号 R1…R8 见 `10-rulers.md`。**本腿没有勾任何框，没有改任何产码。**

⚠ 引数口径先说清（派单点名的"枚数"事故）：本腿全部判语建立在**两枚单测台件的 raw**（`go test ./cmd/wisp/ -run 'TestAC247LiveMicrophoneLevelsReachTheBallSeam' -count=1 -v`，findings `:16` 那把尺，**带 `-v`、带 `-run` 过滤，不是整包名册**）与**一枚常驻 exe 的 stdout**（c1c2 raw）之上；本腿**没有引任何"整包红枚数"**，所以"同树两发红枚数漂 ±1"那一类问题在这五格里不起作用。本腿的 grep/sed/wc/awk 诸尺均为**该目录/该文件的整族穷举**（R1 已具名 ls 全族）。

---

## 格 1 · 票 247 `AC#2` ⇒ 【不成立】

**票面逐字原句**（`.scratch/wisp/issues/247-…md:25`）：
> - [ ] **AC#2 球上的数来自声音，不来自命令行**：一发真机读数＝对麦克风说话与不说话时 `SetAudioLevel` 收到的值**不同**，并给出两形的采样数与出处；⛔ 不许用 `balldebug` 的合成包络冒充（现量 2 那条尾巴要剪掉，不是把它做大）。

**支撑它的凭据（逐字＋行号）**：
- `…a1-run1-quiet.raw.md:18` `AC#2 reading: quiet-A: samples=93 mean=0.367691 max=0.409285 levels_per_s=31.00`
- `…a1-run1-quiet.raw.md:19` `AC#2 reading: sound: samples=189 mean=0.382482 max=0.413813 levels_per_s=31.50`
- `…a1-run1-quiet.raw.md:20` `AC#2 reading: quiet-B: samples=94 mean=0.383328 max=0.408784 levels_per_s=31.33`
- `…a1-run1-quiet.raw.md:21` `resident_audio_247_live_windows_test.go:195: the loud phase is not distinguishable from silence: … (speakers off, or routed to an endpoint this microphone cannot hear)`
- 出处那一半有：`…a1-run1-quiet.raw.md:4` `AC#2 出处: capture="麦克风 (Realtek(R) Audio)" render="扬声器 (Realtek(R) Audio)" pinned_thread=14896`
- 判据本体（产码现读）：`cmd/wisp/resident_audio_247_live_windows_test.go:194` `if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max() {`

**本腿独立尺**：R4 ⇒ run1 `sound/quiet-A` = **1.0402**（判据要 ≥2 且 max 更大：max 那半确实过，0.413813 > 0.409285）；run2 = **1.0334**。★**决定性一枚**：R4 同条件漂移 run1 `quiet-B/quiet-A` = **1.0425 > 1.0402**——"安静 vs 安静"的差比"响 vs 安静"的差还大，所以那 4% 不可归因给声音。

**路径算不算（派单点名那一问，本格先落地）**：**算**。AC#2 原文只禁一件东西——`balldebug` 的合成包络；它要的是"一发真机读数 + 采样数 + 出处"，那两发满足：台件 `:147` 调的是生产装配 `assembleCapture`、真麦、`WISP_LIVE_MIC=1`（台件 `:129`/`:130` 是硬闸）。⇒ 测试内那一发在 247 `AC#2` 这里是合格凭据，本格**不是**因为路径而挂，是因为读数而挂。

**最坏能被怎么糊过去**：①把 `:181` 那行 `levels=376`／`levels_per_s≈31` 读成"球收到了来自声音的数"——那是 `AC#1` 的"有人调用"，不是本格的"值不同"；②按字面"值不同"判绿（0.382482 ≠ 0.367691 数值上真的不同）——本腿试了，这条路只要把判据读成词面型就通，而票面括弧（"现量 2 那条尾巴要剪掉"）指的是**来源**归因，1.0425 的同条件漂移正好证明这 4% 不是归因。**本格判【不成立】而非【不可判】的理由**：两形都有采样数与出处，读数本身给出了否定答案；缺的不是窗口，是票 297 那一问（喂给尺的那串数到底是什么）。⇒ 勾框后果＝这格**继续开着**，且**不该由下一次窗口自动解决**；判"不可判"才会把它错误地挂到"等机主再来一次"。

---

## 格 2 · 票 247 `AC#4` ⇒ 【成立】（带两条具名限界）

**票面逐字原句**（`247-…md:27`）：
> - [ ] **AC#4 隐私默树一格都不许动**：`voice.enabled`／`wake_word.enabled` 的**默认值与读取处**一字不改，且真机读数要证明"默认档下麦克风不会被这条腿打开"（或反过来具名说出它打开了、由哪一行决定）。

**支撑它的凭据（逐字＋行号）**：
- `…c1c2-resident.raw.md:23` `time=2026-10-09T17:27:04.947+08:00 level=INFO msg="audio: capture armed but muted at boot; device not opened" config_source=C:\Users\swq\AppData\Roaming\wisp-dev\config.toml path=T`
- `…c1c2-resident.raw.md:25` `wisp: 采集腿：采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 C:\Users\swq\AppData\Roaming\wisp-dev\config.toml）`
- ★来源＝**机主自己那份 config.toml**（不是临时配置目录），这一条是本格比 `AC#2` 更硬的地方。
- "反过来具名说出它打开了、由哪一行决定"也有：`…c1c2-resident.raw.md:38`（17:46:21.482 `gesture=tray-mute outcome="已取消静音：设备已交接，采集线程在跑…"`）＋ findings `:56` 把四个开关时刻逐字列表。

**本腿独立尺**：
- 默树未动：`internal/config/schema.go:245`（`VoiceSection`，`:244` 起）= `Enabled bool \`toml:"enabled" default:"true"\``；`:197`（`WakeWord`，`:196` 起）= `… default:"false"`；`git log --oneline -S'Enabled bool \`toml:"enabled" default:"true"\`' -- internal/config/schema.go` ⇒ 只命中 `fe126e27`（schema 落地那一笔），票 247 全程没碰。`mic_muted_default` 同尺同结果。
- 读取处现状：`cmd/wisp/resident_audio_windows.go:257` `ra.voiceEnabled = c.Voice.Enabled`、`:262` `if !c.Voice.Enabled {`、`:273` `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`；`WakeWord.Enabled` 的产码命中只有 `internal/config/manager.go:380`/`:395`（热加载比对，不开设备）。
- 机主那份的生效值（R8）：`[voice] enabled = true`、`[voice.wake_word] enabled = false`、`[audio] mic_muted_default = true` ⇒ 三枚全部**与编译默认一致**，"默认档"这一前提在那一发窗口里是真的。
- 代码侧闭环：`internal/audio/gate.go:97`/`:108` 关门态不 start inner ⇒ 设备根本没开。

**两条具名限界（不许被读成缺陷，但要写进表）**：
1. 那一发真机读数**没有直接打** `voice_enabled`（R8：`grep -c "voice_enabled"` 在 c1c2 raw ＝ **0**；那句 Warn 在 `resident_audio_windows.go:252`–`:253`，只在配置文件读不到的兜底支走）。`voice.enabled` 那一半是靠"配置文件现量＋`:257`/`:262` 代码分支"反推的。
2. 本格的两枚默树里 `wake_word.enabled` 今天**没有生产读取处**（只有 config 层比对），所以"读取处一字不改"对它是空真；这一条票面没写、本腿替它写明，免得被读成"KWS 那条路也验过了"。

**最坏能被怎么糊过去**：拿"门在静音位"外推成"这条腿永远不会开麦"——同一场窗口 17:46:21 一按托盘它就开了；本腿试了这条（读 `gate.go` 的 `SetMuted` 支 `:168`–`:190`），并确认票面第二支（"或反过来具名说出它打开了、由哪一行决定"）要求的就是把这个说干净，凭据说了。⇒ 判【成立】。

---

## 格 3 · 票 247 `AC#6` ⇒ 【不可判(欠 三形里的任两形真机读数)】＋"没跑"归谁

**票面逐字原句**（`247-…md:29`）：
> - [ ] **AC#6 降级照跑**：拔麦／无权限／设备被占三形里任取两形，本进程**仍要把球与任务管线跑起来**并且把损失说响亮（现量：票 128 定的是"没有数据根才拒绝启动"，别把它扩大）。

**凭据里有没有**：**没有读数，只有一句自陈**。`2026-10-09-live-window-findings.md:87` 逐字：
> | 票 247 `AC#6` 真设备两形（禁用设备／撤权限） | **没跑**：那要改他 Windows 的系统设置，⛔ 未经他一句话不动 | 票 247，排程见票面 |

**本腿独立尺**：把那一行与票面三形对照 ⇒ **口径不一致**：票面三形是「拔麦／无权限／设备被占」，findings 那一行只列了「禁用设备／撤权限」两形（≈"无权限"那一支的两种做法）。逐形可及性：
- **无权限／禁用设备**：确实要动机主系统设置，他不在键盘前 ⇒ 理由成立。
- **设备被占**：那一形**既不需要他的设置、也不需要他人**——一枚抢同一采集端点的探针就能造形。凭据里没有任何一句说它为什么没跑 ⇒ ★**findings 那一行把整格归给"要改他 Windows 的系统设置"是说轻了**（三形里至少一形今天可自造而未做）。
- **拔麦**：本机 capture 端点是 `麦克风 (Realtek(R) Audio)`（a1 raw `:4`），是不是有"可拔"之物凭据里没写；本腿不替它编理由。

**归谁**：**编排者**。派单对本腿明令零 go 命令、不起 `wisp.exe`、不碰麦克风 ⇒ 按本仓规矩（需要真窗/需要跑才能裁、而那资源在派单里就被禁的格，不算腿的欠账）这格记在编排者名下。**不是**机主的欠账，除非只算"禁用设备／撤权限"那一形。

**最坏能被怎么糊过去**：拿 c1c2 里 17:46:44 那次"已静音：采集已关闭"当降级形（那**不是**降级，是用户主动关门）；或拿 findings `:22`/`:86` 那种"任务入口未启用、进程照常"当"损失说响亮"（那是 stdin 形状，不是设备形状）。本腿两条都试过：raw `:63` 与 `:21` 逐字读法不同，谁也证不了 `AC#6`。

---

## 格 4 · 票 290 `AC#3` ⇒ 【不可判(欠 用户路径上的电平采样数与出处；另欠热键形与 tasklist 现量)】

**票面逐字原句**（`.scratch/wisp/issues/290-gate-mounted-nobody-can-open-it.md:36`）：
> - [ ] **AC#3 默认档一字不改，且两形真机读数都要有**：① **默认双击启动**：`mic_muted_default` 仍为 true、装配走 `WithStartMuted(true)`、球的 verdict 仍打 `:238` 那句"设备未打开"；② **用户主动开门之后**：`SetAudioLevel` 真的收到随声音变化的数（此形与票 247 `AC#2`/`AC#4` 的欠账**同批发**，⛔ 不许拿测试里那枚临时配置（现量最后一条）充当"用户能走到"）。两形各给采样数与出处；跑真机前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 必须为 **0**。

**形①——到位，三处带注**：
- `mic_muted_default` 仍为 true：凭据 `…c1c2-resident.raw.md:25` 逐字含 `[audio] mic_muted_default=true，来源 C:\Users\swq\AppData\Roaming\wisp-dev\config.toml`；本腿 R8 现读了那份文件（`mic_muted_default = true`）。
- 装配走 `WithStartMuted(true)`：产码现量 `resident_audio_windows.go:273` 传的是**变量** `c.Audio.MicMutedDefault`，真机取值 true。⇒ 按"生效值"成立、按"字面量"不成立，这一条要在表里写平。
- "球的 verdict 仍打 `:238` 那句"：**那句还在，行号漂了**。本腿 `sed -n '238p'` ⇒ `cfgPath := filepath.Join(dataDir, configFileName)`；`grep -n 设备未打开` ⇒ 现活 `:298`（另一枚同词在 `:172` 的关门文案）。⇒ 判据的**内容**满足，判据的**指针**过期（R3）。
- ⚠ 本格硬条件"跑真机前两发 tasklist 必须为 0"：**盘上零读数**（R7：runsheet `:45`/`:46` 只有计划，raw 与 logs 里都没有 tasklist 输出），且 findings `:92` 作者自陈没确认。间接证据只有一枚：raw `:7` `single instance = true` ＋ 17:27:04 成功起住。

**形②——没到位，且不许追认**：
- 开门那一刻有（凭据 `:38`/`:63`/`:65`/`:79` 四时刻逐字，走的是 `wisp.exe` 托盘手势，生产调用点 `resident_audio_windows.go:162 gate.SetMuted(!gate.Muted())`）。
- "采集确在跑"有，但**只证到"设备已交接、包在出"这一层**（R5：`push` 在 `stream.Drain()` 之后；共享模式采集在静音房间照样按点交包）⇒ 这把尺证不了"麦克风在被有效读"，更证不了"随声音变化"。findings `:61` 那句"开门窗口里设备确实在采集"作为"包在出"是**站得住的**，作为"设备在被有效采集"是**说过一层的**；同一句里还有**一处算术混基线**：它写"第二扇 `722 → 1073` ⇒ ≈369 帧 / 11.8 秒"，而 `1073-722 = 351`（×32 ms = 11.23 s；11.8 s 对应的是 `1073-704 = 369`，把关门前沿的 18 帧也算进了第二扇）。
- ★**本格要的"采样数与出处"在用户路径上是零**：常驻那条腿**没有任何一处把电平数值打到盘**（现读：`resident_audio_windows.go:122 ra.levels.Add(1)`、`:139 rb.b.SetAudioLevel(level)` 只计数与投递，不打印数值；findings `:86` 自认"常驻路径零打印器"）。唯一带采样数/mean/max 的读出来自**测试内装配 + `:139` 那枚临时配置**（`c.Audio.MicMutedDefault = false`）——那恰是票面括弧明令不许拿来充当"用户能走到"的形状。
- ⇒ **"文案相同不等于路径相同"这一条本腿支持，但结论要落在更前面一句**：连"文案"都没声称完成——`resident_audio_windows.go:174` 那句逐字带 `（球屏上会不会呼吸是票 68 那一格，本票不声称）`，产码文案自己拒称。

**为什么是【不可判】不是【不成立】**：在**票面要求的那条路径**上，形② 从未取过数 ⇒ 谈不上被证伪；"随声音变化"被证伪的是**另一条路径**（测试内那两发，格 1 已判）。判【不成立】会把这格错误地并入票 297 那一问（仪器问题）；判【不可判】才把欠账指对：**欠一枚常驻路径的电平打印器 + 机主在场再开一扇的读数**（外加热键形——findings `:63` 已具名"热键形不算已读"，且票 296 那条根因把这一形挡在今天之外）。

**最坏能被怎么糊过去**：①"编排者说拿到了"＋两处逐字产品文案相同 ⇒ 追认；②把 dropped_frames 那把尺读成"设备真在采集声音"；③把测试内那发的 `samples=93/mean/max` 挪作用户路径读数。三条本腿都逐条试过（R5/R6/R7），任何一条都撑不住。

---

## 格 5 · 票 291 `AC#2` ⇒ 【不可判(欠 B② 纯人声窗一形、B③ 的理论 RMS 仪器、B④ 合成夹具)】

**票面逐字原句**（`.scratch/wisp/issues/291-level-ruler-does-not-strip-dc.md:17`）：
> - [ ] **AC#2 把"该读多少"变成读数而不是意见**：现跑同一枚真麦（`WISP_LIVE_MIC=1` 那族台件），采**四形**各 ≥3 秒并逐形给 `mean/max`：① 完全安静 ② 机主正常音量说话 ③ 外放已知幅度音调（先算出该音调的理论 RMS）④ 人为注入纯 DC＋已知正弦的合成帧（仓内夹具，⛔ 不许据此判生产）。⇒ 判据＝**"说话"与"安静"两形之间存在一个既有的、不动用的判据能分开**——分开不了就具名说"分开不了"，⛔ 不许挑一个能过的门槛凑绿。

**四形名册（每形带票面逐字定义，读数出处＝a1 两发 raw `:18`–`:20`）**：详见 `30-roster-and-debts.md` §1。摘要：**① 有数（两窗×两发）**／**② 零**／**③ 半（mean/max 在，"已知幅度＋理论 RMS"零）**／**④ 零**。

**支撑/否定它的凭据（逐字＋行号）**：
- ① 成立到"≥3 秒"这一尺：`…a1-run1-quiet.raw.md:18` `samples=93 … levels_per_s=31.00` ⇒ 93/31 ≈ 3.0 s（`levels_per_s` 那把尺由台件 `:190` 的 `1 / audio.FrameDuration.Seconds()` 定义，与 `frame=32ms`（raw `:3`）一致）。
- ② 为零的**结构性**凭据（不是编排者的自陈）：`cmd/wisp/resident_audio_247_live_windows_test.go:142` `loud := &livePhase{name: "sound"}` ＋ `:171` `stopNoise := playAlarmThroughSpeakers(t)` ＋ `:232` `media := filepath.Join(os.Getenv("WINDIR"), "Media", "Alarm01.wav")` ⇒ "响"那一窗**硬编码同时外放**，人声只会叠进同一窗；台件 `:234` 自己写着 `AC#2 then needs a person speaking into the microphone` 却没有实现那一形。
- findings 的具名更正在盘：`:42`（"B②（纯人声窗）没有到手……B③/B④ 缺仪器、今天根本没跑"）。

**本腿独立复核这条更正——既说重了一处、也说漏了一处**：
- **说重了（B③）**：更正写"B③……今天根本没跑"，而 B③ 的**测量那一半其实跑了**（外放 `Alarm01.wav` 那一窗，189 样本、mean/max 齐）；真正缺的是票面括弧里"先算出该音调的理论 RMS"那一半（findings `:83` 自己写的就是这一半："没有任何件打印一枚 wav 的 RMS"）。⇒ 名册应记 **B③＝半形有数**，不是零。
- **说轻了（B①）**：更正把 B① 记为"到手"，但 B① 的两窗**互差 1.0425**（R4：run1 quiet-B/quiet-A），比"响 vs 静"的效应还大 ⇒ "完全安静"这一形自带的离散度已超过本格要找的判据灵敏度；把 B① 干净地记成"到手"会让读者以为基线可信。本腿记 **B①＝有数但基线不自洽**，并具名 quiet-B 比 sound 更高这一反序（凭据 `…a1-run1-quiet.raw.md:20` vs `:19`）。
- **B④＝零**，与更正一致（无夹具、无能吃仓内帧的打印器）。

**判语理由**：票面判据点名要用**①与②两形**去分；②今天**没有任何一形**，所以"既有判据能否分开"这一问在这张票上**尚未被回答**——既不是绿，也不是"具名说分开不了"（那句"分开不了"必须由**四形读数**支撑，不能由替代形状支撑）。⇒ 【不可判(欠 B②/B③半形/B④)】。**⛔ 本腿不建议动 `resident_audio_247_live_windows_test.go:194` 那枚 2× 门槛**，也**不把现象读成"麦克风坏了"**：现象已由票 297 立票（`297-…md:1` 标题即为"安静房间的 delivered level 就有 mean≈0.3677、六窗 max/mean 全挤在 1.08–1.11"）。★顺带纠一处**票 297 标题里的过期数**：R4 实测六窗 max/mean 全距是 **1.0531–1.1131**，"全挤在 1.08–1.11"对 3/6 窗不成立（findings `:38` 同错，标题抄了它）。

**最坏能被怎么糊过去**：①把 run2 的 `sound`（外放＋人声叠窗）当 B②"机主正常音量说话"；②把"编排者已在 findings 里具名更正"当作"这格已自证清白、可以勾"；③拿 ③那一窗的 mean/max 冒充 B③ 完整形（缺理论 RMS 就等于缺票面要求的自变量）。三条本腿都试了，①②③都不能使本格变绿。
