# 票 82 — 独立对抗验收裁决表（acceptor-ticket82）

**被验收票**：`.scratch/wisp/issues/82-risk-sync-family-posix-tiering.md`（`Status: coded-awaiting-review`）
**代码落点**：`42ed13f`、`a25430e`、`629ce3c`
**验收基线**：`git archive 629ce3c` → 仓外纯净快照 `C:\Users\swq\AppData\Local\Temp\wisp82acc-a82x`（会话后缀 `a82x`）
**容器**：docker `golang:1.27`（容器内 `go version go1.27.1 linux/amd64`），`WISP_ENV=test`，命名卷 `wisp82acc-a82x-gomod` / `wisp82acc-a82x-gobuild`
**验收者立场**：本表把票 82 的判据当**攻击目标**，重点攻击"分层变成藏东西"。所有数字为**本代理自己那一次读数**，不引用票面数字当证据。

> **本文件是渐进写入的**：下面的"执行状态"块标明哪些组已做完、哪些组没做。
> 若本代理中途死亡，接手者请**只补未完成组**，不要重跑已完成组的结论。

## 执行状态

| 组 | 内容 | 状态 |
|---|---|---|
| A | AC#3(i) ubuntu 侧归零 | **已做完** |
| B | AC#4 变异（剥 tag 两向） | 未做 |
| C | AC#2b 文件名 vs 构建约束 + 两侧集合差异 | 未做 |
| D | AC#1 定性表可验性（抽 3 条读码） | 未做 |
| E | AC#1b `rules_test.go:198` 两侧读数 | 未做 |
| F | AC#5 门禁 + AC#3(ii) Windows 侧 | 未做 |
| G | 假绿四项逐跑点名 | 部分（A 组内已取 SKIP 名单，未定结论） |
| H | 总判词 + 票 70 结论 + 改名建议 | 未做 |

---

## A 组 — AC#3(i)：ubuntu 侧那 8 条是否真的归零

### A.1 命令原文

```
git archive 629ce3c | tar -x -C /tmp/wisp82acc-a82x          # 仓外纯净快照
docker run --rm -v ".../wisp82acc-a82x:/src" \
  -v wisp82acc-a82x-gomod:/go/pkg/mod -v wisp82acc-a82x-gobuild:/root/.cache/go-build \
  -w /src -e WISP_ENV=test golang:1.27 sh -c \
  'go test -count=1 ./internal/risk/ ./internal/tools/'
```

### A.2 逐包结论行原文（`-count=1`，非 -v）

```
ok  	github.com/CarlosShao/wisp/internal/risk	2.401s
ok  	github.com/CarlosShao/wisp/internal/tools	10.477s
```

### A.3 四个数（同快照、同容器、`-count=1 -v`）

| 包 | rc | `=== RUN` 总 | `=== RUN` 顶层 | `--- PASS` 总 | `--- FAIL` | `--- SKIP` | `[no test files]` | `no tests to run` |
|---|---|---|---|---|---|---|---|---|
| `internal/risk` | 0 | **124** | **71** | **123**（顶层 70） | **0** | **1** | 0 | 0 |
| `internal/tools` | 0 | **56** | **44** | **54**（顶层 42） | **0** | **2** | 0 | 0 |

`internal/risk` 的 `--- SKIP` 逐条点名：`TestSyncRegistryProbeLive`（1 条）。
`internal/tools` 的 `--- SKIP` 逐条点名（**不属于本票改动面，但既然是 portable 裸 `go test` 就要记账**）：
`TestD34WriteMatrix`、`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`。

### A.4 那 8 条的逐条读数（ubuntu，`-count=1 -v`）

票面点名的 8 条**一条都不出现在 FAIL 里**（`--- FAIL` 总数 = 0），8/8 逐条 `--- PASS`：

```
--- PASS: TestWriteGatePlainLocalWriteNotFlagged (0.01s)
--- PASS: TestWriteGateNotSelectedByPayloadKey (0.01s)
--- PASS: TestWriteGateEveryPathTargetJudged (0.01s)
--- PASS: TestWriteGateAllPathsNonSyncStaysExempt (0.01s)
--- PASS: TestSyncNormalNewFileWriteNotFlagged (0.01s)
--- PASS: TestSyncFallbackNotDisarmableByWeakRoot (0.01s)
--- PASS: TestSyncEnvConfiguredRoots (0.01s)
--- PASS: TestSyncDotDotTailFailsClosed (0.01s)
```

### A.5 红数没增加

`--- FAIL` = 0，`[no test files]` = 0，`no tests to run` = 0，顶层 RUN = 71（> 0）。
POSIX 层 4 条（`TestSyncNoGradeIsConfirmedOnPosix`、`TestSyncMembershipDecidesOffProfilePosix`、
`TestSyncRegistryProbeIsAStubHere`、`TestExfilSyncWritePosixSpellingInvariant`）在 ubuntu 名单里各出现 1 次；
windows 层 3 条（`TestSyncConfirmedGradesDisarmFallbackWindows`、`TestWriteGatePlainTargetInsideProfileWindows`、
`TestExfilSyncWriteWindowsSpellingInvariant`）在 ubuntu 名单里各出现 **0** 次 —— 见 C 组的两侧集合对比。

### A.6 A 组判词

**AC#3(i) = PASS。** 8 条在 ubuntu 上从红转绿，且转绿是靠"POSIX 上真跑得动"而不是靠搬走：
同一批名字仍出现在 ubuntu 的 `=== RUN` 名单里（A.5 那 8 行本身就是 RUN+PASS 双证据）。
