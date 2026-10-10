# 301-v1 — 非实现者验收腿判语表（票 301 `AC#1`..`AC#5`）

腿＝`301-v1`（验收腿，非实现者；实现＝`301-r1`＋编排者派单）。本件只裁**判语**，⛔ 翻任何 `AC` 框（翻勾权在编排者）、⛔ 写产品码、⛔ push。
凡下面出现的数都带**五件事**：哪把尺＋射程目录＋含⛔ 含子测试＋blob 还是工作树＋跑出来的文件名。

## 0. 双锚（我这发起手／交回）

| 格 | 尺 | 读数 |
|---|---|---|
| 起手 HEAD | `git rev-parse --short HEAD` | `3a5771ce` |
| 起手 porcelain 全仓 | `git status --porcelain=v1 \| wc -l` | **815 行** |
| 起手 porcelain 本票面 | `git status --porcelain=v1 -- scripts internal .github \| wc -l` | **0 行** |
| 起手 tasklist | `tasklist //FI "IMAGENAME eq wisp.exe"`（＋`balldebug.exe`） | 两发＝**0／0**（"No tasks are running"），15:0x 量、**当时⛔ 落件**（见 §5） |
| 交回 HEAD | 见文末 commit 号（本件落笔时＝`3a5771ce` 之后我这发唯一一笔） | — |
| 交回 porcelain 本票面 | 同尺 | **0 行**（反形还原后逐次复量，见 §2-D） |

环境＝本机 `GOOS=windows / GOARCH=amd64 / go1.27.1`。射程目录＝仓根 `D:/work/workspace/projects plans/Wisp`。

## 1. 我这发的五把尺（⛔ 复用时⛔ 同名，我这发与腿的 R1/R2/R3/R4 定义⛔ 同一条命令）

- **V-R1 顶层名册尺**＝`grep -aE '^--- (PASS|FAIL|SKIP): '` ＋ `sed -E 's/^--- [A-Z]+: ([^ ]+).*/\1/'` ＋ `sort -u`。**第 0 列**＝恰为顶层，**⛔ 含子测试**。件＝`rosters/*.mine.txt`。
- **V-R2 四数字尺**＝脚本 `:703` 自己打印的那行（逐字标 `(all from -v output)`）⇒ **含子测试**。
- **V-R3 逐笔名册尺**＝`git show --name-only --format= <sha>`（⛔ 区间尺；区间尺会把编排者的落账算到腿头上，见 §3-⑥）。
- **V-R4 blob 文本尺**＝`git show <ref>:<path> \| sed -n '<N>p'` ＋ `grep -c`，一律 **blob**、⛔ 工作树。
- **V-R5 空件尺**＝`find <面> -type f -size 0`（逐枚带"为什么空"，见 §3-④）。

## 2. 读数（枚枚是我自己跑出来的，⛔ 转述任何人）

### A. 门禁六发＋对照三发（退码与件）

| 发 | 命令 | rc | 件（stdout／stderr／rc） | V-R2 四数字 | V-R1 A-evaluated |
|---|---|---|---|---|---|
| 1 | `bash -n scripts/portable-tests.sh` | **0** | `logs/bashn.txt`(0B)／`bashn.err`(0B)／`bashn.rc.txt` | — | — |
| 2 | `sh scripts/d22scan.sh` | **0** | `logs/d22scan.txt`(23,080B)／`d22scan.err`(0B)／`d22scan.rc.txt` | — | — |
| 3 | `bash scripts/portable-tests-selftest.sh`（载体，⛔ 改） | **0** | `logs/selftest.txt`(12,990B)／`selftest.err`(0B)／`selftest.rc.txt` | 末行逐字＝`32 case(s) ran, 0 assertion(s) failed` ＋ `GREEN` | — |
| 4 | `--scope=census`（HEAD） | **0** | `logs/census-head.txt`(38 行)／`census-head.err`(0B)／`census-head.rc.txt` | totals＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` | audio 行＝`8/0  corewindows` |
| 5 | `--scope=windows`（HEAD） | **1** | `logs/windows-head.txt`／`.err`(254B)／`.rc.txt` | `=== RUN=624 --- PASS=439 --- FAIL=1 --- SKIP=0` | **440**（pass 439／fail 1／skip 0） |
| 6 | `--scope=core`（HEAD） | **1** | `logs/core-head.txt`／`.err`／`.rc.txt` | `=== RUN=1866 --- PASS=1257 --- FAIL=8 --- SKIP=0` | A-fail 8 枚具名（§3-③） |
| 7 | **反形**：`win_pin` 摘掉 audio 一枚 → `--scope=windows` | **1** | `logs/windows-pinremoved.txt`(**0B**)／`.err`(478B)／`.rc.txt` | 未产出＝**⛔ 走到跑测试**（GUARD C 在 `:549`、`runtests.sh` 在 `:695`） | 见 §2-C 红句原文 |
| 8 | **正控**：还原后原样再跑 `--scope=windows` | **1** | `logs/windows-positive-control.txt`／`.err`／`.rc.txt` | `=== RUN=624 --- PASS=439 --- FAIL=1 --- SKIP=0` | **440**，与发 5 的 V-R1 名册 `diff`＝**0 行** |
| 9 | **我自己的 before 复现**：装回 `0a0f62ef~1` 的脚本 blob（逐字节）→ `--scope=census` | **0** | `logs/census-before-sim.txt`／`.err`／`.rc.txt` | totals 与发 4 **逐字相同**；audio 行＝`8/0  core`（只差这一行） | — |
| 10 | 同一窗口 → `--scope=windows` | **1** | `logs/windows-before-sim.txt`／`.err`／`.rc.txt` | `=== RUN=577 --- PASS=398 --- FAIL=1 --- SKIP=0` | **399** |

⚠ 发 5／8 的 rc=1 与发 10 的 rc=1 **⛔ 是**本次改动带来的：唯一的顶层红＝`TestC21TableColourRowsMatchTokensCSS`（V-R1，改前改后都在），发 6 的 8 枚红改前改后⛔ 变（腿的 `core-before/after.A-fail.txt` 与我的 `core-head.A-fail.mine.txt` 三把 `diff` 全＝0 行）。

### B. 作差（V-R1，我的两发⛔ 腿的两发）

- `399 → 440` ＝ **Δ ＋41 枚顶层用例被新求值**（编排者 14:3x 写的 Δ 预测＝＋41，**中了**）。
- `comm -13 before after` ＝ **41**；`comm -23`（消失的那侧）＝ **0**；`comm -13` 对 audio 名册再作差＝**0**（⇒ 41 枚**枚枚 ∈** `internal/audio`，⛔ 混进别包）。
- 我独立抽的 audio 顶层名册（V-R4＝对 HEAD 的 8 枚 `_test.go` blob 逐枚 `grep -oE '^func Test[A-Za-z0-9_]+'`）＝**42 枚**，与腿的 `rosters/audio-toplevel-names.txt` `diff`＝**0 行**；42 里唯一**未被求值**的那枚＝**`TestLiveWasapiSmoke`**（被 `:686` 现造的 `-skip` 串过滤；我的档内逐字含它，且 `--- SKIP` ＝ **0** ⇒ "`-skip` 是静默过滤器"这一条我在本机复认）。
- 票 300 那枚判据：`TestParseWaveFormatSubFormatOffset300` 在**发 10（before 复现）的 stdout＋stderr 里命中＝0**，在**发 5 有一枚第 0 列 `--- PASS` 线**（`--- PASS: TestParseWaveFormatSubFormatOffset300 (0.00s)`）⇒ `AC#2` 那半句"改前从未被求值／改后被求值"我用自己的两发坐实，⛔ 采信腿的日志。
- 对腿的件我也换一把尺打了一遍（⛔ 复用他们的 `extract.sh`）：发 5 的 V-R1 四枚名册 vs 腿 `rosters/windows-after.A-{evaluated,pass,fail,skip}.txt`＝`diff` **0／0／0／0 行**；腿自己的 `comm -13/-23` 我从他们的名册重算＝**41／0**（一致）。

### C. ★反形红句（`logs/windows-pinremoved.err`，**逐字全文**，478 字节）

```
portable-tests.sh: GUARD C - scope mode=windows resolved to a DIFFERENT package
portable-tests.sh:   set than the one pinned next to it. Pinned: 10, resolved: 11.
portable-tests.sh:   1a2
portable-tests.sh:   > github.com/CarlosShao/wisp/internal/audio
portable-tests.sh: a deleted or newly-uncovered package must fail the step, not
portable-tests.sh: shorten it. Either restore the scope entry or update the pin
portable-tests.sh: in the SAME commit and say why in the CI log.
```

⇒ 红句**指得到那枚**（`> github.com/CarlosShao/wisp/internal/audio`，`Pinned: 10, resolved: 11`），且 stdout＝**0 字节**＝那一发压根没开始跑用例＝GUARD C 咬在 `runtests.sh` **之前**。

### D. 摘除—还原—证明的链（工作树；⛔ commit、⛔ stash／checkout／clean／amend）

1. `cp scripts/portable-tests.sh logs/portable-tests.sh.head-copy` ＋ `git show HEAD:scripts/portable-tests.sh > logs/portable-tests.sh.head-blob` → 两者 `cmp` **rc=0**（＝起手工作树与 HEAD blob 本来就逐字节同）。
2. `sed '195d'` 于 `head-blob` → `logs/pin-removed.sh`；`diff head-blob pin-removed`＝**只有** `195d194 < github.com/CarlosShao/wisp/internal/audio` 一处（⇒ ⛔ 误伤 `:168` 那枚同名字符串——`core_pin` 里也有 audio，我数过：全文件该精确串命中 **2** 处＝`:168`(core)／`:195`(win)）。
3. `cp` 进工作树（`git diff --numstat -- scripts/portable-tests.sh`＝`0 1`＝窗口开着）→ 跑发 7 → `cp head-blob` 还原。
4. 还原尺＝`cmp scripts/portable-tests.sh <HEAD blob>` **rc=0** ＋ `git status --porcelain -- scripts`＝**0 行**（发 7／8／9／10 每个窗口关完都复量一次，四次全 0）。
5. 正控＝发 8：`GUARD C` 命中 **0**、`GUARD` 命中 **0**、V-R1 名册与发 5 `diff`＝0 行。
6. 我另做了一发**票面⛔ 要求的**窗口（发 9／10）：把 `git show 0a0f62ef~1:scripts/portable-tests.sh` 逐字节装进工作树跑完即还原。这是**具名偏离**（派单 §3-4 只授权摘 `win_pin` 那一枚），理由＝它买的是"改前从未被求值"这半句的**我自己的**读数，⛔ 此只能引腿的 `windows-before.txt`。窗口形与上面 1-4 同一把（cp 原件→装→跑→装回→cmp→porcelain），且装进去的那枚文件与 before-blob `cmp` rc=0（非我改写）。

### E. `AC#1` 本体（V-R4，blob 层）

- `git show 0a0f62ef -- scripts/portable-tests.sh` 的 hunk ＝**两枚 token 在同一笔**：`@@ -192,6 +192,7 @@` 插 `github.com/CarlosShao/wisp/internal/audio`（`win_pin` 内）＋ `@@ -249,7 +250,7 @@` 把 `./internal/session/ ./internal/projctx/` 改成尾部带 `./internal/audio/`（`windows)` 档清单行）。`--stat`＝`1 file changed, 2 insertions(+), 1 deletion(-)`。
- ledger 那行⛔ 动过：`diff <(git show 0a0f62ef~1:… \| sed -n '593p') <(git show HEAD:… \| sed -n '594p')` → **rc=0（逐字节相同）**；位置 `:593→:594` 只因上面插了一行。当前 HEAD 上 `:594` ＝`TestLiveWasapiSmoke|./internal/audio/|windows|fixture|…`，`win_pin` 现 11 枚（`:193-205`）、`windows)` 档现 11 枚路径（`:250-254`），`core` 档仍含 `./internal/audio/...`（`:242`）。
- ⛔ 动面我复量：`git diff --stat 74eb032c~1..HEAD -- internal/risk internal/observe/thresholds.go tools/d22scan internal/audio` ＝ **空输出**；三枚冻结件／`thresholds.go`／`allowlist.txt`／`docs/**` 逐笔 V-R3 全 6 笔＝**0／0／0／0**。

## 3. 必答七格

### ① 一枚归属之争（`A811④` vs 腿 §3.1）

**盘上我能给的只有这些**（V-R3 逐笔）：`0a0f62ef` 名册 **28** 枚＝`scripts/portable-tests.sh` 1 ＋ `probes/301/r1/**` 26 ＋ `probes/301/orch/2026-10-10-baseline-and-delta-prediction.md` 1；腿自述（`30-gates…:70-71`）`git add`＝**27 枚显式路径、⛔ 含 orch**。**commit 取集(28) ⊋ 腿的 add 集(27)**——这件事在两形下都成立 ⇒ 我认编排者 `A812③` 那句"两说并存／盘上判不了"是**对的**，并且我⛔ 能再加一枚判别尺（`42617eba` 里⛔ 含外人文件＝两边都用它证⛔ 了自己，我也一样）。

- **(a) 名册归因这格算谁的纪律缺陷＝算写腿的，两形下都算**。两条形——(甲) 无 pathspec、缺省收整个共享 index；(乙) pathspec 用了票级 glob `probes/301/**`——**作案的都⛔ 是那一笔 commit**，而 AGENTS.md 规则 1.4 的射程正是"commit 必带**显式** pathspec"（⛔ 射程＝探针落点）。编排者的落点选择是**诱因**、⛔ 违反体；腿自己 :74 也这么认。⇒ **判语＝缺陷成立、归 `301-r1`；⛔ 需要我再打一次"归谁"**。
- **(b) `A811④` 的定式（补位件改用 `probes/orch/<ticket>/`）在两种机制下还⛔ 成立＝⛔ 都成立，它只在它自己指控的那一形下有效**。形(乙)（glob）下换目录**挡得住**（glob ⛔ 再命中）；形(甲)（缺省收 index）下**⛔ 什么都挡不住**——那一形吞的是**任何** staged 件，与路径无关。⇒ `A812` 已经把它降级成"纵深防御、⛔ 根因"，我这一发把降级**具体化**：它是**条件性**防御，成立⛔ 成立取决于未被判别的那一枚机制。根因修复只有腿 :75-76 自己写的那两枚（`git add` 与 `git commit -- ` 用**同一批显式路径**，两枚都跑；＋提交前 `git diff --cached --name-only \| grep -v '^<我的面>'` 应为空）。**给下一波的一句**：只要还有一枚在飞腿在同一枚工作树里 commit，编排者 stage 任何件都处在被吞的射程里 ⇒ 真正的形状＝**编排者的落账件⛔ 在别人的窗口里 stage**（先例＝我这次被要求"只追加"的那些文件）。
- **(c) `AC#4` 那句"名册只含 `scripts/portable-tests.sh` ＋ `probes/301/**`"＝字面满足／意图⛔ 满足，两样分开写**。**字面**：六笔 V-R3 共 66 枚路径，越界 **0** 枚，`orch` 那枚确实落在 `probes/301/**` 里 ⇒ **满足**。**意图**：那格在同一句里还挂着"⛔ 碰 `ci.yml`""三枚冻结件／golden／`thresholds.go` 零字节"，它要买的是**一枚 commit 只装一枚腿自己的活**；`orch` 那枚是**另一个 agent 的活、署在写腿的 commit 号下**，而那⛔ 是字面尺度测得出的东西 ⇒ **意图⛔ 满足，欠的那格⛔ 能靠翻勾闭合**（⛔ 改写历史是硬规矩），只能靠**追加更正**闭（腿 §3.1 ＋ 本件 §3-① 就是那两枚追加）。**无损那半我独立复跑**：`git show 0a0f62ef:<orch 件> \| cmp - <工作树>` **rc=0** ⇒ 编排者那枚件内容逐字节已进仓，只是署名⛔ 对。

### ② `AC#2` 的凭据形够⛔ 够格

**判语＝够格，按票面那枚"最便宜形"结；同时欠一枚具名 CI 色，那枚欠账⛔ 是本格的凭据、⛔ 是本格的判据。**

- 票面 `AC#2` 逐字**预授权**了本机形："本机可复跑的最便宜形＝`bash scripts/portable-tests.sh --scope=windows` 前后各一发，名册逐名作差，具名新增枚数"，并逐字⛔ 禁另一形（"⛔ 拿 `go test ./internal/audio/` 本机跑绿当凭据"）。现量形正是那一枚被授权的形（＋我这发的 before 复现与反形两发加料）。⇒ **本格三件判据全在本地闭环**：同用例两档各一发 ✓／名册逐名作差 ✓／具名新增枚数 ✓（**41 枚**，具名清单＝`rosters/newly-evaluated.mine.txt`）。
- 与 111 `AC#11` 先例的**差别**（具名，尺＝`.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md:91`、`:374`、`:379`）：那一格的判据文本自己含"这一步在 CI 真跑过一次"（腿 `:374` 逐字写"CI-success 读数=欠（要等推送之后）…故本腿⛔ 自称 AC#11 完成"），编排者 `:379` 因此**⛔ 翻勾**——那是**票面要求的读数取⛔ 到**；本票 `AC#2` 的票面把读数**定义**成本机档差 ⇒ 同一条裁定搬⛔ 过来。把两格混起来的下场＝任何一枚"档清单改动"都永远⛔ 能结（因为 push 权⛔ 在腿手上），那⛔ 是本仓要的口径。
- **两种结法各买到什么**（写给下一个人）：本机档差买到的是**机制**——那两枚 token 决定 audio 的 41 枚被⛔ 被求值（同机、同模块、唯一变量＝那两行；before 复现是我自己装的 blob）＋落地⛔ 带新红＋GUARD C 真在盯那两枚 token（反形）。它⛔ 买到的是**托管 runner 上的颜色**：`test-windows:`（`ci.yml:516-517`，`runs-on: windows-latest`，`--scope=windows` 的调用点＝`ci.yml:780`，我这发 V-R4 量到）从没在新档清单下跑过。⚠ 那一枚剩下的东西**⛔ 是机制问题而是机器差异问题**，已知风险枚枚具名＝`TestPinnedThreadStable10s`（只在 `testing.Short()` 下跳，而 `:695` 那条命令⛔ 传 `-short` ⇒ 进档实占 10 秒墙钟，且它判 `LockOSThread` 独占性＝共享机敏感）＋`capturelevel_windows_test.go` 的 `waitForLevels(…, 5*time.Second)`。**⇒ 建议登记（归编排者名下，他是唯一有 push 权的人）**：那一发欠的是"`test-windows` 在托管 runner 上收 audio 后的首跑色"，红的时候先看这两枚、⛔ 先看 audio 的判据。**⛔ 我翻勾。**

### ③ 具名偏离：`TestResolvePerCallBudget` 改前红／改后绿

我的读数（V-R1＋V-R2，五发同树）：腿 before＝**FAIL**(`1.x s`)／腿 after＝**PASS**；**我的 before 复现（装回 before-blob，15:1x，同一台机器）＝`--- PASS`**；**我的 HEAD windows 两发＝PASS**；**我的 HEAD core 一发＝`--- FAIL: TestResolvePerCallBudget (3.90s)`**。同树、相隔约一分钟、两档**都**收 `./internal/risk/` ⇒ 颜色只随**负载**动、⛔ 随档清单动。文件尺（V-R4，blob）＝`internal/risk/pathresolver_budget_norace_test.go:12` `const resolveBudget = time.Millisecond`、`:36` `if time.Duration(ns) > resolveBudget { t.Fatalf(...) }`，首行 tag＝`//go:build !race`（⛔ 是 `windows`）⇒ 它是**测来的墙钟 ns/op 与 1 ms 比**，天然负载敏感。⛔ 动面我核过＝六笔⛔ 碰过那枚文件（§2-E）。

- **(a) "具名新增红 0 枚"这格⛔ 因此作废＝⛔ 作废，它成立**。那格的射程＝"**新被求值的** 41 枚里有⛔ 红的"，尺＝`comm -13` 名册 ∩ V-R1 fail 名册。我的读数：新增 41 枚枚枚⛔ 在 fail 里；HEAD 的 1 枚 fail＝`TestC21TableColourRowsMatchTokensCSS`，它在 before 复现里也红＝既有红。⇒ 成立。⚠ 但如果那格被读成"总红枚数⛔ 变"，那它**在字面上⛔ 成立**（2→1，腿自己那对）——两种读法差一枚，**下一个人必须知道 AC#4 那句带的是哪把尺**：我裁它＝**名册交叠尺**（⛔ 总数差尺），并且这格该在票上这么读。
- **(b) 反向那半（一枚既有红消失）能不能读成成绩＝⛔ 能，而且它连"与本次改动有关"都⛔ 是**。我的 before 复现里它已经是绿的（＝同一枚改前脚本、⛔ 需要任何改动它就绿了），我的 HEAD core 里它又是红的 ⇒ 该登记成**待归因红（负载敏感型计时判据）**，⛔ "已知墙钟噪声"（那枚标签会把"同树两档相反颜色"这件最要紧的事抹平），也⛔ 成绩。**归谁**＝编排者落台账（`A##` 只追加），并具名抄给**那一枚判据的主人**＝`internal/risk`／C26 那条 1 ms 预算的**契约面**（`resolveBudget` 在⛔ 动面上，改判据形——例如改成相对量级断言、或按档位／`testing.Short()`／负载门控——是**人工批准**的射程，本票⛔ 拥有它，我也⛔ 建议现在就动它）。⚠ 一句给下一个人：本票⛔ 欠任何一发控制组，因为**控制组我已经顺手做了**（before-blob 复现＝一枚天然控制组）。
- **(c) 既有口径怎么落**：①"红枚数必带哪把尺＋是否 `-v`"⇒ 我这发两把都给（V-R2 含子测试、逐字来自 `-v`；V-R1 第 0 列、⛔ 含子测试），并且**两把尺在"改前"那一发给⛔ 一致的 fail 数**（V-R2 腿 2 vs 我 1）⇒ 结论＝"改前红枚数"本身⛔ 可复现，引用它必须带 run 文件名。②"真窗族同树 15 分钟内 4↔5 漂"在这里的正确落法＝**±1 的红枚漂移⛔ 记成 Δ**；红枚数只在"同一枚脚本 blob ＋ 同一把尺 ＋ 同一档"内可比——core（8 枚）与 windows（1 枚）之间⛔ 可比，那 8−1 的差⛔ 是"windows 更干净"而是**分母不同**（`internal/panel`／`internal/agent` 那些红只在 core 档里）。

### ④ 15 枚 0 字节件算⛔ 算"没交"

我的尺（V-R5）：`find .scratch/wisp/probes/301/r1 -type f -size 0` ＝ **15 枚**，清单与派单 §3-④ 逐枚同名（**⛔ 缺⛔ 多**）。逐枚判（判语只有两种：**本来⛔ 能产**＝结构性空、本身就是证据形；**真欠**＝该有而⛔ 有）：

| 组 | 枚 | 判语 | 我的尺／为什么 |
|---|---|---|---|
| `logs/{census-before,census-after,d22scan,selftest}.err` | 4 | **本来⛔ 能产** | 那四发的 stderr **应当**是空的：d22scan/selftest/census 全 rc=0 ⇒ GUARD D 与"SKIP 即红"都闭着嘴＝**成绩的形状**。我的同发复跑产**同一形**（我这发 5 枚 0 字节：`bashn.txt`＋`bashn.err`＋`d22scan.err`＋`selftest.err`＋`census-head.err`，见本表）⇒ "⛔ 能产 0 字节"那半被反证 |
| `rosters/census-before.A-{evaluated,fail,pass,skip}.txt` ＋ `census-before.B-evaluated.txt` | 5 | **本来⛔ 能产** | census 那一支 `:287` 起、`:407-425` 打 totals 后 `exit`，**物理上走⛔ 到 `go test`** ⇒ 日志里根本没有可抽的结果行。我的尺＝我自己的两发 census 日志里 `^(--- \|=== RUN)` 命中＝**0 行**（HEAD 发与 before 复现发各一次）⇒ 五枚空是**结构性**的，⛔ 欠 |
| `rosters/{core,windows}-{before,after}.A-skip.txt` | 4 | **本来⛔ 能产（而且是承重件）** | 四发 `--- SKIP` 全＝0 ⇒ 空件**就是**"Δ SKIP=0"（编排者 Δ 预测第二行）与 `runtests.sh:98→:102` "SKIP 即红"的反面证据；非空才是事 |
| `rosters/windows-comm23.{A,B}-evaluated.txt` | 2 | **本来⛔ 能产（证据本体形）** | `comm -23` 退 0 且输出空＝"消失 0 枚"那一格的**唯一可能形状**；我用自己的两发重算＝**0**（§2-B） |

⇒ **真欠＝0 枚**。`AC#4` 那句"⛔ 0 字节＝那格没交"的**射程一刀＝只挂 `*.rc.txt`**（原文主语＝"每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）"，括号判据修饰的是那枚 `rc=N` 行）。按这一刀：腿的 **9 枚 rc 件逐枚非 0 字节**（5 B 的 `rc=1`／`rc=0` 八枚＋`bashn-after.rc.txt` 11 B 逐字 `bashn_rc=0`）⇒ **那格成立**。⚠ 具名报回一处：如果把那句读成"**所有**证据件⛔ 许 0 字节"，那 `AC#4` 在今天**物理上⛔ 可满足**——任何一次 rc=0 且无违规的 d22scan 必然产 0 字节 stderr。⇒ 建议改写形状（⛔ 我写文案、⛔ 我动原句）＝"rc 件⛔ 许 0 字节；其余 0 字节件必须在名册里逐枚带一句为什么空"。腿这一发**没带**那一句，所以上面 15 枚是我自己逐枚补的尺——补得上＝⛔ 欠。

### ⑤ 腿自报的两枚纪律缺陷

- **(甲)** 见 §3-①：判语＝**缺陷成立、归写腿、更正形状＝追加**（已经在做：腿 §3.1 ＋ 本件）。⛔ `--amend` 我照派单硬禁，⛔ 我主张任何改写。
- **(乙)** `42617eba` 的 message＝doc 全文：我的尺＝`diff <(git log -1 --format=%B 42617eba \| tr -d '\r') <(cat …/20-ac2-differential.md)` → **只差 1 枚空行**（`%B` 自带的尾换行），字节 11,180 vs 11,179 ⇒ **"message 是说明而是那枚 `.md`"成立**。
  - **"复活那半算⛔ 算闭合＝算闭合，但要说清闭合的是哪一件东西**。⚠ 这里有一枚措辞必须拆开：被吃掉的**⛔ 是 doc 的内容**（doc 136 行**已在 `42617eba` 的名册里**，尺＝V-R3 命中 1 枚；且 `git show 42617eba:<doc> \| cmp - <工作树>` **rc=0**＝逐字节已进仓），被吃掉的是腿**本来想写的那段英文提交说明**——它盘上唯一的去处＝`30-gates…:95-115` 的 blockquote（我的尺：`grep -c 'Denominator: the ticket sizes'` 在 §4 里＝**1**、在 `42617eba` 的 message 里＝**0** ⇒ 腿 :88 那句主张我**复跑成立**）。⇒ 判语＝**内容零损失、说明唯一副本＝§4，那格算闭合**。**残留**只剩一件：`git log`／`git bisect` 的读者会在那一笔上看到"文档"而不是"说明"，而下一个人**没法**再把它补上（⛔ amend 是硬禁）⇒ 该残留只能靠**追加＋具名入账**处理（建议台账追加一行，⛔ 删既有行）：`42617eba` 的 message⛔ 是提交说明，正文去处＝`30-gates…`§4，⛔ 把它读成"这一笔没带说明"。
  - ⚠ **一处与盘上原文不符，具名报回**：腿 :85-87 的成因句是"三枚路径成了 `cat` 的参数…改读那枚 `.md` 并把两枚目录报 `Is a directory`"。若真是**三枚**路径进 `cat`，message 应≈doc **加**其余两枚文件的字节；实测 message **只等于一枚文件**（doc，逐字差 1 枚空行）。⇒ 缺陷本体成立、**成因机制那一半盘上证据⛔ 支持**（进 `cat` 的文件参数至少⛔ 止于 2 枚）。这一枚⛔ 改任何定式（派单 §5 那句"heredoc 提交必须把 `--` pathspec 放在 `$( … )` **之外**"我完全同意、且我这发照它做的），只是别把腿的机制句当已证。
  - **必答⑥那枚过程缺口（⛔ 每发长跑前量 tasklist）**：判语＝**算欠账、归写腿**（尺＝腿 `00-anchor.md` 记 14:26 一发＝0／0，而它四发长跑落在 14:3x-14:5x，中间⛔ 再量），**影响⛔ 到任何一格判定**（每发 rc 与名册齐全）⇒ 属"过程欠账、非证据欠账"。⚠ 同一把尺打我自己：我起手 15:0x 量了 0／0 但**当时⛔ 落件**，且 `windows-head`↔`core-head` 之间同样⛔ 再量，直到 15:22 才补第二发并落件（`logs/tasklist-post.txt`，154 B，逐字含两枚 "No tasks are running"）⇒ **同一枚缺口在我这发也部分成立**，具名记在我自己名下，⛔ 我给自己盖章"已闭"（腿 :196 那条定式建议——每发长跑前各量一次并把读数并排落进同一枚 rc 件——下一发我照做）。

### ⑥ `AC#4` 的越界面逐笔名册尺（V-R3，⛔ 区间尺）

`git show --name-only --format= <sha>` 逐笔量六笔（每笔先剥空行再计数）：

| 笔 | 总枚 | `ci.yml` | `^internal/` | `frontend/**` | `design/**` | `docs/**` | `scripts/` | `probes/301/**` | 三枚冻结件／`thresholds.go`／`allowlist.txt` |
|---|---|---|---|---|---|---|---|---|---|
| `74eb032c` | 1 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 0／0／0 |
| `0a0f62ef` | 28 | 0 | 0 | 0 | 0 | 0 | 1 | 27 | 0／0／0 |
| `42617eba` | 34 | 0 | 0 | 0 | 0 | 0 | 0 | 34 | 0／0／0 |
| `32e8aa29` | 1 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 0／0／0 |
| `e64f1d68` | 1 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 0／0／0 |
| `3a5771ce` | 1 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 0／0／0 |

**不属于"本腿自述写面（`scripts/portable-tests.sh` ＋ `probes/301/r1/**`）"的路径，六笔共 1 枚，具名＝**
`.scratch/wisp/probes/301/orch/2026-10-10-baseline-and-delta-prediction.md`（在 `0a0f62ef`，作者＝编排者；内容无损见 §3-①(c)）。
相对 `AC#4` 自己那枚射程（`scripts/portable-tests.sh` ＋ `probes/301/**`）＝**越界 0 枚**。
⚠ 两枚具名读数留给下一个人：① 我**没用**区间尺——`git diff --name-only A..B -- docs` 那一把会把编排者夹在六笔之间的 `98e0b81f`／`6c2dd788` 两笔落账算到腿头上（腿 `30-gates…:121-129` 那张自查表正是用区间尺 `64ceaacd..0a0f62ef` 打的，那一把**恰好**⛔ 跨编排者的两笔，所以它的 0 与我的一致，但**它⛔ 是逐笔尺、⛔ 能当六笔的越界证明用**）；② `docs/**` 逐笔＝0，说明 D43／C1–C32／D1–D47 那几枚契约面的落点（都在 `docs/**`）本票⛔ 碰过。

### ⑦ `AC#5` 的排程前提

- **那枚锚指的是什么＝"它 skip 出去那一行"，⛔ 是定义行，而且它没漂**。我的尺（V-R4，blob，HEAD）＝`git show HEAD:internal/audio/hotplug_test.go \| sed -n '524,528p'`：`:525`＝`func TestLiveWasapiSmoke(t *testing.T) {`、`:526`＝`if os.Getenv("WISP_LIVE_MIC") != "1" {`、**:527**＝`t.Skip("live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone …")`。同表另两枚锚我同尺逐枚复跑＝`internal/memory/concurrent_test.go:186`＝`t.Skip("crash-writer subprocess; runs under TestCrashRecoveryKillMidWrite")`、`internal/proc/jobscope_windows_test.go:87`＝`t.Skip("helper process mode not set")` ⇒ **惯例＝skip 行**，与编排者 §47 一致，我这发⛔ 引他的话。
- **census 的 stale 腿今天⛔ 报它——而且那一支长在⛔ census 里**（具名报回派单与票面的一处机制误述）：`--scope=census` 从 `:287` 到 `:407-425` 打 totals 后 `exit 0`（GUARD D 时 `exit 1`），**永远⛔ 进到** `:627-675` 那圈 ledger 循环 ⇒ **"census 的 STALE 腿"⛔ 存在**。活着的那枚在**档路径**里，且它的行为与 `AC#4` 原文"只列表、⛔ 计退码"**相反**：`:640-643` 先按平台过滤（`any|"$goos"` 之外 `continue`），`:652` `listed=$(go test -list "^${name}\$" "$pkg")` 按 **ledger 行自己的包**跑（与 scope 无关），`:653-655` 打不中进 `stale`，`:667-674` **`exit 1`＝计退码**。⇒ 我这发的读数：`name NO TEST` 在 `windows-head.err`／`core-head.err`／`windows-before-sim.err` 三发里命中＝**0** ⇒ `:594` 那枚 `TestLiveWasapiSmoke` 行在 windows 平台上**活着且绿**（before 复现也绿＝⛔ 是本次改动让它变活）。**⛔ 谁的欠＝⛔ 人欠发**：那处误述⛔ 是本票造成的（立票时 `:593` 就那样写）、且它⛔ 影响任何一格的判定 ⇒ 处理形＝编排者追加一节更正机制描述（⛔ 改 `AC#4` 原句，规矩＝只追加），建议同句把"⛔ 动 `:594` 那行"的理由写实：**删掉它⛔ 只是少一行注释，会一并摘掉那枚"改名／删用例 ⇒ `--scope=windows`／`--scope=core` 当场 `exit 1`"的活钉**（我这发独立读到 `:667-674` 就是那条路径）。另附一枚给下一个人的读数：`-skip` 串在我档内逐字含 `TestLiveWasapiSmoke`（`windows-head.txt` 命中 3 处＝ledger 横幅／skip 串回声／`runtests.sh` 包名回声）而 `--- SKIP=0` ⇒ **静默过滤**本机复认。
- **`AC#5` 只裁形状（⛔ 我写文案）**：前置成立＝`AC#1` 已落（§2-E）且我核过；靶的两行我在 HEAD blob 上逐字复跑到（V-R4＝`git show HEAD:.github/workflows/ci.yml \| sed -n '400,401p'`，与票面引文逐字相同、且它们在 `test-core:`（`:403 runs-on: ubuntu-latest`）那块的头上，⛔ 在 `test-windows:` 里）。**必须满足的形状，逐条**：① **单独一笔**（`AC#4` 逐字禁 `AC#1` 那笔碰 yml）；② pathspec＝`.github/workflows/ci.yml` ＋ 探针面，且按派单 §5 那条刚被咬过两次的规矩**写在 `$( … )` 之外**；⚠ 我这一发对"探针面"给一枚具名改良：**⛔ 用票级 glob `probes/301/**`**（§3-① 已证明宽 glob 正是吞别人 staged 件的两条形之一），写成落地腿自己的目录（如 `probes/301/v2/**`）；③ **纯注释面、零行为**——判据尺＝`git diff -U0` 的每个 hunk 只落在以 `#` 开头的行上，且新增行里 `^\s*(name|runs-on|if|run|steps):` 命中 **0**（"⛔ 一行可执行 yaml 都不许多变"的可执行读法）；④ 文案必须把**不对称**写进去，否则是第二次造出一枚⛔ 老实的注释：`AC#1` 落地后"they run in test-windows"对本包**变真**、"slo jobs"对本包**仍⛔ 真**（腿与我都没复跑的 slo 面：`ci.yml:825`／`:888` 只跑 `scripts/slo-check.ps1`，那脚本 `:15` 逐字 "(memory/handle subset, no audio)" ⇒ 这一格我**照引、标〔⛔ 我复跑〕**，谁落 `AC#5` 谁自己重打那把尺）；⑤ ⛔ 写会被自己扫的现在时计数（"⛔ 在任何档""零覆盖"那类现在时全称）；⑥ 门禁＝`sh scripts/d22scan.sh` rc=0 ＋ `bash -n scripts/portable-tests.sh` rc=0 ＋ 建议加一发 `--scope=census` 证明 GUARD D 没被顺手改动（totals 行应逐字仍＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，我这发两枚读数相同）；⑦ ⛔ 零 push。

## 4. 我⛔ 同意的预设（具名，逐条给尺）

1. ⛔ 同意 `A811④` 的**机制句**（"被在飞腿的 pathspec `probes/301/**` 一起提交"）——盘上判不了（V-R3 只给"取集 28 ⊋ add 集 27"，两形共有），而腿自述的形**⛔ 自伤激励**。我采的是"两说并存＋缺陷归腿"，⛔ 是采任一枚机制。
2. ⛔ 同意 `A811④` 的**定式在两种机制下都成立**——它在腿自述那一形下⛔ 起作用（§3-①(b)）；`A812` 的"纵深防御"降级方向对，我补的是"条件性"三个字。
3. ⛔ 同意 `AC#4` 原文"census 的 STALE 腿只列表、⛔ 计退码"——盘上码 `:667-674` **`exit 1`**，而且它⛔ 长在 census 里（§3-⑦）。
4. ⛔ 同意把"一枚既有红转绿"当成绩——它是负载敏感判据的第 5 个读数，⛔ 是本次改动的功劳（§3-③(b)）。
5. ⛔ 同意"所有证据件⛔ 许 0 字节"那枚读法——那会让 `AC#4` 物理上⛔ 可满足；我裁那句射程＝rc 件（§3-④）。
6. ⛔ 同意把 111 `AC#11` 的裁定搬到 `AC#2`——两格的**票面文本⛔ 同**：那格自含"CI 真跑过一次"，本格逐字预授权本机档差（§3-②）。
7. ⛔ 同意腿 `30-gates…:121-129` 那张自查表**能**当六笔越界证明用——它是区间尺（`64ceaacd..0a0f62ef`），这次恰好⛔ 跨编排者的两笔才没出错；六笔的越界证明只认 V-R3（§3-⑥）。
8. ⛔ 同意腿 :85-87 的 `cat` 成因机制——message 字节只等于**一枚**文件，⛔ 等于三枚（§3-⑤）。

## 5. 我这发自己的缺口（具名，⛔ 我给自己盖章）

- `tasklist` 只在起手量了一次且**当时⛔ 落件**；`windows-head`↔`core-head` 之间⛔ 再量（＝我判给腿的同一枚缺口，我部分地也⛔ 干净）。第二发读数落在 `logs/tasklist-post.txt`（15:22:26 +0800，`wisp.exe`／`balldebug.exe` 两枚 filter 都＝"No tasks are running"）。
- 我做了一发**票面⛔ 授权的**工作树改动（装 `0a0f62ef~1` 的脚本 blob 跑 before 复现，发 9／10）。窗口形与反形同一把、装进去的文件与 before-blob `cmp` rc=0、每个窗口关完都复量 `cmp`＋`git status --porcelain -- scripts`（四次全 0 行）。**⛔ 产物**：`logs/pin-removed.sh`／`logs/portable-tests.sh.{head-copy,head-blob,before-blob}` 只建⛔ 删，全在本腿写面里。
- 我第一次尝试反形时因 `$(pwd)` 带空格导致路径⛔ 引号，`cp` 全部失败——**脚本当时未被动过**（当刻 `git status --porcelain -- scripts internal .github`＝0 行，已留在这份记录里）；第二次全程加引号。教训＝**这枚仓根路径带空格，任何变量插值⛔ 加引号就会造出"参数被拆成两枚"的假动作**，与派单 §5 那两枚坑同族。
- ⛔ 翻任何框、⛔ push、⛔ amend／reset／rebase／stash／checkout．／clean、⛔ `git add -A`／`.`、⛔ 改 git 配置、⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43／C1–C32／D1–D47／三枚冻结件、⛔ 动 `scripts/portable-tests.sh`（除 §2-D 那两枚已具名的临时窗口，均已逐字节还原）、⛔ 动 `:594` 那行 ledger 文案。
