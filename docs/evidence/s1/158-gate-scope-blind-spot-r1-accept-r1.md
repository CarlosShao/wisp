# 票 158 非实现者验收表 r1 —— 六格一次裁，锚钉 `c7f638c`

- 我是**验收方**，没写过被验的码（`SPEC-12 §4.3` #1/#3：裁决者≠实现者）。
- 被验物＝`docs/evidence/s1/158-gate-scope-blind-spot-r1.md`（**我的被验物，不是我的权威**）：锚点上现量 **1077 行**；`git diff --stat c7f638c..HEAD -- 那枚文件` ＝ **空**（HEAD 往前走的那一枚只加了派单文件）⇒ 我从盘上读到的那一份就是被验的那一版。
- 派单＝`.scratch/wisp/dispatches/2026-09-26-203x-accept-158-six-cells-pinned-c7f638c.md`。
- 票面＝`.scratch/wisp/issues/158-the-gate-stands-outside-the-shape-…-package.md`（锚点上现量 **84 行、6 枚 `- [ ]`、0 枚 `- [x]`**；本表末尾 §9 由我给出勾法）。

## 0. step 0 四件原文（派单 §0 要求的那四件）

```
$ date
Sat Sep 26 20:28:39 CST 2026
$ git rev-parse --short HEAD
ec5ecc0
$ git status --porcelain            # 全树：23 行已跟踪脏（design/** 那一堆、probes/152/my152.py、
                                    # 152-…-accept-r1.md 那 3 行未提交自校）＋ 21 枚未跟踪，原文与
                                    # 被验件 §A 那一串同形（差＝别人这一段新落的 probes/156、probes/158/r2 等）
$ git log --oneline -3
ec5ecc0 派单 20:3x: 票 158 非实现者验收（六格一次裁，锚钉 c7f638c）
c7f638c 票 158 r2 · 交件后自纠两处（追加，§A–§H 已入库那一行一字未改：git diff numstat deleted=0）
226d050 ledger(A310): 十四枚调研全部收官、全线封批——Q-56 带回三枚独立见证；票158 的病因改档

$ git rev-parse c7f638c
c7f638c5e9e3b66073a66d680b5e8c0b85a5b7e4        # ← 本表全部判据的锚点
```

- **锚点外还核了三件事**，免得后面每一发都要解释口径：
  `git diff --name-only c7f638c..HEAD -- '*.go'` ＝ **空**（HEAD 前进的那一枚没动任何 Go 码）；
  `git diff c7f638c -- internal cmd tools scripts` ＝ **空**（工作树在这四支上就是锚点那一版）；
  `git status -sb` 首行＝`## dev...cnb/dev [ahead 285]` ⇒ 本程与实现程都是**只 commit 未 push**。
- 工具版本现读（不抄被验件）：`go version go1.27.1 windows/amd64`；
  `$(go env GOPATH)/bin/gofumpt.exe --version` ＝ **`v0.12.0 (go1.27.1)`**（与被验件 §A 末、§D.5 那格同值，这一枚对得上）。
- 争用：开测前后各一枚 `Get-Process` 都命中 `Runner.Listener`（pid 3952，2028 与 2048 两次读数同一枚 pid）
  ⇒ **本机 self-hosted runner 同机在跑，争用未排除**。本表**不含任何延迟／阈值结论**，所有 `ok …Ns` 一律不当证据。
- 本程写面只有两处：本文件 ＋ `.scratch/wisp/probes/158/accept-r1/**`。两枚半件（`probes/152/my152.py`、
  `docs/evidence/s1/152-…-accept-r1.md`）**未碰、未还原、未 commit、未补完**；`frontend/**`、`design/**` **未写一个字节**。

## 1. 三把判据仪器先各自打过正控（派单 §1）

> 规矩：**打不出响声的尺上任何"零"都不算证据。** 三把尺各自的正控读数都在下面，全部本程现跑。

### 1.1 尺一＝门本体 G5（`.scratch/wisp/probes/154/gate-clauses.sh`）——逐发复算，不转抄 rc

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh c7f638c   > probes/158/accept-r1/gate-at-anchor.txt ; rc=0（脚本自身）
全文 124 行；按节数：grep -n '^## ' → 11 枚节标题（G1／G1b／G2／G3／G4／G5 主尺／G5 正控说明／G5-正控／G5 负一负说明／G5-负一负／两枚附）
```

| 腿 | 本程现量（锚点 `c7f638c`） | 被验件自报 | 对得上？ |
|---|---|---|---|
| **G5 主尺** | `UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)`、未成对枚数＝**1**、该腿 `# rc=1` ⇒ **响** | §2.2 同一行判语 | **是** |
| **G5 正控** | 名册**非空**（`bridge.go:642` 开＋`bridge.go:695` 合，两行都在）、`git grep rc=0`、未成对枚数＝**0** ⇒ **不响** | §2.3（它引 `:691`，那是注释面之前的行号；本程在锚点上读到的是 **`:695`**，位移与被验件 §B 末那条"注释净加 4 行"一致） | **是** |
| **G5 负一负** | 未成对枚数＝**5**（`panel_assets` ＋ 4 枚 `internal/risk/*_test.go`），该腿 `# rc=5` | §2.6 末"负一负＝5" | **是** |

**这把尺自己有没有牙（派单 §2 那两问，本程现量）**：

1. **成对判据剔注释那一行（`pair()` 的 `grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)'`）承不承重？** 承。
   同一射程、同一条腿，只把那行去掉：负一负 **5 → 6**，多出的正是
   `internal/tools/bridge_scope_open_ticket158_test.go`（它唯一的 `OpenScope` 出现在**注释**里）。
   原文＝`probes/158/accept-r1/g5-with-risk.txt` 旁边那一发（本件 §1.1 表下第二条命令）。
   ⇒ 这不是"摘掉导致编不过"那种假牙：**外部可见读数（枚数）从 5 变到 6。**
2. **`internal/risk/*` 那条排除承不承重（对被验件 §2.4"其实不承重"那句）？** 确实不承重，本程独立复算：
   把 `internal/risk` 放回射程（仍排 `*_test.go`）⇒ 开方文件 3 枚（`panel_assets`／`provenance.go`／`bridge.go`），
   未成对仍＝**1**（只有 `panel_assets.go`）⇒ 判语一字未变，与 §2.4 末那条"顺带量到的口径事实"同向。
3. **主尺射程里今天有没有"字符串字面量假响"（§2.4 登记的结构性残余）？** 没有：
   本程同尺现量 roster 行数＝**3**、剔注释后调用点行数＝**3** ⇒ 射程内一枚注释/字符串混入都没有。
   那条残余是"未来会假响"的形状，不是今天的读数——被验件就是这么写的，判**如实**。

（一处**过期字**要记：门本体脚本第 96 行附近那条正控说明注释仍写"`bridge.go`（`:642` 与 `:691` 同文件成对）"，
`:691` 在锚点上已是 `:695`。它在**注释**里、不进任何读数（判语是动态算的），**不承重**，但下一位照它去 `sed -n '691p'` 会读到别的东西。）

### 1.2 尺二＝AC#4 那 28 枚禁改面的名册尺（按 commit 号集合 ＋ `--no-walk`）

本程自己跑生成器（`.scratch/wisp/probes/158/r2b/zero-byte-census.sh`，读全文核过：**只往 stdout 打，不写仓里任何东西**），
锚点仍取被验程自己的 step-0 `bb08d1f`，读数落 `probes/158/accept-r1/census-at-my-anchor.txt`：

```
## 1  枚数=5（ef34118 8225cdc ee41e63 3122518 c7f638c；反向自证 8 枚编排者提交不在集合里，逐枚看 subject 前缀都是 ledger(A3xx)/派单）
## 2  名册去重后枚数=49 ;  带 --no-walk = 52 行 ;  不带 --no-walk = 4885 行
## 4  28 支冻结面 越界枚数 全 =0，且每支"该面已跟踪文件数"非零（PLAN.md 1 / docs/specs 14 / internal/risk 37 /
     panel 20 / agent 58 / observe 25 / thresholds.go 1 / llm/golden 3 / 含 golden 的路径 6 / testdata/golden 52 /
     allowlist 1 / slo-check.ps1 1 / tools/d22scan 6 / frontend 85 / design 30 / 151·154·156 证据件与票面各 1 …）
## 5  第二口径（字节级 numstat 打同一串正则）⇒ 只剩一枚 HIT: internal/tools/bridge.go
## 8  §5 那把尺的正控（同一串正则，打在已知动过冻结面的三枚真提交上）：
     12181923 命中=2  FIRE internal/risk/assessor_test.go, internal/risk/rules_scale.go
     00bbb76f 命中=3  FIRE internal/observe/{observer_cost_test.go,sampler.go,thresholds.go}   ← thresholds.go 本尊点得出
     cbbbdf85 命中=1  FIRE internal/agent/testdata/golden/spill-tool.sse                          ← golden 点得出
```

- `--no-walk` 那一发本程自己量到的是 **52 对 4885**（被验程在 3 枚号上量到 41 对 4870）。两次的**形状**一致：
  不带旗标就沿祖先链吞整仓历史 ⇒ **被验件那条"不许拿区间 diff 当名册尺"的理由成立**，本程也照做（本表全部普查按号集合）。
- **预测核**：被验件 §C.5 末与 §I.3 留下一条可复算预测（"枚数 3→4 时名册 40→48，§5 仍只剩 bridge.go，28 支仍全 0；再加一枚＝同样形状再走一遍"）。
  本程在 **5 枚**上复算：枚数 5、名册 **49**、§5 仍只剩 `bridge.go`、28 支仍全 0 ⇒ **预测兑现**。
  多出的那 1 枚＝`probes/158/r2b/msg-selffix.txt`（自纠那一枚的提交信息件），**没有一枚落在冻结面**。
- 完整性另一问：49 枚名册逐枚分类 ＝ 47 枚在 `probes/158/r2b/` ＋ 本证据件 ＋ `bridge.go`，
  **写面之外的路径 0 枚**（`sed` 三条排除后为空集，读数在 `…/census-at-my-anchor.txt` §4 那三段 WRITEFACE）。

### 1.3 尺三＝AC#5 那三道工具（不跑 vs 跑了空输出是两件事，一律打 rc）

```
$ $(go env GOPATH)/bin/gofumpt.exe --version ; rc=0        -> v0.12.0 (go1.27.1)
$ $(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm ; rc=0   -> 0 行（stderr 0 字节）
$ $(go env GOPATH)/bin/gofumpt.exe -l internal/tools/ cmd/wisp/ ; rc=0       -> 0 行
$ go vet ./internal/tools/ ./cmd/wisp/ ; rc=0                                -> 输出 0 字节
$ sh scripts/d22scan.sh ; rc=1                                               -> 唯一 finding：frontend/src/components/harness/right-rail.tsx:108 [emoji]
$ (cd tools/d22scan && go run . -root "$PWD/../..") ; rc=1                   -> 1 finding(s)，同一条；八枚作用域 examined 全非零
     bans #1-5 internal/=205  cmd/=23 · ban #6 frontend/=85 · ban #7 internal/tools/=18
     ban #8 design/=39  frontend/=85  internal/=414  cmd/=45            （原文 probes/158/accept-r1/d22scan-direct.txt）
```

⇒ 三道里两道**绿**（gofumpt／vet，且都是"跑了、rc=0、输出为空"而不是"没跑"），第三道 **rc=1**——
它是被验件 §D.5 那枚"未裁"的同一发。那一发怎么判，见 §5.2（派单 §3① 的三件，那里**分三件裁**）。
## 2. 格 AC#1 —— 判定：**成立**（并把被验程判"闭不了"的那一腿闭掉了）

票面 AC#1 三发：①未修码上摘 `bridge.go:559` 的 `b.OpenTask(dec.TaskID)` 必须量到**零枚红**；
②加一枚本包断言使同一变异在 `internal/tools` 就红；③复装后本包回全绿。

### 2.1 本程自己的三发读数（全部 `-overlay`，工作树一个字节没写；全程未与 `-cover*` 同用）

变异体与被变异输入都从 **`git cat-file blob c7f638c:…`** 生成（不从脏工作树取），
`diff` 自证只有一枚 hunk：`559d558 < b.OpenTask(dec.TaskID)`（原文 `probes/158/accept-r1/make-rig-a1.sh` 的输出）。

```
① 未装守卫（把 r1 那枚守卫用例对编译器替换成空 stub）＋ M1 变异：
   $ go test -count=1 -v -overlay=…/ovl/r1-stubmut.json ./internal/tools/      rc=0
     RUN=115 PASS=79 FAIL=0 SKIP=0                     ← **零枚红，量到了**
   对照（只拿掉守卫、不变异）：$ … -overlay=…/ovl/r2-stub.json                 rc=0
     RUN=115 PASS=79 FAIL=0 SKIP=0                     ← stub 自己不引入红，①的读数不是桩件的产物
② 装守卫（锚点上那枚已入库的用例）＋ 同一枚 M1：
   $ go test -count=1 -v -overlay=.scratch/wisp/probes/158/accept-r1/overlay-noopentask-a1.json ./internal/tools/   rc=1
     RUN=116 PASS=79 FAIL=1 SKIP=0
     --- FAIL: TestSensitiveReadAcrossTheBridgeOpensItsTaskScope   （日志第 55 行，与被验件 §1.3 同一枚行号）
③ 复装＝不带 overlay 再跑一次（工作树从未被写，`git diff c7f638c -- internal cmd tools scripts` 空）：
   $ go test -count=1 -v ./internal/tools/                                       rc=0
     RUN=116 PASS=80 FAIL=0 SKIP=0     ok 14.399s
```

**"overlay 真被编译器吃"这一证我没有转抄**——本程自己打了 canary：把 `:559` 改名成 `b.OpenTask_CANARY_A1` 后
`go build -overlay=…overlay-canary-a1.json ./internal/tools/` ⇒
`.\.scratch\wisp\probes\158\accept-r1\bridge.canary-a1.go:559:4: b.OpenTask_CANARY_A1 undefined`（rc=1），
不带 overlay 的 `go build ./internal/tools/` rc=0。⇒ ②那一红确实是这发变异造的，不是静默忽略。

### 2.2 先答派单 §2 那一句："这一发在未修码上响不响"

**响。** 而且今天仍能响：①那一发在被验件 §1.2 里是 115/79/0/0，本程在**锚点码**上用桩件覆盖守卫复算，
同为 **115/79/0/0**、rc=0、零枚红 ⇒ 判据不是恒真、也不是"修完之后再回头看就不响了"的装饰。

⇒ **这里我要翻被验件一处**：§E.4-2 写"要重跑那一腿，`git stash` 之外没有任何办法在共享树里临时删别人用例
⇒ 只能由非实现者验收程在**快照**里跑（§F.6）"。**这句是错的**：`go test -overlay` 把
`internal/tools/bridge_scope_open_ticket158_test.go` 对编译器替换成 `package tools` 一枚空桩即可，
**不删任何人的证据、不动工作树、不进快照**。被验程 §F.6 留给"下一位"的那条重路因此不必走：
AC#1 那一腿今天**已由本程闭掉**，读数是上面的 ①。

### 2.3 承重那句的两问：三枚判据各自有没有牙

被验件 §1.3 说这枚用例"三枚判据"（①桥的开账本 `b.scopes`；②`CloseTask` 审计行 `was_open=true`；③反向对照：`fs.list` 不得开账）。
本程不是读它怎么说，是**逐味摘掉再打变异**（`ovl/r3..r7`，全部现跑）：

| 发 | 变异 | 守卫 | 读数 | 说明 |
|---|---|---|---|---|
| r3 | **M2**＝把 `mark()` 的门改成"每个 task 都开账"（票面点名的退化形状） | 原样三味 | rc=1，**3 枚红**，本用例第一句响在 **`:95` 判据③**："非敏感源 fs.list 也开了 scope" | 判据③**会响** ⇒ 不是恒真腿 |
| r4 | 同一枚 M2 | **摘掉判据③** | rc=1，**仍 2 枚红**（`TestFSListSummarizesADirectory`、`TestAtomicWriteKillsMidWrite`） | ⇒ **判据③没有独立的牙**：我构造的那一发退化变异，摘了它本包照样红 |
| r5 | **M1**（摘 `b.OpenTask`） | 摘掉判据① | rc=1，本用例红在 `:122`（判据②那一句 `was_open=false`） | ①②对 M1 互为第二道 |
| r6 | 同一枚 M1 | 摘掉判据② | rc=1，本用例红在 `:112`（判据①那一句） | 同上 |
| r7 | **M3**＝`OpenTask()` 里把 `b.prov.OpenScope(taskID)` 换成 `_ = open`（桥的账照记，risk 层不注册） | 原样三味 | **internal/tools 116/80/0/0 rc=0 全绿**；再单跑 **cmd/wisp：144/84/0/0 rc=0 也全绿**（`run-cmdwisp-m3.log`，dll 注入走 shell 路径形，`=== RUN`=144≠0 ⇒ 真跑到了） | **本表新查出的残余**，见 §2.4 |

**结论的形状**：摘掉任意一味，**M1 与 M2 都仍会红** ⇒ 这三味是**纵深**，不是三枚各自单点承重。
被验件对判据③的那句"**挡住**一种看起来也满足了判据 1 的退化改法"（引自该文件 `:16-17` 的注释）
应读成"**它是最先响的那一道**"，不是"唯一一道"——**这一处算 MINOR 文案，不算退回**：
AC#1 的判据本体（"同一变异在本包就红"）成立，且被 r5/r6 证明任摘一味仍红。

### 2.4 残余（本程新造出来的一发，票 158 没要求、也不在该票判据里）

M3 这一发的意思是：**守卫量的是桥自己的那本账（`b.scopes`），不是 risk 层的 scope 注册。**
两包 228 枚作用域里的用例对 M3 全部沉默（读数见上表 r7）。机理在 `internal/risk/provenance.go`（冻结面，本程只读）：
`OpenScope`（`:339`）只是把 `p.scopes[id]` 注册成 `nil`，而 `Mark` 走 `p.scopes[scopeID] = append(...)`（`:394`）
自己就会把这一格建出来 ⇒ 在"先 OpenTask 后 Mark"这条路上，那枚注册**唯一**的可见差别是
`Mark` 里 `:386-388` 那句 "scope %q not open …; taint stored" 的日志会不会被印出来（M3 会让它每次都印）。
⇒ 所以这不是一枚安全洞（污点表本身同形），而是一枚**"有外部可见读数变化、但零枚红"**的形状，
和票 158 立项时那枚病（"唯一的守卫长在另一个包里"）**同族不同腿**。
**归口**：不在本票六格里——它要么落在 `internal/risk` 自己的测试面（该面本票冻结），要么由编排者记进台账再决定开票。
本程只登记、不修（`internal/risk/**` 不在我写面）。
## 3. 格 AC#2 —— 判定：**成立**

票面 AC#2 要四件：一枚子句（形状照 G1–G4：一条命令、名册基准、正控先跑）＋ 甲（未修码响不响）＋
乙（正控不响）＋ 丙（假阳性自拆，逐枚判读）。三发读数本程**逐发复算**在 §1.1，全部对上；这里只补 §1.1 没处的两件事。

**丙那一发的三处自拆，本程独立重打**（不是核它"写没写"，是核它写的读数在锚点上还成不成立）：

| 被验件登记的假阳性 | 本程现量 | 判定 |
|---|---|---|
| FP-1 成对判据第一版不剔注释 ⇒ 被**注释**点中 | 同一射程去掉剔注释那一行：负一负 **5→6**，多出的正是 `internal/tools/bridge_scope_open_ticket158_test.go` | **真事**，且那一行的修复**承重**（改变可见读数） |
| FP-2 正控第一版 pathspec 写成 `internal/tools/**/*.go` ⇒ 装饰腿 | `$ git grep -lEw OpenScope c7f638c -- internal/tools/**/*.go` ⇒ **rc=1、零命中**；同尺 `internal/tools/*.go` ⇒ rc=0 命中 `bridge.go` | **真事**：那一版正控**根本产不出读数**，被验件把它换掉是对的 |
| FP-3 每文件计数被 `<锚>:` 前缀污染 ⇒ `开=0 close=0` 自相矛盾 | 锚点上现量主尺射程：`panel_assets.go` 开方调用点＝**1**（与点名同帧，不再自相矛盾） | **真事**，修复后可见 |

**这一发买的是哪一档**（票面 AC#2 那句"拿不到『今天响』是可接受答案，但要按防忘记／防回归两档写清"）：
甲那发今天**响**（未成对＝1），所以那枚 fallback 今天不必动用；被验件 §2.2 末仍写清了当前档＝**防忘记**、
`panel_assets.go` 拿到收尾后同一把尺自动换成**防回归**。**这一格对 AC#3 的关系要单列一句，见 §8.③**——
它不构成"同一枚事实抵成两格"。

**一处 MINOR（不据此退回）**：主尺今天响（该腿 `# rc=1`），但 `gate-clauses.sh` **全文 rc=0**——
脚本没有聚合位，退出码就是最后那行附录 `git grep` 的码。⇒ "响"只活在 stdout 里，
`bash gate-clauses.sh && echo green` 今天**照样绿**。这一形 G1–G4 早就有（AC#2 明令"形状照 G1–G4"，
所以照做不是越界），但**AC#3 选的正是"让它响着"**那一支，"响着"的音量因此是本票的事 ⇒ 记在 §4。
