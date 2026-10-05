# evidence.md — 票 236 AC#6 的 2026-10-05 增量那一格（写腿 `236-r1`）

本腿只做两格：把 `tools/d22scan/selftest.go` 格式化到 gofumpt 干净（纯排版），并给出
**修完之后 tracked 分母还剩几枚、是不是全在 `.scratch`**（编排者摆 `Q-66` 缺的那枚数）。
写面＝`tools/d22scan/selftest.go` ＋本目录 `.scratch/wisp/probes/236/r1/**`，别的一枚未动。
工单票面 AC 框、`docs/reports/pending-and-issues.md`、`.github/workflows/ci.yml`、任何台件、
`tools/d22scan/main.go` 的九条禁令射程、`allowlist.txt` ＝本腿全程零字节。

原始读数码件在 `logs/`（只建不删，`issues/README` 规则 8；`/tmp` 下的归档与临时件一枚未删）。

---

## §0 起手锚与三发容量读数（这棵树枚数随锚点漂，所以每一把都带锚与时刻）

| 时刻（+08） | 事件 | 锚点（短） | tracked 全量 | tracked `.go` |
|---|---|---|---|---|
| `12:38:32` | 起手，`git archive HEAD` 到 `/tmp/arc236` | `c308be29` | 5867 | 924 |
| `12:42:13` | 动笔前复量（别的腿入库 2 枚件） | `0793e35e` | 5869 | 924 |
| `12:44:50` | **本腿 commit** | `b3eabab7` | 5869 | 924 |
| `12:45:43` | 取修后名册时复认（HEAD 已漂） | `3862e0e2` | — | — |
| `12:47:05` | 修后名册在漂移后的 HEAD 上复跑 | `14dc6f48` | 5906 | 924 |

- 起手 `git rev-parse HEAD` ＝ `c308be29f6f8a37121938aae7edefb3de4b42893`，分支 `dev`（`date` 现读 `2026-10-05 12:38:32 +0800`）。
- 归档完整性自证（每一发都做，尺＝`find <arc> -type f | wc -l` 对 `git ls-tree -r --name-only <sha> | wc -l`）：
  `/tmp/arc236-post`（锚 `b3eabab7`）＝ 5869 对 5869；`/tmp/arc236-pre`（锚 `0793e35e`）＝ 5869 对 5869；
  `/tmp/arc236-head`（锚 `14dc6f48`）＝ 5906 对 5906。两数相等＝归档取到了整树（含每枚 `go.mod`，
  那决定 gofumpt 的语言档；普查件 `269/a1/census.md` §9 第 2 条就是栽在只抽 `*.go` 上）。
- 尺本体：`gofumpt` 用本机 `$(go env GOPATH)/bin/gofumpt` ＝ `D:\work\base\gopath\bin\gofumpt`，
  `--version` 现读 **`v0.12.0 (go1.27.1)`**、rc=0（与普查件 §4 第 5 点、CI 日志里 `go install mvdan.cc/gofumpt@latest`
  下载的那版同版 => 本件的数与 CI 那台检出同尺）。
- 起手时工作树脏度（【禁】不是本腿的）：`git status --porcelain | wc -l` ＝ **697**；
  `git diff --cached --name-only | wc -l` ＝ **35**（逐枚看过，全在 `.scratch/wisp/probes/257/v1/` 下＝在飞腿 `257-v1` 的暂存料）。
  commit 1（那枚产码格式化）**零 `git add`**、只带显式 pathspec，commit 之后复量 `git diff --cached --name-only | wc -l` ＝ **35（原封不动）**，
  这一把是"没卷走别人 staged 文件"的凭据。commit 2／3 要新建与改名 `logs/` 下的件，用了**逐枚点名**的 `git add <路径>`／`git mv <旧> <新>`
  （`add -A`、`add .` 全程零枚，三笔 commit 每一笔都带 pathspec；两笔的载荷逐枚核过＝只含 `236/r1/**`，见 §6）。
- 在飞腿地界：`cmd/wisp`（`257-v1`）、`internal/agent/approval`（`259-r2`）＝**未读、未跑、未归因**。
  本腿跑过的 Go 仪器只有 `tools/d22scan` 那一枚模块自己那一套（§4）。

---

## §1 第 0 步：先验前提，再动笔（两把尺的原文读数）

### 尺 0-1（真脏，归档形状）＝**成立**

命令（`/tmp/arc236`，锚 `c308be29`，`12:38:32` 建、`12:39` 跑）：

```
mkdir -p /tmp/arc236 && git archive HEAD | tar -x -C /tmp/arc236
cd /tmp/arc236 && "$(go env GOPATH)/bin/gofumpt" -l tools/d22scan/selftest.go 2>&1; echo rc=$?
```

原文读数：

```
tools\d22scan\selftest.go
rc=0
```

**它真吐出了一行**，且 **rc=0**（＝"要格式化"，不是"解析不过"；解析不过才是 rc=2 那一形）。
`2>/dev/null` 没挂，stderr 一起接（`2>&1`），所以"空读数＝命令没跑到"这一坑在本发里不成立。

同一枚文件在 `b3eabab7` 之前那枚锚 `0793e35e` 的归档里复跑（`/tmp/arc236-pre`，`12:45`）＝
**同样吐出一行、rc=0**（尺见 §3 末条 2）。另加一发 stdin 形（保字节、不带路径）：
`git show 0793e35e:tools/d22scan/selftest.go | gofumpt -l` -> 打 `<standard input>`、rc=0 => 脏在字节里，不在路径形状里。

**行尾污染排除**（这是本机那把尺多算 5 枚假脏样的那一族坑，`A620`／`A615`）：同一枚文件的三形逐字节对齐——

```
CR 字节数：盘上 0 ／ HEAD blob 0 ／ 归档文件 0
字节数：  25589 ／ 25589 ／ 25589
md5：     bbc0bd946e1117e437c58939e867112e（三形全等）
```

尺＝`tr -dc '\r' < f | wc -c`、`wc -c`、`md5sum f <(git cat-file blob HEAD:f)`。
=> 这一枚的脏**不是行尾**，是真格式；本腿随后跑的全集名册（§3）也一律取归档形。

### 尺 0-2（它脏得是不是故意的＝有没有一枚自证把这件事钉进断言）＝**不是故意的载体**

派单给的三把尺逐把跑完，加两把我自己补的：

```
grep -rn 'selftest.go' --include='*_test.go' --include='*.go' tools .scratch scripts cmd internal docs
  -> 1 命中：tools/d22scan/selftestsamples.go:4（一句散文注释"…live in selftest.go; this file is the samples…"）
grep -rln 'selftest' tools/d22scan | head
  -> main.go / selftest.go / selftestsamples.go / selftest_test.go / d22scan.exe（二进制）
grep -rn -i 'gofumpt\|gofmt\|align' tools/d22scan/*.go
  -> 0 命中（该模块里没有任何一把尺判格式）
git grep -nE '(sha256|md5|wc -c|bytes|字节).{0,80}selftest\.go|selftest\.go.{0,80}(sha256|md5)' HEAD
  -> 2 命中，逐枚读原文（见下表 P1／P2）
git grep -nE 'selftest\.go:[0-9]' HEAD
  -> 4 枚件点名行号（下表 P3 族），逐枚复量行号未漂移（尺见下）
```

**"故意载体"排查表（每一处候选钉，逐枚读原文后的判定）**

| # | 出处（`file:line`） | 钉的是什么 | 是不是本枚格式脏的载体 | 凭据 |
|---|---|---|---|---|
| P1 | `tools/d22scan/selftest.go:108` ＋ `:177-241`（`readRosters`） | 名册读的是 **`main.go`** 的字节（`selfBansSourceFlag = "main.go"`），正则 `:117`／`:119`／`:125` | **否**——它读别的文件，本文件的排版不进名册 | `selftest.go:104-108` 原文；普查件 `266/a1/census.md:62` 同一读法（"读 `tools/d22scan/main.go` 自己的字节抽名册（`selftest.go:177-241`）"） |
| P2 | `tools/d22scan/selftest.go:359`（`hash()`） | 去重尺＝`c.file + NUL + c.src + NUL + c.allow`，对象是 `selfCases`（`selftestsamples.go` 的字面量），**不是本文件的字节** | **否** | 该行原文＋`auditSelfCases` `:347-354`；`selfCases` 定义在 `selftestsamples.go` |
| P3 | `selfSkeleton()`（`:365-379`）的两枚键 | 它们是投喂给夹具的种子路径（`docs/readings.md`／`internal/probe/roster.md`），被 `selftestsamples.go:387` 那枚 wantSilent 引用 | **否**——引用的是**路径字符串**与**值字符串**，两者逐字节未动；改的只是键与值之间的空格 | `:374-375` 那两行注释逐字写着"the fixture must seed them so those citations EXIST"＝要的是**存在性**，不是排版 |
| P4 | `.scratch/wisp/probes/161/r4/removed/sha256-before.txt`／`-after.txt` | `7ea3f1d8…` 钉的是 **`selftestsamples.go`**（另一枚文件），`overlay.json`/`overlay2.json` 替换的也是它 | **否**——与本枚无关 | 两枚 txt 原文（路径段就是 `selftestsamples.go`） |
| P5 | `.scratch/wisp/probes/212/v1/verdict-first-instance.md:8` | 记 `selftest.go` 的 md5 `bbc0bd94…`，语义＝"212-v1 那一腿全程对跟踪树零改动"（起止全等） | **否**——那是**那一腿自己**的零改动凭据（历史快照，无任何脚本复算它），不是"这枚文件必须是这串字节"的承诺 | 该行原文；全仓 `grep -rn 'bbc0bd94…'` 只有那一枚件指向它（复量见 §5 第 3 条，我给了新旧对照） |
| P6 | `.scratch/wisp/probes/212/r2/fix-and-readings.md:122-126`／`:228` | 逐字把本枚写成 **"具名一处过期项（非本腿造成、本腿未顺手改）"**，并点名差异就是 `selfSkeleton()` 里那两行的对齐 | **反向成立**：它是**要求这一修**的登记，不是禁止这一修的钉 | 两段原文我都读了，引句在 §7 第 4 条 |
| P7 | `.github/workflows/ci.yml` ＋ `probes/161/r5/attrib.sh` | 格式门只判"干不干净"，`TRACKED_DIRTY` 是**现算**的（`attrib.sh:315/352/411`），没有任何硬编码枚数期待本枚是脏的；`attrib.sh --self-test`（`:263-300`）递的是**文本行**，其 tracked 正控用的是 `tools/d22scan/main.go`，不是 `selftest.go` | **否** | `attrib.sh` 逐行读到 `:300`；本腿未执行它（§5 第 4 条具名） |
| P8 | selftest 表本身（`-self-test` 的 40 项双向） | 断言全在 `selftestsamples.go` 的字面量与 `main.go` 的名册上 | **否**——改前改后 40 项输出逐字同形（只差随机临时目录名，尺与读数＝§4 末条） | `logs/arc236-selftest-delta.txt`（1,434 字节、8 行、全部含 `d22scan-selftest-<随机>`；剔除该串后 0 行） |

**两把尺同时成立的结论**：真脏（尺 0-1，三形互证＋行尾排除）＋ 不是故意的载体（尺 0-2，八处候选逐枚读原文）
=> 允许动笔。若任一枚不成立，本腿按派单停手；这里两枚都成立，所以 §2 往下走。

**动笔前的额外一层保险（行号不漂）**：这次改动 `@@ -373,8 +373,8 @@` 是**等长替换**（2 行换 2 行），文件行数 640 对 640，
所以全仓那些 `selftest.go:<N>` 的引用**一枚都不会因本 commit 失效**。我逐枚复量了三处点名行：

| 引用出处 | 钉的行号 | 修后该行内容（`git show b3eabab7:… \| sed -n '<N>p'`） | 判定 |
|---|---|---|---|
| `issues/212-…-done.md:76` | `:332` | `holes = append(holes, fmt.Sprintf("tag %q has only an expect-silent sample - …` | 仍指同一句 |
| `probes/132/c2/census.md:154` | `:406` | `full := filepath.Join(root, filepath.FromSlash(rel))` | 仍指同一句 |
| `docs/evidence/s1/161-gates-accept-r2.md:148` | `:355` | `	return holes` | 仍指同一句（那格突变目标原样在位） |
| `issues/266-….md:39` | `:226-236`／`:548-554` | `:226`＝`// Every emission site must be accounted for by one of the two shapes above.`；`:548`＝`holes := auditSelfCases(selfCases, numbered, emitted)` | 两段区间未漂 |

---

## §2 改动逐行 before-after（对照表：删的每一行＝加的每一行的同句不同空格）

尺：`git show b3eabab7 -- tools/d22scan/selftest.go`，`git diff --numstat` 于 commit 前现读＝
**`2 2 tools/d22scan/selftest.go`**（2 增 2 删），`git show b3eabab7 | grep -c '^@@'` ＝ **1 枚 hunk**
（＝改动处数，多一枚就是顺带动了别处）。

| 行号 | 改前（`cat -A` 形，`^I`＝Tab，空格照写） | 改后 | 语义载荷 |
|---|---|---|---|
| 376 | `^I^I"docs/readings.md":` + **12 空格** + `"readings\n",` | `^I^I"docs/readings.md":` + **9 空格** + `"readings\n",` | 键串、值串、逗号逐字节同 |
| 377 | `^I^I"internal/probe/roster.md":` + **4 空格** + `"roster\n",` | `^I^I"internal/probe/roster.md":` + **1 空格** + `"roster\n",` | 同上 |

`cat -A` 原文（四行成对，删行在前、增行在后）：

```
-^I^I"docs/readings.md":            "readings\n",
-^I^I"internal/probe/roster.md":    "roster\n",
+^I^I"docs/readings.md":         "readings\n",
+^I^I"internal/probe/roster.md": "roster\n",
```

- 为什么是这两行：`selfSkeleton()` 里 `:374-375` 那两行注释把 map 字面量的**对齐段截成两组**。
  第一组最长键是 `tools/d22scan/allowlist.txt`（`:373`，仍在原对齐位上、本腿未动），
  第二组最长键是 `internal/probe/roster.md`（25 字符）=> gofumpt v0.12.0 要求第二组按**它自己那组**的最长键补齐，
  于是 `docs/readings.md` 少 3 空格、`internal/probe/roster.md` 少 3 空格。合计少 **6 字节**（25589 对 25583，实测）。
- **改法不是我手敲的**：`gofumpt -w` 对一份副本（`logs/arc236-selftest-orig.go.pristine` -> `logs/arc236-selftest-fixed.go.pristine`）跑完，
  `diff -u` 只有上面那一个 hunk，才把结果落到工作树 => 判它脏的那把尺就是把它改干净的那把尺，不留"我理解格式"的余地。
- 逐枚改动行是否全在格式＝**是**（2 对 2，成对同句；键/值/字符串字面量/注释/行数/行尾全不变，`0 枚 CR` 改后复量仍为 0）。
- 未做的事：未顺手调整第一组对齐、未动 `:373`、未动该文件任何注释文字（注释里那句"Ticket 212's wantSilent case cites these two"逐字保留）、
  未动其余 19 枚台件、未动 `main.go` 的九条禁令与 `allowlist.txt`。

---

## §3 那枚数（关键读数）：修完之后 tracked 分母还剩几枚（名册＋口径＋锚点＋时刻）

**口径（三件必须一起说，缺一把就答错题）**

1. **形状**＝归档（`git archive <sha> \| tar -x`），不是本机工作树（本机那把尺被 CRLF 污染，`A620` 现场复现过：本机 32 枚、归档 19–20 枚）。
   两数对齐凭据＝归档 `find -type f` 与 `git ls-tree -r --name-only | wc -l` 相等（§0）。
2. **分母**＝该锚点 tracked 全集 `git ls-tree -r --name-only -z <sha> | grep -zE '\.go$'` ＝ **924 枚**（`0793e35e`／`b3eabab7`／`14dc6f48` 三锚同为 924），
   逐枚前缀归档路径后 `xargs -0 gofumpt -l`；**stdout 与 stderr 两路都要收**（解析错误走 stderr，漏收就少数一枚，普查件 §8 第 2 条栽的就是这个）。
3. **锚点**＝**本腿那笔 commit 之后**：`b3eabab7aefe637918b1f53dc4adaf0d88d231a2`（`12:44:50` 入库），取数时刻 **`2026-10-05 12:45:58 +0800`**。
   取数时 HEAD 已漂到 `3862e0e2`（`12:45:43` 现读），所以我在漂移后的 `14dc6f48`（`12:47:05`）又复跑一发对差（下条 4）。
   本腿后面两笔（commit 2／3）里有过一次我自己造的分母污染与自纠（详见 §6 与 §7 第 8 条），故
   **交付数以 §8 第 6 条那一发为准**（锚 `00629ea9`，含全部三笔，13:07:09 现读：**19 枚、非台件 0 枚、与本名册 `diff` 空**）。
   两枚锚给出的答案相同（19），差别只在中间那枚 `cfd97636` 上短暂是 20（我自己顶上去的，已改名纠正）。

**答案：19 枚。每一枚都落在 `.scratch/wisp/probes/**` 里。非台件＝0 枚（没有"哪几枚不在、归谁"这一格，那一格今天为空）。**

| 口径尺 | 修前（锚 `0793e35e`，`12:45:58` 同把尺） | 修后（锚 `b3eabab7`，同一把尺） |
|---|---|---|
| stdout 枚数 | 19 | **18** |
| stderr 枚数（解析错误行） | 1 | 1 |
| 去重后不同文件 | **20** | **19** |
| 其中在 `.scratch/wisp/probes/**` | 19 | **19** |
| 其中不在 `.scratch` | **1**（`tools/d22scan/selftest.go`） | **0** |
| `xargs` 聚合退码 | 123 | 123 |
| tracked `.go` 分母 | 924 | 924 |

差集尺（`diff` 修前名册 对 修后名册，件在 `logs/census-pre-roster.txt`／`logs/census-post-roster.txt`）＝**恰好一行**：
`20d19 < tools/d22scan/selftest.go`。=> 本枚改动对分母名册的影响是**减一且只减这一枚**，没有多咬别处。

**修后名册逐枚（19 枚，锚 `b3eabab7`；归属＝路径里的票号，`票面在`＝`.scratch/wisp/issues/<NN>-*.md` 真存在，尺＝逐枚 `ls \| grep -c "^<NN>-"`＝全 1）**

| # | 路径 | 归属票 | 票面在 | 脏因（沿普查件 §5 尺一＋本腿复认） |
|---|---|---|---|---|
| 1 | `.scratch/wisp/probes/33/p1/q1/main.go` | 33 | 有 | 纯格式（台件） |
| 2 | `.scratch/wisp/probes/33/p1/q2/main.go` | 33 | 有 | 纯格式 |
| 3 | `.scratch/wisp/probes/33/p1/q3/main.go` | 33 | 有 | 纯格式 |
| 4 | `.scratch/wisp/probes/163/a1/main.go` | 163 | 有 | 纯格式 |
| 5 | `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | 174 | 有 | 纯格式 |
| 6 | `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | 183 | 有 | 纯格式（突变台件副本） |
| 7 | `.scratch/wisp/probes/185/c1/mut/fs_broken.go` | 185 | 有 | **解析错误** `4:1: imports must appear before other declarations`（单发 rc=2；来自 stderr 那一路） |
| 8 | `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | 185 | 有 | 纯格式 |
| 9 | `.scratch/wisp/probes/197/r1c/pre/subagent_197.go` | 197 | 有 | 纯格式（产品码改前逐字节快照，`A457` 记过其承重） |
| 10 | `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go` | 197 | 有 | 纯格式（同上） |
| 11 | `.scratch/wisp/probes/212/v1/mut/main-noq9.go` | 212 | 有 | 纯格式 |
| 12 | `.scratch/wisp/probes/220/r1/prechange/prechange220_probe_test.go` | 220 | 有 | 纯格式（改前快照） |
| 13 | `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | 222 | 有 | 纯格式 |
| 14 | `.scratch/wisp/probes/224/v2/dialect_probe_test.go` | 224 | 有 | 纯格式 |
| 15 | `.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go` | 235 | 有 | 纯格式 |
| 16 | `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go` | 241 | 有 | 纯格式（文件名自称坏样本） |
| 17 | `.scratch/wisp/probes/241/v1/probe_v1_readings_test.go` | 241 | 有 | 纯格式 |
| 18 | `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` | 259 | 有 | 纯格式（目录名自称 gofumpt 负控；`50e0299c` 10-05 12:09 入库＝普查件名册第 20 枚那一枚） |
| 19 | `.scratch/wisp/probes/263/v1/src/main.go` | 263 | 有 | 纯格式（同目录 30 字节 `go.mod` 决定语言档，整树归档才读得到它） |

`tools/d22scan/**` 名下：**0 枚**（该模块在本锚点上归档形与本机形都干净，尺见下）。

补充读数 1（锚点漂移下的同一把尺）：`14dc6f48`（`12:47:05`，比我的 commit 又多 37 枚 tracked 件、`.go` 仍 924）
=> 名册 **19 枚、非 `.scratch` 0 枚**，与 `b3eabab7` **`diff` 空**（`logs/census-head-roster.txt`）。
=> 这枚 19 **不是只在"我这笔 commit 那一瞬间"成立**；它在那之后的 HEAD 上复认了一次。注意：它仍**不是常量**：
`259/r1` 那种负控入库就把数往上顶（`A620` 的 19->20 正是这么发生的），编排者引用它必须带锚点＋时刻。

补充读数 2（gofumpt 单枚点名，两形分开报，别混）：

| 形 | 锚／状态 | 命令 | 读数 |
|---|---|---|---|
| 归档形 | `0793e35e`（修前） | `/tmp/arc236-pre` 里 `gofumpt -l tools/d22scan/selftest.go` | `tools\d22scan\selftest.go`、rc=0（脏） |
| 归档形 | `b3eabab7`（修后） | `/tmp/arc236-post` 里同命令 | **空**、rc=0（干净） |
| 本机形 | 工作树，修后 | `gofumpt -l tools/d22scan/selftest.go` | **空**、rc=0（干净） |
| 本机形 | 工作树，修后（整模块目录走查） | `gofumpt -l tools/d22scan` | **空**、rc=0（该模块本机形全干净） |
| stdin 形 | pre blob 逐字节 | `git show 0793e35e:tools/d22scan/selftest.go \| gofumpt -l` | `<standard input>`、rc=0（脏在字节里） |

补充读数 3（那一格不在我写面、但编排者摆 `Q-66` 要用的"链头红落在哪一步"，我只量不改）：
把 CI `ci.yml:168` 那条命令原样在**修后归档**里跑（`cd /tmp/arc236-post && gofumpt -l . tools/d22scan tools/mockllm`，`12:47:26`）＝
**rc=2**、stdout 18 行、stderr 1 行、非 `.scratch` **0 枚**、名册与 tracked 全集那 19 枚 `diff` 空。
=> 修完之后 `:168` **照红**，但红的理由**从"19 枚台件＋1 枚真违规"变成"19 枚全台件"**；
`rc=2` 那一路是 `185/c1/mut/fs_broken.go` 的解析错误顶的（不是枚数），所以"红形"仍是 `bash -e` 当场掐死那一步那一族
（普查件 §2 乙 第 4 条(一) 的机理，本腿复认）。`ci.yml:184`（(A) 尺）本腿**未执行**，它与 `:168` 吃同一批字节这一条我在同一棵归档上复算过名册（19 枚＝同一批），
它的退码映射（`attrib.sh:342-345` 命中 `A_RC>1` 走 exit 2）是按脚本文本读的，不是跑出来的 => 具名进 §5 第 4 条。

---

## §4 门禁读数（自己跑、带时刻；改前改后各一发）

| 门 | 命令 | 修前时刻／读数 | 修后时刻／读数 |
|---|---|---|---|
| D22 静态门 | `sh scripts/d22scan.sh` | `12:40:49` **clean**、rc=0 | `12:48:46` 与 `12:50:10`（复跑取全量输出）**clean**、rc=0；逐行自报 `bans #1-5 internal/=228`、`cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 design/=39`、`frontend/=85`、`internal/=512`、`cmd/=103` |
| 路径长度预算 | `bash scripts/check-path-length-budget.sh --with-self-test` | `12:41:28` **VERDICT GREEN**、rc=0（`denominator read: 5867 tracked paths`、over-budget=57、roster=57、not in roster=0、longest=180） | `12:49:13` 与 `12:50:10` **VERDICT GREEN**、rc=0（同一组数；分母随锚漂到 5906 那一发在本腿 commit 之后，见 §0） |
| 该模块 vet | `cd tools/d22scan && go vet ./...` | `12:41:34` rc=0（零输出） | `12:43:13` rc=0（零输出） |
| 该模块 test | `go test -count=1 ./...` | `12:41:35` **ok** `github.com/CarlosShao/wisp/tools/d22scan` 23.887s、rc=0 | `12:43:14` **ok** 19.087s、rc=0 |
| 该模块自证 | `go run . -self-test` | `12:40:30` rc=0，`roster read from main.go = 9 numbered ban(s) … + 1 finding type(s); 40 cases, 10 tag(s) covered, both directions required per tag`，`clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)` | `12:43:34` rc=0，**同一行逐字相同**（9 ban／40 cases／10 tags／20+20／clean 行一致） |
| 格式尺（点名那一枚） | `gofumpt -l tools/d22scan/selftest.go` | 归档形打印该文件（§1 尺 0-1） | **归档形与本机形各报一次＝两形都空**（表在 §3 补充读数 2） |

**项数没掉**：派单按 `A620` 写的是"34 项双向全过"，我改前的基线（`12:40:30`，还没碰过码）就是 **40 项**、改后仍 **40 项**、rc=0。
=> 那枚 34 是过期登记（票 212 那批加了 ban #9 的样本，`212/r2` 现跑名册行写的是 `37 cases`，我这一发是 `40 cases`；件＝`docs/evidence/s1/212-comments-phantom-citation-v2.md:190`）。
判据"项数掉了就是被本腿改坏"不成立（未掉），差异**具名进 §5 第 1 条**并回报编排者，我不动台账那行。

**自证输出等价尺**（这是"语义零动"最硬的一发）：把改前（`12:40:30`）与改后（`12:43:34`）两发的 `-self-test` **全量 stdout+stderr** 落盘对差：

```
diff logs/arc236-selftest-before.txt logs/arc236-selftest-after.txt -> 8 行 / 1,434 字节
剔除含 d22scan-selftest-<随机> 的行之后 -> 0 行
```

=> 两发的实质读数**逐字相同**，差异全在 `os.MkdirTemp` 的临时目录名里（那两行是 unparseable 那条 finding 打印的夹具绝对路径）。
两发都是 69 行、40 枚判据行（含 ` OK `）、38 枚 `ban` 前缀＋2 枚 `type unparseable`、26 枚 `pinned because` 注行。

**未跑清单（守地界）**：`cmd/wisp` 的测试、`internal/**`（含 `internal/agent/approval`）、`frontend/**` 与 `design/**`
（后两枚连读都没读）、`go build`（除该模块 self-test 里那条既有 `go build -o <TempDir>` 由 `go test` 自己带出的）、
`attrib.sh`、`gh`／任何 push。

**终尺（本腿四笔全部入库之后复跑一遍整套，锚 `c216fe4b`，13:10:53 至 13:11:58 +0800）**：

| 门 | 时刻 | 读数 |
|---|---|---|
| `sh scripts/d22scan.sh` | `13:10:53` | **clean**、rc=0（`no D22 ban violations`，逐行自报的 live scope 与 §4 改后那一发同数） |
| `bash scripts/check-path-length-budget.sh --with-self-test` | `13:11:25` | **VERDICT GREEN**、rc=0 |
| `go vet ./...`（`tools/d22scan`） | `13:11:31` | rc=0、零输出 |
| `go test -count=1 ./...`（`tools/d22scan`） | `13:11:32` | **ok** 25.904s、rc=0 |
| `go run . -self-test` | `13:11:58` | rc=0、判据行 **40** 枚（含 ` OK `）、`clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)` |

=> 四笔入库后没有把任何一门改红；`-self-test` 那 40 项与本腿起手基线（`12:40:30`，改前）逐字同形。
终尺输出落档＝`logs/final-gate-d22scan.txt`、`logs/final-gate-pathlen.txt`、`logs/final-selftest.txt`（第五笔）。

---

## §5 判不动／量不到／与我无关但影响读数的（逐枚具名）

1. **派单里"34 项"那枚数过期**：实测基线 40 项（改前就 40）。我**没改**、也**不能改**台账 `A620` 那行；这一格归编排者。
   凭据尺＝§4 那发的名册行原文。
2. **枚数是移动靶**：本腿 7 分钟内看到 HEAD 走 `c308be29 -> 0793e35e -> b3eabab7(我) -> 3862e0e2 -> 14dc6f48`，
   tracked 全量 5867 -> 5869 -> 5906（`.go` 分母全程 924）。=> §3 那枚 **19** 只对 `b3eabab7` 与 `14dc6f48` 两个锚成立；
   再有腿入库 `.scratch/**` 下的 `.go` 台件就会往上顶。引用它必须带锚点＋时刻（这一条普查件 §5 末条已现场证明过一次）。
3. **我这枚改动让一处历史 md5 不再可复现**（唯一一处，主动具名）：`probes/212/v1/verdict-first-instance.md:8` 记
   `selftest.go` = `bbc0bd946e1117e437c58939e867112e`（那是 212-v1 那腿"起止全等＝我没动跟踪树"的凭据，当时为真、现在仍是对**那次运行**的真陈述，
   没有任何脚本复算它）。改后新值：`b56c7636473283cf453fee9178e53498`（盘上与 `b3eabab7` 的 blob 全等，md5 现读两形）。
   旧值覆盖区间＝`5e8748b3`（10-03 首入库）至 `0793e35e`；新值自 `b3eabab7` 起。
   => 谁日后按那枚旧 md5 去"复核 212-v1 的零改动"会读到不一致，这一对照表就是给那条读数留的桥。我不改那枚件（`只建不删`／不编辑别人已提交件）。
4. **`attrib.sh` 本腿未执行**：它 `:137-140` 要从 `$(go env GOPATH)/bin` 取 gofumpt、`ci.yml:184` 那步依赖 `:170` 的 `go install`。
   我用等价命令复现了它的分母（tracked 全集 × gofumpt -l，§3 那一把），所以"它读到同一批 19 枚"是**实测**，
   而"它会以 exit 2 拒答（`fs_broken.go` 的 rc=2 把 `xargs` 聚合顶到 123）"是**按脚本 `:342-345` 文本读**的，不是跑出来的。
   要定案需真发 run 或授权我跑那把脚本本体——本腿按派单"不许改 `ci.yml`、不许加 `if:`"停在量不到这一格。
5. **CI 那台的真实颜色我量不到**：本腿不 push、不读 run（派单未给 `gh` 尺）。§3 补充读数 3 给的是**同一把尺在归档形状上的复算**，
   与 `ci.yml:170` 的 `gofumpt@latest` 存在版本漂移风险（本机 v0.12.0；上游改版会动"要格式化"集合，普查件 §7 第 5 点同一条坑）。
6. **形 I／形 II／形 III 的裁定不属本腿**：本腿只把那一枚非台件真违规修掉，
   台件一枚未动（形 II 未被裁）、`ci.yml` 零字节（形 I 要 `Q-66`）、`if:` 一枚未加（形 III 未与 `Q-66` 分离落地）。
   "剔完台件还剩几枚"这一问今天的答案是 **0 枚**，那正是形 I 单独做能不能让 `:168` 变绿的关键料——
   但 `rc=2` 那一路（`fs_broken.go` 的解析错误）不在"格式脏"这一档里，**它会不会让形 I 之后仍然红，属编排者裁形时另要的一枚数**，
   本腿只把两档分开放在这里（18 枚纯格式＋1 枚解析错误），不替那一裁下结论。
7. **`docs/reports/**`／票面 AC 框**：一字未动；翻勾与摆 `Q-66` 归编排者。本件不新立 `A##`/`Q##`，也不改既有行。

---

## §6 commit 名册（只 commit，绝不 push；每笔都带显式 pathspec）

| 笔 | sha（短） | 时刻 | 内容 | pathspec | numstat |
|---|---|---|---|---|---|
| 1 | `b3eabab7` | `2026-10-05 12:44:50 +0800` | `style(161/236-r1): tools/d22scan/selftest.go 格式化到 gofumpt 干净（纯排版，语义零动）` | `-- tools/d22scan/selftest.go` | `1 file changed, 2 insertions(+), 2 deletions(-)`；`^@@` 计数＝1 |
| 2 | `cfd97636` | `2026-10-05 12:58:58 +0800` | `docs(evidence/236-r1): 修完那一枚非台件真违规之后的名册与读数（Q-66 缺的那枚数＝19 枚、全在 .scratch）` | `-- .scratch/wisp/probes/236/r1` | 16 枚新建（`evidence.md`＋`msg-*.txt`＋`logs/**`），0 删 0 改 |
| 3 | `00629ea9` | `2026-10-05 13:06:54 +0800` | 自纠：两枚 `.go` 副本改成 `.go.pristine`（只改后缀、不动字节）＋本件 §0／§2／§3／§6／§7／§8 更正 | `-- .scratch/wisp/probes/236/r1` | 4 枚变动＝2 枚 rename（相似度 100％）＋`evidence.md`＋`msg-selffix.txt`；77 增 14 删 |
| 4 | 本笔 | 落档时刻见 commit 本体（`git log -1` 现读） | §8 第 6 条的真实读数落档（`logs/final-roster.txt` 等 4 枚）＋§3 锚点段补一句"交付数取哪一发" | `-- .scratch/wisp/probes/236/r1` | 4 枚新建 `.txt`＋本件更正（零产码、零 rename） |

- commit 1 的父＝`0793e35e`（现读 `git rev-parse b3eabab7^`）；commit 2 的父＝`6c96a425`（现读 `git rev-parse cfd97636^`，
  那是 257-v1 在 12:56:34 的结案笔；12:44 到 12:58 之间 HEAD 走了 4 枚，全是别人的入库，锚漂是常态）。
- **我这枚 commit 2 自己造过一次分母污染（自纠，具名）**：为留"改前／改后逐字节对照"，我把两枚 Go 源副本
  以 **`.go` 后缀**放进了 `logs/` 并入库 => tracked `.go` 分母从 924 顶到 **926**，
  且那枚改前副本（路径以污染锚 `cfd97636` 为准＝`arc236-selftest-orig.go`，commit 3 起改名 `arc236-selftest-orig.go.pristine`，
  下面两处的历史路径同理）在归档形上被 `gofumpt -l` **点名**（现读：打印该路径、rc=0）。
  => 我一边交付"非台件 0 枚"，一边亲手往名册里加了一枚新的脏台件——这正是本票要治的那一形（`A620` 的 19->20 同源）。
  修法走本仓既有先例（票 267 `pristine/models.go` -> `models.go.pristine`；普查件 §2 乙-1 尺＝在册 17 枚 `.pristine`）＝**只改后缀、不动字节**，
  两件 md5 逐枚全等：改前改后都 `bbc0bd946e1117e437c58939e867112e`（orig）与 `b56c7636473283cf453fee9178e53498`（fixed）。
  改后复尺：`git ls-files -z '*.go' | tr -dc '\0' | wc -c` ＝ **924**、本目录下 `.go` 在册 ＝ **0 枚**（终尺见 §8 末条）。
  这一枚**未被删**（`issues/README` 规则 8 只建不删）、改名前后字节全等；而 §3 那枚 **19** 是在 `b3eabab7` 与 `14dc6f48` 两锚上取的，
  当时我的 `.go` 副本还没入库 => **那枚数不受本次污染影响**；受影响的是"12:58 之后在新锚上复跑"的人，故此段把机制写全。
- 那 35 枚别人的 staged 料在 commit 1 之后仍在索引里（**12:44 现读** `git diff --cached --name-only | wc -l` ＝ 35，逐枚名全在 `probes/257/v1/`）
  => 没有"先 `git add` 再裸 commit"，也没卷走别人。这 35 枚此后由 **257-v1 自己在 `6c96a425`（12:56:34）结的案**，
  索引里今天只剩 2 枚＝本腿 commit 3 那两枚改名（现读，见下条），不是别人的料被动过。
- commit 2 用 `git add <逐枚点名的路径>` ＋ `git commit -F msg -- <同一批点名 pathspec>`
  （新件非入库不可；两把都点名到枚，未用 `-A`／`.`；尺＝commit 2 的 `--name-only` 16 行**全部**在 `236/r1/**` 下，
  与当时索引里那 35 枚 `257/v1/**` 交集＝**0 枚**，逐枚 `grep -c 257` 现读那两笔（1 与 2）都是 0）。
- 未 push（`git status -sb` 本腿全程只 commit）；`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`add -A`／`add .`／`rm` 零次；
  `git mv` 2 次（只动后缀，commit 3 内，逐枚 md5 全等）。
- 临时件全部只建不删：`/tmp/arc236`、`/tmp/arc236-pre`、`/tmp/arc236-post`、`/tmp/arc236-head`、`/tmp/arc236-fix`（五份整树归档）、
  `/tmp/arc236-selftest-{orig,fixed,before,after,delta}.*`、`/tmp/census-{pre,post,head}.{out,err}`、`/tmp/golist-*.z`、
  `/tmp/gate-{d22scan,pathlen}-post-full.txt`、`/tmp/ci168-post.roster`；其中关键件已复制进 `logs/`，
  `logs/` 里被当"可重跑凭据"引用的东西本腿一枚不许删。
- 仓内未落**整树归档**副本（尺：`git status --porcelain -- tools/d22scan` 在 commit 1 前后都只那一枚 M；
  工作树里除本腿那枚产码文件外无新增改动；`/tmp/arc236*` 五棵树全部留在 `/tmp`）
  ——但 commit 2 确实把两枚**单件**副本带进过仓库目录（就是上面自纠那一格），本行的原口径"一枚未落"是错的，已按实改写。
- **`logs/` 两枚件的带内符号具名（不是我写的）**：`logs/arc236-selftest-before.txt` 与 `logs/arc236-selftest-after.txt` 各有 1 行含一枚
  带圈小写字母 a（U+24D0，落在这件的禁带里，本件不用字形本体、只点名码位），
  那是 `-self-test` 自己打印的 `selftestsamples.go` 注行原文（票 260 那一族的形号）。这两枚是**仪器逐字节输出**，
  按"读数码件保真"我一枚不改；本腿自己写的文字（本件＋`msg-fix.txt`）带内符号＝0 枚
  （尺＝`grep -cP '[\x{2190}-\x{2BFF}]' evidence.md msg-fix.txt`＝各 0）。

---

## §7 自我对抗（真改过的东西，逐条给"原来怎么读、现在怎么读、尺是什么"）

1. **"项数掉了"我一度判成本腿改坏**：基线 `-self-test` 打 40 项而派单写 34，我第一反应是自己动坏了；
   实际**改前那一发（12:40:30，尚未碰码）就是 40**。=> 改判为"派单里的 34 是过期登记"（§5 第 1 条）。
   尺＝改前改后两发的名册行逐字对比（同一行、同一数）。
2. **`grep -c '^d22scan -self-test: ban'` ＝38 我一度当成"少了 2 项"**：那把尺数的是**前缀**，不含 2 枚 `type unparseable`（finding type 不是 ban）。
   => 换成 `grep -c ' OK '`＝**40**＝`len(selfCases)`，与 `selftest_test.go:54` 那枚断言同尺同值。
   这条如果没纠，我会把"项数对不上"错报给编排者。
3. **名册我先只收了 stdout 的 18 枚**：第一版数成"修后剩 18 枚"，那**漏了 stderr 上那枚解析错误**（`fs_broken.go` 走 stderr）。
   => 补成 stdout＋stderr 两路、按文件去重＝**19 枚**，并与普查件 §5 尺一的名册逐枚对齐（差集只有 `selftest.go`）。
   尺＝`logs/census-post.out`（18 行）＋`logs/census-post.err`（1 行）＋`logs/census-post-roster.txt`（19 行）。
4. **"它可能是故意的载体"我一度按字面接受 `212/r2` 那句"版本相依"**：那枚件（`:122-126`、`:228`）读起来像"别动它，等 CI 版本定"。
   逐字重读后确认它的自我定性是 **"具名一处过期项（非本腿造成、本腿未顺手改）"＋"只上报，不并入"**，是**要求这一修**的登记，不是禁止这一修的钉
   （它自己写的差异位置＝`selfSkeleton()` 里那两行对齐，与我这发 `-d` 读数逐字相同）。=> 不改停手，并把 P6 写进 §1 排查表。
5. **行尾那把我一度准备用本机尺交差**：派单要求归档形，我仍先在工作树跑了一次 `gofumpt -l`（脏/干净混着 CRLF 那一族假读数）。
   => 改成三形互证（盘上／blob／归档 md5 全等＋CR 各 0 枚）后才把"真脏"当成立，名册一律取整树归档。
   这条与普查件 §9 第 1 条同族（工具没骗我，是我的读法骗我）。
6. **归档我第一发抽在旧锚**：`/tmp/arc236` 建在 `c308be29`（12:38），而 HEAD 在 12:42 已漂到 `0793e35e`。
   修前对差若拿这两发比就会把"别人入库的 2 枚件"算进我的账。=> 重新归档 pre／post **成对**（`0793e35e` 对 `b3eabab7`），
   两锚 tracked 全量与 `.go` 分母逐枚对齐（5869 对 5869、924 对 924），差集才敢写成"恰好一行"。
7. **`ci.yml:184` 的退码我一度写成实测**：那是脚本语义推断。=> 降级为"名册同一批＝实测；退码映射＝按文本读"，
   具名进 §5 第 4 条，不冒充跑过的读数。
8. **我自己造了一枚新的脏台件（这一条是本腿的缺陷，不是读数误差）**：commit 2 把改前/改后的两份 Go 源副本
   以 `.go` 后缀放进 `logs/` 并入库 => tracked `.go` 分母 924 -> **926**，其中那份**改前**副本就是修前的脏字节本体，
   在归档形上被 `gofumpt -l` 点名（现读打印该路径、rc=0）。
   原来怎么读："我只是复制读数码件，且都在 `.scratch` 下，不碰别人写的东西＝合规"。
   现在怎么读：**"落在 `.scratch` 下"不等于"无害"**——本票的命题就是"入库的 `.go` 台件会顶格式门的分母"，
   我一边交"非台件 0 枚"、一边把名册从 19 顶回 20，是同一种病的第三次发作（前两次＝票 185 的 `fs_broken.go`、`259/r1` 的负控）。
   尺是什么：`git ls-tree -r --name-only <锚> | grep -cE '\.go$'`（924/926 两发）＋`git archive <锚> | tar -x` 后
   `gofumpt -l .scratch/wisp/probes/236/r1/logs/`（打印那一枚）。
   处置：commit 3 按本仓既有先例（票 267 `models.go` -> `models.go.pristine`）**只改后缀、不动字节**，md5 逐枚全等；
   未删任何件，`logs/` 里两份对照原件逐字节仍在（改名后仍可 `cmp` 回 §2 那四行）。
   我为什么没在 commit 2 之前发现：我当时的自检只扫了"符号带／占位符／pathspec"，**没有把"新入库的 `.go` 会不会进那把我自己交付的尺的分母"**
   列成一条尺——这条恰好是本件 §3 口径第 2 条写明的东西。=> 已把该尺补进 §8 末条，交付前对新锚复跑一次。

---

## §8 交件判语

1. **第 0 步两把尺都成立**：真脏（归档形单发打印该文件、rc=0；三形 md5 全等、CR 各 0 枚，排除行尾假读数）；
   不是故意的载体（八处候选 P1–P8 逐枚读原文，零枚把本枚的字节／哈希／"gofumpt 会点它"钉成断言；
   selftest 的名册读 `main.go`、去重尺读 `selfCases`，两处都在本文件之外）。=> 动笔合规。
2. **格一完成**：`tools/d22scan/selftest.go` 在归档形与本机形都 `gofumpt -l` 空；改动 2 增 2 删、1 枚 hunk、
   逐枚全在格式（成对同句、只差键值间空格），文件行数 640 未漂、行尾 0 枚 CR、`d22scan` 那一套四读数全绿且自证输出实质零漂移。
3. **格二那枚数＝19**：锚 `b3eabab7`（本腿 commit）、时刻 `2026-10-05 12:45:58 +0800`、口径"整树归档 × tracked 全集 924 × gofumpt v0.12.0 × stdout+stderr"；
   **19 枚全部落在 `.scratch/wisp/probes/**`，非台件 0 枚**（"不在的是哪几枚、分别归谁"＝今天这一格为空）；
   在漂移后的 HEAD `14dc6f48`（12:47:05）复认一次，名册 `diff` 空。
   => 交给 `Q-66` 的那半句现在是：**"形 I 剔除台件之后，非台件剩 0 枚"**（修之前是 1 枚），
   而**"形 I 之后 `:168` 会不会绿"本腿不替它答**——因为 `185/c1/mut/fs_broken.go` 那枚走 stderr 的解析错误（单发 rc=2）不属"格式脏"那一档，
   它把 `xargs` 聚合顶到 123 这一条我在修前修后各测一次（都是 123）。档与不档的分别是：18 枚纯格式＋1 枚解析错误，共 19 枚，全在台件目录里。
4. **禁区自证**：`.scratch/**` 别人的台件零字节（我那三笔 commit 的 pathspec 只到 `-- tools/d22scan/selftest.go` 与 `-- .scratch/wisp/probes/236/r1`）、`ci.yml` 零字节、
   新增 `if:`/`continue-on-error` 零枚、未删任何件（commit 3 是 `git mv` 改后缀，md5 逐枚全等）、`main.go` 九条禁令与 `allowlist.txt` 零字节、
   普查名册里其余 19 枚一枚未顺手修、票面与台账零字节、`frontend/**`／`design/**` 连读都没读。
5. **本腿不判**：形 I／II／III 的取舍、`Q-66` 要不要摆与怎么措辞、`:184` 那步的真实退码、票面 AC 框翻勾、台账 `A##` 追加——全归编排者。
6. **交付前那把我漏掉的尺，补在这里（新锚复跑，已跑出数）**：在**含 commit 3 的锚** `00629ea9ac0d7366d20525837aed8cff778a96cf`
   （13:06:54 入库）上重跑 §3 那把尺（整树归档 × tracked 全集 × gofumpt `-l` × stdout＋stderr），现读时刻 **`13:07:09 +0800`**：

   | 尺 | 读数 |
   |---|---|
   | 归档完整性（`find -type f` 对 `git ls-tree -r --name-only \| wc -l`） | 5936 对 5936（相等＝整树取全，含每一枚 `go.mod`） |
   | tracked `.go` 分母 | **924**（回到 commit 1 那一发的数） |
   | `xargs` 聚合退码 | 123（＝`fs_broken.go` 那一枚 rc=2 顶的，与修前修后同形） |
   | 名册枚数（stdout 18 ＋ stderr 1，去重） | **19** |
   | 其中不在 `.scratch/wisp/probes/**` | **0** |
   | 与 `b3eabab7` 那发名册的 `diff` | **空**（同一批 19 枚） |
   | 我这枚 commit 带进 tracked 集的 `.go` 载荷 | **0 枚**（`git ls-tree -r --name-only <锚> \| grep 236/r1 \| grep -cE '\.go$'`） |
   | 单枚点名（`gofumpt -l tools/d22scan/selftest.go`，该锚归档内） | **空**、rc=0（干净） |

   对照读数（那把我漏掉的尺为什么必须补）：同一把尺在**污染锚** `cfd97636`（我的 commit 2）上＝
   分母 **926**、名册 **20 枚**，多出来的一枚逐字是 `.scratch/wisp/probes/236/r1/logs/arc236-selftest-orig.go`
   （该路径只存在于污染锚 `cfd97636`；commit 3 起它叫 `arc236-selftest-orig.go.pristine`，两枚名都以各自锚为准，别拿今日盘上尺去复核它）
   （尺与名册＝`logs/polluted-roster-at-cfd97636.txt`，20 行）。交付数以 `00629ea9` 这一发为准：**19 枚、非台件 0 枚**。
   落档件＝`logs/final-roster.txt`（19 行）、`logs/final-gofumpt-stdout.txt`（18 行）、`logs/final-gofumpt-stderr.txt`（1 行）、
   `logs/polluted-roster-at-cfd97636.txt`（20 行，那一次污染的现场名册）。
