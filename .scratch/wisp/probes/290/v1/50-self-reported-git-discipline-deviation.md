# 290-v1 · 50 自报一条本腿造成的 Git 纪律偏差（⛔ 不改写历史，只追加记录）

## 事实

起手时索引里有一枚**别腿（编排者）已 staged 但未提交**的删除：
```
git status --porcelain（本腿第 1 次调用，逐字）：
  D  .scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md
```
成因：`5cff604e`（编排者的"票 290 改短名"那一笔）用**显式 pathspec 只提交了新名那一侧**——
本腿现跑尺：`git cat-file -e 5cff604e:".scratch/wisp/issues/290-nothing-…-the-ball.md"` → **rc=0（旧名当时仍在 HEAD）**；
`git show --stat -M 5cff604e` 只列 1 枚文件（`76 insertions`）。⇒ 删除那一侧留在了索引里等下一笔。

本腿第 1 笔用的是 `git add -- <我的件>` ＋ **`git commit`（不带 pathspec）**——
`git commit` 不带 pathspec 时提交的是**整个索引**，于是把那枚别腿的 staged 删除一起收了：
```
git show --name-status --no-renames 2f525f41   →
  D  .scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md
  A  .scratch/wisp/probes/290/v1/00-anchor-and-static-scales.md
```
⇒ **判语：这一处违反派单"commit 必带显式 pathspec"的字面**（`add` 带了，`commit` 没带），
并因此**替编排者提交了编排者的那一笔删除**。

## 影响面（本腿现跑尺）

- **零内容丢失**：旧名的完整内容在 `5cff604e` 及更早每一笔里都在，`git cat-file -e HEAD:<旧名>` 现在 rc=128 只是"HEAD 不再带它"，
  而**这正是编排者那一笔改名的意图**（他们 `:76` 那条 Progress log 逐字：改完现跑 over-budget 57 == roster 57 ⇒ VERDICT GREEN）。
- **磁盘状态未被本腿改动过一枚**：起手时旧名就**不在工作树**（所以 status 才写 `D ` 在第一列），本腿没删过任何文件。
- **别人的在飞件一枚未混入**：`git status --porcelain | grep -E "^ D design/|^ M .gitignore" | wc -l` ＝ **17**（16 枚 `design/` ＋ `.gitignore`），
  与本腿起手时同一个数；`.scratch/wisp/probes/{152,161,242,268}/**` 那批也都还在原处未提交。
  名册行数 759 → 758，差的那一行正是被 `2f525f41` 收走的 staged 删除。
- **本腿后三笔都是干净的**：
  ```
  d33a7235 → A 10-mutations-and-cites.md
  49c4ce4f → A 20-questions-n5-n6.md
  1a50cf93 → A 30-gates-and-rosters.md, A 40-ac-verdicts.md
  ```
- **零 push**（`git log --oneline origin/dev..HEAD | wc -l` ＝ 622，含本腿四笔，全在本地）。

## ⛔ 本腿没有做的更正动作（按规矩只能追加）

⛔ 未 `--amend`、⛔ 未 `reset`、⛔ 未 `rebase`、⛔ 未把旧名恢复回去（恢复＝在树里再放一枚 106 字符的票面副本，
会把 `scripts/check-path-length-budget.sh` 重新打成 RED，而且那是**票面动作**、不归验收腿）。

## 给下一腿的台件写法（一条，具体）

`commit` 那一枚**也要带 pathspec**，`add` 带了不算：
```
git add -- .scratch/wisp/probes/<票号>/<腿>/xx.md
git commit -q -m "…" -- .scratch/wisp/probes/<票号>/<腿>/xx.md      ← 只提交这条 pathspec
```
并且**起手先跑一枚 `git status --porcelain | grep -E "^[MADRC]" | head`**（第一列非空格＝索引里已有别腿的 staged 件），
看见就**先具名报给编排者再落第 1 笔**，别让它搭自己的车。
