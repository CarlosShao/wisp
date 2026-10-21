# 票 124 落地腿 124-r2 交件 part3（sec 8..9）

> 前接 `docs/evidence/s1/124-r2-impl.md`（sec 0..4）与 `124-r2-impl-part2.md`（sec 5..7）。

## 8. 门禁（我自己现跑，逐条贴末行；四数分开报，NOT 相减）

跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`（宿主侧）；容器侧同形。

```
$ GOFLAGS= go build ./...              -> rc=0   输出 0 行（无末行）
$ go vet ./internal/tools/             -> rc=0   输出 0 行
$ go test -count=1 ./internal/tools/   -> rc=0   末行逐字：
      ok  	github.com/CarlosShao/wisp/internal/tools	25.864s
      （同形状四数由 -v 那发分开量：RUN=255  TOP-PASS=212  TOP-FAIL=0  TOP-SKIP=4
        SUB-PASS=39  SUB-FAIL=0  SUB-SKIP=0，六个独立计数，无一由相减得出）
$ sh scripts/d22scan.sh                -> rc=0   末行逐字节录：
      d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42,
      ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116,
      ban #8 internal/=610, ban #8 cmd/=144; ban #8 emoji coverage: design/ 39 text files; frontend/ 116 text files;
      internal/ 610 Go files, comments and _test.go included; cmd/ 144 Go files, comments and _test.go included
$ sh scripts/check-path-length-budget.sh -> rc=0 末两行逐字节录：
      check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a
      reason, and the roster equals the tree
      check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 60
      paths on this tree, recomputed just now in both units
      （denominator: tracked paths=17510  over-budget=60  covered by roster=60  not in roster=0）
$ gofmt -l internal/tools              -> 0 行 (rc=0)
$ /d/work/base/gopath/bin/gofumpt.exe -l internal/tools -> 0 行 (rc=0)
      gofumpt 版本按本机实测报 v0.12.0 (go1.27.1)（ci.yml 是 @latest 未钉版本，各批各报自己手上那枚）
```

**台账八 scope 一枚不降**：`d22scan` 那行里 `bans #1-5 internal/=244`、`cmd/=42`、`frontend/=116`、
`internal/tools/=25`、`design/=39`、`frontend/=116`、`internal/=610`、`cmd/=144`
vs 批 1 收尾记的 `…/608/143` => `internal/` 610（+2）、`cmd/` 144（+1）皆**共树漂移**（票 133/134 等兄弟腿的
`_test.go` 进树），本腿**零新增 `.go` 文件**（15 枚全是既有文件的等量替换，`git status` 无 `??` 在 `internal/tools/`）。

**静默侧**：`gofmt`/`gofumpt` 两条尺在**改动落地后的工作树**上跑，皆 0 行；
`.gitattributes` 的 `*.go text eol=lf` 已逐枚核过（15 枚 `file` 输出全 `UTF-8 text`、无 `CRLF`）。

## 9. 红名册逐名

### 9.1 本腿出口（软链形，`POST3`）：余 3 名 = 票 123 豁免（sec 7 已逐名钉死）

```
TestFSReadOnlyNeverOpensACard
TestL1WriteGoesThroughTheRealBlockWindow
TestLateVetoRendersTheApprovalLayersAppliedStepsReport
```

### 9.2 新增红 = 0 名（双向差集）

```
comm -13 PRE-only-in-link.txt POST3-only-in-link.txt   =>  0 行
comm -13 PRE-LINK-fail.txt    POST3-LINK-fail.txt      =>  0 行   (含两形都红那一档：本包 0 枚)
```

### 9.3 变异（`MUTATION-124-R2`）连带翻红的 24 名 —— 逐名 + 首行（同一族 helper 的其它被喂用例，属同族、NOT 新缺陷）

这 24 名**不在** 35 枚分母里（改前软链形就绿），它们由**前一批 2b-2 已转**的共享调用点喂根
（`fsBridge`/`fsBridgeWith`/`taskBridge`/`judgedTaskBridge`/`tier172Calls`/`openStore` 等），
所以摘掉解析层会连带翻红。**逐名报红行**（每名一行，摘自 `MUT-LINK`）：

```
TestAtomicWriteKillsMidWrite
    （子测试 a_clean_write_lands_and_round_trips / new_file_target_does_not_appear_at_all 红，subfail 名单见 sec 9.4）
TestAuditAllowlistScopeDoesNotInventAVerdict
    （子测试 judged_in_range_stays_true 红，subfail 名单见 sec 9.4）
TestCanonicalizeReturnsAPathTheOSCanOpen
    pathshape_portable_test.go:51: InAllowlist("…/001/payload.txt") = false for a file inside the only allowed root
    "…/001": R2 would send every plain write to L2
TestDeclaredRiskIsOnlyAFloor
    bridge_test.go:618: in-allowlist read: out={Text:L2 审批未通过，已拒绝执行 …}
TestFSEditRoutesToApprovalNotTheL1Window
    fs_edit_ac5_gate_r4_test.go:251: in-scope fs.write (new file) must use the L1 window only: window=0 approval=1
TestFSListHonoursItsCap
    fs_test.go:127: out={Text:L2 审批通道尚未接入（票 21），已拒绝执行 …}
TestFSListSummarizesADirectory
    fs_test.go:101: 同形（同上那句）
TestFSMoveSameVolume
    fs_write_test.go:655: out={Text:L2 审批未通过，已拒绝执行 …}
TestFSReadReturnsTheFileAndTaintsIt
    fs_test.go:24: 同形（L2 审批通道尚未接入）
TestFSReadTaintFeedsR4
    fs_test.go:80: tainted exfil judged L0, want L2 (decision={… Level:L0 RulesHit:[] …})
TestFSTrashGoesToTheRecycleBin
    fs_write_test.go:501: the refusal must say why: "L2 审批未通过，已拒绝执行"
TestFSWriteSilentLossIsNotReported
    fswrite_silentloss_ac1_test.go:97: memory.Open: memory: create data dir: winsec: refusing to seal
    /varlink/w124tmp/…/001/data: … path is not provably resolved …
TestLoopPassesDeclaredL1WriteThroughTheGate
    loop_approval_test.go:167: memory.Open: … refusing to seal /varlink/w124tmp/…/001/db: … not provably resolved …
TestLoopStillRefusesL1WhenNoGateIsRegistered
    loop_approval_test.go:236: 同形（…/001/db）
TestPathCanonicalizerAccountsForRewrittenRoots
    paths_rewrite_ticket102_test.go:64: InAllowlist("…/001/proj/a.txt") = false although the expanded root is a
    real tree; the fix must not turn into option (A)
TestRealToolCallWritesRewriteAccountIntoAudit
    bridge_path_account_ticket105_test.go:112: precondition: the read must actually run, got "L2 审批通道尚未接入（票 21），已拒绝执行"
TestSensitiveReadAcrossTheBridgeOpensItsTaskScope
    bridge_scope_open_ticket158_test.go:95: 对照腿 fs.list 没有跑成：out={Text:L2 审批通道尚未接入（票 21），已拒绝执行 …}
TestTicket107AllowlistBoundaryIsComponentWise
    paths_ticket107_portable_test.go:99: InAllowlist("…/001/proj/a.txt") = false, want true: a file under the
    allowed root must pass
TestTicket107AllowlistJudgmentTwoShapes
    paths_ticket107_portable_test.go:57: shape R root="…/001/proj" target="…/001/proj/a.txt" … InAllowlist=false
TestTicket107bProbeCLinkInsideAllowedRootStaysOutside
    paths_ticket107b_probes_test.go:159: premise signals for "…/proj/esc" -> "…/outside": lstat(modeLink=true)
    readlink(…) eval(link)="/realpriv/w124tmp/…/outside"
TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses
    paths_twocontainments_252_r2_test.go:106: 同形 premise signals（…/proj/esc -> …/outside）
TestToolCallRowsAreComplete
    bridge_test.go:651: memory.Open: … refusing to seal /varlink/w124tmp/…/001/data: … not provably resolved …
TestWorkspaceRootAnswersWithTheAccountThatProducedTheNarrowing
    paths_workspace_account_181r3_test.go:83: the narrowed workspace stopped containing its own tree
TestWorkspaceSwitchNarrowsWhatTheAssessorJudges
    paths_workspace_test.go:53: premise: an in-allowlist write with no workspace chosen assessed "L2" ([R1 R2]), want L1
```

两个文案族：机制甲 `not provably resolved`／`refusing to seal`（6 名，点在 `memory.Open` 那侧）、
机制乙 `L2 审批通道尚未接入`／`L2 审批未通过`／`InAllowlist = false`（18 名，走规范化器那一条）。
=> **乙类走的不是 winsec 底线**，本腿按批 1 的 `next=` 只用调用方解析，**未碰 `internal/risk/**`**（零 hunk）。

### 9.4 变异态子测试红名（6 枚，逐名）

```
Test20r2FirstUseAskFlowGoesOnlyThroughConfig/approved:_the_append_lands_via_Manager.AddAllowedDir_and_the_next_assembly_passes
TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips
TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all
TestAuditAllowlistScopeDoesNotInventAVerdict/judged_in_range_stays_true
TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup
TestTicket224GrantSourceThatCannotAnswerMeansAsking/reports_that_it_could_not_look
```

### 9.5 三态名册读数（一句话收敛）

```
PRE      only-in-link = 38   （= 35 可转 + 票 123 三枚豁免）
POST3    only-in-link =  3   （逐名 = 那三枚豁免）        newly-red = 0
MUT      only-in-link = 62   （38 基线逐名全回归 + 24 同族连带）   MUT-PLAIN = 0
RESTORED only-in-link =  3   （与 POST3 逐名双向 0 行差）
```

## 10. 未验证 / 交回编排者

1. **票 123 那三枚 300 s 腿本腿零贡献**：它们在 `wiring_test.go`，本腿一枚未改；`POST3` 软链形它们仍红、
   时延仍 300.04s。=> 本包"软链形红名 = 0"这个票级终判据**只能到 3、到不了 0**，除非票 123 落地。
2. **变异连带 24 名**已逐名登记（sec 9.3），它们是前一批 2b-2 已转调用点的覆盖面，**不是本腿新增红**；
   若裁腿要求"变异只许翻红本腿那 35 名"，我这一发做不到——两批用的是同一族 helper，摘掉就一起红。
3. **本腿不新写用例**：派单要求"每处改动配一枚今天响、不修不响的用例"，本包的兑现方式是
   **那 35 枚本来就是今天响的用例**（`PRE-LINK` 逐名红）+ 反向判据由 `MUTATION-124-R2` 逐名回归证明。
   如果裁腿要的是"新写一枚从没红过的用例"，那超出本腿写面（本包不存在能把软链形失败复现成绿的新用例形状）。
4. **两形 `RUN` 差**：本包四台恒 255、SKIP 恒 4，票面 15:3x 那条副作用**不复现**（已核）。
5. **`cmd/wisp` 与 `internal/config` 不在本腿写面**，本腿零 hunk；`internal/config` 那两枚 POSIX 摇摆不在本包名册。
