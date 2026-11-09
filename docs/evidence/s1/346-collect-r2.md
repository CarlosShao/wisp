# 346-collect-r2 收件件：票 346 七格逐格复跑（本波被授权翻框腿）

- 腿＝`collect-r2`（编排者收件，本波唯一被授权翻框）；2026-11-09 11:5x +08；起手 HEAD `5071d8825`。
- 起手三核：票面七格现读全文（⛔ 信转述）；交件件三枚现读（`346-a1-census.md`／`346-c1-census.md`／`346-r1-impl.md`）；`git status --porcelain` 无 346 名下写面。
- 门禁读数跑在 `git clone --shared` 导出树（`/d/tmp/wisp/20261109/346collect/export`，HEAD 同串）；原始 CI 日志在盘＝`D:\tmp\wisp\20261105\346r1\run-39101060753.log`（2358566B）⇄`run-39124122121.log`（2795862B）。
- 本腿 0 枚 `go test`、0 枚 `.go` 字节、⛔ push。

## 七格逐格读数与裁定

### AC#0 起手现跑＝翻（读数逐枚自取）

| 行 | 本腿现跑读数 | 与票面/交件对 |
|---|---|---|
| blob 同一性 | 两枚仪器在 `a42b35af`⇄`a5080612`⇄HEAD 三点逐枚同串：`dfaddac8`／`0b0e5173`；`winsec_other.go` 三锚 `e039da10`、`winsec_windows.go` 三锚 `63714a83`；`winsec.go` 两锚同 `67783ae0`、HEAD `cbafe712`＝变过（票 342r1/352r1/353r1/348r2/400r1 那 30 笔，⛔ 归 346 名下） | 同票面第 17 行＋a1 补边界 |
| 翻脸事实 | run 39101060753 两枚顶格 `--- PASS (0.46s)/(0.08s)`；run 39124122121 两枚顶格 `--- FAIL (0.63s)/(0.11s)`；`adec1db1` 无 openleaf（`cat-file -e` rc=128） | 同票面第 18 行，含时长逐枚同 |
| 红句①②③ | 三句逐字在红发 log 各命中（①`answered 0 … landed 6 flips … never reached`；②`refused=1998, sealed-own=0 … proves nothing`；③`refused-at-act=0 … pinned identity gate`） | 同票面第 19–20 行（③是 `t.Logf` 非 `t.Errorf`，按 c1 更正口径） |
| 分母 | 绿发顶格 FAIL＝14 行（core 9＋win 5）／12 唯一名；红发＝13 行（core 9＋win 4）／11 唯一名；两发并集 14 唯一名 | ★票面第 21 行字面读数（core9+win3）⛔ 复现＝c1 已更正、本腿独立复现更正读数 |
| 其余在册 | 归因尺跑两发 log＝`unique red names=14  UNREGISTERED=0`（正控 `TestC26…` 在册命中） | 票面"⛔ 一枚归新账"复现 |
| lint 步 | 两发 gofumpt 步各点名 **20 枚** `.go`、全在 `.scratch/wisp/probes/**`（非 probes 命中 0）、两发作差＝空 | 同票面第 23 行 |
| 同族尺（现跑） | `git grep -n "RED instrument" -- internal/winsec/`＝**31 行命中／12 个文件**（比 c1 的 26/11 又长＝票 354 族 `leftstanding_354_windows_test.go` 新进，名册逐趟腐烂坐实）；本票两枚仍在 `openleaf_322_other_test.go:206`／`posix_seal_window_120_other_test.go:225/:233` | 票面"逐趟必须重跑"复现；同族普查在册（a1 12/8 ⇒ c1 26/11 ⇒ 本腿 31/12） |

**verdict=翻**（判据本体六行＋同族尺全复现；唯一偏差＝分母字面读数，票面第 5 行明写"行号/读数只当读数⛔ 当判据"，更正读数已在 c1 件并在本腿独立复现）。

### AC#1 三形代价表＝翻（表已出、抽点复现；选形仍⛔ 在非实现者外）

表＝`346-c1-census.md` 第二节（甲/乙/丙并排在 `:137-143`）。本腿四处承重落点现跑复现：
① `--scope=winsec` 档在仓＝`scripts/portable-tests.sh:431`（`winsec_scope=./internal/winsec/...`）＋票 251 AC#1 注释 `:469`；② CI ⛔ 叫它＝`grep "scope=winsec" .github/workflows/ci.yml` rc=1、`winsec-tests.sh` 只挂 windows 档（`ci.yml:872`）；③ POSIX 接缝在＝`winsec_other.go:50 var applyDescriptor = applyDescriptorPOSIX`；④ 三枚分母常量在＝`posix_seal_window_120_other_test.go:100/:101`（4/500）＋`openleaf_322_other_test.go:35`（3000）；⑤ 整包塞同一枚 `go test -v -count=1`＝`tools/d22scan/runtests.sh:211`（⛔ 带 `-p`）。

**verdict=翻**（判据要求＝表未出⛔ 开工⛔ 选形；表已出且承重现量抽点全对上。★选形本身仍归编排者裁，甲=`.github/**` 须人工批准，本翻不含选形）。

### AC#2 去牙反形＝⛔ 翻，需 Linux

两枚仪器首行 `//go:build !windows`（a1/r1 各自现读、本腿不再抄），摘弃权守卫那一发必须在仓外导出树 **Linux 档**跑；本机 Windows 结构性不编译＝既⛔ 红也⛔ 假绿可量。且该格以"已选形"为前提（选形未裁）。**具名差什么：一台 Linux（容器或 CI `test-core`）＋选形裁定。**

### AC#3 归因工具＝翻（RED⇄GREEN 两半本腿亲手复跑，导出树）

- 工具在盘＝`.scratch/wisp/probes/346/a1/ci-red-attribution.sh`＋`selftest.sh`（r1 提交 `1f298db54` 落盘，17 枚文件）。
- **RED-before**：`TOOL=/d/tmp/wisp/20261105/346r1/baseline.sh bash selftest.sh` ⇒ **rc=1、FAIL (8 of 16 checks)**（基线 sha256 `db1f1694` 与 r1 所记同串）。
- **GREEN-after**：`bash selftest.sh` ⇒ **PASS (34 checks)、rc=0**。
- **两枚真 run 现跑**（`--log` 两枚原始 log，导出树仓根）：`logs=2  顶格 FAIL (job x name)=16`、`POSITIVE CONTROL ok`、本票两枚逐字 `=> REGISTERED`（`具名文件在册票=322,346,369`／`120,322,342,346,369`）、`unique red names=14  UNREGISTERED=0`、rc=0——与 r1 件第二节读数逐字同。

**verdict=翻**。

### AC#4 双平台成对＋三发＝⛔ 翻，需 Linux 容器＋CI

同 AC#2 的结构性理由（Windows 档⛔ 编译两枚仪器）；且"本机 Windows 档"那一轴对两枚不存在，成对读数只能 CI `test-core` ⇄ Linux 容器。台账 `A1407` 已具名"346 暂缓，其绿读数须 CI，本机不可闭环"。**具名差什么：Linux 容器（⛔ bind 宿主 Windows 卷，按 a1 (iv) 配方）＋同台面 ≥3 发独立复跑＋红⇄绿⇄弃权散布名册。**

### AC#5 禁区＝翻（全部可验项现跑为净）

- r1 提交 `1f298db54` 触碰面＝`probes/346/a1/**` 12 枚新增＋`docs/evidence/s1/346-r1-impl.md` 1 枚＝**零禁区文件**。
- 工作树禁区脏检＝`git status --porcelain -- internal/winsec/ internal/observe/thresholds.go go.mod go.sum .github/ scripts/ tools/d22scan/allowlist.txt docs/PLAN.md .scratch/wisp/issues/`＝**0 枚**。
- 票 322/120 既有钉＝两枚仪器 HEAD blob 仍 `dfaddac8`／`0b0e5173`（=三锚同串那对）；`winsec.go` HEAD 变更全部来自他票提交（342r1/352r1/353r1/348r2/400r1），346 名下腿 0 枚 `.go` 字节。

**verdict=翻**。

### AC#6 纪律格＝翻（可验项全过）

- 票面触碰史＝`git log -- <票346文件>` 只有 4 笔：立票 `05d37577e`／a1 只追加 `268bbb3c5`／编排者收件 `220056a0d`／collect-r2 状态盘点 `84f4a1d0a`；**实现腿 r1 零触碰票面**（其 STOP-note ④ 具名交回，⛔ 违令）。
- r1 提交 17 枚文件逐枚非 0 字节、零字节 0、`.out` 0、U+00A7 命中 0、HEAD blob **CR 字节 0**（全 LF）。
- 提交内容恰为其自报写面＝显式 pathspec 的可观测结果；⛔ push 与中文⛔ PowerShell 内联无反证（其证据件为纯 ASCII 命令行产出的读数件）。

**verdict=翻**（写面许可中"票尾一条 bullet"r1 未用＝许可非义务，其具名交回在案）。

## 结论

七格 **5 翻（AC#0/#1/#3/#5/#6）＋ 2 照实报需环境（AC#2 需 Linux、AC#4 需 Linux 容器＋CI 双平台）**。票⛔ 改 `-done`（两格开着）。AC#2/#4 的下一步归属＝编排者（需选形裁定＋Linux 台面，另派或并入票 342 落地波）。
