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
   （台账归口：见本表 §5「对编排者的不服」第 2 条——我不动 `docs/reports/**`，这一条要编排者落 `A##`／`R##`。）

### 1.4 本格未做／边界

- **未运行** r6 的 `flip-declaration.sh` 与 `faketree.sh`：两枚都把日志写进 `$here/logs`＝`probes/161/r6/logs/**`（派单 §4 明令我"只能读"）。
  我没有改路径去跑它们（那会动别人的台件），而是**自己另建同等两发**（§1.1／§1.2），证据强度不降。
- **未裁**票 158 那枚 G5 活形状（`cmd/wisp/panel_assets.go:232`）——票面红线明令它不是本票靶子，我照旧只登记"它在响"。
- `run.go:555-556` 自己写着"宿主内部直接过桥派发 ⇒ 记录在册的开放端"。它是**一枚与本 4 枚不同的残余**（正向、且结构上在本尺射程外）：
  本锚上 `git grep -nE 'DEFERRED|DEFER\(' <锚> -- cmd/wisp/run.go internal/tools/bridge.go` 只命中 `bridge.go:627 DEFERRED(C25-loop-wiring)`，
  指向的是 ticket 19 的环路接线，**不是这一枚**。⇒ 我不把它算进"真漏"，但**建议编排者核它有没有账**（我无 `docs/reports/**` 写面）。
- AC#6 这一格我的档位：§1.1／§1.2／§1.3 的**全部读数＝〔我本轮现跑过〕**；r6 自己的日志＝〔盘上有件，我未复算〕（我只抽查了它的两形设计，并另建了同等复现）。

## 2. AC#4 —— 反向判据：摘掉那对样本里"违规"那一枚，CI 会不会红

全部读数＝**〔我本轮现跑过〕**。方法：把 `tools/d22scan` 整枚模块按锚上的已跟踪清单 `cp --parents` 到 **仓外** `/tmp/wisp161v2/d22`
（`git diff --name-only <锚> -- tools/d22scan` ＝空 ⇒ 副本＝锚上字节；`sha256sum selftestsamples.go` ＝ **`7ea3f1d8f664ac61a0e3287b5fc2fc4c6a0bc94871a0fcc9005f57790d85daf5`**，
与 r4 记的"摘前＝摘后"那一枚哈希逐字相同 ⇒ 我这副本与 r4 当时那棵树同源）。**一字节生产码／一字节仓内件都没改**（删除只在 `/tmp` 的副本里）。
摘的那一味＝`selftestsamples.go:190-195`（ban #5 `mirror-hash` 名册里**唯一**的 `wantRing`；摘后 `grep -n 'mirror-hash'` 只剩 `:191 wantSilent`，`sha256` 变 `268f25ab…`）。

### 2.1 基线（同一副本，未摘）

- `go run . -self-test` ⇒ **rc=0**，末行逐字 `d22scan -self-test: clean - all 34 direction checks passed (19 expect-ring, 15 expect-silent)`。
- `go build -o d22base.exe .` ⇒ 二进制 rc=0，末行同上。
- `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ `PASS=28 FAIL=0 SKIP=6 === RUN=76`、**rc=1**。
  ⚠ 这枚 rc=1 **不是新增红**，是 `runtests.sh` 对 `SKIP` 的合法拒绝（脚本头第 2 条规矩："SKIP is NOT a pass"），
  而跳过的 6 枚全是"要在真仓里才跑"的用例（`TestCheckRootAcceptsRealRepo`／`TestScannerSelfScanOfRealRepoIsGreen`／`TestScopeReportMatchesRealCoverage`／
  `TestRealRepoBan8CoversFrontendTreeAtBan6sCount`／`TestLedgerCountsMatchAnIndependentWalk`／`TestRealRepoLedgerIsHonest`）⇒ **28＋6＝34**，与盘上／编排者复跑的 34 对得上。
  我这侧的分母与真树不同一事，**登记成"读数受副本环境影响"，不登记成缺陷**。

### 2.2 两形分开答（派单 §1 明令不许合并）

| 形 | 命令 | 摘样本后 | 逐字凭据 |
|---|---|---|---|
| **二进制那一形** | `go build` ＋ `d22rem.exe -self-test` | **rc=2**，且**一行用例都没跑**（输出里 `pinned because` 行数＝0） | `FATAL the table does not cover the tool (1 hole(s)); an unrun self-test is not a green self-test` ＋ `HOLE tag "mirror-hash" has only an expect-silent sample - a ban whose violating sample was deleted cannot be distinguished from a ban that was never implemented` |
| **CI 那一形** | `go run . -self-test`（＝`ci.yml:105` 那一步的原命令，`working-directory: tools/d22scan`） | **rc=1**，stdout 末行 `exit status 2` | 同一份 FATAL/HOLE 文本，退码却是 1 |

⇒ **本格的答案是两句话**：
（一）**CI 那一步会不会红＝会**（非零即红；我按锚现读那一步**不带 `if:`、不带 `continue-on-error`**，`sed -n '83,107p'` 全段只有 `name`／注释／`run`／`working-directory`）。
（二）**CI 看到的是 1 而不是 2** ⇒ "空心拒答(2)"与"真命中(1)"在退码上**同色**，丢的是"这次是拒答还是命中"那一格诊断，诊断只活在 stdout 文本里。
⇒ 编排者 09-27 那条"别替它们合并"的更正**我独立复现成立**（`go run` 把子进程任何非零退码压成 1，本仓的门都走 `go run`：`ci.yml:105` 与 `scripts/d22scan.sh:54`）。
（顺带：**第三种颜色**——`go test` 那一形在同一条摘过的树上 `PASS=28→25 / FAIL=0→3`，翻色三枚＝
`TestSelfTestEntryPassesEveryCase`／`TestSelfTestRosterAuditRejectsAHollowTable`／**`TestSelfTestFlagIsWiredInTheBuiltBinary`**（它自己 `go build` 再要求 rc=0，`selftest_test.go:176-205`）。
⇒ CI 里第 74 步（正控）与第 83 步（自检）**两步同色变红**，且 74 排在 83 前面。r4 只报两枚的理由我判是它的**证法**：`-overlay` 只喂给它自己的 `go test` 进程，而那枚用例内部另起一次 `go build`、不吃 overlay ⇒ 在它那一形下天然免疫。
**这一条的档位＝〔对源码的读＋我这侧的实测对照〕，我没有复跑 overlay 那一形**；我不把它算成 r4 读错，算成**它的证法在爆炸半径上少报一枚用例**。）

### 2.3 承重那一问：摘掉哪一味，会存在一发变异从此打不红？摘掉"空心就拒答"那一味呢？

变异体＝把 ban #5 的发射点 `main.go:771` 改成 `if false {`（禁令永不发火，其余一字未动）；
空心审计＝`selftest.go:355` 的 `return holes` 改成 `return nil`（一发 `sed`，不碰其它判据）。

| 发 | 名册 | 空心审计 | 变异 #5 打哑 | 结果 |
|---|---|---|---|---|
| (a) | 完整 | **在** | 是 | **rc=1**，逐字 `ban #5 mirror-hash ring FAIL … did NOT ring: no "mirror-hash" finding anywhere in the fixture …` ⇒ **变异打得红** |
| (b) | 摘掉 #5 唯一 ring | 在 | 否 | **rc=2**（CI 形 rc=1）＝审计把"摘样本"这个动作本身拦下 |
| (c) | 完整 | **摘掉** | 否 | **rc=0 clean 34 项** ⇒ 单摘这一味**今天不改变颜色**（名册完整时它是惰性的） |
| (d) | 摘掉 #5 唯一 ring | **摘掉** | **是** | **rc=0**，末行 `clean - all 33 direction checks passed (18 expect-ring, 15 expect-silent)` ⇒ **这一发变异从此打不红，实测存在** |

⇒ **答"摘掉哪一味"**：只有"空心审计"那一味能造出逃逸，且**必须与某一 tag 唯一的 `wantRing` 样本一起摘**（(c) 单独不逃逸、(b) 单独被拦）。
⇒ **答"摘掉空心就拒答那一味呢"**：它买的**不是**"防这一程把禁令改哑"——那是 ring case 自己的活，(a) 已证它打得红；
它买的是**"防下一程把样本摘掉、从而让 ring case 消失"**这一枚时序。删掉它＝把这扇门交给"没有人会去摘样本"这个假设，
而 (d) 就是那个假设破产时的读数（**门自己印 clean，一行 #5 的东西都没看过**）。
⇒ 所以本格的"承重操作定义"是：**tag 级配对的空心审计 ＋ 每 tag 的 ring case，两味互为备份，任何单独一味都不是承重**；
"摘掉任一单独一味今天不改变颜色"这件事我量了（(c)），它**不是**装饰的证据，它是备份的形状。

**AC#4 我的裁决＝〔成立〕**（票面那句"把违规样本摘掉 CI 会不会红"答得出、且我复跑；"装饰"这一支被 (d) 排除，因为**要同时删两味**才逃逸，而这两味都有用例钉着）。
⚠ 但我要把 (d) 单独递给编排者：**今天没有任何一把尺在"两味同时被摘"时响**——`go test` 那三枚里只有 `TestSelfTestRosterAuditRejectsAHollowTable` 钉审计，
若下一程连那条断言一起摘，CI 无读。这不是 AC#4 的缺陷（判据问的是"摘掉样本"，不是"摘掉判据"），是**下一程最省事的那条路**。

## 3. AC#7 —— 取证机械自己撞全仓门那一枚：归因仪器与那道顺序

全部读数＝**〔我本轮现跑过〕**（被验物＝`.scratch/wisp/probes/161/r5/attrib.sh` ＋ `ci.yml:152-172` 那一步；r4 的 `probes/161/r4/attrib.sh` 我只读不跑）。
我的两枚改件在**我自己的写面** `.scratch/wisp/probes/161/v2/`（`attrib-noknown.sh`），另两枚在仓外 `/tmp/wisp161v2/emptyrepo/probe/`；**r5 原件一字节未动**（跑它 = 只读）。

### 3.1 问①：摘掉"票面 `.scratch/wisp/issues/<NN>-*.md` 真存在"那一味，仪器会不会变成自证式恒绿？——**会，但不是静默的**

我把 `ticket_known()` 改成恒 `printf '1'`（＝只看路径形状、不查票面），其余一字未动：

- 递同一发行给**原尺**：`UNATTRIBUTABLE .scratch/wisp/probes/999/bad-sample.go (path claims ticket 999 but .scratch/wisp/issues/999-*.md does not exist …) <== REAL INJURY`，**rc=1**。
- 递给**摘掉那一味的尺**：`ticket-999 .scratch/wisp/probes/999/bad-sample.go (untracked, parse_err=0)`，**rc=0** ⇒ **任何假票号都被判"归得出"**＝票面担心的那一枚恒绿形状，实测成立。
- **但它自己当场抓住**：`sh attrib-noknown.sh --self-test` ⇒ `cases=8 failures=4`、**rc=1**（四枚 FAIL 全是该响的：`999` 两发／`777` 一发／`internal/tools/not_a_bench_path.go` 一发被读成 `ticket-` 且 rc=0；
  三枚该静的照旧 PASS，含"另一张真票 169 的路径"⇒ **不是"喂什么都响"**）。

⇒ 归因规则那**两条都要成立**这一味是承重的（删了立刻自证），且这一味**有成对样本钉着**——这就是票面要的"删掉它哪一发会重新漏"的可复算答案。
⚠ **一枚我要报给编排者的洞**：CI 里 `attrib.sh` 只出现一次（`ci.yml:172`，参数 `--tracked-only`）⇒ **`--self-test` 不在 CI**。
所以"归因规则被改坏"这一发**今天只有本地跑自测才会红**，CI 不红。票 161 AC#2 那句"把自检装成规矩（一枚 CI 入口把这些自检全部跑一遍）"对**这把尺没有做到**（对扫描器做到了：`ci.yml:105`）。⇒ 建议登记成残余，归口具名（本程无 `docs/reports/**` 写面）。

### 3.2 问②：甲形那一步有没有空心保护？——**有，且只有一行；那一行在 CI 里没有守卫**

在**仓外**另建一棵空 git 树（`/tmp/wisp161v2/emptyrepo`：`git init` ＋ `go.mod` ＋ `probe/attrib.sh`，**零枚已跟踪 `.go`**），跑 CI 的原命令：

- 原件：`attrib.sh: (A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing` ⇒ **rc=2**。
- 摘掉那一味（删 `attrib.sh:338` 的 `[ "$TRACKED_GO" -gt 0 ] || …`，其余不动）：**同一棵树、同一条命令** ⇒
  `(A) tracked-tree  lines=0 files=0  rule: MUST be empty` ＋ `GREEN (tracked-only) - (A) empty over 0 tracked .go files` ⇒ **rc=0**。

⇒ 甲形那一步**不是第二枚恒绿**（派单 §1 明令的那一发我亲手量了）；但也⇒ **它的空心保护只有这一行，而 `--self-test` 的 8 发全是 `classify_line` 的样本、没有一发钉这条守卫**
⇒ 删掉它 CI **不红**。这一枚与 3.1 的洞同因（自测没进 CI），**修法同一发**：给 `attrib.sh --self-test` 加一发"递 0 枚文件必须 rc=2"的合成样本并把自测接进 CI（**我不做，也不建议顺手做**——它动 `ci.yml`，是票 134 定过的地界）。

### 3.3 问③：那句"说得出删掉它哪一发会重新漏"——我亲手删了两发

上面 3.1（删票面存在性检查）与 3.2（删空心守卫）**都是我自己执行的删除**，两发各自当场变色：前者让假票号静、自己的 `--self-test` 4/8 红；后者让零输入变 GREEN、rc 从 2 掉到 0。
⇒ 这句**答得出，且答的是读数不是形容词**。顺序那一半（`gofumpt -l .`／`d22scan.sh`／`go vet` 必须在台件全部入库之后跑最后一次）：
在 CI 里由"CI 跑的是检出树"这一事实**天然成立**（甲形那一步排在 `gofmt (gofumpt)` **之后**、`git show <锚>:ci.yml` 现读第 136→152 行，位置＝r5 自陈"就坐在既有那步之后，没往上顶"，与锚一致）；
本地那一半（r5 `logs/e1..e4`"台件全部入库后再跑一次"）＝**〔日志＋归档，我未复算〕**。

### 3.4 那枚已知常红会不会吃掉后面步骤的读数——**不会（按作业答）**

`internal/panel/tokens_fourway_test.go` 属 `internal/panel`，本锚上跑它的只有 `scripts/portable-tests.sh`（现读 `:141`／`:179` 两处点名），
而那枚脚本被 `test-core`／`test-windows` 作业调（`ci.yml:341`／`:511`）⇒ **与 `lint` 作业（我全部 AC#4／AC#7 读数所在地）不是同一枚作业**，
GitHub 里作业互相不阻塞，所以它吃不到我的读数。
⚠ 同一作业**内部**确实会吃：`lint` 的六步没有 `if:` ⇒ 上面一步红，后面不跑（"谁先红"的因果就是 ci.yml 注释自己写的那件事）——
这也是为什么 §2.2 我把"74 步与 83 步同时红、且 74 在前"记下来。
**我没跑它、没修它、没 Skip 它、没动任何断言**（派单 §4）；"它今天到底红不红"这一条我的档位是〔盘上有件／票面与台账点名，**我未复算**〕。

## 4. 我没测什么／哪些读数受并程影响

- **并程影响（此刻另一程在写 `internal/tools/**`）**：本表**没有一处**读工作树的 `internal/**`／`cmd/**`——生产码读数全走 `git show <锚>:<path>` 与 `git grep <锚>`；
  仪器读数是 `git diff --name-only <锚>` 先证明"盘上＝锚上同字节"再跑的；`go test`／`go run` 全在 `/tmp` 副本里。
  唯一的并程痕迹＝`runtests.sh` 那发的 `PASS=28/SKIP=6`（副本不是 git 检出 ⇒ 六枚真仓用例跳过），已在 §2.1 记账，**不影响任何一条裁决**。
- **没测**：CI runner 一次未跑（我只在本地跑 CI 步里的原命令）；r4 的 `-overlay` 证法未复跑（§2.2 那句"为什么少报一枚"是源码读＋我这侧对照）；
  r5 的 `tracked-dirty-proof.sh`／`logs/e1..e4` 未复算；r6 的 `faketree.sh`／`flip-declaration.sh` **未运行**（它们写 `probes/161/r6/logs/**`，在我的写面外；我另建了同等两发）；
  `internal/panel/tokens_fourway_test.go` 未跑（见 §3.4）；票 158 那枚 G5 活形状（`cmd/wisp/panel_assets.go:232`）未裁（票面红线）。
- **删除命令**：全程 `rm -rf`／`rm`／`rmdir`／`del`／`git rm`／`clean`／`--amend`／`reset`／`rebase`／`stash`／`checkout .` **零次**。
  （唯一被执行过的删除是 `selftest.go` 里 `runSelfCase` 对自己 `os.MkdirTemp` 的 `defer os.RemoveAll`，由派单指定的命令带出，在 `/tmp` 内、不是我发起。）
  仓外副本（`/tmp/wisp161v2/**`）与其中的三棵合成树按"只建不删"留着等清点。

## 5. 对编排者的不服（三格裁完，按派单 §8 第 6 项）

1. **你那两枚 `>`（09:1x／10:0x）我判：09:1x 那枚成立且是这一格的骨架，10:0x 那枚**方向对但落点偏了**。**
   10:0x 说 r6"把未裁的问题钉成事实"。我裁完的读数是：**4 枚里没有一枚是真漏**（§1.3），所以它钉住的不是"漏"这个事实；
   可它确实钉歪了另一样——**G7 的 `ring` 今天为真的理由是两枚注释 token ＋一枚函数签名**，
   而台件注释（`gate-clauses.sh:288-291`）把这三枚写成了"反向 3 枚＝retention／goroutine／shutdown"，读者会以为那是三处待判的资源。**没裁的问题是"是不是真漏"，更窄的问题是"这三枚里有多少是代码"** ——后者才是这一格该红而没人量的地方。
2. **"那 4 枚"这个数我接了，但它的构成我改写**：`4 = 1 枚跨文件分工（真读数） + 1 枚签名/字符串（词表缺一味：`Go` 不并入方） + 2 枚尾随注释（根本不是代码）`。
   ⇒ 派单 §2 第一支（"任何一枚真漏 ⇒ 单立账"）**不触发**；但我建议另登一枚**仪器精度**账（`pair()` 的剔注释只剔整行 ⇒ 同行尾随注释算调用点），
   归口＝`probes/154/gate-clauses.sh` 的尺本体（票 154/158 那条线），**不是**票 161、更不是票 169。
3. **"这轮不进 CI"我判：对聚合那把尺成立；但你顺带把归因尺的自测留在了 CI 外，那一枚是漏的。**
   §3.1／§3.2 两发我实测：`attrib.sh` 的两味承重保护（票面存在性检查／空心守卫）**被删之后 CI 都不红**，因为 CI 只跑 `--tracked-only`、从不跑 `--self-test`。
   这与"进 CI 要先答哪一腿的红该阻塞推送"（你的理由）不是同一件事——`--self-test` 是**这把尺自己的成对样本**，
   它进 CI 不需要新的阻塞语义，只要"自测红＝这台件坏了＝红"。票面 AC#2 对扫描器正是这么做的（`ci.yml:105`）。⇒ 这条我不推翻你的裁定，我推翻它的**适用范围**。
4. **你那句"像把未裁的问题钉成事实"里的"恒响也算装饰"我不用**：恒响至少是**会响**；本格的形状反过来——**声明会响、但响的原因是注释**。
   ⇒ 我的判据写法（可复算）：`sh probes/154/gate-clauses.sh <锚>` 取 `## G7 主尺` 段的每一枚 `UNPAIRED-REV`，再对每枚 `git show <锚>:<file> | sed -n '<行>p'`，
   **肉眼分"代码里的词／注释里的词"**；分不出来时不许登记"该响"，要登记"该响，因为 N 枚里有 M 枚是注释"。

## 6. 三格档位（只许派单 §8 那四档）

| 格 | 裁决 | 档位 |
|---|---|---|
| AC#4 反向判据（两形＋承重那一问） | **成立**（两形分开：CI 形 rc=1／二进制形 rc=2；(d) 一发证明"摘样本＋摘审计"能逃逸 ⇒ 两味互为备份） | 〔我本轮现跑过〕（overlay 那一形的对照＝〔源码读＋我这侧实测，未复跑 overlay〕） |
| AC#6① 聚合退码 | **成立**（两向翻声明各打红一发、两枚同时翻⇒rc=2、还原⇒0 且名册只差生成时刻一行） | 〔我本轮现跑过〕 |
| AC#6② 两族成对普查 | **成立，但带一枚仪器精度残余**（合成树三向我自己另建并复跑：phase1 两族正向各响 1 枚／phase2 回到 0 且 `git grep rc=0`／phase3 删样本后 0 且 `git grep rc=1`）；那 4 枚反向**逐枚判为形状使然、零枚真漏** | 〔我本轮现跑过〕（r6 自己的两份日志＝〔盘上有件，我未复算〕） |
| AC#7①②③ | **成立**（三问各亲手删／复跑一发：删票面存在性⇒假票号变静＋`--self-test` 4/8 红；删空心守卫⇒零输入 GREEN rc=0 vs 原件 rc=2；顺序那半在 CI 天然成立） | 〔我本轮现跑过〕（r5 `logs/e1..e4`＝〔日志＋归档，我未复算〕） |

⇒ **票面六格里本表判成立的三格（AC#4／AC#6／AC#7）我一格都没勾**（派单 §4）；勾与不勾是编排者的动作，翻勾凭据＝本表 §2／§3／§1 各节。

`next=` 交回编排者，按优先级：
① 三格翻勾凭据已给（本表 §2／§1／§3）；
② 建议登两枚残余（我无 `docs/reports/**` 写面）：**R-a＝`attrib.sh --self-test` 不在 CI ⇒ 归因规则与空心守卫删了 CI 不红**（§3.1/§3.2）；
**R-b＝`pair()` 只剔整行注释 ⇒ 同行尾随注释被算成调用点，G7 今天 3 枚反向里 2 枚是注释**（§1.3，归口票 154/158 的尺本体）；
③ 若要接聚合门进 CI（撤销口令"把 161 的聚合门接进 CI"），**先修 §1.3 那条"枚数不进判据"**——否则新长的未成对枚数躲在 `ring` 声明里，聚合永远 `ok`；
④ `probes/161/v2/**` 里我的两枚改件（`attrib-noknown.sh` 等）是**故意改坏的仪器**，别当仪器用；
⑤ 未做：`-overlay` 那一形的对照复跑、`internal/panel` 那枚常红的现跑（两枚都属"不必为本格裁决而跑"）。
