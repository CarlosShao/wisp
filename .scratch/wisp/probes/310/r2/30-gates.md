# `30-gates.md` — §3 逐把尺（每件自落 `rc=N`；⛔ 0 字节尺件）

台面＝工作树（起手 HEAD `73719706`＋本腿代码提交 `78782c20` 的三枚 shell 增量＋本腿件）。"应当为 0"的链一律 `|| true` 并把 stderr 与命中数并排落盘。

## 0. 起手先跑的一把（⛔ 建新件后才量）
`sh scripts/check-path-length-budget.sh`（无 `--with-self-test`）在写任何 r2 件之前先跑一把 ⇒ `rc=0`；随后带自测再跑一把 ⇒ 同 `rc=0`（见 §2）。四枚目标文件 blob 在 `8e3cb370`⇄现 `HEAD` 逐枚 SAME（见 `00-anchor.md`），且 `git diff --name-only 8e3cb370 HEAD -- scripts tools`＝空 ⇒ 改前读数对现 HEAD 成立。

## 1. 两枚硬闸
| 闸 | 命令（逐字） | `rc=` | 件（末行自落） |
|---|---|---|---|
| D22 扫描 | `sh scripts/d22scan.sh` | **`rc=0`**（＝§1① 的改后凭据；`d22scan: clean`，且那一发带**恰一行** `GOTMPDIR=`，见 `20-paired-readings.md` ①） | `logs/99-d22scan-final-rc.txt`（`d22scan_final_rc=0`）＋`logs/10-tooth1-after-rc.txt` |
| 路径长度预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | **`rc=0`** | `logs/30-pathlen-rc.txt`（`pathlen_rc=0`）＋`logs/30-pathlen.stdout.txt`（`tracked paths=11630  over-budget=57  covered by roster=57  not in roster=0`，`VERDICT GREEN`；`control 2/3 ok`） |

## 2. `--scope=census` totals：改前（仓外导出树，无本腿步子）⇄改后（工作树）
```
改前: portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0   rc=0  件 logs/30-census-before-totals.txt
改后: portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0   rc=0  件 logs/30-census-after-totals.txt
作差: diff 改前 改后  -> totals_diff_rc=0（逐字未变）  件 logs/30-census-totals-diff.txt
```
与票面基线 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` **逐字相同**。

## 3. winsec 档（⛔ `--scope=winsec` 拼法，用正典拼写）
```
bash scripts/winsec-tests.sh ./internal/winsec/   rc=0   件 logs/32-winsec-rc.txt
颜色: winsec-tests.sh: four numbers: === RUN=101 --- PASS=58 --- FAIL=0 --- SKIP=0   （逐字未变）
★本发带恰一行 GOTMPDIR（winsec 自开根）: winsec-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/winsec-20261012-135756-111
   GOTMPDIR 行计数=1（件 logs/32-winsec-gotmpdir-count.txt）；GUARD 3 审计行在位（logs/32-winsec-guard3.txt）
```
⛔ 把根换到 `D:` 卷跑本档（正典禁令、票 311 射程）——本腿全程用默认 `%LOCALAPPDATA` 同卷根，改前⇄改后同根形状。

## 4. cli 档（`bash scripts/portable-tests.sh --scope=cli`；既有真窗/真麦/驻留一族＝⛔ 本腿射程、⛔ 由本腿归因；本腿只证"没多也没少"）
**"没多也没少"＝用一把台面无关的尺定，⛔ 拿两枚不同台面作差**（工作树被别家未提交的 `design/**` 删除污染 ⇒ 导出树-pristine ⇄ 工作树＝两枚台面）：
```
本腿唯一能进 --scope=cli 这条链的改动＝runtests.sh 的根步子；portable 先 export GOTMPDIR ⇒ 守卫 -z 恒假 ⇒ runtests 整块跳过。
实测：`runtests.sh: ticket 310` 行数   改前(HEAD stepless)=0   改后(本腿 step)=0   （件 logs/36-cli2-before.stdout.txt / logs/38-cli2-after.stdout.txt）
⇒ 带步子 ⇄ 无步子喂给 go test 的 GOTMPDIR 与命令逐字节相同 ⇒ 本腿对 cli 结果＝可证零效应，⛔ 能添/削任何一枚红。
```
- 工作树 改后 `--scope=cli` 两发：`RUN=398 PASS=284 FAIL=3 SKIP=1`（shot1 `logs/34-cli-after-rc.txt`）／`RUN=398 PASS=283 FAIL=3 SKIP=2`（shot2 `logs/38-cli2-after-rc.txt`）；红族＝`TestAC13ColdStart…`/`TestAC14GoSideEval…` + (`TestAC4FocusReturn…` ⇄ `TestPanelHostRealWindowHop…` 随真窗时序翻动) ＋ SKIP `TestAC247LiveMicrophone…`/`TestPanelHostLatencyPercentiles…`＝本机既有真窗/真麦/驻留族，⛔ 本腿归因（两发之间成员翻动＝该族自带 flaky）。
- 具名记一笔：曾试图拿"仓外导出树 改前"当基线——第一发六件套（无 design）得 `FAIL=18`，第二发加 `design/docs` 仍 `FAIL=19`；成因＝导出树是 pristine HEAD、工作树被别家 design/** 删除改写，两枚⛔ 同台面 ⇒ 该作差⛔ 算凭据，只作台面假红留档（件 `logs/34-cli-baseline-note.txt`、`logs/38-cli-inertness-verdict.txt`、`logs/34-cli-runtests-inert.txt`）。

## 5. `%TEMP%` 顶层名册尺（PowerShell 走 `-File`，⛔ 内联 `$_`）
```
尺: powershell -NoProfile -File logs/roster-test.ps1
起跑前: COUNT=59  all_top_level_dirs=2839   件 logs/00-temp-roster-start.txt（含逐枚 CreationTime）
交件前: COUNT=59  all_top_level_dirs=2851   件 logs/99-temp-roster-end.txt
comm -13(新增)=0  comm -23(消失)=0          件 logs/99-roster-comm-new.txt / -gone.txt（正文已注 EMPTY-OK；真值 new=0/gone=0 重算见 zz-roster-*-names.txt）
本腿造的专用根全落 %LOCALAPPDATA%\wisp-test-gotmp（<entry>-20261012-* 共 15 枚），%TEMP% 顶层今日 leg-root 枚数=0（件 logs/99-leg-root-placement.txt）
★⛔ 那枚 TestProbeStepByStep3843095677（已裁"归属不明、⛔ 动、永久"）与在册 59 枚一枚没删/没改名/没碰。
```

## 6. 四枚 `tasklist`（bash 里跑，⛔ 进 `.ps1`＝假 0）＋家族尺（PowerShell `-File`）
```
起跑前 / 交件前（件 logs/00-gate-{start,end}-tasklist.txt）:
  wisp.exe rc=0 hits=0 | balldebug.exe rc=0 hits=0 | mockllm.exe rc=0 hits=0 | wisp.test.exe rc=0 hits=0
  正控同行 svchost.exe=84（证明解析器有牙、0 是真 0 不是假 0；stderr_bytes=0）
家族尺（Get-CimInstance *.test.exe，件 logs/00-gate-{start,end}-family.txt）: FAMILY_TEST_EXE_COUNT=0（errbytes=0）
```

## 7. ★格式名册（⛔ 恒真）＋代用尺具名
- `command -v gofumpt`＝**NOT_FOUND**（本机⛔ 装 gofumpt、⛔ 联网）⇒ **代用尺＝随 Go 装的 `gofmt -l`**（射程同一枚 `cmd internal tools scripts`），比 gofumpt **松**。
```
工作树: gofmt -l cmd internal tools scripts  -> 4 枚: cmd/wisp/models.go, internal/agent/approval/pending_read.go, internal/agent/tools.go, internal/risk/provenance.go   件 logs/36-fmt-worktree.txt
仓外导出树(HEAD原文): 先证 find '*.go' = 683 非空，再 gofmt -l -> 0 枚                件 logs/36-fmt-head.txt + logs/10-export-tree-path.txt
那 4 枚归谁: git diff --name-only HEAD -- '*.go' = 0 枚（件 logs/36-godiff.txt 注 EMPTY-OK）⇒ 4 枚⛔ 由本腿造成（本腿只动 shell），＝本机检出面 CRLF 违 `eol=lf` 那一族（`Q-310-r1-4`/`A859`，⛔ 本票射程）。
```
⚠ 代用尺买到⛔ 真·gofumpt 颜色（`Q-310-v1-3` 同款缺口）：那"0 枚"半边⛔ 能推广到 gofumpt；CI 色那半押在推送批次。
