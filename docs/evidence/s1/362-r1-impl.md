# 票 362 落地腿 `362-r1` 交件 —— `ci.yml` 的 gofumpt 步改按目录集传参（形〔甲′〕：改调用面、⛔ 洗任何字节）

- 腿＝`362-r1`（**写腿／实现者**）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
- 写面实动：`.github/workflows/ci.yml`（**只那一处调用面 ＋ 那步注释**）＋ 本件。`internal/`／`cmd/`／`scripts/`／`tools/` **零字节**。
- 台面：起手锚 `51bb042b`（导出树 1）／`8d3ad434`（交件时 HEAD，导出树 2）。★两树之间 `.go` 改动 **0 枚**（`git diff --name-only 76ae3bbb..HEAD -- '*.go'` 空）⇒ 下面 20 枚在两根树上逐枚同名。
- 环境：`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；`gofumpt` 一律**绝对路径** `/d/work/base/gopath/bin/gofumpt.exe`（`v0.12.0 (go1.27.1)`）。

## 0 ★机主批准（逐字）

> 「好的，这三件事我都允了，同意，你可以继续了」

出处＝台账 `docs/reports/pending-and-issues.md` **`A1161`①**（逐字同上）；该条批准范围第②件＝「**票 362** 改 **`.github/**`**」。派单存档＝`.scratch/wisp/probes/orch/20261023/txt/dispatch-three-legs-after-approval.md` §②（提交 `76ae3bbb`）。

## 1 ★改了什么（逐行 diff）

`git diff --stat -- .github` ⇒ **`1 file changed, 35 insertions(+), 1 deletion(-)`**。**唯一的代码行改动**：

```diff
-          OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"
+          OUT="$(gofumpt -l internal cmd tools scripts frontend)"
```

另 **+34 行纯注释**（那步原本**没有注释**）。行号现跑复认（⛔ 信票面）：改前步名 `:172`／调用面 `:176`；改后注释 `:172-205`／步名 `:206`／调用面 `:210`。逐字见 §6。

## 2 ★开工前：票面前提逐条现跑（⛔ 信票面）

| 票面断言 | 我现跑 | 判 |
|---|---|---|
| `git ls-files '*.go'` ＝ 1156 | **1160** | ★漂 +4，且**先于派单**（`76ae3bbb` 已是 1160） |
| 产码只有五目录、根目录无 `.go` | internal 621／cmd 149／tools 18／scripts 11／frontend 1 ＝ **800**；`.scratch` **360**；根 **0**；**其它首段 0 行** | ✔ 800+360=1160 **无漏类**（停手①不触发） |
| 旧形 ⇒ 20 枚、非 `.scratch` 0 | 20／**0**（§3） | ✔ |
| 新形 ⇒ 0 枚 | **0**（§4） | ✔ |
| 那步在 `ci.yml:172` | 步名 `:172`、**调用面在 `:176`** | ✔ 现认 |
| `.github` 里有第二处 `gofumpt -l .` | `gofumpt -l` 共 **4 行**＝**1 处真调用**（原 `:176`）＋3 行注释 | ✔ **停手②字面不触发**（语义那半见 §9①） |

## 3 ★甲尺：仓外导出树，改前 20 行逐字

venue＝`git archive HEAD | tar -x` ⇒ **17607 件**、`.go` 1160、**`*.go` 的 CR 计数＝0**（`tr -cd "\r" | wc -c`）、**无 `.git`**。
shape＝`OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"`（**只取 stdout**，与那步逐字同形；stderr 另记）。
reading＝**rc=0、stdout 20 行、唯一路径 20、stderr 0 行**；归属两把尺（`grep -vc '^\.scratch[\\/]'`／`tr '\\' '/' | grep -vc '^\.scratch/'`）**都得 0**。

```
.scratch\wisp\probes\163\a1\main.go
.scratch\wisp\probes\174\c2\zz174c2_wiring_pair_windows_test.go
.scratch\wisp\probes\183\r2\mut\task-boxset-off.go
.scratch\wisp\probes\185\r1\mut-m1\hostpath_185.go
.scratch\wisp\probes\197\r1c\pre\subagent_197.go
.scratch\wisp\probes\197\r1c\pre\subagent_197_test.go
.scratch\wisp\probes\212\v1\mut\main-noq9.go
.scratch\wisp\probes\220\r1\prechange\prechange220_probe_test.go
.scratch\wisp\probes\222\v1\mutations\m11-budget-50ms-gated\subagent_222_test.go
.scratch\wisp\probes\224\v2\dialect_probe_test.go
.scratch\wisp\probes\231\v1\mutation\config_reload_branch_after_prefix.go
.scratch\wisp\probes\235\v1\mut\prereadoff-m3_222_test.go
.scratch\wisp\probes\241\v1\posctl\badly_formatted.go
.scratch\wisp\probes\241\v1\probe_v1_readings_test.go
.scratch\wisp\probes\253\v2\denoms253v2.go
.scratch\wisp\probes\259\r1\gofumpt-negctl\probe.go
.scratch\wisp\probes\263\v1\src\main.go
.scratch\wisp\probes\33\p1\q1\main.go
.scratch\wisp\probes\33\p1\q2\main.go
.scratch\wisp\probes\33\p1\q3\main.go
```

## 4 ★改后 0 行（同一树；命令**从文件里抠出来**再跑）

- 同一树（`51bb042b`）跑新形 ⇒ **rc=0、0 行、stderr 0 行** ⇒ `STEP WOULD PASS`。
- **从 `ci.yml:210` 抠出参数串**（`-l internal cmd tools scripts frontend`）用绝对路径跑**第二棵 pristine 导出树**（`8d3ad434`，17607 件、CR=0）⇒ **rc=0、0 行**。
- ★**自查：我这把尺错过一次、已改**。第一版 `OUT="$($CMD)"` 未加引号 ⇒ 词分拆后跑的是**裸 `gofumpt`**（本机 rc=127）、stderr 进 `/dev/null` ⇒ **空输出被读成"0 行、PASS"**。改成绝对路径＋显式 rc 检查（rc=127 即拒答）后，同一发读成 **1 行**（§5）。**"量不到"⛔ 等于"0"。**

## 5 ★"产码目录里故意坏枚"被列出 ＋ ★反向判据

造在**仓外导出树内**（⛔ 仓内）：`internal/zz_362r1_planted_unformatted.go`（两空格复合字面量缩进＝本仓 `attrib.sh` 已量的 TRACKED-DIRTY 形）＋镜像**干净**对照 `internal/zz_362r1_planted_clean.go`。

- **新形必须列它**：new ⇒ **`internal\zz_362r1_planted_unformatted.go`、1 行**；**干净那枚 ⛔ 被列**（两向都有读数）。
- 旧形同树 ⇒ **21 行**；`旧−新` 差集 ⇒ **恰那 20 枚台件**（`comm -23`，逐行在 `logs/05-reverse-and-plant.txt`）。
- **反向判据**：改回 `.` ⇒ **20 行重新出现**，两棵**独立**导出树各一发、非 `.scratch` 皆 0。
- ★**射程**：这条证明"新形的分母里若有坏字节，它**会响**"；⛔ 证明"CI 会全绿"（见 §9①）。

## 6 ★注释改后逐字（`:172-205`，含"分母语义变了"那段）

```
      # CALL SURFACE CHANGED (ticket 362, AC#1 form (A'), owner-approved edit of
      # .github/**; the approval is docs/reports/pending-and-issues.md A1161 (1),
      # and the dispatch that carries it is commit 76ae3bbb). Three claims:
      #
      # WHY. This step used to read `gofumpt -l . tools/d22scan tools/mockllm`,
      # which walks the WHOLE tree - and `.scratch/**` is the directory this repo
      # fills with DELIBERATELY broken samples. At the anchor 20 TRACKED .go files
      # under `.scratch/wisp/probes/**` are not gofumpt-clean on purpose (two exist
      # only as negative controls: a file named badly_formatted.go, and
      # gofumpt-negctl/probe.go). Formatting them would destroy what they were
      # created to prove, so the old criterion was satisfiable only by deleting or
      # laundering somebody else's evidence. It now reads
      # `gofumpt -l internal cmd tools scripts frontend`.
      #
      # WHAT THAT DOES TO THE DENOMINATOR. "The denominator of the formatting gate
      # is the TRACKED set, not the working tree" is the claim the census step
      # below is built on, and it is no longer true OF THIS STEP: the denominator
      # here is now the PRODUCTION DIRECTORY SET named above, not every tracked
      # .go. Measured at the anchor: of 1160 tracked .go files, 800 sit in those
      # five directories (internal 621, cmd 149, tools 18, scripts 11, frontend 1),
      # 360 sit under .scratch/ and 0 sit at the repo root, and no tracked .go
      # file's first path segment is anything else - so the five directories are
      # exactly "production" and .scratch/ is exactly what leaves this
      # denominator. The reading is still taken over tracked bytes on a checkout;
      # what narrowed is which tracked bytes are counted.
      #
      # WHAT IT GIVES UP, stated rather than left to be discovered: a bench sample
      # that somebody one day places under one of those five directories WOULD be
      # scanned here and would redden this step. That count is 0 today - all 20
      # dirty tracked files are under .scratch/ - and it should stay 0, because
      # those directories are where the gate has to bite: a sample parked there is
      # either a mistake or belongs in .scratch/. Unchanged: no -w, the first
      # non-empty output still exits 1, and the census step below is a separate
      # ruler that still reads every tracked .go - this edit does not move it.
```

CR＝0（工作树与 HEAD blob 都是 0）；非 ASCII 行数**改前 14＝改后 14**（我这 35 行**纯 ASCII**，与原文件英文风格一致）；全文件最长行 176 字符（` :358`，非本笔）。

## 7 ★门禁（自己现跑；每个脚本首行都先 unset 那几个 env）

- `sh scripts/d22scan.sh` ⇒ **rc=1**，1 条：`internal/risk/mode.go:178: [phantom-citation] comment cites repo path "tools/gate.go" which does not exist on disk`。
  **⛔ 不是本腿的**：该文件 ` M`（未提交、mtime 09:48:01）、`(tools/gate.go:43-47)` 是 `git diff` 的 **`+` 行**、**HEAD blob 无此句**；同刻 `internal/tools/bridge.go` 也 ` M`、另有两枚新 `ticket363_*_test.go` ⇒ ★**`363-r1` 的写窗**。
  **隔离读数**：同一扫描器打**仓外导出树（HEAD 字节）** ⇒ **rc=0 `clean`**；它自报 8 个 scope 里**没有 `.github/`** ⇒ 本腿的 `ci.yml` 对 d22scan **结构上不可见**。前置正控（脚本第 1 步）**PASS=35 FAIL=0 SKIP=0／RUN=77**（先证它红得起来）。
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0 GREEN**，**四数分开**：**分母 tracked paths=17606**／**over-budget=61**／**roster 覆盖=61**／**不在 roster=0**。band：122..180 中段 **57**、过 180 **4**、**墙区间 (206,217]＝0**；longest 185 字符、worst full path 229（帽 165）。
- `go build ./...` ⇒ **rc=0**。

## 8 ★红名册／★逐枚提交越界＋正控

- **红名册（本腿）＝（无）、新增红＝0**。读数形状：本腿逐笔 `git show --name-only` **零 `.go`** ⇒ 结构上无法新增红；`go build ./...` rc=0 作编译面读数。**整包 `go test` 未取**，理由具名＝主检出此刻是 `363-r1` 的写窗（§7），现在取出的名册量的是它的半成品；按"突变窗口⛔ 混入读数"不能记本腿账（已列 §10）。
- 提交逐枚（`git show --name-only --format='%h %s'`；三笔各**单路径**）：

| # | 提交 | 文件 | 越界 |
|---|---|---|---|
| ① | `2dc05b81` | `.github/workflows/ci.yml` | **0** |
| ② | `8ea04e82` | `docs/evidence/s1/362-r1-impl.md`（新建） | **0** |
| ③… | 本件后续每一笔（号见交件消息；尺件 `logs/11-commit-scope.txt`） | 同上**单路径** | **0** |

- ★**行数读数自查（我又错一次、已改）**：第二笔提交消息写"121 行"＝ PowerShell `Measure-Object -Line` 的读数（**它跳过空行**）；`wc -l`＝**153**、`git show --numstat` 同 ⇒ **本件 153 行**，≤160。

- **正控（同一把越界尺）**：`a782196e`（真含 1 枚 `.go`＝`internal/winsec/leafrace_342r2_windows_test.go`）⇒ 尺**报出该枚越界**；空控 `fccf31a9`／`15d00c78`（零 `.go`）⇒ **0** ⇒ 尺不是恒 0。
- ★纪律现记：我 `add` 后 `git diff --cached --name-only` 有 **2 行**——索引里另有一枚**别人**已 staged 的 `.scratch/wisp/issues/164-…-done.md`。两处都用**显式 pathspec**，`git show --name-only HEAD` 逐字**只列 `ci.yml`**，那枚**没被带进**、也仍在索引里（未被我动）。

## 9 ★★两条须编排者裁的问题（本腿⛔ 自行假设、⛔ 越写面）

① ★**本腿只把两道格式门里的一道变绿。** 紧随其后的 census 步（161 AC#7 form A；改后步名 `:223`、`run:` 在 `:274`）跑 `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`，其 ruler A 逐字是 `git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l`（`attrib.sh:540-544`）＝**全 tracked 集** ⇒ 那 20 枚仍在它分母里 ⇒ **它今天仍会以同样 20 个名字 rc=1**；`attrib.sh:637` 那句 "The CI gofmt step and this step read the same denominator" 随本笔**失真**。后果：票 `AC#3` 的「`attrib.sh --tracked-only` 必须转绿（rc=0）」在形〔甲′〕下**取不到**——要它绿只能①洗字节（本腿硬禁）②改 `attrib.sh`（⛔ 写面外）。**要不要一起动，请裁。**
② 那句原话「the denominator of the formatting gate is the TRACKED set, not the working tree」**逐字活在另一道步的注释里**（改前 `:194-195`／改后 `:228-229`），**不在**派单点名的"那步"（原 `:172`，那步原本无注释）。我按写面只改了我这步的注释、**没动**它；要不要同笔改、以及要不要改 `attrib.sh:637` 那句，请裁。

## 10 ★未测／未做

- **CI 真跑读数**：⛔ push ⇒ 无 run，本腿**不自称"CI 绿"**。
- **整包 `go test` 红名册**：未取（§8 理由）。**甲乙两尺都报**，**乙只作诊断**：工作树旧形 **49 行**、新形 **3 行**＝`cmd\wisp\models.go`／`internal\agent\tools.go`／`internal\agent\approval\pending_read.go`——三枚都是**工作树 CRLF 假象**（导出树同尺非 `.scratch` ＝ 0；台账 `A1161` 记的也是这同一族"乙 3 行纯 CRLF 假象"）。
- **`scripts/slo-freshness.sh` 未跑**（要 GitHub API）；只核它那把**文本尺**不被本笔移动：`^[[:space:]]+if:` **改前 37＝改后 37**（HEAD blob／工作树各一发）。★它 `:196-199` 那串行号（`:178…:473`、`slo-full` 写作 `:537`）**在本腿之前就已腐坏**（现量 `slo-full:` 在 `:979`、改前应为 `:945`）⇒ 本笔 +34 让它更旧；修它属 `scripts/**`＝⛔ 我的写面。

★**读数件全在** `D:\tmp\wisp\20261023\362r1\logs\`（`01-census`／`02-export`／`03-before-old-shape`／`04-after-new-shape`／`05-reverse-and-plant`／`06b-file-shape-fixed`／`07-gates`／`08-isolated-d22scan`／`09-loose-ends`／`10-second-tree`；尺件在 `sh\`）。
