# 71 — 对抗验收（AC#7，非实现者）：七格 1:1 裁决表

- **票**：`.scratch/wisp/issues/71-gates-must-self-report.md`（AC#1–AC#7）
- **本腿**：审计腿，**非实现者**（实现是 `4ce5c86ff` / `5532436f1` 两条腿；本文件不改它们任何一字）
- **开工 HEAD**：`6aa244b67`（`dev`，与 `origin/dev` 同点）
- **写面**：`docs/evidence/s1/71-ac7-adversarial.md` ＋ `.scratch/wisp/probes/71ac7/**`（探针日志与二进制）
- **本文件是**：AC#7 要的那张表本体。每行五列，读数全部**本腿现跑**（不抄票面、不抄实现者自账）。README 规则 6：缺行即 FAIL。

## 0. 结论先给

**七格闭合 5／7**：AC#1、#2、#3、#4、#7 **满足**；AC#5、#6 **部分**；**0 格不满足**。

**那发「真实的红」是仓外导出树里的红，不是 GitHub runner 上的红。** 三条独立读数（§4 全展开）：
1. 红的现场是 `/tmp/wisp-71-export-1563`——`git archive 4ce5c86ff` 导出的**仓外**树（`tools/d22scan/export_red_proof.sh:105-111` 明文拒绝仓内目录），本腿复跑 `export_red_proof.sh -r 4ce5c86ff` ⇒ `RESULT PASS`、`20124 tracked paths, 0 missing`、`bans #1-5 internal/ 404 → 405 → 404`、整树 sha256 `14d07cc4…cb5975` 三跑同值。
2. **GitHub 上那个 revision 一次都没跑过**：`gh api repos/CarlosShao/wisp/commits/4ce5c86ff…/check-runs` ⇒ `total_count = 0`。没有 run，就没有 runner 上的红。
3. 反过来，**含这条改动的已推送 commit 在 runner 上是绿的**：`gh run view 43591051106` ⇒ head `445ea5b98`（`git merge-base --is-ancestor 4ce5c86ff 445ea5b98` = true），6 个 job 全 `success`，其中 `lint` = `success`。

⇒ AC#5 的**接线**成立，**「一次真实的红」这一半只到「CI 那条命令在仓外导出树上红」**，不到「CI 红的那一行」。按台账 **A1868** 的规矩（说 CI 必须给 run 号），本文件**不把导出树红说成 CI 红**。

**这票能不能结**：**不能。** 剩两格都是「部分」，且都能一句话翻，但都**不在本腿写面内**：
- AC#5：要么由有 push 权的编排者推一次真违规并贴 run/job 行（AC 原文要的），要么由 owner 把 AC#5 的文字改成「仓外导出树跑 CI 那条命令」——那是**改契约文字**，只有 owner 能批。
- AC#6：`A30⑤` 自己的文字没被改，且 `:796-798` 仍写 `(c) … 未闭`；那一句**今天已失实**（balldebug 那半张票 2026-11-05 已交，本腿 8/8 PASS、rc=0/1/2 全量复跑）。补一句即可，属台账追加。

## 1. 读数环境（两次会咬人的树状态，先说清）

1. **工作树不是 HEAD**：活体 `6aa244b67` ＋ 其他代理**未跟踪/未提交**的文件。19:53 那次 `runtests.sh -C tools/d22scan ./...` 是 **PASS=44 FAIL=2 / rc=1**，两条红都指向别人在飞的未跟踪件 `cmd/wisp/panel_preview_192.go:80`（`phantom-citation`）；同一文件 4 分钟后被其作者改掉，19:57 的扫描即 `clean`、20:07 的整包复跑 **PASS=46 FAIL=0**。⇒ **引用本票任何「整包全绿」前必须先看当时的 `git status`**，否则会把别人的在飞件记到这张票头上。
2. **两种树的数不同不是矛盾**：工作树扫描跳 gitignore（`frontend/dist/` 等），导出树无 `.git` ⇒ 扫描器在 stderr 大声打 `gitignore rules NOT APPLIED`（A207 的响亮方向）并全量扫。故下表两列数各自成立：工作树 `405 / 79 / 123 / 51 / 21 / 39 / 123 / 1016 / 285`，导出树 `404 / 78 / 122 / 51 / 21 / 14 / 122 / 1012 / 283`。

## 2. 1:1 裁决表（README 规则 6，七格 = 七行）

| AC 编号 | 票面原文要求（逐字摘） | 盘上读数（file:line ＋本腿现跑数字） | 我判 | 依据一句话 |
|---|---|---|---|---|
| **AC#1** | `cmd/balldebug -diff` 有一个**阈值门**模式——当任一行 `px_delta_ge8` 低于给定的最小可见性阈值、或缺行、或 `timers` 非 `no`（在要求零定时器的态上）时**非零退出**；阈值**写死在票面/配置里**并带来源（D43/SPEC-08），**不许为了让某次跑变绿而现场调低**（D22 精神）。 | 门与来源：`cmd/balldebug/diffgate_windows.go:57` `diffGateFloorPx8 = 800`、`:62` 来源整句、`:74-83` 低于 800 直接拒绝（不是夹紧）。本腿现跑真二进制（`go build ./cmd/balldebug` ⇒ rc=0）：坏表 `62-diff-baseline`（`px_delta_ge8=0`）⇒ **rc=1** `1 violation(s)`；`timers=yes` 的 Sleeping 行 ⇒ **rc=1**；重复行 ⇒ **rc=1** `2 rows for one state`;合格表 `62-diff-signoff` ⇒ **rc=0** `PASS … 6 rows`;`-diff-gate-min-px8 10` ⇒ **rc=2** 并打印来源；`-diff-gate-min-px8 100000` ⇒ **rc=1** 逐行 6 条 FAIL。测试 `sh tools/d22scan/runtests.sh ./cmd/balldebug/` ⇒ **rc=0，PASS=8 FAIL=0 SKIP=0**。CI 覆盖：`.github/workflows/ci.yml:1350` 跑 `portable-tests.sh --scope=windows`，包清单 `scripts/portable-tests.sh:456` 含 `./cmd/balldebug/`。 | **满足** | 三种构造（低像素／零定时器态出定时器／重复行）都真退出 1，合格表退 0，调低被 2 拒，阈值和来源都在代码里而不是在报告里。 |
| **AC#2** | 阳性对照（**这条不做本票不算完**）：构造一次**故意不成像**的跑（把球隐藏或把阈值调到不可能达到），证明 AC#1 的门**会红**。今天的教训正是"没红"被当成"过"。 | 本腿两发真进程读数（不是读测试名）：① 票 62 记录的不成像真表 `docs/evidence/s1/62-diff-baseline/diff-table.txt`（`px_delta_ge8=0`）⇒ **rc=1**，`gate: FAIL 01: Sleeping: px_delta_ge8=0 < 800`；② 把阈值抬到不可能 `-diff-gate-min-px8 100000` 跑被接受过的 `62-diff-signoff` 表 ⇒ **rc=1**，6 行逐行 FAIL（2120/1620/4856/3230/4848/3220 全 < 100000）。单元缝同向：`cmd/balldebug/diffgate_windows_test.go:26 TestDiffGateRedOnANonImagingRun` PASS（含在 8/8 里）。 | **满足** | AC#2 明列的两种构造都各有一发 rc=1，且都是本腿用真二进制、真退出码跑出来的。 |
| **AC#3** | 一个可复用的测试运行器（脚本或 `cmd/` 下的小工具），跑 `go test -tags <X> -v <pkg>` 后**断言 `--- PASS/FAIL/SKIP` 行数 > 0 且 `-run` 模式实际匹配到测试**；匹配 0 条时**非零退出**并打印它匹配到了什么。**不得**把 SKIP 计入"通过"。⇒ 落地物 `tools/d22scan/runtests.sh`。CI 的三处接入：`lint` 的阳性对照步、`test-core` 的 `TestLayoutForTestEnv`、`test-windows` 的 `TestPathResolverJunctionWindows`。 | 断言本体：`tools/d22scan/runtests.sh:221-225`（数 `^=== RUN`/`^--- PASS`/`^--- FAIL`/`^--- SKIP`/`[no tests to run]`）、`:234-239`（SKIP 一律 rc=1）、`:240-244`（0 顶层结果 = 匹配空 ⇒ rc=1）、`:229-232`（真失败原样传码）。本腿三发真跑：① `runtests.sh ./internal/proc/ -run TestLayoutForTestEnv` ⇒ **rc=0** `PASS=1 FAIL=0 SKIP=0, === RUN=3`；② 漂移名 `-run TestLayoutForTestEnvRenamedAway` ⇒ **rc=1** `PASS=0 FAIL=0 SKIP=0, === RUN=0, '[no tests to run]'=1`（同一条命令不经脚本时 go 自己打印 `PASS`/`ok`）；③ 本腿自建 SKIP fixture（`.scratch/wisp/probes/71ac7/skipfixture`）⇒ **rc=1** `SKIP=1` 并打 `SKIP is not a pass`。三处接线：`ci.yml:119`（lint 阳性对照）、`ci.yml:798`（test-core）、`ci.yml:1375`（test-windows），三行都是 `bash tools/d22scan/runtests.sh`。 | **满足** | 运行器的三条判据各有本腿一发真退出码（0／1／1），三处 CI 调用点逐处有行号。 |
| **AC#4** | `d22scan` **逐作用域**自报 `examined N files`（含 `frontend/` 这类子作用域），任何**声明了但恒为 0** 的作用域 ⇒ 致命退出；或者**删掉那个作用域**并在 allowlist 留一行说明。现存的 `allowlist.txt` 只允许**变短或不变**，变长需我在票面批。 | 自报本体：`tools/d22scan/main.go:428 declaredScopes`、`:544` 那行 `d22scan: scope %-24s examined %3d %s`、`:494 emptyLiveScope`、`:379 undeclaredKeys`、`:514 driftedAbsentScope`。本腿现跑：工作树 `go run . -root <abs>` ⇒ **rc=0，9 条作用域行**（`405 / 79 / 123 / 51 / 21 / 39 / 123 / 1016 / 285`，无 `[NOT COVERED]`）；导出树（`4ce5c86ff`）同一条命令 ⇒ **rc=0，9 行**（`404 / 78 / 122 / 51 / 21 / 14 / 122 / 1012 / 283`）。「不是常量」的证明：`tools/d22scan/selfreport_delta_test.go:47-114`（9 个用例，每个作用域一枚、路径互不相同）、`:150-171`（用例表 ↔ `declaredScopes()` 双向覆盖守卫）、`:196-214`（同一条作用域行 n0→n0+1、当次 rc 0→1、红句点名）；本腿 `runtests.sh -C tools/d22scan ./... -run TestEveryDeclaredScopeCountsTheFileItIsGiven` ⇒ **rc=0，9 个子测试全 PASS**，日志逐条 `examined 3->4 / 1->2 / …`、`rc 0 -> 1`。0 文件作用域致命：`TestBuiltBinaryGoesRedEndToEnd/empty_live_scope_exits_2` 在导出树跑里给出 **rc=2**（`.scratch/wisp/probes/71ac7/logs/10-export-logs-excerpt.txt` 行 191）。allowlist：`8d56a948d` / `4ce5c86ff` / `HEAD` 三处条目数均 **5**，且 `4ce5c86ff` 只动 3 个文件、**没碰** allowlist。 | **满足** | 9 条作用域各报数、每条都随「塞一个坏文件」+1 且当次红、0 文件活作用域 rc=2、allowlist 5→5。 |
| **AC#5** | 把上面三处接进 CI（`lint` 或 `test-windows`），并给出**一次真实的红**：在分支上故意造一条不合格产物、贴出 CI 红的那一行，然后回滚。**只贴"本地能红"不算**。 | **接线成立**：`ci.yml:119`／`:798`／`:1375` 三处 `runtests.sh`；扫描步 `ci.yml:174` 在阳性对照 `:119` **之后**（次序对）。**红存在但不是 runner 红**：本腿复跑 `sh tools/d22scan/export_red_proof.sh -r 4ce5c86ff -k` ⇒ **RESULT PASS**、`20124 tracked paths, 0 missing`、基线 **rc=0**（control `PASS=46`）、塞 `internal/zz71export/leak.go` 后 **rc=1**（`bans #1-5 internal/` **404→405**，红句 `internal/zz71export/leak.go:5: [bare-goroutine] …`，control `PASS=44 FAIL=2`，包装器 `d22scan.sh: FAIL - control rc=1, scan rc=1`）、删除后 **rc=0** 且整树 sha256 `14d07cc4…cb5975` 与基线**逐字相同**、数回 404。红的现场是 `/tmp/wisp-71-export-1563`（仓外导出），**不是** GitHub runner：该 revision `check-runs` = **0**；含它的 runner 发 `43591051106`（head `445ea5b98`）**lint=success**。 | **部分** | 接线与「可复跑的、点名文件的红、还原后同 hash」都在；但 AC# 字面要的是「贴出 CI 红的那一行 / 只贴本地能红不算」，这发红跑在仓外导出树上，不等于 CI 红（A1868 规矩），故只到部分。 |
| **AC#6** | 本票结束时，`docs/reports/pending-and-issues.md` 的 **A30⑤** 追加一句"已闭于票 71"，并把本票新挖到的同类空仪器（如果有）继续登记在 A30 之后，**不要改已完成票的证据文件**。 | 追加的事实存在：`docs/reports/pending-and-issues.md:37251`「**A1870（11-14 19:5x）：A30⑤ 已闭于票 71**」，`:37257` 也写了未闭两格；后半句（新挖空仪器继续登记）也在 `:37259`。**但**：AC 点名的落点（`[A30]` 那条的 ⑤ 项，台账记法 `A30⑤`）本体在 `:781`，其 `:786` 的标题仍是「🟡 ⑤ **部分闭于票 71**（2026-09-21 …）」，且 `:796-798` 仍写「(c) **`cmd/balldebug -diff` 恒返回 nil：未闭**」——这一句今天**已失实**（该半张票 2026-11-05 交于 `08c4c51e2`，本腿 `runtests.sh ./cmd/balldebug/` **8/8 PASS rc=0**、真二进制 rc=1/0/2 全有）。A1870 自称「不编辑 A30 那条老项目符号」，即 AC 字面点的那个位置一字未动。 | **部分** | 闭案句写在了 A1870 并指向 A30⑤，但 AC 点名的 A30⑤ 本体没改、且仍写着「未闭」，读到那一行的人被告知的反而是旧事实。 |
| **AC#7** | 对抗验收由非实现者执行，验收报告含与上述 AC **1:1 的裁决表**（README 规则 6），每行给出可复现命令与真实 exit code。 | 本文件即交付物：七行五列表在 §2，逐格命令与真实退出码在 §3（全部本腿现跑）；执行者与实现者不同人（实现 `4ce5c86ff`/`5532436f1`，本腿只写这一个件与探针目录）。 | **满足** | 七格齐全、每格有命令与真实 rc、判者 ≠ 实现者。 |

**七格汇总**：满足 **5**（#1 #2 #3 #4 #7）／部分 **2**（#5 #6）／不满足 **0**。

## 3. 逐格复跑命令与真实 exit code（AC#7 的「每行命令」）

以下全部在本腿会话内现跑，日志落在 `.scratch/wisp/probes/71ac7/logs/`。

| AC | 命令（在仓根执行） | 本腿读数 | 日志 |
|---|---|---|---|
| #1 | `go build -o .scratch/wisp/probes/71ac7/bin/balldebug.exe ./cmd/balldebug` | rc=0 | `06-balldebug-build.txt` |
| #1 | `./.scratch/wisp/probes/71ac7/bin/balldebug.exe -diff-gate-table docs/evidence/s1/62-diff-baseline/diff-table.txt` | **rc=1**（`1 violation(s)`，`Sleeping: px_delta_ge8=0 < 800`） | `07-ac1-ac2-process.txt` |
| #1 | `… -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt` | **rc=0**（`PASS`，`examined 6 rows`） | 同上 |
| #1 | `… -diff-gate-min-px8 10 -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt` | **rc=2**（拒绝低于 floor，打印来源） | 同上 |
| #1 | `… -diff-gate-table .scratch/wisp/probes/71ac7/fixtures/timers-yes.txt` | **rc=1**（`timers=yes but SPEC-08 §2 grants this state no animation timer`） | `09-ac1-timers-dup.txt` |
| #1 | `… -diff-gate-table .scratch/wisp/probes/71ac7/fixtures/duplicate-row.txt` | **rc=1**（`2 rows for one state`） | 同上 |
| #1/#2 | `sh tools/d22scan/runtests.sh ./cmd/balldebug/` | **rc=0**，`PASS=8 FAIL=0 SKIP=0, === RUN=8` | `05-balldebug-runtests.txt` |
| #2 | `… -diff-gate-min-px8 100000 -diff-gate-table docs/evidence/s1/62-diff-signoff/diff-table.txt` | **rc=1**（6 行逐行 FAIL） | `07-ac1-ac2-process.txt` |
| #3 | `sh tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv` | **rc=0**，`PASS=1 FAIL=0 SKIP=0, === RUN=3` | `04-ac3-runtests.txt` |
| #3 | `sh tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnvRenamedAway` | **rc=1**，`PASS=0 FAIL=0 SKIP=0, === RUN=0, '[no tests to run]'=1` | 同上 |
| #3 | `sh tools/d22scan/runtests.sh -C .scratch/wisp/probes/71ac7/skipfixture ./...` | **rc=1**，`SKIP=1`（`SKIP is not a pass`） | 同上 |
| #4 | `cd tools/d22scan && go run . -root "<abs repo>"` | **rc=0**，9 条 `scope … examined N`（工作树） | `08-worktree-scan.txt` |
| #4 | `sh tools/d22scan/runtests.sh -C tools/d22scan ./... -run TestEveryDeclaredScopeCountsTheFileItIsGiven` | **rc=0**，9 个子测试 PASS，逐条 `n0 → n0+1`、`rc 0 -> 1` | `01-delta-run.txt` |
| #4 | `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**，`PASS=46 FAIL=0 SKIP=0, === RUN=103`（20:07；19:53 那次因别人在飞件为 `PASS=44 FAIL=2 / rc=1`） | `02-…txt` / `11-…txt` |
| #5 | `sh tools/d22scan/export_red_proof.sh -r 4ce5c86ff -k` | **rc=0（RESULT PASS）**：`20124/0` 条、基线 rc=0、塞入 rc=1 且点名、还原 rc=0 同 sha256 | `03-export-red-proof.txt` ＋ `10-export-logs-excerpt.txt` |
| #5 | `gh api repos/CarlosShao/wisp/commits/4ce5c86ff…/check-runs --jq .total_count` | **0**（该 revision 在 runner 上没有 run） | 会话读数（§4） |
| #5 | `gh run view 43591051106 --json headSha,conclusion,jobs` | head `445ea5b98`、6 job 全 `success`（含 `lint`）；`merge-base --is-ancestor 4ce5c86ff 445ea5b98` = true | 会话读数（§4） |
| #6 | `grep -n "A1870\|A30⑤" docs/reports/pending-and-issues.md` | `:37251`（A1870 收件）、`:781`/`:796-798`（A30⑤ 本体仍写「部分闭」「(c) 未闭」） | `12-ledger.txt` |
| #7 | 本文件 | 七行五列齐 | — |

## 4. 专门一节：那发红到底红在哪（导出树 vs runner）

三条独立读数，任一条单独成立就足以否掉「CI 红」的说法：

1. **红的现场是仓外树，且脚本自己拒绝仓内目录。** `tools/d22scan/export_red_proof.sh:104-118` 生成/校验导出目录，`:107-111` 明写「导出目录在仓内 ⇒ 直接 die」；`:130` 打 `exporting 4ce5c86ff … to /tmp/wisp-71-export-1563`；`:195` 每次跑都是 `(cd "$export_dir" && sh scripts/d22scan.sh)`。本腿复跑日志的现场行：`d22scan.sh: scan of /tmp/wisp-71-export-1563`（`10-export-logs-excerpt.txt` 行 310/325），塞入后红句在 `internal/zz71export/leak.go:5`，还原后同树同 sha、`d22scan.sh: OK - control rc=0, scan rc=0`。**GitHub 的 checkout 路径是 `/home/runner/work/...` 或 `D:\a\...`，永远不是 `/tmp/wisp-71-export-<pid>`。**
2. **那个 revision 在 runner 上没有 run。** `gh api repos/CarlosShao/wisp/commits/4ce5c86ffacfcc35ae1af37280708df9036c2a13/check-runs` ⇒ `total_count = 0`；`gh run list --commit 4ce5c86ff…` 空。没有 run ⇒ 不存在「runner 上的红」。
3. **含这条改动的已推送 commit，runner 上是绿的。** run `43591051106`（head `445ea5b98`，2026-11-14T11:39Z 完成）6 job 全 `success`，其中 `lint` = `success`；`git merge-base --is-ancestor 4ce5c86ff 445ea5b98` 返回 0（真）。`lint` 正是跑 `runtests.sh -C tools/d22scan ./...`（`:119`）与 `sh scripts/d22scan.sh`（`:174`）的作业。

**所以本文件对 AC#5 的判词只有一句**：接线 ✔；「一次真实的红」＝**仓外导出树里的红**，**不是** runner 上的红，**不按 CI 红记**。票面、台账 A1870、实现者证据 `docs/evidence/s1/71-ac45-…md:117` 三处也都是这么写的——本腿独立复现后**同意**，并附上第 2/3 条它们当时没有的 runner 侧读数。

## 5. 两格「部分」各自差什么（要什么才能翻）

| 格 | 差的东西 | 谁能做 | 翻的判据（一句话） |
|---|---|---|---|
| **AC#5** | runner 上的一发红：在分支上推一条不合格产物、贴出 CI 红的那一行、再回滚 | 有 push 权的人（本腿与实现腿都被禁 push）；或 owner 改 AC#5 文字，把「仓外导出树跑 CI 那条命令」写成判据 | 给出 `run id ＋ job 行 ＋ 该 job 的 conclusion=failure ＋ 失败步骤里的红句`；按 A1868 的三件套，缺 run 号不算 |
| **AC#6** | `A30⑤` 本体那一句：把 `:796-798` 的「(c) 未闭」改成已闭（或在 A1870 之后加一条显式推翻它的追加，并让 `A30⑤` 那行指过去） | 台账维护者（本腿写面不含台账） | `grep -n "A30⑤" docs/reports/pending-and-issues.md` 的命中里，不再有一行说 balldebug 未闭；且 `:786` 不再写「部分」 |

撤销口令（任一条要复原）：**「71 第 5/6 格按原样留着」**。

## 6. 不声称 / 已知口径问题（宁缺毋造）

1. **不声称 GitHub CI 红过**：见 §4。本文件没有任何一句把导出树红说成 CI 红。
2. **不声称 `bans #1-5 internal/` 这行名字对**：该行同时承载 ban #9（phantom-citation）的计数（`tools/d22scan/main.go` 的 `selfCoverForTag`），名字比承载的 ban 少说 #9。编排者已在票面裁定「排进 71 收尾小活」，本腿**不动**（改文案会动 9 个子测试钉住的字符串，属 `tools/d22scan/**` 写面，不是本腿的）。
3. **探针范围**：本腿的「不是常量」证明跑在 `t.TempDir()` 造的形状树上（`scan_test.go:723 liveFixture`），断言打在**内部计数器**与**打印出来的那一行**上；**进程退出码**那条路由 `TestBuiltBinaryGoesRedEndToEnd` 单独覆盖，但它不读逐作用域的数。两者合起来覆盖 AC#4 的两半，但任何一半单独都不够——记在这里，免得下一个人只看见一条。
4. **数会变**：本文件所有数字带时刻（19:5x–20:0x）。工作树是共享的，`git status` 干净 ≠ HEAD 干净（票面「编排者追加」第 3 条），引用前先用 §3 的命令重跑。
5. **没跑的东西**：`-diff` 的真桌面跑（`-diff -diff-gate` 需要的交互桌面本机 CI 环境不提供）不在本腿读数内；AC#1/#2 的进程级读数走的是 `-diff-gate-table` 这条**同评测器**的入口（`diffgate_windows.go:15-22` 说明它存在的理由）。

## 7. 三把尺（本腿收件人可复跑）

```
# 1. 七格里最关键的那条：自报数不是常量（9 个子测试，rc=0）
sh tools/d22scan/runtests.sh -C tools/d22scan ./... -run TestEveryDeclaredScopeCountsTheFileItIsGiven

# 2. 那发红的现场（仓外导出、三跑、同 hash），看完自己判断它是不是 CI 红
sh tools/d22scan/export_red_proof.sh -r 4ce5c86ff -k

# 3. 两个「部分」各差什么：runner 侧有没有 run、A30⑤ 现在怎么写
gh api repos/CarlosShao/wisp/commits/4ce5c86ffacfcc35ae1af37280708df9036c2a13/check-runs --jq .total_count
grep -n "A30⑤" docs/reports/pending-and-issues.md
```
