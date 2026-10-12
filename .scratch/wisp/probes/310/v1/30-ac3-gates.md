# `30-ac3-gates.md` — `310-v1` 对票 310 `AC#3`（门禁与越界）的判语

**判语一句**：`AC#3` **成立**（四把门禁我自己逐把跑到 rc=0、totals 行跑前⇄跑后逐字未变、逐笔越界尺十笔全在允许清单内），⚠ 带三处具名缺口（格式那格是**代用尺**、`Q-310-r1-4` 的归因**方向写反了**、读数件里 28 枚从未被交付件点名）。

台件＝`ac3-gates-310v1.sh`；主读数件＝`logs/30-ac3-run.txt`（21 行逐字）。⛔ 抄腿的 rc——下面每一个 `rc=` 都是我这一刻的读数。

## 1. 四把门禁（逐把我自己看的复跑）

| 把 | 命令（逐字） | 我的 rc | 逐字读数 | 件 |
|---|---|---|---|---|
| ① D22 扫描 | `sh scripts/d22scan.sh` | **`d22scan_rc=0`** | 末行 `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 i…`；stdout 251 行、stderr **0** 字节；★同一发里 `GOTMPDIR=` 命中＝**0**、`DEGRADED` 命中＝**0**（第五枚入口没有那一步＝`10-ac1-verdicts.md` §③ 的现量） | `logs/30-d22scan.txt`／`.err` |
| ② 路径长度预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | **`pathlen_rc=0`** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；17 行、stderr 327 字节（那句⛔ 报错、是它的自述行）；★我**建自己任何一枚新件之前**先跑过一次＝同一 rc=0（件＝开工第一发） | `logs/30-pathlen.txt`／`.err` |
| ③ census totals | `bash scripts/portable-tests.sh --scope=census` **跑前取一次、跑后再取一次作差** | **`census_before_rc=0`／`census_after_rc=0`／`census_totals_diff_rc=0`（差件 0 字节）** | 跑前＝跑后＝逐字 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（票面基线四枚字段全等；⚠ 我那行 `census_before_totals_matches_ticket_baseline=NO` 是**我自己的尺艺**——件里带前缀 `portable-tests.sh: census totals: `，字面 `=` 当然⛔ 等，字段级比⛔ 等＝**对**）；两发各带一行 `ticket 310 AC#1: GOTMPDIR=`（`gotmpdir_line=1`）⇒ **这一步在 census 发里也跑过** | `logs/30-census-{before,after}.txt`／`-totals.txt`／`30-census-totals-diff.txt` |
| ④ winsec 档 | `bash scripts/winsec-tests.sh ./internal/winsec/` | **`winsec_rc=0`** | `portable-tests.sh: four numbers (all from -v output): === RUN=101  --- PASS=58  --- FAIL=0  --- SKIP=0`；同行 `winsec_gotmpdir_line=ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/winsec-20261012-125118-1774`（那枚根**在盘上**、我没删）；⚠ 拼法照正典：**⛔** `--scope=winsec` | `logs/30-winsec.txt`／`.err` |

★**拼写雷我这一程又量了一遍**（⛔ 当门禁读数，只当"那把 GUARD 有牙"的证据）：`bash scripts/winsec-tests.sh --scope=winsec` ⇒ 逐字 `wrongspelling_rc=2`＋stderr 首行 `winsec-tests.sh: GUARD 1 - the package list for this step does not name ./internal/winsec/; it is [--scope=winsec].`（件＝`logs/30-winsec-wrongspelling.err`）⇒ 与编排者 10-12 那发**逐字同 rc、同第一行**。⚠ 附带一条我从盘上取的副作用：那一发**在被 GUARD 打回之前已经建好了根**（`winsec-20261012-125140-1926`，`12:51:41`）＝**步子跑在守卫之前**，那枚空根留在 `%LOCALAPPDATA%` 里，我没删（记在下面 §5 的枚数里）。

## 2. 格式名册：两把并排（工作树 ⇄ 导出树）＋代用尺那一支算⛔ 算掩盖

| 把 | 我的命令（逐字射程） | 枚数 | 逐字名册 |
|---|---|---|---|
| ① 工作树 | `gofmt -l cmd internal tools scripts` | **4** | `cmd\wisp\models.go`、`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`、`internal\risk\provenance.go`（`logs/30-fmt-worktree.txt`，rc=0、stderr 0 字节） |
| ② 导出树（＝`git archive HEAD go.mod go.sum internal cmd tools scripts \| tar -x`，六件套） | 树内 `gofmt -l cmd internal tools scripts` | **0** | 空；★先证取出来的东西非空＝**`go_files_in_tree=682`**（⛔ "取空了所以 0"，`logs/30-fmt-headtree.txt`） |
| ③ `gofumpt` 到底有没有 | `command -v gofumpt` | — | **`GOFUMPT=NOT_FOUND`** ⇒ 本机确实⛔ 装（腿那句复认 ✓） |

- **"⛔ 装 gofumpt、代之以 `gofmt -l`"算⛔ 算掩盖**：**⛔ 算掩盖，但它买到的东西比 CI 那把尺少。**具名、可复核（本腿同两条腿复跑、同数），且⛔ 用它把任何一格说成绿；⚠ 方向要说清：`gofmt` 比 `gofumpt` **松**，所以"工作树 `gofmt` 都报 4 枚"这一半⛔ 会被 gofumpt 洗白、而"导出树 0 枚"这一半**⛔ 能证明 gofumpt 也 0 枚** ⇒ 真·gofumpt 颜色仍欠（`Q-310-v1-3`，与腿的 `Q-310-r1-4` CI 半边同一笔）。
- ★**"工作树 4 枚 vs HEAD 0 枚"那处差异的归属：腿说"差在检出面形状（行尾族）"——机制方向对、归因写反了，本腿现量**（尺＝`tr -cd '\r' \| wc -c` 三种来路并排，件＝`logs/30-blob-eol.txt`）：
  ```
  file=cmd/wisp/models.go                CR_bytes worktree=334  head_blob=0  git_archive_tree=0
  file=internal/agent/approval/pending_read.go CR_bytes worktree=131 head_blob=0 git_archive_tree=0
  file=internal/agent/tools.go           CR_bytes worktree=225  head_blob=0  git_archive_tree=0
  file=internal/risk/provenance.go       CR_bytes worktree=1115 head_blob=0  git_archive_tree=0
  ```
  ＋`.gitattributes` 逐字有 `*.go text eol=lf`、`git config core.autocrlf=true`。⇒ **HEAD 的 blob 是 LF、导出树是 LF、只有本机工作树带 CRLF** ⇒ 那 4 枚**⛔ 是"blob ⇄ 检出的行尾差异"（那两者行尾相同），而是"本机检出的字节违反自家 `.gitattributes` 的 `eol=lf`"**；`git diff --name-only HEAD -- '*.go'`＝**0 枚**（件＝本段，与腿同数）是因为 git 做了归一化、内容比得相等。后果写明白：**在本机，任何吃 LF 的尺（CI 的 gofumpt、`gofmt -l`）都会一直报这 4 枚**，与票 310 无关、⛔ 本票射程（腿登记成 `Q-310-r1-4` 是对的，**但那一行的文字要更正**，归编排者落账；本腿⛔ 动别人的件一字）。

## 3. ★7 枚 0 字节件：ⓐ"合法的空"⇄ⓑ"那格没交"，**分开写**

尺＝`find -size 0` 枚数＝**7**（`logs/10-r1-zero-byte-names.txt` 逐字：`00-gate-310r1.err 10b-new.txt 10b2-new.txt 10b2-r1-new.txt 10b2-root.txt 10c-run.run.txt 12-roster-new.txt`）；"被交付件点名了吗"尺＝`grep -c -F <名> .scratch/wisp/probes/310/r1/*.md`。

| 件 | 命中 | 被点名 | 判 |
|---|---|---|---|
| `00-gate-310r1.err` | stderr | 3 次；`00-anchor.md:76` 同行逐字 `wc -c logs/00-gate-310r1.err → 0 （empty above＝0 字节＝起手闸门无 stderr）` | **ⓐ 合法的空**（自证行就在同一件里，最强那一形） |
| `10b2-r1-new.txt` | b3 的顶层差集 | 1 次；`20-paired-readings.md:33` 同行 `**0 枚**（logs/10b2-r1-new.txt 为空，empty above＝0 枚）` | **ⓐ 合法的空** |
| `12-roster-new.txt` | 入口档跑下来的顶层差集 | 1 次；`30-gates.md:43` | **ⓐ 合法的空**（我独立复现＝58→58、差集空），⚠ **那句写成了"差集⛔ 空＝0 新增"＝字面读是"差集⛔**不**空"**，与 0 字节件自相矛盾 ⇒ 措辞缺陷、⊘ 交付缺陷，归编排者改（本腿⛔ 动腿件） |
| `10b-new.txt` | b1 的顶层差集 | 1 次；b1 行写 `**0 枚**`＋`10b-new.txt（空）` | **ⓐ 合法的空**（自证是"（空）"那半形，比 `empty above＝0 枚` 弱一档，仍算同行有计数） |
| `10b2-new.txt` | 作废那发（b2）的顶层差集 | **0 次**（只有 `logs/10b2-*` 那种"族名"覆盖到它） | **ⓐ、但按名没交**：b2 整发作废，这枚空件⛔ 支撑任何主张；缺的是"逐枚点名＋同行自证"那一步 |
| `10b2-root.txt` | 作废那发（b2）的"根里条目"名册 | **0 次** | **ⓐ、但按名没交**：它空＝那枚根当时⛔ 有残渣＝**b2 作废的那半证据**；腿在 b2 行里写了 `root_entries=0` 的数字、⛔ 点名这枚件 |
| `10c-run.run.txt` | a4 台件那发的 stdout | **0 次**；同批 `10c-run.run.txt.err`＝134 字节 | **ⓑ 那一格没交**（唯一一枚我判"没交"的：⛔ 一条主张指向它、也⛔ 一行自证说"它应当为空"；测试二进制被杀时 stdout 为空是**可能**的，但那⛔ 是被证过的） |

⇒ **两类混⛔ 成"都是空壳"**：ⓐ **6 枚**（其中 2 枚按名没点名、1 枚措辞自相矛盾），ⓑ **1 枚**（`10c-run.run.txt`）。

## 4. ★逐笔越界尺（⛔ 区间尺，十笔各跑一遍）

尺＝`git show --name-only --format= <每笔一枚 hash>`；允许清单＝`.scratch/wisp/probes/310/r1/**` ＋ `scripts/portable-tests.sh`／`winsec-tests.sh`／`wisp-cli-tests.sh`。件＝`logs/30-percommit-scope.txt`＋逐笔 `logs/30-pc-<hash>.txt`。

```
commit=a21dcd5b paths=5  foreign_outside_allowlist=0
commit=6a8f537a paths=1  foreign_outside_allowlist=0
commit=2c0a0666 paths=1  foreign_outside_allowlist=0
commit=14af45be paths=1  foreign_outside_allowlist=0
commit=92333974 paths=1  foreign_outside_allowlist=0
commit=e0218be4 paths=1  foreign_outside_allowlist=0
commit=369a9e17 paths=1  foreign_outside_allowlist=0
commit=381bbc44 paths=1  foreign_outside_allowlist=0
commit=c5131048 paths=62 foreign_outside_allowlist=0
commit=3122306e paths=1  foreign_outside_allowlist=1
```
合计 74 枚路径、**按派单 §3 那张字面允许清单＝越界 1 枚**；那一枚逐字＝`.scratch/wisp/issues/310-test-binary-deaths-leave-temp-test-dirs-accumulating.md`（`3122306e`＝票尾那条 append-only Progress bullet），而**同一张派单 §3 的交付件清单要求交件人写那条 bullet** ⇒ **按交付件清单＝越界 0 枚**。⛔ 两种读法混着说：**清单本身漏了票 310 那一行**，本腿把这一处**具名报回**（`Q-310-v1-1`）；下一程照我做法：允许清单＝`probes/<腿>/**` ＋ 三枚脚本 ＋ **本票票面一枚**。
★**别家那批**：`.gitignore`／`go.mod`／`design/**` 未提交改动⛔ 进过这十笔（我的 `grep -v -E` 尺逐笔命中 0）；我这一程起手暂存区＝0 枚、全程⛔ add 过任何一枚别家路径。

## 5. ★第三种死法的归口（`Q-310-orch-N`）——本格只回答那一句

我的现量（同一把尺 `grep -c 'TempDir RemoveAll cleanup'`，同行落盘）：**默认 `%LOCALAPPDATA%` 兄弟根那一发＝0 行**（件＝`logs/30-ac3-run.txt` 的 `winsec_tempdir_removeall=0` 行，＋我件里那句 `empty above = 0 lines = self-cleanup never refused on this root shape`）；**换卷根那一发＝4 行**（⛔ 我重跑——正典 §1 明令⛔ 换卷跑 winsec 档，我引的是编排者 12:3x 那两发的件 `probes/orch/20261012/logs/orch-winsec-A-dtmp.txt`，引前已现读那一行）。同批两枚根的条目数我也复量：`winsec-20261012-122119-659` 里逐字 4 枚（`TestGateJudgesThePrivateSetByResolvedSID2872432645`／`…a_name_spelling_cannot_2948749545`／`…runner_descriptor_is_ju3713951561`／`TestSealNarrowsAndNamesThePrincipalItRemovedBySID3143495113`，件＝`logs/20-root-contents.txt`）⇄ 默认根那发 **0 枚**。

⇒ **要那一句的答案**：**票面现量那句"测试正常跑完 ⇒ 0 残留"要加射程，加的那一格落＝票 311 的 `AC#3`**。射程句给成（措辞归编排者，本腿只给判据形状）：**"测试正常跑完 ⇒ 0 残留"只在"专用根与 `%TEMP%` 同卷同 ACL 族"那一枚根形上成立；把根换到另一枚卷时，`t.TempDir()` 的自我收尾**自己**会被 ACL 拒（`TempDir RemoveAll cleanup: … Access is denied.`），正常跑完照样漏 ⇒ 残渣⛔ 只来自"进程被杀"这一种死法。** 落点⛔ 票 310 的 `AC#2`（那是成对两发那格，塞进去就把"名字/枚数"那一形搅浑），⛔ 顺手实现清理器（那是产码活、且是乙形，票面已裁"⛔ 乙单独采用"）；票 310 现量节那句**原句⛔ 改**（规矩：⛔ 动票面上任何人写的原句），更正与射程由编排者落台账、票 311 `AC#3` 具名处置。

## 6. 我这一程造了什么（数了、⛔ 动）

`%TEMP%` 顶层：起手 **58** → 我 (a) 那发之后 **59**（多出来的那枚逐字＝`TestRoundTripLoadMarshalLoad4222770164`，`created=2026-10-12 12:48:51.544`，**归我、可证＝本仓用例名**）→ 四把门禁跑完**仍是 59**（`roster_end_rc=0 count=count=59`）⇒ **门禁那四发顶层 0 新增**（与腿的 `12-roster-*` 同形，我独立复现）。
`%LOCALAPPDATA%\wisp-test-gotmp`：起手 **6** 枚子目录 → 现量 **11** 枚（件＝`logs/30-roots-c.txt` 逐字带时间）＝我这一程建的 **5** 枚：`v1-20261012-124853-148`（我 (b) 那发的根，里面 1 枚残渣）、`portable-20261012-125038-1212`／`portable-20261012-125141-1936`（两次 census）、`winsec-20261012-125118-1774`（门禁④）、`winsec-20261012-125140-1926`（拼写雷那一发在被 GUARD 打回之前建的）。
`D:\tmp\wisp-test-gotmp`：子目录**一枚没多**（起手 6 → 现量 6，件＝`logs/30-roots-d.txt`；`ls -1` 数到 8 条是因为那目录里还躺着腿的两枚**读数件**，见 `20-ac2-paired.md` §2 b4 行）。
**在册那 58 枚（含 `TestProbeStepByStep3843095677`＝已裁"归属不明、⛔ 动、永久"）我一枚没删、没改名、没碰；两枚专用根的残渣一枚没清。**
