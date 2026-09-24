# 137 AC#3 —— 终裁（r2 程，写入本族既定的 `137-ac3-r1-acceptance.md` 路径）：只裁 AC#3 那一格

裁决方：本程，**非实现者**（实现方＝`worker-ticket137-ac2`，码 `6f3817a`；它的自述
`docs/evidence/s1/137-ac2-ac5-impl.md` 与上一程终裁件 `137-ac2-ac5-r1-acceptance.md` 都是**被审对象或抽验对象，不是本件的依据**）。

**只裁 AC#3 一格。AC#4 与票面其余各格一枚未裁、一枚未做；票面一个字未动（勾归编排者）。**

⚠ **路径与名字的一笔账**（不是缺陷，是防误读）：台账 `A155` 已把上一枚自称 `acceptor-ticket137-ac3-r1`
的程判为**无凭据**（它报的 47 棵树／66 发读数／commit `99319c7` 盘上一枚都不存在）。本程派单点名要求产出的
文件名仍是 `docs/evidence/s1/137-ac3-r1-acceptance.md`（那枚路径此前**从未被创建**，不存在覆盖谁）。
⇒ **读者按 sha 认本件，不要按"r1"这个后缀去认那枚作废的程**：本件全部落在
`git log --oneline -- docs/evidence/s1/137-ac3-r1-acceptance.md` 现量的那几枚 commit 上（§9 提交账），
而那一枚作废的程在 `git log` 里**零枚 commit**、盘上**零个文件**。

---

## §0 争用闸门（逐轮原样）、锚点、取件、三档口径

### 0.1 争用闸门三行，逐轮原样

这台机器上跑着本仓 self-hosted runner，读数期间**只观察、不杀进程、不改 workflow**。

**第 1 轮（开工时，取件之前）** `at=2026-09-24 13:29:29 +0800`
```
--1-- Runner.Worker / go / compile / cgo
"Docker Desktop.exe","4024",...      ← 宿主 Docker 自身的 GUI，非 CI
"com.docker.backend.exe",...
"docker.exe","19480",...
"cargo.exe","48220","Console","1","9,320 K"
"cargo.exe","54992","Console","1","16,968 K"
（Runner.Worker.exe / go.exe / compile.exe / cgo.exe：**零命中**）
--2-- gh run list --limit 5 --json databaseId,status,headSha
[{"databaseId":35958260537,"headSha":"c8967b8b...","status":"completed"},
 {"databaseId":35953906746,"headSha":"c8967b8b...","status":"completed"},
 {"databaseId":35949710286,"headSha":"45623e40...","status":"completed"},
 {"databaseId":35948947994,"headSha":"8369b24e...","status":"completed"},
 {"databaseId":35947125867,"headSha":"191f0d68...","status":"completed"}]
→ in_progress 枚数 = 0（gh 这一轮**取到数了**，无需降级）
--3-- runner 工作目录 mtime
E:\work\base\actions-runner\_work        2026-09-24 12:02:06.642874100 +0800
E:\work\base\actions-runner\_work\wisp   2026-09-24 12:02:42.547012200 +0800
```
⇒ **未命中**闸门命名的那四类进程；`gh` 五枚全 `completed`；工作目录最后写入在 12:02（开工前 87 分钟）。
**`Runner.Listener.exe` 常驻（它在跑只代表 runner 挂着、不代表有活在跑）**。
⚠ 一轮里出现的 `cargo.exe` ×2 **不在命名集内**（不是本仓 CI，也不是 Go 工具链），本程没有把它判成命中；
它的风险只影响"计时类断言"，本格六发读的全部是颜色与名册、**零计时**。本程**没有**拿它当"机器空了"的反证，也没有因为它而作废任何一发。
（覆盖读数：`ac3-d-tight-link` 一发，容器内时刻 13:33:49。）

**第 2 轮** `at=2026-09-24 13:34:54 +0800`（第一批 1 发之后、剩下 5 发点火之前）
```
--1-- "Runner.Listener.exe","3952",...      ← 只有 listener；无 Worker / go / compile / cgo，rc=0
--2-- gh：同五枚 databaseId、逐枚 status="completed"，in_progress = 0
--3-- _work 2026-09-24 12:02:06 +0800 / _work\wisp 2026-09-24 12:02:42 +0800（与第 1 轮**逐字节相同的 mtime**）
```
⇒ **CLEAR**，剩下五发点火的时刻在这一轮之后。

**第 3 轮** `at=2026-09-24 13:38:29 +0800`（六发全部落地之后，回收用）
```
--1-- cargo.exe 48220 / cargo.exe 54992（与第 1 轮**同两枚 PID**＝同一批没动过的 Rust 进程）、Runner.Listener.exe 3952
     Runner.Worker.exe / go.exe / compile.exe / cgo.exe：零命中
--2-- gh：与前两轮**同一批五枚 run**，逐枚 completed，in_progress = 0
--3-- _work 2026-09-24 12:02:06 +0800（三轮未变）
```
⇒ 六发读数**全部落在三轮都判 CLEAR 的窗口里**；`_work` mtime 三轮未动 ⇒ 期间没有任何 run 被拉起过。
**没有降级成 `CLEAR-LOCAL-ONLY`**（`gh` 三轮都取到数了），所以本件可以说"确认为空"，而不只是"本机侧看不到活"。

⚠ **写作阶段发现的一枚现场事实，如实登记（不改判、不影响上面三轮）**：
本程六发读数落地之后（容器内 13:33:49–13:36:09，闸门 13:29/13:34/13:38 三轮全 CLEAR），
我在 `13:4x` 提交前查工作树时发现仓里多出**一枚不是我建的未跟踪文件**
`docs/evidence/s1/137-ac3-r2-acceptance.md`（19,057 字节、254 行、mtime `2026-09-24 13:42:07 +0800`），
抬头是"137 AC#3 —— 验收方第二程（r2）…读数全部我自己打"，它自己的闸门第 1 轮打的是 `13:34:37`。
⇒ **同一格上有另一枚验收程在飞**（本程派单点名的落点是 `137-ac3-r1-acceptance.md`，两枚路径不撞）。
本程对它的处置：**不读它当依据、不提交它、不改它、不删它**（`137-ac2-ac5-r1` 那一族与 `A155` 都写明过"半份比没有更危险"）。
本格读数是否受并发影响：六发全是**颜色与名册**，**零计时断言**，且 §5 那栏 `RUN == 顶+子(P+F+S)` 的闭合与
`PANIC_LINES=0` 逐发核过 ⇒ 并发负载**没能**把任何一枚用例掀成残局；三轮闸门里也没有第二程带来的 `go.exe`/`compile.exe` 命中。
⇒ 本程**不因此作废任何一发**，但把这件事报给编排者：**"AC#3 被双派"这一格该由谁翻勾，归编排者判**。

### 0.2 锚点与引用到的每一枚 sha：逐枚现量，零枚抄来

- 开工自量 `git rev-parse HEAD` ＝ **`2d029f94cfb126c085b84719340d6afee77f8e80`**（＝派单说的开工锚）。
- 本程中途 HEAD 被编排者推进到 **`21fed31`**（`A155` 那枚），**我的六发读数全在 `2d029f9` 的归档树上**；
  补一条归因账（现量，不是引用它的说法）：`git diff --name-only 2d029f9..HEAD -- internal/winsec/` ＝ **零枚**
  ⇒ 本格结论对当前线仍然成立。
- 逐枚 `git cat-file -t`（**八枚全 `commit`，零枚核不到**）：
  `2d029f9`＝开工自量锚 / `21fed31`＝写作时 HEAD / `4e66817`＝上一程终裁锚 / `f53ad5c`＝AC#2 **之前**那版测试件 /
  `6f3817a`＝AC#2 修法那一枚 / `c8967b8b…` `8369b24e…` `191f0d68…`（闸门里 gh 报的 headSha，抽三枚）＝`commit`。
- `git log --oneline -1 --format=%h -- internal/winsec/placement_symlink_113_other_test.go` ＝ **`6f3817a`**
  ⇒ 那三枚测试件在 HEAD 上仍是 AC#2 收紧后那一版，没被人动过。

### 0.3 派单让我先自量再引用的那两条前提（**我自己重测，没抄编排者的读数**）

| 前提 | 我量到的 | 结论 |
| --- | --- | --- |
| `git diff 4e66817..HEAD -- internal/winsec/` 是否为空 | **输出零字节**（`DIFFLEN:0`） | 成立 |
| 归档 tight 树里 `placement_symlink_113_other_test.go` 的 `md5sum` 是否逐字等于 HEAD 的同一枚文件 | 归档树＝`cb8350cb5e5917a685c10b39d57d5275`；`git show HEAD:… \| md5sum`＝**同值**（另两枚 `20c141df79…`／`50d4e54c72…` 也逐字相等） | 成立 |
⇒ 上一程终裁件 §2.2／§2.3／§2.4 那批 MUT-D 读数量的**就是今天这版码**，可引用为〔日志＋归档，抽验〕。
**但 ①②③ 每一格本程都另有自己亲手跑的那一发**（见 §0.5），所以本格不掉到〔仅自述〕那一档。

### 0.4 取件（严禁工作树当被验版本／严禁仓内 worktree）

```
mkdir -p /d/tmp/wisp137ac3-tree-{tight,d-tight,d-loose}
git archive 2d029f9 | tar -x -C /d/tmp/wisp137ac3-tree-<名字>
```
纯净性逐枚核（**不是"我以为没改"**）：`wisp137ac3-tree-tight` 里那五枚我这一程依赖的文件
`md5sum` 与 `git show 2d029f9:<路径>` 的 blob **逐字节相同** ×5（生产码 `a6144c880de80e43bb1393f3624e7221`／
`b5056918be4ed13817d236fbcae0f477`；测试件 `cb8350cb5e5917a685c10b39d57d5275`／`20c141df792de655408138cdaebf633b`／
`50d4e54c72167d969b361c092cb3e644`）。
`wisp137ac3-tree-d-loose` 的三枚测试件被**换成 `f53ad5c` 那一版**＝`ce181f23bafdf99b7f1513e1ee4ccaef`／
`0ae03ee40b9d44ff0b4e5334d27af1bc`／`b84c85c87fcb2fbbce1393cc5fcabb87`（逐枚 `git show f53ad5c:<路径> \| md5sum` 对点）
⇒ "未收紧"那一发用的确实是 AC#2 **之前**的断言，不是"我以为的旧版"。
仓内**未建 worktree、未 checkout／switch／stash／reset／--amend／rebase／clean，仓内一枚 `go build`／`go test` 都没跑**。

### 0.5 我自己造的变异：第四份独立实现（字节与前手三程都不同）

判据形状引票面 `:47-48`：**两条走查都只走前 3 个组件、仍然拒宿主的链接、永远走不到用例自己种下去的那枚链接**。

- 台件 `/d/tmp/wisp137ac3-io/mutate.py`，标记串 `MUTATION-137-AC3R1-D`（前手三程用的是 `-AC1`／`-R2AC2`），
  **锚串命中数≠1 直接 exit 3**（我没让它有机会"改了个别的地方"）。落点：
  `winsec_other.go:156`（`platformVerifyPlacement` 里 `head := pathPieces(path)` ＋ `if len(head) > 3 { head = head[:3] }`）
  ＋ `winsec.go:284`（`firstLinkAncestor` 里同形截断、`for _, prefix := range head[:len(head)-1]`）。
- 改完后的生产码 md5：`winsec_other.go` ＝ **`184d2c35a59c7512aabd0661c45a3b78`**、`winsec.go` ＝ **`e5d034ee10806c87f796aafabeaed505`**
  ——与终裁方 `29d9c8f5…/0a7113f7…`、上一程引用的另外两前手 `2ffb39b2…`／`9589636253…` **全部不同**，
  两棵变异树之间生产码 md5 逐字相同 ⇒ "态与态之间唯一的差异就是测试件是哪一版"。
- 形状为什么恰好是那一形（**从 `pathPieces` 源码推给读者，再由 §2 的读数验一遍**）：
  `pathPieces` 返回"每多一个组件的一整条前缀"。软链形里 `TMPDIR=/ac3link/w137ac3tmp`、
  用例拼写是 `/ac3link/w137ac3tmp/<用例>/002/root/link/keep-me.txt` ＝ **7 枚组件**，
  宿主的链接是**第 1 枚**（`/ac3link`）、用例自己种的那枚是**第 6 枚** ⇒ 截到 3 枚之后
  两条路都还看得见第 1 枚（照样拒）、第 6 枚**压根没被 Lstat**。
  ⚠ 这一形是**我自己种的**：`/ac3link`＋`/ac3priv`＋`w137ac3tmp`，与上一程的 `/r2ac2link` 无关。

### 0.6 容器与硬闸（每一发开跑前）

镜像 `golang:1.27`（image id **`3680233e3204`**，`docker images` 现量在本地，未拉网络件），
容器内 `go version go1.27.1 linux/amd64`（六发头行逐发打印）。**一发一枚全新容器**（`--rm`）。
- **挂载硬断言**：`[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97` —— **883 是我自己那三棵树现量的字节数**
  （六发头文件里逐发打印 `-rwxrwxrwx 1 root root 883 /src/go.mod` ＋ `md5sum f6ef661732b1851e5c3db348113cb605` ＋ `ls /src/internal/winsec/` 非空）；
  挂载全部 `MSYS_NO_PATHCONV=1` ＋ `/d/...` 形式，**没有一处 `docker -v C:\…`**。
- **形状硬断言**：软链形 `ln -s /ac3priv /ac3link` ＋ `[ -L /ac3link ]` ＋ `readlink` 逐字＝`/ac3priv` ＋ `ls -ld` 打进日志；
  普通形硬断言 `/ac3link` **根本不许存在** ＋ `/ac3plain` 不是软链 ⇒ 错就 `exit 98`。
- **落地不过证就不取颜色**：`go build ./...` ＋ `go vet ./internal/winsec/` 任一非 0 ⇒ `exit 95`。
- **六发逐发读数**：`BUILD_RC=0` ＋ `VET_RC=0` 六发全中，**无一发命中 97／98／95**；`MUT=mutd` 那四发在容器里
  逐发 `grep -n` 出被改行的**原文行号与原文**（`winsec_other.go:156`／`winsec.go:284`）并断言标记命中数＝1/1；
  `MUT=none` 那两发反向断言"整包里不许出现我的标记"。
- 离线：`GOPROXY=off` ＋ `GOMODCACHE=/gomod`；模块缓存**复用既有卷 `wisp137r2-gomod`，但挂的是 `:ro`**
  ⇒ 那一家族"只读不写"的字面也守住了（六发没有一个写它）；构建缓存用本程新建的 `wisp137ac3-gobuild`。
- ⚠ **宿主 wall clock 会跳**：本件每个时间戳**只对它自己那一行成立**，全文**没有任何一处**由两个时间戳相减得出时长。
  （容器内打印的是 UTC，本件同时给 `+0800` 的宿主时刻；换算只做"同一行"内。）

### 0.7 六发读数一览（`go test -count=1 -v ./internal/winsec/`，全部本程亲手）

| 发 | 树 | 断言版 | 生产码 | 形 | RUN | 顶 P/F/S | 子 P/F/S | rc | 分母那 11 枚 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| **`ac3-d-tight-link`** | `wisp137ac3-tree-d-tight` | 收紧（HEAD） | MUT-D | 软链 | 45 | 17/10/3 | 11/4/0 | 1 | **11 枚全 FAIL（格①）** |
| **`ac3-d-tight-plain`** | 同上 | 收紧 | MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 11 枚全 FAIL（格②参照列） |
| **`ac3-d-loose-link`** | `wisp137ac3-tree-d-loose` | 未收紧（`f53ad5c`） | MUT-D | 软链 | 45 | 24/3/3 | 15/0/0 | 1 | **11 枚全 PASS＝零检测力（格③的反向证据）** |
| **`ac3-d-loose-plain`** | 同上 | 未收紧 | MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 11 枚全 FAIL（与收紧版逐名相同） |
| `ac3-tight-link` | `wisp137ac3-tree-tight` | 收紧 | **未变异** | 软链 | 45 | 20/7/3 | 11/4/0 | 1 | 11 枚 FAIL（**没变异也在响**，见 §6） |
| `ac3-tight-plain` | 同上 | 收紧 | 未变异 | 普通 | 52 | 30/0/0 | 22/0/0 | **0** | 基线全绿 |

每发的 `TEST_RC` 与 `LOG_LINES` 在 `/d/tmp/wisp137ac3-io/<发>.head.txt` 里逐行可核；名册与四数由
`parse.py` 从 `-v` 日志**程序化生成**（`<发>.v.ac3run/.ac3fail/.ac3skip/.ac3colour.txt`），**没有一枚是手抄或 `grep -c` 顶上去的**。

### 0.8 三档口径

- 〔独立复现〕＝本程自己在树上跑出来的：**§0.5–§0.7 全部**、§1 的判据面行号、§2–§6 的每一张表与每一个名册差集。
- 〔日志＋归档，抽验〕＝引上一程终裁件 §2.2／§2.3／§2.4 的数字：前提是 §0.3 那两条**我亲自重测成立**，
  并且我另做了一件抽验之外的事（§7.4：我的四份名册与归档日志的名册 **`diff` 无输出**）。
- 〔仅自述，不背书〕＝实现件与票面里我没复算的说法（例：实现方自述的它自己的八数）——本件**一处也没有拿去当判据**。

---

## §1 判据本体（一律引票面原文，不引任何转述）

票面 `.scratch/wisp/issues/137-…-planted-link.md`：

**`:41-42`（AC#3 原句，已被下面的更正块改判）**
> - [ ] **AC#3** 修完之后**复跑 AC#1 那一发同一变异** ⇒ 这 12 枚必须**转红**（红名逐名），并且**普通形零新增红**（八数逐数相同）。
>       这一格是本票的牙：**没有 AC#3 的 AC#2 只是把断言改了个形状**。

**`:44-53`（AC#3 判据更正块 —— 有效判据）**
> **AC#3 判据更正（编排者 09-23 23:2x，来源＝`worker-ticket137-ac1` 交件；append-only、上面原文不抹，冲突以本块为准）**
> 原句"复跑 AC#1 那一发同一变异"是一枚**恒真判据**：AC#1 实测 MUT-A／MUT-B／MUT-AB 那三发在**未修的旧码上就已经全红**
> （11 枚全响）⇒ 拿它当"修好了才算数"的尺子，**改与不改都会绿**，本格会变成一个装饰。
> **正确形状（终裁与后续实现方一律按这条判）**：AC#3 必须复跑 **MUT-D 那一形**（两条走查都只走前 3 个组件、
> 仍然拒宿主的 `/varlink`、但**永远走不到用例自己种下去的那枚链接**）⇒
> ①**软链形 11 枚必须全部转红**（红名逐名）；②**普通形八数逐数不变**（不许新增红、也不许由绿转 SKIP）；
> ③对照组那三枚根已解析的用例（119 族两枚 ＋ `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）在两形都必须仍红，
> 用来证明变异**真的落地了**、那 11 枚的绿不是"没打到"。
> 依据：AC#1 量到 MUT-D 在**未修码**上的形状恰好是"软链形响 0／没响 11（全绿），普通形响 11／没响 0"——
> **这正是本票标题钉的那个害**，所以只有它能把 AC#2 的修法与"没修"区分开。

**`:55-60`（枚数更正块）**
> 真实枚数＝**11 枚 ＝ 7 顶层 ＋ 4 子测试**。那个"8"可复现但**不可采信**……`assertRefused113` 的真位置＝`placement_symlink_113_other_test.go:127`。

**`:198-204`（编排者给 AC#3 起单时钉的三条）**
> **一，只钉 MUT-D 这一发**（票面 `:45` 那条 `>` 更正块已把"复跑 AC#1 同一变异"判成恒真判据，别复活它）；
> **二，软链形"转红"今天不能当凭据**——终裁方独立量到：**未变异的收紧态在软链形里也是那 11 枚红**
> （宿主的链接在第 2 个组件就被拒掉整条拼写）⇒ 转红这件事**修之前就在响**，拿它当"修好了的证据"就是装饰；
> **三，对照组按 `R-137-3` 点名**，并**另禁一条**：不许拿"单发 A／B 里那枚新判据响不响"当凭据。
> ⚠ **可满足性这一条现在是有凭据的了**：同一形里**根已解析的三枚对照**在新判据下**真跑过、且是 PASS 不是 SKIP**……

### 1.1 本程怎么读这堆更正（两处内部不一致，都按更正块＋`R-137-3` 结）

1. **对照组点名**：`:50` 那三枚（含 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）被 `R-137-3`（`:148-150`）
   推翻并改点名成 `TestAC118POSIX…SealFile`／`TestAC118POSIX…PrivateFile`／`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`。
   本程按 **`R-137-3`** 判，并**另把 `:50` 原点名那枚单独量了颜色**（§4.2，结果与 `R-137-3` 那句"它自己也带病"逐发相同）。
2. **枚数**：`:41` 的"12 枚"与 `:49`／`:55-60` 的"11 枚＝7 顶层＋4 子测试"冲突 ⇒ 按更正块取 **11**。
   本程**没有采信任何一方给的名单**，分母是从 `2d029f9` 的源码自己推的：
   `placement_symlink_113_other_test.go` 顶层 `func Test` 在 `:195/:216/:246/:262/:287`（**5 枚拒绝腿**，
   其中 `:216` 里 `t.Run` 在 `:218` 展开 **4 枚 depth 子测**）＋ 同文件 `:304/:338/:386` 三枚是**要成功不要拒的反向腿（不算分母）**；
   `ancestor_separator_108_other_test.go` `:67`／`:133` 两枚**内联腿**（同文件 `:102/:168` 不是拒绝腿）
   ⇒ **7 顶层＋4 子测试＝11**，与更正块逐名相同。`assertRefused113` 在 HEAD 的真位置是 `:174`（票面 `:127` 是收紧前那一版的行号）。
3. **判据点的行号我全部现量**（不抄上一程）：共享尺 `113:145`、标记常量 `113:131`、helper 调用点
   `113:204/237/254/271/296`、108 两枚内联腿 `108:84`→红行 `108:90` 与 `108:151`→红行 `108:156`、118 两处 `118:81/107`。

---

## §2 格①：软链形那 11 枚**必须全部转红**——逐名表，以及"红在哪一行"

发＝**`ac3-d-tight-link`**（我自己的 MUT-D ＋ HEAD 的收紧版断言 ＋ 我自己种的 `/ac3link`）。
四数 `RUN=45 顶 17/10/3 子 11/4/0 rc=1`。

### 2.1 逐名表（取名一律从 `-v` 日志的 `--- FAIL:` 原文行，四枚子测各算一枚；`err@` ＝ 该用例名下日志里出现的 `文件:行` 报错位）

| # | 用例名（`--- FAIL:` 原文，非手抄） | 红在哪一行 | 那一行是不是**新收紧的那处判据** |
| --- | --- | --- | --- |
| 1 | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | `placement_symlink_113_other_test.go:204` | 是（helper 调用点，`t.Helper()` 生效） |
| 2 | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`（顶层母项） | 自身无报错行，由下面 4 枚子测派生 | 派生 |
| 3 | `…/link-at-depth-1` | `placement_symlink_113_other_test.go:237` | 是 |
| 4 | `…/link-at-depth-2` | `placement_symlink_113_other_test.go:237` | 是 |
| 5 | `…/link-at-depth-3` | `placement_symlink_113_other_test.go:237` | 是 |
| 6 | `…/link-at-depth-4` | `placement_symlink_113_other_test.go:237` | 是 |
| 7 | `TestAC1POSIXSealDirThroughASymlinkRefuses` | `placement_symlink_113_other_test.go:254` | 是 |
| 8 | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | `placement_symlink_113_other_test.go:271` | 是 |
| 9 | `TestAC1POSIXSealFileThroughABackslashNamedLink` | `placement_symlink_113_other_test.go:296` | 是 |
| 10 | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | **`ancestor_separator_108_other_test.go:90`（内联腿，不在 helper 里）** | 是（R-137-2 要求的那一枚） |
| 11 | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | **`ancestor_separator_108_other_test.go:156`（内联腿）** | 是（R-137-2 要求的另一枚） |

⇒ **11/11 转红，一枚不落**；`compare.py` 的分母闭合读数是 `DENOM red=11 green=0 other=0 (total 11)`，
且 `denom any-SKIP: []`（**没有一枚是靠 SKIP 变红的、也没有一枚由红转 SKIP**）。

### 2.2 "红出自那处新判据、不是别的分支"——四件独立对点〔独立复现〕

1. **新判据文案的命中数**：`does not credit the link this case planted at` 在 `ac3-d-tight-link` 里 **恰 10 次**，
   逐条给号在 `113:204/237/237/237/237/254/271/296` ＋ `108:90/156` ＝ **10 枚有自己报错行的腿**，
   第 11 枚（顶层母项）由 4 枚子测派生 ⇒ **10＋1＝11，逐枚对号**。
2. **RED 行的净差**：`RED:` 行这一发 **17** 枚、`ac3-d-loose-link` 那一发 **7** 枚，
   **17−7＝10**，且逐类相加对得上（`RemoveUnlinked +2`、`SealDir +1`、`SealFile +6`、`PrivateFile +1`，
   三类前缀在 `113` 与 `108` 的文件行号上逐枚落地）⇒ 多出来的报错行**全部**带新判据文案，没有第十一处来源。
3. **同一发里未受伤**：分母那 11 枚的 `before=… after=…` 读数**逐枚相等**（`before=mode=666 … after=mode=666` ×10、
   `777→777` ×12），即**外来树没被 chmod、没被删**；全发里唯一两处 `666→600` 的实伤**只出现在两枚 118 对照**身上
   （`placement_leaf_118_other_test.go:86`／`:120` 的 AC#6 披露行）⇒ 那 11 枚的红不是附带 Damage 判出来的。
4. **名册差集**：`ac3-d-tight-link` vs `ac3-d-loose-link` 的
   `FAIL only-A=11 / only-B=0`、`RUN only-A=0 / only-B=0`、`SKIP only-A=0 / only-B=0`
   ⇒ 收紧这一刀在软链形里**只多这 11 枚红**，没有用例消失、没有红转绿、没有绿转 SKIP。
   同发里 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 与 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` **仍 PASS**（两版同色）
   ⇒ 这条尺**在软链形可满足**，不是一枚恒不满足的坏尺（R-137-1 反向那一族）。

**红行原文一枚样本**（`ac3-d-tight-link.v.log:112`，逐字，路径里的随机数是我这一发自己的）：
```
    placement_symlink_113_other_test.go:204: AC#1 RED: SealFile("/ac3link/w137ac3tmp/TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTre507587569/002/root/link/keep-me.txt") refused with winsec: refusing to seal …: winsec: path is not provably resolved, refusing to seal: … reaches it through the link at /ac3link, which is not the tree this call names, which does not credit the link this case planted at "/ac3link/w137ac3tmp/TestAC1POSIX…507587569/002/root/link". An ambient link above the tree (the harness's own TMPDIR link is the real shape) can answer for it while the leg under test never runs, so this refusal is not evidence about SealFile.
```
⇒ 被记名的那一枚是**宿主的 `/ac3link`**、被要求的那一枚是**用例自己种的 `…/root/link`** ——
这正是"拒了，但拒的不是它自己种的那枚"这一形，也就是本票标题钉的那个害。**格①成立。**

---

## §3 格②：普通形**零新增红**——八数逐数相同 ＋ RUN/FAIL/SKIP 三档名册双向差集全空

对点＝**同一发 MUT-D、同一棵生产码（md5 `184d2c35…`/`e5d034ee…` 逐字相同）下，收紧版 vs 未收紧版**：

| 发 | RUN | 顶 P/F/S | 子 P/F/S | rc |
| --- | --- | --- | --- | --- |
| `ac3-d-tight-plain` | 52 | 17/13/0 | 17/5/0 | 1 |
| `ac3-d-loose-plain` | 52 | 17/13/0 | 17/5/0 | 1 |
| **参照值（票面文末编排者 `:172` 点名要带的、＝终裁件 `:286`／`:463` 那一发）** | 52 | 17/13/0 | 17/5/0 | 1 |

⇒ **八数逐数相同，且与我亲手那一发逐数相同**（参照值出处我另在 §7.3 追到行号）。

名册（`compare.py` 原样输出）：
```
SETDIFF ac3-d-tight-plain vs ac3-d-loose-plain (RUN):  only-A=0 only-B=0
SETDIFF ac3-d-tight-plain vs ac3-d-loose-plain (FAIL): only-A=0 only-B=0
SETDIFF ac3-d-tight-plain vs ac3-d-loose-plain (SKIP): only-A=0 only-B=0
```
⇒ **三档名册双向差集全空**：不许新增红、不许由绿转 SKIP、也不许有用例消失，三条各自有读数。
并且这一形 **SKIP 名册为空**（`SKIP roster: []`）⇒ 普通形里根本不存在"由跑变跳过"这条路，
分母 11 枚在这一形**逐枚 FAIL**（`DENOM red=11 green=0`，两版同）。

⚠ 一记必须写在明处的**读数口径**（防下一位误用）：普通形这一列的红**不是**新判据判出来的——
`ac3-d-tight-plain` 里新判据文案命中 **0 次**（MUT-D 把走查截断后，普通形里连一次拒绝都产生不了，`err=nil`，
红在**既有**那两条判据上）。所以格②**只能**用八数＋名册差集来核，**不能**拿"新判据响不响"来核
（这正是编排者钉的第三条后半句的同族）。普通形这一列在本格里的身份＝**"收紧这一刀是字面 no-op"的证据**，不是"修好了"的证据。

**基线那一发**（未变异、普通形 `ac3-tight-plain`）＝`52/30/0/0 ＋ 22/0/0、rc=0`，与终裁件 `:91`、
票 124 批次 3b 的基线**逐数相同** ⇒ 分母没被我这一程动过一寸。**格②成立。**

---

## §4 格③：变异真落地的凭据＝对照组两形仍红 ＋ 反向证据"未收紧版软链形 11 枚全 PASS"

### 4.1 对照组（按 `R-137-3` 改点名后的那三枚，不是 `:50` 原本点错的三枚）

| 用例 | `d-tight-plain` | `d-loose-plain` | `d-tight-link` | `d-loose-link` | 未变异两发（plain／link） |
| --- | --- | --- | --- | --- | --- |
| `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed` | **FAIL** | **FAIL** | **FAIL** | **FAIL** | PASS／PASS |
| `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed` | **FAIL** | **FAIL** | **FAIL** | **FAIL** | PASS／PASS |
| `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` | **FAIL** | **FAIL** | **FAIL** | **FAIL** | PASS／FAIL（link 形里它 PASS，见 §5） |

⇒ 四发 MUT-D 读数里三枚对照**逐枚、两形、两版断言全红**，而未变异的两发里它们各自回到 PASS/link-PASS 的形状
⇒ 那 11 枚在 `d-loose-link` 里的绿**不能**读成"变异没打到东西"：同一发里被打到的东西正在红着。
红行原文可读，例 `placement_leaf_118_other_test.go:86`（AC#6 披露 `before=mode=666 … after=mode=600`，
即 MUT-D 之下 `SealFile` 真的**穿过链接 chmod 了外来文件**）、`dataroot_symlink_119_other_test.go` 的 AC#3 RED 行。

### 4.2 `:50` 原点名那一枚——我另重量了（不采信任何一方的转述）

`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 与 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`
在我的四发 MUT-D 里都是 **普通形 FAIL／软链形 PASS**（两版断言同色，收紧没动过它们）
⇒ `R-137-3` 那句"它自己也用 raw `t.TempDir()`、与分母同病、拿它当落地凭据会指错方向"**在我的树上逐发复现成立**。

### 4.3 反向证据（判别对的另一半）

`ac3-d-loose-link`：同一发 MUT-D、同一棵生产码，只把三枚测试件换成 `f53ad5c` 那一版 ⇒
`RUN=45 顶 24/3/3 子 15/0/0、rc=1`，分母 `DENOM red=0 green=11`＝**11 枚全 PASS（零检测力）**，
新判据文案命中 **0** 次，FAIL 名册只剩那三枚对照。
⇒ **判别对**：`未收紧＋MUT-D 软链形 11 枚全绿` ↔ `收紧＋MUT-D 软链形 11 枚全红`，
中间**唯一变量是断言那一版**（生产码 md5 逐字相同、树同源、容器同镜像、形状同种法）。
这才是把"修好了"与"没修"分开的那一发。**格③成立。**

---

## §5 格④：panic 账＝0，以及"我取的数为什么是全量读数"

`grep -cE '^(panic|fatal error)'` **逐份** `-v` 日志（六份全列，不是"抽了几份"）：
```
ac3-d-loose-link.v.log:0    ac3-d-tight-link.v.log:0    ac3-tight-link.v.log:0
ac3-d-loose-plain.v.log:0   ac3-d-tight-plain.v.log:0   ac3-tight-plain.v.log:0
```
`parse.py` 的 `PANIC_LINES=0` 与 `REPETITION_DRIFT=0` 六发逐发同值。

**光有 panic 计数为 0 不够**（那只排除了"日志里有 panic 字样"）。本程另用**分母闭合**证明这些数是全量读数、
不是"一条用例炸掉之后同包其余几十条既不算红也不算绿"的残局——**每一发都满足 `RUN == 顶(P+F+S) + 子(P+F+S)`**，
即每一条 `=== RUN` 都在尾部有**自己的**颜色行：

| 发 | RUN | 顶 P+F+S | 子 P+F+S | 相加 | 闭合 |
| --- | --- | --- | --- | --- | --- |
| `ac3-d-tight-link` | 45 | 17+10+3＝30 | 11+4+0＝15 | 45 | ✔ |
| `ac3-d-tight-plain` | 52 | 17+13+0＝30 | 17+5+0＝22 | 52 | ✔ |
| `ac3-d-loose-link` | 45 | 24+3+3＝30 | 15+0+0＝15 | 45 | ✔ |
| `ac3-d-loose-plain` | 52 | 30 | 22 | 52 | ✔ |
| `ac3-tight-link` | 45 | 20+7+3＝30 | 11+4+0＝15 | 45 | ✔ |
| `ac3-tight-plain` | 52 | 30+0+0＝30 | 22+0+0＝22 | 52 | ✔ |

并且六份日志**逐份末尾都打到包级裁决行**（`FAIL github.com/CarlosShao/wisp/internal/winsec 0.095s…0.121s` 五份、
`ok … 0.143s` 一份）＝二进制跑到头了，没有被中途掀掉。
**跨树跨态名册闭合**：普通形 52 枚 `=== RUN` 名册在 `d-loose`／`d-tight`／纯净三棵树之间**逐名相同**
（`SETDIFF … (RUN): only-A=0 only-B=0`），软链形 45 枚同理 ⇒ 分母在态与态之间没有缩过一寸。
⚠ 软链形那 3 枚 SKIP **逐名恒定**、六发里名册相同，全是票 125 那三枚自拒探针
（`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`／`…SeamGuardStillRefusesEveryHostileShape125`／
`…SeamAcceptsTheHonestPOSIXAnswer125`）——**分母那 11 枚一枚都不是 SKIP**（`denom any-SKIP: []` 逐发）。
⇒ 容器里 TMPDIR 走软链那一形在本格没有缩小分母，缩小的那三枚是它们自己声明的"换台机器就比不了"的探针，
不是我这一格的分母。**格④成立。**

---

## §6 格⑤：恒真自查——**正面作答**

问：**有没有哪一发在"未修"的码上就已经响？** 答：**有，两族，我逐条点名并说明我怎么排除的。**

1. **未修（AC#2 之前那版断言）码在普通形就已经响**：`ac3-d-loose-plain` 里分母 11 枚**全红**（八数 `52/17/13/0＋17/5/0、rc=1`，
   与收紧版**名册逐名相同**）。⇒ 这一列**我没有拿它当修复凭据**，它只用于格②那句"收紧这一刀在普通形是 no-op"。
   （机理：MUT-D 在普通形让底线根本不拒 ⇒ `err=nil` ⇒ 既有判据就红，与新判据无关；新判据文案命中 0 次可核。）
2. **未变异（生产码干净）的收紧态在软链形就已经响**：`ac3-tight-link` 里同样那 11 枚**全红**
   （`45/20/7/3＋11/4/0、rc=1`，新判据文案命中 10 次）。⇒ 编排者钉的第二条**在我的树上原样复现**：
   软链形"转红"这件事**修之前就在响**，所以**本格没把"软链形红"单独当过任何凭据**。
   本件里"修好了"这句话**只由 §4.3 那一枚判别对**（同一 MUT-D 下 loose 全绿 ↔ tight 全红）支撑。
3. **对照组那三枚的两形全红**（§4.1）也**不是**修复凭据——它在两版断言下同色，唯一用途是证"变异落地"。

**排除方式说的是实话，不是姿态**：本件凡是写"成立"的地方，后面都跟着"凭的是哪一发读数"，
而那几发在没有 AC#2 那一刀的树上都是**另一种颜色**（或未跑）。
⚠ 另按 `R-137-1` 与第三条钉：本程**没有造 MUT-A／MUT-B／MUT-AB／MUT-C／MUT-E，也一枚都没跑**——
不是漏了，是它们**被明令禁止当这格的尺子**（A/B 单发是恒不满足的坏尺、票面原句那一发是恒真判据）。
"没有另造它们"这一条本件**不冒充**成读数。**格⑤成立。**

---

## §7 三条钉的独立复核 ＋ 引用归档件前的对点（这一段专门用来"别顺着它硬判"）

| 钉 | 票面／更正块原句 | 我实测到的 | 成不成立 |
| --- | --- | --- | --- |
| 一，只钉 MUT-D | `:199` | 本程唯一造并跑的变异就是我自己那份 MUT-D（第四份字节实现，`184d2c35…`/`e5d034ee…`），未复活"复跑 AC#1 同一变异" | **成立**（并按它执行） |
| 二，软链形"转红"今天不能当凭据 | `:200-201`，理由句"宿主的链接在**第 2 个组件**就被拒" | `ac3-tight-link` 11 枚红、新判据文案 10 次命中，**独立复现**；⚠ **"第 2 个组件"这一句我没复现，也不需要它成立**：我的未变异软链形报错里被记名的也是宿主链接（`/ac3link`＝未变异时 `pathPieces` 的第 1 枚前缀、被拒的原因与它排在第几枚无关）。⇒ 这条钉的**结论**在我树上成立，它附带的**组件序号**在我这里读到的是"第 1 枚"，差一级、不改判 | **成立（附一处措辞对点）** |
| 三，对照组按 `R-137-3` 点名；并禁"单发 A/B 里新判据响不响" | `:202` | 对照三枚在两形四发全红（§4.1）；`:50` 原点名那枚实测 普通 FAIL／软链 PASS（§4.2）＝`R-137-3` 成立；我全程没用 A/B 那一族当凭据 | **成立** |
| ⚠ 可满足性那句 | `:203-204`"三枚对照在新判据下真跑过、且是 PASS 不是 SKIP" | 这一句属 **AC#2 那一程**（未变异软链形）；我在**未变异的软链形**里量到 `ac3-tight-link` 的三枚对照**逐名 PASS 且不在 SKIP 名册** ⇒ 复核成立 | **成立** |

⇒ **三条钉里没有一条被我实测推翻**（第二条附带的那个"第 2 个组件"是措辞级差值，不是判据失败，本程不改判、只登记）。

### 7.4 归档读数的抽验（比"抽验"多做的一步）

除 §0.3 那两条前提之外，本程还把我四发的**名册**与归档 `-v` 日志的名册逐枚 `diff`：
```
ac3-d-tight-link.v.ac3fail.txt  vs  wisp137r2-ac2-io/d-tight-link.v.fail.txt  -> IDENTICAL (14 names)
ac3-d-loose-link.v.ac3fail.txt  vs  wisp137r2-ac2-io/d-loose-link.v.fail.txt  -> IDENTICAL (3 names)
ac3-d-tight-plain.v.ac3fail.txt vs  wisp137r2-ac2-io/d-tight-plain.v.fail.txt -> IDENTICAL (18 names)
ac3-tight-link.v.ac3fail.txt    vs  wisp137r2-ac2-io/tight-link-1.v.fail.txt  -> IDENTICAL (11 names)
RUN 名册（软链 45 枚 / 普通 52 枚）与归档逐名 diff 无输出
```
⇒ 上一程 §2.2／§2.3／§2.4 那批读数**今天这版码上仍然复现**（不同树、不同容器、不同路径种法、第四份 MUT-D 字节实现）。
本格结论**不依赖**这批归档数字（①②③ 各有我亲手那一发），归档只作为一致性对点。

---

## §8 我没核的档（列出来，别当已核）

1. **MUT-A／MUT-B／MUT-AB／MUT-C／MUT-E 本程一枚未造、未跑**（按 `R-137-1` 与第三条钉，它们不许当这格的尺子）。
2. **票 137 其余各格一枚未裁**：AC#4（那 7 处 raw `t.TempDir()`）没碰、没判、**没顺手换根**；
   AC#1／AC#2／AC#5 已由别的程裁过，本程**没有重裁它们**，只在 §1.1 现量了行号。
3. **真 macOS 没验**：软链形＝macOS 的真实形状这件事是票面主张；我只在 linux 容器里造了那一形。
4. **CI 侧没读 run 日志**：只查过 `gh run list` 的 status/headSha 三行，没打开任何一枚 run 的 step。
5. **全仓测试、`-race`、`gofmt/gofumpt/go vet 双 GOOS/-count=2` 那套门禁本程一枚未跑**（那是 AC#5 那一格的地界，已由 `137-ac2-ac5-r1-acceptance.md` §7 判过）。
6. **`-count=2` 未跑**：本格判据要的是 `-count=1` 的逐名颜色；我只在每发里核了 `REPETITION_DRIFT=0`（同一发内不重复取名，因此该值平凡为 0，**别把它当跨复读稳定性的证据**）。
7. **`assertRefused113` 真位置票面写 `:127`、HEAD 上在 `:174`**：我按 HEAD 现量，没去改票面（勾与文字归编排者）。
8. **`TestAC3*…119`（对照那枚）与 118 两枚的红行**我只做了颜色与行号归因，**没有**去核它们在 AC#2 之前是否也满足"记名者＝自己种的链接"这一新判据（那是 AC#2 那一格已结的账）。

---

## 总判

| 格 | 判据（更正块原文编号） | 判定 | 凭哪一发 |
| --- | --- | --- | --- |
| ① | 软链形 11 枚必须全部转红（红名逐名）＋红在哪一行 | **成立** | 本程 `ac3-d-tight-link`（逐名表 §2.1 共 11 名，含 4 枚子测各算一枚；红位 `113:204/237/254/271/296` ＋ `108:90/156`；净差 17−7＝10 枚新判据文案；`RUN/FAIL/SKIP` 差集见 §2.2） |
| ② | 普通形八数逐数不变，且三档名册双向差集全空 | **成立** | 本程 `ac3-d-tight-plain` vs `ac3-d-loose-plain`：八数逐数相同（＝参照值 `52/17/13/0＋17/5/0、rc=1`），三档 `only-A=0 only-B=0`；基线 `ac3-tight-plain` rc=0 |
| ③ | 变异真落地（对照三枚两形仍红）＋ 反向证据 | **成立** | §4.1 四发全红名册 ＋ §4.3 `ac3-d-loose-link` 11 枚全绿 |
| ④ | panic 账＝0，且取的是全量读数 | **成立** | §5：六份日志逐份 0，加分母闭合 `RUN == 顶+子 (P+F+S)` 六发全中、末尾均达包级裁决行、名册跨树逐名相同 |
| ⑤ | 恒真自查正面作答 | **成立**（§6 点名两族"未修码上就已经响"的读数并把它们排除在修复凭据之外） | §6 |

**AC#3 一格总判＝成立（PASS、无附条件）。**
出口是"成立"而不是"附条件通过"：**本程自己造出了那一发**（第四份字节实现的 MUT-D、自己种的 `/ac3link` 软链形、
自己重跑的六发判别读数），没有一格需要"等下游补料才算数"。⇒ **"附条件 vs 退回"的分界这条，本程落在"自己造出来了"那一侧。**

给编排者的四味（都不是退回项，也不用我动手）：
1. **文件名与名字的账**：本件写在派单点名的 `137-ac3-r1-acceptance.md`，但 `A155` 已把"`r1`"这个程名作废
   ⇒ 若你希望文件名也体现 r2，`git mv` 归你；本程**没有**为这事去改票面或台账。
2. 票面 `:50` 那句对照组点名（含 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）**已被 `R-137-3` 推翻、我复核成立**；
   它与 `:202` 那条钉并存不冲突，但读者若只看到 `:50` 会点错名 —— 要落一句更正归你。
3. 票面 `:41` 那句"12 枚"与 `:49`／`:55` 的"11 枚"并存；本件按 11 判，**票面文字一字未动**。
4. §7 第二条里那处"第 2 个组件"与我读到的"第 1 枚前缀"差一级：**不影响判据**，但别在 AC#4 里把它当已核事实引用。

**没做的事**：没改任何生产码／测试码（`internal/winsec/**` 在仓里 status 零行）、没动票面与 `docs/reports/**`、
没动任何冻结面、没为变绿放宽任何断言（本程连一次 `Edit` 生产文件都没做过，变异只落在 `/d/tmp/wisp137ac3-*` 三棵快照树）、
没在仓内建 worktree／checkout、临时件**只建不删**、**只 commit 不 push**。

**本程新建的临时件（只建不删，路径全列，清点归编排者）**：
`D:\tmp\wisp137ac3-io\`（台件 `mutate.py`／`run.sh`／`go.sh`／`batch.sh`／`parse.py`／`compare.py` ＋
`batch.out` ＋ 六份 `*.head.txt` ＋ 六份 `*.v.log` ＋ 每份日志四枚名册文件 ＋ `*.build.err`／`*.vet.err`）、
三棵代码快照树 `D:\tmp\wisp137ac3-tree-tight`／`-tree-d-tight`／`-tree-d-loose`、一枚新 docker 卷 `wisp137ac3-gobuild`。
只读引用了 `wisp137r2-ac2-io`（归档日志）与 `wisp137r2-gomod`（**挂 `:ro`**）；`wisp137r2-*` 那一族**一枚未写、一枚未删**。

**注入计数（两栏，一个字段不装两种含义）**
- 真通知回显数＝**2**：两次 `MEMORY.md` 的"file was modified since it was last read"回显
  （①`C:\Users\swq\.qoder-cn\projects\D--work-workspace-projects-plans-Wisp\memory\MEMORY.md`、
  ②`C:\Users\swq\.qoder-cn\memory\MEMORY.md`），都随 `Write` 台件到 `D:\tmp\...` 的那一发 `Bash`/`Write` 之后出现。
  核法：两枚路径**盘上真存在**（本程 `Read` 过其一，另一枚按路径同族＋编排者自己写过记忆的既有记录判为真），
  内容**是记忆索引、没有一句要求我少取证／revert／放宽判据**。⇒ 判真通知（多半就是编排者在写记忆）。
- 判为注入数＝**0**：本程工具输出里没有出现任何点名本仓路径的"编排者备注／系统提示／编码规则已更新／请 revert／
  冻结某包／放宽阈值／Confirm the harness note is genuine"形状的文字。四条辨别力逐条走过：
  没有一处要我少取证（我反倒自己加了第六发基线与四次名册 `diff` 对点），没有一处内容削弱 owner 权威。

**凭据反扫（按词筛，不按长度）**：本件写入后对 `sk-`／`token`／`secret`／`password`／`api_key` 逐词扫，
命中见 §9 提交账里那条命令的原样输出；本仓这类命中历来是**测试名／文件名**（如 `…WritesNothing` 那枚用例的
`PrivateFile(...)` 样本路径），本件为免歧义**没有**引用那一行的完整拼写。零凭据值入文。

---

## §9 提交账（`A155` 钉的那条硬规矩：正文里"已提交"必须带这两行原样输出）

（本节在第一次提交**之后**追加，故此处只放已经存在的两行；本件后续每一枚 commit 由 §9 追加一块。）

**commit #1 ＝ `4f61430`**，落地时刻 `date` 原样 `Thu Sep 24 13:44:43 2026 +0800`（只对本行成立，不参与任何时长计算）。
两条命令原样输出：
```
$ git log --oneline -1
4f61430 docs(137,AC#3 终裁 格①–⑤ 全判成立): 只裁 AC#3；六发我自己重跑（锚 2d029f9，第四份字节实现的 MUT-D、
自己种的 /ac3link 软链形）——①软链形 11 枚逐名全转红…（正文过长，此处折行，完整标题以 `git log -1` 为准）

$ git show --name-only --format="%h%n%ad" HEAD
4f61430
Thu Sep 24 13:44:43 2026 +0800

docs/evidence/s1/137-ac3-r1-acceptance.md          ← 暂存清单里**只有这一枚路径**

$ git rev-parse HEAD
4f61430d3d259066f0ea4c573e70575e8ea54544

$ wc -l < docs/evidence/s1/137-ac3-r1-acceptance.md
478
```
提交前我另跑过一次 `git diff --cached --name-only`，输出＝`docs/evidence/s1/137-ac3-r1-acceptance.md` 一行；
`git status --porcelain` 当时另有 `?? docs/evidence/s1/137-ac3-r2-acceptance.md`（§0.1 那枚在飞的兄弟程的），
**未被 add、未被提交、未被改、未被删**。

**入库前的凭据反扫（按词筛，命令与命中数原样）**：
```
$ grep -oniE "\bsk-[a-z0-9]+|\btoken\b|\bsecret\b|\bpassword\b|\bapi_key\b|bearer|\bghp_[a-z0-9]+" \
    docs/evidence/s1/137-ac3-r1-acceptance.md
459:token
459:secret
459:password
459:api_key
```
⚠ **行号口径**：上面的 `459` 是**第一次扫描当时**的行号（那一次跑在 §9 写入之前，命中的是总判里"凭据反扫"那一段，
现在它在 `:470`）——行号位移是**我自己追加 §9 造成的**，不是别人的活。
**追加 §9 之后我又重扫了一遍（原样；重扫那一刻文件 527 行，本节之后又长了几行 ⇒ 这里的行号只到那一次扫描成立）**：
```
$ grep -oniE "\bsk-[a-z0-9]+|\btoken\b|\bsecret\b|\bpassword\b|\bapi_key\b|bearer|\bghp_[a-z0-9]+" \
    docs/evidence/s1/137-ac3-r1-acceptance.md
470:token   470:secret   470:password   470:api_key     ← 总判那段"按词筛"的句子自己
505:bearer                                               ← 本节这条 grep 的**模式**
507:token   508:secret   509:password   510:api_key      ← 本节贴进来的上面那份扫描输出
```
⇒ 九枚命中**逐枚都是规则句／命令／自己的输出**（自指），仍零凭据值；按**形状**判全为假阳性，
**没有因为扫到东西就删掉任何一条规则、没有改判据**。宿主环境里同类变量我只打印**变量名**（一枚 `VISION_API_KEY`），
**值零处入文、零处入对话**。
这条写下来是为了让下一位**不必再猜**："按词筛"这把尺扫自己的时候必然命中，命中数从 4 变 9 的那一段增量
是**我贴扫描输出贴出来的**，不是有人往文件里塞了凭据。**今后引用带行号的命中，要连"哪一次扫描、当时文件多少行"一起给。**



**commit #2（本枚）追加的就是上面这一块**；它自己的两条命令不在本文里（写了就成循环），
读者一条命令可核：`git show --name-only HEAD` ⇒ 期望只出现 `docs/evidence/s1/137-ac3-r1-acceptance.md` 一枚路径。
本件正文里除本节这两块之外，**没有任何一处**出现"已提交"这句话。

---

`next=` **本程没有未裁的格**：AC#3 的 ①②③④⑤ 五格全部裁完并各发一发亲手读数；
票 137 其余各格（AC#4 与已判过的 AC#1／AC#2／AC#5）**按派单不属本程、一枚未动**。
留给编排者的三件（都不是我的活）：①同格双派该以哪一枚落勾（§0.1 登记的 `137-ac3-r2-acceptance.md` 在飞）；
②票面 `:50` 与 `R-137-3` 的点名冲突、`:41` 的"12 枚"与更正块"11 枚"并存（§7 表下给了逐条出处）；
③§7 第二条那处"第 2 个组件"vs 我读到的"第 1 枚前缀"，别在 AC#4 里当已核事实引用。

