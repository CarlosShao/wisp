# `90-final.md` — 票 310 增量程 `310-r2` 收尾一页（实现腿）

**我是谁**：实现腿（写码者）＝票 310 `Q-310-orch-P`／`Q-310-orch-Q` 两笔欠账。⛔ 裁决者＝下一程非实现者腿；⛔ 我翻任何框、⛔ 动 `Status:`、⛔ 改名 `-done`、⛔ push、⛔ 删任何目录。锚＝`00-anchor.md`（起手 `73719706`/`dev`，派单锚 `27176c2c` 为其祖先）。

## 1. 提交号 + 每笔 pathspec 枚数（`git add` 与 `git commit` 双显式 pathspec，走 `--pathspec-from-file`，清单入库）
| 笔 | hash | pathspec 枚数 | `git show --name-only --format=` 越界 |
|---|---|---|---|
| 代码（三枚 shell） | `78782c20` | 3 | `scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`／`tools/d22scan/runtests.sh`＝**3 枚，名册外 0** |
| 票面 bullet | `5ef679b5` | 1 | 票 310 一枚＝**1 枚**（`git show --numstat`＝1/0，纯追加） |
| 证据件 `probes/310/r2/**` | `a2b5a8f4` | 124 | 全在 `.scratch/wisp/probes/310/r2/**`，名册外 **0** |
| 逐笔越界读数（自指补笔） | `6709b0fc` | 2 | `logs/40-postcommit-nameonly-evidence.txt`＋`logs/commit-list-evidence2.txt`，名册外 **0** |

★逐笔 `git show --name-only` 越界读数（⛔ 区间尺）＝`logs/40-postcommit-nameonly-evidence.txt`；总越界尺＝`git diff --name-only 27176c2c..6709b0fc -- .github internal cmd *_test.go thresholds.go allowlist.txt scripts/slo-check.ps1`＝**0 枚**（件 `logs/42-forbidden-audit.txt`），`'*.go'` 越界＝**0 枚**（件 `logs/42-gochg.txt`），别家未提交（`.gitignore`/`go.mod`/`design/**`）进我任何一笔＝**0 枚**。

- 允许清单（派单 §4）＝`probes/310/r2/**` ＋票面 bullet ＋ `scripts/portable-tests.sh`／`scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`／`tools/d22scan/runtests.sh`。
- ★`scripts/portable-tests.sh` 本腿 **0/0 未改**（它 r1 起就带 windows+空-LOCALAPPDATA 的 `else`，早已"说实话"，是这一形的参照；blob 未变），⛔ 越出允许清单、⛔ 为凑"改了四枚"去动它。
- 提交后本腿索引清空：`git diff --cached --name-only`＝空（每笔都用 `--pathspec-from-file` 限定，⛔ 让共享索引跨腿窗口挂着别家的东西）。

## 2. §1 两形各自"响在哪一行"
- **①（P，第五枚入口）**：改后 `sh scripts/d22scan.sh` 里那一行 ＝ `tools/d22scan/runtests.sh:95`（`echo "runtests.sh: ticket 310 AC#1: GOTMPDIR=$GOTMPDIR - this run's ..."`）实跑打出 `runtests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/d22scan-20261012-135432-33 ...`，恰 1 行、根当场 `[ -d ]` 存在、⛔ 落 `%TEMP%` 顶层（件 `logs/10-tooth1-after-verifyroot.txt`）。改前该发 GOTMPDIR 命中＝0（工作树 HEAD 原文 rc=0；导出树 rc=1＝三枚 self-scan 用例假红，均命中 0）。
- **②（Q，静默支）**：那一行新告警＝`winsec-tests.sh:74` / `wisp-cli-tests.sh:87` / `runtests.sh:103`（+ portable `:158` 原已响），内容 `…: ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty, … GOTMPDIR stays unset.`；抽掉 `LOCALAPPDATA` 那发：改前 winsec **0 命中**（跑到 GUARD1、rc=2、stderr 268 B 非空＝⛔ 拿"没输出"当证据）、改后同台面 rc 仍=2 且打出该告警、GOTMPDIR 仍不设。runtests 独立一发打出告警后 go test 照跑 rc=0（⛔ fail-closed／⛔ 拒绝跑测试，件 `logs/24-tooth2-after-rt-rc.txt`）。四枚告警正文归一后逐字节 SAME（件 `logs/15-warn-body-normalized.txt`）。

## 3. 欠账名册（本腿⛔ 能交的，具名）
- `Q-310-r2-1`＝**ubuntu／非 windows 那一支的实测⛔ 交**：本机 `go env GOOS=windows`，四枚块的 windows 判定为假时整块跳过、行为与改前逐字节相同（本腿在 Windows 上永远走⛔ 到非 windows 支）；那一支要变绿需 CI/容器色。（与票 `Q-310-v1-2` 同一笔，本腿未扩射程去动它。）
- `Q-310-r2-2`＝**cli 的 empty-LOCALAPPDATA 告警未做独立实跑**：`wisp-cli-tests.sh --scope=cli` 抽掉 `LOCALAPPDATA` 会连 cmd/wisp 全套一起跑（⛔ 本腿归因、且非本腿两枚牙所需）。本腿只以"正文与 winsec/runtests 逐字节同"这把同一形尺＋ winsec/portable/runtests 三枚实测覆盖；此⛔ 是实测、是代用尺，具名。
- `Q-310-r2-3`＝**真·gofumpt 颜色⛔ 交**：本机 `gofumpt NOT_FOUND`⇒ 代用 `gofmt -l`（比 gofumpt 松）；"导出树 0 枚"⛔ 等于"gofumpt 也 0 枚"。（`Q-310-v1-3` 同款缺口。）
- `Q-310-r2-4`＝**托管 runner 上 `%LOCALAPPDATA` 兄弟根可不可写⛔ 我能交**（需 CI 色；本机可写已由本腿多枚根现量，降级/空支的 DEGRADED 我也各自造出过一发）。
- `Q-310-r2-5`＝**导出树改前 `--scope=cli` 假红未复现到"与改后同三名"**：工作树被别家未提交 `design/**` 删除污染 ⇒ 导出树(pristine)⇄工作树＝两枚台面，名字作差⛔ 算凭据；本腿改用"runtests 在委托路 0 行输出＝对 cli 可证零效应"这把台面无关的尺定"没多也没少"（件 `logs/38-cli-inertness-verdict.txt`）。既有真窗/真麦/驻留红族⛔ 由本腿归因。

## 4. 我⛔ 做到的（纪律自证）
⛔ push（只 commit）；⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；`git add` 与 `git commit` 两处均带显式 pathspec（`--pathspec-from-file`，逐枚列）；⛔ `add -A`/`.`/整目录扫。⛔ 翻 AC 框、⛔ 动 `Status:`、⛔ 把票改名 `-done`。⛔ 动 `.github/**`／`ci.yml`（`git diff --name-only HEAD -- .github internal cmd`＝空，件 `logs/15-no-forbidden.txt` 注 EMPTY-OK）。⛔ 动任何 `*_test.go` 断言、⛔ 动 `internal/**`（`git diff --name-only HEAD -- '*.go'`＝0，件 `logs/36-godiff.txt`）。⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43／C1–C32／D1–D47。⛔ 动票 304 的 Job Object 面、⛔ 动票 311 归因。⛔ 删任何目录/件：`%TEMP%` 顶层在册 59 枚一枚没删/没改名/没碰（`comm -23`＝0 消失）；`TestProbeStepByStep3843095677` 仍按"归属不明、⛔ 动、永久"。⛔ 动非 windows 那两支的行为。共享工作树里别家未提交（`.gitignore`/`go.mod`/`design/**`/别人 `probes/**`：152/161/242/268）**⛔ 进我任何一笔、⛔ revert、⛔ 修改**；起手 `git diff --cached`＝空，全程⛔ 需要 `restore --staged`。临时件只建不删（本腿件全在 `probes/310/r2/**`；`/tmp` 下三枚导出树/探针脚本留档⛔ 删）。读数件⛔ 写进任何专用根（全落 `logs/`；数根条目用具名 `Get-ChildItem -Directory`，件 `logs/99-leg-root-placement.txt`）。⛔ 0 字节尺件（原 24 枚空件逐枚就地注 `EMPTY-OK + 同伴 rc=`，件 `logs/zz-empty-files.txt`＝真值剩 0）。

## 5. 我⛔ 做到的与代用的尺 / 自报越格
- 代用尺具名：`gofumpt`→`gofmt -l`（松，见 `Q-310-r2-3`）；`cli empty-LOCALAPPDATA`→同一形尺代实测（`Q-310-r2-2`）。
- ★**自报越格**（做了⛔ 派单明确要求的，⛔ 自盖"没影响"）：
  1. 本腿把 `winsec`/`cli`/`runtests` 三枚告警正文**归一化对齐**（逗号、`this step has no measured-safe root` 措辞）以达到"四枚同一形"——派单只说"同一形"，对齐正文是实现该要求所必需，⛔ 新增行为、⛔ 动 portable 一字。
  2. 本腿对 `runtests.sh` 那枚新增块**改了注释换行**（把被硬折断的 `GOTMP DIR`/`scripts/ winsec-tests.sh` 复原成不断词的整句）——纯注释面，⛔ 动逻辑；属"把新加的块写干净"，仍列出自报。
  3. 本腿**没跑** winsec 全档在"抽掉 `LOCALAPPDATA`"台面下的那一发（改后用 GUARD-reject 窄台证同一支响）；`--scope=cli` 也未跑空-`LOCALAPPDATA` 独立发（见 `Q-310-r2-2`）。
  4. 本腿多造了几枚 `%LOCALAPPDATA%\wisp-test-gotmp\*-20261012-*` 专用根（d22scan/winsec/portable/cli 等，含一枚 `_probe_inner.sh` 的 `PROBE-` 试验根）与 go-build 缓存——**只建不删**，逐枚留盘。
  5. 起手发现工作树里 `git status` 已有一枚未跟踪的 `-` 文件与 `.scratch/.scratch/` 等别家/历史残渣——本腿⛔ 碰、⛔ 提交、⛔ 清。
