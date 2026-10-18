# 344 — C21 那张 token 对照表仍把**你已废弃的 v2 磨砂玻璃稿**当"权威真相源"（20 行引它＋一枚机器检查硬要它在场）；另：球侧三枚玻璃色值今天⛔ 任何绘制代码读＝死码

- Status: **已立，未派**（2026-10-18 15:4x 编排者立；来源＝机主 15:2x 亲口「磨砂玻璃是我故意放弃的，前端已更新到 v3，千万别改成老版本」＋编排者逐枚现跑；账＝台账 `A1038`）
- ★**本票⛔ 授权任何腿放回 `design/assets/tokens.css`、⛔ 授权任何腿改 `frontend/**`／`design/**`**（那两片归机主的外部 agent，本编队零写面；机主原话已明说不要旧稿）
- 类型：**契约邻接（C21 真相源指向）＋死码登记**。★**改 C21 的"真相源是谁"＝契约变更＝须机主一句明示＋本台账 `A##` 记录**（`AGENTS.md §1.1`／`SPEC-12 §4.1`）；本票只**摆清事实与选项**，⛔ 任何腿自己选。
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 现量（编排者 2026-10-18 15:3x–15:4x 逐枚现跑；引用前重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| ★**旧稿是被故意废弃的、不是丢的** | 机主 15:2x 亲口：「磨砂玻璃是我故意放弃的，因为前端已经更新到了 v3 版本了…千万别改成老版本」；旁证＝`frontend/src/styles/theme.css:2` 逐字 `Panel theme shell (fourth generation, owner pivot 2026-09-25)`、`frontend/scripts/gen-tokens.mjs:1` 逐字 `design-token generator - FOURTH GENERATION (owner pivot, 2026-09-25)` | 现读两文件头 |
| ★**v3 里已经没有玻璃这一族** | `grep -ci glass`：`theme.css`＝**0**、`tokens.generated.css`＝**0** | 两枚文件逐枚数 |
| 被删的那份自称是什么 | `design/assets/tokens.css`（`72b65c61^` 里 287 行）标题逐字 `Wisp / 一缕 — Design Tokens · v2「磨砂玻璃」`＋`唯一样式真相源` | `git show "72b65c61^:design/assets/tokens.css" \| head -6` |
| ★**但对照表仍把它当权威** | `docs/evidence/s1/c21-native-tokens.md:5` 逐字 `权威真相源：design/assets/tokens.css（C21 契约参考实现）`；全表 298 行、**20 行引 `tokens.css`**；`:16` 那行更把它列为"CSS 同源行（40 暗＋38 亮＝78 配色行）"的真相源 | `grep -c "tokens\.css"` |
| ★**机器检查硬要它在场、不许跳** | `internal/ball/tokens_table_test.go:1364` 逐字 `const c21TokensCSSPath = "design/assets/tokens.css"`；检查 5（`:1348` 起，票 74 AC#4 所建）对每一枚引 CSS 的配色行三方对撞；CI 红句逐字 `tokens_table_test.go:1468: read design/assets/tokens.css: … no such file or directory - the CSS leg of this check must never skip` | `sed -n '1348,1372p'`＋CI run `39101060753` 日志 |
| ★**球侧三枚玻璃色值＝死码** | `internal/ball/tokens.go:77/78/79` 定义 `GlassRing`／`GlassHi`／`GlassHiSoft`、`:132-134`（暗）与 `:183-185`（亮）各赋一次值；★**`tokens.go` 之外零读者**（`git grep -l ".GlassRing" HEAD -- "*.go" ":!*_test.go" ":!internal/ball/tokens.go"`＝**0**，三枚各自＝0） | 逐枚现跑 |
| ⚠ 别误判成"整个玻璃族都死了" | 另有一族**活着**：`renderer_windows.go:242/250/292/418` 与 `liquid.go:319` 真读 `v.Glass`（一枚 bool 开关）与 `lk.Caustic`×`GlassCaustic`＝"液面玻璃 look"那套（票 62）⇒ **死的只是那三枚 CSS 同源色值**，⛔ 把整族当死码 | 现读引用点 |
| 这正是那枚检查员自己记载的债 | `tokens_table_test.go:10-13` 逐字：建它就是因为"有一行（`tokens.go CountdownFontPx`）**已经漂到零消费者而任何东西都看不见**" | 现读文件头 |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：上面八行逐枚重跑（一律取具名提交树；旧稿只⛔ **读**、⛔ 放回工作树）。与表不符⇒当场具名更正。★再加一枚：`tokens.go` 里**所有** `Palette` 字段逐枚问"本包外有没有读者"，把死码名册数全（⛔ 只交那三枚就走＝票 336 那族"逐枚锚"的教训）。
- [ ] **AC#1 ★真相源指向这一维＝摆给机主拍，⛔ 任何腿自选**：给三形代价表——**甲**＝改对照表＋检查 5 的指向，认 v3 的真相源（候选＝`frontend/src/styles/theme.css`＋`tokens.generated.css`；★但这两枚是**前端侧生成物**，把它当**原生侧** C21 的真相源要不要成立，得说清）；**乙**＝**原生侧自持**（`tokens.go` 即 C21 原生真相源，CSS 那条腿按票 74 的"范围块"正式作废、检查 5 改成"无 CSS 出处即不查"＝⛔ 造一个假出处）；**丙**＝维持指向、由机主那侧补一份 v3 的 `design/assets/tokens.css`（★那属 `design/**`＝**机主外部 agent 写面**，本编队⛔ 动）。★每形现量：动几枚文件、要不要改 `PLAN.md` 的 C21 正文（＝契约变更）、会不会把票 74／票 12 AC#7 已勾的判语变成假话。
- [ ] **AC#2 死码那一支（可与 AC#1 分开走）**：三枚（或 AC#0 数出的全量）零读者字段——**甲**＝删；**乙**＝留着但把对照表那 20 行的 CSS 列改成具名"（v2 废弃，无消费者）"；**丙**＝留着不动只登记。★⛔ 用"删字段"顺手把表里的行也抹了（那是把债洗成不存在）；⛔ 为凑"有读者"去造一处调用。
- [ ] **AC#3 ⛔ 放宽检查换绿**：⛔ 把"CSS 腿绝不允许跳过"改成"文件不在就算过"（那是把看门狗的嘴捂住＝票 227 `AC#6` 那族"没人钉"的反面教材）；⛔ 删用例、⛔ 裸 `t.Skip`；⛔ 动 `internal/ball/tokens.go` 的**值**（改值＝改视觉＝须机主签收）；⛔ 动 `frontend/**`／`design/**` 一字。
- [ ] **AC#4 颜色与门禁**：`./internal/ball/` 两平台成对（★已知在册存量红 `TestC21TableColourRowsMatchTokensCSS`＝本票要处理的那枚，改前⇄改后逐名点名它转绿或仍红、⛔ 拿别的红抵数）；红名册逐名作差；门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR`（票 311 AC#1 shape B GUARD＝rc=2 射程拒绝，⛔ 读成红）。⛔ `./cmd/wisp` 整包（车道按波次申请）、⛔ `go test ./...`。
- [ ] **AC#5 ★契约变更那一笔的落档形状**：若机主拍甲／乙 ⇒ 编排者先落一枚 `A##`（记文件／行／理由／边界／撤销口令），腿⛔ 自行落档、⛔ 先动 `PLAN.md`；若机主拍"维持指向"（丙）⇒ 本票退回、红归 `design/**` 那侧。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票**七格** AC#0…AC#6 由非实现者裁）、⛔ 改票名／`Status:`／`-done`；写面＝`docs/evidence/s1/c21-native-tokens.md`＋`internal/ball/tokens_table_test.go`（＋必要时 `tokens.go` 的字段名册面）＋`.scratch/wisp/probes/344/r1/**`＋本票尾**一条**；`add` 与 `commit` 两处 pathspec 逐枚手写字面路径（⛔ 从共享 index 反推＝`A1027`／`A1029`）；中文正文⛔ 走 PowerShell（`A1030`）；⛔ push。

## 为什么现在不派

★**AC#1 卡在机主那一拍上**（真相源指向＝契约邻接）⇒ 没有他一句明示，任何腿动它＝越权。AC#2 那支（死码）零契约影响、可先派，但它与 AC#1 同动 `internal/ball/**`＋同一枚对照表 ⇒ ⛔ 拆两腿同波踩同一面。⇒ 释放条件＝机主拍 AC#1 那一形（或明说"先只办死码那支"）。

**撤销口令**：「撤 344」（⇒ 本票退回未派；CI 那枚 `TestC21…` 红照旧在册、⛔ 无人处理）。
