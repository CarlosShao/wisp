# 176-r2 复核表 — 后台起跑口：派单前提过期的独立现量（零产码）

- 腿：**176-r2（写手位收到的派单是「实现起跑口＋加用例」，本程起手三核发现该活已由 176-r1 于 09-28 落地并提交 ⇒ 本程转为独立复核腿，零产码、未勾任何框）**。本程不是 176-r1 的实现者，本表全部读数为本程现跑。
- 锚点（本程 §0 现跑）：`a1c3f245b2b35d01ee31ade15908283641e1f5ac`（短号 `a1c3f245b`）/ 分支 `dev` / `git status --porcelain -- internal/tools/` **空**。
- 台件：`.scratch/wisp/probes/176/launch-verify/`（两枚突变副本＋两枚 overlay JSON，只建不删；突变经 `go test -overlay` 落地，跟踪文件零改动，表 §4）。

## §0 结论一句

**票 176 的「后台起跑口」已在盘上**：产码 `internal/tools/task_backfill.go`（commit `06dae8316`）＋判据 `internal/tools/ticket176r1_start_port_test.go`（commit `00f640a22`）＋组合根接线 `cmd/wisp/run.go:1813`（`RunAsync` → `Wait` → `MarkRoot` → `Backfill`）。本程未重复实现任何一块；今天新增的只有**新鲜突变两发**（牙仍咬的现证）与本表读数。AC 勾与 Status 的处置归编排者（AC#1 的作废条件 176-r1 自己早已如实触发并上报，见票面 09-28 09:3x 段）。

## §1 AC#1 两枚数现量（排测试口径，命令原文）

| 数 | 票面立票时 | 本程现量（2026-11-07） | 命令 |
|---|---|---|---|
| `RunAsync` 生产调用点 | 0 | **2**：`cmd/wisp/run.go:1813`（176-r1 的起跑口本体）＋ `internal/tools/subagent_197.go:372`（票 197 的子代理腿） | `grep -rn "\.RunAsync(" --include="*.go" internal/ cmd/ \| grep -v _test.go` |
| `TaskRoster.Record` 生产写者 | 0 | **4**：`internal/tools/task_backfill.go:141`（176-r1 回填）＋ `internal/tools/task.go:411`（`PublishSubagent`，票 197）＋ `internal/tools/task.go:424`（`MarkRoot`，调用点 `cmd/wisp/run.go:1829`）＋ `internal/tools/subagent_197.go:462`（票 197） | `grep -rn "Roster.Record(\|\.Record(taskID" --include="*.go" internal/ cmd/ \| grep -v _test.go` |
| `ArtifactPath:` 字面量赋值点 | 0 | 仍 **0**（真路径走 `rec.ArtifactPath = sp.Path`，`task_backfill.go:138`） | `grep -rn "ArtifactPath:" --include="*.go" internal/ cmd/ \| grep -v _test.go` |

⇒ **两枚都 ≥1：票面 AC#1 的「本票作废、回来改 164/175 的条件」停手条件早已触发**；本程不复述裁定，只交读数。

## §2 硬前置票 175 的落地接口今天真在盘上（⛔ 信票面，全部现跑）

- 票面件名带 `-done` 后缀：`.scratch/wisp/issues/175-task-output-returns-external-content-with-no-c25-provenance-mark-because-marking-is-name-gated-and-the-name-is-not-on-the-roster-done.md` 在场。
- 名册：`internal/risk/provenance.go:94` `SrcTaskOutput = "task.output"`、`:102` 入名册切片（八枚变九枚那笔，commit `962ba9d1`）。
- 桥级：`internal/tools/bridge.go:744` 的早退闸门仍是名册制，`:737` 注释逐字 `task.output joined that roster`；判据件 `internal/tools/ticket175r2_stamp_live_test.go`（commit `58db7400`）在场，本程全包测试里随包跑过（§3）。

## §3 门禁（全部现跑；env 先 `unset GOTMPDIR TMPDIR TEMP TMP`、`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`）

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/tools/` | **空**，rc=0 |
| `go vet ./internal/tools/` | rc=0 |
| `go build ./internal/tools/` | rc=0 |
| `go test ./internal/tools/ -count=1` | **ok 34.995s**，rc=0 |
| 同包 `-v -run 176r1` | 四腿逐枚 **PASS**（L-1 0.03s／L-2 0.03s／L-3 0.01s／L-4 0.00s） |
| `sh scripts/d22scan.sh` | **clean**，rc=0（ban#7 internal/tools/=48；ban#8 internal/=942、cmd/=226）；随行 `runtests.sh` PASS=45 FAIL=0 SKIP=0 |

## §4 ★牙今日新鲜两发（`go test -overlay`，跟踪文件零改动；每发之后 `git status --porcelain -- internal/tools/` 复验为空）

| 变异 | 改法 | 响的腿 | 读数 |
|---|---|---|---|
| **M3** | 删 `task_backfill.go:123-125` 的 `Stopped(ctx)` 守卫（突变副本 `mut3_task_backfill.go`） | L-3 `TestNoRosterWriteAfterCancel176r1` | **FAIL rc=1**，三发断言齐响：`Look() 命中`／`Count() = 1, want 0`／`跳过回填却没交代理由` |
| **M1** | `rec.ArtifactPath = sp.Path` → `_ = sp.Path`（突变副本 `mut1_task_backfill.go`） | L-1 `TestBackfillFilesTheRealSpilledArtifact176r1` | **FAIL rc=1**，`超长输出没有产物路径：这就是 PLAN.md:2564 的「只截不指」` |

⇒ **176-r1 判据的牙在今天的锚上仍然咬**（不是过期读数）；对照：无 overlay 时同两枚 PASS（§3）。

## §5 AC#4 形状禁区自证（对本程自己的交付）

- (i) 本程**零产码**：没有给 `Loop` 加任何方法，也没有新写任何带 `taskID` 形参的函数（`git diff` 对 `internal/agent/`、`cmd/wisp/` 均空）。
- (ii) 没有把宿主 artifacts 写入暴露成受门控 Tool（没有 `Register|Decl{|Entry{` 新增）。
- (iii) 本程新增代码里没有 `go func(`（零产码，突变副本是既有文件的 sed 派生）。
- (iv) 本程没有引入 `filepath.Clean|Abs`。
- ⚠ G2 陷阱自证：本程在 `cmd/wisp/**` 零字节，未新写任何含 `OpenTask|CloseTask` 字面的说明。

## §6 没做到的／留给编排者的

- **「按票面实现起跑口」这一条本程没做**——不是做不到，是**已有人做完并在盘上**（§0），重复实现＝无收益复杂度（⓪⓪⓪ 第 3 条）。派单前提需编排者更正（今晨筛档 `roster-triage-20261107.md` 把 176 当「未落地」排派，与盘上事实不符）。
- AC 框（AC#2-#5）的翻勾、票面 Status、以及 AC#1 作废条件触发后的改票动作（票 164/175 的条件回改），**全部归编排者**，本程零触碰。
- 本波 `internal/tools/**` 只有本腿一枚写腿（176-r1 是 09-28 的历史腿，非本波在飞）。
