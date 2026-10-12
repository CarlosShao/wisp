# 30-gates.md — 票 311 落地腿 `311-r1b` 格三：AC#4 门禁四数 ＋ 越界 ＋ 三态

台面三态（⛔ 任何"改前"那侧都⛔ 用共享工作树跑——`cmd/wisp/**` 此刻由 `295-r1` 在飞）：
`T_pre`＝`git archive 2567e3a9`（起手 HEAD，件 `00-head.txt`）；
`T_post`＝`git archive 0f888a4f`（本腿那四笔之后）；`WT`＝共享工作树。
两棵导出树各先证**非空**：`go_files_in_tree=`**1001**、六件套 `go.mod go.sum internal cmd tools scripts` 齐（件 `03-tree-pre.txt`；同尺跑在 `T_post` 上的输出在 `21-strip-proof.txt` 头部与 `43` §2）。

| # | 门禁 | 改前（件／尺／读数） | 改后（件／尺／读数） | 判定 |
|---|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | `11-pre-build.txt`（`T_pre`）逐字 `rc=0`、件 252 字节 | `30-G1-post-build.txt`（`T_post`）逐字 `rc=0` | ✅ |
| 2 | `sh scripts/d22scan.sh`（★⛔ 共享工作树） | `11-pre-d22scan.txt`（`T_pre`）`rc=0`＋`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77` | `30-G3-post-d22scan.txt`（`T_post`＝块已在 runtests 里）**同逐字串**、`rc=0` | ✅ ★那一格具名读数见下 §1 |
| 3 | `sh scripts/check-path-length-budget.sh --with-self-test` | `00-len-gate-start.txt`（WT，本腿建件**之前**）`VERDICT GREEN`、`len_gate_rc=0`、名册 `over-budget=57`（⛔ 本腿造的） | `30-G7-post-len.txt`（WT，本腿所有件都建完之后）`VERDICT GREEN`、`rc=0`、同一枚 57 | ✅ |
| 4 | `bash scripts/portable-tests.sh --scope=census` totals 逐字未变 | `11-pre-census.txt`（`T_pre`）`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`、`rc=0` | `30-G2-post-census.txt`（`T_post`）＋`30-G6-post-census-worktree.txt`（WT）**同一串**；尺＝三发取那一行后 `sort -u`＝**1** 枚（分母＝3 发），与派单给的现量基线逐字同（件 `55-audit.txt`） | ✅ 作差＝0 |
| 5 | `bash scripts/winsec-tests.sh ./internal/winsec/` 改前改后各 ≥2 发、取交集逐名作差 | `10-winsec-pre-run1/2.txt`（`T_pre`，默认根）两发逐字 `=== RUN=101 --- PASS=58 --- FAIL=0 --- SKIP=0`、`rc=0` | `30-G4/G5-post-winsec-run1/2.txt`（`T_post`）同逐字串、`rc=0`；名册尺两把并排（`^ *--- PASS`=101 含子例 ⇄ `^--- PASS`=58 只顶层，⛔ 相加⛔ 互替），交集 58⇄58，`comm -23`（少了）**0 枚**、`comm -13`（多了）**0 枚**（件 `42-winsec-roster.txt`） | ✅ |

⛔ 那串 `--scope=winsec` 拼写我⛔ 拿它当档位跑过：它在 `24-order-guard-after-guard1.txt`／`24b-order-degraded-before-guard1.txt` 两发里出现，
而且**故意**——那两发要的正是"让 GUARD 1 先说话"的形状（`rc=2` 属 GUARD 1，⛔ 属档位，票面 `:12` 那句坑已在读数里具名）。

## 1. ★"块在 runtests 里、门仍绿"那一格（派单第 5 节末条）

`30-G3-post-d22scan.txt`（`T_post`，`sh scripts/d22scan.sh`，六件套齐、`.go` 1001 枚证非空）同发四件：
① `runtests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/d22scan-20261012-153442-2483`（默认形状，根建在 C:）；
② `ticket 311` 命中 **0**（分母＝同发那行 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`，件 24442 字节）；
③ 该台面的 `tools/d22scan/runtests.sh` 里块**在场**（`mk-block.sh has`=1、`grep -c 'ticket 311 AC#1 shape B'`=4，件 `43` §3；`sha256` 逐枚在 `21-strip-proof.txt`）；
④ `rc=0`。
⇒ 默认形状穿过自己新写的守卫＝照旧绿，守卫一行也没多。

凡红的发我按定式先判"是不是台面故障"：`28-F3`（`GOOS=linux`）与 `25-C2`（裸相对根）两枚 `rc=1` 都⛔ 计入牙、⛔ 计入档位判定，
红因逐字是 `go: creating work dir: … The system cannot find the path specified.`（`25-C2`）与非本机目标跑不起来（`28-F3`），件名与行都在 `41-shot-table.txt`。

## 2. 格式名册（三态并排，射程＝本腿写面那四枚 shell 入口）

件 `43-format-roster.txt`。四枚文件 × 三态（工作树 ⇄ `git show HEAD:` blob ⇄ `T_post` 导出树）：
`bash -n` 全 **0**、stderr 0 字节；`wt_eq_blob_lines=`**0**、`blob_eq_tree_lines=`**0**（逐字节同）；`crlf_bytes_in_wt=`**0**；
`ticket 311 AC#1 shape B` 行枚数三态都是 **4**（四枚文件各 4）。
gofmt/gofumpt 那两把本腿**无对象**：尺＝逐笔 `git show --name-only --format=` 过滤 `\.go$`＝**0 枚**（分母＝本腿总 4 枚文件）
⇒ ⛔ 为凑名册去编一枚"整仓 gofmt 读数"（那数今天由着在飞态决定，归 `295-r1`/机主，⛔ 本腿账）。

## 3. 越界（逐笔尺，⛔ 区间尺）

| 笔 | 文件 | 枚数 |
|---|---|---|
| `e4e787fb` | `scripts/winsec-tests.sh` | 1 |
| `f8b1f446` | `scripts/portable-tests.sh` | 1 |
| `10734fbd` | `scripts/wisp-cli-tests.sh` | 1 |
| `0f888a4f` | `tools/d22scan/runtests.sh` | 1 |

本腿产码面合计 **4** 枚、`.go` **0** 枚；证据面全部落在 `.scratch/wisp/probes/311/r1b/**`（`txt/` 仪器 8 枚＋`ps1/` 尺件 3 枚＋`logs/` 66 枚＋判语 3 枚＝本腿目录 80 枚件）；
票 311 末尾 1 条 Progress-log bullet（1 加／0 删）。
`internal/winsec/**`、`cmd/wisp/**`、`go.mod`、`.gitignore`、`design/**`、`.github/**`、三枚冻结件、`allowlist.txt`、`thresholds.go`、golden、D43、C1–C32、D1–D47 ⇒ 逐笔命中 **0**（分母＝上面 4 枚文件）。

## 4. 交付自查（件 `55-audit.txt`）

0 字节件＝**0** 枚；缺 `rc=` 行的件＝**0** 枚（分母＝`logs/` 66 枚；那一发先读到 1 枚是 `55` 自己在数自己，重跑尺已具名改正）；
`git status --porcelain -- <四枚入口>`＝**0** 行（工作树⇄HEAD 一致）；`git diff --cached --name-only`＝**0** 行（⛔ 别人的在飞态挂在索引上）；
⛔ push（本腿全程零推送）；名册尺起手 `temp_test_star=`**59** ⇄ 交件 **59**（票面在册数，⛔ 动过那 59 枚），
`wisp.exe`／`balldebug.exe` 起手与交件都 **0**（同发正控 `svchost`=87 ⇄ 85、`all_processes`=363 ⇄ 353 ⇒ 那两枚 0 有分母）。
起手盘上有 1 枚活着的 `wisp.test`（PID 24892，`CreationTime 15:14:37`），归属尺逐字指着 `…\wisp-test-gotmp\295r1-…-postG\go-build…\wisp.test.exe`
＝**`295-r1` 的那发**，⛔ 本腿的；本腿⛔ 起过 `wisp.test`／`audio.test`（件 `02-start-roster.txt`／`50-end-roster.txt`，交件时家族尺 `family_test_exe=`0）。
## 5. 落件与提交次序自查（本腿自己的账，⛔ 区间尺——区间里全是别腿的件）
尺＝逐笔 git show --name-only --format=（§3 那张表）；对照：一把**区间**尺 git diff --name-only 2567e3a9..HEAD 在这一刻给 **31** 枚，其中 **25** 枚⛔ 属本腿（＝295-r1 与编排者刚落的件＋四枚入口）⇒ 本腿的账只能逐笔量。
⚠ 一处差点读错的 0：那把区间尺换成 git diff --name-only 起手HEAD..我的笔 -- .scratch/wisp/probes/311/r1b/ 切出来是 0 枚，
它读的是『区间里还没提交的件』＝本腿的未跟踪证据件⛔ 进 diff，⛔ 等于"本腿一枚件都没写"（80 枚件在盘上，尺＝find . -type f | wc -l，件 55-audit.txt）。
staged_lines_at_delivery=0  my_surface_dirty_lines=0
rc=0
