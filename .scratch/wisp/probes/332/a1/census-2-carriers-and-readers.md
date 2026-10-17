# 332-a1 件2＝AC#2 逐枚可达性（三形入口逐环 `file:line`）＋AC#3 载体正控＋筛检尺

锚＝`7535d83b`；行号全部现取自 `git show 7535d83b:<路径>`，⛔ 抄旧普查行号。
射程口径先钉死：本件的"可达"＝**枚枚走到 marshal 出的快照字节里**（不是"屏幕上已出现"），
最后一跳今天走到哪儿、哪一环已腐烂，见 §④。

## ① 三形入口逐环（现量链，环环带 `file:line`）

**(a) 磁盘／工作区真名** — 链 today 通，且**同一链上有两枚被洗、四枚不被洗**
1. `cmd/wisp/panel_pump.go:128-132` `workspaceView()` → `panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())`
2. `internal/panel/workspace.go:114-119` 填 `Set/Spelling/Canonical`（`:115` 只对空串做 `strings.TrimSpace`）
3. `internal/panel/pump.go:275` `NewComposerState(mode, workspace, nil, maxAttachment)`（★附件那一参 today 传 `nil`）
4. `internal/panel/pump.go:309` `NewSnapshot(...)` → `internal/panel/composer.go:145` 唯一生产调用 `scrubComposerForDisplay`
   → `internal/panel/display_scrub_216.go:228` 只洗 `Workspace.Canonical`
5. marshal＝`internal/panel/pump.go:359`（`Publish`）与 `:344`（`Marshal`）两处 `json.Marshal`
6. 出口＝`cmd/wisp/run.go:744` `Out: rt.bookPanelSnapshot` → 见 §④
同族另两条**今天根本不经过第 4 环**的载体（★票 327 普查只把它们当"相邻面"，本腿按字段实量列为独立载体）：
- `instructions.files[].path`：`cmd/wisp/run.go:733` `Instructions: rt.instructionBundle` → `cmd/wisp/panel_pump.go:153`
  → `internal/panel/instructions_200.go:147` `Path: f.Path` → `internal/panel/pump.go:314`（**赋值在 `NewSnapshot` 之后** ⇒ 不过显示边界）
- `git.*`（`repoRoot`/`branch`/`detachedSha`/`currentWorktree`/`worktrees[].path|branch`/`branches[]`）：
  `cmd/wisp/run.go:718` `Git: rt.gitView` → `cmd/wisp/panel_pump.go:273-274` → `internal/panel/git.go:251` `ReadGitForWorkspace`
  ＋`:543` `GitWorktree{Path, Branch, Sha}` → `internal/panel/pump.go:281` `composer.Git = p.src.Git()`（`scrubComposerForDisplay` 不碰 `Git`）
  ★这一维的字符集完全由**读盘**决定：`internal/panel/git.go:24` 自陈 "NO EXTERNAL PROCESS. Everything is os.Stat / os.ReadFile / os.ReadDir"，
  尺＝`git grep -n 'exec\.Command' 7535d83b -- internal cmd | grep -v _test` 命中 5 行、**无一条在 `internal/panel`**（正控＝同尺在 `cmd/wisp/doctor.go:316` 有命中）。

**(b) 模型正文／工具 reason** — 链 today 通，两枚字段都不被洗
1. `internal/llm/events.go:121-124` `StreamEvent{Kind, Text}`（`Text` 注名服务 `EvTextDelta`/`EvReasoningDelta`）
2. `cmd/wisp/run.go:1343` `c.stream.Append(e.TaskID, e.Text)` → `internal/panel/pump.go:524/532` `StreamLog`
3. → `internal/panel/composer.go:59` `Results []ResultChunk json:"results"`（`NewSnapshot` 只按值搬，⛔ 洗）
4. 工具侧 reason/args：`internal/tools/bridge.go:301` `dec.Reason = verdict.Reason`
   → `internal/panel/approval.go:77-83`（`Tool`/`Args`/`Reason` 三枚逐字进卡，`:83` `Reason: decision.Reason`）
   → `cmd/wisp/run.go:715` `Verdicts: rt.liveVerdicts` → `cmd/wisp/panel_pump.go:62-66` `rt.gate.Queue().LiveApprovals()`
   → `internal/panel/pump.go:241-247` `cards = append(cards, v.CardView())` → `pending[].tool|args|reason`（⛔ 洗）
5. 同一条 marshal 出口（a 的第 5 环）。
⚠ 唯一碰到这一支的洗法是 `cmd/wisp/approval_reply.go:571/575` 的 `sanitizeReplyText`，
射程＝**操作员自己敲的拒绝理由**（`:599` 自陈），且它只换 `unicode.IsControl`（＝`Cc`），对 `Cf` 零作用（量过：155 枚里 `IsControl` 0 枚）。

**(c) 页面／入向** — 链 today 通（★这一条是票 327 普查之后才接上的，见 §④ 腐烂清单 ③）
1. 门名＝`cmd/wisp/panel_host_windows.go:80` `panelDispatchBinding = "wispDispatch"`；绑定回调 `:930` `reply, _ := m.dispatchRaw(ctx, raw)`
2. `cmd/wisp/panel_host_windows.go:941-948` `dispatchRaw` → `:948` `return m.disp.Handle(ctx, raw)`
3. `internal/panel/composer_dispatch.go:176` `Handle` → `:198-209` `dispatch`，`:209` `d.Workspace.HandleWorkspaceRequest(ctx, req)`
4. `internal/panel/workspace_handlers.go:108` → `:125` `RequestWorkspaceSwitch(h.Scope, req.Path, h.Audit)`（`req.Path`＝页面敲的串）
5. `internal/panel/workspace.go:140` → `:169` 填 `Spelling`（用户拼写原文）＋`Canonical`（C26 解析结果）→ 回到 (a) 第 4 环出向
   ⇒ **`Canonical` 过显示边界（名单内 15 枚会被显出来，名单外 155 枚原样穿）、`Spelling` 一字不洗直接出向**（尺＝`display_scrub_216.go:228/234` 两处赋值就是洗的全部）。
6. 常驻装配＝`cmd/wisp/panel_resident_windows.go:223-228` `newResidentComposerDispatch` + `NewWorkspaceScopeHolder()`；
   CLI 那一形传 `nil`（`cmd/wisp/panel_inbound.go:223`）⇒ 按名拒绝，但常驻腿这条路是真的。
7. 入向名册＝7 枚（`internal/panel/bridge.go:42,43,44,45,66,67,75`）。
8. **今天断的一支**：`Attachments[].Name`（唯一另一枚被洗字段）无源——
   `internal/panel/attachments.go:199` `Ingest` 的生产调用者 **0 枚**
   （尺＝`git grep -n '\.Ingest(' 7535d83b -- cmd internal | grep -v _test | wc -l` → **0**；
   ★正控＝**同一把尺形、同一射程**对 `json\.Marshal` 命中 **37** 行 ⇒ 那枚 0 是真 0、不是尺失灵；见 `logs/ruler-ingest-zero-with-control.log`）
   ＋`cmd/wisp/panel_inbound.go:301-302` `Attachment: nil` / `Message: nil`
   ＋(a) 第 3 环那个 `nil` 参 ⇒ 那枚字段是"洗法已就位、源与出向 0"（这一条与 `327-a1` 同读数，**未腐烂**）。

## ② AC#2 逐族可达读数（155 枚分 11 族，族号＝件1 的名册）

**先交一条决定性的"无筛选"尺**（量的）：尺＝`git grep -n 'ContainsAny\|IndexAny\|ContainsRune' 7535d83b -- internal cmd | grep -v _test`
命中 **18** 行（现量枚数），逐行读其字符集合＝全部只点 ASCII 标点或 `\x00\r\n\t` 一类控制字节
（唯一碰"名字"的那几枚＝`internal/panel/attachments.go:301`（只拒 `\x00\r\n\t`，且那一支源 0 读者）、
`internal/memory/artifacts.go:336-337`（只拒 `/ \ :`）与 `internal/secret/refs.go:39`（只拒空白与 `=`/`:`））
⇒ **三形入口上没有任何一处按 `Cf` 筛字符**，
`internal/risk/pathresolver.go` 里字符类拒绝 **0 处**（尺＝`git show 7535d83b:internal/risk/pathresolver.go | grep -c 'unicode\.\|IsPrint\|ContainsAny\|IndexAny'` → **0**；
★正控＝同一把正则对 `display_scrub_216.go` 命中 **5** 行 ⇒ 那枚 0 是"真的不筛"、不是尺失灵）。
⇒ 下面每一族拿到的载体枚数**相同**，本腿就逐枚给枚数、⛔ 不写成"仅理论形"。

| 族（枚数） | 载体(a) 磁盘真名 | 载体(b) 模型正文／reason | 载体(c) 页面入向 | `file:line` 落点 |
|---|---|---|---|---|
| A 前缀标记（13） | **有** | **有** | **有** | (a)`pump.go:314`/`git.go:251`＋(b)`run.go:1343`＋(c)`workspace_handlers.go:125`→`workspace.go:169` |
| B `U+061C`（1） | **有** | **有** | **有** | 同上，逐枚同链 |
| C `U+180E`（1） | **有** | **有** | **有** | 同上；★`unicode.IsSpace(U+180E)`＝false（量的，155 枚里 `IsSpace` 0 枚）⇒ `workspace.go:115` 的 `TrimSpace` 去不掉它 |
| D `U+200E/F`（2） | **有** | **有** | **有** | 同上；与已洗的 `U+202A–U+202E` 同一支 bidi，今天照穿 |
| E 数学作用符 `U+2061–U+2064`（4） | **有** | **有** | **有** | 同上＋★`U+2061` 本腿真盘＋marshal 走过一遍（§③） |
| F `U+206A–U+206F`（6） | **有** | **有** | **有** | 同上＋★`U+206B` 真盘＋marshal 走过一遍（§③） |
| G 行间注释 `U+FFF9–FFFB`（3） | **有** | **有** | **有** | 同上 |
| H 圣书体（16） | **有** | **有** | **有** | 同上（`U+13439–U+1343F` 为 16.0 新进 `Cf`，本腿码位表已含） |
| I 速记（4） | **有** | **有** | **有** | 同上 |
| J 音乐（8） | **有** | **有** | **有** | 同上 |
| K Tag（97） | **有** | **有** | **有** | 同上；★这一族的官方名逐枚＝"TAG ＋ 一枚 ASCII 字符名"（`U+E0021`–`U+E007E` 95 枚一对一，尺＝名表前缀逐枚数） |
| **求和** | **155/155** | **155/155** | **155/155** | 无任一族写"无已知载体"；★唯一写"无已知载体"的是**字段**不是字符：`Attachments[].Name`（§①(c) 第 8 条） |

逐族载体枚数之所以齐平，是上面那把"无字符类筛选"的量尺的直接结果；
本腿⛔ 据此下"要不要扩"的判语，也⛔ 替任何一族担保"屏幕上一定看得见"（射程止于 §④）。

## ③ AC#3 两枚名单外字符的**真盘上正控**＋marshal 字节 dump

唯一执行过的代码＝三枚仓外单文件纯 stdlib 探针（⛔ import 本仓包、⛔ `go test`/`go build`/`go vet`）：
`D:\tmp\332a1\cfprobe\main.go`（尺②）、`D:\tmp\332a1\carrier\main.go`（本节①②段）、
`D:\tmp\332a1\diskroster\main.go`（本节③段），另加 `D:\tmp\332a1\mncount\main.go`（尺④补尺）；逐字源码已入库＝`logs/probe-sources.log`。
其中 `scrubForDisplay`/`isDisplayWashCf`/`writeControlEscape`/`needsDisplayWash` 四段＝`display_scrub_216.go:106-136/164-175/179-199/211-220` 的
**逐字重打**（函数体一字未改）；`isIgnorableTaintRune`＝`taintmatch.go:39-45` 逐字重打（只为量，⛔ 动原件一字）。

| 码位 | 盘上真名（写→`os.ReadDir` 读回） | 盘上 UTF-8 字节 | 过显示边界 | marshal 字节（同 `pump.go:359` 那把尺） | marshal sha256 |
|---|---|---|---|---|---|
| `U+2061`（名单外） | `photo<2061>report.txt` 建名成功、readdir 命中 1 | `70686F746F E281A1 7265706F72742E747874` | `needsDisplayWash=false` ⇒ **走快路、返回同一枚串** | `2270686F746F E281A1 7265706F72742E74787422`＝`"photo⁡report.txt"`，**⛔ 转义形** | `da1e771280e63fd9…b3af889b` |
| `U+206B`（名单外） | `photo<206B>report.txt` 建名成功、readdir 命中 1 | `70686F746F E281AB 7265706F72742E747874` | 同上 | `2270686F746F E281AB …`＝`"photo⁫report.txt"`，⛔ 转义形 | `bb7a8ae1915fe249…ffdd4db` |
| `U+202E`（名单内正控） | 建名成功（盘允许） | `…E280AE…` | `needsDisplayWash=true`、被换形 | `…5C5C 75323032 65…`＝`"photo\\u202ereport.txt"` ⇒ **`E2 80 AE` 在 marshal 字节里 0 命中** | `158adf0dba164bed…c29bb5` |
| `U+001B` ESC（`Cc` 正控） | **盘上写不出来**：`The filename, directory name, or volume label syntax is incorrect` | — | `needsDisplayWash=true`、被换形 | raw＝`…5C7530 303162…`（`json` 把 `Cc` 转义成 `\u001b`）；洗后＝`"photo\\x1breport.txt"` | `59fc78ba9ca0b01d…813661` |

⇒ 两枚名单外字符在"盘上真名 → 逐字重打的显示边界 → `json.Marshal`"这条链上**逐字节原样出现**（`E2 81 A1`／`E2 81 AB` 在位），
而同一探针里的两枚名单内／`Cc` 正控**各自被换形** ⇒ 那把尺对"洗"敏感，上面的"穿过"不是探针失灵造成的假读数。
★顺带一枚本腿量到的不对称（只报数、⛔ 引申）：**`Cc` 反而不能作 NTFS 真名**——对照 6 枚
`U+0000/U+0009/U+000A/U+000D/U+001B/U+001F` 全部 `WRITE_FAIL`，而 `U+0085`/`U+009B`（`C1`）与 `U+00AD`/`U+200B`/`U+FEFF`（名单内 `Cf`）全部 `READBACK_OK`。

**逐枚真盘尺扩展到全部 155 枚**（量的，`logs/ac2a-disk-legality-per-cf-codepoint.log`）：
155 枚一枚一枚写真名 `p<cp>q.txt` → `os.ReadDir` 读回 → 逐枚比对＝**155/155 `READBACK_OK`、`fail=0`**；
名册末行 `DISK_REMAINING=164`＝155 枚 ＋ 5 枚对照成功行（`U+0085 U+009B U+202E U+200B U+FEFF`）＋ 该目录里探针自己的 4 枚文件（`main.go`/`go.mod`/两份输出），
第三来源逐枚对合（`D:\tmp\332a1` 下 `python os.listdir` 尺）＝`P_COUNT=155`、`MISSING_ON_DISK=[]`、`EXTRA_ON_DISK=[]`、`CF_UNIQUE=True`
⇒ 载体 (a) 对**全部 155 枚**都是今天就写得出来、读得回来的真名，且 155 枚同目录并存互不覆盖。

**四把生产者侧筛检尺**（量的，逐枚跑满 155 枚，`logs/ac3-carrier-disk-and-marshal-dump.log` 末行）：
`SCREEN total_outside=155 scrub_leaves_raw=155 unicode_is_space=0 taint_ignores=0 unicode_is_control=0 json_escapes=0`
⇒ 名单外那 155 枚里：被 `scrubForDisplay` 换形的 **0** 枚；被任何 `TrimSpace` 一类去掉的 **0** 枚；
被判定层 `isIgnorableTaintRune` 忽略的 **0** 枚（＝两表在 `Cf` 轴上只交在那 6 枚，尺④）；
是 `Cc` 的 **0** 枚；被 `encoding/json` 转义成 `\uXXXX` 的 **0** 枚。

## ④ 最后一跳今天走到哪儿＋票 327 普查件的腐烂清单（★增量料）

现量链条：`cmd/wisp/run.go:744` `Out: rt.bookPanelSnapshot` ⇒ `internal/panel/pump.go:364-367` 那句
"快照已构造但没有出口（最后一公里在票 33/35 名下）"的错误支**今天不走**（它只服务 `src.Out == nil` 那一形）→
`cmd/wisp/panel_pump.go:447 publishPanelSnapshot` → `:463 rt.pushPanelSnapshotToHost()` →
`cmd/wisp/panel_snapshot_deliver_329.go:95` `_, data, seen := rt.lastPanelSnapshot()`（★第一枚生产读者）→
`:105 h.evalPanelJS(panelSnapshotEvalJS(data))` → `:66` 逐字包裹成 `window.wispOnPanelSnapshot(<marshal 字节原文>)` →
`cmd/wisp/panel_host_windows.go:331/:336` `PanelManager.evalPanelJS` → `w.Dispatch(func(){ w.Eval(js) })`。
页面那一侧仍无接收器：尺＝`git grep -c 'wispOnPanelSnapshot' 7535d83b -- frontend | wc -l` → **0 枚文件**（现量 0 命中）；
★正控＝同一把尺去掉目录限定后命中 **16** 行（全在 `cmd/wisp/panel_snapshot_deliver_329.go` 与 `.scratch`/`docs` 的记录里）⇒ `frontend` 那枚 0 是真 0。

**`327-a1` 读数里今天已腐烂的那几环**（逐枚具名，⛔ 我改它一字）：
1. ★**腐烂最硬的一环**：它写的"快照字节今天**没有进 WebView 的生产路径**／读它的 `lastPanelSnapshot` 生产调用者 0 枚"
   ⇒ 现量读者 **1 枚**＝`panel_snapshot_deliver_329.go:95`，且 `Out` 已非 nil（`run.go:744`）。它那句"最后一公里没接"作为读数**已过期**；
   今天准确的射程＝**marshal 字节已进 WebView2 的 `Eval` 调用串**，`⛔ 像素已出现` 那一半仍成立（页面接收器 0 命中）。
2. 行号漂移：它引 `display_scrub_216.go:169`（`Canonical`）/`:175`（`Name`）⇒ 现量同一处赋值在 **`:228`/`:234`**。
3. ★"外层 `RequestWorkspaceSwitch` 在 HEAD 上非测试零调用者"（它在 §⑤ 的上游尺，票 181/186 时期的读数）
   ⇒ 现量非测试调用者 **1 枚**＝`internal/panel/workspace_handlers.go:125`，且常驻链一路可达页面门
   （`panel_host_windows.go:948` → `composer_dispatch.go:209`）。⇒ 本文 §①(c) 那条"敲进去的 `Spelling` 一字不洗直接出向"的载体**今天才真的通**。
4. 它那一族候选只有 **22 枚**、且当时显示边界 `Cf` 洗 0 枚；今天 `Cf` 全集 170 枚、边界已洗 15 枚 ⇒ 待答集合＝155 枚（本文的名册）。
5. **未腐烂、本腿复跑同读数**的两条：`Attachments[].Name` 那枚字段"洗法已就位、源与出向 0"（`Ingest` 生产读者 0、`panel_inbound.go:301` `nil`、`pump.go:275` 传 `nil`）；
   `results/pending/git/instructions` 诸字段不过显示边界（现量与 §① 逐环同）。
