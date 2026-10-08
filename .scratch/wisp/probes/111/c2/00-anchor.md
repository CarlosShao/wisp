# 111-c2 起手锚（取数＋归因腿）

落笔时刻 `2026-10-08 11:1x +0800`。分支 `dev`（共享工作树）。写面**只新建 `.md`**。
本件每一枚数都是**这一发现量**，⛔ 没有抄派单正文里编排者转述的任何数。

## 1. HEAD 锚（尺：`git log -1 --format=%H`）

```
052b393fe5f403ba0ad8d43debc44ce167b79724
rc=0
```

派单要求「不早于 `052b393f`」。尺（`git merge-base --is-ancestor 052b393f HEAD; echo rc=$?`）：

```
rc=0        # 052b393f 是 HEAD 的祖先（本发里 HEAD 就等于它，不是"晚于"是"同一枚"）
```

⇒ **本腿的锚点＝`052b393f` 本身**，不是它之后的某一枚。后面所有行号/名册都锚在这一发。

## 2. 远端位置（⛔ 未 push，本腿只 commit）

尺 `git rev-list --count <remote>/dev..HEAD`：

```
origin/dev..HEAD = 312     rc=0
cnb/dev..HEAD    = 492     rc=0
```

两端 tip 现量（尺 `git rev-parse <remote>/dev`）：

```
origin/dev = cc31526165734e612de297848bb2080bd459ccba   rc=0
cnb/dev    = c6cf66e64849955bf92a4b3356c096ea81c35563   rc=0
```

⚠ 具名注意：`cnb/dev` 的 tip **就是派单里说的 run 300 那枚 `c6cf66e`**，而 `origin/dev` 落后 HEAD 312 枚。
GitHub Actions 跑的是 `origin` ⇒ **CI 上任何一发的 SHA 都不可能是本机 HEAD `052b393f`**，
也大概率不等于 `origin/dev` 当前 tip（本地 fetch 可能是旧的）。这一格在 Q5 里会用量到。

远端名册（尺 `git remote -v`）：

```
cnb     https://cnb.cool/CarlosShao/wisp (fetch|push)
origin  https://github.com/CarlosShao/wisp.git (fetch|push)
rc=0
```

## 3. 本腿射程边界（自检）

- ⛔ 不跑 `go build`/`go vet`/`go test`（编译与测试面归 `255-r1`/`35-r7`）。本腿至今只跑过
  `git` / `ls` / `date` / `gh` / `grep` / `sed`，**未跑过 `go env`/`go list`**（自报：零次）。
- ⛔ 不开任何真窗、不动产码/测试/`.github/**`（只读）、不动 `docs/reports/**`。
- ⛔ 大日志不进上下文：先落仓外 `D:/tmp/wisp111c2/`，在仓外切，仓内只落 `.md` 切片结论。
- 票面 `111-*.md` 只在文末追加一节；追加后自数未勾框应仍为 **5**（尺见 `attribution.md`）。
