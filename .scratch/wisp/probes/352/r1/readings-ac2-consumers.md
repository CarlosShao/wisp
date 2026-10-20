# 352-r1 · AC#2 消费者面（五包＋`internal/winsec`）成对逐名作差＝新增红 0 枚

★⛔ `cmd/wisp` 整包（本波明令＝归编排者那窗，且会撞 `351-r2` 写面）；★⛔ `./internal/agent/**` 子包（＝`230-r1` 写面）
⇒ 本腿只打 `./internal/agent/` 这一枚包自身，且⛔ 共享工作树、只打在仓外导出树上。
台面＝两棵自建树：基线 `treeB`＝`git archive bf509bb8`（＝本腿起手现认锚，`winsec.go` blob `07300bff…`）
⇄ 落地 `treeL`＝`git archive ecffb93b`（那一行那笔，blob `f3c1b68e…`）；两棵 `.go` 枚数都＝777（非空尺）；
每发跑后 `winsec.go` sha256 前后对回＝`RESTORE=OK`（★⛔ `git checkout`＝还原靠仓外 pristine 拷贝）。
尺＝`sh/03-consumer-pair.sh`（跑色）＋`sh/05-perpkg-diff.sh`（逐名作差）。
★四数只在 `^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP` 取、⛔ 混用⛔ 相减；★新增红只认**逐名 comm**，⛔ 拿整包四数差当"我的字节改了 N 枚测试"。

## 表一 · Windows 本机档（`go1.27.1 windows/amd64`，GOTMPDIR 反斜杠形落 `%LOCALAPPDATA%` 同卷短根）

逐字＝`logs/ac2-host-pair-s1-treeB.raw.txt` ⇄ `logs/ac2-host-pair-s2-treeL.raw.txt`；行区间切片归属（`^(ok|FAIL)[[:space:]]+github\.` 汇总行为右边界）。

| 包 | 基线侧 RUN/PASS_TOP/FAIL_TOP/SKIP | 落地侧 RUN/PASS_TOP/FAIL_TOP/SKIP | 包色 | 新增红（只在 L） | 消失红（只在 B） |
|---|---|---|---|---|---|
| `internal/secret` | 29/21/0/0 | 29/21/0/0 | ok ⇄ ok | **0** | 0 |
| `internal/config` | 195/114/0/0 | 195/114/0/0 | ok ⇄ ok | **0** | 0 |
| `internal/memory` | 80/43/0/2 | 80/43/0/2 | ok ⇄ ok | **0** | 0 |
| `internal/agent`（★只这一枚包自身，⛔ `/**`） | 114/90/0/0 | 114/90/0/0 | ok ⇄ ok | **0** | 0 |
| `internal/panel` | 268/163/7/0 | 268/163/7/0 | FAIL ⇄ FAIL（同一批 7 枚） | **0** | 0 |
| `internal/winsec` | 119/66/2/0 | 119/66/2/0 | FAIL ⇄ FAIL | 1（＝`342r1` 竞态族，见下节＝复量 3 发后新增红 **0**） | 1（同族另一枚） |

`internal/panel` 那 7 枚红＝两侧**逐名同一批**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／
`TestComposerContractTypesMatchFrontend`／`TestComposerRenderFixtureTellsTheTruth`／`TestReadGitOnThisRepositoryIsSelfConsistent`／
`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`）＝存量形（前端契约/长度账本那族），⛔ 我新造。

## 表二 · POSIX 容器档（`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋TMPDIR 容器自 `/tmp`）

逐字＝`logs/ac2-col-s1-treeB.raw.txt` ⇄ `logs/ac2-col-s2-treeL.raw.txt`。

| 包 | 基线侧 | 落地侧 | 包色 | 新增红 | 消失红 |
|---|---|---|---|---|---|
| `internal/secret` | 4/4/0/0 | 4/4/0/0 | ok ⇄ ok | **0** | 0 |
| `internal/config` | 196/111/2/0 | 196/113/0/0 | FAIL ⇄ ok | **0** | 2（＝`mtime` 那族摆动件，见下节） |
| `internal/memory` | 79/42/0/2 | 79/42/0/2 | ok ⇄ ok | **0** | 0 |
| `internal/agent` | 113/89/0/0 | 113/89/0/0 | ok ⇄ ok | **0** | 0 |
| `internal/panel` | 268/163/7/0 | 268/163/7/0 | FAIL ⇄ FAIL（同一批 7 枚） | **0** | 0 |
| `internal/winsec` | 74/46/0/0 | 74/46/0/0 | ok ⇄ ok | **0** | 0 |

★容器档 `internal/winsec` 两棵都全绿（74/46/0/0）＝本票落地⛔ 造红＝POSIX 半边是最干净的一味。

## 撞红的两处成对复量（★公共块令＝安静复量 3 发再归因，⛔ 记成我新造、⛔ 把摆动当闭环）

### ① `internal/winsec` 本机档那一枚逐名换位（尺＝`logs/ac2r{1,2,3}-winsec-s{1,2}-*.rednames.txt`）

| 发 | 基线侧红名 | 落地侧红名 | 新增红（L-only） |
|---|---|---|---|
| 首发（表一那对） | `…NeverTruncates…342r1` ＋ `…LeafIdentityReCheck…342r1` | `…NeverDeletes…342r1` ＋ `…LeafIdentityReCheck…342r1` | 1 |
| 复量 1 | `…LeafIdentityReCheck…` ＋ `…NeverDeletes…` | `…LeafIdentityReCheck…` | **0** |
| 复量 2 | `…NeverTruncates…` ＋ `…LeafIdentityReCheck…` | `…NeverTruncates…` | **0** |
| 复量 3 | 三枚全红 | `…NeverTruncates…` ＋ `…LeafIdentityReCheck…` | **0** |
| ★三发并集作差 | 三枚 | 两枚（⛔ 一枚是并集里没有的） | **新增红 0** |

⇒ 那三枚名**全部住在 `internal/winsec/ancestor{open,unlink}_34{2,8}_windows_test.go`＝票 342/348 的 `342r1` 竞态钉**＝派单摆动名册在册那族（`342-v2`/`342-r2` 都具名量过"落地侧仍有 0~2/400 漏出、基线侧每发都量到"＝⛔ 恒绿、有牙）。
★复量三发里两侧都出现过同一批名＝两侧并集相同 ⇒ **本票那一枚字节新增红 0 枚**；首发那枚 1/1 换位＝竞态形状随台面漂移（公共块：竞态读数⛔ 要求逐字相同，复现的是量级与形状）。
★另证：`AC#1` 那把尺（R1..R13 十三枚读值）在两棵树上⛔ 出现任何 `--- FAIL`＝`FAIL_TOP=0` 两侧同。

### ② `internal/config` 容器档那两枚消失红（尺＝`logs/ac2r{1,2,3}-col-config-s{1,2}-*.rednames.txt`）

| 发 | 基线侧红名 | 落地侧红名 | 新增红（L-only） |
|---|---|---|---|
| 首发 | 两枚都红 | 0 | 0（＝⛔ 我造红，方向相反） |
| 复量 1 | 两枚都红 | 0 | **0** |
| 复量 2 | `…NothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent` | 两枚都红 | 1 |
| 复量 3 | `…NothingWritten…` | 两枚都红 | 1 |
| ★两侧并集 | 两枚 | 两枚 | **新增红 0**（★且落地侧并集＝基线侧并集＝同一批名） |

逐字红句＝`writeguard_handedit_227_test.go:136: mtime did not move on a stale-schema file that logged 'added nothing' - the migration-on-read boundary is no longer being exercised`
＋`writeguard_unnailed_341_test.go:239: mtime did not move on the write branch`。
⇒ 两枚都＝派单摆动名册里点名的 **`internal/config` mtime 那枚**族＝容器文件系统 mtime 粒度（同一秒内落盘⇒mtime⛔ 动）＝与 `wrapPath` 那枚格式串⛔ 关：
★两侧都各自红过、两侧并集相同＝⛔ 我的字节所造；★方向还是"落地侧更干净"，⛔ 据此判绿当凭据（只当"⛔ 我造红"的反证）。

## 票面具名连带（逐枚，尺＝`logs/perpkg-diff.txt` 末节的纯字符串尺）

| 具名 | 基线侧 | 落地侧 | 结论 |
|---|---|---|---|
| `TestAC5FailedSealRefusesTheWrite`（＝`private_fail_test.go:37/50/61/72/75` 五枚 `ErrNotSealable` 判定的宿主钉） | PASS | **PASS** | ★仍绿 |
| ↳ 子测试 `exclusive_artifact` ⇄ `replacement_write` ⇄ `directory_chain` ⇄ `seal_file_and_dir_direct` | 4 枚全 PASS | **4 枚全 PASS** | ★逐枚仍绿（两档都如此＝`logs/perpkg-diff.txt` 两侧各 4 行） |
| `TestAC5FailureIsNotSwallowedByTheHappyPath`（同文件那枚 `fs.ErrExist` 判定） | PASS | **PASS** | ★仍绿 |
| `TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1`（四枚生产调用者＋覆盖写承诺） | PASS（6 枚子测试逐枚 PASS） | **PASS（6 枚逐枚 PASS）** | ★⛔ 受影响 |
| 读 `ErrUnresolvedPath` 的 `openleaf_322` 族＝`TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1` | 容器 PASS | 容器 **PASS** | ★⛔ 受影响（本机档＝该件 `!windows`，两侧都 ABSENT＝⛔ 编译，弃权形照实登记） |
| 读 `ErrUnresolvedPath` 的 `placement_symlink_113` 族＝8 枚逐枚（`…SealFileThroughASymlinkRefuses…`／`…RefusesALinkAncestorAtEveryDepth`／`…SealDirThroughASymlinkRefuses`／`…PrivateFileThroughASymlinkRefusesAndWritesNothing`／`…ThroughABackslashNamedLink`／`…StillNarrowsAPlainFileInsideTheNamedTree`／`…DoesNotFoldABackslashIntoASeparator`） | 容器逐枚 PASS | 容器逐枚 **PASS** | ★⛔ 受影响 |
| `openleaf_guard_322`（`…LeafGuardRefusesLinkNames…`＋`…RouteHasTheLeafGuardSpeaking…`） | 本机 PASS | 本机 **PASS** | ★⛔ 受影响（容器档＝`windows` 件，两侧 ABSENT） |

★弃权两族照实（公共块在册＝`internal/winsec` 两枚 `//go:build !windows` 仪器＝本机/CI windows 档⛔ 编译＝既⛔ 绿也⛔ 弃权）＝
`ABSENT` 那几枚⛔ 算红、⛔ 算绿、⛔ 当闭环；两档合起来每枚具名件都在**它该编译的那一档**量到 PASS。

## AC#2 结论

**落地⇄基线成对逐名作差＝新增红 0 枚**（五包两档 12 格逐格 0；`internal/winsec` 那枚首发换位经 3 发复量归因＝两侧并集相同＝⛔ 我所造）。
★⛔ 放宽任何断言、⛔ 删用例、⛔ `t.Skip`、⛔ 报"全绿"（`internal/panel` 那 7 枚存量红两侧都在册、照实逐名登记）。
