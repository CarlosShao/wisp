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

- [ ] **AC#8**（09-23 22:1x 编排者追加，来源＝AC#1 实现方交件里"两条事故"的第一条）`internal/observe` 里有一枚**越界 panic 会把同包全部用例一起吃掉**的测试：
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

- [ ] **AC#9**（09-23 22:3x 编排者追加，来源＝`acceptor-ticket136-r1` §7 `R-136-1`，**它已从断言做成读数**）
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

> **AC#8② 措辞更正（编排者自纠，append-only、原文不抹；09-23 23:3x）**
> 上面 AC#8② 我写的"**只红这一枚、其余全部照常跑完**"这句**按实测不成立，而且不该成立**：M3 那一发变异同时把 `sampling` 那道门开掉
> ⇒ 另有若干枚**本就该红**（实现方实测：改前那发在 panic 之前已有 4 枚红、改后共 6 枚红）。
> ⇒ 本格真正要钉的性质是我那句话的**前半句**：**"不再由一枚 panic 代答——每枚各红各的、且 `RUN` 名册与基线逐名相同"**。
> **终裁请按这条性质判，不要按"只红一枚"判。** 实现方没有回头改我的措辞去对上读数、而是按实测交并报回，这个做法是对的。

- [ ] **AC#11**（09-23 23:3x 编排者追加，来源＝`worker-ticket136-ac8-ac9` 报回的第三条）本包有一枚**既有偶发失败**：
      `TestNoopTaskReturnsToBaseline`（`goroutine_test.go:33`，红在 `PerTask mid-task = 2, want 3`），
      实现方约 **27 发**整包读数里命中 **2 发**（其中一发落在**最终未变异的工作树**），且自锚点起本包只有它自己那两枚文件动过
      ⇒ **不是本程造的**。（⚠ 它先前在票面写的"1/9"已由它自己 append-only 更正成 2/27，见本票 log。）
      ⇒ 它和"吞读数"那一族**同姓**：一条用例偶发红，会把"这一包到底测没测"变成掷硬币，而门禁只看得见"这一次红没红"。
      **结案判据（可重算）**：① 先量复现率并给出 n（同一棵纯净快照连跑 **≥20 发**整包 `-v`，逐名记命中数；**别用单次读数下结论**）；
      ② 根因归到**计数窗口**那一类（goroutine 起落与断言之间缺 happens-before 一类），**不许**归成"机器负载高"这种不可核的说法；
      ③ 修法只许把"等待"变成**有判据的等待**（轮询到条件成立＋超时上界），⚠ **绝对不许**用 `time.Sleep` 把窗口糊过去，
      也不许 `Skip` 它或调高阈值换绿；④ 修完同一发复跑 ≥20 发，命中数必须为 **0**，两次的逐名读数都留在证据里。

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
