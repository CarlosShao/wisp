# 276 — d22scan 探针件归口 / 未格式化门 / 豁免边界（写腿 `leg 276`）

> 只产码，⛔ 未改任何票面 md（AC 框、`Status:`、票名一字未动），⛔ 零 push，零 `add -A`。
> 起手锚 `bd88a8ed1`（`dev`）。写面严格限于 `tools/d22scan/**` ＋ `scripts/**` ＋ `.scratch/wisp/probes/276/**`。
> 门禁用 Git Bash 的 `sh`（`D:\work\soft\Git\bin\sh.exe`；`C:\Windows\system32\bash.exe` 是 WSL 存根，⛔ 不可用）。

---

## §0 结论先行

- **AC#0 已交付**：17 枚探针件逐枚归口并给出**引用者**（尺＝`git grep -F` 对完整路径，⛔ 不按文件名猜）。
  实测 **17/17 枚全部有引用者**，其中 **3 枚是正控/负控或逐字节凭据**（见 §1）。
  按票面 AC#0 最后一句"报不出引用者的那一枚不许格式化"——**没有一枚落在这一支**。
  ⚠ **但本腿一枚都没有格式化**，理由是 §2 顶回的**门形状已变**：这条门今天不吃 `.scratch/**`。
- **AC#1 已交付（正控成立）**：在仓外临时 git 树（`D:\tmp\ctl276-ac1`）种一枚**新的**未格式化 tracked `.go`
  ⇒ 门**仍红且点名它**（`TRACKED-DIRTY scripts/planted_unformatted.go`，rc=1）。见 §3。
- **AC#2 已交付（名单有边界，四向突变）**：名单摘一枚 ⇒ 分母从 3 掉到 2、只剩名单外那枚红；
  名单多遮一枚 ⇒ 只点名单外那枚；名单遮住**全部** ⇒ 门**拒答 rc=2**（不是空绿）。见 §4。
- ⚠ **遭遇的真正阻碍（必须编排者裁）**：**票面 §"现量"第 1 条已过期**。那条门今天的 (A) 尺**自带
  `:(exclude).scratch/**`**（票 362-r2 落的），所以 17 枚探针件**今天根本不在门的分母里**。
  票面"17 枚跨 14 个目录让 CI 常红"这一前提**在本锚点上不成立**。见 §2。

---

## §1 AC#0 — 17 枚逐枚归口（含引用者与正/负控判定）

尺一：`git ls-files --error-unmatch <path>` ⇒ **17/17 TRACKED**（全在 `.scratch/wisp/probes/` 下）。
尺二：`gofumpt -l` 逐枚（HEAD 归档 `D:\tmp\arc276`，避开本机 CRLF）⇒ **17/17 DIRTY**，与票面数逐枚对上。
尺三：`git grep -n -F "<full path>"` ⇒ 每枚的引用者如下。

| # | 路径（剥 `.scratch/wisp/probes/`） | 归属票 | 引用者（具名） | 正/负控或逐字节凭据？ |
|---|---|---|---|---|
| 1 | `163/a1/main.go` | 163 | `probes/212/a2/work/citing-files-go.txt:49`、`den-go.txt:74`；`issues/361-*.md:35` | 否（只读核台件程序） |
| 2 | `174/c2/zz174c2_wiring_pair_windows_test.go` | 174 | **`probes/174/c2/overlay.json:1`**（`-overlay` 夹具）＋`212/a2` 名册三处 | 夹具（overlay 载体），非正负控 |
| 3 | `183/r2/mut/task-boxset-off.go` | 183 | **`probes/183/r2/overlay-boxset.json:3`**（突变副本）＋`212/a2` 名册 | **突变载体**（牙的载体） |
| 4 | `185/r1/mut-m1/hostpath_185.go` | 185 | **`probes/185/r1/ov-m1.json:1`**（突变副本）＋`212/a2` 名册 | **突变载体** |
| 5 | `197/r1c/pre/subagent_197.go` | 197 | `issues/236-*-done.md:138`（`cmp` 逐字节 IDENTICAL 声明）＋`212/a2` 名册 | **逐字节凭据**（`d1c440ad` 的改前快照） |
| 6 | `197/r1c/pre/subagent_197_test.go` | 197 | `236/a1/census.md:163`、`236/a2/census.md:19`（25,429 字节、同一条 `:93` 声明） | **逐字节凭据**（同上） |
| 7 | `212/v1/mut/main-noq9.go` | 212 | **`probes/212/v1/mut/ov-noq9.json:3`**（`tools/d22scan/main.go` 的 overlay 变异）＋`236/r1/evidence.md:197` | **变异载体（本工具自己）** |
| 8 | `220/r1/prechange/prechange220_probe_test.go` | 220 | `212/a2` 名册三处＋`236/r1/evidence.md:198`（"纯格式（改前快照）"） | 改前快照 |
| 9 | `222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | 222 | `probes/172/r1/logs/11-unfixed-grep-census.txt:104`＋`212/a2` 名册 | 突变副本 |
| 10 | `224/v2/dialect_probe_test.go` | 224 | **`probes/224/v2/overlay.json:3`**＋`docs/evidence/s1/224-session-grant-v2.md:269` | 夹具 |
| 11 | `231/v1/mutation/config_reload_branch_after_prefix.go` | 231 | **`probes/231/v1/overlay-branch-after.json:1`**（+`make-branch-after.ps1:4` 生成器） | **突变载体** |
| 12 | `235/v1/mut/prereadoff-m3_222_test.go` | 235 | `probes/172/r1/logs/11-unfixed-grep-census.txt:115`＋`212/a2` 名册 | 突变副本 |
| 13 | `241/v1/probe_v1_readings_test.go` | 241 | `212/a2/den-go.txt:210`＋`docs/evidence/s1/241-audio-level-producer-v1.md:49` | 探针件（⛔ 非票面点名的两枚台件之一） |
| 14 | `263/v1/src/main.go` | 263 | `issues/275-*-done.md:103`（`go.mod` 上下文那枚）＋`236/r1/evidence.md:205` | 否 |
| 15 | `33/p1/q1/main.go` | 33 | `212/a2` 名册三处＋`236/r1/evidence.md:187`（"台件"） | 台件 |
| 16 | `33/p1/q2/main.go` | 33 | `issues/35-panel-bridge-c17.md:93/:100/:121`（**唯一 Go→页调用点**的指认） | **唯一证据点** |
| 17 | `33/p1/q3/main.go` | 33 | `212/a2` 名册三处＋`236/r1/evidence.md:189` | 台件 |

**票面 §"现量"第 2 条那两枚台件**（⛔ 本票不许格式化，本腿亦零字节）：
`241/v1/posctl/badly_formatted.go`（票 241 正控）、`259/r1/gofumpt-negctl/probe.go`（票 259 负控）——
两枚均在 HEAD 归档上 `gofumpt -l` **点名**（见 §2 的 20 枚全名册），且均**不在**本票 17 枚之内。

---

## §2 ★顶回：门形状已变，票面"现量"第 1 条过期

### 2.1 现读：这条门的 (A) 尺带一条具名排除

`.scratch/wisp/probes/161/r5/attrib.sh:577`：

```
A_N_EXCLUDE=':(exclude).scratch/**'
```

`:579-585` 两次把它交给 `git ls-files`，`:583` 还把三个数打印出来。**现跑**（本锚，工作树）：

```
   (A) denominator: handed to gofumpt = 1065; total tracked .go = 1425; subtracted by the one named exclusion :(exclude).scratch/** = 360
```

⇒ **360 枚 `.scratch` 下的 tracked `.go` 今天不进门的分母**，17 枚全在其中的 360 里。
落这条排除的是**票 362-r2**（`attrib.sh:544-576` 逐字写着它的动机与"HOW TO REVERT"）。
`ci.yml:340` 跑的就是 `attrib.sh --tracked-only`，即这条带排除的 (A) 尺。

### 2.2 本锚今天真实的红名册（不在 `.scratch`）

**现跑** `attrib.sh --tracked-only`：

```
   TRACKED-DIRTY	internal/agent/approval/pending_read.go	(last commit touching it: d03d166fe 10-04 10:37 CarlosShao)	<== CI IS RED ON THIS
   TRACKED-DIRTY	internal/agent/tools.go	(last commit touching it: 5413f46d4 10-04 09:19 CarlosShao)	<== CI IS RED ON THIS
attrib.sh: RED (tracked-only) - tracked-dirty=2 files_dirty=2, over 1065 tracked .go files
```

⚠ **这 2 枚不是真违规**：同一锚的 **HEAD 归档**（`git archive HEAD` → `D:\tmp\arc276`）上
`gofumpt -l internal/agent/approval/pending_read.go internal/agent/tools.go` ＝ **空、rc=0**，
`gofumpt -d` 亦空。⇒ 这 2 行是本机工作树的 **CRLF 污染**，不是格式缺陷
（与 `236-*.md:149` 已记过的"本机 32 枚 / 归档 19–20 枚"同一族，`A615` 亦记过 `cmd/wisp/models.go` CRLF）。
⇒ **判 CI 那道门的颜色只许用归档形状**，本机这 2 行不许当 CI 读数。

### 2.3 归档上的 20 枚全名册（＝`.scratch` 里今天真脏的全部）

尺：`D:\tmp\arc276` 上 `gofumpt -l .` ⇒ **20 行**。其中 **17 枚＝本票 17 枚**，另 3 枚为：

| 路径 | 是什么 |
|---|---|
| `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go` | 票 241 **正控**（票面第 2 条具名） |
| `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` | 票 259 **负控**（票面第 2 条具名） |
| `.scratch/wisp/probes/253/v2/denoms253v2.go` | 票 253 的分母台件（⛔ 不在本票 17 枚名册内，票面未列） |

**票面第 1 条现量与实测的差异（具名）**：票面写"去掉它的样本＝19 枚 / 再去掉两枚台件＝17 枚"；
本锚归档实测 **20 枚 = 17 + 2 具名台件 + 1 枚票面未列的 `253/v2/denoms253v2.go`**。
⇒ 17 这个数**仍成立**，但它今天**不代表"门的分母"**，只代表"`.scratch` 里脏的、且不属于另两枚具名台件的那一批"。

### 2.4 因此：AC#0 的"处置"落点是空的（本腿不格式化，理由具名）

票面 AC#0 要的是"逐枚归口，有引用者的可格式化"。但既然 17 枚**今天不在门分母里**：

- 格式化它们**不买任何绿**（门本来就绿，除非撞上 §2.2 那个 CRLF 假红）；
- 却**必付** §1 表里 6 枚的代价：`197` 两枚是 `d1c440ad` 的**逐字节快照**（`236-*.md:138` 用 `cmp` 声明的凭据）、
  `212/v1/mut/main-noq9.go` 是**本工具自己的 overlay 变异副本**、`183`/`185`/`231`/`224`/`174` 是 overlay 夹具。
  ⇒ 格式化＝**改动别人证据件的字节去买一个不需要的绿**，这正是票面"丙"的反面、`236-a2 §1 枚 5` 已量的代价。

⇒ 本腿的处置＝**零格式化**，并把"门形状已变"作为顶回交给编排者。
**若编排者仍要丙落地**（例如要撤掉 362-r2 那条排除），§1 的逐枚引用者表即为现成的执行清单，⛔ 但撤排除属改门语义＝须裁定。

---

## §3 AC#1 — 新种一枚未格式化 tracked `.go`，门必须仍红且点名（正控）

**做法**（⛔ 不原地编辑共享工作树，全部在仓外临时 git 树）：`D:\tmp\ctl276-ac1`，`git init`，
`.gitattributes` 写 `*.go text eol=lf`，两枚 tracked `.go`：`internal/clean.go`（干净）＋
`scripts/planted_unformatted.go`（**故意未格式化**，文档注释 bullet 紧贴 + 复合字面量括号不独占行）。
经 `ATTRIB_ROOT`（`attrib.sh:136-141`，仪器自带的仓外机制）交给**真门**跑 `--tracked-only`。

**现读**：

```
== (A) tracked tree: git ls-files -z '*.go' ':(exclude).scratch/**' | xargs -0 .../gofumpt.exe -l
   (A) denominator: handed to gofumpt = 2; total tracked .go = 2; subtracted by the one named exclusion = 0
   TRACKED-DIRTY	scripts/planted_unformatted.go	(last commit touching it: bb2dde8 11-03 22:03 ctl)	<== CI IS RED ON THIS
attrib.sh: RED (tracked-only) - tracked-dirty=1 files_dirty=1, over 2 tracked .go files
rc=1
```

- 种下的那枚**被逐名点名**，rc=1 ⇒ **正控成立**：门没有被写成"恒绿"。
- 同树里那枚干净的 `internal/clean.go` **不出现** ⇒ 红的成因就是新种那枚，不是噪声。
- ⚠ 一次**具名弯路**（留档，防后人重踩）：第一版用 `Set-Content` 写 `internal/clean.go`，PowerShell 落了 CRLF，
  于是**干净文件也被点名**（`lines=2 files=2`）。改成 `WriteAllText` 显式 LF 后复跑才得 `lines=1 files=1`。
  ⇒ 这正是 §2.2 那条 CRLF 陷阱的**主动复现**，也证明本节的正控不是撞出来的。

---

## §4 AC#2 — 豁免名单的边界（定向突变，四向）

**名单本体**：真门里唯一的"名单"就是 `A_N_EXCLUDE`（`attrib.sh:577`）。突变一律用
**仓外副本**（`D:\tmp\ctl276-ac2\M*.sh`，原 `attrib.sh` 一枚字节未动——`git status --porcelain` 对该文件空）。
**试验树**：`D:\tmp\ctl276-ac2`，三枚 tracked `.go`：`internal/clean.go`（干净）、
`listed/in_list.go`（脏、"名单内"）、`unlisted/not_in_list.go`（脏、"名单外"、**与 `listed` 同形状**）。

| 突变 | 名单 | 期望 | **现读** |
|---|---|---|---|
| 基线 | `.scratch`（不动） | 两枚都红 | `lines=2 files=2`，rc=1；分母 3/3 ✅ |
| **摘一枚**（M7） | `+.scratch +listed`，**摘掉 `unlisted`** | 只点 `unlisted` | `TRACKED-DIRTY unlisted/not_in_list.go` 唯一；分母 3→**2**；rc=1 ✅ |
| **放一枚**（M2，`listed` 名册 + `exclude`-first） | 只遮 `listed` | 只点 `unlisted` | 分母 3→**2**，只点 `unlisted`，rc=1 ✅ |
| **遮住全部**（M6） | `+.scratch +listed +unlisted +internal`（分母归零） | **拒答**，⛔ 不许空绿 | `ruler saw 0 tracked .go files ... refusing to report 'empty' from a ruler that was handed nothing`，**rc=2** ✅ |

⇒ **名单确有边界**：摘/放各一枚都让点名集与分母数**按预期变**；遮空时门**拒答**而不是报绿
（`attrib.sh:581` 的 `TRACKED_GO -gt 0` 闸门真在承重）。⛔ 本腿未改 `exit 2` 语义、未放宽任何断言。

### ★ 三条具名的仪器事实（本腿踩出来，供后续腿省一次弯路）

1. **本机 git-bash 里 `A='one' 'two'` 不拼接**：它把第一个词赋给 `A`、把 `'two'` 当**新命令**执行
   （实测 `probe4.sh` 原文与读数：`A='one' 'two'` ⇒ `A=[]` ＋ `line 2: two: command not found`；`B='x' 'y' 'z'` 同形同结果）。
   ⇒ 写**多条** pathspec 必须用
   **单引号包住整串**，再在调用点**不加引号地展开**才会 word-split 成多个参数。
2. **`git ls-files` 收到一个"空格连接"的单参数不会拆**：`'a b'` 被当**一条**字面路径 ⇒ 静默不排除任何枚
   （实测：`git ls-files '*.go' ':(exclude)listed/** :(exclude)unlisted/**'` ⇒ 仍 3 行）。
3. **`:(glob,exclude){a,b}/**` 花括号不展开**（实测 `git ls-files '*.go' ':(glob,exclude){listed,unlisted}/**'`
   ⇒ 3 行，一枚未排除）⇒ 这个仓的路径风格下**只能用逐条 `:(exclude)<path>`**，⛔ 别用花括号批量。

---

## §5 门禁读数（本腿现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./tools/d22scan`（模块内 `go build .`） | **rc=0** |
| test | `go test ./tools/d22scan/...` | **ok 52.253s，rc=0** |
| vet | `go vet .`（`tools/d22scan`） | **rc=0** |
| d22scan clean | `cd tools/d22scan && go run . -root "<repo>"` | **clean**，rc=0（`no D22 ban violations`；bans #1-5 internal/=348 / cmd/=52，ban #8 design/=39 frontend/=116 internal/=843 cmd/=188） |
| 真门正控 | `attrib.sh --self-test` | **GREEN**，`cases=10 failures=0`，rc=0 |
| 真门 (A) 尺 | `attrib.sh --tracked-only` | **rc=1**，`tracked-dirty=2`（⚠ 两行均为本机 CRLF 假红，归档上为空，见 §2.2） |

---

## §6 本腿**未**做的事（明确写出，⛔ 别当已做）

1. **未格式化任何一枚探针件**（理由见 §2.4：门今天不吃 `.scratch`，格式化只付代价不买绿）。
2. **未动 `attrib.sh` 一枚字节**（`git status --porcelain -- .scratch/wisp/probes/161/r5/attrib.sh` 空）；
   所有名单突变都在 `D:\tmp` 的副本上，原仪器的 `exit 2`、`A_N_EXCLUDE`、分母表达式一字未改。
3. **未动 `ci.yml`／`internal/**`／`cmd/**`／`thresholds`／`golden`／任何票面 md**。
4. **未把本机 2 行红当 CI 读数**（§2.2 已具名它是 CRLF 假红）。
5. 未裁定"撤不撤 362-r2 那条排除"——那是改门语义，属编排者/机主的口，本腿只交读数与引用者表。
