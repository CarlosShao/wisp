# 票 124 落地腿 124-r2 交件 part2（sec 5..9）

> 前接 `docs/evidence/s1/124-r2-impl.md`（sec 0..4）。

## 5. 三读数（逐字贴）＋ sha256 还原

台面同 sec 2（`golang:1.27` / `go1.27.1 linux/amd64` / `GOPROXY=off` / **每形一枚新容器** / `--rm` /
`exit 90/97/98` 软链形 ＋ `exit 99` 普通形）。**每一发进容器第一件事**打
`go.mod bytes=1026` + `md5 0c266f4d…` + `resolve.go md5 b6876a5e…`（与票 124 前各批逐字同字）。

```
发            RUN  TOP-PASS TOP-FAIL TOP-SKIP SUB-PASS SUB-FAIL  npr  rts   形状
PRE-LINK      255   174       38       4        36       3        18   18   link
PRE-PLAIN     255   212        0       4        39       0         0    0   plain
POST3-LINK    255   209        3       4        39       0         0    0   link
POST3-PLAIN   255   212        0       4        39       0         0    0   plain
MUT-LINK      255   150       62       4        33       6        22   22   link
MUT-PLAIN     255   212        0       4        39       0         0    0   plain
RESTORED-LINK 255   209        3       4        39       0         0    0   link
```

（`npr` = `not provably resolved` 行数，`rts` = `refusing to seal` 行数；两串**不是**判据，只作机制旁证。
`panic` 与 `test timed out` 七发**各 0**。）

**变异红**（`MUTATION-124-R2`）：把两枚 helper 的函数体改回 `return t.TempDir()`。
**先证落地，再读数**（非测试容器 `w124r2-mutproof`）：

```
$ grep -n 'MUTATION-124-R2'  /src/internal/tools/tempdir_resolved_124_test.go
28:	return t.TempDir() // MUTATION-124-R2
$ grep -c 'return proc.SealableRoot(t.TempDir())'  两枚 helper  =>  0 / 0     (解析那一跳确已摘掉)
$ grep -c 'sealableTempDir124(t)\|sealableTempCanonical124(t)'  调用点   =>  28   (调用点仍在，是重点)
$ go build ./internal/tools/   -> BUILD-RC=0
$ go vet   ./internal/tools/   -> VET-RC=0
```

> 变异体里另加 `var _ = proc.SealableRoot` 守 import —— 否则 `go vet` rc=1 是 **unused import 的编译错**，
> NOT 行为读数（批 1 第一发就撞在这上面，`124-r1-impl-part2.md` sec 6 有账）。加了守之后 rc=0，红全部来自行为。

读数 =>  **软链形 only-in-link 38 -> 62**：**改前那 38 名逐名全部回归**（`comm -23` **0 行**），
另有 **24 名**被同一族 helper 连带翻红（名单与首行见 sec 9）。
**同一发在普通形 0 红**（`MUT-PLAIN` 八数与 `PRE-PLAIN`/`POST3-PLAIN` 逐数相同）=> 这层解析只在软链形起作用。

**还原绿**（`RESTORED-LINK`）：树从 `mut3` 复制后**用提交版两枚 helper 覆盖回去**，先核 sha256 再读数：

```
$ sha256sum internal/tools/tempdir_resolved_124_test.go  internal/tools/tempdir_resolved_124x_test.go
3d949a47ea844aca3f032969d36bc513d974b8dc5dcc7c132fead6f9036ccb6a *internal/tools/tempdir_resolved_124_test.go
9c4b7381d47a629d59aab996e3a5a3541bf893319eead5915eabf79613b22a14 *internal/tools/tempdir_resolved_124x_test.go
$ sha256sum <restored3 的那两枚>            => 同两串，逐字相同
$ git show HEAD:internal/tools/tempdir_resolved_124_test.go  | sha256sum
3d949a47ea844aca3f032969d36bc513d974b8dc5dcc7c132fead6f9036ccb6a *-
$ git show HEAD:internal/tools/tempdir_resolved_124x_test.go | sha256sum
9c4b7381d47a629d59aab996e3a5a3541bf893319eead5915eabf79613b22a14 *-
$ grep -c 'MUTATION-124-R2' <restored3>     => 0 / 0     (变异标记已清)
$ grep -c 'return proc.SealableRoot(t.TempDir())' <restored3> => 1 / 1   (解析那一跳回来了)
$ diff -rq post3 restored3   =>  rc=0（空）  (还原态与"原样"POST3 逐字节同一棵树)
```

=> **还原态 = 提交版 = HEAD blob，三串 sha256 逐字相同**。`RESTORED-LINK` 与 `POST3-LINK` 的 only-in-link
**逐名差集双向 0 行**（`comm -23` 0、`comm -13` 0），八数逐数相同 => **还原绿成立**。

**原样绿**：`POST3-PLAIN` 八数 = `PRE-PLAIN` 八数（`212/0/4` + `39/0`）**逐数相同**，
且 `POST3-PLAIN` 只红 0 枚 => 普通形一枚都没多红。关键按派单**只认名册**：`POST3 only-in-plain` = **0 名**、
`POST3 newly-red` = **0 名**（sec 6），NOT 拿枚数相等当判据。

**新用例**：本腿**不新写用例文件**。派单写面允许"`internal/tools/**` 的 `*_test.go` ＋ 该包**已有**的
sealable-tempdir helper"；反向判据由 sec 5 那发变异逐名兑现——摘掉解析层 => **38 名基线红全数回归、逐名报红行见 sec 9**
（不是新增一枚"看起来覆盖了"的用例）。**零新增走法**：没有第四份解析，两枚 helper 一枚未改（sha256 见上）。

## 6. 出口判据：名册集合（NOT 枚数相等）

```
$ comm -23 POST3-LINK-fail.txt POST3-PLAIN-fail.txt   =>  3 名（= 票 123 三枚豁免，见 sec 7）
$ comm -13 POST3-LINK-fail.txt POST3-PLAIN-fail.txt   =>  0 名
$ went-green  = comm -23 PRE-only-in-link.txt POST3-only-in-link.txt  =>  35 名
$ NEWLY-RED   = comm -13 PRE-only-in-link.txt POST3-only-in-link.txt  =>   0 名   <= 本腿的出口前半句
$ 子测试 only-in-link：PRE 3 -> POST3 0
```

**35 枚 went-green**（逐名，全部在本腿 15 枚改动文件的覆盖面内）：
Test165PlanApprovalRingsAgainForAnUnlistedWrite, Test172DocCopyThatLiesReddensTheRuler,
Test172FrozenTierMatchesTheTierTheGateGives, Test20r2DirectionEngineAdjudicatesFileLevelAdditions,
Test20r2FirstUseAskFlowGoesOnlyThroughConfig, Test221ParentStopsItsOwnChildRowAndStreamSettle,
Test283IdentityChainThroughTheRealLoop, Test285RosterRowsCarryDistinctCorrPerCall,
Test287DepthRefusalWinsOverThePoolRefusal, Test287RootCallerWithTheSameShapeStillDerives,
Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop,
TestBackfillFilesTheRealSpilledArtifact176r1, TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2,
TestHealthyPointerStaysSilent, TestHealthyReplyStillMatchesThePreFixTemplate,
TestHostMintedPointerRereadStaysCleanOnRealBridge175r2, TestLeakProbeTokenOnALegalPathStaysSilent174r3,
TestLongOutputPointerRecoversEveryByte, TestNonContentResultStaysUnmarked175r2,
TestPointer183BackfilledArtifactRereadStaysClean, TestPointer183SiblingArtifactInTheSameDirectoryStillHits,
TestPointer185AModelNominatedFileIsNotRostered, TestPointer185AnotherTaskCannotBorrowTheRoster,
TestPointer185ForeignFileCarryingTheHostPathStillBlocksTheReread, TestPointer185PagedRereadsKeepWorking,
TestPointer185SecondRereadOfTheSameArtifactStaysClean, TestPointer185SiblingHostArtifactIsNotLaundered,
TestPointerPast256KiBIsNotFullyReadable, TestPointerToMissingCopyFileSpeaks,
TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1, TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow,
TestTicket224GrantIDColumnRoundTripsThroughTheRealStore, TestTicket224GrantSourceThatCannotAnswerMeansAsking,
TestTicket224LiveGrantStopsTheL1Question, TestTicket224PartialPathCoverageStillAsks

**"变绿 vs 被跳过"逐枚排除**：`POST3-PLAIN` 里这 35 名**逐名 = PASS**（不在 skip.txt 里）；
四台 `SKIP` 名册**逐名同一组 4 枚**（`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`、
`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`、`TestD34WriteMatrix`、
`TestWorkspaceSwitchRefusesAJunctionToOutside`，皆 windows-tag 用例在 POSIX 上自跳，**两形相等**、不属"软链形自缩分母"那档）；
四台 `RUN` 行恒 **255**、`SUBPASS+SUBFAIL` 四台 `39/3 -> 39/0` 无"子测试从未被创建"的缩水。

**两形 `RUN`/`SKIP` 差逐包解释**：本包四台 `RUN` **恒 255**、顶层 `SKIP` **恒 4**、子 `SKIP` **恒 0**
=> 票面 15:3x 那条"软链形本身会缩小分母"在 `internal/tools` **不复现**（软链形红落在用例体内，NOT 死在父用例的 `Open`）。

## 7. 豁免（逐名，注明"归票 123"）

`POST3 only-in-link` 余 **3 名**，逐名即票 123 那族 300 s 腿，本发一手量到时延 `300.04s` 各一枚：

| 名 | 行号 | 时延 | 归口 |
|---|---|---|---|
| `TestL1WriteGoesThroughTheRealBlockWindow` | `wiring_test.go:112` | 300.04s | **归票 123** |
| `TestLateVetoRendersTheApprovalLayersAppliedStepsReport` | `wiring_test.go:182` | 300.04s | **归票 123** |
| `TestFSReadOnlyNeverOpensACard` | `wiring_test.go:291` | 300.04s | **归票 123** |

三枚**按账豁免**：NOT 计入 35 枚分母、NOT 顺手修、NOT 放宽断言。失败文案逐字是
`审批超时（300 秒未确认），C18 一律判拒绝` / `the report is missing "取消不是原子的"`，
**都不是「未解析根」那一族**（`npr`/`rts` 在 `POST3-LINK` 已归 0，而这三枚照红）。
`wiring_test.go` 本腿**一枚未改**（不在 sec 3 那 15 枚里；`git diff --name-only` 无它）。

**另登记一条与本腿零因果的既有账**：`internal/config` 那两枚 POSIX-only 摇摆
（`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`／`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`）
**不在本包名册**，本腿未碰、未计入任何出口判据。

> 本件续 `124-r2-impl-part3.md`（sec 8..10：门禁 rc、红名册逐名、未验证项）
> 与其后的 `124-r2-impl-part4.md`（sec 11：提交号、逐枚 `.go`、越界与正控）。
