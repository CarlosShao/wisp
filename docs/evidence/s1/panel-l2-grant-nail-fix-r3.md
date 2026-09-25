# panel-l2-grant-nail-fix-r3 —— 实现者：把 F-R2-3 那条"借来的负控"提升成常驻测试，并把 F-R2-1 那条约束写进文件头

- 时刻 / 锚点：2026-09-25 12:1x +08 开工，进场 `git rev-parse --short HEAD` = **`1785a77`**。
  共享树，本程取三门读数的锚点是 **`cb60b94`**，写本件时 `dev` 已到 `b293784`。
  ⇒ 本件每个读数旁边写的是**它自己那一发的 sha / blob**，不是"当前 HEAD"。
- 工单来源：`docs/evidence/s1/panel-l2-grant-nail-accept-r2.md` **§7.2**（F-R2-3 的形状与它自己量过的修法）
  与 **§2.4 + §7.4 第 1 句**（F-R2-1：这批新判据给生产码偷偷加的那条约束）。总裁＝成立（附条件入账）。
- 身份：**本程是实现者**。§2 那套反向对照是实现者自证，**承重结论要另一程复算**，本程不给自己判成立。
- 地界：只写两枚路径 —— `internal/panel/l2_grant_boundary_test.go` 与本件。
  **未碰** `tools/d22scan/**`（另一程在飞）、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、
  `internal/risk/**`、`rules_gateway.go`、`allowlist.txt`、`thresholds.go`、任何 golden、
  `internal/panel/` 下别的文件、台账 `docs/reports/pending-and-issues.md`、工单票面、`Q-49` 的勾。
- 纪律：变异一律落在**仓外副本** `D:\tmp\panel-l2-nail-fix-r3\headcopy`，仓库树任何时刻都没进变异态
  （`run.sh` 每发行尾现跑 `git status --porcelain -- internal/panel` 与 `git hash-object internal/panel/bridge.go`）；
  临时件只建不删（未 `rm`/`rmdir`）；未 `add -A`/`.`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、未 push。

---

## §0 简报给的六枚前提：五枚复算对、一枚对不上（对不上的那枚没让我改走法，但改了口径）

| 说法 | 本程复算 | 判 |
|---|---|---|
| 被验物 1964 行 | `wc -l` = **1964** | 对 |
| tripwire 在 `:1778-1780`（简报）／`:1779-1781`（验收件 §7.2） | `grep -n 'this is no longer a clean tree'` → **1779**（`if` 在 1778、`}` 在 1780）⇒ 简报对，验收件整段 +1 | 对简报成立；**本程按字符串定位，没按行号动手** |
| `grantRouteWords` 在 `:133`、11 枚词 | `var grantRouteWords = []string{` 在 **133**，枚数 **11** | 对 |
| `bridge.go` 应仍是 `d2cd6362…` | 进场 / 每发变异之后 / 交件时共 12 次 `git hash-object`，**全部 = `d2cd6362ecc6941a0cee58073a2bf8ab6669a590`** | 对 |
| 被验物＝验收件锚定那一版（没人抢先动过） | 工作树 blob = `git rev-parse HEAD:…` = **`3b7a2cb1…`** = 验收件 §0.1 那枚 | 对 |
| "C21 那两枚红"是包内唯一的先存在红 | 干净树现量**三枚**红名（§1.3） | **对不上** |
| "M14 摘掉 tripwire 后只剩先存在的 C21 红" | 形状成立、**分母变了**：今天摘掉后剩的是那**三枚**（`out/T-C-nostanding-M14.txt`） | 按本程读数记 |

§1.3 那一枚**不是本程造的**，也已被编排者记过账（见那条末尾的 `A237②`），本程只是复算到同一处。

### §0.1 进场读数（`1785a77`，全部现跑）

```
git rev-parse --short HEAD                                -> 1785a77
git status --porcelain                                    -> 只有 design/ 的 16 枚删除 + 4 枚 design/doubao 未跟踪
                                                             （别家在飞；本程未碰、未恢复、未 commit）
wc -l internal/panel/l2_grant_boundary_test.go            -> 1964
git hash-object internal/panel/l2_grant_boundary_test.go  -> 3b7a2cb1d4105df1f279f70c605ccf6a71f7868a
git rev-parse HEAD:internal/panel/l2_grant_boundary_test.go -> 3b7a2cb1d4105df1f279f70c605ccf6a71f7868a
git hash-object internal/panel/bridge.go                  -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
grep -n 'this is no longer a clean tree'                  -> 1779
```

---

## §1 三门：改前 / 改后（四数之外比名册差集）

### §1.1 台子（副本与被验物逐字节同一）

```
cd "D:\work\workspace\projects plans\Wisp" && tar --exclude=./.git --exclude=./third_party \
  --exclude=./build --exclude=./node_modules --exclude=./dist -cf - . \
  | (cd D:/tmp/panel-l2-nail-fix-r3/headcopy && tar -xf -)     -> real 0m6.373s, 311M

headcopy  internal/panel/bridge.go                   d2cd6362ecc6941a0cee58073a2bf8ab6669a590
工作树    internal/panel/bridge.go                   d2cd6362ecc6941a0cee58073a2bf8ab6669a590
交付钉    internal/panel/l2_grant_boundary_test.go   9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7（2081 行，见 §6）
```

### §1.2 `go test ./internal/panel/ -count=2 -v`（锚点 `cb60b94`）

尺子＝验收件 §0.2 声明的那把（`=== RUN` 第 3 字段 → 去 `[N]`/`#NN` → `sort -u`，**含子测试**）。

```
改前（钉 = 3b7a2cb1，out/before-count2.txt）
  RUN=196  TOPPASS=98  SUBPASS=92  FAIL=6  SKIP=0  ^panic:=0  distinct=98
  红名各 ×2：TestC21DesignTokensFourWayAgree
             TestPlantedRendererDoorShapesGoRed
             TestTheRendererHoldsExactlyOneDoorToTheHost
改后（钉 = 交付版，out/final3-count2.txt）
  RUN=198  TOPPASS=100  SUBPASS=92  FAIL=6  SKIP=0  ^panic:=0  distinct=99
  红名各 ×2：与改前**逐枚同**（那三枚）
名册差集 diff out/before-count2.txt.roster out/final3-count2.txt.roster
  76a77
  > TestRealGuardRefusesEveryAssemblableApprovalRouteName
新增用例内的那一行读数（同发，×2 同值）
  l2_grant_boundary_test.go:1317: behavioural sweep of the running guard:
                                  asked=528 assembled route names, answeredByRealGuard=0, hits=[]
```

⇒ 名册 **`+1 / −0`**：没有哪枚原有用例消失。`FAIL=6` 改前后同值 ⇒ 三枚先存在的红**未修、未跳、未放宽**。
⚠ 与验收件 §6.3 那发（锚点 `0b95e9e`）的差：`TOPPASS 102 -> 98`、`FAIL 2 -> 6`，差的正是 §1.3 那两枚。

### §1.3 前提对不上的一枚：今天包内有**三枚**先存在的红，不是"C21 两枚"

现量归因（不是推理）：

```
git show HEAD:frontend/src/lib/panel.ts | grep -n 'view.request'
  -> 252:   sendRequest("panel.view.request", { source: "panel-view", to });
git status --porcelain -- frontend        -> （空）⇒ 这一枚是**已提交**的，不是脏树
git log --oneline -1 -S'panel.view.request' -- frontend/src/lib/panel.ts
  -> d61281c feat(前端·换屏层 Q1=甲): 左边一竖条图标换页上线，顶部四枚文字标签页(乙)退役
红句原文：composer_test.go:522 "the renderer names a route the Go side does not answer:
          src/lib/panel.ts:252 … names \"panel.view.request\""
```

⇒ 前端那一程在 `d61281c` 加了第六枚出站方法名、Go 侧白名单没跟，于是
`TestTheRendererHoldsExactlyOneDoorToTheHost` 与 `TestPlantedRendererDoorShapesGoRed` 从今天起红。
两枚都定义在 `internal/panel/composer_test.go`，**不在本程地界**；本程未修、未跳、未放宽，
只在改前/改后各带一次以证明既不是本程造的、也没被本程放绿。
**编排者已记过同一处账**：`docs/reports/pending-and-issues.md:6124`（`A237②`）逐名点过这两枚用例、
同一红因行 `:522`/`:524`；`Q-50`（同文件 `:1084`）是它的前置判断，`4712ea6` 已把"今天完全惰性"那句收回。
⇒ 本程那句"前提对不上"**只针对简报的字面**（简报说"C21 那两枚"），**不新造账**。

另：`TestC21DesignTokensFourWayAgree` 今天的红句是
`read design/assets/tokens.css: … cannot find the path` —— 那是**别家在飞的 `design/` 删除**
（`git status` 里那 16 枚 `D`；`tokens.css` 在 HEAD 里跟踪着）⇒ 这一枚"先存在的红"今天的**机制**是脏树，
不是验收件当年记的那枚内容缺陷。本程一个字未动。**编排者若按"C21 两枚红"核账会数不上，特此报备。**

### §1.4 另两门 + 契约轴（锚点 `cb60b94`，交付态）

```
gofmt -l internal/panel/     -> 空，rc=0
go vet  ./internal/panel/    -> 空，rc=0
sh scripts/d22scan.sh        -> rc=0
  d22scan: examined 225 production Go files under internal/ and cmd/
  bans #1-5 internal/=203  cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
  ban #8 design/=32  frontend/=46  internal/=407  cmd/=39
  d22scan: clean - no D22 ban violations
  （另：脚本自带的 tools/d22scan 测试跑完 = runtests.sh: OK … PASS=30 FAIL=0 SKIP=0）
```

⚠ 抄验收件 §6.3 的口径：bans #1–5 **不含 `_test.go`**，本程改的恰好是一枚 `_test.go`
⇒ 这道门对本次改动只有 ban #8 那一枚有分母。`tools/d22scan/**` 一字节未碰，也**没有**在根目录
`go vet ./tools/d22scan/`（它是独立模块）。

契约轴（同一把尺，零容忍面 + 一枚正控）：

```
git diff --numstat HEAD -- internal/panel/
  -> 130  14  internal/panel/l2_grant_boundary_test.go        （只此一枚路径，生产码零字节）
git diff --numstat HEAD -- internal/observe/thresholds.go \
     internal/risk/rules_gateway.go tools/d22scan/allowlist.txt \
     internal/risk/ docs/PLAN.md docs/specs/ '*.sse' internal/panel/bridge.go
  -> （空）＝零命中
正控（同尺、同区间、指自己那枚文件）：上面那行 130/14 就是它
```

先证明被点的三面**存在**才敢说"零命中"（同尺现跑）：

```
git ls-files | grep -iE 'thresholds\.go$|allowlist\.txt$|rules_gateway\.go$'
  -> internal/observe/thresholds.go   internal/risk/rules_gateway.go   tools/d22scan/allowlist.txt
git ls-files | grep -c '\.sse$'   -> 52 枚 golden（一字节未动：上面那发 numstat 空输出里根本没有它们）
```

### §1.5 断言尺（`scripts/count.py`，A = 进场那一版 `3b7a2cb1`，B = 交付版）

```
t.Fatalf(          A=31   B=33   delta=+2
t.Fatal(           A=5    B=5    delta=+0
t.Errorf(          A=38   B=40   delta=+2
t.Error(           A=2    B=1    delta=-1     ← 唯一的减少，就是被提升走的那枚借来的负控
t.Logf(            A=14   B=15   delta=+1
t.Skip             A=0    B=0    delta=+0
func Test          A=6    B=7    delta=+1
assertions(total)  A=76   B=79   delta=+3
行数                1964 -> 2081
```

**commit 1 删掉的 14 行逐枚看过**（commit 2 另删 2 行注释、加 3 行，是同一处头段措辞）（`out/commit-81ad6fd.diff` 与 `out/commit-121006d.diff` 里存着两份 `-` 行原文）：
**11 行是注释与一行 `t.Logf` 的措辞改写**（facet 1 那颗列表项 1 行 + 植物 F 头段 9 行 + 那行 `t.Logf` 1 行），
另外 **3 行**就是那一组 `if knownComposerMethod("panel.review.allow") { t.Error(…) }` ——
三行一起搬进 §2 那枚常驻测试。11 + 3 = 14，**没有第 15 枚删除**。
⇒ "只增不减"按**判据**读：没有一枚判据失去承接（`t.Error` 少的那一枚，其谓词由新测试的 528 问 +
`:1310` 那枚专门点名它的 `t.Errorf` 承接），且总数 `76 -> 79`。**这一处请验收程按判据轴复判，别只按枚数轴。**

---

## §2 F-R2-3：借来的负控 → 常驻顶层测试

### §2.1 交付的形状（文件 `internal/panel/l2_grant_boundary_test.go`）

| 位置（交付版行号） | 是什么 |
|---|---|
| `:102` | 文件头新增一段：**F-R2-1** 那条约束（见 §3） |
| `:46`（`WHAT THIS FILE PINS` facet 1 那颗） | 补一句：facet 1 有第二枚**常驻**的、不读源码的半边 |
| `:94`（`DOES NOT COVER` 那两段之后） | 补一句限定：运行期拼名字那枚洞**已被行为扫收一大截**，剩下的网格之外部分仍开 |
| `:1239` | `var grantRoutePrefixes` —— **8 枚命名空间**（不是路由名） |
| `:1247` | `var grantRouteSuffixes = []string{"", ".request", ".now"}` |
| `:1255` | `knownComposerRefusesAssembledGrantNames(t)` —— 运行期拼接 + 问真守卫，返回 `(asked, hits)`；内含两枚 `t.Fatalf`：`:1258` 词表为空、`:1266` "守卫连自己声明的四枚路由都不答"（死守卫反空转控）。两枚都带 `t.Helper()` ⇒ **报出来的行号是调用点 `:1308`**，不是它们自己 |
| `:1307` | **`TestRealGuardRefusesEveryAssemblableApprovalRouteName`** —— 顶层常驻测试 |
| `:1310` | `t.Errorf`：点名被回答的那一枚候选名 |
| `:1315` | `t.Errorf`：`panel.review.allow` 被回答时补一句"植物 F 与 answeredElsewhere 现在描述两棵树"（**共红诊断**，不能单独触发） |
| `:1317` | `t.Logf`：`asked=… answeredByRealGuard=… hits=…` |

**它问的名字一枚都不是硬编码**：`8 前缀 × 11 词 × 3 后缀 × 2 单复数 = 528`，
全部由 `p + "." + form + suffix` 运行期拼出，逐枚交给**真的** `knownComposerMethod`。
分母与网格形状照抄验收件 §7.2 那枚探针（`prefixes / grantRouteWords / {"", ".request", ".now"} / {w, w+"s"}`），
**故意的**：这样两程读数能逐字对账；乘积写死在 `t.Logf` 里，人事后可复算。

### §2.2 零误伤 + 抓得住（**仓外副本**现跑，`-count=1 -v`，钉 = 交付版）

| 行 | 台子 | 除三枚先存在的红之外 | 本测试的读数 |
|---|---|---|---|
| T-A | 交付态 + 干净守卫 | **无** | `asked=528 answeredByRealGuard=0 hits=[]` → PASS（**零误伤**） |
| T-B | 交付态 + **M14**（运行期拼 `panel.review.allow`） | `TestRealGuardRefuses…` **红** | `hits=[panel.review.allow]`，红句点名该名字 |
| T-D | 交付态 + **M16**（运行期拼 `panel.review.ratify`） | `TestRealGuardRefuses…` **红** | `hits=[panel.review.ratify]` |
| T-G | 交付态 + **M17 死守卫控**（守卫连那四枚声明路由都不答） | 本测试**红**，红因是 `:1266` 那枚反空转 `t.Fatalf`（`t.Helper()` ⇒ 读数报在调用点 `:1308`；原文 `out/T-G-delivered-M17deadguard.txt` 第 195 行）；同发还红 `TestAnsweredPanelRoutes…`、`TestComposerEnvelopeAcceptsItsFourRequests`、`TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests`（＋它那两枚 subtest） | `the running guard refuses its own declared route "panel.mode.request": the sweep below would report a clean boundary because nothing is answered any more, which is not the same fact` |
| 加分 | 前缀集 8 → 14（`scripts/widen.py on`，**未交付**） | 无 | `asked=924 answeredByRealGuard=0 hits=[]` → PASS（网格可加宽且不引入误伤） |

**与验收件自量对账**：它记 clean `528 / 0`、M16 树 `528 / hits=[panel.review.ratify]` ⇒
本程 T-A / T-D **逐字同值**（同一构造、同一分母）。**没有对不上的读数。**

### §2.3 反向对照：单点摘掉本程这枚常驻测试

摘法 = `scripts/plant.py nail standing`：从该测试的 doc 注释起、到它自己的 `t.Logf` + 闭括号止一次切掉
（两枚锚点各断言"恰出现一次"、切完断言 `func TestRealGuardRefuses…(` 已不在文件里、零变化即拒绝报告），
产出 `delta_lines=-32`。**`knownComposerRefusesAssembledGrantNames` 那枚 helper 故意留在原地**——
"摘掉测试"不许顺手把它的工具也摘走，而一枚没人叫的 helper 不能替任何人承重。

| 行 | 台子 | 红名（全部，`-count=1`） | 读数 |
|---|---|---|---|
| T-C | 摘常驻测试 + **M14** | 只剩 `TestC21DesignTokensFourWayAgree`、`TestPlantedRendererDoorShapesGoRed`、`TestTheRendererHoldsExactlyOneDoorToTheHost` | **没有任何别的测试顶着** ⇒ M14 回到"今天没人拦"，与本程新测试互为反证 |
| T-E | 摘常驻测试 + **M16** | 同上三枚 | 同上 |
| T-F | 摘常驻测试 + 干净守卫 | 同上三枚；`distinct` **99 → 98**，`A 减 F` 恰差 `TestRealGuardRefusesEveryAssemblableApprovalRouteName`，`F 减 A` **空** | 摘除本身不改别的任何颜色 ⇒ 真是**单点** |
| 补 | 摘常驻测试 + M17 死守卫控 | `TestAnsweredPanelRoutes…`、`TestComposerEnvelope…` 等**照样红** | ⇒ 反空转那半**不是**本程独有的承重；诚实记下：**承重只属 M14/M16 那一族** |

⇒ 按本仓那句定义（摘掉它，存在一发变异从此打不红 ⇒ 承重）：**这枚是钉**。
⚠ 验收件 §7.2 那句"只剩先存在的 C21 红"在今天要读作"只剩先存在的**三枚**红"（§1.3），**形状结论不变**。

### §2.4 植物 F 那三行：本程选**删**，理由是"留下＝两枚判据互相抵账"

选了哪支：**删掉**植物 F subtest 里的
`if knownComposerMethod("panel.review.allow") { t.Error(…) }` 三行，禁令整枚搬进 §2.1 的常驻测试；
并把该 subtest 的头段注释与 `t.Logf` 改成"这里只主张 pool 那半边，问运行守卫那半住在
`TestRealGuardRefusesEveryAssemblableApprovalRouteName`"（`:1864` 起那段、以及 `:1898` 那行日志）。

**为什么不是留下**：两枚判据问的是同一枚谓词（"真守卫不许答 `panel.review.allow`"）。
并存时把本程的常驻测试单点摘掉，M14 会**照样被植物 F 咬红**，"新测试承不承重"这一格就永远量不出来。
这一支本程**量了**不是推的（`scripts/plantf.py back` 把那三行放回副本，再摘常驻测试，跑 M14）：

```
row G（= scripts/runG.sh：交付版钉 → nail standing −32 行 → plantf back +3 行 → bridge M14 +352B/9 行）
      原文 out/G-keepboth-nostanding-M14.txt
红名：TestPlantedGrantWiringGoesRedInASnapshot
      TestPlantedGrantWiringGoesRedInASnapshot/F_a_route_named_outside_the_guard_still_reaches_the_pool
      + 那三枚先存在的
红句：l2_grant_boundary_test.go:1867: the real knownComposerMethod answers panel.review.allow -
      this is no longer a clean tree, and the negative control below is void
```

⇒ 并存那一发**仍红**＝两枚互相顶账，判据只能写"不承重"。删掉之后 T-C 才给得出真读数。

⚠ **本程自己的一发台件失败读数，不藏**：这一发我先命名成 `row U`（`out/U-keepboth-nostanding-M14.txt`），
跑法是 `plantf back` 之后再 `nail standing`。而 `plant.py nail standing` 是**从 `backup/` 重写整枚文件**的，
它把上一步塞回植物 F 的那三行又抹掉了 ⇒ 那一发实测**等同 T-E**（现复算：
`grep -c 'no longer a clean tree' out/U-keepboth-nostanding-M14.txt` = **0**，红名里没有植物 F）。
发现方式＝当场那句 `plantf.py forth` 报 `anchor count = 0, expected 1`（找不到"旧形状"⇒ 文件不是我以为的那版）。
`row G` 是把顺序摆正（先摘测试、再放回那三行）之后**重跑**的那发，上面引的凭据就是它，不是 U。

**"摘掉本程那枚常驻测试之后有什么变红或变不响"**（按交付态答）：
- **变不响**：M14、M16 两形（T-C / T-E）——今天唯一目击者就是被摘那枚。
- **变红**：没有。摘除是单点（T-F：`distinct 99→98`、`FAIL` 枚数不变、别的测试一枚都不改颜色）。
- **植物 F 那枚 subtest**：摘掉常驻测试后它在干净树与 M14 树下**都还是绿的**
  （`out/T-B-delivered-M14.txt` 里 `--- PASS: …/F_a_route_named_outside_the_guard_still_reaches_the_pool`，交付版钉那一发）
  ⇒ 它的语义已收窄成"快照里另一枚函数写下的名字进不进 pool"，主张的那半没变弱，
  也**不再借**运行时那半；`pkg.answered["panel.review.allow"]` 那枚 `t.Errorf` 原样留着（它问的是快照 AST，不是真守卫）。

### §2.5 一处**没做**的事，写明免得被读成"做过了"

副本里那发 `widen.py`（8 → 14 枚前缀、`asked=924`、`hits=0`）**没有交付**：交付集就是 8 枚，
与验收件 §7.2 那枚探针同分母。理由是"两程账对得上"比"网格宽一点"有用；
要拓是改 `grantRoutePrefixes` 一行，`t.Logf` 的 `asked=` 会自己把新分母报出来。

---

## §3 F-R2-1：那条对生产码的新约束现在写在文件头 `:102` 起

**先测后写**：写那段之前先在副本现量它到底踩响几枚测试（`scripts/mdec.py`，只动副本的 `bridge.go`）。

```
row V（原文 out/V-delivered-MDEC.txt）：在副本的 bridge.go 里加一枚**合法**的
`var m map[string]any` + json.Unmarshal（即验收件 §2.4 的 MDEC 形；副本跑完立刻 `mdec.py restore`，
复算 `bridge.go = d2cd6362…`）
本件红名（4 枚）：TestAnsweredPanelRoutesCarryNoApprovalDecision
                  TestNoInboundEnvelopeCanBindAnApprovalVerdict
                  TestJSONKeyDerivationAgreesWithEncodingJSON
                  TestPlantedGrantWiringGoesRedInASnapshot
红句（四处同一句，交付版行号）：:1162 / :1351 / :1554 / :1740
  "a JSON decode destination in 1 place(s) cannot be enumerated: bridge.go:111:
   decodes into &m (type \"map[string]any\"), which is not a same-package struct …
   Judge it in a test that can see that type, or route the bytes through a same-package struct -
   do not let this file report the boundary clean"
其余红名 = §1.3 那三枚先存在的（本程未动）；本程那枚常驻测试在此发**仍绿**（它不读 decode）
```

⇒ **与验收件的一处不符**：§7.4 第 1 句写"三枚 ban 测试 `t.Fatalf`"，本程现量是**四枚**
（facet 4 的快照那枚走的是同一道 `requireReadableInstrument`，它也算响）。
文件头那一段因此按**四枚**写，并点名是 facet 1 / facet 2 / 三向对照 / facet 4 快照。
这段是**文档**：不改判据、不放宽任何断言，方向仍是 fail-closed，
出路（落到同包 struct，或登记进 `inboundTypeRegistry` 再在看得见那类型的测试里判）逐字抄在头里。

---

## §4 本程**没有**测什么（不靠沉默读成通过）

- **没跑 `go test ./...`**（派单明令 + 另两程在飞）：本件所有"全包"字样都指 **`internal/panel` 一包**。
- **没修、也没验**§1.3 那两枚 `composer_test.go` 红（`panel.view.request` 缺 Go 侧登记；账在 `A237②`）：
  不在本程地界。本程只在改前/改后各带一次。
- **没测**"路由名从 config / 网络来"那一族：M14/M16 是**运行期拼接**，不是外部输入；
  **没测**词表外的结论键（M18 那一半）⇒ 那两形仍归文件头 `DOES NOT COVER` 与验收件 §7.4 第 4 句（F-R2-4）。
  ⇒ **本程没有把 F-R2-4 闭上，也没有主张闭上。**
- **没测**交付网格之外的名字：`wisp.review.ok`、`composer.allow` 之类**不在**那 528 枚里；
  924 枚那发只在副本跑过一次，未交付。
- **没做**并发 / 线程 / 进程级测量，**没跑**前端四道门（eslint/vitest/build/lint），**没扩** `emojiRe`
  或任何 d22scan scope 一字节，**没在 Linux 容器里量过任何东西**。
- **没读** `design/**` 脏树内容作任何凭据（只记下"它有 16 枚未提交删除"这一枚事实，用来解释 C21 今天的红句）。
- **没写**台账、票面、`Q-49` 的勾；**没给自己判成立**（§2 是实现者自证，需要另一程复算）。
- 临时件只建不删：`D:\tmp\panel-l2-nail-fix-r3\` 下 `backup/` 2 枚、`out/` 50 枚读数与名册、`scripts/` 8 枚（`plant.py` `plantf.py` `mdec.py` `widen.py` `run.sh` `runG.sh` `roster.sh` `count.py`），
  **未 `rm`**；仓库树内未建 worktree、未 checkout。

---

## §5 给对抗验收的：要复算的是这六句

1. **是钉不是镜子**：单点摘它（`plant.py nail standing`，helper 留在原地）+ M14 ⇒
   只剩三枚先存在的红（`out/T-C-nostanding-M14.txt`）；M16 同（`out/T-E-nostanding-M16.txt`）。
2. **零误伤**：`asked=528 / hits=0`，分母可复算 = 8×11×3×2，前缀与后缀集就写在 `:1239` / `:1247`。
3. **抓得住 M16**：`out/T-D-delivered-M16.txt`，红句点名 `panel.review.ratify`。
4. **名册**：`+1 / −0`（§1.2 那三行 diff）；`t.Skip` 仍 0；断言尺 `76 -> 79`、`t.Error 2 -> 1`
   ——少的那一枚就是被提升走的负控；`§1.5` 已把它按"判据轴 vs 枚数轴"写清，请**按判据轴复判**。
5. **三门与契约轴**：§1.4 原文；`bridge.go` 十二次同值 `d2cd6362…`。
6. **本程量出与上游不同的两枚数**（按本程读数写、未回改成它的口径）：
   先存在的红是**三枚**不是两枚（§1.3，账已在 `A237②`）；MDEC 踩响**四枚**不是三枚（§3）。

---

## §6 交件态

| 步 | 内容 |
|---|---|
| 进场锚点 | `1785a77`（三门读数取在 `cb60b94`；写本件时 HEAD 已到 `b293784`，共享树别家在飞） |
| 被验物（改前） | `internal/panel/l2_grant_boundary_test.go` blob `3b7a2cb1d4105df1f279f70c605ccf6a71f7868a`（1964 行）——与验收件 §0.1 那枚同值 |
| 本程交付（改后） | 同文件 blob **`9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7`**（2081 行），两枚 commit 累加 **`131/14`**（`git diff --numstat 1785a77 121006d -- 该文件`）；
自进场锚点起 `git log --oneline 1785a77..HEAD -- 该文件` **只有本程那两枚** ⇒ 没有别人在这格上动过手 |
| commit 1 | **`81ad6fd`** `test(panel l2门钉 r3): 把 F-R2-3 那三行…`，`git show --name-only` **只有 1 枚路径** `internal/panel/l2_grant_boundary_test.go`，`130/14`；那一发里的 blob = `9c06f3f6…` |
| commit 2 | **`121006d`** `test(panel l2门钉 r3 更正): 头段那句 "三枚 ban 测试 t.Fatalf" 按 §3 现量改成四枚`——只改注释，`git show --name-only` 同样只有那一枚路径，`3/2`；交付 blob = **`9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7`**。**没有 `--amend`**，按本仓规矩追加新 commit（`81ad6fd` 之后别家提了 `4712ea6`/`b293784`，本程未动它们） |
| commit 3 | 本件自己那一枚，路径只有 `docs/evidence/s1/panel-l2-grant-nail-fix-r3.md`。⚠ **本程不给它预先写 sha**——那正是验收件 §0.4 记过的那枚"文档自指悬空"缺陷的形状；要账就在交件后现跑 `git log --oneline -1 -- docs/evidence/s1/panel-l2-grant-nail-fix-r3.md` |
| 生产码 | `internal/panel/bridge.go` `d2cd6362…` 全程同值；每一发变异之后当场复算 `git status --porcelain -- internal/panel`（只有本件那一枚 ` M`，commit 3 之后为空） |
| 未做 | 未 push、未 `add -A`/`.`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、未 `rm`、未写台账/票面/`Q-49` |

---

## §7 交件后追加的两条（本仓规矩：已提交的行不改写，要更正就往下追加）

1. **§6 那行"commit 3 不预先写 sha"现在量到了**：本件那一枚 commit = **`062d389`**
   （`git show --name-only 062d389` → 只有 `docs/evidence/s1/panel-l2-grant-nail-fix-r3.md` 一枚路径）。
   本程在 dev 上的三枚 = `81ad6fd` → `121006d` → `062d389`，`§6` 表里那句"到 `b293784`"是**写本件当时**的读数，
   现在已经往后走了好几枚别家的 commit（共享树，`git log --oneline -1` 现跑＝`062d389`，其父＝`2dbb6af` 别家的票 142 验收）。
2. **committed 态把三门又跑了一遍**（不是引用 §1 那发，是又一发）：
   `gofmt -l internal/panel/` 空 rc=0、`go vet ./internal/panel/` 空 rc=0、
   `go test ./internal/panel/ -count=2 -v` = `RUN=198 TOPPASS=100 SUBPASS=92 FAIL=6 SKIP=0 ^panic:=0 distinct=99`、
   红名仍是 §1.3 那三枚（各 ×2）、新测试那行仍是 `asked=528 assembled route names, answeredByRealGuard=0, hits=[]`
   （原文 `out/committed-count2.txt`）；`git hash-object internal/panel/bridge.go` = `d2cd6362ecc6941a0cee58073a2bf8ab6669a590`、
   交付钉 = `9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7`、`git status --porcelain -- internal/panel` **空**。

---

## §8 §1.4 那发 `d22scan.sh` 的时效（又一处"行号/分母是读数不是常量"）

`§1.4` 记的那一发跑在**头段改准之前**的那版钉上（blob `9c06f3f6`，2080 行）。在交付版（`9fd6defb`，2081 行）上
本程又跑了一遍（原文 `out/d22scan-committed.txt`）：**rc=0**，逐枚分母与 §1.4 的差别只有一处，
而那一处**不是本程造的**——是别家在 `cmd/` 里多提了一枚文件：

```
bans #1-5 internal/=203  cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
ban #8 design/=32  frontend/=46  internal/=407  **cmd/=40**（§1.4 那发是 39）
d22scan: clean - no D22 ban violations
```

⇒ 引本件任何一枚 d22scan 分母时**连锚点一起引**（台账里已有这条规矩：`引用"某区间几枚红"必须连口径一起引并带锚点 sha`）。
本程未碰 `tools/d22scan/**` 一字节：`git log --oneline 1785a77..HEAD -- tools/d22scan/` 的命中全属于别家。
