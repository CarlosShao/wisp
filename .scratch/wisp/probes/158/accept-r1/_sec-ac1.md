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
