# 216-r1 读数件 — 票 216 AC#2／AC#3／AC#4 落地腿（写码腿自证）

> 腿＝`216-r1`（落地写码腿），⛔ 裁决。本件只报**我这一发量到了什么**，不替验收腿下判语。
> 起手 HEAD 现跑＝`83173370b48dcc0438b647b6161cb64116e2a042`；放腿时刻编队已推进到
> `ca0e30fce32425e31b9882c5d947e56a3b17d851`（＝我第一笔的父，下面所有"改前锚"都用它，
> ⛔ 抄派单里的任何号）。
> 我的两笔：`b7676daf`（产码四枚）→（其间别腿落了 `ddffbf09`／`5326f65e`／`e13456e4`）→ `92e8cc57`（测试一枚）。
> ⚠ 我自己犯过一次锚错并当场纠：第一版"改前树"我取的是 `HEAD~2`＝`5326f65e`，`diff` 文件清单
> 发现它**已含我的 `display_scrub_216.go`**（别腿在我两笔之间落了三次 commit，`HEAD~2` 不是我的父）。
> 已用 `git rev-parse b7676daf^` ＝`ca0e30fc` 重建树重跑，上面那两发读数出自重建后的树。

## 0 写面自证（逐枚 `git show --name-only` 可复现）

- `b7676daf` 只含：`internal/panel/composer.go`、`internal/panel/display_scrub_216.go`（新）、
  `internal/panel/pump.go`、`internal/panel/workspace.go`。
- `92e8cc57` 只含：`internal/panel/display_scrub_216_test.go`（新）。
- 本件＋`.scratch/wisp/probes/216/r1/**`＋票 216 末尾**一条** bullet＝本腿全部写面。
- `cmd/wisp/**` 零字节；`frontend/**`、`design/**` 零字节；⛔ 新增加入 Snapshot 名册的键；
  ⛔ 新 C17 方法名；⛔ 新导出 `panel.*` 名（本腿所有新符号一律包内小写）。
- ★**必答项：本腿动过的既有测试件＝零枚。**两笔里没有任何一枚 `_test.go` 被修改，只有新增一枚
  测试件；既有断言一枚没改、一枚没删、⛔ `t.Skip`。（尺：上面两条 `git show --name-only`，
  `grep -c '_test.go'` 在 b7676daf 为 0。）

## 1 AC#2 — 一枚纯函数＋显示前统一洗

**形状**（`internal/panel/display_scrub_216.go`）：`scrubForDisplay(string) string`，
`unicode.IsControl`（C0 `U+0000-U+001F` ＋ C1 `U+007F-U+009F`）逐枚转**可见转义**，非法 UTF-8
字节转原始字节 `\xNN`；`\n \r \t \0` 用助记形，其余 `\xNN`（≥`U+0100` 留 `\uXXXX` 支路）。

- **为什么转义而不是删**：删会把两半拼成一句新的话——`photo\x1b已允许.png` 删掉控制字符后读起来
  像 `photo已允许.png`，那句假的"已允许"仍然在，而唯一能证明它被注入过的字节没了。转义保留每一个
  其它字节的位置，把不可见的那枚印出来。
- **为什么不用 `strconv.Quote`**：它把每枚反斜杠翻倍，菜单里每一枚正常 Windows 路径会变成
  `C:\\Users\\...`，即为了罕见的敌意枚把常见的诚实枚读废。
- **为什么落在显示边界而不是各生产者**（票面原话）：自家已有两份互不相同的实现
  （`cmd/wisp/approval_reply.go` 的 `sanitizeReplyText` 换空格 vs `internal/risk/rules_shell.go`
  的 `findMetachars` 三字符可见转义），再开"各生产者各自洗"就是第四份。且 `WorkspaceView.Canonical`
  在别处是**当路径读的**（`ProjectInstructionLoadRequestFor`、`ReadGitForWorkspace`），在生产者侧洗
  会改掉功能面。唯一调用点＝`NewSnapshot`（`composer.go`，"the one push the panel receives per state
  change"），泵与任何直接构造都过它；洗的是**副本**，不改调用方切片（`ForAgent` 读同一批行）。
- **名册（只装 AC#1 实测今天真出向又裸奔的两枚）**：`AttachmentRef.Name`、`WorkspaceView.Canonical`。
  插件名／技能名／命令说明＝票 213／215 射程，今天零源零出向 ⇒ ⛔ 为它们造键、⛔ 填内容，
  它们将来进 Snapshot 时由同一枚调用覆盖。

**成对读数（同台面，同一枚探针，两枚仓外导出树）**：尺＝
`go test ./internal/panel/ -run '^TestProbe216HostileDisplayNameAtTheDisplayBoundary$' -v -timeout 10m`，
`GOTMPDIR='C:\Users\swq\AppData\Local\Temp\wisp216r1'`（先 mkdir 再用反斜杠形）。
建树七件套 `go.mod go.sum internal cmd tools scripts frontend`；两侧 `go build ./internal/panel/` 均
`rc=0`；`.go` 枚数＝改前 743 ／ 改后 744（差 1 枚＝我的产码件）。探针件：
`.scratch/wisp/probes/216/r1/pair/probe_216_pair_test.go`（不指名任何 216 新符号，故两树同码可编译）。

- **改前（`ca0e30fc` 树）＝红，rc=1，"穿过"那一句逐字**：
  `probe_216_pair_test.go:79: PROBE go-row="photo\x1b已允许.png" raw-ESC-in-go-row=true wire-row="photo\x1b已允许.png" raw-ESC-in-wire-row=true`
  `probe_216_pair_test.go:84: the menu row reached the page with a live control character U+001B at offset 5: "photo\x1b已允许.png" (this is the AC#2 穿过 reading; the display boundary washed nothing in this tree)`
- **改后（`92e8cc57` 树）＝绿，rc=0**：
  `PROBE go-row="photo\\x1b已允许.png" raw-ESC-in-go-row=false wire-row="photo\\x1b已允许.png" raw-ESC-in-wire-row=false` ⇒ `--- PASS`
- ⚠ **口径具名**：JSON 自己会把 `U+001B` 转义成 `\u001b`，所以"快照字节里没有裸 ESC"这一枚**改前也成立**
  （我在仓内件里保留了那枚断言，但它不是牙）。**真正的尺是解码后的值**（页面画的就是它）。

## 2 AC#3 — 正控，且这颗牙自己证明会响

**正控（真装配）**：`internal/panel/display_scrub_216_test.go` 的
`TestAHostileAttachmentNameWalksTheRealBrokerIntoAWashedMenuRow` ＝ 真 `memory.Open(t.TempDir())`
（真 SQLite＋真 `winsec` 封目录）＋真 `memory.ValidArtifactName` 守卫＋真 PNG 字节＋真 sniff，
名字＝`"photo\x1b已允许.png"`，经 `NewComposerState`＋`NewSnapshot` 出快照。
另断言入口那侧**仍带原始 ESC**（`ref.Name` 含 `U+001B`）——证明洗落在显示边界、没落到生产者。
盘上读数：`AC#3 menu row as the page receives it: "photo\\x1b已允许.png"`。

**去牙（把纯函数那行改成原样返回）**：在 `scrubForDisplay` 首行插 `return s`，同台面跑五枚指名用例：

- 红 3 枚（逐字，取自 `display_scrub_216_test.go`）：
  - `:72: scrubForDisplay("photo\x1b已允许.png") = "photo\x1b已允许.png", want "photo\\x1b已允许.png"`
    与 `:75: output "photo\x1b已允许.png" still carries a control byte (0x1b) at offset 5`
    （同一件里 NUL／CR LF／TAB／DEL／`a\u0085b`／`a\u0090b`／`a\u009fb`／`\x1b\x1b` 各子用例逐枚红）
  - `:107: the attachment row reached the menu as "photo\x1b已允许.png", want the ESC shown as a visible escape`
    与 `:110: the canonical path reached the menu as "C:\\data\\\u0085reports", want the C1 shown as a visible escape`
  - `:157: menu row = "photo\x1b已允许.png", want the ESC shown as backslash-x-1-b and every other byte kept`
- 仍绿 2 枚（具名，⛔ 把它们当牙）：`TestTheWashReturnsTheSameStringWhenThereIsNothingToWash`
  （干净路径原样返回，去牙天然成立）与 `TestTheClaimInANameChangesNoVerdict`（它钉的是"判定不动"，
  对洗的存在**必须不敏感**——这是它对去牙绿的设计原因，不是漏牙）。
- **还原**：`cp` 回写后 `git hash-object internal/panel/display_scrub_216.go` ＝
  `6ea6d3ca04d7eb07a0ae9f2f8cfc6f4d96b346c0`，与提交树 `git rev-parse HEAD:internal/panel/display_scrub_216.go`
  **同一枚**；`git diff --stat` 该件为空；`grep -c MUTATION` ＝ 0。
  还原后同三枚再跑＝`ok github.com/CarlosShao/wisp/internal/panel 0.132s`。

**"名字里出现'允许'两个字不改变任何判定"（权限判定那一支的同发读数）**：
`TestTheClaimInANameChangesNoVerdict` 用真 `risk.NewRiskAssessor()`（工具桥跑的同一枚）经
`NewApprovalCardView` 三发同屏——原始名／含"已允许"名／洗过的名：
盘上读数 `real assessor verdicts: plain=L2/[R1 R8] hostile=L2/[R1 R8] washed=L2/[R1 R8]`
⇒ Level 与 RulesHit 三发相同，且 ⛔ 因文件名降级（`L0` 那条断言在）。
另一半：同字节两枚独立店、两枚 broker，除显示名外逐项相等（`Stored/MIME/Kind/Artifact/SizeBytes/
Deduplicated/Reason`），`Artifact` 仍是内容哈希名。

## 3 AC#4 — "没有"与"读不到"拆成两句话

改前坍缩形（实测，非推断）：`pump.go` 在 `PumpSources.Workspace == nil`（**没人读过这一维**）时
直接渲染 `UnsetWorkspaceView()`，即说出"未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定"
——那是"没有"那句话。本腿把它拆成两枚具名句（`workspace.go` 的常量对
`workspaceNotChosenReason` / `workspaceNoReadingReason`＋`unreadWorkspaceView()`），
泵无读数腿时走后者；`Set` 仍 false、`Canonical` 仍空（⛔ 反过来发明一个范围）。
`NewSnapshot` 另补一枚：空 `ComposerState` 的工作区行若无内容且无原因，填"读不到"那句
（与同件里 `Mode.Current=="" -> unknown`、`Git.Kind=="" -> GitViewNotProbed()` 同一纪律）。

盘上字节（`TestWorkspaceAbsentAndWorkspaceUnreadAreTwoSentences` 的 `t.Logf`，同发两枚）：

- 没有＝`未选择工作区：本轮按 [fs] allowed_dirs 授权的目录判定`
- 读不到＝`工作区这一维没有读数：本机的快照泵没有被接入工作区读数腿（PumpSources.Workspace=nil），
  所以「未选择工作区」这句话我没有资格说——这里是没有读数，不是没有工作区（票 216 AC#4；
  票 200 的 not_run 是同一纪律）`

第五枚状态按**盘上真值 `not_run`** 走（⛔ 派单那个 `na`）：`TestTheFiveInstructionStatesAreStillTheDiskSpelling`
把 `loaded/none_found/off/refused/not_run` 五枚字面钉住，并把 `none_found`（查过·没有）与 `not_run`
（没读数）比成两段：`InstructionsSectionFromBundle(&projctx.Bundle{})` vs `(nil)`。
⛔ 新增任何 Snapshot 名册的键（会撞 313 那枚键钉），全部复用既有 `reason` 键。

## 4 AC#5 — 本腿不做，且"零改动"这事实的读数

尺（同形，两把）：`grep -rnE 'Contains\([^)]*Error\(\)|strings\.Index\([^)]*[Ee]rr|HasPrefix\([^)]*Error\(\)|
HasSuffix\([^)]*Error\(\)|MatchString\([^)]*\.Error\(\)|switch [^(]*\.Error\(\)'`
- 射程＝本腿两枚新件 ⇒ **0 命中**；射程＝`git grep 92e8cc57 -- internal cmd tools`（⛔ `_test`）⇒ **0 命中**
  （与 AC#1 普查同读数＝形状在产码里**缺席**，不是零尺未测）。
- 正控（证这把尺不是恒 0）：同 commit 树 `git grep -nE 'strings\.Contains' -- internal` ＝ **903 命中**。
- 本腿对 AC#5 的改动＝**零枚**（未新增、未删除、未改写任何按错误文案分支的形状）。

## 5 `./internal/panel` 整包四数（显式带超时）

尺＝`go test ./internal/panel/ -v -timeout 30m`（带的就是 **`-timeout 30m`**；默认 10m 对别枚重包不够，
本包 11.2 s 未触及），`GOTMPDIR` 反斜杠形，两枚仓外导出树同法同台面。

| 发 | 树 | rc | 顶层 PASS | 顶层 FAIL | 子测试 PASS | SKIP |
|---|---|---|---|---|---|---|
| 改前 | `ca0e30fc` 导出树 | 1 | 141 | 8 | 76 | 0 |
| 改后 | `92e8cc57` 导出树（＝HEAD） | 1 | 148 | 8 | 90 | 0 |

名集合逐枚相减＝**空**（`diff fails-before.txt fails-after.txt` rc=0）⇒ 本腿引入新红 **0 枚**；
顶层 +7 枚＝本腿新增的 7 枚用例，子测试 +14 枚＝那枚表驱动的 14 个字符类子用例。

那 8 枚既有红逐枚归因（同发两树同名同句；⛔ 计入本票判据、⛔ 为它们动断言）：
1. `TestApprovalCardViewJSONKeysMatchFrontendTypes` — `approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`（键钉，Q-51／A383，页侧未声明）
2. `TestComposerContractTypesMatchFrontend` — `composer_test.go:78: 同一枚`（同上）
3. `TestComposerRenderFixtureTellsTheTruth` — `composer_test.go:873: the render fixture has no block for "no host attached" ...`
4. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` — `backpressure_bound_35r8_test.go:309: ... live heap +1201...`（本机内存仪器）
5. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8` — `:336: ... 1x opened 4000 streams ...`（本机内存仪器）
6. `TestGitDimensionHasNoModelCallableTool` — `git_test.go:369: read docs/PLAN.md: open ...\tree-final\docs\PLAN.md: The system cannot find the path specified.`＝**台面缺件**（七件套不含 `docs/`）
7. `TestC21DesignTokensFourWayAgree` — `tokens_fourway_test.go:441: read design/assets/tokens.css: ...tree-final\design\assets\tokens.css ...`＝**台面缺件**（不含 `design/`；机主工作树里该件还是删除态）
8. `TestReadGitOnThisRepositoryIsSelfConsistent` — `git_test.go:600: kind = "not_a_repo" ... : this repository has a .git directory`＝**台面缺件**（导出树无 `.git`）

⇒ 其中 3 枚是导出树缺件、2 枚是本机内存仪器、3 枚是 by-design 的键钉／文案钉。⛔ 本腿把它们读成绿。

## 6 门禁四数（逐枚落 rc）

| 门 | 尺 | rc | 读数 |
|---|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l` 逐枚列本腿 5 个文件（⛔ 用 PATH 里那枚，PATH 形本机会 rc=127） | **0** | 零输出（改前后各跑一次，含测试件） |
| d22scan | `sh scripts/d22scan.sh` | **0** | `runtests.sh: OK ... PASS=35 FAIL=0 SKIP=0` ＋ `examined 281 production Go files`、ban #8 `internal/` 575 件（注释与 `_test.go` 都在射程）＋ `clean - no D22 ban violations` |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `positive control PASSED`；`denominator: tracked paths=14199 over-budget=60 covered=60 not-in-roster=0`；`VERDICT GREEN` |
| go vet | `go vet ./internal/panel/`（射程具名＝本包；`go vet ./...` 全模块归 CI） | **0** | 零输出（去牙那一发中途我另有 `unreachable code` 属变异体，⛔ 计入门禁读数） |

## 7 残余（具名上交，⛔ 本腿自行动手）

1. **"两份合一"没做**：`sanitizeReplyText`（`cmd/wisp/approval_reply.go:609`）与 `findMetachars`
   （`internal/risk/rules_shell.go:57`）仍在各自位置。把本件的纯函数交给 `cmd/wisp` 复用**需要一枚新导出的
   `panel.*` 名**，那是契约面（本腿禁区），且两侧输出形状本就不同（换空格 vs 可见转义），合一＝裁定不是编辑。
2. **Cf 类没洗**：`unicode.IsControl` 只管 Cc。`U+202E`（RTL OVERRIDE，正是票面"把后半段名字盖掉"那形）、
   `U+200B`、`U+FEFF` 今天**穿过**本件。扩不扩是裁定，⛔ 我和没被授权的腿自己 widening。
3. **`"attachments"` 那一段还是"一句表两件事"**：`ComposerState.Attachments` 空数组同时表达
   "本轮没附件"与"这台机器的泵没有附件读数腿"（`PumpSources` 至今没有附件 reader，泵恒传 nil）。
   拆它需要**一枚新键**（`attachmentsKnown` 形），会撞 313 的键钉与页侧 `panel.ts`（Q-51）⇒ 本腿按规矩
   **停手上报**，未加。同处的 `attachmentError`（`ComposerState.AttachmentReason`）今天**零写入者**
   （尺：`grep -rn "AttachmentReason" --include=*.go internal cmd tools` 只命中定义那一枚）＝键在、值恒缺。
4. **`ApprovalCardView.Tool/Args[]/Reason` 仍逐字出向**（AC#1 普查的"相邻面"），不在本票 composer 名册里，
   本腿⛔ 顺手扩。
5. **`cmd/wisp/panel_pump.go:126`（`rt.workspaceView()`）**在 `rt.paths == nil` 时也回 `UnsetWorkspaceView()`
   ＝ AC#4 那个坍缩在 cmd 侧的另一半；那面不是本腿写面，具名留给编排者。
6. 本腿**未翻任何 AC 框**、未改票名／`Status:`、未加 `-done`；票 216 末尾只追加一条 bullet。
   ⛔ push（只有编排者推）。
