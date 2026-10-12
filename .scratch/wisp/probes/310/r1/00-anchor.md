# `00-anchor.md` — 腿 `310-r1`（实现者，票 310 `AC#1`＝甲-G）起手锚

生成时刻 `2026-10-12 11:5x +08`（现量 `date` 见 §2）。⛔ 本腿宣布任何 `AC` 成立（判语归下一程非实现者腿）。

---

## 1. 派单要求逐字抄回的三段（抄全才算读了；正文＝`.scratch/wisp/probes/orch/20261012/txt/310r1-dispatch.txt`）

### §0 第 3～4 步（逐字）

> 3. **起手闸门逐字跑**（正典＝`.scratch/wisp/probes/orch/gate-roster-304.md` §1：**逐字抄进 `00-anchor.md` 再跑**）：四条 `tasklist //FI "IMAGENAME eq <名>"`（`wisp.exe`／`balldebug.exe`／`mockllm.exe`／`wisp.test.exe`）＋那条家族宽松尺（单引号包住整条 PS 命令）＋**新加那一行 `%TEMP%` 下 `Test*` 枚数**。⚠ 命中 0 也要落 stderr 与 rc；★**家族尺的正控必须同行落**（正典 §1 那句：同管线跑一枚必非零的名字，证明这把尺⛔ 恒 0）。
> 4. **第 1 笔 commit 立刻做**：只交 `00-anchor.md`，形式＝`git add -- <显式路径>` 之后 `git commit -F <你自家 logs 里的消息件> -- <同样的显式路径>`。**⛔ 裸 `git commit`**（提交的是整个 index）；**add 带 pathspec、commit 也带 pathspec**；⛔ `git add -A`／`git add .`／`git add <目录>` 一把扫。

### §2 第 5 条（逐字）

> 5. ★**删目录＝不可逆**：本票任何一步⛔ 批量删、⛔ 清 `%TEMP%` 顶层、⛔ 动那**在册 55 枚存量**（含 `A857` 里我裁"归属不明"的那枚 `TestProbeStepByStep3843095677`，⛔ 动、永久）；你自己的专用根**收尾只写路径**，要删也只能按你自己那份逐枚名册、且件里具名"这一枚是我这一发建的"。

### §3 里带 ★ 的那几条（逐字）

>   - ★两发正控缺一⛔ 算凭据；⚠ **"0 枚"⛔ 是读数**——要么同批发证明坏世界被造出来，要么当场标〔作废〕并写非零 rc（编排者 10-11 因这个形态自认作废过两发、10-12 又因台件把自家表头冲掉作废过一发）。
> - `30-gates.md`＋逐件读数（★每把⛔ 0 字节、每件自落一行 `rc=N`；若你引用一枚同伴件说"它是空的、不是没交"，先 `wc -c` 再写）：`sh scripts/d22scan.sh` rc=0；`sh scripts/check-path-length-budget.sh --with-self-test` rc=0；**格式名册两把并排**（工作树 ⇄ `git show HEAD:` blob，各写射程目录＋枚数，先证取出的东西非空再信它的 0）；`bash scripts/portable-tests.sh --scope=census` 的 **totals 行逐字未变**（基线 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，跑前跑后各一发作差、⛔ 拿旧数当基线）；★**你这套改动本身要跑一遍那四枚入口件**（至少 `--scope=census` 与一条窄档），具名报"跑法变了没、颜色变了没"。⛔ 整包 `go test ./...`（会把真窗/音频用例卷进来）。

---

## 2. 锚点与在飞态（派单 §0 第 1～2 步的现量）

```
git rev-parse --short HEAD        = 7e6c1579
git rev-parse --abbrev-ref HEAD   = dev
git status --porcelain -- internal cmd tools scripts  行数 = 0
git -c core.quotePath=false diff --cached --name-only | grep -c ''  = 0   （grep 退码 1＝命中 0 枚，不是尺坏）
```

- 暂存区起手＝**0 枚**（⛔ 别家的在飞态，本腿无需归位动作；若后续发现别家路径只写"我看到 N 枚、我没动"）。
- ⚠ 机主侧未提交改动在 `.gitignore`／`go.mod`／`design/**` 那几枚（派单点名⛔ commit ⛔ revert ⛔ 修改）；本腿只用显式 pathspec 提自己那三件套＋`scripts/**`，不碰它们一根指头。
- `date` 现量＝件首那一行；件＝`logs/00-gate-310r1.txt`（18 行）＋`logs/00-gate-310r1.err`（**0 字节＝起手闸门 stderr 空**，`wc -c` 于 §4）。

## 3. 起手闸门（正典 §1 逐字抄进来，再跑）

正典四条＋那条家族宽松尺＋票 310 新加那一行，逐字：

```
tasklist //FI "IMAGENAME eq wisp.exe"          # 命中行数＝枚数
tasklist //FI "IMAGENAME eq balldebug.exe"     # 命中行数＝枚数
tasklist //FI "IMAGENAME eq mockllm.exe"       # 命中行数＝枚数（票 304 起必数：测试起的"假大脑"会挂着一枚回环套接字不退）
tasklist //FI "IMAGENAME eq wisp.test.exe"     # 命中行数＝枚数（`go test -c` 的产物，同理会留）
```

```
powershell -NoProfile -Command 'Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*.test.exe" } | Select-Object ProcessId,ParentProcessId,Name,CommandLine | Format-List'
```

```
powershell -NoProfile -Command '(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like "Test*" }).Count'   # 枚数＝%TEMP% 下测试临时目录残留数
```

本腿现量（件＝`logs/00-gate-310r1.txt`，stderr 并排落 `logs/00-gate-310r1.err`）：

| 尺 | 读数 | rc |
|---|---|---|
| `tasklist //FI "IMAGENAME eq wisp.exe"` | 命中 **0** 行 | `rc=0` |
| `tasklist //FI "IMAGENAME eq balldebug.exe"` | 命中 **0** 行 | `rc=0` |
| `tasklist //FI "IMAGENAME eq mockllm.exe"` | 命中 **0** 行 | `rc=0` |
| `tasklist //FI "IMAGENAME eq wisp.test.exe"` | 命中 **0** 行 | `rc=0` |
| 家族宽松尺（正典形，逐字 Format-List） | 命中 **0** 枚（`family_count=0`，计数在 PS Inside 里算好再打印＝正典 §1 指定的正形） | `family_verbatim_rc=0`／`family_count_form_rc=0` |
| ★**同管线正控**（证明这把尺⛔ 恒 0） | `posctl_svchost_count=`**85**（`svchost.exe` 同一形管线） | `posctl_rc=0` |
| `%TEMP%` 下 `Test*` 枚数（票 310 `AC#4` 那一行） | **55** ＝编排者 10-12 复跑 (b) 两发之后的在册存量；`%TEMP%` 目录总数＝**2783** | `temp_rc=0` |
| 起手名册逐枚名单 | `logs/00-temp-roster-start.txt`＝**55 行** | `roster_rc=0` |

- ★家族尺的正控与 0 读数**同行落**（上面那两行同批、同一管线）；⛔ 拿"空输出"当"0 枚"。
- ⚠ **数到 ≠ 归我**：起手台面**全 0**（wisp 系进程零枚），所以本腿没有需要答"归属三问"的非零计数；`Test*` 那 55 枚是**在册存量**，派单 §2 第 5 条已写死⛔ 动，本腿一枚都不碰、一枚都不删。

## 4. stderr 与 0 字节自证

```
wc -c logs/00-gate-310r1.err   →  0     （empty above＝0 字节＝起手闸门无 stderr）
wc -l logs/00-gate-310r1.txt   →  18
wc -l logs/00-temp-roster-start.txt → 55
```

## 5. 起手即顶回（派单／票面与盘上不符的两处，具名报回，⛔ 照抄）

1. **名册第四枚与盘上不符**：`scripts/slo-check.ps1` **零枚 `go test` 调用点**——它起的是 `scripts/build.ps1`（`go build`，件内 `:296-314`）与 `wisp.exe`；`go.exe` 只出现在 `:361` 那枚"争用侦测的进程名清单"里。⇒ 本腿⛔ 在那一枚里造一步（造了就是无的放矢的语义改动）。逐枚行号名册见 `10-impl.md` §1。
2. **第五枚不在名册上**：`scripts/d22scan.sh`（经 `tools/d22scan/runtests.sh:75` 起 `go test`）是**本票四枚之外**另一枚会起 `go test` 的入口件。派单写"唯一落点＝那四枚"⇒ 本腿**不动它**，把这一支漏面具名登记成欠账（`Q-310-r1-2`，见 `90-final.md`）。
