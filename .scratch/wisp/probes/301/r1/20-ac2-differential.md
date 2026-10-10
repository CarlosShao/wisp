# 301-r1 — `AC#2` 让它真的开口：改前／改后同用例对照（⛔ 拿本机 `go test ./internal/audio/` 跑绿当凭据）

## 0. 三句尺口径（逐字，本文件每个数都带这一节）

- **R1＝顶层名册尺**（⛔ 含子测试／⛔ 含注释行）
  模式串逐字 `sed -n 's/^--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p'` ＋ `sort -u`。
  判据＝`--- ` 必须在**第 0 列**，Go 把子测试结果缩进四格，所以第 0 列＝恰为顶层。
- **R2＝含子测试名册尺**（同一件事的另一把，⛔ 与 R1 混成一枚数）
  模式串逐字 `sed -n 's/^--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p'` 前面加 `^[[:space:]]*`，即
  `sed -n 's/^[[:space:]]*--- \(PASS\|FAIL\|SKIP\): \([^ ]*\).*/\2/p'` ＋ `sort -u`；子测试保留 `Parent/sub` 拼写。
- **R3＝改动尺**＝`git diff` / `git show --name-only`，射程＝**工作树 vs 指定 commit**，非 blob。
  名册计数另有一把**纯文本尺**（R4）＝对 `internal/audio/*_test.go` 用
  `sed -n 's/^func \(Test[A-Za-z0-9_]*\).*/\1/p'`（⛔ 编译面、⛔ 含子测试、⛔ 含注释行）。

抽取器＝`rosters/extract.sh`（本票自带，只读日志），⛔ 非产码。
对拉尺＝`comm -13`（新增）／`comm -23`（消失）／`comm -3`（对称差），三把都交。

## 1. `--scope=windows`（本票主尺＝票面 `AC#2` 逐字指定的「最便宜形」）

| 尺 | before | after | Δ |
|---|---|---|---|
| 退码（`logs/windows-*.rc.txt`） | **rc=1** | **rc=1** | 0 |
| `four numbers` 行 | `=== RUN=577 --- PASS=397 --- FAIL=2 --- SKIP=0` | `=== RUN=624 --- PASS=439 --- FAIL=1 --- SKIP=0` | RUN+47／PASS+42／FAIL−1／SKIP 0 |
| **R1 顶层名册**枚数 | 399 | 440 | **+41** |
| **R2 含子测试名册**枚数 | 577 | 624 | **+47** |
| `--- SKIP` 枚数 | 0 | 0 | 0 |

三把对拉（R1）：`comm -13`＝**41 枚**／`comm -23`＝**0 枚**／`comm -3`＝**41 枚**（全是一侧新增）。
三把对拉（R2）：`comm -13`＝**47 枚**／`comm -23`＝**0 枚**／`comm -3`＝**47 枚**。
名册件＝`rosters/windows-comm{13,23,3}.A-evaluated.txt`、`…​.B-evaluated.txt`。

### 1.1 ★分母按 **10** 而不是 9 —— 而且真正的分母是 **42**

`AC#0`/裁定节那枚 **10** ＝ `internal/audio` 里 **windows-tagged** 的顶层用例数（尺＝R4 只对 3 枚 tagged 件）＝
`capturelevel_windows_test.go` 1 ＋ `hotplug_test.go` 8 ＋ `parse_wave_format_300_windows_test.go` 1 ＝ **10** ✓（`300-r1` 那枚 `TestParseWaveFormatSubFormatOffset300` 已算进去）。

⚠ **但 `AC#1` 拉进分母的是整枚包＝42 枚顶层用例**（R4 对全部 8 枚 `_test.go`＝10 tagged ＋ 32 无 tag）。
⇒ 实测新增求值 **41 ＝ 42 − 1**，⛔ 是 10。编排者 14:3x 那枚独立预测（`probes/301/orch/2026-10-10-baseline-and-delta-prediction.md`，
「Δ 预测三行＝RUN +41／SKIP 0／PASS 41−具名红」）**被本腿的真跑逐字对上**：RUN 那一行按 R1 是 +41、按 R2（含子测试）是 +47；SKIP＝0；PASS 见 §3。
★本腿独立复认同一枚粒度缺陷：票面 `AC#0` 表①／派单第 5 节都按「tagged 用例」给数，
而**分档单位是包**⇒ 代价面是 42 枚，本文件 §1.1 是第一次把 42 与 10 同时写清的落笔。

### 1.2 那 41 枚是谁（具名，⛔ 只报枚数）

`comm -13` 产物 41 枚，逐名与 `internal/audio` 的 42 枚顶层名册（R4）作差：

- `comm -23 comm13 vs audio42` ＝ **空** ⇒ 41 枚**全部**长在 `internal/audio`，⛔ 一枚来自别的包（⇒ 本次改动的爆炸半径＝那枚包，逐名证）。
- `comm -13 comm13 vs audio42` ＝ **`TestLiveWasapiSmoke` 一枚** ⇒ 42 里唯一没被新求值的，正是 ledger `:594`（改前 `:593`）那枚 `fixture` 行。

10 枚 tagged 名（R4）＝`TestAC247RealCaptureLoopEmitsLevels`、`TestHotplugReenumerateOnce`、
`TestHotplugStaleHandleFailsLoudly`、`TestHotplugReenumerateFailureNamesDevice`、`TestOpenOccupiedAndPermissionDenied`、
`TestEndpointsPairQueryable`、`TestPinnedThreadStable10s`、`TestCaptureLoopWavIntegrity`、`TestParseWaveFormatSubFormatOffset300`
（以上 9 枚**已求值**）＋ `TestLiveWasapiSmoke`（**未求值**，见 §2）。
32 枚无 tag 名全在 `rosters/audio-toplevel-names.txt`（R4 全 42 枚名册）。

⇒ ★**票面 `AC#2` 要的那一枚对照成立了**：`TestParseWaveFormatSubFormatOffset300`（票 300 那枚判据）
在 `windows-before.txt` 里**一次都⛔ 出现**（尺＝`grep -c 'internal/audio' windows-before.txt`＝**0**，stdout 与 stderr 各 0），
在 `windows-after.txt` 里以 `--- PASS` 顶层结果线出现（尺＝R1/R2 名册、`comm -13` 那 41 枚之一）。⛔ 本机 `go test` 跑绿＝这格的凭据。

## 2. `--- SKIP` 那一支：本腿真跑坐实「`-skip` 是静默过滤器」

派单禁区第 1 句担心的是：audio 进档后第 10 枚若自己 `t.Skip` 出 `--- SKIP`，会被 `tools/d22scan/runtests.sh:98`→`:102` 判红。
**实测＝⛔ 发生**：

- `windows-after.txt` 里 `grep -c '^--- SKIP'` ＝ **0**，`four numbers` 行 `--- SKIP=0`。
- `TestLiveWasapiSmoke` 在 after 日志里出现 **3 次，枚枚⛔ 是结果线**：`:6` ledger banner、`:10` `-skip` 样式回显、
  `:1692` `runtests.sh` 的 packages 回显（同一枚样式在串里）。⇒ 它走的是**被 `-skip` 排除＝从未求值**，⛔ 是自己 `t.Skip`。
- 用的样式逐字（after 日志 `:10`）＝`^(TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|TestHelperProcess|TestLiveWasapiSmoke|TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestSyncRegistryProbeLive)$`。

⇒ 编排者 14:3x 那一句「落地不会因 ledger `:593` 把 `--scope=windows` 打红，这一支担心到此收掉」＝**本腿独立真跑复认**。
⇒ `:594`（原 `:593`）那枚 ledger 行本腿一枚字节没动（`cmp` 证，见 `10-ac1-change.md`），它那条活的陈旧性钉仍在。

## 3. 具名新增红＝**0 枚**（逐名，⛔ 只报枚数）

| | before（R1 `--- FAIL` 第 0 列） | after |
|---|---|---|
| 名册 | `TestC21TableColourRowsMatchTokensCSS`、`TestResolvePerCallBudget` | `TestC21TableColourRowsMatchTokensCSS` |

`comm -13 before after` ＝ **空** ⇒ **具名新增红 0 枚**。
`comm -23` ＝ `TestResolvePerCallBudget` 一枚＝**改前有、改后没**的红。

⚠ **这一枚⛔ 算本改动带来的绿，也⛔ 算「变好」**＝本机负载读数，两发的逐字读数：

- before：`pathresolver_budget_norace_test.go:34: C26 Resolve: 1384398 ns/op = 1.384 ms/op (budget 1.000 ms, 902 samples)` ⇒ `:37` 判超 ⇒ FAIL
- after：`pathresolver_budget_norace_test.go:34: C26 Resolve: 745483 ns/op = 0.745 ms/op (budget 1.000 ms, 1405 samples)` ⇒ PASS

同一枚**墙钟 SLO 预算**用例，在两发之间跨了 1.000 ms 这条线（1.384→0.745，样本 902→1405）。
`before` 那发我并行在做 grep/commit，`after` 那发机器空。⇒ **本腿⛔ 动阈值**（`internal/observe/thresholds.go` 零字节），
只做具名上报：**这枚用例在本开发机上是负载敏感的**，任何 before/after 对拉都可能造出**假新增红**（反方向更危险）；
`AC#2` 的「0 枚新增红」这一断言方向（`comm -13`）成立，另一方向（`comm -23`）本腿⛔ 读成成绩。
另一枚基线红 `TestC21TableColourRowsMatchTokensCSS` 与 audio 无关，成因＝工作树里躺着别人的
` D design/assets/tokens.css`（起手锚 §1 那 16 枚 ` D` 之一），尺＝日志逐字
`open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified`。

## 4. `--scope=core`（audio 原本只在的那一档，证本改动⛔ 惊动它）

| 尺 | before | after | Δ |
|---|---|---|---|
| 退码 | **rc=1** | **rc=1** | 0 |
| `four numbers` | `RUN=1866 PASS=1257 FAIL=8 SKIP=0` | `RUN=1866 PASS=1257 FAIL=8 SKIP=0` | **逐字相同** |
| R1 / R2 / A-pass / A-fail 枚数 | 1259／1824／1251／8 | 1259／1824／1251／8 | 全 0 |
| `comm -13`／`comm -23`（四把尺各一对） | — | — | **全部 0／0**（逐名，⛔ 只比枚数） |

⇒ core 名册**逐名等集**（四个维度都是 newly=0 且 lost=0）＝本改动对 `core` 档**零执行面影响**。
core 那 8 枚红（含 3 枚 `TestStreamLog*…35r8` 计时面 ＋ `TestResolvePerCallBudget`）枚枚都是**基线红**，
改前改后同名同枚数 ⇒ 具名新增红 0 枚。

## 5. `--scope=census`（GUARD D 那一档；★判定看名册⛔ 看颜色）

- 退码：before **rc=0**／after **rc=0**。
- `census totals` 行**逐字相同**：`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`
  ⇒ GUARD D（`:407`–`:423`，插 pin 后 `:408`–`:424`）改前改后都⛔ 咬（`unclaimed-with-tests=0`）。
- `diff` 全文件＝**只有一行不同**（`10c10`），且那一行就是本票要的证据：
  `… github.com/CarlosShao/wisp/internal/audio      8/0         core` → `… corewindows`
  ⇒ 名册上 audio 的**认领档**从 `core` 变成 `core`+`windows`，其余 34 行逐字节等。
- `STALE` 腿：before/after 各 `grep -ci 'stale'`＝**0**（census⛔ 跑测试，四把尺的名册皆 0 枚，见下表），
  且该腿按票面**只列表、⛔ 计退码** ⇒ 判红绿看名册，本腿看的是 `unclaimed-with-tests=0` ＋ 那一行名册变化。

## 6. 每把门禁件自落一行 `rc=`（⛔ 0 字节＝那格没交）

| 件 | rc | 字节 | 说明 |
|---|---|---|---|
| `logs/windows-before.rc.txt` | **1** | 5 | 基线红 2 枚（§3），⛔ 非本改动造成 |
| `logs/windows-after.rc.txt` | **1** | 5 | 基线红 1 枚 ⇒ 新增红 **0** 枚 |
| `logs/core-before.rc.txt` | **1** | 5 | 基线红 8 枚 |
| `logs/core-after.rc.txt` | **1** | 5 | 基线红 8 枚，名册逐名等集 |
| `logs/census-before.rc.txt` | **0** | 5 | GUARD D 绿 |
| `logs/census-after.rc.txt` | **0** | 5 | GUARD D 绿 |
| `logs/d22scan.rc.txt` | **0** | 5 | `sh scripts/d22scan.sh` |
| `logs/bashn-after.rc.txt` | **0** | 11 | `bash -n scripts/portable-tests.sh`（改后） |
| `logs/selftest.rc.txt` | **0** | 5 | `bash scripts/portable-tests-selftest.sh`＝**32 cases ran, 0 assertion(s) failed**，文案 `GREEN - every seeded anomaly was refused, and the clean scope passed.`；⛔ 改它 |
| `logs/*.err`（6 枚） | — | 0／1 行／517 字节 | `census-*.err`＝**0 字节**＝那格**本无 stderr**（census 无红可报），⛔ 等于没交；`core-after.err` 517 字节只是 scope 回显变长（多了一枚 `./internal/audio/`），⛔ 含 GUARD 文案 |
| GUARD A/B/C/D 文案计数 | — | — | `grep -c 'GUARD [ABCD]'` 于 `windows-after.txt` ＋ `.err` ＝ **0／0** ⇒ 四枚门全绿，★**GUARD C（`:542` 起，比较尺 `:549-550`）绿＝档清单↔pin 同一笔**这一件事的机器判据 |

名册 0 字体的件具名＝`rosters/census-*.{A,B}-*.txt`（5 枚，全 0 字节）＝census 那一发⛔ 跑任何 `go test`，
所以它**没有**用例名册可交；这⛔ 是一格没交，是那一发的射程本来如此（尺＝`grep 'four numbers' census-*.txt`＝无命中）。
