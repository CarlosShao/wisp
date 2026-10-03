# 212-v1 验收裁决（第一实例，非实现者；碰撞档案版）

> **碰撞声明（先读这个）**：`.scratch/wisp/probes/212/v1/verdict.md` 在本实例工作期间被**并发的"212-v1 第二任"实例**覆写并先行 commit（`234e7301`，21:25，commit message 自称"第二任；复用前任基线"）。本实例即它所称的"前任"——**本实例未死**，21:08 起手、读数全程在手。为防继续互覆，本实例全量裁决落本件；`verdict.md` 归第二任续写，本件只在其尾部留一行指针。两实例并发跑同一票＝编排侧重复派单（A578 只记派一枚 212-v1），归编排者裁撤哪一枚；**本件是其一的完整读数，判语自足，不依赖 verdict.md**。

## §0 起手锚

- 会话起手锚：`dffd9456d25ac11559f9f476b5afa848a9b0526b`（分支 dev）；起手 `2026-10-03 21:08 +0800`，本件落盘 `2026-10-03 21:5x +0800`。
- 实现件状态核验：`tools/d22scan/main.go` 最后触碰 commit＝`5e8748b3`；本腿全程对 `tools/ internal/ cmd/` 跟踪树零改动；起止 md5 全等：main.go `51f0f69f5391462053b5a9400ec1529b`、selftest.go `bbc0bd946e1117e437c58939e867112e`、selftestsamples.go `b5d8fb8ffa115d4600f5038ff4a2a0f7`。
- 工作树有大量**他人**未提交痕迹（design/assets 删除、probes 杂件）——零触碰零还原。transient `tools/d22scan/main.go.prime`：起手 status 有、session 中段消失，非本实例创建，来源判不动（§7）。
- **基线复用声明**：未复用任何前任读数（本实例即第一实例）；`gate-baseline.txt`（21:1x 建，untracked）是**本实例**自跑的门禁全量读数。第二任骨架声称"前任死于服务错误留下 gate-baseline.txt"——**与树不符**：该文件是本实例产出且从未被 commit（`git ls-files` 实证），"死"不存在（§6 #9）。
- **派单指向纠错**：派单"票面『编排者翻勾节』后的 A574/A578 段"——票面（56 行）无该段；A574/A578 在 `docs/reports/pending-and-issues.md`（A574=11256 行、A578=11283 行），已逐条对（§6 #1）。
- 纪律：只建 `.scratch/wisp/probes/212/v1/**`；产码零字改；勾框零碰；frontend/**、design/** 零读零引；commit 显式 pathspec、不 push。

## §1 恒真两问（`mut/`，`go build -overlay` 突变台；跟踪树零改动，md5 复量还原）

**① 摘掉 ban #9 发射点**（删 `for _, cg := range f.Comments {…}` 段＝`mut/main-noq9.go`）：
- 真仓全扫：exit **0**、verdict **clean**（`mut/noq9-out.txt`）——摘门后本仓零红，门"看起来"照样绿 ⇒ 本仓对 ban #9 的全部红都由这一段发射。
- 喂**原版名册**跑 -self-test：exit **1**、35/36，唯一 FAIL＝**ban #9 phantom-citation 的 wantRing 用例**（`internal/probe/cites.go` 引未种的 `docs/evidence/s1/212-citation-ruler.md` ⇒ "did NOT ring"），bans #1-8 用例全绿（`mut/noq9-selftest.txt`）。
- 喂**突变版名册**跑 -self-test：exit **2** FATAL，2 枚 HOLE＝"ban #9 documented but no emission site can print it"＋"sample names a tag no emission site can print"（`mut/noq9-selftest-mutroster.txt`）——名册审计独立抓到断肢。
- **答**：具名红用例恰好 1 枚＝ban #9 ring 向；名册审计路径另有 2 枚 HOLE。判据有牙、牙的缺失自己会叫。

**② 反形**（"存在的引用必红"：`if exists {continue}` → `if !exists {continue}`＝`mut/main-inv.go`）：
- 真仓全扫：exit **1**、**160 枚** finding **全部** `[phantom-citation]`、bans #1-8 零枚（`mut/inv-out.txt`）——`docs/reports/pending-and-issues.md`、`internal/tools/bridge.go`、`internal/ball/statevisual.go` 等真文件被"打红"。
- 喂原版名册跑 -self-test：exit **1**、34/36，ban #9 **两向同 FAIL**（ring 向未种文件不响；silent 向种下的 `docs/readings.md`/`internal/probe/roster.md` 被问罪）（`mut/inv-selftest.txt`）。
- **答**：反形不是"也全绿"，是**全红**且被自测双向拒绝。派单问句的"会不会也全绿"答案为否。

## §2 symRefRe 三发探针（fixture `fix/r2probe/`，读数 `fix/scan-pristine.txt`、`fix/scan-ellipsis.txt`）

| 探针 | token | 读数 | 判 |
|---|---|---|---|
| (a) | `internal/tools.Result.AppliedSteps` | **静默** | 与实现声称一致；所指为真字段（bridge.go:575），排除语义正确 |
| (b) | `internal/tool`（非 .go 后缀、无点无大写段） | **红**（p2 问罪） | symRefRe **不**排除它：`^…(组)*\.[A-Z]` 需一个 `.大写`，无点即不匹配；它消红靠 tools.go:89 **改注释**，不靠正则 |
| (c) | `internal/Build/x.go`（大写目录段、路径不存在） | **红**（p3 问罪） | symRefRe **不误排除**：`/B` 处组断、`internal` 后无 `.` ⇒ 不匹配 ⇒ Stat ⇒ 缺失即红。读正则推演与实测一致：**symRefRe 不会把真文件引用误排除**（小写盘面＋`.大写` basename 才可能，见下）；对照 p4 小写路径同红、p6/p7 存在目标静默 ✓ |

**误排除面普查**（`symref-sweep.txt`，全 5129 tracked 路径）：repoPathRe 前缀 ∧ symRefRe 交集＝**0 枚**；symRefRe 全局命中仅 `docker/builder.Dockerfile`、`docker/mockllm.Dockerfile`——`docker/` 不在 repoPathRe 前缀表 ⇒ 不可达。**今天无任何真文件引用被误排除**。main.go:846-848 "a .go file is never excluded" 经推演＋p3 实测为真。NTFS 大小写不敏感的 Stat 残洞记 §7。

**pre-image 普查**（`preimage-phantom-census.txt`，`5e8748b3^` 树 266 产码 Go 文件，含尾随注释）：distinct phantom token **11**＝真② **8**（`scripts/check-pathclean-ban.sh`、`tools/gate.go`、`tools/bridge.go`×2 site、`internal/tool`、`internal/engines`、`tools/agent/cmd`、`internal/provider`、`docs/evidence/s1/62-`）＋API 形 **3**（`internal/buildinfo.Name`、`internal/proc.WithRegistry`、`internal/tools.Result.AppliedSteps`，symRefRe 排除、实测静默）。⇒ commit 说"首发 10 枚""10 枚真②全部改注释"两个 **10** 均与树不符（§6 #2）。

**附**（AC#3 关键探针）：`docs/evidence/s1/152-...-accept-r1.md`（省略号缩写形，p10.go）在产码注释里**被问罪**——`…`/`...` 在 repoPathRe 字符类内 ⇒ 整 token 匹配 ⇒ Stat 失败 ⇒ 红（`fix/scan-ellipsis.txt`）。见 §5。

## §3 AC#5 八枚 ban 读数对账（`gate-baseline.txt`，exit 0 verdict clean）

| scope | A578 记 | 本腿实跑 | 判 |
|---|---|---|---|
| bans #1-5 internal/ | 228 | 228 | 一致 |
| bans #1-5 cmd/ | 38 | 38 | 一致 |
| ban #6 frontend/ | 85 | 85 | 一致 |
| ban #7 internal/tools/ | 23 | 23 | 一致 |
| ban #8 design/ | 39 | 39 | 一致 |
| ban #8 frontend/ | 85 | 85 | 一致 |
| ban #8 internal/ | 498 | 498 | 一致 |
| ban #8 cmd/ | 96 | 96 | 一致 |

**八枚逐名零漂移**；headline 266 同跑实录；包测试 stage PASS=34 FAIL=0（含 `TestSelfTestEntryPassesEveryCase`）；净版 -self-test **36/36 clean（20 ring/16 silent）**（`selftest-clean.txt`）与 A578/commit 记载逐字一致。

## §4 注释更正逐枚复核（diff 实测：**8 文件 9 hunk**；票面称"10 枚"见 §6 #2）

| # | 文件:处 | 更正后的话 | 真值核查 | 判 |
|---|---|---|---|---|
| 1 | internal/risk/pathresolver.go:28 | "enforced by tools/d22scan's pathresolver-bypass ban; see tools/d22scan/main.go" | ban #2 tag 字面在 main.go:726（名册 :15），件在 | **真** |
| 2 | internal/agent/approval/pending_read.go:43 | "ticket 17's frozen-contract note (enforced by tools/d22scan's pathresolver-bypass ban since the scanner landed)" | 半句字面真；但**所指迁移**：原句要指"冻结注所在文件"，真身＝`internal/tools/gate.go:13`（"ticket 17's frozen-contract note" 就在那），补前缀逐字修即可；"冻结注由 Clean/Abs ban 执法"是语义嫁接 | **真但带条件**（更诚实形＝internal/tools/gate.go） |
| 3 | queue.go:38 | "internal/tools/bridge.go routes cancels by orDefault(CorrelationID, TaskID)" | bridge.go:519/526/590/772 实证 | **真** |
| 4 | queue.go:468 | "internal/tools/bridge.go's orDefault(why, …) on the reject branches" | bridge.go:445/448/461/464 实证 | **真** |
| 5 | internal/agent/tools.go:89 | "error class internal-tool" | **`internal-tool` 不存在于 D37 枚举（PLAN.md:2778-2796）与 observe 常量（ClassTool="tool"/ClassInternal="internal"）**；实质真：loop.go:743 → guard.go:271 → 默认 `ClassInternal`、模型不可自纠（loop.go:748 旁证自认） | **仍存疑**（斜杠触发修掉了、类名仍是造的；诚实形＝D37 class internal） |
| 6 | internal/ball/statevisual.go:68 | "the ticket-62 evidence tables (the ticket-62 evidence tables under docs/evidence/s1) were measured at 96 DPI" | 62-visual-spec-draft.md:6 "3440×1440 @ 96 DPI…没有一条像素数字是在 >96 DPI…上量的" 实证；docs/SLO.md A.2 在 | **真**（括号内同短语重复两遍＝编辑瑕疵） |
| 7 | internal/risk/provenance.go:63 | "no *plugin.DisposalScope reaches the plugin agent command in production today" | 实质真（grep cmd/ `DisposalScope`＝0 文件）；但 "the plugin agent command" **无法机读证伪**（internal/plugin/doc.go 无 agent 字样、cmd/ 无该命令）——比原 `tools/agent/cmd` 更模糊 | **仍存疑** |
| 8 | internal/tools/bridge.go:252 | "D37 class internal-provider" | 同 #5：`internal-provider` 不在枚举与常量；实质真（loop.go:743-748 记不可自纠类） | **仍存疑**（诚实形＝按实际 wrap 定 internal 或 provider） |
| 9 | cmd/wisp/models.go:20 | "engines directory under internal/" | internal/ 列表实测无 engines | **真** |

**小结：5 真＋1 真但带条件＋3 仍存疑**（三枚的斜杠/路径 token 已不再触发 ban #9，但"改成实话"只完成一半——实质半句真、名目半句仍是造的/糊的）。

## §5 AC 格判语（⛔ 勾归编排者，本腿只给判语）

- **AC#2（落点）＝成立但带条件**。判据落 `tools/d22scan`、独立第五形（ban #9 独立编号/tag/发射点，名册 1..9 连续由自测 Rule 0 钉），射程只产码注释面＝与 A574 初裁逐条相符；§3 证明既有 8 枚读数零动；§1 证明有牙且断肢自报。**条件**：票面 AC#2 原文"由非实现者裁"——落点初裁（A574）与实现（A578 代笔）同为编排者一人，裁决程序违反票面；本验收腿（非实现者）现予**追认**补正程序瑕疵。附带：ban #9 搭 bans #1-5 的 walk，self-report 那两行 scope 标签 "bans #1-5 internal/" 字面不再穷尽该 walk 服务的 ban 集（cosmetic 非违规）。
- **AC#3（③形只登记不问罪）＝不成立**。③的省略号形（票面现量点名的拼法）在**产码**注释里被 ban #9 打红（§2 附行实测）——判据与名册第 9 行及 scanGoFile 自述"shorthand forms are a separate prescription, not a violation"**自相矛盾**。今天全仓 clean 是**侥幸**：该形唯一现存标本在 `cmd/wisp/slo_report_144_windows_test.go:851`＝_test.go，而 ban #9 搭的 walkGo 跳过 _test.go。silent 样本（cites-ok.go）只覆盖散文式缩写（"the 152 ruler"/"…那一族"，本无路径 token 可匹配），**不覆盖**省略号路径形。不成立的精确范围：**仅省略号-路径形**；散文缩写确实不问罪。修法（把省略号字符移出 repoPathRe 或匹配后剥除）属**仪器射程变更**须走批准，本腿不动产码一字。
- **AC#5（不许伤）＝成立**。八枚读数逐名零漂移（§3）；正控两向三重实证：净版自测 36/36；摘门 ⇒ ban#9 ring 向 FAIL（35/36）＋名册审计 FATAL×2；反形 ⇒ 两向同 FAIL（34/36）＋真仓 160 枚全红。判据对净仓零新增红、对突变体红得有名字。

## §6 推翻清单

1. **派单**："票面翻勾节后的 A574/A578 段"——票面无此段，在台账（§0 已补偿读取）。
2. **commit/A578 计数**："首发 10 枚分类"＋"10 枚真②全部改注释"——pre-image 普查＝**11** distinct token（8 真②＋3 API 形）；diff 实改 **8** token（9 hunk/8 文件）。两个 "10" 与树都对不上。
3. **commit 定性**："tools/gate.go …错拼类名"——错类：那是**少写 internal/ 前缀**、真身存在（internal/tools/gate.go:13 即票 17 冻结注）；最诚实修法（补前缀）未用。
4. **main.go:802-805 注释**：symRefRe 举例含 "class internal/tool"——该 token **不**被 symRefRe 匹配（无点，§2(b) 实测红）；仪器自述与行为不符（本票治的就是撒谎注释）。
5. **票面 AC#3/A574 初裁**："③形只登记不问罪"——省略号形被问罪（§5）；clean 依赖 _test.go 豁免这个未声明的巧合。
6. **selftestsamples.go silent 样本 note**："shorthand is not convicted"——只对散文缩写成立，对票面点名的省略号拼法不成立（覆盖面声明过强）。
7. **statevisual.go 新句**括号内同短语重复两遍——编辑瑕疵非假话。
8. **碰撞**：第二任骨架（`234e7301`）声称"前任死于服务错误，留基线 gate-baseline.txt"——本实例未死；gate-baseline.txt 是本实例 21:1x 产出且从未 commit，"复用"实为读了同名未跟踪文件而误记出处。其二任 §4 表把 `cmd/wisp/config_readers_255.go`、`cmd/wisp/resident_approval_*246*_test.go`（**258-r1** 的 hunk，A578 记在 258-r1 名下）计入票 212 的"十枚注释更正"——混入；212 的更正集是 8 文件 9 hunk（§4）。其表项 9/10 的"更正前→更正后"对 212 而言不存在对应 diff。
9. **编队账**：A578 记"派 212-v1＋258-v1 两枚验收"——实际 212-v1 起了**两枚并发实例**（本实例 21:08、第二任 21:19），同一 mandated 落点互覆；归编排者裁撤哪一枚、账面记重复派单。
10. **（小）A574 记 census"产码侧真②=1 枚"**与 commit"首发 10 枚"之间的差（9 枚）无解释——commit 说"全仓扫出的其余 8 枚"，两数合并仍对不上 11 的普查值（口径差可容，去向未留痕）。

## §7 判不动／量不到

- "首发 10 枚"的口径无法复原（读数在已死的 r1 会话里）；本件能给的最硬数＝`5e8748b3^` 树 266 产码文件、11 distinct token（含尾随注释位）普查与 8+3 分解。
- transient `tools/d22scan/main.go.prime`：起手 status 有、中段消失，非本实例创建——共享工作树另有在飞者，来源判不动。
- Windows/NTFS 大小写不敏感 Stat：注释引 `internal/Build/x.go` 而 `internal/build/x.go` 存在时 os.Stat 命中 ⇒ 静默。今天仓内无此形（sweep 交集 0），实害量不到，记已知形状。
- ban #9 搭 walkGo 的 _test.go 豁免是既有行为非本票决策；"③形今天不打红"悬在该豁免上——未来任何产码注释里的省略号引用都会现形（判据与 AC#3 的冲突点，非本腿可修项）。
- 第二任实例仍在飞（r3-* 文件持续更新中）：其最终判语与本件的分歧无法在本件内裁定，留给编排者对读。

## 附：本实例产出物（本 commit，显式 pathspec）

`verdict-first-instance.md`（本件）；`gate-baseline.txt`、`selftest-clean.txt`、`symref-sweep.txt`、`preimage-phantom-census.txt`；`mut/`：probe212.py、main-noq9.go、main-inv.go、ov-noq9.json、ov-inv.json、noq9-out.txt、inv-out.txt、noq9-selftest.txt、inv-selftest.txt、noq9-selftest-mutroster.txt；`fix/scan-pristine.txt`、`fix/scan-ellipsis.txt`、`fix/r2probe/**`（探针 fixture 种子）。突变体 .exe 不 commit（probe212.py 可重建）。

**终态 status：交件**。恒真两问①红名＝ban#9 ring 向 1 枚（＋名册 HOLE×2）；②反形＝160 枚全红非全绿。symRefRe：(a) 静默 (b) 红（正则不排除、靠改注释消红）(c) 红（不误排除）；误排除面 0。AC#5 八枚零漂移。十枚更正＝8 文件 9 hunk：5 真＋1 带条件＋3 存疑。AC#2 成立但带条件（追认补正程序瑕疵）、AC#3 不成立（省略号形被问罪）、AC#5 成立。勾框零碰。
