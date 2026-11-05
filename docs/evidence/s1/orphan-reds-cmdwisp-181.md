# 181 — cmd/wisp 孤儿红修复 II（98 C1 余册：夹具缝，27 枚 POSIX 转绿 + 10 枚 STOP）

**Status**：交付件（写腿自证）。**写面**：`cmd/wisp/**`（改 4 枚、新增 1 枚）＋ 本文件；**未 push**。
**锚定树**：`7eb0e3f4b`（派单给的 HEAD，= 180 件 `ae8136261` 并进去之后的那一枚）。
**本腿提交**：`331a4b1cd`（工作树带本腿改动时取的全部读数；读数期间 HEAD 前进过一枚**docs-only** 提交，
见 §9）。
**仪器**：Windows 侧 `go1.27.1 windows/amd64`（`PATH` 前置 `build;third_party/sherpa-onnx`）；
POSIX 侧 `docker.io/library/golang:1.27-bookworm`（容器内 `linux/amd64`、`CGO_ENABLED=1`），
快照一律 `git archive 7eb0e3f4b` + overlay（不建 worktree、不 checkout）。
**时间**：2026-11-05 13:5x–14:2x +08。

## 0. 一句话结论

`docs/evidence/s1/orphan-reds-cmdwisp-180.md` §5 留下的**37 枚 C1 红**，本腿按票分两堆：

- **27 枚修好**（POSIX 上真跑、真绿）：12×6、197×5、200×2、151×2、35×3、365×2、370×2、174r8×2、378×1、220×1、329×1。
  改法**逐字沿用 180 的缝**：夹具在 Windows 上仍走 DPAPI（同 ref、同 blob、同 `Store` 调用、同字节），
  其余平台换成 `env:` ref + `t.Setenv`。
- **10 枚 STOP**（63×8、248×2）：这 10 枚的红**不在夹具里，在生产路径里**——
  `wisp secret set` 与 `StoreCredential` 自己调 `secret.Store(dpapi:…)`，POSIX 上的 fail-closed 桩
  （C28 / `internal/secret/protect_other.go`）就是它的正解。改夹具一格也动不了它，见 §5（**未改一行**）。

整包读数（POSIX，同一仪器、同一快照、同一命令）：

| | PASS | FAIL | SKIP | `=== RUN` |
|---|---|---|---|---|
| 修前（`7eb0e3f4b`） | 188 | **40** | 0 | 310 |
| 修后（同树 + 本腿 overlay） | **215** | **13** | 0 | 310 |

**27 = 40 − 13**；`=== RUN` 前后都是 **310**，用例名清单逐枚相同（§4 第 3 条给可核尺）。
剩下 13 枚 = C4（票 131 的**故意红**，1）＋ C5（我这台 `git archive` 无 `.git` 的**假红**，1）＋
C6（票 119 孤儿，`406` 已立，1）＋ **本件 STOP 的 10 枚**（63×8、248×2）。
**0 断言改动、0 契约改动、0 `//go:build windows`、0 `t.Skip`、0 用例改名**。

## 1. 180 的余册，本腿点了哪几堆

180 §5 逐字留的名册与 98 §6 的逐枚表对得上；本腿把 37 枚按**夹具归属**重排后，它们的形状只有四种
（前三种本腿修，第四种 STOP）：

| 夹具 | 谁的红 | 枚数 |
|---|---|---|
| `run_test.go` 的 `newRunFixture` | 12×6、197×5、200×2、151×2、35×3、220×1、329×1 | **20** |
| `tasklogstate_365_test.go` 的 `newProbe365`（370 直接复用） | 365×2、370×2 | **4** |
| `task_output_pointer_notice_174r8_test.go` 的 `new174r8EmptyRootsFixture` / `tasklog_fkcheck_378_test.go` 的 `fk378WriteConfig` | 174r8×2、378×1 | **3** |
| 生产路径自带（不是夹具） | 63×8、248×2 | **10**（STOP） |

20 + 4 + 3 = **27**。20 枚那堆是本次最高价值的一刀：**一处夹具、七张票的用例同时离开红**。

## 2. 改了什么（四枚夹具 + 一枚新缝）

新增 `cmd/wisp/secret_fixture_181_test.go`（唯一新文件，`_test.go`，不进产品二进制），只补 180 没覆盖的两种形状：

```
fixtureKeyRef181(t, dataDir, dpapiRef, key) string
  windows  → 与今天逐字节相同：secret.NewStore(dataDir) + st.Store(dpapiRef, key)，返回 dpapiRef 本体
  其余平台 → t.Setenv(fixtureKeyEnv180, key)，返回 "env:WISP_TEST_FIXTURE_KEY"
               （复用 180 的常量，全包只此一个变量名）

fixtureKeyDrop181(t, dataDir, ref)
  windows  → st.Delete(ref)          // 那几枚用例原本内联写的就是这一对
  其余平台 → os.Unsetenv(fixtureKeyEnv180) + t.Cleanup 复原
               （用 unset 而不是空值：Resolve 对"没设"和"设成空"给的是两句不同的话）
```

与 180 的 `fixtureKeyRef180(t, dir, blobID, key)` 的分工：180 那枚收 blob id、自己拼 `"dpapi:"+id`；
本腿这四是**夹具已经把整条 ref 写死在 config 串里**（再写第二遍进 `Store`），所以收整条 ref 原样穿透——
Windows 上 config 里的字节与 blob 名都仍是原来那个，返回值就是入参本身。

四处夹具改成"先用缝拿到 ref，再把 ref 写进 config"，`Store` 块从夹具里搬进缝（Windows 分支）：

| 文件 | 夹具 | 被它扶的票 |
|---|---|---|
| `cmd/wisp/run_test.go` | `newRunFixture`（顺带给 `runFixture` 加 `keyRef` 字段） | 12（6）、197（5）、200（2）、151（2）、35（3）、220（1）、329（1） |
| `cmd/wisp/tasklogstate_365_test.go` | `newProbe365` | 365（2）、370（2） |
| `cmd/wisp/task_output_pointer_notice_174r8_test.go` | `new174r8EmptyRootsFixture` | 174r8（2） |
| `cmd/wisp/tasklog_fkcheck_378_test.go` | `fk378WriteConfig` | 378（1） |

**唯一一处不是"ref 从哪来"的改动**（必须同批改，否则缝会造出一枚新红）：
`run_test.go` 的 `TestMissingBlobFailsUnconfiguredNeverSilently` 原先把"这枚 ref 现在指不到东西"写成
`st.Delete("dpapi:…")`。缝一改，POSIX 上 ref 是 `env:`，删 blob 什么也没删掉（变量还在），
用例会反过来红在"缺凭据却 exit 0"。改法：`fixtureKeyDrop181(t, f.dir, f.keyRef)`——
Windows 上仍是 `Delete` 同一个 blob，POSIX 上把变量 unset。
**四条断言（exit 2 / 名字必须是 Unconfigured / 不得有请求到达 provider / 输出不得带 key）一字未动。**

合计：**改 4 枚、新增 1 枚**。`git diff --numstat` 里删掉的每一行都在上面四枚的夹具块与这条前提行内。

## 3. RED before / GREEN after（原文摘录）

**RED before**（POSIX，`7eb0e3f4b`，无 overlay，整包；37 枚 C1 逐枚同一句根因，日志点＝用例的调用行）：

```
--- FAIL: TestRunTextTaskTextPathEndToEnd (0.86s)
    run_test.go:169: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestRunPacketCarriesTheSubagentItsRosterRowFed (0.70s)
    subagent_carrier_197_test.go:233: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestThePumpIsDrivenNotJustAssembled (0.61s)
    panel_pump_test.go:414: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestBookedSnapshotReachesHostBitForBit329 (0.80s)
    panel_snapshot_deliver_329_test.go:77: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: Test365RunRowsAStuckTaskAsStuck (0.72s)
    tasklogstate_365_test.go:186: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: Test370TruncatedRunRowsAnInterruptedTask (0.98s)
    task370_ac2_test.go:99: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: TestTaskOutputPointerNotice174r8StaysQuiet (1.32s)
    task_output_pointer_notice_174r8_test.go:165: secret: DPAPI is only available on Windows; use an env: ref on this platform
--- FAIL: Test378SpawnedChildToolCallsHaveAParentRow (0.00s)
    tasklog_fkcheck_378_test.go:395: secret: DPAPI is only available on Windows; use an env: ref on this platform
```

**GREEN after**（POSIX，同一快照 + 本腿 overlay，先点名 27 枚那一发）：

```
--- PASS: TestRunTextTaskTextPathEndToEnd (…)
--- PASS: TestRunTextTaskFailNextIsClassified
--- PASS: TestHostDispatchThroughTheAssembledBridge
--- PASS: TestComposedGateBlocksAWriteForTwoSeconds
--- PASS: TestRunTextTaskKeyResolvesInTheStore
--- PASS: TestMissingBlobFailsUnconfiguredNeverSilently
--- PASS: TestRunPacketMarksTheRosterRowACardIsHolding
--- PASS: TestRunPacketCarriesTheSubagentItsRosterRowFed
--- PASS: TestRunPacketReportsTheStreamLogPastItsBound
--- PASS: TestRunPacketCarriesTheStreamDroppedCountPastTheHardCeiling
--- PASS: TestRunPacketLightsTheRosterRowOfItsLiveL1Window
--- PASS: TestRunPacketCarriesTheLoadedInstructionFiles
--- PASS: TestRunPacketSaysOffAndSaysItDifferentlyThanNothingFound
--- PASS: Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands
--- PASS: TestCompositionRootClosesTheLoopTasksTaintScope
--- PASS: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit
--- PASS: TestRunBooksWithASnapshotOfItsLiveQueue
--- PASS: TestSnapshotWorkspaceSectionReportsTheNarrowing
--- PASS: TestThePumpIsDrivenNotJustAssembled
--- PASS: TestBookedSnapshotReachesHostBitForBit329
--- PASS: Test365RunRowsAStuckTaskAsStuck
--- PASS: Test365RunKeepsItsWordForTheControls
--- PASS: Test370TruncatedRunRowsAnInterruptedTask
--- PASS: Test370TheInterruptedRowNeedsTheCeiling
--- PASS: TestTaskOutputPointerNotice174r8SpeaksTheRealReason
--- PASS: TestTaskOutputPointerNotice174r8StaysQuiet
--- PASS: Test378SpawnedChildToolCallsHaveAParentRow
ok  github.com/CarlosShao/wisp/cmd/wisp  109.902s
PASS=29 FAIL=0 SKIP=0
```

（`PASS=29` 里多出的两枚是同文件被同一 `-run` 前缀带上的 `Test365TaskLogStateSpeaksOneWordOnBothPaths`
与 `Test365TruncatedCannotBeReachedThroughThisCommand`，它们修前也绿。）

整包：`PASS=215 FAIL=13 SKIP=0`（修前 `PASS=188 FAIL=40 SKIP=0`）。**SKIP 两发都是 0**。

## 4. 这为什么不叫"放宽/摘尺"（五条可核的尺）

1. **没有删用例、没有加 build tag、没有改名**：`git diff --name-status` 只有
   `M run_test.go / M task_output_pointer_notice_174r8_test.go / M tasklog_fkcheck_378_test.go /
   M tasklogstate_365_test.go / A secret_fixture_181_test.go`；`_windows_test.go` 一枚没动。
2. **没有改任何断言的判据**：四枚夹具只改"ref 从哪来"，`fixtureKeyDrop181` 只改"怎么让 ref 指不到东西"。
   可核尺：`git diff -U0 -- cmd/wisp | grep -E '^[-+][^-+]' | grep -Ec 't\.(Errorf|Fatalf|Error|Fatal|Skip|Logf)'`
   = **10**，且 10 行全是 `-  t.Fatal(err)`（被删的 `NewStore`/`Store` 块里的错误出口），**加号侧 0 行**。
3. **用例一枚没增没减没改名**：两发整包的 `=== RUN` 都是 **310**；把两发的 `^--- (PASS|FAIL|SKIP)`
   名字排序后 diff，差异**只有 27 行 FAIL→PASS**（同一批名字），无新增、无消失、无 SKIP。
4. **键值没换**：四处夹具原本都写同一枚 `fakeStoreKey`（`run_test.go:53`），缝里传的就是它。
5. **Windows 侧字节零变化**：缝的 Windows 分支与旧代码同序同参（`NewStore` → `Store(ref, key)`），
   config 里写回的也是同一枚 ref；模板里 `api_key_ref = %q` 对同一枚纯 ASCII ref 渲染出的引号串
   与原字面量逐字节相同。所以 Windows 的 DPAPI 覆盖**没有减少**，变绿只发生在"DPAPI 根本不存在"的平台上
   （那里今天只有 fail-closed 桩，无覆盖可减）。
   附带说明：夹具在 POSIX 上不再**先**建 `secrets/` 目录，但组合根自己会建（`cmd/wisp/run.go:477`
   的 `secret.NewStore(s.dataDir)`），所以被跑起来的那一发磁盘形状不变。

## 5. STOP：63×8 与 248×2（不属夹具缝，未改一行）

这 10 枚的红**不在夹具里**，在**生产路径**里；两处都死在同一个 fail-closed 桩上（POSIX 逐字）：

| 票 | 用例 | 死在哪一行 | 逐字 |
|---|---|---|---|
| 63 | `TestSecretSetGetListRoundTrip` | `secret_test.go:180` | `set --from-stdin exit = 1 … wisp secret set: storing dpapi:…: secret: DPAPI is only available on Windows; use an env: ref on this platform` |
| 63 | `TestSecretFromStdinWritesNoIntermediateFile` | `secret_test.go:569` | 同上（`storing dpapi:no-tmp`） |
| 63 | `TestSecretFailurePathsLogAndPrintNoPlaintext`（3 子用例） | `secret_test.go:608/627/689` | 同上（`corrupt` / `portable` / `late`） |
| 63 | `TestSecretUnsetRefusesWhileReferenced` | `secret_test.go:783` | 同上（`deepseek`） |
| 63 | `TestSecretSameNameUnderThreeEnvsIsThreeBlobs` | `secret_test.go:963` | 同上（`shared-key`） |
| 63 | `TestSecretPortableModeUsesTicket06Seam` | `secret_test.go:1053` | 同上（`portable-own`） |
| 63 | `TestSecretEndToEndConfigRefResolvesAtRequestTime` | `secret_test.go:1086` | 同上（`mockllm`） |
| 63 | `TestSecretOverwriteIsAnnounced` | `secret_test.go:1196` | 同上 |
| 248 | `TestAC2CredentialSentinelAppearsInNoArtifact` | `panel_config_248_test.go:164` | `the credential write was refused: 密钥没有写进凭据存储（配置未改动，未录入任何密钥）：secret: DPAPI is only available on Windows; use an env: ref on this platform` |
| 248 | `TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef` | `panel_config_248_test.go:272` | 同上（`credential write:` 前缀） |

**为什么夹具改不动它们**：

- 63 的红是 `secretCmd.openStore()`（`cmd/wisp/secret.go:296`）→ `secret.Store(...)` → `protect`。
  这两张票的**被测对象就是 DPAPI 写入本身**：blob 要在盘上、密文要不可读、三个 env 要三份**不同**密文、
  损坏 blob 要出 P13 句、portable 要能读自己写的 blob。把 config 换成 `env:` ref 根本不在这些用例的路径上
  （它们不读 config 的 ref，是 `set` 自己铸 `dpapi:` ref）。
- 248 的红是 `configStore.StoreCredential`（`cmd/wisp/panel_config_store.go:327`）→ `s.secrets.Store(ref, value)`，
  `ref` 由 `secret.NewRef()` 现铸；它的 config 早已是 `env:DEEPSEEK_KEY`，红的是**凭据写入**这一刀，
  而且两条断言（回话里必须出现 `dpapi:` / config.toml 必须落 `api_key_ref = 'dpapi:`）明文要求 DPAPI 真写了。

**只剩两条路，两条都超出本腿写面**：

① 给两枚文件加 `//go:build windows` / 改 `_windows_test.go`（`internal/secret/store_test.go` 的票 70 先例）——
   派单明文禁止，98 AC#2 也把"把用例搬出 CLI 分母"划成禁区。
② 在 `internal/secret` 里开一个 protector/Store 缝（或把 `configStore.secrets`/`secretCmd` 的 store 改成接口
   并在测试里装替身）——这是 **C28 / SPEC-02 §6 的平台设计面**（DPAPI-only、fail-closed、无明文回退），
   加缝就是改冻结契约的形状，属派单说的"需要 internal/contract 改动"⇒ STOP。

**本腿对这两张票的处置：一行没改，红原样保留，交编排者裁**（建议开两张票：63 的"DPAPI 主体用例的平台口径"、
248 的"凭据写入腿在非 Windows 上的可测面"，归属票 06/08 的契约主人，不是夹具腿）。

## 6. 门禁读数（逐条命令 + 读数）

全部在本机、工作树带本腿改动时跑：

```
$ go build ./cmd/... ./internal/... ./tools/...
rc=0

$ go test ./cmd/wisp -run 'TestRunTextTask|TestHostDispatchThroughTheAssembledBridge|
    TestComposedGateBlocksAWriteForTwoSeconds|TestMissingBlobFailsUnconfiguredNeverSilently|
    TestRunPacket|Test197Subagent|TestCompositionRootClosesTheLoopTasksTaintScope|
    TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit|Test365|Test370|
    TestTaskOutputPointerNotice174r8|Test378SpawnedChildToolCallsHaveAParentRow|
    TestRunBooksWithASnapshotOfItsLiveQueue|TestSnapshotWorkspaceSectionReportsTheNarrowing|
    TestThePumpIsDrivenNotJustAssembled|TestBookedSnapshotReachesHostBitForBit329|
    TestTicket255|DispatchGate' -count=1 -timeout 420s
ok  github.com/CarlosShao/wisp/cmd/wisp  123.372s      # 上限 420s，余量 297s

$ /d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/
rc=0   （空）

$ go vet ./cmd/wisp/
rc=0

$ sh scripts/d22scan.sh
rc=0   # d22scan: clean - no D22 ban violations
```

**与 180 的对照**：180 那条 gate 撞到过 240s 上限只剩 19s 余量；本腿按派单用 `-timeout 420s`，
实测 123.4s（本腿这 27 枚里最重的是 365/370 的四发真实 `wisp run`，远轻于 101 的 L1/L2 窗口）。

## 7. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"

# A. Windows 侧（DPAPI 路径 —— 本腿改动的"零变化"尺）
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go test ./cmd/wisp -run 'TestRunTextTask|TestHostDispatchThroughTheAssembledBridge|TestComposedGateBlocksAWriteForTwoSeconds|TestMissingBlobFailsUnconfiguredNeverSilently|TestRunPacket|Test197Subagent|TestCompositionRootClosesTheLoopTasksTaintScope|TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit|Test365|Test370|TestTaskOutputPointerNotice174r8|Test378SpawnedChildToolCallsHaveAParentRow|TestRunBooksWithASnapshotOfItsLiveQueue|TestSnapshotWorkspaceSectionReportsTheNarrowing|TestThePumpIsDrivenNotJustAssembled|TestBookedSnapshotReachesHostBitForBit329|TestTicket255|DispatchGate' \
  -count=1 -timeout 420s

# B. POSIX 侧（本腿修的对象；docker 与 180/98 同一只，快照 git archive + overlay）
bash .scratch/leg181/linux-red.sh '<上面那串 -run>'
bash .scratch/leg181/linux-full-before.sh   # 修前：期望 PASS=188 FAIL=40 SKIP=0
bash .scratch/leg181/linux-full-after.sh    # 修后：期望 PASS=215 FAIL=13 SKIP=0

# C. 门禁
go build ./cmd/... ./internal/... ./tools/...
/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/
sh scripts/d22scan.sh
```

（`.scratch/leg181/` 下：`linux-red.sh`、`linux-full-before.sh`、`linux-full-after.sh`、
`overlay/cmd/wisp/{4 枚改 + 1 枚新}`；原始日志 `/d/tmp/leg181/{full-before.log,full-after.log,run-after.log}`。
临时件只建不删。）

## 8. 本件没有建立的事（宁缺毋造）

- **没有修 63×8 与 248×2**（§5）：它们需要 `internal/secret`/`cmd/wisp` 生产面的缝，超出"夹具"写面。
  本件不主张它们"不可修"，只主张"夹具缝修不了"。
- **没有跑 Windows 整包**，只跑了 §6 那一组（123.4s ok）。所以"Windows 整包全绿"本件**不主张**；
  主张的是"本腿改的这 27 枚在 Windows 上仍绿"。
- **没有取 CI 读数**：全部读数是本机（Windows）与容器（POSIX），**不得**引成 CI 结论。
- **没有在 POSIX 上验 `-count=2`**：本件是普查/修复口径，不是门禁口径（180 §8、98 §7 同样留白）。
- **C6 / 票 119 仍红**（`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`）：`406` 已立
  （Blocked by 三选一），需要 `internal/risk` 或替票 119 裁定 POSIX 哨兵，不是夹具缝。本腿原样保留。
- **C4 / 票 131 的故意红、C5 的假红（无 `.git`）** 仍在这 13 里，理由与 180 §5 逐字相同，本腿不重复论证。
- **`internal/secret` 的两个形状没有改**：`protect_other.go` 的 fail-closed 桩、`refs.go` 的两种 ref
  都是冻结契约（C28 / C27），本件只消费它们。

## 9. 读数期间 HEAD 动过一枚（docs-only，两树 `cmd/wisp` 逐字节相同）

派单锚定 `7eb0e3f4b`；**修前**整包那一发（13:5x 起）的 `git archive` 取的就是它。
**修后**整包那一发（14:1x 起）取快照时，别腿的 `b93467e43`（14:12:49，`342-r3` 证据件）已经落地，
所以那一发取到的是 `b93467e43`。两棵树对 `cmd/wisp` 的差：

```
$ git diff --stat 7eb0e3f4b b93467e43 -- cmd/wisp
（无输出，0 行）
$ git diff --name-only 7eb0e3f4b b93467e43
docs/evidence/s1/342-r3-impl.md
```

⇒ 两张整包读数可以并排引：被检对象（`cmd/wisp`）在两枚提交之间逐字节相同，变的只有别人的一份证据件。
本腿的 Windows gate 跑在**工作树**上（= `7eb0e3f4b` 的 `cmd/wisp` + 本腿 5 枚文件），与上面两棵树都不冲突。

