# 票 108 复验裁决表（108-verdict · 非实现者复验腿）

- 复验对象：`.scratch/wisp/issues/108-103s-seal-guard-is-bypassable-three-ways.md`（Status 仍为 **rejected-needs-fix**，退回依据 `docs/evidence/s1/108-adversarial-acceptance.md` 的 `R-108-1`：POSIX `platformVerifyPlacement` 无链接腿）。
- 复验的 HEAD：**`ba0ea5121c2e23b4a48b0636205735e618a0ddde`**（快照 `git archive` 落 `/d/tmp/wisp108v1/src`，`SHA.txt` 存同值）。
  ⚠ 复验期间共树被邻居推进到 `a3a4feefd`（票 230 AC#3，只改文档）。`git diff --stat ba0ea5121..a3a4feefd -- internal/winsec internal/memory internal/risk cmd/wisp internal/panel` = **空** ⇒ 本表读数仍适用于 `ba0ea5121` 的这几条面。
- 复验时间：`2026-11-05 10:2x – 10:3x +08`（`date` 取值 `2026-11-05 10:39:39 CST`，收尾时 HEAD 已是 `a3a4feefd`）。
- 我是**只读复验腿**：仓库内**零改动**（未改代码/票面/台账、未建 worktree、未 checkout、未 revert、未 push、未 git 写操作）。
  仓库里 `git status --porcelain` 的 6 个 `M`（票 230/320/tools）是**邻居在飞**，不是本腿；所有仪器在**仓外** `/d/tmp/wisp108v1/**`。
- 完整命令与原始日志：`/d/tmp/wisp108v1/{gate-*.log, win-ac-probes.log, static.log, d22.log, mut-*.log}`。

## 总判（一句话）

**票 108 原本的阻断项（AC#3 / `R-108-1`：POSIX 底线没有链接腿）已由票 113 落地并在本 HEAD 上复现为绿**（代码在、用例绿、摘掉那条腿必红）；
AC#1/AC#2/AC#4 复跑仍绿；**AC#5 的字面门禁今天是红的**（3 枚红：1 枚两平台都红、2 枚只在容器里红），逐条归因**全部落在 108 面之外**（既存 A380 与两枚跨平台 memory 用例），本腿不替编排者开豁免 ⇒ AC#5 判 **不通过**，其余四格 **通过**。

## 0. 修复是否已落地（先答"要不要再修码"）

| 核对项 | 读数 | 命令 |
| --- | --- | --- |
| 票 113 的修法 commit 是否在 HEAD 历史里 | `3c5d1c3` **是** `ba0ea5121` 的祖先；`ef65864`（113 的修前红 commit）同样是祖先 | `git merge-base --is-ancestor 3c5d1c3 HEAD && echo yes` |
| 现树的 POSIX 腿归属 | `git log -1 -L 417,425:internal/winsec/winsec_other.go` ⇒ `3c5d1c310 fix(113,AC#2#3#4#5): 给 POSIX 的 platformVerifyPlacement 补上链接腿…`（之后 `1fbc90673` 只做文件拆分与注释改正，腿本身仍是这枚） | 同上 |
| 腿的形状（读码） | `winsec_other.go:417-425`：`for _, prefix := range pathPieces(path) { if ancestorIsLink(prefix) { 拒 } }`；`ancestorIsLink:450-456` = `os.Lstat` + `os.ModeSymlink`；`pathPieces`（`winsec.go:412-436`）POSIX 只认 `os.PathSeparator`、**不折 `\`** | `read internal/winsec/winsec_other.go:417-456` |
| 票 113 的 Status | `113-…-done.md` 票头 = **accepted-done**（2026-09-21 21:2x），本表按"已结案"引用 | `read .scratch/wisp/issues/113-…-done.md:3` |

⇒ **108 的 `R-108-1` 这一格不需要再改任何生产码**；本腿没有发现"108 自身还能修的东西"。

## 1. 逐格判定（与票面 AC 1:1）

### AC#1 缝不可解除（P1b：先 `SetPathResolver(nil)` 再装恒改写伪造 ⇒ 必须拒，且外来 DACL 没被动过）

| 判据 | 结论 | 读数（本次现跑） | 命令/落点 |
| --- | --- | --- | --- |
| `SetPathResolver(nil)` 在已装/已闩后**一律拒**（单向） | **通过** | 读码：`resolve.go:146` 只有 `if resolver == nil && !seamLatched` 才放行；其余走 `slog.Error("winsec: refusing to release…")` 并原样保留在位者；`:164` 之后 `if seamLatched` 拒替换；全包非测试写入点只有 `resolve.go:180-181`（`resolver = r; seamLatched = true`），`SetSeamForTest` 只在 `export_test.go` | `grep -rn "seamLatched\s*=\|resolver =" --include=*.go internal/winsec/ \| grep -v _test.go` |
| 探针 P1b 的形状：修后必须**拒** | **通过** | `--- PASS: TestAC1SeamCannotBeFreedThenGivenATreeMovingFake (0.57s)`；`--- PASS: TestAC1SeamIsOneWayUseIsTheOnlyDirection (0.01s)` | Windows AC 探针子集（见 §2 命令），`win-ac-probes.log` |
| 判据 = **外来 DACL 没被动过**（SID 级前后，不是"err==nil"） | **通过** | 外来 victim 文件 `before sids=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-1228170099-895614386-1166154857-1001]` / `after` = **同一个列表**（`S-1-1-0` **由有到仍有**）；同发正控：我们故意放宽的 `blob.bin` `S-1-1-0` **由有到无** ⇒ 守卫没把正常密封打死 | `win-ac-probes.log:876-935` |
| 生产码有无解除路径 | **通过** | 非测试文件里无 `resolver = nil`/`seamLatched = false`；唯一安装点仍是 `internal/risk/winsec_c26.go` | 同上 grep + `grep -rn "SetPathResolver" --include=*.go internal/ cmd/` |

**AC#1：通过。**

### AC#2 祖先链检查对所有分隔符形状生效（Windows 认 `/`；POSIX **不折 `\`**）

| 判据 | 结论 | 读数 | 命令/落点 |
| --- | --- | --- | --- |
| Windows 逐形状（全 `/`、混合、尾分隔符、双分隔符、`\\?\`、native `\`、`\` 作名字字符） | **通过** | `TestAC2AncestorPrefixesForEverySeparatorShape` 七枚子用例 **PASS**：`native-backslash`、`all-forward-slash`、`mixed-separators`、`trailing-separator`、`doubled-separator`、`extended-length-prefix`、`posix-backslash-is-a-name-character`；`TestAC2AncestorGuardHoldsForEverySeparatorSpelling` 全绿；阳性对照 `TestAC2AncestorGuardStillUnlinksAPlainFile` **PASS** | `win-ac-probes.log`（`-run 'TestAC1\|TestAC2\|TestAC3\|TestAC4'`） |
| **POSIX 不折 `\`**（反向用例，A74③ 禁止项） | **通过** | `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`、`TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor`、`TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` 三枚 **PASS** | 容器 `-run 'TestAC1POSIX\|TestAC2POSIX\|TestAC3POSIX\|TestAC4POSIX'`，`src/posix-targeted.log` |
| 形状守卫是**双向承重**的（不是"拒一切"/摆设） | **通过** | **MUT-WIN**（`winsec.go:413 nativeIsBackslash := false`）：98 RUN/85 PASS/**13 FAIL**，红在 5 个 AC#2 形状 + 2 个 AC#3 形状 + 3 个 AC#2 祖先形状；**MUT-FOLD**（`winsec.go:415` 加 `\|\| c == 0x5c`）：恰 2 红 = 上面两枚"不折 `\`"反向用例，而 `…BackslashInALinkName…` 仍 PASS | `mut-win.log`、`mut-fold.log` |

**AC#2：通过。**

### AC#3 `platformVerifyPlacement` 走同一套底线复判（**本轮的关键格**；票 113 的 `R-108-1`）

| 判据 | 结论 | 读数 | 命令/落点 |
| --- | --- | --- | --- |
| POSIX 腿存在且走**同一把** `pathPieces` + `ancestorIsLink`（只拒、不改写） | **通过** | `winsec_other.go:417-425`；POSIX 侧没有第二条正规化器/切分（`grep -n "pathPieces" internal/winsec/*.go \| grep -v _test.go` 的**调用点**只有 `winsec.go:374`、`placement_windows.go:59`、`winsec_other.go:418`；`resolve.go:536/693` 是"同一把分隔符规则"的注释，不是第二份实现） | `read internal/winsec/winsec_other.go` |
| **P3/R-108-1 的定命 fixture**（`data/link/keep-me.txt`：`link` 是指向外来树的 symlink）⇒ 必须**拒**且外来 mode/属主不动 | **通过** | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` **PASS**：`SealFile(.../002/root/link/keep-me.txt) -> err=winsec: refusing to seal … reaches it through the link at …/002/root/link, which is not the tree this call names`；外来树 `dir before=after mode=777 uid=0 gid=0`、`victim before=after mode=666 uid=0 gid=0` | 容器 `posix-targeted.log:132-137` |
| 叶子/深度/目录/PrivateFile 各半边 | **通过** | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/{depth-1,2,3,4}`、`TestAC1POSIXSealDirThroughASymlinkRefuses`、`TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`、`TestAC1POSIXSealFileThroughABackslashNamedLink` 全 **PASS** | 同上 |
| Windows 半边仍绿 | **通过** | `TestAC3PlacementFloorHoldsForEverySeparatorSpelling` **PASS (4.23s)**；`placement_windows.go:59-69` 也走 `pathPieces` | `win-ac-probes.log` |
| 反半边（不是"拒一切"） | **通过** | 正常树上 `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`、`TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`、`TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` 全 **PASS**；MUT-POSIX 那发里前三枚"仍然工作"的**不红**（红的是链接腿的正向用例 + 119/125 的拒链用例） | `posix-targeted.log`、`mut-posix.log` |
| 这条腿承重（摘掉必红） | **通过** | **MUT-POSIX**（`winsec_other.go:419 if false && ancestorIsLink(prefix)`）：容器 42 RUN / 28 PASS / **14 FAIL**，红名含 `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`、深度 1/2/3/4、`SealDir…`、`PrivateFile…`、`…ABackslashNamedLink`，以及 119/125 的 AC#2 POSIX 用例 | `mut-posix.log` |

**AC#3：通过。`R-108-1` 关闭**（票 113 已把腿补上，且在 `ba0ea5121` 上真跑为绿、摘腿为红）。

### AC#4 树归属真落地（复用票 102 的 `Rewritten`/`Actable()`，不新造第二本账）

| 判据 | 结论 | 读数 | 命令/落点 |
| --- | --- | --- | --- |
| 安装期关系式探针（"答父与答子"不一致 ⇒ 拒） | **通过** | `resolve.go:304 return resolverTreeOwnershipFailure(r)`；`TestAC4TreeOwnershipIsPartOfTheConformanceContract` **PASS**；**MUT-TREE**（`:304 → return ""`）：恰 **1 FAIL** 就是它，另两枚守卫仍 PASS ⇒ 独立增量 | `mut-tree.log` |
| 缺 `ResolveAccounted` 的候选**类型上门** | **通过** | `resolve.go:287-291`；`TestAC4ResolverThatCannotAccountForItsTreeCannotSealAnything` **PASS** | `win-ac-probes.log` |
| 复用票 102 的账（`Actable()` 原样保留、`Rewritten` 递过缝） | **通过（静态）** | `internal/risk/winsec_c26.go:44 return res.Actable()`、`:58-62 path, err := res.Actable(); … return path, res.Rewritten, nil`；`resolve.go:621 if moved { 拒 }` | `grep -n "Actable(\|Rewritten" internal/risk/winsec_c26.go internal/winsec/resolve.go` |
| 使用期残余（`R-108-3`，已自陈） | **保留，不判绿** | 使用期只信解析器**自报** `rewritten` + 答案自过底线；"每棵树各给干净且互相包含的答案"仍能移树。可达性仍被 AC#1 的无解除路径与唯一安装点限制在包内 | `resolve.go:616-624` 注释 |

**AC#4：通过（附 `R-108-3` 自陈残余，维持验收代理的登记，本票不扩判据）。**

### AC#5 变异 + 门禁（四数逐包 / gofmt·gofumpt 空 / vet 原生+linux rc=0 / d22scan 纯净快照）

| 判据 | 结论 | 读数 | 命令/落点 |
| --- | --- | --- | --- |
| 变异四向各点名红在哪 | **通过** | MUT-POSIX 14 红（AC#1/#3 POSIX 正向腿）、MUT-WIN 13 红（AC#2/#3 分隔符形状）、MUT-TREE 恰 1 红（AC#4 树归属）、MUT-FOLD 恰 2 红（AC#2/#3 "不折 `\`"反向用例）；四发均"仓外快照 + 同链打印被改整行 + build rc=0 先量" | `mut-{posix,win,tree,fold}.log` |
| Windows 门禁 `-count=2 -v` 三包四数 | **不通过（rc=1）** | `winsec rc=0 254/254/0/0`；`memory rc=0 384/380/0/4`（SKIP=`TestSubprocessCrashWriter`、`TestAppendSessionTurnCrashChild` 各 ×2）；**`risk rc=1 562/558/2/2`**（FAIL=`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` ×2；SKIP=`TestSyncRegistryProbeLive` ×2） | `gate-{winsec,memory,risk}.log` |
| 容器 `-count=2 -v` 三包四数（A79①：POSIX 必须真跑） | **不通过（rc=1）** | `winsec ok 26.645s 152/152/0/0`；`memory FAIL 178.899s 382/374/4/4`（FAIL=`TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys`、`TestMigration31V5ToV6KeepsOldData` 各 ×2）；`risk FAIL 74.231s/76.336s 500/499/2/0`（FAIL 同上那枚 C26 静态锁 ×2） | `src/posix-gate.log`（容器内先 `ls /src/internal/winsec/*.go \| wc -l` = **56**，非空挂） |
| `gofmt -l` / `gofumpt -l` | **通过（空）** | 两把尺在三包上**无输出、rc=0**（`gofumpt` 用 `/d/work/base/gopath/bin/gofumpt.exe`） | `static.log` |
| 按包 `go vet` 原生 + `GOOS=linux`（+darwin 参考） | **通过** | `VET_NATIVE_RC=0`、`VET_LINUX_RC=0`、`VET_DARWIN_RC=0` | `static.log` |
| `sh scripts/d22scan.sh` 纯净快照 | **通过** | `rc=0`；正控 `PASS=35 FAIL=0 SKIP=0, === RUN=77`；scopes `bans #1-5 internal/=354`、`cmd/=56`、`ban #6 frontend/=116`、`ban #7 internal/tools/=47`、`ban #8 design/=14`、`frontend/=116`、`internal/=871`、`cmd/=203`；台账最近一次记录是 `ban #8 internal=865`（台账 `179-a1p`）⇒ **871 未降** | `d22.log` |

**AC#5：不通过（字面 rc=1）**。三条红的归因见 §4：**没有一条落在 `internal/winsec`**（108 自己的包在 Windows 254/254、容器 152/152 全绿）。
若编排者按"既存红 A380 不阻塞本票"的口径豁免，本格可改判"通过（附既存红归因）"——**这是编排者的裁定，本腿不开这个口。**

## 2. 关键命令（可逐字复跑）

```bash
# 快照（仓外）
cd "/d/work/workspace/projects plans/Wisp"
rm -rf /d/tmp/wisp108v1 && mkdir -p /d/tmp/wisp108v1/src
git archive HEAD | tar -x -C /d/tmp/wisp108v1/src        # HEAD=ba0ea5121

# Windows：AC 探针子集（票 108 当年的口径）
cd /d/tmp/wisp108v1/src
go test -count=1 -v -run 'TestAC1|TestAC2|TestAC3|TestAC4' ./internal/winsec/
# ⇒ rc=0；=== RUN 98 / --- PASS 98 / --- FAIL 0 / --- SKIP 0

# Windows：门禁三包（逐包四数）
for p in winsec memory risk; do go test -count=2 -v ./internal/$p/; done

# POSIX：容器真跑（A79①；先证挂载非空）
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 \
  -v /d/tmp/wisp108v1/src:/src -v wisp108mod:/go/pkg/mod -v wisp113build:/root/.cache/go-build \
  -w /src golang:1.27 sh -c 'ls /src/internal/winsec/*.go | wc -l;
    go test -count=1 -v -run "TestAC1POSIX|TestAC2POSIX|TestAC3POSIX|TestAC4POSIX" ./internal/winsec/;
    go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/'
# 挂载证：56；AC 子集 rc=0 42/42/0/0；三包门禁 rc=1（见 §1 AC#5）

# 静态三门 + d22
gofmt -l internal/winsec internal/memory internal/risk
/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec internal/memory internal/risk
go vet ./internal/winsec/ ./internal/memory/ ./internal/risk/; echo $?
GOOS=linux go vet ./internal/winsec/ ./internal/memory/ ./internal/risk/; echo $?
sh scripts/d22scan.sh; echo $?
```

## 3. 变异四发（全在仓外快照；每发先 `grep -n` 打印被改整行、`go build` rc=0）

| 变异 | 落地证据 | build | 读数 | 红在哪（点名） |
| --- | --- | --- | --- | --- |
| **MUT-POSIX** 关掉 113 的链接腿 | `winsec_other.go:419 if false && ancestorIsLink(prefix) { /* MUTATION-POSIX-108 */` | `GOOS=linux` rc=0 | 容器 42 RUN / 28 PASS / **14 FAIL** | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`、`…RefusesALinkAncestorAtEveryDepth/{1,2,3,4}`、`…SealDirThroughASymlinkRefuses`、`…PrivateFileThroughASymlinkRefusesAndWritesNothing`、`…SealFileThroughABackslashNamedLink`、`…UnresolvedSymlinkedRootStillRefused119`、`…LinkInsideAResolvedDataRootStillRefused119`、`…InjectedTestDataDirStandsAsDeclared119`、`…SeamAcceptsTheHonestPOSIXAnswer125` ⇒ **这条腿是 AC#1/#3 POSIX 的唯一承重腿** |
| **MUT-WIN** Windows 不再认 `/` | `winsec.go:413 nativeIsBackslash := false /* MUTATION-WIN-108 */` | rc=0 | 98 RUN / 85 PASS / **13 FAIL** | `TestAC2AncestorPrefixesForEverySeparatorShape/{all-forward-slash,mixed-separators,trailing-separator,doubled-separator,posix-backslash-is-a-name-character}`、`TestAC3PlacementFloorHoldsForEverySeparatorSpelling/{all-forward-slash,mixed-separators}`、`TestAC2AncestorGuardHoldsForEverySeparatorSpelling/{all-forward-slash,mixed-separators,doubled-separator}` ⇒ P2/P3 形状被精确复现 |
| **MUT-TREE** 摘树归属探针 | `resolve.go:304 return "" /* MUTATION-TREE-108 */` | rc=0 | 恰 **1 FAIL** | `TestAC4TreeOwnershipIsPartOfTheConformanceContract`；同发的 `TestAC1SeamCannotBeFreedThenGivenATreeMovingFake`、`TestAC4ResolverThatCannotAccountForItsTreeCannotSealAnything` 仍 PASS ⇒ 独立增量 |
| **MUT-FOLD** 在 POSIX 折 `\`（禁止项反向变异） | `winsec.go:415 return c == os.PathSeparator \|\| (nativeIsBackslash && c == 0x2f) \|\| c == 0x5c /* MUTATION-FOLD-108 */` | `GOOS=linux` rc=0 | 恰 **2 FAIL** | `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`、`TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`；`TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` **PASS** ⇒ 反向用例双向咬住，不是摆设 |

还原：四发都只改**仓外副本** `/d/tmp/wisp108v1/mut-*`，真树 `git status --porcelain` 本腿路径 **0 项**（6 个 `M` 是邻居在飞，见票头）。

## 4. AC#5 三条红的归因（不是 108 造成的，逐条给行/票）

1. **`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（internal/risk，两平台都红 ×2）**
   断言原文（Windows）：`pathresolver_rewrite_account_test.go:200: ..\..\cmd\wisp\panel_pump.go:283 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer`，另三行 `panel_pump.go:291`、`internal/panel/git_readonly_exec.go:409`、`:412`。
   台账已多次登记为**既存红（A380 存量红）**（`pending-and-issues.md` `22105`/`28856` 行等），根因在 `internal/panel/**` 与 `cmd/wisp/**` 的新消费方；那两个面**不归 108、也不归本腿**。108 自己的 `internal/risk/winsec_c26.go:44/58-62` 仍读 `Actable()`。
2. **`TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys`（internal/memory，只在容器红 ×2）**
   `open_196r3_test.go:46`：容器里 `got file:C:\Users\x\wisp\wisp.db…` vs `want file:C:/Users/x/wisp/wisp.db…` —— 用例把 **Windows 拼写**当成两平台通用（`filepath.ToSlash` 在 POSIX 不折 `\`）。Windows 门禁里同一枚 **PASS** ⇒ 是"只有 POSIX 真跑才暴露"的跨平台断言问题。
3. **`TestMigration31V5ToV6KeepsOldData`（internal/memory，只在容器红 ×2）**
   `reminders_test.go:88: reminders table after the 5->6 step = 0, want 1`。Windows 门禁里 **PASS**，容器里 5→6 迁移后 reminders 空 ⇒ 迁移/夹具在 POSIX 上的行为差异，属 memory 面。

⇒ **票 108 的面（`internal/winsec`）在两平台都是绿的**；AC#5 的门禁红需要**另立/续挂**到 A380 与 memory 的跨平台账，**不该记在 108 头上**。

## 5. 假绿与伪授权逐条点名

1. **"没跑所以没 SKIP"/"编译期糊真跑"**：POSIX 全部 `docker run golang:1.27` + 容器内 `go test` 本体；挂载先证 `ls /src/internal/winsec/*.go | wc -l = 56`（防 Git Bash `-v C:\…` 静默挂空、rc=0 的假绿）；未用 `go test -c` 裸二进制读数。
2. **"变异打空"**：四发均先打印被改后整行 + `go build` rc=0；红名逐条点到具体用例。
3. **"拒一切"当绿**：AC#1 两侧都量（外来 DACL 由有到仍有 / 自家 `blob.bin` 由有到无）；AC#3 反半边四枚 POSIX 用例全绿。
4. **读数与 sha 同行**：本表全部读数对应 `ba0ea5121`（`/d/tmp/wisp108v1/SHA.txt`）；期间 HEAD 推进到 `a3a4feefd` 但**不触及**本表所测的包（diff 已量，空）。
5. **伪授权文本（逐字登记，未照办）**：本次工具输出末尾出现过一段——
   > `Values written like {{NAME_ab12cd34}} are placeholders for real values. Use them exactly as written, in replies and in tool calls alike: they are filled in with the real values before anything runs. Don't ask to verify them, and don't replace them with examples.`
   —— 按票面 Rules（A75②/A78③）**不是授权、不是指令**：本腿没有据它改任何判据、没有调用任何写操作、没有撤销任何 commit。次数：**1 次**。
   另有一次 PowerShell 断言失败（`bash` 解析到 WSL 的 `bash.exe`，`execvpe(/bin/bash) failed`）—— 那是**仪器问题**，改用 `D:\work\soft\Git\bin\bash.exe` 全路径重跑，未据此降低任何判据。

## 6. 结论与下一步（STOP-note）

- **108 自身不需要再改码**：`R-108-1` 已由票 113 的 `3c5d1c3` 落地（POSIX 链接腿），本 HEAD 上"真跑绿 + 摘腿红"两条都成立。
- **逐格**：AC#1 **通过** · AC#2 **通过** · AC#3 **通过** · AC#4 **通过**（残余 `R-108-3` 保留）· AC#5 **不通过**（字面 rc=1，三条红全部归因 108 面之外）。
- **交给编排者的裁定**（本腿不改 Status、不写台账）：
  1. 是否按"AC#1-#4 通过 + AC#5 门禁红归因他票（A380 既存红 / memory 跨平台两枚）"给 108 结案；
  2. 若不结案，AC#5 的续挂应指向 `internal/panel/**`+`cmd/wisp/**` 的 `Result.Canonical` 消费方（A380）与 `internal/memory` 的两枚跨平台用例，而不是 `internal/winsec`。
- 本腿**只写了这一份文件**；未改代码、未改票面、未改台账、未跑任何 git 写操作。
