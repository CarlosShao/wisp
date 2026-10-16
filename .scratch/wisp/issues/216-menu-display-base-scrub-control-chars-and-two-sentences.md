# 216 — **菜单的显示底座**：任何一行字只要**不是用户自己打的**，显示前必须洗过控制字符；而且**"没有"和"读不到"必须是两句话**

- Status: **待派**。横跨票 213／214／215 的**每一枚字符串**，所以单独立一枚（owner 22:1x 的原话是「**聊天框这块的内容不少的**」，这一枚就是那块地基）。
- 来源：调研腿 `survey-plusmenu`（`docs/reports/survey-2026-09-28-composer-plus-menu.md`，86486 字节）第 4 题与"第 4 枚票"那一节；两家各写了一份、**写法不同**，且都**专门对文件名下手**。

## 现量（起手逐条复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| 别家有这枚底座 | Step-Code：转成可见转义（`consent.ts:73-92`）；MiniMax：删／替换（`terminal-text.ts:15-26`） | 见上面那份产出件（**行号自己现读**） |
| 我们有没有 | **〔未量〕**——起手先跑 `grep -rn --include=*.go -iE "sanitiz|control.?char|escape.*terminal|\\\\u001b" internal/ cmd/` 点数，**量不到就写"没做全仓反向取证"**（不许写"没有"） | 本票 AC#1 |
| 我们已有的相邻件 | 文件名守卫 `NameGuard`（只管名字合法性，**不管显示安全**） | `internal/panel/attachments.go:165` |

## 为什么这一枚值得单独立（不是洁癖）

菜单里每一行都可能印出**别人控制的名字**：附件文件名、插件名、技能名、目录名。一个带控制字符的名字，在菜单里就是**给用户看的注入面**——
它不改权限、不越界，但它能让用户**看见一行不是真的话**（例如把后半段名字盖掉、或伪造出"已允许"的样子）。
这与我们已经防住的"内容不改权限"（票 200 那枚 hostile-file 判据）是**同一族的两半**：**一半管它说什么不算数，一半管它印出来长什么样**。

## 判据

- [x] **AC#1**：先做**现量普查**（我们到底有没有这枚底座、有几处各写一份），出逐枚名册；**没名册不许进实现**（这条是本仓加尺的固定前置，见票 212 同形）。
- [ ] **AC#2**：**一枚纯函数**（洗控制字符／可选可见转义），**菜单里所有非用户手打的字段都过它**：附件名、插件名、技能名、目录名、命令说明、错误原因。
  ⚠ **落点选"显示前统一洗"，不选"每个生产者各自洗"**——后者会长出 N 份互不相同的实现（别家已经是两份）。
- [ ] **AC#3（正控，必须有）**：造一枚文件名里带控制字符（如 `ESC`＋"已允许"）的附件／插件 ⇒ 菜单里那一行**必须原样不可信**（洗过），且**不许因为这个名字里出现了"允许"两个字而改变任何判定**。
- [ ] **AC#4**：**"没有"与"读不到"分两句话**：清单里每一行要么给内容，要么给"这一枚我读不到，因为……"（票 200 那五状态 `loaded/none_found/off/refused/not_run` 是同一纪律的先行版）。
  ⚠ 不许用空数组／缺键同时表达这两件事（那是 `A412` 点名的"键在值恒缺"）。
- [ ] **AC#5**：错误提示**按错误码分支**，**不许按英文错误串里有没有某个数字**（openchamber 那形 `ChatInput.tsx:2205` 靠 `413` 子串决定弹哪句＝反面教材，我 22:3x 只读到腿的转述、**未逐字复核**，落地前自己再看原文）。

## 禁区

不改权限判定链路（本票纯显示层）；不新增 C17 方法名；`frontend/**`、`design/**` 零写（**界面侧那份由 owner 带给界面那支**，本票只管 Go 侧给出去的字已经洗过）；
`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 一字不动；只 commit 不 push、显式 pathspec、禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；仓内不删东西。
— 2026-10-09 12:0x 编排者腿 boxes-1 搬框：本节原为散句 bullet，现改成 - [ ] 形以进全勾尺；⛔ 判据文字一字未改、⛔ 未翻任何框。
- **AC#1 现量普查已交（只读腿 `216-a1`，`2026-10-16`，尺读提交树快照 `71002d0`）**：名册见 `docs/evidence/s1/216-census.md`。三句头条——① **无"显示前统一洗"底座，且自家已长出两份互不相同的控制字符处理**（`sanitizeReplyText` 换空格 vs `findMetachars` 只转 `\n\r\0` 三字符可见转义），`unicode.IsPrint`／ANSI `\x1b`／`\p{C}` 全仓 0；② 今天真出向又裸奔的菜单字段＝**附件名 `AttachmentRef.Name`**（上游拒绝 `ContainsAny(name,"\x00\r\n\t")` 只挡 4 字符、**ESC/C1 穿过不洗**）与**目录名 `WorkspaceView.Canonical`**；插件/技能/命令说明三枚零源零出向（票 213/215 待派）。③ **"没有 vs 读不到"在显示段已是两句话**（先例五态是 `internal/panel/instructions_200.go` 真代码，非仅票面），但**派单把第五枚状态写成 `na`、盘上真值是 `not_run`**——实现腿按 `not_run` 走；产码里**按错误子串分支 0 枚**（`Contains/Index/HasPrefix/regexp/switch on Error()` 实测 0，票面 `ChatInput.tsx:2205`/`413` 属外部 openchamber、我方无此文件）。⛔ 本腿未做 AC#2/3/4/5、未翻任何框、未改任何产码。
- [2026-10-16 10:4x +08] 编排者翻 **AC#1**（这一格是"先出名册"的普查格，凭据＝只读腿 `216-a1` 名册件 `docs/evidence/s1/216-census.md` 现量 **134 行／14673 字节** ＋我逐条复跑；尺一律读**提交树快照 `71002d0`**，⛔ 混在飞工作树）：四条复现＝① `unicode.IsControl` 整树（`internal`＋`cmd`＋`tools`＋`frontend`）命中 **1 枚**＝`cmd/wisp/approval_reply.go:617`（在 `sanitizeReplyText` 内，2 枚生产调用者）；② 第二份各写各的＝`findMetachars`（`internal/risk/rules_shell.go:57`，只把 `\n\r\0` 转成可见转义）；③ 今天真出向又裸奔的那枚＝附件显示名：上游 `checkDisplayName` 的拒绝只挡 `"\x00\r\n\t"` **四枚字符**（`internal/panel/attachments.go:301`），而 `sanitizeDisplayName`（同件 `:310`）只做 trim＋取末段 ⇒ **ESC／C1 类控制字符今天能穿过进菜单**；④ 插件名／技能名／命令说明三枚今天**零源零出向**（`internal/panel/composer.go` 里 `plugin|skill|command` 命中 0 枚，属票 213／215 的射程）。★**它顶回我派单的一处错、我复跑同意**：我把票 200 那五状态的第五枚写成 `na`，盘上真值＝**`not_run`**（尺＝`git grep -nE 'InstructionsStatus[A-Za-z]+ =' 71002d0 -- internal/panel/instructions_200.go` ⇒ `loaded`／`none_found`／`off`／`refused`／`not_run` 五枚逐枚对上；全仓 `\bna\b` 只命中 `ticket84_no_owner_test.go` 里的测试局部变量）⇒ 名册按盘上真值算、我派单那句⛔ 抹平、只在这里更正。另注口径：腿的正控写作"`.Error()` 114"，我这两把尺读出的是 **374 行／144 枚文件**（`git grep -n` 数行 vs `-c` 数文件）⇒ 只证"这条尺不是恒 0"，⛔ 与它的数相减。next＝**`216-r1` 落地腿（待派，⛔ 现在放）**：面收成 `internal/panel/**` 一枚纯函数＋菜单字段过它＋正控；⛔ 动 `cmd/wisp`、⛔ 新增加出的名册键＝会动 313 键钉与 JSON 键钉 ⇒ 要加先停手上报。**释放条件具名＝本波唯一写腿 `220-r2` 交件之后**（`cmd/wisp/**` 与 `internal/panel/**` 是同一块共享编译面、每波只许一枚写腿，这条是我自己立的规矩，⛔ 我自己破）。
- [2026-10-16 12:2x +08] 落地腿 **`216-r1`（写）交件**：两笔 `b7676daf`（产码 `internal/panel/composer.go`／`display_scrub_216.go`／`pump.go`／`workspace.go`）＋ `92e8cc57`（测试件 `internal/panel/display_scrub_216_test.go` 一枚新件），读数件 `docs/evidence/s1/216-display-scrub-r1.md`（176 行／15799 字节）。**AC#2**＝一枚包内纯函数 `scrubForDisplay`（C0＋C1 逐枚转**可见转义**、非法 UTF-8 字节转原始 `\xNN`；⛔ 删——删会把两半拼成一句新的话；⛔ `strconv.Quote`——它把每枚反斜杠翻倍、正常 Windows 路径读废；幂等），唯一调用点落在显示边界 `NewSnapshot`（票面"显示前统一洗"那格，⛔ 各生产者各自洗；`Canonical` 在别处是当**路径**读的），名册＝AC#1 实测今天真出向又裸奔的两枚（`AttachmentRef.Name`、`WorkspaceView.Canonical`），洗的是副本。**AC#3**＝真 `memory.Open`＋真 `memory.ValidArtifactName`＋真 sniff 把 `"photo\x1b已允许.png"` 从入口走到快照，牙证：去牙（纯函数首行改 `return s`）⇒ `TestScrubForDisplay…`／`TestTheSnapshotBoundary…`／`TestAHostileAttachment…` 三枚红（红句逐字在件里 §2），还原后同三枚 `ok 0.132s`、还原哈希 `6ea6d3ca04d7eb07a0ae9f2f8cfc6f4d96b346c0` ＝ 提交树同名 blob；判定不动那半用真 `risk.NewRiskAssessor()` 三发同屏 `plain=L2/[R1 R8] hostile=L2/[R1 R8] washed=L2/[R1 R8]`。**AC#4**＝泵在 `Workspace` 读数腿缺失时不再说出"未选择工作区"那句，改说具名的"这一维没有读数"（`workspaceNotChosenReason`／`workspaceNoReadingReason` 常量对＋`unreadWorkspaceView()`），空 `ComposerState` 的工作区行也必给一句；第五枚状态按盘上真值 **`not_run`** 走（⛔ 派单那个 `na`），并新钉一枚测试把五态字面与 `none_found` vs `not_run` 两段关系钉住。**成对台面读数**：同一枚探针、两枚仓外导出树（⛔ 共享工作树跑改前那发），改前 `ca0e30fc` 树红（`U+001B at offset 5`＝穿过）⇄ 改后 `92e8cc57` 树绿（洗掉）；`./internal/panel` 整包带 **`-timeout 30m`**：改前 141 PASS／8 FAIL ⇄ 改后 148 PASS／8 FAIL，**红名集合逐枚相减＝空 ⇒ 本腿引入新红 0 枚**（8 枚既有红逐枚归因在件 §5：3 枚导出树缺件 `docs/`／`design/`／`.git`、2 枚本机内存仪器、3 枚 by-design 键钉）。**门禁四数各自 rc**＝gofumpt 0／`scripts/d22scan.sh` 0／长度门 0／`go vet ./internal/panel/` 0。**本腿动过的既有测试件＝零枚**（⛔ 改松任何一枚既有断言、⛔ 删用例、⛔ `t.Skip`）。`cmd/wisp/**` 与 `frontend/**`、`design/**` 零字节；⛔ 新增加入 Snapshot 名册的键、⛔ 新 C17 方法名、⛔ 新导出 `panel.*` 名。**本腿未翻任何框、未改票名／`Status:`、未加 `-done`、⛔ push。** 残余具名上交（件 §7 逐枚）：①"两份合一"（`sanitizeReplyText`／`findMetachars`）需一枚新导出名＝契约面＋两侧形状本就不同 ⇒ 归编排者裁定另派；② Cf 类（`U+202E` RTL OVERRIDE 正是票面"盖掉后半段名字"那形、`U+200B`、`U+FEFF`）今天**仍穿过**，扩不扩是裁定；③ `"attachments"` 段还是"一句表两件事"（泵至今无附件 reader，空数组同说"没附件"与"没读数腿"），拆它要**新键** ⇒ 按规矩**停手上报**未加，同处 `attachmentError` 键在值恒缺（零写入者）；④ `ApprovalCardView.Tool/Args[]/Reason` 仍逐字出向（AC#1 的"相邻面"，不在本票名册）；⑤ `cmd/wisp/panel_pump.go:126` 的 `rt.workspaceView()` 在 `rt.paths == nil` 时同样说"未选择工作区"＝AC#4 坍缩在 cmd 侧的另一半，不是本腿写面。**AC#5 本腿零改动**：按错误子串分支那形在本腿两枚新件 0 命中、在 `git grep 92e8cc57 -- internal cmd tools`（⛔ `_test`）亦 0 命中（形状缺席，非零尺未测；正控 `strings\.Contains` 同树 903 命中），此格留给非实现者验收腿裁。
