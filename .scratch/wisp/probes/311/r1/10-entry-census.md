# 10-entry-census.md — 票 311 落地腿 `311-r1` 格一：入口普查（现量枚数）＋守卫落点

锚（起手自己现跑）＝`git rev-parse --short HEAD`＝**`23c19afa`**（14:4x 起手那把）；
本件写于 **`d5ef9cb9`**（写到本节时 HEAD 又走到 `35744137`）——起手之后 HEAD 走了三笔
（`ef8d16b3`＝编排者收 310-r2＋A862、`d5ef9cb9`＋`35744137`＝票 295 的 295-r1），三笔都**没碰我的写面**：
尺＝`git diff --name-only 23c19afa..HEAD`＝**36 行**（`diff_file_count=36`），同一把尺过滤到
`-- scripts tools/d22scan`＝**0 枚**，`git status --porcelain -- scripts tools/d22scan`＝`write_surface_lines=0`（件 `logs/02-head-move.txt`）。
写面复核（起手三件第 1 条）＝`git status --porcelain -- scripts tools/d22scan`＝**空**（14:4x 与 14:5x 各跑一次，都空）。
⇒ ⛔ 停手条件没成立，开工。

## 1. 今天一共有几枚入口会设 `GOTMPDIR`？＝**4 枚**（现量，⛔ 抄派单那个数）

尺＝`git grep -n "export GOTMPDIR" HEAD -- scripts tools`（射程＝`scripts`＋`tools`，枚数＝命中行去重到文件）
⇒ `export_gotmpdir_file_count=4`（件 `logs/01-entry-census.txt` ruler A/B）。四枚逐枚：

| # | 文件 | 设 root 那行（内容锚，逐字） | 该文件里真正起测试的那一行 |
|---|---|---|---|
| 1 | `scripts/portable-tests.sh` | `:147  export GOTMPDIR="$gotmp_run"`（根名 `portable-…`，base 走 `:144 gotmp_base=${WISP_TEST_GOTMP_ROOT:-…}`） | `:764  sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern"`（`go test -list` 在 `:686`，属 ledger 守卫链、不起测试） |
| 2 | `scripts/winsec-tests.sh` | `:65  export GOTMPDIR="$gotmp_run"`（根名 `winsec-…`，`:63` 一行内联 base） | `:128  bash "$portable" "${scope[@]}" 2>&1 \| tee "$capture"` |
| 3 | `scripts/wisp-cli-tests.sh` | `:78  export GOTMPDIR="$gotmp_run"`（根名 `cli-…`，`:76` 内联） | `:146  bash "$portable" --scope=cli` |
| 4 | `tools/d22scan/runtests.sh` | `:94  export GOTMPDIR="$gotmp_run"`（根名 `d22scan-…`，`:92` 内联） | `:116  go test -v -count=1 "$@" >"$out" 2>&1` |

派单的候选名单＝`winsec-tests.sh`／`wisp-cli-tests.sh`／`portable-tests.sh`／`runtests.sh`——现跑**四枚全中、一枚不多**。

⚠ **`78782c20` 的 commit 标题写着 "5th go-test entry"，与本表的 4 枚⛔ 矛盾吗＝不矛盾，具名在此**：
那枚"第 5 枚"是 `tools/d22scan/runtests.sh` 自己注释里的说法，尺＝现读 `runtests.sh:69-74` 逐字
`# This is the fifth entry that really starts `go test`: scripts/d22scan.sh:51 …`＋`grep -n "fifth" tools/d22scan/runtests.sh`＝1 命中。
它数的是"**起 `go test` 的入口第五枚**"（`scripts/d22scan.sh:51 sh "$root/tools/d22scan/runtests.sh" -C tools/d22scan ./...`
＝那条链上的第五枚调用点），⛔ 是"设 `GOTMPDIR` 的第五枚"。同一枚注释 `:72-74` 自己写死了口径：
`# The block is copied verbatim from scripts/winsec-tests.sh (only the entry name and the echo prefix change)
# so all four entries that start tests present ONE shape, not four dialects.` ⇒ **四枚、同形**，与本票口径一致。
分母/正控：`scripts/d22scan.sh` 全文 `GOTMPDIR` 命中 **0**（尺＝`grep -c GOTMPDIR scripts/d22scan.sh`；
它靠 `runtests.sh` 那一枚，正是注释里"the delegated run inherits it"那条契约），
`scripts/portable-tests-selftest.sh` 的 `GOTMPDIR` 命中亦 **0**（件 `logs/01-entry-census.txt` ruler D：`selftest_gotmpdir_hits=0`、
`selftest_gotest_noncomment=0`，它只 `cp` 那四枚进 shadow 树、自己⛔ 设根）。

## 2. 新守卫落在每枚入口的哪两步之间（＝派单第 0.3 条的顺序口径）

今天量到的事实（编排者 10-12 亲跑，台账 `docs/reports/pending-and-issues.md:16871` 第④行）＝
**310 那枚 `DEGRADED` 分支跑在脚本自己的 GUARD 之前**。⇒ 本腿的新守卫必须**在守卫链之后、在起测试之前**。

| 入口 | 守卫链（现读行号） | 本腿守卫插入点 | 之后才起测试 |
|---|---|---|---|
| `scripts/winsec-tests.sh` | `:61-78` 310 根块（含 DEGRADED）→ GUARD 1 `:89-100` → `GOOS` 空 `:102-106` → 非 windows `:107-114` → `go list` `:116-120` | **`:120` 的 `fi` 之后、`:122` 那行 `echo "…platform=…"` 之前** | `:128` 委派 `portable-tests.sh` |
| `scripts/portable-tests.sh` | 310 根块 `:141-166` → GUARD C `:612-631` → GUARD A `:635-652` → ledger `:666-744` → `:758` 打出 skip pattern | **`:758` 之后、`:760 capture=$(mktemp …)` 之前** | `:764` `sh "$strict"` |
| `scripts/wisp-cli-tests.sh` | 310 根块 `:74-91` → 非 windows `:93-101` → DLL 面 `:106-113`／`:115-120`／`:121-132` | **`:132` 的 `fi` 之后、`:134` 那段 PATH staging 注释之前** | `:146` 委派 `--scope=cli` |
| `tools/d22scan/runtests.sh` | 参数守卫 `:45-66` → 310 根块 `:68-107`（含 DEGRADED） | **`:107` 的 `fi` 之后、`:109 out=$(mktemp …)` 之前** | `:116` `go test` |

自证（派单要求"用一次造形把三行读数的先后顺序打出来"）＝件 `logs/23-order-degraded-then-guard.txt`
与 `logs/24-order-guard-after-guard1.txt`，读数在 `30-teeth.md` §4。

## 3. 改前基线（格四用；⛔ 拿旧数当基线，全部 14:5x 现跑）

件＝`logs/10-pre-build.txt`／`11-pre-gate-d22-len.txt`／`12-pre-census.txt`／`13-pre-winsec-run1.txt`／
`14-pre-winsec-run2.txt`／`15-pre-shape3.txt`（每枚自落 `rc=` 行，无一枚 0 字节）。

- `GOFLAGS= go build ./...` ⇒ 逐字 **`build_rc=0`**（elapsed 5s）
- `sh scripts/d22scan.sh` ⇒ 逐字 **`d22scan_rc=0`**（elapsed 24；同发那行 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`）
- `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ 逐字 **`len_rc=0`**＋`VERDICT GREEN`（分母 `tracked paths=11826`）
- `bash scripts/portable-tests.sh --scope=census` ⇒ 逐字 **`census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`**、**`census_rc=0`**
  ⇒ 与派单给的对照数（`packages=35 …=7 …=7 …=0`）**逐字同**；⚠ 它是现跑的，⛔ 抄来的
- `bash scripts/winsec-tests.sh ./internal/winsec/` 改前两发（默认根形状）⇒ 两发都逐字
  `=== RUN=101  --- PASS=58  --- FAIL=0  --- SKIP=0`、`winsec_rc=0`；那一发用的根逐字
  `GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/winsec-20261012-145224-556` ⇄ `…winsec-20261012-145253-1850`
- 形③ 改前：`GOOS=linux sh tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ **`shape3a_rc=1`**（3s；`go env GOOS` 在设了 `GOOS=linux` 后逐字 `linux`）；
  `LOCALAPPDATA= sh tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ **`shape3b_rc=0`**（22s，同发逐字打出 310 那行
  `DEGRADED - LOCALAPPDATA is empty, so this step has no measured-safe root to build …GOTMPDIR stays unset.`）

## 4. 起手名册（编队三问）与归属

尺①＝`tasklist //FI "IMAGENAME eq wisp.exe"`（走 bash）⇒ 逐字 `INFO: No tasks are running which match the specified criteria.`、
`wisp_tasklist_lines=0`，正控同行 `svchost_tasklist_lines=84`（⛔ 把空输出读成 0 枚）。
尺②＝`powershell -File ps1/roster.ps1`（`Get-Process | Where-Object { $_.ProcessName -like '*.test' }`，
`Get-ChildItem` 那把带 `-ErrorAction SilentlyContinue`）⇒ 起手逐字 `temp_test_star=59`、
`temp_all_dirs_positive_control=2880`、`family_test_exe=1`、`svchost_positive_control=84`、`all_processes_positive_control=369`，
那枚逐字 `family_pid=28152 name=wisp.test started=2026-10-12 14:43:52`。
⇒ **起手就有 1 枚非我这发的 `*.test` 活着**：`wisp.test`（`cmd/wisp` 那枚测试二进制）14:43:52 起，
比我任何一发都早，本腿⛔ 起过 `wisp.test`／`audio.test`（后者 14:53:53 也出现了＝同一批邻居腿）。
归属定式（起进程前名册 ⇄ 杀完后名册作差＋逐枚 `CreationTime`）在 `90-final.md` §4 收口：
本腿自己两枚入口产的名字是 `d22scan.test` 与 `winsec.test`，每次入口退出后名册里都读不到它们。
`%TEMP%` 顶层 `Test*`＝**59**＝票面在册数，⛔ 动、⛔ 删；本腿全程 59 ⇄ 59 未变（件 `logs/15-pre-shape3.txt`
`pre3a temp_test_star=59` ⇄ `post3b temp_test_star=59`，正控 `2887` 同行）。
