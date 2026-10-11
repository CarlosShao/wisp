# 307-a1 · 00-anchor（起手第 1 笔）

腿名 `307-a1`，只读普查腿。本文件是进场台面读数，落笔时刻在任何长跑命令之前。

## 台面读数〔现跑 · 2026-10-10〕

尺与读数（每一条下面都给了原样命令与 `rc`）：

| 项 | 读数 |
|---|---|
| `git rev-parse HEAD` | `1eaa91f66502fdf56c33f220d2c5911d6182bbb3` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- internal cmd` 行数 | `0` 行（即 `internal/` 与 `cmd/` 两目录工作树对 HEAD 干净） |
| `wc -c frontend/dist/index.html` | `1044` 字节 |
| `frontend/dist/index.html` mtime | `2026-10-10 08:51:10 +0800` |

派单未给具体 HEAD 快照号，故不存在「派发号 vs 盘上号」冲突；盘上号即上面这枚，本文件以盘上为准锚定。

## 射程与写面声明

- 写面**只有** `.scratch/wisp/probes/307/a1/**`（本目录）。
- 产码 / `_test.go` / 票面 / 台账 / `docs/**` / `frontend/**` / `design/**`：零字节，不动。
- 任务本体＝票 `.scratch/wisp/issues/307-syncdirs-home-stays-short-while-candidate-folds-long-profile-write-reads-not-sync.md` 的 `AC#0` 那一格（只读数、不碰框、不写「本格闭合」）。

## 自证尺预告（终局逐笔跑）

`git show --name-only --format=%H <每个自己的 commit>`，逐笔核名册是否全在自家区间内、新建 `.go`＝0 枚、`.out`＝0 枚。

## 原始读数留档

见 `logs/00-anchor-readings.log`（同一批命令的 stdout 原样）。
