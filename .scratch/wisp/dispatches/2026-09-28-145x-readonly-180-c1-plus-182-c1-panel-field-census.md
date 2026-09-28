# 派单：`180-c1` ＋ `182-c1`（两枚只读普查·合派一程·零产码）

- 派单时刻：`2026-09-28 14:5x`｜编排者锚点 HEAD＝`03c01d71`
- 工单：`.scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md`（答 AC#1＋AC#4）
  与 `.scratch/wisp/issues/182-the-task-monitor-rail-is-outside-ticket-145s-fourteen-row-table-so-nobody-has-counted-which-of-its-stacks-have-a-go-side-source.md`（答 AC#1＋AC#2＋AC#3＋AC#5）
- 性质：**只读普查** ⇒ **AC 框一枚不许勾**、**产码零字节不动**、**不选落点**（落点归编排者）。

## 0. 为什么两枚并一程
两枚问的是同一件事的两面：**"界面上那一维，Go 里到底有没有真值、有没有读者"**。分开派会各跑一遍同样的 grep。⇒ 一张表两节，别写两份。

## 1. 写面（只这些）
✅ 新建 `docs/evidence/s1/180-182-panel-fields-census-c1.md`（骨架先落盘）｜✅ 新建 `.scratch/wisp/probes/180/c1/**` 与 `probes/182/c1/**`｜✅ 追加票 180／票 182 各一段 Progress log（只追加、不改原句、不勾框）。
⛔ `internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量 **一字不动**；起手与终态各跑 `git status --porcelain -- internal/ cmd/`，两条都必须为空，不为空**停下报我**。
⛔ `frontend/**`／`design/**` 不读不写不引不转述（票 182 只判"Go 侧有没有源"，界面怎么画归另一枚会话）；不属于你的脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`design/**`、`docs/evidence/s1/152-*.md`、`probes/183/**`）**不动不提交不评论**。
⛔ 不许跑任何删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；只 commit 不 push、显式 pathspec、禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`。
⚠ **重跑任何既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看它往哪枚路径写数**（就地覆盖型 sink 会洗掉别家票逐行引用的读数，09-28 实测＝`A367`）。

## 2. 票 180 要答的（两问）

**AC#1 归属**：面板窗口尺寸今天由谁决定？现读并逐处指认：① 配置里那枚字段（尺 `grep -n "width\|Width" internal/config/schema.go`）；② 生产码里**谁读它**（尺 `grep -rn "\.<字段名>" --include=*.go internal/ cmd/ | grep -v _test.go`）；③ 窗口/WebView 的实际尺寸今天从哪来（`grep -rn "WebView2\|CreateWindow\|bounds\|Rect" --include=*.go internal/panel/ internal/winsec/ cmd/ | grep -v _test.go`）。⇒ 结论必须是三档之一并给证据：**规格要求生效**／**规格只要求显示不要求生效**／**规格根本没要求**（引 `docs/PLAN.md`／`docs/specs/SPEC-08*` 的**原句**，不许转述成自己的话）。
**AC#4 同族普查**：`internal/config/schema.go` 里**还有几枚**"带 `default` 但生产零读者"的字段？尺要现跑（对每一枚字段名跑一遍读者 grep），**枚数与逐枚名都列出来**；⚠ 别把"测试里有读者"当"生产有读者"（`_test.go` 一律排除）。

## 3. 票 182 要答的（四问）

**AC#1 先数堆**：把"任务监控"那一栏拆成**可点名的堆**（子代理状态／环境信息／git 审查／工作区文件树／内置终端／……），**枚数现量、别目测**。
**AC#2 逐堆三档定性**：**有源**（Go 里已能取到真值，缺载体）／**无源但规格要求**（要新写）／**规格没要求**（别造）。每档**带一条现跑尺**（`grep -rn "<这一维的候选符号>" --include=*.go internal/ cmd/ | grep -v _test.go`），并写明"取到真值要从哪枚函数走"。
**AC#3 逐堆指认归属票**：能并进既有票的**具名并格**（票 145／181／186／187／35…），并不了的**明说"没票认领 ⇒ 要立票"**。⚠ 指认前必须 `grep` 那张票的 AC 标题确认它真写着这件事——**拿一张近亲票填空是这仓犯过两次的错**（`A359`/`A366` 同族）。
**AC#5 雷区清单**：逐堆标"这一堆里有没有任何一枚**由面板发起的批准／写动作**"。⚠ 这一维连着已定案口径（"由面板侧来源的 L2 允许"是禁止项、`Q-49` 那一族钉过五枚门）⇒ **只报现状与风险点，不许提"顺手把批准接上"**。

## 4. 门禁（只读程也取数，按 `A363` 薄规矩**不充当 AC 结案凭据**）
`sh scripts/d22scan.sh`（基线 rc＝0、`ban #8 internal/` **examined=433**）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只有 `G6neg`＝票 178）；`go test -count=1 ./internal/config/ ./internal/panel/`（⚠ `./cmd/wisp/` 本机测不到东西＝票 98，别算进你的红）。⚠ **不要改 `gate-clauses.sh` 那把尺本身**（另有程在跑名册差集）。终态读数取在你自己最后一枚 commit 之后。

## 5. 表骨架（第 ≤6 枚调用内落盘并 commit 第一枚）
十二节：① 起手锚＋写面闸门 ② 票 180 AC#1 归属三问 ③ 票 180 AC#4 同族枚数与逐枚名 ④ 票 182 AC#1 堆数与逐堆名 ⑤ AC#2 三档定性表 ⑥ AC#3 归属指认表（含"没票认领"那一列） ⑦ AC#5 雷区清单 ⑧ **面板要的字段清单**（交编排者转前端会话） ⑨ 本程没测什么（逐名） ⑩ 门禁终态 ⑪ 被拒调用＋零删除自证＋工具调用终值 ⑫ next＝落地腿还缺什么、哪几枚要人先批准。

## 6. 硬顶与增量交付
工具调用**硬顶 45 枚**（两票合派），**到第 32 枚停止新探索**；**每答完一问 commit 一枚**；表没落盘之前不许继续取数；判"必须改产码才能答完"＝停手上报。

## 7. 结束消息必须回给我
① 票 180 的三档结论＋同族"零读者"字段枚数与逐枚名；② 票 182 的堆数、逐堆三档、逐堆归属（含"没票认领"清单）；③ 雷区那一列的逐堆答案；④ 面板字段清单（枚数＋来源）；⑤ 两条 `git status --porcelain -- internal/ cmd/` 为空的自证；⑥ 门禁四数＋红腿名册；⑦ 被拒调用逐条；⑧ 有没有跑删除命令；⑨ 工具调用终值；⑩ next＝还缺什么、哪几枚要人先批准。
