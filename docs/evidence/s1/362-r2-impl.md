# 票 362 续跑腿 `362-r2` 交件 —— 第二道门：`attrib.sh` 的 ruler A 加**具名排除** `.scratch/**`（⛔ 洗任何字节）

- 腿＝`362-r2`（**写腿／实现者**）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
- 写面实动：`.scratch/wisp/probes/161/r5/attrib.sh`（**ruler A 的输入面 ＋ 随它失真的分母文字**）＋ `.github/workflows/ci.yml`（**`362-r1` 那步注释里已腐坏的计数与指向 ＋ `:436` 邻域那句**）＋ 本件。
- ★**`.go` 字节零枚**：`git diff -- '*.go'` 里那 5 枚 `internal/memory/**` 是**我 10:03 开工前 `git status` 就已在场的别人的未提交改动**（§2②）。
- 台面锚＝**导出树 `938216b8e`**（`git archive` 建于 10:0x，自带 index、17609 件、`.go` **1162**、`*.go` 的 CR＝**0**、`git status` 干净）。★HEAD 被别的腿推着走（`630ac8b21`→`938216b8e`→`8fa66ce28`）⇒ 可复跑命令写**具名提交号**、⛔ `HEAD`。
- 环境：`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；`gofumpt` 一律绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`（`v0.12.0 (go1.27.1)`）。尺件＝`D:\tmp\wisp\20261023\362r2\logs\01…13b` ＋ 同名 `.sh`。

## 0 ★机主批准（逐字）

> 「好的，这三件事我都允了，同意，你可以继续了」

出处＝台账 `docs/reports/pending-and-issues.md` **`A1161`①**（`362-r1` 件 §0 已逐字引过同一条）。本腿据此动 `.github/**`。

## 1 ★改前：仓外导出树上 `--tracked-only` 的逐字（`rc=1`、lines=20 files=20）

shape（★**⛔ `cd` 进树**；root 取自**脚本自身目录**，`A1139①`）＝`ATTRIB_ROOT=<树> sh <主检出>/.scratch/wisp/probes/161/r5/attrib.sh --tracked-only`
```
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 .../gofumpt.exe -l   (tracked .go files handed to it: 1162)
attrib.sh: (A) tracked-tree  lines=20 files=20  rule: MUST be empty
attrib.sh: RED (tracked-only) - tracked-dirty=20 files_dirty=20. The CI gofmt step and this step read the same denominator; this one additionally refuses to answer when handed 0 files.
== rc=1 ==       （20 个 TRACKED-DIRTY 行体例：`   TRACKED-DIRTY<TAB>路径<TAB>(last commit touching it: …)<TAB><== CI IS RED ON THIS`）
```
20 枚的**路径字段逐字**（`awk -F'\t' '{print $2}'`；顺序＝`gofumpt` 打印顺序）：
```
.scratch/wisp/probes/163/a1/main.go
.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go
.scratch/wisp/probes/183/r2/mut/task-boxset-off.go
.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go
.scratch/wisp/probes/197/r1c/pre/subagent_197.go
.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go
.scratch/wisp/probes/212/v1/mut/main-noq9.go
.scratch/wisp/probes/220/r1/prechange/prechange220_probe_test.go
.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go
.scratch/wisp/probes/224/v2/dialect_probe_test.go
.scratch/wisp/probes/231/v1/mutation/config_reload_branch_after_prefix.go
.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go
.scratch/wisp/probes/241/v1/posctl/badly_formatted.go
.scratch/wisp/probes/241/v1/probe_v1_readings_test.go
.scratch/wisp/probes/253/v2/denoms253v2.go
.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go
.scratch/wisp/probes/263/v1/src/main.go
.scratch/wisp/probes/33/p1/q1/main.go
.scratch/wisp/probes/33/p1/q2/main.go
.scratch/wisp/probes/33/p1/q3/main.go
```
归属两把尺（★**`[\\/]`、⛔ 正斜杠单写**）：`tr '\\' '/' | grep -c '^\.scratch/'` ＝ **20**、`grep -vc` ＝ **0**。

## 2 ★开工前三条硬检查

| 检查 | 现跑读数 | 判 |
|---|---|---|
| ① `git diff --cached --name-status` | 只有 `M .scratch/wisp/issues/362-….md`（票面本身）；**无 `D`／`A`** | ✔ 不触发"别人的 staged 被卷走" |
| ② 写面内别人的未提交改动 | `git status --porcelain -- .github/workflows/ci.yml .scratch/wisp/probes/161/r5/attrib.sh` ＝ **空** | ✔ 写面干净 |
| ③ 票面新鲜度 | `git ls-files '*.go'` ＝ **1162**（票面 `AC#0` 记 1156、`362-r1` 记 1160 ⇒ **又漂**）；五目录 **802**（internal 623／cmd 149／tools 18／scripts 11／frontend 1）、`.scratch` **360**、根 **0**、其它首段 **0** ⇒ 802+360＝1162 无漏类 | ✔ 前提成立、只是数字旧 |

## 3 ★改了什么（逐行 diff）

**(a) 排除面**——`attrib.sh` 的**全部代码改动就这 5 行 ＋ 1 行打印**：

```diff
+A_N_EXCLUDE=':(exclude).scratch/**'
-    TRACKED_GO=$(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ') || TRACKED_GO=0
+    TRACKED_GO_ALL=$(git ls-files -z '*.go' | tr -dc '\0' | wc -c | tr -d ' ') || TRACKED_GO_ALL=0
+    TRACKED_GO=$(git ls-files -z '*.go' "$A_N_EXCLUDE" | tr -dc '\0' | wc -c | tr -d ' ') || TRACKED_GO=0
+    [ "$QUIET" = 1 ] || echo "   (A) denominator: handed to gofumpt = $TRACKED_GO; total tracked .go = $TRACKED_GO_ALL; subtracted by the one named exclusion $A_N_EXCLUDE = $((TRACKED_GO_ALL - TRACKED_GO))"
-    A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?
+    A_OUT=$(git ls-files -z '*.go' "$A_N_EXCLUDE" | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?
```
★形状＝**git 自己的 pathspec exclude magic**（减法由 git 做、⛔ 我的 `grep -v`）：`git ls-files -z '*.go' ':(exclude).scratch/**'` ⇒ **802**（1162−360，现跑）。四个前置探针（`logs\03`）：`:!` 短形同 **802**；指向不存在路径的排除 ⇒ **1162**（不静默变小）；`:(exclude)*` 一刀切 ⇒ **0** ⇒ 既有 `TRACKED_GO -gt 0` 空尺守卫**当场 rc=2**（不是绿）。

**(b) 三处分母文字**：`ci.yml:172-216`（那道 gofumpt 步：计数换成现量 1162/802/623 ＋ 写明"这些数会腐、引用前现推"；把 `:205` "census 步仍读全部 tracked `.go`" 换成 `:209` 起的 `WHAT MOVED ONE LEG LATER` 段；`:174` `Three claims:`→`Four claims:`）＋ `ci.yml:446-452`（`362-r1` 记的 `:436` 邻域那句）＋ `attrib.sh:677/679/680`（GREEN／RED 三行）。逐字见 §8。

## 4 ★改后：同一棵树、同一 shape ⇒ `rc=0`、**0 行**

```
== (A) tracked tree: git ls-files -z '*.go' ':(exclude).scratch/**' | xargs -0 .../gofumpt.exe -l
   (A) denominator: handed to gofumpt = 802; total tracked .go = 1162; subtracted by the one named exclusion :(exclude).scratch/** = 360
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty
attrib.sh: GREEN (tracked-only) - (A) empty over 802 tracked .go files (1162 tracked in total; the one named exclusion :(exclude).scratch/** subtracted 360)
== rc=0 ==
```
`TRACKED-DIRTY` 行数现量 ＝ **0**；树未被本步改动（前后同 `b39ac669`、`git status` 0 行）。

## 5 ★反向判据：撤掉具名排除 ⇒ 那 20 枚**逐名**重现

抄本 `D:\tmp\wisp\20261023\362r2\attrib-revert.sh`，与真件 `diff -u` **只 3 行**：赋值换成具名标记串 ＋ 两条 `git` 命令去掉 `"$A_N_EXCLUDE"` 实参（`logs\06` 逐字）。`sh -n` rc=0。

```
rc=1 ; (A) tracked-tree  lines=20 files=20 ; (A) denominator: handed to gofumpt = 1162; total tracked .go = 1162; ... subtracted ... = 0
```
`diff n-before2.txt n-revert2.txt` ⇒ **空 ⇒ IDENTICAL 20/20**（名单＝§1 那 20 行，一字不差、顺序同）。★这条同时是"排除面没顺手改别的东西"的判据：改前/改后只差这一个实参。

## 6 ★正控：产码目录里故意坏的 `.go` **必须被列、`rc=1`**

造在**仓外导出树**、三枚都 `git add` 成 tracked（ruler A 只认 index；`logs\08b` 逐字）：

| 枚 | 路径 | 期望 | 实测 |
|---|---|---|---|
| P1 | `internal/zz_362r2_posctl_unformatted.go`（两空格缩进复合字面量＝**故意坏**，产码目录） | **列它、rc=1** | ✔ **LISTED** |
| P2 | `cmd/zz_362r2_posctl_clean.go`（gofumpt-clean，产码目录）／P3 `.scratch/wisp/probes/999/posctl_dirty.go`（故意坏、在排除面内） | 两枚都不列 | ✔ 都不列 |

```
   (A) denominator: handed to gofumpt = 804; total tracked .go = 1165; subtracted by the one named exclusion ... = 361
   TRACKED-DIRTY	internal/zz_362r2_posctl_unformatted.go	...	<== CI IS RED ON THIS
attrib.sh: (A) tracked-tree  lines=1 files=1      == rc=1 ==
五条断言：PASS P1 被列 / PASS P2 未列 / PASS P3 未列 / PASS rc=1 / PASS 恰好 1 行
```
还原后复跑 ⇒ `tracked .go=1162`、porcelain 0 行、**`rc=0`、0 行**。（三枚是我刚造又自己撤的，形状同 `tracked-dirty-proof.sh:13-14` 的 "dirtied, measured, restored, measured again"。）
★**首发那段的 (c) 是我自己的尺错**：拿 `grep -F "/$P1"` 打整行、而行首是 `"   TRACKED-DIRTY<TAB>"` ⇒ 三枚全判"未列"。**仪器当时就是对的**；错在我的判据只看形状不看字段（`ruler-scope-checklist.md` 第 2 问）。改按 `awk -F'\t'` 取 `$2` 逐字比 ⇒ 全 PASS。

## 7 ★`--self-test` 仍绿（票 161 的载具、票 171 的 roster nails、票 275 的打印控）

```
   PASS  want rc=2+roster-names-the-sample  got rc=2 rows=1  throwaway tracked tree left at: /tmp/wisp-275-r2-rosterctl.AynU4o
   PASS  want ring=5 silent=3 named=3  got ring=5 silent=3  all named must-ring samples present
attrib.sh: --self-test cases=10 failures=0        attrib.sh: SELF-TEST GREEN ...        == rc=0 ==
```
★票 275 那格**顺带成了本笔的内建正控**：它把 `rosterctl/zz_roster_ctl_unformatted.go`（**`.scratch/` 之外**）造进一棵一次性 tracked 树，要求 gofumpt **rc=2** 且 stdout 有 `A-ROSTER` 行**点名它** ⇒ 排除面若吃掉产码形状路径，这格当场红。

## 8 ★三处注释改后逐字（`sed -n` 可复跑；`+` 行全文见 §3 的 `git diff`）

① `ci.yml:186-216`（那道 gofumpt 步）：`…and it is no longer true OF THIS STEP as written: the`／`…Measured at the 362-r2 anchor (2026-10-23, re-derived with`／`\`git ls-files\`, not carried over): of 1162 tracked .go files, 802 sit in`／`those five directories (internal 623, cmd 149, tools 18, scripts 11,`／`frontend 1), 360 sit under .scratch/ and 0 sit at the repo root…`／`…These counts rot as the tree grows: re-derive them before citing them (362-r1 wrote 1160/800/621 here, and by 362-r2 all three were already wrong).`；新增段 `ci.yml:209-216` 逐字：`WHAT MOVED ONE LEG LATER (ticket 362-r2) … It now reads the tracked set MINUS one named exclusion (.scratch/**), so the two steps agree on the bytes today (802 each) by two DIFFERENT definitions - a named INCLUSION here, a named EXCLUSION there - and they stop agreeing the moment a tracked .go falls inside one and outside the other. That divergence is a finding, not a rounding difference: when it happens the census step names the file and this step does not.`

② `ci.yml:446-452`（`362-r1` 记的 `:436` 邻域那句）逐字：`WHY THIS IS THE STEP THAT WAS MISSING, in the ticket's own words: the census`／`step above (\`gofmt (gofumpt) - the tracked set is the denominator\`) runs`／`\`attrib.sh --tracked-only\`, i.e. ONE of that instrument's two rulers. Since`／`ticket 362-r2 that ruler's denominator is the tracked set MINUS one named`／`exclusion (.scratch/**), not every tracked .go - so it is a ruler that can`／`go green on this tree at all, instead of one that is red on the same 20`／`deliberately-broken bench samples forever.`

③ `attrib.sh:680`（`362-r1` 点名的 `:637` "same denominator" 那句**已不成立**；行号现跑核过）：`attrib.sh: RED (tracked-only) - DENOMINATOR: this ruler and the CI gofmt step above agree today by two DIFFERENT definitions (a named EXCLUSION here, a named INCLUSION of the five production directories there), and they agree only while no tracked .go falls inside one and outside the other. When they stop agreeing, this is the ruler that names the file; that divergence is the finding, not a rounding difference. This one additionally refuses to answer when handed 0 files.`
（另 `:677` GREEN、`:679` RED 首行都带上三数分母；`attrib.sh:43-49` 与 `:548-576` 是这套排除的说明段。）

## 9 ★`sh -n` 与门禁 rc（每条首行先 unset 那几个 env；`logs\04`／`logs\12`）

| 项 | rc | 读数 |
|---|---|---|
| `sh -n attrib.sh` | **0** | ★shell 改完必验；`dash -n` 也 **0**；全件 CR＝**0** |
| `GOFLAGS= go build ./...` ／ `go vet ./...` | **0**／**0** | 两把都跑过 |
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（8 个 scope；**无 `.github/`、无 `.scratch/`** ⇒ 本笔结构上不可见） |
| `sh scripts/check-path-length-budget.sh` | **0** | 四数分开：分母 tracked **17612**／over-budget **61**／roster 覆盖 **61**／不在 roster **0**；墙区间 **0** |
| `sh attrib.sh --self-test` | **0** | §7 |
| `sh attrib.sh --tracked-only`（**甲**＝导出树）／（**乙**＝主检出） | **0**／**1** | 甲见 §4；乙 **3 行**，见 §10 |
| 格式门**甲**（导出树） | **0** | 五目录 ⇒ 0 行；`:(exclude).scratch/**` ⇒ 0 行；全 tracked 无排除 ⇒ **20 行、全 `.scratch`（20／0）** |
| 格式门**乙**（工作树，**只作诊断**） | 0 | 五目录 ⇒ **3 行**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` |
| `python -c "yaml.safe_load(ci.yml)"` | **0** | 6 个 job 全在；步名与 `run:` 行**一枚未动**（`:217`／`:221`／`:234`／`:285`／`:493`） |

## 10 ★★必须上报的读数：**乙尺（主检出）仍 `rc=1`，但那 3 行是票 361 那族的 CRLF 假红、⛔ 本笔的**

主检出跑 `attrib.sh --tracked-only`（不带 `ATTRIB_ROOT`）⇒ **`lines=3 files=3`**：`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`。三条判据逐条现跑（`logs\11`）：① 三枚**含 CR**（389／131／225），`tr -d '\r'` 后 `gofumpt -l` ⇒ **空** ⇒ **纯 CRLF 假红、字节本身干净**；② `git status --porcelain` 三枚**皆空**、`git diff` **0 行**、`git show HEAD:<p> | tr -cd '\r' | wc -c` ⇒ **0** ⇒ **blob 是 LF、脏的只是本机检出**；③ 同一条尺在**甲**（导出树／LF）上 ⇒ **0 行、rc=0**。

⇒ ★**`AC#3`「`attrib.sh --tracked-only` 必须 `rc=0`」在"入库字节"这条尺上取到了（甲＝0）；在我这张 Windows 主检出上取不到，而取不到的原因与 `.scratch` 台件无关**（同样三枚 `362-r1` 件 §10 已记，台账 `A1161` 同族）。★这也证明 `attrib.sh` 头部 "form (A) reads tracked bytes" 在本机**不成立**：`git ls-files` 给的是**路径**、`gofumpt` 开的是**盘上字节** ⇒ 有 `core.autocrlf` 的检出上它量的是工作树。**要不要断言化（改读 blob／对这三枚做 eol 归一）＝⛔ 本腿写面、请裁。** 同一把乙尺改前改后：**34 行 → 3 行**（34 ＝ 31 `.scratch` ＋ 3 CRLF；31 比甲尺的 20 多，因 `.scratch` 里另有 11 枚台件盘上 CRLF、blob 是 LF —— 正是"本机 ≠ CI"那一族）。

## 11 ★未做／须裁

- **未 push、未跑 CI** ⇒ ⛔ 不自称"CI 绿"；所有 rc 都是本机现跑。**`go test -count=1` 未跑**：本腿 **`.go` 零字节**（§3），且主检出此刻是别的腿的写窗（`internal/memory/**`、`internal/winsec/**` 在场未提交）⇒ 现在取的名册量的是别人半成品；要这格读数请明示。
- ★**`attrib.sh` 的 `summary()` 与 `ci.yml:234` 步名未动**（"the tracked set is the denominator"）：它读 index、⛔ 工作树，这层对立**仍成立**；缩窄已在 §8① 写明。改步名＝改 yaml 的报告身份，⛔ 本腿"只注释"的写面；若要一致请裁。
- ★**`ci.yml:236-238`**（"the step above and this one read the same bytes and must agree"）**仍成立**：五目录 **802** ＝ tracked−`.scratch` **802**，且"两边都不在"的 tracked `.go` ＝ **0**（`logs\03` Q5b）⇒ 不是假话，故未动。本件未拆 `-part`（≤160 行）。
