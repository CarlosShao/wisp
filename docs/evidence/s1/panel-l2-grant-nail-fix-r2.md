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

