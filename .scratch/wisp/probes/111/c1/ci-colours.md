# 111-c1 · 票 111 CI 读数取数件（盘上现量，非判语）

- 腿：`111-c1`（取数腿，非实现者；⛔ 未碰任何产码/测试/`.github`）
- 起手锚：`.scratch/wisp/probes/111/c1/00-anchor.md`（第一笔已单独进仓）
- 取数时刻：`2026-10-08 11:0x +0800`
- 载体：**GitHub Actions**（`gh` 走 `CarlosShao/wisp`）。仓内只有 `.github/workflows/ci.yml`(992 行) 与 `slo-fresh.yml`(83 行)；`find . -maxdepth 2 -iname '*cnb*'` 零命中 ⇒ **cnb 远端没有任何 CI 载体定义可读**，它只是镜像。
- 大输出全部落**仓外** `D:/tmp/wisp111c1/`，只在仓外切；未把整发日志读进上下文。

---

## Q1 远端与未推量（现量）

```
$ git remote -v
cnb     https://cnb.cool/CarlosShao/wisp (fetch|push)
origin  https://github.com/CarlosShao/wisp.git (fetch|push)
```

| 远端 | `refs/heads/dev` tip（`git ls-remote`，走网络） | `refs/heads/master` tip |
|---|---|---|
| `origin` | `cc31526165734e612de297848bb2080bd459ccba` | `5d777f78a3820f94d40c2afdc0da5a4c91dc089a` |
| `cnb` | `c6cf66e64849955bf92a4b3356c096ea81c35563` | `5d777f78a3820f94d40c2afdc0da5a4c91dc089a` |

本地 HEAD 相对每枚远端的未推枚数（逐枚 `git rev-list --count <remote>/dev..HEAD`）：

| 量时点 | HEAD | vs `origin/dev` | vs `cnb/dev` |
|---|---|---|---|
| 锚（10:54） | `7a367a08…` | **295** | **475** |
| 本件落笔（11:0x，含本腿第 1 笔与并行腿的笔） | `e6c3be17…` | **304** | **484** |

- 两枚远端均 `--left-right --count` = `0 <N>` ⇒ 本地**领先且不分叉**（落后 0 枚）。
- **推送状态只按"远端 tip 逐字等于本地 HEAD"这把尺**：`cc31526…` ≠ `e6c3be17…`、`c6cf66e…` ≠ `e6c3be17…` ⇒ **两枚远端都未推平；本腿未 push，也没有任何 run 跑在未推 SHA 上。**

rc=0（Q1 全部命令 rc=0）

---

## Q2 最近 10 发 run 名册（`gh run list --limit 10` + `gh api repos/CarlosShao/wisp/actions/runs?per_page=12`）

| run# | run id | workflow | 触发 | 状态 | 结论 | created（UTC） | headSha | 触发人 |
|---|---|---|---|---|---|---|---|---|
| 305 | 37703959747 | ci | schedule | completed | **failure** | 2026-10-07T23:44:30Z | `cc31526165` | CarlosShao |
| 50 | 37700819596 | slo-fresh | schedule | completed | success | 2026-10-07T23:11:35Z | `cc31526165` | CarlosShao |
| 49 | 37629934350 | slo-fresh | schedule | completed | success | 2026-10-07T13:37:11Z | `cc31526165` | CarlosShao |
| 48 | 37580160337 | slo-fresh | schedule | completed | success | 2026-10-07T06:11:21Z | `cc31526165` | CarlosShao |
| 304 | 37545246395 | ci | schedule | completed | **failure** | 2026-10-06T23:12:27Z | `cc31526165` | CarlosShao |
| 47 | 37511647900 | slo-fresh | schedule | completed | success | 2026-10-06T18:29:01Z | `cc31526165` | CarlosShao |
| 46 | 37424357373 | slo-fresh | schedule | completed | success | 2026-10-06T06:33:02Z | `cc31526165` | CarlosShao |
| 303 | 37406757402 | ci | push | completed | **failure** | 2026-10-06T02:58:39Z | `cc31526165` | CarlosShao |
| 302 | 37406422380 | ci | push | completed | **failure** | 2026-10-06T02:54:30Z | `b948bcb84e` | CarlosShao |
| 301 | 37405698188 | ci | push | completed | **failure** | 2026-10-06T02:45:36Z | `ec84cff16c` | CarlosShao |

- 最近 10 发里 `ci` 占 5 发，**5/5 failure**；`slo-fresh` 占 5 发，5/5 success。
- 再往前扩量核对（同一条尺，`--limit 100`）：**run 206 → 305 共 100 发 `ci`，conclusion 直方图 = `{failure: 100}`，`success: 0`**；窗口 2026-09-23T04:17:52Z 到 2026-10-07T23:44:30Z。⇒ 可读窗口内 `ci` **从未有过一发绿**。
- 所有 run 的 `headSha` 都不晚于 `cc31526`（= `origin/dev` tip）。**没有任何一发 run 跑在本地未推的 296+ 枚上。**

rc=0（`gh run list` 第 1 发因 `runNumber`/`actor` 字段名不存在而 rc=1，改正字段后 rc=0；`--workflow ci --limit 20` 一发遇 `dial tcp 198.18.0.19:443 connectex ... failed to respond` rc=1，**重试一发即通 rc=0**，未据首失败下结论）

---

## Q3 `test-windows` 逐步颜色（最近那发已终态 = run 305 / `37703959747`）

run 305 `test-windows`（job id `113073784763`，结论 failure，runner **`GitHub Actions 1000001277` / labels `['windows-latest']`＝GitHub 托管，不是机主这台机器**）：

| step | 名称 | 颜色 |
|---|---|---|
| 1 | Set up job | success |
| 2 | Run actions/checkout@v4 | success |
| 3 | Run actions/setup-go@v5 | success |
| 4 | Windows ACL sealing gate (internal/winsec's own tests, ticket 110) | success |
| 5 | Cache third_party (deps.toml key) | success |
| 6 | cgo build smoke (build.ps1 fetch-deps + mingw link + doctor) | success |
| **7** | **cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)** | **failure** |
| 8 | Package coverage census (ticket 111 AC#1 + GUARD D) | success |
| **9** | **Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)** | **failure** |
| 10 | PathResolver junction placeholder (real cases tickets 18/20) | success |
| 18/19 | Post Cache third_party / Post setup-go@v5 | skipped |
| 20/21 | Post checkout / Complete job | success |

同一普查把尺放宽到**最近 12 发已终态 `ci` run（294–305）**，`test-windows` 主步（剔 setup/post）三类计数：

| 步 | 出现发数 | success | failure | skipped |
|---|---|---|---|---|
| Windows ACL sealing gate | 12 | 12 | 0 | 0 |
| Cache third_party | 12 | 12 | 0 | 0 |
| cgo build smoke | 12 | 12 | 0 | 0 |
| **cmd/wisp CLI tests** | 12 | **0** | **12** | 0 |
| Package coverage census | **5** | 5 | 0 | 0 |
| Portable windows tests | 12 | 0 | 12 | 0 |
| PathResolver junction placeholder | 12 | 12 | 0 | 0 |
| **winlive compile gate** | **0** | 0 | 0 | 0 |

口径自守：**`success` 不等于"这道门拦过"，只有 `failure` 才是拦下的读数**；上表的"出现过 success 或 failure"只用来判"这一步真在 CI 里执行过"。

其余 5 个 job 的 12 发普查（同一 `jobs-*.json` 集，供台账看"有牙/没牙"）：

| job | 步 | 出现 | success | failure | skipped |
|---|---|---|---|---|---|
| lint | D22 scanner positive control (seeded red) | 12 | 9 | 3 | 0 |
| lint | D22 scanner self-test (161 AC#2) | 12 | 9 | 0 | 3 |
| lint | D22 seven-ban + emoji scan | 12 | 9 | 0 | 3 |
| lint | Tracked path-length budget (262) | 8 | 7 | 0 | 1 |
| lint | **gofmt (gofumpt)** | 12 | 0 | **9** | 3 |
| lint | gofmt (gofumpt) - tracked set denominator (161 AC#7 A) | 12 | 0 | 0 | **12** |
| lint | go vet (module) | 12 | 0 | 0 | **12** |
| lint | go vet (tools/d22scan module) | 12 | 0 | 0 | **12** |
| lint | **staticcheck** | 12 | 0 | **12** | 0 |
| lint | mockllm module vet | 12 | 12 | 0 | 0 |
| test-core | Portable package tests (core scope) | 12 | 0 | **12** | 0 |
| test-core | compose up / probe mock-llm / env fork / stop | 12 | 12 | 0 | 0 |
| lint-frontend | 全部 9 主步（含 4 步 render evidence） | 12 | 12 | 0 | 0 |
| slo-smoke | Build wisp.exe | 12 | 12 | 0 | 0 |
| slo-smoke | SLO smoke gate | 12 | 8 | 4 | 0 |
| slo-full | Build wisp.exe (deps cached) | 12 | 10 | 0 | 2 |
| slo-full | SLO full gate (six states + settle + leak) | 12 | 8 | 2 | 2 |

- `slo-full` 是唯一跑在**自托管**runner 上的 job：`runner_name = wisp-selfhosted-01`，labels `['self-hosted', 'wisp-slo']`。其余 5 个 job 全在托管（`ubuntu-latest` / `windows-latest`）。
- run 305 两枚红的**具名因**（非本票格子，但入台账可省后人一遍）：`gofmt (gofumpt)` 的红只有 1 条命中 `.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`（该 mutant **是被跟踪的**，`git log --diff-filter=A` = `4813567e` 2026-09-28 ⇒ 它在"跟踪集＝分母"的射程里）；`staticcheck` 的红是 `U1000/ST1012/SA1019/SA9009/S1011` 一族（该发日志里 `.go:行:列:` 形状共 55 行）。

rc=0

---

## Q4 `cmd/wisp` 那一步到底跑没跑（run 305，逐步取证）

**答：跑了，而且红。三问逐个答：**

1. **有没有一步真的执行了 go test 打到 `./cmd/wisp/`？** ⇒ 有。run 305 `test-windows` step 7（`Run bash scripts/wisp-cli-tests.sh`，日志边界 = 第 7407–9597 行；下一发 run 300 同步步在第 2617–4742 行）。
2. **红名是什么？** ⇒ run 305（`cc31526`）该步的 6 枚顶层红（逐字）：
   - `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` (1.34s) — `always_write_no_clobber_226_test.go:77: the run did not store the rule it was answering for, so the surviving hand edit proves nothing`
   - `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` (1.07s) — `approval_always_201_test.go:141: config.toml did not gain the stored rule ...`
   - `TestTicket223PermissionDeniedSitsInItsOwnSentence` (1.08s) — `config_reload_perm_223_windows_test.go:95: attempt 1: icacls denied nothing, so this case cannot show the permission sentence`
   - `TestRunPacketCarriesTheLoadedInstructionFiles` (1.07s) — `instructions_200r2_test.go:167: the packet's instructions carry no entry for ...`
   - `TestPanelHostRealWindowHopAndLifecycle` (18.94s) — `panel_host_windows_test.go:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window`
   - `TestAC14GoSideEvalPushReachesThePage` (0.62s) — `panel_resident_windows_test.go:869: Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK"`
   步末四数（脚本自己打的）：`=== RUN=335  --- PASS=230  --- FAIL=6  --- SKIP=2`，`FAIL github.com/CarlosShao/wisp/cmd/wisp 392.928s`，`portable-tests.sh: strict runner exited 1 for scope=[./cmd/wisp/]`，`##[error]Process completed with exit code 1.`
   对照第二发（run 300，`c6cf66e`）：同样 6 枚红，前 4 枚同名，后两枚换成 `TestPanelHostRealWindowHopAndLifecycle` (3.95s) **＋** `TestPanelHostLatencyPercentilesAC2` (0.00s)（这枚在 run 305 是 SKIP，在 run 300 是 FAIL）；四数 `=== RUN=323 PASS=220 FAIL=6 SKIP=1`。
3. **本机那个 `0xc0000135` 坑在 CI 上是不是这形？** ⇒ **不是。**去管道看末几行＋整步计数：
   - `grep -c '0xc0000135'` 在 cmd/wisp 步切片 = **0**，在**整发 run 305 的失败日志（12 190 行 / 2 013 255 字节）里也 = 0**；run 300 的切片同 0。
   - 该步首行有 DLL 就位证据：`wisp-cli-tests.sh: scope=./cmd/wisp/ dll_dir=/d/a/wisp/wisp/third_party/sherpa-onnx (pinned: onnxruntime.dll sherpa-onnx-c-api.dll sherpa-onnx-...)`（`sherpa|dll` 命中 12 行），且前置 step 5 `Cache third_party` 与 step 6 `cgo build smoke` 均 success。
   - `--- FAIL` 有 7 行（6 顶层＋1 包行）、`--- PASS` 328 行 ⇒ 用例**真的跑起来了**，不是"加载期就死、零 `--- FAIL`"那一形。
   - 2 枚 SKIP 具名：`TestPanelHostLatencyPercentilesAC2`（`panel_host_windows_test.go:985: no cold/hot sample recorded in this process`）、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`panel_resident_windows_test.go:317: AC#13 has no subject in this tree: the embed resolves no entry`；同场 `panel_host_gate_test.go:109` 给出形状读数 `shape=anchor-only built=false entry-bytes=0 entry-err=panel: embedded assets are not built (run npm run build)`）。

⇒ 这一格今天的**取得到的读数**＝"CI 真跑了 `cmd/wisp`，12/12 发红，红名册如上"；**"CI 真跑了 `cmd/wisp` 并且绿"今天没有读数**（`success` 计数 = 0/12，且没有任何 run 跑在未推 SHA 上）。

rc=0

---

## Q5 winlive / 真开窗那族在 CI 的颜色（票 111 `AC#11`）

**yaml 里是哪一步（HEAD 原文，只读）**：`.github/workflows/ci.yml:595-644`，`test-windows` 腿内、紧跟 `:593` 的 `run: bash scripts/wisp-cli-tests.sh` 之后：

- `:595` `- name: "winlive compile gate (go vet -tags winlive, ticket 111 AC#11)"`
- `:642` `shell: bash` · `:643` **`if: ${{ !cancelled() }}`（带了）** · `:644` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`
- 无 `continue-on-error`（HEAD 全文件 `^[[:space:]]+continue-on-error:` = 0 命中）、无作业级 `if:`（`^    if:` = 0 命中）。

**两把尺的数（⚠ 别混）**：

| 尺 | 命令 | HEAD | pushed tip `cc31526` |
|---|---|---|---|
| 结构尺·只 `!cancelled()` | `grep -cE '^[[:space:]]+if:[[:space:]]*\$\{\{[[:space:]]*!cancelled\(\)[[:space:]]*\}\}$'` | **12** | 9 |
| 结构尺·全部 `if:` 形状 | `grep -cE '^[[:space:]]+if:'` | **13**（= 12 `!cancelled()` ＋ 1 `always()`，后者是 `:461` test-core `Stop compose services`） | 10 |
| 词频尺（含注释，⛔ 不是结构尺） | `grep -c 'cancelled'` | **27** | — |

逐条列（HEAD）：12 枚步级 `if: ${{ !cancelled() }}` 的**语句行号**是 `:234 :294 :331 :387 :522 :560 :570 :592 :643 :732 :768 :784`，它们各自主人的 `- name:` 行是 `:184 :244 :320 :335 :516 :559 :566 :573 :595 :646 :735 :771`；`always()` 那枚语句在 `:462`，它的主人 `- name: Stop compose services` 在 `:461`（test-core 腿）。**没有任何一步能凭 `grep` 被读成"跑过"**——名册里 `- name: go vet (module)`(`:237`) 与 `- name: go vet (tools/d22scan module)`(`:240`) 以及另外 22 枚步（共 24 枚）**不带任何 `if:`**，其中这两枚 `go vet` 的 12 发真实颜色是 `[skipped]×12`＝**从未执行**（这正是"在 yaml 里"≠"跑过"的那一形；`skipped` 的成因见 N5）。

**出现过颜色没有？** ⇒ **没有。零。**
- 该步由 commit `351e5a5e8582c1b0ad1092a1d91a73ef5c21e04f`（2026-10-08）加进 `ci.yml`。`git merge-base --is-ancestor 351e5a5e cc31526` ⇒ **NOT ancestor**（不在已推 tip 里）。
- `git show cc31526:.github/workflows/ci.yml | grep -cE 'winlive|-tags'` = **0**。
- 步名差集（HEAD 有 / pushed tip 无）恰为两枚：`winlive compile gate (...)` 与 `Portable tests carrier self-test (ticket 111 AC#3 - GUARD D's positive control)`（`ci.yml:335`，lint 腿）。
- 最近 12 发已终态 `ci` run（294–305）的 job/step 全量里，名字含 `winlive` 的步出现 **0 次**。

⇒ 票面 `AC#11` 的 **ⓑ 反形在 CI 上今天无读数**；`ⓐ` 只有本机读数（票面/前腿件 `probes/111/r6/logs/vet-base.txt`、`orch-vet-winline.txt`，本腿未复测、⛔ 未跑 Go）。`ⓒ` 射程口径本腿不重述判语。

**"真开窗"那族在 CI 的颜色（同一步里已经取得到的两枚，具名）**：
- `TestPanelHostRealWindowHopAndLifecycle` = **failure**（run 305 与 run 300 各一枚，托管 windows runner；末行读数 `cold bring-up measured on this box: -1.000 ms` ⇒ "the message channel did not come up on the real window"）。
- `TestAC14GoSideEvalPushReachesThePage` = **failure**（run 305；`title=""` vs want `PUSHED-33R5-OK`）。
- `TestPanelHostLatencyPercentilesAC2` = run 300 **failure** / run 305 **skipped**；`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` = 两发都 **skipped**（原因同上：embed 无 entry）。
- ⚠ 这几枚都**不是** `winlive` 标签族（`winlive` 那 12 枚今天连编译都没被 CI 查过），它们是 `cmd/wisp` 常规 windows 测试；按 `AC#11 ⓒ`⛔ 不许把它们读成 winlive 的 CI 载体，也不许把它们读成票 35 `:52`/`:75(c)` 的真窗读数已闭合。

rc=0

---

## Q6 今天任何腿都取不到的格（具名，逐格写清欠哪一发读数）

| # | 欠的格（具名） | 欠的具体读数 | 为什么今天取不到 |
|---|---|---|---|
| N1 | 票 111 `AC#11` winlive 编译门的**首次 CI 颜色** | 一发 `ci` run（workflow `ci`、`test-windows` 腿）里名为 `winlive compile gate (go vet -tags winlive, ticket 111 AC#11)` 的步给出 `success` 或 `failure`；该 run 的 `headSha` 必须包含 `351e5a5e` 或其后 | 零 push：`351e5a5e` 不在 `origin/dev` tip（`cc31526`），`grep -cE 'winlive\|-tags'` 在已推 yaml = 0；12 发可读 run 里该步出现 0 次 |
| N2 | 票 111 `AC#11` **ⓑ 反形在 CI 上**的颜色 | 同一门在一发真实 run 里因一枚 winlive 文件坏而 `failure` 的读数（票面逐字要"那一步"红，不是"那条命令"红） | 同 N1；本腿⛔不改产码也不跑 Go，且这一步连一次 `success` 都还没有 |
| N3 | 票 111 `AC#4` 的另一半：**"CI 真的跑了 `cmd/wisp` 并且绿"** | `cmd/wisp CLI tests` 步一发 `success` 读数（含步末四数 `FAIL=0`） | 该步 12 发全 failure（红名册见 Q4）；要它绿必须先修那 6 枚红，然后**推送**才可能有读数 |
| N4 | 票 111 `AC#3` GUARD D 载体自检（`Portable tests carrier self-test`，`ci.yml:335`）的**首次 CI 颜色** | lint 腿该步在一发真实 run 里 `success`/`failure` | 同 N1：该步只存在于未推 SHA（步名差集已证） |
| N5 | `ci.yml` 里 `go vet (module)`(`:237`) 与 `go vet (tools/d22scan module)`(`:240`) 的**任一颜色** | 这两步在一发真实 run 里出现过 `success` 或 `failure`（今天 12/12 `skipped`＝从未执行） | 它们不带 `if:`，前面 `gofmt (gofumpt)` 12 发里 9 failure / 3 skipped ⇒ 永远被吞；需要推一发且让 `gofmt (gofumpt)` 绿，或给这两步补 `!cancelled()`（＝改 `.github`，本腿⛔无权） |
| N6 | `gofmt (gofumpt) - the tracked set is the denominator (161 AC#7 A)` 的颜色 | 同上，12/12 `skipped` | 该步带 `!cancelled()`(HEAD `:234`)，但这把 `if:` 是在 294–305 这些已推 SHA **之后**才落的两枚新增（pushed tip 结构尺 9 vs HEAD 12）⇒ 读数欠推送 |
| N7 | **`ci` 工作流的一发绿**（任何门"拦干净了"的全局读数） | 任意 `ci` run `conclusion=success` | 可读窗口 100/100 failure（run 206–305）；当前 HEAD 未推，且 lint/test-core/test-windows 三条腿同时红 |
| N8 | 本地 HEAD（`7a367a0` / 本件落笔时 `e6c3be17`）上的**任何** CI 读数 | 一发 `headSha` 等于未推 HEAD（或其任一祖先≥ `351e5a5e`）的 run | 零 push ⇒ 这样的 run 不存在；`gh` 只能读已推 SHA 的 run |
| N9 | **cnb 远端那一侧的 CI** | `cnb/dev` tip `c6cf66e` 上任何流水线读数 | 仓内没有 cnb 流水线定义（`find -maxdepth 2 -iname '*cnb*'` 0 命中，只有 `.github/workflows/` 两枚 yml）；`gh` 只读 GitHub ⇒ 无 CI 载体可取 |
| N10 | 票 35 `:52` / `:75(c)` 的**真开窗读数** | 冷/热开窗毫秒级样本（`TestPanelHostLatencyPercentilesAC2` 需要 lifecycle 同场跑出样本） | 按 `A690` 走〔仅本机可量〕那一支，CI 上没有载体；CI 侧今天只取得到它的 `skipped`/`failure`（见 Q5），⛔ 不许当成真窗读数 |

rc=0

---

## 具名顶回（编排者转述 vs 盘上/CI 原文，一律以原文为准）

1. **「10-06 曾推过一批（约 175 枚），之后本地又新增了很多」** ⇒ 盘上现量：相对 `origin/dev` 未推 **295 枚**（锚时点）/ **304 枚**（落笔时点），相对 `cnb/dev` **475 / 484 枚**。"175" 不进任何结论。
2. **「哪几格因为 runner 就是机主这台开发机而取不到」** ⇒ `test-windows` / `slo-smoke` 跑在 **GitHub 托管 `windows-latest`**（`runner_name = GitHub Actions 1000001277/1000001278`），`lint`/`test-core`/`lint-frontend` 跑在托管 `ubuntu-latest`；**只有 `slo-full` 是自托管**（`wisp-selfhosted-01`，labels `['self-hosted','wisp-slo']`）。⇒ winlive 门的取不到**只因零 push**，不是因 runner＝开发机；推上去它就会在托管 windows 上给出颜色。
3. **「数 `if: !cancelled()` 时两种尺差 1」** ⇒ 盘上实测：结构尺 `!cancelled()` = **12**、结构尺全部 `if:` = **13**（差 1，与那句一致）、词频尺 `grep -c 'cancelled'` = **27**（差 15，因为注释命中）。三种尺的公式与命令都已在 Q5 写出，台账引用时请具名尺子。
4. **票文件名 `111-ci-tests-20-of-33-packages.md` 的 "20/33"** ⇒ 本腿未复核该分母（复核要 `go list ./...` 计数，属 `33-v4` 独占面）；本件不使用该数。

---

## 本腿自报：读命令与豁免

- ⛔ 未跑 `go build` / `go vet` / `go test`；**也未跑 `go env` / `go list`**（本件任何结论都不依赖它们）。Go 编译面与卫生仪器留给 `33-v4`。
- 用过的写操作：`git add -- <显式 path>` ＋ `git commit -- <显式 path>`（共 2 笔：锚与本件＋票末追加），⛔ 无 `add -A`/`add .`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`，⛔ 未 push，⛔ 未开真窗，⛔ 未改 `.github/**`（只 `git show` 与只读 `grep`）。
- 大输出落点（仓外，只建不删）：`D:/tmp/wisp111c1/` 下 `run305-logfailed.log`(2 013 255 B/12 190 行)、`run300-logfailed.log`(1 922 082 B/11 672 行)、`wispcli-step.log`、`winport-step.log`、`wispcli-run300.log`、`run305-jobs-api.json`、`jobs-*.json` × 11、`apiruns.json`、`ci20.json`、`ci100.json`、`steps-head.md`、`steps-pushed.md`。⛔ 未在仓库目录内建 worktree/checkout。
- 未在件里内联任何整发日志；日志一律先落仓外再 `grep`/`wc`。

rc=0
