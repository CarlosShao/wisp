# 票 260 — 面板里那把"取消"热键**改了没用**：配置值从来不会被注册，借用那把键写死裸 Esc；而丢借用**零症状**（卡片到点必执行）

**立票时刻**：2026-10-03 09:4x +0800，锚点 HEAD `7de96e61`（`dev`）
**来路**：只读普查腿 `258-a2`（`.scratch/wisp/probes/258/a2/census.md`，606 行／55,101 字节，`3736f0dd`→`5d4f5343`）的"外延 A"；**编排者本人复量三条硬断言后才立票**（逐条见下）；关联＝票 245（裸 Esc 全局热键，已裁乙）／票 258（形 A 落地腿的邻格，⛔ 本票不是 258 的残余）。
**性质**：⚠ **算产品行为**——用户会以为那枚键能用。

## 现量（编排者本机自跑，2026-10-03 09:4x；⚠ 引用前先重跑）

1. **Win32 注册全仓产码只有一枚调用点**：`grep -rn "pRegisterHotKey" --include=*.go internal cmd | grep -v _test` ⇒ 命中 `internal/ball/hotkey_windows.go:391`（调用体）＋`internal/ball/win32_windows.go:42`（proc 声明，不是调用）。
2. **它只被两遍用到，两遍都不注册配置来的 cancel**：① idle 那遍**显式跳过 cancel 槽**（`hotkey_windows.go:458` 逐字注释 "The one slot an idle ball must not register"，配合 `:13`／`:82-85` 的设计说明）——这是**票 245 裁乙**的既有裁定（稳态不绑、确认那两三秒才借），⛔ 本票**不推翻它**；② 借用那遍的加速器**写死裸 Esc**：`hotkey_windows.go:118` `vkEscape = 0x1B`、`:521` `escBorrowAcc`（注释自述 "VK_ESCAPE with nothing…"）、借还成对体在 `internal/ball/ball_windows.go:875-901`。
3. ⇒ **推论（我复量后认它成立）**：用户在 `config.toml` 把 `[hotkey] cancel` 改成任何非 Esc 的组合，**那枚键在今天任何一条路径上都不会成键**——idle 不绑（设计）＋借用不吃配置（缺陷）。⇒ 任何"四枚热键都以配置为准"的句子**今天只对三枚成立**（summon／mute／panel）。
4. **丢借用的后果不是"少个键"，是"卡片到点必执行"**：`internal/agent/approval/gate.go:319`→`:332` 逐字 `ANSWER-EXPIRED decision=timeout->execute`（`:338`），窗口 3s＝`internal/agent/approval/queue.go:116`；且 `Standby` 不进 `Problems()`（`internal/ball/hotkey_windows.go:365`）⇒ **零症状、零告警**。

## 要建什么

- [x] **AC#0（第一格，读数不是码）**：判死"借用那遍到底读不读 `[hotkey] cancel` 的配置值"——逐处给 `file:line`（起点＝现量 2 那四处），并列出**若要吃配置**需要哪些零件（`Accelerator` 从哪来、mods 怎么带、`escBorrowAcc` 那枚常量是谁的、改它会不会撞上"票 245 只在确认窗借裸 Esc"那条裁定的**边界**）。⛔ 本格不许改产码；量不到的那处要具名写"量不到"。
- [x] **AC#1 只在 AC#0 交完并由编排者落一枚具名 `A##` 选形之后才许动**：两条候选形由 AC#0 的读数决定——ⓐ 让借用读配置值（则"按 Esc 取消"这项默认行为**可能随配置变**，要与票 245 那句"确认那两三秒借的那把键"逐字对齐，⛔ 不许把 245 的裁定悄悄改掉）；ⓑ **维持写死**，但把"这枚配置项今天不生效"**说到明处**（建球/回执侧具名，⛔ 不许留"四枚都吃配置"这类句子）。两形都要给"改了配置之后，用户看得见什么"那一句。
  ⚠ **翻勾口径（2026-10-04 13:5x，凭据＝非实现者验收腿 `260-v1` 的 `.scratch/wisp/probes/260/v1/verdict.md` 462 行／50,435 字节／占位 0）＝成立·带条件**：形ⓐ 真落地（配置值→球回执那一段**有钉**，本腿 10／10 PASS），但**"球回执→念给人看的那句话"这一跳今天没有人兜**——唯一碰到 `residentCancelKeySpelling` 的用例只喂 `nil`（`resident_cancel_key_wording_260r3_windows_test.go:201/:204`），其余全经 `cancelKeyRead` 读口种子。⇒ 本格勾的是**行为与句子各自成立**，⛔ 不许读成"端到端（改了配置→界面上那句话）实测过"；那一形与 winlive 那批一起留〔未实测〕。**票 245 未被悄悄改掉**：idle 那一遍仍逐字 "The one slot an idle ball must not register"（`internal/ball/hotkey_windows.go:484-491`）＋ `TestConfiguredCancelStillNeverBoundWhileIdle260` PASS，借还仍是"丢"不是"重绑"。
- [ ] **AC#2 零症状那一格要有声**：丢借用（或借用没成键）这件事今天**不进 `Problems()`、不打日志**；判据＝装一发**会响**的尺（种"借用请求发出但 Win32 没成键"⇒ 指名那一步必须红，⛔ 不许 `t.Skip`、不许把 SKIP 读成通过），并且⚠ 要与票 258 落地腿的 rebind×借还自钉**分开**（那是"改热键时正有借用在飞"，本格是"借的那把键根本没成"，⛔ 两格不许合并成一句"热键接好了"）。
  ⛔ **验收判"不成立"（2026-10-04 13:5x，`260-v1` 本腿亲手种的这一发）**：摘掉 `internal/ball/ball_windows.go:892-899` 那处失败回执的写入 ⇒ `internal/ball` 除了**既有环境红** `TestC21TableColourRowsMatchTokensCSS`（`design/assets/tokens.css` 被人删，早于本票）之外**没有任何一枚为接线层而红**。⇒ 既有那两枚尺的射程是**回执/`Problems()` 语义层**（形状＝测试自己 `withCancel(cancelFailedLine(...))` 把失败行塞回去），**非 winlive 真调用 `TakeEscForCancel` 的用例＝零**（尺＝`grep -rn TakeEscForCancel --include=*_test.go internal cmd`，调用点全在 `//go:build windows && winlive` 那四枚文件里，编排者复跑到）。四程逐件表态亦为"没做"。⇒ 本格**不追认、不降格**，派 **`260-r5`** 补这一枚会响的尺（写面按 `internal/ball/**` 起，要动 `cmd/wisp` 先停手上报）。
- [x] **AC#3 越界检查**：`git diff` 出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` 任一路径 ⇒ 直接退回。⛔ **D43 转移表一个字不许动**；⛔ 不许新增状态。
  ✔ **成立·无条件（凭据＝`260-v1` §AC#3，本腿自拉名册）**：probes/issues 之外共 **15 枚**；`260-r4` 那 **7 枚**与我（编排者）拉的**逐名一致**（含我代提的 `eb496ea6` 只多碰 `approval.go`）；禁区逐名 grep＝**0 命中**；`HotkeyStatus` 枚举零新增；`internal/agent/approval/**` 里 grep `wisp/internal/ball`＝**0**（⇒ A598 §2 那条"不许新开 `approval→ball` 依赖边"守住）；`channelNames` 前后两态（**11 行@12:06 → 8 行@`79c579e2` 之后**，`gate.go:314`/`:477`/`report.go:142` 三处直读全改走 `channelLabel`）＝两个数都真、取数时刻不同。票面框在本腿交件前仍是 **4 未勾／1 已勾**＝四程与本腿零触碰。

## 禁区（本票全程）

- ⛔ **票 245 已裁的那条"稳态不绑、只在确认那两三秒借"属既有裁定**：本票要动它的**边界**＝人工批准（`SPEC-12 §4.1`），编排者单方面派腿＝跑歪模式。AC#1 若量出"必须碰 245 那条"，⛔ 停下来上报，不许自己改判。
- ⛔ 不动 `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件；不许为变绿放宽断言。
- ⛔ `frontend/**`／`design/**` 两层禁令；⛔ 凭据值不进对话/日志/表。
- Git：只 commit 不 push；显式 pathspec 写在 commit 命令上；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 排程与串行

写面＝`internal/ball`（＋`cmd/wisp` 的回执/文案侧）。⚠ **`cmd/wisp` 此刻 `255-r2` 在飞、`internal/panel` 此刻 `253-r5` 在飞** ⇒ AC#1 落地腿**按住**；AC#0 是只读普查，可先派（⛔ 禁跑 `go build`/`vet`/`test`，与写腿同机即互洗）。
**与票 258 的关系**：258 的落地腿管"改完配置后重注册"，本票管"cancel 这枚槽根本没吃配置"——**相邻但不互相吞并**，⛔ 谁也不许把对方的格子并进自己一句结案。撤销本票＝「260 撤」。

---

## 编排者翻勾节（2026-10-03 21:5x，锚 `d8cbb948`；260-a2 交件收档）

**AC#0 翻勾**，凭据＝`.scratch/wisp/probes/260/a2/census.md`（192 行零待填，零 Go 命令）：
- **判定＝借用遍不读 `[hotkey] cancel` 配置值。** 我抽验中心凭据：`internal/ball/hotkey_windows.go:524` 逐字 `func escBorrowAcc() Accelerator { return Accelerator{Mods: modNoRepeat, VK: vkEscape} }`＝零参常量字面；借用链四枚签名零键参。成立。
- **258 桥关系判定**：链已通到借用侧**取值字段**（`ball_windows.go:124/:126` 已喂活），借用注册路径是链上唯一没接的一段 ⇒ **AC#1 若选 ⓐ（借用读配置），258 链零改动自动获益**——两票耦合比我立票时预估的浅。
- **★推翻 a1 两处**（照单入账）：①「出厂 DefaultHotkeys 在 :171、配置从未被读」已被 258-r1 推翻（现 `:269 Hotkeys: cfg`），票 260 结论收窄为"cancel 一枚"；② a1 §5.6 将来时已成本体——idle 日志 `binding` 字段已是配置值、句内 "Esc" 写死＝改键用户今天看见自相矛盾句（**AC#1 形 ⓑ 的头号必改落点**）。
- **245 边界判定**：撞不到已生效裁定（裁的是机制与两枚落点，未裁加速器来源；默认档逐位不变）；相抵的是 245 **未勾**的 AC#6②——只报不判，归人工批准面。
- **AC#1 选形**：等 260-a2 的代价表＋形 ⓑ 头号落点（idle 日志自相矛盾句）一起裁——**不急于本轮**；AC#2（丢借用要有声）与 258 的 rebind×借还自钉是两格，票面已写死不许合并。

## 编排者选形裁定（2026-10-04 **09:2x**，账 A588；落笔时 HEAD 现取＝0693a077〔174-a4 骨架枚〕。⚠ 我这一节初稿把钟点写成 09:3x、锚写成 79b3982b，两处都偏，更正在 A588）＝**形ⓐ（借用那遍吃配置值）＋三条硬条件**

上一节那句"不急于本轮"到此结清。判据来自只读腿 260-a2 的代价表（.scratch/wisp/probes/260/a2/census.md，192 行零待填），中心凭据我自己复跑过（internal/ball/hotkey_windows.go:524 零参返回字面量；:578-586 那行回执已经带 Binding 字段）。

**为什么选 ⓐ**：260-a2 判死"258 那条链已把配置值喂到借用侧的**取值字段**（internal/ball/ball_windows.go:124/:126），借用注册是链上唯一没接的一段"⇒ 选 ⓐ 时**票 258 零改动直接获益**，而选 ⓑ 等于在一条已经通到门口的链旁边写一句"这枚配置今天不生效"。本项目的口径一直是"改了配置要真生效"（票 258 选形 A 同因）。⚠ 这句是我（编排者）的判断，不是腿的判语；它给的是两形的代价，没给推荐。

- **条件①（默认档不许漂）**：改完之后，**用户没配 [hotkey] cancel 时借到的必须与今天逐位相同**（260-a2 已量：ParseAccelerator 吃 "Esc" 得 Mods=0x4000/VK=0x1B ＝ 今天的字面量）。这条要有一枚**常驻判据**钉住，不许只写在注释里。⇒ 票 245 那条"稳态不绑、只在确认那两三秒借裸 Esc"的裁定**在默认档一字不动**；用户自己改了键，借的就是他那枚，这不算悄悄改 245。
- **条件②（说到做到）**：借到哪枚键，卡片/回执那行就得印哪枚（今天 Binding 字段已在回执里，来自 hotkey_windows.go:578-586）⇒ ⛔ 不许出现"页面写着按 Esc、实际借的是别的键"这种说了不做的形。
- **条件③（配坏了要响亮）**：260-a2 点名的新零件 P6（用户把 cancel 配成解析不了的值，今天这条分支根本不存在）必须落地成**具名退回默认 Esc 并在回执里说这句**，⛔ 不许静默吞掉解析失败。

**写面与排程**：落地腿 260-r1 只许动 internal/ball/**（＋证据件目录）。⚠ **现在不派**：写码腿已达 3 枚（212-r2 在 tools/d22scan＋internal/agent＋internal/tools，258-r2 在 cmd/wisp，262-r1 在 scripts/ 与 .github/workflows/ci.yml），按"写码 3 枚（极限 4）"的下界**按住等一枚交件再补位**。AC#2（丢借用要有声）与 AC#1 是两格，⛔ 不许落地腿顺手把 AC#2 一起做了。

---

## 编排者追加一格（2026-10-04 11:0x，账 **A595**；落地腿 **260-r3** 已在飞）

形ⓐ 落地后（260-r1，`eb2c0173`→`ba2a8358`），**借哪枚键**已经吃配置了，但**给人看的那句话还写死着"Esc"**＝A588 条件②禁的那一形。这一格是本票的**最后一格**，授权凭据＝**A595 的具名解冻五样齐**（文件／行／理由／边界／撤销口令「260 文案撤回」），⛔ 本票其余四格的条件都不因这一格改变。

- [x] **AC#4（条件②的文案侧收尾，⛔ 不动逻辑）**：改完之后——① **用户没配 `[hotkey] cancel` 时那十枚句子照旧逐字念 `Esc`（零漂移）**；② **用户改了配置时那句话必须念他配置里那枚键**，且不再出现 `Esc`。射程＝**10 枚字符串**（`internal/agent/approval/approval.go` 两枚：`:90` "按 Esc 键"、`:112` "Esc 取消不可用"；`cmd/wisp/resident_approval_windows.go` 八枚：`:271`/`:279`/`:287`/`:309`/`:311`/`:495`/`:584`/`:586`，⚠ 行号以派单时刻现跑为准）。
  ✔ **成立·带条件（凭据＝`260-v1` §AC#4，本腿自己两档现跑）**：① **默认档零漂移**成立——但那把"字面尺"的**射程要按我（编排者）现跑更正**：`cmd/wisp` 的**产码包内**今天没有 `Esc` 字符串字面量了，⚠ **不等于"`cmd/wisp` 目录零命中"**（`cmd/wisp/testdata/esclistener/main.go:162/:163/:190/:218` 那枚**台件**里还有，它是票 246 那族"第二个进程收不收得到 Esc"的实验装置，不归本票、也不是卡片文案），验收腿那句"零枚含 `Esc` 字符串字面量"**写宽了地界**，按此读；默认档的**渲染**是由本腿用 read-out 突变取出原样字节（仍念 `Esc`）。② 配置档四枚面孔种子后念 `Ctrl+Alt+Q` 且不再出现 `Esc`，并有 CTRL／CTRL2 两发实证会红。③ "第十枚停手上报"核为**动作层真话**（`2ae018be` 的 diff 只碰 `approval.go:112`，没碰那张名册）；今天那一枚＝未加载时「取消键通道：快捷键取消不可用」、已加载时「按 \<读口\> 键」⇒ `A598` 点名的那枚自相矛盾在实跑渲染里**不存在**。⚠ **条件＝端到端（改了配置→界面上那句话真的变了）今天仍〔未实测〕**：那一跳要真开窗，欠机主那一个词（见 AC#1 那条口径，同一格不许两处读成已验）。
  - ⛔ **只改"指代那枚键的名词"**：派发／加载／借还的任何一行逻辑都不动；票 245 那条"稳态不绑、只在确认那两三秒借"的**时机与存在性**一字不许动。
  - ⛔ **不许新造第二条拼键名的事实来源**：印出来的键名必须来自借用侧同一枚链（`ConfiguredHotkeys().Cancel` ＋ `internal/ball/hotkey_windows.go` 里 260-r1 已落地的那份 receipt 拼法），⛔ 不许在 `cmd/wisp` 里再拼一遍组合键字符串——那正是本票要治的形状。若 `internal/agent/approval` 拿不到 ball（它不该新增一条依赖边）⇒ **停手上报**，⛔ 不许自行决定加依赖。
  - ⚠ **三枚属"通道根本没加载"分支**（`:271`／`:279`／`:112`）诚实形是"取消键无处可借／不可用"，⛔ 不是换成另一枚键名。
  - **撞钉（A595 §2 已由编排者做完预检）**：全仓**词面**钉只有 1 枚——`cmd/wisp/resident_task_source_live_246_windows_test.go:69` 的 `const vetoDoneClaim = "按 Esc 否决了卡片"`（用在 `:131` 等子进程 stdout）。该文件 tag＝`windows && winlive`，而 CI 的 cmd/wisp 步**不带 `-tags winlive`**（尺＝`grep -nE "tags|winlive" scripts/wisp-cli-tests.sh .github/workflows/ci.yml`＝零命中，两枚文件都存在）⇒ **CI 永不执行它**；本机在无 sherpa PATH 时 `0xc0000135` 且**没有 `--- FAIL` 行**⇒ 今天也量不到。⇒ 结论：**改了文案必须同批改这枚常量**，而这一条目前只是**读码判断＝〔预测，非实测〕**，不许写成"已验证"。另：`cmd/wisp/resident_approval_246_windows_test.go:343`（tag 只有 `windows`）钉的是"没有可否决的确认项"那句，⛔ 不在射程内、不会被改到。
  - **正控两格（⛔ 不许合成一句"文案接好了"）**：默认档零漂移一枚＋种一发改配置后那句话跟着变一枚；后者要给红句逐字。反向⛔ 不许做成词面尺。
  - ⛔ **winlive 一律不跑**（机主今天尚未批准真开窗／真注册热键）；只许 `go vet -tags winlive ./cmd/wisp/` 证明编得动。
