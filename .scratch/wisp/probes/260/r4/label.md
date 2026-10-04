# 260-r4 证据件 — 卡片上那行通道标签改成读取时求值（具名解冻 `A598 §2`，撤销口令「260 标签撤回」）

> 本腿的授权来源＝`docs/reports/pending-and-issues.md` 的 **A598 §2**（具名解冻五样齐），前身＝`A595`；
> 射程＝票 260 的 **AC#4 那一格的第十枚**（`internal/agent/approval/approval.go` 的 `ChannelEsc` 标签）。
> 全程未跑 winlive；只跑过 `go vet -tags winlive ./cmd/wisp/`。

## 1. 起手读数与锚点（现跑，⛔ 不是抄派单）

| 项 | 尺原文 | 读数 |
|---|---|---|
| 钟点 | `date "+%Y-%m-%d %H:%M:%S %z"` | `2026-10-04 12:12:44 +0800`（终值采集 `12:26:15`） |
| 起手 HEAD | `git rev-parse --short HEAD` | `1678c9be`（＝编排者那枚 A598 入账枚） |
| 写面脏度 | `git status --porcelain -- internal cmd tools scripts .github docs/evidence` | **空**（0 枚）⇒ 与派单"此刻都空"一致 |
| 未推枚数 | `git rev-list --count @{u}..HEAD` | 起手 `74` → 交件时 `76`（本腿 2 枚） |
| tracked 分母 | `git ls-files \| wc -l` ＝ `5498`；`git ls-tree -r 1678c9be --name-only \| wc -l` ＝ `5495` | 本腿 **+3**（`label.md` ＋两枚新测试件） |

派单点名的四处行号，我**开工前自己重跑**过（尺＝`grep -n ... internal/agent/approval/approval.go` ＋ `grep -rn --include=*.go "channelNames" internal cmd`）：

- `approval.go:87` 逐字 `// channelNames are the display labels the card and the ball strip show.`
- `approval.go:88` 逐字 `var channelNames = map[Channel]string{` ／ `:90` 逐字 `ChannelEsc:   "按 Esc 键",`
- `approval.go:205` 逐字 `st.Text = channelNames[ch] + "：" + unavailableText(ch)` ← **本腿要结的那枚自相矛盾**
- 全仓 `channelNames`＝**11 行**（`approval.go` 的 87/88/160/191/203/205/215 ＋ `gate.go:314` ＋ `gate.go:477` ＋ `report.go:142` ＋ r3 测试件 1 行提及）⇒ **产码读取点＝3 处，与 A598 §2 逐名一致**
- 唯一词面钉＝`internal/agent/approval/cancel_key_wording_260r3_test.go:92` 逐字 `if st.Text != "按 Esc 键"`
- 拼法唯一来源＝`cmd/wisp/resident_approval_windows.go:139 func residentCancelKeySpelling(b *ball.Ball) string`；同文件 `:687` 逐字 `slog.Warn(fmt.Sprintf("approval: L1 窗口挂起期间取消键 %s 未借到，本张卡片无法用 %s 否决", key, key),`
- 那枚 `slog.*（fmt.Sprintf` 宽尺现跑（`grep -rn --include=*.go "slog\.[A-Za-z]*(fmt.Sprintf" internal cmd | grep -v _test`）＝**1 枚＝:687** ⇒ 派单那句"仓内第一枚、0 先例"为真，我复核过。
- `internal/ball/hotkey_windows.go` 最后一次改动＝`git log -1 --format=%h` ⇒ **`b1e59d63`（＝260-r1 那批，不是本腿）**：本腿对 `internal/ball/**` 零写入。

## 2. 授权面与射程核对（`A598 §2` 具名解冻五样逐条对账）

| 五样 | 授权原文要点 | 本腿做了什么 | 尺与读数 |
|---|---|---|---|
| **文件** | `approval.go`（把 `channelNames` 那张 `var` 表改成**读取时求值**的形）＋ `gate.go:314`／`gate.go:477`／`report.go:142` 三处调用点 | 四枚全改，改后落点：`approval.go:106`（表）/`:156`（注入面）/`:167`/`:182`/`:190`/`:204`（求值与两枚标签函数）/`:319`+`:321`（`Statuses()` 两支）；`gate.go:314`、`gate.go:477`、`report.go:147` | `grep -rn --include=*.go "channelNames" internal cmd` 现跑＝**8 行**，其中**产码跨文件读取点＝0**（`gate.go` 零命中、`report.go` 只剩我那段注释）⇒ 三处解冻点都改走了 `channelLabel` |
| **理由** | 条件②"说到做到"覆盖卡片每一行人能读到的字，且今天这一枚会自相矛盾 | 见 §3／§5 格④：未加载那一支现在**不指枚任何键名** | 变异 M2 的红句逐字进 `logs/mutations.txt` |
| **边界** | ⛔ 只动"这枚标签从哪儿取值"；⛔ 不动门／审计／回执任何**判定逻辑**；⛔ 不许新增 `approval → internal/ball` 依赖边；⛔ 不许在 `gate.go` 里再拼一遍组合键 | `gate.go` 只改了两行的 **`%s` 实参**（`channelNames[...]` → `channelLabel(...)`），`g.channels.check(...)`、`q.reject(...)`、`ui.Update(...)` 一字未动（`git diff --numstat`＝`2 2`）；`report.go` 只换 `TextFor()` 里那一枚实参（`6 1`）；**依赖边＝0** | `grep -rn '"github.com/CarlosShao/wisp/internal/ball"' internal/agent/approval` ⇒ **0 命中**（含 `internal/ball` 字样的 8 行全是注释，逐名可见）；组合键字符串**没有任何一处被二次拼出**：`channelLabel` 只做 `"按 %s 键"` 的包裹，实参来自注入 |
| **撤销口令** | 「260 标签撤回」 | 未用；本腿交件＝按此形留在盘上，编排者可凭该口令退回 | — |
| **排程** | `internal/agent/approval`＋`cmd/wisp` 本腿独占；与 `263-v1`（`scripts/**`）不同包可并行 | 起手两包脏度 0；全程只我写 | 见 §7 争用记录（`263-v1` 确实在争用，一枚时序红已归因） |

**同批小账（编排者裁的 A598 §2 末段）**：`slog.Warn(fmt.Sprintf(...))` 改回 attr 形 ⇒ 见 §6，语义与措辞未动（消息文本的必然变化写进那一节，不当零改动报）。

**禁区逐名现跑**（`git diff --name-only 1678c9be..HEAD`，probes 之外＝7 枚）：
`internal/agent/approval/{approval.go,gate.go,report.go,cancel_key_wording_260r3_test.go,cancel_key_label_260r4_test.go}` ＋ `cmd/wisp/{resident_approval_windows.go,resident_cancel_key_label_260r4_windows_test.go}`。
⇒ `internal/ball/**` **零命中**、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）**零命中**、`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` **零命中**。
票 260 框现跑：`grep -c "^- \[ \]"` ＝ **4**、`grep -c "^- \[x\]"` ＝ **1** ⇒ 与派单给的"4 未勾／1 已勾"一字不差，**本腿一枚没碰**。

## 3. 形状选型（读取时求值＋装配根注入；含被排除的形）

**选定的形**：`approval` 侧一枚**包级注入面**——`SetCancelKeySpelling(CancelKeySpelling) CancelKeySpelling`（`approval.go:156`）把一枚 `func() string` 存进 `sync.RWMutex` 守护的槽（`:148-153`），标签在**每一枚读取点**通过 `cancelKeySpelling()`（`:167`）现取现拼。
先例三枚，全部具名：
1. `internal/winsec/resolve.go:128 func SetPathResolver(r C26Resolver)`——"提供方包不许被导入"的包级注入面就是这仓已有的形状；
2. 票 246 的门由装配根注入（`approval.Options{UI, Channels, Logf}`，`cmd/wisp/resident_approval_windows.go:220-227`）；
3. 票 253 那条"panel 为真相源、由装配根注入函数、不新开 `tools → panel` 依赖边"。

**被排除的三形（写清楚，免得下一枚腿再走一遍）**：
- ⓐ `internal/agent/approval` 直接 `import internal/ball` ⇒ 新增包级依赖边＝人工批准面（AGENTS.md §1.2／A598 §2 边界原文"碰到就停手上报"）。**没走**。
- ⓑ 装配根在启动时把一个**字符串**写进包级 `var`（"算一次"形）⇒ 这正是票 260 立票那一形：票 258 的热重载桥（`cmd/wisp/resident_ball_windows.go:313` 的 `ball.NewHotkeyReloader(...)`）会把配置里的键换掉，算一次就陈旧。**没走**——存的是函数，值每次渲染现取（判据③＋变异 M1 的红句为凭）。
- ⓒ 把标签塞进 `Options`／`ChannelRegistry` 实例 ⇒ `gate.go:314`／`:477` 拿得到 `g.channels`，但 `report.go` 那处是**值方法** `CancellationReport.TextFor()`（`report.go:122-124` 明写"Exported on the value so a host can re-render a stored report without the gate's state"），实例化路径拿不到闸门 ⇒ 要么给 D31 的类型加字段（改动面比解冻面更大），要么留第二枚来源。**没走**。
⇒ 一处包级机制覆盖四处读取点，且 `nil` 有明确、被钉住的落点（§5 第五格＋§8）。

**⛔ 无权限声明（写进产码注释 `approval.go:136-141`）**：这枚注入面只回答"人看的那句念哪枚键"，**任何判定都不读它**——可用性由 `ChannelRegistry` 决定、否决要过 `check()`、`allow` 要令牌证明。判据＝§5 的第五格（装了回 `F13` 的读口，未加载通道仍 `ErrChannelUnavailable`、窗口仍不被否决、任何答复都不是 `AnswerAllow`）。

## 4. 取值链（单一真相源怎么走；每步带 `file:line`）

```
[hotkey] cancel (config.toml)
  └─ internal/ball/hotkey_windows.go:588 resolveCancelBorrow   ← 260-r1 已落地，本腿零改动
       └─ :671 cancelBorrowedLineFor / :451 cancelIdleLine     （回执的 cancel 行 Binding ＝交给 RegisterHotKey 的同一枚）
            └─ cmd/wisp/resident_approval_windows.go:139 residentCancelKeySpelling(*ball.Ball)
                 └─ :152 (ra *residentApproval).cancelKeySpelling()  ← 注入的读口优先，否则读球回执；空则 ball.DefaultHotkeys().Cancel(:74)
                      └─ :259 approval.SetCancelKeySpelling(ra.cancelKeySpelling)   ← ★本腿新增的唯一一根线
                           └─ internal/agent/approval/approval.go:167 cancelKeySpelling()
                                ├─ :182 cancelKeyLabel() = fmt.Sprintf("按 %s 键", ...)
                                ├─ :190 channelLabel(ch)         ← 已加载／已发生的通道
                                └─ :204 channelStatusLabel(ch, loaded) ←  availability 行的前半，两支分形
                                     ├─ approval.go:319/:321 Statuses()  → ui.Prompt(p.Channels, ui.go:40) → 卡片／球带／控制台（run.go:1369）
                                     ├─ gate.go:314  L1 窗口否决审计句
                                     ├─ gate.go:477  L2 卡片否决句
                                     └─ report.go:147 回执「（否决通道：…）」
```

- **没有第二条拼键名的路**：`cmd/wisp` 侧仍只有 `residentCancelKeySpelling` 一枚读者（r3 已建），本腿只是把这枚读者**交给 approval**，不是再包一份。
- **默认档那枚拼法今天逐字仍是 `Esc`**：`A588` 条件①钉过 `{Mods:0x4000, VK:0x1B}`；`ball.DefaultHotkeys().Cancel`＝`"Esc"`（前提尺在 `TestTicket260R4ShippedConstructorInstallsTheReader` 第一行）。
- **零漂移指键名不指修饰形容词**：按 `A598 §2` 的重述执行——r3 去掉的「裸」**没有回滚**（`resident_approval_windows.go:692-699` 那段注释原样留着讲理由）。

## 5. 判据四格（⛔ 不合成一句）＋第五格＋红句逐字

新增测试件两枚：`internal/agent/approval/cancel_key_label_260r4_test.go`（321 行，5 枚用例）＋ `cmd/wisp/resident_cancel_key_label_260r4_windows_test.go`（174 行，3 枚用例）；另**有意改动** r3 的一枚既有用例（本节末那格）。

| 格 | 用例 | 断言形状（⛔ 词面尺） |
|---|---|---|
| **① 默认档零漂移** | `TestTicket260R4DefaultLabelIsTheOldSentence`（approval）＋ `TestTicket260R4ShippedConstructorInstallsTheReader` 前半（cmd/wisp） | 未装读口／装的是出厂 `Esc` 时，`Statuses()` 里 `ChannelEsc` 那行**逐字节等于** `按 Esc 键`；另三枚标签（`单击悬浮球`／`面板拒绝`／`说取消词`）逐字在场 |
| **② 改了配置⇒标签跟着变且不再出现 `Esc`** | `TestTicket260R4SeededKeyReplacesEscOnEveryFace`（**四枚读取点逐一**：卡片行／L1 否决句／L2 否决句／回执通道行）＋ cmd/wisp 同枚用例的状态句 | 注入 `Ctrl+Alt+Q` 后四枚文本**都含该串**且**都不含 `Esc`**；红句按面具名（`卡片通道行 …`／`L1 否决审计句 …`／`L2 否决审计句 …`／`回执报告的否决通道行 …`） |
| **③ 热重载之后标签不陈旧** | `TestTicket260R4LabelFollowsTheReaderAcrossAReload` | 读口背后是一枚可动的 holder：第一枚渲染得 `按 Esc 键`，把 holder 换成 `Ctrl+Alt+W` 后**同一对象再渲染一次**必须得新键；并计数"取值函数被调用次数 ≥ 渲染次数"⇒ 任何"算一次"的形当场红。**本包侧能测到的就是这一格；端到端那半〔未实测〕，见 §8** |
| **④ 未加载那一支的组合文本自相一致** | `TestTicket260R4UnloadLineNamesNoKey`（approval，两档都跑）＋ `TestTicket260R4NoBallHostNamesNoKeyOnTheCard`／`TestTicket260R4SeamIsNotAProductionShortcut`（cmd/wisp，走真 `bindBallHost`） | **选型＝乙形：前半不指枚任何键名。** 逐字判据：`标签：不可用句` 的形仍在（`SplitN("：")` 两段），前半＝`取消键通道`，且前半**既不含 `Esc`、也不含配置里那枚键、也不含"按"字**；后半仍逐字 `快捷键取消不可用`（那是 r3 的档一句，本腿不许动）；同一只寄存器在**已加载**时必须反过来（含当前键、不含"不可用"）⇒ 两支不许被后来者折回一支 |
| **第五格（边界）** | `TestTicket260R4SpellingSeamCarriesNoAuthority` | 装一枚回 `F13` 的读口：未加载通道 `Veto` 仍 `ErrChannelUnavailable`、`EventWarning` 里不出现 `F13`、窗口**不被**这次否决答复（非阻塞 `select` 读 `res` 为空）、真正加载了的那枚答复**不是** `AnswerAllow`；读口回空串⇒落回出厂键而不是印出 `按  键` |

**为什么④选乙形、不选甲形（"前半指的就是当前那枚键"）**：甲形要求未加载时也念"当前那枚键"，而这台机器**唯一的键名来源是球自己的热键回执**，回执只在窗口存在时才有；无窗口 ⇒ `residentCancelKeySpelling(nil)` 落回 `ball.DefaultHotkeys()`＝`Esc` ⇒ 甲形正好在"卡片弹不出来那一刻"印回 `按 Esc 键：快捷键取消不可用`＝**把编排者要我结的那枚矛盾原样复制一遍**。乙形与 r3 在 `cmd/wisp` 两枚档一句子（`取消键无处可借`／`取消键通道保持未加载`）同族，也是 A595 §1 边界④那条规则的延续。

**那一枚"改了它就响"的既有钉（`:92`）——⛔ 派单的预测被我的实测推翻，逐字记账**：
- 派单说它会红。**实测不红**：`--- PASS: TestTicket260R3LoadedCancelLabelIsTheNamedResidual (0.00s)`（`logs/approval-targeted.txt`，本腿改完之后跑的）。
- 理由不是它瞎，是**判据①（默认档零漂移）要求默认档那枚串一字不动**，而该用例构造时没有装读口 ⇒ 它读到的一直是同一枚 `按 Esc 键`。⇒ 我**仍然有意改了它**，改的是**判据所在的那一格**而不是期望值：
  - 期望串 `按 Esc 键` **一字未动**（不是放宽）；
  - 新增的是**对偶判据**："装配根装了读口时，这行必须**不再是**那枚出厂串"（`cancel_key_wording_260r3_test.go:94-102`，红句在 `:102`），红句逐字见 `logs/mutations.txt` 的 M1；
  - 同文件的头注释里那段"本包看不见配置、要么导 ball 要么裸写包级表"的**停手上报文字**被换成"编排者已裁（A598 §2）＋落成的形"，因为把那笔已还的账留在注释里就是下一枚腿的假事实来源。
- 结论：**它现在是"默认档钉＋对偶钉"两格，比原来的单格更严**；若后来者把标签算回一次，M1 那发红句会同时抓到本件与 r3 件。

**五发突变（自证有牙，`logs/mutations.txt`，全部还原后 `git status` 该文件为空）**：

| 发 | 改的哪一格 | 结果（终值） |
|---|---|---|
| M1 `cancelKeySpelling()` 变成"只回出厂那枚"（＝标签被算一次／装了等于没装） | ②③④对偶格 | approval `rc=1 FAIL=4 PASS=3`；cmd/wisp `rc=1 FAIL=1 PASS=7`。红句逐字含 `热重载后那行陈旧了：got "按 Esc 键", want "按 Ctrl+Alt+W 键"`、`取值函数被调用 0 次，少于两枚渲染所需的次数`、`改了配置以后卡片那行没念那枚键："按 Esc 键"（应含 "Ctrl+Alt+Q"）` |
| M2 `channelStatusLabel` 不再分形（未加载也指枚键名） | ④ | approval `rc=1 FAIL=1`；cmd/wisp `rc=1 FAIL=2`。红句逐字含 `未加载那行的前半还在指枚键名（含 "Ctrl+Alt+Q"）："按 Ctrl+Alt+Q 键：快捷键取消不可用"` ⇒ **反向判据不是词面尺**：种一枚别的键照样红 |
| M3 `gate.go:314` 退回 `channelNames[...]` | ② 的 L1 面 | approval `rc=1 FAIL=1`，红句 `L1 否决审计句 没有念出装配根那枚键："用户在 L1 确认窗口中通过「取消键通道」否决了本次操作…"` |
| M4 删掉装配根那根线（`:259` 的注入） | 生产接线 | cmd/wisp `rc=1 FAIL=1 PASS=7`，红句 `改了配置以后卡片那行没念那枚键："按 Esc 键"（应含 "Ctrl+Alt+Q"）` ⇒ **r3 的 M4 盲区（"摘掉读球回执那一步，用例全绿"）本腿没有整块复现：装没装读口这一格现在有尺了**；回执**内容**那一半仍〔未实测〕，见 §8 |
| M5 `report.go` 那处退回 `channelNames[...]` | ② 的回执面 | approval `rc=1 FAIL=1`，红句 `回执报告的否决通道行 没有念出装配根那枚键："…（否决通道：取消键通道）…"` |

## 6. 同批小账：`slog.Warn(fmt.Sprintf(...))` 改回 attr 形

| 项 | 逐字 |
|---|---|
| 改前（`cmd/wisp/resident_approval_windows.go:687`） | `slog.Warn(fmt.Sprintf("approval: L1 窗口挂起期间取消键 %s 未借到，本张卡片无法用 %s 否决", key, key), "corr", p.CorrelationID, "why", "取消键位被占用或注册被拒，见热键报告")` |
| 改后（`:712-714`） | `slog.Warn("approval: L1 窗口挂起期间取消键未借到，本张卡片无法用该键否决", "corr", p.CorrelationID, "cancel_key", key, "why", "取消键位被占用或注册被拒，见热键报告")` |
| 窄尺 | `grep -rn --include=*.go "slog.Warn(fmt.Sprintf" internal cmd \| grep -v _test` ⇒ 改前 **1**（＝:687）／改后 **0** |
| 宽尺 | `grep -rn --include=*.go "slog\.[A-Za-z]*(fmt.Sprintf" internal cmd \| grep -v _test` ⇒ 改前 **1**／改后 **0**（"仓内第一枚"这句已被结清，不再有人拿它当先例） |
| **语义与措辞** | 事件、级别（`Warn`）、两条事实（"键没借到"＋"这张卡片没法用那枚键否决"）、`corr`／`why` 两枚属性、键名来源（同一枚 `u.ra.cancelKeySpelling`）**一字未动**；给人看的那句 stdout（`:715` 的 `fmt.Printf("wisp: 卡片 %s 的取消键未借到（桌面已有占位者），按 %s 不会否决它\n", ...)`）**原样保留**，键名照旧念出来 |
| ⛔ 不当零改动报的那处必然变化 | attr 形要求消息是**常量**，所以 `%s … %s` 变成"该键"两字，键名从消息串**移到 `cancel_key` 属性**。这正是 r3 件 §5 结尾写的那枚"另一形漂移更大"的代价，编排者按 A598 §2 末段选了它；我没有偷偷把它写成"没改文案"。 |

## 7. 门禁终值读数（逐名，本机 2026-10-04 12:2x +08，读数存 `logs/`）

| 门 | 命令 | 终值 |
|---|---|---|
| build | `go build ./...` | **rc=0** |
| vet | `go vet ./cmd/wisp/ ./internal/agent/approval/` | **rc=0**（无输出） |
| vet·winlive | `go vet -tags winlive ./cmd/wisp/` | **rc=0**（`logs/vet-winelive.txt` 0 字节＝无诊断）⛔ 只编不跑，全程未跑 winlive |
| 定向·approval | `go test -count=1 -v -run 'TestTicket260' ./internal/agent/approval/` | **rc=0，7 枚具名 PASS／0 FAIL／0 SKIP**（`--- PASS: TestTicket260R4DefaultLabelIsTheOldSentence`／`...R4SeededKeyReplacesEscOnEveryFace`／`...R4LabelFollowsTheReaderAcrossAReload`／`...R4UnloadLineNamesNoKey`／`...R4SpellingSeamCarriesNoAuthority` ＋ r3 两枚仍绿）＝`logs/approval-targeted.txt` |
| 整包·approval | `go test -count=1 -v ./internal/agent/approval/` | **rc=0，top-level `--- PASS`=58、`--- FAIL`=0、`--- SKIP`=1**（含子测试共 77 行 `PASS:`）。唯一 SKIP＝`TestDefaultDeadlineWallClockMeasurement`——**起手即在**（r2/r3 两轮记录过同一枚名），不是本腿造成的 |
| 定向·cmd/wisp | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v -run 'TestAC246\|TestTicket256\|TestTicket260' ./cmd/wisp/` | **rc=0，PASS=28／FAIL=0／SKIP=0**（17.9s；含 260-r3 那 5 枚与 246/256 全套）＝`logs/cmdwisp-targeted.txt`。⚠ 命令形状按派单：`./cmd/wisp/` 带尾斜杠＋带 PATH（不带 PATH 会 `0xc0000135` 且没有 `--- FAIL` 行＝用例根本没跑，本腿没复现那次） |
| 整包·cmd/wisp | `PATH=... go test -count=1 ./cmd/wisp/` | **一枚红**＝`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (41.29s)`，红句 `approval_always_201_test.go:159: ... 审批超时（40 秒未确认），C18 一律判拒绝，已自动拒绝` ⇒ **争用归因，不当缺陷上报、不拿这次读数当终值**：单独重跑同一枚＝`--- PASS: TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card (1.15s)`（`logs/cmdwisp-rerun-always201.txt`，rc=0）。编队里 `263-v1` 正在 `scripts/**` 拉起 `wisp.exe` 真进程；该用例恰是"40s 内要有人答复"的时序形，1.15s→41.29s 这个差不是本腿的代码能解释的量级（本腿零判定路径改动，见 §2 边界行） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；八枚 scope 行逐字 `228/38/85/23/39/85/503/100` ⇒ 与派单给的 12:0x 终值（`internal/=502`、`cmd/=99`）**各 +1**，归因＝本腿新增的两枚测试件（`cancel_key_label_260r4_test.go` 进 `internal/`、`resident_cancel_key_label_260r4_windows_test.go` 进 `cmd/`），与 `A597 §1`／`A598 §1` 记过的同源同理 |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **rc=0、`VERDICT GREEN`**；分母 `tracked paths=5498`（锚点 5495，+3＝本腿）；`over-budget=57／covered=57／not in roster=0`。⚠ 口径按派单写清：`hat budget 165` 是**全路径**量、`(206,217]` 是**墙上相对**量，两者不同尺、不作比较；本次最坏全路径读数 `224 chars` |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l <本腿碰过的 7 枚文件>` | **空输出**（rc=0）。⛔ 不是 `gofmt -l`：`gofmt -l internal/agent/approval cmd/wisp` 另报了 3 枚我没碰的文件（`pending_read.go`／`queue.go`／`models.go`），那是**先于本腿存在**的、与 gofumpt 不同的尺，本腿不动它们 |

**7.1 交件前的复跑（最后一枚测试件改动之后，`2026-10-04 12:37:14 +0800`，本机）**：`go vet ./cmd/wisp/ ./internal/agent/approval/` ⇒ rc=0；`go vet -tags winlive ./cmd/wisp/` ⇒ rc=0；`gofumpt -l <7 枚碰过的文件>` ⇒ **空**（`$(go env GOPATH)/bin/gofumpt.exe`；⚠ 从根目录直接敲 `gofumpt` 在这台机器上不在 PATH＝r3 件里那枚 rc=127 的同一个坑，写出来免得下一枚腿再当成格式红）。

**7.2 交件枚之后再复跑（唯一改动＝§9.1.1 那处幻影注释，`12:41～12:45 +0800`，读数存 `logs/*-postedit.txt`）**

| 门 | 命令 | 复跑读数（现跑） |
|---|---|---|
| build | `go build ./...` | **rc=0**（12:41） |
| vet | `go vet ./cmd/wisp/ ./internal/agent/approval/` | **rc=0、无输出** |
| vet·winlive | `go vet -tags winlive ./cmd/wisp/` | **rc=0、无输出**（⛔ 只编不跑） |
| 定向·approval | `go test -count=1 -v -run 'TestTicket260' ./internal/agent/approval/` | **rc=0，`--- PASS`=7／FAIL=0／SKIP=0**，`ok ... 0.036s`＝`logs/approval-targeted-postedit.txt`（与 §7 那格同名同数） |
| 整包·approval | `PATH=... go test -count=1 -v ./internal/agent/approval/` | **rc=0，`--- PASS`=58／FAIL=0／SKIP=1**，`ok ... 0.376s`＝`logs/approval-pkg-postedit.txt`。唯一 SKIP 仍＝`TestDefaultDeadlineWallClockMeasurement`（起手即在的那枚，§7 已具名） |
| 定向·cmd/wisp | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v -run 'TestAC246\|TestTicket256\|TestTicket260' ./cmd/wisp/` | **rc=0，`--- PASS`=28／FAIL=0／SKIP=0**，`ok ... 17.192s`＝`logs/cmdwisp-targeted-postedit.txt` |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**，八枚 scope 行逐字 `228/38/85/23/39/85/503/100`＝**与 §7 终值同**（`internal/=503`、`cmd/=100` 那各 +1 的归因不变＝本腿两枚测试件）＝`logs/d22scan-postedit.txt`；与 `logs/d22scan-final.txt` 差异仅 6 行子测试计时（0.03↔0.05s 量级），⛔ 不是名集变化 |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **rc=0、`VERDICT GREEN`**，`over-budget=57／covered=57／not in roster=0`，最坏全路径 `224 chars`；分母现跑 `tracked paths=5508`＝`logs/path-length-postedit.txt`。**归因更正（这条差点写错，逐名落尺）**：5498→5508 那 **+10 不全归本腿** 是我为 `263-v1` 兜底的猜测，⛔ 错了——尺一 `git log --since="2026-10-04 12:25" --pretty=%h` ⇒ 窗口内只有 `f8259f5a` 一枚提交；尺二 `git show --name-status --pretty=format: f8259f5a \| grep -c '^A'` ⇒ **恰好 10 枚 A**（九份 `logs/*.txt` ＋ `mutate260r4.py`），全在本腿台件目录 `​.scratch/wisp/probes/260/r4/` 下 ⇒ **那 +10 就是本腿自己的交件枚**，`263-v1` 到此刻零枚进树（尺三 `git status --porcelain -- scripts` 现跑为空）。写清免得下一枚腿把这 +10 读成"别人在动分母"。 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l <7 枚碰过的文件>` | **空输出、rc=0**＝`logs/gofumpt-postedit.txt`（0 字节） |

**7.2.1 复跑没覆盖的门（具名，⛔ 不写成通过）**：`PATH=... go test -count=1 ./cmd/wisp/`（整包）在注释枚之后**没有重跑**——理由不是省事：那一枚改动按尺 `git diff -- internal/agent/approval/approval.go` 现跑读数是 `9 insertions / 4 deletions` 且**逐行都以 `//` 开头**（纯注释，产码零变化），而 §7 那格的整包终值本来就是"重跑同一枚单用例 PASS(1.15s)"之后的判定；`263-v1` 仍在 `scripts/**` 起真进程，此刻重跑整盘只会互洗读数。⛔ 谁要把这句当"整包已绿"用，那是误读：整包绿的是 §7 那次（争用枚除外，已归因）。

## 8. 判不动的地方／量不到的地方（全部具名，标〔未实测〕）

1. **格③的端到端那一半〔未实测〕**：真窗里改 `[hotkey] cancel` → 票 258 的桥 `RebindHotkeys` → 回执的 cancel 行改写 → 标签跟着变，这条链的**球侧**要 winlive（`ball.NewHotkeyReloader` 的 binder 只能是 `*Ball` 或假 binder；假 binder 拿不到 `HotkeyReport` 的私有 `bindings`，构造不出带非默认 Binding 的真回执 ⇒ 用假的就是"用 mock 代替真的"，本腿不做）。**本腿能证到的是这一半**：`approval` 不在自己这一侧缓存标签（调用计数＋前后两枚渲染不同串），以及装配根确实把读口装上了（M4 有牙）。⛔ 不许写成"热重载已验证"。机主那个词（今天能不能真开一次球窗注册一次热键）仍欠，归 `260-v1`。
2. **r3 的 M4 盲区没有整块消失，只消了一格**：现在钉住的是"装没装读口"；"读口读到的那枚值确实等于交给 `RegisterHotKey` 的那枚"仍只在读码层（§4 链的 1-2 步，`internal/ball` 由 260-r1/r2 自钉，本腿零改动）。
3. **`approval` 侧那枚 `defaultCancelKeySpelling = "Esc"`**（`approval.go:121`）是**有意的落点**而不是第二份真相源：它只在"宿主从未装读口"时到达，产线里唯一会加载取消通道的宿主（`resident_approval_windows.go:361`）先经过 `:259` 的注入；漂移由 `TestTicket260R4ShippedConstructorInstallsTheReader` 的第一枚前提尺（`ball.DefaultHotkeys().Cancel == "Esc"`）钉住——ball 的默认一动，那格当场红。
4. **控制台腿（`cmd/wisp/run.go`）不在射程，且按读码不构成矛盾〔读码推导，非实测〕**：它给闸门的是**空寄存器**（`run.go:560` 注释自述"NewChannels() stays empty and runSpec.replyVeto stays unset"；调用点 `run.go:802` 传 `s.replyVeto`），`approval_reply.go:512` 的 `SetLoaded` 被 `vetoChannel != ""` 挡着 ⇒ 取消通道在它那一侧**不加载**，卡片走 `Statuses()` 的未加载支（打印点 `run.go:1369`，`取消方式：%s（%s）`）。按读码，那一行改前是 `按 Esc 键：快捷键取消不可用（未加载）`、改后是 `取消键通道：快捷键取消不可用（未加载）`＝**同一枚矛盾在控制台腿也被同一处消掉了**，但本腿没有真跑过控制台那条腿的那一行（`go test ./cmd/wisp/` 里没有任何用例断言那句，尺＝`grep -rn --include=*.go "取消方式" internal cmd` ⇒ 只有 `run.go:1369` 一枚命中）⇒ 记〔读码〕，不写"已实测"。若将来哪枚宿主**加载了通道却没装读口**，它会印出厂 `按 Esc 键`——那一形今天产线不存在（见 §9.5 最后一条），具名留给 `260-v1` 判要不要做成 fail-closed。
5. **`report.go` 的导出重渲染面有一枚时移**：`CancellationReport.TextFor()` 可被宿主在事后重算（`report.go:122-124`），届时那枚历史回执会按**当时**配置的键渲染。默认档配置不变 ⇒ 无影响；已入账的 `Text` 字段仍是否决当时写下的串（`String()` 优先用它）。要不要把标签冻进结构体＝改 D31 的类型形状，⛔ 超出解冻面，具名交给 `260-v1`。
6. **格③的"调用次数"判据只在本包有效**：它数的是 `approval` 调了几枚读口，不能证明球侧回执被重写。写在这格是为了让"算一次"那一形红得干净，不是热重载证据。
7. **未复跑的门**：全仓 `go test ./...`、`staticcheck`、`slo-check`——派单门禁不含，且此刻 `263-v1` 在 `scripts/**` 起真进程，跑整盘只会互洗读数。⛔ 没跑的一律不写成通过。

## 9. 越界自证＋遗留与欠账（含结尾那一句自问）

**9.1 diff 构成（现跑 `git diff --numstat 1678c9be..HEAD -- internal cmd`）**
```
33   7   cmd/wisp/resident_approval_windows.go
174  0   cmd/wisp/resident_cancel_key_label_260r4_windows_test.go   （新）
120  4   internal/agent/approval/approval.go
321  0   internal/agent/approval/cancel_key_label_260r4_test.go     （新）
39  16   internal/agent/approval/cancel_key_wording_260r3_test.go   （有意改，见 §5 末）
2    2   internal/agent/approval/gate.go
6    1   internal/agent/approval/report.go
```
两枚 commit：`b2e4cf00`（骨架枚，12:1x）、`79c579e2`（落地枚，12:2x）；本件与其 logs／突变台件随后另发一枚。⛔ 未 push。

**9.1.1 一处我自抓的幻影注释（不藏着，单独记）**：`79c579e2` 里 `approval.go:118` 那句注释点名了一枚**不存在的用例** `TestTicket260R4FallbackSpellingStillMatchesBallsDefault`（我写注释时先起了名、后来把两枚真实尺并进了别的用例）。交件前自查（尺＝把本腿碰过的 8 枚文件里出现的所有 `Test...` 名字逐个 `grep "func <名>("` 回数）⇒ 12 枚引用**全部命中真实定义**才算过，那一枚是第一遍没过时修的：注释改成点名真实的两枚（`TestTicket260R4ShippedConstructorInstallsTheReader` 的前提 check ＋ `TestTicket260R3CancelKeyComesFromTheBallChain`），并随 §7.2 那次复跑一起进树。**产码行为一字未动**（纯注释）。

**9.2 纪律自证（尺与读数）**
- `grep -c "t.Skip"` 两枚新件＝**0 / 0**（`tools/d22scan/runtests.sh:98` 把 SKIP 判红，本腿不碰这条路）。
- 没有放宽任何既有断言：唯一被改的既有用例是 §5 末那格，期望串一字未动、多加了一枚对偶判据；M1 的红句证明它现在有牙。
- 凭据值不进对话/日志/表：本件与 logs 里没有任何密钥／路径型凭据（`cancel_key` 属性存的是键名，不是秘密）。
- Git：`add` 与 `commit` 同发、每枚都带显式 pathspec；无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；突变还原用备份回写＋`git status` 核空（`logs/mutations.txt` 每发末行＝"(空＝与 HEAD 一致)"）；临时件只建不删；仓内零 worktree。
- 编队：与 `263-v1` 不同包并行，未碰 `scripts/**`；一枚时序红按派单"先重跑再下结论"处理（§7），没当缺陷上报、也没拿争用读数当终值。

**9.3 明确没做**：票 260 的 AC 框一枚没碰（现跑 4 未勾／1 已勾）；AC#2（丢借用要有声）没顺手做；`internal/ball/**` 没动；票 245 的时机与存在性没动（`bindBallHost`／`TakeEscForCancel`／借还体一字未改）；winlive 一次没跑；`gate.go` 判定逻辑一字未动。

**9.4 交给 `260-v1`／编排者的账（具名）**
1. §8-1 的端到端〔未实测〕两格（真窗改键→回执→卡片那行；真窗 reload 桥→标签不陈旧）——机主那个词一给就能补。
2. §8-4 的"加载了通道却没装读口要不要 fail-closed"＝形选择，不是我该定的。
3. §8-5 的 `TextFor()` 时移＝要不要把标签冻进 D31 结构体（改类型形状，需新的具名解冻）。
4. r3 件 §5 那句"本仓第一枚 `slog message built with fmt.Sprintf`"**已过期**（本腿按 A598 §2 改回 attr 形，宽尺现跑 0 枚）——r3 那份证据件我不改（临时件只建不删），这条更正归编排者入账。
5. 本腿给 `A595 §2` 那枚 winlive 词面钉（`resident_task_source_live_246_windows_test.go:69`）留的旧账**没被触发**：那枚常量钉的句子（`按 Esc 否决了卡片 …`）来自 `vetoDoneLine`，本腿没动它；winlive 仍只 `go vet` 过。

**9.5 结尾那一句自问（必答，答不出来就不交）**

> 我这改完之后，"卡片上那行"和"卡片弹不出来那一刻那行"有没有任何一枚还在指一枚机器没握的键？

**答：没有。** 逐条落到尺上——
- **卡片上那行（通道已加载）**：只可能由 `resident_approval_windows.go:361` 加载，而它前面是 `:259` 的注入 ⇒ 标签必然来自**这枚球自己的回执 cancel 行**（＝交给 `RegisterHotKey` 的同一枚），无窗就没有"已加载"这一支（`bindBallHost` 两枚档一分支：`取消键无处可借`／`取消键通道保持未加载`，均不指枚键名）。
- **卡片弹不出来／键没借到那一刻那行（通道未加载）**：`Statuses()` 走 `channelStatusLabel(ch,false)`＝前半 `取消键通道`，后半 `快捷键取消不可用` ⇒ **前半不指枚任何键名**（M2 红句证明默认档和 `Ctrl+Alt+Q` 档都抓得住）。
- **回执／审计那两行**（`gate.go:314`、`gate.go:477`、`report.go:147`）：能走到它们的否决都必须先过 `check()`（`gate.go:421`）＝通道确实加载 ⇒ 念的是当前那枚键。
- **仍在指一枚出厂 `Esc` 的唯一一路**：宿主从未装读口却把通道标成已加载——今天产线不存在这一路（`grep` 现跑：产码里 `SetLoaded(..., true)` 只有两枚站点，`resident_approval_windows.go:361`（装了读口）与 `approval_reply.go:512`（`vetoChannel != ""` 才走，控制台腿传的是空值＝不走））；这条边界被 `TestTicket260R4SeamIsNotAProductionShortcut` 与第五格钉着，并被 §8-4 具名留给编排者。
- 尺（`grep -rn --include=*.go "按 Esc 键\|取消键通道" internal cmd`）：全仓 20 枚命中，**产码只有 5 枚**，逐名是 `approval.go:108`（表里那枚**通道名** `"取消键通道"`，不是键名）、`approval.go:124`／`approval.go:201`（两行注释，讲的是"旧串长什么样"）、`cmd/wisp/approval_reply.go:504`（注释）、`cmd/wisp/resident_approval_windows.go:354`（档一那句 `取消键通道保持未加载`，指通道不指键）⇒ **产码里没有任何一处再写死"按 Esc 键"这枚字面串**；余下 15 枚全在 `_test.go` 的期望常量与本腿/本包的注释里。出厂默认那枚 `"Esc"` 只活在 `approval.go:121` 那枚具名常量上（§8-3 讲了它为什么不是第二份真相源）。
