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
| AC#2 | PASS（B 那半边有一处判语不成立，见 §2.4） | run/job/step/结论四项我自己从 API 读到；那枚 run 的 `slo-check.ps1` blob 我复算＝`560186fa` |

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

## 2. AC#2 真跑出一枚 run id + job id + step 名 + 结论（PASS；B 那半边见 §2.4）

### 2.1 四项我自己从 API 取的（不引实现方的读数）

```
$ gh api repos/CarlosShao/wisp/actions/runs/35810714576
  name=ci event=push branch=dev sha=58302cc status=completed conclusion=failure created=2026-09-23T02:31:45Z
$ gh api repos/CarlosShao/wisp/actions/jobs/107021435150
  name=slo-full conclusion=success runner labels=self-hosted,wisp-slo
  started=02:31:49Z completed=02:34:30Z          # 窗口 161 s（票面那格写 162 s，差 1 秒，口径含首尾，不改原文只登记）
  step 5 = "SLO full gate (six states + settle + leak)" conclusion=success started=02:32:17Z completed=02:34:17Z
```

⇒ run id + job id + step 名 + 结论四项齐，**这一格按本仓口径成立**。
⚠ 一条我必须自己指出而不是顺着实现方：**整枚 run 的结论是 `failure`**（别的 job 红）。AC#2 那格引的是**第 5 步的结论**，
这没有错，但"这枚 run 绿"这种说法在本票任何地方都不成立——本仓 `ci` 徽章长期全红（`A115②`），
`slo-full` 一枚 job 的绿只是其中一步的绿。

### 2.2 blob 等价（AC#2 那格的"改动真的在跑"）

`gh api repos/CarlosShao/wisp/commits/<sha>` 逐枚取 head_sha，再 `git rev-parse <sha>:scripts/slo-check.ps1`：

| 用途 | run | sha | `slo-check.ps1` blob | `thresholds.go` blob |
|---|---|---|---|---|
| AC#2 那枚 run | `35810714576` | `58302cc` | `560186fa` | `e2677b11` |
| AC#6 前置取证（改前红） | `35817761098` | `df4a60a` | `560186fa` | `e2677b11` |
| AC#6 争用⇒无结论 | `35825185739` | `44ab500` | `6e2ba550` | `e2677b11` |
| AC#6 安静⇒真取样 | `35826548877` | `3f17504` | `6e2ba550` | `e2677b11` |

`560186fa` 与 `e951dfa`（AC#4 那枚）的 blob **同一枚**（`git rev-parse e951dfa:scripts/slo-check.ps1` = `560186fa…`）⇒
票面 AC#2 那句"那枚 run 里的 ps1 与我 `e951dfa` 的 blob 逐字节相同"**成立**。
`6e2ba550` 与被验版本 `6effb7e`/`3f17504` 的 ps1 **同一枚**⇒ 两种绿都是**当前这码脚本**跑出来的，不是旧脚本。

### 2.3 两种绿我各自读了步级原文（`gh api .../actions/jobs/<id>/logs`，取法写明）

争用那发（job `107065251117`）第 5 步内，逐字：

```
2026-09-23T06:06:32.5959005Z slo-check.ps1: NO CONCLUSION (machine-contended) - subset=full refused to sample, no numbers were produced
2026-09-23T06:06:32.5975951Z slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: go.exe pid=46680 started=2026-09-23 14:06:19 path=D:\work\base\go\bin\go.exe
2026-09-23T06:06:32.5991238Z slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: link.exe pid=45912 …
2026-09-23T06:06:32.6008624Z slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: gcc.exe pid=25184 …
2026-09-23T06:06:32.6030546Z slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: cc1.exe pid=54776 started=… path=
2026-09-23T06:06:32.6054940Z slo-check.ps1: NO CONCLUSION (machine-contended): 4 reason(s), 0 state file(s) written, slo-report.json NOT written
2026-09-23T06:06:32.7134223Z slo-check.ps1: NO CONCLUSION (machine-contended): exit 0
2026-09-23T06:06:33.9810367Z ##[warning]No files were found with the provided path: build/slo/slo-report.json. No artifacts will be uploaded.
```

同枚 run 的产物表我单独取了一次（`gh api .../runs/35825185739/artifacts`）：**只有 `slo-smoke-report 2026-09-23T06:07:46Z` 一枚**，
没有 `slo-full-report` ⇒ "争用不给钉续命"这一形在真 CI 上成立（这是 AC#6 后半的关键之一，见 §6）。

安静那发（job `107069434922`，sha `3f17504`）第 5 步内，逐字：

```
2026-09-23T06:24:30.0319569Z slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 28%
2026-09-23T06:24:45.3189869Z slo-check.ps1: state Sleeping exit=0 pass=True
2026-09-23T06:25:00.5316293Z slo-check.ps1: state Armed exit=0 pass=True
2026-09-23T06:25:15.7511029Z slo-check.ps1: state Warm exit=0 pass=True
2026-09-23T06:25:31.0408151Z slo-check.ps1: state Conversation exit=0 pass=True
2026-09-23T06:25:46.3450384Z slo-check.ps1: state PanelOpen exit=0 pass=True
2026-09-23T06:26:01.5819050Z slo-check.ps1: state WorkPeak exit=0 pass=True
2026-09-23T06:26:13.7390318Z slo-check.ps1: settle exit=0 pass=True
2026-09-23T06:26:25.0032952Z slo-check.ps1: leak exit=1 flipped_to_fail=True
2026-09-23T06:26:25.0582596Z slo-check.ps1: report written to E:\work\base\actions-runner\_work\wisp\wisp\build\slo\slo-report.json (all_pass=True)
```

我还把那枚 run 上传的 artifact **下载并解开了**（三枚：`10735461474` / `10735960453` / `10737938571`，
见 `/tmp/wisp134-acc-r1-d/`；三枚 = `10737938571`（run `35832874239`，07:42:24Z）/
`10735461474`（就是上面这枚安静 run `35826548877` 的）/ `10729737755`（run `35810714576`，AC#2 那枚））——
每枚 `slo-report.json` 645 KB、`subset=full`、`results` 六条
（Sleeping/Armed/Warm/Conversation/PanelOpen/WorkPeak）每条 `exit=0 pass=True` 且带 10 个字段的真 report、
`settle.pass=True`、`leak_fixture.flipped_to_fail=True`、`all_pass=True` ⇒ 钉读到的"有效样"今天确实是六态真样，
不是空壳（这一条正面回答 §6.2 的"假 report"问题在当前树上的真值）。
顺带一格实现方没写的：`35810714576`（AC#2 那枚 run）自己也上传了 `slo-full-report`
（artifact `10729737755`，`generated_at=2026-09-23T02:34:17Z`，六态齐、`all_pass=True`）⇒ AC#2 那格不必只靠
"precheck ok" 那行日志间接推断它取到了样。

### 2.4 AC#2 的 B 那半边：`schedule` 的"还没到点"这句判语，一半不成立（新发现）

`gh run list --repo CarlosShao/wisp --event schedule --limit 10` ⇒ **空**（本仓历史上零枚 schedule run）。

- `ci.yml` 的 `37 19 * * *`（19:37Z）：现在 08:37Z，**确实还没到点** ⇒ 这条判语成立。
- `slo-fresh.yml` 的 `23 */6 * * *`（00:23/06:23/12:23/18:23Z）：`gh api .../actions/workflows/slo-fresh.yml`
  回 `"created_at":"2026-09-23T11:41:11.000+08:00"` = **03:41:11Z**，也就是这枚 workflow 在 dev 上已注册 2 h 42 m，
  **06:23:00Z 那一格已经过去、没有产出任何 run** ⇒ 证据 §5/§8.7 的 U2 把两枚 cron 一起写成"都还没到点"，
  对 `slo-fresh.yml` 这一半**不成立**。
  我不据此断言"cron 永远不响"：新加入默认分支的 scheduled workflow 在第一格被 GitHub 跳过是有的（也可能是高负载推迟）。
  但它已经不是"还没到点"，而是"**到点了没响**"，这两者在本仓口径下差一整格证据（AC#3/AC#6 的钉自己那枚时钟从没走过）。
  下一格是 **12:23Z**（距今约 3 h 45 m）⇒ 归给 §8 的 `R-134-4`：到点后再读一次，仍 0 就说明 B 那半边没在交付它承诺的东西。

