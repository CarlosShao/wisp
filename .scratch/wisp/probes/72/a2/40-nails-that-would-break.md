# 72-a2 — 会碰坏哪些既有钉（逐枚带尺；尺名与射程目录都写在条目里）

射程说明：凡"引用名册"类尺，射程目录＝`docs/**`、`.scratch/wisp/issues/**`、`tools/**`、`internal/**`、`cmd/**`；
**⛔ 算进账的＝`.scratch/wisp/probes/**` 与 `.scratch/ci-logs/**` 那批历史读数件**（它们是一次性凭据，⛔ 会被改名追改，只⛔ 再拿它们当"当前名册"）。
名册行数一律配一把逐行打印的尺对拉（AGENTS §1.2 那枚 `\r` 陷阱：本程所有取名尺都先 `tr -d '\r'`）。

---

## 1. ★★票 72 自己 `AC#5` 的**逐字节自证尺**——本腿认为这是下一程**最该先裁的一枚**

票面 `:82` 逐字＝
> `git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go` **rc=0**
> （和票 18 那版逐字节相同），没加 `t.Skip`、没加 build tag、锚点拼写没"对齐"、期望值仍是 `ClassA`

- 本腿现跑＝**`rc=0` 成立**（件 `logs/ac5-byte-ruler.txt`，两发并排：junction 文件 rc=0／anchor_spelling 文件 rc=1，后者是票 72 自己在 `a04d3e2` 之后新建的，⛔ 被钉）。
- ★撞车点＝**甲族 7 枚里有 4 枚住在这枚被钉的文件里**（#4/#5/#6/#7，住址尺见 `10-seven-cases.md`）。⇒ **修那 4 枚⛔ 碰这枚文件就修⛔ 了**；一碰，`AC#5`（一枚**已勾**的框）的书面取证方式当场过期。
- ⚠ 这⛔ 是"所以⛔ 能修"。这是**一道必须人工过的闸**：`AC#5` 的实质是"⛔ 放宽判据"，而它用的**仪器**恰好就是那枚要改的文件。处置只有两条，都得编排者点头：
  (a) 由编排者在票 72 尾部**追加一节**（⛔ 改 `:82` 原句）改写"取证尺"的表述——例如换成"`AC#5` 的实质判据改由'两台面 16 枚同跑＋红位集合⛔ 变'来作保"；
  (b) 把纠正落在**新文件**里、把那 4 枚原样留着（⇒ 别名台面上它们继续红，CI 永远⛔ 绿，那⛔ 是选项）。
- 顺带：同一枚文件的**注释**也在替这句话背书——`pathresolver_anchor_spelling_windows_test.go:13-16` 逐字＝`These cases sit NEXT TO pathresolver_junction_windows_test.go, which is the failing instrument and is deliberately left byte-for-byte untouched: its expected values, its anchors and its build tags are the measurement, not the thing under edit.`
  ⚠ **注释里的这句现在与盘上事实冲突**：它说 junction 文件的"anchors"是测量物，而那 4 枚的 anchors 恰恰就是本票族要改的东西。注释一律当待验断言——这一枚**已验为过期**（尺＝上面那把 byte diff ＋ 本腿两台面颜色）。

## 2. ★★`TestOverrideOnlyAcceptsTheResolvedForm`（`anchor_spelling:157`）＝放行侧窄规则的正钉

- 它是唯一防"把 `overrideApplies` 改回遍历全部形"的仪器；票面 `:90-92` 与 `:241-243`、log `:171-173` 三处记它的变异**实测过**转红。
- 撞车点＝甲族 #3（`:237`）与 #7（`:291`）最省事的修法正是动它守的那条规则。⇒ 见 `30-answer-B-line.md` 形 IV 的具名判语。
- ★本腿的量的一枚**加分读数**：这枚用例**在别名台面上自己就是绿的**（件 `logs/aliased-all16-bothfiles.txt`，`--- PASS: TestOverrideOnlyAcceptsTheResolvedForm (0.06s)`）⇒ 那条窄规则在别名世界里**行为正确且已被验证**，⛔ 需要任何"为它放宽"的空间。

## 3. `TestCanonicalInputGainsNoSecondForm`＝票 72 `AC#5` 的**爆炸半径钉**本身

- 票面 `:93-98` 把它从"反驳通道"降级成"爆炸半径证明"，并逐字写"这句话现在是可跑的，⛔ 是口头承诺"。
- 撞车点＝若按 `30-answer-B-line.md` 形 III 里那两条被禁的修法（删名册元素／改阈值），被稀释的正是**这一句话的射程**，而⛔ 是一枚断言的松紧。
- ⚠ 本腿量到：这一枚是 7 枚里唯一**⛔ 携带任何判级**的（它一句 `Classify`/`Gate` 都⛔ 调）⇒ "它的红能不能当'判级被证伪'的证据"这一问，答案**按构造＝⛔ 能**；任何把它说成"票 72 的修复在这里失效了"的表述都属过度声称。

## 4. 票 72 `AC#3①` 与 round-2/round-3 的**变异名册**（三处逐枚点名）

- 票面 `:70-72`、log `:139-141`、log `:146-150` 三处把 `TestClassifyAnchorSpellingIsNotVerdict`／`TestCanonicalInputGainsNoSecondForm`／`TestAListWinsWhereBothTablesHit` 与另两枚（`TestClassifySpellingInvariance`／`TestClassifyFailClosedWhenSpellingUnprovable`）**写死成"退 `formsOf` ⇒ 这三枚红"的名册**。
- 撞车点＝改这三枚的**期望侧／前提侧**之后，那一发变异**必须重跑**才⛔ 是"名册过期"。⇒ 本腿按结构预判（⛔ 仪器，⛔ 授权跑）：#1 换成候选乙"入参也走 `alias` 形"、#3 换成已解析形入参 ⇒ 两枚仍应对"退 `formsOf`"有信号；**#2 有真风险**（它唯一的判据就是那个 `1`，把输入的短拼换成已展开拼法之后，`formsOf` 退化它**可能反而更绿**）。⇒ 落地腿**必须**为 #2 另配一枚能响的钉，或具名交"`AC#3①` 名册里这一枚⛔ 再算"。
- 同一批证据还钉着 `TestPathResolver*` 那几枚"⛔ 算占位"（票面 `:13-16`），⇒ 改 #4/#5/#6 的比较行⛔ 动 `:137`/`:140`/`:166`/`:184` 这四行判级（详见 `30-answer-B-line.md` 形 I）。

## 5. 台账与裁决表里的**逐名引用**（改名＝让它们指空）

尺＝`git grep -n -E '<7 枚名>' HEAD -- docs .scratch/wisp/issues`；射程＝`docs/**` ＋ `.scratch/wisp/issues/**`；**含注释行**（那是 .md）；HEAD blob。
- `docs/reports/pending-and-issues.md`＝**8 处命中**（行号快照 `:2845`、`:4093`、`:4094`、`:5439`、`:6504`、`:12729`、`:15336`、`:15863`）。台账规矩＝`A##` **只追加⛔ 删** ⇒ ⛔ 能改这几行原句，只能追加更正。
- `docs/evidence/s1/ci-read-a1fd5bf-r1.md`＝**3 处**（窄射程尺：`git grep -c -E '<三枚名>' HEAD -- docs/evidence`）；`docs/evidence/s1/141-ac2-stock-inventory-r1.md`＝1、`141-q46c-accept-r1.md`＝1、`docs/evidence/s3/19-adversarial-fix-round.txt`＝2。
- 票面本身：`72-…md`＝**11 处**、`307-…md`＝2、`236-six-cells-that-only-surface-at-the-reading-layer.md`＝1、`75-linux-path-canonical-shape-done.md`＝3。
- ⇒ **判语句：名字⛔ 能改**（这 7 枚的名字是 4 本账的公共键）。改的⛔ 是名字，是**体内**。
- ⚠ 一处**别拿来当甲族名册**的台账行（本腿逐字读了 `:4093-4095` 三行）：那一行的名册里除这 7 枚外还点着 `TestComposedGateBlocksAWriteForTwoSeconds`、`TestTicket101*` 三枚、`TestSync*` 五枚 ⇒ 它是**另一批（更早的）CI 红名册**，⛔ 等于甲族 7 枚、也⛔ 等于乙族 5 枚。⇒ 引用台账那一行"证明甲族是 7 枚"＝**用错名册**；甲族名册的权威出处＝票 72 `:250-252` 追加节与票 307 追加节 `:75`（同族先例＝票 307 追加节 `:61` 那枚⛔ 存在的 `TestAListWinsBothTablesHit`：⛔ 存在的名字⛔ 能追认任何格）。

## 6. `internal/tools/mode.go:98` ＝`risk.Gate(c, overrides)` 的**生产 caller**（★本腿新增的一枚钉）

- 尺＝`git grep -n -E '\bGate\(|\bClassify\(' HEAD -- internal cmd tools | grep -v _test.go`＝**6 行**；生产 override 表那一支＝`internal/tools/mode.go:91 overrides := b.confirmations()` ＋ `:93 c, err := b.paths.Canonicalize(p)` ＋ `:98 d := risk.Gate(c, overrides)`。
- 撞车点＝**任何去动 `overrideApplies` 的"修法"现在⛔ 再是"只动测试"**：它会改生产分档（`BlacklistNote.AlreadyUnlocked` ⇄ `Unlockable`，`internal/tools/mode.go` 的 `switch d.Class`/`if d.Allow` 那一段）。
- 同时它证伪票面 `:245-247` 那句"今天还⛔ 生产 caller 传 override 表"⇒ 见 `20-answer-A-cause.md` A.4。

## 7. build tag／档位／分母（票 301 包级 guard D ＋ 票 302 射程 ＋ 票 306 CI core 档）

- 两枚文件都是 `//go:build windows`（第 1 行，`git show HEAD:<path> | sed -n '1p'`）。⇒ 加 tag／搬文件⛔ 能把它们移出 `test-windows` 那一步，⛔ 票面 `:205` 直接禁止（"路径拼写⛔ 是平台 API 限制，是我们的 bug"），⛔ 票 307 追加节 `:133` 已经把"给那批加 tag 会把它们从 core 档分母里摘掉"写成先例。
- ★一枚**具体的仪器边界**（本腿现量，射程＝`tools/d22scan/main.go`）：`walkGo` 在 `main.go:687`（`strings.HasSuffix(path, "_test.go")` 那支）**把 `_test.go` 排除在 ban #1–#5 之外**，**只有 ban #8（emoji）的射程**明确"comments and `_test.go` INCLUDED"（`main.go:40-46` 原文）。⇒ 后果两条，都对落地腿有用：
  1. 落地腿⛔ 在这两枚文件里塞进任何 emoji 形状（含注释！`✓` U+2713、`≤` U+2264 那一族），⛔ 是"测试文件⛔ 被查"；
  2. `filepath.Clean|Abs` 那一条（ban #2 `pathresolver-bypass`）**⛔ 会因测试文件而红**，但 `allowlist.txt` 只有 **7 行**且**第 6 行就是 `internal/risk/pathresolver.go`**（"C26 sanctioned lexical step"）⇒ **⛔ 动那枚冻结件 ＆ 落地腿⛔ 在 `internal/risk` 任何新文件里引入 `filepath.Clean|Abs`**（AGENTS §1.1 明令⛔ 动 `allowlist.txt`，而规则本身⛔ 因为仪器⛔ 看测试就作废）。
- 新建 `.go` 会挪 `gofumpt`／`d22scan` 的分母（本腿⛔ 建，落地腿才允许；先例＝票 307 `51a0a0f2` 新建 `syncdirs_alias_307_windows_test.go` 时配了 `17 3` 的产码改动一起交）。

## 8. ★**三枚取短名的手已经在同包里**——⛔ 允许落地腿新造第四枚

尺＝`git grep -n -E 'func (getShortPath|tryShortPath|alias307Pair)\(' HEAD -- internal/risk`＝**3 命中**。
- `pathresolver_junction_windows_test.go:58 getShortPath`（`syscall.GetShortPathName`；失败→`""`→上层 `:128` **`t.Skip`**）
- `syncdirs_redteam_windows_test.go:29 tryShortPath`（`windows.GetShortPathName`；`(string, bool)`）
- `syncdirs_alias_307_windows_test.go:38 alias307Pair`（票 307 刚落地的；⛔ Skip，取⛔ 到短形**即 `t.Fatal`**）
⇒ 票 115 定式（"种不出即 Fatal"）只有第三枚守。落地腿的正确姿势＝**复用第三枚**，⛔ 顺手把第一枚的 `t.Skip` 改掉（那一枚＝票 307 边界节 `:34` 明写"登记、⛔ 本票顺手修"、并指向票 72 的射程⇒ 它属票 72 **的另一格**，⛔ 塞进"改这 7 枚期望侧"这一程，否则就是票 307 `AC#5` 那句"⛔ 同一件事记两本账"的反面教材）。

## 9. `TestResolvePerCallBudget`（墙钟那一族）

- 票 72 log `:155-162` 已经把这枚钉成"**红的时候先做纯净树 A/B 再定罪**"的形状。落地腿若给这族用例加"每次向 OS 要短形 ＋ 建目录"的夹具（候选乙），**同包总耗时会上去**；`Resolve` 本体⛔ 动 ⇒ 那枚预算用例本身⛔ 该被影响，但共树编队下它可能**因环境而红**。⇒ 处置照 log `:161` 那句定式，⛔ 改判据（票面 `:162` 原话"我⛔ 动它（改判据=放宽门，⛔ 是本票的活）"）。

## 10. CI 那批色要**同批作差**（编排者自己的账）

票 307 追加节 `:183` 逐字＝`推送⛔ 排（A813 仍⛔ 用，51a0a0f2 起那批 CI 色同批作差）`，且那批里点名票 300／303／111 欠的 CI 色。
⇒ 本票落地笔一旦进 `internal/risk`，那批作差里就多一枚自变量。⇒ **落地腿⛔ push、⛔ 单独取 runner 色**；runner 色归编排者那一发（`docs/evidence/s1/ci-read-a1fd5bf-r1.md` 与 `ci-runner-readings-2026-09-21.md` 两枚裁决表就是那批的旧账）。

## 11. 三枚冻结件与 SLO/golden（⛔ 射程，写在这里只为防止落地腿顺手）

`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`、`internal/observe/thresholds.go`、`tools/d22scan/allowlist.txt`、D43 转移表、C1–C32、D1–D47、SLO 阈值、golden。
- ⚠ 一枚**同名相近**的坑：`internal/perm/ticket90_persist_test.go` 与 `internal/panel/l2_grant_boundary_test.go` 都不在本票射程，但 `L2`／`grant`／`persist` 那几个词会出现在本票族的新注释里 ⇒ ⛔ 因为"看起来相关"就跨包动手。

## 12. 在飞的 `307-v1`（本腿开工时它在飞，交件时已落）

- 起手锚时 HEAD＝`c6fc78f0`，`c04ced26` 是 `307-v1` 的起手锚；本腿工作期间 HEAD 走到 `f6f4eca8`，其中 `ed1d2bad` ＝ **`307-v1` 的 `AC#1`＋`AC#2` 判语**（`git log --oneline -5` 现量）。
- 那枚 commit 标题里已经写着"(d) **3/5 无信号裁归票 72**"⇒ **`307-v1` 已经把一枚读数指回本票族**（与 `10-seven-cases.md` 各枚第 4 问那一列直接相关）。⚠ 本腿⛔ 读那枚件的正文（它落在 `probes/307/v1/**`，本腿射程外、且它是裁决件⛔ 是料），只把它**存在**与**指回票 72**这一件事登记，**下一程必须自己去读它**。
- 期间产码面⛔ 动过：尺＝`git diff --name-only 7db4db4b HEAD -- internal`＝**0 行**（本腿锚点到自己交件之间，`internal/` 一枚文件都没变）⇒ 本程所有读数仍对 HEAD 有效。
