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
