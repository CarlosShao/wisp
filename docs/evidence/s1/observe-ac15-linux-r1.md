# observe AC#15 — Linux 侧的分母·率·与 2s 上界的牙（r1）

agent: `worker-observe-ac15-linux-r1`
face: `docs/evidence/s1/observe-ac15-poll-r1.md` §7.5 第 2/3/7 条（本程就是去还那三条债）
角色: **测量程，不是实现程**。`internal/observe/**` 一字未动（改前改后 md5 相同，见 §0.5）；
所有注入只发生在**仓外快照**（`D:\tmp\observe-ac15-linux-r1-s1\snap-*`，只建不删）。
本程唯一写进仓库的路径＝本文件。

要收的三格（派单编号）：
- **(i)** 2s 上界是在 **Windows** 上现量推的，而本包在 CI 上**只有 ubuntu 一枚测试分母**（台账 `A213⑤`）——先复验那句分母，再在 Linux 上把上界的牙重新量一遍；
- **(ii)** 改后**没有忙窗读数**（poll 程 §7.5 第 2 条自己认的）；
- **(iii)** 改后 `p50 == max`＝**重开路径可能一次都没走过**（同条第 1 款的另一读法）——本程把它走一遍并报出重开枚数。

---

## 0. 锚点、台件、与挂载证明

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-25 10:03 +0800            （§0-§3 写于 09:46-10:05，各小节下面各自贴取数时刻）

$ git rev-parse --short HEAD
888bbd5                            （快照 pin＝39cc381，取于 09:52）
$ git diff --name-only 39cc381..HEAD -- internal/observe/ | wc -l
0                                  ⇒ 快照内容＝写本程时 HEAD 的 internal/observe 内容
$ git status --porcelain -- internal/observe/ | wc -l
0                                  ⇒ 工作树里本包没有被别人改动的未提交量
$ git log --oneline -1 -- internal/observe/
333dfe3 evidence(136 AC#15 r1 §4+§5)   ⇒ 本包最后一枚 commit 仍是 poll 程那两枚之一
```

台件（全部现跑，不抄简报）：

```
$ docker images | grep -i golang
golang:1.27                  3680233e3204       1.31GB          327MB   U
golang:1.27-alpine           4cb7ac979db5        381MB         75.4MB
⇒ 派单指定的 `golang:1.27` 在本地存在，本程用的就是它（没有拿 alpine 顶替）。

$ docker exec ... uname -a
Linux 6.6.114.1-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Mon Dec  1 20:46:23 UTC 2025 x86_64
$ cat /etc/os-release | head -1       PRETTY_NAME="Debian GNU/Linux 13 (trixie)"
$ nproc                               12
$ go version                          go version go1.27.1 linux/amd64   （＝go.mod 的 toolchain go1.27.1）
$ go env GOPROXY                      https://proxy.golang.org,direct   （官方 proxy，没有镜像站）
$ ls -ld /tmp ; getconf TMPDIR ; go env TMPDIR
drwxrwxrwt 1 root root ... /tmp        （真目录、不是软链——台账记过"容器软链 TMPDIR 会自己缩小分母"那一形，本程排除它）
```

⚠ **台件要说实话的一条**：这台 laptop 的 Docker 是 **WSL2 内核上的 Debian 容器**，
不是 GitHub 那枚 `ubuntu-latest` runner。**CPU 数只量能核到的一半**：本机 daemon 现量
`docker info` ⇒ `linux 12`＝**12 枚 vCPU**；`ubuntu-latest` 公开规格是 4 枚 vCPU——
**那一枚本程在本机核不到、也没有核**（要核得去读 run 日志里的 runner 规格，属另一格）。
内核确实不同（本程＝`6.6.114.1-microsoft-standard-WSL2`，Debian 13 userland）。
本程交的是"**Linux 这一侧的调度形状**"，**不等于** CI runner 的形状；
**这一条现在就记在账上**，文末"本程没有测什么"再逐条展开。

### 0.4 挂载陷阱：**两种失败形状本程都真撞到了**，不是引台账

派单点名的这一坑（`docker -v` 静默挂空且 rc=0 ⇒ "干净的绿其实测的是空目录"）**在本机当场成立**。两发都贴原文。

**形状一：MSYS 路径转换开着（默认）——`-w` 被换算成宿主目录，docker 直接报错。**
```
$ docker run --rm -v "/d/tmp/does-not-exist-trap-test:/repo" -w /repo golang:1.27 sh -c '...'
docker: Error response from daemon: the working directory 'C:/Users/swq/.qoder-cn/bin/git/repo'
is invalid, it needs to be an absolute path
```
（`/repo` 被换成了"当前 git 目录下的 repo"。这一发**响亮失败**，不算假绿。）

**形状二：转换关掉（`MSYS_NO_PATHCONV=1`）＋ Windows 式源路径，但源目录不存在——Docker Desktop 给你**
**造一个空目录、容器里 `ls` 得到 `total 4` 的空 `/repo`、rc=0。**这就是假绿的那一发：**
```
$ MSYS_NO_PATHCONV=1 docker run --rm -v "D:/tmp/does-not-exist-trap-test:/repo" -w /repo golang:1.27 ls -la /repo
total 4
drwxrwxrwx 1 root root 4096 Sep 25 01:47 .
drwxrwxrwx 1 root root 4096 Sep 25 01:47 ..
rc=0
$ ls -d /d/tmp/does-not-exist-trap-test      ← 宿主上那枚目录被真的创建出来了（本程只建不删，留在盘上）
```
⇒ **任何"容器里全绿"的读数，必须先证明挂载非空。**本程的三条证明（同一容器、同一命令串里跑）：
```
$ docker exec wisp-obs-ac15-linux-r1 sh -c 'ls /repo | wc -l ; go list ./internal/observe ; md5sum ...'
20
github.com/CarlosShao/wisp/internal/observe      ← go list 在挂载目录里解析出真包（空目录产不出这行）
19de4896ea4b5122fbbc31418a85c59c  /repo/internal/observe/window_wait_136_test.go
19de4896ea4b5122fbbc31418a85c59c  D:\work\workspace\projects plans\Wisp\internal\observe\window_wait_136_test.go
87896a7721ba546e3d2c10aa9eb21388  /scratch/snap-post/internal/observe/sampler_settle_coverage_136_test.go
87896a7721ba546e3d2c10aa9eb21388  /repo/internal/observe/sampler_settle_coverage_136_test.go
```
容器内 `/repo` 是**只读挂载**（`:ro`）⇒ 本程结构上写不进仓库；快照树挂在 `/scratch`（可写，仓外）。

### 0.5 挂载与快照的另一枚自捉（同形第二发，就地认）

第一版 `queue.sh` 在**重建容器时把 `/out` 那枚挂载换成了 `/scratch`**，于是容器里 `/out/run.sh` 不存在，
整批 11 枚 batch **秒退**、每枚都印 `sh: 0: cannot open /out/run.sh`——**如果我的计数脚本把"没有读数"当成
"0 命中"，这一批就会以全绿交出去**。实际形状是 `run.sh` 根本没能启动、`meta.log` 里没有任何
`RUN=` 行 ⇒ 计数器产不出数，不是产出 0。处置：`ln -sfn /scratch/out /out` 后整批重跑（`QUEUE v2`），
第一版日志留在 `out/queue.log` 里不抹。
⇒ **这条正好是派单"仪器会说谎"那一节的第四发**：假绿不必是空目录，也可以是一个**没跑起来的取数器**。

### 0.6 争用归因（给同时在跑的验收程）

- 容器名 `wisp-obs-ac15-linux-r1`，**没有 `--cpus` 限制**＝可用满 daemon 的 12 枚 vCPU（`docker info` 现量 `linux 12`）；
- 忙窗那一段**在容器内起了 `nproc`＝12 枚纯 CPU 自旋**（`sh -c "while :; do :; done"`，起止时刻与逐批负载见 §4）；
- 本程一次只跑一枚重活（串行队列，无 `&`／`-parallel`／`-cpu` 之外的并发）；
- daemon 上还有别人的常驻容器（`docker ps` 现量）：`clipsync*`／`union-proxy` 等 7 枚——**台账提醒过"容器负载不落宿主进程名单"**，
  所以本程每批的负载标签取的是**容器内 `/proc/loadavg`**，不是 Windows 的进程名扫描（两味不同，不可与 poll 程的"九枚名单 0 枚"直接比）。

---

## 1. `A213⑤` 复验：本包在 CI 上确实只有 ubuntu 一枚**测试**分母

派单要求"re-verify it"。复验＝自己读那三枚文件，不引台账句子。

| 断言 | 现量出处 | 判定 |
|---|---|---|
| `test-core` 跑在 ubuntu，且执行 `--scope=core` | `ci.yml:224` job `test-core:`／`ci.yml:225` `runs-on: ubuntu-latest`／`ci.yml:288` `run: bash scripts/portable-tests.sh --scope=core` | ✅ |
| core scope 里**有** `internal/observe` | `scripts/portable-tests.sh:175`（scope 数组含 `./internal/observe/...`）＋ `:140`（`core_pin` 里逐字一行 `github.com/CarlosShao/wisp/internal/observe`） | ✅ |
| `test-windows` 跑在 windows，且执行 `--scope=windows` | `ci.yml:334/335`（`windows-latest`）＋ `ci.yml:458` | ✅ |
| windows scope 里**没有** `internal/observe` | `scripts/portable-tests.sh:184-188`（八枚包，无一枚 observe）＋ `:151-160` `win_pin` 九行里无 observe | ✅ |
| 没有第三枚**跑测试**的 job 问津本包 | `ci.yml` 全部 job 现列：`lint`(66 ubuntu)／`test-core`(225)／`test-windows`(335)／`slo-smoke`(480 windows)／`slo-full`(538 self-hosted)／`lint-frontend`(612 ubuntu)。`lint:121` 是 `go vet ./...`＝**编译并类型检查测试文件，不执行用例**；两枚 slo job 跑的是 `wisp slo` 二进制、不是 `go test` | ✅ |

⇒ **`A213⑤` 成立**：`internal/observe` 的用例在 CI 上**只在 ubuntu-latest 的 `test-core` 一步被执行**，
Windows 那一侧的 scope 结构性不含本包。poll 程（§7.5 第 3 条）与本程读到的同一张盘。

---

## 2. Linux 四数＋名册：与 Windows 基线**同值**，且同值是有结构理由的

命令与派单指定逐字同形：`go test ./internal/observe/ -count=2 -v`（未改动的工作树，只读挂载）。
两发独立复跑（中间隔了一次容器重建，见 §0.5）：

| 发 | 取数时刻（容器内） | `=== RUN` | PASS | FAIL | SKIP | `^panic:` | rc | 包时 | 不同名 |
|---|---|---|---|---|---|---|---|---|---|
| L-1 | 01:49:01→01:49:12 | **142** | **142** | **0** | **0** | **0** | 0 | 8.294s | **71** |
| L-2（重建后） | 02:03:09→02:03:19 | **142** | **142** | **0** | **0** | **0** | 0 | 8.749s | **71** |

Windows 基线（poll 程 §1.1／§8，同命令）：**142 / 142 / 0 / 0**，`^panic:` 0，**71 枚不同名**。
⇒ **Linux 与 Windows 四数逐枚同值；名册枚数同值。**

### 2.1 派单那句"Linux 名册本该更大"在本包**不成立**——理由是结构性的，不是没看

派单的预期是："`*_other_test.go` 与 `//go:build !windows` 只存在于 Linux ⇒ Linux 名册预期更大"。
本包现量：

```
$ grep -rn "go:build" internal/observe/*.go
internal/observe/treemetrics_other.go:1://go:build !windows      ← 全包唯一一枚平台约束文件
$ ls internal/observe/*other*test*
No such file or directory                                        ← 本包没有 *_other_test.go
$ grep -rn "NullTreeReader" --include=*.go .                     （除声明处外）0 命中
$ grep -h '^func Test' internal/observe/*_test.go | wc -l        71   （15 枚 _test.go，与名册同值）
```
`treemetrics_other.go` 里只有一枚 `NullTreeReader`（返回 `ErrTreeUnsupported` 的占位读数器），
**零枚用例、零个引用者**——它的作用只是让本包在 ubuntu 上**编译得过**（文件头自己写着
"must compile (ubuntu test-core)"）。⇒ 所以 Linux 侧**没有**任何一枚"Windows 不跑"的用例，
名册同值＝71 是**两 OS 共同的 71**，不是两枚碰巧相等的数。
⇒ **"两向差集"这一味本程因此给不出非空答案**（差集按定义是空），这条在 §5 里单独标成"未测／结构性不适用"，
不写成"已核对差集为空＝守恒"那种听起来更强的话。

### 2.2 SKIP 与并发的形状（Linux 侧现量，与 poll 程同形）

- `grep -rn 't\.Skip' internal/observe/` ＝ **3 处，全部在注释行**（`sampler_test.go:36`、`:357`、`window_wait_136_test.go:94`）⇒ 本包**结构上产不出 SKIP**，
  所以 Linux 的 `SKIP=0` 不是"被跳过了没看见"（台账那条"变绿还是被跳过"的分界在这里是**可证的**）。
- `t.Parallel` 全包 **0 枚**（15 枚 `_test.go` 逐枚 `grep -c` 全 0）⇒ 名册内不存在并发用例，负载标签可以直接归因。

---

## 3. 上界的牙（Linux）：2s 在真例仍绿、在真饿仍响亮红；**变的是"要多饿才算饿"**

`settleWindowWaitBound = 2 * time.Second`（`window_wait_136_test.go:69`），推导建在 Windows 的
p99＝100/200ms 上（同文件 `:44-59`，引的是 poll 程 §1.2 那两行表）。本程**不改代码**，只问两句话：
**(a)** 真例（一次停顿、随后救援）在 Linux 还绿吗；**(b)** 真饿（读数永远不够）在 Linux 还响亮红吗。
再加一问：**(c)** Linux 的调度让"饿"这件事需要多大的物理扰动。

### 3.1 一枚 100ms 窗在 Linux 交付几枚读——这是 2s 那句算术的**原料换了**

靶子腿 `TestCheckSettleHalfTheReadsFailedReportsItsLoss` 判据是 `want=4`。
Linux 无注入时（`snap-instr` 里在 `awaitWindow` 每次尝试打一行 `AWAIT want= try= n=`，纯仓外仪器）：

```
AWAIT want=4 try=1 n=10    × 3 发（sweep-none 控制批）＋ × 2 发（建树时的单发 sanity）
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
```
⇒ **Linux 的 100ms 窗交付 10 枚读**，判据要 4 枚 ⇒ **余量 6 枚**。
对照盘上已有的 Windows 读数：poll 程 §4 用**同一枚夹具、同一支 60ms 第一读停顿**量到的是
`only 2 reads taken`（红），而本程在 Linux 上用同一支 60ms 停顿量到的是 **5 枚**（绿）。
⇒ **同一枚物理扰动，Linux 比 Windows 多活 3 枚读**＝Linux 的 tick 粒度（Go runtime timer／WSL2）
比 Windows 的 15.6ms 分辨率细，一个窗里塞得下的拍数多。**这是本程对"上界推导的原料"最重要的一条更正**：
2s 那枚数在 Linux 上不是"10 倍 p99"这一味在保护，而是**每窗 2.5 倍的读数余量**在保护（见 §3.4）。

### 3.2 灵敏度曲线（Linux·净窗·第一读停顿 `n=1`·每窗都重停，`-count=3`）

| 第一读停顿 | 该窗交付的读枚数 n | 尝试枚数 | 判定 | 腿时 |
|---|---|---|---|---|
| 0（无注入） | **10** | 1 | PASS | 0.10s |
| 30ms | 8 | 1 | PASS | 0.10s |
| 50ms | 6 | 1 | PASS | 0.10s |
| 55ms | 6 | 1 | PASS | 0.10s |
| **60ms** | **5** | 1 | **PASS** | 0.10s ← Windows 在这一档是 `only 2 reads taken`／**RED** |
| 65ms | 5 | 1 | PASS | 0.10s |
| **70ms** | **4** | 1 | PASS（**恰好压线**） | 0.10s |
| **80ms** | **3** | **20** | **FAIL（上界到期）** | 2.01–2.02s |
| **90ms** | **1** | **20** | **FAIL（上界到期）** | 2.02s |

逐档原文在 `D:\tmp\observe-ac15-linux-r1-s1\out\sweep-snap-instr-first-1-*.log`（9 枚日志，只建不删）。
⇒ **Linux 上"饿到判据不满足"的门槛在 70ms 与 80ms 之间**；Windows 报告里那一发 60ms 就饿到了 2 枚。
⇒ 曲线是**单峰无回弹**的：n 随停顿单调下降（10→8→6→5→4→3→1），没有任何一档出现"该窗时好时坏"——
  在本机这个 CPU 数与这份负载下，边界是**可复算的**，不是掷硬币。

### 3.3 确定性前后对（Linux 版 §4）：**同一枚物理事件，改前红、改后绿**

注入换成**瞬时**那一形（`mode=global n=1 ms=80`：整个进程只停第一枚读，下一枚窗干净）——
这正是归档那四枚目击的形状（"某些窗少收了几枚拍"），也是 poll 程 §4 那发的 Linux 等价物。

**改前的树（`snap-pre`＝`ee5a25e` 版 5 枚测试文件＋无 `window_wait_136_test.go`，`grep await` 现量 0）：**
```
sampler_settle_coverage_136_test.go:214: precondition broken: only 3 reads taken, half-and-half needs a window to lose in
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)      ← 第二发没有停顿，自然绿
FAIL	github.com/CarlosShao/wisp/internal/observe	0.211s
```
⇒ 站点 `:214`、句子 `only 3 reads taken, half-and-half needs a window to lose in`
**与归档四枚目击逐字同形**（poll 程 §1.1 引的红句原文；归档里 2 枚写 `only 2`、若干写 `only 3`）。
⇒ **"这枚红是 Windows 专属"这一支被本程否掉**：同一枚形状在 Linux 上造得出来（在 80ms 那一档）。

**改后的树（`snap-instr`，同一份补丁）：**
```
AWAIT want=4 try=1 n=3
AWAIT want=4 try=2 n=10
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.20s)
```
**改后的树（`snap-post`＝逐字节等于交付码的快照，无仪器）：**
```
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.20s)
--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
ok  	github.com/CarlosShao/wisp/internal/observe	0.313s
```
⇒ **(a) 真例在 Linux 仍绿**，且**等待路径真的被执行了**：**1 枚重开**（`try=1 n=3 → try=2 n=10`），
腿时从 0.10s 涨到 0.20s＝多烧一整枚窗＝重开的时长签名。这正面回答派单 (iii)：
**"p50==max＝一次没走"在 Linux 也不是宿命——本程用人造瞬时饿窗让它走了，并数出 1 枚重开/发**。
⚠ **与 Windows 那一发的可比性要说清**：poll 程 §4 的仪器写的是"把 `alternatingTree.ReadTree` 的
**前 N 枚读**各睡 60ms"，但那支 `arm_stall.py` 在**它自己的**临时件里（`D:\tmp\observe-ac15-poll-s1\`），
本程读不到它的源码 ⇒ **"前 N 枚"是"每枚新窗的前 N 枚"还是"整个进程的前 N 枚"这一味本程判不了**。
本程因此**两味都造了**：`mode=first`（每枚窗的第一读都停，＝"永久"那一支）与
`mode=global`（整个进程只停第一读，＝"瞬时"那一支）。两味的读数分别是 §3.2（永久⇒到期红）与
本节（瞬时⇒1 枚重开转绿）。Windows 那发 0.43s（＝3 枚饿窗＋1 枚满窗）落在哪一味，
**本程不能替它定**——它更像"永久那一味＋交付枚数在阈值附近抖"，因为瞬时形只需一枚重开就够
（Linux 实测 0.20s），凑不出 0.43s。这条归到 §6"本程没有测什么"里点名。


### 3.4 上界到期那一发：**逐字节复刻 Windows 的读数**

`mode=all ms=60`（每枚读都停 60ms＝永远救不回来）打在交付码的快照上：
```
sampler_settle_coverage_136_test.go:237: precondition broken: 16 windows opened inside the 2s monotonic
  bound (clock.go Timeout) all fell short of 4; counts seen: [2 2 2 2 2 2 2 2] (+8 more, all of them short of 4)
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (2.11s)
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (2.10s)
FAIL	github.com/CarlosShao/wisp/internal/observe	4.218s
```
poll 程 §4 末那发（Windows，同一支仪器）是：
```
precondition broken: 16 windows opened inside the 2s monotonic bound (clock.go Timeout) all fell short of 4;
  counts seen: [2 2 2 2 2 2 2 2] (+8 more, all of them short of 4)
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (2.10s)
```
⇒ **(b) 真饿在 Linux 仍然响亮红**，而且**尝试枚数（16）、每窗拿到的枚数（[2×8]+8 more）、报出的腿时（2.10–2.11s）
三味逐枚同值**。这一支不需要"按 OS 重定上界"——它是被 `deadline=+100ms` 钉住的**窗长**决定的，两 OS 同形。
⇒ 注意 §3.2 里 `ms=80/90` 那两档是 **20 枚尝试**（每窗 ~100ms 且只停第一读），而 `mode=all` 是 **16 枚**
（每窗被逐读停顿拖到 ~131ms）：尝试枚数不同**是因为窗长不同，不是因为预算不同**。两形都在同一枚 2s 上界里收场。

### 3.5 那么 2s 该不该改？**本程的答案是"不该"，但理由换了，且有一条要登记的余量**

- 支持"不改"的两味**现量**：①真饿那一发逐字节复刻（§3.4）；②真例那一发**只需 1 枚重开**、
  用量是预算的 **0.20s / 2s＝10%**（§3.3）。
- **换掉的原料**：poll 程给 2s 的说法是"最慢腿 p99 的 10 倍"。Linux 上这条更弱的同时，有一条更强的：
  每窗交付 **10 枚读** vs 判据 **4 枚**（§3.1）。要让一枚窗少到 3 枚，需要 **≥80ms** 的第一读停顿；
  Windows 在 60ms 就到了 2 枚。**同一次调度抖动在 Linux 上造成的损失更小** ⇒ 上界的"饿穿"门槛在 Linux 更高，
  不是更低。**没有一枚本程的量读支持把 2s 调小或调大。**
- **要登记的余量（本程给的算术，用的全是 Linux 现量）**：家族里 nominal 最长的腿
  `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` 是**两枚 200ms 窗＝400ms**（§3.6 的净窗家族表），
  它的预算只有 **2s/400ms ≈ 5 枚窗**；靶子腿是 **2s/100ms＝20 枚**。
  ⇒ "预算≈20 次中位尝试"这句在**家族里不均匀**——那枚 400ms 腿只有 5 次、两枚 200ms 腿只有 10 次。
  忙窗会不会把那枚 400ms 腿真真按到 5 枚窗以外，是 §4 的读数要说的事；
  **本程不在 §3 里替它预判**，§4 落地后在 §5 给"要不要按最长腿／按 OS 加余量"的答复。

### 3.6 家族 12 枚腿的 Linux 净窗时长表（`-count=10 -v`，本程建树时的 pilot 批，01:51:35→01:51:52）

判读规则与 poll 程同：**只认顶层判定行的 `(0.10s)` 字段**与 `--- FAIL`／`ok`／`FAIL`，
不认 `t.Logf` 那些同样带 `file:line:` 前缀的行。

| 腿 | n | p50 | p90 | p99 | max | 红 | 窗形状·判据 |
|---|---|---|---|---|---|---|---|
| `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=3 |
| `TestCheckSettleHalfTheReadsFailedReportsItsLoss` **靶子** | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=4 |
| `TestCheckSettleZeroFootprintDropsAreCountedToo` | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=3 |
| `TestCheckSettleFullyMeasuredWindowReportsNoLoss` | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=3 |
| `TestCheckSettleZeroTrustworthySamplesFailsClosed` | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=1 |
| `TestCheckSettleTrustworthyReadsAreRecorded` | 10 | 100ms | 100 | 100 | 100 | 0 | 1×100ms，want=1 |
| `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | 10 | 200ms | 200 | 200 | 200 | 0 | 1×200ms，want=1 |
| `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | 10 | 200ms | 200 | 200 | 200 | 0 | 1×200ms，want=2 |
| `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | 10 | **400ms** | 400 | 400 | 400 | 0 | **2×200ms**，各 want=1 |
| `TestSampleStateAllMetricsAndVerdicts` | 10 | 120ms | 120 | 120 | 120 | 0 | 1×120ms，want=5 |
| `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable` | 10 | 60ms | 60 | 60 | 60 | 0 | 1×60ms，want=3 |
| `TestSamplerGoroutineAccountingFollowsRegistry` | 10 | 30ms | 30 | 30 | 30 | 0 | 1×30ms，want=3 |

合计 **120 枚腿执行／12 枚腿／红 0／`precondition broken` 0**，日志 `out/pilot-family-q.log`（包时 16.228s）。
⇒ **Linux 的档位与 Windows（poll 程 §6.3：100/120/200/400/30/50/60ms）逐档同形**；
⇒ **`p50 == max` 在 Linux 净窗里同样成立**＝这 120 发**零重开**，poll 程那句"改后可能没走过等待路径"
在 Linux 的净窗里**没有被改善**——这恰是派单 (iii) 的形状，所以本程才另造 §3.3 那发人造饿窗去走那条路；
⇒ 13 处调用点／12 枚腿本程自量（`grep -rnE 'await(Settle|State)Reads\(t, ' internal/observe/*_test.go`
⇒ 13 行、去重 ⇒ 12 枚腿），站点行号 `coverage:143/237/321/374`、`gate:166/217/264/267`、
`sampler_test:93/335`、`zerosample:152`、`settle_zerosample:43/140`——**与 poll 程 §7.7 那张更正表逐枚相同**，
本程没有另开一枚行号读数去推翻它。

