# 181-v3 对抗验收证据（票 181 AC#7；非实现者腿）

- 腿名：181-v3（裁决者 ≠ 实现者，`SPEC-12 §4.3` #1/#3、D22 双角色）。
- 起手锚（`00-anchor.md`，commit `bf185665`）：2026-10-08 18:36:43 +0800；HEAD `0e9e00c3da07830435e4b8e21282c724e750866f`（A728 落账）。
- 判语：**AC#7 成立但附条件**（条件逐条在 §4；两轴分开写，⛔ 不压成一句"已成立"）。
- 射程声明：本程零产码改动（除种突变那一行、逐发还原）、零翻框、零 `-done`、零 push、票面/台账一字未动。`frontend/**` 与 `design/**` 未读未引。三枚冻结件未读未碰（射程判断，非内容引用）。

## 1. 两轴各自的判据（先各一句话）

- **渲染层（成立）**：改写账户沿真生产链 `cmd/wisp/panel_pump.go:87`（`panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())`）→ `WorkspaceViewFromRoot`（`internal/panel/workspace.go:85`，`:91` 逐字 `Reparse: res.Reparse, Rewritten: res.Rewritten,`）→ 线上 `rewritten` 键真被填；本腿三发突变各自打死一种坏填法（M1/M2/M3 全红），"字段恒 false"那半被推翻。
- **真进程层（今天仍不可达）**：活进程里唯一记账点＝`internal/panel/workspace.go:122` 的 `scope.SetWorkspaceRoot(actable, res)`，其前一道 `res.Actable()`（`:120`）对任何改写先拒（`internal/risk/pathresolver.go:93-99`）；`RequestWorkspaceSwitch`（定义 `:111`）在 HEAD 上**非测试零调用者**；面板侧 socket 处理器归票 186（`internal/panel/composer_dispatch.go:73-77`，接口 `:77` 的派发点 `:186` 无生产实现装配）。⇒ 该字段在真实运行进程里今天**恒 false 仍然成立**。

## 2. 起手复跑的现量（我跑我自己的）

1. 被点名的测试**存在**：`TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll` 定义在 `internal/panel/instructions_200_test.go:60`（函数名本程现取）。
2. 基线（未突变，`-count=1`）：`internal/panel` 焦点套 `-run 'TestWorkspaceViewCarriesEveryHalfOfTheAccount|TestWorkspaceAccountReachesEveryConsumer|TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll|TestScopeThatChangesItsAccountIsSurfaced'` ⇒ **rc=0**、11 枚 `=== RUN`、4 枚全 PASS（日志 `/d/tmp/181-v3/baseline-panel.txt`）；`internal/tools` `-run 'TestWorkspaceRoot|TestSetWorkspaceRootRefusesAnAccountAboutAnotherTree|TestClearWorkspaceDropsTheAccountWithTheRoot'` ⇒ **rc=0**、6 枚 `=== RUN`、4 枚全 PASS（`/d/tmp/181-v3/baseline-tools.txt`）。
3. `16901acb` 是 HEAD 祖先（`git merge-base --is-ancestor` rc=0）；`git show --stat --format='' 16901acb -- internal/risk/` ⇒ **空且 rc=0**（对 `internal/risk/` 零 diff）；其 stat 件＝`internal/panel/{git.go,git_test.go,instructions_200.go,workspace.go,workspace_account_181r3_test.go,workspace_test.go}`＋`internal/tools/{paths.go,paths_workspace.go,paths_workspace_account_181r3_test.go,paths_workspace_test.go}`。
4. 真进程层静态尺（本程现跑）：
   - `grep -rn --include='*.go' 'RequestWorkspaceSwitch(' internal/ cmd/ | grep -v _test.go` ⇒ 仅 `internal/panel/workspace.go:111`（定义自身），**非测试调用者 0 枚**。
   - `SetWorkspaceRoot` 非测试调用点 ⇒ 仅 `internal/panel/workspace.go:122`（接口声明 `:54`、实现 `internal/tools/paths_workspace.go:120`）。
   - `ResolveWorkspace` 非测试调用者 ⇒ 仅 `internal/panel/workspace.go:116`（即在零调用者的 `RequestWorkspaceSwitch` 内）。
   - `workspaceAccount` 写点 ⇒ 仅 `internal/tools/paths_workspace.go:135`（`SetWorkspaceRoot`）／`:146`（`ClearWorkspace`）。
   - `Rewritten:` 填充点（非测试全仓）⇒ `internal/panel/workspace.go:91`、`:141`（下文 §5 M4 焦点）＋`internal/risk/pathresolver.go:116`（账的来源，非视图）。

## 3. 三发突变（每发四件套；四件＝种前 hash／盘上改行复量／红句 file:line／还原 hash 等值＋scoped porcelain 空）

突变一律用 `sed -i '<行>s/…/…/'` 打在盘上，改后 `sed -n '<行>p'` 复量；测试带 `-count=1`、带 `PATH="$PWD/../../third_party/sherpa-onnx:$PWD/../../third_party/onnxruntime:$PATH"`（无注 PATH 会 `0xc0000135`＋0 条 `=== RUN`，本程两种情形都真实出现，作废的那次不充当读数）。

### M1 生产者停止填真值（view 构造器：`Rewritten: res.Rewritten,` → `Rewritten: false,`）

- 落点：`internal/panel/workspace.go:91`（`WorkspaceViewFromRoot` 内）。改后盘上复量：`\t\tReparse: res.Reparse, Rewritten: false,`。
- 种前 hash＝`4db510ce5594e24e31c2f0ef275fff0e6f8bf4bc`；还原后 hash＝同值（**RESTORED_EQUAL**）；`git status --porcelain -- internal/panel/workspace.go` 空。
- 测试 `-run 'TestWorkspaceViewCarriesEveryHalfOfTheAccount|TestWorkspaceAccountReachesEveryConsumer'` ⇒ **rc=1**（`/d/tmp/181-v3/m1-panel.txt`），红句逐字：
  - `workspace_account_181r3_test.go:83: rewritten = false for the account C26 marked rewritten ({Set:true Spelling:%WISP181R3ACC%\proj Canonical:C:\elsewhere\proj Reparse:false Rewritten:false Reason:已收窄到该工作区，但这条路径被 C26 的展开步改写过（env）：下面的坐标 C:\elsewhere\proj 是改写后的树，不是使用者命名的那串 %WISP181R3ACC%\proj，范围外的路径按 R2 判定}): the producer dropped the book again, and every consumer of this field is decoration`
  - `workspace_account_181r3_test.go:150: the view built from a scope that reports a rewritten account is not rewritten: {Set:true Spelling:%WISP181R3ACC% Canonical:C:\Users\…\AppData\Local\Temp\TestWorkspaceAccountReachesEveryConsumer2122847867\001\accounted tree …Rewritten:false…}`（临时目录路径中段以 … 缩，断言文本逐字）
  - `--- FAIL: TestWorkspaceViewCarriesEveryHalfOfTheAccount` ；`--- FAIL: TestWorkspaceAccountReachesEveryConsumer`。

### M2 把真值填成常量（`Rewritten: res.Rewritten,` → `Rewritten: true,`）

- 落点同 `:91`。改后复量：`\t\tReparse: res.Reparse, Rewritten: true,`。
- 种前 hash＝`4db510ce…`；还原后等值（RESTORED_EQUAL）；scoped porcelain 空。
- 同套 `-run` ⇒ **rc=1**（`/d/tmp/181-v3/m2-panel.txt`），红句逐字：
  - `workspace_account_181r3_test.go:66: a clean account rendered as an accounted one: {Set:true Spelling:D:\work\Wisp Canonical:D:\work\Wisp Reparse:false Rewritten:true Reason:已收窄到该工作区：范围外的路径按 R2 判定}`
  - `workspace_account_181r3_test.go:107: a reparse traversal rendered as a rewrite: {Set:true Spelling:D:\work\link Canonical:D:\work\link Reparse:true Rewritten:true …联接点…}`
  - `--- FAIL: TestWorkspaceViewCarriesEveryHalfOfTheAccount`；⚠ 端到端那枚 `TestWorkspaceAccountReachesEveryConsumer` **保持 PASS**——正例不因常量 true 而红，两向区分成立。

### M3 真生产者 WorkspaceRoot() 返回前焊常量 false（`return acc` → `acc.Rewritten = false; return acc`）

- 落点：`internal/tools/paths_workspace.go:74`（`WorkspaceRoot()` 返回那处）。改后复量：`\tacc.Rewritten = false; return acc`。
- 种前 hash＝`7e75e36f7abceebfd821969beeaf709f2992c02d`；还原后等值（RESTORED_EQUAL）；`git status --porcelain -- internal/tools/paths_workspace.go` 空。
- 测试 `-run 'TestWorkspaceRoot|TestClearWorkspaceDropsTheAccountWithTheRoot'` ⇒ **rc=1**（`/d/tmp/181-v3/m3-tools.txt`），红句逐字：
  - `paths_workspace_account_181r3_test.go:133: rewritten = false after recording C26's own rewritten Result ({Canonical:C:\Users\…\AppData\Local\Temp\TestWorkspaceRootReportsARealRewrittenAccountRecordedByC263205244231\001\expanded Reparse:false Resolved:true Spelling:%WISP181R3ACC% Rewritten:false Rewrites:[env]}): the producer is filling the field with a constant again, and 181-r2's consumer branch is dead code`
  - `--- FAIL: TestWorkspaceRootReportsARealRewrittenAccountRecordedByC26`。⚠ 与台账 `A618` 转述的它自己 M1 红句同枚断言文本——本行是**本腿自己跑的**（同句在我这份 rc=1 日志里现取）。
- 说明：M3 打的是真生产者的源头；`TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing`（清洁用例）在 M3 下保持 PASS＝该发是单向焊死，另一向由它守住。

## 4. 判语与条件（三选一：成立但附条件）

- **判语＝AC#7 成立但附条件**。AC#7 要求的"区分『填了账户』与『字段恒 false』"已由本非实现者腿独立复现（M1/M2/M3 三发各自红、红句在上；基线未突变全绿）⇒ **不是拿"编译通过"当证明**，且不是采信实现者自跑。
- 条件①（票面既有、非本修复引入；本程静态证据 §1/§2.4 全）：**真进程可达性仍不成立**——活进程唯一记账点被 `Actable()` 闸挡死、`RequestWorkspaceSwitch` 零生产调用者、面板 handler 归票 186 ⇒ 该字段在运行进程里今天恒 false；要真出现 true，只能等票 186 的 workspace handler 落地，或人工批准改 `internal/risk` 判级语义（禁区）。⛔ 两轴不许混写成"已修完"。
- 条件②（本腿新发现，§5）：`internal/panel/workspace.go:141` 成功路径填充点**今天没有任何仪器**（M4 不红）。

## 5. 新恒真面（M4：某发没红，具名）

- 落点：`internal/panel/workspace.go:141`（`RequestWorkspaceSwitch` 成功路径的 `Rewritten: res.Rewritten,`）。种突变＝改为 `Rewritten: true,`（盘上复量过）。
- 种前 hash＝`4db510ce…`；两次种（panel 一轮、cmd/wisp 一轮）后均还原等值、scoped porcelain 空。
- 读数：**rc=0、全绿、无一红**——
  - panel 焦点 7 枚（`TestWorkspaceSwitch*` 六枚＋`TestWorkspaceViewRenderedForSnapshots`＋`TestScopeThatChangesItsAccountIsSurfaced` 中的七枚，日志 `/d/tmp/181-v3/m4-panel.txt`）全 PASS；
  - cmd/wisp 两处调用点测试（`TestSnapshotWorkspaceSectionReportsTheNarrowing`／`TestRunPacketCarriesTheLoadedInstructionFiles`，日志 `/d/tmp/181-v3/m4-cmd.txt`）全 PASS。
- 结论（具名）：**该支今天没有仪器**。`RequestWorkspaceSwitch` 成功路径的该字段因 `Actable()` 闸而结构上恒 false（`res.Rewritten==true` 根本到不了 `:139`，`:144-146` 那枚 `if res.Rewritten` 文案分支同为死枝），把该填充点改成常量 true（一句成功却自称"篡改写"的假话）没有任何用例能抓到。⇒ 与 `:91`（两向都有牙）不同，`:141` 属单侧无牙面；票 186 若把该路放开携带改写账户，需先补这段仪器。
- 同发旁证：`Rewritten:` 非测试全仓填充点只有 `workspace.go:91`／`:141` 两处（§2.4），其余为账来源 `risk/pathresolver.go:116` 与读取 `risk/syncdirs.go:143`。

## 6. 本腿没试出来的（⛔ 不许省）

1. **没做**"活进程黑盒证不可达"（如起真 `wisp run` 试图以改写拼写录工作区）——真进程层读数是静态尺（调用者/闸门/派发点），不是行为实验。
2. **没跑**整包 `./internal/panel/`（在册红：本程两枚 flood/dropped＋既有四枚；按派单禁令不当读数、也未为其放宽任何断言）；`./internal/risk/`／`d22scan`／`gate-clauses` 未跑（本程零产码终态、无门禁增量；不把它们充当绿）。
3. **没跑**整包 `./cmd/wisp/`（只跑了上面具名两枚）；未核 `181-r3` 四笔 commit 的完整族谱（只核了 `16901acb` 祖先性、空 risk diff、件清单）。
4. **没验证** M4 恒真面在其它包（若未来别家测试断言该 view）是否会有牙——今天全仓 `Rewritten:` 填充点已枚举，未见第三处。
5. 未碰：票面、台账、AC 框、`-done`、push、`frontend/**`、`design/**`、三枚冻结件；工作树里他腿在飞的 `probes/181/r3/msg-*.txt` 等 `??`/` M` 项一枚未动。

## 7. 本腿 commit 清单（只 commit、不 push）

- `bf185665` 181-v3 起手锚（`probes/181/v3/00-anchor.md`）。
- `<本件>` 181-v3 证据件（`probes/181/v3/evidence.md`）。
