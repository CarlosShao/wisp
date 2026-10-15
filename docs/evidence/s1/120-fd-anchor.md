# 120 r2 交件件 —— 封 mode 那一跳从"按名字查完再按名字改"换成"逐层钉住描述符＋同一对象复验"

**票**＝`.scratch/wisp/issues/120-posix-seal-is-check-then-act-154-of-4000-retries-hit-the-window.md`
（票面 AC 框一格没翻、文件名没改、判语一格没写）。
**本件不是裁决表**：写者≠判者，终裁归非实现者。本件只摆四样：改了什么形状、成对读数的分子分母、
反控怎么做的（含它没覆盖到哪一支）、每把尺的 `rc=`。凡"看起来更严"的地方一律写明代价与已裁依据。
**程名**＝`120-r2`（落地腿，写码）。

## §0 锚点自量（全部现跑）

| 号 | 是什么 | 怎么量的 |
|---|---|---|
| `58de465a4369c49011a7472670e91cb7b3c75b10` | 起手 `git rev-parse HEAD`（dev） | 现量；成对读数的"改前那棵树"就是 `git archive 58de465a` 导出的字节 |
| `137442cdecbb…c13e2313` | 改后 `internal/winsec/winsec_other.go` 的 sha256 | `sha256sum` 现量；快照里那份 clean 字节与仓内这份同值（§3 反控用它作还原基准） |
| `34cc6bd8ca0d…47f00ec0` | 改前（=58de465a）同一枚文件的 sha256 | `git show 58de465a:internal/winsec/winsec_other.go` 接 `sha256sum` 现量 |
| `f7327671ab97…db7d4741` | 反控变异体（`sameObjectPinned` 被去牙）的 sha256 | 变异后现量，见 `120-control-driver.log` |

改动面＝两枚文件，全在派单划的写面内：
`internal/winsec/winsec_other.go`（425 行，产码）、`internal/winsec/posix_seal_window_120_other_test.go`
（486 行，仪器）。`winsec.go`／`winsec_windows.go`／`internal/risk/**`／`cmd/wisp/**`／阈值／golden／
`allowlist.txt` 一律未碰。

## §1 落进去的形状（next=2）

产码侧三条腿，函数名即落点（行号现量自上述改后文件）：

| 函数 | 行 | 干什么 |
|---|---|---|
| `applyDescriptor`（变量） | :23 | 仍指 `applyDescriptorPOSIX`；AC#5 那枚注入缝（无 tag 的 `private_fail_test.go` 会换掉它）保持被咨询 |
| `applyDescriptorPOSIX` | :83 | `pinPlacement(path)` 拿到钉住的 fd → `chmodPinned(fd, …)`，中间不再按名字走第二遍 |
| `pinPlacement` | :118 | `os.OpenRoot(pinBase)` 打底，逐层 `Root.OpenRoot(comp)` **先钉**、再 `next.Stat(".")` 与 `parent.Lstat(comp)` **后验**；叶子用 `OpenFile(leaf, O_RDONLY\|syscall.O_NONBLOCK, 0)` 钉 |
| `sameObjectPinned` | :183 | 同一对象复验：先拒"名字上带链接"，再 `os.SameFile(pinned, now)` 不同即拒 |
| `placementRefusal` | :214 | 拒绝句同时带 `ErrUnresolvedPath` 与 `ErrNotSealable`（两个 `%w`），仪器据此把"查前拒"与"钉后拒"分开数 |
| `chmodPinned` | :224 | `f.Chmod` 打在**钉住的那枚 fd** 上，再 `f.Stat` 读回校验；沿用原错误文案 `mode is %v, want %v (%s stores no usable permission bits)` |
| `sealHandle` | :256 | 先走 `applyDescriptor`（注入缝），再对 `f` 本身做 `SameFile(f.Stat(), os.Lstat(f.Name()))` 复验，最后 `f.Chmod`；⛔ 复验之后再按名字解析一次 |

顺序是刻意的：**钉 → 复验 → 动**，不是"验 → 钉"。两个交错（钉后换、复验后换）都落在拒绝路上；
残余"换回另一枚真对象"那形由 `SameFile` 兜住。
注释里写死两件事：`O_NOFOLLOW`／`os.Root` 的链接规则⛔当防线（实测：目标留在根内的兄弟相对链接
会被跟、绝对目标与带 `..` 的目标才被拒；`Root.Chmod` 自述有竞态）；⛔ 最后一次复验之后按名字再解析。

## §2 成对读数（同一 fixture 家族、同一镜像 golang:1.27、同一台容器；缺一不成凭据）

**改前一发**＝`git archive 58de465a` 导到仓外快照 `base-head-newinstr`，只把本腿这枚测试文件拷进去，
`WISP_120_BASELINE=1`（读者档）。日志 `120-prefix-newinstrument-reader.log`。
**改后一发**＝仓内工作树，严格默认档（无 env），`go test -count=2 -v ./internal/winsec/`。
日志 `120-postfix-gate-count2-v.log`。分子＝`hit-foreign`（印章落在外人树上），分母＝`attempts`。

| 形 | 改前 hit-foreign/attempts | 改后第 1 发 | 改后第 2 发 |
|---|---|---|---|
| 1 叶子换位 `leaf-swap` | **7**/2000 | 0/2000 | 0/2000 |
| 2 祖先换绝对链接 `ancestor-swap-absolute-link` | **8**/2000 | 0/2000 | 0/2000 |
| 3 祖先换兄弟相对链接 `ancestor-swap-relative-sibling-link` | **11**/2000 | 0/2000 | 0/2000 |
| 合 | **26**/6000 | 0/6000 | 0/6000 |

改后每发的**活性**（⛔ 只数 hits，命中 0 与仪器死了是两种 0）：
第 1 发 `refused-at-check=1003/114/138`（其中 `refused-at-act=11/8/23`）、`sealed-own-object=864/128/181`；
第 2 发 `refused-at-check=993/133/164`（`refused-at-act=7/5/51`）、`sealed-own-object=860/109/191`。
`racer` 自报 flip 错误＝**0**（改前那枚仪器是 1，见 §4）。整包四数：`RUN=110 PASS=110 FAIL=0 SKIP=0`，
无 `RED`／`FLIP BROKEN`／`INSTRUMENT DEAD` 行。
判据是逐样本的：默认档下任一样本 `hit-foreign>0` 即 `t.Errorf`，不平均、不摊平。

## §3 反控（物理改写＋sha256 还原；⛔ `-overlay`）

`go test -overlay` 换的只是编译期输入，仪器是运行时读盘上产码的 AST/行为台件，用它变异会恒绿＝假读数。
这一轮按另一枚腿刚立的规矩做：**盘上改写**。改写发生在仓外快照 `base-fixed`，仓内工作树自始至终没被
变异过（`git status --porcelain -- internal/winsec/winsec_other.go` 只有 ` M`，即本腿那一枚改动）。
变异＝把 `sameObjectPinned` 里两枚判据各改成 `… && false`（只动这一枚函数，其余字节不变）。

```
sha256 restored=137442cd…c13e2313   sha256 mutated =f7327671ab97…db7d4741
sha256 after restore=137442cd…c13e2313   RESTORE MATCHES
repo sha256=137442cd…c13e2313  snapshot-clean sha256=137442cd…c13e2313
```

去掉牙之后读数（`120-control-gate-neutered.log`，`control rc=1`）：形 3
`total … hit-foreign=25 of-which-reported-nil=25`，四枚样本全红，红句逐字（含文案）：

```
AC#1 RED sample 1/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 8 time(s), 8 of them while returning nil (window open on THIS sample; every sample is judgement, no averaging)
AC#1 RED sample 2/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 4 time(s), 4 of them while returning nil (…)
AC#1 RED sample 3/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 9 time(s), 9 of them while returning nil (…)
AC#1 RED sample 4/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 4 time(s), 4 of them while returning nil (…)
--- FAIL: TestPOSIX120SealWindowAncestorRelativeLinkIsClosed (0.29s)
```

**反控没覆盖到的两支，明写**：
1. 形 1、形 2 在去牙后仍 `hit-foreign=0`（`--- PASS`）。原因不是它们假，而是这两形种的是**绝对目标**链接，
   `os.Root` 自己就拒（`openat …: path escapes from parent`），`refused-at-act=0` 那行红句正是为此而写：
   `AC#1 shape=leaf-swap: every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate`。
   形 3（兄弟相对链接，`sub -> other`，不含 `..`，Root 会跟）是唯一只有同一对象复验拦得住的形。
2. `sealHandle` 自己那半 `SameFile` 复验**没有具名单元台件**——变异只改了 `sameObjectPinned`，
   `sealHandle` 一支没被去牙。归 §7。

## §4 仪器自报缺陷两处（next=1）

| 缺陷 | r1 的自报 | 这一轮 | 读数变化 |
|---|---|---|---|
| 1 叶子赛手用"删了再建"腾名字，绝大多数样本撞 `openat … no such file or directory` | 形 1 `other=1725/2000` | 改成交替两名 `rename(2)` 原子换位（`.staged-a`/`.staged-b`），叶子上永远有一枚真对象 | `other` 由 1725 → **52**（同形同 attempts，见两份日志） |
| 2 祖先 flip 步不幂等，一次 flip 记 `recorded flip errors: 1`（`mkdir …: file exists`） | 红在 `120-prefix-reader-reading.log:31` | `clearByKind` 幂等：`Remove(blob)` → `Remove(dir)` → `Symlink` | flip 错误 **1 → 0**；默认档下 `flipErrs>0` 会 `t.Errorf`（`AC#1 RED instrument …a racer that cannot flip is not evidence about the window`） |

新增第 3 枚形（自加，理由写在件头）：见 §3 第 1 点——⛔ 加这枚，形 1/2 的 0 可以靠 stdlib 的拒答蒙过去。

## §5 两处代价（next=2(c)，按台账口径写成"已裁"，⛔ 写成待定）

1. **r 位收紧**＝`docs/reports/pending-and-issues.md` 的 **A990 裁决 2，已接受**。方向是拒，不是松
   （`applyDescriptorPOSIX` 从不放宽 mode）；这条腿 `!windows`，出货平台零动。代价现量见下表：钉 fd 要
   能读，于是 `0200`/`0000` 文件、`0300` 目录从"能修好"变成"拒绝"。
   **uid 0 会把这代价读成零**（容器默认 uid 0 时改前改后全 `SEALED`＝失真），所以刻意降权到 uid 1000 复测：

| 形（permface 台件） | uid 0 改前 | uid 0 改后 | uid 1000 改前 | uid 1000 改后 |
|---|---|---|---|---|
| `file-0200` | SEALED→600 | SEALED→600 | SEALED→600 | **REFUSED**，盘上仍 200 |
| `file-0000` | SEALED→600 | SEALED→600 | SEALED→600 | **REFUSED**，盘上仍 000 |
| `dir-0300` | SEALED→700 | SEALED→700 | SEALED→700 | **REFUSED**，盘上仍 300 |
| `dir-0500`（有 x 无 w 读得到） | SEALED | SEALED | SEALED | SEALED→700 |
| `handle-0200 PrivateFile` | SEALED→600 | SEALED→600 | SEALED→600 | **REFUSED** |
| `dir-0300-parent PrivateDirAll` | err=nil | err=nil | err=nil | **REFUSED**（`openat noexec-parent: permission denied`） |

   拒绝句形如 `winsec: cannot apply a private security descriptor: /tmp/permface…/file-0200: openat file-0200: permission denied`。
   整包在 uid 1000 下改前改后都是 `ok`（`120-uid1000-package-base-head.log` / `-base-fixed.log`）＝
   这套名册里没有"递一枚读不到的文件"的用例，收紧不会把今天的用例推红。
2. **`privateFile` 的开那一跳（先 `Create` 再 `Open` 再封）⛔ 本票没修**：票面 120 的射程是 placement 的
   check-then-act，开 leg 的形状另立 **票 322、未修**。产码注释里也是这么写的
   （`winsec_other.go:82`："…is booked separately as ticket 322 and is unfixed at this anchor"）。
   ⛔ 本件不得被读成"120 之后私钥文件全程原子"。

## §6 门禁逐把（`rc=` 现量；日志同目录）

| 尺 | 口径 | rc |
|---|---|---|
| 容器整包 | `docker run golang:1.27`：`go test -count=2 -v ./internal/winsec/` | **0**（四数 110/110/0/0） |
| AC#3 名册 | 票 113 自己的 `-run` 四族 POSIX 前缀（`TestAC1POSIX`／`TestAC2POSIX`／`TestAC3POSIX`／`TestAC4POSIX`），容器 | **0**（31/31/0/0，0 红 0 Skip） |
| 120 三枚具名 | `-run "TestPOSIX120"`，容器，默认严格档 | **0** |
| `gofmt -l` | 两枚文件 | **0**（空列表） |
| `gofumpt -l` | `/d/work/base/gopath/bin/gofumpt.exe` | **0**（空列表） |
| `go vet ./internal/winsec/` | host（windows，`winsec_other.go` 不带 tag 时不参与） | **0** |
| `GOOS=darwin GOARCH=arm64 go vet ./internal/winsec/` | host 交叉 | **0**（这枚才真读到本腿产码） |
| `sh scripts/d22scan.sh` | host，干净快照 | **0**，`clean - no D22 ban violations`；分母：`bans #1-5 internal/=239, cmd/=41, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, frontend/=116, internal/=563, cmd/=136`，`examined 280 production Go files`（⛔ 一枚分母没掉） |

台件卫生：快照一枚 `mktemp` 家一棵树；测量前先数 `.go` 证非空（`base-head` 1073、
`base-head-newinstr` 1072（差 `permface/main.go` 一枚）、`base-fixed` 1073；`internal/winsec/` 各 33 枚）。
无静默 `.out`；`GOTMPDIR` 全程 `C:\Users\swq\AppData\Local\wisp-test-gotmp\120r2gotmp`。

## §7 没做完的（逐枚具名，⛔ 摊平）

- `next=2 之 sealHandle 的同一对象复验`：⛔ 具名单元台件（§3 第 2 点）。要补一枚"fd 已开、名字随后被换走"
  的用例，且反控要单独去它的牙。
- `next=4 的射程外半边`：默认档只钉住 `internal/winsec` 在 Linux 容器里的读数。`internal/risk`、`cmd/wisp`
  等消费者在 Linux CI 上会不会因 §5 的 r 位收紧变色，本腿没量（写面之外）。
- `next=5（r1 转下来的）票面措写`：120 票面那句"154/4000"仍是改前口径，改措写要人批，本腿只追加 1 条 bullet。
- 形 3 的 `refused-at-act` 在改后读数里 23/51，形 1、形 2 也非零，但⛔ 没有任何一发保证"每形每发都 >0"——
  那一行只 `t.Logf` 不 `t.Errorf`。要不要把它升成硬判据，归判者。

## §8 复跑

```sh
cd 'D:/work/workspace/projects plans/Wisp'
export GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\120r2gotmp'
go test -count=2 -v ./internal/winsec/                       # 容器里同一条
WISP_120_BASELINE=1 go test -count=1 -v -run TestPOSIX120 ./internal/winsec/   # 读者档（改后仍可切回）
WISP_120_ATTEMPTS=2000 go test -count=1 -v -run TestPOSIX120 ./internal/winsec/  # 放大分母
sh .scratch/wisp/probes/120/r2/gate2.sh   # 快照构建 + AC#3 名册（后半 uid 那几发有已知 bug，见 gate3.sh）
sh .scratch/wisp/probes/120/r2/gate3.sh   # uid 0 / uid 1000 的 permface 四发 + 整包两发
sh .scratch/wisp/probes/120/r2/control2.sh  # 反控：去牙 → 红 → 还原并证 sha
```

盘上读数件（全部留底，⛔ 删）：`120-prefix-reader-reading.log`（r1 形状的复现）、
`120-prefix-newinstrument-reader.log`（成对读数前半）、`120-postfix-gate-count2-v.log`（后半）、
`120-ac3-roster-113-names.log`、`120-control-gate-neutered.log` + `120-control-driver.log`、
`120-{root,uid1000}-permface-base-{head,fixed}.txt`、`120-uid1000-package-base-{head,fixed}.log`、
`120-d22scan-postfix.log`、`120-host-rulers.log`、`exp-readout.log`/`exp2-readout.log`（`os.Root` 语义实测）。
仓外：`/d/tmp/probe120/r2/{base-head,base-head-newinstr,base-fixed,permface,exp,exp2}`，
`base-fixed/internal/winsec/winsec_other.go.bak-120r2-control` 留底未删。
`gate2.sh` 后半四发死于已知 bug（`MSYS_NO_PATHCONV=1: command not found`），失败件按"只建不删"留着，
由 `gate3.sh` 补跑；`control.sh` 首版死于 `cp` 到不存在的 `/work/.scratch`，同理留着。
