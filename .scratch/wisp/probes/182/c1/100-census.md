# 182-c1（二轮）任务监控栏逐堆普查 — 交件事（2026-10-08，锚 `17a54338`）

> 只读腿 `182-c1` 二轮交件。⛔ 零产码、零 go 命令、`frontend/**` 未读未引（按本票 AC#7 从严）。
> 件内每把尺的原文都能原样复跑；尺的 rc 行在配套件：`10-ticket-rulers-recheck.txt`／`20-rail-tickets-boxes.txt`／`30-go-side-readings.txt`／`40-misc-and-controls.txt`／`50-demo-registry.txt`／`60-cost-minefield-controls.txt`（同目录）。
> **查重纪律**：占用条／排队序号／停止三堆 167-c2 已量（`.scratch/wisp/probes/167/c2/census.md`，锚 `0ce6cb91`）——本件**只引用其尺与读数，未重跑**；唯一发现的数值冲突具名报在 §5.3。

## 1. AC#1 栏定义与堆数（现量）

口径＝编排者 09-28 17:5x 裁定：**这一栏的堆数＝原型 `design/doubao/demo/**` 注册表口径，不是产品真身**。owner 的截图是**外部产品界面，不是本仓规格**（票面 AC#1 要求明写，此处明写）。

- 主尺（50-demo-registry.txt D1a）：`grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` ⇒ **7 行命中**（7 枚注册：files/plugins/review/terminal＋context/activity/approval）。⚠ 口径：这是**命中行数**，`RbPanels.register(` 是构造函数调用形状，一枚注册一行。
- 对照尺（D2）：`grep -rn --include=*.js "^\s*label:" design/doubao/demo/` ⇒ 8 行，剔 `app.js:300` 命令面板假阳＝**7**，两尺合一。
- 与 09-28 `182-a1` 对照：**成员一致**（K1–K7 同名），本轮 demo 注册表没增没减（demo 目录工作树有 12 行在飞改动，本读数是工作树值，非 HEAD）。
- **另 3 枚"有名无位"**（不在注册表、文本有据）：K8 子代理（票 188 AC#1 已裁"进产品"）、K9 后台任务（demo 混写进活动面板时间线，出处＝a1 尺 `rightbar.js:90`，本程未复跑）、K10 浏览器（出处＝`issues/77:364/:384`，转引 c1/a1，未复跑）。
- ⇒ **这一栏可点名总堆数＝7（有位）＋3（有名无位）＝10**，与 a1 读数一致。

## 2. AC#2 逐堆三档定性＋尺＋命中数（尺原文都可在本仓复跑）

| 堆 | 三档定性 | 尺原文→命中数 | 归属票（§3） | 雷区（§4） |
|---|---|---|---|---|
| K1 上下文·路径/附件 | **有源＋载体已在** | `awk '/^type ComposerState struct/,/^}/' internal/panel/composer.go \| grep 'json:'`＝**11 行命中**（workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError 等）；正控＝同尺有输出 | 票 92 已落；不新开 | 无 |
| K1 上下文·占用条(used÷total) | **无源**（引 167-c2，未重跑） | 167-c2 §1.2/§1.4：ComposerState 无 used/total、PumpSources 无占用 reader、`Loop.Budgets()` 3 枚产码调用者无一取 `.ContextWindow` | 票 145 扩行＋owner 裁分子分母（167-c2 §4②在册） | 无 |
| K1 上下文·token/成本 | **有源缺载体** | 源链尺 `grep -rn "AddCost\|AddUsage" internal cmd --include=*.go \| grep -v _test.go`＝**5 行**（`cost.go:38`／`guard.go:148/:155`／`loop.go:442-443`）＋`grep -rn "publishUsage" internal cmd --include=*.go \| grep -v _test.go`＝**3 行**（定义 `internal/agent/loop.go:983`，调用 `:445/:449`）；消费尺 `grep -rni "cost" internal/panel/*.go cmd/wisp/panel_pump.go \| grep -v _test.go`＝**9 行全是注释散文、0 枚字段**（正控＝同尺能命中注释） | 票 145 行 12「成本」（`docs/evidence/s1/145-snapshot-field-census-r1.md:182` 现读仍在）＋AC#2 未勾 | 无 |
| K2 活动 | **有源缺载体（细维）** | 粗载体 Results 已在（顶层 Snapshot 尺：`awk '/^type Snapshot struct/,/^$/…' \| grep -c 'json:'`＝**6**）；细事件维转引 a1 尺 `internal/agent/sink.go:64/:85`（本程只核文件存在，未复跑行号） | 票 145（十四行表扩行；AC#2 未勾） | 无 |
| K3 审批 | **有源＋载体已在**；序号/停止维**无载体** | Pending 字段（G1 尺内 `json:"pending"`）；序号/停止＝引 167-c2 §2/§3（`LiveApproval.Position` 真源在、`NativeVerdict` 无字段、入向 6 枚无 stop） | 票 167（其落地票）＋票 35（resync 前例形） | **有**：面板批准＝见 §4 第 1 枚 |
| K4 文件树 | **无源但规格已立**（宿主侧读面零；模型侧工具在册） | 宿主尺 `grep -rn "ReadDir\|WalkDir" internal/panel cmd/wisp --include=*.go \| grep -v _test.go`＝**3 行**，逐枚过目＝`assets.go:105`（内嵌资产树）＋`git.go:465/:508`（`.git` 枚举），**无一枚工作区目录树**（正控＝同尺能命中）；模型侧尺 `grep -n "func FSListDecl" internal/tools/fs.go`＝`:309`（1 行）＋注册尺 `grep -n "BuiltinFSEntries(" cmd/wisp/run.go`＝`:515`（1 行） | 票 190（AC#1 已勾设计核，AC#2–6 未勾） | 现状无实雷：`data-*` 尺（D3）rb-files 只命中 dir/file/lucide＝展开类；删除/移动若做＝owner 逐枚批（190 AC#3 具名） |
| K5 插件 | **规格真空**（"右栏插件列表"没人要求过） | `ls internal/plugin`＝`disposal.go`/`disposal_test.go`/`doc.go` 3 件；`grep -rn "Installed\|func List" internal/plugin --include=*.go \| grep -v _test.go`＝**0 行**（⚠ 此尺 rc 被管道 head 吃，正控未打——记 §6 缺陷，相应结论标〔不可判〕从严） | **无票认领**（池内 grep「插件」＋右栏形成本程未跑第二把；候选＝立新票或摆 owner） | **有**：`data-plg-toggle`（D3：`rb-plugins.js:65/75/77`＝写配置＋改档位）；入向已有 `config.get/config.set`（P5：`bridge.go:66-67`）但插件启用在不在其可写键集**本程未量**——待人拍板 |
| K6 审查（git） | 分支/工作树＝**有源＋载体已在**；未提交枚数/diff＝**无源但规格已立** | 载体尺 `grep -n "json:" internal/panel/git.go`＝**15 行**（branch/worktrees/repoRoot/switchBlocked…）＋`grep -n "Git GitView" internal/panel/composer.go`＝`:246`＋reader 尺 `grep -n "gitView" cmd/wisp/panel_pump.go`＝`:228`；diff 尺（票面原尺）`grep -rn "git diff\|diff --numstat" --include=*.go internal/ cmd/`＝**0**（rc=1；正控＝同形尺 `git diff\|\.git/HEAD` 命中 `git_test.go:72…`，尺能响）；uncommitted 尺 `grep -rn "uncommitted\|numstat\|Uncommitted" internal/panel internal/tools cmd/wisp --include=*.go \| grep -v _test.go`＝**0**（正控＝同形尺 `json:"branch"` 命中 `git.go`） | 显示＝票 181（已落，AC#1–5 勾）；**内容审查＝票 189**（六框全未勾） | **有**：diff 接受/回滚、面板提交/推送（§4 第 3/4 枚；入向名册 P5 今天 6 枚、无一枚是动作腿） |
| K7 终端 | 地基**无源但规格在册**（D34 有 `shell.session` 行）；面板载体规格已立（票 191 裁甲） | 名册尺 `grep -n "Builtin.*Entries(" cmd/wisp/run.go`＝**3 行**（FSEntries `:515`／TaskEntries `:544`／SubagentEntries `:780`，**无 Shell 一支**）；`grep -rn "shell.session\|shell.exec" internal cmd --include=*.go \| grep -v _test.go`＝**5 行全注释/自陈**（`spill.go:26`、`internal/config/unwired.go:63/:64/:70/:76`、`rules_shell.go:14`；正控＝同尺有输出）；PTY 尺 `grep -rni "\bpty\b\|conpty" internal cmd --include=*.go \| grep -v _test.go`＝**1 行**且是 `cmd/wisp/testdata/esclistener/main.go:117` 的 `ptY` 假阳（非产码） | 地基＝票 163（AC#1b 未勾）；面板载体＝票 191（AC#1–6 未勾） | **有**：面板发起执行（D3：`data-terminal-add` `rb-terminal.js:26`）＝191 整张票就是这一枚的裁定票；甲支约束＝同一套 risk 门控（191 AC#3），乙支未批（＝C17 契约变更） |
| K8 子代理 | **有源＋载体已在（本框最大订正）**；根行状态维**缺** | 载体尺 `grep -n "Status\|status" internal/panel/subagent_roster_197.go`＝**20 行封顶**（`Status/StatusKnown/StatusReason` 三件套 `:114-120`）＋装配尺 `grep -n "Tasks: rt.taskRosterState" cmd/wisp/run.go`＝`:724`、reader `panel_pump.go:146`；缺维尺（引 167-c2 §1.3）＝`MarkRoot` 不写 State（`internal/tools/task.go` 注释现读在 `:353` 一带，rc 见 30-go-side-readings G3b） | 票 188（AC#1 已裁"进"；AC#2 状态维未勾）＋票 197（已落）＋票 196（词表/schema；其框行格式本程尺没解析到，记缺陷） | 无（188 AC#3 明令面板不许写这一维） |
| K9 后台任务 | 同 K8：**名册有源**、状态维缺、demo 无独立位 | `grep -n "^type TaskRoster struct" internal/tools/task.go`＝`:176`（1 行命中＝尺能响）＋`StateAnswer`＝`:158`；demo 位＝转引 a1 `rightbar.js:90`（未复跑） | 票 188 AC#2／票 196（词表）；"停止"维＝票 167（引 167-c2 §3，四枚入向零生产调用方名册） | 无（现状无面板按停通道——167-c2 §3.3：入向 6 枚无一与停止有关；新入向＝C17 契约面） |
| K10 浏览器 | **无源；已从"规格真空"升级为"有票但被扣"** | 尺 `grep -rlni "browser\|浏览器" internal cmd --include=*.go \| grep -v _test.go`＝**3 个文件**（`risk/blacklist.go`＝浏览器凭据拉黑分类、`panel_assets.go`/`panel_host_windows.go`＝WebView2 宿主措辞）——**逐枚过目无一枚是"栏内浏览显示"读面**（正控＝同尺有输出） | **票 168**（`built-in browser as a disable-able plugin`，七框全未勾，AC#1 被票 50/51 扣着不许开工）；文本出处＝`issues/77:364/:384` | 潜在：168 AC#5 自己写明"不许新增面板可给出允许的出口"——本程只指认，不替它裁 |

## 3. AC#3 归属结论（指认前逐张票读框，框态见 20-rail-tickets-boxes.txt）

- **票面那句"票池里今天没有任何一张票认领它"已过期**（它写于 09-28 立票时）：本程现量尺 R1 `grep -rlniE "内置终端|文件树|任务监控" .scratch/wisp/issues/*.md`＝**8 枚命中**（77/182/186/188/189/190/191/196），188/189/190/191 就是 09-28 下午按本票 §4–§8 那轮普查立出来的落地票。⇒ K3 审批序号/停止＝票 167、K6 显示＝票 181（已落）、K6 审查＝票 189、K4＝票 190、K7 载体＝票 191（甲裁 A379）、K8/K9＝票 188/197/196、K10＝票 168。
- **仍无人认领、须立票或摆 owner 的只剩两堆**：K5 插件列表（Go 零读面＋文本零认领）与 K1 占用条的**产品裁定前置**（分子分母口径，167-c2 §4②在册，不算新缺陷）。
- 可并格：K1 token→票 145 行 12（现读 `docs/evidence/s1/145-snapshot-field-census-r1.md:182` 逐字"| 12 | :3486 | 成本 |…"，AC#2 未勾）；K2 细维→票 145 扩行（同上，十四行表现量界已漂，见 §5.1-7）。
- **⚠ 共享实现特别问（AC#3，交编排者裁，本程不选）**：宿主侧读面与模型侧工具（`fs.list` 注册在 `run.go:515`、`FSListDecl` 在 `fs.go:309`）——**共享**＝D34 的风险分级（fs.list L0/L2 越界）渗进面板通道，面板显示语义被工具判级牵着走，且 190 AC#1 已定"树腿不接受面板传根"；**不共享**＝第二份枚举实现必然漂（票 189 AC#2 用"同一枚文件族不另起一份实现"防的就是这一形，git.go 今天已是先例：宿主读面独立于 tools 面）。两败具名，**归你裁**。

## 4. AC#5 雷区清单（逐堆标"有没有面板发起的批准/写动作"；单列待人拍板，⛔ 不混进只读交付）

现状合规锚（本程复核行号未漂）：入向名册 6 枚＝`internal/panel/bridge.go:42-45` 四枚 `panel.*`＋`:66-67` 两枚 `config.*`，**名册里没有 decide/allow/commit/push/revert/stop**（尺 P5）；`bridge.go:63-65` 注释逐字引 `AGENTS.md §1.2 keeps panel-side L2 "allow" off the table entirely`。

| # | 雷 | 出处（尺） | 定性 |
|---|---|---|---|
| 1 | 面板批准/允许按钮（K3） | demo `rightbar.js:143/:154/:160-161`＝approval 面板今天**只跳屏**（C3 尺）；Go 入向无 decide/allow（P5） | **待拍板**：任何"允许"入向＝`AGENTS §1.2` 逐字禁的形状＋`Q-49` 一族 |
| 2 | 插件启用开关（K5） | `rb-plugins.js:65/75/77 data-plg-toggle`＝写配置＋改档位（D3 尺）；`config.set` 入向存在但键集未量 | **待拍板** |
| 3 | diff 接受/回滚（K6） | 票面雷区名单（c1 §12 转记）＋189 AC#3 逐字"不许出现动作腿" | **待拍板** |
| 4 | 面板 git 提交/推送（K2 环境信息内） | 入向名册 0 枚对应（P5）；181 AC#3 钉"四枚 panel.* 一枚不许加" | **待拍板**（C17 契约面指认，不代改） |
| 5 | 面板终端发起执行（K7） | `rb-terminal.js:26 data-terminal-add`（D3）；**已有裁定**＝191 AC#0 裁甲（A379），乙支未批 | 甲支下走同一套门控；不再是无主雷，但**执行面本身**保持具名 |
| 6 | 附件移除（K1） | `rightbar.js:74 data-attach-del`（C3 尺）；a1 已具名（`attachment.remove` 不在名册） | **待拍板** |

## 5. 前人读数订正与冲突具名（本轮复量 vs 票面/167-c2/a1）

**5.1 票面现量（锚 `039efb47`）逐条复核**：
1. #1"只命中 1 枚"→ 现量 **8 枚**（§3）。**票面那句已过期。**
2. #3"Snapshot 四枚（composer.go:57-62）＋ComposerState 六枚（:200-207）"→ 现量 **Snapshot 六枚**（`:57-92` 区，C6 尺计数 6）、**ComposerState 11 枚 JSON 字段**（`:230-270` 区）。过期。
3. #4a `fs.go:309 FSListDecl` → **未漂**（R3：:309 仍是 func 行）。
4. #4b shell.session 未注册 → **仍真**；但自陈文件现为 `internal/config/unwired.go:63-76`（09-28 c1/a1 引作裸 `unwired.go:63-76`，包属具名补正）。
5. #4c `git diff` 尺 0 → **仍 0**（R4 rc=1＋正控 C2）。
6. #5 引 `bridge.go:42-45` → **未漂**（R6/P5）。
7. 09-28 c1 日志的 `run.go:345/:365` 注册行 → **漂至 `:515/:544`**（P8）；`loop.go:976 publishUsage` → **漂至 `:983`**（P2）。
8. 票 145 行来源 `PLAN.md:3473-3488` → **漂**：现量窗口内只剩 8 行状态（R2），十四行全表实界＝**`PLAN.md:3481-3494`**（P7 awk 计数 14），census r1 引的"行 12 成本 @ :3486"对应现位 `:3492`。**编排者要转告任何按 :3473-3488 开刀的台件。**

**5.2 状态性订正（相对 09-28 两轮）**：K8/K9 从"名册在、无状态维、无载体"升级为**载体已落地**（票 197：`Snapshot.Tasks`＋`TaskRowView.Status/StatusKnown/StatusReason`＋`run.go:724` 装配）；K10 从"规格真空须立票"升级为**票 168 已认领（被 50/51 扣闸）**；K6 显示半从"纸面分工"升级为**票 181 已落地**（a1 已报，本程复量对格：15 枚 json 字段/gitView reader 在位）。

**5.3 与 167-c2 的冲突（具名报回，不默默取其一）**：其 §1.2 写"`ComposerState` 14 枚 JSON 字段，逐枚：…"，但它自己逐枚列出的是 **11** 枚；本程现量＝**11**（G1 尺，范围 awk 计数）。占用结论本身（无 used/total 格）与 167-c2 **一致**，冲突只在枚数表述，疑为其笔误。Snapshot 顶层 6 枚与其一致（C6）。

## 6. 我这把尺的缺陷（哪几发没验、哪几处从读码推）

1. **管道尾 rc**：部分尺经 `| head`/`| grep -v` 后记录的 rc 是管道尾状态而非 grep 本体（G7b、R4 首跑等）；凡 0 命中负向句，本件都在配套件里给了**同形尺正控**或改用 `grep-rc=$?` 直读（R4 重跑 p1rc=1 是 PIPESTATUS[0]，可信）。K5「Installed|func List 0 行」那一发**正控没打**⇒ 相应结论本件按〔有读面无证据，从严记不可判〕处理（定性不受影响：目录清单 3 枚文件本身就是负向正证）。
2. **未复跑而转引**：K9 的 `rightbar.js:90`、K10 的 `issues/77:364/:384`、K2 的 `sink.go:64/:85` 三发转引 a1（`.scratch/wisp/probes/182/a1/census.md`），本程只核了 `internal/agent/sink.go` 存在性，**行号未复核**。
3. **票 196 框态没读到**：`- [ ]` 形制尺在 196 上 0 命中（C4），它的 AC 框可能用了别的列表符号——归属表里 196 那一格只到"票存在且主题对格"级。
4. **demo 读数是工作树值**：`design/doubao/demo` 有 12 行在飞改动（他程名下），D 系尺若 HEAD 版与工作树版不同，本件没做双版对照（注册表 7 枚与 a1 一致是旁证，不是证明）。
5. **K2"活动"细维没现量**：只定了粗载体 Results 在位＋转引 a1；"工具起止/流事件进快照哪一格"这一跳本程没量。
6. **config.set 可写键集没量**（§4 第 2 枚雷的射程）。
7. **枚数口径**：本件所有"命中数"均注明是**行数**；`RbPanels.register(`/`FSListDecl`/`TaskRoster` 等类型全名＋调用形状均按派单口径写明（构造函数调用点与注册点分开写）。
