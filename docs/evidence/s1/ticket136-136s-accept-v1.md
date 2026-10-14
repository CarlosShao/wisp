# ticket 136 — AC#2 / AC#3 / AC#4 非实现者复核（leg 136-v），裁决表 v1

**身份** = 非实现者（实现者＝136s；本腿只读 `scripts/**`、零写 `probes/136/s/**`）。
**锚** = 起手现读 `git rev-parse --short HEAD` = **`1740fd48`**（派发链顶为 `f8b92e7e`；漂移逐报见 §0）。
**方法** = 仓外四树（`old`=`b38567a9` / `pre`=`f706dedd` / `head`=提交版，`cmp` 与 `b6fc818d` 逐字同 /
`wt`=工作树含 `:-20`）＋本腿自建 **33-world 夹具＋自写 `gh` 桩＋自写驱动器**（`probes/136/v/` 只存读数与表；
驱动器/桩/夹具在仓外 `D:\tmp\wisp136v\`，只建不删）。⛔ 未以重跑 136s 的 `drive.sh` 充当交付；对跑一次仅作对照（§4.3）。
**总裁**：**AC#2 成立（PASS、无附条件）；AC#3 成立（PASS、无附条件）；AC#4 成立（PASS、无附条件）**；
**零放宽**（四树 132 发 rc，无一条由非绿变绿）＋**门禁全绿**；**收口项（`:-20`）算术成立、方向正确，且已被现量证明为必要项**
——该行**至今只在工作树、未入库**，而 `origin/dev` 上带新码的 **12:07:32Z schedule run 已经照新逻辑变红**（§5.2）。

## §0 锚、漂移、面（每条现量，无一条推理）

| # | 事项 | 读数 |
|---|---|---|
| 0.1 | HEAD 漂移 | `1740fd48`（133-r1）≠ 派发链顶 `f8b92e7e`；`git diff f8b92e7e..HEAD -- scripts/` **空** |
| 0.2 | `b38567a9..HEAD -- scripts/` | **两枚**：`slo-freshness.sh`（136s 四提交）＋`check-path-length-budget.sh`（`41071cdc`＝编排者 178 收口，1+/1-）⇒ 派单"只许一枚"在今天的范围上不成立，**多出那枚零 136s 归属** |
| 0.3 | numstat 归因 | AC#2 `bb47e9e7` **133/22**（派单那对数字＝这一枚）；AC#3 `966c2aa9` 41/6；AC#4 `f706dedd` 25/3；seam `b6fc818d` 15/2；顺序小计 214/33；**全范围净 `--numstat` = 211+/30-**（后笔改写了前笔插入的行，只有净值是两棵树之间的事实）；`-U0` 全范围 hunk=20 |
| 0.4 | P1 面 | `.github/**` 在范围内零改动；`ci.yml` blob `80b7a393` 两树同 ⇒ 四棵树共用同一 ci.yml |
| 0.5 | 提交面 | 136s 七枚提交逐枚 `--name-only`：四枚产码**只含** `scripts/slo-freshness.sh`；三枚交付只含 `probes/136/s/**`＋票文件；**无外来路径** |
| 0.6 | 收口未入库 | 我现场量：`git show origin/dev:scripts/slo-freshness.sh` 第 138 行仍是 `:-3`；同时 `merge-base --is-ancestor` 逐枚证实 **五个 136s 提交都已在 `origin/dev`** ⇒ 红面已经暴露（§5.2） |
| 0.7 | 两处错位目录（我自造，只建不删，登记） | `D:\d\tmp\wisp136v\worlds\*`（首跑 POSIX 路径滑落）与 `D:\work\workspace\projects plans Wisp\.scratch\wisp\probes\136\v\`（首跑少一个斜杠；该 twin **无 .git**、且早已有他腿同病内容 `probes/272`／空 `probes/orch`）；真交付已写回真仓，twin 件留存 |

## §1 逐格判语

### §1.1 AC#2（`R-134-5`，来源谓词）＝ **成立**

判据逐条对读（全部本腿现跑；逐字原文见 `probes/136/v/ac2-verbatim.txt`，全量日志 `out/<tree>--<world>.log`）：

| 判据 | 旧树（`b38567a9`） | 新树（`head`） | 判 |
|---|---|---|---|
| ① 三个 forge 形（名字对来源错） | `forge-branch` / `forge-workflow` / `forge-event` 逐发 **rc=0**，且逐字给它背书：`P3 newest VALID SAMPLE: ... artifact_id=501/502/503` | 逐发 **rc=1**；红句逐字点名三档谓词之一：`artifact 501 on branch 'feature/evil' rejected (branch must be main or dev; ticket 136 AC#2)` ／ `workflow '.github/workflows/chore.yml'` ／ `event 'pull_request'`，并收在 `FAIL slo-full-sample-never: 1 ... none carries the source predicate ... 1 rejected` | 过 |
| ② 真 artifact 缺字段 | `missing-ts` **rc=0** 且把列读串位成 `created_at=920 artifact_id= run=520`（leading-TAB 被 IFS 折叠 ⇒ 移位）还判 `age: 0 day(s) (0 s)`；`missing-source` **rc=0** 给 `artifact_id=530 run=0` 背书 | 两形各 **exit 2**：`carries no created_at (run=920 artifact=520) - 'no timestamp' is not 'endlessly fresh'` ／ `carries no workflow_run fields (run='0' branch='' artifact=530)` | 过 |
| ③ 取不到数据六形 | k1..k6 = `2/2/2/2/1/1` | 同 `2/2/2/2/1/1`（**无一 pass**；k5/k6 是"确无样本"的正当 rc=1 红） | 过 |
| ④ 新缝 k4b/k7 | `k4b`（总数查询失败）**rc=0** 静默取旧路径；`k7`（run 查询失败）**rc=0** | 各 **exit 2**（`the listing-size query failed` ／ `the run look-up for 950 failed`） | 过 |
| ⑤ `mixed`（伪造更新样 510＋真样 511 同页） | rc=0 但背书**伪造件**：`newest VALID SAMPLE ... artifact_id=510` | rc=0 且：先逐字 `artifact 510 on branch 'release/forge' rejected`，再 `newest VALID SAMPLE ... artifact_id=511`，计数行 `2 record(s) considered, 1 rejected by the source predicate` | 过 |
| ⑥ 被拒 lookalike 逐枚打印＋计数 | 无此概念 | 见 ⑤；零合格样本时由既有 token `slo-full-sample-never` 点名三档谓词＋`N rejected` | 过 |

**一处口径（136s 已在其脚本头写明，我复量一致）**：`SLO_FULL_LAST_SAMPLE` 注入的记录是**声明式 fixture**，
不再过来源谓词（`seam-forged-aged` 两树 rc=0）。我对生产面做了核查：`grep` 全仓 `.github/**` 与
`scripts/**`（除该 pin 自身）**没有任何地方设置** `SLO_FULL_LAST_SAMPLE`／`SLO_FULL_LAST_TRIGGER`／
`SLO_FRESH_NOW`／`SLO_FULL_ARTIFACT_PAGES`（rc=1 无命中）⇒ 该 seam 无生产入口，谓词守的 API 扫描路就是唯一生产输入，
此豁免**不构成洞**。
**bin6fc818d 补的另一半也在**：seam 路的空首字段旧树 rc=0（`artifact_id= run=35831084009`，900 s 判新鲜）⇒ 新树 exit 2
（`an injected ... record starts with an empty field`）；`pre`↔`head` 全 33 世界/四树的**唯一差**就是这一发（§4.1）。

### §1.2 AC#3（`R-134-7`，分页上限静默缩面）＝ **成立**

| 世界（cap=3，除注明） | 旧树 | 新树（`head`） | 判 |
|---|---|---|---|
| `p-green`（total=350>300，第 1 页有新鲜真样） | **rc=0**：`P3 newest VALID SAMPLE ... artifact_id=570`＋`OK` —— 部分扫描当整读 | **exit 2**，红句逐字 `P3 slo-full-scan-incomplete: the artifact listing holds 350 entries but this pin reads only 3 page(s) x 100 = 300. ... Refusing to conclude (exit 2); raise SLO_FULL_ARTIFACT_PAGES ...` | 过 |
| `p-notfound`（total=350，无样） | **rc=1**：`no uploaded slo-full-report artifact at all (scanned 3 page(s) x 100 ...)` —— 把"没扫到"读成"没找到" | **exit 2**（同上一红句） | 过 |
| `cap-zero`（`PAGES=0`） | **rc=1**（"scanned 0 page(s) ... no artifact found"） | **exit 2**（`holds 188 entries but this pin reads only 0 page(s) x 100 = 0`） | 过 |
| `ctrl188` / `spread2`（真样在**第 2 页**） | 各 rc=0 | 各 rc=0；`spread2` 取到 `artifact_id=601 run=1001` ⇒ 跨页 max 与新谓词兼容、对照组不受累 | 过 |
| 触发条件（本腿自加边界） | 旧树无 total 概念 | cap=1：`total=100` **过**（rc=0）／ `total=101` **响**（exit 2）⇒ 条件就是**严格 `listing_size > cap`** | 过 |
| `live349`（total=349；**默认 cap**） | rc=0 | `head`（默认仍 3）**exit 2**；`wt`（默认 20）rc=0，逐字 `scanned 20 page(s) x 100 of the artifact listing (listing size 349)` | 过 |
| >2000 会响（本腿自加边界，不必"代证"） | — | `wt`：`total=2000` **rc=0**（`listing size 2000`，扫描 20 页全读到）／`total=2001` **exit 2**（`holds 2001 entries but this pin reads only 20 page(s) x 100 = 2000`）⇒ 恰在 2001 首响 | 过 |

**一处必要口径**：`p-green`/`p-notfound` 的"新树 exit 2"以 **cap<total** 为条件（这就是 AC#3 判据自身的条件）。
在**收口后的默认 cap=20** 上 total=350 已不是"部分扫描"，两形回到 rc=0/rc=1 —— 这是**正确行为**（20 页覆盖 350 是真整读），
不是放宽：判据"压到小于总数 ⇒ 必须报扫描不完整"在 (0,188)、(1,101)、(3,350)、(3,349)、(3,2000)、(20,2001) 六组 (cap,total)
上逐组现量成立。
**现量 live（真 API，非桩）**：`head` 版 **rc=2**（listing 350→351 时各一发，逐字同上红句）；`wt` 版 5 发中 **4 发 rc=0**
（`newest VALID SAMPLE artifact_id=11862826438 run=38484004521`，`123 record(s) considered, 0 rejected`），1 发因空尾页
`page=7: EOF` 瞬断而 exit 2（fail-closed 行为如设计；见 §5.3）。

### §1.3 AC#4（`R-134-8`，未来戳夹成 0 天）＝ **成立**

| 世界 | 旧树 | 新树 | 判 |
|---|---|---|---|
| `future-sample`（seam 样 2027-01-01） | **rc=0**，逐字 `created_at=2027-01-01T00:00:00Z`＋`age: 0 day(s) (0 s)`（负值夹 0） | **exit 2**：`FAIL slo-full-timestamp-in-future: an injected slo-full-report record claims created_at=2027-01-01T00:00:00Z, which is later than 'now' (1790154000 s); a future sample cannot be aged (ticket 136 AC#4).` | 过 |
| `future-job`（P2 记录 2027-01-01） | **rc=0**＋`age: 0 day(s) (0 s)` | **exit 2**（P2 侧同 token：`the newest slo-full job record claims ... Check the clock of whatever produced it.`） | 过 |
| `future-scan`（**扫描路**未来戳；本腿自加——三处 `ts > now` 的第三处） | **rc=0**（背书 630，`age 0 (0 s)`） | **exit 2**（`a slo-full-report record claims created_at=...`） | 过 |
| `stale-sample` / `stale-job` | rc=1 / rc=1（`age: 10 day(s)`） | rc=1 / rc=1，逐字同 | 过 |
| `real`（对照） | rc=0 | rc=0 | 过 |

## §2 两树（共四树）逐发 rc 表 + 零放宽

**全表**：`probes/136/v/rctable.tsv`（33 世界 × `old/pre/head/wt`）。`old→head` 迁移直方：

| old→head | 枚数 | 世界 |
|---|---|---|
| 0→0 | 7 | real bnd-at ctrl188 mixed p2-live-path seam-forged-aged spread2 |
| 0→1 | 3 | forge-branch forge-workflow forge-event |
| 0→2 | 12 | future-sample/-job/-scan k4b k7 missing-ts missing-source bnd-over d-cap2000 live349 p-green 等 |
| 1→1 | 4 | k5 k6 stale-sample stale-job |
| 1→2 | 3 | cap-zero d-cap2001 p-notfound |
| 2→2 | 4 | k1 k2 k3 k4-pages |

**由非绿变绿的世界：0 枚**（`old rc≠0 → head rc=0` 为空集；`pre → head` 亦 0 枚）。**旧树本来就红的全部仍红**（stale-*、k5、k6、k1–k4）。
`p2-live-path`（P2 走 runs+jobs 真 API 路、无 seam）两树 rc=0 ⇒ 桩不是只认 seam。

## §3 门禁（全由本腿现跑）

| 门 | rc | 读数 |
|---|---|---|
| `sh -n`（head＋wt） | 0 / 0 | 静默 |
| `dash -n`（head＋wt） | 0 / 0 | 静默（CI 用 ubuntu `/bin/sh`＝dash） |
| shellcheck **v0.11.0**（docker，`-s sh`，与 `slo-fresh.yml:78-83` 同形） | 0 / 0 | 无输出 |
| `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations` |
| 长度门 `--with-self-test` | 0 | `VERDICT GREEN`，over-budget=60／roster=60／not-in-roster=0 |
| 跨 shell 全矩阵（dash，33×4） | — | 与 bash 表 **132/132 逐发相同** |

## §4 旁证

1. **b6fc818d 未弱化任何既有世界**：`pre`（f706dedd）↔ `head`（b6fc818d）逐世界对比，33×4 中唯一差异＝
   `seam-empty-first` **0→2**（收紧方向）；且没有任何世界 rc 变小。
2. **`%SystemDrive%` 缓存杂目录**：`find . -maxdepth 4 -iname '*SystemDrive*'` 仓内**空**；移出件在
   `/d/tmp/wisp136s/stray/SystemDrive-var`（仓外，现量存在）。
3. **136s 自己那套 harness 的对跑**（仅对照，非交付）：`WS` 改指 `/d/tmp/wisp136v/theirws` 重跑其 `drive.sh`，
   22 世界 ×2 树的 rc 与其 `logs/gates.txt` **1:1 相同**（含 `future-sample rc=2`、`k4b 0→2`、`k7 0→2`）。

## §5 编排者收口项（`:-20` 与注释）复核

### §5.1 算术独立复算

- `20 pages = 20 x 100 = 2000` 逐字成立。
- 注释写的 **348**（A953，11:38Z 现量）是**带时刻的事实**：我同一会话内现量 **349 → 350 → 351**（+3 在约 42 分钟内，
  push 活跃期），live 跑自身读到 350/351 ⇒ 该数字是油表、不是恒量；算术容差极大（351 ≪ 2000）。
- **"~7 months"**：按注释自己的 ~8/day：(2000−348)/8 = 206.5 d = 6.8 mo（≈7）成立；按我推导的窗口均值
  **(351−188)/21.1 d = 7.73/d**：(2000−351)/7.73 = 213 d = 7.0 mo，同样成立。
- **但速率是有窗口的量**：newest-100 窗口 ≈ 5.24/d（10.3 mo）；全 listing 生涯均值 ≈ 14.7/d（3.7 mo）；
  今天实测 burst ≈ 103/d（16 d）；票面历史 09-21 尖峰 129/d（约 13 d）。⇒ "~7 months" 只在"~8/d"这个前提下为真。

### §5.2 必要性已被**生产现量**证明（本格最重要的读数）

- 事实链：五个 136s 提交已在 `origin/dev`（默认分支）＋**该分支第 138 行仍是 `:-3`**（收口未提交）
  ＋ 现量 listing 350/351 > 300 ⇒ 下一发定时**必红**——而它**已经发生了**：
  **run `38561478699`，event=schedule，createdAt=2026-10-14T12:07:32Z，conclusion=FAILURE**，失败步
  `slo-full freshness pin (ticket 134 AC#3)`，逐字：
  `slo-freshness: P3 slo-full-scan-incomplete: the artifact listing holds 350 entries but this pin reads only 3 page(s) x 100 = 300. ... Refusing to conclude (exit 2);`
  ＋ `##[error]Process completed with exit code 2.`
  （上一发 05:04:04Z 在同一 workflow 上是 success、旧 sha `e3788f08`；原始日志存 `out/prod-run-38561478699.log`。）
- 同一会话对 **`wt` 版**真 API 复跑 5 发：**4 发 rc=0**（含 `0 rejected` 的最新真样）⇒ **这一行字的入库就是这条红的解**。
- ⇒ **具名建议（裁定归编排者，本腿零改动）**：①**尽快把 `:-20`（含注释）随同波次提交并推 `origin/dev`**，
  否则每 6h 的 schedule 会连红；若想立刻把颜色扳回，可再点一发 `workflow_dispatch`。②可选（需另立写腿、另走验证）：
  把扫描循环上界收成 `min(cap, ceil(total/100))`——今天 350 枚只需 4 页调用，而现默认会打 16–20 发**空尾页**调用；
  每发空尾页都是一个"瞬断即 exit 2"的暴露点（§5.3 实证 1 发）。此收紧不碰 AC#3 保证（`total>cap` 仍 exit 2，早停只发生在 total≤cap 时）。③可选措辞：
  注释里给速率标窗口（"since 09-23 平均 ~8/day；burst 可上百"），免下一位把 348/8d 当恒量。

### §5.3 两枚实测瞬断（登记，不改变任何判语）

- live `wt` 首量一发死于**空尾页** `Get "...artifacts?per_page=100&page=7": EOF` ⇒ exit 2（只存在于我终端、未存盘，逐字见
  `probes/136/v/live-readings.txt` 的同文登记）；随后 4 发全绿。
- 我自己的 raw 查询（12:18:15Z）也吃到同一形 EOF，3 分钟后重试即好 ⇒ 属网络/GitHub 侧瞬断，不是脚本或桩的产物；
  **fail-closed 把瞬断翻成 exit 2 是设计方向**，但 cap=20 的 16–18 发空尾页调用把暴露面放大了（见 §5.2②）。

## §6 具名没做（本格的边界）

- ⛔ 未动 `scripts/**` 与 `probes/136/s/**`（零写）；⛔ 未翻任何票框、未改 Status；⛔ 未 commit 收口（裁定权在编排者）。
- 未观察 18:23Z 的下一发 schedule（时间未到）；12:07Z 那发是既成事实、已引原文。
- 未在 Linux 容器里跑矩阵（改走本机 `dash` 全矩阵 132/132 代 CI shell；容器面不在三格判据内）。
- 未测 `SLO_FULL_MAX_AGE_DAYS`／`SLO_FULL_SAMPLE_MAX_AGE_DAYS` 的非默认取值（不属三格判据）。
- 未复现 136s 的全部 21+1 世界的**逐字**输出（对跑核 rc 表 1:1 相同；逐字层用我自己的 33 世界对照到判据面）。
- 瞬断 EOF 的复现率只到 n=5（1/5）；未做统计意义上的率估计。

## 总判

**AC#2 成立、AC#3 成立、AC#4 成立（各 PASS、无附条件）**；零放宽、门禁全绿、收口算术成立且必须立即入库（§5.2）。
证据件：`probes/136/v/`（7 件：本目录 `00-anchor-and-gates.md`、`method.md`、`rctable.tsv`、
`ac2-verbatim.txt`、`ac3-verbatim.txt`、`ac4-verbatim.txt`、`live-readings.txt`）。
