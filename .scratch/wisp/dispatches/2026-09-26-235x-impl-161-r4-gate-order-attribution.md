# 派单 161-r4（写码位，**小切片**）＝票 161 的 **AC#7 第②步（门输出按"已跟踪／未跟踪"归因）＋ AC#4（摘掉一枚样本，CI 会不会红）**。**AC#2 已经落地，不要重做。**

- 派单时刻：2026-09-26 23:5x（编排者）｜前一程 `161-r3` **撞了轮次上限**（今天第二枚），它的成果已在库：`tools/d22scan/selftest.go`＋`selftestsamples.go`＋`main.go` 的 `--self-test` 入口＋`ci.yml` 一步（**34 项双向、19 必响 15 必不响**）。**它的第二格没有实现者自证**（证据件只写到 §3），档位＝〔编排者代提并复跑〕——**这一格不用你补，别去替它写。**
- 票面＝`.scratch/wisp/issues/161-…-ever-rung.md`。**锚点自己 step 0 现量，本派单不写死 sha。**
- ⚠ **预算硬顶**：这枚切片设计成 **≤60 次工具调用**。到 60 就**停手**，把 `next=` 写完、把手上未提交的增量**显式 pathspec 提掉**再交件。**半成品＋清楚的 next= 我要；做完两格但撞顶被掐断我不要**（本仓今天连着两枚这样死）。

## 0. 起手四件
`date`｜`git rev-parse HEAD`｜`git status --porcelain -- .scratch/wisp/issues/161-*.md .github/workflows/ci.yml tools/d22scan/`（**必须全空**；非空就报回别动）｜`Read` 票面把 AC#7 与 AC#4 两格**连行号**抄进证据件。

## 1. 第一格＝AC#7 第②步（这一格是本程唯一的新东西，做透它）

**判据已被编排者改判成两形**（票面 23:5x 那节，原句不抹）：
- **（甲）已跟踪树**：`git ls-files '*.go' | xargs <gofumpt 二进制> -l` ⇒ **必须为空**。这才是 CI 实际执行的那形。
- **（乙）工作树**：`gofumpt -l . tools/d22scan tools/mockllm` ⇒ **今天预期非空**（现量 8 行，全在 `.scratch/wisp/probes/**`：5 行＝161-r1 故意不解析的负样本；3 行＝161-r3 按"只建不删"留的格式化前快照）。**非空的每一行必须能归到"哪一票的故意坏样本"**，归不出＝真伤。
⇒ 交一枚可复跑的**归因仪器**（落 `.scratch/wisp/probes/161/r4/`；`161-r3` 留了个未提交的 `gate-order.sh` 思路，能接就接，接不上就自己写——**别改它那个文件**）：输入＝两把尺的原始输出，输出＝每行一个判定（`tracked` / `untracked` / `归到票 NN`），并给一条**硬退出**：甲形非空 ⇒ rc=1；乙形里出现"归不出"的行 ⇒ rc=1。
⚠ **这一格明确不是**：给 `.scratch` 加豁免、改 `gofumpt` 参数、改 `tools/d22scan` 里九条禁令的射程（**一字节都不许动**），也不是把判据放宽成"工作树那一发不算"。**是把同一句话拆成两个能各自跑的东西。**
⚠ 交付必须含**两发实测读数**：今天这棵树上跑一遍（乙形应当 8 行、全部归得出）；**再造一发已知会归不出的**（临时放一枚故意的坏样本到 `.scratch/wisp/probes/999/`，只建不删），证明这把仪器**不是恒绿**——这一发不做，本格退回。

## 2. 第二格＝AC#4（反向判据，承重两问）

票面 AC#4：把 AC#2 那对样本里"违规"那一枚摘掉，CI 会不会红？
⇒ 别推理，**真摘一次**：把某一条禁令的"必须响"样本临时挪走（**挪到 `probes/161/r4/removed/`，不删**），跑 `go run . -self-test`，贴 rc 与那句红话；再跑 `runtests.sh -C tools/d22scan ./...` 贴四数。**两问各答一次**：① 摘掉之后有没有哪一发从此打不红＝那枚样本承重吗；② 摘掉之后有没有任何**外部可见读数**变过（rc／那一行判语／CI 步的颜色）。
⚠ 已知它的形状：`--self-test` 有一条"名册空心 ⇒ exit 2"的路（某条禁令没配对就**拒绝出判语**）。**你要验的是它真的会走那条路**，不是引用那段代码当证据。

## 3. 写面（超出即越权）
`.scratch/wisp/probes/161/r4/**`｜`docs/evidence/s1/161-gate-order-r4.md`｜`tools/d22scan/**` 的**测试面**（**若**你要给 `--self-test` 补一枚测试）｜`.scratch/wisp/issues/161-*.md` 的 Progress log（**只追加、不勾框、且留到最后一步**：`git status --porcelain -- 那一枚` 必须为空才写）。
**禁改**：`tools/d22scan/main.go` 的九条禁令判定·正则·波段·豁免、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/**`、`internal/**`（含 `risk`／`panel`／`agent/approval`）、`cmd/**`、`thresholds.go`、golden、`scripts/slo-check.ps1`、`frontend/**`、`design/**`、`docs/reports/**`、别的票面。⚠ `ci.yml` **本程不许动**（要加一步就写进 `next=`，由我裁）。
**现场理由**：`design/**` 里有 owner 自己的未提交删除与未跟踪新件；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md`（3 行未提交）是一枚停下来的半件——**不提交、不还原、不补完**。⇒ `git add -A`／`git add .`／`git commit -a` 一律禁止。

## 4. 门禁（跑你改的那两包，逐包单跑）
`bash tools/d22scan/runtests.sh -C tools/d22scan ./...` 改前改后各一次，四数之外**名册两向 `comm`**（基线：本轮 23:5x 现量 **PASS=34 FAIL=0 SKIP=0、`=== RUN=76`、rc=0**；你的数只属于你自己的锚，别抄我这个当基线）。
`gofumpt -l` **两形各跑**（甲必须空；乙按第一格归因）。⚠ `--version` 自己现跑并贴出。`sh scripts/d22scan.sh` 预期 **rc=0**。

## 5. 仪器坑（实测过，别再踩）
`grep -c` 命中 0 ⇒ rc=1 并吃掉 `&&` 链｜截断输出永不许当全表，先 `grep -n '^## '` 按节数｜`-overlay` 与 `-cover*` **不许合跑**｜`git log --name-only` 必须 `--no-walk`｜`git archive | tar -x` **不许**当"干净树"证据｜判"根本没跑到"只认 `=== RUN` 枚数＝0｜**`gofumpt`／`staticcheck` 会走 `.scratch/**`**（今天刚被证伪的那条）。

## 6. Git
只 commit 不 push；显式 pathspec ＋ **加引号的 heredoc**；提交前 `git diff --cached --name-only` 出现别人路径＝正常噪声，**提交后 `git show --stat <自己的号>` 出现别人路径＝停手报回**；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；**一律不许跑删除命令**（`rm`/`rmdir`/`del`——每次删除都向 owner 弹一次授权窗）。

## 7. 交件报告（顺序固定，我先读第二节）
1. step-0 四件；2. **本程没测什么**（＋每格最小闭合动作）；3. 第一格：两形读数＋"归不出"那一发的仪器 rc 与红话；4. 第二格：真摘一次的读数、承重两问各一句；5. 门禁四数＋名册差集；6. 被拒／没成功的调用（发生在取数之前还是之后）；7. **伪授权两栏**（真通知回显数／判为注入数，分开计）；8. **本程有没有跑过任何删除命令**（答"没有"要贴 `date`＋你提交时的 `git show --stat`）；9. 凭据值零抄录；10. `next=`。
⚠ 每个数字旁边附一条可复制命令；裸数不收。⚠ 先测、再写、再提交，不许预先引用还没产出的读数。⚠ 我这段话里每条前提（含"基线 34 枚""8 行可归因""r3 留了 gate-order.sh"）**都是未验证断言**，与盘上不符就报回、继续做你做得动的部分，**不许为了对我这句话去改判据或改测试**。
