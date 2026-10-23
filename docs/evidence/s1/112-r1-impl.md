# 112-r1 —— winsec 步在 CI runner 上第一次跑就红（写腿交件）

**本件是读数与交件，不是裁决。** 裁决表 `112-*.md` 由验收方出；票面 `Status:`／框／`-done` 一字未动。

## 1. 票面那一步 / 那句红（逐字）

票面第 3 行（逐字）：`**Status:** ready-for-review（2026-09-21 20:3x 编排者建；来源=**run \`35595651898\` / job \`106319703680\`**，票 110 新加的 "Windows ACL sealing gate" 步）`

票面 10–18 行（逐字，`gh api .../jobs/106319703680/logs`）：

```
--- FAIL: TestAC3JunctionInputIsRefusedNotSealed (0.09s)
  --- FAIL: TestAC3JunctionInputIsRefusedNotSealed/existing_directory_behind_the_link
  --- FAIL: TestAC3JunctionInputIsRefusedNotSealed/missing_directory_under_the_link
--- FAIL: TestC26PipelineIsWiredIntoWinsec (0.00s)
--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.02s)
FAIL github.com/CarlosShao/wisp/internal/winsec  13.707s
winsec-tests.sh: === RUN=71  --- PASS=32  --- FAIL=3  --- SKIP=0
```

**今天这一步仍然是红的**（正面形状判据：`gh api repos/CarlosShao/wisp/actions/runs/39884431726/jobs` 的 `steps[].number/status/conclusion` 逐条读，不是从任务失败倒推）：

- run `39884431726`（head `23b663c2`，`conclusion=failure`，`2026-10-23T01:31Z` 起）/ job `120189170762`（`test-windows`）/ **step 4** `Windows ACL sealing gate (internal/winsec's own tests, ticket 110)` = `completed` / **`failure`**
- 步内四数逐字：`winsec-tests.sh: four numbers (all from -v output): === RUN=122  --- PASS=69  --- FAIL=2  --- SKIP=0`
- 步 rc：`winsec-tests.sh: gate is RED (rc=1) for scope=[./internal/winsec/]`

## 2. 复认：票面**那三条**今天已经 PASS（逐字，同一份 job log）

```
--- PASS: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.05s)
--- PASS: TestC26PipelineIsWiredIntoWinsec (0.02s)
--- PASS: TestAC3JunctionInputIsRefusedNotSealed (0.12s)
    --- PASS: .../existing_directory_behind_the_link
    --- PASS: .../missing_directory_under_the_link
    --- PASS: .../with_only_the_built-in_verifier
```

⇒ 票面 65–114 行那两段 Progress log 是准的：9-21 两腿把三条愈掉了（①落在 104 那条根上）。**票面第 1 行标题描述的状态已经过期**（本仓票面腐坏又一例，只登记、不改票名）。

## 3. 门今天红在哪（逐字，同一 job step 4）

```
--- FAIL: TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1 (0.52s)
  ancestorunlink_348_windows_test.go:131: AC#2 RED instrument: no leaf leg answered (guard=0, landing's re-check=0) while reparseAt proved 110 flipped leaf cycles across 400
--- FAIL: TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1 (0.79s)
  openleaf_guard_322_windows_test.go:241: AC#3 RED instrument: 400 raced attempts produced guard-leg refusals=0 and racer flips=411 (flipErrs=0, firstFlipErr=""), so nothing below proves the guard is on this route
```

两条都不是「判定错」，是**「仪器什么都没量到」**（RED instrument 的开场断言）。

## 4. 「本地绿、CI 红」的具名差异（逐条量过，⛔ 猜）

| 对象 | 本机（`go test ./internal/winsec/ -run ... -count=3 -v`） | runner（job 120189170762） | 尺 |
|---|---|---|---|
| 322 r1 叶子守卫的**命中数** | `guardRefusals=8` / `8` / `15`（每次 attempts=400） | `guardRefusals=0`（attempts=400） | 用例自打的 `AC#3 local paired reading` 行 |
| 322 r1 的 flips | 319 / 382 / 408 | 411 | 同上 |
| 342 r1「哪条腿答的」里的 **C26 稳态拒绝** | `205 attempts` / 400 | `288 attempts` / 400 | `ancestoropen_342_windows_test.go:48` 的「leg that answered」行 |
| 342 r1 AC#1 的**截断分子** | `truncated=2` / 400 ⇒ 红 | `truncated=0` ⇒ 绿 | 同 `:48` 行 |
| 342 r1 AC#2 的 guard / re-check | `guard=0, re-check=0` | `guard=0, re-check=0` | `ancestorunlink_348_windows_test.go:131` |

**具名差异＝「哪条腿先答」的速率差，不是能力差。** `winsec.go:167` 的 `preOpenLeafGuard` 只在 `resolveString` 答完之后、`pinnedParent` 之前那一段里才可能轮到；runner 上 C26 稳态拒绝把 72% 的尝试先吃掉（本机 51%），于是守卫 400 发一次没轮到。**这条用例的红/绿是胜率，不是性质** ⇒ 门色随机器快慢抛硬币。

本机整包三次，**三次红名不同**（同一份 HEAD、同一晚）：`TestAC1...342r1`＋`TestAC2...342r1`（第一发）→ `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（第二发，`RUN=127 PASS=75 FAIL=1 SKIP=0`）。这就是「门恒红且红名会漂」的实测形状。

## 5. 为什么这不是「环境差异」，而是门的结构问题

- 票 **342 未结**（`342-...-no-ticket-owns-it.md` 无 `-done`，`Status: 已立，未派`）⇒ 它的「修前必红」探针就住在 CI 这道门里 ⇒ **342 落地前这道门不可能绿**。
- 票 **322 已 `-done`**，但它那枚「守卫在不在生产路上」的证明**只能靠抢窗口**：`platformVerifyPlacement` 的 `pathPieces` 走**包含叶子**（`pathpieces_108_test.go:84` 现读 `got := pathPieces(r.input)`，`:85` 打的读数里 `pathPieces("C:\\a\\b\\keep.txt")` 第三段就是叶子）⇒ 链接**先站着**时永远是稳态腿先答；**确定性证明「守卫在 gap 里」需要产码注入点**。
- 票 **346 已把这条机制具名**（「弃权报成红 ⇒ 整仓 CI 随机红」），但它的**射程只到两枚 POSIX 仪器**（票面第 4 行逐字：`射程＝\`internal/winsec/openleaf_322_other_test.go\` ⇄ \`internal/winsec/posix_seal_window_120_other_test.go\` 的**弃权那一支**`；且 346-a1 现量两枚首行＝`//go:build !windows`）⇒ **Windows 这两枚弃权仪器今天无主**（本腿登记，不认领）。

## 6. 选了哪一形 ＋ 为什么它最有价值

甲（修好它）＝要动产码注入点（＝346 的形乙，票面自己写「须另票＋人工批准」），且 342 未结、同文件有在飞腿；乙（显式 skip）＝票 112 `AC#2` 与票 346 `AC#2` 两处**明文禁手**；⇒ **丙（改成 runner 上可跑的形）**，但只落在**我能独自证成的那一半**：

**把「叶子腿在不在生产路上」的稳态半边，从「抢窗口的胜率」改成「不靠运气的确定性分母」**，并把「哪条腿答的」逐字记账。

★为什么最有价值：①它**今天就能量**（不需要 runner、不需要产码、不需要批准）；②它把 runner 上 `guardRefusals=0` 从「谜」变成「可解释、可复算的结构事实」——这正是 346 那条「归因吃轮次」的病根；③它**反向可证**（下面第 8 节三读数），而胜率形**反向也证不动**。⛔ 它**不**把门改绿，也**不**冒充绿：票 112 `AC#5` 本腿**未达成**（见第 10 节）。

## 7. 改了什么（逐行）

新件**一枚**：`internal/winsec/leafstand_112_windows_test.go`（232 行 / 11917 字节 / LF / 四枚用例），`//go:build windows`、`package winsec`。**产码 0 行、`ci.yml` 0 行、`scripts/**` 0 行、skip 0 处、既有一枚断言未放宽**。

- `passthrough112r1`（`Resolve`＋`ResolveAccounted`）：经 `export_test.go` 的 `SetSeamForTest` 装一个「原样回答」的解析器，把**地板**单独暴露出来量（真 C26 会把地板遮住——正是本票要治的「遮住」本身）。首版只写了 `Resolve`，实跑被 `does not implement RewriteAccounted` 拒 ⇒ 这条拒本身是票 108 `AC#4` 的守卫在工作，已在注释里记。
- `legThatAnswered112r1`：把答的腿**按逐字串**归到三枚（地板 `traverses a reparse point at`／C26 `path traverses a reparse point`／守卫 `whyGuardRefusal322r1`），第四种一律「未具名 ⇒ 失败」——⛔ 拿「err != nil」当通过（120-r3 量过那是假绿）。
- `plantLeafLink112r1`：临时树里**显式种**一条指向树外见证文件的文件链接，种完**先量自己**（`reparseAt` 不为真就 `Fatalf("measured nothing")`）。
- 四枚用例逐名见第 8 节；`AC#3` **不**断言 `ErrUnresolvedPath`，因为 `resolveAccounted` 按设计原样透传解析器的错（「refusal 要归 C26 而不是本包的意见」），那条哨兵只在**地板腿**成立——这是读数，不是放宽。

## 8. 反向判据：三读数 ＋ 还原核

台面：`git archive HEAD` 导出树 `D:\tmp\wisp\20261023\112r1\export-asis`（仓外，⛔ 工作树突变窗口），新件拷进去；每发前 `go build ./internal/winsec/` 都 rc=0（编译失败不算变异）。

| 读数 | 命令 | rc | 逐字 |
|---|---|---|---|
| **原样绿** | `go test ./internal/winsec/ -run 112r1 -count=1 -v` | **0** | `=== RUN` 4 / `--- PASS` 4（AC#1–AC#4 全绿） |
| **变异红 M1**（`placement_windows.go:65` `if false && isReparsePoint(prefix)`＝把地板的叶子那一跳摘掉） | 同上 | **1** | `leafstand_112_windows_test.go:147: AC#1 RED: the leg that answered was preOpenLeafGuard (the gap after the answer), want the placement floor's reparse walk ("traverses a reparse point at")` ＋ `--- FAIL: TestAC1AStandingLeafLinkIsRefusedByTheFloorBeforeAnyByte112r1`（另三枚仍绿＝红在守卫，不是噪声） |
| **变异红 M2**（M1 ＋ `winsec.go:167-169` 的 `preOpenLeafGuard(resolved)` 整跳删掉） | 同上 | **1** | `leafstand_112_windows_test.go:144: AC#1 RED: the standing leaf link was answered by some other failure: open C:\...\tree\leaf-link-112r1.bin...` |
| **还原绿** | 同上 | **0** | 4/4 PASS |

**还原核**（先问「谁算的、什么算法、几位」）：`sha256sum` ＝ **SHA-256、64 位十六进制**；⛔ 与 `git hash-object`（SHA-1、40 位）混比。两枚被变异文件的**还原后**读数与 `git show HEAD:<path> | sha256sum` 逐字相同：

- `internal/winsec/placement_windows.go` `5956c679de6fb62b627d04da79b37835a5435ad77d57b29dc6718fb8c8cf7d64`
- `internal/winsec/winsec.go` `6bed1b5bfae9bf37f7b5312e7040a290a874246406f5eb4756a1b42e75a6e52e`
- ⚠ 中途踩过一次坑并已改正：M2 之后我先跑了一次「还原」，可那时树里是 **M2 的产物**（`6ee1d7d6…`／`de673122…`）⇒ 两次再摘再还原之后才是上表那对值。**「还原核」必须在最后一发变异之后量**。
- 工作树侧：`git status --porcelain -- internal/winsec/placement_windows.go internal/winsec/winsec.go` 为空、`git diff --stat` 为空 ⇒ **产码一字节未动**。

## 9. 门禁 rc（共享工作树，`GOFLAGS=`、只 `unset GOTMPDIR`）

- `gofumpt.exe -l internal/winsec/leafstand_112_windows_test.go` ⇒ **空**（rc=0；已 `-w` 过一次）
- `sh scripts/d22scan.sh` ⇒ **rc=0**，逐字 `d22scan: clean - no D22 ban violations`，`ban #8 internal/ examined 635 Go files, comments and _test.go included`
- `go vet ./internal/winsec/` ⇒ **rc=0**
- `go test ./internal/winsec/ -count=1 -v` ⇒ **rc=1**，`RUN=127 PASS=75 FAIL=1 SKIP=0`，红名 `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（**票 342 那族、本腿之前就有**；同一晚三发的红名各不相同，见第 4 节）⇒ 本腿**没有新造红名**，也**没有**把任何红名变绿
- `go test ./internal/winsec/ -run 112r1` ⇒ **rc=0**（本腿四枚）

## 10. 未达成 ＋ 谁该接

- **票 112 `AC#5` 未达成**：这一步今天仍是 `failure`。⛔ 拿本地绿替代。按票面第 43 行的规矩：**欠编排者 push 后复跑**；且即便 push，`342` 未落地前它仍会红。
- 要让这道门**有确定结论**，只有两条路，都**越出本腿写面**：①**票 342 落地**（把 `no leaf leg answered` 那族变成确定性判据）；②**票 346 的甲或乙**——甲＝在 `ci.yml` 里给这些竞态仪器**安静的台面**（要机主批 `.github/**`），乙＝在 `winsec.go` 加**测试注入点**把抢窗口改成确定性交错（票 346 `AC#5` 明文：须另票＋人工批准）。本腿**两项都没做**，只把「Windows 这两枚弃权仪器今天无主」登记在此。
- 共树：`cmd/**`、`internal/risk/**`、`internal/agent/**`、`scripts/**`、`.github/**` **零改动**；工作树里别人未提交的 `cmd/wisp/run.go`／`internal/agent/loop.go`／`docs/reports/pending-and-issues.md` 我**未 add、未碰**（本笔 `add`/`commit` 两处 pathspec 只点本件与新件两枚字面路径）。⛔ push。
