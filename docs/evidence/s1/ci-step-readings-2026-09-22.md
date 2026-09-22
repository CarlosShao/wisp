# CI 步级读数 — 2026-09-22（`ci-reader-s22`）

只读取数。本表所有结论都指到 **run id + job id + step 号 + step 名**，不指 `ci.yml` 行号。
取证方法沿用本仓已登记的口径（`docs/evidence/s1/111-ci-step-readings-attempt2.md` 的四判据表）。

## 0. 日志真伪四判据（每份日志逐枚贴，先立此存照）

取法：`gh api -i repos/CarlosShao/wisp/actions/jobs/<job>/logs` 读响应头，
`gh api repos/CarlosShao/wisp/actions/jobs/<job>/logs` 读 body（**不带 `-i`**，带 `-i` 时 body 会混进头块，
按 `tail -n +N` 切会切出 221–222 字节的假差值 —— 本代理第一遍就踩了这个，见 §6）。

| 日志（step 归属） | job id | HTTP | `Content-Length` | 实际 body 字节 | 行数 | 首行 |
| --- | --- | --- | --- | --- | --- | --- |
| `test-windows`（run 35608530583） | `106361546809` | `HTTP/1.1 200 OK` | 403969 | 403969 | 3006 | `2026-09-21T13:53:48.8126112Z Current runner version: '2.3…'` |
| `test-windows`（run 35609239145） | `106364544488` | `HTTP/1.1 200 OK` | 404087 | 404087 | 3005 | `2026-09-21T14:02:11.8804618Z Current runner version: '2.3…'` |
| `test-windows`（run 35616790753） | `106389418059` | `HTTP/1.1 200 OK` | 413232 | 413232 | 3044 | `2026-09-21T15:07:30.8917956Z Current runner version: '2.337.0'` |
| `test-windows`（run 35706551384） | `106676758365` | `HTTP/1.1 200 OK` | 407431 | 407431 | 3012 | `2026-09-22T08:45:37.8794368Z Current runner version: '2.3…'` |
| `test-windows`（run 35719170857） | `106717758183` | `HTTP/1.1 200 OK` | 408773 | 408773 | 3022 | `2026-09-22T11:01:49.0319440Z Current runner version: '2.3…'` |
| `test-windows`（run 35723172814） | `106730524368` | `HTTP/1.1 200 OK` | 408880 | 408880 | 3022 | `2026-09-22T11:44:50.7118746Z Current runner version: '2.3…'` |

六枚全部 `HTTP/1.1 200 OK` 且 **body 字节 == `Content-Length`**，无一份是 22 字节空 ZIP（`504b 0506`）。
每份 `##[group]Run bash scripts/winsec-tests.sh` 均在（下引行号即在其中）。

其余被引用的 15 份日志，同一法验：

| 日志（归属） | job id | HTTP | `Content-Length` | 实际 body 字节 | 行数 | 首行（截） |
| --- | --- | --- | --- | --- | --- | --- |
| `test-windows`（run 35605531736，**§3.1 分离腿**） | `106351877745` | `HTTP/1.1 200 OK` | 407564 | 407564 | 3003 | `2026-09-21T13:26:35.5147722Z Current runner version: '2.337.0'` |
| `lint`（run 35606321404，票85a **之前**） | `106355017673` | `HTTP/1.1 200 OK` | 52711 | 52711 | 554 | `2026-09-21T13:35…Z Current runner version: '2.3…'` |
| `lint`（run 35608530583） | `106361546952` | `HTTP/1.1 200 OK` | 59620 | 59620 | 618 | 同上口径 |
| `lint`（run 35609239145） | `106364544594` | `HTTP/1.1 200 OK` | 59620 | 59620 | 618 | 同上 |
| `lint`（run 35706551384） | `106676758176` | `HTTP/1.1 200 OK` | 59778 | 59778 | 620 | 同上 |
| `lint`（run 35719170857） | `106717758061` | `HTTP/1.1 200 OK` | 59710 | 59710 | 619 | 同上 |
| `lint`（run 35723172814） | `106730524758` | `HTTP/1.1 200 OK` | 59706 | 59706 | 619 | 同上 |
| `lint`（run 35737627205） | `106778845192` | `HTTP/1.1 200 OK` | 59695 | 59695 | 619 | 同上 |
| `test-windows`（run 35603195107，**票111 反例**） | `106343955813` | `HTTP/1.1 200 OK` | 234781 | 234781 | 1477 | `2026-09-21T13:0…` |
| `test-core`（run 35706551384） | `106676758315` | `HTTP/1.1 200 OK` | 517761 | 517761 | — | `2026-09-22T08:4…` |
| `test-core`（run 35737627205） | `106778845234` | `HTTP/1.1 200 OK` | 530058 | 530058 | 4035 | `2026-09-22T14:0…` |
| `slo-full`（run 35723172814） | `106730524584` | `HTTP/1.1 200 OK` | 11284 | 11284 | 107 | `2026-09-22T11:44:51.0582749Z Current runner version: '2.337.0'` |
| `slo-full`（run 35719170857，**绿对照**） | `106717758014` | `HTTP/1.1 200 OK` | 26316 | 26316 | 286 | `2026-09-22T11:01:51.1003160Z …` |
| `slo-full`（run 35737627205） | `106778845289` | `HTTP/1.1 200 OK` | 23108 | 23108 | 253 | `2026-09-22T14:03:23.8445728Z …` |
| **`slo-full`（run 35706551384）** | **`106676758319`** | **`HTTP/1.1 404 The specified blob does not exist.`** | 215 | 215 | 2 | **不是日志**：`<?xml version="1.0" encoding="utf-8"?><Error><Code>BlobNotFound</Code>…`，`X-Ms-Error-Code: BlobNotFound` ⇒ 见 §5-a |

## 1. run ↔ sha 对应表

对应关系**全部用 `gh api repos/CarlosShao/wisp/compare/<sha>...<headSha>` 判包含**，不凭时间猜。
`status=ahead / behind=0` ⇒ 该 run 的 head **含**这枚 sha；`status=behind / ahead=0` ⇒ **不含**。

| run id | headSha | createdAt (UTC) | status / conclusion | 含 `e563a61`(85a) | 含 `c6dbbf9`(115) | 含 `ab50a99`(今日最新) |
| --- | --- | --- | --- | --- | --- | --- |
| 35603195107 | `399a783` | 09-21 13:03:41 | completed / failure | ✗（behind=16）| ✗ | ✗ |
| 35603607270 | `9141d4c` | 09-21 13:07:40 | completed / failure | ✗（behind=14）| ✗ | ✗ |
| 35605531736 | `88d8956` | 09-21 13:25:37 | completed / failure | ✗（behind=5）| ✗ | ✗ |
| 35606321404 | `d3cc9ed` | 09-21 13:32:59 | completed / failure | **✗（behind=2，票85a 之前最后一枚）** | ✗ | ✗ |
| 35608530583 | `a13c6dc` | 09-21 13:53:44 | completed / **failure** | **✓（最早含 85a，ahead=5）** | **✗（behind=5，实测）** | ✗ |
| 35609239145 | `c9907bb` | 09-21 14:00:25 | completed / **failure** | ✓（ahead=7） | ✗（behind=3） | ✗ |
| **35616790753** | `d1056c0` | 09-21 15:07:27 | completed / **failure** | ✓（ahead=18） | **✓（最早含 115 修）** | ✗ |
| 35706551384 | `e475ce0` | 09-22 08:45:29 | completed / **failure** | ✓（ahead=37） | ✓（ahead=27） | ✗ |
| 35719170857 | `8ec04f3` | 09-22 11:01:43 | completed / **failure** | ✓（ahead=55） | ✓（ahead=45） | ✗ |
| 35723172814 | `09b5285` | 09-22 11:44:46 | completed / **failure** | ✓（ahead=56） | ✓（ahead=46） | ✗ |
| **35737627205** | `ab50a99` | 09-22 14:03:17 | completed / **failure**（本代理首读时是 `in_progress`，读完已收口 `updatedAt=14:11:37Z`） | ✓（ahead=75） | ✓（ahead=65） | **✓（今日最新 head）** |

被取消的 8 枚（`35601246522 / 35601648451 / 35604249114 / 35604909648 / 35605067027 / 35605359280 /
35605937530 / 35608940827`）**全部 `jobs | length == 0`** ⇒ 它们不是失败，是**从没起过 job**
（见 §5、§6）。

## 2. 四项判据逐条

### 2.1 票 115 的结案 run —— **拿到，且逐字命中**

**步名**：`Windows ACL sealing gate (internal/winsec's own tests, ticket 110)`
**run 35616790753 / job `106389418059` / step4**（`##[group]Run bash scripts/winsec-tests.sh` 在该日志第 **186** 行）

原文（行号 = 该 job 日志内绝对行号）：

```
1452:2026-09-21T15:09:11.6769329Z runtests.sh: OK - packages=[./internal/winsec/ -count=1 -skip ^(...)$] top-level: PASS=42 FAIL=0 SKIP=0, === RUN=82, '[no tests to run]'=0
1453:2026-09-21T15:09:11.8012066Z portable-tests.sh: four numbers (all from -v output): === RUN=82  --- PASS=42  --- FAIL=0  --- SKIP=0
1456:2026-09-21T15:09:12.1356475Z winsec-tests.sh: four numbers (all from -v output): === RUN=82  --- PASS=42  --- FAIL=0  --- SKIP=0
```

⇒ 交件方要的逐字形状 **`=== RUN=82 --- PASS=42 --- FAIL=0 --- SKIP=0` 成立**。
这是**最早一枚 head 含 `c6dbbf9` 的 run**（`compare/c6dbbf9...d1056c0` = `status=ahead behind=0`）。
上一枚 run `35609239145` 的 head `c9907bb` **不含** `c6dbbf9`（`status=behind ahead=0`，差 3 枚），
它的读数仍是 `PASS=40 FAIL=2`（job `106364544488` 第 **1463** 行）——不是回归，是**没带修**。

两枚用例名在该步 `--- PASS` 列里的原文行：

```
283:2026-09-21T15:09:11.3010647Z === RUN   TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree
288:2026-09-21T15:09:11.3035440Z --- PASS: TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree (0.06s)
289:2026-09-21T15:09:11.3036276Z === RUN   TestNoticeAttributionKeepsTwoTreesApart
290:2026-09-21T15:09:11.3037013Z --- PASS: TestNoticeAttributionKeepsTwoTreesApart (0.10s)
```

基线那两枚红就是它们俩（run `35608530583` / job `106361546809`）：

```
293:2026-09-21T13:55:32.1691669Z --- FAIL: TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree (0.07s)
296:2026-09-21T13:55:32.1697506Z --- FAIL: TestNoticeAttributionKeepsTwoTreesApart (0.04s)
```

**"42 有前提"这条，本代理核完了：前提在这枚 run 上成立。**
`git log a13c6dc..d1056c0 -- internal/winsec` 只有两枚（`189cb1e` 票 119、`c6dbbf9` 票 115），
其中 `189cb1e` 动的 `dataroot_symlink_119_other_test.go` 首行是 **`//go:build !windows`** ⇒ 在 windows 那一步**不进分母**；
`c6dbbf9` 动的是 `notice_attribution_115_windows_test.go`，`+func Test` 净增 **0**（改的是既有 5 处比对，不加用例）。
⇒ 该步 RUN 仍是 82，`40→42 / FAIL 2→0` **全部**来自这两枚用例转 PASS，无一枚新用例混进分子。

**但**：交件方担心的"分母被别人加过"**确实已经发生了**，只是发生在这枚结案 run **之后**——见 §3。
所以往后任何一枚 run 都**不可能**再报出 42；**票 115 的结案证据只能钉在 run `35616790753` / job `106389418059` / step4 这一格**，
不能拿今天任何一枚 run 去补。

### 2.2 票 85a `lint`/staticcheck 自报分母 —— **拿到了，形状在，但 `findings` 不是 ≈37**

**步名**：`staticcheck`（`lint` job 的 step9）
`##[group]Run go install honnef.co/go/tools/cmd/staticcheck@2026.2.1` 在场 ⇒ 钉版本这一步**真的执行了**。

自报行原文（各 run 的 `lint` job 日志内绝对行号）：

```
run 35608530583 / job 106361546952 / step9 / L596:
  staticcheck self-report: version=staticcheck 2026.2.1 (0.8.1) modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
run 35609239145 / job 106364544594 / step9 / L596:
  staticcheck self-report: version=staticcheck 2026.2.1 (0.8.1) modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
run 35706551384 / job 106676758176 / step9 / L598:  … modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
run 35719170857 / job 106717758061 / step9 / L597:  … modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
run 35723172814 / job 106730524758 / step9 / L597:  … modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
run 35737627205 / job 106778845192 / step9 / L597:  … modules=3 packages=34 findings=42 toolchain-crash-lines=0 (step exit is 1)
```

逐字段裁：

| 字段 | 期待 | 实读 | 裁 |
| --- | --- | --- | --- |
| `modules` | 3 | 3 | ✓ |
| `packages` | …（未钉） | 34 | 见下方分模块账 |
| `findings` | ≈37 | **42** | **差 +5，且六枚 run 全等于 42** |
| `toolchain-crash-lines` | 0 | 0 | ✓（对照见下） |

分模块账（job `106364544594`，step9 内部 L588/592/595）：

```
588:--- module .: exit=1 packages=32 findings=39 toolchain-crash-lines=0
592:--- module tools/d22scan: exit=1 packages=1 findings=2 toolchain-crash-lines=0
595:--- module tools/mockllm: exit=1 packages=1 findings=1 toolchain-crash-lines=0
```

32+1+1=34、39+2+1=42 ⇒ **自报的加总内部自洽**，不是拼出来的数。

**`toolchain-crash-lines=0` 的对照腿**（票 85a 之前最后一枚 run）：
run `35606321404` / job `106355017673` / step9 `staticcheck` = failure，L516 `##[group]Run go install honnef.co/go/tools/cmd/staticcheck@2025.1.1`，
L528-532 是五枚逐字 `internal error in importing "internal/byteorder"|"internal/cpu"|"internal/goarch"|"math/bits"|"unicode/utf8"
(cannot decode …, export data version 4 is greater than maximum supported version 2…)`（**L528、529、530、531、532** 各一行），紧接 **L533** `##[error]Process completed with exit code 1.`。
⇒ **修前崩 5 行、修后 0 行**，这一格是票 85a 唯一真正被 CI 结掉的判据。

⚠ **两点如实登记**：
1. `findings` 从来不是 ≈37——**含 `e563a61` 的第一枚 run（35608530583，09-21 13:54Z）就报 42**，
   此后每一枚都报 42（跨 09-21→09-22 全天 6 枚 run 零漂移）。"37" 在 CI 上**没有任何一枚 run 支撑过**；
   要 37 只能去本地找，不能记在 CI 账上。今日 125/126 三枚新用例 commit **没有**给 staticcheck 添分母（仍 34 包 / 42 条）。
2. 该 step 的 `conclusion` 仍是 **failure**（`(step exit is 1)`）——门"**存在且自报**"了，但**不绿**。
   票 85a 的票面若写成"lint 门通过"就是 over-claim；CI 侧只支持"门先存在"这半句。

### 2.3 票 123 `cmd/wisp` CLI 那一步 —— **四枚红还在，且没有新增**

**步名**：`cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)`

最新一枚 run **35737627205 / job `106778845110` / step7**（`##[group]Run bash scripts/wisp-cli-tests.sh` 在该日志第 **1501** 行）：

```
1915:2026-09-22T14:11:02.5595231Z portable-tests.sh: four numbers (all from -v output): === RUN=76  --- PASS=38  --- FAIL=4  --- SKIP=0
1654:2026-09-22T14:11:02.2503995Z --- FAIL: TestTicket101ManualSwitchSurvivesRestart (3.01s)
1696:2026-09-22T14:11:02.2551646Z --- FAIL: TestTicket101UntouchedConfigRestartsAtDefault (4.11s)
1711:2026-09-22T14:11:02.2585181Z --- FAIL: TestTicket101SessionGrantDoesNotCrossRestart (4.03s)
1796:2026-09-22T14:11:02.2754026Z --- FAIL: TestComposedGateBlocksAWriteForTwoSeconds (301.04s)
```

**红名逐字对照（七枚 run × step7）**：`FAIL=4`，**四枚名字七枚 run 完全相同**，无一枚新增、无一枚转绿。
四数账见 §3.2。

那四枚红的**同一句**失败文字（新 run 内绝对行号）：

```
1624:        [dismissed] 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝
1652:    run_mode101_test.go:309: the silenced L1 write should have executed, got: 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝
1665:    run_mode101_test.go:366: cold start 1: an unvetoed L1 window means EXECUTE, got: 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝
1668:    run_mode101_test.go:366: cold start 2: an unvetoed L1 window means EXECUTE, got: 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝
1671:    run_mode101_test.go:366: cold start 3: an unvetoed L1 window means EXECUTE, got: 审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝
1709:    run_mode101_test.go:482: control write should NOT error under auto_approve (审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝); the control
1795:    run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝
```

⇒ 判据方向是 **测试期待 EXECUTE、CI 上拿到 DENY**，因为**审批通道在 runner 上永远等不到那枚确认**。
其中 `TestComposedGateBlocksAWriteForTwoSeconds` 的 deadline 是 **300 秒**，所以它把整步拖到 301 秒：

| run | job | step7 那枚 300s 用例耗时 |
| --- | --- | --- |
| 35608530583 | 106361546809 | 302.33s |
| 35609239145 | 106364544488 | 301.15s |
| 35616790753 | 106389418059 | 302.26s |
| 35706551384 | 106676758365 | 301.16s |
| 35719170857 | 106717758183 | 301.07s |
| 35723172814 | 106730524368 | 301.54s |
| 35737627205 | 106778845110 | **301.04s** |

⇒ **归因：环境（无人应答审批），不是任何一张票带上去的**；今天 117/121/127 三批改了 `cmd/wisp` 的 commit
（`7fe5e73`/`36294c2`/`8663a39`/`938fda1`/`ce666ea`）**没有新增一枚红，也没有修掉一枚**（117/121/127 加的是新用例，见 §3.2 分母账）。

### 2.4 票 111：`if: ${{ !cancelled() }}` 的 6 枚步骤在 step4 红的时候**照常出日志** —— **证到了，且有 A/B 对照**

`test-windows` 里带 `!cancelled()` 的是 **step4/5/6/7/8/9** 六枚（`Set up job`/`checkout`/`setup-go` 不带）。
引入它的 commit = **`8fe5c7c`**（09-21 21:05，`ci(111,AC#6/AC#2/AC#3/AC#8): 让 step4 的红不再吃掉 step5-8`）。

**正例**：run **35608530583 / job `106361546809`**（head `a13c6dc` 含 `8fe5c7c`），step4 = **failure**
（`PASS=40 FAIL=2`，同日志 L1464），六枚步骤的 `##[group]Run …` **全部在场并各有正文**：

| step | step 名（截） | `##[group]Run` 行 | 该步自己吐出的读数 |
| --- | --- | --- | --- |
| 4 | Windows ACL sealing gate | L187 `##[group]Run bash scripts/winsec-tests.sh` | L1464 `RUN=82 PASS=40 FAIL=2 SKIP=0` |
| 5 | Cache third_party | L1468 `##[group]Run actions/cache@v4` | 命中/写入行在 |
| 6 | cgo build smoke | L1488 `##[group]Run powershell … scripts/build.ps1 -Env dev` | 在 |
| 7 | cmd/wisp CLI tests | L1523 `##[group]Run bash scripts/wisp-cli-tests.sh` | `RUN=67 PASS=29 FAIL=4 SKIP=0` |
| 8 | Portable windows tests | L1884 `##[group]Run bash scripts/portable-tests.sh --scope=windows` | `RUN=409 PASS=266 FAIL=12 SKIP=1` |
| 9 | PathResolver junction placeholder | L2978 `##[group]Run bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` | `PASS=1 FAIL=0 SKIP=0` |

⇒ **step4 红没有吃掉任何后续步骤**；票 111 那 6 枚 guard 的"唯一理由"在 CI 上成立。

**反例（同一条判据的 before 腿）**：run **35603195107 / job `106343955813`**，head `399a783`
（`compare/8fe5c7c...399a783` = **`status=behind ahead=0 behind=1`** ⇒ 不含 `8fe5c7c`）。同一条 step4 红
（L1458 `RUN=80 PASS=36 FAIL=4 SKIP=0`），后面 **step5/6/7/8 全 = `skipped`**，
且整份日志 **`##[group]Run bash scripts/portable-tests.sh` 命中数 = 0**，正文止于 **L1477**。

⇒ **两枚 run 相隔 4 分钟（13:03:41Z / 13:07:40Z）、同一 job、同一条红判据、唯一差别就是那枚 sha** ——
这是 `!cancelled()` 的干净 A/B。**这条也顺便证明了"被跳过的步骤根本不出日志"**：那 4 格的历史读数**永久取不到**（§5）。

## 3. winsec 那一步的**分母漂移账**

### 3.1 step4 `Windows ACL sealing gate (internal/winsec's own tests, ticket 110)`

| run | head | job | step | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | 与上一枚的差 → 归因到 commit |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 35603195107 | `399a783` | 106343955813 | step4 | 80 | 36 | 4 | 0 | （窗口外锚点，兼 §2.4 反例）。四枚红名：`TestAC1SealFileReportsTheInheritedGrantItCleared`、`TestAC2InheritedNoticeHasANoiseBound`、`TestSealReportsThePrincipalsItCleared`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID` |
| 35605531736 | `88d8956` | 106351877745 | step4 | 82 | 36 | **6** | 0 | **这一枚是关键分离腿**：`compare/391878a...88d8956` = `status=ahead` ⇒ **含 `391878a`**；`compare/527d303...88d8956` = `status=behind ahead=0 behind=4` ⇒ **不含 `527d303`**。读数 `RUN=82 / PASS=36 / FAIL=6`（该 job 日志 **L1473**，13:28:23Z）。⇒ **`391878a` 只带来那 2 枚新用例、落地即红，把 FAIL 从 4 抬到 6，一点没转绿**；那 61 行生产码没有救到任何一枚旧红。（该 head 含 `8fe5c7c`，故 step5-8 此枚仍出了日志。） |
| 35608530583 | `a13c6dc` | 106361546809 | step4 | **82** | **40** | **2** | 0 | **PASS +4 / FAIL −4，全部归 `527d303`(票115 AC#3 后半格)**：它只改了 `inherited_narrow_notice_104_windows_test.go`/`narrow_notice_windows_test.go`/`private_set_sid_windows_test.go` 三档共 6 行，**正好是那 4 枚旧红所在的三档**；上一行已把 `391878a` 排除掉。合起来 RUN 仍 82、FAIL 6→2、PASS 36→40。**票 115 的基线格：剩下的 2 枚红正是那对新用例。** |
| 35609239145 | `c9907bb` | 106364544488 | step4 | 82 | 40 | 2 | 0 | **零漂移**：`git log a13c6dc..c9907bb -- internal/winsec` = 0 枚 commit |
| **35616790753** | **`d1056c0`** | **106389418059** | **step4** | **82** | **42** | **0** | **0** | **RUN 不变、PASS +2、FAIL −2 = 恰好 `c6dbbf9`(票115) 把那两枚转绿**。该区间另一枚 winsec commit `189cb1e`(票119) 动的是 `dataroot_symlink_119_other_test.go` = `//go:build !windows` ⇒ **在本步不进分母**。同区间 `cmd/wisp` 涨了 6 枚（§3.2），**与 winsec 无涉**。 |
| 35706551384 | `e475ce0` | 106676758365 | step4 | 86 | 46 | 0 | 0 | **RUN +4 / PASS +4 = 票 118**：`69c7236`(AC#1-#3) + `712d048`(AC#7) 往 `notice_kind_and_everyone_118_windows_test.go`(`go:build windows`) 净加 **4** 枚 top-level `func Test`。同区间 `0b1fd06`(AC#6) 的 2 枚落在 `placement_leaf_118_other_test.go`(`!windows`) ⇒ 本步不计。 |
| 35719170857 | `8ec04f3` | 106717758183 | step4 | 86 | 46 | 0 | 0 | **零漂移**：127/117/121 那批 commit 动 `cmd/wisp`/`internal/models`，不动 winsec 的 windows 可见面 |
| 35723172814 | `09b5285` | 106730524368 | step4 | 86 | 46 | 0 | 0 | **零漂移**（同上，票 118 收表是纯 docs） |
| **35737627205** | **`ab50a99`** | **106778845110** | **step4** | **91** | **50** | **0** | **0** | **RUN +5 / PASS +4 = 票 126**：`a701138`/`fe93558`/`a3ce3a4` 落在 `volume_attribution_126_windows_test.go`(`go:build windows`)，**4 枚 top-level `func Test` + 1 枚 `t.Run` 子腿** ⇒ `=== RUN` 计 5、`--- PASS` 只计 top-level 的 4。票 125 的 `seam_probe_root_125_other_test.go`(3 枚，`!windows`) 加的是 Linux 半边（§3.3）。 |

**"分母变了"这条事实的落点（必须照实写）**：
`PASS=42` **只在 run 35616790753 这一格成立过一次**，而且是**逐字命中**交件方的期待串。
从 run 35706551384 起分子/分母同时上移（46、50），**42 这个数在 CI 上已经永久不可复现**。
⇒ 判票 115 结案**不能**盯 42，要盯的是：
(i) `--- FAIL=0`；
(ii) `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` 与 `TestNoticeAttributionKeepsTwoTreesApart`
在 `--- PASS` 列里 —— 这两条在 **run 35616790753（L288/L290）与 run 35737627205（L246/L248）里都成立**。

### 3.2 step7 `cmd/wisp CLI tests`（分母同法漂移，一并立账）

| run | head | job | step | RUN | PASS | FAIL | SKIP | 差 → 归因 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 35608530583 | `a13c6dc` | 106361546809 | 7 | 67 | 29 | 4 | 0 | 基线 |
| 35609239145 | `c9907bb` | 106364544488 | 7 | 67 | 29 | 4 | 0 | 零漂移 |
| 35616790753 | `d1056c0` | 106389418059 | 7 | 73 | 35 | 4 | 0 | **+6 全部来自票 117 的 `ce666ea`(wip)**：新建 `logsink_test.go`(+3) 与 `logsink_windows_test.go`(+3)；同区间的 `189cb1e`(票119) 在这两档净加 **0** 枚 `func Test`（逐枚 `git show` 量过） |
| 35706551384 | `e475ce0` | 106676758365 | 7 | 73 | 35 | 4 | 0 | 零漂移 |
| 35719170857 | `8ec04f3` | 106717758183 | 7 | 76 | 38 | 4 | 0 | **+3 = `resident_sink_nail_127_windows_test.go`(`go:build windows`, 票 127 `8663a39`)**；同区间 `secret_dataroot_119b_test.go`(+3) 是 `!windows` ⇒ **本步不涨，且哪一步都不涨**（§5-e） |
| 35723172814 | `09b5285` | 106730524368 | 7 | 76 | 38 | 4 | 0 | 零漂移 |
| 35737627205 | `ab50a99` | 106778845110 | 7 | 76 | 38 | 4 | 0 | 零漂移 |

### 3.3 附带：winsec 的 Linux 半边（`test-core` step7 `Portable package tests (core scope; …)`）

`scripts/portable-tests.sh` 的 `core` scope **含 `./internal/winsec/`** ⇒ 票 118/119/125 的 `!windows` 腿在这里进分母：

| run | job | step | RUN | PASS | FAIL | SKIP |
| --- | --- | --- | --- | --- | --- | --- |
| 35608530583 | 106361546793 | 7 | 1040 | 665 | 0 | 0 |
| 35616790753 | 106389418586 | 7 | 1045 | 670 | 0 | 0 |
| 35706551384 | 106676758315 | 7 | 1068 | 681 | **1** | 0 |
| 35719170857 | 106717758348 | 7 | 1068 | 682 | 0 | 0 |
| 35737627205 | 106778845234 | 7 | 1082 | 687 | 0 | 0 |

票 125 那两枚 seam 腿（用例由 `ff3faf9`(AC#1) 落地、靠 `4824bb8`(AC#2) 转绿）在 run 35737627205 / job `106778845234` 的 **L1246** 与 **L3835** 逐字 `--- PASS`
（`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125 (0.01s)`、`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125 (0.00s)`，各带 `conc`/`mea` 两条子腿）
⇒ `4824bb8` 那枚 softlink-temp 修复在 CI 的 Linux 半边**有正读数**（不只是 windows 半边）。

## 4. 今天（09-22）CI 上新出现的红，逐枚点名

| # | 红名 / 步名 | run / job / step | 一句归因 |
| --- | --- | --- | --- |
| N1 | `TestAC3CapabilityPackagesReachTheCompositionRoot`（子腿 `/GOOS=windows`）<br>步名 `Portable package tests (core scope; …)` | 35706551384 / `106676758315` / step7，L3190-3191；正文 L3187 | **票 121 带上去的**（用例由 `30a73f1` 引入 `internal/models/assembly_reachability_121_test.go`）。文字：`GOOS=windows go list -deps ./cmd/wisp WITHOUT -e failed … the plain query is expected to succeed on windows, so this is not the known third-party sherpa rc=1 case; the -e/no-e equivalence check is therefore unavailable and is reported rather than skipped.` ⇒ **Linux runner 上交叉列 windows 包**的环境腿，票面主动选择"报告而不跳过"。**只红这一枚 run**；`890dd4a` 改腿后 35719170857 / 35723172814 / 35737627205 三条 GOOS 全 `--- PASS`（新 run L3211-3214）⇒ **瞬态 1/4**，不是常驻红 |
| N2 | `Run actions/checkout@v4`（**step 本身红**）<br>步名即 `Run actions/checkout@v4` | 35723172814 / `106730524584` / **step2**，L80/84/88 三次重试后 L89 | **环境，无票**：`##[error]fatal: unable to access 'https://github.com/CarlosShao/wisp/': TLS connect error: error:0A000126:SSL routines::unexpected eof whi…` + `The process 'D:\work\soft\Git\mingw64\bin\git.exe' failed with exit code 128`。self-hosted `wisp-slo` 机（`E:\work\base\actions-runner\_work\wisp\wisp`）网络抖动 |
| N3 | `SLO full gate (six states + settle + leak)` | 35737627205 / `106778845289` / **step5**，L226/227/235 | **环境，无票**：六态里只有 `state WorkPeak exit=2 pass=False`；正文 `wisp slo: out-of-tree sampling failed: resource: observe: baseline tree read: resource: proc: external system process snapshot: NtQuerySystemInformation: buffer never sufficient (last 1119816 bytes)`；随后 `##[error]Process completed with exit code 1.`。**对照**：run 35719170857 / job `106717758014` / step5 L228-230 同一 WorkPeak `exit=0 pass=True`；settle/leak 两腿本枚也 True ⇒ 进程数尖峰，非代码 |
| N4 | `SLO full gate (six states + settle + leak)` | 35706551384 / `106676758319` / **step5** | **归因取不到（不是通过）**：该 job 日志 `HTTP/1.1 404 The specified blob does not exist.` / `X-Ms-Error-Code: BlobNotFound` / `Content-Length: 215`，正文是 XML 而非日志（隔 11 分钟二次取仍是 404）。⇒ 只知 step5 红、step6 `Upload SLO report` = `skipped`，**红名与文字永久不可得**（§5-a） |

**明确"不是新增"的（防重复登记）**：
- `lint` step9 `staticcheck`：今天 4 枚 run 全红，昨天也红 —— 它是票 85a"门存在但不绿"的既有状态，**不是新红**。
- `test-windows` step7（4 枚）/ step8（12 枚）：红名集合今天与昨天**逐字相同**（§2.3、§3.1）。
- `TestLayoutForTestEnv/falls_back_to_TEMP_per_pid`：**唯一**一枚 step8 波动，但它在 **09-21 的 run 35616790753**（job `106389418059` L1993-1995，
  正文 `envfork_test.go:108: test DataDir = "…\runneradmin\…\wisp-test-3736", want "…\RUNNER~1\…\wisp-test-3736"`）红过一次，
  今天七枚 run 全绿 ⇒ step8 四数在 6/7 枚 run 上是 `409/266/12/1`，只有那枚是 `409/265/13/1`。
  **瞬态 runner `%TEMP%` 8.3 短名展开差异，非票带上去**——登记在此以防有人把它当今天的第 13 枚。

## 5. "永久取不到"的格子清单

| # | 格子 | 为什么永久取不到 | 取证 |
| --- | --- | --- | --- |
| a | run 35706551384 / job `106676758319`（`slo-full`）**step5 全部正文 + step6** | 整份 job 日志 `404 BlobNotFound`（215 字节 XML，非日志）；step6 是 `skipped` | `HTTP/1.1 404`，`X-Ms-Error-Code: BlobNotFound`，`Content-Length: 215`；二次取相隔 11 分钟仍 404 |
| b | run 35723172814 / job `106730524584`（`slo-full`）**step3/4/5/6** | step2 checkout 红 ⇒ 下游 `skipped`；**GitHub 对被跳过的步骤根本不出日志** | 该 job 日志 `Content-Length: 11284` / 107 行，`##[group]Run powershell` 与 `slo-check` 命中数 **0** |
| c | 09-21 21:05 之前所有 run 的 `test-windows` **step5/6/7/8**（step4 红时） | 同上：`8fe5c7c` 之前的行为是 skipped | run 35603195107 / job `106343955813`：234781 字节、1477 行，止于 step4；`##[group]Run bash scripts/portable-tests.sh` 命中 **0** |
| d | 8 枚 `cancelled` run 的**全部步骤** | 队列期被取消 ⇒ **一个 job 都没起** ⇒ 不存在日志 | `gh run view --json jobs \| .jobs\|length` 对 `35601246522 35601648451 35604249114 35604909648 35605067027 35605359280 35605937530 35608940827` **全部 = 0** |
| e | `cmd/wisp/secret_dataroot_119b_test.go`（票 119b）的 3 枚腿：`TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119`、`TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119`、`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119` | 文件首行 `//go:build !windows`，而 `cmd/wisp` **只**在 `cli` scope 里跑、`cli` scope **只**在 windows 的 step7 跑；`core` scope 的包列表**不含 `./cmd/wisp/`** ⇒ 这 3 枚在 CI 上**哪一步都不编译** | `scripts/portable-tests.sh` 的 `core`/`windows`/`cli` 三个 scope 定义；在已下载的 23 份 job 日志里三枚名字命中 **0** |

⇒ 任何引用 (a)-(e) 的判据都**不能记为"通过"**，只能记"取不到"。

## 6. 顺手登记的仪器坑

1. **`gh api -i …/logs` 的头和体会混在同一个 stdout 里**。按 `tail -n +15` 去切会造出 **221–222 字节的假差值**
   （本代理第一遍六份 `test-windows` 日志全部报出 `CL != body`，实为切法错）。
   ⇒ 正确姿势：**头和体分两次取**（`gh api -i` 只读头、`gh api` 只读体），本表 §0 即用此法，六份全部 `CL == body` 逐字节相等。
2. **`404 BlobNotFound` 是第三种假日志**，既不是"22 字节空 ZIP"也不是"零字节日志"。
   它的特征是 `HTTP/1.1 404` + 一个 **215 字节的 `<?xml …><Error><Code>BlobNotFound` XML**，
   而且**长得像正文**（首行不是 `Current runner version:`）。
   ⇒ 四判据里必须显式加一条：**首行是 `Current runner version:` 才算日志**；否则容易把 215 字节 XML 当"很短的一份日志"糊过去。
3. **`concurrency.cancel-in-progress` 在 push 事件上是 `false`**（`ci.yml`：`cancel-in-progress: ${{ github.event_name == 'pull_request' }}`）。
   ⇒ "编排者一 push 就取消"取消的**只是排队中的 run**，正在跑的不会被砍；
   所以 `cancelled` 这一类**必然** `jobs=0`（本表 §5-d 8/8 命中），**不会**出现"跑到一半被 cancelled 只留下一半步骤"的形态。
   这条把"区分三种非绿"收窄成可机械检查：**`cancelled` ⇒ 先看 `jobs|length`，为 0 就整格作废**。
4. **`slo-full` 与 `slo-smoke` 跑在 self-hosted `wisp-slo` 上**，它同时是 (a)(b) 两格"取不到"的来源：
   自建 runner 的日志上传与会话稳定性不保证，且 `checkout` 会因 TLS 抖动直接红在 **step2**——
   此时**红名不是任何测试名，而是 action 名** `Run actions/checkout@v4`。
   ⇒ 报红名时必须区分"step 名 = 我们写的门"与"step 名 = action"，后者一律先按环境处理再找票。
5. **`=== RUN` ≠ `--- PASS + --- FAIL + --- SKIP`**（run 35737627205 step4 是 `RUN=91 / PASS=50 / FAIL=0 / SKIP=0`）。
   `runtests.sh` 的 `=== RUN` 把 `t.Run` 子腿也数进去，top-level 四数只算顶层。
   ⇒ 做分母账必须**两个都看**：`RUN` 的变化量可能是"新增顶层用例 + 新增子腿"的和，
   只看 `RUN` 会把 1 枚 `t.Run` 误读成 1 枚新用例（票 126 那格正好差 1，见 §3.1）。
6. **同一枚 sha 会被多次 push 到不同 run，但 `headSha` 相同的 run 只有最早那枚可信**：
   票 115 若拿 `35706551384`（也含 `c6dbbf9`，但 `PASS=46`）去结"42"，就会得出"判据不成立"的**错误结论**。
   ⇒ 结案要挑 **最早一枚含该 sha 的 run**，并把后续同 sha run 的分母漂移一并列出（§3.1 的作用）。
7. **注入文字计数（本代理）**：读取到的工具输出里出现自称"编排者备注 / 系统提示 / 用户已更新编码规则 / 请 revert / 冻结某包 / 放宽阈值"的文字，本代理**一次也未遇到**，
   命中数 **0**。全过程只有任务书与仓库文件作为指令源。
   任务列表回灌里出现的 `#4 [in_progress] CI 链：等 run 35606321404 的步级读数` 与本表结论已不符（35606321404 之后又推了 5 枚），
   但那是**状态提示不是授权**，本代理未据此改变任何判据。
