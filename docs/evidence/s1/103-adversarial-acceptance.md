# 票 103 独立对抗验收裁决表 —— 密封缝守卫（R-c，今天可利用）与 `RemoveUnlinked` 祖先检查 + tripwire（R-b，陷阱档）

**验收代理：** `acceptor-ticket103`（独立对抗，未参与实现，未看过实现方会话）
**验收日期：** 2026-09-21
**被验收的码：** 修前红用例 `184af22`；生产码 `0717bf2`（`internal/winsec/resolve.go`、`winsec.go`、`winsec_windows.go`、`winsec_other.go`）
**票面：** `.scratch/wisp/issues/103-seal-seam-is-unguarded-and-removeunlinked-follows-junctions.md`
**开工时仓库：** branch `dev` @ `8b6f691`；收工时 HEAD 已被别人推到 `4b4cfc9`（票 93/95/99 的活）——两枚 103 commit 均在 HEAD 里：
`git merge-base --is-ancestor 0717bf2 HEAD` ⇒ `0717bf2 in HEAD ✓`、`184af22 in HEAD ✓`。
**共树登记：** `git status --porcelain` 里 `internal/config/parse.go`、`migrate.go`、`internal/risk/syncdirs*_test.go`、`.github/workflows/ci.yml` 是别人的未提交活 ⇒ 本表所有 Go 命令一律**按包**跑，或在 `/tmp` 的 `git archive HEAD` 快照里跑（快照 `C:/Users/swq/AppData/Local/Temp/wisp-ac103`、`wisp-ac103b`、`wisp-ac103c`，带会话后缀；**仓内未建 worktree、未 checkout**）。`internal/winsec/` 顶端无未结验收挡路（票 89 的 `01e7007` 已结）。

**档位图例**
- 〔独立复现〕本代理亲自敲命令、亲自读到 rc 与原文行。
- 〔日志＋归档，我抽验〕实现方给了读数，本代理抽验其中可核部分。
- 〔仅自述，不背书〕只有实现方日志、本代理未能复现 ⇒ **必写补救动作**。

---

## 裁决总表（与票面 AC 1:1）

| AC | 票面判据要点 | 本代理主证据（一句话） | 结论 | 档位 |
|---|---|---|---|---|
| **AC#1（R-c）** | 缝只能装一次 + 只能装真解析器；两条腿；"什么都没发生"要 SID 级读数 | 两条腿与 SID 级前后对照本代理全部复现（外来 `S-1-1-0` 前后都在）；但**换个注册顺序**（先 `nil` 释放再装恒改写型伪造）能过探针并真的改掉别人 DACL（P1b） | **通过但有条件**（R-103-1、R-103-4） | 〔独立复现〕 |
| **AC#2（R-b）** | junction 输入 ⇒ 拒 + 不删目标；memory 侧"下降即红"tripwire | 平台用例 + tripwire 均复现（含真机 `mklink /J`）；POSIX 侧本代理在**真 Linux 容器**里量到同一守卫生效；但**分隔符改写**（`/`）绕过 `firstLinkAncestor`，外国文件仍被删且返回 nil（P2） | **通过但有条件**（R-103-2、R-103-3） | 〔独立复现〕 |
| **AC#3** | 变异三向：去守卫⇒红；被拒⇒仍有一条绿；tripwire 改下降⇒红 | 三向本代理**自己重做**：MUT-1 两 FAIL（红在守卫）；MUT-2 恰 1 FAIL 且 witness 绿；MUT-3 tripwire FAIL；还原 5 文件 `diff -q` 全 CLEAN 且快照复跑 `ok` | **通过** | 〔独立复现〕 |
| **AC#4** | 四包回归 + 逐条点名 SKIP/FAIL；gofmt/gofumpt；按包 vet；`sh scripts/d22scan.sh` | 三门本代理复跑全绿；`internal/risk` 在我这轮 **FAIL 1 条**（`TestResolvePerCallBudget`，隔离复跑 PASS ⇒ 票 86 负载假红，两样本都报）；d22scan rc=0、各 scope 文件数不降；**CI 侧 ticket 无关的红仍在（tools/secret/windows 环境）** | **通过但有条件**（假红 1 条已点名；CI 未"过"，见 §AC4-E） | 〔独立复现〕＋〔日志＋归档，我抽验〕 |

**总判：PASS WITH CONDITIONS（通过但有条件）。** AC#1/AC#2 的**判据本身**都立住了（红→绿、SID 级读数、真机 junction、tripwire 咬得住下一步），三向变异与门禁本代理独立复现；但两处守卫都**留有可复现的门**（P1b 注册顺序、P2/P3 分隔符），故不给无条件绿。四格 AC 框**不建议由编排者无条件勾选**：AC#1、AC#2 应勾在"通过但有条件 + `R-103-1..4` 各自开票"之下。

---

## §AC1 —— R-c：`SetPathResolver` 的三道守卫

### A. 票面原句
> seam 只能被**装一次**、且装的必须是**真解析器**：给注册口加守卫（幂等/一次性 + 类型上不给伪造留门，或伪造时**响亮失败并审计**）。判据用例两条腿：**装第二个 ⇒ 红**；**装一个恒说 OK 的 ⇒ winsec 必须拒**，且"什么都没发生"不算绿（要能证明**外来 DACL 没被改**，取 SID 级读数）。

### B. 形状保留（票 94 的账没清零）〔独立复现〕
```
grep -rn "SetPathResolver" --include=*.go .        → 生产侧唯一安装点 internal/risk/winsec_c26.go:21（init）
grep -rn "filepath.Abs\|filepath.Clean" internal/winsec/  → 只剩两处注释（resolve.go:15、winsec.go:159），无调用
```
`ResolvePath` 仍是 `ResolvedPath` 的唯一 mint（`internal/winsec/resolve.go:198-227`），`privateDirAll` 仍收 `ResolvedPath`（`internal/winsec/winsec.go:160-201`）。**没有倒回 `filepath.Abs`。**

### C. 两条腿 + 真机 junction + SID 级读数
命令：`go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/risk/`（真树，rc 见 §AC4）
```
=== RUN   TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone   → --- PASS (0.88s) / --- PASS (0.47s)   （-count=2 两样本）
=== RUN   TestAC1SeamIsSingleUse                                       → --- PASS (0.00s) ×2
=== RUN   TestAC1RefusedInstallLeavesTheSealWorking                    → --- PASS (0.69s)/(0.42s)
```
真机链接确实在位（不是我推断的）：`seam_guard_windows_test.go:139: junction in place: ...\002\data\link (attributes 0x410)`。
**"什么都没发生"的读数（SID 级，`icacls`，不是返回码）**：
```
seam_guard_windows_test.go:176: foreign file DACL before: sids=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]
seam_guard_windows_test.go:187: foreign file DACL after:  sids=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]
seam_guard_windows_test.go:166: PrivateDirAll through the junction refused as: winsec: refusing to seal ...\data\link\artifacts:
                                risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions
```
⇒ `S-1-1-0`（Everyone）**由有到仍有**，且 leg 2 的拒因**点名 reparse**（不是"偶然红"）。这条判据我认可：**通过**。

### D. 守卫是不是幂等的？（两种都可接受，但必须"写出来并钉住"）〔独立复现〕
读码：`internal/winsec/resolve.go:113-125` —— 缝上已有人时，**同名 `%T` 再装 ⇒ 静默保留第一份 + `slog.Debug`（line 123-124），不报错**；不同名 ⇒ ERROR 审计 + 保留在位者（line 116-121）。语义写在了注释（line 78-80）与代码里 ⇒ **"写出来"这一半满足**。
"钉住"这一半**不满足**，我用覆盖率仪器量的，不是猜的：
```
go test -count=1 -coverprofile=/tmp/cov-winsec.out ./internal/winsec/   rc=0
grep "resolve.go:1[0-2][0-9]\." /tmp/cov-winsec.out
  resolve.go:116.3,116.66   1 1      ← 不同名拒装分支：跑过
  resolve.go:123.3,124.9    2 0      ← "installed again identically" 分支：**执行次数 0**
  total: (statements) 78.6%
```
⇒ **`R-103-4`**：恒等重装这条分支树内**没有任何用例走进去**，它对"幂等语义"的钉法是漂的（今天静默保留，明天有人改成报错，没有用例会红）。

### E. 攻守卫：换类型 / 换注册顺序〔独立复现，探针只存在于 /tmp 快照〕
快照文件 `internal/winsec/ac103probe_windows_test.go`、`ac103probe2_windows_test.go`（仓外，仓库未收）。

**P1（同顺序，缝已被真解析器占位）——拒得对：**
```
TestAC103P1... → --- PASS
ERROR winsec: refusing to replace the already installed sealing path resolver
  installed=risk.c26Pipeline attempted=winsec_test.constRewritingResolver reason="the seam is single-use..."
```
顺带量到**类型名撞名不可利用**：`%T` 撞名只会走到"静默保留在位者"，伪造永远进不来（方向是 fail-safe）。

**P1b（换个注册顺序：先 `SetPathResolver(nil)` 释放缝，再装恒改写型伪造）——守卫留了门：**
```
go test -count=1 -v -run 'TestAC103P1b|TestAC103P3' ./internal/winsec/  → rc=1
P1b seam after the unconditional nil reset: <floor>
P1b seam after offering the const-rewriting fake: winsec_test.constRewritingResolver
P1b R-103-1: the conformance probe ACCEPTED a fake that rewrites every answer to one fixed clean absolute spelling
P1b SealFile(...\data\blob.bin) -> <nil>
P1b foreign SID before: [S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]
P1b foreign SID after:  [S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]
P1b R-103-1 CONFIRMED at SID level: S-1-1-0 stripped from the foreign file
icacls(foreign) after:  ... NT AUTHORITY\SYSTEM:(F) / BUILTIN\Administrators:(F) / DESKTOP-LVS7839\swq:(F)   ← Everyone 那条没了
--- FAIL: TestAC103P1bFreeTheSeamThenInstallTheRewritingFake (2.58s)
```
⇒ **`R-103-1`**。机理：`nil` 重置是**无条件允许**的（`resolve.go:97-104`），一次性守卫因此可被**顺序**解除；一致性探针的接受集是"拒 **或** 交回底线认的改写"（`resolve.go:157-162`），恒改写型伪造恰好落在"改写"那一侧；`ResolvePath` 对答案再跑一遍底线（`resolve.go:222`）只挡得住**带链接/`..`/空段/`\?\`**的答案，挡不住"改写到另一棵干净树"。
**档位说明（诚实）**：这需要包外已能调用导出符号，与票 94 PROBE A 同一威胁模型 ⇒ 不是"比原缺陷更难"的门，而是**同一威胁模型下守卫没覆盖到的第二类伪造**；票面 AC#1 的"装的必须是**真解析器**"这一句按字面判 ⇒ 未完全达成。

### F. AC#1 结论
判据的**两条腿 + SID 级读数**本代理独立复现 ⇒ **通过**；守卫留了 P1b 那一扇门、恒等重装未钉 ⇒ **有条件**。条件＝`R-103-1`、`R-103-4`（各应自开一票）。

---

## §AC2 —— R-b：`RemoveUnlinked` 祖先检查 + `internal/memory` tripwire

### A. 票面原句
> 给 `RemoveUnlinked` 一条**平台用例**：输入一个 junction/符号链接 ⇒ **必须拒、返回错误且不删目标**；并在 `internal/memory/removeStray` 一侧补一条**"下降进链接就会红"的 tripwire**（断言遍历不下降）。⚠ 修语义若需要动 `internal/memory` ⇒ 停手登记交回，本票只做守卫与 tripwire。

### B. 平台用例（真机 junction，只在临时目录造，测完由 `t.TempDir()` 清）〔独立复现〕
```
=== RUN   TestAC2RemoveUnlinkedRefusesAPathThroughAJunction  → --- PASS (0.46s)/(0.19s)
=== RUN   TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink   → --- PASS (0.38s)/(0.21s)   ← 阳性对照：独立链接照旧解链
```
守卫的原文读数（我自己在快照里复跑控制腿拿到的）：
```
P2 control (backslash) refused as: winsec: entry is a link to something else, not private data
  ...\data\link\sub\keep-me.txt: the spelling reaches it through the link at ...\data\link,
  and whatever lives behind that link is not this tree's data to delete
```
`firstLinkAncestor`（`winsec.go:249-282`）+ 平台谓词（`winsec_windows.go:506` `isReparsePoint` / `winsec_other.go:87-93` `ModeSymlink`）——**不含叶子**的祖先链，注释与实现一致。

### C. tripwire 挡不挡得住"接上之后的第一步"（票面 AC#3 第③向，我自己做的）
见 §AC3-C：**MUT-3 下 `TestAC2ReclaimWalkNeverDescendsIntoAJunction` FAIL**，红在
`removeStray(junklink): ... winsec: entry is a link ... reaches it through the link at ...\artifacts\junklink`
⇒ 一旦有人把 `removeStray` 的遍历接成下降的，**立刻响亮红**（且因为守卫在位，是"拒"不是"别人树里的文件静默消失"）。同时票 79 既有的 `TestStrayRemovalDoesNotFollowLinks` 在 MUT-3 下**仍 PASS** ⇒ 新 tripwire 是真增量（实现方这点自述属实，我复现了）。

### D. 攻守卫：分隔符改写〔独立复现，快照〕
```
TestAC103P2SeparatorSurgeryOnRemoveUnlinked → --- FAIL (0.24s)
P2 RemoveUnlinked("C:/Users/.../002/data/link/sub/keep-me.txt") -> err=<nil> ;
   victim stat after=GetFileAttributesEx ...\victim\sub\keep-me.txt: The system cannot find the file specified.
P2 R-103-2 CONFIRMED: spelling ... reached through the junction and deleted the foreign file (err=<nil>)
P2 RemoveUnlinked("C:\\...\\data\\link/sub/keep-me.txt") -> err=<nil> ; 外国文件同样消失（第二条腿）
```
⇒ **`R-103-2`**。机理：`firstLinkAncestor` 只按 `filepath.Separator`（Windows = `\`）切段（`winsec.go:250-259`），全 `/` 或混合拼写在 Windows 上被 OS 当分隔符解析，但在守卫眼里**只有一个组件＝叶子**，叶子按设计不查 ⇒ 一个祖先都不 Lstat。
**档位说明**：`internal/memory/removeStray` 用 `filepath.Join` ⇒ 今天只会喂反斜杠拼写（我读了 `artifacts.go:155/170/180`），与 R-b 同档＝**今天够不到的陷阱**，不是新可利用面。但票面 AC#2 判据写的是"输入一个 junction/符号链接 ⇒ 必须拒、返回错误且不删目标"，按字面 ⇒ 未无条件成立。

### E. 顺手量到的另一半（不是 103 的码，但在同一判据里）〔独立复现〕
```
TestAC103P3FloorReparseWalkVersusSlashSpellings → --- FAIL (0.33s)
P3 control (backslash) refused: winsec: ... traverses a reparse point at ...\data\link, which is not the tree this call names
P3 SealFile("C:/Users/.../002/data/link/sub/keep-me.txt") -> <nil>
P3 foreign SID before=[S-1-1-0 ...] after=[S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]
P3 R-103-3 CONFIRMED: the slash spelling got past the floor's reparse walk and rewrote the foreign DACL
```
⇒ **`R-103-3`**：内置底线 `platformVerifyPlacement`（`internal/winsec/placement_windows.go:52`）同样只按 `\` 切段做 reparse 走查。**该文件不在 `0717bf2` 的 diff 里**（`git show --stat 0717bf2` ⇒ 只 resolve.go/winsec.go/winsec_other.go/winsec_windows.go）⇒ **不是 103 引入的回归**，是票 94 底线的既有洞，被 103 的守卫"同一个错法"复制了一份。登记交回，本代理不修。

### F. 越界检查（AC#2 不许动 `internal/memory` 语义）〔日志＋归档，我抽验〕
`git show --stat 0717bf2` 的 4 个文件**全在 `internal/winsec/`** ⇒ "一行 `internal/memory` 生产码都没动"**属实**。tripwire 是 `_test.go`，在允许面内。

### G. AC#2 结论
平台用例 + tripwire 两半都达成且**能被别人重跑**；守卫留有 `R-103-2`（分隔符）⇒ **通过但有条件**。

---

## §AC3 —— 变异三向（本代理自己重做，不引用实现方读数）

仪器：`/tmp/wisp-ac103b`＝`git archive HEAD` 快照；每发**同一条 `&&` 链里 `grep -n` 打印被改后整行** ⇒ `go build` rc=0 先量到（编译失败不算变异）⇒ `go test -v` 数 `=== RUN`。**非 `-v` 看不到 PASS，全部用 `-v`。**

### MUT-1（AC#3①：去掉注册守卫 ⇒ AC#1 红）
```
115:	if resolver != nil && false {
152:	for _, probe := range []string{} {
build rc=0
=== RUN   TestAC1BaselineProductionPaths                                      --- PASS (1.46s)
=== RUN   TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites               --- PASS (6.02s)
=== RUN   TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone
    seam_guard_windows_test.go:149: AC#1 leg 1 RED: a pass-through rubber stamp was installed into the sealing seam (now winsec_test.rubberStampResolver)
    --- FAIL (1.44s)
=== RUN   TestAC1SeamIsSingleUse
    seam_guard_windows_test.go:233: AC#1 RED: the seam is not single-use - winsec_test.narrowOnlyResolverB replaced winsec_test.narrowOnlyResolver
    --- FAIL (0.00s)
=== RUN   TestAC1RefusedInstallLeavesTheSealWorking                           --- PASS (0.75s)
```
⇒ **5 条 `=== RUN`，3 PASS / 2 FAIL**，两枚红**都点名在守卫上**。（比实现方报的 3 条多，因为我的 `-run` 前缀 `TestAC1` 多带了两枚票 89 基线用例，两枚都 PASS——不影响判据。）

### MUT-2（AC#3②：注册了但被拒、只拆审计 ⇒ 必须仍有一条绿）
```
107:		slog.Debug("winsec: refusing to install a path resolver into the sealing seam",
build rc=0
--- PASS: TestAC2SealedDirCoversFilesItNeverTouched (2.44s)
--- PASS: TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs (0.68s)      ← 票 89 那枚，没被改松
--- PASS: TestAC2MigrationBackupIsPrivate (1.32s)
--- PASS: TestAC2PreExistingMigrationBackupIsRepaired (0.42s)
    seam_guard_windows_test.go:152: AC#1 leg 1 RED: refusing the fake left no audit record naming it; log was:   ← 唯一红
--- FAIL: TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone (0.44s)
--- PASS: TestAC1SeamIsSingleUse / TestAC1RefusedInstallLeavesTheSealWorking
--- PASS: TestAC2RemoveUnlinkedRefusesAPathThroughAJunction / TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink
```
⇒ **8 条 `=== RUN`，7 PASS / 恰 1 FAIL**，且红只咬审计断言 ⇒ **MUT-1 的红在守卫身上，不在噪声身上**；witness `TestAC1RefusedInstallLeavesTheSealWorking` 绿。

### MUT-3（AC#3③：tripwire 的"不下降"改成"下降" ⇒ 该用例红）
落地行 + 新增下降递归（`internal/memory/mut3_descend_windows.go`，`os.Stat`+`os.ReadDir` 跟链接）：
```
155:	err := descendFollowingLinks103(filepath.Join(s.artifactsDir, name),
build rc=0
=== RUN   TestAC2ReclaimWalkNeverDescendsIntoAJunction
    artifacts_junction_tripwire_windows_test.go:68: removeStray(junklink): memory: reclaim artifacts stray "junklink":
      memory: remove artifacts stray file "junklink\\sub\\keep-me.txt": winsec: entry is a link to something else,
      not private data ...\artifacts\junklink\sub\keep-me.txt: the spelling reaches it through the link at ...\artifacts\junklink
--- FAIL: TestAC2ReclaimWalkNeverDescendsIntoAJunction (0.18s)
=== RUN   TestStrayRemovalDoesNotFollowLinks
--- PASS: TestStrayRemovalDoesNotFollowLinks (0.17s)     ← 票 79 既有用例钉不住下降，实现方这点属实
```
⇒ 三向齐。**AC#3 结论：通过。**

### 还原证明
```
CLEAN  internal/winsec/resolve.go / winsec.go / winsec_windows.go / winsec_other.go / internal/memory/artifacts.go   （逐个 git show HEAD:<f> | diff -q）
diff -rq 新 archive(wisp-ac103c) vs wisp-ac103b → 只差我的两枚探针文件 + 别人在期间落地的 ci.yml / syncdirs*_test.go / portable-tests.sh（HEAD 前移）
go test ./internal/winsec/ ./internal/memory/ → ok winsec 7.127s / ok memory 13.519s
```
仓库目录内未建 worktree、未 checkout。

---

## §AC4 —— 门禁复跑（本代理亲自敲的）

### A. `go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/risk/` ⇒ **rc=1**
```
ok    github.com/CarlosShao/wisp/internal/winsec   20.829s
ok    github.com/CarlosShao/wisp/internal/memory   30.163s
ok    github.com/CarlosShao/wisp/internal/secret    1.400s
FAIL  github.com/CarlosShao/wisp/internal/risk     11.315s
--- FAIL: TestResolvePerCallBudget (1.50s)
```
**逐条点名 FAIL（1 条）**：`internal/risk/TestResolvePerCallBudget`。隔离复跑第二个样本：
```
go test -count=1 -run TestResolvePerCallBudget -v ./internal/risk/
    pathresolver_budget_norace_test.go:34: C26 Resolve: 779329 ns/op = 0.779 ms/op (budget 1.000 ms, 1677 samples)
--- PASS: TestResolvePerCallBudget (1.48s)   → ok risk 1.564s
```
⇒ **票 86 已登记的负载假红**（四包并发压测时预算 1.000 ms/op 被越过），**不是 103 的缺陷**；两样本都报，不重测到运气好。
**逐条点名 SKIP**：`--- SKIP: TestSubprocessCrashWriter`（×2，`internal/memory`，既有环境条件用例）与 `--- SKIP: TestSyncRegistryProbeLive`（×2，`internal/risk`）。
⇒ 实现方"`=== RUN` 102 条、SKIP 恰 1 条"的核对：`go test -count=2 -v ./internal/winsec/ ./internal/memory/` 我量到 **`=== RUN` 204 行 = 每轮 102 条**、`--- PASS` 120、`--- FAIL` 0、`--- SKIP` 2 行（1 个名字 ×2）⇒ **102 是"每轮/去重"数，不是 `-count=2 -v` 的原始行数**（他确实用了 `-v`，但两个口径不能混）；**"SKIP 恰 1 条"只对 winsec+memory 成立**，他 AC#4 命令里的 `./internal/risk/` 还有第 2 个 SKIP 名（`TestSyncRegistryProbeLive`），票面漏点名 ⇒ 记 **`R-103-5`（文档/口径，非代码）**。

### B. 静态门
```
gofmt -l internal/winsec/ internal/memory/                     → 空，rc=0
$(go env GOPATH)/bin/gofumpt.exe -l internal/winsec/ internal/memory/ → 空，rc=0（本机有二进制，跑了，不是"未跑"）
go vet ./internal/winsec/ ./internal/memory/                   → rc=0
GOOS=linux  go vet ./internal/winsec/ ; ./internal/memory/     → rc=0 / rc=0
GOOS=darwin go vet ./internal/winsec/ ; ./internal/memory/     → rc=0 / rc=0
GOOS=linux/darwin go build ./internal/winsec/ ./internal/memory/ → rc=0 / rc=0
```
按包量（**未**在仓根跑整树 `go vet ./...`／`go build ./...`：那在 Windows 当前树上必红，坏在别人的 `internal/config/parse.go:214 undefined: winsec` + `migrate.go "os" imported and not used`，不是我发现的缺陷）。

### C. `sh scripts/d22scan.sh`（与 CI 那一行逐字同形：`ci.yml:68 run: sh scripts/d22scan.sh`）
```
d22scan rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31
d22scan: examined 217 production Go files under internal/ and cmd/
bans #1-5 internal/=197  cmd/=20 | #6 frontend/=40 | #7 internal/tools/=17 | #8 design/=16 frontend/=40 internal/=346 cmd/=26
d22scan: clean - no D22 ban violations
```
**台账各 scope 文件数不降**：197/20/40/17/16/40 与实现方一致，`#8 internal/` 342→**346（升）**、`#8 cmd/` 26 持平（升的 4 个文件是别人新落的 `_test.go`）。未从仓根 `go run ./tools/d22scan`。

### D. 票 89 的那枚没被改松〔独立复现〕
`TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs` 在 MUT-2 批次与 winsec 全包（`ok 20.829s`）里都 **PASS**。

### E. "本机全绿 ≠ CI 过了"——CI 步级读数（run id 必报）
`gh api repos/CarlosShao/wisp/actions/runs?per_page=6`：
- **run `35587986855`（run_number 148，head `b2fa2ed`，created 2026-09-21T10:17:51Z，status=completed conclusion=failure）⇒ 这是第一个真跑完且包含 `0717bf2` 的 run**（`git merge-base --is-ancestor` 已证）。步级（`.../runs/35587986855/jobs`）：
  - `test-core` ⇒ **failure**，红在步 `Portable package tests (agent/llm/config/memory/observe/secret/risk/statemachine/...)`；但该步里 **`ok internal/memory 10.748s`、`ok internal/secret 0.006s`、`ok internal/risk 1.130s`**，唯一 `--- FAIL: TestPathCanonicalizerAccountsForRewrittenRoots`（`internal/tools 9.198s`）⇒ **票 102/106 地界，与 103 无关**。
  - `test-windows` ⇒ failure，`FAIL internal/secret 0.101s`，红因 `verifyPrivate: DACL names principals outside the private set ([LA LA])`（runner 账号 SID 映射为 `LA`，白名单只认 SY/BA/`u.Uid`）⇒ `winsec_windows.go` 的 `verifyPrivate` **不在 `0717bf2` 的 diff 里** ⇒ 既有环境问题，非 103 引入（但它挡住了"CI 过了"这句话）。
  - `lint` ⇒ failure 只在 `staticcheck` 步（`D22 positive control`、`D22 seven-ban`、`gofmt (gofumpt)`、`go vet (module)`、`go vet (tools/d22scan)` 全 success）。
  - `slo-smoke`/`slo-full`/`lint-frontend` success。
- 最新完成的是 **run `35589408388`（head `56fce50`，conclusion=failure）**，比 103 的码更新；**当前 HEAD（`4b4cfc9`＋本表）没有跑完过的 run** ⇒ **欠编排者：push 后复跑，并把 148 之外的步级读数补齐**。
- 本机全绿不外推：**`internal/risk` 与 `internal/memory` 在 CI ubuntu 上确实跑过并给出 ok（run 35587986855，步级 `Portable package tests`）**。

### F. AC#4 结论
**通过但有条件**：三门（gofmt/gofumpt/d22scan）+ 按包 vet（含 linux/darwin）+ 四包回归本代理全部复跑；一处假红（票 86 已知）与两个 SKIP 已逐条点名；"CI 过"这句**不成立**（148 是 failure，红因不属本票，且 HEAD 已前移）⇒ 条件＝**push 后复跑**（编排者的面）。

---

## 对实现方点名的三处复核要求，逐条直答

### ① `%VAR%`／前导 `~` 的**改写型**伪造仍归票 102 的 `Actable()`——边界守住了吗？
**判：这条边界本身守住了，但实现方对它的描述窄了一格。** 证据：
- 生产路径的改写型伪造**确实被 102 的账拦在 seam 之前**：`internal/risk/winsec_c26.go:39-44` ⇒ `Resolve(...)` 后走 `res.Actable()`，`internal/risk/pathresolver.go:93-100` ⇒ `Rewritten==true` 即 `ErrRewrittenPath`，错误经 `winsec/resolve.go:210` 的 `%w` 包成密封失败 ⇒ **winsec 不会把"被 `%VAR%`/`~` 改写到别人树"当作合法输入继续封**。静态那把锁也在：`internal/risk/pathresolver_rewrite_account_test.go:174` 钉住 `winsec_c26.go` 必须读 `Actable(`。
- 但"winsec 看不出善意改写与劫持改写的区别"这句话**不止适用于 `%VAR%`/`~`**：P1b 证明，**任何** resolver（含伪造者）交回的"干净但换了树"的答案，winsec 一律照单封（`resolve.go:222` 的再验证只认拼写形状，不认"答案是否还在调用者命名的树里"）。⇒ 归票 102 的是**入口展开**那一半；**缝上答案的树归属**那一半没有票接着 ⇒ 我把它记进 `R-103-1` 的 remediation（需要自己的票，不属 102 也不属 103 现判据）。
- 引证 file:line：`internal/winsec/resolve.go:198-227`（mint + 再验证）、`internal/winsec/winsec.go:231-237`（`RemoveUnlinked` 故意不走 `ResolvePath`，靠操作本身 + 祖先检查）、`internal/risk/winsec_c26.go:39-44`。

### ② POSIX 侧只能编译期验——到底验没验到？
**验到了，而且不必只靠 CI：本代理用真 Linux 容器把它量了出来。**
```
docker run --rm -e CGO_ENABLED=0 -v <snapshot>:/src -w /src golang:1.27 go test -count=1 -v ./internal/winsec/ ./internal/risk/ ./internal/memory/ ./internal/secret/
INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1        ← Linux 上探针集合是 1 个形状
--- PASS: TestAC103POSIXRiskResolverPassesTheProbe   ← **探针没有把 risk 的纯词法解析器拒掉**
ERROR winsec: refusing to install ... resolver=winsec_test.posixStamp reason="it answered the hostile shape \"/tmp/../wisp-103-conformance-probe\" by passing it through unchanged"
--- PASS: TestAC103POSIXRemoveUnlinkedAncestorGuard  ← POSIX 祖先守卫拒了 link/keep-me.txt 与 ../ 变体，外国文件在位
ok internal/winsec 0.065s / ok internal/risk 4.588s / ok internal/memory 12.794s / ok internal/secret 0.017s（RC=0，去变异后的干净复跑）
（注：首轮 memory 报 [build failed] 是**我自己 MUT-3 探针残留** `undefined: descendFollowingLinks103`，已还原后重跑为 ok —— 如实登记，不是仓库缺陷。）
```
CI 侧对得上：**run `35587986855` 步级** `ok internal/risk 1.130s`、`ok internal/memory 10.748s`（ubuntu）。
⇒ **②判：不再是未验面**，Linux 上探针接受 risk 的解析器（因为它把 `/tmp/../x` 词法改写成 `/x`，底线认这个改写）。**残余缺口**只剩 darwin：本代理容器只量了 linux，`GOOS=darwin` 只到"编译 + vet rc=0"；且探针在 Linux 只有 1 个形状（相对路径那条只在 Windows 加，`resolve.go:143-145`）⇒ 记 **`R-103-6`**（darwin 无运行期读数 + Linux 探针少一个形状）。

### ③ 祖先检查 fail-closed ⇒ 契约从严还是可用性回归？
**判：方向正确的契约从严，但"正常安装路径不受影响"这一半没有用例证明。** 证据与推理：
- 从严是**契约要求的**：票面 AC#2 原句就是"必须拒、返回错误且不删目标"，且 `winsec_other.go:84-86` 写明了 POSIX 为何也要查祖先。
- **正常路径不受影响的正面读数**：`TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink` PASS（叶子是链接 ⇒ 照旧解链）+ 真实回归里 `internal/memory`（`ok 30.163s`，含票 18/79 的 LRU/回收用例）与 `internal/secret`（`ok 1.400s`）全绿 + **Linux 容器里 `internal/memory` 全绿（`ok 12.794s`）** ⇒ 普通数据根（自己 `t.TempDir()` 底下、无链接祖先）的回收没被打断。
- **但**没有一条用例钉住"数据根**真的**挂在别人 symlink 底下 ⇒ 报的是拒删错误而不是崩溃/静默"这一形态；而 `winsec_other.go:57-64` 自己就点过同形风险（"会打断 `/tmp` 或 HOME 在 symlink 后面的 macOS 安装"）——祖先检查在 POSIX 用 `ModeSymlink`，macOS 上 `/tmp`、`/var` 本身就是 symlink ⇒ **在 macOS 上以 tmp 为数据根的回收会开始报错**。方向仍是拒不是删（安全），但对票 79 的回收路由是语义收紧。⇒ 记 **`R-103-7`**（补一条 macOS/POSIX 形态用例，说明"报错可接受/需豁免"由编排者裁）。

---

## 缺陷登记（验收期一律不修，改它该有自己的票）

- **`R-103-1`（安全边界，AC#1）** 一次性守卫可被**注册顺序**解除：`SetPathResolver(nil)` 无条件允许（`resolve.go:97-104`），释放缝后恒改写型伪造过一致性探针（`resolve.go:157-162`）⇒ `SealFile` 把自己树上的请求改到别人树上并返回 nil；实测 `icacls` 前后 `S-1-1-0` 由有到无（§AC1-E P1b）。**修法方向**（不代为实现）：探针除"拒/底线认的改写"外还须要求**答案与输入的树归属关系**（或禁止对空缝安装非 risk 的解析器、或 `nil` 重置只在 init 期允许）。
- **`R-103-2`（守卫可绕，AC#2）** `firstLinkAncestor` 只按 `filepath.Separator` 切段（`winsec.go:250-259`）⇒ Windows 上 `/` 或混合分隔符拼写一个祖先都不查，`RemoveUnlinked` 沿 junction 删掉外来文件且返回 nil（§AC2-D，两腿均复现）。今天 `removeStray` 用 `filepath.Join` ⇒ 陷阱档。
- **`R-103-3`（既有洞，非 103 引入）** `platformVerifyPlacement`（`placement_windows.go:52`）同样只按 `\` 走查 reparse ⇒ 纯底线（未链接 risk）二进制里 `SealFile("C:/.../link/sub/keep-me.txt")` 返回 nil 且真改掉外来 DACL（§AC2-E）。**该文件不在 `0717bf2` diff 内**；票 94 的底线的账，需自开票。
- **`R-103-4`（钉语义，AC#1）** 恒等重装分支（`resolve.go:123-124`）覆盖率 **0** ⇒ "静默保留第一份"这一幂等语义没有任何用例钉住。
- **`R-103-5`（口径/文档）** 实现方 AC#4 的"`=== RUN` 102、SKIP 恰 1"：102 是每轮数（`-count=2 -v` 原始 204 行），且"SKIP 恰 1"不含其命令里 `./internal/risk/` 的 `TestSyncRegistryProbeLive`。另：票 93 刚落地的口径是"SKIP 不算 ok"（`c8c828e`）⇒ 台账类计数口径今后应以"每轮/原始行"二选一词写明。
- **`R-103-6`（POSIX 残余）** darwin 只有编译期读数；Linux 探针集合只有 1 个形状（相对路径那条 Windows-only，`resolve.go:143-145`）。补救：CI 或本地容器补 darwin 运行期（macOS runner）+ 考虑 Linux 也加相对形状（需先量 `risk` 会不会因此被误拒——本代理 Linux 读数是 `probes_passed=1` 接受）。
- **`R-103-7`（可用性形态，AC#2 ③）** 无"数据根位于 symlink 底下 ⇒ 报 ErrIsReparsePoint 而非崩溃/静默"的用例；POSIX（尤其 macOS `/tmp`、`/var`）下回收会开始报错，方向为拒删，需裁决是"从严即可"还是"要豁免通道"。

---

## 工具输出里的非授权文本登记
本轮所有工具输出末尾**未出现**自称"编排者备注/停手/冻结某包"的文本；收到的两条后台任务通知为系统事件（含 "NOT user input" 标记），按规矩不作为授权处置，未据此改变任何判据。（台账 A75②）

## 我没有改的东西
`docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`（含 `pathresolver*.go`、`rules_gateway.go`）、`tools/d22scan/**`、`allowlist.txt`、`internal/memory` 生产码、票面 Status 与文件名 —— 全部零改动；本表唯一写入是 `docs/evidence/s1/103-adversarial-acceptance.md` 与票面文末追加段。
