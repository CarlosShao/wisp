# 票 254 — 第五档 `winsec` 今天**在 CI 里没有调用者**：`winsec-tests.sh` 仍走显式路径，而 `ci.yml` 的形状属契约级（251-r1 按派单停手没动）

**立票时刻**：2026-10-02 10:2x +08，锚点 HEAD `0cdd3dd3`（`dev`）
**来路**：`251-r1` 交件时的两条"甲（谁补得上）"＋一条"乙（停手未动）"；台账 `A523`。**不夹进票 251**（那票四格已做完，本票是"把新门接上电源"这一步）。
**性质**：仪器／CI 口径。owner 用软件撞不到；⚠ 但**不接上，票 251 那道新 GUARD C 在 CI 上永远不会被触发**。

## 现量（编排者本机自跑，2026-10-02 10:2x）

1. `bash scripts/portable-tests-selftest.sh` ⇒ **`18 case(s) ran, 0 assertion(s) failed`、`rc=0`**（载具已含第五档的两形正控）。
2. `bash scripts/portable-tests.sh --scope=winsecfoo` ⇒ 逐字 `unknown --scope=winsecfoo (known: core, windows, cli, winsec, census)`、**`rc=2`** ⇒ 第五档进名册了、未知档仍响亮拒答。
3. 腿报（我未复跑，属〔仅自述〕级，落地腿起手要复认）：`scripts/winsec-tests.sh:94` 那一行**把包路径显式传给** `portable-tests.sh`，CI 侧由 `.github/workflows/ci.yml:439` 调 `winsec-tests.sh` ⇒ **今天没有任何一步以 `--scope=winsec` 的名义去跑那档**。

## 要建什么

- [x] **AC#1 把第五档接上电源（二选一，先交代价再动手）**：ⓐ 让 `scripts/winsec-tests.sh:94` 改点 `--scope=winsec`（写面只在 `scripts/`，最小）；ⓑ 保持显式路径不变，但**证明显式路径模式今天也真被 GUARD C 对账**（要给出"配钉后拿 glob 解析集"那一条的**同种子双形读数**：拆包必响、正常必绿，⚠ 这正是 `251-r1` 推翻我 AC#1 原判据那一格——**目录形会照绿，所以"我在显式路径上"不等于"我被对账了"**）。判据＝两种形都要有"改前必红"的发数摆着，选哪一支都要具名凭据。
- [x] **AC#2 `core` 档里那行 `./internal/winsec/`（`scripts/portable-tests.sh:211`〔行号待验〕）与第五档的关系要定死**：⛔ 不许两档各自认领同一枚包而对账结果不同；要改 glob 就**同笔**动 `core_pin`，并给"漏计／双计"两形各一枚正控。
- [x] **AC#3 unknown 消息里那份 known 列表（从 4 名扩到 5 名）的口径裁清**：要么"消息文本属射程、不许扩"（则第五档与它不能并存，需改文案生成方式），要么"列表随档自动生成、扩名不算改动"。**这一条由编排者裁，不由腿自选**；腿只交两形的代价与现量。

## 禁区

- ⛔ **`ci.yml` 的形状不属本票射程**：动 workflow 触发档＝契约级（票 134 的 C+B 形状、见 [[wisp-ci-selfhosted-topology]]）⇒ **停手上报，不许"顺手加一步"**。
- ⛔ 不许把 GUARD C 的**双向对账**改成单向或跳过；⛔ 不许把 unknown scope 的 `rc=2` 改软；⛔ 不许为了让第五档"被调用"而往 `core` 档塞包名（那是把两档的钉搅在一起）。
- ⛔ 不许动 `internal/**` 的任何测试或产码（本票射程只有 `scripts/` ＋证据件）。
- `frontend/**`／`design/**` **既不读也不引**；grep/find 必写显式根（`scripts docs tools cmd internal .scratch`），⛔ 不用 `.` 当根。
- Git 纪律：只 commit 不 push；显式 pathspec；⛔ `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 交件要求

起手 `date -Iseconds`＋`git log -1`＋`git status --porcelain scripts/`（⛔ `scripts/` 此刻若有别的写腿在飞就停手上报——票 250／251 都在这面上）；**第一轮内先落骨架并 commit**；跑到第 100 轮前先把「门禁读数」与「判不动的地方」两节写满并 commit（⛔ 空着不算交件）。占位符自查用字符类尺在命令行跑、⛔ 尺的字面文本不落进被扫文件；计数尺一律 `| wc -l` 收尾。⛔ **一枚 AC 勾选框都不许碰**（翻勾归编排者）。写明**它推翻编排者上面哪一句**（上面三行现量与所有行号都是待验断言）。

## 排程

- 写面＝`scripts/`（此刻**已空**：`251-r1` 交件退出）。
- ⚠ 与 `197-r3`（写 `cmd/wisp`）无关，可并发；⛔ 但**不许与另一枚会跑整包 `go test` 的腿并发**（`248-v1` 正在跑 `cmd/wisp`，本腿不跑 Go 测试即可）。

## 编排者裁定与翻勾（2026-10-02 11:19:06+0800 现量，台账 `A529`）

⚠ **这一节的定性要说实话：三格是〔编排者现跑自勾〕，不是独立验收腿裁的。** 理由具名＝同一道题上三枚腿（`248-v1`／`248-v1b`／…）今天连续死于模型服务中断，我判断再派一枚 -v 腿等它的成本高于收益；⛔ 这不等于"独立验收没欠"——**欠的是一枚没跑的 -v 腿，不是这句话可以不说**。我自己在最新代码上重跑的两把尺（不引腿的读数）：

- **HEAD 那一发**：`bash scripts/portable-tests-selftest.sh all` ⇒ 逐字 **`27 case(s) ran, 0 assertion(s) failed`、rc=0**，收尾句 `GREEN - every seeded anomaly was refused, and the clean scope passed.`。⇒ 支撑 **AC#1／AC#2 成立**（载具 27 次调用全绿；新增的 19–22 号场景与票 251 的 18 次调用同跑同绿）。
- **改前那一发**：`git show d253703a^:scripts/portable-tests.sh` ＋ 同法取 `winsec-tests.sh`（`d253703a^` 现量＝`ed459d09`），`CARRIER_WINSEC_TESTS=… bash scripts/portable-tests-selftest.sh all …` ⇒ 逐字 **`27 case(s) ran, 9 assertion(s) failed`、RED**。9 条落点我逐条看过，全在本腿新增的四枚场景里（`FAIL exit 0, expected 2`＝GUARD 3 没按；`FAIL exit 0, expected 1`＝GUARD C 看不见那行；`the two tiers disagree…`＝两档各说各话；`core branch carries 1 private spelling(s)…`＝拼法分家；`the mutant seed took 0 line(s)…`＝枚内正控），**票 251 那 18 次调用在该基零红**＝基只动了本票那一枚变量。⇒ 支撑**这两格不是"文案变长"而是"牙装上了"**。
- **AC#3 我裁＝乙（列表随档自动生成，扩名不算改动）**，凭据是我自己 sed 的两行：`scripts/portable-tests.sh:246` 逐字 `echo "portable-tests.sh: unknown --scope=$mode (known: ${tiers// /, })" >&2`，而 `:192` 逐字 `tiers='core windows cli winsec census'`——⇒ **那份 known 列表今天本来就是 `tiers=` 的派生物，不是第二份名册**。甲那一支要么把 `:246` 硬写成四名＝**重造票 251 AC#4 刚拆掉的"两份名册各说各话"**，要么把第五档摘出名册 ⇒ `--scope=census` 从此 rc=1 且 `unknown --scope=winsec` 不再点名第五档。乙的代价写清在此：**新增一档时载具 case 17 会随档红一次**（要人看一眼，不静默），且本仓台账里那句"4 名"的旧串按**只追加**更正，不改原话。
- ⛔ **零产码改动**这一条我也复认了：`git diff d253703a..HEAD -- scripts/` 只有载具一枚文件动过，`ci.yml` 一个字节没动。

**两处残余（本票不勾死，具名归口）**：
1. **载具不在 CI 里跑**：`grep -c selftest .github/workflows/ci.yml` 现量 **0** ⇒ 上面 19–22 号那四枚场景今天是**笔记本尺**；要接进 CI＝动 workflow 触发档＝契约级（`A523` 那句"`ci.yml` 一个字不许动"仍生效）⇒ **默认不接**，接不接要 owner 一句话。
2. **GUARD 3 红句里 `scope=[${target}…]` 是拼出来的**，父脚本 `target` 丢尾斜杠时那句会跟着漂（腿的突变 D 实测：链仍 `rc=2`，**不会静默变绿**）。⇒ 我裁**本票不动 `winsec-tests.sh`**——失败方向是"红句不好看"而不是"该响的不响"，属可容忍残余；改它另立一手。

---

## 凭据归口（搬运笔，非新验收）

1. 本节＝**搬运**：凭据早已在盘，只是没归口到本票名下；⛔ 本节不新增任何验收读数、不改变本票任何一枚勾的状态。
2. 出处＝只读普查腿 `evidence-close-6`（件 `.scratch/wisp/probes/evidence-close/6/backfill-worklist.md`，§4.3 名册）＋编排者裁定。
3. ⛔ 本票**不该被本节视为已验收**。

归口的凭据路径（⛔ 本节不替它们写任何判语，判语在原件里）：

- `docs/evidence/s1/248-settings-write-path-v1.md`（`:68`／`:278` 两处出现 `254-r1b` 名号）
- `.scratch/wisp/probes/254/r1b/logs/` 那 14 份读数件，逐枚名册（本腿 `ls -1` 现量＝**14 份**，与编排者口径一致）：
  `carrier-all-at-HEAD.txt`／`mutA-chain-split-teeth.txt`／`mutC-guard3-scope-clause.txt`／`mutD-target-spelling-drift.txt`／
  `portable-tests-at-d253703a-parent.sh`／`portable-tests-at-f5f9cc34.sh`／`prefix-chain-explicit-split-audit-bites.txt`／
  `prefix-chain-widened-scope-loses-audit.txt`／`prefix-core-and-winsec-see-the-same-split.txt`／
  `prefix-core-dir-form-hides-the-split.txt`／`prefix-full-at-d253703a-parent.txt`／`prefix-full-at-f5f9cc34.txt`／
  `winsec-tests-at-d253703a-parent.sh`／`winsec-tests-at-f5f9cc34.sh`

存在性复量（搬运腿 `backfill-7a`，2026-10-06 现跑 `test -f`／`ls -1`／`wc -l`／`grep -n`，⛔ 未引任何一枚的判语）：
`docs/evidence/s1/248-settings-write-path-v1.md` **在盘、279 行**（与编排者给的 279 一致），`:68`／`:278` 两行**都在盘上逐字含 `254-r1b`**；
`.scratch/wisp/probes/254/r1b/logs/` **在盘、14 份**——以上两处存在性**全部对得上**。
⚠ **一处计数差具名登记（⛔ 不照抄、不替换、不自己补凭据）**：编排者那句写作"`:68`／`:278` **两处**出现 `254-r1b` 名号"，
本腿 `grep -n '254-r1b' docs/evidence/s1/248-settings-write-path-v1.md` 复量命中 **3 行**＝`:68`／`:179`／`:278`。
⇒ 给定的两行**逐字成立**，多出的 `:179`（同件里另一句并发声明）**不列进本票凭据**，只把这一枚计数差报给编排者裁。
