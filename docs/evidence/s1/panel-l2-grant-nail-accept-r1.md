# panel-l2-grant-nail-accept-r1 —— 非实现者验收：`d88c356` 那枚 Go 侧 L2「允许」门钉

- 时刻 / 锚点：2026-09-25 10:0x +08 开工，验收时 `dev` @ `43a6043`（HEAD 在我脚下继续动，每节末尾刷新）。
- 被验物：`internal/panel/l2_grant_boundary_test.go`（`d88c356` 落地 945 行，`660ffa1` 就地改 4/2 行 → 现 947 行）。
  实现者自己的记录＝`docs/evidence/s1/panel-l2-grant-nail-r1.md`（§0–§9）。
- 身份：**本程不是实现者**，未写一字节生产码、未改任何既有文件。可写路径只有本件。
- 地界：未碰 `internal/observe/**`、`tools/d22scan/**`（另两路 agent 在跑，本程只对它们做了只读检查）、
  `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`。
- 方法：实现件里的每一条读数都当**待复测的断言**，不当证据。变异一律先备份到仓外、`cp` 还原、
  `git hash-object` 前后同值；没有用 `checkout`/`reset`/`stash`，没有 push，没有 `add -A`。

---

## §0 闸门复跑（真跑，非转述）+ 口径裁决

### §0.1 后态（真树、HEAD `dd2b528`，只读）

```
go test ./internal/panel/ -count=2 -v   ->  rc=1
RUN=180  PASS=100  FAIL=2  SKIP=0  ^panic:=0  distinct=90
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)     x2
```

### §0.2 前态（同一棵工作树的**仓库外副本**，`tar` 排除 `.git` 后把
`internal/panel/l2_grant_boundary_test.go` 一枚文件移到我自己的 scratch 里，其余逐字节同树）

```
go test ./internal/panel/ -count=2 -v   ->  rc=1
RUN=152  PASS=90   FAIL=2  SKIP=0  ^panic:=0  distinct=76
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)     x2
```

⇒ **实现件 §2/§5 的四数、panic、distinct、两枚红，逐枚复现，一字不差。**
并且这一发前态不是我引用 `004c6ec` 的旧账：`git diff --numstat 004c6ec HEAD -- internal/panel/`
= `947 0 internal/panel/l2_grant_boundary_test.go` **一枚路径、只增不减** ⇒ 前后态的差**只有这枚文件**，
可比性是构造出来的，不是假设的。台账 `A2xx`（`:5873`）独立记过 `RUN=152/PASS=90/FAIL=2/SKIP=0` 同一发，
两路对上。

### §0.3 名册差集（可复算）

```
comm -13 before after  ->  14 枚（5 顶层 + 9 子测试，逐名见本件末）
comm -23 before after  ->  0 枚（名册没有缩水）
```

14 枚逐名：`TestAnsweredPanelRoutesCarryNoApprovalDecision`、
`TestNoInboundEnvelopeCanBindAnApprovalVerdict`、
`TestGrantWireShapesAreRefusedAtTheDoor`（+2 子）、
`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`（+3 子）、
`TestPlantedGrantWiringGoesRedInASnapshot`（+4 子）。

### §0.4 口径裁决：那个 71 是**别包的数**，本件没算错，是**派单把两包并成一把尺**

- 我这发 `distinct` 的规则（写在明处）：`grep -o '=== RUN[[:space:]]*[^ ]*'` 取第 3 字段 →
  去掉 `#N` 后缀 → `sort -u` → **含子测试**。前 76 / 后 90。
- 另一把尺：`grep -h '^func Test' internal/panel/*_test.go | wc -l` → 后 **51**、前 **46**
  （实现件 §7 用的正是这一把，51−5=46 成立，我现量对上）。
- ⚠ 台账里那枚 `71`（`docs/reports/pending-and-issues.md:5897`，"distinct 名册 **71→71** 未缩"）
  量的是 **`internal/observe`**，不是 `internal/panel`。⇒ **本件 declared 的"口径不一致"在盘上不成立**：
  两个数不是同一枚包的同一把尺，没有矛盾需要调解。真正该记的是**派单把两个包的读数当同一分母**递给了我这一程
  （与 `A217①(a)`、`A219②`、`A220③` 同族：**引用别人的读数不带"哪包哪锚"**）。
  我按上面两把尺各自报数，**不采纳"71"当 panel 的任何基线**。

### §0.5 那两枚 FAIL 的归因（不是我替它说话，是复算）

- 名字唯一：`TestC21DesignTokensFourWayAgree` 跑两遍 ⇒ 2 行。前态同枚同名同数 ⇒ **本批之前就在红**。
- 红因复现（不是转述）：`internal/panel/tokens_fourway_test.go:50` 写死
  `tokensCSSPath = "design/assets/tokens.css"`，而 `git status --porcelain` 现量该路径是 ` D`
  （owner 2026-09-24/25 把 16 枚 `design/` 旧原型移出工作树，**未提交**，属他的地界）。
- 本程**没修、没跳、没放宽**，也没把它算进本批的账。台账 `A208③` 的 `P1`＝interim 留红，与实现件同向。
- `grep -rn "t.Skip" internal/panel/` → **1 枚**，在 `attachments_test.go:191`（先存在，与本批无关）
  ⇒ **本件没有新增 Skip**，两向可判。

锚点刷新：本节全部读数取自 `dev` @ `dd2b528`（真树）与同刻的仓库外副本。

---

## §1 最要紧那一格：这枚钉**今天**买到了什么

实现件 §0.1 说这条边界"结构性没接线"，并据此把严重度从"当下止血"改成"买未来"。**我独立复算四路，四条全部成立，
其中两条我量到了比它更硬的形状。**

| # | 待验断言 | 我用的命令 | 现量 |
|---|---|---|---|
| 1 | `ParseComposerRequest` 无生产调用者 | `grep -rn ParseComposerRequest --include=*.go .` | 定义 `bridge.go:77`；**非注释非测试命中 = 0**。全部命中＝定义 1 + 测试调用（`bridge_test.go:66,83,86,92,102,112`、`l2_grant_boundary_test.go:715,717,734`）+ 注释 4 处（`bridge.go:16,74,106`、`composer_handlers.go:37`、`cmd/wisp/run.go:225`）⇒ **成立** |
| 2 | `rt.modeWrites` 只装配、无人读 | `grep -rn modeWrites --include=*.go .` | **全仓 3 行**：`run.go:223`（注释）、`:229`（字段声明）、`:378`（`rt.modeWrites = &panel.ModeWriteHandler{…}` 赋值）。**读它的一处都没有** ⇒ **成立**，且形状就是"装配了、从此没人碰" |
| 3 | `HandleModeRequest` 只有测试调用者 | `grep -rn HandleModeRequest --include=*.go . \| grep -v _test.go` | 非测试命中＝`composer_handlers.go:109`（注释）＋`:111`（方法定义）。全仓 12 命中里其余 10 枚在 `composer_handlers_test.go` ⇒ **成立** |
| 4 | **Go 里没有任何 WebView2 接收点** | `grep -rni "WebMessage\|add_WebMessage\|OnWebMessage\|TryPostWebMessage" --include=*.go .` ⇒ **0 命中**；`grep -rn PostMessageW --include=*.go .` ⇒ 只有 `internal/ball/{sta,tray,live_windows_test,win32_windows}.go` 的 Win32 消息（球/托盘，不是面板 IPC） | **成立，且比实现件说的更硬**（见下） |

**两发我加上去的硬度（都不在原证据件里）：**

(a) **主模块连 WebView2 的依赖都没有。** `grep -c webview go.mod go.sum` → **0 / 0**；
`go list ./...` 里**没有**任何 `scripts/spike/**` 包（它是**独立模块**，`scripts/spike/go.mod`）。
⇒ 全仓唯一真的拉起过 WebView2 的代码是 `scripts/spike/webview2-latency/main.go`（S0 的延迟探针，
`github.com/jchv/go-webview2`），它**不在被发布的模块里、也没注册任何消息回调**。
所以"Go 侧收不到网页发来的东西"不是"忘了接线"这一种偶然，而是**依赖图上今天就没有这条路**——
tickets 33/35 接线时要同时新增依赖，那是一枚会被 `go.mod` diff 看见的动作，比"顺手加个 case"响得多。

(b) **界面那一头今天仍在往这扇不存在的门发东西。** `frontend/src/lib/panel.ts:169-186`
（`requestApprovalResolution` → `postMessage({method:"panel.approval.request", correlationId, outcome})`）
在 53a1359 之后**一字未改**，我用文件比对量的：
`diff <(git show 53a1359^:frontend/src/lib/panel.ts \| sed -n '169,186p') <(sed -n '169,186p' frontend/src/lib/panel.ts)` → **rc=0（零差异）**。
而今天 `frontend/src/components/l2-approval-card.tsx:98` 仍在调它（拒绝键与关闭键，`send("refuse")`）。
⇒ **入站方向今天仍是"有人发、没人接"**，只是发的 `outcome` 值只剩 `"refuse"`。
这条不改变严重度判定（门还是没接），但它把"这枚钉防的是将来"收窄了一步：
**将来第一个接上的人，接到的第一个真实入站 JSON 里就会带 `outcome` 键**，
而这正是面 2 的 AST/反射两把尺**当场会红**的那一发（§3 我用变异复算过）。 ⇒ 严重度结论**维持**。

### §1.1 裁定：**成立（一枚合法的死人开关 / dead-man's pin），但必须带那句"今天不保护任何东西"——而它已经在了**

派单问的是一支二选一。我量的形状：

- 判"成立"的根据：这枚钉**测的就是"接线之后会怎样"的那两类东西**（哪些方法被答、那个封套能绑哪些键），
  而不是"当下有没有人在敲门"。它今天在绿位、且被 §2–§5 的变异证明过咬得动（不是恒绿）。
  按 `SPEC-12`/D22 的思路，防"接线那一秒无人报警"正是它被批下来的理由（台账 `A217②` 原话：
  "哪天真给那扇门接线，禁令就破了，而且破的时候没有任何仪器会响"）。
- **但"下一位读者会不会把 `Q-49` 划成已覆盖"这件事不能靠运气。** 我逐条查了三处登记面，**三处都在**：
  1. **文件自己**：`l2_grant_boundary_test.go:22-27`「The only thing that made it inert was Go-side」
     ＋ `:54-61`「WHAT THIS FILE DOES *NOT* COVER - read this before citing it as coverage」
     ＋ 子测试 D（`:885-900`）把"看不见"做成一条**会因扩射程而红**的断言。
     ⚠ 但我量到一件它**没有**的：文件里**没有一处写明"今天这条入站路径根本没接线、本钉当下收益为零"**——
     它写的是"历史那枚按钮今天 inert 是因为 Go 不答它的路由"，读者完全可能读成"门是活的、在正确地拒绝"。
     "**结构性没接线**"这句只在证据件 §0.1 与台账 `A220②` 里。⇒ **这一处我判"附条件"**：
     条件是**文件头缺一句"nothing is wired today"**（详见本件 §8 的修法建议；本程不改）。
  2. 台账 `A220②`：编排者自己已经主动把 severity 降回"今天不止血／实际收益＝零"，并在 `A220⑥` 列了不覆盖清单。
  3. 证据件 §0.1 末「严重度结论」段。
  ⇒ **"被划成已覆盖"的防线有两枚（台账＋证据件）站得住，文件内那枚差一句**。总裁见 §8。

锚点刷新：本节现量取自 `dev` @ `43a6043`。

---

## §2 变异台（本程怎么种的、怎么收的）

全部变异只落在**两枚文件**：`internal/panel/bridge.go`（生产码，5 发）与
`internal/panel/l2_grant_boundary_test.go`（被验物自己，3 发）。纪律：

- 动手前备份到仓外 `D:\tmp\panel-l2-accept-r1\backup\`，两枚原件的 `git hash-object` 与我进场读到的一致：
  `bridge.go = d2cd6362ecc6941a0cee58073a2bf8ab6669a590`（与实现件 §3.2 记录的**同值**，两路互印）、
  测试件 `= 5b618d14b742f60cac22d2e58a1f2863dd311c63`（＝HEAD blob）。
- 每发一条命令内完成"变异 → 跑 → `cp` 还原"，还原后当场 `git hash-object` 复算 ⇒
  **8 发全部回到 `d2cd6362` / `5b618d14`，`git status --porcelain -- internal/panel/` 每发之后为空**。
  未用 `checkout`/`reset`/`stash`，未 push，未碰 `frontend/**`、`design/**`、`internal/observe/**`、`tools/d22scan/**`（只跑了它们的只读检查）。
- ⚠ **一次失败的变异我没当成读数**：M10 第一次落地时我猜错了缩进，python 的 `assert` 当场拒写，
  那一发跑出来的"5 枚全绿"是**干净树**的读数，已作废重跑（重跑结果见 §2 表的 M10 行与 §3.4）。
- **"变异真的生效了吗"这一问我用探针量，不靠推测**：临时新建 `internal/panel/zz_accept_l2probe_test.go`
  打印 `knownComposerMethod(...)` 的真实返回值（该文件是 `_test.go`，被 `goSourceFiles` 排除，不污染 AST 扫描），
  跑完 `mv` 到我自己的 `removed/` 里留着。两向都有：干净树上 4 枚候选名全 `false`（负控），变异树下目标名 `true`（正控）。
  ⇒ 下面每一条"没红"都配一条"但是 Go 真的答了/真的绑了"的探针读数。

| 发 | 形状 | 结果 |
|---|---|---|
| MA | 实现件 §3.2 那枚（封套长 `Outcome json:"outcome,omitempty"` **且**守卫答 `panel.approval.request`），**整包跑** | 红＝4 枚本件 + `C21`(先存在留红) ⇒ **先存在 46 枚因它而红＝0 枚**，§6 复算 |
| M1 | 守卫 switch 内联字面量 `"panel.approval.resolve"` | **红**（3 枚） |
| M1b | 同上但换成**候选名单里没有**的名字 `"panel.review.approve"` | **红**（3 枚，含 `no Method* constant` 那一句）⇒ 派生真在起作用，不是硬编码名单在起作用 |
| M2+M3 | 同一条 case 里放 `dynRoute`（包级 **var**）与 `"panel.review."+routeTail`（常量拼接） | **4 枚生产面全绿**；唯一红是快照测试的**陈旧性引信**（`plant B has no anchor…`），不是"发现了送字门"。探针：两枚名 `true` |
| M4 | 第二枚函数 `extraComposerGateM4` 自己的 switch 答 `"panel.review.allow"`，`knownComposerMethod` 末尾 `return` 它 | **5 枚全绿、整包 rc=0**。探针：`knownComposerMethod("panel.review.allow") = true` |
| M6 | 新入站封套 `acceptM6Envelope{Cmd json:"cmd"; Outcome json:"outcome"}` + 真 `json.Unmarshal` 的处理函数 | **5 枚全绿**。探针：`Cmd="panel.review.allow" Outcome="grant"` 真的绑上了 |
| M7 | 新入站封套 `acceptM7Envelope{Method json:"method"; Outcome json:",omitempty"}`（**它确实绑 `method`**） | **5 枚全绿**。探针：线上键 `"outcome"` 绑定成功 `Outcome="grant"` |
| M8 | 同一形状放到**正主** `ComposerRequest` 上（`Outcome string json:",omitempty"`） | **反射那半红、AST 那半 0 findings** ⇒ 两台"逐字同规矩"的仪器当场分家（见 §3.2） |
| M9 | 词表砍掉 `"outcome"` | 自校验 1 子红 + 快照 A/C 红；三面生产面绿（今天真树没有 outcome 键可找） |
| M10 | 反射那半**去掉对内嵌结构的递归** | **只有自校验的第 3 枚植物红**，其余 4 枚顶层全绿 ⇒ §5 的承重凭据 |
| M11 | 路由词表加一枚 `"mode"`（放宽） | 面 1 红 + 面 2 红 + 自校验红 ⇒ 自校验的"合法路由不该被命中"那一半是**镜子** |
| M12 | `ParseComposerRequest` 不再问守卫（删掉那 3 行 `if !knownComposerMethod(...)`），守卫函数原样留着 | **只有面 3 红**（4 枚线上形状被放行），面 1/2/4 全绿 ⇒ 面 3 的承重凭据 |

---

## §3 四面逐格判

### §3.1 facet 1「被答方法集是派生的、不是抄的」——**附条件成立**

**先说它站住的部分（我给足）：** AST 那台仪器**确实**不是抄名单。M1b 是决定性的一发：
我用的是**它 11 枚硬编码候选里没有的名字** `panel.review.approve`，
红句是 `l2_grant_boundary_test.go:569`（派生集的词表判）与 `:604`（"没有一枚 `Method*` 常量声明它"），
**不是** `:590` 那条候选名单的读数（虽然那条也同时红）。⇒ 派生链（`switch` case 标签 → 经本包 `const` 解析）
是活的、能咬的，这一条实现件没吹。

**再说 finding。四条攻击的路径与读数：**

| 攻击 | 造形 | 打红了吗 | 判 |
|---|---|---|---|
| (i) case 里内联字面量 | M1 / M1b | **红了** | 站住 |
| (ii) 路由经**变量/计算串**可达 | M2+M3 | **没红**（4 枚生产面全绿，只有陈旧性引信叫） | **F-1** |
| (iii) **第二枚 switch/if 链**在别处答同一个方法名 | M4 | **完全没红**（5 枚全绿、rc=0） | **F-2（最重）** |
| (iv) 路由**在 `knownComposerMethod` 之外**被答（守卫不再被问） | M12 | 面 3 红、面 1/2/4 绿 | **F-3**（面 3 兜住了行为，词表面没兜住） |

根因是**两句可复算的代码事实**，不是我的推测：

- `routeNameOf`（`:429-441`）只认两种 case 标签：`*ast.BasicLit`（字符串字面量）与
  `*ast.Ident` **且**能在 `pkg.consts` 里查到。而 `pkg.consts` 只在 `decl.Tok == token.CONST` 时填
  （`:294-296`）⇒ 包级 **`var`** 与**任何表达式**（`"a"+b`、函数调用）解析不出来，`routeNamesInFunc`
  就把它**当不存在**（不报"读不懂"，是静默丢弃）。
- 那条互印**只朝一个方向**跑：`:564-567` 是 `for _, name := range answered { if !knownComposerMethod(name) … }`
  ＝ **派生 ⊆ 运行期**。反方向"运行期还多答了谁"**在结构上无法被这行发现**（它没有别的枚举源）。
- `parseBoundaryPackage` 只把**函数名恰好等于 `knownComposerMethod`** 的那一枚函数里的 switch 当路由来源（`:318`）。
  `guardSeen` 那道引信只在**改名/删掉**时叫（`scanGrantBoundary:485`），**被绕过时它不叫**——
  M4 就是"被绕过"，M12 是"被不闻"，两发都绿。

⇒ 文件头那句 "the set of inbound methods Go actually answers - derived from this package's own route guard …
and **agrees with** knownComposerMethod at runtime"（`:40-42`）在盘上是**过强的**：
真实语义是"某一枚具名函数的 switch 里、能写成字面量或包级常量的那批标签，且单向 ⊆"。
**这一格我判附条件成立**：正向（"抄的名单看不见新门"）被 M1b 证否了、是真派生；
但"看不见的那类门"存在，且其中 M4 那一类正是 tickets 33/35 接线时最可能写的形状（处理函数注册表）。

### §3.2 facet 2「入站封套绑不出结论键」——**附条件成立**，另带一枚**仪器缺陷**

它的结构判据是"绑 `method` 键的类型即入站"（`inboundEnvelopes:447-474`，用 `strings.EqualFold(f.JSONKey,"method")`），
再向下闭包取内嵌/嵌套类型。两条攻击：

- **F-4（判据本身）**：M6 造了一枚**绑 `cmd` 不绑 `method`** 的入站封套（`{Cmd json:"cmd"; Outcome json:"outcome"}`），
  带真 `json.Unmarshal` 的处理函数。探针证明线上键真的绑得进去（`Outcome="grant"`），
  **5 枚测试全绿**。⇒ "换一枚键名承载路由"就离开射程。这不是吹毛求疵：文件自己在注释里说
  "路由字段自己就是标记"（`:446`）——那枚标记**恰好就是它可以被换掉的那一枚**。
- **F-5（仪器缺陷，比 F-4 更硬）**：`jsonOr`（`:377-382`）在**有 tag 但 tag 名为空**
  （`json:",omitempty"`）时返回 `""`，而 `encoding/json` 在这种写法下**回落到 Go 字段名**，
  线上键 `"outcome"` 照样绑得进来（M7 探针实测 `Outcome="grant"`）。反射那半算对了
  （`bindableKeysOf:152-172`：`jsonName==""` ⇒ `key = f.Name`）。
  **M8 把这处分家钉出来了**：同一形状放进正主 `ComposerRequest` ⇒ **反射红、AST 那半 0 findings**。
  ⇒ 两台自称"同一条规矩的两种写法"（文件头 `:44-46`、证据件 §1 表第二行"反射 + AST 双写"）
  **今天并不等价**，AST 那台在空 tag 名上是**假阴性**。
  ⚠ 方向要说准：这个缺陷**只会少报、不会多报**，所以它不是"放水"，是"自称双写、实为一台半"。
- 顺带一条**不是缺陷的观察**：反射只对 `ComposerRequest`/`ModeRequest`/`AttachmentPayload` 三枚跑，
  AST 才覆盖全包 —— 两台各自的覆盖面不重叠（M6/M7 都在 AST 的覆盖面里而都漏），
  所以"双写"买到的冗余比它字面承诺的小。

### §3.3 facet 3「行为面拒收历史线上形状」——**成立**，但"历史字节"那一枚不是它的牙齿

先复推它承诺的"同一串字节"（派单要求我别信转述）：

- `git show 53a1359^:frontend/src/components/l2-approval-card.tsx` 的 `:161-167` **逐字复现**
  （`:163` = `onClick={() => send("grant")}`，`:166` = 「本次允许」；53a1359 的 numstat = `0 7` 纯删 7 行），
  `send` 在 `:96-100`、`:98` 调 `requestApprovalResolution(view.correlationId, outcome)`。
- `frontend/src/lib/panel.ts:169-186` 在 53a1359 里**一字未改**：
  `diff <(git show 53a1359^:…panel.ts|sed -n '169,186p') <(sed -n '169,186p' …panel.ts)` → **rc=0**。
  ⇒ 发出的字节 ＝ `{"method":"panel.approval.request","correlationId":"<id>","outcome":"grant"}`，
  **键集与键序都和钉住的第 1 枚一致**。所以**"钉的是历史那串字节"这句成立**（保真度没问题）。

但**诊断力**这一层有一条要报：

- **F-6**：M12（把门整个打开）之下，5 枚钉住的线上形状里 **4 枚变红、第 1 枚仍绿**——
  仍绿的正是文件注释里点名"verbatim from the deleted button"的那一枚，
  因为它**不带 `requestId`/`source`**，守卫不问了它仍会被身份检查拒掉。
  ⇒ "我们拒收历史形状"这面旗实际由**第 2 枚**（补了身份字段的那枚）扛着；
  第 1 枚是**装饰性的 witness**。实现件在注释里诚实写了"plus the identity fields the current envelope demands"，
  所以这不是隐瞒，但下一位读者会以为"历史字节被打过靶"。
- 一条小口径：`:746-749` 那句 `strings.Contains(string(back), "\""+key+"\"")` 是**整串子串匹配**，
  键名与**值**都能命中（M8 的红句里它命中的其实是 `"Outcome":"grant"` 那个值位）。
  方向是**多报不是少报**，判可接受，但它不是"只按键判"的尺，措辞别写成只按键。
- **承重证据（面 3 独有的那一发）**：M12 之下**只有面 3 叫**。
  删掉面 3 ⇒ "主机不再问守卫"这一类变异在全树零仪器覆盖。⇒ 面 3 = **承重**。
  （M8/M9 的红句也都在面 3，但那些枚别的面也叫，只有 M12 是面 3 独有。）

> **§3.3 就地更正（追加，原句不删；这是本程自己的一枚假阳性结论）**：
> 上一版那句"M12 之下**只有面 3 叫**、删掉面 3 这类变异全树零仪器覆盖"**是错的**。
> 我当时只跑了本件的 5 枚测试就下了全称结论；补跑**整包**（`M12-fullv.txt`）现量：
> 红的顶层是 3 枚 —— `TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests`（**先存在**，
> `bridge_test.go:77`）、`TestGrantWireShapesAreRefusedAtTheDoor`（本件）、`TestC21DesignTokensFourWayAgree`（留红）。
> ⇒ 先存在那枚的子测试 `:78` "a request that names an approval decision is not a composer method"
> **也在这一发上叫**（它断言 `approval.decide` 被拒且错误句里含"不是面板 composer 通路"）。
> **改判**：面 3 对 M12 **不是唯一**防线，是**第二**防线；先存在那枚是**单一拼写的字面匹配**
> （它只写了 `approval.decide` 一个名字，所以 MA 那发答 `panel.approval.request` 时它不响，见 §4.2）。
> 面 3 的承重**判据本身不变**（M12 确实只有行为面红、词表三面绿），但"全树零仪器覆盖"那句作废。

### §3.4 facet 4（本件的词表自校验）——**附条件成立：三枚植物里两枚承重，"合法路由不被命中"那一枚是镜子**

派单给的 4 号格子问的是"承重还是镜子"，我这仓的定义＝**摘掉它，有没有哪发变异从此打不红**。三发定向变异：

| 发 | 造形 | 本件的自校验 | 其余三面 + 快照牙齿 | 判 |
|---|---|---|---|---|
| **M10** | 反射那半**去掉对内嵌结构的递归**（`bindableKeysOf` 的 `f.Anonymous` 分支不再下钻） | **只有第 3 枚植物红**：`:683 planted a grant-carrying envelope (a verdict smuggled inside an embedded struct) and the predicate reported nothing - the check is decoration` | **面 1/2/3/4 全绿**（含快照 A/B/C——植物 A 是**直挂字段**，不是内嵌偷渡，所以咬不到这一发） | **承重（唯一防线）** |
| M9 | 词表砍掉 `"outcome"` | 第 1 枚植物红 | 快照 A + C 也红（`:833`/`:863`） | 冗余，但方向对（两处都会叫） |
| M11 | 路由词表**放宽**加一枚 `"mode"` | "cry wolf" 那一枚红（`:697`） | **面 1（`:569`）与面 2 的 AST（`bridge.go:97`）同时红**，用的是同一句谓词、同一批 4 枚名字 | **镜子**（摘掉它这发照样红） |

⇒ **结论：`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` 整体＝承重**，凭据是 M10 那一发
（**摘掉这枚测试，M10 那类"内嵌偷渡"变异从此在全包零反应**）；
但它的**第四条腿（4 枚合法路由不被命中）是镜子**——它与面 1 的词表判是同一个谓词打在同一个名单上，
冗余不加分。"三枚植物必须被点名 + 四枚合法不许被点名"这两向里，**只有前两向（植物）在挣钱**。
实现件 §1 表里那句"谓词自己不能退化成匹配不到任何东西"是准的；
它没有声称第四枚腿独挡任何变异，所以这不算吹，只算**可以删掉一行而不丢任何牙齿**。

另记一条 M10 的**顺带产物**：全包唯一在"内嵌偷渡"这条轴上有覆盖的，就是这枚自校验。
M10 之所以是活的可疑形状，是因为 `ComposerRequest` 今天**真的内嵌**了 `AttachmentPayload`（`bridge.go:70`），
"把结论键藏进内嵌类型"是那枚封套上唯一一处不需要新写 `method` 键就能加字段的地方。

锚点刷新：§2–§3.4 的变异与读数全部取自 `dev` @ `43a6043`（`bridge.go` 每发之后 `=d2cd6362…`、
测试件每发之后 `=5b618d14…`、`composer_test.go` W1 前后 `=ecfd0f7a…` 同值、
`git status --porcelain -- internal/panel/` 逐发为空）。

---

## §4 那一格侧读："真变异之下 46 枚先存在测试零反应"——**复现成立**，但它是**一半辩护词、一半被说过头**

### §4.1 复算（MA 发，整包、`-count=1 -v`，输出 `out/MA-full.txt`）

```
变异 = ComposerRequest 长 Outcome json:"outcome,omitempty" + 守卫 case 追加 "panel.approval.request"
红的顶层（5 枚）：TestAnsweredPanelRoutesCarryNoApprovalDecision / TestNoInboundEnvelopeCanBindAnApprovalVerdict
                 / TestGrantWireShapesAreRefusedAtTheDoor / TestPlantedGrantWiringGoesRedInASnapshot   <- 本件 4 枚
                 / TestC21DesignTokensFourWayAgree                                                    <- 先存在，改前已红
先存在 46 枚中因这枚变异而红 = 0 枚        （46 = 全包 51 枚 func Test - 本件 5 枚，我现量）
```

⇒ **实现件 §7/§8 那句"因变异而红的先存在测试＝0 枚"独立复现，一字不差**。
且它 §8 自己把措辞从"46 枚一枚都没红"收窄成"没有一枚是因它而红"——**那个收窄是对的**，
我在 §0.5 量的 C21 就是"改前改后都红"的那枚，别把它算进任何一方的账。

### §4.2 判：是这枚钉的**辩护词**，也是对**同族三枚仪器**的一份**具名指控**（不是对整个测试包的）

我不接受"46 枚零反应 ⇒ 测试包很糟"这种整体归因（"整树 rc=1＝工具链假象"那族的最省事藏身处，反方向同理）。
把这一格拆成三枚**具名**仪器，各判一条：

| 具名仪器 | 它对不对得上这个形状 | 我量到的证据 |
|---|---|---|
| `TestTheRendererHoldsExactlyOneDoorToTheHost`（`composer_test.go:502`） | **对得上轴、方向反了**。它第 (iv) 条判"每枚 `panel.*` 字面量都是 Go 答的那条路"，而它判"Go 答哪条路"用的**神谕是一份抄来的白名单** `composerRouteLiterals()`（`:409-419`），里面白名单了 `"panel.approval.request"`（`:417`） | 它今天**只读 `frontend/src`**、全文没有一处读 Go 的路由守卫 ⇒ MA 那发改 Go，**它结构上不可能响**（这就是派单问的那格，答案：不该由它响，它响不了）。**W1 变异**（只删 `:417` 那行白名单，别的一动不动）⇒ 它当场红、点名 `src/lib/panel.ts:181 method: "panel.approval.request"`，还原后 `ecfd0f7a…` 同值。⇒ **那枚白名单今天正在为一根活线放行**：界面仍在发这条路（拒绝键/关闭键，`l2-approval-card.tsx:98` 我现读），发的字节里就带 `outcome` |
| `TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests`（`bridge_test.go:77`，子测 `:78` 名 "a request that names an approval decision is not a composer method"） | **对得上，但只钉了一枚拼写**。它是全包唯一先存在的、判"审批结论词不许被这条路答"的行为腿 | MA 之下它**不响**（它测 `approval.decide`，变异答的是 `panel.approval.request`）；**M12 之下它响**（见 §3.3 更正）。⇒ 这正是 owner 那句"要结构性、不要字面匹配"的**具体理由**：字面匹配的单枚拼写已经被写过一次了，它挡不住换拼写的下一发 |
| `TestFrontendNeverNamesAnApprovalDecision`（`frontend_hygiene_test.go:281`，ban #6 的包内拷贝） | **两向都看不见**：判据是正则 `approval\.decide`（`:73`），走的是 `frontend/**` | 历史违规文件里 `approval.decide` **0 命中**（我先用 `grep -cF 'send("grant")'` = **1** 证明尺子活着，再量那枚 0），今天 `frontend/src/` 亦 0；`tools/d22scan/main.go:149` 同一条正则、消费点 `:814`（两枚行号我按 HEAD 现读，与实现件 §4.2 引的一致） |

⇒ **我的结论形状**：这枚钉**是**它的辩护词（它填的确实是全包唯一那个洞：读"这条路上递的是什么字"）；
但"46 枚零反应"**不能读成"全树零覆盖"**——同一格里 `bridge_test.go:78` 是一枚**已经存在、只是太窄**的腿，
`composer_test.go:417` 是一枚**已经存在、正在放行**的白名单。
**最该有反应而没有的那一枚是 `composer_test.go:417`**，而修它要动的是前端那半
（`ApprovalOutcome` 收窄 + C17 白名单），台账 `A217⑤`/`A220⑥` 已记成待切片卡
⇒ **本程不改、不扩，只把"零反应"的三个层次分开钉住**。

### §4.3 一条我自己造的读数纪律（写给下一位，不是谁的缺陷）

§3.3 那条更正就是这一格的后半课：**"整包有没有别的反应"必须整包跑**。
我第一发 M12 只 `-run` 了本件 5 枚就下了"全树零仪器覆盖"的全称结论，那是**我这程的假阳性**，
补跑整包当场推翻（`bridge_test.go:77` 也红）。⇒ 固定动作：
**凡结论句式里出现"只有／零／唯一"，取数命令必须覆盖全部分母，不许用 `-run` 子集。**

---

## §5 契约轴 + 门（六项逐项现量）

| 项 | 判据 | 现量 |
|---|---|---|
| 落地范围 | 本批只允许碰自己的路径 | `git show --numstat d88c356` = `126 0` 证据件 + `945 0` 测试件；`660ffa1` = `201 1` + `4 2` ⇒ **两枚 commit 只带这 2 枚路径，生产码零字节** |
| 断言只增不减 | `internal/panel` 区间内删除行 | `git diff 004c6ec HEAD -- internal/panel/` 删除行 **0 行**（净态 +947/−0）；`660ffa1` 里那 2 行被删的是 `:645-648` 的 `t.Logf`（日志、非断言），改后仍打印同一句、只是加了条件 |
| 无新增 `t.Skip` | `grep -rn "t.Skip" internal/panel/` | **1 枚**，在 `attachments_test.go:191`（先存在）；本件 **0 枚** |
| 阈值 / golden / `thresholds.go` | `git diff --name-only 004c6ec HEAD` 全清单里**无** threshold/golden 字样 | **空** ⇒ 一字节未动；区间内被改的 25 枚路径全在 `design/doubao/**`(14，别家)、`docs/**`(7)、`internal/panel`(1)、`tools/d22scan`(2，别家) |
| `allowlist.txt` | `git rev-parse 004c6ec:… HEAD:…` | 两枚 blob **同值** `6b61fad57085a52a62d4ff513209a817b96a3067` |
| `emojiRe` | 区间内 `tools/d22scan/main.go` **根本没出现在改动清单**里（那里只动了 `gitignore.go`、`scan_test.go` = 另一程） | **未动**，与 `A221⑦` 的"五锚同值"同向 |

### §5.1 门与格式（真跑）

```
gofmt -l internal/panel/  ->  空（rc=0）
go vet  ./internal/panel/ ->  空（rc=0）
sh scripts/d22scan.sh     ->  rc=0
   runtests.sh: OK packages=[./...] top-level: PASS=29 FAIL=0 SKIP=0, === RUN=69
   d22scan: clean - no D22 ban violations
   live scope: bans#1-5 internal/=203 cmd/=22 / #6 frontend/=40 / #7 internal/tools/=18
              / #8 design/=32 frontend/=40 internal/=407 cmd/=39
```

- `ban #8 internal/ = 407` ⇒ 与 `A220⑤`（`406 -> 407`，因本件新增一枚 `*_test.go`）**同值**：分母进新文件，不是回归。
- `design/=32`（台账 `A218⑥` 写 30）＝本区间内 `design/doubao/**` 被另一程动过 ⇒ **与本钉无关，只登记口径**，
  正是 `A220⑤` 那句"引用时现跑一遍、把当时 HEAD 一起写进判据本体"的形状。
- `PASS=29 / RUN=69`（实现件 §5 记 `28 / 68`）＝别家 `tools/d22scan` 那批又落地一枚测试，不是本件造成。
- 本件新文件**未被 ban #8 点名**（整树 `clean`，零 finding）。

### §5.2 两处**指向不存在测试**的引用（一枚本批新抄，一枚先存在）

- **`TestComposerMethodNamesMatchFrontend` 这枚测试在仓里不存在。**
  正控先跑（同族真名 `TestComposerContractTypesMatchFrontend` → 3 命中）；本名 `grep -rn` 命中 **3 处引用、0 处定义**
  （`grep -rn "func TestComposerMethodNamesMatchFrontend"` → 0）。三处引用：
  `internal/panel/bridge.go:33`（**先存在**，票 92 留的）、
  **`internal/panel/l2_grant_boundary_test.go:604`（本批新抄进红句里的）**、实现件 §3.2 的引用块。
  ⇒ 实际干那件事的测试叫 `TestFrontendComposerRequestsMatchTheEnvelope`（`bridge_test.go:131`）。
  **后果说准**：这条红句是给"未来把路由直接写进 switch、绕过命名门"的人看的，它把人指向一枚不存在的门禁。
  判 **缺陷（措辞级），不判退回**——断言本身有效（`:604` 在 M1b 下真红过，§3.1 引的就是它）。
- 同族一处（**先存在**，不属本批）：`composer_test.go:415-416` 那句
  "The approval card's route, whose exact spelling is pinned by `TestFrontendComposerRequestsMatchTheEnvelope`"
  **不成立**——我读了那枚测试全文（`bridge_test.go:131-175`），它从不提 `panel.approval.request`，
  只核四枚 composer 方法在界面出现、核 `bridge.postMessage` 计数 == 2。
  ⇒ 那行白名单**既在放行、又自称有主**（与 §4.2 第一行合起来读）。**本程不改**（动它是 `composer_test.go`/票 92 的地界）。
