# 300-a2 · 起手锚（第 1 笔）

生成时刻：本腿第 2 次工具调用。权威来源＝`git` 对象层，本文件只录读数，不新造规矩。

## 锚点

```
$ git log -1 --format='%h %ad' --date=iso-strict
4e00357f 2026-10-10T17:34:47+08:00
rc=0
```

起手锚＝**终态锚待复核**（见末尾）。

## 台面是否干净（cmd internal scripts .github docs）

```
$ git status --porcelain -- cmd internal scripts .github docs | wc -l
0
rc=0
```

```
$ git status --porcelain | wc -l
820
rc=0
```

解读（读码推的，非量到的）：受影响的五个台面目录**当时 0 枚脏**；全仓 820 枚脏条目在别处（含 `.scratch/**` 工单/探针面）。⇒ 本腿读到的 `HEAD` 与这五个目录的工作树**同一世代**，但**在飞写腿一 commit，HEAD 就会动**，故每件都带锚。

## 进程三连（tasklist）

```
$ tasklist > /tmp/tl_a2.txt ; echo rc=0
rc=0
wisp.exe=0
balldebug.exe=0
msedgewebview2.exe=24
```

（尺＝`grep -c -i -e '<name>'` 对 `tasklist` 全量输出；命中 0 时 `grep -c` 退 1，未接进 `&&` 链。）

## 本腿的尺清单（全程只用这些形状；一律对象层）

1. `git log -1 --format='%h %ad' --date=iso-strict` —— 锚点世代
2. `git status --porcelain -- cmd internal scripts .github docs | wc -l` —— 台面脏数
3. `tasklist` × `grep -c -i -e` —— 进程三连
4. `git show HEAD:<path>` —— 权威表/脚本/CI/测试原文（⛔ 拿工作树当 HEAD）
5. `git grep -n <pattern> HEAD -- internal cmd scripts .github` —— 内容锚名册级 grep（跨 `internal/**` + `cmd/wisp/**`）
6. `git ls-tree -r --name-only HEAD -- internal cmd/wisp` + `cat-file -e` —— 名册与存在性（判文件在不在只认 `cat-file -e`）
7. `git show HEAD:scripts/portable-tests.sh` / `git show HEAD:.github/workflows/ci.yml` × `grep -n 'portable-tests.sh --scope'` —— CI 可见性（调用点⛔ 行号引，只引内容锚）
8. `wc -c` / `wc -l` —— 件非空自证（0 字节＝那格没交）

## 禁止面自陈

- ⛔ Go 编译面：未跑 `go build`/`go test`/`go vet`/`go run`/`go list`。`go env` 累计调用次数＝**0**（截至本笔）。
- ⛔ 产码 / ⛔ 改票面·台账·HANDOVER / ⛔ push / ⛔ `git add -A` / ⛔ `--amend|reset|rebase|stash|checkout .|clean`
- 写面＝仅 `.scratch/wisp/probes/300/a2/**` 新建 `.md`/`.txt`（含本目录 `snap/` 下的 HEAD 快照拷贝，只建不删）。

## rc

`rc=0`
