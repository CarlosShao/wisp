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
| P7 | **AC 框归你，本件只报盘上事实**：246-r2 若改了 `cmd/wisp` 的文件名／钩子注册点，我的 ① 文件清单要不要重抽？ | 派单前重跑 U20/U24，以 246-r2 终态为准（我建议） | 直接用本件清单 | 让 247-r1 自己摸 | 甲（E9 的代价可控） |
| P8 | **第二枚协程（读 `buf` → 算电平 → `SetAudioLevel`）在 D38(b) 名册里没有槽**（`observe/goroutine.go:42-44` 的 `ResidentNames` 固定 6 名；`:65-82` 对名册外名字返 `CategoryUnknown`；`:270-273` 立刻 `slog.Warn(… leak symptom)`；`internal/audio/level.go:85-88` 逐字写着 "D38b leaves no spare resident slot for a level-only goroutine"） | 电平在**已有的 `audio-capture` 协程内**算完再调 `SetAudioLevel` ⇒ 名册零膨胀。代价：把一次 `FrameLevel`＋一枚跨包调用带进 D38(a) 的 pinned 线程，而 `wasapimic_windows.go:149-152` 的循环自述只列 open/wait/drain/resample/push/hotplug | 新起协程并**复用** `"audio-capture"` 这个名字 ⇒ 仪器不红，但 `Resident` 计数变 2，`PLAN.md:2831` 那句"常驻 6 个"的语义被稀释 | 本票只接**生产者半**（AC#1 的"非测试 importer ≥1"由"采集有人 Start"兑现），消费协程等 ASR/KWS 票——那时读 `buf` 的人本来就必须存在（`audio.go:23-33` 的 C8 契约就是这么设计的） | 甲在语义上最干净但碰 D38(a) 的明文清单；丙最省且不改契约。**我倾向丙**，因为甲／乙都要动一枚已冻结的口径（D38(a) 线程职责 或 D38(b) 名册计数），而那是契约面 |

## ① 落点：采集协程归常驻腿还是归 `internal/audio` 自起

起手钟点 `23:04–23:11 +0800`。凡 `cmd/wisp` 的行号＝〔起手读数，可能被 246-r2 改动〕。

### 现状（三条边都靠尺，不靠印象）

| 事实 | 尺 | 读数 |
|---|---|---|
| `internal/audio` 的直接依赖里**只有 `internal/observe` 一枚 wisp 包**，不含 `internal/ball` | U6 | `context encoding/binary fmt …/internal/observe golang.org/x/sys/windows log/slog math os runtime strconv sync sync/atomic syscall time unsafe` |
| `internal/ball` 的直接依赖是 `observe`＋`statemachine`，**不含 `internal/audio`** | U7 | 同上格式 |
| `cmd/wisp` 已 import `internal/ball`／`config`／`observe`／`proc`／`statemachine`，**不含 `internal/audio`** | U8 | 名册 41 项，逐条看过 |
| 全仓产码里**没有一枚**构造采集器、门、有界 channel 或调用电平函数 | U2 + U14 | 两次都空（grep 退出码 1） |
| `internal/ball` 的非测试 importer 有 4 枚（含常驻腿） | U5 | `cmd/balldebug/main.go` · `cmd/wisp/resident_approval_windows.go` · `cmd/wisp/resident_ball_228_test.go` · `cmd/wisp/resident_ball_windows.go` |

### 形甲：采集协程归常驻腿（`cmd/wisp`，票 246 那枚进程）

**要动的文件**（起手读数，23:10:51 现跑 `ls cmd/wisp/*.go | grep -v _test` 得到 27 枚，下面只列要碰的）：

1. **新增** `cmd/wisp/resident_audio_windows.go`（`//go:build windows`）：建 `NewWASAPIMicrophone()`（`internal/audio/wasapimic_windows.go:52`）→ 包 `NewHalfDuplexGate(inner, PathT, WithStartMuted(...))`（`gate.go:84`／`:57`）→ `NewBoundedFrames()`（`audio.go:65`）→ 读 buf、`FrameLevel`（`level.go:118`）→ `b.SetAudioLevel(float32)`（`internal/ball/liquid_windows.go:42`）→ 一枚响亮 verdict。
2. `cmd/wisp/resident_windows.go`：在 `startResidentBall(rt.Registry, …)`（`:135`〔起手读数〕）之后接一段，并把退出钩子挂上（见下面"谁 join"）。
3. **一笔藏不住的额外代价**：常驻腿**今天完全不读 config**。尺：`grep -rln "wisp/internal/config" cmd/wisp/*.go`（23:10:37，排 `_test`）命中 `config_reload.go` · `models.go` · `panel_inbound.go` · `providers.go` · `run.go`——**`resident_windows.go`／`resident_ball_windows.go` 一枚都不在**。所以 `audio.mic_muted_default` 想生效，得先把 config 读进这条腿（或从 `run.go` 把已加载的配置传下去）。⇒ 见 ② 与 ⑦-P1，这是甲形最贵的一笔，不是可选装饰。
4. **不需要**非 Windows 桩：`cmd/wisp/resident_other.go` 已 `os.Exit(2)` 拒绝起常驻（逐字 "the resident process requires Windows … refusing to start."）；`internal/audio/wasapi_other.go` 提供 `WASAPIMicrophone` 的非 Windows 占位（`Start` 必失败），所以类型面在两个 GOOS 都存在。**但**占位体没有 `Err()`／`ThreadID()`／`Endpoints()`（我逐行看了那份文件）⇒ 消费文件若调这三枚，**必须** windows-tagged。

**新增的包级依赖边**：**恰好 1 枚**，`cmd/wisp → internal/audio`。
传递名册增量也用尺算了，不是估的：U10 `comm -23 deps-audio.txt deps-cmdwisp.txt` ⇒ **只有 1 名**（`github.com/CarlosShao/wisp/internal/audio`）；`internal/audio` 其余 143 名传递依赖（`internal/observe`、`x/sys/windows`、标准库）**已经在 `cmd/wisp` 的 276 一名册里**。
⇒ 这枚进程今天的"零新边"叙事（`resident_windows.go:116-121` 注释逐字 "Zero new package-level dependency edges"〔起手读数〕）会被本票**合法地**打破一枚，票面要提前写明，别让验收腿当成越界。

**协程 owner 走不走现成的 `observe.Root`／`rt.Registry`**：

- 采集协程**已经**有 owner：`wasapimic_windows.go:78` ＝ `observe.Default.Spawn("audio-capture", "audio", nil, …)`（U22 现跑）。⚠ 它**写死 `observe.Default`**，而不是常驻腿在用的 `rt.Registry`（`internal/proc/boot_windows.go:79`：没给 `WithRegistry` 时 `rt.Registry = observe.Default`；给了就是另一枚，`resident_ball_windows.go:103-111`〔起手读数〕的注释明说"reg is the runtime's registry, not observe.Default"）。`internal/audio/audio.go:180 SpawnCapture(registry *observe.Registry, run …)` 那枚**收注册表**的入口产码零调用者（U14）。⇒ 两枚注册表并存时，`boot_windows.go:115` 的 `ResidentOverBaseline` 自检看不见采集协程（假绿，⑥-E6，交裁 ⑦-P4）。
- 名字合法、不越名册：`observe/goroutine.go:42-44` 的 `ResidentNames` 含 `"audio-capture"`，`ResidentBaseline = 6`（`:39`）（U23）。
- ⛔ 裸 `go func(`：仪器确实会点红——`tools/d22scan/main.go:700-707`（"bare \`go func(\` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)"，`:700` 那行注释说明它不只匹配 FuncLit）。`observe.Registry.run`（`goroutine.go:286`）是全仓唯一被许可的 `go` 语句（`:20-24` 注释逐字 "This file contains the ONLY sanctioned `go` statement of the codebase"）。
- **电平消费侧那枚"读 buf"的协程在 D38(b) 名册里没有槽**（本件新发现，硬约束）：`ClassifyGoroutine`（`goroutine.go:65-82`，`return CategoryUnknown` 在 `:81`）对名册外名字返回 `CategoryUnknown`，`Spawn`（`:262` 定义）随即在 `:270-273` 打 `slog.Warn("goroutine outside the D38 roster (leak symptom)")`，`RosterReport`（`:405`）把它列进 `Unknown`（`ResidentOverBaseline` 判据在 `:423` 附近，逐字 `rep.Resident > ResidentBaseline`）。票 241 的作者已经预见这件事并把它写进代码：`internal/audio/level.go:85-88` 逐字 —— *"a level reader is called from the consumer's own thread, and D38b leaves no spare resident slot for a level-only goroutine"*。
  ⇒ 三条出路（我不裁，交 ⑦-P8）：① 电平在**已有的 `audio-capture` 协程内**算完再 `SetAudioLevel`（名册零膨胀，代价是把一次跨包调用带进 D38(a) 的 pinned 线程，而 `wasapimic_windows.go:149-152` 的循环注释逐字只列 open/wait/drain/resample/push/hotplug）；② 新起协程并复用 `"audio-capture"` 这个名字（仪器不红，但 `Resident` 计数变 2，`PLAN.md:2831` 的"常驻 6 个"语义被稀释）；③ 本票只接生产者半（AC#1），消费协程等 ASR/KWS 票——那时读 buf 的人本来就必须存在。

**D38(e) 十步一字不动的前提下，退出时谁 join 它**：第 4 步 `stop-audio` 就是 join 点，接缝**现成**：

- 十步在代码里与 `PLAN.md:2855-2865`（第 1 步在 `:2855`、第 4 步在 `:2859`、第 10 步在 `:2865`，23:15:25 现跑 `grep -n`）一字对齐：`internal/proc/shutdown.go:11-31`（注释）＋ `:36-47`（`StepStopAudio` 是第 4 枚）＋ `:49-60`（名字 `stop-audio`）＋ `:163` `run(StepStopAudio, hooks.StopAudio, 0)`（U21）。
- 注册面已经通了：`proc.Runtime.RegisterShutdownHook(step, hook)` 在 `boot_windows.go:188`，`shutdown_hooks.go:51` 把 `StepStopAudio` 列进可挂钩名册，`:118-119` 落到 `hooks.StopAudio`（U20）。本仓**唯一**产码注册者目前是 `cmd/wisp/resident_windows.go:156` 的 `proc.StepCancelTasks`（票 246 AC#4）〔起手读数〕。
- join 的动作本来就在 `Stop()` 里：`wasapimic_windows.go:94-115` —— `cancel()` 之后 `select { case <-done: return nil; case <-time.After(tm.Remaining()): slog.Error("audio capture thread did not exit within 2s (abandoned wait)") }`，上限用 `observe.NewTimeout(2 * time.Second)`（D42#9 单调，**不是**墙钟差）。⇒ 钩子体只需要 `gate.Stop()`／`mic.Stop()`，**不需要新造第 11 步**。
- ⚠ 诚实边界：我**没有跑到**那枚 switch 分支（禁 build/test），`shutdown_hooks.go:118` 是我读到的分支不是跑到的分支（⑥-E5）。⇒ 交 `247-v1` 钉一条：注册后 `StepRecord.Name == "stop-audio"` 且记录**不是** skipped。

### 形乙：`internal/audio` 自己起协程并推到球

乙分两支，代价不同，都量过：

- **乙-1 直推**（audio 里持有／调用 `*ball.Ball`）：新增**直接边** `internal/audio → internal/ball`，今天不存在（U6/U7）。传递名册**多 2 名**：U28 `comm -23 deps-ball.txt deps-audio.txt` ⇒ `internal/ball`、`internal/statemachine`。两个后果：
  1. 采集层从此依赖渲染层＋D43 状态机；
  2. `internal/audio` 在 `wasapi_other.go:6-8` 的明文承诺「the package must still compile on other GOOS values so accidental cross-packages references fail loudly, not silently」**当场破**——`internal/ball` 的 `ball_windows.go`／`d2d_windows.go`／`sta_windows.go` 全是 windows-only。
- **乙-2 注入**（audio 自己起协程、算电平，出口是传入的 `func(float32)` 或 audio 侧声明的 sink 接口）：**边数与甲完全相同**（还是 `cmd/wisp → internal/audio` 一枚；U29 `comm -13 deps-ball.txt deps-audio.txt` 显示 ball 那侧只多 `internal/audio` 一名，说明 audio 不需要碰 ball 的任何依赖）。⚠ 关键：**乙不比甲少动 `cmd/wisp`**——`mic.Start(ctx, buf)` 总得有人调（`wasapimic_windows.go:69` 是同步等开成功的，见 ②），第 4 步钩子也总得有人注册。乙只是把"第二枚协程"的所有权从常驻腿搬进 `internal/audio`，并把 ⑦-P8 那条名册约束从 cmd 侧搬进 audio 侧。

### 我的建议（仍要编排者裁，我只摆代价）

**甲**。三条，全带尺：

1. **边**：甲 1 枚（U10）／乙-1 1＋1 枚且带渲染依赖（U28）／乙-2 与甲同（U29）。
2. **AC#3 只过一枚 float32**：甲天然满足——`FrameLevel` 的返回类型是 `(float32, error)`（`level.go:118`），`SetAudioLevel` 的入参是 `float32`（`liquid_windows.go:42`）；`[]byte` 只在 `internal/audio` 内部与那枚有界 channel 里活（`audio.go:29 Start(ctx, buf chan<- []byte)`）。⛔ 任何"把 samples 递过界"的写法在签名面上就红（能力型判据，不是词面型）。
3. **dt 不该由新协程算**：球自己拥有节律时钟（`liquid_windows.go:77-85 motionDt`），新代码里若写 `time.Now().Sub(...)` 会被 `d22scan` 的 `wallclock-timeout`（`main.go:144` 那条 regexp `\.Sub\(time\.Now\(\)\)`，报在 `:764-766`）点红 ⇒ 甲形里读数的人只交 `float32`，不交时间。

## ② 隐私闸门：接完之后麦克风在什么条件下真的打开

### 三枚开关的真实字段名与默认值（不引派单那句话，只引 `internal/config`）

| 语义名 | Go 字段 | `toml` 路径 | 默认值 | 出处（现读凭据） |
|---|---|---|---|---|
| `voice.enabled` | `VoiceSection.Enabled` | `voice.enabled` | **`true`** | `internal/config/schema.go:244`（struct 头）＋ `:245` `Enabled bool \`toml:"enabled" default:"true"\`` ；section 挂点在 `:114 Voice VoiceSection \`toml:"voice"\`` |
| `wake_word.enabled` | `WakeWord.Enabled` | `voice.wake_word.enabled` | **`false`** | `schema.go:196`（struct 头）＋ `:197` `Enabled bool \`toml:"enabled" default:"false"\`` |
| **`audio.mic_muted_default`**（本件发现的**第三枚**，也是唯一一枚真挡得住麦的） | `AudioSection.MicMutedDefault` | `audio.mic_muted_default` | **`true`** | `schema.go:269`（struct 头）＋ `:276-277` `// MicMutedDefault starts every session muted.` / `MicMutedDefault bool \`toml:"mic_muted_default" default:"true"\`` |

默认值不是文档口径而是**代码事实**：`internal/config/defaults.go:54-60 NewDefaults()` 走 `applyDefaults`（`:64`），逐叶读 `default:"…"` 标签（`:73`），`reflect` 的 Bool 分支在 `:96`。⇒ 这三枚默认值由标签单点决定（D36 规则 3），**我没有动过任何一格**（本件零产码改动）。

### 这两枚开关今天**有没有**被任何产码用来决定开麦

尺 U13（`grep -rn "Voice\.Enabled\|WakeWord\.Enabled" --include=*.go . | grep -v _test`，23:06）⇒ 产码命中**只有 3 处，全在 `internal/config/manager.go`**：`:370`（reload 分层比较）、`:384`、`:385`（把新值抄回活配置）。**用途是"改了要不要重载"，不是"要不要开麦"**。⇒ 派单说的"方向相反"这一层成立，但更准的说法是：**这两枚今天都还没有方向**，因为没有任何消费者读它们做设备决策。

### 具体形状：**双击 `wisp.exe`、什么都没配的那台机器上，这条腿会不会自己开始采音？**

**不会。但不是被默认值挡住的——是因为今天根本没有那条路径。** 三重，每重都带尺：

1. **常驻腿不 import `internal/audio`**：U2（`grep -rl … internal cmd tools` 空，退出码 1）＋ U8（`cmd/wisp` 名册里没有它）。⇒ 那条腿上**没有能开麦的对象**可被构造。
2. **常驻腿连 config 都不读**：`grep -rln "wisp/internal/config" cmd/wisp/*.go` 排测试后的命中是 `config_reload.go`/`models.go`/`panel_inbound.go`/`providers.go`/`run.go`（23:10:37〔起手读数〕），`resident_windows.go` 与 `resident_ball_windows.go` **不在其中**。⇒ `voice.enabled` 的值在"一双击"那条腿上**从未被读取**。
3. **唯一今天喂球的产码是合成数**：U3 ⇒ `cmd/balldebug/main.go:418`／`:477`，而 balldebug 不 import audio（U2/U5）⇒ 与设备无关。

### 那"接完之后"呢——决定权在**接线人装没装门**，不在默认值

| 接法 | 设备会不会在启动时打开 | 决定它的行 |
|---|---|---|
| **装了门**，并把 `mic_muted_default` 的值传进 `WithStartMuted`（`internal/audio/gate.go:57`） | **不会**。`HalfDuplexGate.Start`（`gate.go:99-112`）的明文是 `if g.effectiveOpen() { g.open = g.openInnerLocked() }`；`effectiveOpen()`（`gate.go:234-241`）第一行 `if g.muted { return false }`；`openInnerLocked()`（`:247-252`）根本不跑 ⇒ 内层 `mic.Start` 不跑 ⇒ `wasapi_windows.go:210 Open` 不跑 | `gate.go:108`（那个 `if`）× `gate.go:235-237`（muted 先返回）× `schema.go:277`（默认 true） |
| **没装门**，直接 `mic.Start(ctx, buf)` | **会，而且是一双击就开**。`WASAPIMicrophone.Start`（`wasapimic_windows.go:69-88`）在 `:78` 起 pinned 协程后 `:88 return <-started`——**同步等开成功**；`started` 只有在 `run` 里 `openCurrent()` 成功之后才送 nil（送 nil 在 `wasapimic_windows.go:182`，两处失败早送 `started <- err` 在 `:170`／`:178`） | `wasapimic_windows.go:88` ＋ `:182` |

⇒ **一句话结论**：挡住采音的既不是 `voice.enabled=true` 也不是 `wake_word.enabled=false`，而是"没人接线"这件事本身；接线时唯一可用的默认是 `audio.mic_muted_default=true`（`schema.go:277`），**而它只在门被装上时才生效**。这就是 ⑦-P1 要裁的东西（甲＝强制 `WithStartMuted(c.Audio.MicMutedDefault)` 且 `voice.enabled=false` 时不构造采集器）。

### "如果改了默认值会怎样"（本票零改动，只讲清）

- `audio.mic_muted_default` → `false`：`WithStartMuted(false)` ⇒ `effectiveOpen()` 为真 ⇒ **装了门的接法也会在启动时开设备**。这是三枚里唯一**直接改变开麦行为**的默认。
- `voice.enabled` → `false`：**今天零后果**（U13：无决策读取处）。只有当接线人按 ⑦-P1 甲把它变成"不构造采集器"的条件后，它才开始挡；也就是说它的约束力是**接线造成的，不是默认值自带的**。
- `voice.wake_word.enabled` → `true`：**今天同样零后果**，且 KWS 落地属唤醒词票（票面 AC#7 明列"唤醒词／ASR／TTS 都不属于本票"）。附带一条同向证据：进 `Armed` 的 `EvKwsEnabled` 产码零发射者（U16）。
- ⚠ 别把"零后果"读成"无关"：`voice` 与 `wake_word` 的 `Enabled` 都在 reload-tier（`schema.go:240-243` 注释＋`manager.go:384-385`），改了会走 `rep.Reload = append(rep.Reload, "voice")`（`manager.go:412`）。

## ③ 降级：设备被占／无权限／无麦／热插拔现在返回什么

下面每条都是**用例的断言本体**（`sed` 逐行读），不是用例名字。文件：`internal/audio/hotplug_test.go`。

| 形态 | 用例（行号） | 断言本体（逐行读到的判据） | 错误落在哪一行 |
|---|---|---|---|
| **设备被占**（别的程序独占） | `TestOpenOccupiedAndPermissionDenied` `hotplug_test.go:378`（注释 `:376-377`） | ① `mic.Start(ctx, buf)` **必须返回 error**（`:403-406`，逐字 "Start must fail when the device cannot be opened (D42#12)"）；② `observe.ClassOf(err) == observe.ClassAudioDevice`（`:408-411`）；③ 错误串**必须含设备名** `"Busy Mic"`（`:413-415`）；④ 必须含 `inUseGuidance`（`:417-419`）；⑤ `mic.Stats().LastError != ""`（`:420-422`，逐字 "must record the open failure"） | `device.go:100-101`（`hrAUDCLNTDeviceInUse = 0x8889000A`，码在 `:77`）→ 文案 `device.go:89 inUseGuidance`（逐字 "capture device is held in exclusive mode by another application; close that application or pick another device"） |
| **无权限** | 同一枚用例的第二个 case（`:392-395`） | 同上五条，guidance 换 `privacyGuidance` | `device.go:98-99`（`hrEAccessDenied = 0x80070005`，码在 `:72`）→ 文案 `device.go:86 privacyGuidance`（逐字给出 Windows 设置路径 "Settings > Privacy & security > Microphone"） |
| **热插拔：默认设备换掉** | `TestHotplugReenumerateOnce` `:199` | 初始 `openCount()==1`（`:222-224`）；帧先流出（`:225`）；`watcher.fireChange()` 后帧**继续流**（`:229-231`）；`openCount()` **恰为 2**＝只重举一次（`:232-235`）；第二枚 open 用的是**新默认设备 dev2**（`:236-242`）；`Stats().Reopens == 1`（`:243-245`）；再 fire 一次 ⇒ `Reopens==2 && openCount==3`（`:246-251`）；收尾 `mic.Err()` 必须为 nil（`:254-256`） | 循环体 `wasapimic_windows.go:194-203`（`case <-m.hotplug:` 在 `:194`，重开成功 ⇒ 换流＋换重采样器 `:200-203`）＋ `reopenOnce` `:270-281`（成功 `slog.Info` 且 `meter.reopened()`） |
| **热插拔：重开失败（陈旧句柄）** | `TestHotplugStaleHandleFailsLoudly` `:260` | 终态错误经 `mic.Err()` 出来：分类 `ClassAudioDevice`（`:303-306`）、**含设备名** `"Dying Mic"`（`:308-310`）、含 `inUseGuidance`（`:311-313`）；`Stats().LastError` 非空（`:320-322`）；`openCount()==2`＝只试一次（`:323-325`）；**并且 buf 从此不再出帧**（`:328-340`，注释逐字 "no silent stale flow"） | `wasapimic_windows.go:195-199`（重开失败 ⇒ `meter.fail` ＋ `postErr` ＋ `return`）；`handleStreamError` `:240-256`；`postErr` `:283-287`（非阻塞，永不卡循环） |
| **热插拔：重举本身就失败** | `TestHotplugReenumerateFailureNamesDevice` `:334` | 终态错误**必须含上一枚设备名** `"Unplugged Mic"`（`:361-364`）＋分类 `ClassAudioDevice`（`:365-368`） | `reopenOnce` `wasapimic_windows.go:271-275`，错误串逐字 "hotplug reopen failed; keeping no stale handle (previous device …)" |
| **无麦／枚举失败** | （**没有专属用例**，我如实说） | `openCurrent()`（`wasapimic_windows.go:258-264`）把 `DefaultCaptureDevice()` 的错误 `observe.Wrap(ClassAudioDevice, err, "capture device enumeration failed")` 包成终态；枚举层的原始错误在 `mmdevice_windows.go:86-87`（逐字 `"GetDefaultAudioEndpoint(flow=…) failed (HRESULT)"`）与 `:93-94`（`GetId failed`） | ⚠ **一处文案缺口**：设备不存在时 `GetDefaultAudioEndpoint` 返 `hrAUDCLNTNotInit = 0x88890001`（码在 `device.go:74`），而 `DeviceError` 的 switch（`device.go:97-110`）里**没有 NotInit 这一枚 case** ⇒ 走 `:109` 的 default，只给码不给建议。**这条是我读 switch 读出来的，没跑到**（⑥-E5 同源）；且它走的是 `Wrap` 而非 `DeviceError`，所以实际串形要以真机读数定，我不替它编 |

### 本进程该怎么"照跑并把损失说响亮"——仓里有现成参照系，不用新造

| 参照 | 行号 | 它的形状（读到的，不是概括的） |
|---|---|---|
| **球起不来 ≠ 停止启动** | `cmd/wisp/resident_ball_windows.go:139-145`〔起手读数〕 | `ball.New` 失败 ⇒ 设 `rb.verdict`（逐字 "this process has NO floating ball window: %v"）＋ `slog.Error` ＋ `fmt.Printf`，然后 **`return rb`，不是 `os.Exit`**；同文件 `:125-126` 注释逐字 "Errors are returned as a verdict, never as a failure to boot"；调用侧 `resident_windows.go:132-134` 注释逐字 "A ball that will not come up does not stop the boot" |
| **日志目录开不了 ≠ 停止启动** | `cmd/wisp/resident_windows.go:62-69`〔起手读数〕 | `slog` 没装上就 `fmt.Fprintf(os.Stderr, …)`，注释逐字 "Loud, and it does not stop the app" |
| **只有"没有数据根"才拒绝启动**（票 128 的口径，不许扩大） | `internal/proc/boot_windows.go:80-84` ＋ `cmd/wisp/resident_windows.go:45-48`〔起手读数〕 | `DefaultLayout(env)` 失败 ⇒ `return nil, err`；`runResident` 拿到 err ⇒ `os.Exit(1)`。同族的硬拒只有 Job Object（`boot_windows.go:91-95`）与单实例互斥（`:98-107`）。**音频设备不在这一列** |
| （票 128 本身） | `.scratch/wisp/issues/128-resolvedatadir-…md` | 标题与 AC#2/AC#3 是"refusal"（拒绝启动）口径——我**没有**在本件把它读成"任何失败都可拒起"，反向引用见上面第三行 |

⇒ 建议形态（**零默认值改动、零拒绝启动扩大**）：采集侧任一形态失败 ⇒ ① 球继续停在 `Sleeping`、任务管线继续跑；② verdict／日志把 `observe.ClassAudioDevice` 的分类**和 `DeviceError` 已经编码进串的 guidance 原文**一起说出来，并具名是哪一形（被占／无权限／无设备／热插拔失效）——因为文案已经在 `device.go:86/89/103/105/107` 五处备好了，不需要新写；③ 顺带把 `mic.Stats()`（`wasapimic_windows.go:133`）的 `Reopens`／`FramesDropped`／`LastError` 说响亮（D38d 的"静默丢帧禁止"口径在 `audio.go:80-89`＋`audio.go:136-142`）。

⛔ **不能走状态机那条边**：`internal/statemachine/table.go:97` 的 D43 #14 是 `From: StateListening, Event: EvAudioDeviceLost, To: StateError`——`From` **只有 `Listening`**。而常驻腿今天到达的状态只有 `Sleeping` 与 `Confirming`（U24：`resident_approval_windows.go:443` `b.SetState(stateForCardLevel(p.Level))` 与 `:492` `b.SetState(statemachine.StateSleeping)`〔起手读数〕）。启动期采集失败**没有合法边**可走；硬造一枚 `EvAudioDeviceLost` 的 Sleeping 起点＝改 D43 转移表＝契约面。⇒ ⑦-P5 交裁，并明确建议"乙一律不要"。

⚠ SPEC-05 §3.4 的"分类＋可见"口径我是**借**来的（`docs/specs/SPEC-05-agent-core.md:76-88`，尺 U27 现跑：那张表的主语是 LLM provider 的 429/5xx/401/额度，**没有一行讲音频设备**）。借得对不对写进 ⑥-E12；本节的硬凭据一律挂在 `D42#2/#12` ＋ `observe.ClassAudioDevice` ＋ `device.go` 那五处文案上。

## ④ 第六环（形状）：`liquidDriven` 当前实际含哪几态

### 逐行读数（`internal/ball/liquid.go:62-70`，读的是 switch 里的六个标识符本身）

```
:62 func liquidDriven(s statemachine.State) bool {
:63   switch s {
:64     case statemachine.StateListening, statemachine.StateThinking,     ← 含 ①Listening ②Thinking
:65       statemachine.StateSpeaking, statemachine.StateConversation,     ← 含 ③Speaking ④Conversation
:66       statemachine.StateWarm, statemachine.StateActing:               ← 含 ⑤Warm ⑥Acting
:67       return true
:68   }
:69   return false                                                        ← 其余全部为 false
:70 }
```

⇒ **当前实际包含的正好 6 枚**：`Listening` · `Thinking` · `Speaking` · `Conversation` · `Warm` · `Acting`。
⇒ **不含**：`Armed`、`Muted`、`Sleeping`、`Confirming`、`AwaitingApproval`、`Settling`、`Downloading`、`Unconfigured`、`NoNetwork`、`Error`、`Queued`、`Stuck`、`WatchdogAlert`、`FirstRun`（状态全名册 **20 枚**，逐行读自 `internal/statemachine/states.go:10-30`，同文件 `:8-9` 注释逐字 "The 20 states, names exactly as in D43 / SPEC-08"；23:14:27 现跑）。
⇒ 同一份代码里的自述（`liquid.go:59-61`）逐字：*"liquidDriven reports whether a state receives audio-driven motion. Only session states do; Sleeping/Armed/Muted and the failure helpers are static glass frames by D32 discipline."* ——**"Armed 不进驱动表"是被当成 D32 纪律写下来的，不像是漏了**。这一句是本节判读的支点（但注释不算凭据，所以我另找三条实码，见下）。

调用点三处（尺：`grep -rn "liquidDriven" internal/ball cmd | grep -v _test`）：`liquid_windows.go:56`（`applyLevel` 的第一道闸）、`liquid_windows.go:162/:164`（`motionStateChanged` 决定 beginSession/endSession）、`liquid.go:119`（`enterState` 决定 border 权威）。

### 判"接好了"要不要看 `Armed` 收到电平？我的判读：**本票不该以它为判据**

三条，每条都带尺：

1. **`Armed` 今天不可达**。进入它的唯一合法边是 `internal/statemachine/table.go:55`（D43 #5：`From: StateSleeping, Event: EvKwsEnabled, To: StateArmed`；spec 侧同一行在 `docs/specs/SPEC-08-ui-ball-panel.md:93`）。而 `EvKwsEnabled` 在产码里**零发射者**（U16：命中只有 `events.go:18` 定义与 `table.go:55` 表格行）；常驻腿只 `SetState` 到 `Sleeping`／`Confirming`（U24）。⇒ "看 Armed 收到电平"这七个字**没有可执行的取证路径**，除非实现者另造一枚状态入口——那是伪造，不是接线。
2. **改表＝契约面**。`docs/specs/SPEC-08-ui-ball-panel.md:61`（U26 现跑）逐字 `| Armed | opacity 0.6，直径 44px，静态 |`。把 Armed 放进 `liquidDriven` 而不同批改这行"静态"，代码与冻结表就互相打脸；而同批改 `docs/specs/**` 是本票 **AC#7 明令退回**的路径。⇒ 这条路只能在编排者台账落一枚 `A##` 之后走（票面禁区第一条）。
3. **地界与票 62／74 重叠**（如实收窄，别说过头）：`internal/ball/tokens_test.go:141` 注释逐字点「SPEC-08 §2.1 anchors. The Sleeping size/opacity, **Armed**, Muted and the …」，`:155-156` 断言 `VisualFor(p, 56, StateArmed, 0).Opacity == 0.6`，`:298`／`:318` 把 Armed 列进静态态名册（U17）。⚠ 但 `applyTo`（`liquid.go:312-324`）只写 `LiquidLevel`/`LiquidAngle`/`SummonFlow`/`BorderAlpha`，**不写 `Opacity`** ⇒ 那枚断言**未必真红**（⑥-E8）。我说的是**同一枚状态的视觉真相归票 62／74**这个地界事实，不是"必然失败"。

### 真要"补 Armed 那一格"，要动的文件（只列，不建议本票做）

`internal/ball/liquid.go:62-67`（switch 加名）＋ **必须同批** `docs/specs/SPEC-08-ui-ball-panel.md:61` 的"静态"二字＋ 复核 `internal/ball/tokens_test.go:298/:318` 的名册与 `internal/ball/statevisual.go:158`／`:288` 的 Armed 映射（两处 `case statemachine.StateArmed:`，U16 读到）＋ 另给 `EvKwsEnabled` 找一枚发射者（那是唤醒词票的地界，AC#7 禁入）。

### 本节最重要的产出（给编排者裁 AC#2 用）

即使本票把电平接进常驻腿，**"球动了"今天仍然不可见**，两道独立的闸都在代码里、都不在本票的地界：

1. `prototypeVisuals` **默认关**：`internal/ball/statevisual.go:102` `var prototypeVisuals bool`；打开它的产码调用者**只有** `cmd/balldebug/main.go:122`（U15）；`SetAudioLevel` 在它关的时候**直接 return**（`internal/ball/liquid_windows.go:51-54`）。⇒ 而 `statevisual.go:93` 注释逐字写明"翻默认值是 ticket 68 AC#2"。
2. 常驻腿能到达的状态（`Sleeping`／`Confirming`）**两枚都不在 `liquidDriven` 里** ⇒ 电平即使送达也会被 `applyLevel` 记进 `b.liqRaw` 后丢弃（`liquid_windows.go:56-63`）。

⇒ 所以本票**能**证到的是"真麦克风的电平进了球的门"（AC#2 那两枚不同的 `float32` 读数），**证不到**"球在呼吸"。把"接好了"的判据挂在 `Listening` 上也不行——那枚状态本票不产生（`EvSummon` 的边在 `table.go:51`，发射者同样不在本票地界）。这条判断直接进 ⑦-P2／P3，请裁。
