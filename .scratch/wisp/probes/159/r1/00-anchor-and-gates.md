# 票 159 r1（落地腿·记录级四格）— 锚点·四格读数·交件·面核

- 角色＝**实现者**（判语归非实现者；⛔ 本腿不翻框、不动 `Status:`/票名/`-done`）。
- 派单＝`.scratch/wisp/probes/orch/20261013/txt/159r1-dispatch.txt`（已存档）＋`159r1-send.txt`。
- 工单＝`.scratch/wisp/issues/159-record-level-cleanup-c1-c5-left-by-ticket-156-acceptance.md`（只读＋尾 1 bullet）。
- **锚（逐次现读，漂移逐报）**：起手首读 `a26cb9ab` → `59e2b3ed` → `966c2aa9` → `f706dedd`（`401191b7` 后现读 14 笔：派单笔 `df120e50`＋318 收口族＋136s 族＋178/20 族；本腿 commit 前的最终锚见回执）。
- **写面**＝`docs/evidence/s1/156-exited-asks-os-r1.md`（**追加**）＋`.scratch/wisp/probes/159/r1/**`＋票 159 尾 1 bullet。
  **其余零触碰**（含 `docs/reports/pending-and-issues.md`、`cmd/wisp/**`、`scripts/**`、`frontend/**`、`design/**`、`internal/panel/**`、`probes/154`）。

## 四格读数（逐格，原文在各 logs）

### AC#1（C1）"分开 commit"在 AC#4① 上未兑现——账＋判据
- 验收程点的名 `0e95353` 现量复核＝**是它**（`git merge-base --is-ancestor 0e95353 HEAD` rc=0）：
  `git show --name-only 0e95353`（25 文件）一枚里同时含——
  `cmd/wisp/slo_windows.go`（M；diff 有 `+ return s.queryProcess()`、`+ func (s *sloSubject) queryProcess()`＝AC#2 的码）、
  `cmd/wisp/slo_exit_os_156_windows_test.go`（A＝AC#2/AC#3 的码）、
  `cmd/wisp/slo_report_144_windows_test.go`（`+// Case 13 - AC#3's ruling is ⓐ: …`＝动作①；与现 `:835` 逐字节同，md5 `a756e044366a5c957bef470027a86bc0`）。
- `-S "Case 13 - AC#3's ruling"` 三步现量：`b417d310` 加 → `10e3585` 删 → `0e95353` 补；标题＝`wip(156 checkpoint·编排者代提停在轮次上限的程): …`（**代提**）。
- ⇒ 票面 `:46`"分开 commit"对 ① 未兑现。**账＋判据原文＝`ledger-entry.txt`①**（⛔ 台账本腿零触碰，编排者代贴）。
- 判据（票面逐字）：凡票面写"分开 commit"的格，验收时要能用 `git show --name-only` 指出该动作所在的那一枚 commit 里**不含**其它格的码改。
- ⚠ 不拆已落地的码（改写历史）；本格产物＝账＋判据，不是新 diff。
- 原文＝`logs/159r1-ac1-commit.txt`。

### AC#3（C3）数字归属＋悬空指针
- **自己复算，与票面数吻合**（无分歧可报）：r1 那发 `gate-post`＝`^=== RUN` **143**／`^---` **83**（FAIL=0）；名册两向 comm only-in-post＝**4 枚**（15/16/17/18）、only-in-pre＝0；
  `144/84＋5 枚`＝**`gate-r3-pre`**（RUN=144、`^---`=84、verdicts 里 `TestSLO156*`=5；r4 的 `gate-r4-post` 同形）。r1 上文"改前"那对 `139/79` 现量正确。
  case 19 在 `gate-post` 日志 **0 命中**（当时未落；`0e95353` 才引入）⇒ r1 当轮不存在"第 5 枚"。
- "见第 7 格"指针＝**在 r1 写作时不存在**（r1 停在轮次上限，计划中的第 7 格从未落笔；今天那个"第 7 格"＝r2 伪授权两栏 `:627`，不含该读数）。
  原文两句（`:283`/`:285`）**一字未删**；真实落点（现量）：门禁四数＝`.scratch/wisp/probes/156/gate-post/cmdwisp.log`；
  用例逐枚＝`cmd/wisp/slo_exit_os_156_windows_test.go:100/:141/:183/:232/:270`＋`roster-post-cmdwisp.txt:57-60`。
- 追加节已落 `docs/evidence/s1/156-exited-asks-os-r1.md` 末尾（标题带"编排者／收口程追加，非 r1 原文"）。
- **同形登记（只报不改）**：r1 正文另有同类前向引用——"第 4 格"×8（`:159/:166/:236/:247/:268/:273/:290/:298`）、"§9／第 9 格"×6（`:107/:177/:212/:234/:259/:303`）——
  都指 r1 从未落笔的计划节（变异表／未测章），今日相关内容在 r2 第 5 格与各日志里；⛔ 本腿未动、报回编排者裁。
- 原文＝`logs/159r1-ac3-recount.txt`。

### AC#4（C4）可执行判据＋`wisp slo` 老账（同一条记录）
- (a) `go list -deps ./cmd/wisp/ | grep -c 'CarlosShao/wisp/frontend'` ⇒ **1**（承载 `//go:embed all:dist` 的包在依赖图上；`-f …EmbedPatterns` 同发打出 `[all:dist]`；`-test` 变体同两枚）。
- (b) `git diff-tree --no-commit-id --name-only -r <commit> -- frontend/dist` ⇒ **输出非空**＝动了嵌入源目录（空结果 rc 仍 0，⛔ 别拿 rc 当尺）。
- 正控 `9318bb81` ⇒ 命中 `frontend/dist/.gitkeep`（**响**）；负控 `a2a3b3c`（只动 `frontend/src/**` 3 枚）⇒ 空。
- dist↔src（如实）：`npm run build`（`tsc -b && vite build`）产 `dist/`；`dist/*` gitignore 只留 `.gitkeep` 锚（vite 的 `keepDistAnchor` 回写锚）；
  工作树今天另有未入库 build 产物（`frontend/dist/index.html`＋`assets/`，`check-ignore` 命中）——commit 级判据看不见它们。
- 老账并排：`wisp slo` 今天不被任何 `go test` 执行（复验 3 条 grep：`"slo"` 2 命中皆注释；`cmdSLO` 2 命中皆注释；`exec.Command` 邻域 2 枚皆文件名/测试名）；
  真调用方＝`scripts/slo-check.ps1`（`wisp slo -state`）＋`.github/workflows/slo-fresh.yml`。
- 坏尺复核：`grep -rln 'frontend' cmd/wisp/ | wc -l` ⇒ **11**（票 159 当时 2；树演化）⇒"搜字样"尺今天更坏。
- 记录原文＝`ledger-entry.txt`②；读数原文＝`logs/159r1-ac4-yardstick.txt`。

### AC#5（C5）case 19 归类＋同类扫描
- 追加一句已落 r1 文件末尾：`case 19 归类＝**夹具守卫**，不是"会响的检"`。
- 同类扫描（尺＋输出＝`logs/159r1-ac5-scan.txt`）：**命中 2 枚**——
  ①票面点名：`:298` 判定句"（15/17/18/19）"vs `:276` 表行 19"**夹具自己**"；
  ②同类·计数面候选：`:579`"这 8 枚"vs `:577-578` 枚举 **9** 项（只列不改、报回；⛔ 未默默只补被点名的那一枚）。
- 扫过、按尺不命中逐条写明：`:266`／`:268`／`:285`／`:470`／`:471`／`:528`／`:541-545`＋r2 第 4/5 格。

## AC#2 押后（具名）
- `cmd/wisp/slo_report_144_windows_test.go`（`:742`／`:749`／`:743`＋`:752` 三处时态）属 `cmd/wisp` 面；
  本波该面归 **133-r1** ⇒ 本腿**不做**、**具名押后**（派单缺口）。

## 交件（6 枚，每件 ≤200 行；逐枚行数）
| 路径 | 行数 | 是什么 |
|---|---|---|
| `.scratch/wisp/probes/159/r1/00-anchor-and-gates.md` | 本件 | 锚＋四格读数＋交件面 |
| `.scratch/wisp/probes/159/r1/ledger-entry.txt` | 37 | AC#1②／AC#4 的台账原文草案（编排者代贴） |
| `.scratch/wisp/probes/159/r1/logs/159r1-ac1-commit.txt` | 116 | AC#1 原文（锚/漂移＋0e95353 全量） |
| `.scratch/wisp/probes/159/r1/logs/159r1-ac3-recount.txt` | 50 | AC#3 原文（四发门禁＋名册＋case19 时间线） |
| `.scratch/wisp/probes/159/r1/logs/159r1-ac4-yardstick.txt` | 73 | AC#4 原文（判据/正负控/老账/坏尺） |
| `.scratch/wisp/probes/159/r1/logs/159r1-ac5-scan.txt` | 96 | AC#5 原文（尺＋逐枚对＋结论） |
（票 159 尾 bullet＝对票面的追加，不算件。）

## 临时件（只建不删；未入库；按 `issues/README` 规则 8 留盘）
| 路径 | 体积 | 被证据件引用？ |
|---|---|---|
| `logs/159r1-ac0-anchor.txt` | 1912 B | 内容已并入 `ac1-commit`；本文只登记 |
| `logs/159r1-ac1-case13head.txt` | 583 B | 同上（并入 `ac1-commit`） |
| `logs/159r1-ac1-diff-detail.txt` | 2589 B | 同上 |
| `logs/159r1-ac3-gatepre-and-grepanomaly.txt` | 547 B | 并入 `ac3-recount` |
| `logs/159r1-ac4-embed-and-dist.txt` | 3990 B | 手工重放前的首轮；读数与 `ac4-yardstick` 同形 |
| `logs/159r1-ac4-golist.txt` | 1669 B | 同上 |
| `logs/159r1-ac4-positive-control.txt` | 3119 B | 同上 |
| `logs/159r1-ac4-slo-account.txt` | 8291 B | 同上 |
| `logs/159r1-ac4-dist-src-relation.txt` | 1485 B | 同上 |

## 门禁（面核）与"没做"
- commit＝**只 commit 不 push**；`add`/`commit` 两处**显式 pathspec**；⛔ 未用 `add -A`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`。面核＝`git show --name-only <本腿 commit>`（回执逐枚给）。
- 没做：AC#2（具名押后）；没跑 `go test` 全套（本腿无此面）；没动任何冻结面；没动台账／HANDOVER／injection-timeline。
