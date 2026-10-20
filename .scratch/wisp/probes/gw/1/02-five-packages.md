# gw-1 味②＝票 352 `AC#2` 五包成对色＋逐名 `comm`（★只测不写·零 `.go` 字节）

台面＝两棵**仓外导出树**：基线侧 `gwB`＝锚 `0fb399c3`（＝`ecffb93b^`＝那一行落地的父）⇄ 落地侧 `gwL`＝锚 `a19d5bf2`（＝本味起手的今天 HEAD）。
尺＝本腿 `sh/02-build-tree.sh`（建树＋非空尺）＋`sh/05-run-pair.sh`（跑色／逐名 comm）＋`sh/14-config-remeasure.sh`（撞红复量）＋`sh/16-five-verdicts.sh`（具名逐枚色）。
★量的（跑出来的到达数），⛔ 算的。⛔ 共享工作树跑 `go test`、⛔ 整包 `./cmd/wisp`（票 352 票面逐字归编排者那窗）、⛔ `./...`。
四数＝`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`，★⛔ 混用⛔ 相减。逐名册底片＝仓外 `logs/*.rednames.txt`，原始逐发＝`logs/*.raw.txt`。
树非空尺＝`gwB .go=777` ⇄ `gwL .go=781`（量的，`find -name '*.go' | wc -l`）；`DLL build=3 ⇄ sherpa=3 ⇄ dist 文件 4`；`cmd/wisp .go=139` 两侧同。
★`SHAPE=VOID` 尺＝`RUN=0` 才判 VOID＝**本次 20 发全 `SHAPE=OK`、`VOIDKIND=none`＝台损 0 枚**。

## 一、Windows 本机档（五包 × 两侧；每枚包单独一发）

| 包 | 基线侧 `RC/RUN/PASS_TOP/FAIL_TOP/SKIP` | 落地侧 同形 | 逐名新增红 | 逐名消失红 |
|---|---|---|---|---|
| `internal/secret` | 0/29/21/0/0 | 0/29/21/0/0 | **0** | 0 |
| `internal/config` | 0/195/114/0/0 | 0/195/114/0/0 | **0** | 0 |
| `internal/memory` | 0/80/43/0/2 | 0/80/43/0/2 | **0** | 0 |
| `internal/agent`（★只包自身，⛔ `/**`＝230-r1 写面那族） | 0/114/90/0/0 | 0/114/90/0/0 | **0** | 0 |
| `internal/panel` | 1/268/163/**7**/0 | 1/268/163/**7**/0 | **0** | 0 |

★`internal/panel` 那 7 枚＝两侧**逐名同一批**（尺＝名集 `comm`，⛔ 四数）＝存量形、⛔ 本腿造：
`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestComposerContractTypesMatchFrontend`／
`TestComposerRenderFixtureTellsTheTruth`／`TestReadGitOnThisRepositoryIsSelfConsistent`／
`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`。
★两侧名集 symdiff＝0（`comm -13`⇄`comm -23` 各 0 命中，正控先行＝尺在同射程读得出 7 枚）。
`internal/memory` 那 2 枚 SKIP＝两侧逐名同名（`TestSubprocessCrashWriter`／`TestAppendSessionTurnCrashChild`）＝存量弃权、⛔ 算绿⛔ 算红。

## 二、POSIX 容器档（`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋`MSYS_NO_PATHCONV=1`＋TMPDIR 容器自 `/tmp`）

| 包 | 基线侧 `RC/RUN/PASS_TOP/FAIL_TOP/SKIP` | 落地侧 同形 | 逐名新增红 | 逐名消失红 |
|---|---|---|---|---|
| `internal/secret` | 0/4/4/0/0 | 0/4/4/0/0 | **0** | 0 |
| `internal/config` | 1/196/112/**1**/0 | 1/196/112/**1**/0 | 首发 1（＝换位，见下节复量归 **0**） | 首发 1 |
| `internal/memory` | 0/79/42/0/2 | 0/79/42/0/2 | **0** | 0 |
| `internal/agent` | 0/113/89/0/0 | 0/113/89/0/0 | **0** | 0 |
| `internal/panel` | 1/268/163/**7**/0 | 1/268/163/**7**/0 | **0** | 0 |

## 三、撞红那一处的三发复量＋并集归因（★公共块令＝⛔ 拿单发换位当新增红、⛔ 把摆动当闭环）

首发＝基线侧红 `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` ⇄ 落地侧红 `TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`。

| 发 | 基线侧红名 | 落地侧红名 | 单发新增红 |
|---|---|---|---|
| 首发 | `…HandDeletedSchema…` | `…NothingWritten…` | 1（换位） |
| r1 | `…Unnailed…`＝`HandDeletedSchemaVersion…` 同枚（逐字见下） | `NothingWritten…` | 1（同形换位） |
| r2 | **两枚都红**（`FAIL_TOP=2`） | `NothingWritten…` | **0** |
| r3 | `HandDeletedSchemaVersion…` | **两枚都红**（`FAIL_TOP=2`） | 1（`NothingWritten…`） |
| ★三发并集 | 两枚 | 两枚（`cmp` ＝**同一批名**） | **新增红 0** |

逐字红句（两侧同族两枚，尺＝`-v` 原始件里的 `file:line:` 行）＝
`writeguard_unnailed_341_test.go:239: mtime did not move on the write branch` ⇄
`writeguard_handedit_227_test.go:136: mtime did not move on a stale-schema file that logged 'added nothing' - the migration-on-read boundary is no longer being exercised`。
⇒ 两枚＝公共块摆动名册在册的 **`internal/config` mtime 那枚**族（容器 fs mtime 粒度），⇄ 三发里两侧**都各自红过且并集相同**＝⛔ 字节所造；★方向也不能反读成"落地更干净"＝单发 r2/r3 恰好相反。⇒ 本味判据（新增红）＝**0 枚**。

## 四、票面具名连带（逐枚点名两侧色；★尺＝`sh/16` 锚定 `^(    )?--- (PASS|FAIL|SKIP): <名> `，⛔ 拿四数顶逐枚）

`private_fail_test.go:37/50/61/72/75` 五枚 `ErrNotSealable` 判定 ⇄ 宿主用例⇄子测试映射（行号尺＝对**入库 blob**`git cat-file blob`取那一行逐字＝五枚全部命中 `errors.Is(err, ErrNotSealable)`）：

| 判定行 | 逐字 | 宿主用例／子测试 | 基线侧（host⇄container） | 落地侧（host⇄container） |
|---|---|---|---|---|
| `:37` | `if !errors.Is(err, ErrNotSealable) {` | `TestAC5FailedSealRefusesTheWrite/exclusive_artifact` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| `:50` | `if !errors.Is(err, ErrNotSealable) {` | ↳ `/replacement_write` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| `:61` | `} else if !errors.Is(err, ErrNotSealable) {` | ↳ `/directory_chain` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| `:72` | `if err := SealFile(p); !errors.Is(err, ErrNotSealable) {` | ↳ `/seal_file_and_dir_direct` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| `:75` | `if err := SealDir(dir); !errors.Is(err, ErrNotSealable) {` | ↳ 同枚子测试（★具名一处＝票面把 :72⇄:75 算两枚判定、仪器只有四枚子测试名＝两枚同住一支；本腿⛔ 圆判） | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| 宿主顶层 | — | `TestAC5FailedSealRefusesTheWrite` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| 同文件另枚 | — | `TestAC5FailureIsNotSwallowedByTheHappyPath` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |
| 覆盖写承诺 | — | `TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1` | **PASS ⇄ PASS** | **PASS ⇄ PASS** |

读 `ErrUnresolvedPath` 的两族 ⇄ 两侧色（★`ABS`＝该 build tag 在这一档⛔ 编译＝弃权形，⛔ 算红⛔ 算绿）：

| 具名 | host 基线 ⇄ 落地 | container 基线 ⇄ 落地 |
|---|---|---|
| `openleaf_322`＝`TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1` | ABS ⇄ ABS（⛔ 编译） | **PASS ⇄ PASS** |
| `placement_symlink_113` 八枚逐名（`…SealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`／`…RefusesALinkAncestorAtEveryDepth`／`…SealDirThroughASymlinkRefuses`／`…PrivateFileThroughASymlinkRefusesAndWritesNothing`／`…SealFileThroughABackslashNamedLink`／`…StillNarrowsAPlainFileInsideTheNamedTree`／`…StillWorksNextToAndThroughRealDirectoriesAndLinks`／`…DoesNotFoldABackslashIntoASeparator`） | ABS ⇄ ABS | **逐枚 PASS ⇄ 逐枚 PASS** |
| `openleaf_guard_322` 两枚（`TestAC3WindowsLeafGuardRefusesLinkNames…`／`…PrivateFileRouteHasTheLeafGuardSpeaking…`） | **PASS ⇄ PASS** | ABS ⇄ ABS |

⇒ ★两档合起来每枚具名都在**它该编译的那一档**量到 PASS＝⛔ 受影响；★具名件两侧字节同一性（尺＝`git rev-parse <锚>:<path>`，量的）＝`private_fail_test.go`＝`e95d48aa` 三锚同串、`openleaf_overwrite_322_test.go`＝`39026065` 同、`openleaf_322_other_test.go`＝`dfaddac8` 同、`placement_symlink_113_other_test.go`＝`3a18a3d4` 同、`openleaf_guard_322_windows_test.go`＝`b2b71590` 同 ⇒ 逐名作差只反映产码。

## 五、味② 锚⇄时点自证（★那味＝读数时点⇄字节时点对得上）

起手锚＝`git rev-parse HEAD`＝**`a19d5bf2`**（22:49 放第一发 host 成对色前现量）⇄ 交件侧 HEAD＝**见 `04` 件**（其间漂过的笔全是编排者⇄别腿的**非 `.go`** 件）。
★漂没漂尺（量的）＝`git diff --name-only a19d5bf2..HEAD -- "*.go"`＝**0 枚**；另加一把更硬＝七枚被测／相关目录的 **tree hash** 两侧同串（`internal/winsec 68e765de…`／`internal/secret 2d2c53b5…`／`internal/config 3eeed0b3…`／`internal/memory 5a5cf82c…`／`internal/agent ef30090f…`／`internal/panel c8b5e209…`／`cmd/wisp 34af5e0c…`＝`SAME=YES`）⇒ ★⛔ 任一发需要重跑：读数所依字节与我锚定那次 HEAD 逐目录同串。
跑色前后 `winsec.go` 物理哈希对回＝两侧树各自未被写动（`gwB`＝`07300bff…` blob ⇄ `gwL`＝`757b16a8…` blob，尺＝`git rev-parse <锚>:internal/winsec/winsec.go`）。

## 六、味② 台面登记

`0xc0000135`＝**0**／`build failed`＝**0**／`[build failed]`⇄`setup failed`＝**0**／`walk`＝**0**／`no test files`＝**0**／`SHAPE=VOID`＝**0**（20 发全 OK）；台根＝`%LOCALAPPDATA%\wisp-test-gotmp\gw1-<时间戳>`（★⛔ `D:\tmp\wisp` 作 `GOTMPDIR`＝109-r1 量到的 icacls 继承假红那一味；`icacls` 现量＝该根只有三条 `(I)(OI)(CI)(F)`、**无** `BUILTIN\Users` 继承）；★跑 panel 相关带 `PATH=<tree>/third_party/sherpa-onnx;<tree>/build`（`cygpath -w` 两枚段＋分号，⛔ `cygpath -wp` 转整条 PATH）。
