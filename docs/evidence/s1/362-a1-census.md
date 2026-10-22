# 票 362 普查腿 `362-a1` 交件（`AC#0` 只读现读＋归因）

- 腿＝`362-a1`（只读腿：⛔ 落 `.go`、⛔ 改票面／`Status:`、⛔ 翻框、⛔ push、⛔ 改 `.github/**`）
- 时刻 2026-10-22 19:19–20:1x +08｜起手 HEAD `dac97d70`，收件 HEAD **`095699f2`**（本腿期间 HEAD 动了 4 笔，见 §7）
- 尺：`gofmt`＝`go1.27.1 windows/amd64`｜`gofumpt v0.12.0 (go1.27.1)`｜`GOFLAGS=`、`unset GOTMPDIR TMPDIR TEMP TMP`、⛔ 设 `CGO_ENABLED`
- 读数原件（仓外）：`D:\tmp\wisp\20261022\362a1\**`（`per-file-archive.tsv`／`birth.tsv`／`stability.tsv`／`t1..t3-archive-*.txt`／`ci172-*.txt`／`d22scan2.txt`／`pathlen2.txt` 等）
- **结论：票面全部前提成立、四条必做全部量到；★没有触发任何停手条件。**

## 1 前提新鲜度（逐条现跑，对 **HEAD blob** 复量）

| 票面断言 | 我的读数（HEAD `095699f2`） | 判 |
|---|---|---|
| 20 枚 `.go` 入库字节本身脏 | `git ls-files '*.go'`＝**1156**，逐枚 `git show HEAD:<p> \| gofumpt -l` ⇒ **20 枚 `DIRTY_BLOB`** | ✅ |
| 仓外 `git archive` 导出树跑 `--tracked-only` ⇒ rc=1、lines=20 files=20 | 导出树（17577 件、`git status --porcelain`＝**0 行**、20 枚 CR 合计 **0**）⇒ `rc=1`、`lines=20 files=20` | ✅ |
| 与 HEAD blob 逐枚 sha256 相同（20/20） | 逐枚 `sha256sum` ⇒ **20/20 IDENTICAL** | ✅ |
| 与已推送 `626fcc7b` 的 blob 逐枚相同（20/20） | **20/20 IDENTICAL**；且 `git merge-base --is-ancestor 626fcc7b origin/dev` ⇒ **YES**（确在已推侧） | ✅ |
| 其中 15 枚连 `gofmt -l` 都判脏 | **15 DIRTY ＋ 5 CLEAN**（`comm -13` 空 ⇒ gofmt 名单是 gofumpt 名单的真子集） | ✅ |
| `ci.yml:172` 那步在同一树同样 20 行 | **rc=1、lines=20、名字逐枚 `diff` IDENTICAL** | ✅ |
| `360-v1` 补的：20 枚 `not_probes=0`、全在 `.scratch/wisp/probes/**` | `probes=20`／`not_probes=0`；16 个 `probes/<NN>/` 各自 `issues/<NN>-*.md` **全部存在（16/16）** | ✅ |
| 票面 `ci.yml:172` 逐字＝`gofumpt -l . tools/d22scan tools/mockllm` | `git show HEAD:.github/workflows/ci.yml \| sed -n '172,181p'` 逐字同（step 名 `gofmt (gofumpt)`、`if: ${{ !cancelled() }}`） | ✅ |
| 来路两笔（`5182cb74`／`0c85c503`）在盘 | `5182cb74` 2026-10-22 16:48（361-a1 两件）／`0c85c503` 19:09（360-v1 两件），四件均 `git cat-file -s HEAD:` 可读 | ✅ |

## 2 ★20 枚逐名 ＋ 逐枚三列

判据：在**仓外导出树**（＝HEAD blob 字节）上 `gofmt -l <p>`／`gofumpt -l <p>` 逐枚单跑；「属哪个 probes」＝路径第 1 段。

| # | 路径（仓库相对） | `gofmt -l` | `gofumpt -l` | probes |
|---|---|---|---|---|
| 1 | `.scratch/wisp/probes/33/p1/q1/main.go` | **DIRTY** | **DIRTY** | 33 |
| 2 | `.scratch/wisp/probes/33/p1/q2/main.go` | **DIRTY** | **DIRTY** | 33 |
| 3 | `.scratch/wisp/probes/33/p1/q3/main.go` | **DIRTY** | **DIRTY** | 33 |
| 4 | `.scratch/wisp/probes/163/a1/main.go` | **DIRTY** | **DIRTY** | 163 |
| 5 | `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | CLEAN | **DIRTY** | 174 |
| 6 | `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | **DIRTY** | **DIRTY** | 183 |
| 7 | `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | **DIRTY** | **DIRTY** | 185 |
| 8 | `.scratch/wisp/probes/197/r1c/pre/subagent_197.go` | CLEAN | **DIRTY** | 197 |
| 9 | `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go` | CLEAN | **DIRTY** | 197 |
| 10 | `.scratch/wisp/probes/212/v1/mut/main-noq9.go` | **DIRTY** | **DIRTY** | 212 |
| 11 | `.scratch/wisp/probes/220/r1/prechange/prechange220_probe_test.go` | CLEAN | **DIRTY** | 220 |
| 12 | `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | **DIRTY** | **DIRTY** | 222 |
| 13 | `.scratch/wisp/probes/224/v2/dialect_probe_test.go` | **DIRTY** | **DIRTY** | 224 |
| 14 | `.scratch/wisp/probes/231/v1/mutation/config_reload_branch_after_prefix.go` | **DIRTY** | **DIRTY** | 231 |
| 15 | `.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go` | **DIRTY** | **DIRTY** | 235 |
| 16 | `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go` | **DIRTY** | **DIRTY** | 241 |
| 17 | `.scratch/wisp/probes/241/v1/probe_v1_readings_test.go` | **DIRTY** | **DIRTY** | 241 |
| 18 | `.scratch/wisp/probes/253/v2/denoms253v2.go` | **DIRTY** | **DIRTY** | 253 |
| 19 | `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` | **DIRTY** | **DIRTY** | 259 |
| 20 | `.scratch/wisp/probes/263/v1/src/main.go` | CLEAN | **DIRTY** | 263 |

四数（⛔ 相减）：`gofumpt DIRTY = 20`／`gofmt DIRTY = 15`／`gofmt CLEAN = 5`／`not_probes = 0`。16 个票号、20 枚件。
★**那 5 枚 `gofumpt`-only 的规则逐枚现读**（`git show HEAD:<p> | gofumpt -d`）：
`174` 复合字面量尾逗号＋缩进；`197/subagent_197.go` 与 `197/subagent_197_test.go` 复合字面量换行；`220` 嵌套复合字面量缩进；`263` **`0644` → `0o644`**。
★另核（不属本票判据、只作诊断）：20 枚在**主检出**逐枚 `git status --porcelain`＝**CLEAN**、`git ls-files --eol`＝**`i/lf w/lf attr/text eol=lf`**、工作树 `CR=0` ⇒ **⛔ 是 CRLF 假红**，是真脏。

## 3 ★正控（硬闸，双向）—— 造在**仓外导出树**，⛔ 仓内

树：`D:\tmp\wisp\20261022\362a1\archive2`（`git archive HEAD` 导出＋自建 `.git`／index，与 CI 同形）。两枚均纯 LF（`CR=0`）、`git add` 后 tracked。

| 枚 | 内容 | `gofmt -l` | `gofumpt -l` | 判 |
|---|---|---|---|---|
| `internal/census362a1/known_clean.go` | 正常缩进、无 `0644`、单 import | **不列** | **不列** | ✅ 尺不恒脏 |
| `internal/census362a1/known_dirty.go` | 错缩进 `x := 1` ＋ 未排序 import | **列** | **列** | ✅ 尺不恒绿 |

三发读数（`ATTRIB_ROOT` 指该树，见 §8 陷阱）：

```
(1) 基线           tracked .go=1156  rc=1  (A) lines=20 files=20  TRACKED-DIRTY 行=20
(2) ＋两枚正控     tracked .go=1158  rc=1  (A) lines=21 files=21  TRACKED-DIRTY 行=21
                   ★点名 known_dirty.go=1、known_clean.go=0
                   ci.yml:172 形状同发：rc=1 lines=21，dirty=1 clean=0
(3) 撤两枚正控     tracked .go=1156  rc=1  (A) lines=20 files=20  TRACKED-DIRTY 行=20（回基线）
```

⇒ **两把尺都活着、双向都被钉住**：干净枚两尺都不列且 ⛔ 被点名，脏枚两尺都列且**逐名点名**；`20 → 21 → 20` 可逆。

## 4 归因：这 20 枚是怎么进来的

`git log --diff-filter=A --format='%H' -- <p>` 逐枚取**首次提交**，再 `git show <c>:<p> | gofumpt -l` 判**入库当时**是否已脏：

- ★**`DIRTY_AT_BIRTH = 20 / 20`、`CLEAN_AT_BIRTH = 0`** ⇒ **一枚都不是后来被弄脏的，全是"生下来就脏"**。
- ★**20 枚逐枚 `git log -- <p>` 的提交数＝1**（`1 × 20/20`）⇒ 自入库起**再没被任何提交碰过**。
- ★**20 枚逐枚 `git show --name-status -M <c> -- <p>` 全为 `A`**（纯新增，**无一枚是 rename 带进来的**）。
- **16 笔**首次提交（20 枚散在 16 笔里）；**集中在少数几笔的只有三笔**：

| 笔数 | 提交 | 日期 | 枚 | 主题（截断） |
|---|---|---|---|---|
| **3** | **`7a956604`** | 10-02 | `probes/33/p1/q1..q3` | `probe(33-p1) 第二轮三问真跑交件…q1/q2/q3三枚探针+14份日志+census五节` |
| **2** | **`c7abc84b`** | 09-28 | `probes/197/r1c/pre/` 两枚 | `ledger(A415 收 197-r3b：载体层入库 10 枚文件…)` |
| **2** | **`c1e1a92e`** | 09-30 | `probes/241/v1/` 两枚 | `evidence(241-v1) §3 AC#3 成立…＋突变台件落盘` |
| 1 | `e13196bd` | 09-27 | `probes/163/a1/main.go` | `163-a1 只读设计核裁定表…台件 probes/163/a1` |
| 1 | `6479d540` | 09-28 | `probes/174/c2/…` | `docs(evidence): 174-c2 wiring cost census + read-only before/after probe` |
| 1 | `78ea1d19` | 09-28 | `probes/183/r2/mut/…` | `183-r2 交件表＋票面 Progress log 追加＋本程读数落盘（零产码改动）` |
| 1 | `ee1e2118` | 09-28 | `probes/185/r1/mut-m1/…` | `185-r1: 分页续读——豁免改绑…` |
| 1 | `f23586fa` | 09-29 | `probes/222/v1/mutations/…` | `222-v1 台件与原始读数落盘（只建不删）…` |
| 1 | `48e09742` | 09-29 | `probes/235/v1/mut/…` | `evidence(235-v1): AC#1 裁完＝成立…` |
| 1 | `961e23b0` | 09-30 | `probes/224/v2/…` | `probes(224-v2) 台件留档：N#5 的 overlay 方言形态探针…` |
| 1 | `54742783` | 10-03 | `probes/212/v1/mut/…` | `212-v1 第一实例验收交件（非实现者…）` |
| 1 | `480b970d` | 10-03 | `probes/220/r1/prechange/…` | `220-r1: 名册那一格改成 corr 与 taskID 双查…` |
| 1 | `8c25bfab` | 10-04 | `probes/263/v1/src/main.go` | `263-v1 交件：验收表九节写满…` |
| 1 | `50e0299c` | 10-05 | `probes/259/r1/gofumpt-negctl/probe.go` | `票 259 腿 259-r1 第 3 步：证据件 §2 到 §9 写满…` |
| 1 | `fdb437c3` | 10-06 | `probes/231/v1/mutation/…` | `probes(231-v1 遗产代提)：验收腿撞 150 轮帽、29 枚台件…编排者代提` |
| 1 | `41040e89` | 10-14 | `probes/253/v2/denoms253v2.go` | `ticket 253 (253-v2): AC#2 非实现者验收裁决…` |

⇒ **成因一句话**：这些是**台件（deliberately-broken 样本）**，被各腿**原样入库**——`probes/241/v1/posctl/badly_formatted.go` 与 `probes/259/r1/gofumpt-negctl/probe.go` 两枚**名字本身就是"故意不格式化"**。入库时**没有任何一步对 `.scratch/**` 跑过 `gofumpt -w`**，`ci.yml:172` 那步也从未在 `.scratch` 上变红过（本票之前无人量）。
★时间跨度 09-27 → 10-14（**18 天、16 笔、13 条不同腿**）⇒ **不是一次事故，是一条长期没有闸的通道**。

## 5 ★`ci.yml:172` 那步今天会不会红 —— **会**

在**仓外导出树**上按 `ci.yml:174-181` 的形状逐字跑（`gofumpt -l . tools/d22scan tools/mockllm`，非 `--tracked-only`）：

```
gofumpt would reformat:      <- 20 行名字
STEP_RC=1   LINES=20
tracked .go handed to the walk = 1156
```

- **rc=1、lines=20**（两棵独立导出树各跑一发，`archive2` 与当前 HEAD 的 `archive3`，**都是 rc=1/20**）。
- ★那 20 个名字与 §2 名册 **`diff` 逐字 IDENTICAL** ⇒ CI 的 gofumpt 步与 `attrib.sh --tracked-only` **在同一批字节上同时红、同一批名字**。
- ★导出树 `git status --porcelain`＝**0 行**、20 枚 `CR=0` ⇒ **CI 侧没有 CRLF 可以赖**；这是**真红**，⛔ 检出副本假红。
- 归因成立 ⇒ **⛔ 触发停手条件②**。

## 6 门禁（本腿现跑，`GOFLAGS=`／已 unset 四个 temp 变量）

| 门 | 命令 | rc | 末行逐字 |
|---|---|---|---|
| D22 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; …` |
| 路径预算 | `sh scripts/check-path-length-budget.sh` | **0** | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

另：`denominator: tracked paths=17577  over-budget=61  covered by roster=61  not in roster=0`。

## 7 纪律格

- **本腿新增 `.go` ＝ 0**（正控两枚造在仓外导出树）。三把尺同判：
  ①`git status --porcelain -- '*.go' | grep -c census362a1` ⇒ **0**（正控名不进本仓）；
  ②本仓 `find . -name '*.go' -newermt '2026-10-22 19:19'`（本腿起手时刻）⇒ **1 枚**，是 **`internal/config/writeguard_adopt_360_test.go`**，属**并发写腿 `360-r2`**（`095699f2`），⛔ 本腿；
  ③`.go` 未跟踪 77 枚（未忽略）＋325 枚（已忽略）＝**402**，1156＋402＝**1558**＝盘上遍历数 ⇒ 无一枚是本腿造的
- 写面：`docs/evidence/s1/362-a1-census.md`（本件）＋仓外 `D:\tmp\wisp\20261022\362a1\**`；⛔ 改 `internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`
- ⛔ 改票面、⛔ 改 `Status:`（HEAD 仍 `ready-for-agent`）、⛔ 翻任何框、⛔ push（`origin/dev..HEAD`＝404，全在本地）
- `add`／`commit` 两处**显式 pathspec**；⛔ `add -A`／`.`、⛔ amend／reset／rebase／stash／checkout ./clean
- ★本腿期间 HEAD 动了 4 笔（`dac97d70`→`d6adbd8e`→`ceb9e64c`→**`095699f2`** 360-r2）⇒ **§1/§5 的两把尺在最终 HEAD `095699f2` 上重跑过：仍是 20 枚、名字逐字同**（⛔ 漂移）

## 8 ★两条附注（下一腿会踩，照实记）

1. ★★**`attrib.sh` 的 root 取自脚本自身目录（`:137` `d=${ATTRIB_ROOT:-$here}`），⛔ 取自 `cwd`**。
   所以「`cd` 进导出树再跑 `sh <主检出>/…/attrib.sh`」**实际量的是主检出**（我第一次就这么跑，得到 **35**，那是主检出读数，⛔ 导出树读数）。
   正确形状：`ATTRIB_ROOT=<导出树> sh <attrib.sh> --tracked-only` ⇒ **20**。`361-a1` 报的 20 是对的，但**这个陷阱它没写下来**。
2. ★**乙尺（工作树）跟踪集读数＝35**（`git ls-files '*.go' | xargs gofumpt -l`），⛔ 55／⛔ 50。
   本机 `core.autocrlf=true` ＋ `.gitattributes:4 *.go text eol=lf`；同一树文件系统遍历＝**1558** 枚 `.go`
   ＝ 跟踪 **1156** ＋ 未跟踪未忽略 **77** ＋ 未跟踪已忽略 **325** ⇒ **遍历口径与跟踪口径必须分开报**（`A1134⑦a` 复现）。

## 9 一句话

★**仓库里真有 20 枚入库字节就没格式化的 `.go`，全在 `.scratch/wisp/probes/**` 的 16 个票的台件目录里，15 枚 `gofmt` 也判脏、5 枚 `gofumpt`-only，20 枚全部"生下来就脏、之后没被碰过"（16 笔首次提交、横跨 09-27→10-14），而 `ci.yml:172` 那步在纯 LF、无未跟踪状态的导出树上 rc=1、20 行、名字逐字同 —— 这是真债、⛔ 本机/CI 不同形。**
