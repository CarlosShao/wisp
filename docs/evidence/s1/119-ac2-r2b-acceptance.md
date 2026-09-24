# 票 119 · AC#2 复判（第二程 r2b）· 对抗验收表

- 验收者：`acceptor-ticket119-ac2-r2b`（非实现者，与实现方不是同一程）
- 只裁：**AC#2 复判**。AC#1／AC#3／AC#4／AC#5／AC#6 不裁、不碰；票面一枚勾不许翻。
  `AC#7` 只回答 §3 的顺手一问 A（只判定，不改正文、不实现）。
- 票面：`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md`
- 第一轮表（别人的，只当起点）：`docs/evidence/s1/119-adversarial-acceptance.md`
- 本机 `date` 现量：开工 `2026-09-24 16:03:15 +0800`（后续每节写"几点"前重新量）

## §0 建件回执

> 本节按要求**先建件、先提交**，再往下取证。三条硬要求的第 1、2 步都成功了：
> `Write` 没被拒、`git add`／`git commit` rc=0。**所以本程继续干活。**

开工锚点（自量，不读别人的号）：

```
$ git rev-parse HEAD
1b323ba494307c80e6fe6ddc5f065b9cbff8705a
$ git cat-file -t HEAD
commit
$ git rev-parse --abbrev-ref HEAD
dev
$ git diff --cached --name-only          # 提交前暂存清单核对
（空）
```

第一次 commit（`git add -- docs/evidence/s1/119-ac2-r2b-acceptance.md` ＋
`git commit -q -F - -- docs/evidence/s1/119-ac2-r2b-acceptance.md`）之后，
要求的那两条命令的**原样输出**：

```
$ git log --oneline -1
3a49745 evidence(137 AC#4 r1 终裁 §0): 锚点自量 a9c4d58 + b1010ff..HEAD -- internal/winsec/ 为空 + 闸门 + 被验版本盘上身份

$ git show --name-only HEAD
commit 3a497457cfc5ea4564749cbbf80b620cdc210b71
Author: CarlosShao <1933942520@qq.com>
Date:   Thu Sep 24 16:03:45 2026 +0800

    evidence(137 AC#4 r1 终裁 §0): 锚点自量 a9c4d58 + b1010ff..HEAD -- internal/winsec/ 为空 + 闸门 + 被验版本盘上身份

    非实现者终裁表开工。§0 只报四件事：
    - 锚点 git rev-parse HEAD = a9c4d58f99b134b0ab3253b80469dde92bb2692a（cat-file -t = commit、分支 dev）
    - 派单要求的那条未动证明：git diff b1010ff..HEAD -- internal/winsec/ 为空（rc=0）⇒ 不需要停手
    - 闸门：gh run list 三枚全 completed、每发前 GATE-PRE 的 golang-containers=0、docker server 29.6.2 现量、
      golang:1.27 本机已有未 pull
    - 被验版本盘上身份：swap 树 = git archive b1010ff（RAWROOTS=0 SWAPPEDROOTS=9）、
      base 树 = git archive a02da50（= 9c0f546^，比实现方的 182daed 更严的父；两者 winsec 逐字节相同）；
      两树 internal/winsec/ 现地 diff = 42 变动行、其中 root := 命中正好 7 对

    另登记本程一枚仪器自拒（发生在读数之前、非树错）：空 GOMODCACHE 撞 gate 95，
    以 :ro 源 cp -a 建自己的 wisp137ac4r1-gomod（726M→CP_RC=0→726M）后复跑同发才开批。
    工作树里 design/** 那 18 行别人的活一枚未动。AC#4 的勾没翻。

docs/evidence/s1/137-ac4-r1-acceptance.md
```

⚠ **这两条读数量的不是我这枚 commit**，而是**同机另一枚验收程（票 137 AC#4）在 22 秒后压上来的那枚**。
共享工作树里这是常态，所以把两件事都留下：上面那对是"当时 HEAD 的原样输出"，下面这对是**我这枚的盘上身份**。

```
$ git show --name-only 20f8c22      # 我这枚骨架 commit 自己
20f8c22 16:03:23 docs(evidence): 票 119 AC#2 复判表建件（r2b 第一段，仅抬头＋§0 标题）

docs/evidence/s1/119-ac2-r2b-acceptance.md

$ git log --format='%h %p %ad %s' --date=format:'%H:%M:%S' -3     # 父子链
3a49745 20f8c22 16:03:45 evidence(137 AC#4 r1 终裁 §0): 锚点自量 a9c4d58 + b1010ff..HEAD -- internal/winsec/ 为空 + 闸门 + 被验版本盘上身份
20f8c22 1b323ba 16:03:23 docs(evidence): 票 119 AC#2 复判表建件（r2b 第一段，仅抬头＋§0 标题）
1b323ba d55c00d 16:02:04 docs(A163; 一枚假交件): 同一程连报两枚互不相同的"完成"（692 行/六枚号、314 行/三枚号），而盘上零产物、九枚号全部不存在——票 119 复判至今仍未做
```

⇒ 骨架那枚＝ **`20f8c22`**，只动 `docs/evidence/s1/119-ac2-r2b-acceptance.md` 一条路径，
父＝开工锚点 `1b323ba`（＝ `A163` 那枚登记"上一程假交件"的 commit）。三枚号都过了 `git cat-file -t`／
`git log` 链核对。**〔独立复现〕**

### §0.1 本程遇到的一次工具输出异常（按注入登记，不据此停手）

第一次跑 `git log --oneline -1`（命令前 40 字：`git log --oneline -1; echo "=== show ==="; gi`，
工具＝`Bash`）时，回显给我的**不是**我这枚，是一行我从未见过的号与主语：

```
5f7e104 docs/evidence/s1: 票 119 AC#2 复判表建件（骨架，§0 第一次提交）
```

现量反扫：`git cat-file -t 5f7e104` → `fatal: Not a valid object name`；`git branch --contains 5f7e104` 同错；
dev 近 12 枚里**没有**这个号也没有这个主语。同批第二条命令（`echo` 打出来了、两条 git 却给了**空输出**）
也是同一发异常。
⇒ 形状＝记忆里登记的"第 8 代注入＝冒充我们自己的锚点"（假 sha 混进工具输出），
且**票 119 上一程的死法就是报了九枚全部不存在的号**。本程处置：**只登记、不服从、锚点一律自己 `rev-parse` 现量**，
后续所有 sha 落笔前先 `git cat-file -t`。异常本身不改变任何判据。

### §0.2 争用闸门（开工时现量）

```
$ docker ps --format '{{.ID}} {{.Image}} {{.Status}}'
e8e44eb98212 golang:1.27 Up 12 seconds          ← 同机 acceptor-ticket137-ac4-r1 在跑
eee61b846e26 union-api-proxy:local Up 8 days
...（其余为与本机项目无关的常驻容器）
```
⇒ **golang 容器非 0 ⇒ 本程先做不取数的活**（读票面、读第一轮表、量 `git log/show --numstat` 的三档改动账），
取数批次等 `gh run list --branch dev` 无 `in_progress` 且 `docker ps` 里 golang 归零再开，撞负载的批次作废重跑。

### §0.3 同格双派登记（不是我的活，但我看见了）

`git log --oneline --grep='119'` 现量：`6e04d1a`（本机 **16:05:53**，比我 `fc71fa6` 晚 2 分钟）落了一枚
`docs/evidence/s1/119-ac1-ac2-r2-acceptance.md`（78 行，只到 §0，主语"119 复判r2"）——
**与本程同格（AC#2）的第二枚验收程正在跑**。处置：那枚文件我**一字未读其结论、未改、未替它翻勾**；
两边快照目录不同名（本程一律 `wisp119r2b-*`）。**同格两份表不互相抵账**，谁的读数算谁的。

---

## §1 AC#2 的三档改动账（`internal/winsec/winsec_other.go`）

现量时刻 `date` ＝ `2026-09-24 16:0x +08`（§0 那条 16:03:15 之后、本节 git 读数之前）。
本程锚点（本节全部读数取自这一枚树）：`git rev-parse HEAD` 在 §0 量到 `1b323ba`，
提交本节后随共树移动到别人的号；**`internal/winsec/winsec_other.go` 自 `034080c` 起零漂移**（下面第 4 条读数）。

**① 那一枚文件的完整历史（`git log --oneline -- internal/winsec/winsec_other.go`，本程自己跑）**

```
034080c docs(119 返修,R-119-2/R-119-3/R-119-10): 注释收窄成能站住的形状，并把 R-119-3 那一刀裁清（纯注释，判定分支 0 hunks）
189cb1e fix(119,AC#1-AC#6): 票 113 那条 POSIX 链接腿不再把"OS 自己把数据根写成软链"当攻击——语义裁 ②，winsec 判定分支一字未动
3c5d1c3 fix(113,AC#2#3#4#5): 给 POSIX 的 platformVerifyPlacement 补上链接腿——复用 pathPieces + ancestorIsLink，只许拒不许改写
0717bf2 fix(103,AC#1#2) / 7910bcd fix(94,AC#1-AC#3) / de15a6b feat(89,AC#2,AC#3) / b994a2c feat(89,AC#2)   ← 更早，与本格无关
```
⇒ **票 119 的返修只有一枚动了这个文件：`034080c`（09-22 17:11）**；`189cb1e` 是第一轮被退回的那枚**交件**本身。

**② 逐枚 `git show --numstat`（只贴 `internal/winsec/` 那几行，全文件清单在下面一条）**

| commit | 时刻 | `winsec_other.go` | 同一枚里其它 `internal/winsec/` 路径 |
|---|---|---|---|
| `189cb1e` 交件 | 09-21 22:49 | **21 增 / 6 删** | `dataroot_symlink_119_other_test.go` 271 增 / 0 删（新文件） |
| `33c8acd` 返修 R-119-9 | 09-22 17:03 | 未动 | `dataroot_symlink_119_other_test.go` **76 增 / 20 删** |
| `36294c2` 返修 R-119-1 | 09-22 17:03 | 未动 | 无（动的是 `cmd/wisp/secret.go` 15/0 + 新 `cmd/wisp/secret_dataroot_119b_test.go` 259/0） |
| `4f19ec6` 返修（票面） | 09-22 17:08 | 未动 | 无（只动词面 71/0） |
| `034080c` 返修 R-119-2/3/10 | 09-22 17:11 | **41 增 / 14 删** | 无 |

累计尺（本程自己跑，三枚号都是 `git diff --numstat <a>..<b> -- internal/winsec/winsec_other.go` 的原样输出）：

```
ce666ea..189cb1e   21  6
ce666ea..034080c   53  11
ce666ea..HEAD      53  11      ← 与上一行逐字相同
034080c..HEAD      （空输出）   ⇒ 该文件自返修以来零漂移（票 137 AC#4 那批只动了 _test.go，没动它）
```

⚠ **一处措辞精度（记账，不改判）**：票面 `:303` 写"本轮 ③④ 又加了 **32/5**"。
`32/5` 是**相对交件态的净累计差**（53−21 / 11−6），`41/14` 是 `034080c` 那一枚 commit 自己的 numstat
（改注释时重排/替换会同时产生增与删）。两个数都对，**但它们是两把不同的尺**——
下一个引这句的人必须带上"哪把尺"，否则又是一次归因腐坏。**〔独立复现〕**

**③ 三档分类（只有第一档算实质）**

| 档 | 内容 | 行数（`ce666ea..HEAD`） | 本程怎么量的 |
|---|---|---|---|
| **一档·判定分支** | `platformVerifyPlacement` 里 `if ancestorIsLink(prefix) {` 那三行，及包内任何可执行语句 | **0 / 0** | 见下面两把尺 |
| **二档·注释** | 包文档"两条成本"块、`cleanSpelling119` 的说明、被证伪那句全称的收窄 | **53 增 / 11 删（＝全部）** | 两把尺的差集 |
| **三档·其它** | import、签名、常量、`//go:build` 行、空白形状 | **0 / 0** | 两把尺的差集 |

两把尺（**都是本程自己现跑，不引第一轮**）：

```
尺 A（剥尽注释与空行后逐行比）：
  git show ce666ea:…winsec_other.go | grep -vE '^[[:space:]]*($|//)'   = 59 行
  git show HEAD    :…winsec_other.go | grep -vE '^[[:space:]]*($|//)'   = 59 行
  diff 两文件 | grep -c '^[<>]'  = 0        diff -u | grep -c '^@@' = 0        ← 0 hunks
尺 B（票面自己那把，换我重跑）：
  git diff -U0 ce666ea..HEAD -- …winsec_other.go | grep '^[+-]' | grep -v '^[+-][+-][+-]'
        | grep -vcE '^[+-][[:space:]]*//|^[+-][[:space:]]*$'   = 0
  同尺分区间：ce666ea..189cb1e = 0        189cb1e..034080c = 0
```
（注：`grep -c` 命中 0 时 rc=1，那是"无匹配"不是"命令失败"；上面引的是**计数值**。）
判据分支仍在位（`git show HEAD:… | grep -n`）：`platformVerifyPlacement` 在 **:155**、
`if ancestorIsLink(prefix) {` 在 **:157**。⇒ **一档零行，二档吃掉全部 53/11，三档零行。**〔独立复现〕

**④ 同包其余生产文件（第一轮 §四 的那几枚 md5，本程自己重跑，并按"到哪一版"分栏）**

| 文件 | `ce666ea` | `189cb1e` | `034080c` | `HEAD` |
|---|---|---|---|---|
| `internal/winsec/winsec.go` | `a6144c88…e7221` | 同 | 同 | **同**（四版一致，复跑两次相同） |
| `internal/winsec/winsec_other.go` | `daf9d799…032ff2` | `457b556c…68cdc6` | `b5056918…0f477` | **同 `034080c`**（`034080c..HEAD` 零 diff） |
| `internal/winsec/winsec_windows.go` | `fb20bca5…296b6` | 同 | 同 | **同** |
| `internal/winsec/resolve.go` | `7eb8a754…33074` | 同 | 同 | **`b6876a5e…129c3` ＝ 已变** |

⇒ 第一轮记的"到交件态四枚一字未动"仍然成立；但**"resolve.go 到 HEAD 仍是 `7eb8a754…`"今天不成立**。
本程现量 `git diff --numstat 034080c..HEAD -- internal/winsec/resolve.go` = **205 增 / 16 删**，
来自**三枚别的票**（`git log --oneline 034080c..HEAD -- internal/winsec/resolve.go`）：

```
a45b2e9 09-22 23:41 feat(129,AC#1+AC#2): 危害量成读数——过缝的候选把 S-1-1-0 从受害者树上剥掉了; 绝对性进两枚比较(只更严)
4824bb8 09-22 20:58 fix(winsec,125,AC#2): 守门人造探针前先解析它问 OS 的那枚答案——软链 temp 不再把整条 C26 挡在缝外
a701138 09-22 19:59 fix(winsec,126): 把卷段纳入 sameTree/answerInsideTree——缝守那条腿实测会被跨卷作答的候选走过
```
**这不是票 119 的债，但它是一枚会咬人的引文**：谁再拿"resolve.go md5 与第一轮一致"当 AC#2/AC#4 的旁证，
**必须带上"到哪一版为止"**（`…到 `034080c` 为止一致；HEAD 起已被 126/125/129 换掉）。**〔独立复现〕**
（⚠ 本程第一次跑那张四栏 md5 时，`HEAD` 那行回显给我的却是旧串 `7eb8a754…`；
我按 `md5sum` 复跑三次 ＋ `git diff --numstat` 一条才发现是**回显错**、不是树错。
异常已进 §6 注入/异常栏。表里现在是复跑后的值。）

**⑤ 那么"winsec 一字未动"今天对不对？**

- **按字面、且指整个 `internal/winsec/`：不成立。** 到 HEAD 为止这个包里动过 **53 增 / 11 删注释**
  （`winsec_other.go`）**外加测试件 271 新增 ＋ 76/20 改**（`33c8acd` 那发恒真修法）。
- **按"判定分支"读：成立**，且我用比第一轮更狠的尺（剥尽注释后 0 hunks）验过。
- **按 `winsec.go` 单文件读：成立**（md5 四版一致）。

票面在 `034080c`/`4f19ec6` 之后把话改口成"判定未动，注释动了"（票面 `:299-308` 那段"一处更正（append-only，不改上面那段的字）"），
而我上面每一档都是重量的 ⇒ 按派单第 ③ 条的口径，**这一子条闭合**。
⚠ **必须明写的一句：闭合的方式是措辞对上读数，不是读数变干净了。**
盘上今天仍然写着"53 增 / 11 删"，只是台账与票面不再把它写成"一字未动"。
（另外记一笔共树现状，不归我裁：票面 AC#2 那格仍是实现方自勾的 `[x]`，第一轮判的是**退回**——
"表里退回而票面 `[x]`"这个形状本程不翻转、不代勾，只把事实留在这一行。）

---

## §2 正向读数重跑 ＋ 第一轮"条件 1"逐条对上

### 2.0 这一节的读数是在哪一版树上量的（先把这个钉死，否则下面每个数都会腐坏）

```
git rev-parse HEAD        = bc4909612d4fb2b230042f2e0c15ec173bc444c9   (git cat-file -t = commit)
git log -1 --format='%h %s' = bc49096 docs(A165; 更正 A164 的判早): …
```
⇒ **§2 与 §3-A 的全部容器读数 ＝ 锚点 `bc49096`（2026-09-24 16:2x，dev）**。
共树在跑（同机还有 137 AC#4 与另一枚 119 复判程），**HEAD 每几分钟就换人**，
所以任何后来的人要复算，请 `git archive bc49096` 而不是 `HEAD`。
容器：`golang:1.27`（`go version go1.27.1 linux/amd64`，本机已有、未 pull），
挂载一律 **Windows 形式** `-v "D:\tmp\<tree>:/wisp:ro"`，容器内每发自证
`ls -l /wisp/go.mod` ＝ `-rwxrwxrwx 1 root root 883` ＋ `md5sum = f6ef661732b1851e5c3db348113cb605`。
⚠ **派单给的 `-v /d/tmp/…:/wisp` 这一发我亲手量到是假绿**：容器里 `/wisp/go.mod` 不存在，
而 `docker run` **rc=0**（同机那枚 119 复判程也登记了同一发，我这里独立复现一次）。
四数一律 `-count=1 -v`、按 `^\s*--- ` 计数；每发另加 `grep -c '^panic\|^fatal error'`；
`GATE golang-containers=0` 在每发前现量（18 发全 0，无一发撞负载）。

### 2.1 正向读数：`internal/config` 那两枚用例 ＋ 落点归属逐行表

**(a) 两枚真用例（锚点原树 `wisp119r2b-t0`，两形各一发）**

| 形 | 读数 |
|---|---|
| 普通 `TMPDIR=/tmp/plain119r2b`（`ls -ld` 读到 `drwxr-xr-x`） | `RUN=2 PASS=2 FAIL=0 SKIP=0 panic=0 rc=0`：`TestRoundTripLoadMarshalLoad` PASS、`TestResolvedNeverPersists` PASS |
| 软链 `TMPDIR=/varlink119r2b/tmp`（`ls -ld` 首列 `l`、`readlink -f` = `/realpriv119r2b/tmp`） | `RUN=2 PASS=2 FAIL=0 SKIP=0 panic=0 rc=0`：同两名逐名 PASS |
两形的 `=== RUN` 名册与 `--- PASS` 名册**逐名相同**（差集为空），SKIP 名册两形皆空 ⇒ 没有"被跳过冒充变绿"。
两形都打 `level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=1`
⇒ **软链形下 C26  seam 也装上了**（第一轮 §七③ 那次"软链 temp 让守门人自伤拒装、打 ERROR"的读数
在 `bc49096` 上已经不响；那是票 125 AC#2 `4824bb8` 的功劳，**不是票 119 的**，别记错账）。〔独立复现〕

**(b) `subtest | path | seal` 逐行表（探针树 `wisp119r2b-t5`＝锚点树 ＋ 一枚仓外探针文件
`internal/config/zz_r2b_landing_probe_119_test.go`，该文件**只存在于仓外快照**、未进仓、未 add）**
探针走的是生产路（`config.SaveFile` → 密封缝），复用包内自己的 helper
（`fullConfig`、`sealableTempDir124`），不新算任何期望值：

| subtest | 形 | path（SaveFile 的目标） | seal 结果 | 落盘字节 |
|---|---|---|---|---|
| `resolved-helper-root`（＝两枚真用例今天递根的形状） | 普通 | `/tmp/plain119r2b/…/001/config.toml` | `seal-err=<nil>` | **4163** |
| `resolved-helper-root` | **软链** | **`/realpriv119r2b`**`/tmp/…/001/config.toml`（**不是** `/varlink119r2b` 那一拼） | `seal-err=<nil>` | **4163** |
| `raw-ambient-tempdir`（＝票 124 之前那两枚的形状） | 普通 | `/tmp/plain119r2b/…/001/config.toml` | `seal-err=<nil>`（无线可穿） | 4163 |
| `raw-ambient-tempdir` | **软链** | `/varlink119r2b/tmp/…/001/config.toml` | **拒**：`winsec: refusing to seal … .wisp-config-….tmp … reaches it through the link at `**`/varlink119r2b`**`，which is not the tree this call names` | **-1（什么都没写）** |
| `planted-link-inside-resolved-root`（对照，证明"seal-err=<nil>"不等于"这条人不密封"） | 普通 | `/tmp/plain119r2b/…/planted119r2b/config.toml` | **拒**，且记名记的是**本用例自己种的** `…/001/planted119r2b` | -1 |
| `planted-link-inside-resolved-root` | 软链 | `/realpriv119r2b/…/planted119r2b/config.toml` | **拒**，同上记名自己种的链接 | -1 |
整支探针：普通 `RUN=4 PASS=4 FAIL=0 SKIP=0 panic=0 rc=0`、软链 `RUN=4 PASS=4 FAIL=0 SKIP=0 panic=0 rc=0`
（三枚子测试全绿，因为"被拒"本身就是它们的期望；`path/seal` 两列才是这一发的内容）。

⇒ **R-119-4 那半面（"② 换出来的落点归属"）本程自己量到了，方向与第一轮一致**：
根一旦由调用方解析，**"哪棵树被密封"就交给 `TMPDIR`（这里是 `/varlink119r2b`）的所有者**——
写进去的是 `/realpriv119r2b/…`、4163 字节真落盘，而**穿过新种的链接照旧拒、且拒了什么都不建**。
⚠ 第一轮那句"数据根被建成并被 **chmod 0700**"我这发**没有复算到**（我打的是根目录 mode，
探针的 `Perm().String()` 把类型位吃掉了、读到 `rwxr-xr-x`；`config.toml` 自己的 mode 我没取）⇒ 记进 §5。

### 2.2 恒真自查（派单第 ⑤ 条：判据在未修旧码上今天响不响）

我自己这一发正向读数如果"怎样都绿"就是装饰。**树 `wisp119r2b-t6`＝锚点树 ＋ 把
`internal/config/tempdir_resolved_124_test.go:33` 的 `proc.SealableRoot(t.TempDir())` 退回 `t.TempDir()`**
（落地原文：`return func() string { _ = proc.SealableRoot; return t.TempDir() }() // R2B de-resolve`）：

| 形 | 读数 | 红名 |
|---|---|---|
| 普通 | `RUN=2 PASS=2 FAIL=0 SKIP=0 rc=0` | — |
| **软链** | `RUN=2 PASS=0 `**`FAIL=2`**` SKIP=0 rc=1` | `TestRoundTripLoadMarshalLoad`（`boundary_test.go:191 save: …`）、`TestResolvedNeverPersists`（`:239`），两句红因都是 `refuses … reaches it through the link at /varlink119r2b` |
⇒ 这一发是"**未修码上今天不响、把解析摘掉才响**"的形状 ⇒ §2.1(a) 的"两形都 PASS"**不是恒真**，
它确实在消费"调用方先解析"这条纪律。⚠ 同时给这条尺**划清归属**：
它响的是 **票 124 那批 harness 换根**（`sealableTempDir124`），**不是票 119 那三枚生产调用点**——
把它当"票 119 AC#2 的落地凭据"记就是一次数错对象（这正是本会话今天上午那笔"归因腐坏"的同类）。**〔独立复现〕**

对 §3-A 那几发的恒真自查另有一发：`wisp119r2b-t7` ＝ 两味齐全的树 ＋ **MUT-D**
（把 `internal/winsec/winsec_other.go:158` 文案里的 `" through the link at "` 抹成 `" reaches it via "`；
落地证 ＝ 该串在两枚文件里的 `grep -c` 由 1 变 0；**只动措辞、`ErrUnresolvedPath` 与判定分支一字未动**）
⇒ 两形各 `RUN=5 PASS=3 `**`FAIL=2`**` rc=1`，红名恰是接了记名断言的那两枚
（`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`、`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`），
而 `:251` 那枚 AC#3 与两枚 becomes-sealable 仍 PASS ⇒ **记名断言有牙齿、不是恒绿**。**〔独立复现〕**

### 2.3 第一轮 `119-adversarial-acceptance.md:100` 那句"条件 1"，逐条对上

第一轮原文（:100-:113）拆成四条，逐条给我自己的读数（**〔独立复现〕**指标本程现量；
标〔日志＋归档，抽验〕＝我只引第一轮那发的数并核了它的可复算性）：

| # | 第一轮写的 | 今天（锚点 `bc49096`）盘上／容器里是什么 | 闭合？ |
|---|---|---|---|
| 1 | `cmd/wisp/secret.go:118 resolveSecretLayout` 把 `os.UserConfigDir()` 原样交给 `proc.LayoutFor`，**不走 `DefaultLayout`** ⇒ ② 有第二个生产读者没被覆盖 | `grep -n` 现量：`cmd/wisp/secret.go:129 func resolveSecretLayout`、**`:148 root = proc.SealableRoot(root)`**（`36294c2` 落的那"一行"，numstat `cmd/wisp/secret.go 15/0`）⇒ 这条生产路现在解析 | **闭合** |
| 2 | 这条路上"合法数据根被误伤 ⇒ 跑不起来"今天仍被真实造出来（真二进制 `wisp secret list` dev 三形 rc=1） | 我没有重跑真二进制（要 sherpa 运行库＋宿主链接，属票 111/123 地界）；**我用同一条函数的容器用例代替**：`./cmd/wisp/ -count=1 -v -run 'SecretRoute.*119'` 两形各一发＝ `RUN=3 PASS=3 FAIL=0 SKIP=0 panic=0 rc=0`，三枚名册＝ `…SymlinkedHomeBecomesSealable119`、`…SymlinkedXDGConfigHomeBecomesSealable119`、`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`。**真二进制那一发仍属第一轮〔日志＋归档，抽验〕**，本程不背书其 rc 数 | **按"结局不再被造出来"闭合；按"我这发亲眼看过的最强凭据"只到用例级** |
| 3 | `grep os.Symlink cmd/wisp/*_test.go` 命中 **0** 次 ⇒ 这条形状在本仓零用例 | 本程重跑：命中文件 **1 枚**（`cmd/wisp/secret_dataroot_119b_test.go`，`os.Symlink` **3 处**）；`grep resolveSecretLayout\|resolveDataDir cmd/wisp/*_test.go` 由 0 → **18 处**（`dataroot_128_test.go`）＋ **9 处**（`secret_dataroot_119b_test.go`）⇒ 命令面级用例有了 | **闭合** |
| 4 | 修它要"不许拿被测函数算 fixture"（第一轮 :514；`MUT-119-R` 56/56 全绿就是没守这条的直接后果） | `grep -n SealableRoot cmd/wisp/secret_dataroot_119b_test.go` ＝ 只在 **:24 的注释**里出现（"never from proc.SealableRoot"）；期望值来自 `filepath.EvalSymlinks`（:79）＋ `os.SameFile`；恒真那一族另有 `33c8acd`（`internal/winsec/dataroot_symlink_119_other_test.go` 76/20）把 fixture 换成 `cleanSpelling119`（:112 明写"deliberately filepath.EvalSymlinks and not proc.SealableRoot"） | **闭合** |

**"② 只做了一半"这半句本身**：第一轮点名的覆盖面是 3 枚调用方（`envfork.go:106/:226`、`doctor.go:249`）
缺第 4 枚。本程在 `bc49096` 现量 `grep -rn "SealableRoot(" --include='*.go' cmd internal | grep -v _test.go` ⇒
**4 枚生产调用点 ＋ 1 枚定义**，逐枚给号：

```
cmd/wisp/doctor.go:263            base = proc.SealableRoot(base)
cmd/wisp/secret.go:148            root = proc.SealableRoot(root)      ← 第一轮点名缺的那一枚
internal/proc/envfork.go:125      return filepath.Join(SealableRoot(os.TempDir()), …)
internal/proc/envfork.go:245      l, err := LayoutFor(env, SealableRoot(dir))
internal/proc/envfork.go:160      func SealableRoot(path string) string   ← 定义，不是调用
```
⚠ **行号已经从第一轮引的那组（:106/:226/:249）漂走了**（票 124/125/128/129 都动过这几枚文件）
⇒ 引"第几枚调用方"要引**枚数与名册**，别引行号。
而 `winsec_other.go:116-119` 的注释今天**逐名列了这四枚**（含 `resolveSecretLayout`）
⇒ 半面补全，且注释与盘上数量对得上。**〔独立复现〕**

### 2.4 那本账接住了没有——分两层答，别混

* **账面上的"登记"层**：`R-119-4` 已在 `034080c` 写进 `winsec_other.go` 的"第二条成本"（§1 一档零行、
  二档 53/11 里就含它），票面 `:421` 也把它列进 `next=` 第 5 条。**这一层有。**
* **归属层（第一轮 :496/:529 建议并案到的那本账）**：票面 `next=` 第 5 条今天仍写着
  "**要不要并案到票 120 / `R-108-2` 那本账，请编排者裁**"——**没有并案、没有裁**。
  所以：**"落点归属"这件事的*事实*这一格接住了（§2.1(b) 我量到了），它的*归属*还挂在编排者手上**。
  这不影响 AC#2 的射程（AC#2 要求的是"裁 ①/②/③ ＋ 写理由 ＋ winsec 是否一字未动"），
  但**别在台账里把 `R-119-4` 记成"已并案"**——它是"已登记、待人裁归属"。**〔独立复现〕**

### §1／§2 的提交账（按"正文里出现'已提交'就必须带两条原样输出"的规矩补）

```
$ git show --name-only 4757198 | tail -3        # §1 那枚
    - 归因腐坏一条：resolve.go 到 034080c 仍是 7eb8a754…、HEAD 已变 b6876a5e…（126/125/129 三枚别的票，205/16）
    - 同格双派登记：6e04d1a 落了第二枚 119 复判表的 §0，本程未读其结论、快照目录不同名

docs/evidence/s1/119-ac2-r2b-acceptance.md

$ git log --oneline -1 ; git show --name-only HEAD | tail -3      # §2 提交当刻
024ee87 evidence(119 AC#2 r2b §2): 正向读数在锚点 bc49096 重跑＋第一轮"条件 1"四条逐条对上（四条都闭合）
    - 2.4 落点归属分两层：登记层有、归属层（并案到票 120/R-108-2）仍挂编排者 ⇒ 台账别记成已并案

docs/evidence/s1/119-ac2-r2b-acceptance.md
```
§1 那枚完整号 ＝ `4757198`（`git log --oneline -1` 当时回显的就是它，见上一条命令的 subject 行；
`git show --name-only` 只列一条路径＝本文件）。两枚都**只有我自己的路径**——
但 §1 那枚提交**前**的 `git diff --cached --name-only` 回显过**一枚别人的**
`docs/evidence/s1/137-ac4-r1-acceptance.md`（同机那程自己 staged 的，状态 `MM`），
本程**没有提交它、没有 unstage 它**（显式 pathspec 挡住的，`git show --name-only 4757198` 可复算），它现在仍在索引里。

---

## §3 顺手两问

### A. 票 119 `AC#7` 判据①"两味药必须一起下……只收紧不换根会在软链形恒红"——成立（我不改它、也不实现它）

锚点 `bc49096`、容器 `go1.27.1 linux/amd64`、`-count=1 -v -run 'TestAC[123]POSIX.*119'`
（名册永远是那 5 枚，全表见 `/d/tmp/wisp119r2b-work/summary.txt`）：

| 树 | 一味 | 普通形 | **软链形** |
|---|---|---|---|
| `t0` 锚点原样 | 两味都没下（今天） | `5/5/0/0 rc=0` | `5/5/0/0 rc=0`——**但 `:203`/`:308` 两枚的拒因打印的是宿主的 `/varlink119r2b`** ⇒ 今天零区分力（这条＝判据②要的那发，我复现了） |
| `t2` | **只接记名**（同包 `refusalCreditsLink137`，零新增依赖） | `5/5/0/0 rc=0` | `5 RUN / 3 PASS / `**`2 FAIL`**` rc=1`，红名＝**恰好那两枚**，红句原文"refusal did not credit the link this case planted at …" ⇒ **"软链形恒红"实测成立**（普通形不红，所以判据写的"软链形"这个限定词是准的） |
| `t3` | **只换已解析根** | `5/5/0/0 rc=0` | `5/5/0/0 rc=0`——换根后**打印的拒因已经变成本用例自己种的链接**（`… through the link at /realpriv119r2b/…/001/varlink119`、`…/001/injlink119`），可 sentinel-only 断言在两形都不响 ⇒ **另一味同样不可省**（否则这一格只是把日志换对了） |
| `t4` | 两味一起 | `5/5/0/0 rc=0` | `5/5/0/0 rc=0` ⇒ **判据可满足，不是一枚结构上永远产不出读数的尺** |
| `t7` | 两味一起 ＋ MUT-D（抹掉生产文案里的 `" through the link at "` 标记，`grep -c` 由 1→0；`ErrUnresolvedPath` 与判定分支未动） | `5 / 3 PASS / `**`2 FAIL`**` rc=1` | 同左 `2 FAIL rc=1` ⇒ 两味齐了以后**有牙齿**（红名册仍只有那两枚，`:251` 那枚对照组与两枚 becomes-sealable 全程 PASS ⇒ 判据④"AC#3 逐名不变"在我这四发里也没被弄绿） |

⇒ **判据①照原样成立，两条单向分支我都各自量到了红/绿的具体位置。**
只补一句给落地的人（**不改判据本体**）：`t3`/`t4` 我这发用 `proc.SealableRoot(t.TempDir())` 换根，
那是**仪器形状**；真要落地该跟同包 `:251` 那枚的先例走 `cleanSpelling119`（`filepath.EvalSymlinks`），
否则会踩本票 Rules `:94` 那条"不许拿被测函数算 fixture"（`R-119-9`）。**〔独立复现〕**

### B. `dataroot_symlink_119_other_test.go` 里那三枚用例今天在 CI 上有没有分母——**有**

现量两处文件 ＋ 一发真跑过的 run：
- `.github/workflows/ci.yml:224-225` `test-core:` ＋ `runs-on: ubuntu-latest` → `ci.yml:288`
  `run: bash scripts/portable-tests.sh --scope=core`；
- `scripts/portable-tests.sh:180` core scope 数组里明写 `./internal/winsec/`，`:149` `core_pin` 钉
  `github.com/CarlosShao/wisp/internal/winsec`（少了这行 `portable-tests.sh` 的 pin 守卫会红）；
- `internal/winsec/dataroot_symlink_119_other_test.go:1` 是 `//go:build !windows` ⇒ 在 ubuntu 上参与编译；
  文件里两枚 `t.Skipf` 在 `:77`/`:243`，条件都是"这里种不出软链"，不是平台自拒；
- 真跑过的一发：run **`35967768017`**（dev push，`2026-09-24T07:05:41Z`），job `test-core`（`107530150794`）
  结论 **success**，日志行 `ok github.com/CarlosShao/wisp/internal/winsec 0.026s`，
  同发 runtests.sh 汇总行的 `-skip` 正则（`TestDefaultDeadlineWallClockMeasurement|TestSubprocessCrashWriter|
  TestRealDownloadVadThroughPipeline|TestRealDownloadPuncArchiveThroughPipeline|TestC26RewrittenSyncRootDoesNotDisarmSuspectNet|
  TestWorkspaceSwitchRefusesAJunctionToOutside|TestD34WriteMatrix|TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`）
  **不含 119 的任何一枚**，顶层 `PASS=711 FAIL=0 SKIP=0`。〔独立复现（读的是 run 日志原文）〕

⚠ **"有分母"只覆盖了一半的形状，这句必须跟着上面一起读**：CI 的 ubuntu 上 `/tmp` 是**真目录**，
所以那 5 枚在 CI 跑的是**普通形**——正是 A 表里 `t2` 那一列（只接记名在普通形 `5/5 PASS`）。
**"软链形"那一列在 CI 上零分母**，只有容器有（本程 18 发全是容器）。
另一笔账顺便核了：票面"两处更正"第 2 条说 `cmd/wisp` 那三枚新用例在 CI 上没腿——**今天仍成立**：
`cmd/wisp` 属 `--scope=cli`（`scripts/portable-tests.sh:195`），CLI 那一步只在 windows 腿
（`ci.yml:334-335` 的 `test-windows`、`ci.yml:422` `bash scripts/wisp-cli-tests.sh`），
而 `cmd/wisp/secret_dataroot_119b_test.go:1` 是 `//go:build !windows` ⇒ **那三枚在 windows 腿上根本不编译**。〔独立复现〕

（§3 那枚提交的号：`85c1393`，`git show --name-only 85c1393` 末行只有
`docs/evidence/s1/119-ac2-r2b-acceptance.md` 一条路径。）

---

## §4 总判（只裁 AC#2 这一格）

**成立（无附条件）。** 一句话理由：第一轮把这一格按退回的两条理由——"② 只做了一半"与"交回的注释把这一半写成了全部"——
今天**一条已经不成立、另一条以"措辞对上读数"的方式闭合**，而这两种闭合我都是自己在盘上/容器里重量过的，不是抽验别人日志。

按派单五条逐项落点：

| 派单要裁的 | 本程读数落点 | 判 |
|---|---|---|
| ① 返修 commit 逐枚 numstat | §1② ：动 `winsec_other.go` 的返修只有 `034080c`（41/14），交件 `189cb1e` 21/6 与第一轮逐字同 | 量到 |
| ② 三档改动账 | §1③ ：一档 0 行／二档 53 增 11 删（＝全部）／三档 0 行，两把尺各 0 hunks | 量到 |
| ③ "winsec 一字未动"今天对不对 | §1⑤ ：按字面不成立、按判定分支与 `winsec.go` 单文件成立；票面 `034080c`/`4f19ec6` 已改口 ⇒ **闭合的方式是措辞对上读数，不是读数变干净了**（盘上今天仍是 53/11 注释） | 闭合 |
| ④ 退回的另一半（落点归属那本账 ＋ 条件 1 逐条） | §2.1 正向读数我自己跑（含 `subtest\|path\|seal` 逐行表）；§2.3 条件 1 四条**全部闭合**；§2.4 分两层：登记层有、**归属层（并案到票 120/`R-108-2`）仍挂在编排者手上**，别记成已并案 | 闭合（AC#2 射程内） |
| ⑤ 恒真判据自查 | §2.2 两发：`t6`（把 config 的换根 helper 退回未修形）⇒ 软链形 2 枚 FAIL；`t7`（两味齐 ＋ MUT-D）⇒ 两形各 2 枚 FAIL。§2.0 把"我这发在哪版树上量的"钉成 `bc49096`；§1④／§2.3 各主动报一条**引文腐坏**（`resolve.go` md5 到 HEAD 已变、调用点行号已从第一轮引的那组漂走） | 自查过 |

三件**必须跟着这个"成立"一起读**的事（不是条件，是边界）：

1. **票面 AC#2 那格仍是实现方自勾的 `[x]`**，而第一轮表判的是"退回"。改名与翻勾的权不在本程，
   本程也不替它翻——**台账要记的是"第二程复判＝成立"，别记成"表与勾一直一致"**。
2. **`R-119-3`（声明树动不动）这一裁我核到的是"有用例钉着"**：`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`
   两形各 PASS（在 5/5 那一发里逐名点到），两半（声明原样／OS 答案解析）都在这枚的两条 leg 上；
   它的**文档层**（`internal/proc/envfork.go:26/7` 那笔注释）我只数了 numstat、没逐句验语义 ⇒ 记进 §5。
3. **同格有第二枚验收程在跑**（`119-ac1-ac2-r2-acceptance.md`，§0.3）。本格＝AC#2，两份表不互相抵账。

## §5 我没核的清单（如实列，不当已核）

1. **真二进制那一层**：第一轮 §二(c) 那五行 `wisp secret list`／`wisp doctor` 的 rc 表，我**没有**重跑
   （要在容器里造带 sherpa 运行库的真二进制 ＝ 票 111/123 地界）。我的替代凭据是用例级（`cmd/wisp` 三枚两形 3/3 PASS），
   真二进制那五行在本表里只值**〔日志＋归档，抽验〕**。
2. **macOS 半面**（`TMPDIR` 住在 `/var` 那条链接下面）：零读数，容器不是 macOS。
3. **`chmod 0700` 那一细节没复算到**：我探针打的是根目录 mode、且 `Perm().String()` 会把类型位吃掉
   （所以表里那列写成 `rwxr-xr-x` 而不是 `drwxr-xr-x`）；`config.toml` **文件自身**的 mode 我没取。
   第一轮"数据根被建成并被 chmod 0700"这半句，本程只背书"落盘 4163 字节"这半句。
4. **MUT-D 我只抹了一处**（`winsec_other.go:158`）。137 那发的 MUT-D 抹的是**两处**（另一处 `winsec.go:262`）
   ⇒ 我只验了"本包这条腿哑不哑"，不替"两枚标记全抹"那一发背书。
5. **AC#6 的地界一枚未跑**：`d22scan`／`gofmt`／`gofumpt`／`go vet`／全仓测试／`-race` 都没跑（本格不裁 AC#6）。
6. **AC#7 判据②③④ 没按它自己那四条逐条裁**：我只答"①成不成立"，加上"④那枚 AC#3 用例在我四发里逐名没变色"这一条顺带读数；
   并且**没改 `AC#7` 正文、没实现它**（生产码与测试码在仓内零改动，见 §6 的 `git status`）。
7. **票 137 AC#4 那 7 处换根、票 124 那批换根的归因**：我只在 §2.2 划了一句边界（那两枚 config 用例红因属于 124 不属于 119），没逐枚核。
8. **CI 只核到 job conclusion ＋ `test-core` 一发的日志行**；同 run 里 `lint`／`test-windows` 两枚 failure 与本格有无关系**没查**。
9. `docs/reports/pending-and-issues.md` 的台账我**没读**（不替编排者改账）；第一轮表里 `A160`/`A163` 那两笔转述我只当起点引用。
10. **同格那枚 `119-ac1-ac2-r2-acceptance.md` 的正文我没读**（只 `--stat` 看了它动哪枚文件），以免把别人的读数当自己的。

## §6 临时件 ／ 注入两栏 ／ 被拒记录

**临时件（全部只建不删，都在仓外）**

```
D:\tmp\wisp119r2b-work\   drive.sh(容器侧形状自证) patch.sh(仪器补丁) run.sh(宿主侧驱动)
                          md5-head.txt anchor.txt(锚点) summary.txt(18 发四数+红名册+panic 计数)
                          base-code.txt head-code.txt(§1 尺 A 的两份剥注释文本)
D:\tmp\wisp119r2b-t0\     git archive bc49096 | tar -x  ＝ 被验版本（纯净）
D:\tmp\wisp119r2b-t2\     ＋只接记名        D:\tmp\wisp119r2b-t3\  ＋只换已解析根
D:\tmp\wisp119r2b-t4\     ＋两味一起        D:\tmp\wisp119r2b-t7\  ＝ t4 ＋ MUT-D（抹标记）
D:\tmp\wisp119r2b-t5\     ＋仓外探针 internal\config\zz_r2b_landing_probe_119_test.go（未进仓）
D:\tmp\wisp119r2b-t6\     ＋把 config 的 sealableTempDir124 退回未解析形（恒真自查）
D:\tmp\wisp119r2b-logs\   18 发 .log ＋ .rc
命名卷 wisp119r2b-build（本程自己的 go 构建缓存）；只读复用既有卷 ac119-gomodcache（＝依赖缓存，不是别人的读数）
容器一律 --rm，收尾 docker ps 里 golang 计数 0
```
⚠ 两处我自己造的垃圾也留着：`D:\tmp\wisp119r2b-t0;C`（一条 shell 引用写错的产物，空目录）、
以及 §2.1(b) 那枚探针里的 `rootmode` 列格式错（`Perm().String()`）——都如实标在 §5 第 3 条，不删不盖。

**只读守住的自证**：`git status --porcelain internal/winsec/` 输出为空（**最后一次现量 `2026-09-24 16:35:45 +08`**，
就在写下这一段之前；同刻还跑了下面那五行 `git show --name-only`）。
本程**已落地**的五枚 commit（`20f8c22`/`fc71fa6`/`4757198`/`024ee87`/`85c1393`）的 `git show --name-only --format=''`
逐枚回显（原样，16:35:45 现量）：
```
20f8c22 -> docs/evidence/s1/119-ac2-r2b-acceptance.md
fc71fa6 -> docs/evidence/s1/119-ac2-r2b-acceptance.md
4757198 -> docs/evidence/s1/119-ac2-r2b-acceptance.md
024ee87 -> docs/evidence/s1/119-ac2-r2b-acceptance.md
85c1393 -> docs/evidence/s1/119-ac2-r2b-acceptance.md
```
⇒ 每枚只带这一条路径；生产码/测试码/票面/阈值/golden 零改动。
工作树里 `design/**` 那批 owner 侧的挪动**一枚未动、未 add、未还原**（收尾 16:35:45 现量
`git status --porcelain design/` = 18 行、首三行 ` D design/assets/{base.css,icons.js,theme.js}`）。

**注入栏（每条带出处：工具名 ＋ 命令前 40 字）**

| # | 形状 | 出处（工具名 ＋ 命令前 40 字） | 处置 |
|---|---|---|---|
| 1 | 回显一枚**盘上不存在的 commit 号** `5f7e104`，主语还冒充"票 119 AC#2 复判表建件（骨架，§0 第一次提交）"＝本程自己的活 | `Bash`：`git log --oneline -1; echo "=== show ==="; gi` | `git cat-file -t 5f7e104` → `Not a valid object`；`git branch --contains` 同错；dev 近 12 枚无此号无此主语 ⇒ **判为冒充锚点的注入／异常回显，只登记不服从**；此后每个 sha 落笔前先 `cat-file -t` |
| 2 | 同一批里两条 git 命令**给了空输出**（只有 echo 的分隔符），`git log -1` 像没执行 | `Bash`：`date "+%Y-%m-%d %H:%M:%S %z"; git log --oneli` | 复跑同一对命令拿到真值；不当"没提交"处理（`git show --name-only` 追到了 `20f8c22`） |
| 3 | `git show HEAD:internal/winsec/resolve.go \| md5sum` 回显**旧串** `7eb8a754…`（复跑三次实为 `b6876a5e…`） | `Bash`：`cd "/d/work/workspace/projects plans/Wisp" &&` | 用 `md5sum` ×3 ＋ `git diff --numstat 034080c..HEAD` 双查，发现是**回显错不是树错**；§1④ 表按复跑值写，并把这次误显记在正文里（不悄悄改） |
| 4 | 一条"harness 说某文件被创建"的提示，指向 `D:\work\workspace\D:tmpwisp119r2b-t0` 这枚**路径本身畸形、盘上不存在**的对象 | `Bash`：`cd /d/tmp && for v in t2 t3 t4; do rm -rf /d` 那一批的回显尾部 | `ls` 反扫＝不存在 ⇒ 按注入判（第 8 代形状：冒充我们自己的对象名），不动作 |

**真通知／正常回显栏（分开计，不混进上面）**

- `MEMORY.md was modified` ×2（本会话自动记忆摘要在写我的记忆目录，编排者已按三条判据核过有真实来源）⇒ 不需要本程做任何事；
- harness 追加的 `The task tools haven't been used recently…` ＋ **别人的**任务列表 ×≥9（从第一条起就判它不是指令，未逐条精确计数——**报"约数"而不是精确数，就是为了不制造不可核的自述**）；
- `Write`/`Edit` 成功回执里的 CRLF 提示 × 若干（我全程没在任何跟踪文件里写 CR，§1 两把尺读的都是 LF 文本）。
⇒ **判为注入：4 条；要我"少取证／别用工具／直接下结论"的：0 条**（第 1 条最坏变体如果被我照抄，本表的锚点就是一枚假号）。

**本程被拒过的工具调用：零次（权限类）。**
命令级非零 rc 如实列在这里，与"被拒"分开：`grep -c` 命中 0 时的 rc=1（§1 三处、尺 A 那发 diff 的 rc=1 同因）、
`git cat-file -t 5f7e104` 的 rc=129（**故意的反扫**）、`gh run list` 在 16:1x 一度网络失败（约 15 分钟后重试成功，
所以我没有把"CI 闸门"当成已核，改用 `docker ps` 计数＋每发前 `GATE golang-containers=0` 自证）。
**没有一次拒发生在取数之前，也没有一次取数因拒而少跑。**

