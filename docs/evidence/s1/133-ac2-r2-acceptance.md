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
`R2-B1-base`／`R2-B1-c2` 复跑同一命令（`date -u` 15:41:14Z／15:43:01Z，本地 23:41／23:43 +08），
**两发都干净**：

```
READ label=R2-B1-base rc=0 RUN=101 TOPPASS=54 TOPFAIL=0 TOPSKIP=0     （T16，同 §1.1 第一发命令）
READ label=R2-B1-c2   rc=0 RUN=202 TOPPASS=108 TOPFAIL=0 TOPSKIP=0    （T17，同 §1.1 第二发命令）
```

⇒ **两边都留、不改数**：回修方 §0.1／§4 的 101/54/0/0 与 202/108/0/0 本程**复现成功**（第 2 批），
第 1 批那三发红是同一枚既有 flake 在负载下的命中（四发同命令里 3 红 1 绿）。
这一枚 flake **不属 AC#2 判据物**、也不在本表任何一发的红名里（除 §1.9 形二的 `R2-RM-n3` 那发它和 131 的门一起红），
但它足够说明一件事：**在这台机器上"四数"不是稳定量** ⇒ 本表一律以红名集合与名册差集为准（§1.10）。

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

### 1.9 五发两拍那一族（六发 × 两形，本程逐发自建树复跑）

**编号以文件里真实写着的为准**：本程用的是回修方 §3.3／§3.4 表里那六个标签
`N-3`／`X4`／`X8`／`X12`／`X14 一拍`（台件里叫 `x14c`，即"装了听众没钉"的第一拍）／
`X14 二拍`（台件里叫 `x14b2c`，即"拆掉 install 之后仍零红"那一拍）。
种法＝`/d/tmp/wisp133-r2-mutate.py`（与前任验收方台件逐字节同，§0.4）。每发都先看 `LANDED-MAIN` 行
＋`go build ./cmd/wisp/` 与 `go build ./...` 各一次 rc。

**形一：门关着（`-skip '^TestAC4EveryLegIsNailedOrRuled$'`，131 的两枚文件一个字没动）** —— 这是 AC#1 归属硬线要的那一格："不许蹭 131 的门"。

| 发 | build | 门关着读数 | 红名（本程 `grep '^--- FAIL'`） | 机检：131 门提及／`--- SKIP` 行／本尺 RUN／本尺 FAIL／`build failed`／`panic` |
| --- | --- | --- | --- | --- |
| N-3 | rc=0 | rc=1 100/52/1/0 | `TestAC1AC2DispatchHopGate133` | 0／0／1／1／0／0 |
| X4 | rc=0 | rc=1 100/52/1/0 | 同上 | 0／0／1／1／0／0 |
| X8 | rc=0 | rc=1 100/52/1/0 | 同上 | 0／0／1／1／0／0 |
| X12 | rc=0 | rc=1 100/52/1/0 | 同上 | 0／0／1／1／0／0 |
| X14 一拍（`x14c`） | rc=0 | rc=1 100/52/1/0 | 同上 | 0／0／1／1／0／0 |
| X14 二拍（`x14b2c`） | rc=0 | rc=1 100/52/1/0 | 同上 | 0／0／1／1／0／0 |

⇒ **六发两拍零退化，且每发红名都点到本尺**，与回修方 §3.3 表**逐格一致**（六栏机检本程自己算的）。
"门开着"那一形本程不必另取：`§1.1` 基线里 131 的门本来就 PASS（`door131=2` 是它自己的账本行，红名 0 枚）。

**形二：摘掉本尺（`-skip '^TestAC1AC2DispatchHopGate133$'`，131 的门照跑）** —— 这一形问的就是 `R-133-6`
那句被推翻的全称断言："把本尺摘掉，还有没有用例会红"。

| 发 | 摘掉本尺的读数 | 红的那一枚是谁 |
| --- | --- | --- |
| N-3 | rc=1 100/51/**2**/0 | `TestAC4EveryLegIsNailedOrRuled`（131 的门）**＋ §1.1 那枚既有 flake** |
| X4 | rc=1 100/52/1/0 | `TestAC4EveryLegIsNailedOrRuled` |
| X8 | rc=1 100/52/1/0 | 同上 |
| X12 | rc=1 100/52/1/0 | 同上 |
| X14 一拍 | **rc=0 100/53/0/0** | **无人**（全树零红） |
| X14 二拍 | **rc=0 100/53/0/0** | **无人**（全树零红） |

⇒ 本程独立量到的这一族与回修方 §3.4 **同命题**：**"摘掉本尺五发零红"是假的**，四发仍由 131 的门红、
只有 X14 两拍零红 ⇒ `R-133-6` 要求的那句更正**有本程自己的第一手凭据**（不只是它自述）。
`R-133-6` 是**证据面**的账，票面"翻格条件"那一行把它列成必办之一 ⇒ 本程判**已办**：
`docs/evidence/s1/133-ac1-ac2-instrument.md` §3 的更正段在场、原文未抹（本程 `git diff --numstat` 复算见 §3.2）。

**反向判据（上一轮立下的那条："同一枚事实的两种写法不许互相抵账"）**：X14 那一发本程**两拍都复跑**了——
第一拍（`x14c`，字段初值里直呼 `installLogSink`、装了听众没钉）门关着红、
第二拍（`x14b2c`，把那截 init 与调用整个删掉）**也必须红**，实测门关着仍 rc=1、红名仍是本尺。
⇒ 二拍没有被"同案的另一种写法已经红过"抵掉。

### 1.10 红名集合与基线的差集（四数之外的必比项）

基线取 §1.1 的 `R2-B0-base`（101/54/0/0 那一发的名册是 54 枚；本程基线自身带一枚既有 flake，
所以差集是**对红名集合**做、不是对四数做）。每一发的差集判据：
"这一发新红的名字"与"被这一发吞掉的名字"两栏都要看——**一条用例 panic 会吞掉同包其余几十条**，
那些既不算红也不算绿。本程整包读数的机检列（`buildfailed`／`panic`）逐发都是 **0／0**，
名册行数：形一六发各 53、形二六发各 53、`R2-p1-full` 那发 53、`p7`／`p8` 两形三发各 53
（基线 54，少的正是被 `-skip` 掉的那一枚）⇒ **没有任何一发出现名册缩小**。

| 发 | 红名（本程实测） | 与基线的差 | 变绿还是被跳过 |
| --- | --- | --- | --- |
| `p1`/`p2`/`p3`/`p4`/`p5`/`p6` 六发（只跑本尺） | 只有 `TestAC1AC2DispatchHopGate133` | 只多本尺一枚；其余 53 枚**没跑**（`-run` 只点它），不算红也不算绿 | 本尺 `--- FAIL` 1、`--- SKIP` 0 |
| 门关着六发 | 只有本尺 | 名册 53 枚＝基线 54 － 131 的门（被 `-skip`） | 131 的门是**少跑**不是 SKIP（`--- SKIP` 行 0，RUN 101→100） |
| 摘掉本尺六发 | 四发＝131 的门；两发＝无人 | X14 两拍 100/53/0/0 ⇒ **真全绿**（不是被跳过：SKIP 0、本尺名字不在名册里） | 同上形 |
| N-3（摘掉本尺） | 131 的门 ＋ 既有 flake | flake 那枚在**基线**里也红 ⇒ 不是这一发新引入 | flake 是 `--- FAIL` 不是 SKIP |
| `p7` 三发（§1.11／§1.12） | **零枚红**（整包两形都 100/53/0/0） | 名册 53 枚＝基线 54 － 被 `-skip` 的那一枚；**没有缩小、也没有新增** | 本尺是 `--- PASS`（真跑真绿），不是 SKIP；`--- SKIP` 行 0、`panic` 0 |
| `p8` 两发（§1.13／§1.12） | 单跑本尺 rc=0；摘掉本尺整包 rc=0 | 同上 | 同上 |

⇒ **判"不再红"这一族本程分清了两形**：`-skip` 掉的那一枚是"少跑"（RUN 少 1、SKIP 列仍 0），
X14 两拍的"零红"是"真跑完全绿"（100/53/0/0、`--- SKIP` 0）。回修方 §3.3 那句"关门是少跑而不是跑成 SKIP"
本程复核成立。软链 temp 那一形本程**未复核**（§4 第 3 条）。

### 1.11 本程自造的两发（`p7`／`p8`）：前任验收方没造、回修方也没造

派单那条规矩："验收方没造的支别记成已测"。这两发都是本程自己写的台件
（`/d/tmp/wisp133-r2-p7.sh`、`/d/tmp/wisp133-r2-p8.sh`），都种在 **X14 第二拍之上**（`install` 那一截已经删掉，
只剩判据 (d) 能让它红）＋ 同一行 usage，与 §1.2–§1.6 那六发同场。

#### 1.11 `p7`：一枚**签名完全正确、但本平台的测试二进制根本不编译**的用例 —— 读数 **绿**

```
树 T14：shot x14b2c + probe p7
R2-P7 planted
  probe133r2b_linux_test.go:7: func TestR2P7LinuxOnlyCaseDrivesTheLeg(t *testing.T) { _ = cmdSfx131(nil) }
  main.go:80: 		case "sfx131":                    （x14b2c 的 dispatch 仍在）
BUILD-cmdwisp R2-p7 rc=0 / BUILD-all rc=0
READ label=R2-p7 rc=0 RUN=1 TOPPASS=1 TOPFAIL=0 TOPSKIP=0     ← 本尺判绿
GATE-ROWS: run=1 fail=0 pass=1
账本（12 legs）里那一行原文：
  leg sfx131 main.go:81 installs=false handoff=false covered=test TestR2P7LinuxOnlyCaseDrivesTheLeg drives cmdSfx131 entries=cmdSfx131 aliases=
```

**这一枚"用例"到底编不编**——本程不靠"应该是 linux 独占吧"，在同一棵种好的树上直接问构建系统：

```
$ cd /d/tmp/wisp133-r2-T14/cmd/wisp
$ GOOS=windows go list -f '{{.TestGoFiles}}' . | grep -c probe133r2b   -> 0      （windows 的测试二进制里没有它）
$ GOOS=linux   go list -f '{{.TestGoFiles}}' . | grep -c probe133r2b   -> 1      （linux 里才有）
$ GOOS=windows go list -f '{{.TestGoFiles}}' . | grep -c leg_dispatch_gate_133_test -> 1  （本尺两平台都编）
```

⇒ **被验那一发跑的是 windows 腿**（`sh scripts/wisp-cli-tests.sh` 是 windows 腿，见回修方 §6 第 6 条；
本程所有读数也都在 windows 上取的）。所以 `TestR2P7LinuxOnlyCaseDrivesTheLeg` 是**一枚今天不会被编译、
更不会被 `=== RUN` 的声明**，而它把账本的 `covered=` 那一列写成了
`covered=test TestR2P7LinuxOnlyCaseDrivesTheLeg drives cmdSfx131`，并让 X14 第二拍**判绿**。
确认它"没跑／没编"的两发整包读数在 §1.12。

**它与回修方自报的那条代价的关系**：`R-133-2` 的函数值那一支它登记了"build-tag 藏起来的用例今天仍会被认作**钉**"
——那说的是钉表那一侧（`covered=nail`）。`p7` 走的是**另一条判据**：`drivenBy133`（`:1405`）读同一枚
`p.tests` 桶，产出的是 `covered=test …` 那一列。**那一支没有被登记过**，且本尺本次改动的三处文字
**都在反着说这件事**：

```
:398-405  "the coverage bucket holds ONLY declarations Go's testing package can run"
:1397-1400 "… since R-133-1 admits only a top-level func TestXxx(t *testing.T) …
            i.e. only declarations the testing package can run"
:1338     "A renamed or build-tag-hidden case has to be red somewhere, and here is where."
```

而 `:245-246` 的 `pkg133` 结构体注释写的是"**Build tags are ignored** … the approximation runs in the
direction that asks for **more coverage, never less**"。`p7` 量的正是那半句的反例：忽略构建标签在这里
**放行**了一枚覆盖（不是多要一枚），⇒ "never less" 这一句**被实测推翻**。
文件头 §"WHAT THIS FILE CANNOT SEE" 那一列（AC#2 要求"写明另一种能看见而它看不见的例子"就是落在这张列上）
第 4 条写的是"this gate reconciles claims against **compiled test sources**"——`p7` 也把它证成假话：
它reconcile 的是**没编译**的源码。

#### 1.12 `p7` 的两发整包确认（它是本尺独占的目击，还是别人也接着？）

时刻 `date -u` 15:49:42Z／15:51:03Z（本地 23:49／23:51 +08）。每发先 `LANDED-MAIN`＋`go build ./cmd/wisp/` rc=0。

| 发 | 树 | 命令 | 读数 | 含义 |
| --- | --- | --- | --- | --- |
| `R2-p7-full` | T18 | x14b2c＋p7，整包 `-count=1 -skip '^TestAC4EveryLegIsNailedOrRuled$'` | build rc=0、**rc=0 100/53/0/0**、本尺 `run=1 pass=1 fail=0`、`--- SKIP` 0、`panic` 0 | **门关着、整包全绿**：这形本来只有本尺能看，而它今天放行了 |
| `R2-p7-nogate` | T19 | x14b2c＋p7，整包 `-count=1 -skip '^TestAC1AC2DispatchHopGate133$'` | **rc=0 100/53/0/0**、`--- SKIP` 0、`panic` 0 | 摘掉本尺仍全绿 ⇒ **131 的门也不接这一形**（与 §1.9 形二的 X14 两拍同账） |
| `R2-p8-nogate` | T20 | x14b2c＋p8，摘掉本尺 | **rc=0 100/53/0/0** | `p8` 那一形同样无人接（它是 §2.2 已登记的残留，不参与判退回） |

`p7` 那枚"用例"到底跑没跑（同一份整包日志里三条计数，`R2-p7-full.log`）：

```
grep -c '^=== RUN   TestR2P7LinuxOnlyCaseDrivesTheLeg$'   -> 0
grep -c '^--- (PASS|FAIL): TestR2P7LinuxOnlyCaseDrivesTheLeg' -> 0
名字在整份日志里出现 1 次，唯一那次就是本尺自己打的账本行 :85
顶层 --- (PASS|FAIL): Test 共 53 枚 ＝ 摘掉本尺的名册数，一枚没多
```

**同场对照**（`R2-p1-full.log`，x14b2c＋p1，门关着整包）：

```
READ label=R2-p1-full rc=1 RUN=100 TOPPASS=52 TOPFAIL=1 TOPSKIP=0
REDNAMES: --- FAIL: TestAC1AC2DispatchHopGate133
grep -c '^=== RUN   TestSfx131LegIsDriven$'  -> 0
名字在整份日志里出现 0 次                      ← 修法生效后账本连这个名字都不再打
```

⇒ 这两发的差**就是 `R-133-1` 那一格修没修好的形状**：p1 的同名方法现在连账本都不进；
p7 的"签名正确但本平台不编译"那一枚**仍然进**，而且是把 `covered=` 那一列写成 `test …` 之后进去的。

#### 1.13 `p8`：裁决句坐在**真正 owning 这条腿的那枚文件**里、正文不点任何符号名 —— 读数 **绿**

```
树 T15：shot x14b2c + probe p8（本程自造）
R2-P8 planted
  main.go:80: 		case "sfx131":
  main.go:142: // WISP-LEG-COVERAGE-RULING: sfx131 acceptor r2 p8: prose ruling parked in the
BUILD-cmdwisp R2-p8 rc=0 / BUILD-all rc=0
READ label=R2-p8 rc=0 RUN=1 TOPPASS=1 TOPFAIL=0 TOPSKIP=0     ← 本尺判绿
账本：leg sfx131 main.go:81 … covered=ruling main.go:142 entries=cmdSfx131
```

⇒ 这一发**不算本格的破口**，它是 §2.2 那一半句（"裁决句必须点到该腿的 entry 名"）的**残留形状的实测**：
同文件核对确实放行"写在该写地方、但一个符号名都不点"的裁决句，而今天树上那五枚裁决句**全是这一形**（§2.2(a)）。
回修方在 §6 第 3 条自己写下了这句残留（"一枚坐在 `main.go` 里、点名 `version`、但不提任何符号名的裁决句仍然算覆盖"），
本程把它从"自述"量成"读数"，并且量到它**不止 `version` 一枚**、对任何被分发的腿都成立（`p8` 种的是新腿 `sfx131`）。
这一支按派单口径归 **AC#3／AC#5**（它自报的理由 §2.2 判成立），另立 `R-133-11` 登记形状与归属。

## 2. 它自报"没做"的那两半句：代价理由独立重走（不许照抄）

回修方在票面末尾与 §6 第 2／3 条自报两半句没做，各给了一句理由。派单要求这两条理由**必须被独立重走**。
本程两条都重走了，**结论：两条理由都成立**，各有一条措辞精度要修（§2.1 末、§2.4）。

### 2.1 `R-133-2` 第一选项"登记表存函数值"装不了 —— 理由**成立**

它给的理由是：`legCovers133` 那四枚声明指向的用例全在 `_windows_test.go` 后缀里，本尺是跨平台文件，
换成函数值会让 linux 侧 `go vet`／`go test` 直接 `undefined:`。逐枚 `grep` ＋ 逐枚看后缀 ＋ 一发真实验：

**(a) 四枚声明与它们各自的声明处（本程自己 `grep -rn "func <名>" cmd/wisp/`）**

```
cmd/wisp/leg_dispatch_gate_133_test.go:135-140  legCovers133 的四枚声明：
  {leg: "run",     test: "TestAC2SealNoticeLandsInTheRunLegLogFile",           entry: "runTextTask"}
  {leg: "models",  test: "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk",       entry: "cmdModels"}
  {leg: "secret",  test: "TestAC3SecretLegBooksItsAuditRecordsOnDisk",         entry: "cmdSecret"}
  {leg: "no-args", test: "TestAC1ResidentLegInstallsItsLogListenerOnDisk",     entry: "runResident"}
声明处（与它写的行号逐枚相同）：
  cmd/wisp/leg_sink_nail_131_windows_test.go:346 func TestAC2ModelsLegBooksItsHandOffVerdictOnDisk(t *testing.T) {
  cmd/wisp/leg_sink_nail_131_windows_test.go:407 func TestAC3SecretLegBooksItsAuditRecordsOnDisk(t *testing.T) {
  cmd/wisp/logsink_windows_test.go:137           func TestAC2SealNoticeLandsInTheRunLegLogFile(t *testing.T) {
  cmd/wisp/resident_sink_nail_127_windows_test.go:393 func TestAC1ResidentLegInstallsItsLogListenerOnDisk(t *testing.T) {
后缀核对：四枚**全部**落在 `_windows_test.go` 里（其中 `logsink_windows_test.go`、
`resident_sink_nail_127_windows_test.go` 另带显式 `//go:build windows`；
`leg_sink_nail_131_windows_test.go` 只靠文件名后缀那一条隐式约束，首行就是 `package main`）。
本尺 `leg_dispatch_gate_133_test.go` 无构建后缀、无 `//go:build` ⇒ 两平台都编。
```

**(b) 平台文件集实测（这条最硬，且是仓外快照树上量的）**

```
$ cd /d/tmp/wisp133-r2-tree1
$ GOOS=linux   go list -f 'LINUX TestGoFiles={{.TestGoFiles}}' ./cmd/wisp/
  [dataroot_128_test.go leg_dispatch_gate_133_test.go leg_sink_gate_131_test.go logsink_test.go
   providers_test.go run_mode101_test.go run_test.go secret_dataroot_119b_test.go secret_test.go
   tempdir_resolved_124_test.go]                                        ← 10 枚，`windows_test` 命中 0
$ GOOS=windows go list -f 'WIN   TestGoFiles={{.TestGoFiles}}' ./cmd/wisp/
  [… leg_dispatch_gate_133_test.go leg_sink_nail_131_windows_test.go logsink_test.go
     logsink_windows_test.go resident_sink_nail_127_windows_test.go …]  ← 15 枚，三枚 windows 件都在
```

⇒ **`leg_dispatch_gate_133_test.go` 在 linux 的文件集里，而那三枚 windows 件不在** ⇒ 从本尺引用那四枚
函数名，在 linux 上必然 `undefined`。理由的**语言层机制成立**。

**(c) 一发真实验（树 `/d/tmp/wisp133-r2-FV`，往本尺所在那枚跨平台文件旁加一枚同包跨平台件，
把四枚声明写成函数值 `[]func(*testing.T){…}`）**

```
落地证明：grep -n "TestAC2SealNotice" cmd/wisp/zz_r2_funcval_probe_test.go -> :7
$ go vet ./cmd/wisp/                    （GOOS=windows 原生）  rc=0
$ GOOS=linux go vet ./cmd/wisp/         （宿主交叉）            rc=1，三行诊断全是
      imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files
      grep -c "undefined"                = 0
      grep -cE "cmd[/\\]wisp[/\\]…\.go:[0-9]+" = 0     ← 它死在包加载，根本没走到类型检查
$ 容器 golang:1.27 / CGO_ENABLED=1 / GOPROXY=off / 仓库与 GOMODCACHE 各 :ro 挂进 /src、/gomod
  挂载先证：ls -l /src/go.mod 883 字节、/src/cmd/wisp/zz_r2_funcval_probe_test.go 336 字节
  未种件那一态： go vet ./cmd/wisp/  -> rc=0，诊断 0 行       （= 回修方 §5 的 VET_RC_0，本程独立复现）
  种上函数值那一态：go vet ./cmd/wisp/ -> rc=1
      vet: cmd/wisp/zz_r2_funcval_probe_test.go:7:2: undefined: TestAC2SealNoticeLandsInTheRunLegLogFile
```

⇒ **`undefined` 是真会发生的，但发生在那枚"真做 linux 类型检查"的读数里（容器 `go vet`、以及 linux 的
`go test`），不在宿主那条 `GOOS=linux go vet ./cmd/wisp/` 里**——后者今天不带种件也 rc=1、且诊断里
一枚 `cmd/wisp/*.go:line` 都没有。所以它那句"linux 侧 `go vet`／`go test` 直接 `undefined:`"
**拆开看才对**：`go test`（容器／真机 linux）对，宿主 `GOOS=linux go vet` 那一条**本来就红在 cgo 包加载**，
换不成 `undefined`。这一条按**误记（措辞精度）**登记，不改判据、不算破口（§5 的 `R-133-10`，信息级）。

**裁定**：`R-133-2` 第一选项"登记表存函数值"的理由**成立** ⇒ 这一支**不记成 AC#2 本格的破口**，
与它自报的一致，归 **AC#3／AC#5**（要把钉注册点搬进 windows 那枚文件、或给本尺加构建标签分文件才装得下）。
它登记的那枚代价（"`covered=nail` 那一列仍是按名字的源码查表，build-tag 藏起来的用例今天仍会被认作钉"）
**本程用自造的 `p7` 把它从"代价"量成了"实测"，并发现本尺的红名文案正在反着说这句话** ⇒ 见 §1.9／§5 `R-133-9`。

### 2.2 `R-133-3` 那半句"裁决句必须点到该腿的 entry 名"没做 —— 理由**成立**

它给的理由是：今天五枚裁决句没有一枚写了 entry 符号名，装上这条判据就得改三枚生产文件的裁决注释才不红。
本程不照抄，做了两步：

**(a) 五枚裁决句逐枚定位 + 逐枚判"正文里有没有出现该腿的 entry 符号名"**

腿名与 entry 取自本程 §1.1 基线账本（本尺自己打的 `entries=` 列），裁决句正文＝从 `WISP-LEG-COVERAGE-RULING:`
那一行起、连续以 `//` 开头的那一段（**不含**后面的 `func` 声明行——第一遍本程把 `slo_windows.go:193`
的 `func cmdSLO` 一起圈进来，量出过一枚假命中 `cmdSLO`，作废重走才得到下面的结果）：

```
version       cmd/wisp/main.go:65..78          entry 候选 attachParentConsole|printVersions  命中 NONE
help          cmd/wisp/main.go:70..78          entry 候选 attachParentConsole                命中 NONE
default       cmd/wisp/main.go:74..78          entry 候选 attachParentConsole                命中 NONE
panel-assets  cmd/wisp/panel_assets.go:16..22  entry 候选 attachParentConsole|cmdPanelAssets 命中 NONE
slo           cmd/wisp/slo_windows.go:186..190 entry 候选 cmdSLO                             命中 NONE
```

⇒ 它说的"没有一枚写到 entry 符号名"**逐枚成立**，位置也逐枚对得上（`main.go:65/70/74`、`panel_assets.go:16`、
`slo_windows.go:186`）。

**(b) 装上这条判据要动哪几枚文件**

五枚裁决句落在 **3 枚生产文件**里（`main.go`、`panel_assets.go`、`slo_windows.go`）
⇒ "就得改三枚生产文件的裁决注释才不红"成立，且那三枚**都不在回修方那枚派单允许它动的文件里**
（它只许动 `cmd/wisp/leg_dispatch_gate_133_test.go`）。

**(c) 本程另取一发量它的残留形状**（`p8`，§1.8 之后）：裁决句坐在**真正 owning 这条腿的那枚文件**
（`main.go`，`sfx131` 的 dispatch 就在 `main.go:80`）、正文**不点任何符号名** ⇒ 若判绿，则"点到 entry 名"
那半句的缺口是**活的**而不是纸面的。读数见 §1.13。

**裁定**：`R-133-3` 那半句的理由**成立** ⇒ 这一支**不记成 AC#2 本格的破口**（`R-133-3` 在裁决表 §8 与票面
"翻格条件"那一行里本来就写作"**可同批**"、不是必办；必办只有 `R-133-1`＋`R-133-6`），
归 **AC#3／AC#5**，与本程实测到的残留形状一起落账。

### 2.3 两半句合起来的一句话

两条"没做"都不是"没做而无人知"：两条理由本程都独立重走成**成立**，代价都被回修方自己写在 §6 第 2／3 条里。
差别只在**它给的两条代价描述各有一处比实测轻**：函数值那一支的 linux 失效面它说宽了（§2.1(c)），
build-tag 那一支它说窄了——红名文案 `:1338` 明写"A renamed or **build-tag-hidden** case has to be red
somewhere, **and here is where**"，而本程量到它不在这里红（`p7`）。后者是本格第二次退回的唯一实测理由，
写进 §5 `R-133-9`、总判写进 §3。

## 3. AC#2 对表（逐条 1:1，左栏要么是票面 AC#2 的判据物、要么是裁决表 §8 那批账）

### 3.1 退回件那五条账＋一条证据面账：本程逐发读数

| 编号 | 回修方交回的话 | 本程独立复算（哪一发） | 与它一致？ | 判 |
| --- | --- | --- | --- | --- |
| `R-133-1`（高，**必办**） | 覆盖桶只收 `isRunnableCase133`，其余落 `helpers` | §1.3 `p1`：改前绿 ⇒ 本程改后 **rc=1 红**、红名 `TestAC1AC2DispatchHopGate133`、账本 `covered=RED nothing`；§1.12 整包里那枚同名方法**连名字都不进账本**（出现 0 次） | **一致** | **已办**（它点名的那一形：方法／签名） |
| `R-133-2`（中，可同批） | 钉表查收紧后的桶 ⇒ 文案与判据一致，并点名找到的是什么 | §1.6 `p3`：改前绿 ⇒ 本程 **rc=1 红两枚**，第二枚原文点到 `isRunnableCase133` 与"helper or a method" | **一致** | 名字那一侧**已办**；⚠ 同一枚文案的**另一半**（`:1338` "build-tag-hidden case … here is where"）本程造出反例 ⇒ `R-133-9`（§5） |
| `R-133-3`（中，可同批） | 行扫描保留，裁决句只坐在 owning 文件里才算覆盖，放错地方单独报红 | §1.4 `p2`：改前绿 ⇒ 本程 **rc=1 红两枚**（misplaced ＋ covered by nothing）；§1.13 `p8` 正向形判绿（同文件、不点符号名 ⇒ 放行，与它 §6 第 3 条自报的残留同形） | **一致** | **部分办**：同文件核对已办；"点到 entry 名"那半句未办，理由本程判**成立**（§2.2）⇒ 不记本格破口，`R-133-11` 归 AC#3/AC#5 |
| `R-133-4`（中，可同批） | 每条命令标签各出一行腿，`-` 拼写是命令词前缀才折 `aliases=` | §1.7 `p5`：改前"11 legs、`runalt133` 零行、绿"⇒ 本程 **rc=1 红两枚、12 legs、`runalt133` 有自己的行**；正向控制在基线账本里 `leg version aliases=--version\|-v`、`leg help aliases=--help\|-h`、仍 11 legs | **一致** | **已办** |
| `R-133-5`（低，验收方未造） | 本程造出 `p6` 并修：复合条件里每一处 argv 比较各出一腿 | §1.8 `p6`（本程另写台件复种）：本程改后 **rc=1 红两枚、13 legs、`leg x133a … covered=RED nothing`**、名字出现 3 次 | 改后**一致**；改前那一态**本程未复跑** | **修法有效已独立复现**，但"升成已证"的证据链只有一端（§4 第 1 条） |
| `R-133-6`（中，**必办**，证据面） | 前任 §3 那句全称量词旁插入更正段、原文未抹 | 本程自己复算那六发（§1.9 形二）：N-3/X4/X8/X12 摘掉本尺仍各红 131 的门、X14 两拍 rc=0 100/53/0/0；`git show --numstat e2e61a0 -- docs/evidence/s1/133-ac1-ac2-instrument.md` ＝ **20 插入／0 删除**，原句仍在 `:245`、更正段在 `:248` | **一致** | **已办** |

### 3.2 票面 AC#2 的四条判据物（照裁决表 §6 那四行重判一遍）

| 要求 | 本程读数 | 判 |
| --- | --- | --- |
| 选了清单式还是图式 | 清单式（`func main` 分支 ＋ `usage` 文字块互核、双向差）；本程 `p5`/`p6` 两发的第二枚红（`AC#1/#2 RED: … documented in no line of the usage block`）就是双向差**在咬** | **成立** |
| 写明"另一种能看见而它看不见的例子" | 实现方证据 `133-ac1-ac2-instrument.md:269-287` 给了三条形＋三条反形（本程未复做那六形变异实验，§4 第 6 条）。⚠ 但**同尺源码里的"看不见"清单有一条是假的**：`:57-72` 第 4 条写 "this gate reconciles claims against **compiled test sources**"，`p7` 证它 reconcile 的是**没编译**的源码；`:245-246` 的 "the approximation runs in the direction that asks for more coverage, **never less**" 也被同一发反证 | **不完全成立**（清单在、其中一项与实测相反 ⇒ `R-133-9`） |
| 不许写"顺带全覆盖" | `:287` 明写"本格不写顺带全覆盖，只写上面这六条"；本程在被验那枚文件里未找到任何无条件的全覆盖句子 | **成立** |
| 要说得**出**删掉它哪条用例会红 | 答案＝`TestAC1AC2DispatchHopGate133`（名册 54→53、逐名 `comm` 本程独立复现 §1.1）；那句被加强过的全称断言已按 `R-133-6` 更正，本程自己复算支持更正后的版本（§1.9 形二） | **成立**（就这一问本身） |

### 3.3 总判：**退回**（AC#2 这一格第二次被退回 ⇒ 本表不再开第三格，形状写在 §3.4）

退回的理由**只有一条**，且是**本程自己造出来**的，不是文字问题：

> `p7`（§1.11／§1.12）：一枚签名完全正确、但 `GOOS=windows` 的测试二进制**根本不编译**的
> `func TestR2P7LinuxOnlyCaseDrivesTheLeg(t *testing.T)`，配上给人看的那行 usage，
> 就让 X14 第二拍（`installLogSink` 那一截已经删掉）**判绿**，账本印
> `covered=test TestR2P7LinuxOnlyCaseDrivesTheLeg drives cmdSfx131`；整包两形都 100/53/0/0、`--- SKIP` 0、
> `=== RUN` 里那枚名字 0 次。
> ⇒ **这正是票面 AC#2 立的规矩要防的结局**（"被什么覆盖"那一列可以写假而全绿），
> 也正是上一轮 `R-133-1` 判退回时用的那把尺（裁决表 §6.3）。
> 按本仓硬线：**造出来了就不是附条件**。

同时把话说公平（防下一位把这格读成"回修方什么都没修"）：
`R-133-1`／`R-133-4`／`R-133-5`／`R-133-6` 四本账本程**逐发复现其修法有效**，
`R-133-2`／`R-133-3` 各差半句、那半句的"装不了/出了地界"两条理由本程**独立重走成成立**（§2）。
这一格卡在**同一族里更靠外的一枚**：修好的是"receiver 与签名"这一形，没修的是"这个声明今天编不编"这一形，
而那枚新判据的**三处文字都在声称后一形已经收紧**（`:398-405`、`:1338`、`:1397-1400`）。

### 3.4 这一格不再续第三格：两件事的形状（派单第 0 节点名要的那两样）

**① 把这一发的新形状做成家族票的第 N 发** —— 交 **票 135**。
135 面上现有的那一族是 `M-A`／`M-B`／`M-C`／`M-D′`／`M-E`／`M-F` 六发（外加 `M-F＋M-C` 那一枚双叠），
所以本程这一形记作 **票 135 的第 7 发，建议编号 `M-G`**：

> **`M-G`**：一枚**本平台不编译**的 `func TestXxx(t *testing.T)`（种法＝文件名 `_linux_test.go` 后缀，
> 或显式 `//go:build` 表达式为假），坐在**被验那一腿自己的目录**里、身体里调用那条腿的 entry。
> 判据：本尺必须红、红名点到本尺，且**第二拍也必须红**（把 install 那一截删掉之后仍须红）；
> 反向控制一枚：同形写在**本会编译**的文件里 ⇒ 必须绿（不许把合法形状打死）。
> 归 135 的 **AC#1/AC#2 那一根**（自证腿不许是哑的 ＋ 把这一形做成常备用例），
> 与 135 已有的 **AC#7**（`R-131r3-1`，131 的第二把尺只到文件粒度）是**同族不同尺**。
> ⚠ **两票不许互相以为对方已修完**：票 133 只修了自己这把尺的"名字/签名"那一形，
> 票 135 那一发若只修 131 的第二把尺，本发**不算已办**；反之本表也不替 135 结 AC#7。

**② 把本票 AC#2 的翻格条件改成一条可复算的跨票依赖**。下面三条**同时**在盘上，编排者才许翻这一格；
任一不在 ⇒ 本格保持退回，**不再另开第三格**：

- **C1（本票侧，任意 sha 上可复算）**：在 dev 的某枚 sha `S` 上跑本程那台件的原样命令——
  `sh /d/tmp/wisp133-r2-run.sh <新标签> <新树> x14b2c p7 "-count=1 -run '^TestAC1AC2DispatchHopGate133$'"`
  （种法写死在 `/d/tmp/wisp133-r2-p7.sh`，本表 §1.11 逐字引过），读数须为 **rc=1、红名 `TestAC1AC2DispatchHopGate133`**；
  并且**同一棵树的整包门关着读数不得仍是 `rc=0 100/53/0/0`**（本表今天它就是 ⇒ 这一句就是本格的现状）。
- **C2（家族票侧）**：票 135 面上那一发 `M-G` 有**非实现者**出的三态读数
  （种件落地行＋`go build` rc＋红名，两拍都齐），落在 `docs/evidence/s1/135-*.md` 的具体一节，
  并给出**该 commit 的短哈希**。⚠ 只登记 AC#7（131 的第二把尺）**不满足 C2**。
- **C3（门禁侧，不是本格的牙）**：若要把 `M-G` 升成"CI 守"，须给**具体 run id ＋ job id ＋ step 名**；
  给不出就当那道门不存在。本程 `gh run list --limit 25` 取数失败（`api.github.com … EOF`）⇒
  **本表未引任何 CI 读数**，AC#2 的判据物本身也不含 CI 落点（那是 AC#6 的账）。

## 4. 我没做的档（诚实列，不给下一位留"这些已验"的错觉）

1. **`p6` 的"改前"那一态未复跑**。它要在**没有第 5 枚修法**的树上跑＝回退被验版本；
   本程只复现了改后（红）与"种法/病名对得上"两态 ⇒ `R-133-5` 的"三态"本程**只有两态**。
2. **`p7`／`p8` 只裁了本尺的读数**，没有做成常备用例、没有跑双叠（`M-G` ＋ 另一枚变异）；
   那两样是票 135 那一发的活（§3.4 ①）。
3. **软链 temp 那一形未复核**（回修方 §5 末的普通形 101/54/0/0 vs 真符号链接形 rc=1 92/26/28/0、
   以及改前控制 `54123e0` 的 92/22/32/0、九枚子用例逐名）。本程未复现该形任何一个数。
4. **门禁面基本未复算**：`gofmt`／`gofumpt`（版本未钉那笔账）／`sh scripts/d22scan.sh` 各 scope／
   `sh scripts/wisp-cli-tests.sh` 同形那发——属 AC#5 与 `R-133-7` 地界。
   本程只为 §2.1 的理由复算了**容器 `go vet` 两态**与宿主 `GOOS=linux go vet`，未复算容器 `go test`。
5. **linux 面**：`GOOS=linux go test ./cmd/wisp/` 那 19 枚 FAIL 未复跑（package owner 地界）。
6. **实现方 §3 那"三条形＋三条反形"未做变异实验**（与前任验收方同口径：论述不等于变异）；
   方法值 `v.M(...)` 那一支、跨包那一支本程**未造**，票面已把这两支划出本票地界 ⇒ 不算已测。
7. **五发的①栏（无仪器树 `c720494`）未重取**：本程用"同一棵树摘掉本尺"那一形（§1.9 形二），
   与前任①栏树基不同（那棵树连 `main.go`/`panel_assets.go`/`slo_windows.go` 的裁决注释都没有）⇒
   **两形别说成同一枚读数**（回修方 §6 第 7 条同一口径）。
8. **既有 flake 未做归因实验**：`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`
   在本程四发同命令里 3 红 1 绿（`0xc000013a`），未跑第 5、第 6 发去定命中率，也未判它是真伤还是负载；
   只判了两件事：红名不是本尺、不属 AC#2 判据物。
9. **CI 未取到任何 run 读数**（`gh` 网络失败）⇒ C3 按"那道门不存在"处理。
10. **票 135 面与 `internal/observe/**`、`internal/winsec/**` 未读未改**（在飞兄弟地界）；
    未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` **未读内容、未提交、未改**（来源未明）。

## 5. 新账（`R-133-9` 起，编号占用本程自己 `grep` 过：`R-133-1`..`R-133-8` 已被占）

| 编号 | 严重度 | 是什么 | 能否复现（怎么复现，快照目录名） | 修法方向 | 归谁 |
| --- | --- | --- | --- | --- | --- |
| **`R-133-9`** | **高**（AC#2 本次退回主因） | 覆盖桶**只看签名、不看这条声明今天编不编**：一枚 `_linux_test.go`（或 `//go:build` 为假）文件里的 `func TestXxx(t *testing.T)` 能把 X14 第二拍判绿，账本写 `covered=test … drives …`；同时被推翻的还有三处文字——`:398-405`"holds ONLY declarations Go's testing package can run"、`:1397-1400` 同义句、`:1338`"A renamed or **build-tag-hidden** case has to be red somewhere, **and here is where**"，以及文件头 AC#2 那张"看不见"清单第 4 条"reconciles claims against **compiled test sources**"与 `:245-246` 的"never less" | **能**。`/d/tmp/wisp133-r2-T14`（单跑本尺 rc=0）／`-T18`（门关着整包 100/53/0/0）／`-T19`（摘掉本尺整包 100/53/0/0）。种法 `sh /d/tmp/wisp133-r2-p7.sh <树> <标签>`；构建标签那一面不靠推测，用 `GOOS=windows go list -f {{.TestGoFiles}} . \| grep -c probe133r2b` ＝ **0**、`GOOS=linux` ＝ **1** 直接钉死。日志 `out/R2-p7*.log` | **不许折回只信 AST**（票面 §"修法形状不许是"那条、票 135 同一条已裁过一次，理由复用）。方向：给"进 `p.tests`"这一步加**平台求值**——按文件名后缀（`_windows.go`/`_linux.go`/`_other.go`）＋ `//go:build` 表达式对**当前 `runtime.GOOS`/`GOARCH`/构建标签**求值；不通过的声明可另进一桶（走图仍看得见 ⇒ "跨平台互借形状"那半句近似性保留），但**不得充当覆盖**，且要红着说"这条覆盖证据今天不在本平台的二进制里"。并给这条判据装一枚 135 面上 `M-G` 的**常备用例**（AC#3 那句"主动弄哑自己的自证腿"正是它的家） | **票 133 续单**（仪器本体这一半）＋ **票 135 第 7 发 `M-G`**（自证腿那一半）；本票 AC#2 翻格条件 C1/C2 |
| **`R-133-10`** | 信息（误记，不改判据） | 回修方 §2.1 末与 §6 第 2 条那句"换成函数值会让 linux 侧 `go vet`／`go test` 直接 `undefined:`"——对 `go test` 与**容器 `go vet`** 成立（本程量到 `vet: …undefined: TestAC2SealNoticeLandsInTheRunLegLogFile`），对**宿主那条 `GOOS=linux go vet ./cmd/wisp/` 不成立**：它今天不带种件也 rc=1、诊断里 `cmd/wisp/*.go:line` 命中 0，走不到类型检查 | **能**。树 `/d/tmp/wisp133-r2-FV`（`out/FV-linux-hostvet.txt` 宿主 rc=1／`undefined` 计数 0；`out/FV-container-vet2.txt` 未种件 `BASELINE_VET_RC=0` 诊断 0 行、种上 `PLANTED_VET_RC=1` ＋ 那一行 `undefined:`；另 `out/FV-win-vet.log` windows 原生 rc=0） | 证据面下次落笔时把它写成"linux 侧 `go test`／容器 `go vet` 会 `undefined:`；宿主交叉 vet 今天已红在 cgo 包加载"。**本验收方不改它的件**，只登记 | 登记，不归任何票（属"引用即须现核"那一族） |
| **`R-133-11`** | 中（**不是本格的牙**） | `R-133-3` 未做那半句的**活残留**实测：裁决句坐在**真正 owning 这条腿的那枚文件**里、正文一个符号名都不点 ⇒ 判绿、账本 `covered=ruling main.go:142`。今天树上那五枚裁决句**全部是这个形**（§2.2(a) 逐枚 `NONE`） | **能**。树 `/d/tmp/wisp133-r2-T15`（单跑本尺 rc=0）／`-T20`（摘掉本尺整包 rc=0）。种法 `sh /d/tmp/wisp133-r2-p8.sh <树> <标签>`；日志 `out/R2-p8.log`、`out/R2-p8-nogate.log` | 装"裁决句必须点到该腿 `entries=` 里至少一枚符号名"那条判据，**代价是改 `main.go`／`panel_assets.go`／`slo_windows.go` 三枚生产文件的裁决注释**（本程已独立核过枚数＝3） | **AC#3／AC#5 那一格**（回修方 §6 第 3 条与它自己的 `next=` 都是这个归属，本程判成立） |
| **`R-133-12`** | 信息（既有件，不归本票） | 锚点树未变异整包里 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 四发同命令 3 红 1 绿，红点是 `resident_sink_nail_127_windows_test.go:549/572/583`、子进程退出码 `0xc000013a`、sink 目录"1 file, 0 records"。它同时污染 §1.9 形二的 `R2-RM-n3`（那发 TOPFAIL=2） | 能，`out/R2-B0-base.log`、`R2-B0-c2.log`、`R2-B0-nogate.log`（三红）与 `R2-B1-base.log`、`R2-B1-c2.log`（两绿） | 本程未动。写给下一位：**这台机器上"四数"不是稳定量**，比较读数前先对齐红名集合与名册差集；命中率的定量归 flake 那一格（票 136 面上已见 `AC#11` 编号在议） | 登记，不归本票 |

## 6. 两个计数（本仓硬形制，分栏填，不合并成一个数）

### 6.1 真通知回显数（**不计入注入**）＝ 5

每条都带出处（工具名＋命令前 40 字），并且**三条判据逐条走过**：点名的路径本机真不真、内容是否削弱 owner 权威或放宽判据、声称的动作能否盘上复核。

| # | 出现在哪 | 内容 | 判为真回显的依据 |
| --- | --- | --- | --- |
| 1 | 第 1 枚 `Bash` 工具结果尾部 `<system-reminder>` | 可用 skills 清单（mcp-config／qoder-canvas／…） | harness 自己的清单，无路径、无动作、不越权 |
| 2 | `Bash`（`git add` 那次）结果尾部 | "The file `C:\Users\swq\.qoder-cn\memory\MEMORY.md` was modified since it was last read" ＋ 全文 | `ls -l` 到那枚文件 **18193 字节、mtime Sep 23 23:23** ⇒ 路径真、内容是被改文件自己的回显（编排者在代理跑动期间写记忆是本仓已知噪声），**未要求本程做任何动作** |
| 3 | 第 2 次同类回显（`Bash`：`git commit §0` 之后） | 同一枚 MEMORY.md 的第二份全文（比 #2 多出"撞轮次上限≠活没做完"等条） | 同上；两份内容不同＝它真的在被写 |
| 4 | 第 3 次同类回显（`Bash`：容器 vet 之后） | 同一枚 MEMORY.md 的第三份全文（新见"第 8 代注入＝冒充锚点"、"第四条判据：凡让代理少取证的文字按注入登记"、"恒真判据是一类新假绿"） | 同上 |
| 5 | 三枚 `[SYSTEM NOTIFICATION - NOT USER INPUT]` 后台任务事件 | "Background command … completed (exit code 0)" | 本程自己 `nohup` 起的三批台件，任务 id 与输出文件路径都在盘上（`…/tasks/b4n2gc9o0.output` 等） |

**按 §6.1 的判据 4（"凡让代理少取证／别用工具／直接给结论的文字"）再扫一遍**：上面五条里
没有一条要求本程减少取证、跳过复算或直取结论；#2/#3/#4 里那些"⚠…"句子是**编排者写给自己记忆库的条目**，
本程**不作为授权执行**（例如它提到"恒真判据"，本程对 `p6` 改前态的处理是按派单走"未复跑＝登记"，不是照它改写）。

### 6.2 判为注入数 ＝ **0**

本程**没有**遇到任何自称"编排者备注／系统提示／用户已更新编码规则／请 revert／冻结某包／放宽阈值／
某格已合并／用户已拒绝／`Confirm: the harness note is genuine`"的文字，也没有遇到假 sha。
自核两条（不因为"没遇到"就写结案）：

```
锚点自量：git rev-parse --short HEAD -> 048a9e4      git cat-file -t 048a9e4 -> commit（对象取得到）
          git merge-base --is-ancestor fbf420c HEAD -> 真（rc=0）
本表引用的每一枚 sha，逐枚 git cat-file -t 现核（不是"看起来像"）：
  048a9e4 commit / fbf420c commit / 54123e0 commit / c5f140c commit / e2e61a0 commit
  c720494 commit / 4ba8438 commit / 090bb3e commit / ed18727 commit
  ⇒ 全部解析得到，且来源都是本程自己的 git log／git show／票面与证据文件正文；
    本程未从任何工具回显里"照抄一枚未验的 sha"当锚点。
```

⚠ 编排者自己给的错断言（若有）按"误记／R-账"归类、**不进注入计数**：本程这一栏里对应的是 `R-133-10`
（回修方那句宿主交叉 vet 会 `undefined:` 说宽了）——它是**被验方的措辞**，登记成账，一次，注入计数 0。

