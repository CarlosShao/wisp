# 111-c1 起手锚（取数腿 c1 · 票 111 CI 读数格）

- 腿代号：`111-c1`
- 落笔时刻：`2026-10-08 10:54:35 +0800`
- 本件性质：起手锚（第一笔进仓）。下面所有数**全部现量**，未抄任何人口头转述。

## 1. HEAD（现量）

```
$ git log -1 --format=%H
7a367a08292ef94908d86656608974e81c42f0d4
```

- 当前分支：`dev`

## 2. 工作树脏行数（现量）

```
$ git status --porcelain -- cmd internal docs .github scripts | wc -l
0
```

- 具名口径：上面这五个受monitor路径**当前干净**（0 行）。
- 对照量（整仓）：`git status --porcelain | wc -l` = `755`（共享工作树里有别的腿的活；本腿不碰）。

## 3. 远端名册与 tip（现量，`git ls-remote` 走网络）

```
$ git remote -v
cnb     https://cnb.cool/CarlosShao/wisp (fetch)
cnb     https://cnb.cool/CarlosShao/wisp (push)
origin  https://github.com/CarlosShao/wisp.git (fetch)
origin  https://github.com/CarlosShao/wisp.git (push)
```

| 远端 | `ls-remote refs/heads/dev`（tip 现量号） | 本地 tracking ref | 两者是否逐字相等 |
|---|---|---|---|
| `origin`（GitHub） | `cc31526165734e612de297848bb2080bd459ccba` | `cc31526165734e612de297848bb2080bd459ccba` | 相等 |
| `cnb` | `c6cf66e64849955bf92a4b3356c096ea81c35563` | `c6cf66e64849955bf92a4b3356c096ea81c35563` | 相等 |

（`refs/heads/master` 两枚远端均为 `5d777f78a3820f94d40c2afdc0da5a4c91dc089a`，本腿不据此判任何事。）

## 4. 本地未推枚数（现量，逐枚 `rev-list --count`）

```
$ git rev-list --count refs/remotes/origin/dev..HEAD
295
$ git rev-list --count refs/remotes/cnb/dev..HEAD
475
```

- `git rev-list --left-right --count refs/remotes/origin/dev...HEAD` → `0	295`
- `git rev-list --left-right --count refs/remotes/cnb/dev...HEAD` → `0	475`
- 即：本地 HEAD 领先两枚远端，**且落后 0 枚**（无分叉）。

## 5. 推送状态判定（只用"远端 tip 逐字等于本地 HEAD"这一把尺）

- `origin/dev` tip `cc31526…` ≠ HEAD `7a367a0…` ⇒ **未推平**。
- `cnb/dev` tip `c6cf66e…` ≠ HEAD `7a367a0…` ⇒ **未推平**。
- 本腿**不 push**（令：机主口径暂时不再推）。⇒ 凡"需要一次新推送才有的 CI 读数"，本腿一律只能取到**历史那批发过的 run**，取不到的具名列在交付件里。

## 6. 与转述的冲突（具名顶回）

- 编排者转述：「10-06 曾推过一批（约 175 枚），之后本地又新增了很多」。
- 盘上现量：相对 `origin/dev` 未推 **295 枚**，相对 `cnb/dev` 未推 **475 枚**。
- 处理：以盘上现量为准；"175" 视为过期口径，不进任何结论。

## 7. 本腿纪律自报

- ⛔ 不 `go build` / `go vet` / `go test`（Go 编译面与卫生仪器由 `33-v4` 独占）。若只读用到 `go env` / `go list` 会在件里自报。
- ⛔ 不 push、不开真窗、不动 `.github/**`（只读）。
- ⛔ 不在仓库目录内建 worktree/checkout；大输出先落仓外 `D:/tmp/wisp111c1/` 再切。
- 只新建 `.md` 件；每件自落 `rc=N` 一行。
- Git：只 commit 不 push；每笔显式 pathspec；⛔ `add -A` / `add .` / `--amend` / `reset` / `rebase` / `stash` / `checkout .` / `clean`。

rc=0
