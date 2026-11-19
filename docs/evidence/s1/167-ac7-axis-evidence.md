# 167 AC#7 契约轴实证 — `internal/agent/approval/**` 归零＋Seq 铸号回 panel 侧（写腿 agent-144，2026-11-19）

- 票：`.scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md` **AC#7**
- 派单：`.scratch/wisp/probes/orch/20261119/txt/167-c-ac7-axis-dispatch.txt`
- 锚点：`e1049dbf5`（票 167 编排者收 AC#2＋AC#3 那笔，`e1049dbf5~1` ＝ 本票任何 Go 改动之前）。上一腿 `2c43ec586`/`6b5b0ccb1` 在 `internal/agent/approval/pending_read.go` 加了 `Seq` 字段＋读数（11 行新增），**越了 AC#7 的零字节轴**——本腿把它逐字节撤掉，并把"取消不重排"那枚稳定显示序号改在 panel 侧自铸。
- 写面：`internal/panel/**`（票 145 已批局部解冻内：composer/pump 族＋一户新文件）、`cmd/wisp/panel_pump.go`（唯一）、`docs/evidence/s1/167-ac7-axis-*`；**`internal/agent/approval/**` 目标＝零改动（达成，见 §1）**；`frontend/**` 零字节

## 1. 撤那一处 — 逐字节还原，blob 级实证

- 还原方法：派单允许的 `git show e1049dbf5~1:internal/agent/approval/pending_read.go > 该文件`，经 `cmd /c` 重定向（PowerShell 管道会改编码，`cmd` 重定向是字节直写）。⛔ 没碰任何别的文件——没有 restore/checkout。
- **blob 级判据（现跑）**：`git rev-parse e1049dbf5~1:internal/agent/approval/pending_read.go` ＝ `ef71c24512e10ba69037a818f376e0db8fa60f2d`；`git hash-object` 工作树文件 ＝ **同一枚** `ef71c24512e10ba69037a818f376e0db8fa60f2d` ⇒ 逐字节相等。
- `git diff e1049dbf5~1..HEAD -- internal/agent/approval/` 提交后＝**空读数（0 枚文件）**——AC#7 轴闭合，§3 的 diff 实测为准。

## 2. Seq 改在 panel 侧铸 — 机制（不造第二条队列、不复制审批逻辑）

- **新枚** `internal/panel/pump_queue_seq_167.go`：`QueueSeqLedger`——一张 `correlationId -> seq` 稳定映射＋一枚单调计数器。`Assign(correlationID)`：首次见到某 id ⇒ 铸下一个单调号并记下；再次见到 ⇒ 返回同一个号；空 id ⇒ 铸 0（卡不带 seq 键，omitempty 吃掉——不给无 id 的行编造号码）。互斥锁一枚（`sync.Mutex`）。
- **为什么这不是第二条队列**：ledger 不存条目、不判任何事、不给任何审批放行，没有任何 route 进/出审批面；它只是泵旁边的编号簿，和 `StreamLog` 给流行记键是同一类记账。审批逻辑（入队/取消/过期/放行）全部还在 `internal/agent/approval` 原处，一字未动。
- **position 仍用现源**：host 映射 `nativeVerdictsFromGate` 照旧逐字读队列自己的 `Position`（1-based FIFO 现位），本腿没动 position 的任何来源。
- **铸号点在泵里，一泵一簿**：`internal/panel/pump.go` `Snapshot()` 组 pending 卡时逐卡 `card.Seq = p.seqs.Assign(v.CorrelationID)`；`NewSnapshotPump` 各建一枚 ledger。两台生产泵（run 泵＋resident 泵）各自铸各自页面的序号史。`NativeVerdict.CardView()` 本身**不铸号**（CLI 卡路径无泵，seq/position 双 omitempty＝Q-51 的零形状不变）。
- **`NativeVerdict.QueueSeq` 载体字段已删**：撤掉 `LiveApproval.Seq` 后它永久读 0，留着就是一枚假槽——面板的铸号根本不经过它。这是**删**上一腿引入的载体，不是删既有的任何出口。

## 3. AC#7 轴实证 — `git diff --name-only e1049dbf5~1 <commit>` 逐面实测

> 本节读数＝提交前对暂存树现跑（`git diff --name-only e1049dbf5~1`，与提交后的 diff 逐枚相同；总 33 枚，其中上一腿占 28＋台账 1，本腿只落 5＋证据件 1）：

- `internal/agent/approval/**`：**0 枚**（逐字目标达成）
- `docs/PLAN.md`：0 枚 · `docs/specs/**`：0 枚 · `internal/risk/**`：0 枚 · `thresholds.go`：0 枚 · golden：0 枚 · `allowlist.txt`：0 枚 · `design/**`：0 枚（该面 grep 计数＝0）
- `internal/panel/**` 全 6 枚（**逐枚列出**；前 4 枚＝上一腿的、本腿零字节碰，后 2 枚＝本腿）：
  1. `internal/panel/approval.go`（上一腿加的 `seq`/`position` 两键＋注释；本腿只改注释语义——键集、tag、类型零改动）
  2. `internal/panel/composer.go`（上一腿：occupancy/failedSend 四键，纯新增字段；本腿零字节碰）
  3. `internal/panel/composer_test.go`（上一腿：两枚键的 in-册 注册；本腿零字节碰）
  4. `internal/panel/occupancy_queue_167_test.go`（上一腿建的测试文件；本腿改测法）
  5. `internal/panel/pump.go`（本腿：`SnapshotPump` 新增私有 `seqs` 槽＋`NewSnapshotPump` 建簿＋`Snapshot()` 逐卡铸 seq——**新增字段/私有槽**，无删既有字段、无新增出向键、无新增方法名；`NativeVerdict` 撤掉上一腿的 `QueueSeq` 载体）
  6. `internal/panel/pump_queue_seq_167.go`（本腿新增，唯一新文件——`QueueSeqLedger` 原语，票 145 已批 composer/pump 族内）
- `cmd/wisp/**` 全 13 枚（共享编译面逐枚；前 12 枚＝上一两腿的，`panel_pump.go` 同时含本腿唯一改动——撤 `QueueSeq: it.Seq` 一行＋注释改写，只减不增键；`panel_occupancy_queue_167_windows_test.go` 含本腿改测法）：`crash_selfrescue_167.go`/`_test.go`、`doctor.go`、`panel_failedsend_167.go`/`_test.go`、`panel_message_453.go`/`_windows.go`、`panel_occupancy_queue_167_windows_test.go`、`panel_pump.go`、`panel_resident_pump.go`、`panel_stop_tail_167_test.go`、`resident_windows.go`、`run.go`
- docs 4 枚：两腿证据件＋台账追加＋本件。

## 4. 牙与读数 — 三样现跑（本腿会话内，2026-11-19）

1. **d22scan**（`./scripts/d22scan.sh`，含正控自测）：**0 findings**，末行 `d22scan.sh: OK - control rc=0, scan rc=0`，rc=0。
2. **载具**（`scripts/wisp-cli-staged-tests.sh`，串行一次一发，跑前 `Get-Process wisp.test` 确认空车道；bash 用 `D:\work\soft\Git\bin\bash.exe`）：
   - `-test.run '^Test167'` ⇒ **11 RUN / 11 PASS / 0 FAIL**，`ok github.com/CarlosShao/wisp/cmd/wisp 4.00s`，**rc=0**。含 `Test167QueueSeqStableAcrossCancel`（取消第 2 条后第 3 条仍是 3——panel 铸号版，卡 3 的 seq=3、pos=2 双断言全在）。
   - 邻族回归（同载具一发）：`TestRosterCancelDoor|Test453|Test192MessageDoor|Test214|Test38ForceChip|TestTaskCardFields|Test65ResidentPump|TestThePumpIsDriven|TestAC8SnapshotClamp|Test148|Test48C18|TestAC7TornPacket|TestAC7PumpLock|TestFrontendHasNoEmoji|TestApprovalCardView` ⇒ 全 PASS，rc=0（读数见回话）。
3. **`go build ./...`** ⇒ **rc=0**（改后先跑，提交前复跑一遍读数为准）。
4. **panel**（`go test ./internal/panel -run '^Test167|^TestOccupancy' -count=1`）⇒ `ok ... 0.062s` rc=0；`TestQueueNumbersTravelAndZeroStaysHome` 单发 −v ⇒ PASS（正向：泵铸 seq=1/位=2 上线；重建同号=1；无 id ⇒ 双键缺席——**断言未放宽**：缺席方向从"队列没编过号"改为"无 id 可铸"，都还钉 0/缺席形状）。
5. **整包复跑**：`go test ./internal/agent/approval -count=1` ⇒ `ok ... 5.952s` rc=0（**还原后的 approval 包原样全绿**）；`go test ./internal/panel -count=1` ⇒ `ok ... 23.935s` rc=0（reconciliation 两把尺带在册键全绿）。

## 5. 断言未放宽 — 逐枚对照

- `Test167QueueSeqStableAcrossCancel`：**原有逐条全保留**（深度 3→取消→live=2 不丢条；card 3 seq==3 不重排；pos==2 会动；card 1 1/1；wire `seq:3`/`position:2` 在线）。**新增**（更强，不是更弱）：取消前先验三卡铸号 1/2/3；同一泵重建 seq 不变（重排即红的正向钉）。读法从 `LiveApprovals().Seq` 改为泵出向（机制搬家，被测命题不变：**取消第 2 条后第 3 条仍是 3**）。
- `TestQueueNumbersTravelAndZeroStaysHome`：双向钉都还在——带 id 卡 seq/position 都上线（值 1/2），无 id 卡双键缺席；**加了**重建同号钉。没有一条断言被删或放宽。
- `TestOccupancy*` 三发：本腿零字节碰，读数照旧。
- 变异（钉的证明）：撤掉泵里的 `card.Seq = p.seqs.Assign(...)` 铸号行即红（wire seq 缺席/重排钉都咬）；未单独跑变异发，钉句在测试正文可逐字读。

## 6. 边界（本腿没做、也不许做的）

- `internal/agent/approval/**`：**0 字节改动**（还原后 blob 相等实证）。队列自己的 admission 计数器（`qitem.Seq`，queue.go 内部）与 `LiveApprovals` 白名单面全部保持原状——panel 的 ledger 不是它的替身，也不碰它的状态。
- `frontend/**` 零字节（composer_test.go 在册的 judged page-side gap 原样：页面侧声明仍是界面腿的活）。
- ⛔ thresholds/golden/allowlist/PLAN/specs/D43/design：0 字节；⛔ 放宽断言：无；⛔ push：未推；⛔ UI emoji：d22scan clean。
