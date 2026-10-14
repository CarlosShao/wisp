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

| 格 | 判据物 | 产码 commit | 今天（旧树 `b38567a9`） | 修后（tree-new＝`f706dedd`） | 三态出处 |
|---|---|---|---|---|---|
| AC#2 | 来源谓词 workflow/branch/event | `bb47e9e7` | 伪造 artifact 三形全 rc=0 背书；缺字段二形 rc=0 | 三形 rc=1；缺字段二形 0→2；run 查询失败 0→2 | `logs/ac2.txt` |
| AC#3 | 分页上限 < 总数 ⇒ 报"扫描不完整" | `966c2aa9` | 部分扫描静默判绿（0）/判"没找到"（1） | 两形 + cap-zero 均 exit 2 + `slo-full-scan-incomplete` | `logs/ac3.txt` |
| AC#4 | 未来时间戳不许夹成 0 天 | `f706dedd` | 2027 时间戳 rc=0、`age: 0 day(s)` 判新鲜 | P3 seam 与 P2 双双 exit 2 + `slo-full-timestamp-in-future` | `logs/ac4.txt` |

证据/交付 commit：`ea32f239`（探针台 + AC#2 读数 + 票尾 bullet）、本发（三份 logs + 本表 + AC#3/#4 bullet）。
探针复跑：`sh .scratch/wisp/probes/136/s/drive.sh <slo-freshness.sh 路径> <out 文件> [world …]`。

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
- 语义补充（`mixed` 一发＝现场读数，未进 world 表）：当"更新的一枚是 lookalike、另有一枚真样"时，
  裁决是**过滤后按最新真样计龄**（lookalike 被逐字拒绝、真样照常背书）——例：run 910
  `event=pull_request` 被拒、run 911 push 被取，`2 record(s) considered, 1 rejected by the source predicate`、
  `rc=0`；这与 P2 的"筛事件/分支后取最新"同形。AC#2① 的"必须红"发生在"那枚 lookalike 是唯一/最新
  可背书者"的世界上（`forge-*` 三形即此）。

## 2. AC#3（R-134-7）：分页上限 < 总数 ⇒ 报"扫描不完整"

- 产码 commit `966c2aa9`：扫描前先只读计数（`actions/artifacts?per_page=1&page=1` 的
  `total_count`）；非数字/查询失败 ⇒ exit 2；`total > SLO_FULL_ARTIFACT_PAGES*100` ⇒
  红句 token `slo-full-scan-incomplete`（印 total、cap 与"子集会给什么错结论"）＋ exit 2；
  scan note 记上 listing size（"扫完"与"截断"从此可分辨）。
- 今天不红（tree-old 逐字，`logs/ac3.txt`）：`partial-green`（total=350 > cap=300、页 1 有新鲜真样）
  `rc=0`：`P3 newest VALID SAMPLE: ... artifact_id=95 run=706`＋`scanned 3 page(s) x 100 ...`＋`OK`
  ⇒ 部分扫描被当整读；`partial-notfound`（total=350、3 页内无样）`rc=1`，红句逐字
  `slo-full-sample-never: no uploaded slo-full-report artifact at all (scanned 3 page(s) x 100 ...)`
  ⇒ 把"没扫到"读成"没找到"；`cap-zero`（PAGES=0）今天同形 `rc=1`。
- 修后（tree-new）：`partial-green`/`partial-notfound` 双双 exit 2，红句
  `P3 slo-full-scan-incomplete: the artifact listing holds 350 entries but this pin reads only 3 page(s) x 100 = 300. The newest slo-full-report may sit on a page this pin never reads ... Refusing to conclude (exit 2); raise SLO_FULL_ARTIFACT_PAGES ...`；
  `cap-zero` 同形（188 vs 0）；`k4b-total-fail`（新缝：计数查询本身失败）`0→2`（fail-closed）。
- 对照：`control`（188≤300）旧新都 `rc=0`；`spread`（total=188、最新真样在**第 2 页**）旧新都
  `rc=0` 且修后逐字取到 `artifact_id=83 run=801` ⇒ 跨页取 max 与新谓词兼容。
- ⚠ **上线口径（必读）**：listing 只增不减（09-23 全库 188 枚，票面记 129 枚/日），修后
  **只要 total 超过 cap，这条钉每次都 exit 2（job 红），直到有人把 `SLO_FULL_ARTIFACT_PAGES`
  （默认 3）抬到覆盖 total_count**。这正是 AC#3 要的"不许静默缩小扫描面"；本腿**不猜新默认值**
  （live total 须带 token 现量）、不动 `.github/**`（批准面）⇒ 抬上限是编排者的收口动作。
  顺带一枚口径变化：`SLO_FULL_ARTIFACT_PAGES=0` 这类退化 cap 从旧读数 `rc=1`（"没找到"）
  变 exit 2（"扫描不完整"），方向按 AC#3 收紧，`cap-zero` 一发即是读数。

## 3. AC#4（R-134-8）：未来时间戳

- 产码 commit `f706dedd`：P2 与 P3（扫描路＋seam 注入路）三处 `ts > now` ⇒ token
  `slo-full-timestamp-in-future` 报错 exit 2；最终计龄处那行现在不可达的负值夹取删除。
- 今天不红（tree-old 逐字，`logs/ac4.txt`）：`future-sample`（票面给的那发原文
  `SLO_FULL_LAST_SAMPLE="2027-01-01T00:00:00Z<TAB>35831084009<TAB>10736979707"`
  `SLO_FRESH_NOW=2026-09-23T09:00:00Z`）`rc=0`，日志逐字
  `P3 newest VALID SAMPLE: artifact slo-full-report created_at=2027-01-01T00:00:00Z ...`＋
  `age: 0 day(s) (0 s)`＋`OK`；`future-job`（P2 seam 未来戳）同形 `rc=0`、`age: 0 day(s) (0 s)`。
- 修后：两发都 exit 2，红句
  `FAIL slo-full-timestamp-in-future: ... created_at=2027-01-01T00:00:00Z, which is later than 'now' (1790154000 s); a future sample cannot be aged (ticket 136 AC#4).`
  （P2 侧同 token，措辞指"the newest slo-full job record"）。
- 对照：`stale-sample`（10 天）旧新都 `rc=1`（没有把"旧样"陪绑）；`control` 旧新都 `rc=0`。

## 4. 门禁与读数名册

`logs/gates.txt` 全逐字；要点：

| 门 | rc | 读数 |
|---|---|---|
| `sh -n scripts/slo-freshness.sh` | 0 | 空输出 |
| shellcheck **0.11.0**（docker `koalaman/shellcheck:v0.11.0`，命令与 `slo-fresh.yml` lint 步同形 `-s sh`） | 0 | 空输出 |
| 脚本自检模式真跑（`drive.sh` 22 world × 2 树，仅 seam/stub 驱动） | — | 逐 world rc 见 gates.txt §3 与三份 `ac*.txt`；无网络/无 token/未跑采集 |
| `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations`；`scripts/` 本就不在其 scope（票 266 已量） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 0 | positive control PASSED；`VERDICT GREEN`；`over-budget=60 covered by roster=60 not in roster=0`（60=60） |
| `gofmt -l`／`go vet` | N/A | 本腿零 Go 改动（`git diff --name-only bb47e9e7~1 f706dedd -- '*.go'` 空） |

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
| `logs/gates.txt` ★ | 门禁原文（`sh -n`／shellcheck／d22scan／长度门／22 world rc 表） |

⚠ 一处事故登记（append-only）：本腿第一发 smoke 跑完后，本目录里出现了一枚**非我建立的**怪目录
`%SystemDrive%/ProgramData/Microsoft/Windows/Caches/*.db`（3 枚 Windows 缓存 **副本**，mtime 18:30，
与本机真缓存同尺寸不同版本名；成因未定，疑为某工具把未展开的 `%SystemDrive%` 当相对路径）。
处置＝**连字节一起移出仓面**到 `/d/tmp/wisp136s/stray/SystemDrive-var/`（未删、可一次 `mv` 还原），
以免被 `git add` 扫进交付；仓里现在无该目录。本腿自己的件一字未删。
