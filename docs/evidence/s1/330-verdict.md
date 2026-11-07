# 票 330 裁决表（330-r1 · cmd/wisp G1 对接钉 · 非实现者）

- 日期：2026-11-07 12:1x +08；腿＝330 裁决腿（与 330-r1 写腿、330-r2、330-v1 均非同一腿）。
- 料：票面全文＋`docs/evidence/s1/330-r1-dock.md`＋commit `edad424fe` diff＋`.scratch/wisp/probes/330/r1/**`。
- 锚：读数取当前共享工作树（HEAD `532f936d6`，交付件 `edad424fe` 之后零枚别腿字节改动、tracked 树 clean）。
- ★纪律：本腿只读 `cmd/wisp/**` 与 `internal/**`（两发突变当场 cp 还原），写面只有本件＋票尾 Progress log 一行；框全归编排者，本腿⛔ 翻框。

## 逐格判语（PASS／部分／FAIL ＋ 本腿自己的读数）

| AC | 判语 | 本腿读数 |
|---|---|---|
| AC#0 | **PASS** | 三把尺全真：commit `edad424fe` 在盘（`git show` 8 文件＋764 行）；`docs/evidence/s1/330-r1-dock.md` 30 行在盘；判据核心一发本腿独立重放（见 AC#1）＝第三把尺。 |
| AC#1 | **PASS** | ★G1 那发**本腿亲重放**（⛔ 抄读数）：probe 基线副本 sha 逐位等基线后，`sed 's/plan\.present/plan.presentX/g'` 打 `internal/agent/testdata/golden/plan-present.sse`（`%12` 名变 `plan.presentX`）＋`internal/agent/plan_present_doorbell_165r3_test.go`（`:38` 常量变 `plan.presentX`）、真名册零字节 ⇒ 同命令 `go test ./cmd/wisp -run 'Test330' -count=1` ＝ **FAIL rc=1**，红在 `plan_present_roster_dock_330r1_test.go:128`，红句具名：`G1 对接钉：plan-present.sse 线上名 "plan.presentX" 不在装配出来的真名册（22 枚）里：[clipboard.read … window.manage]`——165-v2 报的那发绿**确已变红**＝本票唯一判据成立。当场 `cp` 基线副本还原：两枚 sha256 逐位回 `7d26fa889b59fc3a…`／`6542773d3d2edcc6…`，`git status --porcelain internal/agent` ＝ 0 行。正向发亲跑：`--- PASS`，rc=0，读数行 `:134`＝`对接读数：线上名 "plan.present" ＝ 真名册（22 枚）里的一行；Resident=true`（与实现腿自报逐字同）。 |
| AC#2 | **PASS** | 尺亲跑：`grep -c '"plan\.present' cmd/wisp/plan_present_roster_dock_330r1_test.go` ＝ **0**。`git show edad424fe` 增行里带引号名字的行全落在 prose／票尾 bullet／probe 基线副本（既有字节的拷贝，非新字面量）；`cmd/wisp` 面新增 `.go` 仅一枚且零枚名字面量。两半出处具名：线上名＝golden `:12` 字节经 `wireNamesFromSSE330r1` 结构化抽出（本件 `:39-80`）；真名册＝`newReplyHost→runTextTask→assembleRuntime` 真装配的 `rt.bridge.Registry().List()`（本件 `:112-114`），名字最终出自 `internal/tools/plan_present.go:72`。⛔ 第二份硬编码名册：判语不成立。 |
| AC#3 | 维持（票面已 [x]；本票写面外） | `internal/agent` 面，330-v1 已独立裁成立（`:131` 冗余门、牙在 `rowOfTool165r3`）；本腿未复测、⛔ 背书。 |
| AC#4 | 维持（票面已 [x]；本票写面外） | 同上，330-r2 造世界＋330-v1 重放成对两发已裁；本腿未复测、⛔ 背书。 |
| AC#5 | **部分** | 实现腿未跑整包（它如实报了）。本腿按派单补一发：`unset GOTMPDIR TMPDIR TEMP TMP && export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"` 后 `go test ./cmd/wisp -count=1 -timeout 30m` ＝ **rc=1**（306.5s）。死点 6 枚具名（⛔ 当判据、只记死点）：`TestTicket223HandEditedFsLooseningCostsAnL2Card`／`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128`／`TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord`／`TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot`／`TestPanelHostRealWindowHopAndLifecycle`／`TestPanelHostLatencyPercentilesAC2`——全是真窗／真进程族，无一枚挂在本票交付件上（`Test330DockPinsTheWireNameToTheAssembledRoster` 在整包内自己绿）。终态日志无 panic、无 `ok` 总结行（非 `-v` 跑，RUN/PASS/SKIP 四数未取）。票面要求的正式台面（仓外导出树八件套＋起手红名册作差＝0）**今天仍不成立**⇒ 判「部分」：点名判据绿、四数作差未达。 |
| AC#6 | **PASS**（门禁本腿复跑） | 绝对路径 `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/` 列空 rc=0；`go vet ./cmd/wisp` rc=0；`go build ./cmd/wisp` rc=0；`sh scripts/d22scan.sh` rc=0 clean（本腿新读数存 `/d/tmp/330v-d22scan.log`，ban #8 cmd/=224 与实现腿读数同族）。 |
| AC#7 | **部分** | 本腿亲核 `git log`：近 80 枚＋`git log --all --grep=329` 零枚票 329 落地提交（命中全是别票信息里的字样）＝**329 未落地**，与实现腿声明一致；票面字面「本票必须排在票 329 之后」的前提今天不存在。撞面事实核清：本波 `cmd/wisp` 写面仅 `edad424fe` 一枚落地（上一枚同面提交＝`766bf71f4` 票 170 甲支），其后 3 枚提交零 `cmd/wisp` 字节 ⇒ 并发写同面**未实际发生**。框判定权归编排者（见翻框建议）。 |

## 翻框建议（⛔ 本腿不翻，归编排者）

1. **AC#1 可翻**：唯一判据（G1 发红）本腿亲重放成立，正向绿＋还原核验闭环。
2. **AC#2 可翻**：尺＝0 且出处两半具名，无第二份字面量。
3. **AC#6** 已 [x]，本腿复跑四门全 rc=0，维持。
4. **AC#0 可翻**：三把尺全真。
5. **AC#5 建议不按原字面翻**：四数作差＝0 的正式台面没跑过（实现腿与本腿都只跑了共享工作树）；可按「点名判据 Test330 绿＋整包死点 6 枚具名归真窗族」部分翻，或留给仓外导出树正式台面——编排者拍。
6. **AC#7 编排者拍**：两个事实摆明——①329 未落地（前提死了，派单字面无法满足）；②同面无并发撞（本波 cmd/wisp 仅一枚写提交）。若按「串行目的（不撞面）已达成」翻是解释性翻框，票面字面未满足；本腿不替拍。
7. **AC#3/#4** 已 [x]，330-v1 判语在档，维持不动。

## 本腿边界（如实）

- `internal/**` 仅两发当场突变＋cp 还原（sha 逐位回基线、git status 空），零枚提交字节。
- `cmd/wisp/**`、`frontend/**`、`design/**`、thresholds／golden 判据面零字节改。
- 四数（RUN/PASS/SKIP/PANIC）未取：整包一发按派单口径只记死点、⛔ 当判据，起手红名册作差未做。
