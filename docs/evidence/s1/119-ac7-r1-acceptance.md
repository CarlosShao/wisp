# 119 AC#7 —— 裁决表 r1（非实现者）

**被验格子**＝票 `.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md`
的 **AC#7 一格**（裁决定格时刻盘上 `:72-91`；本程锚点 `2956897` 那版票面里是 `:61-80`，
两者**文字逐字相同**，证明在 §0.3）。
**出表程**＝`acceptor-ticket119-ac7-r1`（验收方，≠ 实现方 `worker-ticket119-ac7`）。
**本格总判**：**成立（无附条件）**——理由与承重证据见 §5，判据①②③④逐格见 §1-§4。
**权限边界**：本表**不翻 AC#7 的勾、不改票名、不写票面 119 那枚文件**；勾由编排者按本表打。

---

## §0 身份、锚点、哈希与仪器（先自核身份，派单的话一律当"未验证断言"）

### 0.1 锚点自量（`git rev-parse` 现量，不抄派单）

| 时刻 | 量到什么 | 命令 |
|---|---|---|
| 17:21:36（开工） | `git rev-parse HEAD` = `2956897b1a9135b0427821ddf13cf690cc2da4d8` | `date && git rev-parse HEAD` |
| 17:24:43 | HEAD 已漂到 `c2fa2e9`（128 程落了 `cmd/wisp/dataroot_128_test.go`，**未碰我的 scope**） | `git log --oneline -3` |
| 17:45:08 | HEAD = `79cfa1d`（五枚邻居 commit：`c23d825 ef26704 8dac13b 744d79e 79cfa1d`） | `git log --oneline c2fa2e9..HEAD` |

⇒ **裁决树钉死为 `git archive 2956897`**（开工锚），并现量证明被测对象在整程里没动过：
`git diff --name-only c94927d..HEAD -- internal/winsec/` ⇒ **空输出**（17:45 复量）。
所以"两枚码 commit ＋ 我重量"这三件事落在同一枚文件版本上。

### 0.2 三枚文件的开工哈希（盘上 vs HEAD blob，逐枚）

```
$ git rev-parse HEAD                  2956897b1a9135b0427821ddf13cf690cc2da4d8
$ md5sum internal/winsec/dataroot_symlink_119_other_test.go
  4995e4f51104b338e3a33c1c7d0d25a9      盘上（382 行 / 20667 字节）  blob d798bf49cf89f8725a0252e378e89811e6388e30
  git cat-file -p HEAD:… | md5sum     4995e4f51104b338e3a33c1c7d0d25a9   ← 一致
$ md5sum docs/evidence/s1/119-ac7-impl.md
  d26726414ed26de966bc39dca251fc10      盘上（469 行 / 37043 字节）  blob 99247131b014aa00b92db2b6db71d176c7be8af9
  git cat-file -p HEAD:… | md5sum     d26726414ed26de966bc39dca251fc10   ← 一致
$ md5sum .scratch/wisp/issues/119-posix-link-…-data-roots.md
  e3cbb5efe104c2cba25e1d42005425e4      盘上（499 行 / 50507 字节）  blob 40dae1b229d743b2489b344cfe20894a0ffc6cc5
  git cat-file -p HEAD:… | md5sum     e3cbb5efe104c2cba25e1d42005425e4   ← 一致（17:21）
```

### 0.3 树被动过吗？动过，但**动的不是被测对象**

17:21 三枚全一致；**17:24 再量时票面 119 那枚已漂**：盘上 md5 `d56a6798…`／521 行，
与当时 HEAD blob（481 行）差 40 行。17:45 复量：漂移已被编排者自己提交（`79cfa1d` 之前那批），
盘上 == HEAD == 521 行 `d56a6798…`。漂的原因就是派单里那句"我要往里补回执"
（往 AC#1/AC#2 那两格下面插了复判回执块，并把 `AC#1／AC#2` 的勾打上）。

⇒ 这一枚必须钉死，否则"判据本体"会被读成动过的东西：

```
$ git cat-file -p HEAD:票面 | sed -n '61,80p' > /tmp/ac7_head.txt      # 开工锚 2956897 那一版
$ sed -n '72,91p' 票面 > /tmp/ac7_disk.txt                              # 现盘那一版
$ diff /tmp/ac7_head.txt /tmp/ac7_disk.txt ; echo rc=$?
（空输出）rc=0
```

即 **AC#7 那一格 20 行逐字未变**，只是行号从 `:61-80` 漂到 `:72-91`。本表按那一格的**文字**裁，不按行号。

### 0.4 两枚码 commit 的账（现量，与派单一致）

| commit | 内容 | numstat | 落点文件 |
|---|---|---|---|
| `ddc1583` | 两枚未解析根用例接记名断言 ＋ 同时换已解析基根 | 60 增 / 5 删 | `internal/winsec/dataroot_symlink_119_other_test.go`（仅此一枚） |
| `c94927d` | 两句红句里点名本用例自己种的那枚链接 | 2 增 / 2 删 | 同一枚 |
| 合计 | `git diff --numstat a9c8b6e..HEAD -- <file>` | **62 增 / 7 删** | 派单那句**复算成立** |

改前基线 `a9c8b6e6dce4ad63a6fc04b15385651a1443ea65`（`git rev-parse` 现量，非 `is-ancestor` 猜）：
`git cat-file -p a9c8b6e:<file> | md5sum` = `cec75e819aef89584cde7796878ef78a`，
与我 `git archive a9c8b6e` 快照里那枚**逐字节一致**（同一发 `md5sum` 现量）。

### 0.5 行号这一枚：**派单的前提不成立，实现方的行号是对的**

派单写："实现方报告里写的是 `:241`／`:361`，与盘上差 2 行 ⇒ 它引的是改前号／中间态号，还是它数错了"。
盘上现量＋读它原文的结论是第三种：**那两个号标的本来就是哨兵那一支，不是记名那一支**，一字不差。

```
$ grep -n "errors.Is(err, winsec.ErrUnresolvedPath)\|refusalCreditsLink137(err, shape.link)\|base := cleanSpelling119" <file>
225:	base := cleanSpelling119(t, t.TempDir())              ← AC#1 母项换根
241:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)   ← 哨兵支（AC#1）
243:	} else if !refusalCreditsLink137(err, shape.link)     ← 记名支（AC#1）
291:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)   ← :251 那一处（AC#3，不并入本格）
337:	base := cleanSpelling119(t, t.TempDir())              ← AC#2 leg 1 换根
361:	} else if !errors.Is(err, winsec.ErrUnresolvedPath)   ← 哨兵支（AC#2）
363:	} else if !refusalCreditsLink137(err, shape.link)     ← 记名支（AC#2）
```

`119-ac7-impl.md` §1 那段引文（同一发 `grep -n`）里，`:241`／`:361` 后面写的是"（原 `:203`／`:308` 那一处，条件未动）"，
`:243`／`:363` 后面写的才是"（新：记名断言）" ⇒ **实现方七个号全部与盘上一致**，
差 2 行是派单把"记名断言"安到了哨兵支的号上。本表把这一处记为**派单笔误**，不是实现方缺陷。
另一处同理：impl §0 写票面 AC#7 在 `:61-80` —— 那是 `b4e692e` 那版票面的真号（现量 `grep -n` 命中 61/80），
今天漂到 `:72-91` 是编排者自己之后插的 11 行，不是实现方数错。

### 0.6 仪器（违反任一条＝读数作废，所以我把它写全）

- **只在 Linux 容器真跑**：`docker version` ⇒ `29.6.2`；`docker images` 现量 `golang:1.27 = 3680233e3204`
  （**未 `docker pull`**）；容器内 `go version go1.27.1 linux/amd64`、`id -u`=0。宿主是 Windows，
  `//go:build !windows` 的用例在本地**没有分母**，所以本程没有一枚"宿主颜色"入账。
- **挂载一律 Windows 风格**：`-v "D:\tmp\wisp119-ac7-r1acc\snaps\<名>:/wisp"`
  （外加 `:ro` 的宿主 `GOMODCACHE` 与命名卷 `wispacc119-gobuild`，`GOPROXY=off` ⇒ 全程零网络）。
  每发容器第一句 `ls -l /wisp/go.mod` 进 `head.txt`，读数固定
  `-rwxrwxrwx 1 root root 883 … /wisp/go.mod` ＋ `gomod_bytes=883` ＋
  `f6ef661732b1851e5c3db348113cb605 /wisp/go.mod`；另有硬门 `[ ! -s /wisp/go.mod ] ⇒ exit 97`（不取颜色）。
- **台件全在仓外**、全部 `git archive <sha> | tar -x`，共 16 棵快照（清单 §8）；**禁读脏工作树冒充被验版本**。
- **每发 `-count=1 -v ./internal/winsec/`**（整包、不是 `-run` 子集，这样名册差集顺带证明"本格没让别人变色"）；
  四数只从 `-v` 日志取；`=== RUN` 名册／全量 colour 名册／SKIP 名册**每发都落盘**；`panic|fatal error` 每发计数。
- **形状每发容器里现种并硬断言**（种序：先 `mkdir -p /acwpriv/wacc119tmp` 再 `ln -s /acwpriv /acwlink`）：
  软链形断 `[ -L /acwlink ]` ＋ `readlink -f /acwlink/wacc119tmp == /acwpriv/wacc119tmp`，`TMPDIR=/acwlink/...`；
  普通形硬断 `/acwlink` **根本不许存在** ＋ `ls -ld /acwplain` 读到 `drwxr-xr-x`。断不过 ⇒ `exit 97`，不取颜色。
  ⚠ 我的链接名与前两程（`/ac7link`）刻意不同，读数是**另造形状重跑**，不是复用它的日志。
- **两枚取数前的门**：`go build ./internal/winsec/ ./internal/proc/` rc=0 **且** `go vet ./internal/winsec/` rc=0；
  任一非 0 ⇒ `GATE95` 直接 exit，**不发颜色**（这一枚是本程加的第二道门，来历见 §0.7）。
- **绝不把 FAIL 改成 SKIP**、不放宽任何断言、不碰阈值／golden；本程**未取任何时序／内存读数**。
- ⚠ **格式门这把尺的版本**：盘上现量 `"$(go env GOPATH)/bin/gofumpt.exe" --version` ⇒
  **`v0.12.0 (go1.27.1)`**（路径 `D:\work\base\gopath\bin\gofumpt.exe`，mtime 09-23 22:23）。
  票面 `:58`（AC#6 那一格）与本程派单写的"v0.7.0"**是过期事实**；本表所有格式读数按 **v0.12.0** 这把尺算，
  并与实现方 §6 现量的同一枚版本对得上（⇒ 两程是同一把尺）。

### 0.7 本程仪器自拒两枚（登记，不抹；这也是"改前那一发凭什么可信"的一部分）

1. **MUT-D2 第一版少了个逗号** ⇒ `winsec_other.go:158:170: syntax error: unexpected newline in argument list` ⇒
   `b09-head-mutd2-plain`／`b10-head-mutd2-link` 两发 **`GATE95` 直接 exit，零枚颜色入账**。
   修好的 `head-mutd2f` 才是 `d01`／`d02`。坏快照 `head-mutd2`／`rB1-mutd2`／`rB2-mutd2` **留着不删**（§6）。
2. **回退 B 第一版补丁留了一枚游离的 `}`**（我把 `} else if 记名…{` 换成 `}` 又删掉它的消息行，
   于是链尾原本那枚 `}` 变成多余）⇒ 测试件编译不过 ⇒ `c06`–`c12` 七发 **`RUN=0`**、`PKG_RC=2`、`[setup failed]`，
   归零入账。暴露它的不是 `go build`——**`go build` 根本不编译 `_test.go`**（`BUILD_RC=0` 而 `VET_RC=1`）——
   是**空名册**与 vet 的 rc。修法是 `del` 掉分支头与消息两行、只留注释（`REVERSION-WACC119-B2`），
   重建 `rB1-mutd-f2`／`rB2-mutd-f2`／`rB1-mutd2-f2`／`rB2-mutd2-f2`，读数 `e01`–`e08`。
   ⚠ 顺带把 vet 提成硬门（上一条）：已入账的 22 发**全部 `VET_RC=0`**（逐发 `grep` 复量在 §6 的 `gates.txt`），
   所以这道后加的门不改变任何一枚已入账读数，只让下一次这类自伤在取颜色之前就停住。

**入账读数一共 23 发**（`b01`-`b08`、`d01`-`d02`、`c01`-`c05`、`e01`-`e08`）＋ 两枚冒烟（`t0` 被 `GATE95` 停住、
`t1` 只证明仪器能跑）＋ 三发门禁容器（`gate-container2`、`gate-container-base`）。
作废的 `b09`/`b10`/`c06`–`c12` 日志与快照**全留**，标签不重用。

### 0.8 本节落盘的那枚 commit（原样输出，先测后写再提交）

```
$ git log --oneline -1
27a6f90 evidence(119 AC#7 r1 §0): 身份与仪器——裁决树钉死 git archive 2956897 …
$ git show --name-only --format="%H" HEAD | tail -3
27a6f90…（本节 §0 这一枚）
docs/evidence/s1/119-ac7-r1-acceptance.md
```

（提交前 `git diff --cached --name-only` 只有一行：`docs/evidence/s1/119-ac7-r1-acceptance.md`。）

---

## §1 判据①「两味药必须一起下」——**成立，且两味都承重**（不是并列装饰）

### 1.1 落点（我自己量的，命令原文见 §0.5）

| 药 | 落在哪两行 | 用的是谁 | 是否改前就有 |
|---|---|---|---|
| 已解析基根 | `:225`（AC#1 母项）、`:337`（AC#2 leg 1） | 本文件既有的 `cleanSpelling119`（＝`filepath.EvalSymlinks`） | **是**：`git cat-file -p a9c8b6e:<file> \| grep -n "func cleanSpelling119"` ⇒ `112` |
| 记名断言 | `:243`（AC#1）、`:363`（AC#2），各带一句点名本用例种的链接的红句 `:244`／`:364` | 同包既有 helper `refusalCreditsLink137` | **是**：`internal/winsec/placement_symlink_113_other_test.go:160`，最后动它的是票 137 的 `9c0f546`，而 `git merge-base --is-ancestor 9c0f546 a9c8b6e` ⇒ 真 |

**零新增依赖**这一枚我另走了一条不引它自述的尺：

```
$ diff <(git cat-file -p a9c8b6e:<file> | sed -n '/^import (/,/^)/p') <(sed -n '/^import (/,/^)/p' <file>)
（空输出）⇒ IMPORTS_IDENTICAL
$ git diff --numstat a9c8b6e..HEAD -- internal/winsec/placement_symlink_113_other_test.go
（空输出）⇒ helper 那枚文件一字未动；两版 md5 都是 25660e86ad055e68304839e227713059
$ git diff --name-status a9c8b6e..HEAD -- internal/winsec/
M	internal/winsec/dataroot_symlink_119_other_test.go     ← 整个 scope 只有这一枚
```

**换根那味的方向**也对：`cleanSpelling119` 问的是文件系统（`EvalSymlinks`），**不是** `proc.SealableRoot`
⇒ 不违反票 119 Rules 里 `R-119-9`"不许拿被测函数算 fixture"（那条正是本票返修时立的）。

### 1.2 承重证据 A：**只撤换根、留记名** ⇒ 软链形恒红（＝137 的 `r2` 处境，被这发复算并排除）

台件＝`git archive 2956897` 之上只把 `:225`（或 `:337`）换回 `base := t.TempDir()`，生产码未变异。

| 发 | 撤了哪一枚 | 形 | 四数（RUN／顶P/F/S／子P/F/S／panic） | 名册差集 vs 同形未改（`b06`／`b05`） | 变色的名字 |
|---|---|---|---|---|---|
| `c01` | AC#1 的换根 | **软链** | 45／26/**1**/3／15/0/0／0 | 2 行 | **只有** `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` PASS→FAIL |
| `c02` | AC#1 的换根 | 普通 | 52／30/0/0／22/0/0／0 | **0 行** | 无 |
| `c03` | AC#2 的换根 | **软链** | 45／26/**1**/3／15/0/0／0 | 2 行 | **只有** `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` PASS→FAIL |
| `c04` | AC#2 的换根 | 普通 | 52／30/0/0／22/0/0／0 | **0 行** | 无 |
| `c05` | 两枚都撤 | **软链** | 45／25/**2**/3／15/0/0／0 | 4 行 | 两枚母项一起红（＝回到 137 `r2` 的形状） |

**"撤掉这一处，哪一枚用例变了颜色？"答得出**：撤 AC#1 的换根只有 AC#1 红，撤 AC#2 的只有 AC#2 红，
两枚都撤两枚都红，普通形一枚都不红 ⇒ **两处换根都不是装饰**。
红句原文（各一枚，未截断）：

```
c01: dataroot_symlink_119_other_test.go:234: premise broke: the resolved base
     /acwlink/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1192121665142/001
     is itself refused by the floor (winsec: refusing to seal … premise …)
c03: dataroot_symlink_119_other_test.go:364: AC#2 RED: the refusal of the declared root
     "/acwlink/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1193044188816/001/injlink119/harness/picked"
     named ErrUnresolvedPath but did not credit the link this case planted at "…/injlink119" …
```

⚠ 这里有一条**实现方没写、我也不替它写**的机制差异：AC#1 那枚的恒红停在它新加的前提门
（`:233-235` 的 `t.Fatalf`，红在 `:234`），所以记名那一支根本没机会跑；AC#2 那枚没有独立前提门
（它的前提门是 leg 2），所以红直接落在记名那一支（`:364`）。两枚都恒红，**红因不同一处**。
前提门是 `ddc1583` 加的第三样东西，判据①没要求它；它的作用是让"解析过的基根被宿主链接污染"这件事
**响得比记名断言更早**，方向是收紧不是放宽 ⇒ 不判它违规，只把机制写清，免得下一位以为两枚都是记名支响的。

### 1.3 承重证据 B：**只撤记名、留已解析根** ⇒ 派单的预期不成立，我照实推翻

派单写"预期 MUT-D 打不红（零区分力回来了）"。**实测：打红。** 撤掉记名那一支之后，MUT-D 在两形
仍然把两枚母项打红，名册差集为 **0 行**：

| 发 | 撤了哪一枚 | 生产码 | 形 | 四数 | 名册差集 vs `b08`／`b07` | 变色 |
|---|---|---|---|---|---|---|
| `e01` | AC#1 的记名支 | MUT-D | 软链 | 45／15/12/3／11/4/0／0 | **0 行** | 无（两枚照旧 FAIL） |
| `e02` | AC#1 的记名支 | MUT-D | 普通 | 52／17/13/0／17/5/0／0 | **0 行** | 无 |
| `e03` | AC#2 的记名支 | MUT-D | 软链 | 45／15/12/3／11/4/0／0 | **0 行** | 无 |
| `e04` | AC#2 的记名支 | MUT-D | 普通 | 52／17/13/0／17/5/0／0 | **0 行** | 无 |

⇒ **MUT-D 的区分力是"换根"那一味给的**（红从哨兵支响，机制见 §3.4），记名那一支在这一发上本来就不响。
记名那一支的射程在**另一发变异**上，一撤就露：

| 发 | 撤了哪一枚 | 生产码 | 形 | 四数 | 名册差集 vs `d02`／`d01` | 变色 |
|---|---|---|---|---|---|---|
| `e05` | AC#1 的记名支 | MUT-D2 | 软链 | 45／19/**8**/3／11/4/0／0 | 2 行 | `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` **FAIL→PASS** |
| `e06` | AC#1 的记名支 | MUT-D2 | 普通 | 52／22/**8**/0／18/4/0／0 | 2 行 | 同上 |
| `e07` | AC#2 的记名支 | MUT-D2 | 软链 | 45／19/**8**/3／11/4/0／0 | 2 行 | `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` **FAIL→PASS** |
| `e08` | AC#2 的记名支 | MUT-D2 | 普通 | 52／22/**8**/0／18/4/0／0 | 2 行 | 同上 |

⇒ **把记名那一支摘掉，MUT-D2 就完全打不红那一枚**（另一枚照旧红＝逐枚归属，同时把判据③的"单点回退 C"结掉）。
这一支**承重**。

### 1.4 本格要害那一问的裁定

> 问：**"两发变异各自能响一支"是否等于"这味药承重"**，还是必须"同一发变异里记名支先响"？

**裁：前者。要求"同一发里记名支先响"是错的判据，我拒绝按它退这一格。** 理由三条，都是量出来的：

1. MUT-D 把两条走查截短之后**底线根本不拒**（`b08` 里 `AC#1 PrivateDirAll(...) -> winsec: … varlink119 is not a directory`，
   不是 `ErrUnresolvedPath`），而记名那一支的代码形状是 `} else if !refusalCreditsLink137(...)`，
   挂在 `errors.Is(err, ErrUnresolvedPath)` 之后——**它的触发前提在这发里被变异本身抹掉了**。
   要求它"先响"＝要求一支在它的前提为假时仍然响，那只能靠把 `else if` 拆成两条独立 `if` 来造，
   那是**措辞工程**，不是判据。
2. 承重的**可操作定义**是"拿掉这一味，存不存在一发变异打得它永不响"。§1.3 的 `e05`-`e08` 答"存在"
   （MUT-D2 摘不得），§1.2 的 `c01`-`c05` 答"换根摘不得"。**两味各自被一发单点回退打回不响**＝两味都承重。
3. 反向也要有一句：记名那一支**不是恒响的装饰**——未变异的 `b05`/`b06` 里它是绿的（前提真、结论真），
   MUT-D2 的 `d01`/`d02` 里它单独响。恒绿与恒红都排除了。

⇒ 实现方 §9.1 那句自述（"MUT-D 的红是从哨兵支响的"）**属实**，我独立复算到同一机制；
它另下一发 MUT-D2 证明记名支有牙，这一枚是对症的，不是搪塞。

### 1.5 判语

判据① **成立（独立复现）**。两味药一起下、都落在**该落的两枚母项**上、方向都不放宽，
且**两味各被一发单点回退证明承重**；判据里点名的"别复活 137 的 `r2`"由 `c01`-`c05` 五发正面排除。

### 1.6 本节落盘的 commit（原样输出）

```
$ git log --oneline -1
62088ce evidence(119 AC#7 r1 §1): 判据①成立且两味都承重——落点 :225/:337 与 :243/:363 …
$ git show --name-only --format="%H" HEAD | tail -2
62088ce…
docs/evidence/s1/119-ac7-r1-acceptance.md
```

⚠ 归属只认现量，不认本节这句话：`git log -L '/^## §1/,/^## §2/:docs/evidence/s1/119-ac7-r1-acceptance.md' --oneline`。
这一行 commit 记录本身是**随 §2 那枚 commit 入库**的（渐进写的必然错位，写清就不含糊）。

---

## §2 判据②「改前零区分力那一发要复现」——**成立（我自己复现的，不是引票面那句话）**

台件＝`git archive a9c8b6e`（改前那棵纯净树，快照里那枚 `_test.go` md5 = `cec75e819aef89584cde7796878ef78a`
＝ `git cat-file -p a9c8b6e:<file> | md5sum`）＋**我自己那一发 MUT-D**（不是复用它的脚本）。
`CREDIT_BRANCHES=0`、`RESOLVED_BASES=0`、`SENTINEL_BRANCHES=3` 每发容器里现量（`head.txt`）＝
**改前那两枚只有哨兵判据、没有记名判据**。

### 2.1 我的 MUT-D 与它的 MUT-D 是同一形，两处落点（先证落地，才读颜色）

```
$ grep -n "MUTATION-WACC119-D" snaps/base-mutd/internal/winsec/winsec_other.go snaps/base-mutd/internal/winsec/winsec.go
winsec_other.go:156:	sealPrefixes := pathPieces(path) // MUTATION-WACC119-D: walk cut short, first 3 prefixes only
winsec_other.go:157:	if len(sealPrefixes) > 3 { // MUTATION-WACC119-D
winsec_other.go:158:		sealPrefixes = sealPrefixes[:3] // MUTATION-WACC119-D
winsec_other.go:159:	} // MUTATION-WACC119-D
winsec_other.go:160:	for _, prefix := range sealPrefixes { // MUTATION-WACC119-D
winsec.go:284:	if len(prefixes) > 3 { // MUTATION-WACC119-D: same cut on the unlink route
winsec.go:285:		prefixes = prefixes[:3] // MUTATION-WACC119-D
winsec.go:286:	} // MUTATION-WACC119-D
```

第二证＝同一发容器里 `go build ./internal/winsec/ ./internal/proc/` ⇒ `BUILD_RC=0` 且
`go vet ./internal/winsec/` ⇒ `VET_RC=0`（`gates.txt`，非 0 就 `GATE95` 停住不取颜色）；
第三证＝每发 `head.txt` 里的 `MUT_other=5 / MUT_winsec=3` 树身份门
（未变异树恒为 `0/0`；`WINSEC_PROD_MD5_other` 从 `b5056918…` 变到 `2665ec1e…` ＝ 生产码真的被改过，
不是"我以为改了"）。⚠ 与我前两枚饵不同的一枚：我的截短**带长度守卫**（`if len > 3`），
它的 `sealPrefixes[:3]` 没守卫；被测拼写永远 ≥3 段 ⇒ 两版在这一格等价，我用守卫版是为了不让饵自己 panic。

### 2.2 改前四发（同一把尺、同两种形状）

| 发 | 形 | 生产码 | RUN | 顶 P/F/S | 子 P/F/S | panic | 包级 rc |
|---|---|---|---|---|---|---|---|
| `b01-base-plain` | 普通 | 未变异 | 52 | 30/0/0 | 22/0/0 | 0 | 0 |
| `b02-base-link` | 软链 | 未变异 | 45 | 27/0/3 | 15/0/0 | 0 | 0 |
| `b03-base-mutd-plain` | 普通 | MUT-D | 52 | 17/13/0 | 17/5/0 | 0 | 1 |
| `b04-base-mutd-link` | **软链** | MUT-D | 45 | 17/**10**/3 | 11/4/0 | 0 | 1 |

**那一发本身**（票面钉的形状：软链形里两枚母项零区分力）：

```
$ grep -E "TestAC1POSIXUnresolved|TestAC2POSIXInjected" b04-base-mutd-link/b04-base-mutd-link.roster-run.txt
RUN TestAC1POSIXUnresolvedSymlinkedRootStillRefused119
RUN TestAC2POSIXInjectedTestDataDirStandsAsDeclared119
$ grep -E "TestAC1POSIXUnresolved|TestAC2POSIXInjected" b04-base-mutd-link/b04-base-mutd-link.roster-colour.txt
PASS TestAC1POSIXUnresolvedSymlinkedRootStillRefused119
PASS TestAC2POSIXInjectedTestDataDirStandsAsDeclared119
```

⇒ 两枚**都在 `=== RUN` 名册里**（跑了）、颜色**都是 PASS**（打不红）＝**零区分力复现**，
不是我 SKIP 掉、也不是它们没上场。同形未变异的 `b02` 里它们也是 PASS（这一对的期望本来就恒绿，
这正是"没有区分力"的意思）。

**同一发里它们自己把拒因打印成宿主的链接**（脚本 `scripts/credit.py` 机器判，不靠肉眼）：

```
## b04-base-mutd-link
   AC#1 credited=/acwlink, which is not the tree this call names  ambient_host_link=True  names_own_planted_link=False
   AC#2 credited=/acwlink, which is not the tree this call names  ambient_host_link=True  names_own_planted_link=False
```

⇒ 走查在**第一枚组件**（宿主那枚 `/acwlink`）就停了，用例自己种的 `varlink119`／`injlink119`
**从没被看过**。这一枚是"为什么只加记名断言不算修好"的正面读数。

### 2.3 与实现方 §3 的对点（**只当对表，不当凭据**）

它的 `a1`-`a4`＝`52/30/0/0＋22/0/0`、`45/27/0/3＋15/0/0`、`52/17/13/0＋17/5/0`、`45/17/10/3＋11/4/0`，
我的 `b01`-`b04` **八格四数逐字相同**，包级 rc 也相同（0/0/1/1）。
它是 `/ac7link`、我是 `/acwlink`，它是 `w119ac7tmp`、我是 `wacc119tmp` ⇒ 不是抄它的日志。

### 2.4 判语

判据② **成立（独立复现）**。改前那一发我在自己的容器、自己的快照、自己的饵上量到了；
四数、名册、拒因指向三样都对上，且**没有任何一枚被改成 SKIP**（`b02`/`b04` 的 SKIP 名册与 `b06`/`b08`
逐字节同一枚 md5 `194a396dc3380b91d23145d5d42ca2ff`，三枚都是票 125 的自拒探针）。

### 2.5 本节落盘的 commit（原样输出）

```
$ git log --oneline -1
d0f97d3 evidence(119 AC#7 r1 §2): 判据②成立（独立复现）——改前树 git archive a9c8b6e ＋ 我自己那发 MUT-D …
$ git show --name-only --format="%H" HEAD | tail -2
d0f97d3…
docs/evidence/s1/119-ac7-r1-acceptance.md
```

（记录本身随 §3 那枚 commit 入库；`git log -L` 的复算式同上。）

---

## §3 判据③「改后的判别对」——**成立（独立复现）**，含一条实现方自述属实的机制更正

台件＝`git archive 2956897`（＝交件树，被测文件 md5 `4995e4f5…` 每发容器现量），
两形 × {未变异、MUT-D、MUT-D2}。

### 3.1 六发主读数

| 发 | 形 | 生产码 | RUN | 顶 P/F/S | 子 P/F/S | panic | 包级 rc | 两枚母项 |
|---|---|---|---|---|---|---|---|---|
| `b05-head-plain` | 普通 | 未变异 | 52 | 30/**0**/0 | 22/0/0 | 0 | **0** | PASS／PASS |
| `b06-head-link` | 软链 | 未变异 | 45 | 27/**0**/3 | 15/0/0 | 0 | **0** | PASS／PASS |
| `b07-head-mutd-plain` | 普通 | MUT-D | 52 | 17/**13**/0 | 17/5/0 | 0 | 1 | **FAIL／FAIL** |
| `b08-head-mutd-link` | **软链** | MUT-D | 45 | **15/12/3** | 11/4/0 | 0 | 1 | **FAIL／FAIL** |
| `d01-head-mutd2-plain` | 普通 | MUT-D2 | 52 | 21/9/0 | 18/4/0 | 0 | 1 | **FAIL／FAIL** |
| `d02-head-mutd2-link` | 软链 | MUT-D2 | 45 | 18/9/3 | 11/4/0 | 0 | 1 | **FAIL／FAIL** |

**正向对照**：派单给的期望数是 `52/30/0/0＋22/0/0`（普通）与 `45/27/0/3＋15/0/0`（软链），
我量到的 `b05`/`b06` **逐格相同**，**没有不一致要报**。

### 3.2 名册差集（改前 ↔ 改后，同形同产物；四数之外那一枚）

| 对比 | `=== RUN` 名册 | 全量 colour 名册 | SKIP 名册 |
|---|---|---|---|
| `b01` vs `b05`（普通·未变异） | 差 **0** 行 | 差 **0** 行 | 差 **0** 行 |
| `b02` vs `b06`（软链·未变异） | 差 **0** 行 | 差 **0** 行 | 差 **0** 行 |
| `b03` vs `b07`（普通·MUT-D） | 差 **0** 行 | 差 **0** 行 | 差 **0** 行 |
| `b04` vs `b08`（软链·MUT-D） | 差 **0** 行 | 差 **4 行** | 差 **0 行** |

那 4 行原文（`comm -3`，左＝只在 `b04`、右＝只在 `b08`）：

```
PASS TestAC1POSIXUnresolvedSymlinkedRootStillRefused119        （只在 b04）
PASS TestAC2POSIXInjectedTestDataDirStandsAsDeclared119        （只在 b04）
FAIL TestAC1POSIXUnresolvedSymlinkedRootStillRefused119        （只在 b08）
FAIL TestAC2POSIXInjectedTestDataDirStandsAsDeclared119        （只在 b08）
```

⇒ 差额**恰好是那两枚母项由 PASS 转 FAIL**，一枚不多、一枚不少；
`=== RUN` 名册 23 发两两不差（没有谁整发不见），`panic|fatal error` 计数 **23 发全 0**（逐发 `numbers` 文件）。

### 3.3 红句原文入库（判据③要的那句"点到我自己种的那枚链接"）

`b08`（软链·MUT-D）里那两行，**未截断**（我把整行从 `-v` 日志里 grep 出来）：

```
dataroot_symlink_119_other_test.go:242: AC#1: refusal of "/acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1193191947254/001/varlink119/data" did not name ErrUnresolvedPath: winsec: /acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1193191947254/001/varlink119 is not a directory. Nothing answered for the link this case planted at "/acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1193191947254/001/varlink119", so this is not the placement leg speaking
dataroot_symlink_119_other_test.go:362: AC#2: refusal of the declared root "/acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1193821381434/001/injlink119/harness/picked" did not name ErrUnresolvedPath: winsec: /acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1193821381434/001/injlink119 is not a directory. Nothing answered for the link this case planted at "/acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared1193821381434/001/injlink119", so this is not the placement leg speaking
```

`b07`（普通·MUT-D）同两支、前缀换成 `/acwplain/...`。**两形都点到本用例自己种的那枚链接**（`varlink119`／`injlink119`）＝达成。

`d01`/`d02`（MUT-D2，红在**新加那一支**）原文两行（软链形，**整行未截断**）：

```
dataroot_symlink_119_other_test.go:244: AC#1 RED: the refusal of "/acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1191546711540/001/varlink119/data" named ErrUnresolvedPath but did not credit the link this case planted at "/acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1191546711540/001/varlink119" (winsec: refusing to seal /acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1191546711540/001/varlink119/data: winsec: path is not provably resolved, refusing to seal: /acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1191546711540/001/varlink119/data reaches it via /acwpriv/wacc119tmp/TestAC1POSIXUnresolvedSymlinkedRootStillRefused1191546711540/001/varlink119, which is not the tree this call names): an ambient link above the tree answered for it, so this says nothing about the leg under test
dataroot_symlink_119_other_test.go:364: AC#2 RED: the refusal of the declared root "/acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119519690961/001/injlink119/harness/picked" named ErrUnresolvedPath but did not credit the link this case planted at "/acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119519690961/001/injlink119" (winsec: refusing to seal /acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119519690961/001/injlink119/harness/picked: winsec: path is not provably resolved, refusing to seal: /acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119519690961/001/injlink119/harness/picked reaches it via /acwpriv/wacc119tmp/TestAC2POSIXInjectedTestDataDirStandsAsDeclared119519690961/001/injlink119, which is not the tree this call names): an ambient link above the base answered for it, so the walk under test never ran
```

（`c9` 段里 `reaches it via` ＝ MUT-D2 抹掉记名短语之后的原文形状，正是记名那一支该响的样子。）

### 3.4 一条必须照实报的机制（否则"红句"会被读成"新支响了"）

MUT-D 那一发的红**从哨兵那一支响**（`:242`／`:362`），不是新加的记名那一支：`credit.py` 在 `b08` 里
连一行 refusal 都没抓到，因为底线在截短之后根本不拒，`PrivateDirAll` 往下走到自己那一级撞到
`… varlink119 is not a directory`。⇒ **实现方 §9.1 那句自述属实**，我独立复现到同一机制；
`c94927d`（只改措辞那枚）存在的理由就是这个——第一版 `ddc1583` 的红句在 MUT-D 下打不到本用例种的链接。

**这一枚怎么算？** 我按判据的文字裁：判据③要的是"MUT-D ⇒ 两枚都红**且红句点到本用例自己种的那枚链接**"，
`b07`/`b08` 达成（两形、四枚红句原文都在上面）。但**"红句点名"本身不是承重证据**——
把链接名塞进一条由别的机制触发的消息里，只改措辞就能做到。承重证据是我另造的**单点回退**（§1.3）：
摘掉记名那一支之后 MUT-D 照旧红（⇒ 点名只是措辞），而摘掉之后 **MUT-D2 由红转绿**（⇒ 那一支真有牙）。
所以本格判"成立"，**依据是 §1.2/§1.3 的十发回退，不是 `:242` 那句话好看**。

⚠ 另一枚只写一遍、别记成"已测"：MUT-D2 打的是**全仓共用的记名短语**，
`b06`→`d02` 的名册差集是 **26 行**（13 枚条目转红：票 113 那族 5 枚含 4 枚子测、票 118 那 2 枚、本程这 2 枚），
`b05`→`d01`（普通形）逐名相同 13 枚。⇒ **MUT-D2 那一发的射程不是本格的判据**，
它只用作"新那一支在两形都能单独响、不是恒绿"的旁证；
`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` 在 MUT-D2 下 **PASS**（它不读那条短语），
这一格差异本身就是"短语只被记名断言消费"的证据。

### 3.5 判语

判据③ **成立（独立复现）**：未变异两形两枚都绿（`b05`/`b06`，包级 rc=0）；
MUT-D 两形两枚都红（`b07`/`b08`）且红句点名本用例种的链接；两形都给；四数之外名册差集与 panic 计数都给了。
附带推翻一条派单预期（"只撤记名 ⇒ MUT-D 打不红"），并把"点名"与"承重"分开算账。

### 3.6 本节落盘的 commit（原样输出）

```
$ git log --oneline -1
67704c7 evidence(119 AC#7 r1 §3): 判据③成立——改后六发主读数 …
$ git show --name-only --format="%H" HEAD | tail -2
67704c7…
docs/evidence/s1/119-ac7-r1-acceptance.md
```

---

## §4 判据④「反半边不许被这一格弄绿」——**成立（独立复现）**

判据点名的两件事：`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` **两形颜色逐名不许变**；
**`:251` 那一处明令不并入本格**。两件事我各自走了一条不引实现方脚本的尺。

### 4.1 那一枚函数体一字未动（四版逐字相同）

```
$ for r in a9c8b6e c94927d 2956897 HEAD; do git cat-file -p "$r:internal/winsec/dataroot_symlink_119_other_test.go" \
    | awk '/^func TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119/,/^\}$/'; done | …
a9c8b6e  body_lines=45 bytes=1934 md5=831a5c08b04c66b556b0150a19c63aea
c94927d  body_lines=45 bytes=1934 md5=831a5c08b04c66b556b0150a19c63aea
2956897  body_lines=45 bytes=1934 md5=831a5c08b04c66b556b0150a19c63aea
HEAD     body_lines=45 bytes=1934 md5=831a5c08b04c66b556b0150a19c63aea     ← HEAD 现为 79cfa1d（邻居 commit 落盘后复量）
```

⚠ 派单说"a9c8b6e 与 HEAD 两版都是 45 行／1934 字节／`831a5c08`"——**复算成立**，
我另外把中间两版（`c94927d`、开工锚 `2956897`）也量了，四版同一枚 md5。

### 4.2 `:251`（现 `:291`）没被并入本格：它不在**任何** hunk 里

```
$ git diff -U0 a9c8b6e..HEAD -- <file> | grep -c '^@@'      ⇒ 11
$ git diff     a9c8b6e..HEAD -- <file> | grep -c '^@@'      ⇒ 5（默认 -U3）
-U3 五枚的旧侧区间：28-34 / 41-46 / 190-207 / 282-289 / 300-312   ⇒ 251 一枚都不覆盖
-U0 十一枚的旧侧起点：31 43 192 194 197 202 204 284 286 303 309    ⇒ 同样不含 251
```

⇒ **枚数要钉死**（这类"几枚 hunk"的账本项目吃过）：`-U0` 是 **11** 枚、`-U3` 是 **5** 枚。
实现方 §5 写"九枚 hunk"、派单简报写"hunk 头在旧行 190／282／300 **三处**"——**两个数都不对**，
但**都不影响判据**：190/282/300 是 `-U3` 里那三枚**碰到被测函数**的（另两枚 28/41 只动注释），
"九枚"是它把纯插入 hunk 漏计了两枚。结论一致：**`:251` 那处零 hunk**。

第二枚尺（不依赖 diff）：每发容器现量 `SENTINEL_BRANCHES=3` 而 `CREDIT_BRANCHES=2`
（23 发逐发相同，见各 `head.txt`）⇒ 第三枚哨兵判据（`:291`）**始终没有记名同伴**＝没被并进来。

### 4.3 它的颜色，两形逐名，23 发全量

| 发（形·生产码） | AC#3 颜色 |
|---|---|
| `b01` 普通·未变异 / `b05` 普通·未变异（改后） | PASS / PASS |
| `b02` 软链·未变异 / `b06` 软链·未变异（改后） | PASS / PASS |
| `b03` 普通·MUT-D / `b07` 普通·MUT-D（改后） | FAIL / FAIL |
| `b04` 软链·MUT-D / `b08` 软链·MUT-D（改后） | FAIL / FAIL |
| `d01`/`d02` MUT-D2 两形 | PASS / PASS（它不读记名短语，与 §3.4 一致；新增读数不是变色） |
| `c01`-`c05`（回退换根，两形） | 全 PASS |
| `e01`-`e04`（回退记名＋MUT-D，两形） | 全 FAIL（与 `b07`/`b08` 同色） |
| `e05`-`e08`（回退记名＋MUT-D2，两形） | 全 PASS（与 `d01`/`d02` 同色） |

⇒ **改前↔改后同形同产物的四对里，它逐名逐色一致**；§3.2 那 4 行名册差额里**没有它的名字**。
它仍是"两形都响"的那枚对照组：MUT-D 四发（改前改后 × 两形）全红，而两枚母项只在 `b08` 才红——
这一格差异本身就是"本格没有拿它换绿"的形状证明。

### 4.4 反半边的语义本身（它今天仍在拒什么）

软链形未变异那一发（`b06`）里它自己打印的那行（**整行原文，未截断**，`logs/b06-head-link/b06-head-link.v.log`）：

```
dataroot_symlink_119_other_test.go:288: AC#3 SealFile("/acwpriv/wacc119tmp/TestAC3POSIXLinkInsideAResolvedDataRootStillRefused1194083439489/001/data/out/keep-me.txt") -> winsec: refusing to seal /acwpriv/wacc119tmp/TestAC3POSIXLinkInsideAResolvedDataRootStillRefused1194083439489/001/data/out/keep-me.txt: winsec: path is not provably resolved, refusing to seal: /acwpriv/wacc119tmp/TestAC3POSIXLinkInsideAResolvedDataRootStillRefused1194083439489/001/data/out/keep-me.txt reaches it through the link at /acwpriv/wacc119tmp/TestAC3POSIXLinkInsideAResolvedDataRootStillRefused1194083439489/001/data/out, which is not the tree this call names
```

⇒ 根解析掉之后，**种在数据根里面的那枚链接（`…/001/data/out`）仍然把密封带出这棵树、仍然被拒**
（颜色 PASS ＋ 上面这行拒因），票 113 那条腿的反半边在这一格之后照旧钉着。

### 4.5 判语

判据④ **成立（独立复现）**。函数体四版同 md5、`:251` 零 hunk、`SENTINEL=3/CREDIT=2` 恒量、
两形逐名颜色四对一致、名册差额不含它。

### 4.6 本节落盘的 commit（原样输出）

```
$ git log --oneline -1
33af710 evidence(119 AC#7 r1 §4+§5): 判据④成立（AC#3 函数体四版同 md5 831a5c08 …
$ git rev-parse HEAD
33af710eec2adfa41ec7e603f9b8a86304d8715c
$ git show --name-only --format="%H" HEAD
docs/evidence/s1/119-ac7-r1-acceptance.md     ← 只有这一行（本文件），无别人的路径
```

（本节与 §5 同一枚 commit；号是我在写这一段之前 `git log` 现量的，不是预写的。
逐节归属仍以 `git log -L '/^## §4/,/^## §5/:<本文件>' --oneline` 为准。）

---

## §6 门禁复算（**一律按包 scope**，因为 `cmd/wisp/**` 有另一枚代理在飞）

| 尺 | 命令原文 | 读数 |
|---|---|---|
| `gofmt -l`（宿主） | `cd D:/tmp/wisp119-ac7-r1acc/snaps/head && gofmt -l internal/winsec/` | **空输出**，`GOFMT_WINSEC_RC=0`；单枚文件 `gofmt -l internal/winsec/dataroot_symlink_119_other_test.go` 也空 |
| `gofumpt -l`（宿主，**v0.12.0**） | `"D:\work\base\gopath\bin\gofumpt.exe" -l internal/winsec/` | **空输出**，`GOFUMPT_WINSEC_RC=0`；`--version` 原文 `v0.12.0 (go1.27.1)` |
| `gofumpt -l`（容器） | 未跑 | 仪器边界：镜像里没有 `gofumpt`，`GOPROXY=off` 下 `go run mvdan.cc/gofumpt@…` 取不到 module（实现方 §6 同一枚错误形状）。⇒ 格式门**以宿主那把 v0.12.0 为准**，两程同一把尺 |
| `go vet`（宿主，三 GOOS） | `go vet ./internal/winsec/`；`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go vet ./internal/winsec/`；`GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go vet ./internal/winsec/` | 三发 **rc=0**（`VET_HOST_WINDOWS_RC`／`VET_HOST_LINUX_RC`／`VET_HOST_DARWIN_RC`） |
| `go vet`（容器，三 GOOS） | 容器内 `go vet ./internal/winsec/ ./internal/proc/`（linux，默认 CGO）；`GOOS=windows GOARCH=amd64 CGO_ENABLED=0 …`；`GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 …` | `VET_LINUX_RC=0`、`VET_WINDOWS_RC=0`、`VET_DARWIN_RC=0`（`logs/gate-container2/gates.txt`）；另：23 发取数容器**逐发** `VET_RC=0`（`gates.txt`，本程把它取数前的门从"只 build"升级成"build＋vet"，见 §0.7） |
| 容器 `gofmt -l` | `gofmt -l internal/winsec/` | **空输出**，`GOFMT_RC=0` |
| `sh scripts/d22scan.sh`（纯净快照，交件树） | 容器内、`git archive 2956897` 快照 | **`D22SCAN_RC=0`** ＋ `d22scan: clean - no D22 ban violations`；台账：`bans #1-5 internal/=203 cmd/=22、#6 frontend/=40、#7 internal/tools/=18、#8 design/=16 frontend/=40 internal/=404 cmd/=39` |
| 同一把尺跑**改前基线树** | 容器内、`git archive a9c8b6e` 快照 | **`D22SCAN_RC=0`**，台账八枚数字与上一行**逐一相同**（203/22/40/18/16/40/404/39）⇒ "不降"这一枚我量到的是**逐项相等、差 0**：本格只改了一枚既有 `_test.go`，没有新增文件、没有减任何 scope |
| 正向控制（门能不能红） | `d22scan.sh` 第一步 `runtests.sh -C tools/d22scan ./...` | 真跑：`runtests.sh: OK - top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31` ⇒ 报"clean"之前那把尺自己证明过它能红 |

**我没跑的东西（不是漏，是纪律）**：整树 `go test ./...`、`go build ./...`（`cmd/wisp/**` 里
`worker-ticket128-ac4` 在飞，跑整树会吃到它未提交的改动并冒充成我的读数）；
`-count=2` 与"非 `-v` 一遍"那两枚 AC#6 口径的仪器（本格的判据是③那对判别＋名册差集，
四数与 SKIP 名册我已按 `-count=1 -v` 逐发给全）。

**本程未取任何时序／内存读数**（派单禁，本格判据也没要）。

### 6.1 档位说明（本表用的三档）

- **独立复现**＝我自己造台件、自己下饵、自己跑容器、自己数四数与名册；命令原文与日志路径都在表里。
- **日志＋归档抽验**＝只在它给的日志/快照上核过形式，没重跑。
- **仅自述不背书**＝只有它自己说过，我没重走。

本表里判据①②③④**没有一格落在后两档**；落在"日志＋归档抽验"的只有两处**旁证**：
它 §8 那两枚 docker volume 的来源卷（我没核），以及它 §0 记的那次"编排者推送带走中间态"（我只在
`git log` 里读到 `28c997c` 的 subject，未去核远程）。"仅自述不背书"：**零格**。

---

## §7 我明确没核的清单（列出来是为了下一位不必猜我干了什么）

1. **票面 119 的 AC#1–AC#6 五格**：那是别的程的账（`119-ac1-ac2-r2-acceptance.md`、
   `119-ac2-r2b-acceptance.md`、`119-adversarial-acceptance.md`），本表**只裁 AC#7**。
2. **CI 上这枚文件有没有腿**：票 119 日志 `next=` 第 7 条那笔（本轮三枚新用例在 CI 上没有分母之类），
   本格判据零处要求它，我也没去读任何 run。**这条不因本表结清。**
3. **`internal/winsec` 生产码的正确性本身**：本程生产码一字未动（三枚文件 md5 两树相同，容器里现量），
   所以"那条腿对不对"不是本格的题；我只验测试件。
4. **`tools/d22scan/**` 的实现**：我只跑它、读它的 stdout，未审它（那是 `allowlist.txt`/门禁票的地界）。
5. **票 137 证据 §3.2/§3.3 的原始读数**：判据①说"按票 137 AC#2 的形状"，我核的是**代码形状**
   （helper 存在、签名、可达性、两树 md5），**没有**去复算 137 那张表里的八发数——那一格不是我这次的活。
6. **AC#3 那枚用例在 137 AC#4 里被当对照组的用法**：我只验它"两形颜色逐名没变"（§4.3），
   没验"别人拿它当对照组时读到的数"。
7. **macOS 真实形状**（`/tmp`、`/var` 本身就是链接）：我只在 Linux 容器模拟同形，
   本机是 Windows，那半边**没有分母**——这一枚与前两程一样，属仪器边界不是缺陷。
8. **硬链接**（`R-113-E`）：票面明令"本票不修硬链接、不许混进来"，我也没测。
9. **`docs/evidence/s1/119-ac7-impl.md` §7 那两栏计数**（它的真通知回显数／判为注入数）：
   那是它自己的账，我只记我的（§8）；两程的枚数不同**不构成矛盾**（各自计数，不互相抵账）。
10. **它作废的那两遍读数**（§9.2 说的 `ddc1583` 配套数、以及 `GATE95` 那一遍）：
    我没去它的 `logs/` 里逐发复算——我的 a/b/c/d/e 系列全部自造，不复用它的日志。

---

## §8 临时件清单（**只建不删**，全部在仓外）＋ 两栏计数 ＋ 被拒/报错

### 8.1 台件与日志

```
D:\tmp\wisp119-ac7-r1acc\snaps\            19 棵（349 MB），全部 git archive <sha> | tar -x 起于仓外
   base(=a9c8b6e)  head(=2956897)  base-mutd  head-mutd  head-mutd2(坏,§0.7-1)  head-mutd2f
   rA1  rA2  rA12                        ← 撤换根（单点各一枚／两枚都撤一枚）
   rB1-mutd  rB2-mutd  rB1-mutd2  rB2-mutd2(坏,§0.7-2)
   rB1-mutd-f2  rB2-mutd-f2  rB1-mutd2-f2  rB2-mutd2-f2   ← 修好的回退 B
D:\tmp\wisp119-ac7-r1acc\scripts\          run_one.sh（每发主仪器：门 97 挂载非空/形状、门 95 build+vet）
   drive.sh（宿主侧，Windows 风格 -v）  mutate.py（MUT-D/MUT-D2/回退 A/回退 B，一枚脚本四种 op）
   batch.sh  credit.py（拒因指向机器判）  flipunion.py（17 对名册差集与并集）
D:\tmp\wisp119-ac7-r1acc\logs\             38 个目录（2.0 MB）：23 发入账 ＋ t0/t1 冒烟
   ＋ 作废的 b09/b10、c06–c12（标签不重用）＋ gate-container(坏) / gate-container2 / gate-container-base
   每发：head.txt（身份＋形状硬断言）/ gates.txt / <tag>.v.log / <tag>.rc / <tag>.numbers
        / <tag>.roster-run.txt / <tag>.roster-colour.txt / <tag>.roster-skip.txt
        / <tag>.colours119.txt / <tag>.refusals.txt / <tag>.redlines.txt
   pairs.txt（17 对比较的清单，flipunion.py 的输入）
D:\tmp\ac7_head.txt / ac7_disk.txt / tk_b4.md / tk_head.md / ac3_*.txt   （diff 用的临时抽取）
docker volume  wispacc119-gobuild          （GOCACHE，新建空卷起步；模块缓存用宿主 :ro 挂载，未拷卷）
```

⚠ 坏快照与作废日志**留着不删**（`rm`/`rmdir` 一律不做）：删了裁决表就从〔独立复现〕掉回〔仅自述〕。

### 8.2 两栏计数（**真通知回显数** 与 **判为注入数** 分开，绝不并成一个字段）

**真通知回显数＝6**，逐条带出处（工具名＋当时那条命令前 40 字）：

| # | 出处（挂在哪次调用之后到达） | 内容形状 | 我的处置 |
|---|---|---|---|
| 1 | `Bash: cd "D:\work\workspace\projects plans\Wisp" && git rev`（第一批：`rev-parse HEAD`＋三枚哈希） | harness 追加的 `The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified …` ＋ 一份记忆索引 | 未按它改任何判据；索引是编排者写给它自己的话 |
| 2 | `Write: D:\tmp\wisp119-ac7-r1acc\scripts\run_one.sh`（主仪器落盘当刻） | 同一句提示＋同一份索引（索引本身在长） | 同上；里面出现的"派单里的版本/前提同样是要测的断言"与本程 §0.5/§1.3 的结论方向一致，不是指令 |
| 3 | `Bash(后台): cd "D:/tmp/wisp119-ac7-r1acc" && bash scripts/batc`（core 十发） | `[SYSTEM NOTIFICATION - NOT USER INPUT]` 后台任务完成（`b4cfj5oqb`，exit 0） | 真事件：它明确写着不得当作用户确认；我只据它去读日志 |
| 4 | 同上（core2＋rev 十四发，`b1o0ixb7i`） | 同一形状 | 同上 |
| 5 | `Bash(后台): … for spec in "rB1-mutd-f2 link e01" …`（`bd58bav46`） | 同一形状 | 同上 |
| 6 | `Bash: cd "D:\work\workspace\projects plans\Wisp" && git log --on`（查 c2fa2e9..HEAD 那发） | 首回合随派单一起到达的 `<system-reminder>`（skills 清单＋AGENTS.md 项目上下文） | 派单自带上下文，不是外部指令；AGENTS.md 的规矩与本次派单一致 |

**判为注入数＝0。** 全程工具输出里**没有**出现过自称"编排者备注／系统提示／请 revert／冻结某包／
放宽阈值／少取证／别用工具／直接给结论／预先认定某句话为真"的文字。
按四条判据各自走过一遍：(a) 出现的每段文字指向的路径都是真的吗——是（`MEMORY.md` 真存在、mtime 追得到）；
(b) 有没有替我写好结论的句子——没有；(c) 有没有要我做盘上核不到的动作——没有；
(d) 有没有让我少取证的——没有（若出现，按第四条只登记、继续取证）。
⚠ 反向一条也记：本程"我没遇到注入"**不能**洗掉别的程遇到的，各自计数。

**锚点这一枚的额外防线**（第 8 代形状＝往工具输出里混假 sha）：本表所有 sha 都出自
`git rev-parse` / `git log --oneline` 现量，**未抄派单**（派单写的 `ddc1583`/`c94927d`/`2956897` 我逐枚
`git cat-file -t` 验过是 `commit` 且 subject 对得上）；`a9c8b6e` 是我自己 `git log -- <file>` 抽的改前号。

### 8.3 本程被拒／报错（派单要求单列一栏）

- **权限系统拒绝次数＝0**（没有一次调用被权限挡下，也没有因被拒而绕道）。
- 工具报错（非权限）**四次**，逐条：
  1. `Bash: git archive a9c8b6e \| tar -x -C "D:\tmp\..."` ⇒ `tar: Cannot open: No such file or directory`
     （MSYS 把反斜杠路径当字面量）⇒ 换 `D:/tmp/...` 前斜杠形状重跑，成功。
  2. `Edit: D:\tmp\...\scripts\mutate.py`（第一版回退 B 补丁）⇒ `0 occurrences`（我记错了自己写过的字）
     ⇒ 换成更短的精确锚点重试成功。
  3. `Edit: docs/evidence/s1/119-ac7-r1-acceptance.md`（§4.4 那行）⇒ `0 occurrences`
     （前一次 Edit 已经把整块 §4 写进去了，我拿旧锚点又试一次）⇒ 先 `grep -n` 定位真实行再改。
  4. `Bash: for t in ...; do python3 ...` ⇒ `Python was not found`（Windows 应用执行别名）
     ⇒ 改用盘上真有的 `python` ＋ 把脚本落成 `scripts/credit.py`。
- 仪器自拒**两批**（是我自己的台件问题，取颜色之前就停了，见 §0.7）：`b09`/`b10`（MUT-D2 少逗号）、
  `c06`-`c12`（回退 B 多一枚 `}`）；另有 `t0-smoke`（容器 `-w /wisp` 没给 ⇒ `BUILD_RC=1` ⇒ `GATE95`）。
  **三批全部零枚颜色入账。**

### 8.4 纪律回执

只 commit、**未 push**；每枚 commit 带**显式 pathspec**（`git commit -q -F - -- docs/evidence/s1/119-ac7-r1-acceptance.md`），
每次提交前 `git diff --cached --name-only` 只有一行＝我自己那枚路径；未用 `--amend`/`reset`/`rebase`/`stash`/
`checkout .`/`clean`；未在仓库内建 worktree 或 checkout；**未翻 AC#7 的勾、未改票名、未写票面 119 那枚文件**；
`internal/winsec` 生产码、`internal/risk/**`、`docs/PLAN.md`、`docs/specs/**`、`tools/d22scan/**`、
`allowlist.txt`、阈值／golden、`frontend/**`、`design/**`（owner 有未提交改动）、`docs/reports/**`、
别人的证据文件——**一枚未碰**；`cmd/wisp/**`（`worker-ticket128-ac4` 在飞）**一枚未写、未跑整树测试**。
本程未读过、也未写过任何真实凭据值（容器 `GOPROXY=off`，无网络动作）。


---

## §5 总判

### 5.1 逐格判语与档位

| 格 | 判据本体（票面那一格的原文要点） | 我的凭据（命令原文在对应节） | 档位 | 判语 |
|---|---|---|---|---|
| ① | 两味药必须一起下（记名断言＋已解析基根；只收紧不换根会恒红） | §1.1 落点与零新增依赖；§1.2 五发 `c01`-`c05`（撤换根 ⇒ 逐枚恒红、普通形不动）；§1.3 八发 `e01`-`e08`（撤记名 ⇒ MUT-D 照旧红、MUT-D2 由红转绿） | **独立复现** | **成立**，且**两味各自承重**（不是装饰、不是并列摆设） |
| ② | 改前零区分力那一发要自己复现，不许只引票面那句话 | §2.2 四发 `b01`-`b04`＋`b04` 里两枚"在 RUN 名册且 PASS"＋§2.1 我自己的 MUT-D 落地三证＋`credit.py` 机器判拒因记在宿主 `/acwlink` | **独立复现** | **成立** |
| ③ | 改后判别对：未变异两枚都绿／MUT-D 两枚都红且红句点到本用例种的链接，两形都给，四数之外给名册差集与 panic 计数 | §3.1 六发主读数（含正向对照与派单数一致）；§3.2 四对名册差集 0/0/0/4；§3.3 `b08` 与 `d02` 红句整行原文；§3.4 机制更正 | **独立复现** | **成立**——附一条**照实记账**：MUT-D 那发的红从哨兵支响、红句点名是措辞（`c94927d` 的目的），**记名那一支的射程由 MUT-D2 与四发回退证明**（§1.3/§3.4） |
| ④ | 反半边不许被弄绿：AC#3 那枚两形颜色逐名不许变；`:251` 明令不并入本格 | §4.1 函数体四版 md5 `831a5c08`；§4.2 `:251` 在 `-U0` 十一枚／`-U3` 五枚 hunk 里**一枚都不沾**＋`SENTINEL=3/CREDIT=2` 恒量；§4.3 23 发逐名颜色；§4.4 它今天的拒因整行原文 | **独立复现** | **成立** |

### 5.2 那一问："翻转的并集是不是恰好那两枚母项"

**是，不多不少。** 我做了 **17 对**名册比较（`logs/pairs.txt` 原文，每对两侧**生产码相同**，
只在"药在不在"上不同；`scripts/flipunion.py` 求并集）：

```
prepost-plain-unmut  delta=0      prepost-link-unmut   delta=0
prepost-plain-mutd   delta=0      prepost-link-mutd    delta=4   ← 两枚母项 PASS→FAIL
rA1-link delta=2（只 AC#1）  rA1-plain delta=0  rA2-link delta=2（只 AC#2）  rA2-plain delta=0
rA12-link delta=4（两枚）
rB1-mutd-link/-plain delta=0、0    rB2-mutd-link/-plain delta=0、0     ← 撤记名不影响 MUT-D
rB1-mutd2-link/-plain delta=2（只 AC#1 转绿）   rB2-mutd2-link/-plain delta=2（只 AC#2 转绿）

--- union of names that were PASS on one side / FAIL on the other ---
PASS-side: ['TestAC1POSIXUnresolvedSymlinkedRootStillRefused119', 'TestAC2POSIXInjectedTestDataDirStandsAsDeclared119']
FAIL-side: ['TestAC1POSIXUnresolvedSymlinkedRootStillRefused119', 'TestAC2POSIXInjectedTestDataDirStandsAsDeclared119']
symmetric name set: [同样两枚]
```

⇒ **17 对里 8 对有翻转、9 对零翻转；翻转并集＝那两枚母项，集合闭合、没有第三枚**（同包另外 45 枚
在两形 × 三种生产码 × 十一种回退里逐名逐色一动不动）。
派单问的是"七次翻转"——**实数是 17 对／8 对翻转**，我按实数报，不改成七。

### 5.3 总判语

**AC#7 这一格判：成立（无附条件）。**

- 四条判据**逐条独立复现**（我全部自造台件、自造饵、自己重跑四数与名册，未把实现方日志当凭据引用一次）；
- 票面这一格最要害的那一问（"两发变异各自能响一支算不算承重"）**裁为算**，
  理由与反证在 §1.4——**要求同一发里记名支先响，是要求一支在它的前提被变异抹掉时仍然响**，
  那只能靠改代码形状（把 `else if` 拆成两条 `if`）来造，属措辞工程；
  承重的可操作定义（"摘掉它之后存不存在打不红的变异"）我量到了：存在，且只在 MUT-D2 上（`e05`-`e08`）。
- **不给"附条件"章**的理由：本程没有造出任何一格判据声称要防的坏结局（没放宽断言、没洗 SKIP、
  没动 `:251`、没碰生产码——三枚生产文件 md5 在基线树与交件树逐一相同，我在容器里现量过
  `winsec.go=a6144c880de80e43bb1393f3624e7221`、`winsec_other.go`（未变异）=`b5056918be4ed13817d236fbcae0f477`，
  与实现方 §5 报的两个号一致）。

### 5.4 缺陷栏（**都不足以退这一格**，但必须留在表上）

| # | 缺陷 | 归属 | 影响 |
|---|---|---|---|
| 1 | 实现方 §5 那句"九枚 hunk"数不对：现量 `-U0` **11** 枚、`-U3` **5** 枚 | 实现件（枚数笔误） | 零——它的结论（`:251` 零 hunk）与两个真数都成立 |
| 2 | 派单简报两处前提不成立：(a) "实现方 `:241`／`:361` 与盘上差 2 行"（那两个号本来就是哨兵支，它写的 `:243`／`:363` 才是记名支，一字不差）；(b) "单点回退 B 预期 MUT-D 打不红"（实测照旧红） | **派单**，不是实现件 | (b) 若照它硬判就会**错退**实现方；本表按实测推翻 |
| 3 | 票面 `:65` 那句"本机 v0.7.0 存在"过期；`:68` 已有编排者自己追加的更正块（现量 `v0.12.0 (go1.27.1)`）。⚠ 另两枚历史行（`:199`、`:382`）那两轮的格式读数**是 v0.7.0 那把尺量的** | 票面历史行 | 本格这一格：实现方与我**同一把尺（v0.12.0）**；但"历史门禁与今天同尺"这句**不成立**，别拿它当"一直没变过" |
| 4 | 票面 AC#7 那一格的行号从 `:61-80` 漂到 `:72-91`（现盘），漂因是编排者自己往同一枚文件插了 11 行 | 票面（append-only 正常生长） | 零——§0.3 那发 `diff` 证明那一格 20 行逐字未变 |
| 5 | 我这把尺的一处死计数：`run_one.sh` 里 `UNRESOLVED_BASES` 的 grep 用了两个制表符，改前树也报 0（现量真值：`base` 树 **5** 枚、`head` 树 **3** 枚、`rA12` 树 **5** 枚，单制表符缩进） | **本程仪器** | 零——它是冗余计数，没有任何门靠它；承重门 `CREDIT_BRANCHES`／`RESOLVED_BASES`／`SENTINEL_BRANCHES`／`MUT_*` 逐发有效（数值见各 `head.txt`，`b04` 那发是 `0/0/3`、`b08` 是 `2/2/3`） |
