# 票 135 AC#8（家族第 7 发 `M-G`）实现方证据 —— 变异自证三态 ＋ 门禁账

**执行方**：`worker-ticket135-ac8-r2`（**接续程**，非 `worker-ticket135-ac8-mg`）
**被验版本**：commit **`8369b24`**（`8369b24e20eb85b5c9afb45198499dc4b9d9d804`，含上一程的 `fa35557`）
**本格判据来源**：票 135 面 AC#8（09-24 00:0x 编排者追加）＋ 票 133 第二任验收方的
`docs/evidence/s1/133-ac2-r2-acceptance.md` §1.11／§1.12／§3.3／§3.4（`R-133-9`＝`M-G` 的出处）
**本程地界**：一枚新建证据文件（本文件）。**生产码零接触、上一程那把尺零改动**（§1.4 字节为证）。
**本程不翻任何勾**（AC#8 由编排者按本表翻）。

## 0. 争用闸门（每一批读数之前跑一次，逐轮原样登记）

### 0.1 为什么本程把这道闸门当第一条规矩

这台机器上 **self-hosted CI runner 与编队同机**：`.github/workflows/ci.yml:538` 的 `slo-full`
那一枚 job 的 `runs-on` 是 `[self-hosted, wisp-slo]`。编排者在 **10:50:36 +08** 推了一次 dev
（推的正是被验那枚 `8369b24`，`gh run list` 的 `headSha=8369b24e20eb85b5c9afb45198499dc4b9d9d804`），
当场起了一枚 `ci` run ⇒ **"编队里没有别的代理"不等于"机器是空的"**。

### 0.2 固定动作（三行）与本程用的台件

派单写死的三行：

```
tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'
tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'
gh run list --limit 3 --json status,conclusion,name,createdAt
```

本程把它们包成一枚**有上界的轮询**：`/d/tmp/wisp135mg-r2-gate/poll.sh`
（`INTERVAL=45`s × `MAXROUNDS=26` ≈ 19 分钟一枚上界；命中即 `BUSY` 继续等，
全空即 `CLEAR` 退出 0；到点上界仍忙则退出 3 并照实报。
**没有用 `time.Sleep` 糊窗口**——等的是"条件成立"，且带次数与时长上界。
**没有杀过任何 runner 进程、没有改过 `.github/workflows/**`、没有取消过任何 run**：只观察。）

### 0.3 逐轮原始读数

**手工头两轮（本程自己跑的三行原样）**

```
$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:51:50 +0800
Runner.Listener.exe          50052 Console                    1     86,552 K
Runner.Worker.exe            36068 Console                    1     98,024 K      ← Worker 活着 = 远程 run 在本机落地
$ tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'          （无输出，rc=1）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
[{"conclusion":"","createdAt":"2026-09-24T02:50:36Z","name":"ci","status":"in_progress"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:24:40Z","name":"ci","status":"completed"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:14:55Z","name":"ci","status":"completed"}]

$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:55:01 +0800
Runner.Listener.exe          50052 Console                    1     86,936 K      ← Worker 已退，Listener 常驻（本机 runner 空转）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
（同上：仍有一枚 in_progress 的 ci）
```

**轮询逐轮**（全文在 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log`，下面按轮摘 VERDICT 行＋该轮时刻；
每轮的三行完整输出都留在那枚日志里）

```
---- round=1 at=2026-09-24 10:56:55 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
---- round=1 at=2026-09-24 10:57:10 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
```

⚠ 本节此刻**只有这两轮**：本程后面每一批发作之前都按 0.2 那三行复跑一次，
VERDICT 与该轮时刻记在**对应读数那一节**里（§3 起），不在这里预填。
**读到本文件时若 §2…§7 任一节缺失或未结，那一节的读数就是没采**，不是"采了没写"。

### 0.4 判读规则与本程的处置

- `Runner.Listener.exe` **常驻不算命中**：它是等活的进程，只有 `Runner.Worker.exe` 代表"真有 job 在这台机器上跑"。
  这一条本程按派单原文的形状写（派单第 2 节：「有 Worker 就说明远程 run 活着」）。
- **`gh run list` 里有一枚 `in_progress` 的 ci＝命中**，即使该 run 此刻落在 GitHub 托管的 runner 上
  （本程现量过那枚 run 的 job 分布：`gh run view 35948947994 --json jobs` ⇒ `lint`/`slo-full`/
  `lint-frontend`/`test-core`/`slo-smoke` 已 completed、`test-windows` in_progress，
  `runs-on: windows-latest` 即 GitHub 托管）——
  **本程不赌"它不在这台机器上"，只要闸门命中就重采**。
- ⚠ **本程自己犯过一次、并且当场作废的一发**（不藏）：`§3` 的第一发 `R2-B1`
  在 `GATE-CHECK 2026-09-24 11:02:20 worker=1` 之下取的——编排者 11:01:39 又推了一枚 dev
  （`45623e4`），当场起了新一枚 `ci` run 与一枚 `Runner.Worker.exe`（PID 38372，本程 11:02:47／11:03:01 两次
  `tasklist` 都还在）。⇒ **那一发按派单口径不作数**，本程重新按闸门等到 CLEAR 之后再发整批；
  `R2-B1` 的原始日志留在盘上（`/d/tmp/wisp135mg-r2-out/R2-B1.log` 等）当"这一形在未作废之前长什么样"的对照，
  **但本表任何结论都不引用它**。
- ⚠ 一条本程自己造出来的噪声要交代：**一旦本程开始 `go test`，第 2 行的 `^go\.exe|compile\.exe` 必然命中**
  （那是本程自己的编译）。⇒ 闸门只在**每批读数开始之前**跑一次，批次内的命中不属争用。

**结论（§0）**：闸门**确实命中了**——头两轮一枚 `Runner.Worker.exe`、加上一直挂着的 `in_progress` ci run。
本程因此**先等再采**，等待用的是有上界的轮询，逐轮留档。
〔独立复现〕本程自己跑的三行与 `poll.sh`；日志 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log` 可逐轮重看。

## 1. 锚点、取件、以及"本程没有动那把尺"

### 1.1 锚点（开工第一步自量，不抄派单）

```
$ git rev-parse --short HEAD            8369b24          （开工时；本程第一枚 commit 之后才是 087ef2e）
$ git rev-parse --abbrev-ref HEAD       dev
$ git cat-file -t 8369b24               commit           ← 逐枚现核，见 §8 末
$ git cat-file -t fa35557               commit
$ git merge-base --is-ancestor fa35557 8369b24   真（rc=0）⇒ 被验版本含上一程那把尺
```

被验版本按派单钉在 **`8369b24`**。⚠ 本程跑到一半时 dev 又前进了两次
（编排者 11:0x 提 `45623e4`＝`docs(Q-45 自决,A152)` 并推了，那是他们改 `docs/PLAN.md` 的 commit，
与本票无关、本程未读未改未提交）。**本程全部读数仍在 `8369b24` 的归档树上**，不受影响。

### 1.2 取件（仓外快照，零 worktree／零 checkout／零工作树内编译）

```
$ mkdir -p /d/tmp/wisp135mg-r2-tree && git archive 8369b24 | tar -x -C /d/tmp/wisp135mg-r2-tree
   ARCHIVE-OK
$ sha1sum  /d/tmp/wisp135mg-r2-tree/cmd/wisp/leg_dispatch_gate_133_test.go   23b443ac34787aa9ea60e181b9b8c789f468cd56
$ git show 8369b24:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum          23b443ac34787aa9ea60e181b9b8c789f468cd56
$ git show 8369b24:cmd/wisp/main.go                     | sha1sum            bba113648726af8d47371a00c17048aa94b1055d
$ sha1sum  /d/tmp/wisp135mg-r2-tree/cmd/wisp/main.go                          同上（逐字同）
$ cp <repo>/third_party/sherpa-onnx/*.dll <tree>/third_party/sherpa-onnx/
     onnxruntime.dll 17799168 / sherpa-onnx-c-api.dll 4605952 / sherpa-onnx-cxx-api.dll 259584
     （三枚 dll 不在归档里＝被 gitignore，缺了测试二进制 0xc0000135，照票 133 验收方 §0.3 的同法补）
```

**没有读脏工作树当被验内容**：工作树里那枚尺与 `8369b24` 的 blob 逐字同（`23b443ac…`，本程
在仓外只读地比过一次），但**所有读数取的都是归档树**，且每发一棵**新**树（`cp -r` 自
`wisp135mg-r2-tree`，树已存在即 `TREE-EXISTS-REFUSING` 拒跑）⇒ 结构上不存在"上一轮的 `.bak`
盖回当前树"那一类作废读数。

### 1.3 旧尺（＝`M-G` 洞的载体）也自取了一份，并核到与前一手同源

```
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go > /d/tmp/wisp135mg-r2-oldgate.go
$ sha1sum /d/tmp/wisp135mg-r2-oldgate.go        906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   906201f4a3995d10b0a65910aa4cc68e1f745a9d
```

⇒ 这枚 `906201f4…` 与票 133 第二任验收方 §0.3 记录的
`git show 048a9e4:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum` **逐字同**。
也就是说：**本程拿来当"未修那一侧"的旧尺，与造出 `p7` 那一发、量到 `100/53/0/0` 全绿的那枚尺，是同一枚字节**
⇒ §4 的旧尺对照发才配当"洞复现"的证据，而不是另一棵树的另一回事。

### 1.4 本程对那把尺做了什么：**零**

- 归档树的尺 `23b443ac…` ＝ 工作树的尺 ＝ `8369b24` 的 blob ⇒ **一字未动**。
- 本程**不判断需要改这把尺**：派单第 1 节那条"若必须改才许结案就停手报回"的停手线**没有被碰到**。
  它引用的证据路径 `docs/evidence/s1/133-ac1-ac2-instrument.md` 本程也现核存在（`ls` 命中）。
- 所有变异都落在**快照副本**里（`/d/tmp/wisp135mg-r2-T*`），生产码与测试码在工作树里零改动。

### 1.5 工作树里那一枚不是本程的改动（登记，不提交）

开工第一刻 `git status --porcelain` 是**空**。本程 §0 提交之后，工作树出现
`M docs/PLAN.md`（4 增 4 删，`D1–D46`→`D1–D47`／`C1–C31`→`C1–C32` 那类指针计数修正），
**不是本程写的**；本程未改它、未提交它，随后它由编排者自己以 `45623e4` 入库。
本程每一枚 commit 都带显式 pathspec（`git add -- docs/evidence/s1/135-ac8-mg-impl.md` ＋
`git commit -q -F - -- <同一路径>`），`git diff --cached --name-only` 每枚都只有本文件。

## 2. 被读的那把尺：形状登记（读数在 §3 起）

`8369b24` 里那枚 `cmd/wisp/leg_dispatch_gate_133_test.go` 的 M-G 判据（本程逐枚 `grep -n` 现量行号）：

```
291    coveredBy []string                                   ← leg133 带上"这行账本 credit 了谁"
1538   const rosterChildEnv135 = "WISP_135_RUN_ROSTER_CHILD" ← 子进程防重入闸
1543   const rosterReadTimeout135 = 2 * time.Minute          ← 上界，不是 sleep
1548   var rosterName135 = regexp.MustCompile(`^Test[A-Za-z0-9_]*$`)
1570   func compiledRunRoster135() (map[string]bool, int, error)
1601   func headRunes135(s string, n int) string
1614   func flagValue135(name string) string
1632   func (p *pkg133) runRosterReds135(...) []string       ← 判据 (c)
1663   func (p *pkg133) runRosterDisclosure135(...) string   ← 替换掉那句 guessed 的失明披露
```

**读数来源是运行中那枚二进制，不是解析器**（这条就是票面禁止"折回只信 AST"的那一条）：

```
1574   exe, err := os.Executable()
1580   cmd := exec.CommandContext(ctx, exe, "-test.list", ".*")
```

全文件里 `//go:build` / `_linux_test.go` / `runtime.GOOS` 只出现在**注释与红名文案**里
（`:462`、`:1560`、`:1569`、`:1651`、`:1912` 的 `goos133()`），**没有任何一处代码去求值构建标签**
⇒ 它没有把第二种"只问签名"换成第三种"只问字节"，而是引入了第三种**独立**读数。本程判：方向合规。

`runRosterReds135` 的两条 fail-closed 也现量到（这两条决定 §3 起怎么读它的红）：
`rerr != nil` ⇒ 直接红（"读不到名册"不是"没人欠覆盖"）；名册里没有本尺自己的名字 ⇒ 也红。

⚠ 本程在这把尺里读到一处**它自己写明的不对称**，登记、不改、不替它辩护：
`covered=test` 那一桶**红**，而 registry 钉（`covered=nail`）落在名册外时**只披露不弄红**
（理由＝那四枚住在 `*_windows_test.go`，在 linux 形上必落名册外，弄红＝在本平台主张一笔覆盖债）。
⇒ 本程的 §5 反向控制与 §7 门禁账都只测 **windows 腿**；linux 腿的运行读数**本程取不到**，
写进 §8 的"没做的档"，不当已验。

〔独立复现〕上面每一行行号与 sha 都是本程自己在 `8369b24` 的归档树／仓外只读命令上量的。
〔日志＋归档，抽验〕"上一程为何这么修"的理由文本来自 `fa35557` 的 commit message 与本文件源码注释，
本程只核了字节与行为，不背书它的措辞。
