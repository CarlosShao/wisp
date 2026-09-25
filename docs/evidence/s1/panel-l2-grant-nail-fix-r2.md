# panel-l2-grant-nail-fix-r2 —— 实现者：把验收 r1 造出来的三枚洞补上（(a)-(e)）

- 时刻 / 锚点：2026-09-25 10:3x +08 开工，进场 `dev` @ `310816f`。
- 工单来源：`docs/evidence/s1/panel-l2-grant-nail-accept-r1.md` §7.3（三件修法 (a)(b)(c)+F-5）与 §7.4
  （文件头补 "nothing is wired today" + 改掉 `:604` 指向不存在测试的引用），
  总裁＝**退回（附条件入账）**、`Q-49` 不勾。
- 身份：**本程是实现者**，只写 `internal/panel/l2_grant_boundary_test.go` 与本件。
  生产码零字节（见 §9 的 `git diff --numstat`）。`frontend/**`、`design/**`、`tools/d22scan/**` 未碰。
- 纪律：变异一律落在**仓外副本** `D:\tmp\panel-l2-nail-fix-r2\headcopy`，仓库树从未进入变异态；
  临时件只建不删，未 `rm`、未 `checkout`/`reset`/`stash`/`clean`、未 `add -A`、未 push、未 `--amend`。
  台账 `docs/reports/pending-and-issues.md` 一字未写，`Q-49` 未勾。

---

## §0 变异台先跑"改之前"：五发全部复现验收的读数

### §0.1 平台与可比性

派单给了两条腿（仓外副本 / 就地 `cp` 还原 + `hash-object` 作证）。我选**仓外副本**那条，
因为它让仓库树在任何时刻都没有处于变异态，而不是"平均几秒没处于"。

```
cd "D:\work\workspace\projects plans\Wisp" && tar --exclude=./.git --exclude=./third_party \
  --exclude=./build --exclude='./frontend/node_modules' --exclude='*/node_modules' -cf - . \
  | (cd /d/tmp/panel-l2-nail-fix-r2/headcopy && tar -xf -)      -> 3.3s, 155M
```

副本与被验物的逐字节同一性（`git hash-object`，四枚全等，且与验收 §2 记的两枚同值＝两路互印）：

```
working tree  bridge.go                   d2cd6362ecc6941a0cee58073a2bf8ab6669a590
HEAD          internal/panel/bridge.go    d2cd6362ecc6941a0cee58073a2bf8ab6669a590
headcopy      internal/panel/bridge.go    d2cd6362ecc6941a0cee58073a2bf8ab6669a590
backup        bridge.go                   d2cd6362ecc6941a0cee58073a2bf8ab6669a590

working tree  l2_grant_boundary_test.go   5b618d14b742f60cac22d2e58a1f2863dd311c63
HEAD          …/l2_grant_boundary_test.go 5b618d14b742f60cac22d2e58a1f2863dd311c63
headcopy      …/l2_grant_boundary_test.go 5b618d14b742f60cac22d2e58a1f2863dd311c63
```

副本基线复跑（`go test ./internal/panel/ -count=1 -v`，`out/baseline-copy.txt`）与验收 §0.1 同值：

```
RUN=180  顶层 PASS=100  子 PASS=78  FAIL=2  SKIP=0  panic=0  distinct=90
--- FAIL: TestC21DesignTokensFourWayAgree   (先存在留红，A208③，本程未修未跳)
```

⇒ 下面每一发的"没红"都是在这把尺上读的，不是推的。

### §0.2 探针（不是我推理，是打印出来的运行期真值）

每发变异同时生成一枚 `internal/panel/zz_fix_l2probe_test.go`（`_test.go` 被 `goSourceFiles` 排除，
不污染 AST 扫描），打印 `knownComposerMethod` 的真实返回值与真解包结果。
**负控**：干净树（`NONE`）上 4 枚候选名里只有 `panel.mode.request` 是 `true`，
`panel.review.allow` / `panel.review.grant` / `panel.approval.request` / `index.html` 全 `false`。

### §0.3 改之前的五发读数（本件的 5 枚钉 = 全绿的那些）

| 发 | 形状 | 本件红的测试 | 探针真值 | 判 |
|---|---|---|---|---|
| `M2M3` | 守卫 case 追加包级 `var acceptM2Route = "panel.review.allow"` 与拼接标签 `"panel.review." + acceptM3Tail` | **只有 `TestPlantedGrantWiringGoesRedInASnapshot`**，红句是陈旧性引信：`l2_grant_boundary_test.go:798: plant B has no anchor…`；`TestAnswered… / TestNoInbound… / TestGrantVocabulary… / TestGrantWire…` 四枚**全绿** | `knownComposerMethod("panel.review.allow") = true`、`("panel.review.grant") = true`、`ParseComposerRequest err=<nil> parsed={Method:panel.review.allow …}` | 复现验收 F-1（静默丢弃） |
| `M4` | 第二条链：`knownComposerMethod` 末尾 `return acceptM4Gate(m)`，新函数 `case "panel.review.allow": return true` | **零枚红，整包 rc 只剩 C21** | `knownComposerMethod("panel.review.allow") = true`、`ParseComposerRequest err=<nil>` | 复现验收 F-2 |
| `M6` | 新入站封套 `acceptM6Envelope{Cmd json:"cmd"; Outcome json:"outcome"}` + 真 `json.Unmarshal` 处理函数 | **零枚红** | `acceptM6Dispatch verdict = "grant"`（线上键真的绑进 Go） | 复现验收 F-4 |
| `M7` | `acceptM7Envelope{Method json:"method"; Outcome json:",omitempty"}` + 真 Unmarshal | **零枚红** | `acceptM7Dispatch verdict = "grant"` | 复现验收 F-5 的入站版 |
| `M8` | 同一形状放进正主 `ComposerRequest`（`Outcome string \`json:",omitempty"\``） | 反射那半红（`:624 ComposerRequest.Outcome binds "Outcome"`）、面 3 红（`:747`/`:752` 各 4 枚）；**AST 那半 0 findings** —— 判据是 `grep -c 'l2_grant_boundary_test.go:643' before-M8.txt` = **0**（那一行就是 `if len(findings) > 0 { t.Errorf(...) }` 的报点，它一声没出） | `parsed={… Text: Outcome:grant …}` | 复现验收 F-5（两台仪器分家） |
| `M13` | M4 的第二条链 **+** M6 的 `cmd`/`outcome` 封套（合体） | **零枚红**（本件 5 枚全绿，整包唯一红仍是 C21） | `knownComposerMethod("panel.review.allow") = true`、`ParseComposerRequest err=<nil>`、`acceptM13Dispatch verdict = "grant"` | **复现验收 §6 的决定性一发** |

原文读数：`D:\tmp\panel-l2-nail-fix-r2\out\before-{NONE,M2M3,M4,M6,M7,M8,M13}.txt`。
每发之后当场 `python mutate.py restore`，副本 `bridge.go` 回到 `d2cd6362…`，
仓库 `internal/panel/bridge.go` 每发之后复算仍是 `d2cd6362…`（**七发七次同值**，日志在上）。

⇒ 这一节的全部意义：下面 §2–§5 那几把新尺子如果哪一枚咬不动上面某一发，
它就不是钉，是镜子；改完之后同样的种子必须当场红（§7）。

---

## §1 进场锚点与地界

- 进场 HEAD `310816f`（`dev`），被验物 blob `5b618d14…`，生产码 blob `d2cd6362…`。
- 本程写的路径只有两枚：`internal/panel/l2_grant_boundary_test.go` 与本件。
  逐枚 commit、每次 commit 前 `git diff --cached --name-only` 现核（读数在 §9）。
- 未跑前端（`npm`/`vitest`/`tsc` 一律未跑），未跑 `go test ./...`（两路 agent 在飞）。

---

## §2 修法 (a)：不可解析的 case 标签**不许静默丢弃**（治 F-1 / M2M3）

### §2.1 改了什么

`routeNamesInFunc` 旧形状（`5b618d14` 的 `:406-427`）把 `routeNameOf` 返回 `false` 的标签
`if s, ok := …; ok { out = append(…) }` 一句吃掉，全包零条痕。现在：

- `routeLabelsInFunc(fn, consts, fset, rel) ([]string, []unresolvedLabel)` —— 解析不出来的标签
  带着 `file:line` 与 **`go/printer` 打回来的原文表达式**返回，不再消失；
- 新增 `guardReadabilityProblem(dir, pkg) string` 一枚纯函数（返回"这棵树的 answered 集为什么不可信"，
  没问题时返回 `""`），`requireReadableGuard(t, dir, pkg)` 只负责把它变成 `t.Fatalf`。
  拆成两半是**为了牙齿可测**：植物要断言"消息里点了那两枚表达式的名"，而 `t.Fatalf` 会把断言它的那枚
  子测试一起打死（Go 没有"吸收一发子测试失败"的写法）。
- `scanGrantBoundary` 与面 1 各自的三枚/两枚 fail-fast 收进同一枚 `requireReadableGuard`，
  所以**每一次真扫描都走这条路**，植物只是额外证明消息体存在。
- 新增植物 E（`TestPlantedGrantWiringGoesRedInASnapshot/E_…`）：把 M2M3 的形状种进快照
  （`case …, acceptE2Route, "panel.review." + acceptE3Tail:` + 一枚包级 `var` + 一枚 `const` 尾巴），
  断言 `len(pkg.dropped)==2`、`answered` 仍是 4 枚（**不许靠"不再读守卫"通过**）、
  且 problem 里必须同时出现 `acceptE2Route` 与 `"panel.review." + acceptE3Tail`。

### §2.2 现量（`-count=1 -v`，锚点＝本节末的 HEAD）

```
--- PASS: TestPlantedGrantWiringGoesRedInASnapshot/E_a_guard_route_reached_through_a_var_or_a_concatenation_is_named,_not_dropped
    l2_grant_boundary_test.go:976: loud as required, and the names are not invented:
          bridge.go:99: "panel.review." + acceptE3Tail
          bridge.go:99: acceptE2Route
```

⇒ 表达式是被 printer 打出来的原文，不是我抄的字符串（同一行两个标签都指到 `bridge.go:99`，
因为种子的 case 标签写在同一行——这正说明**只有行号不够，必须点名表达式**）。

### §2.3 (a) 能治什么、治不了什么（写清楚，别让它替谁说话）

- 治：`case <包级 var>:`、`case "a" + b:`、`case f():` —— 任何**写在正主守卫里**而 AST 读不懂的标签，
  从此当场 `t.Fatalf`（M2/M3 的形状）。
- 不治：标签在**另一枚函数**里（M4/M13 的形状）。那一发 `knownComposerMethod` 自己的 switch 一字未动，
  drop 列表是空的。⇒ 这是 §3 那把新尺的活，(a) 不是全部。

---

## §3 修法 (b)：神谕换成"包内每一枚像路由的字面量，都拿去问**真的**守卫"（治 F-2 / M4 / M13）

### §3.1 先独立推演验收的断言，再动手

验收 §7.3(b) 的原话是"**我推演它咬得住 M4 与 M13**（`"panel.review.allow"` 就是包内的字面量、守卫真答它）"，
并自己标了"这是未验证断言"。我独立推了三条，其中**一条它没说、而它说的两条我复算成立**：

| 待验 | 推演 | 现量（§3.3） |
|---|---|---|
| M4：字面量在包内？ | `acceptM4Gate` 的 `case "panel.review.allow":` 就在 `bridge.go`（非 `_test.go`，正是 AST 的分母） | 命中 `written at bridge.go:125` |
| M4：真守卫答它？ | `knownComposerMethod` 末尾 `return acceptM4Gate(m)` ⇒ 运行期 `true` | 探针 `= true`（§0.3）；红句直接引用真函数 |
| M13：同一发？ | M13 的路由半**就是** M4 的形状（第二条链 + 同一个名字），所以 (b) 应当只治它的**路由半** | 红句同形（`bridge.go:123`）；**封套半仍 0 findings**，那一半是 §4 的 (c) |
| 验收没说的一条 | 反过来还有一格：**守卫答了、名字却从没被整枚拼出来**（`"panel.review." + tail` 且 tail 在别处）⇒ pool 比 answered 小，整把尺退化成装饰 | `answeredOutsidePool` 那一判就是它，见 §3.2 第三条 |

⇒ **(b) 的判据不需要换**，但我把它从"一条"扩成了"三条互印"（见下），因为只加验收写的那一条会留下
"pool 小于 answered 时静默退化"这个新洞。**没有换方案，是加了一格。**

### §3.2 改了什么

- `pkg.literals []routeLiteral`：`collectRouteLiterals` 在**每枚文件解析后单独走一遍**，
  收全包（含别的函数、含 const/var 初始化式）里所有 `routeShapedName` 通过的字符串字面量，带 `file:line`。
  `routeShapedName`＝"至少两段、点分、只含字母数字下划线（段内不许以 `-` 开头）"，
  所以 struct tag（含 `:`）、mime（含 `/`）、带目录的路径、printf 动词全都进不来。
  **它只是筛子不是判据**：判据是真守卫，混进来一枚也无所谓。
- `routeNamePool(pkg)`：字面量 ∪ 包级 string 常量的值（守卫可能只经由一枚常量答一个名字），名字 → 出处列表。
- `poolJudgedByRealGuard(pool, answered)` 两向：
  1. `answeredElsewhere` —— 真守卫答了、**派生集没看见** ⇒ 红（M4/M13 那一发）；
  2. `answeredGrantDoor` —— 真守卫答了、**名字本身是审批形状** ⇒ 红（禁令本体）。
- `answeredOutsidePool` 第三向（防退化）：派生集里若有名字**不成枚出现在包里** ⇒ 红，
  外加面 1 里 `len(pool) < len(answered)` 直接 `t.Fatalf`。
  这一格是验收那版没有的：**神谕换成一枚更小的尺子时，最容易的就是没人发现它变小了。**
- 植物 F（`…/F_a_route_named_outside_the_guard_still_reaches_the_pool`）：把第二条链种进快照，
  断言 pool 收下了 `panel.review.allow` 并点到 `l2-plant-f.go:7`，同时断言
  `pkg.answered` **没有**它（证明这一发的确是旧枚举的盲区）。快照只解析不编译，
  所以运行期那半的真红在 §6 的真树变异里，这一枚钉的是"神谕看得见别的函数"这一格 ingredient。

### §3.3 现量：M4 与 M13 的真红

读数状态＝**(a) 已提交（`9d85789`）、(b) 在当时的测试树里**（blob 我在 commit 后当场复算，见本节末）。
副本 `headcopy` 与仓库只差这一枚被同步过的 `_test.go`，`bridge.go` 三处仍同值 `d2cd6362…`。

```
M4  --- FAIL: TestAnsweredPanelRoutesCarryNoApprovalDecision
    :817 knownComposerMethod answers "panel.review.allow" (written at bridge.go:125) but the guard's own
         case list, which is what this file enumerates, never named it: ... the answered set read here is
         INCOMPLETE, not clean (D33/F2, R20)
    :820 knownComposerMethod answers "panel.review.allow" (written at bridge.go:125), an approval decision
         addressed from the panel, and the name was found written in this package rather than guessed here
    +  --- FAIL: …/F_a_route_named_outside_the_guard_still_reaches_the_pool   （它的负控：真守卫现在答这个名字）
M13 --- FAIL: TestAnsweredPanelRoutesCarryNoApprovalDecision    （同两行，出处 bridge.go:123）
    仍 0 findings 的那一半：:878 "…AST scan of internal/panel: 0 findings"  ⇒ 那是 §4 (c) 的活
```

M2M3 那一发同时红了三枚（`TestAnswered… / TestNoInbound… / TestPlanted…`），红句是 (a) 的
`the inbound route guard … has 2 case label(s) this instrument cannot resolve to a string:
bridge.go:99: "panel.review." + acceptM3Tail / bridge.go:99: acceptM2Route`。
⇒ 验收 §0.3 里"4 枚生产面全绿、只有陈旧性引信叫"那一格，现在是**三枚叫、而且叫的是门**。

本节的锚点：`TestAnswered…` 里那两枚红句的行号（817/820）与 `:878`、`:839` 都按当时代码现读，
commit 后本件的行号会继续漂，引用时以 `l2_grant_boundary_test.go` 的**句子**为准，不要拿行号当句子的指纹。


---

## §4 修法 (c)：入站封套的判据换成"本包里作为 JSON decode 目的的那枚类型"（治 F-4 / M6 / M13 的封套半）

### §4.1 先复算验收给的两个数，再换尺

| 待验断言 | 我的命令 | 现量 |
|---|---|---|
| "全包只有 `bridge.go:79` 一处 json.Unmarshal" | `grep -rn "json.Unmarshal" internal/panel/*.go \| grep -v _test.go` | **1 处**，`bridge.go:79`，目的 `&r`（`var r ComposerRequest`）⇒ 成立 |
| 换尺之后**今天**的入站集合会不会变 | 临时探针 `zz_seedcomparison_probe_test.go`（跑完 `mv` 进 `removed/`）在同一棵树上并排算两把判据的闭包 | **逐枚相同**：`old-criterion inbound = [AttachmentPayload AttachmentRef ComposerRequest]`、`new-criterion inbound = [AttachmentPayload AttachmentRef ComposerRequest]`、`identical=true`（`t.Errorf` 那一支没触发），且 `decodes=1` 就是 `bridge.go:79 dst=&r type=ComposerRequest` ⇒ **这不是收紧也不是放宽，是把同一格外延换了个判据**，M6/M13 从此进射程（M7 见 §4.3，它不归 (c) 管） |
| 反例：为什么不能"查包里所有 struct" | `internal/panel/approval.go:54` `SessionOverrideBlocked bool json:"sessionOverrideBlocked"` 归一化含 `override`、`:58` `DecidedBy json:"decidedBy"` 含 `decide`，两枚都在 `ApprovalCardView`（**出站渲染视图**，不是任何 decode 的目的） | 派单警告成立，**而且我量到第二枚**（验收只点了 `decidedBy`）。⇒ 那一版今天会红两枚存量合法键，未采纳 |

⇒ 验收 §7.3(c) 的**判据选择成立**，但它给的理由（"M6/M7/M13 三枚植物都带自己的 json.Unmarshal ⇒ 三枚全在射程内"）
**只对了两枚**：M7 换了尺之后**仍然全绿**，因为它逃的不是种子判据、是 `jsonOr` 那枚空 tag 名（F-5）。
这一格我用现量把它钉在 §4.3，没有顺着验收的话写。

### §4.2 改了什么

- `boundaryPackage` 多三格：`pkgVars`（包级 var 的声明类型）、`decodes []decodeSite`、
  `inboundSeeds map[string]bool` + `decodeProblem string`。
- 解析改成**两遍**：第一遍逐文件收 struct/const/路由字面量/包级 var 并把 `*ast.File` 攒下；
  第二遍 `collectDecodeSites` 走每一枚函数体（含 `FuncLit`、含参数表），先收集局部件名→类型
  （`var x T`、`x := T{}`、`x := &T{}`、`x := new(T)`，同名两型记冲突），再把每一枚 decode 调用的
  目的表达式解析成类型。**为什么必须两遍**：目的可以是**另一个文件**里的包级 `var`，
  而"它到不到位得上同包 struct"要等 struct 集齐才能判（`classifyDecodes`）。
- 认哪些调用：`json.Unmarshal(x, &v)`；以及在**导了 `encoding/json` 的文件**里、receiver 文本含 `json`
  的 `.Decode(&v)`（`json.NewDecoder(r).Decode(&v)` 与先存成 `d := json.NewDecoder(r)` 两种写法都收）。
  `hex.Decode` / `base64.*.DecodeString` 进不来（名字与 receiver 都不满足）。今天全包 1 处，
  所以这一条对存量是零风险；写在这儿是为了**下一枚** decode 落地时自动进射程。
- **读不懂的目的**（`map[string]any` 这种能绑任意键的、别包的 struct、找不到声明的变量、冲突件名）
  ⇒ `instrumentBlindnessProblem` 点名 `file:line + 表达式 + 类型文本`并 `t.Fatalf`。
  原 `guardReadabilityProblem`/`requireReadableGuard` 因而改名
  `instrumentBlindnessProblem`/`requireReadableInstrument`：它现在判的是**两类失明**
  （守卫读不懂的标签、decode 目的枚举不出键），不只守卫。
- `inboundEnvelopes` 的种子从"绑 `method` 键"改成 `pkg.inboundSeeds`；**内嵌/嵌套闭包原样保留**
  （验收 §3.4 那发 M10 证明闭包承重，这一格不许动）。
- 植物 G：`{Cmd json:"cmd"; Outcome json:"outcome"}` + 自己的 `json.Unmarshal`。三条断言钉死它是谁：
  seed 认它（`inboundSeeds`）、它**不**绑 method 键（否则新旧两把判据分不开，植物就白种了）、
  扫描当场点名 `"outcome"`。

### §4.3 现量（副本树 = 仓库 test 文件同 blob；`bridge.go` 每发还原 `d2cd6362…`）

```
M6  --- FAIL: TestNoInboundEnvelopeCanBindAnApprovalVerdict
    :1191 the panel's inbound Go boundary has a grant-carrying face:
            bridge.go:124: inbound envelope acceptM6Envelope can bind the JSON key "outcome" -
            a verdict a page can set, ... : Outcome string `json:"outcome"`
    +  --- FAIL: TestPlantedGrantWiringGoesRedInASnapshot （它的前置引信先叫：
            :1333 "the snapshot of the real package is already dirty before any planting" —— 真树已脏，
            植物不能再用，这一发是**正确**的反应，不是误伤）
M13 两半同时红：:1133/:1136（(b) 那半）+ :1191 "bridge.go:133: inbound envelope acceptM13Envelope …"
M7  仍然**零枚红** —— 它的封套就是 acceptM7Envelope、就是 decode 目的、已经在射程里，
    但 AST 把 `json:",omitempty"` 算成空键名 ⇒ 词表判的是 ""。**这一发归 (c')**，见 §5。
M8  仍 0 findings（`grep -c 'has a grant-carrying face' afterC-M8.txt` = 0，改前 before-M8.txt 亦 0）
    ⇒ 同上，(c) 治不了 F-5。
```

⇒ 三发对照把话说死了：**(c) 治 M6/M13 的封套半，(c') 治 M7/M8，两枚各自承重、互相替不了。**

---

## §5 修法 (c′)：`jsonOr` 的空 tag 名（F-5），并把"两台仪器一条规矩"做成三向互印

### §5.1 缺陷本体与最小修

`json:",omitempty"` 的**名字部分是空**，`encoding/json` 于是回落到 Go 字段名（线上键 `Outcome`）。
旧 `jsonOr(goName, jsonName, tagged)` 判的是 `if tagged { return jsonName }` ⇒ 返回 `""`，
AST 那半就此**看不见这枚键**（验收 M8：反射红、AST 0 findings）。修成一格：

```go
if tagged && jsonName != "" { return jsonName }   // 空名字 = 回落到字段名
return goName
```

顺手把反射那半的一处**同族但方向相反**的不一致钉平（这条验收没量到，是我读它邻行读出来的）：
`bindableKeysOf` 旧写法是 `if tag == "-" { continue }`，比的是**整枚 tag**，
所以 `json:"-,"`／`json:"-,omitempty"` 在它眼里是"一枚名叫 `-` 的键"。
现在两半都按**名字段**判（`name == "-"` 即跳过），一台尺一个规矩。
⚠ 代价说清：`json:"-,"` 在 `encoding/json` 里其实是"键名就叫 `-`"，我现在两边都当它不存在——
但一名为 `-` 的键**装不进任何审批词**（`carriesGrantWord("-")` = false），
所以这个取舍对本禁令是无害的，写进 §10"没测什么"。

### §5.2 三向互印 `TestJSONKeyDerivationAgreesWithEncodingJSON`（新顶层测试）

F-5 的病不在"某个键算错"，在**两台自称同一条规矩的仪器可以分开**。所以补的不是一个用例，是一枚三向裁判：

- (i) `jsonOr` 的**规则表**：无 tag／有名 tag／`json:",omitempty"`／`json:"-"`，
  外加一枚 `f5Probe` struct（含 `Nameless string \`json:",omitempty"\``、`Untagged`、`Skipped json:"-"`、
  `Recursive json:"rec,-"`）逐键核反射那半；
- (ii) **AST 键表 == 反射键表**，逐枚 inbound 类型比（只比本层键，不做递归，两把尺比同一格东西才判得动）；
- (iii) **第三个裁判是 `encoding/json` 自己**：`DisallowUnknownFields()` 分得出
  "没有字段应答这个名字"与"有字段应答了、值不合适"，两者后者＝可绑。
  对每一枚列出的键问它"真能绑吗"，再对 **24 枚审批拼写**（`outcome/Outcome/allow/Allow/allowOnce/
  AllowOnce/approved/Approved/grant/Grant/verdict/Verdict/decision/Decision/decide/Decide/
  bypass/Bypass/override/Override/permit/Permit/authorize/Authorize`，两向拼写是因为回落规则让
  `Outcome` 同时接 `outcome`）问它"这也能绑？"。

外加一枚**防腐**断言：`pkg.inboundSeeds` 里任何一枚类型若不在 `inboundTypeRegistry()` 里 ⇒ `t.Fatalf`。
（反射枚举不了一枚包里有哪些类型，这份名单只能手写；手写的东西会烂，所以给它配了引信——
M7 那一发当场把这枚引信踩响了。）

### §5.3 现量

干净树（`-count=1`）：

```
ComposerRequest:   7 keys agreed by both instruments [attachments method path requestId source text to]
ModeRequest:       2 keys [correlationId to]      AttachmentPayload: 4 [dataBase64 declaredMime name sizeBytes]
AttachmentRef:     9 [artifact deduplicated id kind mime name reason sizeBytes stored]
每枚：N listed keys confirmed bindable, 24 verdict spellings confirmed not bindable
```

变异树（同一副本，`bridge.go` 每发还原 `d2cd6362…`）：

```
M7 --- FAIL: TestNoInboundEnvelopeCanBindAnApprovalVerdict
   :1211 the panel's inbound Go boundary has a grant-carrying face:
           bridge.go:125: inbound envelope acceptM7Envelope can bind the JSON key "Outcome" ...
   --- FAIL: TestJSONKeyDerivationAgreesWithEncodingJSON
   :1418 decode destination acceptM7Envelope has no reflection twin in inboundTypeRegistry …(引信，见 §5.2)
   +  --- FAIL: TestPlantedGrantWiringGoesRedInASnapshot （真树已脏的前置引信）
M8 --- FAIL: TestNoInboundEnvelopeCanBindAnApprovalVerdict   （AST 那半从此不再 0 findings）
   --- FAIL: TestJSONKeyDerivationAgreesWithEncodingJSON/iii_…
   :1466 ComposerRequest binds the wire key "outcome" … / binds the wire key "Outcome" …（2 条，两向拼写各一）
   --- FAIL: TestGrantWireShapesAreRefusedAtTheDoor/every_answered_route_drops_a_smuggled_verdict
```

⇒ 改前 M7 是**零枚红**（§0.3），改后 3 枚红；M8 改前反射红而 AST 静默，改后 AST 也红、
且**运行中的解码器自己出来作证**。植物 H 把同一形状钉在快照里（断言 AST 键表就是
`[Outcome method]`，改前是 `["", method]`）。

### §5.4 (c′) 与 (c) 的分工（不许混着领功）

派单说 (c) 的推论是"M6/M7/M13 都带自己的 `json.Unmarshal` ⇒ 换尺之后**三枚全在射程内**"。
现量是：**换尺之后三枚确实都进了射程，但只有 M6/M13 当场红**；
M7 在射程里仍然**静默**，因为射程只决定"看不看这枚 struct"，词表判的是"这枚 struct 的键算出来叫什么"。
⇒ §4 那半句我按读数改过（`cc6dc16`），(c′) 是独立承重的一枚，不是 (c) 的零头。

---

## §6 修法 (d)(e)：文件头补"今天什么都没接线"，`:604` 那处指向不存在测试的引用改掉

### §6.1 (d) 一句大白话，且每个字都可复算

派单要的是"nothing is wired today"这一句，验收 §1.1 判它是"下一位把 `Q-49` 划成 covered 的入口"。
我写成一段并把**每个断言都配一枚可复算命令**，不写"据称"：

```
grep -rn ParseComposerRequest --include=*.go .   ->  ./cmd/wisp/run.go(1,注释) ./internal/panel/bridge.go(4)
                                                      bridge_test.go(6) composer_handlers.go(1,注释)
                                                      l2_grant_boundary_test.go(6)
   非测试命中里除去定义本身(bridge.go:77)，其余四处全是注释  =>  生产调用者 = 0
grep -c webview go.mod                            ->  0        (主模块今天连 WebView2 依赖都没有)
sed -n '225,227p' cmd/wisp/run.go -> "Nothing calls it yet - the WebView2 \"event -> ParseComposerRequest\"
                                       hop does not exist in this tree (tickets 33/35)"
```

头段落同时改了一句**过去可能被读反**的措辞：绿"不等于 panel 拿不到批准权"，绿的意思是
"这道边界**被接线**时不可能不红"。并把"两枚补强之后仍然开着的洞"写进 `WHAT THIS FILE DOES *NOT* COVER`：
运行期拼出来的路由名（config／网络上来的）与词表之外的结论键名（`proceed`/`yes`/`ok`）——
这两格**不是我没做，是判据本身只看得见"写下来的东西"**，写明白比留白好。

### §6.2 (e) 指向不存在测试的那句

旧句（`5b618d14` 的 `:604`）：`… past the naming gate TestComposerMethodNamesMatchFrontend`。
那枚测试**定义 0 处**（派单已核；我进场先复算：见下）。新句不再把人往不存在的地方领，
而是把**这一行为什么是唯一一道闸**说清：

```
Go answers %q but no Method* constant declares it: a route written straight into the guard's case list.
Nothing else in this package would notice - TestFrontendComposerRequestsMatchTheEnvelope (bridge_test.go:131)
is the test that reads the frontend for these names, and it only ever iterates the four declared constants,
so a fifth route that never became a constant is invisible to it. This line is the only gate on that shape
```

- 指向的是**存在**的那枚：`grep -rn "func TestFrontendComposerRequestsMatchTheEnvelope" internal/panel/` = 1 处定义
  （`bridge_test.go:131`），且它的行为与我写的描述对得上（我读了全文：它 `for _, m := range []string{Method…}`
  四枚常量去 grep `frontend/src/lib/panel.ts`，**从不读守卫**）。
- "只迭代那四枚常量"这半句就是 M1b 那发的红句能立住的原因，不是我替它编的赞美。

### §6.3 顺手核到但**不修**的一处生产码失效引用

`internal/panel/bridge.go:33` 同一枚不存在的名字（票 92 留的，先存在）：

```
// Methods the composer route answers. Renaming one on either side goes red in
// TestComposerMethodNamesMatchFrontend, which greps the frontend for them.
```

派单明令我**别碰** `bridge.go:33`（生产码、别家地界）。现量：全仓 `grep -rn TestComposerMethodNamesMatchFrontend`
= **1 处（就是这一行）+ 我改掉的这处引用**，定义 **0 处**。
⇒ 上报为遗留缺陷：注释承诺了一枚不存在的门禁，实际干这活的是 `bridge_test.go:131`；
修它要动生产文件的注释，属票 92/接线切片卡那一程，**不在本程**。
