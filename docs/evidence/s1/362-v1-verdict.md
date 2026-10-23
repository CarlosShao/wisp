# 票 362 终裁腿 `362-v1` —— `AC#2`（落地＋牙）**能翻**

- **腿**：`362-v1`＝**非实现者裁腿**（裁决者≠实现者）；裁的对象＝`362-r2` 三笔 `8ebbc1491`／`ade547ac8`／`0754fa7d9`
- **裁时**：2026-10-23 11:0x–11:2x +08。**决定性读数钉在 `2bed76832` 的树**（我导出的仓外树）；此后 HEAD 一直在走（落笔时 `7be24aaea`，`2bed76832..7be24aaea` 共 **7** 笔），而 `git diff --name-only 2bed76832..7be24aaea` **`.go` 0 枚、`ci.yml` 0 行、`attrib.sh` 0 行**，`tracked .go` 仍是 **1165** ⇒ **读数不因 HEAD 移动而失效**。
- **写面**：只本件 ＋ `D:\tmp\wisp\20261023\362v1\**`（仓外）。**所有突变都造在仓外导出树上**，仓内产码／门禁／`.github/**`／`.scratch/wisp/probes/**` **一字未动**；⛔ 洗任何 `.go` 字节。
- **一句话**：**两道门都真的不再对故意坏的台件报红，而产码目录的真脏两门都还看得见** ⇒ `AC#2` 的"落地＋牙"两条都成立、反向判据 20/20 IDENTICAL、`--self-test` 10/0/rc=0、`ci.yml` 语义零改动、三笔越界 0。

---

## §1 票面前提新鲜度（对 HEAD blob 复量）

| 量 | 票面（`AC#0`） | `362-r2` 报（其锚树 `938216b8e`） | **我现量（`2bed76832` 的树）** |
|---|---|---|---|
| `tracked *.go` | 1156 | 1162 | **1165** |
| 五个产码目录 | — | 802 | **805**（internal 626／cmd 149／tools 18／scripts 11／frontend 1） |
| `.scratch/**` | 20 枚脏 | 360 | **360** |
| 根／其它目录 | 0 | 0 | **0**（805＋360＝1165，与总数**对上**） |
| 判脏（LF 导出树上跑 r2 之前的原件） | 20 | 20 | **20**，且 **20/20 都在 `.scratch/wisp/probes/**`** |
| 那 20 枚的 **blob CR** | — | — | **逐枚 0**（`git show HEAD:<p> \| tr -cd '\r' \| wc -c`） |

⇒ ★**票面的实质前提（20 枚入库就脏、全在台件目录、母数是 tracked 集）在 HEAD 上仍成立**；只有**计数在长**。
★**一处要报的偏**（⛔ 阻塞翻框）：`ci.yml:191` 逐字写 "of **1162** tracked .go files, **802** …（internal **623**…）"，
而**本笔自己那棵树**（`0754fa7d9`）就是 **1165／805／internal 626**（我把它的树 `git archive` 到仓外数的）。
1162/802 是它量过的 `938216b8e`——**比它自己的基 `6e19f4a43` 旧 12 笔**。⇒ 那句话的**锚点选错**，不是假话，但按 commit 复现会对不上。

## §2 `AC#2` 判点一：**两道门都不再对台件报红**（我自己跑，仓外）

- 树：`git archive 0754fa7d9` → `T_R2`（17616 件）／`git archive 2bed76832` → `T_NOW`（17618 件）；各**自带 `.git`＋index**、`*.go` 的 CR 总数 **0**、`git status` **0** 条。两树 `.go` 名册**逐名相同**。
- **门①**（`ci.yml:217-226`，按 `run:` 行**逐字**）`OUT="$(gofumpt -l internal cmd tools scripts frontend)"` ⇒ `T_R2` **0 行**、`T_NOW` **0 行**（rc 皆 0）。
- **门②**：`ATTRIB_ROOT=<树> sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`
  （★**⛔ `cd` 进树再跑** —— `attrib.sh:141` 是 `d=${ATTRIB_ROOT:-$here}`，`cd` 进树跑的是**主检出**，那是尺的射程错）
  ⇒ `T_R2` **rc=0／lines=0／files=0**、`T_NOW` **同**；分母打印**两树皆** `handed=805 total=1165 subtracted=360`。
- ⇒ **CI 形状（新检出、纯 LF）上两道门都转绿**。
- ★**射程声明**：门①我只跑了 `run:` 块里**判据那行**（`OUT="$(gofumpt -l …)"`），**⛔ 跑 `go install mvdan.cc/gofumpt@latest`**（要网、且会装"当时最新"）；用的是本机 **`gofumpt v0.12.0 (go1.27.1)`**。⇒ 两门读的是**同一枚二进制**，"CI 装到更新版会不会改判"是**既存**面、⛔ 本笔引入。

## §3 `AC#2` 判点二（本票最要紧）：**产码目录的真脏还看得见**

仓外 `T_PLANT`（`T_R2` 的拷贝）`git add` 成 tracked 的三枚：

| 控 | 路径 | 故意坏？ | 在哪 | 门①输出 | 门② |
|---|---|---|---|---|---|
| **P1** | `internal/agent/approval/zz_362v1_posctl_unformatted.go` | **是** | **产码目录** | **列出 1 行**（`internal\agent\approval\…`，★反斜杠）⇒ 走 `exit 1` 分支 | **rc=1**、TRACKED-DIRTY **1 行**、逐名同上 |
| P2 | `internal/agent/approval/zz_362v1_posctl_clean.go` | 否（干净） | 产码目录 | **不列** | **不列** |
| P3 | `.scratch/wisp/probes/999/zz_362v1_posctl_dirty.go` | **是** | **排除面内** | **不列** | **不列** |

- 撤三枚后：门① **0 行**、门② **rc=0／0 行** ⇒ **红是那枚 P1 造的，不是恒红**。
- ★再补一条**反向射程**控：把 `A_N_EXCLUDE` 人为放大成 `:(exclude)*`（仓外变体）⇒ **rc=2 拒答**
  （`ruler saw 0 tracked .go files after its one named exclusion (:(exclude)*) - refusing to report 'empty' from a ruler that was handed nothing`）
  ⇒ **空名册不会被当绿**、且**分母三联数会跟着动**（排除只能"移动一个打印读数"，⛔ 静默缩名册）。
- ⇒ ★★**"为了让台件不红而把真脏弄瞎"没有发生**：排除面是**具名 `:(exclude).scratch/**`**，减法是 **git 自己做的**（⛔ `grep -v`——那会把"分母"和"判据"混成一件事）。

## §4 反向判据：撤两处排他 ⇒ **20 枚逐名重现、diff 空**

两个**仓外**变体：`pre` ＝ `6e19f4a43` 上的**真原件**（r2 之前，662 行）；`rev` ＝ 今天的件（705 行）**只撤掉排他**——
`sed` 命中 ` "$A_N_EXCLUDE"` **3 行**：`580`／`585` 是**两处 `git ls-files` 实参**（功能面），`574` 是**注释里**同一串（`HOW TO REVERT` 那句，无功能影响）；两变体 `sh -n` 皆 rc=0，**diff 就这 3 行**。

| 树 | `pre` | `rev` | 名册 diff |
|---|---|---|---|
| `T_R2` | rc=1 lines=**20** files=**20** | rc=1 lines=**20** files=**20** | **空 ⇒ IDENTICAL** |
| `T_NOW` | rc=1 lines=**20** files=**20** | 同 | **空 ⇒ IDENTICAL** |

另：`pre(T_R2)` vs `pre(T_NOW)` 也 **IDENTICAL**。20 枚逐名（`rev`＝`pre`，两个独立来源同一份；下列每行的公共前缀都是 `.scratch/wisp/probes/`）：

```
 1 33/p1/q1/main.go                                11 222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go
 2 33/p1/q2/main.go                                12 224/v2/dialect_probe_test.go
 3 33/p1/q3/main.go                                13 231/v1/mutation/config_reload_branch_after_prefix.go
 4 163/a1/main.go                                  14 235/v1/mut/prereadoff-m3_222_test.go
 5 174/c2/zz174c2_wiring_pair_windows_test.go       15 241/v1/posctl/badly_formatted.go
 6 183/r2/mut/task-boxset-off.go                   16 241/v1/probe_v1_readings_test.go
 7 185/r1/mut-m1/hostpath_185.go                   17 253/v2/denoms253v2.go
 8 197/r1c/pre/subagent_197.go                     18 259/r1/gofumpt-negctl/probe.go
 9 197/r1c/pre/subagent_197_test.go                19 263/v1/src/main.go
10 212/v1/mut/main-noq9.go                         20 220/r1/prechange/prechange220_probe_test.go
```

⇒ 这 20 枚与 `docs/evidence/s1/362-a1-census.md:29-48` 那张普查表**逐名相同**（第三源交叉核过）⇒ **不是"我这一把尺自己跟自己一致"**。

## §5 `--self-test` 仍绿

`sh .scratch/wisp/probes/161/r5/attrib.sh --self-test` ⇒ **`cases=10 failures=0`、rc=0**；
roster nails 逐项 **`PASS want ring=5 silent=3 named=3 / got ring=5 silent=3 / all named must-ring samples present`** ⇒ **票 161 的载具未破**。

## §6 `ci.yml` 的 YAML 语义零改动

- ★**派单给的命令逐字跑 ⇒ rc=1** —— 但**不是 YAML 的错**：本机 `python` 默认编码 **`cp936`**，`open()` 不带 `encoding` 会在**第 41 行**那个 `加`（UTF-8 `e5 8a a0`）上抛 `UnicodeDecodeError`。显式 `encoding='utf-8'` ⇒ **rc=0**、**6 job／53 step**（`lint, test-core, test-windows, slo-smoke, slo-full, lint-frontend`）。
- 非 ASCII 行：`ade547ac8^` **14** ／ `ade547ac8` **14** ／ HEAD **14** ⇒ **本笔一个非 ASCII 字节都没夹进来**；三件（`ci.yml`／`attrib.sh`／`362-r2-impl.md`）**CR=0**。
- **步名与 `run:` 行一枚未动**：三份（`ade547ac8^`／`ade547ac8`／HEAD）**去掉行号后**逐行 `diff` **空**（79 → 79 → 79 行）；**全部 `run:` 行**（before-r2 vs HEAD）**逐字相同**。行号整体 +11 是注释块长了，⛔ 语义。
- `ade547ac8` 的 `+/-` 行里 **非注释行 0**（`grep -v '^[+-][[:space:]]*#'` 后为空）⇒ **纯注释笔**。

## §7 门禁复跑（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；`go env GOTMPDIR` 为空 ⇒ `d22scan` 311 守卫不会 `rc=2`）

| 命令 | rc | 读数 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | 无输出 |
| `sh scripts/d22scan.sh` | **0** | `examined 290 production Go files`；ban #1–#8 全 clean |
| `sh scripts/check-path-length-budget.sh` | **0 GREEN** | 11:25:37 复跑：`tracked paths=17620 over-budget=61 covered by roster=61 not in roster=0` ⇒ `VERDICT GREEN`。★**11:1x 第一发是 `rc=1 RED`**，唯一不在 roster 的是 `.scratch/wisp/issues/365-two-taskLogState-…-statusstuck-…-done?.md`（相对 126／名 **105**）；该件由 **`8beec66b1`**（编排者立票 365）新增、`6e19f4a43` 上**不存在**、三笔 r2 **一笔未碰**；编排者 **`a957dfdc0`（A1170）**已把它改名成 `365-two-taskLogState-functions-disagree.md` ⇒ **现在绿**。**两次都与本笔无关** |

**格式门甲乙两尺（分开报、⛔ 相减）**
- **甲（判据＝新检出／入库字节）**：`T_R2`／`T_NOW` ⇒ 门① **0 行**、门② **rc=0／0 行**。
- **乙（诊断＝主检出工作树，11:22:34 取）**：门② **rc=1／4 行**；同树门① **4 行**（同一批）。四枚＝
  `cmd/wisp/models.go`（blob_cr=0／wt_cr=389／diff 0 行）、`internal/agent/approval/pending_read.go`（0／131／0）、
  `internal/agent/tools.go`（0／225／0）——**3 枚 CRLF 假红**；＋`internal/proc/envfork_test.go`（blob_cr=0、**wt_cr=0**、
  **`git diff` 1 行／`status` 1 行／mtime 11:13:39**）——**别腿在飞编辑**、⛔ CRLF、也⛔ 真格式债。乙尺**不参与**判据。

## §8 新增红 0／写面纪律

- 三笔逐笔 `git show --name-only`：**各 1 文件、`.go` 0**；`8ebbc1491`→`.scratch/wisp/probes/161/r5/attrib.sh`；
  `ade547ac8`→`.github/workflows/ci.yml`；`0754fa7d9`→`docs/evidence/s1/362-r2-impl.md`（**154 行**，与自报一致）⇒ **越界 0**。
- **正控（同尺打一笔含 `.go` 的真提交）**：`5413f46d4` ⇒ 5 文件／**5 枚 `.go`** ⇒ 这把尺**看得见 `.go`**，三笔的 `go_files=0` **不是尺瞎**。
- **新增红**：门①／门② 在甲尺上**由红转绿**；门禁三条**现全绿**（`go build` 0／`d22scan` 0／路径长度 0，后者是本腿跑完后别腿改名转绿的，两次都⛔ 本笔）⇒ **本笔新增红 0**。

## §9 编排者那两条裁，我独立判

**A)「本机那 3 行是 `autocrlf` 工作树假象 ⇒ CI 上为 0 ⇒ `attrib.sh` 在 CI 上 `rc=0`」＝成立**，三向复核：
① 三枚 `blob_cr=0`、`wt_cr=389/131/225`、`git diff` **0 行**、`git status` **0 行** ⇒ 脏**只在盘上**；
② **新导出树**（`git archive`＋新 index）上三枚**根本不出现**、门② `rc=0`；
③ ★**机制比"autocrlf 不生效"更硬**：`.gitattributes` 逐字有 **`*.go text eol=lf`** ⇒ **新检出的 `.go` 必是 LF**（与本机 `core.autocrlf=true` 无关；本机那 3 枚 CRLF 是"编辑器直接写盘"的，`git add` 归一后**永不进 index**）。
★但我采纳它的**"⛔ 归一／⛔ 改尺读 blob"**这条裁时，拿到一条**比它更硬的理由**：乙尺这轮**恰好**抓到一枚"工作树里的真脏"（§7 的 `envfork_test.go` 在飞编辑）——**改成读 blob 就完全看不见它**；而"已入库字节"那一面**甲尺已经量了** ⇒ 改读 blob 是**用一条能看见真脏的射程，换一条重复的射程**。

**B)「`AC#3` 按入库字节判＝成立」＝成立**（甲尺 rc=0）。★须写清：**取到它靠两道门一起改**——只改 `ci.yml` 时 `attrib.sh` 仍以同样 20 个名字 `rc=1`（§4 的 `pre` 尺把这 20 个名字原样复现了）。

**C) 它自报的「乙尺 34 行 → 3 行（34 ＝ 31 `.scratch` ＋ 3 CRLF）」＝数对、但会腐**：
我 **10:50 直接量到** `pre` 尺 **34 行＝31 枚 `.scratch` ＋ 3 枚产码**（逐名；那 3 枚正是它点名的那 3 枚）⇒ 同刻扣掉台件应得 **3**，与它自报一致；
今天 **11:22 直接量** `post` 尺＝**4**（多了 `envfork_test.go`）⇒ **"34→3" 成立、"落差是 pre-existing" 也成立**；
★但**这个数只在某一刻成立**（共享工作树会动），引用**必须带时刻**（结论不受影响：甲尺仍 0）。

## §10 `AC#2` 判语：**能翻**

- **落地**：`AC#1` 拍的〔甲′〕"改调用面、⛔ 洗任何字节"在**两道门**上都落了（`ci.yml:221` 目录集；`attrib.sh:577/580/585` 具名排除），**`.go` 字节一枚未动**（三笔 `go_files=0`）。
- **牙**：门①／门②在 P1 上**都红且都逐名**；P2／P3 都不列；撤控即回绿；放大排除面 ⇒ **rc=2 拒答**；《反向 20/20 IDENTICAL》＋《`--self-test` 10/0》两枚护栏都在。
- **缺什么（都不阻塞翻框，只防下游误读）**：①`ci.yml:191` 的 1162/802/623 是 `938216b8e` 的数，本笔自己树上应是 **1165/805/internal 626**；
  ②乙尺读数**不稳定**（共享工作树：同时刻另有 **84 枚**未提交 `.go`），任何"乙＝N 行"都要带时刻；
  ③路径长度门那发 RED（票 365 文件名 105）已被编排者 `a957dfdc0` 改名解决、现绿 ⇒ **本票无残留**。

## §11 我自己的两处尺错（先记，别让人踩）

1. 我第一次判归属写成 `[\/]\.scratch[\/]` ⇒ **前导 `.scratch` 前面没有分隔符**，命中 **0**，把 31 枚台件全算成"非台件"（清单第 4 例的变体）。正解＝先 `tr '\\' '/'` 再 `grep '^\.scratch/'`（31＋4＝35，与总行数对上）。
2. 我照派单跑 `python -c "…yaml.safe_load(open('ci.yml'))"` 得 **rc=1**，差一点记成"YAML 坏"（清单第 2 例"输出形状"）。真因＝**cp936 默认编码**，⛔ YAML。

## §12 提交

本件 `docs/evidence/s1/362-v1-verdict.md`（**收手 146 行**；第一笔时 143 行、本笔补 §12 到 146；**≤160 ⇒ 未拆 `-part`**）。
- 本笔第一笔：**`74fc2e4fb`**（父 `eac86c4e7`），`git show --name-only` ＝ **1 文件**、`.go` **0**。
- **提交后复跑**（`git archive 74fc2e4fb` → 仓外 `T_POST`，17621 件、`tracked .go` **1165**、`*.go` 的 CR **0**）：
  门① **0 行**；门② **rc=0／0 行**，分母 `handed=805 total=1165 subtracted=360` ⇒ **绿是提交后的绿，不是提交前的**。
