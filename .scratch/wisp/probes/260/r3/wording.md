# 260-r3 证据件 — 票 260 条件②剩下的 10 枚「Esc」人看文案（具名解冻 `A595`）

> 授权来源＝`docs/reports/pending-and-issues.md` **A595** §1 的五样（文件／行／理由／边界／撤销口令），
> 票面＝`.scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md` **AC#4**（编排者 11:0x 由 `60767578` 补的那一格，
> 内容＝A595 §1 五样逐字搬，与本腿派单同一口径）。
> 结论先说：**10 枚里落地 9 枚，第 10 枚（`internal/agent/approval/approval.go:90`）按票面那条"拿不到 ball 就停手上报"停下并具名上报，见 §8。**
> 交件时刻：`2026-10-04 11:1x +0800`；本腿三枚 commit＝`85331955`（骨架）／`2ae018be`（九枚＋判据＋winlive 钉同步）／`170e0459`（:584 那一枚形改回句子级）。

## 1. 起手读数与锚点（现跑）

- 时刻：`2026-10-04 10:44 +0800`（`date "+%Y-%m-%d %H:%M %z"`）；起手 HEAD＝`57a33804`（`git rev-parse --short HEAD`，分支 `dev`）。
- 尺①：`grep -nE '"[^"]*Esc[^"]*"' internal/agent/approval/approval.go` ⇒ **2 枚**，逐字：
  - `:90` `	ChannelEsc:   "按 Esc 键",`
  - `:112` `		return "Esc 取消不可用"`
- 尺②：`grep -nE '"[^"]*Esc[^"]*"' cmd/wisp/resident_approval_windows.go` ⇒ **8 枚**，行号 `271/279/287/309/311/495/584/586`（逐字见 §4）。
  ⇒ 与 `A595` §1 记的行号**一致**（那一趟是 `10:19:36`／锚 `84dbda52`；中间落的 `260-r2` 只动 `internal/ball`，没挪这两枚文件的行距）。⚠ 派单前重跑这件事本腿做了，别把这里的号当常量。
- 口径：只数「字符串字面量内部含 `Esc`」的行；注释与标识符（`ChannelEsc`／`vetoByEsc`／`TakeEscForCancel`／`escLoad`）不计。⚠ 这把尺**会被注释里的 `"Esc"` 骗到**（它分不清注释与字符串），所以本腿改完后把注释里的引用一律写成「」形，让终值读数不带噪声：改后 `grep -nE '"[^"]*Esc[^"]*"'` ⇒ `approval.go` **1 枚**（只剩 `:90`）、`resident_approval_windows.go` **0 枚**。
- 起手写面脏度（`git status --porcelain internal/agent/approval cmd/wisp`）＝**空**（只有 `?? .scratch/wisp/probes/260/r3/`，那是我自己在建的证据目录）。⇒ 与派单说的"另有 `212-v3`／`263-r1` 在飞"不冲突：那两枚的写面是 `tools/d22scan`＋台件与 `scripts/slo-check.ps1`，与本腿零交集；终值自查见 §9。

## 2. 授权面与射程核对（A595 §1 五样逐条对账）

| 五样 | 授权原文（要点） | 本腿落点 |
|---|---|---|
| 文件／行 | `approval.go` 2 枚＋`resident_approval_windows.go` 8 枚 | 8 枚全改；`approval.go` 改 `:112`，`:90` 停手（§8） |
| 理由 | 形ⓐ 后借的键吃 `[hotkey] cancel`，句还念 Esc＝条件②禁的那一形 | §3 取值链把念的键接到回执同一行 |
| 边界① 只改指代那枚键的名词，⛔ 不动派发／加载／借还 | 未动任何派发／加载／借还行（尺见 §9 的 diff 构成） |
| 边界② ⛔ 不碰 `gate.go`、⛔ 不碰 `internal/ball/**` | `git diff` 名册里没有这两处路径（§9）；`internal/ball` 只被**读** |
| 边界③ 票 245 的时机与存在性一字不动 | 稳态不绑／只在确认那两三秒借那段代码本腿一行没碰；判据＝本腿不新增任何注册调用（§9） |
| 边界④ `:271`／`:279`／`:112` 诚实形＝「取消键无处可借／不可用」，⛔ 不是换一枚键名 | 三枚都按档一处理，且判据钉住"两支不许指枚任何键名"（§4/§5） |
| 撤销口令 | 「260 文案撤回」＝把这 10 枚改回原句、一枚 commit，⛔ 不许留半改状态 | 撤 = `git revert 170e0459 2ae018be`（两枚，顺序反向），⛔ 别 revert `85331955`（那是证据件骨架） |

⚠ 票面 AC#4 第①格写的是"**那十枚句子**照旧逐字念 `Esc`"，而同一格第三条 bullet 又写"三枚属'通道根本没加载'分支的诚实形是取消键无处可借／不可用"。⇒ 本腿按**更具体的那条**读：① 的零漂移只加在**指枚键名的句子**上，档一三枚的判据换成"不再指枚任何键名"。这一处口径解读写在这里，`260-v1` 若判它读错了，纠正只需要改 §5 的用例名册，不需要改产码形。

## 3. 取值链（单一真相源怎么走）

产码里"念哪枚键"只有一个出口，链是 260-r1/r2 已落地的那条，本腿没新造第二个来源：

1. 配置值进 `internal/ball`：`b.boundCfg` / `b.cancelBinding`（`internal/ball/ball_windows.go:840 ConfiguredHotkeys()`；`:821` 那行 `b.cancelBinding = cfg.Cancel` 是 260-r2 的既有落点）。
2. 借键那遍把配置值**解析一次**：`internal/ball/hotkey_windows.go` 的 `resolveCancelBorrow(bind)`（`:588`）⇒ 返回一枚 `cancelBorrow{Binding, Acc, Fallback}`；**同一个** `borrow` 既交给 `takeEscWithAcc`（⇒ `RegisterHotKey`），也交给回执构造 `cancelBorrowedLineFor(borrow)`（`:671`），后者写进 `HotkeyReport` 的 cancel 那一行的 `Binding` 字段。被拒那一支同理走 `cancelFailedLineFor`（保留尝试过的那枚 `Binding`，Status＝`HotkeyError`）。
3. 稳态那一遍的 cancel 行由 `cancelIdleLine(bind)`（`:451`）填，`Binding` ＝用户配的那枚（配坏了＝Unparsable＋`Err`；没配＝Disabled、`Binding` 为空）。
4. `cmd/wisp` 的读口＝`cmd/wisp/resident_approval_windows.go:139 residentCancelKeySpelling(*ball.Ball)`：只读 `b.HotkeyReport().Bindings()` 里 `Name=="cancel"` 那一行的 `Binding`；`Binding` 为空时落回 **`ball.DefaultHotkeys().Cancel`**（`internal/ball/hotkey_windows.go:74`，就是 `resolveCancelBorrow("")` 用的那枚默认，不是本包抄的字面量）。包装它的是 `:152 (*residentApproval).cancelKeySpelling()`。
5. 八枚句子全部经由第 4 步那一个函数取值，⛔ 本包不再拼任何组合键字符串、不再判"哪枚键才生效"。

两个本腿自己定的小形，具名：
- `:152` 上挂了一枚 **测试读口 `cancelKeyRead`（`:113`）**：`newResidentApproval*` 永远不装它（`TestTicket260R3ProductionPathHasNoSeam` 钉住"产码构造后该字段必为 nil"），它存在的唯一理由是真窗属 winlive（§6）。
- `:363 vetoRejectedLine`／`:367 vetoDoneLine`＝把 `vetoByEsc` 的两句抽出：抽出前它们要活 L1 窗口才走得到（＝真窗＝winlive），抽出后判据不需要窗。**行为与参数一字未变**（`vetoByEsc` 仍返回同样两句、同样以 `ra.cancelKeySpelling()` 填第一枚 `%s`）。这一格在 §9 里按"形改动"如实入账，不当"零改动"。

## 4. 十枚逐句分类（三档）＋改前／改后逐字

行号：改前＝锚 `57a33804`，改后＝`170e0459`。默认档渲染列＝把 `%s` 按 `Esc` 填之后的人看字。

**档二（通道已加载 ⇒ 念当前那枚键；① 逐字零漂移）**

| # | 改前逐字 | 改后逐字 | 这句在说什么事实 | 钉法 |
|---|---|---|---|---|
| 1 | `:287` `ra.residentAuditf("resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；" + "单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）")` | `:350` 同句两个 `Esc` → 两个 `%s`，实参都是同一枚 `key := ra.cancelKeySpelling()` | 装配完的常驻腿：取消通道加载了，而且"本票只落这一条通道" | `TestTicket260R3DefaultWordingIsTheOldSentence`（逐字等旧句）＋`…SeededKeyReplacesEsc`（种 `Ctrl+Alt+Q`） |
| 2 | `:309` `return fmt.Sprintf("按 Esc 否决卡片 %s 被拒：%v", card.CorrelationID, err)` | `:364` `return fmt.Sprintf("按 %s 否决卡片 %s 被拒：%v", key, corr, err)`（由 `vetoRejectedLine` 承担，`vetoByEsc` 传 `ra.cancelKeySpelling()`） | 用户按了取消键、否决被路由拒绝 | 同上两枚（默认档逐字旧句）＋`…ProductionPathHasNoSeam`（不装读口的产码路径端到端） |
| 3 | `:311` `return fmt.Sprintf("按 Esc 否决了卡片 %s（%s / %s）：该调用未执行，答案已入审计", card.CorrelationID, card.Level, card.Tool)` | `:368` `按 %s 否决了卡片 %s（%s / %s）…`，第一枚 `%s` ＝ `key` | 否决成功那句——**就是 winlive 那枚常量镜像的句子** | 同上两枚；端到端那一支属〔未实测〕（§6） |
| 4 | `:495` `return fmt.Sprintf("审批门已装配进本进程（取消通道：Esc 已加载；等待中的确认项：%d）", len(ra.cards.Pending()))` | `:578` `取消通道：%s 已加载`，`%s` ＝ `ra.cancelKeySpelling()` | 启动报告那行 | 同上两枚 |
| 5 | `approval.go:90` `ChannelEsc:   "按 Esc 键",` | **未改**（停手上报，§8） | 卡片／通道名册上那枚"已加载"的标签 | 本腿不钉它；`TestTicket260R3LoadedCancelLabelIsTheNamedResidual` 钉的是**今天的字**，改了要同时改账（§8） |

**档三（借键被 Win32 拒 ⇒ 念"尝试过的那枚键"，仍在加载分支里）**

| # | 改前逐字 | 改后逐字 | 事实 | 钉法 |
|---|---|---|---|---|
| 6 | `:584` `slog.Warn("approval: L1 窗口挂起期间裸 Esc 未借到，本张卡片无法用 Esc 否决", "corr", …, "why", …)` | `:687` `slog.Warn(fmt.Sprintf("approval: L1 窗口挂起期间取消键 %s 未借到，本张卡片无法用 %s 否决", key, key), "corr", …, "why", …)` | 卡片挂起期间借键被拒的告警行 | 〔未实测〕：这句在 `Prompt` 里，要真球窗才进得去（§6）。默认档渲染与旧句**只差一枚词「裸」**，理由与红句见 §5 末 |
| 7 | `:586` `fmt.Printf("wisp: 卡片 %s 的取消键未借到（桌面已有占位者），按 Esc 不会否决它\n", p.CorrelationID)` | `:689` 同一句 `按 %s 不会否决它\n`，多传 `key` | 给人看的那一行 | 〔未实测〕同上；键名取自同一枚 `key`，与 :687 同源同值 |

**档一（通道根本没加载 ⇒ 不指枚任何键名，A595 §1 边界④）**

| # | 改前逐字 | 改后逐字 | 为什么不是换一枚键名 |
|---|---|---|---|
| 8 | `:271` `"四条否决通道全部保持未加载（Esc 无处可借，卡片无处可呈）"` | `:327` `"四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）"` | 这一支连球窗都没有：此刻**没有任何一枚键被借过**，念用户配置的键名＝advertise 一枚按不动的键（＝同支句子自己点名的 B1 禁止形状） |
| 9 | `:279` `"Esc 通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"` | `:336` `"取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）"` | 同上：有窗但装配根没注入取消执行者，通道仍**未加载** |
| 10 | `approval.go:112` `return "Esc 取消不可用"` | `approval.go:119` `return "快捷键取消不可用"` | `unavailableText` 只在**未加载**时被取；新形与另三通道同构（语音取消不可用／面板取消不可用（票 37 未接入）／悬浮球取消不可用），对任何配置都成立 |

档一的连带效果（具名，⛔ 不是本腿新造的洞）：`check()` 把 `unavailableText` 装进 `ChannelError.Msg`，而 `vetoByEsc` 的 `:309`／`:364` 用 `%v` 把它印出来 ⇒ 端到端那句现在读作 `按 Esc 否决卡片 corr 被拒：快捷键取消不可用: esc`（实跑读数，`…ProductionPathHasNoSeam` 记录）。旧形是 `…：Esc 取消不可用: esc`。这一处正是 `TestTicket260R3ProductionPathHasNoSeam` 里那枚"不许还带旧形"的断言在管的。

## 5. 判据（两格）＋红句逐字

新判据七枚（五枚在 `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go`，两枚在 `internal/agent/approval/cancel_key_wording_260r3_test.go`），全部非 winlive、不需要真窗；形状照派单第 4 条指的 `newResidentApproval()`＋句构造器那条既有合法种子建（先例＝`cmd/wisp/resident_approval_246_windows_test.go:338-345`）。

**格①（默认档零漂移）** `TestTicket260R3DefaultWordingIsTheOldSentence`（`:78`）：四枚旧句逐字抄自锚点（`git show 57a33804:cmd/wisp/resident_approval_windows.go` 的 `:287/:309/:311/:495`）存成常量 `old260r3LoadedAudit / old260r3RejectedSaid / old260r3DoneSaid / old260r3StatusLine`，断言**逐字相等**；并且先钉一枚前提尺 `ball.DefaultHotkeys().Cancel == "Esc"`（默认值一动，这格当场红，不许悄悄换基准）。

**格②（种一发别的组合键）** `TestTicket260R3SeededKeyReplacesEsc`（`:117`）：把取消键换成 `Ctrl+Alt+Q`，四枚句子必须**出现那枚新键**且**不再出现 `Esc`**。红句逐字（取自突变 M1 的实跑，`logs/mutations.txt`）：

```
resident_cancel_key_wording_260r3_windows_test.go:130: 否决成功句 没有念出配置里那枚键："按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计"（应含 "Ctrl+Alt+Q"）
resident_cancel_key_wording_260r3_windows_test.go:133: 否决成功句 仍在念 Esc："按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计"
--- FAIL: TestTicket260R3SeededKeyReplacesEsc (0.00s)
```

M2（装配回执回写死）的红句：

```
resident_cancel_key_wording_260r3_windows_test.go:145: 改过键以后的装配句没有念那枚键："wisp: [audit] resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）\n"
resident_cancel_key_wording_260r3_windows_test.go:148: 改过键以后的装配句仍在念 Esc：（同上）
--- FAIL: TestTicket260R3SeededKeyReplacesEsc (0.00s)
```

**档一那一格（⛔ 不并进①②）** `TestTicket260R3UnloadBranchesNameNoKey`（`:159`）：两支未加载分支在**两种档**下都必须"说通道空、不指枚键名"——既不许念 `Esc`，也不许念配置里那枚键。M3（把 `Esc 无处可借` 放回 `:327`）的红句：

```
resident_cancel_key_wording_260r3_windows_test.go:168: 没有球窗那一支没有说「取消键无处可借」（旧句逐字 "四条否决通道全部保持未加载（Esc 无处可借，卡片无处可呈）"）：got "wisp: [audit] …（Esc 无处可借，卡片无处可呈）\n"
resident_cancel_key_wording_260r3_windows_test.go:171: 没有球窗那一支还在指枚一枚键：（同上）
resident_cancel_key_wording_260r3_windows_test.go:190: 未加载的两支指枚了键名 "Esc"：…
--- FAIL: TestTicket260R3UnloadBranchesNameNoKey (0.00s)
```

`approval.go:119` 那一枚的档一尺＝`TestTicket260R3EscUnavailableLineNamesNoKey`（另带三枚对照句"不许被顺带动"尺）；M5（改回旧句）的红句：

```
cancel_key_wording_260r3_test.go:57: 取消通道不可用句 = "Esc 取消不可用", want "快捷键取消不可用"（档一：只说通道不可用，不指枚键名）
cancel_key_wording_260r3_test.go:60: 取消通道不可用句还在指枚 Esc："Esc 取消不可用"
--- FAIL: TestTicket260R3EscUnavailableLineNamesNoKey (0.00s)
```

反向判据**没有做成词面尺**：⛔ 没有"grep 文件里还剩几个 Esc 字样"这种断言；上面每一枚都是"渲染出来的句子必须念哪枚键／不许念哪枚键"，并且 M1／M2／M3／M5 各打红一发证明它们有牙。

**M4 与它照不到的那一半（⛔ 不写成已验证）**：M4＝把 `residentCancelKeySpelling` 里"读球回执行"那一步短路（`if false && bd.Name=="cancel"`），让它永远落回默认常量 ⇒ **五枚用例全绿**（`--- PASS` ×5，rc=0）。⇒ 这就是本腿判据的真实边界：**"链上球侧那半"在本腿量不到**，能执行的只有"标签→句子"那半。球侧那半的依据＝读码（§3 的第 2／3／4 步逐处 `file:line`）＋ winlive 那枚常量（§6），结论一律标〔未实测〕，交给 `260-v1` 用真窗那一路补跑。

**`#6`（:584）的单字漂移**：旧句 `裸 Esc` 里的「裸」＝"没有任何修饰键"，用户配 `Ctrl+Alt+Q` 时那两个字本身是新的一句谎（同一族缺陷），所以「裸」必须走；默认档渲染因此与旧句**差两个字**（`approval: L1 窗口挂起期间取消键 Esc 未借到，本张卡片无法用 Esc 否决` vs 旧 `…期间裸 Esc 未借到…`）。另一形（键名进 attr、消息常量不变）漂移更大，被换掉了；代价是本仓第一枚 `slog.Warn(fmt.Sprintf(...))`（尺＝`grep -rn "slog\.[A-Za-z]*(fmt.Sprintf" --include=*.go internal cmd` ⇒ 改前 0 枚，改后 1 枚＝这一处）。这一处由编排者裁：要 attr 形就一句改回，产码其余不动。

## 6. winlive 那枚钉的同步（`vetoDoneClaim`）＝判据，不是实测

- 改前逐字（锚 `57a33804`，`cmd/wisp/resident_task_source_live_246_windows_test.go:69`）：`const vetoDoneClaim = "按 Esc 否决了卡片"`
- 改后逐字（现 `:87`）：`var vetoDoneClaim = "按 " + ball.DefaultHotkeys().Cancel + " 否决了卡片"`，使用处 `:149`（`pollUntil127(200, func() bool { return leg.stdout.has(vetoDoneClaim) })`）一字没动。
- 为什么跟着默认常量：台架 `prepareResidentHarness` 写的 `config.toml`（同文件 `:324-347`，逐字从 `schema_version = 2` 到 `allowed_dirs = …`，**不含 `[hotkey]` 段**——尺①＝`sed -n '326,347p' <该文件> | grep -c "\[hotkey\]"` ⇒ **0**；尺②＝`grep -n "\[hotkey\]" <该文件>` ⇒ 2 命中，两枚都是本腿刚写进 `vetoDoneClaim` 注释里的那句"no [hotkey] section"，没有一枚是配置项）⇒ 该台架跑的就是默认档 ⇒ 与 `resolveCancelBorrow("")`/`ApplyHotkeyDefaults` 同一枚来源；如果这枚常量继续写死 `"按 Esc …"`，那它是在**预言**产码不会再变，而不是在**跟着**产码——正是要避免的"必谎的尺"。
- ⚠ 这一格只能是判据〔预测，非实测〕，三条尺本腿全部复跑：
  1. CI：`grep -nE "tags|winlive" scripts/wisp-cli-tests.sh .github/workflows/ci.yml` ⇒ 本腿复跑 **0 命中**，两枚文件都在（`ls -l` ＝ `scripts/wisp-cli-tests.sh` **45,198** 字节、`.github/workflows/ci.yml` **6,429** 字节）⇒ 0 是尺的结论、不是尺没跑（与 `A595` §2 那两枚字节数一致）；
  2. 本机不带 sherpa PATH：`go test -count=1 -run TestAC246ChannelNeedsBothWindowAndExecutor ./cmd/wisp/` ⇒ 本腿**当场复现** `exit status 0xc0000135` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.044s`，输出里**没有 `--- FAIL` 行**（票 98 在册那一族死法：加载期缺 DLL，用例根本没进）⇒ 今天也量不到；
  3. 本腿**没有跑过任何 winlive**：只跑了 `go vet -tags winlive ./cmd/wisp/`（rc=0，证编得动，读数在 §7）。
  ⇒ 所以"这枚常量现在跟句子同源、不会红"是**读码判断**。预测内容：真窗那一路若跑，`:149` 仍应等到这句；若某日它等不到，先查是不是台架开始自己写 `[hotkey] cancel` 而没同步这枚常量（改后注释里已把这句话写进文件）。
- 另核一枚非 winlive 的近邻：`cmd/wisp/resident_approval_246_windows_test.go:343`（tag 只有 `windows`）钉的是「没有可否决的确认项」那句 ⇒ 不在 10 枚射程内；本腿改后该用例仍绿（它在 §7 那 39 枚 PASS 里，名册见 `logs/targeted-cmdwisp-final.txt`）。

## 7. 门禁终值读数（逐名，2026-10-04 11:1x +08，本机）

| 门 | 命令（逐字） | 终值 |
|---|---|---|
| 构建 | `go build ./...` | rc=0 |
| vet | `go vet ./cmd/wisp/ ./internal/agent/approval/` | rc=0 |
| vet winlive（只证编得动，⛔ 未跑） | `go vet -tags winlive ./cmd/wisp/` | rc=0 |
| 定向（带 PATH） | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v -run 'TestAC246|TestTicket256|TestTicket260' ./cmd/wisp/` | **rc=0**，top-level `--- PASS`=39、`--- FAIL`=0、`--- SKIP`=0；本腿五枚全 PASS（`logs/targeted-cmdwisp-final.txt:96,101,111,114,119`）；`ok github.com/CarlosShao/wisp/cmd/wisp 17.7s` |
| approval 包全量 | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./internal/agent/approval/` | rc=0，`--- PASS`=72、`--- FAIL`=0；唯一 `--- SKIP: TestDefaultDeadlineWallClockMeasurement`＝`internal/agent/approval/ticket84_no_owner_test.go:222` **起手即在**的那枚（非本腿造，本腿没动它） |
| D22 扫描 | `sh scripts/d22scan.sh` | **rc=0 clean - no D22 ban violations**；`ban #8 internal/=502`、`ban #8 cmd/=99`（编排者 10:40 终值是 501／98 ⇒ 两枚 +1 恰是本腿新增的两枚测试件文件本身，不是新增违规）；`bans #1-5 internal/=228`、`cmd/=38` 不变 |
| 路径预算 | `sh scripts/check-path-length-budget.sh` | rc=0，逐字 `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；`longest=180 chars relative`；`tracked paths=5442`（编排者 10:40 是 5316 ⇒ 增量为今天别的腿进树的路径，与本腿无关）；`over-budget=57 covered=57 not in roster=0`、`in the wall interval=0` |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt" -l <碰过的五枚文件>` | **空**（输出 0 字节，`logs/gofumpt-final.txt`）。⚠ 同目录那枚 `logs/gofumpt.txt` 里是 50 字节的 `command not found`＝本腿第一次直接敲 `gofumpt` 的 rc=127（它装在 `$(go env GOPATH)/bin`，不在这条 shell 的 PATH 上），**不是格式红，也不是门红** |
| ⛔ 未跑 | `slo` 全量／`-Subset full` | 未跑（本机 self-hosted runner 的活，派单明令避开）；winlive 全部未跑 |

**7.1 进树后复跑（证据件 `0de321fa` 落盘之后，`2026-10-04 11:21:34 +0800`，读数存 `logs/path-length-postcommit.txt`／`logs/d22scan-postcommit.txt`）**：
`sh scripts/check-path-length-budget.sh` ⇒ **rc=0**、逐字 `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`、
`tracked paths=5460`（比上表那枚 5442 多 14＝本腿刚进树的 14 枚读数与台件本身，⚠ 上面 §7 表里那行是**进树前**的读数）、`longest=180 chars relative`、`over-budget=57 covered=57 not in roster=0`、`in the wall interval=0`；
`sh scripts/d22scan.sh` ⇒ **rc=0**、逐字 `clean - no D22 ban violations`、`ban #8 internal/=502`、`ban #8 cmd/=99`（与上表一致：`.scratch/**` 不在 d22scan 的扫描射程，进树不改分母）。
⇒ **真终值以本节为准**；`go build`／两档 `go vet`／`gofumpt -l`／两包定向测试在本腿最后一枚 commit 之后未再改动产码，故沿用上表读数。

## 8. 判不动的地方／量不到的地方（含 `:90` 的停手上报）

**8.1 唯一没落地的第十枚＝`approval.go:90` `ChannelEsc: "按 Esc 键"`（停手上报，等编排者裁）**

票面 AC#4 与派单都写了"拿不到 ball 就停手上报，⛔ 不许自行决定加依赖"。本腿复量：这枚标签在包外有**四处无锁直读**——`internal/agent/approval/approval.go:203`／`:205`（`Statuses()` 里"已加载念标签、未加载念标签：不可用句"那两行）、`gate.go:314`、`gate.go:477`（两句否决审计「通过「%s」否决了本次操作」）、`report.go:142`（「否决通道：%s」）。⇒ 三形与代价：

- **甲＝装配根注入**（`cmd/wisp` 已经把同一枚 `key` 交进 `approval`，本包只把它写进 `channelNames[ChannelEsc]`）：⛔ 但 `gate.go:314`／`:477`／`report.go:142` 直读那张 map、没有任何锁 ⇒ 写它要么需要把这三处改走一枚带锁的 getter（**越出 A595 §1 的写面**：边界②点名不许碰 `gate.go`，`report.go` 也不在名册里），要么接受"装配期写一次、运行期不再写"——而"运行期不再写"正好在 258 的 `[hotkey]` 热重载路径上变陈旧 ⇒ 用户改完键，卡片标签还念旧键＝**票 260 立票的那一形重演**。⇒ 甲不是一枚腿能干净落下的形。
- **乙＝换成语档无关的通用形**（例：「按取消键（[hotkey] cancel）」）：不动依赖、不加锁，但**默认档不再逐字念 `Esc`** ⇒ 直接撞 AC#4 ① 的零漂移格。
- **丙＝今天这样不改**（本腿选的）：留一枚已知会谎的标签，但**具名入账＋钉住**：`internal/agent/approval/cancel_key_wording_260r3_test.go` 的 `TestTicket260R3LoadedCancelLabelIsTheNamedResidual` 钉的是**今天的字**，句里写明"这枚是欠账尺、不是批准"，谁改了标签它就响，逼改的人同时改 §8 这笔账。

⇒ **要人拍的那一句**：甲要不要开（开＝同时解冻 `gate.go` 那两行＋`report.go` 那一行，或批准"装配期一次写＋热重载时再写"并接受热重载前的窗口期）；还是乙（同意 ① 对这一枚让路）；还是丙继续挂着。**撤销口令不变（「260 文案撤回」不涵盖 :90，因为它没改）。**
连带影响具名：`:90` 未改 ⇒ 档一那句未加载句（`approval.go:205` 拼出来的）今天仍读 `按 Esc 键：快捷键取消不可用`——后半已经不说谎，前半还说。**⇒ 边界④那三枚里，`approval.go` 这一枚只能做到"一半"，另一半卡在 :90。** 这一格是本腿交件里最该被攻的地方，别当它不存在。

**8.2 量不到的（全部〔未实测〕，指名怎么补）**

| 量不到的东西 | 为什么 | 谁能量 |
|---|---|---|
| `:584`／`:586`（借键被拒那两枚句子）的渲染 | 在 `ballCardUI.Prompt` 里，要真球窗 | winlive（机主给词）或 `260-v1` 造真窗台架 |
| "球回执 → `cmd/wisp` 念的键名"这一段（§3 第 4 步） | M4 实证：摘掉它本腿五枚用例全绿 | winlive：真改 `[hotkey] cancel` 后看 `:311`／`:495` 那两句念什么 |
| `vetoDoneClaim` 改了会不会让那发 winlive 用例由绿变红／由红变绿 | CI 无 `winlive` tag＋本机 `0xc0000135`（§6 三把尺） | 同上 |
| 真机上"改了 `[hotkey] cancel` 之后卡片与回话念同一枚键" | 需要真窗＋真注册热键 | 机主给词后的 winlive |
| `approval` 包那枚 SKIP 是不是被本腿挪动 | 本腿只跑了包全量一次，`ticket84_no_owner_test.go:222` 起手即在 | 验收腿比 `git diff` 即可判（本腿未碰那枚文件） |

**8.3 判据自身的边界**：`cancelKeyRead` 读口能证明"句子跟着标签走"，⛔ 不能证明"标签跟着球走"（M4 就是这件事的红字）。`TestTicket260R3ProductionPathHasNoSeam` 把读口在生产路径上钉成 nil，是防这一枚种子面变成产码捷径的那道闸。

## 9. 我越界了什么＋遗留与欠账

**9.1 名册与 diff（现跑）**：`git diff --numstat 57a33804..HEAD -- <三枚改动的产码/测试件>` 的本腿构成＝
- `cmd/wisp/resident_approval_windows.go`（改）：10 枚人看句子里的 8 枚＋取值口＋两句抽出＋读口字段。
- `internal/agent/approval/approval.go`（改）：1 枚（`:112`→`:119`）＋注释。
- `cmd/wisp/resident_task_source_live_246_windows_test.go`（改）：常量形同步＋import 一行（`internal/ball` 在本包产码里早已 import ⇒ **不是新包间依赖边**）。
- 新增：`cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go`、`internal/agent/approval/cancel_key_wording_260r3_test.go`。
- 终值尺：`grep -cE '"[^"]*Esc[^"]*"'` ⇒ 两枚文件合计 **1** 枚（就是 `:90` 那枚欠账）。
- ⛔ 名册之外零路径：`git diff --name-only 57a33804..HEAD` 里 `gate.go`／`report.go`／`internal/ball/**`／三枚冻结件／`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` **一枚都没有**（`tools/d22scan`、`scripts/slo-check.ps1` 也不在：那是 `212-v3`／`263-r1` 的写面）。

**9.2 三处"没超出射程但确实改了形"的自陈（不当零改动报）**：
1. `vetoRejectedLine`／`vetoDoneLine` 两枚函数＝为判据抽出，行为与参数一字未变；若验收判这属"为测试造形状"，撤法是把两句内联回去、并放弃 §5 里那两枚句子级用例（会同时失去 M1 那发牙）。
2. `residentApproval` 新增 `cancelKeyRead func() string` 字段＝测试种子面，产码构造路径恒 nil（有尺钉）。
3. `slog.Warn(fmt.Sprintf(...))`＝仓内第一枚该形（0 先例 → 1），理由与替代形写在代码注释里，可一句换回 attr 形。

**9.3 明确没做的**：AC 框一枚没碰（票 260 现 4 枚未勾：AC#1/AC#2/AC#3/AC#4，翻勾归编排者）；没有 `t.Skip`；没有放宽任何既有断言（`vetoDoneClaim` 是**跟着句子改**不是拆断言——它等的仍是同一句默认档渲染）；没跑 winlive；没跑 `-Subset full`。

**9.4 交给 `260-v1`／编排者的账（具名）**：
1. `:90` 三形（甲／乙／丙）拍一枚——**这是本腿唯一的未落地格**，且它使档一那枚未加载句只剩"一半诚实"（§8.1）。
2. §2 末那处口径读法（AC#4 ① 的"十枚"vs 边界④的"三枚"）请裁一句：本腿按"更具体那条"执行了。
3. :584 那两个字「裸」的单字漂移（§5 末）请裁：接受、还是换 attr 形。
4. 〔未实测〕那张表（§8.2）四格，只有真窗能销。
5. M4 那发"摘掉读球那一步仍全绿"＝本腿判据的真实盲区，验收腿若要攻文案链，攻击点应该在这里，不在 §5 那七枚里。
