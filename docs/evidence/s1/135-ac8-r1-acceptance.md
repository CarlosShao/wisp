# 票 135 AC#8（家族第 7 发 `M-G`）**非实现者**终裁 —— `acceptor-ticket135-ac8-r1`

**执行方**：`acceptor-ticket135-ac8-r1`（**裁决方，不是实现方**；两枚实现程 `fa35557`／`8369b24` 系的自述在本表里一律按"被审对象"处理，不当依据）
**被验版本**：commit **`c8967b8`**（`c8967b8bb90f93dc0d7e14a8825ce714fdcc8364`；`git cat-file -t` 现核＝`commit`，且＝本地 `dev`＝`cnb/dev`，`rev-list --count @{u}..HEAD`＝**0**）
**判据来源（不是被审对象）**：
- 票面 `.scratch/wisp/issues/135-…-16-16-green.md` **AC#8** 那一格（09-24 00:0x 编排者追加）＋该票最后两条编排者 log（断点登记／"本格没翻勾、等你"）
- 来源判据 `docs/evidence/s1/133-ac2-r2-acceptance.md` §1.11／§1.12／§3.4（`R-133-9`＝`M-G` 出处，票 133 AC#2 第二次退回）
- 票 133 面 `## AC#2 第二格复算` 那节的 **C1／C2／C3**
**被审对象（不作依据）**：`docs/evidence/s1/135-ac8-mg-impl.md`（§0–§9，736 行）
**本程地界**：一枚新建证据文件（本文件）。**生产码零接触、测试码零接触、那把尺零改动**（§1.3 字节为证）。**本程不翻任何勾、不动票面。**
**取件形状**：所有读数落在仓外归档树 `/d/tmp/wisp135ac8r1-tree` 及其逐发新副本；**仓内零 `go build`／零 `go test`／零 worktree／零 checkout／零 amend/reset/stash/clean**。

---

## §0 争用闸门（每一批读数之前一次，逐轮原样）

### 0.1 本程开工头一发（`Bash` 工具第一次调用）

⚠ **仪器缺陷先登记**：这一发本程没在输出里打 `date`（把三行包在一条命令里时漏了时间戳）。同一状态在 12:08:01 的第二发里带了时刻，故本条仍能定位、只是**它自己的时刻是自述不了**。形状记下来：**闸门输出必须自带 `date`**，否则那一发的"何时看过"只能引用别人的钟。

```
--- gate 1: runner ---
Runner.Listener.exe          43288 Console                    1     98,396 K
Runner.Worker.exe            39500 Console                    1     96,832 K
rc_runner=0
--- gate 2: go toolchain ---
rc_go=1                                  ← 无输出，grep 退出码 1
--- gate 3: gh runs ---
[{"conclusion":"","createdAt":"2026-09-24T04:01:58Z","status":"in_progress"},
 {"conclusion":"failure","createdAt":"2026-09-24T03:01:39Z","status":"completed"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:50:36Z","status":"completed"}]
```

⇒ **命中**：一枚 `Runner.Worker.exe`（PID 39500）＋一枚 `in_progress` 的 run。**这一批发不作数。**

### 0.2 第二发（12:08:01，取件之后、写作之前）

```
at=2026-09-24 12:08:01 +0800
--1--
Runner.Listener.exe          43288 Console                    1     98,752 K
--2--
（无输出）
--3--
[{"conclusion":"","createdAt":"2026-09-24T04:01:58Z","status":"in_progress"},
 {"conclusion":"failure","createdAt":"2026-09-24T03:01:39Z","status":"completed"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:50:36Z","status":"completed"}]
```

⇒ `Runner.Worker.exe` 已退，Listener 常驻；**第 3 行那枚 `in_progress`（`04:01:58Z`＝`12:01:58 +08`）还在** ⇒ 按派单"任一命中即不发"**仍判 BUSY**。

### 0.3 第三发（12:09:50，包成 `gate.sh` 之后的第一发；**这一发不是读数批**）

```
---- gate at=2026-09-24 12:09:50 +0800 ----
[1] tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'
Runner.Listener.exe          43288 Console                    1     98,768 K
[2] tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'
(none)
[3] gh run list --limit 3 --json status,conclusion,createdAt
[{"conclusion":"","createdAt":"2026-09-24T04:01:58Z","status":"in_progress"}, …]
VERDICT worker-hits=0 listener-hits=1 toolchain-hits=0 in_progress_runs=1
VERDICT=BUSY
```

这一发只做了一件事：`go version`（`go1.27.1 windows/amd64`）与 §1 的哈希核对。**没有任何红绿读数取自争用窗口。**

### 0.4 台件与判读规则（本程自取，不沿用实现方的 `wisp135mg-r2-gate`）

- `/d/tmp/wisp135ac8r1-gate/gate.sh`：把派单那三行原样打印，并给 `VERDICT worker-hits=… listener-hits=… toolchain-hits=… in_progress_runs=…`。
- `/d/tmp/wisp135ac8r1-gate/poll.sh`：**有上界的轮询**（`INTERVAL=45` × `MAXROUNDS=26` ≈ 19 分钟一枚上界；逐轮全文进 `gate-rounds.log`；到点上界仍忙则 `exit 3` 并照实报）。**不是 `time.Sleep` 糊窗口**——等的是"三行全空"这个条件，且带次数与时长上界。
- 判读口径（本程自己定的，写清楚，别和实现方的口径混）：
  - `Runner.Worker.exe` 命中＝**真有 job 落在本机**；`Runner.Listener.exe` 常驻**不单独算命中**（它是等活的进程），但**照打不误**。
  - `gh run list` 里任一 `in_progress`＝**命中**（本程不赌"它大概不在这台机器上"）。
  - `gh` 命令本身报错（`EOF`／网络）＝**命中并照实登记为"该维未取到"**，不写成"确认为空"。
  - ⚠ 本程一开始 `go test` 之后，第 2 行必然被**本程自己的编译**命中 ⇒ 闸门只在**每批开始之前**跑，批内命中不属争用（这条是实现方 §0.4 先说过的形状，本程独立采用并在此具名）。
- ⚠ **本程未杀过任何进程、未改过 `.github/workflows/**`、未取消过任何 run**；只观察。
- ⚠ **本机墙钟不可信**（派单点名近两天两次量到前跳 8h 以上）⇒ 本程**不用时间戳相减算任何耗时**；所有时刻都是现 `date` 现抄，只作"这一发在闸门之后取的"排序用。

〔独立复现〕§0 每一发都是本程自己跑的三行原样；逐轮日志 `/d/tmp/wisp135ac8r1-gate/gate-rounds.log`、轮询 stdout `/d/tmp/wisp135ac8r1-gate/poll-run*.out` 在盘。

---

## §0.5 ⚠ 本程自己犯的一枚仪器错，以及它把哪些读数作废了（append-only，不抹上一条）

### 0.5.1 错在哪

上面 §0.4 那台 `gate.sh`（v1）只做了一半的事：**它把 `VERDICT` 打在批头，但没有任何东西读那个 verdict**。派单写的是"任一命中 ⇒ 那一批发不作数"，而我的驱动是"先跑闸门、下一行就开炮"——**判 BUSY 也照开**。

第二半是：`gh run list` 从 **12:17 起持续 EOF**（本程 12:33 又连试三发，全 `failed to get runs: Get "https://api.github.com/…/actions/runs?…": EOF`）。v1 把"gh 报错"并进 BUSY 是对的（不许把它读成"确认为空"），但它**和"真有 job 在本机"共用一个词**，于是本程看见 12:14 那批 CLEAR、后面几批的 BUSY 没逐批打开看，就一路把炮开了下去。

```
batch-1.log  gate at=12:12:55  worker=0 toolchain=0 in_progress=0        VERDICT=CLEAR
batch-2.log  gate at=12:14:20  worker=0 toolchain=0 in_progress=0        VERDICT=CLEAR
batch-3.log  gate at=12:17:05  worker=0 toolchain=0 in_progress=0        VERDICT=BUSY   ← 因 gh EOF
batch-4.log  gate at=12:19:49  worker=0 toolchain=0 in_progress=0        VERDICT=BUSY   ← 因 gh EOF
batch-5.log  gate at=12:22:46  worker=0 toolchain=0 in_progress=0        VERDICT=BUSY   ← 因 gh EOF
batch-6.log  gate at=12:28:14  worker=0 toolchain=0 in_progress=0        VERDICT=BUSY   ← 因 gh EOF
```

⇒ **batch-3…batch-6 的全部读数（`MG-*`／`MGF-*`／`U-*`／`1B-*`／`CP-*`／`RST`／`FC`／`NH-*`）按派单口径不作本格的凭据。**
本程**没有**为腾出窗口杀过进程、改过 workflow、取消过 run；这三行只观察。⚠ 顺带一枚要登记的现场事实：闸门输出里 `Runner.Listener.exe` 的 PID 从 **43288 变成 3952**（内存 99MB→82MB）⇒ **本机 runner 服务在中途重启过一次**。本程不知道是谁动的、也不需要知道，但它是"别把远程取不到当成机器空了"的第二条理由。

### 0.5.2 修法（改的是本程自己的台件，不是被审对象）

- `/d/tmp/wisp135ac8r1-gate/gate2.sh`：把两种失败**拆开**——`BUSY-LOCAL`（worker／toolchain 命中＝真有活在机器上）／`BUSY-REMOTE`（远程确有 in_progress）／`LOCAL-CLEAR-REMOTE-UNKNOWN`（本地空、`gh` 取不到 ⇒ exit 2，**明写"这不是确认为空"**）。
- `/d/tmp/wisp135ac8r1-gate/batch.sh`：**轮询到条件成立才点火**（`INTERVAL=45` × `MAXROUNDS=26`；`gh` 连续 `UNK_LIMIT=3` 轮取不到才降级为 `CLEAR-LOCAL-ONLY` 点火，并把那句降级写在日志里；到上界仍不满足 ⇒ `exit 3`、**批不发**）。
- `/d/tmp/wisp135ac8r1-batchR.sh`：**BATCH-R**＝把甲…辛那一整套在点火闸门之下重跑一遍（19 发读数＋2 发邻桶探针），并在批内每 4–6 发复跑一次 gate2。

### 0.5.3 处置

- **本格所有承重读数一律以 BATCH-R 的 `R-*` 为准**（§2 起）。
- 第一遍那些日志与树**留在盘上不删**（`issues/README` 规则 8），档位在本表里标〔仅自述，不背书〕——对本程自己同样适用：它们是本程自己采的，但采的窗口不合格，所以**不背书**。
- ⚠ 形状钉在这里，别只当本程的手抖：**"闸门打印了 verdict 但没人读它"是一类会自我安慰的仪器**——它长得完全像"我在守规矩"，输出里每一批都真的有一行 VERDICT。判据是问一句"**这一行如果被印成 BUSY，点火的那一步会停下来吗**"。下一条写闸门的人请直接读 `batch.sh`，不要读打印器。

〔独立复现〕0.5.1 那六行 `batch-*.log` 与 `gate-rounds.log` 都在盘上可逐轮重看；0.5.2 两枚台件的判读规则写在文件注释里。


---

## §1 锚点、取件、以及"那把尺本程一字未动"

### 1.1 锚点自量（不抄派单），引用到的每一枚 sha 逐枚 `git cat-file -t`

```
$ git rev-parse --short HEAD                                  c8967b8
$ git rev-parse --abbrev-ref HEAD ; @{u}                      dev ; cnb/dev
$ git rev-list --count @{u}..HEAD                             0        ← 被验版本＝已推的 dev tip
$ git cat-file -t c8967b8                                     commit   ← 派单钉的那枚，先核类型
```

本表引用到的每一枚都当场过 `git cat-file -t`（**第 8 代注入的形状就是假 sha 混进工具回显**；本程全部来源是自己 `rev-parse`／`log`／`show`，无一枚抄回显）：

```
c8967b8=commit   c8967b8^=4113cd7=commit   fa35557=commit   fa35557^=d3e3a43=commit
4113cd7=commit   3c2bd48=commit   087ef2e=commit   4f1b73a=commit   faa7cf0=commit
4f344ab=commit   02370fe=commit   4963ae2=commit
```

### 1.2 取件（仓外归档树）

```
$ mkdir -p /d/tmp/wisp135ac8r1-tree && git archive c8967b8 | tar -x -C /d/tmp/wisp135ac8r1-tree
   extract_rc=0
$ find /d/tmp/wisp135ac8r1-tree -type f | wc -l                1079
$ cp <repo>/third_party/sherpa-onnx/{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll \
     /d/tmp/wisp135ac8r1-tree/third_party/sherpa-onnx/          ← 归档里没有（gitignore）
```

⚠ **工作树此刻是干净的**（`git status --porcelain` 空，`count_dirty=0`）——本程**仍**只读归档树，不以此为由省掉取件纪律。

**dll 不是可选项，且本程把它做成"看得见"的**：`scripts/wisp-cli-tests.sh:73,109` 是 CI 自己的形状——它 `export PATH="$root/third_party/sherpa-onnx:$PATH"`，并在缺目录时 `exit 1` 明说"缺 dll 会以 `0xc0000135` 死在**加载之前**、看着像包坏了"。本程的单发驱动 `/d/tmp/wisp135ac8r1-run.sh` 逐字采同形（`export PATH="$tree/third_party/sherpa-onnx:$PATH"`，MSYS 形式而非 `pwd -W` 形式，理由写在那枚脚本的注释里）。

### 1.3 "那把尺一字未动"——**同一种算法三向对**（派单 §5 点名的仪器形状）

被审对象 `135-ac8-mg-impl.md` 引 `sha1sum`＝`23b443ac…`，编排者引 `git hash-object`＝`bde61ddd…`。**两个都对、输入不同**：`git hash-object` 先拼 `blob <长度>\0` 头再摘要。本程因此**只挑一种算法（`sha1sum`，纯文件字节）做三向对**，不去怀疑任何一方造假：

```
$ sha1sum /d/tmp/wisp135ac8r1-tree/cmd/wisp/leg_dispatch_gate_133_test.go
  23b443ac34787aa9ea60e181b9b8c789f468cd56      ← git archive c8967b8 的树
$ sha1sum "<repo>/cmd/wisp/leg_dispatch_gate_133_test.go"
  23b443ac34787aa9ea60e181b9b8c789f468cd56      ← 工作树
$ git show c8967b8:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum
  23b443ac34787aa9ea60e181b9b8c789f468cd56      ← blob 本体
（对照，另一种算法，只列一次不再混用）
$ git hash-object <同一枚树内文件>    bde61ddd2b5a1225ae68b31df196c9a33695e56e
$ git hash-object <旧尺 oldgate.go>  5febd1d0a1db52145a3aa4566ffc48f9639a4dca
```

⇒ **`c8967b8` 里那把尺＝工作树那把尺＝归档树那把尺，逐字节同。**

再看"接续程到底动没动它"——这是本程自己要的一条，不是引用它的自述：

```
$ git diff --numstat fa35557 c8967b8 -- cmd/wisp/leg_dispatch_gate_133_test.go
  （空输出）                       ← 第二程对那把尺确实零改动
$ git diff --numstat fa35557^ c8967b8 -- cmd/wisp/leg_dispatch_gate_133_test.go
  274     28    cmd/wisp/leg_dispatch_gate_133_test.go
$ git diff --name-only fa35557^ c8967b8 | grep -v '^docs/\|^\.scratch/\|^AGENTS.md\|^README.md'
  cmd/wisp/leg_dispatch_gate_133_test.go        ← 唯一一枚非文档改动
```

⇒ **`M-G` 那一发的全部码改动＝这一枚测试文件**（`fa35557`），后面七枚 commit 只带证据文件。

### 1.4 "未修那一侧"的旧尺也核到字节同源（不然 §4 的"洞复现"就不是同一枚洞）

```
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go > /d/tmp/wisp135ac8r1-oldgate.go
$ sha1sum /d/tmp/wisp135ac8r1-oldgate.go                     906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   906201f4a3995d10b0a65910aa4cc68e1f745a9d
```

⇒ 与票 133 第二任验收方在 `133-ac2-r2-acceptance.md` §"AC#2 第二格复算"里记的 `048a9e4:…leg_dispatch_gate_133_test.go | sha1sum = 906201f4…` **逐字同**。
也就是说：**本程当"未修侧"用的那把尺，与造出 `p7`、量到 `100/53/0/0` 全绿的那把尺是同一枚字节** ⇒ §4 的旧尺判绿才配叫"同一枚洞的复现"，而不是另一棵树上的另一回事。

### 1.5 植物"不进本轮编译"这一维——用**第二枚仪器**量，不靠文件名猜

派单 §5 那一条（"任何否定式结论先确认仪器形状"）在这里最容易被糊过去：`git ls-tree` 不加 `-r` 会数出假的 0，"这枚文件不在版本里"这种话不能靠肉眼。本程用 `go list` 直接问构建：

```
（/d/tmp/wisp135ac8r1-SELFTEST，plant=tag、leg=leg2；**这一发只跑 go list，没跑 build/test，不算读数批**）
$ grep -n 'case "probe135ac8":\|^  wisp probe135ac8' cmd/wisp/main.go
  46:  wisp probe135ac8 ticket 135 AC#8 acceptance probe leg, not a product
  108:	case "probe135ac8":
$ grep -c installLogSink cmd/wisp/probe135ac8leg.go        0        ← 第二拍：安装器一个字都不剩
$ go list -f '{{range .TestGoFiles}}…' .  | grep -c probe135ac8     0        ← 本轮测试文件集里查无此件
$ go list -f '{{.IgnoredGoFiles}}' .
  [console_other.go notify_other.go probe135ac8tag_test.go resident_other.go
   secret_dataroot_119b_test.go slo_other.go versioninfo_other.go]           ← 植物在 IgnoredGoFiles 里
```

⇒ "这枚用例今天编不进去"不是本程从文件名后缀**推**出来的，是 `go list` **量**出来的：它在 `TestGoFiles` 里 0 命中、在 `IgnoredGoFiles` 里点名。这枚读数同时是本程 §4 那几发的**独立对照仪器**（它不属于被审的那把尺，也不是 go/parser 走源码）。

〔独立复现〕§1 每一行都是本程自己跑的（仓外只读命令＋`go list`）。
〔日志＋归档，抽验〕"两程各自引的哈希都对"这条的**对方那一侧**数字出自被审对象 §9 与编排者 log，本程只自证了 `sha1sum` 三向。

---

## §1.6 判据 ② 的静态面：新读数**没有**折回"只信 AST"（票面写死的那条不许走的路）

这一节不需要读数，只需要把"它到底新在哪"盘清楚。`fa35557`（＝`M-G` 那一发的唯一一枚码改动）新增的函数**恰好五枚**：

```
$ diff <(grep -oE '^func …' oldgate.go|sort) <(grep -oE '^func …' c8967b8 那枚尺|sort) | grep '^>'
> func (p *pkg133) runRosterDisclosure135
> func (p *pkg133) runRosterReds135
> func compiledRunRoster135
> func flagValue135
> func headRunes135
```

**它有没有开始"自己判断这一枚文件今天编不编"？** 没有——那正是票 133 与本票各裁过一次的"折回只信 AST"：

```
$ grep -nE 'go/build|build\.Context|MatchFile|Constraint|_linux|_windows|_aix|HasSuffix\(.*_test|runtime\.GOOS|goos133' <新尺> | grep -v '^\s*[0-9]*:\s*//'
  438:  isTest := strings.HasSuffix(name, "_test.go")      ← 旧尺就有：只是"这是不是测试文件"，不是"今天编不编"
  1748: if e.IsDir() || !strings.HasSuffix(name,".go") || strings.HasSuffix(name,"_test.go")   ← 旧尺就有
  1912: func goos133() string { return runtime.GOOS }      ← 旧尺就有（oldgate.go 里 grep -c = 1）
  1665 / 1688: 只出现在披露文案的 GOOS=%s 里               ← 用于**说话**，不用于**判**
  1651: 红名句子文本里的 `_linux_test.go`/`//go:build`     ← 用于**给修法举例**，不用于判
```

⇒ 全尺**没有一处**对文件名后缀或 `//go:build` 表达式求值。新读数**唯一**的来源是一枚子进程：

```
1574  exe, err := os.Executable()
1580  cmd := exec.CommandContext(ctx, exe, "-test.list", ".*")
```

⇒ **判据 ② 的静态面成立**：它不是"把 AST 教得更聪明"，是**多了一把独立的眼睛**。

**两条旁证（本程自己找的，不是引用实现方的话）**：

1. 全仓还有谁也写 `covered=` 这种账？`grep -rn 'covered=' --include=*.go --include=*.sh --include=*.ps1 .` 除 `leg_dispatch_gate_133_test.go` 之外**零命中** ⇒ 本格修的是**唯一一处**会印 `covered=test` 的地方，不是一堆平行账本里挑一枚。
2. 那把尺的注释（`:1566-1569`）声称 `scripts/portable-tests.sh` 早就在读同一个对象。本程去核了**那句话指着的真文本**：`scripts/portable-tests.sh:30`（"re-verified against the compiled test binary for THIS platform via `go test -list`… bury it behind a build tag, and the entry goes stale and this step goes red"）、`:346`（`go test -list '.*' …`）、`:381`（逐名 `go test -list "^${name}\$"`）⇒ **注释引的是真行为，不是"预先引用尚未产出的读数"那一类假绿前身**。同一读数形状在仓里已有一枚步级用户，本格的修法与既有实践同形。

〔独立复现〕本节每一行都是本程在归档树／`git show` 上自己量的。

---

## §2 本格到底要证什么（**判据本体，引原文，不引实现方的转述**）

票面 AC#8（`c8967b8` 的 `.scratch/wisp/issues/135-…md:82-93`）要的四条：

| 票面 | 原文要点 | 本程哪一节答它 |
|---|---|---|
| ① | 自己重造 p7 那一形（**不许抄前两程的红名**），给出"整包绿 ＋ 名册里查无此名"的原文 | §4（单跑＋整包两形都在） |
| ② | 修法方向只许"每个 `covered=test` 的名字必须在本轮 `=== RUN` 名册里逐名对上"，**不许折回只信 AST** | §1.6（静态面）＋ §5（行为面） |
| ③ | 修后两拍：门关着零退化（八数逐数不变）＋ p7 必须转红、**红名点到"这名没跑过"而不是"缺标记"** | §5、§6 |
| ④ | 全程 `-count=1 -v` 取四数并列名册差集，`SKIP=0`、`panic=0`、`build-failed=0` 逐发给出 | §3 起每一张表的 `EIGHT` 行 |

来源判据（票 133 第二任验收方 `133-ac2-r2-acceptance.md` §3.4 ①，本程逐字引）：

> **`M-G`**：一枚**本平台不编译**的 `func TestXxx(t *testing.T)`（种法＝文件名 `_linux_test.go` 后缀，**或显式 `//go:build` 表达式为假**），坐在**被验那一腿自己的目录**里、身体里调用那条腿的 entry。
> 判据：本尺必须红、红名点到本尺，且**第二拍也必须红**（把 install 那一截删掉之后仍须红）；反向控制一枚：同形写在**本会编译**的文件里 ⇒ 必须绿（不许把合法形状打死）。

⇒ 本程因此**两枚种法都造**（`//go:build` 表达式为假 ＋ 文件名后缀），后缀那一支本程**故意不跟前两程同名**（他们用 `_linux_test.go`，本程用 `_aix_test.go`）。
⚠ 本程**没有**跑票 133 那套台件（`/d/tmp/wisp133-r2-run.sh`／`-p7.sh` 在盘上还在，本程只 `ls` 到、未执行）——那是别人的可重跑凭据，且 **C1 是票 133 自己那一格的账，不由本程背书**（本程只在 §13 说清"本程哪一发在实质上是 C1 那句话的等价物"）。

### 2.1 本程的种法（名字一枚都没抄）

| 维度 | 票 133 验收方 | 票 135 第一程 | 票 135 第二程 | **本程（裁决方）** |
|---|---|---|---|---|
| 腿 | `sfx131` | `sfx135g` | `sfx135r2` | **`probe135ac8`** |
| entry | `cmdSfx131` | — | `cmdSfx135r2` | **`cmdProbe135Ac8`** |
| 落地文件 | — | — | `sfx135r2leg.go` | **`probe135ac8leg.go`** |
| 植物文件 | `probe133r2b_linux_test.go` | `probe135mg_linux_test.go` | `probe135r2m_linux_test.go` | **`probe135ac8tag_test.go`**（`//go:build probe135ac8r1_never_satisfied`）／**`probe135ac8_aix_test.go`**（隐式 GOOS=aix） |
| 用例名 | `TestR2P7LinuxOnlyCaseDrivesTheLeg` | `TestMGGp1…` | `TestR2MgLinuxTaggedCaseDrivesThePlantedLeg` | **`TestAc8r1UncompiledTagCaseDrivesThePlantedLeg`**／**`TestAc8r1AixSuffixCaseDrivesThePlantedLeg`** |
| 反向控制 | — | — | `TestR2MgCompiledCaseDrivesThePlantedLeg` | **`TestAc8r1CompiledCaseDrivesThePlantedLeg`** |

台件：`/d/tmp/wisp135ac8r1-mutate.py`（种法）、`/d/tmp/wisp135ac8r1-run.sh`（单发驱动：**先证落地再读数**——`grep -n` 落地行＋`installLogSink` 命中数＋`go list` 的 `TestGoFiles` 命中数，再 `go build ./cmd/wisp/` 与 `go build ./...` 双 rc，两 rc 非 0 就 `exit 7` 且**一个读数都不产生**）。

### 2.2 ⚠ 本程第二枚仪器坑（自己踩的，形状要钉住）：`grep -c "=== RUN   $name"` 会把**指控本身**数成一次运行

被审那把尺的红句里**逐字含**一串 `no "=== RUN   TestAc8r1…"`。本程第一版驱动用不带行首锚的 `grep -c "=== RUN   $pn"` 去数那枚植物跑没跑，于是把红句自己数成 `===RUN-hits=1`，看着像"植物其实跑了"。
⇒ 驱动已改成 `grep -c "^=== RUN   $pn$"`（行首锚），并从**留在盘上的原始日志**重算，未重发也不需重发。核对：`R-MGN2.log` 里那枚名字 log-hits=**2**（账本行 1 ＋ 红句 1）、**anchored `=== RUN` 行 = 0**、名册里也 0（`in-roster=0`）。
形状：**"某串文本出现次数"这种数法，凡是那一串可能被别人（尤其是判你红的凶手）复述，就必须锚死行首。** 本程第一遍的 `MG-*` 那批里这一栏是脏的——那批本来就因 §0.5 的闸门问题不作凭据，这里再记第二笔账。

---

## §3（甲）基线，本程自采——**不引用实现方那 `101/54/0/0`**

批：BATCH-R，点火闸门 `12:37:57 GATE=CLEAR-LOCAL-ONLY`（本地两维连 3 轮 0；`gh` 取不到已按 §0.5.2 明写成"不是确认为空"）。

| 发 | 尺 | 命令 | rc | RUN | 顶层 P/F/S | 子测 P/F/S | SKIP 行 | build-failed | panic | 名册枚数 |
|---|---|---|---|---|---|---|---|---|---|---|
| `R-BLNEW` | 新 | `-count=1 -v -run '^TestAC1AC2DispatchHopGate133$'` | **0** | 1 | 1/0/0 | 0/0/0 | 0 | 0 | 0 | 1 |
| `R-BLOLD` | 旧 | 同上 | **0** | 1 | 1/0/0 | 0/0/0 | 0 | 0 | 0 | 1 |
| `R-BLFULL` | 新 | `-count=1 -v`（整包，门开着） | **0** | **101** | **54/0/0** | **47/0/0** | 0 | 0 | 0 | 54 |
| `R-BLCLOD` | 新 | `-count=1 -v -skip '^TestAC4EveryLegIsNailedOrRuled$'`（门关着） | **0** | 100 | 53/0/0 | 47/0/0 | 0 | 0 | 0 | 53 |

⇒ **本程自己复现了 `101/54/0/0 ＋ 子测 47/0/0`、红名 0 枚、SKIP 0、panic 0、build-failed 0**（与被审对象 §3 同数，但**这一行不是引它，是本程自己那四发的读数**）。

**逐名名册**（`R-BLFULL`，54 枚，`/d/tmp/wisp135ac8r1-out/R-BLFULL.runnames`，全量 54 行在盘不重列）。附带一枚比对——把被审对象 §3 列的那 54 枚抽出来与本程这份**逐名对**：

```
$ sed -n '218,244p' docs/evidence/s1/135-ac8-mg-impl.md | tr -s ' \t' '\n' | grep -E '^Test[A-Za-z0-9_]+$' | sort -u > impl54.txt
$ comm -23 R-BLFULL.runnames impl54.txt   ->  （空）
$ comm -13 R-BLFULL.runnames impl54.txt   ->  （空）
```

⇒ 两份名册逐枚同名。⚠ 这一行**只**说明"本程自采的基线与他们抄的那份是同一棵树同一套名字"，**不是**"他们的读数被本程复现"——本程复现的是上面那四发自己的数。

**名册差集（基线两形）**：
```
comm -23 R-BLFULL.runnames R-BLCLOD.runnames  ->  TestAC4EveryLegIsNailedOrRuled
comm -13 R-BLFULL.runnames R-BLCLOD.runnames  ->  （空）
```
⇒ 门关着那一形是**少跑一枚**（RUN 101→100、`--- SKIP` 行仍 **0**），**不是**跑成 SKIP——票 133 验收方 §1.10 那条"判不再红先分清变绿还是被跳过"的口径，本程独立采到同一形状。

**新尺在未变异基线上不响**（恒真自查的第一半）：`R-BLNEW` 的披露行原文——
```
run-roster disclosure: GOOS=windows, 54 startable cases read from this binary itself
(`-test.list '.*'`); case names this round's ledger credited through `covered=test`:
1 distinct, 1 of them startable …
```
⇒ 今天账本里那 1 枚 `covered=test` 名字**在**名册里 ⇒ 新判据今天不响。

〔独立复现〕四发全部本程自采（树 `T-R-BLNEW`／`T-R-BLOLD`／`T-R-BLFULL`／`T-R-BLCLOD`；日志同名 `.log`＋`.verdicts`＋`.runnames`）。

---

## §4（乙）**旧尺判绿＝病真的在**：整包绿 ＋ 名册里查无此名，原文

旧尺＝`fa35557^`（＝`d3e3a43`）那枚，`sha1sum 906201f4…`（§1.4 已核与票 133 验收方同一枚字节）。

### 4.1 落地证明（`R-MGFO2C`，`/d/tmp/wisp135ac8r1-out/R-MGFO2C.verdicts` 原文摘录）

```
GATE-FILE-SWAPPED-TO /d/tmp/wisp135ac8r1-oldgate.go
GATE-MODE old gate-sha1sum=906201f4a3995d10b0a65910aa4cc68e1f745a9d …
LANDED-LEG leg2
LANDED-PLANT tag name=TestAc8r1UncompiledTagCaseDrivesThePlantedLeg
LANDED-MAIN:
  46:  wisp probe135ac8 ticket 135 AC#8 acceptance probe leg, not a product
  108:	case "probe135ac8":
  probe135ac8leg.go lines=9 installLogSink-hits=0        ← 第二拍：安装器一个字节都不剩
  9:func TestAc8r1UncompiledTagCaseDrivesThePlantedLeg(t *testing.T) {
GOLIST-TESTGOFILES cmd/wisp count=15 probe135ac8-hits=0  ← 独立仪器量到：本轮构建不收这枚文件
BUILD-cmdwisp R-MGFO2C rc=0
BUILD-all   R-MGFO2C rc=0
READ label=R-MGFO2C gate=old leg=leg2 plant=tag args='-skip ^TestAC4EveryLegIsNailedOrRuled$' rc=0
```

### 4.2 读数

| 发 | 尺 | 形 | rc | 八数 | 红名 | 那枚植物的名字 | 账本里那一行（逐字） |
|---|---|---|---|---|---|---|---|
| `R-MGO2` | 旧 | 单跑本尺，第二拍＋`//go:build` 假 | **0** | 1/1/0/0＋0/0/0 | **无人** | log 1 次／`=== RUN` **0 次** | `covered=test TestAc8r1UncompiledTagCaseDrivesThePlantedLeg drives cmdProbe135Ac8` |
| `R-MGO2B` | 旧 | 单跑本尺，第二拍＋`_aix` 后缀 | **0** | 1/1/0/0＋0/0/0 | **无人** | log 1 次／`=== RUN` **0 次** | `covered=test TestAc8r1AixSuffixCaseDrivesThePlantedLeg drives cmdProbe135Ac8` |
| `R-MGFO2C` | 旧 | **整包、门关着**，第二拍＋`//go:build` 假 | **0** | **100/53/0/0**＋47/0/0 | **无人** | log **1 次**／`=== RUN` **0 次**／名册 **0 次** | 同 `R-MGO2` 那一行 |

`R-MGFO2C` 那"1 次"是本尺自己打的账本行（`grep -n` 在 `R-MGFO2C.log:10`）。SKIP 行 0、panic 0、build-failed 0。

⇒ **票面 ① 成立，且是本程第一手**：一枚名字对、签名对、身体里真断言、`go list` 亲口说"本轮不收这枚文件"的用例，在旧尺上拿到 `covered=test … drives …` 背书，**整包 100/53/0/0 全绿、门关上也是全绿、它在 `=== RUN` 里 0 次**。
⇒ 病不是文字问题：**两枚不同的藏法（`//go:build` 表达式为假、文件名隐式 GOOS）都判绿**，本程各造一发。

〔独立复现〕`T-R-MGO2`／`T-R-MGO2B`／`T-R-MGFO2C` 三棵树、三份日志在盘。

---

## §5（丙）**新尺判红＝牙在**：同一形、同一棵树形状，只换尺

| 发 | 尺 | 形 | rc | 八数 | 红名 | 那枚植物 |
|---|---|---|---|---|---|---|
| `R-MGN2` | 新 | 单跑本尺，第二拍＋`//go:build` 假 | **1** | 1/0/1/0＋0/0/0 | `TestAC1AC2DispatchHopGate133` | log 2 次／`=== RUN` **0 次**／名册 0 次 |
| `R-MGN2B` | 新 | 单跑本尺，第二拍＋`_aix` 后缀 | **1** | 1/0/1/0＋0/0/0 | `TestAC1AC2DispatchHopGate133` | log 2 次／`=== RUN` **0 次**／名册 0 次 |

**账本行与旧尺逐字同**（`covered=test … drives cmdProbe135Ac8`，§4 与 §5 那两行一模一样）⇒ **变的只有判，没有读数口径**。
"同一枚尺上只差那一枚植物"的四格对照，本程自己采齐：

```
旧尺:  R-U2O 第二拍＋usage（无植物） rc=1 covered=RED nothing
       R-MGO2 第二拍＋usage＋植物    rc=0 covered=test <植物名>     ← 钉一枚不存在的名字就能过闸
新尺:  R-U2N 第二拍＋usage（无植物） rc=1
       R-MGN2 第二拍＋usage＋植物    rc=1                          ← 牙在
```

### 5.1 红句口径核对（票面 ③ 要的那一句，本程核实现方说法对不对）

`R-MGN2.log` 第 3 行原文（本程 `grep` 出，不改一字）：

```
leg_dispatch_gate_133_test.go:229: AC#1/#2 RED: leg "probe135ac8" (main.go:108) is booked in
the ledger below as covered by the case "TestAc8r1UncompiledTagCaseDrivesThePlantedLeg", and
this round's test binary has no such case: the roster read from the running binary lists 54
startable cases and "TestAc8r1UncompiledTagCaseDrivesThePlantedLeg" is not one of them, so no
"=== RUN   TestAc8r1UncompiledTagCaseDrivesThePlantedLeg" line exists in this run or can exist
in it. The name is declared at probe135ac8tag_test.go:9, which is the point - the declaration
is in the sources and the sources are not the build.
```

对照**另一种说法**（"缺标记"那一支，`R-U2N.log` 的红句，同一把新尺、同一拍、只差没有植物）：

```
AC#1 RED: leg "probe135ac8" (main.go:108) is dispatched by func main and covered by nothing:
no nail in this gate's registry, no test case in this directory that drives a symbol belonging
to this leg alone, and no WISP-LEG-COVERAGE-RULING sentence naming it.
```

⇒ **两句话在盘上就是两枚不同的分支、两种不同的字**：一支说"这轮没有、也不可能有它这一行 `=== RUN`"（日志前缀 `leg_dispatch_gate_133_test.go:229`＝`runRosterReds135` 那圈 `t.Errorf`），一支说"缺覆盖/缺标记"（日志前缀 `leg_dispatch_gate_133_test.go:209`＝`coverageReds133` 那圈）。两枚行号本程从 `.log` 里现 `grep` 出来，不是从源码推的。
⇒ **实现方 §9 那句"红名点到的是前者"——本程独立核＝对。**

〔独立复现〕`T-R-MGN2`／`T-R-MGN2B`／`T-R-U2N`／`T-R-U2O`；红句可从那四份 `.log` 里逐字重 `grep`。


