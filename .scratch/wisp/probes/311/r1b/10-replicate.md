# 10-replicate.md — 票 311 落地腿 `311-r1b` 格一：把 3b857e53 那枚换卷拒绝块复制进其余三枚入口

**锚**：起手 `HEAD`＝`2567e3a9`（件 `logs/00-head.txt`，尺＝`git log -1 --format='%H %ad %s'`，台面＝工作树）。
起手写面尺＝`git status --porcelain -- <那四枚>`＝**0 行**，且 `git diff --name-only 3b857e53..HEAD -- scripts tools`＝**0 枚**
⇒ 停手条件（别腿在写同一批入口）没成立。交件时我这四笔＝`e4e787fb`（winsec 注释）／`f8b1f446`（portable）／
`10734fbd`（cli）／`0f888a4f`（runtests），逐笔 `git show --name-only --format=`＝**各 1 枚文件**（⛔ 区间尺，件 `logs/43-*.txt`）。

**三态台面**（下文每个数都注明）：
- 工作树＝`D:/work/workspace/projects plans/Wisp`（⛔ 任何"改前"读数都不在这里取——`cmd/wisp/**` 此刻由 `295-r1` 在飞）
- `T_pre`＝`git archive 2567e3a9` 导出的**仓外**树＝改前那侧（件 `logs/03-tree-pre.txt`：`go_files_in_tree=`**1001**，六件套齐全，块枚数 winsec=1／其余三枚=0）
- `T_post`＝`git archive 0f888a4f`（含我这四笔）导出的**仓外**树＝改后那侧（`go_files=`**1001**，四枚块 `has`=1，`sha256` 逐枚在件 `logs/21-strip-proof.txt`）

## 1. 插入点：复跑死腿那张表，逐枚现量（尺＝`grep -n` 内容锚，件 `logs/01-insertion-points.txt`）

行号只作定位辅注，判据是**内容锚**；三枚目标的锚都是**该文件里唯一命中**（尺＝`grep -cE`，件 `logs/13-replicate-insert.txt`）：

| 入口 | 复认的锚（唯一性＝1 枚） | 落在"哪一步之后／哪一步之前" | 死腿量到的号 | 我现量的号 |
|---|---|---|---|---|
| `scripts/portable-tests.sh` | `^capture=` | skip pattern 那行 `:758` 之后、`sh "$strict"` 那一委派之前 | 758/760/764 | 758/760/764（同） |
| `scripts/wisp-cli-tests.sh` | `^# PATH is handed over` | 非 windows＋三枚 DLL 面守卫的 `fi`（`:132`）之后、PATH staging 之前 | 132/134/146 | 132/134/146（同） |
| `tools/d22scan/runtests.sh` | `^out=` | 310 根块（含 DEGRADED 分支）收尾 `fi`（`:107`）之后、`go test -v -count=1`（`:116`）之前 | 107/109/116 | 107/109/116（同） |

⇒ 死腿那张表**一枚不差**，我的块就落在它指的缝里；块自己的起始行（插完之后）＝`portable:760`／`cli:134`／`runtests:109`（`grep -n -F '# ---- TICKET 311 AC#1'`，件 `logs/14-replicate-diff.txt`）。

## 2. 复制这件事是**拷贝**、不是抄写（尺＝件 `logs/13`＋`logs/14`）

`txt/mk-block.sh` 按内容锚从 `scripts/winsec-tests.sh` 切出那块（注释头到块自己的**列零 `fi`**，块内每一枚 `fi` 都带缩进），
只换一处＝`sed` 把 4 行 `winsec-tests.sh:` 前缀换成该入口自己的前缀，再插到唯一锚之前。三枚各 `block_lines=`**71**、
`prefixed_echo_lines=`**4**。四枚 `bash -n` 全 rc=0、stderr **0 字节**（分母：同发 `lines=` 枚数 280/909/218/223）。

四枚块逐行 diff（`T_post`/工作树同一形状，件 `logs/14-replicate-diff.txt`）：

| 对 | differing_lines | 差在哪 |
|---|---|---|
| winsec ⇄ portable | **8** | 4 枚 echo 行的前缀（INCONCLUSIVE 1 ＋ GUARD 3），一进一出＝8 |
| winsec ⇄ cli | **8** | 同上 |
| winsec ⇄ runtests | **8** | 同上 |

**"哪枚 gate"那一格我⛔ 改**（派单允许我改，我选择不改，理由具名）：块里唯一带档位身份的那句是
`this gate's COLOUR is a function of that volume - the same unchanged internal/winsec code reads PASS=58 FAIL=0 … PASS=56 FAIL=2 …`，
它报的是**本机量过的那一对颜色**（票 311 现量＋`311-a1`＋编排者复跑）。把它换成 `cmd/wisp` 或 `tools/d22scan` 的档位＝
给一枚⛔ 存在过的颜色编造凭据（`cmd/wisp` 换卷色今天本机根本没量，见 `20-teeth.md` §6）。所以四枚**逐字同形、只差前缀**，
`A862` 给票 310 那块立的"源行比较"口径在这里同样成立。

## 3. 那段注释指名过不存在的件（`3b857e53` 〔未验证半成品〕的由来）

块里 `Ordering readings:` 指的两枚 `probes/311/r1/logs/23-order-degraded-then-guard.txt`／`-24-order-guard-after-guard1.txt`
盘上⛔ 存在（死腿 15:0x 断连，尺＝`ls -1 .scratch/wisp/probes/311/r1/logs/`＝只有 `00/01/02/10/11/12/13/14/15` 九枚，
逐枚名在 `probes/311/r1/` 那批已入库的件里）。我这枚把它们**真跑出来了**，落在我自家目录同名两枚
（`logs/23-…`／`logs/24-…`，加一枚代证的 `logs/24b-…`），并把注释的两条路径改指 r1b＋具名写了"r1 那两枚从没落盘"这笔
（＝`e4e787fb`，7 行→8 行、纯注释面，numstat 5 加 4 删）。读数判语在 `20-teeth.md` §4。

## 4. 顺手那一件（派单第 9 节）：票 310 那块，我同意／我不同意哪半句

**先确认那块今天还是 310-r2 交来的样子**：`git diff --numstat 3b857e53^ 3b857e53 -- scripts/winsec-tests.sh`＝**71 加 0 删**、
单枚 hunk `@@ -121,0 +122,71 @@`⇒ 纯插入，既有一行⛔ 动；更硬的一把＝blob 号
`git rev-parse 23c19afa:scripts/winsec-tests.sh`＝`31f64e0c…` ⇄ `3b857e53^` 同文件＝`31f64e0c…` **同一枚 blob**。

**判语（件 `logs/40-a862-310block-check.txt`，尺＝同一枚 `mk-block.sh` 把 START 锚换成 `# ---- TICKET 310 AC#1` 之后整块对拉）**：

- ✅ **"跑在 GUARD 之前"＝同意**，四枚都坐实（内容锚行号，改后台面）：
  `winsec 310@46 < 第一枚 : GUARD@95`、`portable 99 < 179`、`cli 60 < 98`、`runtests 68 < 174`；
  并且**真机一发**：`LOCALAPPDATA= WISP_TEST_GOTMP_ROOT= bash scripts/winsec-tests.sh --scope=winsec`（`T_post`）
  ⇒ 逐字先打 `[310 DEGRADED 行]`、再由 `[GUARD 1]` 收尾 `rc=2`（件 `logs/24b-…`）＝A862 ⑤ 那一发在我这枚台面上复现。
  ⚠ 一处具名收窄：`runtests.sh` 的 310 块**之前**还有 4 枚 `exit 2` 参数守卫（`:46/:51/:61/:65`），它们的句子⛔ 带
  `GUARD` 这个词，所以"第一枚 : GUARD@174"那把尺看不见它们——这句"跑在 GUARD 之前"对 runtests 只能说"跑在
  **自称 GUARD 的那些**守卫之前"。winsec/cli/portable 各只有 1 枚 `exit 2` 在 310 之前（`41-44`／`54-57`／`93-96`，都是"委派对象那枚脚本⛔ 在"的前置检查）。
- ⚠ **"四枚同形"＝按 A862 自己那句收窄版同意，按字面版不同意**。整块对拉的差行＝
  `winsec⇄cli` **27** 行、`winsec⇄runtests` **47** 行、`winsec⇄portable` **89** 行（件 `logs/40`，逐字差集同件）。
  A862 ⑥ 那把尺是 `grep -a DEGRADED` 的**句首**比较，它已经自带"我这把量的是源码行、⛔ 渲染后的整句"那句收窄 ⇒ 收窄版我认。
  但字面"四枚同形"够不着两件事，具名登记：
  1. **`portable` 多一枚分支**：非 windows 那一支它打 `GOTMPDIR left as-is on GOOS=…` 一行，另三枚同支**一声不响**
     （尺＝`logs/28-F3-runtests-gooslinux.txt` 310 命中 **0**，分母＝同发 `runtests.sh: go test exited 1 …` 行在；
     ⇄ `29-P3` 改前同发亦 0 ⇒ 两枚台面一致，⛔ 是我这块造成的）。这是**行为**差、⛔ 是换行差。
  2. **同句不同字**：空 `LOCALAPPDATA` 那句的交叉引用，`portable` 写 `see the note above`、另三枚写
     `see the note in scripts/portable-tests.sh`；`cli` 那枚还把 `nothing said that the redirection did not happen`
     写成 `nothing said the redirection did not happen`。渲染出来读得出的差异，句首尺读不出。
  ⇒ 对本票的实际后果：**我这枚 311 块四枚只差前缀（§2 那 8 行），比票 310 那块更同形**；
  票 310 那块的这三处方言**⛔ 由本腿修**（修它＝动 310 的写面、且 `A862` 那格的盖章人是编排者本人，口令写着
  「310-r2 补一枚验收腿」）⇒ 我只把读数交到这儿。

## 5. 我自报两枚仪器缺陷（都当场抓到、当场改，⛔ 留着当"跑过"）

1. `txt/run-shot.sh` 第一版把相对件名留到 `cd` 之后再写 ⇒ 四枚件**一枚没落**（`No such file or directory`），
   我按 stderr 抓到、改成先绝对化＋`[ -d workdir ]` 当场拒绝（件 `logs/09-runner-selftest.txt` 好支；坏支＝**故意**指
   一枚不存在的树，尺返回 `REFUSED workdir is not a directory: […]`＋`exit 96`，且**盘上没落任何件**——所以这一发没有 09b
   那枚文件，"没落盘"本身就是那条读数，命令逐字记在这儿：
   `sh txt/run-shot.sh logs/09b-… "/d/tmp/311r1b/nonexistent-tree" -- pwd`）。
   同批还有一枚是我自己的引号漏洞：`$L/xx.txt` 没加引号，被仓库路径里那个空格**词分裂**成三枚参数。
2. `txt/mk-block.sh` 第一版把块锚当 awk **动态正则**用，`(shape B` 那枚裸括号＝unmatched group ⇒ extract 切出 0 行，
   而 `insert` 仍然"成功"了——它往三枚文件各塞了**一枚空行**（尺＝`git diff --stat`＝三枚各 `1 +`）。
   处置：三枚空行按内容逐枚删回（`git diff` 对那三枚＝**0 行**才算复原），再把尺换成**固定串** `index()==1`／`$0=="fi"`、
   并给 `insert` 加"切出少于 40 行就 REFUSED"、给 `strip` 加"没切到东西就 REFUSED"两道闸（第二轮它果然用
   `exit 97` 拦住了我，见件 `logs/13` 逐字 `REFUSED: extracted block is only 0 lines`）。

⚠ 这两枚都⛔ 影响交付物的正确性，但**第一枚差点让我把"四枚件不存在"读成"那四发⛔ 跑"**——记在这儿是因为它离"假读数"只差我没去 `ls` 那一下。
