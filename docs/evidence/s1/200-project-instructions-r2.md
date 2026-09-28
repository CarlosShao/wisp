# 200-r2 — 改写账户被新消费者打红／面板那枚键没人填／一枚格式：交件与读数

写码程：`200-r2`（写腿）。工单：`.scratch/wisp/issues/200-nothing-reads-the-projects-instructions-for-ai-although-all-five-harnesses-do.md`。
派单：`.scratch/wisp/dispatches/2026-09-28-194x-impl-200r2-rewrite-account-and-instructions-carrier.md`。
上一腿证据件：`docs/evidence/s1/200-project-instructions-r1.md`。
**裁决与缺口审计必须由另一个 agent 做**（`SPEC-12 §4.3` #1/#3、D22 双角色）；本件是实现者自述，不是裁决表。

代码 commit：**`a90677e9`**（9 枚 pathspec，`git commit` 带显式路径，无 `add -A`/`.`，无 push）。
本件与票面那一跳的更正另作一枚 commit。

---

## 0. 起手锚点与名册（自取，不信派单行号）

- 起手 `git rev-parse --short HEAD` = **`d1c440ad`**；`date` = `Mon Sep 28 19:40:03 CST 2026`；交件时刻 2026-09-28 20:1x +08。
- 起手期间 `dev` 上另两枚 commit 落了盘（`7ea14ce3`、`64c1eea0`，都是 197-r1b），**我没碰 `internal/tools/**` 一字**，
  也没把它们的件带进我的 commit（我的 9 枚 pathspec 逐名在 §5）。

### 0.1 撞钉预检（派单 §0.2 那条命令原样，`PATH=$PWD/third_party/sherpa-onnx:$PWD/build:$PATH`）

`go test -count=1 ./internal/risk/ ./internal/panel/ ./internal/agent/ ./internal/projctx/ ./cmd/wisp/`
→ 逐字读数在 `.scratch/wisp/probes/200/r2/preflight-gate.txt`（19:44:44 起）：

| 包 | 起手读数 |
|---|---|
| `internal/risk` | **FAIL 2 枚**：`TestResolvePerCallBudget`（`1139468 ns/op = 1.139 ms/op (budget 1.000 ms, 1116 samples)`）＋ `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（红句点名 `..\..\cmd\wisp\run.go:686`）|
| `internal/panel` | **FAIL 4 枚**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree` |
| `internal/agent` | ok 3.852s |
| `internal/projctx` | ok 0.508s |
| `cmd/wisp` | **ok 104.063s**（DLL 路径带上了，不是 `0xc0000135`，用例真跑了）|

⚠ **与派单 §0.2 那点名册差一枚**：派单写 `internal/risk` 起手只有 `TestC26…` 一枚红，我这发起手命令里
`TestResolvePerCallBudget` **同时是红的**（同一发、同一命令、19:44:44）。处置见 §2(c)。

### 0.2 起手绿名册（"今天绿的用例名"，逐名存档）

- 四包 `-v`（`internal/risk`/`internal/panel`/`internal/projctx`/`internal/agent`，19:47）：
  **PASS 297 枚／FAIL 5 枚／SKIP 1 枚**，逐名清单 `.scratch/wisp/probes/200/r2/start-names-4pkgs.txt`（303 行）。
  注意：这一发里 `TestResolvePerCallBudget` 是 **PASS** —— 它与 §0.1 那发的差别只有并发包数与机器负载，这条事实记进 §2(c)。
- `cmd/wisp` 单包 `-v`：**PASS 89 枚／FAIL 0 枚**，逐名 `start-cmdwisp-v.txt.names.txt`（原始 `-v` 输出 `start-cmdwisp-v.txt` 留在 bench，未入库）。
  **入库的本程读数件逐名**：`preflight-gate.txt`、`final-gate-5pkgs.txt`、`final-gate-committed.txt`、
  `start-names-4pkgs.txt`、`final-names-4pkgs.txt`、`start-cmdwisp-v.txt.names.txt`、`final-cmdwisp-v.txt.names.txt`、
  `gofumpt-tracked-baseline.txt`、`gofumpt-tracked-after.txt`、`budget-test-readings.txt`、
  `cmdwisp-instructions-readings.txt`、`m1-static-ruler-blind.txt`、`build-vet-final.txt`、`adjacent-pkgs.txt`、
  `d22scan-final.txt`（都在 `.scratch/wisp/probes/200/r2/`）；
  两枚原始 `-v` 逐包转储（`start-v-roster.txt`／`final-v-roster.txt`／`start-cmdwisp-v.txt`／`final-cmdwisp-v.txt`）
  体积 37KB–300KB，**未入库**，从上面的逐名清单可以复跑出同名结果。
- 派单点名的在册常红三枚（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／
  `TestC21DesignTokensFourWayAgree`）起手即在红名册里，逐名相符；那枚"只有一跳、在 `frontend/src/lib/panel.ts`"的红
  `TestApprovalCardViewJSONKeysMatchFrontendTypes` 也在。**四枚我一个都没修、没 Skip、没放宽、没绕**（见 §2(b) 末与 §6）。

---

## 1. 落点逐跳（改后文件的 `file:line`，现读）

| 跳 | 位置 |
|---|---|
| 改写账户的消费（语义修法本体） | `internal/panel/instructions_200.go:116` `ProjectInstructionLoadRequestFor(ws WorkspaceView, globalDir string)`（`:117` 读 `ws.Rewritten`，改写 ⇒ `WorkspaceDir`/`GlobalDir` 都不给、只给一句拒读理由；`:119-124` 那句理由逐字带 `rewritten=true`＋两枚拼法；请求类型本体 `:85-93`）|
| 装配根不再 reach past 账户 | `cmd/wisp/run.go:706`（一次调用拿两枚目录＋理由）、`:707-709`（`rt.auditf` 打"一份都没读"那一行）、`:715-727`（`projctx.Options` 带 `Refused`）、`:728`（`rt.setInstructionLoader`）。**`cmd/wisp/run.go` 里 `.Canonical` 命中数 1 → 0**（`grep -c '\.Canonical' cmd/wisp/run.go`＝0，`grep -c 'Actable(\|\.Rewritten'` 也＝0：这个文件今天不读账户也不读权威路径，因为它根本不再拿路径）|
| 加载器的拒读分支 | `internal/projctx/projctx.go:182`（`Options.Refused`）、`:247-255`（`load()` 早退，优先于 `Enabled`）、`:108-113`＋`:125`（`SkipReason` 名册与字段）|
| 载体形状（五枚状态） | `internal/panel/instructions_200.go:25-41`（五枚状态名）、`:71-82`（`InstructionsSection{status,reason,files}`）、`:158-199`（`InstructionsSectionFromBundle`，含未知跳过名的 default 支）、`:201-213`（`WithInstructions` 保留且改成带状态的形，仍复制切片不 alias）；字段落点 `internal/panel/composer.go:74` |
| 泵的第一跳（读者） | `internal/panel/pump.go:143`（`PumpSources.Instructions func() *projctx.Bundle`）、`:232-241`（读者存在 ⇒ 必发键；不存在 ⇒ 仍四枚键）|
| 泵的第二跳（生产装配） | `cmd/wisp/run.go:475`（`Instructions: rt.instructionBundle`）、`cmd/wisp/panel_pump.go:93-101`（`setInstructionLoader`，`instrMu` 守）、`:108-119`（`instructionBundle` 读加载器的 live bundle，不塑形）|
| 运行时字段 | `cmd/wisp/run.go:254-262`（`instrMu` / `instrLoader`）|
| 判据 | `internal/panel/instructions_200_test.go:60/142/235/287` 四枚；`cmd/wisp/instructions_200r2_test.go:111/186` 两枚端到端 |

派单给的行号复算：`internal/panel/workspace.go:104-111`（构造 view 那段）与 `84feec44:internal/panel/composer.go:68`
（`Instructions []ProjectInstructionFile json:"instructions,omitempty"`）**都对**；`cmd/wisp/run.go:686` 那枚红句行号也对。
本程行号一律现读，未照抄。

---

## 2. 三格：改前 → 改后读数

### (a) 票 102 的在册判据：被打红 ⇒ 全绿，且修的是语义不是尺

- **未修码的读数（"哪一枚 commit 算未修码"的直接答复）**：HEAD `d1c440ad` 里含 `84feec44`，起手那一发
  `internal/risk` 的 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` **FAIL**，红句逐字：
  `..\..\cmd\wisp\run.go:686 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer`
  （`pathresolver_rewrite_account_test.go:200`；存档 `preflight-gate.txt:5-6`）。
  ⇒ **未修码＝`84feec44` 那一枚 commit**（派单的因果对照我复算相符：该 commit 之前 run.go 里 `.Canonical` 0 命中、之后 1 命中）。
- **改后**：同一枚尺 `go test -count=1 -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg ./internal/risk/`
  → `--- PASS (0.04s)`（19:53:55）；全包 `ok`（终态两轮 7.215s / 6.270s，§4）。
- **为什么这不是"塞一句 `.Rewritten` 过尺"**：
  1. 落点在**持视图的包**里，不在消费者手里摸字段——照本仓已有的两枚先例：`cmd/wisp/panel_pump.go:129-131`
     明写 "It is handed over as the whole view, not as its path string"，
     `internal/panel/git.go:240-245` 的 `ReadGitForWorkspace` 就是同一形状；
  2. **算而且 branch**：`ws.Rewritten` 为真 ⇒ 工作区**与数据目录都不给**（`GlobalDir` 也是空），
     `projctx.load()` 在新加的 `Refused` 早退分支上一份文件都不打开；
  3. **响亮**：装配根一行 `rt.auditf`（`run.go:707-709`）＋每轮加载器 `Manifest()` 把同一句打到 `projctx:` 日志；
  4. **行为正控真跑**（不是词面）：`TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll`
     （`internal/panel/instructions_200_test.go:60`）造两枚真目录、各放一份真 `AGENTS.md`，
     喂 `Rewritten=true` 的视图 ⇒ 断 `req.WorkspaceDir==""` 且 `req.GlobalDir==""` 且真加载器
     `Block==""`／`Files==0`／`SkipReason==rewrite_refused`／日志含"一份都没有读"与"rewritten=true"；
     **控制组**同一接线只把那枚字段换成 `false` ⇒ 两份正文都进 `Block`、两枚路径逐字命中、`Files>=2`。
     终态 `-v` 读数：`--- PASS`（§4 名册）。
- **今天这一支在生产里可达吗**：**不可达**，我没有把它说成已闭合的活口子。唯一会给 `WorkspaceView.Rewritten=true` 的
  生产者是 `RequestWorkspaceSwitch`，而它在构造视图之前就用 `res.Actable()` 把改写拼法拒了
  （`internal/panel/workspace.go:85` → `:104-111` 不可达那支）。这与 `internal/panel/git.go:260-266` 的
  "MEASURED TODAY: this branch is unreachable" 是同一笔账，**生产者侧填真值＝票 181 AC#7（台账 task #162 仍 pending）**。
  我这格买到的不是"今天没人能踩"，而是"上游一松，这条路就已经是拒读而不是照读"。

### (b) 面板那枚 `instructions` 键：**选了甲**（把两跳接上＋端到端读数＋把状态分开）

- **改前现量**（派单 §2(b) 我复算全对）：`84feec44:internal/panel/instructions_200.go:43-46` 对 `nil`/空 bundle 一律
  `return nil`，`composer.go:68` 带 `omitempty` ⇒ **`off`、`开着但目录里没文件`、以及 (a) 那种"被改写账户拒读"三种真相
  marshal 成同一枚"键不存在"**；`internal/panel/pump.go` 里 `grep Instructions` 零命中 ⇒ 载体在、装配侧无人喂。
- **为什么不选乙**：乙把那格当场退回，代价是票 200 AC#7 的"面板可见"一半作废；而键已经进了 C17 快照契约
  （票面 `A273`/`A389` 的射程只到"新增字段"），摘掉它同样是契约面移动，且界面那一跳已经写进票面派给 owner。
  甲还有一处乙没有的读数：**功能开没开这件事第一次能从生产 packet 里读出来**（下面两条 e2e）。
- **改后形状**（形状我自己定，派单 §2(b) 授权了；键数不变，仍是一枚 `instructions`）：
  `instructions: { "status": 五枚名之一, "reason": 非 loaded 时逐字带原因, "files": [ …八个字段逐字未变… ] }`，
  五枚状态＝`loaded`／`none_found`／`off`／`refused`／`not_run`；`files` 空时是 `[]` 不是 `null`（渲染侧不必摸 nil）。
  `omitempty` 留着，但它现在**只**表示"这台宿主没有读者"这一枚诚实缺席。
- **端到端读数（`cmd/wisp` 现成 harness，真装配、真 publish、断的是 bytes 不是结构体）**：
  - `TestRunPacketCarriesTheLoadedInstructionFiles`（`instructions_200r2_test.go:111`）：真工作区收窄走
    `panel.RequestWorkspaceSwitch`（生产处理者，不是测试拼的路径），树里放真 `AGENTS.md` ⇒
    `status=loaded reason= files=[instr-ws/AGENTS.md tier=project depth=0 bytes=269]`，逐字段钉
    `tier=project`／`depth=0`／`bytes>0`／`source=="fs.read"`（现成名，未新造）；
    并钉正文 `E200R2HEAD` **不在** wire bytes 里（AC#4 的分层搬到面板线上）；stderr 里那句 `projctx:` 打印也在。
    终态 `-v`：`--- PASS (1.94s)`。
  - `TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound`（`:186`）：`[agent] project_instructions_enabled = false`
    那发 ⇒ `status=off reason=已按你的配置跳过：agent.project_instructions_enabled=false，本轮一份项目说明都没有读（不是没找到，是被配置关掉的）。 files=[]`；
    **控制组**（同 fixture 形状、开关默认开）⇒ `status=loaded files=[004/AGENTS.md tier=global depth=-1 bytes=274]`；
    再钉两枚 packet 逐字节不同。终态 `-v`：`--- PASS (2.49s)`（第一发 5.43s，含控制组）。
  - 三态在 marshal 层面分开的常驻尺：`TestTheInstructionStatesDoNotShareOneWireShape`
    （`internal/panel/instructions_200_test.go:142`）把 `off`/`none_found`/`loaded`/`refused` 四枚真加载器的产物
    逐对比字节（五对全钉），并钉：`off` 必带那句配置跳过、`none_found` **不许**带它（它没被谁关掉）、
    `refused` 必带 `rewritten=true`、`loaded` 之外任何状态不许出现正文。终态 `-v`：`--- PASS`。
  - "键在、值恒缺"这一形按派单的要求被反向钉住：`TestAPumpWithoutAnInstructionsReaderSendsNoKey`（`:235`）
    前半——**没接读者的泵，marshal 出的顶层键仍只有那四枚**（`pump_test.go:123`/`:291` 两枚逐字钉的读数因此一字未变，
    我没动它们，它们仍绿：`TestThePumpBuildsThePacketFromWhatTheHostHolds` / `TestPublishHandsTheBytesToTheAttachedExit`
    在终态 `-v` 里 PASS）；后半——接了读者的泵**必发键**，`not_run` 与 `loaded` 是两枚不同状态而不是"一枚键的有无"。
- **附带那枚我改不了的（派单 §2(b) 末）**：`TestApprovalCardViewJSONKeysMatchFrontendTypes` **仍红，逐字同名册**，
  红句子一字未变：`approval_test.go:129: Go Snapshot emits [instructions] that interface PanelSnapshot does not declare`。
  我把载体从数组换成对象**不会**改变这把尺的读数——它只按 `json` tag 收顶层键名
  （`internal/panel/approval_test.go:174-191` 的 `jsonKeysOf` 不递归）。我**没改尺、没加豁免、没让 Go 少发一枚键来过尺**。
  转绿只有一跳（`frontend/src/lib/panel.ts` 的 `PanelSnapshot`），`frontend/**` 对我是零写面；
  **票面那一跳的文字因为形状变化过期了，我已就地更正**（`.scratch/wisp/issues/200-*.md` 那一节新增 8 行"200-r2 更正"，
  写清 `instructions?:` 现在是一枚对象、五枚 status 名、`files` 是数组、以及"泵没接读者的宿主仍不发键"）。
  `TestComposerContractTypesMatchFrontend` 的红句子同样逐字未变（`instructions` + `git currentModel modelKnown`）。
- **AC#7 现在的真实边界**（不粉饰）：Go 侧线上已经有值、有状态、有原因；**界面侧仍然一个字节都不读**（`panel.ts` 未声明、
  也没有 Go→页面通道，票 33/35 的最后一公里仍开）。所以"面板可见"这半格我只能报到"packet 里看得见"为止。

### (c) 两枚小账

- **gofumpt**：派单点名的那枚 `internal/agent/instructions_test.go` 已过尺。
  - 改前跟踪集读数（`git ls-files -z '*.go' | xargs -0 "$GOPATH/bin/gofumpt.exe" -l`）＝ **8 枚**：
    `internal\agent\instructions_test.go`、`internal\tools\subagent_197.go`、`internal\tools\subagent_197_test.go`
    ＋ 5 枚 `.scratch\wisp\probes\**` 刻意做坏的样本（逐名在 `gofumpt-tracked-baseline.txt`）。
  - 改后（`gofumpt-tracked-after.txt`，20:13:31）＝ **只剩那 5 枚 probes**：
    `probes/163/a1/main.go`、`probes/174/c2/zz174c2_wiring_pair_windows_test.go`、`probes/183/r2/mut/task-boxset-off.go`、
    `probes/185/c1/mut/fs_broken.go`（枚 gofumpt 自己的 parse 报错行）、`probes/185/r1/mut-m1/hostpath_185.go`；
    逐名可按 `attrib.sh` 的归因规则对上真票（`issues/163-*`、`issues/174-*`、`issues/183-*`、`issues/185-*` 四张都在）。
    **`internal/`、`cmd/` 下已无未净件**。
  - 我碰那枚的增量只有一个空行（`git diff` numstat `1/1`，无任何内容改动）：`internal/agent/instructions_test.go:262`。
  - ⚠ **`internal/tools/subagent_197*.go` 那两枚不是我清的**：它们在我起手后由 197-r1b 的 `7ea14ce3` 交件时格式化掉了。
    派单 §2(c) 写"归 197-r1b、你别动"——我确实没动（`git status --porcelain -- internal/tools` 在我 commit 前后都空，
    我的 9 枚 pathspec 里也没有 `internal/tools/**`）。这条账不该记到我头上，故在此具名。
  - `.scratch/wisp/probes/161/r5/attrib.sh --tracked-only` 这枚仪器在本机**跑不出读数**：它把
    `D:\work\base\gopath/bin/gofumpt.exe` 这种混合分隔符路径交给 `xargs -0`，gofumpt 起不来，脚本按 rc=123 判
    "tracked tree not readable by this ruler"（原文在 `attrib.sh:168` 那类守卫里）。**CI 那版（checkout、纯 POSIX 路径）
    是否同形我没测**；我改用同一枚尺手抄命令（`git ls-files -z '*.go' | xargs -0 gofumpt.exe -l`），
    即 CI 步骤 (A) 的原文形式。这枚仪器缺陷归票 161 的 instrument owner，不在我射程内，我只具名不修。
- **`TestResolvePerCallBudget`：我不复现，交回编排者**（派单 §2(c) 给的第二支），但带**改前那发红句的完整命令与读数**，
  以及我改后 9 发的读数表（存档 `budget-test-readings.txt`）：

  | 形状 | 发数 | 读数 |
  |---|---|---|
  | 起手五包并行（派单 §0.2 原命令，改前） | 1 | **FAIL** `1139468 ns/op = 1.139 ms/op (budget 1.000 ms, 1116 samples)` |
  | 四包 `-v`（无 cmd/wisp，改前） | 1 | PASS |
  | 单跑 `-run TestResolvePerCallBudget` | 3 | ok 1.127s / 1.326s / 2.393s（编排者那三发是 1.2/1.4/1.4，同形）|
  | 整包 `./internal/risk/` 单跑 | 3 | ok 4.176s / 2.856s / 3.970s |
  | 五包并行（改后两轮门禁） | 2 | ok 7.215s / 6.270s |
  | 五包并行（commit `a90677e9` 落盘后复跑，§4 第三发） | 1 | ok 4.931s |

  ⇒ **改后绿 9 发／改前红 1 发**（3＋3＋2＋1＝9，起手那发四包 `-v` 的绿在"改前"一侧，不计入）。

  ⇒ 结论按派单的规矩写死：**"今天应该红"这句话我不写**；我复现到的唯一一发红是起手那一发（命令与读数已在盘上），
  之后 9 发全绿。这枚尺是计时尺（`pathresolver_budget_norace_test.go:34`），
  对同机并发负载敏感；**修它要动阈值或动尺，两枚都是硬禁**（AGENTS.md §1.1"thresholds.go 一字节都不许动"），
  所以我把它交回编排者处置，并建议名册口径改成"起手读数逐名抄"，不是"19:3x 那一发抄"（§6 第 1 条）。

---

## 3. 两发变异／正控（都真跑，读数逐字）

| # | 变异（临时改产码，跑完即回滚，`grep "if false"` 现在零命中） | 读数 |
|---|---|---|
| **M1**（(a) 侧：算了却不 branch） | `internal/panel/instructions_200.go:116` 的 `if ws.Rewritten {` 改成 `if false && ws.Rewritten {` —— 正是派单 §2(a) 点名的假修形状 | `TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll` **FAIL**：`refusal lost for a rewritten view: {WorkspaceDir:…\001\moved GlobalDir:…\002\wispdata Refused:}`（两枚目录照旧发出去）；`TestTheInstructionStatesDoNotShareOneWireShape` **FAIL**：`refused: status = "none_found", want "refused"` ＋ `refused and none_found marshal to the same bytes`。**同一发里我把 `internal/risk` 那枚静态尺真跑了一遍**：`go test -count=1 -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg ./internal/risk/` → **`ok 0.319s`（绿）**，存档 `m1-static-ruler-blind.txt` ⇒ **尺看不见的那一形，行为用例看见了**；这也是我对派单 §2(a)"不许只骗过尺"的正面答复：光过尺的修法在这套用例下必红，而过尺本身抓不到它 |
| **M2**（(b) 侧：读者接了却没人消费） | `internal/panel/pump.go:236` 的 `if p.src.Instructions != nil {` 改成 `if false && …` | `TestAPumpWithoutAnInstructionsReaderSendsNoKey` **FAIL** 两处：`a wired reader with no bundle yet lost the not_run state`（bytes 里根本没有 `instructions`）＋ `the reader's bundle did not reach the wire`；**`cmd/wisp` 两枚端到端同时 FAIL**：`the packet carries no instructions key at all: {…}`（来自真装配的真 bytes，两发明文各带整枚 packet）。反向那一发（没接读者的 pump 仍只发四枚键）在同一用例前半钉住，且 `pump_test.go:123`/`:291` 两枚逐字钉在变异前后都绿 ⇒ 键的**出现**是被读者造成的，不是恒在 |

回滚自证：`grep -n "if false" internal/panel/pump.go internal/panel/instructions_200.go` → 零命中；
回滚后 `go build ./...`、`go vet ./...`、五包门禁全部重跑过（§4）。

---

## 4. 门禁终态逐包读数 ＋ 在册红名册逐名比对

终态五包 `-count=1` 连跑两轮（`final-gate-5pkgs.txt`）：

| 包 | 第 1 轮 | 第 2 轮 | 红名（两轮逐字相同） |
|---|---|---|---|
| `internal/risk` | **ok 7.215s** | **ok 6.270s** | 无（`TestC26…` 与 `TestResolvePerCallBudget` 都绿）|
| `internal/panel` | FAIL 4 | FAIL 4 | `TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree` |
| `internal/agent` | ok 2.854s | ok 2.839s | 无 |
| `internal/projctx` | ok 0.338s | ok 0.472s | 无 |
| `cmd/wisp` | ok 87.372s | ok 82.682s | 无 |

**第三发：commit `a90677e9` 落了盘之后，我对着盘上那份 bytes 又跑了一遍同一发命令**（`git status --short -- internal/ cmd/`
当时为空 ⇒ 工作树＝那枚 commit 的字节；读数存档 `final-gate-committed.txt`，20:2x）：
`internal/risk` **ok 4.931s**／`internal/panel` **FAIL 4（逐名同上，一字未变）**／`internal/agent` ok 2.207s／
`internal/projctx` ok 0.154s／`cmd/wisp` **ok 82.400s**。
⇒ 上面那两轮的读数不是"改完注释之前的读数"。

- 邻检（我改的包被别人 consume 的面）：`internal/tools` ok 16.531s、`internal/config` ok 1.455s、
  `internal/agent` ok 2.212s、`internal/llm` ok 41.127s（`adjacent-pkgs.txt`）。
- `tools/d22scan`（`sh scripts/d22scan.sh`）：**clean**，正控 `runtests.sh` PASS=34 FAIL=0 SKIP=0；
  扫描分母逐枚：bans #1-5 internal/=214、cmd/=24，ban #6 frontend/=85，ban #7 internal/tools/=22，
  ban #8 design/=39、frontend/=85、internal/=450（注释与 `_test.go` 都算）、cmd/=49；生产 .go 238 枚。
- `go build ./...` rc=0，`go vet ./...` rc=0（全模块，含测试编译）。

**名册逐名比对（差集，不是"应该没变"）**：

| 集合 | 起手 | 终态 | 差 |
|---|---|---|---|
| 四包 `-v` 顶层用例 | PASS 297／FAIL 5／SKIP 1 | **PASS 302／FAIL 4／SKIP 1** | 那枚 SKIP 起手与终态同一枚：`TestSyncRegistryProbeLive`（不是我 Skip 的，两态一致）。**少了的红**＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（转绿）；**多了的绿**＝我新写四枚 `TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll`／`TestTheInstructionStatesDoNotShareOneWireShape`／`TestAPumpWithoutAnInstructionsReaderSendsNoKey`／`TestWithInstructionsKeepsTheR1CarrierContract`。**`comm -23` 差集为空 ⇒ 没有一枚原本绿的变红**，也没有一枚断言被我放宽或删掉 |
| `cmd/wisp` `-v` | PASS 89／FAIL 0 | **PASS 93／FAIL 0** | 多的四枚：我的 `TestRunPacketCarriesTheLoadedInstructionFiles`／`TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound`，以及 197-r1b 在 `64c1eea0` 带来的 `TestSubagentStreamKeyBuildersAgreeAcrossBothPackages`／`TestSubagentStreamKeyPrefixAgreesAcrossBothPackages`（**不是我写的，逐名分开记账**）|
| 在册常红比对派单 §0.2 | 三枚常红＋一枚待界面 | 名册**逐字不变** | **没有少报一枚、没有多报一枚**（`TestResolvePerCallBudget` 见 §2(c) 单列，不进名册）|

---

## 5. 纪律自证（被拒枚数与写面）

- **被权限/纪律系统拒绝的调用：0 枚**。全程未 `git add -A`/`.`/`-a`、未 `--amend`/`reset`/`rebase`/`stash`/
  `checkout .`/`restore`/`clean`/`merge`/`worktree`、**未 push**、未在本仓目录内删除任何文件。
- 我的 commit `a90677e9` 的 9 枚 pathspec（逐名）：
  `internal/projctx/projctx.go`、`internal/panel/composer.go`、`internal/panel/instructions_200.go`、
  `internal/panel/instructions_200_test.go`（新建）、`internal/panel/pump.go`、`cmd/wisp/run.go`、
  `cmd/wisp/panel_pump.go`、`cmd/wisp/instructions_200r2_test.go`（新建）、`internal/agent/instructions_test.go`。
  合计 `9 files changed, 881 insertions(+), 33 deletions(-)`。
- `git diff --numstat`（提交前自证，删除列逐枚给理由）：
  `issues/200-*.md` **8/0**（纯追加）；`cmd/wisp/panel_pump.go` **31/0**；`cmd/wisp/run.go` **39/14**（14 枚删除逐枚是我
  替换掉的 `ws := rt.workspaceView()` 那 5 行旧块＋旧 `projctx.New` 内联块的 9 行注释/字段行，无一行内容被丢）；
  `internal/agent/instructions_test.go` **1/1**（一个空行）；`internal/panel/composer.go` **13/7**（注释块重写）；
  `internal/panel/instructions_200.go` **148/9**；`internal/panel/pump.go` **21/1**；`internal/projctx/projctx.go` **27/1**。
  删除列全部来自"把旧形状换成新形状"这一格本身，没有一处是删别人的活。
- **别人的脏文件一枚未动、一枚未提交**：`.gitignore`、`.scratch/wisp/probes/152/**`、
  `.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除、`docs/evidence/s1/152-*.md`、
  以及未跟踪的 `docs/reports/survey-2026-09-28-*.md` 两枚——提交后 `git status` 里它们仍在原样。
- **零写面遵守**：`frontend/**`／`design/**` 不写不读不转述（`internal/panel` 那两枚尺自己会读 `panel.ts`，我没打开过它）；
  `internal/tools/**` 零命中；`internal/panel/tokens_fourway_test.go`、`PLAN.md`、`docs/specs/**`、
  `thresholds.go`、golden、`allowlist.txt` 一字未动；本仓那枚根 `AGENTS.md` 只在 r1 当过样本，本程未读未写。
- 未引入任何新依赖：新增 import 只有仓内包（`cmd/wisp/panel_pump.go`＋`internal/projctx`、
  `internal/panel/pump.go`＋`internal/projctx`，方向仍是 panel→projctx→risk，无环；
  `internal/projctx/projctx_test.go` 那侧是外部测试包 import panel，本来就存在）。
- **契约面**：新增的只有 `Snapshot` 那一枚既有键的**内部形状**（数组→对象）与 `projctx` 的
  `Options.Refused`／`Bundle.SkipReason` 两个字段——`C1–C32` 与 `D1–D47` 我一枚都没改；
  C17 快照"新增字段"的批准射程（票面 `A273`/`A389`）我按 r1 的口径理解成"允许加字段"，
  **把已声明字段的类型换掉是不是同一枚授权，归编排者裁**（§6 第 5 条），我没有自行扩大到"可以改契约"。

---

## 6. 我没测什么（具体，≥3）

1. **改写账户那一支在生产里今天不可达**，我只证了"喂进去就不读"（§2(a) 末）：生产者侧填真值＝票 181 AC#7 仍 pending。
   所以本程对票 102 那枚尺的转绿，买到的是"结构性不再 reach past 账户＋上游一松就拒读"，
   **不是**"今天有一条活路被我堵上了"。
2. `cmd/wisp` 端到端**没有跑过 refused / not_run 这两枚状态**：refused 要造一枚真 `%VAR%` 拼法的工作区请求，
   而"面板 → `RequestWorkspaceSwitch`"那一路在本仓仍零调用者（`cmd/wisp/panel_pump.go:79-81` 明写）；
   not_run 在真 run 里我推断不出现（execute() 在第一次 publish 之前接好了加载器），**这是推断，我没量**。
   两枚状态只有 panel 层的真加载器用例钉着。
3. **子代理环（`task.spawn` 的孩子）不带项目说明**这件事我没修也没测，只现读记录：
   `internal/tools/subagent_197.go:266` 用 `agent.New(opt)` 起孩子，`opt` 来自
   `SubagentDeps.BaseOptions`（`cmd/wisp/run.go` 的 `rt.loopOpt`），而 `agent.Options` 里没有 instruction 源
   （它在 `Assembler` 上，只有 `AttachProjectInstructions` 会接，`internal/agent/prompt.go:280`）
   ⇒ 父环带项目规矩、子环不带。写面归 197 腿（`internal/tools/**` 对我是零写面），故只具名不修。
4. **界面那侧一个字节都没测**（`frontend/**` 零写面）：`PanelSnapshot` 仍没声明 `instructions`，那枚尺仍红；
   "面板真画出五枚状态"这句我拿不出证据，且 Go→页面通道本身还不存在（票 33/35）。
5. **多包并行下的计时飘红我无法排除**：`TestResolvePerCallBudget` 起手红过一次、之后 9 发绿，
   我没做"固定并发度×重复次数"的命中率测量（那要么改尺要么改阈值，两枚都禁）。
6. 大小写／卷根／UNC／`\\?\`／junction 这些路径形状我**一个都没新测**——继承 r1 证据件 §6 那笔未覆盖面，
   本程没有改进它。拒读那支在这些形状下的行为等价于"根本不交目录"，逻辑上更保守，但保守不等于测过。
7. **全局级 × 拒读**只证到"两枚目录都不给"（panel 用例），没证到"数据目录里那枚 `AGENTS.md` 在改写工作区下确实没被读进请求"
   ——那需要 wire-body 级断言（票 200 AC#1 那把尺在 `internal/agent`），我借用不了它，因为它是 agent 包的 harness、
   喂不进一枚改写过的视图。
8. `omitempty` 的缺席形我只在 panel 层与 reader-less 的 pump 上证过；
   **生产里"泵接好了但 execute() 之前"那一发 publish**（比如 boot 期若有 publish）我没读过它的 bytes。

---

## 7. 对派单的不服／更正（逐条，能给读数就给读数）

1. **起手名册少点名一枚**。派单 §0.2 写 `internal/risk`＝只有 `TestC26…` 一枚红（19:3x 现量），
   而我按那条命令原样跑的第一发（19:44:44）里 `TestResolvePerCallBudget` 同发为红，读数 1.139 ms/op。
   不是我在争"它也该算在册"，而是**名册口径**该定成"每腿起手自己复量并逐名抄"，否则每腿都会多出一枚没人认领的红。
2. **§2(a) 的落点我走成了更严的一旁**：派单写"在把 `ws.Canonical` 交给 projctx 之前先看 `ws.Rewritten`"，
   我实现成**run.go 不再持有 `.Canonical`**（`grep -c '\.Canonical' cmd/wisp/run.go`＝0），账户由持视图的包读，
   先例是 `internal/panel/git.go:240-245` 与本文件自家注释 `cmd/wisp/panel_pump.go:129-131`
   （"handed over as the whole view, not as its path string"）。
   行为、审计、正控三件都在（§2(a)、§3 M1）。**若编排者要求字面按派单落回 run.go，我按同一判据重落，不需要新证据**——
   只是我判断那一版更弱：它会在文件里留一处"读了权威路径"的事实。
3. **状态枚数我给了五枚，不是派单说的三枚**。多出的两枚不是我加的装饰：
   `refused` 若并进 `none_found`，本格要消灭的那形（两种真相同一枚读数）就在原地复活；
   `not_run` 若不命名，"读者接了但这一轮还没跑过加载"就只能靠 `nil` 猜，那与"没接读者"在 bytes 上又混成一枚。
4. **`internal/tools/subagent_197*.go` 那两枚未格式化件不是我清的**（`7ea14ce3` 是 197-r1b 的）。
   派单 §2(c) 已经写明归它、别我动，我照办了；这里只是把那格"完成"的归属钉回盘上事实。
5. **载体形状换了（数组→对象），使票面写给界面那支的那一行过期**，我按派单 §2(b)"形状你自己定"做了，
   并同批把票面那一节更正了（8 行纯追加）。这里请裁一句：`instructions` 是我 r1 那批"新增字段"授权
   （票面 `A273`/`A389`）之内的形状调整，还是**换掉已声明字段的类型**已越出那枚授权？
   我的立场：键数没变、红句子逐字未变、界面那一跳仍是一行声明 ⇒ 未越出；但这一枚该由编排者定，不该由实现者自证。
6. **`attrib.sh --tracked-only` 在本机跑不出读数**（`xargs -0` 拿到混合分隔符的 gofumpt 路径，rc=123，
   脚本自己判"tracked tree not readable"）。我没改那枚仪器（它是票 161 的交付物），
   改用同一把尺的原文形式 `git ls-files -z '*.go' | xargs -0 gofumpt.exe -l`。**CI 那版是否同形我没测**（与 §6 第 5 条同一枚未覆盖面）。
7. 行号：派单给的 `run.go:686`、`composer.go:68`、`instructions_200.go:41/65`、`workspace.go:104-111`
   在我这发复算**全部对上**，本轮无行号漂可报。
