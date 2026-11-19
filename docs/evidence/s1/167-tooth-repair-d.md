# 167 修牙腿（167-d）——审计点名「构造性为真」的 seq 不重排牙，改成能红的形状（2026-11-19 下午）

- 审计出处：票 167 面上「只读审计腿复核（2026-11-19 14:5x）」节（12 枚提交＝8 通过 / 3 可疑 / 1 不通过）；不通过＝`427dd7955`。派单存档＝`.scratch/wisp/probes/orch/20261119/txt/167-d-tooth-repair-dispatch.txt`（commit `8e9803858`）。
- 审计指控（原文两点）：①它删掉了队列自己的 admission 号断言，把命题换到另一枚铸号器上；②新牙 `Test167QueueSeqStableAcrossCancel` 的 seq 取自**测试自建的泵**＋`QueueSeqLedger` **首次见即铸号** ⇒ **构造性为真、无法红**——「每次 Snapshot 重排号」的账本与真账本在该测试下不可区分。
- 本腿不做裁决（裁决者≠实现者）；本件是实现腿的证据件，读数全部本腿现跑，出处逐条给命令与日志。

## 1. 修法：把铸号与遍历次序分离（可分离的形状）

- `internal/panel/pump_queue_seq_167.go`：`Assign(correlationID)`（一次调用＝首见即铸，命题与遍历融合）拆成两段：
  - `castAll(ids []string)`——泵把**一帧快照 pending 段的整组 correlation id**交给账本（交整组的边界＝这一次调用）；未见的 id 按该组呈现次序铸下一个单调号；**已见的 id 无论在本次组里排第几，保号不动**（这一句就是「取消不重排」命题本身，现在有了独立行使的入口）。
  - `number(correlationID)`——按 id 取回它被铸的号；空 id／未铸 id 取回 0（卡不带 seq 键，omitempty，与原形一致）。
  - nil 账本两段调用均安全（castAll 不做事、number 返 0），与原「nil ⇒ 不铸号」口径一致。
- `internal/panel/pump.go` `Snapshot()`：逐卡 `Assign` 改为先 `castAll(整组 ids)` 再逐卡 `number`；出向键、omitempty、`CardView()` 路径零变化（panel 全包回归 ok，27.332s，见 §4）。
- ⛔ 未碰 `internal/agent/approval/**`（提交后 diff 计数见 §6）；⛔ 未放宽任何断言（旧牙逐条断言原样保留，见 §3 发二逐枚 PASS 名单）；⛔ 未 push；⛔ 无 UI emoji。

## 2. 新用例（三枚，各钉各的命题）

| 用例 | 包 | 钉什么 | 读法 |
|---|---|---|---|
| `Test167SeqReadTwiceFromTheOutboundPayload` | `cmd/wisp` | **不经 `QueueSeqLedger` 任何私有字段**：从泵出向载荷（`pump.Snapshot().Pending`）读两次——取消第 2 条前 `corr3→seq A`；取消后 `corr3→seq A`（同号）且 `corr1` 号不变。两次读数都取自对外载荷 | 端到端：真队列 `approval.New`→三张真 L2 卡→`DecideFromNative` 取消→同一泵两次 Snapshot |
| `Test167MutationRenumberOnEverySnapshot` | `cmd/wisp` | **变异可红载体牙**：突变态（`castAll` 每次重排 1..n 按当前顺序）下端到端读数必须红，点名「card 3 seq 从 3 变 2」；还原态＝普通绿牙 | 突变经 `GOFLAGS=-overlay` 注入（先例＝`docs/evidence/s1/219-ac5-reason-findback.md` §8），盘上树恒为修复后形状 |
| `Test167QueueSeqLedgerNumberSurvivesSetShrinkage` | `internal/panel` | **ledger 级命题牙**（修复拆缝后的直接行使面）：`castAll([1,2,3])` 后 `castAll([1,3])`（集合收缩＝取消的账面形状），幸存者保号（corr-3 仍 3 不变 2）＋号不回收（新 id 铸 4）＋空 id 不铸＋nil 账本安全 | 包内直测，不依赖泵 |

## 3. 两发读数（先证红，再证绿；载具一次一发，每发跑前 `Get-Process wisp.test`＝0）

载具＝`scripts/wisp-cli-staged-tests.sh`（经 `D:\work\soft\Git\bin\bash.exe`，`bash -c '<整条>'` 形传参）。突变体＝`%TEMP%\wisp167d\pump_queue_seq_167.mut.go`（`castAll` 每次清簿重排 1..n）＋`%TEMP%\wisp167d\overlay.json`（仅替换 `internal/panel/pump_queue_seq_167.go` 一枚映射）；overlay vet 双向 rc=0（挂上 0／还原 0）。盘上树全程只有本腿写面文件（`git status --porcelain -- internal/panel cmd/wisp docs/evidence` 复核）。

### 发一（证红）：挂 overlay 突变，单跑变异用例 ⇒ 红，读数逐字点名

```
GOFLAGS="-overlay=C:/Users/swq/AppData/Local/Temp/wisp167d/overlay.json" scripts/wisp-cli-staged-tests.sh -test.run "Test167MutationRenumberOnEverySnapshot" -test.count=1
```

- `=== RUN   Test167MutationRenumberOnEverySnapshot`
- 红句逐字：`panel_occupancy_queue_167_windows_test.go:357: card 3 payload seq after the cancel = 2, want 3 (the mutation made it fall to 2: 取消不重排 is broken and this tooth catches it)`
- `--- FAIL: Test167MutationRenumberOnEverySnapshot (0.02s)`；`FAIL github.com/CarlosShao/wisp/cmd/wisp`；**载具 rc=1**（15:57:57）。日志＝`.scratch/wisp/probes/orch/20261119/167d-red-renumber-mut.log`。
- **即派单点名的读数：card 3 的 seq 从 3 变成 2。**

### 发二（还原证绿）：同一 overlay 挂着、不带突变环境（盘上树本就是修复形状），整族 `^Test167` ⇒ 全绿

```
scripts/wisp-cli-staged-tests.sh -test.run "Test167" -test.count=1
```

- **13 枚全 PASS / 0 FAIL / 0 SKIP**（逐枚）：AC6DoctorFinds 0.00s／AC6QuietWhenNoCrashOnFile 0.00s／AC6NewestFileWins 0.00s／AC6DoctorCarries 0.00s／AC5FailedSendKeeps 0.00s／AC5FailureTouches 0.00s／AC5WireShape 0.00s／AC5FailedSendFaces 0.00s／**QueueSeqStableAcrossCancel 0.02s**／**SeqReadTwiceFromTheOutboundPayload 0.02s**／**MutationRenumberOnEverySnapshot 0.02s**／OccupancyReader 0.00s／AC4BusyStop 2.96s。
- `ok github.com/CarlosShao/wisp/cmd/wisp 3.00s`；**载具 rc=0**（15:59 前后）。日志＝`.scratch/wisp/probes/orch/20261119/167d-green-restored-full167.log`。
- 旧牙 `Test167QueueSeqStableAcrossCancel` 全部原断言（取消前铸号 1/2/3、card3 seq==3、pos==2、live=2 不丢条、wire `seq:3`/`position:2`）原样在、原样绿——**没有一条被放宽**。

### 支撑两发（ledger 级命题牙，`internal/panel` 直跑）

- 突变态：`GOFLAGS=-overlay=… go test ./internal/panel -run '^Test167' -count=1 -v` ⇒ `--- FAIL: Test167QueueSeqLedgerNumberSurvivesSetShrinkage (0.00s)`，红句逐字：`pump_queue_seq_167_test.go:35: corr-3 after the shrink = 2, want 3 (NOT renumbered to 2 by its new place …)`；**rc=1**（日志见 `.scratch/wisp/probes/orch/20261119/167d-panel-mut-red.log`）。
- 还原态：同命令 ⇒ `--- PASS`，`ok github.com/CarlosShao/wisp/internal/panel 0.050s`，**rc=0**（`.scratch/wisp/probes/orch/20261119/167d-panel-restored-green.log`）。

### 过程诚实笔（中间形读数，留档不冒充终态）

- 第一版变异用例带了 `WISP167D_RENUMBER` 环境门：还原发里该用例 `--- SKIP`。**skip≠绿**，判不合格，当场去掉门（用例成为常驻牙：还原态跑＝绿、突变态跑＝红），两发终态读数按 §3 重取。中间发日志留档＝`.scratch/wisp/probes/orch/20261119/167d-mut-active-full167.log`（突变仍挂、整族 `^Test167`：变异用例红 0.02s，**且旧端到端牙 `Test167QueueSeqStableAcrossCancel` 也红了**——`card 3 after the cancel = seq 2 pos 2, want seq 3`＋`wire seq = 2`，其余 11 枚 PASS）。
- **这一发顺带是审计指控被修复的直接反证**：改形前（审计口径），「每次重排」的账本对泵驱动测试不可区分；改形后，同一枚变异让**两枚牙**（新变异牙＋旧端到端牙）当场红。修复没有靠新增门禁开关、没有靠 skip。

## 4. 三样现跑（本腿自己的门禁读数）

1. **d22scan**：`sh scripts/d22scan.sh` ⇒ **0 findings**，clean 行逐字＝`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=422, bans #1-5 cmd/=107, ban #6 frontend/=187, ban #7 internal/tools/=52, ban #10 internal/risk/=21, ban #8 design/=39, ban #8 frontend/=187, ban #8 internal/=1066, ban #8 cmd/=359; ban #8 emoji coverage: design/ 39 text files; frontend/ 187 text files; internal/ 1066 Go files, comments and _test.go included; cmd/ 359 Go files, comments and _test.go included`，末行逐字＝`d22scan.sh: OK - control rc=0, scan rc=0`，rc=0（首跑曾红一枚 phantom-citation＝注释先引了尚未建的本证据件；证据件落盘后复跑即绿，首跑读数留档 `.scratch/wisp/probes/orch/20261119/167d-d22scan-first-red.log`，终跑＝`167d-d22scan-final.log`）。
2. **载具 `^Test167`**：§3 发二，13 枚全 PASS，**rc=0**（含两枚新用例）。
3. **`go build ./...`**：rc=0（读数见 §5）。

## 5. 逐字读数汇总（回话第一行用的那几个数）

- 变异红（载具，15:57:57）：`--- FAIL: Test167MutationRenumberOnEverySnapshot (0.02s)`，`card 3 payload seq after the cancel = 2, want 3`，rc=1。
- 还原绿（载具）：`ok github.com/CarlosShao/wisp/cmd/wisp 3.00s`，13 枚全 PASS（含新用例），rc=0。
- panel 变异红→还原绿：`NumberSurvivesSetShrinkage` FAIL（corr-3 = 2 want 3）rc=1 ⇒ PASS ok 0.050s rc=0。
- panel 全包回归：`ok github.com/CarlosShao/wisp/internal/panel 27.332s`（字节钉全绿）。
- d22scan 末行：`d22scan.sh: OK - both verdicts green (control rc=0, scan rc=0)`。
- `go build ./...`：rc=0。

## 6. 禁区与口径（如实写）

- **`internal/agent/approval/**` 零字节**：本腿 diff 零枚该包文件（提交后 `git show --name-only <本腿> | grep -c agent/approval` ＝ 0，读数随 commit 更新在回话里给）。AC#7 零字节轴仍有效，本腿未申请解冻。
- **如实写这一句（派单第 4 条）**：`internal/agent/approval/queue.go:165 q.seq++` 的 admission 不重排性质**因 AC#7 零字节轴在包内零测试覆盖**（这是轴要求的取舍，不是放宽——该包不得为本票加一行测试）；现在对外可见的号由面板侧铸（本腿拆缝后的 `castAll`/`number`），其不重排由本腿两发用例守：出向载荷两读数（`Test167SeqReadTwiceFromTheOutboundPayload`）＋变异可红载体（`Test167MutationRenumberOnEverySnapshot`），panel 包内另有 ledger 级命题牙（`Test167QueueSeqLedgerNumberSurvivesSetShrinkage`）。
- 突变体只在 `%TEMP%\wisp167d\`，⛔ 不入库；盘上树除本腿写面外零改动；⛔ 未 push、未 amend/reset/rebase/stash/checkout/clean；`git add`/`commit` 显式 pathspec。
- 本腿写面：`internal/panel/pump_queue_seq_167.go`＋`internal/panel/pump.go`（铸号点两段化，票 145 已批 composer/pump 族内）＋`internal/panel/pump_queue_seq_167_test.go`（新增）＋`cmd/wisp/panel_occupancy_queue_167_windows_test.go`（两枚新用例＋头注）＋本证据件＋`.scratch/wisp/probes/orch/20261119/167d-*.log`（读数存档）。
