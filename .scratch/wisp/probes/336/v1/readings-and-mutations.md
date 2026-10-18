# 票 336 终裁腿 336-v1·读数与必答（具名分件，正件＝`336-v1-verdict.md`）

射程说明：下表每一枚数都出自**我自己**建的仓外导出台树（`D:\tmp\wisp\20261018\336v1\tree-*`），⛔ 抄腿、⛔ 抄编排者、⛔ 抄派单正文。四数尺＝`RUN=^=== RUN`（含子测试）⇄`PASS_TOP=^--- PASS`（只顶层）⇄`FAIL_TOP=^--- FAIL`⇄`SKIP=^--- SKIP`，⛔ 混用⛔ 相减。原始日志逐枚＝`.scratch/wisp/probes/336/v1/logs/**`。

## 一、起手锚与三枚关键 blob（量的，`logs/00-*.txt`／`logs/02-roster-raw.txt`）

| 事 | 我的读数 | 尺 |
|---|---|---|
| 起手 HEAD | `bcf90e06`（12:38:5x）；建台时 HEAD 已漂到 `0fd15c6f`，末笔读数时又漂到 `2038c45d` | `git log -1 --format="%H %ct %s"` 现跑（三次，见各 runner 首行） |
| 台面锚 | `e5021339`＝`336-r1` 严格前缀尺末笔（18 笔全序在 `logs/00-commits-336r1.txt`） | `git log --all --grep='^票 336 落地腿 336-r1'`（先按定式采样真实主题形状） |
| `blacklist.go` | 起手锚⇄`e5021339`⇄`ed8ff551`⇄工作树＝**同一串 `c1799cd24afce19e6accef23db23a7d7ee746dec`**；18 笔 distinct＝1 | `git rev-parse <rev>:path` 逐枚＋`git hash-object` 工作树 |
| 牙本体 | `anchor_fallback_teeth_336_windows_test.go`＝`a429cdae29723adfa99de95f0f98cc52ac0fdf08`（160 行）／`…_other_test.go`＝`f5f888de…`（92 行） | `git ls-tree HEAD --` ⇄ 我导出树的 `hash-object` 逐字同串 |
| 写面 | 80 枚唯一路径 ⊆ 三处；越界 0；`.go` 只 2 枚且都是 `_test.go` | 逐笔 `git show --name-only` 去重＋白名单反尺（正控＝换射程命中 2） |
| AC 框 | 未勾 **7**／已勾 **0**（判据区行段 20-26） | `sed -n '20,26p' … \| grep -c '^[[:space:]]*- \['`（⛔ 全档计数，与 `'- \[ \] \*\*AC#'` 形对撞＝同为 7） |

## 二、AC#0 四发我自己重放的逐字结果（`logs/20-ac0-*.log`／`logs/20-ac0-runner.log`）

| 发 | 内容锚落点（diff 首行，只当读数） | 四数＋rc | 红名册 |
|---|---|---|---|
| 基线（未突变） | — | `10/3/0/0 rc=0` | 无 |
| M-V2a | `237c237`（`local…aUnder(...)` → `anchorMatches(...)`） | `10/3/0/0 rc=0` | **无＝全绿** |
| M-V2b | `227c227` | `10/3/0/0 rc=0` | **无＝全绿** |
| M-V2c | `231c231` | `10/2/1/0 rc=1` | `TestAnchorMissReasonNamesCandidateAndMissingAnchor/installing-aws-moves-the-named-anchor-to-kube` |
| M-V3 | `366c366`（`return isUnder(..)\|\|isUnder(..)` → `return true`） | `10/1/2/0 rc=1` | `TestAnchorMissReasonNamesCandidateAndMissingAnchor`＋`TestFileBaseNameClaimChangesNoClassOrReason`；★闸门格 `TestAnchorMissKeepsNonOverridableClassAndNeedL2` 逐字 `--- PASS` |

⇒ 与票面"现量"表**逐字一致**，无更正。八次哈希 `match=YES`；每发 `落点命中枚数 1→0`、`changed_lines=2`。

## 三、AC#2 八发去牙扫描（我重跑，`logs/30-ac2-*`；⛔ 复用腿的树）

成对正控（未突变）＝`rc=0 RUN=17 PASS_TOP=1 FAIL_TOP=0 SKIP=0`，八枚逐枚 `matched=false clause=true`（`logs/30-ac2-00-pristine-teeth.log`）。

| 发 | 操作数 | 落点（读数） | 牙四数 | 牙红名册（子用例逐名） | 在册三枚仪器 |
|---|---|---|---|---|---|
| MU-1a | `aEq` | `:220` | `rc=1 17/0/1/0` | git-credentials | `rc=0 10/3/0/0` |
| MU-1b | `aUnder`（AND 左半） | `:224` | `rc=1 17/0/1/0` | wisp-config-under（class 仍 A，红的是 `:138` 那句"reason 必须点名本锚"） | `rc=0` |
| MU-1c | `aEq`（AND 右半） | `:225` | `rc=1 17/0/1/0` | wisp-config-under ＋ wisp-config-eq（class=B） | `rc=0` |
| MU-1d | `aUnder` | `:227` | `rc=1 17/0/1/0` | ssh（class=B） | `rc=0` |
| MU-1e | `aEq` | `:229` | `rc=1 17/0/1/0` | aws（class=B） | **`rc=1 10/2/1/0`**＝ⓐ |
| MU-1f | `aEq` | `:231` | `rc=1 17/0/1/0` | kube（class=B） | **`rc=1 10/2/1/0`**＝ⓑ |
| MU-1g | `aUnder` | `:235` | `rc=1 17/0/1/0` | dpapi-protect（class=B） | `rc=0` |
| MU-1h | `aUnder` | `:237` | `rc=1 17/0/1/0` | local-credentials（class=B） | `rc=0` |

八枚落点＝`220 224 225 227 229 231 235 237`＝8 枚不同物理行（尺＝我的 runner 自打的 `diff` 头，`logs/09-mu1b-verbatim.log`）；每发 `branch_removals_on_disk=1`（＝只摘一枚，独立于"落点命中"另一把尺）。红名册**并集**＝8 枚唯一子用例名（`logs/61-redname-union-two-rulers.log`；同一把尺打在腿的 8 份跟踪件上＝同样 8 枚 ⇒ 腿票尾那句"只有 7 枚唯一名"数的是**逐发首红**，两把尺都成立、单位不同）。

★另加一发（腿没发，`logs/acY-*.log`）：**M-V3（变严）打在牙上**＝`rc=1 17/0/1/0`、红＝8 枚主行＋8 枚 `offchain` 全红（`FAIL_TOP` 只算顶层＝1，⛔ 与子用例数相减）⇒ `334-v1` 那句"闸门格恒不红"那一支的残余被牙的链外臂接住了。

## 四、AC#3 乙形那一发（派单点名的必做发，`logs/40-ac3-runner.log`）

- 探针来源 blob＝`69a70fa454886029aea9f8f716cb372908467339`；我放进**仓外树**后的落盘 `hash-object`＝**同串**；`//go:build` 枚数＝0；行数＝248。
- 断言普查（我自己的尺）：`t_Log=13／t_Fatalf=7／t_Errorf=0／t_Skip=0` ⇒ **零断言件**（会红的只有 setup 的 `Fatalf`＝台面故障，不承载判据）。
- 发 1（未突变）：`rc=0 RUN=22 PASS_TOP=4 FAIL_TOP=0 SKIP=0`，wall 3.89s。
- 发 2（只摘 `:237` 一枚，突变哈希 `5f8789548f624e12ba2527c016f0f9b040c36e3a`）：**`rc` 仍＝0**、四数**一字不变** `22/4/0/0`，wall 2.55s；降级只在可读行里显形（shape 作差 2 行：`chain:case9-local-creds class=A→B`）。
- 同一棵树同一发打在甲形（牙）上：`rc=1 17/0/1/0`＋红名册逐名 `TestAnchorFallbackTeethPerAnchor336/local-credentials`（wall 1.57s；牙 pristine 对照 1.24s）。
- 线种尺（乙形未突变输出按行种数）：`chain:9／chain-clean:9／offchain:9／real-temp:4／real-appdata:1／OVERRIDE-chain:9／OVERRIDE-real-appdata:1／PREFIX-anchor_root:7／LSTAT:4` ⇒ chain 族＝27（＝腿那句"27 枚具名 want"），全量读数行＝46。
- ⇒ **不敏感这一形复现**：乙形对"摘掉一枚回落支"零反应；甲形同一发当场红并点名。

## 五、AC#4 成对台面＋整包＋门禁（`logs/50-*`／`logs/51-*`／`logs/gates-runner.log`）

- 定向名册（四枚仪器并列）：dshape 改前 `10/3/0/0 rc=0` ⇄ 改后 `27/4/0/0 rc=0`；ci-shape 改前 `10/3/0/0 rc=0` ⇄ 改后 `27/4/0/0 rc=0` ⇒ **两形同色**。
- 交件锚三发（我自己重打，⛔ 抄编排者）：pristine `27/4/0/0 rc=0`（`teeth` 落盘 blob＝跟踪件同串）→ 摘 `:237` 突变 `5f878954…` `rc=1 27/3/1/0` 红名册逐名 `TestAnchorFallbackTeethPerAnchor336/local-credentials` → 还原 `c1799cd2…` `match=YES` 复跑 `27/4/0/0 rc=0`。
- 整包 `-run .`：改前 `230/135/5/1 rc=1` ⇄ 改后 `247/136/5/1 rc=1`；红名册逐名作差 **`diff_rc=0`／0 行**（五枚同名，`logs/50-ac4-full-roster.diff` 已标注为空＝零红差＋正控）；SKIP 两发同一枚 `TestSyncRegistryProbeLive`；`RUN` 差 17 枚逐名对得上牙名册（8 主＋8 offchain＋1 顶层）。
- 台面敏感：同一棵树（after）整包 ci-shape＝`247/140/1/1 rc=1`（红名册只 1 枚）⇄ d-shape＝`247/136/5/1 rc=1`（5 枚）；差集＝那四枚 `TestSyncAncestor…`／`TestSyncFirstActable…116`。**改前那棵树（根本没有 `336` 字节）的 d-shape 同一枚名册同样红这 4 枚** ⇒ 判＝既有台面形状、⛔ 本票新增红（归位＝台面／`A1024` 第③节那族）。`TestC26…` 两棵树一发红句逐字同、四行全指 `internal\panel\composer.go:127`＋`display_scrub_216.go:26/37/218`＝别腿在飞面。
- 门禁四数（只在共享工作树，跑前 `unset GOTMPDIR TMPDIR TMP TEMP`，逐字 `logs/gates-runner.log`）：`gofumpt.exe -l`（新增两枚件）rc=**0**／列出 **0** 枚（正控＝包级列出 1 枚 `internal\risk\provenance.go`＝CRLF 存量假脏）；`sh scripts/d22scan.sh` rc=**0** clean（ban#8 `internal/=590`、`cmd/=143`）；`sh scripts/check-path-length-budget.sh` rc=**0** VERDICT **GREEN**（分母 **15615**、over-budget 60＝rostered 60、not-in-roster 0）；`go vet ./internal/risk/` rc=**0**。
- 我的写面最长路径＝81 字符（分母尺现跑）。

## 六、必答（⛔ 我盖章不算；每一答都指回上面那把尺）

**(i) 票 334 `AC#2` 那半格现在成不成立＝成立（我裁：翻框交给编排者）。** 判据原文那句"任何试图把这一形降成 B 档／可豁免的改动必须被咬红"——`334-v1` 当时判不成立的两枚反例（只摘 `:237`、只摘 `:227`）我这次同一把尺重放：在册三枚仪器**仍全绿**（`10/3/0/0 rc=0`，本表第二节），而牙在同一发上**当场红并点名**（`MU-1d`／`MU-1h`，本表第三节）。"射程从 2 枚扩到 8 枚"这句话**站得住**：我逐枚验＝`leg_red=0` 六枚、`leg_red=2` 两枚（`MU-1e`／`MU-1f`＝旧仪器的 `.aws`＋`.kube`），牙 `rc=1` 八枚八枚 ⇒ 2 枚扩到 8 枚＝**量出来的**。另一半（M-V3 恒不红）也被牙的链外臂接住（`acY` 一发 16 枚子用例红）。⇒ 334 那半格的两支残余都闭合，具备翻框条件；⛔ 我翻、⛔ 我改名、⛔ 我加 `-done`。

**(ii) `:224`⇄`:225` 要不要第三枚世界才算逐枚独立＝不要"第三枚世界"才算达成，但我实测出一枚**便宜**的加件能让它逐枚互不重叠。** 读数：两枚**已经可判别**（红集合 {under} ⇄ {under, eq}）＋落点已是两枚不同物理行 ⇒ 按票面"每枚锚各发一两枚用例，落点必须不同"我判**达成**。代价实测（我自己那把探针尺 `logs/60-m2-*`）：把候选从 `<APPDATA>\Roaming` 挪到更深一层（`<APPDATA>\Roaming\subfolder\nest`），未突变 `class=A`（`rule=A:%APPDATA%-wisp-config.toml`，两枚操作数都 `matched=false clause=true`）；只摘 `:224` ⇒ 该行**落到 `class=B`／`rule=B-tier`**（同发只摘 `:225` 也落到 B）。⇒ 存在一个让 `:224` **单独**把 class 挪走的世界；代价＝**一枚新行**进 `teeth336Roster()`（现有件里加 1 行＋最多 1 枚常量，⛔ 新件；`relCand` 已支持任意 profile 相对目录）＝**0 枚文件／1–3 行**。★注意这一枚是**扩面**，本腿⛔ 自行加（⛔ 动 `internal/risk/**`）＝交编排者拍：要么按现状成立＋具名残余，要么派一枚只加这一行的补牙腿。红名册互不重叠不是票面写着的判据，所以我没把它当缺口（见 `logs/NN-blocked-2.log`）。

**(iii) POSIX 半格那一形＝成立（诚实的缺世界＋会响的绊线）；票 339 那笔债⛔ 算本票新增。** 我自己重测（⛔ 抄腿）：静态尺＝`pathresolver_other.go:10` 逐字 `func resolveHandle(p string) (string, bool) { return "", false }` ⇄ `pathresolver_windows.go:25` 有真身（`logs/23b-posix-mechanism-sanitized.log`）；`uncertainAnchorMiss` 首行 `if cand.root == "" || anchor.root == "" { return false }` ⇒ 这一平台**造不出**"摘这一枚⇒必红"的世界。★并在**Linux 现跑**（docker `golang:1.27`／`go1.27.1 linux/amd64`，我自己的导出树，`logs/24b-posix-linux-run.log`）：绊线件 `--- PASS`（`RUN=11 PASS_TOP=2 FAIL_TOP=2 SKIP=0 rc=1`），四枚 placement 逐枚 `cand_root=EMPTY verdict=B`；两枚红＝票 334 字节 `TestAnchorMiss*` 那两枚（真身行 `anchor_miss_reason_334_test.go:131/:156/:193/:197/:201` 现读逐字是 `fixture guard…want A` 与 `class: got %s, want A` 那几句）；`TestAnchorFallbackTeethPerAnchor336` 命中 **0** 枚（`//go:build windows` 生效＝不是"偷偷少跑"）。339 那笔债归位：那枚文件的引入笔＝`4711955f`（`334-r1`，2026-10-17），腿 18 笔对它命中 **0** ⇒ **票 339 独立成立、⛔ 记在 336 头上**（编排者那句复认我也自己重跑了一次）。台面故障具名：我第一次 docker 发 `-e GOCACHE=/tmp/gocache` 被 MSYS 改写成 `C:/Users/...`，Go 直接死在 "GOMODCACHE entry is relative"（`rc=1`、零枚 RUN）＝**⛔ 计入判据**，改正版把 cache 目录设在容器内。

**(iv) 乙形搬不搬＝不搬；`AC#3` 判成立。** 决定性那一发我复跑＝**摘牙后乙形 `rc` 仍＝0、四数一字不变**（本表第四节）⇒ 它现在对这件事**不敏感**（零断言、13 发 `t.Log` 只在人读的行里显形）。搬＝248 行一件＋逐枚补具名 want（chain 族 27 枚，全量读数 46 行）＋`//go:build` 0 枚 ⇒ 进跟踪面就吃 CI 两平台射程，而它要钉的东西甲形已经用 160 行＋92 行、1.2–1.6 秒、带 tag 钉住了。**甲够。** 乙的材料继续留在只读面（`logs/40-probe-334v1-shapeB.txt`），将来若有人要"把三值对撞当断言"再搬，搬的时候必须先补 want、⛔ 原样入库当零断言件。

**(v) 腿那六处自抓作废里，该退件的＝0 枚；一处升级为"必须具名"但已由腿自己改正。** 逐处我的判语＋尺：① `f67048d9` 正文"闸门格在红名册里"＝写错 ⇒ 腿已具名改口，我也独立重放到 `--- PASS` 侧（本表第二节）＝**已闭**；② `teeth-00` 里 `wisp-config-under` 从 `<profile>\work` 出发是 CERTAIN miss ⇒ 落点改正为 `<APPDATA>\Roaming`＝我现跑该世界逐枚 `matched=false clause=true`（`logs/30-ac2-00-pristine-teeth.log`）＝**改正是对的、⛔ 退件**；③ `MU-1b` 只咬 `:138`、`class` 没降＝**不是作废、是读数**（正是"只看 class 是盲的"），牙咬住了；④ `ac4-10／ac4-11` 两发假读数（perl 表达式写进 `bash -c` 双引号⇒反引号被命令替换吃掉）⇒ 腿已具名作废并另落改正版；★我自己这一族**又犯了一次**（见下"我侧自抓"）＝同族第二日复发，具名给编排者：任何"表达式穿过双引号"的台件都算高危形；⑤ 6 枚 0 字节 artefact（双引号里写 glob 不展开）⇒ 按"0 字节＝没交"⛔ 入库、⛔ 删、真名册在同名非空件里＝我复认成立、⛔ 影响任何判据（我的对应改正＝把每一枚空读数文件**就地写上尺＋正控**，见本面 8 枚 `.rednames`／`.diff`／`g1`／`g4`）；⑥ report 里 `分节号` 那句"现量 0 枚"是假读数 ⇒ 腿改为描述它的那枚字符＋现量 0，我一把独立尺复验＝`report.md` 0 枚、`logs/` 0 枚文件（正控＝同形尺找 `TEETH336` 命中 26）。⇒ **没有一处该退件**：六处都是腿自抓、自改、并留下可重跑的尺，正是本仓要的形状。

**我侧自抓（⛔ 等谁替我圆场）**：① 我的 `run-ac4-surface.sh` 那句 `teeth_red_hits=8`＝**射程错**（grep 整份日志、把 `=== RUN`／`--- PASS` 里的文件名算进去），改正尺＝只算 `--- FAIL:` 行⇒0 枚（`logs/52-ac4-teeth-red-correction.log`，正控＝同尺命中 `TestC26` 1 枚）；② 同一件 ci-shape 的 `GOTMPDIR` 目录名多带一枚 `}`（引号位置），同卷同形、不影响四数，具名留档；③ 我第一次 AC#5 尺**按派单原文用了三条不存在的路径**（`internal/tools/rules_gateway.go`／`internal/risk/thresholds.go`／`internal/risk/tokens_fourway_test.go`），产出 `same=NO／distinct=19`＝坏尺⛔ 读数，已作废并改正（`logs/25-forbidden-faces-SUPERSEDED.txt` ⇄ `logs/07-forbidden-faces-corrected.log`）；④ 我一度把 `blacklist.go` 整份逐字拷进我的写面、又拷了一份带产码注释的 Linux 机制读数 ⇒ 两处都含分节号＝违我面规矩，已改为**具名切片＋OCT-302-247 转写**、原枚就地写明 superseded（⛔ 删）。

## 七、遇未定义的两枚（⛔ 整腿停摆，都按"不依赖它的格"继续）

- `logs/NN-blocked-1.log`＝派单／票面禁区清单里 `rules_gateway.go`／`thresholds.go` 的包路径与树不符（真身＝`internal/risk/`、`internal/observe/`）。我按"取真名、量同一件事"处理并把差异具名交给编排者，⛔ 自行改判据。
- `logs/NN-blocked-2.log`＝"每枚锚各发一两枚用例"要不要红名册逐枚**互不重叠**，票面与判据原文都没写。我按票面字面（"落点必须不同"＋"红名册逐名对撞"）裁＝达成＋具名残余，并把可让它互不重叠的那一枚便宜加件（必答 (ii)）交给编排者拍，⛔ 自行扩面。
