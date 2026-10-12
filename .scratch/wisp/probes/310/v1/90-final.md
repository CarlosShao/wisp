# `90-final.md` — `310-v1` 结论一页（非实现者验收腿，票 310 `AC#1`／`AC#2`／`AC#3`）

**我是谁**：验收腿，⛔ 实现者（实现腿＝`310-r1`）。本格判语**出自非实现者**＝本件与 `docs/evidence/s1/310-ac1-3-verdict.md`。
⛔ 我翻任何一枚框、⛔ 动 `Status:`、⛔ 把票改名 `-done`、⛔ push、⛔ 删盘上任何东西。锚＝起手 `618d984c`/`dev`（件＝`00-anchor.md`）。

## 三格判语（各一句＋支持它的逐字读数在哪）

| 格 | 判语 | 支持它的那一行逐字读数（⛔ 抄腿、⛔ 抄编排者，全是我这一刻现跑） | 件 |
|---|---|---|---|
| **`AC#1`** | **成立**（按已裁射程落干净）；⚠ 但"四枚入口件"这枚数在盘上就⛔ 是全量（真数五枚），漏的那一枚⛔ 动摇本格、⛔ 归本腿落地 | `portable-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013 - ...`（腿声称跑过的 5 枚发全带这行，且**那 5 枚根名逐枚在盘上**）⇄ `grep -rn --include=*.sh --include=*.ps1 'go test' scripts/ tools/`＝`hits=30 files=7`、其中唯一"名册外真起点"＝`tools/d22scan/runtests.sh:75` 逐字 `go test -v -count=1 "$@" >"$out" 2>&1` | `10-ac1-verdicts.md` §①②③ |
| **`AC#2`** | **成立**——"死亡形当场漏一枚、那一步在位则顶层不多且残渣逐字进专用根"被三条独立来路打中；⚠ 腿件 a3 那处**具名凭据⛔ 成立**（那枚目录名我从盘上取不到任何来路），错在名字⛔ 在判据 | 我自己复跑那两发：`a2_top_level_new=1 names_in_diff=TestRoundTripLoadMarshalLoad4222770164`（同行 `panic: test timed out after 2s`＋栈帧＋`created=2026-10-12 12:48:51.544` 且**今天仍在盘上**）⇄ `b_top_level_new=0` ＋ 根里残渣逐字 `  /c/Users/swq/AppData/Local/wisp-test-gotmp/v1-20261012-124853-148/TestRoundTripLoadMarshalLoad3712892429`；腿那枚 `2236627062`＝`grep -rl` 命中 8 枚件、**⛔ 一枚名册件**，那一毫秒在盘上的是 `…1721168370`（`created=11:49:06.507` ⇄ 腿写的 `external_kill_issued at=11:49:06.508`） | `20-ac2-paired.md` §1–§3 |
| **`AC#3`** | **成立**（四把门禁我自己逐把 rc=0、totals 行跑前⇄跑后逐字未变、十笔逐笔越界全在允许清单内）；⚠ 三处具名缺口：格式那格是**代用尺**（`gofumpt` 本机 `NOT_FOUND`）、`Q-310-r1-4` 那处 4⇄0 的**归因方向写反**、7 枚 0 字节件里 **1 枚**（`10c-run.run.txt`）判"那格没交" | `d22scan_rc=0`／`pathlen_rc=0`／`census_totals_diff_rc=0 diff_bytes=0`／`winsec_rc=0`＋`=== RUN=101  --- PASS=58  --- FAIL=0  --- SKIP=0`；拼写雷另发 `wrongspelling_rc=2`；越界＝`commit=3122306e paths=1 foreign_outside_allowlist=1`（那一枚⛔ 是别家路径、**是本票票面的 Progress bullet**＝派单字面清单漏了它一行）；eol 三来路 `CR_bytes worktree=334 head_blob=0 git_archive_tree=0` | `30-ac3-gates.md` §1–§6 |

**★第三种死法的归口（`Q-310-orch-N`，本格只答那一句）**：票面"测试正常跑完 ⇒ 0 残留"那句**要加射程**，落点＝**票 311 的 `AC#3`**（⛔ 塞进票 310 `AC#2`、⛔ 顺手实现清理器、票面原句⛔ 改）。我自己那把同行尺在默认根形上读到 **0 行** `TempDir RemoveAll cleanup`（＝与"默认根正常跑完不漏"一致），换卷那 4 行**我⛔ 重跑**（正典 §1 明令⛔ 换卷跑 winsec 档），只在件里具名引编排者那发。

**大白话（给机主看）**：那一步（把测试的临时目录从 `%TEMP%` 顶层改赶到一枚专用根）确实落进了三枚真正会起测试的脚本里，我亲眼跑了一遍、根名对得上盘上目录，测试进程被杀时 `%TEMP%` 顶层不再多目录、残渣落到专用根；编排者当初那张"四枚脚本"的名单少算了一枚（`d22scan` 那条路），这一条我没动、要下一程决定；腿有一处把"没存在过的目录名"写进了交付件，我用盘上的时间戳尺抓出来了，但"测试死了会漏一枚"这个结论本身经得起复跑。

## 我欠的读数（逐枚编号；⛔ 这些⛔ 成立，只写"没跑"）

- `Q-310-v1-0`＝**起手那枚假 0 形状**（本腿自己撞的、已在 `00-anchor.md` §2.2 具名结清）：四条 `tasklist` 被我第一次放进 `.ps1` 里跑，PowerShell ⛔ 认 bash 的 `//FI` 拼法 ⇒ 正控 `svchost.exe` 一起报 0。改正形＝**`tasklist` 归 bash、`Get-CimInstance` 归 `.ps1`，每条 0 读数同行带正控＋stderr 字节**。
- `Q-310-v1-1`＝**派单 §3 那张越界允许清单漏了本票票面那一枚路径**（`⛔ 面 Progress bullet 是交付件、又在允许清单外）；两种读法（越界 1 枚／0 枚）本腿都写进件里，⛔ 混着说。
- `Q-310-v1-2`＝**ubuntu 那一支的实测⛔ 交**（本机 `go env GOOS=windows`；`portable-tests.sh:163` 那行 `GOTMPDIR left as-is on GOOS=…` 我在 Windows 上永远走⛔ 到，非 Windows 那一支一行都⛔ 跑过）。⛔ 要这一格变绿，那需要 CI／容器色。
- `Q-310-v1-3`＝**真·`gofumpt` 的颜色⛔ 交**（`command -v gofumpt` ＝ `NOT_FOUND`；本腿与腿都用 `gofmt -l` 代尺，而 `gofmt` 比 `gofumpt` 松 ⇒ "导出树 0 枚"⛔ 等于"gofumpt 也 0 枚"）。
- `Q-310-v1-4`＝**第五枚入口⛔ 落⛔ 落**＝功能决策，归编排者一句话（`tools/d22scan/runtests.sh:75`，`scripts/d22scan.sh` 与 `ci.yml:82` 都经它；我已现量：我跑 `sh scripts/d22scan.sh` 那一发里 `GOTMPDIR` 命中＝0 ⇒ 那条路今天照旧漏）。⛔ 本腿动它＝落地活。
- `Q-310-v1-5`＝**`winsec-tests.sh`／`wisp-cli-tests.sh` 在 `LOCALAPPDATA` 为空（或非 Windows）的那一支是静默的**：两枚件的块是合取条件（`… && [ -n "${LOCALAPPDATA:-}" ]`），⛔ `DEGRADED`、⛔ `left as-is`，只靠委派的 `portable-tests.sh:158` 才响。那一支⛔ 一发实测过（我造 `DEGRADED` 时设的是 `WISP_TEST_GOTMP_ROOT`，走的是 portable 那层）。
- `Q-310-v1-6`＝**正典 §1 那行建议补的半句**（`突变台件自报的 dir= 字段只能当"那一形发生了"的动作行、⛔ 当漏掉的那枚目录名`）——**落点⛔ 归本腿**（那是正典面、只能编排者改），本腿只把判据形状写出来。
- `Q-310-v1-7`＝**腿把两枚读数件写进了专用根里面**（`D:/tmp/wisp-test-gotmp/shapec-115000-c2_redirected.run.txt`〔0 字节〕＋`….err`〔35458 字节〕）：票面⛔ 禁"写进根"（只禁删），但它让我起手的 `Get-ChildItem -Directory`＝6 枚 ⇄ `ls -1`＝8 条目，**同一枚"根里有几枚"的尺有两种读数**。要不要立一条"专用根里只准放 `t.TempDir()` 的产物"＝待裁。
- `Q-310-v1-8`＝**托管 runner 上 `%LOCALAPPDATA%` 兄弟根可不可写⛔ 我能交**（那需要 CI 色；本机可写已由我自己那 5 枚根现量，降级分支的 `DEGRADED` 我也自己造出来过一发）。
- `Q-310-v1-9`＝**`AC#2` 里我⛔ 复跑的发**：① 外力 `TerminateProcess` 那一形（腿 a3/a4/b4 三发）我⛔ 重打，我只重跑了 `-timeout` panic 那一形的 (a)(b) 两侧；② a1／b2 两发"帽落在⛔ 正在跑的 `t.TempDir()`"那两种作废形状我⛔ 重打；③ 我复跑的那两发是**直跑 `go test` 自控 `GOTMPDIR`**（机制级），"改前那发经真入口件＋中途杀测试二进制"那一形⛔ 跑过（腿也⛔ 跑过＝`Q-310-r1-6` 的 SIGINT 半边同族）。
- `Q-310-v1-10`＝**换卷那 4 行 `TempDir RemoveAll cleanup`＝本腿⛔ 复现**（正典 §1 明令⛔ 换卷跑 winsec 档，我照它⛔ 跑）⇒ 那一格的独立复现归票 311 台面。
- `Q-310-v1-11`＝**`scripts/slo-check.ps1` 整枚⛔ 跑**（与腿的 `Q-310-r1-5` 同一笔；我的依据只有两把尺：它 `go test` 命中 1 枚＝`:336` 注释、块内零调用点）。

**我没做什么（纪律自证）**：⛔ 改 `scripts/**`／`internal/**`／`cmd/**` 一个字节（我只**跑**它们）；⛔ 动 `ci.yml`／`.github/**`／`AGENTS.md`／`docs/PLAN.md`／`docs/specs/**`；⛔ 动票 304 的 Job Object 面；⛔ 放宽任何断言、⛔ 造 `--- SKIP`；⛔ 删任何目录/件（`%TEMP%` 下在册 58 枚〔现量含我那一枚＝59〕、两枚专用根的残渣、`TestProbeStepByStep3843095677` 全⛔ 动）；⛔ 新造协程；⛔ `--no-verify`（本仓 `.git/hooks/` 只有 `post-checkout`／`post-commit`，钩子缺位⛔ 等于许可）；每笔 `git add` 与 `git commit` **两处都带显式 pathspec**（走 `--pathspec-from-file`，路径逐枚列出）；⛔ `add -A`／`.`／整目录扫；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内⛔ worktree；临时件只建不删、只建在 `.scratch/wisp/probes/310/v1/**`；证据件一律 `.md`／`.txt`／`.tsv`（⛔ `.out` 结尾；⚠ 本腿有一枚台件副产物取名 `10-degraded-probe.out.txt`＝末段仍以 `.txt` 收尾、⛔ 被 `.gitignore` 第 8 行全仓 `*.out` 吃掉，`git check-ignore` 我核过）；`probes/**` 下⛔ 新建 `.go`；别家未提交那 842 行（`.gitignore`／`go.mod`／`design/**`）**一枚没动、没 add、没进我任何一笔**；起手暂存区 0 枚，全程⛔ 需要 `git restore --staged`。
