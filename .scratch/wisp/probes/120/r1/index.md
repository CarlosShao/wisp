# 票 120 r1 交件（agent=120-r1，2026-10-15 16:3x +08，deadline-pause 收尾）

- 锚点（起手现读）：`git rev-parse --short HEAD` = **bf25809a**；本次交件**未动一行生产码**（`internal/winsec/winsec_other.go` 保持原样，`winsec.go`/`winsec_windows.go`/`internal/risk/**` 未碰）。
- 写面：`internal/winsec/posix_seal_window_120_other_test.go`（本票用例，已入库，入库形态＝**reader / 基线读数形态**）。

## 0. 结论先行（逐格）

| 格 | 到哪一步 | 读数/件 |
| --- | --- | --- |
| **AC#1** 可重跑用例 | **仪器已入库**（reader 形态：默认只读不判；严格形态 `WISP_120_ASSERT_CLOSED=1`）＋**基线读数 > 0 已量到**（两形状都读到窗口） | 见 §1；件 `120-repo-reader-check.log` |
| **AC#2** 前置问题 | **已答，判据可执行**（容器实测表见 §2）：票面点名的**完全体做不到**（写面内）；写面内**可做到的替代形态已定** | 见 §2；探针源 3 枚 |
| **AC#3** 113 那批容器用例 | **只有基线读数**（全绿、0 Skip）；**未**在"修复后"的码面上验（本票没落修复） | 见 §3 |
| **AC#4** 是否须改边界话 | 本票不需要改 `winsec.go` 包文档语义（"边界话"）⇒ 该格无绑定登记；但**两条要编排者裁**（§4） | 见 §4 |
| **AC#5** 门禁四数 | **只有基线读数**（`-count=1 -v`，rc=0、RUN=54/PASS=32/FAIL=0/SKIP=0、`gofmt -l` 0 行、`go vet` rc=0）；`-count=2 -v` 四数、gofumpt、d22scan **未跑** | 见 §3 |

**一句话**：这一轮把 AC#2 的可行性判据量实了、把 AC#1 的仪器落成了库里的基线读数器（两处自身缺陷已具名），**修复本身与"修后必须 0"没做**。

## 1. AC#1 — 仪器与基线读数

件：`internal/winsec/posix_seal_window_120_other_test.go`（入库形态见文件头 `STATE OF THIS INSTRUMENT` 段）。

**形状**（票面 P4 的重试形状）：一个 racer goroutine 持续翻转"调用方点名的那个对象"（叶子 / 或某一级祖先目录），在"树内真对象"与"指向外来见证文件的链接"之间来回；主循环对同一拼写反复调 `SealFile`。每次尝试**从见证对象读结论**（外来文件 mode 变 0600 即命中），并另记"命中且返回 nil"（"动了外来树"与"报成功"是两件事）。

**本文件里已具名的两处缺陷**（不藏、不调宽，见文件头原文）：
1. **叶子 racer 换手不原子**：它 remove 再重建叶子 ⇒ 叶子在四个步骤里有两步"不存在"，密封调用多数尝试落在"no such file"空档（首跑 sample 1：`other=89/100`）。修法方向已写在文件头：rename(2) 原子替换 + 两个 staged 名字交替。
2. **祖先 racer 原先把失败经 `t.Fatal` 抬在 goroutine 里**：`t.Fatal` 不终止 racer、只把测试标红（首跑读数因此变成死仪器：`refused=0、sealed-own=100`×4）。入库形态改成**记录**失败（§ 读数里 `recorded flip errors: 1`，那一步不幂等）；翻转仍会到达链接态，所以读数**是活的**，但该步的修法（先清空目录再 rmdir 再 symlink）仍在 `next=`。

**基线读数**（容器 `golang:1.27`、`CGO_ENABLED=0`、uid 0；`internal/winsec` 未改动的码面 = HEAD bf25809a）：

```
go test -count=1 -v ./internal/winsec/       # reader 形态（默认）
rc=0   RUN=54  PASS=32  FAIL=0  SKIP=0
leaf-swap    : attempts=2000 hit-foreign=16 of-which-reported-nil=13 refused-at-check=110 sealed-own=119 other=1755
ancestor-swap: attempts=2000 hit-foreign=12 of-which-reported-nil=8  refused-at-check= 79 sealed-own= 81 other=1828
ancestor-racer recorded flip errors: 1
```

同日三次抓取的两形状命中率在 **0.2% - 0.8%（2000 发里 4 - 16 发）**，逐样本全报、不平均、不取最好（本目录的 log 是最后一次抓取，含每样本行）。对照：票 113 验收方自己的探针在修后码面读到 **154/4000 = 3.85%**；本轮读数低于它**是缺陷 1 的代价**（分子被"空档"吞掉），不是另一个判断。

**可重跑命令**（本机 Git Bash，仓外快照 + 本用例复制进去；本目录 log 即此命令的输出）：
```bash
mkdir -p /d/tmp/probe120/base && git archive bf25809a | tar -x -C /d/tmp/probe120/base
cp internal/winsec/posix_seal_window_120_other_test.go /d/tmp/probe120/base/internal/winsec/
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v /d/tmp/probe120/base:/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 \
  sh -c 'go test -count=1 -v ./internal/winsec/'
```
（严格形态：加 `-e WISP_120_ASSERT_CLOSED=1`，此时命中 > 0 与 liveness 失败都会判红 —— 那是"修后"该有的形状；`WISP_120_ATTEMPTS=<n>` 缩样本。）

## 2. AC#2 — 前置问题的**可执行判据**与结论

探针源（本目录，容器 `golang:1.27` 内 `go run` 实跑）：`ac2-probe-open-root.go`（os.Root 语义）、`ac2-probe-nofollow-and-links.go`（O_NOFOLLOW 是否经 Root 生效 / 各类链接是否被跟随）、`ac2-probe-permission-shapes.go`（权限形状，uid 0 与 uid 1000 各跑一遍）。

实测读数：

| 判据 | 读数 |
| --- | --- |
| `syscall.Openat` 的可编译平台 | `GOOS=linux go build` ok；**darwin / freebsd / openbsd / netbsd / dragonfly / solaris 全部 `undefined: syscall.Openat`** ⇒ stdlib 里"逐级 openat"只有 linux |
| `os.Root`（`go doc`，go1.27.1） | "Methods on Root **will follow symbolic links**, but symbolic links may not reference a location outside the root. Symbolic links must not be absolute." ＋ "On Unix, Root.Chmod/Chown/Chtimes **are vulnerable to a race condition**…" |
| 经 `os.Root` 开**叶子**链接（绝对目标） | 拒：`openat …: path escapes from parent` |
| 经 `os.Root` 开**叶子**链接（相对目标），**显式传 `O_NOFOLLOW`** | **OPENED** ⇒ `O_NOFOLLOW` 经 `Root.OpenFile` 失效（stdlib 自己把链接解了） |
| 经 `os.Root` 走**祖先**链接（相对，如 `../foreign`） | **OPENED（被跟随）** ⇒ Root 不做"任何链接都拒"这条规则 |
| `Root.OpenRoot` 钉住子目录 | fd 锚定 OK（绝对目标链接 → 拒；相对目标链接 → 跟随） |
| 真实文件/目录经 Root 开 + `f.Chmod` + `f.Stat` 回读 | ok（`-rw-------` / `drwx------`，回读同 fd） |
| FIFO + `O_NONBLOCK` 经 Root | ok（不阻塞） |
| 权限形状（**uid 1000**，非 root） | 文件 mode 0200 / 0000 走 `O_RDONLY` → `permission denied`；目录 0300 → `permission denied`；mode 0400 文件 ok。**uid 0 下全部 ok**（容器 root 无视 mode ⇒ 该格在容器里失真，这次是特意降权另跑取的） |

**结论（判据）**：
- **票面点名的完全体（逐级 `openat` + `O_NOFOLLOW` + 按已验证 fd 密封）在本仓写面内做不到**：stdlib 逐级 `openat` 只有 linux（darwin/BSD 全无），而 `os.Root` **不实现**本包底线的那条规则（它跟随"根内相对链接"，`O_NOFOLLOW` 经它失效、且它自陈 `Chmod` 有竞态）。要真关到"任何一级被换成链接都拒"，须**按平台拆文件**或**加 `x/sys/unix` 依赖边** —— 两者都在本腿写面外（按 AC#4 交回编排者裁）。
- **写面内能做到的替代形态（下一步照此落，全部 portable stdlib）**：
  1. **持 fd 的入口**（`sealHandle`，PrivateFile/PrivateFileExclusive 走它）：密封用 `f.Chmod`（fchmod 该 fd）+ `f.Stat` 回读 —— 模式落在**字节将要落到的那枚对象**上；前面加 `os.SameFile(f.Stat(), os.Lstat(f.Name()))` 的**同一对象闸门**（开的时候曾跟随链接、或名字在开之后被换走，两个事实就会不一致 ⇒ 在任何字节写出前拒绝）。`Lstat` 不跟随叶子，所以这一步不需要 walk。
  2. **按名入口**（`applyDescriptorPOSIX`，SealFile/SealDir/PrivateDirAll 走它）：逐级 `Root.OpenRoot(comp)` **钉住**（fd 锚定）＋`Root.Lstat(comp)` **复验**同一对象（`os.SameFile`），叶子同法钉+验，然后 `fchmod` **钉住的那枚 fd**。任何一级在窗口里被换成链接 ⇒ 钉住的对象与名字不是同一枚 ⇒ 拒；末次复验之后不再按名复解析，所以 act 无法被换走。这不依赖 `O_NOFOLLOW` 生效（已实测它经 Root 失效），只依赖"钉住"与"名字-对象同一性"两件事。
- **已量到的代价**（写进注释与票尾）：(a) **r 权限类收紧** —— 0w0/000 文件、缺 r 的目录从"（属主）可修"变"拒"（uid 1000 实测 EACCES；方向是拒、不是放宽）；(b) `privateFile`（**winsec.go**，写面外）的 `os.OpenFile` 在同一个窗口里仍会跟随叶子链接 → 对**外来文件 O_TRUNC 造成破坏**（新闸门拦得住私密**字节**，拦不住那次 open 的截断）⇒ §4 登记。

## 3. AC#3 / AC#5 — 本轮只有基线读数

容器命令与读数（code surface = 未改动的 HEAD 快照，`golang:1.27`、`CGO_ENABLED=0`、uid 0）：

```bash
go test -count=1 -v ./internal/winsec/     # 全包
rc=0   RUN=54  PASS=32  FAIL=0  SKIP=0     # log: 120-baseline-gate-winsec-count1-v.log（HEAD 快照，未含本票用例）
gofmt -l internal/winsec                   # 0 行
go vet ./internal/winsec/                  # rc=0
```
（本票用例入库后的同一条命令：rc=0、RUN=54/PASS=32/FAIL=0/SKIP=0，log: `120-repo-reader-check.log`。）

- **AC#3（113 那批）**：基线快照上 `internal/winsec` 全包 **0 FAIL、0 SKIP**，含 113 交付的 8 个顶层 + 4 个深度子用例与 119 的 4 枚（即 113 验收方那条 `-run 'TestAC1POSIX|TestAC2POSIX|TestAC3POSIX|TestAC4POSIX'` 口径的 16 枚）全部绿。**注意**：同一条 `-run` 口径**今天**会多匹配到票 108/125 同前缀的用例（今天 20 顶层 + 4 子 = 24），本票用例特意**不用** `TestACnPOSIX` 前缀，以免污染那条口径的计数。修复落地后必须重跑该口径并逐枚点名。
- **AC#5 门禁四数**：`-count=2 -v` 四数、`gofumpt -l`、`sh scripts/d22scan.sh`（含台账各 scope）**本轮未跑** —— 本票没有生产码改动可门禁；它们是 `next=` 里修复落地后那一轮的事。

## 4. AC#4 — 登记（交回编排者裁，未自行改任何禁区面）

1. 本票**不**需要改 `winsec.go` 的包文档语义（票面所称"那段边界话"）⇒ 该格无绑定登记；§2 的替代形态全在 `internal/winsec/winsec_other.go` 内。
2. **要编排者裁的两条**（都可写进票 120 的 next 轮）：
   - **`privateFile` 的开在 winsec.go（本腿写面外）**：叶子在"底线检查通过之后、open 之前"被换成链接时，`os.OpenFile(…O_CREATE|O_TRUNC)` 会跟随链接 → 外来文件被**截断**（私密字节被 §2 的 fd 闸门拦住，破坏仍发生）。要关它得动 `winsec.go` 的 open（如 `O_NOFOLLOW` + fd 复核），属另一格。
   - **r 权限类收紧代价**（§2 末）：采用 fd 锚定密封后，"属主但不可读"的对象从可修变拒；需机主/编排者确认可接受（方向为拒，且已实测非 root 读数）。

## 5. 未做（逐条 `next=`，按依赖排序）

```text
next=1 修仪器两处缺陷（都在同一文件内）：
      (a) 叶子 racer 换 rename(2) 原子换手（两枚 staged 名字交替，叶子永不缺席）；
      (b) 祖先 racer 的翻转步改幂等（先 Remove(blob) 再 Remove(dir) 再 Symlink），并去掉那次 recorded flip error。
next=2 落 AC#2 结论里的替代形态（仅 internal/winsec/winsec_other.go）：
      (a) sealHandle: SameFile(f.Stat(), os.Lstat(f.Name())) 闸门 + f.Chmod + f.Stat 回读（AC#5 注入缝 applyDescriptor 仍要被走到，见 private_fail_test.go）；
      (b) applyDescriptorPOSIX: Root.OpenRoot/Root.Lstat 逐级 pin+verify，叶子同法，fchmod 钉住的 fd；
      (c) 注释：把代价（r 权限类、privateFile 的开）按本仓"两条成本"体例写清，不新造承诺。
next=3 门禁（容器，一轮）：-count=2 -v 四数逐条点名；run 'TestAC1POSIX|TestAC2POSIX|TestAC3POSIX|TestAC4POSIX' 逐枚点名（113 那 16 枚，0 红 0 Skip）；gofmt -l + gofumpt -l；go vet；GOOS=darwin go vet ./internal/winsec/；按 113 验收方的口径跑等价 uid-1000 读数；sh scripts/d22scan.sh 快照 rc=0、台账各 scope 不降。
next=4 仪器转严格：把 WISP_120_ASSERT_CLOSED=1 的断言（每样本 hits==0 + liveness）提升为默认，并在证据里记"修前 reader 读数 > 0（本件）/ 修后 0"两行。
next=5 交回编排者：§4 的两条裁 + 若采用 next=2(b) 的替代形态，票面"AC#2 只承诺要查、没承诺原子"那段口径需按新事实重写（属改票面口径，须人工批准）。
```

## 6. 盘上件一览（本目录，7 件）

| 件 | 是什么 |
| --- | --- |
| `index.md` | 本文件 |
| `120-baseline-gate-winsec-count1-v.log` | HEAD 快照（未含本票用例）全包 `-count=1 -v` 原始日志：rc=0、54/32/0/0 |
| `120-repo-reader-check.log` | 本票用例入库后同一命令的原始日志（reader 形态）：rc=0、54/32/0/0 + 两形状逐样本读数 + `recorded flip errors: 1` |
| `ac2-probe-open-root.go` | AC#2 探针：`os.Root` 开叶子/中间链接、`f.Chmod`+`f.Stat` 回读、绝对名与 FIFO（实跑过） |
| `ac2-probe-nofollow-and-links.go` | AC#2 探针：`O_NOFOLLOW` 经 Root 是否生效、绝对/相对链接是否被跟随、`Root.OpenRoot` 行为（实跑过） |
| `ac2-probe-permission-shapes.go` | AC#2 探针：权限形状（uid 0 与 uid 1000 各跑一遍）（实跑过） |
| `../../../../internal/winsec/posix_seal_window_120_other_test.go` | 本票 AC#1 仪器（入库形态＝reader，缺陷已具名） |

探针源在 **仓外目录 `/d/tmp/probe120/exp{,2,3,4}/main.go`** 被我实跑（本目录的 3 枚是它们的副本；`-count=1 -v` 的原始日志另有 `/d/tmp/probe120/out/base-winsec-v.log` 一份）。
