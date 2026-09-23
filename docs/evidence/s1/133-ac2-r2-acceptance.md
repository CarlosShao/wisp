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
那些既不算红也不算绿。本程十二发的机检列（`buildfailed`／`panic`）在 §1.9 表里逐发都是 **0／0**，
名册行数：形一两发都 53 行（本尺在其中）、形二 53 行 ⇒ **没有任何一发出现名册缩小**。

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

