# 35-a1 — C17 出向（Go→页）落点普查（只读腿）

> 腿：`35-a1`。射程：回答「出向那一跳要落在哪、怎么落、会撞谁」，**不选形**。
> ⛔ 零写入（除本目录）、零 `go test/build/vet`（`272-v1` 正在取整包 Go 读数）。

## S0 起手锚（原文）

```
$ date ; git rev-parse --short HEAD ; git status --porcelain -- cmd internal scripts tools .github docs frontend
Wed Oct  7 11:31:51 CST 2026
367e41b3
（六族路径 status 输出为空 = 工作树在这六族上干净）
```

- 分支：`dev`。
- ⚠ 起手即见：六族路径 status **为空**。工作树里 `design/**`、`.gitignore`、`.scratch/**` 的未提交件不属于本腿，
  不还原、不提交、不评价。终态自查按这六族判，必须与上面逐字相同（= 仍为空）。
- 原始输出：`logs/00-anchor.txt`。

## S1-S6

（进行中，逐节追加；本节先占位，避免空件。）
