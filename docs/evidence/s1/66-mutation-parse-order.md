# 66 — AC#1 变异检验：解析器顺序缺陷的回归用例咬合证明

**时间**：2026-09-20T14:24Z–14:31Z（Git Bash，本机 DESKTOP-LVS7839）
**被测对象**：`internal/proc`（`parseSystemProcesses` + 新抽出的唯一解码器 `WalkSystemProcesses`）
**基线 commit**：`86e868d`（缺陷①修复 + 回归用例已入库）

## 做法

变异 = 把 `internal/proc/systemprocs_windows.go::WalkSystemProcesses` 退回修复前的顺序：
读到 `NextEntryOffset` 之后、解码当前项**之前**先 `if next == 0 { return nil }`。
这正是 A14 记录的原始缺陷形态（`internal/proc/treemetrics_windows.go:240-243` 旧代码）。

变异落地证据（`grep -n "MUTATION" -A3 -B3 internal/proc/systemprocs_windows.go`）：

```text
73-			return nil // buffer exhausted: chain ended without a terminator
74-		}
75-		next := *(*uint32)(unsafe.Add(base, offset))
76:		// MUTATION (ticket 66 AC#1 mutation proof): restore the pre-fix order.
77-		if next == 0 {
78-			return nil
79-		}
```

变异撤销后的复核（`grep -c "MUTATION" internal/proc/systemprocs_windows.go` → `0`）。

## 变异态（期望：新用例转红，包内其余用例仍绿）

`go test -count=1 -v ./internal/proc/` → `rc=1`

```text
--- PASS: TestBootTestEnvBootsAllPackages (0.00s)
--- PASS: TestBootSingleInstanceConflict (0.00s)
--- PASS: TestMutexNamesPerEnv (0.01s)
--- PASS: TestLayoutForkMatrix (0.00s)
--- PASS: TestLayoutForProdDev (0.00s)
--- PASS: TestLayoutForTestEnv (0.00s)
--- PASS: TestLayoutForUnknownEnv (0.00s)
--- PASS: TestDefaultLayoutUnknownEnvFails (0.00s)
--- PASS: TestPortableOverride (0.01s)
--- PASS: TestSummaryBadge (0.00s)
--- PASS: TestAcquireRejectsEmptyNames (0.00s)
--- PASS: TestJobScopeKillsChildOnClose (0.03s)
--- PASS: TestJobScopeTreeAccounting (0.05s)
--- PASS: TestStartInJobRejectsBadCommand (0.01s)
--- PASS: TestShutdownOrderAudit (0.00s)
--- PASS: TestShutdownFastPathSkipsOnlyStep7 (0.00s)
--- PASS: TestShutdownBoundedWaitAbandonsAndContinues (0.14s)
--- PASS: TestShutdownNilHooksAreRecordedSkipped (0.00s)
--- PASS: TestSingleInstanceSecondExitsAfterSignallingFirst (0.03s)
--- PASS: TestSystemProcessSnapshotSeesSelf (0.03s)
--- PASS: TestParseSystemProcessesRejectsTruncatedEntry (0.00s)
--- PASS: TestParseSystemProcessesRejectsNonAdvancingEntry (0.00s)
--- PASS: TestWalkSystemProcessesEmptyBuffer (0.00s)
--- FAIL: TestParseSystemProcessesKeepsLastSnapshotEntry (0.00s)
--- FAIL: TestWalkSystemProcessesStopsAtTerminator (0.00s)
--- FAIL: TestParseSystemProcessesSkipsPIDZero (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/proc	0.440s
FAIL
```

失败原文（三条都是本票新增的 A14 用例，逐字）：

```text
=== RUN   TestParseSystemProcessesKeepsLastSnapshotEntry
    systemprocs_windows_test.go:95: parsed 4 of 5 snapshot entries: [1337 60060 4 4004]
--- FAIL: TestParseSystemProcessesKeepsLastSnapshotEntry (0.00s)
=== RUN   TestWalkSystemProcessesStopsAtTerminator
    systemprocs_windows_test.go:150: walker delivered [100], want [100 200] (stop after the terminator entry)
--- FAIL: TestWalkSystemProcessesStopsAtTerminator (0.00s)
=== RUN   TestParseSystemProcessesSkipsPIDZero
    systemprocs_windows_test.go:164: parseSystemProcesses: system process snapshot parsed zero entries
--- FAIL: TestParseSystemProcessesSkipsPIDZero (0.00s)
```

`parsed 4 of 5 snapshot entries: [1337 60060 4 4004]` —— 缺的正是手工放在链末项的 pid `8080`
（`0x1F90`），且断言的是它的 handle 数 / 私有工作集 / 线程数 / user+kernel 时间是否入表，
不只是 `len(out)`（len 只是第一道断言，末项在场后仍逐项比对字段）。

## 复原态（期望：全绿）

`go test -count=1 -v ./internal/proc/` → `rc=0`

```text
--- PASS: TestParseSystemProcessesKeepsLastSnapshotEntry (0.00s)
--- PASS: TestWalkSystemProcessesStopsAtTerminator (0.00s)
--- PASS: TestParseSystemProcessesSkipsPIDZero (0.00s)
--- PASS: TestParseSystemProcessesRejectsTruncatedEntry (0.00s)
--- PASS: TestParseSystemProcessesRejectsNonAdvancingEntry (0.00s)
--- PASS: TestSystemProcessSnapshotSeesSelf (0.02s)
--- PASS: TestWalkSystemProcessesEmptyBuffer (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/proc	0.412s
```

（复原态完整输出 26 条用例全 PASS，与变异态同一集合。）

## 结论

- 三条新用例中，直接钉住 A14 的 3 条在变异态全部转红，包内其余 23 条不受影响 → 用例咬住了缺陷本身。
- 其中 `TestSystemProcessSnapshotSeesSelf`（真机快照必须看得见自己）在变异态仍 PASS 是**符合预期**的：
  它依赖内核当时的链序，属于按进程创建顺序随机的旧症状，因此本票不把它当唯一防线，
  手工构造缓冲的用例才是。

---

# 复测（独立第二人）— agent-ticket66b，2026-09-20T15:02Z

**基线**：`81afe3c`（工作树含 `86e868d` + `00bbb76` + `63b61ef`，变异前后 `git diff` 均为空）。
动机：编排者要求独立复核 AC#1 的变异检验真的咬得住，而不是引用上一位代理的自述。
本次变异比上次**更严**：链序恰好落在末项的进程在本次真机跑里正是采样者自己，
于是 `TestSystemProcessSnapshotSeesSelf` 与两条树外采样用例（`00bbb76` 新增的 AC#3 用例）
也一起转红 —— 上次那轮这两类是 PASS。

## 变异落地证据（先证明变异真的进了文件，再谈红）

`grep -n "MUTATION-66B" internal/proc/systemprocs_windows.go` → `76:`

```text
74-			}
75-		next := *(*uint32)(unsafe.Add(base, offset))
76:		// MUTATION-66B (ticket 66 AC#1 mutation proof, agent-ticket66b):
77-		// temporarily restore the pre-fix ordering — test the chain
78-		// terminator BEFORE decoding the current entry.
79-		if next == 0 {
80-			return nil
81-		}
82-		if offset+systemProcEntryFloor > uintptr(len(buf)) {
```

`git diff --stat internal/proc/systemprocs_windows.go` → `1 file changed, 6 insertions(+)`
（即：这次编辑确实改了磁盘上的文件，不是编辑器幻觉）。

## 变异态原始输出（`go test -count=1 -v ./internal/proc/` → rc=1）

逐条结果行（原文顺序，共 24 PASS / 6 FAIL / 1 SKIP）：

```text
--- PASS: TestBootTestEnvBootsAllPackages (0.00s)
--- PASS: TestBootSingleInstanceConflict (0.00s)
--- PASS: TestMutexNamesPerEnv (0.01s)
--- PASS: TestLayoutForProdDev (0.00s)
--- PASS: TestLayoutForTestEnv (0.00s)
--- PASS: TestLayoutForkMatrix (0.00s)
--- PASS: TestPortableOverride (0.01s)
--- PASS: TestSummaryBadge (0.00s)
--- PASS: TestLayoutForUnknownEnv (0.00s)
--- PASS: TestDefaultLayoutUnknownEnvFails (0.00s)
--- FAIL: TestExternalSamplerReadsSubjectFromOutside (0.04s)
--- PASS: TestExternalSamplerRefusesToMeasureItself (0.00s)
--- PASS: TestExternalSamplerFailClosedOnDeadSubject (0.01s)
--- FAIL: TestExternalSamplerTracksSubjectGrowth (0.04s)
--- SKIP: TestHelperProcess (0.00s)
--- PASS: TestJobScopeKillsChildOnClose (0.03s)
--- PASS: TestJobScopeTreeAccounting (0.05s)
--- PASS: TestStartInJobRejectsBadCommand (0.01s)
--- PASS: TestShutdownOrderAudit (0.00s)
--- PASS: TestShutdownFastPathSkipsOnlyStep7 (0.00s)
--- PASS: TestShutdownBoundedWaitAbandonsAndContinues (0.14s)
--- PASS: TestShutdownNilHooksAreRecordedSkipped (0.00s)
--- PASS: TestSingleInstanceSecondExitsAfterSignallingFirst (0.03s)
--- PASS: TestAcquireRejectsEmptyNames (0.00s)
--- FAIL: TestParseSystemProcessesKeepsLastSnapshotEntry (0.00s)
--- FAIL: TestWalkSystemProcessesStopsAtTerminator (0.00s)
--- FAIL: TestParseSystemProcessesSkipsPIDZero (0.00s)
--- PASS: TestParseSystemProcessesRejectsTruncatedEntry (0.00s)
--- PASS: TestParseSystemProcessesRejectsNonAdvancingEntry (0.00s)
--- FAIL: TestSystemProcessSnapshotSeesSelf (0.01s)
--- PASS: TestWalkSystemProcessesEmptyBuffer (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/proc	0.419s
FAIL
```

六条红名的失败原文（逐字）：

```text
    externalsampler_windows_test.go:30: ReadTree: resource: proc: system snapshot does not contain subject 17620
    externalsampler_windows_test.go:86: first ReadTree: resource: proc: system snapshot does not contain subject 11528
    systemprocs_windows_test.go:95: parsed 4 of 5 snapshot entries: [60060 4 4004 1337]
    systemprocs_windows_test.go:150: walker delivered [100], want [100 200] (stop after the terminator entry)
    systemprocs_windows_test.go:164: parseSystemProcesses: system process snapshot parsed zero entries
    systemprocs_windows_test.go:210: system snapshot does not contain the sampling process 12328 (513 entries)
```

手工构造的那条红得最干净：`[60060 4 4004 1337]` 四个都在，缺的正是链末项 pid `8080`（`0x1F90`）。
`TestHelperProcess` 的 SKIP 是 Go 子进程测试夹具的固有写法（`jobscope_windows_test.go:87`，
helper 模式未置位时不该断言），与本票无关，**不是**测量断言的静默跳过。

值得单独记一笔：变异态下 A14 把 AC#3 的**树外**读取也打回原形
（`system snapshot does not contain subject <pid>`）—— 被采样的子进程新建后同样落在
`ActiveProcessLinks` 末尾，所以「换个地方读」并不天然免疫这个缺陷，靠的是同一份 `WalkSystemProcesses`。
这就是本票坚持「全仓只留一份走链实现」的理由，代码里那处注释不是洁癖。

## 撤变异后的原始输出（`go test -count=1 ./internal/proc/ ./internal/observe/` → rc=0）

```text
ok  	github.com/CarlosShao/wisp/internal/proc	0.700s
ok  	github.com/CarlosShao/wisp/internal/observe	1.128s
rc=0
```

撤除证据：`grep -c "MUTATION" internal/proc/systemprocs_windows.go` → `0`，
且 `git diff --stat internal/proc/systemprocs_windows.go` 无输出（文件与 `81afe3c` 逐字节一致）。

## 结论（AC#1 复测）

变异 = 缺陷原文顺序；6 条用例转红、其余 24 条仍绿；撤变异全绿且磁盘文件回到 HEAD 状态。
**回归用例确实咬住缺陷本身，而不是咬住某一次链序巧合。**
