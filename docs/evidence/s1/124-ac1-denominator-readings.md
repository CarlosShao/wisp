# 票 124 AC#1 — 分母读数（POSIX 容器 + 软链 `TMPDIR`）

**实现方**：`worker-ticket124-ac1` · **日期**：2026-09-23
**票面判据（AC#1 原文）**：容器内以软链 TMPDIR 跑相关包，逐包 `RUN/PASS/FAIL/SKIP` 四数点名 + 红名清单
（**多样本全报，不许只报一次**），并证明同一批用例在 plain 形下 0 红 ⇒ 排除「其实是别的形状红」。
**本轮范围**：只做 AC#1 这一格，**只测不改**（见 §6）。

## 0. 锚点与快照

- 锚定 sha：`7b4c36aa9f69960043806f6267e107aed1688ac7`（`git rev-parse HEAD`，取快照当刻）⇒ 下表全部读数只对这一枚 sha 负责。
- 快照：`git archive <锚> | tar -x -C /d/tmp/w124-ac1-snap` —— **纯净快照，不含共享工作树里他人的未提交半成品**
  （取快照当刻 `git status --porcelain` 里 `cmd/wisp/slo_windows.go` 等是别人的活，一律没进快照）。
  快照共 1030 枚文件；容器内 `md5sum /src/go.mod = f6ef661732b1851e5c3db348113cb605`、
  `md5sum /src/internal/winsec/resolve.go = b6876a5efe759f6e17434d1b50a129c3`。
- 镜像：`golang:1.27`（容器内 `go version go1.27.1 linux/amd64`，uid=0，`Linux ... 6.6.114.1-microsoft-standard-WSL2`）。
- 依赖：复用票 119 的命名卷 `ac119-gomodcache` / `ac119-gocache` ⇒ 容器内**不需要网络**即可编译。
- 开测前查在飞 run（本机即 self-hosted runner）：`gh run list --repo CarlosShao/wisp --limit 5` 当时有 2 枚 `in_progress`
  （run `35826783905`、`35826548877`）⇒ 样本 L1 的前半与它们**同机抢 CPU**（12 枚逻辑 CPU），两枚在 14:36 前后进终态，
  L2/P1 落在无竞争窗口。红/绿不受这件事影响，但**单枚耗时**（尤其 `./internal/tools/`）在 L1 里偏悲观。

## 1. 仪器与跑法（「可重跑」就是这一格）

⚠ **容器挂载的假绿陷阱**（本仓踩过、简报再次点过）：Git Bash 下把 -v 的反斜杠形式路径挂进容器会**静默挂空且 rc=0**（详见下方跑法：一律 `/d/...` + `MSYS_NO_PATHCONV=1`）。
⇒ 本票一律用 `/d/...` 形式 + `MSYS_NO_PATHCONV=1`，并且每枚样本进容器第一件事就是 `ls -l /src/go.mod` 自证挂载真生效
（§2/§3/§4 各贴了那行原文）。

跑法（宿主 Git Bash，一条命令一枚样本）：

```bash
export MSYS_NO_PATHCONV=1
docker run --rm \
  -v /d/tmp/w124-ac1-snap:/src \
  -v /d/tmp/w124-ac1-logs:/logs \
  -v /d/tmp/w124-ac1-h.sh:/h.sh:ro \
  -v ac119-gomodcache:/go/pkg/mod \
  -v ac119-gocache:/root/.cache/go-build \
  golang:1.27 bash /h.sh <link|plain> <TAG> > /d/tmp/w124-ac1-logs/<TAG>.driver.txt 2>&1
```

`/h.sh` 逐字（下一位照抄即可复跑；它只读 `/src`，写只写 `/logs`）：

```bash
#!/usr/bin/env bash
# Ticket 124 AC#1 denominator harness (POSIX, container).
# Runs every ./internal/... and ./cmd/... package with -count=1 -v under one
# TMPDIR shape and books RUN/PASS/FAIL/SKIP per package plus the red names.
#
#   usage: bash /h.sh <link|plain> <tag>
#     link  -> mkdir /realpriv/w124tmp ; ln -s /realpriv /varlink ; TMPDIR=/varlink/w124tmp
#     plain  -> mkdir -p /plainroot/w124tmp ; TMPDIR=/plainroot/w124tmp   (no symlink anywhere)
#
# Nothing here writes to /src; logs go to /logs.

set -u

SHAPE="$1"
TAG="$2"
LOGDIR=/logs
SUM="$LOGDIR/$TAG.summary.txt"

assert_link() {
  case "$(ls -ld /varlink)" in
    l*) ;;
    *) echo "SHAPE_ASSERT_FAILED /varlink is not a symlink: $(ls -ld /varlink)"; exit 97 ;;
  esac
  if [ "$(readlink -f /varlink/w124tmp)" != "/realpriv/w124tmp" ]; then
    echo "SHAPE_ASSERT_FAILED readlink -f /varlink/w124tmp = $(readlink -f /varlink/w124tmp)"
    exit 98
  fi
}

echo "== mount self-proof =="
ls -l /src/go.mod
md5sum /src/go.mod /src/internal/winsec/resolve.go
go version
uname -sra

cd /src || exit 90

if [ "$SHAPE" = link ]; then
  mkdir -p /realpriv/w124tmp
  ln -s /realpriv /varlink || exit 96
  assert_link
  export TMPDIR=/varlink/w124tmp
  echo "SHAPE=$SHAPE TMPDIR=$TMPDIR resolved=$(readlink -f "$TMPDIR")"
  ls -ld /varlink
else
  mkdir -p /plainroot/w124tmp
  export TMPDIR=/plainroot/w124tmp
  echo "SHAPE=$SHAPE TMPDIR=$TMPDIR resolved=$(readlink -f "$TMPDIR")"
  ls -ld /plainroot /plainroot/w124tmp
fi

export WISP_ENV=test
export HOME="${HOME:-/root}"
echo "WISP_ENV=$WISP_ENV HOME=$HOME"

: > "$SUM"

PKGS=$(go list ./internal/... ./cmd/... 2>/dev/null | sed 's#^github.com/CarlosShao/wisp/#./#; s#/$##; s#$#/#')
echo "== package list ($(echo "$PKGS" | wc -l) packages) =="
echo "$PKGS"

for p in $PKGS; do
  slug=$(echo "$p" | tr '/.' '__')
  f="$LOGDIR/$TAG.$slug.txt"
  go test -count=1 -v -timeout 10m "$p" > "$f" 2>&1
  rc=$?
  run=$(grep -c '^=== RUN' "$f" || true)
  pass=$(grep -c '^--- PASS' "$f" || true)
  fail=$(grep -c '^--- FAIL' "$f" || true)
  skip=$(grep -c '^--- SKIP' "$f" || true)
  subfail=$(grep -c '^[[:space:]][[:space:]]*--- FAIL' "$f" || true)
  buildsfail=$(grep -c '^FAIL.*\[build failed\]' "$f" || true)
  printf '%s|rc=%s|RUN=%s|PASS=%s|FAIL=%s|SKIP=%s|SUBFAIL=%s|BUILDFAIL=%s\n' \
    "$p" "$rc" "$run" "$pass" "$fail" "$skip" "$subfail" "$buildsfail" >> "$SUM"
done

echo "== summary $TAG ($SHAPE) =="
cat "$SUM"
echo "== done =="
```

要点：

- 样本 `L1`/`L2` = 软链形两枚独立样本（各一枚全新容器、各一次 `-count=1`）；`P1` = plain 形对照。
  **同一镜像、同一快照、同一 `WISP_ENV=test`、同一 `HOME=/root`，唯一变量是 `TMPDIR` 经不经软链。**
- 四数口径（与票 119/111/125 同一把尺）：`RUN` 数 `^=== RUN`（**含子测试**）；`PASS`/`FAIL`/`SKIP` 数 `^--- PASS`、`^--- FAIL`、`^--- SKIP`（**只顶层**，所以 `RUN != PASS+FAIL+SKIP` 是正常形状）；子测试的红另立一列（`^\s+--- FAIL`），并在红名清单里逐枚点名、标「顶层/子测试」。
- 包面：`go list ./internal/... ./cmd/...` = **30 枚包**逐枚单独跑（不合并跑，避免多包 `-v` 输出交错导致归属不清）。
  票 119 当年只点名「六包」；本票把**全部** `internal/**` + `cmd/**` 拉进分母，就是为了回答「还有谁在红」。
- 形状硬断言（票 119 踩过的那枚坑：`mkdir -p /realpriv /varlink` 会把 `/varlink` 造成**真目录**、链接根本不存在）：
  脚本里 `case "$(ls -ld /varlink)" in l*) ;; *) exit 97;; esac` 且 `readlink -f /varlink/w124tmp` 必须等于
  `/realpriv/w124tmp`；断言不过就当停。每枚样本的驱动原文里都有那两行结果。
- `git archive` 的纯净快照不含未入库的 `third_party/*.dll` ⇒ 那是 **Windows** 腿的坑（`0xc0000135`）；
  本票全在 Linux 容器里跑，`./cmd/wisp/` 正常编译并产出四数（见各表末行），未拷 DLL。

## 2. 软链形样本 L1（第一枚）

形状：容器内 `mkdir -p /realpriv/w124tmp && ln -s /realpriv /varlink && export TMPDIR=/varlink/w124tmp`；`t.TempDir()` 于是长成 `/varlink/w124tmp/TestXxx<随机数>/001`，**调用方一个字都没解析**。

驱动原文开头（挂载自证 + 形状自证，逐字）：

```
== mount self-proof ==
-rwxrwxrwx 1 root root 883 Sep 23 06:26 /src/go.mod
f6ef661732b1851e5c3db348113cb605  /src/go.mod
b6876a5efe759f6e17434d1b50a129c3  /src/internal/winsec/resolve.go
go version go1.27.1 linux/amd64
Linux bc9ab58f3038 6.6.114.1-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Mon Dec  1 20:46:23 UTC 2025 x86_64 GNU/Linux
SHAPE=link TMPDIR=/varlink/w124tmp resolved=/realpriv/w124tmp
lrwxrwxrwx 1 root root 9 Sep 23 06:32 /varlink -> /realpriv
WISP_ENV=test HOME=/root
== package list (30 packages) ==
```

| 包 | rc | RUN | PASS | FAIL | SKIP | 子测试 FAIL | 构建失败 |
|---|---|---|---|---|---|---|---|
| `./internal/agent/` | 1 | 75 | 40 | 18 | 0 | 8 | 0 |
| `./internal/agent/approval/` | 1 | 49 | 28 | 1 | 1 | 0 | 0 |
| `./internal/agent/scheduler/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/audio/` | 0 | 18 | 18 | 0 | 0 | 0 | 0 |
| `./internal/ball/` | 0 | 46 | 40 | 0 | 0 | 0 | 0 |
| `./internal/buildinfo/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./internal/config/` | 1 | 99 | 48 | 6 | 1 | 1 | 0 |
| `./internal/llm/` | 1 | 72 | 52 | 11 | 0 | 6 | 0 |
| `./internal/llm/adaptertest/` | 0 | 14 | 4 | 0 | 0 | 0 | 0 |
| `./internal/llm/anthropic/` | 0 | 27 | 9 | 0 | 0 | 0 | 0 |
| `./internal/llm/golden/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/llm/openaichat/` | 0 | 38 | 16 | 0 | 0 | 0 | 0 |
| `./internal/llm/openairesponses/` | 0 | 26 | 8 | 0 | 0 | 0 | 0 |
| `./internal/memory/` | 1 | 36 | 2 | 33 | 1 | 0 | 0 |
| `./internal/models/` | 0 | 51 | 38 | 0 | 2 | 0 | 0 |
| `./internal/observe/` | 0 | 54 | 54 | 0 | 0 | 0 | 0 |
| `./internal/panel/` | 1 | 76 | 45 | 1 | 0 | 0 | 0 |
| `./internal/perm/` | 1 | 14 | 9 | 5 | 0 | 0 | 0 |
| `./internal/plugin/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/proc/` | 0 | 16 | 12 | 0 | 0 | 0 | 0 |
| `./internal/risk/` | 0 | 138 | 77 | 0 | 1 | 0 | 0 |
| `./internal/secret/` | 0 | 4 | 4 | 0 | 0 | 0 | 0 |
| `./internal/session/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/speech/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/statemachine/` | 0 | 73 | 10 | 0 | 0 | 0 | 0 |
| `./internal/tools/` | 1 | 77 | 40 | 19 | 3 | 2 | 0 |
| `./internal/watchdog/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/winsec/` | 1 | 45 | 14 | 13 | 3 | 4 | 0 |
| `./cmd/llmrecord/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./cmd/wisp/` | 1 | 70 | 15 | 25 | 0 | 8 | 0 |
| **合计（30 枚包）** | - | **1136** | **601** | **132** | **12** | **29** | **0** |

FAIL 行合计（顶层 + 子测试 = 票 119 数「81 条 FAIL 行」那把尺子的口径）= **161**。

### L1 红名清单（逐枚点名，顶层与子测试分列）

**`./internal/agent/`** — 顶层 18 枚 / 子测试 8 枚
- `TestFailedTaskBooksOpenCallRowWithDecision`（顶层）
- `TestCancelledTaskPersistsTerminalRows`（顶层）
- `TestPerToolTimeoutFires`（顶层）
- `TestPerToolTimeoutOfContractHonestToolIsToolClass`（顶层）
- `TestTaskLogAndToolCallRows`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide`（顶层）
- `TestWriteFileExclusiveIsExclusive`（顶层）
- `TestSpillSameIDRetryOverwrites`（顶层）
- `TestSpillAcrossRestartsKeepsRetrySemantics`（顶层）
- `TestSpillCallIDHostileShapesSanitizedToBareNames`（顶层）
- `TestSpillContainmentByDirectoryListing`（顶层）
- `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`（顶层）
- `TestSpillThresholdScalesWithWindow`（顶层）
- `TestSpillTokenBoundary`（顶层）
- `TestSpillArtifactAndStubShape`（顶层）
- `TestSpillThroughLoop`（顶层）
- `TestSpillArtifactRespectsRawCap`（顶层）
- `TestMaxTokensFailsAllToolCallsOfThatMessage`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/separator`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/unc`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence`（子测试）

**`./internal/agent/approval/`** — 顶层 1 枚 / 子测试 0 枚
- `TestTenOpsInOneToolCallGetOneConfirm`（顶层）

**`./internal/config/`** — 顶层 6 枚 / 子测试 1 枚
- `TestRoundTripLoadMarshalLoad`（顶层）
- `TestResolvedNeverPersists`（顶层）
- `TestSaveFileRoundTripsThroughLoad`（顶层）
- `TestMigrateV1Fixture`（顶层）
- `TestMigrateNoVersionKeyAssumedV1`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone/SaveFile_output_reloads`（子测试）

**`./internal/llm/`** — 顶层 11 枚 / 子测试 6 枚
- `TestProbeSuiteMeasuresBrokenFC`（顶层）
- `TestProbeSuiteMeasuresBrokenVision`（顶层）
- `TestProbeSuiteHonestProviderRecordsNoMismatch`（顶层）
- `TestProbeSuiteHonestNegativeIsNotAMismatch`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects`（顶层）
- `TestProbeSuiteCapableOnAllThreeDialects`（顶层）
- `TestProbeSuiteRefusesSilentRuns`（顶层）
- `TestProbeSuiteSinkFailurePropagates`（顶层）
- `TestProbeSuiteAudioStaysUnprobed`（顶层）
- `TestProbeSuiteMeasuresBrokenThinking`（顶层）
- `TestProbeSuiteThinkingCapableIsNotTheSameAsBroken`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-responses`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-responses`（子测试）

**`./internal/memory/`** — 顶层 33 枚 / 子测试 0 枚
- `TestDeleteArtifactRejectsTheFourHostileShapes`（顶层）
- `TestDeletePrivacyItemRejectsTheFourHostileShapes`（顶层）
- `TestArtifactsContainmentByDirectoryListing`（顶层）
- `TestArtifactsLiteralBackslashKeepsItsAsymmetry`（顶层）
- `TestStraySubdirectoryCannotHideBytesFromTheQuota`（顶层）
- `TestStraySubdirectorySharesTheQuotasLRUQueue`（顶层）
- `TestPurgeArtifactsReclaimsStraySubdirectory`（顶层）
- `TestStrayRemovalDoesNotFollowLinks`（顶层）
- `TestConcurrentWritersReaders`（顶层）
- `TestCrashRecoveryKillMidWrite`（顶层）
- `TestProviderHealthProbeRoundTrip`（顶层）
- `TestProviderHealthRecordErrorAndList`（顶层）
- `TestProviderHealthMissingRowIsNoRows`（顶层）
- `TestProfileUpsertAndLRUEvictionLogs`（顶层）
- `TestProfileDeleteAndPurge`（顶层）
- `TestMemoryAddSearchDeletePurge`（顶层）
- `TestTaskLogLifecycleAndErrorClassValidation`（顶层）
- `TestToolCallLifecycleAndValidation`（顶层）
- `TestGrantCostPluginStateDAO`（顶层）
- `TestWriteQueueLazyLifecycle`（顶层）
- `TestAdvPanicWedge`（顶层）
- `TestMigrateSeededV1DatabaseToV2`（顶层）
- `TestPrivacyOpsAllDomains`（顶层）
- `TestRetentionBoundaries`（顶层）
- `TestRetentionSchedule`（顶层）
- `TestArtifactsLRUQuota`（顶层）
- `TestArtifactsDeleteAndPurgeAndGuards`（顶层）
- `TestSchemaContractIntrospection`（顶层）
- `TestMigrateFreshDatabaseCreatesCurrentVersion`（顶层）
- `TestMigrationChainWithBackup`（顶层）
- `TestMigrationFailedStepIsAtomic`（顶层）
- `TestMigrationNewerSchemaIsUnmigratable`（顶层）
- `TestMigrationForeignDatabaseIsUnmigratable`（顶层）

**`./internal/panel/`** — 顶层 1 枚 / 子测试 0 枚
- `TestComposerRenderFixtureTellsTheTruth`（顶层）

**`./internal/perm/`** — 顶层 5 枚 / 子测试 0 枚
- `TestTicket90ManualSwitchSurvivesRestart`（顶层）
- `TestTicket90UntouchedConfigStartsAtTheDefault`（顶层）
- `TestTicket90SessionGrantDoesNotSurviveRestart`（顶层）
- `TestTicket90ConfigKeyChangesWhatTheChainAsks`（顶层）
- `TestTicket90HandEditLooseningGoesThroughD36`（顶层）

**`./internal/tools/`** — 顶层 19 枚 / 子测试 2 枚
- `TestRealToolCallWritesRewriteAccountIntoAudit`（顶层）
- `TestDeclaredRiskIsOnlyAFloor`（顶层）
- `TestToolCallRowsAreComplete`（顶层）
- `TestFSReadReturnsTheFileAndTaintsIt`（顶层）
- `TestFSReadTaintFeedsR4`（顶层）
- `TestFSListSummarizesADirectory`（顶层）
- `TestFSListHonoursItsCap`（顶层）
- `TestAtomicWriteKillsMidWrite`（顶层）
- `TestFSTrashGoesToTheRecycleBin`（顶层）
- `TestFSMoveSameVolume`（顶层）
- `TestPathCanonicalizerAccountsForRewrittenRoots`（顶层）
- `TestTicket107AllowlistJudgmentTwoShapes`（顶层）
- `TestTicket107AllowlistBoundaryIsComponentWise`（顶层）
- `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`（顶层）
- `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges`（顶层）
- `TestCanonicalizeReturnsAPathTheOSCanOpen`（顶层）
- `TestLoopPassesDeclaredL1WriteThroughTheGate`（顶层）
- `TestLoopStillRefusesL1WhenNoGateIsRegistered`（顶层）
- `TestL1WriteGoesThroughTheRealBlockWindow`（顶层）
- `TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all`（子测试）
- `TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips`（子测试）

**`./internal/winsec/`** — 顶层 13 枚 / 子测试 4 枚
- `TestAC5FailedSealRefusesTheWrite`（顶层）
- `TestAC5FailureIsNotSwallowedByTheHappyPath`（顶层）
- `TestPOSIXPrivateFileIsReally0600`（顶层）
- `TestPOSIXPrivateDirIsReally0700`（顶层）
- `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed`（顶层）
- `TestPOSIXMissingFileIsNotAnError`（顶层）
- `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC4POSIXFloorAnswersInsideTheNamedTree`（顶层）
- `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`（顶层）
- `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`（顶层）
- `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC5FailedSealRefusesTheWrite/exclusive_artifact`（子测试）
- `TestAC5FailedSealRefusesTheWrite/replacement_write`（子测试）
- `TestAC5FailedSealRefusesTheWrite/directory_chain`（子测试）
- `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct`（子测试）

**`./cmd/wisp/`** — 顶层 25 枚 / 子测试 8 枚
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128`（顶层）
- `TestAC4EveryLegIsNailedOrRuled`（顶层）
- `TestProvidersProbeRecordsMeasuredThinkingFalse`（顶层）
- `TestProvidersProbeRecordsMeasuredThinkingTrue`（顶层）
- `TestProvidersDiscoverListsWhatTheServerServes`（顶层）
- `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`（顶层）
- `TestTicket101ManualSwitchSurvivesRestart`（顶层）
- `TestTicket101UntouchedConfigRestartsAtDefault`（顶层）
- `TestTicket101SessionGrantDoesNotCrossRestart`（顶层）
- `TestTicket101ModeSwitchUsesTheRealL2Gate`（顶层）
- `TestTicket101UnreadableModeFailsLoudlyAndStrict`（顶层）
- `TestRunTextTaskTextPathEndToEnd`（顶层）
- `TestRunTextTaskFailNextIsClassified`（顶层）
- `TestHostDispatchThroughTheAssembledBridge`（顶层）
- `TestComposedGateBlocksAWriteForTwoSeconds`（顶层）
- `TestRunTextTaskKeyResolvesInTheStore`（顶层）
- `TestMissingBlobFailsUnconfiguredNeverSilently`（顶层）
- `TestSecretSetGetListRoundTrip`（顶层）
- `TestSecretFromStdinWritesNoIntermediateFile`（顶层）
- `TestSecretFailurePathsLogAndPrintNoPlaintext`（顶层）
- `TestSecretUnsetRefusesWhileReferenced`（顶层）
- `TestSecretSameNameUnderThreeEnvsIsThreeBlobs`（顶层）
- `TestSecretPortableModeUsesTicket06Seam`（顶层）
- `TestSecretEndToEndConfigRefResolvesAtRequestTime`（顶层）
- `TestSecretOverwriteIsAnnounced`（顶层）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only`（子测试）


## 3. 软链形样本 L2（第二枚 —— 票面「多样本全报，不许只报一次」）

**两枚样本的关系（先量再写）**：`diff /d/tmp/w124-ac1-logs/L1.summary.txt /d/tmp/w124-ac1-logs/L2.summary.txt` **rc=0**
⇒ 30 行摘要（逐包 rc/RUN/PASS/FAIL/SKIP）完全相同；再把 30 份 `-v` 日志里的 `--- FAIL` / `--- SKIP` 行**去掉耗时后缀**逐名比对
⇒ **枚枚相同**（唯一差异就是 `(0.01s)` 这类耗时）。所以 L2 不是"另一个数"，而是同一副面的第二次复现；
下面 §3 的表与 §2 的表因此同数（161 枚 FAIL 行、六包口径 83 枚）。

驱动原文开头（挂载自证 + 形状自证，逐字）：

```
== mount self-proof ==
-rwxrwxrwx 1 root root 883 Sep 23 06:26 /src/go.mod
f6ef661732b1851e5c3db348113cb605  /src/go.mod
b6876a5efe759f6e17434d1b50a129c3  /src/internal/winsec/resolve.go
go version go1.27.1 linux/amd64
Linux 5d6d9b1e8221 6.6.114.1-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Mon Dec  1 20:46:23 UTC 2025 x86_64 GNU/Linux
SHAPE=link TMPDIR=/varlink/w124tmp resolved=/realpriv/w124tmp
lrwxrwxrwx 1 root root 9 Sep 23 06:46 /varlink -> /realpriv
WISP_ENV=test HOME=/root
== package list (30 packages) ==
```

| 包 | rc | RUN | PASS | FAIL | SKIP | 子测试 FAIL | 构建失败 |
|---|---|---|---|---|---|---|---|
| `./internal/agent/` | 1 | 75 | 40 | 18 | 0 | 8 | 0 |
| `./internal/agent/approval/` | 1 | 49 | 28 | 1 | 1 | 0 | 0 |
| `./internal/agent/scheduler/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/audio/` | 0 | 18 | 18 | 0 | 0 | 0 | 0 |
| `./internal/ball/` | 0 | 46 | 40 | 0 | 0 | 0 | 0 |
| `./internal/buildinfo/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./internal/config/` | 1 | 99 | 48 | 6 | 1 | 1 | 0 |
| `./internal/llm/` | 1 | 72 | 52 | 11 | 0 | 6 | 0 |
| `./internal/llm/adaptertest/` | 0 | 14 | 4 | 0 | 0 | 0 | 0 |
| `./internal/llm/anthropic/` | 0 | 27 | 9 | 0 | 0 | 0 | 0 |
| `./internal/llm/golden/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/llm/openaichat/` | 0 | 38 | 16 | 0 | 0 | 0 | 0 |
| `./internal/llm/openairesponses/` | 0 | 26 | 8 | 0 | 0 | 0 | 0 |
| `./internal/memory/` | 1 | 36 | 2 | 33 | 1 | 0 | 0 |
| `./internal/models/` | 0 | 51 | 38 | 0 | 2 | 0 | 0 |
| `./internal/observe/` | 0 | 54 | 54 | 0 | 0 | 0 | 0 |
| `./internal/panel/` | 1 | 76 | 45 | 1 | 0 | 0 | 0 |
| `./internal/perm/` | 1 | 14 | 9 | 5 | 0 | 0 | 0 |
| `./internal/plugin/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/proc/` | 0 | 16 | 12 | 0 | 0 | 0 | 0 |
| `./internal/risk/` | 0 | 138 | 77 | 0 | 1 | 0 | 0 |
| `./internal/secret/` | 0 | 4 | 4 | 0 | 0 | 0 | 0 |
| `./internal/session/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/speech/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/statemachine/` | 0 | 73 | 10 | 0 | 0 | 0 | 0 |
| `./internal/tools/` | 1 | 77 | 40 | 19 | 3 | 2 | 0 |
| `./internal/watchdog/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/winsec/` | 1 | 45 | 14 | 13 | 3 | 4 | 0 |
| `./cmd/llmrecord/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./cmd/wisp/` | 1 | 70 | 15 | 25 | 0 | 8 | 0 |
| **合计（30 枚包）** | - | **1136** | **601** | **132** | **12** | **29** | **0** |

FAIL 行合计（顶层 + 子测试 = 票 119 数「81 条 FAIL 行」那把尺子的口径）= **161**。

### L2 红名清单（逐枚点名，顶层与子测试分列）

**`./internal/agent/`** — 顶层 18 枚 / 子测试 8 枚
- `TestFailedTaskBooksOpenCallRowWithDecision`（顶层）
- `TestCancelledTaskPersistsTerminalRows`（顶层）
- `TestPerToolTimeoutFires`（顶层）
- `TestPerToolTimeoutOfContractHonestToolIsToolClass`（顶层）
- `TestTaskLogAndToolCallRows`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide`（顶层）
- `TestWriteFileExclusiveIsExclusive`（顶层）
- `TestSpillSameIDRetryOverwrites`（顶层）
- `TestSpillAcrossRestartsKeepsRetrySemantics`（顶层）
- `TestSpillCallIDHostileShapesSanitizedToBareNames`（顶层）
- `TestSpillContainmentByDirectoryListing`（顶层）
- `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`（顶层）
- `TestSpillThresholdScalesWithWindow`（顶层）
- `TestSpillTokenBoundary`（顶层）
- `TestSpillArtifactAndStubShape`（顶层）
- `TestSpillThroughLoop`（顶层）
- `TestSpillArtifactRespectsRawCap`（顶层）
- `TestMaxTokensFailsAllToolCallsOfThatMessage`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/separator`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/unc`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence`（子测试）

**`./internal/agent/approval/`** — 顶层 1 枚 / 子测试 0 枚
- `TestTenOpsInOneToolCallGetOneConfirm`（顶层）

**`./internal/config/`** — 顶层 6 枚 / 子测试 1 枚
- `TestRoundTripLoadMarshalLoad`（顶层）
- `TestResolvedNeverPersists`（顶层）
- `TestSaveFileRoundTripsThroughLoad`（顶层）
- `TestMigrateV1Fixture`（顶层）
- `TestMigrateNoVersionKeyAssumedV1`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone/SaveFile_output_reloads`（子测试）

**`./internal/llm/`** — 顶层 11 枚 / 子测试 6 枚
- `TestProbeSuiteMeasuresBrokenFC`（顶层）
- `TestProbeSuiteMeasuresBrokenVision`（顶层）
- `TestProbeSuiteHonestProviderRecordsNoMismatch`（顶层）
- `TestProbeSuiteHonestNegativeIsNotAMismatch`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects`（顶层）
- `TestProbeSuiteCapableOnAllThreeDialects`（顶层）
- `TestProbeSuiteRefusesSilentRuns`（顶层）
- `TestProbeSuiteSinkFailurePropagates`（顶层）
- `TestProbeSuiteAudioStaysUnprobed`（顶层）
- `TestProbeSuiteMeasuresBrokenThinking`（顶层）
- `TestProbeSuiteThinkingCapableIsNotTheSameAsBroken`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-responses`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-responses`（子测试）

**`./internal/memory/`** — 顶层 33 枚 / 子测试 0 枚
- `TestDeleteArtifactRejectsTheFourHostileShapes`（顶层）
- `TestDeletePrivacyItemRejectsTheFourHostileShapes`（顶层）
- `TestArtifactsContainmentByDirectoryListing`（顶层）
- `TestArtifactsLiteralBackslashKeepsItsAsymmetry`（顶层）
- `TestStraySubdirectoryCannotHideBytesFromTheQuota`（顶层）
- `TestStraySubdirectorySharesTheQuotasLRUQueue`（顶层）
- `TestPurgeArtifactsReclaimsStraySubdirectory`（顶层）
- `TestStrayRemovalDoesNotFollowLinks`（顶层）
- `TestConcurrentWritersReaders`（顶层）
- `TestCrashRecoveryKillMidWrite`（顶层）
- `TestProviderHealthProbeRoundTrip`（顶层）
- `TestProviderHealthRecordErrorAndList`（顶层）
- `TestProviderHealthMissingRowIsNoRows`（顶层）
- `TestProfileUpsertAndLRUEvictionLogs`（顶层）
- `TestProfileDeleteAndPurge`（顶层）
- `TestMemoryAddSearchDeletePurge`（顶层）
- `TestTaskLogLifecycleAndErrorClassValidation`（顶层）
- `TestToolCallLifecycleAndValidation`（顶层）
- `TestGrantCostPluginStateDAO`（顶层）
- `TestWriteQueueLazyLifecycle`（顶层）
- `TestAdvPanicWedge`（顶层）
- `TestMigrateSeededV1DatabaseToV2`（顶层）
- `TestPrivacyOpsAllDomains`（顶层）
- `TestRetentionBoundaries`（顶层）
- `TestRetentionSchedule`（顶层）
- `TestArtifactsLRUQuota`（顶层）
- `TestArtifactsDeleteAndPurgeAndGuards`（顶层）
- `TestSchemaContractIntrospection`（顶层）
- `TestMigrateFreshDatabaseCreatesCurrentVersion`（顶层）
- `TestMigrationChainWithBackup`（顶层）
- `TestMigrationFailedStepIsAtomic`（顶层）
- `TestMigrationNewerSchemaIsUnmigratable`（顶层）
- `TestMigrationForeignDatabaseIsUnmigratable`（顶层）

**`./internal/panel/`** — 顶层 1 枚 / 子测试 0 枚
- `TestComposerRenderFixtureTellsTheTruth`（顶层）

**`./internal/perm/`** — 顶层 5 枚 / 子测试 0 枚
- `TestTicket90ManualSwitchSurvivesRestart`（顶层）
- `TestTicket90UntouchedConfigStartsAtTheDefault`（顶层）
- `TestTicket90SessionGrantDoesNotSurviveRestart`（顶层）
- `TestTicket90ConfigKeyChangesWhatTheChainAsks`（顶层）
- `TestTicket90HandEditLooseningGoesThroughD36`（顶层）

**`./internal/tools/`** — 顶层 19 枚 / 子测试 2 枚
- `TestRealToolCallWritesRewriteAccountIntoAudit`（顶层）
- `TestDeclaredRiskIsOnlyAFloor`（顶层）
- `TestToolCallRowsAreComplete`（顶层）
- `TestFSReadReturnsTheFileAndTaintsIt`（顶层）
- `TestFSReadTaintFeedsR4`（顶层）
- `TestFSListSummarizesADirectory`（顶层）
- `TestFSListHonoursItsCap`（顶层）
- `TestAtomicWriteKillsMidWrite`（顶层）
- `TestFSTrashGoesToTheRecycleBin`（顶层）
- `TestFSMoveSameVolume`（顶层）
- `TestPathCanonicalizerAccountsForRewrittenRoots`（顶层）
- `TestTicket107AllowlistJudgmentTwoShapes`（顶层）
- `TestTicket107AllowlistBoundaryIsComponentWise`（顶层）
- `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`（顶层）
- `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges`（顶层）
- `TestCanonicalizeReturnsAPathTheOSCanOpen`（顶层）
- `TestLoopPassesDeclaredL1WriteThroughTheGate`（顶层）
- `TestLoopStillRefusesL1WhenNoGateIsRegistered`（顶层）
- `TestL1WriteGoesThroughTheRealBlockWindow`（顶层）
- `TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all`（子测试）
- `TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips`（子测试）

**`./internal/winsec/`** — 顶层 13 枚 / 子测试 4 枚
- `TestAC5FailedSealRefusesTheWrite`（顶层）
- `TestAC5FailureIsNotSwallowedByTheHappyPath`（顶层）
- `TestPOSIXPrivateFileIsReally0600`（顶层）
- `TestPOSIXPrivateDirIsReally0700`（顶层）
- `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed`（顶层）
- `TestPOSIXMissingFileIsNotAnError`（顶层）
- `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC4POSIXFloorAnswersInsideTheNamedTree`（顶层）
- `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`（顶层）
- `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`（顶层）
- `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC5FailedSealRefusesTheWrite/exclusive_artifact`（子测试）
- `TestAC5FailedSealRefusesTheWrite/replacement_write`（子测试）
- `TestAC5FailedSealRefusesTheWrite/directory_chain`（子测试）
- `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct`（子测试）

**`./cmd/wisp/`** — 顶层 25 枚 / 子测试 8 枚
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128`（顶层）
- `TestAC4EveryLegIsNailedOrRuled`（顶层）
- `TestProvidersProbeRecordsMeasuredThinkingFalse`（顶层）
- `TestProvidersProbeRecordsMeasuredThinkingTrue`（顶层）
- `TestProvidersDiscoverListsWhatTheServerServes`（顶层）
- `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`（顶层）
- `TestTicket101ManualSwitchSurvivesRestart`（顶层）
- `TestTicket101UntouchedConfigRestartsAtDefault`（顶层）
- `TestTicket101SessionGrantDoesNotCrossRestart`（顶层）
- `TestTicket101ModeSwitchUsesTheRealL2Gate`（顶层）
- `TestTicket101UnreadableModeFailsLoudlyAndStrict`（顶层）
- `TestRunTextTaskTextPathEndToEnd`（顶层）
- `TestRunTextTaskFailNextIsClassified`（顶层）
- `TestHostDispatchThroughTheAssembledBridge`（顶层）
- `TestComposedGateBlocksAWriteForTwoSeconds`（顶层）
- `TestRunTextTaskKeyResolvesInTheStore`（顶层）
- `TestMissingBlobFailsUnconfiguredNeverSilently`（顶层）
- `TestSecretSetGetListRoundTrip`（顶层）
- `TestSecretFromStdinWritesNoIntermediateFile`（顶层）
- `TestSecretFailurePathsLogAndPrintNoPlaintext`（顶层）
- `TestSecretUnsetRefusesWhileReferenced`（顶层）
- `TestSecretSameNameUnderThreeEnvsIsThreeBlobs`（顶层）
- `TestSecretPortableModeUsesTicket06Seam`（顶层）
- `TestSecretEndToEndConfigRefResolvesAtRequestTime`（顶层）
- `TestSecretOverwriteIsAnnounced`（顶层）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only`（子测试）


## 4. plain 形对照 P1（同一批用例、同一枚快照，只把链接抽掉）

形状：`mkdir -p /plainroot/w124tmp && export TMPDIR=/plainroot/w124tmp`；`/plainroot` 是真目录、容器内不存在 `/varlink`。

**这枚对照的关键读数（先把「plain 形 0 红」这句话量准，别照抄票面）**：plain 形**并非全树 0 红** ——
FAIL 行合计 **29** 枚，全部落在两枚包：`./internal/panel/` 1 枚（`TestComposerRenderFixtureTellsTheTruth`）
与 `./cmd/wisp/` 28 枚。后者错误原文里 `DPAPI is only available on Windows` 有 20 处，其余是
`TestAC4EveryLegIsNailedOrRuled`（命令面腿钉台账那一族）等，逐枚点名见下面 §4 的红名清单 ——
**它们 plain 形同红 ⇒ 与本票形状无关，归因不在本票**（DPAPI 那部分票 119 已记成 `R-119-7`）。
⇒ AC#1 要的「同一批用例在 plain 形下 0 红」**成立于软链造成的那 132 枚**（逐名差集见 §5.3：它们在 P1 里一枚都不红）。
⚠ 别把「plain 形 0 红」当成全树读数抄走：全树 plain 是 29 枚，只是这 29 枚不归本票。

驱动原文开头（挂载自证 + 形状自证，逐字）：

```
== mount self-proof ==
-rwxrwxrwx 1 root root 883 Sep 23 06:26 /src/go.mod
f6ef661732b1851e5c3db348113cb605  /src/go.mod
b6876a5efe759f6e17434d1b50a129c3  /src/internal/winsec/resolve.go
go version go1.27.1 linux/amd64
Linux 26fe69cda429 6.6.114.1-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Mon Dec  1 20:46:23 UTC 2025 x86_64 GNU/Linux
SHAPE=plain TMPDIR=/plainroot/w124tmp resolved=/plainroot/w124tmp
drwxr-xr-x 3 root root 4096 Sep 23 06:59 /plainroot
drwxr-xr-x 2 root root 4096 Sep 23 06:59 /plainroot/w124tmp
WISP_ENV=test HOME=/root
```

| 包 | rc | RUN | PASS | FAIL | SKIP | 子测试 FAIL | 构建失败 |
|---|---|---|---|---|---|---|---|
| `./internal/agent/` | 0 | 75 | 58 | 0 | 0 | 0 | 0 |
| `./internal/agent/approval/` | 0 | 49 | 29 | 0 | 1 | 0 | 0 |
| `./internal/agent/scheduler/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/audio/` | 0 | 18 | 18 | 0 | 0 | 0 | 0 |
| `./internal/ball/` | 0 | 46 | 40 | 0 | 0 | 0 | 0 |
| `./internal/buildinfo/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./internal/config/` | 0 | 101 | 55 | 0 | 0 | 0 | 0 |
| `./internal/llm/` | 0 | 72 | 63 | 0 | 0 | 0 | 0 |
| `./internal/llm/adaptertest/` | 0 | 14 | 4 | 0 | 0 | 0 | 0 |
| `./internal/llm/anthropic/` | 0 | 27 | 9 | 0 | 0 | 0 | 0 |
| `./internal/llm/golden/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/llm/openaichat/` | 0 | 38 | 16 | 0 | 0 | 0 | 0 |
| `./internal/llm/openairesponses/` | 0 | 26 | 8 | 0 | 0 | 0 | 0 |
| `./internal/memory/` | 0 | 67 | 35 | 0 | 1 | 0 | 0 |
| `./internal/models/` | 0 | 51 | 38 | 0 | 2 | 0 | 0 |
| `./internal/observe/` | 0 | 54 | 54 | 0 | 0 | 0 | 0 |
| `./internal/panel/` | 1 | 76 | 45 | 1 | 0 | 0 | 0 |
| `./internal/perm/` | 0 | 14 | 14 | 0 | 0 | 0 | 0 |
| `./internal/plugin/` | 0 | 7 | 7 | 0 | 0 | 0 | 0 |
| `./internal/proc/` | 0 | 16 | 12 | 0 | 0 | 0 | 0 |
| `./internal/risk/` | 0 | 138 | 77 | 0 | 1 | 0 | 0 |
| `./internal/secret/` | 0 | 4 | 4 | 0 | 0 | 0 | 0 |
| `./internal/session/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/speech/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/statemachine/` | 0 | 73 | 10 | 0 | 0 | 0 | 0 |
| `./internal/tools/` | 0 | 79 | 62 | 0 | 3 | 0 | 0 |
| `./internal/watchdog/` | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| `./internal/winsec/` | 0 | 52 | 30 | 0 | 0 | 0 | 0 |
| `./cmd/llmrecord/` | 0 | 2 | 2 | 0 | 0 | 0 | 0 |
| `./cmd/wisp/` | 1 | 70 | 19 | 21 | 0 | 7 | 0 |
| **合计（30 枚包）** | - | **1178** | **718** | **22** | **8** | **7** | **0** |

FAIL 行合计（顶层 + 子测试 = 票 119 数「81 条 FAIL 行」那把尺子的口径）= **29**。

### P1 红名清单（逐枚点名，顶层与子测试分列）

**`./internal/panel/`** — 顶层 1 枚 / 子测试 0 枚
- `TestComposerRenderFixtureTellsTheTruth`（顶层）

**`./cmd/wisp/`** — 顶层 21 枚 / 子测试 7 枚
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128`（顶层）
- `TestAC4EveryLegIsNailedOrRuled`（顶层）
- `TestTicket101ManualSwitchSurvivesRestart`（顶层）
- `TestTicket101UntouchedConfigRestartsAtDefault`（顶层）
- `TestTicket101SessionGrantDoesNotCrossRestart`（顶层）
- `TestTicket101ModeSwitchUsesTheRealL2Gate`（顶层）
- `TestTicket101UnreadableModeFailsLoudlyAndStrict`（顶层）
- `TestRunTextTaskTextPathEndToEnd`（顶层）
- `TestRunTextTaskFailNextIsClassified`（顶层）
- `TestHostDispatchThroughTheAssembledBridge`（顶层）
- `TestComposedGateBlocksAWriteForTwoSeconds`（顶层）
- `TestRunTextTaskKeyResolvesInTheStore`（顶层）
- `TestMissingBlobFailsUnconfiguredNeverSilently`（顶层）
- `TestSecretSetGetListRoundTrip`（顶层）
- `TestSecretFromStdinWritesNoIntermediateFile`（顶层）
- `TestSecretFailurePathsLogAndPrintNoPlaintext`（顶层）
- `TestSecretUnsetRefusesWhileReferenced`（顶层）
- `TestSecretSameNameUnderThreeEnvsIsThreeBlobs`（顶层）
- `TestSecretPortableModeUsesTicket06Seam`（顶层）
- `TestSecretEndToEndConfigRefResolvesAtRequestTime`（顶层）
- `TestSecretOverwriteIsAnnounced`（顶层）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders`（子测试）
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)`（子测试）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only`（子测试）


## 5. 软链形到底造了什么：红因 + 与票面「79 到 83 红」对照

### 5.1 软链形到底造出了什么（红因，逐字取自 L1 日志）

- 容器里 `TMPDIR=/varlink/w124tmp`、`readlink -f` 出 `/realpriv/w124tmp`，而 `t.TempDir()` 返回的是**未解析的** `/varlink/w124tmp/TestXxx<随机>/00N`；用例把这枚拼写原样交给密封底线 ⇒ 底线拒（原文一律是 `winsec: refusing to seal /varlink/... : path is not provably resolved` 那一族）。
- L1 全部日志里出现的**被拒路径**（去重）共 **134** 枚：其中 **130** 枚以 `/varlink/` 开头，而这 130 枚**全部**形如 `/varlink/w124tmp/Test<名><随机数>/...`（= 用例自己 `t.TempDir()` 长出来的根，**0 枚**是产品路向 OS 问来的数据根）；折到 `/varlink/w124tmp/<TestNNN>` 这一层是 **123** 枚临时根、去掉每枚临时根尾部的随机数后是 **119** 枚**用例名**（票 119 那句「79 枚去重被拒路径」是同一把尺，对照见 §5.2）。
- **同口径对照**（只为对上票面那句 79）：只数票 119 那六包的日志，去重被拒路径 = **82** 枚，其中 **81** 枚在 `/varlink/` 之下 ⇒ 票面「79 枚」到今天长了 2 枚；**逐包分布我未按 09-22 的树复算，+2 归给谁留给下一手**（今天六包里各枚包的去重数：winsec 22 / agent 21 / memory 32 / config 6 / proc 0 / secret 0）。
- 剩下 **4** 枚不在此列，全部是用例**自己钉拒绝腿**的形状（逐枚点名，含出处包）：
  - `/realpriv/w124tmp/TestAC3POSIXLinkInsideAResolvedDataRootStillRefused1192984876379/001/data/out/keep-me.txt`  ← 出自 `./internal/winsec/`
  - `/realpriv/w124tmp/TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused1192011455855/001/homereal119/.config/wisp-dev/secrets2`  ← 出自 `./cmd/wisp/`
  - `~/wisp102-tilde-probe/artifacts`  ← 出自 `./internal/risk/`
  - `~\wisp102-tilde-probe\artifacts`  ← 出自 `./internal/risk/`

### 5.2 与票面《事实》三条断言的对照

| 票面断言（2026-09-22 建票时） | 本票在锚点 `7b4c36a` 上的实测 | 判定 |
|---|---|---|
| 「六包 `-count=1 -v` 有 81/83 条 FAIL 行」 | `L1` 形下**那六包**（winsec/proc/config/agent/memory/secret）FAIL 行（顶层+子测试）= **83**：`agent` 26、`config` 7、`memory` 33、`proc` 0、`secret` 0、`winsec` 17 | **对上**（同形同尺同数） |
| 「六包 `-count=1 -v` 有 81/83 条 FAIL 行」 | `L2` 形下**那六包**（winsec/proc/config/agent/memory/secret）FAIL 行（顶层+子测试）= **83**：`agent` 26、`config` 7、`memory` 33、`proc` 0、`secret` 0、`winsec` 17 | **对上**（第二枚样本同数） |
| 「79 枚去重被拒路径全部在 `/varlink/w119tmp/Test*` 之下」 | L1 全 30 包口径：被拒路径去重 **134** 枚，其中 **130** 枚以 `/varlink/` 开头，而这 130 枚**全部**形如 `/varlink/w124tmp/Test*`（**0 枚**是产品路向 OS 问来的数据根）| **机制成立、数目已长**：**同口径（只数那六包）今天 = 81 枚**（票面 79 枚 ⇒ +2；逐包分布见 §5.1，+2 的归因本票不裁）；本票把分母扩到全 30 包后 = 130 枚 |
| 「plain 形 0 红」 | `P1` 全树 FAIL 行 = **29**，但**软链造成的那 132 枚在 plain 形一枚都不红**（逐名差集 §5.3）；plain 剩下的 29 枚 = `./internal/panel/` 1 + `./cmd/wisp/` 28（DPAPI 族与命令面腿钉台账，票 119 `R-119-7` 那本账） | **部分成立**：按票面本意（同一批用例）成立；「plain 形 0 红」**不能当成全树读数抄走** |
| 「红是本形状造成的，不是别的形状」 | L1 的 161 枚 FAIL 行里 **132** 枚只在软链形红、**29** 枚两形都红 ⇒ 红不是"别的形状"带来的 | **成立**（§5.3 逐名差集） |

### 5.3 红名按形状拆开（L1 对 P1 逐名差集）

| 包 | L1 红（顶层+子） | P1 红 | 只在软链形红 | 两形都红（与本票形状无关） |
|---|---|---|---|---|
| `./internal/agent/` | 26 | 0 | 26 | 0 |
| `./internal/agent/approval/` | 1 | 0 | 1 | 0 |
| `./internal/agent/scheduler/` | 0 | 0 | 0 | 0 |
| `./internal/audio/` | 0 | 0 | 0 | 0 |
| `./internal/ball/` | 0 | 0 | 0 | 0 |
| `./internal/buildinfo/` | 0 | 0 | 0 | 0 |
| `./internal/config/` | 7 | 0 | 7 | 0 |
| `./internal/llm/` | 17 | 0 | 17 | 0 |
| `./internal/llm/adaptertest/` | 0 | 0 | 0 | 0 |
| `./internal/llm/anthropic/` | 0 | 0 | 0 | 0 |
| `./internal/llm/golden/` | 0 | 0 | 0 | 0 |
| `./internal/llm/openaichat/` | 0 | 0 | 0 | 0 |
| `./internal/llm/openairesponses/` | 0 | 0 | 0 | 0 |
| `./internal/memory/` | 33 | 0 | 33 | 0 |
| `./internal/models/` | 0 | 0 | 0 | 0 |
| `./internal/observe/` | 0 | 0 | 0 | 0 |
| `./internal/panel/` | 1 | 1 | 0 | 1 |
| `./internal/perm/` | 5 | 0 | 5 | 0 |
| `./internal/plugin/` | 0 | 0 | 0 | 0 |
| `./internal/proc/` | 0 | 0 | 0 | 0 |
| `./internal/risk/` | 0 | 0 | 0 | 0 |
| `./internal/secret/` | 0 | 0 | 0 | 0 |
| `./internal/session/` | 0 | 0 | 0 | 0 |
| `./internal/speech/` | 0 | 0 | 0 | 0 |
| `./internal/statemachine/` | 0 | 0 | 0 | 0 |
| `./internal/tools/` | 21 | 0 | 21 | 0 |
| `./internal/watchdog/` | 0 | 0 | 0 | 0 |
| `./internal/winsec/` | 17 | 0 | 17 | 0 |
| `./cmd/llmrecord/` | 0 | 0 | 0 | 0 |
| `./cmd/wisp/` | 33 | 28 | 5 | 28 |

只在软链形红的合计：顶层 110 + 子测试 22 = **132**；两形都红（既有缺口）：**29**。

### 只在软链形红的逐枚点名（对照 L1 vs P1）

**`./internal/agent/`** 只在软链形红（26 枚）
- `TestFailedTaskBooksOpenCallRowWithDecision`（顶层）
- `TestCancelledTaskPersistsTerminalRows`（顶层）
- `TestPerToolTimeoutFires`（顶层）
- `TestPerToolTimeoutOfContractHonestToolIsToolClass`（顶层）
- `TestTaskLogAndToolCallRows`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide`（顶层）
- `TestWriteFileExclusiveIsExclusive`（顶层）
- `TestSpillSameIDRetryOverwrites`（顶层）
- `TestSpillAcrossRestartsKeepsRetrySemantics`（顶层）
- `TestSpillCallIDHostileShapesSanitizedToBareNames`（顶层）
- `TestSpillContainmentByDirectoryListing`（顶层）
- `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`（顶层）
- `TestSpillThresholdScalesWithWindow`（顶层）
- `TestSpillTokenBoundary`（顶层）
- `TestSpillArtifactAndStubShape`（顶层）
- `TestSpillThroughLoop`（顶层）
- `TestSpillArtifactRespectsRawCap`（顶层）
- `TestMaxTokensFailsAllToolCallsOfThatMessage`（顶层）
- `TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id`（子测试）
- `TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/separator`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/unc`（子测试）
- `TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence`（子测试）

**`./internal/agent/approval/`** 只在软链形红（1 枚）
- `TestTenOpsInOneToolCallGetOneConfirm`（顶层）

**`./internal/config/`** 只在软链形红（7 枚）
- `TestRoundTripLoadMarshalLoad`（顶层）
- `TestResolvedNeverPersists`（顶层）
- `TestSaveFileRoundTripsThroughLoad`（顶层）
- `TestMigrateV1Fixture`（顶层）
- `TestMigrateNoVersionKeyAssumedV1`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone`（顶层）
- `TestUnwiredGuardLeavesHonestConfigsAlone/SaveFile_output_reloads`（子测试）

**`./internal/llm/`** 只在软链形红（17 枚）
- `TestProbeSuiteMeasuresBrokenFC`（顶层）
- `TestProbeSuiteMeasuresBrokenVision`（顶层）
- `TestProbeSuiteHonestProviderRecordsNoMismatch`（顶层）
- `TestProbeSuiteHonestNegativeIsNotAMismatch`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects`（顶层）
- `TestProbeSuiteCapableOnAllThreeDialects`（顶层）
- `TestProbeSuiteRefusesSilentRuns`（顶层）
- `TestProbeSuiteSinkFailurePropagates`（顶层）
- `TestProbeSuiteAudioStaysUnprobed`（顶层）
- `TestProbeSuiteMeasuresBrokenThinking`（顶层）
- `TestProbeSuiteThinkingCapableIsNotTheSameAsBroken`（顶层）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteBrokenOnAllThreeDialects/openai-responses`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-chat`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/anthropic`（子测试）
- `TestProbeSuiteCapableOnAllThreeDialects/openai-responses`（子测试）

**`./internal/memory/`** 只在软链形红（33 枚）
- `TestDeleteArtifactRejectsTheFourHostileShapes`（顶层）
- `TestDeletePrivacyItemRejectsTheFourHostileShapes`（顶层）
- `TestArtifactsContainmentByDirectoryListing`（顶层）
- `TestArtifactsLiteralBackslashKeepsItsAsymmetry`（顶层）
- `TestStraySubdirectoryCannotHideBytesFromTheQuota`（顶层）
- `TestStraySubdirectorySharesTheQuotasLRUQueue`（顶层）
- `TestPurgeArtifactsReclaimsStraySubdirectory`（顶层）
- `TestStrayRemovalDoesNotFollowLinks`（顶层）
- `TestConcurrentWritersReaders`（顶层）
- `TestCrashRecoveryKillMidWrite`（顶层）
- `TestProviderHealthProbeRoundTrip`（顶层）
- `TestProviderHealthRecordErrorAndList`（顶层）
- `TestProviderHealthMissingRowIsNoRows`（顶层）
- `TestProfileUpsertAndLRUEvictionLogs`（顶层）
- `TestProfileDeleteAndPurge`（顶层）
- `TestMemoryAddSearchDeletePurge`（顶层）
- `TestTaskLogLifecycleAndErrorClassValidation`（顶层）
- `TestToolCallLifecycleAndValidation`（顶层）
- `TestGrantCostPluginStateDAO`（顶层）
- `TestWriteQueueLazyLifecycle`（顶层）
- `TestAdvPanicWedge`（顶层）
- `TestMigrateSeededV1DatabaseToV2`（顶层）
- `TestPrivacyOpsAllDomains`（顶层）
- `TestRetentionBoundaries`（顶层）
- `TestRetentionSchedule`（顶层）
- `TestArtifactsLRUQuota`（顶层）
- `TestArtifactsDeleteAndPurgeAndGuards`（顶层）
- `TestSchemaContractIntrospection`（顶层）
- `TestMigrateFreshDatabaseCreatesCurrentVersion`（顶层）
- `TestMigrationChainWithBackup`（顶层）
- `TestMigrationFailedStepIsAtomic`（顶层）
- `TestMigrationNewerSchemaIsUnmigratable`（顶层）
- `TestMigrationForeignDatabaseIsUnmigratable`（顶层）

`./internal/panel/` 两形都红 ⇒ 与本票形状无关，本票登记不修：`TestComposerRenderFixtureTellsTheTruth`

**`./internal/perm/`** 只在软链形红（5 枚）
- `TestTicket90ManualSwitchSurvivesRestart`（顶层）
- `TestTicket90UntouchedConfigStartsAtTheDefault`（顶层）
- `TestTicket90SessionGrantDoesNotSurviveRestart`（顶层）
- `TestTicket90ConfigKeyChangesWhatTheChainAsks`（顶层）
- `TestTicket90HandEditLooseningGoesThroughD36`（顶层）

**`./internal/tools/`** 只在软链形红（21 枚）
- `TestRealToolCallWritesRewriteAccountIntoAudit`（顶层）
- `TestDeclaredRiskIsOnlyAFloor`（顶层）
- `TestToolCallRowsAreComplete`（顶层）
- `TestFSReadReturnsTheFileAndTaintsIt`（顶层）
- `TestFSReadTaintFeedsR4`（顶层）
- `TestFSListSummarizesADirectory`（顶层）
- `TestFSListHonoursItsCap`（顶层）
- `TestAtomicWriteKillsMidWrite`（顶层）
- `TestFSTrashGoesToTheRecycleBin`（顶层）
- `TestFSMoveSameVolume`（顶层）
- `TestPathCanonicalizerAccountsForRewrittenRoots`（顶层）
- `TestTicket107AllowlistJudgmentTwoShapes`（顶层）
- `TestTicket107AllowlistBoundaryIsComponentWise`（顶层）
- `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`（顶层）
- `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges`（顶层）
- `TestCanonicalizeReturnsAPathTheOSCanOpen`（顶层）
- `TestLoopPassesDeclaredL1WriteThroughTheGate`（顶层）
- `TestLoopStillRefusesL1WhenNoGateIsRegistered`（顶层）
- `TestL1WriteGoesThroughTheRealBlockWindow`（顶层）
- `TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all`（子测试）
- `TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips`（子测试）

**`./internal/winsec/`** 只在软链形红（17 枚）
- `TestAC5FailedSealRefusesTheWrite`（顶层）
- `TestAC5FailureIsNotSwallowedByTheHappyPath`（顶层）
- `TestPOSIXPrivateFileIsReally0600`（顶层）
- `TestPOSIXPrivateDirIsReally0700`（顶层）
- `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed`（顶层）
- `TestPOSIXMissingFileIsNotAnError`（顶层）
- `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC4POSIXFloorAnswersInsideTheNamedTree`（顶层）
- `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`（顶层）
- `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`（顶层）
- `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`（顶层）
- `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`（顶层）
- `TestAC5FailedSealRefusesTheWrite/exclusive_artifact`（子测试）
- `TestAC5FailedSealRefusesTheWrite/replacement_write`（子测试）
- `TestAC5FailedSealRefusesTheWrite/directory_chain`（子测试）
- `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct`（子测试）

**`./cmd/wisp/`** 只在软链形红（5 枚）
- `TestProvidersProbeRecordsMeasuredThinkingFalse`（顶层）
- `TestProvidersProbeRecordsMeasuredThinkingTrue`（顶层）
- `TestProvidersDiscoverListsWhatTheServerServes`（顶层）
- `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`（顶层）
- `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`（子测试）

`./cmd/wisp/` 两形都红 ⇒ 与本票形状无关，本票登记不修：`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128`、`TestAC4EveryLegIsNailedOrRuled`、`TestTicket101ManualSwitchSurvivesRestart`、`TestTicket101UntouchedConfigRestartsAtDefault`、`TestTicket101SessionGrantDoesNotCrossRestart`、`TestTicket101ModeSwitchUsesTheRealL2Gate`、`TestTicket101UnreadableModeFailsLoudlyAndStrict`、`TestRunTextTaskTextPathEndToEnd`、`TestRunTextTaskFailNextIsClassified`、`TestHostDispatchThroughTheAssembledBridge`、`TestComposedGateBlocksAWriteForTwoSeconds`、`TestRunTextTaskKeyResolvesInTheStore`、`TestMissingBlobFailsUnconfiguredNeverSilently`、`TestSecretSetGetListRoundTrip`、`TestSecretFromStdinWritesNoIntermediateFile`、`TestSecretFailurePathsLogAndPrintNoPlaintext`、`TestSecretUnsetRefusesWhileReferenced`、`TestSecretSameNameUnderThreeEnvsIsThreeBlobs`、`TestSecretPortableModeUsesTicket06Seam`、`TestSecretEndToEndConfigRefResolvesAtRequestTime`、`TestSecretOverwriteIsAnnounced`、`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask`、`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels`、`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders`、`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor`、`TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(corrupted_blob,_non-portable)`、`TestSecretFailurePathsLogAndPrintNoPlaintext/bad_dpapi_decrypt_(portable,_P13)`、`TestSecretFailurePathsLogAndPrintNoPlaintext/store_write_failure_surfaces_the_ref_only`

### 5.4 软链形的第二桩副作用：它还会把用例改成 SKIP、把子测试截短

| 包 | RUN（L1 软链 / P1 plain） | SKIP（L1 / P1） |
|---|---|---|
| `./internal/config/` | 99 / 101 | 1 / 0 |
| `./internal/memory/` | 36 / 67 | 1 / 1 |
| `./internal/tools/` | 77 / 79 | 3 / 3 |
| `./internal/winsec/` | 45 / 52 | 3 / 0 |

（= 两形 `RUN` 或 `SKIP` 有差的**全部** 4 枚包，其余 26 枚包两形同数；`RUN` 含子测试，故软链形更小。）

- `./internal/winsec/` 在软链形下把 3 枚 POSIX 用例**从「跑」变成「SKIP」**（plain 形 0 枚 SKIP）：
  `TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、
  `TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`。跳过原因逐字（L1 日志，`seam_probe_root_125_other_test.go:223`）：

  ```
  seam_probe_root_125_other_test.go:223: fixture: the harness's own temp base "/varlink/w124tmp/TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125557676579/001" already reaches itself through the link at "/varlink", so neither spelling below is the control side and every reading here would measure one shape
  ```

  ⇒ 这族用例**自己认出了这个形状并拒绝测量**（方向是对的），但后果要记下来：**软链形下分母还会缩小**，
  AC#2 若把调用点解析干净，这 3 枚应回到 PASS 侧，别把它们算进「被清掉的红」。
- `./internal/config/` 同理：`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125` 在 plain 形 `--- PASS`、在软链形 `--- SKIP`（票 125 那枚「POSIX 上 C26 在位」的正向钉自己认出了形状）⇒ **两形加起来这枚才有分母，软链形单独跑等于没跑**。其余 `RUN` 差额（memory 36/67、tools 77/79、winsec 45/52）
  是同一族机制的**截短**：设置阶段就被底线拒 ⇒ 后面的子测试不再展开。
  ⇒ **软链形的 `RUN` 比 plain 形小**，所以「六包 83 枚」这类数**天然低估**被这个形状影响的用例总量：
  光 plain 形下 `memory` 就有 67 枚 RUN、软链形只剩 36 枚。

## 6. 改动面自证（本轮只测不改）

- 本轮入库的只有两枚路径：本证据文件 + 票面 `.scratch/wisp/issues/124-*.md` 的 Progress log 追加。
- 禁改列（`internal/winsec/resolve.go`、`internal/risk/**`、`docs/PLAN.md`、`docs/specs/**`、`rules_gateway.go`、
  `tools/d22scan/**`、`allowlist.txt`、`internal/observe/thresholds.go` 与任何阈值/golden）：本轮 0 hunks。
  `internal/winsec/**` 连 `_test.go` 都没动（AC#1 是「只测不改」）。
- `cmd/wisp/**`、`scripts/**`、`.github/workflows/*`、`frontend/**`：一个字节没碰（归在飞代理）。
- 自查命令：`git show --stat <本轮 commit>` ⇒ 文件清单只含上面两枚路径。

## 7. 未验证项与 `next=`

- **macOS 那半仍未实测**（本机无 macOS runner）：票面《事实》说这形状代表 macOS 的真实形状（`TMPDIR` 在 `/var` 下、
  `/var` 是链接），本票只量了 Linux 容器 ⇒ 与票 119 `R-119-8` 第 4 点同账，**不背书**。
- **单枚慢用例**：`./internal/tools/` 的 `TestL1WriteGoesThroughTheRealBlockWindow` 在软链形下单次 ~300s 才 FAIL，
  把整枚包的墙钟拉长（L1 与 L2 同形）；它在 plain 形的结局见 §4。这条要不要一起清是 AC#2 的事。
- **AC#1 之外的格一个没做**：AC#2（改调用点）、AC#3（不许放宽放行侧）、AC#4（变异自证）、
  AC#5（`-count=2` 门禁 + gofmt/gofumpt/vet/d22scan）、AC#6（`absoluteness` 的 POSIX 分母）全部未动。
- `next=`：**交编排者** —— AC#1 的分母已量完（三枚样本齐、形状差集见 §5.3），请裁 AC#2 的方向：把 §5.3 那些调用点按票 119 生产路那条纪律先解析再递给底线，还是逐枚降级为「就是要测拒绝腿」。裁完我再动那一格。
