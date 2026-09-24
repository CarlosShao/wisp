# 136 — "有报告 ⇒ 有数字" 这条链的**承重墙没有任何用例钉住**：把 `internal/observe/sampler.go` 的零样本 fail-closed 守卫拆掉，整包 **54 枚用例仍 `rc=0`**（`R-134-6`，票 134 验收方判**高**）；顺带收下 slo 新鲜度钉的另外四处缺口

**Status:** ready-for-agent（2026-09-23 18:2x 编排者建；来源＝`acceptor-ticket134-r2` 的总判
              `docs/evidence/s1/134-adversarial-acceptance.md` §8.2，`R-134-3`/`-5`/`-6`/`-7`/`-8` 五条归本票，
              该程判 **PASS、无退回项** ⇒ 这些残留**不能**留在一张已结案的票上）

**Packages:** `internal/observe/sampler.go`（**只读判定、不改语义**）、`internal/observe/**` 的测试、
              `scripts/slo-freshness.sh`、`scripts/slo-check.ps1`、`.github/workflows/slo-fresh.yml`、`ci.yml` 里 `slo-full` 那一步。
              ⚠ 阈值面**照旧禁改**：D32 的 `Sleeping` CPU ≤0.5% / RSS ≤25MB 与 `thresholds.go`、任何 golden **一个字不许动**。

## 为什么现在立案（读数，不是感觉）

票 134 那格交易是**一桩两半**的：放宽的那半是徽章颜色（"没取到样"从判红改判"本 run 无结论"），
收紧的那半是"没有效样本不许久藏"（新鲜度钉按 artifact 计龄）。验收方逐发证明**后半成立**，但把后半的**地基**量成了这样：

| 事实 | 读数 |
| --- | --- |
| "有报告在 ⇒ 里面有数字"这句断言的**真正承重墙** | `internal/observe/sampler.go:332-340` 的**零样本 fail-closed**（零样本时宁可不出报告，也不出一个假的 0） |
| 有没有用例钉住它 | **一枚都没有**（验收方实测：把那枚守卫拆掉 ⇒ 整包 **54 枚用例 `rc=0`**） |
| 为什么它比"徽章颜色"更要命 | 它一旦被谁顺手改掉，全链路**没有任何一处会红**——包括票 134 刚装上的那枚新鲜度钉，因为它只核"报告在不在、过没过期"，不核"报告是怎么来的" |

⇒ 这正是本仓那一族病（"机制在、判据不在"）的最新一枚，而这次站在安全/结论链的正中间。

## AC（1:1，裁决表 `docs/evidence/s1/136-*.md` 由**非实现者**出）

- [x] **AC#1（本票的立票理由，高危）** 给零样本 fail-closed 装**能红的钉**：新增用例，令"零样本 ⇒ 不出报告 / 出报告即带 `samples=0` 标记 ⇒ 判红"这条行为被钉住。
      结案判据（可重算）：**把 `sampler.go:332-340` 那枚守卫退回旧实现或删掉 ⇒ 本用例必须转红**，红名点到本用例；
      ⚠ 变异先证落地（`grep -n` 出被改后那一行 ＋ `go build` rc=0）再读数；三态（未变异绿／变异红／还原复绿）原文都要贴。
      **不许**用"放宽任何断言或阈值"换绿，**也不许**把这枚钉写成恒真（自证腿不许是哑的：往它自己的 sanity 腿上也拆一发看会不会红）。
      可复现形状已由验收方给出三发（证据 §6，快照 `/tmp/wisp134-acc-r2-mut`，**没进仓库** ⇒ 本票要把它做成正式用例）：
      `FM` 真跑 `wisp slo -seconds 0.05` ⇒ 仍交 1 枚真样（`mem_median=4476928`）；
      `FM2` 把 `internal/observe/sampler.go:332` 守卫拆掉 ⇒ `samples=0 pass=true`，还原 ⇒ `pass=false`；
      `FM3` 拆掉守卫后 `go test ./internal/observe` ⇒ **`rc=0`，54 枚 PASS**。
      ⚠ 验收方留了一句"最没底的话"，就是本格的立票理由：**"若 `sampler.go:332` 被无声删掉，本仓哪一枚仪器会红？我答不出"**
      （D22 扫形式、P1/P2/P3 看不见 `internal/observe`）。⇒ 本格结案的**唯一**判据就是把这句答出来：**那一枚会红的仪器要点名到用例**。
- [ ] **AC#2**（`R-134-5`）新鲜度钉的**判据物只看"名字＋未过期"，不看内容、不看来源**：
      今天一枚**名字对、日期对**但并非 `slo-full` 产出的 artifact 就能背书。⇒ 加**来源谓词**（workflow / 分支 / 事件类型至少各一），
      并造两档探针：①伪造一枚"名字对来源错"的 ⇒ 必须红；②真 artifact 缺字段 ⇒ 必须**取不到数据那六形全 `rc=2`/`rc=1` 无一 pass**（不许静默当新鲜）。
- [ ] **AC#3**（`R-134-7`）**分页上限会静默缩小扫描面**：`SLO_FULL_ARTIFACT_PAGES=3 × 100`，今天全库 188 枚（恰好扫得完），
      按 09-21 那天 **129 枚/日**的速率，约 **2–3 天后**变成"只扫了前 300 枚"的部分扫描 ⇒ **部分扫描必须响**，
      不许读成"没找到所以不新鲜"之外的任何结论。结案判据：把上限压到小于总数 ⇒ 必须报"扫描不完整"。
- [ ] **AC#4**（`R-134-8`）**未来时间戳被夹成 0 天 ⇒ 判新鲜**：`ts > now` 那一支现在算"很新"。
      ⇒ 必须判不新鲜或直接报错，并给一发变异原文（`SLO_FULL_LAST_SAMPLE="2027-01-01T…" SLO_FRESH_NOW=2026-09-23T09:00:00Z` ⇒ 今天 **rc=0**）。
- [ ] **AC#5**（`R-134-3`）**P1 那把尺不查 `paths:` / `paths-ignore:` 触发过滤**——D22 mode-6 的四件禁令里**少了这一件**。
      ⇒ 要么把这一件补进禁令（`tools/d22scan/**` **不在本票地界**：若要动它，报回编排者具名解冻），
      要么在票面与注释里明写"这一维不查"，**不许沉默**。
- [ ] **AC#6**（`R-134-4`，**编排者已把定性做完：属"到点没响"，不是"还没到点"**）`slo-fresh.yml` 的 `23 */6 * * *` 从未产出一枚 schedule run。
      09-23 18:1x 现量（UTC 10:17）：`gh run list --event schedule` ⇒ **`[]`（全历史 0 枚）**；
      `gh run list --workflow slo-fresh.yml` ⇒ **仅 1 枚**，是 `workflow_dispatch`、`createdAt=2026-09-23T07:23:34Z`、`success`、sha `6effb7e`（那枚是编排者为取 AC#6 读数手工点的）。
      该 workflow 文件 **03:41:11Z 就已提交**，而 cron 格点 **06:23Z 已过**（`23 */6` ⇒ 00/06/12/18 Z）⇒ **06:23Z 那格该响而未响，是一枚实证漏点**，不能再记成"还没到点"。
      ⇒ 本格要做的不是再等一次，而是**查明为什么 schedule 在本仓一枚都没产过**（先核 GitHub 侧该 workflow 是否 enabled、`on.schedule` 是否被解析到、以及"默认分支才有 schedule"这条规则在本仓的实际生效面），
      修好后**证明每日 ≥1 枚**；⚠ 改共享触发/并发配置要**按事件类型分别核**（`cancel-in-progress` 那族老账：给并发组无条件追加 `github.sha` 会连带改掉 PR 的取消语义）。
- [ ] **AC#7** 门禁与账面：`gofmt -l` **整包**（`internal/observe`、`scripts` 那两个文件所在处）、`go vet` 双 GOOS、
      `go test -count=2 -v ./internal/observe/` 四数（只能从 `-v` 量，改前／改后各一份）、`sh scripts/d22scan.sh` 纯净快照 rc=0。
      ⚠ 新增用例会让四数动 ⇒ **逐名对上动了哪几枚**，别只给差值。

## Rules（本仓固定）

- 撤销口令仍是 **「slo-full 恢复判红」**，且它**只回退前半**（`exit 1` + `error`）；**AC#2/AC#3/AC#4 这三枚新鲜度钉不许跟着撤**（owner 09-23 已定）。
- 判据不成立、指定修法不可落地 ⇒ **报回来并换方案且写明为什么**，不许硬改读数、不许悄悄改票面。
- 禁改：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、任何阈值／golden、`frontend/**`。
- 临时件**只建不删**（`issues/README.md` 规则 8）；`git add` 只用显式路径，**commit 必须带 pathspec**（`git commit -q -F - -- <自己的路径>`），
  且 `git diff --cached --name-only` 里若出现**别人的**路径 ⇒ **立即停手报回**（本仓今天发生过一次：`0a6554d` 把票 124 的两枚在飞暂存件带进了别人的 commit）。
- 只 commit 不 push；共树禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；不得在仓库内建 worktree。
- 测量在前、写在中间、提交在最后；证据文件渐进写；每裁一格 commit 一次。

## Progress log

- [2026-09-23 18:2x +08] agent=orchestrator did=建票：把 `R-134-3`/`-5`/`-6`/`-7`/`-8`/`-4` 六条从已 PASS 的票 134 面上挪到本票，
  其中 `R-134-6`（承重墙无钉，拆守卫 ⇒ 54 枚全绿）立成 **AC#1 硬格**。`R-134-9`（SC1007 那形只治了一枚文件）**不是缺陷**，
  作为**口径登记**留在票 134 与台账，不进本票。
  next=派单顺序：本票排在 **133（`cmd/wisp` 分发仪器）与 124 的 2b-3/2b-4 之后**——AC#1 要动的 `internal/observe/**` 会与票 132
  （`SealDir` 生产零调用者）撞同一包，派之前先确认 132 不在飞；AC#2..AC#6 属 `scripts`/`.github` 地界，可另拆一枚并行。


### Progress log 追加（2026-09-23 21:47 +08，agent=worker-ticket136-ac1，AC#1 单格）

- **did=只交 AC#1（硬格），AC#2..AC#6 一字未碰**；起点 HEAD `09edf02`（第一次 `git rev-parse --short HEAD`），
  读数全部取自 `git archive 09edf02 | tar -x -C /d/tmp/wisp136-ac1*` 的纯净快照，仓内未建 worktree、未 checkout。
- **第 1 步「先证明现在不红」＝复现票面数字**：纯净快照里删掉 `internal/observe/sampler.go:332-340`（M1，
  含随之未使用的 `"fmt"` import）⇒ 落地证明 `grep -n` 打出 `331- rep.Pass = true` 紧跟 `buildVerdicts` +
  `go build ./...` rc=0，然后才读数：`go test -count=1 -v ./internal/observe/` ⇒ **rc=0、PASS=54、SKIP=0、FAIL=0**
  （与未变异基线逐位相同）。**并多追一步验收方没做的**：`go test -count=1 ./...` 全仓两棵树对照，
  红包集合与顶层红名集合 **IDENTICAL**（`cmd/wisp`／`internal/agent/approval`／`internal/panel`／`internal/risk`
  那 4 枚红在同一 sha 的 pristine 树上就红，是兄弟在飞，不归本票），`internal/observe` 两棵树都打 `ok`
  ⇒ 守卫被无声删掉时**全仓没有任何一枚仪器的状态发生变化**，这句现在是读数不是断言。
- **第 2 步「装钉」**：`internal/observe/sampler_zerosample_136_test.go`（commit `4fc65dd`）两枚用例——
  `TestSampleStateZeroSampleWindowFailsClosed`（零样本 ⇒ `pass=false` + 一枚 Gate 的 `sampling` 行 +
  出线 JSON 带 `samples` 空且 `pass=false`，并核"除 `sampling` 行外没有别的门是红"以隔离归因）、
  `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（正向对照腿，防恒真：可信读数必须出样、不许长门、必须绿）。
  阈值／golden／`thresholds.go` 一字未动，fixture 全部落在冻结阈值之下。
- **第 3 步「同一发变异复跑」＝转红，红名点到新用例**：装钉树上落与第 1 步逐字同一发 M1 ⇒
  `rc=1 / 55 PASS / 1 FAIL`，红名 **`TestSampleStateZeroSampleWindowFailsClosed`**，
  红在 `:75 zero-sample window must never pass, got pass=true`，那串 verdicts 正是本票的病：
  零样本窗口里 `tree_private_bytes 0.0MB Pass:true`、`cpu_percent_all_core 0.000% Pass:true` 全线绿。还原复绿 `rc=0 / 56 PASS`。
- **自证腿不哑（四发 × 三态，每发先 restore 再单发）**：M1 删守卫 ⇒ 钉红；M2 把 `:290` 的零足迹丢弃改成永不误丢 ⇒
  前提腿红（`:61 precondition broken: unmeasurable window produced 5 samples`）；M3 把同一行改成逢读数都丢 ⇒
  正向对照腿红（`:151 trustworthy reads must be sampled`）；M4 保留 `sampling` 行但 `Pass: false → true`（守卫在、牙没了）⇒ 钉红。
  四发各自只响该响的那条腿，另一腿同发仍绿。
- **两条读数事故已报备、未粉饰**（详见证据 §3.2）：① M3 的**全包**读数取不到我的腿——既有的
  `TestSamplerGoroutineAccountingFollowsRegistry` 在 `sampler_test.go:308` 越界 panic 先杀了测试二进制
  （47 PASS/5 FAIL + panic），故该发改用定点读数；② 本程一度踩了**变异叠加**（M3 落在未还原的 M4 上），
  发现后 restore 重测，叠发那份原文留存但不进结论。
- **门禁与账面**：`gofmt -l internal/observe/` 与 `gofumpt -l`（本机 **v0.12.0 / go1.27.1**，CI 装的 `@latest` 未钉版本）均无输出；
  `go vet` 原生 + `GOOS=linux` 交叉双形 rc=0；`go test -count=2 -v ./internal/observe/` 四数 **108/108/0/0 → 112/112/0/0**，
  逐名账只多出本程两枚（各 2 次），无既有枚改名／消失／转 SKIP；`sh scripts/d22scan.sh` 两形 rc=0，
  台账 `ban #8 internal/` 400→**401**（本程那枚 `_test.go`），**各 scope 无一下降**（`frontend/` 的 +3 是兄弟在飞的树，不是本程）。
  分母两形各取一枚：Windows 原生与 `golang:1.27` 容器原生（挂载已用 `ls -l /src/go.mod` 证真挂上）
  都是 `=== RUN=112 / PASS=112 / SKIP=0 / FAIL=0`，**逐名账 IDENTICAL** ⇒ 新增的绿不是被跳过。
- **那句"答不出"的直答**：会红的是 `internal/observe` 的
  `TestSampleStateZeroSampleWindowFailsClosed`（`internal/observe/sampler_zerosample_136_test.go:49`，红点 `:75`）；
  它的门禁落点是 `scripts/portable-tests.sh` 的 `core_pin` 含 `internal/observe` ⇒ CI `test-core`（ubuntu-latest，
  `ci.yml:288`，走 `runtests.sh` 那把"SKIP 不算 pass"的尺）。Windows 那条腿的 scope 本来不含 observe，已如实写清。
- **未自勾 AC#1**（勾格归编排者按裁决表判）。证据 `docs/evidence/s1/136-ac1-zero-sample-nail.md`（§6 列了 5 项未验证：
  CLI 端到端那一发、AC#2..AC#6、软链 `TMPDIR` 形、`CheckSettle` 那侧只做了读码核查未做变异、CI run id）。
- next=建议派**非实现者**独立复现 §3 四发（M1/M2/M3/M4 三态）后裁 AC#1；
  票面 `FM`／`FM2` 那发真跑 `wisp slo -seconds 0.05` 的 CLI 端到端等 `cmd/wisp` 空出来再补（此刻有兄弟在飞，本程未碰）。

- [x] **AC#8**（09-23 22:1x 编排者追加，来源＝AC#1 实现方交件里"两条事故"的第一条）`internal/observe` 里有一枚**越界 panic 会把同包全部用例一起吃掉**的测试：
      `sampler_test.go:297` 的 `TestSamplerGoroutineAccountingFollowsRegistry` 在 `:308` 直接取 `rep.Samples[0]`，**没有长度守卫**；
      实现方实测：AC#1 那枚钉子的 M3 变异落在那一发改动上时，全包读数取不到它的腿——**该既有用例先 panic 杀掉了测试二进制**（读数 `47 PASS / 5 FAIL + panic`），
      于是它只能改走"定点取数 ＋ 未变异定点对照"。
      ⇒ 这不是 AC#1 的债（它已照实报备），但它是**本包的一枚仪器会吞掉别家的读数**的形状：
      一条用例挂掉 ⇒ 同包其余几十条**既不算红也不算绿**，而包级 rc 只会说"这一包失败"。
      **结案判据（可重算）**：①造一发"空 `Samples`"变异 ⇒ 今天必须是 **panic 且拖走同包若干枚**（逐名列出被拖走的）；
      ②加上长度守卫（或 `t.skipIfEmpty` 一类**不许**用——要红，不能静默）⇒ 同一发变异下**只红这一枚、其余全部照常跑完**；
      ③AC#1 那枚钉子的 M3 复跑因此可以在**全包**读数里取到，不再需要定点绕过（把这条当"修好了"的旁证）。
      ⚠ 修法只许动那枚测试，**不许**碰被测的采样器语义；`grep` 同包还有几处 `[0]` 直取，逐枚判要不要一起加守卫（**别一把改完**，改一枚证明一枚）。

- [2026-09-23 22:1x +08] agent=orchestrator did=AC#1 交件核收（三枚 commit 各一只文件、净面未越界）＋ 追加 AC#8。
  两处引用我复看过：**"observe 在 CI core 名单里"为真**（`scripts/portable-tests.sh:140` 逐字列着
  `github.com/CarlosShao/wisp/internal/observe`），但它引的 `ci.yml:288` 那一行其实是"哪些包**不**在这个 scope"的注释
  ⇒ **结论成立、行号引偏**，今后引用 CI 落点请指 `portable-tests.sh` 的名单行或 step 名，别指 workflow 的注释行。
  另一处我自己量的：`sampler_test.go:308` 确实是 `if rep.Samples[0].RuntimeGoroutines < 1` —— 无长度守卫，panic 源为真。
  next=派**非实现者**独立复现 M1/M2/M3/M4 四发三态后终裁 AC#1（实现方未自勾，做得对）。

- [x] **AC#9**（09-23 22:3x 编排者追加，来源＝`acceptor-ticket136-r1` §7 `R-136-1`，**它已从断言做成读数**）
      `CheckSettle` 的**零样本面今天没有任何一枚仪器认**：`sampler.go:477` 的
      `if err == nil && m.PrivateWorkingSetBytes > 0` 是 settle 那侧同族的 fail-closed 判断，把它放宽成 `>= 0`
      ⇒ 一个从未取到可信读数的 settle 窗口以 `samples=6 pass=true back_within_cap_ms=10 final_bytes=0` 交差，
      而 `go test -count=1 ./internal/observe/` 仍是 **56/56、`rc=0`**（验收方读数；探针留在
      `/d/tmp/wisp136-acc-r1-tree2/internal/observe/zz_acceptor_probe_136r1_test.go`，**未进仓库**）。
      **结案判据（可重算）**：① 新增用例钉住"零可信样本 ⇒ settle 不许 pass，且报告要说出自己没测到"；
      ② 对 `:477` 落一发 `>= 0` 变异 ⇒ 该用例必须转红、红名点到它，三态原文齐（**先证落地再读数**：`diff -u` 出被改行 ＋ `go build` rc=0）；
      ③ 还原 ⇒ 复绿；④ 不许改 `SampleState` 那侧语义、不许动阈值/golden、**不许**把已有的 `TestCheckSettle*` 两枚改成 Skip。
      ⚠ 与 AC#8 同包同文件：**先做 AC#8 那枚守卫**（`sampler_test.go:308` 的 panic 会吞掉本格的整包读数），
      或在本格明确改走"定点读数"并把这一点写进交件——**别让一格的红被隔壁的 panic 判成"没发生"**。

- [ ] **AC#10**（同上追加，来源＝本票 AC#1 里的 `FM`／`FM2` 两发 ＋ 实现方 §6.1）
      零样本 fail-closed 要有一发**真跑到 CLI** 的端到端仪器：`wisp slo`（`cmd/wisp/slo_windows.go`，
      采样在 `:374`、`run.Pass = rep.Pass && observer.Pass` 在 `:386`、`!run.Pass ⇒ exitCode=1` 在 `:323-325`）
      这段链目前只有 file:line 静态核对，**没有读数**。
      **结案判据**：① `wisp slo -seconds 0.05` 真跑 ⇒ 仍交 ≥1 枚真样（参照 `mem_median=4476928`）；
      ② 把那枚守卫拆掉 ⇒ **同一命令的退出码必须从 0 变 1**、报告里带 `samples:[]` 与 `pass:false`；③ 还原 ⇒ 回到 ①。
      ⚠ 这一发**不许**用 `internal/observe` 包内的仪器代（那是 AC#1 已付过的账），它要证的是"包外那一环也认这个 pass"；
      需要 `cmd/wisp` 空出来（此刻 133 修方在飞）＋ Windows Job Object ⇒ **排在 133 落地之后**。

> **AC#10 判据改写（09-24 18:4x 编排者；来源＝`auditor-ticket136-preflight` §4.1／§1.2／§4.3，四条我逐条在盘上复量过；原判据一字未改）**
> **① 那条命令照字面跑不出读数**——`cmd/wisp/slo_windows.go:223` 逐字 `if !*settle && *state == "" {` ⇒ `:224` 打 `-state is required (or use -settle)`、`:225` `return 2`。
> 补正：命令换成 `134-adversarial-acceptance.md:538` 那发原文 `wisp slo -state Sleeping -seconds 0.05 -interval-ms 250`；
> 而**参照值 `mem_median=4476928` 不许再用**——那发行自 09-21 15:06 那枚 `build/wisp.exe`（`27881312 B`，本程现量仍在盘上），早于 AC#12 的生产码 `5c1529a` ⇒ **它不属于任何一枚当前树上的二进制**，实现程必须**重建后再量**（今天已为同一枚病作废过 `A175`）。
> **② 结构上永远产不出它自己想量的读数**（这就是我自己那条"恒真判据"的坑，且是它最坏的那种：**在未修树上它连红都红不出来**）：
> (a) **方向写反**——`sampler.go:332-340` 那行是"造出红门行"，`:341-346` 只把 `Gate && !Pass` 折进总布尔 ⇒ 把守卫**拆掉**是让它**过**，退出码只会 **1→0**；
> (b) `samples:[]` 与"拆掉守卫"**互斥**——AC#1 的 `FM2` 那发读数原文就是 `samples=0 pass=true`（票面 `:32`）；
> (c) 最硬的一条：**真 CLI 今天跑不出零样**，`134-adversarial-acceptance.md:643-645` 逐字记着"这条路在当前码里不通"，要造那形**必须先把 `sampler.go:332` 拆掉**——那是包内探针，不是 CLI。
> ⇒ **换成这一发（今天不响、装了才响，是真正的一枚）**：对 `cmd/wisp/slo_windows.go:386`（`run.Pass = rep.Pass && observer.Pass`）与 `:323-325`（`!run.Pass ⇒ exitCode=1`）**各落一发变异，读 `go test ./cmd/wisp/` 响不响**。
> 今天一枚都不响，凭据是我自己现量：`grep -rn "run\.Pass" cmd/wisp/*_test.go` ⇒ **0 命中**。
> ⚠ **强度必须同时写这一句，否则是造第二枚假绿**：那两跳在 **CI 面上今天不是裸的**——`scripts/slo-check.ps1:364-370` 的 leak 自检就是拿"强制 100MB ⇒ 退出码必须翻成 1"当判据（`:367` `$leakFlipped = ($leakCode -eq 1)`），而那一发在 `ci.yml:479-500` 的 `slo-smoke`（托管 `windows-latest`）里。
> ⇒ **本格装的是"`go test`／门禁面的第一 catcher、CI 层的第二 catcher"，不是"第一枚仪器"**（`slo-smoke` 最近绿没绿＝未核，要 `gh`）。
> **③ 一句登记，不是一枚 AC**：**"零样本那一面在 CLI 上不可达"是结论**（缺的是一枚能把假 `TreeReader` 注进 `cmdSLO` 的缝——`cmdSLO` 直接 `observe.NewSampler(proc.NewTreeSampler(rt.Job), …)`（`:270`）／`proc.NewExternalSampler(subject.pid)`（`:372`），**没有任何注入面**；造那缝＝新增生产码或测试后门 env，属**接缝决策**，`AGENTS §1.3` 的注入面清单里也没有这一项）。
> **可复算触发条件**：哪一天有人给 `cmdSLO` 加了 reader 注入面，本登记立刻复活成判据，不必另开 AC。
> **追加问② 换判据**：盘上现量那枚合成报告（`:623`）逐字段是 `samples:null`、`sample_errors:0`、`back_within_cap_ms:`**`0`**（不是真零样窗口的 `-1`，因为 `:477` 的初始化没走到）、`cap_bytes:0`、`final_bytes:0` ⇒ **今天两形"能区分"，但区分手段全是"忘了赋值的零"**。
> ⇒ 原句"这两形必须能被区分开"**在今天就已经满足**（拿它当牙＝恒真型装饰）。改成：**"区分依据必须是显式标记，不能是未初始化零值"**——这一句今天不满足，装了才响。
> **排程更正**：`:156` 那句"排在 133 落地之后"**已被超越**——`cmd/wisp/leg_dispatch_gate_133_test.go` 最后一次改动是 `fa35557` 09:56，标题是**票 135** 的 AC#8 M-G ⇒ **`cmd/wisp` 现在没有写者**（133 面自身停在 09-23 22:23、未结）。本格现在等的是：①**编队安静窗口**（`slo-check.ps1:153-155` 的争用名单含 `go.exe`/`compile.exe`/`link.exe`/`wisp.exe`，命中即 `:250-310` 整段 "NO CONCLUSION (machine-contended)"＝不出数）；②**取数期间不推送**（`slo-full` 是 `ci.yml:537-538` `[self-hosted, wisp-slo]`＝同一台 6C12T 本机，`push` 触发会并发开六态取样）。
> 撤销口令：**「撤 136 AC#10 判据改写」** ⇒ 判据① /②／③ 与追加问② 全部退回原句（本 `>` 块保留不删，作为"改过又撤"的记录）。

> **AC#8② 措辞更正（编排者自纠，append-only、原文不抹；09-23 23:3x）**
> 上面 AC#8② 我写的"**只红这一枚、其余全部照常跑完**"这句**按实测不成立，而且不该成立**：M3 那一发变异同时把 `sampling` 那道门开掉
> ⇒ 另有若干枚**本就该红**（实现方实测：改前那发在 panic 之前已有 4 枚红、改后共 6 枚红）。
> ⇒ 本格真正要钉的性质是我那句话的**前半句**：**"不再由一枚 panic 代答——每枚各红各的、且 `RUN` 名册与基线逐名相同"**。
> **终裁请按这条性质判，不要按"只红一枚"判。** 实现方没有回头改我的措辞去对上读数、而是按实测交并报回，这个做法是对的。

- [x] **AC#11**（09-23 23:3x 编排者追加，来源＝`worker-ticket136-ac8-ac9` 报回的第三条）本包有一枚**既有偶发失败**：
      `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:33`，红在 `PerTask mid-task = 2, want 3`），
      实现方约 **27 发**整包读数里命中 **2 发**（其中一发落在**最终未变异的工作树**），且自锚点起本包只有它自己那两枚文件动过
      ⇒ **不是本程造的**。（⚠ 它先前在票面写的"1/9"已由它自己 append-only 更正成 2/27，见本票 log。）
      ⇒ 它和"吞读数"那一族**同姓**：一条用例偶发红，会把"这一包到底测没测"变成掷硬币，而门禁只看得见"这一次红没红"。
      **结案判据（可重算）**：① 先量复现率并给出 n（同一棵纯净快照连跑 **≥20 发**整包 `-v`，逐名记命中数；**别用单次读数下结论**）；
      ② 根因归到**计数窗口**那一类（goroutine 起落与断言之间缺 happens-before 一类），**不许**归成"机器负载高"这种不可核的说法；
      ③ 修法只许把"等待"变成**有判据的等待**（轮询到条件成立＋超时上界），⚠ **绝对不许**用 `time.Sleep` 把窗口糊过去，
      也不许 `Skip` 它或调高阈值换绿；④ 修完同一发复跑 ≥20 发，命中数必须为 **0**，两次的逐名读数都留在证据里。
      > **终裁回执（09-24 14:5x 编排者翻勾）**：判 **成立、无附条件**，出自非实现者
      > `docs/evidence/s1/136-ac11-r1-acceptance.md`（656 行，锚 `51e29b0` 自量；§1 改前 **2/30**、§2 改后 **0/30** 且名册逐发闭合、
      > §3 本程自己造的两枚变异 **6/6 红在 `goroutine_test.go:67`** 的超时诊断、§4 修法五查、§5 独立 500 发归因探针、`## 总判`）。
      > 那七枚 commit 全在盘上（`2546b76 4666fe0 e2bb38d 941acf3 869206b 3783cb2 2088e80`，逐枚 `git cat-file -t` 现量 = commit）；
      > 末枚 `2088e80` 是**我代提交**的它的未提交末节——它自己那两枚 `1049513`／`1b93810` 盘上不存在。
      > ⚠ **两枚实现侧的表（`136-ac11-impl.md`／`136-ac11-second-witness-readings.md`）不与本表互相抵账**：
      > 翻勾只凭终裁表那四格；实现侧的自述在终裁表里被独立复算过（2/30 与 0/30 是**它自己重跑的数**，不是抄的）。
      > **它明确没核的档**（照这条读，别把本格当全裁）：本包**其余 64 枚**有没有别的既有偶发失败没量；
      > linux 容器里没复跑（该包纯 Go）；`d667888`／`17b61a4`／`3c3e388`／`d0ca6a3`／`2644207` 五枚号只做到"盘上不存在"。
      > ⚠ 编排者那条"必须具名解冻 `internal/observe/goroutine.go` 才修得了"的复核：**不成立**（与实现方同判，修法只碰 `_test.go`）
      > ⇒ 当时那次解冻请求被我拒绝是对的，这一条现在是**两程独立同判**而不是我一个人的判断。
      > **撤销口令**：回"撤 136 AC#11" ⇒ 我把这一格退回未勾、票面与四份证据一字不动。

- [x] **AC#12**（09-23 23:4x 编排者追加，来源＝`acceptor-ticket136-ac8-ac9-r1` §7 `R-136-7`，严重度＝中）
      **settle 那侧"部分读数失败"是可发生、可编译、且报告里查不到痕迹的**：终裁方实测两发探针——
      "整窗只 1 枚可信"⇒ `samples=1 back=60 pass=true`；"一半读数报错"⇒ `samples=2 pass=true`；
      而 `SettleReport` 出线里 `sample_errors`／`last_sample_error`／`dropped_reads` **三枚字段全部不存在（ABSENT）**，
      同包 `StateReport` 早为同一件事加过字段并写明理由（`:164-168`）。
      ⇒ **本格要的是把那一族的"说过自己没测到"从 `sampling` 侧同形地搬到 settle 侧**，让"这一窗丢了几次读数"能从报告里读出来。
      ⚠ **需要具名解冻生产码 `internal/observe/sampler.go:470-501`**（编排者在本格授权之内；**只此一处**，`SampleState` 那侧、阈值、golden 一律不许动）。
        > **解冻范围更正（09-24 09:3x 编排者，追认实现方报回的那处越界）**：我原先给的 `:470-501` **是我把范围写小了**——
        > 判据① 要求字段带 json 标签出现在报告里，那**必然**要动 `SettleReport` 的声明区（锚点坐标 `:442`）。
        > ⇒ **追认成立，范围补入 `sampler.go:431-443`（`SettleReport` 结构体声明块）**，仍**不含** `SampleState` 那侧（`:247-348`）、不含阈值／golden／`scripts/`／`.github/`。
        > ⚠ 这条记成**我的一课，不是它的一课**：具名解冻要按"**这条判据要落盘必须碰哪几段**"来划，不是按"症状所在的行"来划
        > （同一族旧账：'必须解冻多半是没核依赖方向'——这次是我核漏了"声明区"这一跳）。实现方**没有**自己扩权、也**没有**把越界藏进 hunk 里，
        > 而是 `git diff -U0` 逐行列出哪一枚 hunk 在界内哪一枚在界外再报回 ⇒ **处置正确**。
      **结案判据**：① 报告结构体带上"可信样本数／失败次数／最后一次失败原因"，且有 json 标签；
      ② 上面两发探针各钉一枚用例（"只 1 枚可信"与"一半失败"都必须**说不出** pass 的绝对性，即报告要能自陈丢了几个）；
      ③ 落一发变异（把计数摘掉）⇒ 用例必须转红、红名逐名；④ 还原复绿、三态齐；
      ⑤ ⚠ **若发现要让 `slo-check.ps1` 或任何 golden 跟着改才过 ⇒ 停手报回，不许自己动阈值面**。

- [x] **AC#13**（同上来源，`R-136-9`，严重度＝低）`sampler_test.go:31` 那一支是**第二枚"会吞读数的仪器"**（同族于 AC#8，但形状不是 `[0]` 直取）：
      终裁方判"今天无入口踩得到"为真（两处裸构造分别被 `sampler.go:253-255` 与"不读树"挡住），
      **但它一发探针证明：新写一枚正常采样用例即 `index out of range [-1]`** ⇒ 那是"下一个加测试的人"会踩的雷，
      而它一挂同样会把同包其余读数吞掉。**结案判据**：① 复现那一发（自建一枚最小正常采样用例 ⇒ 拍到 panic 并逐名列出被吞的读数）；
      ② 修法只许动那枚仪器（补长度/空集守卫，**要红不许静默**），⚠ **不许改 `sampler.go` 的语义去绕它**；
      ③ 复跑同包整包 ⇒ panic=0、名册与基线逐名相同。**与 AC#11 分开做**：那一格是偶发红，这一格是"一挂吞一片"。

> **AC#12／AC#13 终裁回执（09-24 10:0x 编排者，两格已翻勾 `[ ]`→`[x]`）**
> 验收方＝`acceptor-ticket136-ac12-ac13-r1`（**非实现者**），锚点 `7d73b5f`，证据
> `docs/evidence/s1/136-ac12-ac13-r1-acceptance.md`（301 行；§0 锚点与快照账／§1 AC#13／§2 AC#12／§3 两格之外的账／§4 总判），
> 分三枚渐进 commit（`ae7d593`→`d3e3a43`→`db9abb0`，逐枚 `git show --name-only` 我核过**只含它自己那一枚文件**）。
> **两格都判「成立（PASS、无附条件）」**，读数全部出自它在仓外快照上自己跑的那几发（〔独立复现〕）；它没复算的档自己列在 §4 末。
> - **AC#13**：它在**未修的旧版本**上自建一条正常采样用例，真拍到 `index out of range [-1]`，
>   量到一次吞掉 **17 条**别人用例的读数（实现方自述 16 ⇒ 见下面「计数口径」那条）；
>   修好后同一条用例**只让自己红、其余照常跑完、一条没吞**；守卫拆掉／改成"默默给个零值"各打一发都按预期变红；**没有 `t.Skip` 任何用例**。
> - **AC#12**：它用自己写的 seam 探针量到两形都满足 **"样本数 ＋ 丢弃数 ＝ 实际取的次数"**；
>   并确认"可信样本数"**等价于 `samples` 的长度**（唯一能塞进脏数据的路径它一打就红＝AC#9 那枚钉还在）；
>   生产码改动**在界内**（界＝上面那枚 `>` 更正块追认的 `:431-443` ∪ `:470-501` 两段）。
> - **它裁掉了我留在票面上的一个歧义**：判据① 我列了三样（可信样本数／失败次数／最后一次失败原因），
>   实现方只加两枚 key、第三枚 `dropped_reads` 没加。验收方判**不需要**（rides 在 `samples` 上、且被上面那发等式量过）。
>   ⇒ **本编排者接受，并把它记成"我判据①写多了"**：**若日后有人按字面要三枚独立 key，那是"改判据"而不是"补料"**，
>   必须具名改本票面再退回，**不许由验收方或实现方任一侧自行解释**（它 §4 第 2 条原话，我照抄在此，免得那句只活在证据里）。

> **⚠ 计数口径（今后引用这一格必须带上）**：AC#13 那句"被吞掉的读数枚数"**实现方量到 16、终裁方量到 17**，
> 差一枚的原因是**探针插入位置差一位**。⇒ 引这个数**必须连"插在哪个位置"一起引**，否则又是一枚会腐坏的数字
> （同一族旧账：引"某区间几枚红"要连口径与锚点 sha 一起引）。**两本账都不作废**，各自标来源：
> `136-ac12-ac13-impl.md` 引 16、`136-ac12-ac13-r1-acceptance.md` §1.1 引 17。

- [ ] **AC#14**（09-24 10:0x 编排者追加，来源＝`acceptor-ticket136-ac12-ac13-r1` §4 第 1 条，严重度＝中）
      **`R-136-7` 的修法方向里有一味既没进 AC#12 的判据、也没人做；本票不许让它作为"半句话"沉在两本账之间。**
      原文（`136-ac8-ac9-r1-acceptance.md:463`「修法方向」栏，逐字）：
      「给 `SettleReport` 补 `sample_errors`/`last_sample_error` 两枚字段，**并产出一枚与 `sampling` 同形的自陈门行**
      （要动 `sampler.go:470-501`，**生产码**）；配一发"部分未测"的钉子」。
      ⇒ 前半句＝AC#12 已结；**"自陈门行"这一维今天不存在**：
      采样侧 `StateReport` 有 `Verdicts []Verdict json:"verdicts"` ＋ `Pass bool // all Gate verdicts pass`
      （`internal/observe/sampler.go:189`，构造在 `:331`/`:335`），而 `SettleReport` **没有 `Verdicts` 字段**
      （实现方 `136-ac12-ac13-impl.md` §5-6 如实列了"没造门行，加它是更大的出线形状改动"）。
      **这一维要的东西**：让"这一窗丢了几次读数"成为**一行会自己判失败的门**，
      而不只是一个**得知道去看才看得见的数字**（同一族教训：`fail-closed 兜住 ≠ 功能通`、"一步存在却从不产出结论"）。
      **结案判据**：① `SettleReport` 出线带上与 `sampling` 同形的门行（每行点名"检查的是哪一道门、过没过"），
      且"部分未测"这一形**必须由那一行自己说出不 pass**，不许只靠 `pass=false` 一个总布尔代答；
      ② 配一发钉子：一半读数报错 ⇒ 该门行红、红名逐名；把门行摘掉 ⇒ 钉子必须转红（⚠ 先证变异落地再读数）；
      ③ 还原复绿、三态齐；④ 名册差集：不得有绿转 SKIP。
      ⚠ **停手上报线（比 AC#12⑤ 更严）**：本格动的是 **settle 报告的出线形状**——
      若发现要让 `scripts/slo-check.ps1`、任何 golden、或 `SPEC-02 §3`（schema 契约级、不得自行增删字段）跟着改才过
      ⇒ **停手报回，不许自己动**；那已超出"补一枚仪器"，属**契约面变更**，须走 `SPEC-12 §4.1` 人工批准。
      ⚠ **具名解冻本票面暂不给**：派单时由编排者按"**这条判据要落盘必须碰哪几段**"现划，
      **不许照抄 AC#12 那两段的坐标**（那是我已记在册的"我把范围划小的一课"）。
      ⚠ **排程**：排在 AC#10 之后（AC#10 要 `cmd/wisp` 空出来，而本格读它的出线）。

> **AC#14 具名解冻＋三裁定（09-24 18:4x 编排者；来源＝`auditor-ticket136-preflight` §3.3／§3.4／§3.5，全部我在盘上复量；本格原判据一字未改）**
> **先认一枚**：`:250-251` 那句"派单时由编排者按'这条判据要落盘必须碰哪几段'现划"——我上一版派单简报里把 **AC#15 的 `:278` 地界句（"只许动 `internal/observe/**_test.go`"）当成了 AC#14 的**。
> 现量：AC#14 是 `:232-252`，通篇**没有**那句；后果很实在——照它派，实现程会以为不许动生产码，而这一格**非动不可**（下面第二行就是证据），它要么停手报回、要么自己越界。**这条记在我身上，不记在它名下。**
> **① 具名解冻（本票面此前一个字没给，现在给）**：只解冻 `internal/observe/sampler.go` 的**两段**——
> `:431-456`（`SettleReport` 声明块，`Verdicts` 字段今天不在里面）＋ `:462-521`（`(*Sampler).CheckSettle` 体内，`:517-520` 那三布尔的折叠处在里面）。
> **边界（越界即停手报回）**：`StateReport` 那一侧**冻结**——`:189-190`（`Verdicts` ＋ `Pass bool // all Gate verdicts pass`）、`:331`（`buildVerdicts`）、`:341-346`（`Gate && !Pass` 折叠）一字节不许改；
> 要给 settle 造行就**另起一枚构造函数**，不许把 `buildVerdicts` 做成两用的（那是把已结案的 `StateReport` 形状拖进本格的射程）。
> **为什么非动不可**（盘上现量，不是推理）：`sed -n '431,456p' internal/observe/sampler.go | grep -c Verdicts` ⇒ **0**；`:520` 逐字 `rep.Pass = memOK && backInTime && releaseOK`——**三个布尔的与，与门行无关** ⇒ "补一枚会自己判失败的行"在这枚结构体上**无法只用测试实现**。
> **② 门行做成 `Gate: true`，但落地顺序被一枚读数钉死**（这一支是**门禁行为变更**，不是"补仪器"）：
> 现量后果链＝`scripts/slo-check.ps1:381` `$allPass = ($failingStates.Count -eq 0) -and $settlePass` → `:396` `if (-not $allPass) { exit 1 }` ⇒ **`wisp slo -settle` 在"部分未测"的真窗口上从 exit 0 变 exit 1，颜色会打到 CI 的 `slo-smoke`（托管，每推送/PR）与 `slo-full`（本机 self-hosted）**。
> 所以我**不批准"先做成 Gate 再看"**：**实现程第一发必须先量**——编队安静窗口内跑 `wisp slo -settle` **≥5 次**、逐次记 `sample_errors`。
> 全 0 ⇒ 按 `Gate: true` 落地（这一形今天不会误伤存量合法窗口）；**任一发非 0 ⇒ 停手报回**，因为那等于把一台有争用的机器变成"固定红、且红因与 SLO 无关"＝**我自己吃过的那枚信号量之病**（`ci` 连红两天让新真伤失去信号），那一步要人拍板、不由实现程自决。
> **并且**：门行做成 Gate 时，**红句必须把 `sample_errors` 点名印出来**（"丢了 N 次读数"），否则下一位看不出这枚红是"测不满"还是"内存超了"。
> **③ 与 AC#15 串行，不并行**：两格共用 `internal/observe/sampler_settle_coverage_136_test.go`——AC#14 判据② 要的正是**推翻**那里 `:208-210` 那句 `if !rep.Pass { t.Fatalf("disclosure leg, not a verdict leg: this window still passes, ...") }`（它现在逐字要求"丢一半读数那一形**仍然 pass**"，与 AC#14 方向相反）。
> ⇒ **AC#14 先做**，且**派单预先授权它改那一处断言**、交件里必须点名"断言被动了、理由是票面 `:243-244`"——不然下一位验收方会拿"实现方改了既有断言"这一条把它判成放水（这正是我记在册的"判放水只看两条"里的那一条）。
> AC#15（那里 `:189` 的 `precondition broken: only %d reads taken` 偶发红）**排在 AC#14 之后**，且 `AC#14` 落地后它的靶子形状可能已经变了 ⇒ 到时要现量重划，不许照抄今天的行号。
> **④ 排程反过来：AC#14 先、AC#10 后**。`:252` 那半句"**而本格读它的出线**"按盘上量**不成立**：`grep -rn "Verdicts" cmd/wisp/*.go` ⇒ 只命中 `providers_test.go` 的 `probeVerdicts`（`wisp providers probe` 那一族，与 SLO 无关）⇒ `cmd/wisp` 今天**一枚 `Verdicts` 都不读**。
> 只有当 AC#10 顺手把"根本没测"标进 `slo_windows.go:623` 时那枚耦合才存在，而按上面的改写 AC#10 收敛成"只装仪器、不动生产码" ⇒ **耦合消失，两格不同目录、可并行**；
> 但 AC#10 要占机（真取样）、AC#14 不要，**机时是稀缺资源** ⇒ 先派 AC#14（只缺一张纸，那张纸上面已经给完了），AC#10 等编队安静窗口。
> 撤销口令：**「撤 136 AC#14 解冻」** ⇒ 解冻与②③④ 一并作废（本格退回"不可派＝缺一张授权"那个状态，票面原文不动）。


> **AC#14 交回后的编排者裁定（09-24 19:2x；实现件 `aef82f5`＋报告 `docs/evidence/s1/136-ac14-impl.md` 318 行；下面每一条我都在 `98665ef` 上复量过，没有一条按自述收下）**
> **交付态＝门行按"记录行"落地、`Gate` 那一支没做**：`sampler.go:556` `const settleCoverageRowGates = false`，`:537` 折进门行（`foldSettlePass(memOK && backInTime && releaseOK, rep.Verdicts)`）。
> **它没抗命——是我那张纸划错了范围。** 上面第③条我写"预先授权它改**那一处**断言"（`:208-210`）：**现量三处同形**（`git show ca2c55e:...coverage_136_test.go` 的 `:143-145`／`:208-210`／`:252-254` 逐字都是 `if !rep.Pass { t.Fatalf(...) }`，
> 其中 `:144` 是 "this leg pins disclosure, not the verdict; a covered-enough window must still pass"、`:253` 那枚 `t.Fatalf` 只印 `report=%+v`），
> ⇒ **任何抓得住"丢一半读数"的门，必然同时把那三形一起抓红**（覆盖率严格更差，没有只红一处的中间形状）；放过另两形只剩"非单调阈值"＝**新造阈值**，那是 `AGENTS §1.1` 的禁面。
> 它按上面第②条的分支落在记录行，并留了**一行开关**：只翻 `settleCoverageRowGates` ⇒ 实测**恰好红 2 枚**（它的变异 S1）。**这一支我暂时不批**，理由不是代价不明，而是"翻完之后**必须仍绿的有哪几枚**"这一半没人量过——`HEAD` 上还剩 `:143`、`:280` 两枚要求 pass，另有 `:307` 一枚是"**健康**窗口必须 pass"（翻 Gate 后**它必须还是绿的**，否则就是折叠过头）。
> ⇒ 判据补一条：**终裁程必须独立翻一次那枚布尔**，并给出"红名册**恰好**＝哪几枚 ＋ 哪几枚必须不红"；答不出"必须不红的那几枚"就退回。
> **两处读数偏差登记（都不改变结论，但都会被下一位当事实引用）**：
> (α) 它自述 `git diff -U0` **只有两枚 hunk**（`@@ -452,7`/`@@ -514,9`）⇒ 现量 **四枚**：`@@ -455 +455,13`、`@@ -516,0 +529`、`@@ -520 +533,5`、`@@ -522,0 +540,58`。
> 逐枚判过归属：**全部落在解冻面内或属纯插入**——`:455` 在 `:431-456` 声明块内，`:516`/`:520`/`:522` 在 `:462-521` 体内或其尾插，
> 被删的只有两行（`Pass bool json:"pass"` 那行的注释、`rep.Pass = memOK && ...` 那行本体），`buildVerdicts` 与 `thresholds.go` **零行改动**（`git diff ca2c55e..HEAD -- internal/observe/thresholds.go` 空；新码里两处 `buildVerdicts` 只出现在注释里）⇒ **没越界，但"只动两处"这句不成立**。
> (β) 我 `A180⑪` 与派单里那份"争用名单"是 **8 枚**，真值 **15 枚**（`scripts/slo-check.ps1:155-157`：`go`/`gofmt`/`cgo`/`compile`/`asm`/`link` ＋ `gcc`/`g++`/`cc1`/`cc1plus`/`as`/`ld` ＋ `wisp`/`wisp-cli`/`staticcheck`）——
> 那个 8 是我从预检程的表里抄的、它自己漏了工具链那六枚。⇒ **闸门比我写的更严**，方向是好消息，账要改：**这已经是今天第二次"抄来的名单比盘上短"**（上一次是 `A181③` 的 golden 名册 58 枚）。
> **另两处它替我查出来的错，逐条认**：① `lint` 在 `ca2c55e` 就是 `failure`（不是我某轮说的"新红"）；② 宿主直跑 `go test ./cmd/wisp/` 以 `0xc0000135` 收场＝**缺 sherpa DLL**，与本仓 09-23 记的那枚加载期形状同一病，不是新缺陷。
> **AC#15 那格里我引的两句反向证据已被普查程推翻并已换理由**（账 `A181③`，此处只留指针不重述）：仓里**确实有**被跟踪的 SLO 报告 JSON（`docs/evidence/s1/66/66-full-subset-slo-report.json`、`66-settle-1.json`）与 `docs/SLO.md:541` 一处逐字引用，
> **真正的尺是"零枚仪器读它们"**；golden 名册 58 枚＝52 `.sse` ＋ 6 枚 `.go` 比较器。**⇒ 今后判"改这个形状会不会弄坏存量"，一律交"读者名册"，不交"文件落在哪个目录"。**
> **本格状态：`[ ]` 未勾**（实现件自述也未勾，与盘上一致）；`next=`＝派**非实现者终裁**（含上面那条"独立翻布尔＋给出必须不红的名册"的新判据），终裁回来后我再就 `settleCoverageRowGates` 出一次裁定。

> **AC#14 的 `Gate` 那一支：编排者批准（09-24 19:5x；终裁表＝`docs/evidence/s1/136-ac14-r1-acceptance.md`，429 行／7 枚单路径 commit `ca09955`→`30e19ef`，档位①—⑤全〔独立复现〕）**
> **批准内容＝动 1 枚布尔（`sampler.go:556` `settleCoverageRowGates` false→true）＋ 按 `:208` 已建立的形状改写另 2 枚断言（现盘 `:143-145`、`:280-282`）。实现件无需返工。**
> **为什么现在敢批**：上一轮我扣着不批的理由是"翻完之后必须仍绿的有哪几枚"没人量过（`A182⑤`）。终裁程在**仓外快照**里独立翻了那枚布尔（先证落地＋`go build` rc=0），交回**两张名册**：
> **(A) 红名册恰好 2 枚**＝`…SingleTrustworthyRead…`（`:144`）与 `…ZeroFootprintDrops…`（`:281`）——**不是"红了一片"**；
> **(B) 它点名 9 枚必须仍绿**（含 `:307` 那枚"健康窗口必须 pass"与 `:208` 那枚已改写腿），其余 60 枚守恒，实测全绿，`SKIP 0`、名册 71 守恒。
> 承重判定它还做了一发我没想到的：**在 `gate=true` 的快照上重做 `M3a`（撤折叠调用点）⇒ 红 3 枚** ⇒ 折叠那一行在**终态**里承重（我 `A182⑧` 预挡的那条误判没有发生）。另自加一发 **`M4` 摘掉红句里的字段点名 ⇒ 红 2 枚** ⇒ "红句必须印 `sample_errors`"这一味也承重。
> **两处我自己要更正**：① 我 19:2x 那段里**没有 `⑤` 这个编号**（只有 ①—④），"两张名册"那条判据实际落在票面 `:300`——我给终裁程派单时指了一枚不存在的编号，它自己找回了正确的落点；② 我说"三处同形"——**原始枚数是 4**（`ca2c55e` 上 `:143/:208/:252/:279`），第 4 枚正是健康窗；**"三处该被门抓 ＋ 一处必须放过"才是完整形状**，我上一格只数了前三。
> 它还补了两条比我更硬的理由：**比例判据必然非单调**（想只放过某些丢读形，只剩"丢几个算多"这种新造阈值）；**按丢读原因分不开 `:143` 与 `:208`**（两处的 `LastSampleError` 逐字相同）⇒ 票面那句"只改一处"在盘上**不可满足**，`A182③` 确认成立。
> **⚠ 一枚残留风险，我明写不藏**：终裁程**没有复现 §1.2 那六发 `sample_errors` 全 0**（它如实标了〔仅自述〕）。我 19:5x 自己用实现程留下的**同一枚二进制**（`D:\tmp\wisp136ac14\bin\wisp.exe`，mtime 18:53）跑了六发：**逐发 `sample_errors=0`、`samples=40`、`pass=true`、`back_within_cap_ms` 267-288**（读件 `D:\tmp\wisp136gate\s1..s6.json`，取数前后各扫一次争用名单为空）。
> 但**这枚 exe 是从 `ca2c55e` 编的、比 AC#14 的码（19:08）早**，它的 JSON 里**根本没有 `verdicts` 这个键**（我核过：实现程那六发的原始件同样无此键）⇒ **我这六发复现的是"真窗口会不会丢采样"这个物理问题**（该问题只依赖采样循环，而 `git diff ca2c55e..HEAD` 的四枚 hunk 全在 `:455`/`:516-522` 及其尾插，**没碰 `:287-294` 那段读取循环**），**不是"带闸的终态二进制会不会变红"**。
> ⇒ **落地程必须先重建再量一遍**：**翻完布尔、改完断言之后，用含门行的新 exe 跑 ≥6 发 `wisp slo -settle`**，逐发记 `sample_errors` 与 **exit code**；**任一发非 0 ⇒ 停手把开关退回 false 并交回读数**（那一支会把 `slo-smoke`/`slo-full` 打成固定红＝我这几天亲历的信号量之病）。
> **判据（落地程必须逐条答）**：①翻布尔前 `grep -n` 出被改那行＋`go build` rc=0；②改后 `go test -count=2 -v ./internal/observe/` 四数＋名册差集**必须仍 71 守恒、0 红 0 跳**（那两枚断言是**换成会自己判失败**、不是删掉）；③`M1` 摘门行 ⇒ 钉子转红；④新 exe 的 ≥6 发真窗口读数（含 exit code）；⑤门禁全套（`gofmt`/`gofumpt` 现量 **v0.12.0**、`go vet`、`d22scan` 纯净快照＋台账不降）。
> 撤销口令：**「撤 136 Gate 批准」** ⇒ 开关退回 `false`（其余形状不动），本 `>` 块保留作"批过又撤"的记录。



> **AC#10 追加两问（09-24 10:0x 编排者，来源＝`acceptor-ticket136-ac12-ac13-r1` §4 第 3、4 条；原判据一字未改）**
> 终裁方量到两件事**不属 AC#12 的债、但 AC#10 结案时必须逐名回答**：
> ① **两枚新 key 今天零生产读者**：`SampleErrors`／`LastSampleError` 在 `internal/observe/` 之外**无任何引用**，
>    `cmd/wisp` 只序列化不读，`slo-check.ps1` 只取 `.pass` 与 `.settle.free_os_memory_count`。
>    ⇒ AC#10 那一发端到端要答的是：**真报告里 `sample_errors` 到底非零过没有**（不是"字段存在没有"）。
>    ⚠ 若答案是"永远零"⇒ 按本仓教训记为"一步存在却从不产出结论"，**不许**用"AC#12 判据没要求读者"糊过去。
> ② **`cmd/wisp/slo_windows.go:623` 那枚合成失败报告**里 `sample_errors=0` 的语义是**"根本没测"**、不是"零丢失"
>    —— 与"测满零丢"在出线上**同形**，目前只靠 `pass=false` 挡住静默绿。
>    ⇒ 落点在 `cmd/wisp/**`（AC#10 的地界、也与**票 133** 的腿清单相邻）⇒ AC#10 结案判据加一条：
>    **这两形必须能被区分开**，区分手段写进证据，不许写成"靠 pass 字段就够"。

- [ ] **AC#15**（09-24 14:5x 编排者追加，来源＝`worker-ticket136-ac11-r2` 交件里"顺带发现"那条，严重度＝低但同族）
      **本票 AC#12 今天新装的那枚用例自己会偶发红**：`TestCheckSettleHalfTheReadsFailedReportsItsLoss`
      （`internal/observe/sampler_settle_coverage_136_test.go:189` 那一族，红句形态＝`only 2 reads taken` 这类"窗口内没取满读数"）
      在**盘上可查的 240 发整包 `-v` 里红 1 发**（14:3x 第二 witness 报回并留在它的表里；我核过那枚文件与那一段行号真实存在）。
      **为什么本格不是"再一个偶发失败"就算了**：它与 AC#8／AC#11／AC#13 是**同一族**——
      一条用例的偶发红会把"这一包到底测没测"变成掷硬币，而门禁只看得见"这一次红没红"；
      且这一枚**是本轮加固动作自己带进来的**（AC#12 为补"部分读数失败"而装），**谁装的谁负责留可复算的收口**，不许留给下一位当"未知红"。
      **结案判据（照 AC#11 的四条，一条不许省）**：
      ① 先量复现率并给 **n**，且**分口径写**（整包单发 `-count=1` vs 同进程 `-count=N` 连发是两个分母，**禁止加总成一个百分比**——AC#11 已为此钉过一条）；
      ② 根因归到**取读数的窗口/间隔与断言之间缺的那条判据**，⛔ **不许**收在"机器负载高"这类不可核的说法上；
      ③ 修法只许把"等待"变成**有判据的等待**（轮询到条件成立＋单调超时上界）；
      ⛔ **绝对禁止**：`time.Sleep` 糊窗口／`Skip` 这条用例／放宽或调高任何阈值／删断言（本仓口径：跳过＝把"没测"洗成"通过"）；
      ④ 改完同一命令复跑 **≥30 发命中必须 0**，两批逐名读数都留在证据里，并给**名册双向差集为空**（证明 0 命中不是"用例变少／被跳过／被 panic 吞掉"）。
      ⚠ **地界**：只许动 `internal/observe/**_test.go`；`sampler.go` 的 `Settle`/`checkSettle` 语义若必须变 ⇒ **停手报回**（那与 AC#14 的出线形状、`SPEC-02 §3` 的 schema 契约同侧，属人工批准面）。
      ⚠ **排程**：**不排在今日队里**，挂在 AC#11 结案之后、与 AC#10／AC#14 同批裁；本格是"登记"不是"承诺日期"。
      **撤销口令**：回"撤 136 AC#15" ⇒ 我把它退回成第二 witness 表里的一条报名（票面这一格删掉，报名与出处不动）。

> **AC#15 这一格里有三处要更正（09-24 18:4x 编排者；原句一字不抹，出处＝`auditor-ticket136-preflight` §3.2／§3.4，我复量一致）**
> **(a) 符号名是我的错，盘上不存在那两个名字**：`grep -n 'Settle' internal/observe/sampler.go cmd/wisp/slo_windows.go` 现量 ⇒ 生产码只有 **`(*Sampler).CheckSettle`（`sampler.go:462`）**，cmd 侧另有 **`runSettle`（`slo_windows.go:599`）**；`Settle`／`checkSettle` 两个符号**一枚都没有**（`Settle` 只是 `slo_windows.go:118` 那个字段名）。今后本族引它一律写实名。
> **(b) "`SPEC-02 §3` 同侧"这句引用不成立**：`docs/specs/SPEC-02-data-storage.md:22` 的 §3 标题确实是"Schema（契约级，不得自行增删字段）"，可它 `:23-` 起讲的是 **`CREATE TABLE schema_meta / profile / …` 的 SQLite DDL**，**不管 SLO 报告的出线形状** ⇒ SLO 报告的形状落不进它的射程。（**是不是另有契约在管**——`PLAN.md` D32/D39 那一族、`SPEC-10`、`SPEC-11`——预检程只排除了 `SPEC-02 §3` 这一处，**没做全仓契约面普查**，登记为未决，AC#14 落地前补那一查。）
> **(c) 真正该停手的那一支，这一格反倒没写**：另两处触发线今天**不触发**——`slo-check.ps1` 对 settle 只读 `.pass`（`:351`）与 `.settle.free_os_memory_count`（`:353`）、**全仓没有严格解析**（`DisallowUnknownFields` 只有 `internal/config/parse.go:72,123` 两处，都不在这条链上）⇒ 加字段不弄坏它；golden 也不动（`git ls-files | grep -i golden` 名册里全是 LLM 的 `.sse`，带 `back_within_cap_ms` 的文件都在 `build/` 而 `build/` 被 `.gitignore:14` 忽略 ⇒ 仓里没有 SLO 报告 golden）。
> ⇒ 该站住的是**门行做成 `Gate: true` 会把 exit 色打到 CI**（账在 AC#14 格下方 `>` 第②条，那里给了"先量 `sample_errors` 再决定"的前置）。本格的停手线**照旧有效**，只是理由换成那一支。
> **(d) 排程**：本格改**串行在 AC#14 之后**（共用文件，账在 AC#14 格 `>` 第③条），不再与它"同批"。


- [2026-09-23 22:3x +08] agent=orchestrator did=**AC#1 翻勾**（`[ ]`→`[x]`）＋ 落笔 AC#9/AC#10 ＋ **更正我自己上一条登记里的一处错**。
  **AC#1 终裁依据**＝验收方 `acceptor-ticket136-r1`（锚 `e8190bf`，四枚 commit `5b28331`/`fa65384`/`5125899`/`ae6a01e`，
  372 行证据 §0-§9 齐）判 **PASS、无附条件**：十发变异各响各的腿（M1 删守卫⇒红点 `:75`；M2⇒前提腿 `:61`；M3⇒正向对照腿 `:151`；
  M4 翻 `Pass`⇒红点 `:75`；M6 降级成非门／M7 **改名**⇒红在 `:84`；M8 条件反写；M9 空 Samples 独立路；M10 settle 侧），
  每发先证落地再读数、还原后 56/56 复绿也证；**不是自己算自己**（断言链逐跳落到 `sampler.go:290-298/332-346`、`thresholds.go:84-206`）；
  **CI 落点真进真能红**（容器原生 `golang:1.27`：pristine `ok (own line)`，M1 ⇒ `FAIL (own line)` ＋ 定点 strict runner `56/55/1/0 rc=1`）。
  它明确没复算的档我照实留着不背书：linux 全仓两态、AC#7 门禁账、`-count=2`/`-race`、端到端（⇒ 已落成 AC#10）、CI run id。
  ⚠ **AC#8 只做完判据①**（"panic 是既有件、不是 AC#1 的债"这一条被独立复算成立：M3 与 M9 两发独立造出同一形状 ⇒
  吞读数的原因是"空 `Samples` ＋ `:308` 无守卫"，被拖走 4 枚逐名、`blame 739bb15f 2026-09-20`），②③未做、**一字未修** ⇒ AC#8 本格不翻。
  **更正（append-only，原文不抹）**：上一条里我写的"实现方引的 `ci.yml:288` 其实是注释行 ⇒ 行号引偏"**不成立**。
  我刚刚自己 `grep -n` 重走：`:288` 逐字是 `run: bash scripts/portable-tests.sh --scope=core`，注释块是 `:267-287` —— **实现方引对了，是我记错**。
  验收方把它按"误记不是注入"处理并登记为 `R-136-5`，判得对；今后引用 CI 落点指 `portable-tests.sh:140/:175` 与 step 名 `ci.yml:266`。
  还有一处**两边都留、谁也没改谁**的读数分歧：全仓既有红包我量到 4 枚、验收方量到 3 枚（差 `internal/agent/approval`，
  它那边该包 `ok 0.720s`；`git diff 09edf02..e8190bf` 未动过那四条路径）⇒ 能说的只有"负载/环境差，各以锚点为准"，**不合并成一句结论**。
  next=派**一名非实现者**做 AC#8②③ ＋ AC#9（同包同文件、按顺序）；AC#10 排在 133 落地之后（要 `cmd/wisp` 空出来）。

- [22:47 +08] agent=worker-ticket136-ac8-ac9 did=**AC#8②③ 交件**（锚点自己量＝`1d38206`；`cmd/wisp/**` 整块未碰）。
  ②：`internal/observe/sampler_test.go` 的 `TestSamplerGoroutineAccountingFollowsRegistry`（`:297`）在 `rep.Samples[0]`
  直取前加长度守卫（commit `79ddd49`，净面只有这一枚文件、10 增 0 删）；空 `Samples` 时 `t.Fatalf` 红在 `:315`，
  **没有 `t.Skip`、没有静默返回**，`sampler.go` 一字未动（`git diff 1d38206..HEAD -- internal/observe/sampler.go` 无输出）。
  同包 `[0]` 直取按 `grep -rn '\[0\]' internal/observe/*_test.go` 全量普查＝**17 处**，逐枚判：**只改 1 处**，
  其余 16 处的理由（同行 `len(x) != N ||` 短路／前置 `t.Fatalf` 长度守卫／`SplitN` 天然非空）写在证据 §2.2；
  另登记一枚**同族但不属 `[0]` 形状**的相邻件：`sampler_test.go:31` 的 `f.mu[len(f.mu)-1]`，实测当前无入口能踩到，未改。
  前后对照（同一发 M3＝`sampler.go:290` `<= 0`⇒`>= 0`，两棵独立仓外纯净树，每发先无条件 restore 再打、还原后 `diff -q` 证逐字相同）：
  改前 tree0＝`git archive 1d38206` ⇒ `rc=1 / RUN=52 / PASS=47 / FAIL=5 / SKIP=0 / panic=1`，拖走 4 枚逐名（含本票 AC#1 两腿）；
  改后 tree1＝`git archive 79ddd49` ⇒ `rc=1 / RUN=56 / PASS=50 / FAIL=6 / SKIP=0 / panic=0`，名册与基线 `diff` 无输出。
  ③：M3 的整包 `-v` 读数里直接取得到 AC#1 那枚钉——腿 A `TestSampleStateZeroSampleWindowFailsClosed` 绿、
  腿 B `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable` 红在 `sampler_zerosample_136_test.go:151`，
  **本程未跑过任何 `-run` 定点绕过**；AC#1 的断言一字未改。三态齐（56/56/0/0 → 上述红 → 还原 56/56/0/0）。
  ⚠ 一条措辞按实测报回：票面②那句"只红这一枚"**不成立也不该成立**（M3 同时开掉 `sampling` 门，另有 5 枚本就该红，
  改前 panic 之前已红 4 枚可自证）；本格真正被钉住的性质是"不再由一枚 panic 代答、每枚各红各的、名册完整"。
  未自勾 AC#8。证据 `docs/evidence/s1/136-ac8-ac9-impl.md` §0-§3。
  next=接着做 AC#9（同序：守卫已在位，整包读数不再被吞）；两格终裁与翻勾派**非实现者**，门禁读数与"未做的档"一并留在证据里。

- [23:02 +08] agent=worker-ticket136-ac8-ac9 did=**AC#9 交件**（①②③，`date -u` 原文 `Wed Sep 23 15:02:12 UTC 2026` ＋8）。
  ①：新增 `internal/observe/sampler_settle_zerosample_136_test.go`（commit `2f291d0`，153 增 0 删，两枚用例与 AC#1 同形）——
  腿 A `TestCheckSettleZeroTrustworthySamplesFailsClosed`（前提腿 `:63/:66` ⇒ 归因隔离 `:73/:76/:79`（released 计数与内存比较三项**全满足**，
  所以 `pass=false` 只能来自"从未记录可信读数"）⇒ 钉 `:85/:88` ⇒ 出线自陈 `:101-113`（`samples` 空＋`pass:false`＋`back_within_cap_ms:-1`））；
  腿 B `TestCheckSettleTrustworthyReadsAreRecorded` 正向对照防恒真（与已有 `sampler_test.go:243`/`:276` 不重复：那两枚都没断言过"可信读数会落进 `Samples`"）。
  ②：`sampler.go:477` 的 `> 0` ⇒ `>= 0`（M10）**先证落地再读数**——`477: if err == nil && m.PrivateWorkingSetBytes >= 0 {` ＋ `diff -u` 单行 ＋ `go build ./...` rc=0，
  然后整包 `-v`：`rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0 / panic=0`，红名**只一枚**＝`TestCheckSettleZeroTrustworthySamplesFailsClosed`，红点 `:66`，
  消息里就是那串病形（`BackWithinCapMS:20`、5 枚 `TreePrivateBytes:0` 的"样本"、`FinalBytes:0`、`Pass:true`）；腿 B 同发仍绿。
  **改前对照我也复算了**：没有本程这枚钉的树（tree1＝`79ddd49`）落同一发 M10 ⇒ `56/56、rc=0`（复跑两次同数）⇒ 验收方 `R-136-1` 那句在我的锚点系上成立。
  ③：还原 ⇒ `diff -q` 逐字相同（并另证与仓库里的生产码逐字相同）⇒ 复绿 `58/58`。**三态齐**。
  ④约束：`SampleState` 那侧语义一字未动（`git diff 1d38206..HEAD -- internal/observe/` 只有我那两枚文件）、阈值/golden/`thresholds.go` 未动、
  已有的 `TestCheckSettle*` 两枚未改未 Skip（整包读数里它们同发 `--- PASS`）。**本程生产码零改动**（`go build ./...` rc=0，`sampler.go` 与锚点逐字相同）。
  自加一发 M11（`:486` 不记 `Samples`）证腿 B 不哑：只红腿 B（`:142`）、腿 A 仍绿。
  门禁：`gofmt -l internal/observe/` 无输出；`gofumpt -l`（本机 **v0.12.0 / go1.27.1**，CI 是 `@latest` 未钉版本 ⇒ 该读数只在改版前有效）无输出；
  `go vet` 双 GOOS 逐错误行归因（原生 `./internal/observe/` 与 `./...` 都 rc=0；`GOOS=linux ./...` rc=1 的**唯一**一条落在
  `cmd/wisp → sherpa-onnx-go-linux@v1.13.8 build constraints exclude all Go files`，剔掉 `cmd/wisp` 后只剩 `cmd/balldebug` 那一条形同，
  再剔掉两枚 `cmd/` 后**30/30 非 cmd 包 linux 交叉零输出 rc=0**）；`-count=2 -v` 四数 **112/112/0/0 → 116/116/0/0**，逐名差集只有本程新增两枚 ×2 轮、
  无改名无消失无转 SKIP；`sh scripts/d22scan.sh` 两形 rc=0，`ban #8 internal/` 401→**402**（本程那枚 `_test.go`）、各 scope 无一下降，
  并核清 `frontend/` 的 40→43 是工作树里未跟踪的 `frontend/dist/` 三枚构建产物、**不是**兄弟在飞的树（`git diff --name-only 1d38206..HEAD -- frontend/` 无输出）。
  ⚠ 三条报回（不硬改读数、不放宽断言）：①**本格未动生产码**，"报告要说出自己没测到"钉在**现有哨兵**（空 `samples`＋`-1`＋`pass=false`）上——
  若终裁要求像 `SampleState` 那样给 `SettleReport` 产出一枚显式 gate 门行，那要改 `sampler.go:477-501`，本程按派单停下没动；
  ②量到一枚**既有 flake** `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:33`，`PerTask mid-task = 2, want 3`，观测 1/9，非本程造的），
  它和"吞读数"同族，建议单立一格；③派单提的"`GOOS=linux` 交叉会停在 cgo／`pathresolver` 那形"这一发**实测未成立**（3f 含 `internal/risk` 零输出），只登记差异。
  未自勾 AC#9。证据 `docs/evidence/s1/136-ac8-ac9-impl.md` §4-§7（§6 是 12 条未做的档，含 linux 全仓两态、CI run id、端到端）。
  next=两格都交完 ⇒ 派**非实现者**按票面 AC#8②③／AC#9①②③ 终裁并翻勾；同时请裁 §6.5（要不要补生产码那半件事）与 §4.5（flake 是否另立一格）。AC#10 仍排 133 之后。

- [23:06 +08] agent=worker-ticket136-ac8-ac9 did=**上一条的一处读数更正（append-only，原文不抹）**：我在 AC#9 那条里写既有 flake `TestNoopTaskReturnsToBaseline`"观测 1/9"，那是当时的事实；交完两格后我对**最终工作树**（`f08c247`）又取了一发未变异的整包 `-v`，同一枚仪器**再次命中同一红点**：
  `rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0`，红名 `TestNoopTaskReturnsToBaseline`、`goroutine_test.go:33: PerTask mid-task = 2, want 3`（原文 `/d/tmp/wisp136ac8-ac9-final-v.txt`；随后连跑 3 次 58/58 复绿，`/d/tmp/wisp136ac8-ac9-final-r1..3.txt`）。
  累计：本程约 27 发整包读数命中 2 发（1 发在未变异工作树、1 发在 tree1＋M10）。AC#8／AC#9 的任何判据不因此改动——两格的三态都是在**红名为本程那枚仪器**的读数上判的，且这两发里 `SKIP=0`、`panic=0`。证据 §4.5 与 §5-G4 已同步补记。
  next=同上（终裁派非实现者）；这条只是把"flake 只出现过一次"这句从我的账里撤掉，别让下一个人按 1/9 去估风险。

- [2026-09-24 00:0x +08] agent=orchestrator did=**AC#8／AC#9 两格翻勾**（依终裁表，两格各判 **PASS、无附条件**）＋ 落笔 **AC#12**（`R-136-7`，settle 侧"丢了几次读数"报告里查不到）与 **AC#13**（`R-136-9`，第二枚会吞读数的仪器）。
  **终裁依据**＝`acceptor-ticket136-ac8-ac9-r1`（536 行、§0-§10、十一枚渐进 commit `9bed3e4`…`2df8ac2`；`internal/observe/**` 一字未写、未装任何工具、未翻任何勾）。
  三处值得记的复算：①它**用的是我 23:3x 更正后的性质**，并且**独立复核了那条更正的前提**（6 枚红名逐枚查因果，判"无一枚误伤、无一枚靠 panic 顺序侥幸绿过"）；
  ②AC#9 那枚钉的"够不够"它**判得出来**：未改码上 `Samples` 空 ⇒ `-1` 摘不掉 ⇒ `pass=false` 结构不可破，
    而**四种可编译的改法**（V1 初值 `-1`⇒`0`／V3 摘足迹判据／V4 把 `samples` 收成 `omitempty`／V6 `backInTime` 改拿 `FinalBytes` 当证据）
    **逐发都被钉认出来、红点各不相同**（`:85`/`:66`/`:102`/`:85`）；V2 编不过 ⇒ 该发作废、它没算进成绩 ⇒ **本格不要求作者动生产码**；
  ③它**推翻了我一句派单前提的适用范围**并纠了实现方一处腐坏引用（见下面"更正两条"）。
  **本格不扩权**：AC#12 需要解冻生产码 `sampler.go:470-501`，我在格内具名授权**只此一处**（阈值／golden／`SampleState` 那侧一律不许动），
  并写明"若必须动 `slo-check.ps1` 或任何 golden 才过 ⇒ 停手报回"。
  ⚠ **更正两条**：(a) AC#12/AC#13 上面标的"23:4x"不准——两格实际落笔跨越零点，真实区间是 **23:5x–00:0x**（我取表时 `date -u` 已是 16:00 UTC）；
  (b) 实现方 §8 那句"全仓 `grep` 'No tools needed' 零命中 ⇒ 那句话不在仓库任何文件里"**已腐坏**：终裁方与我各量到 **4 枚文件命中**
  （含我自己写的 `injection-timeline.md` §11 与本台账 A138/A139）⇒ **是我们把原文逐字登记进去的**，那条判据从此不可用。
  ⇒ 替代判据（已实测可行）：**问"这句话是不是仓内既有文字"要用 `git log -S'<原文>' --reverse` 看首枚引入提交**，
  本次首现＝`172c7aa`（代理自己的登记提交）⇒ 仍证明它不来自仓内既有内容。**`grep -l` 计数会随我们每次逐字登记单调变坏，不能当"来源清白"的证据。**
  next=派 `worker-ticket136-ac12-ac13`（两格同包、AC#13 先做，它不依赖生产码）；AC#10 仍等 `cmd/wisp` 的 133 终裁表回。

- [2026-09-24 08:49 +08] agent=worker-ticket136-ac12-ac13 did=**AC#13 交件**（派单自量锚点 `45c8d1c`；`git cat-file -t` ＝ commit；本程引用到的 9 枚外来 sha 逐枚现核，全 commit，无一枚抄来当快照锚点）。
  ①复现：纯净快照 `tree0`＝`git archive 45c8d1c` 基线 `rc=0/RUN=58/PASS=58/FAIL=0/SKIP=0/panic=0`；我自己写一枚最小正常采样用例（裸 `&fakeTree{}` ＋ `SampleState`，不写任何断言）⇒ `rc=1 / RUN=43 / PASS=42 / FAIL=1 / SKIP=0 / panic=1`，panic 源逐字 `sampler_test.go:31` ← `sampler.go:269`（`index out of range [-1]`）；**被吞 16 枚逐名**列在证据 §1.1（含 AC#1 那两腿与 AC#8 刚修的那枚），并如实记一条位置依赖（同一枚探针放包尾就吞 0 枚——危害由"谁排在后面"决定）。
  ②修法：只动那枚仪器（`sampler_test.go` **22 增 0 删**，唯一 hunk 在 `ReadTree` 的"脚本耗尽"那一支：`len(f.mu)==0` ⇒ 返回 `errFakeTreeNoScript`，走 `TreeReader` 的错误通道；**无 `t.Skip`、无静默零值返回**）；`git diff 45c8d1c..f4c7062 -- internal/observe/sampler.go` **无输出** ⇒ 未改采样器语义去绕它。另加三腿钉 `sampler_faketree_guard_136_test.go`（111 增 0 删）。commit `f4c7062`。
  ③复跑：`tree1`＝`git archive f4c7062` 基线 `rc=0 / 61 / 61 / 0 / SKIP0 / panic0`，名册与 tree0 基线逐名差集＝**只多我三枚**；**同一发探针落在修好的树上** ⇒ `rc=1 / RUN=62 / PASS=61 / FAIL=1 / SKIP=0 / panic=0`，红名只一枚＝探针自己（红点 `sampler_test.go:93`，消息里带 seam 那句话），名册缺 **0** 枚（改前缺 16）。钉不哑：MA 摘守卫 ⇒ 我的前提腿红在 `sampler_faketree_guard_136_test.go:46` 且 **panic 复发、再吞 19 枚逐名**（§1.5）；AB 把守卫改成静默零值 ⇒ 两腿红（`:57`/`:103`）、`panic=0`、名册完整；每发先无条件 restore、还原后 22/22 枚 `diff -q` 逐字相同＋复绿 61/61。**未跑任何 `-run` 定点绕整包读数**。
  约束：与 AC#11 分开做——那枚 flake 本程只命中未修（`-count=2` 第一发红在 `goroutine_test.go:29 live count mid-task = 2, want 3`，票面记的是同函数前一跳 `:33`，只登记差异）；`cmd/wisp/**` 一字未碰；快照只在 `D:\tmp\wisp136ac1213-*`（只建不删）；未翻任何勾。门禁原文（gofmt／gofumpt **v0.12.0 (go1.27.1)** 用既有二进制、未 `go install`／vet 双 GOOS 逐错误行归因：原生整树 rc=0 零行、`GOOS=linux` 整树 rc=1 唯一 3 行诊断落在外部模块 `sherpa-onnx-go-linux@v1.13.8` 的 build constraint，剔两枚点名 cmd 后 **31/31 非 cmd 包零输出 rc=0**／`-count=2` 四数与逐名差集／`d22scan` 两形 rc=0、`ban #8 internal/` 402→403、各 scope 无一下降）在证据 §4。
  next=AC#13 终裁派非实现者；本格未自勾、未动票面判据。接着交 AC#12（同包，需具名解冻生产码）。

- [2026-09-24 09:16 +08] agent=worker-ticket136-ac12-ac13 did=**AC#12 交件**（四判据齐，同一锚点系）。
  ①`SettleReport` 带出 `sample_errors`／`last_sample_error`（`json` 标签齐，与 `StateReport:164-168` 同形；"可信样本数"按同族形状 rides 在 `samples` 上——这条解释明写在证据 §2.1 请终裁过目，不当它已成立）。⚠ **一处越出授权范围，报回**：解冻写的是 `sampler.go:470-501`，但判据① 必然要动 `SettleReport` 的声明区（锚点坐标 `:442`）；本程按判据做了最小插入，`git diff -U0` 旧侧三处 hunk 逐行列在 §2.2 表里（`@@ -442` 越界、`@@ -477` 与 `@@ -479,0` 在授权内；numstat `21 2`）；`SampleState` 那侧 `:247-348` **零 hunk**、`thresholds.go`／任何 golden／`scripts/`／`.github/` 在 `45c8d1c..HEAD` 全范围 `git diff --name-only` **无输出**。commit `5c1529a`（生产码）＋`f53ad5c`／`c03aee3`（钉）。
  ②终裁方那两发探针各钉一枚（`TestCheckSettleSingleTrustworthyReadReportsItsLoss`／`TestCheckSettleHalfTheReadsFailedReportsItsLoss`），另补两腿：另一支丢弃路（零足迹）也必被数到、正向对照防恒真。**同一份只读 JSON 的探针跑在两棵树上**给出改前／改后对照：改前 `samples=1 | sample_errors ABSENT | last_sample_error ABSENT | pass=true`、改后同形 `sample_errors=9` 带原因；"一半报错"那形改前 `samples=6 | 三枚 key 全 ABSENT`、改后 `sample_errors=4`。⚠ 如实：`dropped_reads` 那第三枚 key 两形都仍 ABSENT（本程按"两枚同形搬"办，未加），`dropped_reads`/独立计数字段若被判据①字面要求 ⇒ 记 §5-5 缺口。
  ③落变异四发、各响各腿（每发先 restore、先 `diff -u`＋`grep -n`＋`go build ./...` rc=0 证落地再取整包 `-count=1 -v`）：MD 摘 err 支计数 ⇒ `rc=1/65/62/3` 红名 `…SingleTrustworthyRead…:134`＋`…HalfTheReadsFailed…:197`（第三枚红是既有 flake，非本格仪器）；ME 摘零足迹支 ⇒ 只 `…ZeroFootprintDrops…:247` 红；MF 摘原因 ⇒ `:141`/`:206` 两枚红；MC（自加）放宽 `<=0`⇒`<0` ⇒ **隔壁 AC#9 的钉仍红在 `sampler_settle_zerosample_136_test.go:66`**＋本程腿 3 `:244`，证我重排 `:477` 没把那一族的牙弄钝。
  ④三态齐：未变异 `rc=0/65/65/0/SKIP0/panic0` → 上表四发红 → 还原 23/23 枚 `diff -q` 逐字相同（并证与仓库生产码逐字相同）⇒ 复绿 65/65；`-count=2 -v`＝`130/130/0/0`、去重名册与 `-count=1` 逐名相同；本格九发读数 `SKIP` 逐发 0。
  ⑤**没碰 `scripts/slo-check.ps1`、没碰任何 golden／阈值**（判据⑤未触发）；唯一 Go 侧消费者 `cmd/wisp` 只把测试二进制按新 `observe` 重新链接（`go test -c` rc=0）**未执行**（快照树缺 sherpa DLL，且 `cmd/wisp` 归票 133 在飞）。`SettleReport` 消费面核查在 §3：出线 key 是超集（11→13、无改名无删除）、那条链上没有严格解析器、无仪器读旧归档报告。
  未做的档 12 条诚实列在 §5（含容器 linux 面、端到端、`-race`、CI run id、R-136-8、第三枚 key）；一条仪器级异常（本机墙钟前跳 8h29m，三条量法）与一枚**本程自造的错 sha `f4c7620`**（`git cat-file -t` ⇒ 不存在，落盘前删净、现文件与 committed 版各 `grep`＝0）都按规格登记在 §5-12／§7。两个计数：真通知回显 **7**、判为注入 **0**（§7 逐条带出处；`HANDOVER.md` 回显里"136 两枚在飞"那句我盘上核成"一名代理做两格"的措辞误差，不据它等待也不据它改做）。
  next=两格终裁派非实现者；**先裁 §2.2 那处 struct 越界**（接受 ⇒ 请把解冻范围补成含 `:431-443`；撤销 ⇒ 回退 `5c1529a` 的 struct hunk）；AC#12 的端到端仍归 AC#10／`cmd/wisp` 空出来之后。

- [2026-09-24 09:4x +08] agent=orchestrator did=**核收 AC#13／AC#12 两格交件**（该程撞到 150 轮上限，但按提交序列与产物判＝**两格都交完**：`45c8d1c`→`5c1529a`／`f53ad5c`／`c03aee3`→`1f7875d`→`c6bb91c`→`e25d2a7`，
  工作树对它**零未提交残留**，通知里那句 `result` 只是中段碎片、未采信）＋**裁掉它要我裁的那处**：越界**成立、追认**，
  解冻范围补入 `sampler.go:431-443`（见 AC#12 内那格 `>` 更正块；**记成我把范围划小的一课，不是它越权**）。
  **本格两枚勾仍不翻**（AC#12／AC#13 终裁归非实现者）。三件值得点名的好行为：
  ①它**自报删过树内临时副本**（违反"临时件只建不删"）但把原件／驱动器／全部 `-v` 日志／五棵快照树都留着、并写明"可一次 `cp` 复原"、不当"顺手清理"混过去（`e25d2a7`）；
  ②它**自己造出一枚不存在的 sha `f4c7620`** 并用 `git cat-file -t` 抓到、落盘前删净、现文件与 committed 版各 `grep`＝0（`§5-12`）
    ⇒ 今天新立那条"取 sha 先核存在"**连自造错值都能挡住**，这是它第一次被证明有效、且是被代理**对自己**用上的；
  ③**隔壁那枚钉没被弄钝**：它自加一发 MC（放宽 `<=0`⇒`<0`）专门验"我重排 `:477` 之后 AC#9 那枚钉仍红在 `:66`" ⇒ 主动证"我改的没把兄弟那面搞哑"，这一手以前都是等验收方发现的。
  两处如实报回我收下并登记：`dropped_reads` 第三枚 key **没加**（它按"两枚同形搬"办、把缺口记进 §5-5 请终裁判，不自己扩）；flake 命中在 `goroutine_test.go:29`、票面记的是同函数 `:33`（同一支不同跳，只登记差异）。
  另：它独立量到本机墙钟前跳 **8h29m**（三条量法），与票 137 验收方量到的 8h32m **同量级互证** ⇒ 那条"不许拿时间戳相减当耗时"成立。
  next=派 `acceptor-ticket136-ac12-ac13-r1` 终裁两格（含 §5-5 那枚 key 该不该补、"可信样本数 rides 在 `samples`"这个解释成不成立）；
  AC#10（端到端）待 `cmd/wisp` 空出——那里现在正被票 135 的 M-G 那一程改。

- [2026-09-24 10:1x +08] agent=orchestrator did=**AC#12／AC#13 两格翻勾**（`[ ]`→`[x]`）＋ 落笔 **AC#14** ＋ 给 AC#10 追加两问。
  **依据**＝`acceptor-ticket136-ac12-ac13-r1`（非实现者，锚 `7d73b5f`，301 行证据 §0–§4、三枚渐进 commit 逐枚核过只含它自己那枚文件）
  判两格 **成立、PASS、无附条件**；读数全部是它自己在仓外快照上跑出来的（〔独立复现〕），没复算的档它自己列在 §4 末、我不背书。
  **我按它的裁定接受的三件事**：① `dropped_reads` 不补（rides 在 `samples` 上、被"样本数＋丢弃数＝实取次数"那发等式量过）
  ⇒ 记成**我判据①写多了**，并把它那句"要字面三枚 key 属**改判据**、须具名改票面、不许任一侧自行解释"**照抄进票面**（不许只活在证据里）；
  ② 生产码改动判**界内**（界＝`:431-443` ∪ `:470-501`，即上面那枚我自己"把范围划小"的更正块）；
  ③ 16 vs 17 那枚差判为**探针插入位置差一位** ⇒ 两本账都不作废，但**引这个数必须连位置口径一起引**（新增一条 `>` 口径块）。
  **我不接受、另立一格的一件事**：`R-136-7` 修法方向里那句"**产出一枚与 `sampling` 同形的自陈门行**"
  —— 它**从来没进过 AC#12 的判据**、实现方 §5-6 也如实说没做。⇒ 落成 **AC#14**（含"这一维要的东西"是为什么：
  让"这一窗丢了几次"成为**一行会自己判失败的门**，而不是一个**得知道去看才看得见的数字**）。
  ⚠ AC#14 的**具名解冻我一个字没给**，派单时现划，明令**不许照抄 AC#12 那两段的坐标**。
  ⚠ 另给它加了一条比 AC#12⑤ **更严的停手线**：碰 `slo-check.ps1`／golden／`SPEC-02 §3` 任一 ⇒ 上报（那是契约面，不是仪器面）。
  **顺带**：它报的两件"不属本格债"的事（两枚新 key 零生产读者、`cmd/wisp/slo_windows.go:623` 合成报告 `sample_errors=0`
  语义＝"根本没测"而非"零丢失"、与"测满零丢"出线同形）**追加进 AC#10 判据**（原判据一字未改）。
  本票现值：**5 勾／9 未勾**（9 里含刚落笔的 AC#14）。**未结、不改名 `-done`。**
  入库前我对那份证据做了凭据反扫（按**词**筛＝api key/secret/token/password/Bearer/私钥头 全 0 命中；
  本机候选凭据串 2 枚、文档命中 0）⇒ 干净。
  next＝①`cmd/wisp` 现在仍被票 135 M-G 占（jsonl mtime 10:01＝活着），**AC#10／AC#14 都排在它后面**；
  ②AC#11（flake 2/27 那格）可独立起，不占 `cmd/wisp`；③等 `acceptor-ticket137-ac2-ac5-r1` 交回。
- [2026-09-24 15:0x +08] agent=orchestrator did=**AC#11 翻勾**（`[ ]`→`[x]`）＋ **代提交终裁方自己的未提交末节** ＋ **自报一枚我造的现场事故**。
  **AC#11 终裁依据**＝非实现者 `acceptor-ticket136-ac11-r1`（锚 `51e29b0` 由它 `rev-parse` 自量；证据 656 行、§0–§6＋`## 总判`；
  七枚 commit `2546b76 4666fe0 e2bb38d 941acf3 869206b 3783cb2 2088e80` 我在翻勾前逐枚 `git cat-file -t` 现量＝全 commit）判 **成立、无附条件**：
  ① 改前 **2/30**（`git archive 51e29b0` 那棵纯净树连跑 30 发整包 `-v`、逐名表、红句原文入库，没用任何单次读数）；
  ② 归因它**自己重跑 500 发**而不是抄实现方：`Count()` 6/500 < `PerTask` 12/500、**不存在**任何 `perTask>count` 形状、缺的那一跳只落在两枚空体腿上；
  ③ 修法只碰那一枚 `_test.go`，新增行里 `time.Sleep` **0** 枚、`Skip` **0** 枚、`want 3`／`!= 3`／`before+1` 逐字节未动；
  ④ 改后 **0/30** 且**名册逐发闭合**（30 发顶层是同一枚 65 枚集合、`PASS+FAIL` 恒 65、`SKIP`／`panic`／非零 rc 全 0、目标用例逐发都有 `--- PASS` 行）
  ⇒ 那个 0 **不是**把分母做小或把用例跳掉换来的。**附加那一发它自己造出来了**：两枚变异各 3 发、**6/6 红在 `goroutine_test.go:67` 的超时诊断**、
  每发 `(2.00s)` 有界收尾 ⇒ 新装的那道等待不是装饰（"附条件 vs 退回"的分界在这一格是"造出来了"）。
  **它明确没核、我不背书的三档**：本包其余 64 枚有无别的既有偶发失败没量；linux 容器那半边没复跑；
  实现侧报的 5 枚号（`d667888`／`17b61a4`／`3c3e388`／`d0ca6a3`／`2644207`）只做到"盘上不存在"。
  ⚠ **两枚实现侧的表（`136-ac11-impl.md`／`136-ac11-second-witness-readings.md`）不与终裁表互相抵账**：翻勾只凭终裁表那四格。
  ⚠ **我上一条登记（09-24 10:1x）那句"AC#11 可独立起"里的解冻请求，已被这一程独立复核为不成立**
  （`f06a8d0` 只带测试文件；`Snapshot()`／`NewTimeout()` 改前就存在）⇒ 我当时拒绝解冻是对的，且这一条现在是**两程独立同判**、不是我一个人的判断。
  **代提交那一枚 `2088e80`**：它交件时自报的 `1049513`／`1b93810` 两枚号盘上不存在＝**它的 `git commit` 那一步静默失败**，
  末节（§6.2 凭据自扫回执）留在工作树里。我先替它跑完它自己那三条反扫（候选真值枚数 1／文件命中 0／危险形状 `[]`，
  四枚 32 位十六进制串逐枚判为文件 md5 不是凭据），扫空才代提，**只在它的那一枚文件上带 pathspec**。
  **我造的这枚事故（append-only、不改写）**：上一条我打算登记的 python 追加**把本票文件清空成了 0 字节**——
  它 `open(p,'wb')` 之后再 `nl.join(list[str])` 抛 `TypeError`（bytes 拼 str），**截断已经发生、异常在写盘之前那一步之后**；
  更糟的是我当时只看到 `勾= 0 未勾= 0` 就去忙别的了，把"两个 0"当成了"文件还没写完"而不是"文件是空的"。
  ⇒ 处置：`git cat-file blob HEAD:<本文件> > <本文件>` 还原（398 行、5 勾／10 未勾、AC#15 那一格在 `:253`，逐量对上），
  然后**改用 Edit 工具**做翻勾与这次追加。代价：AC#11 那一格的勾我重做了一遍，**票面没有丢任何历史内容**（HEAD 里都在）。
  ⇒ 记成规矩：**往带中文的跟踪文件里追加，一律走 Edit 工具**；python 只许用于"读＋计数"，
  真要用它写就必须"先读到内存→在内存里拼完→确认非空→再开写句柄"，且**不许在异常后假设文件还在**。
  ⚠ **同一时刻另一处不是我造的变化**：`design/**` 16 枚 tracked 文件在工作树里被移走（现在盘上是 `design/old/`（未跟踪，内容在）
  与一枚空的 `design/doubao/`，移动时刻 15:00）⇒ **不是我动的**（我这程只碰过本票文件），我没有还原它、也没提交它；
  按规矩前端/设计原型归 owner 与外部 agent，这一处由他们收。⚠ 这一条的实际风险是"谁下一枚不带 pathspec 的 commit 会把 design/ 从库里删掉"，
  我在台账里同步登记，提醒所有在飞程。
  本票现值：**6 勾／9 未勾**（未勾＝AC#2/#3/#4/#5/#6/#7/#10/#14/#15）。**未结、不改名 `-done`。**
  next＝①**AC#15 不单独起**，与 AC#10／AC#14 同批（那两格要 `cmd/wisp`，AC#14 的解冻范围要现量重划）；
  ②`cmd/wisp` 现在空不缺人但**先派 137 AC#4**（不占 `cmd/wisp`、且已写明两格不许并做）；
  ③AC#2..AC#7 六格是 09-20 那批的旧账，排在 137 AC#4 与 136 AC#10/#14/#15 之后。
- [2026-09-24T10:41Z] agent=orchestrator did=**只读预检程交回后的四格裁定入库**（`auditor-ticket136-preflight`，表 `docs/evidence/s1/136-ac10-ac14-preflight-r1.md`，306 行／锚 `7cf8075`；**它不翻勾、不裁成立**，只答"这两格能不能派、按什么顺序派"）。**它推翻我简报里四处前提，逐条我在盘上复量一致才落**：P1 我把 **AC#15 的 `:278` 地界句**（"只许动 `internal/observe/**_test.go`"）当成 AC#14 的——现量 AC#14 是 `:232-252`、通篇没那句；P2 "会不会开六态窗口"＝**不开任何窗**（`proc.Boot` 在 `boot_windows.go:57-113` 只做 env 自检/`DefaultLayout`/`OpenJobScope`/`AcquireSingleInstance`，`internal/proc` 全包的 `webview/panel/ball/CreateWindow` 只出现在注释与可为 nil 的钩子字段名上），**但会真取样＋另起两枚子 `wisp.exe`＋每发前 2 秒静置**（`slo_windows.go:374`/`:451`/`:201`）⇒ 编队空那条照旧成立、问法要换；P4 `Settle`/`checkSettle` **两个符号盘上不存在**（实名 `CheckSettle` `sampler.go:462`、`runSettle` `slo_windows.go:599`），票面 `:278` 同错、原句不抹另起 `>` 更正；P5 **"等 133"那半句已被超越**（`cmd/wisp/leg_dispatch_gate_133_test.go` 最后改动是 `fa35557` 09:56，标题属**票 135**）⇒ `cmd/wisp` 现在没有写者。另两处它没提我自己核到：`SPEC-02 §3` 那句引用射程不对（它讲 SQLite DDL，不管 SLO 出线）；"全仓没有 `_test.go` 引用 `slo-check.ps1`"实为 **1 枚命中**（`internal/observe/sampler_test.go:15` 是**注释行**，不影响结论，但按"计数要现跑"记在它名下）。**四格裁定**：①**AC#10 判据② 结构上永远产不出它想量的读数**（方向写反：`sampler.go:332-340` 是"造红门行"，拆掉它 = 让它过 ⇒ 退出码只会 1→0；`samples:[]` 与拆守卫互斥（`FM2` 原文 `samples=0 pass=true`）；且 `134-...:643-645` 已逐字记"这条路在当前码里不通"）⇒ 换成"对 `slo_windows.go:386`／`:323-325` 各落一发变异、读 `go test ./cmd/wisp/` 响不响"，今天 `grep -rn "run\.Pass" cmd/wisp/*_test.go` **0 命中**＝不响、装了才响；同时把"零样表面在 CLI 不可达"**登记为结论＋可复算触发条件**（谁给 `cmdSLO` 加了 reader 注入面谁复活它），不留一枚永不响的 AC；②AC#10 判据① 命令**缺 `-state`**（`:223-225` 现量 return 2）⇒ 换成 `134-...:538` 那发原文，且参照值 `mem_median=4476928` **作废**（行自 09-21 15:06 那枚 `build/wisp.exe`，早于 `5c1529a`）；③**AC#14 给出具名解冻**：`internal/observe/sampler.go:431-456` ＋ `:462-521`，边界＝`StateReport` 那侧（`:189-190`/`:331`/`:341-346`）冻结、`buildVerdicts` 不许改成两用；④**门行做成 `Gate: true` 我先不批"直接做"**——它会把 `wisp slo -settle` 的 exit 0→1，经 `slo-check.ps1:381`→`:396` 打到 `slo-smoke`＋`slo-full`（后者是这台 6C12T 本机），所以**实现程第一发必须先量** ≥5 次 `wisp slo -settle` 的 `sample_errors`：全 0 才按 Gate 落地，任一发非 0 **停手报回**（一台有争用的机器被做成"固定红且红因与 SLO 无关"＝我吃过两天的信号量之病，那一步要人拍板），且 Gate 那枚红句必须把 `sample_errors` 点名印出。**排程反过来**：AC#14 先（不占机、纸已给全）、AC#10 后（要编队安静窗口＋取数期间不推送）；`:252` 那半句"本格读它的出线"按 `grep -rn Verdicts cmd/wisp/*.go`（只命中 `providers_test.go` 的 `probeVerdicts`，同名不同物）**作废**，另起 `>` 不悄改。**新增一枚文件级串行**：AC#14 与 AC#15 共用 `sampler_settle_coverage_136_test.go`，而 AC#14 判据② 要推翻的正是那里 `:208-210` 的 `if !rep.Pass { t.Fatalf("disclosure leg, not a verdict leg…") }` ⇒ **派单预先授权它改那一处断言并点名理由**，否则会被下一位判成放水；AC#15 串行在 AC#14 之后。本票勾数不变（**5 勾／10 未勾**，预检不产勾），`next=` 改成：①派 **AC#14 实现程**（解冻已给完、不占机）；②AC#10 等编队安静窗口；③AC#15 排 AC#14 之后；④AC#2..AC#7 六格旧账仍排最后

- [2026-09-24 19:1x +08] agent=worker-ticket136-ac14 did=**AC#14 那一格交件（门行落地，`Gate` 那一步按读数停手报回）**。
  锚点链：取基线时 `ca2c55e` → 编队落了 `2e6d171`（＝票面 `:325(b)` 要的那一查，结论与本格同向：
  `136-slo-output-contract-census-r1.md:80` "没有任何一份文档在 `SPEC-12 §4.1` 意义上管 SLO 报告的出线形状"）→
  本程码 `aef82f5` → 证据 `d6c83de`＋本次。证据 `docs/evidence/s1/136-ac14-impl.md`（§0 锚点与基线四数／
  §1 那六发读数／§2 逐段落点与冻结面自查／§3 门禁／§4 六发变异／§5 跨包逐名对照／§6 残留／§7 纪律）。
  **落了什么**：`SettleReport` 加 `Verdicts []Verdict json:"verdicts"`（票面 `:277` 划的 `:431-456` 内），
  `buildSettleVerdicts`（**另起一枚构造函数**，`buildVerdicts`/`thresholds.go` 一字节未动）产出一枚与 `sampling` 同形的行，
  `measured` 逐字沿用 `%d valid / %d errors`，红句按**字段名**印 `sample_errors=` ＋ 最后一次丢读原因；
  折叠处 `foldSettlePass` 复刻 `:341-346`。探针 6 枚全在新文件 `sampler_settle_gate_136_test.go`（自带 fixture，不挤 AC#15）。
  **预授权那一处按点名要求报**：`sampler_settle_coverage_136_test.go` 原 `:208-210` 那句
  `if !rep.Pass` ＋ `t.Fatalf("disclosure leg, not a verdict leg…")` **被本程改掉**，理由＝判据①／②；
  改后写成**与 `gate` 取值无关**的形状（行必须存在／行必须自己说 not-pass／红句必须印 `sample_errors`／
  失败 gate 行与 `pass=true` 不得并存），所以翻不翻布尔它都成立——**这是刻意的，不是漏**。
  那枚文件的头注释 `:25-27`（"a partially covered window still passes today"）本程**没动**，它对另外两形今天仍是真的，
  改准它属 AC#15 射程，登记不越界。
  **Gate 那一步（票面 `:281-285`）的读数**：编队安静窗口内 `wisp slo -settle` **六发**（≥5），
  `sample_errors` 逐发 **0**（`exit=0 pass=true samples=40 back_within_cap_ms 263-269`），
  每发前后各量一次 `slo-check.ps1:153-155` 那份名单（**盘上 15 枚，不是简报里那 8 枚**）⇒ **12 次全空**、`retries=0`；
  取数前核 `gh run view 35989390166` ⇒ `slo-full`（这台机）已 **completed/success**、`test-windows` 在托管侧。
  **但本程没有做成 `Gate: true`**——盘上推翻一处前提：票面 `:287` 记"AC#14 要推翻的既有断言"＝**一处**，
  现量是**三处同形**（`:143-145` / `:208-210` / `:252-254`，都是"disclosure not verdict"），
  而任何抓得住"丢一半"的判据必然同时抓住另外两形（覆盖率严格更差：kept=1/lost>=9 对照 kept≈lost），
  要放过另两形只剩"只红在恰好一半"那种非单调阈值＝**新造阈值**，禁。⇒ 扩授权不由实现程自决，
  按 `:284` 那一档落**记录行**（`settleCoverageRowGates=false`）交回裁。**代价已量成读数，不靠推理**：
  打在已含 `:208-210` 改写的树上只翻那一枚布尔（变异 S1）⇒ 恰好红 **2 枚**
  （`TestCheckSettleSingleTrustworthyReadReportsItsLoss`、`TestCheckSettleZeroFootprintDropsAreCountedToo`），
  改法＝本程已用的同款三段，生产码只动 `:556` 一枚字面量。**批准口令＝翻 `settleCoverageRowGates` ＋ 授权那 2 枚断言。**
  **门禁**：`-count=2 -v ./internal/observe/` 两遍各 **142/142/0/0**（顶层 **71 枚 ×2**，基线 65×2=130，净增本程 6 枚）；
  两遍名册 `comm -3` 空、基线→改动后无一名消失或转 SKIP；`grep -ci panic` 两遍各 8（全是用例名带 `Panic` 的两枚，真 `^panic:`=0）；
  gofmt/gofumpt（盘上现量 **v0.12.0 (go1.27.1)**）/`go vet` 全清；`sh scripts/d22scan.sh` 在 `git archive` 仓外纯净快照 **rc=0**、
  八 scope `203/22/40/18/16/40/405/39` 对基线 `…/404/39` 唯一变化＝`ban #8 internal/` ＋1＝本程那枚新 `_test.go`（台账 `:4392` 同形制），一枚不降；
  `GOOS=linux CGO_ENABLED=0` build＋vet 各 rc=0（容器读数没取，§6-2）。
  **变异六发**（先 `grep -n` 证落地＋`go build ./...` rc=0 再读数）：M1 摘门行⇒红 4 枚逐名、M2 门行不看丢读⇒红 2 枚、
  M3a **单点只撤折叠调用点⇒71 枚全绿＝那一处在 HEAD 形状下就是装饰，本程明说不遮**、M3b 同处回退＋gate=true⇒立刻红 3 枚（它变承重的瞬间已有钉子）、
  S1 见上、正向对照未改树 71/71。**跨包**：`cmd/wisp` 整包改动前后各 95/46/**8**/0，红名逐名 `comm` 双向空（本程零影响），
  那 8 枚属票 128/winsec 族、`ca2c55e` 上就在；另记宿主直跑 `go test ./cmd/wisp/` 会以 `0xc0000135`（DLL 缺位）收场，
  与简报 `⚠` 那条同源，AC#10 出去前该知道。`lint` 在 `ca2c55e` 那发 run 上就是 failure（先于本程）。
  **本格状态**：判据① 落了一半（行会自己说不 pass，但不否决总布尔）、②③④ 全量到位 ⇒ **本程不翻勾，`AC#14` 仍未结**。
  禁改面全部没碰（`cmd/wisp/**`、`internal/proc/**`、`internal/winsec/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、
  `internal/risk/**`、`tools/d22scan/**`、阈值／golden／`thresholds.go`、`scripts/slo-check.ps1`、`.github/workflows/**`、票 128/135/140/119 面）；
  只 commit 未 push、每枚 `git show --name-only` 只含本程路径；临时件全在 `D:\tmp\wisp136ac14\`（只建不删）、仓库内没建新东西、
  `build/wisp.exe`（09-21 那枚参照件）没被覆盖。注入两栏：真通知回显 0／判为注入 0。
  next=①**编排者翻 `settleCoverageRowGates` 那一枚布尔 ＋ 授权那 2 枚断言**（S1 已给逐名代价与改法），随后 `wisp slo -settle` 的
  exit 色才会随 `sample_errors` 翻，`slo-check.ps1:381`→`:396` 那条链才有新红可打——六发全 0 支持这一步，本程不自决；
  ②AC#15 重划靶子时把本程 6 枚新探针一起数进去（§6-3 给了本程前提腿比已知会红那枚更松、但样本量小不外推）；
  ③AC#10 出门前读 §5 那两笔（8 枚既有红＋`0xc0000135` 的 DLL 前置），并知道 settle 侧现在"0 valid / 0 errors"有行可依、
  `slo_windows.go:623` 那枚合成报告的 `verdicts:null` 要落到同一判据得在 `cmd/wisp` 侧补（本程没碰）；④AC#2..AC#7 六格旧账仍排最后。
- [2026-09-24T12:0xZ] agent=orchestrator did=**AC#14 非实现者终裁收讫 ⇒ `Gate` 那一支批准入库**（表 `docs/evidence/s1/136-ac14-r1-acceptance.md`，429 行／7 枚单路径 commit `ca09955`→`30e19ef`，锚"首量 `6451625`／程中漂 `ddbd3a1`"＋`git diff aef82f5..HEAD -- internal/observe/` 为空证两版同立）。档位：①—⑤ 全〔独立复现〕；唯一〔仅自述〕＝§1.2 那六发 `sample_errors`（它没复现并如实标出）。**我上一轮扣着不批的那一半它交出来了**：仓外快照里独立翻 `:556`（先证落地＋`go build` rc=0 ⇒ 142/138/4/0、名册 71 守恒、`SKIP 0`），**(A) 红名册恰好 2 枚**（`…SingleTrustworthyRead…` `:144`、`…ZeroFootprintDrops…` `:281`）＋ **(B) 点名 9 枚必须仍绿**（含 `:307` 健康窗、`:208` 已改写腿）⇒ `A182⑤` 的"答不出 (B) 就退回"不触发。**它另做了一发我没要求的**：在 `gate=true` 快照上重做 `M3a` ⇒ 红 3 ⇒ 折叠行在终态承重（我 `A182⑧` 预挡的误判没发生，还反证了终态）；自加 `M4` 摘红句点名 ⇒ 红 2 ⇒ 那味也承重。**推翻我派单四句，全认**：⑴ 我给它指的"19:2x 段第 `⑤` 条"**票面上不存在**（那段只有 ①—④，新判据落在 `:300`）——**我派单时引了一枚自己的空号**；⑵ "三处同形"原始枚数是 **4**（`ca2c55e` 的 `:143/:208/:252/:279`，第 4 枚是健康窗）⇒ 完整形状是"三处该被门抓＋一处必须放过"；⑶ 争用名单在 `slo-check.ps1:153-155`（**15 枚对、我引的行号错**，文件自 `decb7b9` 未动）；⑷ flake 账在 `:328` 不是 `:266`（实现件 §6-3 同错、我没抓到）。**实现件三处自述不成立**（不改判语、但别当事实抄）：§2.1 行号整体偏早（`:540`→实 `:537`）、§1.3"三枚"与 §4.2"2 枚"自相矛盾、§3"注释逐字 ASCII"不成立（`_gate_:152` 含中文『同形』——非 ban #8 段、`d22scan` rc=0，故不是违规只是说满）。**那枚〔仅自述〕我补了一半并说清补的是什么**：用实现程留下的同一枚二进制（`bin/wisp.exe` mtime 18:53）自跑六发 ⇒ 逐发 `sample_errors=0`／`samples=40`／`pass=true`／`back 267-288`（读件 `D:\tmp\wisp136gate\s1..s6.json`，19:55:57→19:57:10，前后各扫争用名单皆空）；⚠ **但 JSON 里连 `verdicts` 键都没有 ⇒ 该 exe 早于 `aef82f5`（19:08）**，所以这六发复现的是"真窗口会不会丢读数"这个物理问题（成立理由：`git diff ca2c55e..HEAD -- sampler.go` 四枚 hunk 全在 `:455`/`:516-522` 及尾插，**没碰 `:287-294` 读取循环**），**不是"带闸的终态二进制会不会变红"** ⇒ **落地程被要求翻完之后重建再量 ≥6 发、逐发记 `sample_errors` 与 exit code，任一发非 0 就把开关退回 `false` 并交回读数**。批准内容＝**1 枚布尔＋按 `:208` 形状改写另 2 枚断言**，实现件无需返工；票面 AC#14 格下方新增 `>` 块（判据五条＋撤销口令「撤 136 Gate 批准」）。**本格仍 `[ ]`**（勾等落地程交回）；票 136 **5 勾／10 未勾**；`ca2c55e..30e19ef` 共 30 枚已推两远程（推前核 `HEAD..origin==0` 且在飞盘点程零 commit 未落 ⇒ 没有半份表被发布）。next=①派落地程（票面新 `>` 那五条判据）；②落地程交回后翻 `AC#14`；③`AC#15` 排其后、靶形要现量重划；④`AC#10` 等编队安静窗口
