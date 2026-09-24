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
   secret_daroot_119b_test.go slo_other.go versioninfo_other.go]           ← 植物在 IgnoredGoFiles 里
```

⇒ "这枚用例今天编不进去"不是本程从文件名后缀**推**出来的，是 `go list` **量**出来的：它在 `TestGoFiles` 里 0 命中、在 `IgnoredGoFiles` 里点名。这枚读数同时是本程 §4 那几发的**独立对照仪器**（它不属于被审的那把尺，也不是 go/parser 走源码）。

〔独立复现〕§1 每一行都是本程自己跑的（仓外只读命令＋`go list`）。
〔日志＋归档，抽验〕"两程各自引的哈希都对"这条的**对方那一侧**数字出自被审对象 §9 与编排者 log，本程只自证了 `sha1sum` 三向。
