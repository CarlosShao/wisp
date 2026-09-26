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
## 4. 格 AC#3 —— 判定：**成立**（选了支 B"让它响着并登记"，本程核的是它有没有拿修辞抵账）

票面 AC#3 的三件"不算收"逐件量（不采信被验件的自查表，本程自己重打）：

| 那件"不算收" | 本程现量 | 判定 |
|---|---|---|
| ① 只在注释里写"这是 CLI、进程退出就干净了" | `git grep -nE 'Defer\|disposal\|DisposalScope\|os\.Exit' c7f638c -- cmd/wisp/panel_assets.go` ⇒ **rc=1、零命中**；`cmd/wisp/main.go` 里 `case "panel-assets"` 在 **`:103`**、`os.Exit(cmdPanelAssets(args[1:]))` 在 **`:105`**（两处行号与被验件 §3.1 的 R-3/R-4 **一字对上**）⇒ "有界"的凭据确实长在**另一枚文件的分发行**上，开方文件里什么可机读形状都没有 | 它没写那句空话，它写的是**一枚尺＋两条熄火路径** |
| ② 加进白名单而不给可机读判据 | §1.2 名册普查：`tools/d22scan/**`、`allowlist.txt` 越界 **0** 枚；G5 的 `pair()` 里没有任何文件名豁免（读全文核过） | 没有白名单 |
| ③ 造一枚今天不响也永远不会响的恒真检 | §1.1 三腿：主尺**今天响**（1）、正控名册**非空**且未成对 0、负一负 5；§1.1-2 那发"把 internal/risk 放回射程"复算 ⇒ 判语不变 | 不恒真 |

**熄火条件**（被验件 §3.3 那两条）本程复核为**可复算**：读 `bash .scratch/wisp/probes/154/gate-clauses.sh` 的 `## G5` 段"未成对枚数"，
路径 1＝`cmd/wisp/panel_assets.go` 同文件出现 `CloseScope`（或委托桥的 `CloseTask`，`pair()` 会把委托 hint 印出来）——
本程现量该文件**两族都零枚** ⇒ 今天确实没走通；路径 2＝owner 给一条**谓词**而非文件名豁免（＝`Q-56` 那条线，agent 不能单方面选）。
⇒ 支 B 的内容不是"我先不修"，是"**这条命令每次跑都会点名，除非这两条形之一兑现**"。**成立。**

**音量这一处要记（MINOR，不据此退回）**：见 §3 末——G5 那一腿 rc=1，但门本体**全文 rc=0**（`bash gate-clauses.sh; echo $?` ⇒ 0）。
所以"响着"今天是**打印级**的响，不是**退出码级**的响。形状照 G1–G4 是 AC#2 明令的，故不判越界；
但请编排者在归口时把"要不要给门本体一枚聚合 rc"记成独立一问——**它决定 AC#3 这一支的响对不跑门的人是否存在**（被验件 r1 §5-2 自己就登记了"门没接 CI"这一条，两件事同向）。

### 4.3 单列一行：那一处**由编排者发起的注释面扩权**（它不属于六枚框中的任何一枚）

被验件 §H-③ 明确请求"别并进 AC#2"。本程按独立一行裁：**成立**。

- 批准射程＝`internal/tools/bridge.go` 那一枚门指路句；本程自己重打三件自证：
  `git show ef34118 --numstat` ＝ **7 加 / 3 删**、**非 `//` 开头的改动行＝0 枚**、
  改动行里出现 `/*` 或 `*/` 的枚数＝**0**（票面 `:81` 点名的"把代码包起来"那一形没有）。
- **它闭合的是不是真缺陷**：是。锚点上现读 `:663-670` 那八行——新句同时点名
  `.scratch/wisp/probes/154/gate-clauses.sh`（可跑生成器，路径存在）、G1/G1b 管桥这条腿、G5 管**不过桥**那一腿，
  并把两份 write-up 分开指。被验件 §B 的②号理由我核过：`docs/evidence/s1/154-host-id-never-closed-r1.md` 的 `### 2.3` 一节里
  确实写着 G1/G1b/G2/G3/G4（本程 `awk` 按节取，命中 12 行）⇒ 那枚指针**今天不再骗人**。
- 新句里那句事实断言"Today no production code dispatches on the bridge except loop.go"本程也打了尺：
  G1b 那一腿在锚点 `# rc=1`（生产码 `.Execute(` 排除 `loop.go`/`bridge.go` 后零命中）⇒ **注释没说过头话**。
- r2 在 §B 末登记的那枚仪器坑（未提交时门本体看不见注释 ⇒ 假等价）本程复算方式不同但同判：
  本程所有读数一律 `git cat-file blob c7f638c:…`／`git grep … c7f638c`，从不信工作树。
## 5. 格 AC#4 —— 判定：**成立**（契约轴零字节；本程把自己的号集合也套进同一把尺重跑）

主尺读数全在 §1.2（**28 支冻结面越界全 0**、每支"该面已跟踪文件数"非零、名册 49 枚、
`--no-walk` 52 对 4885、§8 正控三枚全命中）。这里补 §1.2 没处的四件：

1. **"零字节"与"零枚文件"不是一件事，本程把两样都量了。**
   `git log --no-walk --pretty=tformat: --name-only <五枚号>` 的 49 枚名册逐枚分类后：
   `probes/158/r2b/` 47 枚 ＋ 本票证据件 1 枚 ＋ `internal/tools/bridge.go` 1 枚，
   排除这三支后**剩空集**（`rc=1`，原文见本表 commit 序列旁那条命令）⇒ 冻结面上**连一枚文件名都没有**，
   更不必谈字节。另：`git show` 逐枚 numstat（本程自跑，原文在上）里删除列只出现在
   `internal/tools/bridge.go`（3 行，全注释）与证据件的自纠（0 删）。
2. **追加性（被验件反复主张的那一条）本程按范围量**：
   `git diff --numstat 6f127f1 c7f638c -- docs/evidence/s1/158-gate-scope-blind-spot-r1.md` ＝ **652 加 / 0 删**
   ⇒ r1 的第 0–7 节确实一字未改写；逐枚看四枚证据件提交分别是 282/0、80/0、236/0、54/0，
   `c7f638c` 那枚"自纠"的删除列＝**0** ⇒ 它标题里那句"已入库那一行一字未改"**成立**，不是修辞。
3. **§C.6 那两处"被自己扫到"的计数，本程按提交逐枚复算，对得上且现在会漂**：
   `t.Skip` 的**新增行**在 r2 量它自己那三枚号时＝**2**（ef34118 0 ＋ 8225cdc 1 ＋ ee41e63 1），
   本程在五枚号上量到 **4**（AC#4 那一枚又说了两遍"没有 t.Skip"）。`DEFERRED(D-` 同理。
   ⇒ 这不是它报错了，是**这把尺的口径会随集合变大**；被验件把"不拆文件会把否认读成做了"写进 §E.3/§C.6 是对的。
   要紧的那一半本程单独钉了：**`*.go` 里新增的 `DEFERRED(D-` ＝ 0 枚**
   （`git log --no-walk -p … -- '*.go' | grep -cE '^\+.*DEFERRED\(D-'` ⇒ 0）⇒
   AGENTS.md §1.1 那条"DEFERRED(D-xx) 与 SPEC-12 §5 登记表 1:1"没有被本程添乱。
4. **这把尺的覆盖面在哪**：`FROZEN` 那 28 支里 `internal/tools/**` 只以 `bridge.go` 一枚出现（因为测试面是**写面**）。
   ⇒ 单看 §4 的逐支账，它抓不到"往 `internal/tools/` 新加一枚生产文件"这种形状；
   真正兜住它的是 §2 的**全名册**（49 枚逐枚归类、空集自证）。这一条是**这把尺设计对的**的证据，
   不是缺口——但下一位如果只抄 §4 那张表当"零字节"的全部，会以为逐支账自成闭环。本程因此把 49 枚逐枚归类重跑了一遍。

**另记一件与本程有关的事**（免得下一位复算时被误导）：那把尺的 §6「工作树侧」今天会列出
**我的** `probes/158/accept-r1/**`（未跟踪）。那是验收程的台件，不是 r2 的越界；
复算 AC#4 时请只读 §1–§5、§8，§6 那一行按"当前树里谁在动"理解。
## 6. 格 AC#5 —— 判定：**四子项成立 ＋ 一子项未裁 ⇒ 这一枚框本程不勾**（最小闭合集合在 §8.①(c)）

票面 AC#5 五子项：①两包改前改后各一次、四数之外名册两向 `comm`；②`cmd/wisp` 要 dll 注入（CI 同形）；
③`gofumpt -l . tools/d22scan tools/mockllm` 空（版本自读）；④`go vet` 两包空；⑤`sh scripts/d22scan.sh` **rc=0** 且各作用域 `examined N` 非零。

**① 四数与名册（本程自己重跑"改后"，并和被验程的两份名册做字节级对照）**

```
$ go test -count=1 -v ./internal/tools/            rc=0   RUN=116 PASS=80 FAIL=0 SKIP=0     ok 14.399s
$ bash scripts/wisp-cli-tests.sh                   rc=0   RUN=144 PASS=84 FAIL=0 SKIP=0
$ 名册尺（^[[:space:]]*--- (PASS|FAIL|SKIP): 全名、去时长、LC_ALL=C sort）
   本程 cli 名册 144 行  md5 = c77310b70b17def6102f661ad543ca34
   被验件 roster-cli-before.txt  = c77310b70b17def6102f661ad543ca34
   被验件 roster-cli-after.txt    = c77310b70b17def6102f661ad543ca34
   本程 tools 名册 116 行  md5 = 31e18e28f1f03bd5ba387edf99c3ae17
   被验件 roster-tools-before.txt / -after.txt = 同一枚 31e18e28…
```

⇒ 三枚来源不同的名册**字节级同一枚值**（我这一跑是独立进程、独立 PATH、独立 tip）。
它同时把被验件 §I.2 那两条"改前＝改后、`cmp rc=0`、两向 `comm` 皆空"的自纠坐实了：
我算不出差集，因为**没有差集**；而且我用的正是它自纠后那条尺（`LC_ALL=C` 一路到底 ＋ 一把字节级尺在旁边）。

**② dll／PATH 那一坑（本程自己一正一负都打了）**：正控＝上面那一跑 `=== RUN`=144≠0；
负控＝本程**故意**不注入 dll 直接 `go test -count=1 -v ./cmd/wisp/` ⇒ `RUN=0`（原文
`probes/158/accept-r1/cli-negative-nodll.log`（字面两行＝`exit status 0xc0000135` 与 `FAIL`，`=== RUN` 计数＝**0**）
⇒ "判根本没跑到只认 `=== RUN` 枚数＝0"这一条**独立成立**，被验件 §D.4 那两发不是修辞。

**③④ gofumpt / vet**：读数在 §1.3。版本是**本程自己 `--version` 现读**的
`v0.12.0 (go1.27.1)`（与被验件 §A 末、§D.5 同值 ⇒ 它没抄旧版本），全仓 `-l . tools/d22scan tools/mockllm` **0 行 rc=0**，
`go vet ./internal/tools/ ./cmd/wisp/` **0 字节 rc=0**。两枚都是"跑了且空"，不是"没跑"。

**⑤ `sh scripts/d22scan.sh`＝本程唯一不勾这一格的原因。**
票面这一子项写的是 **rc=0** 且 `examined N` 非零。本程现量：**rc=1**（后文三口径证明它是**被验那一版的真红**，
不是脏树假象——见 §8.①(a)），`examined N` 那一半**非零**（八枚作用域全在 §1.3 列出）。
⇒ **一格两半，一半今天不成立。** 被验程判它"未裁·受阻"，并且写明"本程不许它变绿"（§E.2 那行）。
**这个处置本程支持**，理由不是"它可怜"，是三条机器读数：那枚字形所在文件的**工作树字节＝锚点 blob**（§8.①(a) 的 hash 对照）、
`frontend/**` 与 `tools/d22scan/**` 与 `allowlist.txt` 全在它的冻结名单里（§1.2 越界 0 枚反过来证明了它确实没能力修）、
它没有放宽断言/加 `t.Skip`/动阈值（§5 第 3 件：`*.go` 里新增 `t.Skip`/`DEFERRED(D-` 全 0）。

**AC#5 的档位因此是**：①②③④ **成立**；⑤ **未裁**（不是"不过"，也不是"过"）；
**整格按票面文字不成立**（票面要求 rc=0），但**责任面不在本程/实现程手里** ⇒ 框**不勾**，
残余**归口**与最小闭合集合在 **§8.①(b)(c)**，本表 §9 的勾法里 AC#5 那一枚保持 `- [ ]`。

**顺带把 §D.3 那枚"票面数过期"的判语复算**：被验件说票面 AC#5 第③条那句"声明 82／本机跑到 79"今天应为
**87 声明／84 编译／84 跑到**，差集 3 枚全在 `secret_dataroot_119b_test.go` 的 `//go:build !windows` 之下。
本程现量：CI 同形那一跑 `=== RUN`=**144**、顶层 PASS=**84**（与之一致），且我**没跑** `go test -list`（那一步与它同形，
不重复占机器）；`84 跑到`这一枚与我独立一跑对上 ⇒ **那句"数过期、道理成立"判成立**。
道理那一半本程也照做：**我没有**把"这两包全绿"写成"全仓无影响"（`go test ./...` 本程一枚没跑，见 §10）。
## 7. 格 AC#6 —— 判定：**成立**（两向本程都自己数过；两处 MINOR 是"口径话说过头"，不是读数错）

**正向（一枚框 → 哪一程 → 那一发在未修码上响不响）**：票面在锚点上 **84 行、`- [ ]` 6 枚、`- [x]` 0 枚**，
六枚框所在行 `24/28/32/34/38/43`（本程自己 `grep -n '^- \[ \] \*\*AC#'`，与被验件 §E.0 那六个数**一字不差**）。
"这一发在未修码上响不响"这一问，本程不是复述而是重打：AC#1 在 §2.2（今日重造变异 **116/79/1/0**，
红名与被验件 §1.3 同一枚、日志同一枚行号 **55**）；AC#2/AC#3 在 §1.1（主尺 1／正控 0 且名册非空／负一负 5）；
AC#4 的形状不是"响不响"而是"有没有往冻结面落字节"，读数在 §1.2（28 支全 0）。

**票面点名的三枚行号，本程在锚点全部复量**：

| 票面引的 | 本程现量（`git cat-file blob c7f638c:…` / `git grep … c7f638c`） | 漂没漂 |
|---|---|---|
| `bridge.go:559` | `	b.OpenTask(dec.TaskID)` | **未漂** |
| `panel_assets.go:232` | `	prov.OpenScope(taintSourceScopeID)` | **未漂**（G5 今天仍只点它） |
| `task_scope_close_151_test.go:109/:113` | `t.Errorf` 在 **62／103／109／113** 四枚（票面只点后两枚） | **未漂**，被验件 §0-P3 的自纠方向对 |

**反向（本程动过的每一类 → 哪一枚框；孤儿工作清点）**——本程按号集合重算，不采信它的表：

- 五枚号的全部名册路径落三支写面之外**一枚都没有**（把 `.github|tools/d22scan|scripts|docs/specs|docs/PLAN|internal/risk|
  internal/panel|internal/agent|internal/observe|frontend|design` 十类当筛子打上去 ⇒ **无输出**）。
  ⇒ §E.3 末那张"本票不解决的事四条逐项自查"里"不把门接进 CI"、"不碰 frontend/design"两行**成立**。
- `Q-56` 那一行它算得对，本程逐枚归名复算：五枚号合计新增 **8** 行含 `Q-56`，分布＝
  证据件 2（都是"本程没答 Q-56"这两句自己）、`msg-ac6.txt` 1（提交信息件）、`probes/158/r2b/bridge.{canary-r2,noopentask-r2}.go` 各 1
  （生产注释整份抄成的探针副本）、三份 `gate-*.txt` 各 1（被抄进去的节标题 `## G3 Q-56 那一支落地…`）。
  **`Go` 生产码里 0 枚** ⇒ "不答 Q-56"成立；被验件在它自己三枚号上报的那 **5** 枚＝本程这里的
  探针副本 2 ＋ 读数件 3 ⇒ **同一个数、两把尺对得上**，它那句"不拆文件会把'抄了什么'读成'答了什么'"是本格最有价值的一条口径登记。
- **票面框 6 枚实现方一枚没勾**：`.scratch/wisp/issues/` 在五枚号的名册里 **0 枚**（§1.2 逐支账 `本票票面 越界=0`）⇒ 对。

**两处 MINOR（都属"口径话说过头"，读数不受影响，本程因此仍判成立）**：

1. §E.0 末那句"⇒ 本程**一枚 `sed -n 'A,Bp'` 都没用**"说过头了：同一份交付物 §B 里就有
   "改后原文（现量 `sed -n '663,669p' internal/tools/bridge.go`）"——那是一枚**行窗**。
   本程按锚点复量那段注释：663–670 才是它抄在下面的那 **8** 行，`663,669p` 只有 7 行（末行 `// \`git grep\` away…` 会缺）。
   ⇒ **它抄的原文逐字是对的**（本程把 8 行与 blob 逐行比过，全同），**但它给的那条复算命令少一行**。
   这条正好撞在它自己 §I 那节"引原文就是要能逐字复算"的判据上——只是这次是**命令**复算不出**引文**，不是引文错。
2. §E.2 里 AC#2 那一行写"锚点 `192ad56`（本程改前）与 `ef34118`（本程改后）两跑判语一字未变"——
   本程不按它的锚点，改按**我的**锚点 `c7f638c` 重跑（§1.1）：判语同样是"主尺 1／正控 0／负一负 5"，
   但**名册里 `bridge.go` 的合方行在我这一跑是 `:695`**（它 §2.3 引的是 `:691`）。
   两件事不冲突（它引的是注释面之前那一版），**可被验件正文里没有一句把这枚漂移标在 §2.3 旁边**，
   只有 §B 末藏着一句"名册里唯一变的是行号 691→695"。下一位拿 §2.3 的 `:691` 去 `sed` 会读到别的东西
   ——与本表 §1.1 末那处"门本体脚本注释里仍写 `:691`"是同一枚病，**登记一次即可**。

**AC#6 自己要求的那一半（"没判的明写未裁＋最小闭合集合"）**：被验件 §E.4 列了三处。本程逐处裁定——
第 1 处（`d22scan` 子项）**维持未裁**，本程把它的三口径复核做在 §8.①；
第 2 处（AC#1 的"未修码零枚红"那一腿）**已被本程闭掉**，且它给的"必须进快照"那句是错的（§2.2）；
第 3 处（票面框不自勾）**成立**，勾法由本表 §9 给。
