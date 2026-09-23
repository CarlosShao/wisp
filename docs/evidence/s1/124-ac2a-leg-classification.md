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
