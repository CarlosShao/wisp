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
