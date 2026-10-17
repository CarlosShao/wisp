# 328-r1 落地件②：台面／颜色成对／门禁四数／AC#3 四面／git 自证

腿＝`328-r1`。起手锚＝`648eda04`（现认），交件侧锚＝`f705b082`。四数尺＝`^=== RUN`（含子测试）／`^--- PASS`（只顶层）／
`^--- FAIL`／`^--- SKIP`，**逐枚点名、未混用、未相减**；一律 `-count=1 -timeout 30m`（显式）；⛔ `go test ./...`、⛔ 跑 `./cmd/wisp`。

## 一、台面（仓外 `git archive` 导出树，两枚同一配方）

配方脚本＝`.scratch/wisp/probes/328/r1/logs/build-tree.sh`（八件套 `go.mod go.sum internal cmd tools scripts frontend models`
＋仓根 `docs/`＋手工三件 `build/*.dll`／`third_party/sherpa-onnx/*.dll`／`frontend/dist`）。

| 读数 | 起手侧树 `wisp-328r1-before-tree`（锚 `648eda04`） | 交件侧树 `wisp-328r1-after-tree`（锚 `f705b082`） |
|---|---|---|
| `.go` 枚数（证非空） | **746** | **748**（＝我新增两枚测试文件） |
| `docs/PLAN.md` | present | present |
| `build/*.dll`／`sherpa/*.dll`／`frontend/dist` | 3／3／present | 3／3／present |
| 8.3 短名别名（`cmd /X`，`D:` 卷） | `WI8846~1` | `WI98B9~1` |
| `GOTMPDIR` | `D:\tmp\328r1-gotmp`（先 `mkdir -p`） | 同 |
| 名册那枚 blob | `35ce6ffb8603…`（＝普查件对锚同一串） | `2da6c596369c…`（＝本腿改后） |

⛔ 共享工作树跑 `go test`／`go build`：本腿所有颜色都在上面两棵树里跑；共享树只跑门禁四数（下一节）。

## 二、颜色：改前⇄改后成对（同配方、同机、同 `-timeout 30m`）

| 包／射程 | RUN | PASS_TOP | FAIL_TOP | SKIP | rc |
|---|---|---|---|---|---|
| `./internal/panel` 改前（整包，`before-panel-full.log`） | 257 | 160 | 7 | 0 | 1 |
| `./internal/panel` 改后（整包，`after-panel-full.log`） | 268 | 163 | 7 | 0 | 1 |
| `./internal/agent` 改前（整包，`before-agent-full.log`） | 112 | 88 | 0 | 0 | 0 |
| `./internal/agent` 改后（整包，`after-agent-full.log`） | 114 | 90 | 0 | 0 | 0 |
| `./internal/panel` 定点（本腿四形＋票 220 那族，`after-panel-targeted.log`） | 19 | 5 | 0 | 0 | 0 |
| `./internal/agent` 定点（`ⓐ` 正控两支，`after-agent-targeted.log`） | 2 | 2 | 0 | 0 | 0 |
| `./internal/panel` 定点 **改前**（`:141/:142` 那族表，`before-panel-220only.log`） | 8 | 2 | 0 | 0 | 0 |

作差是"我新加了几枚用例"的读数（panel `+11 RUN / +3 PASS_TOP`＝3 枚顶层＋8 枚子测试；agent `+2/+2`），
⛔ 拿四数相减当判据。★改前那一行（RUN=8／PASS_TOP=2）是 `ⓑ` 的起手侧：`child-1#retry` 当时钉 **false** 且绿。

**终态红名册 ⇄ 起手作差＝0**：尺＝`grep '^--- FAIL' <两枚日志> | sed 's/ ([0-9.]*s)//' | sort -u` 作差，
`redroster_diff_rc=0`，两侧各 **7 枚名字、逐名同集**（名册＝`logs/red-before.txt`／`logs/red-after.txt`，
尺与输出逐字＝`logs/red-roster-and-ac3-extra.log`；`./internal/agent` 两侧红名册各 0 枚，同日志）。七枚名字与归属（红句逐字在各日志里）：

1. `TestApprovalCardViewJSONKeysMatchFrontendTypes` → `approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`
   ＝普查件点名的**起手即红前端契约键**那一枚（与卡片无关；本腿⛔ 加键 ⇒ 它没有多出第二句）。
2. `TestComposerContractTypesMatchFrontend` → `composer_test.go:78:` 同句同族 ＝普查件点名的另一枚起手即红。
3. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` → `backpressure_bound_35r8_test.go:278/309/336`（heap／fan-out 实测读数）。
4. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8` → 同文件（`:336` 那句 `pinned dropped-name runes 211994 > the cap 36864`）。
5. `TestComposerRenderFixtureTellsTheTruth` → `composer_test.go:873: the render fixture has no block for "no host attached" …`。
6. `TestReadGitOnThisRepositoryIsSelfConsistent` → `git_test.go:600: kind = "not_a_repo" …: this repository has a .git directory`
   ＝**台面局限**（导出树里没有 `.git`），改前改后同红 ⇒ 作差里互相抵掉；⛔ 为它放宽任何断言。
7. `TestC21DesignTokensFourWayAgree` → `tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\…\design\assets\tokens.css: The system cannot find the path specified.`
   ＝**台面局限**（派单那套手工三件不含 `design/`，那是机主外部 agent 的写面），改前改后同红。
   ★这枚是 AC#3 里"三枚冻结件"之一，本腿**一个字节没动**它（diff 见下一节），它红在台面射程上，不是红在我身上。

## 三、门禁四数（只在共享工作树跑，逐枚 `rc=`；日志＝`gates-328r1.log`）

- `gofumpt -l <本腿四枚文件>`（真身 `/d/work/base/gopath/bin/gofumpt.exe`）→ 无输出，`gofumpt_rc=0`
- `sh scripts/d22scan.sh` → `d22scan: clean - no D22 ban violations`，`d22scan_rc=0`
  （射程现量：ban #8 `internal/` 586 Go files、`cmd/` 143、`design/` 39 text files ⇒ 不是"空射程"那种 `exit 2`）
- `sh scripts/check-path-length-budget.sh` → `VERDICT GREEN`，`length_rc=0`（分母 tracked paths=**15035**，over-budget=60 全在册）
- `go vet ./internal/panel/ ./internal/agent/` → `vet_rc=0`

## 四、`⓮`＝AC#3 四面 diff 逐枚报 0（日志＝`ac3-four-faces.log`）

尺＝对本腿**每一枚自己的提交**跑 `git diff --stat <c>^ <c> -- <面>` 数行数（共享树里同一段还有别腿落笔，
whole-range 作差只当归属读数用）：四枚提交 × 四面＝16 格，**每格 `diff_lines=0 rc=0`**。
四面＝门控档位（`internal/risk/mode.go`＋`internal/agent/approval/gate.go`＋`…/queue.go`）／面板侧 L2 出处
（`internal/risk/provenance.go`＋`internal/panel/approval.go`）／D43 转移表（`internal/statemachine/`）／
三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）。
归属复核（量的）：`git log --oneline 648eda04..HEAD -- <四面>` 逐面**均无提交**（连别腿也没碰）；
`git status --porcelain -- <四面>`＝空。⛔ 把 L2 卡片文案复用权限判定文案（票 242/244 那族）：本腿没动任何卡片文案，
`ApprovalCardView` 字段名册一字未改（`grep -c TaskID internal/panel/approval.go` 交件侧仍 0）。

## 五、git 自证（⛔ push；⛔ `add -A`／`.`；⛔ amend/reset/rebase/stash/checkout ./clean；⛔ 仓内 worktree）

起手三件（`logs/start-anchor.txt`／`logs/start-status.txt`）：`date`＝`Sat Oct 17 12:38:56 CST 2026`，
起手锚＝`648eda040deabd7906bb449d1d47274bf731d202`，`git status --porcelain`＝**1179 行**（其中跟踪改动 3 枚＝
票 328／334／335 的 `Status:` 行＝**编排者在飞字节**，本腿⛔ 把它们记进任何提交；两枚新测试文件在写面内起手为未跟踪＝零枚）。

四枚提交（逐枚 `git show --name-only`，全在写面内）：

| 提交 | 文件 | 行数 |
|---|---|---|
| `b2d55342` | `internal/panel/subagent_roster_197.go` | +66／-10 |
| `1696d01b` | `internal/panel/subagent_blocked_220_test.go` | +12／-4 |
| `ffebb6d1` | `internal/panel/subagent_corr_task_prefix_328_test.go`（新增） | +89 |
| `f705b082` | `internal/agent/corr_roster_join_328_test.go`（新增） | +171 |

期间别的提交在飞（例：编排者 `fb7a4d1a` 落 `docs/reports/HANDOVER.md`，就夹在我第 1、2 笔之间）⇒ 每笔都只用
逐枚显式 pathspec，`git add -- <一枚文件>` ＋ `git commit -F <正文文件> -- <同一枚文件>`；正文先 Write 成
`.scratch/wisp/probes/328/r1/commit-msg-*.txt`（中文正文走 `-F`，⛔ 双引号内联）。
凡引用别腿字节一律取具名提交树（`git show <锚>:<路径>`；本腿读过的别面＝`cmd/wisp/panel_pump.go`、
`internal/agent/approval/queue.go`、`internal/agent/approval/pending_read.go`、`cmd/wisp/panel_snapshot_deliver_329.go`，
全部 `git show f705b082:…`／`git show 648eda04:…`，⛔ 工作树在飞字节）。

收尾（`logs/end-status.txt`＝**1238 行**，其中跟踪改动 2 枚，见 `logs/end-status-summary.txt`）：
` M .scratch/wisp/issues/328-…-join-key.md`＝**本腿那一枚只追加的 bullet**（`git diff` 逐字＝+2 行：一空行＋我那一枚，
⛔ 第二个人的字节；那一枚 `Status:` 改动已由编排者在 `ed6fa14e` 入库，本腿追加前先跑一把尺证该文件对 HEAD 为空）；
` M .scratch/wisp/issues/335-….md`＝**别腿在飞字节，本腿⛔ add ⛔ commit**（逐枚 pathspec 挡住）。
本腿产码／仪器四枚文件全部入库，`git status --porcelain -- internal/panel internal/agent`＝**空**（现量 rc=0 无输出）。

## 六、本腿没做的事（具名，⛔ 当已做）

- ⛔ 翻 AC 框（AC#2/#3/#4 七格留给非实现者裁）；⛔ 改票名／`Status:`／判据文字／加 `-done`。
- ⛔ push；⛔ 跑 `./cmd/wisp` 整包（`228-r3b` 在取证）⇒ 那包的名册读数＝零，静态界见件①第五节。
- ⛔ 动 `internal/risk/**`、`frontend/**`、`design/**`、SLO 阈值／golden／`thresholds.go`、三枚冻结件、`go.mod`／`go.sum`。
- ⛔ 新 JSON 键／新导出 `panel.*` 名／新 C17 方法名；⛔ 用 `go build`／`go vet` 的 rc=0 当"接上了"（本腿判据是字节形状＋调用点）。
- ⛔ 杀活实例 PID 21892；⛔ 删任何临时件（两棵导出树、日志、正文文件全留）。
