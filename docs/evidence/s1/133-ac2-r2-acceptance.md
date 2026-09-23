# 票 133 AC#2 第二格复算 —— 非实现者验收方的第二次裁决（同一格第二次退回即终止）

**执行方**：`acceptor-ticket133-ac2-r2`（验收位，非那枚尺的作者，也非回修方）
**被验的说话**：回修方 `worker-ticket133-ac2-fix` 在票面末尾那 52 行里的那一句 —— "它把退回件的五条账修完了"
**判据物**：票面 AC#2 那一格 ＋ `docs/evidence/s1/133-adversarial-acceptance.md` §3.1／§6／§8
**回修方证据**：`docs/evidence/s1/133-ac2-fix.md`
**本表只裁 AC#2**。AC#1／AC#3／AC#4／AC#5／AC#6 一律未裁。**本表不翻任何勾**（AC#2 由编排者按本表翻）。

## 0. 锚点与树

### 0.1 锚点（开工第一步自己量）

```
$ git rev-parse --short HEAD            048a9e4
$ git rev-parse --abbrev-ref HEAD       dev
$ git show -s --format="%H %ad" --date=iso 048a9e4  048a9e4c0239aaa80722abc20841bad8a395ade3 2026-09-23 22:58:09 +0800
$ git merge-base --is-ancestor fbf420c HEAD && echo yes      yes
$ date -u                               Wed Sep 23 15:00:09 UTC 2026  （开工；+8 ＝ 23:00）
```

派单写"HEAD 约 `048a9e4`"，本程自量也是 `048a9e4`。回修方 `next=` 要求"在 `fbf420c`（或其后代）上复跑"，
`fbf420c` 是 `048a9e4` 的祖先 ⇒ 条件成立。

### 0.2 锚点之后的提交：逐枚判断有没有动 `cmd/wisp/**`

```
$ git log --oneline 048a9e4 ^fbf420c      共 11 枚，全部列出（见下）
048a9e4 docs(122,AC#7 锚点精度)          docs/evidence + .scratch 票面
6e027e4 docs(133,AC#2 回修收尾落盘)      .scratch 票面 133 + docs/evidence/s1/133-ac2-fix.md
21c8def docs(133,AC#2 回修证据 §5/§3.6)  docs/evidence/s1/133-ac2-fix.md
2f291d0 test(136,AC#9①)                  internal/observe/sampler_settle_zerosample_136_test.go
595abd3 docs(136,票面 log)               .scratch 票面 136
36443f2 docs(136,AC#8②③ 实现方证据)      docs/evidence + .scratch 票面 136
ca2b34a docs(133,AC#2 回修证据 §3.3/§3.4) docs/evidence/s1/133-ac2-fix.md
79ddd49 test(136,AC#8②)                  internal/observe/sampler_test.go
1d38206 docs(A134,A135;122,136 翻勾)      .scratch 票面 ×3 + docs/evidence ×2 + docs/reports
7bdfbbb docs(133,AC#2 回修证据 §2.5/…)    docs/evidence/s1/133-ac2-fix.md
ae6a01e docs(136,AC#1 验收方 三处措辞)    docs/evidence/s1/136-ac1-adversarial-acceptance.md
```

机器核过，不靠上面的目测：

```
$ git diff --name-only fbf420c..048a9e4 | grep -c "^cmd/wisp/"     0
$ git diff --stat  fbf420c..048a9e4 -- cmd/wisp/                   （空输出）
```

⇒ **`cmd/wisp/**` 自 `fbf420c` 起一字未动**。本表所有读数取于 `048a9e4` 的归档树，
其中被测那枚尺与 `fbf420c` 逐字节同（见 §0.3），回修方在 `fbf420c` 上取的 §3.3／§3.6 读数
因此与本程同树可比。锚点上唯一未提交面是 `docs/evidence/s1/136-ac8-ac9-impl.md`（兄弟在飞）
与一枚来源未明的未跟踪件（§6 末），**都不在 `cmd/wisp/**`，本程未读、未改、未提交**。

### 0.3 快照树（只在仓外建，**只建不删**；仓内零 worktree／零 checkout）

```
$ mkdir -p /d/tmp/wisp133-r2-tree1 && git archive 048a9e4 | tar -x -C /d/tmp/wisp133-r2-tree1
$ git show 048a9e4:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ sha1sum /d/tmp/wisp133-r2-tree1/cmd/wisp/leg_dispatch_gate_133_test.go 906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ ls -l cmd/wisp/ | grep -iE "131x|probe|\.off|\.bak"                （空 ⇒ 归档里没有种件残留）
$ cp <repo>/third_party/sherpa-onnx/*.dll <树>/third_party/sherpa-onnx/
    onnxruntime.dll 17799168 / sherpa-onnx-c-api.dll 4605952 / sherpa-onnx-cxx-api.dll 259584（缺即 0xc0000135）
$ go build ./cmd/wisp/   rc=0
```

`906201f4a3995d10b0a65910aa4cc68e1f745a9d` 与回修方 §3.3 自报的 `ARCHIVE-MODULE` 行**逐字同**
⇒ 本程量的就是它交回的那一枚尺。

**每发读数一棵新树**：本程不复用回修方／前任验收方的树，`/d/tmp/wisp133-r2-run.sh` 在每次读数前
`cp -r tree1 wisp133-r2-T<n>`，树已存在则拒跑（`TREE-EXISTS-REFUSING`）。
这样做是为了绕开回修方 §3.2 自己踩到的那个坑 —— `probe.py` 的 `clean()` 会把上一轮留下的
`leg_dispatch_gate_133_test.go.p1.bak` 盖回当前树，让"改后的树"读出"改前的读数"。
本程每发的树都是干净归档的副本，**没有可被盖回的 `.bak`**，所以那类作废读数在本程结构上不可能出现。

### 0.4 台件（复用的只有"种法定义"，没有复用任何人的读数）

前任验收方的两枚台件与本程逐字节同（本程自己 `sha1sum` 过，不是照抄回修方的自述）：

```
$ sha1sum /d/tmp/wisp133-acc-r1-mutate.py /d/tmp/wisp133-ac2fix-mutate.py
  cd88168c05cc14db3410913524ec97992029dc73  （两枚相同）
$ sha1sum /d/tmp/wisp133-acc-r1-probe.py   /d/tmp/wisp133-ac2fix-probe.py
  c44704ab47a89e66325407c4c8104932053e5e63  （两枚相同）
```

⇒ 回修方 §0 那句"复跑台是验收方台件的逐字节副本"**成立**，本程 `p1`／`p2`／`p3`／`p4`／`p5` 与原
裁决表 §3.1 是**同一发件**。本程自己的副本：`/d/tmp/wisp133-r2-{mutate,probe}.py`（`sha1sum` 同上）。
`p6` 是回修方自造的（`/d/tmp/wisp133-ac2fix-p6.sh`），本程按它交接段里写死的形状
**另写一份** `/d/tmp/wisp133-r2-p6.sh` 复种（§1.6 判它造得对不对）。
`p7`／`p8` 是本程自造（§1.7／§1.8），**前任验收方与回修方都没造过** ⇒ 按本仓规矩不记成别人已测。

日志与逐发红名全在 `/d/tmp/wisp133-r2-out/`；树 `/d/tmp/wisp133-r2-T01…T12`（＋后续），**保留不删**。

### 0.5 本程的硬约束（照派单，写给下一位读者）

- `cmd/wisp/**` 只读；被测那枚尺一个字没改（§1 每发的落地证明都是"改树"，改的是快照副本）。
- 在飞兄弟不碰：`worker-ticket136-ac8-ac9` 的 `internal/observe/**`、`worker-ticket137-ac1` 的
  `internal/winsec/**` —— 本程**未读它们未提交的半成品**，`internal/observe/**` 里 `2f291d0`／`79ddd49`
  两枚已提交的件本程也不需要（不在 AC#2 判据物上）。
- 禁改清单（`docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／阈值／golden／`frontend/**`…）零接触。
- **不跑任何计时类断言**（本机另有代理在飞、`slo-full` 会随 push 自启抢 CPU）：本程所有读数都是
  四数／红名／rc／字节数，无 D32 那两枚数。
- 凭据值零接触：本表只可能出现变量名。
## 1. 逐发复算（先证落地，再取读数；红名一律是本程自己 `grep` 出来的）

**通则**：下面每一发都是**一棵新树**（`cp -r tree1 T<n>`，树名已存在即拒跑），先发 `LANDED`/`PROBE applied`
行证种件落地，再 `go build ./cmd/wisp/` 与 `go build ./...` 各取一次 rc，才允许读红绿。
本程**不沿用回修方与前任验收方的任何红名**——每发的 `--- FAIL:` 行都是本程日志里 `grep -E '^--- FAIL'` 打出来的。
时间戳：`date -u` 15:08:19Z–15:17:16Z（本地 23:08–23:17 +08）为第 1 批，23:1x–23:3x 为第 2 批（逐节标）。

### 1.0 判据五枚修法在本程锚点树上的字节（"改后"是哪一形，先钉死）

```
cmd/wisp/leg_dispatch_gate_133_test.go:398   case isTest && isRunnableCase133(fd, testingPkg):
cmd/wisp/leg_dispatch_gate_133_test.go:412       p.helpers[fd.Name.Name] = append(p.helpers[fd.Name.Name], decl)
cmd/wisp/leg_dispatch_gate_133_test.go:453   func isRunnableCase133(fd *ast.FuncDecl, testingPkg string) bool {
cmd/wisp/leg_dispatch_gate_133_test.go:1100  func (p *pkg133) caseLabelLegs133(cc *ast.CaseClause) (keys, aliases, reds []string) {
cmd/wisp/leg_dispatch_gate_133_test.go:1173  func (p *pkg133) classifyCond133(e ast.Expr) (keys, aliases []string, kind condKind133) {
cmd/wisp/leg_dispatch_gate_133_test.go:1220  func splitCommandLabels133(values []string) (keys, aliases []string) {
cmd/wisp/leg_dispatch_gate_133_test.go:1348  adjacent := p.rulingAdjacentTo133(leg)
cmd/wisp/leg_dispatch_gate_133_test.go:1390  reds = append(reds, p.misplacedRulingReds133(legs)...)
cmd/wisp/leg_dispatch_gate_133_test.go:1542  func (p *pkg133) rulingAdjacentTo133(leg *leg133) string {
cmd/wisp/leg_dispatch_gate_133_test.go:1556  func (p *pkg133) legOwnsFile133(leg *leg133, file string) bool {
cmd/wisp/leg_dispatch_gate_133_test.go:1574  func (p *pkg133) misplacedRulingReds133(legs []*leg133) []string {
```

⇒ 回修方 §2.1–§2.5 声称的五枚修法**都在被验那枚文件里**（字节在场 ≠ 判据有效，有效性看 §1.1–§1.9）。

### 1.1 独立基线（未变异，`048a9e4` 归档树）

| 发 | 树 | rc | RUN | 顶层 PASS | FAIL | SKIP | 红名 | SKIP 行 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `R2-B0-base` `-count=1 -v` | T01 | **1** | 101 | 53 | **1** | 0 | `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` | 0 |
| `R2-B0-c2` `-count=2 -v` | T11 | 1 | 202 | 107 | 1 | 0 | 同上那一枚 | 0 |
| `R2-B0-nogate` `-count=1 -v -skip '^TestAC1AC2DispatchHopGate133$'` | T12 | 1 | 100 | 52 | 1 | 0 | 同上那一枚 | 0 |

⚠ **一处与本程被验方读数不一致，两边都留、不改数**：回修方 §0.1／§4 的基线是 101/54/0/0 与 100/53/0/0，
本程在锚点树上量到的是**同一枚既有 flake 红了三次**（`resident_sink_nail_127_windows_test.go:549,572,583`：
常驻腿子进程退出码 3221225786 = `0xc000013a` STATUS_CONTROL_C_EXIT，sink 目录里"1 file(s), 0 record(s)"）。
归因线索（本程量的，不是猜的）：
`--- PASS` 的那一发 4.58 s、`--- FAIL` 的两发 5.13 s／6.72 s；本机此刻有兄弟代理在跑整包
（`worker-ticket136-ac8-ac9` 在 `internal/observe/**`），而本票那枚尺 `TestAC1AC2DispatchHopGate133`
在三发读数里**都是 PASS**（`R2-B0-base` 的 `GATE-ROWS: run=1 pass=1`、count=2 里 `run=2 pass=2`）。
⇒ 本程**不把它当 AC#2 的破口**（红的是票 127 的常驻钉，与判据 (d) 无关），但**也不许拿"我这发绿"当它的读数**：
所以 §1.2–§1.9 的每一发都改判**红名集合与基线的差集**（§1.10），四数只作旁证。第 2 批里本程另取一发
`R2-B1-base`／`R2-B1-c2` 复跑同一命令，看这一形稳不稳（见 §1.11）。

**名册**：`R2-B0-base` 顶层 54 枚、`R2-B0-nogate` 顶层 53 枚；本程自己 `comm`：
新增只有 `TestAC1AC2DispatchHopGate133` 一枚、消失 0 枚、共有 53 枚一字未动
（`/d/tmp/wisp133-r2-out/R2-B0-base.roster` 54 行、`R2-B0-nogate.roster` 53 行）。
⇒ 回修方 §4 的逐名账、前任 §7.1 的 101/54→100/53 在本程独立复现（"摘掉本尺少的那一枚就是本尺"这一问的答案仍是 `TestAC1AC2DispatchHopGate133`）。

### 1.2 `p4` 对照（诚实的那一发：只加给人看的 usage 行）

```
树 T02：shot x14b2c + probe p4
LANDED-MAIN x14b2c:  main.go:80: 		case "sfx131":
PROBE p4 applied
BUILD-cmdwisp R2-p4 rc=0 / BUILD-all rc=0
READ label=R2-p4 rc=1 RUN=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing: …
    dispatch ledger … (12 legs, 4 claims in this gate's registry)
    leg sfx131 main.go:81 installs=false handoff=false covered=RED nothing entries=cmdSfx131 aliases=
```

⇒ **与回修方 §3.6 一致（红）**，与前任 §3.1 的"对照发本来就红"一致。12 legs ＝ 基线 11 ＋ 被分发的 `sfx131`。

### 1.3 `p1` 同名方法（`R-133-1` 的那一发）

```
树 T03：shot x14b2c + probe p1
PROBE p1 applied
  probe133_test.go:10: func (phantomRecv133) TestSfx131LegIsDriven() { _ = cmdSfx131(nil) }
BUILD-cmdwisp R2-p1 rc=0 / BUILD-all rc=0
READ label=R2-p1 rc=1 RUN=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing …
    ledger (12 legs, 4 claims) / leg sfx131 … covered=RED nothing
```

⇒ **改前绿（前任 §3.1、回修方 §1 两次阅读都是 rc=0）⇒ 本程改后红**，红名点到本尺。与回修方 §3.6 一致。
"这枚方法到底跑没跑"另有整包一发 `R2-p1-full`（§1.9）。

### 1.4 `p2` 不相邻裁决句（`R-133-3`）

```
树 T04：shot x14b2c + probe p2
PROBE p2 applied
  doctor.go:341: // WISP-LEG-COVERAGE-RULING: sfx131 acceptor probe: planted in an unrelated production file, far from any code that owns this leg.
BUILD-cmdwisp R2-p2 rc=0 / BUILD-all rc=0
READ label=R2-p2 rc=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: doctor.go:341 carries WISP-LEG-COVERAGE-RULING: for leg "sfx131", but doctor.go owns
      neither this leg's dispatch (main.go:81) nor any symbol it calls (main.go, sfx131x.go). …
    AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing …
```

⇒ **两枚红，改前绿 ⇒ 改后红**，与回修方 §3.2／§3.6 一致。
**正向控制本程另取一发**（§1.5）：裁决句改到这条腿自己的文件里仍须绿，否则"收紧"就是把两形一起打死。

### 1.5 裁决句写在 owning 文件里的正向控制（本程自取，回修方给过同形读数）

回修方 §3.2 的 `POST-p2b`：把同一行裁决句从 `doctor.go` 末尾移进 `sfx131x.go:22` ⇒ rc=0 绿、
账本 `covered=ruling sfx131x.go:22`。本程在 §1.8 的 `p8` 里取了**更强的同族正向形**
（裁决句坐在 `main.go`，即这条腿 dispatch 所在的那枚文件，且**不点任何 entry 符号名**），
读数是绿 ⇒ "同文件核对"确实在放行"写在该写地方的裁决句"，收紧没有把合法形状打死；
代价是 `R-133-3` 那半句"点到 entry 名"没有判据兜住（§2.2 独立重走）。

### 1.6 `p3` 钉表收 helper 名（`R-133-2`）

```
树 T05：shot x14b2c + probe p3
PROBE p3 applied
  leg_dispatch_gate_133_test.go:140: 	{leg: "sfx131", test: "phantomHelper133", entry: "cmdSfx131"},
BUILD-cmdwisp R2-p3 rc=0 / BUILD-all rc=0
READ label=R2-p3 rc=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: leg "sfx131" (main.go:81) is dispatched by func main and covered by nothing …
    AC#1 RED: this gate claims leg "sfx131" is nailed by "phantomHelper133", which is not
      a test function in this directory's sources: declares probe133_test.go:6, which this file
      sorts with isRunnableCase133: a helper or a method, not a case Go's testing package runs. …
```

⇒ **两枚红，改前绿 ⇒ 改后红**，与回修方 §3.1／§3.6 一致；**红名文案与判据这一发是对得上的**
（它说的是"不是用例"，查的确实是收紧后的桶）。⚠ 同一枚文案的**另一半**（"A renamed or build-tag-hidden
case has to be red somewhere, and here is where."，`:1338`）本程**复现出反例**，见 §5 的 `R-133-9`。

### 1.7 `p5` 同 case 多标签（`R-133-4`）

```
树 T06：shot none（pristine main.go 上手种，与前任 extra.sh、回修方 §3.2 同法）+ probe p5
PROBE p5 applied
  main.go:80: 		case "run", "runalt133":
BUILD-cmdwisp R2-p5 rc=0 / BUILD-all rc=0
READ label=R2-p5 rc=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: leg "runalt133" (main.go:80) reaches installLogSink on this path: dispatch -> runTextTask -> installLogSink
    AC#1/#2 RED: leg "runalt133" is dispatched by func main and documented in no line of the usage block …
    ledger (12 legs, 4 claims)：
      leg run        main.go:80  installs=true  covered=nail TestAC2SealNoticeLandsInTheRunLegLogFile -> runTextTask
      leg runalt133  main.go:80  installs=true  covered=RED sink with no nail
```

⇒ **改前"账本仍 11 legs、`runalt133` 零行、判绿" ⇒ 本程改后 12 legs、第二条标签有自己的行、两枚红**，
与回修方 §3.2／§3.6 一致。别名折叠的正向控制本程在 §1.1 基线账本里直接读到：
`leg version … aliases=--version|-v`、`leg help … aliases=--help|-h`、`11 legs` ⇒ 没有把 CLI 惯例别名当成第二条命令。

### 1.8 `p6` 复合 argv 条件（`R-133-5`，**回修方自造**的那一发）

先答派单那一问：**这一发造得对不对、防的是不是 `R-133-5` 那件事本身。**

- 前任 §8 给 `R-133-5` 登记的字节是：`classifyCond133` 的字面量分支 `key, found = v, true` **没有** `!found` 守卫
  （last-match），而 no-args 分支带 `!found`（first-match）⇒ 病名是"复合条件里只留一枚命令名"。
- 回修方的种件是在 `switch args[0] {` 之前插
  `if len(args) > 0 && (args[0] == "x133a" || args[0] == "version") { attachParentConsole(); printVersions(""); return }`，
  **并把 `x133a` 故意留在 usage 块之外**。⇒ 形状正是"同一枚复合条件里两枚字面量"：
  第一枚 `x133a` 是被 last-match 吃掉的那一枚（若判据松，它整条不进账本），
  第二枚 `version` 是活下来的那一枚；`x133a` 不在 usage ⇒ 一旦它出账本，双向差**必然**撞上，
  于是"没进账本"与"进了但没人管"两形可分。本程核它种得对不对，用的是**改前的反证**：
  同一枚种件在本程树上会让账本把 `version` 那枚腿从 `main.go:79` 那一行也报出来，
  而 `x133a` **整份读数零出现** ⇒ 病名与字节对上。
  ⚠ 本程未复跑"改前"那一发（那要在**没有第 5 枚修法**的树上跑，等价于回退被验版本）：
  §4 诚实列进"我没做的档"。本程判它"造得对"的**独立**凭据是下面这条改后读数里
  **`x133a` 与 `version` 各占一枚腿**这一事实本身。
- **它防的确实是 `R-133-5` 那件事本身**，不是别的东西：判据面是"每一枚可归类的 argv 比较各出一条腿"，
  与本程读到的 `classifyCond133`→`splitCommandLabels133`（`:1173`／`:1220`）一致。

```
树 T07：shot none + probe p6（本程自己另写的 /d/tmp/wisp133-r2-p6.sh，形状按交接段写死的种件）
R2-P6 applied (compound early if; x133a deliberately NOT in the usage block)
  main.go:79: 	if len(args) > 0 && (args[0] == "x133a" || args[0] == "version") {
BUILD-cmdwisp R2-p6 rc=0 / BUILD-all rc=0
READ label=R2-p6 rc=1 RUN=1 TOPPASS=0 TOPFAIL=1 TOPSKIP=0
REDNAMES:    --- FAIL: TestAC1AC2DispatchHopGate133
    AC#1 RED: leg "x133a" (main.go:79) is dispatched by func main and covered by nothing …
    AC#1/#2 RED: leg "x133a" is dispatched by func main and documented in no line of the usage block …
    ledger (13 legs, 4 claims)：
      leg x133a  main.go:79  installs=false handoff=false covered=RED nothing entries=attachParentConsole|printVersions
X133A 名字在整份读数里出现 3 次（回修方 §3.5 的 NAME-APPEARANCES=3，本程同数）
```

⇒ **本程改后红、13 legs、`x133a` 有自己的行**，与回修方 §3.5／§3.6 逐字同形。
⚠ **但这一发的"改前"是回修方自证的**：按派单那条规矩（验收方没造的支别记成已测），
`R-133-5` 从"未造的读码断言"升成"已证"这件事，**证据链只有一端**（回修方的 `T7` 那一发 12 legs、`x133a` 零出现）。
本程没有第二端 ⇒ 判"修法有效已独立复现（改后红）"、"三态只复现两态"。

