# 176 裁决表（176-verdict）—— 裁决腿 2026-11-07 18:1x +08

- 本腿身份：**裁决腿**（非 176-r1 实现者、非 176-r2 复核腿），具裁决特权（可翻 AC 勾、改 Status）。
- 锚点：本腿 §1 现跑 `git log -1` ⇒ `29b3501d5`（2026-11-07 18:06:08 +0800）；工作期间 wave65/66（ball）在飞 commit 落盘，与本腿零交集；`git status --porcelain -- internal/ cmd/` **空**。
- 上一腿（176-r2）读数只当线索，**全部重跑**（见 §1）；两处更正见 §4。

## §0 裁定＝乙：按实落地翻勾结票（Status: done，票面改名加 `-done` 后缀）

理由（铁律 ⓪⓪⓪：判据＝最有价值＋最符合盘上事实，⛔「翻勾省事」也 ⛔「不敢改票」）：

1. **本票的核心交付是本票自己的腿做的**：起跑口三件产码/判据/接线的 commit 全部自报名号 176-r1——`41e8873bd`（`cmd/wisp/run.go` +20/−1，execute() 改调现成 `Loop.RunAsync`）＋`06dae8316`（`internal/tools/task_backfill.go` +117，名册写侧）＋`00f640a22`（`internal/tools/ticket176r1_start_port_test.go` +274，四枚常驻判据）。本腿逐一 `git show --stat` 现核，commit 均真在。
2. **AC#1 的作废条款写在实现之前**，它防的形状是「别的码把起跑口造出来了、本票还挂着没人认领」——那时才该「作废＋回改 164/175」。本票的触发者是**自己交付成功**：主流 harness 语义里这叫 completed（close as completed），不叫 withdrawn（close as not planned）。按字面把自己作废会让台账说谎：池里一枚 withdrawn 票，盘上挂着以它命名的三个活 commit 和今天仍在咬的判据（§3 牙读数）。
3. **「回改 164/175 的条件」这半句是真活但不在本腿写面**：164:57、175:19/:53/:67/:83 还留着「RunAsync 零调用点／Record 生产零写者」的历史 log——均为历史 log 且自带「归票 176」指针，176 结票（-done）后指针落在真话上，trail 自洽；逐字回改属 164/175 票面（超出本腿写面 `176-*.md`）⇒ 登记给编排者（§5）。
4. **运营价值**：`-done` 后缀是防重领唯一键——上一轮 roster-triage 恰因 176 未结票而错派（A1605 自认第 24 处筛档错）；结票防第三次空派。

## §1 独立读数（全部本腿现跑，⛔ 抄上腿）

| # | 读数项 | 结果 |
|---|---|---|
| 1 | `ls -la` 三件在场 | `internal/tools/task_backfill.go` 7784B（mtime 10-22）／`internal/tools/ticket176r1_start_port_test.go` 12047B（mtime 10-21）／票面 `.scratch/wisp/issues/176-no-spawner-zero-production-call-sites.md` 25289B；且三件均被 git 跟踪（`git ls-files` 命中） |
| 2 | **AC#1 两枚数（排测试口径）** | `grep -rn '\.RunAsync(' --include='*.go' internal/ cmd/ \| grep -v _test.go` ⇒ **2 枚**：`cmd/wisp/run.go:1834`（`bg := loop.RunAsync(ctx, task)`，注释块起于 ≈:1813）、`internal/tools/subagent_197.go:372`（`bg := child.RunAsync(childCtx, prompt)`）。`grep '\.Record(' 排测试` 原始 9 枚，其中 **TaskRoster 生产写者＝4 枚**：`task_backfill.go:141`／`task.go:411`／`task.go:424`／`subagent_197.go:462`；另 5 枚是 grants.Record（approval/gate.go:805）、memoryL1.Record（run.go:1853）、approval 卡/授权面（approval_reply.go:136、resident_approval_windows.go:293/:922）——**别的类型的 Record，不计入** |
| 3 | `go test ./internal/tools -count=1`（env 先 unset GOTMPDIR/TMPDIR/TEMP/TMP、PATH 加 build+sherpa-onnx） | **ok 33.962s** |
| 4 | 判据四枚真名（`grep '^func Test'`） | `:69 TestBackfillFilesTheRealSpilledArtifact176r1`（L-1）／`:128 TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1`（L-2）／`:217 TestNoRosterWriteAfterCancel176r1`（L-3）／`:249 TestTerminalWriteCarrierStaysOutOfScope176r1`（L-4） |
| 5 | 176-r1 三 commit 触碰面（`git show --stat`） | `41e8873bd`：仅 `cmd/wisp/run.go`（+20/−1）；`06dae8316`：仅 `internal/tools/task_backfill.go`（+117）；`00f640a22`：仅 `internal/tools/ticket176r1_start_port_test.go`（+274）——**契约轴零字节**（AC#5 证据） |
| 6 | AC#4 四条自证 | (i) `grep -rn "func (l \*Loop) [A-Z][A-Za-z0-9]*(.*taskID" 生产` ⇒ **0 命中**；(ii) `grep 'Tool\|Register' task_backfill.go` ⇒ **空**（不是受门控 Tool）；(iii)(iv) `grep 'go func(\|filepath.Clean\|filepath.Abs' task_backfill.go` ⇒ **0 命中** |
| 7 | AC#2 对应物 | `internal/tools/ticket175r2_stamp_live_test.go` 被跟踪（commit `58db7400` 真在，`git log --oneline -1` 命中）；175 本票已 `-done` |
| 8 | **牙（本腿自己的 overlay json，非 r2 的）** | 对照（无 overlay，`-run 176r1`）⇒ **ok 0.415s**；**M3** 摘 `Stopped(ctx)` 守卫（`:123-125`）⇒ **L-3 FAIL**，三断言齐响（`Look() 命中`／`Count() = 1, want 0`／`跳过回填却没交代理由`）；**M1** `rec.ArtifactPath = sp.Path` → `_ = sp.Path`（`:138`）⇒ **L-1+L-2 FAIL**，`超长输出没有产物路径：这就是 PLAN.md:2564 的「只截不指」`；每发之后 `git status --porcelain -- internal/tools/` **空** |

## §2 AC#1 作废条款为什么不按字面走（一句版）

作废条件确已触发（**2 ≥ 1、4 ≥ 1**），但触发者是本票自己的腿的**成功交付**；条款防的是「被别人抢先交付还挂着」的重复劳动，不是「交付完成后给自己开倒车」。结 done ＋ §5 把「回改 164/175」登记给编排者，两半义务都落地。

## §3 牙的跑法细节（可复算）

- 变异副本＝r2 现派生的 `.scratch/wisp/probes/176/launch-verify/mut{1,3}_task_backfill.go`；本腿 `diff` 现核：mut3 与现产码**只差** `:123-125` 三行守卫、mut1 **只差** `:138` 一行赋值（其余逐字节同）。
- 本腿**没有**用 r2 的 overlay json，自己写了两枚：`.scratch/wisp/probes/176/verdict-leg/overlay-mut{1,3}.json`（临时件，只建不删）。
- ⚠ 09-28 的老台件 `probes/176/r1/mut/*.overlay.json` 在今天 HEAD 上跑是 **[build failed]**（`task_state_188_test.go:190/:206` 找不到 `TaskBackfill.State`）——成因＝`62d7d3b49`（188-r2 片①）给 `TaskBackfill` 加了 `State` 字段，老 mut 副本还是 09-28 版。**这不是牙钝了，是台件过期**；任何人再复跑须用现派生副本。

## §4 对上一腿（176-r2）读数的两处更正

1. **行漂移**：r2/派单写的 `cmd/wisp/run.go:1813`，现行号＝**1834**（`:1813` 起是那段「起跑口」注释块，文件被后续腿改过）。
2. **判据件取代**：175 票面引的 `bridge_mark_provenance_ticket175_test.go` **已不在盘**，被 `ticket175r2_stamp_live_test.go`（`58db7400`，五腿）取代；AC#2 的对应物以后者为准。

## §5 给编排者的两条

1. **「回改 164/175 的条件」**：164:57、175:19/:53/:67/:83 的历史 log 仍写「RunAsync 零调用点／TaskRoster.Record 生产零写者」——都是历史 log 且自带「归票 176」指针（176 结票后指针自洽）。如要逐字回改，属 164/175 票面写面，由编排者定夺；不改也不产生误导（指针能走通）。
2. 台账登记 **A1606**（本裁决＋撤销口令）。

## §6 门禁（终态读数，结票 commit 之前现跑）

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/tools/` | **空**，rc=0 |
| `go vet ./internal/tools/` | rc=0 |
| `go test ./internal/tools -count=1` | **ok 69.067s**（同日早前第一次跑 ok 33.962s，两次全绿） |
| `sh scripts/d22scan.sh` | **clean**，rc=0（bans #1-5 internal/=379、cmd/=61；#7 internal/tools/=48；#8 internal/=944、cmd/=230、frontend/=116、design/=39；随行 runtests PASS=45 FAIL=0 SKIP=0） |

- ⛔ push；零产码；本腿写面＝票面 `176-*.md`＋本表＋台账 A1606，`internal/**`、`cmd/**`、`frontend/**`、`design/**`、thresholds/golden 零字节。
