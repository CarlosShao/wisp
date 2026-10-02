# vm-draw-1：schedule 那发（`36940536372`）的逐名读数 ＝ "机器/负载抽签 vs 这版码就是红" 的分水岭

> 本腿是**只读取证腿**。零 `go test` / 零 `go build` / 零 `go vet` / 零 exe 执行 / 零 CPU 负载
> （题面硬边界 #1，盘上三枚写码腿 `198-r1`／`33-r8`／`250-r1` 在跑毫秒级计时）。
> 零 rerun、零 cancel、零 workflow 触发、零 push、零为采样而推码（硬边界 #2）。
> 全部结论只来自：`gh run/api` 的**读**操作 ＋ 读盘上已有的 CI 日志 ＋ `git show/log/diff/grep`（只读）。
> 搜索根**逐发显式限定**在 `cmd internal tools docs .scratch` 或单文件；**`frontend/**`／`design/**` 零读取、零引用、零写入**（硬边界 #3）。

---

## §0 起手锚

### 0.1 本机锚（自取，未拿题面任何 sha 当锚）

| 尺 | 读数 |
|---|---|
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git rev-parse HEAD` | `3b23613b0e681fb6c5116b464e1d3fabedbe626f` |
| `git log -1 --format='%H %ad' --date=iso` | `3b23613b… 2026-10-02 09:13:14 +0800`（＝前一枚腿 `ci223-1` 那一枚） |
| `date` | `2026-10-02 09:15:20 +0800` ＝ `2026-10-02T01:15:19Z` |
| `gh auth status` | Logged in as account `CarlosShao`（keyring），protocol https，token scopes `gist, read:org, repo` |
| `gh repo view --json defaultBranchRef` | `{"name":"dev"}` ⇒ **默认分支＝`dev`**；仓名 `wisp`，owner `CarlosShao` |
| `git remote -v` | `origin` = `https://github.com/CarlosShao/wisp.git`、`cnb` = `https://cnb.cool/CarlosShao/wisp` |

⚠ 起手锚 `3b23613b` **不等于**任何一发 CI 编译的 commit。本腿引 CI 里的行号一律来自 CI 日志本身；
引源码行号一律钉在 `8ae4c23e`，钉法 `git show 8ae4c23e:<file> | grep -n <串>`（0.-1 那条坑本腿照防）。

⚠ **写面过程中 HEAD 果然漂走了**（前一枚腿报过的那枚坑，本腿复现）：交件前复量
`git rev-parse HEAD` = **`93f2b0c94e79363b2493d6645dbf3641b4fe336c`**（起手是 `3b23613b`）。
⇒ 本文件里**任何**一个源码行号都不是"在 HEAD 上读的"，全部出自 `git show 8ae4c23e:<file>` 或 `git log/diff 0589fd9c..8ae4c23e`。
下一枚腿照这条：**在共享树里引行号前先自问"我读的是哪个 sha 的那份字节"。**

### 0.2 三发 run 的元数据（本腿自取，`gh run view --json`，取数 `2026-10-02 09:15 +08`）

| run | workflow | event | headBranch | headSha | status | conclusion | createdAt | updatedAt | attempt |
|---|---|---|---|---|---|---|---|---|---|
| `36789659174` | `ci` | **schedule** | `dev` | `0589fd9cc3c48eab85950475d25505f447c0a425` | completed | failure | `2026-09-30T23:10:15Z` | `2026-10-01T01:32:14Z` | 1 |
| `36889094435` | `ci` | **push** | `dev` | `8ae4c23e2ab4e1bda87ba6d280430db75527670b` | completed | failure | `2026-10-01T16:03:34Z` | `2026-10-01T16:17:44Z` | 1 |
| `36940536372` | `ci` | **schedule** | `dev` | `8ae4c23e2ab4e1bda87ba6d280430db75527670b` | completed | failure | `2026-10-01T23:23:21Z` | `2026-10-02T00:26:02Z` | 1 |

**复认题面：`36940536372` 的 `workflowName:"ci"`／`event:"schedule"`／同 `headSha 8ae4c23e`／`completed`／`failure` 五字全对。**
补一条题面没写的：**"基线"那一发 `36789659174` 也是 `event:"schedule"`**（不是 push）。
⇒ 三发里两发是定时、一发是推送；`conclusion:"failure"` 是**整发**的判语（三发都红，红在别的格），
**与本枚用例的 PASS/FAIL 是两回事**，别把 run 级 failure 读成用例级 failure。

### 0.3 本腿写面（硬边界 #5，只建不删）

唯一写面目录＝`.scratch/wisp/probes/orchestrator/vm-draw-1/`。产出件：

| 件 | 是什么 |
|---|---|
| `read.md` | 本文 |
| `logs_sched_all.txt` / `.err` | `gh run view 36940536372 --log` 全文（6 612 行 / 1 073 604 B，err 0 B） |
| `job_<作业名>.txt` / `.err` | 六发**逐作业** `gh run view --log --job <id>`（可用性证据本体，6 枚 err 全 0 B） |
| `clean_tw_sched.txt` | `test-windows` 去时间戳/去 TSV 列后的日志（4 454 行 / 551 748 B） |
| `clean_tw_push.txt` / `clean_tw_base.txt` | 同尺剥 `ci-delta-1/raw_tw_new.txt`、`raw_tw_base.txt`（4 460 / 4 069 行） |
| `clean_all_sched.txt` / `clean_all_push.txt` / `clean_all_base.txt` | 三发整份日志剥壳（6 612 / 6 338 / 10 923 行） |
| `clean_tc_sched.txt` / `clean_tc_push.txt` | `test-core` 那两发的整作业日志 |
| `strip.awk` / `stepstats.awk` / `whichstep.awk` / `cmduisp_fail.awk` / `winscope_fail.awk` / `beforeours.awk` / `tsv_tw_{push,base}.txt` / `up_*.txt` / `ws_*.txt` / `sched_*_cmdwisp.txt` | 尺与本腿名册读数（临时件只建不删） |

---

## §1 那一发 schedule **读到了**，而且六发日志全在

### 1.1 逐作业点名（run `36940536372`）

尺：`gh run view 36940536372 --json jobs` ＋ `gh api repos/CarlosShao/wisp/actions/runs/…/jobs` ＋
**逐作业 `gh run view --log --job <id>`，可用性看 rc 与 stderr 字节数**（不是看 `--log-failed`，那把尺少算，见 1.3）。

| 作业 | job id | 结论 | runner_name | labels | started → completed (UTC) | 日志 | 字节 | stderr |
|---|---|---|---|---|---|---|---|---|
| `test-windows` | 110630736700 | **failure** | `GitHub Actions 1000001162` | `windows-latest` | 23:23:24 → 23:35:23 | **在** | 796 721 | 0 B |
| `slo-full` | 110630736901 | success | `wisp-selfhosted-01` | `self-hosted,wisp-slo` | **10-02T00:23:22 → 00:26:01** | **在** | 34 638 | 0 B |
| `slo-smoke` | 110630736919 | success | `GitHub Actions 1000001163` | `windows-latest` | 23:23:25 → 23:25:28 | **在** | 36 802 | 0 B |
| `lint` | 110630736958 | failure | `GitHub Actions 1000001164` | `ubuntu-latest` | 23:23:25 → 23:24:36 | **在** | 118 462 | 0 B |
| `test-core` | 110630737185 | failure | `GitHub Actions 1000001166` | `ubuntu-latest` | 23:23:24 → 23:24:20 | **在** | 37 159 | 0 B |
| `lint-frontend` | 110630737210 | success | `GitHub Actions 1000001165` | `ubuntu-latest` | 23:23:24 → 23:23:45 | **在** | 49 822 | 0 B |

⇒ **六枚作业日志全部可读，`log not found` 零枚。** 整份日志（`--log` 不带 `--job`）的 stderr 也是 0 字节 ⇒ 交叉印证。
本枚要的两枚（`test-windows` ＝两枚目标用例的所在作业）**日志完整**。**这一发放买到了它该买的东西。**

### 1.2 "没日志"在本仓是什么形状——顺手把对照钉住

| run | 缺日志的作业 | 尺的读数 |
|---|---|---|
| `36789659174`（base） | 无 | `ci-delta-1/logs_base_all.err` ＝ **0 B** |
| `36889094435`（push） | **`slo-full`（job `110459845223`）** | `ci-delta-1/logs_new_all.err` ＝ **28 B**，逐字 `log not found: 110459845223` |
| `36940536372`（sched，本腿） | 无 | 六枚 per-job `.err` ＋ 整份 `.err` 全 **0 B** |

⇒ 上一发 `slo-full` 那次"没日志"是**自托管 runner 侧**的形状（`wisp-selfhosted-01`，本腿在 §4 里把它和 host 身份对上）。
⛔ 本腿没有把任何一次"没日志"读成红或绿。

### 1.3 本腿自己撞到的读数坑（下一枚腿照这个防）

1. **`gh run view --log --job` 的输出是 TSV：`<作业名>\t<步骤名>\t<时间戳><正文>`。** 时间戳**不在行首**，
   拿 `sed 's/^[0-9]\{4\}-…Z //'` 直接剥行首**剥不动**（本腿第一发就这么着，产出一枚 4 537 B 的废件
   `clean_tw_push.txt` 的第一版）。必须按 TSV 取第 3 列起再剥。**本腿留的可用尺＝`strip.awk`**。
2. 同一个坑的**反向形态**：`ci-delta-1/raw_tw_{base,new}.txt` 与 `logs_{base,new}_all.txt` **不是同一种壳**
   （前者是裸时间戳行、无 TSV 列；后者有）。⇒ 拿一枚 awk 通吃三发会静默吞掉两发的内容而**不报错**。
   本腿的处置：裸壳那两枚另用 `sed` 剥，并把两法产出的**同一枚用例判决行并排比对**过一遍（§2 表里三发行号
   1885 / 1896 / 1885 就是这么互证的）。
3. **`gh api repos/…/actions/jobs/<id>` 的 `.runner.name` 是 `null`**（本腿实测六枚全 null），
   前一枚腿那句 `sl=GitHub Actions 1000001157` 出自**另一条路**：
   `gh api repos/CarlosShao/wisp/actions/runs/<run>/jobs` 的 **`.runner_name`**（下划线形）。
   用错那条就取不到机名，"是不是同一台机器"这一问直接答不了。
4. ⚠ **前一枚腿 `ci223-1` 报的 `raw_testcore_new.txt` 是一把少算的尺**：本腿从**整份日志**数
   push 那发的 `test-core` ＝ `=== RUN=3 / --- PASS=1 / --- FAIL=0`，而那份切片件数出 **0 / 0 / 0**。
   ⇒ 它少了一枚 preflight。见 §4.2。

### 1.4 机器身份：三发是**三台**不同的 VM（这一条决定 §3 的成色）

尺：`gh api .../runs/<r>/jobs --jq '.jobs[] | .name + .runner_name'`（列 `test-windows`）＋ 各作业日志头部的
`Worker ID` / `Azure Region` / `Image` / runner 版本（本腿逐发 `grep -aoE`）。

| run | 事件 | 作业号（runner_name） | Worker ID | Azure Region | Hosted Compute Agent | **Runner Image 版本** | runner 版本 |
|---|---|---|---|---|---|---|---|
| `36789659174` | schedule（`0589fd9c`） | `GitHub Actions 1000001150` | `{165450a5-9cd6-4bfa-8abb-f42678c74be5}` | `northcentralus` | `20260828.587` | **`20260922.246.2`** | `2.337.0` |
| `36889094435` | push（`8ae4c23e`） | `GitHub Actions 1000001157` | `{b4efc40e-f741-49f3-b35e-f03cebedce44}` | **`westus2`** | `20260901.588` | **`20260925.250.1`** | `2.337.0` |
| `36940536372` | schedule（`8ae4c23e`） | `GitHub Actions 1000001162` | `{48178208-bca3-4d08-b56b-201759e90be3}` | `northcentralus` | `20260901.588` | **`20260925.250.1`** | `2.337.0` |

⇒ **三枚 runner 名不同、三枚 Worker ID 不同 = 三台不同的 hosted-compute agent。不是同一台机器。**
所以这一发的判别力**不打折**（题面 2 "若同机要打折"那一支不适用）。
三处要把话说死：
- schedule 那发落回 `northcentralus`（与基线同区），push 那发在 `westus2` ⇒ 区不是因，**同区内也会抽签**。
- **`windows-2025-vs2026` 镜像与 runner 版本（`2.337.0`）三发逐字相同**，但**镜像小版本换了**：
  基线 `20260922.246.2` → 两发新码 **都是 `20260925.250.1`**（Hosted Compute Agent 同步 `20260828.587` → `20260901.588`）。
  ⇒ **这条题面与前一枚腿都没报，必须入账**："基线绿 vs 新发红"这一对上，除代码之外**还有一枚镜像更新在同时变**。
  ⚠ 它**不动 §3 的判语**：§3 用的那一对（push vs sched）**同镜像、同 SHA、不同 VM**，镜像这一维在两发间是**常量**。
  它动的是"基线→新发"的归因：那一档差异**分不出"码"与"镜像"**（§5 C2 已据此改写）。
  ⛔ 本腿没有任何 CPU/内存遥测，也没有 0922→0925 那张镜像的 changelog，**不猜镜像里变了什么**。

---

## §2 逐名读数（本条唯一有意义的尺）

**口径先写死（不写口径这些数不可比）：**
- **判红只认行首 `--- FAIL: Test`**（剥壳后的行首，`^-{3} FAIL: Test`）。`=== RUN` 与 `t.Logf` 都带 `file.go:NN:` 前缀，
  **不当判决行**（前一枚腿在 `clean_new.txt:2170` 上栽过：那是 RUN 行、真判决在 `:2185`）。
- 判决行的**耗时取括号原值**，不四舍五入、不换算。
- 行号＝**本腿盘上件的行号**，与本腿件的生成尺一起给（§0.3），换一枚腿复算必须换尺重出。
- 同一枚用例的**每一次出现都点名**（含"另一发作业里还有第二枚读数"这种，见 2.3）。

### 2.1 `TestTicket223HandEditedFsLooseningCostsAnL2Card`（`cmd/wisp/config_reload_223_test.go:302/:305` 那两行哨兵）

所在作业/步骤：**`test-windows` → `##[group]Run bash scripts/wisp-cli-tests.sh`**（即 `./cmd/wisp/`）。

| 发 | run / 事件 / headSha | runner 作业号 | 判决行（逐字） | 耗时 | 本腿件:行 |
|---|---|---|---|---|---|
| 基线 | `36789659174` / schedule / `0589fd9c` | `GitHub Actions 1000001150` | `--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.13s)` | **2.13s** | `clean_tw_base.txt:1885` |
| 新发 A | `36889094435` / push / `8ae4c23e` | `GitHub Actions 1000001157` | `--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.43s)` | **2.43s** | `clean_tw_push.txt:1896` |
| **新发 B（本腿取的那一发）** | `36940536372` / schedule / `8ae4c23e` | `GitHub Actions 1000001162` | `--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.28s)` | **2.28s** | `clean_tw_sched.txt:1885` |

**同码（`8ae4c23e`）两发：一发 FAIL、一发 PASS。⇒ 见 §3。**

哨兵本身（`config_reload_223_test.go:NNN` 报出行）逐发点名，尺
`grep -aoE "config_reload_223_test\.go:[0-9]+" <件> | sort | uniq -c`：

| 发 | 读数 |
|---|---|
| 基线 | **零命中** |
| 新发 A（push） | `1 config_reload_223_test.go:303` / `1 config_reload_223_test.go:306` |
| **新发 B（sched）** | **零命中** ⇒ `:302` 与 `:305` 两个 `if !strings.Contains(...)` 都没响 |

⇒ 与前一枚腿 `ci223-1` §0.3 的复认**同形**（那两行确实是 push 那发唯一的两条报文）。
⚠ 题面写的哨兵是 `:302/:305`，CI 报的是 `:303/:306` —— **两个都对，不是矛盾**：`:302`/`:305` 是
`if` 那一行，`:303`/`:306` 是 `t.Errorf` 那一行，`t.Errorf` 报的是**自己所在行**。本腿在 `8ae4c23e` 上现读确认。

### 2.2 `TestResolvePerCallBudget`（`internal/risk/pathresolver_budget_norace_test.go`，纯挂钟钉）

所在作业/步骤：**`test-windows` → `##[group]Run bash scripts/portable-tests.sh --scope=windows`**。
（⚠ 它**不**在 `test-core`，也**不**在 `wisp-cli-tests.sh`；本腿逐作业点名确认，两处独立路数同读数。）

| 发 | run / 事件 / headSha | runner 作业号 | 判决行（逐字） | 耗时 | 仪器自己吐的那一行（逐字） | 本腿件:行 |
|---|---|---|---|---|---|---|
| 基线 | `36789659174` / schedule / `0589fd9c` | `…1150` | `--- PASS: TestResolvePerCallBudget (1.81s)` | 1.81s | `pathresolver_budget_norace_test.go:34: C26 Resolve: 961784 ns/op = 0.962 ms/op (budget 1.000 ms, 1677 samples)` | `clean_tw_base.txt:3391`（Logf 在 `:3389`） |
| 新发 A | `36889094435` / push / `8ae4c23e` | `…1157` | `--- FAIL: TestResolvePerCallBudget (2.51s)` | 2.51s | `pathresolver_budget_norace_test.go:34: C26 Resolve: 1241491 ns/op = 1.241 ms/op (budget 1.000 ms, 1902 samples)`<br>`:37: C26 budget breach: Resolve averages 1.241491ms per call, budget 1ms` | `clean_tw_push.txt:3756`（Logf `:3753`、breach `:3755`） |
| **新发 B** | `36940536372` / schedule / `8ae4c23e` | `…1162` | `--- FAIL: TestResolvePerCallBudget (2.55s)` | 2.55s | `pathresolver_budget_norace_test.go:34: C26 Resolve: 3346261 ns/op = 3.346 ms/op (budget 1.000 ms, 560 samples)`<br>`:37: C26 budget breach: Resolve averages 3.346261ms per call, budget 1ms` | `clean_tw_sched.txt:3750`（Logf `:3747`、breach `:3749`） |

**这一枚在两发同码上都红 ⇒ 与 §2.1 不同形。** 但它的"红"不是断言之争，是**仪器直读**：

| 发 | headSha | ns/op | samples | 判决 |
|---|---|---|---|---|
| 基线 | `0589fd9c` | 961 784 | **1 677** | PASS（0.962 / 预算 1.000 ⇒ **只剩 3.8% 余量**） |
| 新发 A | `8ae4c23e` | 1 241 491 | **1 902** | FAIL |
| 新发 B | `8ae4c23e` | **3 346 261** | **560** | FAIL |

⇒ **同一枚 SHA、同一枚测试二进制，两台 VM 测出 2.69 倍的每次调用成本、3.4 倍的吞吐（1902 vs 560 samples）。**
`samples` 就是 `testing.Benchmark` 在固定目标时长内塞进多少次迭代 —— **它是机器速度的自报数，不是推断**。
本腿把它当作"这台机器当时就是慢"的**一级读数**用。

行号钉锚（尺 `git show 8ae4c23e:internal/risk/pathresolver_budget_norace_test.go | grep -nE …`，本腿实跑）：
`:12 const resolveBudget = time.Millisecond`、`:34 t.Logf(...)`、`:36 if time.Duration(ns) > resolveBudget {`、
`:37 t.Fatalf("C26 budget breach: …")`、`:40/:41` 是第二道 50-call 包络（**这一发没响**，`:37` 已 Fatal）。
⚠ 题面给的"`:12/:36`"——`:12` 对，**判决报文实际出自 `:37`**（`if` 在 `:36`）。前一枚腿 §3(c) 把
`t.Fatalf` 记在 `:36` 同一行，也是差一行。**换尺不换人，读数就得带行号来源。**

### 2.3 第二枚 `TestResolvePerCallBudget` 样本：**只有基线那发有**

尺：整份日志按作业列 `--- (PASS|FAIL): TestResolvePerCallBudget`（`awk -F'\t' '$1=="test-core" …'`）。

| 发 | `test-windows`（scope=windows） | `test-core`（scope=core，ubuntu） |
|---|---|---|
| 基线 `0589fd9c` | PASS 1.81s @ `…1150` | **PASS (1.97s)** @ `GitHub Actions 1000001152`，逐字 `pathresolver_budget_norace_test.go:34: C26 Resolve: 724 ns/op = 0.001 ms/op (budget 1.000 ms, 1675071 samples)`（`clean_all_base.txt:7476/:7478`） |
| 新发 A push | FAIL 2.51s | **无样本**（`test-core` 整段没跑 `go test`，见 §4） |
| 新发 B sched | FAIL 2.55s | **无样本**（同上） |

⚠ **这一格顺带钉出一个新事实（题面与前一枚腿都没说）：**同一份码（基线 `0589fd9c`）在
Linux 上测出 **724 ns/op / 1 675 071 samples**、在 Windows runner 上测出 **961 784 ns/op / 1 677 samples** ——
**差 1 328 倍**。⇒ 这枚 1 ms 预算钉在 Windows 上**几乎没有余量**（基线 0.962 ms ＝ 预算的 96.2%），
它的 Windows 档**天然是一枚按机器抽签翻面的用例**，与代码改没改**无关**。
⇒ 题面那句"它只可能因机器慢而红"本腿**复认成立，而且比题面更强**：它在 Windows 档离预算只有 3.8%，
一次普通的机器抖动就够翻它。用它当"码改坏"的证据是**错的用法**。

### 2.4 两发同码的**全步名册**比对（防"只是这一枚巧合同向"）

尺：`stepstats.awk`（按 `##[group]Run` 归属）与 `cmduisp_fail.awk`/`winscope_fail.awk`（步内 `--- FAIL:`/`--- SKIP:` 逐名）。

| 步骤（`test-windows`） | 基线 `0589fd9c` | 新发 A push | 新发 B sched |
|---|---|---|---|
| `scripts/winsec-tests.sh` | RUN=101 PASS=58 FAIL=0 SKIP=0 | RUN=101 PASS=58 FAIL=0 SKIP=0 | RUN=101 PASS=58 FAIL=0 SKIP=0 |
| `scripts/wisp-cli-tests.sh`（`./cmd/wisp/`） | RUN=191 PASS=110 FAIL=7 SKIP=0 | RUN=241 **PASS=147 FAIL=11 SKIP=2** | RUN=241 **PASS=148 FAIL=11 SKIP=1** |
| `portable-tests.sh --scope=windows` | RUN=454 PASS=311 FAIL=12 SKIP=1 | RUN=468 PASS=320 FAIL=13 SKIP=1 | RUN=468 PASS=320 FAIL=13 SKIP=1 |
| `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` | 1/1/0/0 | 1/1/0/0 | 1/1/0/0 |

（`PASS/FAIL/SKIP` 取 `portable-tests.sh` 自己那行"four numbers"，与本腿 awk 逐名数**同读数**：
sched `3146:… === RUN=241  --- PASS=148  --- FAIL=11  --- SKIP=1`、push `3148:… 241/147/11/2`。）

- `--scope=windows` 那 468 枚：push 与 sched 的 **FAIL+SKIP 逐名集合完全相同**（15 行，`diff` 空输出）。
  基线→新发唯一的**名字**变化就是 `TestResolvePerCallBudget`（PASS→FAIL）。
- `./cmd/wisp/` 那 241 枚：push 与 sched 的差集**只有两处**——
  `--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card`（push 独有，sched 里它 PASS）；
  `TestPanelHostLatencyPercentilesAC2` **push 里 SKIP ×2、sched 里 FAIL**（于是 FAIL 数都是 11，但成员不同）。
  其余 10 枚 FAIL **两发同名同判**。
- **跑序上游同样零差**：本枚之前该步的顶层 `=== RUN` 三发**都是 11 枚、逐名逐序 `diff` 空输出**
  （`up_base.txt` / `up_push.txt` / `up_sched.txt`；本腿的尺只数 `^=== RUN   Test[^/]*$`，
  与前一枚腿那个"13 枚"的口径不同，因为它把 `=== RUN` 的子测试也带进去了）。
- 包级耗时（`^(ok|FAIL)\s+github.com/…` 行）：`cmd/wisp` **424.719s → 579.297s → 553.675s**；
  `internal/risk`（windows scope）**6.244s → 6.247s → 9.430s**。
  ⇒ sched 比 push **整体更快**（cmd/wisp 快 4.4%）**却**在 `internal/risk` 上慢 50%。
  **机器不是一根统一的比例尺，是每段各自的抽签。**

### 2.5 `TestPanelHostLatencyPercentilesAC2` 那一格（§2.4 差集里的第二枚，顺手钉成旁证）

同码、两发、不同形态：

- push（`…1157`）：`panel_host_windows_test.go:544: cold bring-up measured on this box: -1.000 ms …` ＋
  `:546: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window`
  ⇒ 同进程没记账样本 ⇒ `:794 t.Skipf("no cold/hot sample recorded in this process …")` ⇒ **SKIP**。
- sched（`…1162`）：`:544: cold bring-up measured on this box: 3496.853 ms …` ＋
  `:549: cold bring-up 3496.9 ms exceeds D32 panel cold budget 1500 ms` ⇒ 有样本 ⇒
  `:808 … cold P50=3496.853 P95=3496.853 (budget 1500)` ⇒ `:814` `t.Errorf` ⇒ **FAIL (0.00s)**。

（尺：`git show 8ae4c23e:cmd/wisp/panel_host_windows_test.go` 现读 `:785 func TestPanelHostLatencyPercentilesAC2`、
`:793 if len(cold) == 0`、`:813 if coldP95 > 1500` —— 判决报文出自 `:794`/`:814`。）

⇒ **一台 VM 上 WebView2 冷拉起压根没回来（-1.000 哨兵），另一台上回了但花 3.5 s（预算 1.5 s）。**
同码两发：**一枚 SKIP、一枚 FAIL**，且两枚都不是产品语义。这是第三枚"按机器翻面"的用例。

---

## §3 分水岭判定：**甲**

> 判据（题面 §3，只允许三种之一）：同一枚 SHA 在两发里，`TestTicket223HandEditedFsLooseningCostsAnL2Card`
> 一发红一发绿 ⇒（甲）；两发都红 ⇒（乙）；读不到 ⇒（丙）。

### 判定＝**甲：同码一发红一发绿 ⇒ "码改坏"这一支死，"机器/负载抽签"这一支活。**

凭据（三样，缺一不成立，逐条可复核）：

1. **同码。** 两发的 `headSha` 逐字相同：`8ae4c23e2ab4e1bda87ba6d280430db75527670b`
   （尺：`gh run view --json headSha`，两发输出字符串相等；本腿 §0.2 表已并排贴出）。
2. **不是同一台机器，所以"绿"不是机器被换回来造成的假绿；而"红"也不是机器的确定属性。**
   三枚 runner 名（`…1150` / `…1157` / `…1162`）＋ 三枚 Worker ID 全不同（§1.4）。
3. **一发红一发绿，逐字判决行在此：**
   - `--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.43s)` ← push `36889094435`（`clean_tw_push.txt:1896`）
   - `--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.28s)` ← schedule `36940536372`（`clean_tw_sched.txt:1885`）
   - 哨兵面同时翻：push 响 `config_reload_223_test.go:303` 与 `:306`；sched **两条都零命中**。
   ⇒ 前一枚腿 §3(d) 靠"读盘上日志证渲染发生过"判死真回归，**本腿用第二次采样把它从'推理判死'升级成'采样判死'**。

**加分支（不在题面三种判据里，但同一发免费买到，必须报）：**
- 同一次 run 里 **另一枚**按机器翻面的用例（`TestPanelHostLatencyPercentilesAC2`：push SKIP ↔ sched FAIL，
  且一台冷拉起 -1.000、另一台 3496 ms）与**同码两发都红**的挂钟钉（`TestResolvePerCallBudget`，
  1.241 vs 3.346 ms/op、1902 vs 560 samples）**一起存在**。
  ⇒ "同一次 run 里多台用例按机器抽签翻面"不是本枚孤例，**是那一档用例的常态**。
- `TestResolvePerCallBudget` 两发都红**不复活"码改坏"那一支**，因为它的被测码路径在区间内**一字未动**，本腿新量了一层更强的：
  `internal/risk` 的第一方依赖闭包（`internal/observe`、`internal/plugin`、`internal/winsec`）
  在 `0589fd9c..8ae4c23e` 区间内 **各 0 枚 commit**，`internal/risk` 本身也 **0 枚**。
  （尺：`git log --oneline 0589fd9c..8ae4c23e -- <包> | wc -l` 逐包四发；导入闭包尺
  `git grep -h -oE '"github.com/CarlosShao/wisp/[^"]+"' 8ae4c23e -- internal/risk | sort -u`。）
  区间内 `go.mod`/`go.sum` 各 **1 枚** commit（`697b4fae`），diff 只是新增两枚 **indirect**
  `github.com/jchv/go-webview2` / `go-winloader`（票 33 面板宿主），**不在 Resolve 路径上**。
  ⇒ 这枚钉的红只能落在**环境侧**（机器或镜像，二者在"基线→新发"这一对上没分开，见 §5 C2）。
  ⚠ **这一条是"排除码"，不是"定位到机器"**——本腿不许把 C2 那三个假说里的任何一个写成已证。
  而 §3 判定本身**不依赖**这枚钉：它靠的是 §2.1 那一对同镜像、同 SHA、不同 VM 的 PASS↔FAIL。

### 本腿**没有**判成、也提醒编排者别顺手判成的两件事

1. **没判"push 那台 `…1157` 有问题"。** sched 那发是**第三台**机器、**第二台**给出绿/红不一致的样本，
   手里只有 push=1 红 / sched=1 绿。要说"哪台机器坏"需要**同机重复采样**（本腿禁推码、禁 rerun，买不到）。
2. **没判"可以放心划绿"。** （甲）只杀掉"新代码把产品改坏"这一支。
   `cmd/wisp/config_reload_223_test.go:301` 那次未加界的 stdout 单读**照旧是一枚真缺陷**
   （前一枚腿 §4.1 已判死反命题、并指出同文件 `:443` 同一形状、`:129-139` 早已登记过 1/25、1/5、1/15 的实测频率）。
   本腿这次是**第 3 次**为它添独立频率读数：**3 发里 1 次红**（同码 2 发：1 红 1 绿；基线 1 发：绿）。
   ⛔ 改断言/加等待属"改契约＝人工批准"，本腿不动、也不建议（§5 C6）。

---

## §4 顺带能钉死的另一格：`test-core` 这一发**同样采不到**，`slo-full` 这一发**日志回来了**

### 4.1 `test-core`（`bash scripts/portable-tests.sh --scope=core`）＝**这一发仍然零 `go test`**

尺：`awk -f strip.awk job_test-core.txt` 后逐名数判决行 ＋ 看 GUARD C 块。

| | 基线 `0589fd9c` | push `8ae4c23e` | **sched `8ae4c23e`（本腿）** |
|---|---|---|---|
| `test-core` 整作业 `=== RUN` | **1368**（`clean_all_base.txt`，本腿按 TSV 逐作业数） | **3** | **3** |
| `--- PASS:` / `--- FAIL:` | 922 / 4 | 1 / 0 | **1 / 0** |
| 那 3 枚 RUN 是什么 | （core scope 全量在跑） | 只有 preflight：`tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`（1 顶层 + 2 子测试，PASS） | **同左，逐字同形** |
| `--scope=core` 步 | 跑了 | **GUARD C 打死** | **GUARD C 打死** |
| 步开始 → 报错 | — | `16:04:19.774` → `16:04:23.083`（**3.31 s**） | `23:24:12.744` → `23:24:15.342`（**2.60 s**） |
| GUARD C 报文（逐字） | — | `portable-tests.sh: GUARD C - scope mode=core resolved to a DIFFERENT package` / `portable-tests.sh:   set than the one pinned next to it. Pinned: 25, resolved: 35.` | **两条逐字相同** |
| 差异前缀 | — | `25a26,35`，多出的 10 行全是 `> go: downloading …`（`dustin/go-humanize`、`google/uuid`、`mattn/go-isatty`、`modernc.org/{libc,mathutil,memory,sqlite}`、`golang.org/x/{crypto,sys}`、`remyoudompheng/bigfft`） | **同一批 `go: downloading` 噪声**（本腿 `clean_tc_sched.txt:343+` 逐条复点） |

⇒ **题面那句"上一发它 3 秒内死在 GUARD C，`Pinned: 25, resolved: 35`，整段 1365 个 `=== RUN` 缺失"——本腿逐字复认。**
1365 ＝ 1368 − 3，两发**同一把尺、同一个缺口**。
⇒ **"core scope 到底跑没跑 `go test`"这一格：schedule 这一发同样采不到，读数没有变化。**
⇒ 票 250（在飞的那枚 `scripts/` 腿）要修的正是这个：`scripts/portable-tests.sh` 把 `go list` 的 **stderr 并进了分母**
（`delta.md` §第 5 条已挂名 `GUARD C stderr 并进分母（冷缓存即炸）`）。
本腿不动它（硬边界 #5），只补一条：**这一发仍是冷模块缓存的形状**，两台 ubuntu VM（`…1160` / `…1166`）连续两发同炸 ⇒
**它不是抽签，是确定性缺陷**（这条比前一枚腿的"要一次温缓存的跑"更省：证据已经在盘上）。

### 4.2 `slo-full`（自托管 `wisp-selfhosted-01`）＝**这一发日志在，而且它根本不跑 `go test`**

- 可用性：push 那发**缺**（`log not found: 110459845223`）；**schedule 这发在**（34 638 B，结论 `success`）。
- 步骤名册（`awk -f strip.awk job_slo-full.txt | grep '^##\[group\]Run'`）只有四步：
  `actions/checkout@v4` → `actions/setup-go@v5` →
  `powershell … scripts/build.ps1 -Env dev` →
  `powershell … scripts/slo-check.ps1 -Subset full -SecondsPerState 6` → `actions/upload-artifact@v4`。
  `=== RUN` / `--- PASS:` / `--- FAIL:` **全 0 枚**，作业日志 280 行。
  产物 `Artifact slo-full-report.zip … Final size is 20290 bytes … Artifact ID 11202101597`。
- 排程形状值得单独记一笔：整发其它五枚作业 23:23 起跑、`test-windows` 23:35:23 收，
  **`slo-full` 却在 `2026-10-02T00:23:22Z` 才开始**（晚 48 分钟）—— 自托管 runner 排队。
  run 的 `updatedAt` `00:26:02Z` 就是它收尾的时间。
- 顺带钉死一条本腿没打算找的事实：那枚自托管 runner 的日志尾部逐字写着
  `Copying 'C:\Users\swq\.gitconfig' to …` / `D:\work\soft\Git\cmd\git.exe` / `E:\work\base\actions-runner\_work\wisp\wisp`
  ⇒ **`wisp-selfhosted-01` 就是这台开发机的物理机本身。**
  含义：`slo-full` 那一档读数**不是"CI 环境"**，它和"本机"共用一台机器；
  拿 `slo-full` 与 `test-windows` 做"CI vs 本机"对照是**拿两台不同性质的机器比**，
  前一枚腿推翻的 T3（"CI 才红"形不对）在这一格上还要再加一层。

---

## §5 判不动的地方（每条：为什么判不动 ＋ 谁能判）

**C1 —— "push 那台 VM（`…1157`）是不是特别差"判不动。**
手里只有每台 VM 各 1 枚样本：`…1150` 绿 / `…1157` 红 / `…1162` 绿。**n=1 不能定"这台机器坏"。**
谁能判：**同机重复采样**。⛔ 本腿不许 rerun、不许触发 workflow、不许推码（硬边界 #2）。
能做的只有编排者在**下一次 schedule**（同一 SHA 若还在 dev 头）自然拿到第 4、第 5 枚样本，
或允许一枚带 push 的腿专门复采样。**这条不需要 `go test`，比 §5 任何一把 `go test` 尺都便宜——本腿只是没权用它。**

**C2 —— 基线→新发那一档 29% 的挂钟涨幅（0.962 → 1.241 ms/op）判不动，而且现在有两枚候选因分不开。**
`internal/risk` 第一方闭包四包各 0 枚 commit、`go.mod` 只加两枚不在路径上的 indirect（§3 加分支）
⇒ **代码侧本腿指认不出任何一枚 commit**。但 §1.4 新量到：基线镜像 `20260922.246.2`、两发新码**都是** `20260925.250.1`
⇒ "**镜像更新**"与本枚用例的"机器抽签"在这一对上**同时变、无法分离**。
再加上手里没有 CPU 遥测、没有 Defender 扫描事件、没有 runner 宿主争用数，
本腿**三个假说（码 / 镜像 / 机器）一个都杀不掉**。
谁能判：同机 `-count=25`（前一枚腿 §5.2/§5.4）拿频率；拿得到 0922→0925 镜像 changelog 的人；
或**任何一枚"同 SHA ＋ 旧镜像"或"基线 SHA ＋ 新镜像"的 run**——本腿在盘上三发里**没有**这种交叉样本。
⚠ 本腿**不许**把它写成"就是机器慢"。
⛔ 但要说死：这一格**不影响 §3 的判语**，因为 §3 那一对（push vs sched）镜像是常量、只有 VM 在变。

**C3 —— "VM 规格/宿主负载换了没有"判不动（比 C2 更彻底）。**
三发能看到的身份**全不同**（runner 名、Worker ID），能看到的镜像声明**基线独有的小版本差**（§1.4），
`Current runner version` 与 `Microsoft Windows Server 2025 / 10.0.26100 / Datacenter` 三发**逐字相同**；
但这些行**不含 vCPU 数、不含宿主争用**。⇒ 本腿只能说"身份全不同、OS 声明全相同、镜像小版本基线独旧"。
谁能判：只有 Actions 侧拿 runner 宿主指标的人（本腿 token scopes 是 `repo, read:org, gist`，读不到）。
**本腿顺手把前一枚腿 §6 C3 那枚悬着的反例量完了**（那腿只知道它快了、没在第三发上取数）：
`TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` 三发**全 PASS**，耗时
`1.96s`（基线 `…1150`）→ `1.34s`（push `…1157`）→ `1.22s`（sched `…1162`）——
**在这台"更抖"的机器上它反而一枚比一枚快**（累计 −38%）。
⇒ 前一枚腿那句"不能写 VM 更慢、只能写更抖"**本腿复认并加强**：那 6 枚跨过 1 s tick 的用例涨、这枚掉，
**"整体慢 x%"这一形在本仓读数里根本立不住**。它为什么掉判不动（要读它的等待结构，属另一枚票的账）。

**C3b —— `cmd/wisp` 整包在两发同码上快了 4.4%（579.297s → 553.675s）、`internal/risk` 却慢了 50%（6.247s → 9.430s）。**
同一枚作业、同一台机器内、同一份码，两段的**方向相反** ⇒ 本腿**给不出"哪一段该被当作机器的代表数"**。
谁能判：只有逐段拿同机重复采样；本腿禁跑、禁 rerun。
（这条的实际后果：**任何"包级耗时涨了所以机器慢了"的一步推论都是坏尺**，包括前一枚腿 §3(c) 第 4 条那种用法。）

**C4 —— 本枚的**真实翻红频率**判不动。**
本腿能给的只有 **3 枚独立样本**（`0589fd9c` 1 枚绿、`8ae4c23e` 2 枚：1 红 1 绿）＝ **同码 1/2**。
这不等于 1/25、1/5、1/15 那串频率里的任何一档，**不能并账成"实测 1/2"**（样本量差两个数量级）。
谁能判：前一枚腿 §5.2 那把 `-count=25`（要 `go test`，本腿禁），或攒够 n≥20 的 CI 自然采样。

**C5 —— `TestResolvePerCallBudget` 该不该修判不动（且本腿明确不动）。**
§2.3 量到它在 Windows 档**基线就只剩 3.8% 余量**、同一份码在 Linux 档是 0.001 ms/op（预算的 0.07%）。
这不是"红/绿"问题，是**那枚 1 ms 预算从没按 Windows 标定过**。
本腿判不动的两点：(i) 这算不算 AC#6（票 18）本来的意图；(ii) 平台分档预算要不要改。
谁能判：票 18 的 owner ＋ 编排者。**⛔ 本腿一字未动 `thresholds.go`／golden／任何断言（硬边界 #4）**，
也**不**建议"为了变绿放宽断言"——本腿只是把余量数字报出来。

**C6 —— 该不该把 `cmd/wisp/config_reload_223_test.go:301` 那次未加界单读改成 `awaitStdout` 判不动。**
本腿这次的 PASS 只是**支持**前一枚腿 §4.1 那条"断言读早了"的链，
**不构成"可以不改"**（同码 2 发里它就红了 1 次，同文件 `:443` 是同一形状的漏网之鱼，还有一批次级同形 `:253/:274/:348/:353/:378/:421`）。
谁能判：票 223 的 owner 与编排者（改契约＝人工批准）。裁下来才轮到写码腿动那枚文件。

**C7 —— `internal/panel` 那 4 枚"基线红、push 消失"在本腿这一发是什么样子，判不动（但可以缩一半）。**
本腿**没读 `frontend/**`**（硬边界 #3），但本腿能确认一件更窄的：
`test-core` 两发（push / sched）**都是 `=== RUN=3 / PASS=1 / FAIL=0`**，
⇒ 那 4 枚在**两发新发里都没有被测到**，"消失"仍然是**分母没了**、不是修好（前一枚腿 §1.3 复认成立）。
判不动的一半：`frontend/dist` 里那批桩件是不是唯一因、以及票 250 修完 GUARD C 后它们会不会照旧红。
谁能判：面板/前端侧的腿 ＋ 票 250 落地后**同 SHA 再取一发 core scope**。

**C8 —— `slo-full` 这一发为什么"回来了"判不动。**
push 那发 `log not found`、基线与本发都在，同一名 runner（`wisp-selfhosted-01`）。
本腿能看到的只有"它这次排队排到 00:23 才开始、这次跑成功、日志传上来了"。
分不出是 runner 掉线、日志上传失败、还是作业被抢——**这些都不在 `gh` 的只读面里**。
谁能判：那台自托管 runner 的宿主日志（本腿没有、也不该去碰）。
⛔ 本腿**没有**把 push 那发的"没日志"读成红或绿（定式：一步红／一步没跑＝读数永久采不到）。

**C9 —— 本腿全部尺里唯一一把"没自跑"的：前一枚腿 §5.1–§5.4 四把 `go test` 尺。**
硬边界 #1 禁本腿跑测试/构建 ⇒ 本腿**一个字都没跑**。
本腿自跑的只有：`gh`（只读端点）／`git show|log|diff|grep|rev-parse|remote`／`sed|awk|grep|diff|comm|wc|sort|uniq|ls|find`。
凡本腿写进 §2/§3/§4 的数，**都由这些尺产生并至少用两条独立路数互证过一次**
（判决行 = per-job TSV 与整份日志 TSV 各数一遍；push/base = 裸壳 `sed` 剥一次；
GUARD C = per-job 与整份日志各一次；步级计数 = 本腿 awk 与 CI 自己的 "four numbers" 行各一次）。
谁能判：下一枚允许跑 `go test` 的腿。

---

## §6 本腿推翻题面的地方（具名，全部现量复认）

**推翻 V1 —— 题面："它的同码二次采样……`status:"completed"`、`conclusion:"failure"`"（把整发 failure 当成本枚的上下文暗示）；
题面 §0 只给了"新发 run `36889094435`（push）FAIL"，暗示两发同码的差别在"event"。**
⇒ 准确形不是"push vs schedule"，是**三台 VM**：`…1150` / `…1157` / `…1162`（§1.4）。
event 与结果**不相关**——基线（schedule）绿、push（push）红、本发（schedule）绿；
两发 schedule 编译的是**不同的** SHA。**"换 event"不是自变量，"换机器"才是。**
⇒ 对下一枚腿的实际后果：**再拿 event 做归因轴会白跑一发。**

**推翻 V2 —— 题面："基线 run `36789659174`（`headSha 0589fd9c`）PASS，新发 run `36889094435`（`headSha 8ae4c23e`）FAIL"——
题面把基线描述成一发与"新发"对偶的普通样本，暗示手里是"1 枚旧码绿 + 1 枚新码红"。**
⇒ 实际盘上基线那发的 `test-core` **还跑了 core scope**（1368 枚 `=== RUN`），
所以 `TestResolvePerCallBudget` 在基线**有两个样本**（windows PASS 1.81s ＋ core PASS 1.97s），
新发**两个都没有**。⇒ 三发的可用样本数是 `2 / 1 / 1`，不是题面读起来的 `1 / 1`。
（谁造成的：§4.1 的 GUARD C，票 250 那一格。）

**推翻 V3 —— 题面："另一枚同族的挂钟用例 `TestResolvePerCallBudget`（`internal/risk/pathresolver_budget_norace_test.go:12/:36`…）
请一并逐名读，它是最硬的旁证：它只可能因机器慢而红。"**
⇒ **"最硬的旁证"这个定位要改。它不是旁证，它自己就是本发的主角，而且它**不支持"码改坏死"这一步**——
它在两发同码上**都红**。真正把（甲）钉死的是 §2.1 那枚 ticket223 的 PASS/FAIL 翻转。
它给的**是另一件题面没要的东西**：它**自报了机器速度**（1902 vs 560 samples、1.241 vs 3.346 ms/op，同码），
并且它证明这枚钉**在 Windows 档本来就没有余量**（基线 0.962/1.000）。
⇒ 下一枚腿别拿它的"两发都红"去复活真回归那一支（§3 加分支已把它的被测码路径证到 0 改动）。

**推翻 V4 —— 题面："本仓已知：有的作业会 `log not found`——上一发 `slo-full` 就是这样，那时只有 runner 掉线注解"。**
⇒ schedule 这发**六枚作业日志全在、`log not found` 零枚**（§1.1、§1.2）。
并且本腿从**回来的那份 `slo-full` 日志尾部**看到它的 host 路径是
`C:\Users\swq\.gitconfig` / `D:\work\soft\Git\cmd\git.exe` / `E:\work\base\actions-runner\…`
⇒ **`wisp-selfhosted-01` 就是这台开发机**（§4.2）。
后果：拿 `slo-full` 当"CI 那一侧"的对照样本是**类别错**——它和"本机"同机。这条题面与前一枚腿都没有。

**推翻 V5 —— 题面："如果同一发里 `slo-full`／`test-core` 那两段的可用性也能读到，就把'`core` scope 到底跑没跑 `go test`'再取一次样……"**
⇒ 取到了，读数是**"仍然没跑"**，而且本腿把它**从抽签改判成确定性缺陷**：
两台不同的 ubuntu VM（`…1160` push / `…1166` sched）在**同一枚 SHA** 上
都在 GUARD C 上 2.6–3.3 秒打死、`Pinned: 25, resolved: 35` 逐字相同（§4.1）。
⇒ 它**不是**"冷缓存偶然抽到"，票 250 不用等同机复采样就能动手。

**推翻 V6 —— "基线 vs 新发"这一对不是干净的代码对照**（题面 §背景与前一枚腿都按"只有码变了"来叙述这一对）。
本腿从作业日志头部现量：基线 `test-windows` 的 Runner Image 版本 **`20260922.246.2`**、
Hosted Compute Agent `20260828.587`；两发新码**都是** `20260925.250.1` / `20260901.588`
（尺：`grep -aoE "(Version: 2026[0-9.]+|Current runner version: '[^']*'|Worker ID: \{[^}]*\}|Azure Region: [a-z0-9]+")"` 逐发）。
⇒ **"基线绿→新发红"这一对上，代码与 CI 镜像同批变了。**
所以：拿它俩做任何"归因到区间内某枚 commit"的推论都**先天带一枚未控变量**。
本腿的 §3 判定**因此刻意只用同镜像的那一对**（push ↔ sched，两边都 `20260925.250.1`）。
⇒ 对下一枚腿的实际后果：**别再用 `0589fd9c` 那发当"同镜像基线"用**，那一格在盘上不存在。

**推翻 V6 —— 题面 §2 的括号："并注明各自的 runner 作业号（是不是同一台机器？若同机，这一发的判别力要打折）"。**
⇒ **不同机**（V1/§1.4），所以**不打折**。但要说死打折的反面：本腿也**因此买不到**"同机重采样"那格
（那才是能定"哪台机器坏"的尺）——见 §5 C1。

---

## §7 一句话交付

（甲）。**同一枚 SHA `8ae4c23e`、两台不同 VM（`1000001157` 红 / `1000001162` 绿）**：
`TestTicket223HandEditedFsLooseningCostsAnL2Card` 一发 `--- FAIL … (2.43s)`、一发 `--- PASS … (2.28s)`，
哨兵 `config_reload_223_test.go:303`/`:306` 同步从"各响一次"变成"零命中"
⇒ **"这版码就是红"这一支死，"机器/负载抽签"这一支活**；前一枚腿 `ci223-1` §3(d) 那条靠读盘闭合的判死，
本腿**升级为同码二次采样判死**。
同发免费买到的旁证：`TestResolvePerCallBudget` 自报 1.241 vs 3.346 ms/op（1902 vs 560 samples，同码同二进制）
与 `TestPanelHostLatencyPercentilesAC2` SKIP↔FAIL（-1.000 vs 3496.853 ms）——
**同一只手在这一次 run 里同时翻了至少三枚按挂钟吃饭的用例，而它们被测的代码在区间内一字未动。**
---

---

## 附 A：本腿**没有**入库的中间件与它们的复算尺（临时件只建不删 ⇒ 全在盘上，只是不进 commit）

本 commit 只收：`read.md` ＋ 七把 awk 尺 ＋ 六枚 `job_*.txt`/`job_*.err`（§1 可用性证据本体 ＋ 决定性那发的原始日志）
＋ `clean_tw_sched.txt`/`clean_tc_sched.txt` ＋ 逐名小件（`up_*`／`ws_*`／`sched_*_cmdwisp.txt`）。
下列件**体积大且可一步再生**，留在盘上不进库：

| 未入库件 | 一步复算尺 |
|---|---|
| `logs_sched_all.txt`（1 073 604 B） | `gh run view 36940536372 --log` ＝ 六枚 `job_*.txt` 的并集，内容不增 |
| `clean_all_sched.txt` / `clean_all_push.txt` / `clean_all_base.txt` | `awk -f strip.awk logs_sched_all.txt` ／`awk -f strip.awk ../ci-delta-1/logs_new_all.txt`／`…logs_base_all.txt` |
| `clean_tw_push.txt` / `clean_tw_base.txt` | `sed -E 's/\r$//; s/^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9:.+-]+Z ?//' ../ci-delta-1/raw_tw_new.txt`（base 换 `_base`）；**裸壳件，不能用 `strip.awk`**（§1.3 第 2 条那枚坑） |
| `tsv_tw_push.txt` / `tsv_tw_base.txt` | `sed 's/^/test-windows\tSTEP\t/' ../ci-delta-1/raw_tw_new.txt` |
| `clean_slofull_sched.txt` | `awk -f strip.awk job_slo-full.txt` |

原始件 `ci-delta-1/raw_tw_{base,new}.txt`、`raw_testcore_new.txt`、`logs_{base,new}_all.txt` **都是前一枚腿已入库的件，本腿一字未动、只读。**

## 附 B：本腿全部结论所依赖的尺（逐条，可单发复算）

| # | 尺（全部只读） | 用在哪一节 |
|---|---|---|
| R1 | `git rev-parse --abbrev-ref HEAD`／`git rev-parse HEAD`／`git log -1 --format='%H %ad' --date=iso`／`date`／`gh auth status`／`gh repo view --json defaultBranchRef,name,owner`／`git remote -v` | §0.1 |
| R2 | `gh run view <r> --json databaseId,workflowName,name,event,headBranch,headSha,status,conclusion,createdAt,updatedAt,displayTitle,attempt` | §0.2 |
| R3 | `gh run view <r> --json jobs --jq '.jobs[] \| [...] '`（**`id` 不是合法字段，用 `databaseId`**） | §1.1 |
| R4 | `gh api repos/CarlosShao/wisp/actions/runs/<r>/jobs --jq '.jobs[] \| .name + " \| runner_name=" + (.runner_name // "NULL") + " \| labels=" + ((.labels//[]) \| join(","))'`（**`.runner.name` 六枚全 null，必须用下划线形**） | §1.1／§1.4 |
| R5 | `gh run view 36940536372 --log --job <id> > job_<name>.txt 2> job_<name>.err`，逐枚看 `rc` 与 `wc -c` 两侧 | §1.1 |
| R6 | `grep -aoE "(Worker ID: \{[^}]*\}\|Azure Region: [a-z0-9]+\|Current runner version: '[^']*'\|Version: 2026[0-9.]+\|Microsoft Windows Server [0-9]+)" <件> \| sort -u` | §1.4 |
| R7 | `awk -f strip.awk <件>` ＋ `grep -anE '^ *--- (PASS\|FAIL\|SKIP): (TestTicket223HandEditedFsLooseningCostsAnL2Card\|TestResolvePerCallBudget)'` | §2.1／§2.2 |
| R8 | `grep -aoE "config_reload_223_test\.go:[0-9]+" <件> \| sort \| uniq -c` | §2.1 |
| R9 | `awk -f stepstats.awk <件>`（按 `##[group]Run` 归属出步级 RUN/PASS/FAIL/SKIP）＋ `grep -aE '=== RUN=[0-9]+.*PASS=[0-9]+'`（与 CI 自己那行 "four numbers" 并排） | §2.4 |
| R10 | `awk -f cmduisp_fail.awk <件> \| sort`（cmd/wisp 步内 FAIL+SKIP 逐名）／`awk -f winscope_fail.awk`（windows scope 步内）＋ `diff` | §2.4／§2.5 |
| R11 | `awk -f beforeours.awk <件>` ＋ `diff`（三发跑序上游逐名逐序） | §2.4 |
| R12 | `git log --oneline 0589fd9c..8ae4c23e -- <包> \| wc -l` 逐包四发 ＋ `git grep -h -oE '"github.com/CarlosShao/wisp/[^"]+"' 8ae4c23e -- internal/risk \| sort -u` ＋ `git diff 0589fd9c 8ae4c23e -- go.mod` | §3 |
| R13 | `git show 8ae4c23e:internal/risk/pathresolver_budget_norace_test.go \| grep -nE "resolveBudget\|ns/op\|Fatalf\|Logf\|func Test"`；同法读 `cmd/wisp/panel_host_windows_test.go:785-815` | §2.2／§2.5 |
| R14 | `awk -f strip.awk job_test-core.txt` ＋ 步级 `=== RUN`/`--- PASS` 计数 ＋ `grep -an "GUARD C" -B12 -A3` ＋ 整份日志按 `$1=="test-core"` 独立数一遍 | §4.1 |
| R15 | `awk -f strip.awk job_slo-full.txt \| grep -aE '^##\[group\]Run\|^##\[error\]\|^##\[warning\]'` ＋ `tail -25` 读 host 路径 | §4.2 |

---

## §7b 交付口径后半：入账姿势（与上方 §7 同节，因插了附 A/B 才落在末尾）

⛔ 本枚**不构成**"可以划绿"的证据：`cmd/wisp/config_reload_223_test.go:301` 那次未加界 stdout 单读仍是真缺陷，
本腿按硬边界 #4/#5 一字未动、不建议自动修（§5 C6）。
入账姿势：**"未定性"升级为"已排除真回归 ＋ 已量化到同码 1/2 采样翻面"**，
补差集里那枚的名建议写作 `CI 抖动／同码二采样判死真回归`，出处＝本文件 §2.1 与 §3。


