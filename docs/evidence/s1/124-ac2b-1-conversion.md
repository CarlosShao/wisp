# 票 124 AC#2b 批次1（`worker-ticket124-ac2b-1`）— memory 33 + config 7 接上票 119 那条纪律

**日期**：2026-09-23 · **性质**：转换方交件（改 `.go`，但只改 `_test.go`）
**本批范围**：清点账 `docs/evidence/s1/124-ac2a-leg-classification.md` §5.4 的第 45-51 行（config 7 枚）与 §5.7 的第 74-106 行（memory 33 枚），共 **40 枚**，账上全标「可转」。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2b 的放行与分批」（09-23 16:33 编排者）那张表 + 五条共用结案判据。

## 0. 锚点 / 快照 / 跑法（可复核）

- 开工首读 `git rev-parse HEAD` = `158f7298662beeed491bab41a837d65c5d28629f`（下文记 **pre 锚**）。
  ⚠ 共树期间 HEAD 会漂：本次测量途中先后漂到 `658fb19`、`b22f898`。已核 `git diff --name-only 158f729 HEAD` 只命中 `docs/evidence/s1/134-adversarial-acceptance.md` 一枚，且 `git diff --stat 158f729 HEAD -- internal/ cmd/ go.mod go.sum` **零 hunk** ⇒ pre 读数量的就是被改前的同一版本。
- 改件 commit：**`fafe2b4a6299202df477fde0e358028011ff37d3`**（`test(124,AC#2b-1)`，9 枚文件，+96/-14）。
- 全部测量都在 `git archive <sha> | tar -x` 的纯净快照里做；工作树里别人的半成品一个字都不读。

| 快照目录 | 内容 | 文件数 |
|---|---|---|
| `/d/tmp/wisp124-2b1-pre` | `git archive 158f729`（改前基线） | 1035 |
| `/d/tmp/wisp124-2b1-post` | `git archive fafe2b4`（改后） | 1038 |
| `/d/tmp/wisp124-2b1-probe2` | 改后 + 判据④的 `t.Logf` 取字串探针（19 处，**未入库**） | 1039 |
| `/d/tmp/wisp124-2b1-probe3` | 同 probe2 的修正版（两枚 `_` 弃值处把原 `t.Errorf` 一起搬了过去，编译过） | 1039 |
| `/d/tmp/wisp124-2b1-mut` | 改后 + 判据⑤变异 `MUTATION-124-2B1`（拆掉新加的解析层，**未入库**） | 1038 |

- 容器 `golang:1.27`（`go1.27.1 linux/amd64`，容器内 `CGO_ENABLED=1`），复用票 119 的命名卷 `ac119-gomodcache` / `ac119-gocache` ⇒ 离线可编。
- 挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`；**每枚样本进容器第一件事**打 `ls -l /src/go.mod`（`-rwxrwxrwx 1 root root 883`）+ `md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum resolve.go` = `b6876a5efe759f6e17434d1b50a129c3` — 三枚与 AC#1/AC#2a 逐字同字 ⇒ 非空挂自证（Git Bash 下 `docker run -v "C:\…"` 会静默挂空且 rc=0＝假绿）。
- 形状硬断言（防「`/varlink` 被 mkdir 成真目录、链接根本没建成」那枚假绿）：软链形 `exit 97`（`/varlink` 不是 symlink）/ `exit 98`（`readlink -f /varlink/w124tmp != /realpriv/w124tmp`）；普通形**另加** `exit 99`（普通形里 `/varlink` 必须根本不存在、`/plainroot` 不许是链接）。本批 8 发容器全部通过形状断言、无一发 97/98/99。
- 退出码一律 `rc=$?` 单独一行取（dash 下 `${PIPESTATUS[0]}` 是 Bad substitution）。
- 跑法脚本（新建，**未改** AC#2a 那枚 `/d/tmp/wisp124-ac2a-h.sh`，只作形状参考）：`/d/tmp/wisp124-2b1-h.sh`（`bash /h.sh <link|plain> <tag>`，`COUNT=2` 切门禁那一发）、`/d/tmp/wisp124-2b1-m.sh`（变异台：先 grep 证落地 + `go build`/`go vet` rc，再读红名）、`/d/tmp/wisp124-2b1-probe-patch.py`（判据④探针）、`/d/tmp/wisp124-2b1-mutate.py`（判据⑤变异）。日志全在 `/d/tmp/wisp124-2b1-logs/`。
- 开测前查在飞 run（本机就是 self-hosted runner）：16:36 `gh run list --limit 5` 见 2 枚 `in_progress`（`35837937954`、`35837787364`，都是 docs push 的 `ci`）⇒ 与我的跑同机抢 CPU；收尾 17:2x 复看已全部 `completed`。争用只影响单枚耗时，不影响红绿与四数（两形同码只差 `TMPDIR` 一个变量；`-count=2` 那一发的四数与 `-count=1` 逐数成 2 倍关系，也侧面排掉争用改判）。

## 1. 对派单数字的复核（派单给的每个数字都是断言）

| 派单断言 | 实测 | 结论 |
|---|---|---|
| 本批 40 枚 = memory 33 + config 7 | pre 锚软链形逐名逮到 memory `--- FAIL` **33** 枚、config 顶层 `--- FAIL` **6** + 子测试 `--- FAIL` **1** = **7** 枚；40 枚与账 §5.4/§5.7 的行号逐名对得上 | **成立** |
| 这 40 枚全是「可转」（拒绝腿档 = 0 枚） | 40 枚在 pre 锚普通形**逐枚 = PASS**（§2 表 PRE-P 列无一 FAIL/SKIP）⇒ 无一能自证是「未解析根必须被拒」的那条腿 | **成立** |
| `TestCrashRecoveryKillMidWrite` 把根经 `WISP_CRASH_DIR` 交给子进程 ⇒ 交出去那一份也要解析 | pre 锚软链形 `--- FAIL (30.01s)`、`concurrent_test.go:254: subprocess only committed 0 rows within budget (want >= 20)`；改后软链 `--- PASS (0.06s)` 且子进程真开了库（`concurrent_test.go:312: crash recovery: 35 committed rows recovered, integrity ok, wal=0 bytes`） | **成立，且已落地** |
| `internal/config` 那 1 枚 `TestAC2POSIX…125` 系 seam 自拒探针不算进清掉的红 | 该枚与其两枚子测试在改前改后的软链形都 `SKIP`／不展开，**从未出现在 40 枚名册里**；本批也没动它（`c26_seam_posix_125_test.go` 不在这枚 commit 的文件清单内） | **成立** |

## 2. 判据①：逐枚转绿且点名（`-v` 才有 PASS 名）

四台对照（同一批用例，`-count=1 -v`）。PRE-L=改前软链、PRE-P=改前普通、POST-L=改后软链、POST-P=改后普通。

| 包 | 用例名 | PRE-L(改前软链) | PRE-P(改前普通) | POST-L(改后软链) | POST-P(改后普通) |
|---|---|---|---|---|---|
| `memory` | `TestAdvPanicWedge` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestArtifactsContainmentByDirectoryListing` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestArtifactsDeleteAndPurgeAndGuards` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestArtifactsLRUQuota` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestArtifactsLiteralBackslashKeepsItsAsymmetry` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestConcurrentWritersReaders` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestCrashRecoveryKillMidWrite` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestDeleteArtifactRejectsTheFourHostileShapes` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestDeletePrivacyItemRejectsTheFourHostileShapes` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestGrantCostPluginStateDAO` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMemoryAddSearchDeletePurge` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrateFreshDatabaseCreatesCurrentVersion` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrateSeededV1DatabaseToV2` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrationChainWithBackup` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrationFailedStepIsAtomic` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrationForeignDatabaseIsUnmigratable` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestMigrationNewerSchemaIsUnmigratable` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestPrivacyOpsAllDomains` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestProfileDeleteAndPurge` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestProfileUpsertAndLRUEvictionLogs` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestProviderHealthMissingRowIsNoRows` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestProviderHealthProbeRoundTrip` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestProviderHealthRecordErrorAndList` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestPurgeArtifactsReclaimsStraySubdirectory` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestRetentionBoundaries` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestRetentionSchedule` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestSchemaContractIntrospection` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestStrayRemovalDoesNotFollowLinks` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestStraySubdirectoryCannotHideBytesFromTheQuota` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestStraySubdirectorySharesTheQuotasLRUQueue` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestTaskLogLifecycleAndErrorClassValidation` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestToolCallLifecycleAndValidation` | FAIL | PASS | **PASS** | PASS |
| `memory` | `TestWriteQueueLazyLifecycle` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestMigrateNoVersionKeyAssumedV1` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestMigrateV1Fixture` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestResolvedNeverPersists` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestRoundTripLoadMarshalLoad` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestSaveFileRoundTripsThroughLoad` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestUnwiredGuardLeavesHonestConfigsAlone` | FAIL | PASS | **PASS** | PASS |
| `config` | `TestUnwiredGuardLeavesHonestConfigsAlone/SaveFile_output_reloads` | FAIL | PASS | **PASS** | PASS |

**逐包四数**：

| 台 | 形状 | `internal/memory` | `internal/config` |
|---|---|---|---|
| PRE-L（改前） | 软链 | `rc=1 RUN=36 PASS=2 FAIL=33 SKIP=1 SUBPASS=0 SUBFAIL=0` | `rc=1 RUN=99 PASS=48 FAIL=6 SKIP=1 SUBPASS=43 SUBFAIL=1` |
| PRE-P（改前） | 普通 | `rc=0 RUN=67 PASS=35 FAIL=0 SKIP=1 SUBPASS=31 SUBFAIL=0` | `rc=0 RUN=101 PASS=55 FAIL=0 SKIP=0 SUBPASS=46 SUBFAIL=0` |
| POST-L（改后） | 软链 | `rc=0 RUN=67 PASS=35 FAIL=0 SKIP=1 SUBPASS=31 SUBFAIL=0` | `rc=0 RUN=99 PASS=54 FAIL=0 SKIP=1 SUBPASS=44 SUBFAIL=0` |
| POST-P（改后） | 普通 | `rc=0 RUN=67 PASS=35 FAIL=0 SKIP=1 SUBPASS=31 SUBFAIL=0` | `rc=0 RUN=101 PASS=55 FAIL=0 SKIP=0 SUBPASS=46 SUBFAIL=0` |

⇒ 本批红名数 **40 → 0**（软链形），两包 `rc` 从 `1` 到 `0`。日志里 `refusing to seal` / `not provably resolved` 的行数：改前软链 memory 32 + config 6 = **38**，改后软链 **0 + 0 = 0**。

## 3. 判据②：普通形一枚都不许多红 + `git diff` 证判定分支一字未动

**逐数相同**：POST-P 与 PRE-P 两包**八个数逐数相同**（memory `67/35/0/1`+`31`、config `101/55/0/0`+`46`）。普通形红名数 **0 → 0**，一枚都不许多红。

`git show fafe2b4` 的**删除侧全文只有 14 行**，逐行抄（`-` 号省略）：

- `dir := t.TempDir()` — config/boundary_test.go 2 处、config/loader_test.go 1 处、config/unwired_test.go 1 处（子测试内，缩进一级）
- `root := t.TempDir()` — memory/artifacts_path_invariant_test.go
- `dir := filepath.Join(t.TempDir(), "data")` — memory/concurrent_test.go、memory/migrate_v1seed_test.go、memory/schema_test.go 5 处
- `// openTestStore opens a Store in a fresh t.TempDir subdirectory.` — 注释一行

**新增侧**（除两枚新文件外只有 23 行）：13 行是 `dir/root := sealableTempDir124(t)`（或其 `filepath.Join(…, "data")` 形），10 行是注释。**判定分支、断言、阈值、golden 一行未动** — 删除行与新增行一一对应，全部落在「递给底线的根怎么拼」这一件事上。⇒ 票面 AC#3「不许拿放行侧放宽换绿」在本批成立：解析层在调用方（测试）这一侧，底线的判定码一个字没碰。

范围核对：`git show --numstat fafe2b4` = 9 枚文件、全部 `_test.go`；`git diff --name-only 158f729 fafe2b4` 里 `internal/`+`cmd/` 非 `_test.go` 命中 **0 枚**。禁改列 `internal/winsec/resolve.go`、`internal/risk/**` **零 hunk**；`cmd/wisp/**` **零 hunk**（2b-4 地界）。

## 4. 判据③：两形各取一枚 + RUN/SKIP 差逐包解释

两形各一枚已完成（§2 的 POST-L / POST-P，另 `-count=2` 各一枚见 §7）。逐包解释：

- **`internal/memory`**：改前软链 `RUN=36` vs 普通 `RUN=67`（−31）；**改后软链 `RUN=67` = 普通形**，缺口自己合上。那 31 枚逐名列在 `/d/tmp/wisp124-2b1-logs/restored-memory-subtests.txt`：`TestDeleteArtifactRejectsTheFourHostileShapes/{separator,separator_backslash_literal,dotdot,dotdot/bare_and_empty,dotdot_backslash_literal,drive_letter,unc}` 7 枚、`TestDeletePrivacyItemRejectsTheFourHostileShapes/…` 6 枚、`TestSchemaContractIntrospection/{columns,indexes}×9 表` 18 枚。它们改前**根本没被创建**（父用例在 `Open` 那一步 `t.Fatalf` ⇒ 子测试不展开），不是 SKIP；改后 **31 枚逐枚 = PASS**（脚本核全部，非抽样）。
  SKIP 名册改前改后**逐名同一枚**：`TestSubprocessCrashWriter`（`concurrent_test.go:186`，子进程角色在父进程里自己 `t.Skip`）⇒ 与解析层无关，也没因解析层增减。
- **`internal/config`**：软链 `RUN=99` vs 普通 `RUN=101`，**改前改后这个差一字不变**（99/101）。差的正是票 125 那枚 seam 自拒探针 `TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125` 与其两枚子测试 `control_plain_temp`、`measured_symlink_spelled_temp`：探针在 fixture 里看到「harness 的 temp base 穿过链接回到自己」就 `t.Skipf`，子进程不启动 ⇒ 两枚子测试不计入 RUN。名册级差集用 `comm` 双向核过：普通形有、软链形没有的只有这 2 枚；软链形有、普通形没有的 **0 枚**。
  SKIP 名册改前改后**逐名同一枚**（就是那枚 125 探针），它**不在**本批 40 枚里，也不计为清掉的红。
- **怎么排除「被跳过冒充变绿」（逐枚）**：40 枚的状态取自 `-v` 日志里的**逐名行**，不是包级 `ok`。判据是 `POST-L.<pkg>.res.txt` 每枚名字旁边写的是 `PASS` 还是 `SKIP` — 40 枚**逐枚 = `PASS`**（§2 表 POST-L 列，无一 SKIP/FAIL）；再核改后软链的**用例名集合**与改前普通形**逐名相等**（`comm` 双向为空 ⇒ 既没多出一枚新用例名，也没少一枚）。

## 5. 判据④：本批「另一种拒」复算 — 逐枚实拿字符串

**取字串仪器**：`/d/tmp/wisp124-2b1-probe3` = 改后快照 + 19 处 `t.Logf`（只打印、不改断言；两处 `_` 弃值的把原 `t.Errorf` 一起搬走，免得探针自己吞断言）。探针台四数与 POST-L **逐数相同**（memory `67/35/0/1`+31、config `99/54/0/1`+44）⇒ 探针没改变任何结局。以下字符串逐字抄自该台 `-v` 日志的 `ACQ124 …` 行（形状＝软链形）。

账上点名的两枚 `want ErrSchemaUnmigratable`：

1. `TestMigrationNewerSchemaIsUnmigratable` 实拿：`schema_version=99 (target 2): database was written by a newer Wisp binary; downgrading is not supported`，同枚并读 `IsSchemaUnmigratable=true`。
2. `TestMigrationForeignDatabaseIsUnmigratable` 实拿：`schema_version=-1 (target 2): database has tables but no schema_meta watermark table`，`IsSchemaUnmigratable=true`（该枚第二半 `error should name the missing watermark` 在同一次跑里也通过）。

两枚**都不是**「未解析根」那句。整批 40 枚的软链形日志里 `not provably resolved` / `refusing to seal` 命中数 = **0**。

其余断言里带「拒 / Err / Rejects / Refuses / must fail / never」的枚，逐枚实拿：

3. `TestDeleteArtifactRejectsTheFourHostileShapes`（父）+ 6 枚形状子测试：`memory: invalid artifact name "nested/canary-separator.txt" (bare file names only)`、`"nested-backslash\canary-literal.txt"`、`"../canary-dotdot.txt"`、`"..\canary-dotdot-literal.txt"`、`"/realpriv/w124tmp/TestDeleteArtifactRejectsTheFourHostileShapes84453475/001/Users/carlos/Documents/canary-drive-letter.txt"`、`"\\localhost\$ealpriv/…/canary-unc.txt"` — 逐条都读 `Is(ErrInvalidArtifactName)=true`、`Is(ErrNotFound)=false`（后者是它自己钉的「不许碰到文件系统」那一半）。
4. `TestDeletePrivacyItemRejectsTheFourHostileShapes`：同六枚形状、`route=DeletePrivacyItem`，`memory: invalid artifact name …` 同格式，`Is(ErrInvalidArtifactName)=true`。
5. 裸名腿 `…/dotdot/bare_and_empty` 实拿四条：`memory: invalid artifact name: required (bare file names only)`（空名）、`".."`、`"."`，以及 `"...."` 拿到 `memory: invalid artifact name "...." (a trailing dot or space is a different path on Windows)`。
6. `TestArtifactsDeleteAndPurgeAndGuards`：重复删拿 `artifact "a.txt": memory: not found`；五枚 traversal 各拿 `memory: invalid artifact name …`（`../evil.txt`、`sub\dir.txt`、`..`、空名 → `required`、`a:b.txt`）。
7. `TestProfileDeleteAndPurge` `profile 1: memory: not found`；`TestMemoryAddSearchDeletePurge` `memory 1: memory: not found`；`TestTaskLogLifecycleAndErrorClassValidation` `task_log task-1: memory: not found`；`TestToolCallLifecycleAndValidation` `tool_call 1: memory: not found`；`TestGrantCostPluginStateDAO` 两条 `plugin_state plugin-a: memory: not found`（`Is(ErrNotFound)=true`）与 `memory: revoke grant 999: approval_grant 999: memory: not found`；`TestPrivacyOpsAllDomains` `profile 99999: memory: not found`；`TestRetentionBoundaries` `task_log t-31: memory: not found`（`Is(ErrNotFound)=true`）。
8. `TestWriteQueueLazyLifecycle` 实拿 `memory: store is closed`（它要的是 `ErrStoreClosed`）。
9. `TestAdvPanicWedge` 实拿 `internal: write command panicked: boom: adversarial write panic`，`type=*observe.Error`（它要的是 observe 内部类错，不是密封拒）。
10. `TestStrayRemovalDoesNotFollowLinks` 实拿 `outsideDirErr=<nil> outsidePreciousErr=<nil>` ⇒ 目标树在软链形下也站得住；它改前红在 `openTestStore` 的开库那一步，根本没跑到链接腿。
11. `config` 的 `TestUnwiredGuardLeavesHonestConfigsAlone/tightening a guarded key past its default is still rejected` 实拿：`config: config.toml: risk.shell_enabled is written but does nothing: no shell.exec tool is registered in internal/tools, so nothing reads this flag. it lands with the shell.exec tool itself (SPEC-07 sec 3, S3, default-disabled per D14). Remove the key or set it back to its default - accepting it silently would be a lying config option.` ⇒ 它要的是那把 unwired 守卫的拒，不是底线的拒。
12. `config` 的 `TestResolvedNeverPersists`（「解析出的明文永远不许落盘」）实拿 `carriesPlaintext=false carriesRef=true`。

**结论**：本批没有一枚的断言被「未解析根」那句替代；账上「可转 ⇒ 接解析后仍成立」这一判在本批 40 枚上逐枚复核成立。

## 6. 判据⑤：变异自证（票面 AC#4）— 三态原文

变异 = **拆掉本批新加的那一层解析**（两枚 `sealableTempDir124` 的函数体改成把 `t.TempDir()` 原样交回去；`proc` 引用留一行 `_ =` 保编译）。台子 `/d/tmp/wisp124-2b1-mut`，跑法 `/d/tmp/wisp124-2b1-m.sh`。

```
== mutation landed (grep) ==
internal/memory/tempdir_resolved_124_test.go:38:	_ = proc.SealableRoot // MUTATION-124-2B1: the resolution layer this batch added is removed here
internal/config/tempdir_resolved_124_test.go:33:	_ = proc.SealableRoot // MUTATION-124-2B1: the resolution layer this batch added is removed here
resolution calls still present (want 0):
internal/memory/tempdir_resolved_124_test.go:0
internal/config/tempdir_resolved_124_test.go:0
go build rc=0
go vet rc=0
```

三态读数（同快照血统、同容器、同形状断言）：

| 台 | 形状 | 四数 | 红名 |
|---|---|---|---|
| **POST（未变异）** | 软链 | memory `rc=0 RUN=67 PASS=35 FAIL=0 SKIP=1 SUBPASS=31`；config `rc=0 RUN=99 PASS=54 FAIL=0 SKIP=1 SUBPASS=44` | **无**（0 枚） |
| **MUT-L（拆掉解析层）** | 软链 | memory `rc=1 RUN=36 PASS=2 FAIL=33 SKIP=1 SUBPASS=0 SUBFAIL=0`；config `rc=1 RUN=99 PASS=48 FAIL=6 SKIP=1 SUBPASS=43 SUBFAIL=1` | **40 枚**，逐名 `diff` 与改前基线 PRE-L 的名册**完全相同**（`memory: IDENTICAL name sets`、`config: IDENTICAL name sets`）；红串里 `not provably resolved` 命中 memory 32 + config 6 行，样例逐字：`refusing to seal /varlink/w124tmp/TestRoundTripLoadMarshalLoad17465406/001/.wisp-config-3986474556.tmp: the installed risk.c26Pipeline answered "/varlink/…"` |
| **MUT-P（同一发变异）** | 普通 | memory `rc=0 RUN=67 PASS=35 FAIL=0 SKIP=1 SUBPASS=31`；config `rc=0 RUN=101 PASS=55 FAIL=0 SKIP=0 SUBPASS=46` | 无（四数与 PRE-P/POST-P 逐数相同 ⇒ 这层解析只在软链形起作用，普通形是恒等操作） |

⇒ 红名**点到用例自己**（40 枚逐枚点名，不是包级 `FAIL`），且变异前 `grep` 证落地、`go build` / `go vet` rc=0 先于读数。

## 7. 门禁读数（票面 AC#5）

- **`-count=2 -v` 四数**（只能从 `-v` 量；`COUNT=2` 走同一支 `/h.sh`）：
  - 软链形 POST-L2：`./internal/memory/|rc=0|RUN=134|PASS=70|FAIL=0|SKIP=2|SUBPASS=62|SUBFAIL=0|SUBSKIP=0`、`./internal/config/|rc=0|RUN=198|PASS=108|FAIL=0|SKIP=2|SUBPASS=88|SUBFAIL=0|SUBSKIP=0`
  - 普通形 POST-P2：`./internal/memory/|rc=0|RUN=134|PASS=70|FAIL=0|SKIP=2|SUBPASS=62|SUBFAIL=0|SUBSKIP=0`、`./internal/config/|rc=0|RUN=202|PASS=110|FAIL=0|SKIP=0|SUBPASS=92|SUBFAIL=0|SUBSKIP=0`
  - 两包两形 **FAIL=0、SUBFAIL=0**；每个数都是 `-count=1` 那一发的 2 倍 ⇒ 无 flake、无缓存读数。
- **`gofmt -l`**：容器内对 `fafe2b4` 快照 `gofmt -l internal cmd` ⇒ **0 行**；宿主对全仓同一命令 ⇒ **0 行**。
- **`"$(go env GOPATH)/bin/gofumpt.exe" -l`**：真跑，版本行 `v0.7.0 (go1.27.1)`，`gofumpt -l internal cmd` ⇒ **0 行**（不是「未跑」，也没有错误原文可引）。
- **`go vet` 双 GOOS（逐错误行归因，不整体归因成 cgo 假象）**：
  - `GOOS=windows`（宿主原生，`go vet ./...`）：**rc=0**，错误 **0 行**。
  - `GOOS=linux`（容器原生、`CGO_ENABLED=1`，在 `fafe2b4` 快照里 `go vet ./...`）：**rc=0**；本批两包单跑亦 `rc=0`。
  - `GOOS=linux`（宿主交叉、`CGO_ENABLED` 默认 0，`go vet ./...`）：**rc=1**，全部 2 处错误逐行归因，两处在 pre 锚 `158f729` 的快照里**逐字相同**地存在 ⇒ 与本批无关的既有形状，本批未新增一行：
    1. `package github.com/CarlosShao/wisp/cmd/wisp → imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx → imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in <GOMODCACHE>/github.com/k2-fsa/sherpa-onnx-go-linux@v1.13.8` — 点到的是**外部模块目录**，且 `cmd/wisp/**` 本批零 hunk（2b-4 地界）。
    2. `package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in <tree>/cmd/balldebug` — 该 main 全文件带 windows tag，交叉到 linux 无文件可编。
    - 把这 2 枚包从清单剔掉后重跑 `GOOS=linux go vet $(go list ./... | grep -vE 'cmd/(wisp|balldebug)$')` ⇒ **rc=0、0 行**。⇒ 「rc=1」只由这两枚既有形状贡献，没有一行指向 `internal/memory` 或 `internal/config`。
- **`sh scripts/d22scan.sh`**（纯净快照，`set -eu`、无跳步）：
  - `fafe2b4` 快照：**rc=0**；正向对照那一步 `packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31`，真扫 `examined 225 production Go files`。
  - `158f729` 快照（同一支脚本，对照组）：**rc=0**。
  - 台账八 scope 逐数对照（pre → post）：`bans #1-5 internal/ 203→203`、`bans #1-5 cmd/ 22→22`、`ban #6 frontend/ 40→40`、`ban #7 internal/tools/ 18→18`、`ban #8 design/ 16→16`、`ban #8 frontend/ 40→40`、`ban #8 internal/ 390→`**`392`**、`ban #8 cmd/ 37→37` ⇒ **一枚都不降**，唯一变化是 ban #8 `internal/` +2（本批新增的两枚 `_test.go` 进树）。
  - 与票 131 那本账引用的八数（`203/22/40/18/16/40/390/37`）逐数一致，本批只把第 7 数推到 392。

## 8. 改动面与本批用的 helper

- **用的既有 helper = `proc.SealableRoot`**（`internal/proc/envfork.go:160`），即票 119 生产路 `TestDataDir` / `DefaultLayout` / `cmd/wisp` 的 `resolveDataDir` 走的那一枚。
- 新增两枚**委托函数** `sealableTempDir124(t)`（memory、config 各一枚，函数体一行 `return proc.SealableRoot(t.TempDir())`）。**零新解析路数**：仓内「存在前缀 + `EvalSymlinks` + 尾段原样接回」那本副本账（票 125 的 `R-125-2`，现计三枚：`proc.SealableRoot`、`internal/tools/paths.go`、`winsec/resolve.go` 的 `resolveProbeRoot`）**没有增加第四枚**。之所以每包一枚委托而非共用一枚：Go 的 `_test.go` 不能跨包引用，而把测试仪器提成一枚生产包会新增一条 `memory/config → proc` 的生产依赖边；两枚委托各 1 行、语义同一。
- 递根点共 **13 处**（memory 8、config 5），其中 3 处在既有 helper 内（`openTestStore`、`inv76Fixture`、`writeConfigFile`），一次改就覆盖 40 枚里的 27 枚。
- 未动：`internal/config/c26_seam_posix_125_test.go`（票 125 自拒探针，要的就是未解析那一形）、`internal/config/manager_test.go`（软链形不红，不在本批账面）、`internal/config/private_acl_windows_test.go`（windows-only）、`internal/winsec/**`、`internal/risk/**`、`cmd/wisp/**`。

## 9. 未验证项 + `next=`

**未验证项（照实）**：

1. **macOS 那半仍未实测**（无 macOS runner）— 本批只在 Linux 容器软链 `TMPDIR` 复现；与 AC#1/AC#2a 同一本未付账。
2. **`TestCrashRecoveryKillMidWrite` 的子进程只核到「行」这一层**：父侧断言的是提交行数与 WAL 截断；子进程 `Open` 收到的拼写是否**逐字**等于解析后的根，我只由「改前 0 行 / 改后 35 行 + 本批对那枚用例只改了交给子进程的那一个表达式」推得，没在子进程里打出那枚字符串（要打就得动那枚用例的日志或生产路径，超出本批地界）。
3. **CI 的 ubuntu 腿看不见这一族**（`/tmp` 是真目录）⇒ 本批的绿在 CI 上与改前同数；这仍是票 124 立票时那条「CI 恒绿」的结构事实，未变。
4. **全树两形归零未量**：判据⑤那条终判据（软链形红名数＝0）要 2b-1..2b-4 全落地后才量；本批只裁本批这 40 枚。
5. **`-count=2` 只跑本批两包**；其余 7 包属 2b-2/3/4，本批一字未动（`git diff --name-only` 可核）。

**`next=`（交编排者）**：

- 本批 40 枚已归零 ⇒ 账面余 **91 枚可转**（tools 20 + llm 17 + perm 5 + agent/approval 1 = 43 给 **2b-2**；agent 26 + winsec 17 = 43 给 **2b-3**；cmd/wisp 5 给 **2b-4**，照 16:33 那批次排，且 `acceptor-ticket131-r3` 交回前不动 `cmd/wisp/**`）。
- ⚠ **给 2b-2**：`internal/tools` 的乙类红（`InAllowlist=false` → 升 L2 → `审批未通过`）走的是**路径规范化器**，不是 `winsec` 底线；接解析时先确认「规范化那条入口该不该被测试侧同样喂解析后的根」，别顺手改 `internal/risk/**`（禁改列）。`TestL1WriteGoesThroughTheRealBlockWindow` / `TestLateVetoRendersTheApprovalLayersAppliedStepsReport` 两枚按账归口票 123，**不在**清零目标里。
- ⚠ **给 2b-3（这一族最危险）**：`internal/winsec` 那 17 枚是**底线自己的用例**，票 119 的 `R-119-9` 已写明「不许拿被测函数算 fixture」；`leafLinkTo118` / `cleanSpelling119` 是**故意**用 `filepath.EvalSymlinks` 而不是 `SealableRoot`（后者幂等 ⇒ 期望值恒真）。在那一包直接调 `proc.SealableRoot` 还会拆掉 `envfork.go:147-152` 的边界话（`nothing in internal/winsec calls it`）。⇒ 那一包**别照抄本批这一枚委托**；要么按票 125 的先例在本包自带走法，要么停下来报我。另：`internal/winsec` 软链形有 3 枚 `TestAC2POSIX…125` 从「跑」变「SKIP」，判据③的逐包解释要在那一包重做一遍。
- ⚠ **给 2b-4**：`cmd/wisp` 那 5 枚全在 `providers_test.go` 的 fixture 与 `secret_test.go` 的一枚子测试；同包另有 28 枚两形都红的 DPAPI/命令面腿，**别混进本票的红名数**（AC#1 §5.3 已排除）。
- ⚠ **通用**：`internal/tools`/`llm`/`perm`/`agent` 各自加一行委托是可以的（语义与本批同一枚 `proc.SealableRoot`）；但**若出现第四枚「自己走一遍软链」的实现，那就是 `R-125-2` 的账要翻第三页**，按判据腐坏处理。
- 建议：**AC#2 本格的翻格留到 2b-4 交回后按「软链形红名数＝0（除票 123 那两枚）」一次性复算**；本批只 append 一格进度，不翻 AC#2 / AC#2b 的格。

## 10. 临时件清单（只建不删；唯一例外见末行）

保留：`/d/tmp/wisp124-2b1-{pre,post,probe2,probe3,mut}/`、`/d/tmp/wisp124-2b1-logs/`（`-v` 原始日志 8+ 份、`.summary.txt`、逐名 `.res.txt` 名册、`restored-memory-subtests.txt`、`batch-table.md`、`VET-*.txt`、`D22-{PRE,POST}.txt`、`PROBE3-L.*`）、`/d/tmp/wisp124-2b1-h.sh`、`/d/tmp/wisp124-2b1-m.sh`、`/d/tmp/wisp124-2b1-probe-patch.py`、`/d/tmp/wisp124-2b1-mutate.py`。
⚠ 如实登记一处偏离：`/d/tmp/wisp124-2b1-probe/`（第一版探针树，被我自己的脚本写进一枚 0 字节 `.go` 导致整包编不过）被 `rm -rf` 销毁过一次，随后改用新目录 `probe2`/`probe3` 重跑。**毁的是我自己两分钟前建的坏仪器，不是任何被验版本**；pre / post / probe3 / mut 四棵快照全部在盘可复核。
