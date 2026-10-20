# 票 230 `AC#2` 落地腿 `230-r1` 交件主报（★纯注释面·可执行字节变动＝0）

- 起手锚＝本腿第一笔 `git log -1` 现认＝`c16ff47ec465749f1977efa94a2b024cc8931e27`
  （`票 352 落地腿 352-r1 笔④…`，16:29:06 +08）；分支＝`dev`（`git branch --show-current` 现跑）。
- 本腿提交＝两笔：`c489eacc`（笔①）＋`c73df208`（笔②）。主题逐字以「票 230 落地腿 230-r1」起头，
  名册尺（锚定行首）＝`git log --all --grep='^票 230 落地腿 230-r1' --format='%h'`＝**2 笔**，⛔ 换形。

## 1. 定位尺＝内容锚（⛔ 用绝对行号；票面写的 `:247` 已漂）

尺（起手侧，落笔前，对 `HEAD` 跑）：

    git grep -n 'only reader of q.alias' HEAD -- internal/agent/approval/queue.go
    => HEAD:internal/agent/approval/queue.go:262:// stays an unknown correlation id. This function is the only reader of q.alias

同一枚谎话在同文件有**两**处（第二把尺把两处一起抓到）：

    git grep -n 'only reader' HEAD -- internal/agent/approval/queue.go
    => :87   // to (qitem.names) plus the item that owns it. Its only reader is
    => :262  // stays an unknown correlation id. This function is the only reader of q.alias

⇒ 笔①改 `:262` 那枚函数文档（`lookupForRefusalLocked`，定义在 `:265`）；笔②改 `:87` 那枚字段文档
（`Queue.alias`，声明在 `:91`）。两枚都在 `internal/agent/approval/queue.go`，⛔ 碰第二枚文件。

## 2.「其余读者」名册（逐枚 `file:line`，尺＝`git grep -n 'q\.alias'`）

起手锚 `c16ff47e` 上的全部命中（6 枚，⛔ 含注释自身以外零枚别处）：

    queue.go:212  set := q.alias[n]              <- indexLocked   体内（定义 :210）＝读＋建桶
    queue.go:215  q.alias[n] = set               <- indexLocked   体内＝写（新键挂桶）
    queue.go:225  set := q.alias[n]              <- unindexLocked 体内（定义 :223）＝读
    queue.go:231  delete(q.alias, n)             <- unindexLocked 体内＝删空桶
    queue.go:262  (注释自己那句假话)
    queue.go:272  if set := q.alias[corr]; …     <- lookupForRefusalLocked 体内＝唯一"把回答的名字解析成条目"的读

⇒ **同文件另有读者＝两枚＝`indexLocked`（`queue.go:210`）⇄`unindexLocked`（`queue.go:223`）**；
票面 `:197`/`:210`/`:216` 那三枚号今天漂到 `:212`/`:225`/`:231`（票面 `:13` 已自己具名"别信这里的行号"）。
全仓复尺＝`git grep -n '\.alias\[' <锚> -- internal/ cmd/ tools/`（⛔ 只看本包）＝**4 枚**命中：
起手锚 `c16ff47e` 上＝`:212`/`:215`/`:225`/`:272`（两次读桶＋一次挂桶＋解析那一读），
HEAD 上＝`:215`/`:218`/`:228`/`:283`（同一批字节被我的注释行数顶下去）；
删除那枚 `delete(q.alias, n)` 不含 `[`，由上面那把 `q\.alias` 尺抓到＝共 6 枚命中。
零枚在 `queue.go` 之外 ⇒ 字段是私有的，别包读不到（⛔ 新导出名，本腿也没造）。

调用侧那半句（"Queue.reject is its only caller"）**今天为真**，尺＝

    git grep -n 'lookupForRefusalLocked' HEAD -- internal/  => 唯一非注释调用点 queue.go:493
    定义 queue.go:491  func (q *Queue) reject(corr, reason string) error

⇒ 改法走票面 `AC#2` 的**第一支**（把"唯一读者"改成说实话）：限定"唯一"的射程＝唯一*解析回答*的读者，
点名其余两枚维护者并具名它们为什么不构成解析（各自只搬传进来的那一枚 `*qitem` 自己的桶），
已经变真的那一半（唯一调用者＋allow 侧走 `lookupForAllowLocked`、其体内只 `return q.byID[corr]`）
原样保住，⛔ 把一条被 87/97 钉住的安全陈述洗成含糊。

## 3. 改后注释原文＝见件 `01-rulers.txt` 第 0 节（逐字，带行号，取自 `git cat-file blob HEAD:…`）

## 4. 判据尺（可执行字节变动＝0）——本腿自己算的三把，全部在 `01-rulers.txt`

1. `git show -U0 <commit>` 的每枚 `^[+-]` 变动行都以 `//` 开头：笔①＝14 枚全 `//`；笔②＝7 枚全 `//`；
   起手锚→HEAD 累计＝21 枚全 `//`；非注释变动行＝**0 枚／0 枚／0 枚**。
2. 剥掉整行注释后的正文（`grep -vE '^[[:space:]]*//'`）起手锚 ⇄ HEAD：421 行 ⇄ 421 行、
   `diff` rc=0／0 字节、两侧 sha256 逐字同串 `fadaf8d1…5089f`。⇒ 这就是"可执行字节变动＝0"的量证。
3. `git show --name-only` 逐枚只含 `internal/agent/approval/queue.go` 一枚（越界 0）；
   `--numstat`＝笔① `11 3`、笔② `5 2`，全部落在那一枚文件的注释面上。

⛔ 把 `Queue.reject` 分支改宽／⛔ 动 `Queue` 任何授权判定／⛔ 放宽任何断言／⛔ 翻任何框／⛔ push——全部零字节。

## 5. 门禁（只在共享工作树；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

    GOFLAGS= go build ./...                     rc=0
    sh scripts/d22scan.sh                       rc=0  clean（ban#8 internal/ 603 枚含注释与 _test.go）
    gofmt   -l internal/agent/approval/queue.go rc=0  列 0 枚
    gofumpt -l internal/agent/approval/queue.go rc=0  列 0 枚（全路径 /d/work/base/gopath/bin/gofumpt.exe）
    裸 gofumpt                                  rc=127 具名上报＝不在 PATH，⛔ 当"跳过"
    go vet ./internal/agent/...                 rc=0
    sh scripts/check-path-length-budget.sh      rc=0  VERDICT GREEN
（逐字输出＋并排两把的仓内脏件名册见 `02-gates.txt`；⛔ 洗别人已知脏件：那 4 枚 tracked 脏件都是 CRLF 存量，
本腿零字节，`git ls-files --eol` 逐枚具名。）

## 6. 颜色（定向 `./internal/agent/...`，⛔ 整包 `./cmd/wisp`、⛔ `./...`、⛔ `./internal/winsec/`）

两棵仓外导出树（唯一差异＝`queue.go`，尺＝`diff -rq` 只报那一枚）：
base＝`3eeaa2e2`（笔①的父＝编排者 A1064 那笔）queue.go blob `ce63814b…`；
now＝`c73df208`（HEAD）blob `890e1245…`。
⚠ 可比性另补一把尺：起手锚 `c16ff47e` ⇄ base `3eeaa2e2` 的 `queue.go` blob **同哈希**
`66fec7ae375344d7bc741270b42c4ec5050cf41f`（＝中间那两笔编排者的笔⛔ 碰过这枚文件，
我的 base 与我起手看见的字节逐字相同）；`merge-base --is-ancestor c16ff47e c489eacc` rc=0。

    四数（RUN / PASS_TOP / FAIL_TOP / SKIP，尺＝^=== RUN / ^--- PASS / ^--- FAIL / ^--- SKIP）
    base = 212 / 168 / 0 / 1      now = 212 / 168 / 0 / 1      TEST_RC 两侧＝0
    顶层结果名册 symdiff＝0 字节；两侧 sha256 同串 f4ac5a80…8ebb（169 枚逐名）
    被注释点名的两枚用例在两侧都绿：TestAnAliasCanNeverBuyAnAllow / TestEveryRefusalRouteOnAnUnknownEntryStillRefuses
    唯一 SKIP＝TestDefaultDeadlineWallClockMeasurement（两侧同枚＝存量弃权，⛔ 当闭环）

逐字读数与归一化作差（剥时间戳后仍有 3 行差异＝`t.TempDir` 的随机尾号，⛔ 要求逐字相同，
按公共块"只认到达数分子＋复现形状"那味具名登记）见 `03-colour.txt`。

## 7. 跨腿耦合的自证

- ⛔ 动 `internal/winsec/**`／`cmd/wisp/**`：两枚提交的 `--name-only` 各只含 `queue.go` 一枚（越界 0）。
- `3eeaa2e2`→`HEAD` 之间别腿在我包内的落字＝零：`git log c16ff47e..HEAD` 逐枚 `--name-only` 复验，
  `352-r1` 笔⑤（`06aa2abf`）只动它自己的 `.scratch/wisp/probes/352/r1/**`（我⛔ 用 `A..B` 区间数"谁动了某文件"）。
- `queue.go` 与票 220 同包＝本腿是该包今天唯一写腿（票 220 今天⛔ 派）；⛔ `-w` 扫别人文件、⛔ 攒一批。

## 8. 台面自纠（具名报编排者，一条）

`cygpath -wp` 整条 PATH 那味本腿没踩：定向色只在仓外导出树跑 `go test`，全程**不动 PATH**
（`internal/agent/...` 的导入图里没有 sherpa／CGO 那两味 ⇒ 不需要 DLL 搜索路径），
所以 A1057④ 那形"零枚 `--- FAIL` 的整包死"在本趟不构成风险；若下一波有人要在仓外树跑 `.exe` 取 DLL，
仍按派单那一味只在那时转。

## 9. 未做／不自盖（⛔ 把没做的当做了）

- ⛔ 翻框：票 230 终量仍 0 勾／5 未勾（尺见 `01-rulers.txt` 第 4 节）。`AC#2` 归**非实现者**裁。
- ⛔ 顺手做 `AC#1`（另派 `115-v1`）／⛔ 碰 `AC#3`／`AC#4`／`AC#5`（`AC#5` 要整包色＝归编排者裁后排）。
- 一处**具名请编排者裁**（⛔ 我自扩射程）：`gate.go:661-662` 那句"The alias index is not in this path"
  与本枚字段注释同族、今天**为真**（`sessionSubject`⇄`allowScoped` 都只走 `lookupForAllowLocked`，
  现读 `queue.go:468`/`:396`），本腿⛔ 动它一个字节。另 `approval.go:497` 那句"only callers"与本票无关、也没动。
