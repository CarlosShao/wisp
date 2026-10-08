# 197-a1 起手锚（只读普查腿）

- 腿号：`197-a1`
- 任务：全仓普查「要求合并/折叠面板快照分片」这一族说法，判与 `internal/panel/pump.go`
  「Overflow TRUNCATES and never merges」那条已裁规则是否同链冲突。
- 写面：只新建 `.md`（本目录 `.scratch/wisp/probes/197/a1/`）。⛔ 不改任何产码/测试/票面/台账，
  ⛔ 不勾任何 AC 框，⛔ 不新建 `.sh`/`.ps1`/`.txt`/`.out`。
- 跑面：⛔ `go build`/`go vet`/`go test`（编译面由 `255-r1`、`35-r7` 在飞）。
  `go env`/`go list` 若使用会在交件里自报。
- Git：只 commit 本目录自己的 `.md`，显式 pathspec，⛔ `git add -A`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`，⛔ push。

## 锚（现量，见下方逐条 rc）
```
date_local=2026-10-08 11:18 +0800
branch=dev rc_branch=0
HEAD=052b393fe5f403ba0ad8d43debc44ce167b79724 rc_head=0
HEAD_line=052b393fe5f403ba0ad8d43debc44ce167b79724 2026-10-08 11:14 rc=0
is_ancestor_052b393f_rc=0
own_path_status_lines=1 rc=0
worktree_dirty_files_total=35 rc=0
```

