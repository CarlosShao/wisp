# 票 260 — 非实现者验收腿 `260-v1` 对抗验收表

> 本腿＝`260-v1`，未参与 `260-r1`／`r2`／`r3`／`r4` 任何一程。本表只放本腿自己现跑的读数与本腿的判语。
> 起手时刻 `2026-10-04 13:16:10 +08`，起手 HEAD＝`515ca5c5`（`dev`，ahead 81）；交件时刻 `13:4x`，其间 HEAD 被别人推到 `0b938773`（`a69b19ef`＝A599、`35e852f9`＝立票 265、`0b938773`＝265-a1 骨架），
> 本腿全程零产码改动、零票面改动（票 260 面最后一条 commit 仍是 `60767578`，见 §5）。
> ⛔ 本腿不翻任何勾；四枚框的翻勾全部归编排者。
> 落点＝`.scratch/wisp/probes/260/v1/verdict.md`；读数档同目录 `logs/`；两台突变脚本在本目录。
> **凭据标注口径**：本节以下每一格都写清是本腿自跑（〔自跑〕）还是转述别人（〔仅它的读数，本腿未复跑〕）。
> 四程证据件（`probes/260/r1/landing.md`、`r2/sync.md`、`r3/wording.md`、`r4/label.md`）在本腿眼里全部是**待验断言**，没有一句被当凭据引用。

---

## 一、起手台本与写面基线

| 项 | 读数 | 取法 |
|---|---|---|
| 起手钟点 | `2026-10-04 13:16:10 +0800` | 〔自跑〕`date` |
| 起手 HEAD | `515ca5c5`（`## dev...cnb/dev [ahead 81]`） | 〔自跑〕`git rev-parse`／`git status --porcelain --branch` |
| 写面起手态（`-- internal cmd tools scripts .github docs`） | **空（0 行）** | 〔自跑〕`logs/baseline-status-internal-cmd.txt`（`wc -l`＝0） |
| 关键六枚文件起手 md5 | `resident_approval_windows.go 7a26c7a990dbd2351bdf9898b5bdc192`／`approval.go 9326bd2b…050`／`gate.go 8d93476e…618`／`report.go 5a753677…16d`／`hotkey_windows.go be79cbda…ba`／`ball_windows.go 2488b2ac…d7` | 〔自跑〕`logs/baseline-md5-keyfiles.txt` |
| 票 260 框现量 | **4 未勾／1 已勾** | 〔自跑〕`grep -c '^- \[ \]'`＝4、`'^- \[x\]'`＝1 |
| 全树脏文件（不归本腿） | `design/**` 16 枚删除（含 `design/assets/tokens.css`）、`.gitignore`、`probes/**` 旧日志 | 〔自跑〕`git status --porcelain`；⛔ 本腿未提交、未还原任何一枚 |
| 写面交件态（`-- internal cmd tools scripts .github docs`） | **空**，且六枚 md5 逐枚等于起手值 | 〔自跑〕§2／§4 每次突变后的还原核对 |

⚠ 环境态对本腿读数有一处直接影响：`design/assets/tokens.css` 在工作树里是**被删**的（别人所为），所以 `internal/ball` 整包里
`TestC21TableColourRowsMatchTokensCSS` 恒红（`tokens_table_test.go:1468`，红句见 `logs/gate-test-ball-verbose.txt`）。
本腿把这枚红**归口为环境／别人的写面**，⛔ 不计入票 260 的账，也不去修它（见 §9 记我 J-3）。

---

## 二、★头一格：`260-r3` 自报的 M4 盲区——本腿亲手重跑那一发突变

### 2.1 本腿做了什么

突变体＝`cmd/wisp/resident_approval_windows.go:139 residentCancelKeySpelling` 里**"读球回执"那一步整块摘掉**：
原体（现读 `:139-148`）

```go
if b != nil {
    for _, bd := range b.HotkeyReport().Bindings() {
        if bd.Name == "cancel" && bd.Binding != "" {
            return bd.Binding
        }
    }
}
return ball.DefaultHotkeys().Cancel
```

被换成 `_ = b; return ball.DefaultHotkeys().Cancel`。脚本＝`mutate_m4_v1.py`（在 finally 里从**内存里的原始字节**还原）。

### 2.2 三条的读数

**① 摘掉之后是不是真的全绿 ⇒ 是的，全绿，而且不止它说的五枚。**〔自跑〕`logs/M4-run-targeted-260.txt`、`logs/M4-run-full-cmdwisp.txt`

- 定向 `-run TestTicket260 ./cmd/wisp/`：**rc=0，8 枚全 PASS、0 FAIL、0 SKIP**——
  `TestTicket260R4ShippedConstructorInstallsTheReader`／`R4NoBallHostNamesNoKeyOnTheCard`／`R4SeamIsNotAProductionShortcut`／
  `R3DefaultWordingIsTheOldSentence`／`R3SeededKeyReplacesEsc`／`R3UnloadBranchesNameNoKey`／`R3CancelKeyComesFromTheBallChain`／`R3ProductionPathHasNoSeam`。
  ⚠ `r3` 自报的是"本批那**五枚**"，本腿数到的是 **8**（`r4` 的三枚也在摘掉之后全绿）。
- 整包 `go test ./cmd/wisp/`：**rc=0，`ok … 194.642s`**（不是 `0xc0000135`＋`0.0xxs` 那枚"根本没跑"的形状——那两发本腿都现量过，见 §7.4）。
- 还原核对：`resident_approval_windows.go` md5 回到 `7a26c7a9…192`，`git status --porcelain -- internal cmd` 空。

**② 如果红，指名哪枚用例红、红句逐字 ⇒ 本发不红；本腿另外跑了三发对照，全部按预期红，所以②的"全绿"不是瞎尺的产物。**〔自跑〕`logs/M4b-controls-console.txt`

| 发 | 突变内容 | 结果 | 红句逐字 |
|---|---|---|---|
| **M4** | 摘掉读球回执那一步 | **rc=0 全绿**（8／8） | —（这就是盲区本身） |
| **M4b** | 读回执的循环保留，但把 `bd.Name == "cancel"` 改成 `== "summon"`（＝读**错一枚键**的回执） | **rc=0 全绿**（定向 8／8；整包 `ok … 188.089s`） | —（比 M4 更糟：连"读错行"都不响） |
| **CTRL** | 删掉装配根那根注入线（`:259 approval.SetCancelKeySpelling(ra.cancelKeySpelling)`） | **rc=1，FAIL=1／PASS=7** | `resident_cancel_key_label_260r4_windows_test.go:104: 改了配置以后卡片那行没念那枚键："按 Esc 键"（应含 "Ctrl+Alt+Q"）`＋`:107: 改了配置以后卡片那行还在念 Esc："按 Esc 键"` ⇒ **`r4` 的接线格有牙，这句我复跑复认** |
| **CTRL2** | 把 `vetoDoneLine` 的键名重新写死成 `"Esc"` | **rc=1，FAIL=1／PASS=7** | `resident_cancel_key_wording_260r3_windows_test.go:130: 否决成功句 没有念出配置里那枚键："按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计"（应含 "Ctrl+Alt+Q"）`＋`:133: … 仍在念 Esc …` ⇒ **`r3` 的 ②格有牙，复认** |
| **D** | `internal/agent/approval/approval.go:126 defaultCancelKeySpelling` 由 `"Esc"` 改 `"F13"` | **rc=1，FAIL=3** | `cancel_key_label_260r4_test.go:93: 默认档已加载的取消通道标签漂移：got "按 F13 键", want "按 Esc 键"`＋`TestTicket260R4SpellingSeamCarriesNoAuthority`＋`TestTicket260R3LoadedCancelLabelIsTheNamedResidual` ⇒ **approval 侧那枚重复的出厂默认字面量不是无人看管的第二来源** |

CTRL／CTRL2／D 三发都红 ⇒ **本腿这把尺抓得到人尽皆知的真值**，于是 M4／M4b 的"全绿"是一枚**真实测量**，不是我的驱动坏了。
（⚠ 本腿自己的第一版驱动**确实坏过一次**，把产码留在突变态：见 §9 记我 J-1。）

**③ 全绿的情况下，AC#1「配置值真被读到」这格的凭据到底落在谁身上——具名。**

这格不是一枚整体，本腿把它拆成两段分开落名（这也就是 `r3`／`r4` 都没写满的那半句）：

- **段一「配置值 → 球自己的回执行」＝有钉，而且本腿现跑全绿。** 兜它的具名钉＝
  `internal/ball/hotkey_cancel_borrow_260_test.go` 的 `TestBorrowFollowsConfiguredBinding260`／`TestBorrowReceiptSharesOneSourceWithRegistration260`／`TestBorrowDefaultIsBitIdentical260`／`TestBorrowUnparsableFallsBackLoudly260`
  ＋ `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go` 的 `TestBorrowWishFollowsSeed260r2`／`TestBorrowWishDefaultCellIsBitIdentical260r2`／`TestBorrowWishRefusedSeedIsNotSilent260r2`
  ＋ 既有 `TestCancelBorrowFailureIsAProblemLine`／`TestConfiguredCancelStillNeverBoundWhileIdle260`。本腿 §7 整包跑：**10／10 PASS**（`logs/gate-test-ball-verbose.txt`）。
  ⚠ 这枚钉的射程止于 `HotkeyReport` 的**内容**，它不经过 `cmd/wisp`。
- **段二「球回执 → `cmd/wisp` 念给人看的那句」＝没人兜。** M4 与 M4b 两发（摘掉读／读错行）都杀不死任何用例 ⇒
  `residentCancelKeySpelling` 里那个循环在本仓的**可执行判据射程之外**；唯一碰到这枚函数的具名钉
  `TestTicket260R3CancelKeyComesFromTheBallChain`（`:201`／`:204`）**只喂 `nil`**，走的是 `if b != nil` 之外的 fallback 支，
  所以它对"读回执"这件事的敏感度＝零。其余 `r3`／`r4` 用例全部经 `cancelKeyRead` 读口种子（`ra.cancelKeySpelling()` 第一行就返回读口），
  本腿逐枚读码确认：**没有任何一枚用例用非 nil 的 `*ball.Ball` 调过 `residentCancelKeySpelling`**（现量尺＝
  `grep -n "residentCancelKeySpelling(" cmd/wisp/*.go` ⇒ 产码 1 处调用（`:156`）＋测试 2 处（都在 `nil` 上））。
- ⇒ **判语落点**：AC#1 的"配置值真被读到"**在球侧成立**（有钉、现跑绿），**在"读给用户看"那半不成立**（本腿两发突变实证无人兜）。
  这不是"实现者说了不算"——是**本腿自己摘了一遍、绿了、并且用三发对照证明过我的摘法不是空操作**。
  `r4` 件 §8 第 2 条自己写了"r3 的 M4 盲区没有整块消失，只消了一格（现在钉住的是装没装读口）"，本腿复跑后**这句为真**；
  但 `r4` 的测试件头部写的是"the M4 blind spot of r3 was about this exact wiring, from the other side"——**这句在盘上偏硬**：
  它堵的是"注入线"那一格（CTRL 复认有牙），没有堵"读回执"那一格（M4 复跑仍全绿）。具名差一处，归 §8 的条件。

---

## 三、AC#1（形ⓐ：借用那遍真读 `[hotkey] cancel` 的配置值）

**判语＝成立（带条件）。** 凭据全部本腿现跑。

### 3.1 现在真读不读——代码位置逐枚（现读行号，⛔ 不抄四程给的号）

| 落点 | 现读坐标 | 本腿判定 |
|---|---|---|
| 取值口 | `internal/ball/hotkey_windows.go:588 resolveCancelBorrow(bind string)` | `bind==""`→出厂裸 Esc；能解析→**原样吃配置**；解析失败→裸 Esc＋`Fallback` 具名（P6） |
| 交给 Win32 | `:609 takeEscWithAcc` → `:631 takeEscBorrow` → `internal/ball/ball_windows.go:892`（`resolveCancelBorrow(b.cancelBinding)` 于 `:886`，同一枚 `borrow` 既给 `RegisterHotKey` 也给回执行） | 一处解析、两处用 ⇒ 不是"各拼一遍" |
| 回执行 | `:671 cancelBorrowedLineFor`／失败形 `:695 cancelFailedLineFor`／idle 形 `:451 cancelIdleLine` | 回执 `Binding` 字段＝配置原样 |
| 旧写死形是否还在 | `:554 escBorrowAcc()` 仍在，但**只作为默认档返回值**（`:590`／`:596`），不再是借用唯一体 | 与"写死裸 Esc"不是一回事，条件①要的正是它 |

### 3.2 两档现量（"用户看得见的那句话"到底是什么）

形状先例本腿**现跑取号**（⛔ 派单给的 `:42` 已被 `r2` 的 125 增行顶走）：
`grep -n 'Cancel: *"Ctrl+Alt+V"' internal/ball/hotkey_live_test.go` ⇒ **`:46`**（`liveHotkeys()`）与 **`:605`**。

**默认档（机主没配 `[hotkey] cancel`）**——〔自跑，全部来自本腿现跑的 `./cmd/wisp/` 整包 `-v` 日志与一发读突变〕

- 借到的加速器＝`{MOD_NOREPEAT 0x4000, VK 0x1B}`，逐位钉在 `TestBorrowDefaultIsBitIdentical260`＋`TestBorrowWishDefaultCellIsBitIdentical260r2`（本腿跑＝PASS）。
- 装配回执句（**实跑渲染**，`logs/executed-cancel-sentences.txt`）：
  `audit: resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）`
- 否决成功句（**实跑渲染**，经本腿的读-out 突变 C 从红句里读出原样字节，见 `logs/C-readout-default-sentence.txt`）：
  `按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计`
- 卡片那行（`TestTicket260R4DefaultLabelIsTheOldSentence` PASS＋D 发的红句 `got "按 Esc 键"`）：**`按 Esc 键`**。
- 通道没加载那两支（实跑渲染）：`四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）`／`本进程有悬浮球窗口，但装配根没有注入取消执行者，取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）`。
  ⇒ 这一支**不指枚任何键名**＝A595 §1 边界④要求的档一形，本腿现量到字面。

**配置档（`[hotkey] cancel` 配一枚非 Esc 的键）**——〔自跑，但见"射程"那一栏〕

- 球侧：`Ctrl+Alt+K`／`Ctrl+Alt+Space`／`F9` 三枚种子 ⇒ `RegisterHotKey` 被 handed 的加速器逐枚等于配置解析值，且 `reg.bindsEsc()` 为假；
  `Ctrl+Alt+V`（winlive 夹具那枚）⇒ `{mods 0x4003, VK 0x56}`。本腿跑＝`TestBorrowFollowsConfiguredBinding260`／`TestBorrowWishFollowsSeed260r2`／`TestLiveBorrowHelperPremiseIsConfigDependent260` **PASS**。
- 回执行 `Binding` 字段＝配置那枚串（`TestBorrowReceiptSharesOneSourceWithRegistration260` PASS，内含"打印的拼法解析回来必须等于注册出去的加速器"这半）。
- 给人看的那句：本腿能量到的是**经读口种子**那一半——实跑渲染
  `audit: resident-approval: 审批门已装配进常驻进程，取消通道 Ctrl+Alt+Q 已加载（本票只落 Ctrl+Alt+Q 一条通道；…）`，卡片那行渲染成 `按 Ctrl+Alt+Q 键`（`TestTicket260R4SeededKeyReplacesEscOnEveryFace` PASS）。
- ⛔ **产码端到端那一半＝〔未实测〕**：真开球窗＋真注册热键才有"配置的键→回执→句子"整条（winlive），机主今天没给那个词，本腿没开窗、也不许开。
  且 §2 的 M4 实证：**这一半今天就算接错也没人会响**。

**票 245 那条裁定有没有被悄悄改掉——本腿逐处查，判：没有。**

- 稳态不绑：`internal/ball/hotkey_windows.go:484-491` 的 idle 遍仍逐字 `// The one slot an idle ball must not register.` 后 `continue`；
  执行证据＝`TestConfiguredCancelStillNeverBoundWhileIdle260` PASS（内含 `reg.attemptsOf(hkCancel) != 0` 即 `t.Fatalf("ticket 245 RED: …")`，本腿跑绿）。
- 只在确认那两三秒借：唯一注册路径仍是 `TakeEscForCancel`（`internal/ball/ball_windows.go:881`，`:869` 注释仍写 "this is the ONLY path that registers the cancel id (ticket 245)"）；
  归还仍是**丢**不是重绑（`:916-931`，`releaseEscWith` 只 `unreg(hkCancel)`）。
- 改了键之后借的是他那枚＝A588 条件①写的"用户自己改了键，借的就是他那枚，这不算悄悄改 245"；默认档逐位不变由上面那两枚钉兜住。
- ⇒ ⛔ 本腿**没有**在票 260 四程里找到任何一处把"只在确认那两三秒借裸 Esc"读成"只在确认那两三秒借配置值"之外的越界；
  本票动的是**来源**，不是**时机与存在性**。`TestBorrowUnparsableFallsBackLoudly260` 里那句"配坏了仍要借默认键"也是对时机的守（不借＝改判 245），不是新形。

---

## 四、AC#2（零症状那一格要有声）

**判语＝不成立（本格没人做，且判据换成反形今天仍全绿）。**

### 4.1 票面要的判据 vs 盘上有的东西

票面 AC#2 要求：种**「借用请求发出但 Win32 没成键」** ⇒ **指名那一步必须红**。
本腿自己种了这一发（突变 A）：`internal/ball/ball_windows.go:892-899` 里"借被拒 → 把失败形写回报纸"那一步整块摘掉，
只留 `return`（`slog.Error` 仍在 `takeEscWithAcc` 里，即"打日志"这一半没摘）。〔自跑〕`logs/A-silent-borrow-failure-ball.txt`：

```
[A-silent-borrow-failure-ball] rc=1
   --- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)
       tokens_table_test.go:1468: read design/assets/tokens.css: … cannot find the path specified
```

⇒ **整包唯一的红＝§一 那枚环境红（别人的 `design/**` 删除造成的既有基线红，与票 260 无关）**。
**没有任何一枚用例因为"产码不再把没成键写进回执"而红。** 也就是说：AC#2 要的那把"会响的尺"，
在**"借用发出→回执→Problems()"这条接线**上今天不存在。

### 4.2 "抵一枚既有尺"抵到了什么、盖不盖得到这一形

- `internal/ball` 里"有声"这半**确实有尺**，而且**有牙**——本腿跑对照突变 B（`hotkey_windows.go:400-402` 的 `Problems()` `default:` 支改成 `continue`，即"attempted-but-refused 不再进 Problems"）：〔自跑〕`logs/B-control-problems-goes-silent.txt` 红句逐字
  `hotkey_cancel_borrow_260_test.go:278: the refused borrow does not name the key it tried: ""`
  `hotkey_status_test.go:280: the failed borrow produced no user-visible line: ""`
  ⇒ 两枚红＝`TestBorrowFailureNamesTheKeyItTried260`（260-r1 新增）＋`TestCancelBorrowFailureIsAProblemLine`（票 245 既有，本腿核：引入 commit `1068eb9e`，早于本票）。
  **这两枚的射程＝回执/Problems 的语义层**，不是 `TakeEscForCancel` 那条接线层。
- ⚠ 两枚尺的形状都是**测试自己手工把失败形塞回报纸**（`…withCancel(cancelFailedLine(syscall.EINVAL))`），
  没有任何一枚经由 `Ball.TakeEscForCancel` 走一遍。本腿尺：`grep -rn "TakeEscForCancel" --include=*.go` ⇒
  非 winlive 的调用者＝**零**（其余全在 `windows && winlive` 文件里：`hotkey_cancel_borrow_live_260_test.go`／`hotkey_live_test.go`／`interaction_live_test.go`／`live_windows_test.go`，
  tag 现读见 `logs/build-tags.txt`）＋产码 `cmd/wisp:671`／`cmd/balldebug:635`。
- 票面另一句"要与票 258 的 rebind×借还自钉**分开**"：本腿核——`resident_hotkey_258_windows_test.go` 那族没被并进本票任何一句判据，
  四程的具名用例（`TestBorrow…260`／`TestBorrowWish…260r2`／`TestTicket260R3/R4…`）没有一枚用"热键接好了"这种合句结案。⇒ 这半守住了。
- **四程对 AC#2 的表态（本腿逐件核）**：`r1` 件 §5 表里逐字 "**AC#2 … 本腿没做**，任务书写明不许顺手做"；
  `r2`＝winlive 尺同步（AC#1 邻格）；`r3`＝文案；`r4`＝卡片标签。⇒ **没有任何一程声称做过 AC#2，也没有"用既有尺抵"这句被写成判据。**
  本腿因此**不是**判它"追认了瞎尺"，而是判它**未落地**——但这正是票面 AC#2 那一格本身，**不能翻勾**。

### 4.3 判语

AC#2 ＝**不成立**（未做）。归口：要么给 `Ball` 的借到／归还路径开一枚**非 winlive 可注入的注册 seam**（把 `takeEscBorrow` 的 `registerFn` 变成字段），
要么正式把它记为"只有 winlive 能验"并让机主那个词来结。⛔ 本腿不裁形、不改码，只具名上报。

---

## 五、AC#3（越界检查：名册本腿自己拉）

**判语＝成立。**

### 5.1 名册尺与逐程名册

尺＝`git show --name-only --format= <每枚>`，四程**逐枚全量**（本腿比派单多核了三枚：`d9bb6817`＝r1 落地、`eb2c0173`＝r1 证据件骨架、`60767578`＝票面 AC#4），
聚合去重档＝`logs/all-paths-sorted.txt`；逐枚原文＝`logs/ac3-rosters-all-commits.txt`。

**probes／issues 之外，四程全部被触文件＝15 枚**（逐枚本腿拉出）：

| 程 | probes 之外被触文件 | 枚数 |
|---|---|---|
| `r1`（`eb2c0173`／`d9bb6817`／`b1e59d63`／`ba2a8358`） | `internal/ball/ball_windows.go`、`internal/ball/hotkey_windows.go`、`internal/ball/hotkey_cancel_borrow_260_test.go`、`internal/ball/hotkey_cancel_borrow_live_260_test.go` | 4 |
| `r2`（`84dbda52`／`e3e19e8e`／`971a5f15`／`57a33804`） | `internal/ball/hotkey_live_test.go`、`internal/ball/hotkey_cancel_borrow_expect_260r2_test.go` | 2 |
| `r3`（`85331955`／`2ae018be`／`170e0459`／`0de321fa`／`1e9ebaeb`／`da874122`） | `cmd/wisp/resident_approval_windows.go`、`cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go`、`cmd/wisp/resident_task_source_live_246_windows_test.go`、`internal/agent/approval/approval.go`、`internal/agent/approval/cancel_key_wording_260r3_test.go` | 5 |
| `r4`（`b2e4cf00`／`79c579e2`／`f8259f5a`／`eb496ea6`〔编排者代提〕） | `cmd/wisp/resident_approval_windows.go`、`cmd/wisp/resident_cancel_key_label_260r4_windows_test.go`、`internal/agent/approval/approval.go`、`internal/agent/approval/cancel_key_label_260r4_test.go`、`internal/agent/approval/cancel_key_wording_260r3_test.go`、`internal/agent/approval/gate.go`、`internal/agent/approval/report.go` | 7 |

**与编排者现跑那份对账**：`r4`＝**7 枚，逐名一致**（你给的 `cmd/wisp/resident_approval_windows.go`、`cmd/wisp/resident_cancel_key_label_260r4_windows_test.go`、
`internal/agent/approval/{approval.go, cancel_key_label_260r4_test.go, cancel_key_wording_260r3_test.go, gate.go, report.go}`）。
`eb496ea6` 只多碰 `approval.go`（已在集内）。⇒ **对上了，本腿没有多出或缺少的名。**

### 5.2 逐名对授权集

| 文件 | 授权出处 | 本腿判 |
|---|---|---|
| `internal/ball/hotkey_windows.go`／`ball_windows.go` | A588 选形裁定"写面＝`internal/ball/**`" | 在集内 |
| `internal/ball/hotkey_cancel_borrow_260_test.go`／`…_live_260_test.go` | 同上（`internal/ball/**`，测试面） | 在集内 |
| `internal/ball/hotkey_live_test.go`／`…_expect_260r2_test.go` | A594／A597 §2"裁同步不拆"（`internal/ball/**`） | 在集内 |
| `cmd/wisp/resident_approval_windows.go` | A595 §1 边界②之外的写面（`cmd/wisp` 回执／文案侧，票面"要建什么"写明） | 在集内 |
| `internal/agent/approval/approval.go` | A595 §1（`:90`／`:112` 两枚）＋A598 §2（`channelNames` 改读取时求值） | 在集内 |
| `…/cancel_key_wording_260r3_test.go` | r3 自件；r4 改它是**加严**（见 5.4） | 在集内 |
| `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go`／`…/cancel_key_label_260r4_test.go`／`cmd/wisp/resident_cancel_key_label_260r4_windows_test.go` | 各程自件判据 | 在集内 |
| `cmd/wisp/resident_task_source_live_246_windows_test.go` | 票面 AC#4 撞钉节逐字"改了文案**必须**同批改这枚常量"（`60767578`） | 在集内（**唯一被授权改的别票测试件**） |
| `internal/agent/approval/gate.go`／`report.go` | A598 §2 具名解冻"`gate.go:314`／`gate.go:477`／`report.go:142` 三处调用点" | 在集内 |

### 5.3 禁区逐名扫描（本腿自跑尺）

对 `logs/all-paths-sorted.txt` 跑一条 grep：`docs/PLAN.md`／`docs/specs/`／`internal/observe/thresholds.go`／`golden`／`tools/d22scan/allowlist.txt`／`^frontend/`／`^design/`／
`tokens_fourway_test`／`l2_grant_boundary_test`／`ticket90_persist_test`／`scripts/slo-check.ps1` ⇒ **命中行数＝0**（读数在终端记录，档同目录）。
另：`git status --porcelain -- <禁区>` 里 `design/**` 那 16 枚删除**不是四程做的**（属编排者点明的"不归你"那批）。

- **D43 转移表**：`docs/PLAN.md` 未被任何一程碰 ⇒ 一字未动。
- **不许新增状态**：`internal/ball` 的 `HotkeyStatus` 枚举（现读 `hotkey_windows.go:244-262`：Live／Disabled／Unparsable／Taken／Error／Standby）
  本腿按 `git diff 0693a077 HEAD -- internal/ball/hotkey_windows.go` 拉过：唯一形变是 `b.Status = HotkeyError` → `line.Status = HotkeyError`（赋值对象换了，不是新枚）。⇒ **没新增状态。**

### 5.4 两处本腿特别核过的"看起来像越界/像放宽"

1. **`r4` 改了 `r3` 的测试件**（`internal/agent/approval/cancel_key_wording_260r3_test.go`，39 增／16 删）。
   本腿逐行读 diff：它把 `TestTicket260R3LoadedCancelLabelIsTheNamedResidual` 从"只钉今天那句"改成
   "**先钉装了读口时那句必须不再是 `按 Esc 键`，再钉默认档仍逐字 `按 Esc 键`**"——
   **加严**，不是放宽；默认档那半枚断言一字未松（连 `t.Errorf` 的话都还是"若有意的换形请先改账再动它"）。⇒ 判：合法且更硬。
2. **`r4` 在 `internal/agent/approval` 里留了第二枚出厂默认字面量** `approval.go:126 defaultCancelKeySpelling = "Esc"`。
   A598 §2 的边界是"⛔ 不许在 `gate.go` 里再拼一遍组合键"——本腿核：`gate.go`/`report.go` 现在只调 `channelLabel(v.Channel)`／`channelLabel(r.Channel)`（现读 `gate.go:314`、`:477`、`report.go:147`），
   全仓**没有任何一处解析或拼组合键**；那枚 `"Esc"` 是**默认键名**的回退值，不是拼法。
   本腿又现跑突变 D 证它不是无人看管的重复值（**3 枚红**，见 §2.2 表末行）。⇒ 判：不构成"第二处再拼一遍"，但**它是一枚与 `ball.DefaultHotkeys().Cancel` 重复的事实**，
   两者的一致由两枚各自独立的钉（`ball` 默认→cmd/wisp 前提检查；approval 回退→r4 ①格）分别钉成 `"Esc"`，**没有一枚跨包直读比较**。这一处本腿只具名不判罪。
3. **票面框**：`grep -c '^- \[ \] '`＝**4**、`'^- \[x\]'`＝**1**，票面最后一条 commit 仍是 `60767578`（编排者那枚 AC#4），
   `git log --since="2026-10-04 13:15" -- <票面>`＝**空** ⇒ **本腿与四程都没碰过票面框**（`r1`–`r4` 各自件里的"票面未碰"本腿复认）。

### 5.5 依赖边（A598 §2 那句"⛔ 不许新增 `internal/agent/approval → internal/ball`"）

本腿现跑：`grep -rln "wisp/internal/ball" internal/agent/approval/` ⇒ **零命中**（只有注释里出现 `internal/ball` 字样，逐条核为散文）。
拼法入口＝`approval.SetCancelKeySpelling(CancelKeySpelling)`（`approval.go:161`），由装配根在 `cmd/wisp/resident_approval_windows.go:259` 注入 `ra.cancelKeySpelling`。⇒ **边为零，注入形状成立。**

---

## 六、AC#4（条件②文案侧，⛔ 不动逻辑）

**判语＝成立（带条件：条件＝§2 那半枚"读回执无人兜"必须留在账上，⛔ 不许写成端到端已验证）。**

### 6.1 ① 默认档零漂移——那十枚句子照旧逐字念 `Esc`

尺（本腿自跑，A595 §1 同形）＝`grep -nE '"[^"]*Esc[^"]*"' <两枚文件>`，取数时刻 `13:2x`（HEAD `515ca5c5`）：

- `cmd/wisp/resident_approval_windows.go` ⇒ **零命中**（八枚句子里的字面 `Esc` 全没了）。
- `internal/agent/approval/approval.go` ⇒ **两枚命中**：`:120` 是注释、`:126` 是 `const defaultCancelKeySpelling = "Esc"`。
  ⇒ 这两枚**不是**"给人看的那句话"，一枚是散文、一枚是出厂默认值本身。

⚠ **"字面里没有 Esc"不等于"渲染出来没有 Esc"**——默认档渲染必须实跑取。本腿两路现量：

- 判据路：`TestTicket260R3DefaultWordingIsTheOldSentence`（PASS）＋`TestTicket260R4DefaultLabelIsTheOldSentence`（PASS）＋
  `TestTicket260R3LoadedCancelLabelIsTheNamedResidual`（PASS）——三枚钉的都是**改前逐字旧句常量**。
- 渲染路（本腿自己把原样字节读出来）：见 §3.2 默认档那五句，全部实跑渲染含 `Esc`／`按 Esc 键`。
  其中否决成功句是本腿用一发读-out 突变（改 `r3` 件里那枚旧句常量为一枚哨兵，让红句把 `got` 打出来）取到的，红句逐字：
  `默认档否决成功句漂移：got "按 Esc 否决了卡片 corr-260r3（L1 / fs.write）：该调用未执行，答案已入审计", want "260-v1 READ-OUT SENTINEL"`
  ⇒ 突变只动测试文件里一枚常量、当场还原（`logs/AC2-readout-console.txt` 末三行 `matches_start=True`）。
- 三档分形（A595 边界④）本腿逐枚现读并对到行号（⚠ 行号＝本腿现读，A595 派单时的 `:271/:279/:287/:309/:311/:495/:584/:586` 已被 `r3`/`r4` 的增行顶位）：
  档一（通道没加载⇒不指枚键名）＝`bindBallHost` 无球支「取消键无处可借」＋无执行者支「取消键通道保持未加载」＋`approval.go` 的「快捷键取消不可用」；
  档二（已加载⇒念当前那枚）＝现读 `:367-369`（装配回执）、`:382`／`:386`（`vetoRejectedLine`／`vetoDoneLine`，`key` 形参进 `%s`）、`:596-597`（状态句）；
  档三（借被拒⇒念尝试过那枚）＝现读 `:712`（`slog.Warn`）＋`:715`（`wisp: 卡片 %s 的取消键未借到（桌面已有占位者），按 %s 不会否决它`）。
  ⇒ 十枚各有归属，**没有一枚被"统一替换"了事**；`r3` 件里那句"三档分形"本腿**复算为真**。

### 6.2 ② 配置档必须念配置里那枚键、且不再出现 `Esc`

- 卡片四枚面孔（可用性行／L1 否决句／L2 否决句／D31 回执报告行）：`TestTicket260R4SeededKeyReplacesEscOnEveryFace` PASS；
  本腿另用 CTRL 发（删注入线）实证这格**会红**，红句见 §2.2。⇒ 这半**有牙**。
- `cmd/wisp` 那几枚：`TestTicket260R3SeededKeyReplacesEsc` PASS＋本腿 CTRL2 发实证 `vetoDoneLine` 回写死 `"Esc"` 即红。⇒ 有牙。
- ⛔ **产码端到端（配置 → 球回执 → `residentCancelKeySpelling` → 句子）＝〔未实测〕**：winlive 才跑得动，且本腿 M4／M4b 实证这一半**今天没有可执行的负控**。
  所以 ②的凭据形状是"读口种子＋球侧独立钉"，**不许写成"改配置后卡片与句子真念那枚键已被实测"**。

### 6.3 `r3` 那句"第十枚停手上报"是真是假——本腿判：真话，且今天那一枚念对了

- 尺＝`git show 2ae018be -- internal/agent/approval/approval.go` 的 diff 体：`r3` 只改了 `:112` 那一枚（逐字 `-		return "Esc 取消不可用"` → `+		return "快捷键取消不可用"`），
  **没有碰 `channelNames` 那张表** ⇒ "第十枚（卡片通道名）我停手"是**动作层为真**，不是嘴上说说。
- 今天那一枚念什么（本腿现读＋现跑）：`approval.go:106-111` 表里 `ChannelEsc: "取消键通道"`＝**槽位名**；
  渲染分两支——已加载⇒`channelLabel`→`fmt.Sprintf("按 %s 键", cancelKeySpelling())`，默认档实跑＝**`按 Esc 键`**（D 发红句 `got "按 F13 键"` 反证它真的从这里取）；
  未加载⇒`channelStatusLabel` 落表里那枚槽位名，实跑＝**`取消键通道：快捷键取消不可用`**。
  ⇒ **读起来对不对**：对。A598 §2 点名的那枚当天新造的自相矛盾（「按 Esc 键：快捷键取消不可用」）在本腿现量的未加载渲染里**不存在**；
  `TestTicket260R4NoBallHostNamesNoKeyOnTheCard`／`R4UnloadLineNamesNoKey` 还把前半句里出现 `Esc`／种子键／`按` 都钉成红（PASS）。
- "裸"那个修饰词：`vetoDoneLine`／状态句今天都不含"裸"字；A598 §2 已批准"默认档比旧句少两个字"并把它定义为**修饰词不计入零漂移**。
  ⇒ 本腿按这条口径判零漂移（键名没漂），⛔ 没看成放宽。（`logs/executed-cancel-sentences.txt` 里那句 `本票只落 Esc 一条通道` 也说明"票 246 那批旧句"仍逐字在。）

### 6.4 逻辑没被动（AC#4 的 ⛔ 边）

- 派发／加载／借还：`r3`＋`r4` 对 `internal/ball/**` **零命中**（名册 §5.1）；`cmd/wisp` 侧 `vetoRejectedLine`／`vetoDoneLine` 是从 `vetoByEsc` 里**提出来的纯句子函数**，
  调用参数与形状不变（现读 `:405-413`），`bindBallHost` 的 `SetLoaded(approval.ChannelEsc, true)` 位置不变（`:361`）。
- `TestTicket260R4SpellingSeamCarriesNoAuthority` 本腿跑＝PASS（D 发红过 ⇒ 它不是摆设），钉的是"装了读口也不许把未加载通道报成可用／不许可否决"。

---

## 七、门禁读数

### 7.1 `sh scripts/d22scan.sh`

〔自跑〕**rc=0**，终行逐字含 `d22scan: clean - no D22 ban violations`；八枚 scope 现值：

```
bans #1-5 internal/=228  bans #1-5 cmd/=38  ban #6 frontend/=85  ban #7 internal/tools/=23
ban #8 design/=39  ban #8 frontend/=85  ban #8 internal/=503  ban #8 cmd/=100
```

自测段终行＝`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`。
（⚠ 本腿按派单提醒只走脚本自带的正控路径，**没有**以为它跑了整包；也**没**从根模块 `go run ./tools/d22scan`——A597 已记那枚形状坑。）

**红名集合逐名比先例**：A598 §1 记的是 `ban #8 internal/=502`、`cmd/=99` ⇒ 本腿读到 **503／100**，各 **+1**。
归因（本腿现跑尺，⛔ 不是顺手动 `tools/d22scan`）＝`git show --name-status --format= 79c579e2 | grep ^A` ⇒
**`A cmd/wisp/resident_cancel_key_label_260r4_windows_test.go`**（`cmd/` +1）＋**`A internal/agent/approval/cancel_key_label_260r4_test.go`**（`internal/` +1）。
两枚都是 `r4` 自己新增的判据件 ⇒ **与 A597 §1 那两枚 +1 同源同理**，非违规、非仪器被动。

### 7.2 `sh scripts/check-path-length-budget.sh --with-self-test`

〔自跑〕**rc=0**，逐字关键行（`logs/gate-path-length-v1.txt`）：

```
control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:
denominator: tracked paths=5536  over-budget=57  covered by roster=57  not in roster=0
longest=180 chars relative (.scratch/wisp/issues/252-the-two-spellings-…-r2-l2.md)
worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])
bands: over the hat=57  of which in the 122..180 middle=57  past the old debt line(180)=0  in the wall interval=0  roster entries=57
VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
```

`tracked paths=5536` 比 A598 时点的 5492 大 ⇒ 差值含本腿自己新增的台件（本表＋两台脚本＋约 20 枚读数档），**⛔ 本腿不为此改任何名册**；
`in the wall interval=0` ⇒ 没进那堵墙，与 `A596` 的单位口径自洽（本腿没有拿 224 去比 217）。

### 7.3 定向 `go test` 终值（三包，带 sherpa PATH）

| 包 | rc | PASS（顶层） | FAIL（顶层） | SKIP | 备注 |
|---|---|---|---|---|---|
| `./cmd/wisp/` | **0** | **218**（含子用例 314） | **0** | **0** | `ok … 187.750s`；本腿的数与 `r3` 件的"PASS=39"**口径不同**（本腿逐枚数 `--- PASS` 行，见 §9 J-4），load-bearing 的"0 FAIL／0 SKIP"一致 |
| `./internal/agent/approval/` | **0** | **58**（含子用例 77） | **0** | **1** | 唯一 SKIP＝`TestDefaultDeadlineWallClockMeasurement`（`ticket84_no_owner_test.go:222`），引入 commit `1068eb9e`（09-21），且 `git show --name-only` 逐枚核：**票 260 四程零命中该文件** ⇒ 起手即在，本腿**没把它读成通过** |
| `./internal/ball/` | **1** | **65** | **1** | **0** | 唯一 FAIL＝`TestC21TableColourRowsMatchTokensCSS`（读不到 `design/assets/tokens.css`）＝§一 的环境红，先于本票且与票 260 无关；本腿在突变 A／B 之外**未**跑绿过它，也**没**动它的断言 |

四程新增／同步的 260 具名用例逐枚现读状态：`internal/ball` 10 枚（§3.2 列名）**全 PASS**；
`cmd/wisp` 8 枚（§2.2 列名）**全 PASS**；`internal/agent/approval` 的 r3／r4 具名用例含上面 D 发红过的那三枚在内**全 PASS**（rc=0）。

`t.Skip` 现值尺（⛔ 谁都不许用 SKIP 抵读数）：`internal/ball/hotkey_live_test.go`＝**4**（`r2` 自述"改前改后都是 4"⇒ 本腿复量一致）；
其余七枚 260 相关测试件＝**0**（`logs/tskip-counts.txt`）。

### 7.4 sherpa PATH 那枚坑——本腿两发都现量过

- ⛔ 不带 PATH：`go test -count=1 -run 'TestTicket260CancelKeyWording|TestTicket260CancelKeyLabel' ./cmd/wisp/` ⇒
  逐字 `exit status 0xc0000135` ＋ `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.032s` ＋ **没有任何 `--- FAIL` 行**（`logs/test-cmdwisp-WITHOUT-sherpa-path.txt`）⇒ **用例压根没跑，这不是写手报错**。
  （⚠ 本腿这发的 `-run` 名字是自己顺手写的，跑出来本来就 0 枚匹配——但 `0xc0000135` 出现在**加载期**，与匹配数无关，形状仍复认；见 §9 J-5。）
- ✅ 带 PATH（`export PATH="$(pwd)/third_party/sherpa-onnx:$PATH"`，即 `wisp-cli-tests.sh:95` 那个"$dll_dir 原形、⛔ 不用 `pwd -W`"口径）：同一枚定向 ⇒ **rc=0，8 枚 PASS**（`logs/test-cmdwisp-260r3r4-targeted.txt`）。
- 本腿所有 `cmd/wisp` 读数（含 §2 的 M4／M4b／CTRL／CTRL2）都在**带 PATH** 下取，且每发都核过 `0.0xxs` 不是绿、也核过 `ok … 18x–19x s` 才是真跑过。

### 7.5 `gofumpt`／`go vet`（只这三包范围）

- `gofumpt` 不在 PATH（第一发 rc=127＝**工具缺失不是判据**，⚠ 与 `r3` 件的 `msg-gofumpt-note.txt` 同形）；改用 `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp internal/ball internal/agent/approval` ⇒
  点名 **3 枚**：`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/approval/queue.go`。
  本腿逐名对账（`logs/gofumpt-attribution.txt`）：三枚**都不在票 260 名册里**（`grep -qxF` 对 `all-paths-sorted.txt` 逐枚＝not in roster），
  最后触碰分别是 `5e8748b3`（10-03，"212-r1＋258-r1 编排者代笔收尾"）与 `d03d166f`（10-04 10:37，`212-r3`）。
  ⇒ **归票 212／258，不归本票；⛔ 本一枚都不许"顺手"修**。本腿另把 `r3`／`r4` 全部被触文件逐枚过一遍 gofumpt ⇒ **零点名**。
- `go vet ./cmd/wisp/ ./internal/agent/approval/ ./internal/ball/` ⇒ **rc=0**（空输出）。
- `go vet -tags winlive ./cmd/wisp/ ./internal/ball/` ⇒ **rc=0**（证明 winlive 那批编得动；⛔ 本腿**没有跑**任何 winlive，见 §8）。

---

## 八、够格翻哪几枚框（判断归本腿，翻勾归编排者）

| 框 | 本腿判语 | 建议 | 条件（若要翻） |
|---|---|---|---|
| **AC#0** | 不属本腿（已勾，编排者的） | — | — |
| **AC#1** | **成立（带条件）** | 够格翻 | 随翻勾必须一起落一句：「配置值→给人看那句」的**读回执那一半无 executable 负控**（M4／M4b 本腿实证全绿），该缺口要么另立一格、要么等 winlive；⛔ 不许让 AC#1 的勾把这半读成"已实测" |
| **AC#2** | **不成立（未做）** | **必须停勾** | 本腿种的那发"借被拒不再写回执"整包无人响（§4.1）。四程里没有一程声称做过它，⛔ 也不许用 `TestCancelBorrowFailureIsAProblemLine`／`TestBorrowFailureNamesTheKeyItTried260` 这两枚**语义层**尺去抵"接线层"这一格 |
| **AC#3** | **成立** | 够格翻 | 无条件。本腿自己拉的名册、禁区逐名零命中、D43 零命中、枚举零新增、依赖边零、票面框四程＋本腿零触碰；`r4` 改 `r3` 测试件那处经逐行读 diff 判为**加严** |
| **AC#4** | **成立（带条件）** | 够格翻 | ①默认档零漂移＝本腿实跑取到原样字节（§6.1）；②配置档念那枚键＝只在**读口种子**与**球侧独立钉**两层成立，端到端那一层〔未实测〕。随翻勾必须写明："文案侧已验；`配置→句子`的产码端到端未验，且 §2 那发突变证明它今天不会响" |

**winlive 归属**（机主那个词没给 ⇒ 本腿一律没跑，⛔ 也没用读码冒充实测）：
`internal/ball/hotkey_live_test.go`（含同步后的 `requireEscBorrowed`＋`:46`／`:605` 那枚 `Ctrl+Alt+V`）、
`internal/ball/hotkey_cancel_borrow_live_260_test.go`、`cmd/wisp/resident_task_source_live_246_windows_test.go` 的 `vetoDoneClaim` 词面钉、
以及"真改配置后卡片／句子／借到的键三处一致"这条端到端——**全部〔未实测〕**。
本腿能给的只有：winlive 那批**编得动**（`go vet -tags winlive` rc=0）＋默认档那枚词面常量与新句**逐字一致**（读码层，非实测；`r2` 件也是这么标的，本腿未复算为"已验"）。

---

## 九、记我 与 判不动／量不到

### 9.1 记我（本腿自己写错／跑歪的尺，逐条留档不当笔误藏）

| 号 | 本腿错在哪 | 影响与更正 |
|---|---|---|
| **J-1** | **第一版突变驱动（`mutate_m4_v1.py` 的姊妹脚本 `mutate_m4b_and_controls_v1.py`）把产码留在了突变态**：`finally` 里我写的是 `src = load()`——**重读了已被突变的当前文件**再写回去，于是"还原"还原的是突变体；三发之后 `cmd/wisp/resident_approval_windows.go` 停在 CTRL2（`vetoDoneLine` 写死 `"Esc"`）的形，md5 `cd00b28e…` ≠ 起手 `7a26c7a9…` | caught by my own check：驱动末行自己打了 `final md5 … matches_start = False`。本腿立刻用**任何突变之前**抽好的 HEAD 副本（`logs/head-copy-resident_approval_windows.go`，md5 恰为起手值）还原，并跑 `md5sum` ＋ `git status --porcelain -- internal cmd tools scripts .github docs`（空）复认。**读数没有被骗**：三发都从同一份 `main()` 起始读出的原始字节派生，互不叠加。**代价**：本仓的产码在这一段窗口里处于"别人看到的不是盘上的东西"的状态——这一条完全可避免（第一版 `mutate_m4_v1.py` 就是把原始字节存在闭包里、还原正确）。第二版驱动（`mutate_ac2_and_readout_v1.py`）改成"开局把全部原始字节存字典、finally 从字典还原、末了逐枚 md5 对起手值"，三枚全部 `matches_start=True`。⛔ 本腿未提交任何一次中间态 |
| **J-2** | 把 `gofumpt` 的输出读错了文件名：`gofumpt.exe` 在 Windows 上打的是**反斜杠**路径（`cmd\wisp\models.go`），我第一轮把它认成了 `cmd/wisp/pending_read.go` 并去 `git log` 那枚不存在的路径（返回空） | 更正尺＝`cat -A` 看原始字节 ＋ `tr '\\' '/'` 转真实名，逐枚 `-f` 存在性核对：真身是 `cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/approval/queue.go`。**如果我当时照第一轮的读法写进判语，就会给票 212 之外凭空记一枚不存在文件**——这条与 `A596`/`A598` 那两枚"读错单位的数"同族，只是这次是我 |
| **J-3** | 一开始把 `internal/ball` 那枚 `TestC21TableColourRowsMatchTokensCSS` 的红当成"要去归因的疑点"，没先看工作树 | 更正尺＝`ls design/assets/tokens.css` ⇒ 不存在；`git status` ⇒ `design/**` 16 枚 `D`（编排者已写明不归我）。⇒ 它是**环境红**，本腿既没修它也没把它记成票 260 的账；但在 §4.1 那发的 rc=1 里它一度看起来像"突变 A 被抓到了"——**如果本腿没有逐枚读红名集合，就会把 AC#2 判成成立**，这是这轮最险的一次误读 |
| **J-4** | 拿本腿的 PASS 计数去和 `r3` 件的"PASS=39"对比时先以为对不上＝谁错了 | 更正＝两者**口径不同**（本腿逐枚数 `--- PASS`／`^--- PASS` 两种粒度都打了；件里的数出自 `portable-tests.sh` 那套"顶层锚定行"计数）。本腿在 §7.3 把两种口径都写出来，并只把 **0 FAIL／0 SKIP** 当承重读数。⛔ 没有据此判任何一程读数不实 |
| **J-5** | §7.4 那发"不带 PATH"的复现里，我顺手写的 `-run 'TestTicket260CancelKeyWording|TestTicket260CancelKeyLabel'` **两个名字都不存在**（真名是 `TestTicket260R3…`／`TestTicket260R4…`） | 本腿核对过：这不推翻那一发的结论（`0xc0000135` 是**加载期**，早于任何用例匹配），但**这枚尺的形状不干净**，故同时跑了带 PATH 的同命令（0.063s／8 PASS）做对照并写明。若我把那发当成"用例红名的证据"用，就已经犯了本腿正在判别人的错 |
| **J-6** | 派单里两处坐标我照抄了一次，随后自己现跑更正：① `internal/ball/hotkey_live_test.go:42` 的 `Cancel: "Ctrl+Alt+V"` 现读是 **`:46`**（`r2` 加了 125 行顶位）；② `A595`/`r3` 给的 `:584`／`:586` 现读是 **`:712`／`:715`** | 本腿在 §3.2／§6.1 只写现跑号，⛔ 没把任何一枚行号当常量传下去 |
| **J-7** | 本腿一度准备把 `defaultCancelKeySpelling` 直接判成"第二处拼键名" | 更正＝先跑尺再判：本腿跑了突变 D（**3 枚红**）才下结论"不是无人看管的重复值"；同时在 §5.4 保留"它与 `ball.DefaultHotkeys().Cancel` 是重复事实、无跨包直读比较"这一句只具名不判罪。**没有**用"我觉得这是形状违规"去退回一枚已授权的落地 |

### 9.2 判不动／量不到（逐条具名＋归口）

| 号 | 本格 | 为什么动不了 | 归口 |
|---|---|---|---|
| U-1 | 配置档端到端那句话（真窗里改了 `[hotkey] cancel` 之后卡片／回执／借到的键是否三处一致） | winlive，机主那个词没给；本腿⛔不开窗、不注册热键 | 与 `258-v2` 同批等那一个词 |
| U-2 | "读回执那一步接错了会不会响"（M4／M4b 那一格） | 本腿**已量到"不会响"**，但"该用什么形补"要么改产码形状（给 `cmd/wisp` 一枚吃回执值的纯函数）要么上 winlive——两者都不在本腿权限内 | 编排者裁：另立一格或并回 AC#2 的接缝缺口 |
| U-3 | AC#2 的补形（`Ball` 借还路径的可注入 `registerFn` seam） | 属产码改动，⛔ 本腿不动；且它落在 `internal/ball/**`，会与他人写面冲突 | 编排者派腿；本腿只交 §4 的红句当靶 |
| U-4 | `internal/agent/approval.defaultCancelKeySpelling` 与 `ball.DefaultHotkeys().Cancel` 的一致性能不能一枚跨包直读钉住 | 那需要 `approval` 导入 `ball`（＝A598 §2 明令禁止的新边）或第三枚装配根尺；本腿没有实现任何一种，因此**不许据"做不到"排除候选形** | 具名〔候选形未验证〕；要结这格得先有人造一枚候选形再判 |
| U-5 | `resident_task_source_live_246_windows_test.go:69` 那枚词面常量同步后真窗会不会红 | 本腿没跑 winlive；CI 也不带 `-tags winlive`（本腿现读该文件 tag＝`windows && winlive`，`logs/build-tags.txt`） | 仍是〔预测非实测〕，与 `A595` §2、`A597` §2 同一格 |
| U-6 | `gofumpt` 那三枚（`cmd/wisp/models.go`／`pending_read.go`／`queue.go`）为什么没被 CI 挡 | 不在票 260 射程，也不在本腿写面；本腿没去查 CI 是否跑 gofumpt | 归票 212／258 的下一次碰这三包；本腿只具名 |
| U-7 | 整包 `go test ./...` 会不会有别的红 | ⛔ 派单明令不跑（会撞别人的包）；本腿只跑三包 | 编排者推送窗口的 CI |

### 9.3 交件自量尺（本腿自跑，读数存档 `logs/placeholder-selfcheck.txt`）

- 尺＝`sh .scratch/wisp/probes/260/v1/placeholder-check.sh`（模式串只存在脚本里，⛔ 不写进本表，否则本表自己就命中自己）；
  它对本表跑四枚编排者点名的形状＋两枚同族标记，判据＝**命中 0 行**。
- ⚠ 本表里有两枚字样不是占位、而是**产码原句引用**，逐具名留下不删：§6.1 那句「桌面已有占位者」是 `cmd/wisp/resident_approval_windows.go:715`
  渲染文案的一部分（本腿现读的字节），另一处是本节标题里"占位"这个词本身。⇒ 四枚模式串在这两处**零命中**，尺跑过存档。
- 本腿全程未 `t.Skip`、未放宽任何断言、未动阈值／golden／`thresholds.go`／`tools/d22scan/**`／`scripts/slo-check.ps1` 一字节；未跑 `slo-check.ps1`。
