# 票 271 — task 家族枚数尺选形（甲）落地凭据 · leg 271-v1

- HEAD 起手 `ac2c1d903`（派单写 `edab04ab6`，共享树上其他腿已推进；锚点行全部现量核对，见 §2）
- 写面：`internal/tools/taskfamily_roster_271_tripwire_test.go`（新）＋ `internal/tools/task_cancel_221_legs_test.go`（纯注释与文案，判据一字未动）＋ 本件
- ⛔ 未动：`docs/reports/pending-and-issues.md`（编排者台账）、`task.go`（生产面一字未动，含 `:23` 标记与 `:279` 注释）、未注册 `task.list`

## 1. AC#0 选形＝甲（升格成明写的警报线），理由按铁律判据

| 形 | 动哪行 | "新增第三枚那天接手的人会看到几枚红" | 顶红现有钉 | 本腿裁决 |
|---|---|---|---|---|
| 甲 | `task_cancel_221_legs_test.go` 注释＋Errorf 文案＋新正控文件；判据 `!= 2` 一字不动 | 1 枚红＝哨兵本身，红句带工序说明书，且该红被 resident 正控证明"假红但设计内" | 0 枚 | **选它** |
| 乙 | 改 `:184` 判据为名集形（须具名解冻，AC#4） | 0 枚红——但"家族变大"维度从此零传感（271-a1 普查 §1：全仓数枚数的只有这一把；236r3 名集钉只认 task.list 一枚负向名） | 0 枚代码钉，但动 `-done` 票的尺，程序上无书面解冻先例 | 不选 |
| 丙 | 不动（A648 已裁过一次） | 同甲的红，但红句无说明书、M-2c 读数永远停在转述 | 0 枚 | 不选（重开条件已触发） |

**为什么甲最有价值（铁律 2 判据：最有价值＋主流 harness，非动作最小）**：

1. **A648 选丙的唯一理由是"没真跑过"**——票面 §Progress log 19:5x 那条：三张代价表的几枚红"全部是转述"。本腿把那次裁决里唯一承重的未跑读数（236r3 交件头注亲口承认的 M-2c："NOT measured by r3c ... 假红预测停在 argument from source, not a reading"）补成 **resident 实测**：`Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire` 每次跑套件都重证"B-3 台件下枚数哨兵红（红句具名）＋名集钉绿"两半。丙给不出这个，乙也不用给（它直接把哨兵删了）。**甲是三形里唯一能把"假红那发要能复现"（AC#1）变成常驻判据的**。
2. **乙买到的东西名册里已经有**：`Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment`（236r3 AC#2 钉）就是名集形、亲口注明 "Deliberately NOT a count, so a legitimately approved third task tool cannot false-red"。把 `:184` 也改成名集＝重复覆盖＋削掉唯一枚数维传感器＋动别人 `-done` 票的尺。乙是"用唯一传感器换一份重复件"。
3. **主流 harness 形态**：刻意哨兵＋失败文案即文档＋**正控证明警报是活的**（dead tripwire 反模式＝从不验证警报会响）。甲正是这个形状；乙是"把警报器拆了换一个已存在的警报器"。
4. 程序代价：甲零判据改动 → AC#4 具名解冻**不触发**；乙要先落台账解冻五件套（台账归编排者，本腿禁区）。

## 2. 现量核对（派单锚点全部现跑）

- 尺本体：`internal/tools/task_cancel_221_legs_test.go` 起手 `:183-185`（本轮加注释后 Errorf 漂至 `:196`；判据 `if len(names) != 2` 逐字未动）
- 名册：`internal/tools/task.go:600-605`，`BuiltinTaskEntries` 现量 `:602`/`:603` 两行（`task.output`/`task.cancel`）——**一字未动**
- 生产调用者：`cmd/wisp/run.go:544`（未动）
- 236r3 名集钉：`internal/tools/tasklist_deferred_236r3_teeth_test.go:150`（未动；本腿复用其 `registeredTaskNames236r3` 分母助手与其测试函数做还原自证）

## 3. 本腿三枚 resident 件（全绿）

1. `Test271B3UnrelatedThirdToolTripsOnlyTheCountTripwire` —— **AC#1 的 B-3 台件，实测两半**：
   - half 1（警报是活的）：overlay 台件（`task.go` 名册追加一枚无关 `task.note`，仅 `Name()`，236-v1 B-3 同形）下 `Test221TaskCancelIsRegisteredAtItsFrozenLevel` **FAIL**，红句逐字含 `task 家族注册了 3 枚` 与 `task.note`（读数可归属到哨兵本身）；
   - half 2（红是假红）：同一台件下 `Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment` **PASS**，其 Logf 报出现名册含 task.note。
   - 台件机制：`//go:embed task.go` 字面拷贝做锚点手术（锚点命中数≠1 即 t.Fatal，拒绝在漂移锚点上盲切）→ temp 目录 → `go test -overlay`；**共享树零编辑**。
2. `Test271OverlayOnlyDisciplineRosterBackAtBaseline` —— **AC#1 还原自证，resident**：(a) 盘上 task.go 字节==编译期字节（有人原地动刀即红）；(b) 真名集钉（直调其测试函数，不复制谓词）在未 overlay 名册上通过；(c) 起手 2 枚只 Logf 记录不断言——本腿不新造第二把枚数红（AC#2）。
3. `task_cancel_221_legs_test.go` 注释块＋Errorf 文案 —— **甲形的说明书**：响了先分清红型（task.list＝真伤找 236r3；别的第三枚＝回票 271 重选形），⛔ 不许就地放宽。**判据、 Fatalf、其余断言一字未动**（gofumpt 前后 md5 差＝注释＋文案两处，`git diff` 可查）。

## 4. AC#3 门禁四数＋还原自证（全绿，本腿 2026-11-06 11:0x-11:1x 现跑）

| 门 | 读数 |
|---|---|
| `go vet ./internal/tools/` | rc=0 |
| `go build ./cmd/... ./internal/... ./tools/...` | rc=0 |
| `go test ./internal/tools -count=1 -timeout 420s` | 起手 `ok 23.8s` → 终态 `ok 35.8s`（增量＝本腿两枚 overlay 件 ≈10s）；**逐名红名册起手＝终态＝空**，作差＝0 |
| `gofmt -l` / `gofumpt -l`（两枚被碰文件） | 零输出 |
| `sh scripts/d22scan.sh` | rc=0（中途真红过一发：Errorf 字符串里的 `⛔` U+26D4 命中 ban#8"字符串不豁免"——已改为纯文字，注释里的符号豁免不动；这就是 d22scan 该抓的抓、修完即绿的完整轨迹） |
| `git status --porcelain -- internal cmd` | 起手＝3 枚未跟踪在途件（`cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`、`internal/tools/ticket23_image_channel_test.go`）；终态 commit 后应原样＝起手（本腿三件入库，未碰他人物件） |

## 5. 诚实边界（本腿未做到/未验证的）

1. **这枚哨兵在 CI 历史上有没有真响过**＝仍然一条证据也没有（271-a1 §4 第①条，本腿同样取不到 CI 历史）；本腿补的是"它现在响不响、响得对不对"。
2. **批准维零尺**＝三形都没买到，甲也不例合法性的代码验证——"第三枚的票是不是真的"仍是流程（台账）辖区。
3. **AC#0 的裁决权归属**：票面写"由非实现者出裁定表"，本派单明令"AS THE LEG evaluate … write your pick + why in evidence"——按派单执行（派单＝编排者当场指令）；本件即裁定表，编排者验收时可按 `SPEC-12 §4.3` 另派对抗验收。
4. **台账动作归编排者**（本腿禁区）：建议按甲形登记——重开条件已由本派单触发并落地；撤销口令＝"回退本 commit（三文件）即回 ac2c1d903 原样，哨兵回到无说明书状态，甲形说明书在 git 历史可取回"。

## 6. 交付物

- `internal/tools/taskfamily_roster_271_tripwire_test.go`（新，resident 正控＋还原自证）
- `internal/tools/task_cancel_221_legs_test.go`（注释＋文案，判据一字未动）
- 本件
- commit：本腿只 commit 不 push（编排者核后推 origin/dev）。
