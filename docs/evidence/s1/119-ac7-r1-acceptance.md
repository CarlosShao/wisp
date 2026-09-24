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
- **台件全在仓外**、全部 `git archive <sha> | tar -x`，共 16 棵快照（清单 §6）；**禁读脏工作树冒充被验版本**。
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
