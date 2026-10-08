# 242-v2 起手锚（验收腿，非实现者）

跑时刻 `date "+%Y-%m-%d %H:%M %z"` ⇒ `2026-10-08 11:52 +0800`，rc=0。

- 第 1 发尺 `ls -la .scratch/wisp/probes/242/v2/` ⇒ `No such file or directory`，rc=2 = **空**，与派单现量一致。
  父目录 `.scratch/wisp/probes/242/` 现量三件：`precheck.md` `r1/` `v1/`。`v1/`＝10-03 死腿遗产，
  **本腿不引为凭据、不删**。
- `git log --oneline -1` ⇒ `0c9726f9 111-ciif1 起手锚：HEAD 实测 6320a16a（派单写的 e05b8a2b 取不到，顶回）+ ci.yml 起手干净 + 773 条在飞改动登记`，rc=0。
  ⚠ **顶回派单第 0 条**：派单写 HEAD 现量 `e05b8a2b`，实测 HEAD＝`0c9726f9`（另一枚腿 111-ciif1 已落笔）。
  `e05b8a2b` 是否仍是祖先未验；本腿全部读数锚在 `0c9726f9`。
- `git status --porcelain` ⇒ 778 行（含大量不属于本腿的在飞改动）。全文落
  `logs/git-status-start.txt`（rc=0）。**只登记，不提交、不还原。**
- Q1 三把尺：
  - `git merge-base --is-ancestor 4bf7e683 HEAD` ⇒ ancestor_rc=0（**在祖先里**）。
  - `git show --name-status --format= 4bf7e683` ⇒ 一枚 `A internal/agent/approval/ticket242_binding_test.go`；
    `--numstat` ⇒ `113  0  <同一枚>`。零产码、单文件，**与派单转述一致**，rc=0。
  - 工作树 vs HEAD blob：`git hash-object` ＝ `1eaba0972097d714882fd81f9fd7a33984b4fcfc` ＝ `git rev-parse HEAD:<file>` ⇒ SAME，rc=0（工作树无未提交改动）。
  - ⚠ **顶回派单第 1 条（重要）**：派单问「真只动那一枚文件吗……且没被后续改动过吗」。**后半是错的**：
    `git log --oneline -- <file>` ⇒ 除 `4bf7e683` 外还有 `b6b1d6a4`（票 259 腿 259-r1）改过同一枚文件。
    `git rev-parse 4bf7e683:<file>` = `a4932b4466d15a59dd70478177b110741a715992` ≠ HEAD blob。
    差异 `git diff 4bf7e683 HEAD -- <file>` = 12 增 5 删，落 `logs/q1-delta-4bf7-to-HEAD.txt`（rc=0）。
    差异内容＝`spend` 返回值由 `bool` 迁到 `grantDenial`（票 259 AC#2 的机械适配），**判定条件等价**
    （`d == denialNone` 对旧 `true`），断言文案零移动。⇒ 本腿审的是 HEAD 版，不是 `4bf7e683` 版；两版判据不同物，已在正文具名。
