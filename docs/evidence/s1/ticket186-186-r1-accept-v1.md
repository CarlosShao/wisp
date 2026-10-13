# 票 186 裁决表 — 落地腿 `186-r1`（非实现者验收 · 腿 `186-v1`）

- 裁决腿：`186-v1`（非实现者）。锚：`HEAD=044b9066`（交件 `341ca42e` 是其祖先，`git merge-base --is-ancestor`=YES）；分支 `dev`。
- 交件：`341ca42e`（票尾 22:3x 完整记录；该腿死于 ECONNRESET）。裁决时刻 2026-10-13 23:1x–23:4x。
- 读数落点：`.scratch/wisp/probes/186/v1/logs/`（件 00–19）。攻击性读数全部带还原证（仓外备份 `/c/Users/swq/AppData/Local/Temp/wisp186v1-backup/` ＋ sha256 ＋ `git diff` 空）。
- 树会呼吸：现跑期间 `cmd/wisp/panel_config_store.go` 被并行腿 `313-r2` 改、`internal/tools` 被 `172-r1` 改；本表锚在 `341ca42e` 的对象层（被改动/被扫描文件 blob hash 逐枚对上）。

## 0. 基线（读数）

| 读数 | 结果 | 件 |
|---|---|---|
| 三枚具名仪器（编排者前置＋本腿复跑） | `ok` rc=0 | `logs/01-baseline-named3.txt` |
| `./internal/panel/` 整包 | 恰 **5 枚红**＝已知红名册五枚（ApprovalCardViewJSONKeys…／C21DesignTokens…／ComposerContractTypes…／StreamLogDropped…／StreamLogFlood…）；**无 186 新增红** | `logs/02-baseline-panel-pkg.txt` |
| 被扫描/被改文件 blob hash | 与 `341ca42e` 逐枚相同（`5b3e5e73`／`4db510ce`／`5df8f348`） | `logs/03-mutation-pre-sha.txt` |
| 两枚在册 `GitSwitchBlockedReason` 非空判据 | 各现跑一发：`TestReadGitSyntheticSymrefRepo`／`TestReadGitNonRepoSaysSoOutLoud` 均 PASS | `logs/10-two-registered-git-tests.txt` |

## 1. 突变台（攻"判据真假"：每枚仪器弄坏必须红）

| # | 自造突变（⛔ 不抄交件的形） | 仪器读数 | 还原证 | 件 |
|---|---|---|---|---|
| M1 | `git.go` 尾种 `var plantedWrite186 = os.WriteFile`（写形＋第三枚文件） | **红**（红句具名 `git.go:582`） | sha256 `1419e6e1…` 逐字同＋diff 空 | `04-mutation-M1-writefile.txt` |
| M2 | `workspace.go` 加 `"os/exec"`＋`_ = exec.Command("git","reset","--hard")` | **红**（红句具名 `workspace.go:29/185`） | sha256 `0330e416…`＋diff 空 | `05-mutation-M2-reset.txt` |
| M3 | `workspace_handlers.go` 种 `const plantedRoute186 = "panel.worktree.switch"` | 守卫**红**（"a second door"） | sha256 `2ea25cd2…`＋diff 空 | `07-mutation-M3-route-literal.txt` |
| M4 | `composer_dispatch.go` 种**第二条路由**（第二个 switch 含 `case MethodWorkspaceRequest:`） | 守卫**红**（2 case labels ≠1） | sha256 `d001fb22…`＋diff 空 | `08-mutation-M4-second-route.txt` |
| M5 | `git.go` 种一行注释含 `.SetWorkspaceRoot(`（第三调用点名册） | 守卫**红**（`[git.go workspace.go workspace_handlers.go]`） | sha256 `1419e6e1…`＋diff 空 | `09-mutation-M5-third-site.txt` |
| M6 | 包内**新增第四枚产码文件**含 `exec.Command` | 仪器**绿**＝射程盲区（只扫 3 枚具名文件） | 文件已删、`git status` 空 | `16-M6-blindspot.txt` |
| M7 | `git.go` 种 `var plantedTemp186 = os.CreateTemp` | 仪器**绿**＝regex 漏 `CreateTemp` 一族 | sha256 `1419e6e1…`＋diff 空 | `17-M7-createTemp-evasion.txt` |
| — | 还原后复跑三具名仪器 | `ok` | — | `03-mutation-pre-sha.txt`（尾段） |

M6/M7＝具名残差（不驳回任何已声称判据；两枚都**不在**仪器的自称射程/正控形里，记 §5 供编排者排票）。

## 2. AC#1–AC#9 逐格判语

- **AC#1 成立**（读面由票 181 落，本腿复认；本腿现跑独立复认：`.git/HEAD`=`ref: refs/heads/dev` 对真值 `dev`；`.git/worktrees` 目录 2＋主仓 1＝3 对 `git worktree list --porcelain` 3；递归 `refs/heads`=4 对 `git branch --list`=4；`packed-refs` NO；`exec.Command("git")` 非 test＝**0**）。尺＝`logs/15-ac1-ac5-rerun.txt`。
- **AC#2 成立**（复用已在册的 `panel.workspace.request`，白名单现读仍 4 枚 `panel.*`；M3 证明自造第 5 枚名字会被守卫打红）。将来"新建工作树"确需新名⇒先入档（见 §4 欠账①）。
- **AC#3 成立**（三问各带判据）：(i) 选"照切"＋上界＝scope 一个字段＋一行审计；M1/M2 证明"这条腿不能动树"的仪器两向之外我自造两向仍红；残差 M6/M7 记 §5。(ii) 同源＝指针级行证链：`run.go:501` 唯一 canonicalizer → `:516`/`:747`（fs 工具与 bridge 同用）→ `panel_pump.go:87`（显示读同一 `rt.paths`）→ `resident_task_source_windows.go:300/311`（`src.run=run`＋`ws.Bind(run.paths)`）→ `workspace_handlers.go:188-199/244-254`（holder 存/转同一 scope）→ `workspace.go:122`（唯一收窄点 `SetWorkspaceRoot`）；`Bind` 的 `!=` 对 `*PathCanonicalizer`（`paths.go:65` 返指针）＝指针等价。件＝`logs/06-same-source-lines.txt`。(iii) 两形审计（`SWITCH`/`SWITCH-REFUSED`）＋handler 未达腿时的 `WORKSPACE-REFUSED`，均在册测试绿。
- **AC#4 成立**（票 181 `ReadGit` 的 `not_a_repo`+人话；`TestReadGitNonRepoSaysSoOutLoud` 绿；`GitSwitchBlockedReason` 改字面后两枚在册非空判据复跑绿，语义新读法已写进 ask 表）。
- **AC#5 成立**（`TestGitDimensionHasNoModelCallableTool` 绿＝D34 零 git 行／`internal/tools` 零 `git.*` 名／白名单恰 4 枚；`exec.Command("git")` 非 test=0；交件未动 `docs/PLAN.md`、`internal/tools`）。
- **AC#6 成立**（五样全"保留"；我抽 2＋加 1 现跑复认：③ C25 名册键＝`dec.TaskID`（`internal/tools/bridge.go:650`）与 `OpenScope(taskID)`（`:856`），`internal/risk/provenance.go` 全文无 workspace/树 符号；⑤ `BudgetsFor(ctxWindow)`（`budgets.go:85`）＋`loop.go:222`；② 产物目录＝`rt.spec.dataDir/artifacts`（`run.go:1013/1109`，装配期定、不在工作区根下）。件＝`logs/13-ac6-spotcheck.txt`；交件件 `03`（行号缺文件名，见 §5 小缺陷）。
- **AC#7 成立**（零冻结件命中：`PLAN.md`／`docs/specs`／`allowlist`／`thresholds`／golden／`internal/risk`／`frontend`／`design` 全不在 `341ca42e` 路径表；`workspace.go` **零字节改动**＝强于"只许追加"；`composer_dispatch.go` 改动经逐行核＝纯注释 3 处；`git.go` 除注释外仅 const 一行）。件＝`logs/11-file-scope.txt`。
- **AC#8 附条件成立**（交件自落 `04a`/`04b`；本腿复跑 panel 整包＝恰 5 枚已知红、无 186 新红。⚠ 条件：① 票尾指针 `logs/04-gates.txt` **不存在**，真实读数在 `05-gate-clauses.txt`＋`06-final-cmdwisp.txt`，且两枚都在工作树 **未提交**（`??`）；② `04b` 里 `TestAC2SealNoticeLandsInTheRunLegLogFile` 红＝`GOTMPDIR=C:\tmp` 环境形（`06` 在 `%TEMP%` 下绿；我 `C:/tmp` 单跑 3/3 红，因 `.wisp-config-*.tmp` 的 winsec 密封记录多一条），非 186 因果链；③ `05` 的腿名册另有 `G1b` 声明漂移（`subagent_197.go:591`，票 197 的面）、`06` 另有 `TestTicket223ModeLoosening…` 红＋158.8s 死亡（票 314 死区）——均未在已知红名册具名。）；`d22scan` 半格未落任何读数（见 §5"没跑"）。
- **AC#9 成立**（`docs/evidence/s1/186-panel-ask.md`＝7 控件＋逐字信封＋4 口径＋4 具名缺口；交件零 `frontend/**`、零 `design/**`；现跑：页面 4 条请求尺＝4（`panel.ts:246/255/272/300`）、白名单 4 枚、页面死接线复认（`git-branch.tsx:59` `return null`、`App.tsx:216-219` 硬写空视图＋空回调））。

## 3. 两枚具名欠账按票面 AC 判

1. **"新建工作树无方法"**：**不构成本票 AC 缺口**（票面交付 3 只要求"列出并切到已有"的工作树/本地/分支；AC#2 明写"若必须新增白名单枚⇒具名报回、编排者落 `A##` 后再写"，即"先报不做"就是合规形）。**附条件**：AC#2 原句要求"具名报回要加哪一枚"，交件只写"需新 C17 名"、**未给拟用名字**；且"跑它会咬 `d22scan` ban #1/#4 的管线"表述不精确（九条 ban 无一条直管"起外部进程"；真跑会被 #1/#4 扫到它周围的 goroutine/超时写法）。建议编排者落 `A##` 时把名字与门归属一并登记。
2. **"成功回执＋常驻腿快照推送＝H10 最后一跳未做"**：**不属本票 AC**（H10＝票 33/35 的出向最后一跳；AC#9 明说页面另一半归别的会话；本票义务＝读面＋请求腿＋快照字段＋同源）。**附条件**：端到端"点击后立刻回显"在本票无法验收（常驻腿无快照推送、router 成功返回空串），ask 表 §4.1/§4.2 已具名；本票的"如实"义务（同源字段＋审计行）已落。

## 4. 具名残差 / 名册缺口（不驳回 AC，供编排者排票）

- **M6 射程**：`TestTheWorkspaceSwitchLegCannotMutateARepo` 只扫 3 枚具名文件（`workspace.go`/`workspace_handlers.go`/`git.go`）；包内**新增第四枚产码文件**可带 `exec.Command` 而不红。建议该腿扩面或改成"扫包内全部非 test 文件"。
- **M7 漏网**：`repoMutatingCallRe` 不含 `os.CreateTemp`/`os.Symlink`/`os.Truncate` 一族（`os.Create` 的 `\b` 使其不匹配 `CreateTemp`）；已覆盖的主形（spawn/`WriteFile`/`OpenFile`/`Remove*`/`Rename`/`Mkdir*`/`Chmod`）够钉住 09-27 事故机理，建议补词表。
- **守卫读 3 对注释敏感**（M5：新增注释里出现 `.SetWorkspaceRoot(` 即红）——保守向，如实登记。
- **交件件 `05`…**：真实名是 `05-gate-clauses.txt`/`06-final-cmdwisp.txt`，票尾误写为 `logs/04-gates.txt`；两枚未提交。
- **交件件 `03` 可读性**：`provenance.go`/`bridge.go` 等读数行号**缺文件名**（如"③ 650: …"），复核需自己补尺（本腿已代跑，见 `logs/13`）。
- **`workspace_handlers.go:52` 措辞**："ClearWorkspace is internal"不准确——它是导出方法、仅"零生产调用方"（我现跑＝只有注释与测试命中）；且它只把收窄退回配置根（仍在 allowlist 内），"不许放宽"的结构性结论不受影响。
- **判据外的两条非 186 红**：`04b` 的 AC2SealNotice（环境形）、`06` 的 Ticket223（未具名）＋314 死区死亡；`05` 的 G1b 声明漂移。归账不属本票，但"已知红名册"应补全。

## 5. 哪格没牙 / 没跑（具名）

- 没牙（具名残差）：M6（新文件不在射程）、M7（`CreateTemp` 一族漏网）。其余被攻仪器均"弄坏必红"（M1–M5）。
- 没跑：`sh scripts/d22scan.sh` 与 `bash .scratch/wisp/probes/154/gate-clauses.sh` 本腿**未复跑**（交件 `05` 有后者的锚定读数，但未提交；前者全仓零读数）；`./cmd/wisp/` 整包本腿未复跑（只单跑一枚定性红；交件 `06` 有未提交终读）；真窗端到端（H2/H3/H10 未落，无页面可看）；常驻腿真机 boot 未拉起（无头）。

## 6. 翻框建议

建议翻：**AC#1、AC#2、AC#3、AC#4、AC#5、AC#6、AC#7、AC#9**（成立）。
**AC#8 建议先补条件再翻**：把 `05`/`06` 入库（或编排者接受本表定性）＋把 AC2SealNotice/Ticket223/G1b 三条红具名归账（均非 186）。落地腿的"两枚具名欠账"经本表判＝一枚属"先报不做"合规形（欠名字的具名）、一枚不在本票 AC。
⛔ 本表未动任何 AC 框、未动 Status。
