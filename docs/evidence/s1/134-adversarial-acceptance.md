# 票 134 对抗验收（AC#1..AC#6）— `acceptor-ticket134-r1`，2026-09-23

验收方：`acceptor-ticket134-r1`（独立对抗验收方，非实现方）。票面：
`.scratch/wisp/issues/134-slofull-autostarts-on-every-push-on-this-laptop-but-making-it-manual-only-would-kill-the-gate.md`。
本文件与实现方的三份证据（`134-ac2-ac3-…` / `134-ac4-…` / `134-ac5-…` / `134-ac6-…`）互不替代：
下面每一格要么是**我自己重跑的原文**，要么写明"只复核了实现方的读数、复算方式是 X"。

## 锚定（禁读脏工作树冒充被验版本）

- 开工时 `git rev-parse HEAD` = **`6fdb39dd17011c30bd3ba406cd235bba775643a1`**（分支 `dev`）。
- 被验版本 = **`6effb7e`**，快照 = `git archive 6effb7e | tar -x -C /tmp/wisp134-acc-r1`（仓库内零 worktree、零 checkout）。
- `6effb7e` 是一枚**只改文档**的 commit（`git show --stat` 只有票面 + 本证据树里的 `134-ac6-…md`）；
  AC#6 的**码**在 `decb7b9`（前半）/ `44ab500`（后半）/ `3f17504`（shellcheck 那格），三枚都是 `6effb7e` 的祖先。
  我核过等价：`3f17504` 与 `6effb7e` 上 `scripts/slo-check.ps1`=`6e2ba550`、`scripts/slo-freshness.sh`=`9882f1c4`、
  `.github/workflows/ci.yml`=`c5a98063`、`.github/workflows/slo-fresh.yml`=`18739b2d` 四枚 blob 逐一相同 ⇒ 以
  `6effb7e` 的树作被验面是**超集**，没有"文档已更新而码没更新"的错锚。
- 简报的状态断言复核（**三条不成立，如实报**）：
  1. 「`6fdb39d` 之后的提交只碰 `docs/` 与票面」——**不成立**（口径也错）：`6fdb39d` 就是 HEAD，其后**零枚**提交；
     而我真正该核的是 `6effb7e..HEAD`，它**碰了 `cmd/wisp/**`**：
     `git diff --numstat 6effb7e..HEAD -- scripts .github internal cmd` =
     `754 13 cmd/wisp/leg_sink_gate_131_test.go` / `17 9 cmd/wisp/leg_sink_nail_131_windows_test.go` /
     `6 3 cmd/wisp/slo_windows.go`（票 131 续单在 15:2x 之前入库了）。
     ⇒ 我的裁决一律取自 `6effb7e` 快照，不用工作树；这条不改变结论，但简报那句照字面读会以为"树是干净的"。
  2. 「工作树里 `cmd/wisp/**` 有三枚未提交改动」——**已过期**：15:37 我 `git status --porcelain` 只看到
     1 枚 `M .scratch/wisp/issues/131-…md`，那三枚 `cmd/wisp` 文件已被它们的主人提交。我一个字节没碰。
  3. 「`docs/evidence/s1/124-*` / `131-followup-*` 是未跟踪件」——**已过期**：`5742433` / `c49327c` 已把它们入库
     （`git status` 里不再出现）。同样一个字节没碰。
- 开工前查在飞 run（本机就是 self-hosted runner）：`gh run list --limit 5` 当时有 `35831816899`（`ci`，in_progress）。
  我的时间窗内这台机器又跑了 `35832874239`（07:39:40Z 起，8m58s）⇒ 见下面 §2 与 §7：我有两发本地形状**被这枚 run
  的 `wisp.exe` 撞成 `machine-contended`**，那两发我没有当作读数用。

## 本程临时件清单（只建不删；收尾请一次性清理）

| 目录 | 用途 |
|---|---|
| `/tmp/wisp134-acc-r1` | `6effb7e` 纯净快照（d22scan / YAML / 解析器都在这里跑） |
| `/tmp/wisp134-acc-r1-a` | 变异副本：`slo-stub.cmd`（假 wisp.exe）、`three-states.sh`（红绿红三态）、`lock-test.ps1`（陈旧 report 那一形）、`log/slo-check-{OLD,NEW,REVERTED}.ps1`、`out-*` |
| `/tmp/wisp134-acc-r1-b` | 变异副本：`run-p3.sh` + `run-p3b.sh`（P3 独立复造 W1..W8 / X1..X8）、`log/*.txt`、`scrub_token.py` |
| `/tmp/wisp134-acc-r1-c` | 变异副本：`yaml_check.py`、`direction_check.py`、`cite_audit.py`、`log/ci.yml.orig`（P1 塞 `if:` 的那发已还原） |
| `/tmp/wisp134-acc-r1-d` | 从真 CI 下载的三枚 `slo-full-report` artifact（zip + 解开的 `slo-report.json`）+ `inspect_artifacts.py` |
| `/tmp/wisp134-acc-r1-e` | **未变异**快照副本：`slo-stub.cmd` + `retry-leakpass.sh` + `out-E-*.log` + `retry-log/` |
| `/tmp/wisp134-acc-r1-h` | shellcheck 的 docker 挂载面：`sc/{canary-bad.sh, slo-freshness-44ab500.sh, slo-freshness-6effb7e.sh}` + 各自 `.out` |
| `/tmp/wisp134-acc-r1-probe` | `proc-probe.ps1`（争用探针）、`parse-check.ps1`、`token-count.ps1` |
| `/tmp/wisp134-acc-r1-g-removed-lines.txt` | AC#6 前半那 16 行删除行的原文（逐行列，见 §6.1） |
| `/tmp/wisp134-acc-r1-f-d22scan.log` | 快照上 `sh scripts/d22scan.sh` 的完整输出 |

脱敏交代：我自己的第一版 P3 跑批脚本把 `run "$@"` 连 `GH_TOKEN=<值>` 一起回显进了
`/tmp/wisp134-acc-r1-b/log/summary.txt`。已用 `scrub_token.py`（python `str.replace`，只打印文件名）就地覆写成
`<REDACTED-TOKEN>`，并对我全部临时件目录按"词"筛过一遍（`gho_` 前缀）⇒ 命中 0。**值没有一个字进本报告**。
第二版脚本 `run-p3b.sh` 改成打印前先用 `case` 判是否含 `GH_TOKEN` 并替换。

## 裁决表（与 AC 编号 1:1；**渐进写**：每裁一格补一行 + 补一节，各一枚 commit）

| AC | 判定 | 一句话依据（原文/`file:line` 见同名小节） |
|---|---|---|
| AC#1 | PASS | 形状＝C 为主 + B 为辅，票面 `:179-186` 有 owner 原话与为什么不选 A；动 `ci.yml` 的第一枚 commit 在 AC#1 之后 |

## 1. AC#1 形状由 owner 认（PASS）

票面 `:179-186`（`agent-ticket134` 10:2x 那行）记的是 owner 09-23 10:2x 原话「**这个也都按照你说的来吧**」⇒
定形 **C 为主 + B 为辅、不撤销、不走 A**，并写了为什么不选 A（A 把"读数可能被污染"这支可逆的毛病换成"D32 两条
硬阈值不再有自动结论"这支不可逆的毛病）、撤销口令「就要 A」仍有效。这一条是**票面里可读的授权记录**，不是自述：
它落在 `f6d9ce8`（09-23 10:20 `ticket(134): AC#1 closed - shape is C primary + B secondary, no threshold edits`）。

**"没有 AC#1 不许动 `ci.yml`" 的时序我自己核了一遍**（`git log --format="%h %ad %s" --date=format:"%m-%d %H:%M" -- .github/workflows/ci.yml`，
票 134 那一段只有两枚）：

```
decb7b9 09-23 13:51 feat(134 AC#6 前半): machine-contended 改判"本 run 无结论"…
b9b2072 09-23 10:41 feat(134 AC#2,AC#3): daily schedule for slo-full plus a pin that goes red when it stops being triggered
```

`10:41` 与 `13:51` 都晚于 AC#1 落票面的 `10:20` ⇒ **时序成立**：动触发表之前，形状已经被 owner 认过。
（`7ad6eb4 09-23 09:52` 那枚也碰过 `ci.yml`，但那是 A104/A105/Q-25 那一批的注释与别的 job，不属本票的触发段；
本票的触发表改动面我用 §5 的解析器读数复核过：非注释行只有 `+  schedule:` 与 `+    - cron: '37 19 * * *'`。）

阈值侧：`internal/observe/thresholds.go` 在 `6effb7e`、`b723978`、`3f17504`、`7b4c36a`、`ff4d27b`、`decb7b9`、
`44ab500`、`HEAD` **八枚锚点上都是同一枚 blob `e2677b11b5a5adff8b4eb87c36d31286a274471d`**（`git rev-parse <r>:…` 逐枚算）。

