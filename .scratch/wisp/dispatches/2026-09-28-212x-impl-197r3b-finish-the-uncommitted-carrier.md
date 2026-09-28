# 派单 2026-09-28 21:2x — 写腿 `197-r3b`：**接着半成品做完**（前一腿撞到 150 轮上限，代码全在工作树里未提交、证据件没写）

## 0. 现场（我 21:2x 现量，别重做已完成的部分）

- 工作树里未提交的改动＝**票 197 载体层的产码**：`cmd/wisp/panel_pump.go`、`cmd/wisp/run.go`、`internal/panel/composer.go`、
  `internal/panel/pump.go`、`internal/tools/subagent_197.go`（+120/−约 20）、`internal/tools/subagent_197_test.go`（+43），
  加一枚**未跟踪的新测试件** `cmd/wisp/subagent_carrier_197_test.go`。合计 **372 插入 / 33 删除**（6 枚文件）。
- 我的现跑读数：`internal/tools` **ok 16.471s**；`internal/panel` **红＝恰好起手那 4 枚在册**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／
  `TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）⇒ **它没打红新名册、也没改尺**。
  `cmd/wisp` 我起了后台跑但**读数还没回来**（我此刻盘上无可引的 `cmd/wisp` 终态数）⇒ **你自己跑它**，别信任何"应该已经绿"。
- **别人在飞/在脏的文件你不许碰也不许提交**：`docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md`（`M`，不归本票）、
  `.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`design/**` 那 16 枚删除、`.scratch/wisp/.scratch/**`。

## 1. 你的活（三件，按顺序，别扩范围）

1. **读懂那 372 行**（工单正文仍是 `.scratch/wisp/dispatches/2026-09-28-210x-impl-197r3-carrier-roster-and-streams-into-snapshot.md`，判据与禁区都在里面）。
   先自己判断它**离"名册与每枚子代理的流真进快照"还缺哪一跳**——**缺什么补什么，别推倒重写**。
2. **把判据钉成钉**：载体层要有**生产路径**读数（真装配 ⇒ 快照里出现子代理、带 `taskID`／D43 状态／父任务／它自己那条流），
   还要一条**拔掉泵必红**的反向正控（这条前一腿没来得及写）。**状态只许 D43 那 20 枚名字**；**取消不级联**；
   **不许新增 C17 `PanelBridge` 方法名**（要入向就停手具名报我）；**不许把东西加成无条件顶层键**（`pump_test.go:123`／`:291` 那两枚"四枚顶层键"钉今天没红，
   说明前一腿走的是嵌进现有键或 `omitempty`——**保持这个选择**，并把它写清是哪一种）。
3. **提交并写证据件**：一次 commit 带**显式 pathspec**（那 6 枚＋新测试件逐枚点名），
   证据件 `docs/evidence/s1/197-subagent-carrier-r3.md`，六节按上一份派单第 3 节。

## 2. 纪律（与上一份派单逐字一致，这里只重列最要紧的）

只 commit 绝不 push；禁 `git add -A`/`.`/`-a`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/`merge`/`worktree`；
仓内绝不删东西（写坏只准 `git cat-file blob HEAD:<path> > <path>`）；`frontend/**`、`design/**` 不写不读不转述；
`internal/panel/tokens_fourway_test.go`、`internal/tools/bridge.go`、`internal/agent/budgets.go`、`thresholds.go`、golden、`allowlist.txt`、`PLAN.md`、`docs/specs/**` 一字不动；
测试必带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（否则 `0xc0000135`＝一枚不跑）；`gofumpt` 在 `"$GOPATH/bin/gofumpt.exe"`，交件时 `internal/`＋`cmd/` 未净必须＝0；
中文长段用 Edit/Write，别用 shell heredoc；写跟踪文件后 `git diff --numstat` 自证删除列＝0。

**轮次纪律**（前一腿就是死在没管轮次）：**先 commit 产码，再写证据件**；证据件**分段落盘**（写一节存一节），
不要憋到最后一口气写一枚 600 行的大件——上一腿就是憋到最后才断的。
