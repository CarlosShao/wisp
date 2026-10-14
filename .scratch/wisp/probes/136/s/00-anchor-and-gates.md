# 票 136 写腿 136s —— 新鲜度钉三缺口落地（AC#2 / AC#3 / AC#4 ＝ R-134-5/7/8）

- 腿身份＝实现者（判语归非实现者腿）；⛔ 未翻框、⛔ 未动 `Status:`/票名、⛔ 未 push。
- 锚（起手现读）＝`b38567a9`（与派发锚相同，零漂移）；读数窗口 2026-10-14 18:2x–20:xx +08。
- 面＝`scripts/slo-freshness.sh`（产码 三笔）＋`.scratch/wisp/probes/136/s/**`（探针）。
  `scripts/slo-check.ps1` **一字未动**（理由见 §5「具名没做」第 3 条）。
- 探针方法：`drive.sh` 用脚本自己的 `SLO_*` 缝（`SLO_FRESH_NOW` / `SLO_FULL_LAST_TRIGGER` /
  `SLO_FULL_LAST_SAMPLE` / `SLO_FULL_ARTIFACT_PAGES`）＋ PATH 首位一枚 `gh` 桩
  （134 验收 K 形同款仪器；`ghjq.py` 只认脚本实际传的 filter，不认就 exit 1 ⇒ 钉子落 exit 2），
  **不真跑采集**（无网络、无 token、⛔ 不跑 `wisp slo`、⛔ 不跑 `slo-check.ps1`），
  时钟钉死 `SLO_FRESH_NOW=2026-09-23T09:00:00Z`，每发都注入新鲜 P2 记录 ⇒ 红只能来自 P3。
  两棵树：tree-old＝`git show HEAD:scripts/slo-freshness.sh`；tree-new＝修后工作树副本。

## 0. 结论先行

| 格 | 判据物 | 今天（旧树） | 修后（tree-new） | 三态出处 |
|---|---|---|---|---|
| AC#2 | 来源谓词 workflow/branch/event | 伪造 artifact 三形全 rc=0 背书 | 三形 rc=1；缺字段二形 0→2；run 查询失败 0→2 | `logs/ac2.txt` |
| AC#3 | 分页上限 < 总数 ⇒ 报"扫描不完整" | 部分扫描静默判绿/判"没找到"（0/1） | 两形 exit 2 + `slo-full-scan-incomplete` | `logs/ac3.txt` |
| AC#4 | 未来时间戳不许夹成 0 天 | 2027 时间戳 rc=0 判新鲜 | P3/P2 双双 exit 2 + `slo-full-timestamp-in-future` | `logs/ac4.txt` |

## 1. AC#2（R-134-5）：来源谓词

- 产码 commit `bb47e9e7`（`scripts/slo-freshness.sh`）：扫描行加第 4 列 `head_branch`；
  字段完整性（`created_at`/`workflow_run.id`/`head_branch`）缺一即 exit 2；
  分支谓词 `main|dev`＋workflow（`.github/workflows/ci.yml`）＋事件
  （`push|schedule|workflow_dispatch`，与 P2 同一套）对 branch 合格候选按时间倒序核，
  首个通过即停；被拒 lookalike 逐枚打印＋计数，零合格样本用既有 token
  `slo-full-sample-never` 点名三档谓词。
- 今天不红（tree-old 逐字）：`forge-branch`/`forge-workflow`/`forge-event` 三形各
  `== ... rc=0`，日志里 `P3 newest VALID SAMPLE: ... artifact_id=90/91/92 ...` 之后
  `age: 0 day(s)`＋`OK - ...` ⇒ 名单对、来源错的 artifact 当场背书（H7 同形）。
  `missing-source` 今天 rc=0（`run=0`，静默当新鲜）；`missing-ts` 今天更糟：
  读串位成 `created_at=704 artifact_id= run=93`，仍 `age: 0 day(s)` → rc=0。
- 修后（tree-new 红句点名）：`forge-branch` ⇒
  `P3 source predicate: artifact 90 on branch 'feature/evil' rejected (branch must be main or dev; ticket 136 AC#2)`
  ＋`FAIL slo-full-sample-never: 1 slo-full-report artifact(s) were found but none carries the source predicate ...`（rc=1）；
  `forge-workflow` ⇒ `workflow '.github/workflows/chore.yml' ... rejected`；`forge-event` ⇒
  `event 'pull_request' ... rejected`；`missing-source` ⇒
  `carries no workflow_run fields (run='0' branch='' artifact=94) ... refusing (ticket 136 AC#2)`（exit 2）；
  `missing-ts` ⇒ `carries no created_at (run=704 artifact=93) ... refusing (ticket 136 AC#2)`（exit 2）。
- 取不到数据"六形"（AC#2② 的那一族）修后逐形 rc：k1 无 token 2／k2 无 repo 2／
  k3 runs 查询失败 2／k4 artifacts 查询失败 2／k5 空列表 1／k6 空列表＋旧样 1
  ⇒ **无一 pass**；另加 k4b（总数查询失败）0→2、k7（run 查询失败）0→2（新缝也 fail-closed）。
- 还原/对照：`control`（真·dev+ci.yml+push）旧新都 rc=0，修后打印
  `2 record(s) considered, 0 rejected by the source predicate`；`stale-sample` 旧新都 rc=1。

## 2. AC#3（R-134-7）：分页上限不许静默缩小扫描面

<!-- 待填（本格落地时补：产码 commit、今天读数、修后红句、对照） -->

## 3. AC#4（R-134-8）：未来时间戳

<!-- 待填 -->

## 4. 门禁与读数名册

<!-- 待填 -->

## 5. 具名没做

- ⛔ **AC#5**（`tools/d22scan/**` 的 `paths:` 那一件）——那面要**具名解冻**才动得了，本腿未碰。
- ⛔ **AC#6**（CI schedule 与本仓 `.github/**`）——`slo-fresh.yml`/`ci.yml` 属**批准面**，本腿零字节改动。
- ⛔ **AC#7 里属 `internal/observe` 的部分**——那半已结（AC#1/#8/#9/#12/#13/#14/#15），本腿未碰。
- ⛔ **AC#10**（真跑 CLI＋机时）——本腿未跑任何真采样。
- ⛔ `scripts/slo-check.ps1` 一字未动：三条判据的落点全在 `slo-freshness.sh`（AC#2 的替代修法
  "把最新 artifact 下下来核内容"未被派单选中，来源谓词那条才是本格判据），ps1 侧零需求。
- ⛔ 未动：`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`tools/d22scan/**`、`allowlist.txt`、
  任何阈值/golden、`.github/**`、`frontend/**`、`design/**`。
- ⛔ 未 push、未 `add -A`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；临时件只建不删。

## 6. 读数名册（`.scratch/wisp/probes/136/s/`，★＝追库，全 ≤200 行）

| 件 | 是什么 |
|---|---|
| `00-anchor-and-gates.md` ★ | 本表 |
| `drive.sh` ★ | 探针驱动器（world 矩阵 + gh 桩，复跑：`sh drive.sh <script> <out> [world…]`） |
| `ghjq.py` ★ | gh 桩的 `--jq` 替身（仓内无 jq 时用；只认脚本实际传的 filter） |
| `logs/ac2.txt` ★ | AC#2 两树逐字（三态） |
| `logs/ac3.txt` ★ | AC#3 两树逐字 |
| `logs/ac4.txt` ★ | AC#4 两树逐字 |
| `logs/gates.txt` ★ | 门禁原文（`sh -n`／shellcheck／d22scan／长度门） |
