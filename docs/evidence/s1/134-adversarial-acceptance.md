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
| AC#3 | PASS（一处覆盖面缺口登记为 `R-134-3`） | 我在**被验版本那棵树**上用 `SLO_FRESH_CI` 指向四枚 ci.yml 变异件，P1 四发各红一次、干净件绿；树文件 `diff` 与 `git diff --numstat` 证明真 `ci.yml` 一字节未动 |
| AC#4 | PASS | 三条拒绝理由我各造一发（真 `wisp.exe` 并发 / `RUNNER_TEMP` 前缀命中 / CPU 54% 自然撞），安静时也真放行过；判色改了之后**拒绝出数的条件一字未动**（可执行行差分 + 六处条件行号回读） |

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
## 3. AC#3 反静默死用例：那枚钉有没有牙（PASS；一处缺口见 §3.3）

AC#3 的核心判据是"能红的钉"，不是"存在的钉"。我不复用实现方的 11 发，自己在**被验版本的树**上重造。
两枚仪器：`p1-teeth.sh`（静态探针 P1，变异件全部经 `SLO_FRESH_CI` 指过去，不改树文件）与
`run-p3b.sh`（动态探针 P2/P3 的三态）。

### 3.1 P1 四发结构变异（`/tmp/wisp134-acc-r1-c/log/p1-summary.txt`，逐字尾部）

```
ci-clean       rc=0 |
ci-no-cron     rc=1 | slo-freshness: FAIL slo-full-trigger-missing: … has no 'cron:' schedule entr
ci-no-push     rc=1 | slo-freshness: FAIL slo-full-trigger-missing: … has no 'push:' trigger; slo-
ci-coe         rc=1 | slo-freshness: FAIL slo-full-trigger-missing: the slo-full job now carries 'continue-on-error:' (D22 mode-6: a skippable job …
ci-paths       rc=0 |                                        <- 这一发是缺口，见 §3.3
tree copy still byte-identical to the snapshot ci.yml
```

另外一发我在 §6.2 里跑的（同一枚 P1，另一种形状）：给 `slo-full:` 的 job 体塞 `if: github.repository == 'CarlosShao/wisp'`
⇒ `X7 rc=1 slo-full-trigger-missing: the slo-full job now carries 'if:'`，把 `ci.yml` 还原（`cp` 回原文件 + `diff` 判空）⇒
`X8 rc=0`。**这五发全部落在含 AC#6 `if-no-files-found: warn` 的那棵树上** ⇒ 实现方"AC#6 的 `warn` 改动之后 P1 仍有牙"
那句自述**成立**（我用的是另一组变异件，不是复跑它那几个）。

### 3.2 P2/P3 的三态与两枚旋钮互不顶替（我自己的编号 X/W）

```
X1 rc=1  +10 天世界（真 API 的最新 report 当样、注入新鲜 job 记录）
        FAIL slo-full-sample-stale: the last VALID SLO SAMPLE is 10 day(s) old (> 3) …（P2 那行明写 age: 0 day(s) 绿）
X2 rc=0  同一世界只把 SLO_FULL_SAMPLE_MAX_AGE_DAYS=30 -> 绿（红确实是阈值咬的，不是常数红）
X3 rc=1  还原 -> 又红
X4 rc=1  同一世界只把 P2 的 SLO_FULL_MAX_AGE_DAYS=9999 -> 仍红 slo-full-sample-stale（P3 不是 P2 的别名）
X5 rc=1  反过来：只放宽 SLO_FULL_SAMPLE_MAX_AGE_DAYS=30、把 job 记录摆老 -> 仍红 slo-full-stale（P2 也没被顶掉）
X6 rc=0  真 API 全扫描（今天真实态）：P3 newest VALID SAMPLE: artifact slo-full-report created_at=2026-09-23T07:42:24Z
        artifact_id=10737938571 run=35832874239 / scanned 3 page(s) x 100 … 76 valid-sample record(s) considered
W1 rc=1  "每推都被触发、零份 report"那一形（SLO_FULL_LAST_SAMPLE=none）
        FAIL slo-full-sample-never: no uploaded slo-full-report artifact at all …
W2 rc=1  同一世界把 P2 放宽到 9999 天 -> 仍然 slo-full-sample-never（接缝关不掉探针）
W6 rc=2  无 token（`env -u GH_TOKEN -u GITHUB_TOKEN` + GITHUB_REPOSITORY 有值）
        slo-freshness: no GH_TOKEN/GITHUB_TOKEN - the freshness probes cannot look (this is not a pass)
```

⇒ 实现方交回项里我最在意的三发（"只放宽 P2 仍红"、"`slo-full-sample-never` 会红"、"无 token ⇒ `rc=2`"）**我逐发自己跑过，全部成立**，
读数是上面这三段，不是它的。注入面我也核过：`.github/workflows/slo-fresh.yml` 里四个测试接缝
（`SLO_FRESH_NOW`/`SLO_FULL_LAST_TRIGGER`/`SLO_FULL_LAST_SAMPLE`/`SLO_FRESH_CI`）**一个都没有被生产 workflow 设**
（`yaml_check.py` 逐条打印，只有 `SLO_FULL_SAMPLE_MAX_AGE_DAYS` 出现在注释文字里），
生产步的 env 只有 `GH_TOKEN` ⇒ 上面那些"注入"进不了 CI 那一条路。

### 3.3 我抓到的一处覆盖面缺口（不推翻 AC#3 的勾，登记归单）

P1 静态探针查的是：顶层 `on:` 在不在、`push:` 在不在、branches 里有没有 main/dev、有没有 `cron:`、
`slo-full:` job 体内有没有 `if:`/`continue-on-error:`、以及还跑不跑 `-Subset full`。
它**不查 `paths:` / `paths-ignore:` 触发过滤**（§3.1 的 `ci-paths` 那一发 `rc=0` 就是这件事的读数：
我在 `on.push` 下面加了 `paths: [internal/**]`，钉当场看不见）。票面 AC#6 与 D22 mode-6 的禁令原文是
"`if:` / `continue-on-error` / skip / **路径过滤**"四件并列，所以这是一条**禁令未被探针覆盖**的缺口。
后果不是"门立刻静默死"：那种形状下 3 天后 P2/P3 仍会红（我 X5/W1 那两发就是它的下游），
代价是**最多 3 天的软化期没人响亮说"触发器被动过了"**，而且红的时候报的是"样旧"不是"触发器被改"，归因会绕。
登记为 `R-134-3`（见 §8），本格仍判 PASS：AC#3 要求的那枚"能红的钉"我自己造出来了、也确实红。

## 4. AC#4 有效性检查只往"更常拒绝出数"改 + 真能触发（PASS）

AC#4 有三件要做：拒绝这一形**真能触发**（要贴读数）、拒绝时**不产出任何数字**、以及"只许往更常拒绝的方向改"。
我不复用 `134-ac4-machine-contended-readings.md` 的三发，自己造了四发，其中三发打在**未变异的快照副本**
（`/tmp/wisp134-acc-r1-e/scripts/slo-check.ps1`，`git hash-object` = `6e2ba550d555a373b717ff7717f81b73f9cbd6ba`，
与被验版本同一枚字节）：

1. **理由一（外来工具链/wisp 进程）是真并发撞上的**（跑的是 copy a，但那一条用的是名探针，与被验字节同一份逻辑）：
   ```
   slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: wisp.exe pid=47928 started=2026-09-23 15:40:55 path=E:\work\base\actions-runner\_work\wisp\wisp\build\wisp.exe
   slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: wisp.exe pid=33324 started=2026-09-23 15:40:57 path=…build\wisp.exe
   slo-check.ps1: machine-contended reason: foreign toolchain/wisp process present: wisp.exe pid=49732 started=2026-09-23 15:40:59 path=…build\wisp.exe
   ```
   那三枚 pid 是**当时在飞的 run `35831816899` 的 slo-full 自己**——这一发不是我摆拍的，是本机 runner 真在取样时撞的。
   CI 侧同一条形我也自己从日志读到（`go.exe`/`link.exe`/`gcc.exe`/`cc1.exe`，§2.3 逐字）。
2. **理由二（runner 工作根下的进程）——我自己造的世界**（`three-states.sh`，跑的是从 git 里 `show` 出来的原始 blob，
   阈值与条件一行的没改）：把 `RUNNER_TEMP` 指到 `C:\Windows` 之后，拒绝行逐字：
   ```
   slo-check.ps1: machine-contended reason: process running under the runner work root: sihost.exe pid=5892 started=2026-08-31 19:28:23 path=C:\Windows\system32\sihost.exe
   slo-check.ps1: machine-contended reason: process running under the runner work root: svchost.exe pid=7500 …
   ```
3. **理由三（机器整体 CPU >=50%）自然撞上的，未变异副本**（`out-E-leakpass.log` 头部）：
   ```
   slo-check.ps1: machine-contended reason: machine-wide cpu utilisation 54% over a 1s window (>= 50%)
   slo-check.ps1: NO CONCLUSION (machine-contended): 1 reason(s), 0 state file(s) written, slo-report.json NOT written
   ```
4. **安静时它不是"永远拒绝"**：同一枚未变异副本三连跑，`precheck ok` 分别读到 `cpu max 28%` / `40%` / `43%`
   （`out-E-fail.log` / `out-E-nofile.log` / `out-E-pass.log`）⇒ 拒绝是有条件的。

**拒绝 ⇒ 盘上真没数字**：第 3 发跑完 `OutDir` 里我只 `ls` 到 `slo-no-conclusion.json`，
`state-*.json` / `settle.json` / `slo-report.json` 都不存在；那枚记录里 `state_files_written=0`、
`slo_report_written=false`、`d32_evaluated=false`、`verdict=no-conclusion`。

**"只许往更常拒绝出数"这条方向没被 AC#6 反转**（AC#4 交给 AC#6 的那道硬边界）：
`direction_check.py` 把 `b723978`（AC#6 之前）与 `6effb7e`（被验版本）两枚 blob 各自剥掉注释与
`Write-Host`/`Add-Content`/`Set-Content` 这类上色与落记录行之后做 unified diff——**可执行行的差三处**：
争用分支里 `exit 1` → `exit 0`、新增的 `$noConclusion` 记录对象（纯数据）、step summary 的两行文字。
所有**判定条件**在被验版本里逐条回读，行号一起给：

```
scripts/slo-check.ps1:144  if ($allProcs.Count -lt 2) {            # 进程表读不动 => Fail, 仍然 exit 1
scripts/slo-check.ps1:153  $loadNames = @('go.exe', 'gofmt.exe', …  # 名单一字未增删
scripts/slo-check.ps1:179  if ($loadNames -contains $name) { … }
scripts/slo-check.ps1:182  if ($exePath.StartsWith($prefix, …)) { … }
scripts/slo-check.ps1:211  $cpuBusyPct = 50
scripts/slo-check.ps1:225  if ($cpuMax -ge $cpuBusyPct) { … }
scripts/slo-check.ps1:229  if ($reasons.Count -gt 0) {              # 进这一支之前没有任何取样发生
scripts/slo-check.ps1:396  if (-not $allPass) { exit 1 }            # 出了数而数不过仍然红
```

**D32 两条阈值**：`internal/observe/thresholds.go:19` `memCapSleeping int64 = 25 << 20`、
`:26` `cpuLimitSleeping = 0.5` ⇒ 在被验版本与 `b723978`/`3f17504`/`7b4c36a`/`ff4d27b`/`decb7b9`/`44ab500`/`HEAD`
上**同一枚 blob `e2677b11…`**（§1 那行，八枚锚点逐枚 `git rev-parse`）。再加两条实现方没交过的硬证：
`git log b9b2072~1..6effb7e -- internal/observe/thresholds.go cmd/wisp/slo_windows.go` = **空**；
本票五枚带码的 commit（`e951dfa`/`b9b2072`/`decb7b9`/`44ab500`/`3f17504`）文件清单加起来只有四枚路径
（`scripts/slo-check.ps1`、`scripts/slo-freshness.sh`、`.github/workflows/ci.yml`、`.github/workflows/slo-fresh.yml`）
⇒ AC#6 的具名解冻地界没越。ps1 里 `0.5%`/`25MB` 六处命中（`:29 :135 :241 :242 :256 :257`）全是注释与打印文本，
没有一处是判据。

**AC#6 对 `ci.yml` 的那处超界改动，我单独看了**（票面 `:8-10` 自己报备的那格）：
`git diff b9b2072 6effb7e -- .github/workflows/ci.yml` 只有两个 hunk，都在 `slo-full` 那枚 job 体内——
一是那段注释，二是 `Upload SLO report` 步的 `if-no-files-found: error` → `warn`。
不可避性我判**成立**：门步骤保持无条件（§5 解析器读数：`slo-full` 五枚 step 全 unconditional），
而"零匹配时既不失败、也不创建空 artifact"这一条不是我信它的注释，是我在真 CI 上量到的——
§2.3 那枚争用 run 的产物表里**没有** `slo-full-report`、日志里有那条 `##[warning]No files were found …`。
`slo-smoke` 那一步仍留 `error`（同一枚 diff 里没出现它）⇒ 放宽没有被顺手推广到 hosted 侧。
