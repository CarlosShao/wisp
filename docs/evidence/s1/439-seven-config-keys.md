# 票 439 落地腿 `439-w1` — 7 枚零读者键逐枚复核与归因（2026-11-13）

- 性质：**写腿**（产码零字节——7 枚复核后**零枚判甲、零枚判乙**，本程的活＝AC#0 复核＋AC#1 归因＋台账回填；
  交付物＝本件＋票 439 Progress log＋commit，无代码改动）。
- 起手 HEAD：`673da11df`。基线：`git status --porcelain internal/config` ⇒ **空输出**（干净）。
  全树 `git status --porcelain` 只有 `.scratch/wisp/probes/**` 等他人/历史未跟踪台件与少数未跟踪
  `design/`、`frontend/`、`docs/evidence/` 散件——**已跟踪产码零处未提交改动**，无在飞写腿撞面。
- 写面声明：`internal/config/**` 本程**零触碰**（AC#3 契约轴零改动：`docs/PLAN.md`、`docs/specs/**`、
  `thresholds.go`、golden、`allowlist.txt` 全部原样）；`frontend/**`、`design/**` 零触碰；
  `build/mut_rvw.go` 未删。本程只新建本件＋追加票 439 的 Progress log。

## 1. AC#0 逐枚复核（本程现跑，⛔ 不照抄 180 的读数）

尺＝票面那把：`grep -rn "\.<字段名>" --include=*.go internal/ cmd/`；段内有歧义的字段换**段内立尺**
（如 `Voice\.TTS\.Speed`），另跑**裸词补尺**（`VetoWords`、`HideOnFullscreen`、`SLOSampleIntervalSec`、
`ThemeLight|ThemeDark`、`Language`）防"前缀不在同一表达式"的漏计。排除类按票面：struct 定义、
段拷贝（reload plan 的 apply/commit 段）、比较（tier 判分）、枚举/词汇校验、名册（roster）字符串、测试赋值。

| # | 字段（toml 键） | 粗尺命中（行） | 命中逐条判类 | 真读 | 与 180 表的差集 |
|---|---|---|---|---|---|
| 1 | `app.language` | **13** | manager.go:379 重启档比较；manager.go:387 注释；speech/cloudasr.go:515-516 与 tools/web_search.go:130-131 是**别家字段**（`ep.Language`/`q.Language`）；其余 9 行全测试（firstrun_198r2_test:253、watchdog/config_poll_test:299、boundary_test:79、manager_test:94、tiers_223r3_test:145/221/246） | **0** | **无实质差集**（180 同判零真读；粗尺绝对数差异＝180 排除了 `_test.go` 行，本程把测试行单列） |
| 2 | `app.theme` | **22**（`\.Theme\b`，另 `ThemeLight|ThemeDark` 补尺） | manager.go:378 比较＋:383 段拷贝；validate.go:51/53 枚举校验；schema.go:65-66/102 是 config 自己的词表常量；internal/ball 的 `ThemeDark/ThemeLight` 是 **ball 包自己的类型**（tokens.go:23-24，同名不同物）；其余 15 行全测试 | **0** | **无实质差集**（180 同判零真读） |
| 3 | `ball.hide_on_fullscreen` | **3**（裸尺 `HideOnFullscreen`） | schema.go:194 注释＋:195 字段定义；boundary_test.go:86 测试赋值 | **0** | **无实质差集**（180"全树仅定义＋赋值＋注释"逐条复现；341 文件的引用在 180 件里归类为注释面，本程粗尺未见额外行） |
| 4 | `voice.wake_word.veto_words` | **10**（裸尺 `VetoWords`） | schema.go:224-225 定义；manager.go:434 段拷贝＋:441 比较；speech/kws.go:74 注释＋:83 `kwsVetoWords` **空表定义**＋:291 `Veto: kwsVetoWords[...]`（消费的是 speech 包内 map，⛔ 不是 config 字段）；resident_voice_windows.go:1136 注释（"kwsVetoWords is empty"）；其余 2 行测试 | **0** | **无实质差集**（多出的 kws.go:291 是 speech 内部 map 的读者，不是 config 读者——恰是"接线半边已在、数据半边缺"的证据，见 §2#4） |
| 5 | `voice.tts.provider` | **2**（段内立尺 `TTS\.Provider`） | manager.go:408 reload 档比较＋:423 段拷贝 | **0** | **无实质差集** |
| 6 | `voice.tts.speed` | **4**（段内立尺 `TTS\.Speed`） | manager.go:435 段拷贝＋:442 比较；config_readers_255.go:333 名册字符串（引读点原文）；**internal/speech/cloudcatalog.go:119 唯一 config 读点** | **0**（读点是**死读点**） | **一处精化**：下游消费其实**在盘**——`cloudtts.go:140 [Speed: t.ep.ttsSpeed()]` 把 endpoint 的 Speed 写进云端请求体（`cloudtts.go:276-284` wire body、`cloudendpoint.go:148-152` 0→1.0 归一）。死读点死的**不是 Speed 没人消费**，而是 `ResolveCloudChain` 全仓非测试调用者 0（`config_readers_255.go:333` 名册同判、`config_receipt_255_test.go:365` 有牙钉住）⇒ 链不活、整段悬空。这把"接 61 即活"的判词钉得更死（见 §2#6） |
| 7 | `observe.slo_sample_interval_sec` | **3**（裸尺 `SLOSampleIntervalSec`） | schema.go:670-671 定义；boundary_test.go:185 测试赋值（unwired_test.go:237 是 TOML 夹具串，也是测试件） | **0** | **无实质差集**（180"全树仅 schema 定义＋boundary_test 赋值"复现） |

**AC#0 结论：7 枚全部复现"零生产读者"，与票 180 的读数无一处矛盾；一处精化（#6 的下游消费面）。**

## 2. AC#1 逐枚归因（甲乙丙丁）

**总判：0 甲 · 0 乙 · 4 丙（#1 #2 #3 #5 #6 中五枚里四枚主判丙）· 2 丁（#4 #7）＋ #2 的球侧半边附丁。**
逐枚理由（全部 file:line 本程现读）：

### #1 `app.language`（default zh-CN）→ **丙：归票 39（config editor GUI，open）**

- 三档机器全在盘：`tiers.go:60 [TierRestart]`、`manager.go:373-392 planApp`（restart 档比较＋OnRestartPending）、
  `cmd/wisp/config_reload.go:346-375 reportRestartPending`（逐字把 `app.language` 列进 restartTierKeys、
  句子里点名"界面语言"）＋ `restart_tier_keys_255r2_test.go` 名对名钉住（改句子＝动 255 的测试契约，本票不动）。
- 缺的是**重启时消费它去定 UI 语言的壳**。今天的 UI 文案全是硬编码中文（D23 的 i18n 推迟行
  `SPEC-12 §5:70 [DEFERRED i18n/非中文 → S8]`），语言选择点的活随**配置编辑器 GUI**（票 39，SPEC-08 §5.2/D36，
  Status ready-for-agent）落地；180 普查件 §2 末尾同样把 language/theme 指到票 39。
- 为什么不是甲：规格要求的是"这枚键按 restart 档对待"（已在盘），**没有一条规格要求本票此刻接出
  "UI 语言切换"的消费者**；消费面是 GUI 票的地界，且有具名 open 票认领 ⇒ 丙。

### #2 `app.theme`（default dark）→ **主判丙：归票 39（GUI/设置面）；球侧半边附丁（摆机主，见 §3）**

- hot 档机器在盘：`tiers.go:59`、`manager.go:378-385 planApp themeChanged`（真 hot 应用段）、
  `validate.go:51-53` 词表校验（dark|light|auto）。
- 消费缝**两条都在盘**：①GUI/设置面归票 39（同 #1 的理由）；②球侧：`ball.Options.Theme`（ball_windows.go:83）
  → `renderer SetTheme`（renderer_windows.go:107）→ `PaletteFor`（tokens.go:223）整链现成，**但 cmd/wisp 的
  ball.New 调用点（resident_ball_windows.go:429-471）没有喂 Theme**——grep `App\.Theme` 在 cmd/wisp 的非测试
  命中 0。
- 为什么主判丙不判甲：球的皮跟不跟 `app.theme` 走是一个**产品/视觉决定**，且 owner 2026-09-25 已把
  C21 权威转向第四代主题（tokens.go:73-75 注释逐字"c21's live design authority moved with the owner's
  2026-09-25 pivot to the fourth-generation theme"）——二值词表（dark/light/auto）与第四代主题的关系
  **规格没写清**；且球侧接线没有任何 open 票认领（62 review、68 review 都不是 theme 认领票）。
  ⇒ GUI 半边丙（39），球侧半边丁（§3 摆机主）。本票两半都不动。

### #3 `ball.hide_on_fullscreen`（default true）→ **丙：归票 62/68（球视觉族，均 open）＋ 功能半边缺的事实**

- hot 档（`SPEC-03:27`、`PLAN.md:2741`）。**功能半边不在**（本程现量）：`internal/ball` 无任何
  "前台窗口是否全屏"的判定（grep `[Ff]ullscreen` cmd/ 0 命中；internal/ball 的 8 处 monitor 命中全是
  monitorFromWindow/工作区几何），`ball.Options` 也没有 HideOnFullscreen 字段——想接都没有对象，
  接之前要先造"全屏判定"这个新功能件。
- 归口：票 83 已把球体视觉键族归口"票 62/65/68＋39"；票 180 复核 62（review）与 68（review，
  blocked-on-owner R15）均 open。⇒ 丙，具名 62/68；谁先开工谁把"全屏判定＋隐藏切换"一起接。

### #4 `voice.wake_word.veto_words`（default 取消,停下,别）→ **丁：挂账，激活前置＝模型资产（摆机主）**

- **消费半边已在盘且已接线**（这是 180/423 表没写足的）：`speech/kws.go:291 [Veto: kwsVetoWords[res.Keyword]]`
  分类 → `cmd/wisp/resident_voice_windows.go:1155-1159 onKwsHit`（`hit.Veto` ⇒ `TriggerInterrupt` ＋ 状态行走，
  381-r1 预留的通道）。分支是真的，不是死码。
- **缺的是数据，且 config 喂值救不了**：`kwsVetoWords` 是空表（kws.go:83）是**刻意的**——kws.go:78-82 逐字：
  shipped keywords 表（manifest sha256 46de8c68）只有 8 个唤醒词、无 veto 词，"the B leg adds veto words by
  **editing the keywords file**"。KWS 是转导表数据驱动：config 里的三个词不进转导表就**永远不会被 spot 出来**；
  把 config 值喂进 `kwsVetoWords` 只能给模型已认出的词贴标签＝接线字面上的假绿
  （`.scratch/wisp/probes/orch/20261110/deferred-disposition.md` §一#2 同判，本程独立复核成立）。
- 为什么是丁不是甲：激活它＝**重生成 keywords 转导表（模型资产＋manifest sha256 变更）**——花钱/资产决策，
  摆机主（§3）。挂账五字段草稿已在 disposition §三.1；SPEC-12 §5 的 B1 行（:76，快捷键路径语音否决）
  是**另一码事**（disposition 明说"别混"），不因本枚改写。

### #5 `voice.tts.provider`（default local-sherpa）→ **丙：归票 61 收尾（SPEC-12 §5:81 行 open）**

- 消费件**全部在盘**：`speech/cloudchain.go:144 NewCloudTTSChain`（:97 实装 `TtsEngine`）、
  `cloudtts.go:117 CloudTTS.Speak`、`cloudcatalog.go:40 ResolveCloudChain`（非测试调用者 0）、
  KeyResolver 缝＝`secret.Store.Resolve`（providers.go:125 等多处生产现用）。
- 装配现状（本程现量）：`resident_voice_engines_windows.go:436-448` TTS loader **无条件** `NewMatchaTts`；
  boot 门 `resident_voice_windows.go:1866` 只拦 **ASR** provider，TTS 半边没有门；`voicePlan` 无 TtsProvider 字段。
- 为什么是丙：`tts.go:24-25` 逐字——R5 裁定默认语音路径＝云 TTS（StepFun，票 61）；`SPEC-12 §5:81`
  "云端 ASR/TTS provider（级联，C9）……现随 **票 61** 落地"仍 open ⇒ 具名认领＝**票 61 收尾**
  （disposition §一#3 判"可立刻接"）。本票不接：那是一张完整收尾票（provider 分派＋boot 门 TTS 半边＋
  KeyResolver 接线＋成对判据族），且 `cmd/wisp` 是共享编译面（每波只许一枚写腿），不是顺手一行。

### #6 `voice.tts.speed`（default 1.0）→ **丙：随 #5 搭车（票 61 收尾）；本地半边钉死不接是终态**

- 唯一 config 读点 `cloudcatalog.go:119` 死读（`ResolveCloudChain` 非测试调用者 0）；
  **但下游消费在盘**：`cloudtts.go:140` 把 `ep.ttsSpeed()` 写进云端请求体（§1#6）⇒ 61 收尾把
  `ResolveCloudChain` 接上组合根的那一刻，speed 自动词活，零额外工作量（disposition §一#4"搭车 1"同判）。
- 本地半边：`tts.go:34-48` `ttsSpeed = 1.0` 是**测量纪律常数**（15-b/15-c；注释逐字"the honest way to
  make speech longer is to speak more text, not to resample the budget"）——旋钮化会破坏"同证据两次运行
  可比"。本地不旋钮化是**终态**，五字段挂账草稿在 disposition §三.2。

### #7 `observe.slo_sample_interval_sec`（default 5）→ **丁：字段语义未定案（摆机主；编排者已裁挂账 DEFERRED）**

- 全树零读者（§1#7）。任何接法都动冻结面：`SPEC-03 §4.3` 冻结 watchdog 1s tick 双职责——
  生产实现就是 `cmd/wisp/config_reload.go:101 [configReloadPollInterval = time.Second]` 那条 tick；
  `internal/watchdog` 的 `SustainTicks=60` 按 1s tick 定义 60s，改节拍＝重述 D32 sustain 语义。
- 契约项 C1 已由编排者裁定"**不拆**"（台账 `pending-and-issues.md:34012` 逐字：
  "slo_sample_interval_sec＝**挂账 DEFERRED**（不改契约文字，标准五字段登记）⇒ O 不产生新的机主契约项"）。
- 为什么仍是丁：**字段本身的语义**（"SLO 采样节奏"到底＝同循环每 N tick 采样，还是别的定义）没有任何
  契约文字写清——接它必须先改 `SPEC-03 §4.3`/D32 文字（改契约＝人工批准）。⇒ 摆机主三选一（§3），
  维持挂账是现状合法态。

## 3. AC#2 / AC#3

- **AC#2（成对判据）：零甲 ⇒ 无判据可跑，如实记**。7 枚没有任何一枚判"该接读者"，所以本程没有
  "设两个不同值 ⇒ 结果随之变"的落地对象；反向判据（不许改 default 交差）自动满足——
  `git diff HEAD -- internal/config/schema.go` ⇒ **空**（`default:` 标签 70 枚一字未动）。
- **AC#3（契约轴）**：`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`
  **零字节改动**（本程 diff 只有本件与票面 log）。需要**契约文字/机主裁定**的三处，具名摆机主：

| # | 摆机主的裁定 | 选项 | 不答的代价 |
|---|---|---|---|
| 1 | **`app.theme` 的皮跟谁走**（球侧半边）：D36 说 theme=hot、D29 说球有 system-theme follow（tokens.go:17），但二值词表（dark/light/auto）与第四代 C21 主题的关系没写清，球侧接线无认领票 | ①批"球吃 app.theme"→开小票接线（缝在盘：Options.Theme＋SetTheme）；②批"theme 只归 GUI/面板"→归票 39 并在 D36 备注消费对象（改契约文字须批）；③维持现状挂票 39 | 键继续零读者；GUI 票 39 的 hot 徽章对 theme 撒谎 |
| 2 | **`voice.wake_word.veto_words` 要不要激活**：接线被模型资产卡死（keywords 转导表无 veto 词，重生成＝资产＋manifest sha256 变更＋`internal/speech` 词集注入点） | ①批资产重生→开 B 腿（判据/五字段已在 disposition §三.1）；②不批→维持挂账（现状） | 功能分支永远打不响（onKwsHit 的 veto 半边在盘等数据）；config 键继续零读者 |
| 3 | **`observe.slo_sample_interval_sec` 的语义**：①给语义（同循环 N-tick 采样？）→须改 SPEC-03 §4.3/D32 文字＝改契约须批；②批删键→也是改契约（D36/SPEC-03 表里有它）须批；③维持挂账（编排者已裁 DEFERRED，现状） | 同左三选一 | 键继续零读者；写任何值都白改 |

## 4. AC#4 门禁五条读数（commit 前工作树现跑；本程 Go 零改动 ⇒ 读数＝基线）

1. `go build ./cmd/wisp/` ⇒ **rc=0，无输出**。
2. `go vet ./cmd/wisp/` ⇒ **rc=0，无输出**。
3. `go test -count=1 ./internal/config/ ./internal/panel/` ⇒
   `ok github.com/CarlosShao/wisp/internal/config …s` ＋ `ok github.com/CarlosShao/wisp/internal/panel …s`，**rc=0 全绿**。
4. `gofumpt -l internal cmd tools scripts frontend` ⇒ **空输出**。
5. `sh scripts/d22scan.sh` ⇒ 末行 `d22scan: clean - no D22 ban violations; …`，**rc=0 干净**。
   ⚠ `./cmd/wisp/` 包测试按票 98 本机测不到，不进本腿红账。

（逐条原始输出见本程会话记录；五条命令逐字见票面 Progress log。）

## 5. AC#5 边界声明

**本票收：0 枚接线、0 枚删字段**（7 枚零甲零乙；收的是复核＋归因＋台账回填）。
**不收、留回台账的 7 枚去向：**

| 字段 | 去向 |
|---|---|
| `app.language` | 票 39（open，GUI 配置编辑器） |
| `app.theme` | 票 39（GUI 半边）＋ 球侧半边摆机主（§3#1） |
| `ball.hide_on_fullscreen` | 票 62/68（open；接线前置＝"全屏判定"功能半边，谁开工谁接） |
| `voice.wake_word.veto_words` | 摆机主（§3#2）；挂账五字段草稿 disposition §三.1 |
| `voice.tts.provider` | 票 61 收尾（SPEC-12 §5:81 行 open） |
| `voice.tts.speed` | 随票 61 收尾搭车；本地半边终态不接（disposition §三.2） |
| `observe.slo_sample_interval_sec` | 摆机主（§3#3）；编排者已裁挂账 DEFERRED（台账 34012） |

台账回填：`docs/reports/pending-and-issues.md` 的追加由**编排者**做（本腿只交证据件与票面 log，
不直接写台账——与票 180 腿同界）。

## 6. 本程没做（明写）

- **没改任何 Go 产码**（7 枚零甲零乙 ⇒ 无可落）；`internal/config/**` 零触碰。
- **没跑任何成对判据**（零甲 ⇒ 无对象；不是漏跑）。
- **没动 SPEC-12 §5 登记表**（veto/slo 两枚的"挂账→登记行"是台账写面，归编排者）。
- **没跑真窗/winlive**（无判据可跑，不适用）。
- #6 的精化（云端请求体消费 Speed 在盘）只改判词的**精度**，不改判类。

## 7. 一句结论

7 枚全部复现零生产读者、且**没有一枚是"接一行就好"**：4 枚有具名 open 票认领（39／62·68／61 收尾），
2 枚卡在只有机主能答的地方（veto 的模型资产、slo 的契约语义），theme 还欠"皮跟谁走"一句话。
本票以零产码收口，把每枚的去向钉到了 file:line。
