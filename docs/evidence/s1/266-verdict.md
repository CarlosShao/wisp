# 票 266 裁决（非实现者 / 只读腿）— 外牙 `scripts/slo-nail-census.sh` ＋ lint 作业的 CI 步

- 裁决腿：`266-verdict` ＝ **非实现者**（实现腿 `266-r1`/`266-r2`/`266-ci`，普查腿 `266-a1`）。本腿**只读**，唯一写入＝本件。
- 时刻：2026-11-05 20:43–20:53 +08。
- 现跑 HEAD：`34637ae6c`（派单写 `7eb64d5ce`）。漂移＝**两笔票 51 的 commit**：`git diff --stat 7eb64d5ce..34637ae6c` ＝ 3 件、全在票 51 面
  （`51-tier2-goja.md`→`-done` 改名、`docs/evidence/s1/51-verdict.md`、台账 5 行）；`-- scripts tools .github` **0 行** ⇒ **266 交件面未漂**。
  266 的五枚交件 commit（`6589a0625` 脚本 / `78e77888f`+`38d05fcf4`+`0fa1527f9` 证据 / `6e5957756` CI 步）**逐枚**现跑
  `git merge-base --is-ancestor <sha> 7eb64d5ce` ＝ **YES ×5**。
- 环境：Git Bash；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；⛔ 全程**未** `CGO_ENABLED=0`。
- 判据：**盘上现跑**，不采信证据件自述；证据件只当「尺在哪」的索引。⛔ 未改代码、未改工单、未改别腿证据件。
- ⛔ 本腿**没有**执行任何写树动作：变异全在 `mktemp -d` 台件上（本尺自带的 `--with-self-test` 也是台件），工作树一字节未动。

## 1. 现跑总读数（本件所有结论的底数）

| # | 命令 | 读数（取数时刻） |
|---|---|---|
| 1 | `grep -c '"scripts' tools/d22scan/main.go` | **0**，`rc=1`（20:44）⇒ `scripts/**` 不在 d22scan 射程 |
| 2 | 两枚钉真身现读（`grep -n`，20:44） | 能力钉定义 `scripts/slo-check.ps1:214`／ok 句 `:230`；词面钉定义 `:233`／ok 句 `:258`；调用 `:291`（shape）／`:292`（instrument）；首枚采样锚 `:533` ⇒ **与 r2 §1（18:24 读数）逐号相同** |
| 3 | `sh scripts/slo-nail-census.sh --with-self-test` | **rc=0**；`positive control PASSED - 32 control(s) ran, 0 failed`；真读 `nail 1/2 …:233/:258/:291`、`nail 2/2 …:214/:230/:292`、`PASS - 2/2 … :533`（20:45） |
| 4 | `dash scripts/slo-nail-census.sh --with-self-test`（CI 的 `sh` 是 dash，本机 `/usr/bin/dash` 在） | **rc=0**，末三行与第 3 行同（20:50） |
| 5 | **本腿自搭台件**：`mktemp -d` 里放 `scripts/{slo-nail-census.sh,slo-check.ps1}` 两枚拷贝，`sed` 删两枚调用行 | **rc=1**，逐字 `slo-nail-census-nail-call-missing` ×2（两枚都点名）＋ `FAIL 2 declaration violation(s)`；还原拷贝 ⇒ **rc=0** `PASS - 2/2 … :533`（20:45） |
| 6 | **本腿自量的一枚更小绕过**：台件里把 `declared_nail_count` 改 1、两张表各删一行（单文件 3 行编辑） | **rc=0**；再删掉被删那枚钉的调用行 ⇒ **仍 rc=0**（20:48）。见 §2 AC#1 残余面 |
| 7 | `sh scripts/d22scan.sh`（正控 `runtests.sh` ＋ 真扫） | **rc=0**；`clean - no D22 ban violations`；分母 `bans#1-5 internal/=363 cmd/=58`、`ban#6 frontend/=116`、`ban#7 internal/tools/=48`、`ban#8 design/=39 frontend/=116 internal/=905 cmd/=211`；`runtests.sh: PASS=35 FAIL=0 SKIP=0`（20:46:49→20:47:56） |
| 8 | `sh scripts/check-path-length-budget.sh --with-self-test` | **rc=1**（**既有红，非本票**）；`denominator read: 18791 tracked paths`；`over-budget=65 covered by roster=64 not in roster=1`；self-test `3/3 PASSED`；红名**恰两枚**且都是票 179 那对（20:46:57→20:47:05） |
| 9 | `git status --porcelain -- scripts tools .github` | **空**（20:44 与 20:51 各一次；本腿全程只读） |
| 10 | `git diff --stat 6589a0625..HEAD -- scripts tools .github` | **仅** `.github/workflows/ci.yml | 14 +`（＝`6e5957756`），`scripts/`、`tools/` **0 行** |
| 11 | `.github/workflows/ci.yml` YAML 解析（utf-8） | 6 个 job 未变；`lint.steps` **15 枚**；census 步 `index=6`，`keys=['if','name','run']`（**无** `continue-on-error`），前一步 `Tracked path-length budget (ticket 262)`、后一步 `gofmt (gofumpt)` |
| 12 | 证据件 §6 围栏 vs `ci.yml:209-221` | **13/13 行 byte-exact = True**（20:47） |
| 13 | 正控自锁：台件里把一条调用**复制**一行（"健康拷贝必须绿"这条控制必红） | `--with-self-test` ⇒ **rc=1**、`5 control(s) failed`、末行 `positive control RED - the census is not run, because a ruler that cannot go red does not get to print a reading`；**真读未被印出**（20:52） |

## 2. 逐 AC 裁决

| AC | 票面要求（简） | 证据件 | 本腿现跑尺 | 判 |
|---|---|---|---|---|
| **#0** | 落地腿现跑复认 §现量 1/2/3、现读两枚钉真身行号 | r2 §1（＋r1 票尾 bullet） | §现量 1：`grep -c '"scripts'`＝0/rc=1，`declaredScopes()`（`tools/d22scan/main.go:403` 起）只列 `internal/ cmd/ frontend/ internal/tools/ design/`，且 ban #9 自己的射程注释 `main.go:50-66` 逐字写 `tools/**` 与 `_test.go` 不在内、纳入是 owner 批准的 scope change ⇒ 现读复认。§现量 3：P3 常量 `scripts/slo-freshness.sh:137`（`SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3`）、判定 `:549`（`-gt`）＋`:550`、只在默认分支评估的入口 `slo-fresh.yml:46-52`＋`:41-44` 自陈 ⇒ 现读复认。§现量 4：`:214/:230/:233/:258/:291/:292/:533` 现读，与 r2 §1 同号。⚠ **§现量 2（删钉后门 rc=0 全绿）门级未现跑** | **PASS（1 项具名替代）**：§现量 2 的门级重跑没做（r1 自陈、r2 归 AC#2），本腿用两把现读替代其**作用面结论**——① 全仓对 `slo-nail-census` 的引用只有 `ci.yml:211/:221` 两行（`grep -rn` 20:47，其余命中全在证据件）；② P3 看产物年龄、对钉存在零敏感。⇒ 「删钉今天无人发现」成立。若编排者按字面读 AC#0，**这是唯一没现跑的分句**（成本＝整条 PowerShell 门 + 采样） |
| **#1** | 外牙形状落在被检文件之外、具名 `file:line`；答"凭什么发现不了删检查也删声明"，别当牙 | r2 §1 | 声明在 `scripts/slo-nail-census.sh:110-112`（`declared_nail_count=2`／`nail_names`／`nail_ok`），被检件 `scripts/slo-check.ps1` 全文无这些声明（现 grep 0 命中）；判据 1–5 各自指向被检件真身行号（见 §1 表第 2 行）。绕过面 r2 具名三层：删两枚调用必响（本腿实测 rc=1）；删 census 文件＋CI 步则全静（本腿 `grep -rn` 复认无别的调用点）；静态尺不证可达 | **PASS（附 1 枚残余面）**：形状在外、`file:line` 齐、明说"不是关上的门"——三条都满足。但 r2 §1 层 2 那句「**体积最小的一发＝删文件+删 CI 步**」**被本腿实测推翻**：真正更小的是**就地缩声明**（单文件 3 行 ⇒ rc=0；再删那枚钉的调用行 **仍 rc=0**，台件探针，20:48）。⇒ 这枚绕过**没被具名**，"最小"那句是错的；不据此判 AC#1 失败（该 AC 要的是别把牙说成门，而它逐字说了），但**必须登记**，见 §4 |
| **#2** | 正控两发（删钉必响逐字／还原绿）＋还原证（`git status --porcelain -- scripts tools .github` 空） | r2 §2、§5 | 发①：本腿自搭台件删两枚调用 ⇒ **rc=1**，逐字 `slo-nail-census-nail-call-missing`（两枚钉各一行，含各自定义行号 `:233`/`:214`）＋ `FAIL 2 declaration violation(s) - …（ticket 266）`；发②：还原 ⇒ **rc=0** `PASS - 2/2 nails … 采样 :533`。台件与真件字节一致（`cp`）。还原证：`git status --porcelain -- scripts tools .github` **空**（20:51）。32 发正控 `sh` 与 `dash` 各 rc=0、0 failed | **PASS** |
| **#3** | 四条偷懒形状逐条自查并点名"我没做哪一种" | r2 §3 | ① 动阈值/golden：`git show --name-only 6589a0625` ＝ 仅 `scripts/slo-nail-census.sh`；`git diff 6589a0625^..HEAD -- internal/observe/thresholds.go` **空**；脚本里无 SLO 数值比较。② 放宽 `Set-StrictMode`：`slo-check.ps1` 自 `cb8dceccc`（263-r1）未动，`:82 Set-StrictMode -Version 2.0` 在。③ 改跳过：CI 步无 `continue-on-error`/`if: false`；脚本无可跳过的步，正控红时**拒出读数**（§1 表第 13 行实测：`positive control RED - the census is not run`，rc=1）。④ 吞：`grep -n "|| true"` ＝ **0**；`2>/dev/null` ＝ **2 处**（`:306`/`:429`）且两处都被 `|| cannot_decide` 接住；`catch` 只在 `:481` 注释里 | **PASS** |
| **#4** | 卫生门读数不扩大（d22scan／路径长度门／票 212 ban #9），红名逐名比对带时刻；动 scope 须具名报 | r2 §4、ci 证据 §2B | d22scan **rc=0 `clean`**（分母见 §1 表第 7 行；`ban #9` 的 `TestBan9PrefixTablesAreTheSameList` 在正控组内 PASS；ban #9 射程只 `internal/`+`cmd/` 生产 Go，`scripts/` 本不在内）。路径门 **rc=1** 18791/65/64/1，红名**逐名＝票 179 那对**（"over budget and NOT in the roster"＋"roster entry … drifted away"，名与 r2 18:27/18:29 两刻相同）⇒ **不扩大**。分母 18773（r2 交件）→**18791** 是**别的腿**在这之后新增的 tracked 路径（`git diff --stat 6589a0625..HEAD -- scripts tools .github` 只有 ci.yml 14 行），而 `over-budget 65→65`、`roster 64→64`、`not in roster 1→1` **不变**。`allowlist.txt`／路径门 roster／`slo-check.ps1` 现跑 `git diff` 全空 | **PASS**（路径门的红是**既有红**，归票 179 改名 `-done` 而 roster 写旧名，非本票；本腿未翻任何豁免） |

## 3. CI 步（形 ①-A 的另一半）现核

| 项 | 读数 |
|---|---|
| 在场与位置 | `ci.yml:209-221`，`lint` 作业 `steps[6]`，夹在路径门（`[5]`）与 `gofumpt`（`[7]`）之间；6 个 job、15 步（原 14） |
| 形状 | `keys=['if','name','run']`，**无** `continue-on-error`、**无** `if: false`；`on:`/`paths:`/并发/`runs-on`/既有步一字未动（`git diff --numstat` ＝ `14 0`，本腿复认 diff 只有这一处） |
| 与证据件一致 | 与 r2 §6 围栏 **13/13 行 byte-exact** |
| 会跑 | 命令本机 `sh` 与 `dash` 各 rc=0（CI 的 `sh` 是 dash）；`lint` 作业 `timeout-minutes: 20`（`ci.yml:103`），该步 ~1s，无 hang 面 |
| 买到什么 | 从 `6e5957756` 起「有人删/注释掉 `slo-check.ps1` 两枚调用行」⇒ 每次 push/PR 的 lint 红，且红句具名（`…-nail-call-missing`/`…-ordering`/`…-declaration-drift`/`…-cannot-decide`） |
| 没买到什么 | ①"那两行 ok 句真被打印过"仍无人断言（形 ② 未做，每次 PASS 的 NOTE 逐字复述）；②"删 census 文件＋删 CI 步"这一发仍全静；③ 就地缩声明不响（§4） |

## 4. 残余面与登记（⛔ 不影响判勾，编排者要收）

1. **证据件「体积最小的一发」这句是错的**（AC#1 附注）：r2 §1 层 2 写「删 `slo-nail-census.sh` ＋ 删 CI 步」是最小绕过；本腿实测**更小的一发是就地缩声明**——把 `declared_nail_count` 改 1、两张表各删一行（单文件 3 行）⇒ `rc=0`；若同时删掉那枚钉的调用行，**仍 `rc=0`**（探针 A/B，20:48）。
   ⇒ 建议：票尾追加一条 bullet 具名此事（或立下一票把"声明表被就地编辑"纳入射程，例如外部件读声明表行数/名字的哈希）；⛔ 不要据此把 AC#1 判红，该 AC 的字面（形状在外、具名、别当门）已满足。
2. **"删文件＋删 CI 步"绕过**：现读全仓只有 `ci.yml:211/:221` 两处引用，删这两处＋删脚本文件即全静；P3 按天钝且 dev 上 inert（`slo-fresh.yml:41-44` 自陈）。此条 r2 §1 已具名，本腿复认。
3. **PASS 的 NOTE 已在每次真读时打印**（"pins shape only … one commit deleting this file together with the CI step … would leave nothing ringing"），未藏。
4. **既有红不归本票**：路径门两枚红名＝票 179 改名 `-done` 后 roster 漂离；d22scan 现跑 rc=0 无红。本腿**未**动 roster／allowlist／阈值。
5. 本腿未跑 `go build`（本票零 Go 改动；d22scan 现跑已编译 `tools/d22scan`），未跑 PowerShell 门本体（理由同 AC#0 脚注）。

## 5. 一句话可结

**可结（CLOSE）**：外牙 `scripts/slo-nail-census.sh`＋lint 作业的 CI 步两半都在盘上、都现跑过——删两枚调用 ⇒ rc=1 逐字 `slo-nail-census-nail-call-missing`、还原 ⇒ rc=0 `PASS - 2/2`、32 发正控 0 失败（sh/dash 各一发）、CI 步 byte-exact 且可跑、AC#0–#4 读数逐条对上、d22scan `clean`、路径门红名集合与交件时**逐名相同**；两处已具名残余（「删文件+删步」全静、**就地缩声明比证据件所称的最小绕过更小且 rc=0**）不属本票要关的门，**建议收件并在票尾追加第 4 节第 1 条那条 bullet**。

> 撤销口令（若编排者要销本裁决）：删除本文件即可，盘上无别的变化；票 266 的两半交件不受影响。
