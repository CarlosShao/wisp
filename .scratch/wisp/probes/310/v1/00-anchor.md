# `310-v1` 起手锚＋名册三问（非实现者验收腿，票 310 `AC#1`/`AC#2`/`AC#3`）

派单原文＝`.scratch/wisp/probes/310v1-dispatch.txt`（已逐字读）。本程是**验收腿**：⛔ 产码、⛔ 翻框、⛔ 动 `Status:`、⛔ 改票名、⛔ push、⛔ 删盘上任何东西。

## 1. 起手锚（本刻现读，逐字）

| 尺（逐字可重跑） | 我的读数 |
|---|---|
| `git rev-parse --short HEAD` | `618d984c` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- internal` 行数 | `0` |
| `git status --porcelain -- cmd` 行数 | `0` |
| `git status --porcelain -- tools` 行数 | `0` |
| `git status --porcelain -- scripts` 行数 | `0` |
| `git status --porcelain -- docs` 行数 | `0` |
| `git status --porcelain -- .scratch/wisp/issues` 行数 | `0` |
| `git status --porcelain`（全仓，只为数不归属）行数 | `842` |
| `git diff --cached --name-only`（暂存区枚数） | `0` |

**大白话**：我开工这一刻，HEAD 在 `618d984c`、分支 `dev`，派单点名的六枚射程目录**干净**（没人在里面留未提交改动），暂存区里**没有**别家的在飞改动（`0` 枚）——所以不需要 `git restore --staged`，我也⛔ 需要清空任何人的东西。
全仓另有 `842` 行未提交（机主/别家的 `.gitignore`、`go.mod`、`design/**`、`.scratch/**` 那一批）——**我看见了 842 行，我一枚没动、一枚没 add、一枚没 commit**。读数件＝`.scratch/wisp/probes/310/v1/logs/00-porcelain-all.txt`。

实现腿 `310-r1` 的落点（票面 `:52` 那句 + 我下面会自己复跑）＝三枚入口件 `scripts/portable-tests.sh`／`scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`。

## 2. 起手名册（正典＝`.scratch/wisp/probes/orch/gate-roster-304.md` §1，先读了那一节才跑）

四条精确名各一发＋家族尺＋`%TEMP%` 下 `Test*` 枚数；**每条 0 读数同行带正控（`svchost.exe`）与 stderr 字节数**；PowerShell **一律写成 `.ps1` 用 `-File` 跑**（⛔ 内联 `-Command` 放 `$_`／反引号／`-join`）。

### 2.1 四条精确名（件＝`logs/00-tl-summary.txt`，台件＝`gate-310v1.sh`，`sh_rc=0`）

```
exact_name=wisp.exe rc=0 hit_lines=0 stdout_bytes=64 stderr_bytes=0 positive_control_svchost=83
exact_name=balldebug.exe rc=0 hit_lines=0 stdout_bytes=64 stderr_bytes=0 positive_control_svchost=83
exact_name=mockllm.exe rc=0 hit_lines=0 stdout_bytes=64 stderr_bytes=0 positive_control_svchost=83
exact_name=wisp.test.exe rc=0 hit_lines=0 stdout_bytes=64 stderr_bytes=0 positive_control_svchost=83
exact_name=svchost.exe rc=0 hit_lines=83 stdout_bytes=6632 stderr_bytes=0 positive_control_svchost=83
```

`empty above＝0 枚` 那四行逐字已在件里（每件 `stdout_bytes=64`＝tasklist 那句 "没有运行的任务"，⛔ 是空输出）。

### 2.2 家族尺（件＝`logs/00-gate-start.txt` 第 7～8 行，台件＝`gate-310v1.ps1`，`ps_rc=0`，那发 stderr `0` 字节＝`logs/00-gate-start.err`）

```
canonical_count=0
family_positive_control_svchost=83 family_positive_control_all_exe=343
```

计数在 **PS Inside** 算好再打印（正典 §1 那条"假 0"事故的改正形）。⚠ **我这一发第一次踩了同一枚雷**：我最初把那四条 `tasklist` 也放进 `.ps1` 里跑，PowerShell 不认 bash 的 `//FI` 拼法，正控 `svchost.exe` 一起报 **0**＝一枚假 0 台面；发现后我把那四条挪回 bash 重跑（§2.1 那五行才是算数的读数），`.ps1` 里只留 `Get-CimInstance` 那把。**记我一笔**：正控同行这条规矩第一次跑就救了我，具名＝`Q-310-v1-0`（已在本件结清，无需欠账）。

### 2.3 `%TEMP%` 下 `Test*` 枚数（同一把尺，件＝`logs/00-gate-start.txt` 第 11 行）

```
temp_root=C:\Users\swq\AppData\Local\Temp\ test_count=58 positive_control_all_temp_dirs=2807
```

正典 §1 那条现量序列（…→`310-r1` 交件 **58**→编排者复跑 winsec 两发之后仍是 **58**）⇄ **我这一刻 58** ⇒ **没有再长**（我这一程到本刻为止⛔ 起过任何 `go test`；下面若起，会在各格件里续数）。

### 2.4 名册三问（正典 §1 末："数到 ≠ 归我"）

1. **谁起的**：本刻四条精确名＋家族尺全 0，**没有一枚活着的进程要我答这一问**（0 台面）。`Test*` 那 58 枚是**在册存量**：正典 §1 那条序列与台账 `A858` 已逐枚归过来路（`310-a1` 3 枚＋编排者自己 11:04/11:05 两枚＋`310-r1` 3 枚＋更早的存量），**本腿一枚不起、一枚不认领新增**。
2. **父 PID 还在不在**：本发无非零计数可用作旁证；且按正典那条例，这一问⛔ 当归属用（Windows 复用 PID），本腿⛔ 靠它判任何东西。
3. **本腿认下哪几枚**：**0 枚进程、0 枚目录**。我唯一"认下"的是**别家的**：58 枚在册目录（⛔ 动、含 `A857` 裁"归属不明、⛔ 动、永久"的 `TestProbeStepByStep3843095677`）＋两枚专用根的**全部残渣**（下面 §2.5）。

### 2.5 两枚专用根（本腿只数、⛔ 动；件＝`logs/00-gate-start.txt` 第 71～84 行逐字）

```
root=C:\Users\swq\AppData\Local\wisp-test-gotmp exists=1 child_dirs=6
    child=portable-20261012-115401-2013 entries=1 created=2026-10-12 11:54:01
    （另 5 枚 entries=0）
root=D:\tmp\wisp-test-gotmp exists=1 child_dirs=6
    child=winsec-20261012-122119-659 entries=4 created=2026-10-12 12:21:19
    （另 5 枚，其中 probe-114728-128 / shapec-115000 各 entries=1）
```

`…122119-659 entries=4` ⇄ `A858` §5 那把 `Get-ChildItem … | Measure-Object`＝**4** 逐字同数；默认 `%LOCALAPPDATA%` 兄弟根 `winsec-20261012-122147-814 entries=0` ⇄ `A858` §5 的 **0** 逐字同数。**我一枚没删、没碰、没改名**。

## 3. 与派单转述的冲突（要求"以盘上为准并具名报回"）

- 派单 §1 说被测件"落点＝三枚 `scripts/*.sh` 各一段 `GOTMPDIR` 块，`git diff --numstat`＝69/21/20 加、0 删" ⇒ 这一条我**⛔ 采信转述**，`10-ac1-verdicts.md` 里我自己现跑 `git log`/`git diff` 对数。
- 派单 §1 说 `%TEMP%` 在册 **58** 枚 ⇄ 盘上现量 **58** ⇒ 一致。
- 派单 §0 未写锚号（按规矩＝我这一刻的读数 `618d984c`）；票面/台账里被派单转述的三处（`A858` 的 `2236627062` 从未存在、`38631165` 是被 `cut` 截短的名、第五枚入口件 `tools/d22scan/runtests.sh:75`）**我逐一现跑复核**，判语在 `10-ac1-verdicts.md`／`20-ac2-paired.md`。
