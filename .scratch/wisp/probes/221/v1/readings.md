# 221-v1 逐发读数（非实现者验收腿）

锚点：`git rev-parse HEAD` = `29215a93696a2127b4558c2417554c6f04d14af6`（本腿自取）
可迁移性：`git diff c44b30c4..HEAD -- internal cmd` 为空；`internal/tools/task.go` md5 `4138177e29ffff427776a73eb0d61a9b` ＝ HEAD 同值；`internal/tools/subagent_197.go` md5 `06caf8f5465ff1c47c310da27fe3ffab` ＝ HEAD 同值。

## 起跑卫生

- 每次 `go test` 前：`tasklist //FI "IMAGENAME eq go.exe" | grep -c go.exe` = `0`（逐次现跑，见各节末）。
- 一次只跑一发；突变全部 `-overlay`，全程 `git status --porcelain -- internal cmd` = 空（下表末行）。

## 名册（asis 对照）

命令：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools -run 'Test221|Test197|TestEveryRegisteredToolIsClassifiedForMarking175r2' -v -count=1`
日志：`logs/mut-asis.txt`
读数：rc=0 ／ 两层结果行 20 枚 ／ `--- FAIL` 0 枚
绿名册（20）：Test197SpawnPublishesIdentityRow, Test197RowExistsBeforeFirstChildModelCall, Test197SubagentPoolNeverExceedsBridgeCeiling, Test197SpawnDescriptionNamesTheRealPoolCap, Test197SubagentPoolCapsAtBridgeCeiling, Test197FullPoolRefusesNextSpawnWithReadableReason, Test197ChildCannotDeriveSubagent, Test197ConclusionCarriesTaskOutputTaint, Test197CancelIsPerRowAndNeverCascades, Test197SubagentHasNoSelfApprovalOutlet, Test197StreamKeyShapeIsLiteral, Test221TaskCancelIsRegisteredAtItsFrozenLevel, Test221DeferredMarkerForCancelLiftedButListStillMarked, Test221SpawnDescriptionPromisesOnlyWhatIsTrue, Test221ParentStopsItsOwnChildRowAndStreamSettle, Test221SubagentCannotStopSiblingOrItself, Test221ParentCancellationStillDoesNotCascade, Test221TaskCancelUnderNoGateStopsAtTheWindow, Test221EveryPromisedTaskNameIsRegistered, TestEveryRegisteredToolIsClassifiedForMarking175r2

单跑八枚（`logs/t221-baseline.txt`）：rc=0 ／ 8 枚全 PASS。

## 突变（12 发，`mutations/*.json`，逐日志 `logs/mut-<名>.txt`）

尺：`go test ./internal/tools -run 'Test221|Test197|TestEveryRegisteredToolIsClassifiedForMarking175r2' -v -count=1 -overlay <json>`

| 发 | 改动 | rc | 结果行 | 红 | 红名册 |
|---|---|---|---|---|---|
| m1-cancel-unregistered | 摘 `task.cancel` 注册行 | 1 | 20 | 6 | EveryPromisedTaskNameIsRegistered / TaskCancelIsRegisteredAtItsFrozenLevel / ParentStopsItsOwnChildRowAndStreamSettle / SubagentCannotStopSiblingOrItself / ParentCancellationStillDoesNotCascade / UnderNoGateStopsAtTheWindow |
| m2-no-parent-check | 归属判定支改 `if false &&` | 1 | 20 | 1 | SubagentCannotStopSiblingOrItself |
| m3-no-selfstop-check | 自停支改 `if false &&` | 1 | 20 | 1 | SubagentCannotStopSiblingOrItself |
| m4-no-empty-caller-guard | `caller == ""` 支改 `if false &&` | 0 | 20 | 0 | （无尺，缺格） |
| m5-no-nil-roster-guard | `Roster == nil` 支改 `if false &&` | 0 | 20 | 0 | （无尺，缺格＝实现腿自陈的 M7，本腿证实且更宽一枚） |
| m6-cancel-before-refusal | 拒绝判定之前先 `Roster.Cancel(target)` | 1 | 20 | 1 | SubagentCannotStopSiblingOrItself |
| m7-level-l0 | `Declared: risk.L1` → `risk.L0` | 1 | 20 | 2 | TaskCancelIsRegisteredAtItsFrozenLevel / UnderNoGateStopsAtTheWindow |
| m8-old-promise-returned | 乙形①②旧裸许诺放回 Description() | 1 | 20 | 1 | SpawnDescriptionPromisesOnlyWhatIsTrue |
| m9-old-receipt-returned | 乙形③旧回执放回 | 1 | 20 | 1 | ParentCancellationStillDoesNotCascade |
| m10-cascade-returns | `WithCancel(WithoutCancel(ctx))` → `WithCancel(ctx)` | 1 | 20 | 3 | ParentCancellationStillDoesNotCascade / Test197CancelIsPerRowAndNeverCascades / Test197FullPoolRefusesNextSpawnWithReadableReason |
| m11-census-answer-deleted | 删票 175 普查表里 `task.cancel` 那行答案 | 1 | 20 | 1 | TestEveryRegisteredToolIsClassifiedForMarking175r2 |
| m12-finalize-wrong-state | `finalize` 落账写成 `Settled`（丢 `Muted`） | 1 | 20 | 1 | ParentStopsItsOwnChildRowAndStreamSettle |

m1 那发的 AC#1 红句逐字：
「说明书对模型许诺了一枚不存在的工具：task.spawn 的 Description() 写着「task.cancel 只有派生它的那枚父任务能用它单独停孩子——子代理停兄弟、停自己都一律被拒」，而 "task.cancel" 没有注册进这次装配的并集（并集 9 枚：fs.delete / fs.edit / fs.list / fs.move / fs.read / fs.trash / fs.write / task.output / task.spawn）」

### m13 / m14（AC#4 那把尺）＝直接 overlay 不可见，改用"换读路径＋正控"测出结论

直接 overlay `task.go` 的两发（`logs/mut-m13-list-marker-lifted.txt`／`logs/mut-m14-cancel-marked-again.txt`）都是 **rc=0／0 枚红**，原因是尺的实现读盘：
`os.ReadFile("task.go")`（`internal/tools/task_cancel_221_legs_test.go:224`）——`-overlay` 只改编译期视图，不改运行时读到的那块盘。⇒ 这两发不作证据。

换攻法（仍零写工作树）：overlay 只替换**测试文件里的读路径**，让尺的判定逻辑逐字跑在突变体文本上；突变体一律 `git cat-file blob HEAD:` 生成（`mutations/*.go`），
生成器 `gen_mutations.py`／`gen_ac4_mutations.py`，执行器 `run_teeth.py`。

| 发 | 尺读到的是 | 读数（`logs/teeth-*.txt`） |
|---|---|---|
| teeth-control | HEAD `task.go` 原样副本 | **PASS**（装法本身通，非空转） |
| teeth-m13-list-marker-lifted | 摘掉 `task.list` 那行 DEFERRED 标记（不接线） | **PASS＝不响** ⇒ (b) 支无牙 |
| teeth-m14-cancel-marked-again | 把 `task.cancel` 的 DEFERRED 行放回（注册仍在） | **FAIL**：`task_cancel_221_legs_test.go:241`「task.go 里还有 1 行把 task.cancel 标成 DEFERRED…」 ⇒ (a) 支有牙 |

m13 不响的原因（现读对拉）：尺数的是"含 `DEFERRED` 的行里含 `task.list`"，`task.go:279` 的叙述行「…task.list is DEFERRED (§7 :1531)…」同样命中；
HEAD 命中 2 行 → m13 突变体命中 1 行（`grep -n 'DEFERRED' … | grep -c 'task.list'`）⇒ `listMarked == 0` 那一支永远等不到。
teeth 三发跑完 `git status --porcelain -- internal cmd` = ''（工作树零写入）。

## AC#5 整包（本腿自发）

命令：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./cmd/wisp ./internal/...`
日志：`logs/ac5-full.txt`，rc=1（跑到终态）
- 包级：`grep -cP '^ok[[:space:]]'` = 23 ／`grep -cP '^FAIL\t'` = 3（cmd/wisp 118.938s、internal/ball 0.241s、internal/panel 2.982s）／无测试文件 5（合计 31 枚包）
- 名级红：`grep -cE '^[[:space:]]*--- FAIL'` = 6：
  cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues(5.15s, `config_reload_223_test.go:379` "the operator is not told the loosening did not take effect")、
  internal/ball/TestC21TableColourRowsMatchTokensCSS、internal/panel 4 枚
  （TestApprovalCardViewJSONKeysMatchFrontendTypes / TestComposerContractTypesMatchFrontend / TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme / TestC21DesignTokensFourWayAgree）
- panel/ball 那 5 枚的失败正文逐字：`read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified` ⇒ 别人在工作树里删掉的那棵 `design/**`，在册常红，本腿一字未动。
- `internal/tools` 那行：`ok github.com/CarlosShao/wisp/internal/tools 17.435s`。
- 隔离复量（`logs/recheck223-1..3.txt`）：3/3 rc=0（5.055s / 5.019s / 5.034s）⇒ 那枚 cmd/wisp 红判为计时红〔待复量已过〕。

## 门禁三件（终态树）

- `"$GOPATH/bin/gofumpt.exe" -l internal/tools/ cmd/wisp/` = 空（rc=0）
- `go vet ./internal/tools/ ./cmd/wisp/` = rc=0
- `bash scripts/d22scan.sh` = rc=0，`logs/d22scan.txt` 尾行 = "clean - no D22 ban violations; live scope work: bans #1-5 internal/=219, bans #1-5 cmd/=29, ban #6 frontend/=85, ban #7 internal/tools/=22, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=462, ban #8 cmd/=63"

## 现读点尺（AC#2 与 AC#4）

- `grep -rn '\.Cancel(' --include='*.go' internal/ cmd/ | grep -v _test` = 12 枚命中，逐枚读后名册那枚的调用者 = 1（`internal/tools/task.go:740`），其余 11 枚为 `context.CancelFunc` / `RunningTask.Cancel` / `replyRoot` / `reloadRoot`；**命中里没有一行是注释**。
- `grep -n 'DEFERRED' internal/tools/task.go` = 3 行（`:23` task.list 标记、`:32`/`:279` 叙述），无一含 `task.cancel`。
- `grep -c 'DEFERRED(D-'` 在本票 6 枚写面上全部 = 0。
- `git diff a7993a9b..HEAD --name-only -- internal cmd` = 6 枚；`-- docs/specs docs/PLAN.md` = 0 行。
- `git show c44b30c4 -- internal/tools/task.go | grep -c '^-.*task.cancel.*DEFERRED'` = 1；`… | grep -c '^+.*{Tool: taskCancel{d: d}'` = 1（同发）。
- `git show --stat 453eebab` = 票面 31 行 + probe plan 82 行，零产码。
- `grep '^-' run.go diff | grep -vcP '^-\s*//'` = 0；`grep '^+' … | grep -vP '^\+\s*//'` = 空。
- 名册方法枚举（`grep -n 'func (r \*TaskRoster)' internal/tools/task.go`）= Record / WatchRow / Look / Count / Descendants / PublishSubagent / MarkRoot / TryAcquireSubagentSlot / InFlightSubagents / RunningSubagentIDs / AttachCancel / DetachCancel / Cancel ⇒ **无删行 API**。
- 生产审计接线：`cmd/wisp/run.go:584` `Logf: rt.auditf`、`:571` `Journal: mem`；sink 实现在 `internal/tools/bridge.go:992`（每发调用一行）与 `:1030`（落库 `memory.ToolCall{TaskID, Tool, ArgsJSON, …}`）。
- 孩子目录可见性：`internal/tools/subagent_197.go:549` `const subagentHiddenTool = "task.spawn"`，`:552-567` 只按这一枚名字过滤 ⇒ 孩子目录里有 `task.cancel`。
- 任务 id 铸造：`internal/agent/loop.go:1123 newTaskID()`（crypto/rand，`rand.Read` 失败退化成 `time.Now().UnixNano()` 派生）；`Loop.Run`/`RunAsync`（`:332`/`:321`）不收调用方 id；`internal/tools/task_backfill.go:38` 明写 call id 是 model-supplied；`grep -rn newTaskID --include='*_test.go'` 只命中 `internal/agent/compress_trace_test.go:524-527`（形状 36 字符/4 短横）与 `:557`（两次铸造不相等）。
