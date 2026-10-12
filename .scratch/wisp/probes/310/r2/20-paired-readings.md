# `20-paired-readings.md` — §1 两形改前⇄改后成对读数（每行带件名与 `rc=`）

台面：⛔ 拿旧数当基线。改前那一发在**仓外导出树**用 HEAD 原文跑（`git archive HEAD go.mod go.sum internal cmd tools scripts frontend`，`go_files_in_tree=683` 非空先证、`archive_rc=0`，件 `logs/10-export-tree-path.txt`）；⛔ 在共享工作树跑改前那一发。★本腿另证四枚目标文件 blob 在 `8e3cb370`⇄现 `HEAD 73719706` 逐枚 SAME（`git diff --name-only 8e3cb370 HEAD -- scripts tools`＝空）⇒ 建导出树时的 HEAD 原文＝现 HEAD 原文，改前读数对现 HEAD 仍成立。
专用根默认落 `${LOCALAPPDATA}/wisp-test-gotmp/<entry>-…`（与 `%TEMP%` 同卷），⛔ 换卷、⛔ 指进仓内。

## ① `Q-310-orch-P`：`sh scripts/d22scan.sh` 的 `GOTMPDIR=` 成对
| 形 | 台 | 读数 | `rc=` | 件 |
|---|---|---|---|---|
| 改前 | 工作树（runtests.sh＝HEAD 原文，`grep ticket 310`＝0、`grep GOTMPDIR`＝0） | `GOTMPDIR=` 命中 **0 行**（且 d22scan 正控跑满：`d22scan: clean`） | `worktree_d22scan_preedit_rc=0` | `logs/05-worktree-d22scan-preedit-rc.txt`／`.stdout.txt` |
| 改前 | 仓外导出树 | `GOTMPDIR=` 命中 **0 行**；★整发 `rc=1`＝`TestScannerSelfScanOfRealRepoIsGreen`/`TestScopeReportMatchesRealCoverage`/`TestRealRepoLedgerIsHonest` 三枚读**全仓**的 self-scan 用例在六件套导出树上假红（⛔ temp-dir 无关），`gotmpdir_hits=0` 才是牙信号 | `tooth1_before_rc=1`（假红，具名 `logs/10-tooth1-before-diagnosis.txt`） | `logs/10-tooth1-before-rc.txt`／`-diagnosis.txt` |
| 改后 | 工作树 | **恰一行** `runtests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/d22scan-20261012-135432-33 - this run's t.TempDir() family lands here; nothing deletes it.` | `tooth1_after_rc=0`，`gotmpdir_eq_lines=1` | `logs/10-tooth1-after-rc.txt`／`10-tooth1-after-gotmpdir-line.txt` |

★那一行的根**当场在盘上存在**且⛔ 落 `%TEMP%` 顶层（逐条 `[ -d ]`＋PowerShell `Test-Path -PathType Container`）：
```
bash  -d /c/Users/.../d22scan-20261012-135432-33  -> bash_dir_test=EXISTS rc=0
verify-root.ps1 -> ROOT_EXISTS_DIR=True  ROOT_UNDER_TEMP=False  ROOT_TEMPNAME_IN_TEMP_ROSTER=False
名册尺 Get-ChildItem -Path $env:TEMP -Directory  (件 logs/99-temp-roster-end.txt，COUNT=59 未含此根)
```
件：`logs/10-tooth1-after-verifyroot.txt`（`verify_ps_rc=0 errbytes=0`）。

## ② `Q-310-orch-Q`：抽掉 `LOCALAPPDATA` 那一发的成对
| 形 | 台 | 读数 | `rc=` | 件 |
|---|---|---|---|---|
| 改前 | 工作树 `env -u LOCALAPPDATA bash scripts/winsec-tests.sh ./internal/proc/`（`winsec-tests.sh` 合取 `&& [ -n LOCALAPPDATA ]`＝改前形） | 整块静默：stdout **0 B**、stderr 只有 GUARD 1 两行（`this step does not name ./internal/winsec/`）、`ticket310/GOTMPDIR/DEGRADED` 命中 **0**（件 `logs/20-tooth2-before-hits.txt` 记 0） | `tooth2_before_worktree_rc=2` | `logs/20-tooth2-before-wt-rc.txt`／`-wt.err`(268B) |
| 改前 | 仓外导出树（同抽法、同拼法） | 同：0 命中、GUARD 1 rc=2 | `tooth2_before_rc=2` | `logs/20-tooth2-before-rc.txt`／`.err` |
| 改后 | 工作树，同一台面同一抽法 | 打出**那一行新告警**：`winsec-tests.sh: ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty, … GOTMPDIR stays unset.`；`GOTMPDIR=` 命中 **0**（根仍未设）；stderr 268→509 B（多的就是那一行）；`rc` 与改前**同数** | `ws_after_rc=2 (before=2)` | `logs/20-tooth2-after-ws-rc.txt`／`-ws-warnline.txt` |

★同一告警在另外两枚入口也各自响（证四枚同一形，非恒真）：
```
portable（对照，改前已响、本腿⛔ 动它一字）： env -u LOCALAPPDATA bash scripts/portable-tests.sh --scope=census
  -> DEGRADED-LOCALAPPDATA-empty 一行响、GOTMPDIR=0、census totals 逐字未变   pb rc=0（before=0）件 logs/22-tooth2-after-pb-*.txt、22-tooth2-portable-before-line.txt（改前同行）
runtests（本腿新增那一枚）： env -u LOCALAPPDATA sh tools/d22scan/runtests.sh -C tools/d22scan ./...
  -> DEGRADED-LOCALAPPDATA-empty 一行响、GOTMPDIR=0，且 go test 照跑 rc=0（`OK PASS=35 FAIL=0 SKIP=0`）＝⛔ fail-closed／⛔ 拒绝跑测试  rt rc=0 件 logs/24-tooth2-after-rt-*.{txt}
```
- 抽法自证：`env -u LOCALAPPDATA sh -c 'echo ${LOCALAPPDATA:-UNSET}'`＝`UNSET`（改前／改后同）。
- ⛔ 把"没输出"当证据：改前那一发有完整 stderr（268 B GUARD 1 行）＋`rc=2`；命中 0 是 grep 结果并排落盘，非空文件。
- ★cli 那一枚（`wisp-cli-tests.sh`）未跑"抽掉 LOCALAPPDATA"那发（跑了会连 cmd/wisp 全套一起，⛔ 由本腿归因、且非本腿牙所需）；其告警正文＝winsec/runtests 逐字节同（件 `logs/15-warn-body-normalized.txt`），本腿⛔ 声称实测、只给同一形尺——此为具名代用（见 `90-final.md`）。
