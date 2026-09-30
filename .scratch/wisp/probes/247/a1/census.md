# 247-a1 — 采集栈接球：四问普查件（只读腿）

- 腿：`247-a1`（只读普查）· 票：`.scratch/wisp/issues/247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md` AC#0
- 取数时刻：`2026-09-30 23:04 +0800` 起（`date` stdout），HEAD `f0707c7b`，分支 `dev`
- 本件**不改任何产码**；写点只有这一枚文件。
- ⚠ 本文凡引用 `cmd/wisp` 的行号一律标〔起手读数，可能被 246-r2 改动〕。
- ⛔ 本文不提任何词面型仪器；`frontend/**`／`design/**` 未读、未引。

## 目录（先 ⑤⑥⑦ 后 ①-④，按派单硬要求）

- ⑤ 我为了这份表跑了哪些尺、每条的真实读数
- ⑥ 我可能写错的条目（对抗我自己；必交物）
- ⑦ 判不动的地方（要编排者裁：甲／乙／不做）
- ① 落点：采集协程归常驻腿还是归 `internal/audio` 自起
- ② 隐私闸门：接完之后麦克风在什么条件下真的打开
- ③ 降级：设备被占／无权限／无麦／热插拔现在返回什么
- ④ 第六环（形状）：`liquidDriven` 当前实际含哪几态，补 `Armed` 要动什么

---

## ⑤ 尺清单与真实读数

全部现跑于 `2026-09-30 23:04–23:08 +0800`（`date` 的 stdout，逐条附钟点），HEAD `f0707c7b`，GOOS=windows／GOARCH=amd64（尺 U18 现量）。
命令逐字可复跑；`⇒` 后面是 stdout 的真实读数，不是我概括的话。

| # | 尺（原样命令） | 取数钟点 | 读数（原样） |
|---|---|---|---|
| U1 | `grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . \| grep -v "/internal/audio/" \| grep -vc _test` | 23:04:47 | `0` |
| U1b | 同上但去掉末段 `\| grep -vc _test`（要名册不要数） | 23:04:47 | 空（grep 退出码 1，即一枚都没有） |
| U2 | `grep -rl "CarlosShao/wisp/internal/audio" --include=*.go internal cmd tools` | 23:07:08 | 空（退出码 1）＝**产码三棵树里零 importer** |
| U3 | `grep -rn "SetAudioLevel" --include=*.go internal/ cmd/` | 23:07:02 | 产码调用点只有 2 枚：`cmd/balldebug/main.go:418`、`:477`；定义 `internal/ball/liquid_windows.go:42`；其余 4 处是注释（`internal/audio/level.go:112`、`internal/ball/ball_windows.go:108`、`internal/ball/liquid.go:13`、`liquid_windows.go:14`/`:24`/`:35`）；测试调用 4 处全在 `internal/ball/live_windows_test.go:637/648/656/685` |
| U4 | `ls internal/speech/` | 23:04:12 前后 | `doc.go`（一枚） |
| U5 | `grep -rl "CarlosShao/wisp/internal/ball" --include=*.go . \| grep -v "/internal/ball/"` | 23:04:47 | `./cmd/balldebug/main.go` · `./cmd/wisp/resident_approval_windows.go` · `./cmd/wisp/resident_ball_228_test.go` · `./cmd/wisp/resident_ball_windows.go` |
| U6 | `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/audio` | 23:04:47 | 直接依赖含 **`internal/observe` 一枚 wisp 包**，其余全是标准库＋`golang.org/x/sys/windows`。**不含 `internal/ball`** |
| U7 | 同 `./internal/ball` | 23:04:47 | 直接依赖含 `internal/observe`＋`internal/statemachine`。**不含 `internal/audio`** |
| U8 | 同 `./cmd/wisp` | 23:04:47 | 含 `internal/ball`、`internal/config`、`internal/observe`、`internal/proc`、`internal/statemachine`…；**不含 `internal/audio`** |
| U9 | `go list -deps ./internal/audio \| sort`／`./internal/ball`／`./cmd/wisp`，行数 | 23:05 | `144` / `145` / `276` 枚（名册文件已留在本探针目录：`deps-audio.txt`／`deps-ball.txt`／`deps-cmdwisp.txt`，按"临时件只建不删"保留） |
| U10 | `comm -23 deps-audio.txt deps-cmdwisp.txt` ＝「`cmd/wisp` 若 import `internal/audio` 会多出哪几名」 | 23:05 | **只有 1 名**：`github.com/CarlosShao/wisp/internal/audio`。其余 143 名 `internal/audio` 的传递依赖（含 `internal/observe`、`x/sys/windows`）**已在 `cmd/wisp` 名册里** |
| U11 | `grep -rn 'toml:"enabled" default' internal/config/schema.go` ＋ `grep -n "type VoiceSection struct\|type WakeWord struct\|type AudioSection struct"` | 23:06 | `:197`（`WakeWord.Enabled`，`wake_word`，`default:"false"`）· `:223`（`AEC.Enabled`，`default:"true"`）· `:233`（`Realtime.Enabled`，`default:"false"`）· `:245`（`VoiceSection.Enabled`，`voice`，`default:"true"`）· `:371`／`:530`（别的 section，与本票无关）。struct 头：`WakeWord` 在 `:196`、`VoiceSection` 在 `:244`、`AudioSection` 在 `:269` |
| U12 | `grep -rn "MicMutedDefault\|mic_muted_default" --include=*.go .` | 23:06 | `internal/config/schema.go:276-277` ＝ **`MicMutedDefault bool` `toml:"mic_muted_default" default:"true"`**（第三枚开关，默认**静音**）· `internal/audio/gate.go:55-56`（`WithStartMuted` 的映射说明）· `internal/config/boundary_test.go:103` |
| U13 | `grep -rn "Voice\.Enabled\|WakeWord\.Enabled" --include=*.go . \| grep -v _test` | 23:06 | 产码命中**只有 3 处，全在 `internal/config/manager.go`**：`:370`（reload 计划的比较）、`:384`、`:385`（把新值抄回活配置）。**没有任何一处产码读它来决定开不开麦** |
| U14 | `grep -rn "NewWASAPIMicrophone\|NewHalfDuplexGate\|NewBoundedFrames\|FrameLevel\|SpawnCapture\|NewWavInjector" --include=*.go internal cmd`（排 `internal/audio/` 与 `_test.go`） | 23:07:08 | **空（退出码 1）**＝采集器、门、有界 channel、电平函数**没有一枚产码构造／调用者**（第一版这把尺没排 `.scratch`，被 `241/v1` 的变异体 fixture 淹了 30 行假命中，已重跑，见 ⑥-E2） |
| U15 | `grep -rn "EnablePrototypeVisuals" --include=*.go .` | 23:07 | 产码调用者**只有 `cmd/balldebug/main.go:122`**；其余全在 `internal/ball/*_test.go`。定义与默认值：`internal/ball/statevisual.go:102` `var prototypeVisuals bool`（＝**默认关**），开关在 `:106` |
| U16 | `grep -rn "EvKwsEnabled\|StateArmed" --include=*.go internal cmd \| grep -v _test` | 23:07:29 | `EvKwsEnabled` 只出现在 `internal/statemachine/events.go:18`（定义）与 `internal/statemachine/table.go:55`（D43 #5 那一行）。**产码里零发射者**。`StateArmed` 的产码命中：`internal/ball/statevisual.go:158`/`:288`（视觉映射）、`internal/statemachine/states.go:13`/`:41`、`table.go:61/65/69/75/236`、`cmd/balldebug/main.go:54`/`:485`（手填状态名册）。**`cmd/wisp` 一枚都没有** |
| U17 | `grep -rn "Armed" internal/ball/tokens_table_test.go internal/ball/tokens_test.go` | 23:07:39 | `tokens_table_test.go` 里 **零命中**；`tokens_test.go:141` 注释点名「SPEC-08 §2.1 anchors … Armed」，`:155-156` 断言 `VisualFor(p,56,StateArmed,0).Opacity == 0.6`，`:298`／`:318` 把 Armed 列进静态态名册 |
| U18 | `go env GOOS GOARCH` ＋ `GOOS=windows go list -deps ./cmd/wisp \| sort` 与裸跑的 `diff` | 23:07 | `windows` / `amd64`；两份名册**逐字相同**（276 行，diff 空）⇒ 本票所有名册读数不受 GOOS 前缀影响 |
| U19 | `grep -rn "^func " internal/audio/wasapimic_windows.go`（＋同法过 `audio.go`／`gate.go`／`device.go`／`wasapi_windows.go`） | 23:05 | 采集栈公开面：`NewWASAPIMicrophone` `:52` · `Start` `:69` · `Stop` `:94` · `Err` `:119` · `ThreadID` `:130` · `Stats` `:133` · `Endpoints` `:137` · pinned 循环 `run` `:149`；`SpawnCapture` `audio.go:180`；门：`NewHalfDuplexGate` `gate.go:84` · `Start` `:99` · `SetSpeaking` `:130` · `SetMuted` `:168` · `effectiveOpen` `:234` |
| U20 | `grep -rn "RegisterShutdownHook\|StepStopAudio" --include=*.go internal/proc cmd` | 23:06 | 接缝**已经存在**：`internal/proc/boot_windows.go:188` `func (rt *Runtime) RegisterShutdownHook(step ShutdownStep, hook …) error`；`internal/proc/shutdown_hooks.go:51`/`:118-119` 把 `StepStopAudio` 落到 `ShutdownHooks.StopAudio`；`shutdown.go:163` `run(StepStopAudio, hooks.StopAudio, 0)`。产码里当前**唯一**注册者是 `cmd/wisp/resident_windows.go:156`（`proc.StepCancelTasks`，票 246 AC#4）〔起手读数，可能被 246-r2 改动〕 |
| U21 | `sed -n '11,31p' internal/proc/shutdown.go` ＋ `sed -n '278,290p' docs/PLAN.md`（D38(e) 十步） | 23:06 | 代码里的十步与 `PLAN.md` 一字对齐；第 4 步 `stop-audio` 的名字逐字 `audio capture thread stops -> device closed   <- BEFORE step 5`。第 2 步 `stop-hotkey-kws`。名册：`shutdown.go:49-60` |
| U22 | `grep -rn "observe.Default.Spawn\|Registry.Spawn\|residency" internal/audio/*.go`（看采集协程走谁） | 23:06 | `internal/audio/wasapimic_windows.go:78` ＝ `observe.Default.Spawn("audio-capture", "audio", nil, …)`（**写死 `observe.Default`**，不接 `rt.Registry`）；`audio.go:180-182` 另有一枚 `SpawnCapture(registry *observe.Registry, …)` 收注册表，**产码零调用者**（U14） |
| U23 | `grep -n "ResidentNames\|ResidentBaseline\|OnDemandNames" internal/observe/goroutine.go` | 23:06 | `:39` `ResidentBaseline = 6` · `:42-44` `ResidentNames` 含 `"audio-capture"` · `:47` `OnDemandNames = []string{"kws-infer"}` · `:213` `type Registry` · `:262` `Spawn` · `:247` 附近 `var Default = NewRegistry()` |
| U24 | `grep -rn "SetState" --include=*.go internal/ball cmd/wisp \| grep -v _test` | 23:07 | 定义 `internal/ball/ball_windows.go:311`；产码调用者 `cmd/wisp/resident_approval_windows.go:443`（`stateForCardLevel(p.Level)`）与 `:492`（`StateSleeping`）〔起手读数，可能被 246-r2 改动〕⇒ **常驻那条腿目前只会把球按在 `Sleeping`／`Confirming` 两态** |
| U25 | `date "+%Y-%m-%d %H:%M:%S %z"` ＋ `git log --oneline -3` | 23:04:12 | `2026-09-30 23:04:12 +0800`；HEAD `f0707c7b`（本件起手时）；`git status --short` 显示 `246-r2` 正在动的 `cmd/wisp/**` **未出现在我的写点**（我零产码改动） |
| U26 | `grep -n "Armed" docs/specs/SPEC-08-ui-ball-panel.md` | 23:09:20 | `:61` 逐字 `\| \`Armed\` \| opacity 0.6，直径 44px，静态 \|`（**"静态"二字在冻结表里**）· `:93` D43 #5 的 spec 侧「`Sleeping` \| KWS 开启 \| `Armed`」· `:44-45` 编排者自己留下的**未裁契约-代码分歧**（Armed 直径 44 vs 代码 56） |
| U27 | `sed -n '76,100p' docs/specs/SPEC-05-agent-core.md` | 23:09:20 | 标题 `### 3.4 重试与失败语义（§14.2，失败必须可见且可区分）`。表里三行对本票有参照价值：`重试耗尽 → 明确报错，**不得静默降级** → Error + 面板给原因与建议`、`401/403 → 不重试 → Unconfigured + 引导配置面板`、`额度耗尽 → 区分「网络问题」与「账户问题」`。⚠ 这张表的**主语是 LLM/provider 调用**，不是音频设备——我把它的"分类＋可见"口径借过来用，借得对不对标在 ⑥-E12 |

⚠ 本票**没有**跑过 `go build`／`go vet`／`go test`／任何 `./...`（派单禁令，`246-r2` 在写同棵树）；U6–U10、U18 全部是**指名单包**的 `go list`，`internal/audio`／`internal/ball` 不在它的地界里，`cmd/wisp` 只做了 `go list -deps`（纯解析，不吃类型检查的一半）。

## ⑥ 我可能写错的条目（对抗我自己）

| # | 我写下的哪一条 | 为什么可能错 | 若错了会怎样 |
|---|---|---|---|
| E1 | 「`internal/audio` 非测试 importer＝0」 | U1 那把尺**认的是"文本里出现 import 路径"**，不是编译器意义上的 import。注释里写一句 `github.com/CarlosShao/wisp/internal/audio` 就会被计数 | 若将来某处只在注释里提到它，尺会说"有人调用了"。⇒ 编排者验收 AC#1 时**要换一把**：`go list -f '{{.ImportPath}} <- {{join .Deps " "}}'` 反查，或直接 `go list -deps` 前后名册差（我在 U10 用的就是后者，读数一致＝1 名，**这把不受注释污染**） |
| E2 | U14／U2 的空读数 | 我第一版 U14 用 `grep -rn … .`，被 `.scratch/wisp/probes/241/v1/mut-*.go`（30 行变异体 fixture）淹掉，看着像"有人调用"。改成 `internal cmd tools` 作用域才干净 | 如果我留着第一版，会误报"电平函数已有产码用户"。已重跑（23:07:08，退出码 1）。**但作用域收窄也带来反向风险**：我**没扫** `tools/` 之外的顶层 `.go`（顶层无产码，`ls` 只见 md/目录），也没扫 `scripts/spike`（票 AC#7 禁区，我故意不碰） |
| E3 | 「默认档下麦克风不会被打开」 | 我读的是**当下**（`cmd/wisp` 零 audio importer，U2/U14）。这句话在 `247-r1` 落地后**立刻失效**。而且 `mic_muted_default=true`（U12）挡住的前提是**接线的人真的把它传给 `WithStartMuted`**——门存在不等于门被装上（`gate.go:99-112` 只有包了门才不开内层） | 若编排者把它当成"永久安全"读，接完之后可能一双击就开麦。**这是本件最要紧的一条**，见 ⑦ 的 P1 |
| E4 | 「`cmd/wisp` 接 audio 只多 1 名依赖」 | U10 是**传递闭包差**，不是"包级直接边"差。包级边也确实只多 `cmd/wisp → internal/audio` 一枚（U8 名册里 audio 不在、U6 显示 audio 只依赖 observe，而 observe 已在 U8）。但如果落地时把电平消费者写进 `internal/ball`，或让 `internal/audio` 反过来 import `internal/ball`，**多出的是 1 名＋一条环** | 若我算错，票 246 引以为傲的「零新边」叙事会被这票打破却没人预警。⇒ ① 里我把两形的边**分开列**并各标尺 |
| E5 | 「第 4 步 `stop-audio` 的接缝现成可用」 | U20 证明 `RegisterShutdownHook(StepStopAudio, …)` 这条路**通**，但我**没执行过**（禁 build/test）。`shutdown_hooks.go:118` 的 `case StepStopAudio` 是我读到的分支，不是跑到的分支 | 若那枚 switch 分支其实漏了（例如只写进 roster 没落到 hooks），① 的方案 A 退出侧会落空。⇒ 交 `247-v1` 用一条测试钉：注册后 `StepRecord.Name == "stop-audio"` 且状态不是 skipped |
| E6 | 「`audio-capture` 名字合法、不超常驻 6」 | U23 表明确实在名册里。但 U22 显示 `wasapimic_windows.go:78` **写死 `observe.Default`**；如果 `247-r1` 用 `proc.WithRegistry(observe.NewRegistry())` 起进程（测试就这么干，见 `cmd/wisp/resident_approval_246_windows_test.go:133`〔起手读数〕），采集协程会记进**另一枚**注册表 | 后果不是崩，是**假绿**：常驻名册检查（`boot_windows.go:115` `ResidentOverBaseline`）看不见它，泄漏查不出。⇒ ⑦ P4 |
| E7 | 「`liquidDriven` 当前含 6 态、不含 Armed」 | 读的是 `liquid.go:62-70` 的 switch，机械、可复核。风险在另一头：我把 `Muted` 也归进"不含"，是靠 `:60` 的**注释**（"Sleeping/Armed/Muted … static"），不是靠 switch——注释不算凭据 | 若 `Muted` 的归属要落契约，我得能用 switch 的六个名字自证。已核对：`Listening/Thinking/Speaking/Conversation/Warm/Acting`，`Muted` 确实不在其中（同一份读到的代码，`liquid.go:64-66`） |
| E8 | 「Armed 那一格会撞票 74 的 C21 表镜像」 | U17 只证明 `tokens_test.go:155` 钉 Armed 的 **opacity**，而 `applyTo`（`liquid.go:312-324`）改的是 `LiquidLevel/LiquidAngle/SummonFlow/BorderAlpha`，**不改 opacity**。所以"会不会真红"我**没跑过**，不敢断言 | 我说的是"地界重叠"（同一枚状态的视觉真相），不是"测试会红"。措辞在 ④ 里已收窄，别把它读成后者 |
| E9 | 所有 `cmd/wisp` 行号（U5/U20/U24 与 ① 的文件清单） | `246-r2` 正在改 `run.go`／`resident_windows.go`／`resident_approval_windows.go`。我已逐条标〔起手读数〕，但**没预见到它会新增 audio 相关代码** | 若 246-r2 顺手挪了钩子注册点，① 的"要动哪几枚文件"名单会变。⇒ 派 `247-r1` 前该重跑 U20/U24 |
| E10 | 「`internal/speech` 只有 `doc.go` ⇒ ASR/KWS 一块没写」 | `ls` 只证**文件名**。若 ASR 藏在别的包（如 `internal/models`），我的推断过强 | 已交叉：票 247 现量 3 与项目记忆「240-c1 三分类」同口。且 U16 显示 KWS 事件零发射者，方向一致。**但唤醒词到底该落在哪枚包不是本票的结论**，我没往下普查（越界）。 |
| E11 | 「SPEC-08 §2.1 的 Armed 行写着静态」 | 我读的是 `docs/specs/SPEC-08-ui-ball-panel.md:61`（尺 U26，23:09:20 现跑）那行 `| Armed | opacity 0.6，直径 44px，静态 |`。"静态"是不是等于"永不接受电平"是我加的语义 | 若编排者认为"静态"只约束尺寸/opacity 不约束运动，④ 的结论会从"要人工批准"降为"直接可写"。⇒ 已作为 P2 交给裁 |
| E12 | 「SPEC-05 §3.4 要分类＋可见，可以借给音频降级当参照系」 | U27 读到的那张表，**主语是 LLM provider 调用**（429/5xx/401/额度），表里**没有一行**讲音频设备。派单叫我"参照系"用它，我照做，但**它是类比不是明文** | 若 `247-v1` 按"SPEC-05 §3.4 没写音频"判我这条参照无效，③ 的判据要改挂 `D42#2/#12`＋`observe.ClassAudioDevice` 那两处**明文**（`internal/audio/device.go:94-109`，尺 U19 邻域已读）。⇒ ③ 里我把两处明文的权重写在其上，SPEC-05 只当口径类比 |
| E13 | 「`cmd/wisp` 现在只会把球停在 `Sleeping`／`Confirming`」 | U24 只扫了 `SetState` 的产码调用者两枚〔起手读数，246-r2 正在改这枚文件〕。若 246-r2 新增别的状态来源，我这句话当场过期 | ① 方案 A 里"电平在 Sleeping 会被 `liquidDriven` 拒收"这一环就变了。⇒ 已标〔起手读数〕并写进 P7：派单前重跑 |

## ⑦ 判不动的地方（要编排者裁：甲／乙／不做）

| # | 判不动的是什么 | 甲 | 乙 | 不做 | 我的倾向（不算数） |
|---|---|---|---|---|---|
| P1 | **默认档下"这双腿谁都不开麦"这件事，要不要在本票里加一道硬门？**（U13：`voice.enabled`／`wake_word.enabled` 目前**没有任何产码读取处**；U12：真正会挡的是 `mic_muted_default=true` → `WithStartMuted`） | 接线时**必须**把 `c.Audio.MicMutedDefault` 传给 `WithStartMuted`，并把 `voice.enabled=false` 视为"根本不构造采集器"。零默认值改动，AC#4 可证 | 只在文档／日志里说明"默认静音"，代码靠门 `gate.go:99` 的自然行为挡 | 本票不碰，留到唤醒词票 | 甲。理由：乙的"自然行为"取决于接线人**是否包了门**，不包门就没有挡（② 详解） |
| P2 | **`Armed` 要不要收电平？**（`liquid.go:62` 的 switch 不含它；`docs/specs/SPEC-08-ui-ball-panel.md:61` 写着 Armed 静态（尺 U26）；`table.go:55` 进 Armed 的事件 `EvKwsEnabled` **产码零发射者**，U16） | 本票不动 `liquidDriven`，"接好了"的判据改挂 `Listening`（D43 #4 从 Sleeping 经 EvSummon 进 Listening，是今天真能到的态） | 本票补 Armed 那一格 ⇒ **同时改 SPEC-08 §2.1 那行**＝契约面，须台账落一枚 `A##` | 把 Armed 留给唤醒词票（KWS 有了发射者再说） | 甲＋不做（本票不补，Armed 明确归口唤醒词票）。因为**乙要先有人改冻结表**，而那不该由写码腿顺手做 |
| P3 | **"球随声音呼吸"今天看不见的另一道闸：`prototypeVisuals` 默认关**（U15：产码只有 `cmd/balldebug:122` 打开；`liquid_windows.go:52` 在它关的时候把 `SetAudioLevel` **直接 return**） | 本票在常驻腿里调 `ball.EnablePrototypeVisuals(true)`，这样 AC#2 才看得见 | 本票不碰，AC#2 改用**读数型证据**（`Ball` 内部 `liqRaw`/`PrototypeVisualsEnabled()`）而不是屏幕 | 等票 68 AC#2（`statevisual.go:93` 写明"翻默认值"归它） | 这一条**超出普查权限**：它决定 AC#2 的取证形态，而 AC 框归编排者。交你裁，我只把行号摆出 |
| P4 | **采集协程记进哪枚注册表？**（U22：`wasapimic_windows.go:78` 写死 `observe.Default`；`audio.go:180` 那枚收 registry 的 `SpawnCapture` 零调用者；常驻腿用 `rt.Registry`，`resident_ball_windows.go:111`〔起手读数〕） | 把 `WASAPIMicrophone.Start` 改成走传入的 registry（动 `internal/audio` 产码，票 241 的地界？需你定） | 常驻腿只用默认注册表启动（不给 `WithRegistry`），并在测试里接受这一处不同 | 现在不动，登记台账 | 甲，但**这是"改已交付的票 13 产物"**，不是我该定的。见 E6 的假绿后果 |
| P5 | **降级要不要走 D43 #14 那条边？**（`table.go:97`：`Listening` + `EvAudioDeviceLost` → `Error`。但采集失败发生在**还没进 Listening 之前**，今天球停在 `Sleeping`，U24 常驻腿只会 `SetState(Sleeping/Confirming)`） | 设备失败**不推状态**，只把 `observe.ClassAudioDevice` 的分类＋guidance 文案**响亮地**打印／落日志（票 128 口径：只有没有数据根才拒绝启动，不扩大） | 给 `EvAudioDeviceLost` 找一枚 Sleeping 起点的合法边——但我翻表**没找到**（#14 的 From 只有 Listening），那要改 D43 转移表＝契约面 | 只 `slog.Error`，用户可见面留白 | 甲。**乙一律不要**：改 D43 表只能编排者在台账落 `A##` |
| P6 | **AC#3 的"只过一枚 `float32`"由谁保证？** 电平函数在 `internal/audio`（`level.go:118 FrameLevel → (float32, error)`），消费口在 `internal/ball`（`liquid_windows.go:42`）。中间那枚读 channel 的人放哪？ | 放 `cmd/wisp`（采集腿自己 `FrameLevel` 再 `SetAudioLevel`）＝零新边，见 ① 方案 A | 新造一枚 `internal/audio/levelfeed` 小包 → 新增两条包级边 | 让 `internal/audio` import `internal/ball` 直接推 → **audio→ball 这条边今天不存在**（U6/U7），且把渲染拖进采集的生命周期 | 甲。乙要多 2 枚直接依赖＋传递名册变化，我用 U10 同法可量，别凭感觉 |
| P7 | **票面 AC#0 说"由编排者裁后再派 `247-r1`"，但排程又要求 247-r1 排在 246-r2 之后**——若 246-r2 的终态把 `cmd/wisp` 的文件名／钩子注册点改了，我的 ① 文件清单要不要重抽？ | 派单前重跑 U20/U24，以 246-r2 终态为准（我建议） | 直接用本件清单 | 让 247-r1 自己摸 | 甲（E9 的代价可控） |

## ① 落点（填写中）

## ② 隐私闸门（填写中）

## ③ 降级（填写中）

## ④ 第六环（形状）（填写中）
