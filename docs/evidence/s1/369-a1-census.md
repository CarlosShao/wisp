# 369-a1 普查腿（只读）——票 346 的射程只到 POSIX：Windows runner 上那两枚的归属与颜色

- 腿：`369-a1`＝**只读**（⛔ 落任何 `.go`／⛔ 改票面·台账·`Status:`／⛔ push／⛔ 翻框）。做的是票 369 `AC#0`。
- 锚：`HEAD=127784149ea3b6c45896f88caeb7b5384f681fdd`（2026-10-23 20:30 +0800，`git log -1` 现跑）。
- 台面：本机 `go1.27.1 windows/amd64`；`GOFLAGS=`；**只** `unset GOTMPDIR`（`TMPDIR/TEMP/TMP` 一字未动，`TEMP=C:\Users\swq\AppData\Local\Temp`）。
- 写面：本件＋`D:\tmp\wisp\20261024\369a1\**`（原始件逐枚在盘，下列每条读数都指回其中一枚）。
- 与 runner 树的关系：CI 两发都在 `23b663c2`（本机 HEAD 的**祖先**：`git log HEAD..23b663c2` 空、`git merge-base --is-ancestor 23b663c2 HEAD` rc=0）。**被测四枚文件的 blob 与 runner 逐枚同串**；`winsec.go` 差 `23/2`（新增 `removeRefusedLeaf` 接缝 var，默认值＝原 `unlinkPinnedLeaf`）⇒ 行为等价、**分母不同**（见 §5.4）。

## 0. 一句话

★**票 346 的射程是两枚 `//go:build !windows` 的弃权仪器**（现跑复认、未漂）⇒ ★**它们不进入 Windows 测试二进制**（`go list` 机械判据）⇒ ★**Windows runner 上那两枚，票 346 改完仍按胜率抛硬币**（§4 三条判据）。
★**但票面前提有两处要收窄**（§9）：①"今天红"是**逐发**的——同一 `23b663c2` 的 22:34 那一发这两枚**绿**、winsec 步 `success`；②"342 未结 ⇒ 门不可能绿"严格读**不成立**（同一反例），应写"不可能**稳定**绿"。

## 1. 三条前提逐条复核（⛔ 抄票面）

### 1.1 前提 A：票 346 射程＝两枚 POSIX 文件的弃权那一支

- 票 346 `:4` 逐字：`射程＝internal/winsec/openleaf_322_other_test.go ⇄ internal/winsec/posix_seal_window_120_other_test.go 的弃权那一支` ⇒ **现读成立**。
- 弃权支现址（`git show HEAD:<path>` 逐行读）：
  - `openleaf_322_other_test.go:205` `if liveness == 0 || flips == 0 {` ⇒ 红句 `:206`。
  - `posix_seal_window_120_other_test.go:231` `if totals.refused == 0 || totals.sealedOwn == 0 {` ⇒ 红句 `:233`；同族另有 `:223`／红句 `:225`。
- ★**"不止两枚"具名上报（§9 第①条）**：同一形状在 `!windows` 文件里现量 **4 处／3 文件**（＋`sealhandle_privateroute_120_other_test.go:212`）；在 `//go:build windows` 文件里现量 **7 处／6 文件**（runner 树 `23b663c2` 内）／**8 处／7 文件**（本机 HEAD，多 `deleteseam_348r2_windows_test.go:131`）。尺＝逐枚读 `t.Errorf|Fatalf("…RED instrument` **前四行内的守卫条件**；逐名见 §9 第①条。

### 1.2 前提 B：今天红在 Windows runner 上的那两枚，红因＝"仪器什么都没量到"

- 我**自己取**的 job log（`gh run view 39884431726 --job 120189170762 --log`，落 `D:\tmp\…\20/21-*.log`）逐字：
  - `ancestorunlink_348_windows_test.go:131: AC#2 RED instrument: no leaf leg answered (guard=0, landing's re-check=0) while reparseAt proved 63 flipped leaf cycles across 400 (flipErrs=0, first "")`
  - `openleaf_guard_322_windows_test.go:241: AC#3 RED instrument: 400 raced attempts produced guard-leg refusals=0 and racer flips=411 (flipErrs=0, firstFlipErr="")`
  - 同 job step 4 四数：`RUN=122 PASS=69 FAIL=2 SKIP=0`、结果行 `FAIL … 10.834s`、`winsec-tests.sh: gate is RED (rc=1)`。
- ★**同 SHA 的反例**（`gh run view 40034438521 --job 120690332367 --log`，head 也是 `23b663c2`）：step 4＝**success**、四数 `RUN=122 PASS=71 FAIL=0`、两枚 `--- PASS`；322 那枚读数 `guardRefusals=1`，342 那枚 split `preOpenLeafGuard=1, re-check=0`。
  ⇒ ★**0 命中 ⇒ 红、1 命中 ⇒ 绿**＝**胜率**，本腿现量坐实（⛔ 仅采信 `112-r1`）。
- ★`112-r1` 两处引数偏差（不影响本票前提，但影响它的读数可信度）：①它引 `proved 110 flipped leaf cycles`，同一 job log 实测 **63**（尺：`Select-String 'flipped leaf cycles'` 命中 **1** 行）；②它把 `288/400` 标在 `ancestoropen_342_windows_test.go:48`，实测 288 在 `ancestorunlink_348_windows_test.go:94`（`:48` 的两形是 375／281）。

### 1.3 前提 C：本机⇄runner 的具名差异（"哪条腿先答"的速率差）

| 量 | 本机（本腿现跑） | runner（旧发 39884431726） | 判 |
|---|---|---|---|
| 322 叶子守卫命中（attempts=400） | `guardRefusals=11/2/8/0/1`（`-count=5`）＋`7`（`-count=1`）＋`13`（整包） | **0**（旧发）／**1**（新发） | 方向成立；**逐发散布**，本机也能出 0 |
| C26 稳态拒绝（`ancestorunlink_348:94`，symlink-dir 形） | **174/400** | **288/400** | 方向成立（44% ⇄ 72%）；`112-r1` 报本机 205＝同量级 |
| 342 叶子守卫 split（`ancestorunlink_348:136`） | `(0,0)`／`(3,1)`／`(0,1)`／`(4,0)`／`(0,1)` | `(0,0)`（旧发）／`(1,0)`（新发） | 同上 |

### 1.4 行号漂没漂（逐条）

| 票面引 | 现值（HEAD blob） | 判 |
|---|---|---|
| 346:4 两枚文件 | 同名同 blob（`dfaddac8…`／`0b0e5173…`） | **未漂** |
| 346:19 `openleaf_322_other_test.go:206` | `:206` 逐字同句 | **未漂** |
| 346:20 `posix_seal_window_120_other_test.go:370` | `:370`＝`judgeRun120(t, totals, flips, firstFlip)` **调用点**；红句在 `:233`／`:225` | ★**不是腐烂、是尺差**（调用点 vs 句址）；该文件 blob 自 `6a6b1b18d`（2026-10-15）未变 |
| 369:38 `winsec.go:167` | `:167` 逐字 `if err := preOpenLeafGuard(resolved); err != nil {`（`:163` resolveString、`:179` pinnedParent） | **未漂** |
| 369:13/14 两红句 | `ancestorunlink_348_windows_test.go:131`／`openleaf_guard_322_windows_test.go:241` | **未漂** |

## 2. 票 346 射程那两枚的 build tag（逐枚现跑）

- `internal/winsec/openleaf_322_other_test.go:1` ⇒ **`//go:build !windows`**（214 行，`package winsec`）。
- `internal/winsec/posix_seal_window_120_other_test.go:1` ⇒ **`//go:build !windows`**（486 行，`package winsec_test`＝**外部测试包**）。
- ★尺的陷阱（⛔ 拿 `grep -l` 直接当数）：`grep -l '//go:build !windows' internal/winsec/*.go` ⇒ **18** 命中，但**真 tag 只有 15 枚**；3 枚假阳性＝`pinnedleaf_342_fallback.go:1` 是 `!windows && !unix`（前缀命中）＋`pinnedleaf_342_windows.go:9`／`tempdir_resolved_124_windows_test.go:19` 是**注释里提到**（⛔ 代码）。尺＝只取**首行**且 `^//go:build`。
- 机械判据（⛔ 推断）：`go list -f '{{range .TestGoFiles}}…'` 在 windows ⇒ 两枚 **0 命中**；`GOOS=linux` ⇒ 322 在 `TestGoFiles`、120 在 `XTestGoFiles`。

## 3. Windows 上真正在跑的那两枚（build tag ＋ CI 位置）

- `internal/winsec/ancestorunlink_348_windows_test.go:1` ⇒ **`//go:build windows`**；用例 `:107` `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`；弃权守卫 `:130`、红句 `:131`。
- `internal/winsec/openleaf_guard_322_windows_test.go:1` ⇒ **`//go:build windows`**；用例 `:158` `TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1`；弃权守卫 `:240`、红句 `:241`。
- CI 位置（三条尺逐条现跑）：`.github/workflows/ci.yml:649` job `test-windows` → `:660` step 名 `Windows ACL sealing gate (internal/winsec's own tests, ticket 110)` → `:701 run: bash scripts/winsec-tests.sh`；`gh api repos/CarlosShao/wisp/actions/jobs/<id>` 逐条读 ⇒ **job 内 step 号＝4**（1 Set up／2 checkout／3 setup-go／4 本步），是该 job **第一枚 `run:` 步**（job 共 8 个 `- name:` 步）。
- 步的射程＝**整包**：`scripts/winsec-tests.sh:38 target=./internal/winsec/`、`:206` 委派 `scripts/portable-tests.sh`（无 `-run` 过滤）⇒ 这两枚**逐发都跑**，红/绿直接决定 step 4 的颜色。

## 4. 判语：票 346 改完，Windows 上那两枚还红吗

★**不会变绿；红/绿分布一字不变。** 判据三条（都可现跑，⛔ 猜）：

1. **构建约束（机械）**：两枚 POSIX 文件在 windows 的 `TestGoFiles`＋`XTestGoFiles` 里 **0 命中**（§2 尾）⇒ 346 的任何字节改动**不进入** Windows 测试二进制；反过来，Windows 两枚在 `GOOS=linux` 下同样 0 命中。
2. **红句归属**：Windows 两枚的红句住在**各自文件**（`:131`／`:241`），346 射程两枚文件里 **0 命中**（正控：在各自文件里各命中 1）⇒ 346 就算把两枚 POSIX 的弃权支整段重写，也改不到这两句的触发条件（`st.legs[...]`／`guardSpoke` 都是 Windows 文件里的计数器）。
3. **346 三形的落点**：甲＝改 **CI 那一步**——但那两枚 POSIX 仪器跑在 **ubuntu 的 `test-core`**（`scripts/portable-tests.sh:324` 的 `core` scope 含 `"$winsec_scope"`＝`./internal/winsec/...`），**⛔ 是 test-windows 的 step 4**；乙＝产码注入点（346 自己写"须另票＋人工批准"，且 `AC#5` 明文⛔ 动 `winsec.go`）；丙＝只在两枚 POSIX 文件里打读数。**三形都不触 Windows 那两枚。**

⇒ ★**本票不该关**（第②条"未定义即停"的判据：Windows 那两枚**不能**被 346 覆盖）。

## 5. 正控：在本机跑 Windows 上那两枚

台面：本机 `go1.27.1 windows/amd64`，`GOFLAGS=`，只 `unset GOTMPDIR`；每发 `-timeout 30m`；原始件 `10/11/12/30-*.txt`。

1. `-count=1` 两枚：**都 PASS**；读数 `guardRefusals=7`（322）、`preOpenLeafGuard=0, re-check=1`（342）。
2. `-count=5`（共 10 发）：**`RUN=10 PASS=8 FAIL=2`** ⇒ ★**两枚各自在本机红过**：
   - 342 那枚 **shot 1 红**，红句逐字 `guard=0, landing's re-check=0`（`provedReparse=72`）；
   - 322 那枚 **shot 5 红**，红句逐字 `guard-leg refusals=0 and racer flips=118`；
   - 其余读数 `guardRefusals=11/2/8/1`、split `(3,1)/(0,1)/(4,0)/(0,1)`。
   ⇒ ★**"guard=0" 与 runner 那两句红句，在本机就能看见**（⛔ 需要 runner、⛔ 是环境差）。
3. 整包 `-count=1`：**`RUN=127 PASS=75 FAIL=1 SKIP=0` rc=1**；红名＝`TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（`leafrace_342r2_windows_test.go:359` 的**归因红**：6 发答在名册不具名的腿上）⇒ ★**红名漂移**在本机也可见，且是**第三枚红名、第四种红形**。
4. **分母核对**：本机 `127` ⇄ runner `122`，差 **5** ＝ 本机多两枚文件（`leafstand_112_windows_test.go` 4 枚＋`deleteseam_348r2_windows_test.go` 1 枚，均不在 `23b663c2`）；被测四枚文件 blob 与 runner **逐枚同串**；`winsec.go` 差 `23/2`（`removeRefusedLeaf` var，默认＝原函数）⇒ 本机读数与 runner **可比**，但**分母不同**（⛔ 相减当同一把尺）。

## 6. 票 342 的状态（本票要不要等）

- `git ls-files '.scratch/wisp/issues/342-*'` ⇒ **1 枚**：`342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md`，**无 `-done`**；`Status:` 逐字 `**已立，未派**` ⇒ ★**未结**（票面前提成立）。
- 它的探针**在门里**：包内 `342` 具名顶层用例 **10 枚**（`ancestoropen_342_{unix,windows}` 2／`ancestorunlink_348_{unix,windows}` 4／`ancestorwindow_342r2_windows` 3／`leafrace_342r2_windows` 1）；step 4 跑整包 ⇒ 它们逐发都跑（本机整包那枚红名就是 `342r2`）。
- ★★**但"门在 342 落地前不可能绿"严格读不成立**：同一 `23b663c2` 的 run `40034438521` step 4＝**success**（§1.2）。⇒ 口径应收窄成**"不可能稳定绿；红是胜率"**。
- ★**本票要不要等 342**：★**不必等**——342 落地改的是**祖先支／删除支**，不改叶子守卫"轮不轮得到"；但 342 落地会改同包的**分母与腿名册** ⇒ 若 369 落地腿要在 342 之后量，必须重量，⛔ 复用本件读数。

## 7. 门禁 rc

（本件写完后现跑；命令逐条、rc 逐条，⛔ 走管道取 rc）

- `sh scripts/d22scan.sh` ⇒ **rc=0**，逐字 `d22scan: clean - no D22 ban violations`（`ban #8 internal/ examined 639 Go files`／`cmd/ 150`；本件是 `.md`、⛔ 在射程内）。原始件 `40-d22scan.txt`／`.rc`。
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0**，逐字 `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（`roster entries=61`）。原始件 `41-pathlen.txt`／`.rc`。
- `git diff --cached --name-status`（提交前核索引）⇒ **逐行只有一行**：`A	docs/evidence/s1/369-a1-census.md`（⛔ 出现本笔不该有的 `D`／`A`；工作树里别人的 `M`／`??` 一枚未 add）。
- 说明：本腿写面只有一枚 `.md`；`gofumpt` 对本件不适用（⛔ `.go`，也⛔ 用 `gofumpt -l` 判 `HEAD` 的换行）。

## 8. 提交号与 `.go` 枚数

- 提交：`ed7ecc42001194b6957c28a71b5d906316251a3d`（首笔＝本件全部正文；**本笔只补录该号**，⛔ 动别的字）。
- `git show --name-only --format=%H ed7ecc420` ⇒ 1 枚：`docs/evidence/s1/369-a1-census.md`；**`.go`＝0 枚**。
- 纪律：`add`／`commit` 两处 pathspec 逐枚手写；⛔ `add -A`／`.`；⛔ amend／reset／rebase／stash／`checkout .`／`clean`；⛔ push；提交正文用 `-F <真临时文件>`（⛔ 反引号、⛔ `$()`）。

## 9. 未定义即停：具名上报（四条）

1. **"票 346 射程那两枚其实不止两枚"** ⇒ 射程＝2 枚文件**成立**；但**同族弃权支**：`!windows` **4 处／3 文件**（`openleaf_322_other:206`／`posix_seal_window_120_other:225`／`:233`／`sealhandle_privateroute_120_other:212`）、`windows` **7 处／6 文件**（`openleaf_guard_322_windows:241`／`ancestorunlink_348_windows:101`／`:131`／`ancestoropen_342_windows:56`／`ancestorwindow_342r2_windows:366`／`leafrace_342r2_windows:364`／`nametopup_353_windows:124`）／本机 HEAD 多 `deleteseam_348r2_windows:131` ⇒ ★Windows 侧无主的**不止那两枚**；本票若只治 2 枚，同族其余 5 枚仍会逐发红。
2. **"Windows 那两枚能被 346 覆盖"** ⇒ **不成立**（§4）⇒ 本票**该开、不该关**。
3. **"票 342 已结"** ⇒ **不成立**（无 `-done`）；但"门不可能绿"这句要按 §6 收窄。
4. **归属口径**：`348-…-done.md:44` 已把这句 `no leaf leg answered (guard=0, landing's re-check=0)` 记成**"票 348 的仪器地板红"**（并注明 `deleted=0`、⛔ 损伤红）⇒ ★"无主"应写成**"无开着的票"**：348 已结、346 射程外、`359-…-done.md:56` 只把它当读数。

---

**本件射程**：只报读数与判据；⛔ 选形（`AC#1` 归编排者）、⛔ 动任何产码字节、⛔ 翻框、⛔ 改票名／`Status:`／`-done`。
