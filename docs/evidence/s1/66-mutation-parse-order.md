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
