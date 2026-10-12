# `30-gates.md` — 票 310 `AC#3` 门禁与越界（腿 `310-r1`，实现者自跑；判语归非实现者腿）

台面＝工作树（HEAD `7e6c1579` ＋ 本腿三枚 `scripts/**` 增量 ＋ 本腿件）。每件自落 `rc=N`；⛔ 0 字节尺件；"某枚计数应当为 0"的链一律 `|| true` 并把 stderr 与命中数并排落盘。

## 1. 两枚硬闸

| 闸 | 命令（逐字） | 读数 | rc | 件 |
|---|---|---|---|---|
| D22 扫描 | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=231, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ...` | **`rc=0`** | `logs/20-d22scan.txt`（末行 `d22scan_rc=0`） |
| 路径长度预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` ＋ `unit cross-check: the byte view and the character view select the SAME 57 paths` | **`rc=0`** | `logs/20-pathlen.txt`（末行 `pathlen_rc=0`） |

新建件之前也跑过这把长度门（本腿末段名全部≤100、相对路径≤121、全路径≤165；最长一枚＝`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\310\r1\logs\...` 那批，由上面那一发 `rc=0` 一并覆盖）。

## 2. ★格式名册两把并排（工作树 ⇄ HEAD blob）

⚠ 本机⛔ 装 `gofumpt`（`command -v gofumpt` 空，本腿⛔ 联网装），**代用尺＝随 Go 一起装的 `gofmt -l`**，射程同一枚 `cmd internal tools scripts`；替换这件事本身写在这里给裁决腿审。

| 把 | 射程 | 枚数 | 说明 |
|---|---|---|---|
| ① 工作树 | `gofmt -l cmd internal tools scripts` | **4** | 逐名册＝`logs/20-fmt-worktree-names.txt`：`cmd\wisp\models.go`、`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`、`internal\risk\provenance.go`（四枚都 `git ls-files --error-unmatch`＝tracked、且⛔ 在本腿的改动清单里） |
| ② HEAD blob | `git archive HEAD cmd internal tools scripts go.mod go.sum \| tar -x` 到仓外树后再 `gofmt -l cmd internal tools scripts` | **0** | ★先证取出的东西非空：`head_tree_go_files=`**682**（`logs/20-fmt-head.txt`，`archive_rc=0`）⇒ 那枚 0 是"取出来的字节本来就干净"，⛔ "取空了所以 0" |
| 差的那 4 枚归谁 | `git diff --name-only HEAD -- '*.go'` ＝ **0 枚**（`logs/20-fmt-godiff.txt`，`empty above＝0 枚`） | — | ⇒ 工作树⇄HEAD 这一处⛔ 由本腿造成（本腿只动三枚 `.sh`，见 §4）；4 枚的成因＝检出面与 blob 的行尾／工作树形状，归尺名册、⛔ 本票射程（登记 `Q-310-r1-4`） |

## 3. `--scope=census` 的 totals 行：跑前／跑后各一发，作差

**⛔ 拿旧数当基线** ⇒ "跑前"那一发用 HEAD 的**脚本原文**在仓外导出树里现跑（`git archive HEAD go.mod go.sum internal cmd tools scripts frontend \| tar -x` → `go_files=`**683**、`grep -c 'ticket 310' `＝**0** 枚（证那一枚树里⛔ 有本腿的步子）、`archive_rc=0`）。

```
跑前（HEAD 导出树，无本腿步子）：packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0    rc=0   件 logs/30-census-prepatch.txt
跑后（工作树，本腿步子在场）    ：packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0    rc=0   件 logs/30-census-postpatch2.txt
作差尺：diff logs/30-census-prepatch-totals.txt logs/30-census-postpatch-totals.txt  →  totals_diff_rc=0（逐字未变）
```

（⚠ 第一发跑前读数在 `/d/tmp/310r1-tree.Ttzxdu` 那枚树上⛔ 成立——导出⛔ 带 `frontend/`，`go list ./...` 退 1，件 `logs/30-census-prepatch.txt`(旧) 已被带 `frontend` 的那一发覆盖为 rc=0 版本；那一发按"造不出世界即作废"的规矩⛔ 计凭据，只把成因写在这里。改后工作树的第一发＝`logs/30-census-postpatch.txt`，totals 同样逐字未变。）

★**跑法变了没／颜色变了没**（派单那句具名回报，本腿跑了三枚入口件；第四枚 `slo-check.ps1`⛔ 跑、理由见 `90-final.md` 的 `Q-310-r1-5`）：

| 入口 | 跑法 | 颜色（改前⇄改后） | 件 |
|---|---|---|---|
| `scripts/portable-tests.sh --scope=census` | 多打一行 `ticket 310 AC#1: GOTMPDIR=...`（39 行 vs 38 行），totals 一字未改 | 两侧 `rc=0` | `logs/30-census-prepatch.txt`／`logs/30-census-postpatch2.txt` |
| `scripts/winsec-tests.sh`（窄档，真起 `go test`） | 同上多一行 GOTMPDIR；`--scope=winsec` 解析、GUARD A/B 全部照跑 | 改前（HEAD 导出树）`RUN=101 PASS=58 FAIL=0 SKIP=0` ⇄ 改后（出厂根）`RUN=101 PASS=58 FAIL=0 SKIP=0` ⇒ **颜色一字未变**；★但把根建在**另一枚卷**（第一稿 `D:/tmp/...`）时会变：`PASS=56 FAIL=2` | `logs/11-winsec-headtree.txt`／`logs/12-winsec-shipped.txt`／`logs/11-winsec.txt` |
| `scripts/wisp-cli-tests.sh`（cmd/wisp 档，经 `--scope=cli`） | 同上多一行 GOTMPDIR | 改后 `RUN=398 PASS=285 FAIL=2 SKIP=1 rc=1` ⇄ 同台面把根**指回 OS 临时根**（＝改前的位置）`RUN=398 PASS=285 FAIL=2 SKIP=1 rc=1` ⇒ **同一红、同三名**：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14GoSideEvalPushReachesThePage`、`--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam`（本机既有红／真窗／真麦那一族，⛔ 本票射程，⛔ 由本腿归因） | `logs/12-cli-shipped.txt`／`logs/12-cli-noredirect-cf.txt` |
| `%TEMP%` 顶层枚数（这几发入口跑下来） | 名册尺同 §2 那把 | `logs/12-roster-before.txt`=58 → `logs/12-roster-after.txt`=**58**，差集⛔ 空＝**0 新增**（`empty above＝0 枚`，`logs/12-roster-new.txt`） | 见左 |

## 4. ★越界尺（逐笔，⛔ 区间尺）

`git show --name-only --format=` 逐笔自扫（本腿每一笔都是 `add`＋`commit` 双 pathspec）：允许清单＝`.scratch/wisp/probes/310/r1/**` ＋ `scripts/portable-tests.sh`／`scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`。越界枚数与每笔读数记在 `90-final.md` §3（本腿交件后当场逐笔数过，见那一张）。共享工作树里机主/别家那批未提交改动（`.gitignore`／`go.mod`／`design/**`）⛔ 提交、⛔ revert、⛔ 修改，也⛔ 进过本腿任何一笔。
