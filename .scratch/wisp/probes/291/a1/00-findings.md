# 票 291 · 只读普查腿 `291-a1` — AC#1（契约面判定）＋ AC#3（最小改动落点表）

生成时刻：`2026-10-09 14:0x +08`。两次 `date` 原样输出（⛔ 本文所有时刻以 stdout 为准，不以票面推断）：

```
$ date
Fri Oct  9 14:05:48 CST 2026        # 开工第一发
Fri Oct  9 14:08:54 CST 2026        # AC#1 字样名册现量那一发
```

分支 `dev`；开工时 `git status --short` 显示工作树**已有别人未入库的改动**（`.gitignore`、`.scratch/wisp/probes/{152,161,242,268}/**`、`design/**` 下一批 `D`/`M`）——本腿**一字节未碰、未还原、未提交**。本件 commit 之后 HEAD＝`67b4b57d`（骨架枚）。

本腿性质＝**只读**。⛔零产码改动、⛔零 `docs/**` 改动、⛔零 `go build`／`go vet`／`go test`／`./...`（Go 编译/测试面此刻由写腿 `289-r1` 独占，它的"改前基线"要求树是干净的 HEAD）。
⇒ 因此 §2.2 那张"会被打红的钉"名册**全部来自读断言原文＋算术**，凡读码定不了的，都按红线具名写在 §2.3，⛔ 没写成结论。

---

## §0 射程与命令面（可重跑）

写点唯一＝`.scratch/wisp/probes/291/a1/00-findings.md`（本件；目录 `291/a1/` 由本腿新建）。件名 `.md`（根 `.gitignore` 全仓忽略 `*.out`，故不用 `.out`）。
文件是否存在一律只认 git 对象层：`git ls-tree -r --name-only HEAD <dir>`／`git cat-file -e HEAD:<path>`。`frontend/**` 与 `design/**` 零读零写（`design/**` 只出现在 §0 引用的 `git status` 摘要行里，⛔ 未打开其中任何文件）。

跑过的命令全名册（编号供下文 §引用）：

| 号 | 命令（在仓根跑） | 用途／产出 |
|---|---|---|
| C1 | `git status --short` / `git rev-parse --abbrev-ref HEAD` / `date` | 开工卫生与钟 |
| C2 | `Read .scratch/wisp/issues/291-…dc…md` | 票面 25 行逐字读 |
| C3 | `git ls-tree -r --name-only HEAD docs/` | docs 面盘点 |
| C4 | `grep -n -E "电平\|音量\|响度\|[Ll]evel\|RMS\|rms\|volume\|loudness\|0\.\.1\|满刻度\|full.?scale\|dBFS\|呼吸\|breath\|包络\|envelope" docs/PLAN.md` | AC#1 主搜（命中 40 行） |
| C5 | 同一枚 18 项 alternation，对 `docs/specs/SPEC-04-voice-pipeline.md`／`SPEC-05-agent-kernel.md`／`SPEC-08-ui-ball-panel.md` 逐枚 `git cat-file -e`＋grep | 三枚 spec 面；**SPEC-05-agent-kernel.md 在 HEAD 里不存在**（见 §4 冲突①） |
| C6 | `git ls-tree -r --name-only HEAD docs/specs/` ＋对 `SPEC-04` 加宽字样（`振幅\|幅度\|波动`）＋`wc -l` | 真名＝`SPEC-05-agent-core.md`；SPEC-04 **142 行、18 项字样零命中**（加宽后仍零命中） |
| C7 | `grep -n -E "VAD\|energy\|能量\|speech\|ASR\|frame\|帧\|采集\|麦克风" docs/specs/SPEC-04-voice-pipeline.md` | **证明 C6 的零命中不是"文件是空的"**：SPEC-04 有 25 行在讲采集/VAD/ASR，⛔ 没有一行讲电平 |
| C8 | `grep -n -E "…level terms…" docs/specs/SPEC-05-agent-core.md` | 零命中 |
| C9 | `Read internal/audio/level.go`（全 148 行） | AC#3 位点 |
| C10 | `Read internal/audio/level_test.go`（全 399 行） | AC#3 钉名册 |
| C11 | `grep -n -E "D-1\|envelope\|电平\|音量\|level\|Level" docs/specs/SPEC-12-roadmap-governance.md` | **零命中**（标度问题不在 SPEC-12 §5 登记表里） |
| C12 | `grep -n -E "DEFERRED\|D-1" docs/PLAN.md` | PLAN 的 DEFERRED 表 12 行，⛔ 无一枚讲电平/包络 |
| C13 | `sed -n '1,60p' internal/audio/capturelevel_windows_test.go` ＋ `grep -n -E "12345\|0\.37\|DC\|dc\|expect"` | DC 夹具：全文件**只有 `:23` 一枚 `12345`，零枚字面期望值** |
| C14 | `grep -n -E "2\|Skip\|WISP_LIVE_MIC\|silent\|loud\|×\|倍" cmd/wisp/resident_audio_247_live_windows_test.go` | 252 行；定位 ×2 门槛与三处闸 |
| C15 | `sed -n '126,205p' cmd/wisp/resident_audio_247_live_windows_test.go` | ×2 断言逐字 |
| C16 | `grep -n -E "level\|Level\|float32\|sink\|Sink" cmd/wisp/resident_audio_windows.go` ＋ `sed -n '24,32p;105,135p;205,215p'` | 投递那一跳的形状 |
| C17 | `grep -rn -E "SilenceLevelGate\|SetAudioLevel\|AudioLevel" internal/ball/*.go` ＋ `sed -n '25,35p' liquid.go` ＋ `sed -n '28,42p' liquid_windows.go` | 消费者侧门限与入口 |
| C18 | `grep -n -B2 -A6 -E "输入\|音量\|波动" docs/specs/SPEC-08-ui-ball-panel.md` | ⓒ 许诺原文 |
| C19 | `grep -n -B3 -A10 -E "D-1" .scratch/wisp/probes/240/c1/census.md` | D-1 登记原文（`:479`） |
| C20 | `grep -n -E "直流\|电平\|RMS\|标度" docs/reports/pending-and-issues.md` | 台账现量：`直流` 只在票 291 的挂号条目里出现，⛔ 无既往裁定 |
| C21 | `git ls-tree … docs/evidence/s1/ \| grep -E "/(241\|247\|290\|291)-"` ＋ `grep -rn -E "直流\|DC offset\|契约面" docs/evidence/s1/241-*` | **241-v1 §1.7 已就 DC 出过判语**（`:108-111`，标题在 `:106`） |
| C22 | `grep -c -n "Progress log" <票 291>` ＋ `grep -rln "## Progress log" .scratch/wisp/issues/` | 票 291 **零枚** Progress log 标题（见 §4 冲突③） |
| C23 | `sed -n '100,125p' docs/evidence/s1/241-audio-level-producer-v1.md` | 判语逐字 |
| C24 | `grep -n -E "3\.9\|dBFS\|0\.368462\|0\.376740\|12345" .scratch/wisp/probes/247/v1/30-ac2-level-ruler.md` | 前腿读数（⛔ 不是本腿读数） |
| C25 | 逐字样 `grep -c -F` 在 4 枚面上出**字样名册表**（§1.1） | 覆盖率证明 |
| C26 | `grep -n "^#\{1,4\} " docs/PLAN.md`（取 3560 之前 5 枚）＋`grep -n -E "^#{2,4} .*(D5\|D16\|D18\|D32\|D47)"` | 确认 `:3554` 落在 §17.7 悬浮球原型；D 表节界 |
| C27 | `grep -n -E "\| C8 \|\|…C12 …C17 …" docs/PLAN.md` | C 表音频/球条目原文（`:1358`/`:1362`/`:2434`/`:2968`） |
| C28 | `sed -n '60,100p' internal/audio/capturelevel_windows_test.go` ＋ 全仓 `grep -rn -E "0\.3767\|0\.3684\|32768" --include=*.go` ＋ `grep -rn -E "level (==\|!=) 0\|== 0\.0\|not zero" --include=*_test.go internal/audio cmd/wisp` | 全仓**零枚字面电平断言**；零枚"电平必须非零"断言 |
| C29 | `sed -n '1,20p;540,560p' internal/ball/tokens_table_test.go` | `SilenceLevelGate` 是 **C21 token 行的机器核对钉** |
| C30 | `grep -n -E "C21" docs/PLAN.md` ＋ `grep -n -E "SilenceLevelGate\|LevelAttackTauMs\|电平" docs/evidence/s1/c21-native-tokens.md` | C21＝冻结契约（`PLAN.md:1371`）；token 表 `:162-164` 三行把"每单位电平"的视觉增益钉成契约值 |

---

## §1 AC#1 — 这枚定义是不是契约面

### 1.1 字样名册与搜索面（⛔ 不许只答"没找到"：先证明搜过）

命令＝C25（`grep -c -F "<字样>"`，枚枚单独跑，**计数单位是"含该字样的行数"**，不是出现次数）。面＝四枚：`docs/PLAN.md`、`docs/specs/SPEC-04-voice-pipeline.md`、`docs/specs/SPEC-05-agent-core.md`（票面写的 `SPEC-05-agent-kernel.md` 不存在，见 §4①）、`docs/specs/SPEC-08-ui-ball-panel.md`。

| 字样 | PLAN.md | SPEC-04 | SPEC-05 | SPEC-08 |
|---|---|---|---|---|
| `电平` | **0** | **0** | **0** | **0** |
| `音量` | 8 | 0 | 0 | 1 |
| `响度` | 0 | 0 | 0 | 0 |
| `level` | 3 | 0 | 0 | 0 |
| `Level` | 18 | 0 | 0 | 0 |
| `RMS` | 0 | 0 | 0 | 0 |
| `rms` | 0 | 0 | 0 | 0 |
| `volume` | 1 | 0 | 0 | 0 |
| `loudness` | 0 | 0 | 0 | 0 |
| `0..1` | 0 | 0 | 0 | 0 |
| `满刻度` | 0 | 0 | 0 | 0 |
| `full scale` / `full-scale` | 0 / 0 | 0 / 0 | 0 / 0 | 0 / 0 |
| `dBFS` | 0 | 0 | 0 | 0 |
| `呼吸` | 13 | 0 | 0 | 4 |
| `breath` | 0 | 0 | 0 | 0 |
| `包络` | 0 | 0 | 0 | 0 |
| `envelope` | 0 | 0 | 0 | 0 |
| （加宽）`振幅`/`幅度`/`波动` | — | 0/0/**0** | — | — |

零命中不是"文件空"：`wc -l` 现量 `SPEC-04` **142 行**、`SPEC-08` 有 240＋行；C7 证明 SPEC-04 确实在讲这条链——`:41` 逐字 `采集：WASAPI 共享模式；设备原生采样率 ≠16k 时**进程内重采样**`、`:42` `帧长对齐 VAD 输入（512 样本 @16k ≈ 32ms）`、`:53` VAD 行。**它规定了帧形状（512/32ms），没有规定电平**。

命中里**七枚 `音量` 是假阳性**（是"系统音量"这个工具面，不是麦克风电平）：`PLAN.md:358`（无 `system.*`（音量/亮度/媒体键…））、`:378`（`sysinfo` 电量/网络/音量）、`:380`（`media` 播放/暂停/音量 L1）、`:1667`、`:2554`（`system.get` …音量…）、`:2555`（`system.set` 音量/亮度/静音）、`:3125`（场景脚本"音量调到 30%"）。
`Level` 的 18 行也是假阳性：`RiskLevel`（`:1351`/`:1352`/`:1364`/`:1374`/`:2523`/`:2531`/`:2980`…）、意图分类 Level 1/2/3（`:278`/`:280`/`:295`/`:1774`）、`:198` macOS `NSWindow` + `level = .floating`、`:2747` `[observe] level`（**日志级别**）。
⇒ **四枚契约面里，"麦克风电平"这件事在文本上的唯一落点＝ 1 行**（下 ⓑ/ⓒ）。

### 1.2 ⓐ 有没有明文把电平定义成"含直流的全波 RMS"？

**契约面（`docs/PLAN.md` ＋三枚 spec）：没有。一处都没有。**
判据＝上面 §1.1 的零命中名册：`电平`／`RMS`／`rms`／`dBFS`／`满刻度`／`full scale`／`0..1` 在四枚面上**全部零命中**；C 表里与音频/球相关的条目逐枚读过——
`docs/PLAN.md:1358`：`| C8 | \`AudioSource\` | 真实麦克风 / **wav 注入（测试钩子）** / \`[Aec]\` 仅接口位 | D16, 测试地基 |`
`docs/PLAN.md:1362`：`| C12 | \`BallState\` | **20 态** + **完整转移表（D43，第四轮才真正写出）** + 每态超时 | §2 + D43 |`
C17（`:2434`/`:2968`）＝面板方法白名单，与电平无关。
`D5`（`PLAN.md:139`）／`D16`（`:444`）／`D18`（`:497`）／`D32`（`:2221`）／`D47`（`:3241`）五节的**节界已由 C26 钉出**，里面没有一行的主题是电平标度。

**但"含直流的 RMS"这句话确实存在——它在代码与裁决表里，共三处：**

1. `internal/audio/level.go:17-20`（逐字）：
   > `// THE SCALE. level = sqrt(mean(sample^2)) / LevelFullScale over one C8 seam`
   > `// frame, i.e. plain unweighted RMS of the int16 samples referenced to the`
   > `// largest magnitude int16 can hold. It is a physical number, not a perceptual`
   > `// one: no gain, no compression, no noise floor, no smoothing.`
2. `internal/audio/level.go:37-39`（逐字，**最硬的一枚**）：
   > `//	  DC counts: RMS does not remove a DC offset, so a frame held at -32768`
   > `//	  reads 1.0. Consumers that need "voice only" gate it themselves (the ball`
   > `//	  does, with SilenceLevelGate).`
3. `docs/evidence/s1/241-audio-level-producer-v1.md:108-111`（非实现者验收腿**已经裁过 DC**，逐字；`### 1.7` 标题在 `:106`）：
   > `level.go:37-39 声称 "DC counts: RMS does not remove a DC offset, so a frame held at -32768 reads 1.0"。`
   > `本腿实测：整帧恒 \`-32768\` ⇒ \`1\`；整帧恒 \`20000\`（**纯直流、零交流**）⇒ \`0.6103515625\` ＝ \`20000/32768\` 精确。`
   > `⇒ 声称与实读一致。物理后果也一并登记在这里：**采集中链上任何直流偏置都会读成"响"**，`
   > `本票选择"不除 DC、并在注释里点名由消费者自己开门"，这与票面"线性、无增益、无压缩、无噪声地板"是同一句话的另一面，不是缺陷。`

⇒ **判语（严格版，别读松）**：去直流**不是** `D1–D47`／`C1–C32` 的契约变更（AGENTS.md §1.1 那串名册里没有一行的射程覆盖"电平标度"），**但它是"已被非实现者验收腿裁过并写进裁决表"的判断**——`241-v1:111` 那句"不是缺陷"与本票标题"这是一枚真缺陷"**正面冲突**。
另外两处既往记录构成同一方向的重量：
- `.scratch/wisp/probes/240/c1/census.md:479` 把标度问题登记成 **D-1**，逐字：`| **D-1** | **RMS 到 0..1 的标度到底是什么？** 线性 \`rms/32767\`？还是要加增益/压缩？ | … | 这是**观感口径**，不是代码事实。…定标度要有真声音＋真桌面，本腿零跑进程、零读视觉层 |` ——⚠ 注意 D-1 问的是**增益/压缩**，**没问过直流**；这一维从来没人裁过，是**盲区**而非"已定案"。该 D-1 ⛔ **不在** `SPEC-12 §5` 登记表里（C11 零命中），也不在 `PLAN.md` 的 DEFERRED 表里（C12）⇒ 它是一枚**普查件内的待定案项**，不是登记过的推迟项。
- `docs/reports/pending-and-issues.md:9790`（编排者裁语，逐字）：`**默认动作＝维持线性、不加增益、不偷偷压曲线**（任何偷偷的曲线都会被打红：\`TestLevelIsLinearInAmplitude\` ＋端到端那枚严格 0.5）` ＋ `**撤销口令**＝「241 增益改要摆」`。

**给不懂技术的机主的那句后果（⛔ 不是"要不要人工批准"的抽象问题）**：尺减不减直流，决定**机主能不能感觉到球随声音动**。今天不减直流＋本机麦带 DC ⇒ 静默时球收到的是 `0.368462`，而球侧那道判"没人在说话"的门限是 `SilenceLevelGate = 0.06`（`internal/ball/liquid.go:30`，逐字 `SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"`）——**0.368 远大于 0.06，这道门今天永远开着**，球的"有没有人说话"在真机上是常亮而不是会关。这一句是"改动能不能被感觉到"的实测面，也是必须摆给机主的原因。

### 1.3 ⓑ 有没有一处只定义成"随响度单调上升的 0..1 标量"？

**有，而且契约面总共就这一种形状。**四枚面的全部相关原文：

- `docs/PLAN.md:3554`（§17.7 悬浮球原型，节界由 C26 确认）：
  > `| \`Listening\` | 直径 56px + **外环随音量波动**（原型用 3 层同心环错峰呼吸示意）+ \`--accent\` 描边 |`
- `docs/specs/SPEC-08-ui-ball-panel.md:63`（同一件事的 spec 版）：
  > `| \`Listening\` | 56px + 外环随音量波动 + accent 描边 |`

两枚都**只规定"随音量波动"**＝单调上升的形状，⛔ 没有任何一处规定这个标量的**直流那一维**（既没说含直流，也没说不含）。
`docs/reports/pending-and-issues.md:970` 还留着一枚未闭的措辞账（逐字）：`| **Q-12** \`Listening\` 写"外环随音量"、实为时基 | **改措辞**（环=时基；电平驱动光晕与转速） | 视觉上确有"随音量"的观感，不值得为文案改渲染 |` ⇒ 契约文本连"哪一层像素随电平动"都没定死，只定了"电平驱动光晕与转速"这一单调关系。

唯一把数值形状写出来的地方是**代码文档，不是契约**：`internal/ball/liquid_windows.go:35-36` 逐字
> `// SetAudioLevel feeds one capture-frame envelope (RMS of the last 512 samples,`
> `// normalised 0..1) to the liquid.`

（`240-c1` 的 `census.md:355` 正是拿这一行当"除数全仓无定义"的证据；`census.md:479` 的 D-1 也源于此。）

**契约面里唯一带"每单位电平"数值的，是 C21 token 表的三行**（C30，`docs/evidence/s1/c21-native-tokens.md:162-164`，逐字节选）：
> `:162` `| 音频包络 → 液体**几何**（票 62 建，票 74 换指向） | 每单位电平可见半径收缩 0.12 / 每单位唤起爆发扩散 0.10 | \`LiquidGatherPerLevel\` / \`SummonFlowSpread\` | …`
> `:163` `| 音频包络 → 液体**转速** rad/s（\`liquid.go\`，票 62，无 CSS 对应） | 基准 0.45 + 每单位电平 4.20 + 唤起爆发期 2.20 | …`
> `:164` `| 包络门限与时间常数（\`liquid.go\`，票 62，无 CSS 对应） | 静默门 0.06 / 静音保持 150ms / attack 60ms / release 420ms | \`SilenceLevelGate\` / …`

C21 是**冻结契约行**（`docs/PLAN.md:1371`：`| **C21** | **\`DesignTokens\`** | 单一 token 源（色板/圆角/阴影/字体与字号阶梯/缓动曲线/间距）…`）。
⇒ 这三行说的是"**每单位电平**产生多少像素/转速、超过多少算静默"＝**对单位电平的导数与门限**，⛔ 不是"单位电平本身是什么"。去直流**不需要动这些数字本身**，但会**换掉它们的物理含义**（0.06 从此是"去直流后的 0.06"）。这一格是 AC#1 判语里唯一真正靠向契约的边，具名留给编排者。

### 1.4 ⓒ SPEC-08 球那侧对"输入"的许诺原文

C18＋C6：`SPEC-08` 全篇 `电平` **零命中**、`RMS` **零命中**、`0..1` **零命中**、`float32`/`512` **零命中**（`grep -n -E "回调|512|float32|输入面|接口|C17|C12"` 只命中 `:3` 追溯行、`:85` 转移表标题、`:156`/`:235` 的 C17 面板白名单、`:25` 的 C21 token 要求）。
它对"输入"的全部许诺就是 `:63` 那一行 `56px + 外环随音量波动 + accent 描边`，加上动画纪律（`:19-22`，逐字）：
> `- **动画纪律**：\`Sleeping\` 态完全静态（不启动定时器，CPU≈0）；仅`
> `  \`Listening/Thinking/Acting/Speaking/Confirming\` 启动动画定时器，帧率上限 **30fps**，`
> `  CPU 计入工作态预算。…`

⇒ **SPEC-08 许诺的是"可验证的形状"（随音量波动），不是"可验证的数值"**：它⛔没许诺任何标度定义、⛔没许诺静默时读 0。
所以严格讲：**今天真机那条 `0.368462` 并不违反 SPEC-08 的字面**（尺仍然单调，球仍然"随音量波动"）；它违反的是**产品意图**（球分不开"安静"与"有人在说话"）与票 247 的 `×2` 判据。⇒ 本票的凭据**不能引 SPEC-08 说"契约被破坏"**，只能引：① `liquid.go:30` 的门限语义在真机上失效（0.368 > 0.06），② 票 247 `:194` 的门槛不可达。**这一句必须写清，否则对抗验收会抓我们把"观感问题"说成"契约违反"。**

### 1.5 AC#1 判语（三问各一句）

- ⓐ **没有**任何契约文档把电平定义成"含直流的全波 RMS"——该定义只存在于 `internal/audio/level.go:37-39` 的代码注释、`:17-20` 的标度段，以及裁决表 `docs/evidence/s1/241-audio-level-producer-v1.md:111` 那句"不是缺陷"。
- ⓑ **只有这一种**：契约面对电平的全部许诺＝"随音量波动"（`PLAN.md:3554`／`SPEC-08:63`）＝单调 0..1 标量，⛔ 未规定直流那一维 ⇒ 按票面 `AC#1` 的判据，**去直流属实现细节、本票可自办**。
- ⓒ `SPEC-08` 球侧原文只有 `:63` 那行 `56px + 外环随音量波动 + accent 描边`；它许诺形状不许诺数值，⇒ **不能拿它当"契约被破坏"的凭据**。

**但"可自办"要附三个具名条件**（⛔ 不是本腿自己加闸，是盘上已有的三笔账）：
1. 它会**推翻一枚已记录的非实现者判语**（`241-v1:111`"不是缺陷"）⇒ 按 AGENTS.md §0 第 3 条／`SPEC-12 §4.3` #1，**改动的缺口审计与验收必须由另一枚 agent 做**，且这一格应由编排者把 `241-v1` 的判语与票 291 的标题并排摆出来。
2. 它会**换掉 C21 token 表 `:162-164` 三行的物理含义**（数字不动、单位动）⇒ 若日后要重定 `SilenceLevelGate = 0.06` 的值，**那才是碰冻结面**（`internal/ball/tokens_table_test.go:552` 机器核对钉）。本票范围内⛔不许动那一枚常数。
3. 台账 `:9790` 那句"维持线性、不加增益、不偷偷压曲线"与撤销口令「241 增益改要摆」**只覆盖了增益/压缩那一维，没覆盖直流**；`240-c1` D-1（`census.md:479`）问过标度但⛔没问过 DC ⇒ **直流是一枚无人裁过的盲区**，本票不是"翻案"而是"补裁"。这一句差别要写给机主，否则他会以为自己要推翻自己已经拍过的板。

---

## §2 AC#3 — 最小改动的落点表（⛔ 本格零落地，一字节产码未改）

### 2.0 先摆最小改法本身

今天的求和段（`internal/audio/level.go:96-107` 逐字）：
```go
func LevelOfSamples(samples []int16) float64 {
	if len(samples) == 0 {
		return MinLevel
	}
	var sumSquares int64
	for _, s := range samples {
		v := int64(s)
		sumSquares += v * v
	}
	rms := math.Sqrt(float64(sumSquares) / float64(len(samples)))
	return rms / LevelFullScale
}
```
去直流＝两趟：先 `sum` 求 `mean`，再累加 `(v-mean)²`。**形状代价**（读码可得，⛔ 不必真跑）：
- 两趟 ⇒ 循环次数从 512 变 1024（仍在一枚帧上，D38b 名册不变、不新增协程、不用时钟）。
- `int64` 精确性论证**会失效**：`:89-90` 逐字 `The sum of squares accumulates in int64, which is exact for any seam frame (512 * 32768^2 = 2^39, far inside int64)` 与 `:91` `so the same input always produces the same bits.` ⇒ 去直流后必须先在 `int64` 里精确算 `mean`（`sum` 最大 512×32768＝2^21，int64 内精确），偏差要摊成整数（例如 `sum/n` 的定点化）才能保住"same input same bits"，否则 `TestLevelSameInputSameBits`（`level_test.go:215-234`，逐枚比 `math.Float32bits`）的**论证**先变假。⇒ 这是改法要还的第一笔债。

### 2.1 位点逐字引文（五枚文件）

| # | 文件 | 位点 | 逐字原文（关键行） |
|---|---|---|---|
| L1 | `internal/audio/level.go` | `:96-107` | 上面那段（`var sumSquares int64` … `return rms / LevelFullScale`） |
| L2 | 同上 | `:50-54` | `// LevelFullScale is the denominator of the level scale: the largest magnitude` / `const LevelFullScale = 32768.0` |
| L3 | 同上 | `:66-76` | `// SineLevelTolerance is the named tolerance … Measured worst case over the frames in` / `level_test.go: 2.31e-5. 5e-5 is that bound with headroom …` / `const SineLevelTolerance = 5e-5` |
| L4 | 同上 | `:78-82` | `// MinLevel and MaxLevel are the closed ends of the scale.` / `MinLevel = 0.0` / `MaxLevel = 1.0` |
| L5 | 同上 | `:118-128` `FrameLevel` | `func FrameLevel(frame []byte) (float32, error) {` … `return float32(LevelOfSamples(samples)), nil` |
| L6 | 同上 | `:137-147` `DecodeFrame` | `samples[i] = int16(uint16(frame[2*i]) \| uint16(frame[2*i+1])<<8)` |
| L7 | `internal/audio/capturelevel_windows_test.go` | `:21-24` | `samples := make([]int16, frames*FrameSamples)` / `samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude` |
| L8 | 同上 | `:48-52`（**期望值是派生不是字面**） | `want := float32(LevelOfSamples(samples[i*FrameSamples : (i+1)*FrameSamples]))` / `if l != want { t.Fatalf("frame %d: level = %v, want %v", i, l, want) }` |
| L9 | `cmd/wisp/resident_audio_247_live_windows_test.go` | `:128-131`（skip 闸） | `if os.Getenv("WISP_LIVE_MIC") != "1" {` / `t.Skip("AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")` |
| L10 | 同上 | `:194-198`（**×2 门槛原句**） | `if loud.mean() <= quietA.mean()*2 \|\| loud.max() <= quietA.max() {` / `t.Fatalf("the loud phase is not distinguishable from silence: %s \| %s \| %s"+` |
| L11 | 同上 | `:190-193`／`:199-201` | `if want := 1 / audio.FrameDuration.Seconds(); loud.perSecond() < want*0.6 \|\| …`／`if quietB.mean() > loud.mean() {` |
| L12 | `cmd/wisp/resident_audio_windows.go` | `:112-120`／`:128-133` | `func (ra *residentAudio) levelOut(out func(float32)) func(float32) {`／`func (rb *residentBall) setAudioLevel(level float32) {` … `rb.b.SetAudioLevel(level)` |
| L13 | 同上 | `:26-30`（注释里唯一写死形状的地方） | `- level (P6 甲 + P8 甲): computed on the capture side with audio.FrameLevel` / `inside the existing audio-capture goroutine and handed out as ONE float32` |
| L14 | 同上 | `:211`／`:213` | `audio.WithLevelSink(ra.levelOut(out)),`／`gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))` |

### 2.2 会被打红的既有钉（名册＋为什么）

⛔ **本表靠读断言原文出，没有跑任何测试**；凡"读码定不了"的都进 §2.3。

**A. 结构事实（决定整张表的形状）：去直流只对"均值非零"的帧有效。**仓里绝大多数电平钉用的是**对称零均值**信号（`±amp` 方波、整周期正弦），减不减直流**数值上等价**⇒ 它们⛔不会红。真正会红的只有载着直流的那几枚：

| 钉 | 位点 | 会不会红 | 为什么／它守的是哪一句 |
|---|---|---|---|
| `TestLevelFullScaleSquareEndpoints`（**1.0 端点那发**） | `internal/audio/level_test.go:74-88` | ⛔ **必红**（读码可判） | `hot[i] = -32768` 整帧恒幅＝**纯直流**：`mean = -32768` ⇒ 去直流后 AC 能量＝0 ⇒ 读 `0.0`，而断言逐字 `if l := LevelOfSamples(hot); l != MaxLevel { t.Fatalf("all-minus-32768 frame = %v, want exactly %v", l, MaxLevel) }`。它守的是 `level.go:25-26` 那句 `level 1.0 exactly: every sample sits at the largest magnitude an int16 can carry (-32768), which is the definition of LevelFullScale`——**这一枚是整张表里唯一"改定义必改判据文字"的硬钉**：1.0 的定义本身就被钉成"恒幅帧"，去直流后它不再是 1.0。 |
| 同一枚测试的**另三发**（`:96` 全幅方波、`:105` ±1 LSB 方波、`:110-114` 半/四分之一幅映射、`:116-125` 不出 `[0,1]`） | `level_test.go:90-125` | ⛔ 不会红 | `squareFrame(amp)` 是 `±amp` 交替（`:21-31`），**mean 恰好为 0** ⇒ 去直流是恒等变换；`FullScaleSquareLevel`/`LevelLSB` 是**常数算术**（`32767/32768`、`1/32768`），不跑函数。 |
| `TestLevelSilentFrameIsExactZero` | `level_test.go:48-67` | ⛔ 不会红 | 全零帧 mean 也是 0；`got != MinLevel` 与负零钉（`:57`）都保持。 |
| `TestLevelSineAgainstRootTwo` | `level_test.go:132-166` | ⚠ **判"不会红"，但属"读码算出、必须真跑确认"** | 被测量是整周期正弦（`:33-36` 注释逐字 `Whole cycles are what make mean(sin^2) exactly 1/2`）；mean≈0，去直流带来的方差减少是 `rms_ac = sqrt(rms² − mean²)`。整数舍入残留的 `mean ≤ 0.5 LSB` ⇒ 位移量级 `0.5²/(2×23170) ≈ 5e-6`，**低于** `SineLevelTolerance = 5e-5` ⇒ 预测在容差内。⛔ 这条是**算出来的**，见 §2.3 T1。 |
| `TestLevelIsLinearInAmplitude` | `level_test.go:172-192` | ⛔ 不会红 | 输入全是零均值（整周期正弦／`±amp` 方波）；线性与单调两发都保持。它守 `ledger:9790` 那句"任何偷偷的曲线都会被打红"。 |
| `TestLevelSineToleranceIsNamedAndBounded`（票面点名的**双边 tolerance 钉**） | `level_test.go:198-212` | ⛔ **不会红（就它自己而言）** | 三发断言逐字是 `SineLevelTolerance >= LevelLSB`、`SineLevelTolerance <= 1e-3`、`gap > 100*SineLevelTolerance`——**前两枚是常数间比较**，去直流不动 `LevelLSB` 也不动 `5e-5`。⚠ **但票面 `:12` 那句"`MinLevel/MaxLevel` 是闭区间端点"值得纠一句**：`MaxLevel = 1.0` 这个端点的**唯一达到方式**就是恒幅帧（纯 DC），所以上面第一枚红与这枚端点定义是**同一件事**，⛔ 不是两枚独立的钉。 |
| `TestLevelSameInputSameBits` | `level_test.go:215-234` | ⚠ 取决于改法（不取决于定义） | 只要 `mean` 用 `int64` 精确求和后定点化，仍"same input same bits"；若直接写 `float64` 均值再平方相减，比较的是 `math.Float32bits` ⇒ 可重现性靠实现细节撑着。见 §2.3 T2。 |
| `TestFrameLevelRejectsNonSeamFrames`／`TestLevelFrameIsOneSeamFrame`／`TestSeamCodecRoundTrip` | `level_test.go:239-271`／`:369-380`／`:276-303` | ⛔ 不会红 | 它们钉的是**帧形状**（`FrameBytes`/`FrameSamples`/小端/奇数长拒绝），与标度无关。`FrameLevel`（`level.go:118-128`）只是转发，不动签名 ⇒ `:268` 那发"整帧静默必须被接受"仍过。 |
| `TestLevelOverBoundedChannelFromWavInjector`（端到端那枚严格 0.5） | `level_test.go:310-364` | ⛔ 不会红 | 夹具是 `±16384` 交替方波（`:314-316`），零均值；`if float64(l) != 0.5` 保持。这枚正是 `ledger:9790` 点名的"端到端那枚严格 0.5"。 |
| **DC 夹具那发**（票面说"改 DC 夹具期望值"） | `capturelevel_windows_test.go:21-24` ＋ `:48-52` | ⛔ **不会红——票面这一格说反了**（见 §4②） | 期望值是 `want := float32(LevelOfSamples(...))` **现场派生**，跟着被测函数一起变 ⇒ 恒绿。C13 全文件 `grep -E "12345\|0\.37\|DC\|expect"` **只命中 `:23`**；C28 全仓 `grep -rn -E "0\.3767\|0\.3684" --include=*.go` **零命中**（只有 `.scratch/wisp/probes/241/v1/mut-*.go` 那几份 `level.go` 的历史突变副本命中 `32768`，属探针、只建不删）。**真实代价不是"要改期望值"，而是"夹具失效"**：去直流后它变成 `want == 0` 对着 `0` 比，`// a steady DC level: RMS is exactly this magnitude` 这句注释同时变假，而它是本票现量里**唯一一枚仓内纯 DC 见证**。⇒ 落点表在这里要加的不是"改数值"而是"**新增一枚显式钉住 DC-读出-0 的用例**，否则'安静读 0'这件事在仓里没有任何见证"。 |
| 球的静默门 | `internal/ball/liquid.go:30` 常数 ＋ `:215 if raw < SilenceLevelGate` ＋ `internal/ball/tokens_table_test.go:552` `"SilenceLevelGate": SilenceLevelGate` | ⛔ **不会红，但语义会静默搬家**（⛔ 本票不许动） | 球侧⛔不调 `audio`（C17 名册：`grep -rn SilenceLevelGate internal/ball/*.go` 的命中全在 ball 包内＋`cmd/balldebug`），所以尺怎么改都不会让 ball 的用例红；`live_windows_test.go:637/648/656/685` 喂的是**字面量 0.9**，`liquid_test.go:60` 比的是 `m.level > SilenceLevelGate`，都与生产者无关。⚠ 但 `0.06` 是 **C21 冻结 token 行**（`c21-native-tokens.md:164`，`PLAN.md:1371`）——去直流后 `0.06` 的**含义**变了（今天 0.06 含直流、真机上永远达不到；改后 0.06 是纯 AC 门限，可能永远开着反方向的问题）。**这一格是"数字不动、单位动"，必须具名交编排者，⛔ 不许写腿顺手重定 `0.06`。** |
| 票 247 的 ×2 门槛 | `resident_audio_247_live_windows_test.go:194-198` | ⚠ **会红吗：预测不会红，但会变得几乎不判事**（必须真跑，§2.3 T3） | 去直流后 `quietA.mean()` 预计趋近 0 ⇒ `loud.mean() <= quietA.mean()*2` 的右边趋近 0 ⇒ **任何正读数都过**。⛔ 这就是"把尺改好，判据自动变松"——**本票禁的是这条**（票面 `:23` ⛔不许为了让 AC#2 绿而放宽 2×）。⇒ 落点表必须写：改尺**必须**同时把 ×2 换成**绝对量级**判据（例：有声窗 AC 读数 ≥ 某 dBFS 门），而**门槛的具体数值⛔不许由写腿自定，要非实现者裁**（`SPEC-12 §4.3` #1）。同文件 `:190-193`（节拍 `[0.6,1.6]×31.25/s`）与 `:187`（非空）不受影响；`:199-201`（`quietB.mean() > loud.mean()` 判"这是噪声不是读数"）在尺变干净后**更容易触发失败**——注意它是**反向**保护：去直流后两窗静默差会变小，这一枚的敏感度**上升**。 |
| 投递那一跳 | `cmd/wisp/resident_audio_windows.go:112-120`／`:128-133` | ⛔ 不会红 | 形状上**只有 `func(float32)` 与 `float32` 类型**，⛔没有任何数值假设；`:115` `ra.levels.Add(1)` 是计数不是量级。唯一写死语义的是 `:26-30` 那段注释与 `:264` 的 verdict 文案（`每帧一枚 float32 交给球`），**改它只动注释、不动判据**。⇒ 具名结论：这一跳**不在最小改动落点表里**（本票不需要碰它）。 |

**B. 落点表枚数（最小改法真要动的文件）**：**1 枚产码（`internal/audio/level.go`：`:96-107` 求和段 ＋ `:17-20`/`:37-39` 标度文档段）＋ 3 枚测试件**（`level_test.go:74-88` 的 1.0 端点那发、`capturelevel_windows_test.go:23` 的注释与失效夹具、`resident_audio_247_live_windows_test.go:194` 的 ×2 判据）**＋ 1 枚必须新增的 DC→0 见证用例**。
⛔ 不动的：`cmd/wisp/resident_audio_windows.go`、`internal/ball/**`（含 `SilenceLevelGate` 与 C21 token 表）、`thresholds.go`、golden、`allowlist.txt`、D43 表、C 表。

### 2.3 读码定不了、必须真跑的（⛔ 本腿一枚都没跑）

| 号 | 定不了什么 | 为什么定不了 | 谁跑 |
|---|---|---|---|
| T1 | 五枚"预测不会红"的钉里，`TestLevelSineAgainstRootTwo`／`TestLevelSameInputSameBits`／`TestLevelIsLinearInAmplitude` 的**实际逐名红/绿** | 我对 `mean≈0` 的位移是**算的**（5e-6 < 5e-5），不是量的；`wholeCycleSine` 的整数舍入残差均值实际是多少，读码读不出来 | 票面 AC#3 自己写的尺：`go test ./internal/audio/ -count=1 -v` 改前/改后逐名红名册作差（此刻 Go 面被 `289-r1` 独占） |
| T2 | 两趟求和后的**位级可重现性**是否真保住 | 取决于 `mean` 的定点化写法（写什么形态都"可能"红或不红），⛔ 属实现细节，只有跑得出 | 同上 |
| T3 | 票 247 `:194` 那枚 ×2 在去直流后的**实际松紧**，以及 `:199` 反向那枚的敏感度是否真会上升 | 要真麦读数（AC#2 那一格），⛔ 不许用合成夹具推 | AC#2 腿（机主在场） |
| T4 | 去直流会不会**顺带改变 `frames_dropped`／节拍**（两趟循环多一倍样本遍历，落在 pinned `audio-capture` 线程上） | `ledger:9753` 那类代价**只能实测**；`resident_audio_247_live_windows_test.go:190-193` 的 `[0.6,1.6]×31.25/s` 窗宽够不够，读码不知 | 落地腿门禁（AC#4 那一发） |

---

## §3 AC#2 — 具名：这一格**没做**

票面 `AC#2` 要的是**同一枚真麦四形现读**（完全安静／机主说话／已知幅度外放／合成 DC＋正弦），各 ≥3 秒并逐形给 `mean/max`。
**本腿零做，三个具名理由**：
1. 它要 `WISP_LIVE_MIC=1` 跑真机用例 ⇒ 撞本腿硬约束（⛔ 零 `go test`；此刻 Go 编译/测试面由写腿 `289-r1` 独占，它的改前基线要求干净 HEAD）。
2. 它要**机主在场**（`resident_audio_247_live_windows_test.go:18` 起那段就写着"nobody is speaking into it, the loud phase reads like the quiet ones"；`:234` 还有 `t.Skipf("no alarm sample to play out loud (%v): AC#2 then needs a person speaking into the microphone", err)` 这一发人闸）。
3. ⛔ 本腿**没有**拿仓内合成夹具或历史日志冒充真机读数。唯一出现的真机数字是从**前腿件里逐字引**的，并明确标它的出处：
   - `.scratch/wisp/probes/247/v1/30-ac2-level-ruler.md:14` ⇒ `resident_audio_247_live_windows_test.go:184: AC#2 reading: quiet-A: samples=93 mean=0.368462 max=0.402494 levels_per_s=31.00`
   - 同件 `:92` ⇒ `Its level is \`12345/32768 = 0.376740\``；`:94` ⇒ `measured floor 0.3685 corresponds to a DC of \`0.368462 × 32768 = 12,074\` LSB — within 2% of that`
   - 同件 `:53` ⇒ `\`sqrt(0.7369² − 0.3685²)\` = **0.6382 FS = −3.90 dBFS at the microphone**. That is clipping`
   ⇒ 票面/派单里那两枚数（0.368462、0.376740）与 −3.9 dBFS 反推**复现自前腿记录、不是本腿现量**，引用时请带这层出处。

⚠ **从 `247-v1` 的记录里额外浮出一枚、AC#2 那一发要当回事的事实**（本腿只是转述，⛔ 不据此判任何东西）：该件 `:44-45` 两行逐字
`| sound − quiet-A | \`sqrt(0.381574² − 0.368462²)\` = **0.09917** FS (−20.1 dBFS) |`
`| quiet-B − quiet-A（纯漂移，中间什么都没放） | \`sqrt(0.382262² − 0.368462²)\` = **0.10178** FS |`
⇒ 那发**外放警报声的 AC 增量（0.09917）比两窗静默之间的纯漂移（0.10178）还小**。这一条与"设备增益/底噪"那一面同源（票面 `:9` 的 max/mean 1.06–1.09 就是它），说明：**即便尺改了直流，两窗静默漂移仍可能吞掉真实信号**；AC#2 那一发要判的"能不能分开"很可能答案是**分开不了**——按票面 `:17` 的原话就得具名写"分开不了"，⛔ 不许挑一个能过的门槛凑绿。

---

## §4 与派单转述不一致之处（以原文为准，具名指出）

① **`docs/specs/SPEC-05-agent-kernel.md` 不存在**。派单让我"在 SPEC-05-agent-kernel.md 里逐字搜"；`git cat-file -e HEAD:…` 报 MISSING、`grep` 报 No such file。HEAD 里的真名＝`docs/specs/SPEC-05-agent-core.md`（`git ls-tree -r --name-only HEAD docs/specs/`，C6）。票面 `AC#1` 只写"SPEC-05"，⛔ 没写错；错在派单的转述。我按真名搜了，结果零命中。

② **票面 `:11` 说"DC 夹具的期望值"要改，代码原文说它不需要改**。票面 `AC#3`（`:18`）也列了"`capturelevel_windows_test.go` 的 DC 夹具期望值"。原文是 `want := float32(LevelOfSamples(...))`（派生，`:48-52`），全文件没有一枚字面 DC 期望值（C13），全仓没有 `0.3767`/`0.3684` 字面（C28）。⇒ 真实代价不是"改期望值"而是"夹具失去见证力＋注释 `// a steady DC level: RMS is exactly this magnitude` 变假"。**以代码为准；这一处要按 §2.2 那行的写法落进落点表。**

③ **派单要我"在票 291 的 `## Progress log` 末尾追加一行"——票 291 没有这一节，且票面 ⛔ 不授权改票本身**。C22：`grep -c -n "Progress log" <票291>` ⇒ **0**（同族旧票如 `01-build-chain-done.md` 确有 `## Progress log`，所以这不是搜索姿势错，是这枚票没建）。票面 `AC#4`（`:19`）逐字：**只读程＝零产码改动＋写点唯一 `.scratch/wisp/probes/291/**`**；票面末行还写着 `⛔ 零翻框、零 push`。
⇒ 按 AGENTS.md §0「未定义即停」与 §1.1 的"以原文为准"，本腿**没有**给票面加 `## Progress log` 标题（那是新建结构，不是追加），⛔ 没动票面任何一行，⛔ 没勾任何 AC 框。派单要的第 3 笔 commit（票面那一行）**未做**，交回编排者裁：要么由编排者给票 291 补上 Progress log 结构，要么承认这一发只交探针件。本件的三笔预算＝`67b4b57d`（骨架）＋正文枚＋（若编排者裁定要票面行）由**有授权的一方**补。

④ **票面 `:12` 把 `MinLevel/MaxLevel` 与 `SineLevelTolerance` 并列成"相邻的既有钉"，读码看它们是同一枚**。`MaxLevel = 1.0` 唯一的达到路径就是一整帧恒幅＝纯直流（`level.go:25-26`＋`level_test.go:74-88`），所以"改去直流"必然动 `MaxLevel` 的端点语义，而 `SineLevelTolerance`（双边那把尺）**不动**。⛔ 别把这一格读成"两枚独立的钉"。

⑤ **本腿没做的事，具名**：票面 `AC#2`（§3）、`AC#4` 的门禁那一发（`GOFLAGS= go build ./...`／`sh scripts/d22scan.sh`／`gofmt -l`／`gofumpt -l`——⛔ 全属 Go 面，此刻归 `289-r1`）。
