# 票 216 AC#1 —— 菜单显示底座的控制字符洗普查（只读·名册）

> **本件射程**：只做票 216 的 **AC#1**——四问的逐枚名册＋尺＋射程，让下一枚实现腿照着进。
> **本件之外都没做**：不造纯函数（AC#2）、不做正控（AC#3）、不翻任何框、不改任何产码；
> `frontend/**`、`internal/**` 零字节。本腿不产任何 `.go`。

## 0. 锚点、来源口径、零改动自证

### 0.1 来源＝提交树快照（⛔ 混在飞工作树）
并发写腿 `220-r2` 正在动 `cmd/wisp`（`panel_pump.go`＋`run.go`＋一枚新测试件），且工作树本身脏（`go.mod`、`design/**` 删除、`-` 等一堆在飞件）。
⇒ 本件**所有读数出自 pin 死的提交树**：进场 `git rev-parse --short HEAD` = **`71002d0`**（`2026-10-16 10:16:04 +0800`，分支 `dev`）。
每一把尺都写成 `git grep ... 71002d0 -- <path>`，**读提交对象、无视工作树在飞态**；编排者日后原样复跑即得同一读数（不受别的腿再 commit 影响）。
⚠ 例外具名：`sanitizeReplyText` 落在 `cmd/wisp/approval_reply.go`（属 `cmd/wisp` 面），但我读的是 `71002d0` 里那一份、非工作树，故稳定。

### 0.2 尺（只用 grep/git grep/sed/ls/wc）
- 取数一律 `N=$(... || true)` 避开 `grep -c` 命中 0 时 rc=1 截断 `&&`；报 0 前均跑过**已知存在的正控**（见各问"正控"）；⛔ 用 `| head` 退码当尺的退码。
- 一把尺曾因 ERE 字符类里的方括号报**假 0**（`len\(.*\) == 0` 我第一版写 `[A-Za-z0-9_.\[\]]` 命中 0），已换 `len\(.*\) == 0` 复跑＝417 命中，正控 `== 0` 在别处亦命中——**本文所有"0 命中"都是配了正控后的真 0**。

### 0.3 零改动自证
本腿写面仅三处：`docs/evidence/s1/216-census.md`（本件）、`.scratch/wisp/probes/216/a1/`（含 `logs/` 原始尺转写）、票 216 末尾追加**一条** bullet。产码/契约轴/阈值/golden/`frontend/**`/`internal/**`/`design/**` 全部只被 `git show`／`git grep`／`sed -n` 读，未被写。

---

## 问 1：有没有"洗控制字符"这枚底座？有几处各写一份？

**结论**：**没有"显示前统一洗"的共享底座**；有 **2 处各写一份的、语义不同的控制字符处理**＋1 处**只挡 4 个字符的拒绝**＋3 处**名叫 sanitize 却不洗控制字符**的改名件。这正是票面警告的"别家已经两份"之病——**我们自家已经长出了两份**。

### 1.A 真洗控制字符的（两处，各写一份）
| 枚 | 文件∶行（`71002d0`） | 洗哪一类 | 洗法 | 生产调用者枚数 |
|---|---|---|---|---|
| `sanitizeReplyText` | `cmd/wisp/approval_reply.go:609`（IsControl 判在 `:617`） | **全 C0/C1**（`unicode.IsControl`） | 控制字符→空格，折叠空白，截断 | 2（同文件 `:571 reject`、`:575 panelReject`） |
| `findMetachars` | `internal/risk/rules_shell.go:57`（转义在 `:78-88`） | **仅 `\n\r\0` 三枚**（shell 元字符集内） | 转**可见转义**（`\n`→`\n`、`\r`→`\r`、`0`→`\0`） | 1（同文件 `:34` R6 风险理由） |

> 两份**语义不同**（一份"换成空格"、一份"转可见转义"）、**覆盖面不同**（全控制字符 vs 三字符子集）、**落点不同**（审批理由 vs 风险理由），且都**生产者侧各写各的**——⛔ 覆盖菜单里任何一行。`sanitizeReplyText` 注释自陈"票 216's control-character wash，落在产文那侧"，即已知这枚不解决本票。
> 尺（逐枚复跑）：
> `git grep -nE 'unicode\.IsControl' 71002d0 -- internal cmd tools frontend`
> `git grep -nE 'func (sanitizeReplyText|findMetachars)' 71002d0 -- cmd internal`
> `git grep -nE '\b(sanitizeReplyText|findMetachars)\(' 71002d0 -- internal cmd tools | grep -v 'func '`（调用者，测试另计）

### 1.B 只挡 4 个字符的拒绝（非洗）
- `internal/panel/attachments.go:301`：`if strings.ContainsAny(name, "\x00\r\n\t") { return errors.New("contains control characters") }`——对附件显示名做**整条拒绝**，字符集**只有 NUL/CR/LF/TAB**。**ESC(`\x1b`)、BEL(`\x07`)、其它 C0/C1 一律放过**。
  尺：`git grep -nE 'ContainsAny\(name, "' 71002d0 -- internal/panel/attachments.go`；正控：`git grep -nE 'strings\.ContainsAny' 71002d0 -- internal cmd tools | grep -v _test`（14 命中，语法确能命中）。

### 1.C 名叫 sanitize、实则不洗控制字符
| 枚 | 文件∶行 | 实际做 | 是否碰控制字符 |
|---|---|---|---|
| `sanitizeDisplayName` | `internal/panel/attachments.go:310`（调用者 `:256`，1 枚） | `TrimSpace`＋取末段（`filepathBaseLike`） | **否** |
| `sanitizeBlobID` | `internal/secret/migrate.go:217`（调用者 `:141`，1 枚） | 非 `[A-Za-z0-9._-]`→`_`，为**文件系统名** | 间接（控制字符落 default→`_`），但**非显示层** |
| `sanitizeArgs` | `internal/tools/bridge.go:1290` | 截断到 2048 字节供 forensics 列 | **否** |
| `escapeLike` | `internal/memory/models.go:216` | SQL `LIKE` 通配 `%_\` 转义 | **否**（假朋友，排除） |

### 1.D 假朋友 / 不存在（逐枚具名）
- `scrub`（9 命中）**全是前端"拖动调值滑块"（scrub-to-change inspector）**，与控制字符无关。尺：`git grep -niE 'scrub' 71002d0 -- internal cmd tools frontend`。
- `unicode.IsPrint` **0 命中**（真 0，正控见 0.2）。`\x1b` **0**、`\u001b` **0**、`\p{C}` **0**——**全仓（含前端）无一处处理 ANSI 转义序列/通用控制类**。`\x7f` 仅 1（在敌意测试 fixture 里）。
- `strconv.Quote` 25 命中：除 2 枚工具件（`tools/d22scan/gitignore.go:362`、`tools/mockllm/sse.go:123`）外全在 `_test.go`；它是"转可见转义"机制但**只用于测试/工具输出**，无一处喂菜单显示。尺：`git grep -nE 'strconv\.Quote' 71002d0 -- internal cmd tools frontend | grep -v _test`。

---

## 问 2：菜单里"非用户手打"的字段有哪些、今天从哪个源出

**出向通道口径**：C17（`internal/panel/bridge.go`）名册 7 枚——`MethodModeRequest/WorkspaceRequest/AttachmentAdd/MessageSend`（`panel.*`，`:42-45`）＋`MethodConfigGet/ConfigSet/ModelsList`（`config.*`，`:66-67/75`）——都是**入向请求方法**。菜单字段真正的**出向通道是 `Snapshot`**（`internal/panel/composer.go:57`，"the one push the panel receives per state change"）经 pump 走同一根 C17 桥管道。所以下表"出向"列指该字段寄居的 `Snapshot` 段。

| 字段 | Go 侧源（文件∶行） | 中间有无洗 | 寄居的出向段（json 键） | 今天状态 |
|---|---|---|---|---|
| 附件名 | `AttachmentRef.Name`；由 `sanitizeDisplayName(src.DisplayName)` 产（`internal/panel/attachments.go:256/310`），上游仅 `:301` 拒绝 NUL/CR/LF/TAB | **仅 trim+末段**；ESC/C1 **既不放行处不洗** | `Snapshot.Composer.Attachments`（`"attachments"`）；另 `internal/panel/bridge.go:108` 消息信封（`omitempty`） | 有源、裸奔 |
| 目录名(用户拼写) | `WorkspaceView.Spelling`（`internal/panel/composer.go:214`，注释自陈"the user's own string, verbatim"） | **无**（明示逐字） | `Snapshot.Composer.Workspace`（`"workspace"`） | 有源·用户手打 |
| 目录名(权威解析) | `WorkspaceView.Canonical`（`:216`，C26 解析后的树） | **无** | 同上 | 有源、裸奔 |
| 插件名 | `internal/plugin/disposal.go:74 Name`（处置记录）；**不进快照**（`git grep -niE 'plugin' 71002d0 -- internal/panel/composer.go`＝0） | 无（也无出向） | — | 票 215 待派·零出向 |
| 技能名 | 无源：技能连配置键都不存在（票 215 现量 `grep -c Skill internal/config/schema.go`＝0；本腿复核见下） | — | — | 零源·零出向 |
| 命令说明 | 无源：斜杠命令目录产码 0 命中（自证见下） | — | — | 票 213 待派·零源 |

**"待洗名册"（今天真出向、且非用户手打、且无控制字符洗）**：
1. `AttachmentRef.Name`（`"attachments"`）——最危险：`:301` 只拒 4 字符，**含 ESC/C1 的附件名能穿过拒绝、原样进菜单**。
2. `WorkspaceView.Canonical`（`"workspace.canonical"`）——Go 解析的路径、非用户手打、无显示洗。
3. （`Spelling` 是用户自己打的，按票面"非用户手打"定义不入待洗名册，列出仅为区分。）
> 相邻面（超出加号菜单本票射程，具名备查）：`ApprovalCardView.Tool/Args[]/Reason`（`internal/panel/approval.go:39+`，`json:"pending"`）也逐字出向（`Reason` 注释"verbatim"、`Args` 直接 append），同样无控制字符洗。尺：`git grep -nE 'type ApprovalCardView struct' 71002d0 -- internal/panel`。

---

## 问 3："没有"与"读不到"现在是同一句话还是两句话

**结论**：在显示段这条链上**已是两句话**，且先例（票 200 五态）**在代码里、不是只在票面**。

### 3.A 五态词盘上真存在（逐枚尺）
`internal/panel/instructions_200.go`：`:27 InstructionsStatusLoaded="loaded"`、`:30 NoneFound="none_found"`、`:32 Off="off"`、`:37 Refused="refused"`、`:40 NotRun="not_run"`。
逐枚尺（每条应 ≥1 命中）：
`git grep -nE '\b(loaded|none_found|not_run)\b' 71002d0 -- internal cmd`
`git grep -nE 'InstructionsStatus[A-Za-z]+' 71002d0 -- internal/panel/instructions_200.go`
⚠ **派单把第五枚写成 `na`——与盘上不符**：`git grep -nE '\bna\b' 71002d0 -- internal cmd tools frontend` 仅 3 命中且全在 `internal/agent/approval/ticket84_no_owner_test.go`（一枚测试局部变量名）。**盘上真值第五枚＝`not_run`**（与票 216 AC#4、`composer.go:66` 注释一致）。实现腿按 `not_run` 走，勿采 `na`。

### 3.B "两句话"先例已落进代码
- `instructions_200.go:29` 注释逐字："This is NOT the same fact as `off`"——`none_found`(真查过·没有) 与 `refused`/`not_run`(读不到) 被拆成互斥状态。
- 快照 `Instructions`/`Tasks` 用**指针＋常发键**：`composer.go:24-31` 明写"bare list plus `omitempty` made four different truths marshal as one missing key"，200-r2 已把列表换成带 `Status` 的段。
- git 维四态 `internal/panel/git.go:68-71`：`repo / not_a_repo / permission_denied / unreadable`（"没有仓库" vs "读不到"分列）。
- 诚实半布尔：`ModelKnown`、`CredentialKnown`、`ReasonKnown`（`approval.go`，`:reasonKnown`）——"没人读过"与"读了是空"分开。

### 3.C 今天还有没有"读不到⇒渲染成没有"的坍缩形
显示路径**已主动反坍缩**：`Snapshot` 各显示段（Composer/Instructions/Tasks/Git）都带显式状态/原因，无靠缺键表达。
残余 `slice + omitempty` 命中 11 处（尺：`git grep -nE '\[\][A-Za-z0-9_]+ +.json:[^,]+,omitempty' 71002d0 -- internal cmd tools | grep -v _test`），逐判：
- 面板侧仅 `internal/panel/bridge.go:108 Attachments`（入向**消息**信封、非显示快照）；
- 余下 `internal/llm/*/request.go`（tools/stop/include）＝**发给 provider 的线路**、`internal/models/manifest.go`、`internal/plugin/disposal.go`——均不在加号菜单显示链上。
⇒ 问 3 无"活着的坍缩形"可点名；若实现腿要新增菜单字段，**必须沿用状态＋原因形，勿再引入 slice+omitempty**。

---

## 问 4：错误提示有没有按英文错误串里的子串决定弹哪句

**结论**：**产码里 0 枚**按错误子串分支——"子串决定弹哪句"这一形在我方 Go/前端**不存在**（不是"零尺未测"，是配了正控的实测 0）。票面点名的反面教材 `ChatInput.tsx:2205` 靠 `413` 子串属**外部件 openchamber**，我方仓无此文件。

### 4.A Go 侧尺（逐枚 0 命中，均配正控）
| 反面形 | 尺 | 命中 |
|---|---|---|
| Contains on err | `git grep -nE 'Contains\([^)]*Error\(\)' 71002d0 -- internal cmd tools \| grep -v _test` | **0**（含测试则 159，全在 `_test.go`，见 4.C） |
| Index on err | `git grep -nE 'strings\.Index\(.*([Ee]rr|Error\(\))' 71002d0 -- internal cmd tools \| grep -v _test` | **0**（`strings.Index(` 本身有 18 文件命中＝正控在） |
| HasPrefix/Suffix/Cut on Error() | `git grep -nE '(HasPrefix\|HasSuffix\|CutPrefix\|CutSuffix)\([^)]*Error\(\)' 71002d0 -- internal cmd tools \| grep -v _test` | **0** |
| regexp/switch on Error() | `git grep -nE '(MatchString\|switch ).*\.Error\(\)' 71002d0 -- internal cmd tools \| grep -v _test` | **0** |
正控：`.Error()` 在非测试 Go 里 114 处命中（尺 `git grep -nE '\.Error\(\)' 71002d0 -- internal cmd tools | grep -v _test | wc -l`）——语法确能命中，故上面 0 为真 0。

### 4.B 前端侧（只读列，同形）
- 我方无 `ChatInput.tsx`；有 `frontend/src/components/composer.tsx`、`ai-native/chat-composer.tsx`（尺：`git ls-tree -r --name-only 71002d0 -- frontend/src/components`）。
- `.includes(`/`.indexOf(` 大量存在（正控）但**对三位数字状态零命中**、对 `(message|errMsg|error).includes/indexOf/match` **零命中**：
  `git grep -nE '\.(includes|indexOf)\(.{0,20}[0-9]{3}' 71002d0 -- frontend/src`
  `git grep -nE '(message|errMsg|error)\.(includes|indexOf|match)\(' 71002d0 -- frontend/src`
- 我方前端无 `413` 字样：`git grep -nE '\b413\b' 71002d0 -- frontend`（真 0）。

### 4.C "能不能被一次错误文案改动打红"
- **产码分支点：0 枚**，故**不存在可被一次错误文案改动打红的产码分支**（形状缺席，非未测）。
- 唯一会被文案改动打红的是 159 处 `_test.go` 里的 `strings.Contains(err.Error(), <字面>)`——它们是**故意钉住错误文案的测试断言**（措辞一改即红＝预期行为），不是"决定弹哪句"的分支逻辑。尺：`git grep -nE 'Contains\([^)]*Error\(\)' 71002d0 -- internal cmd tools | wc -l`（=159）。
- 正路存在且普遍：按 typed code / `errors.Is` / iota 枚举分支，56 个非测试文件命中（尺：`git grep -cE 'errors\.Is\(|code ==|\.Code\b|ErrKind' 71002d0 -- internal cmd tools | grep -v _test | wc -l`）。

---

## 给实现腿的一句话进法（不越位，只指路）
AC#1 名册已出：**底座无、且自家已有两份互不相同的控制字符处理**（`sanitizeReplyText` 换空格 / `findMetachars` 三字符可见转义），加号菜单里今天真出向又裸奔的是 `AttachmentRef.Name`（拒绝只挡 NUL/CR/LF/TAB，ESC/C1 穿过）与 `WorkspaceView.Canonical`；插件/技能/命令三枚今天零源零出向。AC#2 的"纯函数＋显示前统一洗"落点，应替换掉这两份各写各的、并补 `AttachmentRef.Name` 的 ESC/C1 缺口。⛔ 本件不下 AC#2 的判决，只把地基量清楚。
