# 327-a1 §① 字段名册：今天真能带着 `Cf` 走到界面上的字段

腿 `327-a1`（只读普查）｜锚＝具名提交 `44f51fc8`（读数一律取自**工作树现行字节**，票面 33 行点名）｜2026-10-16
⛔ 本件不含判语、⛔ 不含修法建议（修法三形归票 327 `AC#1`，由编排者裁）；决定性尺的逐字转储在 `logs/readings.txt`（摘录），全量 497 行在同一次运行的仓外件 `D:\tmp\327a1\logs\readings-raw.txt`。
原始命令逐枚同屏附在下面。凡"屏幕上哪一行"那一列引的 `frontend/**` 行号＝**读数不是判据**：那棵树此刻由机主的外部前端 agent 在飞改写（起手 `git status` 可见 ` M frontend/**`）。

## 0. 显示边界今天到底洗了什么（尺 R1，`grep -rnE 'scrubComposerForDisplay|scrubForDisplay|needsDisplayWash|writeControlEscape' --include=*.go .`，rc=0）

生产调用点**只有 1 枚**：`internal/panel/composer.go:145` `Composer: scrubComposerForDisplay(composer)`（其余命中全在 `display_scrub_216.go` 自身、`display_scrub_216_test.go` 与 `.scratch/wisp/probes/216/r1/restore-copy-*` 那枚还原副本）。
洗的字段集合＝`display_scrub_216.go:169`（`Workspace.Canonical`）＋`:175`（`Attachments[].Name`）**两枚，逐字到此为止**——同一枚函数里没有第三行赋值。
⇒ **名册外的 `Snapshot` 字段今天一律裸奔**（`Pending`／`Results`／`Instructions`／`Tasks` 都是 `Snapshot` 的兄弟字段，不在这枚函数的射程）。

字符集读回：`display_scrub_216.go:114 case unicode.IsControl(r)`／`:155 unicode.IsControl(r)`；`:52–58` 的注释逐字写着不洗 Cf（"it does not wash Cf (format) characters such as U+202E RTL OVERRIDE or U+200B ZERO WIDTH SPACE"）。
`needsDisplayWash` 对 22 枚候选 `Cf` 逐枚**不响**（含 `U+202E`、`U+200B`、`U+FEFF`），对 `ESC`/`LF` 两枚 `Cc` **响**＝正控，见 §② 探针表（同一枚判据在件里逐字重打，⛔ import 本仓包）。

## 1. 票 216 `AC#1` 那枚名册，逐枚四问

| 字段 | ① 字节从哪来 | ② 走不走 `scrubForDisplay` | ③ `needsDisplayWash` 对含 `U+202E` 的样例 | ④ 屏幕上哪一行渲染 |
|---|---|---|---|---|
| **附件名** `AttachmentRef.Name`（`json:"attachments"`） | 页面 picker 的 `file.name`（磁盘真名）：`frontend/src/lib/panel.ts:298 sendAttachmentBytes(file.name,…)` → `AttachmentPayload.Name` → `attachments.go:403 DisplayName: p.Name` → `:202/:256` | **走**（`display_scrub_216.go:175`） | **不响**（`IsControl(U+202E)=false`）⇒ `scrubForDisplay` 首行 `return s`，字节原样 | **今天到不了屏幕，两头都断**：① 入向腿在装配根就被写成 nil——`cmd/wisp/panel_inbound.go:301 Attachment: nil`（`:302 Message: nil` 同形，旁边注释"both still refused by name when they arrive"），`MethodAttachmentAdd` 分支（`composer_dispatch.go:210–214`）因此只回 `unattached`；② `AttachmentBroker.Ingest` 生产调用者 **0 枚**（尺 R40：`grep -rn '.Ingest(' --include=*.go .` rc=0，19 枚命中全在 `_test.go` 与 `.scratch` 探针）；③ 泵写死 `NewComposerState(mode, workspace, nil, maxAttachment)`＝`pump.go:275`，附件切片恒 `nil` ⇒ `"attachments"` 段今天恒空。**名册里的这一枚＝洗法已就位、源与出向 0**（同票 216 `216-r1` 残余 ③） |
| **插件名** | **零源**：`internal/panel/**` 里 `plugin` 只命中 `config_handlers.go:109` 的锁定族名与注释（尺 R41） | 不适用（无字段） | 不适用 | 无。票 213 未落地 |
| **技能名** | **零源**：同枚尺，`skill` 在 `internal/panel/**` 非测试字节里 0 命中 | 不适用 | 不适用 | 无。票 215 未落地 |
| **目录名** `WorkspaceView.Canonical`（`json:"canonical"`） | **磁盘真名**：`cmd/wisp/panel_pump.go:128 WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())` → `workspace.go:119,120` 直取 C26 `res.Canonical`（NTFS 目录名允许 `U+202E`，`PathResolver` 侧不洗字） | **走**（`display_scrub_216.go:169`） | **不响** ⇒ 原样穿过这道边界 | `frontend/src/components/composer.tsx:150`（`title=`）与 `:152`（行内文本）；键在 `frontend/src/lib/panel.ts:85 canonical: string` 有声明 |
| **命令说明**（今天形＝审批卡上的 argv/命令行） | **模型正文**：LLM tool-call 的 args JSON → `internal/llm/events.go:239 Args = append(…, ev.ArgsDelta…)` → `internal/tools/bridge.go:382 Args: req.Args` → `cmd/wisp/panel_pump.go:68 Args: panelArgs(it.Decision)` → `panelArgs` `cmd/wisp/panel_pump.go:297–312`（`argv`／`command`／**整段原始 JSON `[]string{raw}` :311**） | **不走**——`ApprovalCardView` 是 `Snapshot.Pending`，不在 `scrubComposerForDisplay` 射程（该函数只碰 `Composer` 段） | 不适用（它压根没走到这枚判据） | `frontend/src/components/l2-approval-card.tsx:152`（展开态 `{[tool, …args].map(…)}`＋`:160 {JSON.stringify(element)}`）／`:164`（收起态 `[tool, …args].join("\n")`）／`:305`（`<ParamsBlock args={view.args} …>`）；入口链 `frontend/src/App.tsx`→`approval-screen.tsx:103 {pending.map((view) => (` |
| **错误原因**（今天有 5 枚活字段） | 逐枚：`ApprovalCardView.Reason`＝**模型正文＋判定文本**（`internal/risk/rules_shell.go:51` 把 `ShellArgv[0]` 逐字写进 `"R6: 程序不在 shell 白名单: %s"`）；`WorkspaceView.Reason`＝**磁盘真名**（`workspace.go:124–126` 把 `res.Canonical`＋`res.Spelling` 拼进句子）；`GitView.Reason`／`InstructionsSection.Reason`／`TaskRowView.StatusReason`＝各自读数的句子 | **全都不走**：`scrubComposerForDisplay` 只碰 `Workspace.Canonical` 与 `Attachments[].Name`，**`Workspace.Reason` 与 `AttachmentReason` 都是同一段的兄弟字段、没被洗**（`:169/:175` 那两行之外无赋值） | `Workspace.Reason`/`AttachmentReason` 到了判据跟前：**不响**；卡与 git/instr/tasks 的 `Reason`：没走到判据 | `reason` 在 `frontend/src/lib/panel.ts` 命中 6 处（卡与 composer 行）；`instructions`/`git`/`repoRoot`/`branch`/`worktrees` 在该文件命中 **0** ⇒ 这几枚 Go 侧已出向、页面侧今天没有声明（`App.tsx:217` 给 `GitBranchChip` 的是硬编码空对象 `{current:"", isRepo:false,…}`，注释说"Reads its strings from the snapshot"、字节里没读——⛔ 本腿不判它，只录读数） |

## 2. 名册外、但今天真在出向的 `Cf` 载体（必答题要用，逐枚带尺）

| 字段 | 源 | 洗？ | 屏幕 | 尺 |
|---|---|---|---|---|
| `ResultChunk.Text` | **模型正文**（SSE 增量）：`cmd/wisp/run.go:1328–1332 consoleSink.Publish` `case agent.EvTextDelta: … c.stream.Append(e.TaskID, e.Text)` → `internal/panel/pump.go:524–538 Append` 里 `cur.Text += text` **逐字累积**（无字符筛） | **不洗**（`Snapshot.Results` 不在 `Composer` 段） | `frontend/src/App.tsx:461–463 {chunk.text}`；`App.tsx:432 chunk.text.slice(0, 18)`（turn rail 标签）；`result-stream.tsx:37 {chunk.done ? chunk.text : <RevealText text={chunk.text}/>}` | R1 命中集里没有 `results`；`pump.go:344,359 json.Marshal` |
| `TaskRowView.Label` | **模型正文**：`internal/tools/subagent_197.go:183` 清单声明 `description` 进名册当 label → `:246 label := strings.TrimSpace(a.Description)` → `:251 capRunes(label,120)` → `:333/:368 PublishSubagent(taskID, parentID, label)` → `cmd/wisp/panel_pump.go:217 Label: rec.Out.Label` | **不洗**（`Snapshot.Tasks` 不在 `Composer` 段） | 页面侧 `panel.ts` 无 `tasks` 键（尺 R52 `grep -rn 'tasks' frontend/src/lib/panel.ts` **rc=1**，同把尺对 `canonical`/`args`/`text` 命中 1/1/5 ⇒ 非恒 0 尺）⇒ **出向已到 JSON、屏幕今天不读** | `internal/panel/subagent_roster_197.go:108` |

## 3. `Snapshot` 的键与"洗没洗"总账（尺＝`grep -n 'json:' internal/panel/composer.go`）

`Snapshot` 六枚键（`composer.go:58–91`：`pending`/`results`/`composer`/`generatedAt`/`instructions`/`tasks`）里，**过显示边界的只有 `composer` 一枚**，而这枚里面又只有 `workspace.canonical` ＋ `attachments[].name` 两处赋值。
⇒ 一句话读数：**今天"洗过"的字段 2 枚，"出向且没洗"的字段 7 枚**（`pending[].tool/args/reason`、`results[].text`、`composer.workspace.reason`、`composer.attachmentError`、`tasks.rows[].label`、`instructions.files[].path`、`git.*`）。是否要扩、扩到哪几枚＝`AC#1` 的裁定，⛔ 本腿选型。
