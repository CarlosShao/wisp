# 161-gates-accept-r2 —— 票 161 的 AC#4／AC#6／AC#7 三格，非实现者裁决表（代号 161-v2）

- 派单：`.scratch/wisp/dispatches/2026-09-27-100x-accept-161-v2-ac4-ac6-ac7.md`
- 被验的三格：AC#4（反向判据）· AC#6①②（聚合退码＋两族成对普查）· AC#7（取证机械撞全仓门的顺序＋归因仪器）
- 我是**非实现者**：本表**不勾任何框**、不改任何码、不放宽任何断言、不 `t.Skip`。
- 实现者分别是 161-r4（AC#4／AC#7②）、161-r6（AC#6）、161-r3／161-r5（AC#7①与 CI 那一步）；本程与它们无 implements 关系。

## 0. 锚点与被验版本（起手三件）

- **锚＝`30121fef17a42791eadab463c4773325a98838e2`**，`git rev-parse HEAD` 现量（09-27 09:5x），`git cat-file -t <锚>` ＝ **`commit`**。分支 `dev`。
- **脏工作树的处置**：本锚上 `git diff --name-only <锚> -- .scratch/wisp/probes/ .github/workflows/ci.yml tools/d22scan scripts/ internal/panel/` 只返回一枚
  `.scratch/wisp/probes/152/my152.py`（别程的半件，本程未动）⇒ **本表被验的全部仪器在盘上＝锚上同一份字节**，
  这条读数本身就是"我没读脏树"的凭据。
- 并程写面（票 162 正在写 `internal/tools/**`）：本程**一字节都没读过工作树里的 `internal/**`／`cmd/**`**——
  所有生产码读数走 `git show <锚>:<path>` 与 `git grep <锚>`（⇒ 按定义读的是**提交树**）。
  成对普查那两族的尺（`gate-clauses.sh`）也是 `git grep "$A"`，与工作树无关，脚本第 11 行自陈同一件事，本程核 its behaviour 而非 its words。
- **未复跑／不在我射程**：`frontend/**`、`design/**` 未查未引（不是本编队的写面）；`probes/161/r1..r6/**`、`probes/154/**` 只读。

## 1. AC#6 —— 聚合退码那一腿 ＋ 两族成对普查两腿

档位标记：**〔我本轮现跑过〕**＝我自己在锚上跑出来的读数；**〔日志＋归档，抽验〕**＝我读了实现者的日志并抽查其中一条自己复跑。

### 1.1 第一问：聚合有没有牙 —— **成立〔我本轮现跑过〕**

我没有跑 r6 的 `flip-declaration.sh`（理由见 §1.4：它把日志写进 `probes/161/r6/logs/`，那是别人的已提交件，本程写面外）。
我自己做的是：**把 `gate-clauses.sh` 逐字节 `cp` 到 `/tmp/wisp161v2/flip/base.sh`，用 `sed` 只翻 `want` 行，`cmp` 确认改动非空，再显式把锚当 `$1` 递进去跑**（不改 `$A` 的默认取值方式）。

| 发 | 翻动 | 声明 vs 实测 | 聚合退码 |
|---|---|---|---|
| 基线 | 未翻 | 14 腿全 ok | **rc=0** |
| 甲（该静而实测响） | `want G6 ring rev` → `want G6 quiet rev` | `# BAD 腿=G6 声明=quiet 实测=ring` | **rc=1** |
| 乙（该响而实测不响） | `want G4 quiet` → `want G4 ring` | `# BAD 腿=G4 声明=ring 实测=quiet` | **rc=1** |
| 丙（同时翻两枚） | 甲＋乙 | 两条 `BAD` | **rc=2** ⇒ 退码是**说谎腿的计数**，不是布尔 |
| 还原 | 未翻（原字节） | 14 腿全 ok | **rc=0**，且名册与基线**逐行只差 `# 生成时刻` 那一行** |

⇒ "某腿点响、聚合仍退 0" 那一形（票 158 交来的原始缺陷）**在这把尺上已经不成立**：聚合吃的是"声明与实测不符"，
两向都能被一发翻动打红，且响而不符的腿（G2／G5／G5neg／G6／G6neg／G7／G7neg 共 7 枚今天按设计该响）**不**把退码顶起来。
基线 14 腿逐枚 `ok`、`腿数＝14 声明与实测不符＝0` 与我同锚另一次独立跑（`sh` 原件、未翻）一致。
读数落盘：`.scratch/wisp/probes/161/v2/aggregate-flips.txt`。

### 1.2 第三问：会响证据（甲）是不是真把尺指向了那棵树 —— **是〔我本轮现跑过〕**

我**另建**了一棵仓外合成树（`/tmp/wisp161v2/faketree`，`git init`，三枚提交），
pathspec 与主尺同形（`internal/**/*.go`＋`cmd/**/*.go`），样本是**我写的**、不是 r6 的 `samples/`：
`internal/fake/taskopen.go`（只 `OpenTask`）、`internal/fake/defersonly.go`（只 `x.Defer(...)`，文件里没有 `DisposalScope` 这个词）。
在**那棵树的 cwd 里**跑同一份 pristine 尺（`gate-clauses.sh` 第 13 行 `cd "$(git rev-parse --show-toplevel)"` ⇒ `$A` 解析到合成树的 HEAD）：

- **phase1（开了没合）**：G6 主尺 `# UNPAIRED internal/fake/taskopen.go (开方调用点=1)`、`# 未成对枚数＝1`、`# rc=1`；
  G7 主尺 `# UNPAIRED internal/fake/defersonly.go`、`# 未成对枚数＝1`、`# rc=1`。⇒ 两族**正向**都点得到新造的未成对样本。
- **phase2（同一文件补上合方）**：两族 `未成对枚数＝0` **且 `# git grep rc=0`** ⇒ 名册非空的那种 0＝"真的配上了"。
- **phase3（在合成树里删掉样本并提交）**：两族仍 `未成对枚数＝0` 但 **`# git grep rc=1`** ⇒ 与 phase2 的 0 **分得开**。
  落盘：`.scratch/wisp/probes/161/v2/faketree-phase{1,2,3}-g6g7-mainlegs.txt`。

⇒ 票面 AC#6② 那句"加腿之后必须能拿一发新造的未成对样本让它响；响不了＝装饰"，**我这一发独立复现成立**（不依赖 r6 的样本、不依赖它的日志）。
r6 自己的 `faketree.sh` 我**未运行**（同上，它写自己的 logs 目录），它的 phase 读数属〔日志＋归档，我未复算，但上面那发主证我自己跑过〕。

### 1.3 第二问：那 4 枚反向读数，逐枚裁（这是派单留给我打的洞）

先说**判据**（可复算，三问都答完才算数）：
（一）被点名的那一行**是不是代码**；（二）如果不是配对的另一半，**负责关它的那一行在哪个文件的哪一行**；
（三）那条负责路径**是否在所有出口都会走**。三问全过＝形状使然；第（三）问过不了＝真漏，单立账。

| # | 读数（锚上现量） | 裁 | 复算凭据 |
|---|---|---|---|
| 1 | `cmd/wisp/run.go`（G6 反向，合方调用点=1） | **形状使然（跨文件分工），不是真漏** | 合方＝`cmd/wisp/run.go:563 rt.bridge.CloseTask(taskID)`，它在 `admitTask` 返回的 revoke 闭包里（`run.go:557-565`）；开方＝`internal/tools/bridge.go:559 b.OpenTask(dec.TaskID)`，在 `mark()` 里**按调用**惰性开。负责走的出口＝`internal/agent/loop.go:366 defer revokeAdmission()`（我读了 355-367 段确认是 `defer`，panic/cancel/error 出口都覆盖）。两侧不同包这一点由 `bridge.go:629-632` 的注释自己写明（"open per CALL, close per TASK … nothing keeps the two sides in step by construction"）。⇒ 反向那一枚是**判据想要的形状的定义**（开方在射程外的定义本体文件里），不是没关。 |
| 2 | `internal/memory/retention.go`（G7 反向，合方调用点=2） | **形状使然，且是判据词表缺一味** | 两行合方＝`retention.go:98 func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, ...)`（签名）＋ `:100 panic("...requires a DisposalScope")`（字符串）——**都不是"拿了一个没收尾的 scope"**。真正的收尾登记在 `retention.go:103 scope.Go("retention-job", ...)`，而 `Go` 是 disposal 注册动词之一（`internal/plugin/disposal.go:226`），**G7 的开方词表只写 `Defer(Named)?`（`gate-clauses.sh:293-294`），所以这一枚按构造永远不会成对**。负责方＝构造 scope 的人；而我现量 `git grep -nEw NewDisposalScope <锚> -- internal cmd 排测试排定义本体` ⇒ **rc=1，生产码今天一枚构造点都没有**（与 r6 台件注释同一读数）。⇒ 该关没关这一支**今天不成立**。 |
| 3 | `internal/observe/goroutine.go`（G7 反向，合方调用点=1） | **不是代码，是同行尾随注释** | 被算进调用行的那一行＝`goroutine.go:34 CategoryTemporary GoroutineCategory = "temporary" // inside a DisposalScope`。`pair()` 只剔"行首是 `//`"的整行（`gate-clauses.sh:107`），**尾随注释豁免不了**（r6 自己在 `:290-291` 记下了这味粗糙）。该文件除 `:54` 的整行注释外**没有任何 `DisposalScope` 的代码出现点**。⇒ 零枚真漏，也零枚真读数。 |
| 4 | `internal/proc/shutdown.go`（G7 反向，合方调用点=1） | **同上，尾随注释** | `shutdown.go:85 ReleaseSpeechSessions func(ctx context.Context) error // step 5 (speech, via DisposalScope)`；该文件另两处（`:23`、`:161`）是整行注释、已被剔。⇒ 零枚真漏。 |

**结论：4 枚里没有一枚是真漏 ⇒ 不需要单立泄漏账**（派单 §2 第一支不触发）。
但**这一格不能就这么结案**，因为我裁出来的东西比"是形状"更难听：

> **G7 那条主尺今天的 `ring` 声明是真的，可它为真的理由不在判据想量的那个东西上。**
> 3 枚反向里 2 枚是**注释 token**、1 枚是**函数签名/字符串**；代码级的"挂了收尾却没绑到 scope"一枚都没有。
> 也就是说：G7 现在**不是**一枚会漏掉真漏洞的尺（合成树那发证明它正向吃得住），
> 但它是一枚**靠注释响**的尺，而且它的 `n` 把正向与反向**加在同一个数里**（`gate-clauses.sh:142`），
> 所以聚合只看"枚数≥1"，**看不出响的是哪一向**。

**可复算的会响条件**（把"为什么该响"从形容词变成命令；未来任何一程可原样复跑）：

- **G6 该响** ⇔ 存在 `f ∈ {git grep -nEw 'OpenTask|CloseTask' <锚> -- internal/**/*.go cmd/**/*.go :!*_test.go :!internal/tools/bridge.go 的行，剔行首注释}` 的文件集，使
  （正向）`f` 含 `OpenTask` 而不含 `CloseTask`，或（反向）`f` 含 `CloseTask` 而不含 `OpenTask`。
  锚 `30121fe` 上：正向 0 枚（开方在被排掉的定义本体里，**这条射程结构上没有开方可数**）／反向 1 枚＝`cmd/wisp/run.go`。
  ⇒ **这一腿的"该响"今天只由反向一枚撑起，且那枚的原因就是"开方在射程外"**；正向的会响证明是合成树 phase1（我已复跑）。
- **G7 该响** ⇔ 同一公式，开方 `Defer(Named)?`／合方 `(New)?DisposalScope`，排 `*_test.go` 与 `internal/plugin/disposal.go`。
  锚 `30121fe` 上：正向 0 枚（生产 `Defer` 调用点全在定义本体）／反向 3 枚＝上表 #2#3#4，**其中 #3#4 不是代码行**。
  ⇒ 复算命令：`sh .scratch/wisp/probes/154/gate-clauses.sh <锚>` 取 `## G7 主尺` 段的 `UNPAIRED-REV` 行，再对每一枚跑
  `git show <锚>:<file> | sed -n '<line>p'`，**肉眼分"代码里的词"与"注释里的词"**——这一步是本格今天缺的那一步，
  尺自己没有做（它只在整行级剔注释）。

**这条声明会在什么情况下被未来一程误用成豁免**（派单要求的反面，我写成三条）：

1. **"已知该响"被读成"已知在册、可以不理"**：聚合退码只看"声明==实测"，所以**在声明仍是 `ring` 的整段期间，这一腿对任何新长的未成对枚数都退 0**。
   今天 `n` 只有 1／3 枚，明天变成 1／9 枚，聚合还是 `ok`——**枚数不进判据，只有 0/非 0 进**。这是本格里最实际的一枚豁免风险。
2. **把 `ring` 钉成"这 3 个文件有活洞"**：上表裁的是**锚上的读数**，不是那三个文件的许可。谁拿声明去替 `retention.go`/`shutdown.go` 的新增未配对背书，就是把我这张表反着用。
3. **判据被修好时反向咬人**：若下一程把 `pair()` 的剔注释改成"同行尾随 `//` 也剔"（**正确的修法**）或把 `Go` 并进开方词表，
   G7 立刻变 `quiet` 而声明还是 `ring` ⇒ 聚合 **rc=1**，看起来像"引入了回归"，实际是尺的精度变好了。
   `gate-clauses.sh:174` 那句"逐枚都是 09-27 在锚上现量登记的今天真值"已经承认声明是**锚点级快照**，⇒ 结论：
   **任何改动 `pair()`/`diffsets()` 词表或剔注释规则的一程，必须在同一发提交里重登记 `want` 行**，否则聚合红的是"声明过期"。
   （台账归口：见本表 §4「对编排者的不服」末段——我不动 `docs/reports/**`，这一条要编排者落 `A##`。）

### 1.4 本格未做／边界

- **未运行** r6 的 `flip-declaration.sh` 与 `faketree.sh`：两枚都把日志写进 `$here/logs`＝`probes/161/r6/logs/**`（派单 §4 明令我"只能读"）。
  我没有改路径去跑它们（那会动别人的台件），而是**自己另建同等两发**（§1.1／§1.2），证据强度不降。
- **未裁**票 158 那枚 G5 活形状（`cmd/wisp/panel_assets.go:232`）——票面红线明令它不是本票靶子，我照旧只登记"它在响"。
- `run.go:555-556` 自己写着"宿主内部直接过桥派发 ⇒ 记录在册的开放端"。它是**一枚与本 4 枚不同的残余**（正向、且结构上在本尺射程外）：
  本锚上 `git grep -nE 'DEFERRED|DEFER\(' <锚> -- cmd/wisp/run.go internal/tools/bridge.go` 只命中 `bridge.go:627 DEFERRED(C25-loop-wiring)`，
  指向的是 ticket 19 的环路接线，**不是这一枚**。⇒ 我不把它算进"真漏"，但**建议编排者核它有没有账**（我无 `docs/reports/**` 写面）。
- AC#6 这一格我的档位：§1.1／§1.2／§1.3 的**全部读数＝〔我本轮现跑过〕**；r6 自己的日志＝〔盘上有件，我未复算〕（我只抽查了它的两形设计，并另建了同等复现）。
