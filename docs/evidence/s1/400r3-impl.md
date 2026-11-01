# 400-r3 交付件：票 400 残余①（wire 携 DroppedCount——Go 侧四点挂接全落）

- 腿＝`400-r3`（票 400 残余①执行腿、写腿·实现者）；开工锚＝`556ad1b7b`（dev，恰为派单锚点；索引零已跟踪改动、写面 `internal/panel`＋`cmd/wisp` 无他人脏）。
- 施工图＝`docs/evidence/s1/400r2-halt.md`（400-r2 停手件）断点链路四处现读；**本腿开工时对该四处全部现读复核，逐字一致**（⇒ 停手条件①未触发）：①`StreamLog.DroppedCount()` accessor（pump.go:750-757，`droppedTotal` 总数）②`panel.TaskRosterState` 无 DroppedCount 字段（subagent_roster_197.go:94-101，六字段）③`TaskRosterSection` wire 六键无 `droppedCount`（:150-172）④render `taskRosterSectionFrom` 无拷贝行（:210-223）＋`panel_pump.go:250` 装配点=`state.DroppedStreamKeys = rt.stream.DroppedKeys()`。装配点内容与 400-v1 现读一致，对不上的只是「只动装配点即可」前提——本腿按 halt 件建议形状施工，未采信 400-v1 残余①措辞（编排者裁定③）。
- **wire 性质具名（停手条件②边界自查）**：本腿四点是**纯类型面补齐＋接线**——新增一个 state 字段、一个 wire json 键、一行 render 拷贝、一行装配赋值；未动 panel 契约语义（DroppedStreamKeys 的语义、TaskRosterSectionFrom 的排序/空数组口径、C17 桥、已知 Composer 键、`TestAPumpWithoutARosterReaderSendsFourKeys` 钉的顶层四键一字未动）。

## 1. 写面（逐枚，commit pathspec 同此 4 枚）

1. `internal/panel/subagent_roster_197.go` **三处**：
   - `TaskRosterState`＋`DroppedCount int`（`DroppedStreamKeys` 之后，注释具名「log 自身总数、ledger 只装得下最近一段」）；
   - `TaskRosterSection`＋`DroppedCount int `json:"droppedCount"``（注释沿用 DroppedStreamKeys 的口径：**Zero means nothing was dropped, never "unread"**；点名「named keys 永远只是总数最近一段」）；
   - `taskRosterSectionFrom`＋拷贝行 `DroppedCount: state.DroppedCount,`。
2. `cmd/wisp/panel_pump.go`：`:250` 装配行后＋一行 `state.DroppedCount = rt.stream.DroppedCount()`（带两行注释：不接这行，「some of them」会在页面上被读成「all of them」）。
3. `cmd/wisp/subagent_carrier_197_test.go`：①`carrier197Wire`（JSON 键直读、无 Go 侧帮衬）＋`DroppedCount int `json:"droppedCount"``；②既有 `TestRunPacketReportsTheStreamLogPastItsBound` 内照 ：561 形状加同族断言（`sect.DroppedCount` vs `driven.stream.DroppedCount()`）；③**新案 `TestRunPacketCarriesTheStreamDroppedCountPastTheHardCeiling`**（见 §3）。
4. 本证据件 `docs/evidence/s1/400r3-impl.md`。

## 2. ★TS diff 三行（⛔ 本腿自作、只写证据件——机主批后由编排者转界面支）

`frontend/src/lib/panel.ts` `TaskRosterSection`（现形状 ：221-228，本腿现读）恰好**三行新增**：

```diff
--- a/frontend/src/lib/panel.ts
+++ b/frontend/src/lib/panel.ts
@@ TaskRosterSection（:221 起）
   streamTruncated: boolean;
   streamElidedRunes: number;
   droppedStreamKeys: string[];
+  /** DroppedCount is every key the log stopped tracking at the hard ceiling,
+      including the ones the naming ledger no longer holds. Zero means none. */
+  droppedCount: number;
 }
```

- ★**Go 发 TS 未跟窗口＝knowingly 具名**：本腿落地即「Go wire 恒发 `droppedCount`、TS 名册未跟」——400-r1 AC#1 治的病以** knowingly 一小窗**复发，按编排者裁定②执行（400-v1 附条件③：此类跳点先取机主预批）。窗口内不红：对账尺（TestApprovalCardViewJSONKeysMatchFrontendTypes／TestComposerContractTypesMatchFrontend）只对**顶层键**、`TaskRosterSection` 嵌套类型无仪器覆盖（400-r2 停手件 §未走的路 现读）；窗口关闭＝机主批上面三行（台账已挂「droppedCount TS 键预批（可与追认并回）」）。

## 3. 读数①：wire 后读数（flood 场景非零）

- **仪器**：新案 `TestRunPacketCarriesTheStreamDroppedCountPastTheHardCeiling`——真装配（真 `taskRosterState`／真 render／真 `publishPanelSnapshot`／真 packet 字节），flood 用 **run 自己的 StreamLog**（`Append`＝run.go:1704 同一方法；硬上限=DefaultStreamKeys×2=64 键，64 枚**真孩子环路**在本套件不可达——本文件文末「未测清单」第 3 条原话，故 flood 内容合成、装配真——该条已如实改写记新状）。
- **读数（`logs/06-cmdwisp-consumer.log`）**：`--- PASS: TestRunPacketCarriesTheStreamDroppedCountPastTheHardCeiling (2.15s)`；wire 字节（packetTasks197 现印）尾部逐字：`..."droppedStreamKeys":["flood-400r3-000","flood-400r3-001",...,"flood-400r3-007"...],"droppedCount":9}`＋`hard ceiling crossed: dropped total=9 named=9`——**DroppedCount 到 Go wire 面、flood 下非零（9），与 log 自身 `DroppedCount()`=9 逐字相等**；named(9)≤total(9)（本 flood 9 次逐出全在 ledger 64 容量内，两侧同账）。
- 对照形状：既有 33 流案（只过软界 32）同一断言读 0==0（log 未逐出、双侧一致）——非零读数由硬上限案承载。

## 4. 读数②：反向判据（摘 wire ⇒ 面板侧零读数红；仓外 overlay、三读数＋还原核）

- **突变**＝仓外副本摘 `panel_pump.go` 装配（本腿加的 2 注释＋1 赋值行整块，diff=3 删现核），`-overlay` 构建（`mut/overlay-droppedcount-removed.json`），二进制不落工作树。
- **读数一（红）**：`mut-removed.test.exe` 同案 rc=1，红句逐字两句——`subagent_carrier_197_test.go:635: packet counts 0 dropped, want the log's own 9: the reader has to be the log, and a count the ledger lost room for still has to ride along`＋`:639: packet names 9 dropped keys but counts 0: the named half can never exceed the total`（named 半仍在＝单变量：摘的只是计数装配线；**「面板侧零 DroppedCount 读数红」逐字成立**）。
- **读数二（还原绿）**：同树无 overlay 干净二进制同案 rc=0 PASS（`dropped total=9 named=9`）。
- **读数三（盘上零动）**：`cmd/wisp/panel_pump.go` SHA256 before=after=`E53A5F0BCD70282A09AD7FED2F74131ADD1B0B2D7B5CD9E0EEA5E96D2126883A` 逐位同；`git status --porcelain --untracked-files=no` 全程仅本腿三枚产码改动。
- **正控对照**：突变件 rc=0 编译（State 字段仍在、render 照拷零值）⇒ 红咬的是装配线本身，不是编译面。

## 5. 读数③：零回归（panel 全套＋cmd/wisp 消费面＋六族＋262/640/902）

- **panel 全套**：`go test ./internal/panel -v -count=1` rc=0 ⇒ **RUN=271 PASS_TOP=173 FAIL_TOP=0 SKIP=0——与 400-r1 刚落的基线逐位同、零破**（`logs/05-panel-full.log`）。
- **cmd/wisp 消费面**（roster/streamlog/instructions/197/329 全族，-run＝`TestRunPacket|Test197|TestSubagentStreamKey|TestBookedSnapshotReachesHostBitForBit329|TestRunBooksWithASnapshotOfItsLiveQueue|TestSnapshotWorkspaceSectionReportsTheNarrowing|TestThePumpIsDrivenNotJustAssembled`）：rc=0 ⇒ **RUN=16 PASS=16 FAIL=0 SKIP=0**＝派单 15 族全绿＋本腿新案 1 枚（`logs/06-cmdwisp-consumer.log`，80.4s）。
- **六族 staged 合跑**（no-tag 二进制＋三 DLL 贴 bin＋397-r4 同方六 env＋cwd=包目录，`-test.run 'Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384'`）：**PASS=34 FAIL=1 SKIP=0——与 397-r4/392-r1 基线 34/1/0 逐位归位**（`logs/18-sixfamilies.log`）；唯一红 `Test384ResidentWatchdogEndToEnd`(3.34s)＝合跑负载时窗既有形状，**单发复跑 rc=0（1.17s PASS）归因核**（`logs/19-384-single.log`）⇒ 红名集合零新增。
- ★**`262/640/902` 逐位同**：`resident_voice_15e_windows_test.go:428: LATENCY openDelay=262ms chunk=640ms ballLatency=902ms (openDelay+chunk) seamWaitToFirstText=5.51s finalLatency=4.096s speechMS=3392 closedBy=end`。

## 6. 读数④：四数／两尺／d22scan 导出树＋正控

- **四数**：panel 271/173/0/0；消费面 16/16/0/0；d22scan 自检 RUN=77 PASS=35 FAIL=0 SKIP=0（与 397-r4 基线 35/0/0 逐位同）；六族 34/1/0。
- **两尺**（对象=`internal/panel`＋`cmd/wisp`）：甲 `gofmt -l`＝**空**、乙 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`）＝**空**，双 rc=0（`logs/01-gofmt-l.txt`／`02-gofumpt-l.txt`）。
- **d22scan 实扫**（landing tree，Git Bash `scripts/d22scan.sh`）：**rc=0 clean**；分母 bans#1-5 internal/=288、cmd/=52、ban#6 frontend/=116、ban#7 internal/tools/=25、ban#8 design/=39/frontend/=116/internal/=734/cmd/=184（随树漂移如实记；`logs/13-d22scan-landing.log`）。
- **d22scan 导出树**（`git archive HEAD` 仓外树＋本腿三件 SHA 等值 overlay 进树，tree/tools/d22scan `go run . -root tree`）：**rc=0 clean**（design/=14＝gitignored 面导出差，397-v2 §6⑤ 同形伪影；其余分母逐位同；`logs/16-d22scan-exported.log`）。
- **正控**＝d22scan 自检（`TestScanDetectsAllSeededViolations` 在场 PASS；种子违规 35/0/0）——尺活着、双向被钉。
- `go vet ./internal/panel ./cmd/wisp` rc=0；`go build ./internal/panel ./cmd/wisp` rc=0（`logs/03`/`04`）。

## 7. 纪律格与写面申报

- 仓内改动恰 §1 的 4 枚（3 产/测码＋本证据件）；⛔ 未动：panel 契约语义/C17 桥（bridge.go 零字节）、`frontend/src/lib/panel.ts`（TS diff 只在本件 §2）、thresholds/golden、testdata 仓根、design、票面、台账、specs、PLAN、D43；⛔ 下载、⛔ deps.toml；⛔ push（只 commit、推送归编排者）；`add`/`commit` 显式 pathspec 逐枚手写。
- 仓外临时件 `D:\tmp\wisp\20261101\400r3\**`（mut/bin/logs/tree/cases，只建不删）：突变副本＋overlay json＋双对照二进制＋no-tag 六族二进制＋全部 log。
- 环境：GOTMPDIR 全程 unset；cmd/wisp 测试二进制 DLL 以仓根前置 PATH／bin 同目录贴靠两形跑通；pwsh 直跑测试二进制时 `-test.*` 旗标须 `=` 形（首次按空格分参报 `flag provided but not defined: -test`，改形即过——非码面，如实记）。
- 未定义即停三条款：①断点链路与现读逐字一致（未触发）②wire 未动契约语义、纯类型面＋接线（未触发）③DroppedCount 已到 Go wire 面（读数①非零 9）（未触发）。
