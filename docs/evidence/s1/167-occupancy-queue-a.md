# 167 片A — 占用"未知"具名态 ＋ 排队序号/当前第几条（写腿交件，2026-11-19 10:5x）

- 派单：`.scratch/wisp/probes/orch/20261119/txt/167-a-occupancy-queue-dispatch.txt`（同刻入库＝commit `b3562a37c`，存档即派）。
- 起手锚：`b3562a37c`（派单存档笔；其父 `c6793d560`）。写面全部落在派单许可面内：`internal/agent/approval/pending_read.go`（只读出口）、`internal/panel/**`（carrier＋尺＋测试）、`cmd/wisp/**`（装配接线＋测试）、本件。
- ⛔ 本腿零字节未落 `frontend/**`（两个新卡键 seq/position 与三枚占用键的 page 半边走 `composerTSGapLedger` 登记在册，见 §5）；⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design（D32 只读）；⛔ 未动他人在飞件（`composer-states.html`／`44-costmeter-c23-done.md`／`291` 测试／`346-a2-dispatch.txt` 四枚 ` M` 本腿全程未碰，收口 `git status` 复认仍在）。

## 1. 先量后写：四枚今天的出口现量（派单点名的那把 grep）

尺（现跑）：`grep -rn "occupancy|context.*token|ContextWindow|queued|queue" --include=*.go internal cmd | grep -v _test | head -30` 等价现量（pwsh Select-String 同词面、非测试）＋逐枚追到函数体。四枚判语：

| 枚 | 判语 | 凭据（全部现量，带 file:line） |
|---|---|---|
| 占用（分子） | **有源、未接出口** | `internal/panel/run_feed_145.go:277` `RunFeed.Usage`（Cost 累计口径，`internal/agent/cost.go:52` `AddUsage` 注释逐字 "Counters SUM across turns"）；喂点 `cmd/wisp/run.go:2928`（EvUsage → runFeed/taskLive/stream 三读者） |
| 占用（分母） | **有源、未接出口** | `llm.Endpoint.ContextWindow`（`internal/llm/resolver.go:56`）；loop 的两步读＝`internal/agent/loop.go:237-239`（0 则 `info.MaxContextWindow`，`internal/llm/provider.go:51`）；`internal/agent/budgets.go:86-88` 的 128k 只是缩放基线（⛔ 不当分母编数用） |
| 排队序号（稳定号） | **有源、未接出口** | `internal/agent/approval/queue.go:31-34` `qitem.Seq`＋`:165` `q.seq++` 入队一次铸定，drop 不改他人 Seq |
| 当前第几条 | **有（process 内）** | `queue.go:293-301` `position()`（1-based FIFO 现位）；出向面 `pending_read.go` `LiveApproval.Position`——但卡面 `ApprovalCardView`（`internal/panel/approval.go:39-86`）零枚队列键，出向只到 cmd/wisp 内部就断（`nativeVerdictsFromGate` 未带） |

⇒ **没有一枚"已有完整出口"**：两枚都是"源在、接线断"，故按派单**复用现源接线**，未硬造第二套（Seq 复用 `qitem.Seq`，分母复用 loop 自己的两步读，分子复用 runFeed 已喂的累计行）。

## 2. 落地形状

**AC#3 排队（两枚键，一位不改稳定号）**：
- `pending_read.go` `LiveApproval` 增 `Seq uint64`（`LiveApprovals` 填 `it.Seq`）——纯只读出口，判据面（`ticket148_panel_whitelist_test` 的 Decision 反射尺）不涉：Seq 在 LiveApproval 包装层，与既有 `Position` 同位。
- `internal/panel/pump.go` `NativeVerdict` 增 `QueueSeq`/`Position`，`CardView()` 逐字带上。
- `internal/panel/approval.go` `ApprovalCardView` 增 `seq`/`position` 两枚 `omitempty` 键（Q-51 已答的那条分水岭：无 omitempty 的 int 恒发 0＝CLI 卡面必带 `"position":0` 的退化形；1-based＋稳定号均非零故 omitempty 安全）。
- 接线：`cmd/wisp/panel_pump.go` `nativeVerdictsFromGate` 逐字带 `QueueSeq: it.Seq, Position: it.Position`（两条 pump——run 与 resident——共用此映射，自动同享）。
- 语义分工写死在注释：`seq`＝入队号（取消不重排），`position`＝当前第几条（随队头推进而动）——两枚键各说各的真话。

**AC#2 占用（三枚键，具名未知态）**：
- `internal/panel/composer.go` `ComposerState` 增 `contextTokens`/`contextWindow`（两枚 omitempty int64）＋`contextOccupancyKnown`（恒发 bool＝具名未知态，false 即"未知"枚举字段，对齐 `modelKnown`/`credentialKnown` 家族）。
- `internal/panel/pump.go` 增 `OccupancyReading{UsedTokens, WindowTokens}`＋`PumpSources.Occupancy` reader（reader-gated 家族同一形状）；`Snapshot()` 填充规则＝**分子或分母任一未知 ⇒ 整枚 unknown**（known=false，两数一律不上线，⛔ 0、⛔ 编窗）。
- 接线：`cmd/wisp/panel_pump.go` `contextOccupancy()`＋`contextWindowTokens()`（分母＝`rt.endpoint.ContextWindow`，0 则 `rt.provs[0].Info().MaxContextWindow`，再无则 nil——**与 loop 自己缩放用的两步读同源**；分子＝`rt.runFeed.Section().Usage` 的累计行，in+out，未落 EvUsage 前 nil）；`run.go` PumpSources 挂 `Occupancy: rt.contextOccupancy`。

## 3. 牙（四发，逐发种→红→还原，读数全部本腿现跑）

| 发 | 突变位 | 命中的钉（红句要点） | 结果 |
|---|---|---|---|
| ①分子填0 | `internal/panel/pump.go` Snapshot 占用支：occ==nil ⇒ 填 `ContextTokens=0`＋`Known=true` | `TestOccupancyUnknownNumeratorIsANamedStateNotZero`：`contextOccupancyKnown = true ... want the named unknown state false` | 红 |
| ②分母编数（面板侧） | 同支：WindowTokens<=0 ⇒ 填 128000＋known | `TestOccupancyUnknownDenominatorNeverInventsOne`：三行全红（known=true＋`contextTokens=5200`/`contextWindow=128000` 上线） | 红 |
| ②b分母编数（装配根侧） | `cmd/wisp/panel_pump.go` `contextWindowTokens` 尾部返 128000 | `Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured`：`reader with a usage event but an unknown window = &{UsedTokens:540 WindowTokens:128000}, want nil` | 红 |
| ③取消重排 | `pending_read.go` `Seq: it.Seq` → `Seq: uint64(q.position(it))` | `Test167QueueSeqStableAcrossCancel`：`card 3 Seq after the cancel = 2, want 3`＋`wire seq = 2 ... want 3` | 红 |
| ④取消丢条 | `queue.go` `dropLocked` 改成连删下一条 | 同枚：`live rows after the cancel = 1, want exactly 2` | 红 |

还原后复认：working tree `grep MUTATION 167`（*.go）＝**0 枚**；`git diff` 上四枚产码件无突变行残留。

## 4. 正控（绿读数）

- `internal/panel`：`TestOccupancyBothKnownTravelVerbatim`（wire `contextOccupancyKnown:true`＋`contextTokens:5200`＋`contextWindow:128000` 逐字）／`TestQueueNumbersTravelAndZeroStaysHome`（反向种 2/3 ⇒ 键在值对；未编号 verdict ⇒ wire 零键——Q-51 正向种 0 那发）。
- `cmd/wisp`：`Test167QueueSeqStableAcrossCancel` 全链（真队列 `approval.New` → 三张真 L2 卡经生产入口 `PendingApproval` 停靠 → `DecideFromNative` 取消第 2 条 → 第 3 条 `Seq=3` 不动＋`Position=2`（诚实的当前位次）＋`live=2` 条不少 → pump 映射 → wire `seq:3/position:2`）／`Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured`（无 reader ⇒ known=false 且两键缺席）。
- 两把契约尺（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContract*`）绿：新键走 `composerTSGapLedger` 登记（见 §5），未放宽任何断言。

## 5. 契约轴与在册缺口

- `frontend/src/lib/panel.ts` 一字节未动；Go 侧新键五枚（卡面 `seq`/`position`，composer `contextTokens`/`contextWindow`/`contextOccupancyKnown`）已按 415/422 的既定程序登记进 `internal/panel/composer_test.go` 的 `composerTSGapLedger`（`"ApprovalCardView"`、`"ComposerState"` 两条，各带判语）——登记不是跳过：page 半边一旦声明同名键即 stale-red，届时删登记走普通双向尺。
- ⚠ 既有红（非本腿造成，未动）：`internal/panel` `TestFrontendHasNoEmoji` 红＝`frontend/scripts/render-hiddenpause36.tsx:7` 有 `⛔`（U+26D4）；该件为**已入库**件（commit `3de613364`，HEAD 祖先已核 `merge-base --is-ancestor` rc=0）⇒ 此红在本腿起手前就在 HEAD 上。本票硬约束禁写 `frontend/**`，故只报不动（撤销口令形：把那一行的 `⛔` 改成文字即可复绿，归前端腿）。
- `run_feed_145.go:76-80` 既有注释逐字："pending[].position and context.occupancy are ticket 167's territory"——本腿落的就是这两枚，未新增第三枚 pending[] 键。

## 6. 门禁读数（收口三样＋包相关）

- `sh scripts/d22scan.sh`（11:01 现跑）：**`d22scan.sh: OK - control rc=0, scan rc=0`，0 findings**（positive control 全绿、self-test 46/46）。
- 载具＝`sh scripts/wisp-cli-staged-tests.sh`（先查 `Get-Process wisp.test`＝0；正形 `-test.timeout=90m -test.v`，台件 `.scratch/wisp/probes/orch/20261119/167-a-carrier-full.log`，全程 2096s，rc=1）。⚠ 首跑曾用裸 `go test ./cmd/wisp` 撞默认 10m timeout（0xc0000135 另见 PATH 排障，`scripts/wisp-cli-tests.sh` 注释），已改回载具正形重跑，裸跑读数不作数。**整包红名册 5 枚，逐枚定性**（定向复跑现跑，载具空档一发）：
  - `Test37TimeoutWarnsAt270AutoRejectsAt300ReplayWorksBadgeClears`／`Test39NativeLooseningCardAcceptedAppliesAndAudits`／`Test40ModelProbeFailureCarriesTheRealErrorNeverOK`／`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`——**隔离单跑全绿**（同一发 `-run` 五枚，仅 TestAC4 红）⇒ 整包负载 flake（Test39 红句＝Temp config.toml "being used by another process"；Test40＝死端点探测 66s；与我的 diff 面零交集）。
  - `TestAC4FocusReturnToPriorWindowGap33r5`——隔离也红＝**本机真窗环境既有红**（2026-11-18 整包名册同枚：`probes/orch/20261118/39-deflake-fullpkg2.txt` 唯一红；红句自带锚 `head b3562a37c`；真窗焦点测试，桌面状态相关，非本腿面）。
  - **本腿两枚在载具整包内 PASS**：`Test167QueueSeqStableAcrossCancel`（0.02s）／`Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured`（0.00s）；且既有钉 `panel_pump_test.go:265` 的真实包字节里新键已在线上：卡面 `"seq":1,"position":1`、composer `"contextOccupancyKnown":false`（reader 未接＝具名未知随包发布，两数缺席）。
- `go test ./internal/agent/approval -count=1` ＝ ok 7.362s；`go test ./internal/panel -count=1 -skip TestFrontendHasNoEmoji` ＝ ok 36.647s（skip 仅为 §5 既有红，非放宽）。
- `go build ./...` ＝ **rc=0**。
