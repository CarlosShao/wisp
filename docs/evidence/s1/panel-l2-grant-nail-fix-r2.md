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
