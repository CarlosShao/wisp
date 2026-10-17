# 334-r1 事故自报（当场抓到、具名上报，⛔ 改写历史去遮）

时刻＝2026-10-17 12:0x +08。提交＝**`8823549f`**（我那一笔"证据五"）。

## 现象

`git show 8823549f --name-status` 里除我写面那五件之外，多出一枚**不属于本腿写面**的路径：

```
R100  .scratch/wisp/issues/332-cf-invisible-runes-outside-display-boundary-reachability-census.md
   -> .scratch/wisp/issues/332-cf-invisible-runes-outside-display-boundary-reachability-census-done.md
```

＝票 332 那枚"结案改 `-done`"的重命名（别的腿的写面、别的腿的判据动作）落进了**我的提交**。

## 因（现跑的尺，⛔ 猜）

- 我这八笔全程按规矩 `git add -- <逐枚显式文件>`，但**最后一步写成了裸 `git commit -m`**（没带 pathspec）。
- 这个仓是**一枚共享工作树＋一枚共享索引**：`332` 结案腿在我 commit 之前把那枚重命名 `git add` 进了索引、
  还没轮到 commit；我的裸 commit 就把索引里**全部**在飞条目一起收走了。
- 前七笔逐枚 `git show --name-only` 都只含本腿写面（那几刻索引里没有别人 staged 的东西）；
  第八笔撞上＝时序巧合，⛔ 我这把尺"add 用了显式 pathspec"就足以自证清白——它只管 add，不管 commit。

## 损害面（量过才写）

- **内容零损失**：`R100` ＝100% 相似＝纯改名，我这笔没带走他们任何一行正文；
  他们随后那笔 **`dc367a2a`**（parent 就是我的 `8823549f`）在同一枚件上做了 `M`（六格翻勾），
  盘上现状 `…-done.md` 未勾数＝**0**、已勾＝**6**，`git status --porcelain -- .scratch/wisp/issues/332*` 干净＝已入库。
- 损害＝**账目归属**：一枚 `-done` 防重领键（票面规矩里它是"只由裁腿翻"的动作）的**产生记录**写在我的提交号下。
  谁按 `git log --name-only` 追"谁改的 332 票名"，会先撞到我这笔 ⇒ 这是要给编排者补进台账 `A##` 的一条，
  ⛔ 由我盖章算结案。
- ⛔ `reset`／⛔ `rebase`／⛔ `amend`／⛔ revert 遮它：派单与 `AGENTS.md §1.4` 都禁改写已提交历史，
  且那几把尺在共享工作树里会吞掉别人在飞的字节。**要更正就追加**＝本件＋最终回报里那句具名上报。

## 本腿剩下的 commit 立刻改用的安全形

`git commit -F <正文文件> -- <逐枚显式 pathspec>`（把要交的字节点名在 **commit** 上，⛔ 只点名在 add 上、
⛔ 裸 commit、⛔ `-- .`／`-A`／目录形）。本件就是按这形交的（验证＝该笔 `--name-only` 只含 `.scratch/wisp/probes/334/r1/**`）。

## 给编排者的三行（不替它裁）

1. 票 332 的 `-done` 键实际产生于 `8823549f`（我的裸 commit 吞了 `332` 结案腿 staged 的重命名），
   正文改名与六格翻勾仍在 `dc367a2a`（他们自己那笔）⇒ 332 的结案账该记谁，请你按 `issues/README` 规矩裁。
2. 本票 334 的写面自证要加一句限定：`⛔ commit 带 pathspec` 这条我前七笔满足、第八笔未满足＝已具名上报。
3. 建议往派单模板里补一条：**同一棵共享工作树里，子代理的 commit 一律 `git commit -F 正文 -- <显式件名>`**，
   ⓵ 裸 commit 会被别人的 staged 字节污染，⓶ `-- .` 会把别人的在飞工作树字节记进提交（这本仓已入账过）。
