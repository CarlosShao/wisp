# 119-`next=` 三问审计 r1 —— 只读取证（`auditor-ticket119-next`）

**被派**：台账 `docs/reports/pending-and-issues.md:5373`（A170 节内，09-24 16:4x）——
"替我把 `next=` 第②条 `R-119-5` 到底结没结、第⑤条 `R-119-4` 该不该并案、第⑥条那两枚 staged deletion
的归属三条量清楚，只新建 `119-next-audit-r1.md`"。
**角色**：本程**只读审计**。全程 `go test`／容器 **0 次**、`internal/winsec/**` **0 字节改动**、
commit **0 枚**、push **0 次**；唯一新建文件＝本文件。

## §0 锚点自量（本程所有"盘上原文"读数的取数时刻与版本）

```
$ git rev-parse HEAD
b1ea719e494ed3d62624b5d00d6fc77cc0ca60d1
$ git branch --show-current
dev
$ date '+%Y-%m-%d %H:%M %z'
2026-09-24 17:00 +0800
$ git status --porcelain
 D design/assets/base.css …（design/** 共 16 枚未暂存删除，owner 自己挪的已结案账，本程一字未碰）
?? design/doubao/
?? design/old/
 M internal/winsec/dataroot_symlink_119_other_test.go   ← 取数中段现量（17:00），＝ `worker-ticket119-ac7` 在飞的那一枚
```

- 本程引用的每一枚源码原文都先证明**工作树＝HEAD**：
  `git diff HEAD --stat -- internal/winsec/resolve.go internal/config/c26_seam_posix_125_test.go` 现量**空输出**（0 行）。
- ⚠ 取数中段（17:00）盘上出现一枚**不是我下的**未暂存改动：`internal/winsec/dataroot_symlink_119_other_test.go`
  （派单 A170:5372 点名的在飞写码位）。本程全部结论**不依赖**该文件的盘上版本；开工那一刻（第一次 `git status`）
  它还是干净的。此处登记只为说明"共享树在动"，本程没有读它、没有改它。
- 本程引用的每一枚别程读数都**连它量数的树锚点一起引**（这是本项目今天刚吃过亏的那条规矩）。

---

## §Q1 `R-119-5`（`next=` 第 2 条）：两半各自今天在不在

票 119 原文（`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md:415-418`）：

> 2. `R-119-5`（POSIX 上"C26 解析器在不在位"零用例 + 软链 temp 下守门人自伤拒装）**今天仍在**：
>    本轮真二进制复算里，`WISP_ENV=test TMPDIR=/varlink/w119tmp` 那一形**仍打那条 ERROR**
>    （`refusing to install a path resolver into the sealing seam … /varlink/wisp-103-conformance-probe …`），
>    dev 三形则打 `INFO … probes_passed=1`。它不在本票射程（`resolve.go` 禁改、`internal/risk` 冻结）。

### 先钉死一个时间事实：上面那句"今天仍在"对 HEAD 已经是过期描述

- 该 `next=` 段 `git blame -L 411,426` 现量＝ `9d252f2`（CarlosShao，**09-22 17:21**）；
  它所属轮次的纯净快照门禁跑在 `git archive 034080c`（票面 `:355`；`034080c`＝09-22 17:11，`git merge-base --is-ancestor` 两枚皆 ANCESTOR）。
- 票 125 那味药落地于 `4824bb8`（`fix(winsec,125,AC#2)`，**09-22 20:58**，ANCESTOR of HEAD）。
  ⇒ 票 119 写"仍打 ERROR"时，药还没到；**它的读数在它的锚点（`034080c` 树，`resolve.go` 停在 `a505607` 时代、无 `resolverProbeRoot`）上是真的**，
  对 HEAD 则已被票 125 的交件与终裁表覆盖。两头都不算谁的错，但**裁 R-119-5 不能拿 :416 那行当今天的读数**。

### 半（a）"POSIX 上'C26 解析器在不在位'零用例"——**已不在（钉子今天在 HEAD 上）**

- 交件：`internal/config/c26_seam_posix_125_test.go`，`//go:build !windows`（盘上 `:1`），
  钉子 `TestAC1POSIXSeamHoldsC26Pipeline125`（`:54`，零 Skip 条件）＋ 第二形用例
  `TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`（`:117`，harness 自身经链接时 `:120` 才 skip——R-125-4 那一族，不是恒绿）。
- HEAD 含它：`git ls-tree --name-only HEAD internal/config/ | grep 125` ⇒ 命中（盘上也在）。
- **它与被验版逐字节相同**：`git log --oneline -- internal/config/c26_seam_posix_125_test.go` ＝ 仅 `4824bb8`/`ff3faf9`；
  `git diff 4824bb8 HEAD -- internal/config/c26_seam_posix_125_test.go internal/winsec/seam_probe_root_125_other_test.go` 现量**空输出**。
- 裁决（`docs/evidence/s1/125-ac1-ac4-r1-acceptance.md`，`acceptor-ticket125-r1`，非实现者）：AC#1 **成立**；
  其 §2.1（表 `:225-260`）是**本程自己跑的**三台对照——`pre`＝`git archive ff3faf9`（未修生产码）＋交件用例、`post`＝`git archive 4824bb8`——
  且其 §1 用自造 `MUT-R1A` 证明过钉子会红。原票面"今天唯一断言在 `resolve_windows_test.go`、POSIX 那枚是 Skip"
  （票 125 `:18`）是**建票时**的史实，被钉子文件自己头部注释（`c26_seam_posix_125_test.go:7-15`，锚 `81b4d5f`）逐字复述。

### 半（b）"软链 temp 下守门人自伤拒装"——药在 HEAD 的码里，效果读数只锚到 `4824bb8`，HEAD 级没人量过

**当前盘上（＝HEAD，diff 已证）`internal/winsec/resolve.go` 决定这一形结局的行**：

| 行 | 原文（节选） | 作用 |
|---|---|---|
| `:198` | `shapes := []string{resolverProbeRoot() + sep + ".." + sep + "wisp-103-conformance-probe"}` | 敌意探针**不再直接踩未解析 `os.TempDir()`**（修前旧形＝`os.TempDir()` 原样拼，见票 119 `:201` 与票 125 `:12` 引的 `:195/197`/`:258/260` 旧行号） |
| `:242-244` | `func resolverProbeRoot() string { return resolveProbeRoot(os.TempDir()) }` | "OS 给的答案先解析"的接点位 |
| `:253-279` | `func resolveProbeRoot(path string) string { …os.Lstat 走到最长存在前缀 → filepath.EvalSymlinks(cur) :273 → 未存在尾段原样接回 }`；失败方向 `:264`/`:275` `return path` | **拿不回实名就退回未解析原拼写＝今天的形状**，绝不静默跳过探针 |
| `:340-341` | `parent := resolverProbeRoot()` | 第二处探针根（树归属腿）同走解析 |
| `:158` | `slog.Error("winsec: refusing to install a path resolver into the sealing seam",` | ERROR 分支本身还在，只等**候选自己不诚实**时才响 |
| `:182` | `slog.Info("winsec: sealing path resolver installed", "resolver", name(r),` | 软链 temp 形今天的期望去处（`probes_passed=1`） |

**"效果达成"的原判读（终裁表自己量的，锚点全具名）**：

- 表 §0.3（`:57-60`）钉被验主版本＝**`4824bb8`**，pre＝**`ff3faf9`**（＝`4824bb8^`，快照 `grep -c resolverProbeRoot resolve.go`=0 证未修），
  门禁树＝**`bd50c63`** 且 §0.4（`:64-66`）证 `git diff --name-only 4824bb8 bd50c63 -- internal/ cmd/ scripts/ tools/` **空输出**（码逐字节＝`4824bb8`）。
- §2.1（`:225-260`）钉子四读数：软链形 `pre --- FAIL` ⇒ `post --- PASS`；现场行逐字
  `pre … AC#1 RED: winsec.PathResolverInstalled() = <nil>` ⇒ `post … INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1`；
  shapes 行 `pre ["<BASE>/varlink125/…"]` ⇒ `post ["<BASE>/real125/…"]`（敌意 `..` 一字未少）。
- **真二进制级**（正对着票 119 `:416-417` 那一问）：表 §3.1（`:420-436`，本程自己容器 build）
  `改前 ff3faf9`（19,092,792 B）`TMPDIR=<软链> wisp run "hi"` ⇒ stderr 第 1 行 `ERROR winsec: refusing to install …`；
  `改后 bd50c63（码＝4824bb8）`（19,097,696 B）同一形 ⇒ `INFO winsec: sealing path resolver installed … probes_passed=1`。
  ⇒ **"软链 temp 那一形不再打那条 ERROR"在被验锚点成立，且由非实现者复现。**
- 拒绝侧没放宽：§2.2 自造 `MUT-R1B`（摘掉解析刀⇒红）与 `MUT-R1C`（删候选答案底线复算⇒红）点名在案。

**从 `4824bb8` 到 HEAD 之间隔着一枚什么**：

- `git log --oneline -20 -- internal/winsec/resolve.go` 现量：`resolve.go` 在 `4824bb8` 之后**只被 `a45b2e9`**（票 129，09-22 23:41，ANCESTOR）动过。
- `git diff 4824bb8 a45b2e9 -- internal/winsec/resolve.go` 现量 **50 增/3 删、全在三个 hunk**：`sameTree` 加 `|| !sameAbsoluteness(a,b)`、
  新增 `sameAbsoluteness`（`IsAbs(a)==IsAbs(b)`）、`answerInsideTree` 加同一腿——注释自述方向"**only gets one more requirement to satisfy … only ever refuses more**"。
  **`:198/:242/:253/:340` 这条探针根链路一字未动**（diff 里零命中）。
- 最近的 HEAD 系旁证：票 137 终裁表（锚 `a9c4d58f99b134b0ab3253b80469dde92bb2692a`，09-24 15:47，ANCESTOR；winsec 码面与交付面 `b1010ff` 逐字节相同）
  §4 在 27 发容器跑里**普通形下 125 那三枚探针腿恒 0 SKIP**、且 `…SeamAcceptsTheHonestPOSIXAnswer125` 在普通形 MUT-D 下**转红抓破口**（`:406/:419-427`）
  ⇒ post-a45b2e9 树上这些腿活着且响。**但 137 跑的树里没有软链 temp × 真二进制的 ERROR-shape 那一发。**
- 终裁表自设边界（不许绕过）：票 125 `:42`"**`a45b2e9`（09-22 23:41）之后的 HEAD 版 `resolve.go`——本表不许被拿去签 HEAD**"；
  表 §6 第 3 条（`:658-661`）同文："**读者不许拿本表去签 HEAD**"。
- 授权面：动 `resolve.go` 那 82 行由台账 **A162**（`docs/reports/pending-and-issues.md:5190`，09-24 15:5x）今日署名**追认**（"今天生效"，非补写 09-22 批条），撤销口令「撤 A162 追认」。

**另注**：票 125 本身**还没翻 `-done`**（票头 Status: open；编排者 A160~A170 段明写"改名这件事留着"，卡在票 124 合流腿与"CI 哪条腿有分母"）——
但"票 125 未结案"≠"R-119-5 的药不在 HEAD"，两件事分开裁。

### Q1 裁定句

> **`R-119-5` 只能结一半。** 半（a）"POSIX 零用例"**可以今天就结**（凭据：`internal/config/c26_seam_posix_125_test.go`@HEAD 与被验版逐字节相同 ＋ `125-ac1-ac4-r1-acceptance.md` §1/§2.1 判 AC#1 成立，锚 `4824bb8`）；
> 半（b）"软链 temp 自伤拒装"**药已上 HEAD 的码**（盘上 `resolve.go:198/:242/:253/:340` 原文，本进程验工作树＝HEAD），
> "那一形不再打 ERROR"的**最强读数是终裁表自己容器 build 的真二进制 pre/post，锚 `ff3faf9`/`bd50c63`（码面＝`4824bb8`）**，
> 而 `4824bb8`→HEAD 之间唯一一笔 `resolve.go` 改动（`a45b2e9`）经我逐 hunk 静态核对**不触碰探针根**、方向只更严——
> **但终裁表具名拒绝签 HEAD，HEAD 级的软链真二进制读数今天没有任何一程量过**。
> 要整格结案，判法二选一：①你裁定"锚 `4824bb8` 的终裁读数＋本程 §Q1 的静态桥（命令原文齐）合并即足够"；②等票 119 AC#7 的在飞写码（`dataroot_symlink_119_other_test.go` 换根那批）落地后，随票 125/124 合流腿补一发 HEAD 树软链形真二进制复算再结。
> 无论选哪条，票面 `:415-418` 那句"**今天仍在**"本身已是过期描述（写于 09-22 17:21，药到 20:58），**不该再按字面续命**。

本节节档：半（a）〔独立复现（我给命令原文）〕；半（b）"码在 HEAD"＝〔独立复现〕、"效果"＝〔日志＋归档，抽验〕——
终裁表是**它自己容器里现量的**（非转述实现方），我核了它的全部锚点 sha、祖先关系（`4824bb8`/`ff3faf9`/`bd50c63`/`81b4d5f`/`a45b2e9`/`a9c4d58` 六枚
`git merge-base --is-ancestor <c> HEAD` 全 ANCESTOR，本程现量）与 `resolve.go` 的逐 hunk diff；HEAD 级效果读数缺失一事如实登记，**不背书成已有**。

复现命令（Linux 侧须容器，本程**未跑**，只给出票面原文的构造法；静态部分本机可跑）：

```
git log --oneline -20 -- internal/winsec/resolve.go
git diff 4824bb8 a45b2e9 -- internal/winsec/resolve.go
git diff 4824bb8 HEAD -- internal/config/c26_seam_posix_125_test.go internal/winsec/seam_probe_root_125_other_test.go
git merge-base --is-ancestor 4824bb8 HEAD && git merge-base --is-ancestor a45b2e9 HEAD && git merge-base --is-ancestor a9c4d58 HEAD
# 效果复算（容器）：终裁表 §2.1 构造法—— pre=git archive ff3faf9 + 4824bb8 的两枚交件用例； post=git archive 4824bb8
# HEAD 级缺口补测：  git archive HEAD + `ln -s <real> <link>`；TMPDIR=<link>/tmproot go test -count=1 -v ./internal/config/（看 AC#1 钉与 INFO 行）
```

---

## §Q2 `R-119-4`（`next=` 第 5 条）：要不要并案到票 120 / `R-108-2`

### 三处的台账/票面原句（逐字，含行号）

- **`R-108-2` 的账本体**（台账 `:1218`，A84②）：
  > **A84② `R-108-2` 我裁定了，而且裁的是"不改"。** 验收代理问：要不要把"外来"定义成"越出调用方数据根"。
  > **不要**——winsec 的守卫是"**拼写/祖先链/树归属**"这一级，"这棵树归谁"是**调用方数据根纪律**（票 76/95）的职责；
  > 把两层混到底线里会让守卫变成第二套策略而不是防止绕过的工具。做法：把它**写进文档**（`doc.go`），作为票 113 的 **AC#6 追加**

- `R-108-2` 答复落点的更正（台账 `:1552`，A85②）：
  > 票 113 的 AC#6（我给 `R-108-2` 的答复落点）写的是 `internal/winsec/doc.go` —— **这个文件在本仓不存在**，包文档注释在 `winsec.go` 第 1 行起。

- **票 120 的账本体**（台账 `:1272`，A89⑤）：
  > **票 120**＝同一族形状**关得不够严**（check-then-act，验收方 4000 次重试抓到 **154 次**，且这条没进它的"成本"注释）
  > …119/120 并成一张必然有一条判据被另一条盖掉；…**并票的边界按"要不要动生产码"划，不按主题划**。

- **`R-119-4` 的定义与建议归属**（`docs/evidence/s1/119-adversarial-acceptance.md:496`，§八表）：
  > | `R-119-4` | **② 换出来的一半面**：环境变量所有者现在能决定密封落在哪棵树（实测 rc=0 且数据根被建成并被 chmod 0700，修前 rc=1 什么都不建）。穿过链接仍拒（实测），所以这是**落点归属**问题，不是可穿性问题 | `internal/risk` C26 那本账（票 102/105 地界），或票 120 的"先查后动"族 |

  同表 §九行 5（`:529`）：
  > | 5 | `R-119-4`（软链把落点让给环境变量所有者）与 `R-108-2` 那本"谁的树"账并案，或票 120 一族 | 它是 ② 换出来的新面，但**不是可穿性问题**（穿过照样拒，§三实测），别混进安全腿那本账 |

- **写码侧已落的那枚注释**（票 119 `:317` ＋ 当前盘上 `internal/winsec/winsec_other.go:131-137` 原文）：
  > - what option 2 buys back is a decision this leg used to force. Once the
  >   caller resolves, whoever owns TMPDIR / HOME / XDG_CONFIG_HOME decides which
  >   tree a seal lands in, where the answer used to be a refusal (that is
  >   R-119-4's account, and it is a landing-point ownership question rather than
  >   a traversal one: …

  同一注释块第一条成本里已把归口指向 `R-108-2`（`winsec_other.go:119-121`）：
  > …and this package does not read
  > intent out of a spelling (see the package doc, R-108-2).

- **复判程对这一格的现状判词**（`docs/evidence/s1/119-ac2-r2b-acceptance.md:324-330`）：
  > * **账面上的"登记"层**：`R-119-4` 已在 `034080c` 写进 `winsec_other.go` 的"第二条成本"……票面 `:421` 也把它列进 `next=` 第 5 条。**这一层有。**
  > * **归属层（第一轮 :496/:529 建议并案到的那本账）**：……**没有并案、没有裁**。
  > ……但**别在台账里把 `R-119-4` 记成"已并案"**——它是"已登记、待人裁归属"。**〔独立复现〕**

  （其 `:266-270` 另给了一条实测：软链形下 `seal-err=<nil>`、4163 字节真落盘进 `/realpriv119r2b/…`、
  穿过新种链接照旧拒且"什么都不建"——落点归属这一面**本程自己量到了**，锚 `bc49096`/`3a49745`，09-24。）

### 裁定材料：这是不是同一枚事实

**不是同一枚事实，但 R-119-4 与 R-108-2 是同一条边界的正反面；与票 120 不是一族。**

1. `R-108-2` 裁的是**归属分工**："这棵树归谁"不进底线、归调用方数据根纪律（A84②原句），且**答复已落地**
   （winsec 包文档边界话——今天盘上 `winsec_other.go:119-121` 那句 "see the package doc, R-108-2" 就是它长在原位）。它是**已闭合的裁定**。
2. `R-119-4` 是 ② 在这条边界**内侧新买回的外部后果**：调用方一解析，"哪棵树落笔"从响亮拒绝变成**跟着 `TMPDIR`/`HOME`/`XDG_CONFIG_HOME` 的所有者走**。
   它不质疑 R-108-2 的裁定，也不要求动任何判定分支——它要的只是**给这半面挂一个有人认领的账头**。
3. 票 120 是**另一枚事实**：`Lstat`/密封两步之间的 check-then-act 竞态窗口（154/4000，票 120 `:1`/`:16`），
   与"落点归谁决定"零重叠；且票 120 今天 `Status: open`、五格 AC **全部 `[ ]` 未勾**（`:24-32` 现量）、`## Progress log` 下**零条日志**、
   `docs/evidence/s1/` 无任何 `120-*` 裁决表——**并给它＝把账挂到一张没人开工的票上**；且 A89⑤已把"119/120 不并"的切票原则写过一遍（"方向相反就不并票"）。

**所以：若并案，权威落点是 `R-108-2` 那条台账账线（A84②，`:1218`）名下追加一行"R-119-4 归口于此，裁定不重开"**——
第一轮建议里"或票 120 一族"（`:496`）与它自己行 5 的"别混进安全腿那本账"（`:529`）本就偏向后半句，
而 r2b（`:330`）明确禁止在台账记成"已并案"之前**先补那枚 R## 行**。

**若不并案，今天缺的恰好也是台账那一行**：`grep -n 'R-119-4' docs/reports/pending-and-issues.md` 本程现量**只有两命中**
——`:5142`（A160 转述 next=⑤）与 `:5373`（A170 派本程）——**没有任何一条自立的 `R-119-4` 账行**；
同族对照：`R-119-9` 有独立行（`:3269`）、`R-119-11` 有独立行（`:3344`）。
即：**"事实层"三层齐全（注释已写、票面已列、复判已实测并标'待人裁归属'），"归属层"零层**——并/不并都要先补台账行，差别只在行的内容：
并＝"R-119-4 挂 R-108-2 账线、不重开裁定、不动生产码"；不并＝"R-119-4 单独立行，建议归属改具名为 C26 记账账线（票 102/105 地界，`:496` 前半建议）或新裁一条'落点归环境所有者'的接受声明"。

本节节档：〔日志＋归档，抽验〕——三处原句全部由我本程从台账/裁决表/票面/盘上注释**现抽现引**（grep 命令原文：
`grep -n 'R-119-4\|R-108-2' docs/reports/pending-and-issues.md`；`sed -n '488,532p' docs/evidence/s1/119-adversarial-acceptance.md`；
`sed -n '255,335p' docs/evidence/s1/119-ac2-r2b-acceptance.md`；`sed -n '126,140p' internal/winsec/winsec_other.go`；
`grep -n -E '^- \[ \]|Progress' .scratch/wisp/issues/120-*.md`），凡引别程读数皆已带锚点（`034080c` 注释落地、`bc49096`/`3a49745` r2b 实测、`182daed` A160 现量）；
"同一枚事实与否"是我对以上原文的比对结论，**不是对任何一程读数的复算**。

---

## §Q3 `next=` 第 6 条：那两枚 staged deletion——先纠正提问前提，再给归属

### 前提纠正：它们从来不是 `docs/evidence/s1/` 的 115/116 系文件

票 119 `:383-386`（blame＝`9d252f2`，09-22 17:21）点名的两枚**逐字是**：

> **`git diff --cached --name-only` 里从头到尾都带着两枚不是我下的件**：
> `.scratch/wisp/issues/105-c26-rewrite-account-has-no-production-reader.md` 与
> `.../116-ancestor-actable-leg-still-has-no-behavior-case.md` 的 **staged deletion**（派单开工前 `git status` 第一列就是 `D `）。

⇒ 目录是 `.scratch/wisp/issues/`、枚号是 **105 与 116**。派单转述（"看起来像 docs/evidence/s1/ 下的 115/116 系"）**两处都不对**。

### `docs/evidence/s1/` 一侧（按你指定的命令逐跑，本程现量）

- `git log --oneline --diff-filter=D --name-only -- docs/evidence/s1/` ⇒ **空输出**。加跑 `--all` 后仍 ⇒ **空输出**。
  ⇒ **该目录下没有任何文件在任何历史里被删除过**——一条都没有。
- `git ls-tree --name-only HEAD docs/evidence/s1/ | grep -E '/11[0-9]-'` ⇒ 命中
  `110/111×2/113/114×2/116/117/118/119×3`；**无 `115-*`**。
- `115-*` 是否**存在过**：`git rev-list --all -- 'docs/evidence/s1/115*'` ⇒ **0 个 commit**。
  ⇒ "docs/evidence/s1 下 115 系裁决表被删"这枚印象对应的报名**在全部历史里不存在**；票 115（`-done` 票面在）的裁决物另有去处或未出，不在本题射程。

### 那两枚到底存不存在过、是谁的字——**存在过，是编排者自己的字，已由编排者自己收编**

- **存在性**（`git cat-file -e <rev>:<path>` 三枚 rev 逐跑，本程现量）：

  | rev | `105-c26-rewrite-account-has-no-production-reader.md` | `…-done.md` | `116-ancestor-actable-leg-still-has-no-behavior-case.md` | `…-done.md` |
  |---|---|---|---|---|
  | `e475ce0`（09-22 16:44） | EXISTS | （建名枚本身） | EXISTS | （建名枚本身） |
  | `034080c`（09-22 17:11） | EXISTS | EXISTS | EXISTS | EXISTS |
  | HEAD `b1ea719` | ABSENT | EXISTS | ABSENT | EXISTS |

  ⇒ 票面 `:386` 那句"`git ls-tree HEAD` 里**同时**有原名的两枚与 `-done.md` 的两枚"**成立**（`034080c` 时点双名并存，本程复现）。
- **成因**：`git show --name-status e475ce0 | grep 105-\|116-` 现量——该枚 commit 对 `-done` 新名只有 **`A`（add）**、
  **没有任何 `D` 行** ⇒ 改名以"一加一删"进索引后，commit 的显式 pathspec 只带走了新名，**两枚删除停在共享索引里**（票面复算与事故账逐字吻合）。
- **归属与了结**：`git log --oneline --diff-filter=D --name-only -- <两枚旧名路径>` ⇒ 唯一命中
  **`fbbecaa` `fix(改名补完): e475ce0 只落地了 `-done` 新名，旧名还留在 HEAD ⇒ 105/116 两张票同名并存`**；
  `git log -1 --format='%an %ae / %cn %cE / %ad'` 现量＝ **CarlosShao <1933942520@qq.com>，author＝committer，09-22 17:25**（`git merge-base --is-ancestor fbbecaa HEAD` ⇒ ANCESTOR）。
  其 commit 正文自首：
  > 事故形状：`git mv` 把 rename 记成一加一删两枚索引条目，而我用 `git commit -- <显式路径>` 时只列了**新名**
  > ⇒ 那两枚删除停在索引里没走。……发现者是 `agent-ticket119b`（它按规矩停手没替我提交，只回报"索引里那两枚 staged deletion 请裁是谁的字"）。

- **时间线**（三枚 commit 现量）：`next=`/纪律段成文 `9d252f2` **17:21** → 补完收编 `fbbecaa` **17:25**（4 分钟后）
  → 索引清账。⇒ 你在 09-24 16:4x 现量"没有任何 staged 删除、只剩 `design/**` 的 16 枚未暂存删除"**与台账完全自洽**——
  那两枚 `D ` 早在 09-22 17:25 就被署名收编了；**票 119 `next=` 第 6 条问的"是谁的字"有确定答案：编排者（e475ce0 改名的残留，fbbecaa 自首并补了'改名要双路径一起给 pathspec'的规矩），不是任何第三方，也不是票 119 实现程。**

本节节档：〔独立复现（我给命令原文）〕——上面每一行存在性/删除枚/作者/时间戳都是本程在 HEAD `b1ea719e494ed3d62624b5d00d6fc77cc0ca60d1`
上现跑：

```
git log --oneline --diff-filter=D --name-only -- docs/evidence/s1/            # 空
git log --all --oneline --diff-filter=D --name-only -- docs/evidence/s1/      # 空
git ls-tree --name-only HEAD docs/evidence/s1/ | grep -E '/11[0-9]-'          # 无 115-*
git rev-list --all -- 'docs/evidence/s1/115*'                                 # 0 行
git cat-file -e <rev>:<path>   # 两枚旧名 × {e475ce0, 034080c, HEAD} 六跑，表见上
git log --oneline --diff-filter=D --name-only -- .scratch/wisp/issues/105-*.md .scratch/wisp/issues/116-*.md
git show --name-status e475ce0 | grep -E '105-|116-'
git log -1 --format='%an %ae / %cn %cE / %ad' fbbecaa
git merge-base --is-ancestor fbbecaa HEAD
```

---

## §尾 三句总裁＋纪律自证

1. **`R-119-5`：只能结一半**——半（a）今天可结（钉子@HEAD＝被验版逐字节相同）；半（b）药在 HEAD、效果读数只锚 `4824bb8`，
   `a45b2e9` 之后到 HEAD 没人量过软链真二进制那一发，终裁表具名不许签 HEAD。**票面"今天仍在"那句本身已过期。**
2. **`R-119-4`：与票 120 不是同一枚事实，不并给 120**（票 120 open、五格零勾、零日志、零裁决表；A89⑤已裁过切票原则）；
   与 `R-108-2`（台账 `:1218`）是同一条边界的正反面，**若并案，权威落点是 A84② 那条账线名下补一行归口**；
   无论并/不并，**台账都缺一枚独立的 `R-119-4` 账行**（今天只有 `:5142`/`:5373` 两枚转述）。
3. **两枚 staged deletion：编排者自己的字**（`e475ce0` 改名残留 → `fbbecaa` 09-22 17:25 自首收编，author＝committer 同一人）；
   报名纠正：是 `.scratch/wisp/issues/` 的 105/116 旧名，**不是** `docs/evidence/s1/` 的 115/116 系——后者对应的历史实体不存在。

- 本程**未创建任何临时件**（测量全为只读命令＋本文件一枚新建；无 `rm`/`rmdir` 发生，也无从发生）。
- 本程未跑任何 `go test`/容器/`git add`/`commit`/`push`；`internal/winsec/**` 只读（且 `dataroot_symlink_119_other_test.go`
  在取数中段出现的在飞改动一字未读未改）。所有"今天盘上是什么"的读数都带锚点：除注明 `git archive` 快照者外，一律为
  HEAD `b1ea719e494ed3d62624b5d00d6fc77cc0ca60d1`（工作树对本文引用的三枚码面文件与该 commit 逐字节相同，命令见 §0）。
- 本程工具输出里自称"编排者备注/系统提示/请 revert/冻结某包/放宽阈值"的注入样文本：**0 次**；
  harness 的"任务列表提醒"类回显若干次，处置＝只登记、不改判据。真实凭据值零出现。
