# collection-r2 · 收件腿复跑判语（2026-11-09，编排者收件腿 collect-r2＝本波唯一被授权翻框的腿）

- 腿＝收件/裁决者（非实现者）；判据＝收件三把尺（提交真在／件真在／**自己复跑读数**）。
- 复跑台面＝`git clone --shared` 导出树 `/d/tmp/wisp/20261109/collect-r2-tree`，锚 **`329e76981`**（clone 时 dev 顶；比开工时共享树顶 `706fec5b4` 新一枚＝411 腿刚落，如实记）。
- 门禁（导出树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`）：`gofumpt -l cmd internal tools scripts`＝0（`/d/work/base/gopath/bin/gofumpt.exe` v0.12.0）／`go vet ./cmd/wisp ./internal/proc` rc=0／`go build ./...` rc=0／`sh scripts/d22scan.sh` clean／真 Linux 容器 `golang:1.27` `go vet ./internal/proc` rc=0。
- 本腿写面＝四票的 AC 框＋票尾 Progress log＋本件；AC 原文零字节（只翻勾＋追加）；零 `.go`；LF 核 CR=0；只 commit 不 push。

## 票 230（`230-four-cells-left-unfinished-inside-closed-tickets.md`）— AC#3 翻／AC#4 翻／AC#5 不翻

**AC#3＝翻。** 四把尺现跑：①票 21 现量未勾 **0**／已勾 **5**（现名 `21-approval-gates-minimal-done.md`）且 `blacklist_overrides` 命中 **0**；②`PLAN.md:2407` 逐字「**但豁免必须按文件、必须由人点、必须留痕。**」在；③`80:61` 交接句逐字在（作废走票 230 名下记录、原句未抹＝append-only 成立）；④读数件 `230-ac3-handoff-disposition.md`（`a3a4feefd`）＋复核件 `230-cells-closure.md`（`d0b9b95d1`）tracked 在库。与两腿报读数对不上处＝零。

**AC#4＝翻。** ①owner 面＝台账 `A1549`（2026-11-05 17:4x）在＝机主裁定走 B（原话「走 B 比较好…开放口子给用户」）⇒ 本格字面判据「把这枚摆给 owner」已兑现且已答；②接线在库＝`eacf043d0`（删 unwired 响亮失败＋`overrides.go` 103 行＋216 行测试）→`0d809934e`（装配措辞「已接线」）逐笔 `git show --stat` 真在；③盘上形状＝`unwired.go:125`「consumed: …B tier only, A tier refused…armed by cmd/wisp (armPersistedExemptions…)」、`risk.Gate/GateStanding` 真吃 `bOverrides`（`blacklist.go:83/:84`）；④我自跑＝`go test ./internal/config -count=1` ok 2.791s＋`./internal/risk` 定向 8 枚全 PASS（`TestGateStandingPersistedExemptionAllowsButSaysSo`／`…DoesNotUnlockA`／`TestOverrideOnlyAcceptsTheResolvedForm`／`TestBListDefaultDenyAndOverride` 等＝票 80 接收两格字面行为 (a)(b)(c) 逐格绿）。

**AC#5＝⛔ 不翻（如实）。** `230-cells`（11-08）照实未达终态：cmd/wisp 包级超时 1500.217s、两发两死点（真窗/热键活体测试逐个无限挂起，各归各票）。票不 `-done`。

## 票 140（`140-job-level-wisp-env-…-windows-runner.md`）— 四格全翻 ⇒ 改名 `-done`（下一笔 commit）

- **AC#1＝翻**：读数件 `140-ac1-realrun-shaping-r1.md`（`edab04ab6`）tracked，§5 三轴表＋§5.1 轴4 恒真点名＋§5.2 轴5 名册＋§2.3 CI 分母地图＝09-24 终裁③口径全覆盖；我复跑＝buildinfo＋secret 两形各一发 **27/27 PASS 逐名同色**（env 可归因换色＝0）；doctor 两形探针亲测（APPDATA 重定向一次性根）＝UNSET→`<APPDATA>\wisp-dev`（OS 被问）vs TEST→`C:\Users\swq\wisp-test-<pid>`（OS 不被问）＝job 级 env 使 OS 读永不被咨询成立。
- **AC#2＝终裁翻**：**收 (b) 为主＋(c) 三笔具名登记；⛔ (a)**（代价已量＝本机即 `[self-hosted, wisp-slo]` runner、摘 env 后 `probeWritable` 真写 `%APPDATA%\wisp-dev` 跨 run 复用同一工作树＝不免费）。落地两处逐字核过：ticket224 先栽后钉＋`wisp.db` 文件系统牙（`:53-102`）＋`dataroot_128_test.go` 文件系统期望（R-140-1 改码，`1afa3c280`）；前置 (i)(ii)(iii)＝票 123 族 HEAD **0 命中**（现跑）＝不在分母；ci.yml 一字未动、提案 diff 落档＝STOP 格守约。
- **AC#3＝翻**：变异 M1（摘 pin 行）导出树亲跑 ⇒ ticket224 **FAIL 16.32s**（红句逐字「the pinned WISP_ENV=test did not reach the child: boot 1 left no wisp.db in the injected data dir…」）；逐字还原（diff 空）后复跑 ok 15.793s。M2/M3 系腿报留档、我未重打（如实注明）。
- **AC#4＝翻**：scoped 9 枚顶层 `-count=2`＝**18 读数 0 FAIL 0 panic**（ok 34.340s）；gofumpt 0；vet 宿主 rc0；**真 Linux 容器 vet `./internal/proc` rc0**；`GOOS=linux go vet ./cmd/wisp` rc=1＝sherpa build constraint 已知形状（红句逐字核过、与被审对象无关）；d22scan clean（导出树＝纯净快照）；票 123 族 0 命中＝无放宽对象。

## 票 138（`138-cancelling-a-task-kills-nothing-…-no-terminatejobobject.md`）— AC#2 翻／AC#3 翻 ⇒ 四格闭环 ⇒ 改名 `-done`（下一笔 commit）

**Status 前提先核**：票面 `:3` `Status: blocked 等 Q-55` ＝**过期读数**——台账 `A1571`（2026-11-06 17:xx）逐字「B(211/215/256/345/350/354/174/177/**138**/150/148)全批 ⇒ 解锁 138」＝Q-55 已获机主批复；其后 11-07 r1/r2 两腿落的 AC#2 裁定＋修法＋牙时间线自洽。⇒ AC#2/AC#3 不再被「等批准」挡。

**AC#2＝翻。** 裁定书两半 tracked（`138-cascade-ruling.md`＝`8429d0e60`→`b3e990300`→`bf224be19`→`b576c10a3`→`9287d1179` 五笔＋`138-cancel-kills-cascade.md`＝r2 `c10a7d834`）。四条件逐条核：① risk/agent/proc/cmd 零触碰（六笔 stat 逐笔看＝只 manifestslot＋evidence＋票面）；② 无 go.mod/go.sum 改动；③ Windows API 在 `_windows.go` 半边（`command_job_windows.go:146` `windows.TerminateJobObject` 生产 1 处）、非 Windows 半边不静默（`command_job_other.go:478` 降级注释＋`command.go:444-445` Spawn 前打印真实降级）；④ cancelled 文案零改动。裁形状＝「每轮一枚自己的 Job Object、取消时级联」这一族由票 52 的 per-call Job 占位实现（每子进程 Job＋级联 Kill）；「不叠 per-task 层」的裁定我认：层级复杂买不到真价值（铁律③）。

**AC#3＝翻。** 变异 M1（`stopwatch.go:23` `WithDeadline(parent,…)`→`WithDeadline(Background,…)` 摘父链）导出树亲跑 ⇒ `Test138TaskCancelKillsTheProcessTree` **FAIL 15.22s 红在 `command_138_cancel_test.go:94`**（红句逐字「Exec did not return within 15s of the caller's cancel: the cancel cascade never reached the child」）＝与腿报 15.24s 同点同形；逐字还原（diff 空）后复跑 ok 1.387s。无 t.Skip、无阈值旋钮。

**AC#1/AC#4** 前已勾（09-26 验收程翻），不复述。

## 票 346（`346-two-race-window-instruments-…-blob-identical-flip.md`）— 七格本波**未裁未翻**，状态报明

- 两趟在盘核过：a1 普查件 `346-a1-census.md`（`6b959ee27`）＋编排者收件 log（10-19 22:5x，明确判语归后续非实现者裁）。
- AC#0 的 blob 同一性尺我现跑复现＝两枚仪器在 `a42b35af`⇄`a5080612`⇄HEAD 三点逐枚同串（`dfaddac8…`／`0b0e5173…`）、`openleaf_322` 在 `adec1db1` 不存在（rc=128）＝与票面逐字同。
- 剩余复跑面（各需专门趟）：CI run 日志逐名作差（39101060753 vs 39124122121，原始件 `D:\tmp\wisp\20261018\ci<id>.log`）、同族弃权仪器 12 处/8 文件名册核、三形代价表对票 342 对齐、AC#2 去牙反形（仓外导出树物理字节）、AC#4 双平台成对＋三发独立复跑。
- 本波时间盒收了 230/140/138 三票；346 状态＝**已立、普查第一趟收讫、七格待裁**。

## 纪律自证

- 四票 AC 原文零字节（只翻勾＋追加）；Status/票名原句未抹（138/140 的 `-done` 改名走 `git mv`，票面 Status 原句留档、判语里具名过期读数与解扣出处）。
- `git add`/`commit` 两处显式 pathspec 点名具体文件；零 `-A`/amend/reset/rebase/stash/checkout ./clean；零 push。
- 本腿触碰路径全集＝四票工单文件＋本件＋138/140 的 `git mv`；零 `.go`、零 thresholds/golden、零 winsec/frontend/forks/PLAN/specs/gate-clauses。
- 改动文件 LF 核：四票＋本件 CR=0（commit 前现跑）。
