# 票 360 落地腿 `360-r1` 交付件（AC#2 落地＋牙）

- 时刻＝2026-10-22 17:3x–18:0x +08｜起手锚＝`git rev-parse HEAD` 现认 `5ae6716c`（派单提交）｜本腿两笔：`d68ddfae`、`fc7c85cc`。
- 形＝〔甲〕**只加一枚用例**（编排者 `AC#1` 已拍）。**零产码**：`allowdirs.go`／`writeguard.go`／`manager.go`／`settings.go`／`permmode.go` **零字节**（`sha256` 见 §5）。
- 写面＝`internal/config/writeguard_adopt_360_test.go`（新用例）＋本件。临时件全在仓外 `D:\tmp\wisp\20261022\360r1\**`。

## 1. 开工前（派单「⛔ 跳步」三条，逐条现跑）

1. **`git status --porcelain` 写面**：`-- internal/config docs/evidence/s1` **零行**（写面内无别人的未提交改动）⇒ 可开工。
2. **票面前提逐条现跑**（对 HEAD blob 现认，⛔ 抄 `360-a1`）：
   - `statOwnWrite` 分层＝**5 行 ＝ 1 定义 ＋ 1 代码调用 ＋ 3 注释**，与票面逐枚同：定义 `allowdirs.go:157`（体 `:157-163`，只写 `seenMtime`／`seenSize`／`seenValid`）／调用 `writeguard.go:179`（在 `mergeWrite` 的 `if len(diverged)==0` 支内）／注释 `allowdirs.go:150`／`writeguard.go:105`／`writeguard_226_test.go:186`。
   - 写回入口 3 枚 → `mergeWrite`：`allowdirs.go:142`／`permmode.go:79`／`settings.go:303` 均 `m.mergeWrite(...)`；汇聚 `writeguard.go:111`，`:118 readConfigFile` 读盘当 base、`:163 SaveFile` 从不回读。
   - Adopt 全仓一条：`manager.go:142 CheckAndReload`。⇒ **两条路互不相连**，前提今天成立。
3. **复跑 `360-a1` 的复现＋正控两组读数**（仓外 `-overlay` 探针 `zz360r1_probe_test.go`，仓内零 `.go`；件 `logs\01-repro-and-control.txt`）：
   - `READ A2 after writeback: file.ball.size=60 mem.ball.size=56`
   - `READ A4b diffKeyPaths(memory, file-after-write) n=1 [ball.size]`
   - `READ B1 after Adopt: report.Hot=[ball] mem.ball.size=60 file.ball.size=60`
   - `READ C0 clean write: file.ball.size=56 mem.ball.size=56 poll-report=<nil>`
   ⇒ **今天复现**（与 `360-a1` 同形同值），**不触发**派单的停手上报条②。

## 2. 新用例（一枚，题目就是这一跳）

`internal/config/writeguard_adopt_360_test.go`

- **`TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit`**（提交 `d68ddfae`，两态打印由 `fc7c85cc` 补）
  - 题目＝**写回之后内存仍持旧值**：同一 `Manager`（开机 `ball.size=56`），手改文件到 `60`，跑一次受守卫写 `AddAllowedDir` 后**先不 reload**，直接读两侧。
  - 三条断：`file.ball.size == 60`（手改还在盘上）／`mem.ball.size == 56`（**没被 Adopt**）／`diffKeyPaths(memory,file) == [ball.size]`（分歧就这一枚键）。
  - ★**自带两态正控**（同一台面、同一 manager 再跑一次 Adopt 路 `CheckAndReload`）：断 `mem=60 && file=60 && report.Hot` 含 `ball`；并**显式断两态读数不同**（`adoptedMem != memSize`），台子若分不出两态即 `Fatal` 判 vacuous。读数起点设成不可能的 `-1`，**摘掉那一臂就必红**（见 §3 M1）。

**为什么它今天响、而不是恒真**：这两条断言今天**已被别的票、别的格的顺带断言覆盖**（普查腿实测：写回路径真 adopt 时当场红 9 枚：`settings_248_test.go:206`／`writeguard_226_test.go:89,:224,:248,:259,:285`／`writeguard_handedit_227_test.go:50`／`writeguard_unnailed_341_test.go:183,:327`／宿主 `always_write_no_clobber_226_test.go:119`），但**没有任何一枚以这一跳为题目、且自带两态正控**。本用例补的就是那个题目 ⇒ 它今天**真的能红**（§3 三发变异逐发证）。

## 3. 反向判据（三读数逐字，全在仓外导出树 `tree2` 上跑；仓内零 `.go` 改动）

变异前先证落地（`grep -n`）＋ `go build ./...` rc=0，再读数。**串行**。

| 发 | 造法（落地证明） | build rc | 读数（逐字） |
|---|---|---|---|
| **原样绿** | 无变异（`logs\14-bench-pristine-green.txt`） | 0 | `STATE 1 after writeback: file.ball.size=60 mem.ball.size=56` / `STATE 2 after one Adopt pass: mem.ball.size=60 file.ball.size=60 report.Hot=[ball]` / `--- PASS: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit` / `ok ... 0.103s` ⇒ **rc=0** |
| **M1 变异红**（摘掉 Adopt 那一臂） | `sed -i "73,81d"`；`grep -n CheckAndReload` 只剩注释 `:19`，代码调用 **0** | 0 | `STATE 1 ... file.ball.size=60 mem.ball.size=56` / `STATE 2 after one Adopt pass: mem.ball.size=-1 file.ball.size=-1 report.Hot=[]` / **红行 `writeguard_adopt_360_test.go:76`**：`the Adopt control did not produce the other state: mem.ball.size=-1 file.ball.size=-1 report.Hot=[], want 60 / 60 / [ball]; without it the readings above cannot tell two states apart` / `--- FAIL` ⇒ **rc=1** |
| **M2 变异红**（产码侧显式 Adopt：`writeguard.go:169` 后插 `m.cur = deepCopyConfig(base)`，`grep -n` 命中 `:170`） | 同上 `grep -n MUTATION-360r1` ⇒ `:170` | 0 | `STATE 1 ... file.ball.size=60 mem.ball.size=60` / **红行 `:60`** `the writeback adopted the hand edit into memory: mem.ball.size=60, want the 56 this process loaded; adopting here is the reload path's job, behind ConfirmLocked (D36 rule 1)` / **红行 `:65`** `the memory/file divergence after the writeback = [], want exactly [ball.size]` / `STATE 2 ... report.Hot=[]` / **红行 `:85`** Adopt 控制未产出另一态 ⇒ **rc=1** |
| **M3 变异绿（照实记账）**：把状态 1 三条断言改成恒真（`:55`／`:59`／`:64` 的 `if false`） | `grep -n "if false"` ⇒ 3 命中 | 0 | `STATE 1 ... 60/56`、`STATE 2 ... 60/60 [ball]`、`--- PASS` ⇒ **rc=0**。★这不是缺陷：派单判据原文是「摘掉 Adopt 臂**或**把断言改成恒真 ⇒ 必须转红」；**M1（摘臂）已满足**。M3 单独改恒真仍绿，是因为正控臂与"两态不同"那两条仍在咬 ⇒ **记成"判据的另一半由 M1 承担"**，⛔ 声称 M3 也红。 |
| **M3＋M2 组合红**（恒真断言 ＋ 产码 adopt） | `grep -n MUTATION-360r1b` ⇒ `:170` | 0 | `STATE 1 ... file.ball.size=60 mem.ball.size=60` / `STATE 2 ... report.Hot=[]` / **红行 `:85`** ⇒ **rc=1**。⇒ 正控臂**独立成牙**，不是靠状态 1 那三条撑着的。 |

- **还原绿**：两枚件从仓外原样副本拷回 ⇒ `go build ./...` rc=0、`--- PASS` ⇒ rc=0（`logs\20-restored-green.txt`）。
- **整包（导出树还原后）**：`go test -count=1 ./internal/config/` ⇒ `ok github.com/CarlosShao/wisp/internal/config 2.444s` rc=0。

## 4. 还原核 `sha256`（三处逐字相同）

| 件 | (1) 提交版 HEAD blob | (2) 还原后的仓外台 `tree2` | (3) 主检出工作树 |
|---|---|---|---|
| `internal/config/writeguard_adopt_360_test.go` | `481fb1626f4307fe8ddca4d9d67dc43a7d20aadcc7cc8787cab5ee5071462c7e` | 同串 | 同串 |
| `internal/config/writeguard.go`（**零改动**） | `efe12661e71b940f1fcf99d722f1c25e4403431997e90041491d5eae337d47bd` | 同串 | 同串 |

（`writeguard.go` 那串与 `360-a1` 记的 `efe12661…47bd` 同值 ⇒ 本腿确实一个字节没动它。）

## 5. 门禁各 rc（主检出；`GOFLAGS=`／unset `GOTMPDIR TMPDIR TEMP TMP`／未设 `CGO_ENABLED`）

- `go build ./...` ⇒ **rc=0**
- `go vet ./internal/config/` ⇒ **rc=0**
- `go test -count=1 -v ./internal/config/` ⇒ **rc=0**，四数**分开报**：`RUN=196`／`PASS=115`（顶层）／`FAIL=0`／`SKIP=0`，另 `^panic:` **0**（⛔ 相减）
- `sh scripts/d22scan.sh` ⇒ **rc=0**，末行逐字 `d22scan: clean - no D22 ban violations`（分母 `ban #8 internal/ 617 Go files` 是**文件数**、⛔ 违规数）
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0**，`VERDICT GREEN`、`tracked paths=17571  over-budget=61  covered by roster=61  not in roster=0`（读数取自本件入库**前**那棵树；本件自身入库后分母变 **17572**，同一把尺在最终 HEAD 上复跑仍 `rc=0`／`VERDICT GREEN`／`not in roster=0`，见 §6 末行）
- **格式门禁两把尺（裁定件 `pipe-format-gate-reads-tracked-blob-not-worktree.md`）**：
  - **甲＝入库 blob 尺（判据）**：`git archive HEAD` 到仓外 `tree-final`，对 1156 枚 `.go` 跑 `/d/work/base/gopath/bin/gofumpt.exe -l`（`v0.12.0 (go1.27.1)`）⇒ **列 20 枚**，**全部在 `.scratch/wisp/probes/**`**（本仓既有的故意做坏样本，属别的票的台件、**不在本票写面**）；**本腿新增件 0 命中**（`grep -c writeguard_adopt_360_test.go` ⇒ `0`）。★**对照读数**：同一把尺打父提交 `HEAD~1` 的导出树 ⇒ **同样 20 枚、`diff` 逐字相同** ⇒ **本腿新增 0 枚**。
  - **乙＝工作树尺（只作诊断）**：整树 ⇒ **列 55 枚**、rc=123（`xargs` 因列表过长分批后最后一批 rc，⛔ 当判据）；**本腿新增件 0 命中**；单文件直跑 `gofumpt -l internal/config/writeguard_adopt_360_test.go` ⇒ **空、rc=0**。⇒ 甲绿乙红的差＝本机检出副本（`core.autocrlf` 假红），按裁定件 §4.2 **登记一句、⛔ 判该票失败**。
  - ⚠ **CI 那一步真正调用的仪器照实登记**：`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **rc=1**，`(A) tracked-tree lines=35 files=35`。**这 35 与甲尺那 20 的差我已查清、⛔ 留悬念**（件 `logs\25..27`）：`comm` 两向 ⇒ **只在 attrib 有 = 15 枚、只在 blob 有 = 0 枚**；那 15 枚**逐枚现量 CR 数**（工作树 CR 185/176/…/1115，`HEAD:` blob CR **全 0**）⇒ **全是本机 `autocrlf` 的 CRLF 检出副本**，与裁定件 §二那张表同形。⇒ **35 ＝ 20 枚真脏 blob ＋ 15 枚 CRLF 假红**，两把尺**在"真脏"上完全一致**（`comm -13` 空）。
    那 20 枚真脏里：**20 枚全在 `.scratch/wisp/probes/**`**（甲尺名单，别的票的故意做坏样本）；attrib 的 35 行里另有 `cmd/wisp/models.go`／`internal/agent/{approval/pending_read,tools}.go`／`internal/risk/provenance.go` 四枚产码——**这四枚正是上表那 15 枚 CRLF 假红之四**，**不在甲尺名单里**。
    **这 35 枚无一属本票写面、无一由本腿引入**（父提交对照 0 增量、本腿新增件在两把尺上均 0 命中）⇒ 属**既有 CI 状态**，本腿⛔ 修（在写面外）、⛔ 隐瞒、⛔ 据此判本票失败。

## 6. 写面尺（逐笔 `git show --name-only` 现核）

| 提交 | 主题 | 文件数 | `.go` | 越界 |
|---|---|---|---|---|
| `d68ddfae` | 票 360 落地腿 360-r1: 加一枚用例…（零产码） | 1（`internal/config/writeguard_adopt_360_test.go`） | 1 | **0** |
| `fc7c85cc` | 票 360 落地腿 360-r1: 用例自印两态读数… | 1（同件） | 1 | **0** |

- **正控尺**（同尺打一笔含 `.go` 的真提交 `ecebc7a9`）⇒ 文件 4、`.go` **2**（`cmd/wisp/always_write_no_clobber_{226,227}_test.go`）⇒ **尺不空转**。
- **负控**（纯文档笔 `5ae6716c`）⇒ `.go` **0**。
- 区间内另有编排者的 `1cfe94bb`（只动 `docs/reports/pending-and-issues.md`、`internal/config` **0** 枚）⇒ **⛔ 算本腿写面**。
- 只 commit、**⛔ push**；`add`／`commit` 两处均**显式 pathspec**、新文件**逐个文件名字面点名**；未 amend／reset／rebase／stash／checkout ./clean；未翻任何框、未改票名／`Status:`／`-done`／台账。

## 7. 未做／边界（照实记账）

- **⛔ 改产码**：`allowdirs.go`／`writeguard.go`／`manager.go`／`settings.go`／`permmode.go` 零字节（§4 两处 sha 可核）。**未触发**派单的停手上报条①。
- ★〔甲〕**⛔ 意味着"这一跳不重要"**：本件只钉"**今天**这一跳不该被产码改变"（一改就撞 D36 rule 1 `ConfirmLocked` 先行）。将来要改产码＝**独立行为变更、须重新立案**，⛔ 在本票射程内做。
- **M3 单独恒真仍绿**已照实登记（§3），判据的"摘臂"那一半由 **M1** 承担，⛔ 声称三发全红。
- **`cmd/wisp` 整包未跑**（派单口径 `A1126`／`A1127`：整包读数⛔ 当判据；本腿未需要）。
- 格式门禁**第三把仪器**（`attrib.sh --tracked-only`）的既有 35 枚红已登记，**归属与是否该修由编排者裁**（在写面外）。
- **最终 HEAD 复跑（本件入库后）**：`go build ./...` rc=0／`go vet ./internal/config/` rc=0／`go test -count=1 ./internal/config/` rc=0（`RUN=196 PASS=115 FAIL=0 SKIP=0`）／`d22scan` rc=0／`pathlen` rc=0（`tracked paths=17572  over-budget=61  covered by roster=61  not in roster=0`、`VERDICT GREEN`）／甲尺 `git archive HEAD` ⇒ **仍 20 枚、与父提交 `diff` 逐字相同、本腿新增件 0 命中**。⇒ **§5 的每个 rc 在最终 HEAD 上原样成立**。
