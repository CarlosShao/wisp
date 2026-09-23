# 票 124 AC#2a — 软链形 132 枚红名的「断言方向账」（逐枚：可转 / 拒绝腿 / 待裁）

**清点方**：`worker-ticket124-ac2a` · **日期**：2026-09-23
**本格性质**：纯清点。**零 `.go` 改动**（判据之一）；只产这一份账，供 AC#2b 决定能动哪几枚、必须留红哪几枚。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2 的裁定」AC#2a 段（15:3x 编排者）。

## 0. 锚点 / 快照 / 跑法（可复核）

- 锚定 sha：`bcb03aa00c260d328048616acf6b59ffcad912ee`（开工 `git rev-parse HEAD`）。
  ⚠ 记一笔：任务简报让我「开工先记 HEAD」——**第一次 `git rev-parse HEAD` = `6fdb39d`，两秒后再读已变 `bcb03aa`**（他人正往共享树 `dev` 提交）。本账锁定 `bcb03aa` 做全部读数。
  AC#1 的锚点是 `7b4c36a`。两锚之间 132 枚包里**只有 `cmd/wisp/leg_sink_gate_131_test.go` 等票 131 的码改过**（tools/memory/winsec/agent/config/llm/perm 的测试源码两锚逐字同），⇒ 见 §6 对账。
- 纯净快照：`git archive bcb03aa | tar -x -C /d/tmp/wisp124-ac2a`（1033 枚文件；容器内 `ls -l /src/go.mod` = `-rwxrwxrwx 1 root root 883`，`md5sum resolve.go = b6876a5efe759f6e17434d1b50a129c3`，与 AC#1 同字 ⇒ 挂载真生效）。
- 容器：`golang:1.27`（`go1.27.1 linux/amd64`）；复用票 119 命名卷 `ac119-gomodcache`/`ac119-gocache` ⇒ 离线可编。
- 跑法同 AC#1（`/d/...` + `MSYS_NO_PATHCONV=1`；进容器先 `ls -l /src/go.mod`；形状硬断言 `exit 97`/`exit 98` 沿用）。两枚样本各跑**受影响 9 包** `-count=1 -v`：`A2L`＝软链形（`ln -s /realpriv /varlink`+`TMPDIR=/varlink/w124tmp`），`A2P`＝普通形（`TMPDIR=/plainroot/w124tmp`）。
- 开测前查在飞 run（本机 self-hosted runner）：`gh run list` 见 run `35831816899` `in_progress`（一枚 docs push 的 ci）⇒ 与我的跑同机抢 CPU；**只影响单枚耗时，不影响红绿**（普通形/软链形红名差集与是否争用无关，两形同容器同码只差 `TMPDIR` 一个变量）。

## 1. 三档总数

| 档 | 枚数 | 说明 |
|---|---|---|
| **可转**（接上「先解析再递底线」后断言照样成立） | **131** | 含 setup 被未解析根挡住（断言要「成」）与「要的是别的拒绝、普通形已满足」两小类；见逐枚判据 |
| **拒绝腿**（要的就是「未解析根必须被拒」这句 ⇒ AC#2b 不许变绿） | **0** | 一枚都没有——与 AC#2裁定那句硬前提一致：132 枚普通形全绿 ⇒ 无一能自证是拒绝腿 |
| **待裁 / 归因待票 123 裁** | **1** | 仅 `TestL1WriteGoesThroughTheRealBlockWindow`（软链形 ~300.04 s 才 FAIL，300 s 恰为 C18 审批超时常量，不像本票「未解析根」这一族；按票面 AC#2裁定 附带发现 2 单列、不替它下结论） |

**合计 131 + 0 + 1 = 132**，与 AC#1 §5.3 那 132 枚逐名对齐（对账见 §6）。

## 2. 按包分布表

| 包 | 132 里枚数 | 可转 | 拒绝腿 | 待裁/归因 |
|---|---|---|---|---|
| `./internal/agent/` | 26 | 26 | 0 | 0 |
| `./internal/agent/approval/` | 1 | 1 | 0 | 0 |
| `./internal/config/` | 7 | 7 | 0 | 0 |
| `./internal/llm/` | 17 | 17 | 0 | 0 |
| `./internal/memory/` | 33 | 33 | 0 | 0 |
| `./internal/perm/` | 5 | 5 | 0 | 0 |
| `./internal/tools/` | 21 | 20 | 0 | 1（`TestL1Write`，归因待 123） |
| `./internal/winsec/` | 17 | 17 | 0 | 0 |
| `./cmd/wisp/` | 5 | 5 | 0 | 0 |
| **合计** | **132** | **131** | **0** | **1** |

## 3. 「因为别的原因红」的点名（不是未解析根那一族）

- **在 132 之内**：只有 `./internal/tools/TestL1WriteGoesThroughTheRealBlockWindow`——软链形 `--- FAIL (300.04s)`、普通形 `--- PASS (3.01s)`。300 s 这个数是 C18 审批超时常量，红因是「进审批窗口、没人在旁边点确认 ⇒ 300 s 超时」，不是「未解析根被拒」。**归因待票 123 裁**（本票不替它下结论）。
- **在 132 之外、但本形新显形**：`./internal/tools/TestLateVetoRendersTheApprovalLayersAppliedStepsReport`——软链形 `--- FAIL (300.02s)`、普通形 `--- PASS (3.00s)`，与上枚同病（C18 审批超时族）。它**不在 AC#1 的 132 清单里**，是本快照新量出的一枚 only-in-link（见 §6 多一枚的解释）。同样应归口票 123，不进本票 AC#2b 的转换范围。
- 其余 130 枚的红因都是「未解析根」，只是**表层文案分两种**：(a) 直接被密封底线拒（`winsec: refusing to seal … path is not provably resolved … through the link at /varlink`）；(b) 先被 `risk.c26Pipeline` / 路径规范化器判为「那棵树不在盘上、不授权」⇒ 写升级 L2 ⇒ `L2 审批未通过，已拒绝执行`。**两者根因同为未解析根**（普通形绿、软链形红，且把根解析后各自回到普通形结局），归 AC#2b 可转；差别只在文案，逐枚判据里已写明。

## 4. 证据来源标注（哪些跑出来、哪些读码定）

- **跑出来**（`A2L`/`A2P` 两形 `go test -v` 原始日志，`/d/tmp/wisp124-ac2a-logs/`）：每枚「软链形红 / 普通形绿」的事实、每枚红**拿到的拒绝字符串**（col ②的落点与 §3 的「别的原因」判定）——即「它是未解析根那一族，还是超时/别因」这一刀，全部由跑出来的失败文案定，不靠读码猜。
- **读码定**（快照 `bcb03aa` 的 `_test.go`）：col ②的**断言原文摘录**（源码里的 `t.Fatalf(...)` / `if …` 那句）、以及「该断言到底要『成』还是要『拒』、解析后成不成立」的语义判断。
- 二者互校：对断言含「拒/refuse/Err…」字样的每一枚（winsec AC5/AC118/AC3 反向腿、memory 迁移腿、tools 回收站腿、llm `RefusesSilentRuns` 等），都单独回看了它在软链形拿到的字符串，确认那句「拒」**不是**「未解析根必须被拒」，而是普通形本就拿到的另一种拒（注入 / leaf-link / 回收站 / `ErrSchemaUnmigratable`）⇒ 解析根后仍成立 ⇒ 判可转。**没有任何一枚的断言是「就是要拿未解析的根去试底线拒绝」。**

## 5. 逐枚明细（132 行，不合并）

### 5.1 `./internal/winsec/`（17 枚：13 顶层 + `TestAC5FailedSealRefusesTheWrite` 的 4 子测试）— 全可转

> 这一包是密封底线自己，最像「拒绝腿」，故先出、逐字读源码 + 回看软链形字符串。结论：全可转——它们要的是「自己的树必须封得住」或「另一种拒（注入 / leaf-link）」，普通形已满足，解析根后不洗绿。

- **1 `private_fail_test.go:26` `TestAC5FailedSealRefusesTheWrite`**
  摘录：`dir := t.TempDir()` → 子测试断言 `if !errors.Is(err, ErrNotSealable) { t.Fatalf("error does not name the refusal: %v", err) }`（`private_fail_test.go:37-38`）
  判定：**可转**。判据：软链形它拿到的是 `error does not name the refusal: winsec: refusing to seal … path is not provably resolved … /varlink`（A2L 逐字）——它要的是**注入的** `applyDescriptor` 失败（`ErrNotSealable`），未解析根只是**抢在前面**把返回换成另一个 sentinel。解析根后注入腿重新成为操作性拒因 ⇒ 回到普通形绿，注入断言不被洗。

- **2 `private_fail_test.go`（子测试）`TestAC5FailedSealRefusesTheWrite/exclusive_artifact`**
  摘录：`err := PrivateFileExclusive(p, …); if err == nil { t.Fatal("the write succeeded although the seal was refused") }; if !errors.Is(err, ErrNotSealable) { … }`（`private_fail_test.go:33-38`）
  判定：**可转**。判据：同上——A2L 该子测试红在 `:38 error does not name the refusal: … not provably resolved`，要的是 `ErrNotSealable`（注入），非未解析根。解析后成立。

- **3 `TestAC5FailedSealRefusesTheWrite/replacement_write`**
  摘录：`err := PrivateFile(p, []byte("dpapi-shaped secret"), 0o600); if err == nil { t.Fatal("PrivateFile accepted a refused seal") }; if !errors.Is(err, ErrNotSealable) { t.Fatalf("error does not name the refusal: %v", err) }`（`private_fail_test.go:46-51`）
  判定：**可转**。判据：A2L 红在 `:51 error does not name the refusal: … not provably resolved`，同注入腿。解析后成立。

- **4 `TestAC5FailedSealRefusesTheWrite/directory_chain`**
  摘录：`if err := PrivateDirAll(p, 0o700); err == nil { t.Fatal("PrivateDirAll reported a private directory it could not seal") } else if !errors.Is(err, ErrNotSealable) { … }`（`private_fail_test.go:59-62`）
  判定：**可转**。判据：A2L 红在 `:62 … not provably resolved`，要注入的 `ErrNotSealable`。解析后成立。

- **5 `TestAC5FailedSealRefusesTheWrite/seal_file_and_dir_direct`**
  摘录：`if err := SealFile(p); !errors.Is(err, ErrNotSealable) { t.Fatalf("SealFile: %v", err) }; if err := SealDir(dir); !errors.Is(err, ErrNotSealable) { … }`（`private_fail_test.go:72-76`）
  判定：**可转**。判据：A2L 红在 `:73 SealFile: … not provably resolved`。解析根后 `os.WriteFile` 落盘、注入的 `SealFile/SealDir` 失败成为唯一拒因 ⇒ 成立。

- **6 `private_fail_test.go:99` `TestAC5FailureIsNotSwallowedByTheHappyPath`**
  摘录：`if err := PrivateFileExclusive(p, []byte("written while the seal works")); err != nil { t.Fatalf("PrivateFileExclusive: %v", err) }`（`private_fail_test.go:102-104`）
  判定：**可转**。判据：断言要「**成**」（happy path 必须写成功）。A2L 红在 `:103 PrivateFileExclusive: … not provably resolved`。解析后写成功 ⇒ 绿；反向腿，与拒无关。

- **7 `private_other_test.go:20` `TestPOSIXPrivateFileIsReally0600`**
  摘录：`if err := PrivateFile(p, []byte("tool output"), 0o644); err != nil { t.Fatalf("PrivateFile: %v", err) }` 后 `if got := info.Mode().Perm(); got != 0o600 { t.Errorf("artifact mode is %v, want -rw-------", got) }`（`private_other_test.go:25-34`）
  判定：**可转**。判据：要「成 + 权限位 0600」。A2L 红在 `:26 PrivateFile: … not provably resolved`（没走到权限读回）。解析后成立。

- **8 `private_other_test.go:37` `TestPOSIXPrivateDirIsReally0700`**
  摘录：`if err := PrivateDirAll(root, 0o755); err != nil { t.Fatalf("PrivateDirAll: %v", err) }`（`private_other_test.go:39-41`）
  判定：**可转**。判据：要「成 + 0700」。A2L 红在 `:40 PrivateDirAll: … not provably resolved`。解析后成立。

- **9 `private_other_test.go:60` `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed`**
  摘录：`if err := PrivateDirAll(root, 0o700); err != nil { t.Fatal(err) }` … 核心 `if _, err := os.Stat(innocent); err != nil { t.Errorf("TARGET DELETED - removal followed the symlink: %v", err) }`（`private_other_test.go:72-74`、`:90-92`）
  判定：**可转**。判据：它测的是 `RemoveUnlinked` **不跟着软链删目标**（另一件事）；软链形红在 `:73` 的 setup `PrivateDirAll … not provably resolved`，压根没跑到链接腿。解析后 setup 成、链接腿照常测 ⇒ 目标不被删的断言成立。

- **10 `private_other_test.go:97` `TestPOSIXMissingFileIsNotAnError`**
  摘录：`if err := RemoveUnlinked(filepath.Join(t.TempDir(), "gone")); err != nil { t.Errorf("RemoveUnlinked of a non-existent entry: %v", err) }`
  判定：**可转**。判据：要「**不报错**」。软链形里 `RemoveUnlinked` 对未解析 `/varlink/...` 返回未解析根错误 ⇒ 断言被反触发。解析后「缺文件不算错」的幂等承诺成立。

- **11 `ancestor_separator_108_other_test.go:94` `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`**
  摘录：`if err := winsec.RemoveUnlinked(stray); err != nil { t.Errorf("AC#2 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file: %v", stray, err) }`（`:111-113`）
  判定：**可转**。判据：反向腿——要「自己的树里的文件必须能删」，不许因反斜杠被折成隔符而误拒。软链形红因未解析根又叠了层拒。解析后成立（普通形已绿）。

- **12 `ancestor_separator_108_other_test.go:152` `TestAC4POSIXFloorAnswersInsideTheNamedTree`**
  摘录：`if err := winsec.PrivateDirAll(dir, 0o700); err != nil { t.Fatalf("PrivateDirAll(%s): %v", dir, err) }`（`:158-160`）；核心 `if got.String() != dir { t.Errorf("AC#4 RED: the floor answered %q for %q …") }`
  判定：**可转**。判据：要「成 + 底线答案留在命名树内」。A2L 红在 setup `PrivateDirAll … not provably resolved`。解析后成立。

- **13 `placement_leaf_118_other_test.go:75` `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`**
  摘录：helper `leafLinkTo118` 的 `if same, err := filepath.EvalSymlinks(link); err != nil || same != target { t.Fatalf("the leaf link does not answer with its own target: %q vs %q (err=%v)", same, target, err) }`（`placement_leaf_118_other_test.go:50-52`）；测试体 `assertRefused113(t, "SealFile", link, winsec.SealFile(link))`（`:81`）
  判定：**可转**。判据：它**确实**断言一个拒——但拒的是「leaf 位置那条指向外人的软链」（`assertRefused113` 要 `ErrUnresolvedPath`）。软链形它红在 helper 前置 `:78 the leaf link does not answer with its own target: "/realpriv/…" vs "<T>"`（`EvalSymlinks` 穿 `/varlink` 解析出 `/realpriv` 而 `target` 是未解析拼写），根本没跑到 `SealFile`。解析根后 leaf-link 腿照旧被拒（普通形已证）⇒ 它要的那句拒**不洗绿**。可转。

- **14 `placement_leaf_118_other_test.go:96` `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`**
  摘录：同 13 的 helper 前置；测试体 `err := winsec.PrivateFile(link, []byte("top secret"), 0o600); assertRefused113(t, "PrivateFile", link, err)`（`placement_leaf_118_other_test.go:106-107`）
  判定：**可转**。判据：A2L 红在 `:99 the leaf link does not answer with its own target`（同 13，前置挡住）。它要的拒是 leaf-link 的 `PrivateFile` 拒（普通形满足），非未解析根。可转。

- **15 `placement_symlink_113_other_test.go:253` `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`**
  摘录：`if err := winsec.SealFile(plain); err != nil { t.Errorf("AC#3 RED: a plain file inside the named tree was refused: %v", err) }`（`:268-270`）
  判定：**可转**。判据：反向腿——命名树内的普通文件**必须能封**（不许一律拒）。软链形被未解析根误拒 ⇒ 反触发。解析后成立。

- **16 `placement_symlink_113_other_test.go:287` `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`**
  摘录：`if err := winsec.PrivateDirAll(deep, 0o700); err != nil { t.Fatalf("AC#3 RED: PrivateDirAll refused its own tree: %v", err) }`（`:297-299`）
  判定：**可转**。判据：反向腿（sibling 链/子链不该过拒）。A2L 红在 setup 未解析根。解析后成立。

- **17 `placement_symlink_113_other_test.go:335` `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`**
  摘录：`if err := winsec.SealFile(mine); err != nil { t.Errorf("AC#3 RED: POSIX folded the backslash in %q into a separator and refused the tree's own file: %v", mine, err) }`（`:354-356`）
  判定：**可转**。判据：反向腿。A2L 红在 `:355 … refused the tree's own file: … not provably resolved`。解析后成立。

### 5.2 `./internal/agent/`（26 枚：18 顶层 + 8 子测试）— 全可转

> 这 26 枚的软链形失败串在 A2L 里**只有两种**：`winsec: refusing to seal … the installed risk.c26Pipeline answered …` 与 `… path is not provably resolved … through the link at /varlink`（无第三因）。每一枚的红都落在**开库/开 spill 存储的 setup**（`memory.Open` / `NewSpiller`+`Prepare`）上——底线把未解析的 `t.TempDir()` 拒了 ⇒ 库开不成 ⇒ 测试真正要断言的行为根本没跑起来。断言原文全是「必须成」形 ⇒ 全可转。

- **18 `forensics_test.go:22` `TestFailedTaskBooksOpenCallRowWithDecision`**
  摘录：`t.Fatalf("memory.Open: %v", err)`（`forensics_test.go:27`）
  判定：**可转**。判据：setup 开库被未解析根拒（A2L 串含 `not provably resolved`）；测的是「失败任务落 open-call 行」，要成。解析根后开库成 ⇒ 断言成立。
- **19 `forensics_test.go:82` `TestCancelledTaskPersistsTerminalRows`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:87`）。判：同上（开库 setup 被拒），要成 ⇒ 可转。
- **20 `guard_test.go:151` `TestPerToolTimeoutFires`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:155`）。判：开库 setup 被拒；测工具超时行为，要成 ⇒ 可转。
- **21 `guard_test.go:211` `TestPerToolTimeoutOfContractHonestToolIsToolClass`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:215`）。判：同上 ⇒ 可转。
- **22 `loop_golden_test.go:299` `TestTaskLogAndToolCallRows`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:303`）。判：开库 setup 被拒，要成 ⇒ 可转。
- **23 `spill_name_injectivity_test.go:99` `TestSpilledBytesSurviveANameThatUsedToCollide`（顶层）**
  摘：`dir := t.TempDir()` → `sp := NewSpiller(dir, spill79Budget)` → 子测试 `t.Fatalf("Prepare(%q): %v", pair.idA, err)`（`:114`）。判：顶层红由子测试聚合；`Prepare` 落 spill 到未解析根被拒，要成 ⇒ 可转。
- **24 `spill_name_injectivity_test.go:277` `TestWriteFileExclusiveIsExclusive`**
  摘：`t.Fatalf("first exclusive write failed: %v", err)`（`:282`）。判：exclusive 写进未解析根被拒；断言要「写成且二次冲突」，要成 ⇒ 可转。
- **25 `spill_name_injectivity_test.go:312` `TestSpillSameIDRetryOverwrites`**
  摘：`t.Fatalf("first spill: spilled=%v err=%v", a.Spilled, err)`（`:318`）。判：spill setup 被未解析根拒，要成 ⇒ 可转。
- **26 `spill_name_injectivity_test.go:352` `TestSpillAcrossRestartsKeepsRetrySemantics`**
  摘：`t.Fatal(err)`（`:355`，spill/重开 setup）。判：setup 被拒，要成 ⇒ 可转。
- **27 `spill_path_invariant_test.go:102` `TestSpillCallIDHostileShapesSanitizedToBareNames`（顶层）**
  摘：`root := t.TempDir()` → 子测试 `t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)`（`:112`）。判：顶层红由子测试聚合；`Prepare` 写 spill 被未解析根拒。⚠ 名字带「Hostile Shapes」但那是 **callID 文件名注入形状**（`separator`/`dotdot`/`drive_letter`/`unc`），与路径未解析根无关，断言要「净化成裸名且写成功」⇒ 可转。
- **28 `spill_path_invariant_test.go:199` `TestSpillContainmentByDirectoryListing`**
  摘：`t.Error(e)`（`:241`，遍历 spill 目录的不变式）。判：spill 根未解析 ⇒ 落盘/列目录不变式破；要成 ⇒ 可转。
- **29 `spill_path_invariant_test.go:344` `TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`**
  摘：`t.Fatalf("memory.Open: %v", err)`（`:349`）。判：开库 setup 被拒，要成 ⇒ 可转。
- **30 `spill_test.go:24` `TestSpillThresholdScalesWithWindow`**
  摘：`t.Fatalf("4096 window: spilled=%v err=%v, want spill (threshold shrank past the payload)", sp.Spilled, err)`（`:62`）。判：要「必须 spill 成功」；未解析根使 spill 落不下 ⇒ 可转。
- **31 `spill_test.go:71` `TestSpillTokenBoundary`** — 摘：`t.Fatal(err)`（`:88`，spill setup）。判：setup 被拒，要成 ⇒ 可转。
- **32 `spill_test.go:97` `TestSpillArtifactAndStubShape`** — 摘：`t.Fatal(err)`（`:111`）。判：setup 被拒，要成 ⇒ 可转。
- **33 `spill_test.go:186` `TestSpillThroughLoop`** — 摘：`t.Fatalf("tiny window must spill the same result, got %+v", resTiny.ToolLog)`（`:214`）。判：要「loop 里 tiny window 也 spill 成」；未解析根挡落盘 ⇒ 可转。
- **34 `spill_test.go:241` `TestSpillArtifactRespectsRawCap`** — 摘：`t.Fatal(err)`（`:249`）。判：setup 被拒，要成 ⇒ 可转。
- **35 `truncation_test.go:26` `TestMaxTokensFailsAllToolCallsOfThatMessage`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:26`）。判：开库 setup 被拒，要成 ⇒ 可转。
- **36 子测试 `TestSpilledBytesSurviveANameThatUsedToCollide/slash_id_then_bare_id`**（父 `spill_name_injectivity_test.go:99`）
  摘：`t.Fatalf("Prepare(%q): %v", pair.idA, err)`（`spill_name_injectivity_test.go:114`）。判：`Prepare` 落 spill 被未解析根拒；断言要「碰撞名前后 spill 都存活」，要成 ⇒ 可转。
- **37 子测试 `TestSpilledBytesSurviveANameThatUsedToCollide/bare_id_then_backslash_id`**（同父/同 assert）
  摘：`t.Fatalf("Prepare(%q): %v", pair.idA, err)`（`spill_name_injectivity_test.go:114`）。判：同上，未解析根挡 spill，要成 ⇒ 可转。
- **38 子测试 `TestSpilledBytesSurviveANameThatUsedToCollide/dotdot_id_then_word_id`**（同父/同 assert）
  摘：`t.Fatalf("Prepare(%q): %v", pair.idA, err)`（`spill_name_injectivity_test.go:114`）。判：同上 ⇒ 可转。
- **39 子测试 `TestSpillCallIDHostileShapesSanitizedToBareNames/separator`**（父 `spill_path_invariant_test.go:102`）
  摘：`t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)`（`spill_path_invariant_test.go:112`）。判：`Prepare` 被未解析根拒；形状是 callID 名净化、非路径根 ⇒ 要成，可转。
- **40 子测试 `TestSpillCallIDHostileShapesSanitizedToBareNames/dotdot`**（同父/同 assert）
  摘：`t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)`（`spill_path_invariant_test.go:112`）。判：同上 ⇒ 可转。
- **41 子测试 `TestSpillCallIDHostileShapesSanitizedToBareNames/drive_letter`**（同父/同 assert）
  摘：`t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)`（`spill_path_invariant_test.go:112`）。判：同上 ⇒ 可转。
- **42 子测试 `TestSpillCallIDHostileShapesSanitizedToBareNames/unc`**（同父/同 assert）
  摘：`t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)`（`spill_path_invariant_test.go:112`）。判：同上 ⇒ 可转。
- **43 子测试 `TestSpillCallIDHostileShapesSanitizedToBareNames/encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence`**
  摘：`t.Fatalf("Prepare(%q): %v", id, err)`（`spill_path_invariant_test.go:166`）。判：同上，setup 被未解析根拒，要成 ⇒ 可转。

### 5.3 `./internal/agent/approval/`（1 枚）— 可转

- **44 `batch_test.go` `TestTenOpsInOneToolCallGetOneConfirm`**
  摘录：`p := ui.wait(t); if p.Level != "L1" { t.Fatalf("level=%q：授权目录内的声明 L1 写应留在 L1", p.Level) }`（`batch_test.go:164-167`）
  判定：**可转**。判据：软链形它红在 `level="L2"`（A2L `:166`）——写目标在未解析 `t.TempDir()` 里 ⇒ 规范化器不认它「在授权树内」⇒ 声明 L1 被升级到 L2。随后批合并拿不到、`spawn` 的协程卡在待答 ⇒ cleanup 里 `fakes_test.go:193 派生的测试协程未在 3s 内退出` + `:277 闸门退出时仍有 1 个待审批项未清理`。**3 s 是 t.Fatalf 之后 cleanup 的连锁，非独立病因**；根因仍是未解析根→不授权→升级 L2。解析根后目录被认作授权、写留 L1、10 项并一次确认 ⇒ 回到普通形绿。（此枚在 AC#1 两枚样本 L1/L2 同现，非争用 flake。）

### 5.4 `./internal/config/`（7 枚：6 顶层 + 1 子测试）— 全可转

> 逐字：A2L 里 config 的红串是 `save: config: config.toml temp file seal: winsec: refusing to seal /varlink/…`（写 config 到未解析根被拒），断言全是「必须 save/load/migrate 成」⇒ 全可转。⚠ 注意本包 §5.4 附带发现：`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125` 在软链形从「跑」变「SKIP」——它**不在本 132 里**（不属这 7 枚），是票 125 自己认出形状拒测的正向钉，别把它当「被清掉的红」。

- **45 `boundary_test.go:186` `TestRoundTripLoadMarshalLoad`** — 摘：`t.Fatalf("save: %v", err)`（`:191`；实拿 `…config.toml temp file seal: refusing to seal /varlink/…`）。判：save 写未解析根被拒，要成 ⇒ 可转。
- **46 `boundary_test.go:234` `TestResolvedNeverPersists`** — 摘：`t.Fatal(err)`（`:239`，save 步）。判：setup 写被拒，要成 ⇒ 可转。
- **47 `loader_test.go:275` `TestSaveFileRoundTripsThroughLoad`** — 摘：`t.Fatalf("save: %v", err)`（`:292`）。判：同上 ⇒ 可转。
- **48 `migrate_test.go:48` `TestMigrateV1Fixture`** — 摘：`t.Fatalf("migrating load: %v", err)`（`:53`）。判：迁移前要写/开未解析根被拒，要成 ⇒ 可转。
- **49 `migrate_test.go:112` `TestMigrateNoVersionKeyAssumedV1`** — 摘：`t.Fatalf("file without schema_version is treated as v1: %v", err)`（`:116`）。判：写 seed 文件被未解析根拒 ⇒ 可转。
- **50 `unwired_test.go:55` `TestUnwiredGuardLeavesHonestConfigsAlone`（顶层）** — 摘：子测试 `t.Fatal(err)`（`unwired_test.go:90`）。判：顶层红由子测试聚合；写/重载 config 被未解析根拒，要成 ⇒ 可转。
- **51 子测试 `TestUnwiredGuardLeavesHonestConfigsAlone/SaveFile_output_reloads`** — 摘：`t.Fatal(err)`（`unwired_test.go:90`，SaveFile 产物重载）。判：SaveFile 落未解析根被拒 ⇒ 可转。

### 5.5 `./internal/llm/`（17 枚：11 顶层 + 6 子测试）— 全可转

> 11 枚顶层红全部落在 `pf := newProbeFixture(t)`——fixture 内 `memory: create data dir: winsec: refusing to seal /varlink/… not provably resolved`。测的是 probe 套件对能力的判读（要成），与拒无关。⚠ `TestProbeSuiteRefusesSilentRuns` 名字带 "Refuses"，但它红在**同一 fixture setup 被拒**（`probe_health_test.go:391`），它要「拒」的是「probe 不许静默跑」这一应用行为、不是未解析根 ⇒ 可转（解析后 fixture 开得成，那条"拒静默"照旧成立）。

- **52 `probe_health_test.go:150` `TestProbeSuiteMeasuresBrokenFC`** — 摘：`pf := newProbeFixture(t)`（`:151`）。判：fixture 开库被未解析根拒，测能力判读要成 ⇒ 可转。
- **53 `probe_health_test.go:213` `TestProbeSuiteMeasuresBrokenVision`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **54 `probe_health_test.go:251` `TestProbeSuiteHonestProviderRecordsNoMismatch`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **55 `probe_health_test.go:281` `TestProbeSuiteHonestNegativeIsNotAMismatch`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **56 `probe_health_test.go:306` `TestProbeSuiteBrokenOnAllThreeDialects`（顶层）** — 摘：子测试 `pf := newProbeFixture(t)`（`:309`）。判：顶层红由 3 方言子测试聚合 ⇒ 可转（见 59–61）。
- **57 `probe_health_test.go:352` `TestProbeSuiteCapableOnAllThreeDialects`（顶层）** — 摘：子测试 `pf := newProbeFixture(t)`（`:355`）。判：同上 ⇒ 可转（见 62–64）。
- **58 `probe_health_test.go:390` `TestProbeSuiteRefusesSilentRuns`** — 摘：`pf := newProbeFixture(t)`（`:391`）。判：红在 fixture setup 被未解析根拒；"Refuses Silent"是 probe 行为断言、非根，要成 ⇒ 可转。
- **59 `probe_health_test.go:420` `TestProbeSuiteSinkFailurePropagates`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **60 `probe_health_test.go:457` `TestProbeSuiteAudioStaysUnprobed`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **61 `probe_thinking_test.go:84` `TestProbeSuiteMeasuresBrokenThinking`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **62 `probe_thinking_test.go:120` `TestProbeSuiteThinkingCapableIsNotTheSameAsBroken`** — 摘：`pf := newProbeFixture(t)`。判：同上 ⇒ 可转。
- **63 子测试 `TestProbeSuiteBrokenOnAllThreeDialects/openai-chat`**（父 `probe_health_test.go:306`）— 摘：`pf := newProbeFixture(t)`（`:309`）。判：fixture 开库被未解析根拒，测方言能力判读要成 ⇒ 可转。
- **64 子测试 `TestProbeSuiteBrokenOnAllThreeDialects/anthropic`**（同父/同 assert）— 摘：`pf := newProbeFixture(t)`（`:309`）。判：同上 ⇒ 可转。
- **65 子测试 `TestProbeSuiteBrokenOnAllThreeDialects/openai-responses`**（同父/同 assert）— 摘：`pf := newProbeFixture(t)`（`:309`）。判：同上 ⇒ 可转。
- **66 子测试 `TestProbeSuiteCapableOnAllThreeDialects/openai-chat`**（父 `probe_health_test.go:352`）— 摘：`pf := newProbeFixture(t)`（`:355`）。判：同上 ⇒ 可转。
- **67 子测试 `TestProbeSuiteCapableOnAllThreeDialects/anthropic`**（同父/同 assert）— 摘：`pf := newProbeFixture(t)`（`:355`）。判：同上 ⇒ 可转。
- **68 子测试 `TestProbeSuiteCapableOnAllThreeDialects/openai-responses`**（同父/同 assert）— 摘：`pf := newProbeFixture(t)`（`:355`）。判：同上 ⇒ 可转。

### 5.6 `./internal/perm/`（5 枚）— 全可转

> 逐字：A2L perm 红串 = `SaveFile: config: config.toml temp file seal: winsec: refusing to seal /varlink/…`。测三档权限持久化（要成），与拒无关。

- **69 `ticket90_persist_test.go:138` `TestTicket90ManualSwitchSurvivesRestart`** — 摘：`path := writeConfig(t, dir, risk.ModeAskEveryStepName)`（`:140` 处 `t.Fatalf("SaveFile: %v", err)`）。判：写 config 到未解析根被拒 ⇒ 可转。
- **70 `ticket90_persist_test.go:181` `TestTicket90UntouchedConfigStartsAtTheDefault`** — 摘：`t.Fatalf("SaveFile: %v", err)`（`:190`）。判：同上 ⇒ 可转。
- **71 `ticket90_persist_test.go:220` `TestTicket90SessionGrantDoesNotSurviveRestart`** — 摘：`t.Fatalf("memory.Open: %v", err)`（`:227`）。判：开库 setup 被未解析根拒 ⇒ 可转。
- **72 `ticket90_persist_test.go:291` `TestTicket90ConfigKeyChangesWhatTheChainAsks`** — 摘：`path := writeConfig(t, dir, tc.key)`。判：写 config 被拒 ⇒ 可转。
- **73 `ticket90_persist_test.go:336` `TestTicket90HandEditLooseningGoesThroughD36`** — 摘：`path := writeConfig(t, dir, risk.ModeAskEveryStepName)`。判：同上 ⇒ 可转。

### 5.7 `./internal/memory/`（33 枚，全顶层）— 全可转

> 逐字：A2L 里 33 枚红中 32 枚直报 `Open: memory: create data dir: winsec: refusing to seal /varlink/… not provably resolved`（DAO 全在 `openTestStore(t)` / `inv76Fixture(t)` 开库这步被拒，真正要断言的 DAO/quota/迁移行为根本没跑）；第 33 枚 `TestCrashRecoveryKillMidWrite` 报的是 `subprocess only committed 0 rows`——根因也是未解析根（见 85 行判据）。⚠ 两枚名字带 "Rejects…Hostile Shapes"、三枚带 "Unmigratable"，都**不是**未解析根腿：前者的"拒"是 artifact 删除的路径不变式形状（且它们红在 fixture setup、根本没到删除步），后者要的是 `ErrSchemaUnmigratable`（另一种错），普通形已满足 ⇒ 全可转。

- **74 `artifacts_path_invariant_test.go:314` `TestDeleteArtifactRejectsTheFourHostileShapes`** — 摘：`s, _, _, _, shapes := inv76Fixture(t)`（`:315`）。判：红在 fixture 开库被未解析根拒；"四种 hostile shape" 是删除路径不变式、非根，要成 ⇒ 可转。
- **75 `artifacts_path_invariant_test.go` `TestDeletePrivacyItemRejectsTheFourHostileShapes`** — 摘：`s, _, _, _, shapes := inv76Fixture(t)`（`:345`）。判：同上（红在开库）⇒ 可转。
- **76 `artifacts_path_invariant_test.go:389` `TestArtifactsContainmentByDirectoryListing`** — 摘：`s, root, dataDir, userDir, shapes := inv76Fixture(t)`（`:390`）。判：开库被拒，测 containment 要成 ⇒ 可转。
- **77 `artifacts_path_invariant_test.go:560` `TestArtifactsLiteralBackslashKeepsItsAsymmetry`** — 摘：`s, root, _, _, shapes := inv76Fixture(t)`（`:561`）。判：开库被拒 ⇒ 可转。
- **78 `artifacts_stray_test.go:87` `TestStraySubdirectoryCannotHideBytesFromTheQuota`** — 摘：`s := openTestStore(t, WithLogger(logger))`（`:89`）。判：开库被未解析根拒，要成 ⇒ 可转。
- **79 `artifacts_stray_test.go` `TestStraySubdirectorySharesTheQuotasLRUQueue`** — 摘：`s := openTestStore(t)`（`:136`）。判：同上 ⇒ 可转。
- **80 `artifacts_stray_test.go` `TestPurgeArtifactsReclaimsStraySubdirectory`** — 摘：`s := openTestStore(t, WithLogger(logger))`（`:179`）。判：同上 ⇒ 可转。
- **81 `artifacts_stray_test.go:223` `TestStrayRemovalDoesNotFollowLinks`** — 摘：`s := openTestStore(t)`（`:224`）。判：红在开库被未解析根拒；它测的是「清游离子目录不跟软链」（普通形已绿），与根无关 ⇒ 可转。
- **82 `concurrent_test.go:34` `TestConcurrentWritersReaders`** — 摘：`s := openTestStore(t)`（`:35`）。判：开库被拒 ⇒ 可转。
- **83 `dao_providerhealth_test.go` `TestProviderHealthProbeRoundTrip`** — 摘：`s := openTestStore(t)`（`:13`）。判：开库被拒，测 provider health DAO 要成 ⇒ 可转。
- **84 `dao_providerhealth_test.go` `TestProviderHealthRecordErrorAndList`** — 摘：`s := openTestStore(t)`（`:56`）。判：同上 ⇒ 可转。
- **85 `concurrent_test.go:220` `TestCrashRecoveryKillMidWrite`**
  摘：`t.Fatalf("subprocess only committed %d rows within budget (want >= %d)", committed.Load(), killAfter)`（`concurrent_test.go:254`）
  判定：**可转（带一处落地提醒）**。判据：软链形它红在「子进程只提交了 0 行」——父测把 `dir := t.TempDir()/data`（未解析 `/varlink`）经 `WISP_CRASH_DIR` 交给子进程，子进程 `Open(dir)` 被未解析根拒 ⇒ 0 行 ⇒ 30 s 超时。红串不是字面 "not provably resolved"（子进程 stderr 未被扫），但**根因仍是未解析根**。AC#2b 转换时须把**传给子进程的 `WISP_CRASH_DIR` 也换成解析后的根**，否则这条接不上（普通形已绿，非拒绝腿）。
- **86 `dao_providerhealth_test.go` `TestProviderHealthMissingRowIsNoRows`** — 摘：`s := openTestStore(t)`（`:117`）。判：开库被拒 ⇒ 可转。
- **87 `dao_test.go` `TestProfileUpsertAndLRUEvictionLogs`** — 摘：`s := openTestStore(t, WithLogger(logger))`（`:55`）。判：同上 ⇒ 可转。
- **88 `dao_test.go` `TestProfileDeleteAndPurge`** — 摘：`s := openTestStore(t)`（`:136`）。判：同上 ⇒ 可转。
- **89 `dao_test.go` `TestMemoryAddSearchDeletePurge`** — 摘：`s := openTestStore(t)`（`:164`）。判：同上 ⇒ 可转。
- **90 `dao_test.go` `TestTaskLogLifecycleAndErrorClassValidation`** — 摘：`s := openTestStore(t)`（`:250`）。判：同上 ⇒ 可转。
- **91 `dao_test.go` `TestToolCallLifecycleAndValidation`** — 摘：`s := openTestStore(t)`（`:305`）。判：同上 ⇒ 可转。
- **92 `dao_test.go` `TestGrantCostPluginStateDAO`** — 摘：`s := openTestStore(t)`（`:379`）。判：同上 ⇒ 可转。
- **93 `dao_test.go` `TestWriteQueueLazyLifecycle`** — 摘：`s := openTestStore(t, WithLogger(logger), WithRegistry(...))`（`:467`）。判：同上 ⇒ 可转。
- **94 `dao_test.go` `TestAdvPanicWedge`** — 摘：`s := openTestStore(t, WithLogger(logger), WithRegistry(...))`（`:526`）。判：同上 ⇒ 可转。
- **95 `migrate_v1seed_test.go` `TestMigrateSeededV1DatabaseToV2`** — 摘：`t.Fatalf("open v1: %v", err)`（`:26`）。判：open v1 写未解析根被拒 ⇒ 可转。
- **96 `privacy_test.go:16` `TestPrivacyOpsAllDomains`** — 摘：`s := openTestStore(t)`（`:17`）。判：开库被拒 ⇒ 可转。
- **97 `retention_test.go:112` `TestRetentionBoundaries`** — 摘：`s := openTestStore(t)`（`:113`）。判：同上 ⇒ 可转。
- **98 `retention_test.go` `TestRetentionSchedule`** — 摘：`s := openTestStore(t, WithLogger(logger))`（`:191`）。判：同上 ⇒ 可转。
- **99 `retention_test.go:227` `TestArtifactsLRUQuota`** — 摘：`s := openTestStore(t, WithLogger(logger))`（`:229`）。判：同上 ⇒ 可转。
- **100 `retention_test.go:271` `TestArtifactsDeleteAndPurgeAndGuards`** — 摘：`s := openTestStore(t)`（`:272`）。判：同上 ⇒ 可转。
- **101 `schema_test.go:171` `TestSchemaContractIntrospection`** — 摘：`s := openTestStore(t)`（`:172`）。判：开库被拒 ⇒ 可转。
- **102 `schema_test.go` `TestMigrateFreshDatabaseCreatesCurrentVersion`** — 摘：`s := openTestStore(t)`（`:323`）。判：同上 ⇒ 可转。
- **103 `schema_test.go` `TestMigrationChainWithBackup`** — 摘：`t.Fatalf("open v1: %v", err)`（`:386`）。判：open v1 被未解析根拒 ⇒ 可转。
- **104 `schema_test.go:456` `TestMigrationFailedStepIsAtomic`** — 摘：`t.Fatal(err)`（`:460`，setup open/migrate）。判：开库/迁移 setup 被拒 ⇒ 可转。
- **105 `schema_test.go` `TestMigrationNewerSchemaIsUnmigratable`** — 摘：`t.Fatal(err)`（`:512`，setup 写 seeded v-newer 库）。判：开库被未解析根拒，要成 ⇒ 可转。
- **106 `schema_test.go` `TestMigrationForeignDatabaseIsUnmigratable`** — 摘：`t.Fatalf("error = %v, want ErrSchemaUnmigratable", err)`（`:579`）。判：软链形它拿到的 err 是未解析根拒（非 `ErrSchemaUnmigratable`）⇒ `want` 落空。它要的错误是**迁移语义的 `ErrSchemaUnmigratable`**、不是未解析根；解析根后能开那枚 foreign 库并正确返回该错（普通形已绿）⇒ 可转。

### 5.8 `./internal/tools/`（21 枚：19 顶层 + 2 子测试）— 20 可转 + 1 归因待票 123 裁

> ⚠ 本包的红**分两种机制**（都是未解析根，但文案不同），逐字都回看过：
> (甲) 直接被密封底线拒：`memory.Open / composeLoop → refusing to seal /varlink/… not provably resolved`（109、124、125 等 setup 腿）；
> (乙) 路径规范化器/C26 判「那棵树不在盘上、不授权」⇒ `InAllowlist=false` ⇒ 写升级 L2 ⇒ `L2 审批未通过，已拒绝执行`（RealToolCall/FS*/DeclaredRisk/Ticket107*/Canonicalize 等）。
> 乙类的红串**不含**字面 "not provably resolved"，但根因同为未解析根（普通形全绿；解析根后规范化器认树、写留 L1）。除 `TestL1Write`（丙类，300 s C18 审批超时，见下行）外，全可转。

- **107 `bridge_path_account_ticket105_test.go:89` `TestRealToolCallWritesRewriteAccountIntoAudit`**
  摘：`t.Fatalf("precondition: the read must actually run, got %q", out.Text)`（`:112`；实拿 `got "L2 审批通道尚未接入（票 21），已拒绝执行"`）。判：乙类——读目标在未解析根、不被认授权→升 L2→被拒；断言要「读必须真跑」，要成 ⇒ 可转。
- **108 `bridge_test.go:608` `TestDeclaredRiskIsOnlyAFloor`**
  摘：`t.Fatalf("in-allowlist read: out=%+v err=%v", out, err)`（`:618`；`out={…L2 审批未通过，已拒绝执行…}`）。判：乙类（in-allowlist 读被升 L2 拒），要成 ⇒ 可转。
- **109 `bridge_test.go:649` `TestToolCallRowsAreComplete`**
  摘：`store := openStore(t, filepath.Join(root, "data"))`（`:651`；实拿 `memory.Open: create data dir: refusing to seal /varlink/… not provably resolved`）。判：甲类 setup 被拒，测 tool_call 行完整性要成 ⇒ 可转。
- **110 `fs_test.go:16` `TestFSReadReturnsTheFileAndTaintsIt`** — 摘：`t.Fatalf("out=%+v err=%v", out, err)`（`:24`；`out={…L2 审批通道尚未接入，已拒绝执行…}`）。判：乙类，要成 ⇒ 可转。
- **111 `fs_test.go:50` `TestFSReadTaintFeedsR4`**
  摘：`t.Fatalf("tainted exfil judged %v, want L2 (decision=%+v)", d.Level, d)`（`:80`；实拿 `tainted exfil judged L0, want L2`）。判：乙类——因源读被未解析根挡住→没打上 taint→外渗被判 L0（要的是 L2 这档**风险判级**，非未解析根的拒）；解析后读得成、taint 生效、判回 L2（普通形绿）⇒ 可转。
- **112 `fs_test.go:93` `TestFSListSummarizesADirectory`** — 摘：`t.Fatalf("out=%+v err=%v", out, err)`（`:101`）。判：乙类，要成 ⇒ 可转。
- **113 `fs_test.go:118` `TestFSListHonoursItsCap`** — 摘：`t.Fatalf("out=%+v err=%v", out, err)`（`:127`）。判：乙类，要成 ⇒ 可转。
- **114 `fs_write_test.go:360` `TestAtomicWriteKillsMidWrite`（顶层）** — 摘：子测试 `t.Errorf("a new file is the L1 row, so the window route was expected (window=%d)", w)`（`:424`）。判：乙类（新文件本应 L1、被升 L2 无窗口）；顶层红由 2 子测试聚合 ⇒ 可转（见 126、127）。
- **115 `fs_write_test.go:477` `TestFSTrashGoesToTheRecycleBin`**
  摘：`if !strings.Contains(out.Text, "回收站") { t.Errorf("the refusal must say why: %q", out.Text) }`（`:500-501`；实拿 `the refusal must say why: "L2 审批未通过，已拒绝执行"`）
  判定：**可转（关键辨析）**。判据：它**确实**断言一个拒——但那是「POSIX 无回收站时 fs.trash 必须拒并说明带『回收站』」（`recycleBinSupported()==false` 分支），**不是**未解析根的拒。软链形里未解析根先把写升 L2、被审批拒（文案 `审批未通过`），抢在回收站检查前 ⇒ 该断言拿错文案。普通形里回收站腿的拒照拿（已绿），解析根后回到那条拒 ⇒ **它要的那句拒不被洗绿**。可转。
- **116 `fs_write_test.go:642` `TestFSMoveSameVolume`** — 摘：`t.Fatalf("out=%+v err=%v", out, err)`（`:655`；`L2 审批未通过，已拒绝执行`）。判：乙类，要成 ⇒ 可转。
- **117 `paths_rewrite_ticket102_test.go:17` `TestPathCanonicalizerAccountsForRewrittenRoots`** — 摘：`t.Errorf("InAllowlist(%q) = false although the expanded root is a real tree; the fix must not turn into option (A)", c)`（`:64`）。判：乙类（未解析根致 expanded root 不在盘确认→不授权），断言要 InAllowlist=true，要成 ⇒ 可转。
- **118 `paths_ticket107_portable_test.go:27` `TestTicket107AllowlistJudgmentTwoShapes`** — 摘：`t.Errorf("shape R (root already resolved) InAllowlist(%q) = false, want true", …)`（`paths_ticket107_portable_test.go:63`）。判：乙类，要成（InAllowlist 应为真）⇒ 可转。
- **119 `paths_ticket107_portable_test.go:80` `TestTicket107AllowlistBoundaryIsComponentWise`** — 摘：`t.Errorf("InAllowlist(%q) = false, want true: a file under the allowed root must pass", inside)`（`:99`）。判：乙类，要成 ⇒ 可转。
- **120 `paths_ticket107b_probes_test.go:144` `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`**
  摘：`InAllowlist("<T>") = false, want true: "<T>" is a real tree the operator named`（`:170`；前置 `:165` 记 `C26 expanded it onto <T> but that tree is not confirmed on disk, authorizing nothing`）
  判定：**可转**。判据：名字 "StaysOutside" 讲的是「允许根内那条软链的探针应在根外」这一 allowlist 边界形状（非未解析根的拒）；它红在 `InAllowlist=want true 却 false`——未解析根致 C26 判该树不在盘、不授权。解析根后该树被确认，边界判定回到普通形（绿），"stays outside" 那半照测 ⇒ 可转。
- **121 `paths_workspace_test.go:41` `TestWorkspaceSwitchNarrowsWhatTheAssessorJudges`** — 摘：`t.Fatalf("premise: an in-allowlist write with no workspace chosen assessed %q (%v), want L1", …)`（`:53`；实拿 `assessed "L2" ([R1 R2]), want L1`）。判：乙类（未解析根→不在授权目录→升 L2），断言要留 L1，要成 ⇒ 可转。
- **122 `pathshape_portable_test.go:28` `TestCanonicalizeReturnsAPathTheOSCanOpen`** — 摘：`t.Fatalf("InAllowlist(%q) = false for a file inside the only allowed root %q: R2 would send every plain write to L2", …)`（`:51`）。判：乙类，要成 ⇒ 可转。
- **123 `loop_approval_test.go:165` `TestLoopPassesDeclaredL1WriteThroughTheGate`** — 摘：`f := composeLoop(t, dir, true)`（`:167`；setup 内 `memory.Open … refusing to seal /varlink/…`）。判：甲类 setup 被拒，要成 ⇒ 可转。
- **124 `loop_approval_test.go:226` `TestLoopStillRefusesL1WhenNoGateIsRegistered`**
  摘：`f := composeLoop(t, dir, false)`（`:228`；实拿 `memory.Open: create data dir: refusing to seal /varlink/… not provably resolved`）
  判定：**可转**。判据：名字 "StillRefuses" 要拒的是「无门时 L1 写被拒」这一**闸门语义**（普通形已满足）；它红在 `composeLoop` 的 `memory.Open` setup 被未解析根拒、根本没走到无门断言。解析根后 setup 成、"无门仍拒" 照测 ⇒ 可转，不洗绿。
- **125 `wiring_test.go:112` `TestL1WriteGoesThroughTheRealBlockWindow`**
  摘：`t.Fatalf("the window ran out unopposed, which MEANS execute: %+v", out)`（`:124`；实拿 `审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`，`--- FAIL (300.04s)`、普通形 `--- PASS (3.01s)`）
  判定：**归因待票 123 裁**（不替它下结论，按票面 AC#2裁定 附带发现 2）。判据：它红在 **300 s C18 审批超时常量**（文案逐字「审批超时（300 秒未确认），C18 一律判拒绝」），是「用例假设有人在旁边点确认」那一族（票 123）在 POSIX 软链形下的显形，**不是**「未解析根被拒」这一族。本账把它从三档里单拎出来；AC#2b 不得把它当"未解析根"接解析变绿，也不得当拒绝腿留红——归口待票 123。
- **126 子测试 `TestAtomicWriteKillsMidWrite/new_file_target_does_not_appear_at_all`**（父 `fs_write_test.go:360`）
  摘：`t.Errorf("a new file is the L1 row, so the window route was expected (window=%d)", w)`（`:424`）。判：乙类（未解析根→升 L2、不走 L1 窗口），要成 ⇒ 可转。
- **127 子测试 `TestAtomicWriteKillsMidWrite/a_clean_write_lands_and_round_trips`**（同父）
  摘：`t.Fatalf("out=%+v err=%v", out, err)`（`:438`；`L2 审批未通过，已拒绝执行`）。判：乙类，要成 ⇒ 可转。

### 5.9 `./cmd/wisp/`（5 枚：4 顶层 + 1 子测试）— 全可转

> 逐字：这 5 枚红串都直报 `memory: create data dir: winsec: refusing to seal /varlink/… not provably resolved`——全在 fixture/store setup 这步。⚠ 同包另外 28 枚是**两形都红**（DPAPI/命令面腿那本旧账，AC#1 §5.3 已排除在本票 132 之外），此处不列、不判。

- **128 `providers_test.go:133` `TestProvidersProbeRecordsMeasuredThinkingFalse`** — 摘：`pf := newProvidersFixture(t)`（`:133`）。判：fixture 开库被未解析根拒 ⇒ 可转。
- **129 `providers_test.go:179` `TestProvidersProbeRecordsMeasuredThinkingTrue`** — 摘：`pf := newProvidersFixture(t)`（`:179`）。判：同上 ⇒ 可转。
- **130 `providers_test.go:208` `TestProvidersDiscoverListsWhatTheServerServes`** — 摘：`pf := newProvidersFixture(t)`（`:208`）。判：同上 ⇒ 可转。
- **131 `providers_test.go:228` `TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`** — 摘：`pf := newProvidersFixture(t)`（`:228`）。判：同上 ⇒ 可转。
- **132 子测试 `TestSecretFailurePathsLogAndPrintNoPlaintext/unset_name_(no_such_blob)`**（父在 `secret_test.go`）
  摘：`t.Errorf("the failure must name the ref, got: %s", q.errb)`（`secret_test.go:679`；实拿 `got: … secret: create /varlink/…: refusing to seal … not provably resolved`）
  判定：**可转**。判据：它测「unset 不存在的名字时，失败信息要点名那个 ref、不漏明文」（错误卫生断言）；软链形里前置 `secret set` 落未解析根被拒（`:689 set: exit 1 …`）⇒ blob 没建成、unset 拿到的错误是密封拒而非「ref 不存在」的干净错误 ⇒ `must name the ref` 落空。普通形里该子测试绿（同包 `unset_name_(no_such_blob)` 未进 both-red 的 28 枚），解析根后回到普通形 ⇒ 可转。

## 6. 与 AC#1 那 132 枚的对账（逐名，含「变绿 vs 被跳过」与多/少一枚）

**口径**：本账在 `bcb03aa` 纯净快照、受影响 9 包 `-count=1 -v` 两形各一枚（`A2L`/`A2P`）。逐名差集（`link FAIL` 减 `plain FAIL`）：

| 包 | 本账 only-in-link | AC#1 §5.3 only-in-link | 差 |
|---|---|---|---|
| `internal/agent` | 26 | 26 | 0 |
| `internal/agent/approval` | 1 | 1 | 0 |
| `internal/config` | 7 | 7 | 0 |
| `internal/llm` | 17 | 17 | 0 |
| `internal/memory` | 33 | 33 | 0 |
| `internal/perm` | 5 | 5 | 0 |
| `internal/tools` | **22** | 21 | **+1（`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`）** |
| `internal/winsec` | 17 | 17 | 0 |
| `cmd/wisp` | 5 | 5 | 0 |
| **合计** | **133** | **132** | **+1** |

- **AC#1 的 132 枚逐名全在**：脚本核对 `list132.txt`（从 AC#1 §5.3 逐字抄）——132 枚在 `A2L` 全 `--- FAIL`、在 `A2P` 全 `--- PASS`，**零枚对不上**。⇒ AC#1 那本名册在 `bcb03aa` 上仍逐字成立。
- **多的一枚 = `TestLateVetoRendersTheApprovalLayersAppliedStepsReport`（tools）**：软链形 `--- FAIL (300.02s)`、普通形 `--- PASS (3.00s)`。它与 AC#1 已单列的 `TestL1WriteGoesThroughTheRealBlockWindow`（同为 `300.0x s`、同为 `审批超时（300 秒未确认），C18 一律判拒绝`）**同属 C18 审批超时族**，红因不是未解析根。AC#1 只逮到 `TestL1Write` 一枚、漏了这枚。
  - 旁证：AC#1 记 tools 软链形 `RUN=77`（比 plain 的 79 少 2，因 300 s 腿在 setup 就被拒、子测试不展开），本账 tools 软链形 `RUN=79`——差别正在这两枚审批腿：本快照里它们**跑满了 300 s**（进了窗口、没被 setup 提前挡），于是各自计入 `RUN` 并 FAIL。这不影响「未解析根那一族」的 132 枚名册，只是把 `TestLateVeto` 也带进了 only-in-link。
  - ⇒ 这条**不是**「少/多一枚未解析根红」，而是「C18 超时腿是否跑满」的形态差；`TestLateVeto` 应和 `TestL1Write` 一起归口票 123，不进 AC#2b 转换范围。
- **「变绿还是被跳过」逐包核（§5.4 附带发现，本账独立复现）**：软链形相对普通形新出现的 `SKIP` 只有 4 枚——`internal/winsec` 3 枚（`TestAC2POSIXSeam{ProbeShapes,Guard,Accepts}…125`）+ `internal/config` 1 枚（`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`）。这 4 枚**全不在 132 名册里**（是票 125 那族自认形状、拒测的 seam 探针），与本账要判的 132 无交集。⇒ 本账 132 枚（含特殊 `TestL1Write`）的「普通形不红」**全部是 `PASS`、无一靠 `SKIP` 蒙过去绿**（脚本核：133 枚 only-in-link 在 `A2P` 状态逐枚 = `PASS`）。若 AC#2b 把根解析干净，这 4 枚 seam 探针应回到 `PASS` 侧，别算进「被清掉的红」。

## 7. 未验证项 + `next=`

**未验证项**：
1. **macOS 那半仍未实测**（本机无 macOS runner）——本账只在 Linux 容器软链 `TMPDIR` 复现；票面《事实》说这形状代 macOS（`/var` 是链接）那一形，与票 119 `R-119-8` 同账，不背书。
2. **`TestL1Write` / `TestLateVeto` 的 C18 归因**：本账只**坐实其红因是 300 s 审批超时、非未解析根**（文案逐字为证），**未**裁它该在票 123 怎么修——归口待 123。
3. **`memory.TestCrashRecoveryKillMidWrite` 的可转性有一处前提**（行 85 判据）：接解析时须把交给子进程的 `WISP_CRASH_DIR` 也换成解析后的根；本账**未在改码后复算**（AC#2b 的活），故该枚判「可转（带落地提醒）」，非「已验可转」。
4. **"+2 枚被拒路径归属"**（AC#1 §5.1，票面 79→今 81）：本账不动它——它属 AC#2b 复算范围（票面 AC#2裁定 附带发现 3）。

**`next=`（交编排者 / 供 AC#2b）**：
- 本账三档：**可转 131 / 拒绝腿 0 / 待裁 1（`TestL1Write` 归票 123）**。⇒ **AC#2b 可对全部 131 枚「可转」接上票 119 那条「先解析再递底线」纪律**；**没有任何一枚需要作为拒绝腿保留红**（结案判据里"软链形红名数 = 拒绝腿档枚数"这一项，若 131 全接解析，理论上归零；余下红名应只剩 `TestL1Write`/`TestLateVeto` 两枚 C18 超时腿，而那两枚归票 123 裁，不在本票清零目标内）。
- ⚠ AC#2b 三条硬提醒（账上已逐枚标注）：① `TestCrashRecoveryKillMidWrite` 要解析**交给子进程**的那个根；② 断言里带 "Refuses/Rejects/want Err…" 的 6 枚（winsec AC5/AC118/AC3 反向、memory 迁移腿、tools 回收站、llm RefusesSilentRuns）接解析后**必须仍拿到它们各自要的那种拒**（不是未解析根），复算时逐枚点红名自证；③ 归零结论**普通形 + 软链形各一枚**，且两形 `RUN/SKIP` 差要逐包解释（票面 AC#2裁定 附带发现 1 已升为跑法要求）。
