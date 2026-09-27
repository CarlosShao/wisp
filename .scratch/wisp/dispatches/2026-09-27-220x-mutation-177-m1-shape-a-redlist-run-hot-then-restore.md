# 派单 177-m1（**AC#1b：甲形变异台件——把"读码推导的红名单"变成现跑读数**）＝票 177 的 AC#1 那句"不许推理，要跑台件"的欠账；本程**量完即还原**，交付物里**一枚产码都不许留下**

- 派单时刻：09-27 22:0x｜锚点＝**你自己 step 0 现量**（别抄我的号）。
- 票面＝`.scratch/wisp/issues/177-stamping-external-content-and-r4-scans-path-arguments.md`（真名很长，用 `ls .scratch/wisp/issues | grep '^177-'` 取；**AC#1b 那一格连 21:5x 那六条追加一起抄进你的表**）。
- 上游表＝`docs/evidence/s1/177-c25-r4-path-exemption-c1.md`（**它是被告的推导，不是证据**；它的 §4.1 六行红名单整片标〔读码推导〕——**你这一程的存在理由就是去量它**）。它的读数件＝`.scratch/wisp/probes/177/c1/readings.md`。
- 已批准的形状（owner 原话"行，都按推荐就完事了"，台账 `A349`）＝**甲：精确豁免**，只豁免"宿主自己写下过的那几条具体路径"；**目录级／前缀级／配置项级一律不在批准范围内**（`Q-61` 那一行的边界，逐字）。
- ⚠ **预算硬顶 ≤60 次工具调用**；**先测→再还原→再写→再提交**；**不许勾任何框**（AC#1/AC#1b 都由我裁）。

## 0. 起手五件（缺一件就停手上报）

`date`（**每条带时间的账之前重新取一次**）｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手，一个字不写）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/ cmd/`（**必须空**；不空＝停手上报，别在脏树上做变异）｜基线 `go test -count=1 ./internal/risk/ ./internal/tools/`（**逐包，禁全仓 `./...`**）四数＋名册口径分开报（顶层 PASS vs 含子测试 RUN）。

## 1. 排他约束（本程唯一一枚在树里改产码的程）

`174-v1`（验收位）**必须先交完**我才会发这一单。**你这一段期间编队里不得再有第二枚写手或取数腿**——你要临时改 `internal/risk/**`，任何并发程在你变异窗口内跑出的红名集合都不可归因。若你起手发现树不干净或有别的写手在动 `internal/risk/**`：**停手上报，不要自己清场**。

## 2. 主菜：甲形变异（AC#1b 的第一问）

**只做形状，不做产品**：把甲形实现成"在一枚 mark 里，把宿主声明的那一段连续窗口排除在片段索引之外"，落点按上游表 §6 的最小路线（**行号自己现读，别照抄**）：

1. `internal/risk/taintmatch.go` 建索引那一圈（`newFragmentIndex`）——**只跳过被声明的窗口**；`f.src` 全文与 `strings.Contains` 复核那两半**不许动**（动了就会造出假命中，那是另一枚缺陷）。
2. `Mark` 侧：跨度定位**必须在归一化之后**（`normalizeTaint` 会丢空白与零宽字符，用原始下标一定偏）；**空判定必须排在排除之前**（否则你会亲眼看到 `TestMarkEmptyAfterNormalization` 红，那是形状选错，不是发现）。
3. **不许**给 `agent.Loop` 加收任何带 `taskID` 的导出方法（`probes/154/gate-clauses.sh` 的 **G3** 安静腿就是来抓这个的）；**不许**在 `risk.PathResolver` 之外用 `filepath.Clean|Abs`（d22scan ban 2）；**不许**新造"按 scope 维护一条 path 列表"（那是名册，上游表 §2 末尾明写不推荐，本程只是量形状、不替它选路）。
4. **给豁免喂路径的载具已定**（我这一轮裁的，别再自行改路）：**不给 C1 `Result` 加收字段**——契约面不动；台件里用一个包内的最小通道（哪怕是未导出的临时参数）把"P 是哪一段"送到 `Mark`。

要答的问题（每条都要可复制命令＋逐名红集合）：

- **Q1 上游表 §4.1 那六行，哪几枚真红？** 逐枚点名（`TestFragmentMatchExactBoundary`／`TestFragmentIndexShortSourceNeverMatches`／`TestFragmentHashCollisionCannotFakeHit` 会不会**编译红**、`TestMarkEmptyAfterNormalization`、那三枚枚数类 `provenance_test.go:857`／`fs_test.go:35`／`provenance_test.go:916` 在"不为 P 另加一枚 mark"的前提下响不响）。
- **Q2 ⚠ 过界哨（这一问最重要）：`provenance_test.go:136`／`:139` 那两枚 fail-closed 响不响？** 上游表判"甲形下它们**不该红**，它们红＝豁免越界的信号"。**若它们红：立刻停手回禀，不许改判据、不许放宽那两枚、不许把它降级成"预期 casualties"**（AGENTS §1.1＋票面 AC#4）。
- **Q3 形状自证（防"恒不变"）**：拿一枚你确认**会**红的用例做正控——把你的排除窗口临时改宽成"整段 mark 不入索引"，看 R-1/R-3 那一族（上游表 `probes/177/c1/fixture-reverse-criterion.md` 的三发形状）响不响。**一枚都不会红的判据＝装饰**，那一形就地作废并写明。
- **Q4 d22scan／gate-clauses 两向**：变异态跑一次、还原后再跑一次，把两次的 `腿数＝… 声明与实测不符＝…` 与 ban 枚数都贴出来（⚠ **别跑 `probes/161/r6/flip-declaration.sh`**，它跑一次就脏跟踪日志）。

## 3. 次菜（预算够才做）：乙形的一发对照

不要求完整实现乙形。只要造出**那一发结构性冲突**并留证据：**同一条路径字符串**，一发来自合法续读、一发来自外来内容——在"按值放行"的匹配侧下**二者逐字相同**，因此按值放行必然同时放过外来那一发。形式自由（一枚临时用例或一段可跑的台件都行），但**必须真跑过**并贴读数；跑不出来就写"未证，仍是推导"。

## 4. 还原协议（这一步之后才许提交）

变异**一律不入库**。逐文件 `git cat-file blob HEAD:<那枚路径> > <那枚路径>`（禁 `checkout`／`restore`／`revert`／`clean`），然后：
`git status --porcelain -- internal/ cmd/` **必须空并贴出来** ＋ `go test -count=1 ./internal/risk/ ./internal/tools/` **必须复绿**。
⚠ **中途任何一步失败（编译不过、panic、预算耗尽）：第一件事就是还原**，别留半棵树给别人。还原后 `git diff --numstat` 对产码应为空。

## 5. 门禁（还原之后自己重跑，别引上游的数）

`go test -count=1 ./internal/risk/`＋`./internal/tools/`｜`sh scripts/d22scan.sh`（应 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线两向 `comm` 差集）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（rc 与"腿数／不符枚数"逐字贴）。⚠ 引用"某区间几枚红"必须**连口径（哪几个包）一起引并带锚点 sha**。

## 6. 写面（超出即越权）

`docs/evidence/s1/177-mutation-shape-a-redlist-m1.md`（新建，你的表）｜`.scratch/wisp/probes/177/m1/**`（台件与读数；⚠ **`logdir` 取脚本自身目录**，不继承 CWD）｜票面 177 的 Progress log（**只追加、不勾任何框**；写前 `git status --porcelain -- 那枚票面` 必须为空）。
**禁改**：`internal/**` 与 `cmd/**`（§2 的临时变异只在 §4 协议内存在，且**绝不入提交**）、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`（**含 `:1532` 那行刚落的 DEFERRED，一字节不许动**）、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、别人的票面与证据件、`probes/**` 既有台件（只读，含 `probes/177/c1/**`）。⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）⇒ **不提交、不还原、不补完、不评论**。
**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`；heredoc 必须加引号（`<<'EOF'`）；含反引号的中文段走 Edit 工具，别用双引号串。**只 commit，绝不 push。**⚠ **更正自己写错的东西只许追加新提交——本仓禁 `--amend`**。⚠ **绝不在仓库目录内建 worktree 或 checkout**，也别建第二份 `.git`。

## 7. 交件报告（顺序固定）

1. step-0 五件＋排他约束是否满足；2. **本程没测什么**；3. **Q1 逐枚红集合**（变异态 vs 还原态两向）；4. **Q2 过界哨读数**（那两枚 fail-closed 各贴逐字结果）；5. Q3 正控（你的宽一档变异真的响了吗）；6. Q4 两把尺两向；7. 次菜那一发（或未做的原因）；8. 门禁四枚＋名册差集；9. **还原证明**（空 status ＋ 复绿 ＋ `git show --stat` 逐枚列出你名下提交**不含任何产码**）；10. 被拒／没成功的调用（取数前还是后）；11. 有没有跑过删除命令；12. 伪授权两栏（各带出处）；13. 凭据值零抄录；14. `next=`（含一句：甲形落地前还缺哪几格）。
⚠ **每个"几枚"旁边附可复制命令**。**先测→再还原→再写→再提交。** 我这段话里每一条（含"六行红名单""那两枚不该红""三枚枚数类不响"）**都是未验证断言**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
