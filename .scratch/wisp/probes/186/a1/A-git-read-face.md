# A | 票 186 AC#1 三支现跑（只读普查腿 186-a1，本仓真 git 树，HEAD=`b54f1740`）

- 程：`186-a1`（只读普查，零产码）。分支 `dev`。跑法：全部在本仓 `D:/work/workspace/projects plans/Wisp` 现跑；本仓自带 3 枚工作树（主仓＋2 linked）＝好样本。
- 规范：每行＝结论 ＋ 尺（可复制可重跑）＋ 读数。

## A① 认 git：读 `.git/HEAD` 的祖先回溯法

- **结论：读文件够，不必起 `git` 进程；且必须做"逐级向上回溯"＋"目录形/文件形两吃"。**
- 尺（根目录）：`test -d .git && cat .git/HEAD`；对照真值 `git rev-parse --abbrev-ref HEAD`。
  - 读数：`.git IS DIR`；`HEAD` ＝ `ref: refs/heads/dev`；真值 ＝ `dev` ⇒ 一致。
- 尺（子目录回溯）：`for d in internal/panel internal .; do test -d "$d/.git" && echo "$d DIR"; done`
  - 读数：`internal/panel` 无 `.git`、`internal` 无 `.git`、`.` ＝ `.git IS DIR` ⇒ 只有根命中。
  - ⇒ 从子目录起读**必须**逐级上溯到命中为止（或到盘根），不能只看一层。
- 尺（文件形＝linked worktree）：`test -f D:/wt/fe/.git && cat D:/wt/fe/.git`
  - 读数：`gitdir: D:/wt/fe/.git` 的反向——`.git` 是**文件**（内容 `gitdir: <主仓>/.git/worktrees/fe`），逐级探测必须同时吃 `test -d` 与 `test -f` 两形。
- 尺（detached 裸 sha 形）：`cat .git/worktrees/wt-head/HEAD`
  - 读数：`bb61dc5c0c5b15783db8736ea5e6fccd6fa50e3f`（40 位裸 sha、无 `ref: ` 前缀）⇒ 判据＝"HEAD 首 5 字节是不是 `ref: `"。
- **一句话：三形（目录/symref、文件/worktree、裸 sha/detached）全部只读文件系统可判，零外部进程。**

## A② 工作树清单：`ls .git/worktrees/` ＋ 主仓拼接

- **结论：目录法给"路径＋分支"两样都拿得到；但清单必须"主仓 ＋ `.git/worktrees/*`"两堆拼——只枚举目录会漏掉用户当前所在的主工作树。**
- 尺（枚数）：`ls -1 .git/worktrees/ | wc -l` ⇒ 读数 **2**（`fe`、`wt-head`）。
- 尺（逐枚 路径＋分支）：`for w in fe wt-head; do cat .git/worktrees/$w/HEAD; cat .git/worktrees/$w/gitdir; done`
  - `fe`：`HEAD` ＝ `ref: refs/heads/dsh/feat/frontend-p0-v2`；`gitdir` ＝ `D:/wt/fe/.git` ⇒ 工作树根 ＝ 该 `.git` 文件的父目录 `D:/wt/fe`。两样齐。
  - `wt-head`：`HEAD` ＝ 裸 sha（detached）；`gitdir` ＝ `D:/tmp/wisp136instr-r1/wt-head/.git` ⇒ 根 `D:/tmp/wisp136instr-r1/wt-head`。两样齐（分支维＝"detached"）。
  - 附带态：`test -f .git/worktrees/$w/locked`（两枚均 NO）、`prune`（两枚均 NO）。
- 尺（真值对照）：`git worktree list --porcelain` ⇒ 读数 **3 条**：
  - `D:/work/workspace/projects plans/Wisp` → `refs/heads/dev`（**主仓，不在 `.git/worktrees/` 里**）
  - `D:/tmp/wisp136instr-r1/wt-head` → `detached`
  - `D:/wt/fe` → `refs/heads/dsh/feat/frontend-p0-v2`
- ⇒ 目录法 2 枚 ＋ 主仓（读 `.git/HEAD` 与 `.git` 父目录）拼 1 枚 ＝ 3 枚，与真值逐枚吻合。
- **⇒ `git worktree list --porcelain` 不必跑**（"路径＋分支"两维不缺；locked/prunable 若要显示，文件形＝看 `locked`/`prune` 文件在不在）。

## A③ 分支清单：`refs/heads/` ＋ `packed-refs` 两堆合并

- **结论：必须"递归"走 `refs/heads`；本仓 `packed-refs` 不存在，只走 loose 就读全（合并逻辑今天给不出真数据）。**
- 尺（单层——**会少报**）：`ls -1 .git/refs/heads | wc -l` ⇒ 读数 **3**（`dev`、`dsh`、`master`）。
- 尺（递归——**对**）：`find .git/refs/heads -type f | wc -l` ⇒ 读数 **4**。
  - 逐枚：`.git/refs/heads/dev`、`.git/refs/heads/dsh/feat/frontend-p0`、`.git/refs/heads/dsh/feat/frontend-p0-v2`、`.git/refs/heads/master`。
  - ⇒ 带斜杠的分支名在文件系统里是**嵌套目录**（`dsh/feat/...`）；任何单层列举在本仓今天就会少报一枚。
- 尺（真值对照）：`git branch --list | wc -l` ⇒ 读数 **4**，与递归法相等。真值输出里 `*`＝当前分支 `dev`、`+`＝`dsh/feat/frontend-p0-v2` 已被一枚 linked worktree 检出。
- 尺（packed）：`test -f .git/packed-refs` ⇒ 读数 **NO**（本仓根本没有该文件）⇒ 合并逻辑只是格式陈述，**未被真数据验过**（承接 09-28 普查 §⑨#1）。

## A④ 不跑 `git` 二进制能不能读全——逐支答

- ① 认 git：**能**（`.git/HEAD` 一条文件 ＋ 逐级回溯；目录/文件两形都吃）。不起进程。
- ② 工作树清单：**能**（`.git/worktrees/*` 逐枚读 `HEAD`＋`gitdir`，＋主仓读 `.git/HEAD` 与父目录拼一枚）。不起进程。
- ③ 分支清单：**能**（递归 `refs/heads`；`packed-refs` 存在时按行 `<sha> refs/heads/<name>` 并读、同名 loose 覆盖 packed）。本仓无 packed ⇒ 只走 loose 即读全。不起进程。
- **⇒ 三支全部"不跑 `git` 二进制可读全"，本票这一维没有"非跑外部 `git` 不可"的缺维。**
- **若哪天真要起外部 `git`（本程结论＝不必），管它的门是哪条（现读具名）：**
  - `tools/d22scan/main.go` 的九条 ban 逐名（`:8-72` 头注释）：#1 bare-goroutine、#2 pathresolver-bypass、#3 plaintext-key、#4 wallclock-timeout、#5 mirror-hash、#6 panel-approval、#7 internal-artifact-tool、#8 emoji、#9 phantom-citation。
    - **没有任何一条管"起外部进程"**（与 09-28 普查 §⑫#5 一致）。
  - 真跑会咬到的是：**#1**（等子进程不能裸 `go`，要走 `observe.Registry.Spawn`）＋**#4**（超时必须单调钟、⛔ 墙钟差）。
  - `tools/d22scan/allowlist.txt` 现读全文（912 B，6 行）：条目＝`pathresolver-bypass`×3（`internal/memory/open.go`、`internal/models/manifest.go`、`internal/risk/pathresolver.go`）、`mirror-hash`×1（`internal/models/downloader.go`）、`bare-goroutine`×1（`internal/observe/goroutine.go`）。**零条关于 `exec`／`git`**（尺：`grep -n -i 'exec\|git' tools/d22scan/allowlist.txt` ⇒ 空输出）。
  - ⇒ 结论：**跑外部 `git` 归 ban #1／#4 管（不是 allowlist 条目），且今天全仓零次 ⇒ 属一枚新的宿主能力。** 本程现跑三支均证明**不必**引入它。

## A⑤ 一枚与 09-28 普查相比的现量差（供落地腿留意）

- 主工作树 HEAD 现值 ＝ `b54f1740cf377287164729443bfb5e2de9261f4a`（`git worktree list --porcelain` 与 `git rev-parse HEAD` 一致），工作树数仍 3 枚、`fe` 仍指向 `dsh/feat/frontend-p0-v2`、`wt-head` 仍 detached ⇒ 样本形状未变，20 天前的三支结论今天仍成立。
