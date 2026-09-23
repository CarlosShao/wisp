# 票 129 — AC#4（变异自证）+ AC#5（门禁）读数表

锚定 sha：`a71b2d8`（`internal/winsec/**` 自 `db9fafc`（AC#3）以来一字未动，实测
`git diff --stat db9fafc a71b2d8 -- internal/winsec` 空输出）。
本段由第二任代理接续（前任最后写入 09:03、静默三小时、无未提交 WIP），**只做 AC#4/AC#5**，
AC#1..AC#3 的判据与代码不在本段改动范围内。

工作面：仓库树零改动，全部变异落在 `git archive a71b2d8` 解出的快照目录
`/tmp/wisp129c-s23-*`（带本会话后缀 `s23`；仓库内零 worktree / 零 checkout，A38④）。

## STEP 0 — 基线重量（不复用前任的数）

`git archive a71b2d8 | tar -x -C /tmp/wisp129c-s23-base`，在快照目录里
`go test -count=2 -v ./internal/winsec/`（rc=0，日志 `/tmp/wisp129c-s23-base-run.log`）：

四数按 `scripts/winsec-tests.sh` 里那四条 grep 的**逐字形状**数（`^=== RUN` / `^--- PASS` /
`^--- FAIL` / `^--- SKIP`，即只数顶层行）：

| 读数 | 值 |
| --- | --- |
| `=== RUN` | 202 |
| `--- PASS` | 116 |
| `--- FAIL` | 0 |
| `--- SKIP` | 0 |
| `grep -c '(cached)'` | 0 |
| 顶层 `--- PASS` 去重名字数 | 58 |

⇒ 与前任在 AC#3 那一格自述的 `RUN=202 PASS=116 FAIL=0 SKIP=0`（单次 101/58）同形，
但我这次的 202/116 与那枚 58 枚的名字集合是**重新量的**（落盘
`/tmp/wisp129c-s23-base-pass.txt`），后面 AC#4 第④条的 diff 全部以这一枚为基线，不引用前任的数。

（AC#4 的三态读数与 AC#5 的门禁读数在下方逐格追加。）
