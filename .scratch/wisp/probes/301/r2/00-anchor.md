# 301-r2 / 00-anchor — 起手双锚 + tasklist 现量（写腿，commit-first）

ticket 301 `AC#5` 落地腿（只剩注释面）。本文件是第 1 笔（写腿 commit-first，⛔ 适用"第 45 次调用才 commit"那道只读腿闸门）。

## 三行读数（现量）

```
date                                  = 2026-10-10 15:35:49 +0800 (Saturday)
git rev-parse --short HEAD            = 7e32af9d
git status --porcelain -- scripts internal .github docs | wc -l   = 0
git rev-parse --abbrev-ref HEAD       = dev
```

## tasklist 现量（同一时刻，两条尺）

第 1 发（`tasklist /FI "IMAGENAME eq wisp.exe"` + `head -5`）：

```
wisp.exe                     29388 Console                    1     18,504 K
wisp.exe                     16504 Console                    1    119,280 K
```

第 1 发计数行：`wisp.exe=3` ／ `balldebug.exe=0`

第 2 发（写本件前复核，尺＝`tasklist /FI "IMAGENAME eq wisp.exe" /FO CSV /NH | grep -c '^"wisp.exe"'`）：

```
wisp.exe count     = 0
balldebug.exe count = 0
raw                = INFO: No tasks are running which match the specified criteria.
```

⚠ 读数⛔ 一致（第 1 发 3 枚、第 2 发 0 枚，中间隔约 2 分钟）⇒ 期间有 wisp.exe 实例退出，**本腿没有杀任何进程**（上一腿也不是本腿起的，按纪律只记录）。
⇒ **每发长跑前各量一次**（§3-7），读数并排落进同一枚 rc 件；任一发若⛔ 为 0，该发**⛔ 跑**，只落 `rc=NOT_RUN` + 读数具名说明，绝⛔ 拿"应该没事"顶过去。

## 本腿要动的唯一靶（内容锚，引逐字片段；行号仅记测量时刻位置）

`.github/workflows/ci.yml`，`test-core:` 块头上两行注释，逐字：

```
  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope
  # by platform, not skipped: they run in test-windows / slo jobs.
```

测量时刻位置＝`:400-401`（下一行 `:402` 为 `  test-core:`、`:403` 为 `    runs-on: ubuntu-latest`）；行号锚会漂，**本腿后续所有引用一律内容锚**。
