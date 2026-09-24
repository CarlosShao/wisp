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

### 5.2 红的**归属**（防止"其实是别的分支响的"这一类附条件）

一枚红必须有唯一可指的出处。本程把每一发的 `t.Errorf` 前缀（`file:line`）逐枚 `grep` 出来数：

| 发 | 该发里出现的红分支（`file:line` × 枚数） | 该分支是谁 |
|---|---|---|
| `R-MGN2` | `leg_dispatch_gate_133_test.go:229` × **1** | `runRosterReds135`（＝check (c)，**新那一支**） |
| `R-MGN2B` | 同上 × 1 | 同上 |
| `R-MGFN2C`（整包·门关着） | 同上 × 1 | 同上 |
| `R-MGFO2C`（旧尺同形） | **零枚** | —（旧尺没有 `:229` 这一支） |
| `R-U2N`／`R-1BN` | `…:209` × 1 | `coverageReds133`（**旧那一支**，与改动无关） |
| `R-CPN`／`R-MGF131` | 零枚 | — |

⇒ **M-G 那几发的红**只由新那一支判出、每发恰好一枚、没有旁支混进来；**对照那两发**的红只由旧支判出。⇒ "红名点到的是这名没跑过"这句话是**有出处的**，不是裁决方替它挑的。

---

## §6（丁）⚠ 票面写死的那条反向判据：**拆掉 install 的第二拍也必须红**——两粒度都量，并比名册差集

这一发是票 133 AC#2 第二次退回时加的形状（`133-ac2-r2-acceptance.md` §3.4 ①"且**第二拍也必须红**"），**不做＝本格不成立**。本程两粒度都做了。

### 6.1 第二拍"install 真被拆掉"的落地证明（不是本程嘴上说拆了）

```
R-MGFN2C 落地行：  probe135ac8leg.go lines=9 installLogSink-hits=0
                   108:	case "probe135ac8":            ← 腿还在分发里
                   GOLIST-TESTGOFILES cmd/wisp count=15 probe135ac8-hits=0
```
⇒ 那条腿**还在被 `func main` 分发**、身体里**一个安装器字节都没有**、植物**不进本轮构建**。三样都是盘上行，不是推。

### 6.2 两粒度读数（新尺）

| 粒度 | 发 | 命令 | rc | 八数 | 红名 | 名册枚数 | SKIP/panic/build-failed |
|---|---|---|---|---|---|---|---|
| 单跑本尺 | `R-MGN2` | `-run '^TestAC1AC2DispatchHopGate133$'` | **1** | 1/0/1/0＋0/0/0 | `TestAC1AC2DispatchHopGate133` | 1 | 0/0/0 |
| 整包·门关着 | `R-MGFN2C` | `-skip '^TestAC4EveryLegIsNailedOrRuled$'` | **1** | **100/52/1/0**＋47/0/0 | `TestAC1AC2DispatchHopGate133` | 53 | 0/0/0 |
| 整包·两门都活 | `R-MGFN2O` | 无筛 | **1** | 101/53/1/0＋47/0/0 | `TestAC1AC2DispatchHopGate133` | 54 | 0/0/0 |

⇒ **第二拍在整包形上也红**（旧尺同一形是 `R-MGFO2C` rc=0 100/53/0/0）⇒ 没有被"另一种写法已经红过"抵掉。

### 6.3 第一拍**独立**红（另一拍，不互相抵账）

| 发 | 尺 | 形 | rc | 账本那一行（逐字） |
|---|---|---|---|---|
| `R-1BN` | 新 | 第一拍（装了听众）＋植物，单跑 | **1** | `installs=true … covered=RED sink with no nail` |
| `R-1BO` | 旧 | 同上 | **1** | 同上（两枚尺逐字同） |

⇒ 第一拍那枚红**两把尺都响**（它响在 `installs` 那一支，比这次改动老）⇒ **本程按票面口径：它不能当 (c) 的修复凭据**，只当"改完没把第一拍弄哑"的证据。它落在 `installs` 那一支、`R-MGN2` 落在 `runRosterReds135`（`:229` vs `:209`），**两拍各自独立、不互相抵账**——这一条本程是用两发独立读数量的，不是引用谁的断言。

### 6.4 名册差集（票 133 §1.10 那口径：既看"新红的"，也看"被吞掉的"）

```
R-MGFN2C vs R-BLCLOD(53/53)  only-in-shot=[]  only-in-base=[]
R-MGFN2O vs R-BLFULL(54/54)  only-in-shot=[]  only-in-base=[]
R-MGFO2C vs R-BLCLOD(53/53)  only-in-shot=[]  only-in-base=[]
```

⇒ **零缩小、零新增**；每一发 `--- SKIP` 行 0、`panic` 0、`build-failed` 0（§3–§6 全部 19 发的 `EIGHT` 行都带这三栏，无一非零）。
⇒ 逐名换色只有一枚：`TestAC1AC2DispatchHopGate133` `--- PASS`→`--- FAIL`（`R-BLCLOD` 的 `REDNAMES:` 空、`R-MGFN2C` 的那一栏只有它）。植物那枚用例**没进名册**（正是被验的那件事），也**没吞掉**任何一条别的读数。
⇒ **票面 ③ 的"门关着零退化"另有硬证**：把新尺装在**未变异**的树上，门关着仍是 `R-BLCLOD` 100/53/0/0、门开着仍是 `R-BLFULL` 101/54/0/0（§3），与旧尺基线（`R-BLOLD` rc=0）一致 ⇒ 这把尺今天**没有改变任何一格既有读数**（BATCH-S 的 `R-BLOLDF`／`R-BLOLDC` 再补一整包粒度的旧尺基线，见 §7.2）。

〔独立复现〕`T-R-MGFN2C`／`T-R-MGFN2O`／`T-R-MGFO2C`／`T-R-1BN`／`T-R-1BO`＋`analyseR.txt` 第 3 节。

---

## §7（戊）反向控制：同一形写进**这轮会编译**的文件 ⇒ 必须判绿

| 发 | 尺 | 植物 | `go list` 的 `TestGoFiles` | rc | 八数 | 披露行名册枚数 | 账本那一行 |
|---|---|---|---|---|---|---|---|
| `R-CPN` | 新 | `probe135ac8ok_test.go`（无平台约束） | count=**16**，probe135ac8-hits=**1** | **0** | 1/1/0/0＋0/0/0 | **55** startable | `covered=test TestAc8r1CompiledCaseDrivesThePlantedLeg drives cmdProbe135Ac8` |
| `R-CPO` | 旧 | 同一发 | count=16，hits=1 | **0** | 1/1/0/0 | —（旧尺无名册读数） | 同一行 |

⇒ **新尺不是"一味喊红"**：植物编得进去 ⇒ 名册 54→**55**、判绿；藏得起来 ⇒ 名册 54、判红。两发其余字节同（同一棵基线树、同一条腿、同一个 entry、同一行 usage、同一拍）。
⇒ 且**旧尺本来也放行这一发**（`R-CPO` rc=0）⇒ 这枚绿不是新尺偏门造出来的。

**一枚本程自己踩出来的附带控制（形状要留，别当漂亮）**：`MG-N2A`／`MG-O2A`（第一遍批内）本意是"文件名后缀那一支"，但本程把文件写成 `probe135ac8aix_test.go`——`_aix` 前面没有下划线，**不匹配 `<name>_<GOOS>_test.go>` 形状**，于是那枚文件**照常进构建**（驱动里那行 `GOLIST … probe135ac8-hits=1` 当场把它量了出来，而新尺对它判绿＝行为正确）。⇒ 那两发**不是** M-G 一发，而是**又变成一枚反向控制**；本程修成 `probe135ac8_aix_test.go` 后重发＝`R-MGO2B`／`R-MGN2B`（§4/§5 用的是这两发）。这枚坑值得钉：**"文件名后缀自带平台约束"这件事不能靠肉眼看文件名，必须 `go list` 量**——否则裁决方自己会造出一发假红/假绿。

### 7.2 BATCH-S：把 (戊) 与"装新尺零退化"提到**整包**粒度（点火闸门 `12:54:04 GATE=CLEAR-LOCAL-ONLY`）

单跑那一形只能证"startable"，证不了"started"——那要把整包跑一遍（披露行自己就是这么划的）。本程补了六发：

| 发 | 尺 | 形 | rc | 八数 | 红名 | 名册枚数 | 那枚植物 |
|---|---|---|---|---|---|---|---|
| `R-BLOLDF` | 旧 | 未变异·整包门开着 | **0** | 101/54/0/0＋47/0/0 | 零枚 | 54 | — |
| `R-BLFULL` | 新 | 同一形（§3） | **0** | 101/54/0/0＋47/0/0 | 零枚 | 54 | — |
| `R-BLOLDC` | 旧 | 未变异·整包门关着 | 0 | 100/53/0/0＋47/0/0 | 零枚 | 53 | — |
| `R-CPNF` | 新 | **编译得进去的植物**·整包 | **0** | **102/55/0/0**＋47/0/0 | 零枚 | **55** | `=== RUN`×**1**、`--- PASS`×1 ⇒ **真跑了** |
| `R-CPFO` | 旧 | 同一形 | 0 | 102/55/0/0＋47/0/0 | 零枚 | 55 | 同上 |
| `R-MGOFC` | 旧 | `_aix` 后缀植物·整包门关着 | **0** | 100/53/0/0＋47/0/0 | 零枚 | 53 | `=== RUN` **0 次** |
| `R-MGNFC` | 新 | 同一形 | **1** | 100/52/1/0＋47/0/0 | `TestAC1AC2DispatchHopGate133` | 53 | `=== RUN` **0 次** |

名册差集（三条，都是双向算）：

```
R-BLOLDF  vs R-BLFULL  ->  双向皆空          ⇒ 装上新尺，今天这棵树的整包名册一枚不多一枚不少
R-CPNF    vs R-BLFULL  ->  only-in-shot=[TestAc8r1CompiledCaseDrivesThePlantedLeg]  only-in-base=[]
                                              ⇒ (戊) 那句"名册从 54 变 55"在整包形上成立，且多出的正是本程那枚
R-MGNFC   vs R-BLCLOD  ->  双向皆空          ⇒ 藏起来那一枚既没进名册、也没吞掉任何一条（门关着 53 枚逐名同）
R-MGOFC   vs R-BLCLOD  ->  双向皆空          ⇒ 旧尺那一发同上
```

⇒ **三件事各自独立钉住**：① 新尺**不改今天任何读数**（两粒度的旧/新基线名册逐名同、八数逐数同）；② 合法形状**判绿且真跑**（`R-CPNF` 那枚植物在整包里 `=== RUN`／`--- PASS` 各 1 次，不是"被跳过"）；③ `_aix` 那一支（实现方 §8.1 自己列为**未测支**）在旧尺上同样放行、在新尺上同样转红、且整包门关着也红。

〔独立复现〕`T-R-CPN`／`T-R-CPO`／`T-R-MGO2B`／`T-R-MGN2B`；`T-MG-N2A`／`T-MG-O2A` 两棵作废树留盘不删（§0.5.3 同一处置）。

---

## §8（己）恒真自查：**本程独立重走**，并正面回答"有没有哪一发在未修旧码上就已经响"

票面这一问的形状是：**"修完后复跑 X 必须红"里的 X 若在未修旧码上已经全响，那一发就不是修复凭据。** 本程每一发动手前都问了一句，答案按发列在盘上（旧尺那一列全部是本程自己发的 `R-*O*` 发，不是抄实现方 §4.6）。

| 这一发（本程自造） | 在**未修**那侧（旧尺 `906201f4…`）响不响 | 在**修了**这侧（新尺 `23b443ac…`）响不响 | 能不能当"修好了"的凭据 |
|---|---|---|---|
| 第二拍＋`//go:build` 假植物（`R-MGO2` → `R-MGN2`） | **不响**（rc=0，账本给它 `covered=test`） | **响**（rc=1，红名逐点到那枚名字，`:229`） | **能——本格的承重墙** |
| 第二拍＋`_aix` 文件名后缀植物（`R-MGO2B` → `R-MGN2B`） | **不响**（rc=0） | **响**（rc=1，同一支 `:229`） | **能**（同族第二枚藏法，且是**实现方自己 §8.1 第 3 条列为未测支**的那一支） |
| 整包·门关着（`R-MGFO2C` → `R-MGFN2C`） | **不响**（rc=0 100/53/0/0） | **响**（rc=1 100/52/1/0） | **能**（反向判据那一发的整包形） |
| 第二拍＋只加 usage、无植物（`R-U2O` → `R-U2N`） | **响**（rc=1，`covered=RED nothing`，`:209`） | 响 | **不能**——两侧都响。它只证"这条腿今天确实没人覆盖"，是对照 |
| 第一拍（`R-1BO` → `R-1BN`） | **响**（rc=1，`installs=true covered=RED sink with no nail`） | 响 | **不能**当 (c) 的凭据；只算"改完没把第一拍弄哑" |
| 编译得进去的植物（`R-CPO` → `R-CPN`） | 不响 | **不响（必须不响）** | 不是修复凭据，是"不许打死合法形状"的凭据 |
| 基线（`R-BLOLD`／`R-BLNEW`／`R-BLFULL`／`R-BLCLOD`） | 不响 | **不响** | 恒真自查第一半：新判据今天不响 ⇒ 不是装饰 |
| `rerr != nil` 那一支（`R-FC-CHILDENV`，只注入环境变量、树未变异） | （旧尺无此支，无从响） | **响**（rc=1，`could not be read … Red, not skipped`） | 它证的是**这条判据不靠"看不见就当没事"**；不是 M-G 本形的凭据 |

**正面回答那句问：** **有**——`R-U2N`／`R-U2O`（无植物对照）与 `R-1BN`／`R-1BO`（第一拍）**两发在未修旧码上就已经响**。本程已把它们**明确排除在修复凭据之外**（上表第 4、5 行）。
⇒ **本格不只剩这种发**：承重的是"旧侧不响、新侧响"那三行（`R-MGO2→R-MGN2`、`R-MGO2B→R-MGN2B`、`R-MGFO2C→R-MGFN2C`），本程各各第一手。

〔独立复现〕八行全部对上盘上的 `R-*.verdicts`；旧侧那一列是本程自己发的，未引用实现方 §4.6 任何一格。

---

## §9（庚）还原：变异面逐枚列、还原发与快照源码逐枚同、仓库残留 0 枚

本程的还原是**结构性**的——每一发一棵新 `cp` 的树，工作树从头到尾没种过件。仍按票面把三样量了。

### 9.1 变异面逐枚（`diff -rq 归档树 ↔ 那一发的树`，全文在 `analyseR.txt` 第 8 节）

```
R-BLNEW / R-BLFULL / R-BLCLOD / R-RST      （无差异，只剩 wisp.exe 构建副产物）
R-BLOLD                                    只差 cmd/wisp/leg_dispatch_gate_133_test.go（旧尺换上来）
R-U2N / R-1BN                              main.go differ | + probe135ac8leg.go          (+ tag 植物，仅 1BN)
R-MGN2 / R-MGN2B                           main.go differ | + probe135ac8leg.go | + 植物一枚
R-MGO2 / R-MGO2B / R-MGFO2C / R-1BO / R-U2O / R-CPO   同上，再多一行“尺文件 differ”（旧尺那侧）
R-CPN                                      main.go differ | + probe135ac8leg.go | + probe135ac8ok_test.go
R-MGFN2C / R-MGFN2O / R-MGF131             main.go differ | + probe135ac8leg.go | + probe135ac8tag_test.go
R-NH1                                      只差 cmd/wisp/leg_sink_nail_131_windows_test.go（邻桶探针，§11）
```
⇒ **没有任何一发把改动漏在快照之外**：每棵树的差异面恰好＝本程种的那几枚文件（外加旧尺那侧多一行尺文件差异），零意外路径、零仓内路径。

### 9.2 还原发复绿 ＋ 与基线那发逐字比

```
$ R-RST: 新尺、无种件、-count=1 -v -run 本尺  ->  rc=0，1/1/0/0＋0/0/0，SKIP 0、panic 0、build-failed 0
$ diff -rq /d/tmp/wisp135ac8r1-tree /d/tmp/wisp135ac8r1-T-R-RST      ->  （除 wisp.exe 外零差异）
$ sed -E 's/time=[^ ]+//; s/\([0-9.]+s\)//; s/[0-9.]+s$//' R-BLNEW.log > bl.norm ; 同一法做 R-RST.log
$ diff -u bl.norm rst.norm      ->  **空输出，diff_rc=0**
```
⇒ 还原发与基线发**只差时间戳与构建耗时两栏**，其余逐字相同（11→11 legs、4 claims、名册枚数、披露句全同）。

### 9.3 仓库侧三查（收尾时再量一次，同一组命令）

```
$ git status --porcelain                          （只有本文件那一枚路径，且已带 pathspec 提交）
$ sha1sum cmd/wisp/leg_dispatch_gate_133_test.go  23b443ac34787aa9ea60e181b9b8c789f468cd56  ← 与 §1.3 取件时同
$ grep -rIl 'probe135ac8\|TestAc8r1' . --exclude-dir=.git
  ./docs/evidence/s1/135-ac8-r1-acceptance.md     ← 唯一命中＝本程证据文件自己；生产码/测试码残留 **0 枚**
```

〔独立复现〕§9 三样本程自量，命令原样在上面。

---

## §10（辛）实现方自报的那条，本程独立重走：**131 的那道门接不接这一形**

被审对象 §4.5 报："摘掉本尺整包 rc=0 ⇒ 131 的门不接这一形（还读到它打了这行腿又 PASS）"。这一条决定**本格的必要性与归谁**——若 131 的门其实接得住，这一格就是重复装钉。本程自己发一发：

```
R-MGF131：新尺（**装好但整包时被 -skip 掉**）、第二拍＋M-G 植物、-count=1 -v -skip '^TestAC1AC2DispatchHopGate133$'
  -> rc=0，100/53/0/0 ＋ 子 47/0/0，SKIP 行 0、panic 0、build-failed 0
  -> 名册 53 枚 vs R-BLCLOD 53 枚：only-in-shot=[TestAC4EveryLegIsNailedOrRuled] only-in-base=[TestAC1AC2DispatchHopGate133]
     （只是"哪一枚门被 -skip 掉"的对调，零缩小）
  -> 那枚植物名字在**整份日志里出现 0 次**
  -> 131 的门自己打了这一行腿、然后 PASS：
     leg probe135ac8  main.go:108   install=false records=false ruled=false nails=-   -> no records   (:91)
     --- PASS: TestAC4EveryLegIsNailedOrRuled (0.02s)                                  (:97)
```

⇒ **判：131 的门真的不接这一形**——不是"它没看见"（它看见了这条腿、把它归成 `no records`），也不是"它被跳过"（`--- SKIP` 行 0，它是 `--- PASS` 真跑真绿）。
⇒ 两枚门的分工是**盘上事实**：131 的门问"这条腿到没到听众、有没有记录"；本形里那条腿**既不到听众也没记录**（第二拍就是把 install 拆掉的那一拍），于是在 131 那里落进"无记录"这一合法桶。能把这条腿钉住的只剩"账本里那枚 `covered=test` 的名字今天编不编"这一问，而那一问**只有新尺答**。
⇒ **本格的必要性成立、不是重复装钉。**

〔独立复现〕`T-R-MGF131`＋`analyseR.txt` 第 3、5 节。

---

## §11 本程自造的一枚**邻桶**探针（不是本格判据要的，但它关乎"这格还剩多少"）

实现方 §2 登记了那把尺自己写明的一处不对称：`covered=test` 落名册外 ⇒ **红**；`covered=nail`（这枚尺自己那 4 条登记表 `legCovers133`）落名册外 ⇒ **只披露不弄红**。本程不去辩护也不去复述，直接**造一发量它**：给那枚住过 windows 的钉文件 `leg_sink_nail_131_windows_test.go` 头上插一行本机器永不满足的 `//go:build`（其余一字不动）。

```
R-NH1（单跑本尺，未变异树＋只那一行改动，go list 量到它不再进 TestGoFiles）
  -> rc=0（**本尺判绿**），1/1/0/0，红名零枚
  -> 披露行：run-roster disclosure: GOOS=windows, 51 startable cases …; this gate's registry:
     4 claims, **2 of them outside this round's roster - TestAC2ModelsLegBooksItsHandOffVerdictOnDisk,
     TestAC3SecretLegBooksItsAuditRecordsOnDisk. A registry nail the build does not take is
     disclosed here and not reddened …
  -> 而账本那两行照印：leg models … covered=nail TestAC2ModelsLeg… -> cmdModels ／ leg secret … covered=nail TestAC3SecretLeg…

R-NH2（同一棵树、整包两门都活）
  -> **rc=1**，51 枚顶层跑、50 PASS/1 FAIL、子 45/0/0、SKIP 0、panic 0
  -> 唯一红名＝ TestAC4EveryLegIsNailedOrRuled（**131 的门**），红句原文：
     AC#4 RED (the instrument, not the code): zero nails registered in this test binary, so the
     gate has nothing to reconcile and would pass on an empty list. GOOS=windows compiled no nail
     file … This reading stays red on purpose: a gate that cannot see is not a gate that has seen
     nothing to complain about.
```

⇒ **三句结论，分开说**：
1. 邻桶那一形**本尺确实放行**（rc=0）⇒ 实现方 §2 那句不对称**本程独立重走＝成立**，不是文字托词。
2. 但**整包在这一形上不会静默全绿**——131 的门**接**它（它的钉是 `registerLegNail131` **运行时注册**，藏起文件就＝注册数为 0 ⇒ fail-closed 红）。⇒ 这一形**不是**双叠假绿，本格不因它扣分。
3. 但两枚门的分工在这里正好**反**了一回（§10：`covered=test` 那形只有 133 接；§11：`covered=nail` 这形只有 131 接）。⇒ **建议登记成家族新的一发（下枚编号 `M-H`），形状＝"两枚门各接一桶，但没有一桶是两枚门都接"**；本程只量不修，也不把它当本格的破口。

〔独立复现〕`T-R-NH1` 一棵树两份日志（`R-NH1.log` 单跑／`R-NH2.log` 整包）；落地行＝`NH-GOLIST-in-build=0`、`NH-BUILD-all rc=0`（`batchR.log`），变异面只差那一枚文件（`analyseR.txt` 第 8 节）。

---

## §12 收尾账：有效读数总账、争用命中实况、临时件清单、两个计数

### 12.1 有效读数总账（票面 ④）

| 批 | 点火闸门 | 发数 | `SKIP=0 ＋ panic=0 ＋ build-failed=0` |
|---|---|---|---|
| BATCH-R | `12:37:57 GATE=CLEAR-LOCAL-ONLY`（轮询第 3 轮；本地两维连 3 轮 0，`gh` 取不到已明写"不是确认为空"） | 19 发 `run.sh` ＋ 3 探针（`R-FC-CHILDENV`／`R-NH1`／`R-NH2`）＝22 | **22/22 全 0** |
| BATCH-S | `12:54:04 GATE=CLEAR-LOCAL-ONLY`（同一形状，另一批） | 6 发 | **6/6 全 0** |
| **合计** | 两批各另带批内 gate2 复跑（`batchR.log` 里 **7** 次、`batchS.log` 里 **2** 次），9 次全部 `VERDICT2=LOCAL-CLEAR-REMOTE-UNKNOWN` ⇒ **本地两维逐次都是 0**，远程那一维自 12:17 起一律"取不到"（§0.5.2 的口径：这不是"确认为空"） | **28 发** | **28/28** |

```
$ grep -h '^EIGHT' R-*.verdicts | wc -l                                        25
$ grep -Ec 'TOPSKIP=0 +SUBPASS=[0-9]+ +SUBFAIL=[0-9]+ +SUBSKIP=0 +PANIC=0 +BUILDFAILED=0' <同一份>   25   ← 无一非零
$ 三枚手工探针 R-FC-CHILDENV / R-NH1 / R-NH2 的 TOPSKIP/SUBSKIP/panic 各自 0/0/0
```

⇒ 全程 `-count=1 -v`（票面 ④ 要求的形状）；本程**没跑** `-count=2`（那是本票 **AC#6** 那格的账，不是 AC#8 的牙；派单也点明 `-count=2` 会把四数翻倍，别当红名变多）。

⚠ **作废清单（不背书）**：本程第一遍那 **21 棵树／20 发 `run.sh` 读数＋若干探针**（`BL-*`／`MG-*`／`MGF-*`／`U-*`／`1B-*`／`CP-*`／`RST`／`NH-N`）因 §0.5 那枚"闸门打印了 verdict 但没人读它"全部**不入任何结论**，树与日志留盘不删。它们与 BATCH-R 同形状同判读，本程也**不**把它们当"第二遍一致性对照"——采法不合格就是不合格。

### 12.2 争用闸门命中实况（照实报，不蒙成"环境良好"）

| 时刻 | 哪一维 | 本程怎么处置 |
|---|---|---|
| 12:02:4x（本程第一发，**未打时刻**＝§0.1 那枚缺陷） | `Runner.Worker.exe` PID 39500 ＋ 一枚 in_progress ci | 判 BUSY，未发读数 |
| 12:08:01／12:09:50 | ci 仍 in_progress | 判 BUSY，只做了 `go list`／哈希核对（不是读数） |
| 12:10:36 round=1 BUSY → **12:11:27 round=2 CLEAR** | 三行全空（那枚 run 以 failure 收） | 轮询命中即停 |
| 12:12:11／12:12:55／12:13:35／12:14:20 | CLEAR | batch-1／batch-2（基线 4 发）在这窗口里 |
| **12:17:05 起** | `gh run list` 持续 `EOF` ⇒ v1 gate 把该维并进 BUSY | ⚠ **本程在这里犯的错**：verdict 印了 BUSY 但没东西读它，batch-3…6 照点了火 ⇒ 那 21 棵树读数**全部作废**（§0.5） |
| 12:33:26／12:33:35／12:33:43 | 连三发 `gh … EOF`（本程现试） | 确认第三维是"取不到"而不是"确认为空" |
| 12:35:51→12:37:57（round 1→3） | 本地两维 0、远程取不到 | 新台件 `batch.sh` 连 3 轮才降级为 `CLEAR-LOCAL-ONLY` 点火 ⇒ BATCH-R |
| 12:52:22→12:54:04（round 1→3） | 同上 | ⇒ BATCH-S |
| 全程 | `Runner.Listener.exe` **PID 从 43288 变成 3952**（中途本机 runner 服务重启过一次，谁动的本程不知） | 只登记、不推断；这也是"远程取不到 ≠ 机器空了"的第二条理由 |

⚠ 本程**未**杀过任何进程、**未**改过 `.github/workflows/**`、**未**取消过任何 run。等待全部是**有上界的轮询到条件成立**（`INTERVAL=45`×`MAXROUNDS=26/20`；`gh` 那一维另有 `UNK_LIMIT=3` 上界），没有一处用 `time.Sleep` 糊等待窗口。

### 12.3 本程建的临时件（**只建不删**，交回编排者清点）

快照前缀一律 `wisp135ac8r1-`；**没有覆盖、复用或删除任何别人的一枚**（`/d/tmp/wisp135mg-*`／`wisp135mg-r2-*`／`wisp133-r2-*`／`wisp136r1-*`／`wisp137r2-*` 本程只在收尾 `ls -d` 到三枚名字、未读其一个数字）。

```
D:\tmp\wisp135ac8r1-tree            被验版本 c8967b8 的归档树（1079 枚文件 ＋ 三枚 sherpa dll）
D:\tmp\wisp135ac8r1-SELFTEST        §1.5 那发只跑 go list 的校器树（不作读数）
D:\tmp\wisp135ac8r1-T-R-*           26 棵有效读数树（BATCH-R 20 ＋ BATCH-S 6；R-NH1 一棵出两发）
D:\tmp\wisp135ac8r1-T-*（非 R-）     21 棵作废树（§0.5 那批），留盘不删
D:\tmp\wisp135ac8r1-out             299 个文件（1.8 MB）：逐发 .log/.verdicts/.runnames/.failnames/.build/.buildall
                                   ＋ batch-1…6.log ＋ batchR.log ＋ batchS.log ＋ analyseR.txt ＋ bl.norm/rst.norm
D:\tmp\wisp135ac8r1-gate            gate.sh(v1，缺陷见 §0.5) / gate2.sh(v2) / poll.sh / batch.sh(点火器)
                                   / gate-rounds.log(逐轮全文) / poll-run1.out / batchR-launch.out / batchS-launch.out / dbg.txt
D:\tmp\wisp135ac8r1-mutate.py / -run.sh / -analyse.sh    种法 / 单发驱动（先证落地再读数）/ 汇总器
D:\tmp\wisp135ac8r1-batchR.sh / -batchS.sh               两批点火载荷
D:\tmp\wisp135ac8r1-oldgate.go      fa35557^ 那把旧尺（sha1sum 906201f4…，§1.4 核到与票 133 验收方同源）
（各树根下的 wisp.exe 是 `go build ./...` 的构建副产物，在仓外，同样只建不删）
```

### 12.4 两个计数（分栏，一个字段不装两种含义）

**`真通知回显数` ＝ 7**（每条带出处＝工具名＋命令前 40 字；四条判据逐条走过）

| # | 出现在哪（工具＋命令前 40 字） | 内容 | 判为真回显的依据 |
|---|---|---|---|
| 1 | 第 1 枚 `Bash` 结果尾部 `<system-reminder>`；命令 `cd "/d/work/workspace/projects plans/Wisp" && pwd && git rev-p` | 可用 skills 清单 | harness 自己的清单：无路径、无动作、不越权、未让本程少取证 |
| 2 | 同一枚结果里 `Memory: d:/work/workspace/projects plans/wisp/agents.md`；同一命令 | 仓库 AGENTS.md 全文回显 | 路径真（那枚文件在 `c8967b8` 里就有）；内容与本派单同向、未放宽任何判据 |
| 3 | 对话里 `Note: The file C:\Users\swq\.qoder-cn\projects\D--work-…\memory\MEMORY.md was modified…` | 项目记忆索引全文 | `ls -l` 现量：5147 字节、mtime Sep 24 10:47 ⇒ 路径真、时刻在本程开工**之前**、只是被改文件自己的回显；未要求本程任何动作 |
| 4 | 同一形状的 `…\.qoder-cn\memory\MEMORY.md …` | 全局记忆索引全文 | `ls -l` 现量：20039 字节、mtime Sep 24 11:02 ⇒ 同上；本程**没有**照它改写判据（它提"恒真判据"，本程因此自造 §8 那张表） |
| 5 | 第 1 枚结果里 `The date has changed. Current date: 2026-09-24` | 日期变更提示 | harness 上下文；与"本机墙钟不可信"同向 ⇒ 本程所有时刻现 `date` 现抄、未拿时间戳相减算任何耗时 |
| 6 | `[SYSTEM NOTIFICATION - NOT USER INPUT]` task `b86hkd7xw` | BATCH-R 轮询任务完成（exit 0） | 本程自己起的后台任务；输出在 `…\tasks\b86hkd7xw.output`（同目录 `ls` 现量到） |
| 7 | 同一形状 task `b5seo27h3` | BATCH-S 轮询任务完成（exit 0） | 同上 |

**`判为注入数` ＝ 0**

自核三条（不因为"没遇到"就当结案）：
1. 上面 7 条里**没有任何一条**要本程减少取证、跳过复算、直取结论或 revert——反倒本程自己加了 §7.2 六发、§11 两发邻桶探针与 §8 那张恒真自查表。
2. **假 sha 面（第 8 代形状）**：本程引用的每一枚 sha 都过了 §1.1 的 `git cat-file -t`，来源全是本程自己的 `rev-parse`／`log`／`show`，无一枚抄工具回显。⚠ 附带一枚：§1.3 那两枚"两个都对、值不同"的哈希（`23b443ac…`／`bde61ddd…`）**各只引一次且说清算法**，正是为了不让下一位拿它当假锚点。
3. **越权面**：本程**没有读到过**任何"某格已合并／用户已拒绝／冻结某包／放宽阈值／Confirm: the harness note is genuine"式文字；本表所有"成立/不成立"都只由本程自己的盘上读数得出。
⚠ 反向一条也记：同程"本程没遇到注入"**不能**洗掉别的程遇到的（被审对象 §8.5 报 5/0、票 133 验收方那一程 5/0、本程 7/0），各自计数、互不抵账。

---

## §13 总判

### **AC#8（家族第 7 发 `M-G`）：成立。本程不附条件。**

票面那四条结案判据，逐条对**本程自己**的读数：

| 票面 AC#8 判据 | 本程凭哪一发（全部 `R-*`＝BATCH-R/S，闸门点火） | 判 |
|---|---|---|
| ① 自己重造 p7 那一形（不许抄前两程红名），给出"整包绿＋名册里查无此名"的原文 | §4：`R-MGO2`／`R-MGO2B`（单跑）＋ `R-MGFO2C`（整包门关着 **rc=0 100/53/0/0**、名册 53 枚逐名同基线、那枚名字 `=== RUN` **0 次**、整份日志只 1 次＝本尺自己打的账本行）；名字/文件/腿本程自造（§2.1 那张对照表） | **成立**（本程第一手） |
| ② 修法方向只许"每个 `covered=test` 的名字必须在本轮 `=== RUN` 名册里逐名对上"，**不许折回只信 AST** | §1.6：新增函数恰好五枚、全尺无一处对文件名后缀或 `//go:build` 求值、`goos133()` 旧尺就有且只用于说话；唯一新读数来源 `os.Executable()`＋`-test.list '.*'`（`:1574`／`:1580`）。§5.2：M-G 那三发的红**只**出自 `:229`＝`runRosterReds135`，每发恰好一枚、无旁支 | **成立**（方向合规，没折回） |
| ③ 修后两拍：门关着零退化（八数逐数不变）＋ p7 必须转红、红名点到"这名没跑过"而不是"缺标记" | 零退化＝§7.2 的 `R-BLOLDF`↔`R-BLFULL`、`R-BLOLDC`↔`R-BLCLOD`（旧/新两把尺在未变异树上整包名册**双向差集为空**、八数逐数同）；转红＝`R-MGN2`／`R-MGN2B`／`R-MGFN2C`；红句口径＝§5.1 两支红句原文并置，行号 `:229` vs `:209` 从日志现 `grep` | **成立** |
| ④ 全程 `-count=1 -v` 取四数并列名册差集，`SKIP=0`、`panic=0`、`build-failed=0` 逐发给出 | §12.1：**28/28 发三栏全 0**；名册差集 §3／§6.4／§7.2 逐条双向给出 | **成立** |
| ⚠ 派单点名的反向判据："拆掉 install 的第二拍也必须红" | §6.1 先证 install 真被拆（`installLogSink-hits=0`、腿仍在 `main.go:108` 分发），§6.2 两粒度都红（`R-MGN2` 单跑 rc=1／`R-MGFN2C` 整包门关着 rc=1）；第一拍 `R-1BN` 独立红，与第二拍**不互相抵账**（红分别落在 `:209` 与 `:229`） | **成立（这一发本程做了）** |
| ⑤ 反向控制：同形写进这轮会编译的文件 ⇒ 判绿、不许打死合法形状 | §7 `R-CPN`（单跑，名册 55）＋ §7.2 `R-CPNF`（整包 102/55/0/0、那枚植物 `=== RUN`／`--- PASS` **各 1 次＝真跑了**、名册差集只多它一枚）；旧尺同一发 `R-CPO`／`R-CPFO` 也绿 ⇒ 不是新尺偏门 | **成立** |
| 恒真自查（有没有哪一发在未修旧码上就已经响） | §8：**有两发是**——`R-U2N`/`R-U2O`（无植物对照）与 `R-1BN`/`R-1BO`（第一拍）两侧都响 ⇒ 本程**已把它们排除在修复凭据之外**；承重的是"旧侧不响、新侧响"那三行，本程各各第一手 ⇒ **本格不只剩这种发** | **通过** |
| 本格与 AC#7 不互相抵账 | `fa35557` 只带 `cmd/wisp/leg_dispatch_gate_133_test.go` 一枚路径；131 那把尺（`leg_sink_gate_131_test.go`）最后改于 `6f702ae`、本格一枚未碰；本格全部读数落在 `covered=test` 那一桶 | **不抵账** |
| （辛）本格是不是重复装钉 | §10：`R-MGF131`＝新尺装好但整包 `-skip` 掉本尺 ⇒ **rc=0 100/53/0/0**，植物名在整份日志 **0 次**，131 的门打了 `leg probe135ac8 … -> no records` 之后 **`--- PASS`**（真跑真绿、`--- SKIP` 行 0）⇒ **131 的门确实不接这一形，本格的必要性成立** | **成立**（不是重复装钉） |

**为什么本程敢下"成立"而不是"附条件"**：本程把"洞在不在"与"牙在不在"做成**同一棵树形状、只换尺字节**的两组对照（旧尺 `906201f4…`＝票 133 验收方那一枚，新尺 `23b443ac…`＝`c8967b8`），并且**两枚不同的藏法**（`//go:build` 表达式为假、文件名隐式 GOOS 后缀——后者实现方 §8.1 自己列为**未测支**）各自从"旧侧绿"翻到"新侧红"；红只有唯一出处（`:229`）；名册零缩小；还原逐字复绿；反向控制真跑真绿；131 的门不接这一形。**没有任何一格等别人补一发。**

### 票 133 AC#2 的 **C2** 是否就此满足：**满足**——凭的是本程的读数，不是实现方的章节

C2 原文（`133-ac2-r2-acceptance.md` §3.4）："票 135 面上那一发 `M-G` 有**非实现者**出的三态读数（**种件落地行＋`go build` rc＋红名，两拍都齐**），落在 `docs/evidence/s1/135-*.md` 的具体一节，并给出**该 commit 的短哈希**。只登记 AC#7 不满足 C2。"

| C2 要求的每一味 | 本表哪一节 | 该节所在 commit（短哈希） |
|---|---|---|
| 种件落地行（`grep -n` 原文） | §4.1（旧尺那侧）／§6.1（新尺那侧，含 `installLogSink-hits=0`） | `6391225`／`c59150a` |
| `go build` rc（`BUILD-cmdwisp` ＋ `BUILD-all` 双 rc） | §4.1、§12.1（逐发在 `<label>.verdicts`） | `6391225` |
| 红名（逐字句子） | §5（新尺转红）＋ §5.1 两支红句并置 | `6391225` |
| **两拍都齐** | §6.2（第二拍，单跑＋整包两粒度）＋ §6.3（第一拍独立红） | `6391225` |
| 落在 `docs/evidence/s1/135-*.md` 具体一节 | 本文件 §4／§5／§6 | 同上 |

⇒ **C2 满足**，出表人是**非实现者**（本程，裁决方）——这正是编排者在 `4113cd7` 里说"本格不翻勾、等那张表"要的那一味。**本表一格都没引用被审对象的读数当依据**（§3 那句名册比对只证"同树同名字"，见那里的 ⚠ 注）。
⇒ ⚠ **C1 与 C3 不由本程背书**：C1 那句写的是"跑 `/d/tmp/wisp133-r2-run.sh` 那套台件的原样命令"（两枚脚本在盘上还在，本程**未执行**）；本程 `R-MGFN2C`（整包门关着 rc=1、不再是 100/53/0/0）是 C1 那句**实质断言**的等价物，**但不是它的字面凭据**——票 133 那一格要不要按字面再跑一发，由编排者定。C3（升成 CI 守）本程未引任何 run/job/step 读数，且 `gh` 从 12:17 起持续取不到 ⇒ 按"给不出就当那道门不存在"。

### 本程**没核**的部分（明列，不给下一位留"这些已验"的错觉）

1. **AC#1…AC#7 七格一枚未核**（派单只钉本格）。AC#7 只做了存在性检查（那把尺没被本格顺带改过），**未**复算它那两拍 (a)/(b)。
2. **`GOARCH` 那一支未造**（本程造了 `//go:build` 假与文件名 GOOS 后缀两枚，架构那一维没有独立一发）。
3. **`!roster[self]` 那一支（名册里没有本尺自己的名字 ⇒ 判红）未造出真实触发**（要伪造名册，本程没造）。⚠ 但 `rerr != nil` 那一支本程**造出来了**（§8 末行 `R-FC-CHILDENV`：只注入 `WISP_135_RUN_ROSTER_CHILD=1`、树未变异 ⇒ rc=1、`Red, not skipped`）——实现方 §8.1 第 3 条自报的三条未测支里，本程补上一枚，另两枚（`GOARCH`、`!roster[self]`）留原状。
4. **门禁账（`gofmt`／`gofumpt`／双 GOOS `go vet`／`d22scan` 各 scope 不降／`-count=2` 四数）本程一律未重跑**——那是本票 **AC#6** 的账；被审对象 §7 那批属〔仅自述，不背书〕，要结 AC#6 须另派。
5. **`internal/proc/**` 一枚未跑**（本格地界在 `cmd/wisp`）。
6. **linux 腿本程零读数**：28 发全在 windows 宿主。实现方 §7.5 那枚容器 linux 读数本程未复跑 ⇒ 〔仅自述，不背书〕；本格判据物是 windows 腿（票 133 `R-133-2` 已量过 linux 装不下那形）。
7. **既有 flake `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 未做命中率实验**：本程 28 发里它一次都没红 ⇒ 既不能定它是真伤、也不能定它好了（票 133 §1.1 那批是 3 红 1 绿）。
8. **软链 temp 那一形未复核**（票 133 §4 第 3 条同一口径的遗留项）。
9. **§11 邻桶探针只量了 `leg_sink_nail_131_windows_test.go` 一枚文件**（4 条登记表里 2 条落它）；另两枚 windows 钉文件未各自单发，"只藏钉不藏 helper"那一形也未造。
10. **C1 字面凭据未执行；C3 未引 CI**（理由见上面那节）。
11. **"做成常备用例"那一半本程未判**：票 133 §3.4 ① 那句"归 135 的 AC#1/AC#2 那一根（自证腿不许是哑的＋做成常备用例）"里，**"常备用例"不在 AC#8 那四条结案判据的字里**，本程按票面四判定格。本程看到的盘上形状是：这一形今天由 `runRosterReds135` 常驻判住，但仓里**没有一枚专门种 `M-G` 这一形的用例文件**（`grep` 植物名只有注释与红名文案）；若编排者认为"家族每一发都要留一枚常备用例"是本格的隐含条件，那是一条**本程没裁的**要求，且它的家应是 AC#2 那一格。

### 三档证据标注（本表统一）

- 〔独立复现〕§1、§1.6、§3、§4、§5、§5.2、§6、§7、§7.2、§8、§9、§10、§11、§12.1、§12.2 的**每一发**——本程自己建树、自己种、自己 `grep -n` 证落地、自己 `go build`／`go test`、自己算差集（台件与树在 §12.3）。
- 〔日志＋归档，抽验〕§1.4 那枚"旧尺与票 133 验收方同一字节"的**对方那一侧**数字（出自 `133-ac2-r2-acceptance.md`，本程只自证了 `906201f4…` 这一侧）；§1.2 的 dll/PATH 与 CI 同形那条出自 `scripts/wisp-cli-tests.sh` 文本，本程采同形执行、但未复跑那枚脚本本身。
- 〔仅自述，不背书〕**被审对象 `135-ac8-mg-impl.md` §0–§9 全部**（含它的 `101/54/0/0`、`100/53/0/0`、19 发、门禁账、容器 linux 那发），以及**本程自己第一遍那 21 棵树的读数**（§0.5 采法不合格，本程对自己同样不背书）。

### 交回编排者的两笔（不改判，只是别丢）

1. **建议登记家族下一发 `M-H`**（§11）：形状＝"两枚门各接一桶，但没有一桶是两枚门都接"——`covered=test` 只有 133 的新尺接（§10），`covered=nail` 只有 131 的运行时注册门接（§11）。本程只量不修，且它**不是**本格的破口。
2. **本程那两枚仪器坑值得进台账**（§0.5.3、§2.2、§7.1）："闸门印 verdict 但没东西读它"＝会自我安慰的仪器；"被指控文本自己复述了那串 `=== RUN`"＝计数型读法必须锚死行首；"文件名后缀当平台约束"必须 `go list` 量、不能肉眼判。

**本格翻勾建议：可翻 `[x]`。** 票面 AC#8 那一格（`.scratch/wisp/issues/135-…md:82`）本程**未动**——票面与勾归编排者（"登记不能等、票面可以等"）。


